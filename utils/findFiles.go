package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

// Searches recursively for full path(es) of a given filename. Returns a 1-elem
// array if there is just one occurrence, or an array with several elements otherwise.
func FindFiles(justTheFileName, rootDir string) ([]string, error) {
	var targetPaths []string

	err := filepath.Walk(rootDir, func(currentPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if the current file matches the target file name
		if info.Mode().IsRegular() && info.Name() == justTheFileName {
			targetPaths = append(targetPaths, currentPath)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("Error searching for files: %v", err)
	}
	return targetPaths, nil
}
