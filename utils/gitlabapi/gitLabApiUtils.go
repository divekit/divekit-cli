package gitlabapi

import (
	"divekit-cli/divekit/ars"
	"fmt"

	"github.com/apex/log"
	"github.com/xanzy/go-gitlab"
)

type GitLabClient interface {
	UserExists(username string) (*gitlab.User, bool, error)
	CreateOnlineRepositories(groupDataMap map[string]*ars.GroupData, configContent ars.RepositoryConfigContentType) error
}

type gitLabType struct {
	client *gitlab.Client
}

func NewGitLabClient(token, baseURL string) (GitLabClient, error) {
	client, err := gitlab.NewClient(token, gitlab.WithBaseURL(baseURL))
	if err != nil {
		return nil, fmt.Errorf("error creating GitLab client: %w", err)
	}
	return &gitLabType{client: client}, nil
}

func (g *gitLabType) UserExists(username string) (*gitlab.User, bool, error) {
	users, _, err := g.client.Users.ListUsers(&gitlab.ListUsersOptions{Username: &username})
	if err != nil {
		return nil, false, fmt.Errorf("error listing users: %w", err)
	}
	if len(users) == 0 {
		return nil, false, nil
	}
	return users[0], true, nil
}

func (g *gitLabType) CreateOnlineRepositories(groupDataMap map[string]*ars.GroupData, configContent ars.RepositoryConfigContentType) []error {
	errors := make([]error, 0)
	for _, groupData := range groupDataMap {
		var validUsers []*gitlab.User
		for _, record := range groupData.Records {
			username, ok := record["username"]
			if !ok {
				continue
			}
			user, exists, err := g.UserExists(username)
			if err != nil {
				errors = append(errors, err)
				continue
			}
			if exists {
				validUsers = append(validUsers, user)
			}
		}

		if len(validUsers) == 0 {
			log.Infof("No valid users found for %s; skipping repository creation.", groupData.Name)
			continue
		}

		repoName := groupData.Name
		project, _, err := g.client.Projects.CreateProject(&gitlab.CreateProjectOptions{
			Name:        &repoName,
			NamespaceID: &configContent.Remote.TestRepositoryTargetGroupId,
		})
		if err != nil {
			message := fmt.Sprintf("error creating repository for %s: %v", repoName, err)
			errors = append(errors, fmt.Errorf(message))
			log.Errorf("error creating repository for %s: %w", repoName, err)
		}

		for _, user := range validUsers {
			accessLevel := gitlab.AccessLevelValue(gitlab.DeveloperPermissions)
			_, _, err := g.client.ProjectMembers.AddProjectMember(project.ID, &gitlab.AddProjectMemberOptions{
				UserID:      &user.ID,
				AccessLevel: &accessLevel,
			})
			if err != nil {
				log.Errorf("Failed to add user %s to project %s:\n\t%v", user.Username, repoName, err)
			}
		}

		log.Infof("Repository %s created successfully", repoName)
	}

	return errors
}
