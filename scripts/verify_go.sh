#!/usr/bin/env bash
# Go host build + test gate (see design_docs/planned/w-world-library-m1.md M6).
#
# The durable LOCAL mirror of CI's go-verify job: `go build ./...` then
# `go test ./... -count=1` from the repo root, resolved relative to this script
# so it runs from any cwd (like verify_ail.sh).
#
# ANTI-FALSE-GREEN GUARD (Standing Rule / V27 / B1): the host/replay tests
# `t.Skip()` silently when AILANG_BIN is unset (pinnedBinary), so a bare
# `go test` would report `ok` with the load-bearing replay assertions never
# running — the exact silent-skip class the mission forbids. This gate FAILS
# LOUDLY if AILANG_BIN is unset OR the binary it names does not report v0.30.0,
# turning a would-be false-green into a red local gate that mirrors CI (where
# the go-verify job exports AILANG_BIN before invoking this script). Locally,
# the operator exports AILANG_BIN=/tmp/ailang-v0300/ailang.
set -euo pipefail
cd "$(dirname "$0")/.."

if [ -z "${AILANG_BIN:-}" ]; then
  echo "✗ AILANG_BIN is unset — host/replay tests would t.Skip() silently and this gate would be false-green." >&2
  echo "  Export the pinned released binary, e.g. AILANG_BIN=/tmp/ailang-v0300/ailang" >&2
  exit 1
fi

ver="$("$AILANG_BIN" --version 2>&1)" || {
  echo "✗ AILANG_BIN=$AILANG_BIN could not be executed for a version check" >&2
  exit 1
}
if ! echo "$ver" | grep -q 'v0.30.0'; then
  echo "✗ AILANG_BIN=$AILANG_BIN does not report v0.30.0 (got: $ver) — replay goldens are v0.30.0-scoped." >&2
  exit 1
fi
echo "── AILANG_BIN=$AILANG_BIN ($ver)"

echo "── go build ./..."
go build ./...

echo "── go test ./... -count=1"
go test ./... -count=1

echo "✓ go gate PASSED: build clean, tests pass with pinned AILANG_BIN ($ver)"
