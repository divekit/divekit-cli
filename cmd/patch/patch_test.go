package patch

import (
	"bytes"
	"divekit-cli/cmd"
	"divekit-cli/divekit/origin"
	"divekit-cli/utils/api"
	"divekit-cli/utils/errorHandling"
	"divekit-cli/utils/fileUtils"
	"divekit-cli/utils/logUtils"
	"fmt"
	"github.com/apex/log"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xanzy/go-gitlab"
	"os"
	"strings"
	"testing"
	"unicode/utf8"
)

var (
	client             *gitlab.Client    // Interacts with the GitLab API
	repositoryIds      map[string]string // Access an id by using a repository name as key
	homePath           string
	testOriginRepoName string
)

func TestPatch(t *testing.T) {
	testCases := []struct {
		name           string
		arguments      PatchArguments  // input
		generatedFiles []GeneratedFile // expected
		error          error           // expected
	}{
		{
			"[dry run] patch with no args should fail",
			PatchArguments{
				false,
				"",
				"",
				"",
				"",
				[]string{},
			},
			[]GeneratedFile{},
			&cmd.InvalidArgsError{},
		},
		{
			"[dry run] patch only with a patch file arg should fail",
			PatchArguments{
				false,
				"",
				"",
				"",
				"",
				[]string{"ProjectApplication.java"},
			},
			[]GeneratedFile{},
			&origin.OriginRepoError{},
		},
		{
			"[dry run] patch with a non existing patch file arg should fail",
			PatchArguments{
				false,
				"",
				testOriginRepoName,
				homePath,
				"",
				[]string{},
			},
			[]GeneratedFile{},
			&PatchFileError{},
		},
		{
			"[dry run] patch with an invalid home path should fail",
			PatchArguments{
				false,
				"",
				testOriginRepoName,
				"invalid_path",
				"",
				[]string{"ProjectApplication.java"},
			},
			[]GeneratedFile{},
			&fileUtils.InvalidPathError{},
		},
		{
			"[dry run] patch with an invalid origin repo name should fail",
			PatchArguments{
				false,
				"",
				"invalid_name",
				homePath,
				"",
				[]string{"ProjectApplication.java"},
			},
			[]GeneratedFile{},
			&fileUtils.InvalidPathError{},
		},
		{
			"[dry run] patch with an invalid log level should fail",
			PatchArguments{
				false,
				"invalid_level",
				testOriginRepoName,
				homePath,
				"",
				[]string{"ProjectApplication.java"},
			},
			[]GeneratedFile{},
			&logUtils.LogLevelError{},
		},
		{
			"[dry run] patch with an invalid distribution name should fail",
			PatchArguments{
				false,
				"",
				testOriginRepoName,
				homePath,
				"invalid_distribution",
				[]string{"ProjectApplication.java"},
			},
			[]GeneratedFile{},
			&origin.OriginRepoError{},
		},
		{
			"[dry run] patch a file",
			PatchArguments{
				true,
				"",
				testOriginRepoName,
				homePath,
				"",
				[]string{"ProjectApplication.java"},
			},
			[]GeneratedFile{
				{
					"ST1_Test_group_8063661e-3603-4b84-b780-aa5ff1c3fe7d",
					"src/main/java/thkoeln/st/st1praktikum/ProjectApplication.java",
					Code,
					nil,
					nil,
				},
				{
					"ST1_Test_group_86bd537d-9995-4c92-a6f4-bec97eeb7c67",
					"src/main/java/thkoeln/st/st1praktikum/ProjectApplication.java",
					Code,
					nil,
					nil,
				},
				{
					"ST1_Test_group_8754b8cb-5bc6-4593-9cb8-7c84df266f59",
					"src/main/java/thkoeln/st/st1praktikum/ProjectApplication.java",
					Code,
					nil,
					nil,
				},
			},
			nil,
		},
		{
			"[dry run] patch a wildcard file",
			PatchArguments{
				true,
				"",
				testOriginRepoName,
				homePath,
				"test",
				[]string{"$DonationClassName$.json"},
			},
			[]GeneratedFile{
				{
					"ST1_Test_tests_group_446e3369-ed35-473e-b825-9cc0aecd6ba3",
					"src/test/resources/milestones/milestone5/objectdescriptions/SponsoringAgreement.json",
					Test,
					nil,
					nil,
				},
				{
					"ST1_Test_tests_group_9672285a-67b0-4f2e-830c-72925ba8c76e",
					"src/test/resources/milestones/milestone5/objectdescriptions/SponsoringAgreement.json",
					Test,
					nil,
					nil,
				},
			},
			nil,
		},
		{
			"[dry run] patch two wildcard files with variations in test group",
			PatchArguments{
				true,
				"",
				testOriginRepoName,
				homePath,
				"test",
				[]string{"$E04M01Name$_E04M01.java", "$E07M03Name$_E07M03.java"},
			},
			[]GeneratedFile{
				{
					"ST1_Test_tests_group_446e3369-ed35-473e-b825-9cc0aecd6ba3",
					"src/main/java/thkoeln/st/basics/exercise/E04Methods/NumberOfVowels.java",
					Test,
					nil,
					nil,
				},
				{
					"ST1_Test_tests_group_446e3369-ed35-473e-b825-9cc0aecd6ba3",
					"src/main/java/thkoeln/st/basics/exercise/E07Lists/CanBalance.java",
					Test,
					nil,
					nil,
				},
				{
					"ST1_Test_tests_group_9672285a-67b0-4f2e-830c-72925ba8c76e",
					"src/main/java/thkoeln/st/basics/exercise/E07Lists/IsBalanceable.java",
					Test,
					nil,
					nil,
				},
			},
			nil,
		},
		{
			"[dry run] patch two wildcard files with variations in code group",
			PatchArguments{
				true,
				"",
				testOriginRepoName,
				homePath,
				"",
				[]string{"$E04M01Name$_E04M01.java", "$E07M03Name$_E07M03.java"},
			},
			[]GeneratedFile{
				{
					"ST1_Test_group_86bd537d-9995-4c92-a6f4-bec97eeb7c67",
					"src/main/java/thkoeln/st/basics/exercise/E07Lists/IsBalanceable.java",
					Code,
					nil,
					nil,
				},
				{
					"ST1_Test_group_8754b8cb-5bc6-4593-9cb8-7c84df266f59",
					"src/main/java/thkoeln/st/basics/exercise/E07Lists/IsBalanceable.java",
					Code,
					nil,
					nil,
				},
				{
					"ST1_Test_group_8063661e-3603-4b84-b780-aa5ff1c3fe7d",
					"src/main/java/thkoeln/st/basics/exercise/E07Lists/IsBalanceable.java",
					Code,
					nil,
					nil,
				},
				{
					"ST1_Test_group_86bd537d-9995-4c92-a6f4-bec97eeb7c67",
					"src/main/java/thkoeln/st/basics/exercise/E04Methods/VowelsInString.java",
					Code,
					nil,
					nil,
				},
				{
					"ST1_Test_group_8063661e-3603-4b84-b780-aa5ff1c3fe7d",
					"src/main/java/thkoeln/st/basics/exercise/E04Methods/VowelsInString.java",
					Code,
					nil,
					nil,
				},
				{
					"ST1_Test_group_8754b8cb-5bc6-4593-9cb8-7c84df266f59",
					"src/main/java/thkoeln/st/basics/exercise/E04Methods/NumberOfVowels.java",
					Code,
					nil,
					nil,
				},
			},
			nil,
		},
		{
			"[dry run] patch an uxf file",
			PatchArguments{
				true,
				"",
				testOriginRepoName,
				homePath,
				"test",
				[]string{"LDM.uxf"},
			},
			[]GeneratedFile{
				{
					"ST1_Test_tests_group_446e3369-ed35-473e-b825-9cc0aecd6ba3",
					"LDM.jpg",
					Test,
					nil,
					nil,
				},
				{
					"ST1_Test_tests_group_9672285a-67b0-4f2e-830c-72925ba8c76e",
					"LDM.jpg",
					Test,
					nil,
					nil,
				},
				{
					"ST1_Test_tests_group_446e3369-ed35-473e-b825-9cc0aecd6ba3",
					"LDM.uxf",
					Test,
					nil,
					nil,
				},
				{
					"ST1_Test_tests_group_9672285a-67b0-4f2e-830c-72925ba8c76e",
					"LDM.uxf",
					Test,
					nil,
					nil,
				},
			},
			nil,
		},
		{
			"[dry run] patch files with solution deletion",
			PatchArguments{
				true,
				"",
				testOriginRepoName,
				homePath,
				"test",
				[]string{"Factorial.java", "Main.java"},
			},
			[]GeneratedFile{
				{
					"ST1_Test_tests_group_446e3369-ed35-473e-b825-9cc0aecd6ba3",
					"src/main/java/thkoeln/st/basics/exercise/E08Testing/Factorial.java",
					Test,
					[]string{"UnsupportedOperationException"},
					[]string{"SHOULD BE DELETED", "ArithmeticException"},
				},
				{
					"ST1_Test_tests_group_9672285a-67b0-4f2e-830c-72925ba8c76e",
					"src/main/java/thkoeln/st/basics/exercise/E08Testing/Factorial.java",
					Test,
					[]string{"UnsupportedOperationException"},
					[]string{"SHOULD BE DELETED", "ArithmeticException"},
				},
			},
			nil,
		},
		{
			"patch a file to one repository",
			PatchArguments{
				false,
				"",
				testOriginRepoName,
				homePath,
				"test",
				[]string{"$E04M01Name$_E04M01.java"},
			},
			[]GeneratedFile{
				{
					"ST1_Test_tests_group_446e3369-ed35-473e-b825-9cc0aecd6ba3",
					"src/main/java/thkoeln/st/basics/exercise/E04Methods/NumberOfVowels.java",
					Test,
					nil,
					nil,
				},
			},
			nil,
		},
		{
			"patch two wildcard files with variations to multiple repositories",
			PatchArguments{
				false,
				"",
				testOriginRepoName,
				homePath,
				"",
				[]string{"$E06M05Name$_E06M05.java", "$E02M04Name$_E02M04.java"},
			},
			[]GeneratedFile{
				{
					"ST1_Test_group_8063661e-3603-4b84-b780-aa5ff1c3fe7d",
					"src/main/java/thkoeln/st/basics/exercise/E06Arrays/LengthOfUniqueArray.java",
					Code,
					nil,
					nil,
				},
				{
					"ST1_Test_group_86bd537d-9995-4c92-a6f4-bec97eeb7c67",
					"src/main/java/thkoeln/st/basics/exercise/E06Arrays/LengthOfUniqueArray.java",
					Code,
					nil,
					nil,
				},
				{
					"ST1_Test_group_8754b8cb-5bc6-4593-9cb8-7c84df266f59",
					"src/main/java/thkoeln/st/basics/exercise/E02Conditions/GetDayByNumber.java",
					Code,
					nil,
					nil,
				},
			},
			nil,
		},
		{
			"patch an uxf file to multiple repositories",
			PatchArguments{
				false,
				"",
				testOriginRepoName,
				homePath,
				"",
				[]string{"E2.uxf"},
			},
			[]GeneratedFile{
				{
					"ST1_Test_group_8063661e-3603-4b84-b780-aa5ff1c3fe7d",
					"src/main/resources/milestones/milestone3/exercises/exercise2/E2.jpg",
					Code,
					nil,
					nil,
				},
				{
					"ST1_Test_group_86bd537d-9995-4c92-a6f4-bec97eeb7c67",
					"src/main/resources/milestones/milestone3/exercises/exercise2/E2.jpg",
					Code,
					nil,
					nil,
				},
				{
					"ST1_Test_group_8754b8cb-5bc6-4593-9cb8-7c84df266f59",
					"src/main/resources/milestones/milestone3/exercises/exercise2/E2.jpg",
					Code,
					nil,
					nil,
				},
				{
					"ST1_Test_group_8063661e-3603-4b84-b780-aa5ff1c3fe7d",
					"src/main/resources/milestones/milestone3/exercises/exercise2/E2.uxf",
					Code,
					nil,
					nil,
				},
				{
					"ST1_Test_group_86bd537d-9995-4c92-a6f4-bec97eeb7c67",
					"src/main/resources/milestones/milestone3/exercises/exercise2/E2.uxf",
					Code,
					nil,
					nil,
				},
				{
					"ST1_Test_group_8754b8cb-5bc6-4593-9cb8-7c84df266f59",
					"src/main/resources/milestones/milestone3/exercises/exercise2/E2.uxf",
					Code,
					nil,
					nil,
				},
			},
			nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			generatedFiles := testCase.generatedFiles
			dryRunFlag := testCase.arguments.dryRun
			distributionFlag := testCase.arguments.distribution

			deleteFilesFromRepositories(t, generatedFiles, dryRunFlag)
			_, err := executePatch(testCase.arguments)

			checkErrorType(t, testCase.error, err)
			if err == nil {
				matchGeneratedFiles(t, generatedFiles, distributionFlag)
				checkFileContent(t, generatedFiles)
				checkPushedFiles(t, generatedFiles, dryRunFlag)
			}
		})
	}
}

