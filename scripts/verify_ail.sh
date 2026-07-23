#!/usr/bin/env bash
# ailang-code verify gate (see charter Repo Profile).
# Runs `ailang ai-check` — type check + Z3 contract verification, JSON output, one pass —
# on EVERY .ail module in the repo. ai-check is a strict superset of `ailang check`, so a
# single command covers both legs of the gate; modules without contracts pass verify with
# zero obligations.
# Module paths must match file paths relative to the source root, so we cd into each root
# before checking. Extend ROOTS as source trees are added.
set -uo pipefail
cd "$(dirname "$0")/.."

ROOTS=(design_docs)
fail=0
count=0
for root in "${ROOTS[@]}"; do
  [ -d "$root" ] || continue
  while IFS= read -r -d '' f; do
    rel="${f#"$root"/}"
    count=$((count + 1))
    echo "── ai-check $f"
    if ! (cd "$root" && ailang ai-check "$rel"); then
      echo "✗ FAILED: $f"
      fail=1
    fi
  done < <(find "$root" -name '*.ail' -print0 | sort -z)
done

echo "checked $count module(s)"
if [ "$count" -eq 0 ]; then
  echo "✗ no .ail modules found — the gate would be vacuous; failing loudly"
  exit 1
fi
exit $fail
