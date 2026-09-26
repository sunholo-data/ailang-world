# w-host-subprocess-cleanup-boundary (row 24) sprint plan

- Sprint: `w-host-subprocess-cleanup-boundary` · queue row **24** · Go host code, no `.ail`
  change · design doc
  [`w-host-subprocess-cleanup-boundary.md`](w-host-subprocess-cleanup-boundary.md), the
  quorum-approved r2 (carve-out). Read its §3 (D1–D6), §5 (AC1–AC10), §7 (test-process
  safety), §11 (milestones) and §12 (mutation matrix) first.
- Base: `17c62b0` (detached worktree `.planner-wt-iter193`). `git diff da5f77b 17c62b0 --
  host/ cmd/` is **empty** (only the design doc changed), so the design's r1/r2 measurements
  carry over to this tree.
- Planned 2026-09-26 by mission-world iteration 193 (sprint-planner `pi:openrouter/moonshotai/kimi-k3` via `mission_pi_run.sh`, verdict `ok`, 436 s, $0.87 metered; attribution corrected by the controller).
- **The plan is measured, not estimated.** The iteration-192 planner built the complete
  design as uncommitted edits in this worktree (`M host/broker/handlers.go`,
  `M host/capsule/capsule.go`, new `host/procbound/`, `host/proctest/`,
  `host/capsule/cleanup_test.go`, `host/broker/cleanup_test.go`) and ran a 46-row mutation
  drill: `~/.ailang/state/world-iter192/plan/mutations-r2.txt` (**46/46 KILLED**, named
  failing tests) with harness `…/plan/mutate-r2.py`. This iteration verified the prototype
  compiles and passes (§8) and wrote this plan; the prototype files were **not** modified.
- Prototype size (measured, `git diff --numstat` + `wc -l`): `handlers.go` +40/−4,
  `capsule.go` +74/−8, `procbound.go` 76, `procbound_test.go` 137, `proctest.go` 121,
  `proctest_test.go` 24, capsule `cleanup_test.go` 231, broker `cleanup_test.go` 343.

## §1 Scope

In: the two overflow/timeout cleanup paths of `host/capsule` and `host/broker` — group-wide
overflow kill through one `killGroup` seam per package, kill error joined behind the typed
error with an unconditional ESRCH filter, one cleanup deadline (`execTimeout +
2·pipeCloseGrace`) covering pipe-drain close AND the bounded direct-child wait via the new
`host/procbound` (exact atomic reservation, `MaxOutstanding = 8`), and the guarded helpers in
`host/proctest`. Out (design §6/§10): `WaitDelay` (measured inert), archive/replay/pkgproj
(follow-on rows A and B), row 25's arm, descendant tracking, any `.ail`/driver change.

## §2 Baseline

| # | Command | Result |
|---|---|---|
| BL1 | `git merge-base --is-ancestor da5f77b HEAD` + `git diff --stat da5f77b HEAD -- host/ cmd/` | ancestor; **empty diff** — design base and this base agree on all code |
| BL2 | `go version`; `$HOME/.pinned-ailang/ailang --version` | `go1.26.6 darwin/arm64`; `AILANG v0.41.0` |
| BL3 | capsule fixture helpers the new test file uses | `archivedFixture`, `archiveExecutable`, `source` pre-exist in `host/capsule/capsule_test.go:34,40,59` |
| BL4 | `grep -c t.Parallel host/capsule/*_test.go host/broker/*_test.go` | 0 hits (required: both packages mutate the `killGroup` global and census `procbound.Outstanding()`) |

## §3 Milestones

The design's three milestones (§11) are kept, with measured re-assignments (DV5). Each
boundary **compiles and passes** cumulatively: M1 → M1+M2 → M1+M2+M3. Staging rule: where a
prototype test or function references `procbound` before M3 lands it, the milestone applies
the prototype text **minus exactly the enumerated lines**, which are applied at M3 as their
own hunks. No code is rewritten; every staged omission is listed below.

### M1 — `host/proctest` (guards) + group-wide overflow kill + joined kill error + ESRCH filter

Production hunks (from the prototype diff):
- **new `host/proctest/proctest.go`, staged**: package doc, `Check`, `Signal`, `Alive`,
  `DeadWithin`, `ReadPid`, `ReapOnCleanup`. **Omit until M2**: `EscapeeShell`,
  `RunEscapeeIfRequested` and the `flag` import. M1 imports: `errors, fmt, os, strconv,
  strings, syscall, testing, time`.
