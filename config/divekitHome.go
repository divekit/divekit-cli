package config

import (
	"divekit-cli/cmd"
	"divekit-cli/utils"
	"github.com/apex/log"
	"os"
)

var (
	DivekitHomeDir string
)

// cmd.DivekitHomeFlag is the home directory of all the Divekit repos. It is set by the
// --home flag, the DIVEKIT_HOME environment variable, or the current working directory
// (in this order).
func InitDivekitHomeDir() {
	log.Debug("config.InitDivekitHomeDir()")
	setDivekitHomeDirFromVariousSources()
	utils.OutputAndAbortIfErrors(utils.ValidateAllFilePaths(DivekitHomeDir))
	log.WithFields(log.Fields{
		"DivekitHomeDir": DivekitHomeDir,
	}).Info("Setting Divekit Home Dir:")
}

func setDivekitHomeDirFromVariousSources() {
	if cmd.DivekitHomeFlag != "" {
		log.Info("Home dir is set via flag -m / --home: " + cmd.DivekitHomeFlag)
		DivekitHomeDir = cmd.DivekitHomeFlag
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
