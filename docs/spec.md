# list-codes Specification

_Last updated: 2025-07-08_

This document describes how `list-codes` determines what to include in the generated documentation/output.

## 1. Language Detection

1. **Signature scan** – The project root is recursively scanned for files or directories listed in `utils.PROJECT_SIGNATURES`.  
   Each signature maps to one or more programming languages or frameworks.  
   Matching a signature marks the associated language as a **primary language**.
2. **Fallback scan (extension counting)** – If no primary languages were detected, the tool counts file-extension usage based on `utils.EXTENSIONS` and treats the most frequent extensions as **fallback languages**.

## 2. File Collection Logic

The collection stage runs in three passes:

1. `collectDependencyFiles` – Reads dependency / configuration files defined in `utils.FRAMEWORK_DEPENDENCY_FILES` that relate to the detected languages.
2. `collectReadmeFiles` – Reads all `README.md` files.
3. `collectSourceFiles` – Reads actual source code files.

### 2.1 Previous behaviour (≤ v0.10.0)

`collectSourceFiles` processed **only** files that belonged to primary languages (or, if none were detected, to fallback languages).  Files in other languages were skipped by default.  For example, a project containing HTML and JavaScript but without a `package.json` would output only the HTML files.

### 2.2 New default behaviour (≥ v0.11.0)

Starting with this version, `collectSourceFiles` now **includes all recognised source files**, regardless of whether their language is primary, fallback, or neither.  The only exclusions are:

* Files larger than `utils.MaxFileSizeBytes`
* Files flagged by `utils.IsTestFile`
* Files/paths filtered by the include/exclude parameters

This change ensures that mixed-language projects (e.g. HTML + JS) are fully covered without requiring additional configuration.

## 3. Exclusion Rules

* Directories / files in `utils.DefaultExcludeNames` are always skipped.
* Paths explicitly provided via the CLI `--exclude-path` option are skipped.
* Test files are detected via:
  * Directory fragments (`/test/`, `/spec/`, etc.)
  * Filename keywords (`*_test`, `*_spec`, etc.)
  * Custom patterns in `utils.EXCLUDE_TEST_PATTERNS`

## 4. Size Limit

Single files exceeding `utils.MaxFileSizeBytes` (default 1 MiB) are skipped and counted in the “skipped” total.

## 5. Output Format

Collected content is grouped per language under Markdown `###` headings and fenced code blocks annotated with a lowercase language hint (e.g. \`\`\`javascript).  The final document contains:

1. Project structure tree
2. README files
3. Dependency / config files (under “Dependency Files” section)
4. Source files grouped by language

---

> **Note**: Future improvements may introduce an explicit “--languages” CLI option to re-enable language filtering when desired. 