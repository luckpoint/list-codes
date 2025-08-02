package utils

import (
	"path/filepath"
	"strings"
)

// GetLanguageByExtension detects the language from a file name or extension.
func GetLanguageByExtension(fileNameOrExtension string) string {
	if filepath.Base(fileNameOrExtension) == "Dockerfile" {
		return "Dockerfile"
	}
	ext := strings.ToLower(filepath.Ext(fileNameOrExtension))
	if ext == "" {
		ext = strings.ToLower(fileNameOrExtension)
	}

	for lang, extensions := range EXTENSIONS {
		for _, e := range extensions {
			if ext == e {
				return lang
			}
		}
	}
	return ""
}

