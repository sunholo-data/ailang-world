# Sprint plan — w-exec-probe-cross-profile-net6-hoist (iteration 340, row 242)

Approved design: [w-exec-probe-cross-profile-net6-hoist.md](w-exec-probe-cross-profile-net6-hoist.md),
r3, including §11's controller narrow-refinement carve-out. Base:
`45bbc95af54bb41c622f469c0cd2c8e40af5f331` (detached planner worktree).
One milestone, **M1**, estimated 0.25 day. Status: **planned**. The planner only banks the plan;
the executor applies the diff and checks hashes; the controller runs listener tests and mutations
outside the sandbox and commits. No role instruction here authorizes the planner to implement or
perform Git writes. The prior iteration's shape is the reference, not its multi-milestone scope.

M1 binds every profile's arm-6 listener before the first goroutine and closes them in the parent
only after `wg.Wait()`. A failed bind stays in its owner's listener and is judged by the existing
control leg. The private parameter plumbing changes; exported signatures and arm diagnostics do
not. This belongs at the Go effect boundary: socket ownership and goroutine ordering cannot be
an AILANG policy package. No new seam or production context root is introduced.

The executor must apply the banked patch without rewriting either Go file. Design §6 R4 and
§11 r3 explicitly require V11 byte-exact; Linux `SO_ACCEPTCONN` is a follow-up only if CI exposes
a false match. Do not run `gofmt -w`, embellish the tests, or alter the ownership design.

## M1 — exact files and identities

The implementation file list is exactly these three paths. The new file has no pre-change
bytes, so its pre-change identity is **ABSENT**, not the hash of an empty file.

| File | Change | Pre-change sha256 | Post-change sha256 |
|---|---|---|---|
| `host/broker/exec_probe.go` | Pinned hoist and private parameter plumbing | `96f33172667d27defa0c52735d5cff1ad74d4e907ab00a4c056c4e5a26f698a4` | `19de369381579387bd9d83aceae28a5912717dab00b283454ca0067d41d087bd` |
| `host/broker/exec_probe_net6_test.go` | Pinned 216-line T1/T2 test file | `ABSENT` | `400faf0ffc25a3f0a79302cc9436349da6e445b998744b5767f6d706475a64aa` |
| `.github/workflows/ci.yml` | Append two names to the existing reverse-leg tests= line | `a11ff4ab05a41c9864299edb8f636675e3abb5dcdcba974d43fb6a405940710f` | `41a0a58611565fe148faaac0327dc34793550f12ddef74e7cc3c33a250c4697d` |

Bank: `design_docs/verification/world-iter340/plan-M1.diff.txt`, sha256
`3d08d627134fd1062e066ce287ec64a96f038f108dea67faee96537eddf708a2`.
It starts with the complete designer patch, byte for byte, then appends exactly one CI hunk.
The CI addition is on the existing `$tests` line ending in
`TestExecProbeChildListenerDrainsTheHello`, with the two new names appended in T1/T2 order.
The four planner deliverables are documentation/evidence artifacts, outside this three-file
implementation list. The existing `host/broker/exec_probe_test.go` stays byte-identical.

## Acceptance gates

Run the following blocks at the applied tree root, in Bash with `set -euo pipefail`.
A0/A1 run once before application; A2–A7 run after application. The controller may apply A1 before
handing the tree to the executor; do not apply twice. A hash mismatch is a finding: stop and
report it, never repair the bank or implementation by hand. Re-run A2 after every mutation run.

Environment setup for the controller's gates (all evidence paths remain inside this tree):

```bash
set -euo pipefail
export AILANG_BIN="$HOME/.pinned-ailang/ailang"
"$AILANG_BIN" --version
go version
node --version
export EVIDENCE="$PWD/design_docs/verification/world-iter340/controller-M1"
mkdir -p "$EVIDENCE/tmp"
export TMPDIR="$EVIDENCE/tmp"
uptime > "$EVIDENCE/load.txt"
ps -Ao pid,pcpu,comm > "$EVIDENCE/processes.txt"
```

