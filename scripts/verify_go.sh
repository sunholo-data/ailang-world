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
# the operator exports AILANG_BIN=$HOME/.pinned-ailang/ailang (moved off /tmp 2026-09-03:
# macOS wipes /private/tmp on boot, and it took the pin with it — see the charter Repo Profile).
set -euo pipefail
cd "$(dirname "$0")/.."

check_evidence_manifest() {
  python3 - "$1" "$2" <<'PY'
import collections, json, sys

REQUIRED_EVIDENCE_TESTS = {
    "TestAttackerChosenValidatorCannotMintForHostAuthority",
    "TestConstructorNamesActuallyUsedUnorderedTimeouts",
    "TestConstructorPinsBusyTimeoutBelowObjectReadTimeout",
    "TestConstructorRefusesEmptyCompilerVersion",
    "TestConstructorRefusesEmptyRequiredIdentities",
    "TestConstructorRefusesNilReader",
    "TestConstructorRefusesNonPositiveObjectReadTimeout",
    "TestConstructorRefusesUnknownBusyTimeout",
    "TestConstructorRefusesUnsetCompilerIdentity",
    "TestDecodeProposalCapsBeforeParse",
    "TestEnvelopeCanonicalRoundTripAndMACDeferral",
    "TestEnvelopeCarriesTheReportItAlreadyDecoded",
    "TestEnvelopeStrictRefusals",
    "TestFailedProofReportIsRefused",
    "TestIncompleteProofReportIsRefused",
    "TestInvalidProofRefIsRefused",
    "TestMalformedProofReportIsRefused",
    "TestMismatchedProofSubjectIsRefused",
    "TestMismatchedProofToolIsRefused",
    "TestMissingProofReportIsRefused",
    "TestNestingDepthBombWithinByteCapIsRefused",
    "TestOtherwisePerfectReportWithWrongMACIsUnauthenticated",
    "TestOtherwisePerfectReportWithoutMACIsUnauthenticated",
    "TestOversizeProofReportIsRefused",
    "TestPayloadHashMismatchIsRefused",
    "TestProofReportCanonicalRoundTrip",
    "TestProofReportCaps",
    "TestProofReportStrictRefusals",
    "TestProposalStrictRefusals",
    "TestPublicAuthoritySurfaceIsFrozen",
    "TestReaderWaitBoundsCannotBeLostThroughWrapper",
    "TestRealStoreBlockedObjectReadReturnsWithinObjectReadTimeout",
    "TestTruncatedTailReportIsRefusedNotPanicked",
    "TestValidatorMintIdentitiesAreDistinct",
    "TestWrongInterfaceIsRefused",
    "TestWrongSemanticIDIsRefused",
    "TestZeroValueForgeryCannotResolve",
}
EXACT_EVIDENCE_TESTS = int(sys.argv[2])

if not REQUIRED_EVIDENCE_TESTS:
    sys.stderr.write("verify_go.sh: FATAL INSTRUMENT FAILURE: REQUIRED_EVIDENCE_TESTS is empty\n")
    sys.exit(1)
events = []
try:
    with open(sys.argv[1], encoding="utf-8") as stream:
        for line in stream:
            if line.strip():
                events.append(json.loads(line))
except Exception as exc:
    sys.stderr.write("verify_go.sh: FATAL INSTRUMENT FAILURE: cannot parse evidence test JSON: %s\n" % exc)
    sys.exit(1)

package_passes = [e for e in events
                  if e.get("Action") == "pass" and not e.get("Test")
                  and e.get("Package", "").endswith("/host/evidence")]
if not package_passes:
    sys.stderr.write("verify_go.sh: FATAL INSTRUMENT FAILURE: zero passing host/evidence packages discovered\n")
    sys.exit(1)

terminal_passes = [e["Test"] for e in events
                   if e.get("Action") == "pass" and e.get("Test") and "/" not in e["Test"]]
if not terminal_passes:
    sys.stderr.write("verify_go.sh: FATAL INSTRUMENT FAILURE: observed terminal named-test pass set is empty\n")
    sys.exit(1)

counts = collections.Counter(terminal_passes)
duplicates = sorted(name for name, count in counts.items() if count != 1)
observed = set(counts)
missing = sorted(REQUIRED_EVIDENCE_TESTS - observed)
extra = sorted(observed - REQUIRED_EVIDENCE_TESTS)
skipped = sorted({e["Test"] for e in events if e.get("Action") == "skip" and e.get("Test")
                  and "/" not in e["Test"] and e["Test"] in REQUIRED_EVIDENCE_TESTS})
failed = sorted({e["Test"] for e in events if e.get("Action") == "fail" and e.get("Test")
                 and "/" not in e["Test"] and e["Test"] in REQUIRED_EVIDENCE_TESTS})
if missing or extra or skipped or failed or duplicates or len(observed) != EXACT_EVIDENCE_TESTS:
    sys.stderr.write("evidence test set differs from REQUIRED_EVIDENCE_TESTS\n")
    sys.stderr.write("  missing=%s\n  skipped=%s\n  failed=%s\n  duplicate=%s\n  extra=%s\n" %
                     (missing, skipped, failed, duplicates, extra))
    sys.stderr.write("  observed_unique=%d exact_required=%d\n" %
                     (len(observed), EXACT_EVIDENCE_TESTS))
    sys.exit(1)

print("   ✓ all %d required top-level evidence tests passed exactly once" % EXACT_EVIDENCE_TESTS)
PY
}

