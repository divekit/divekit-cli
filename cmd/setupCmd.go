package cmd

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/apex/log"
	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"divekit-cli/divekit/ars"
	"divekit-cli/utils"
)

var (
	ShowDetails bool
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
	setupCmd.Flags().StringP("naming", "n", "", "name template for the repositories to be created")
	setupCmd.Flags().StringP("group-by", "g", "", "group by column name")
	setupCmd.Flags().StringP("table", "t", "", "path to the table file")
	setupCmd.Flags().BoolVarP(&ShowDetails, "details", "d", false, "Show detailed output for each group")

	patchCmd.MarkPersistentFlagRequired("originrepo")
	rootCmd.AddCommand(setupCmd)
}

// Checks preconditions before running the command
func setupPreRun(cmd *cobra.Command, args []string) {
	ars.Repo = ars.NewARSRepo()
	if ars.Repo == nil || ars.Repo.Config.RepositoryConfigFile == nil || ars.Repo.Config.RepositoryConfigFile.ReadContent() != nil {
		log.Fatal("ARSRepo or its Config is not properly initialized or failed to load")
	}
}

func setupRun(cmd *cobra.Command, args []string) {
	log.Debug("setup.run()")

	// get the naming pattern
	if ars.Repo == nil || ars.Repo.Config.RepositoryConfigFile == nil {
		log.Fatal("ARSRepo or its Config is not properly initialized")
	}

	if utils.DryRunFlag {

		yellow := color.New(color.FgYellow).SprintFunc()
		grey := color.New(color.FgHiBlack).SprintFunc()

		fmt.Println()
		fmt.Println("Simulated repository names:")

		groupDataMap, err := ars.NameGroupedRepositories(
			ars.WithGroups(ars.Repo.Config.RepositoryConfigFile.Content.Repository.RepositoryMembers),
			ars.WithNamingPattern(ars.Repo.Config.RepositoryConfigFile.Content.Repository.RepositoryName),
		)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Error("Error naming repositories")
			return
		}

		for _, groupData := range groupDataMap {
			// print group name
			fmt.Printf("\t%s", yellow(groupData.Name))
			if ShowDetails {
				fmt.Println()
				for _, record := range groupData.Records {
					// collect keys
					keys := make([]string, 0, len(record))
					for key := range record {
						keys = append(keys, key)
					}
					sort.Strings(keys)

					// collect details
					recordDetails := make([]string, 0, len(record))
					for _, key := range keys {
						value := record[key]
						recordDetails = append(recordDetails, fmt.Sprintf("%s %s", grey(fmt.Sprintf("%s:", key)), value))
					}
					fmt.Println("\t  ", strings.Join(recordDetails, ", "))
				}
			} else {
				fmt.Printf(" %s\n", grey(fmt.Sprintf("(%d rows)", len(groupData.Records))))
			}
		}
		fmt.Println()
		os.Exit(0)
	} else {
		log.Error("Error: 'setup' command is not yet implemented")
		os.Exit(1)
	}
}
