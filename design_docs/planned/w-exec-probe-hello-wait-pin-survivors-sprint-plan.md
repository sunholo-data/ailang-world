# Sprint plan: w-exec-probe-hello-wait-pin-survivors (iteration 339, row 241)

Authority: `design_docs/planned/w-exec-probe-hello-wait-pin-survivors.md` r1, all sections, particularly §2–§4, §6 and §7; design quorum r1 PROCEED 3/3. Status: planned; outside-sandbox acceptance remains for the controller. One milestone M1, test-only; estimate ~0.15d, one executor run. No implementation is performed by this planner.

Planner HEAD `82305a9e633a104d9597aebc5ba03c2cbe4bf2b9`; code/design base `dcb2366e`. No git add, commit, stash, checkout, branch or other git metadata writes. The controller owns commits. Apply operations below change scratch/executor files only.

M1 installs the exact banked T12 and T10/helloConn changes and the one CI-list addition. T12 counts real dial cadence against 20 reference periods (gate 6), floors recorded gaps at the production constant, and caps timestamps at 80. T10 compares to setAt, removing the stop hand-off lag. Keep every design byte, including comments. The accepted 2×–4× band and survivor-2/4 equivalences remain accepted; no production seam or stronger gate is authorized.

## M1 scope and exact payload

Only `host/broker/exec_probe_test.go` (+59/-4) and `.github/workflows/ci.yml` (1/1) may change outside design_docs. Append ` TestExecProbeReversePollCadence` immediately after `TestExecProbeReverseHelloWaitBounds` in the tests= line of the “Row 140 Workspace.Exec real-srt tests, verbose (linux bubblewrap)” step. Do not change another test, any production byte in `host/broker/exec_probe.go`, any .ail file, go.mod, any other workflow or row 242's net6 work. Mutation edits live exclusively in scratch copies and must never reach the executor tree.

| Artifact | sha256 |
|---|---|
| plan-M1.diff.txt | b75e4cf0517633d71e70c88468a93748e58c95cfedfe9ef257d244b4dde146de |
| designer-proto_test.go.patch.txt (patch bytes) | 43221064d23bebcf830e38d261eb1a9d0c6153c1106dc409617dee7b4185f040 |
| patched host/broker/exec_probe_test.go | bbfccede19236d17b81d40401cc079aeb9987ff356c6a08affe6bc2a159e80d0 |
| patched .github/workflows/ci.yml | a11ff4ab05a41c9864299edb8f636675e3abb5dcdcba974d43fb6a405940710f |
| unchanged host/broker/exec_probe.go | 96f33172667d27defa0c52735d5cff1ad74d4e907ab00a4c056c4e5a26f698a4 |
| plan-mutants.sh.txt | 6854f0190ee5361b666f1674e9c9316e7b97462f79b7ed3d2c2ee555dcd03f5d |

The design header's `bbfccede…` belongs to the PATCHED FILE, not the named patch. The actual patch hash is `43221064…`; planner measured both. Never verify the patch using the file hash.

Run the following from the executor repository root in one `/bin/bash` session. The shell is macOS Bash 3.2. No associative arrays or empty-array expansion under set -u; no `timeout` binary. Python 3, Perl, pinned AILANG and the existing Go caches are required. Every broker/full-suite go test inherits `AILANG_BIN=$HOME/.pinned-ailang/ailang`. Do not replace a failed gate or edit the gate to manufacture green.

```bash
set -eu
set -o pipefail
WT=$(pwd -P)
X=$(mktemp -d /tmp/world-iter339-executor.XXXXXX)
export AILANG_BIN=$HOME/.pinned-ailang/ailang
D=design_docs/verification/world-iter339/plan-M1.diff.txt
python3 - "$D" <<'HASH'
from pathlib import Path
import hashlib,sys
assert hashlib.sha256(Path(sys.argv[1]).read_bytes()).hexdigest() == 'b75e4cf0517633d71e70c88468a93748e58c95cfedfe9ef257d244b4dde146de'
HASH
# Refuse unrelated code drift, including pre-existing changes, before applying.
git diff --quiet dcb2366e -- . ':!design_docs'
python3 - <<'BASE'
from pathlib import Path
import hashlib
expected={'.github/workflows/ci.yml': 'daf57b6e168f616cd99e370df8ded5c853df72207c3cc5ec248110702c0ecead', 'host/broker/exec_probe_test.go': 'e35784db17623df899e573096edf41ea4b205dbb4975701380594082e8181c5e', 'host/broker/exec_probe.go': '96f33172667d27defa0c52735d5cff1ad74d4e907ab00a4c056c4e5a26f698a4'}
for p,h in expected.items(): assert hashlib.sha256(Path(p).read_bytes()).hexdigest()==h,p
BASE
git apply --check "$D"
git apply "$D"
python3 - <<'POST'
from pathlib import Path
import hashlib
expected={'.github/workflows/ci.yml': 'a11ff4ab05a41c9864299edb8f636675e3abb5dcdcba974d43fb6a405940710f', 'host/broker/exec_probe_test.go': 'bbfccede19236d17b81d40401cc079aeb9987ff356c6a08affe6bc2a159e80d0', 'host/broker/exec_probe.go': '96f33172667d27defa0c52735d5cff1ad74d4e907ab00a4c056c4e5a26f698a4'}
for p,h in expected.items(): assert hashlib.sha256(Path(p).read_bytes()).hexdigest()==h,p
print('post hashes OK')
POST
```


