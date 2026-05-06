package utils

import (
	"path/filepath"
	"testing"
)

func TestNormalizeAbsolutePath(t *testing.T) {
	t.Run("cleans absolute path without changing root", func(t *testing.T) {
		root := t.TempDir()
		input := filepath.Join(root, "src", "..", "main.go")

		got, err := normalizeAbsolutePath(input)
		if err != nil {
			t.Fatalf("normalizeAbsolutePath(%q) returned error: %v", input, err)
		}

		want := filepath.Join(root, "main.go")
		if got != want {
			t.Fatalf("normalizeAbsolutePath(%q) = %q, want %q", input, got, want)
		}
	})

	t.Run("resolves relative path to absolute", func(t *testing.T) {
		got, err := normalizeAbsolutePath(".")
		if err != nil {
			t.Fatalf("normalizeAbsolutePath returned error: %v", err)
		}

		if !filepath.IsAbs(got) {
			t.Fatalf("normalizeAbsolutePath(.) = %q, want absolute path", got)
		}
	})
}

func TestPathWithinRoot(t *testing.T) {
	root := t.TempDir()

	testCases := []struct {
		name string
		path string
		want bool
	}{
		{
			name: "root itself",
			path: root,
			want: true,
		},
		{
			name: "direct child",
			path: filepath.Join(root, "main.go"),
			want: true,
		},
		{
			name: "nested descendant",
			path: filepath.Join(root, "src", "pkg", "main.go"),
			want: true,
		},
		{
			name: "sibling with common prefix",
			path: root + "-backup",
			want: false,
		},
		{
			name: "parent directory",
			path: filepath.Dir(root),
			want: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := pathWithinRoot(root, tc.path); got != tc.want {
				t.Fatalf("pathWithinRoot(%q, %q) = %v, want %v", root, tc.path, got, tc.want)
			}
		})
	}
}