// deleteFilesFromRepositories deletes the specified files from their respective repositories.
// Prior to testing, it is necessary to delete these files to ensure that they are actually pushed to the
// repositories, given that they are initially included in the repositories.
func deleteFilesFromRepositories(t *testing.T, files []GeneratedFile, dryRunActive bool) {
	if dryRunActive {
		t.Log("Dry Run flag set: SKIP DELETING remote files")
		return
	}

	for _, file := range files {
		repositoryId := repositoryIds[file.RepoName]
		filePath := file.RelFilePath
		if err := api.DeleteFileByRepositoryId(client, repositoryId, filePath); err != nil {
			t.Logf("The file %s does not exist in the repository `%s` (id: %s): %v",
				filePath, file.RepoName, repositoryId, err)
		} else {
			t.Logf("Deleted file %s from repository `%s` (id: %s)", filePath, file.RepoName, repositoryId)
		}
	}
}

// executePatch executes the patch command with the given arguments and returns the output and the error
func executePatch(args PatchArguments) (string, error) {
	root := createCmd()
	patchWithArgs := append([]string{"patch"}, args.splitIntoStrings()...)

	buffer := new(bytes.Buffer)
	root.SetOut(buffer)
	root.SetErr(buffer)
	root.SetArgs(patchWithArgs)

	err := root.Execute()
	output := buffer.String()

	return output, err
}

