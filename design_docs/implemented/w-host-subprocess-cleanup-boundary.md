# w-host-subprocess-cleanup-boundary — make the output-overflow kill process-group-wide in capsule and broker, surface a failed kill behind the typed overflow error, and bound the whole cleanup (pipe drain AND direct-child wait) by one deadline, even when every kill fails; `cmd.WaitDelay` is measured inert here and is not added

- Status: **implemented** (row 24, iter-193: M1 `4a25e39`, M2 `d1583ee`, M3 `e129c63`, CI `719dd76`; evaluator `sonnet` PASS 98/100, zero blocking) · was **planned (revision 1)** · Date: 2026-09-26 · Designer: claude-opus-5-5 (iter-192) · Base commit: `87912f7` (r0), `da5f77b` (r1) · Owning queue row: **24** · Scope class: daemon robustness (clause-2) · Verify profile: go-host (no `.ail` change)
- **Revision 1 (quorum r1), 2026-09-26.** Round 1 BLOCKED 2 reject / 1 pass, with all three reviewers present. Both objections were upheld after measurement (§15):
  - oc-kimi-k3: the AC4 fixture raced cold perl at a 1 s timeout, and perl on CI was unmeasured. **Perl is removed.** The escapee is now the test binary itself in a `setsid` helper mode, at `ExecTimeout` 3 s. It passes `-count=3` and under dash.
  - gpt6-astra: a failed kill leaves `Wait` unbounded. This is TRUE and pre-existing; at r0 it held `Run` for the child's full 30.39 s. **Option A was taken:** a bounded wait (new `host/procbound`), a named `ErrCleanupIncomplete`, a background reaper that owns the eventual `Wait`, and a process-wide admission cap.

  Mutation matrix: **27/27 killed**, including all 15 r0 rows re-anchored to the r1 diff.
- Prototypes:
  - r0: scratch worktree `../.design-wt-iter192`.
  - r1: `../.design-wt-iter192b` at `da5f77b`.

  Both are removed. Diffs (`proto-final.patch`, `proto-r1.patch`), shipped-style tests (`r1/`), probes, the mutation harnesses (`mutate.py`, `mutate-r1.py`) and every raw output cited below are banked under `~/.ailang/state/world-iter192/design/`.
- **Why this doc has a process-safety section (§7).** Iterations 190 (twice) and 191 attempted this design and were killed silently, together with the driver and controller. Their probes did `syscall.Kill(gc, SIGKILL)` in `t.Cleanup`, where `gc` came from a helper that returned **-1** when the grandchild's pid file never appeared. `kill(-1, SIGKILL)` signals every process the user owns. The banked `*.HAZARD-kill-minus-1` files were not copied or run. Every probe in this doc was written fresh behind a pid guard, and that guard is a shipped requirement (D4, AC5, X14/X15).

## §1 Problem

This section restates row 24 using this iteration's measurements. Every claim has a row in §13.

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
  - Broker returns promptly for two reasons. Its `io.LimitedReader` stops at limit+1 without needing EOF, and `Wait` on a `StdoutPipe` does not wait for readers. So `runBounded` returns, the ctx is cancelled by `defer cancel()` **after** `Wait` has reaped the child, and the group `Cancel` never fires. The grandchild survives the call.
- **P3. The row's unattributed 33 s residue is ESTABLISHED: `archive.Archive`'s `--version` probe waits for a grandchild's whole lifetime** (V12).
  - Take a fixture that forks the grandchild *before* its `--version` branch, so every invocation forks one. `archiveExecutable` then took **30.371 s**, and `Run` took **3.015 s**. The test total was **33.40 s** against the row's 33.38 s.
  - Mechanism: `probeVersion` (`archive.go:470`) uses `bytes.Buffer` stdout. `Wait` therefore awaits os/exec's copier goroutine, which needs EOF. The direct child has already exited, so the 10 s `probeTimeout`'s kill targets a dead process. There is no `Setpgid` and no `WaitDelay`.
  - So the probe is bounded by the **grandchild's** lifetime (30.37 s), **not** by its own 10 s `probeTimeout`. The comment at `archive.go:69-75` predicts this.
  - The row's 33 s was archive (≈30 s) + Run (≈3 s). **Nothing in the test's cleanup waits.**
- **P4. The same shape exists at `replay.go:327`** (V13). `runPinnedTransition` with a fake interpreter that forks `sleep 3 &` and exits at once took **3.235 s**, against a **0.205 s** no-fork control. It returned `err=nil`. This is outside this row's two packages; see §8.
- **P5. A descendant that LEFT the group is bounded by nothing today, and `cmd.WaitDelay` does not help in either package** (V14, V15).
  - A `setsid` escapee holding capsule's pipes kept `Run` for **8.06 s** at HEAD against a 3 s `ExecTimeout`, i.e. the escapee's lifetime.
  - With the group-wide prototype: 8.17 s with `WaitDelay = 0`, 8.11 s with `WaitDelay = 2 s`.
  - `WaitDelay` is inert here because its expiry is delivered only to `Wait`, and both packages call `Wait` only **after** their reads finish.
- **P6 (r1, gpt6-astra). If the kills fail, `Wait` is unbounded, and this is PRE-EXISTING.**
  - At HEAD, broker's overflow path is `_ = cmd.Process.Kill(); _ = cmd.Wait()` with no `WaitDelay`. Capsule's `collectOutput` calls `child.Wait()` = `cmd.Wait()` (V1, V3).
  - Measured on the r0 design (group kill + pipe close, unbounded `Wait`). With an injected `killGroup` that returns EPERM **without signalling** (so both the overflow kill and `Cancel` fail) and a direct child that `exec sleep 30`s after overflowing, `Run` held for **30.39 s**, the child's whole lifetime (V28).
  - r0's D3 claimed a cleanup bound it did not deliver in this case. r1 covers it (D3, D6).

## §2 Decisions (Q1–Q5 and the r1 objections), each by measurement

