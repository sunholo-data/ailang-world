# Sprint plan — `w-finish-scope-expiry-boundary` (iteration 327, row 236, red on `dev`)

**Design doc:** [`design_docs/planned/w-finish-scope-expiry-boundary.md`](w-finish-scope-expiry-boundary.md)
(r3: r2 after quorum r1, plus the quorum-r2 fixes applied by the controller; §10). **Sprint base:**
worktree `/Users/voightkampff/dev/sunholo-data/ailang-world/.wt-world-iter327`, branch
`sprint/world-iter327-finish-scope-expiry`, HEAD `f7767f12`. Its Go tree equals the design's
measurement base `4d43d8e3`: `git diff --name-only 4d43d8e3 HEAD` lists 21 files, all under
`design_docs/` (V-P1). **Planner:** mission-control iteration 327 (planner role). **Executor:**
claude-haiku, so the whole code change is one machine-applied diff (§M1.3) and every command is
spelled out. **Copy code and commands from fenced blocks, never from table cells.**
**Estimate:** about 0.1 day of change; 1.5–2 h wall for the executor, dominated by package runs.
**Two milestones**, each independently green. **Risk:** low (one unexported field, three reads
moved onto it, tests). The planner prototyped the exact planned code, every acceptance command and
every §4 mutant in throwaway copies outside the repo (`/tmp/w327p/base`, `/tmp/w327p/fix`, both
`git archive HEAD`); results in §9 and §10.

**This prototype run is the first verification of the r2 diff.** Design §3/§4/AC-3 say the
P/Q/R/TM1 kill attributions and the package-green result were measured on the pre-rename (r1)
prototype only. The planner's run on the exact planned files (`finishNow`, r2) re-measured all of
them: AC-1/AC-1b, all 16 bound mutants, S1–S3, W1, TM1 and the whole `host/daemon` package plain
and `-race` (§10 V-P5…V-P9). The executor's run is the second.

**Scope (the only files you may edit or create):**
- `host/daemon/daemon.go` (M1; the field + `New` wiring)
- `host/daemon/episode_finish.go` (M1; three reads, the `time` import removed)
- `host/daemon/episode_finish_test.go` (M1; T1 rewritten, T2 and T3 new, AC11 `future` arm)
- `design_docs/verification/world-iter327/executor-report.md` (M2; new)
- `design_docs/verification/world-iter327/executor-*.txt` (M1/M2; new evidence logs)

**No `.ail` file. No CI file. No `docs/` file. No `tools/launchd/*`. No other `.go` file.**
**Do NOT edit `design_docs/world-mission*.md`.** Row 236 does not exist in the charter at base
(`grep -cE '^236\. ' design_docs/world-mission.md` → `0`; control `grep -nE '^235\. '` → `6304:235.`,
V-P10). The **controller** writes row 236. Do not move the design doc to `implemented/`; the
controller does that too. Do not touch the `designer-*`, `controller-*`, `gate1-2-readings.txt` or
`quorum/` files already in `design_docs/verification/world-iter327/`.

---

## 0. The one-paragraph version

`TestEpisodeFinishScopeExpiry` builds its `+3601` arm from its own `time.Now()` and the handler
re-reads the clock later (`episode_finish.go:168-169` at base), so a whole-second edge between the
two reads turns the expected 403 into a 200 (CI run 41093918891 attempt 1, race gate). AC11's
`future` arm (`+301` against the proof's `> now+300`) has the same shape. The fix (design §2,
option iii) adds an unexported `finishNow func() int64` to `Daemon`, wired by `New` to the wall
clock, and moves the finish route's three clock reads onto it. The tests pin it at
`finishClockAt = 2_000_000_000` and test every exact edge:
- **T1** `TestEpisodeFinishScopeExpiry` (rewritten): scope at now −1/0/+1/+3600/+3601 → 403/403/200/200/403.
- **T2** `TestEpisodeFinishOperatorProofExpiryBoundary` (new): proof at now −1/0/+1/+300/+301 → 403/403/200/200/403.
- **T3** `TestEpisodeFinishRecordRechecksScopeExpiry` (new): scope approved at T+3600; record clock T−1/T/T+3599/T+3600 → 403/200/200/403.
- **AC11 `future`**: `+301` → `+3600` (wall clock, refused with margin in the safe direction).

Milestones:
- **M1:** AC-1/AC-1b red at base (harness), apply the diff, compile/vet, AC-1/AC-1b green, AC-2 mutation drill.
- **M2:** AC-3 (package plain + race, `-count=200` race, `cmd/ailang-worldd -run Finish`), AC-4 (full gate + `verify_ail.sh`), executor report.

---

## 1. Environment (on EVERY command)

```sh
cd /Users/voightkampff/dev/sunholo-data/ailang-world/.wt-world-iter327
export PATH=/opt/homebrew/bin:$PATH
export AILANG_BIN=$HOME/.pinned-ailang/ailang     # reports "AILANG v0.41.0" (measured, V-P2)
mkdir -p /tmp/w327x/tools                         # the executor's scratch dir; NEVER inside the worktree
PLAN=design_docs/planned/w-finish-scope-expiry-boundary-sprint-plan.md
x() { awk -v h="$1" '$0==h{f=1;next} f&&/^```$/{f=0} f' "$PLAN"; }   # extracts one fenced block by its exact opening line
```

The shell is **zsh**. Run every multi-line script and loop with `bash` (the tool scripts below have
a `#!/bin/bash` line and are run as executables).

| Pin | Measured this session | Why it matters |
|---|---|---|
| `go` | `go version go1.26.6 darwin/arm64` | The CI pin |
| `AILANG_BIN` | `AILANG v0.41.0` (the CI pin) | Without it daemon e2e tests SKIP and `verify_ail.sh` cannot run |
| `patch` | `/usr/bin/patch` | Applies the §M1.3 diff |

**Instrument rule.** Never write `cmd | tail; echo rc=$?`. Redirect to a file, put `echo "rc=$?"` on
the same line after `;`, then read the file. A `-run` that matches nothing exits 0, so read every
verdict from counted `--- PASS:` / `--- FAIL:` lines.

**Compile fence.** `go build ./...` does **not** compile `_test.go` files. The compile fence is
`go vet ./host/daemon/` (or `go test -run '^$' ./host/daemon/`).

**Executor constraints (binding).**
1. **No git write operations of any kind:** no `commit`, `add`, `stash`, `checkout -- file`,
   `restore`, `reset`, `apply`. The controller commits. `git diff`, `git status`, `git show`,
   `git archive` are allowed.
2. Edit only the *Scope* files, only through the §M1.3 diff. Tool scripts, snapshots and logs go
   under `/tmp/w327x/`, **never** inside the worktree. The only file a tool script creates inside
   the worktree is the S-mutant scratch `host/daemon/zz_mut_wall.go`, which `allmuts.sh` deletes and
   reports (`scratch removed: yes`).
3. Never edit `tools/launchd/*`. Never touch the main checkout
   `/Users/voightkampff/dev/sunholo-data/ailang-world` or `~/.ailang/state/mission-v1*`. Never edit
   `design_docs/world-mission*.md`. Do not run `ailang messages`.
4. Every mutant is applied **alone** by `runmut.sh`, which restores the file and prints
   `sha=<before>/<after>`; the two must be equal. A mutant with `build_err` ≠ 0 is a **protocol
   error, never a kill**. Do not `gofmt` a mutant.
5. Do not add `t.Parallel()`. The tests write `d.finishNow` between synchronous handler calls; a
   parallel subtest would race it (design §2, concurrency).
6. **Bound every wait**: wrap long runs in `perl -e 'alarm N; exec @ARGV' …` and give `go test` an
   explicit `-timeout`. Run long tests in the **foreground**. Never end your turn announcing a wait.
