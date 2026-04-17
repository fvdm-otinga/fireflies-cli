#!/usr/bin/env bash
# scrub-fixtures.sh — detect or redact secrets in go-vcr fixture files.
#
# Usage:
#   scripts/scrub-fixtures.sh [--check] [<dir>]
#
# Arguments:
#   --check   exit 1 if any sensitive pattern is found (CI mode, no changes made)
#   <dir>     directory to scan (default: testdata/fixtures)
#
# Patterns detected:
#   - Bearer tokens:         Bearer [A-Za-z0-9._-]+
#   - Email-like values:     strings containing @ in YAML value position
#   - 36-char UUIDs:         xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx
#   - Explicit API key env:  FIREFLIES_API_KEY=...
#
# In default (replace) mode, matching strings are replaced with REDACTED
# in-place using sed. In --check mode, matches are printed and exit 1 is
# returned if any are found.

set -euo pipefail

CHECK_ONLY=0
FIXTURES_DIR="testdata/fixtures"

for arg in "$@"; do
  case "$arg" in
    --check) CHECK_ONLY=1 ;;
    *)       FIXTURES_DIR="$arg" ;;
  esac
done

if [[ ! -d "$FIXTURES_DIR" ]]; then
  echo "scrub-fixtures: directory not found: $FIXTURES_DIR" >&2
  exit 0
fi

# Patterns as extended-regex strings (ERE, compatible with grep -E and sed -E)
declare -a PATTERNS=(
  'Bearer [A-Za-z0-9._/+=-]+'
  '[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}'
  '[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}'
  'FIREFLIES_API_KEY=[^[:space:]]+'
)

FOUND=0

for pat in "${PATTERNS[@]}"; do
  # Find files with matches
  while IFS= read -r file; do
    if [[ -z "$file" ]]; then continue; fi
    if [[ $CHECK_ONLY -eq 1 ]]; then
      echo "SENSITIVE match in $file (pattern: $pat):" >&2
      grep -nE "$pat" "$file" >&2 || true
      FOUND=1
    else
      # In-place replacement — macOS sed requires '' for -i
      sed -i'' -E "s|${pat}|REDACTED|g" "$file"
      echo "scrubbed: $file (pattern: $pat)"
    fi
  done < <(grep -rlE "$pat" "$FIXTURES_DIR" 2>/dev/null || true)
done

if [[ $CHECK_ONLY -eq 1 && $FOUND -ne 0 ]]; then
  echo ""
  echo "ERROR: sensitive data found in fixtures. Run scripts/scrub-fixtures.sh to redact." >&2
  exit 1
fi

if [[ $CHECK_ONLY -eq 1 && $FOUND -eq 0 ]]; then
  echo "scrub-fixtures: no sensitive patterns found in $FIXTURES_DIR"
fi
