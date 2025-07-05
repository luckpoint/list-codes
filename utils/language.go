package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
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

// DetectProjectLanguages detects the programming languages used in the project.
func DetectProjectLanguages(folderPath string, debugMode bool, includePaths, excludeNames, excludePaths map[string]struct{}) ([]string, map[string]int) {
	detectedLanguages := make(map[string]struct{})
	signatureMap := make(map[string]string)
	for lang, sigs := range PROJECT_SIGNATURES {
		for _, sig := range sigs {
			signatureMap[sig] = lang
		}
	}

	filepath.WalkDir(folderPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			PrintWarning(fmt.Sprintf("Error accessing path %s: %v", path, err), debugMode)
			return nil
		}

		absPath, absErr := filepath.Abs(path)
		if absErr != nil {
			PrintWarning(fmt.Sprintf("Could not get absolute path for %s: %v", path, absErr), debugMode)
			return nil
		}

		if d.IsDir() && shouldSkipDir(absPath, d.Name(), true, includePaths, excludeNames, excludePaths) {
			return filepath.SkipDir
		}

		if lang, ok := signatureMap[d.Name()]; ok {
			detectedLanguages[lang] = struct{}{}
			PrintDebug(fmt.Sprintf("Detected '%s' via file/dir '%s' in '%s'.", lang, d.Name(), filepath.Dir(path)), debugMode)
		}

		ext := strings.ToLower(filepath.Ext(d.Name()))
		if lang, ok := signatureMap[ext]; ok {
			detectedLanguages[lang] = struct{}{}
			PrintDebug(fmt.Sprintf("Detected '%s' via extension '%s' in '%s'.", lang, ext, filepath.Dir(path)), debugMode)
		}

		return nil
	})

	primaryDetectedLangs := make([]string, 0, len(detectedLanguages))
	for lang := range detectedLanguages {
		primaryDetectedLangs = append(primaryDetectedLangs, lang)
	}
	sort.Strings(primaryDetectedLangs)

	if len(primaryDetectedLangs) > 0 {
		return primaryDetectedLangs, nil
	}

	PrintDebug("No strong signatures found, falling back to extension counting.", debugMode)
	extensionCounts := make(map[string]int)
	filepath.WalkDir(folderPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			PrintWarning(fmt.Sprintf("Error accessing path %s: %v", path, err), debugMode)
			return nil
		}

		absPath, absErr := filepath.Abs(path)
		if absErr != nil {
			PrintWarning(fmt.Sprintf("Could not get absolute path for %s: %v", path, absErr), debugMode)
			return nil
		}

		if d.IsDir() {
			if shouldSkipDir(absPath, d.Name(), true, includePaths, excludeNames, excludePaths) {
				return filepath.SkipDir
			}
			return nil
		}

		if !d.IsDir() {
			if shouldSkipDir(absPath, d.Name(), false, includePaths, excludeNames, excludePaths) {
				return nil
			}
			lang := GetLanguageByExtension(d.Name())
			if lang != "" && !IsTestFile(path) {
				extensionCounts[lang]++
			}
		}
		return nil
	})

	if len(extensionCounts) == 0 {
		PrintWarning("No source files detected for fallback.", debugMode)
	}

	return nil, extensionCounts
}
