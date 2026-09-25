# w-daemon-lock-wait-not-deadline-bound (row 22) sprint plan

- Sprint: `w-daemon-lock-wait-not-deadline-bound` · queue row **22** · design doc
  [`w-daemon-lock-wait-not-deadline-bound.md`](w-daemon-lock-wait-not-deadline-bound.md), the
  quorum-approved r2 (arm B). Read its §3 (D1/D2), §5 (AC1–AC5), §6 (mutation matrix) and §7
  (residuals) first.
- Base: `675639607115ab1a8e037eb525893efea6ad3042` (`dev`). No `host/` file changed between the
  design's base `a5090f0` and this base (`git diff --stat a5090f0 6756396 -- host/` is empty), so
  the design's measurements still apply to this tree.
- Planned 2026-09-25 by mission-world iteration 189 (sprint-planner, claude-opus-5-5).
- **The plan is measured, not estimated.** The planner built the whole design in a scratch
  worktree at `/Users/voightkampff/dev/sunholo-data/.planner-wt-iter189` (a sibling of the repo,
  now removed). Every file:line anchor, LOC figure and mutation verdict below comes from that
  prototype. The banked artifacts are under `~/.ailang/state/world-iter189/planner/`:
  - `proto_final.diff`: the full prototype diff (M1 and M2);
  - `lock_wait_order_test.go`: the final test file, including the DV3 assertion;
  - `mut/mut.py`: the mutation harness;
  - `mut/run1.txt` and `mut/run2.txt`: the harness output before and after DV3;
  - `base_*.txt`, `proto_*.txt` and `final_*.txt`: the gate logs.
- Size: about **0.3 d**. Go production +29/−1 (23 of it logic, 6 of it a comment), a new test
  file of 107 lines, and comment rewrites +23/−9. No `.ail` change, no constant change, and no
  change to the status contract.

## §1 Scope

**IN:** design D1 (the startup ordering check in `daemon.New`), the three AC tests, the D2
comment edits, and the three planner deviations DV1–DV3 (§4).

**OUT (the design's own non-goals, §9):** arm A (per-request busy capping); changing a
lock-blocked read from 500 to 503; changing either constant's value; any driver change. The
residuals in the design's §7 stay declared and are not fixed here.

## §2 Baseline on the pristine tree (measured before any edit, requirement 4)

The acceptance commands were run on the pristine tree at `6756396` in the scratch worktree,
with Go `go1.26.6 darwin/arm64` and `AILANG v0.41.0` at `$HOME/.pinned-ailang/ailang`.

| # | Command | Result at base |
|---|---|---|
| BL1 | `go vet ./...` | rc=0, no output |
| BL2 | `AILANG_BIN=$HOME/.pinned-ailang/ailang go test ./... -count=1` | rc=0, **21 `ok`, 0 `FAIL`**, 58 s (`host/daemon` 9.9 s, `host/store` 8.8 s) |
| BL3 | `AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_ail.sh` | rc=0, `✓ verify gate PASSED: 11 required identities verified, 40 named tests pass` |
| BL4 | `go test ./host/daemon/ -run 'ProductionBusyTimeoutConfigured\|BusyTimeoutAtOrAbove\|CheckReadOrdering' -count=1 -v` | `testing: warning: no tests to run` · `ok … [no tests to run]`. **The AC selector is vacuous at base**, so a green run of it proves nothing until the test file exists. Executors must read `=== RUN` lines, not rc. |
| BL5 | `gofmt -l host/ cmd/` (instrument, **not** a gate) | already lists `host/store/schema_version_test.go` and `host/store/store.go` at base. This sprint touches neither, so a non-empty `gofmt -l` afterwards is not the sprint's doing. Check only the touched files: `gofmt -l host/daemon/ host/store/writer_lock.go` must print nothing beyond those two pre-existing store files. |

The design's pre-fix baselines B1 and B2 (§6: nothing in the repo red-lights an R2 retune today,
21/21 green; an R1 retune reds only incidentally, through store and evidence fixtures) were
measured by the designer and controller at `a5090f0`/`c4fd591`. Since `host/` has not changed,
they carry over and are not re-run here.

## §3 Milestones

### M1: startup ordering check and its tests (AC1–AC3)

**Production diff (measured): `host/daemon/daemon.go` +23/−0.** Anchors are at base `6756396`.

