package cmd

import (
	"fmt"
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
	Run: setupRun,
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

func setupRun(cmd *cobra.Command, args []string) {
	log.Debug("setup.run()")

	naming, err := cmd.Flags().GetString("naming")
	if err != nil {
		log.Errorf("Failed to get 'naming' flag: %v", err)
		return
	}

	groupBy, err := cmd.Flags().GetString("group-by")
	if err != nil {
		log.Errorf("Failed to get 'group-by' flag: %v", err)
		return
	}

	table, err := cmd.Flags().GetString("table")
	if err != nil {
		log.Errorf("Failed to get 'table' flag: %v", err)
		return
	}

	if utils.DryRunFlag {

		yellow := color.New(color.FgYellow).SprintFunc()
		grey := color.New(color.FgHiBlack).SprintFunc()

		fmt.Println()
		fmt.Println("Simulated repository names:")

		groupDataMap, err := ars.GroupAndNameRepositories(
			ars.WithTablePath(table),
			ars.WithNamingPattern(naming),
			ars.WithGroupBy(groupBy),
		)
		if err != nil {
			fmt.Println("Error grouping and naming repositories:", err)
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
	} else {
		log.Error("Error: 'setup' command is not yet implemented")
	}
}
