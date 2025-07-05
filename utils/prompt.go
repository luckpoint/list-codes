package utils

import (
	"fmt"
)

// PromptTemplates contains predefined prompt templates for LLM analysis (legacy, kept for compatibility)
var PromptTemplates = PromptTemplatesEN

// GetPromptTemplateNames returns a sorted list of available prompt template names
func GetPromptTemplateNames() []string {
	var templates map[string]string
	if IsJapanese() {
		templates = PromptTemplatesJA
	} else {
		templates = PromptTemplatesEN
	}
	
	names := make([]string, 0, len(templates))
	for name := range templates {
		names = append(names, name)
	}
	return names
}

// GetPrompt returns the prompt text based on the given prompt parameter.
// If the prompt parameter matches a predefined template name, it returns that template.
// If the prompt parameter appears to be a file path, it reads the prompt from the file.
// Otherwise, it returns the prompt parameter as-is (custom prompt).
func GetPrompt(promptParam string, debugMode bool) (string, error) {
	return GetPromptI18n(promptParam, debugMode)
}

// FormatWithPrompt formats the output content with the specified prompt.
// If prompt is empty, returns the content as-is.
func FormatWithPrompt(prompt, content string) string {
	if prompt == "" {
		return content
	}
	return fmt.Sprintf("%s\n\n%s", prompt, content)
}