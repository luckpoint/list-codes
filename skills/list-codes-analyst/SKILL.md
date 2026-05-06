---
name: list-codes-analyst
description: Analyzes source code repositories using list-codes and ghq. Use this when the user asks to analyze, review, explain, or refactor a codebase (local or remote).
---

# list-codes-analyst

This skill automates the process of gathering codebase context and analyzing it with an LLM by integrating `ghq` (repository management) and `list-codes` (source code formatting for LLMs).

## Core Workflow

### 1. Repository Identification
- Extract the repository identifier (e.g., `github.com/user/repo` or `user/repo`) from the user's request.
- Use `run_shell_command` with `ghq list --exact --full-path <repo>` to check if the repository is already available locally.

### 2. Handling Missing Repositories
- If the repository is **not** found locally:
    1. **MUST** ask the user for permission: "The repository `<repo>` is not found locally. May I clone it using `ghq get <repo>`?"
    2. If the user approves, execute `run_shell_command("ghq get <repo>")`.
    3. Retrieve the new local path using `ghq list --exact --full-path <repo>`.

### 3. Intent Analysis & Prompt Selection
- Determine the most appropriate `list-codes` prompt template based on the user's goal:
    - **Explain**: Use `--prompt explain` (Overview, architecture, components).
    - **Code Review**: Use `--prompt review` (Comprehensive quality check).
    - **Refactor**: Use `--prompt refactor` (Design patterns and structure).
    - **Security**: Use `--prompt security` (Vulnerability assessment).
    - **Bug Hunting**: Use `--prompt find-bugs` (Logic errors and edge cases).
    - **Custom**: If the user provides specific instructions (e.g., "Check how authentication is implemented"), use a custom string: `--prompt "Check how authentication is implemented"`.

### 4. Execution
- Run the analysis pipeline using `run_shell_command`:
  ```bash
  list-codes -f <repo_path> --prompt <selected_prompt> | gemini
  ```
- **Note**: Replace `gemini` with the current agent's name if applicable (e.g., `claude`, `codex`).

## Trigger Examples
- "Analyze the expressjs/express repository."
- "Explain what's happening in luckpoint/list-codes."
- "Can you do a security review of this remote repo: user/project?"
- "Clone github.com/facebook/react and find bugs in it."
