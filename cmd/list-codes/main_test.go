package main

import (
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/luckpoint/list-codes/utils"
	"github.com/spf13/cobra"
)

func TestResolvePathPatternsGlob(t *testing.T) {
	tempDir := t.TempDir()

	readme := filepath.Join(tempDir, "README.md")
	goFile := filepath.Join(tempDir, "main.go")
	if err := os.WriteFile(readme, []byte("# README"), 0644); err != nil {
		t.Fatalf("failed to create README.md: %v", err)
	}
	if err := os.WriteFile(goFile, []byte("package main"), 0644); err != nil {
		t.Fatalf("failed to create main.go: %v", err)
	}

	got := utils.ResolvePathPatterns(tempDir, []string{"*.md"}, "include", false)
	if _, ok := got[readme]; !ok {
		t.Fatalf("expected README.md to be resolved from glob pattern")
	}
	if _, ok := got[goFile]; ok {
		t.Fatalf("did not expect main.go to be included for *.md")
	}
}

func TestResolvePathPatternsDoubleStar(t *testing.T) {
	tempDir := t.TempDir()

	docsRoot := filepath.Join(tempDir, "docs")
	nested := filepath.Join(docsRoot, "guide")
	if err := os.MkdirAll(nested, 0755); err != nil {
		t.Fatalf("failed to create docs dirs: %v", err)
	}

	topLevelMD := filepath.Join(docsRoot, "prompts.md")
	nestedMD := filepath.Join(nested, "intro.md")
	otherFile := filepath.Join(docsRoot, "index.html")

	for path, content := range map[string]string{
		topLevelMD: "# prompts",
		nestedMD:   "# intro",
		otherFile:  "<html/>",
	} {
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("failed to create file %s: %v", path, err)
		}
	}

	got := utils.ResolvePathPatterns(tempDir, []string{"docs/**/*.md"}, "include", false)
	if _, ok := got[topLevelMD]; !ok {
		t.Fatalf("expected top-level markdown under docs to match docs/**/*.md")
	}
	if _, ok := got[nestedMD]; !ok {
		t.Fatalf("expected nested markdown under docs to match docs/**/*.md")
	}
	if _, ok := got[otherFile]; ok {
		t.Fatalf("did not expect non-markdown file to match docs/**/*.md")
	}
}

func TestResolvePathPatternsNoMatchKeepsLiteralPattern(t *testing.T) {
	tempDir := t.TempDir()

	got := utils.ResolvePathPatterns(tempDir, []string{"missing/*.md"}, "include", false)
	expectedLiteral := filepath.Join(tempDir, "missing", "*.md")
	expectedLiteralAbs, err := filepath.Abs(expectedLiteral)
	if err != nil {
		t.Fatalf("failed to build expected absolute path: %v", err)
	}

	if _, ok := got[expectedLiteralAbs]; !ok {
		t.Fatalf("expected unresolved glob pattern to be kept as literal include path")
	}
}

func TestJoinSetSortsKeys(t *testing.T) {
	got := joinSet(map[string]struct{}{
		"zeta":  {},
		"alpha": {},
		"beta":  {},
	})

	if got != "alpha, beta, zeta" {
		t.Fatalf("joinSet returned %q, want sorted keys", got)
	}
}

func TestPromptCompletionReturnsSortedPromptNames(t *testing.T) {
	got, directive := promptCompletion(nil, nil, "")

	if directive != cobra.ShellCompDirectiveNoFileComp {
		t.Fatalf("promptCompletion directive = %v, want %v", directive, cobra.ShellCompDirectiveNoFileComp)
	}
	if !slices.IsSorted(got) {
		t.Fatalf("promptCompletion returned unsorted names: %#v", got)
	}
	if len(got) == 0 {
		t.Fatalf("promptCompletion returned no prompt names")
	}
}
