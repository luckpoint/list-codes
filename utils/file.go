package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// IsTestFile determines if the specified file path is a test file.
// It now accepts a debug flag to print its reasoning.
func IsTestFile(filePath string, debug bool) bool {
	normalizedPath := filepath.ToSlash(strings.ToLower(filePath))
	fileName := filepath.Base(normalizedPath)

	// Helper function for logging the decision-making process in debug mode
	logDecision := func(decision bool, reason string) bool {
		if debug && decision {
			PrintDebug(fmt.Sprintf("IsTestFile: Path '%s' identified as test file. Reason: %s", filePath, reason), true)
		}
		return decision
	}

	// 1. Check if the file is in a test directory
	for _, testDirPattern := range EXCLUDE_TEST_DIRS {
		if strings.Contains(normalizedPath, testDirPattern) {
			return logDecision(true, "File is in an excluded directory ('"+testDirPattern+"')")
		}
	}

	// 2. Check against specific filename patterns (most reliable check)
	for _, pattern := range EXCLUDE_TEST_PATTERNS {
		if strings.HasSuffix(fileName, strings.ToLower(pattern)) {
			return logDecision(true, "Filename matches pattern '"+pattern+"'")
		}
	}

	// 3. Check for common keyword-based naming conventions
	baseName := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	for _, keyword := range EXCLUDE_TEST_KEYWORDS {
		// Combined check for prefixes, suffixes, and contains logic
		if strings.HasPrefix(baseName, keyword) || strings.HasSuffix(baseName, keyword) ||
			strings.HasPrefix(baseName, keyword+"_") || strings.HasSuffix(baseName, "_"+keyword) ||
			strings.HasPrefix(baseName, keyword+"-") || strings.HasSuffix(baseName, "-"+keyword) ||
			strings.Contains(baseName, "_"+keyword+"_") || strings.Contains(baseName, "-"+keyword+"-") {
			return logDecision(true, "Filename ('"+baseName+"') contains keyword '"+keyword+"'")
		}
	}

	return false
}

// shouldSkipDir determines whether to skip directories or files based on a set of rules.
// The logic prioritizes user intent:
// 1. User-defined exclusions (--exclude) always result in a skip.
// 2. If --include is used, an item is kept if it's explicitly included or is a parent of an included item.
//    This rule overrides the default dotfile exclusion.
// 3. If not explicitly included, any item starting with a '.' is skipped by default.
// 4. If --include is used and an item is not on the include list (or a parent), it is skipped.
func shouldSkipDir(fullPath, name string, isDir bool, includePaths, excludeNames, excludePaths map[string]struct{}) bool {
	absPath, err := filepath.Abs(fullPath)
	if err != nil {
		PrintWarning(fmt.Sprintf("Could not get absolute path for %s: %v", fullPath, err), true)
		return true // Failsafe skip
	}

	// Priority 1: Explicit --exclude options and hardcoded names always cause a skip.
	if _, ok := excludeNames[name]; ok {
		return true
	}
	if _, ok := excludePaths[absPath]; ok {
		return true
	}

	// Priority 2: Whitelisting logic for --include.
	// If an item is explicitly included or is a parent of an included item,
	// it should NOT be skipped, even if it's a dotfile.
	if len(includePaths) > 0 {
		isNeededForInclude := false
		for incPath := range includePaths {
			// A path is needed if it's the included path itself, or a parent of it.
			// e.g., include "a/b", current is "a" -> strings.HasPrefix("a/b", "a") -> true
			// e.g., include "a/b", current is "a/b" -> strings.HasPrefix("a/b", "a/b") -> true
			if strings.HasPrefix(incPath, absPath) || strings.HasPrefix(absPath, incPath) {
				isNeededForInclude = true
				break
			}
		}

		if isNeededForInclude {
			return false
		}
	}

	// Priority 3: Default exclusion for dotfiles and dot-directories.
	// This runs if the item was not whitelisted by the --include logic above.
	if strings.HasPrefix(name, ".") && name != "." && name != ".." {
		return true
	}

	// Priority 4: If --include is active, skip anything not covered by the whitelist.
	if len(includePaths) > 0 {
		return true // Not needed for the include list, so skip.
	}

	// If no rules caused a skip, process the item.
	return false
}

