package cmd

import (
	"divekit-cli/utils"
	"github.com/apex/log"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

const (
	// ARS
	ARS_REPO_NAME                  = "divekit-automated-repo-setup"
	ARS_REPOSITORY_CONFIG_FILENAME = "repositoryConfig.json"
	ARS_CONFIG_DIR_NAME            = "resources\\config"
	ARS_INDIVIDUAL_REPOSITORY_NAME = "resources\\individual_repositories"

	// Origin repo
	DIVEKIT_DIR_NAME       = ".divekit_norepo"
	DISTRIBUTIONS_DIR_NAME = "distributions"
)

var (
	// Flags
	AsIfFlag       bool
	VerboseFlag    bool
	DebugFlag      bool
	OriginRepoName string
	DivekitHome    string

	// global vars
	OriginRepoFullPath              string
	DivekitDirFullPath              string
	DistributionsRootDirFullPath    string
	ARSRepoFullPath                 string
	ARSRepositoryConfigFileFullPath string
	ARSIndividualRepoFullPath       string

	rootCmd = &cobra.Command{
		Use:   "divekit",
		Short: "divekit helps to create and distribute individualized repos for software engineering exercises",
		Long: `Divekit has been developed at TH Köln by the ArchiLab team (www.archi-lab.io) as
universal tool to design, individualize, distribute, assess, patch, and evaluate
realistic software engineering exercises as Git repos.`,
		PersistentPreRun: persistentPreRun,
		SilenceErrors:    true,
		SilenceUsage:     true,
	}
)

func init() {
	log.Debug("divekit.init()")
	rootCmd.PersistentFlags().BoolVarP(&AsIfFlag, "asif", "a", false,
		"just tell what you would do, but don't do it yet")
	rootCmd.PersistentFlags().BoolVarP(&VerboseFlag, "verbose", "v", false,
		"be extra chatty with all details of the execution")
	rootCmd.PersistentFlags().BoolVarP(&DebugFlag, "debug", "g", false,
		"debug mode, printing all debug and trace messages")
	rootCmd.PersistentFlags().StringVarP(&OriginRepoName, "originrepo", "o", "",
		"name of the origin repo to work with")
	rootCmd.PersistentFlags().StringVarP(&DivekitHome, "home", "m", "",
		"home directory of all the Divekit repos")
}

func persistentPreRun(cmd *cobra.Command, args []string) {
	utils.DefineLoggingConfig(VerboseFlag, DebugFlag)
	log.Debug("divekit.persistentPreRun()")

	DivekitHome = getHomeDir()
	OriginRepoFullPath = filepath.Join(DivekitHome, OriginRepoName)
	DivekitDirFullPath = filepath.Join(OriginRepoFullPath, DIVEKIT_DIR_NAME)
	DistributionsRootDirFullPath = filepath.Join(DivekitDirFullPath, DISTRIBUTIONS_DIR_NAME)
	ARSRepoFullPath = filepath.Join(DivekitHome, ARS_REPO_NAME)
	ARSRepositoryConfigFileFullPath =
		filepath.Join(ARSRepoFullPath, ARS_CONFIG_DIR_NAME, ARS_REPOSITORY_CONFIG_FILENAME)
	ARSIndividualRepoFullPath = filepath.Join(ARSRepoFullPath, ARS_INDIVIDUAL_REPOSITORY_NAME)

	utils.OutputAndAbortIfErrors(
		utils.ValidateAllDirPaths(OriginRepoFullPath, DivekitDirFullPath, DistributionsRootDirFullPath, ARSRepoFullPath,
			ARSIndividualRepoFullPath))
	utils.OutputAndAbortIfErrors(
		utils.ValidateAllFilePaths(ARSRepositoryConfigFileFullPath))

	log.WithFields(log.Fields{
		"OriginRepoFullPath":              OriginRepoFullPath,
		"DivekitDirFullPath":              DivekitDirFullPath,
		"DistributionsRootDirFullPath":    DistributionsRootDirFullPath,
		"ARSRepoFullPath":                 ARSRepoFullPath,
		"ARSRepositoryConfigFileFullPath": ARSRepositoryConfigFileFullPath,
		"ARSIndividualRepoFullPath":       ARSIndividualRepoFullPath,
	}).Info("Setting global variables:")
}

// DivekitHome is the home directory of all the Divekit repos. It is set by the
// --home flag, the DIVEKIT_HOME environment variable, or the current working directory
// (in this order).
func getHomeDir() string {
	if DivekitHome != "" {
		log.Info("Home dir is set via flag -m / --home: " + DivekitHome)
		return DivekitHome
	}
	envHome := os.Getenv("DIVEKIT_HOME")
	if envHome != "" {
		log.Info("Home dir is set via DIVEKIT_HOME environment variable: " + envHome)
		return envHome
	}
	workingDir, _ := os.Getwd()
	log.Info("Home dir set to current directory: " + workingDir)
	return workingDir
}

func Execute() error {
	log.Debug("divekit.Execute()")
	return rootCmd.Execute()
}
