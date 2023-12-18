@echo off
cls
SET PATH=%PATH%;.\assets\binaries
echo.

echo %cd%

:: Check if that output dir exists
if not exist "output" mkdir "output"

:: Cleanup any previous output files
del /F /Q "output\divekit-windows-386.msi" 2> nul
del /F /Q "output\divekit-windows-amd64.msi" 2> nul

echo ---------------------------- Cloning repositories ----------------------------
git clone https://github.com/divekit/divekit-automated-repo-setup ./repositories/divekit-automated-repo-setup
git clone https://github.com/divekit/divekit-repo-editor ./repositories/divekit-repo-editor

:: Create the executable and verify its creation
go build -C ../../ -o ./scripts/installer/divekit.exe
if not exist "divekit.exe" echo Could not build "divekit.exe" & exit /B

echo ------------------------- Checking Divekit Version --------------------------
:: NOTE: divekit.exe version is currently not supported; instead, the version is temporarily simulated.
:: How it could be called: divekit.exe version > out

:: Get the first line of the command `divekit version`
echo divekit version 1.0.0> out

:: Set a variable called DIVEKIT_VERSION with the current divekit version
for /f "tokens=*" %%i in (out) do set DIVEKIT_VERSION=%%i
set DIVEKIT_VERSION=%DIVEKIT_VERSION:divekit version =%
echo Divekit version is '%DIVEKIT_VERSION%'

:: Substitute every instance of DIVEKIT_VERSION in the configuration file with the current version.
powershell -Command "(gc ./config/divekit-windows-amd64.wxs) -replace 'DIVEKIT_VERSION', '%DIVEKIT_VERSION%' | Out-File -encoding ASCII divekit-windows-amd64-proc.wxs"
echo.

echo ---------------------------- Creating x64 MSI ----------------------------
:: Purpose: heat.exe is a tool in the WiX Toolset that is used to automatically generate WiX authoring (XML code)
:: for a directory tree. It harvests the files and components from a source directory and creates WiX authoring
:: that can be used in a WiX source file (.wxs). In this scenario, the cloned repositories are created as WiX
:: authoring and are then referenced in the config file (divekit-windows-amd64.wxs).
heat dir ./repositories/divekit-automated-repo-setup -scom -frag -srd -sreg -ke -gg ^
    -o ./ars-dir.wxs -cg CMP_ARS_DIR -dr ARS_DIR
heat dir ./repositories/divekit-repo-editor  -scom -frag -srd -sreg -ke -gg ^
    -o ./repo-editor-dir.wxs -cg CMP_REPO_EDITOR_DIR -dr REPO_EDITOR_DIR

:: Purpose: candle.exe is the WiX compiler. It takes WiX source files (.wxs) as input and compiles them into
:: intermediate object files (.wixobj). This step verifies the correctness of the WiX source and
:: prepares it for linking.
candle divekit-windows-amd64-proc.wxs > nul
candle ars-dir.wxs -arch x64 > nul
candle repo-editor-dir.wxs -arch x64 > nul

:: Purpose: light.exe is the WiX linker. It takes the intermediate object files (.wixobj) produced by candle
:: and combines them into a Windows Installer package (.msi). The linker resolves references between components,
:: features, and other elements, creating a complete and executable installer.
light -ext WixUIExtension -cultures:en-us divekit-windows-amd64-proc.wixobj -o divekit-windows-amd64-proc.msi ^
    ars-dir.wixobj -b ./repositories/divekit-automated-repo-setup ^
    repo-editor-dir.wixobj -b ./repositories/divekit-repo-editor

:: echo ---------------------------- Creating x86 MSI ----------------------------
:: Support x86 later

:: Rename and move the msi file to output
ren divekit-windows-amd64-proc.msi divekit-windows-amd64.msi > nul
move divekit-windows-amd64.msi output > nul

:: Cleanup build artifacts
del /F /Q divekit-windows-amd64-proc.wxs 2> nul
del /F /Q divekit-windows-amd64-proc.wixpdb 2> nul
del /F /Q divekit-windows-amd64-proc.wixobj 2> nul
del /F /Q ars-dir.wxs 2> nul
del /F /Q ars-dir.wixobj 2> nul
del /F /Q repo-editor-dir.wxs 2> nul
del /F /Q repo-editor-dir.wixobj 2> nul
del /F /Q "divekit.exe" 2> nul
del "out" 2> nul
powershell -Command "rm ./repositories -r -force"