## Acceptance commands (all required)

Run sequentially from WT; save each command's stdout/stderr and rc under X. An alarm/timeout is an instrument result, never a passing gate. Every explicit Go test uses -timeout 30m; the 2100-second outer alarm leaves time to emit Go's timeout trace. Preserve diagnostics and stop M1 for controller triage on a load-bearing gate mismatch. The mutation round deliberately continues to collect all rows and returns nonzero if any mismatch remains.

**A1 focused**

```bash
perl -e 'alarm 2100; exec @ARGV' go test -timeout 30m -count=1 -v -run '^(TestExecProbeReversePollCadence|TestExecProbeReverseHelloWaitBounds)$' ./host/broker/
```

Expected: rc 0; exactly 2 top-level PASS, T10 slow-child and slow-mux PASS (4 PASS at all levels), one CADENCE line, 0 FAIL/SKIP; dials >= 6; fastest >= execProbePollEvery = 50 ms.

**A2 package**

```bash
perl -e 'alarm 2100; exec @ARGV' go test -timeout 30m -count=1 -v ./host/broker/
```

Expected: rc 0; 336 top-level PASS, 30 top-level SKIP, 0 FAIL (design V15). -v is added to the §4 step-2 command so the census is observable.

**A3 mutation**

```bash
/bin/bash "$WT/design_docs/verification/world-iter339/plan-mutants.sh.txt" "$WT" "$X/mutants" > "$X/mutation-round.log" 2>&1
```

Expected: rc 0; all 61 execution rows match their exact rc and FAIL set, all 59 mutated executions restored(cmp ok), no site count other than 1. 58 distinct configurations: 29 production + 29 test/pairs; controls C0/C1 and duplicate M1-broad account for remaining executions. See full table and expected.json generated by script. A socket-denied row has no verdict and must be rerun outside sandbox.

**A4 HEAVY-LOAD**

```bash
/bin/bash "$X/heavy-load.sh" "$X/heavy-load.log" > "$X/heavy-load-summary.log" 2>&1
```

Expected: CONTROLLER-RUN OUTSIDE SANDBOX, rc 0; yes-alive-before=160, yes-alive-after=0; 100 T12 and 100 T10 top-level PASS, 100 CADENCE samples, minimum dials >= 6, smallest fastest >= 50 ms; uptime before and after. Use the exact script below, replacing design §4 step 4 (32 yes).

**A5a format**

```bash
perl -e 'alarm 600; exec @ARGV' gofmt -l host/broker
```

Expected: rc 0; no output.

**A5b vet**

```bash
perl -e 'alarm 600; exec @ARGV' go vet ./host/broker/
```

Expected: rc 0; no output.

**A6a build and full suite**

```bash
perl -e 'alarm 600; exec @ARGV' go build ./... && perl -e 'alarm 2100; exec @ARGV' go test -timeout 30m ./...
```

Expected: both rc 0; no FAIL. Full suite has explicit -timeout 30m; never rely on the default 10m timeout.

**A6b AILANG**

```bash
perl -e 'alarm 1200; exec @ARGV' ./scripts/verify_ail.sh
```

Expected: rc 0; final verify gate PASSED banner, all required manifest identities and named tests accepted.

**A7 shape and hashes**

