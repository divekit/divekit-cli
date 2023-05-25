package cmd

import (
	"divekit-cli/divekit"
	"divekit-cli/divekit/origin"
	"divekit-cli/utils"
	"github.com/apex/log"
	"github.com/spf13/cobra"
)

var (
	// Flags
	OriginRepoNameFlag string
	LogLevelFlag       string
	DivekitHomeFlag    string

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
	rootCmd.PersistentFlags().BoolVarP(&utils.DryRunFlag, "dry-run", "0", false,
		"just tell what you would do, but don't do it yet")
	rootCmd.PersistentFlags().StringVarP(&LogLevelFlag, "loglevel", "l", "info",
		"log level (warn, info, debug, error)")
	rootCmd.PersistentFlags().StringVarP(&OriginRepoNameFlag, "originrepo", "o", "",
		"name of the origin repo to work with")
	rootCmd.PersistentFlags().StringVarP(&DivekitHomeFlag, "home", "m", "",
		"home directory of all the Divekit repos")
}

func persistentPreRun(cmd *cobra.Command, args []string) {
	utils.DefineLoggingLevel(LogLevelFlag)
	log.Debug("divekit.persistentPreRun()")
	divekit.InitDivekitHomeDir(DivekitHomeFlag)
	origin.InitOriginRepo(OriginRepoNameFlag)
}

func Execute() error {
	log.Debug("divekit.Execute()")
	return rootCmd.Execute()
}
