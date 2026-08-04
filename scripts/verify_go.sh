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

# This deny-list is the measured set: go1.26.0-go1.26.5 on darwin/arm64.
# Future go1.26.6 or go1.27.x versions are not covered here; the canary in this
# gate is the version-agnostic detector for any version that miscompiles the shape.
ACTIVE_GO=$(go env GOVERSION)      # observe; never assign
case "$ACTIVE_GO" in
  go1.26.0|go1.26.1|go1.26.2|go1.26.3|go1.26.4|go1.26.5)
    echo "verify_go.sh: FATAL: active toolchain $ACTIVE_GO miscompiles host/store/scan.go's" >&2
    echo "  array-literal shape (see design_docs/verification/w-race-gate-blindspot/). Pin a" >&2
    echo "  known-good toolchain, e.g. GOTOOLCHAIN=go1.25.6." >&2
    exit 1 ;;
esac

echo "── race-detector known-positive control"
go version
set +e
race_control_output="$(cd design_docs/verification/w-race-gate-blindspot/racecontrol && go run -race . 2>&1)"
set -e
printf '%s\n' "$race_control_output"
if ! grep -q 'WARNING: DATA RACE' <<<"$race_control_output"; then
  echo "verify_go.sh: FATAL: the race detector is not armed; every 0-races result in this gate is void" >&2
  exit 1
fi

echo "── go build ./..."
go version
go build ./...

echo "── go test ./... -count=1"
go version
go test ./... -count=1

# Measured at 78 s wall on darwin/arm64 at 7550ee9 under load ~4.5;
# host/broker was the 76.9 s critical path. The doc's ~179 s was not reproduced.
echo "── go test ./... -count=1 -race -timeout 8m"
go version
python3 - <<'PY'
import os, signal, subprocess, sys
cmd = ["go", "test", "./...", "-count=1", "-race", "-timeout", "8m"]
p = subprocess.Popen(cmd, start_new_session=True)
try:
    sys.exit(p.wait(timeout=600))
except subprocess.TimeoutExpired:
    os.killpg(os.getpgid(p.pid), signal.SIGKILL)
    print("verify_go.sh: FATAL: -race leg timed out after 600s", file=sys.stderr)
    sys.exit(124)
PY

echo "✓ go gate PASSED: build clean, plain and race tests pass with pinned AILANG_BIN ($ver)"