// checkErrorType checks if the expected error type matches with the actual error type
func checkErrorType(t *testing.T, expected error, actual error) {
	errorHandling.IsErrorType(t, expected, actual)
}

// matchGeneratedFiles checks if the found file paths match with the expected files and
// throws an error when there are any differences.
func matchGeneratedFiles(t *testing.T, expectedFiles []GeneratedFile, distributionFlag string) {
	distribution := newDistribution(distributionFlag)
	actualFilePaths := getGeneratedFilePaths(t, distribution)

	var expectedFilePaths []string
	for _, expectedFile := range expectedFiles {
		expectedFilePaths = append(expectedFilePaths, expectedFile.getFilePath())
	}

	require.ElementsMatchf(t, expectedFilePaths, actualFilePaths, "The expected file paths do not match with the found file paths")
}

// checkFileContent checks if the content of the files is correct.
func checkFileContent(t *testing.T, files []GeneratedFile) {
	for _, file := range files {
		byteContent, err := os.ReadFile(file.getFilePath())
		if err != nil {
			t.Fatalf("Could not read the file %s: %v", file.getFilePath(), err)
		}

		// Generated files should be UTF-8 encoded in order to check their content
		if !utf8.Valid(byteContent) {
			t.Logf("The file %s is not UTF-8 encoded: SKIPPING content check for this file.", file.getFileName())
			continue
		}

		content := string(byteContent)

		// Any file should contain at least one character
		assert.NotEmptyf(t, content, "The content of the file %s is empty", file.getFilePath())

		// Asserts that the content does not contain any delimiters,
		// because they indicate that wildcards have not been substituted correctly
		delimiter := "$"
		assert.NotContainsf(t, content, delimiter, "The file %s contains a %s delimiter", file.getFilePath(), delimiter)

		// Asserts that the content contains all included keywords
		for _, keyword := range file.Include {
			assert.Containsf(t, content, keyword, "The file %s should contain the keyword: %s", file.getFilePath(), keyword)
		}

		// Asserts that the content does not contain any excluded keywords
		for _, keyword := range file.Exclude {
			assert.NotContainsf(t, content, keyword, "The file %s should not contain the keyword: %s", file.getFilePath(), keyword)
		}
	}
}