- **new `host/proctest/proctest_test.go`**, verbatim: `TestSignalRefusesDangerousPidsWithoutSending` (AC5).
- `host/capsule/capsule.go`: add `killGroup` var + comment; `cmdChild.Kill` →
  `return killGroup(c.cmd.Process.Pid)` + comment (struct keeps only the `cmd` field; `Wait`
  stays `c.cmd.Wait()` until M3); `Cancel` → `return killGroup(cmd.Process.Pid)`; add
  `var killErr error` and `killOnce.Do(func() { killErr = child.Kill() })`; overflow branch
  gains the `killFailure` computation with the ESRCH-filter comment and
  `withCleanupFailures(&OutputLimitError{Limit: limit}, killFailure, runErr)`; timeout branch
  → `withCleanupFailures(&TimeoutError{Limit: execTimeout}, nil, runErr)`; add
  `withCleanupFailures` **staged**: full 3-arg signature, body **without** the
  `if errors.Is(waitErr, procbound.ErrCleanupIncomplete)` clause (M3 hunk; `waitErr` is a
  legal unused parameter at M1). No `procbound` import at M1.
- `host/broker/handlers.go`: overflow branch replaces
  `_ = cmd.Process.Kill(); _ = cmd.Wait(); return nil, &HandlerOutputOverflowError{…}` with
  the `errs` slice + group kill + ESRCH filter, **staged**: keep `_ = cmd.Wait()` where the
  prototype has the 3-line `if waitErr := wait(); errors.Is(…ErrCleanupIncomplete)` clause
  (that clause and the `wait` closure arrive at M3). No `procbound` import at M1.

Test hunks — new `host/capsule/cleanup_test.go` and `host/broker/cleanup_test.go`, staged:
- capsule: `overflowLoop`, `scriptedInterpreter`, `runScripted`,
  `TestOverflowKillReachesForkedGrandchild` (AC1), `TestOverflowKillErrorJoinedBehindTypedError`
  (AC2), `TestOrdinaryOverflowCarriesNoESRCH` (AC9), `TestZombieGroupKillNotJoined` (AC9).
- broker: `overflowLoop`, `runScript`, same four tests.
- **Staged omissions from AC1 (both packages), applied at M3**: `base :=
  procbound.Outstanding()`; the post-call `Outstanding() != base` block; the
  `errors.Is(err, procbound.ErrCleanupIncomplete)` negative block. **Omit until M2**:
  `TestEscapeeHelper`. **Omit until M3**: `recordFailingKill` and the AC7/AC8/AC10/Start-failure
  tests. Imports are adjusted mechanically to the staged contents (M1 drops `procbound`,
  `sync`, `strings`, `sync/atomic`).
- AC3: the existing `TestF6OutputCap*`, `TestOutputCollection*`,
  `Test*OutputCapWritesFailureRecord` and `assertHandlerFailureRecord(…, ErrHandlerOverflow)`
  pass unedited (verified on the full prototype, §8).

