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
			name:         "File Outside Include Scope (Additive Behavior)",
			path:         "main.go",
			isDir:        false,
			includePaths: map[string]struct{}{abs("src"): {}},
			wantSkip:     false,
			reason:       "With additive include behavior, normal files should still be processed even when includes are active.",
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
		// New test cases for additive behavior
		{
			name:         "Additive: Normal File with Include Active",
			path:         "service.go",
			isDir:        false,
			includePaths: map[string]struct{}{abs(".env"): {}},
			wantSkip:     false,
			reason:       "Normal source files should still be processed when --include is active (additive behavior).",
		},
		{
			name:         "Additive: Normal Directory with Include Active",
			path:         "src",
			isDir:        true,
			includePaths: map[string]struct{}{abs(".github"): {}},
			wantSkip:     false,
			reason:       "Normal directories should still be processed when --include is active (additive behavior).",
		},
		{
			name:         "Additive: Dotfile Not Included with Include Active",
			path:         ".dockerignore",
			isDir:        false,
			includePaths: map[string]struct{}{abs(".env"): {}},
			wantSkip:     true,
			reason:       "Dotfiles not explicitly included should still be skipped even with --include active.",
		},
		{
			name:         "Additive: Dot Directory Not Included with Include Active",
			path:         ".vscode",
			isDir:        true,
			includePaths: map[string]struct{}{abs(".github"): {}},
			wantSkip:     true,
			reason:       "Dot directories not explicitly included should still be skipped even with --include active.",
		},
		{
			name:         "Additive: Multiple Includes - First Match",
			path:         ".env",
			isDir:        false,
			includePaths: map[string]struct{}{
				abs(".env"):        {},
				abs(".dockerignore"): {},
			},
			wantSkip: false,
			reason:   "File should not be skipped when it matches any of multiple includes.",
		},
		{
			name:         "Additive: Multiple Includes - Second Match",
			path:         ".dockerignore",
			isDir:        false,
			includePaths: map[string]struct{}{
				abs(".env"):        {},
				abs(".dockerignore"): {},
			},
			wantSkip: false,
			reason:   "File should not be skipped when it matches any of multiple includes.",
		},
		{
			name:         "Additive: Multiple Includes - No Match",
			path:         ".gitignore",
			isDir:        false,
			includePaths: map[string]struct{}{
				abs(".env"):        {},
				abs(".dockerignore"): {},
			},
			wantSkip: true,
			reason:   "Dotfile should be skipped when it doesn't match any includes.",
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

			gotSkip := shouldSkipDir(fullPath, name, tc.isDir, tc.includePaths, tc.excludeNames, tc.excludePaths, nil)

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
			
			output := GenerateDirectoryStructure(tempDir, 10, false, includePaths, excludeNames, excludePaths, tt.includeTests, nil)
			
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

// Note: GitIgnore interaction tests are covered by existing gitignore_test.go
// which already tests the interaction between includes and gitignore patterns

func TestShouldSkipDirWithExcludeAndAdditive(t *testing.T) {
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
			name:         "Additive: Include vs Exclude Names - Exclude Wins",
			path:         "src",
			isDir:        true,
			includePaths: map[string]struct{}{abs("src"): {}},
			excludeNames: map[string]struct{}{"src": {}},
			wantSkip:     true,
			reason:       "Exclude by name should override include (exclude has higher priority).",
		},
		{
			name:         "Additive: Include vs Exclude Paths - Exclude Wins",
			path:         "src",
			isDir:        true,
			includePaths: map[string]struct{}{abs("src"): {}},
			excludePaths: map[string]struct{}{abs("src"): {}},
			wantSkip:     true,
			reason:       "Exclude by path should override include (exclude has higher priority).",
		},
		{
			name:         "Additive: Normal File with Excludes Active",
			path:         "main.go",
			isDir:        false,
			includePaths: map[string]struct{}{abs(".env"): {}},
			excludeNames: map[string]struct{}{"vendor": {}},
			excludePaths: map[string]struct{}{abs("node_modules"): {}},
			wantSkip:     false,
			reason:       "Normal files should not be affected by excludes that don't match them.",
		},
		{
			name:         "Additive: Included File with Different Exclude",
			path:         ".env",
			isDir:        false,
			includePaths: map[string]struct{}{abs(".env"): {}},
			excludeNames: map[string]struct{}{"vendor": {}},
			excludePaths: map[string]struct{}{abs("node_modules"): {}},
			wantSkip:     false,
			reason:       "Included files should not be affected by excludes that don't match them.",
		},
		{
			name:         "Additive: Non-included Dotfile with Excludes Active",
			path:         ".dockerignore",
			isDir:        false,
			includePaths: map[string]struct{}{abs(".env"): {}},
			excludeNames: map[string]struct{}{"vendor": {}},
			wantSkip:     true,
			reason:       "Non-included dotfiles should still be skipped by default dotfile rule.",
		},
		{
			name:         "Additive: Multiple Excludes - Both Name and Path",
			path:         "vendor",
			isDir:        true,
			includePaths: map[string]struct{}{abs("vendor"): {}, abs(".env"): {}},
			excludeNames: map[string]struct{}{"vendor": {}},
			excludePaths: map[string]struct{}{abs("vendor"): {}},
			wantSkip:     true,
			reason:       "Directory should be excluded when matched by both exclude name and path (exclude overrides include).",
		},
		{
			name:         "Additive: Complex Priority Test",
			path:         ".env",
			isDir:        false,
			includePaths: map[string]struct{}{abs(".env"): {}, abs(".github"): {}},
			excludeNames: map[string]struct{}{"temp": {}},
			excludePaths: map[string]struct{}{abs("vendor"): {}},
			wantSkip:     false,
			reason:       "Included dotfile should not be skipped when excludes don't match it.",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Initialize maps if they are nil to avoid nil pointer issues
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

			gotSkip := shouldSkipDir(fullPath, name, tc.isDir, tc.includePaths, tc.excludeNames, tc.excludePaths, nil)

			if gotSkip != tc.wantSkip {
				t.Errorf("shouldSkipDir() for '%s' returned skip=%v, want %v. Reason: %s", tc.path, gotSkip, tc.wantSkip, tc.reason)
			}
		})
	}
}

