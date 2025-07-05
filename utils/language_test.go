package utils

import (
	"path/filepath"
	"sort"
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

// TestDetectProjectLanguages はDetectProjectLanguages関数のユニットテストです。
func TestDetectProjectLanguages(t *testing.T) {
	tests := []struct {
		name           string
		setupFiles     map[string]string
		includePaths   map[string]struct{}
		excludeNames   map[string]struct{}
		excludePaths   map[string]struct{}
		expectedLangs  []string
		expectedCounts map[string]int
	}{
		{
			name: "Go project with go.mod",
			setupFiles: map[string]string{
				"go.mod":  "module myapp",
				"main.go": "package main",
			},
			expectedLangs:  []string{"Go"},
			expectedCounts: nil,
		},
		{
			name: "Python project with pyproject.toml",
			setupFiles: map[string]string{
				"pyproject.toml": "[project]",
				"app.py":         "print('hello')",
			},
			expectedLangs:  []string{"Python"},
			expectedCounts: nil,
		},
		{
			name: "Mixed project with fallback",
			setupFiles: map[string]string{
				"file.js":  "",
				"file.py":  "",
				"test.py":  "", // Should be excluded from counts
				"notes.md": "",
			},
			expectedLangs:  nil,
			expectedCounts: map[string]int{"Javascript": 1, "Python": 1, "Markdown": 1},
		},
		{
			name: "Project with excluded directory",
			setupFiles: map[string]string{
				"go.mod":              "module test",
				"main.go":             "",
				"node_modules/dep.js": "",
			},
			excludeNames:   map[string]struct{}{"node_modules": {}},
			expectedLangs:  []string{"Go"},
			expectedCounts: nil,
		},
		{
			name: "Project with excluded path",
			setupFiles: map[string]string{
				"go.mod":           "module test",
				"main.go":          "",
				"src/temp/file.js": "",
			},
			excludePaths:   map[string]struct{}{"src/temp": {}},
			expectedLangs:  []string{"Go"},
			expectedCounts: nil,
		},
		{
			name: "Project with included path",
			setupFiles: map[string]string{
				"main.go":         "",
				"src/app/file.js": "",
				"src/lib/file.py": "",
			},
			includePaths:   map[string]struct{}{"src/app": {}},
			expectedLangs:  nil,
			expectedCounts: map[string]int{"Javascript": 1},
		},
		{
			name:           "Empty project",
			setupFiles:     map[string]string{},
			expectedLangs:  nil,
			expectedCounts: map[string]int{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir, cleanup := setupTestDir(t)
			defer cleanup()

			for path, content := range tt.setupFiles {
				createTestFile(t, filepath.Join(dir, path), content)
			}

			// excludePathsとincludePathsの絶対パスを解決
			resolvedExcludePaths := make(map[string]struct{})
			for p := range tt.excludePaths {
				absPath, err := filepath.Abs(filepath.Join(dir, p))
				if err != nil {
					t.Fatalf("Failed to resolve absolute path for exclude: %v", err)
				}
				resolvedExcludePaths[absPath] = struct{}{}
			}

			resolvedIncludePaths := make(map[string]struct{})
			for p := range tt.includePaths {
				absPath, err := filepath.Abs(filepath.Join(dir, p))
				if err != nil {
					t.Fatalf("Failed to resolve absolute path for include: %v", err)
				}
				resolvedIncludePaths[absPath] = struct{}{}
			}

			actualLangs, actualCounts := DetectProjectLanguages(dir, true, resolvedIncludePaths, tt.excludeNames, resolvedExcludePaths)

			if tt.expectedLangs != nil {
				if len(actualLangs) != len(tt.expectedLangs) {
					t.Errorf("DetectProjectLanguages() primary langs count mismatch: expected %d, got %d. Expected: %v, Got: %v", len(tt.expectedLangs), len(actualLangs), tt.expectedLangs, actualLangs)
				} else {
					sort.Strings(actualLangs)
					sort.Strings(tt.expectedLangs)
					for i := range actualLangs {
						if actualLangs[i] != tt.expectedLangs[i] {
							t.Errorf("DetectProjectLanguages() primary langs mismatch: expected %v, got %v", tt.expectedLangs, actualLangs)
							break
						}
					}
				}
			} else if tt.expectedCounts != nil {
				if len(actualCounts) != len(tt.expectedCounts) {
					t.Errorf("DetectProjectLanguages() fallback counts count mismatch: expected %d, got %d. Expected: %v, Got: %v", len(tt.expectedCounts), len(actualCounts), tt.expectedCounts, actualCounts)
				} else {
					for lang, count := range tt.expectedCounts {
						if actualCounts[lang] != count {
							t.Errorf("DetectProjectLanguages() fallback count for %s mismatch: expected %d, got %d", lang, count, actualCounts[lang])
						}
					}
				}
			} else {
				t.Errorf("Test case misconfigured: either expectedLangs or expectedCounts must be set.")
			}
		})
	}
}