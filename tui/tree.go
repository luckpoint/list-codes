package tui

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/luckpoint/list-codes/utils"
)

type NodeState int

const (
	Unchecked NodeState = iota
	Checked
	Partial
)

type TreeNode struct {
	Name     string
	Path     string // relative path from root
	IsDir    bool
	State    NodeState
	Children []*TreeNode
	Parent   *TreeNode
	Expanded bool
	Depth    int
}

type BuildTreeOpts struct {
	IncludeTests    bool
	MaxDepth        int
	ExcludePatterns []string
	IncludePatterns []string
}

func BuildTree(rootPath string, opts BuildTreeOpts) (*TreeNode, error) {
	absRoot, err := filepath.Abs(rootPath)
	if err != nil {
		return nil, err
	}

	gi, _ := utils.NewGitIgnoreMatcher(absRoot)

	excludeNames := make(map[string]struct{})
	for k := range utils.DefaultExcludeNames {
		excludeNames[k] = struct{}{}
	}

	root := &TreeNode{
		Name:     filepath.Base(absRoot),
		Path:     ".",
		IsDir:    true,
		State:    Unchecked,
		Expanded: true,
		Depth:    0,
	}

	buildChildren(root, absRoot, absRoot, excludeNames, gi, opts, 1)
	return root, nil
}

func buildChildren(parent *TreeNode, dirPath, rootPath string, excludeNames map[string]struct{}, gi *utils.GitIgnoreMatcher, opts BuildTreeOpts, depth int) {
	if opts.MaxDepth > 0 && depth > opts.MaxDepth {
		return
	}

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return
	}

	sort.Slice(entries, func(i, j int) bool {
		// directories first, then alphabetical
		if entries[i].IsDir() != entries[j].IsDir() {
			return entries[i].IsDir()
		}
		return strings.ToLower(entries[i].Name()) < strings.ToLower(entries[j].Name())
	})

	for _, entry := range entries {
		name := entry.Name()
		fullPath := filepath.Join(dirPath, name)

		if utils.ShouldSkipEntry(fullPath, name, entry.IsDir(), nil, nil, excludeNames, nil, gi) {
			continue
		}

		relPath, _ := filepath.Rel(rootPath, fullPath)
		relPath = filepath.ToSlash(relPath)

		node := &TreeNode{
			Name:   name,
			Path:   relPath,
			IsDir:  entry.IsDir(),
			State:  Unchecked,
			Parent: parent,
			Depth:  depth,
		}

		if entry.IsDir() {
			node.Expanded = false
			buildChildren(node, fullPath, rootPath, excludeNames, gi, opts, depth+1)
		}

		parent.Children = append(parent.Children, node)
	}
}

func FlattenVisible(root *TreeNode) []*TreeNode {
	var result []*TreeNode
	flattenNode(root, &result)
	return result
}

func flattenNode(node *TreeNode, result *[]*TreeNode) {
	*result = append(*result, node)
	if node.IsDir && node.Expanded {
		for _, child := range node.Children {
			flattenNode(child, result)
		}
	}
}

func Toggle(node *TreeNode) {
	switch node.State {
	case Checked:
		setStateRecursive(node, Unchecked)
	default:
		setStateRecursive(node, Checked)
	}
	updateParents(node.Parent)
}

func setStateRecursive(node *TreeNode, state NodeState) {
	node.State = state
	if node.IsDir {
		for _, child := range node.Children {
			setStateRecursive(child, state)
		}
	}
}

func updateParents(node *TreeNode) {
	if node == nil {
		return
	}

	allChecked := true
	allUnchecked := true
	for _, child := range node.Children {
		switch child.State {
		case Checked:
			allUnchecked = false
		case Unchecked:
			allChecked = false
		case Partial:
			allChecked = false
			allUnchecked = false
		}
	}

	if allChecked {
		node.State = Checked
	} else if allUnchecked {
		node.State = Unchecked
	} else {
		node.State = Partial
	}

	updateParents(node.Parent)
}

func ToggleExpand(node *TreeNode) {
	if node.IsDir {
		node.Expanded = !node.Expanded
	}
}

func SetAllState(root *TreeNode, state NodeState) {
	setStateRecursive(root, state)
}

func SetInitialState(root *TreeNode, includePatterns, excludePatterns []string) {
	setInitialStateRecursive(root, includePatterns, excludePatterns)
	updateParents(root)
}

