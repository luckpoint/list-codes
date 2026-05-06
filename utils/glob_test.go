package utils

import (
	"path/filepath"
	"regexp"
	"testing"
)

func TestResolvePathPatterns_BlankPatternsAreIgnored(t *testing.T) {
	tempDir := t.TempDir()

	got := ResolvePathPatterns(tempDir, []string{"", "   ", "main.go"}, "include", false)
	want := filepath.Join(tempDir, "main.go")

	if len(got) != 1 {
		t.Fatalf("ResolvePathPatterns returned %d paths, want 1: %#v", len(got), got)
	}
	if _, ok := got[want]; !ok {
		t.Fatalf("ResolvePathPatterns did not include normalized literal path %q: %#v", want, got)
	}
}

func TestDoubleStarPatternToRegexp(t *testing.T) {
	re := regexp.MustCompile(doubleStarPatternToRegexp("/repo/docs/**/*.md"))

	testCases := []struct {
		path string
		want bool
	}{
		{path: "/repo/docs/README.md", want: true},
		{path: "/repo/docs/guide/intro.md", want: true},
		{path: "/repo/docs/guide/intro.txt", want: false},
		{path: "/repo/docs-backup/README.md", want: false},
	}

	for _, tc := range testCases {
		t.Run(tc.path, func(t *testing.T) {
			if got := re.MatchString(tc.path); got != tc.want {
				t.Fatalf("regexp match for %q = %v, want %v", tc.path, got, tc.want)
			}
		})
	}
}

func TestPatternSearchRoot(t *testing.T) {
	testCases := []struct {
		name    string
		pattern string
		want    string
	}{
		{
			name:    "recursive glob uses fixed directory prefix",
			pattern: filepath.Join(string(filepath.Separator), "repo", "src", "**", "*.go"),
			want:    filepath.Join(string(filepath.Separator), "repo", "src"),
		},
		{
			name:    "plain glob uses containing directory",
			pattern: filepath.Join("docs", "*.md"),
			want:    "docs",
		},
		{
			name:    "no glob returns the path itself",
			pattern: filepath.Join("docs", "README.md"),
			want:    filepath.Join("docs", "README.md"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := patternSearchRoot(tc.pattern); got != tc.want {
				t.Fatalf("patternSearchRoot(%q) = %q, want %q", tc.pattern, got, tc.want)
			}
		})
	}
}