7. Never run two package-scale `go test` commands at once (load widens timing windows in unrelated
   tests; see §M2.3).

---

## 2. Gate list

| ID | Command | Pass |
|---|---|---|
| G1 | `gofmt -l host/daemon/` | empty output |
| G2 | `go vet ./host/daemon/` | rc 0 (**the compile fence**) |
| G3 | §M1.4 finish tests | rc 0, 5 `^--- PASS: TestEpisodeFinish`, 0 FAIL, plain and `-race` |
| G4 | AC-1 / AC-1b on the fix (§M1.5) | each: rc 0, PASS=20, FAIL=0, plain and race |
| G5 | AC-2 (§M1.6) | 21 mutant lines each `rc=1 build_err=0` with the named killer in `fails=`; the S0 control `rc=0`; every `sha=X/X` equal; `scratch removed: yes` |
| G6 | AC-3 (§M2.1) | vet rc 0; package plain and race `ok  	github.com/sunholo-data/ailang-world/host/daemon`, 0 `DATA RACE`; `-count=200` race `ok`; `cmd/ailang-worldd -run Finish` `ok` |
| G7 | AC-4 (§M2.2) | `go build ./... && go test ./...` rc 0 (or every failure classified per §M2.3 as load/pre-existing, none in `host/daemon`); `verify_ail.sh` rc 0 |
| G8 | `git diff --name-only`; `git status --porcelain` | diff: exactly the 3 `host/daemon` files; status adds only `executor-*` under `design_docs/verification/world-iter327/`; 0 `world-mission` lines |

---

## 3. Pristine-tree baselines (record these BEFORE the first edit)

```sh
shasum -a 256 host/daemon/daemon.go host/daemon/episode_finish.go host/daemon/episode_finish_test.go > /tmp/w327x/base.sha
cat /tmp/w327x/base.sha
git status --porcelain
git log --oneline -1
```

Expected (planner, this session):

| File | base sha256 | planned (fixed) sha256 |
|---|---|---|
| `host/daemon/daemon.go` | `5b09012f8d19df9e2d3d15588aee9349cf3112782222f07bac7e71dafcca7104` | `e6b70d7e3f3fb20ef6a89d52b82a87bbdff0a05915c6afaab8956cb05be6c668` |
| `host/daemon/episode_finish.go` | `ea44759896b57b12aeff835469c978e91bfe6236bd22020305658ee51daa3703` | `e429375dc3e0b03841d5bef84eafcdc4b718dd324c287c8df8bee30f42e9bf0b` |
| `host/daemon/episode_finish_test.go` | `094884e036f77e325ea4984109d7d384e3ea8262f87350a80da2913834c70992` | `b872a38cc091e250e570d1bba16afa17044c57ff1d9a69cc63484400429b9d67` |

`git status --porcelain` at base shows only this plan and its JSON (the controller commits them).
If any base sha256 differs, **stop** and report: the tree is not the planned base.

### 3.1 Base RED / GREEN

- **Red at base (deterministic, harness):** AC-1 (`TestEpisodeFinishScopeExpiry`, spin before `:164`) and
  AC-1b (`TestEpisodeFinishOperatorProofAC11`, spin before `:79`), 20/20 each, plain and race.
- **Red at base (flaky):** the same two tests whenever a second edge falls between the test's and
  the handler's clock reads (CI run 41093918891 attempt 1).
- **Absent at base:** T2 `TestEpisodeFinishOperatorProofExpiryBoundary`, T3 `TestEpisodeFinishRecordRechecksScopeExpiry`.
- **Green at base and must stay green:** `TestEpisodeFinishOperatorProofAC11`,
  `TestEpisodeFinishAuthorizationPrecedesBody`, the whole `host/daemon` package, `cmd/ailang-worldd -run Finish`.

---

## 4. Milestones

### Milestone M1 — red at base, the diff, AC-1/AC-1b green, AC-2 mutation drill (+84/−13, three files)

#### M1.1 Write the tool scripts (outside the worktree)

