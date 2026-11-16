#!/usr/bin/env bash
# Row 260 (w-ailang-v0604-contract-eq) sprint gates — plan §5.
# usage: bash gates.sh <tree> <logdir> <m1|m2>
#   m1 = only M1 applied: the `UNRECOGNISED RELEASE ... 'v0.60.4'` notice MUST print (AC1).
#   m2 = M1+M2 applied:   the notice MUST be absent.
# env: GATES_SKIP_VERIFY_GO=1 skips the last leg (full ./scripts/verify_go.sh plain); default runs it.
# macOS /bin/bash 3.2 safe (no declare -A). No `set -e`: every leg runs and records its own rc,
# captured from the command itself, never through a pipe. Exits 1 if any leg failed.
set -uo pipefail
T=${1:?tree}; L=${2:?logdir}; MODE=${3:?m1|m2}
case "$MODE" in m1|m2) ;; *) echo "mode must be m1 or m2" >&2; exit 2 ;; esac
BASE=cbcf80fb
V604=${V604:-/tmp/ailang-v0604-i393/ailang}
PIN=${PIN:-$HOME/.pinned-ailang/ailang}
WANT_CH='sha256:359a82fde07d8e9004786448da37ea3fb07cdc256682efc0fa57240df3c4a72c'
mkdir -p "$L"; cd "$T" || exit 2
L=$(cd "$L" && pwd)
: > "$L/legs.txt"
FAILS=""
s0=$(date +%s)
note() { echo "$*" | tee -a "$L/legs.txt"; }
result() { # result <leg> <rc> <secs> [detail]
  note "LEG $1 rc=$2 $3s ${4:-}"
  [ "$2" -eq 0 ] || FAILS="$FAILS $1"
}
"$V604" --version 2>/dev/null | head -1 | grep -qx 'AILANG v0.60.4' || { echo "V604 is not AILANG v0.60.4: $V604" >&2; exit 2; }
"$PIN" --version 2>/dev/null | head -1 | grep -qx 'AILANG v0.41.0' || { echo "PIN is not AILANG v0.41.0: $PIN" >&2; exit 2; }
note "tree=$T head=$(git rev-parse --short HEAD) mode=$MODE base=$BASE dirty_files=$(git status --porcelain | wc -l | tr -d ' ')"

# ---- AC1 / AC2: verify_ail.sh on the CI lane (v0.60.4) and the pin lane (v0.41.0)
verify_ail_leg() { # verify_ail_leg <leg> <legs-1-2 binary> <expect-notice yes|no>
  local leg=$1 bin=$2 want_notice=$3 s rc miss=""
  s=$(date +%s)
  AILANG_BIN="$bin" WORLD_PKG_AILANG_BIN="$PIN" env -u AILANG_REGISTRY_API_KEY ./scripts/verify_ail.sh > "$L/$leg.log" 2>&1
  rc=$?
  for pat in '✓ 16/16 required world/ identities verified across 24 module(s)' \
             '✓ all 138 execargs sketch named tests pass' \
             '✓ all 40 required named tests pass (failed_tests=0)' \
             '✓ all 503 se-tools named tests pass across 10 modules' \
             '✓ world package gate PASSED: 9/9 steps performed non-zero work' \
             'PUB011 ratchet passed: uncontracted_exports=3'; do
    grep -qF -- "$pat" "$L/$leg.log" || miss="$miss [$pat]"
  done
  if grep -qF "UNRECOGNISED RELEASE: legs 1-2 resolved 'v0.60.4'" "$L/$leg.log"; then got=yes; else got=no; fi
  [ "$got" = "$want_notice" ] || miss="$miss [v0.60.4 notice: want=$want_notice got=$got]"
  [ -z "$miss" ] || { echo "missing:$miss" >> "$L/$leg.log"; [ "$rc" -ne 0 ] || rc=97; }
  result "$leg" "$rc" $(( $(date +%s) - s )) "notice=$got${miss:+ MISSING:$miss}"
}
if [ "$MODE" = m1 ]; then NOTICE604=yes; else NOTICE604=no; fi
verify_ail_leg AC1_verify_ail_v0604 "$V604" "$NOTICE604"
verify_ail_leg AC2_verify_ail_pin "$PIN" no

