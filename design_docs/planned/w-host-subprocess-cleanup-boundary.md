# w-host-subprocess-cleanup-boundary — make the output-overflow kill process-group-wide in capsule and broker, surface a failed kill behind the typed overflow error, and bound the one case a group kill cannot reach (a setsid escapee) by closing the parent's read ends; `cmd.WaitDelay` is measured inert here and is not added

- Status: **planned** · Date: 2026-09-26 · Designer: claude-opus-5-5 (iter-192) · Base commit: `87912f7` · Owning queue row: **24** · Scope class: daemon robustness (clause-2) · Verify profile: go-host (no `.ail` change)
- Prototype: scratch worktree `../.design-wt-iter192` at `87912f7`, now removed. Its diff (`proto-final.patch`), shipped-style tests, probe sources, mutation harness (`mutate.py`) and every raw output cited below are banked under `~/.ailang/state/world-iter192/design/`.
- **Why this doc has a process-safety section (§7).** Iterations 190 (twice) and 191 attempted this design and were killed silently, together with the driver and controller. Their probes did `syscall.Kill(gc, SIGKILL)` in `t.Cleanup`, where `gc` came from a helper that returned **-1** when the grandchild's pid file never appeared. `kill(-1, SIGKILL)` signals every process the user owns. The banked `*.HAZARD-kill-minus-1` files were not copied or run. Every probe in this doc was written fresh behind a pid guard, and that guard is a shipped requirement (D4, AC5, X14/X15).

## §1 Problem

This section restates row 24 using this iteration's measurements. Every claim has a row in §12.

- **P1. Two kill paths in each package; only one is group-wide.**
  - The positive enumeration of kill sites in non-test `host/` and `cmd/` gives 7 lines (V1):
    - `broker/handlers.go:82`: `killGroup` = `syscall.Kill(-pgid, SIGKILL)`, used by `Cancel` at `:106`.
    - `broker/handlers.go:122`: `_ = cmd.Process.Kill()`, which is the **overflow** path.
    - `capsule/capsule.go:180`: `Cancel`, which is group-wide.
    - `capsule/capsule.go:66`: `cmdChild.Kill` = `c.cmd.Process.Kill()`, which is the **overflow** kill used by `collectOutput`.
  - Both files set `Setpgid: true` (`:101`, `:175`), and they are the only two files that do (V1, V4).
  - The row's `capsule.go:165/:193` line numbers date from iteration 99 and have drifted.
- **P2. Capsule's defect is lost promptness; broker's is a leaked process.** Measured at HEAD with my own guarded probes (V10, V11). The fixture forks `sleep 30 &` inheriting stdout, records its pid **before** writing, then overflows a 1024-byte cap.

  | Case (HEAD `87912f7`) | Return | Grandchild at return | Typed error |
  |---|---|---|---|
  | capsule, no grandchild (control) | **14.28 ms** | n/a | `*OutputLimitError` |
  | capsule, grandchild, `ExecTimeout` 3 s | **3.0035 s** | dead (killed by the ctx `Cancel` at the deadline) | `*OutputLimitError` |
  | broker, no grandchild (control) | **6.59 ms** | n/a | `*HandlerOutputOverflowError` |
  | broker, grandchild, `execTimeout` 3 s | **5.83 ms** | **ALIVE**, and still alive after a 2 s poll | `*HandlerOutputOverflowError` |

  - Capsule reproduces the row's figures: 3.005 s against an 11.29 ms control.
  - Broker returns promptly for two reasons. Its `io.LimitedReader` stops at limit+1 without needing EOF, and `Wait` on a `StdoutPipe` does not wait for readers. So `runBounded` returns, the ctx is cancelled by `defer cancel()` **after** `Wait` has reaped the child, and the group `Cancel` never fires. The grandchild survives the call (**controller claim (i) confirmed**).