// GenerateDirectoryStructure generates the project directory structure in Markdown format.
func GenerateDirectoryStructure(startPath string, maxDepth int, debugMode bool, includePaths, excludeNames, excludePaths map[string]struct{}, includeTests bool) string {
	PrintDebug("Generating directory structure...", debugMode)
	var structureLines []string
	structureLines = append(structureLines, "## Project Structure", "```text")

	absStartPath, err := filepath.Abs(startPath)
	if err != nil {
		PrintError(fmt.Sprintf("Could not get absolute path for %s: %v", startPath, err))
		return ""
	}
	rootDisplayName := filepath.Base(absStartPath)
	structureLines = append(structureLines, fmt.Sprintf(". (%s)", rootDisplayName))

	var generateTreeRecursive func(currentPath, prefix string, depth int)
	generateTreeRecursive = func(currentPath, prefix string, depth int) {
		entries, err := os.ReadDir(currentPath)
		if err != nil {
			PrintWarning(fmt.Sprintf("Could not list directory '%s': %v", currentPath, err), debugMode)
			return
		}

		sort.Slice(entries, func(i, j int) bool {
			return strings.ToLower(entries[i].Name()) < strings.ToLower(entries[j].Name())
		})

		var filteredEntries []os.DirEntry
		for _, entry := range entries {
			itemPath := filepath.Join(currentPath, entry.Name())
			if shouldSkipDir(itemPath, entry.Name(), entry.IsDir(), includePaths, excludeNames, excludePaths) {
				continue
			}
			// Skip test files from directory structure
			if !entry.IsDir() && !includeTests && IsTestFile(itemPath, debugMode) {
				continue
			}
			filteredEntries = append(filteredEntries, entry)
		}

		for i, entry := range filteredEntries {
			pointer := "├── "
			if i == len(filteredEntries)-1 {
				pointer = "└── "
			}

			structureLines = append(structureLines, prefix+pointer+entry.Name())

			if entry.IsDir() {
				if maxDepth == 0 || depth < maxDepth-1 {
					extension := "│   "
					if pointer == "└── " {
						extension = "    "
					}
					generateTreeRecursive(filepath.Join(currentPath, entry.Name()), prefix+extension, depth+1)
				} else if depth == maxDepth-1 {
					// Check for subdirectories to print the ellipsis
					nextPath := filepath.Join(currentPath, entry.Name())
					subEntries, err := os.ReadDir(nextPath)
					if err == nil {
						hasVisibleSubItems := false
						for _, subEntry := range subEntries {
							subItemPath := filepath.Join(nextPath, subEntry.Name())
							if shouldSkipDir(subItemPath, subEntry.Name(), subEntry.IsDir(), includePaths, excludeNames, excludePaths) {
								continue
							}
							// Skip test files from directory structure
							if !subEntry.IsDir() && !includeTests && IsTestFile(subItemPath, debugMode) {
								continue
							}
							hasVisibleSubItems = true
							break
						}
						if hasVisibleSubItems {
							extension := "│   "
							if pointer == "└── " {
								extension = "    "
							}
							structureLines = append(structureLines, prefix+extension+"└── ...")
						}
					}
				}
			}
		}
	}

	generateTreeRecursive(absStartPath, "", 0)
	structureLines = append(structureLines, "```")
	PrintDebug("Directory structure generation complete.", debugMode)
	return strings.Join(structureLines, "\n") + "\n\n"
}

