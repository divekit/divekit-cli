package main

import (
	"divekit-cli/cmd"
	_ "divekit-cli/cmd/patch" // NOTE: Each subcommand must include this import to be acknowledged as a valid subcommand.
	"fmt"
	"github.com/apex/log"
	"os"
)

func main() {
	log.Debug("main()")
	err := cmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