- **P3. The row's unattributed 33 s residue is ESTABLISHED: `archive.Archive`'s `--version` probe waits for a grandchild's whole lifetime** (V12).
  - Take a fixture that forks the grandchild *before* its `--version` branch, so every invocation forks one. `archiveExecutable` then took **30.371 s**, and `Run` took **3.015 s**. The test total was **33.40 s** against the row's 33.38 s.
  - Mechanism: `probeVersion` (`archive.go:470`) uses `bytes.Buffer` stdout. `Wait` therefore awaits os/exec's copier goroutine, which needs EOF. The direct child has already exited, so the 10 s `probeTimeout`'s kill targets a dead process. There is no `Setpgid` and no `WaitDelay`.
  - So the probe is bounded by the **grandchild's** lifetime (30.37 s), **not** by its own 10 s `probeTimeout`. The comment at `archive.go:69-75` predicts this: *"needs cmd.WaitDelay and is deliberately out of scope here"*.
  - The row's 33 s was archive (≈30 s) + Run (≈3 s). **Nothing in the test's cleanup waits.**
- **P4. The same shape exists at `replay.go:327`** (V13). `runPinnedTransition` with a fake interpreter that forks `sleep 3 &` and exits at once took **3.235 s**, against a **0.205 s** no-fork control. It returned `err=nil`, stdout `"ok\n"`. The 60 s `execTimeout` cannot bound this, because it kills an already-exited direct child. This is outside this row's two packages; see §8.
- **P5. The bounded-cleanup residual is real only for a descendant that LEFT the group, and `cmd.WaitDelay` does not bound it in either package** (V14, V15).
  - A `setsid` escapee (`perl POSIX::setsid`, `sleep 8`, pid recorded before any output) holds capsule's pipes:
    - at HEAD, `Run` returned at **8.06 s** against a 3 s `ExecTimeout`, i.e. the escapee's lifetime;
    - with the group-wide prototype and `WaitDelay = 0`, **8.17 s**;
    - with `WaitDelay = 2 s`, **8.11 s**.
  - The broker's timeout path with an escapee: **8.05 s** against a 3 s bound.
  - `WaitDelay` is inert here for a structural reason. Its expiry is delivered only to `Wait`, and both packages call `Wait` only **after** their reads finish (`collectOutput`'s `wg.Wait()` then `child.Wait()`; `runBounded`'s `io.ReadAll` then `cmd.Wait()`).

## §2 Decisions (Q1–Q5), each by measurement

| Q | Decision | Measured reason |
|---|---|---|
| **Q1** Overflow kill group-wide | **Yes, in both packages, through one `killGroup` seam per package.** Capsule gains a `killGroup` var identical to broker's. Both `cmdChild.Kill` and `Cancel` call it. Broker's overflow path calls the existing `killGroup`. | Prototype, grandchild case (V16): capsule **14.83 ms** (HEAD 3.0035 s); broker 10/10 runs **4.96–7.12 ms**, grandchild dead at return in all 10 (HEAD: alive). One 1.139 s broker outlier appeared only when two packages ran in parallel; 10 sequential repeats did not reproduce it (V17). |
| **Q2** Bounded cleanup | **Do NOT add `cmd.WaitDelay`** to capsule or broker; it is inert in both (P5). **Add a deadline pipe close instead**: `time.AfterFunc(execTimeout+pipeCloseGrace, close read ends)`, `defer`-stopped, with `pipeCloseGrace = 1 s`. | Measured with a 500 ms grace (V18): capsule escapee+overflow **3.506 s** (was 8.05 s), still `*OutputLimitError`; broker escapee timeout **3.505 s** (was 8.05 s), still `*HandlerTimeoutError`. In both cases the escapee is **alive** at return, which proves the close, not a kill, freed the call. It is unkillable from here (it left the group), so the honest bound is on the *call*, not on the escapee. For in-group grandchildren the group kill already closes the pipe, and the closer is never reached: the control and grandchild cases are unchanged at 12–18 ms with the closer armed (V18). |
| **Q3** `errors.Join` surfacing | **Join the kill error only, behind the typed error**: `errors.Join(&OutputLimitError{…}, fmt.Errorf("capsule: overflow kill: %w", killErr))`, and the same for broker. **Do not join `Wait`'s error.** After our own SIGKILL it is always `signal: killed`, which is noise. No explicit `Close` exists to surface. | The typed errors keep their identities: `errors.As(err, &*OutputLimitError)`, `errors.As(err, &*HandlerOutputOverflowError)` and `errors.Is(err, ErrHandlerOverflow)` all hold through the join (AC3, X6–X9). The kill error is reachable only rarely. A group SIGKILL whose members have all exited but are **unreaped** returns **nil** on darwin (V19). The overflow kill always precedes `Wait` (the reap), so ESRCH is not expected on that path. It is surfaced for real failures such as EPERM, and tested by injection. Linux zombie-group errno is **UNMEASURED** (darwin only). |
| **Q4** Other launch sites | **OUT, as two named follow-on rows (§8).** (a) `archive.go:470` + `replay.go:327`: `bytes.Buffer` + `CommandContext`, no group, grandchild-lifetime-bounded (P3 **30.37 s** vs a 10 s bound; P4 **3.235 s** vs 0.205 s). There `WaitDelay` **would** work, because `Wait` owns the copier goroutines. (b) `pkgproj/iface.go:171` already has `WaitDelay = 2 s` but no group, so a descendant leaks. `pkgproj/pkgproj.go:247` has no context at all; that is **row 109's** scope, not re-owned here. | The fixes differ in shape from capsule/broker (copier-goroutine `Wait` vs caller-drained pipes). Including them would double the conflict surface (§9). |
| **Q5** Row 25 | **Not claimed.** The new fixtures emit 200 × 33 = 6600 bytes. That is well under the **65536-byte** pipe capacity row 25 measured, so no arm drives a child blocked in `write()` on a full pipe. Row 25 stays open and unchanged. | The fixture size is in `cleanup_test.go` (`overflowLoop`), and row 25's measured capacity is in the charter (V20). |

