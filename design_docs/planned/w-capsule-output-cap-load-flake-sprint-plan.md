# Sprint plan — `w-capsule-output-cap-load-flake` (queue row 20, clause 2)

**Design doc:** [`design_docs/planned/w-capsule-output-cap-load-flake.md`](w-capsule-output-cap-load-flake.md)
(rounds 1–3; cleared to plan by `D-WORLD-23` arm A, ratified attended by Mark 2026-08-20T08:01:31Z, `#68` comment `A`)
**Design-doc base:** `eb1eded` · **Sprint base:** `dev` @ `eb1eded`, working tree clean
**Planner:** mission-control iteration 100, `opus` (`fail-closed:env-pin`)
**Estimate:** 1.0 day, 3 milestones, 3 commits · **Risk:** medium
**Scope:** `host/capsule/capsule.go` + `host/capsule/capsule_test.go` only. No `.ail`, no new package, no new file.

---

## 0. The one-paragraph version

The named flaky test proves the wrong thing with the wrong instrument. Its only witness that
an over-limit child is *killed* is a wall-clock stopwatch that is started **outside** the region
`ExecTimeout` governs. This sprint extracts the post-`Start` output lifecycle behind an unexported,
injectable seam, witnesses the kill through a fake child's own counter instead of a clock, and then
retires the stopwatch. The design doc's spine — "the flake no longer reproduces" is vacuous at base —
is honoured: **every acceptance criterion in this plan was run on the pristine tree this session and
is RED there**, and the residual measured value is stated as a margin, never as an outcome.

---

## 1. Environment — load-bearing, on EVERY command

```sh
export PATH=/opt/homebrew/bin:$PATH
export GOTOOLCHAIN=go1.25.6
export AILANG_BIN=/tmp/ailang-v0300/ailang
```

| Pin | Measured this session | Why it is load-bearing |
|---|---|---|
| `GOTOOLCHAIN=go1.25.6` | ambient `go version` = **`go1.26.4 darwin/arm64`**; pinned = `go1.25.6 darwin/arm64` | go1.26.0–go1.26.5 are on `scripts/verify_go.sh`'s DENY-LIST (they miscompile `host/store/scan.go`). A command without the pin reds in a way the executor cannot diagnose from the failure text. |
| `PATH=/opt/homebrew/bin:$PATH` | `gh`/`go` resolve only with it | `rc=127` here is a PATH gap, never a broken toolchain. |
| `AILANG_BIN=/tmp/ailang-v0300/ailang` | `AILANG v0.30.0`; `stat -f %z` = **`91826738`** | `pinnedBinary` (`capsule_test.go:17`) calls `t.Skip` when unset — **every capsule test silently vanishes** and the package reports `ok`. `verify_go.sh` fails loudly on an unset/wrong `AILANG_BIN`; a bare `go test ./host/capsule` does not. |

**Instrument rule (repeated because this repo has paid for it):** never `<cmd> | tail -N; echo rc=$?`
and never `${PIPESTATUS[0]}` (bash-only; **empty** in zsh). Redirect to a file, `echo "rc=$?"` on the
very next line, then read the file.

---

## 2. Pristine-tree baselines — measured first-party at `eb1eded`

Rule 3e requires every acceptance command to be baselined on the unchanged tree and the base result
recorded *as part of the criterion*. All of the following were run by the planner this session.

### 2.1 Repository gates at base — GREEN, so nothing is being laundered

| Gate | Command | Base result |
|---|---|---|
| `.ail` | `./scripts/verify_ail.sh` | **rc=0** — `✓ world package gate PASSED: 9/9 steps performed non-zero work`, `✓ verify gate PASSED: 10 required identities verified, 39 named tests pass` |
| Go | `./scripts/verify_go.sh` | **rc=0** — build clean, plain + race legs pass, driver-drift gate green, race known-positive control armed |
| `gofmt` | `gofmt -l host/ cmd/` | **empty** (clean) |
| capsule package | `go test ./host/capsule -count=1 -v` | **rc=0**, 7/7 PASS, **0 SKIP**, 17.510 s |
| broker capsule-touching arms | `go test ./host/broker -count=1 -run 'Episode\|EverySubprocessSite' -v` | **rc=0**, incl. `.../host/capsule/capsule.go` subtest PASS |

### 2.2 The five ACs at base — all RED

