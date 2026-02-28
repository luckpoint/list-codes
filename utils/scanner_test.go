package utils

import (
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newProjectScannerForTest(root string) *projectScanner {
	return &projectScanner{
		rootPath: root,

		includePaths: map[string]struct{}{},
		excludeNames: map[string]struct{}{},
		excludePaths: map[string]struct{}{},

		processedDepFiles: map[string]struct{}{},

		collectStructure: true,
		collectReadmes:   true,
		collectSources:   true,
	}
}

func flattenSourceMarkdownByLanguage(sourceFileContents map[string][]string) string {
	var all []string
	for _, contents := range sourceFileContents {
		all = append(all, contents...)
	}
	sort.Strings(all)
	return strings.Join(all, "\n")
}

func countSourceEntries(sourceFileContents map[string][]string) int {
	total := 0
	for _, contents := range sourceFileContents {
		total += len(contents)
	}
	return total
}

func TestProjectScannerScan_BasicCollection(t *testing.T) {
	tempDir := t.TempDir()

	createTestFile(t, filepath.Join(tempDir, "README.md"), "# Root README")
	createTestFile(t, filepath.Join(tempDir, "src", "main.go"), "package main\nfunc main() {}\n")
	createTestFile(t, filepath.Join(tempDir, "src", "main_test.go"), "package main\n")
	createTestFile(t, filepath.Join(tempDir, "assets", "logo.png"), "not-a-real-png")

	scanner := newProjectScannerForTest(tempDir)
	scanner.includeTests = false

	result, err := scanner.scan()
	require.NoError(t, err)
	require.NotNil(t, result)

	require.Len(t, result.readmeFiles, 1)
	assert.Contains(t, result.readmeFiles[0], "### README.md\n```markdown\n# Root README")

	sourceMarkdown := flattenSourceMarkdownByLanguage(result.sourceFileContents)
	assert.Contains(t, sourceMarkdown, "### src/main.go")
	assert.NotContains(t, sourceMarkdown, "main_test.go")
	assert.NotContains(t, sourceMarkdown, "logo.png")

	structureMD := buildDirectoryStructureMarkdown(result.rootPath, 10, result.structureChildren)
	assert.Contains(t, structureMD, "README.md")
	assert.Contains(t, structureMD, "src")
	assert.NotContains(t, structureMD, "main_test.go")
	assert.NotContains(t, structureMD, "logo.png")
}

func TestProjectScannerScan_ExplicitIncludeCollectsAssetAsText(t *testing.T) {
	tempDir := t.TempDir()
	assetPath := filepath.Join(tempDir, "assets", "logo.png")
	createTestFile(t, assetPath, "asset-bytes")

	scanner := newProjectScannerForTest(tempDir)
	scanner.includePaths = map[string]struct{}{
		assetPath: {},
	}
	scanner.collectStructure = false
	scanner.collectReadmes = false
	scanner.collectSources = true

	result, err := scanner.scan()
	require.NoError(t, err)

	sourceMarkdown := flattenSourceMarkdownByLanguage(result.sourceFileContents)
	assert.Contains(t, sourceMarkdown, "### assets/logo.png")
	assert.Contains(t, sourceMarkdown, "```text")
}

func TestProjectScannerScan_SizeLimitBoundaries(t *testing.T) {
	origMaxFileSize := MaxFileSizeBytes
	origTotalMaxSize := TotalMaxFileSizeBytes
	t.Cleanup(func() {
		MaxFileSizeBytes = origMaxFileSize
		TotalMaxFileSizeBytes = origTotalMaxSize
	})

	t.Run("max file size includes exact boundary", func(t *testing.T) {
		tempDir := t.TempDir()
		createTestFile(t, filepath.Join(tempDir, "a.go"), "a")
		createTestFile(t, filepath.Join(tempDir, "bb.go"), "bb")

		MaxFileSizeBytes = 1
		TotalMaxFileSizeBytes = 0

		scanner := newProjectScannerForTest(tempDir)
		scanner.collectStructure = false
		scanner.collectReadmes = false
		scanner.collectSources = true

		result, err := scanner.scan()
		require.NoError(t, err)

		sourceMarkdown := flattenSourceMarkdownByLanguage(result.sourceFileContents)
		assert.Contains(t, sourceMarkdown, "### a.go")
		assert.NotContains(t, sourceMarkdown, "bb.go")
		require.Len(t, result.skippedFileMessages, 1)
		assert.Contains(t, result.skippedFileMessages[0], "bb.go")
	})

	t.Run("max total size includes exact boundary", func(t *testing.T) {
		tempDir := t.TempDir()
		createTestFile(t, filepath.Join(tempDir, "a.go"), "a")
		createTestFile(t, filepath.Join(tempDir, "bb.go"), "bb")
		createTestFile(t, filepath.Join(tempDir, "ccc.go"), "ccc")

		MaxFileSizeBytes = 1024
		TotalMaxFileSizeBytes = 3 // a + bb exactly fits, ccc exceeds

		scanner := newProjectScannerForTest(tempDir)
		scanner.collectStructure = false
		scanner.collectReadmes = false
		scanner.collectSources = true

		result, err := scanner.scan()
		require.NoError(t, err)

		sourceMarkdown := flattenSourceMarkdownByLanguage(result.sourceFileContents)
		assert.Contains(t, sourceMarkdown, "### a.go")
		assert.Contains(t, sourceMarkdown, "### bb.go")
		assert.NotContains(t, sourceMarkdown, "ccc.go")
		assert.True(t, result.limitHit)
		assert.EqualValues(t, 3, result.totalFileSize)
	})

	t.Run("zero max file size only allows zero byte files", func(t *testing.T) {
		tempDir := t.TempDir()
		createTestFile(t, filepath.Join(tempDir, "empty.go"), "")
		createTestFile(t, filepath.Join(tempDir, "nonempty.go"), "x")

		MaxFileSizeBytes = 0
		TotalMaxFileSizeBytes = 0

		scanner := newProjectScannerForTest(tempDir)
		scanner.collectStructure = false
		scanner.collectReadmes = false
		scanner.collectSources = true

		result, err := scanner.scan()
		require.NoError(t, err)

		sourceMarkdown := flattenSourceMarkdownByLanguage(result.sourceFileContents)
		assert.Contains(t, sourceMarkdown, "### empty.go")
		assert.NotContains(t, sourceMarkdown, "nonempty.go")
		require.Len(t, result.skippedFileMessages, 1)
		assert.Contains(t, result.skippedFileMessages[0], "nonempty.go")
	})
}

func TestProjectScannerScan_SymlinkLoopDoesNotRecurseInfinitely(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlink behavior differs on Windows")
	}

	tempDir := t.TempDir()
	createTestFile(t, filepath.Join(tempDir, "src", "main.go"), "package main")

	loopPath := filepath.Join(tempDir, "src", "loop")
	if err := os.Symlink(tempDir, loopPath); err != nil {
		t.Skipf("symlink creation is not supported in this environment: %v", err)
	}

	scanner := newProjectScannerForTest(tempDir)
	result, err := scanner.scan()
	require.NoError(t, err)

	structureMD := buildDirectoryStructureMarkdown(result.rootPath, 10, result.structureChildren)
	assert.Contains(t, structureMD, "main.go")
	assert.LessOrEqual(t, strings.Count(structureMD, "loop"), 1)
}