// CollectReadmeFiles collects README.md files in the project and returns their content in Markdown format.
func CollectReadmeFiles(folderAbs string, includePaths, excludeNames, excludePaths map[string]struct{}, debug bool) string {
	PrintDebug("Searching for README.md files...", debug)
	var readmeFiles []string

	filepath.WalkDir(folderAbs, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			PrintWarning(fmt.Sprintf("Error accessing path %s: %v", path, err), debug)
			return nil
		}

		if d.IsDir() {
			if shouldSkipDir(path, d.Name(), true, includePaths, excludeNames, excludePaths) {
				return filepath.SkipDir
			}
			return nil
		}

		if strings.ToLower(d.Name()) == "readme.md" {
			if shouldSkipDir(path, d.Name(), false, includePaths, excludeNames, excludePaths) {
				return nil
			}

			content, err := os.ReadFile(path)
			if err != nil {
				PrintWarning(fmt.Sprintf("Could not read README file '%s': %v", path, err), debug)
				return nil
			}
			fileDisplayName, relErr := filepath.Rel(folderAbs, path)
			if relErr != nil {
				PrintWarning(fmt.Sprintf("Could not get relative path for %s: %v", path, relErr), debug)
				fileDisplayName = path
			}
			fileDisplayName = filepath.ToSlash(fileDisplayName)
			markdownContent := fmt.Sprintf("### %s\n\n```markdown\n%s\n```\n", fileDisplayName, string(content))
			readmeFiles = append(readmeFiles, markdownContent)
		}
		return nil
	})

	PrintDebug(fmt.Sprintf("Found %d README.md file(s).", len(readmeFiles)), debug)
	if len(readmeFiles) == 0 {
		return "# Project README Files\n\nNo README.md files found in the project."
	}
	return "# Project README Files\n\n" + strings.Join(readmeFiles, "\n")
}

// collectDependencyFiles collects dependency files.
func collectDependencyFiles(folderAbs string, primaryLangs []string, fallbackLangs map[string]int, includePaths, excludeNames, excludePaths map[string]struct{}, debug bool) (map[string][]string, map[string]struct{}) {
	depFileContents := make(map[string][]string)
	processedDepFiles := make(map[string]struct{})

	if !debug {
		return depFileContents, processedDepFiles
	}

	PrintDebug("Collecting dependency and configuration files...", debug)
	effectiveLangsForDeps := make([]string, 0) // Initialize as empty slice
	if primaryLangs != nil {
		effectiveLangsForDeps = primaryLangs // Assign if not nil
	}

	// If no primary languages were detected, use fallback languages for dependency file collection
	if len(effectiveLangsForDeps) == 0 {
		for lang := range fallbackLangs {
			effectiveLangsForDeps = append(effectiveLangsForDeps, lang)
		}
		sort.Strings(effectiveLangsForDeps)
	}

	for _, langKey := range effectiveLangsForDeps {
		if depFiles, ok := FRAMEWORK_DEPENDENCY_FILES[langKey]; ok {
			for _, depFilePattern := range depFiles {
				filepath.WalkDir(folderAbs, func(path string, d os.DirEntry, err error) error {
					if err != nil {
						PrintWarning(fmt.Sprintf("Error accessing path %s: %v", path, err), debug)
						return nil
					}
					if d.IsDir() {
						if shouldSkipDir(path, d.Name(), true, includePaths, excludeNames, excludePaths) {
							return filepath.SkipDir
						}
						return nil
					}

					if d.Name() == depFilePattern {
						absPath, _ := filepath.Abs(path)
						if _, ok := processedDepFiles[absPath]; ok {
							return nil
						}
						processedDepFiles[absPath] = struct{}{}

						content, err := os.ReadFile(path)
						if err != nil {
							PrintWarning(fmt.Sprintf("Could not read dependency file '%s': %v", path, err), debug)
							return nil
						}
						fileDisplayName, relErr := filepath.Rel(folderAbs, path)
						if relErr != nil {
							PrintWarning(fmt.Sprintf("Could not get relative path for %s: %v", path, relErr), debug)
							fileDisplayName = path
						}
						fileDisplayName = filepath.ToSlash(fileDisplayName)
						lang := GetLanguageByExtension(d.Name())
						if lang == "" {
							lang = strings.ToLower(d.Name())
						}
						codeBlockLangHint := strings.ToLower(lang)
						codeBlockLangHint = strings.ReplaceAll(codeBlockLangHint, "/", "")
						codeBlockLangHint = strings.ReplaceAll(codeBlockLangHint, "+", "p")
						markdownContent := fmt.Sprintf("### %s\n\n```%s\n%s\n```\n", fileDisplayName, codeBlockLangHint, string(content))
						depFileContents[DependencyFilesCategory] = append(depFileContents[DependencyFilesCategory], markdownContent)
					}
					return nil
				})
			}
		}
	}
	return depFileContents, processedDepFiles
}