| AC | Base measurement (planner-run) | Red? |
|---|---|---|
| AC1 | `go test ./host/capsule -run '^TestOutputCollectionOverflowKillsAndOutranksDeadline$' -count=1 -v` → **go-test rc=0**, body `testing: warning: no tests to run` / `ok ... [no tests to run]`; **AC composite rc=1** (the `=== RUN` grep is false) | ✅ |
| AC2 | same shape, `TestOutputCollectionAtLimitDoesNotKill` → go-test rc=0, AC composite **rc=1** | ✅ |
| AC3 | same shape, `TestOutputCollectionTwoOverflowsKillOnce` → go-test rc=0, AC composite **rc=1** | ✅ |
| AC4 | the doc's `python3` heredoc → **rc=1**, `AssertionError: time.Now()`. Planner instrumented the same slice and confirmed **all four** forbidden tokens are present at base: `time.Now()` ✔, `time.Since(` ✔, `elapsed >= clock` ✔, `dbl("0123456789abcdef"` ✔ | ✅ |
| AC5 | same shape, `TestOutputCollectionCallerReleaseUnblocksReadersAndWait` → go-test rc=0, AC composite **rc=1** | ✅ |

**This is the measurement that makes the doc's `-run` warning concrete:** `go test -run '<no match>'`
exits **0** on this toolchain, four times out of four. The exit code is not the tooth; the
`^=== RUN   <name>$` line is.

### 2.3 Design-doc premises re-derived first-party (all TRUE, none refuted)

Run in single calls carrying their own controls:

- `capsule_test.go:233 func TestF6OutputCapReturnsStructuredOverflow`; same-path positive control
  `grep -c '^func Test'` = **7**; fresh-literal negative control = **0**.
- `capsule.go:238 data, err := io.ReadAll(io.LimitReader(pipe, limit+1))`; `capsule.go:237
  func readCapped(...)`. The doc's round-3 `:237 → :238` correction is confirmed correct.
- `grep -nE 'stdoutPipe|stderrPipe' host/capsule/capsule.go` → **exactly four**: `:168`, `:172`
  (creation), `:197`, `:198` (drain use). Same-path control `killOnce` = **2**.
- Kill sites in `capsule.go`: **two** — `:165` `syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)`
  (cancellation, group-wide) and `:193` `killOnce.Do(func() { _ = cmd.Process.Kill() })` (overflow,
  **not** group-wide). The asymmetry is queue row **24**'s, and this sprint must not close it.
- Queue row 24 `w-host-subprocess-cleanup-boundary` is **OPEN** in the charter (`world-mission.md:3642`).

