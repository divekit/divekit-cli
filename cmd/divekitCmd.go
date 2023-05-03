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
	originRepoNameFlag string
	logLevelFlag       string
	divekitHomeFlag    string

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
	rootCmd.PersistentFlags().BoolVarP(&utils.AsIfFlag, "asif", "a", false,
		"just tell what you would do, but don't do it yet")
	rootCmd.PersistentFlags().StringVarP(&logLevelFlag, "loglevel", "l", "info",
		"log level for the output (warn, info [default], debug, error)")
	rootCmd.PersistentFlags().StringVarP(&originRepoNameFlag, "originrepo", "o", "",
		"name of the origin repo to work with")
	rootCmd.PersistentFlags().StringVarP(&divekitHomeFlag, "home", "m", "",
		"home directory of all the Divekit repos")
}

func persistentPreRun(cmd *cobra.Command, args []string) {
	utils.DefineLoggingLevel(logLevelFlag)
	log.Debug("divekit.persistentPreRun()")
	divekit.InitDivekitHomeDir(divekitHomeFlag)
	origin.InitOriginRepo(originRepoNameFlag)
}

func Execute() error {
	log.Debug("divekit.Execute()")
	return rootCmd.Execute()
}
