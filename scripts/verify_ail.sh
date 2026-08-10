#!/usr/bin/env bash
# ailang-code verify gate (see charter Repo Profile; design_docs/planned/w-m1-ailang-hardening.md §D5).
#
# NON-VACUOUS required-check-manifest gate over every .ail module in the repo. Runs on the pinned
# released `ailang` (v0.30.0). Two facts make exit codes insufficient, so this gate PARSES JSON:
#   - a Z3 encoding error exits 0 SILENTLY (verify.errors>0 but rc 0, V10);
#   - a dropped contract vanishes from verify.results[] with rc 0 (V20);
#   - a deleted test pair still clears an aggregate floor (V21).
# So the gate asserts a HARDCODED manifest of proven-contract and passing-test IDENTITIES (not
# aggregate counts), with exact totals as SECONDARY checks only.
#
# Gate policy is HARDCODED and deliberately NOT env-overridable: there are no MIN_* or timeout
# knobs. AILANG_BIN (the binary PATH, not gate strength) is the ONLY configurable knob. CI installs
# and checksum-verifies its own released binary and exports AILANG_BIN.
#
# Module paths must match file paths relative to the SOURCE ROOT (MOD010). A source root is the
# directory a module's namespace prefix is relative to, so we cd into that base before checking and
# pass the module path relative to it:
#   - design_docs/sketches/foo.ail declares `module sketches/foo` -> base=design_docs, rel=sketches/foo.ail
#   - world/foo.ail              declares `module world/foo`      -> base=.          , rel=world/foo.ail
# ROOTS pairs each swept tree with the base cwd its module names are relative to, as "base|tree".
#
# python3 is REQUIRED (JSON parse; no jq dependency). Its absence fails the script LOUDLY rather
# than passing vacuously. It is present on macOS dev machines and ubuntu-latest CI runners.
#
# NOTE (test-name collisions, V22): `ailang test --format json world/` merges modules into one JSON
# with BARE test names. Two test-bearing modules exist now (world/logepoch + world/contracts) with
# distinct function names ⇒ no collision. A future colliding name would force Leg 2 to per-module
# runs (documented escape hatch).
set -uo pipefail
cd "$(dirname "$0")/.."

AILANG_BIN="${AILANG_BIN:-ailang}"

# ── Resolved-binary announcement for legs 1-2 (charter item 9; iter-66) ───────
# ATTRIBUTION, NOT A GATE — this block cannot change the exit code, and that is deliberate.
#
# verify_go.sh:33 announces its resolved binary; this gate never did, and legs 1-2 default to bare
# PATH `ailang`. So on 2026-08-04 upstream `latest` moved v0.30.0 -> v0.33.0 and the repo's PRIMARY
# .ail gate silently began validating against an unpinned compiler, in violation of CLAUDE.md's
# "never a -dirty dev build" rule; two iterations (51, 52) ran that way and nothing surfaced it
# (charter item 9, iter-53). A recorded prediction of future breakage is not a monitor. This is.
#
# Leg 3 is NOT covered here: verify_world_package.sh:15 resolves its own WORLD_PKG_AILANG_BIN and
# already announces the compiler it pinned by exact bytes. Legs 1-2 are the unannounced half.
#
# The hard version ASSERTION and the CI `latest`->pinned-tag edit are deliberately NOT here. They
# are coupled and human-gated: a hard assert alone would red CI on the next upstream release, with
# no human present. Warning is the half that is safe to land headless.
GATE_PINNED_VERSION="v0.30.0"
_ail_resolved="$(command -v -- "$AILANG_BIN" 2>/dev/null || printf '%s' "$AILANG_BIN")"
_ail_version_line="$("$AILANG_BIN" --version 2>&1 | head -1)"
[ -n "$_ail_version_line" ] || _ail_version_line="<unavailable>"
# Exact token compare, never a substring: `v0.30.0-205-g54d6bd191-dirty` CONTAINS `v0.30.0`, so a
# substring test would grade a 205-commit dirty dev build as pinned.
_ail_version_tok="$(printf '%s\n' "$_ail_version_line" | awk '{print $2}')"
echo "── legs 1-2 AILANG_BIN=$_ail_resolved ($_ail_version_line)"
if [ "$_ail_version_tok" != "$GATE_PINNED_VERSION" ]; then
  echo "⚠ DRIFT: legs 1-2 are validating against '${_ail_version_tok:-<none>}', not the documented"
  echo "  pin $GATE_PINNED_VERSION. This gate does NOT fail on drift (the hard assertion is coupled to"
  echo "  the CI \`latest\`->pinned-tag edit and is human-gated, charter item 9). Export"
  echo "  AILANG_BIN=/path/to/$GATE_PINNED_VERSION/ailang to verify against the pin."
fi

command -v python3 >/dev/null 2>&1 || {
  echo "✗ python3 not found — the gate parses JSON with python3 and cannot run without it" >&2
  exit 1
}

