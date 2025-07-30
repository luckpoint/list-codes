package utils

import (
	"os"
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
			expectedCounts: map[string]int{"Go": 1, "Javascript": 1, "Python": 1},
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

			actualLangs, actualCounts := DetectProjectLanguages(dir, true, resolvedIncludePaths, tt.excludeNames, resolvedExcludePaths, nil)

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

func TestDetectProjectLanguagesWithTestFiles(t *testing.T) {
	// Test that test files are properly excluded from language counts in fallback mode
	tempDir := t.TempDir()
	
	// Create a project where most files are test files
	testFiles := map[string]string{
		"app.py":          "print('main app')",
		"test_app.py":     "import unittest",
		"test_models.py":  "import unittest", 
		"test_views.py":   "import unittest",
		"test_utils.py":   "import unittest",
		"test_api.py":     "import unittest",
		"spec/test_integration.py": "import unittest",
		"tests/test_e2e.py": "import unittest",
		"tests/conftest.py": "import pytest", // This should also be excluded
		"pytest.ini":      "[pytest]", // Config file, not counted
	}
	
	for filePath, content := range testFiles {
		fullPath := filepath.Join(tempDir, filePath)
		err := os.MkdirAll(filepath.Dir(fullPath), 0755)
		if err != nil {
			t.Fatalf("Failed to create directory: %v", err)
		}
		err = os.WriteFile(fullPath, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create file %s: %v", filePath, err)
		}
	}

	includePaths := make(map[string]struct{})
	excludeNames := make(map[string]struct{})
	excludePaths := make(map[string]struct{})

	primaryLangs, fallbackCounts := DetectProjectLanguages(tempDir, true, includePaths, excludeNames, excludePaths, nil)

	// Should not detect Python as primary language because no dependency file exists
	if len(primaryLangs) > 0 {
		t.Errorf("Expected no primary languages (no dependency files), but got: %v", primaryLangs)
	}

	// Should fall back to counting, but test files should be excluded
	expectedPythonCount := 1 // Only app.py should be counted, test files should be excluded
	if fallbackCounts["Python"] != expectedPythonCount {
		t.Errorf("Expected Python count to be %d (excluding test files), but got %d. All files detected: %v", 
			expectedPythonCount, fallbackCounts["Python"], fallbackCounts)
	}

	// Verify that the count represents only the non-test files
	if fallbackCounts["Python"] >= 5 { // If it's 5 or more, test files weren't properly excluded
		t.Errorf("Test files appear to be included in language count. Expected 1, got %d", fallbackCounts["Python"])
	}
}

func TestDetectProjectLanguagesTestFileEdgeCase(t *testing.T) {
	// Test edge case where ALL Python files are test files
	tempDir := t.TempDir()
	
	testFiles := map[string]string{
		"test_only.py":    "import unittest",
		"spec_only.py":    "import unittest", 
		"main_test.go":    "package main", // Go test file
		"service_test.go": "package main", // Go test file
		"go.mod":          "module test",  // Go dependency file - makes Go primary
	}
	
	for filePath, content := range testFiles {
		fullPath := filepath.Join(tempDir, filePath)
		err := os.MkdirAll(filepath.Dir(fullPath), 0755)
		if err != nil {
			t.Fatalf("Failed to create directory: %v", err)
		}
		err = os.WriteFile(fullPath, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create file %s: %v", filePath, err)
		}
	}

	includePaths := make(map[string]struct{})
	excludeNames := make(map[string]struct{})
	excludePaths := make(map[string]struct{})

	primaryLangs, fallbackCounts := DetectProjectLanguages(tempDir, true, includePaths, excludeNames, excludePaths, nil)

	// Go should be detected as primary due to go.mod
	expectedPrimary := []string{"Go"}
	if len(primaryLangs) != 1 || primaryLangs[0] != "Go" {
		t.Errorf("Expected primary languages to be %v, but got: %v", expectedPrimary, primaryLangs)
	}

	// Python count should be 0 since all Python files are test files
	if fallbackCounts["Python"] != 0 {
		t.Errorf("Expected Python count to be 0 (all files are test files), but got %d", fallbackCounts["Python"])
	}
}