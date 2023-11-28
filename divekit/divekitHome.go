package divekit

import (
	"divekit-cli/utils/errorHandling"
	"divekit-cli/utils/fileUtils"
	"github.com/apex/log"
	"os"
)

// Global vars
var (
	DivekitHomeDir string
)

// subcmd.DivekitHomeFlag is the home directory of all the Divekit repos. It is set by the
// --home flag, the DIVEKIT_HOME environment variable, or the current working directory
// (in this order).
func InitDivekitHomeDir(divekitHomeFlag string) {
	log.Debug("config.InitDivekitHomeDir()")
	setDivekitHomeDirFromVariousSources(divekitHomeFlag)
	errorHandling.OutputAndAbortIfErrors(fileUtils.ValidateAllDirPaths(DivekitHomeDir),
		"Could not initialize divekitHomeDir")
	log.WithFields(log.Fields{
		"DivekitHomeDir": DivekitHomeDir,
	}).Info("Setting Divekit Home Dir:")
}

func setDivekitHomeDirFromVariousSources(divekitHomeFlag string) {
	if divekitHomeFlag != "" {
		log.Info("Home dir is set via flag -m / --home: " + divekitHomeFlag)
		DivekitHomeDir = divekitHomeFlag
		return
	}
	envHome := os.Getenv("DIVEKIT_HOME")
	if envHome != "" {
		log.Info("Home dir is set via DIVEKIT_HOME environment variable: " + envHome)
		DivekitHomeDir = envHome
		return
	}
	workingDir, _ := os.Getwd()
	log.Info("Home dir set to current directory: " + workingDir)
	DivekitHomeDir = workingDir
	return
}
