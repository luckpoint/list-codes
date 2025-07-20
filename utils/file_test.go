package utils

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestShouldSkipDir(t *testing.T) {
	projectRoot, cleanup := setupTestDir(t)
	defer cleanup()

	// Helper to create absolute paths for the test cases
	abs := func(p string) string {
		path, err := filepath.Abs(filepath.Join(projectRoot, p))
		if err != nil {
			t.Fatalf("Failed to get absolute path for %s: %v", p, err)
		}
		return path
	}

	testCases := []struct {
		name         string
		path         string // relative to projectRoot
		isDir        bool
		includePaths map[string]struct{}
		excludeNames map[string]struct{}
		excludePaths map[string]struct{}
		wantSkip     bool
		reason       string
	}{
		{
			name:     "Standard File",
			path:     "main.go",
			isDir:    false,
			wantSkip: false,
			reason:   "Standard files should not be skipped by default.",
		},
		{
			name:     "Standard Directory",
			path:     "src",
			isDir:    true,
			wantSkip: false,
			reason:   "Standard directories should not be skipped by default.",
		},
		{
			name:     "Dotfile Default Skip",
			path:     ".env",
			isDir:    false,
			wantSkip: true,
			reason:   "Dotfiles should be skipped by default.",
		},
		{
			name:     "Dot Directory Default Skip",
			path:     ".git",
			isDir:    true,
			wantSkip: true,
			reason:   "Dot directories should be skipped by default.",
		},
		{
			name:     "File in Dot Directory Skip",
			path:     filepath.Join(".git", "config"),
			isDir:    false,
			wantSkip: false, // The new logic only checks immediate names, not paths
			reason:   "Files inside dot-directories are handled by directory traversal logic, not shouldSkipDir.",
		},
		{
			name:         "Explicit Include of Dotfile",
			path:         ".env",
			isDir:        false,
			includePaths: map[string]struct{}{abs(".env"): {}},
			wantSkip:     false,
			reason:       "An explicitly included dotfile should not be skipped.",
		},
		{
			name:         "Explicit Include of Dot Directory",
			path:         ".git",
			isDir:        true,
			includePaths: map[string]struct{}{abs(".git"): {}},
			wantSkip:     false,
			reason:       "An explicitly included dot-directory should not be skipped.",
		},
		{
			name:         "Parent of Included Item",
			path:         ".", // projectRoot itself
			isDir:        true,
			includePaths: map[string]struct{}{abs(filepath.Join("src", "api.go")): {}},
			wantSkip:     false,
			reason:       "A parent of an included item must be traversed, so it should not be skipped.",
		},
		{
			name:         "File Outside Include Scope",
			path:         "main.go",
			isDir:        false,
			includePaths: map[string]struct{}{abs("src"): {}},
			wantSkip:     true,
			reason:       "When includes are active, items not in the include set (or parents of included items) should be skipped.",
		},
		{
			name:         "Exclude by Name",
			path:         "src",
			isDir:        true,
			excludeNames: map[string]struct{}{"src": {}},
			wantSkip:     true,
			reason:       "Items in excludeNames should be skipped.",
		},
		{
			name:         "Exclude by Path",
			path:         "src",
			isDir:        true,
			excludePaths: map[string]struct{}{abs("src"): {}},
			wantSkip:     true,
			reason:       "Items in excludePaths should be skipped.",
		},
		{
			name:         "Exclude Overrides Include",
			path:         "src",
			isDir:        true,
			includePaths: map[string]struct{}{abs("src"): {}},
			excludeNames: map[string]struct{}{"src": {}},
			wantSkip:     true,
			reason:       "Exclusion rules should have higher precedence than inclusion rules.",
		},
		{
			name:     "Nested Dot Directory Skip",
			path:     filepath.Join("docs", ".claude"),
			isDir:    true,
			wantSkip: true,
			reason:   "Nested dot directories should be skipped by default.",
		},
		{
			name:         "Include Nested Dot Directory",
			path:         filepath.Join("docs", ".claude"),
			isDir:        true,
			includePaths: map[string]struct{}{abs(filepath.Join("docs", ".claude")): {}},
			wantSkip:     false,
			reason:       "An explicitly included nested dot-directory should not be skipped.",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Initialize maps if they are nil to avoid nil pointer issues in the function
			if tc.includePaths == nil {
				tc.includePaths = make(map[string]struct{})
			}
			if tc.excludeNames == nil {
				tc.excludeNames = make(map[string]struct{})
			}
			if tc.excludePaths == nil {
				tc.excludePaths = make(map[string]struct{})
			}

			fullPath := filepath.Join(projectRoot, tc.path)
			name := filepath.Base(fullPath)
			if tc.path == "." {
				name = "."
			}

			gotSkip := shouldSkipDir(fullPath, name, tc.isDir, tc.includePaths, tc.excludeNames, tc.excludePaths)

			if gotSkip != tc.wantSkip {
				t.Errorf("shouldSkipDir() for '%s' returned skip=%v, want %v. Reason: %s", tc.path, gotSkip, tc.wantSkip, tc.reason)
			}
		})
	}
}