| Task | Where (base anchor) | What |
|---|---|---|
| M1.1 | after `host/daemon/daemon.go:236` (`var ErrNonLoopbackBind = …`) | Add `ErrUnorderedTimeouts` and `checkReadOrdering`, **byte-for-byte as the design's D1 code block**: 18 lines including the leading blank line. In the prototype these land at `:249` and `:254`, after DV2's +5 comment lines. |
| M1.2 | after `host/daemon/daemon.go:453`, the closing `}` of the `store.Open` error branch, and **before** `d := &Daemon{` | Add the 5-line call block from D1: `checkReadOrdering(s.BusyTimeout(), readDeadline)`, then `_ = s.Close()`, then `&StartupError{Stage: StageStoreOpen, Detail: "the store's lock-retry window is not below the read deadline", Err: err}`. It must come before `archive.Archive` and `registry.Bootstrap`, because M13 kills a late placement. |
| M1.3 | new `host/daemon/lock_wait_order_test.go` (107 lines) | Three tests: `TestProductionBusyTimeoutConfiguredBelowReadDeadline` (AC1), `TestNewRefusesBusyTimeoutAtOrAboveReadDeadline` (AC2, five DSN arms: `disabled (negative)`, `disabled (zero)`, `just below`, `equal`, `above`) and `TestCheckReadOrderingAcceptsDisabledWindow` (AC3). Start from the design's r1 test (`~/.ailang/state/iter189-design/lock_wait_order_test.r1.go`, 101 lines) **plus the DV3 assertion** (§4). The planner's final file is banked at `~/.ailang/state/world-iter189/planner/lock_wait_order_test.go`. |

**Acceptance (M1)**
- **AC1.** `New` on a plain temp path succeeds. `d.store.BusyTimeout()` is `> 0` (the
  non-vacuity guard; N4 kills its removal) and `< d.readDeadline`. The test goes red under R1,
  R2 and R3 even with each constant's own literal pin updated.
- **AC2.**
  - Accepted arms (`-1ms`, `0`, `readDeadline−1ms`):
    - `New` succeeds;
    - `BusyTimeout()` equals the DSN value;
    - a registry head exists (the positive control).
  - Refused arms (`== readDeadline`, `+1ms`):
    - the error is a `*StartupError` with `Stage == StageStoreOpen`;
    - `errors.Is(err, ErrUnorderedTimeouts)`;
    - **(DV3)** `err.Error()` contains `busy_timeout <window> must be below read deadline <readDeadline>`;
    - `store.Open(dsn)` succeeds afterwards, so writer authority was released;
    - there is no registry head.
- **AC3.** `checkReadOrdering` accepts `-1ms` and `0`.
- The **M1 mutation table (§5.1) is re-run by the executor, and every row is KILLED by its
  named test.** The tally goes in the sprint's Verification Log.

**Gate commands (M1)**, run in order. All must be green:
```
go vet ./...                                   # NOT go build: it does not compile _test.go
go test ./host/daemon/ -run 'ProductionBusyTimeoutConfigured|BusyTimeoutAtOrAbove|CheckReadOrdering' -count=1 -v
                                               # must show 3 top-level + 5 subtest === RUN lines (BL4: vacuous at base)
AILANG_BIN=$HOME/.pinned-ailang/ailang go test ./... -count=1
AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_ail.sh
```

**Non-vacuity against M1's own production diff.** Every row of §5.1 labelled "diff" mutates a
line M1 adds: M3, M4, M5, M6, M7, M8, M10, M11, M12a, M12b, M13, M14, N1, N2 and N3. All of
them are killed. M3 is the behavioural revert of the whole call: it removes the check without a
compile error. Deleting the M1 production diff outright only produces a **compile** red,
because the test references `checkReadOrdering` and `ErrUnorderedTimeouts`. That is not
evidence that any assertion carries weight, so this plan does not count it as a kill.

### M2: comment honesty edits (AC4) and the whole-tree gates (AC5)

**Production diff (measured): comments only.** `host/store/writer_lock.go` +11/−5,
`host/daemon/handlers.go` +12/−4, `host/daemon/daemon.go` +6/−1 (DV2).

