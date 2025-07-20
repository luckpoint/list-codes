package utils

import "testing"

// TestConstants はconfig.goの定数のテストです。
func TestConstants(t *testing.T) {
	// 定数の値をテスト
	if MaxStructureDepthDefault != 7 {
		t.Errorf("MaxStructureDepthDefault should be 7, got %d", MaxStructureDepthDefault)
	}
	
	if DependencyFilesCategory != "Dependency Files" {
		t.Errorf("DependencyFilesCategory should be 'Dependency Files', got %q", DependencyFilesCategory)
	}
}

// TestDefaultExcludeNames はDefaultExcludeNamesのテストです。
// Note: Dotfiles are now handled by a general rule in shouldSkipDir
func TestDefaultExcludeNames(t *testing.T) {
	expectedNames := []string{
		"node_modules", "vendor", "target", "build", "dist",
		"__pycache__", "env", "venv",
	}
	
	for _, name := range expectedNames {
		if _, exists := DefaultExcludeNames[name]; !exists {
			t.Errorf("DefaultExcludeNames should contain %q", name)
		}
	}
	
	// 想定されていない項目が含まれていないかチェック
	if len(DefaultExcludeNames) != len(expectedNames) {
		t.Errorf("DefaultExcludeNames should have %d items, got %d", len(expectedNames), len(DefaultExcludeNames))
	}
}

// TestProjectSignatures はPROJECT_SIGNATURESのテストです。
func TestProjectSignatures(t *testing.T) {
	tests := []struct {
		language      string
		expectedSigs  []string
	}{
		{"Go", []string{"go.mod"}},
		{"Ruby", []string{"Gemfile"}},
		{"Javascript", []string{"package.json", "node_modules"}},
		{"Python", []string{"requirements.txt", "setup.py", "pyproject.toml", "Pipfile"}},
		{"Rust", []string{"Cargo.toml"}},
		{"Java", []string{"build.gradle", "settings.gradle", "pom.xml"}},
		{"Markdown", []string{"README.md"}},
		{"Solidity", []string{"hardhat.config.js", "hardhat.config.ts", "truffle-config.js", "foundry.toml"}},
	}
	
	for _, tt := range tests {
		signatures, exists := PROJECT_SIGNATURES[tt.language]
		if !exists {
			t.Errorf("PROJECT_SIGNATURES should contain language %q", tt.language)
			continue
		}
		
		if len(signatures) != len(tt.expectedSigs) {
			t.Errorf("Language %q should have %d signatures, got %d", tt.language, len(tt.expectedSigs), len(signatures))
			continue
		}
		
		sigMap := make(map[string]bool)
		for _, sig := range signatures {
			sigMap[sig] = true
		}
		
		for _, expectedSig := range tt.expectedSigs {
			if !sigMap[expectedSig] {
				t.Errorf("Language %q should have signature %q", tt.language, expectedSig)
			}
		}
	}
}

// TestExtensions はEXTENSIONSのテストです。
func TestExtensions(t *testing.T) {
	tests := []struct {
		language   string
		extensions []string
	}{
		{"Python", []string{".py", ".pyw"}},
		{"Go", []string{".go"}},
		{"Javascript", []string{".js", ".jsx", ".mjs", ".cjs"}},
		{"Typescript", []string{".ts", ".tsx", ".mts", ".cts"}},
		{"Rust", []string{".rs"}},
		{"Java", []string{".java"}},
		{"C++", []string{".cpp", ".hpp", ".cxx", ".hxx", ".cc", ".hh", ".c++", ".h++"}},
		{"Dockerfile", []string{"Dockerfile"}},
		{"Markdown", []string{".md"}},
		{"Solidity", []string{".sol"}},
	}
	
	for _, tt := range tests {
		extensions, exists := EXTENSIONS[tt.language]
		if !exists {
			t.Errorf("EXTENSIONS should contain language %q", tt.language)
			continue
		}
		
		if len(extensions) != len(tt.extensions) {
			t.Errorf("Language %q should have %d extensions, got %d", tt.language, len(tt.extensions), len(extensions))
			continue
		}
		
		extMap := make(map[string]bool)
		for _, ext := range extensions {
			extMap[ext] = true
		}
		
		for _, expectedExt := range tt.extensions {
			if !extMap[expectedExt] {
				t.Errorf("Language %q should have extension %q", tt.language, expectedExt)
			}
		}
	}
}

