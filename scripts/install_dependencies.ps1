#Requires -RunAsAdministrator

function commandExists([string]$cmdName) {
    return [bool](get-command -Name $cmdName -ErrorAction SilentlyContinue)
}

function install([string]$cmdName, [string]$packageId, [string]$packageManager, [bool]$shouldUpgrade) {
    if (-not($shouldUpgrade) -and (commandExists($cmdName))) {
        write-output "Upgrade flag is set to false and the package $packageId is already installed. Skipping upgrade."
        return
    }

    if ($packageManager -eq "winget") {
        winget install --id $packageId --accept-source-agreements --accept-package-agreements
        return
    }
    choco install $packageId -y
}

#  Required to install dependencies
if (-not(commandExists("winget"))) {
    write-output "winget needs to be installed to install dependencies"
    exit
}
if (-not(commandExists("choco"))) {
    write-output "chocolatey needs to be installed to install dependencies"
    exit
}

# Upgrade if already installed
[bool]$shouldUpgrade = $false

# install [cmdName] [packageId] [packageManager] [shouldUpgrade]
install "git" "Git.Git" "winget" $shouldUpgrade

install "go" "GoLang.Go" "winget" $shouldUpgrade

install "java" "Oracle.JavaRuntimeEnvironment" "winget" $shouldUpgrade # umlet requires java runtime

install "node" "OpenJS.NodeJS" "winget" $shouldUpgrade # includes npm to start typescript/ javascript repositories

install "umlet" "umlet" "choco" $shouldUpgrade # required to generate images from uxf files