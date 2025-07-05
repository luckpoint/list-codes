// Package main provides a CLI tool for scanning project folders and generating Markdown summaries.
//
// list-codes is a command-line tool that analyzes project directories and generates comprehensive
// Markdown summaries including:
//   - Project directory structure
//   - Detected programming languages and frameworks
//   - Dependency and configuration files (in debug mode)
//   - Source code files with content
//   - README files
//   - File size statistics and filtering
//
// The tool supports various configuration options including:
//   - Custom file size limits with --max-file-size
//   - Include/exclude path filtering
//   - Debug mode for verbose output
//   - README-only mode
//   - Configurable directory depth limits
//
// Example usage:
//
//	list-codes --folder ./my-project --output summary.md --debug
//	list-codes --readme-only
//	list-codes --exclude node_modules,vendor --max-file-size 2097152
//
// The tool automatically detects project languages based on signature files (like go.mod, package.json)
// and file extensions, then processes relevant source files while excluding test files and
// commonly ignored directories.
package myutilsport