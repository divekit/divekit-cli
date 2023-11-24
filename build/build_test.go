package installer

import (
	"github.com/stretchr/testify/assert"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestBuild(t *testing.T) {
	testCases := []struct {
		name       string
		scriptPath string // input
		filePath   string // expected
		error      error  // expected
	}{
		{"Should build a msi file successfully", ".\\build.bat", ".\\output\\*.msi", nil},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			err := executeScript(testCase.scriptPath)

			assert.True(t, fileExistsWithExtension(testCase.filePath), "Failed to build the MSI File.")
			assert.IsType(t, testCase.error, err)
		})
	}
}

func fileExistsWithExtension(filePath string) bool {
	files, err := filepath.Glob(filePath)
	if err != nil {
		return false
	}

	return len(files) > 0
}
func executeScript(scriptPath string) error {
	cmd := exec.Command(scriptPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