```sh
x '```python mut.py' > /tmp/w327x/tools/mut.py
x '```bash runmut.sh' > /tmp/w327x/tools/runmut.sh
x '```bash harness.sh' > /tmp/w327x/tools/harness.sh
x '```bash allmuts.sh' > /tmp/w327x/tools/allmuts.sh
chmod +x /tmp/w327x/tools/*.sh
shasum -a 256 /tmp/w327x/tools/* > /tmp/w327x/tools.sha; cat /tmp/w327x/tools.sha
```

Expected tool sha256 (planner's round trip, V-P11):

```text
96e27ec707b7941117eea6c345c516ed083b893ebb07295441f3a8936a520be7  /tmp/w327x/tools/allmuts.sh
c2531dca1783fda4705fa1219080c6a513a00715068029425704d9e771183b67  /tmp/w327x/tools/harness.sh
715ab30279a747ef78d2e184386575348b6ee0037ee4906988da895231dd3b35  /tmp/w327x/tools/mut.py
bd8aed15d5ff295de5eef0c290ea19dba229123aedd2d7cf432a9afae7e4eee7  /tmp/w327x/tools/runmut.sh
```

```python mut.py
# usage: mut.py FILE LINE OLD NEW  -- replace OLD with NEW exactly once on LINE (1-indexed)
import sys
path, line, old, new = sys.argv[1], int(sys.argv[2]), sys.argv[3], sys.argv[4]
lines = open(path).read().split('\n')
src = lines[line-1]
assert src.count(old) == 1, (line, src, old)
lines[line-1] = src.replace(old, new)
open(path, 'w').write('\n'.join(lines))
```

```bash runmut.sh
#!/bin/bash
# usage: runmut.sh TREE RUNREGEX COUNT ID FILE LINE OLD NEW [extra go test flags]
tree=$1; re=$2; n=$3; id=$4; f=$5; ln=$6; old=$7; new=$8; shift 8
cd "$tree" || exit 9
cp "$f" "/tmp/w327x/tools/orig.$$"
before=$(shasum -a 256 "$f" | cut -c1-12)
python3 -I /tmp/w327x/tools/mut.py "$f" "$ln" "$old" "$new" || { echo "$id MUT-APPLY-FAILED"; cp "/tmp/w327x/tools/orig.$$" "$f"; exit 9; }
out=$(go test -timeout 180s "$@" ./host/daemon/ -run "$re" -count="$n" -v 2>&1); rc=$?
cp "/tmp/w327x/tools/orig.$$" "$f"; rm "/tmp/w327x/tools/orig.$$"
after=$(shasum -a 256 "$f" | cut -c1-12)
build=$(echo "$out" | grep -cE 'build failed|imported and not used|declared and not used|\[setup failed\]')
fails=$(echo "$out" | grep -E -- '--- FAIL' | sed -E 's/ *--- FAIL: ([^ ]+).*/\1/' | sort | uniq -c | tr -s ' ' | tr '\n' ';')
first=$(echo "$out" | grep -E '_test.go:[0-9]+:' | head -1 | sed 's/^ *//')
echo "$id rc=$rc build_err=$build fails=[$fails] first=[$first] sha=$before/$after"
```

```bash harness.sh
#!/bin/bash
# usage: harness.sh TREE FILE LINE COUNT RUNREGEX [extra go test flags]  -- inserts the boundary spin before LINE, runs, restores
tree=$1; f=$2; ln=$3; n=$4; re=$5; shift 5
cd "$tree" || exit 9
cp "$f" /tmp/w327x/tools/h.orig
before=$(shasum -a 256 "$f" | cut -c1-12)
python3 -I - "$f" "$ln" <<'PY'
import sys, os
p, ln = sys.argv[1], int(sys.argv[2])
L = open(p).read().split('\n')
L.insert(ln-1, '\t' + os.environ.get('SPIN', 'for time.Now().Nanosecond() < 999_950_000 {}') + ' // harness')
open(p,'w').write('\n'.join(L))
PY
out=$(go test -timeout 300s "$@" ./host/daemon/ -run "$re" -count="$n" -v 2>&1); rc=$?
cp /tmp/w327x/tools/h.orig "$f"
after=$(shasum -a 256 "$f" | cut -c1-12)
echo "rc=$rc PASS=$(echo "$out" | grep -cE '^--- PASS') FAIL=$(echo "$out" | grep -cE '^--- FAIL') first=[$(echo "$out" | grep -E '_test.go:[0-9]+:' | head -1 | sed 's/^ *//')] tail=[$(echo "$out" | tail -1)] sha=$before/$after"
```

```bash allmuts.sh
#!/bin/bash
# all §4 mutants (except TM2) against TREE; each alone, restored, sha checked by runmut.sh
T=$1; shift
M=/tmp/w327x/tools/runmut.sh; RE='^TestEpisodeFinish'; P=host/daemon/episode_finish.go; D=host/daemon/daemon.go; TF=host/daemon/episode_finish_test.go
$M $T $RE 1 P1 $P 168 'scope.ExpiresAt <= now' 'scope.ExpiresAt < now' "$@"
$M $T $RE 1 P2 $P 168 'scope.ExpiresAt <= now ||' 'false ||' "$@"
$M $T $RE 1 P3 $P 168 'scope.ExpiresAt > now+3600' 'scope.ExpiresAt >= now+3600' "$@"
$M $T $RE 1 P4 $P 168 'now+3600' 'now+3601' "$@"
$M $T $RE 1 P5 $P 168 'now+3600' 'now+3660' "$@"
$M $T $RE 1 P6 $P 168 'scope.ExpiresAt > now+3600' 'false' "$@"
$M $T $RE 1 Q1 $P 270 'scope.ExpiresAt <= now' 'scope.ExpiresAt < now' "$@"
$M $T $RE 1 Q2 $P 270 'scope.ExpiresAt <= now ||' 'false ||' "$@"
$M $T $RE 1 Q3 $P 270 'scope.ExpiresAt > now+3600' 'scope.ExpiresAt >= now+3600' "$@"
$M $T $RE 1 Q4 $P 270 'now+3600' 'now+3601' "$@"
$M $T $RE 1 Q6 $P 270 'scope.ExpiresAt > now+3600' 'false' "$@"
$M $T $RE 1 R1 $P 116 'p.ExpiresAt <= now' 'p.ExpiresAt < now' "$@"
$M $T $RE 1 R2 $P 116 'p.ExpiresAt <= now ||' 'false ||' "$@"
$M $T $RE 1 R3 $P 116 'p.ExpiresAt > now+300' 'p.ExpiresAt >= now+300' "$@"
$M $T $RE 1 R4 $P 116 'now+300' 'now+301' "$@"
$M $T $RE 1 R6 $P 116 'p.ExpiresAt > now+300' 'false' "$@"
S=$T/host/daemon/zz_mut_wall.go
printf 'package daemon\n\nimport "time"\n\nfunc wallNow() int64 { return time.Now().Unix() }\n' > $S
$M $T $RE 1 S0-control $P 115 'd.finishNow()' 'd.finishNow()' "$@"
$M $T $RE 1 S1 $P 115 'd.finishNow()' 'wallNow()' "$@"
$M $T $RE 1 S2 $P 167 'd.finishNow()' 'wallNow()' "$@"
$M $T $RE 1 S3 $P 269 'd.finishNow()' 'wallNow()' "$@"
rm -f $S; echo "scratch removed: $(test -e $S && echo NO || echo yes)"
$M $T $RE 1 W1 $D 800 'finishNow: func() int64 { return time.Now().Unix() }' 'finishNow: func() int64 { return 0 }' "$@"
$M $T $RE 1 TM1 $TF 179 'pinFinishClock(d, finishClockAt)' '_ = finishClockAt' "$@"
```

`harness.sh` inserts `for time.Now().Nanosecond() < 999_950_000 {}` (one tab, then a `// harness`
comment) as a new line before LINE, runs the test `-count=N`, restores the file and prints
`sha=<before>/<after>`. `SPIN='…'` overrides the spin (the control).

#### M1.2 AC-1 and AC-1b red at base (BEFORE the diff)

```sh
W=/Users/voightkampff/dev/sunholo-data/ailang-world/.wt-world-iter327; H=/tmp/w327x/tools/harness.sh; F=host/daemon/episode_finish_test.go
{ echo "AC1 base plain:    $($H $W $F 164 20 '^TestEpisodeFinishScopeExpiry$')"
  echo "AC1 base race:     $($H $W $F 164 20 '^TestEpisodeFinishScopeExpiry$' -race)"
  echo "AC1 base control:  $(SPIN='for time.Now().Nanosecond() > 100_000_000 {}' $H $W $F 164 20 '^TestEpisodeFinishScopeExpiry$')"
  echo "AC1b base plain:   $($H $W $F 79 20 '^TestEpisodeFinishOperatorProofAC11$')"
  echo "AC1b base race:    $($H $W $F 79 20 '^TestEpisodeFinishOperatorProofAC11$' -race)"
  echo "AC1b base control: $(SPIN='for time.Now().Nanosecond() > 100_000_000 {}' $H $W $F 79 20 '^TestEpisodeFinishOperatorProofAC11$')"
} > /tmp/w327x/ac1-base.txt 2>&1; echo "rc=$?"; cat /tmp/w327x/ac1-base.txt
shasum -a 256 -c /tmp/w327x/base.sha
```

Want: the four non-control lines `rc=1 PASS=0 FAIL=20`, AC-1 first line
`episode_finish_test.go:171: expiry <n>: 200`, AC-1b first line `episode_finish_test.go:94: 200 {"request":…`;
the two controls `rc=0 PASS=20 FAIL=0`; every `sha=094884e036f7/094884e036f7`; `shasum -c` all `OK`.
Planner (V-P4): exactly that. If a non-control line PASSes, **stop**: the hazard was not built.

#### M1.3 The fix: apply the planned diff (all three files at once)

```sh
x '```diff' > /tmp/w327x/plan.diff
shasum -a 256 /tmp/w327x/plan.diff                 # want d36610127fb2ebb075fe828a4f7bb0bc8802a81b71af22ed297a7db0c9f552eb
patch -p1 --dry-run < /tmp/w327x/plan.diff; echo "dry rc=$?"   # want rc 0, 3 "patching file" lines, no "Hunk … offset/fuzz"
patch -p1 < /tmp/w327x/plan.diff; echo "rc=$?"
shasum -a 256 host/daemon/daemon.go host/daemon/episode_finish.go host/daemon/episode_finish_test.go   # want the "planned" column of §3
git status --porcelain                              # no *.orig / *.rej files
```

If the dry run reports any hunk failure, fuzz or offset, **stop** (the base was checked in §3, so
this means the diff was extracted wrongly). If the sha256s differ from §3's planned column, **stop**
and `diff` against the block below. Do not hand-edit.

```diff
--- a/host/daemon/daemon.go
+++ b/host/daemon/daemon.go
@@ -444,6 +444,16 @@
 	// is the only way the expiry branch of the SHIPPED Shutdown path can be
 	// exercised without a ten-second test.
 	drainTimeout time.Duration
+
+	// finishNow is the clock the operator finish route reads (row 236): the
+	// proof expiry, the approve-time scope bound and the record-time
+	// re-check, all in Unix seconds and all under humanMu on the request
+	// path. New always sets it to the wall clock. It is a field so a test can
+	// pin the exact boundary instead of racing two wall-clock reads; a test
+	// may write it only before or between synchronous handler calls. It is
+	// NOT the coordinator's clock (coordinator.Config.Now, wired below): no
+	// finish-route time is ever compared with a coordinator time.
+	finishNow func() int64
 
 	// readDeadline bounds every store read a GET handler performs. New always
 	// sets it from the readDeadline constant; it is a field rather than a
