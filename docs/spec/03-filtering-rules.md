# 03. Filtering Rules

_Last updated: 2026-05-06_

Filtering is centralized in `utils.ShouldSkipEntry()` for the current CLI paths.

## Priority

The skip decision follows this priority:

1. **Explicit exclusions and default excluded names** - `--exclude` matches and names in `utils.DefaultExcludeNames` always skip.
2. **Include whitelist** - when `--include` is active, matching include paths/patterns are allowed and can override dotfile and `.gitignore` filtering.
3. **`.gitignore`** - ignored files/directories are skipped unless allowed by the include whitelist.
4. **Dotfiles/dot-directories** - names starting with `.` are skipped unless allowed by the include whitelist.
5. **Include-only fallback** - when `--include` is active, non-whitelisted files are skipped.

For pattern-only includes such as `**/*.md`, directories are kept traversable so descendants can still match. Dot directories and `.gitignore`-ignored directories are still skipped while traversing pattern-only includes.

## Default Excluded Names

`utils.DefaultExcludeNames` currently excludes these names:

* Dependency/build/runtime directories: `node_modules`, `vendor`, `target`, `build`, `dist`, `__pycache__`, `env`, `venv`
* Asset directories: `assets`, `static`, `images`, `media`, `uploads`

Dotfiles and dot-directories are handled by a separate default rule rather than this map.

Default excluded names have higher priority than the include whitelist. For example, `--include "assets/logo.png"` does not make the `assets` directory traversable because `assets` itself is skipped by name before include rules are evaluated.

## CLI Include/Exclude Patterns

Both `--include` and `--exclude` accept paths and glob-like patterns.

* Patterns containing glob metacharacters (`*`, `?`, `[]`) are converted to slash-separated matcher patterns.
* Non-recursive glob patterns that do not start with `/` and do not contain `**` are anchored to the scan root. For example, `*.md` becomes `/*.md` and matches only root-level Markdown files.
* Recursive patterns use `**`. For example, `**/*.md` matches Markdown files at any depth.
* Non-glob paths are resolved relative to `--folder` unless absolute, converted back to a root-relative slash path, and anchored with `/`.
* Exclusions apply even inside included directories.

Examples:

```bash
list-codes --exclude "*.md"          # root Markdown files only
list-codes --exclude "**/*.md"       # Markdown files at any depth
list-codes --include "src" --exclude "src/secret"
list-codes --include ".github/**"    # overrides dotfile exclusion in normal collection
```

## `.gitignore`

Unless `--no-gitignore` is set, `utils.NewGitIgnoreMatcher()` loads `.gitignore` files under the scan root and applies them during traversal.

* Matching is done with `github.com/sabhiram/go-gitignore`.
* Paths are matched relative to the scan root using slash separators.
* Directory matches include a trailing slash when the caller knows the path is a directory.
* Nested `.gitignore` files are considered when matching descendants.
* Explicit `--include` patterns can allow files/directories that would otherwise be ignored.

## Test File Detection

Test files are excluded unless `--include-tests` is set. `utils.IsTestFile()` checks:

* Directory markers such as `/test/`, `/tests/`, `/spec/`, `/specs/`, `/__tests__/`, `/fixtures/`, `/mocks/`, `/e2e/`, `/integration/`, `/unit/`
* Filename keywords such as `test`, `spec`, `e2e`, `benchmark`, `bench`, `mock`, `fixture`
* Specific suffixes such as `_test.go`, `.test.js`, `.spec.ts`, `Test.java`, `.test.rs`, `.t.sol`

The same test exclusion applies to the project tree and source collection in normal CLI output.

## Asset Filtering

Asset files are excluded from the project tree and source collection by default. `utils.IsAssetFile()` excludes common binary or non-source extensions including images, fonts, audio, video, archives, office documents, PDFs, and executables.

Explicitly included asset files can appear in normal source output when their parent directories are traversable. If the asset extension is not recognized as a source language, the emitted code fence uses `text`.

Back to [spec index](../spec.md).

