#!/usr/bin/env bash
# Deterministically project the canonical World semantic modules into world/core.
set -euo pipefail

cd "$(dirname "$0")/.."

readonly SOURCE_DIR="world"
readonly PACKAGE_DIR="packages/world-core"
readonly DEST_DIR="$PACKAGE_DIR/world"
readonly EXPECTED_MODULE_COUNT=4

readonly ALLOWLIST=(
  "types.ail"
  "contracts.ail"
  "transitions.ail"
  "logepoch.ail"
)

allowlist_file="$(mktemp -t world_package_allowlist.XXXXXX)"
staging_dir="$(mktemp -d "$PACKAGE_DIR/.world-stage.XXXXXX")"
cleanup() {
  rm -f "$allowlist_file"
  rm -rf "$staging_dir"
}
trap cleanup EXIT

# The staging directory was newly created by mktemp and must be empty before projection.
if [ -n "$(find "$staging_dir" -mindepth 1 -print -quit)" ]; then
  printf '%s\n' "REFUSED: staging directory is not empty: $staging_dir" >&2
  exit 1
fi

printf '%s\n' "${ALLOWLIST[@]}" > "$allowlist_file"
allowlisted_count="$(wc -l < "$allowlist_file" | tr -d '[:space:]')"
if [ "$allowlisted_count" -ne "$EXPECTED_MODULE_COUNT" ]; then
  printf '%s\n' "REFUSED: allowlist must contain exactly $EXPECTED_MODULE_COUNT modules; found $allowlisted_count" >&2
  exit 1
fi

# Refuse source-set growth until the projection contract is deliberately updated.
while IFS= read -r source_path; do
  module="${source_path#"$SOURCE_DIR/"}"
  if [ "$module" = "_smoke.ail" ]; then
    continue
  fi
  if ! grep -Fqx -- "$module" "$allowlist_file"; then
    printf '%s\n' "REFUSED: unexpected fifth .ail source module: $source_path" >&2
    exit 1
  fi
done < <(find "$SOURCE_DIR" -maxdepth 1 -type f -name '*.ail' -print | sort)

iterated_count=0
while IFS= read -r module; do
  source_path="$SOURCE_DIR/$module"
  [ -f "$source_path" ] || {
    printf '%s\n' "REFUSED: allowlisted source module is missing: $source_path" >&2
    exit 1
  }
  cp "$source_path" "$staging_dir/$module"
  iterated_count=$((iterated_count + 1))
done < "$allowlist_file"

printf 'allowlisted modules: iterated=%s wc-l=%s\n' "$iterated_count" "$allowlisted_count"
if [ "$iterated_count" -ne "$allowlisted_count" ]; then
  printf '%s\n' "REFUSED: iterated module count $iterated_count differs from allowlist wc -l $allowlisted_count" >&2
  exit 1
fi

# Replace the projection wholesale so stale files can never survive an incremental build.
rm -rf "$DEST_DIR"
mv "$staging_dir" "$DEST_DIR"
printf 'projected %s modules into %s\n' "$iterated_count" "$DEST_DIR"
