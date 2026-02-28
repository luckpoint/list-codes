package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

var doubleStarRegexpCache sync.Map

// ResolvePathPatterns resolves include/exclude patterns against baseFolder.
// Supports filepath.Glob syntax and `**` recursive wildcard patterns.
func ResolvePathPatterns(baseFolder string, patterns []string, flagName string, debug bool) map[string]struct{} {
	resolvedPaths := make(map[string]struct{})

	for _, rawPattern := range patterns {
		pattern := strings.TrimSpace(rawPattern)
		if pattern == "" {
			continue
		}

		fullPattern := pattern
		if !filepath.IsAbs(fullPattern) {
			fullPattern = filepath.Join(baseFolder, fullPattern)
		}

		expanded, err := expandPathPattern(fullPattern)
		if err != nil {
			PrintWarning(fmt.Sprintf("Could not expand %s pattern '%s': %v", flagName, pattern, err), debug)
			continue
		}

		for _, p := range expanded {
			absPath, err := normalizeAbsolutePath(p)
			if err != nil {
				PrintWarning(fmt.Sprintf("Could not resolve absolute path for %s '%s': %v", flagName, p, err), debug)
				continue
			}
			resolvedPaths[absPath] = struct{}{}
		}
	}

	return resolvedPaths
}

func expandPathPattern(pattern string) ([]string, error) {
	cleanPattern := filepath.Clean(pattern)
	if !hasGlobMeta(cleanPattern) {
		return []string{cleanPattern}, nil
	}

	if !strings.Contains(cleanPattern, "**") {
		matches, err := filepath.Glob(cleanPattern)
		if err != nil {
			return nil, err
		}
		if len(matches) == 0 {
			// Keep the literal pattern to ensure include-only mode still takes effect.
			return []string{cleanPattern}, nil
		}
		return matches, nil
	}

	searchRoot := patternSearchRoot(cleanPattern)
	searchRootAbs, err := normalizeAbsolutePath(searchRoot)
	if err != nil {
		return nil, err
	}

	re, err := getCachedDoubleStarRegexp(cleanPattern)
	if err != nil {
		return nil, err
	}

	var matches []string
	walkErr := filepath.WalkDir(searchRootAbs, func(path string, d os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return nil
		}
		normalizedPath := filepath.ToSlash(filepath.Clean(path))
		if re.MatchString(normalizedPath) {
			matches = append(matches, path)
		}
		return nil
	})
	if walkErr != nil {
		return nil, walkErr
	}

	if len(matches) == 0 {
		// Keep the literal pattern to ensure include-only mode still takes effect.
		return []string{cleanPattern}, nil
	}
	return matches, nil
}

func getCachedDoubleStarRegexp(pattern string) (*regexp.Regexp, error) {
	normalized := filepath.ToSlash(filepath.Clean(pattern))
	if cached, ok := doubleStarRegexpCache.Load(normalized); ok {
		return cached.(*regexp.Regexp), nil
	}

	re, err := regexp.Compile(doubleStarPatternToRegexp(normalized))
	if err != nil {
		return nil, err
	}
	if cached, loaded := doubleStarRegexpCache.LoadOrStore(normalized, re); loaded {
		return cached.(*regexp.Regexp), nil
	}
	return re, nil
}

func hasGlobMeta(pattern string) bool {
	return strings.ContainsAny(pattern, "*?[")
}

func patternSearchRoot(pattern string) string {
	idx := strings.IndexAny(pattern, "*?[")
	if idx < 0 {
		return pattern
	}

	prefix := pattern[:idx]
	root := filepath.Dir(prefix)
	if root == "." || root == "" {
		return string(filepath.Separator)
	}
	return root
}

func doubleStarPatternToRegexp(pattern string) string {
	var b strings.Builder
	b.WriteString("^")

	for i := 0; i < len(pattern); {
		ch := pattern[i]
		switch ch {
		case '*':
			if i+1 < len(pattern) && pattern[i+1] == '*' {
				if i+2 < len(pattern) && pattern[i+2] == '/' {
					// `**/` means zero or more nested directories.
					b.WriteString("(?:[^/]+/)*")
					i += 3
					continue
				}
				b.WriteString(".*")
				i += 2
				continue
			}
			b.WriteString("[^/]*")
		case '?':
			b.WriteString("[^/]")
		default:
			b.WriteString(regexp.QuoteMeta(string(ch)))
		}
		i++
	}

	b.WriteString("$")
	return b.String()
}
