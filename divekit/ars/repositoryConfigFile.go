package ars

/**
 * This file an "object-oriented lookalike" implementation for the repositoryConfig.json file.
 * It is used to read and write the repositoryConfig.json file.
 */

import (
	"divekit-cli/utils/errorHandling"
	"divekit-cli/utils/fileUtils"
	"encoding/json"
	"fmt"
	"github.com/apex/log"
	"io/ioutil"
	"os"
)

// struct for the repositoryConfig.json file
type RepositoryConfigFileType struct {
	FilePath string
	Content  struct {
		General struct {
			LocalMode                     bool   `json:"localMode"`
			CreateTestRepository          bool   `json:"createTestRepository"`
			VariateRepositories           bool   `json:"variateRepositories"`
			DeleteSolution                bool   `json:"deleteSolution"`
			ActivateVariableValueWarnings bool   `json:"activateVariableValueWarnings"`
			MaxConcurrentWorkers          int    `json:"maxConcurrentWorkers"`
			GlobalLogLevel                string `json:"globalLogLevel"`
		} `json:"general"`
		Repository struct {
			RepositoryName    string     `json:"repositoryName"`
			RepositoryCount   int        `json:"repositoryCount"`
			RepositoryMembers [][]string `json:"repositoryMembers"`
		} `json:"repository"`
		IndividualRepositoryPersist struct {
			UseSavedIndividualRepositories      bool   `json:"useSavedIndividualRepositories"`
			SavedIndividualRepositoriesFileName string `json:"savedIndividualRepositoriesFileName"`
		} `json:"individualRepositoryPersist"`
		Local struct {
			OriginRepositoryFilePath string   `json:"originRepositoryFilePath"`
			SubsetPaths              []string `json:"subsetPaths"`
		} `json:"local"`
		Remote struct {
			OriginRepositoryId          int  `json:"originRepositoryId"`
			CodeRepositoryTargetGroupId int  `json:"codeRepositoryTargetGroupId"`
			TestRepositoryTargetGroupId int  `json:"testRepositoryTargetGroupId"`
			DeleteExistingRepositories  bool `json:"deleteExistingRepositories"`
			AddUsersAsGuests            bool `json:"addUsersAsGuests"`
		} `json:"remote"`
		Overview struct {
			GenerateOverview     bool   `json:"generateOverview"`
			OverviewRepositoryId int    `json:"overviewRepositoryId"`
			OverviewFileName     string `json:"overviewFileName"`
		} `json:"overview"`
	}
}

// This method is similar to a constructor in OOP
func NewRepositoryConfigFile(path string) *RepositoryConfigFileType {
	log.Debug("ars.repositoryConfigFile() - path: " + path)
	errorHandling.OutputAndAbortIfErrors(fileUtils.ValidateAllFilePaths(path),
		"The path to the repository config file is invalid")
	return &RepositoryConfigFileType{
		FilePath: path,
	}
}

func (repositoryConfigFile *RepositoryConfigFileType) ReadContent() error {
	log.Debug("ars.ReadContent() - filePath: " + repositoryConfigFile.FilePath)
	configFile, err := os.ReadFile(repositoryConfigFile.FilePath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %v", err)
	}
	err = json.Unmarshal(configFile, &repositoryConfigFile.Content)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %v", err)
	}
	repositoryConfigFile.CheckForDeathTraps()
	return nil
}

func (repositoryConfigFile *RepositoryConfigFileType) WriteContent() error {
	log.Debug("ars.WriteContent() - filePath: " + repositoryConfigFile.FilePath)
	updatedConfig, err := json.MarshalIndent(repositoryConfigFile.Content, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}

	err = ioutil.WriteFile(repositoryConfigFile.FilePath, updatedConfig, 0644)
	if err != nil {
		return fmt.Errorf("failed to write updated config file: %v", err)
	}

	return nil
}

func (repositoryConfigFile *RepositoryConfigFileType) CheckForDeathTraps() bool {
	log.Debug("ars.checkForDeathTraps() - filePath: " + repositoryConfigFile.FilePath)
	if !repositoryConfigFile.Content.General.LocalMode && repositoryConfigFile.Content.Remote.DeleteExistingRepositories {
		errorHandling.Confirm(
			"Your repositoryConfig.json sets local mode to false, and sets \"deleteExistingRepositories\" \n" +
				"to true. This means that you'll delete all repositories in the target group. \n" +
				"Are you sure you want to do this?")
	}
	return true
}

func (repositoryConfigFile *RepositoryConfigFileType) Clone() *RepositoryConfigFileType {
	log.Debug("ars.Clone() - filePath: " + repositoryConfigFile.FilePath)
	return repositoryConfigFile.CloneToDifferentLocation(repositoryConfigFile.FilePath)
}

func (repositoryConfigFile *RepositoryConfigFileType) CloneToDifferentLocation(newFilePath string) *RepositoryConfigFileType {
	log.Debug("ars.CloneToDifferentLocation() - newFilePath: " + newFilePath)
	newFile := NewRepositoryConfigFile(newFilePath)
	errorHandling.OutputAndAbortIfError(fileUtils.DeepCopy(repositoryConfigFile, newFile),
		"Could not copy the repository config file")
	newFile.FilePath = newFilePath
	return newFile
}