## §3 Design

### D1 — Capsule: one group-kill seam for both paths

```go
// Kill is the overflow kill. It is group-wide for the same reason Cancel is:
// a forked grandchild inherits the pipes, so killing only the direct child
// leaves the drains blocked until the deadline (queue row 24).
func (c cmdChild) Kill() error { return killGroup(c.cmd.Process.Pid) }

// killGroup SIGKILLs the child's whole process group. Both the overflow kill
// and the ctx Cancel use it; it is a package-level seam so tests can observe
// or fail the kill.
var killGroup = func(pgid int) error {
	return syscall.Kill(-pgid, syscall.SIGKILL)
}
```

`Cancel` becomes `return killGroup(cmd.Process.Pid)`. `collectOutput` records `killErr` inside `killOnce` and joins it behind `*OutputLimitError` (Q3). The `childProcess` interface and `fakeChild` are unchanged. The existing `collectOutput` unit tests pass without edits (V21).

### D2 — Broker: route the overflow kill through the existing `killGroup`

```go
if int64(len(output)) > bounds.maxOutputBytes {
	// Group-wide, like Cancel: a grandchild holding the pipe would otherwise
	// outlive the call, because Wait returns once the direct child is reaped
	// and the deadline's group kill never fires (queue row 24).
	killErr := killGroup(cmd.Process.Pid)
	_ = cmd.Wait()
	var overflow error = &HandlerOutputOverflowError{Limit: bounds.maxOutputBytes}
	if killErr != nil {
		return nil, errors.Join(overflow, fmt.Errorf("broker: overflow kill: %w", killErr))
	}
	return nil, overflow
}
```

`TestHandlerTimeoutKillsTheWholeProcessGroup` wraps `killGroup` and asserts `count == 1`. It exercises no overflow, so the new call site does not change its count (it passes in the full suite, V22).

### D3 — The deadline pipe close (replaces `WaitDelay`)

In both packages, immediately after `cmd.Start()` succeeds:

```go
// pipeCloseGrace bounds the one case the group kill cannot reach: a
// descendant that left the process group (setsid) and still holds the pipes.
// exec.Cmd.WaitDelay cannot do this here: its timer is only consumed by Wait,
// which runs after the reads.
const pipeCloseGrace = time.Second

closer := time.AfterFunc(r.execTimeout+pipeCloseGrace, func() {
	_ = stdoutPipe.Close()
	_ = stderrPipe.Close() // capsule; broker has one merged pipe
})
defer closer.Stop()
```

