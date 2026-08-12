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

# ── Released-binary assertion and announcement for legs 1-2 ──────────────────
_ail_release_verdict() { # $1=version token; prints one stable verdict code
  local tok="$1"
  if [ -z "$tok" ]; then
    printf '%s\n' NO_VERSION_TOKEN
  elif [[ "$tok" =~ (-dirty$|-[0-9]+-g[0-9a-f]+) ]]; then
    printf '%s\n' DEV_BUILD
  elif ! [[ "$tok" =~ ^v[0-9]+\.[0-9]+\.[0-9]+(-[A-Za-z0-9.]+)?$ ]]; then
    printf '%s\n' NOT_A_RELEASE
  else
    printf '%s\n' OK
  fi
}

_refuse() { # $1=stable code; $2=detail
  echo "✗ AILANG_BIN refused [$1]: $2" >&2
  exit 1
}

# Known-positive control: if these disagree, no observation made by this instrument is trustworthy.
while IFS='|' read -r _control_token _control_expected; do
  _control_actual="$(_ail_release_verdict "$_control_token")"
  if [ "$_control_actual" != "$_control_expected" ]; then
    echo "verify_ail.sh: FATAL: release-verdict instrument broken: token '${_control_token:-<empty>}' expected $_control_expected, got $_control_actual" >&2
    exit 1
  fi
done <<'CONTROL'
v0.30.0|OK
v0.30.0-205-g54d6bd191-dirty|DEV_BUILD
|NO_VERSION_TOKEN
CONTROL

if ! _ail_resolved="$(command -v -- "$AILANG_BIN" 2>/dev/null)"; then
  _refuse UNRESOLVABLE "'$AILANG_BIN' is not an executable command"
fi
_ail_version_line="$("$AILANG_BIN" --version 2>&1 | head -1)"
if [ -z "$_ail_version_line" ]; then
  _refuse NO_VERSION_OUTPUT "--version produced no first line"
fi
_ail_version_tok="$(printf '%s\n' "$_ail_version_line" | awk '{print $2}')"
if [ -z "$_ail_version_tok" ]; then
  _refuse NO_VERSION_TOKEN "first --version line has no version token: $_ail_version_line"
fi
_ail_verdict="$(_ail_release_verdict "$_ail_version_tok")"
if [ "$_ail_verdict" = DEV_BUILD ]; then
  _refuse DEV_BUILD "version token '$_ail_version_tok' identifies a development build"
fi
if [ "$_ail_verdict" = NOT_A_RELEASE ]; then
  _refuse NOT_A_RELEASE "version token '$_ail_version_tok' is not an upstream release"
fi

echo "── legs 1-2 AILANG_BIN=$_ail_resolved ($_ail_version_line)"
# Expected-release notice (charter 9/CF-A-2). MEMBERSHIP, not equality against one line.
#
# The warning this replaces compared against the v0.30.0 pin and so fired on EVERY CI run once the
# `9/OD-10` ACCEPT ruling made legs 1-2 track releases — and a warning that always fires is not a
# signal. An equality test against a single recorded observation just MOVES that defect: legs 1-2
# legitimately resolve TWO different releases depending on lane, because CI's job 1 exports no
# AILANG_BIN (ci.yml:87 exports only WORLD_PKG_AILANG_BIN) and takes `releases/latest` off PATH,
# while a local operator exports the documented v0.30.0 pin. Under equality the local lane then
# printed `moved from 'v0.33.0' to 'v0.30.0'` on every run — always-firing again, and false: the
# operator's deliberate pin is not an upstream move.
#
# So the file lists the releases legs 1-2 are EXPECTED to resolve, and the notice fires only on a
# token outside that set. Quiet in both real lanes today; fires once when upstream's `latest` moves,
# which is exactly the 2026-08-04 v0.30.0 -> v0.33.0 event that created charter item 9 and ran
# unnoticed for two iterations. Non-fatal by construction — it can never red CI.
_ail_expected_releases="$(grep -vE '^[[:space:]]*(#|$)' scripts/testdata/ailang_release_observed.txt)"
[ -n "$_ail_expected_releases" ] || {
  echo "verify_ail.sh: FATAL: scripts/testdata/ailang_release_observed.txt lists no releases — an empty set would silence this notice forever rather than firing it." >&2
  exit 1
}
if ! printf '%s\n' "$_ail_expected_releases" | grep -qxF "$_ail_version_tok"; then
  echo "ℹ UNRECOGNISED RELEASE: legs 1-2 resolved '$_ail_version_tok', which is not among the"
  echo "  releases this repo expects ($(printf '%s' "$_ail_expected_releases" | tr '\n' ' ')). The primary"
  echo "  .ail gate has switched compiler — review it, then update scripts/testdata/ailang_release_observed.txt."
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

