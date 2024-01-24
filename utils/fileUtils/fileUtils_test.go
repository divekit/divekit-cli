package fileUtils

import (
	"divekit-cli/utils/errorHandling"
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"reflect"
	"sort"
	"syscall"
	"testing"
)

// testDirs contains various paths leading to test directories intended for mock purposes. These paths should never be
// modified during tests, because it would lead to wrong test results.
var testDirs TestDirs

func TestFindFilesInDirRecursively(t *testing.T) {
	testCases := []struct {
		name         string
		testDir      string   // input
		searchedFile string   // input
		foundFiles   []string // expected
		error        error    // expected
	}{
		{
			"Fail to find a file in a non-existent dir",
			"non-existent-dir",
			"1.txt",
			[]string{},
			syscall.ERROR_FILE_NOT_FOUND,
		},
		{
			"Find nothing with an invalid searched file",
			testDirs.Nested,
			"invalid",
			[]string{},
			nil,
		},
		{
			"Find nothing with an existing dir name as searched file",
			testDirs.Nested,
			"a",
			[]string{},
			nil,
		},
		{
			"Find 1's in a nested dir recursively",
			testDirs.Nested,
			"1.txt",
			[]string{"a/a/1.txt", "b/1.txt", "b/a/1.txt"},
			nil,
		},
		{
			"Find 0's in a nested dir recursively",
			testDirs.Nested,
			"0.txt",
			[]string{"a/a/0.txt", "b/0.txt"},
			nil,
		},
		{
			"Distinguish between file and dir",
			testDirs.EqualFileAndDir,
			"a",
			[]string{"a/a"},
			nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			foundFiles, err := FindFilesInDirRecursively(testCase.testDir, testCase.searchedFile)
			foundFiles = ToRelPaths(foundFiles, testCase.testDir)

			errorHandling.IsErrorType(t, testCase.error, err)
			assert.ElementsMatch(t, testCase.foundFiles, foundFiles)
		})
	}
}

func TestFindAnyFilesInDirRecursively(t *testing.T) {
	testCases := []struct {
		name       string
		testDir    string   // input
		foundFiles []string // expected
		error      error    // expected
	}{
		{
			"Fail to find files in a non-existent dir",
			"non-existent-dir",
			[]string{},
			syscall.ERROR_FILE_NOT_FOUND,
		},
		{
			"Find nothing within an empty dir",
			testDirs.EmptyDir,
			[]string{},
			nil,
		},
		{
			"Find everything in a nested dir recursively",
			testDirs.Nested,
			[]string{"a/a/0.txt", "a/a/1.txt", "a/a/11.txt", "a/a/new_file.json", "b/0.txt", "b/1.txt", "b/a/1.txt"},
			nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			foundFiles, err := FindAnyFilesInDirRecursively(testCase.testDir)
			foundFiles = ToRelPaths(foundFiles, testCase.testDir)

			errorHandling.IsErrorType(t, testCase.error, err)
			assert.ElementsMatch(t, testCase.foundFiles, foundFiles)
		})
	}
}

func TestFindFilesInDir(t *testing.T) {
	testCases := []struct {
		name         string
		testDir      string   // input
		searchedFile string   // input
		foundFiles   []string // expected
		error        error    // expected
	}{
		{
			"Fail to find a file in a non-existent dir",
			"non-existent-dir",
			"1.txt",
			[]string{},
			syscall.ERROR_FILE_NOT_FOUND,
		},
		{
			"Find nothing with an invalid searched file",
			testDirs.Simple,
			"invalid",
			[]string{},
			nil,
		},
		{
			"Find nothing with an existing dir name as searched file",
			testDirs.Simple,
			"a",
			[]string{},
			nil,
		},
		{
			"Find nothing within a nested dir",
			testDirs.Nested,
			"1.txt",
			[]string{},
			nil,
		},
		{
			"Find 1's in a dir",
			testDirs.Simple,
			"1.txt",
			[]string{"1.txt"},
			nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			foundFiles, err := FindFilesInDir(testCase.testDir, testCase.searchedFile)
			foundFiles = ToRelPaths(foundFiles, testCase.testDir)

			errorHandling.IsErrorType(t, testCase.error, err)
			assert.ElementsMatch(t, testCase.foundFiles, foundFiles)
		})
	}
}

func TestTransformIntoRelativePaths(t *testing.T) {
	testCases := []struct {
		name    string
		root    string // input
		absPath string // input
		relPath string // expected
		error   error  // expected

	}{
		{"/a/b/c should be b/c", "/a", "/a/b/c", "b/c", nil},
		{"/b/c should be ../b/c", "/a", "/b/c", "../b/c", nil},
		{"./a/b/c should raise an error", "/a", "./a/b/c", "", fmt.Errorf("")},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			relPath, err := TransformIntoRelativePaths(testCase.root, testCase.absPath)
			relPath = UnifyPath(relPath)

			errorHandling.IsErrorType(t, testCase.error, err)
			assert.Equal(t, testCase.relPath, relPath)
		})
	}
}

