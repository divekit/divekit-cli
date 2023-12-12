package errorHandling

import (
	"bufio"
	"divekit-cli/utils/fileUtils"
	"fmt"
	"github.com/apex/log"
	"os"
	"runtime"
	"strings"
)

// Outputs a list of errors to stderr, and aborts the program if there are any errors
func OutputAndAbortIfErrors(errorsList []error, msg string) {
	log.Debug("errorHandling.OutputAndAbortIfErrors()")
	for _, err := range errorsList {
		outputWithErrorLocation(err, msg)
	}

	if len(errorsList) > 0 {
		os.Exit(1)
	}
}

// Outputs an error to stderr, if there is one, and aborts the program if so
func OutputAndAbortIfError(err error, msg string) {
	log.Debug("errorHandling.OutputAndAbortIfError()")
	if err != nil {
		outputWithErrorLocation(err, msg)
		os.Exit(1)
	}
}

// Confirm asks the user to confirm an action, and aborts if the user doesn't confirm
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

func outputWithErrorLocation(err error, msg string) {
	location := getErrorLocation()
	_, _ = fmt.Fprintf(os.Stderr, "Error at %s %s: %v\n", location, msg, err)
}

// getErrorLocation gets a filename along with the line where the error was triggered
func getErrorLocation() string {
	pc, _, _, _ := runtime.Caller(3)
	file, line := runtime.FuncForPC(pc).FileLine(pc)
	src := fileUtils.GetBaseName(fmt.Sprintf("%s:%v", file, line))
	return src
}
