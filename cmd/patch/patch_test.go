package patch

import (
	"bytes"
	"divekit-cli/cmd"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestPatch(t *testing.T) {
	testCases := []struct {
		name      string
		arguments PatchArguments
		// repoIds // expected
	}{
		{
			"patch one file",
			PatchArguments{
				true,
				"",
				"divekit-origin-test-repo",
				"C:/Users/Thomas/Documents/Praxisprojekt/Git",
				"",
				[]string{"$PersonClass$RegistrationUseCases.java"},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			clearARSGeneratedFiles()
			_, err := executePatch(testCase.arguments)

			assert.NoError(t, err)
			if testCase.arguments.dryRun {
				// todo: check if files have been patched
			} else {
				// todo: check if files have been patched
			}
		})
	}
}

func clearARSGeneratedFiles() {

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