Must pass at the boundary: the four staged M1 tests per package plus both packages' whole
existing suites (`go test ./host/proctest/ ./host/capsule/ ./host/broker/ -count=1 -p 1
-timeout 600s` green; `go vet ./host/...` rc=0).
Mutation rows killed (§5.1): **X1–X9, X14, X15, X16, X17, Z6, Z7, Z8, Z9, Z11, Z12**.
(X3/X4/X8's AC4-named killers are redundant: each is also killed by an M1 test.)

### M2 — the deadline pipe close + the setsid escapee helper

Production hunks:
- `host/proctest/proctest.go`: add the `flag` import, `EscapeeShell`, `RunEscapeeIfRequested`
  (verbatim).
- `host/capsule/capsule.go`: `pipeCloseGrace` const + comment; in `Run`, immediately after
  `cmd.Start()`'s error branch: the `closer := time.AfterFunc(r.execTimeout+pipeCloseGrace,
  …)` block closing **both** pipes, and `defer closer.Stop()`. **Do not** add
  `started := time.Now()` yet (unused at M2 → compile error; it is M3's hunk).
- `host/broker/handlers.go`: `pipeCloseGrace` const + comment; after Start:
  `closer := time.AfterFunc(bounds.execTimeout+pipeCloseGrace, func() { _ = pipe.Close() })`
  + `defer closer.Stop()`.

Test hunks: `TestEscapeeHelper` (both packages, verbatim);
`TestPipeCloseBoundsEscapedDescendant` (both packages, AC4) **staged** — omit the
`errors.Is(err, procbound.ErrCleanupIncomplete)` negative block (M3 hunk; it is the
assertion that kills Y11/Y12 once the bounded wait exists).

Must pass at the boundary: M1 set plus `AILANG_BIN=… go test ./host/capsule/ ./host/broker/
-run '^TestPipeCloseBoundsEscapedDescendant$' -count=3 -p 1 -timeout 300s -v` — 3/3 PASS per
package (prototype measured 4.01–4.31 s per run, §8); full four-package suite green.
Mutation rows killed (§5.2): **X10, X11, X12, X13, X13b**.
AC4 acceptance: `-count=3` locally, then PASS (not SKIP, not FAIL) in the first linux CI `-v`
log (§7).

### M3 — `host/procbound` + bounded direct-child wait + admission + comment

Production hunks:
- **new `host/procbound/procbound.go` and `procbound_test.go`, verbatim** (Wait with
  background-reaper ownership, Admit as a nonblocking CAS reservation, idempotent
  `release`, `ErrCleanupIncomplete`, `ErrCleanupBacklog`, `MaxOutstanding = 8`; tests:
  `TestWaitReturnsResultWithinBoundAndReleases`,
  `TestAbandonedWaitersRetainReservationsUntilReaped`, `TestReleaseIsIdempotent`,
  `TestAdmitIsExactUnderConcurrency`).
- `host/capsule/capsule.go`: `procbound` import; `cmdChild` gains `deadline`/`release`
  fields; `Wait` → `procbound.Wait(c.cmd.Wait, time.Until(c.deadline), c.release)`; `Run`
  gains the `Admit` block, `release()` on Start failure, `started := time.Now()`, and the
  full `cmdChild{cmd: cmd, deadline: started.Add(r.execTimeout + 2*pipeCloseGrace), release:
  release}` literal; `withCleanupFailures` gains the `ErrCleanupIncomplete` clause.
- `host/broker/handlers.go`: `procbound` import; `Admit` block + `release()` on Start
  failure; `cleanupDeadline` + the `wait` closure; overflow branch: the staged
  `_ = cmd.Wait()` is replaced by the prototype's 3-line `waitErr` clause; timeout branch:
  `waitErr := wait()` and the `ErrCleanupIncomplete` join.
- **The one authored (non-prototype) hunk** — the design's §4/§11 `collectOutput` comment
  update, which the prototype omits (DV4). Replace the final sentence "*The residual
  lifecycle work is owned by queue row 24, w-host-subprocess-cleanup-boundary.*" with:
  "*Queue row 24 (w-host-subprocess-cleanup-boundary) closed that residual: the overflow kill
  is group-wide, a failed kill and an incomplete cleanup are joined behind the typed error
  (withCleanupFailures), and one cleanup deadline (execTimeout + 2·pipeCloseGrace) bounds the
  drains (the pipe closer in Run) and the direct-child wait (procbound). Declared residuals:
  a setsid escapee is disowned, not killed; a child that survives a failed kill is reaped
  later by procbound's background waiter. The same grandchild shapes in archive and replay
  are follow-on row A (w-archive-replay-grandchild-bound); the pkgproj descendant leak is
  follow-on row B (w-pkgproj-quality-descendant-leak).*"

Test hunks: `recordFailingKill` (both packages); AC1 gains its three deferred procbound
blocks (both packages); AC4 gains its deferred `ErrCleanupIncomplete` negative block (both);
`TestFailedKillBoundsRunAndChildIsReapedLater` (capsule, AC7);
`TestFailedKillBoundsCallAndChildIsReapedLater` (broker, AC7, `/overflow` + `/timeout`);
`TestFullCleanupBacklogRefusesRun` (both, AC8); `TestStartFailureReleasesReservation` (both);
`TestReservationIsExactUnderConcurrency` (broker, AC10); imports completed (`procbound`,
`sync`, `strings`, `sync/atomic`). At M3 the test files are **byte-identical to the
prototype**.

Must pass at the boundary: the full four-package suite green; AC7's return/reap cycle
(capsule ≈3 s, broker ≈6 s at `ExecTimeout` 1 s); AC8; AC10.
Mutation rows killed (§5.3): **Y1–Y15, Y13b, Z1, Z2, Z3, Z4, Z5, Z10**.

## §4 Deviations from the design

- **DV1 — `procbound.Wait` takes a third argument** `release func()` (D6 sketched
  `Wait(wait, d)`). Forced by the r2 exact reservation: ownership of the reservation must
  travel into `Wait` so the background reaper releases it at the reap (mutations Y14, Z3,
  Z4, Z5 target exactly this).
- **DV2 — the mutation table is 46 rows, not §12's 27 (+5 UNRUN).** X13b mirrors X13 for
  broker (the design's X13 was capsule-only); Y13b is the check-then-add (non-CAS) variant of
  Y13; Z1–Z12 are planner-derived diff-covering rows (Start-failure release, Wait success
  release, no-op release, ESRCH-filter misspellings EPERM/`==`/filter-all, Admit off-by-one).
  All 46 were **executed and KILLED** (`mutations-r2.txt`) — including X16, X17, Y13, Y14,
  Y15, which §12 marked "UNRUN; planner executes". The carve-out's planner-execution
  condition is therefore met and recorded here.
- **DV3 — renamed/added tests vs the r1 names in §12.**
  `TestAbandonedWaitersAreCountedCappedAndReaped` →
  `TestAbandonedWaitersRetainReservationsUntilReaped` (r2 retention semantics); new:
  `TestWaitReturnsResultWithinBoundAndReleases`, `TestReleaseIsIdempotent`,
  `TestAdmitIsExactUnderConcurrency` (procbound units for the r2 reservation),
  `TestStartFailureReleasesReservation` ×2 (kills Z1/Z2), `TestReservationIsExactUnderConcurrency`
  (broker AC10), `TestOrdinaryOverflowCarriesNoESRCH` and `TestZombieGroupKillNotJoined`
  (AC9, both named in the design's r2 text), `TestEscapeeHelper` ×2 (D4's declared hook).
- **DV4 — the prototype omits the §4/§11 `collectOutput` comment update.** The plan supplies
  the text (M3 above); it is the executor's only authored hunk.
- **DV5 — milestone re-assignments.** Y11/Y12 move M2→M3: their target expression
  (`…execTimeout + 2*pipeCloseGrace` as the *wait* deadline) and their killing assertion
  (AC4's negative `ErrCleanupIncomplete` check) only exist once the bounded wait is wired at
  M3; at the M2 boundary the mutation has no anchor and the harness would report NOT-APPLIED.
  Z6–Z9 and Z11/Z12 are assigned to M1: they mutate M1's ESRCH filter and join, and each is
  killed at M1 by `TestZombieGroupKillNotJoined` or AC2.
- **DV6 — staged application (not a code change).** Intermediate milestone trees differ from
  the prototype only by the enumerated deferred lines (§3); at M3 the trees are byte-identical.
  This is the reasoned alternative to an M2+M3 merge: the pipe close compiles and is
  AC4-verifiable without `procbound`, and merging would put ~600 LOC into one milestone, far
  past the design's ≤150-LOC milestone guideline.

## §5 Mutation tables (copied verbatim from `mutations-r2.txt`; every row EXECUTED, 46/46 KILLED)

Harness: `python3 ~/.ailang/state/world-iter192/plan/mutate-r2.py <ids…>`, run from the
worktree root. Each edit must match exactly once (`NOT-APPLIED` = wrong milestone or drifted
anchor, not a pass); BUILD-FAIL is not a kill; restores originals in `finally` (safe for
uncommitted work). The executor re-runs each milestone's rows after applying its hunks.

### §5.1 M1 rows

| # | Mutation | Verdict — named killing tests |
|---|---|---|
| X1 | capsule overflow kill -> direct child | KILLED: TestFailedKillBoundsRunAndChildIsReapedLater†, TestOverflowKillErrorJoinedBehindTypedError, TestOverflowKillReachesForkedGrandchild |
| X2 | broker overflow kill -> direct child | KILLED: TestFailedKillBoundsCallAndChildIsReapedLater†(+†/overflow), TestOverflowKillErrorJoinedBehindTypedError, TestOverflowKillReachesForkedGrandchild |
| X3 | capsule drop Setpgid | KILLED: TestOverflowKillErrorJoinedBehindTypedError, TestOverflowKillReachesForkedGrandchild, TestPipeCloseBoundsEscapedDescendant‡ |
| X4 | broker drop Setpgid | KILLED: TestOverflowKillErrorJoinedBehindTypedError, TestOverflowKillReachesForkedGrandchild, TestPipeCloseBoundsEscapedDescendant‡ |
| X5 | capsule killGroup targets pid not group | KILLED: TestOverflowKillReachesForkedGrandchild |
| X6 | capsule discards the kill error | KILLED: TestFailedKillBoundsRunAndChildIsReapedLater†, TestOverflowKillErrorJoinedBehindTypedError |
| X7 | broker discards the kill error | KILLED: TestFailedKillBoundsCallAndChildIsReapedLater†(+†/overflow), TestOverflowKillErrorJoinedBehindTypedError |
| X8 | capsule kill error REPLACES typed overflow | KILLED: TestFailedKillBoundsRunAndChildIsReapedLater†, TestOrdinaryOverflowCarriesNoESRCH, TestOverflowKillErrorJoinedBehindTypedError, TestOverflowKillReachesForkedGrandchild, TestPipeCloseBoundsEscapedDescendant‡, TestZombieGroupKillNotJoined |
| X9 | broker kill error REPLACES typed overflow | KILLED: TestFailedKillBoundsCallAndChildIsReapedLater†(+†/overflow), TestOrdinaryOverflowCarriesNoESRCH, TestOverflowKillErrorJoinedBehindTypedError, TestOverflowKillReachesForkedGrandchild, TestZombieGroupKillNotJoined |
| X14 | guard admits -1 (pid <= 1 -> pid < -1) | KILLED: TestSignalRefusesDangerousPidsWithoutSending |
| X15 | guard drops own-pid refusal | KILLED: TestSignalRefusesDangerousPidsWithoutSending |
| X16 | capsule drops the ESRCH filter | KILLED: TestZombieGroupKillNotJoined |
| X17 | broker drops the ESRCH filter | KILLED: TestZombieGroupKillNotJoined |
| Z6 | capsule filters EPERM instead of ESRCH | KILLED: TestFailedKillBoundsRunAndChildIsReapedLater†, TestZombieGroupKillNotJoined |
| Z7 | broker filters EPERM instead of ESRCH | KILLED: TestFailedKillBoundsCallAndChildIsReapedLater†(+†/overflow), TestZombieGroupKillNotJoined |
| Z8 | capsule ESRCH filter by == (misses wrapped) | KILLED: TestZombieGroupKillNotJoined |
| Z9 | broker ESRCH filter by == (misses wrapped) | KILLED: TestZombieGroupKillNotJoined |
| Z11 | capsule filters every kill error | KILLED: TestFailedKillBoundsRunAndChildIsReapedLater†, TestOverflowKillErrorJoinedBehindTypedError |
| Z12 | broker filters every kill error | KILLED: TestFailedKillBoundsCallAndChildIsReapedLater†(+†/overflow), TestOverflowKillErrorJoinedBehindTypedError |

### §5.2 M2 rows

| # | Mutation | Verdict — named killing tests |
|---|---|---|
| X10 | capsule pipe closer stopped immediately | KILLED: TestFailedKillBoundsRunAndChildIsReapedLater†, TestPipeCloseBoundsEscapedDescendant |
| X11 | broker pipe closer stopped immediately | KILLED: TestFailedKillBoundsCallAndChildIsReapedLater†(+†/timeout), TestPipeCloseBoundsEscapedDescendant, TestReservationIsExactUnderConcurrency† |
| X12 | capsule closer closes stdout only | KILLED: TestFailedKillBoundsRunAndChildIsReapedLater†, TestPipeCloseBoundsEscapedDescendant |
| X13 | capsule pipeCloseGrace = 1h | KILLED: TestFailedKillBoundsRunAndChildIsReapedLater†, TestPipeCloseBoundsEscapedDescendant |
| X13b | broker pipeCloseGrace = 1h | KILLED: TestFailedKillBoundsCallAndChildIsReapedLater†(+†/overflow +†/timeout), TestPipeCloseBoundsEscapedDescendant, TestReservationIsExactUnderConcurrency† |

### §5.3 M3 rows

| # | Mutation | Verdict — named killing tests |
|---|---|---|
| Y1 | procbound.Wait timer = 1h | KILLED: TestFailedKillBoundsCallAndChildIsReapedLater (+/overflow +/timeout), TestFailedKillBoundsRunAndChildIsReapedLater, TestReservationIsExactUnderConcurrency |
| Y2 | capsule cmdChild.Wait unbounded (HEAD shape) | KILLED: TestFailedKillBoundsRunAndChildIsReapedLater |
| Y3 | broker wait unbounded (compiling spelling) | KILLED: TestFailedKillBoundsCallAndChildIsReapedLater (+/overflow +/timeout), TestReservationIsExactUnderConcurrency |
| Y4 | capsule drops ErrCleanupIncomplete from the join | KILLED: TestFailedKillBoundsRunAndChildIsReapedLater |
| Y5 | broker overflow drops ErrCleanupIncomplete | KILLED: TestFailedKillBoundsCallAndChildIsReapedLater (+/overflow) |
| Y6 | broker timeout drops ErrCleanupIncomplete | KILLED: TestFailedKillBoundsCallAndChildIsReapedLater (+/timeout), TestReservationIsExactUnderConcurrency |
| Y7 | background waiter never releases after the reap | KILLED: TestAbandonedWaitersRetainReservationsUntilReaped, TestAdmitIsExactUnderConcurrency, TestFailedKillBoundsCallAndChildIsReapedLater (+/overflow +/timeout), TestFailedKillBoundsRunAndChildIsReapedLater, TestReleaseIsIdempotent, TestReservationIsExactUnderConcurrency |
| Y8 | Admit never refuses | KILLED: TestAbandonedWaitersRetainReservationsUntilReaped, TestAdmitIsExactUnderConcurrency, TestFullCleanupBacklogRefusesRun, TestReleaseIsIdempotent, TestReservationIsExactUnderConcurrency |
| Y9 | capsule skips Admit | KILLED: TestFailedKillBoundsRunAndChildIsReapedLater, TestFullCleanupBacklogRefusesRun |
| Y10 | broker skips Admit | KILLED: TestFullCleanupBacklogRefusesRun, TestReservationIsExactUnderConcurrency |
| Y11 | capsule wait deadline = execTimeout | KILLED: TestPipeCloseBoundsEscapedDescendant (the M3-added negative ErrCleanupIncomplete assertion) |
| Y12 | broker wait deadline = execTimeout | KILLED: TestPipeCloseBoundsEscapedDescendant (same) |
| Y13 | Admit reverts to the r1 observational check | KILLED: TestAdmitIsExactUnderConcurrency, TestReleaseIsIdempotent |
| Y13b | Admit check-then-add (non-atomic reservation) | KILLED: TestAdmitIsExactUnderConcurrency |
| Y14 | reservation released at return, not retained by the waiter | KILLED: TestAbandonedWaitersRetainReservationsUntilReaped, TestFailedKillBoundsRunAndChildIsReapedLater, TestReservationIsExactUnderConcurrency |
| Y15 | release not idempotent | KILLED: TestReleaseIsIdempotent |
| Z1 | capsule drops release() on Start failure | KILLED: TestStartFailureReleasesReservation |
| Z2 | broker drops release() on Start failure | KILLED: TestReservationIsExactUnderConcurrency, TestStartFailureReleasesReservation |
| Z3 | Wait success path does not release | KILLED: TestAbandonedWaitersRetainReservationsUntilReaped, TestAdmitIsExactUnderConcurrency, TestOverflowKillReachesForkedGrandchild (the M3-added Outstanding assertion), TestReleaseIsIdempotent, TestReservationIsExactUnderConcurrency, TestWaitReturnsResultWithinBoundAndReleases |
| Z4 | capsule passes a no-op release to Wait | KILLED: TestFailedKillBoundsRunAndChildIsReapedLater, TestOverflowKillReachesForkedGrandchild |
| Z5 | broker passes a no-op release to Wait | KILLED: TestFailedKillBoundsCallAndChildIsReapedLater (+/overflow +/timeout), TestOverflowKillReachesForkedGrandchild, TestReservationIsExactUnderConcurrency |
| Z10 | Admit off-by-one (n > MaxOutstanding) | KILLED: TestAbandonedWaitersRetainReservationsUntilReaped, TestAdmitIsExactUnderConcurrency, TestFullCleanupBacklogRefusesRun, TestReleaseIsIdempotent, TestReservationIsExactUnderConcurrency |

† = a killer that exists only from M3; the row still has at least one M1/M2-local killer, as
read from `mutations-r2.txt`. ‡ = an M2 test; redundant for the kill (an M1 test also fails).

**Tally: 46/46 KILLED, 0 survived, 0 build-fail** (full prototype). Equivalent mutant noted
by the design, not run: `Cancel` spelled `syscall.Kill(-pid)` vs `killGroup(pid)`.

## §6 Test discipline and process safety (mandatory)

Iterations 190 (twice) and 191 died — taking the driver and controller with them — because a
probe passed pid **-1**, parsed from a pid file that was never written, to `kill(2)`:
`kill(-1, SIGKILL)` signals every process the user owns. The following are hard rules for
every test in this sprint and for any executor-authored probe:

- Every test-side signal goes through `proctest.Signal(pid, sig, send)` or
  `proctest.ReapOnCleanup`. `proctest.Check` refuses `pid <= 1`, the test process and its
  parent **before** `send` is invoked; `Signal`'s send func is injectable precisely so the
  guard is tested (AC5) without ever signalling.
- A signalled pid comes only from (a) a live `cmd.Process.Pid` (in AC7/AC10: the pgid the
  production code passed to the injected `killGroup`, recorded by `recordFailingKill`), or
  (b) a pid file read by `proctest.ReadPid`, which `t.Fatalf`s on a missing or malformed
  file — never a sentinel.
- **No raw `syscall.Kill(...)` call in any `_test.go`.** Passing `syscall.Kill` as the
  `send` argument to `proctest.Signal` (the prototype's only use) is the sanctioned path —
  it fires only after `Check`. (Pre-existing out-of-scope exception:
  `host/verifygate/mission_config_gate_test.go:55` kills `-cmd.Process.Pid` from a live
  process; untouched.)
- No `kill 0`, `kill -9 -1`, `pkill`, `killall` anywhere (grep-clean, §9 PV6).
- No `t.Parallel` in capsule or broker tests: both mutate the package-global `killGroup`
  seam and read the process-wide `procbound.Outstanding()` census; the harness and all gates
  run with `-p 1`.
- The AC7/AC10 cleanup kills are idempotent-guarded (a `killed` flag) so a reaped pid is
  never re-signalled (pid reuse). The escapee helper writes its pid **atomically after**
  `setsid` (tmp + rename), so the pid file's existence proves the escape happened.
- No wall-clock oracles: every assertion is a state oracle (`Alive`/`DeadWithin` on the
  fixture's own pid), an error-identity check, or an `Outstanding()` count.

## §7 Merge gate (design §11 r2, oc-kimi-k3 verbatim clause)

The sprint PR's **ubuntu-latest CI job must run the cleanup tests with `-v`**
(`go test ./host/procbound/ ./host/proctest/ ./host/capsule/ ./host/broker/ -count=1 -p 1 -v`),
and the controller banks the following from the CI log **before merge**:

- **V34** (AC1 joined-error text on linux, no ESRCH residue) — grep for, all PASS:
  `--- PASS: TestOverflowKillReachesForkedGrandchild` (capsule **and** broker),
  `--- PASS: TestOverflowKillErrorJoinedBehindTypedError` (both),
  `--- PASS: TestOrdinaryOverflowCarriesNoESRCH` (both),
  `--- PASS: TestZombieGroupKillNotJoined` (both).
  Any joined ESRCH residue found on linux **blocks merge, not the sprint**.
- **V35** (AC7 return/reap timings) — grep for, with their `-v` elapsed times:
  `--- PASS: TestFailedKillBoundsRunAndChildIsReapedLater` (capsule; expect ≈3 s =
  `ExecTimeout` 1 s + 2·grace) and `--- PASS: TestFailedKillBoundsCallAndChildIsReapedLater`
  (broker; expect ≈6 s for the two arms).
- **AC4's promoted CI clause** (r2): `--- PASS: TestPipeCloseBoundsEscapedDescendant` in both
  packages — PASS, **not SKIP, not FAIL** (the V30 dash mirror is promoted into AC4).

## §8 Prototype gate results (this iteration, on the unmodified prototype)

| Gate | Result |
|---|---|
| `go vet ./host/...` | rc=0, no output |
| `AILANG_BIN=$HOME/.pinned-ailang/ailang go test ./host/procbound/ ./host/proctest/ ./host/capsule/ ./host/broker/ -count=1 -p 1 -timeout 600s` | rc=0: `ok …/host/procbound 0.263s`, `ok …/host/proctest 0.243s`, `ok …/host/capsule 15.557s`, `ok …/host/broker 55.487s` — no FAIL; nothing sandbox-suspect |
| AC4 stability: `go test ./host/capsule/ ./host/broker/ -run '^TestPipeCloseBoundsEscapedDescendant$' -count=3 -p 1 -v` | rc=0, 3/3 PASS per package: capsule 4.22/4.25/4.31 s, broker 4.01/4.01/4.01 s |
| `python3 …/plan/mutate-r2.py` (iter-192 planner) | 46 rows, all KILLED (`mutations-r2.txt`) |

## §9 Execution notes for the executor

- The executor receives this plan plus the prototype patch. **Milestones are delivered by
  applying the prototype hunks, not rewriting them**; the only authored text is the M3
  `collectOutput` comment (§3 M3, DV4) and the mechanical import adjustments of the staged
  test files (§3 M1). Apply per §3, in order M1 → M2 → M3; stage files by name, never
  `git add -A`. Suggested commits: M1 = proctest (staged) + both packages' M1 hunks + staged
  cleanup tests; M2 = escapee helpers + closers + AC4 tests; M3 = procbound + wiring +
  deferred test hunks + comment.
- `AILANG_BIN=$HOME/.pinned-ailang/ailang` is mandatory for `go test` (`host/verifygate`
  fails loudly without it) and for `verify_ail.sh`.
- Gate commands per milestone (all must be green, in order):

  **M1:**
  ```
  go vet ./host/...
  AILANG_BIN=$HOME/.pinned-ailang/ailang go test ./host/proctest/ ./host/capsule/ ./host/broker/ -run '^(TestOverflowKill|TestSignalRefuses|TestOrdinaryOverflow|TestZombieGroup)' -count=1 -p 1 -timeout 300s
  AILANG_BIN=$HOME/.pinned-ailang/ailang go test ./host/proctest/ ./host/capsule/ ./host/broker/ -count=1 -p 1 -timeout 600s
  python3 ~/.ailang/state/world-iter192/plan/mutate-r2.py X1 X2 X3 X4 X5 X6 X7 X8 X9 X14 X15 X16 X17 Z6 Z7 Z8 Z9 Z11 Z12   # 20/20 KILLED
  ```

  **M2 (applied on top of M1):**
  ```
  go vet ./host/...
  AILANG_BIN=$HOME/.pinned-ailang/ailang go test ./host/capsule/ ./host/broker/ -run '^TestPipeCloseBoundsEscapedDescendant$' -count=3 -p 1 -timeout 300s -v   # 3/3 PASS per package
  AILANG_BIN=$HOME/.pinned-ailang/ailang go test ./host/proctest/ ./host/capsule/ ./host/broker/ -count=1 -p 1 -timeout 600s
  python3 ~/.ailang/state/world-iter192/plan/mutate-r2.py X10 X11 X12 X13 X13b   # 5/5 KILLED
  ```

  **M3 (applied on top of M1+M2):**
  ```
  go vet ./host/...
  AILANG_BIN=$HOME/.pinned-ailang/ailang go test ./host/procbound/ ./host/proctest/ ./host/capsule/ ./host/broker/ -count=1 -p 1 -timeout 600s
  python3 ~/.ailang/state/world-iter192/plan/mutate-r2.py Y1 Y2 Y3 Y4 Y5 Y6 Y7 Y8 Y9 Y10 Y11 Y12 Y13 Y13b Y14 Y15 Z1 Z2 Z3 Z4 Z5 Z10   # 21/21 KILLED
  AILANG_BIN=$HOME/.pinned-ailang/ailang go test ./... -count=1 -p 1
  AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_ail.sh
  ```
  Run the harness from the worktree root. It restores originals in `finally` (safe on
  uncommitted work) — but apply the milestone's hunks fully first; `NOT-APPLIED` means the
  anchor belongs to a later milestone or the text drifted.
- The PR's CI must run the §7 command with `-v`; the controller banks V34/V35 from that log.
- The row's Verification Log records the executor's **own** mutation tallies (not this
  plan's), the per-milestone gate results, and the V34/V35 banked lines.

## §10 Verification Log (planner, iter-193)

| # | Command | Observed |
|---|---|---|
| PV1 | `git status --short`; `git log --oneline -1` | the six prototype paths modified/untracked; HEAD `17c62b0` (r2 carve-out commit) |
| PV2 | `git merge-base --is-ancestor da5f77b HEAD`; `git diff --stat da5f77b HEAD -- host/ cmd/` | ancestor; empty — design base ≡ this base on code |
| PV3 | `go vet ./host/...` | `vet rc=0`, no output |
| PV4 | `AILANG_BIN=$HOME/.pinned-ailang/ailang go test ./host/procbound/ ./host/proctest/ ./host/capsule/ ./host/broker/ -count=1 -p 1 -timeout 600s` | `test rc=0`; `ok …/procbound 0.263s` · `ok …/proctest 0.243s` · `ok …/capsule 15.557s` · `ok …/broker 55.487s`; 0 FAIL (AC3's existing tests included) |
| PV5 | `go test ./host/capsule/ ./host/broker/ -run '^TestPipeCloseBoundsEscapedDescendant$' -count=3 -p 1 -timeout 300s -v` (AILANG_BIN set) | rc=0; capsule PASS ×3 (4.22/4.25/4.31 s), broker PASS ×3 (4.01 s ×3) |
| PV6 | `grep -n 'syscall.Kill' host/*/*_test.go`; `grep -rn 't\.Parallel' host/{capsule,broker}/*_test.go`; `grep -rn 'kill 0\|kill -9\|pkill\|killall' host/` | only as `proctest.Signal(…, syscall.Kill)` callbacks in the new tests (+ one pre-existing verifygate line, live-`cmd.Process` derived); no `t.Parallel`; no shell kills |
| PV7 | `go version`; `$HOME/.pinned-ailang/ailang --version` | `go1.26.6 darwin/arm64`; `AILANG v0.41.0` |
| PV8 | read of `~/.ailang/state/world-iter192/plan/mutations-r2.txt` + `mutate-r2.py` | 46 rows, all KILLED with named failing tests; harness asserts exactly-once anchors, BUILD-FAIL ≠ kill, restores in `finally` |
| PV9 | read of the full prototype (`capsule.go`, `handlers.go`, `procbound{,_test}.go`, `proctest{,_test}.go`, both `cleanup_test.go`) and the design's §3/§5/§7/§11/§12/§15 | basis of §3's per-milestone hunk lists and §4's deviations |
