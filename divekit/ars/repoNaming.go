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

type GroupData struct {
	Records []map[string]string
	Name    string
}

type GroupOption func(*GroupOptions)

type GroupOptions struct {
	TablePath     string
	NamingPattern string
	GroupBy       string
}

func WithTablePath(path string) GroupOption {
	return func(opts *GroupOptions) {
		opts.TablePath = path
	}
}

func WithNamingPattern(pattern string) GroupOption {
	return func(opts *GroupOptions) {
		opts.NamingPattern = pattern
	}
}

func WithGroupBy(groupBy string) GroupOption {
	return func(opts *GroupOptions) {
		opts.GroupBy = groupBy
	}
}

// GroupAndNameRepositories groups students data and applies a naming pattern
func GroupAndNameRepositories(options ...GroupOption) (map[string]*GroupData, error) {

	// Default-Optionen
	opts := &GroupOptions{
		TablePath:     viper.GetString("distribution.table"),
		NamingPattern: viper.GetString("distribution.namingpattern"),
		GroupBy:       viper.GetString("distribution.groupBy"),
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
