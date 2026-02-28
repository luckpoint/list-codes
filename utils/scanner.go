package utils

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type scannedEntry struct {
	name  string
	path  string
	isDir bool
}

type scanResult struct {
	rootPath              string
	structureChildren     map[string][]scannedEntry
	readmeFiles           []string
	sourceFileContents    map[string][]string
	totalFileSize         int64
	skippedFileMessages   []string
	limitHit              bool
	stopCollectingSources bool
}

type projectScanner struct {
	rootPath string

	includePaths map[string]struct{}
	excludeNames map[string]struct{}
	excludePaths map[string]struct{}

	debug        bool
	includeTests bool
	gi           *GitIgnoreMatcher

	processedDepFiles map[string]struct{}

	collectStructure bool
	collectReadmes   bool
	collectSources   bool

	// stopWalkOnSourceLimit preserves collectSourceFiles behaviour where scan
	// stops immediately once total-size limit is exceeded.
	stopWalkOnSourceLimit bool
}

type fileScanMeta struct {
	isTest             bool
	isAsset            bool
	explicitlyIncluded bool
}

func (s *projectScanner) scan() (*scanResult, error) {
	absRoot, err := normalizeAbsolutePath(s.rootPath)
	if err != nil {
		return nil, err
	}

	result := &scanResult{
		rootPath: absRoot,
	}
	if s.collectStructure {
		result.structureChildren = make(map[string][]scannedEntry)
		result.structureChildren[absRoot] = nil
	}
	if s.collectSources {
		result.sourceFileContents = make(map[string][]string)
	}

	walkErr := filepath.WalkDir(absRoot, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			PrintWarning(fmt.Sprintf("Error accessing path %s: %v", path, walkErr), s.debug)
			if d != nil && d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if path == absRoot {
			return nil
		}

		absPath := path
		if shouldSkipPath(absPath, d.Name(), d.IsDir(), s.includePaths, s.excludeNames, s.excludePaths, s.gi) {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if d.IsDir() {
			if s.collectStructure {
				result.structureChildren[filepath.Dir(absPath)] = append(result.structureChildren[filepath.Dir(absPath)], scannedEntry{
					name:  d.Name(),
					path:  absPath,
					isDir: true,
				})
				if _, ok := result.structureChildren[absPath]; !ok {
					result.structureChildren[absPath] = nil
				}
			}
			return nil
		}

		var meta fileScanMeta
		if s.collectStructure || s.collectSources {
			meta = buildFileScanMeta(absPath, s.includePaths, s.includeTests, s.debug)
		}

		if s.collectStructure && shouldIncludeInStructure(meta) {
			result.structureChildren[filepath.Dir(absPath)] = append(result.structureChildren[filepath.Dir(absPath)], scannedEntry{
				name:  d.Name(),
				path:  absPath,
				isDir: false,
			})
		}

		if s.collectReadmes && strings.EqualFold(d.Name(), "readme.md") {
			s.collectReadmeFile(result, absPath)
		}

		if s.collectSources {
			if err := s.collectSourceFile(result, absPath, d, meta); err != nil {
				return err
			}
		}

		return nil
	})

	if walkErr != nil && !errors.Is(walkErr, errTotalSizeLimitExceeded) {
		PrintWarning(fmt.Sprintf("Error during file walk: %v", walkErr), s.debug)
		return result, walkErr
	}

	return result, nil
}

func buildFileScanMeta(absPath string, includePaths map[string]struct{}, includeTests bool, debug bool) fileScanMeta {
	meta := fileScanMeta{
		isAsset:            IsAssetFile(absPath, debug),
		explicitlyIncluded: isExplicitlyIncludedPath(absPath, includePaths),
	}
	if !includeTests {
		meta.isTest = IsTestFile(absPath, debug)
	}
	return meta
}

func shouldIncludeInStructure(meta fileScanMeta) bool {
	if meta.isTest {
		return false
	}
	if meta.isAsset && !meta.explicitlyIncluded {
		return false
	}
	return true
}

func isExplicitlyIncludedPath(absPath string, includePaths map[string]struct{}) bool {
	return len(includePaths) > 0 && isPathRelatedToIncludes(absPath, includePaths)
}

func (s *projectScanner) collectReadmeFile(result *scanResult, absPath string) {
	content, err := os.ReadFile(absPath)
	if err != nil {
		PrintWarning(fmt.Sprintf("Could not read README file '%s': %v", absPath, err), s.debug)
		return
	}

	fileDisplayName := relativeDisplayPath(result.rootPath, absPath, s.debug)
	markdownContent := fmt.Sprintf("### %s\n```markdown\n%s\n```\n", fileDisplayName, string(content))
	result.readmeFiles = append(result.readmeFiles, markdownContent)
}

