#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'EOF'
Usage: ./release.sh vX.Y.Z /path/to/release-notes.md

Creates and pushes a git tag, waits for GoReleaser (CI) to publish the GitHub Release,
then updates the release title and notes using the provided English notes file via gh CLI.

Environment overrides:
  WAIT_TIMEOUT_SEC   Total seconds to wait for the release to appear (default: 300)
  WAIT_INTERVAL_SEC  Seconds between checks (default: 5)
EOF
}

log()   { printf '\033[1;34m[info]\033[0m %s\n' "$*"; }
warn()  { printf '\033[1;33m[warn]\033[0m %s\n' "$*"; }
error() { printf '\033[1;31m[err ]\033[0m %s\n' "$*" 1>&2; }

if [[ ${1:-} == "-h" || ${1:-} == "--help" || $# -ne 2 ]]; then
  usage
  exit 1
fi

TAG="$1"
NOTES_FILE="$2"

if [[ ! $TAG =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
  error "Invalid tag format: '$TAG' (expected vX.Y.Z)"
  exit 1
fi

if [[ ! -f "$NOTES_FILE" ]]; then
  error "Notes file not found: $NOTES_FILE"
  exit 1
fi

# Check required tools
for cmd in git gh; do
  if ! command -v "$cmd" >/dev/null 2>&1; then
    error "Required command not found: $cmd"
    exit 1
  fi
done

# Verify GitHub auth
if ! gh auth status >/dev/null 2>&1; then
  error "GitHub CLI is not authenticated. Run: gh auth login"
  exit 1
fi

# Ensure we are on main and working tree is clean
CURRENT_BRANCH=$(git rev-parse --abbrev-ref HEAD)
if [[ "$CURRENT_BRANCH" != "main" ]]; then
  warn "Current branch is '$CURRENT_BRANCH' (expected 'main')."
fi

if ! git diff --quiet || ! git diff --cached --quiet; then
  error "Working tree is not clean. Commit or stash changes first."
  exit 1
fi

log "Syncing with origin/main..."
git fetch origin main --tags

if git show-ref --tags --verify --quiet "refs/tags/$TAG"; then
  log "Tag $TAG already exists locally."
else
  log "Creating annotated tag $TAG..."
  git tag -a "$TAG" -m "Release $TAG"
fi

log "Pushing tag $TAG to origin..."
git push origin "$TAG"

WAIT_TIMEOUT_SEC=${WAIT_TIMEOUT_SEC:-300}
WAIT_INTERVAL_SEC=${WAIT_INTERVAL_SEC:-5}

log "Waiting for GitHub Release (created by GoReleaser) to appear for $TAG..."
start_ts=$(date +%s)
while true; do
  if gh release view "$TAG" >/dev/null 2>&1; then
    log "Release $TAG is available. Updating title and notes..."
    gh release edit "$TAG" --title "$TAG" --notes-file "$NOTES_FILE"
    log "Release notes updated for $TAG."
    break
  fi

  now_ts=$(date +%s)
  elapsed=$(( now_ts - start_ts ))
  if (( elapsed >= WAIT_TIMEOUT_SEC )); then
    warn "Timed out waiting for release to appear. Creating the release via gh as a fallback..."
    gh release create "$TAG" --title "$TAG" --notes-file "$NOTES_FILE" --latest
    log "Fallback release created for $TAG."
    break
  fi

  sleep "$WAIT_INTERVAL_SEC"
done

log "Done. View release: $(gh release view "$TAG" --json url --jq .url)"