check_mission_config() {
  validator="scripts/mission_decisions.sh"
  charter="design_docs/world-mission.md"
  if [ ! -r "$validator" ]; then
    echo "verify_go.sh: FATAL: mission validator is missing or unreadable: $validator" >&2
    return 1
  fi
  if [ ! -r "$charter" ]; then
    echo "verify_go.sh: FATAL: World charter is missing or unreadable: $charter" >&2
    return 1
  fi
  /bin/bash -n "$validator"
  /bin/bash "$validator" --check --file "$charter"
  echo "   World decision-ledger validated; centralized runtime currency is NOT certified here."
}

if [ "${1:-}" = "--evidence-manifest-check" ]; then
  if [ "$#" -ne 3 ]; then
    echo "usage: $0 --evidence-manifest-check JSON EXACT" >&2
    exit 2
  fi
  check_evidence_manifest "$2" "$3"
  exit $?
fi

if [ "${1:-}" = "--mission-config-check" ]; then
  if [ "$#" -ne 1 ]; then
    echo "usage: $0 --mission-config-check" >&2
    exit 2
  fi
  check_mission_config
  exit $?
fi

if [ "$#" -gt 0 ]; then
  echo "verify_go.sh: unrecognized first argument: $1" >&2
  exit 2
fi

if [ -z "${AILANG_BIN:-}" ]; then
  echo "✗ AILANG_BIN is unset — host/replay tests would t.Skip() silently and this gate would be false-green." >&2
  echo "  Export the pinned released binary, e.g. AILANG_BIN=/tmp/ailang-v0300/ailang" >&2
  exit 1
fi

