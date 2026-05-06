# 06. CUI/TUI Selector

_Last updated: 2026-05-06_

The `select` subcommand opens an interactive CUI/TUI for creating a `.list-codes.yaml` selection file:

```bash
list-codes select
list-codes select ./my-config.yaml
```

The default output/load path is `.list-codes.yaml` under `--folder`.

## Tree Construction

`tui.BuildTree()` builds a visible tree before selection state is applied.

Current tree construction:

* Uses `--folder` as the root.
* Respects `--max-depth`.
* Sorts directories before files, then case-insensitively by name.
* Applies `utils.DefaultExcludeNames`, dotfile filtering, and `.gitignore` filtering.
* Does not pass CLI include/exclude matchers into `ShouldSkipEntry()` while building the visible tree.

This means CUI/TUI include/exclude flags are currently **selection filters**, not full visibility filters. For example, `list-codes --include ".github/**" select` can check visible matching files, but it does not make a hidden dot-directory visible during tree construction.

`--no-gitignore` is not passed into `tui.BuildTree()` in the current implementation; the selector always creates a gitignore matcher while building the tree.

## Initial Selection Filters

After the visible tree is built, `tui.SetInitialState()` applies initial checked/unchecked states:

* If `--include` patterns are present, visible files start unchecked and files matching those patterns become checked.
* If no `--include` patterns are present, visible files with recognized source extensions are checked by default.
* Test files and asset files are unchecked by the default initial selection logic.
* `--exclude` patterns are then applied and uncheck matching visible files.
* Existing config at the selected config path is loaded unless `--no-config` is set. Its `include` patterns are applied as checked, then its `exclude` patterns are applied as unchecked.

`BuildTreeOpts.IncludeTests` is currently passed from the CLI but not consulted by `SetInitialState()`. As implemented, test files remain unchecked in the selector's default initial state even when `--include-tests` is supplied; they can still be manually checked if visible.

## Interaction Model

Node states are:

* `Unchecked`
* `Checked`
* `Partial`

Main keybindings:

* `j/k`, `Ctrl+n/Ctrl+p`, arrow up/down: move cursor
* `Enter`, `l`, arrow right: expand/collapse a directory
* `h`, arrow left: collapse directory or move to parent
* `Space`: toggle selected node; toggling a directory recursively toggles children
* `a`: check all visible tree nodes
* `n`: uncheck all visible tree nodes
* `s` or `w`: save config and quit
* `q` or `Ctrl+c`: quit without saving
* `?`: show help

Parent states are recalculated after toggles, so mixed child states render as partial.

## Saved Config

On save, `tui.GeneratePatterns()` converts the checked tree into include patterns:

* A fully checked directory is saved as `path/**`.
* A checked file is saved as its relative path.
* Unchecked directories are omitted.
* Duplicate patterns are removed.

### Output YAML Shape

The selector writes a `tui.Config` value with YAML `omitempty` fields:

```yaml
include:
  - "src/**"
  - "cmd/root.go"
exclude:
  - "**/*.generated.go"
options:
  include-tests: true
  max-file-size: "2m"
  max-depth: 5
```

The schema can contain these top-level keys:

* `include`: list of root-relative slash-separated file paths or glob patterns.
* `exclude`: list of root-relative slash-separated file paths or glob patterns.
* `options`: optional settings accepted by the shared config loader.

However, the current `select` save path constructs the config as:

```go
cfg := &Config{
    Include: includes,
    Exclude: excludes,
}
```

So `select` currently writes only generated `include` and `exclude` fields, and `options` is never emitted by the selector.

### Generated Include Patterns

`GeneratePatterns()` emits include-oriented patterns from the checked tree:

* Checked non-root directories become `path/**`.
* Checked files become `path/to/file.ext`.
* The root node `.` is never emitted as `./**`; if the whole tree is checked, each checked top-level child is emitted.
* Partial directories are traversed and represented by their checked descendants.
* Unchecked directories and files are omitted.
* Paths are relative to the selector root and use `/` separators.
* Patterns are deduplicated while preserving first occurrence order.

Example output:

```yaml
include:
  - src/**
  - cmd/list-codes/main.go
  - README.md
```

### Current Exclude Behavior

The current generator accepts an `excludes` return value, but `collectPatterns()` never appends to it. As implemented, saving from `select` does not generate `exclude` entries from unchecked descendants.

This means:

* An unchecked file inside a partial directory is represented by omitting that file from `include`, not by writing an `exclude`.
* Existing `exclude` entries can affect the initial checked state when a config is loaded.
* After saving, existing `exclude` entries are not preserved unless the generator starts emitting them in the future.
* Existing `options` entries are also not preserved by selector save.
* If no files are checked, all `omitempty` fields are empty and the YAML marshaler writes an empty mapping.

In practice, saved selector configs are include-oriented snapshots of the currently checked visible tree.

Back to [spec index](../spec.md).