- The timer is armed after `Start`, and the ctx was created before it. So the close always fires at least `pipeCloseGrace` after the deadline's group kill.
- A drain unblocked by the close returns `os.ErrClosed`. The existing precedence then classifies the result: overflow first, then `ctx.Err() == DeadlineExceeded` → `*TimeoutError` / `*HandlerTimeoutError` (V18).
- `pipeCloseGrace` is a `const`. The tests bound it through `ExecTimeout`/`execTimeout`, not by overriding it.
- **Declared limit:** the escapee itself is **not** killed. It left the group and nothing here tracks it. The call is bounded and the process is not (§6).

### D4 — `host/proctest`: the guarded process-observation helpers (shipped, S1–S4)

This is a new non-production package, imported only by `_test.go` files. `host/proctest` does not exist today (V6).

- `Check(pid)` refuses `pid <= 1`, `os.Getpid()` and `os.Getppid()`.
- `Signal(pid, sig, send)` calls `send` only after `Check` accepts the pid. The send function is injectable, so the guard is tested **without ever signalling** -1.
- `Alive`, `DeadWithin`, `ReadPid` and `ReapOnCleanup` go through `Check`.
  - `ReadPid` **t.Fatalfs** on a missing file ("fixture did not record grandchild pid"). It never returns a sentinel.
  - `ReapOnCleanup` SIGKILLs a checked pid in `t.Cleanup`, so no fixture `sleep` outlives its test.
- Fixtures write the grandchild pid **before** producing output (S3), with a pid-file wait loop bounded at 500 × 10 ms for the setsid case.
- The package doc comment records the iteration-190/191 incident as the reason the package exists.
- **Rule for all future host subprocess tests:** no raw `syscall.Kill`/`unix.Kill` on a pid not taken from a live `cmd.Process`. Everything else goes through `proctest`.

### D5 — What does NOT change

- The typed error types, their `Error()` text and their `Unwrap` chains.
- `collectOutput`'s signature and the `childProcess` interface.
- `readCapped`.
- The `ExecTimeout`/`execTimeout` defaults.
- `archive`, `replay` and `pkgproj` (§8).
- No `.ail` file.

## §4 Files

| File | Change | Prototype size |
|---|---|---|
| `host/capsule/capsule.go` | D1 + D3 + comment update at `:203-210` (the residual is now owned here and closed; name the follow-on rows) | +29/−4 |
| `host/broker/handlers.go` | D2 + D3 | +17/−2 |
| `host/proctest/proctest.go`, `proctest_test.go` (new) | D4 | 83 + 24 |
| `host/capsule/cleanup_test.go` (new) | AC1–AC4 capsule arms | 91 |
| `host/broker/cleanup_test.go` (new) | AC1–AC4 broker arms | 79 |

## §5 Acceptance criteria

- **AC1 — the overflow kill reaches a forked grandchild (capsule + broker).**
  - Fixture: `(sleep 10; : > survived) & echo $! > gc.pid`, pid written before the 6600-byte overflow. `ExecTimeout` is **120 s = 12×** the grandchild's 10 s lifetime.
  - Pass requires all of the following:
    - the error is the typed overflow (`errors.As`);
    - the `survived` marker is **absent**, meaning the grandchild did not run to completion;
    - `proctest.DeadWithin(gc, 1 s)`, a 10× separation from the 10 s lifetime (measured death ≤ 3 µs after return, V16).
  - **No elapsed-time oracle.** Against a direct-child-only kill, capsule `Run` returns only after the grandchild finished, so the marker is present. Broker returns at once but the grandchild is alive.
- **AC2 — a failed overflow kill is surfaced behind the typed error.**
  - Injected `killGroup` = `errors.Join(orig(pgid), injected)`. `orig` still runs, so no process is leaked.
  - Pass: `errors.As(typed overflow)` **and** `errors.Is(err, injected)`, plus `errors.Is(err, ErrHandlerOverflow)` for broker.
