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
	eachUser         bool
	distributionName string

	patchCmd = &cobra.Command{
		Use:   "patch",
		Short: "Apply a patch to all repos",
		Long:  `Patch one or several files in all the repos of a certain distribution of the origin repo`,
		Args:  validateDirectory,
		Run:   runPatch,
	}
)

func init() {
	patchCmd.Flags().BoolVarP(&eachUser, "eachuser", "e", true,
		"users in the repo distribution are patched one-by-one, in order to avoid memory overflow")
	patchCmd.Flags().StringVarP(&distributionName, "distribution", "d", "milestone",
		"name of the repo-distribution to patch")
	patchCmd.MarkFlagRequired("originrepo")

	// Add the patch subcommand to the divekit command
	rootCmd.AddCommand(patchCmd)
}

// Check if the directory exists and contains a ".divekit" subfolder
func validateDirectory(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("You need to specify at least one filename to patch.")
	}
	err := utils.ValidateRepoPath(originRepoName, true)
	return err
}

func runPatch(cmd *cobra.Command, args []string) {
	repoDir := originRepoName
	distributionsPath := filepath.Join(repoDir, ".divekit", "distributions")
	configFilename := fmt.Sprintf("%s.repositoryConfig.json", distributionName)
	configFilePath := filepath.Join(distributionsPath, configFilename)

	_, err := os.Stat(configFilePath)
	if os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "File does not exist: %s", configFilePath)
		os.Exit(1)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error checking for file: %s", err)
		os.Exit(1)
	}

	foundFiles, foundErr := utils.FindFiles(args[0], repoDir)
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
}