# ── Hardcoded gate policy (NOT env-overridable) ───────────────────────────────
# The manifest is keyed by (repo-relative module file, bare function name): ai-check emits BARE
# names in verify.results[].function (V17). Sketches carry EMPTY required sets and are excluded from
# the total, so a future contracted sketch can neither mask a required identity nor perturb the total.
GATE_LEG_TIMEOUT_S=120   # wall-clock cap per ai-check module leg (Standing Rule 6, V26)
GATE_TEST_TIMEOUT_S=180  # wall-clock cap for the directory-mode test leg (V26)
export GATE_LEG_TIMEOUT_S GATE_TEST_TIMEOUT_S

# "base|tree": cd into `base`, then find/check .ail files under `tree` (relative to base).
ROOTS=(
  "design_docs|."
  ".|world"
)

# ── Leg 0 — bounded execution (hardcoded deadlines, non-env-overridable, V26) ──
# ai-check -timeout bounds only individual Z3 queries and `ailang test` has no timeout at all (V26),
# so a solver/runner/parse hang would block CI indefinitely. Every binary invocation in BOTH legs
# runs through this wrapper: start the child in its own process group and, on expiry, SIGKILL the
# WHOLE group, print a named TIMEOUT to STDERR, and exit 124. A 124 is the ONE fatal, non-advisory
# exit code (distinct from the Z3-error-exits-0 and counterexample-exits-1 classes the JSON parse
# handles). python3's start_new_session gives us the process group.
run_bounded() {  # $1=timeout_s  $2=out_file  $3..=cmd ;  exit 124 + named msg on expiry
  local t="$1" out="$2"; shift 2
  python3 - "$t" "$out" "$@" <<'PY'
import os, signal, subprocess, sys
t = int(sys.argv[1]); out = sys.argv[2]; cmd = sys.argv[3:]
with open(out, "wb") as f:
    p = subprocess.Popen(cmd, stdout=f, stderr=subprocess.STDOUT, start_new_session=True)
    try:
        sys.exit(p.wait(timeout=t))
    except subprocess.TimeoutExpired:
        os.killpg(os.getpgid(p.pid), signal.SIGKILL)
        sys.stderr.write("✗ TIMEOUT after %ds: %s\n" % (t, " ".join(cmd))); sys.exit(124)
PY
}

# Absolute temp paths: Leg 1 runs run_bounded inside `( cd "$base" ... )`, so the out-file path
# must not be relative to the changed cwd.
tmp_json="$(mktemp -t verify_ail_json.XXXXXX)"
tmp_test_json="$(mktemp -t verify_ail_test.XXXXXX)"
trap 'rm -f "$tmp_json" "$tmp_test_json"' EXIT

