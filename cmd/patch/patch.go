package patch

import (
	"divekit-cli/cmd"
	"divekit-cli/divekit/ars"
	"divekit-cli/divekit/origin"
	"divekit-cli/divekit/patch"
	"divekit-cli/utils/fileUtils"
	"divekit-cli/utils/logUtils"
	"divekit-cli/utils/runner"
	"fmt"
	"github.com/apex/log"
	"github.com/spf13/cobra"
	"path/filepath"
	"strings"
)

var (
	// Flags
	DistributionNameFlag string
	// command state vars
	PatchFiles []string
	ARSRepo    *ars.ARSRepoType
	PatchRepo  *patch.PatchRepoType

	patchCmd = NewPatchCmd()
)

func init() {
	log.Debug("patch.init()")
	initFlags()
	setCmdFlags(patchCmd)
	cmd.RootCmd.AddCommand(patchCmd)
}

func initFlags() {
	DistributionNameFlag = ""
	PatchFiles = []string{}
	ARSRepo = nil
	PatchRepo = nil
}

func NewPatchCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "patch",
		Short:   "Apply a patch to all repos",
		Long:    `Patch one or several files in all the repos of a certain distribution of the origin repo`,
		Args:    validateArgs,
		PreRunE: preRun,
		RunE:    run,
	}
}

func setCmdFlags(cmd *cobra.Command) {
	log.Debug("patch.setCmdFlags()")
	cmd.Flags().StringVarP(&DistributionNameFlag, "distribution", "d", "milestone",
		"name of the repo-distribution to patch")

	cmd.MarkPersistentFlagRequired("originrepo")
}

func validateArgs(c *cobra.Command, args []string) error {
	log.Debug("patch.validateArgs()")
	if len(args) == 0 {
		return &cmd.InvalidArgsError{"You need to specify at least one filename to subcmd"}
	}

	return nil
}

// Checks preconditions before running the command
func preRun(cmd *cobra.Command, args []string) error {
	log.Debug("patch.preRun()")
	var err error
	if ARSRepo, err = ars.NewARSRepo(); err != nil {
		log.Errorf("Could not initialize the ARS repository:", err)
		return err
	}

	if PatchRepo, err = patch.NewPatchRepo(); err != nil {
		log.Errorf("Could not initialize the patch repository:", err)
		return err
	}

	if _, err = origin.OriginRepo.GetDistribution(DistributionNameFlag); err != nil {
		log.Errorf("Could not find any distribution:", err)
		return err
	}

	return nil
}

func run(cmd *cobra.Command, args []string) error {
	log.Debug("patch.run()")
	var err error
	var distribution *origin.Distribution

	if err = definePatchFiles(args); err != nil {
		log.Errorf("Could not define patch files:", err)
		return err
	}

	log.Info(fmt.Sprintf("Found files to patch:\n%s", strings.Join(PatchFiles, "\n")))

	if err = setRepositoryConfigWithinARSRepo(); err != nil {
		log.Errorf("Could not set the repository config within the ARS repository:", err)
		return err
	}

	if err = copySavedIndividualizationFileToARS(); err != nil {
		log.Errorf("Could not copy the saved individualization file to ARS:", err)
		return err
	}

	if err = runner.RunNPMStartAlways(ARSRepo.RepoDir, "Starting local generation of the individualized "+
		"repositories containing patch files"); err != nil {
		log.Errorf("Could not run the ARS repository:", err)
		return err
	}

	if err = copyLocallyGeneratedFilesToPatchTool(); err != nil {
		log.Errorf("Could not copy the locally generated files to the patch tool:", err)
		return err
	}

	if distribution, err = origin.OriginRepo.GetDistribution(DistributionNameFlag); err != nil {
		log.Errorf("Could not find any distribution:", err)
		return err
	}

	if err = PatchRepo.UpdatePatchConfigFile(distribution.RepositoryConfigFile); err != nil {
		log.Errorf("Could not update the patch config file:", err)
		return err
	}

	if _, err = runner.RunNPMStart(PatchRepo.RepoDir, "Actually patching the files to each repository"); err != nil {
		log.Errorf("Could not run the patch repository:", err)
		return err
	}

	return nil
}

