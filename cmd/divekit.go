package cmd

import (
	"github.com/spf13/cobra"
)

var (
	// Flags
	asIf           bool
	verbose        bool
	originRepoName string

	rootCmd = &cobra.Command{
		Use:   "divekit",
		Short: "divekit helps to create and distribute individualized repos for software engineering exercises",
		Long: `divekit has been developed at TH Köln by the ArchiLab team (www.archi-lab.io) as  
universal tool to design, individualize, distribute, assess, patch, and evaluate 
realistic software engineering exercises as Git repos.`,
		SilenceErrors: true,
		SilenceUsage:  true,
	}
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&asIf, "asif", "a", false,
		"just tell what you would do, but don't do it yet")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false,
		"be extra chatty with all details of the execution")
	rootCmd.PersistentFlags().StringVarP(&originRepoName, "originrepo", "o", "",
		"name of the origin repo to work with")
}

func Execute() error {
	return rootCmd.Execute()
}