- **AC3 — typed identities unchanged.** The existing `TestF6OutputCap*`, `TestOutputCollection*`, `Test*OutputCapWritesFailureRecord` and `assertHandlerFailureRecord(…, ErrHandlerOverflow)` tests pass unedited.
- **AC4 — the pipe close bounds a setsid escapee (capsule overflow path; broker timeout path).**
  - Fixture: `perl -MPOSIX -e 'POSIX::setsid(); …write $$…; exec "sleep","30"'`, with `ExecTimeout` 1 s.
  - Pass: the typed error (`*OutputLimitError` / `*HandlerTimeoutError`) **and the escapee is ALIVE at return**. The call was freed while its pipe holder lived, so the close did it.
  - perl is the portable `setsid`: ubuntu's `/bin/sh` is dash, which refuses `set -m` without a tty (V23).
  - **Non-vacuity:** the test `t.Skip`s without perl **only when `CI` is unset**; under `CI` a missing perl is `t.Fatal`. The prototype `t.Skip`s unconditionally and must be tightened. The first CI run must show AC4 **PASS, not SKIP**, in `-v` output. Perl on `ubuntu-latest` is **UNMEASURED** here.
- **AC5 — the pid guard refuses before sending (S4).** `proctest.Signal` with a recording send func refuses -1, 0, 1, own pid and parent pid, and records **zero** sends. Positive control: pid 424242 is passed through once.
- **AC6 — gates.** `go vet ./...`, `AILANG_BIN=$HOME/.pinned-ailang/ailang go test ./... -count=1` and `./scripts/verify_ail.sh` (no `.ail` change; run for the gate record) are green. The capsule and broker suites must contain no `t.Parallel` (both mutate a package-global `killGroup`; broker's AST census already enforces it, and capsule has 0 today, V9).

## §6 What is NOT fixed (declared residuals)

1. **A setsid escapee is not killed**, only disowned: the pipe is closed and the call returns. Killing it would require tracking descendants (a cgroup, or a subreaper), which is out of scope for a ~1 d row.
2. **archive `--version` and replay** are grandchild-lifetime-bounded (P3, P4). This is follow-on row A in §8.
3. **pkgproj `iface.go`** leaks a descendant after `WaitDelay` releases the call. This is follow-on row B in §8. **`pkgproj.go:247`** has no context; that is row 109.
4. **Row 25** (blocked-in-`write()` arm) is untouched (Q5).
5. **Linux zombie-group kill errno** is unmeasured (Q3). If linux returns ESRCH for an exited-but-unreaped group, AC2's real-world path would add a joined `no such process` to an ordinary overflow. The typed error would still be primary. The sprint must run AC1 on CI and inspect the error text; if ESRCH appears, filter `errors.Is(killErr, syscall.ESRCH)` before joining.

## §7 Risk: test-process safety (the reason D4 exists)

| Hazard | Where it bit | Guard in this design | Proof |
|---|---|---|---|
| `kill(-1, SIGKILL)` from a pid parsed from a never-written file | iter-190 ×2, iter-191: driver, controller and designer died silently | `proctest.ReadPid` fatals on a missing file; `Check` refuses ≤1/self/parent before any send | AC5; X14, X15 killed. During this design, the first capsule setsid probe run hit a missing pid file (cold perl was killed by the 1 s deadline before `setsid`) and **failed safely with `fixture did not record grandchild pid`** (V14 note) |
| group kill with an unvalidated target | n/a | Only `-pgid` with `pgid = cmd.Process.Pid` of a `Setpgid:true` child we started. The X3/X4 mutants drop `Setpgid`, so `kill(-childpid)` names no group and returns ESRCH, which is harmless | X3, X4 run safely |
| a fixture `sleep` outliving the test | — | `ReapOnCleanup` on every recorded pid; group kill in production code | the full suite left no probe process running |

No `kill -9 -1`, `kill 0`, `pkill` or `killall` appears anywhere, and no shell `kill -- -$x`.

## §8 Residuals → follow-on queue rows (proposed; the controller files them)

- **Row A — `w-archive-replay-grandchild-bound`.**
  - Sites: `archive.go:470` and `replay.go:327`.
  - Measured exposure: 30.37 s against a 10 s `probeTimeout` (P3), and 3.235 s vs 0.205 s (P4).
  - Shape: `Setpgid` + group `Cancel` + `WaitDelay`. `WaitDelay` **is** effective there, because `cmd.Run` with `bytes.Buffer` makes `Wait` own the copiers.
  - Note that the archive probe runs on every daemon start (the `archive.go:60-64` comment).
- **Row B — `w-pkgproj-quality-descendant-leak`.**
  - Site: `iface.go:171`.
  - It is already bounded by `WaitDelay = 2 s`, but it has no group, so a descendant survives. Exposure is **UNMEASURED** (the prototype did not probe pkgproj). The row must measure it first.

## §9 Conflict surface

| Row | Overlap | Collision? |
|---|---|---|
| 25 (blocked-child kill coverage) | Same `capsule.go` `collectOutput`; row 25 adds a test arm, and this doc changes `cmdChild.Kill`'s target and `collectOutput`'s return | **Textual only.** Row 25's arm observes the kill via the child's state; a group kill still kills the child. Land either first; the second rebases trivially. |
| 109 (CrossCheck in-process bound) | `pkgproj.go:247` | **None.** This doc does not touch pkgproj. |
| 111 (daemon lock-blocked read status) | `host/daemon`, `host/store` | **None.** |
| 23 (store deadline-free residue owner) | `host/store` ctx threading | **None.** |
| 26 (bounded Z3 report producer) | Will add a bounded process-group producer in `host/evidence` | **Design dependency, not a collision.** Row 26 should reuse D1's shape (group kill on overflow **and** deadline) and D4's `proctest`. |

## §10 Non-goals

Descendant tracking; changing exec bounds or defaults; `archive`/`replay`/`pkgproj`; row 25's arm; any `.ail`, charter or driver change.

## §11 Milestones (≈0.5 d total, each ≤150 LOC)

- **M1 — `host/proctest` + group-wide overflow kill + joined kill error** (prod ≈ +35/−6; tests ≈ 180).
  - Acceptance: AC1, AC2, AC3, AC5, AC6.
  - Mutations X1–X9, X14 and X15 are killed.
  - The `capsule.go:203-210` comment is updated.
- **M2 — the deadline pipe close** (prod ≈ +12; tests ≈ 50).
  - Acceptance: AC4 (including the CI non-skip rule) and AC6.
  - Mutations X10–X13 are killed.
  - The first CI run's `-v` output shows AC4 PASS on linux, and the §6.5 errno is inspected.

## §12 Test plan and mutation matrix: every row was RUN against the prototype

Harness: `~/.ailang/state/world-iter192/design/mutate.py`. It applies one string mutation, asserts the mutation text occurs exactly once, runs `go test <pkg> -run '^(TestOverflowKill|TestPipeClose|TestSignalRefuses)' -count=1 -p 1`, and restores. Output: `mutations.txt`.

| # | Mutation | Source | Killed by | Result |
|---|---|---|---|---|
| X1 | capsule `cmdChild.Kill` → `c.cmd.Process.Kill()` (revert the fix) | fix | AC1 (10.53 s: marker present), AC2 | **KILLED** |
| X2 | broker overflow → `cmd.Process.Kill()` (revert) | fix | AC1 (grandchild alive), AC2 | **KILLED** |
| X3 | capsule drop `Setpgid` | fix | AC1, AC4 | **KILLED** |
| X4 | broker drop `Setpgid` | fix | AC1, AC4 | **KILLED** |
| X5 | capsule `killGroup` targets `pgid` not `-pgid` | diff | AC1 | **KILLED** |
| X6 | capsule discards the kill error (`_ = child.Kill()`) | diff | AC2 | **KILLED** |
| X7 | broker discards the kill error | diff | AC2 | **KILLED** |
| X8 | capsule kill error **replaces** `*OutputLimitError` | diff | AC2 (`errors.As` fails) | **KILLED** |
| X9 | broker kill error **replaces** `*HandlerOutputOverflowError` | diff | AC2 | **KILLED** |
| X10 | capsule pipe closer stopped immediately (no bound) | diff | AC4 (30.41 s, escapee dead at return) | **KILLED** |
| X11 | broker pipe closer stopped immediately | diff | AC4 | **KILLED** |
| X12 | capsule closer closes stdout only | diff | AC4 (stderr drain still held) | **KILLED** |
| X13 | `pipeCloseGrace` = 1 h | diff | AC4 | **KILLED** |
| X14 | guard `pid <= 1` → `pid < -1` (admits -1) | S4 | AC5 (recording send, never a real signal) | **KILLED** |
| X15 | guard drops the own-pid refusal | S4 | AC5 | **KILLED** |

**Tally: 15/15 killed, 0 survived.** Equivalent mutant noted, not run: `Cancel` spelled `syscall.Kill(-pid)` vs `killGroup(pid)` is behaviour-identical, so no test can distinguish them.

## §13 Verification Log

All commands were run at `87912f7` (repo) or on the prototype worktree at the same base. Raw outputs are in `~/.ailang/state/world-iter192/design/`.

| V | Claim | Command | Observed |
|---|---|---|---|
| V1 | 7 kill-site lines; the overflow paths are direct-child | `grep -rnE 'Process\.Kill\(\|syscall\.Kill\(\|killGroup\(\|Setpgid' --include='*.go' host cmd \| grep -v _test.go` | `broker/handlers.go:82` `syscall.Kill(-pgid…)`, `:101` Setpgid, `:106` `killGroup(cmd.Process.Pid)`, `:122` `_ = cmd.Process.Kill()`; `capsule/capsule.go:66` `…c.cmd.Process.Kill()`, `:175` Setpgid, `:180` `syscall.Kill(-cmd.Process.Pid…)` (7 lines) |
| V2 | 6 launch sites | `grep -rnE 'exec\.Command(Context)?\(' --include='*.go' host cmd \| grep -v _test.go` | broker `:93`, archive `:470`, capsule `:169`, pkgproj `:247` (`exec.Command`), iface `:171`, replay `:327` |
| V3 | `WaitDelay` is set only in `pkgproj/iface.go` (positive control in the same call) | `grep -rn 'WaitDelay' --include='*.go' host cmd \| grep -v _test.go` | `iface.go:73,84,99,174` (set at `:174`); `archive.go:74` is a comment only; no capsule or broker hit |
| V4 | Only 2 files set `Setpgid` | `grep -rlE Setpgid --include='*.go' host cmd \| grep -v _test.go` | `host/broker/handlers.go`, `host/capsule/capsule.go` |
| V5 | Toolchain | `grep -n '^go ' go.mod`; `go version` | `go 1.26.6`; `go1.26.6 darwin/arm64` |
| V6 | `host/proctest` is new (control: `host/capsule` exists) | `ls -d host/proctest host/capsule` | `No such file or directory`; `host/capsule` |
| V7 | Row 24 is the owner named in code | `grep -n 'row 24' host/capsule/capsule.go host/pkgproj/iface.go host/archive/archive.go` | `capsule.go:209`, `iface.go:72` |
| V8 | Pinned interpreter for the gates | `$HOME/.pinned-ailang/ailang --version` | `AILANG v0.41.0` |
| V9 | No `t.Parallel` in capsule tests; broker census exists | `grep -c t.Parallel host/capsule/*_test.go`; `sed -n 70,78p host/broker/handlers_parallel_guard_test.go` | `0`, `0`; census `t.Fatalf("…race the t.Cleanup restore of the killGroup package global")` |
| V10 | Capsule HEAD: grandchild 3.0035 s vs control 14.28 ms | `go test ./host/capsule/ -run TestP192 -count=1 -v` (HEAD) → `head-capsule.txt` | `control … run=14.27675ms … isOverflow=true`; `runfork … run=3.003549917s … isOverflow=true isTimeout=false … aliveAtReturn=false` |
| V11 | Broker HEAD: prompt return, leaked grandchild | `go test ./host/broker/ -run TestP192Broker -count=1 -v` (HEAD) → `head-broker-and-setsid.txt` | `broker-control run=6.58675ms`; `broker-runfork run=5.831833ms isOverflow=true … aliveAtReturn=true diedWithin2s=false after=2.000639667s` |
| V12 | F5 residue = the archive `--version` probe | same run, `TestP192EarlyFork` (HEAD) | `earlyfork archive=30.371383375s run=3.015283458s`; `--- PASS: TestP192EarlyFork (33.40s)` |
| V13 | Replay waits for a grandchild | `go test ./host/replay/ -run TestP192ReplayGrandchild -count=1 -v` → `replay-probe.txt` | `fork=false elapsed=205.10125ms out="ok\n" err=<nil>`; `fork=true elapsed=3.234901292s out="ok\n" err=<nil>` |
| V14 | Setsid escapee holds capsule `Run` at HEAD | `TestP192SetsidEscape` (HEAD, `ExecTimeout` 3 s) | `setsid … run=8.061772792s … isOverflow=true`. Note: the first attempt at `ExecTimeout` 1 s failed with `fixture did not record grandchild pid` (a safe failure) |
| V15 | `WaitDelay` is inert (prototype, group kill) | `TestP192WD` (WaitDelay 0 vs 2 s) → `proto-wd.txt` | capsule `setsid run=8.174490167s` (0) vs `run=8.110386417s` (2 s); broker `setsid` returns ≈50–90 ms either way, escapee alive (no pipe holder on the overflow path) |
| V16 | Group kill: capsule 14.83 ms, grandchild dead | `proto-wd.txt`, `proto-grace.txt` | capsule `runfork run=14.825583ms … aliveAtReturn=false diedWithin2s=true after=1.333µs` |
| V17 | Broker 10/10 prompt, no leak; 1.139 s was parallel-load noise | `TestP192BrokerRepeat -p 1` → `proto-broker-repeat.txt` | runs 4.96–7.12 ms, `aliveAtReturn=false` ×10; the outlier `run=1.1389945s` came from the parallel two-package run in `proto-wd.txt` |
| V18 | The pipe close bounds the escapee; typed errors kept | `TestP192Grace`, `TestP192BrokerEscapeNoOverflow` (grace 0 vs 500 ms) → `proto-grace.txt` | capsule `setsid run=8.046036209s` → `run=3.506353791s … isOverflow=true aliveAtReturn=true`; broker `escape-nooverflow run=8.049418083s` → `run=3.505477208s err=…timed out after 3s aliveAtReturn=true`; controls 12–18 ms unchanged |
| V19 | Darwin: group kill of an exited-but-unreaped group returns nil | `TestP192ZombieGroupKillErrno` → `head-capsule.txt` | `zombie-group kill errno=<nil> ; after-reap errno=no such process` |
| V20 | Row 25's pipe capacity | `awk '/^25\. /,/^26\. /' design_docs/world-mission.md` | "`BLOCKED after 65536 bytes`" |
| V21 | Existing `collectOutput` tests pass on the prototype | `go test ./host/capsule/ -run 'Overflow\|PipeClose' -v` | `PASS: TestOutputCollectionOverflowKillsAndOutranksDeadline`, `…TwoOverflowsKillOnce`, plus the 3 new tests |
| V22 | Full gates on the prototype | `go vet ./... && AILANG_BIN=$HOME/.pinned-ailang/ailang go test ./... -count=1` → `full-test.txt` | `VET-OK`; 22 `ok` lines, 0 non-ok lines |
| V23 | Dash refuses `set -m`; perl is present locally | `dash -c 'set -m; sleep 3 & …ps -o pid=,pgid=…'`; `command -v perl` | `dash: 1: set: can't access tty; job control turned off`, child pgid = parent pgid (26684); `/usr/bin/perl` v5.34.1 |
| V24 | Mutation matrix | `python3 mutate.py` → `mutations.txt` | 15 × `KILLED` (§12) |

## §14 Controller claims checked

- **(i)** Broker at HEAD is prompt but leaks: **TRUE** (V11).
- **(ii)** Capsule ≈3.005 s → ≈16 ms: **TRUE** (V10 3.0035 s; V16 14.83 ms).
- **(iii)** Earlyfork `archive=30.3 s` is the F5 residue: **TRUE, and now established** (V12: 30.37 + 3.015 = 33.40 s). A further finding: the probe also **exceeds its own 10 s `probeTimeout`**.
- **(iv)** The inherited setsid readings are untrustworthy: **TRUE**. Re-measured safely (V14, V18).
- **The inherited prototype's `WaitDelay` is FALSE as a bound in these two packages.** Its comment reads "bounds cleanup after the deadline's group kill, for a descendant that left the process group". Measured inert (V15). Replaced by D3.
