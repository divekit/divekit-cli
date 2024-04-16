package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/apex/log"
	"github.com/spf13/cobra"

	"divekit-cli/divekit/ars"
	"divekit-cli/divekit/origin"
	"divekit-cli/utils"
	"divekit-cli/utils/dye"
)

var (
	ShowDetails     bool
	originRepo      *origin.OriginRepoType
	distributionKey string
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

	log.Error("Error: 'setup' command is not yet implemented")
	os.Exit(1)
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
