# 02. Language Classification and File Collection

_Last updated: 2026-05-06_

## Language Classification

Current source collection is file-based:

1. `utils.GetLanguageByExtension()` classifies each file by extension or special file name.
2. `Dockerfile` is classified as `Dockerfile`.
3. Laravel Blade templates ending in `.blade.php` are classified as `HTML`.
4. Files without a recognized language are skipped unless they are explicitly included asset files; such assets are emitted as `text`.

The current CLI collection path does **not** use project-wide primary/fallback language detection to decide which source files to collect. `utils.PROJECT_SIGNATURES` still exists as a constants table, but `ProcessSourceFiles()` calls `collectSourceFiles()` with nil primary/fallback language arguments and `collectSourceFiles()` intentionally processes all recognized source files.

## Default File Collection

The default summary path runs these stages:

1. `GenerateDirectoryStructure()` builds the Markdown project tree using the same high-level skip rules as source collection.
2. `collectDependencyFiles()` currently returns no dependency/configuration files. `FRAMEWORK_DEPENDENCY_FILES` has been removed.
3. `collectSourceFiles()` walks the folder and collects recognized source files.
4. `buildMarkdownOutput()` combines size diagnostics, project structure, file snippets, and any dependency/configuration snippets.

### Source File Inclusion Criteria

A file is collected when all of the following are true:

* It passes explicit/default path filtering.
* It is not a test file, unless `--include-tests` is set.
* It is not an asset file, unless explicitly included.
* It has a recognized language according to `utils.EXTENSIONS`, `Dockerfile`, or `.blade.php`; explicitly included assets without a language are emitted as `text`.
* Its size is less than or equal to `utils.MaxFileSizeBytes`.
* Adding it does not exceed `utils.TotalMaxFileSizeBytes` when a total limit is configured.
* It was not already collected as a dependency/configuration file. This set is currently empty because dependency collection is a no-op.

## README-only Mode

`CollectReadmeFiles()` walks the folder and emits only files whose name is `README.md` case-insensitively. Directory traversal and README file inclusion still use `ShouldSkipEntry()`, so explicit excludes, include-only mode, default excludes, dotfile rules, and `.gitignore` rules still apply.

If no README files are found, the output is:

```markdown
# Project README Files

No README.md files found in the project.
```

Back to [spec index](../spec.md).

