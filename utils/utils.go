package utils

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// ParseSize parses human-readable size strings like "1k", "1m", "1g" and returns bytes.
// Supports formats: 123, 123b, 123k, 123m, 123g (case insensitive)
func ParseSize(sizeStr string) (int64, error) {
	if sizeStr == "" {
		return 0, nil
	}

	// Regular expression to match number followed by optional unit
	re := regexp.MustCompile(`^(\d+(?:\.\d+)?)\s*([kmgKMG]?[bB]?)$`)
	matches := re.FindStringSubmatch(strings.TrimSpace(sizeStr))
	
	if len(matches) != 3 {
		return 0, fmt.Errorf("invalid size format: %s (expected formats: 123, 123b, 123k, 123m, 123g)", sizeStr)
	}

	// Parse the numeric part
	value, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return 0, fmt.Errorf("invalid number in size: %s", matches[1])
	}

	// Parse the unit part
	unit := strings.ToLower(matches[2])
	var multiplier int64 = 1

	switch unit {
	case "", "b":
		multiplier = 1
	case "k", "kb":
		multiplier = 1024
	case "m", "mb":
		multiplier = 1024 * 1024
	case "g", "gb":
		multiplier = 1024 * 1024 * 1024
	default:
		return 0, fmt.Errorf("unsupported size unit: %s (supported: b, k, m, g)", matches[2])
	}

	result := int64(value * float64(multiplier))
	if result < 0 {
		return 0, fmt.Errorf("size cannot be negative: %s", sizeStr)
	}

	return result, nil
}

// SaveToMarkdown saves the generated Markdown content to a file or outputs it to stdout.
func SaveToMarkdown(content string, outputPath string) error {
	if outputPath == "" {
		fmt.Println(content)
		return nil
	}
	return os.WriteFile(outputPath, []byte(content), 0644)
}