| Task | Where (base anchor) | What |
|---|---|---|
| M2.1 | `host/store/writer_lock.go:175-180` (the `busyTimeoutMillis` doc comment; the const stays at `:181` at base, `:187` after the edit) | Delete the false sentence "*the request context remains the outer bound … so the context always wins*". The replacement says three things, per design D2: busy_timeout, not the request context, bounds a lock-blocked read, because the driver's interrupt does not break SQLite's busy-retry sleep; `daemon.New` only validates the configured ordering (`ErrUnorderedTimeouts`), which is not a guarantee; and at 2000 ms against 10 s a lock-blocked read was measured ending at about 2.05 s with `SQLITE_BUSY`, reported as a sanitized 500, which is a margin and not a bound. Keep the final "ride out a writer's commit burst" sentence. The prototype text is in `proto_final.diff`. |
| M2.2 | `host/daemon/handlers.go:299-302`, residual (ii)'s tail, from "Today that composition" to "guarantee." | Replace that tail only. The new text says: the configured ordering is now validated at startup by `checkReadOrdering` and pinned by `TestProductionBusyTimeoutConfiguredBelowReadDeadline`, which is configuration validation, not deadline enforcement; the deadline still does not govern the lock-wait regime, and at production settings a lock-blocked read surfaces as **500 Internal at about busy_timeout**, not 503 (design P3); and tests that shrink `d.readDeadline` after `New` bypass the check by construction. **The `LIMITATION(w-daemon-late-read-503)` line (`:283`), residual (i) (`:290-292`) and the D18 part of residual (ii) (`:293-299`) stay byte-identical.** |
| M2.3 (DV2) | `host/daemon/daemon.go:124-127` (the `readDeadline` doc comment, "*This constant does.*") | Change "This constant does." to "This constant does, EXCEPT for a read blocked on a SQLite lock: …". The addition says busy_timeout bounds a lock-blocked wait, and `New` only validates the configured ordering (`checkReadOrdering`): configuration validation, not deadline enforcement. This is 5 new comment lines. See DV2 for why. |

**Acceptance (M2)**
- **AC4 (a text AC).** Its instrument-health controls were measured on base and on the
  prototype. Under S6 a grep is **not** load-bearing and does not claim to be:

  | Control | base `6756396` | prototype |
  |---|---|---|
  | `grep -c 'always wins' host/store/writer_lock.go` | 1 | **0** |
  | `grep -c 'ORDERING nothing in this code asserts' host/daemon/handlers.go` | 1 | **0** |
  | `grep -c 'checkReadOrdering' host/daemon/handlers.go` | 0 | **1** |
  | `grep -c '500 Internal' host/daemon/handlers.go` | 0 | **1** |
  | `grep -c 'LIMITATION(w-daemon-late-read-503)' host/daemon/handlers.go` | 1 | **1** (unchanged) |
  | residual (i) bytes: `sed -n '/(i) a read that COMPLETES/,/deadline;/p' host/daemon/handlers.go \| shasum` | `43a5caf60884…` | **`43a5caf60884…`** (identical) |
  | `grep -c 'EXCEPT for a read blocked' host/daemon/daemon.go` (DV2) | 0 | **1** |

- **AC5.** On the merged M1+M2 tree:
  - `go vet ./...` is clean;
  - the full suite is green;
  - `verify_ail.sh` still reports 11 identities and 40 named tests, unchanged from BL3.

**Gate commands (M2):** the same four commands as M1, run on the merged tree, plus the seven
AC4 greps above.

**Non-vacuity against M2's own production diff.** M2 changes only comments. No test can kill a
comment mutation, and this plan does not pretend one can. The load-bearing evidence for M2 is
the text controls above, with a measured value on each side. The behavioural claims those
comments make are the design's P2/P3/V19/V20 probe measurements. No test asserts them, by
design (§6).

## §4 Deviations from the design (each forced by a measurement)

- **DV1: D1 adds 23 lines, not the "+21 (r1)" in the design's §4.**
  - Measured: `git diff --numstat` on the prototype gives `23 0 host/daemon/daemon.go` for D1
    alone.
  - The design's own banked `prototype.r1.diff` has hunks `@@ -235,6 +235,24 @@` (+18) and
    `@@ -451,6 +469,11 @@` (+5), which is also 23.
  - The §4 figure is a count slip. The code is identical.
