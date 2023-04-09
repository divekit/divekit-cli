package main

import (
	//"github.com/divekit/divekit-cli"
	"divekit-cli/cmd"
	"fmt"
	"os"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
