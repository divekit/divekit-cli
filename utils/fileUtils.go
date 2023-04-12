package utils

import (
	"fmt"
	"github.com/apex/log"
	"os"
	"path/filepath"
)

// Searches recursively for full path(es) of a given filename. Returns a 1-elem
// array if there is just one occurrence, or an array with several elements otherwise.
func FindFiles(justTheFileName, rootDir string) ([]string, error) {
	log.Debug("utils.FindFiles() - justTheFileName: " + justTheFileName + ", rootDir: " + rootDir)
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

// Check if a given file exists (and is a file)
func ValidateFilePath(filePath string) error {
	log.Debug("utils.ValidateFilePath()")
	return ValidatePath(false, filePath)
}

// Check if a given directory exists (and is a directory)
func ValidateDirPath(dirPath string) error {
	log.Debug("utils.ValidateDirPath()")
	return ValidatePath(true, dirPath)
}

// Check if a given path exists (and is a directory or a file, depending on shouldBeDir)
func ValidatePath(shouldBeDir bool, path string) error {
	log.Debug("utils.ValidatePath()")
	fileInfo, err := os.Stat(path)
	if os.IsNotExist(err) {
		return fmt.Errorf("%s does not exist", path)
	}
	if err != nil {
		return fmt.Errorf("Error checking for %s", err)
	}
	if shouldBeDir != fileInfo.Mode().IsDir() {
		errorMessage := "%s is a "
		if shouldBeDir {
			errorMessage += "file, not a directory"
		} else {
			errorMessage += "directory, not a file"
		}
		return fmt.Errorf(errorMessage, path)
	}

	return nil
}

// Validate a list of files
func ValidateAllFilePaths(paths ...string) []error {
	log.Debug("utils.ValidateAllFilePaths()")
	return ValidateAllPaths(false, paths...)
}

// Validate a list of directories
func ValidateAllDirPaths(paths ...string) []error {
	log.Debug("utils.ValidateAllDirPaths()")
	return ValidateAllPaths(true, paths...)
}

// Validate a list of paths (files or directories)
func ValidateAllPaths(shouldBeDir bool, paths ...string) []error {
	log.Debug("utils.ValidateAllPaths()")
	var errorsList []error
	for _, path := range paths {
		err := ValidatePath(shouldBeDir, path)
		if err != nil {
			errorsList = append(errorsList, err)
		}
	}
	return errorsList
}
