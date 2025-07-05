package utils

import (
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// TestIsTestFile はIsTestFile関数のユニットテストです。
func TestIsTestFile(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		expected bool
	}{
		{name: "Python test file", filePath: "my_test.py", expected: true},
		{name: "Go test file", filePath: "my_test.go", expected: true},
		{name: "JavaScript spec file", filePath: "my_spec.js", expected: true},
		{name: "TypeScript test file", filePath: "test.ts", expected: true},
		{name: "Test directory Python", filePath: "tests/my_module/file.py", expected: true},
		{name: "Test directory Go", filePath: "pkg/my_package/test/file.go", expected: true},
		{name: "Spec directory JS", filePath: "src/components/spec/component.js", expected: true},
		{name: "Solidity test pattern", filePath: "MyContract.t.sol", expected: true},
		{name: "Normal Python file", filePath: "main.py", expected: false},
		{name: "Normal Go file", filePath: "main.go", expected: false},
		{name: "Normal JS file", filePath: "app.js", expected: false},
		{name: "Case insensitive keyword", filePath: "TEST_file.py", expected: true},
		{name: "Keyword in middle", filePath: "file_test_data.py", expected: true},
		{name: "Keyword as prefix", filePath: "testdata.txt", expected: true},
		{name: "Keyword as suffix", filePath: "data_spec.json", expected: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := IsTestFile(tt.filePath)
			if actual != tt.expected {
				t.Errorf("IsTestFile(%q): expected %v, got %v", tt.filePath, tt.expected, actual)
			}
		})
	}
}

// TestShouldSkipDir はshouldSkipDir関数のユニットテストです。
func TestShouldSkipDir(t *testing.T) {
	dir, cleanup := setupTestDir(t)
	defer cleanup()

	// テスト用のディレクトリとファイルを作成
	createTestFile(t, filepath.Join(dir, "regular.txt"), "content")
	createTestFile(t, filepath.Join(dir, ".hidden.txt"), "content")
	createTestFile(t, filepath.Join(dir, ".git/config"), "content")
	createTestFile(t, filepath.Join(dir, "node_modules/package/index.js"), "content")
	createTestFile(t, filepath.Join(dir, "src/main.go"), "content")

	tests := []struct {
		name         string
		fullPath     string
		nameOnly     string
		isDir        bool
		includePaths map[string]struct{}
		excludeNames map[string]struct{}
		excludePaths map[string]struct{}
		expected     bool
	}{
		{
			name:     "Regular file",
			fullPath: filepath.Join(dir, "regular.txt"),
			nameOnly: "regular.txt",
			isDir:    false,
			expected: false,
		},
		{
			name:     "Hidden file",
			fullPath: filepath.Join(dir, ".hidden.txt"),
			nameOnly: ".hidden.txt",
			isDir:    false,
			expected: true,
		},
		{
			name:     "Git directory",
			fullPath: filepath.Join(dir, ".git"),
			nameOnly: ".git",
			isDir:    true,
			expected: true,
		},
		{
			name:         "Node modules directory",
			fullPath:     filepath.Join(dir, "node_modules"),
			nameOnly:     "node_modules",
			isDir:        true,
			excludeNames: map[string]struct{}{"node_modules": {}},
			expected:     true,
		},
		{
			name:         "Excluded by name",
			fullPath:     filepath.Join(dir, "build"),
			nameOnly:     "build",
			isDir:        true,
			excludeNames: map[string]struct{}{"build": {}},
			expected:     true,
		},
		{
			name:         "Excluded by path",
			fullPath:     filepath.Join(dir, "src/main.go"),
			nameOnly:     "main.go",
			isDir:        false,
			excludePaths: map[string]struct{}{filepath.Join(dir, "src/main.go"): {}},
			expected:     true,
		},
		{
			name:         "Included path - file in included directory",
			fullPath:     filepath.Join(dir, "src/main.go"),
			nameOnly:     "main.go",
			isDir:        false,
			includePaths: map[string]struct{}{filepath.Join(dir, "src"): {}},
			expected:     false,
		},
		{
			name:         "Not included path - file outside included directory",
			fullPath:     filepath.Join(dir, "regular.txt"),
			nameOnly:     "regular.txt",
			isDir:        false,
			includePaths: map[string]struct{}{filepath.Join(dir, "src"): {}},
			expected:     true,
		},
		{
			name:         "Parent of included path",
			fullPath:     dir,
			nameOnly:     filepath.Base(dir),
			isDir:        true,
			includePaths: map[string]struct{}{filepath.Join(dir, "src"): {}},
			expected:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := shouldSkipDir(tt.fullPath, tt.nameOnly, tt.isDir, tt.includePaths, tt.excludeNames, tt.excludePaths)
			if result != tt.expected {
				t.Errorf("shouldSkipDir(%q, %q, %v, %v, %v, %v) = %v, want %v",
					tt.fullPath, tt.nameOnly, tt.isDir, tt.includePaths, tt.excludeNames, tt.excludePaths, result, tt.expected)
			}
		})
	}
}

