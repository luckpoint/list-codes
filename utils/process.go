package utils

import (
	"fmt"
	"sort"
	"strings"
)

// buildMarkdownOutput builds the final Markdown output.
func buildMarkdownOutput(directoryStructureMD string, depFileContents map[string][]string, sourceFileContents map[string][]string, totalFileSize int64, skippedFileMessages []string, limitHit bool, debug bool) string {
	var outputMDParts []string

	hasSourceFiles := false
	for lang, contents := range sourceFileContents {
		if lang != DependencyFilesCategory && len(contents) > 0 {
			hasSourceFiles = true
			break
		}
	}

	// Show Source Code Size Check section FIRST if we have collected files OR if we have skipped files
	if hasSourceFiles || len(skippedFileMessages) > 0 {
		outputMDParts = append(outputMDParts, "## Source Code Size Check\n")
		
		// Add file size statistics
		fileSizeMB := float64(totalFileSize) / (1024 * 1024)
		maxFileSizeMB := float64(MaxFileSizeBytes) / (1024 * 1024)
		totalLimitMB := float64(TotalMaxFileSizeBytes) / (1024 * 1024)
		
		var statsParts []string
		statsParts = append(statsParts, fmt.Sprintf("%.2f MB collected", fileSizeMB))
		statsParts = append(statsParts, fmt.Sprintf("max file: %.1f MB", maxFileSizeMB))

		if limitHit {
			statsParts = append(statsParts, fmt.Sprintf("scan stopped at %.2f MB total limit", totalLimitMB))
		} else if TotalMaxFileSizeBytes > 0 {
			statsParts = append(statsParts, fmt.Sprintf("max total: %.2f MB", totalLimitMB))
		} else {
			statsParts = append(statsParts, "max total: unlimited")
		}
		
		statsText := fmt.Sprintf("**File Statistics**: %s\n", strings.Join(statsParts, ", "))
		outputMDParts = append(outputMDParts, statsText)

		if len(skippedFileMessages) > 0 {
			sort.Strings(skippedFileMessages)
			
			var skippedList []string
			skippedList = append(skippedList, fmt.Sprintf("\n**Skipped %d file(s) > %.1f MB:**", len(skippedFileMessages), maxFileSizeMB))
			for _, msg := range skippedFileMessages {
				skippedList = append(skippedList, "- "+msg)
			}
			outputMDParts = append(outputMDParts, strings.Join(skippedList, "\n"))
		}
		
		outputMDParts = append(outputMDParts, "\n")
	}

	// Add directory structure as SECOND section (after Source Code Size Check)
	outputMDParts = append(outputMDParts, directoryStructureMD)

	// Add source code language sections after directory structure
	if hasSourceFiles {
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

	// Add dependency files section LAST
	if len(depFileContents[DependencyFilesCategory]) > 0 {
		outputMDParts = append(outputMDParts, "## Dependency and Configuration Files\n")
		sort.Strings(depFileContents[DependencyFilesCategory])
		outputMDParts = append(outputMDParts, depFileContents[DependencyFilesCategory]...)
	}

	return strings.Join(outputMDParts, "\n")
}

// ProcessSourceFiles processes source files and generates a Markdown summary.
func ProcessSourceFiles(folderAbs string, maxDepth int, includePaths, excludeNames, excludePaths map[string]struct{}, debug bool, includeTests bool, gi *GitIgnoreMatcher) string {
	directoryStructureMD := GenerateDirectoryStructure(folderAbs, maxDepth, debug, includePaths, excludeNames, excludePaths, includeTests, gi)
	primaryLangs, fallbackLangs := DetectProjectLanguages(folderAbs, debug, includePaths, excludeNames, excludePaths, gi)

	depFileContents, processedDepFiles := collectDependencyFiles(folderAbs, primaryLangs, fallbackLangs, includePaths, excludeNames, excludePaths, debug, gi)
	sourceFileContents, totalFileSize, skippedFileMessages, limitHit := collectSourceFiles(folderAbs, primaryLangs, fallbackLangs, processedDepFiles, includePaths, excludeNames, excludePaths, debug, includeTests, gi)

	return buildMarkdownOutput(directoryStructureMD, depFileContents, sourceFileContents, totalFileSize, skippedFileMessages, limitHit, debug)
}
