package main

// NOTE: Each subcommand (e.g. patch) must be imported with an underscore to trigger it's init() function, which
// registers the subcommand with the root command.
import (
	"divekit-cli/cmd"
	_ "divekit-cli/cmd/patch"
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