// TestGenerateDirectoryStructure はGenerateDirectoryStructure関数のユニットテストです。
func TestGenerateDirectoryStructure(t *testing.T) {
	tests := []struct {
		name           string
		setupFiles     map[string]string
		maxDepth       int
		includePaths   map[string]struct{}
		excludeNames   map[string]struct{}
		excludePaths   map[string]struct{}
		expectedOutput string
	}{
		{
			name: "Simple structure",
			setupFiles: map[string]string{
				"file1.txt":              "",
				"dir1/file2.txt":         "",
				"dir1/subdir1/file3.txt": "",
			},
			maxDepth: 0,
			expectedOutput: "## Project Structure\n" +
				"```text\n" +
				". (test_project)\n" +
				"├── dir1\n" +
				"│   ├── file2.txt\n" +
				"│   └── subdir1\n" +
				"│       └── file3.txt\n" +
				"└── file1.txt\n" +
				"```\n\n",
		},
		{
			name: "Max depth 1",
			setupFiles: map[string]string{
				"file1.txt":              "",
				"dir1/file2.txt":         "",
				"dir1/subdir1/file3.txt": "",
			},
			maxDepth: 1,
			expectedOutput: "## Project Structure\n" +
				"```text\n" +
				". (test_project)\n" +
				"├── dir1\n" +
				"│   └── ...\n" +
				"└── file1.txt\n" +
				"```\n\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir, cleanup := setupTestDir(t)
			defer cleanup()

			for path, content := range tt.setupFiles {
				createTestFile(t, filepath.Join(dir, path), content)
			}

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

			actualOutput := GenerateDirectoryStructure(dir, tt.maxDepth, true, resolvedIncludePaths, tt.excludeNames, resolvedExcludePaths)
			actualOutput = strings.ReplaceAll(actualOutput, filepath.Base(dir), "test_project")
			expectedOutput := strings.ReplaceAll(tt.expectedOutput, "test_project", "test_project")

			if !reflect.DeepEqual(strings.Split(actualOutput, "\n"), strings.Split(expectedOutput, "\n")) {
				t.Errorf("GenerateDirectoryStructure() output mismatch:\nExpected:\n%s\nActual:\n%s", expectedOutput, actualOutput)
			}
		})
	}
}

// TestCollectReadmeFiles はCollectReadmeFiles関数のユニットテストです。
func TestCollectReadmeFiles(t *testing.T) {
	tests := []struct {
		name           string
		setupFiles     map[string]string
		includePaths   map[string]struct{}
		excludeNames   map[string]struct{}
		excludePaths   map[string]struct{}
		expectedOutput string
	}{
		{
			name: "Single README.md in root",
			setupFiles: map[string]string{
				"README.md": "# Project Readme",
			},
			expectedOutput: "# Project README Files\n\n" +
				"### README.md\n\n" +
				"```markdown\n# Project Readme\n```\n",
		},
		{
			name: "Multiple README.md files",
			setupFiles: map[string]string{
				"README.md":      "# Root Readme",
				"docs/README.md": "## Docs Readme",
			},
			expectedOutput: "# Project README Files\n\n" +
				"### README.md\n\n" +
				"```markdown\n# Root Readme\n```\n\n" +
				"### docs/README.md\n\n" +
				"```markdown\n## Docs Readme\n```\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir, cleanup := setupTestDir(t)
			defer cleanup()

			for path, content := range tt.setupFiles {
				createTestFile(t, filepath.Join(dir, path), content)
			}

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

			actualOutput := CollectReadmeFiles(dir, resolvedIncludePaths, tt.excludeNames, resolvedExcludePaths, true)
			actualOutput = strings.ReplaceAll(actualOutput, "\n", "\n")
			expectedOutput := strings.ReplaceAll(tt.expectedOutput, "\n", "\n")

			if actualOutput != expectedOutput {
				t.Errorf("CollectReadmeFiles() output mismatch:\nExpected:\n%q\nActual:\n%q", expectedOutput, actualOutput)
			}
		})
	}
}

