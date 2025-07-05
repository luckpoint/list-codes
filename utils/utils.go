package utils

import (
	"fmt"
	"os"
)

// SaveToMarkdown saves the generated Markdown content to a file or outputs it to stdout.
func SaveToMarkdown(content string, outputPath string) error {
	if outputPath == "" {
		fmt.Println(content)
		return nil
	}
	return os.WriteFile(outputPath, []byte(content), 0644)
}