func TestCopyFile(t *testing.T) {
	testCases := []struct {
		name        string
		srcFilePath string        // input
		dstDirPath  func() string // input
		srcFile     string        // input
		error       error         // expected
	}{
		{
			"Throw an error for a directory",
			testDirs.Simple,
			CreateTmpDir,
			"/a",
			fmt.Errorf(""),
		},
		{
			"Throw an error for a non-existing file",
			testDirs.Simple,
			CreateTmpDir,
			"/non_existing_file",
			syscall.ERROR_FILE_NOT_FOUND,
		},
		{
			"Should copy a file into another directory",
			testDirs.OneFileDir,
			CreateTmpDir,
			"/1.txt",
			nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// Create a destination directory before each test
			dstDirPath := testCase.dstDirPath()
			// Delete a destination directory after each test
			defer DeleteDir(dstDirPath)

			filePath := testCase.srcFilePath + testCase.srcFile
			err := CopyFile(filePath, dstDirPath)

			errorHandling.IsErrorType(t, testCase.error, err)
			if err == nil {
				compareDirContent(t, testCase.srcFilePath, dstDirPath)
			}
		})
	}
}

func TestCopyAllFilesInDir(t *testing.T) {
	testCases := []struct {
		name    string
		srcDir  string        // input
		dstDir  func() string // input
		srcFile string        // input
		dstFile string        // input
		error   error         // expected
	}{
		{
			"Throw an error for a non-existing file",
			testDirs.Nested,
			CreateTmpDir,
			"non_existing_file",
			"",
			syscall.ERROR_FILE_NOT_FOUND,
		},
		{
			"Copy one file with source file should fail",
			testDirs.Nested,
			CreateTmpDir,
			"/a/a/0.txt",
			"",
			syscall.EISDIR,
		},
		{
			"Copy one file with source file and destination file should fail",
			testDirs.Nested,
			CreateTmpDir,
			"/a/a/0.txt",
			"/a/a/0.txt",
			syscall.ERROR_FILE_NOT_FOUND,
		},
		{
			"Copy all files into another directory",
			testDirs.Nested,
			CreateTmpDir,
			"",
			"",
			nil,
		},
		{
			"Copy files into another directory with a subfolder as target",
			testDirs.Nested,
			CreateTmpDir,
			"",
			"/a/a",
			nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			dstDir := testCase.dstDir()
			defer DeleteDir(dstDir)

			srcPath := testCase.srcDir + testCase.srcFile
			dstPath := dstDir + testCase.dstFile
			err := CopyAllFilesInDir(srcPath, dstPath)

			errorHandling.IsErrorType(t, testCase.error, err)
			if err == nil {
				compareDirContent(t, srcPath, dstPath)
			}
		})
	}
}

func TestValidateAllPaths(t *testing.T) {
	testCases := []struct {
		name        string
		paths       []string // input
		shouldBeDir bool     // input
		error       error    // expected
	}{
		{
			"Throw an error for a non-existing file",
			[]string{"non_existing_file"},
			false,
			&InvalidPathError{},
		},
		{
			"Throw an error for a non-existing file with an existing file",
			[]string{testDirs.OneFileDir + "/1.txt", "non_existing_file"},
			false,
			&InvalidPathError{},
		},
		{
			"Throw an error for an existing file handled as a directory",
			[]string{testDirs.OneFileDir + "/1.txt"},
			true,
			&InvalidPathError{},
		},
		{
			"Throw an error for an existing directory handled as a file",
			[]string{testDirs.Simple + "/a"},
			false,
			&InvalidPathError{},
		},
		{
			"Check one valid path",
			[]string{testDirs.Simple + "/a"},
			true,
			nil,
		},
		{
			"Check multiple valid paths",
			[]string{testDirs.Simple + "/a", testDirs.Nested + "/a", testDirs.Nested + "/b/a"},
			true,
			nil,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			err := ValidateAllPaths(testCase.shouldBeDir, testCase.paths...)
			errorHandling.IsErrorType(t, testCase.error, err)
		})
	}
}

func TestFindUniqueFileWithPrefix(t *testing.T) {
	testCases := []struct {
		name      string
		dir       string // input
		prefix    string // input
		foundFile string // expected
		error     error  // expected
	}{
		{
			"Throw an error for a non-existing dir",
			"non_existing_dir",
			"-",
			"",
			syscall.ERROR_FILE_NOT_FOUND,
		},
		{
			"Throw an error for an empty prefix",
			testDirs.Nested,
			"",
			"",
			fmt.Errorf(""),
		},
		{
			"Do not find a file recursively",
			testDirs.Nested,
			"new",
			"",
			fmt.Errorf(""),
		},
		{
			"Multiple files with the same prefix should throw an error",
			testDirs.Nested + "/a/a",
			"1",
			"",
			fmt.Errorf(""),
		},
		{
			"Find a file",
			testDirs.Nested + "/a/a",
			"n",
			"new_file.json",
			nil,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			path, err := FindUniqueFileWithPrefix(testCase.dir, testCase.prefix)

			errorHandling.IsErrorType(t, testCase.error, err)
			foundFile := GetBaseName(path)
			if foundFile == "." {
				foundFile = ""
			}
			assert.Equal(t, testCase.foundFile, foundFile)
		})
	}
}