- **DV2: a third false-as-written comment that D2 does not list. The plan adds it to M2.**
  - `host/daemon/daemon.go:124-127` says `readDeadline` "*bounds the ELAPSED TIME of every
    store read a GET handler performs … the wait that happens BELOW the transport, inside
    database/sql. This constant does.*"
  - The design's own P2/P3 measurements (V11, V13) show it does not do this for a
    lock-blocked read: 2.046 s under a 300 ms deadline, and 500 at about 2.05 s under 10 s.
  - Found by `grep -rn 'always wins\|ORDERING nothing\|outer bound\|bounds the ELAPSED\|This
    constant does' host cmd docs` (non-test), which hits `daemon.go:124` and `:127`. The
    only other hit, `archive.go:59`, is about `probeTimeout`, which is out of scope.
  - This is the same honesty class as D2 and costs +6/−1 comment lines. Leaving it would keep
    a false dominance claim one screen above the check that disproves it.
  - **Not edited, flagged for the controller (F1):** `daemon.go:498-499` says of the
    projection's `MaxWait: readDeadline`, "*the same proven bound applies*". The bound it
    means is D15's CPU-bound cancellation, which is proven. The lock-wait case is not, and
    that is the same residual. The plan leaves this alone to keep M2 inside D2's footprint.
- **DV3: the design's test does not enforce its own S7 claim. The plan adds one assertion
  (+6 test lines).**
  - Design D3/S7 says: "*The refusal is a `StartupError` whose Detail and wrapped error name
    both durations.*"
  - Two mutations derived from the diff, both run against the design's r1 test exactly as
    written:
    - **N2** (`Err: err` → `Err: ErrUnorderedTimeouts`, which drops the durations) **SURVIVED**;
    - **N3** (`window, deadline` → `deadline, window` in the `Errorf` arguments, which swaps
      the names) **SURVIVED**.
  - Evidence: `mut/run1.txt`, `N2: … SURVIVED ran=3 rc=0`, `N3: … SURVIVED ran=3 rc=0`.
  - Fix: in AC2's refused arm, after the `errors.Is` check, assert that
    `strings.Contains(err.Error(), fmt.Sprintf("busy_timeout %s must be below read deadline %s", tc.window, readDeadline))`.
    The test also gains a `strings` import.
  - After the fix, both are **KILLED** at `lock_wait_order_test.go:83` (`mut/run2.txt`), and
    all 17 other rows stay KILLED.
  - This asserts an operator-facing message on purpose. It is the only place the refusal names
    the colliding values, which S7 counts as the usage surface.

The design is otherwise implemented **exactly**: the D1 code block is byte-for-byte, the
placement and stage are as specified, and there is no new Stage constant and no store or
constant change.

## §5 Mutation tables (every row EXECUTED on the prototype)

Harness: `~/.ailang/state/world-iter189/planner/mut/mut.py`.
- Each edit is literal, and each must match **exactly once**. Anything else prints `NOAPPLY`,
  which would show up as a harness defect, not a pass.
- The harness then runs **only the row's named test(s)** with an anchored
  `-run '^Test…$/^(subtest)$'` in `./host/daemon/ -count=1 -v`.
- It records the verdict, the number of `=== RUN` lines (a 0 would print `VACUOUS`), and the
  first failing assertion's file:line.
- Last, it restores with `git checkout --` against the committed prototype.
- **Executor warning:** that restore discards uncommitted edits to `daemon.go`, `handlers.go`,
  `writer_lock.go`, `context_read_test.go` and `daemon_test.go`. Commit M1 (and M2) **before**
  running it, or adapt the restore. This is the row-82 lesson: a mutation-drill restore
  destroyed uncommitted work.

Final run (`mut/run2.txt`, with DV3 in place), 20 s wall-clock in total:

### §5.1 M1 mutations

