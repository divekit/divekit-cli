#Requires -RunAsAdministrator

write-output "------------------------- Installing Dependencies --------------------------"
& $PSScriptRoot/install_dependencies.ps1
$Env:Path = [System.Environment]::GetEnvironmentVariable("Path","Machine") + ";" + `
    [System.Environment]::GetEnvironmentVariable("Path","User")


write-output "------------------------- Preparing Setup ----------------------------------"
# set home and cli root path
set-location $PSScriptRoot
$cliRootPath= git rev-parse --show-toplevel
set-location "$cliRootPath/.."
$homePath = $PWD -replace '\\', '/'
# add git mingw64/bin folder to $env:Path temporarily. This is needed for envsubst
$gitPath = (Get-Command git).Source.Replace("\cmd\git.exe", "\mingw64\bin")
$env:Path += ";$gitPath"
# load env variables from .env file
set-location $cliRootPath
get-content .env | foreach-object {
    $name, $value = $_.split('=')
    set-content env:$name $value
}
# set test origin variables
$testOriginRepo = (((invoke-webrequest -uri "$env:HOST/api/v4/projects/$env:TEST_ORIGIN_REPO_ID" `
    -method GET `
    -headers @{ "PRIVATE-TOKEN" = "$env:API_TOKEN" } | `
    select-object -expandProperty content) -split ',' | `
    select-string -pattern "http_url_to_repo") -split '"http_url_to_repo":"https://' | `
    select-object -last 1) -replace '.git"', ''
$testOriginRepoName = $testOriginRepo.Substring($testOriginRepo.lastIndexOf('/') + 1)
$env:TEST_ORIGIN_REPO_FILE_PATH = "$homePath/$testOriginRepoName"


write-output "------------------------- Cloning Repositories -----------------------------"
set-location $homePath
git clone "https://github.com/divekit/divekit-automated-repo-setup"
git clone "https://github.com/divekit/divekit-repo-editor"
git clone "https://$( $env:USERNAME ):$( $env:API_TOKEN )@$( $testOriginRepo )"



write-output "------------------------- Setting Up ARS Repository ------------------------"
set-location "./divekit-automated-repo-setup"
git checkout "test_cli"
npm install
new-item -itemType "directory" -path "./resources/test/input"
new-item -itemType "directory" -path "./resources/test/output"
new-item -itemType "directory" -path "./resources/overview"
new-item -itemType "directory" -path "./resources/individual_repositories"
robocopy "./resources/examples/config" "./resources/config"
get-content "./.env.example" | envsubst '$API_TOKEN' | set-content "./.env"



write-output "------------------------- Setting Up Repo Editor Repository ----------------"
set-location "$homePath/divekit-repo-editor"
git checkout "test_cli"
npm install
get-content "./.env.example" | envsubst '$API_TOKEN' | set-content "./.env"
new-item -itemType "directory" -path "./assets/input/code"
new-item -itemType "directory" -path "./assets/input/test"



write-output "------------------------- Setting Up Test Origin Repository ----------------"
set-location $homePath/divekit-origin-test-repo
$distributions = "./.divekit_norepo/distributions"
$repositoryConfig = "repositoryConfig.json"
get-content "$distributions/milestone/$repositoryConfig" | envsubst '$TEST_ORIGIN_REPO_FILE_PATH' | `
    set-content "$distributions/milestone/$repositoryConfig"
get-content "$distributions/test/$repositoryConfig" | envsubst '$TEST_ORIGIN_REPO_FILE_PATH' | `
    set-content "$distributions/test/$repositoryConfig"



write-output "------------------------- Finished Setup -----------------------------------"
