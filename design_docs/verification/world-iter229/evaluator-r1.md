# World iteration 229: independent evaluator verdict (round 1)

- **Subject:** `sprint/w-verifygate-forklock-guard-residuals` at `733dffc` (PR #198). Base: `origin/dev` `2d3a255`. Row 142.
- **Evaluator:** Agent tool, `opus`. It ran in a fresh context with a scratch worktree for its drills (`.eval-world-iter229`, detached at `733dffc`, removed afterwards). It did not touch the executor worktree `.wt-world-iter229`. The executor ran on `sonnet`, so the generator and the judge are different models.
- **Verdict:** **PASS 92/100, zero blocking.**

## Scope (item 4)

`git diff --name-only 2d3a255..733dffc` lists 7 paths:
- `host/verifygate/forklocked_write_test.go`
- `.github/workflows/ci.yml`, with one step body changed
- 5 files under `design_docs/` (the design, the plan, the executor report and 2 quorum JSONs)

There is no production code, no `.ail` and no `tools/launchd/*`.

## Gates (evaluator, at 733dffc; `AILANG_BIN=~/.pinned-ailang/ailang` = AILANG v0.41.0, go1.26.6 darwin/arm64)

| Gate | Result |
|---|---|
| `gofmt -l host/verifygate` | empty |
| `go vet ./host/verifygate/` | rc 0 |
| `go test ./host/verifygate/ -count=1 -v` | rc 0. 61 top-level `--- PASS`, 1 `--- SKIP` (the kernel control, darwin) and 0 FAIL. `ok … 142.406s` |
| `go test -race ./host/verifygate/ -count=1 -run 'TestForkLocked\|TestVerifygateTestWritesAreForkLocked\|TestNoParallel'` | rc 0. 4 PASS (ConcurrentFork 0.27s, WrappersHoldLockThroughClose 0.01s, WritesAreForkLocked 26.93s, NoParallel 0.63s). 0 `DATA RACE`. `ok … 29.276s` |

## Mutation table (scratch worktree)

Each mutant was applied by script and then run against the 4 fast tests. It was reverted with `git checkout -- host/verifygate`, and `git diff --quiet` plus `git status --porcelain` came back clean after every revert. At the end, `git diff --quiet 733dffc` was IDENTICAL.

### iter-228 survivors (all must die)

| # | Mutant | Result | Firing line |
|---|---|---|---|
| E1 | `createTempForkLocked` does raw `os.CreateTemp`+`Write`+`Close` with no `forkLockedDo` | **KILLED** (twice) | probe: `forklocked_write_test.go:274: createTempForkLocked: fork-lock probe stages = [], want [filled closed]`; scan: `forklocked_write_test.go:109: raw os.CreateTemp in func createTempForkLocked` |
| E2 | `forkLockedWrite`: `if mode&0o111 != 0 {` raw OpenFile/fill/Close `}` | **KILLED** (twice) | probe: `writeFileForkLocked / copyFileForkLocked / forkLockedWrite: fork-lock probe stages = [], want [filled closed]`; scan: `forklocked_write_test.go:73: raw os.OpenFile in func forkLockedWrite` |
| E3 | `defer RUnlock` removed; `RUnlock()` before `f.Close()` (and on the error paths) | **KILLED** | `forklocked_write_test.go:278: <each of 5 wrappers>: ForkLock released before Close returned`. ConcurrentFork and the scan stayed PASS, as designed |
| E5b | new helper `func evalH229() { wf := os.WriteFile; _ = wf("x", nil, 0o755) }` in `ail_binary_gate_test.go` | **KILLED** | `ail_binary_gate_test.go:763: raw os.WriteFile in func evalH229` |
| E8b | new helper with `syscall.Open` + `os.NewFile` (and a `syscall` import) in `toolchain_pin_gate_test.go` | **KILLED** | `toolchain_pin_gate_test.go:1780: raw syscall.Open in func evalRaw229` and `:1781: raw os.NewFile in func evalRaw229` |
| E6b | scan skips package-level uses (`if where == "package-level initializer" { continue }`) | **KILLED** | `fixture PkgInit: fixture_PkgInit.go:5\|package-level initializer not reported (got [])` |

### New mutants aimed at the new code

| # | Mutant | Result | Firing line / note |
|---|---|---|---|
| N1 | probe never installed (`forkLockedProbe.Store(&probe)` → `_ = probe`) | **KILLED** | `:273: <5 wrappers>: fork-lock probe stages = [], want [filled closed]` |
| N2 | "closed" probe call moved before `f.Close()` | **KILLED** | `:276: <5 wrappers>: probe at closed stage saw an open file` |
| N3 | exemption callee matched by bare name (`fid.Name != "forkLockedDo"`) instead of `info.Uses[fid] != doObj` | **SURVIVED** | No fixture defines a non-package-scope `forkLockedDo` (a local or shadowed one) whose open closure must not be exempt. Low impact: a shadowed `forkLockedDo` would have to be written on purpose |
| N4 | exemption widened to every selector inside the open closure | **KILLED** | `fixture Smuggle: fixture_Smuggle.go:9\|raw os.WriteFile in func h not reported (got [])` |
| N5 | opener-name check dropped (any `os` func returned from the closure is exempt) | **SURVIVED** | No fixture returns a non-OpenFile/CreateTemp opener. Low impact: the closure still runs under RLock, and the matching `syscall.Open` is still reported |
| N6 | tripwire dir exemption broadened (`dir != "host/verifygate"` → `!strings.HasPrefix(dir, "host/")`) | **SURVIVED** | The known positive uses a synthetic dir `zz`, and the live tree has 0 violations, so no assertion pins the exemption to exactly `host/verifygate`. A one-line change can silently neuter the R2 tripwire for every host package |
| N7 | tripwire write detection: `os` dropped from `writeSel` | **KILLED** | `:723: sweep saw 0 write+fork dirs, want >= 12` |
| N8 | tripwire fork detection: `exec` dropped, and the write+fork floor lowered to 0 | **KILLED** | `:735: known positive not named: []` |
| N9 | CI step: `-run` list drifts from the loop list (`TestNoParallel…` dropped from `-run`). Run locally with `RUNNER_OS=Linux` and the kernel test removed from `tests` so darwin can go green | **KILLED** | Control (no drift): rc 0. Mutant: rc 1, `::error::TestNoParallelWriteForkPackagesOutsideVerifygate did not PASS on linux (SKIP?)` |
| N10 | live escape: new helper `os.CopyFS(dst, os.DirFS("testdata"))` in `ail_binary_gate_test.go` | **SURVIVED** | `os.CopyFS` opens write fds inside package `os`, with no ForkLock. Per `go doc`, it "creates files with mode 0o666 plus any execute permissions from the source", so it is a natural way to stage executable fixtures. It is outside the D2 banned set |
| N11 | live method value `c := r.Create` on `*os.Root` | **KILLED** | `ail_binary_gate_test.go:763: raw os.Create in func evalRoot229` |
| N12 | syscall banned set dropped from the scan | **KILLED** | `fixture RawFd: fixture_RawFd.go:9\|raw syscall.Open in func h not reported` |
| N13 | fork inside `fill` (`exec.Command("/usr/bin/true").Run()` in `writeFileForkLocked`), the R1 deadlock | **KILLED by timeout only** | `panic: test timed out after 30s`. This confirms R1: a self-deadlock is caught only by `-timeout` (120 s in CI), not by any assertion |
| N14 | kernel-control exemption chosen by name prefix (`Test*`) instead of `info.Defs` identity | **KILLED** | `forklocked_write_test.go:314: raw os.OpenFile in func TestKernelRefusesExecOfWriterOpenFile` |
| N15b | `RUnlock()` before `f.Close()`, `RLock()` again after it (the lock is dropped only across Close) | **SURVIVED** | The probe samples the lock at "filled" and "closed", not during Close. A fork in the Close window would inherit the fd. Contrived (nobody re-locks after Close by accident), but it shows that D1's "held at and after Close" is a two-point sample |
| N16 | Uses scan narrowed back to call `Fun` idents only | **KILLED** | `fixture FuncValue: fixture_FuncValue.go:5\|raw os.WriteFile in func h not reported (got [])` |

**Totals:** 6 of 6 iter-228 survivors KILLED. 16 new mutants: 11 KILLED (N13 by timeout only) and 5 SURVIVED (N3, N5, N6, N10, N15b).

N1 and N14 were first applied in a non-compiling form (`declared and not used`) and were re-applied correctly. The killed results above are from the compiling runs.

**Equivalent mutant noted:** an exemption that accepts any `FuncLit` argument, not only `Args[0]`, is type-equivalent. `fill` is `func(*os.File) error` and `hold` is `func()`, so neither can contain `return os.OpenFile(…)` as a sole return. It was not run.

## CI step semantics (item 3)

I read `ci.yml` at `733dffc`; both jobs are `runs-on: ubuntu-latest`, and the default shell is bash.

- **rc handling:** `set -euo pipefail`, then `go test … >"$log" 2>&1 || rc=$?`, `cat "$log"` and `[ "$rc" -eq 0 ] || exit "$rc"`. No pipe can eat go test's rc, and the loop's `grep -q` reads a file.
- **One list:** the `-run` regex `"^(${tests// /|})\$"` is built from the same `$tests` string the loop iterates over. Drill N9 shows that drift is caught.
- **Linux precondition fails loud:** with `RUNNER_OS` unset, rc 1 and `::error::this step must run on linux`.
- **SKIP is caught:** with `RUNNER_OS=Linux` on darwin, rc 1 and `::error::TestKernelRefusesExecOfWriterOpenFile did not PASS on linux (SKIP?)`.
- **Probe race safety:**
  - The probe is published through `atomic.Pointer`.
  - It is installed only by a serial test, with a `t.Cleanup` that restores nil.
  - It filters on the test's own `t.TempDir()` prefix before it touches the captured slices.
  - The wrappers are called synchronously, so the slices are touched only on the test goroutine.
  - The `-race` run is clean.
- **No fork in the RLock region:** the probe does `TryLock`/`Unlock`/`Stat`/`Name` only, and so does the rest of the region.
- **Deadlock residual:** a future fork inside the region is caught only by timeout (N13, R1).

## Blocking findings

None.

## Non-blocking findings

1. **`os.CopyFS` evades the scan (N10).** It writes files, keeps exec bits, and opens its fds inside package `os` without ForkLock. The fix is to add `CopyFS` to `bannedOSFuncs` and to the tripwire's `writeSel`, plus a P-CopyFS fixture. This is the same class as E5b/E8b, and it is now the cheapest live escape.
2. **The tripwire's directory exemption is unpinned (N6).** Nothing asserts that a non-`host/verifygate` real-shaped dir such as `host/x` is named, or that `host/verifygate` itself is exempt. Fix: add a positive keyed `host/zz` and a negative keyed `host/verifygate` to `parallelWriteForkViolations`'s fixtures.
3. **Exemption identity and opener-name are unpinned (N3, N5).** Fix: add a fixture with a local `forkLockedDo := func(…)` (it must report) and one whose closure returns `os.Create(…)` (it must report, if the narrow opener set is intended). Low impact.
4. **The probe is a two-point sample (N15b).** It shows the lock is held at "filled" and at "closed", not continuously across `Close`. This is acceptable for the realistic E3 shape. Record it as a known limit of D1.
5. **AC10's N-OpenClosure half was not observable** (plan override 6). **AC8 was shown** through the local darwin known-positives and the CI PASS loop, not through a scratch PR commit with a forced skip. Both are acceptable as recorded by the executor.
6. **Cost:** `TestVerifygateTestWritesAreForkLocked` takes 26.9 s under `-race` locally, against 0.6–5 s plain. The CI verifygate race leg took 256.8 s on this head, against 238.1 s on the iter-228 PR head and 155.1 s on the merge. That is runner noise plus the extra fixtures. It is not a gate risk at the current budgets, but watch it.
7. **The design doc still reads `Status: PLANNED` in `planned/`.** The controller moves it on landing, as with the parent.

## Score

| Category | Points |
|---|---|
| Tests pass | 20/20 |
| Lint clean (vet, gofmt) | 10/10 |
| Acceptance criteria | 28/30. AC1–AC7, AC9 and AC12 were reproduced first-party or by equivalent. AC8 and AC11 are shown by the CI log below. AC10 is partial |
| Code quality | 12/15. N6 and N10 survive |
| Documentation | 12/15. The executor report is thorough; the doc move is pending |
| Design fidelity | 10/10. D1, D2 and D3 are implemented as specified, and the 6 plan overrides are honoured |
| **Total** | **92/100, PASS** |

## CI evidence (PR head 733dffc)

There are 2 check-runs:
- `ailang-code verify gate`: completed, success.
- `go host build + test gate` (job 111281818460, run 37150052585): completed, success, from 20:01:58Z to 20:16:23Z.

Every step is success, including `Fork-locked write tests, verbose (w-verifygate-etxtbsy linux gate)`.

The verbose step's lines on `ubuntu-latest`. This is the first linux run of the two new tests (AC8/AC11), and the step's loop asserted all five:

```
2026-10-03T20:14:18.0736998Z --- PASS: TestForkLockedWriteBlocksConcurrentFork (0.25s)
2026-10-03T20:14:18.0738695Z --- PASS: TestForkLockedWrappersHoldLockThroughClose (0.00s)
2026-10-03T20:14:18.0740305Z --- PASS: TestKernelRefusesExecOfWriterOpenFile (0.00s)
2026-10-03T20:14:18.0741468Z --- PASS: TestVerifygateTestWritesAreForkLocked (5.31s)
2026-10-03T20:14:18.0742789Z --- PASS: TestNoParallelWriteForkPackagesOutsideVerifygate (0.12s)
```

Verifygate took 195.657 s plain and 256.832 s under race. The two `WARNING: DATA RACE` lines in the log belong to the intentional race-detector known-positive control (`design_docs/verification/w-race-gate-blindspot/racecontrol/main.go`, `exit status 66`), not to this change.

Evidence logs are under `~/.ailang/state/eval229/`: `full.log`, `race.log`, `m_*.log`, `step{A,B,C,D}.log` and `ci_go.log`.
