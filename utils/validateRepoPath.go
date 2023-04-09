package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func ValidateRepoPath(repoPath string, isOrigin bool) error {
	_, err := os.Stat(repoPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("Repo does not exist: %s", repoPath)
	}

	if isOrigin {
		divekitDir := filepath.Join(repoPath, ".divekit")
		_, err = os.Stat(divekitDir)
		if os.IsNotExist(err) {
			return fmt.Errorf(".divekit subfolder not found in: %s", repoPath)
		}
	}
	return nil
}
