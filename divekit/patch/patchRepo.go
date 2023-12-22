package patch

/**
 * This file an "object-oriented lookalike" implementation for the structure of the ARS repository.
 */

import (
	"divekit-cli/divekit"
	"divekit-cli/divekit/ars"
	"divekit-cli/utils/fileUtils"
	"github.com/apex/log"
	"os"
	"path/filepath"
)

// all the paths used in the ARS repository (all as full paths)
type PatchRepoType struct {
	RepoDir         string
	PatchConfigFile *PatchConfigFileType
	InputDir        string
}

// This method is similar to a constructor in OOP
func NewPatchRepo() (*PatchRepoType, error) {
	log.Debug("patch.NewPatchRepo()")
	var err error
	patchRepo := &PatchRepoType{}
	patchRepo.RepoDir = filepath.Join(divekit.DivekitHomeDir, "divekit-repo-editor")
	patchConfigFileName := filepath.Join(patchRepo.RepoDir, "src/main/config/editorConfig.json")
	patchRepo.PatchConfigFile, err = NewPatchConfigFile(patchConfigFileName)
	if err != nil {
		return nil, err
	}

	patchRepo.InputDir = filepath.Join(patchRepo.RepoDir, "assets/input")
	if err := fileUtils.ValidateAllDirPaths(patchRepo.RepoDir, patchRepo.InputDir); err != nil {
		return nil, err
	}

	log.WithFields(log.Fields{
		"patchRepo.RepoDir":   patchRepo.RepoDir,
		" patchRepo.InputDir": patchRepo.InputDir,
	}).Info("Setting patch repo variables:")

	return patchRepo, nil
}

func (patchRepo *PatchRepoType) CleanInputDir() error {
	codeDirPath := filepath.Join(patchRepo.InputDir, "code")
	testDirPath := filepath.Join(patchRepo.InputDir, "test")

	if errCode := os.RemoveAll(codeDirPath); errCode != nil {
		return &os.PathError{
			Err:  errCode,
			Op:   "Remove code input directory",
			Path: codeDirPath,
		}
	}
	if errTest := os.RemoveAll(testDirPath); errTest != nil {
		return &os.PathError{
			Err:  errTest,
			Op:   "Remove test input directory",
			Path: testDirPath,
		}
	}

	return nil
}

func (patchRepo *PatchRepoType) UpdatePatchConfigFile(repositoryConfigFile *ars.RepositoryConfigFileType) error {
	log.Debug("patch.UpdatePatchConfigFile()")
	patchConfigFile := patchRepo.PatchConfigFile
	if err := patchConfigFile.UpdateFromRepositoryConfigFile(repositoryConfigFile); err != nil {
		log.Errorf("Error in patch.UpdatePatchConfigFile():", err)
		return err
	}
	err := patchConfigFile.WriteContent()

	return err
}
