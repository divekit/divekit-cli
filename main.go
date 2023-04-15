package main

import (
	//"github.com/divekit/divekit-cli"
	"divekit-cli/cmd"
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
