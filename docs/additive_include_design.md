# Additive --include Behavior Design

## Overview

This document describes the proposed change to the `--include` option behavior from "whitelist mode" to "additive mode" to better meet user expectations and workflows.

## Current Behavior (Whitelist Mode)

### How it works:
- When `--include` is specified, the tool switches to "include-only" mode
- Only explicitly included files/directories are processed
- All other files are excluded (except parent directories needed for traversal)
- Overrides default exclusions (dotfiles, gitignore patterns)

### Example:
```bash
# Current: Only processes files in .github/ directory
list-codes --include ".github"

# Current: Only processes .env and .dockerignore files
list-codes --include ".env" --include ".dockerignore"
```

### Problems:
1. **Unexpected behavior**: Users expect `--include` to ADD items, not replace the entire inclusion set
2. **Loss of default files**: Source code files are excluded unless explicitly included
3. **Verbose for common use cases**: Need to specify many paths to get normal behavior + extras

## Proposed Behavior (Additive Mode)

### How it would work:
- `--include` adds files/directories to the default inclusion set
- Normal default behavior continues (includes source files, respects language detection, etc.)
- Specified `--include` items are processed in addition to defaults
- Can still override default exclusions (like dotfiles) but only for specified items
- Maintains all existing filtering priority and logic

### Example:
```bash
# Proposed: Processes all normal source files + .github/ directory
list-codes --include ".github"

# Proposed: Processes all normal source files + .env and .dockerignore files
list-codes --include ".env" --include ".dockerignore"

# Proposed: Processes all normal source files + entire .vscode directory
list-codes --include ".vscode/**"
```

## Implementation Design

### High-Level Changes

The core change involves modifying the `shouldSkipDir` function in `utils/file.go` to use additive logic instead of whitelist logic.

### Current Priority System (in shouldSkipDir):
1. **Explicit exclusions** (`--exclude`) - Always cause a skip
2. **Include whitelist** - When `--include` used, only whitelisted items are kept
3. **Gitignore patterns** - Files matching `.gitignore` are skipped  
4. **Default dotfile exclusion** - Files/dirs starting with `.` are skipped
5. **Include-only mode** - When `--include` used, anything not whitelisted is skipped

### Proposed Priority System:
1. **Explicit exclusions** (`--exclude`) - Always cause a skip
2. **Include additions** - Items matching `--include` patterns are never skipped (override defaults)
3. **Gitignore patterns** - Files matching `.gitignore` are skipped
4. **Default dotfile exclusion** - Files/dirs starting with `.` are skipped
5. **No include-only mode** - Default source file inclusion continues

### Key Logic Changes

#### Current Logic (Whitelist):
```go
// In shouldSkipDir function
if len(includePaths) > 0 {
    isNeededForInclude := false
    for incPath := range includePaths {
        if strings.HasPrefix(incPath, absPath) || strings.HasPrefix(absPath, incPath) {
            isNeededForInclude = true
            break
        }
    }
    
    if isNeededForInclude {
        return false  // Don't skip
    }
    // If include paths specified but this isn't needed, skip it
    return true
}
```

#### Proposed Logic (Additive):
```go
// In shouldSkipDir function
if len(includePaths) > 0 {
    isExplicitlyIncluded := false
    for incPath := range includePaths {
        if strings.HasPrefix(incPath, absPath) || strings.HasPrefix(absPath, incPath) {
            isExplicitlyIncluded = true
            break
        }
    }
    
    if isExplicitlyIncluded {
        return false  // Don't skip - explicitly included
    }
    // Continue with normal default logic (don't skip unless other rules apply)
}
```

### Detailed Implementation Steps

#### 1. Modify `shouldSkipDir` function (`utils/file.go` lines 84-137)

**Before (Priority 2 - Include whitelist):**
- If `--include` specified, only whitelisted items are kept
- Introduces early return that skips non-whitelisted items

