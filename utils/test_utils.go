package utils

import (
	"os"
	"path/filepath"
	"testing"
)

// setupTestDir はテスト用の一時ディレクトリを作成し、そのパスを返します。
// クリーンアップ関数も返します。
func setupTestDir(t *testing.T) (string, func()) {
	t.Helper()
	dir := t.TempDir()
	return dir, func() {}
}

// createTestFile は指定されたパスにファイルを作成し、内容を書き込みます。
func createTestFile(t *testing.T, path, content string) {
	t.Helper()
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("Failed to create directories for %s: %v", path, err)
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write file %s: %v", path, err)
	}
}
