package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	buildCLIOnce sync.Once
	cliBinary    string
	buildCLIErr  error
)

type cliRunResult struct {
	stdout string
	stderr string
	err    error
}

func findRepoRoot(t *testing.T) string {
	t.Helper()

	dir, err := os.Getwd()
	require.NoError(t, err)

	for i := 0; i < 8; i++ {
		if _, statErr := os.Stat(filepath.Join(dir, "go.mod")); statErr == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	t.Fatalf("could not find repository root from %q", dir)
	return ""
}

func buildListCodesCLI(t *testing.T) string {
	t.Helper()

	buildCLIOnce.Do(func() {
		repoRoot := findRepoRoot(t)
		tempDir, err := os.MkdirTemp("", "list-codes-cli-*")
		if err != nil {
			buildCLIErr = fmt.Errorf("failed to create temp dir for CLI build: %w", err)
			return
		}

		name := "list-codes"
		if runtime.GOOS == "windows" {
			name += ".exe"
		}
		cliBinary = filepath.Join(tempDir, name)

		cmd := exec.Command("go", "build", "-o", cliBinary, "./cmd/list-codes")
		cmd.Dir = repoRoot
		output, err := cmd.CombinedOutput()
		if err != nil {
			buildCLIErr = fmt.Errorf("go build failed: %w\n%s", err, string(output))
		}
	})

	require.NoError(t, buildCLIErr)
	return cliBinary
}

func runListCodesCLI(t *testing.T, args ...string) cliRunResult {
	t.Helper()

	cmd := exec.Command(buildListCodesCLI(t), args...)
	cmd.Dir = findRepoRoot(t)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	return cliRunResult{
		stdout: stdout.String(),
		stderr: stderr.String(),
		err:    err,
	}
}

func TestCLI_EndToEndExecutionFlow(t *testing.T) {
	projectDir := t.TempDir()
	outputFile := filepath.Join(t.TempDir(), "result.md")

	require.NoError(t, os.WriteFile(filepath.Join(projectDir, "main.go"), []byte("package main\nfunc main(){}\n"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(projectDir, "main_test.go"), []byte("package main\n"), 0o644))

	result := runListCodesCLI(t,
		"--folder", projectDir,
		"--output", outputFile,
		"--include-tests",
		"--prompt", "test",
		"--lang", "en",
		"--max-file-size", "1m",
		"--max-total-size", "10m",
	)

	require.NoError(t, result.err, "stderr: %s", result.stderr)

	outputBytes, err := os.ReadFile(outputFile)
	require.NoError(t, err)
	output := string(outputBytes)

	assert.Contains(t, output, "Please analyze the following codebase and provide testing improvement suggestions:")
	assert.Contains(t, output, "## Project Structure")
	assert.Contains(t, output, "### main.go")
	assert.Contains(t, output, "### main_test.go")
}

func TestCLI_LanguageFlagAffectsHelpOutput(t *testing.T) {
	ja := runListCodesCLI(t, "--lang", "ja", "--help")
	require.NoError(t, ja.err, "stderr: %s", ja.stderr)
	assert.Contains(t, ja.stdout, "list-codesは、指定されたプロジェクトフォルダをスキャンし")

	en := runListCodesCLI(t, "--lang", "en", "--help")
	require.NoError(t, en.err, "stderr: %s", en.stderr)
	assert.Contains(t, en.stdout, "list-codes is a CLI tool that scans a specified project folder")
}

func TestCLI_LanguageFlagAffectsPromptTemplate(t *testing.T) {
	projectDir := t.TempDir()
	outputFile := filepath.Join(t.TempDir(), "result-ja.md")

	require.NoError(t, os.WriteFile(filepath.Join(projectDir, "main.go"), []byte("package main\nfunc main(){}\n"), 0o644))

	result := runListCodesCLI(t,
		"--folder", projectDir,
		"--output", outputFile,
		"--prompt", "test",
		"--lang", "ja",
	)
	require.NoError(t, result.err, "stderr: %s", result.stderr)

	outputBytes, err := os.ReadFile(outputFile)
	require.NoError(t, err)
	output := string(outputBytes)

	assert.Contains(t, output, "以下のコードベースを分析し、テストの改善提案を行ってください：")
}
