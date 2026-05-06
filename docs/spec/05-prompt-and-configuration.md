# 05. Prompt and Configuration

_Last updated: 2026-05-06_

## Prompt and Language Handling

If `--prompt` is set, `utils.GetPrompt()` resolves it through the i18n prompt templates:

* A known template name returns the localized template.
* Any unknown value is treated as custom prompt text.
* The prompt is prepended to the generated Markdown with one blank line between prompt and content.

`--lang` is parsed early before Cobra help text is generated. Supported values are `ja`, `japanese`, `en`, and `english`; the public flag documents `ja|en`.

If `--lang` is not provided, `utils.InitI18n()` uses the system locale and falls back to English.

## Configuration File (`.list-codes.yaml`)

### Auto-loading

In default summary mode, when `--config` is empty and `--no-config` is not set, the tool looks for `.list-codes.yaml` inside `--folder`. If found, it loads that file before building matchers and scanning files.

### Structure

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

Supported `options` fields are:

* `include-tests`
* `max-file-size`
* `max-depth`

`max-total-size`, `prompt`, `output`, `lang`, `debug`, `readme-only`, and `no-gitignore` are not config-file options in the current implementation.

### Merge Rules

* Config `include` patterns are prepended to CLI `--include` values.
* Config `exclude` patterns are prepended to CLI `--exclude` values.
* Config `options` are applied only if the corresponding CLI flag was not explicitly changed.
* `--no-config` disables auto-loading.

Back to [spec index](../spec.md).

