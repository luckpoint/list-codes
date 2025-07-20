package utils

import (
	"fmt"
	"sort"
	"strings"
	"strconv"
	"unicode"
)

// formatNumber adds commas to large numbers for better readability
func formatNumber(n int64) string {
	s := strconv.FormatInt(n, 10)
	if len(s) < 4 {
		return s
	}
	
	// Add commas every 3 digits from right to left
	result := ""
	for i, digit := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			result += ","
		}
		result += string(digit)
	}
	return result
}

// isKanji returns true if the rune is a kanji character
func isKanji(r rune) bool {
	return unicode.In(r, unicode.Han)
}

// isHiragana returns true if the rune is a hiragana character
func isHiragana(r rune) bool {
	return unicode.In(r, unicode.Hiragana)
}

// isKatakana returns true if the rune is a katakana character
func isKatakana(r rune) bool {
	return unicode.In(r, unicode.Katakana)
}

// isAscii returns true if the rune is ASCII (English/numbers/symbols)
func isAscii(r rune) bool {
	return r <= 127
}

// estimateTokenCountFromContent estimates token count based on actual content analysis
func estimateTokenCountFromContent(content string) float64 {
	var tokenCount float64
	
	for _, r := range content {
		if isKanji(r) {
			// Kanji: 2 tokens per character
			tokenCount += 2.0
		} else if isHiragana(r) || isKatakana(r) {
			// Hiragana/Katakana: 1 token per character
			tokenCount += 1.0
		} else if isAscii(r) {
			// ASCII (English): 0.75 tokens per character
			tokenCount += 0.75
		} else {
			// Other Unicode characters: default to 1 token
			tokenCount += 1.0
		}
	}
	
	return tokenCount
}

// estimateTokenCount estimates token count based on file content analysis
func estimateTokenCount(totalFileSize int64) int64 {
	// For backward compatibility, if we can't get content, use simple estimation
	if IsJapanese() {
		// Japanese: approximately 1.5 characters per token (weighted average)
		return int64(float64(totalFileSize) / 1.5)
	} else {
		// English: approximately 0.75 tokens per character
		return int64(float64(totalFileSize) * 0.75)
	}
}

// buildMarkdownOutput builds the final Markdown output.
func buildMarkdownOutput(directoryStructureMD string, depFileContents map[string][]string, sourceFileContents map[string][]string, totalFileSize int64, skippedFileCount int, debug bool, allContent string) string {
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
			statsText = fmt.Sprintf(T("**File Statistics**: %.2f MB total, %d files skipped (> %.1f MB limit)", "**ファイル統計**: 合計%.2f MB、%d個のファイルをスキップ (> %.1f MB制限)"), fileSizeMB, skippedFileCount, maxFileSizeMB)
		} else {
			statsText = fmt.Sprintf(T("**File Statistics**: %.2f MB total (limit: %.1f MB)", "**ファイル統計**: 合計%.2f MB (上限: %.1f MB)"), fileSizeMB, maxFileSizeMB)
		}
		
		// Add token estimation - use content-based calculation if available
		var estimatedTokens int64
		if allContent != "" {
			estimatedTokens = int64(estimateTokenCountFromContent(allContent))
		} else {
			estimatedTokens = estimateTokenCount(totalFileSize)
		}
		tokenText := fmt.Sprintf(T("**Token Estimate**: ~%s tokens (Note: As a rule of thumb, 100 tokens are roughly 75 words.)", "**トークン見積もり**: 約%sトークン (注: 漢字は2トークン、ひらがな・カタカナは1トークン、英語は0.75トークンで計算されています。)"), formatNumber(estimatedTokens))
		
		outputMDParts = append(outputMDParts, statsText+"\n"+tokenText+"\n")
		
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
func ProcessSourceFiles(folderAbs string, maxDepth int, includePaths, excludeNames, excludePaths map[string]struct{}, debug bool, includeTests bool) string {
	directoryStructureMD := GenerateDirectoryStructure(folderAbs, maxDepth, debug, includePaths, excludeNames, excludePaths, includeTests)
	primaryLangs, fallbackLangs := DetectProjectLanguages(folderAbs, debug, includePaths, excludeNames, excludePaths)

	depFileContents, processedDepFiles := collectDependencyFiles(folderAbs, primaryLangs, fallbackLangs, includePaths, excludeNames, excludePaths, debug)
	sourceFileContents, totalFileSize, skippedFileCount, allContent := collectSourceFiles(folderAbs, primaryLangs, fallbackLangs, processedDepFiles, includePaths, excludeNames, excludePaths, debug, includeTests)

	return buildMarkdownOutput(directoryStructureMD, depFileContents, sourceFileContents, totalFileSize, skippedFileCount, debug, allContent)
}
