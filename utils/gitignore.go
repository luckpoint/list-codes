package utils

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/sabhiram/go-gitignore"
)

// GitIgnoreMatcher provides functionality to match paths against .gitignore rules
type GitIgnoreMatcher struct {
	root       string
	tree       map[string]*ignore.GitIgnore // key = directory absolute path
	loadedDirs map[string]struct{}          // absolute dirs already checked
	mu         sync.RWMutex
}

// NewGitIgnoreMatcher creates a matcher and eagerly checks only the root
// directory. Nested .gitignore files are loaded lazily on Match calls.
func NewGitIgnoreMatcher(root string) (*GitIgnoreMatcher, error) {
	absRoot, err := normalizeAbsolutePath(root)
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
	matcher.loadGitIgnoreForDir(absRoot)

	return matcher, nil
}

// SimpleMatcher matches paths against a fixed list of patterns.
type SimpleMatcher struct {
	root     string
	matcher  *ignore.GitIgnore
	patterns []string
}

// NewSimpleMatcher creates a new SimpleMatcher from a list of patterns.
func NewSimpleMatcher(root string, patterns []string) (*SimpleMatcher, error) {
	if len(patterns) == 0 {
		return &SimpleMatcher{
			root:     root,
			matcher:  nil,
			patterns: patterns,
		}, nil
	}

	absRoot, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}

	matcher := ignore.CompileIgnoreLines(patterns...)
	return &SimpleMatcher{
		root:     absRoot,
		matcher:  matcher,
		patterns: patterns,
	}, nil
}

// Match checks if the given path matches any of the patterns.
func (m *SimpleMatcher) Match(path string) bool {
	if m == nil || m.matcher == nil {
		return false
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return false
	}

	// Calculate relative path from root
	relPath, err := filepath.Rel(m.root, absPath)
	if err != nil {
		return false
	}

	// filepath.Rel returns "." if paths are equal
	if relPath == "." {
		return false
	}

	return m.matcher.MatchesPath(filepath.ToSlash(relPath))
}

// HasPatterns returns true if the matcher has any patterns configured.
func (m *SimpleMatcher) HasPatterns() bool {
	return m != nil && len(m.patterns) > 0
}

// Match checks if the given path should be ignored according to .gitignore rules.
// It returns true if the path should be ignored, false otherwise.
// The function applies Git's precedence rules: closer .gitignore files take precedence.
func (m *GitIgnoreMatcher) Match(path string) bool {
	isDir := false
	if info, err := os.Stat(path); err == nil {
		isDir = info.IsDir()
	}
	return m.MatchWithType(path, isDir)
}

// MatchWithType is a faster variant of Match when caller already knows whether
// the path is a directory.
func (m *GitIgnoreMatcher) MatchWithType(path string, isDir bool) bool {
	if m == nil {
		return false
	}

	absPath, err := normalizeAbsolutePath(path)
	if err != nil {
		return false
	}

	// Make sure the path is within our root
	if !pathWithinRoot(m.root, absPath) {
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
	if !isDir {
		currentDir = filepath.Dir(absPath)
	}

	// Collect all applicable .gitignore matchers from current to root.
	var dirs []string
	dir := currentDir
	for {
		dirs = append(dirs, dir)
		if dir == m.root {
			break
		}
		parent := filepath.Dir(dir)
		if parent == dir || !pathWithinRoot(m.root, parent) {
			break
		}
		dir = parent
	}
	for _, d := range dirs {
		m.loadGitIgnoreForDir(d)
	}

	// Apply matchers from root to current directory.
	var matchers []*ignore.GitIgnore
	m.mu.RLock()
	for i := len(dirs) - 1; i >= 0; i-- {
		if gitignore, exists := m.tree[dirs[i]]; exists {
			matchers = append(matchers, gitignore)
		}
	}
	m.mu.RUnlock()

	for _, gitignore := range matchers {
		if gitignore.MatchesPath(relPathUnix) {
			return true
		}
	}

	return false
}

func (m *GitIgnoreMatcher) loadGitIgnoreForDir(dir string) {
	m.mu.RLock()
	_, loaded := m.loadedDirs[dir]
	m.mu.RUnlock()
	if loaded {
		return
	}

	m.mu.Lock()
	if _, loaded := m.loadedDirs[dir]; loaded {
		m.mu.Unlock()
		return
	}
	m.loadedDirs[dir] = struct{}{}
	m.mu.Unlock()

	gitignorePath := filepath.Join(dir, ".gitignore")
	if _, err := os.Stat(gitignorePath); err != nil {
		if !os.IsNotExist(err) {
			PrintWarning("Could not access .gitignore file at "+gitignorePath+": "+err.Error(), true)
		}
		return
	}

	gitignore, err := ignore.CompileIgnoreFile(gitignorePath)
	if err != nil {
		PrintWarning("Could not load .gitignore file at "+gitignorePath+": "+err.Error(), true)
		return
	}

	m.mu.Lock()
	m.tree[dir] = gitignore
	m.mu.Unlock()
}
