package ars

/**
 * This file an "object-oriented lookalike" implementation for the structure of the ARS repository.
 */

import (
	"divekit-cli/divekit"
	"divekit-cli/utils"
	"github.com/apex/log"
	"path/filepath"
)

// all the paths used in the ARS repository (all as full paths)
type ARSRepoType struct {
	RepoDir string
	Config  struct {
		Dir                  string
		RepositoryConfigFile *RepositoryConfigFileType
	}
	IndividualizationConfig struct {
		Dir      string
		FileName string
	}
	GeneratedOverviewFiles struct {
		Dir string
	}
	GeneratedLocalOutput struct {
		Dir string
	}
}

// This method is similar to a constructor in OOP
func NewARSRepo() *ARSRepoType {
	log.Debug("ars.NewARSRepo()")
	arsRepo := &ARSRepoType{}
	arsRepo.RepoDir = filepath.Join(divekit.DivekitHomeDir, "divekit-automated-repo-setup")
	arsRepo.Config.Dir = filepath.Join(arsRepo.RepoDir, "resources\\config")
	arsRepo.Config.RepositoryConfigFile =
		NewRepositoryConfigFile(filepath.Join(arsRepo.Config.Dir, "repositoryConfig.json"))
	arsRepo.IndividualizationConfig.Dir = filepath.Join(arsRepo.RepoDir, "resources\\individual_repositories")
	arsRepo.GeneratedOverviewFiles.Dir = filepath.Join(arsRepo.RepoDir, "resources\\overview")
	arsRepo.GeneratedLocalOutput.Dir = filepath.Join(arsRepo.RepoDir, "resources\\test\\output")

	utils.OutputAndAbortIfErrors(
		utils.ValidateAllDirPaths(arsRepo.RepoDir, arsRepo.Config.Dir, arsRepo.IndividualizationConfig.Dir,
			arsRepo.GeneratedOverviewFiles.Dir, arsRepo.GeneratedLocalOutput.Dir))
	log.WithFields(log.Fields{
		"RepoDir":                      arsRepo.RepoDir,
		"ConfigDir":                    arsRepo.Config.Dir,
		"NewRepositoryConfigFile":      arsRepo.Config.RepositoryConfigFile,
		"IndividualizationConfigDir":   arsRepo.IndividualizationConfig.Dir,
		"GeneratedOverviewFilesDir":    arsRepo.GeneratedOverviewFiles.Dir,
		"GeneratedLocalOutputFilesDir": arsRepo.GeneratedLocalOutput.Dir,
	}).Info("Setting global variables:")
	return arsRepo
}
