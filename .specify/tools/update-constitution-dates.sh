#!/usr/bin/env bash
# update-constitution-dates.sh
# Usage:
#   .specify/tools/update-constitution-dates.sh --ratified 2025-06-13
#   .specify/tools/update-constitution-dates.sh --amended 2025-10-03
#   .specify/tools/update-constitution-dates.sh --ratified 2025-06-13 --amended 2025-10-03

set -euo pipefail
script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
root_dir="$(cd "$script_dir/../.." && pwd)"
constitution_file="$root_dir/.specify/memory/constitution.md"

usage(){
  cat <<EOF
Usage: $0 [--ratified YYYY-MM-DD] [--amended YYYY-MM-DD]
Updates the constitution file's ratification and/or last-amended dates.
EOF
}

if [[ $# -eq 0 ]]; then
  usage
  exit 1
fi

ratified=""
amended=""
while [[ $# -gt 0 ]]; do
  case "$1" in
    --ratified)
      ratified="$2"
      shift 2
      ;;
    --amended)
      amended="$2"
      shift 2
      ;;
    *)
      echo "Unknown arg: $1" >&2
      usage
      exit 1
      ;;
  esac
done

if [[ ! -f "$constitution_file" ]]; then
  echo "Constitution file not found at $constitution_file" >&2
  exit 1
fi

tmpfile="${constitution_file}.tmp"
cp "$constitution_file" "$tmpfile"

if [[ -n "$ratified" ]]; then
  sed -i "s/TODO(RATIFICATION_DATE): confirm original adoption date/${ratified}/g" "$tmpfile" || true
  # Also replace any bare [RATIFICATION_DATE] tokens if present
  sed -i "s/\[RATIFICATION_DATE\]/${ratified}/g" "$tmpfile" || true
fi

if [[ -n "$amended" ]]; then
  sed -i "s/[0-9]\{4\}-[0-9]\{2\}-[0-9]\{2\}/${amended}/1" "$tmpfile" || true
  # Also replace explicit LAST_AMENDED_DATE tokens if present
  sed -i "s/\[LAST_AMENDED_DATE\]/${amended}/g" "$tmpfile" || true
fi

mv "$tmpfile" "$constitution_file"
echo "Updated constitution file: $constitution_file"