@@ -787,7 +797,7 @@
 		cfg:       cfg, store: s, bootstrap: registry.Bootstrap, reads: s, intents: s, commits: s, commitBudget: commitBudget, credentialBudget: credentialBudget, drainTimeout: shutdownTimeout,
 		aguiBudget: aguiRunBudget, aguiBodyBound: aguiBodyBound, aguiTick: aguiTick,
 		aguiStop: make(chan struct{}), aguiSlots: make(chan struct{}, aguiCap),
-		readDeadline: readDeadline, errLog: resolveErrorLog(cfg.ErrorLog),
+		readDeadline: readDeadline, errLog: resolveErrorLog(cfg.ErrorLog), finishNow: func() int64 { return time.Now().Unix() },
 		scanPageSize: integrityScanPageSize, scanRowBudget: integrityScanRowBudget,
 		scanTimeBudget: integrityScanTimeBudget, workspace: workspace,
 		resolver: authority.New(authority.ReadOnlyAt(credPath)), credentialsPath: credPath,
--- a/host/daemon/episode_finish.go
+++ b/host/daemon/episode_finish.go
@@ -14,7 +14,6 @@
 	"regexp"
 	"strings"
 	"syscall"
-	"time"
 
 	"github.com/sunholo-data/ailang-world/host/broker"
 	"github.com/sunholo-data/ailang-world/host/hashref"
@@ -113,7 +112,7 @@
 		return fmt.Errorf("operator proof cannot be read")
 	}
 	var p operatorProof
-	now := time.Now().Unix()
+	now := d.finishNow()
 	if json.Unmarshal(data, &p) != nil || p.Phase != phase || p.ExpiresAt <= now || p.ExpiresAt > now+300 {
 		return fmt.Errorf("operator proof phase or expiry is invalid")
 	}
@@ -165,7 +164,7 @@
 	defer cancel()
 	if input.Phase == "approve" {
 		scope, err := broker.ParseFinishApprovalScope(input.Scope)
-		now := time.Now().Unix()
+		now := d.finishNow()
 		if err != nil || scope.ExpiresAt <= now || scope.ExpiresAt > now+3600 || input.Record != nil {
 			refuse("operator route: invalid or expired finish scope")
 			return
@@ -267,7 +266,7 @@
 		return bad()
 	}
 	scope, err := broker.ParseFinishApprovalScope(request.Scope)
-	now := time.Now().Unix()
+	now := d.finishNow()
 	if err != nil || scope.ExpiresAt <= now || scope.ExpiresAt > now+3600 {
 		return bad()
 	}
--- a/host/daemon/episode_finish_test.go
+++ b/host/daemon/episode_finish_test.go
@@ -76,7 +76,7 @@
 			case "expired":
 				proof = finishProofFile(t, d, "approve", time.Now().Unix()-1, 0o600)
 			case "future":
-				proof = finishProofFile(t, d, "approve", time.Now().Unix()+301, 0o600)
+				proof = finishProofFile(t, d, "approve", time.Now().Unix()+3600, 0o600)
 			case "readable":
 				proof = finishProofFile(t, d, "approve", time.Now().Unix()+100, 0o640)
 			case "world-readable":
@@ -159,19 +159,81 @@
 	}
 }
 