// checkPushedFiles checks if the generated files have been pushed correctly to the corresponding repositories.
func checkPushedFiles(t *testing.T, localFiles []GeneratedFile, dryRunActive bool) {
	if dryRunActive {
		t.Log("Dry Run flag set: SKIPPING remote repository check")
		return
	}

	for _, localFile := range localFiles {
		currentId := repositoryIds[localFile.RepoName]
		remoteFile, err := api.GetFileByRepositoryId(client, currentId, localFile.RelFilePath)

		if remoteFileExists := assert.NoErrorf(t, err, "Could not get file %s from repository `%s` (id: %s): %v",
			localFile.RelFilePath, localFile.RepoName, currentId, err); !remoteFileExists {
			continue
		}

		assert.Equalf(t, localFile.getSHA256(), remoteFile.SHA256,
			"The checksum of the locally generated file does not match with that of the remote file.\n"+
				"The file may not have been pushed correctly to the repository. `%s` (id: %s)\n"+
				"Provided file: %s", localFile.RepoName, currentId, localFile.getFilePath())
	}
}

func createCmd() *cobra.Command {
	rootCmd := cmd.NewRootCmd()
	cmd.SetCmdFlags(rootCmd)

	patchCmd := NewPatchCmd()
	initFlags()
	setCmdFlags(patchCmd)
	rootCmd.AddCommand(patchCmd)

	return rootCmd
}

