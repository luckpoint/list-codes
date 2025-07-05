package utils

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

// setupTestDir はテスト用の一時ディレクトリを作成し、そのパスを返します。
// クリーンアップ関数も返します。
func setupTestDir(t *testing.T) (string, func()) {
	dir, err := ioutil.TempDir("", "test_project")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	return dir, func() {
		os.RemoveAll(dir)
	}
}

// createTestFile は指定されたパスにファイルを作成し、内容を書き込みます。
func createTestFile(t *testing.T, path, content string) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatalf("Failed to create directories for %s: %v", path, err)
	}
	if err := ioutil.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write file %s: %v", path, err)
	}
}