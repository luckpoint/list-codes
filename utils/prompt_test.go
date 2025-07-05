package utils

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
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
	tests := []struct {
		name     string
		template string
		wantErr  bool
	}{
		{"Valid template - explain", "explain", false},
		{"Valid template - find-bugs", "find-bugs", false},
		{"Valid template - refactor", "refactor", false},
		{"Invalid template", "nonexistent", false}, // Should be treated as custom text
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prompt, err := GetPrompt(tt.template, false)
			
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPrompt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			
			if tt.template == "explain" {
				if !strings.Contains(prompt, "プロジェクトの目的と主要機能") {
					t.Errorf("Expected explain template to contain specific text")
				}
			}
			
			if tt.template == "nonexistent" {
				if prompt != "nonexistent" {
					t.Errorf("Expected custom text to be returned as-is, got %s", prompt)
				}
			}
		})
	}
}

func TestGetPrompt_FromFile(t *testing.T) {
	// Create a temporary file with prompt content
	tempDir := t.TempDir()
	promptFile := filepath.Join(tempDir, "test-prompt.txt")
	testContent := "This is a test prompt from file\nWith multiple lines"
	
	err := os.WriteFile(promptFile, []byte(testContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	
	// Test reading from file
	prompt, err := GetPrompt(promptFile, false)
	if err != nil {
		t.Errorf("GetPrompt() error = %v", err)
		return
	}
	
	expectedContent := strings.TrimSpace(testContent)
	if prompt != expectedContent {
		t.Errorf("Expected prompt content '%s', got '%s'", expectedContent, prompt)
	}
}

func TestGetPrompt_FileNotFound(t *testing.T) {
	nonExistentFile := "/path/that/does/not/exist/prompt.txt"
	
	// Since the file doesn't exist, it should be treated as custom text
	prompt, err := GetPrompt(nonExistentFile, false)
	if err != nil {
		t.Errorf("GetPrompt() should not error for non-existent files (treated as custom text), got error: %v", err)
		return
	}
	
	if prompt != nonExistentFile {
		t.Errorf("Expected custom text to be returned as-is, got %s", prompt)
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
	
	prompt, err := GetPrompt(customText, false)
	if err != nil {
		t.Errorf("GetPrompt() error = %v", err)
		return
	}
	
	if prompt != customText {
		t.Errorf("Expected custom text to be returned as-is, got '%s'", prompt)
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