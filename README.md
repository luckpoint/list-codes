![list-codes-logo](https://github.com/user-attachments/assets/2668fa42-2bec-4426-91f0-72c04536aa8f)

# list-codes - Source Code Analysis Assistant for LLMs

**list-codes** is a command-line tool designed to collect and format source code from a specified project, making it easy to analyze with Large Language Models (LLMs).

It streamlines tasks like code reviews, documentation generation, and bug detection by outputting the source code of an entire project or a specific directory as a single text, which can be copied and pasted into an LLM's prompt.

<img width="857" height="591" alt="image" src="https://github.com/user-attachments/assets/475ff692-6acb-4dd2-8103-8e97c24616f2" />

This example shows "list-codes --prompt explain | gemini"  
list-codes gathers code and "--prompt explain" emits a prompt and then passes it to gemini.

When "--prompt explain" is specified, the prompt will be as follows:
```
Please analyze the following codebase and explain the following aspects:

1. **Project purpose and main functionality**
2. **Architecture and design patterns**
3. **Key components and their roles**
4. **Technology stack and dependencies**
5. **Code structure and organization**

Please provide a clear and concise explanation that can be understood by non-technical people.
```

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
- `deps-tree` - Dependency tree generation in Mermaid (external/internal/runtime)
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

1. **Explicit exclusions** (`--exclude`) always win, even inside included folders.
2. **Include whitelist** (`--include`) overrides default dotfile exclusion.
3. **Default dotfile exclusion** applies if not explicitly included.
4. **Gitignore rules** are respected unless overridden by `--include`.

### Glob Patterns and Anchoring

Both `--include` and `--exclude` flags support glob patterns, following a behavior similar to Python's `glob` module:

- `*.md`: Matches `.md` files **only in the project root** (non-recursive).
- `**/*.md`: Matches `.md` files **recursively** in all directories.
- `src/*.ts`: Matches `.ts` files only inside the `src` directory.
- `src/**/*.ts`: Matches `.ts` files recursively within `src` and its subdirectories.

### Example Filtering Scenarios

```bash
# Exclude root Markdown files only
list-codes --exclude "*.md"

# Exclude all Markdown files recursively
list-codes --exclude "**/*.md"

# Include a specific directory but exclude a sub-directory within it
list-codes --include "src" --exclude "src/secret"

# Exclude all test files but include .github workflows
list-codes --exclude "**/*test*" --include ".github/workflows/**"
```

## Output Format

**list-codes** generates a structured Markdown output with the following sections:

1. **Project Structure** - Directory tree visualization showing the project layout
2. **Source Code Files** - Organized by programming language with syntax highlighting
3. **Dependency and Configuration Files** - Package files, configs (shown in debug mode)
4. **Source Code Size Check** - File statistics, size limits, and skipped file information (shown in debug mode)

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
The **Source Code Size Check** section (shown in `--debug` mode) displays:
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

#### Configuration Options
- `--config`, `-c`: Config file path (`.list-codes.yaml`)
- `--no-config`: Disable auto-loading `.list-codes.yaml`

#### Other Options
- `--debug`: Enable debug mode
- `--lang`: Force language (ja|en) instead of auto-detection
- `--version`, `-v`: Show version information
- `--help`, `-h`: Show help message

## Configuration File (`.list-codes.yaml`)

You can create a `.list-codes.yaml` file in your project root to persist include/exclude patterns and options. This file is auto-loaded by default.

```yaml
include:
  - "src/**"
  - "pkg/**"
exclude:
  - "**/*.generated.go"
options:
  include-tests: false
  max-file-size: "1m"
  max-depth: 7
```

CLI flags take priority over config file values. Use `--no-config` to disable auto-loading.

### Interactive File Selector (`select` subcommand)

The `select` subcommand opens a TUI (Terminal User Interface) to interactively browse your project tree and select files. The selection is saved as a `.list-codes.yaml` config file.

```bash
# Open interactive file selector
list-codes select

# Specify output config path
list-codes select ./my-config.yaml
```

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