# ── Leg 1 — ai-check required-check manifest ──────────────────────────────────
# For each swept module: run bounded, capture JSON, parse. Exit codes ADVISORY (JSON is
# authoritative) EXCEPT a 124 TIMEOUT which is FATAL. The python parser routes its failure message
# to STDERR (so it survives the $() stdout capture) and prints only the module's verified count to
# STDOUT on success.
echo "── Leg 1: ai-check required-check manifest"
total_verified=0
checked=0
for entry in "${ROOTS[@]}"; do
  base="${entry%%|*}"
  tree="${entry#*|}"
  if [ "$tree" = "." ]; then searchdir="$base"; else searchdir="$base/$tree"; fi
  [ -d "$searchdir" ] || continue
  while IFS= read -r -d '' f; do
    rel="${f#"$base"/}"          # module path relative to its base (what ai-check wants)
    mod="${f#./}"                # repo-relative path (manifest key), normalized (gemini catch)
    checked=$((checked + 1))
    echo "   ai-check $mod"
    ( cd "$base" && run_bounded "$GATE_LEG_TIMEOUT_S" "$tmp_json" "$AILANG_BIN" ai-check -timeout 5s "$rel" )
    rc=$?
    if [ "$rc" -eq 124 ]; then
      echo "✗ ai-check TIMEOUT on $mod (>${GATE_LEG_TIMEOUT_S}s)" >&2
      exit 1
    fi
    # other exit codes advisory (V10/V20) — the JSON parse below is authoritative
    mod_verified=$(python3 - "$mod" "$tmp_json" <<'PY'
import json, sys
mod = sys.argv[1]
# Hardcoded manifest — keyed by (repo-relative module file, bare function name), V17.
REQUIRED_VERIFIED = {
    "world/transitions.ail": {"applyRevision"},
    "world/contracts.ail":   {"isValidNextWorld"},
    "world/logepoch.ail":    {"sameRef", "servesEntry"},
    "world/types.ail":       set(),
}
try:
    with open(sys.argv[2]) as fh:
        d = json.load(fh)
except Exception as e:
    sys.stderr.write("✗ %s: could not parse ai-check JSON (%s)\n" % (mod, e)); sys.exit(1)

check = d.get("check", {})
verify = d.get("verify", {})
if check.get("passed") is not True:
    sys.stderr.write("✗ %s: check.passed != true\n" % mod); sys.exit(1)
if verify.get("errors", 0) > 0:
    sys.stderr.write("✗ %s: verify.errors == %s (Z3 encoding error, silent under exit codes V10)\n"
                     % (mod, verify.get("errors"))); sys.exit(1)
if verify.get("counterexample", 0) > 0:
    sys.stderr.write("✗ %s: verify.counterexample == %s\n" % (mod, verify.get("counterexample"))); sys.exit(1)

# Required identities (world/ modules only; sketches carry empty sets).
required = REQUIRED_VERIFIED.get(mod, set())
by_fn = {r.get("function"): r.get("status") for r in verify.get("results", [])}
for fn in sorted(required):
    st = by_fn.get(fn)
    if st is None:
        sys.stderr.write("✗ %s: required identity (%s, %s) MISSING from verify.results[] "
                         "(vanished silently, V20)\n" % (mod, mod.split('/')[-1].replace('.ail',''), fn)); sys.exit(1)
    if st != "verified":
        sys.stderr.write("✗ %s: required identity (%s) status is '%s', expected 'verified'\n"
                         % (mod, fn, st)); sys.exit(1)

# Print the count of VERIFIED results (used for the world/-scoped secondary total).
print(sum(1 for r in verify.get("results", []) if r.get("status") == "verified"))
PY
    ) || exit 1
    case "$mod" in
      world/*) total_verified=$((total_verified + mod_verified));;
    esac
  done < <(find "$searchdir" -name '*.ail' -print0 | sort -z)
done

if [ "$checked" -eq 0 ]; then
  echo "✗ no .ail modules found — the gate would be vacuous; failing loudly" >&2
  exit 1
fi

EXACT_TOTAL_VERIFIED=4
if [ "$total_verified" -ne "$EXACT_TOTAL_VERIFIED" ]; then
  echo "✗ expected exactly $EXACT_TOTAL_VERIFIED proven world/ contracts, got $total_verified" >&2
  exit 1
fi
echo "   ✓ $total_verified/$EXACT_TOTAL_VERIFIED required world/ identities verified across $checked module(s)"

# ── Leg 2 — named inline tests (directory mode, V18/V22) ──────────────────────
# One bounded `test --format json world/` run. Names are preserved and merged across modules (V22).
# Exit codes advisory EXCEPT 124. Strip the human banner before the first '{' (V19). Assert every
# REQUIRED_TESTS name is present with status pass and failed_tests==0; len(tests[])==14 is secondary
# (NOT passed_tests — it also counts contract-derived PROPERTIES and is flaky, correction D-B).
# Property-test skips are tolerated (record-param generators absent, V14).
echo "── Leg 2: named inline tests"
run_bounded "$GATE_TEST_TIMEOUT_S" "$tmp_test_json" "$AILANG_BIN" test --format json world/
rc=$?
if [ "$rc" -eq 124 ]; then
  echo "✗ ailang test TIMEOUT (>${GATE_TEST_TIMEOUT_S}s)" >&2
  exit 1
fi
# other exit codes advisory — the JSON parse below is authoritative
python3 - "$tmp_test_json" <<'PY' || exit 1
import json, sys
REQUIRED_TESTS = {  # logepoch (8) + contracts (6)
    "renderRef_test_1", "renderRef_test_2", "sameRef_test_1", "sameRef_test_2",
    "cacheKey_test_1", "cacheKey_test_2", "servesEntry_test_1", "servesEntry_test_2",
    "proposalMatchesWorld_test_1", "proposalMatchesWorld_test_2",
    "verificationMatchesProposal_test_1", "verificationMatchesProposal_test_2",
    "commitAllowed_test_1", "commitAllowed_test_2",
}
EXACT_TOTAL_TESTS = 14
raw = open(sys.argv[1], "rb").read().decode("utf-8", "replace")
i = raw.find("{")                       # strip stdout banner before the first '{' (V19)
if i < 0:
    sys.stderr.write("✗ test leg: no JSON object in output (crashed or empty run)\n"); sys.exit(1)
try:
    d = json.loads(raw[i:])
except Exception as e:
    sys.stderr.write("✗ test leg: could not parse test JSON (%s)\n" % e); sys.exit(1)

tests = d.get("tests", [])
by_name = {t.get("name"): t.get("status") for t in tests}
missing_or_failing = []
for name in sorted(REQUIRED_TESTS):
    st = by_name.get(name)
    if st != "pass":
        missing_or_failing.append("%s=%s" % (name, "MISSING" if st is None else st))
if missing_or_failing:
    sys.stderr.write("✗ test leg: required named tests missing/failing: %s\n"
                     % ", ".join(missing_or_failing)); sys.exit(1)

failed = d.get("failed_tests", 0)
if failed and failed > 0:
    sys.stderr.write("✗ test leg: failed_tests == %s\n" % failed); sys.exit(1)

n = len(tests)
if n != EXACT_TOTAL_TESTS:
    sys.stderr.write("✗ test leg: expected exactly %d named tests[], got %d (secondary check)\n"
                     % (EXACT_TOTAL_TESTS, n)); sys.exit(1)

print("   ✓ all %d required named tests pass (failed_tests=0)" % EXACT_TOTAL_TESTS)
PY

echo "── Leg 3: world package nine-step gate"
./scripts/verify_world_package.sh || exit $?

echo "✓ verify gate PASSED: 4 required identities verified, 14 named tests pass"
