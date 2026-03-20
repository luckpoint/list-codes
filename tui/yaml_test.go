package tui

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadSaveConfig_Roundtrip(t *testing.T) {
	cfg := &Config{
		Include: []string{"src/**/*.go", "cmd/**/*.go"},
		Exclude: []string{"**/*_test.go", "vendor/**"},
		Options: &ConfigOptions{
			IncludeTests: false,
			MaxFileSize:  "1m",
			MaxDepth:     7,
		},
	}

	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "test-config.yaml")

	err := SaveConfig(path, cfg)
	require.NoError(t, err)

	loaded, err := LoadConfig(path)
	require.NoError(t, err)

	assert.Equal(t, cfg.Include, loaded.Include)
	assert.Equal(t, cfg.Exclude, loaded.Exclude)
	assert.Equal(t, cfg.Options.IncludeTests, loaded.Options.IncludeTests)
	assert.Equal(t, cfg.Options.MaxFileSize, loaded.Options.MaxFileSize)
	assert.Equal(t, cfg.Options.MaxDepth, loaded.Options.MaxDepth)
}

func TestLoadConfig_NotFound(t *testing.T) {
	_, err := LoadConfig("/nonexistent/path.yaml")
	assert.Error(t, err)
}

func TestSaveConfig_NoOptions(t *testing.T) {
	cfg := &Config{
		Include: []string{"*.go"},
	}

	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "minimal.yaml")

	err := SaveConfig(path, cfg)
	require.NoError(t, err)

	data, err := os.ReadFile(path)
	require.NoError(t, err)

	content := string(data)
	assert.Contains(t, content, "include:")
	assert.NotContains(t, content, "options:")
}
