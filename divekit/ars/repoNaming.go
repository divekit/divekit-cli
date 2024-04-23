package ars

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"text/template"
	"unicode"

	"github.com/spf13/viper"
)

// GroupData contains the records and the name of a group
type GroupData struct {
	Records        []map[string]string
	RepositoryName string
}

// GroupOption is a function that modifies the GroupOptions
type GroupOption func(*GroupOptions)

// GroupOptions contains options for grouping and naming repositories
type GroupOptions struct {
	TablePath     string
	NamingPattern string
	GroupBy       string
	Groups        [][]string
}

type TemplateData struct {
	Usernames []string
	Group     string
}

// WithTablePath allows to provide a path to a table file
func WithTablePath(path string) GroupOption {
	return func(opts *GroupOptions) {
		opts.TablePath = path
	}
}

// WithNamingPattern allows to provide a naming pattern for the repositories
func WithNamingPattern(pattern string) GroupOption {
	return func(opts *GroupOptions) {
		opts.NamingPattern = pattern
	}
}

// WithGroupBy allows to provide a column name to group by
func WithGroupBy(groupBy string) GroupOption {
	return func(opts *GroupOptions) {
		opts.GroupBy = groupBy
	}
}

// WithGroups allows to provide grouped student ids directly
func WithGroups(groups [][]string) GroupOption {
	return func(opts *GroupOptions) {
		opts.Groups = groups
	}
}

// NameGroupedRepositories takes grouped student ids and applies a naming pattern
func NameGroupedRepositories(options ...GroupOption) (map[string]*GroupData, error) {
	opts := &GroupOptions{
		NamingPattern: viper.GetString("namingpattern"),
	}

	for _, option := range options {
		option(opts)
	}

	groups := make(map[string][]map[string]string)
	for _, group := range opts.Groups {
		var records []map[string]string
		for _, user := range group {
			records = append(records, map[string]string{"username": user})
		}
		id := userGroupIdentifier(group)
		groups[id] = records
	}

	return applyGroupingAndNaming(opts, groups)
}

func userGroupIdentifier(group []string) string {
	return strings.Join(group, "-")
}

// GroupAndNameRepositories groups students data and applies a naming pattern
func GroupAndNameRepositories(options ...GroupOption) (map[string]*GroupData, error) {
	opts := &GroupOptions{
		TablePath:     viper.GetString("table"),
		NamingPattern: viper.GetString("namingpattern"),
		GroupBy:       viper.GetString("groupBy"),
	}

	for _, option := range options {
		option(opts)
	}

	file, err := os.Open(opts.TablePath)
	if err != nil {
		return nil, fmt.Errorf("error opening table file at %s: %v", opts.TablePath, err)
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	header, err := csvReader.Read()
	if err != nil {
		return nil, fmt.Errorf("error reading header from table file: %v", err)
	}

	groups := make(map[string][]map[string]string)
	for {
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading record from table file: %v", err)
		}

		data := make(map[string]string)
		for i, value := range record {
			data[header[i]] = value
		}

		groupName := data[opts.GroupBy]
		groups[groupName] = append(groups[groupName], data)
	}

	return applyGroupingAndNaming(opts, groups)
}

func applyGroupingAndNaming(opts *GroupOptions, groups map[string][]map[string]string) (map[string]*GroupData, error) {
	groupDataMap := make(map[string]*GroupData)

	for groupName, records := range groups {
		fmt.Fprintf(os.Stderr, "Group %s has %d records\n", groupName, len(records))
		data := mapFromRecords(records, "")
		naming, err := applyDynamicTemplate(opts.NamingPattern, data)
		if err != nil {
			return nil, err
		}

		groupDataMap[naming] = &GroupData{
			Records:        records,
			RepositoryName: cleanGitLabProjectName(naming),
		}
	}

	return groupDataMap, nil
}

// mapFromGroup converts a group of student ids to a map with keys "username[0]", "username[1]", ...
func mapFromRecords(records []map[string]string, group string) TemplateData {
	if group == "" {
		group = "username"
	}

	usernames := make([]string, len(records))
	for i, record := range records {
		usernames[i] = record[group]
	}
	return TemplateData{Usernames: usernames, Group: group}
}

func applyDynamicTemplate(namingPattern string, data TemplateData) (string, error) {

	tmpl, err := template.New("naming").Funcs(template.FuncMap{
		"now":           Now,
		"creation":      Creation,
		"hash":          Hash,
		"uuid":          Uuid,
		"autoincrement": Autoincrement,
	}).Parse(namingPattern)
	if err != nil {
		return "", fmt.Errorf("template parsing failed: %w", err)
	}

	var result strings.Builder
	err = tmpl.Execute(&result, data)
	if err != nil {
		return "", fmt.Errorf("template execution failed: %w, data was: %+v", err, data)
	}

	return result.String(), nil
}

// cleanGitLabProjectName cleans up a project name for GitLab
func cleanGitLabProjectName(name string) string {
	name = replaceUmlauts(name)
	name = cleanUpIllegalCharacters(name)
	name = cleanUpHyphens(name)

	return name
}

// cleanUpHyphens removes multi hyphens and leading/trailing hyphens
func cleanUpHyphens(cleaned string) string {
	cleaned = regexp.MustCompile(`\-+`).ReplaceAllString(cleaned, "-")

	cleaned = strings.Trim(cleaned, "-")
	return cleaned
}

// cleanUpIllegalCharacters removes all characters that are not A-Z, a-z, 0-9 or a hyphen
func cleanUpIllegalCharacters(name string) string {
	reg := regexp.MustCompile(`[^a-zA-Z0-9\-]+`)
	cleaned := reg.ReplaceAllString(name, "-")
	return cleaned
}

// replaceUmlauts replaces umlauts with their equivalent ascii representation
func replaceUmlauts(input string) string {
	replacements := map[rune]string{
		'ä': "ae", 'ö': "oe", 'ü': "ue", 'ß': "ss",
		'Ä': "Ae", 'Ö': "Oe", 'Ü': "Ue", 'ẞ': "Ss",
	}
	result := []rune{}
	inputRunes := []rune(input)

	for i, r := range inputRunes {
		if repl, ok := replacements[r]; ok {
			var shouldCapitalize bool
			if i > 0 && unicode.IsLetter(inputRunes[i-1]) && unicode.IsUpper(inputRunes[i-1]) {
				shouldCapitalize = true
			}
			if i < len(inputRunes)-1 && unicode.IsLetter(inputRunes[i+1]) && unicode.IsUpper(inputRunes[i+1]) {
				shouldCapitalize = true
			}

			if unicode.IsUpper(r) && shouldCapitalize {
				repl = strings.ToUpper(repl)
			}

			result = append(result, []rune(repl)...)
		} else {
			result = append(result, r)
		}
	}

	return string(result)
}
