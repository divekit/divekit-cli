#Requires -RunAsAdministrator

function commandExists([string]$cmdName) {
    return [bool](get-command -Name $cmdName -ErrorAction SilentlyContinue)
}

function install([string]$cmdName, [string]$packageId, [string]$packageManager, [bool]$shouldUpgrade) {
    if (-not($shouldUpgrade) -and (commandExists($cmdName))) {
        write-output "$packageId is already installed"
        return
    }

    if ($packageManager -eq "winget") {
        winget install --id $packageId --accept-source-agreements --accept-package-agreements
        return
    }
    choco install $packageId -y
}

function initialize() {
    # install winget if not installed (https://winget.pro/winget-install-powershell/)
    if (-not(commandExists("winget"))) {
        # prevent "microsoft.ui.xaml.2.7" could not be found
        write-output "Downloading microsoft.ui.xaml.2.7 ..."
        $ProgressPreference = 'SilentlyContinue'
        invoke-webrequest `
            -uri https://www.nuget.org/api/v2/package/microsoft.ui.xaml/2.7.3 `
            -outfile xaml.zip -useBasicParsing
        $ProgressPreference = 'Continue'
        new-item -itemtype directory -path xaml
        expand-archive -path xaml.zip -destinationpath xaml
        write-output "Installing microsoft.ui.xaml.2.7 ..."
        add-appxpackage -path "xaml\tools\appx\x64\release\microsoft.ui.xaml.2.7.appx"
        remove-item xaml.zip
        remove-item xaml -recurse

        # prevent "microsoft.vclibs.140.00.uwpdesktop" could not be found
        write-output "Downloading microsoft.vclibs.140.00.uwpdesktop ..."
        $ProgressPreference = 'SilentlyContinue'
        invoke-webrequest `
            -uri https://aka.ms/microsoft.vclibs.x64.14.00.desktop.appx `
            -outfile uwpdesktop.appx -useBasicParsing
        $ProgressPreference = 'Continue'
        write-output "Installing microsoft.vclibs.140.00.uwpdesktop ..."
        add-appxpackage uwpdesktop.appx
        remove-item uwpdesktop.appx

        # get the download URL of the latest winget installer from GitHub
        $ProgressPreference = 'SilentlyContinue'
        $api_url = "https://api.github.com/repos/microsoft/winget-cli/releases/latest"
        $download_url = $(invoke-restmethod $api_url).assets.browser_download_url |
                where-object {$_.endsWith(".msixbundle")}
        # download the installer
        write-output "Downloading winget ..."
        invoke-webrequest -uri $download_url -outFile winget.msixbundle -useBasicParsing
        $ProgressPreference = 'Continue'

        # install winget
        write-output "Installing winget ..."
        add-appxpackage winget.msixbundle

        # remove the installer
        remove-item winget.msixbundle

        # refresh $env:Path to use winget
        $Env:Path = [System.Environment]::GetEnvironmentVariable("Path", "Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path", "User")
    }
    # install choco if not installed
    if (-not(commandExists("choco"))) {
        install "choco" "Chocolatey.Chocolatey" "winget" $false
        # refresh $env:Path to use choco if not installed before
        $Env:Path = [System.Environment]::GetEnvironmentVariable("Path", "Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path", "User")
    }
}

initialize

# Upgrade if already installed
[bool]$shouldUpgrade = $false

# install [cmdName] [packageId] [packageManager] [shouldUpgrade]
install "git" "Git.Git" "winget" $shouldUpgrade

install "go" "GoLang.Go" "winget" $shouldUpgrade

install "java" "Oracle.JavaRuntimeEnvironment" "winget" $shouldUpgrade # umlet requires java runtime

install "node" "OpenJS.NodeJS" "winget" $shouldUpgrade # includes npm to start typescript/ javascript repositories

install "umlet" "umlet" "choco" $shouldUpgrade # required to generate images from uxf files