```bash
git ls-files --others --exclude-standard -- . ':!design_docs' > "$X/untracked.txt"
[ ! -s "$X/untracked.txt" ]
git diff --name-only dcb2366e -- . ':!design_docs' > "$X/shape.txt"
python3 - "$X/shape.txt" <<'SHAPE'
from pathlib import Path
import sys
assert Path(sys.argv[1]).read_text().splitlines()==['.github/workflows/ci.yml','host/broker/exec_probe_test.go']
SHAPE
git diff --quiet dcb2366e -- host/broker/exec_probe.go
git show dcb2366e:host/broker/exec_probe.go > "$X/base-exec_probe.go"
cmp "$X/base-exec_probe.go" host/broker/exec_probe.go
git diff --numstat dcb2366e -- .github/workflows/ci.yml
# Exactly: 1<TAB>1<TAB>.github/workflows/ci.yml
git diff dcb2366e -- .github/workflows/ci.yml
# One changed tests= line in the Row 140 step, only the required appended name.
python3 - <<'POSTCHECK'
from pathlib import Path
import hashlib
expected={'.github/workflows/ci.yml': 'a11ff4ab05a41c9864299edb8f636675e3abb5dcdcba974d43fb6a405940710f', 'host/broker/exec_probe_test.go': 'bbfccede19236d17b81d40401cc079aeb9987ff356c6a08affe6bc2a159e80d0', 'host/broker/exec_probe.go': '96f33172667d27defa0c52735d5cff1ad74d4e907ab00a4c056c4e5a26f698a4'}
for f,h in expected.items(): assert hashlib.sha256(Path(f).read_bytes()).hexdigest()==h,f
print('post hashes OK')
POSTCHECK
```

Expected: rc 0; exactly two named paths, production cmp identical, CI numstat 1/1, all exact post hashes.



## HEAVY-LOAD: controller only, outside sandbox

Write the following exact block to `$X/heavy-load.sh` and run A4 above. It starts its own 160 spinning yes processes, verifies all are alive, kills only those PIDs via a trap, polls their termination with a bounded deadline, reaps them and asserts 0 alive. A Go alarm cannot strand the yes processes. No unrelated ambient process is touched. Record the two uptime readings; 160 added spinners are a repeatable load condition, not a claim that the historical other-session 150–214 ambient workload is identical. No GOMAXPROCS override is applied to this requested gate. Do not run this load gate in the planner sandbox.

```bash
#!/bin/bash
# Controller only, outside the sandbox. CWD is the post-M1 executor tree.
set -eu
set -o pipefail
export AILANG_BIN=$HOME/.pinned-ailang/ailang
LOG=${1:?give an absolute log pathname}
pids=''
cleaned=0
cleanup() {
  if [ "$cleaned" -eq 1 ]; then return; fi
  local pid alive=0
  for pid in $pids; do kill "$pid" 2>/dev/null || :; done
  local until_s=$(( $(date +%s)+10 )) running
  while :; do
    running=0
    for pid in $pids; do
      if kill -0 "$pid" 2>/dev/null; then running=$((running+1)); fi
    done
    if [ "$running" -eq 0 ]; then break; fi
    if [ "$(date +%s)" -ge "$until_s" ]; then
      for pid in $pids; do kill -KILL "$pid" 2>/dev/null || :; done
      break
    fi
    sleep 0.1
  done
  for pid in $pids; do wait "$pid" 2>/dev/null || :; done
  for pid in $pids; do
    if kill -0 "$pid" 2>/dev/null; then alive=$((alive+1)); fi
  done
  echo "yes-alive-after=$alive"
  uptime
  cleaned=1
  [ "$alive" -eq 0 ]
}
trap cleanup EXIT
trap 'exit 130' INT
trap 'exit 143' TERM
uptime
n=0
while [ "$n" -lt 160 ]; do
  yes > /dev/null &
  pids="$pids $!"
  n=$((n+1))
done
alive=0
for pid in $pids; do
  if kill -0 "$pid" 2>/dev/null; then alive=$((alive+1)); fi
done
echo "yes-alive-before=$alive"
[ "$alive" -eq 160 ]
rc=0
perl -e 'alarm 2100; exec @ARGV' go test -timeout 30m -count=100 -race -v -run '^(TestExecProbeReversePollCadence|TestExecProbeReverseHelloWaitBounds)$' ./host/broker/ > "$LOG" 2>&1 || rc=$?
echo "heavy-load rc=$rc"
cleanup
trap - EXIT INT TERM
python3 - "$LOG" "$rc" <<'PY'
from pathlib import Path
from decimal import Decimal
import re,sys
s=Path(sys.argv[1]).read_text(errors='replace')
rows=re.findall(r'CADENCE dials=(\d+) span=\S+ fastest=(\S+) slowest=\S+',s)
units={'h':Decimal(3600),'m':Decimal(60),'s':Decimal(1),'ms':Decimal('.001'),'µs':Decimal('.000001'),'μs':Decimal('.000001'),'us':Decimal('.000001'),'ns':Decimal('.000000001')}
def duration(s):
    parts=re.findall(r'(\d+(?:\.\d+)?)(ns|µs|μs|us|ms|s|m|h)',s)
    assert ''.join(n+u for n,u in parts)==s,s
    return sum((Decimal(n)*units[u] for n,u in parts),Decimal(0))
if re.search(r'operation not permitted|bind:|listen tcp|listen unix',s):
    sys.exit('UNINFORMATIVE UNDER SANDBOX; controller must rerun outside sandbox')
assert int(sys.argv[2])==0,sys.argv[2]
assert len(rows)==100,len(rows)
minimum=min(int(d) for d,f in rows)
fastest=min(duration(f) for d,f in rows)
print('CADENCE samples=%d minimum-dials=%d (gate 6) smallest-fastest=%s ms (execProbePollEvery floor 50 ms)' % (len(rows),minimum,fastest*1000))
assert minimum>=6 and fastest>=Decimal('.050')
for name in ['TestExecProbeReversePollCadence','TestExecProbeReverseHelloWaitBounds']:
    assert len(re.findall(r'^--- PASS: '+name+r' ',s,re.M))==100,name
assert not re.search(r'--- (FAIL|SKIP):',s)
PY
```

