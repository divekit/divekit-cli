package cmd

import (
	"divekit-cli/config"
	"divekit-cli/utils"
	"github.com/apex/log"
	"github.com/spf13/cobra"
)

var (
	// Flags
	AsIfFlag           bool
	VerboseFlag        bool
	DebugFlag          bool
	OriginRepoNameFlag string
	DivekitHomeFlag    string

	// global vars
	OriginRepo *config.OriginRepoType

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
	rootCmd.PersistentFlags().StringVarP(&OriginRepoNameFlag, "originrepo", "o", "",
		"name of the origin repo to work with")
	rootCmd.PersistentFlags().StringVarP(&DivekitHomeFlag, "home", "m", "",
		"home directory of all the Divekit repos")
}

func persistentPreRun(cmd *cobra.Command, args []string) {
	utils.DefineLoggingConfig(VerboseFlag, DebugFlag)
	log.Debug("divekit.persistentPreRun()")
	config.InitDivekitHomeDir()
	if OriginRepoNameFlag != "" {
		OriginRepo = config.OriginRepo(OriginRepoNameFlag)
	}
}

func Execute() error {
	log.Debug("divekit.Execute()")
	return rootCmd.Execute()
}