func (s *projectScanner) collectSourceFile(result *scanResult, absPath string, d os.DirEntry, meta fileScanMeta) error {
	if result.stopCollectingSources {
		return nil
	}
	if _, ok := s.processedDepFiles[absPath]; ok {
		return nil
	}
	if meta.isTest {
		return nil
	}

	language := GetLanguageByExtension(d.Name())
	if meta.isAsset {
		if !meta.explicitlyIncluded {
			return nil
		}
		if language == "" {
			language = "text"
		}
	} else if language == "" {
		return nil
	}

	fileInfo, err := d.Info()
	if err != nil {
		PrintWarning(fmt.Sprintf("Could not get file info for '%s': %v", absPath, err), s.debug)
		return nil
	}

	if fileInfo.Size() > MaxFileSizeBytes {
		PrintDebug(fmt.Sprintf("Skipping file '%s' due to size (%d bytes > %d bytes)", absPath, fileInfo.Size(), MaxFileSizeBytes), s.debug)
		relPath := relativeDisplayPath(result.rootPath, absPath, s.debug)
		fileSizeMB := float64(fileInfo.Size()) / (1024 * 1024)
		skippedMessage := fmt.Sprintf("`%s` (%.2f MB)", relPath, fileSizeMB)
		result.skippedFileMessages = append(result.skippedFileMessages, skippedMessage)
		return nil
	}

	if TotalMaxFileSizeBytes > 0 && result.totalFileSize+fileInfo.Size() > TotalMaxFileSizeBytes {
		PrintDebug(fmt.Sprintf("Stopping scan as total size limit of %d bytes would be exceeded.", TotalMaxFileSizeBytes), s.debug)
		result.limitHit = true
		if s.stopWalkOnSourceLimit {
			return errTotalSizeLimitExceeded
		}
		result.stopCollectingSources = true
		return nil
	}

	result.totalFileSize += fileInfo.Size()

	content, err := os.ReadFile(absPath)
	if err != nil {
		PrintWarning(fmt.Sprintf("Could not read source file '%s': %v", absPath, err), s.debug)
		return nil
	}

	fileDisplayName := relativeDisplayPath(result.rootPath, absPath, s.debug)
	codeBlockLangHint := strings.ToLower(language)
	codeBlockLangHint = strings.ReplaceAll(codeBlockLangHint, "/", "")
	codeBlockLangHint = strings.ReplaceAll(codeBlockLangHint, "+", "p")
	markdownContent := fmt.Sprintf("### %s\n```%s\n%s\n```\n", fileDisplayName, codeBlockLangHint, string(content))
	result.sourceFileContents[language] = append(result.sourceFileContents[language], markdownContent)
	return nil
}

func relativeDisplayPath(root, absPath string, debug bool) string {
	relPath, err := filepath.Rel(root, absPath)
	if err != nil {
		PrintWarning(fmt.Sprintf("Could not get relative path for %s: %v", absPath, err), debug)
		return filepath.ToSlash(absPath)
	}
	return filepath.ToSlash(relPath)
}

func buildDirectoryStructureMarkdown(rootPath string, maxDepth int, children map[string][]scannedEntry) string {
	var structureLines []string
	structureLines = append(structureLines, "## Project Structure", "```text")

	rootDisplayName := filepath.Base(rootPath)
	structureLines = append(structureLines, fmt.Sprintf(". (%s)", rootDisplayName))

	var generateTreeRecursive func(currentPath, prefix string, depth int)
	generateTreeRecursive = func(currentPath, prefix string, depth int) {
		entries := append([]scannedEntry(nil), children[currentPath]...)
		sort.Slice(entries, func(i, j int) bool {
			return strings.ToLower(entries[i].name) < strings.ToLower(entries[j].name)
		})

		var filesToShow []scannedEntry
		var dirsToShow []scannedEntry
		var hasHiddenDirs bool

		for _, entry := range entries {
			if entry.isDir {
				if maxDepth > 0 && depth < maxDepth {
					dirsToShow = append(dirsToShow, entry)
				} else {
					hasHiddenDirs = true
				}
			} else {
				filesToShow = append(filesToShow, entry)
			}
		}

		entriesToShow := append(filesToShow, dirsToShow...)
		showEllipsis := hasHiddenDirs

		for i, entry := range entriesToShow {
			pointer := "├── "
			if i == len(entriesToShow)-1 && !showEllipsis {
				pointer = "└── "
			}

			structureLines = append(structureLines, prefix+pointer+entry.name)

			if entry.isDir {
				extension := "│   "
				if pointer == "└── " && !showEllipsis {
					extension = "    "
				}
				generateTreeRecursive(entry.path, prefix+extension, depth+1)
			}
		}

		if showEllipsis {
			structureLines = append(structureLines, prefix+"└── ...")
		}
	}

	generateTreeRecursive(rootPath, "", 0)
	structureLines = append(structureLines, "```")
	return strings.Join(structureLines, "\n") + "\n\n"
}
