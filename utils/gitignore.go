package utils

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/sabhiram/go-gitignore"
)

// GitIgnoreMatcher provides functionality to match paths against .gitignore rules
type GitIgnoreMatcher struct {
	root string
	tree map[string]*ignore.GitIgnore // key = directory absolute path
}

// NewGitIgnoreMatcher creates a new GitIgnoreMatcher by walking the repository
// and loading all .gitignore files found in the directory hierarchy.
func NewGitIgnoreMatcher(root string) (*GitIgnoreMatcher, error) {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}

	matcher := &GitIgnoreMatcher{
		root: absRoot,
		tree: make(map[string]*ignore.GitIgnore),
	}

	// Walk the directory tree and load .gitignore files
	err = filepath.WalkDir(absRoot, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil // Continue walking despite errors
		}

		if !d.IsDir() && d.Name() == ".gitignore" {
			dir := filepath.Dir(path)
			
			// Load the .gitignore file
			gitignore, loadErr := ignore.CompileIgnoreFile(path)
			if loadErr != nil {
				// Log warning but continue processing
				PrintWarning("Could not load .gitignore file at "+path+": "+loadErr.Error(), true)
				return nil
			}
			
			matcher.tree[dir] = gitignore
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return matcher, nil
}

// Match checks if the given path should be ignored according to .gitignore rules.
// It returns true if the path should be ignored, false otherwise.
// The function applies Git's precedence rules: closer .gitignore files take precedence.
func (m *GitIgnoreMatcher) Match(path string) bool {
	if m == nil {
		return false
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return false
	}

	// Make sure the path is within our root
	if !strings.HasPrefix(absPath, m.root) {
		return false
	}

	// For gitignore matching, we need to use relative path from the repository root
	relPath, err := filepath.Rel(m.root, absPath)
	if err != nil {
		return false
	}
	
	// Convert to forward slashes for gitignore matching
	relPathUnix := filepath.ToSlash(relPath)

	// Start from the file's directory and walk up to root
	currentDir := absPath
	if info, err := os.Stat(absPath); err == nil && !info.IsDir() {
		currentDir = filepath.Dir(absPath)
	}

	// Collect all applicable .gitignore matchers from closest to root
	var matchers []*ignore.GitIgnore
	for {
		if gitignore, exists := m.tree[currentDir]; exists {
			matchers = append(matchers, gitignore)
		}

		// Move up one directory
		parent := filepath.Dir(currentDir)
		if parent == currentDir || !strings.HasPrefix(parent, m.root) {
			break
		}
		currentDir = parent
	}

	// Apply matchers in reverse order (closest first, which matches Git's behavior)
	// If any matcher says to ignore, we ignore (unless a closer one says not to)
	for i := len(matchers) - 1; i >= 0; i-- {
		if matchers[i].MatchesPath(relPathUnix) {
			return true
		}
	}

	return false
}