// TestCollectDependencyFiles は collectDependencyFiles 関数のユニットテストです。
func TestCollectDependencyFiles(t *testing.T) {
	tests := []struct {
		name                   string
		setupFiles             map[string]string
		primaryLangs           []string
		fallbackLangs          map[string]int
		debugMode              bool
		expectedContentsCount  int
		expectedProcessedCount int
		expectedContains       []string
	}{
		{
			name:                   "Debug mode off",
			setupFiles:             map[string]string{"go.mod": "module myapp"},
			primaryLangs:           []string{"Go"},
			debugMode:              false,
			expectedContentsCount:  0,
			expectedProcessedCount: 0,
		},
		{
			name:                   "Debug mode on with primary languages",
			setupFiles:             map[string]string{"go.mod": "module myapp", "go.sum": ""},
			primaryLangs:           []string{"Go"},
			debugMode:              true,
			expectedContentsCount:  2,
			expectedProcessedCount: 2,
			expectedContains:       []string{"### go.mod", "### go.sum"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir, cleanup := setupTestDir(t)
			defer cleanup()

			for path, content := range tt.setupFiles {
				createTestFile(t, filepath.Join(dir, path), content)
			}

			contents, processed := collectDependencyFiles(dir, tt.primaryLangs, tt.fallbackLangs, nil, nil, nil, tt.debugMode)

			if len(contents[DependencyFilesCategory]) != tt.expectedContentsCount {
				t.Errorf("expected %d content entries, got %d", tt.expectedContentsCount, len(contents[DependencyFilesCategory]))
			}
			if len(processed) != tt.expectedProcessedCount {
				t.Errorf("expected %d processed files, got %d", tt.expectedProcessedCount, len(processed))
			}

			actualContent := strings.Join(contents[DependencyFilesCategory], "\n")
			for _, substr := range tt.expectedContains {
				if !strings.Contains(actualContent, substr) {
					t.Errorf("expected content to contain %q, but it did not", substr)
				}
			}
		})
	}
}

// TestCollectSourceFiles は collectSourceFiles 関数のユニットテストです。
func TestCollectSourceFiles(t *testing.T) {
	tests := []struct {
		name              string
		setupFiles        map[string]string
		primaryLangs      []string
		fallbackLangs     map[string]int
		processedDepFiles map[string]struct{}
		expectedCounts    map[string]int
		expectedContains  map[string][]string
	}{
		{
			name:             "Collect based on primary languages",
			setupFiles:       map[string]string{"main.go": "package main", "app.js": ""},
			primaryLangs:     []string{"Go"},
			expectedCounts:   map[string]int{"Go": 1, "Javascript": 0},
			expectedContains: map[string][]string{"Go": {"### main.go"}},
		},
		{
			name:           "Skip test files",
			setupFiles:    map[string]string{"main.go": "", "main_test.go": ""},
			primaryLangs:  []string{"Go"},
			expectedCounts: map[string]int{"Go": 1},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir, cleanup := setupTestDir(t)
			defer cleanup()

			resolvedProcessedDepFiles := make(map[string]struct{})
			for path := range tt.processedDepFiles {
				absPath, _ := filepath.Abs(filepath.Join(dir, path))
				resolvedProcessedDepFiles[absPath] = struct{}{}
			}

			for path, content := range tt.setupFiles {
				createTestFile(t, filepath.Join(dir, path), content)
			}

			contents, _, _ := collectSourceFiles(dir, tt.primaryLangs, tt.fallbackLangs, resolvedProcessedDepFiles, nil, nil, nil, true)

			for lang, expectedCount := range tt.expectedCounts {
				if len(contents[lang]) != expectedCount {
					t.Errorf("expected %d files for language %s, got %d", expectedCount, lang, len(contents[lang]))
				}
			}

			if tt.expectedContains != nil {
				for lang, expectedSubstrings := range tt.expectedContains {
					actualContent := strings.Join(contents[lang], "\n")
					for _, substr := range expectedSubstrings {
						if !strings.Contains(actualContent, substr) {
							t.Errorf("expected content for %s to contain %q, but it did not", lang, substr)
						}
					}
				}
			}
		})
	}
}