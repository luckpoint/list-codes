# list-codes - Source Code Analysis Assistant for LLMs

**list-codes** is a command-line tool designed to collect and format source code from a specified project, making it easy to analyze with Large Language Models (LLMs).

It streamlines tasks like code reviews, documentation generation, and bug detection by outputting the source code of an entire project or a specific directory as a single text, which can be copied and pasted into an LLM's prompt.

## Key Features

- **📁 Source Code Collection:** Recursively collects source code from the specified directory. It automatically respects `.gitignore` to exclude unnecessary files.
- **🔧 Flexible Filtering:** Use `--exclude` and `--include` options to omit or target specific files and directories.
- **🤖 LLM-Ready Prompts:** The `--prompt` option allows you to prepend one of over 15 predefined prompts (e.g., `explain`, `find-bugs`) to make your instructions to the LLM simple and precise.

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
go build -o list-codes cmd/list-codes/main.go
```

## Usage

### Basic Usage

```bash
# Display the source code of the current directory to standard output
list-codes

# Specify a project folder
list-codes --folder /path/to/your/project

# Target only the `src` and `pkg` directories
list-codes --include "src/**,pkg/**"

# Exclude `*.test.go` files
list-codes --exclude "**/*.test.go"

# Add an 'explain' prompt and output to a file
list-codes --prompt explain --output for_llm.txt

# Enable debug mode and force language
list-codes --debug --lang en

# Show version information
list-codes --version
```

### LLM Integration

The output of this tool is intended to be used directly as input for an LLM.

```bash
# Copy the entire project's source code to the clipboard and paste it into an LLM
list-codes --folder /path/to/your/project | pbcopy

# Request a code review for a specific feature
list-codes --prompt refactor --folder ./src/feature | llm-cli

# Generate a project overview
list-codes --prompt explain --folder . > project_overview.txt
# (Then, pass the content of project_overview.txt to the LLM)
```

### Command-Line Options

#### Core Options
- `--folder`, `-f`: Folder to scan (default: current directory)
- `--output`, `-o`: Output Markdown file path
- `--prompt`, `-p`: Prompt template name to prepend to output

#### Filtering Options
- `--include`, `-i`: Folder path to include (repeatable)
- `--exclude`, `-e`: Folder path to exclude (repeatable)
- `--readme-only`: Only collect README.md files
- `--max-file-size`: Maximum file size in bytes (default: 1MB)
- `--max-depth`: Max depth for directory structure (default: 7)

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