## Mutation round

Runnable artifact: `design_docs/verification/world-iter339/plan-mutants.sh.txt`. It makes a complete post-M1 scratch snapshot (including design_docs, needed by TestExecArgsSketchDrift), excludes .git, and checks source post hashes before starting. Every replacement is exact and refuses a site count other than one; each pair checks each site independently. It restores both affected files from the immutable pristine snapshot after every mutated row and asserts cmp, also in an exit/signal trap. Controls mutate nothing and need no restore. It records the raw rc, complete deduplicated FAIL list including parents, top-level counts, uptime, expected result, and a per-row log. Its expected.json contains full names. `SITES_ONLY=1 /bin/bash ...` validates all sites and restores without go test. `ONLY=ID /bin/bash ...` runs a single existing row. Neither option relaxes a gate.

F = `^(TestExecProbeReversePollCadence|TestExecProbeReverseHelloWaitBounds|TestExecProbeReverseLegSilentHolderRefusesSticky|TestExecProbeReversePollRecordsConnectReset|TestExecProbeReverseDialErrIsSilent)$`.
W = `^(TestExecProbeReverse|TestExecProbeForeignListener|TestExecProbeChildListener)`.
T10 = `^TestExecProbeReverseHelloWaitBounds$`.
B = `^TestExecProbeReversePollCadence$`, count=30. All other rows count=1. Every row runs `perl -e 'alarm 2100; exec @ARGV' go test -timeout 30m -count="$count" -v -run "$regex" ./host/broker/` on scratch. No separate unbounded vet: go test's default vet is required for S2-argdrop.

The following expectations come from design §3, not from sandbox failures. T10 “both leaves” means parent plus `/slow-child` and `/slow-mux`; T7 means parent `TestExecProbeReverseLegSilentHolderRefusesSticky` plus indicated leaves. RC 1 is expected for killed mutants; green rows require rc 0 and an empty FAIL list. PC-const2s must red T10 while T12 stays green. M1-5x and M1-45ms are named B-5x/B-45ms in the banked runner. B-const1ms is K-const1ms. Known load-only extra reds are recorded and triaged, never whitelisted.