**Inherited, NOT re-measured by the planner** (carried from the doc's §2, load-bearing nowhere in
this plan's criteria): the phase-split timings (resolve 3–12 µs / verify 37.4–46.4 ms idle), the
32-execution non-reproduction sample, and the standalone `shasum` figures. This plan asserts no
threshold that depends on any of them.

---

## 3. THE PLANNER'S OWN MEASUREMENT — mutant blast radius, RUN not predicted

Every mutation in the doc's §5 table was applied to the **base** source in a sibling worktree
(`/Users/voightkampff/dev/sunholo-data/.planner-wt-iter100`, removed after measurement), built, run
unscoped against `./host/capsule` *and* the two `host/broker` arms that consume `capsule.Run`, then
restored by `cp` from `/tmp/w20_backup/` with the SHA-256 re-asserted (`b3b6a41b…95`) and
`git status --porcelain` empty after each.

Base-analogue forms were used because the helper does not yet exist; the mutation is applied to the
identical statement that will move *into* the helper.

| Mut | Base-analogue applied | build | `go vet` | **Tests it kills at base (unscoped)** | Classification |
|---|---|---|---|---|---|
| **M1** | `if false && errors.Is(*dstErr, errOutputLimit)` @ `:192` | rc=0 | rc=0 | **1**: `TestF6OutputCapKillsChildBeyondOnePipeBuffer` — and *only* via `capsule_test.go:303: overflow took 5.043s, i.e. the full 5s bound — the child was not killed` | **SINGLE-TEST** — `-skip` inverse arm valid (verified: `-skip '^TestF6OutputCapKillsChildBeyondOnePipeBuffer$'` → **rc=0**) |
| **M2** | `if true \|\| errors.Is(...)` @ `:192` | rc=0 | rc=0 | **0 — SURVIVES the entire suite** (`./host/capsule` rc=0, 7/7 PASS; broker arms rc=0) | **UNCOVERED BRANCH AT BASE** |
| **M3** | deadline branch (`:208–210`) hoisted above overflow branch (`:205–207`) | rc=0 | rc=0 | **0 — SURVIVES** | **UNCOVERED BRANCH AT BASE** |
| **M4** | `killOnce.Do(...)` → direct kill | rc=0 | **rc=1** with the doc's `_ = killOnce` | **0 — SURVIVES** | **UNCOVERED** + doc defect, see §6 D2 |
| **M5** | `readCapped` returns `data[:limit], nil` @ `:246` | rc=0 | rc=0 | **2**: `TestF6OutputCapReturnsStructuredOverflow` (`:260 error = <nil> <nil>, want *OutputLimitError`) and `TestF6OutputCapKillsChildBeyondOnePipeBuffer` (`:295 error after 5.05s = *capsule.TimeoutError …`) | **NARROW-BLAST (2)** — enumerate, do not `-skip` |
| **M6** | fixed `return Result{}, nil` immediately after `cmd.Start()` | rc=0 | **rc=1** (`unreachable code`) | **7 across 2 packages**: capsule `F2,F3,F4,F5,F6Structured,F6Kills` (6 of 7; `F1` survives because it refuses **before** `Start`) + broker `TestEpisodeLiveReplayThreeArmsAndEvidence` | **BROAD-BLAST** — enumerated set required, see §6 D3 |
| **M7** | `runErr := cmd.Wait()` hoisted above `wg.Wait()` (`:200–201`) | rc=0 | rc=0 | **0 — SURVIVES** | **UNCOVERED BRANCH AT BASE** |

### What this table decides

1. **Four of the seven mutations (M2, M3, M4, M7) survive the whole existing suite.** They are the
   sprint's real teeth, not decoration. The doc's §5 claims are *correct as claims* and now carry a
   measured base status.
2. **M1's only current killer is the assertion AC4 deletes.** `capsule_test.go:303` is the
   `elapsed >= clock` line. Therefore **AC1 must exist and pass before the stopwatch is removed** —
   this is the hard ordering constraint behind the milestone sequence in §5, and the reason
   "delete the flaky test first, add the good ones after" is a forbidden seam.
3. **M6 is broad-blast and must be recorded as an enumerated set.** Its kill of the *named* arm
   (`TestF6OutputCapReturnsStructuredOverflow`) is real, and the other six members are explained:
   every capsule test except `F1` drives `Runner.Run` past `cmd.Start` and reads its result, and
   `broker/TestEpisodeLiveReplayThreeArmsAndEvidence` asserts `got.Stdout == "capsule-transition\n"`
   from a live `capsule.Run`. `F1` survives *by design* (hash-mismatch refusal precedes `Start`) and
   that survival is itself the control that M6 is severing post-`Start` wiring specifically.
4. **The `-skip` inverse arm is licensed for M1 only.** M5 gets a 2-member enumerated set; M6 gets a
   7-member enumerated set; M2/M3/M4/M7 get single-`-run` arms against the new tests they are
   designed to kill (their post-change blast radius must be measured, not predicted — see §7 T3).

---

## 4. Cross-package coupling the executor must not break

These are **not** in the design doc and were found by the planner. Each is a named non-change.

| Coupling | Where | Constraint |
|---|---|---|
| `host/broker/registry_publish_test.go:1075 enumerateSubprocessSites` | AST-walks every non-test `.go` under `host/` and `cmd/` for `exec.Command`/`exec.CommandContext`, groups by **file**, and `t.Fatal`s if `len(drivers) != len(files)` | The single capsule exec site must stay **in `host/capsule/capsule.go`**. Putting the helper in a new file is fine; putting an `exec.Command*` call in a new file is a **cross-package red** with a confusing message. |
| `host/broker/allowlist_test.go:30` | `go list -deps ./host/broker/... ./host/capsule/...` against a module allowlist | The helper must introduce **no non-stdlib import**. Everything it needs (`bytes`, `context`, `errors`, `io`, `sync`) is already imported. |
| `host/broker/episode_test.go:66–75` | live `capsule.New(...).Run(...)`, asserts `stdout == "capsule-transition\n"`, `len(stderr) == 0` | The refactor is **behaviour-preserving**. Any change to what `Run` returns reds a different package. |
| `capsule.go:193` `cmd.Process.Kill()` (not group-wide) | overflow kill | **Queue row 24 owns this.** The adapter must preserve `cmd.Process.Kill()` exactly. Changing it to `syscall.Kill(-pid, …)` inside this tranche reorders separately-owned work and violates `D-WORLD-23` arm A's disposition. |

---

## 5. Milestone plan — 3 milestones, 3 commits, each green on both gates

> Granularity rationale for the seams chosen, and for the ones rejected, is in §5.4.

### Milestone A — extract the output-collection seam (behaviour-preserving)

**Files:** `host/capsule/capsule.go` · **LOC:** ~+45 net (≈55 moved, ≈30 new)
**Tests changed:** none.

**Work**

1. Add an unexported child boundary and helper in `capsule.go`:
   ```go
   // childProcess is the narrow slice of *os.Process the output lifecycle needs.
   // Production supplies the real child; tests supply a fake whose Kill and Wait
   // write their own observables.
   type childProcess interface {
       Kill() error
       Wait() error
   }
   ```
2. Move `capsule.go:180–219`'s post-`Start` body into
   `func collectOutput(ctx context.Context, stdout, stderr io.Reader, limit int64, execTimeout time.Duration, child childProcess) (Result, error, error)`
   — returning `(Result, runErr, err)` so `Run` retains `execPath` for `ExecError`.
   *(Signature shape is the executor's judgment; §5.5 pins the two invariants it must satisfy.)*
3. `Run` adapts the real process:
   ```go
   res, runErr, err := collectOutput(ctx, stdoutPipe, stderrPipe, r.maxOutputBytes, r.execTimeout, cmdChild{cmd})
   ```
   where `cmdChild.Kill()` is **exactly** `cmd.Process.Kill()` and `cmdChild.Wait()` is `cmd.Wait()`.
4. Move the two existing explanatory comments (`:185–188` F6-must-not-decay-into-F5, `:199`
   reads-before-Wait, `:203–204` overflow-outranks-deadline) **with** the code they explain.
5. Add the **liveness precondition** doc comment above `collectOutput`, in the doc's own words
   (§3): the caller owns the finite cleanup bound and must arrange that cancellation makes both
   supplied readers and `Wait` return; `collectOutput` neither makes an arbitrary `io.Reader`
   cancellable nor makes `Wait() error` bounded; production supplies that property via
   `context.WithTimeout` + `exec.CommandContext` + the group-wide `cmd.Cancel` SIGKILL. State the
   residual owner by name: **queue row 24 `w-host-subprocess-cleanup-boundary`**.

**Acceptance — AC-A (milestone-local, planner-authored; NOT in the design doc)**

| | |
|---|---|
| Criterion | With the helper extracted, the **M6 severance mutant** reds an enumerated set **⊇** `{capsule: F2, F3, F4, F5, F6OutputCapReturnsStructuredOverflow, F6OutputCapKillsChildBeyondOnePipeBuffer; broker: TestEpisodeLiveReplayThreeArmsAndEvidence}` and `capsule/F1` still PASSES |
| Base status | **RUNNABLE AT BASE and MEASURED**: the base-analogue M6 kills exactly that 7-member set (§3) |
| Why it exists | Milestone A adds no new test. Its tooth is that the *existing* suite provably covers the extracted path — otherwise this is a refactor commit with no gate. Under rule 3i the observable (`Run`'s returned error/bytes) is downstream of the seam; the writer is the real interpreter through the real pipes, not the assertion. |
| Also asserted | `git diff --stat host/capsule/capsule_test.go` is **empty** (the milestone changes no test), and `go test ./host/capsule -count=1 -v` shows **7 RUN / 7 PASS / 0 SKIP** |

**Gates:** `verify_ail.sh` rc=0, `verify_go.sh` rc=0, `gofmt -l host/ cmd/` empty. **Commit 1.**

---

### Milestone B — the four deterministic state-machine tests (AC1, AC2, AC3, AC5)

**Files:** `host/capsule/capsule_test.go` · **LOC:** ~+215
**Production code changed:** none. The old timing test is **untouched** in this milestone.

**Work**

1. A `fakeChild` whose observables are written **at the boundary, never beside the assertion**
   (doc §3, final paragraph):
   ```go
   type fakeChild struct {
       mu        sync.Mutex
       killCount int
       killed    bool
       waitCount int
       waitEntered chan struct{} // AC5 only; nil elsewhere
       waitRelease chan struct{} // AC5 only; nil elsewhere
   }
   // Kill writes killCount and killed.
   // Wait INDEPENDENTLY READS killed and returns errWaitedWithoutKill if false.
   ```
   `errWaitedWithoutKill` is a test-local sentinel. The test asserts the counter and the `Wait`
   result; it never writes either.
2. **AC1** `TestOutputCollectionOverflowKillsAndOutranksDeadline`: stdout = `bytes.NewReader` of
   exactly `limit+1` bytes, stderr empty, context **already expired** (`WithTimeout(…, 0)` +
   `cancel()` before the call, so `ctx.Err()` is `context.DeadlineExceeded` — asserted in the test
   *before* invoking the helper, so the pre-condition is measured, not assumed). Assert: captured
   stdout is **exactly `limit` bytes**, `killCount == 1`, `Wait` observed the kill (its returned
   error is **not** `errWaitedWithoutKill`), `errors.As(err, &*OutputLimitError)` true, and
   `errors.As(err, &*TimeoutError)` **false**. No duration assertion anywhere.
3. **AC2** `TestOutputCollectionAtLimitDoesNotKill`: exactly `limit` bytes, live context. Assert
   byte identity, `killCount == 0`, `waitCount == 1`, and the fake's wait result surfaces.
4. **AC3** `TestOutputCollectionTwoOverflowsKillOnce`: both prefilled readers over-limit. Assert
   `killCount == 1`.
5. **AC5** `TestOutputCollectionCallerReleaseUnblocksReadersAndWait`: readers whose `Read` signals
   entry on a channel then blocks on a **test-owned** release channel; `Wait` signals entry then
   blocks on a **separate** test-owned release channel. Helper runs in a goroutine. The test
   asserts the ordered handshake: **both reader entries observed → close reader-release → wait
   entry observed → close wait-release → helper returns** with the expected non-overflow result.
   Every handshake and the final return sit behind a phase-specific `select { case <-ch: case
   <-time.After(watchdog): t.Fatalf("phase X never happened") }`. The watchdog is a **failure
   escape, not a product assertion** — it must never be compared against a performance threshold.
   Sizing: 10 s per phase (three orders of magnitude above the sub-ms handshakes; far under
   `verify_go.sh`'s 8 m race timeout and its 600 s wall cap).

**Acceptance:** AC1, AC2, AC3, AC5 exactly as written in the design doc §4, each run with its
`^=== RUN   <name>$` **and** `^--- PASS: <name> ` greps **and** `rc -eq 0`. Base status for each:
**composite rc=1**, measured (§2.2).

Additionally assert, in the same call as the package run,
`grep -c '^--- SKIP' <log>` = **0** — the `pinnedBinary` silent-skip class does not apply to these
four tests (they use no binary), so a SKIP here would mean the package didn't build the way the
executor thinks it did.

**Gates:** both, plus the **race** leg matters here specifically: AC5 runs the helper in a goroutine
and shares `fakeChild` state across it. `verify_go.sh`'s race leg is the check; if it reds, fix the
fake's locking — **never** by removing the goroutine, which would delete the arm's whole point.
**Commit 2.**

---

### Milestone C — retire the timing/throughput oracle (AC4) + the mutation drill

**Files:** `host/capsule/capsule_test.go` · **LOC:** ~−25 net

**Work**

1. From `TestF6OutputCapKillsChildBeyondOnePipeBuffer`, delete the outer stopwatch (`start :=
   time.Now()`, `elapsed := time.Since(start)`, the `elapsed >= clock` fatal and every `elapsed`
   interpolation) and the 64 KiB throughput fixture (`dbl("0123456789abcdef", 13)`).
2. **The function identifier must survive.** See §6 D1 — AC4's script does
   `s.index('func TestF6OutputCapKillsChildBeyondOnePipeBuffer')`, which raises `ValueError` (rc=1,
   AC4 RED) if the function is deleted. Planner verified this by simulating the deletion.
   Retarget the surviving body at real-interpreter production wiring without a clock: run the pinned
   interpreter with a small over-cap source, assert `*OutputLimitError`, assert **not**
   `*TimeoutError`, assert `len(stdout) <= limit`. Rewrite its doc comment to say plainly: the name
   is retained because AC4 keys on it; the throughput/stopwatch oracle was removed as
   non-deterministic; the kill causality is now witnessed by the helper-level fakes in AC1/AC3,
   because — as measured — a real child that is *not* killed still returns `*OutputLimitError`,
   just slowly, so wall time was the only thing the old assertion could see.
3. Run the mutation drill (§7 T3) and record every result.

**Slice trap (measured):** `TestF6OutputCapKillsChildBeyondOnePipeBuffer` is currently the **last**
`func` in `capsule_test.go`, so AC4's slice `s[start:]` runs to EOF (`end` index = `-1`). AC4's end
marker is the next `\nfunc ` occurrence. If a non-`func` declaration (`type fakeChild struct{…}`, a
`var` block) is placed **after** the named test and **before** the next `func`, those lines fall
*inside* the slice and any banned token in them reds AC4. Put the milestone-B fixtures and tests
**before** the named test, or ensure a `func` immediately follows it.

**Acceptance:** AC4 exactly as written (base: **rc=1**, `AssertionError: time.Now()`, with all four
tokens confirmed present), plus AC1/AC2/AC3/AC5 still green, plus the mutation table in §7 T3
complete with **measured** reds.

**Gates:** both. **Commit 3.**

---

### 5.4 Seams chosen and seams rejected

| Seam | Verdict | Reason |
|---|---|---|
| A (extract) → B (new tests) → C (retire old oracle) | **CHOSEN** | Each commit is green on both gates; no commit lands a red tree; the M1-coverage window (§3 finding 2) is never open. |
| "Delete the flaky test first, then build the replacement" | **REJECTED** | Measured: `capsule_test.go:303` is the **only** thing that kills the M1 mutation today. Deleting it before AC1 exists opens a commit-long window where an unconditional-no-kill implementation passes everything. |
| "Tests first, production seam second" | **REJECTED** | Lands a red tree; `go test ./...` is CI's gate and nothing lands red. |
| "One commit for the whole item" (row 21's shape) | **REJECTED here** | The controller asked for separate, individually-green milestones, and A/B/C each have a distinct, independently-checkable tooth. |
| "Split B into four commits, one per AC" | **REJECTED** | The four tests share one `fakeChild`; splitting lands three commits whose only content is a partially-used fixture. |
| "Fold in `CloseOutput`/`WaitDelay`/context-aware wait/joined cleanup errors" | **FORBIDDEN** | `D-WORLD-23` arm A, ratified attended. Queue row **24** owns it and is OPEN. A milestone that reintroduces it is a scope violation, not a bonus. |

### 5.5 The two invariants the helper signature must satisfy (shape otherwise free)

1. **The precedence decision lives inside the helper.** AC1 asserts `*OutputLimitError` wins over an
   already-expired context, and mutation M3 hoists the deadline branch above the overflow branch. If
   either branch stays in `Run`, AC1 cannot see M3 and the mutation is unkillable by its named row.
   `Run`'s post-helper code must not re-decide precedence.
2. **The helper must not need `*exec.Cmd`.** It takes `childProcess`; the real `*exec.Cmd` is
   adapted at the call site. A helper that accepts `*exec.Cmd` cannot be driven by a fake and every
   AC in this sprint becomes unwritable.

The helper also needs `execTimeout time.Duration` (for `TimeoutError{Limit: …}`) and either
`execPath` or a `runErr` return (for `ExecError{Path: …}`). Either is acceptable.

---

## 6. Design-doc defects found by the planner

A refutation is a success. Four are real, one is cosmetic, and **none** of the controller's
handed-down facts were refuted (all re-derived TRUE — §2.3).

### D1 — AC4 silently **requires** the flaky test to survive, while §3's prose reads like a deletion (BLOCKING AMBIGUITY, resolved in-plan)

§3 says "Replace the timing-bearing assertion in `TestF6OutputCapKillsChildBeyondOnePipeBuffer` with
deterministic package-local tests" — which reads as *delete this test*. But AC4's script opens with
`s.index('func TestF6OutputCapKillsChildBeyondOnePipeBuffer')`. **Measured:** if the function is
deleted, that raises `ValueError: substring not found`, python exits **1**, and AC4 is **RED after
the change**.

The doc's later sentence — "Delete the old outer stopwatch and the 64 KiB throughput fixture **from
the named test**" — is the half that AC4 mechanically enforces, so **that half wins**. §5.C follows
it. Named residual: the surviving function's name ("BeyondOnePipeBuffer") is then stale, because
AC4 also bans the 64 KiB fixture that made it true. Renaming it would red AC4 as written, so a
rename is a **doc amendment, not an executor decision**. Flagged for the evaluator.

### D2 — M4's literal form fails `go vet` (copylocks)

The doc's M4 says "retaining all imports (for example `_ = killOnce; …`)". Measured:

```
host/capsule/capsule.go:193:8: assignment copies lock value to _: sync.Once contains sync.noCopy
```

`go build` and `go test` both still pass (copylocks is not in `go test`'s default vet subset), so
this would not *block* — it would just make the executor's "mutant is clean" evidence noisy or, worse,
make them "fix" the mutant into something that no longer neuters the guard. **Fix, measured
working:** `_ = &killOnce` → `go build ./...` rc=0, `go vet ./host/capsule/` **rc=0**.

### D3 — M6's literal form fails `go vet` (unreachable code)

"after `cmd.Start`, bypass the new helper and return a fixed successful `Result`" as an inserted
early `return` measures as:

```
host/capsule/capsule.go:181:2: unreachable code
```

**Fix:** sever at the **call site** instead — replace the `collectOutput(...)` call with zero-valued
locals plus `_ = collectOutput` (keeps the symbol and imports live, no unreachable statement). Same
severance, vet-clean.

Consequently the executor's **"mutant builds"** proof is `go build ./...` **plus**
`go test ./host/capsule -run '^$' -count=1` (measured rc=0; `go build` does **not** compile
`_test.go`) — with `go vet ./host/capsule/` as a *third* signal, not the definition of "builds".

### D4 — M6 is broad-blast, and the doc's warning does not say so

The doc flags M6 correctly ("must not be recordable as killed by compilation failure, global
timeout, missing test discovery, or a skipped capsule test") but frames it as if it kills one row.
**Measured: 7 tests across 2 packages** (§3). Its criterion is therefore an enumerated failing set
with every member explained, not a `-skip` inverse arm. §7 T3 writes it that way.

### D5 — AC4's "**At base:** fails on all four forbidden timing/throughput tokens" is imprecise (cosmetic)

`assert` raises on the **first** violation, so the command demonstrates one token, not four. The
planner instrumented the same slice and confirms all four *are* present, so the claim is TRUE — the
*command* just does not establish it. Recorded so nobody later mistakes the AC's output for a
four-token measurement. No plan change.

### D6 — the doc's own narrative under-states why the stopwatch has to go

**New measurement, not in the doc:** under mutation M1 (kill neutered), the current test fails at
`capsule_test.go:303` — the `elapsed >= clock` line — *after* `errors.As(err, &overflow)` has already
**succeeded**. That is, a child that is never killed still yields `*OutputLimitError`; it just takes
the full 5 s. So wall-clock time is the **only** channel through which the existing test can observe
the kill at all. This strengthens the doc's case (the fake-child counter is not a convenience, it is
the only non-timing witness available) and it is what forces milestone ordering B-before-C.

---

## 7. Task breakdown

| Task | Milestone | Est. | Output |
|---|---|---|---|
| **T1** — `childProcess` + `collectOutput` extraction, adapter, liveness doc comment | A | 1.5 h | `capsule.go` +45 |
| **T2** — `fakeChild` + AC1/AC2/AC3/AC5 tests | B | 3.0 h | `capsule_test.go` +215 |
| **T3** — mutation drill M1–M7, restore controls | C | 2.0 h | recorded table |
| **T4** — AC4 retirement of the stopwatch + comment rewrite | C | 0.75 h | `capsule_test.go` −25 |
| **T5** — three gate runs (`verify_ail.sh` + `verify_go.sh` per milestone), gofmt, 3 commits | A/B/C | 1.25 h | 3 commits |
| Contingency (race-leg remediation on AC5, watchdog tuning) | — | 1.5 h | — |
| **Total** | | **10 h ≈ 1.0 day** | matches the doc's re-priced 1.0 day |

### T3 — the mutation drill, per-mutant protocol

For **each** of M1…M7, in order, six steps, no shortcuts:

1. `cp /tmp/w20_backup_exec/capsule.go host/capsule/capsule.go` — restore first, so the previous
   mutant can never leak into this one. **Take the backup with `cp` at the start of T3;
   NEVER `git checkout -- <file>`, which would delete the executor's uncommitted work.**
2. Apply the mutation; prove it **landed** by `grep -n` on the mutated literal (a mutation that
   silently didn't apply reports a false "survived").
3. Prove it **builds**: `go build ./...` rc=0 **and** `go test ./host/capsule -run '^$' -count=1`
   rc=0. Record `go vet ./host/capsule/` rc as a third signal (use the D2/D3 vet-clean forms).
4. Run the mutant's **named** command below and record the actual `--- FAIL:` lines.
5. Run the **classification** command below and record the *complete* failing set.
6. Restore by `cp`; assert `shasum -a 256` matches the backup and `git status --porcelain --
   host/capsule/` is empty; re-run the pristine control (`go test ./host/capsule -count=1`) → rc=0.

| Mut | Named row | Named command (must go **RED**) | Classification arm |
|---|---|---|---|
| M1 | AC1 | `go test ./host/capsule -run '^TestOutputCollectionOverflowKillsAndOutranksDeadline$' -count=1 -v` → `--- FAIL` | `go test ./host/capsule -count=1 -v`; expect **AC1 (+ possibly AC3)**, and `F6OutputCapKillsChildBeyondOnePipeBuffer` now **PASSES** (its killer was removed in C) — record which |
| M2 | AC2 | `-run '^TestOutputCollectionAtLimitDoesNotKill$'` → `--- FAIL` on `killCount 1, want 0` | full package; **base-measured survivor**, so any kill at all is new coverage. This is M1's required dual arm (`if false &&` cannot neuter a *skip*) |
| M3 | AC1 | `-run '^TestOutputCollectionOverflowKillsAndOutranksDeadline$'` → `--- FAIL` on `*TimeoutError`, want `*OutputLimitError` | full package |
| M4 | AC3 | `-run '^TestOutputCollectionTwoOverflowsKillOnce$'` → `--- FAIL` on `killCount 2, want 1` | full package. Use `_ = &killOnce` (D2) |
| M5 | AC1 | `-run '^TestOutputCollectionOverflowKillsAndOutranksDeadline$'` → `--- FAIL` | full package + `./host/broker -run 'Episode\|EverySubprocessSite'`. **Enumerated set expected ⊇ {AC1, AC2?, F6OutputCapReturnsStructuredOverflow}** — base-measured 2-member set was `{F6Structured, F6Kills}`; explain every member |
| M6 | `TestF6OutputCapReturnsStructuredOverflow` | `go test ./host/capsule -run '^TestF6OutputCapReturnsStructuredOverflow$' -count=1 -v` → assert `^=== RUN`, `^--- FAIL:`, **`grep -c '^--- SKIP' = 0`**, and that the run is **not** a build failure (step 3 already proved it builds) | **BROAD-BLAST**: full package **and** `./host/broker -run 'Episode\|EverySubprocessSite' -v`. Enumerated set must be **⊇ the 7 members measured at base** (§3) with `capsule/F1` still PASS. Sever at the call site (D3) |
| M7 | AC5 | `-run '^TestOutputCollectionCallerReleaseUnblocksReadersAndWait$'` → `--- FAIL` on the ordered handshake, **via the watchdog's phase-specific diagnostic, not a suite hang** | full package. **Base-measured survivor** |

`-skip` inverse arms: permitted for **M1 only** (proven single-test at base; the arm
`go test ./host/capsule -count=1 -skip '^<named>$'` → rc=0 was verified). For every other mutant the
criterion is the enumerated set.

---

## 8. Risks

| Risk | Likelihood | Mitigation |
|---|---|---|
| AC5's blocking fixture reds the `-race` leg | **medium** — a goroutine plus shared fake state is exactly what the race detector is for | Guard all `fakeChild` fields with its mutex; assert entry/exit through channels, never through unsynchronised bools. `verify_go.sh` arms a known-positive race control, so a green race leg is a measurement. If it reds, fix the fake — never delete the goroutine. |
| AC5 hangs the suite instead of failing | medium | Per-phase watchdog with a phase-specific `t.Fatalf`. **No `t.Parallel()` in these tests.** `verify_go.sh` also caps the race leg at 600 s. |
| The executor renames or deletes `TestF6OutputCapKillsChildBeyondOnePipeBuffer` | **measured-real hazard** | D1. AC4 goes RED (`ValueError`). §5.C is explicit; the plan's JSON carries it as a named non-change. |
| The executor "also fixes" the non-group-wide overflow kill at `:193` | low, but it is one line away from the code being moved | Named non-change (§4). Row 24 owns it; `D-WORLD-23` arm A forbids absorbing it. |
| A new file under `host/capsule/` gains an `exec.Command*` call | low | Reds `host/broker`'s AST enumeration with a driver-count fatal (§4). |
| Ambient `go1.26.4` leaks into one command | **measured real** — it is this rig's default | `GOTOOLCHAIN=go1.25.6` on every command; `host/store/toolchain_canary_test.go` + `verify_go.sh`'s deny-list are the backstops. |
| `AILANG_BIN` unset in a focused run → the capsule package reports `ok` with **everything skipped** | measured-real class (bit the iteration-99 controller) | Assert `grep -c '^--- SKIP' = 0` in the same call as every capsule package run. |
| `rc=127` | recurring | PATH gap, not a broken toolchain: `export PATH=/opt/homebrew/bin:$PATH` first. |
| A trailing command's rc impersonates the gate's | measured real, repeatedly | Redirect to a file; `echo "rc=$?"` on the very next line. |

---

## 9. Out of scope / owed elsewhere

- **Queue row 24 `w-host-subprocess-cleanup-boundary` (OPEN)** owns: `CloseOutput() error`,
  `Wait(context.Context) error`, an explicit bounded cleanup context, `cmd.WaitDelay`,
  `errors.Join` surfacing of kill/close/wait failures, **and** the non-group-wide overflow kill at
  `capsule.go:193`. `D-WORLD-23` arm A is the ratified disposition: this tranche keeps scope, weakens
  its claim to exactly what it proves, and records the residual with a named OPEN owner.
- No total-`Run` latency budget, no redefinition of `ExecTimeout`, no hoist/cache/bound of
  `verifyExecutable` (the 91,826,738-byte read+hash at `capsule.go:134` that precedes
  `context.WithTimeout` at `:152` remains unbounded, and this sprint says so rather than fixing it).
- No `host/broker` change. No `.ail`, Z3, kernel, package, CLI, or `docs/QUICKSTART.md` change —
  S7 is not engaged because no user-facing surface is added (the helper is unexported and the
  exported API is byte-identical).
- S1/S2/S3 disposition, per the doc §1: this is host-boundary effect code, so no Z3 contract and no
  `.ail` artifact is appropriate; "why is this not a package?" is answered because nothing is added
  to `world/` and no exported surface changes.

---

## 10. Definition of done

1. Three commits on `dev`, each individually green on `./scripts/verify_ail.sh` **and**
   `./scripts/verify_go.sh` with the pinned environment.
2. AC1–AC5 pass, each recorded **with its measured base result** from §2.2.
3. AC-A recorded with its enumerated M6 set.
4. M1–M7 each recorded with: mutation-landed proof, builds proof, named-row red, blast-radius
   classification (**run, not predicted**), and a byte-identical `cp` restore control.
5. `gofmt -l host/ cmd/` empty; `git status --porcelain -- tools/launchd/ scripts/mission_decisions.sh`
   empty (no driver drift).
6. The doc defects D1–D6 carried into the log entry, with D1 flagged as an **owed doc amendment**
   (the stale test name) rather than silently absorbed.