Use the controller's existing released AILANG pin, never a dirty build. Node is mandatory for
T1: a skip cannot satisfy A4. Real-srt Linux CI supplies its existing pinned environment; the
PASS-line loop must actually report both added names. A7 verifies wiring only; successful Linux
CI remains required and is **UNMEASURED** by the planner. Keep suite runs and mutation runs
sequential. Take any controller baseline in a pristine archive copy under the tree, not in a
concurrently edited checkout. Save rc immediately for any command whose result needs diagnosis;
`set -e` makes these blocks stop on unexpected nonzero results.

**A0 — base pins, before applying.**

```bash
test "$(git rev-parse 45bbc95af54bb41c622f469c0cd2c8e40af5f331^{commit})" = 45bbc95af54bb41c622f469c0cd2c8e40af5f331
test ! -e host/broker/exec_probe_net6_test.go
shasum -a 256 -c <<'SHA'
96f33172667d27defa0c52735d5cff1ad74d4e907ab00a4c056c4e5a26f698a4  host/broker/exec_probe.go
a11ff4ab05a41c9864299edb8f636675e3abb5dcdcba974d43fb6a405940710f  .github/workflows/ci.yml
SHA
```

**A1 — bank identity and apply, once at the base.**

```bash
shasum -a 256 -c <<'SHA'
3d08d627134fd1062e066ce287ec64a96f038f108dea67faee96537eddf708a2  design_docs/verification/world-iter340/plan-M1.diff.txt
SHA
python3 -I - <<'PYCODE'
from pathlib import Path
b = Path('design_docs/verification/world-iter340')
p = (b / 'designer-proto.patch.txt').read_bytes()
s = (b / 'plan-M1.diff.txt').read_bytes()
assert s.startswith(p), 'designer patch bytes changed'
assert s[len(p):].startswith(b'--- a/.github/workflows/ci.yml\n')
assert s[len(p):].count(b'@@ ') == 1, 'expected exactly one CI hunk'
PYCODE
git apply --check design_docs/verification/world-iter340/plan-M1.diff.txt
git apply design_docs/verification/world-iter340/plan-M1.diff.txt
```

**A2 — post-apply pins.**

```bash
shasum -a 256 -c <<'SHA'
19de369381579387bd9d83aceae28a5912717dab00b283454ca0067d41d087bd  host/broker/exec_probe.go
400faf0ffc25a3f0a79302cc9436349da6e445b998744b5767f6d706475a64aa  host/broker/exec_probe_net6_test.go
41a0a58611565fe148faaac0327dc34793550f12ddef74e7cc3c33a250c4697d  .github/workflows/ci.yml
SHA
```

**A3 — formatting, compile/vet, build, runner syntax.**

```bash
gofmt -l host/broker/exec_probe.go host/broker/exec_probe_net6_test.go > "$EVIDENCE/gofmt.log"
test ! -s "$EVIDENCE/gofmt.log"
go vet ./host/broker/
go build ./...
bash -n design_docs/verification/world-iter340/plan-mutants.sh.txt
```

**A4 — both new tests individually; outside sandbox.**

```bash
for name in TestExecProbeNet6BoundBeforeAnyPick TestExecProbeNet6BindSiteIsBeforeTheFanOut; do
  log="$EVIDENCE/$name.log"
  rc=0
  go test -count=1 -run "^${name}$" ./host/broker/ -v -timeout 120s > "$log" 2>&1 || rc=$?
  cat "$log"
  test "$rc" -eq 0
  /usr/bin/grep -Fq -- "--- PASS: $name " "$log"
  if /usr/bin/grep -q -- '^--- SKIP:' "$log"; then exit 1; fi
done
```

**A5 — mutation matrix; controller outside sandbox, no concurrent suite.**

```bash
bash design_docs/verification/world-iter340/plan-mutants.sh.txt "$PWD" > "$EVIDENCE/mutants.log" 2>&1
cat "$EVIDENCE/mutants.log"
test "$(/usr/bin/grep -c '^EXPECTATION MATCH:' "$EVIDENCE/mutants.log")" -eq 27
if /usr/bin/grep -q '^EXPECTATION MISMATCH:' "$EVIDENCE/mutants.log"; then exit 1; fi
test "$(/usr/bin/grep -c '^RESTORE cmp ok: .*TestExecProbe' "$EVIDENCE/mutants.log")" -eq 27
```

**A6 — full broker and pinned repo gates; controller outside sandbox.**

