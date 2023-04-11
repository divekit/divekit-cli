package config

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	// ARS
	ARS_REPO_NAME              = "divekit-automated-repo-setup"
	REPOSITORY_CONFIG_FILENAME = "repositoryConfig.json"

	// Origin repo
	DIVEKIT_DIR_NAME       = ".divekit"
	DISTRIBUTIONS_DIR_NAME = "distributions"
)

var (
	DivekitHome string
)

func GetDivekitHome() (string, error) {
	divekitHome := os.Getenv("DIVEKIT_HOME")
	return divekitHome, nil
}

func GetFullPathDivekitDir() (string, error) {
	divekitHome, err := GetDivekitHome()
	if err != nil {
		return "", err
	}
	path := filepath.Join(divekitHome, DIVEKIT_DIR_NAME)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", fmt.Errorf("Divekit directory does not exist: %s", path)
	}
	return path, nil
}

func GetFullPathDistributionsDir() (string, error) {
	divekitDir, err := GetFullPathDivekitDir()
	if err != nil {
		return "", err
	}
	path := filepath.Join(divekitDir, DISTRIBUTIONS_DIR_NAME)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", fmt.Errorf("Distributions directory does not exist: %s", path)
	}
	return path, nil
}

func GetFullPathRepositoryConfigFile() (string, error) {
	distributionsDir, err := GetFullPathDistributionsDir()
	if err != nil {
		return "", err
	}
	path := filepath.Join(distributionsDir, REPOSITORY_CONFIG_FILENAME)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return "", fmt.Errorf("Repository config file does not exist: %s", path)
	}
	return path, nil
}
