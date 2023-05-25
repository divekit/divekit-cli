package cmd

import (
	"divekit-cli/divekit/ars"
	"divekit-cli/divekit/origin"
	"divekit-cli/divekit/patch"
	"divekit-cli/utils"
	"fmt"
	"github.com/apex/log"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
)

var (
	// Flags
	DistributionNameFlag string
	OnlyGenerateFlag     bool
	// command state vars
	PatchFiles []string
	ARSRepo    *ars.ARSRepoType
	PatchRepo  *patch.PatchRepoType

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
	patchCmd.Flags().StringVarP(&DistributionNameFlag, "distribution", "d", "milestone",
		"name of the repo-distribution to patch")
	patchCmd.Flags().BoolVarP(&OnlyGenerateFlag, "only-generate", "1", false,
		"only generate the patch files locally, but don't move them to the student repos yet")

	patchCmd.MarkPersistentFlagRequired("originrepo")
	rootCmd.AddCommand(patchCmd)
}

func validateArgs(cmd *cobra.Command, args []string) error {
	log.Debug("subcmd.validateArgs()")
	var err error
	if len(args) == 0 {
		err = fmt.Errorf("You need to specify at least one filename to subcmd.")
	}
	return err
}

// Checks preconditions before running the command
func preRun(cmd *cobra.Command, args []string) {
	ARSRepo = ars.NewARSRepo()
	PatchRepo = patch.NewPatchRepo()

	distribution := origin.OriginRepo.GetDistribution(DistributionNameFlag)
	if distribution == nil {
		log.WithFields(log.Fields{
			"DistributionNameFlag": DistributionNameFlag,
		}).Fatal("Distribution not found")
	}

	if utils.DryRunFlag && OnlyGenerateFlag {
		log.Warn("Running in dry-run mode - this overrides the --just-generate flag")
		OnlyGenerateFlag = false
	}
}

func run(cmd *cobra.Command, args []string) {
	log.Debug("subcmd.run()")
	definePatchFiles(args)
	log.Info(fmt.Sprintf("Found files to patch:\n%s", strings.Join(PatchFiles, "\n")))

	setRepositoryConfigWithinARSRepo()
	copySavedIndividualizationFileToARS()
	utils.RunNPMStart(ARSRepo.RepoDir,
		"Starting local generation of the individualized repositories containing patch files")

	if !OnlyGenerateFlag {
		copyLocallyGeneratedFilesToPatchTool()
		distribution := origin.OriginRepo.GetDistribution(DistributionNameFlag)
		PatchRepo.UpdatePatchConfigFile(distribution.RepositoryConfigFile)
		utils.RunNPMStart(PatchRepo.RepoDir,
			"Actually patching the files to each repository")
	} else {
		log.Info("Skipping patching of the student repositories")
	}
}

func definePatchFiles(args []string) {
	log.Debug("subcmd.definePatchFiles()")
	srcDir := filepath.Join(origin.OriginRepo.RepoDir, "src")
	for index := range args {
		foundFiles, foundErr := utils.FindFilesInDir(args[index], origin.OriginRepo.RepoDir)
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
		relFile, err := utils.TransformIntoRelativePaths(origin.OriginRepo.RepoDir, foundFiles[0])
		log.Debug(fmt.Sprintf("... relative to origin repo: %s", relFile))
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s", err)
			os.Exit(1)
		}
		PatchFiles = append(PatchFiles, relFile)
	}
}

func setRepositoryConfigWithinARSRepo() {
	log.Debug("subcmd.setRepositoryConfigWithinARSRepo()")
	distribution := origin.OriginRepo.GetDistribution(DistributionNameFlag)
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
		filepath.Base(origin.OriginRepo.DistributionMap[DistributionNameFlag].IndividualizationConfigFileName)
	repositoryConfigWithinARSRepo.Content.IndividualRepositoryPersist.SavedIndividualRepositoriesFileName =
		individualConfigFile
	repositoryConfigWithinARSRepo.Content.General.LocalMode = true
	repositoryConfigWithinARSRepo.Content.General.GlobalLogLevel = utils.LogLevelAsString()
	repositoryConfigWithinARSRepo.WriteContent()
}

func copySavedIndividualizationFileToARS() {
	log.Debug("subcmd.copySavedIndividualRepositoriesFileToARS()")
	err := utils.CopyFile(origin.OriginRepo.DistributionMap[DistributionNameFlag].IndividualizationConfigFileName,
		ARSRepo.IndividualizationConfig.Dir)
	if err != nil {
		log.Fatalf("Error copying individualization file to %s: %v", ARSRepo.IndividualizationConfig.Dir, err)
	}
}

func copyLocallyGeneratedFilesToPatchTool() {
	log.Debug("subcmd.copyLocallyGeneratedFilesToPatchTool()")
	log.Info("Copying locally generated files to patch tool...")
	// Copy the generated files to the patch tool
	err := PatchRepo.CleanInputDir()
	if err == nil {
		err = utils.CopyAllFilesInDir(ARSRepo.GeneratedLocalOutput.Dir, PatchRepo.InputDir)
	}
	if err != nil {
		log.Fatalf("Error copying locally generated files to patch tool: %v", err)
		os.Exit(1)
	}
	log.Info("Copying completed.")
}
