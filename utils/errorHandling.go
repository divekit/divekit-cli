package utils

import (
	"bufio"
	"fmt"
	"github.com/apex/log"
	"os"
	"strings"
)

// Outputs a list of errors to stderr, and aborts the program if there are any errors
func OutputAndAbortIfErrors(errorsList []error) {
	log.Debug("utils.OutputAndAbortIfErrors()")
	for _, err := range errorsList {
		_, _ = fmt.Fprintln(os.Stderr, "Error: ", err)
	}

	if len(errorsList) > 0 {
		os.Exit(1)
	}
}

// Outputs an error to stderr, if there is one, and aborts the program if so
func OutputAndAbortIfError(error error) {
	log.Debug("utils.OutputAndAbortIfError()")
	if error != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Error: ", error)
		os.Exit(1)
	}
}

// Asks the user to confirm an action, and aborts if the user doesn't confirm
func Confirm(prompt string) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("%s\n\n(Please type \"yes\" to confirm, or anything else to abort):\n", prompt)

	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Error reading input: %v\n", err)
		os.Exit(1)
	}

	input = strings.TrimSpace(strings.ToLower(input))
	if input != "yes" {
		fmt.Println("Aborting")
		os.Exit(1)
	}
}
