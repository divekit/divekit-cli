package runner

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

// RunNPMStartAlways unconditionally starts the npm start process for the given path.
func RunNPMStartAlways(dirPath, infoMsg string) error {
	logDirPathAndInfoMsg(dirPath, infoMsg)
	return executeCmd(dirPath)
}

// RunNPMStart initiates the npm start process for a specified path, and its execution is contingent on the
// DryRunFlag, determining whether the function should be skipped. The bool return value is set to true if the
// `npm start` operation was executed, regardless of whether it was successful or not.
func RunNPMStart(dirPath, infoMsg string) (bool, error) {
	logDirPathAndInfoMsg(dirPath, infoMsg)
	if DryRunFlag {
		log.Info("'Dry Run' flag set, therefore SKIP RUNNING 'npm start'.")
		return false, nil
	}
	return true, executeCmd(dirPath)
}

func executeCmd(dirPath string) error {
	// Run "npm start"
	cmd := exec.Command("npm", "start")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = dirPath

	if err := cmd.Run(); err != nil {
		log.Errorf("Error running 'npm start': %v", err)
		return err
	}

	return nil
}

func logDirPathAndInfoMsg(dirPath string, infoMsg string) {
	log.Debug("runner.RunNPMStartInDir(): dirPath = " + dirPath)
	log.Info(infoMsg + " by running 'npm start' in " + dirPath + ".")
}
