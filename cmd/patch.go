package cmd

import (
	"divekit-cli/config"
	"divekit-cli/utils"
	"fmt"
	"github.com/apex/log"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	// Flags
	OneUserPerRun    bool
	DistributionName string
	PatchPath        string

	// command state vars
	PatchFiles                   []string
	RootDirToLookForPatchFiles   string
	RepositoryConfigFileFullPath string

	patchCmd = &cobra.Command{
		Use:    "patch",
		Short:  "Apply a patch to all repos",
		Long:   `Patch one or several files in all the repos of a certain distribution of the origin repo`,
		Args:   validateArgs,
		PreRun: preRun,
		Run:    run,
	}
)

func init() {
	log.Debug("patch.init()")
	patchCmd.Flags().BoolVarP(&OneUserPerRun, "oneuser", "1", true,
		"users in the repo distribution are patched one-by-one, in order to avoid memory overflow")
	patchCmd.Flags().StringVarP(&DistributionName, "distribution", "d", "milestone",
		"name of the repo-distribution to patch")
	patchCmd.Flags().StringVarP(&PatchPath, "patchpath", "p", "",
		"directory PatchPath within the origin repo containing the patch files (default: root)")

	// Add the patch subcommand to the divekit command
	rootCmd.AddCommand(patchCmd)
}

// Check if the directory exists and contains a ".divekit" subfolder
func validateArgs(cmd *cobra.Command, args []string) error {
	log.Debug("patch.validateArgs()")
	var err error
	if len(args) == 0 {
		err = fmt.Errorf("You need to specify at least one filename to patch.")
	}
	return err
}

// Checks preconditions before running the command
func preRun(cmd *cobra.Command, args []string) {
	if OriginRepoFullPath == "" {
		err := fmt.Errorf("You need to specify the origin repo with the -o / --originrepo flag.")
		fmt.Fprintf(os.Stderr, "%s", err)
		os.Exit(1)
	}
	configFilename := fmt.Sprintf("%s.%s", DistributionName, REPOSITORY_CONFIG_FILENAME)
	RepositoryConfigFileFullPath = filepath.Join(DistributionsDirFullPath, configFilename)
	RootDirToLookForPatchFiles = filepath.Join(OriginRepoFullPath, PatchPath)

	utils.OutputAndAbortIfErrors(utils.ValidateAllFilePaths(RepositoryConfigFileFullPath))
	utils.OutputAndAbortIfErrors(utils.ValidateAllDirPaths(RootDirToLookForPatchFiles))

	log.WithFields(log.Fields{
		"RepositoryConfigFileFullPath": RepositoryConfigFileFullPath,
		"RootDirToLookForPatchFiles":   RootDirToLookForPatchFiles,
	}).Debug("Setting patch variables")
}

func run(cmd *cobra.Command, args []string) {
	log.Debug("patch.run()")
	definePatchFiles(args)
	log.Info(fmt.Sprintf("Found files to patch:\n%s", strings.Join(PatchFiles, "\n")))
	setARSRepositoryConfig()
	runLocalGeneration()
}

func definePatchFiles(args []string) {
	log.Debug("patch.definePatchFiles()")
	for index, _ := range args {
		foundFiles, foundErr := utils.FindFiles(args[index], RootDirToLookForPatchFiles)
		if foundErr != nil {
			fmt.Fprintf(os.Stderr, "%s", foundErr)
			os.Exit(1)
		}
		if len(foundFiles) == 0 {
			fmt.Fprintf(os.Stderr, "No files found with name %s", args[index])
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
		log.Debug(fmt.Sprintf("Found file %s", foundFiles[0]))
		relFile, err := utils.TransformIntoRelativePaths(OriginRepoFullPath, foundFiles[0])
		log.Debug(fmt.Sprintf("... relative to origin repo: %s", relFile))
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s", err)
			os.Exit(1)
		}
		PatchFiles = append(PatchFiles, relFile)
	}
}

func setARSRepositoryConfig() {
	log.Debug("patch.setARSRepositoryConfig()")
	config.ReadConfigRepository(RepositoryConfigFileFullPath)
	config.ConfigRepository.Local.SubsetPaths = PatchFiles
	config.ConfigRepository.General.LocalMode = true
	config.WriteConfigRepository(ARSRepositoryConfigFileFullPath)
}

func runLocalGeneration() {
	log.Debug("patch.runLocalGeneration()")
	// Store the original directory
	originalDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current directory: %v", err)
	}
	err = os.Chdir(ARSRepoFullPath)
	if err != nil {
		log.Fatalf("Error changing directory to %s: %v", ARSRepoFullPath, err)
	}

	// Run "npm start"
	cmd := exec.Command("npm", "start")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Fatalf("Error running 'npm start': %v", err)
	}

	// Change back to the original directory
	err = os.Chdir(originalDir)
	if err != nil {
		log.Fatalf("Error changing back to the original directory: %v", err)
	}
	log.Info("Execution completed.")
}