| Row | Regex / count | Expected rc | Expected exact FAIL set |
|---|---|---|---|
| C0-broad | W / 1 | 0 | none (green) |
| M1-broad | W / 1 | 1 | T12 |
| M1 | F / 1 | 1 | T12 |
| M1-10x | F / 1 | 1 | T12 |
| M1-60s | F / 1 | 1 | T12 |
| M1-div10 | F / 1 | 1 | T12 |
| M1-1ms | F / 1 | 1 | T12 |
| M1-busy | F / 1 | 1 | T12 |
| M1-ctl50 | F / 1 | 0 | none (green) |
| M1-ctl50+const100 | F / 1 | 1 | T12 |
| PC-const2s | F / 1 | 1 | T10 + both leaves |
| K-const100 | F / 1 | 0 | none (green) |
| K-const25 | F / 1 | 0 | none (green) |
| K-listen1100 | F / 1 | 0 | none (green) |
| K-wait950 | F / 1 | 0 | none (green) |
| C1-broad | W / 1 | 0 | none (green) |
| S2-msg | F / 1 | 0 | none (green) |
| S2-wait600 | F / 1 | 0 | none (green) |
| S2-msg+wait600 | F / 1 | 1 | T7 + all 3 leaves |
| S2-argdrop | F / 1 | 1 | build/vet printf error; no test FAIL rows |
| S4-ceil | T10 / 1 | 0 | none (green) |
| K-wait960 | T10 / 1 | 1 | T10 + both leaves |
| S4-ceil+wait960 | T10 / 1 | 0 | none (green) |
| S4-ceil2poll | T10 / 1 | 0 | none (green) |
| S4-floor100 | T10 / 1 | 0 | none (green) |
| K-wait249 | T10 / 1 | 1 | T10 + both leaves |
| S4-floor100+wait249 | T10 / 1 | 0 | none (green) |
| S4-floor500 | T10 / 1 | 0 | none (green) |
| S4-floor600 | T10 / 1 | 1 | T10 + both leaves |
| S4-arrival | T10 / 1 | 0 | none (green) |
| K-wait5 | T10 / 1 | 1 | T10 + both leaves |
| S4-arrival+wait5 | T10 / 1 | 0 | none (green) |
| S4-nowrote | T10 / 1 | 0 | none (green) |
| P-nohello | F / 1 | 1 | T10 + both leaves, T7 + silent-after-chatty |
| S4-nowrote+nohello | F / 1 | 1 | T7 + silent-after-chatty |
| S4-nosets | T10 / 1 | 0 | none (green) |
| P-nodeadline | T10 / 1 | 1 | T10 + both leaves |
| S4-nosets+nodeadline | T10 / 1 | 1 | T10 + both leaves |
| S4-nosetsnozero+nodeadline | T10 / 1 | 0 | none (green) |
| TM1-nocount | F / 1 | 0 | none (green) |
| TM1-nocount+M1 | F / 1 | 0 | none (green) |
| TM2-countflip | F / 1 | 1 | T12 |
| TM3-nofloor+div10 | F / 1 | 0 | none (green) |
| TM3-floorflip | F / 1 | 1 | T12 |
| TM4-periods1+M1 | F / 1 | 0 | none (green) |
| TM5-threshold+M1 | F / 1 | 0 | none (green) |
| TM6-reflit | F / 1 | 0 | none (green) |
| TM6-reflit+M1 | F / 1 | 0 | none (green) |
| TM8-cap1+div10 | F / 1 | 0 | none (green) |
| TM7-noref | F / 1 | 1 | T12 |
| TM9-nosetAt | F / 1 | 1 | T10 + both leaves |
| TM10-setAtIsDeadline | F / 1 | 0 | none (green) |
| TM10-setAtIsDeadline+wait960 | F / 1 | 0 | none (green) |
| TM11-stopref+wait951 | T10 / 1 | 1 | T10 + both leaves |
| B-2x | B / 30 | 0 | none (green) |
| B-3x | B / 30 | 0 | none (green) |
| B-4x | B / 30 | 0 | none (green) |
| B-5x | B / 30 | 1 | T12 |
| B-45ms | B / 30 | 1 | T12 |
| B-const1ms | B / 30 | 0 | none (green) |
| K-wait951 | T10 / 1 | 1 | T10 + both leaves |

There are **58 distinct mutant configurations (29 production, 29 test-code/pairs)** and **61 executions**: the 58, C0/C1 controls, and an extra M1 under W. This includes every killed row, the positive control, the 2×–4× accepted band (30 runs each), survivor 2 alone and its observable pairs, survivor 4's ceiling/floor/arrival/wrote/sets clauses and all weakening pairs. Design §4 step 3's eight named rows are all included, with added explicit K-wait951 (measured in design V07 rather than the designer A/B/C script).

The design's C0 broad control was rc 1 under load 92, while C1 was rc 0; it records unchanged real-socket flakes, not an acceptable control outcome. Likewise M1-broad had an extra exhaustion red; M1-1ms/TM3-nofloor/TM8 had T7 load reds. This runner expects the causal sets above; preserve any extras and the control outcome for controller disposition. The T12 count gate has measured margin, not a proof against asymmetric starvation, and load can weaken the slow-mutant kill. Do not change timing constants, the reference window, or expectations to make a red round pass.

## Pristine-base application proof

