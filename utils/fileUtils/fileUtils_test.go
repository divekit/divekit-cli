package fileUtils

import (
	"divekit-cli/utils/testUtils"
	"errors"
	"github.com/stretchr/testify/assert"
	"io/fs"
	"log"
	"os"
	"reflect"
	"testing"
)

// testDirs contains various paths leading to config directories intended for mock purposes. These paths should never be
// modified during tests, because it would lead to wrong config results.
var TestDirs testDirs

func TestFindFilesInDirRecursively(t *testing.T) {
	testCases := []struct {
		name         string
		testDir      string   // input
		searchedFile string   // input
		foundFiles   []string // expected
		error        error    // expected
	}{
		{
			"Find 1's in a nested dir recursively",
			TestDirs.Nested,
			"1.txt",
			[]string{"a/a/1.txt", "b/1.txt", "b/a/1.txt"},
			nil,
		},
		{
			"Find 0's in a nested dir recursively",
			TestDirs.Nested,
			"0.txt",
			[]string{"a/a/0.txt", "b/0.txt"},
			nil,
		},
		{
			"Distinguish between file and dir",
			TestDirs.EqualFileAndDir,
			"a",
			[]string{"a/a"},
			nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			foundFiles, err := FindFilesInDirRecursively(testCase.searchedFile, testCase.testDir)
			foundFiles = testUtils.ToRelPaths(foundFiles, testCase.testDir)

			assert.IsType(t, testCase.error, err)
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
			"Find 1's in a dir",
			TestDirs.Simple,
			"1.txt",
			[]string{"1.txt"},
			nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			foundFiles, err := FindFilesInDir(testCase.searchedFile, testCase.testDir)
			foundFiles = testUtils.ToRelPaths(foundFiles, testCase.testDir)

			assert.IsType(t, testCase.error, err)
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
		{"./a/b/c should raise an error", "/a", "./a/b/c", "", errors.New("")},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			relPath, err := TransformIntoRelativePaths(testCase.root, testCase.absPath)
			relPath = testUtils.UnifyPath(relPath)

			assert.IsType(t, testCase.error, err)
			assert.Equal(t, testCase.relPath, relPath)
		})
	}
}

func TestCopyFile(t *testing.T) {
	testCases := []struct {
		name        string
		srcFilePath string // input
		dstDirPath  string // input
		shouldCopy  bool   // expected
		error       error  // expected
	}{
		{
			"Should throw an error and not copy a directory",
			TestDirs.Simple,
			testUtils.CreateTmpDir(),
			false,
			&fs.PathError{},
		},
		{
			"Should copy a file into another directory",
			TestDirs.Simple + "/1.txt",
			testUtils.CreateTmpDir(),
			true,
			nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			defer testUtils.DeleteDir(testCase.dstDirPath)
			err := CopyFile(testCase.srcFilePath, testCase.dstDirPath)

			dstFilePath := testCase.dstDirPath + "/" + testUtils.GetBaseName(testCase.srcFilePath)
			if testCase.shouldCopy {
				assert.FileExists(t, dstFilePath)
			} else {
				assert.NoFileExists(t, dstFilePath)
			}
			assert.IsType(t, testCase.error, err)
		})
	}

}

func TestCopyAllFilesInDir(t *testing.T) {
	//CopyAllFilesInDir()
	// todo
}

func TestValidateFilePath(t *testing.T) {
	//ValidateFilePath()
	// todo
}

func TestValidateAllPaths(t *testing.T) {
	//ValidateAllPaths()
	// todo
}

func TestValidateAllDirPaths(t *testing.T) {
	//ValidateAllDirPaths()
	// todo
}

func TestFindUniqueFileWithPrefix(t *testing.T) {
	//FindUniqueFileWithPrefix()
	// todo
}

func TestListSubfolderNames(t *testing.T) {
	// todo
}

func TestValidateAllFilePaths(t *testing.T) {
	//ValidateAllFilePaths()
	// todo
}

func TestDeepCopy(t *testing.T) {
	//DeepCopy()
	// todo
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

type testDirs struct {
	Simple          string
	Nested          string
	EqualFileAndDir string
}

func createTestDirs() {
	TestDirs = testDirs{
		Simple:          createSimpleDir(),
		Nested:          createNestedDir(),
		EqualFileAndDir: createEqualFileAndDir(),
	}
}

func deleteTestDirs() {
	values := reflect.ValueOf(TestDirs)

	for i := 0; i < values.NumField(); i++ {
		testDir := values.Field(i).String()
		testUtils.DeleteDir(testDir)
	}
}

func createSimpleDir() string {
	// ./
	// ├─ a/
	// │  └─ 1.txt
	// └─ 1.txt

	rootDir := testUtils.CreateTmpDir()

	testUtils.CreateDir(rootDir + "/a")
	for _, file := range []string{"/1.txt", "/a/1.txt"} {
		testUtils.CreateFile(rootDir, file, "")
	}

	return rootDir
}
func createNestedDir() string {
	// ./
	// ├─ a/
	// │  └─ a/
	// │     ├─ 0.txt
	// │     └─ 1.txt
	// └─ b/
	//    ├─ a/
	//    │  └─ 1.txt
	//    ├─ 0.txt
	//    └─ 1.txt

	rootDir := testUtils.CreateTmpDir()

	for _, subDir := range []string{"/a/a", "/b/a"} {
		testUtils.CreateDir(rootDir + subDir)
	}
	for _, file := range []string{"/a/a/0.txt", "/a/a/1.txt", "/b/0.txt", "/b/1.txt", "/b/a/1.txt"} {
		testUtils.CreateFile(rootDir, file, "")
	}
	return rootDir
}

func createEqualFileAndDir() string {
	// ./
	// └─ a/ <- directory
	//    └─ a <- file

	rootDir := testUtils.CreateTmpDir()

	testUtils.CreateDir(rootDir + "/a")
	testUtils.CreateFile(rootDir, "/a/a", "")

	return rootDir
}