func TestIsTestFile(t *testing.T) {
	testCases := []struct {
		name     string
		path     string
		isTest   bool
		isGoTest bool // Specific check for Go's _test.go convention
	}{
		{"Go Test File", "utils/file_test.go", true, true},
		{"Go Test File in subdir", "auth/jwt_test.go", true, true},
		{"JS Test File", "src/components/button.test.js", true, false},
		{"JS Spec File", "src/services/api.spec.ts", true, false},
		{"Python Test File", "tests/test_api.py", true, false},
		{"Python Pytest File", "app/test_models.py", true, false},
		{"Ruby Spec File", "spec/models/user_spec.rb", true, false},
		{"File in test dir", "test/data.json", true, false},
		{"File in tests dir", "/home/user/project/tests/utils/helpers.py", true, false},
		{"Solidity Test File", "contracts/MyContract.t.sol", true, false},
		{"Non-test file", "src/main.go", false, false},
		{"File with 'test' in name", "testing/service.go", false, false},
		{"File with 'spec' in name", "specification/api.md", false, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test generic IsTestFile function
			got := IsTestFile(tc.path, false)
			if got != tc.isTest {
				t.Errorf("IsTestFile(%q) = %v, want %v", tc.path, got, tc.isTest)
			}
		})
	}
}

func TestGenerateDirectoryStructureWithIncludeTests(t *testing.T) {
	tempDir := t.TempDir()
	
	// Create test files structure
	testFiles := map[string]string{
		"main.go":           "package main",
		"service.go":        "package main",
		"main_test.go":      "package main",
		"service_test.go":   "package main", 
		"utils/helper.go":   "package utils",
		"utils/helper_test.go": "package utils",
		"test/data.json":    "{}",
		"tests/setup.py":    "# setup",
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
		name         string
		includeTests bool
		expectFiles  []string
		rejectFiles  []string
	}{
		{
			name:         "exclude tests (default behavior)",
			includeTests: false,
			expectFiles:  []string{"main.go", "service.go", "helper.go"},
			rejectFiles:  []string{"main_test.go", "service_test.go", "helper_test.go", "data.json", "setup.py"},
		},
		{
			name:         "include tests",
			includeTests: true,
			expectFiles:  []string{"main.go", "service.go", "main_test.go", "service_test.go", "helper.go", "helper_test.go", "data.json", "setup.py"},
			rejectFiles:  []string{},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			includePaths := make(map[string]struct{})
			excludeNames := make(map[string]struct{})
			excludePaths := make(map[string]struct{})
			
			output := GenerateDirectoryStructure(tempDir, 10, false, includePaths, excludeNames, excludePaths, tt.includeTests)
			
			for _, expectedFile := range tt.expectFiles {
				if !strings.Contains(output, expectedFile) {
					t.Errorf("Expected output to contain %q when includeTests=%v, but it did not\nOutput:\n%s", expectedFile, tt.includeTests, output)
				}
			}
			
			for _, rejectedFile := range tt.rejectFiles {
				if strings.Contains(output, rejectedFile) {
					t.Errorf("Expected output not to contain %q when includeTests=%v, but it did\nOutput:\n%s", rejectedFile, tt.includeTests, output)
				}
			}
		})
	}
}

func TestCollectSourceFilesWithIncludeTests(t *testing.T) {
	tempDir := t.TempDir()
	
	// Create test files structure
	testFiles := map[string]string{
		"main.go":           "package main\n\nfunc main() {}",
		"service.go":        "package main\n\nfunc Service() {}",
		"main_test.go":      "package main\n\nimport \"testing\"\n\nfunc TestMain(t *testing.T) {}",
		"service_test.go":   "package main\n\nimport \"testing\"\n\nfunc TestService(t *testing.T) {}", 
		"utils/helper.go":   "package utils\n\nfunc Helper() {}",
		"utils/helper_test.go": "package utils\n\nimport \"testing\"\n\nfunc TestHelper(t *testing.T) {}",
		"test/data.json":    "{}",
		"tests/test_setup.py": "# test setup",
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
		name         string
		includeTests bool
		expectFiles  []string  // Files that should be in source collection
		rejectFiles  []string  // Files that should NOT be in source collection
	}{
		{
			name:         "exclude tests (default behavior)",
			includeTests: false,
			expectFiles:  []string{"main.go", "service.go", "helper.go"},
			rejectFiles:  []string{"main_test.go", "service_test.go", "helper_test.go", "test_setup.py"},
		},
		{
			name:         "include tests",
			includeTests: true,
			expectFiles:  []string{"main.go", "service.go", "main_test.go", "service_test.go", "helper.go", "helper_test.go", "test_setup.py"},
			rejectFiles:  []string{}, // JSON files don't have language extension, so won't be collected anyway
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			primaryLangs := []string{"Go", "Python"}
			fallbackLangs := map[string]int{"Go": 10, "Python": 5}
			processedDepFiles := make(map[string]struct{})
			includePaths := make(map[string]struct{})
			excludeNames := make(map[string]struct{})
			excludePaths := make(map[string]struct{})
			
			sourceFileContents, _, _, _ := collectSourceFiles(tempDir, primaryLangs, fallbackLangs, processedDepFiles, includePaths, excludeNames, excludePaths, false, tt.includeTests)
			
			// Convert collected files to a flat string for easier testing
			var collectedFiles []string
			for _, files := range sourceFileContents {
				collectedFiles = append(collectedFiles, files...)
			}
			allOutput := strings.Join(collectedFiles, "\n")
			
			for _, expectedFile := range tt.expectFiles {
				if !strings.Contains(allOutput, expectedFile) {
					t.Errorf("Expected source files to contain %q when includeTests=%v, but it did not\nCollected files:\n%s", expectedFile, tt.includeTests, allOutput)
				}
			}
			
			for _, rejectedFile := range tt.rejectFiles {
				if strings.Contains(allOutput, rejectedFile) {
					t.Errorf("Expected source files not to contain %q when includeTests=%v, but it did\nCollected files:\n%s", rejectedFile, tt.includeTests, allOutput)
				}
			}
		})
	}
}