Planner used a full `git archive dcb2366e` copy, applied the banked prototype patch, generated the one-line CI change, then built ONE unified diff. On a second pristine archive copy `git apply --check` returned **0**, applying returned **0**, the patched test hash equaled **bbfccede19236d17b81d40401cc079aeb9987ff356c6a08affe6bc2a159e80d0**, and production bytes compared equal to base. The read-only git diff against dcb2366e listed **exactly .github/workflows/ci.yml and host/broker/exec_probe_test.go**. CI post hash **a11ff4ab05a41c9864299edb8f636675e3abb5dcdcba974d43fb6a405940710f**; diff hash **b75e4cf0517633d71e70c88468a93748e58c95cfedfe9ef257d244b4dde146de**. No git metadata was created or changed. Reproduce with:

```bash
B=$(mktemp -d /tmp/world-iter339-apply-proof.XXXXXX)
GD=$(git rev-parse --absolute-git-dir)
git archive dcb2366e | tar -xf - -C "$B"
(cd "$B" && git apply --check "$WT/$D" && git apply "$WT/$D")
# Run the post-hash Python block above from B too.
git --git-dir="$GD" --work-tree="$B" diff --name-only dcb2366e -- . ':!design_docs'
git --git-dir="$GD" --work-tree="$B" diff --quiet dcb2366e -- host/broker/exec_probe.go
```


## Planner measurements and findings

All commands ran on a full pristine archive of dcb2366e with plan-M1.diff.txt applied, at `/var/folders/kr/h0wr2sj94vd6ljtmsxv8jkt00000gn/T/world-iter339-plan.hoibo5j4/base`; no code in the planner worktree was edited. go1.26.6 darwin/arm64, pinned AILANG v0.41.0. Exported AILANG_BIN for all broker and full-suite tests. After all tests, post hashes and the full scope diff were rechecked unchanged. Scratch logs and immutable copies remain at `/var/folders/kr/h0wr2sj94vd6ljtmsxv8jkt00000gn/T/world-iter339-plan.hoibo5j4` for local audit; results and commands are recorded here and in the JSON, so the controller does not depend on /tmp surviving.

| Gate | Raw rc / observations | Verdict |
|---|---|---|
| focus | rc 0; 4.33s; 2 top-level PASS, 2 T10 leaves PASS, no FAIL/SKIP; CADENCE dials=21 span=1.037708167s fastest=50.072292ms slowest=55.103041ms | measured |
| package | rc 1; 140.60s; partial raw census 240 PASS / 25 SKIP / 13 FAIL top-level, aborted by httptest bind panic; NOT the required 336/30 census | UNINFORMATIVE UNDER SANDBOX |
| gofmt | rc 0; 0.05s; empty output | measured |
| vet | rc 0; 0.51s; empty output | measured |
| build | rc 0; 1.74s; empty output | measured |
| suite | rc 1; 392.08s; socket-denied package panics, denied controlling TTY and archive git-ls-files issue; see below | UNINFORMATIVE UNDER SANDBOX |
| ail | rc 0; 63.57s; final verify gate PASSED | measured |
| HEAVY-LOAD | Not run in planner sandbox | Controller-run outside sandbox; acceptance pending |

The full-suite go test was attempted independently after build rc 0, even though the package check had socket-denied results, to collect its diagnostics. No failing command was edited or replaced. The full-suite rc and package rc above are instrument records, not acceptance verdicts. Full-suite package failures include `github.com/sunholo-data/ailang-world/cmd/ailang-worldd`, `github.com/sunholo-data/ailang-world/cmd/world-publish`, `github.com/sunholo-data/ailang-world/host/broker`, `github.com/sunholo-data/ailang-world/host/capsule`, `github.com/sunholo-data/ailang-world/host/daemon`, `github.com/sunholo-data/ailang-world/host/pinfetch`, `github.com/sunholo-data/ailang-world/host/pins`, `github.com/sunholo-data/ailang-world/host/ptytest`, `github.com/sunholo-data/ailang-world/scripts/marketingcheck`.

**Base findings:** the exact targeted command below returned rc 1 on pristine base AND post-M1, with the exact FAIL list `TestOrdinaryOverflowCarriesNoESRCH`, `TestExecHandlerTimeoutKillsTheGroup`. No sockets are used by this targeted run. The former said the group leader never became an unreaped zombie; the latter observed exit_code=0 / timed_out=false. `ps` itself is denied by this sandbox (measured `operation not permitted: ps`), so process-observation restrictions are a possible cause, not a proven diagnosis. These are measured base reds requiring controller triage outside sandbox, not something to work around in M1.

