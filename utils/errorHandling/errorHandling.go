package errorHandling

import (
	"bufio"
	"fmt"
	"github.com/cockroachdb/errors"
	"github.com/stretchr/testify/assert"
	"os"
	"strings"
	"testing"
)

// Confirm asks the user to confirm an action. The return value contains true if the action get confirmed.
func Confirm(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("%s\n\n(Please type \"yes\" to confirm, or anything else to abort):\n", prompt)

	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("Error reading input: %v\n", err)
		return false
	}

	input = strings.TrimSpace(strings.ToLower(input))
	if input != "yes" {
		fmt.Println("input does not equal \"yes\"")
		return false
	}

	return true
}

// IsErrorType asserts that the cause of the actual error is of the expected type
func IsErrorType(t *testing.T, expected error, actual error) {
	actualCause := errors.Cause(actual)
	assert.IsTypef(t, expected, actualCause, "The cause of the actual error: %v", actualCause)
}
