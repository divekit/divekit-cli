package ars

import (
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMapFromRecords(t *testing.T) {
	records := []map[string]string{
		{"username": "alice"},
		{"username": "bob"},
		{"username": "john"},
	}
	expected := TemplateData{Usernames: []string{"alice", "bob", "john"}, Group: "username"}
	result := mapFromRecords(records, "username")
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestApplyDynamicTemplate(t *testing.T) {
	templateString := "{{index .Usernames 0}}-project"
	data := TemplateData{Usernames: []string{"john"}}
	expected := "john-project"
	result, err := applyDynamicTemplate(templateString, data)
	if err != nil {
		t.Fatalf("Error should not have occurred: %v", err)
	}
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestReplaceUmlauts(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Replace Umlauts",
			input:    "Gänsefüßchen-Ölprüfer",
			expected: "Gaensefuesschen-Oelpruefer",
		},
		{
			name:     "Replace Umlaut At End Of All Caps String",
			input:    "AEIOÜ",
			expected: "AEIOUE",
		},
		{
			name:     "Replace Umlaut At Beginning Of All Caps String",
			input:    "ÄEIOU",
			expected: "AEEIOU",
		},
		{
			name:     "Replace Umlaut Inside Of All Caps String",
			input:    "AEIÖU",
			expected: "AEIOEU",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := replaceUmlauts(test.input)
			if result != test.expected {
				t.Errorf("Expected %s, got %s", test.expected, result)
			}
		})
	}
}

func TestCleanUpHyphens(t *testing.T) {
	name := "test----project---name----"
	expected := "test-project-name"
	result := cleanUpHyphens(name)
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestCleanUpIllegalCharacters(t *testing.T) {
	name := "test@#project$%^&name"
	expected := "test-project-name"
	result := cleanUpIllegalCharacters(name)
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestCleanGitLabProjectName(t *testing.T) {
	name := "Projekt-ÄÖÜß@@@***"
	expected := "Projekt-AEOEUEss"
	result := cleanGitLabProjectName(name)
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestUserGroupIdentifier(t *testing.T) {
	group := []string{"alice", "bob"}
	expected := "alice-bob"
	result := userGroupIdentifier(group)
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestNameGroupedRepositories(t *testing.T) {
	tests := []struct {
		name        string
		options     []GroupOption
		expected    map[string]*GroupData
		expectError bool
	}{
		{
			name: "single group",
			options: []GroupOption{
				WithNamingPattern("group-{{index .Usernames 0}}-{{index .Usernames 1}}"),
				WithGroups([][]string{{"alice", "bob"}}),
			},
			expected: map[string]*GroupData{
				"group-alice-bob": &GroupData{
					Records: []map[string]string{
						{"username": "alice"},
						{"username": "bob"},
					},
					RepositoryName: "group-alice-bob",
				},
			},
			expectError: false,
		},
		{
			name: "multiple groups",
			options: []GroupOption{
				WithNamingPattern("group-{{index .Usernames 0}}-{{index .Usernames 1}}"),
				WithGroups([][]string{{"alice", "bob"}, {"charlie", "dave"}}),
			},
			expected: map[string]*GroupData{
				"group-alice-bob": &GroupData{
					Records: []map[string]string{
						{"username": "alice"},
						{"username": "bob"},
					},
					RepositoryName: "group-alice-bob",
				},
				"group-charlie-dave": &GroupData{
					Records: []map[string]string{
						{"username": "charlie"},
						{"username": "dave"},
					},
					RepositoryName: "group-charlie-dave",
				},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NameGroupedRepositories(tt.options...)
			if (err != nil) != tt.expectError {
				t.Fatalf("NameGroupedRepositories() error = %v, expectError %v", err, tt.expectError)
			}
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestGroupAndNameRepositories(t *testing.T) {
	tests := []struct {
		name        string
		options     []GroupOption
		expected    map[string]*GroupData
		expectError bool
	}{
		{
			name: "Simple Grouping",
			options: []GroupOption{
				WithTablePath("testdata/grouping.csv"),
				WithNamingPattern("{{.Group}}-project"),
				WithGroupBy("group"),
			},
			expected: map[string]*GroupData{
				"A-project": {
					Records: []map[string]string{
						{"campusID": "alice", "group": "A"},
						{"campusID": "bob", "group": "A"},
					},
					RepositoryName: "A-project",
				},
				"B-project": {
					Records: []map[string]string{
						{"campusID": "john", "group": "B"},
					},
					RepositoryName: "B-project",
				},
			},
			expectError: false,
		},
		{
			name: "Invalid Naming Pattern",
			options: []GroupOption{
				WithTablePath("testdata/grouping.csv"),
				WithNamingPattern("{{.unknown}}-project"),
				WithGroupBy("group"),
			},
			expected:    nil,
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result, err := GroupAndNameRepositories(tc.options...)
			if (err != nil) != tc.expectError {
				t.Fatalf("Expected error: %v, got: %v", tc.expectError, err)
			}
			if !tc.expectError {
				if diff := cmp.Diff(tc.expected, result); diff != "" {
					t.Errorf("Mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}
