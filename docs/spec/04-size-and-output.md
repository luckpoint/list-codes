# 04. Size Management and Output Format

_Last updated: 2026-05-06_

## Size Management

### Human-readable Size Formats

`utils.ParseSize()` accepts:

* Bytes: `123`, `123b`
* Kilobytes: `1k`, `1kb`, `500k`
* Megabytes: `1m`, `1mb`, `2.5m`
* Gigabytes: `1g`, `1gb`, `10g`

An empty size string parses as `0`, which means no total limit for `--max-total-size`.

### Limits

* `--max-file-size` controls the individual file limit. Files above this limit are skipped and tracked for debug output.
* `--max-total-size` controls the total collected source size. When adding the next file would exceed the limit, source scanning stops and `limitHit` is set.

### Size Diagnostics

`## Source Code Size Check` is emitted only in `--debug` mode and only when source files or skipped-size messages exist. It includes:

* Total collected size
* Individual file size limit
* Total size limit or `unlimited`
* A sorted list of skipped files that exceeded the individual file limit
* A `scan stopped at ... total limit` message when the total limit was hit

## Output Format

The generated Markdown has this order:

1. **Source Code Size Check** - debug mode only, and placed before the tree.
2. **Project Structure** - a tree in a `text` code fence.
3. **Source file snippets** - one `### relative/path` section per collected source file.
4. **Dependency and Configuration Files** - emitted last only if dependency/configuration snippets exist; currently this section is normally absent because dependency collection is a no-op.

Source snippets are stored by language internally and then emitted with stable sorting:

* Language keys are sorted alphabetically.
* Snippets within each language are sorted lexicographically.
* No visible `### Go` / `### Python` language grouping headings are emitted.
* Each file heading is immediately followed by its code fence with no blank line between them.
* Code fence language hints are lowercased and normalized by removing `/` and replacing `+` with `p`.

The project tree:

* Starts with `## Project Structure`.
* Shows the root as `. (<root-name>)`.
* Uses Unicode tree characters (`├──`, `└──`).
* Sorts entries case-insensitively.
* Displays files before directories at each level.
* Respects `--max-depth`; hidden deeper directories are represented by `...`.
* Applies path, test, and asset filtering.

Back to [spec index](../spec.md).

