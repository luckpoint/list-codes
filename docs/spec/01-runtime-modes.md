# 01. Runtime Modes and CLI Flags

_Last updated: 2026-05-06_

`list-codes` has three main user-facing modes:

1. **Default summary mode** - scans a folder, emits a project tree, and emits recognized source/configuration files.
2. **README-only mode** - `--readme-only` collects only `README.md` files that pass the normal path filters.
3. **Interactive selector mode** - `list-codes select [config-output-path]` opens the Bubble Tea based CUI/TUI and saves a `.list-codes.yaml` file.

Common flags are registered as persistent flags and are available to the root command and subcommands unless noted:

* `--folder`, `-f`: folder to scan
* `--output`, `-o`: output Markdown file path; stdout when empty
* `--readme-only`: collect README files only
* `--max-depth`: max depth shown in the project tree; default `7`
* `--debug`: print debug/warning diagnostics and include size diagnostics in output
* `--include`, `-i`: include path or glob pattern; repeatable
* `--exclude`, `-e`: exclude path or glob pattern; repeatable
* `--max-file-size`: max individual file size; default `1m`
* `--max-total-size`: max total collected source size; empty means unlimited
* `--prompt`, `-p`: prompt template name or custom prompt text prepended to output
* `--lang`: force help/prompt language (`ja` or `en`)
* `--version`, `-v`: print version
* `--include-tests`: include test files in normal collection
* `--no-gitignore`: disable `.gitignore` filtering
* `--no-config`: disable auto-loading `.list-codes.yaml`

`--config`, `-c` is currently a root-command flag for the default summary path. The `select` subcommand uses its optional positional argument as the config output/load path.

Back to [spec index](../spec.md).