type PatchArguments struct {
	dryRun       bool
	logLevel     string
	originRepo   string
	home         string
	distribution string
	patchFiles   []string
}

func (p PatchArguments) splitIntoStrings() []string {
	result := strings.Split(p.toString(), " ")

	// Return an empty slice if no arguments defined
	if len(result) == 1 && result[0] == "" {
		return []string{}
	}

	return result
}

func (p PatchArguments) toString() string {
	result := ""
	if p.dryRun {
		result += "-0 "
	}
	if p.logLevel != "" {
		result += "-l " + p.logLevel + " "
	}
	if p.originRepo != "" {
		result += "-o " + p.originRepo + " "
	}
	if p.home != "" {
		result += "-m " + p.home + " "
	}
	if p.distribution != "" {
		result += "-d " + p.distribution + " "
	}
	if len(p.patchFiles) > 0 {
		result += strings.Join(p.patchFiles, " ")
	}

	return result
}

type GeneratedFile struct {
	RepoName     string
	RelFilePath  string
	Distribution Distribution
	Include      []string
	Exclude      []string
}

func (g GeneratedFile) getFilePath() string {
	return getARSOutputDir() + "/" + string(g.Distribution) + "/" + g.RepoName + "/" + g.RelFilePath
}

func (g GeneratedFile) getFileName() string {
	return fileUtils.GetBaseName(g.RelFilePath)
}

func (g GeneratedFile) getSHA256() string {
	return fileUtils.GetSHA256(g.getFilePath())
}

type Distribution string

const (
	Code Distribution = "code"
	Test Distribution = "test"
)

func newDistribution(distributionString string) Distribution {
	switch distributionString {
	case "code":
		return Code
	case "":
		return Code
	case "test":
		return Test
	default:
		log.Fatalf("Could not create a distribution: The distribution string is invalid: %s", distributionString)
		return ""
	}
}

func getGeneratedFilePaths(t *testing.T, distribution Distribution) []string {
	outputDir := getARSOutputDir() + "/" + string(distribution)
	foundPaths, err := fileUtils.FindAnyFilesInDirRecursively(outputDir)
	require.NoErrorf(t, err, "Could not find any files in the directory %s: %v", outputDir, err)

	return fileUtils.UnifyPaths(foundPaths)
}
func getARSOutputDir() string {
	if ARSRepo == nil {
		log.Fatalf("Could not find the generated output directory: ARSRepo is nil")
	}

	dir := fileUtils.UnifyPath(ARSRepo.GeneratedLocalOutput.Dir)
	if dir == "" {
		log.Fatalf("Could not find the generated output directory: The directory path is empty")
	}

	return dir
}
func TestMain(m *testing.M) {
	setup()
	e := m.Run() // run the tests
	os.Exit(e)   // report the exit code
}

func setup() {
	fileUtils.LoadEnv()
	host := os.Getenv("HOST")
	token := os.Getenv("API_TOKEN")
	codeGroupId := os.Getenv("CODE_GROUP_ID")
	testGroupId := os.Getenv("TEST_GROUP_ID")
	testOriginRepoId := os.Getenv("TEST_ORIGIN_REPO_ID")

	var err error
	if client, err = api.NewGitlabClient(host, token); err != nil {
		log.Fatalf("", err)
	}
	initRepositoryIds([]string{codeGroupId, testGroupId})
	homePath = fileUtils.GetHomePath()
	testOriginRepoName = getTestOriginRepoName(testOriginRepoId)
}

func initRepositoryIds(groupIds []string) {
	repositoryIds = make(map[string]string)

	for _, groupId := range groupIds {
		repositories, err := api.GetRepositoriesByGroupId(client, groupId)
		if err != nil {
			log.Fatalf(fmt.Sprintf("could not get projects with group id %s: %v", groupId, err))
		}

		for _, repository := range repositories {
			id := fmt.Sprintf("%d", repository.ID)
			repositoryIds[repository.Name] = id
		}
	}
}

func getTestOriginRepoName(testOriginRepoId string) string {
	repository, err := api.GetRepositoryById(client, testOriginRepoId)
	if err != nil {
		log.Fatalf("", err)
	}

	return strings.ToLower(repository.Name)
}