func TestShouldSkipDirWithAdditiveEdgeCases(t *testing.T) {
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
			name:         "Additive: Parent Directory Traversal - Parent of Nested Include",
			path:         "src",
			isDir:        true,
			includePaths: map[string]struct{}{abs(filepath.Join("src", "api", "handler.go")): {}},
			wantSkip:     false,
			reason:       "Parent directories of included items must be traversed and should not be skipped.",
		},
		{
			name:         "Additive: Parent Directory Traversal - Grandparent",
			path:         ".",
			isDir:        true,
			includePaths: map[string]struct{}{abs(filepath.Join("src", "api", "handler.go")): {}},
			wantSkip:     false,
			reason:       "Grandparent directories of included items must be traversed and should not be skipped.",
		},
		{
			name:         "Additive: Deep Nested Include - Intermediate Directory",
			path:         filepath.Join("src", "api"),
			isDir:        true,
			includePaths: map[string]struct{}{abs(filepath.Join("src", "api", "v1", "users.go")): {}},
			wantSkip:     false,
			reason:       "Intermediate directories in the path to included items should not be skipped.",
		},
		{
			name:         "Additive: Sibling Directory of Included Path",
			path:         filepath.Join("src", "lib"),
			isDir:        true,
			includePaths: map[string]struct{}{abs(filepath.Join("src", "api", "handler.go")): {}},
			wantSkip:     false,
			reason:       "Sibling directories should still be processed with additive behavior (not excluded).",
		},
		{
			name:         "Additive: File in Sibling Directory",
			path:         filepath.Join("src", "lib", "utils.go"),
			isDir:        false,
			includePaths: map[string]struct{}{abs(filepath.Join("src", "api", "handler.go")): {}},
			wantSkip:     false,
			reason:       "Files in sibling directories should still be processed with additive behavior.",
		},
		{
			name:         "Additive: Root Level File with Deep Include",
			path:         "main.go",
			isDir:        false,
			includePaths: map[string]struct{}{abs(filepath.Join("src", "api", "handler.go")): {}},
			wantSkip:     false,
			reason:       "Root level files should still be processed with additive behavior.",
		},
		{
			name:         "Additive: Include Directory Itself",
			path:         "src",
			isDir:        true,
			includePaths: map[string]struct{}{abs("src"): {}},
			wantSkip:     false,
			reason:       "Directory that exactly matches include path should not be skipped.",
		},
		{
			name:         "Additive: Include File Itself",
			path:         filepath.Join("src", "handler.go"),
			isDir:        false,
			includePaths: map[string]struct{}{abs(filepath.Join("src", "handler.go")): {}},
			wantSkip:     false,
			reason:       "File that exactly matches include path should not be skipped.",
		},
		{
			name:         "Additive: Empty Include List",
			path:         "main.go",
			isDir:        false,
			includePaths: map[string]struct{}{},
			wantSkip:     false,
			reason:       "With empty include list, normal files should not be skipped (default behavior).",
		},
		{
			name:         "Additive: Dotfile with Empty Include List",
			path:         ".env",
			isDir:        false,
			includePaths: map[string]struct{}{},
			wantSkip:     true,
			reason:       "With empty include list, dotfiles should still be skipped by default.",
		},
		{
			name:         "Additive: Complex Path Matching - Prefix Match",
			path:         "src-backup",  // Should NOT match "src" 
			isDir:        true,
			includePaths: map[string]struct{}{abs("src"): {}},
			wantSkip:     false,
			reason:       "Directory with similar name but not exact match should still be processed (additive behavior).",
		},
		{
			name:         "Additive: Case Sensitivity",
			path:         "SRC",
			isDir:        true,
			includePaths: map[string]struct{}{abs("src"): {}},
			wantSkip:     false,
			reason:       "Case-different directory should still be processed with additive behavior (paths are case sensitive).",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Initialize maps if they are nil to avoid nil pointer issues
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

			gotSkip := shouldSkipDir(fullPath, name, tc.isDir, tc.includePaths, tc.excludeNames, tc.excludePaths, nil)

			if gotSkip != tc.wantSkip {
				t.Errorf("shouldSkipDir() for '%s' returned skip=%v, want %v. Reason: %s", tc.path, gotSkip, tc.wantSkip, tc.reason)
			}
		})
	}
}