// collectSourceFiles collects source code files.
func collectSourceFiles(folderAbs string, primaryLangs []string, fallbackLangs map[string]int, processedDepFiles, includePaths, excludeNames, excludePaths map[string]struct{}, debug bool, includeTests bool) (map[string][]string, int64, int) {
	sourceFileContents := make(map[string][]string)
	var totalFileSize int64
	var skippedFileCount int
	PrintDebug("Processing source files...", debug)

	filepath.WalkDir(folderAbs, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			PrintWarning(fmt.Sprintf("Error accessing path %s: %v", path, debug), debug)
			return nil
		}

		if d.IsDir() {
			if shouldSkipDir(path, d.Name(), true, includePaths, excludeNames, excludePaths) {
				return filepath.SkipDir
			}
			return nil
		}

		if shouldSkipDir(path, d.Name(), false, includePaths, excludeNames, excludePaths) {
			return nil
		}

		absPath, _ := filepath.Abs(path)
		if _, ok := processedDepFiles[absPath]; ok {
			return nil
		}
		if !includeTests && IsTestFile(path, debug) {
			return nil
		}

		language := GetLanguageByExtension(d.Name())
		if language == "" {
			return nil
		}

		// By default, process all source files even if they belong to languages
		// that are not part of the detected primaryLangs. This change broadens
		// coverage so that files like JavaScript alongside HTML are included
		// without requiring additional configuration. If future behaviour needs
		// to restrict languages, the caller should introduce an explicit filter
		// rather than relying on this internal function.
		_ = primaryLangs // keep referenced to avoid unused parameter compile errors
		_ = fallbackLangs
		processThisFile := true

		if !processThisFile {
			return nil
		}

		fileInfo, err := d.Info()
		if err != nil {
			PrintWarning(fmt.Sprintf("Could not get file info for '%s': %v", path, err), debug)
			return nil
		}

		if fileInfo.Size() > MaxFileSizeBytes {
			PrintDebug(fmt.Sprintf("Skipping file '%s' due to size (%d bytes > %d bytes)", path, fileInfo.Size(), MaxFileSizeBytes), debug)
			skippedFileCount++
			return nil
		}

		totalFileSize += fileInfo.Size()

		content, err := os.ReadFile(path)
		if err != nil {
			PrintWarning(fmt.Sprintf("Could not read source file '%s': %v", path, debug), debug)
			return nil
		}
		fileDisplayName, relErr := filepath.Rel(folderAbs, path)
		if relErr != nil {
			PrintWarning(fmt.Sprintf("Could not get relative path for %s: %v", path, relErr), debug)
			fileDisplayName = path
		}
		fileDisplayName = filepath.ToSlash(fileDisplayName)
		codeBlockLangHint := strings.ToLower(language)
		codeBlockLangHint = strings.ReplaceAll(codeBlockLangHint, "/", "")
		codeBlockLangHint = strings.ReplaceAll(codeBlockLangHint, "+", "p")
		markdownContent := fmt.Sprintf("### %s\n\n```%s\n%s\n```\n", fileDisplayName, codeBlockLangHint, string(content))
		sourceFileContents[language] = append(sourceFileContents[language], markdownContent)
		return nil
	})
	return sourceFileContents, totalFileSize, skippedFileCount
}