# STDOUT ONLY (iter-89's class, 4th site, bitten 2026-08-18): a >200MB
# ~/.ailang/state makes every ailang call emit `Observatory: NNNMB` on STDERR
# with a timestamp prefix; capturing 2>&1 here made head-1/awk-2 parse that
# timestamp as the version token and reject the correctly-pinned binary.
# An instrument that captures more than it parses can be voided by a process
# it has nothing to do with — never merge streams into parsed output.
ver="$("$AILANG_BIN" --version 2>/dev/null)" || {
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

echo "── World mission-input gate (decision ledger)"
check_mission_config

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

# --- P1 (queue row 48 / D-WORLD-28): the SELECTED toolchain must be at-or-above the
# root module floor. Below it, the nested race-control module's own floor can exceed
# the toolchain that runs it and the known-positive control at :229 goes silently
# unarmed. The floor is READ from ./go.mod, never hardcoded.
go_version_ge() {   # rc 0 = $1 >= $2 ; rc 1 = $1 < $2 ; rc 2 = a token is not a release version
  awk -v a="$1" -v b="$2" 'BEGIN{
    if (a !~ /^go1\.[0-9]+(\.[0-9]+)?$/ || b !~ /^go1\.[0-9]+(\.[0-9]+)?$/) exit 2
    sub(/^go/,"",a); sub(/^go/,"",b)
    na=split(a,A,"."); nb=split(b,B,".")
    n=(na>nb?na:nb)
    for(i=1;i<=n;i++){x=(i<=na?A[i]+0:0); y=(i<=nb?B[i]+0:0)
      if(x>y) exit 0; if(x<y) exit 1}
    exit 0}'
}
root_go_lines=$(awk '/^go /{n++} END{print n+0}' go.mod)
if [ "$root_go_lines" -ne 1 ]; then
  echo "verify_go.sh: FATAL: root go.mod has $root_go_lines column-0 'go ' lines, want exactly 1;" >&2
  echo "  the root module floor cannot be read, so the race-detector control cannot be bounded." >&2
  exit 1
fi
ROOT_FLOOR="go$(awk '/^go /{print $2; exit}' go.mod)"
# Row 61: the comparator's verdict is consumed by `if` DIRECTLY. There is no intermediate
# variable and no `$?` read on the success path, so a dataflow break between the comparison and
# the branch cannot open the gate. The predecessor laundered the verdict through a reassignable
# `floor_rc`, and ONE inserted line (`floor_rc=0`) reopened ALL THREE refusal branches at once --
# below-floor AND malformed-token -- while the gate went on printing an arithmetically false
# success line. Note what does NOT fix it: dropping the variable and reading `case "$?"` is
# fail-OPEN even unmutated, because the intervening `set -e` is itself a successful command that
# resets `$?` (measured). `$?` is still read below, but only to ATTRIBUTE the refusal: every arm
# of that `case` exits, so a broken attribution is still a refusal.
if go_version_ge "$ACTIVE_GO" "$ROOT_FLOOR"; then
  :
else
  case $? in
    1) echo "verify_go.sh: FATAL: active toolchain $ACTIVE_GO is BELOW the root module floor $ROOT_FLOOR;" >&2
       echo "  the race-detector known-positive control would be disarmed. Pin GOTOOLCHAIN to $ROOT_FLOOR or above." >&2
       exit 1 ;;
    *) echo "verify_go.sh: FATAL: cannot order toolchain tokens (ACTIVE_GO=$ACTIVE_GO ROOT_FLOOR=$ROOT_FLOOR);" >&2
       echo "  at least one is not a well-formed goX.Y[.Z] release version." >&2
       exit 1 ;;
  esac
fi
echo "   ✓ toolchain floor gate: $ACTIVE_GO >= root module floor $ROOT_FLOOR"

echo "── race-detector known-positive control"
go version
set +e
race_control_output="$(cd design_docs/verification/w-race-gate-blindspot/racecontrol && GOTOOLCHAIN="$ACTIVE_GO" go run -race . 2>&1)"
set -e
printf '%s\n' "$race_control_output"
if ! grep -q 'WARNING: DATA RACE' <<<"$race_control_output"; then
  echo "verify_go.sh: FATAL: the race detector is not armed; every 0-races result in this gate is void" >&2
  exit 1
fi

echo "── go build ./..."
go version
go build ./...

echo "── focused host/evidence named-manifest gate (37 exact top-level tests)"
evidence_json="$(mktemp "${TMPDIR:-/tmp}/verify-go-evidence.XXXXXX")"
trap 'rm -f "$evidence_json"' EXIT
set +e
go test -json ./host/evidence -count=1 >"$evidence_json"
evidence_test_rc=$?
set -e
check_evidence_manifest "$evidence_json" 37
if [ "$evidence_test_rc" -ne 0 ]; then
  echo "verify_go.sh: FATAL: host/evidence go test exited rc=$evidence_test_rc" >&2
  exit "$evidence_test_rc"
fi
rm -f "$evidence_json"
trap - EXIT

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
