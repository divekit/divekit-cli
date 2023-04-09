package cmd

import (
	"github.com/spf13/cobra"
)

var patchCmd = &cobra.Command{
	Use:   "patch",
	Short: "Apply a patch",
	Long:  `Apply a patch to the divekit.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Your patch command logic goes here
	},
}

func init() {
	// Add the patch subcommand to the divekit command
	rootCmd.AddCommand(patchCmd)
}
