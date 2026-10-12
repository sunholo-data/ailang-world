#!/usr/bin/env bash
# Regenerates the committed ai-check fixtures from the pinned AILANG binary.
# Run by hand only, never by a test. The ai-check exit code is data, not a
# failure: a refuted or erroring module is a fixture too.
set -euo pipefail
cd "$(dirname "$0")"
BIN="${AILANG_BIN:?AILANG_BIN must name the pinned v0.41.0 binary}"
ver="$("$BIN" --version 2>/dev/null | head -1)"
if [ "$ver" != "AILANG v0.41.0" ]; then
  echo "regen.sh: pinned binary reports '$ver'; want 'AILANG v0.41.0'" >&2
  exit 1
fi
scratch="$(mktemp -d)"
trap 'rm -rf "$scratch"' EXIT
mkdir "$scratch/fixtures"
for n in verifiedpass refuted encerror mixed skippedshape loaderror nocontracts; do
  cp "$n.ail" "$scratch/fixtures/$n.ail"
  ( cd "$scratch" && "$BIN" ai-check -timeout 5s "fixtures/$n.ail" ) > "$n.json" 2>/dev/null || true
done
echo "regen.sh: wrote 7 fixtures"