func TestListSubFolderNames(t *testing.T) {
	testCases := []struct {
		name       string
		folderPath string   // input
		subFolders []string // expected
		error      error    // expected
	}{
		{
			"Throw an error for a non-existing dir",
			"non_existing_dir",
			[]string{},
			syscall.ERROR_FILE_NOT_FOUND,
		},
		{
			"Find subfolders in a simple dir",
			testDirs.Simple,
			[]string{"a"},
			nil,
		},
		{
			"Find subfolders in a nested dir",
			testDirs.Nested,
			[]string{"a", "b"},
			nil,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			subFolders, err := ListSubFolderNames(testCase.folderPath)

			errorHandling.IsErrorType(t, testCase.error, err)
			assert.ElementsMatch(t, testCase.subFolders, subFolders)
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
	log.Println("fileUtils package: Setting up.")
	createTestDirs()
}

func teardown() {
	log.Println("fileUtils package: Tearing down.")
	deleteTestDirs()
}

type TestDirs struct {
	Simple          string
	Nested          string
	EqualFileAndDir string
	EmptyDir        string
	OneFileDir      string
}

func createTestDirs() {
	testDirs = TestDirs{
		Simple:          createSimpleDir(),
		Nested:          createNestedDir(),
		EqualFileAndDir: createEqualFileAndDir(),
		EmptyDir:        createEmptyDir(),
		OneFileDir:      createDirWithOneFile(),
	}
}

func deleteTestDirs() {
	values := reflect.ValueOf(testDirs)

	for i := 0; i < values.NumField(); i++ {
		testDir := values.Field(i).String()
		DeleteDir(testDir)
	}
}

func createSimpleDir() string {
	// ./
	// ├─ a/
	// │  └─ 1.txt
	// └─ 1.txt

	rootDir := CreateTmpDir()

	CreateDir(rootDir + "/a")
	for _, file := range []string{"/1.txt", "/a/1.txt"} {
		CreateFile(rootDir, file, file)
	}

	return rootDir
}

func createNestedDir() string {
	// ./
	// ├─ a/
	// │  └─ a/
	// │     ├─ 0.txt
	// │     ├─ 1.txt
	// │     ├─ 11.txt
	// │     └─ new_file.json
	// └─ b/
	//    ├─ a/
	//    │  └─ 1.txt
	//    ├─ 0.txt
	//    └─ 1.txt

	rootDir := CreateTmpDir()

	for _, subDir := range []string{"/a/a", "/b/a"} {
		CreateDir(rootDir + subDir)
	}
	for _, file := range []string{"/a/a/0.txt", "/a/a/1.txt", "/a/a/11.txt", "/a/a/new_file.json", "/b/0.txt", "/b/1.txt", "/b/a/1.txt"} {
		CreateFile(rootDir, file, file)
	}
	return rootDir
}
func createEqualFileAndDir() string {
	// ./
	// └─ a/ <- directory
	//    └─ a <- file

	rootDir := CreateTmpDir()

	CreateDir(rootDir + "/a")
	CreateFile(rootDir, "/a/a", "/a/a")

	return rootDir
}

func createEmptyDir() string {
	return CreateTmpDir()
}

func createDirWithOneFile() string {
	// ./
	// └─ 1.txt

	rootDir := CreateTmpDir()

	CreateFile(rootDir, "/1.txt", "/1.txt")

	return rootDir
}

// compareDirContent compares relative file paths of two directories and
// makes sure that the content of these files are identical.
func compareDirContent(t *testing.T, expectedDir string, actualDir string) {
	expectedPaths, err := FindAnyFilesInDirRecursively(expectedDir)
	assert.NoError(t, err, "Failed to find files in the expected directory")
	actualPaths, err := FindAnyFilesInDirRecursively(actualDir)
	assert.NoError(t, err, "Failed to find files in the actual directory")

	expectedFiles := ToRelPaths(expectedPaths, expectedDir)
	actualFiles := ToRelPaths(actualPaths, actualDir)

	ok := assert.ElementsMatch(t, expectedFiles, actualFiles)
	if ok {
		// Ensure that the order of files is the same, to compare the content of files with the same name
		sort.Strings(expectedFiles)
		sort.Strings(actualFiles)

		for i := range expectedFiles {
			expectedPath := expectedDir + "/" + expectedFiles[i]
			actualPath := actualDir + "/" + actualFiles[i]
			compareFileContent(t, expectedPath, actualPath)
		}
	}
}

func compareFileContent(t *testing.T, expectedFile string, actualFile string) {
	expectedContent, err := os.ReadFile(expectedFile)
	assert.NoError(t, err, "Failed to read the expected file '%v'", expectedFile)
	actualContent, err := os.ReadFile(actualFile)
	assert.NoError(t, err, "Failed to read the actual file '%v'", actualFile)

	assert.Equal(t, string(expectedContent), string(actualContent),
		"The content of the expected file '%v' does not match with the content of the actual file '%v'",
		expectedFile, actualFile)
}
