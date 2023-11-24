package testUtils

import (
	"log"
	"os"
)

// CreateFile generates a file at a specified path, returning the file path. If desired, the file can be generated
// with content provided as an argument.
func CreateFile(path string, fileName string, fileContent string) {
	f, err := os.Create(path + "/" + fileName)
	if err != nil {
		log.Fatalf("Could not create a file: %v", err)
	}

	_, err = f.Write([]byte(fileContent))
	if err != nil {
		log.Fatalf("Could not write to file: %v", err)
	}

	err = f.Close()
	if err != nil {
		log.Fatalf("Could not close file: %v", err)
	}
}

// CreateTmpDir creates a directory in the temp folder and provides its path as a
// return value. It is the caller's responsibility to remove this folder when it is no longer needed.
func CreateTmpDir() string {
	path, err := os.MkdirTemp("", "divekit_cli_")
	if err != nil {
		log.Fatalf("Could not create directory: %v", err)
	}

	return path
}

// CreateDir creates a directory named path, along with any necessary parents. If path is already a directory,
// CreateDir does nothing.
func CreateDir(path string) {
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		log.Fatalf("Could not create directory: %v", err)
	}
}

// DeleteDir deletes a specified directory along with its files and subdirectories.
func DeleteDir(path string) {
	if err := os.RemoveAll(path); err != nil {
		log.Fatalf("Could not remove directory: %v", err)
	}
}
