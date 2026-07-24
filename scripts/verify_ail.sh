#!/usr/bin/env bash
# ailang-code verify gate (see charter Repo Profile).
# Runs `ailang ai-check` — type check + Z3 contract verification, JSON output, one pass —
# on EVERY .ail module in the repo. ai-check is a strict superset of `ailang check`, so a
# single command covers both legs of the gate; modules without contracts pass verify with
# zero obligations.
#
# Module paths must match file paths relative to the SOURCE ROOT (MOD010). A source root is
# the directory that a module's namespace prefix is relative to, so we cd into that base
# before checking and pass the module path relative to it:
#   - design_docs/sketches/foo.ail declares `module sketches/foo`  -> base=design_docs, rel=sketches/foo.ail
#   - world/foo.ail              declares `module world/foo`       -> base=.            , rel=world/foo.ail
# ROOTS pairs each swept tree with the base cwd its module names are relative to,
# encoded as "base|tree". Extend ROOTS as source trees are added.
#
# The gate binary is configurable via AILANG_BIN (default: `ailang` on PATH). CI installs and
# checksum-verifies its own released binary and exports AILANG_BIN; this script never hardcodes
# a path so it stays CI-safe.
set -uo pipefail
cd "$(dirname "$0")/.."

AILANG_BIN="${AILANG_BIN:-ailang}"

# "base|tree": cd into `base`, then find/check .ail files under `tree` (relative to base).
ROOTS=(
  "design_docs|."
  ".|world"
)
fail=0
count=0
for entry in "${ROOTS[@]}"; do
  base="${entry%%|*}"
  tree="${entry#*|}"
  # searchdir is the tree to walk; when tree is ".", walk the base itself.
  if [ "$tree" = "." ]; then searchdir="$base"; else searchdir="$base/$tree"; fi
  [ -d "$searchdir" ] || continue
  while IFS= read -r -d '' f; do
    rel="${f#"$base"/}"
    count=$((count + 1))
    echo "── ai-check $f"
    if ! (cd "$base" && "$AILANG_BIN" ai-check "$rel"); then
      echo "✗ FAILED: $f"
      fail=1
    fi
  done < <(find "$searchdir" -name '*.ail' -print0 | sort -z)
done

echo "checked $count module(s)"
if [ "$count" -eq 0 ]; then
  echo "✗ no .ail modules found — the gate would be vacuous; failing loudly"
  exit 1
fi
exit $fail