```bash
go test -count=1 -timeout 600s ./host/broker/ > "$EVIDENCE/broker.log" 2>&1
cat "$EVIDENCE/broker.log"
./scripts/verify_ail.sh > "$EVIDENCE/verify-ail.log" 2>&1
cat "$EVIDENCE/verify-ail.log"
./scripts/verify_go.sh > "$EVIDENCE/verify-go.log" 2>&1
cat "$EVIDENCE/verify-go.log"
```

**A7 — CI name wiring; final A2 hash gate again.**

```bash
python3 -I - <<'PYCODE'
from pathlib import Path
s = Path('.github/workflows/ci.yml').read_text()
start = s.index('      - name: Row 140 Workspace.Exec real-srt tests, verbose')
stop = s.index('      - name:', start + 1)
block = s[start:stop]
lines = [l for l in block.splitlines() if l.lstrip().startswith('tests=')]
for name in ('TestExecProbeNet6BoundBeforeAnyPick', 'TestExecProbeNet6BindSiteIsBeforeTheFanOut'):
    assert sum(l.count(name) for l in lines) == 1, name
assert 'for t in $tests; do' in block
assert '^--- PASS: $t ' in block
assert 'go test ./host/broker/ ./host/daemon/' in block
PYCODE
```

A6 uses `verify_go.sh` for the design’s `go build ./...` and `go test ./... -count=1`
requirements, plus its required race/control checks, without duplicating the repo suite.
The base commit object and source hashes are pinned; documentation-only planning commits may
advance HEAD before execution. A6 is the full verification gate, not a claim that all historical skips become passes. Inspect
named verdicts and required identities in its logs and Linux CI. If an unexpected red occurs,
record load, compare the affected package against base, and rerun that named test and package
in isolation outside the sandbox before attributing it. The R2 pass-through flake stays a
reported residual; do not silently widen its assertions or substitute a passing rerun for evidence.

## Mutation matrix

T1 = `TestExecProbeNet6BoundBeforeAnyPick` (real production wiring, two profiles).
T2 = `TestExecProbeNet6BindSiteIsBeforeTheFanOut` (pure lexical source pin).
TP = `TestExecProbePassThroughShimFailsEveryConfinementArm`.
**GREEN = rc 0 plus the named PASS; RED = rc 1 plus the named FAIL and the indicated diagnostic**.
A compile failure, panic, timeout, skip, or zero selected tests is an expectation mismatch.

The table mirrors design §5's named-test expectations; the final column preserves the design's
historical family results for context. The runner never uses the family as a mutation gate.
TP is omitted in all scenarios except MUT-6: its known roughly 5% flake would make even an
expected-green TP row spuriously print MISMATCH on red. Thus CONTROL's TP result from the design
is historical evidence only. MUT-6's expected-red TP row is supported by the design's 3/3 reruns.

| Scenario | Mutation from pristine applied tree | T1 | T2 | TP in runner | Historical `-run ExecProbe` (design, not rerun) |
|---|---|---|---|---|---|
| CONTROL | None | GREEN | GREEN | Omitted (design GREEN) | GREEN |
| HEAD-positive-control | Replace production file with pinned base; keep new tests | RED: missing live port | RED: call-site pin | Omitted | GREEN excluding T1/T2 |
| HEAD-SLEEP | Base plus 2 s sleep before old bind | RED: missing live port | RED: call-site pin | Omitted | GREEN excluding T1/T2 |
| MUT-1 | 2 s sleep in parent before bind loop | GREEN | GREEN | Omitted | GREEN |
| MUT-1b | 2 s sleep at old bind site | GREEN | GREEN | Omitted | GREEN |
| MUT-2 | Rebind in `probe`, defer local close; retain unused parent listeners as designer did | RED: missing live port | RED: two call sites | Omitted | RED |
| MUT-3 | Remove parent bind loop; bind in goroutine before profile work | GREEN | RED: function literal / after go | Omitted | RED |
| MUT-3b | Rebind after `NewExecHandler`; retain parent listeners as designer did | RED: missing live port | RED: two call sites | Omitted | RED |
| MUT-4 | Parent close replaced by `_ = l` | RED: listeners left after return | GREEN | Omitted (initial red unattributed, 3/3 green reruns) | RED |
| MUT-5 | Pass `net6[0]` to every profile | RED: want two distinct ports | GREEN | Omitted | RED |
| MUT-6 | Close all listeners before `wg.Wait()` | RED: missing live port | GREEN | RED: does not report accepts 2/3 legs | RED |
| MUT-7 | Remove parent bind loop; call another file's `bindNet6Late` from goroutine | GREEN | RED: helper call site | Omitted | Not measured |
| MUT-8 (planner addition) | Keep parent binds; rebind at old site via package-level function-valued variable in helper file | RED: missing live port | GREEN: scan skips variable declarations | Omitted | Not measured |

