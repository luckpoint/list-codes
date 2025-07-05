package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// IsTestFile determines if the specified file path is a test file.
func IsTestFile(filePath string) bool {
	normalizedPath := filepath.ToSlash(strings.ToLower(filePath))
	fileName := filepath.Base(normalizedPath)

	// 1. Check if the file is in a test directory
	for _, testDirPattern := range EXCLUDE_TEST_DIRS {
		if strings.Contains(normalizedPath, testDirPattern) {
			return true
		}
	}

	// 2. Check against specific patterns
	for _, pattern := range EXCLUDE_TEST_PATTERNS {
		if strings.HasSuffix(fileName, pattern) {
			return true
		}
	}

	// 3. Check for common test file naming conventions
	baseName := strings.TrimSuffix(fileName, filepath.Ext(fileName))
	for _, keyword := range EXCLUDE_TEST_KEYWORDS {
		if strings.HasPrefix(baseName, keyword) || strings.HasSuffix(baseName, keyword) ||
			strings.HasPrefix(baseName, keyword+"_") || strings.HasSuffix(baseName, "_"+keyword) ||
			strings.HasPrefix(baseName, keyword+"-") || strings.HasSuffix(baseName, "-"+keyword) ||
			strings.Contains(baseName, "_"+keyword+"_") || strings.Contains(baseName, "-"+keyword+"-") {
			return true
		}
	}

	return false
}

// shouldSkipDir determines whether to skip directories or files.
func shouldSkipDir(fullPath, name string, isDir bool, includePaths, excludeNames, excludePaths map[string]struct{}) bool {
	absPath, err := filepath.Abs(fullPath)
	if err != nil {
		PrintWarning(fmt.Sprintf("Could not get absolute path for %s: %v", fullPath, err), true)
		return true // 安全のためスキップ
	}

	if _, ok := excludeNames[name]; ok {
		return true
	}

	if strings.HasPrefix(name, ".") && name != "." && name != ".." {
		if _, ok := DefaultExcludeNames[name]; ok {
			return true
		}
		if !isDir { // All dot files are excluded
			return true
		}
	}

	if _, ok := excludePaths[absPath]; ok {
		return true
	}

	if len(includePaths) > 0 {
		isIncluded := false
		for incPath := range includePaths {
			if strings.HasPrefix(absPath, incPath) {
				isIncluded = true
				break
			}
		}

		if !isIncluded {
			isParentOfIncluded := false
			for incPath := range includePaths {
				if strings.HasPrefix(incPath, absPath) {
					isParentOfIncluded = true
					break
				}
			}
			if !isParentOfIncluded {
				return true
			}
		}
	}

	return false
}

// GenerateDirectoryStructure generates the project directory structure in Markdown format.
func GenerateDirectoryStructure(startPath string, maxDepth int, debugMode bool, includePaths, excludeNames, excludePaths map[string]struct{}) string {
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
			absPath, _ := filepath.Abs(itemPath)
			if shouldSkipDir(absPath, entry.Name(), entry.IsDir(), includePaths, excludeNames, excludePaths) {
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
							subItemPath, _ := filepath.Abs(filepath.Join(nextPath, subEntry.Name()))
							if !shouldSkipDir(subItemPath, subEntry.Name(), subEntry.IsDir(), includePaths, excludeNames, excludePaths) {
								hasVisibleSubItems = true
								break
							}
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

		absPath, _ := filepath.Abs(path)
		if d.IsDir() {
			if shouldSkipDir(absPath, d.Name(), true, includePaths, excludeNames, excludePaths) {
			return filepath.SkipDir
			}
			return nil
		}

		if strings.ToLower(d.Name()) == "readme.md" {
			if shouldSkipDir(absPath, d.Name(), false, includePaths, excludeNames, excludePaths) {
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
					absPath, _ := filepath.Abs(path)
					if d.IsDir() {
						if shouldSkipDir(absPath, d.Name(), true, includePaths, excludeNames, excludePaths) {
							return filepath.SkipDir
						}
						return nil
					}

					if d.Name() == depFilePattern {
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
func collectSourceFiles(folderAbs string, primaryLangs []string, fallbackLangs map[string]int, processedDepFiles, includePaths, excludeNames, excludePaths map[string]struct{}, debug bool) (map[string][]string, int64, int) {
	sourceFileContents := make(map[string][]string)
	var totalFileSize int64
	var skippedFileCount int
	PrintDebug("Processing source files...", debug)

	filepath.WalkDir(folderAbs, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			PrintWarning(fmt.Sprintf("Error accessing path %s: %v", path, debug), debug)
			return nil
		}
		absPath, _ := filepath.Abs(path)

		if d.IsDir() {
			if shouldSkipDir(absPath, d.Name(), true, includePaths, excludeNames, excludePaths) {
				return filepath.SkipDir
			}
			return nil
		}

		if shouldSkipDir(absPath, d.Name(), false, includePaths, excludeNames, excludePaths) {
			return nil
		}

		if _, ok := processedDepFiles[absPath]; ok {
			return nil
		}
		if IsTestFile(path) {
			return nil
		}

		language := GetLanguageByExtension(d.Name())
		if language == "" {
			return nil
		}

		processThisFile := false
		if len(primaryLangs) > 0 {
			for _, lang := range primaryLangs {
				if lang == language {
					processThisFile = true
					break
				}
			}
		} else if len(fallbackLangs) > 0 {
			if _, ok := fallbackLangs[language]; ok {
				processThisFile = true
			}
		} else {
			processThisFile = true
		}

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