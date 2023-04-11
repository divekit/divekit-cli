package cmd

import (
	"divekit-cli/config"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

var (
	// Flags
	AsIf           bool
	Verbose        bool
	OriginRepoName string
	DivekitHome    string

	// global vars
	OriginRepoFullPath       string
	DivekitDirFullPath       string
	DistributionsDirFullPath string

	rootCmd = &cobra.Command{
		Use:   "divekit",
		Short: "divekit helps to create and distribute individualized repos for software engineering exercises",
		Long: `Divekit has been developed at TH Köln by the ArchiLab team (www.archi-lab.io) as
universal tool to design, individualize, distribute, assess, patch, and evaluate
realistic software engineering exercises as Git repos.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			DivekitHome = getHomeDir()
			fmt.Println("Home directory:", DivekitHome)
			OriginRepoFullPath = filepath.Join(DivekitHome, OriginRepoName)
			DivekitDirFullPath = filepath.Join(OriginRepoFullPath, config.DIVEKIT_DIR_NAME)
			DistributionsDirFullPath = filepath.Join(DivekitDirFullPath, config.DISTRIBUTIONS_DIR_NAME)
		},
		SilenceErrors: true,
		SilenceUsage:  true,
	}
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&AsIf, "asif", "a", false,
		"just tell what you would do, but don't do it yet")
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false,
		"be extra chatty with all details of the execution")
	rootCmd.PersistentFlags().StringVarP(&OriginRepoName, "originrepo", "o", "",
		"name of the origin repo to work with")
	rootCmd.PersistentFlags().StringVarP(&DivekitHome, "home", "m", "",
		"home directory of all the Divekit repos")
}

// DivekitHome is the home directory of all the Divekit repos. It is set by the
// --home flag, the DIVEKIT_HOME environment variable, or the current working directory
// (in this order).
func getHomeDir() string {
	if DivekitHome != "" {
		return DivekitHome
	}
	envHome := os.Getenv("DIVEKIT_HOME")
	if envHome != "" {
		return envHome
	}
	workingDir, _ := os.Getwd()
	return workingDir
}

func ValidateOriginRepo() error {
	_, err := os.Stat(OriginRepoFullPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("Origin repo does not exist (have you used the -o / --originrepo flag?): %s", OriginRepoFullPath)
	}
	_, err = os.Stat(DivekitDirFullPath)
	if os.IsNotExist(err) {
		return fmt.Errorf(".divekit subfolder not found in: %s", DivekitDirFullPath)
	}
	return nil
}

func Execute() error {
	return rootCmd.Execute()
}
