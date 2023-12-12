package patch

import (
	"bytes"
	"divekit-cli/cmd"
	"divekit-cli/utils/fileUtils"
	"encoding/json"
	"github.com/apex/log"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"os"
	"regexp"
	"strings"
	"testing"
	"unicode/utf8"
)

const (
	originRepo = "divekit-origin-test-repo"
	home       = "C:/Users/Thomas/Documents/Praxisprojekt/Git"
)

func TestPatch(t *testing.T) {
	testCases := []struct {
		name      string
		arguments PatchArguments // input
	}{
		{
			"patch a file",
			PatchArguments{
				true,
				"",
				originRepo,
				home,
				"",
				[]string{"ProjectApplication.java"},
			},
		},
		{
			"patch two wildcard files",
			PatchArguments{
				true,
				"",
				originRepo,
				home,
				"test",
				[]string{"$PersonClass$RegistrationUseCases.java", "$ObjectClass$CatalogUseCases.java"},
			},
		},
		{
			"patch an uxf file",
			PatchArguments{
				true,
				"",
				originRepo,
				home,
				"test",
				[]string{"E01Solution.uxf"},
			},
		},
		{
			"patch a wildcard uxf files",
			PatchArguments{
				true,
				"",
				originRepo,
				home,
				"test",
				[]string{"Test$PersonClass$.uxf"},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			dryRunActive := testCase.arguments.dryRun
			_, err := executePatch(testCase.arguments)

			assert.NoError(t, err)
			checkGeneratedFiles(t, testCase.arguments)
			if !dryRunActive {
				checkPushedFiles()
			}
		})
	}
}

// Check if the generated files have been successfully created
func checkGeneratedFiles(t *testing.T, args PatchArguments) {
	// These constants are required to determine, for instance, the expected file names
	// e.g. $PersonClass$RegistrationUseCases.java
	const delimiter = "$"
	const suffix = "Class"

	expectedFileNames := getGeneratedFileNames(args.patchFiles, delimiter, suffix)
	outputDir := getGeneratedARSOutputDir(args)
	foundPaths, _ := fileUtils.FindAnyFilesInDirRecursively(outputDir)
	foundFileNames := fileUtils.GetBaseNames(foundPaths...)

	assert.ElementsMatch(t, expectedFileNames, foundFileNames, "expected files do not match the found files")
	checkFileContent(t, foundPaths, delimiter)
}

func checkFileContent(t *testing.T, foundPaths []string, delimiter string) {
	for _, foundPath := range foundPaths {
		content, err := os.ReadFile(foundPath)
		if err != nil {
			log.Fatalf("could not read the file: %s", foundPath)
		}

		// Generated files should be UTF-8 encoded in order to test their content
		if !utf8.Valid(content) {
			log.Warn("The file " + fileUtils.GetBaseName(foundPath) + " is not UTF-8 encoded => " +
				"Skipping content check.")
			continue
		}

		// Generated files should not contain any delimiters,
		// because delimiters indicate that wildcards have not been substituted correctly
		assert.NotContainsf(t, string(content), delimiter, "%s contains a %s delimiter",
			fileUtils.GetBaseName(foundPath), delimiter)
	}
}

// getGeneratedFileNames returns generated file names for the given patch files.
// It supports wildcard names, such as obtaining `CustomerRegistrationUseCases.java`, `UserRegistrationUseCases.java`
// etc. for `$PersonClass$RegistrationUseCases.java`.
// Also, it considers special file extensions, such as `.uxf`, which means additional files gets generated.
func getGeneratedFileNames(patchFiles []string, delimiter string, suffix string) []string {
	var result []string

	for _, patchFile := range patchFiles {
		fileNames := getFileNames(patchFile, delimiter, suffix)
		additionalFileNames := getAdditionalFileNames(fileNames)

		result = append(result, fileNames...)
		result = append(result, additionalFileNames...)
	}

	return result
}

