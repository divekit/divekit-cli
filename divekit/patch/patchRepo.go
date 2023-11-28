package patch

/**
 * This file an "object-oriented lookalike" implementation for the structure of the ARS repository.
 */

import (
	"divekit-cli/divekit"
	"divekit-cli/divekit/ars"
	"divekit-cli/utils/errorHandling"
	"divekit-cli/utils/fileUtils"
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
	patchConfigFileName := filepath.Join(patchRepo.RepoDir, "src/main/config/editorConfig.json")
	patchRepo.PatchConfigFile = NewPatchConfigFile(patchConfigFileName)
	patchRepo.InputDir = filepath.Join(patchRepo.RepoDir, "assets/input")

	errorHandling.OutputAndAbortIfErrors(fileUtils.ValidateAllDirPaths(patchRepo.RepoDir, patchRepo.InputDir),
		"Invalid path have been detected for one or more patchRepo paths")
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
		fmt.Println("Error removing code input directory:", errCode)
		return errCode
	}
	if errTest != nil {
		fmt.Println("Error removing test input directory:", errTest)
		return errTest
	}
	return nil
}

func (patchRepo *PatchRepoType) UpdatePatchConfigFile(repositoryConfigFile *ars.RepositoryConfigFileType) error {
	log.Debug("patch.UpdatePatchConfigFile()")
	patchConfigFile := patchRepo.PatchConfigFile
	err := patchConfigFile.UpdateFromRepositoryConfigFile(repositoryConfigFile)
	if err != nil {
		log.Errorf("Error in patch.UpdatePatchConfigFile():", err)
		return err
	}
	err = patchConfigFile.WriteContent()
	return err
}
