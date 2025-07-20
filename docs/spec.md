# list-codes Specification

_Last updated: 2025-07-20_

This document describes how `list-codes` determines what to include in the generated documentation/output and details the current implementation behavior.

## 1. Language Detection

1. **Signature scan** – The project root is recursively scanned for files or directories listed in `utils.PROJECT_SIGNATURES`.  
   Each signature maps to one or more programming languages or frameworks.  
   Matching a signature marks the associated language as a **primary language**.
2. **Fallback scan (extension counting)** – If no primary languages were detected, the tool counts file-extension usage based on `utils.EXTENSIONS` and treats the most frequent extensions as **fallback languages**.

## 2. File Collection Logic

The collection stage runs in two main passes:

1. `collectDependencyFiles` – Reads dependency / configuration files defined in `utils.FRAMEWORK_DEPENDENCY_FILES` that relate to the detected languages (only in debug mode).
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
* Files matching `.gitignore` patterns

### 3.2 CLI-based Exclusions
* Paths explicitly provided via `--exclude` option (highest priority)
* When `--include` is used, only whitelisted paths are processed

### 3.3 Test File Detection
Test files are automatically excluded unless `--include-tests` is specified. Detection uses:
* **Directory patterns**: `/test/`, `/tests/`, `/__tests__/`, `/spec/`, `/specs/`, `/e2e/`, `/integration/`, `/unit/`
* **Filename keywords**: `test`, `spec`, `e2e`, `benchmark`, `bench`, `mock`, `fixture`
* **Specific patterns**: `_test.go`, `.test.js`, `.spec.ts`, `Test.java`, etc. (see `utils.EXCLUDE_TEST_PATTERNS`)

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
The **Source Code Size Check** section displays:
* Total size of collected files
* Current size limits
* List of skipped files with their sizes
* Whether scanning stopped due to total size limit

## 5. Output Format

The generated Markdown output follows a structured format with sections in this order:

### 5.1 Section Order
1. **Source Code Size Check** - File statistics, size limits, skipped files
2. **Project Structure** - Directory tree visualization
3. **Source Code Files** - Grouped by language with syntax highlighting
4. **Dependency and Configuration Files** - Package files, configs (debug mode only)

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
2. **Include whitelist** (`--include`) - Overrides default dotfile exclusion
3. **Default exclusions** - Dotfiles, build directories, test files
4. **Include-only mode** - When `--include` is used, non-whitelisted items are excluded

## 7. Processing Flow

1. **Initialize configuration** - Parse CLI flags, set size limits
2. **Language detection** - Scan for signatures, count extensions
3. **Directory structure** - Generate tree visualization
4. **Dependency collection** - Collect config files (debug mode)
5. **Source collection** - Collect source files with size tracking
6. **Output generation** - Build structured Markdown output

---

> **Implementation Notes**: 
> - The `ParseSize()` function in `utils/utils.go` handles human-readable size parsing
> - File size tracking uses `errTotalSizeLimitExceeded` for total limit enforcement
> - Test detection logic is implemented in `IsTestFile()` with debug logging support 
