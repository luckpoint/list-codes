# list-codes - Source Code Analysis Assistant for LLMs

**list-codes** is a command-line tool designed to collect and format source code from a specified project, making it easy to analyze with Large Language Models (LLMs).

It streamlines tasks like code reviews, documentation generation, and bug detection by outputting the source code of an entire project or a specific directory as a single text, which can be copied and pasted into an LLM's prompt.

## Installation

### Option 1: Homebrew (macOS/Linux) - Recommended

```bash
brew tap luckpoint/list-codes
brew install list-codes
```

### Option 2: Go Install

```bash
go install github.com/luckpoint/list-codes/cmd/list-codes@latest
```

**Note:** You need to have `$GOPATH/bin` (or `$HOME/go/bin`) in your system's `PATH`.

### Option 3: Download Binary

Download the latest binary for your platform from the [releases page](https://github.com/luckpoint/list-codes/releases).

### Option 4: Build from Source

```bash
git clone https://github.com/luckpoint/list-codes.git
cd list-codes
go build -o list-codes ./cmd/list-codes
```

## Usage

### Basic Usage

```bash
# Display the source code of the current directory to standard output
list-codes
```

### LLM Integration

The output of this tool is intended to be used directly as input for an LLM.

```bash
# Copy the entire project's source code to the clipboard and paste it into an LLM
list-codes | pbcopy

# Request a code refactor to gemini using predefined template
list-codes --prompt refactor | gemini

# Use custom prompt text
list-codes --prompt "Analyze this code for security vulnerabilities and provide recommendations" --folder ./src

## Prompt Templates and Custom Prompts

The `--prompt` option allows you to prepend specialized prompts to your code output, making it easier to get targeted analysis from LLMs.

### Predefined Templates

**list-codes** includes a comprehensive set of predefined prompt templates for common analysis tasks:

- `explain` - Project overview and architecture explanation
- `find-bugs` - Bug detection and error identification
- `refactor` - Code refactoring suggestions
- `security` - Security vulnerability analysis
- `optimize` - Performance optimization recommendations
- `test` - Testing strategy and test case suggestions
- `document` - Documentation improvement suggestions
- `migrate` - Technology stack migration advice
- `scale` - Scalability analysis and recommendations
- `maintain` - Code maintainability improvements
- `api-design` - API design evaluation and suggestions
- `patterns` - Design pattern application opportunities
- `review` - Comprehensive code review
- `architecture` - Architecture analysis and improvements
- `deploy` - Deployment and operations suggestions

Templates are available in both English and Japanese, automatically selected based on your system locale.

## Filtering and Exclusion Behavior

### Automatic Exclusions

**list-codes** automatically excludes certain files and directories to focus on relevant source code:

#### Dotfiles and Dot-directories (Default)
All files and directories starting with a dot (`.`) are excluded by default, including:
- `.git`, `.gitignore`, `.github/`
- `.vscode/`, `.idea/`
- `.env`, `.DS_Store`
- `.terraform/`, `.serverless/`
- `.pytest_cache/`, `.mypy_cache/`, `.ruff_cache/`

#### Build and Dependency Directories
Common build artifacts and dependency directories are also excluded:
- `node_modules/`, `vendor/`, `target/`
- `build/`, `dist/`, `__pycache__/`
- `env/`, `venv/`

#### Gitignore Integration
The tool automatically respects your project's `.gitignore` rules.

### Override Exclusions with --include

You can override the default exclusions by explicitly including specific files or directories:

```bash
# Include .github directory (normally excluded as a dotfile)
list-codes --include ".github/**"
```

### Exclusion Priority

The filtering logic follows a clear priority system:

1. **Explicit exclusions** (`--exclude`) always win
2. **Include whitelist** (`--include`) overrides default dotfile exclusion
3. **Default dotfile exclusion** applies if not explicitly included
4. **Include-only mode**: When using `--include`, only whitelisted items are processed

### Example Filtering Scenarios

```bash
# Exclude all test files but include .github workflows
list-codes --exclude "**/*test*" --include ".github/workflows/**"

# Include only specific source directories, excluding dotfiles elsewhere
list-codes --include "src/**" --include "lib/**"

# Include documentation from typically excluded directories
list-codes --include ".github/**.md" --include "docs/**"
```

## Output Format

**list-codes** generates a structured Markdown output with the following sections in order:

1. **Source Code Size Check** - File statistics, size limits, and information about skipped files
2. **Project Structure** - Directory tree visualization showing the project layout
3. **Source Code Files** - Organized by programming language with syntax highlighting
4. **Dependency and Configuration Files** - Package files, configs (shown in debug mode)

### File Size Management

The tool provides comprehensive file size control to manage output size and processing time:

#### Human-Readable Size Formats
All size parameters accept human-readable formats:
- **Bytes**: `123`, `123b`
- **Kilobytes**: `1k`, `1kb`, `500k`
- **Megabytes**: `1m`, `1mb`, `2.5m`
- **Gigabytes**: `1g`, `1gb`, `10g`

#### Size Limit Types
- `--max-file-size`: Individual file size limit (default: 1m)
- `--max-total-size`: Total collected files size limit (no limit by default)

#### Size Check Output
The **Source Code Size Check** section displays:
- Total size of collected files
- Current size limits
- List of skipped files (when files exceed limits)
- Whether scanning stopped due to total size limit

```bash
# Examples of size management
list-codes --max-file-size 100k --max-total-size 5m
list-codes --max-file-size 2m --max-total-size 50m
list-codes --max-file-size 500k  # No total limit
```

### Test File Handling

By default, **list-codes** automatically excludes test files to focus on production code. Test files are identified by:

- **File patterns**: `*_test.go`, `*.test.js`, `*.spec.ts`, `Test*.java`
- **Directory patterns**: `/test/`, `/tests/`, `/__tests__/`, `/spec/`
- **Keywords**: Files containing `test`, `spec`, `mock`, `fixture`

To include test files in the analysis:
```bash
list-codes --include-tests
```

### Command-Line Options

#### Core Options
- `--folder`, `-f`: Folder to scan (default: current directory)
- `--output`, `-o`: Output Markdown file path
- `--prompt`, `-p`: Prompt text or template name to prepend to output (accepts both predefined templates and custom text)

#### Filtering Options
- `--include`, `-i`: File/folder path to include, overrides default exclusions (repeatable, supports glob patterns)
- `--exclude`, `-e`: File/folder path to exclude, takes highest priority (repeatable, supports glob patterns)
- `--readme-only`: Only collect README.md files
- `--max-file-size`: Maximum file size to include (supports human-readable formats: 1m, 500k, 2g) (default: 1m)
- `--max-total-size`: Maximum total file size to collect (supports human-readable formats: 10m, 1g) - empty means no limit
- `--max-depth`: Max depth for directory structure (default: 7)
- `--include-tests`: Include test files in the output (excluded by default)

#### Other Options
- `--debug`: Enable debug mode
- `--lang`: Force language (ja|en) instead of auto-detection
- `--version`, `-v`: Show version information
- `--help`, `-h`: Show help message

## Build

To build the executable from the source code, run the following command:

```bash
go build -o list-codes ./cmd/list-codes
```

This will create an executable file named `list-codes` in the project root directory.

## Testing

To run the tests for the entire project, use the following command:

```bash
go test ./...
```

### Test Coverage

To generate a test coverage report, run the following commands:

```bash
# Generate coverage profile
go test -coverprofile=coverage.out ./...

# View the report in your browser
go tool cover -html=coverage.out -o coverage.html
```

## License

This project is released under the MIT License.
