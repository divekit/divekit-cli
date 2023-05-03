package utils

/**
 * This file contains utility functions for running external commands.
 */

import (
	"github.com/apex/log"
	"os"
	"os/exec"
)

// Global flags
var (
	AsIfFlag bool // If true, then don't actually run the commands, just output what would be run
)

func RunNPMStart(dirPath, infoMsg string) error {
	log.Debug("utils.RunNPMStartInDir(): dirPath = " + dirPath)
	log.Info(infoMsg + " by running 'npm start' in " + dirPath + ".")
	if AsIfFlag {
		log.Info("'As if' flag set, therefore SKIP RUNNING 'npm start'.")
		return nil
	}
	// Store the original directory
	originalDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current directory: %v", err)
	}
	err = os.Chdir(dirPath)
	if err != nil {
		log.Fatalf("Error changing directory to %s: %v", dirPath, err)
	}

	// Run "npm start"
	cmd := exec.Command("npm", "start")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Fatalf("Error running 'npm start': %v", err)
	}

	// Change back to the original directory
	err = os.Chdir(originalDir)
	if err != nil {
		log.Fatalf("Error changing back to the original directory: %v", err)
	}
	return err
}
