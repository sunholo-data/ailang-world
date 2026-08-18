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
# EXACT TOKEN, NEVER A SUBSTRING (iter-66). This assertion used `grep -q 'v0.30.0'`, which the
# string `v0.30.0-205-g54d6bd191-dirty` SATISFIES — so the anti-false-green guard admitted exactly
# the dirty dev build CLAUDE.md forbids, measured against the real script before the fix. The
# released binary reports the bare token in both arms that matter: `AILANG v0.30.0` locally
# (/tmp/ailang-v0300/ailang) and `AILANG v0.30.0` on the linux runner (CI step log, run
# 31249744703, `go host build + test gate`), so tightening this cannot red CI.
ver_tok="$(printf '%s\n' "$ver" | head -1 | awk '{print $2}')"
if [ "$ver_tok" != 'v0.30.0' ]; then
  echo "✗ AILANG_BIN=$AILANG_BIN does not report exactly v0.30.0 (got: ${ver_tok:-<none>}) — replay goldens are v0.30.0-scoped." >&2
  echo "  A -dirty or -N-g<sha> suffix is a REJECTION, not a match: dev builds are forbidden by CLAUDE.md." >&2
  exit 1
fi
echo "── AILANG_BIN=$AILANG_BIN ($ver)"

echo "── tracked-binary hygiene gate"
# SM.A (squash 13315da) committed a 15.7 MB compiled darwin/arm64 `ailang-worldd`
# Mach-O into the repository root. It passed the codex executor, a sonnet evaluator
# scoring 87/100 with zero blocking findings, the controller's four-gate re-run and
# both CI jobs — because nothing anywhere looked. Compiled artifacts are permanent
# git bloat, are platform-specific in a repo whose whole thesis is byte-exact
# reproducibility, and silently dirty the shared checkout the moment anyone rebuilds
# (which changes controller behaviour at Gate 0, while Principle 0 forbids stashing).
#
# The detector is git's OWN binary classification: `git diff --numstat` prints `-`
# for both the add and delete counts of any blob it considers binary. That is
# portable between darwin and the ubuntu runner in a way `file(1)`'s wording
# ("Mach-O" vs "ELF") is not, and it needs no allowlist — measured at 0eb58f5,
# exactly one of 142 tracked files was binary, and it was the stray artifact.
#
# NON-VACUITY (rule 3a): an enumeration that returns nothing is indistinguishable
# from a clean repo — same silence, same exit path. So the TOTAL file count is
# asserted non-zero in the SAME call, before the binary count is believed. A
# detector that can see nothing fails loudly instead of passing quietly.
empty_tree=$(git hash-object -t tree /dev/null)
binary_numstat="$(git diff --numstat "$empty_tree" HEAD)"
tracked_total=$(printf '%s\n' "$binary_numstat" | grep -c . || true)
if [ "$tracked_total" -eq 0 ]; then
  echo "verify_go.sh: FATAL: the tracked-binary detector enumerated 0 files; the instrument" >&2
  echo "  is broken, so every 'no binaries' result it reports is void." >&2
  exit 1
fi
tracked_binaries="$(printf '%s\n' "$binary_numstat" | awk -F'\t' '$1 == "-" && $2 == "-" { print $3 }')"
if [ -n "$tracked_binaries" ]; then
  echo "verify_go.sh: FATAL: compiled/binary artifacts are tracked in git:" >&2
  printf '%s\n' "$tracked_binaries" | sed 's/^/    /' >&2
  echo "  Remove each with 'git rm --cached <path>' and add it to .gitignore." >&2
  exit 1
fi
echo "   ✓ 0 binary blobs among $tracked_total tracked files"

echo "── mission routing + decision-ledger gate"
/bin/bash tools/launchd/test_mission_routing.sh
/bin/bash -n tools/launchd/mission-control.sh tools/launchd/derive-planner-lane.sh scripts/mission_decisions.sh

# DRIVER DRIFT GATE — D-WORLD-DRIVER-1, RESOLVED B (Mark, attended 2026-08-17).
# The driver is FLEET-owned: changes land here only as fleet-authored commits,
# never as World-controller edits. launchd executes this repo's WORKING TREE
# (dev.ailang.mission-world.plist ProgramArguments), so an uncommitted driver is
# a live driver that exists in no repository — iter-89 measured exactly that
# state lurking for two days, carrying the human decision ledger with it.
# In CI the checkout is clean and this passes; on the rig, mid-propagation dirt
# reds LOUDLY until the fleet commits it. That red is the point, not a nuisance.
# Path-liveness control: prove git is scanning a real tracked set before
# trusting an empty diff — a mistyped path would pass vacuously.
driver_tracked=$(git ls-files tools/launchd/ scripts/mission_decisions.sh | wc -l | tr -d ' ')
if [ "$driver_tracked" -lt 5 ]; then
  echo "verify_go.sh: FATAL: driver drift gate control failed — only $driver_tracked tracked driver files (expected >=5); the gate is not scanning what it claims" >&2
  exit 1
fi
driver_drift=$(git status --porcelain -- tools/launchd/ scripts/mission_decisions.sh)
if [ -n "$driver_drift" ]; then
  echo "verify_go.sh: FATAL: DRIVER DRIFT (D-WORLD-DRIVER-1) — the running driver differs from the committed one:" >&2
  printf '%s\n' "$driver_drift" | sed 's/^/    /' >&2
  echo "  The driver is fleet-owned; land this as a fleet-authored commit. World's controller must not edit or absorb it." >&2
  exit 1
fi
echo "   ✓ driver drift gate: $driver_tracked tracked driver files, working tree matches HEAD"

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