+// finishClockAt is far from the wall clock (2033-05-18), so a handler that reads
+// the wall clock instead of d.finishNow refuses every accept arm below (row 236).
+const finishClockAt int64 = 2_000_000_000
+
+func pinFinishClock(d *Daemon, at int64) { d.finishNow = func() int64 { return at } }
+
+func finishScopeAt(expires int64) string {
+	scope, _ := broker.ParseFinishApprovalScope(testFinishScope())
+	scope.ExpiresAt = expires
+	return broker.FormatFinishApprovalScope(scope)
+}
+
+// TestEpisodeFinishScopeExpiry pins the approve phase's scope bound in
+// handleEpisodeFinish at both exact edges: expired is ExpiresAt <= now, too far
+// is ExpiresAt > now+3600. No arm depends on two clock reads sharing a second.
 func TestEpisodeFinishScopeExpiry(t *testing.T) {
 	d := newHandlerDaemon(t)
-	for _, expiry := range []int64{time.Now().Unix() - 1, time.Now().Unix() + 3601} {
-		scope, _ := broker.ParseFinishApprovalScope(testFinishScope())
-		scope.ExpiresAt = expiry
-		proof := finishProofFile(t, d, "approve", time.Now().Unix()+100, 0o600)
-		w := finishRoute(t, d, FinishRequest{Phase: "approve", Scope: broker.FormatFinishApprovalScope(scope)}, proof, "")
-		if w.Code != 403 {
-			t.Fatalf("expiry %d: %d", expiry, w.Code)
+	pinFinishClock(d, finishClockAt)
+	for _, arm := range []struct {
+		offset int64
+		want   int
+	}{{-1, 403}, {0, 403}, {1, 200}, {3600, 200}, {3601, 403}} {
+		proof := finishProofFile(t, d, "approve", finishClockAt+100, 0o600)
+		w := finishRoute(t, d, FinishRequest{Phase: "approve", Scope: finishScopeAt(finishClockAt + arm.offset)}, proof, "")
+		if w.Code != arm.want {
+			t.Fatalf("scope expiry now%+d: %d, want %d: %s", arm.offset, w.Code, arm.want, w.Body.String())
 		}
 	}
 }
 
+// TestEpisodeFinishOperatorProofExpiryBoundary pins consumeOperatorProof's
+// bound: ExpiresAt <= now and ExpiresAt > now+300 are refused.
+func TestEpisodeFinishOperatorProofExpiryBoundary(t *testing.T) {
+	d := newHandlerDaemon(t)
+	pinFinishClock(d, finishClockAt)
+	input := FinishRequest{Phase: "approve", Scope: finishScopeAt(finishClockAt + 900)}
+	for _, arm := range []struct {
+		offset int64
+		want   int
+	}{{-1, 403}, {0, 403}, {1, 200}, {300, 200}, {301, 403}} {
+		proof := finishProofFile(t, d, "approve", finishClockAt+arm.offset, 0o600)
+		if w := finishRoute(t, d, input, proof, ""); w.Code != arm.want {
+			t.Fatalf("proof expiry now%+d: %d, want %d: %s", arm.offset, w.Code, arm.want, w.Body.String())
+		}
+	}
+}
+
+// TestEpisodeFinishRecordRechecksScopeExpiry pins finishDecisionScope's
+// record-time re-check: an approval whose scope has expired by record time, or
+// that names an expiry more than an hour past the record clock, no longer
+// authorizes a record.
+func TestEpisodeFinishRecordRechecksScopeExpiry(t *testing.T) {
+	d := newHandlerDaemon(t)
+	pinFinishClock(d, finishClockAt)
+	proof := finishProofFile(t, d, "approve", finishClockAt+100, 0o600)
+	w := finishRoute(t, d, FinishRequest{Phase: "approve", Scope: finishScopeAt(finishClockAt + 3600)}, proof, "")
+	var approval FinishApproval
+	if w.Code != 200 || json.Unmarshal(w.Body.Bytes(), &approval) != nil {
+		t.Fatalf("approve: %d %s", w.Code, w.Body.String())
+	}
+	record := &EpisodeFinish{Episode: "ep1", Branch: "world/ep1", Commit: strings.Repeat("a", 40),
+		Remote: "https://github.com/o/r.git", Approval: approval.Approval, Outcome: "push-failed"}
+	for _, arm := range []struct {
+		clock int64
+		want  int
+	}{{-1, 403}, {0, 200}, {3599, 200}, {3600, 403}} {
+		at := finishClockAt + arm.clock
+		pinFinishClock(d, at)
+		proof := finishProofFile(t, d, "record", at+100, 0o600)
+		if w := finishRoute(t, d, FinishRequest{Phase: "record", Record: record}, proof, ""); w.Code != arm.want {
+			t.Fatalf("record at approve%+d: %d, want %d: %s", arm.clock, w.Code, arm.want, w.Body.String())
+		}
+	}
+}
+
 type unreadFinishBody struct {
 	t *testing.T
 }
```

What the diff does (for the report; the diff is authoritative):
- `daemon.go`: the `finishNow func() int64` field and its comment after `drainTimeout` (`:447-456`);
  `finishNow: func() int64 { return time.Now().Unix() },` appended to `New`'s literal line (`:800`).
  The coordinator wiring (`Now: …`) is unchanged.
- `episode_finish.go`: the `"time"` import removed; `now := d.finishNow()` at `:115` (proof),
  `:167` (approve scope), `:269` (record re-check). The bounds are now at `:116`, `:168`, `:270`.
- `episode_finish_test.go`: AC11 `future` `+301` → `+3600` (`:79`); `finishClockAt`,
  `pinFinishClock`, `finishScopeAt`; T1 (`:177`, pin at `:179`, Fatalf `:187`); T2 (`:194`, Fatalf
  `:204`); T3 (`:213`, Fatalf `:232`).

#### M1.4 Compile, vet, the finish tests

```sh
gofmt -l host/daemon/                                         # empty
go vet ./host/daemon/; echo "vet rc=$?"                       # rc 0
perl -e 'alarm 300; exec @ARGV' go test -timeout 280s ./host/daemon/ -run '^TestEpisodeFinish' -count=1 -v > /tmp/w327x/fin-plain.txt 2>&1; echo "rc=$?"
grep -E '^(--- |ok|FAIL)' /tmp/w327x/fin-plain.txt
perl -e 'alarm 400; exec @ARGV' go test -race -timeout 380s ./host/daemon/ -run '^TestEpisodeFinish' -count=1 -v > /tmp/w327x/fin-race.txt 2>&1; echo "rc=$?"
grep -cE '^--- PASS: TestEpisodeFinish' /tmp/w327x/fin-race.txt; grep -cE -- '--- FAIL|DATA RACE' /tmp/w327x/fin-race.txt
```

Want: plain rc 0 with exactly these 5 PASS lines (AC11, ScopeExpiry, OperatorProofExpiryBoundary,
RecordRechecksScopeExpiry, AuthorizationPrecedesBody); race rc 0, `5`, `0`. Planner (V-P5): exactly that.

#### M1.5 AC-1 and AC-1b green on the fix

```sh
W=/Users/voightkampff/dev/sunholo-data/ailang-world/.wt-world-iter327; H=/tmp/w327x/tools/harness.sh; F=host/daemon/episode_finish_test.go
sed -n '79p;179p' $F      # :79 the AC11 future line (+3600); :179 "	pinFinishClock(d, finishClockAt)" inside T1
{ echo "AC1 fix plain:  $($H $W $F 179 20 '^TestEpisodeFinishScopeExpiry$')"
  echo "AC1 fix race:   $($H $W $F 179 20 '^TestEpisodeFinishScopeExpiry$' -race)"
  echo "AC1b fix plain: $($H $W $F 79 20 '^TestEpisodeFinishOperatorProofAC11$')"
  echo "AC1b fix race:  $($H $W $F 79 20 '^TestEpisodeFinishOperatorProofAC11$' -race)"
} > /tmp/w327x/ac1-fix.txt 2>&1; echo "rc=$?"; cat /tmp/w327x/ac1-fix.txt
```

Want: each `rc=0 PASS=20 FAIL=0`, `sha=b872a38cc091/b872a38cc091`. The `ok … <t>s` time is about
20 s or more: the spin really waits for the edge on each of the 20 iterations (that is the
non-vacuity check; a time under 5 s means the spin did not run — stop and check the inserted line).
Planner (V-P6): 22.2 s / 22.9 s / 25.1 s / 23.4 s, all PASS 20.

Snapshot the fixed tree **now** (after G1–G4 are green):

```sh
shasum -a 256 host/daemon/daemon.go host/daemon/episode_finish.go host/daemon/episode_finish_test.go > /tmp/w327x/fixed.sha
cp host/daemon/daemon.go host/daemon/episode_finish.go host/daemon/episode_finish_test.go /tmp/w327x/
```

#### M1.6 AC-2 mutation drill (design §4, every row except TM2)

```sh
perl -e 'alarm 1200; exec @ARGV' /tmp/w327x/tools/allmuts.sh /Users/voightkampff/dev/sunholo-data/ailang-world/.wt-world-iter327 > /tmp/w327x/ac2-mutants.txt 2>&1; echo "rc=$?"
cut -c1-300 /tmp/w327x/ac2-mutants.txt
test -e host/daemon/zz_mut_wall.go && echo "SCRATCH LEFT BEHIND" || echo "scratch absent"
shasum -a 256 -c /tmp/w327x/fixed.sha
```

Each `runmut.sh` line applies one exact replacement on one line (`assert count == 1`), runs
`go test -timeout 180s ./host/daemon/ -run '^TestEpisodeFinish' -count=1 -v` (design AC-2's
command), restores the file by `cp`, and prints `ID rc=… build_err=… fails=[…] first=[…] sha=before/after`.

| ID | Mutant (fixed-tree line) | Named killer (must be in `fails=`, and `first=` as shown) | Planner (V-P7) |
|---|---|---|---|
| P1 | `:168` `scope.ExpiresAt <= now` → `<` | T1 `ScopeExpiry`, `:187: scope expiry now+0: 200, want 403` | killed, as named |
| P2 | `:168` lower arm → `false` | T1, `now-1: 200, want 403` | killed |
| P3 | `:168` `> now+3600` → `>=` | T1 `now+3600: 403, want 200`; also T3 | killed (T1 + T3) |
| P4 | `:168` `3600` → `3601` | T1 `now+3601: 200, want 403` | killed |
| P5 | `:168` `3600` → `3660` | T1 `now+3601: 200, want 403` | killed |
| P6 | `:168` upper arm → `false` | T1 `now+3601: 200, want 403` | killed |
| Q1 | `:270` `<= now` → `<` | T3 `RecordRechecksScopeExpiry`, `:232: record at approve+3600: 200, want 403` | killed |
| Q2 | `:270` lower arm → `false` | T3 `approve+3600: 200, want 403` | killed |
| Q3 | `:270` `> now+3600` → `>=` | T3 `approve+0: 403, want 200` | killed |
| Q4 | `:270` `3600` → `3601` | T3 `approve-1: 200, want 403` | killed |
| Q6 | `:270` upper arm → `false` | T3 `approve-1: 200, want 403` | killed |
| R1 | `:116` `p.ExpiresAt <= now` → `<` | T2 `OperatorProofExpiryBoundary`, `:204: proof expiry now+0: 200, want 403` | killed |
| R2 | `:116` lower arm → `false` | AC11 `expired`; T2 | killed (+ cascade, D4) |
| R3 | `:116` `> now+300` → `>=` | T2 `now+300: 403, want 200` | killed |
| R4 | `:116` `300` → `301` | T2 `now+301: 200, want 403` | killed |
| R6 | `:116` upper arm → `false` | AC11 `future`; T2 | killed (+ cascade, D4) |
| S0 | control: `:115` `d.finishNow()` → itself, scratch `wallNow` file present | none: `rc=0 fails=[]` | rc 0 |
| S1 | `:115` `d.finishNow()` → `wallNow()` | T1, T2, T3 | killed (all three) |
| S2 | `:167` `d.finishNow()` → `wallNow()` | T1, T2, T3 | killed (all three) |
| S3 | `:269` `d.finishNow()` → `wallNow()` | T3 `approve+0: 403, want 200` | killed |
| W1 | `daemon.go:800` `New` wires `return 0` | AC11, `:106: 403` | killed |
| TM1 | `episode_finish_test.go:179` pin → `_ = finishClockAt` | T1 `now+1: 403, want 200` | killed |

Want: 21 lines `rc=1 build_err=0` with the named killer(s) in `fails=`; S0 `rc=0 build_err=0 fails=[]`;
every `sha=` pair equal (`e429375dc3e0` for P/Q/R/S, `e6b70d7e3f3f` for W1, `b872a38cc091` for TM1);
`scratch removed: yes`; final `shasum -c` all `OK`. Any `MUT-APPLY-FAILED` or `build_err` ≠ 0 is a
protocol error: stop and report (do not hand-edit the mutant). Extra failing tests beyond the named
killer are expected only for P3 (T3), R2/R6 (AC11 cascade, D4) and S1/S2 (T1–T3).

### Milestone M2 — AC-3 (package), AC-4 (full gate), executor report

Do not start M2 while a mutant is applied: `shasum -a 256 -c /tmp/w327x/fixed.sha` must say `OK`
for all three files, and `host/daemon/zz_mut_wall.go` must be absent.

#### M2.1 AC-3 (design §5, copied)

Run these **one at a time**, in the foreground.

```sh
go vet ./host/daemon/; echo "vet rc=$?"
s=$(date +%s); perl -e 'alarm 900; exec @ARGV' go test -timeout 880s ./host/daemon/ -count=1 > /tmp/w327x/ac3-plain.txt 2>&1; echo "rc=$? elapsed=$(( $(date +%s) - s ))s"; tail -3 /tmp/w327x/ac3-plain.txt; grep -c -- '--- FAIL' /tmp/w327x/ac3-plain.txt
s=$(date +%s); perl -e 'alarm 1500; exec @ARGV' go test -race -timeout 1450s ./host/daemon/ -count=1 > /tmp/w327x/ac3-race.txt 2>&1; echo "rc=$? elapsed=$(( $(date +%s) - s ))s"; tail -3 /tmp/w327x/ac3-race.txt; grep -c 'DATA RACE' /tmp/w327x/ac3-race.txt
s=$(date +%s); perl -e 'alarm 900; exec @ARGV' go test -race -timeout 880s ./host/daemon/ -run '^TestEpisodeFinish' -count=200 > /tmp/w327x/ac3-race200.txt 2>&1; echo "rc=$? elapsed=$(( $(date +%s) - s ))s"; tail -3 /tmp/w327x/ac3-race200.txt
s=$(date +%s); perl -e 'alarm 600; exec @ARGV' go test -timeout 580s ./cmd/ailang-worldd/ -run Finish -count=1 > /tmp/w327x/ac3-worldd-finish.txt 2>&1; echo "rc=$? elapsed=$(( $(date +%s) - s ))s"; tail -3 /tmp/w327x/ac3-worldd-finish.txt
```

Want: vet rc 0; plain rc 0 and `ok  	github.com/sunholo-data/ailang-world/host/daemon`, 0 FAIL; race
rc 0, the same `ok` line, `DATA RACE` 0; `-count=200` race rc 0 `ok`; `cmd/ailang-worldd` rc 0 `ok`.
Planner (V-P8, V-P9): vet rc 0; plain rc 0, `ok … 335.821s` (337 s wall, load average 14.4), 0 FAIL; race rc 0, `ok … 409.042s` (411 s wall), 0 `DATA RACE`, 0 FAIL; `-count=200` race rc 0, `ok … 52.989s`; `cmd/ailang-worldd -run Finish` rc 0, `ok … 7.861s`, 19 `--- PASS`, 0 FAIL. Expect the package runs to take longer under machine load; the `-timeout`s above leave about 2.5× headroom.

#### M2.2 AC-4 (the full gate)

```sh
s=$(date +%s); perl -e 'alarm 2700; exec @ARGV' bash -c 'go build ./... && go test ./...' > /tmp/w327x/ac4-full.txt 2>&1; echo "rc=$? elapsed=$(( $(date +%s) - s ))s"
grep -E '^(FAIL|ok|---|panic)' /tmp/w327x/ac4-full.txt | grep -v '^ok' | head -40
grep -c '^ok' /tmp/w327x/ac4-full.txt; grep -c '^FAIL' /tmp/w327x/ac4-full.txt
grep -E 'host/daemon' /tmp/w327x/ac4-full.txt
s=$(date +%s); perl -e 'alarm 900; exec @ARGV' ./scripts/verify_ail.sh > /tmp/w327x/ac4-ail.txt 2>&1; echo "rc=$? elapsed=$(( $(date +%s) - s ))s"; tail -3 /tmp/w327x/ac4-ail.txt
```

Want: full gate rc 0 with no `FAIL` line, and `ok  	github.com/sunholo-data/ailang-world/host/daemon`
among the `ok` lines; `verify_ail.sh` rc 0 (no `.ail` is touched). The design did not run AC-4
(design V10), and neither did the planner on the whole repo (V-P12): **your run is the first
full-suite check of this change**; CI on the fix head is the second.

#### M2.3 Telling a pre-existing or load failure from a regression (rule 3e)

The designer saw 13 untriaged failures in **coverage-mode** runs of `cmd/ailang-worldd` (12) and
`cmd/world-publish` (1, `/dev/tty` absent) while an attended trial was loading the machine (design
V15). CI at base `4d43d8e3` is green on its rerun. The planner's own plain baseline of
`go test ./cmd/ailang-worldd/` is V-P13: at base `4d43d8e3` (`git archive` copy, plain, load average 9.7): rc 0, `ok … 90.701s`, 237 PASS / 0 FAIL / 9 SKIP; on the planned code: rc 0, `ok … 89.031s`, 237/0/9, the same 9 SKIPs. So a `cmd/ailang-worldd` FAIL in your AC-4 run is **not** an expected pre-existing failure: classify it with the steps below. (`cmd/world-publish` was not baselined by the planner; the designer's one failure there was `/dev/tty` absent, which is environmental.)

If AC-4 shows any `FAIL`:
1. **Any failure in `host/daemon`, or in a test whose name contains `Finish`, is a stop.** Report it; do not patch around it.
2. Otherwise, for each failing test `<T>` in package `<P>`, re-run it alone in the worktree:
   `perl -e 'alarm 600; exec @ARGV' go test -timeout 580s ./<P>/ -run '^<T>$' -count=3 -v > /tmp/w327x/ac4-rerun-<T>.txt 2>&1; echo "rc=$?"`.
   3/3 PASS → record it as a **load flake** (with the rerun file) and move on.
3. If it fails alone, baseline the **same command at base** in a throwaway copy (never in the worktree):
   ```sh
   rm -rf /tmp/w327x/base && mkdir -p /tmp/w327x/base && git archive 4d43d8e3 | tar -x -C /tmp/w327x/base
   cd /tmp/w327x/base && perl -e 'alarm 600; exec @ARGV' go test -timeout 580s ./<P>/ -run '^<T>$' -count=3 -v > /tmp/w327x/ac4-base-<T>.txt 2>&1; echo "rc=$?"; cd /Users/voightkampff/dev/sunholo-data/ailang-world/.wt-world-iter327
   ```
   Fails at base too → **pre-existing**: record both files, it is not this change's regression.
   Passes at base → **regression**: stop and report.
   Caveat: a `git archive` copy is not a git repository, so a test that shells out to `git`
   (e.g. `host/pins TestReleaseLiteralsExplainedByManifest`, `git ls-files: exit status 128`)
   fails there for that reason alone; such a test passing in the worktree alone (step 2) is enough.
4. The change can only affect a package that imports `host/daemon`
   (`go list -test -f '{{.ImportPath}} {{join .Imports " "}}' ./... | grep -c 'ailang-world/host/daemon\b'`):
   design V15 lists `cmd/ailang-worldd`, `cmd/world-publish` and `host/daemon` itself, and coverage
   shows only `host/daemon`'s own tests execute `episode_finish.go`.

#### M2.4 Executor report and evidence

```sh
for f in ac1-base fin-plain fin-race ac1-fix ac2-mutants ac3-plain ac3-race ac3-race200 ac3-worldd-finish ac4-ail; do
  cp /tmp/w327x/$f.txt design_docs/verification/world-iter327/executor-$f.txt
done
grep -E '^(FAIL|ok|---|panic)' /tmp/w327x/ac4-full.txt > design_docs/verification/world-iter327/executor-ac4-full-summary.txt
cp /tmp/w327x/ac4-rerun-*.txt /tmp/w327x/ac4-base-*.txt design_docs/verification/world-iter327/ 2>/dev/null; true
```

(If you copied any `ac4-rerun-*`/`ac4-base-*` files, rename them with an `executor-` prefix.)
`design_docs/verification/world-iter327/executor-report.md` follows the shape of
`design_docs/verification/world-iter320/executor-report.md`: worktree/branch/base, the base sha256
check, the planned-vs-actual sha256 of the three files, then per milestone a table of every command,
rc and counted result (AC-1/AC-1b base and fix, G1–G3, every AC-2 mutant with its `fails=` and
`sha=`, AC-3's four legs, AC-4's two legs and any §M2.3 classification), and a **Deviations**
section (anything that differed from this plan, with the measurement). State plainly that
`design_docs/world-mission*.md` were not touched.

#### M2.5 Acceptance (M2)

```sh
shasum -a 256 -c /tmp/w327x/fixed.sha                  # OK x3 (no mutant left in)
test -e host/daemon/zz_mut_wall.go && echo LEFT || echo absent
git diff --name-only                                    # exactly the 3 host/daemon files
git diff --stat                                         # 3 files changed, 84 insertions(+), 13 deletions(-)
git status --porcelain                                  # the 3 files, the two plan files, the new executor-* files
git status --porcelain | grep -c 'world-mission'        # want 0
```

---

## 5. Controller-owned items

- **Row 236 charter/queue text** in `design_docs/world-mission.md` and the mission log, from design
  §1/§2 and the executor report.
- Moving the design doc and this plan from `planned/` to `implemented/`.
- Committing, the PR, and reading CI: all 3 jobs green on the fix head, with
  `ok  	github.com/sunholo-data/ailang-world/host/daemon` in the "go race gate" job (design AC-4).
  No CI list edit is needed: no workflow names a finish test (`grep -nE 'EpisodeFinish' .github/workflows/*.yml` → no output; control `grep -c 'TestLiveScriptPure' .github/workflows/ci.yml` ≥ 1, V-P10).

## 6. Risks

| Risk | Mitigation |
|---|---|
| A `Daemon` built by struct literal has a nil `finishNow` and panics on the finish route | Design §3/V18: 11 constructions, only `New`'s reaches the route; no nil guard by design (a silent wall-clock fallback would hide the flake). W1 proves `New`'s wiring is load-bearing |
| Two pinnable notions of "now" in one daemon | Design §2 option (iii): no comparison crosses them; the field's comment says it is not the coordinator's clock |
| A test writes `d.finishNow` while a handler runs | Writes happen only between synchronous `ServeHTTP` calls; constraint 5 (no `t.Parallel`); race runs in M1.4/M2.1 |
| Load-induced failures elsewhere in AC-4 | §M2.3 classification; never patch another test |
| The executor hand-copies the diff and drops tabs or a blank context line | The diff is extracted with `awk` from this file and checked by sha256 before `patch` |

## 7. Definition of done

1. M1: AC-1/AC-1b red at base (4 × FAIL 20, controls PASS 20); the three files equal the planned
   sha256; G1–G3 green; AC-1/AC-1b green on the fix (4 × PASS 20); AC-2: 21 killed as named, S0 green,
   every restore equal, scratch removed.
2. M2: AC-3's four legs green; AC-4 green (or every failure classified per §M2.3, none in `host/daemon`
   or a finish test); executor report and evidence files written; `git diff --name-only` is exactly the
   three `host/daemon` files; no `world-mission*` file touched.

## 8. Deviations from the design (planner-measured; flagged to the controller)

- **D1 (test comments).** `designer-prototype.diff` is the r2 diff for the two production files
  (byte-identical: the planner's `daemon.go`/`episode_finish.go` sha256 prefixes `e6b70d7e3f3f` /
  `e429375dc3e0` equal design V17's). Its test file carried three comment defects, which the plan fixes:
  (a) `finishClockAt`'s comment said `d.now` (the pre-rename name) → `d.finishNow`; (b) T1/T2/T3's
  comments cited base line numbers `:169`/`:117`/`:271`, which are wrong from the moment the `time`
  import goes (the bounds move to `:168`/`:116`/`:270`) → replaced by the function names
  (`handleEpisodeFinish`, `consumeOperatorProof`, `finishDecisionScope`); (c) T3's comment reflowed
  to 80 columns. Comments only; no test logic changed. Consequence: T3's `Fatalf` is at **`:232`**,
  not design §4's `:231`; T1 `:187` and T2 `:204` are as the design says. Test diff `+70/−8`
  (design §0 estimate "about +70/−8"; V9's r2 count said 77 changed lines).
- **D2 (mutant forms made concrete).** The design names mutants by effect. The plan fixes the text:
  "lower arm dropped" = `X <= now ||` → `false ||`; "upper arm dropped" = `X > now+B` → `false`;
  `<=`→`<`, `>`→`>=`, `3600`→`3601`/`3660`, `300`→`301`. The table has no Q5 or R5 (16 bound mutants =
  P1–P6, Q1–Q4+Q6, R1–R4+R6), as design §4 says. Lines are the **fixed-tree** lines (`:116/:168/:270`,
  seam reads `:115/:167/:269`, `daemon.go:800`, test `:179`).
- **D3 (S0 control added).** With the scratch `wallNow` file present and a no-op replacement, the run
  is green (rc 0), so S1–S3's reds come from the mutation, not the scratch file (the designer ran the
  same control, design V8).
- **D4 (R2/R6 cascade).** R2 and R6 also red AC11 subtests `phase`, `readable`, `symlink`,
  `world-readable` (and R2 `expired`, R6 `future`): once one subtest's proof is wrongly accepted, the
  approve writes the approval heads and every later subtest's `refusal wrote …` check fires. The named
  killers are red; the extras are a consequence, not a protocol error. Same as the designer's run.
- **D5 (scratch dir).** The design's `/tmp/w327d` is instantiated as `/tmp/w327x` for the executor
  (the planner used `/tmp/w327p`). The tool scripts are the designer's (`designer-tools.txt`) with
  that path, plus a `build_err` count and an optional extra-flags tail in `runmut.sh`, and restore-on-
  apply-failure.

**Nothing in the design was found false.** Re-measured and held: base sha256 `ea44759896b5` for
`episode_finish.go` (V4); the three reads at base `:116/:168/:270` with bounds at `:117/:169/:271`;
`New`'s literal; AC-1/AC-1b base 20/20 red plain and race, controls 20/20 green (V5); every §4
killer (V8, now on the r2 code); the package green plain and race (V10, now on the r2 code).

## 9. Prototype results summary (r2 code = the planned files)

| Check | Result |
|---|---|
| AC-1 base (`:164`) | plain rc 1 FAIL 20; race rc 1 FAIL 20; first `episode_finish_test.go:171: expiry 1793398298: 200`; control rc 0 PASS 20 |
| AC-1b base (`:79`) | plain rc 1 FAIL 20; race rc 1 FAIL 20; first `episode_finish_test.go:94: 200 {"request":…}`; control rc 0 PASS 20 |
| AC-1 fix (`:179`) | plain rc 0 PASS 20 (22.2 s); race rc 0 PASS 20 (22.9 s) |
| AC-1b fix (`:79`) | plain rc 0 PASS 20 (25.1 s); race rc 0 PASS 20 (23.4 s) |
| AC-2 | 21/21 killed with the named killers, S0 rc 0, 0 build errors, every sha restored |
| AC-3 | vet rc 0; plain `ok … 335.821s`; race `ok … 409.042s`, 0 `DATA RACE`; `-count=200` race `ok … 52.989s`; `cmd/ailang-worldd -run Finish` `ok … 7.861s` (19 PASS). **This is the first package run of the r2 diff** (design AC-3's r3 note) |
| AC-4 | not run by the planner on the whole repo (V-P12) |

## 10. Verification notes (planner, this session, worktree `f7767f12`)

| # | Claim | Command | Observed |
|---|---|---|---|
| V-P1 | Base: the Go tree equals `4d43d8e3` | `git log --oneline -2`; `git status --porcelain`; `git diff --name-only 4d43d8e3 HEAD \| grep -vc '^design_docs/'`; `… \| grep -c '^design_docs/'` | `f7767f12 design(row 236): …`, `4d43d8e3 Merge pull request #432 …`; status empty before the plan files; `0` non-design files, `21` design files |
| V-P2 | Toolchain | `go version`; `$HOME/.pinned-ailang/ailang --version` | `go1.26.6 darwin/arm64`; `AILANG v0.41.0` (commit `24ee108`) |
| V-P3 | Base sha256 and edit sites | `shasum -a 256` of the 3 files; `sed -n '1,20p;105,125p;150,175p;255,280p' episode_finish.go`; `sed -n '436,460p;780,795p' daemon.go` | §3 base column; `time` import `:17`; reads `:116/:168/:270`; `drainTimeout` `:446`; `New`'s literal `:784`, `errLog` line `:790` |
| V-P4 | AC-1/AC-1b at base | `harness.sh /tmp/w327p/base … 164\|79 20 <run> [-race]`, control `SPIN='for time.Now().Nanosecond() > 100_000_000 {}'` | as §9; every `sha=094884e036f7/094884e036f7` |
| V-P5 | Planned code: gofmt, vet, finish tests | `/tmp/w327p/fix` = base + `designer-prototype.diff` (`patch -p1`, clean) + the D1 comment edits; `gofmt -l host/daemon/`; `go vet ./host/daemon/`; `go test ./host/daemon/ -run '^TestEpisodeFinish' -count=1 -v`, and `-race` | gofmt empty; vet rc 0; plain rc 0, 5 PASS (`ok … 0.446s`); race rc 0, 5 PASS, 0 FAIL/`DATA RACE` (`ok … 2.008s`); `diffstat`: `3 files changed, 84 insertions(+), 13 deletions(-)` |
| V-P6 | AC-1/AC-1b on the planned code | `harness.sh /tmp/w327p/fix … 179\|79 20 …` | as §9; every `sha=b872a38cc091/b872a38cc091` |
| V-P7 | AC-2 on the planned code (first run of P/Q/R/TM1 on the r2 diff) | `allmuts.sh /tmp/w327p/fix` (= the §M1.1 block with `/tmp/w327p`) | P1–P6, Q1–Q4, Q6, R1–R4, R6: rc 1, 0 build errors, killers as §M1.6 (T1 `:187`, T3 `:232`, T2 `:204`, AC11 subtests); S0 rc 0; S1, S2: T1+T2+T3; S3: T3 `:232: record at approve+0: 403, want 200`; W1: AC11 `:106: 403`; TM1: T1 `:187: scope expiry now+1: 403, want 200`; sha pairs `e429375dc3e0` ×20, `e6b70d7e3f3f` ×1, `b872a38cc091` ×1, all equal; `scratch removed: yes`; `shasum -c` OK ×3 |
| V-P8 | AC-3 package plain (planned code) | `go test -timeout 570s ./host/daemon/ -count=1` (machine load average 14.4 at start) | rc 0, `ok  github.com/sunholo-data/ailang-world/host/daemon 335.821s` (337 s wall), 0 FAIL |
| V-P9 | AC-3 package race, `-count=200`, worldd Finish (planned code) | `go test -race -timeout 1450s ./host/daemon/ -count=1`; `go test -race ./host/daemon/ -run '^TestEpisodeFinish' -count=200`; `go test ./cmd/ailang-worldd/ -run Finish` | race: rc 0, `ok  github.com/sunholo-data/ailang-world/host/daemon 409.042s` (411 s wall), 0 `DATA RACE`, 0 `--- FAIL`; `-count=200`: rc 0, `ok … 52.989s`; worldd: rc 0, `ok … 7.861s`, 19 `--- PASS`, 0 `--- FAIL` |
| V-P10 | Charter row 236 absent; no CI list names a finish test | `grep -cE '^236\. ' design_docs/world-mission.md`; `grep -nE '^23[5-9]\. ' …`; `grep -nE 'EpisodeFinish\|finish' .github/workflows/*.yml`; control `grep -n 'TestLiveScriptPure' .github/workflows/ci.yml` | `0`; `6304:235.` only; no output; control `:269` |
| V-P11 | The plan's fenced blocks reproduce the planned files and tools byte-for-byte | extract the `diff` block and the 4 tool blocks from this file with §1's `x()`; `git archive HEAD` into `/tmp/w327p/rt`; `patch -p1` there; `shasum`; `allmuts.sh` (extracted, `/tmp/w327x` paths) not re-run | all 4 tool blocks and the diff block `cmp`-identical to the planner's files (extracted under bash and under zsh: same sha256); extracted diff sha256 `d36610127fb2ebb075fe828a4f7bb0bc8802a81b71af22ed297a7db0c9f552eb`; `patch -p1 --dry-run` rc 0 (3 `patching file`, no offset/fuzz), `patch -p1` rc 0; the three files `shasum -c` OK against the planned column of §3; 0 `.orig`/`.rej`. The executor's tools differ from the ones the planner ran only by the path `/tmp/w327p/` → `/tmp/w327x/` (`diff` after `sed` back: identical) |
| V-P12 | Full gate | not run on the whole repo by the planner (a `git archive` copy fails `host/pins` for want of `.git`, prior plan V-P8b); the executor's AC-4 is the first | — |
| V-P13 | `cmd/ailang-worldd` plain baseline at base `4d43d8e3` (rule 3e) | `/tmp/w327p/base` (`git archive`); `go test -timeout 900s ./cmd/ailang-worldd/ -count=1` | base: rc 0, `ok  github.com/sunholo-data/ailang-world/cmd/ailang-worldd 90.701s` (94 s wall, load average 9.7), **237 PASS, 0 FAIL, 9 SKIP**. The same command on the planned code (`/tmp/w327p/fix`): rc 0, `ok … 89.031s`, 237/0/9, identical SKIP set (`TestCLIAgainstSeToolsDaemon`, `TestCLINoGenesisIsNoWorldHead`, 4 × `TestDemoSeed*`, `TestEpisodePtyChildProcess`, `TestMarketingChromeSixShots`, `TestSessionPtyChildProcess`). So the designer's 12 `cmd/ailang-worldd` failures (design V15) do not reproduce in plain mode at base; they belong to coverage mode and/or the concurrent attended load, and remain untriaged |
