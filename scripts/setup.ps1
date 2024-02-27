# Description: This script sets up the environment for this project in order to function properly.
# It calls the scripts installDependencies.ps1, setupRepositories.ps1 and substitutes the api token in the .env files.
# Precondition: The .env file exists and contains the required credentials.
#               Admin rights are required to run this script.

#Requires -RunAsAdministrator

# set home and cli root path
set-location $PSScriptRoot
$cliRootPath = git rev-parse --show-toplevel
set-location "$cliRootPath/.."
$homePath = $PWD -replace '\\', '/'

# install dependencies
& $PSScriptRoot/installDependencies.ps1

# setup repositories
& $PSScriptRoot/setupRepositories.ps1 -destination $homePath

# substitute api tokens for ars and repo editor
$arsPath = "$homePath/divekit-automated-repo-setup"
(get-content "$arsPath/.env.example").replace("YOUR_API_TOKEN", "$env:API_TOKEN") | set-content "$arsPath/.env"
$repoEditorPath = "$homePath/divekit-repo-editor"
(get-content "$repoEditorPath/.env.example").replace("YOUR_API_TOKEN", "$env:API_TOKEN") | set-content "$repoEditorPath/.env"