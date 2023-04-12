package utils

import (
	"fmt"
	"github.com/apex/log"
	"os"
)

func OutputAndAbortIfErrors(errorsList []error) {
	log.Debug("utils.OutputAndAbortIfErrors()")
	for _, err := range errorsList {
		_, _ = fmt.Fprintln(os.Stderr, "Error: ", err)
	}

	if len(errorsList) > 0 {
		os.Exit(1)
	}
}
