package gitlabapi

import (
	"divekit-cli/divekit/ars"
	"fmt"
	"log"

	"github.com/xanzy/go-gitlab"
)

type GitLab struct {
	Client *gitlab.Client
}

func NewGitLabClient(token, remote string) *GitLab {
	client, err := gitlab.NewClient(token, gitlab.WithBaseURL(remote))
	if err != nil {
		log.Fatalf("Error creating GitLab client: %v", err)
	}
	return &GitLab{Client: client}
}

func (g *GitLab) UserExists(username string) (*gitlab.User, bool) {
	users, _, err := g.Client.Users.ListUsers(&gitlab.ListUsersOptions{Username: &username})
	if err != nil || len(users) == 0 {
		return nil, false
	}
	return users[0], true
}

func (g *GitLab) CreateOnlineRepositories(groupDataMap map[string]*ars.GroupData, configContent ars.RepositoryConfigContentType) {
	fmt.Println()
	fmt.Println("Creating repositories online...")
	testRepositoryTargetGroupID := configContent.Remote.TestRepositoryTargetGroupId

	for _, groupData := range groupDataMap {
		var validUsers []*gitlab.User
		for _, record := range groupData.Records {
			if username, ok := record["username"]; ok {
				if user, exists := g.UserExists(username); exists {
					validUsers = append(validUsers, user)
				}
			}
		}

		if len(validUsers) > 0 {
			repoName := groupData.Name
			project, _, err := g.Client.Projects.CreateProject(&gitlab.CreateProjectOptions{
				Name:        &repoName,
				NamespaceID: &testRepositoryTargetGroupID,
			})
			if err != nil {
				log.Fatalf("Error creating repository for %s: %v", repoName, err)
			}

			for _, user := range validUsers {
				accessLevel := gitlab.AccessLevelValue(gitlab.DeveloperPermissions)
				_, _, err := g.Client.ProjectMembers.AddProjectMember(project.ID, &gitlab.AddProjectMemberOptions{
					UserID:      &user.ID,
					AccessLevel: &accessLevel,
				})
				if err != nil {
					fmt.Printf("Failed to add user %s to project %s:\n\t%v\n", user.Username, repoName, err)
				}
			}

			fmt.Printf("Repository %s created successfully\n", repoName)
		} else {
			fmt.Printf("No valid users found for %s; skipping repository creation.\n", groupData.Name)
		}
	}

	fmt.Println("Repositories created successfully")
	fmt.Println()
}