```bash
export AILANG_BIN=$HOME/.pinned-ailang/ailang
go test -timeout 30m -count=1 -v -run '^(TestOrdinaryOverflowCarriesNoESRCH|TestExecHandlerTimeoutKillsTheGroup)$' ./host/broker/
```

**Archive artifact:** full-suite `host/pins/TestReleaseLiteralsExplainedByManifest` reported git ls-files rc 128 because git archive contains no .git. Pristine base archive reproduced rc 1. With read-only `GIT_DIR` pointing to the existing planner metadata and `GIT_WORK_TREE` pointing to the appropriate scratch tree, `go test -timeout 30m -count=1 -v ./host/pins/` returned rc 0 on both base and post-M1. This is a targeted diagnostic using existing metadata; no repository/index was created or modified. The original full-suite result is retained unchanged.

**Mutation measurements:** all 61 site-check executions rc 0, exactly 59 mutation restores/cmp verified; /bin/bash -n rc 0. The actual round recorded 26 informative expectation matches, 35 socket-denied executions, 0 informative mismatches. Complete rc/FAIL lists are in JSON `milestones[0].mutants.planner_measurements`; each raw row log is under `/var/folders/kr/h0wr2sj94vd6ljtmsxv8jkt00000gn/T/world-iter339-plan.hoibo5j4/mutants`. The exact F/W gates cannot be accepted under the sandbox. No claimed green/red causal result is inferred from those rows.

| Mutation execution | Raw rc | Planner classification |
|---|---|---|
| C0-broad | 1 | UNINFORMATIVE UNDER SANDBOX |
| M1-broad | 1 | UNINFORMATIVE UNDER SANDBOX |
| M1 | 1 | UNINFORMATIVE UNDER SANDBOX |
| M1-10x | 1 | UNINFORMATIVE UNDER SANDBOX |
| M1-60s | 1 | UNINFORMATIVE UNDER SANDBOX |
| M1-div10 | 1 | UNINFORMATIVE UNDER SANDBOX |
| M1-1ms | 1 | UNINFORMATIVE UNDER SANDBOX |
| M1-busy | 1 | UNINFORMATIVE UNDER SANDBOX |
| M1-ctl50 | 1 | UNINFORMATIVE UNDER SANDBOX |
| M1-ctl50+const100 | 1 | UNINFORMATIVE UNDER SANDBOX |
| PC-const2s | 1 | UNINFORMATIVE UNDER SANDBOX |
| K-const100 | 1 | UNINFORMATIVE UNDER SANDBOX |
| K-const25 | 1 | UNINFORMATIVE UNDER SANDBOX |
| K-listen1100 | 1 | UNINFORMATIVE UNDER SANDBOX |
| K-wait950 | 1 | UNINFORMATIVE UNDER SANDBOX |
| C1-broad | 1 | UNINFORMATIVE UNDER SANDBOX |
| S2-msg | 1 | UNINFORMATIVE UNDER SANDBOX |
| S2-wait600 | 1 | UNINFORMATIVE UNDER SANDBOX |
| S2-msg+wait600 | 1 | UNINFORMATIVE UNDER SANDBOX |
| S2-argdrop | 1 | matched |
| S4-ceil | 0 | matched |
| K-wait960 | 1 | matched |
| S4-ceil+wait960 | 0 | matched |
| S4-ceil2poll | 0 | matched |
| S4-floor100 | 0 | matched |
| K-wait249 | 1 | matched |
| S4-floor100+wait249 | 0 | matched |
| S4-floor500 | 0 | matched |
| S4-floor600 | 1 | matched |
| S4-arrival | 0 | matched |
| K-wait5 | 1 | matched |
| S4-arrival+wait5 | 0 | matched |
| S4-nowrote | 0 | matched |
| P-nohello | 1 | UNINFORMATIVE UNDER SANDBOX |
| S4-nowrote+nohello | 1 | UNINFORMATIVE UNDER SANDBOX |
| S4-nosets | 0 | matched |
| P-nodeadline | 1 | matched |
| S4-nosets+nodeadline | 1 | matched |
| S4-nosetsnozero+nodeadline | 0 | matched |
| TM1-nocount | 1 | UNINFORMATIVE UNDER SANDBOX |
| TM1-nocount+M1 | 1 | UNINFORMATIVE UNDER SANDBOX |
| TM2-countflip | 1 | UNINFORMATIVE UNDER SANDBOX |
| TM3-nofloor+div10 | 1 | UNINFORMATIVE UNDER SANDBOX |
| TM3-floorflip | 1 | UNINFORMATIVE UNDER SANDBOX |
| TM4-periods1+M1 | 1 | UNINFORMATIVE UNDER SANDBOX |
| TM5-threshold+M1 | 1 | UNINFORMATIVE UNDER SANDBOX |
| TM6-reflit | 1 | UNINFORMATIVE UNDER SANDBOX |
| TM6-reflit+M1 | 1 | UNINFORMATIVE UNDER SANDBOX |
| TM8-cap1+div10 | 1 | UNINFORMATIVE UNDER SANDBOX |
| TM7-noref | 1 | UNINFORMATIVE UNDER SANDBOX |
| TM9-nosetAt | 1 | UNINFORMATIVE UNDER SANDBOX |
| TM10-setAtIsDeadline | 1 | UNINFORMATIVE UNDER SANDBOX |
| TM10-setAtIsDeadline+wait960 | 1 | UNINFORMATIVE UNDER SANDBOX |
| TM11-stopref+wait951 | 1 | matched |
| B-2x | 0 | matched |
| B-3x | 0 | matched |
| B-4x | 0 | matched |
| B-5x | 1 | matched |
| B-45ms | 1 | matched |
| B-const1ms | 0 | matched |
| K-wait951 | 1 | matched |