func TestProjectScannerScan_PermissionDeniedDirectoryIsSkipped(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX permission bits are required for this test")
	}

	tempDir := t.TempDir()
	createTestFile(t, filepath.Join(tempDir, "public.go"), "package main")
	secretDir := filepath.Join(tempDir, "secret")
	createTestFile(t, filepath.Join(secretDir, "hidden.go"), "package main")

	require.NoError(t, os.Chmod(secretDir, 0o000))
	t.Cleanup(func() {
		_ = os.Chmod(secretDir, 0o755)
	})

	if _, err := os.ReadDir(secretDir); err == nil {
		t.Skip("filesystem does not enforce chmod-based read restrictions in this environment")
	}

	scanner := newProjectScannerForTest(tempDir)
	scanner.debug = true

	result, err := scanner.scan()
	require.NoError(t, err)

	sourceMarkdown := flattenSourceMarkdownByLanguage(result.sourceFileContents)
	assert.Contains(t, sourceMarkdown, "public.go")
	assert.NotContains(t, sourceMarkdown, "hidden.go")
}

func TestProjectScannerScan_UnreadableFileIsSkipped(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("POSIX permission bits are required for this test")
	}

	tempDir := t.TempDir()
	readablePath := filepath.Join(tempDir, "readable.go")
	unreadablePath := filepath.Join(tempDir, "unreadable.go")
	createTestFile(t, readablePath, "package main")
	createTestFile(t, unreadablePath, "package main")

	require.NoError(t, os.Chmod(unreadablePath, 0o000))
	t.Cleanup(func() {
		_ = os.Chmod(unreadablePath, 0o644)
	})

	if _, err := os.ReadFile(unreadablePath); err == nil {
		t.Skip("filesystem does not enforce chmod-based read restrictions in this environment")
	}

	scanner := newProjectScannerForTest(tempDir)
	scanner.debug = true

	result, err := scanner.scan()
	require.NoError(t, err)

	sourceMarkdown := flattenSourceMarkdownByLanguage(result.sourceFileContents)
	assert.Contains(t, sourceMarkdown, "readable.go")
	assert.NotContains(t, sourceMarkdown, "unreadable.go")
}

func TestProjectScannerScan_SpecialFileNames(t *testing.T) {
	tempDir := t.TempDir()

	createTestFile(t, filepath.Join(tempDir, "hello world.go"), "package main")
	longPath := filepath.Join(
		tempDir,
		"very-long-segment-aaaaaaaaaaaaaaaaaaaaaaaa",
		"very-long-segment-bbbbbbbbbbbbbbbbbbbbbbbb",
		"very-long-segment-cccccccccccccccccccccccc",
		"deeply_nested.go",
	)
	createTestFile(t, longPath, "package deep")

	expectedMinEntries := 2
	nonUTF8Name := string([]byte{0xff, 0xfe}) + ".go"
	nonUTF8Path := filepath.Join(tempDir, nonUTF8Name)
	if err := os.WriteFile(nonUTF8Path, []byte("package odd"), 0o644); err == nil {
		expectedMinEntries = 3
	}

	scanner := newProjectScannerForTest(tempDir)
	scanner.collectReadmes = false

	result, err := scanner.scan()
	require.NoError(t, err)

	sourceMarkdown := flattenSourceMarkdownByLanguage(result.sourceFileContents)
	assert.Contains(t, sourceMarkdown, "hello world.go")
	assert.Contains(t, sourceMarkdown, "deeply_nested.go")
	assert.GreaterOrEqual(t, countSourceEntries(result.sourceFileContents), expectedMinEntries)
}