| Q | Decision | Measured reason |
|---|---|---|
| **Q1** Overflow kill group-wide | **Yes, in both packages, through one `killGroup` seam per package.** Capsule gains a `killGroup` var identical to broker's. Both `cmdChild.Kill` and `Cancel` call it. Broker's overflow path calls the existing `killGroup`. | Grandchild case (V16, V17): capsule **14.83 ms** (HEAD 3.0035 s); broker 10/10 runs 4.96–7.12 ms, grandchild dead at return (HEAD: alive). |
| **Q2** Bounded cleanup | **No `cmd.WaitDelay`**; it is inert (P5). **One cleanup deadline per call, `Start + execTimeout + 2·pipeCloseGrace` (`pipeCloseGrace` = 1 s), covering both phases:** (a) at `execTimeout + pipeCloseGrace` the parent's read ends are closed, so the drains return; (b) the direct-child `Wait` is bounded by the same deadline through `procbound.Wait` (D6). If the deadline expires before the child exits, the call returns and a background waiter owns the eventual reap. | (a): escapee cases return at `ExecTimeout`+1 s, **4.01 s** (capsule) / **4.01–4.22 s** (broker) at `ExecTimeout` 3 s, typed errors kept, escapee **alive** at return (V25). (b): with both kills failing (EPERM, nothing signalled), capsule `Run` returns in **3.16–3.26 s** at `ExecTimeout` 1 s, i.e. 1 + 2·1 s, where r0 took **30.39 s** (V26, V28). Broker returns in **3.0 s per path** (overflow, timeout; V26). The direct child is **alive** at return, and after a guarded SIGKILL the background waiter reaps it (outstanding count returns to baseline). |
| **Q3** `errors.Join` surfacing | **The typed error stays primary. Behind it the call joins the kill error (overflow path) and `procbound.ErrCleanupIncomplete` (overflow and timeout paths).** With neither present, the typed error is returned unwrapped (identity unchanged). `Wait`'s own `signal: killed` is not joined, because it is noise after our own SIGKILL. | `errors.As` on `*OutputLimitError` / `*HandlerOutputOverflowError` / `*HandlerTimeoutError`, and `errors.Is(ErrHandlerOverflow)`, hold through the join (AC2, AC7; X8/X9, Y4–Y6). A **negative** assertion in AC1/AC4 (no `ErrCleanupIncomplete` when the child exited) kills Y11/Y12, the too-short-deadline mutants. |
| **Q4** Other launch sites | **OUT, as two named follow-on rows (§8).** archive + replay are grandchild-lifetime-bounded (30.37 s vs a 10 s bound; 3.235 s vs 0.205 s); there `WaitDelay` **would** work. `pkgproj/iface.go` leaks a descendant. `pkgproj.go:247` is **row 109's**. | The fixes differ in shape (copier-goroutine `Wait` vs caller-drained pipes). |
| **Q5** Row 25 | **Not claimed.** The fixtures emit 6600 bytes, far below the **65536-byte** pipe capacity row 25 measured. | V20. |
| **r1-A** (astra) cover or narrow? | **Option A, covered in scope.** It fits: `host/procbound` is 58 lines including comments, capsule is +66/−8 in total (r0+r1), broker is +35/−4 (V29). | P6 measured TRUE and pre-existing (30.39 s); the fix is measured in V26 and 6 mutants (Y1–Y6). |
| **r1-A** why waiters cannot accumulate without bound — **the soft-cap half of this answer is WITHDRAWN at r2 (gpt6-astra); the cap is now an exact reservation, D6/AC10** | (i) A waiter outlives a call **only** when the direct child survives its SIGKILL past the deadline. (ii) SIGKILL to our own same-uid, unreaped child group **succeeded in every probe**, including a SIGSTOPped group (`kill=<nil>`, reaped in 1.06–1.16 ms, V27). EPERM cannot arise for a same-uid target, and ESRCH arises only after the reap. So outside injection a waiter appears only for a process that does not act on SIGKILL, such as one in uninterruptible I/O (**not reproduced, UNMEASURED**). (iii) The cap is kept anyway because it is cheap: `procbound.MaxOutstanding = 8` process-wide, with `Admit()` refusing new subprocesses with `ErrCleanupBacklog` before `Start` (AC8; Y7–Y10). It is a **soft** cap: N concurrent callers can each pass `Admit` before any of them abandons a waiter, so the worst case is 8 + the number of concurrent calls. | V26, V27; `TestAbandonedWaitersAreCountedCappedAndReaped`. |
| **r1-B** (kimi) AC4 fixture | **Perl is removed.** The fixture script launches **the test binary itself** (`proctest.EscapeeShell`) in a helper mode (`proctest.RunEscapeeIfRequested`). The helper calls `syscall.Setsid()` (it is a child of `sh`, so not a group leader), writes its pid **atomically after the escape** (tmp + rename, so the file's existence proves the escape happened), and sleeps 30 s. The script waits, bounded, for the pid file **before** writing any output. `ExecTimeout` is **3 s**. | `-count=3`: 3/3 PASS per package, 4.01–4.22 s (V25). The broker arms also pass with the fixture run under **`/bin/dash`**, ubuntu's `/bin/sh` (V30). No perl remains, so the perl-on-CI premise is **N/A** (V31). |

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

- `Cancel` becomes `return killGroup(cmd.Process.Pid)`.
- `collectOutput` records `killErr` inside `killOnce`.
- `collectOutput` returns `withCleanupFailures(typed, killFailure, runErr)`:
  ```go
  // withCleanupFailures keeps the typed error primary (errors.As-reachable) and
  // joins a failed kill and an incomplete cleanup behind it. With neither it
  // returns primary itself, so the typed error's identity is unchanged.
  ```
  It is applied to both the `*OutputLimitError` and the `*TimeoutError` branch.
- `collectOutput`'s signature, the `childProcess` interface and `fakeChild` are unchanged. The existing `TestOutputCollection*` unit tests pass unedited (V21).

### D2 — Broker: route the overflow kill through the existing `killGroup`

```go
errs := []error{&HandlerOutputOverflowError{Limit: bounds.maxOutputBytes}}
if killErr := killGroup(cmd.Process.Pid); killErr != nil {
	errs = append(errs, fmt.Errorf("broker: overflow kill: %w", killErr))
}
if waitErr := wait(); errors.Is(waitErr, procbound.ErrCleanupIncomplete) {
	errs = append(errs, waitErr)
}
if len(errs) == 1 {
	return nil, errs[0]
}
return nil, errors.Join(errs...)
```

The timeout branch joins `ErrCleanupIncomplete` behind `*HandlerTimeoutError` in the same way. `TestHandlerTimeoutKillsTheWholeProcessGroup` wraps `killGroup` and asserts `count == 1`. It exercises no overflow and passes in the full suite (V22, V32).

### D3 — One cleanup deadline: pipe close, then a bounded wait (replaces `WaitDelay`)

In both packages, around `cmd.Start()`:

```go
// pipeCloseGrace bounds the one case the group kill cannot reach: a
// descendant that left the process group (setsid) and still holds the pipes.
// execTimeout+pipeCloseGrace after Start the parent's read ends are closed, so
// the drains return; Wait is then bounded by a further pipeCloseGrace, so Run
// returns within execTimeout+2*pipeCloseGrace even if every kill fails.
// exec.Cmd.WaitDelay cannot do this here: its timer is only consumed by Wait,
// which runs after the drains.
const pipeCloseGrace = time.Second

release, err := procbound.Admit(); if err != nil { return …, err } // D6 exact reservation (r2); release on Start failure or reap
if err := cmd.Start(); err != nil { … }
started := time.Now()
closer := time.AfterFunc(r.execTimeout+pipeCloseGrace, func() {
	_ = stdoutPipe.Close()
	_ = stderrPipe.Close() // capsule; broker has one merged pipe
})
defer closer.Stop()
// capsule: cmdChild{cmd: cmd, deadline: started.Add(r.execTimeout + 2*pipeCloseGrace)}
//   func (c cmdChild) Wait() error { return procbound.Wait(c.cmd.Wait, time.Until(c.deadline)) }
// broker:  wait := func() error { return procbound.Wait(cmd.Wait, time.Until(cleanupDeadline)) }
```

- **The documented API bound is `execTimeout + 2·pipeCloseGrace`**, plus scheduling. It holds even when every kill fails.
- A drain unblocked by the close returns `os.ErrClosed`. The existing precedence classifies the result: overflow first, then `DeadlineExceeded` → timeout.
- **Declared limits:**
  - a setsid escapee is **not** killed, only disowned;
  - a direct child that survives a failed kill is **not** killed by us. It is **reaped** by the background waiter whenever it exits (§6).

### D4 — `host/proctest`: guarded process-observation helpers + the setsid escapee helper (shipped, S1–S4)

This is a new non-production package, imported only by `_test.go` files (V6).

- `Check(pid)` refuses `pid <= 1`, `os.Getpid()` and `os.Getppid()`.
- `Signal(pid, sig, send)` calls `send` only after `Check` passes, and the send function is injectable. The guard is therefore tested **without ever signalling** -1.
- `Alive`, `DeadWithin`, `ReadPid` (**t.Fatalf** on a missing file, never a sentinel) and `ReapOnCleanup` all go through `Check`.
- **r1:** `EscapeeShell(t, pidfile)` returns the `/bin/sh` fragment `'<os.Executable()>' -test.run='^TestEscapeeHelper$' escapee '<pidfile>' &` followed by a pid-file wait bounded at 500 × 10 ms.
- **r1:** `RunEscapeeIfRequested()` is the helper mode. In a normal run (`flag.Arg(0) != "escapee"`) it returns at once. Otherwise it runs `Setsid`, writes its pid via tmp + rename, and sleeps 30 s. Each package using it declares `func TestEscapeeHelper(t *testing.T) { proctest.RunEscapeeIfRequested() }`.
- The package doc comment records the iteration-190/191 incident.
- Rule: no raw `syscall.Kill` in host subprocess tests on a pid not taken from a live `cmd.Process`. AC7's cleanup kill uses the pgid that production's `killGroup` received, i.e. `cmd.Process.Pid`, and still passes it through `proctest.Signal`.

### D5 — What does NOT change

- The typed error types, their `Error()` text and their `Unwrap` chains.
- `collectOutput`'s signature and `childProcess`.
- `readCapped`.
- Exec bounds and defaults.
- `archive`, `replay` and `pkgproj` (§8).
- No `.ail` file.

### D6 — `host/procbound`: the bounded direct-child wait and its ownership (r1)

```go
var ErrCleanupIncomplete = errors.New("subprocess cleanup incomplete: child not reaped by the cleanup deadline")
var ErrCleanupBacklog = errors.New("subprocess refused: too many unreaped children outstanding")
const MaxOutstanding = 8
func Outstanding() int64                         // admitted children + retained waiters (r2)
func Admit() (release func(), err error)          // r2: nonblocking atomic reservation; ErrCleanupBacklog when full
func Wait(wait func() error, d time.Duration) error
```

- `Wait` runs `wait` in a goroutine.
  - If `wait` finishes within `d`, its result is returned.
  - Otherwise `Wait` increments the outstanding count, leaves a goroutine that receives the eventual result and decrements the count, and returns `ErrCleanupIncomplete`.
- **Who reaps:** that background goroutine owns the child's `cmd.Wait`, which is the reap.
- **What bounds the waiters (r2, gpt6-astra's fix applied verbatim):** "Replace Admit's observational check with a nonblocking atomic reservation acquired before cmd.Start. Count both active children and abandoned waits against a documented process-wide limit. Release the reservation on Start failure or actual completion of cmd.Wait; on cleanup timeout, transfer ownership to the background waiter and retain the reservation until reaping completes. Refuse excess admissions with ErrCleanupBacklog without waiting." The limit is `MaxOutstanding` (8) and is now exact, not soft: the r1 soft-cap justification in §2 r1-A and §6.3 is withdrawn. `release` is idempotent (a `sync.Once`), so a reservation is released exactly once whichever path finishes it.
- **ESRCH filter (r2, oc-glm-5-3 + oc-kimi-k3 fixes applied verbatim), in D1 and D2:** join the kill error only when `!errors.Is(killErr, syscall.ESRCH)` — "ESRCH from kill(-pgid) means the group is already empty, so there was nothing to kill; every other errno is still joined". Comment at both sites: the group leader may already be an unreaped zombie when the overflow kill fires. Applied unconditionally on both platforms, so the shipped error contract no longer depends on the unmeasured linux errno (§6.6).
- This is a package rather than two copies. The logic has one owner, and row 26 (bounded Z3 producer) is the next consumer (§9).

## §4 Files

| File | Change | Prototype size (r1, V29) |
|---|---|---|
| `host/capsule/capsule.go` | D1 + D3 + `withCleanupFailures` + comment update at `:203-210` (the residual is now owned and closed here; name the follow-on rows) | +66/−8 |
| `host/broker/handlers.go` | D2 + D3 | +35/−4 |
| `host/procbound/procbound.go`, `_test.go` (new) | D6 | 58 + 41 |
| `host/proctest/proctest.go`, `_test.go` (new) | D4 | 121 + 24 |
| `host/capsule/cleanup_test.go` (new) | AC1, AC2, AC4, AC7, AC8 | 174 |
| `host/broker/cleanup_test.go` (new) | AC1, AC2, AC4, AC7 (overflow + timeout), AC8 | 169 |

## §5 Acceptance criteria

- **AC1 — the overflow kill reaches a forked grandchild (capsule + broker).**
  - Fixture: `(sleep 10; : > survived) & echo $! > gc.pid`, pid written before the 6600-byte overflow. `ExecTimeout` is 120 s (12× the grandchild's lifetime).
  - Pass requires all of the following:
    - the typed overflow (`errors.As`);
    - **no** `ErrCleanupIncomplete`;
    - the `survived` marker absent;
    - `proctest.DeadWithin(gc, 1 s)`, 10× under the 10 s lifetime.
  - No elapsed-time oracle.
- **AC2 — a failed overflow kill is surfaced behind the typed error.**
  - Injected `killGroup` = `errors.Join(orig(pgid), injected)`. The real kill still runs.
  - Pass: `errors.As(typed)` and `errors.Is(injected)`, plus `errors.Is(ErrHandlerOverflow)` for broker.
- **AC3 — typed identities unchanged.** The existing `TestF6OutputCap*`, `TestOutputCollection*`, `Test*OutputCapWritesFailureRecord` and `assertHandlerFailureRecord(…, ErrHandlerOverflow)` tests pass unedited.
- **AC4 — the pipe close bounds a setsid escapee (capsule overflow path; broker timeout path).**
  - Fixture: `proctest.EscapeeShell`, the test binary in `setsid` helper mode, pid written after the escape and before any output. `ExecTimeout` is **3 s**.
  - Pass: the typed error (`*OutputLimitError` / `*HandlerTimeoutError`), **no** `ErrCleanupIncomplete`, and **the escapee ALIVE at return**.
  - No perl and no skip path.
  - **M2 acceptance:** AC4 passes under `-count=3` locally (measured 3/3 per package, V25), then PASS, not SKIP and not FAIL, in the first linux CI run's `-v` output.
- **AC5 — the pid guard refuses before sending (S4).** A recording send func sees **zero** sends for -1, 0, 1, own pid and parent pid. Positive control: pid 424242 passes through once.
- **AC6 — gates.** `go vet ./...`, `AILANG_BIN=$HOME/.pinned-ailang/ailang go test ./... -count=1` and `./scripts/verify_ail.sh` (gate record) are green. No `t.Parallel` in capsule or broker tests: both mutate package globals (`killGroup`) and read the process-wide `procbound.Outstanding()` (V9).
- **AC7 (r1) — every kill fails; the call is still bounded, and the child is reaped later.**
  - `killGroup` is injected to record its pgid (= `cmd.Process.Pid`) and return **EPERM without signalling**, so both the overflow kill and `Cancel` fail.
  - The direct child `exec sleep 30`s.
  - Arms: capsule overflow; broker overflow; broker timeout (no output). `ExecTimeout` is 1 s, so the bound is 3 s, 10× under the child's 30 s.
  - Pass requires all of the following:
    - the typed error and `errors.Is(ErrCleanupIncomplete)`, plus `errors.Is(EPERM)` on the overflow arms;
    - the direct child **ALIVE** at return (state oracle: the call did not wait for it);
    - `Outstanding() == base+1`.
  - Then the test sends a guarded `proctest.Signal(pgid, SIGKILL)` and requires `Outstanding()` back to `base` within 5 s, i.e. the background waiter reaped it.
  - A `t.Cleanup` kills the child if the test fails earlier.
- **AC8 (r1) — a full backlog refuses admission.**
  - The test fills `procbound` to `MaxOutstanding` with blocking waits, then asserts that capsule `Run` and broker `runBounded` each return `ErrCleanupBacklog` without starting a process.
  - It then releases the waits and requires the count to drain.

- **AC9 (r2, oc-glm-5-3) — an ordinary overflow carries no ESRCH.** "An ordinary overflow with an uninjected killGroup and an already-exited child must return the typed error with no syscall.ESRCH anywhere in the unwrapped chain, while the errors.As/Is contracts of AC3 still hold." Plus oc-kimi-k3's `TestZombieGroupKillNotJoined`: "injecting killGroup = ESRCH-only and asserting the returned error is exactly the typed overflow" (both packages). AC2's injected-EPERM arm is unaffected, so no existing acceptance weakens.
- **AC10 (r2, gpt6-astra) — the reservation is exact under concurrency.** "A barrier-synchronized concurrent test with more callers than the limit and injected failed kills, asserting that admitted children and retained waiters never exceed the limit and that reservations are released exactly once." Fixtures obey S1–S3 (guarded cleanup kills only `cmd.Process`-derived pids).

## §6 What is NOT fixed (declared residuals)

1. **A setsid escapee is not killed**, only disowned: the pipe is closed and the call returns. Killing it would require descendant tracking (a cgroup or a subreaper).
2. **A direct child that survives a failed kill is not killed by us.** It is reaped by the background waiter when it exits. The only non-injected route is a process that ignores SIGKILL for a while (uninterruptible I/O), which is **UNMEASURED**.
3. ~~The admission cap is soft~~ — **withdrawn at r2**: the cap is an exact atomic reservation (D6, AC10).
4. **archive `--version` and replay**: follow-on row A (§8). **pkgproj `iface.go`**: follow-on row B. **`pkgproj.go:247`**: row 109.
5. **Row 25** (blocked-in-`write()` arm) is untouched (Q5).
6. **Linux zombie-group kill errno** is still unmeasured locally (darwin: nil, V19; the rig has no linux container runtime — controller-checked, `docker`/`podman`/`colima`/`limactl`/`orb` all absent). It no longer decides the shipped contract: the ESRCH filter (D6) ships unconditionally with AC9 and X16/X17, and oc-kimi-k3's merge gate below makes the linux reading a precondition of merge, not a post-merge patch. **The reactive "inspect the first CI run and patch" clause of r1 is withdrawn.**

## §7 Risk: test-process safety (the reason D4 exists)

| Hazard | Where it bit | Guard in this design | Proof |
|---|---|---|---|
| `kill(-1, SIGKILL)` from a pid parsed from a never-written file | iter-190 ×2, iter-191: driver, controller and designer died silently | `proctest.ReadPid` fatals on a missing file; `Check` refuses ≤1/self/parent before any send | AC5; X14, X15 killed. The r0 capsule setsid probe hit a missing pid file (cold perl killed by a 1 s deadline) and **failed safely** (V14). That is also the reason r1 removed perl |
| group kill with an unvalidated target | n/a | Only `-pgid` with `pgid = cmd.Process.Pid` of a `Setpgid:true` child. The X3/X4 mutants drop `Setpgid`, so `kill(-childpid)` names no group and returns ESRCH, which is harmless | X3, X4 ran safely |
| AC7's cleanup kill | r1 | The pgid is the one production passed to `killGroup` (`cmd.Process.Pid`), sent via `proctest.Signal`; a `killed` flag prevents a second kill after the reap (pid reuse) | AC7 3/3 under `-count=3` |
| a fixture process outliving the test | — | `ReapOnCleanup` on every recorded pid; the escapee helper exits by itself after 30 s at worst | no leftover processes after the suite |

No `kill -9 -1`, `kill 0`, `pkill` or `killall` appears anywhere, and no shell `kill -- -$x`.

## §8 Residuals → follow-on queue rows (proposed; the controller files them)

- **Row A — `w-archive-replay-grandchild-bound`.**
  - Sites: `archive.go:470` and `replay.go:327`.
  - Measured exposure: 30.37 s against a 10 s `probeTimeout` (P3), and 3.235 s vs 0.205 s (P4).
  - Shape: `Setpgid` + group `Cancel` + `WaitDelay`. `WaitDelay` **is** effective there because `cmd.Run` with `bytes.Buffer` makes `Wait` own the copiers. `procbound` is reusable if a failed-kill bound is wanted there too.
  - The archive probe runs on every daemon start.
- **Row B — `w-pkgproj-quality-descendant-leak`.**
  - Site: `iface.go:171`.
  - It is bounded by `WaitDelay = 2 s` but has no group, so a descendant survives. **UNMEASURED**; the row measures it first.

## §9 Conflict surface

| Row | Overlap | Collision? |
|---|---|---|
| 25 (blocked-child kill coverage) | Same `capsule.go`; this doc changes `cmdChild` (adds a `deadline` field) and `collectOutput`'s returns | **Textual only.** Row 25's arm observes the kill via the child's state, and a group kill still kills the child. Its fixture must not construct `cmdChild` without a deadline (a zero deadline would make `Wait` time out at once). Land either first; the second rebases. |
| 109 (CrossCheck in-process bound) | `pkgproj.go:247` | **None** here. `procbound` is available to it. |
| 111, 23 | `host/daemon`, `host/store` | **None.** |
| 26 (bounded Z3 report producer) | A bounded process-group producer in `host/evidence` | **Design dependency, not a collision.** It should reuse D1's shape, `procbound` and `proctest`. |

## §10 Non-goals

Descendant tracking; changing exec bounds or defaults; an exact (CAS) admission cap; `archive`/`replay`/`pkgproj`; row 25's arm; any `.ail`, charter or driver change.

## §11 Milestones (≈0.75 d total, each ≤150 LOC prod+test delta)

- **M1 — `host/proctest` (guards) + group-wide overflow kill + joined kill error.**
  - Acceptance: AC1, AC2, AC3, AC5, AC6.
  - Mutations X1–X9, X14 and X15 are killed.
- **M2 — the pipe close + the `setsid` escapee helper.**
  - Acceptance: AC4 (`-count=3` locally, then PASS on linux CI) and AC6.
  - Mutations X10–X13, Y11 and Y12 are killed.
- **M3 — `host/procbound` + the bounded wait + admission.**
  - Acceptance: AC7, AC8, AC6.
  - Mutations Y1–Y10 are killed.
  - The `capsule.go:203-210` comment is updated to say the residual is closed here, naming rows A and B.

- **r2 additions (carve-out, verbatim reviewer fixes):** M1 also delivers the ESRCH filter with AC9 + `TestZombieGroupKillNotJoined` (mutations X16/X17); M3's `Admit` becomes the atomic reservation with AC10 (mutations Y13–Y15). **M1 exit gate (oc-kimi-k3, verbatim):** "linux verification — run the r1 test set under linux; bank V34 (AC1 joined error text on linux) and V35 (AC7 return/reap timings); any joined ESRCH residue found there blocks merge, not the sprint." On this rig the linux run is the sprint PR's own ubuntu-latest CI job, run with `-v` and banked before merge; oc-glm-5-3's V30 mirror (host/capsule/cleanup_test.go under /bin/dash) is promoted into AC4's acceptance for both packages.

## §12 Test plan and mutation matrix: every row was RUN against the r1 prototype

Harness: `mutate-r1.py`. It applies one string mutation, asserts the mutation text occurs exactly once, runs `go test <pkgs> -run '^(TestOverflowKill|TestPipeClose|TestSignalRefuses|TestFailedKill|TestFullCleanupBacklog|TestWaitReturns|TestAbandoned)' -count=1 -p 1`, and restores. It records the **named** failing tests. Output: `mutations-r1.txt`.

A build failure is not a kill. Y3's first spelling left `cleanupDeadline` unused and failed to compile. It was re-spelled to compile and re-run: `mutation-Y3-detail.txt`, `mutation-Y3-rerun.txt`.

| # | Mutation | Source | Killed by (named failing tests) |
|---|---|---|---|
| X1 | capsule `cmdChild.Kill` → `c.cmd.Process.Kill()` | fix | AC1, AC2, AC7 |
| X2 | broker overflow → `cmd.Process.Kill()` | fix | AC1, AC2, AC7/overflow |
| X3 | capsule drop `Setpgid` | fix | AC1, AC4 |
| X4 | broker drop `Setpgid` | fix | AC1, AC4 |
| X5 | capsule `killGroup` targets `pgid`, not `-pgid` | diff | AC1 |
| X6 | capsule discards the kill error | diff | AC2, AC7 |
| X7 | broker discards the kill error | diff | AC2, AC7/overflow |
| X8 | capsule kill error **replaces** `*OutputLimitError` | diff | AC1, AC2, AC4, AC7 |
| X9 | broker kill error **replaces** `*HandlerOutputOverflowError` | diff | AC1, AC2, AC7/overflow |
| X10 | capsule pipe closer stopped immediately | diff | AC4, AC7 |
| X11 | broker pipe closer stopped immediately | diff | AC4, AC7/timeout |
| X12 | capsule closer closes stdout only | diff | AC4, AC7 |
| X13 | `pipeCloseGrace` = 1 h | diff | AC4, AC7 |
| X14 | guard `pid <= 1` → `pid < -1` (admits -1) | S4 | AC5 (recording send; never a real signal) |
| X15 | guard drops the own-pid refusal | S4 | AC5 |
| Y1 | `procbound.Wait` timer = 1 h (effectively unbounded) | r1 diff | AC7 capsule, AC7/overflow, AC7/timeout |
| Y2 | capsule `cmdChild.Wait` → `c.cmd.Wait()` (the r0/HEAD shape) | r1 diff | AC7 (30.39 s, V28) |
| Y3 | broker `wait` → unbounded `cmd.Wait()` (compiling spelling) | r1 diff | AC7/overflow, AC7/timeout |
| Y4 | capsule drops `ErrCleanupIncomplete` from the join | r1 diff | AC7 |
| Y5 | broker overflow drops `ErrCleanupIncomplete` | r1 diff | AC7/overflow |
| Y6 | broker timeout drops `ErrCleanupIncomplete` | r1 diff | AC7/timeout |
| Y7 | `procbound` never uncounts a reaped waiter | r1 diff | `TestAbandonedWaiters…`, AC7 |
| Y8 | `Admit` never refuses | r1 diff | `TestAbandonedWaiters…`, AC8 |
| Y9 | capsule skips `Admit` | r1 diff | AC8 |
| Y10 | broker skips `Admit` | r1 diff | AC8 |
| Y11 | capsule wait deadline = `execTimeout` (no grace after the pipe close) | r1 diff | AC4 (the negative `ErrCleanupIncomplete` assertion) |
| Y12 | broker wait deadline = `execTimeout` | r1 diff | AC4 (the same) |

| X16 | capsule drops the ESRCH filter (joins every kill error) | r2 diff | AC9, `TestZombieGroupKillNotJoined` — **UNRUN at design time (r2 carve-out); the planner must execute it** |
| X17 | broker drops the ESRCH filter | r2 diff | AC9, `TestZombieGroupKillNotJoined` — **UNRUN; planner executes** |
| Y13 | `Admit` reverts to the observational check (no reservation) | r2 diff | AC10 — **UNRUN; planner executes** |
| Y14 | reservation not retained on cleanup timeout (released at return) | r2 diff | AC10 — **UNRUN; planner executes** |
| Y15 | `release` not idempotent (double release) | r2 diff | AC10 (released-exactly-once assertion) — **UNRUN; planner executes** |

**Tally: 27/27 killed, 0 survived** on the r1 prototype (r0: 15/15). The five r2 rows (X16, X17, Y13–Y15) are specified by the r2 carve-out and were NOT run at design time; the planner must prototype and execute them before the plan is committed. Equivalent mutant noted, not run: `Cancel` spelled `syscall.Kill(-pid)` vs `killGroup(pid)`.

## §13 Verification Log

All commands were run at `87912f7` (repo), or on the r0/r1 prototypes on that base (r1 on `da5f77b`, which differs only by this doc). Raw outputs are in `~/.ailang/state/world-iter192/design/`.

| V | Claim | Command | Observed |
|---|---|---|---|
| V1 | 7 kill-site lines; the overflow paths are direct-child | `grep -rnE 'Process\.Kill\(\|syscall\.Kill\(\|killGroup\(\|Setpgid' --include='*.go' host cmd \| grep -v _test.go` | `broker/handlers.go:82` `syscall.Kill(-pgid…)`, `:101` Setpgid, `:106` `killGroup(cmd.Process.Pid)`, `:122` `_ = cmd.Process.Kill()`; `capsule/capsule.go:66` `…c.cmd.Process.Kill()`, `:175` Setpgid, `:180` `syscall.Kill(-cmd.Process.Pid…)` |
| V2 | 6 launch sites | `grep -rnE 'exec\.Command(Context)?\(' --include='*.go' host cmd \| grep -v _test.go` | broker `:93`, archive `:470`, capsule `:169`, pkgproj `:247`, iface `:171`, replay `:327` |
| V3 | `WaitDelay` is set only in `pkgproj/iface.go` (positive control in the same call) | `grep -rn 'WaitDelay' --include='*.go' host cmd \| grep -v _test.go` | `iface.go:73,84,99,174`; `archive.go:74` is a comment only; no capsule or broker hit |
| V4 | Only 2 files set `Setpgid` | `grep -rlE Setpgid --include='*.go' host cmd \| grep -v _test.go` | `host/broker/handlers.go`, `host/capsule/capsule.go` |
| V5 | Toolchain | `grep -n '^go ' go.mod`; `go version` | `go 1.26.6`; `go1.26.6 darwin/arm64` |
| V6 | `host/proctest` is new (control: `host/capsule` exists); r1 adds `host/procbound` | `ls -d host/proctest host/capsule` | `No such file or directory`; `host/capsule` |
| V7 | Row 24 is the owner named in code | `grep -n 'row 24' host/capsule/capsule.go host/pkgproj/iface.go host/archive/archive.go` | `capsule.go:209`, `iface.go:72` |
| V8 | Pinned interpreter | `$HOME/.pinned-ailang/ailang --version` | `AILANG v0.41.0` |
| V9 | No `t.Parallel` in capsule tests; broker census exists | `grep -c t.Parallel host/capsule/*_test.go`; `sed -n 70,78p host/broker/handlers_parallel_guard_test.go` | `0`, `0`; census `…race the t.Cleanup restore of the killGroup package global` |
| V10 | Capsule HEAD: grandchild 3.0035 s vs control 14.28 ms | `go test ./host/capsule/ -run TestP192 -count=1 -v` (HEAD) → `head-capsule.txt` | `control … run=14.27675ms`; `runfork … run=3.003549917s … isOverflow=true isTimeout=false … aliveAtReturn=false` |
| V11 | Broker HEAD: prompt, grandchild leaked | `go test ./host/broker/ -run TestP192Broker` (HEAD) → `head-broker-and-setsid.txt` | `broker-control run=6.58675ms`; `broker-runfork run=5.831833ms … aliveAtReturn=true diedWithin2s=false` |
| V12 | F5 residue = the archive `--version` probe | `TestP192EarlyFork` (HEAD) | `earlyfork archive=30.371383375s run=3.015283458s`; `--- PASS: TestP192EarlyFork (33.40s)` |
| V13 | Replay waits for a grandchild | `TestP192ReplayGrandchild` → `replay-probe.txt` | `fork=false elapsed=205.10125ms`; `fork=true elapsed=3.234901292s … err=<nil>` |
| V14 | Setsid escapee holds capsule `Run` at HEAD; perl at 1 s failed safely | `TestP192SetsidEscape` (HEAD) | `setsid … run=8.061772792s`; the 1 s attempt: `fixture did not record grandchild pid` |
| V15 | `WaitDelay` is inert | `TestP192WD` → `proto-wd.txt` | capsule `setsid run=8.174490167s` (0) vs `8.110386417s` (2 s) |
| V16 | Group kill: capsule 14.83 ms | `proto-wd.txt` | `runfork run=14.825583ms … aliveAtReturn=false … after=1.333µs` |
| V17 | Broker 10/10 prompt, no leak | `TestP192BrokerRepeat -p 1` → `proto-broker-repeat.txt` | 4.96–7.12 ms, `aliveAtReturn=false` ×10 |
| V18 | r0 pipe close bounds the escapee | `proto-grace.txt` | capsule 8.046 s → 3.506 s; broker 8.049 s → 3.505 s; escapee alive |
| V19 | Darwin: group kill of an exited-but-unreaped group returns nil | `TestP192ZombieGroupKillErrno` | `zombie-group kill errno=<nil> ; after-reap errno=no such process` |
| V20 | Row 25's pipe capacity | `awk '/^25\. /,/^26\. /' design_docs/world-mission.md` | "`BLOCKED after 65536 bytes`" |
| V21 | Existing `collectOutput` tests unaffected | r1 `-count=3` run | `PASS: TestOutputCollectionOverflowKillsAndOutranksDeadline`, `…TwoOverflowsKillOnce`, `…CallerReleaseUnblocksReadersAndWait` ×3 each |
| V22 | r0 full gates | `go vet ./... && go test ./... -count=1` → `full-test.txt` | `VET-OK`; 22 `ok`, 0 FAIL |
| V23 | Dash refuses `set -m` | `dash -c 'set -m; sleep 3 & …'` | `dash: 1: set: can't access tty; job control turned off` |
| V24 | r0 mutation matrix | `mutate.py` → `mutations.txt` | 15 × KILLED |
| V25 | **r1** AC4 with the test-binary escapee at `ExecTimeout` 3 s, `-count=3` | `AILANG_BIN=… go test ./host/capsule/ ./host/broker/ ./host/procbound/ ./host/proctest/ -run '…' -count=3 -p 1 -v` | capsule `TestPipeCloseBoundsEscapedDescendant (4.01s)` ×3; broker `(4.22s)` ×3; all PASS |
| V26 | **r1** AC7: bounded despite failed kills; reaped later | same run | capsule `TestFailedKillBoundsRunAndChildIsReapedLater` 3.16 s / 3.17 s / 3.26 s PASS; broker `TestFailedKillBoundsCallAndChildIsReapedLater (6.02s)` ×3 (two subtests of ≈3 s) PASS; `TestFullCleanupBacklogRefusesRun` PASS ×3 per package |
| V27 | SIGKILL to our own group succeeds, even when SIGSTOPped | probe `r1-zz_probe_kill_test.go` → `r1-kill-own-child.txt` | `stopped=false kill=<nil> wait=signal: killed reaped-after=1.162333ms`; `stopped=true kill=<nil> … reaped-after=1.056458ms` |
| V28 | **P6 baseline**: an unbounded `Wait` (r0/HEAD shape) + failed kills holds `Run` for the child's lifetime | r1 with `cmdChild.Wait` → `c.cmd.Wait()`, `-run '^TestFailedKill'` → `r1-baseline-unbounded-wait.txt` | `--- FAIL: TestFailedKillBoundsRunAndChildIsReapedLater (30.39s)` |
| V29 | r1 size | `git diff --numstat`; `wc -l` | `handlers.go` 35/4; `capsule.go` 66/8; `procbound.go` 58, `_test` 41; `proctest.go` 121, `_test` 24; capsule `cleanup_test.go` 174; broker 169 |
| V30 | Broker fixtures (incl. the escapee) under ubuntu's shell, dash | broker `cleanup_test.go` with `path: "/bin/dash"`, `-count=1 -v` → `r1-dash-broker.txt` | `PASS: TestPipeCloseBoundsEscapedDescendant (4.01s)`, `PASS: TestFailedKillBoundsCallAndChildIsReapedLater (6.02s)`, AC1/AC2/AC8 PASS |
| V31 | Perl is no longer used (the perl-on-CI premise is N/A) | `grep -c perl` over the r1 test files | 0 (the fixture is `proctest.EscapeeShell`) |
| V32 | r1 full gates | `go vet ./... && AILANG_BIN=… go test ./... -count=1` → `r1-full-test.txt` | `VET-OK`; 23 `ok` (the new `procbound`/`proctest` included), 0 FAIL |
| V33 | r1 mutation matrix | `mutate-r1.py` → `mutations-r1.txt` (+ Y3 rerun) | 27 × KILLED, each with named failing tests |

## §14 Controller claims checked

- **r0 (i)–(iv):** all TRUE; see r0 V10–V14.
- **The inherited prototype's `WaitDelay` is FALSE as a bound** in these packages (V15).
- **r1 controller measurement ("pre-existing unbounded `Wait` on failed kills"): TRUE**, reproduced at 30.39 s (V28). This doc now covers it (option A).
- **r1 controller preference ("remove perl via a test-binary `setsid` helper"): adopted and measured** (V25, V30).

## §15 Quorum log

### Round 1: BLOCKED 2 reject / 1 pass, all present

| Reviewer | Verdict | Surface named | Resolution |
|---|---|---|---|
| oc-glm-5-3 | PASS | — | — |
| oc-kimi-k3 | REJECT | **AC4 fixture** (§5 AC4 / M2): a 1 s `ExecTimeout` with a perl fixture that writes its pid only after a cold boot + `setsid`, the configuration V14 recorded failing; perl on ubuntu-latest unmeasured | **Upheld.** Perl removed. The test-binary `setsid` helper (the controller's preference) writes its pid atomically after the escape and before any output. `ExecTimeout` is 3 s. `-count=3` 3/3 per package (V25), passes under dash (V30); the perl V-row is N/A (V31). M2 acceptance adopts the reviewer's CI clause. |
| gpt6-astra | REJECT | **D2/D3 cleanup bound**: a failed kill leaves `Wait` unbounded; AC2 never tests a failed kill with a live direct child | **Upheld; option A.** Measured TRUE and pre-existing (V28, 30.39 s). One cleanup deadline now covers the drain and the direct-child wait (D3). `procbound` (D6) supplies the bounded wait, `ErrCleanupIncomplete`, background reap ownership and a (soft) admission cap with `ErrCleanupBacklog`. The EPERM-without-signalling tests (AC7, capsule + broker overflow + broker timeout) require a return within the documented bound, the typed error and the injected failure discoverable, the child alive at return, and an eventual guarded kill + reap. AC8 covers admission. Mutants Y1–Y12 were added and killed. |

### Round 2: BLOCKED 3 reject / 0 pass, all present — resolved by the narrow-refinement carve-out (controller)

| Reviewer | Verdict | Surface named | Resolution |
|---|---|---|---|
| oc-glm-5-3 | REJECT | **linux ESRCH premise**: every probe ran on darwin; §6.6 pre-authorised an untested reactive ESRCH filter | Applied verbatim: ESRCH filter in D1/D2 (D6 bullet), AC9, X16/X17, V34 on linux, V30 dash mirror promoted into AC4 |
| oc-kimi-k3 | REJECT | **linux ESRCH premise** (same surface as glm) | Applied verbatim: M1 exit gate (V34/V35 on linux block merge), `TestZombieGroupKillNotJoined`, X16 |
| gpt6-astra | REJECT | **D6 admission cap is soft** — `Admit` reserves nothing, so "8 + concurrent calls" is no finite bound | Applied verbatim: nonblocking atomic reservation before `Start`, retained by the background waiter until reaping, AC10 barrier test, Y13–Y15; soft-cap justification withdrawn |

Carve-out conditions (gate-2): every remaining objection carries a concrete reviewer-authored `proposed_fix`, and none disputes the design DIRECTION (group-wide overflow kill, joined kill error, one cleanup deadline). The fixes above are the reviewers' own text, applied by the controller; no controller-invented resolution. Objection surfaces per round: r1 = AC4 fixture, cleanup bound; r2 = linux ESRCH premise (×2), admission-cap exactness — no surface repeated across rounds, so this is not the SPLIT signal. Controller measurement for r2: no linux container runtime exists on the rig, so the linux reading is deferred to the PR's CI as a merge gate rather than asserted.

