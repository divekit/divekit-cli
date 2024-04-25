package ars

import (
	"testing"
)

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
