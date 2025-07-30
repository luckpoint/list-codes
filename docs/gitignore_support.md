# .gitignore Support – Design & Implementation Plan

_Last updated: 2025-07-30_

## 1. Motivation

`list-codes` already honours dot-files, test files, CLI include/exclude flags, and several hard-coded directories.  However, many projects rely on a `.gitignore` to express *project-specific* exclusions (generated artefacts, local builds, editor caches, etc.).  Failing to respect these rules can:

* Increase noise in the generated output
* Waste time/size budget on irrelevant files
* Surprise users who expect tooling to follow the same visibility model as Git

Therefore the tool should **exclude any file or directory matched by the effective `.gitignore` rules**, unless the user explicitly overrides the behaviour.

---

## 2. Functional Requirements

1. **Default behaviour** – Paths ignored by Git must be skipped during all collection phases (directory tree, README scan, dependency scan, source scan).
2. **Nested `.gitignore` files** – Rules are hierarchical; each directory’s rules extend the parent’s.
3. **Negation (`!`) patterns** – Files re-enabled by `!pattern` should not be excluded.
4. **User overrides** –
   * `--exclude` ⟹ still ignored (highest priority)
   * `--include` ⟹ always processed even if `.gitignore` says to ignore
   * New flag `--no-gitignore` disables this feature entirely.
5. **Performance** – Parsing overhead should be <5 ms for typical repos; matching must remain < O(pathLen) per lookup.
6. **Deterministic output** – Same input tree & flags must yield identical results across runs.

---

## 3. High-Level Design

```mermaid
flowchart TD
    Init[CLI parse]
    Init -->|root path| LoadIg[Load .gitignore hierarchy]
    LoadIg --> Matcher[GitIgnoreMatcher]
    Matcher -->|shouldSkipDir?|
        DirTree[GenerateDirectoryStructure]
    Matcher --> SrcScan[collectSourceFiles]
    Matcher --> DepScan[collectDependencyFiles]
    Matcher --> ReadmeScan[CollectReadmeFiles]
```

### 3.1 Library Choice

We will leverage the well-maintained package `github.com/sabhiram/go-gitignore` which already implements:

* Parsing of Git’s glob syntax (including `**`, `!`, trailing `/`, etc.)
* Logical combination of multiple `.gitignore` files
* Efficient path matching via compiled patterns

### 3.2 New Component – `GitIgnoreMatcher`

```go
// utils/gitignore.go
package utils

type GitIgnoreMatcher struct {
    root string
    tree map[string]*ignore.GitIgnore // key = directory absolute path
}

func NewGitIgnoreMatcher(root string) (*GitIgnoreMatcher, error)
func (m *GitIgnoreMatcher) Match(path string) bool // true ➜ ignore
```

Creation walks the repo once, loading every `.gitignore` encountered (depth-first).  `Match` climbs up from `path` to `root`, applying the closest matcher first, mirroring Git’s precedence rules.

### 3.3 API Changes

1. `shouldSkipDir` gains an extra parameter `git *GitIgnoreMatcher`.
2. Public helpers (`GenerateDirectoryStructure`, `collect*`) also accept the matcher.
3. CLI: `--no-gitignore` (bool, default `false`).  When set, the matcher is `nil` and no extra checks run.

### 3.4 Skip Decision Priority (revised)

1. `--exclude` explicit paths/names
2. `--include` whitelist (overrides 3-5)
3. **.gitignore** matcher (new)
4. Dot-file default exclusion
5. General defaults / test-file detection

---

## 4. Implementation Steps

1. **Dependency** – Add in `go.mod`:
   ```
   require github.com/sabhiram/go-gitignore v1.1.0
   ```
2. **New file** `utils/gitignore.go` (see 3.2).
3. **Refactor** `shouldSkipDir` signature:
   ```go
   func shouldSkipDir(fullPath, name string, isDir bool, includePaths, excludeNames, excludePaths map[string]struct{}, gi *GitIgnoreMatcher) bool
   ```
4. **Pass matcher** from `cmd/list-codes/main.go` where configuration is built.
5. **CLI flag** – Add `--no-gitignore` to `cmd/list-codes/lang_wrapper.go` (global flags struct).
6. **Unit tests** –
   * `utils/gitignore_test.go` covering
     * simple ignore pattern
     * negation pattern
     * nested `.gitignore`
     * interaction with `--include`
7. **Docs update** –
   * Update `docs/spec.md` section *3. Exclusion Rules* to mark feature as **Implemented**
   * Reference this document for details.
8. **Benchmark (optional)** – confirm negligible overhead on large trees.

---

## 5. Roll-out & Backwards Compatibility

* The change is **opt-in false** (automatic) but can be disabled via `--no-gitignore` for legacy behaviour.
* No breaking API changes for library consumers – only internal helpers gain parameters.
* The default output should shrink in size for repositories with generated artefacts; CI baselines may need updating.

---

## 6. Open Questions / Future Work

1. Should we respect `git config core.excludesFile` (global ignore)?  ➜ *Out of scope*.
2. Support for `.git/info/exclude`? ➜ Could be added with the same matcher.
3. Cache invalidation for long-running daemon mode.  Not relevant today (CLI is short-lived).

---

## 7. Timeline Estimate

| Task | ETA |
|------|-----|
| Add dependency & matcher | 1 h |
| Integrate into skip logic | 1 h |
| Flag & config plumbing | 0.5 h |
| Unit tests | 1 h |
| Docs & spec updates | 0.5 h |
| Total | **4 h** |

---

**Author**: @shunsuke – 2025-07-30
