package utils

import (
	"path/filepath"
	"strings"
	"testing"
	
	"golang.org/x/text/language"
)

func TestGetPromptTemplateNames(t *testing.T) {
	names := GetPromptTemplateNames()
	
	if len(names) == 0 {
		t.Fatal("Expected some prompt template names, got none")
	}
	
	// Check that we have expected templates
	expectedTemplates := []string{"explain", "find-bugs", "refactor", "security", "optimize"}
	for _, expected := range expectedTemplates {
		found := false
		for _, name := range names {
			if name == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected template '%s' not found in template names", expected)
		}
	}
}

func TestGetPrompt_PredefinedTemplate(t *testing.T) {
	// Ensure we start with English for consistent testing
	currentLang = language.English
	printer = nil
	
	tests := []struct {
		name     string
		template string
		wantErr  bool
	}{
		{"Valid template - explain", "explain", false},
		{"Valid template - find-bugs", "find-bugs", false},
		{"Valid template - refactor", "refactor", false},
		{"Custom text as prompt", "nonexistent", false}, // Custom text should be allowed
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prompt, err := GetPrompt(tt.template, false)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPrompt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if tt.template == "explain" && err == nil {
				if !strings.Contains(prompt, "Project purpose") && !strings.Contains(prompt, "プロジェクトの目的と主要機能") {
					t.Errorf("Expected explain template to contain specific text, got: %s", prompt)
				}
			}
		})
	}
}

func TestGetPrompt_CustomTextAllowed(t *testing.T) {
	// Test that custom text is now allowed
	tempDir := t.TempDir()
	promptFile := filepath.Join(tempDir, "test-prompt.txt")
	
	result, err := GetPrompt(promptFile, false)
	if err != nil {
		t.Errorf("GetPrompt() should allow custom text, got error: %v", err)
	}
	
	if result != promptFile {
		t.Errorf("GetPrompt() should return custom text as-is, got: %s, want: %s", result, promptFile)
	}
}

func TestGetPrompt_FileNotFound(t *testing.T) {
	nonExistentFile := "/path/that/does/not/exist/prompt.txt"
	
	// Should return the file path as-is (custom text), no error
	result, err := GetPrompt(nonExistentFile, false)
	if err != nil {
		t.Errorf("GetPrompt() should allow custom text, got error: %v", err)
	}
	
	if result != nonExistentFile {
		t.Errorf("GetPrompt() should return custom text as-is, got: %s, want: %s", result, nonExistentFile)
	}
}

func TestGetPrompt_EmptyPrompt(t *testing.T) {
	prompt, err := GetPrompt("", false)
	if err != nil {
		t.Errorf("GetPrompt() error = %v", err)
		return
	}
	
	if prompt != "" {
		t.Errorf("Expected empty prompt to return empty string, got '%s'", prompt)
	}
}

func TestGetPrompt_CustomText(t *testing.T) {
	customText := "This is a custom prompt text"
	
	// Should return custom text as-is (new behavior)
	result, err := GetPrompt(customText, false)
	if err != nil {
		t.Errorf("GetPrompt() should allow custom text, got error: %v", err)
	}
	
	if result != customText {
		t.Errorf("GetPrompt() should return custom text as-is, got: %s, want: %s", result, customText)
	}
}

