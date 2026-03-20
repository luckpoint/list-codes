## Release Guide

This document describes the end-to-end release process for `list-codes`, including semantic versioning, tagging, automated releases via GoReleaser, creating English release notes with the `gh` CLI, and Homebrew tap updates. An example for version v0.5.2 is included.

### Prerequisites
- Go toolchain installed
- Write access to the GitHub repository and the Homebrew tap repository (`luckpoint/list-codes` tap)
- Logged in to Git (SSH or HTTPS) and GitHub

### Versioning
- Follow Semantic Versioning: MAJOR.MINOR.PATCH
- Example in this guide: v0.5.2

### 1) Prepare the Release
1. Update documentation and examples if needed.
2. Ensure tests pass locally:
   ```bash
   go test ./...
   ```
3. Commit and merge changes to `main`.

### 2) Tag the Version
Create and push a signed or lightweight tag for the release:
```bash
git pull origin main
git tag v0.5.2
git push origin v0.5.2
```

### 3) Release via GoReleaser (Automated)
Releases are automated with GoReleaser. Pushing a tag (e.g., `v0.5.2`) triggers the pipeline to:
- Build and upload release artifacts
- Create or update the GitHub Release
- Update the Homebrew tap automatically

Optional local validation before pushing the tag:
```bash
# Dry run (no publish), builds locally and validates the pipeline
goreleaser release --clean --skip=publish --snapshot

# Real local release (normally handled by CI on tag push)
# goreleaser release --clean
```

### 4) Create English Release Notes with gh
Use the GitHub CLI to draft and publish release notes in English. Ensure you are authenticated:
```bash
gh auth status || gh auth login
```

There are two common paths depending on whether GoReleaser creates the release entry:

- If GoReleaser already created the release for `v0.5.2`, edit notes/title:
```bash
# Use a notes file (recommended)
gh release edit v0.5.2 \
  --title "v0.5.2" \
  --notes-file docs/release-notes/v0.5.2.md

# Or inline notes (short form)
gh release edit v0.5.2 \
  --title "v0.5.2" \
  --notes "Minor utils fixes and small improvements. No breaking changes."
```

- If the release does not exist yet (rare with GoReleaser), create it:
```bash
# Generate notes from PRs/commits (English titles recommended)
gh release create v0.5.2 \
  --title "v0.5.2" \
  --generate-notes \
  --latest

# Or with your curated notes file
gh release create v0.5.2 \
  --title "v0.5.2" \
  --notes-file docs/release-notes/v0.5.2.md \
  --latest
```

Verification:
```bash
gh release view v0.5.2 --web   # open in browser
gh release view v0.5.2 --json name,tagName,url
```

Note: Go users can install a tagged version via:
```bash
go install github.com/luckpoint/list-codes/cmd/list-codes@v0.5.2
```

### 5) Homebrew Tap (Automated)
Homebrew formula updates are automated. No manual action is typically required.

Verification after a release:
```bash
brew update
brew upgrade list-codes || brew install list-codes
list-codes --version
```

### 6) Post-release Verification
- Homebrew installs/upgrades successfully
- `list-codes --version` prints `v0.5.2`
- `go install ...@v0.5.2` works
- Docs and README reference the new version where relevant

### Release Notes Template (English)
Use this as a checklist when drafting the GitHub Release notes:

```
## Changes
- Briefly summarize the main changes in this release.

## Fixes
- List fixes (reference issues/PRs when applicable).

## Breaking Changes
- Note any breaking changes and migration steps.

## Thanks
- Acknowledge contributors.
```

### Example Notes for v0.5.2
- Minor utils fixes and small improvements
- No breaking changes
- Thanks to all contributors