# Exact Leg-1 module manifest (identity allowlist, NOT a count). An intentional module
# add/remove is a ONE-LINE edit here, in the SAME commit. Repo-relative paths, matching
# the sweep's $mod key.
LEG1_MODULES=(
  design_docs/sketches/effectbroker.ail
  design_docs/sketches/logepoch.ail
  design_docs/sketches/storejournal.ail
  design_docs/sketches/transitions.ail
  design_docs/sketches/worlddapi.ail
  design_docs/sketches/worldkernel.ail
  design_docs/sketches/worldtypes.ail
  world/contracts.ail
  world/logepoch.ail
  world/transitions.ail
  world/types.ail
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

# NUL-aware diagnostic formatter. Reads a NUL-delimited set file and prints ONE
# shell-quoted token per line, so `diff -u` stays line-oriented and a pathological path
# (newline, |, space, glob char) is RENDERED safely rather than PARSED.
_nul_quoted_lines() { # $1=NUL-delimited file
  local _p
  while IFS= read -r -d '' _p; do printf '%q\n' "$_p"; done < "$1"
}

# Absolute temp paths: Leg 1 runs run_bounded inside `( cd "$base" ... )`, so the out-file path
# must not be relative to the changed cwd.
tmp_json="$(mktemp -t verify_ail_json.XXXXXX)"
tmp_test_json="$(mktemp -t verify_ail_test.XXXXXX)"
tmp_mods_actual="$(mktemp -t verify_ail_mods_a.XXXXXX)"
tmp_mods_expected="$(mktemp -t verify_ail_mods_e.XXXXXX)"
trap 'rm -f "$tmp_json" "$tmp_test_json" "$tmp_mods_actual" "$tmp_mods_expected"' EXIT

# ── Leg 1 — ai-check required-check manifest ──────────────────────────────────
# For each swept module: run bounded, capture JSON, parse. Exit codes ADVISORY (JSON is
# authoritative) EXCEPT a 124 TIMEOUT which is FATAL. The python parser routes its failure message
# to STDERR (so it survives the $() stdout capture) and prints only the module's verified count to
# STDOUT on success.
echo "── Leg 1: ai-check required-check manifest"
total_verified=0
checked=0

# ── Leg 1a — enumerate ONCE into parallel indexed arrays (NUL end to end; no delimiter
# is ever embedded in a record, so no path can be mis-parsed by the gate that must reject it).
bases=(); rels=(); mods=()
for entry in "${ROOTS[@]}"; do
  base="${entry%%|*}"
  tree="${entry#*|}"
  if [ "$tree" = "." ]; then searchdir="$base"; else searchdir="$base/$tree"; fi
  [ -d "$searchdir" ] || continue
  while IFS= read -r -d '' f; do
    bases+=("$base")
    rels+=("${f#"$base"/}")      # module path relative to its base (what ai-check wants)
    mods+=("${f#./}")            # repo-relative path (manifest key), normalized (gemini catch)
  done < <(find "$searchdir" -name '*.ail' -print0 | sort -z)
done

# ── Leg 1b — MEMBERSHIP COMPARE, before any ai-check runs.
# Guards come BEFORE the printf writes and read ${#arr[@]}: under `set -u` (line 30) with
# bash 3.2 (the rig's /usr/bin/env bash), "${arr[@]}" on an EMPTY array is an unbound-variable
# ABORT, so a write-then-guard order would kill the script before its own null-case message
# could ever print. ${#arr[@]} is set -u-safe on an empty array.
if [ "${#mods[@]}" -eq 0 ]; then
  echo "✗ swept .ail enumeration was empty — the membership compare would pass vacuously; failing loudly" >&2
  exit 1
fi
if [ "${#LEG1_MODULES[@]}" -eq 0 ]; then
  echo "✗ LEG1_MODULES allowlist is empty — the membership compare would pass vacuously; failing loudly" >&2
  exit 1
fi
printf '%s\0' "${mods[@]}"         | LC_ALL=C sort -z > "$tmp_mods_actual"
printf '%s\0' "${LEG1_MODULES[@]}" | LC_ALL=C sort -z > "$tmp_mods_expected"
if ! cmp -s "$tmp_mods_expected" "$tmp_mods_actual"; then
  echo "✗ swept .ail module set differs from the LEG1_MODULES allowlist — an intentional" >&2
  echo "  module add/remove must edit LEG1_MODULES in scripts/verify_ail.sh in the SAME commit" >&2
  diff -u --label "expected: LEG1_MODULES in scripts/verify_ail.sh" \
          --label "actual:   .ail files swept under ROOTS" \
          <(_nul_quoted_lines "$tmp_mods_expected") <(_nul_quoted_lines "$tmp_mods_actual") >&2
  exit 1
fi
echo "   ✓ swept .ail module set equals the LEG1_MODULES allowlist (${#mods[@]} modules)"

# ── Leg 1c — consume the SAME arrays by index. cwd / run_bounded / absolute-temp semantics
# are untouched; "what is compared" and "what is checked" are the same array objects.
for i in "${!mods[@]}"; do
  base="${bases[$i]}"
  rel="${rels[$i]}"
  mod="${mods[$i]}"
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
