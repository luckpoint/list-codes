package utils

import (
	"testing"
)

func TestParseSize(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
		hasError bool
	}{
		// Basic numbers
		{"123", 123, false},
		{"0", 0, false},
		{"", 0, false},

		// Bytes
		{"123b", 123, false},
		{"123B", 123, false},

		// Kilobytes
		{"1k", 1024, false},
		{"1K", 1024, false},
		{"1kb", 1024, false},
		{"1KB", 1024, false},
		{"2k", 2048, false},
		{"0.5k", 512, false},

		// Megabytes
		{"1m", 1024*1024, false},
		{"1M", 1024*1024, false},
		{"1mb", 1024*1024, false},
		{"1MB", 1024*1024, false},
		{"2m", 2*1024*1024, false},
		{"1.5m", int64(1.5*1024*1024), false},

		// Gigabytes
		{"1g", 1024*1024*1024, false},
		{"1G", 1024*1024*1024, false},
		{"1gb", 1024*1024*1024, false},
		{"1GB", 1024*1024*1024, false},
		{"2g", 2*1024*1024*1024, false},

		// Whitespace handling
		{" 1m ", 1024*1024, false},
		{"1 m", 1024*1024, false},

		// Error cases
		{"abc", 0, true},
		{"1x", 0, true},
		{"1tb", 0, true},
		{"-1m", 0, true},
		{"1.2.3m", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := ParseSize(tt.input)
			
			if tt.hasError {
				if err == nil {
					t.Errorf("ParseSize(%q) expected error, but got none", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("ParseSize(%q) unexpected error: %v", tt.input, err)
				}
				if result != tt.expected {
					t.Errorf("ParseSize(%q) = %d, expected %d", tt.input, result, tt.expected)
				}
			}
		})
	}
}