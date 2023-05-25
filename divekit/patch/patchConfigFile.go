package patch

/**
 * This file an "object-oriented lookalike" implementation for the repositoryConfig.json file.
 * It is used to read and write the repositoryConfig.json file.
 */

import (
	"divekit-cli/divekit/ars"
	"divekit-cli/utils"
	"encoding/json"
	"fmt"
	"github.com/apex/log"
	"io/ioutil"
	"os"
	"time"
)

// struct for the editorConfig.json file
type PatchConfigFileType struct {
	FilePath string
	Content  struct {
		OnlyUpdateTestProjects bool   `json:"onlyUpdateTestProjects"`
		OnlyUpdateCodeProjects bool   `json:"onlyUpdateCodeProjects"`
		GroupIds               []int  `json:"groupIds"`
		LogLevel               string `json:"logLevel"`
		CommitMsg              string `json:"commitMsg"`
	}
}

// This method is similar to a constructor in OOP
func NewPatchConfigFile(path string) *PatchConfigFileType {
	log.Debug("patch.patchConfigFile() - path: " + path)
	utils.OutputAndAbortIfErrors(utils.ValidateAllFilePaths(path))
	log.WithFields(log.Fields{
		"PatchConfigFileType.FilePath": path,
	}).Info("Setting NewPatchConfigFile variables:")
	return &PatchConfigFileType{
		FilePath: path,
	}
}

func (patchConfigFile *PatchConfigFileType) UpdateFromRepositoryConfigFile(repositoryConfigFile *ars.RepositoryConfigFileType) error {
	log.Debug("patch.UpdateFromRepositoryConfigFile() - repositoryConfigFile: " + repositoryConfigFile.FilePath)
	patchConfigFile.Content.OnlyUpdateTestProjects = false
	patchConfigFile.Content.OnlyUpdateCodeProjects = false
	patchConfigFile.Content.GroupIds = make([]int, 2)
	patchConfigFile.Content.GroupIds[0] = repositoryConfigFile.Content.Remote.CodeRepositoryTargetGroupId
	patchConfigFile.Content.GroupIds[1] = repositoryConfigFile.Content.Remote.TestRepositoryTargetGroupId
	patchConfigFile.Content.LogLevel = utils.LogLevelAsString()
	currentTime := time.Now()
	formattedTime := currentTime.Format("2006-01-02 15:04")
	patchConfigFile.Content.CommitMsg = "Patch applied on " + formattedTime
	err := patchConfigFile.WriteContent()
	return err
}

func (patchConfigFile *PatchConfigFileType) ReadContent() error {
	log.Debug("patch.ReadContent() - filePath: " + patchConfigFile.FilePath)
	configFile, err := os.ReadFile(patchConfigFile.FilePath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %v", err)
	}
	err = json.Unmarshal(configFile, &patchConfigFile.Content)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %v", err)
	}
	return nil
}

func (patchConfigFile *PatchConfigFileType) WriteContent() error {
	log.Debug("patch.WriteContent() - filePath: " + patchConfigFile.FilePath)
	updatedConfig, err := json.MarshalIndent(patchConfigFile.Content, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}
	err = ioutil.WriteFile(patchConfigFile.FilePath, updatedConfig, 0644)
	if err != nil {
		return fmt.Errorf("failed to write updated config file: %v", err)
	}
	return nil
}