# ---- AC3: per-file ai-check on both CI compilers: passed=true, verified 6/5/3/2, no other status
s=$(date +%s); : > "$L/AC3.log"; rc=0
for bin in "$V604" "$PIN"; do
  for pair in types:6 contracts:5 transitions:3 logepoch:2; do
    f="world/${pair%%:*}.ail"; want=${pair##*:}
    "$bin" ai-check -timeout 5s "$f" > "$L/AC3.raw.json" 2>/dev/null
    python3 - "$L/AC3.raw.json" "$want" "$("$bin" --version | head -1)" "$f" >> "$L/AC3.log" 2>&1 <<'PY'
import json, sys
p, want, ver, f = sys.argv[1], int(sys.argv[2]), sys.argv[3], sys.argv[4]
try:
    d = json.load(open(p))
except Exception as e:
    print(f"{ver} {f} PARSE_FAIL {e}"); sys.exit(1)
c = d.get("check", {}); res = (d.get("verify") or {}).get("results") or []
st = sorted({r.get("status") for r in res})
n = sum(1 for r in res if r.get("status") == "verified")
ok = c.get("passed") is True and n == want and st == ["verified"] and len(res) == want
print(f"{ver} {f} passed={c.get('passed')} verified={n}/{want} statuses={st} {'OK' if ok else 'FAIL'}")
sys.exit(0 if ok else 1)
PY
    [ $? -eq 0 ] || rc=1
  done
done
rm -f "$L/AC3.raw.json"
result AC3_per_file_counts "$rc" $(( $(date +%s) - s )) "$(grep -c ' OK$' "$L/AC3.log")/8 OK"

# ---- AC4: kernel contract history controls (CI step "Kernel contract and PUB011 mutation controls")
s=$(date +%s)
AILANG_BIN="$PIN" python3 scripts/test_kernel_contracts.py > "$L/AC4.log" 2>&1
rc=$?
grep -qF 'world contract history controls PASSED' "$L/AC4.log" || { [ "$rc" -ne 0 ] || rc=97; }
result AC4_kernel_contracts "$rc" $(( $(date +%s) - s ))

# ---- gofmt on the Go files this sprint changed (vs base)
s=$(date +%s)
CHANGED_GO=$( { git diff --name-only "$BASE" -- '*.go'; git ls-files --others --exclude-standard -- '*.go'; } | sort -u)
out=$( [ -z "$CHANGED_GO" ] || gofmt -l $CHANGED_GO )
rc=0; [ -z "$out" ] || rc=1
printf 'changed go files:\n%s\ngofmt -l:\n%s\n' "$CHANGED_GO" "$out" > "$L/gofmt.log"
result gofmt "$rc" $(( $(date +%s) - s )) "files=$(echo "$CHANGED_GO" | grep -c .)"

# ---- AC5: Go build, vet, seven packages (host/broker re-run ONCE on failure, per design AC5)
s=$(date +%s); go build ./... > "$L/AC5_build.log" 2>&1; rc=$?
result AC5_go_build "$rc" $(( $(date +%s) - s ))
s=$(date +%s); go vet ./cmd/world-publish/... ./host/pkgproj/... ./host/runbook/... > "$L/AC5_vet.log" 2>&1; rc=$?
result AC5_go_vet "$rc" $(( $(date +%s) - s ))
s=$(date +%s); rc=0; okc=0
for pkg in ./cmd/world-publish/... ./host/broker/... ./host/pkgproj/... ./host/runbook/... ./host/replay/... ./host/verifygate/... ./host/pins/...; do
  name=$(echo "$pkg" | tr '/.' '__')
  AILANG_BIN="$PIN" WORLD_PKG_AILANG_BIN="$PIN" env -u AILANG_REGISTRY_API_KEY go test -count=1 "$pkg" > "$L/AC5_test$name.log" 2>&1
  prc=$?
  if [ "$prc" -ne 0 ] && [ "$pkg" = ./host/broker/... ]; then
    note "   host/broker FAILED once (rc=$prc; failing tests: $(grep -E '^--- FAIL' "$L/AC5_test$name.log" | awk '{print $3}' | tr '\n' ' ')) — RE-RUNNING ONCE as design AC5 allows (known load flakes)"
    mv "$L/AC5_test$name.log" "$L/AC5_test$name.first.log"
    AILANG_BIN="$PIN" WORLD_PKG_AILANG_BIN="$PIN" env -u AILANG_REGISTRY_API_KEY go test -count=1 "$pkg" > "$L/AC5_test$name.log" 2>&1
    prc=$?
    note "   host/broker re-run rc=$prc"
  fi
  if grep -qE '^(FAIL|--- FAIL)' "$L/AC5_test$name.log"; then [ "$prc" -ne 0 ] || prc=97; fi
  [ "$prc" -eq 0 ] && okc=$((okc+1)) || rc=1
  note "   $(grep -E '^(ok|FAIL|---)' "$L/AC5_test$name.log" | tail -1)"
done
result AC5_go_test_7pkgs "$rc" $(( $(date +%s) - s )) "$okc/7 ok"

# ---- AC6: identity
s=$(date +%s); : > "$L/AC6.log"
git show "$BASE:scripts/world_package_ready_packet.golden.json" > "$L/AC6.base-golden.json"
(cd packages/world-core && env -u AILANG_REGISTRY_API_KEY "$PIN" publish --dry-run) > "$L/AC6.dryrun.txt" 2>&1
drc=$?
python3 - scripts/world_package_ready_packet.golden.json "$L/AC6.base-golden.json" "$WANT_CH" packages/world-core/ailang.toml "$L/AC6.dryrun.txt" "$drc" >> "$L/AC6.log" 2>&1 <<'PY'
import json, sys, tomllib
g = json.load(open(sys.argv[1])); b = json.load(open(sys.argv[2])); want_ch = sys.argv[3]
m = tomllib.load(open(sys.argv[4], "rb")); dry = open(sys.argv[5], encoding="utf-8").read(); drc = int(sys.argv[6])
checks = [
  ("contentHash == design value", g["contentHash"] == want_ch, g["contentHash"]),
  ("interfaceHash == base", g["interfaceHash"] == b["interfaceHash"], g["interfaceHash"]),
  ("interfaceHashV2 == base", g["interfaceHashV2"] == b["interfaceHashV2"], g["interfaceHashV2"]),
  ("golden version 0.1.2", g["version"] == "0.1.2", g["version"]),
  ("golden compilerVersion unchanged", g["compilerVersion"] == b["compilerVersion"], g["compilerVersion"]),
  ("manifest version 0.1.2", m["package"]["version"] == "0.1.2", m["package"]["version"]),
  ("manifest [release] kind feature", m["release"]["kind"] == "feature", m["release"]["kind"]),
  ("dry-run rc 0", drc == 0, drc),
  ("dry-run Publishing world/core@0.1.2...", "Publishing world/core@0.1.2..." in dry, ""),
  ("dry-run contracts 16/16 verified", "contracts: ✓ 16/16 verified" in dry, ""),
  ("dry-run changelog_section true", "changelog_section true" in dry, ""),
  ("dry-run no gates", "✓ no gates" in dry, ""),
  ("dry-run full contentHash equals golden", g["contentHash"][len("sha256:"):][:17] in dry, ""),
]
bad = 0
for name, ok, val in checks:
    print(f"{'OK  ' if ok else 'FAIL'} {name} {val}")
    bad += (not ok)
sys.exit(1 if bad else 0)
PY
rc=$?
result AC6_identity "$rc" $(( $(date +%s) - s )) "$(grep -c '^OK' "$L/AC6.log")/13 OK"

# ---- AC7: history is not rewritten
s=$(date +%s); : > "$L/AC7.log"; rc=0
d=$(git diff "$BASE" -- host/broker/testdata host/replay/testdata)
[ -z "$d" ] && echo "OK   host/broker/testdata + host/replay/testdata unchanged vs $BASE" >> "$L/AC7.log" || { echo "FAIL testdata changed" >> "$L/AC7.log"; rc=1; }
git show "$BASE:packages/world-core/CHANGELOG.md" | sed -n '/^## 0\.1\.1/,$p' > "$L/AC7.base-tail.md"
sed -n '/^## 0\.1\.1/,$p' packages/world-core/CHANGELOG.md > "$L/AC7.cur-tail.md"
if [ -s "$L/AC7.base-tail.md" ] && cmp -s "$L/AC7.base-tail.md" "$L/AC7.cur-tail.md"; then
  echo "OK   CHANGELOG ## 0.1.1 and ## 0.1.0 byte-identical ($(wc -c < "$L/AC7.cur-tail.md" | tr -d ' ') bytes)" >> "$L/AC7.log"
else echo "FAIL CHANGELOG history sections differ" >> "$L/AC7.log"; rc=1; fi
first=$(grep -m1 '^## ' packages/world-core/CHANGELOG.md)
case "$first" in "## 0.1.2 — "*) echo "OK   newest section: $first" >> "$L/AC7.log" ;; *) echo "FAIL newest section: $first" >> "$L/AC7.log"; rc=1 ;; esac
result AC7_history "$rc" $(( $(date +%s) - s ))

