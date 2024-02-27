package ars

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

func NewRepositoryConfigFile(path string) (*RepositoryConfigFileType, error) {
	log.Debug("ars.repositoryConfigFile() - path: " + path)
	if err := fileUtils.ValidateAllFilePaths(path); err != nil {
		return nil, err
	}

	return &RepositoryConfigFileType{FilePath: path}, nil
}

func (repositoryConfigFile *RepositoryConfigFileType) ReadContent() error {
	log.Debug("ars.ReadContent() - filePath: " + repositoryConfigFile.FilePath)
	configFile, err := os.ReadFile(repositoryConfigFile.FilePath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}
	err = json.Unmarshal(configFile, &repositoryConfigFile.Content)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}
	repositoryConfigFile.CheckForDeathTraps()

	return nil
}

func (repositoryConfigFile *RepositoryConfigFileType) WriteContent() error {
	log.Debug("ars.WriteContent() - filePath: " + repositoryConfigFile.FilePath)
	updatedConfig, err := json.MarshalIndent(repositoryConfigFile.Content, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	err = ioutil.WriteFile(repositoryConfigFile.FilePath, updatedConfig, 0644)
	if err != nil {
		return fmt.Errorf("failed to write updated config file: %w", err)
	}

	return nil
}

func (repositoryConfigFile *RepositoryConfigFileType) CheckForDeathTraps() {
	log.Debug("ars.checkForDeathTraps() - filePath: " + repositoryConfigFile.FilePath)
	if !repositoryConfigFile.Content.General.LocalMode && repositoryConfigFile.Content.Remote.DeleteExistingRepositories {
		if !errorHandling.Confirm(
			"Your repositoryConfig.json sets local mode to false, and sets \"deleteExistingRepositories\" \n" +
				"to true. This means that you'll delete all repositories in the target group. \n" +
				"Are you sure you want to do this?") {
			log.Fatalf("Aborting action")
		}
	}
}

func (repositoryConfigFile *RepositoryConfigFileType) Clone() (*RepositoryConfigFileType, error) {
	log.Debug("ars.Clone() - filePath: " + repositoryConfigFile.FilePath)
	return repositoryConfigFile.CloneToDifferentLocation(repositoryConfigFile.FilePath)
}

func (repositoryConfigFile *RepositoryConfigFileType) CloneToDifferentLocation(newFilePath string) (*RepositoryConfigFileType, error) {
	log.Debug("ars.CloneToDifferentLocation() - newFilePath: " + newFilePath)
	var newFile *RepositoryConfigFileType
	var err error

	newFile, err = NewRepositoryConfigFile(newFilePath)
	if err != nil {
		return nil, err
	}

	if err = fileUtils.DeepCopy(repositoryConfigFile, newFile); err != nil {
		return nil, err
	}

	newFile.FilePath = newFilePath

	return newFile, nil
}
