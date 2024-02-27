# Description: This script installs all dependencies required to build and run this project.
# Precondition: Admin rights are required to run this script.

#Requires -RunAsAdministrator

function refreshEnvPath() {
    $Env:Path = [System.Environment]::GetEnvironmentVariable("Path", "Machine") + ";" +
            [System.Environment]::GetEnvironmentVariable("Path", "User")
}

function commandExists([string]$cmdName) {
    return [bool](get-command -Name $cmdName -ErrorAction SilentlyContinue)
}

function install([string]$cmdName, [string]$packageId) {
    # check if package is already installed
    if (commandExists($cmdName)) {
        write-output "$packageId is already installed"
        return
    }

    choco install $packageId -y
}

function initialize() {
    write-output "------------------------- Installing Dependencies --------------------------"
    if (-not(commandExists("choco"))) {
        # install choco: https://docs.chocolatey.org/en-us/choco/setup#install-with-powershell.exe
        Set-ExecutionPolicy Bypass -Scope Process -Force;
        [System.Net.ServicePointManager]::SecurityProtocol =
        [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; invoke-expression ((New-Object `
        System.Net.WebClient).DownloadString('https://chocolatey.org/install.ps1'))

        # refresh env path to make choco available
        refreshEnvPath
    }
}

initialize

# install [cmdName] [packegeId]

# call git commands
install "git" "git"

# build and run this project
install "go" "golang"

# required for umlet
install "java" "javaruntime"

# generate images from uxf files
install "umlet" "umlet"

# interact with ts/js repositories (e.g. npm)
install "node" "nodejs"

# refresh $env:Path to include newly installed dependencies
refreshEnvPath