package main

import (
	//"github.com/divekit/divekit-cli"
	"divekit-cli/cmd"
	"os"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
