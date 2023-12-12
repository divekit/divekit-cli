package runner

import (
	"divekit-cli/utils/fileUtils"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"os/exec"
	"testing"
)

var jsonPath string // This path leads to a temporary directory with a simulated package.json file for mocking.

func TestRunNPMStart(t *testing.T) {
	testCases := []struct {
		name       string
		path       string // input: path leads to a simulated package.json file and can be empty to raise an error.
		dryRunFlag bool   // input
		executed   bool   // expected
		error      error  // expected
	}{
		{"True dryRunflag should skip execution", jsonPath, true, false, nil},
		{"False dryRunflag should execute", jsonPath, false, true, nil},
		{"True dryRunflag with no path should skip execution", "", true, false, nil},
		{"False dryRunflag with no path should execute and fail", "", false, true,
			&exec.ExitError{}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			DryRunFlag = testCase.dryRunFlag

			executed, err := RunNPMStart(testCase.path, "")
			assert.Equal(t, testCase.executed, executed)
			assert.IsType(t, testCase.error, err)
		})
	}
}

func TestRunNPMStartAlways(t *testing.T) {
	testCases := []struct {
		name  string
		path  string // input
		error error  // expected
	}{
		{"Valid path should be successful", jsonPath, nil},
		{"No path should fail", "", &exec.ExitError{}},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			err := RunNPMStartAlways(testCase.path, "")
			assert.IsType(t, testCase.error, err)
		})
	}
}

func TestMain(m *testing.M) {
	setup()
	e := m.Run() // run the tests
	teardown()
	os.Exit(e) // report the exit code
}

func setup() {
	log.Println("runner package: Setting up.")
	createJson()
}
func teardown() {
	log.Println("runner package: Tearing down.")
	deleteJson()
}

// createJson generates a simulated package.json in a temporary directory and provides its path as a
// return value.
func createJson() {
	jsonPath = fileUtils.CreateTmpDir()
	fileUtils.CreateFile(jsonPath, "package.json", "{ \"scripts\": { \"start\": \"\" } }")
}

// deleteJson removes the package.json file along with the temporary directory generated by createJson.
func deleteJson() {
	fileUtils.DeleteDir(jsonPath)
}
