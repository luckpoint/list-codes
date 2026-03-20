# list-codes Specification

_Last updated: 2026-03-20_

This document describes how `list-codes` determines what to include in the generated documentation/output and details the current implementation behavior.

## 1. Language Detection

1. **Signature scan** – The project root is recursively scanned for files or directories listed in `utils.PROJECT_SIGNATURES`.  
   Each signature maps to one or more programming languages or frameworks.  
   Matching a signature marks the associated language as a **primary language**.
2. **Fallback scan (extension counting)** – If no primary languages were detected, the tool counts file-extension usage based on `utils.EXTENSIONS` and treats the most frequent extensions as **fallback languages**.

## 2. File Collection Logic

The collection stage runs in two main passes:

1. `collectDependencyFiles` – No longer collects dependency files (FRAMEWORK_DEPENDENCY_FILES has been removed).
2. `collectSourceFiles` – Reads actual source code files from all recognized languages.

### 2.1 Current Implementation

`collectSourceFiles` processes **all recognized source files** regardless of whether their language is primary, fallback, or neither. This ensures comprehensive coverage of mixed-language projects.

**File inclusion criteria:**
* File extension matches any language in `utils.EXTENSIONS`
* File size ≤ `utils.MaxFileSizeBytes` (configurable, default: 1MB)
* Total collected size ≤ `utils.TotalMaxFileSizeBytes` (if set)
* Not flagged as test file by `utils.IsTestFile` (unless `--include-tests` is used)
* Not filtered by include/exclude parameters
* Not in excluded directories (`utils.DefaultExcludeNames`)

## 3. Exclusion Rules

### 3.1 Default Exclusions
* Directories/files in `utils.DefaultExcludeNames` (e.g., `node_modules`, `vendor`, `target`, `build`, `dist`, `__pycache__`)
* Dotfiles and dot-directories (`.git`, `.vscode`, `.idea`, etc.) - can be overridden with `--include`
* Files matching `.gitignore` patterns (**Implemented** - See section 3.4)

### 3.2 CLI-based Exclusions
* **Glob Pattern Support**: Both `--include` and `--exclude` support glob patterns (wildcards).
    * **Non-recursive**: Patterns without recursive wildcards (e.g., `*.md`) are **anchored to the project root**. They match files only in the root directory, not in subdirectories.
    * **Recursive**: Use double asterisks (`**`) for recursive matching (e.g., `**/*.md` matches `.md` files in all directories).
    * **Path Resolution**: All patterns are treated as relative to the project root folder.
* **Exclusion Priority**: Explicit exclusions (`--exclude`) take the highest priority.
    * **Recursive Exclusion within Includes**: Exclusions apply even within directories that have been explicitly included. For example, if you include `src` but exclude `src/secret`, the `secret` directory will be skipped.
* **Include-only mode**: When `--include` is used, only whitelisted paths (and their children, unless excluded) are processed.

### 3.3 Test File Detection
Test files are automatically excluded unless `--include-tests` is specified. Detection uses:
* **Directory patterns**: `/test/`, `/tests/`, `/__tests__/`, `/spec/`, `/specs/`, `/e2e/`, `/integration/`, `/unit/`
* **Filename keywords**: `test`, `spec`, `e2e`, `benchmark`, `bench`, `mock`, `fixture`
* **Specific patterns**: `_test.go`, `.test.js`, `.spec.ts`, `Test.java`, etc. (see `utils.EXCLUDE_TEST_PATTERNS`)

### 3.4 .gitignore Support (**Implemented**)
The tool automatically respects `.gitignore` files found in the project directory tree:
* **Hierarchical processing**: Nested `.gitignore` files are processed with proper precedence (closer rules override parent rules)
* **Full pattern support**: Supports all Git ignore patterns including wildcards (`*.log`), directory patterns (`/build/`), and negation (`!important.log`)
* **User override**: The `--include` flag overrides `.gitignore` rules for explicitly whitelisted paths
* **Disable option**: Use `--no-gitignore` to completely disable `.gitignore` processing for legacy behavior
* **Implementation**: Uses the `github.com/sabhiram/go-gitignore` library for Git-compatible pattern matching

## 4. Size Management

### 4.1 Human-Readable Size Formats
All size parameters accept human-readable formats:
* **Bytes**: `123`, `123b`
* **Kilobytes**: `1k`, `1kb`, `500k`
* **Megabytes**: `1m`, `1mb`, `2.5m`
* **Gigabytes**: `1g`, `1gb`, `10g`

