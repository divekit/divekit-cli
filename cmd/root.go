package cmd

import (
	"divekit-cli/divekit"
	"divekit-cli/divekit/origin"
	"divekit-cli/utils/logUtils"
	"divekit-cli/utils/runner"
	"github.com/apex/log"
	"github.com/spf13/cobra"
)

var (
	OriginRepoNameFlag string
	LogLevelFlag       string
	DivekitHomeFlag    string

	RootCmd = NewRootCmd()
)

func init() {
	log.Debug("divekit.init()")
	SetCmdFlags(RootCmd)

	RootCmd.Version = "0.0.1" // todo: get version from git tag
}

func NewRootCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "divekit",
		Short: "divekit helps to create and distribute individualized repos for software engineering exercises",
		Long: `Divekit has been developed at TH Köln by the ArchiLab team (www.archi-lab.io) as
universal tool to design, individualize, distribute, assess, patch, and evaluate
realistic software engineering exercises as Git repos.`,
		PersistentPreRunE: persistentPreRun,
	}
}

func SetCmdFlags(cmd *cobra.Command) {
	log.Debug("divekit.SetCmdFlags()")
	cmd.PersistentFlags().BoolVarP(&runner.DryRunFlag, "dry-Run", "0", false,
		"just tell what you would do, but don't do it yet")
	cmd.PersistentFlags().StringVarP(&LogLevelFlag, "loglevel", "l", "info",
		"log level (warn, info, debug, error)")
	cmd.PersistentFlags().StringVarP(&OriginRepoNameFlag, "originrepo", "o", "",
		"name of the origin repo to work with")
	cmd.PersistentFlags().StringVarP(&DivekitHomeFlag, "home", "m", "",
		"home directory of all the Divekit repos")
}

func persistentPreRun(cmd *cobra.Command, args []string) error {
	log.Debug("divekit.persistentPreRun()")
	if err := logUtils.DefineLoggingLevel(LogLevelFlag); err != nil {
		log.Errorf("Could not define the logging level flag: %v", err)
		return err
	}

	if err := divekit.InitDivekitHomeDir(DivekitHomeFlag); err != nil {
		log.Errorf("Could not initialize the divekit home flag: %v", err)
		return err
	}

	if err := origin.InitOriginRepo(OriginRepoNameFlag); err != nil {
		log.Errorf("Could not initialize the origin repository: %v", err)
		return err
	}

	return nil
}

func Execute() error {
	log.Debug("divekit.Execute()")
	return RootCmd.Execute()
}

type InvalidArgsError struct {
	Msg string
}

func (e *InvalidArgsError) Error() string {
	return e.Msg
}
