package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/apex/log"
	"github.com/spf13/cobra"
	"github.com/xanzy/go-gitlab"

	"divekit-cli/divekit/ars"
	"divekit-cli/divekit/origin"
	"divekit-cli/utils"
	"divekit-cli/utils/dye"
)

var (
	ShowDetails     bool
	originRepo      *origin.OriginRepoType
	distributionKey string
	token           string
	remote          string // e.g. https://git.archi-lab.io/
)

// setupCmd represents the setup command
var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "setup group repositories",
	Long: `Create several repositories for individual students or student groups. 
For example:

divekit setup --naming=praktikum-S{{now "2006"}}-{{.group}}-{{uuid}}-{{autoincrement}} --dry-run --details`,
	Run:    setupRun,
	PreRun: setupPreRun,
}

func init() {
	log.Debug("setup.init()")
	// setupCmd.Flags().StringP("naming", "n", "", "name template for the repositories to be created")
	// setupCmd.Flags().StringP("group-by", "g", "", "group by column name")
	// setupCmd.Flags().StringP("table", "t", "", "path to the table file")
	setupCmd.Flags().BoolVarP(&ShowDetails, "details", "", false, "Show detailed output for each group")
	setupCmd.Flags().StringVarP(&token, "token", "t", "", "GitLab token")
	setupCmd.Flags().StringVarP(&remote, "remote", "r", "", "Remote repository URL (GitLab Instance)")

	patchCmd.MarkPersistentFlagRequired("originrepo")
	rootCmd.AddCommand(setupCmd)
}

// Checks preconditions before running the command
func setupPreRun(cmd *cobra.Command, args []string) {
	ars.Repo = ars.NewARSRepo()
	if ars.Repo == nil || ars.Repo.Config.RepositoryConfigFile == nil || ars.Repo.Config.RepositoryConfigFile.ReadContent() != nil {
		log.Fatal("ARSRepo or its Config is not properly initialized or failed to load")
	}

	originRepoName, err := rootCmd.PersistentFlags().GetString("originrepo")
	if err != nil {
		log.Fatal("Failed to get originrepo flag")
	}

	distributionKey, err = rootCmd.PersistentFlags().GetString("distribution")
	if err != nil {
		log.Fatal("Failed to get distribution flag")
	}

	originRepo = origin.NewOriginRepo(originRepoName)
	if err := checkRepoConfig(originRepo, distributionKey); err != nil {
		log.Fatal(err.Error())
	}
}

func checkRepoConfig(repo *origin.OriginRepoType, distributionKey string) error {
	if repo == nil {
		return fmt.Errorf("OriginRepo is nil")
	}
	if repo.DistributionMap == nil {
		return fmt.Errorf("DistributionMap is nil")
	}
	distribution, ok := repo.DistributionMap[distributionKey]
	if !ok {
		return fmt.Errorf("distribution %s not found in Distribution Folder", distributionKey)
	}
	if distribution.RepositoryConfigFile == nil {
		return fmt.Errorf("no RepositoryConfigFile found for %s", distributionKey)
	}
	if distribution.RepositoryConfigFile.ReadContent() != nil {
		return fmt.Errorf("failed to read content from RepositoryConfigFile for %s", distributionKey)
	}
	return nil
}

func setupRun(cmd *cobra.Command, args []string) {
	log.Debug("setup.run()")

	// get the naming pattern
	if ars.Repo == nil || ars.Repo.Config.RepositoryConfigFile == nil {
		log.Fatal("ARSRepo or its Config is not properly initialized")
	}

	var configContent ars.RepositoryConfigContentType = originRepo.DistributionMap["test"].RepositoryConfigFile.Content

	groupDataMap, err := ars.NameGroupedRepositories(
		ars.WithGroups(configContent.Repository.RepositoryMembers),
		ars.WithNamingPattern(configContent.Repository.RepositoryName),
	)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Error naming repositories")
		return
	}

	if utils.DryRunFlag {
		fmt.Println(dye.Very.Red("Dry run mode enabled. No changes will be made."))

		printExample(groupDataMap, configContent)
		os.Exit(0)
	}

	createOnlineRepositories(groupDataMap, configContent)

	log.Error("Error: 'setup' command is not yet implemented")
	os.Exit(1)
}

func userExists(git *gitlab.Client, username string) (*gitlab.User, bool) {
	users, _, err := git.Users.ListUsers(&gitlab.ListUsersOptions{Username: &username})
	if err != nil || len(users) == 0 {
		return nil, false
	}
	return users[0], true
}

func getNamespaceID(git *gitlab.Client, namespaceName string) (int, error) {
	groups, _, err := git.Groups.ListGroups(&gitlab.ListGroupsOptions{Search: &namespaceName})
	if err != nil {
		return 0, err
	}
	for _, group := range groups {
		if group.Name == namespaceName {
			return group.ID, nil
		}
	}
	return 0, fmt.Errorf("namespace not found")
}

func createOnlineRepositories(groupDataMap map[string]*ars.GroupData, configContent ars.RepositoryConfigContentType) {
	fmt.Println()
	fmt.Println("Creating repositories online...")

	git, err := gitlab.NewClient(token, gitlab.WithBaseURL(remote))
	if err != nil {
		log.Fatalf("Error creating GitLab client: %v", err)
	}

	testRepositoryTargetGroupID := configContent.Remote.TestRepositoryTargetGroupId

	for _, groupData := range groupDataMap {
		var validUsers []*gitlab.User
		for _, record := range groupData.Records {
			if username, ok := record["username"]; ok {
				if user, exists := userExists(git, username); exists {
					validUsers = append(validUsers, user)
				}
			}
		}

		if len(validUsers) > 0 {
			repoName := groupData.Name
			project, _, err := git.Projects.CreateProject(&gitlab.CreateProjectOptions{
				Name:        &repoName,
				NamespaceID: &testRepositoryTargetGroupID,
			})
			if err != nil {
				log.Fatalf("Error creating repository for %s: %v", repoName, err)
			}

			for _, user := range validUsers {
				accessLevel := gitlab.AccessLevelValue(gitlab.DeveloperPermissions)
				_, _, err := git.ProjectMembers.AddProjectMember(project.ID, &gitlab.AddProjectMemberOptions{
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

func printExample(groupDataMap map[string]*ars.GroupData, configContent ars.RepositoryConfigContentType) {
	fmt.Println()
	fmt.Println("Target IDs:")
	fmt.Println("\t", dye.Grey("LIVE:"), dye.Yellow(configContent.Remote.CodeRepositoryTargetGroupId))
	fmt.Println("\t", dye.Grey("TEST:"), dye.Yellow(configContent.Remote.TestRepositoryTargetGroupId))

	fmt.Println("Repository names:")

	for _, groupData := range groupDataMap {

		fmt.Printf("\t%s", dye.Yellow(groupData.Name))
		if ShowDetails {
			fmt.Println()
			for _, record := range groupData.Records {

				keys := make([]string, 0, len(record))
				for key := range record {
					keys = append(keys, key)
				}
				sort.Strings(keys)

				recordDetails := make([]string, 0, len(record))
				for _, key := range keys {
					value := record[key]
					recordDetails = append(recordDetails, fmt.Sprintf("%s %s", dye.Grey(fmt.Sprintf("%s:", key)), value))
				}
				fmt.Println("\t  ", strings.Join(recordDetails, ", "))
			}
		} else {
			fmt.Printf(" %s\n", dye.Grey(fmt.Sprintf("(%d rows)", len(groupData.Records))))
		}
	}

	fmt.Println()
}
