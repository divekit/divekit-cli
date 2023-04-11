package utils

import (
	"fmt"
	"os"
)

func OutputAndAbortIfErrors(errorsList []error) {
	for _, err := range errorsList {
		_, _ = fmt.Fprintln(os.Stderr, "Error: ", err)
	}

	if len(errorsList) > 0 {
		os.Exit(1)
	}
}
