package fileUtils

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"github.com/apex/log"
	"github.com/joho/godotenv"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Searches recursively for full path(es) of a given filename. Returns a 1-elem
// array if there is just one occurrence, or an array with several elements otherwise.
func FindFilesInDirRecursively(rootDir string, justTheFileName string) ([]string, error) {
	var targetPaths []string

	err := filepath.Walk(rootDir, func(currentPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Check if the current file matches the target file name
		if info.Mode().IsRegular() && info.Name() == justTheFileName {
			targetPaths = append(targetPaths, currentPath)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("Error searching for files: %w", err)
	}
	return targetPaths, nil
}

// Same as FindFilesInDirRecursively, but searches for any files (not just the ones with the given names)
func FindAnyFilesInDirRecursively(rootDir string) ([]string, error) {
	var targetPaths []string

	err := filepath.Walk(rootDir, func(currentPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Mode().IsRegular() {
			targetPaths = append(targetPaths, currentPath)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("Error searching for files: %w", err)
	}
	return targetPaths, nil
}

// same as FindFilesInDirRecursively, but without the recursive descent
func FindFilesInDir(rootDir string, justTheFileName string) ([]string, error) {
	files, err := ioutil.ReadDir(rootDir)
	if err != nil {
		return nil, err
	}
	filePaths := make([]string, 0)
	for _, file := range files {
		if !file.IsDir() && file.Name() == justTheFileName {
			filePaths = append(filePaths, filepath.Join(rootDir, file.Name()))
			break
		}
	}
	return filePaths, nil
}

// Transforms an absolute path into a relative paths, relative to a given root
func TransformIntoRelativePaths(root string, absPath string) (string, error) {
	log.Debug("utils.TransformIntoRelativePaths() with root: " + root)
	relPath, err := filepath.Rel(root, absPath)
	if err != nil {
		return "", err
	}
	return relPath, nil
}

func CopyFile(srcFileName, destDirName string) error {
	var file os.FileInfo
	var err error

	if file, err = os.Stat(srcFileName); err != nil {
		return err
	}

	if file.IsDir() {
		return fmt.Errorf("%s can not be a directory", srcFileName)
	}

	srcFile, err := os.Open(srcFileName)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFullPath := filepath.Join(destDirName, filepath.Base(srcFileName))
	destFile, err := os.Create(destFullPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}
	err = destFile.Sync()
	if err != nil {
		return err
	}
	return nil
}

func CopyAllFilesInDir(srcDirName, destDirName string) error {
	return filepath.Walk(srcDirName, func(srcPath string, info fs.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("failed to access path %s: %w", srcPath, err)
		}

		relPath, err := filepath.Rel(srcDirName, srcPath)
		if err != nil {
			return fmt.Errorf("failed to get relative path for %s: %w", srcPath, err)
		}

		destPath := filepath.Join(destDirName, relPath)

		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}

		srcFile, err := os.Open(srcPath)
		if err != nil {
			return fmt.Errorf("failed to open source file %s: %w", srcPath, err)
		}
		defer srcFile.Close()

		destFile, err := os.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, info.Mode())
		if err != nil {
			return fmt.Errorf("failed to create destination file %s: %w", destPath, err)
		}
		defer destFile.Close()

		if _, err := io.Copy(destFile, srcFile); err != nil {
			return fmt.Errorf("failed to copy file content from %s to %s: %w", srcPath, destPath, err)
		}

		return nil
	})
}

// Check if a given file exists (and is a file)
func ValidateFilePath(filePath string) error {
	log.Debug("utils.ValidateFilePath()")
	return ValidatePath(false, filePath)
}

// Check if a given directory exists (and is a directory)
func ValidateDirPath(dirPath string) error {
	log.Debug("utils.ValidateDirPath()")
	return ValidatePath(true, dirPath)
}

// Check if a given path exists (and is a directory or a file, depending on shouldBeDir)
func ValidatePath(shouldBeDir bool, path string) error {
	log.Debug("utils.ValidatePath()")
	var fileInfo os.FileInfo
	var err error

	if fileInfo, err = os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("%s does not exist: %w", path, err)
		}
		return fmt.Errorf("Error while validating path: %w", err)
	}

	if shouldBeDir != fileInfo.Mode().IsDir() {
		errorMessage := "%s is a "
		if shouldBeDir {
			errorMessage += "file, not a directory"
		} else {
			errorMessage += "directory, not a file"
		}
		return fmt.Errorf(errorMessage, path)
	}

	return nil
}

// Validate a list of files
func ValidateAllFilePaths(paths ...string) error {
	log.Debug("utils.ValidateAllFilePaths()")
	return ValidateAllPaths(false, paths...)
}

// Validate a list of directories
func ValidateAllDirPaths(paths ...string) error {
	log.Debug("utils.ValidateAllDirPaths()")
	return ValidateAllPaths(true, paths...)
}

