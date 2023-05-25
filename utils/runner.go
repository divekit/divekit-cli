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
	DryRunFlag bool // If true, then don't actually run the commands, just output what would be run
)

func RunNPMStartAlways(dirPath, infoMsg string) error {
	return runNPMStartWithDryRunCheck(dirPath, infoMsg, false)
}

func RunNPMStart(dirPath, infoMsg string) error {
	return runNPMStartWithDryRunCheck(dirPath, infoMsg, true)
}

func runNPMStartWithDryRunCheck(dirPath, infoMsg string, skipIfDryRun bool) error {
	log.Debug("utils.RunNPMStartInDir(): dirPath = " + dirPath)
	log.Info(infoMsg + " by running 'npm start' in " + dirPath + ".")
	if skipIfDryRun && DryRunFlag {
		log.Info("'Dry Run' flag set, therefore SKIP RUNNING 'npm start'.")
		return nil
	}

	// Run "npm start"
	cmd := exec.Command("npm", "start")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = dirPath
	err := cmd.Run()
	if err != nil {
		log.Fatalf("Error running 'npm start': %v", err)
	}

	return err
}
