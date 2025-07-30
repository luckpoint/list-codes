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

func TestSaveToMarkdownErrorPath(t *testing.T) {
	// Test error path when trying to write to a read-only directory
	t.Run("Write to read-only directory", func(t *testing.T) {
		// Create a temporary directory and make it read-only
		tempDir := t.TempDir()
		readOnlyDir := filepath.Join(tempDir, "readonly")
		err := os.Mkdir(readOnlyDir, 0555) // Read and execute only, no write
		if err != nil {
			t.Fatalf("Failed to create read-only directory: %v", err)
		}
		
		// Try to write to a file in the read-only directory
		testFile := filepath.Join(readOnlyDir, "test.md")
		content := "# Test Content"
		
		err = SaveToMarkdown(content, testFile)
		if err == nil {
			t.Errorf("Expected SaveToMarkdown to return an error when writing to read-only directory, but it did not")
		}
		
		// Verify the file was not created
		if _, statErr := os.Stat(testFile); statErr == nil {
			t.Errorf("Expected file not to be created in read-only directory, but it was")
		}
	})
	
	t.Run("Write to invalid path", func(t *testing.T) {
		// Try to write to a path that contains non-existent parent directories
		// Use a path that's very unlikely to exist and can't be created
		invalidPath := "/this/path/absolutely/does/not/exist/and/cannot/be/created/test.md"
		content := "# Test Content"
		
		err := SaveToMarkdown(content, invalidPath)
		if err == nil {
			t.Errorf("Expected SaveToMarkdown to return an error for invalid path, but it did not")
		}
	})
}