**After (Priority 2 - Include additions):**
- If `--include` specified, whitelisted items override later exclusion rules
- Remove the include-only mode logic
- Let normal default logic continue for non-included items

#### 2. Update Function Documentation

Update comments in `shouldSkipDir` to reflect the new additive behavior:

```go
// shouldSkipDir determines if a directory should be skipped during traversal.
// Priority order (higher priority rules override lower priority):
// 1. Explicit exclusions (--exclude) - always skip
// 2. Include additions (--include) - never skip, overrides defaults  
// 3. Gitignore patterns - skip if matched
// 4. Default dotfile exclusion - skip if starts with '.'
// 5. Normal source file inclusion continues
```

#### 3. Update CLI Help Text

Update the help description for the `--include` flag:

**Current:**
```go
rootCmd.PersistentFlags().StringSliceVarP(&includes, "include", "i", []string{}, "Folder path to include (repeatable)")
```

**Proposed:**
```go
rootCmd.PersistentFlags().StringSliceVarP(&includes, "include", "i", []string{}, "Additional folder/file paths to include beyond defaults (repeatable)")
```

#### 4. Update Tests

Update test cases in `utils/file_test.go` to reflect additive behavior:

- Tests that expect include-only mode should be updated to expect additive mode
- Add new tests that verify default files are still included when using `--include`
- Test cases verifying override behavior for dotfiles should remain the same

#### 5. Update Documentation

Update README.md and any other documentation to reflect the new behavior.

## Migration Considerations

### Backward Compatibility Impact

This is a **breaking change** that will affect existing users who rely on the current whitelist behavior.

### Migration Strategies

#### Option 1: Introduce New Flag (Recommended)
- Keep current `--include` behavior as-is
- Add new `--also-include` or `--add-include` flag with additive behavior
- Deprecate `--include` in favor of the new flag over time

#### Option 2: Feature Flag
- Add a configuration option to switch between whitelist and additive modes
- Default to current behavior initially, migrate to additive as default later

#### Option 3: Direct Change (Simplest)
- Change behavior directly
- Provide clear migration notes in changelog
- Most users likely expect additive behavior anyway

## Examples of New Behavior

### Common Use Cases

#### Adding dotfiles to normal processing:
```bash
# Process all normal source files + .env file
list-codes --include ".env"

# Process all normal source files + .github directory
list-codes --include ".github"

# Process all normal source files + multiple dotfiles
list-codes --include ".env" --include ".dockerignore" --include ".github"
```

#### Adding normally excluded directories:
```bash
# Process all normal source files + node_modules (if needed for analysis)
list-codes --include "node_modules"

# Process all normal source files + build output directory
list-codes --include "dist" --include "build"
```

#### Combined with exclude:
```bash
# Process all normal files + .github, but exclude tests
list-codes --include ".github" --exclude "**/*_test.go"
```

## Benefits of Additive Behavior

1. **Intuitive**: Matches user expectations of "include" meaning "add to"
2. **Flexible**: Allows fine-tuning default behavior without losing it
3. **Efficient**: Common use case of "normal files + some extras" becomes simple
4. **Composable**: Works well with `--exclude` for fine-grained control

## Testing Strategy

### Existing Tests to Update
- Include-only mode tests should verify additive behavior instead
- Verify that normal source files are still included when using `--include`

### New Tests to Add
- Test that default language detection still works with `--include`
- Test that dependency files are still collected with `--include`
- Test interaction between `--include` and `--exclude` in additive mode
- Test that included dotfiles override default exclusions while normal exclusions continue

### Integration Tests
- End-to-end tests with real project structures
- Verify that output includes expected source files + additional included items
- Performance tests to ensure no regression with additive logic

## Implementation Priority

1. **High**: Core logic change in `shouldSkipDir`
2. **High**: Update test suite 
3. **Medium**: Update documentation and help text
4. **Low**: Consider migration strategy for existing users

This design maintains the flexibility and power of the current include system while making it more intuitive and user-friendly for common workflows.