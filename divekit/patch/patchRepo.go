package patch

/**
 * This file an "object-oriented lookalike" implementation for the structure of the ARS repository.
 */

import (
	"divekit-cli/divekit"
	"divekit-cli/utils"
	"fmt"
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
func NewPatchRepo() *PatchRepoType {
	log.Debug("patch.NewPatchRepo()")
	patchRepo := &PatchRepoType{}
	patchRepo.RepoDir = filepath.Join(divekit.DivekitHomeDir, "divekit-repo-editor")
	patchConfigFileName := filepath.Join(patchRepo.RepoDir, "build\\main\\config\\editorConfig.json")
	patchRepo.PatchConfigFile = NewPatchConfigFile(patchConfigFileName)
	patchRepo.InputDir = filepath.Join(patchRepo.RepoDir, "assets\\input")

	utils.OutputAndAbortIfErrors(utils.ValidateAllDirPaths(patchRepo.RepoDir, patchRepo.InputDir))
	log.WithFields(log.Fields{
		"patchRepo.RepoDir":   patchRepo.RepoDir,
		" patchRepo.InputDir": patchRepo.InputDir,
	}).Info("Setting patch repo variables:")
	return patchRepo
}

func (patchRepo *PatchRepoType) CleanInputDir() error {
	codeDirPath := filepath.Join(patchRepo.InputDir, "code")
	testDirPath := filepath.Join(patchRepo.InputDir, "test")
	errCode := os.RemoveAll(codeDirPath)
	errTest := os.RemoveAll(testDirPath)
	if errCode != nil {
		fmt.Println("Error removing code input directory %s:", errCode)
		return errCode
	}
	if errTest != nil {
		fmt.Println("Error removing test input directory %s:", errTest)
		return errTest
	}
	return nil
}