### 4.2 Size Limits
* **Individual file limit**: `--max-file-size` (default: `1m`)
  - Files exceeding this limit are skipped and listed in output
* **Total collection limit**: `--max-total-size` (default: unlimited)
  - Collection stops when total size reaches this limit
  - Displays "scan stopped at X MB total limit" message

### 4.3 Size Reporting
The **Source Code Size Check** section is shown only in debug mode and displays:
* Total size of collected files
* Current size limits
* List of skipped files with their sizes
* Whether scanning stopped due to total size limit

## 5. Output Format

The generated Markdown output follows a structured format with sections in this order:

### 5.1 Section Order
1. **Project Structure** - Directory tree visualization
2. **Source Code Files** - Grouped by language with syntax highlighting
3. **Dependency and Configuration Files** - Package files, configs (debug mode only)
4. **Source Code Size Check** - File statistics, size limits, skipped files (debug mode only)

### 5.2 Content Organization
* **Language grouping**: Source files are organized under `### Language` headings
* **Syntax highlighting**: Code blocks use language-specific highlighting (e.g., `\`\`\`javascript`)
* **File ordering**: Files within each language section are sorted alphabetically
* **Relative paths**: File paths are displayed relative to the project root

### 5.3 Directory Structure
* Uses tree-style visualization with Unicode characters (`├──`, `└──`)
* Respects `--max-depth` setting (default: 7 levels)
* Shows ellipsis (`...`) when directories exceed max depth
* Excludes files/directories based on filtering rules

## 6. CLI Flag Priority

The filtering system follows a clear priority hierarchy:

1. **Explicit exclusions** (`--exclude`) - Highest priority, always excludes
2. **Include whitelist** (`--include`) - Overrides default dotfile exclusion and .gitignore rules
3. **.gitignore patterns** - Files/directories matching .gitignore rules are excluded (unless disabled with `--no-gitignore`)
4. **Default exclusions** - Dotfiles, build directories, test files
5. **Include-only mode** - When `--include` is used, non-whitelisted items are excluded

## 7. Processing Flow

1. **Initialize configuration** - Parse CLI flags, set size limits
2. **Load config file** - Auto-load `.list-codes.yaml` (unless `--no-config` is set)
3. **Load .gitignore rules** - Parse .gitignore files in project hierarchy (unless `--no-gitignore` is set)
4. **Language detection** - Scan for signatures, count extensions
5. **Directory structure** - Generate tree visualization
6. **Dependency collection** - Collect config files (debug mode)
7. **Source collection** - Collect source files with size tracking
8. **Output generation** - Build structured Markdown output

---

## 8. Configuration File (`.list-codes.yaml`)

### 8.1 Auto-loading
* When no `--config` flag is provided and `--no-config` is not set, the tool looks for `.list-codes.yaml` in the scan folder
* If found, it is automatically loaded before CLI flags are processed

### 8.2 Config Structure
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

### 8.3 Priority
* Config file patterns are prepended to CLI flag values (CLI flags take priority)
* Config `options` are applied only when the corresponding CLI flag is not explicitly set

### 8.4 Related CLI Flags
* `--config`, `-c`: Explicit config file path
* `--no-config`: Disable auto-loading

## 9. Interactive File Selector (`select` subcommand)

The `select` subcommand provides a TUI (bubbletea-based) for interactively browsing the project tree and toggling file selection. The result is saved as a `.list-codes.yaml` config file.

### 9.1 Usage
```bash
list-codes select [config-output-path]
```
* Default output path: `.list-codes.yaml` in the scan folder
* Respects `--folder`, `--include`, `--exclude`, `--include-tests`, `--max-depth` flags

### 9.2 Tree Building
* Uses `tui.BuildTree()` which calls `utils.ShouldSkipEntry()` for filtering
* Sorts entries: directories first, then alphabetical (case-insensitive)
* Respects `.gitignore` and default exclusions

### 9.3 State Management
* Three states: `Unchecked`, `Checked`, `Partial`
* Toggling a directory recursively toggles all children
* Parent state is recalculated automatically (all checked / all unchecked / partial)

---

> **Implementation Notes**:
> - The `ParseSize()` function in `utils/utils.go` handles human-readable size parsing
> - File size tracking uses `errTotalSizeLimitExceeded` for total limit enforcement
> - Test detection logic is implemented in `IsTestFile()` with debug logging support
> - `ShouldSkipEntry()` is exported for reuse in the `tui` package
