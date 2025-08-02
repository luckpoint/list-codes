package utils

import (
	"testing"
)

// TestGetLanguageByExtension はGetLanguageByExtension関数のユニットテストです。
func TestGetLanguageByExtension(t *testing.T) {
	tests := []struct {
		name     string
		fileName string
		expected string
	}{
		{name: "Python file", fileName: "script.py", expected: "Python"},
		{name: "JavaScript file", fileName: "app.js", expected: "Javascript"},
		{name: "Go file", fileName: "main.go", expected: "Go"},
		{name: "Dockerfile", fileName: "Dockerfile", expected: "Dockerfile"},
		{name: "Unknown extension", fileName: "document.xyz", expected: ""},
		{name: "No extension", fileName: "README", expected: ""},
		{name: "Case insensitive extension", fileName: "INDEX.HTML", expected: "HTML"},
		{name: "Multiple dots", fileName: "archive.tar.gz", expected: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := GetLanguageByExtension(tt.fileName)
			if actual != tt.expected {
				t.Errorf("GetLanguageByExtension(%q): expected %q, got %q", tt.fileName, tt.expected, actual)
			}
		})
	}
}



