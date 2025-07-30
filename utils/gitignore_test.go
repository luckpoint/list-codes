package utils

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewGitIgnoreMatcher(t *testing.T) {
	// Create a temporary directory
	tempDir := t.TempDir()

	// Test with empty directory (no .gitignore files)
	matcher, err := NewGitIgnoreMatcher(tempDir)
	if err != nil {
		t.Fatalf("NewGitIgnoreMatcher failed: %v", err)
	}
	if matcher == nil {
		t.Fatal("Expected matcher to be non-nil")
	}
	if len(matcher.tree) != 0 {
		t.Errorf("Expected empty tree, got %d entries", len(matcher.tree))
	}

	// Test with a .gitignore file
	gitignoreContent := `*.log
/build/
node_modules/
!important.log
`
	gitignorePath := filepath.Join(tempDir, ".gitignore")
	err = os.WriteFile(gitignorePath, []byte(gitignoreContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create .gitignore file: %v", err)
	}

	matcher, err = NewGitIgnoreMatcher(tempDir)
	if err != nil {
		t.Fatalf("NewGitIgnoreMatcher failed with .gitignore: %v", err)
	}
	if len(matcher.tree) != 1 {
		t.Errorf("Expected 1 entry in tree, got %d", len(matcher.tree))
	}
}

func TestGitIgnoreMatcher_Match(t *testing.T) {
	// Create a temporary directory structure
	tempDir := t.TempDir()
	
	// Create subdirectories
	buildDir := filepath.Join(tempDir, "build")
	srcDir := filepath.Join(tempDir, "src")
	err := os.MkdirAll(buildDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create build directory: %v", err)
	}
	err = os.MkdirAll(srcDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create src directory: %v", err)
	}

	// Create .gitignore file
	gitignoreContent := `*.log
/build/
node_modules/
!important.log
`
	gitignorePath := filepath.Join(tempDir, ".gitignore")
	err = os.WriteFile(gitignorePath, []byte(gitignoreContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create .gitignore file: %v", err)
	}

	// Create test files
	testFiles := []string{
		"test.log",
		"important.log",
		"build/output.txt",
		"src/main.go",
		"node_modules/package.json",
	}
	
	for _, testFile := range testFiles {
		fullPath := filepath.Join(tempDir, testFile)
		dir := filepath.Dir(fullPath)
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			t.Fatalf("Failed to create directory for %s: %v", testFile, err)
		}
		err = os.WriteFile(fullPath, []byte("test content"), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", testFile, err)
		}
	}

	// Create matcher
	matcher, err := NewGitIgnoreMatcher(tempDir)
	if err != nil {
		t.Fatalf("NewGitIgnoreMatcher failed: %v", err)
	}

	// Test cases
	testCases := []struct {
		path     string
		expected bool
		reason   string
	}{
		{filepath.Join(tempDir, "test.log"), true, "*.log should be ignored"},
		{filepath.Join(tempDir, "important.log"), false, "!important.log should not be ignored (negation pattern)"},
		{filepath.Join(tempDir, "build", "output.txt"), true, "/build/ should be ignored"},
		{filepath.Join(tempDir, "src", "main.go"), false, "src/main.go should not be ignored"},
		{filepath.Join(tempDir, "node_modules", "package.json"), true, "node_modules/ should be ignored"},
		{filepath.Join(tempDir, "normal.txt"), false, "normal files should not be ignored"},
	}

	for _, tc := range testCases {
		result := matcher.Match(tc.path)
		if result != tc.expected {
			t.Errorf("Match(%s) = %v, expected %v. Reason: %s", tc.path, result, tc.expected, tc.reason)
		}
	}
}

func TestGitIgnoreMatcher_Match_NestedGitignore(t *testing.T) {
	// Create a temporary directory structure
	tempDir := t.TempDir()
	
	// Create subdirectories
	subDir := filepath.Join(tempDir, "subdir")
	err := os.MkdirAll(subDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create subdirectory: %v", err)
	}

	// Create root .gitignore
	rootGitignore := `*.tmp
/build/
`
	err = os.WriteFile(filepath.Join(tempDir, ".gitignore"), []byte(rootGitignore), 0644)
	if err != nil {
		t.Fatalf("Failed to create root .gitignore: %v", err)
	}

	// Create subdirectory .gitignore
	subGitignore := `*.local
!important.local
`
	err = os.WriteFile(filepath.Join(subDir, ".gitignore"), []byte(subGitignore), 0644)
	if err != nil {
		t.Fatalf("Failed to create subdirectory .gitignore: %v", err)
	}

	// Create matcher
	matcher, err := NewGitIgnoreMatcher(tempDir)
	if err != nil {
		t.Fatalf("NewGitIgnoreMatcher failed: %v", err)
	}

	// Test cases for nested .gitignore behavior
	testCases := []struct {
		path     string
		expected bool
		reason   string
	}{
		{filepath.Join(tempDir, "test.tmp"), true, "*.tmp from root should be ignored"},
		{filepath.Join(subDir, "test.tmp"), true, "*.tmp from root should apply to subdirectory"},
		{filepath.Join(subDir, "test.local"), true, "*.local from subdirectory should be ignored"},
		{filepath.Join(subDir, "important.local"), false, "!important.local should not be ignored"},
		{filepath.Join(tempDir, "test.local"), false, "*.local should not apply to root directory"},
	}

	for _, tc := range testCases {
		result := matcher.Match(tc.path)
		if result != tc.expected {
			t.Errorf("Match(%s) = %v, expected %v. Reason: %s", tc.path, result, tc.expected, tc.reason)
		}
	}
}

func TestGitIgnoreMatcher_NilMatcher(t *testing.T) {
	var matcher *GitIgnoreMatcher
	
	// Nil matcher should not match anything
	result := matcher.Match("/some/path")
	if result != false {
		t.Errorf("Nil matcher should return false, got %v", result)
	}
}

func TestGitIgnoreMatcher_Match_WithIncludePaths(t *testing.T) {
	// Create a temporary directory
	tempDir := t.TempDir()
	
	// Create .gitignore file
	gitignoreContent := `*.log
build/
`
	gitignorePath := filepath.Join(tempDir, ".gitignore")
	err := os.WriteFile(gitignorePath, []byte(gitignoreContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create .gitignore file: %v", err)
	}

	// Create build directory and files
	buildDir := filepath.Join(tempDir, "build")
	err = os.MkdirAll(buildDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create build directory: %v", err)
	}

	buildFile := filepath.Join(buildDir, "output.txt")
	logFile := filepath.Join(tempDir, "test.log")
	normalFile := filepath.Join(tempDir, "main.go")

	// Create test files
	for _, file := range []string{buildFile, logFile, normalFile} {
		err = os.WriteFile(file, []byte("test"), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", file, err)
		}
	}

	// Create matcher
	matcher, err := NewGitIgnoreMatcher(tempDir)
	if err != nil {
		t.Fatalf("NewGitIgnoreMatcher failed: %v", err)
	}

	// Test shouldSkipDir with --include paths that override .gitignore
	includePaths := map[string]struct{}{
		buildFile: {},
	}
	excludeNames := make(map[string]struct{})
	excludePaths := make(map[string]struct{})

	// Test that --include overrides .gitignore rules (Priority 2 > Priority 3)
	shouldSkip := shouldSkipDir(buildFile, "output.txt", false, includePaths, excludeNames, excludePaths, matcher)
	if shouldSkip {
		t.Errorf("File explicitly included via --include should not be skipped even if it matches .gitignore")
	}

	// Test that files not in --include are still subject to .gitignore rules
	shouldSkip = shouldSkipDir(logFile, "test.log", false, includePaths, excludeNames, excludePaths, matcher)
	if !shouldSkip {
		t.Errorf("File matching .gitignore and not in --include should be skipped")
	}

	// Test normal file not in .gitignore and not in --include (should NOT be skipped due to additive behavior)
	shouldSkip = shouldSkipDir(normalFile, "main.go", false, includePaths, excludeNames, excludePaths, matcher)
	if shouldSkip {
		t.Errorf("Normal file should NOT be skipped when --include is active (additive behavior)")
	}
}

func TestGitIgnoreMatcher_ComplexPatterns(t *testing.T) {
	// Create a temporary directory structure for complex patterns
	tempDir := t.TempDir()
	
	// Create complex directory structure
	dirs := []string{
		"logs",
		"a/some/deep/b",
		"temp",
		"logs/important",
		"other/a/b",
	}
	
	for _, dir := range dirs {
		err := os.MkdirAll(filepath.Join(tempDir, dir), 0755)
		if err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}
	}
	
	// Create test files
	testFiles := []string{
		"logs/debug.log",
		"logs/important.log",
		"logs/important/critical.log",
		"a/some/deep/b/file.txt",
		"other/a/b/file.txt",
		"temp/data.tmp",
		"temp/cache/file.cache",
		"important.log",
		"script.sh",
	}
	
	for _, testFile := range testFiles {
		fullPath := filepath.Join(tempDir, testFile)
		dir := filepath.Dir(fullPath)
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			t.Fatalf("Failed to create directory for %s: %v", testFile, err)
		}
		err = os.WriteFile(fullPath, []byte("test content"), 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", testFile, err)
		}
	}

	// Create .gitignore with complex patterns
	gitignoreContent := `# Complex patterns test
logs/
!logs/important.log
a/**/b
temp/
!temp/cache/
*.tmp
!important.log
`
	gitignorePath := filepath.Join(tempDir, ".gitignore")
	err := os.WriteFile(gitignorePath, []byte(gitignoreContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create .gitignore file: %v", err)
	}

	// Create matcher
	matcher, err := NewGitIgnoreMatcher(tempDir)
	if err != nil {
		t.Fatalf("NewGitIgnoreMatcher failed: %v", err)
	}

	// Test cases for complex patterns
	testCases := []struct {
		path     string
		expected bool
		reason   string
	}{
		// Directory patterns with trailing slash
		{filepath.Join(tempDir, "logs", "debug.log"), true, "logs/ should ignore everything in logs directory"},
		{filepath.Join(tempDir, "logs", "important.log"), false, "!logs/important.log should negate the logs/ rule"},
		{filepath.Join(tempDir, "logs", "important", "critical.log"), true, "subdirectories of logs/ should still be ignored"},
		
		// Nested directory wildcards with **
		{filepath.Join(tempDir, "a", "some", "deep", "b", "file.txt"), true, "a/**/b should match nested b directories"},
		{filepath.Join(tempDir, "other", "a", "b", "file.txt"), true, "a/**/b should match a/b pattern anywhere"},
		
		// Directory-only patterns with negations
		{filepath.Join(tempDir, "temp", "data.tmp"), true, "temp/ should ignore directory"},
		{filepath.Join(tempDir, "temp", "cache", "file.cache"), false, "!temp/cache/ should negate temp/ for cache subdirectory"},
		
		// File extension patterns with negations
		{filepath.Join(tempDir, "important.log"), false, "!important.log should negate *.tmp (even though it doesn't match *.tmp, the negation for important.log applies)"},
		{filepath.Join(tempDir, "script.sh"), false, "normal files should not be ignored"},
	}

	for _, tc := range testCases {
		result := matcher.Match(tc.path)
		if result != tc.expected {
			t.Errorf("Match(%s) = %v, expected %v. Reason: %s", tc.path, result, tc.expected, tc.reason)
		}
	}
}