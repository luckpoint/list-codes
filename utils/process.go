package utils

import (
	"fmt"
	"sort"
	"strings"
)

// buildMarkdownOutput builds the final Markdown output.
func buildMarkdownOutput(directoryStructureMD string, depFileContents map[string][]string, sourceFileContents map[string][]string, totalFileSize int64, skippedFileCount int, debug bool) string {
	var outputMDParts []string
	outputMDParts = append(outputMDParts, directoryStructureMD)

	if len(depFileContents[DependencyFilesCategory]) > 0 {
		outputMDParts = append(outputMDParts, "## Dependency and Configuration Files\n")
		sort.Strings(depFileContents[DependencyFilesCategory])
		outputMDParts = append(outputMDParts, depFileContents[DependencyFilesCategory]...)
	}

	hasSourceFiles := false
	for lang, contents := range sourceFileContents {
		if lang != DependencyFilesCategory && len(contents) > 0 {
			hasSourceFiles = true
			break
		}
	}

	if hasSourceFiles {
		outputMDParts = append(outputMDParts, "## Source Code Files\n")
		
		// Add file size statistics
		fileSizeMB := float64(totalFileSize) / (1024 * 1024)
		maxFileSizeMB := float64(MaxFileSizeBytes) / (1024 * 1024)
		var statsText string
		if skippedFileCount > 0 {
			statsText = fmt.Sprintf("**File Statistics**: %.2f MB total, %d files skipped (> %.1f MB limit)\n\n", fileSizeMB, skippedFileCount, maxFileSizeMB)
		} else {
			statsText = fmt.Sprintf("**File Statistics**: %.2f MB total (limit: %.1f MB)\n\n", fileSizeMB, maxFileSizeMB)
		}
		outputMDParts = append(outputMDParts, statsText)
		
		languages := make([]string, 0, len(sourceFileContents))
		for lang := range sourceFileContents {
			if lang != DependencyFilesCategory {
				languages = append(languages, lang)
			}
		}
		sort.Strings(languages)

		for _, lang := range languages {
			outputMDParts = append(outputMDParts, fmt.Sprintf("### %s\n", lang))
			sort.Strings(sourceFileContents[lang])
			outputMDParts = append(outputMDParts, sourceFileContents[lang]...)
		}
	}

	return strings.Join(outputMDParts, "\n")
}

// ProcessSourceFiles processes source files and generates a Markdown summary.
func ProcessSourceFiles(folderAbs string, maxDepth int, includePaths, excludeNames, excludePaths map[string]struct{}, debug bool) string {
	directoryStructureMD := GenerateDirectoryStructure(folderAbs, maxDepth, debug, includePaths, excludeNames, excludePaths)
	primaryLangs, fallbackLangs := DetectProjectLanguages(folderAbs, debug, includePaths, excludeNames, excludePaths)

	depFileContents, processedDepFiles := collectDependencyFiles(folderAbs, primaryLangs, fallbackLangs, includePaths, excludeNames, excludePaths, debug)
	sourceFileContents, totalFileSize, skippedFileCount := collectSourceFiles(folderAbs, primaryLangs, fallbackLangs, processedDepFiles, includePaths, excludeNames, excludePaths, debug)

	return buildMarkdownOutput(directoryStructureMD, depFileContents, sourceFileContents, totalFileSize, skippedFileCount, debug)
}
