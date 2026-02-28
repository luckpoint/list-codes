package utils

import (
	"path/filepath"
	"strings"
)

// normalizeAbsolutePath returns a cleaned absolute path and avoids filepath.Abs
// when the input is already absolute.
func normalizeAbsolutePath(path string) (string, error) {
	if filepath.IsAbs(path) {
		return filepath.Clean(path), nil
	}
	return filepath.Abs(path)
}

// pathWithinRoot reports whether path is root itself or a descendant of root.
func pathWithinRoot(root, path string) bool {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return false
	}
	if rel == "." {
		return true
	}
	prefix := ".." + string(filepath.Separator)
	return rel != ".." && !strings.HasPrefix(rel, prefix)
}
