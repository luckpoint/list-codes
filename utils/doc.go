// Package utils provides utility functions for the list-codes CLI tool.
//
// This package contains core functionality for:
//   - Project language detection and classification
//   - File system operations and directory traversal
//   - Source code collection and processing
//   - Markdown output generation
//   - Configuration management
//   - Logging and debugging utilities
//
// The package is organized into several modules:
//   - config.go: Configuration constants and variables
//   - language.go: Language detection and classification
//   - file.go: File system operations and filtering
//   - process.go: Source code processing and Markdown generation
//   - log.go: Logging utilities
//   - utils.go: General utility functions
//
// Key features include:
//   - Automatic programming language detection based on signature files and extensions
//   - Configurable file size limits and filtering
//   - Support for include/exclude path patterns
//   - Test file detection and exclusion
//   - Comprehensive directory structure generation
//   - Dependency file collection and processing
package utils