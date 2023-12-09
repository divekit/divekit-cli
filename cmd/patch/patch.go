package patch

import (
	"divekit-cli/cmd"
	"divekit-cli/divekit/ars"
	"divekit-cli/divekit/origin"
	"divekit-cli/divekit/patch"
	"divekit-cli/utils/errorHandling"
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
	setCmdFlags(patchCmd)
	cmd.RootCmd.AddCommand(patchCmd)
}

func NewPatchCmd() *cobra.Command {
	return &cobra.Command{
		Use:    "patch",
		Short:  "Apply a patch to all repos",
		Long:   `Patch one or several files in all the repos of a certain distribution of the origin repo`,
		Args:   validateArgs,
		PreRun: preRun,
		Run:    run,
	}
}

func setCmdFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&DistributionNameFlag, "distribution", "d", "milestone",
		"name of the repo-distribution to patch")

	cmd.MarkPersistentFlagRequired("originrepo")
}

func validateArgs(cmd *cobra.Command, args []string) error {
	log.Debug("subcmd.validateArgs()")
	var err error
	if len(args) == 0 {
		err = fmt.Errorf("You need to specify at least one filename to subcmd")
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
		})
		errorHandling.OutputAndAbortIfError(fmt.Errorf("distribution not found"),
			"Could not prepare the patch command")
	}
}

func run(cmd *cobra.Command, args []string) {
	log.Debug("subcmd.Run()")
	definePatchFiles(args)
	log.Info(fmt.Sprintf("Found files to patch:\n%s", strings.Join(PatchFiles, "\n")))

	setRepositoryConfigWithinARSRepo()
	copySavedIndividualizationFileToARS()
	err := runner.RunNPMStartAlways(ARSRepo.RepoDir,
		"Starting local generation of the individualized repositories containing patch files")
	errorHandling.OutputAndAbortIfError(err, "Could run the ARS repository")

	copyLocallyGeneratedFilesToPatchTool()
	distribution := origin.OriginRepo.GetDistribution(DistributionNameFlag)
	PatchRepo.UpdatePatchConfigFile(distribution.RepositoryConfigFile)
	_, err = runner.RunNPMStart(PatchRepo.RepoDir, "Actually patching the files to each repository")
	errorHandling.OutputAndAbortIfError(err, "Could not run the repo editor repository")
}

func definePatchFiles(args []string) {
	log.Debug("subcmd.definePatchFiles()")
	srcDir := filepath.Join(origin.OriginRepo.RepoDir, "src")
	for index := range args {
		println("args[index]:", args[index])
		foundFiles, foundErr := fileUtils.FindFilesInDir(args[index], origin.OriginRepo.RepoDir)
		foundFiles2, foundErr2 := fileUtils.FindFilesInDirRecursively(args[index], srcDir)
		foundFiles = append(foundFiles, foundFiles2...)
		if foundErr != nil || foundErr2 != nil {
			errorHandling.OutputAndAbortIfError(fmt.Errorf(fmt.Sprintf("%s %s", foundErr, foundErr2)), "-")
		}
		if len(foundFiles) == 0 {
			errorHandling.OutputAndAbortIfError(fmt.Errorf("No files found with name "+args[index]), "-")
		}
		if len(foundFiles) > 1 {
			errorMsg := "Error: Multiple files found:\n"
			for _, file := range foundFiles {
				errorMsg += fmt.Sprintf("  - %s\n", file)
			}
			errorHandling.OutputAndAbortIfError(fmt.Errorf(errorMsg), "-")
		}
		log.Debug(fmt.Sprintf("Found file %s", foundFiles[0]))
		relFile, err := fileUtils.TransformIntoRelativePaths(origin.OriginRepo.RepoDir, foundFiles[0])
		log.Debug(fmt.Sprintf("... relative to origin repo: %s", relFile))
		if err != nil {
			errorHandling.OutputAndAbortIfError(err, "Could not transform into relative paths")
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
		})
		errorHandling.OutputAndAbortIfError(fmt.Errorf("distribution not found"),
			"Could not set the repository config within the ars repo")
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
	repositoryConfigWithinARSRepo.Content.General.GlobalLogLevel = logUtils.LogLevelAsString()
	repositoryConfigWithinARSRepo.WriteContent()
}

func copySavedIndividualizationFileToARS() {
	log.Debug("subcmd.copySavedIndividualRepositoriesFileToARS()")
	errorHandling.OutputAndAbortIfError(
		fileUtils.CopyFile(origin.OriginRepo.DistributionMap[DistributionNameFlag].IndividualizationConfigFileName,
			ARSRepo.IndividualizationConfig.Dir), "Error copying individualization file to "+ARSRepo.IndividualizationConfig.Dir)
}

func copyLocallyGeneratedFilesToPatchTool() {
	log.Debug("subcmd.copyLocallyGeneratedFilesToPatchTool()")
	log.Info("Copying locally generated files to patch tool...")
	// Copy the generated files to the patch tool
	err := PatchRepo.CleanInputDir()
	if err == nil {
		err = fileUtils.CopyAllFilesInDir(ARSRepo.GeneratedLocalOutput.Dir, PatchRepo.InputDir)
	}
	errorHandling.OutputAndAbortIfError(err, "Error copying locally generated files to patch tool")
	log.Info("Copying completed.")
}
