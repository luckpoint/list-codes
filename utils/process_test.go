package utils

import (
	"path/filepath"
	"strings"
	"testing"
)

// TestBuildMarkdownOutput は buildMarkdownOutput 関数のユニットテストです。
func TestBuildMarkdownOutput(t *testing.T) {
	tests := []struct {
		name                 string
		directoryStructureMD string
		depFileContents      map[string][]string
		sourceFileContents   map[string][]string
		debugMode            bool
		expectedContains     []string
		expectedNotContains  []string
	}{
		{
			name:                 "All sections present in debug mode",
			directoryStructureMD: "## Project Structure\n```...```",
			depFileContents:      map[string][]string{DependencyFilesCategory: {"### go.mod"}},
			sourceFileContents:   map[string][]string{"Go": {"### main.go"}},
			debugMode:            true,
			expectedContains: []string{
				"## Project Structure",
				"## Dependency and Configuration Files",
				"## Source Code Files",
				"### Go",
				"### go.mod",
				"### main.go",
			},
		},
		{
			name:                "No dependency files section when empty",
			directoryStructureMD: "## Project Structure",
			depFileContents:      map[string][]string{},
			sourceFileContents:   map[string][]string{"Go": {"### main.go"}},
			debugMode:            false,
			expectedContains:     []string{"## Project Structure", "## Source Code Files"},
			expectedNotContains:  []string{"## Dependency and Configuration Files"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := buildMarkdownOutput(tt.directoryStructureMD, tt.depFileContents, tt.sourceFileContents, 1024, 0, tt.debugMode)

			for _, substr := range tt.expectedContains {
				if !strings.Contains(output, substr) {
					t.Errorf("expected output to contain %q, but it did not\nOutput:\n%s", substr, output)
				}
			}
			for _, substr := range tt.expectedNotContains {
				if strings.Contains(output, substr) {
					t.Errorf("expected output not to contain %q, but it did\nOutput:\n%s", substr, output)
				}
			}
		})
	}
}

// TestProcessSourceFiles はProcessSourceFiles関数のユニットテストです。
func TestProcessSourceFiles(t *testing.T) {
	tests := []struct {
		name                string
		setupFiles          map[string]string
		maxDepth            int
		debugMode           bool
		includePaths        map[string]struct{}
		excludeNames        map[string]struct{}
		excludePaths        map[string]struct{}
		expectedContains    []string
		expectedNotContains []string
	}{
		{
			name: "Basic source file processing",
			setupFiles: map[string]string{
				"main.go": "package main\nfunc main(){}",
				"app.py":  "print('hello')",
			},
			maxDepth:  0,
			debugMode: false,
			expectedContains: []string{
				"## Project Structure",
				"main.go",
				"app.py",
				"## Source Code Files",
				"### main.go",
				"```go\npackage main\nfunc main(){}\n```",
				"### app.py",
				"```python\nprint('hello')\n```",
			},
			expectedNotContains: []string{
				"## Dependency and Configuration Files",
			},
		},
		{
			name: "Debug mode with dependency files",
			setupFiles: map[string]string{
				"go.mod":       "module myapp",
				"main.go":      "package main",
				"package.json": "{\"name\":\"test\"}",
			},
			maxDepth:  0,
			debugMode: true,
			expectedContains: []string{
				"## Project Structure",
				"## Dependency and Configuration Files",
				"### go.mod",
				"## Source Code Files",
				"### main.go",
			},
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

			actualOutput := ProcessSourceFiles(dir, tt.maxDepth, resolvedIncludePaths, tt.excludeNames, resolvedExcludePaths, tt.debugMode)
			actualOutput = strings.ReplaceAll(actualOutput, "\n", "\n")

			for _, expected := range tt.expectedContains {
				if !strings.Contains(actualOutput, expected) {
					t.Errorf("ProcessSourceFiles() output missing expected string:\nExpected: %q\nActual:\n%q", expected, actualOutput)
				}
			}
			for _, unexpected := range tt.expectedNotContains {
				if strings.Contains(actualOutput, unexpected) {
					t.Errorf("ProcessSourceFiles() output contains unexpected string:\nUnexpected: %q\nActual:\n%q", unexpected, actualOutput)
				}
			}
		})
	}
}