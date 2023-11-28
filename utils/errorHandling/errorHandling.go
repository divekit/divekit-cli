package errorHandling

import (
	"bufio"
	"divekit-cli/utils/testUtils"
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
		outputErrorWithFuncName(err, msg)
	}

	if len(errorsList) > 0 {
		os.Exit(1)
	}
}

// Outputs an error to stderr, if there is one, and aborts the program if so
func OutputAndAbortIfError(err error, msg string) {
	log.Debug("errorHandling.OutputAndAbortIfError()")
	if err != nil {
		outputErrorWithFuncName(err, msg)
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

func OutputError(err error) {
	log.Debug("errorHandling.OutputError()")
	_, _ = fmt.Fprintln(os.Stderr, err)
}

func outputErrorWithFuncName(err error, msg string) {
	funcName := getPreviousFuncName()
	_, _ = fmt.Fprintf(os.Stderr, "[%s] %s: %v\n", funcName, msg, err)
}
func getPreviousFuncName() string {
	pc, _, _, _ := runtime.Caller(3)
	fullFuncName := runtime.FuncForPC(pc).Name()
	funcName := testUtils.GetBaseName(fullFuncName)
	return funcName
}