**Supplemental socket-free spot checks:** to measure design §4 step 3’s eight behaviors despite the F regex’s socket requirements, each was also applied with single-occurrence checks and restored from a pristine post-M1 copy, then run using the exact command below. This supplements F, not a substitute for the outside-sandbox F acceptance gate. All eight results match the design’s causal expectations.

```bash
export AILANG_BIN=$HOME/.pinned-ailang/ailang
go test -timeout 30m -count=1 -v -run '^(TestExecProbeReversePollCadence|TestExecProbeReverseHelloWaitBounds)$' ./host/broker/
```

| Row | rc | FAIL set (parents included) | CADENCE |
|---|---|---|---|
| M1 | 1 | `TestExecProbeReversePollCadence` | CADENCE dials=3 span=1.017072083s fastest=501.028875ms slowest=501.031958ms |
| M1-div10 | 1 | `TestExecProbeReversePollCadence` | CADENCE dials=178 span=1.012444291s fastest=5.029208ms slowest=9.166917ms |
| M1-ctl50+const100 | 1 | `TestExecProbeReversePollCadence` | CADENCE dials=40 span=2.036669916s fastest=50.020417ms slowest=68.301625ms |
| PC-const2s | 1 | `TestExecProbeReverseHelloWaitBounds`, `TestExecProbeReverseHelloWaitBounds/slow-child`, `TestExecProbeReverseHelloWaitBounds/slow-mux` | CADENCE dials=21 span=40.040428833s fastest=2.000080166s slowest=2.006423667s |
| K-wait951 | 1 | `TestExecProbeReverseHelloWaitBounds`, `TestExecProbeReverseHelloWaitBounds/slow-child`, `TestExecProbeReverseHelloWaitBounds/slow-mux` | CADENCE dials=20 span=1.02971325s fastest=50.022333ms slowest=56.273125ms |
| K-wait950 | 0 | none | CADENCE dials=21 span=1.035919708s fastest=50.052917ms slowest=56.328584ms |
| TM2-countflip | 1 | `TestExecProbeReversePollCadence` | CADENCE dials=21 span=1.019665958s fastest=50.259041ms slowest=53.887ms |
| TM9-nosetAt | 1 | `TestExecProbeReverseHelloWaitBounds`, `TestExecProbeReverseHelloWaitBounds/slow-child`, `TestExecProbeReverseHelloWaitBounds/slow-mux` | CADENCE dials=20 span=1.029208625s fastest=50.055625ms slowest=56.395042ms |

**Design corrections and limits:** the patch/file hash attribution is corrected above. No behavioral design expectation was contradicted by an informative planner mutation result. The 336 PASS / 30 SKIP package expectation remains unverified here; the package aborted before completion and two unrelated base process tests also red. T12’s heavy-load acceptance and all socket-dependent outcomes remain controller responsibilities. The historical 150–214 ambient workload is not identical to 160 added spinners; record the real uptime readings rather than claim equivalence.

All load-bearing commands must be rerun by the controller outside the sandbox before judging M1. A denied loopback bind or listen panic is **UNINFORMATIVE UNDER SANDBOX**, never a passing/failing acceptance verdict. In particular no broker package census can be accepted from an aborted run. Findings already present on base remain findings: preserve their logs, report them, and ask the controller to disposition them; no unrelated code fix or acceptance weakening is in scope.
