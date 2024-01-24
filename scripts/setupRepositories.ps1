# Description: This script clones the necessary repositories to a specified location and configures these repositories.
# Precondition: The .env file exists and contains the required credentials.
#               Admin rights are required to run this script.

#Requires -RunAsAdministrator

param([string]$destination)

write-output "------------------------- Preparing Setup ----------------------------------"

# CHANGE TO "main" WHEN BRANCH IS MERGED
$branchName = "test_cli"

# check if the destination path is passed as an argument
if ( [string]::IsNullOrEmpty($destination)) {
    write-output "The destination path is not passed as an argument"
    exit
}

# set cli root path
set-location $PSScriptRoot
$cliRootPath = git rev-parse --show-toplevel

# check if git is installed
if (-not(get-command -Name "git" -ErrorAction SilentlyContinue)) {
    write-output "Git needs to be installed to run this setup properly."
    exit
}

# check if nodejs is installed (includes npm)
if (-not(get-command -Name "node" -ErrorAction SilentlyContinue)) {
    write-output "NodeJs needs to be installed to run this setup properly."
    exit
}

# check if .env exists
if (-not(test-path -path "$cliRootPath/.env")) {
    write-output "The .env file does not exist. Copy or rename .env.example to .env and substitute `YOUR_GITLAB_USERNAME` and `YOUR_GITLAB_API_TOKEN` with your credentials"
    exit
}

# load env variables from .env file
set-location $cliRootPath
get-content .env | foreach-object {
    $name, $value = $_.split('=')
    set-content env:$name $value
}

# get test origin repository info
$request = invoke-webrequest -uri "$env:HOST/api/v4/projects/$env:TEST_ORIGIN_REPO_ID" `
    -useBasicParsing `
    -method GET `
    -headers @{ "PRIVATE-TOKEN" = "$env:API_TOKEN" }

# check if request was successful
if ($null -eq $request) {
    write-output "Could not invoke a web request successfully. The provided credentials in .env might be wrong"
    exit
}

# create test origin repository variables with the request
$testOriginRepo = ($request.content | convertfrom-json |
        select-object -expandproperty "http_url_to_repo").replace('https://', '').replace('.git', '')


write-output "------------------------- Cloning Repositories -----------------------------"
set-location $destination
git clone "https://github.com/divekit/divekit-automated-repo-setup"
git clone "https://github.com/divekit/divekit-repo-editor"
git clone "https://$( $env:USERNAME ):$( $env:API_TOKEN )@$( $testOriginRepo )"


write-output "------------------------- Setting Up ARS Repository ------------------------"
set-location "$destination/divekit-automated-repo-setup"
git checkout $branchName
npm install
new-item -itemType "file" -path "./resources/test/output/.keep" -force
new-item -itemType "file" -path "./resources/overview/.keep" -force
new-item -itemType "file" -path "./resources/individual_repositories/.keep" -force
robocopy "./resources/examples/config" "./resources/config"


write-output "------------------------- Setting Up Repo Editor Repository ----------------"
set-location "$destination/divekit-repo-editor"
git checkout $branchName
npm install
new-item -itemType "file" -path "./assets/input/code/.keep" -force
new-item -itemType "file" -path "./assets/input/test/.keep" -force


write-output "------------------------- Setting Up Origin Test Repository ----------------"
# define variables
$testOriginRepoName = $testOriginRepo.Substring($testOriginRepo.lastIndexOf('/') + 1)
$testOriginRepoPath = "$destination/$testOriginRepoName"
$distributions = "$destination/$testOriginRepoName/.divekit_norepo/distributions"
$repositoryConfig = "repositoryConfig.json"

# substitute repository config paths for origin-test-repo
(get-content "$distributions/milestone/$repositoryConfig").
        replace("YOUR_TEST_ORIGIN_REPO_PATH", "$testOriginRepoPath") |
        set-content "$distributions/milestone/$repositoryConfig"
(get-content "$distributions/test/$repositoryConfig").
        replace("YOUR_TEST_ORIGIN_REPO_PATH", "$testOriginRepoPath") |
        set-content "$distributions/test/$repositoryConfig"

write-output "------------------------- Finished Setup -----------------------------------"