func getFileNames(patchFile string, delimiter string, suffix string) []string {
	individualRepositoriesFile := ARSRepo.IndividualizationConfig.Dir + "/individual_repositories.json"

	// Check if the patch file contains any delimiters.
	// If not, return the patch file multiplied by the number of members from the individual_repositories.json file
	// e.g. ["file.java"] * 2 => ["file.java", "file.java"]
	if !strings.Contains(patchFile, delimiter) {
		return multiplyStringByN(patchFile, getNumberOfMembers(individualRepositoriesFile))
	}

	// Get individual values from the individual_repositories.json file
	objectName := getIndividualObjectName(patchFile, delimiter, suffix)
	objectValues := getIndividualObjectValues(individualRepositoriesFile, objectName)

	// Generate file names for each individual value
	var generatedFileNames []string
	for _, value := range objectValues {
		generatedFileName := strings.ReplaceAll(patchFile, delimiter+objectName+suffix+delimiter, value)
		generatedFileNames = append(generatedFileNames, generatedFileName)
	}

	return generatedFileNames
}

func getAdditionalFileNames(fileNames []string) []string {
	var result []string

	for _, fileName := range fileNames {
		if strings.HasSuffix(fileName, ".uxf") {
			result = append(result, strings.TrimSuffix(fileName, ".uxf")+".jpg")
		}
	}

	return result
}

func getIndividualObjectName(fileName string, delimiter string, suffix string) string {
	// Define a regular expression to match the name within delimiter...delimiter in the filename
	re := regexp.MustCompile("\\" + delimiter + "(.+)\\" + delimiter)

	// Find the submatches in the file name
	matches := re.FindStringSubmatch(fileName)

	// Check if any matches were found
	if len(matches) < 2 {
		log.Fatalf("could not match the filename: %s", fileName)
	}

	// Get the objectName and remove trailing suffixes
	objectName := strings.TrimSuffix(matches[1], suffix)

	return objectName
}

func getIndividualObjectValues(configPath string, objectName string) []string {
	// Collect all values for the given objectName
	var objectValues []string
	for _, dataBlock := range unmarshalConfig(configPath) {
		if objSelection, ok := dataBlock["individualSelectionCollection"].(map[string]interface{}); ok {
			if objects, ok := objSelection["individualObjectSelection"].(map[string]interface{}); ok {
				if objectValue, ok := objects[objectName].(string); ok {
					objectValues = append(objectValues, objectValue)
				}
			}
		}
	}

	return objectValues
}

func getNumberOfMembers(configPath string) int {
	return len(unmarshalConfig(configPath))
}

func unmarshalConfig(configPath string) []map[string]interface{} {
	// Validate the config path
	if err := fileUtils.ValidateFilePath(configPath); err != nil {
		log.Fatalf("could not validate the config path: %s", configPath)
	}

	// Read the config file
	content, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("could not read the config file: %s", configPath)
	}

	// Unmarshal JSON data into a slice of DataBlock
	var dataBlocks []map[string]interface{}
	if err := json.Unmarshal(content, &dataBlocks); err != nil {
		log.Fatalf("could not unmarshal the config file: %s", configPath)
	}

	return dataBlocks
}
func getGeneratedARSOutputDir(args PatchArguments) string {
	if args.distribution != "" {
		return ARSRepo.GeneratedLocalOutput.Dir + "/test"
	}
	return ARSRepo.GeneratedLocalOutput.Dir + "/code"
}

func multiplyStringByN(s string, n int) []string {
	var result []string

	for i := 0; i < n; i++ {
		result = append(result, s)
	}

	return result
}

// check if the generated files have been pushed to the corresponding repositories
func checkPushedFiles() {

}
func executePatch(args PatchArguments) (string, error) {
	root := createCmd()
	cmdWithArgs := strings.Split("patch "+args.toString(), " ")

	buffer := new(bytes.Buffer)
	root.SetOut(buffer)
	root.SetErr(buffer)
	root.SetArgs(cmdWithArgs)

	err := root.Execute()
	output := buffer.String()
	println(output)

	return output, err
}
func createCmd() *cobra.Command {
	rootCmd := cmd.NewRootCmd()
	cmd.SetCmdFlags(rootCmd)

	patchCmd := NewPatchCmd()
	initVariables()
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

//func TestPreRun(t *testing.T) {
//	// todo
//	//patch = &cobra.Command{Use: "root", RunE: PreRun}
//}
//
//func TestDefinePatchFiles(t *testing.T) {
//	// todo
//	//definePatchFiles()
//}
//
//func TestCopySavedIndividualizationFileToARS(t *testing.T) {
//	// todo
//	//copySavedIndividualizationFileToARS()
//}
//
//func TestCopyLocallyGeneratedFilesToPatchTool(t *testing.T) {
//	// todo
//	//copyLocallyGeneratedFilesToPatchTool()
//}
