package ars

import (
	"divekit-cli/divekit"
	"divekit-cli/utils/fileUtils"
	"fmt"
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

func NewARSRepo() (*ARSRepoType, error) {
	log.Debug("ars.NewARSRepo()")
	var err error

	arsRepo := &ARSRepoType{}
	arsRepo.RepoDir = filepath.Join(divekit.DivekitHomeDir, "divekit-automated-repo-setup")
	arsRepo.Config.Dir = filepath.Join(arsRepo.RepoDir, "resources/config")
	arsRepo.Config.RepositoryConfigFile, err = NewRepositoryConfigFile(filepath.Join(arsRepo.Config.Dir, "repositoryConfig.json"))
	if err != nil {
		return nil, err
	}

	arsRepo.IndividualizationConfig.Dir = filepath.Join(arsRepo.RepoDir, "resources/individual_repositories")
	arsRepo.GeneratedOverviewFiles.Dir = filepath.Join(arsRepo.RepoDir, "resources/overview")
	arsRepo.GeneratedLocalOutput.Dir = filepath.Join(arsRepo.RepoDir, "resources/test/output")

	if err := fileUtils.ValidateAllDirPaths(arsRepo.RepoDir, arsRepo.Config.Dir, arsRepo.IndividualizationConfig.Dir,
		arsRepo.GeneratedOverviewFiles.Dir, arsRepo.GeneratedLocalOutput.Dir); err != nil {
		return nil, fmt.Errorf("the path to the ARS repo is invalid: %w", err)
	}

	log.WithFields(log.Fields{
		"RepoDir":                      arsRepo.RepoDir,
		"ConfigDir":                    arsRepo.Config.Dir,
		"NewRepositoryConfigFile":      arsRepo.Config.RepositoryConfigFile,
		"IndividualizationConfigDir":   arsRepo.IndividualizationConfig.Dir,
		"GeneratedOverviewFilesDir":    arsRepo.GeneratedOverviewFiles.Dir,
		"GeneratedLocalOutputFilesDir": arsRepo.GeneratedLocalOutput.Dir,
	}).Info("Setting global variables")

	return arsRepo, nil
}