// Validate a list of paths (files or directories)
func ValidateAllPaths(shouldBeDir bool, paths ...string) error {
	log.Debug("utils.ValidateAllPaths()")
	var errMsg string
	for _, path := range paths {
		err := ValidatePath(shouldBeDir, path)
		if err != nil {
			errMsg += "\n" + err.Error()
		}
	}

	if errMsg != "" {
		return &InvalidPathError{errMsg}
	}

	return nil
}

func FindUniqueFileWithPrefix(dir, prefix string) (string, error) {
	if prefix == "" {
		return "", fmt.Errorf("Prefix cannot be empty")
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return "", fmt.Errorf("Error reading directory: %w", err)
	}

	matchingFiles := []string{}
	for _, file := range files {
		if !file.IsDir() && strings.HasPrefix(file.Name(), prefix) {
			matchingFiles = append(matchingFiles, file.Name())
		}
	}

	if len(matchingFiles) == 0 {
		return "", fmt.Errorf("No file found with prefix '%s' in directory '%s'", prefix, dir)
	}

	if len(matchingFiles) > 1 {
		return "", fmt.Errorf("Multiple files found with prefix '%s' in directory '%s'", prefix, dir)
	}

	return filepath.Join(dir, matchingFiles[0]), nil
}

func ListSubFolderNames(folderPath string) ([]string, error) {
	files, err := ioutil.ReadDir(folderPath)
	if err != nil {
		return nil, err
	}

	subFolders := make([]string, 0)

	for _, file := range files {
		if file.IsDir() {
			subFolders = append(subFolders, file.Name())
		}
	}

	return subFolders, nil
}

func DeepCopy(srcObject, destinationObject interface{}) error {
	jsonData, err := json.Marshal(srcObject)
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonData, destinationObject)
	if err != nil {
		return err
	}

	return nil
}

// CreateFile generates a file at a specified path, returning the file path. If desired, the file can be generated
// with content provided as an argument.
func CreateFile(path string, fileName string, fileContent string) {
	f, err := os.Create(path + "/" + fileName)
	if err != nil {
		log.Fatalf("Could not create a file: %v", err)
	}

	_, err = f.Write([]byte(fileContent))
	if err != nil {
		log.Fatalf("Could not write to file: %v", err)
	}

	err = f.Close()
	if err != nil {
		log.Fatalf("Could not close file: %v", err)
	}
}

// CreateTmpDir creates a directory in the temp folder and provides its path as a
// return value. It is the caller's responsibility to remove this folder when it is no longer needed.
func CreateTmpDir() string {
	path, err := os.MkdirTemp("", "divekit_cli_")
	if err != nil {
		log.Fatalf("Could not create directory: %v", err)
	}

	return path
}

// CreateDir creates a directory named path, along with any necessary parents. If path is already a directory,
// CreateDir does nothing.
func CreateDir(path string) {
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		log.Fatalf("Could not create directory: %v", err)
	}
}

// DeleteDir deletes a specified directory along with its files and subdirectories.
func DeleteDir(path string) {
	if err := os.RemoveAll(path); err != nil {
		log.Fatalf("Could not remove directory: %v", err)
	}
}
func ToRelPath(absPath string, root string) string {
	relPath, err := filepath.Rel(root, absPath)
	if err != nil {
		log.Fatalf("Could not convert an absolute path into a relative path: ", err)
	}
	return UnifyPath(relPath)
}
func ToRelPaths(absPaths []string, root string) []string {
	var result []string
	for _, absPath := range absPaths {
		result = append(result, ToRelPath(absPath, root))
	}
	return result
}

func GetBaseName(path string) string {
	return ToRelPath(path, filepath.Dir(path))
}

func GetBaseNames(paths ...string) []string {
	var result []string
	for _, path := range paths {
		result = append(result, GetBaseName(path))
	}
	return result
}

// UnifyPath replaces all `\\` with `/`, addressing the variations in path formats across different operating systems.
func UnifyPath(path string) string {
	return strings.ReplaceAll(path, "\\", "/")
}

func UnifyPaths(paths []string) []string {
	var result []string
	for _, path := range paths {
		result = append(result, UnifyPath(path))
	}

	return result
}
func GetSHA256(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf(fmt.Sprintf("could not open the file %s: %v", filePath, err))
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		log.Fatalf("", err)
	}

	return fmt.Sprintf("%x", hash.Sum(nil))
}

func GetHomePath() string {
	projectRootDir := GetProjectRootDir()
	baseDir := GetBaseName(projectRootDir)
	result := strings.TrimSuffix(projectRootDir, baseDir)

	return result
}

func GetProjectRootDir() string {
	bytes, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		log.Fatalf(fmt.Sprintf("Could not get the root dir of this project: %v", err))
	}

	return strings.TrimSpace(string(bytes))
}
func LoadEnv() {
	projectRootDir := GetProjectRootDir()
	err := godotenv.Load(projectRootDir + "/.env")
	if err != nil {
		log.Fatalf("Error loading .env:", err)
	}
}

type InvalidPathError struct {
	Msg string
}

func (e *InvalidPathError) Error() string {
	return e.Msg
}
