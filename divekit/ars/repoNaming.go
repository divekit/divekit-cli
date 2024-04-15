package ars

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
	"text/template"

	"github.com/spf13/viper"
)

// GroupData contains the records and the name of a group
type GroupData struct {
	Records []map[string]string
	Name    string
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

	// Default options
	opts := &GroupOptions{
		NamingPattern: viper.GetString("namingpattern"),
	}

	// Override options
	for _, option := range options {
		option(opts)
	}

	groupDataMap := make(map[string]*GroupData)
	for _, group := range opts.Groups {

		naming, err := applyDynamicTemplate(opts.NamingPattern, mapFromGroup(group))
		if err != nil {
			fmt.Println("Error applying naming pattern:", err)
			return nil, err
		}

		var records []map[string]string
		for _, user := range group {
			records = append(records, map[string]string{"username": user})
		}

		groupDataMap[naming] = &GroupData{
			Records: records,
			Name:    cleanGitLabProjectName(naming),
		}
	}

	return groupDataMap, nil
}

// mapFromGroup converts a group of student ids to a map with keys "username[0]", "username[1]", ...
func mapFromGroup(group []string) map[string]string {
	data := make(map[string]string)
	for i, value := range group {
		data[fmt.Sprintf("username[%d]", i)] = value
	}
	return data
}

// GroupAndNameRepositories groups students data and applies a naming pattern
func GroupAndNameRepositories(options ...GroupOption) (map[string]*GroupData, error) {

	// Default-Optionen
	opts := &GroupOptions{
		TablePath:     viper.GetString("table"),
		NamingPattern: viper.GetString("namingpattern"),
		GroupBy:       viper.GetString("groupBy"),
	}

	// Optionen überschreiben
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

	groupDataMap := make(map[string]*GroupData)
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
		if group, exists := groupDataMap[groupName]; exists {
			group.Records = append(group.Records, data)
		} else {
			naming, err := applyDynamicTemplate(opts.NamingPattern, data)
			if err != nil {
				fmt.Println("Error applying naming pattern:", err)
				return nil, err
			}
			groupDataMap[groupName] = &GroupData{
				Records: []map[string]string{data},
				Name:    cleanGitLabProjectName(naming),
			}
		}
	}

	return groupDataMap, nil
}

func applyDynamicTemplate(namingPattern string, data map[string]string) (string, error) {
	if data == nil {
		data = make(map[string]string)
	}

	tmpl, err := template.New("naming").Funcs(template.FuncMap{
		"now":           Now,
		"creation":      Creation,
		"hash":          Hash,
		"uuid":          Uuid,
		"autoincrement": Autoincrement,
	}).Parse(namingPattern)
	if err != nil {
		return "", err
	}

	var result strings.Builder
	err = tmpl.Execute(&result, data)
	if err != nil {
		return "", err
	}

	return result.String(), nil
}

func cleanGitLabProjectName(name string) string {
	// Konvertiere den String zuerst in Kleinbuchstaben
	name = strings.ToLower(name)

	// Ersetze deutsche Umlaute und ß durch ihre Äquivalente
	replacements := map[string]string{
		"ä": "ae",
		"ö": "oe",
		"ü": "ue",
		"ß": "ss",
	}
	for old, new := range replacements {
		name = strings.ReplaceAll(name, old, new)
	}

	// Ersetze alle unerwünschten Zeichen durch Bindestriche
	reg := regexp.MustCompile(`[^a-z0-9\-]+`)
	cleaned := reg.ReplaceAllString(name, "-")

	// Entferne mehrfache Bindestriche, die durch die Ersetzung entstanden sein könnten
	cleaned = regexp.MustCompile(`\-+`).ReplaceAllString(cleaned, "-")

	// Entferne Bindestriche am Anfang und Ende
	cleaned = strings.Trim(cleaned, "-")

	return cleaned
}