MUT-8 is a new prediction from source inspection, **UNMEASURED** until the controller runs it.
It discharges the missing variable-initializer case in the source pin's all-files claim without
changing the approved bytes. MUT-7 exercises a real function declaration in another file; MUT-8
uses `var planLateNet6 = func() *probeListener { return newProbeListener() }`. Both keep the
production call chain, and T1 still observes the real old-site replacement in MUT-8.

The banked runner takes the tree root as `$1`, resolves it, checks pristine/base hashes, and
runs 13 scenarios as **27 independent named-test rows** (T1+T2 each, plus MUT-6/TP). It applies
one scenario for one test at a time, using `python3 -I` with unique anchors, and invokes
`go test -count=1 -run '^<name>$' ./host/broker/` with verbose evidence and a 120 s timeout.
It prints subject-qualified `EXPECTATION MATCH` or `EXPECTATION MISMATCH`, continues after test
mismatches, and returns nonzero if any mismatch occurred. Every row restores source and test
from saved pristine copies and checks both with `cmp`; helper absence is also checked. EXIT and
signal traps repeat restoration. Restore failure exits 3. Logs and runner caches are retained in
`$1/.plan-mutants.XXXXXX` for review. Do not run it concurrently with other edits or test suites.

S6 load-bearing assertions: HEAD/MUT-2/MUT-3b/MUT-6 pin T1's per-pick port membership; MUT-4 pins
its after-return listener inventory; MUT-5 pins its two-distinct-ports assertion. MUT-2/-3/-3b/-7
pin T2's final call-site assertion, with MUT-7 proving the cross-file scan. MUT-1/-1b are harmless
slow-step controls. T1's existing first-pick/Git-in-flight guard remains a non-vacuity check;
no additional claim is made that this guard itself was mutation-tested here. TP under MUT-6
pins its existing `does not report` checks of raw/proxy accepts. A7 is instrument health only.

## Residuals carried unchanged from design §6

- **R1 — cross-process binds/dials.** Hoisting excludes in-process picks only. The rig's
  `serve.py` dialer reaches loopback listeners about every 31.8 s. Bind-to-acceptsC remains a
  fail-closed exposure; hoisting adds roughly 7 ms median (unprotected) plus protected pre-bind
  Git time, about +1.2% of the measured existing window. The design estimates about +0.02%
  refusal probability per unprotected profile/startup against that dialer. Foreign accepts
  before control refuse via raw/proxy legs; moving control into the parent does not help.
  Cross-process port theft also remains. Linux UNMEASURED.
- **R2 — TP accepts-4 flake.** Design records 3/63 reds, roughly 4.8%. One profile's
  acceptsA-to-acceptsC window is unchanged. The periodic dialer explains the shape but was not
  caught during a flake; MUT-4's initial TP red is unattributed. This row does not fix it.
- **R3 — lexical source pin.** Function-value indirection can defeat a lexical claim; removing
  the one visible direct call reds T2, while raw `net.Listen` behind an ignored parent listener
  can escape it and should red T1. Planner MUT-8 extends this residual to package variables.
- **R4 — descriptor scan.** On Darwin it also matches bound stream sockets without peers that
  are not listening. False green needs a future-port coincidence; missed descriptors beyond
  min(4096, RLIMIT_NOFILE) can cause false red. Preserve the exact patch. Linux is UNMEASURED;
  follow up only if Linux CI reveals a false match, per the verbatim quorum r3 instruction.
- **R5 — lifetime/resource extension.** At most N listeners are held for the entire fan-out,
  including profiles that returned early. V14's suite differences were transient and both
  reran green (V21). V22's mixed bind-failure/early-return outputs were byte-equal with no
  listeners left. Those are historical design evidence, not planner reruns.

