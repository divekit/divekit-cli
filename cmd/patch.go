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
	OneUserPerRunFlag    bool
	DistributionNameFlag string
	PatchPathFlag        string

	// command state vars
	PatchFiles []string
	ARSRepo    *config.ARSRepoType

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
	patchCmd.Flags().BoolVarP(&OneUserPerRunFlag, "oneuser", "1", true,
		"users in the repo distribution are patched one-by-one, in order to avoid memory overflow")
	patchCmd.Flags().StringVarP(&DistributionNameFlag, "distribution", "d", "milestone",
		"name of the repo-distribution to patch")

	patchCmd.MarkPersistentFlagRequired("originrepo")
	rootCmd.AddCommand(patchCmd)
}

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
	ARSRepo = config.ARSRepo()

	distribution := OriginRepo.GetDistribution(DistributionNameFlag)
	if distribution == nil {
		log.WithFields(log.Fields{
			"DistributionNameFlag": DistributionNameFlag,
		}).Fatal("Distribution not found")
	}
}

func run(cmd *cobra.Command, args []string) {
	log.Debug("patch.run()")
	definePatchFiles(args)
	log.Info(fmt.Sprintf("Found files to patch:\n%s", strings.Join(PatchFiles, "\n")))
	setRepositoryConfigWithinARSRepo()
	copySavedIndividualizationFileToARS()
	runLocalGeneration()
}

func definePatchFiles(args []string) {
	log.Debug("patch.definePatchFiles()")
	srcDir := filepath.Join(OriginRepo.RepoDir, "src")
	for index, _ := range args {
		foundFiles, foundErr := utils.FindFilesInDir(args[index], OriginRepo.RepoDir)
		foundFiles2, foundErr2 := utils.FindFilesInDirRecursively(args[index], srcDir)
		foundFiles = append(foundFiles, foundFiles2...)
		if foundErr != nil || foundErr2 != nil {
			fmt.Fprintf(os.Stderr, "%s %s", foundErr, foundErr2)
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
		relFile, err := utils.TransformIntoRelativePaths(OriginRepo.RepoDir, foundFiles[0])
		log.Debug(fmt.Sprintf("... relative to origin repo: %s", relFile))
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s", err)
			os.Exit(1)
		}
		PatchFiles = append(PatchFiles, relFile)
	}
}

func setRepositoryConfigWithinARSRepo() {
	log.Debug("patch.setRepositoryConfigWithinARSRepo()")
	distribution := OriginRepo.GetDistribution(DistributionNameFlag)
	if distribution == nil {
		log.WithFields(log.Fields{
			"DistributionNameFlag": DistributionNameFlag,
		}).Fatal("Distribution not found")
		os.Exit(1)
	}
	repositoryConfigFile := distribution.RepositoryConfigFile
	repositoryConfigFile.ReadContent()
	repositoryConfigWithinARSRepo :=
		repositoryConfigFile.CloneToDifferentLocation(ARSRepo.Config.RepositoryConfigFile.FilePath)
	repositoryConfigWithinARSRepo.Content.Local.SubsetPaths = PatchFiles
	repositoryConfigWithinARSRepo.Content.IndividualRepositoryPersist.UseSavedIndividualRepositories = true

	individualConfigFile :=
		filepath.Base(OriginRepo.DistributionMap[DistributionNameFlag].IndividualizationConfigFileName)
	repositoryConfigWithinARSRepo.Content.IndividualRepositoryPersist.SavedIndividualRepositoriesFileName =
		individualConfigFile
	repositoryConfigWithinARSRepo.Content.General.LocalMode = true
	repositoryConfigWithinARSRepo.Content.General.GlobalLogLevel = utils.LogLevelAsString()
	repositoryConfigWithinARSRepo.WriteContent()
}

func copySavedIndividualizationFileToARS() {
	log.Debug("patch.copySavedIndividualRepositoriesFileToARS()")
	err := utils.CopyFile(OriginRepo.DistributionMap[DistributionNameFlag].IndividualizationConfigFileName,
		ARSRepo.IndividualizationConfig.Dir)
	if err != nil {
		log.Fatalf("Error copying individualization file to %s: %v", ARSRepo.IndividualizationConfig.Dir, err)
	}
}

func runLocalGeneration() {
	log.Debug("patch.runLocalGeneration()")
	log.Info("Starting local generation of the individualized repositories...")
	// Store the original directory
	originalDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current directory: %v", err)
	}
	err = os.Chdir(ARSRepo.RepoDir)
	if err != nil {
		log.Fatalf("Error changing directory to %s: %v", ARSRepo, err)
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
