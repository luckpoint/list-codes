# 07. Processing Flow and Implementation Notes

_Last updated: 2026-05-06_

## Default Summary Mode

1. Parse early `--lang` and initialize i18n.
2. Handle `--version`.
3. Auto-detect and load `.list-codes.yaml` unless disabled.
4. Merge config include/exclude patterns and options with CLI flags.
5. Parse size strings.
6. Build default excluded-name map.
7. Resolve `--folder` to an absolute path.
8. Normalize include/exclude patterns and build `SimpleMatcher` instances.
9. Build `.gitignore` matcher unless `--no-gitignore` is set.
10. Run README-only or default source processing.
11. Apply prompt text if provided.
12. Write to `--output` or stdout.

## Selector Mode

1. Resolve output config path from `[config-output-path]` or default `.list-codes.yaml`.
2. Build `BuildTreeOpts` from `--include-tests`, `--max-depth`, `--exclude`, and `--include`.
3. Build the visible tree.
4. Apply initial selection from CLI include/exclude patterns.
5. Load existing config from the output path unless `--no-config` is set, then apply it.
6. Run the Bubble Tea TUI.
7. Save generated include patterns when `s` or `w` is pressed.

## Implementation Notes

* `utils.ParseSize()` implements human-readable size parsing.
* `utils.ShouldSkipEntry()` is reused by the CLI tree, README collection, source collection, and the TUI tree builder.
* `utils.IsTestFile()` and `utils.IsAssetFile()` are the current content-type exclusion helpers.
* `collectDependencyFiles()` is intentionally a no-op in the current implementation.
* `utils/scanner.go` contains an alternate scanner implementation, but `ProcessSourceFiles()` currently uses the functions in `utils/file.go`.

Back to [spec index](../spec.md).
