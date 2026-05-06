package tui

import (
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTestProject(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()

	// Create directory structure
	dirs := []string{
		"src",
		"src/pkg",
		"cmd",
	}
	for _, d := range dirs {
		require.NoError(t, os.MkdirAll(filepath.Join(dir, d), 0755))
	}

	// Create files
	files := map[string]string{
		"src/main.go":      "package main",
		"src/main_test.go": "package main",
		"src/pkg/util.go":  "package pkg",
		"cmd/root.go":      "package cmd",
		"README.md":        "# Test",
	}
	for name, content := range files {
		require.NoError(t, os.WriteFile(filepath.Join(dir, name), []byte(content), 0644))
	}

	return dir
}

func findTreeNode(root *TreeNode, path string) *TreeNode {
	if root.Path == path {
		return root
	}
	for _, child := range root.Children {
		if found := findTreeNode(child, path); found != nil {
			return found
		}
	}
	return nil
}

func TestBuildTree(t *testing.T) {
	dir := createTestProject(t)

	root, err := BuildTree(dir, BuildTreeOpts{})
	require.NoError(t, err)

	assert.True(t, root.IsDir)
	assert.Equal(t, ".", root.Path)
	assert.True(t, len(root.Children) > 0)
}

func TestBuildTree_InvalidDir(t *testing.T) {
	_, err := BuildTree("/nonexistent/dir/xxxxx", BuildTreeOpts{})
	// BuildTree itself doesn't error on nonexistent because os.ReadDir fails gracefully
	// But it should still return a root node
	assert.NoError(t, err)
}

func TestBuildTree_MaxDepth(t *testing.T) {
	dir := createTestProject(t)

	root, err := BuildTree(dir, BuildTreeOpts{MaxDepth: 1})
	require.NoError(t, err)

	src := findTreeNode(root, "src")
	require.NotNil(t, src)
	assert.Empty(t, src.Children, "children deeper than MaxDepth should not be loaded")
}

func TestFlattenVisible(t *testing.T) {
	dir := createTestProject(t)

	root, err := BuildTree(dir, BuildTreeOpts{})
	require.NoError(t, err)

	// Root is expanded, children are collapsed
	visible := FlattenVisible(root)
	// Should show root + top-level children (cmd, src, README.md)
	assert.True(t, len(visible) >= 4)

	// Expand first dir child
	for _, child := range root.Children {
		if child.IsDir {
			child.Expanded = true
			break
		}
	}

	expanded := FlattenVisible(root)
	assert.True(t, len(expanded) > len(visible))
}

func TestToggle(t *testing.T) {
	dir := createTestProject(t)

	root, err := BuildTree(dir, BuildTreeOpts{})
	require.NoError(t, err)

	// Find a file node
	var fileNode *TreeNode
	for _, child := range root.Children {
		if !child.IsDir {
			fileNode = child
			break
		}
	}

	if fileNode != nil {
		assert.Equal(t, Unchecked, fileNode.State)
		Toggle(fileNode)
		assert.Equal(t, Checked, fileNode.State)
		Toggle(fileNode)
		assert.Equal(t, Unchecked, fileNode.State)
	}
}

func TestToggle_DirPropagation(t *testing.T) {
	dir := createTestProject(t)

	root, err := BuildTree(dir, BuildTreeOpts{})
	require.NoError(t, err)

	// Find a dir with children
	var dirNode *TreeNode
	for _, child := range root.Children {
		if child.IsDir && len(child.Children) > 0 {
			dirNode = child
			break
		}
	}
	require.NotNil(t, dirNode)

	// Toggle dir on - all children should be checked
	Toggle(dirNode)
	assert.Equal(t, Checked, dirNode.State)
	for _, child := range dirNode.Children {
		assert.Equal(t, Checked, child.State)
	}

	// Toggle dir off - all children should be unchecked
	Toggle(dirNode)
	assert.Equal(t, Unchecked, dirNode.State)
	for _, child := range dirNode.Children {
		assert.Equal(t, Unchecked, child.State)
	}
}

func TestToggle_ParentPartial(t *testing.T) {
	dir := createTestProject(t)

	root, err := BuildTree(dir, BuildTreeOpts{})
	require.NoError(t, err)

	// Find a dir with multiple children
	var dirNode *TreeNode
	for _, child := range root.Children {
		if child.IsDir && len(child.Children) >= 2 {
			dirNode = child
			break
		}
	}

	if dirNode != nil {
		// Check only the first child
		Toggle(dirNode.Children[0])
		assert.Equal(t, Partial, dirNode.State)
	}
}

func TestSetInitialState(t *testing.T) {
	dir := createTestProject(t)

	root, err := BuildTree(dir, BuildTreeOpts{})
	require.NoError(t, err)

	SetInitialState(root, nil, nil)

	// .go files should be checked, README.md should be checked (Markdown is a known ext)
	selected, total := CountSelected(root)
	assert.True(t, selected > 0)
	assert.True(t, total > 0)
}

func TestGeneratePatterns(t *testing.T) {
	dir := createTestProject(t)

	root, err := BuildTree(dir, BuildTreeOpts{})
	require.NoError(t, err)

	SetInitialState(root, nil, nil)

	includes, _ := GeneratePatterns(root)
	assert.True(t, len(includes) > 0)
}

func TestCountSelected(t *testing.T) {
	dir := createTestProject(t)

	root, err := BuildTree(dir, BuildTreeOpts{})
	require.NoError(t, err)

	selected, total := CountSelected(root)
	assert.Equal(t, 0, selected)
	assert.True(t, total > 0)

	SetAllState(root, Checked)
	selected, total = CountSelected(root)
	assert.Equal(t, total, selected)
}

func TestSetInitialState_ExcludePatterns(t *testing.T) {
	dir := createTestProject(t)

	root, err := BuildTree(dir, BuildTreeOpts{})
	require.NoError(t, err)

	SetInitialState(root, nil, []string{"*.md"})

	// README.md should still be in the tree but unchecked
	var found bool
	for _, child := range root.Children {
		if child.Name == "README.md" {
			found = true
			assert.Equal(t, Unchecked, child.State, "excluded file should be unchecked")
		}
	}
	assert.True(t, found, "README.md should still be in tree")

	// .go files should still be checked
	selected, _ := CountSelected(root)
	assert.True(t, selected > 0, "non-excluded files should be checked")
}

func TestSetInitialState_IncludeAndExclude(t *testing.T) {
	dir := createTestProject(t)

	root, err := BuildTree(dir, BuildTreeOpts{})
	require.NoError(t, err)

	// Include src/**, then exclude test files
	SetInitialState(root, []string{"src/**"}, []string{"**/*_test.go"})

	for _, child := range root.Children {
		if child.IsDir && child.Name == "src" {
			for _, f := range child.Children {
				if !f.IsDir {
					if f.Name == "main_test.go" {
						assert.Equal(t, Unchecked, f.State, "test file should be excluded")
					} else if f.Name == "main.go" {
						assert.Equal(t, Checked, f.State, "non-test file should be included")
					}
				}
			}
		}
		if child.IsDir && child.Name == "cmd" {
			assert.Equal(t, Unchecked, child.State, "cmd/ should be unchecked (not in include)")
		}
	}
}

func TestSetInitialState_IncludePatterns(t *testing.T) {
	dir := createTestProject(t)

	root, err := BuildTree(dir, BuildTreeOpts{})
	require.NoError(t, err)

	SetInitialState(root, []string{"src/**"}, nil)

	// Files under src/ should be checked, others unchecked
	for _, child := range root.Children {
		if child.IsDir && child.Name == "cmd" {
			assert.Equal(t, Unchecked, child.State, "cmd/ should be unchecked")
		}
		if !child.IsDir && child.Name == "README.md" {
			assert.Equal(t, Unchecked, child.State, "README.md should be unchecked")
		}
		if child.IsDir && child.Name == "src" {
			assert.NotEqual(t, Unchecked, child.State, "src/ should have checked files")
		}
	}
}

func TestSetInitialState_IncludePatterns_Empty(t *testing.T) {
	dir := createTestProject(t)

	root, err := BuildTree(dir, BuildTreeOpts{})
	require.NoError(t, err)

	// Empty include patterns should fall back to language-based auto-detection
	SetInitialState(root, []string{}, nil)

	selected, _ := CountSelected(root)
	assert.True(t, selected > 0, "empty include patterns should use language-based detection")
}

func TestApplyConfig(t *testing.T) {
	dir := createTestProject(t)

	root, err := BuildTree(dir, BuildTreeOpts{})
	require.NoError(t, err)

	cfg := &Config{
		Include: []string{"src/**"},
	}
	ApplyConfig(root, cfg)

	// Files under src/ should be checked
	for _, child := range root.Children {
		if child.IsDir && child.Name == "src" {
			assert.NotEqual(t, Unchecked, child.State)
		}
	}
}

func TestApplyConfig_ExcludeOverridesInclude(t *testing.T) {
	dir := createTestProject(t)

	root, err := BuildTree(dir, BuildTreeOpts{})
	require.NoError(t, err)

	ApplyConfig(root, &Config{
		Include: []string{"src/**"},
		Exclude: []string{"src/main.go"},
	})

	mainGo := findTreeNode(root, "src/main.go")
	utilGo := findTreeNode(root, "src/pkg/util.go")
	require.NotNil(t, mainGo)
	require.NotNil(t, utilGo)

	assert.Equal(t, Unchecked, mainGo.State)
	assert.Equal(t, Checked, utilGo.State)
}

func TestGeneratePatterns_CollapsesCheckedDirectory(t *testing.T) {
	dir := createTestProject(t)

	root, err := BuildTree(dir, BuildTreeOpts{})
	require.NoError(t, err)

	SetAllState(root, Unchecked)
	src := findTreeNode(root, "src")
	require.NotNil(t, src)
	Toggle(src)

	includes, excludes := GeneratePatterns(root)

	assert.Contains(t, includes, "src/**")
	assert.Empty(t, excludes)
	assert.False(t, slices.Contains(includes, "src/main.go"), "checked directories should be represented by a single glob")
}

func TestMatchGlobDoubleStar(t *testing.T) {
	testCases := []struct {
		name    string
		pattern string
		path    string
		want    bool
	}{
		{
			name:    "matches root-level file",
			pattern: "**/*.go",
			path:    "main.go",
			want:    true,
		},
		{
			name:    "matches nested file",
			pattern: "**/*.go",
			path:    "src/pkg/util.go",
			want:    true,
		},
		{
			name:    "respects fixed prefix",
			pattern: "src/**/*.go",
			path:    "cmd/root.go",
			want:    false,
		},
		{
			name:    "rejects wrong extension",
			pattern: "src/**/*.go",
			path:    "src/README.md",
			want:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, matchGlob(tc.pattern, tc.path))
		})
	}
}
