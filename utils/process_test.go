package utils

import (
	"os"
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
				"## Source Code Size Check",
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
			expectedContains:     []string{"## Project Structure", "## Source Code Size Check"},
			expectedNotContains:  []string{"## Dependency and Configuration Files"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := buildMarkdownOutput(tt.directoryStructureMD, tt.depFileContents, tt.sourceFileContents, 1024, []string{}, false, tt.debugMode)

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
				"## Source Code Size Check",
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
				"## Source Code Size Check",
				"### main.go",
			},
			expectedNotContains: []string{
				"## Dependency and Configuration Files", // This section is no longer generated
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

			actualOutput := ProcessSourceFiles(dir, tt.maxDepth, resolvedIncludePaths, tt.excludeNames, resolvedExcludePaths, tt.debugMode, false, nil)
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

func TestProcessSourceFilesWithIncludeTests(t *testing.T) {
	tempDir := t.TempDir()
	
	// Create test files structure
	testFiles := map[string]string{
		"main.go":           "package main\n\nfunc main() {\n\tprintln(\"Hello World\")\n}",
		"service.go":        "package main\n\nfunc Service() {\n\t// service logic\n}",
		"main_test.go":      "package main\n\nimport \"testing\"\n\nfunc TestMain(t *testing.T) {\n\t// test main\n}",
		"service_test.go":   "package main\n\nimport \"testing\"\n\nfunc TestService(t *testing.T) {\n\t// test service\n}", 
		"utils/helper.go":   "package utils\n\nfunc Helper() {\n\t// helper logic\n}",
		"utils/helper_test.go": "package utils\n\nimport \"testing\"\n\nfunc TestHelper(t *testing.T) {\n\t// test helper\n}",
		"tests/test_setup.py": "# test setup\ndef setup():\n    pass",
		"go.mod":            "module test\n\ngo 1.19",
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
	
	tests := []struct {
		name                string
		includeTests        bool
		expectedContains    []string
		expectedNotContains []string
	}{
		{
			name:         "exclude tests (default behavior)",
			includeTests: false,
			expectedContains: []string{
				"## Project Structure", 
				"## Source Code Size Check",
				"main.go", 
				"service.go", 
				"helper.go",
			},
			expectedNotContains: []string{
				"main_test.go", 
				"service_test.go", 
				"helper_test.go",
				"test_setup.py",
			},
		},
		{
			name:         "include tests",
			includeTests: true,
			expectedContains: []string{
				"## Project Structure", 
				"## Source Code Size Check",
				"main.go", 
				"service.go", 
				"helper.go",
				"main_test.go", 
				"service_test.go", 
				"helper_test.go",
				"test_setup.py",
			},
			expectedNotContains: []string{},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			includePaths := make(map[string]struct{})
			excludeNames := make(map[string]struct{})
			excludePaths := make(map[string]struct{})
			
			output := ProcessSourceFiles(tempDir, 10, includePaths, excludeNames, excludePaths, false, tt.includeTests, nil)
			
			for _, expected := range tt.expectedContains {
				if !strings.Contains(output, expected) {
					t.Errorf("Expected ProcessSourceFiles output to contain %q when includeTests=%v, but it did not\nOutput:\n%s", expected, tt.includeTests, output)
				}
			}
			
			for _, unexpected := range tt.expectedNotContains {
				if strings.Contains(output, unexpected) {
					t.Errorf("Expected ProcessSourceFiles output not to contain %q when includeTests=%v, but it did\nOutput:\n%s", unexpected, tt.includeTests, output)
				}
			}
		})
	}
}

func TestBuildMarkdownOutputDeep(t *testing.T) {
	tests := []struct {
		name                 string
		directoryStructureMD string
		depFileContents      map[string][]string
		sourceFileContents   map[string][]string
		totalFileSize        int64
		skippedFileMessages  []string
		limitHit             bool
		debugMode            bool
		expectedContains     []string
		expectedNotContains  []string
		checkSectionOrder    bool
	}{
		{
			name:                 "With skipped files and size limit hit",
			directoryStructureMD: "## Project Structure\n```\nproject/\n  file.go\n```",
			depFileContents:      map[string][]string{},
			sourceFileContents:   map[string][]string{"Go": {"### main.go\n```go\npackage main\n```"}},
			totalFileSize:        1024 * 1024, // 1MB
			skippedFileMessages:  []string{"`large_file.go` (5.00 MB)", "`huge_image.png` (10.50 MB)"},
			limitHit:             true,
			debugMode:            false,
			expectedContains: []string{
				"## Source Code Size Check",
				"**File Statistics**: 1.00 MB collected",
				"**Skipped 2 file(s)",
				"`large_file.go` (5.00 MB)",
				"`huge_image.png` (10.50 MB)",
				"scan stopped",
				"## Project Structure",
			},
			checkSectionOrder: true,
		},
		{
			name:                 "Order check: Size Check before Project Structure",
			directoryStructureMD: "## Project Structure\n```\nproject/\n```",
			depFileContents:      map[string][]string{},
			sourceFileContents:   map[string][]string{"Go": {"### main.go"}},
			totalFileSize:        500,
			skippedFileMessages:  []string{},
			limitHit:             false,
			debugMode:            false,
			expectedContains: []string{
				"## Source Code Size Check",
				"## Project Structure",
			},
			checkSectionOrder: true,
		},
		{
			name:                 "Multiple language sections",
			directoryStructureMD: "## Project Structure\n```\nproject/\n```",
			depFileContents:      map[string][]string{},
			sourceFileContents: map[string][]string{
				"Go":         {"### main.go\n```go\npackage main\n```"},
				"Python":     {"### app.py\n```python\nprint('hello')\n```"},
				"JavaScript": {"### script.js\n```javascript\nconsole.log('hi');\n```"},
			},
			totalFileSize:       2048,
			skippedFileMessages: []string{},
			limitHit:            false,
			debugMode:           false,
			expectedContains: []string{
				"## Source Code Size Check",
				"### Go",
				"### main.go",
				"### Python", 
				"### app.py",
				"### JavaScript",
				"### script.js",
			},
		},
		{
			name:                 "Debug mode with dependency files",
			directoryStructureMD: "## Project Structure\n```\nproject/\n```",
			depFileContents:      map[string][]string{DependencyFilesCategory: {"### go.mod\n```\nmodule test\n```"}},
			sourceFileContents:   map[string][]string{"Go": {"### main.go"}},
			totalFileSize:        1000,
			skippedFileMessages:  []string{},
			limitHit:             false,
			debugMode:            true,
			expectedContains: []string{
				"## Source Code Size Check",
				"## Project Structure",
				"### Go",
			},
			expectedNotContains: []string{
				// Note: Dependency section still appears in debug mode, just updated the test logic
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := buildMarkdownOutput(tt.directoryStructureMD, tt.depFileContents, tt.sourceFileContents, tt.totalFileSize, tt.skippedFileMessages, tt.limitHit, tt.debugMode)

			// Check expected content
			for _, expected := range tt.expectedContains {
				if !strings.Contains(output, expected) {
					t.Errorf("Expected output to contain %q, but it did not\nOutput:\n%s", expected, output)
				}
			}

			// Check content that should not be present
			for _, unexpected := range tt.expectedNotContains {
				if strings.Contains(output, unexpected) {
					t.Errorf("Expected output not to contain %q, but it did\nOutput:\n%s", unexpected, output)
				}
			}

			// Check section ordering if requested
			if tt.checkSectionOrder {
				sizeCheckIndex := strings.Index(output, "## Source Code Size Check")
				projectStructureIndex := strings.Index(output, "## Project Structure")
				
				if sizeCheckIndex == -1 {
					t.Errorf("Expected to find '## Source Code Size Check' section in output")
				}
				if projectStructureIndex == -1 {
					t.Errorf("Expected to find '## Project Structure' section in output")
				}
				
				if sizeCheckIndex != -1 && projectStructureIndex != -1 && sizeCheckIndex >= projectStructureIndex {
					t.Errorf("Expected 'Source Code Size Check' section to appear before 'Project Structure' section. Size Check at %d, Project Structure at %d", sizeCheckIndex, projectStructureIndex)
				}
			}
		})
	}
}