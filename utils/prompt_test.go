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
		{"Invalid template", "nonexistent", true}, // Should return error for unknown template
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

func TestGetPrompt_FromFile(t *testing.T) {
	// Create a temporary file with prompt content
	tempDir := t.TempDir()
	promptFile := filepath.Join(tempDir, "test-prompt.txt")
	
	// Test reading from file path - should return error as it's not a predefined template
	_, err := GetPrompt(promptFile, false)
	if err == nil {
		t.Errorf("GetPrompt() should return error for file paths, got no error")
	}
}

func TestGetPrompt_FileNotFound(t *testing.T) {
	nonExistentFile := "/path/that/does/not/exist/prompt.txt"
	
	// Should return error for file paths that aren't predefined templates
	_, err := GetPrompt(nonExistentFile, false)
	if err == nil {
		t.Errorf("GetPrompt() should return error for file paths, got no error")
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
	
	// Should return error for custom text that's not a predefined template
	_, err := GetPrompt(customText, false)
	if err == nil {
		t.Errorf("GetPrompt() should return error for custom text, got no error")
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
	// Verify that all expected prompt templates exist and have content
	expectedTemplates := []string{
		"explain", "find-bugs", "refactor", "security", "optimize",
		"test", "document", "migrate", "scale", "maintain",
		"api-design", "patterns", "review", "architecture", "deploy",
	}
	
	for _, template := range expectedTemplates {
		t.Run("Template_"+template, func(t *testing.T) {
			content, exists := PromptTemplates[template]
			if !exists {
				t.Errorf("Template '%s' does not exist", template)
				return
			}
			
			if strings.TrimSpace(content) == "" {
				t.Errorf("Template '%s' has empty content", template)
			}
			
			// Basic content validation - should contain some structure
			if !strings.Contains(content, "**") {
				t.Errorf("Template '%s' does not appear to have proper formatting", template)
			}
		})
	}
}

func TestPromptTemplateCount(t *testing.T) {
	expectedCount := 15
	actualCount := len(PromptTemplates)
	
	if actualCount != expectedCount {
		t.Errorf("Expected %d prompt templates, got %d", expectedCount, actualCount)
	}
}