func definePatchFiles(args []string) error {
	log.Debug("subcmd.definePatchFiles()")
	srcDir := filepath.Join(origin.OriginRepo.RepoDir, "src")
	for index := range args {
		println("args[index]:", args[index])
		foundFiles, foundErr := fileUtils.FindFilesInDir(origin.OriginRepo.RepoDir, args[index])
		foundFiles2, foundErr2 := fileUtils.FindFilesInDirRecursively(srcDir, args[index])
		foundFiles = append(foundFiles, foundFiles2...)
		if foundErr != nil || foundErr2 != nil {
			return fmt.Errorf("failed to find files: %w, %w", foundErr, foundErr2)
		}

		if len(foundFiles) == 0 {
			return fmt.Errorf("No files found with name %s %w", args[index], &PatchFileError{})
		}

		if len(foundFiles) > 1 {
			errorMsg := "Error: Multiple files found:\n"
			for _, file := range foundFiles {
				errorMsg += fmt.Sprintf("  - %s\n", file)
			}
			return fmt.Errorf("%s %w", errorMsg, &PatchFileError{})
		}

		log.Debug(fmt.Sprintf("Found file %s", foundFiles[0]))
		relFile, err := fileUtils.TransformIntoRelativePaths(origin.OriginRepo.RepoDir, foundFiles[0])
		log.Debug(fmt.Sprintf("... relative to origin repo: %s", relFile))
		if err != nil {
			return fmt.Errorf("failed to transform into relative paths: %w", err)
		}

		PatchFiles = append(PatchFiles, relFile)
	}

	return nil
}

func setRepositoryConfigWithinARSRepo() error {
	log.Debug("patch.setRepositoryConfigWithinARSRepo()")
	var err error
	var distribution *origin.Distribution
	var repositoryConfigWithinARSRepo *ars.RepositoryConfigFileType

	if distribution, err = origin.OriginRepo.GetDistribution(DistributionNameFlag); err != nil {
		return err
	}

	repositoryConfigFile := distribution.RepositoryConfigFile
	if err = repositoryConfigFile.ReadContent(); err != nil {
		return err
	}

	if repositoryConfigWithinARSRepo, err =
		repositoryConfigFile.CloneToDifferentLocation(ARSRepo.Config.RepositoryConfigFile.FilePath); err != nil {
		return err
	}

	repositoryConfigWithinARSRepo.Content.Local.SubsetPaths = PatchFiles
	repositoryConfigWithinARSRepo.Content.IndividualRepositoryPersist.UseSavedIndividualRepositories = true
	individualConfigFile := filepath.Base(origin.OriginRepo.DistributionMap[DistributionNameFlag].IndividualizationConfigFileName)
	repositoryConfigWithinARSRepo.Content.IndividualRepositoryPersist.SavedIndividualRepositoriesFileName = individualConfigFile
	repositoryConfigWithinARSRepo.Content.General.LocalMode = true
	repositoryConfigWithinARSRepo.Content.General.GlobalLogLevel = logUtils.LogLevelAsString()

	if err = repositoryConfigWithinARSRepo.WriteContent(); err != nil {
		return err
	}

	return nil
}

func copySavedIndividualizationFileToARS() error {
	log.Debug("patch.copySavedIndividualRepositoriesFileToARS()")
	if err := fileUtils.CopyFile(origin.OriginRepo.DistributionMap[DistributionNameFlag].IndividualizationConfigFileName,
		ARSRepo.IndividualizationConfig.Dir); err != nil {
		return err
	}

	return nil
}

func copyLocallyGeneratedFilesToPatchTool() error {
	log.Debug("patch.copyLocallyGeneratedFilesToPatchTool()")
	log.Info("Copying locally generated files to patch tool...")
	// Copy the generated files to the patch tool
	if err := PatchRepo.CleanInputDir(); err != nil {
		return err
	}

	if err := fileUtils.CopyAllFilesInDir(ARSRepo.GeneratedLocalOutput.Dir, PatchRepo.InputDir); err != nil {
		return err
	}

	log.Info("Copying completed.")

	return nil
}

type PatchFileError struct {
	Msg string
}

func (e *PatchFileError) Error() string {
	return e.Msg
}