func TestFormatWithPrompt(t *testing.T) {
	tests := []struct {
		name     string
		prompt   string
		content  string
		expected string
	}{
		{
			name:     "With prompt",
			prompt:   "Test prompt",
			content:  "Test content",
			expected: "Test prompt\n\nTest content",
		},
		{
			name:     "Empty prompt",
			prompt:   "",
			content:  "Test content",
			expected: "Test content",
		},
		{
			name:     "Empty content",
			prompt:   "Test prompt",
			content:  "",
			expected: "Test prompt\n\n",
		},
		{
			name:     "Both empty",
			prompt:   "",
			content:  "",
			expected: "",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatWithPrompt(tt.prompt, tt.content)
			if result != tt.expected {
				t.Errorf("FormatWithPrompt() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestPromptTemplatesExist(t *testing.T) {
	// Verify that all expected English prompt templates exist and have content
	expectedTemplates := []string{
		"explain", "find-bugs", "refactor", "security", "optimize",
		"test", "document", "migrate", "scale", "maintain",
		"api-design", "patterns", "review", "architecture", "deploy",
	}
	
	for _, template := range expectedTemplates {
		t.Run("EnglishTemplate_"+template, func(t *testing.T) {
			content, exists := PromptTemplatesEN[template]
			if !exists {
				t.Errorf("English template '%s' does not exist", template)
				return
			}
			
			if strings.TrimSpace(content) == "" {
				t.Errorf("English template '%s' has empty content", template)
			}
			
			// Basic content validation - should contain some structure
			if !strings.Contains(content, "**") {
				t.Errorf("English template '%s' does not appear to have proper formatting", template)
			}
		})
	}
}

func TestGetPrompt_Japanese(t *testing.T) {
	defer resetI18nPrompt() // Ensure i18n is reset after the test

	SetLanguage("ja", false)
	
	prompt, err := GetPrompt("explain", false)
	if err != nil {
		t.Fatalf("GetPrompt() returned an unexpected error: %v", err)
	}

	if !strings.Contains(prompt, "プロジェクトの目的と主要機能") {
		t.Errorf("Expected Japanese 'explain' prompt, but got: %s", prompt)
	}
	
	// Test another template
	prompt, err = GetPrompt("find-bugs", false)
	if err != nil {
		t.Fatalf("GetPrompt() returned an unexpected error for find-bugs: %v", err)
	}
	
	if !strings.Contains(prompt, "バグ") && !strings.Contains(prompt, "問題") {
		t.Errorf("Expected Japanese 'find-bugs' prompt with Japanese content, but got: %s", prompt)
	}
}

// resetI18nPrompt helper function for testing prompts
func resetI18nPrompt() {
	SetLanguage("en", false) // Reset to English
}

func TestPromptTemplatesExistJapanese(t *testing.T) {
	// Verify that Japanese templates exist and have content
	expectedTemplates := []string{
		"explain", "find-bugs", "refactor", "security", "optimize",
		"test", "document", "migrate", "scale", "maintain",
		"api-design", "patterns", "review", "architecture", "deploy",
	}
	
	for _, template := range expectedTemplates {
		t.Run("JapaneseTemplate_"+template, func(t *testing.T) {
			content, exists := PromptTemplatesJA[template]
			if !exists {
				t.Errorf("Japanese template '%s' does not exist", template)
				return
			}
			
			if strings.TrimSpace(content) == "" {
				t.Errorf("Japanese template '%s' has empty content", template)
			}
			
			// Basic content validation - should contain some Japanese characters
			hasJapanese := false
			for _, r := range content {
				if (r >= '\u3040' && r <= '\u309F') || // Hiragana
					(r >= '\u30A0' && r <= '\u30FF') || // Katakana
					(r >= '\u4E00' && r <= '\u9FAF') {   // Kanji
					hasJapanese = true
					break
				}
			}
			if !hasJapanese {
				t.Errorf("Japanese template '%s' does not appear to contain Japanese characters", template)
			}
		})
	}
}

func TestPromptTemplateCount(t *testing.T) {
	expectedCount := 15
	actualCountEN := len(PromptTemplatesEN)
	actualCountJA := len(PromptTemplatesJA)
	
	if actualCountEN != expectedCount {
		t.Errorf("Expected %d English prompt templates, got %d", expectedCount, actualCountEN)
	}
	
	if actualCountJA != expectedCount {
		t.Errorf("Expected %d Japanese prompt templates, got %d", expectedCount, actualCountJA)
	}
	
	// Verify both template maps have the same keys
	for key := range PromptTemplatesEN {
		if _, exists := PromptTemplatesJA[key]; !exists {
			t.Errorf("Japanese templates is missing key '%s' that exists in English templates", key)
		}
	}
	
	for key := range PromptTemplatesJA {
		if _, exists := PromptTemplatesEN[key]; !exists {
			t.Errorf("English templates is missing key '%s' that exists in Japanese templates", key)
		}
	}
}