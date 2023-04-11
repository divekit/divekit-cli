package cmd

import (
	"divekit-cli/utils"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var (
	// Flags
	oneUserPerRun    bool
	distributionName string
	patchPath        string

	// command state vars
	patchFiles      []string
	dirPathToLookIn string

	patchCmd = &cobra.Command{
		Use:   "patch",
		Short: "Apply a patch to all repos",
		Long:  `Patch one or several files in all the repos of a certain distribution of the origin repo`,
		Args:  validateArgs,

		Run: run,
	}
)

func init() {
	patchCmd.Flags().BoolVarP(&oneUserPerRun, "oneuser", "1", true,
		"users in the repo distribution are patched one-by-one, in order to avoid memory overflow")
	patchCmd.Flags().StringVarP(&distributionName, "distribution", "d", "milestone",
		"name of the repo-distribution to patch")
	patchCmd.Flags().StringVarP(&patchPath, "patchpath", "p", "",
		"directory patchPath within the origin repo containing the patch files (default: root)")
	dirPathToLookIn = filepath.Join(OriginRepoFullPath, patchPath)

	// Add the patch subcommand to the divekit command
	rootCmd.AddCommand(patchCmd)
}

// Check if the directory exists and contains a ".divekit" subfolder
func validateArgs(cmd *cobra.Command, args []string) error {
	var err error
	if len(args) == 0 {
		err = fmt.Errorf("You need to specify at least one filename to patch.")
	}
	/*
		if OriginRepoFullPath == "" {
			err = fmt.Errorf("You need to specify the origin repo with the -o / --originrepo flag.")
		}
		if err != nil {
			utils.ValidateRepoPath(OriginRepoFullPath, true)
		}
	*/
	return err
}

func run(cmd *cobra.Command, args []string) {
	configFilename := fmt.Sprintf("%s.repositoryConfig.json", distributionName)
	configFilePath := filepath.Join(DistributionsDirFullPath, configFilename)

	_, err := os.Stat(configFilePath)
	if os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "File does not exist: %s", configFilePath)
		os.Exit(1)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error checking for file: %s", err)
		os.Exit(1)
	}

	definePatchFiles(args)
}

func definePatchFiles(args []string) {
	for index, _ := range args {
		foundFiles, foundErr := utils.FindFiles(args[index], dirPathToLookIn)
		if foundErr != nil {
			fmt.Fprintf(os.Stderr, "%s", foundErr)
			os.Exit(1)
		}
		if len(foundFiles) == 0 {
			fmt.Fprintf(os.Stderr, "No files found with name %s", args[0])
			os.Exit(1)
		}
		if len(foundFiles) > 1 {
			errorMsg := "Error: Multiple files found:\n"
			for _, file := range foundFiles {
				errorMsg += fmt.Sprintf("  - %s\n", file)
			}
			fmt.Fprintf(os.Stderr, "%s", errorMsg)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stdout, "Found file: %s", foundFiles[0])
		patchFiles[index] = foundFiles[0]
	}
}
