package utils

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

// TestPrintDebug はPrintDebug関数のユニットテストです。
func TestPrintDebug(t *testing.T) {
	tests := []struct {
		name          string
		msg           string
		debug         bool
		expectedEmpty bool
	}{
		{name: "Debug mode on", msg: "test message", debug: true, expectedEmpty: false},
		{name: "Debug mode off", msg: "test message", debug: false, expectedEmpty: true},
		{name: "Empty message with debug on", msg: "", debug: true, expectedEmpty: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture stderr
			oldStderr := os.Stderr
			r, w, _ := os.Pipe()
			os.Stderr = w

			PrintDebug(tt.msg, tt.debug)

			w.Close()
			os.Stderr = oldStderr

			var buf bytes.Buffer
			buf.ReadFrom(r)
			output := buf.String()

			if tt.expectedEmpty && output != "" {
				t.Errorf("PrintDebug() with debug=%v should produce no output, got %q", tt.debug, output)
			}
			if !tt.expectedEmpty && output == "" {
				t.Errorf("PrintDebug() with debug=%v should produce output, got empty", tt.debug)
			}
			if !tt.expectedEmpty && !strings.Contains(output, "DEBUG: "+tt.msg) {
				t.Errorf("PrintDebug() output should contain 'DEBUG: %s', got %q", tt.msg, output)
			}
		})
	}
}

// TestPrintWarning はPrintWarning関数のユニットテストです。
func TestPrintWarning(t *testing.T) {
	tests := []struct {
		name          string
		msg           string
		debug         bool
		expectedEmpty bool
	}{
		{name: "Debug mode on", msg: "warning message", debug: true, expectedEmpty: false},
		{name: "Debug mode off", msg: "warning message", debug: false, expectedEmpty: true},
		{name: "Empty warning with debug on", msg: "", debug: true, expectedEmpty: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture stderr
			oldStderr := os.Stderr
			r, w, _ := os.Pipe()
			os.Stderr = w

			PrintWarning(tt.msg, tt.debug)

			w.Close()
			os.Stderr = oldStderr

			var buf bytes.Buffer
			buf.ReadFrom(r)
			output := buf.String()

			if tt.expectedEmpty && output != "" {
				t.Errorf("PrintWarning() with debug=%v should produce no output, got %q", tt.debug, output)
			}
			if !tt.expectedEmpty && output == "" {
				t.Errorf("PrintWarning() with debug=%v should produce output, got empty", tt.debug)
			}
			if !tt.expectedEmpty && !strings.Contains(output, "Warning: "+tt.msg) {
				t.Errorf("PrintWarning() output should contain 'Warning: %s', got %q", tt.msg, output)
			}
		})
	}
}

// TestPrintError はPrintError関数のユニットテストです。
func TestPrintError(t *testing.T) {
	tests := []struct {
		name string
		msg  string
	}{
		{name: "Error message", msg: "error occurred"},
		{name: "Empty error message", msg: ""},
		{name: "Multi-line error", msg: "line1\nline2"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture stderr
			oldStderr := os.Stderr
			r, w, _ := os.Pipe()
			os.Stderr = w

			PrintError(tt.msg)

			w.Close()
			os.Stderr = oldStderr

			var buf bytes.Buffer
			buf.ReadFrom(r)
			output := buf.String()

			if output == "" {
				t.Errorf("PrintError() should always produce output, got empty")
			}
			if !strings.Contains(output, "Error: "+tt.msg) {
				t.Errorf("PrintError() output should contain 'Error: %s', got %q", tt.msg, output)
			}
		})
	}
}