func setInitialStateRecursive(node *TreeNode, includePatterns, excludePatterns []string) {
	if node.IsDir {
		for _, child := range node.Children {
			setInitialStateRecursive(child, includePatterns, excludePatterns)
		}
		updateDirState(node)
		return
	}

	if len(includePatterns) > 0 {
		node.State = Unchecked
		for _, pattern := range includePatterns {
			if matchGlob(pattern, node.Path) {
				node.State = Checked
				break
			}
		}
	} else {
		lang := utils.GetLanguageByExtension(node.Name)
		if lang != "" && !utils.IsTestFile(node.Path, false) && !utils.IsAssetFile(node.Path, false) {
			node.State = Checked
		} else {
			node.State = Unchecked
		}
	}

	if node.State == Checked && len(excludePatterns) > 0 {
		for _, pattern := range excludePatterns {
			if matchGlob(pattern, node.Path) {
				node.State = Unchecked
				break
			}
		}
	}
}

func updateDirState(node *TreeNode) {
	if len(node.Children) == 0 {
		return
	}

	allChecked := true
	allUnchecked := true
	for _, child := range node.Children {
		switch child.State {
		case Checked:
			allUnchecked = false
		case Unchecked:
			allChecked = false
		case Partial:
			allChecked = false
			allUnchecked = false
		}
	}

	if allChecked {
		node.State = Checked
	} else if allUnchecked {
		node.State = Unchecked
	} else {
		node.State = Partial
	}
}

func ApplyConfig(root *TreeNode, cfg *Config) {
	if cfg == nil {
		return
	}
	// Reset all to unchecked first
	setStateRecursive(root, Unchecked)

	// Apply include patterns
	if len(cfg.Include) > 0 {
		applyPatterns(root, cfg.Include, Checked)
	}

	// Apply exclude patterns (override includes)
	if len(cfg.Exclude) > 0 {
		applyPatterns(root, cfg.Exclude, Unchecked)
	}

	// Recalculate parent states
	recalcParents(root)
}

func applyPatterns(node *TreeNode, patterns []string, state NodeState) {
	if !node.IsDir {
		for _, pattern := range patterns {
			matched, _ := filepath.Match(pattern, node.Path)
			if !matched {
				// Try doublestar-style matching
				matched = matchGlob(pattern, node.Path)
			}
			if matched {
				node.State = state
				break
			}
		}
		return
	}
	for _, child := range node.Children {
		applyPatterns(child, patterns, state)
	}
}

func matchGlob(pattern, path string) bool {
	// Handle ** patterns
	if strings.Contains(pattern, "**") {
		parts := strings.SplitN(pattern, "**", 2)
		prefix := parts[0]
		suffix := strings.TrimPrefix(parts[1], "/")

		if prefix != "" && !strings.HasPrefix(path, prefix) {
			return false
		}

		if suffix == "" {
			return true
		}

		// Match suffix against all possible subpaths
		pathParts := strings.Split(path, "/")
		for i := 0; i < len(pathParts); i++ {
			subPath := strings.Join(pathParts[i:], "/")
			if matched, _ := filepath.Match(suffix, subPath); matched {
				return true
			}
			// Also try matching just the filename part
			if matched, _ := filepath.Match(suffix, pathParts[len(pathParts)-1]); matched {
				return true
			}
		}
		return false
	}
	matched, _ := filepath.Match(pattern, path)
	return matched
}

func recalcParents(node *TreeNode) {
	if !node.IsDir {
		return
	}
	for _, child := range node.Children {
		recalcParents(child)
	}
	updateDirState(node)
}

func GeneratePatterns(root *TreeNode) (includes, excludes []string) {
	// Build a default tree state reference to compute diffs
	collectPatterns(root, &includes, &excludes)
	return dedupPatterns(includes), dedupPatterns(excludes)
}

func collectPatterns(node *TreeNode, includes, excludes *[]string) {
	if !node.IsDir {
		if node.State == Checked {
			*includes = append(*includes, node.Path)
		}
		return
	}

	// If all children are same state, express as directory glob
	if node.State == Checked && node.Path != "." {
		*includes = append(*includes, node.Path+"/**")
		return
	}
	if node.State == Unchecked && node.Path != "." {
		return
	}

	for _, child := range node.Children {
		collectPatterns(child, includes, excludes)
	}
}

func dedupPatterns(patterns []string) []string {
	if len(patterns) == 0 {
		return nil
	}
	seen := make(map[string]struct{})
	var result []string
	for _, p := range patterns {
		if _, ok := seen[p]; !ok {
			seen[p] = struct{}{}
			result = append(result, p)
		}
	}
	return result
}

func CountSelected(root *TreeNode) (selected, total int) {
	countFiles(root, &selected, &total)
	return
}

func countFiles(node *TreeNode, selected, total *int) {
	if !node.IsDir {
		*total++
		if node.State == Checked {
			*selected++
		}
		return
	}
	for _, child := range node.Children {
		countFiles(child, selected, total)
	}
}