| # | Mutation | Anchored to | Named killing test | Verdict (first failing line) |
|---|---|---|---|---|
| R1 | `busyTimeoutMillis` → 15000 **and** `wantBusyTimeoutMS` → 15000 | the row's scenario | AC1 `TestProductionBusyTimeoutConfiguredBelowReadDeadline` | **KILLED** `:24` New refused at store-open |
| R2 | `readDeadline` → `1 * time.Second` **and** the D7 pin in `daemon_test.go` → 1 s | the row's scenario, the other direction | AC1 | **KILLED** `:24` |
| R3 | `busyTimeoutMillis` → 10000 = `readDeadline` (pin updated) | boundary | AC1 | **KILLED** `:24` |
| M3 | call guard becomes `false && err != nil` | diff (the call) | AC2 `/equal`, `/above` | **KILLED** `:71` "New accepted busy_timeout 10s" |
| M4 | `window >= deadline` → `window > deadline` | diff (boundary) | AC2 `/equal` | **KILLED** `:71` |
| M5 | `s.BusyTimeout()` → `2*time.Second` | diff (source of truth) | AC2 `/equal`, `/above` | **KILLED** `:71` |
| M6 | delete `_ = s.Close()` from the refusal | diff (early return) | AC2 `/equal`, `/above`, re-open assertion | **KILLED** `:89` "another process already holds the writer lock" |
| M7 | `Stage: StageStoreOpen` → `StageConfig` in the refusal | diff (stage) | AC2 `/equal`, `/above`, stage assertion | **KILLED** `:75` |
| M8 | `"%w: busy_timeout …"` → `"%v: …"` | diff (sentinel wrap) | AC2 `/equal`, `/above`, `errors.Is` | **KILLED** `:78` |
| M10 | `readDeadline` → `writeTimeout` at the call | diff (neighbouring constant) | AC2 `/equal`, `/above` | **KILLED** `:71` |
| M11 | arguments swapped: `checkReadOrdering(readDeadline, s.BusyTimeout())` | diff (signature) | AC1 **and** AC2 `/disabled_(negative)`, `/disabled_(zero)` | **KILLED** both runs (`:24`; `:56` "New refused an ordered window -1ms") |
| M12a | `window < 0 \|\| window >= deadline` | diff (disabled-handler semantics) | AC2 `/disabled_(negative)` **and** AC3 | **KILLED** both (`:56`; `:104`) |
| M12b | `window <= 0 \|\| window >= deadline` | diff | AC2 `/disabled_(negative)`, AC2 `/disabled_(zero)` **and** AC3 | **KILLED** all three runs (`:56`, `:56` "0s", `:104`) |
| M13 | check moved after `registry.Bootstrap`, just before `d.integrity = …`, via `d.abort(StageStoreOpen, …)`, so the store is still closed | diff (placement) | AC2 `/equal`, `/above`, no-registry-head assertion | **KILLED** `:93` "refused startup left a registry head" |
| M14 | `window == -time.Millisecond \|\| window >= deadline` | diff (r0's sentinel reading) | AC2 `/disabled_(negative)` **and** AC3 | **KILLED** both (`:56`; `:104`) |
| **N1** *(planner, from the diff)* | the refusal returns the bare `err`, not the `*StartupError` wrapper | diff (wrapper shape) | AC2 `/equal`, `/above`, `errors.As` stage assertion | **KILLED** `:75` |
| **N2** *(planner, from the diff)* | `Err: err` → `Err: ErrUnorderedTimeouts` (drops the durations) | diff (wrapped error) | AC2 `/equal`, `/above`, **DV3** message assertion | **KILLED** `:83`. It **SURVIVED** on the design's test (run1), which forced DV3. |
| **N3** *(planner, from the diff)* | `Errorf` arguments `window, deadline` → `deadline, window` | diff (message) | AC2 `/equal`, `/above`, **DV3** | **KILLED** `:83`. It **SURVIVED** on the design's test (run1). |
| **N4** *(planner, from the test's guard)* | `busyTimeoutMillis` → 0 **and** `wantBusyTimeoutMS` → 0. A disabled handler is *accepted* by the check, so only AC1's `> 0` guard can catch it. | AC1's non-vacuity guard | AC1 | **KILLED** `:29` "store BusyTimeout = 0s; the ordering below would pass vacuously" |

M9 was retired by the design in r1, because its branch no longer exists. It is not run.

**Tally (M1).**
- Final run: **19 run, 19 killed, 0 survived, 0 equivalent.** That is the design's 15 rows plus
  4 planner rows.
- Before DV3 (`run1.txt`): 17 killed and 2 survived (N2 and N3). Those two are the measurement
  behind DV3.

### §5.2 M2 mutations

M2 changes only comments, so there is **no executable mutation for a test to kill**, and the
plan claims none (S6: a comment mutation that "survives" is not evidence of a weak test, and a
grep cannot be passed off as a kill). M2's controls are the AC4 before/after table in §3, and
each row there was measured on both trees. One plausible regression is guarded, but only as an
instrument: reverting M2.2 to "ORDERING nothing in this code asserts" is caught by the
`ORDERING nothing` count going from 0 back to 1. That is a text check.

**Overall tally:** 19 executable mutations, all killed, 0 survived, 0 equivalent. M2 has 7
text controls, each measured.

## §6 Test discipline

- **No wall-clock assertions.** None of the three tests measures elapsed time. Every assertion
  is a construction-time refusal, a value comparison or a string containment. As a stability
  check, the AC selector was re-run 20 times back to back (`-count=20`): `ok … 0.438s`, so none
  of the 20 was flaky.
- No test shrinks `d.readDeadline` after `New` to exercise the check. That path bypasses the
  check by construction (design §7).
- Sub-test names contain parentheses, so anchored `-run` patterns must escape them:
  `disabled_\(negative\)`.

## §7 Prototype gate results (final prototype = M1 + M2 + DV1–DV3)

| Gate | Result |
|---|---|
| `go vet ./...` | rc=0, no output |
| `AILANG_BIN=$HOME/.pinned-ailang/ailang go test ./... -count=1` | rc=0, **21 `ok`, 0 `FAIL`**, 54 s. It was also 21/0 before DV3 (53 s). |
| `AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_ail.sh` | rc=0, `11 required identities verified, 40 named tests pass`, identical to BL3 |
| AC selector `-v` | 3 top-level tests and 5 sub-tests `=== RUN`, all PASS |
| `gofmt -l host/daemon/` | empty |
| `git diff --numstat 6756396 <prototype>` | `29 1 host/daemon/daemon.go` · `12 4 host/daemon/handlers.go` · `107 0 host/daemon/lock_wait_order_test.go` · `11 5 host/store/writer_lock.go` |

## §8 Execution notes for the executor

- Suggested commits: **M1** = `daemon.go` D1 hunks plus the test file. **M2** = the three
  comment edits, including DV2's `daemon.go:124-127` hunk. Stage by name. Never use
  `git add -A`.
- `AILANG_BIN` is mandatory for both `go test ./...` (`host/verifygate` fails loudly without
  it) and `verify_ail.sh`.
- The row's Verification Log must record: the M1 mutation tally from the executor's own re-run
  (not this plan's), the AC4 before/after table, and the four gate results.
- Controller findings carried from the design, not acted on here: O1 (no queue row owns
  residual (i) or the 500-vs-503 classification), O2 (stale line numbers in the row text), and
  this plan's F1 (the projection comment, see DV2).

## §9 Verification Log (planner)

| # | Command | Observed |
|---|---|---|
| PV1 | `git worktree add --detach /Users/voightkampff/dev/sunholo-data/.planner-wt-iter189 HEAD` | `HEAD is now at 6756396` |
| PV2 | `git diff --stat a5090f0 6756396 -- host/` | empty: the design's base and this base agree on `host/` |
| PV3 | BL1–BL5 (§2) | as tabulated |
| PV4 | anchors at base: `grep -n` on `ErrNonLoopbackBind`, `This constant does.`, `proven bound applies`, `Today that composition`, `guarantee.`, `connection (w-daemon-read-cancellation`, `instantly with SQLITE_BUSY` | `daemon.go:236`, `:127`, `:498-499` · `handlers.go:299`, `:302` · `writer_lock.go:176`, `:180`; `sed -n 450,454p daemon.go` shows the `store.Open` error branch closing at `:453` |
| PV5 | `python3 mut/mut.py all` on the pre-DV3 prototype | 17 KILLED, **N2 and N3 SURVIVED** (`mut/run1.txt`) |
| PV6 | `python3 mut/mut.py all` on the final prototype | **19/19 KILLED** (`mut/run2.txt`) |
| PV7 | §7 gates on the final prototype | as tabulated (`final_vet.txt`, `final_test.txt`, `final_verify_ail.txt`) |
| PV8 | AC4 controls on `6756396` and the prototype (`git show <rev>:<path> \| grep -c …`) | as tabulated in §3 M2 |
| PV9 | `git worktree remove --force /Users/voightkampff/dev/sunholo-data/.planner-wt-iter189` | removed after banking. The prototype commits were detached scratch commits and never touched `dev`. |
