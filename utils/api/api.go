package api

import (
	"divekit-cli/utils/fileUtils"
	"fmt"
	"github.com/xanzy/go-gitlab"
	"os"
)

var branchName string

func init() {
	fileUtils.LoadEnv()
	branchName = os.Getenv("DIVEKIT_MAINBRANCH_NAME")
}

func NewGitlabClient(host string, token string) (*gitlab.Client, error) {
	git, err := gitlab.NewClient(token, gitlab.WithBaseURL(host))
	if err != nil {
		return nil, fmt.Errorf("could not create a new gitlab client: %w", err)
	}

	return git, nil
}

func GetRepositoryById(client *gitlab.Client, repoId string) (*gitlab.Project, error) {
	project, _, err := client.Projects.GetProject(repoId, nil)
	if err != nil {
		return nil, fmt.Errorf("could not get project with repoId: %s: %w", repoId, err)
	}

	return project, nil
}
func GetRepositoriesByGroupId(client *gitlab.Client, groupId string) ([]*gitlab.Project, error) {
	projects, _, err := client.Groups.ListGroupProjects(groupId, nil)
	if err != nil {
		return nil, fmt.Errorf("could not get projects with groupId: %s: %w", groupId, err)
	}

	return projects, nil
}

func GetFileByRepositoryId(client *gitlab.Client, repoId string, filePath string) (*gitlab.File, error) {
	option := &gitlab.GetFileOptions{Ref: gitlab.Ptr(branchName)}
	file, _, err := client.RepositoryFiles.GetFile(repoId, filePath, option)
	if err != nil {
		return nil, fmt.Errorf("could not get file: %w", err)
	}

	return file, nil
}
func DeleteFileByRepositoryId(client *gitlab.Client, repoId string, filePath string) error {
	option := &gitlab.DeleteFileOptions{
		Branch:        gitlab.Ptr(branchName),
		CommitMessage: gitlab.Ptr("Prepare test [delete]"),
	}
	_, err := client.RepositoryFiles.DeleteFile(repoId, filePath, option)

	return err
}
func GetCommitsByRepositoryId(client *gitlab.Client, repoId string) ([]*gitlab.Commit, error) {
	option := &gitlab.ListCommitsOptions{
		RefName: gitlab.Ptr(branchName),
		All:     gitlab.Ptr(false),
	}
	commits, _, err := client.Commits.ListCommits(repoId, option)
	if err != nil {
		return nil, fmt.Errorf("could not get commits: %w", err)
	}

	return commits, nil
}
func RevertCommitByRepositoryIdAndCommitId(client *gitlab.Client, repoId string, commitId string) error {
	option := &gitlab.RevertCommitOptions{
		Branch: gitlab.Ptr(branchName),
	}
	_, _, err := client.Commits.RevertCommit(repoId, commitId, option)

	return err
}