// TestFrameworkDependencyFiles はFRAMEWORK_DEPENDENCY_FILESのテストです。
func TestFrameworkDependencyFiles(t *testing.T) {
	tests := []struct {
		language string
		files    []string
	}{
		{"Go", []string{"go.mod", "go.sum"}},
		{"Javascript", []string{"package.json", "package-lock.json", "yarn.lock", "pnpm-lock.yaml", "bun.lockb"}},
		{"Python", []string{"requirements.txt", "Pipfile", "Pipfile.lock", "pyproject.toml", "poetry.lock"}},
		{"Ruby", []string{"Gemfile", "Gemfile.lock"}},
		{"Java", []string{"build.gradle", "settings.gradle", "pom.xml"}},
		{"Rust", []string{"Cargo.toml", "Cargo.lock"}},
	}
	
	for _, tt := range tests {
		files, exists := FRAMEWORK_DEPENDENCY_FILES[tt.language]
		if !exists {
			t.Errorf("FRAMEWORK_DEPENDENCY_FILES should contain language %q", tt.language)
			continue
		}
		
		if len(files) != len(tt.files) {
			t.Errorf("Language %q should have %d dependency files, got %d", tt.language, len(tt.files), len(files))
			continue
		}
		
		fileMap := make(map[string]bool)
		for _, file := range files {
			fileMap[file] = true
		}
		
		for _, expectedFile := range tt.files {
			if !fileMap[expectedFile] {
				t.Errorf("Language %q should have dependency file %q", tt.language, expectedFile)
			}
		}
	}
}

// TestExcludeTestKeywords はEXCLUDE_TEST_KEYWORDSのテストです。
func TestExcludeTestKeywords(t *testing.T) {
	expectedKeywords := []string{"test", "spec", "e2e", "benchmark", "bench", "mock", "fixture"}
	
	if len(EXCLUDE_TEST_KEYWORDS) != len(expectedKeywords) {
		t.Errorf("EXCLUDE_TEST_KEYWORDS should have %d keywords, got %d", len(expectedKeywords), len(EXCLUDE_TEST_KEYWORDS))
	}
	
	keywordMap := make(map[string]bool)
	for _, keyword := range EXCLUDE_TEST_KEYWORDS {
		keywordMap[keyword] = true
	}
	
	for _, expectedKeyword := range expectedKeywords {
		if !keywordMap[expectedKeyword] {
			t.Errorf("EXCLUDE_TEST_KEYWORDS should contain %q", expectedKeyword)
		}
	}
}

// TestExcludeTestDirs はEXCLUDE_TEST_DIRSのテストです。
func TestExcludeTestDirs(t *testing.T) {
	expectedDirs := []string{
		"/test/", "/tests/", "/spec/", "/specs/", "/__tests__/", "/__test__/", 
		"/testing/", "/fixtures/", "/mocks/", "/e2e/", "/integration/", "/unit/",
		"test/", "tests/", "spec/", "specs/",
	}
	
	if len(EXCLUDE_TEST_DIRS) != len(expectedDirs) {
		t.Errorf("EXCLUDE_TEST_DIRS should have %d directories, got %d", len(expectedDirs), len(EXCLUDE_TEST_DIRS))
	}
	
	dirMap := make(map[string]bool)
	for _, dir := range EXCLUDE_TEST_DIRS {
		dirMap[dir] = true
	}
	
	for _, expectedDir := range expectedDirs {
		if !dirMap[expectedDir] {
			t.Errorf("EXCLUDE_TEST_DIRS should contain %q", expectedDir)
		}
	}
}

// TestExcludeTestPatterns はEXCLUDE_TEST_PATTERNSのテストです。
func TestExcludeTestPatterns(t *testing.T) {
	expectedPatterns := []string{
		"_test.go", "_spec.rb", ".test.js", ".spec.js", ".test.jsx", ".spec.jsx",
		".test.ts", ".spec.ts", ".test.tsx", ".spec.tsx", ".test.py", ".spec.py",
		"Test.java", "Tests.java", "IT.java", ".test.php", ".spec.php", "_test.rb",
		".test.cs", ".spec.cs", "Test.cs", "Tests.cs", ".test.cpp", ".spec.cpp",
		".test.c", ".spec.c", ".test.rs", ".spec.rs", ".test.kt", ".spec.kt",
		"Test.kt", ".test.swift", ".spec.swift", "Test.swift", ".t.sol", ".test.sol", ".spec.sol",
	}
	
	if len(EXCLUDE_TEST_PATTERNS) != len(expectedPatterns) {
		t.Errorf("EXCLUDE_TEST_PATTERNS should have %d patterns, got %d", len(expectedPatterns), len(EXCLUDE_TEST_PATTERNS))
	}
	
	patternMap := make(map[string]bool)
	for _, pattern := range EXCLUDE_TEST_PATTERNS {
		patternMap[pattern] = true
	}
	
	for _, expectedPattern := range expectedPatterns {
		if !patternMap[expectedPattern] {
			t.Errorf("EXCLUDE_TEST_PATTERNS should contain %q", expectedPattern)
		}
	}
}