## M1 must NOT touch

- Any path outside the three-file implementation list; in particular the existing
  `host/broker/exec_probe_test.go`, other test helpers, or a permanent mutant helper file.
- The approved design, this plan, the JSON, either designer tool, prototype patch, or banked
  `plan-*.txt` while executing. Evidence uses new controller/executor log paths.
- Reverse-leg retry/reset semantics (row 233, D-WORLD-105), pick function, hello wait, cadence,
  child script, arm strings/order/counting/control placement, or exported API signatures.
- Any `.ail`, kernel/package projection, `go.mod`/`go.sum`, QUICKSTART, daemon or startup caller.
- Linux-specific descriptor refinements, new listener seams, `t.Parallel()`, cross-process
  scanner mitigation, or TP-flake fixes. Findings become follow-ups, never patch drift.
- Mission charter/log/status/dashboard/index, shared driver `tools/launchd/*`, skills,
  the main checkout, V1 state, or any Git index/history/worktree metadata. Controller commits.

## Planner verification and findings

The scratch copy was created at this worktree's `.plan-scratch/` with exactly:

```bash
mkdir .plan-scratch
bash -o pipefail -c 'git archive HEAD | tar -x -C .plan-scratch/'
```

Inside it, `git apply --check ../design_docs/verification/world-iter340/plan-M1.diff.txt`
returned **rc 0**. Applying the same file returned **rc 0**. Post-apply sha256 values are the
three measured values above. The complete designer patch remains a byte-exact prefix; its
new-file hunk reconstructs the test file exactly. No Git index or history was written.

| Planner check | Pristine archive | Applied archive |
|---|---|---|
| `go vet ./host/broker/` | rc 0 | rc 0 |
| `go build ./...` | rc 0 | rc 0 |
| `gofmt -l` on both patched Go files | Not applicable to absent new file | rc 0, empty output |
| `/bin/bash -n plan-mutants.sh.txt` | Documentation artifact | rc 0; Bash 3.2 |

Toolchain: go1.26.6 darwin/arm64. Go vet compiles test sources as well as production; build alone
does not prove test compilation. A first vet launch had an invalid nested TMPDIR setup and
stopped before analysis; after creating `.plan-scratch/.planner-tmp`, both actual vet checks
returned 0. Scratch logs/temp files stayed inside the worktree and the scratch tree was deleted.
JSON validation, all eight acceptance blocks’ Bash syntax, embedded Python compilation, the
read-only patch-prefix gate and the CI wiring gate also returned rc 0.
No broker test, package suite, or mutation runner was executed by the planner. Any listener-test
result under this sandbox is **UNINFORMATIVE UNDER SANDBOX**, never a pass or fail. The pinned
AILANG, full-suite, runtime mutations, and Linux gates remain the controller's acceptance work.

Findings, carried without changes to the approved implementation:

1. **T2's all-package-call-site claim is broader than its AST walk.** It examines only
   `*ast.FuncDecl`; package-level `GenDecl` initializers are skipped. MUT-8 predicts T2 GREEN
   and T1 RED for a late factory bind. This is a lexical residual, not an approved-patch edit.
2. **MUT-2 and MUT-3b are labelled moves in the design but are extra binds in its actual tool.**
   They retain the hoisted parent binds and shadow/replace the per-profile pointer, closing
   the replacement locally. The runner intentionally matches the measured tooling. MUT-3/-7
   exercise true relocation. The numeric table expectations remain as designed.
3. **Literal all-process write confinement conflicts with existing tests.** Runner-controlled
   output, TMPDIR, GOPATH, GOCACHE and GOMODCACHE stay under `$1`. Existing `probe()` and
   `plantExecProbe()` explicitly write `/tmp/claude`, independent of TMPDIR. A portable runner
   cannot redirect those writes without changing approved code or providing OS-level path
   isolation. This is unresolved for a literal interpretation of “write nothing outside that
   root”; no planner command ran those tests or wrote that path.
4. `git show e4e3990c --stat` contains the iteration-339 mission records, not the banked files
   themselves. Its log entry identifies planner commit `68580b0a` and their provenance. This
   does not affect the iteration-340 byte pins.
