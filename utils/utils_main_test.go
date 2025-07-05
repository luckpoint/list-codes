package utils

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

// TestSaveToMarkdown はSaveToMarkdown関数のユニットテストです。
func TestSaveToMarkdown(t *testing.T) {
	tests := []struct {
		name           string
		content        string
		outPath        string
		expectErr      bool
		expectedOutput string
	}{
		{
			name:           "Save to file",
			content:        "# Test Content",
			outPath:        "output.md",
			expectErr:      false,
			expectedOutput: "",
		},
		{
			name:           "Print to stdout",
			content:        "# Stdout Content",
			outPath:        "",
			expectErr:      false,
			expectedOutput: "# Stdout Content\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir, cleanup := setupTestDir(t)
			defer cleanup()

			var actualOutput bytes.Buffer
			if tt.outPath == "" {
				oldStdout := os.Stdout
				r, w, _ := os.Pipe()
				os.Stdout = w

				defer func() {
					os.Stdout = oldStdout
				}()

				err := SaveToMarkdown(tt.content, tt.outPath)
				w.Close()
				_, readErr := actualOutput.ReadFrom(r)
				if readErr != nil {
					t.Fatalf("Failed to read from stdout: %v", readErr)
				}

				if (err != nil) != tt.expectErr {
					t.Errorf("SaveToMarkdown() error mismatch: expected error %v, got %v", tt.expectErr, err)
				}

				if actualOutput.String() != tt.expectedOutput {
					t.Errorf("SaveToMarkdown() stdout mismatch:\nExpected: %q\nActual: %q", tt.expectedOutput, actualOutput.String())
				}
			} else {
				fullPath := filepath.Join(dir, tt.outPath)
				err := SaveToMarkdown(tt.content, fullPath)

				if (err != nil) != tt.expectErr {
					t.Errorf("SaveToMarkdown() error mismatch: expected error %v, got %v", tt.expectErr, err)
				}

				if !tt.expectErr {
					readContent, readErr := ioutil.ReadFile(fullPath)
					if readErr != nil {
						t.Fatalf("Failed to read back file %s: %v", fullPath, readErr)
					}
					if string(readContent) != tt.content {
						t.Errorf("SaveToMarkdown() file content mismatch:\nExpected: %q\nActual: %q", tt.content, string(readContent))
					}
				}
			}
		})
	}
}