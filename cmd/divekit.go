package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "divekit",
	Short: "A CLI tool for divekit operations (tbd)",
	Long:  "tbd",
}

func Execute() error {
	return rootCmd.Execute()
}