# ---- M2 only: the two JSON copies byte-identical, the v0.60.4 row present
if [ "$MODE" = m2 ]; then
  s=$(date +%s); rc=0
  cmp -s pins/ailang-releases.json host/pins/ailang-releases.json || rc=1
  python3 -c 'import json,sys; m=json.load(open("pins/ailang-releases.json")); r=[x for x in m["releases"] if x["release"]=="v0.60.4"]; sys.exit(0 if len(r)==1 and [y["role"] for y in r[0]["roles"]]==["observed"] and "assets" not in r[0] else 1)' || rc=1
  grep -qx 'v0.60.4' scripts/testdata/ailang_release_observed.txt || rc=1
  result M2_pins_row "$rc" $(( $(date +%s) - s ))
fi

# ---- the repo's full Go gate (CI step "go build + test gate"): verify_go.sh, plain leg
if [ "${GATES_SKIP_VERIFY_GO:-0}" = 1 ]; then
  note "LEG verify_go SKIPPED (GATES_SKIP_VERIFY_GO=1)"
else
  s=$(date +%s)
  VERIFY_GO_LEGS=plain AILANG_BIN="$PIN" WORLD_PKG_AILANG_BIN="$PIN" env -u AILANG_REGISTRY_API_KEY ./scripts/verify_go.sh > "$L/verify_go.log" 2>&1
  rc=$?
  grep -qF '✓ go gate PASSED (plain)' "$L/verify_go.log" || { [ "$rc" -ne 0 ] || rc=97; }
  result verify_go_plain "$rc" $(( $(date +%s) - s )) "$(grep -E '^(FAIL|--- FAIL)' "$L/verify_go.log" | head -5 | tr '\n' ' ')"
fi

note "GATES_DONE mode=$MODE $(( $(date +%s) - s0 ))s FAILS:${FAILS:- none}"
[ -z "$FAILS" ]
