#Requires -RunAsAdministrator

write-output "------------------------- Installing Dependencies --------------------------"
& $PSScriptRoot/install_dependencies.ps1
# refresh $env:Path to include newly installed dependencies
$Env:Path = [System.Environment]::GetEnvironmentVariable("Path", "Machine") + ";" + [System.Environment]::GetEnvironmentVariable("Path", "User")


write-output "------------------------- Preparing Setup ----------------------------------"
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

# set home and cli root path
set-location $PSScriptRoot
$cliRootPath = git rev-parse --show-toplevel
set-location "$cliRootPath/.."
$homePath = $PWD -replace '\\', '/'

# check if .env exists
if (-not(test-path -path "$cliRootPath/.env")) {
    write-output "The .env file does not exist. Copy or rename .env.example to .env and substitute `$USERNAME` and `$API_TOKEN` with your credentials, before running this script again."
    exit
}

# add git mingw64/bin folder to $env:Path temporarily. This is needed for envsubst
$gitPath = (Get-Command git).Source.Replace("\cmd\git.exe", "\mingw64\bin")
$env:Path += ";$gitPath"

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
if ($request -eq $null) {
    write-output "Could not invoke a web request successfully. The provided credentials in .env might be wrong"
    exit
}

# create test origin repository variables with the request
$testOriginRepo = ($request.content | convertfrom-json |
        select-object -expandproperty "http_url_to_repo").replace('https://', '').replace('.git', '')
$testOriginRepoName = $testOriginRepo.Substring($testOriginRepo.lastIndexOf('/') + 1)
$env:TEST_ORIGIN_REPO_FILE_PATH = "$homePath/$testOriginRepoName"



write-output "------------------------- Cloning Repositories -----------------------------"
set-location $homePath
git clone "https://github.com/divekit/divekit-automated-repo-setup"
git clone "https://github.com/divekit/divekit-repo-editor"
git clone "https://$( $env:USERNAME ):$( $env:API_TOKEN )@$( $testOriginRepo )"



write-output "------------------------- Setting Up ARS Repository ------------------------"
set-location "./divekit-automated-repo-setup"
git checkout "test_cli" #TEMPORARY UNTIL BRANCH IS MERGED
npm install
new-item -itemType "directory" -path "./resources/test/input"
new-item -itemType "directory" -path "./resources/test/output"
new-item -itemType "directory" -path "./resources/overview"
new-item -itemType "directory" -path "./resources/individual_repositories"
robocopy "./resources/examples/config" "./resources/config"
get-content "./.env.example" | envsubst '$API_TOKEN' | set-content "./.env"



write-output "------------------------- Setting Up Repo Editor Repository ----------------"
set-location "$homePath/divekit-repo-editor"
git checkout "test_cli" #TEMPORARY UNTIL BRANCH IS MERGED
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