func TestIsAssetFile(t *testing.T) {
	testCases := []struct {
		name    string
		path    string
		isAsset bool
	}{
		{"PNG Image", "images/logo.png", true},
		{"SVG Image", "icons/button.svg", true},
		{"MP4 Video", "media/intro.mp4", true},
		{"WOFF2 Font", "assets/fonts/inter.woff2", true},
		{"PDF Document", "docs/spec.pdf", true},
		{"Go Source File", "main.go", false},
		{"File without extension", "README", false},
		{"Case insensitive extension", "IMAGE.PNG", true},
		{"JavaScript file", "app.js", false},
		{"CSS file", "style.css", false},
		{"HTML file", "index.html", false},
		{"Unknown extension", "file.xyz", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := IsAssetFile(tc.path, false)
			if got != tc.isAsset {
				t.Errorf("IsAssetFile(%q) = %v, want %v", tc.path, got, tc.isAsset)
			}
		})
	}
}

func TestCollectReadmeFiles(t *testing.T) {
	tempDir := t.TempDir()
	
	testCases := []struct {
		name           string
		setupFiles     map[string]string
		expectedCount  int
		expectedText   string
		excludeNames   map[string]struct{}
		excludePaths   map[string]struct{}
	}{
		{
			name: "Single README file",
			setupFiles: map[string]string{
				"README.md": "# Main Project\nThis is the main readme.",
			},
			expectedCount: 1,
			expectedText:  "# Main Project",
		},
		{
			name: "Multiple README files",
			setupFiles: map[string]string{
				"README.md":     "# Root README",
				"docs/README.md": "# Documentation README",
			},
			expectedCount: 2,
			expectedText:  "# Root README",
		},
		{
			name: "No README files",
			setupFiles: map[string]string{
				"main.go": "package main",
			},
			expectedCount: 0,
			expectedText:  "No README.md files found",
		},
		{
			name: "README in excluded directory",
			setupFiles: map[string]string{
				"README.md":           "# Root README",
				"node_modules/README.md": "# Node modules README",
			},
			excludeNames:  map[string]struct{}{"node_modules": {}},
			expectedCount: 1,
			expectedText:  "# Root README",
		},
		{
			name: "Different file extensions",
			setupFiles: map[string]string{
				"readme.md":  "# Correct readme",
				"readme.txt": "# Wrong extension", // This won't match as it's not exactly "readme.md"
			},
			expectedCount: 1, // Only readme.md will match
			expectedText:  "# Correct readme",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create test directory structure
			testDir := filepath.Join(tempDir, tc.name)
			err := os.MkdirAll(testDir, 0755)
			if err != nil {
				t.Fatalf("Failed to create test directory: %v", err)
			}

			for filePath, content := range tc.setupFiles {
				fullPath := filepath.Join(testDir, filePath)
				err := os.MkdirAll(filepath.Dir(fullPath), 0755)
				if err != nil {
					t.Fatalf("Failed to create directory: %v", err)
				}
				err = os.WriteFile(fullPath, []byte(content), 0644)
				if err != nil {
					t.Fatalf("Failed to create file %s: %v", filePath, err)
				}
			}

			// Initialize maps
			includePaths := make(map[string]struct{})
			excludeNames := tc.excludeNames
			if excludeNames == nil {
				excludeNames = make(map[string]struct{})
			}
			excludePaths := tc.excludePaths
			if excludePaths == nil {
				excludePaths = make(map[string]struct{})
			}

			// Collect README files
			output := CollectReadmeFiles(testDir, includePaths, excludeNames, excludePaths, false, nil)

			// Check if expected text is present
			if !strings.Contains(output, tc.expectedText) {
				t.Errorf("Expected output to contain %q, but it did not\nOutput:\n%s", tc.expectedText, output)
			}

			// Count README sections in output
			sectionCount := strings.Count(output, "### ")
			if sectionCount != tc.expectedCount {
				t.Errorf("Expected %d README sections, got %d\nOutput:\n%s", tc.expectedCount, sectionCount, output)
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
			
			sourceFileContents, _, _, _ := collectSourceFiles(tempDir, primaryLangs, fallbackLangs, processedDepFiles, includePaths, excludeNames, excludePaths, false, tt.includeTests, nil)
			
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

func TestCollectSourceFilesFileSizeLimits(t *testing.T) {
	// Save original values
	originalMaxFileSize := MaxFileSizeBytes
	originalTotalMaxSize := TotalMaxFileSizeBytes
	defer func() {
		MaxFileSizeBytes = originalMaxFileSize
		TotalMaxFileSizeBytes = originalTotalMaxSize
	}()

	tempDir := t.TempDir()
	
	// Create test files with specific sizes
	testFiles := map[string]string{
		"small.go":  strings.Repeat("// comment\n", 10),  // Small file (~110 bytes)
		"medium.go": strings.Repeat("// comment\n", 40),  // Medium file (~440 bytes)
		"large.go":  strings.Repeat("// comment\n", 1000), // Large file (~11KB)
	}
	
	for filePath, content := range testFiles {
		fullPath := filepath.Join(tempDir, filePath)
		err := os.WriteFile(fullPath, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to create file %s: %v", filePath, err)
		}
	}

	t.Run("MaxFileSizeBytes limit", func(t *testing.T) {
		// Set limit to 500 bytes - should skip large.go
		MaxFileSizeBytes = 500
		TotalMaxFileSizeBytes = 0 // No total limit
		
		primaryLangs := []string{"Go"}
		fallbackLangs := map[string]int{"Go": 3}
		processedDepFiles := make(map[string]struct{})
		includePaths := make(map[string]struct{})
		excludeNames := make(map[string]struct{})
		excludePaths := make(map[string]struct{})
		
		sourceFileContents, _, skippedMessages, limitHit := collectSourceFiles(tempDir, primaryLangs, fallbackLangs, processedDepFiles, includePaths, excludeNames, excludePaths, false, false, nil)
		
		// Should not hit total limit
		if limitHit {
			t.Errorf("Expected limitHit to be false when using MaxFileSizeBytes limit, got true")
		}
		
		// Should have skipped large.go
		if len(skippedMessages) == 0 {
			t.Errorf("Expected skipped messages for large files, got none")
		}
		
		// Combine all collected content
		var allContent string
		for _, files := range sourceFileContents {
			allContent += strings.Join(files, "")
		}
		
		// Should contain small and medium files but not large
		if !strings.Contains(allContent, "small.go") {
			t.Errorf("Expected to find small.go in collected files")
		}
		if !strings.Contains(allContent, "medium.go") {
			t.Errorf("Expected to find medium.go in collected files")
		}
		if strings.Contains(allContent, "large.go") {
			t.Errorf("Expected large.go to be skipped due to size limit")
		}
	})

	t.Run("TotalMaxFileSizeBytes limit", func(t *testing.T) {
		// Set individual file limit high but total limit low
		MaxFileSizeBytes = 100000 // 100KB - high individual limit
		TotalMaxFileSizeBytes = 11200 // 11200 bytes total - allow large.go (11000) but not medium.go (440)
		
		primaryLangs := []string{"Go"}
		fallbackLangs := map[string]int{"Go": 3}
		processedDepFiles := make(map[string]struct{})
		includePaths := make(map[string]struct{})
		excludeNames := make(map[string]struct{})
		excludePaths := make(map[string]struct{})
		
		sourceFileContents, _, skippedMessages, limitHit := collectSourceFiles(tempDir, primaryLangs, fallbackLangs, processedDepFiles, includePaths, excludeNames, excludePaths, false, false, nil)
		
		// Should hit total limit
		if !limitHit {
			t.Errorf("Expected limitHit to be true when exceeding TotalMaxFileSizeBytes, got false")
		}
		
		// Should have processed at least one file before hitting limit
		totalFiles := 0
		for _, files := range sourceFileContents {
			totalFiles += len(files)
		}
		
		if totalFiles == 0 {
			t.Errorf("Expected to process at least one file before hitting total size limit")
		}
		
		// Files that exceed the total should not have skipped messages (they stop processing instead)
		_ = skippedMessages // No specific assertion as behavior may vary
	})
}

func TestGenerateDirectoryStructureMaxDepth(t *testing.T) {
	tempDir := t.TempDir()
	
	// Create a deep directory structure
	testFiles := map[string]string{
		"root.go":                     "package main",
		"level1/file1.go":            "package level1",
		"level1/level2/file2.go":     "package level2", 
		"level1/level2/level3/file3.go": "package level3",
		"level1/level2/level3/level4/file4.go": "package level4",
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

	testCases := []struct {
		name             string
		maxDepth         int
		expectedFiles    []string
		unexpectedFiles  []string
		expectEllipsis   bool
	}{
		{
			name:             "Depth 0 - root only",
			maxDepth:         0,
			expectedFiles:    []string{"root.go"},
			unexpectedFiles:  []string{"level1", "file1.go", "file2.go"},
			expectEllipsis:   true,
		},
		{
			name:             "Depth 1 - one level deep",  
			maxDepth:         1,
			expectedFiles:    []string{"root.go", "level1", "file1.go"},
			unexpectedFiles:  []string{"file2.go", "level2", "file3.go"},
			expectEllipsis:   true,
		},
		{
			name:             "Depth 2 - two levels deep",
			maxDepth:         2,
			expectedFiles:    []string{"root.go", "level1", "file1.go", "level2", "file2.go"},
			unexpectedFiles:  []string{"file3.go", "level3"},
			expectEllipsis:   true,
		},
		{
			name:             "High depth - no limit",
			maxDepth:         10,
			expectedFiles:    []string{"root.go", "file1.go", "file2.go", "file3.go", "file4.go"},
			unexpectedFiles:  []string{},
			expectEllipsis:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			includePaths := make(map[string]struct{})
			excludeNames := make(map[string]struct{})
			excludePaths := make(map[string]struct{})
			
			output := GenerateDirectoryStructure(tempDir, tc.maxDepth, false, includePaths, excludeNames, excludePaths, false, nil)
			
			for _, expectedFile := range tc.expectedFiles {
				if !strings.Contains(output, expectedFile) {
					t.Errorf("Expected output to contain %q with maxDepth=%d, but it did not\nOutput:\n%s", expectedFile, tc.maxDepth, output)
				}
			}
			
			for _, unexpectedFile := range tc.unexpectedFiles {
				if strings.Contains(output, unexpectedFile) {
					t.Errorf("Expected output not to contain %q with maxDepth=%d, but it did\nOutput:\n%s", unexpectedFile, tc.maxDepth, output)
				}
			}
			
			if tc.expectEllipsis {
				if !strings.Contains(output, "...") {
					t.Errorf("Expected output to contain ellipsis (...) with maxDepth=%d, but it did not\nOutput:\n%s", tc.maxDepth, output)
				}
			} else {
				if strings.Contains(output, "...") {
					t.Errorf("Expected output not to contain ellipsis (...) with maxDepth=%d, but it did\nOutput:\n%s", tc.maxDepth, output)
				}
			}
		})
	}
}
