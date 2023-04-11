package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type ConfigRepositoryType struct {
	General struct {
		LocalMode                     bool `json:"localMode"`
		CreateTestRepository          bool `json:"createTestRepository"`
		VariateRepositories           bool `json:"variateRepositories"`
		DeleteSolution                bool `json:"deleteSolution"`
		ActivateVariableValueWarnings bool `json:"activateVariableValueWarnings"`
		MaxConcurrentWorkers          int  `json:"maxConcurrentWorkers"`
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

var ConfigRepository ConfigRepositoryType

func ReadConfigRepository(configFilePath string) error {
	configFile, err := os.ReadFile(configFilePath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %v", err)
	}
	err = json.Unmarshal(configFile, &ConfigRepository)
	if err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %v", err)
	}
	return nil
}

func WriteConfigRepository(configFilePath string) error {
	updatedConfig, err := json.MarshalIndent(ConfigRepository, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %v", err)
	}

	err = ioutil.WriteFile(configFilePath, updatedConfig, 0644)
	if err != nil {
		return fmt.Errorf("failed to write updated config file: %v", err)
	}

	return nil
}
