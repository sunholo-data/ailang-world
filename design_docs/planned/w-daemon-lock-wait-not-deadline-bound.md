# w-daemon-lock-wait-not-deadline-bound — assert and pin at startup that busy_timeout is below the read deadline; per-request capping is feasible but deferred

- Status: **planned** · Date: 2026-09-25 · Designer: claude-opus-5-5 (iter-189) · Base commit: `a5090f0` · Owning queue row: **22** · Scope class: daemon robustness (clause-2) · Verify profile: go-host (no `.ail` change)
- Prototype: scratch worktree `../.design-wt-iter189` at `a5090f0`, now removed. Its diff, test file, probes and mutation harness are saved under `~/.ailang/state/iter189-design/`.

## §1 Problem

This section restates row 22 using this iteration's measurements. Every claim has a row in §12.

- **P1. Two timeouts, not linked by any code.**
  - `host/store/writer_lock.go:181` sets `const busyTimeoutMillis = 2000` (V1).
  - `host/daemon/daemon.go:133` sets `readDeadline = 10 * time.Second` (V2).
  - No file under `host/` names both (V3).
  - Each value is pinned as a literal in its own package:
    - the store pins at `context_read_test.go:163` and `read_object_test.go:153`;
    - the daemon pins at `daemon_test.go:226`.
  - Nothing pins their **order** (V2, V9).
- **P2. The deadline does not govern a lock-blocked read.** The pinned driver (`modernc.org/sqlite v1.54.0`) handles context cancellation by calling `sqlite3_interrupt` (`sqlite.go:101`, V10). SQLite's busy-retry sleep does not stop on that interrupt:
  - With a 300 ms deadline and a raw handle holding `BEGIN EXCLUSIVE`, a read returned at **2.045763 s**, carrying `context deadline exceeded`.
  - The unlocked control returned in **128.8 µs** (P-A, V11).
  - These numbers reproduce the row's iteration-94 figure of 2.043 s.
- **P3. At production settings the busy window, not the context, ends a lock-blocked read, so the daemon answers 500, not 503.**
  - Measured through the shipped handler: `daemon.New` on a temp DB with a raw `BEGIN EXCLUSIVE` holder. `GET /v1/head` answered **`500 {"class":"Internal","message":"internal store failure"}` in 2.046774 s**, with `readDeadline` at 10 s and `BusyTimeout()` at 2 s.
  - Same-call control, unlocked: `404` in 145 µs (P-C, V13).
  - So the comment at `writer_lock.go:175-180` ("*the context always wins*") is **false as written**, not just unchecked. When the ordering holds, `busy_timeout` wins, the context is still live, and `timedOut` classifies the error as Internal.
- **P4. What is at risk.** Retune either constant into the wrong order and a lock-blocked read overruns its deadline without warning. Today's gates catch this only in one direction, and only by accident (§6 baseline rows B1/B2):
  - **B2:** lowering `readDeadline` to 1 s, with its D7 pin updated, leaves `./host/daemon/` **green** at HEAD.
  - **B1:** raising `busyTimeoutMillis` to 15 s, with `wantBusyTimeoutMS` updated, is caught only by:
    - the store's *other* literal `want 2s` pin (`read_object_test.go:153`);
    - two `host/evidence` tests whose `ObjectReadTimeout` of 3 s is a **test fixture**. `NewValidator` has no production caller (V6).
  - None of these references `readDeadline`.
- **P5. Who can hold the lock is narrow, and is measured, not assumed.**
  - A `mode=ro` reader holding an open read transaction, which is what `store.OpenReadOnly` renders, does **not** block the writer handle's read: 79 µs (P-B 2a).
  - A lock-bypassing raw handle holding `BEGIN EXCLUSIVE` does block it: 2.044 s (P-B 2b).
  - Production uses the driver's default journal mode (no non-test `journal_mode` pragma, V8).

## §2 Decision: arm B, assert and pin the ordering. Arm A is feasible but deferred, for a measured reason.

Row 22 poses the choice. I prototyped both arms against the pinned driver.

| Arm | Verdict | Measured reason |
|---|---|---|
| **A: the deadline dominates.** Cap the effective busy wait by the remaining context budget. | **Feasible, deferred** | **What works:** pin a connection (`db.Conn(ctx)`), run `PRAGMA busy_timeout = min(configured, remaining)`, run the read, then restore. Under a 300 ms deadline this returned at **315.7 ms**, a 15.7 ms overrun from SQLite's sleep granularity. Under 700 ms it returned at 721.7 ms. Each PRAGMA costs **1.16 µs**, against 5.38 µs per unlocked read (P-A). Sources: no busy-handler API in the driver (`grep -c BusyHandler` → 0, V10), and `RegisterConnectionHook` runs only at connection open. **What it costs, measured:** (1) The restore cannot use the request context, because by then it has expired: `restore-with-req-ctx err=context deadline exceeded`, `restore-with-bg err=<nil>` (P-A). (2) A skipped or failed restore **poisons the pool**. A `db.Conn` that was capped to 37 ms and returned without a restore leaves `PRAGMA busy_timeout` reading **37** on the next pool use, against a configured 2000 (P-B 1). The store is `SetMaxOpenConns(1)` (V8), so that connection is the **writer's** connection: the next `Commit` would fail `SQLITE_BUSY` at once instead of riding out a commit burst. (3) The work lives in the store, not the daemon, because the daemon holds no `*sql.DB`. Every one of the **8** ctx-taking `(*Store)` methods (V14) would have to move from `s.db` to a pinned connection with set and restore. That is the same signature surface row 23 is threading contexts through. |
| **B: the ordering is asserted and pinned.** | **Chosen** | This is the smallest honest change, and it mirrors this repo's own precedent: `evidence.NewValidator` reads `reader.BusyTimeout()` and refuses with `ErrUnorderedTimeouts` (`validator.go:85-90`, V5). The daemon already holds the opened `*store.Store` (`daemon.go:456`, V5), and `(*Store).BusyTimeout()` returns the **effective**, DSN-resolved window (`store.go:293`, V5). Prototype: +23 lines in `daemon.go` and a 97-line test file. **14 of 14 mutations killed, and both pre-fix baseline retunes were measured (B1 caught only by accident, B2 green)** (§6). |

**Why defer A instead of rejecting it.** P5 shows the lock-wait regime is reachable only by a lock-bypassing process, and P3 shows the daemon already bounds it at 2 s. A's real value is a *503 at the deadline* instead of a *500 at the busy window*. That is a status-contract change, which is what the missing successor for residual (i) owns (§11 O1). A's poisoning hazard lands on the writer's only connection, so paying that cost belongs in the item that also changes the status contract. It should not ride on a 0.5 d ordering fix.

## §3 Design

### D1 — Startup refusal in `daemon.New`, from the opened store's effective window

`host/daemon/daemon.go` gains a sentinel, a pure helper, and one call site. The code below is the prototype's exact shape:

```go
// ErrUnorderedTimeouts is the named sentinel for a store whose SQLite lock-retry
// window (busy_timeout) is not strictly below the read deadline. ...
var ErrUnorderedTimeouts = errors.New("daemon: store busy_timeout is not below the read deadline")

func checkReadOrdering(window, deadline time.Duration) error {
	if window < 0 {
		return fmt.Errorf("%w: busy_timeout is unknown (%s)", ErrUnorderedTimeouts, window)
	}
	if window >= deadline {
		return fmt.Errorf("%w: busy_timeout %s must be below read deadline %s", ErrUnorderedTimeouts, window, deadline)
	}
	return nil
}
```

The call sits **immediately after `store.Open` succeeds, before `d` is built**:

```go
	if err := checkReadOrdering(s.BusyTimeout(), readDeadline); err != nil {
		_ = s.Close()
		return nil, &StartupError{Stage: StageStoreOpen,
			Detail: "the store's lock-retry window is not below the read deadline", Err: err}
	}
```

**Design points**

- **The window comes from `s.BusyTimeout()`, not the constant.** An operator can set the window through the `--db` DSN, for example `file:w.db?_pragma=busy_timeout(20000)`. `withBusyTimeout` never overrides an explicit caller value, and `busyTimeoutFromParams` reports the value the driver actually applied, first-wins (`read_object_test.go:247-280`). A check against the constant would let a DSN-configured window of 20 s through. Mutation M5 kills that.
- **The comparison is strict: `window >= deadline` refuses.** At equality, the lock wait ends at the deadline plus up to one SQLite sleep interval, which is an overrun. Mutation M4 kills a `>` comparison.
- **The check runs before any lifecycle side effect.** The refusal must happen before `archive.Archive` and `registry.Bootstrap` write to the store. The test proves that no registry head exists after a refusal, with a control showing an accepted startup does create one. Mutation M13 kills a late placement.
- **Stage `StageStoreOpen`, with no new stage.** The incompatibility belongs to the store the daemon just opened. Adding a Stage constant would widen an enumerated public set for no new operator decision.
- **The `window < 0` branch cannot be reached through `New` today, and I measured that.** `store.Open(":memory:?_pragma=%zz")` fails first with `invalid URL escape` (P-C). The branch is kept because `BusyTimeout()`'s contract allows −1 (`busyTimeoutFromDSN`, `writer_lock.go:246`). It is killed at the helper level (M9, M12).
- **Window 0 is accepted.** An in-memory store has no busy window (control: `store.Open(":memory:")` → `0s`, P-C). With no busy window, a lock conflict fails immediately, so it is trivially ordered.
- **The projection is covered too.** `projection.Config.MaxWait` is also `readDeadline` (`daemon.go:511`, V7). The one check covers both consumers.

### D2 — Honesty edits: retire the false claim and narrow the LIMITATION

- **`host/store/writer_lock.go:175-180`.** Replace "*the request context remains the outer bound, and 2000 ms sits well below the daemon's 10 s read deadline so the context always wins*". P3 shows that sentence is false as written. The replacement says:
  - `busy_timeout`, not the request context, bounds a lock-blocked read, because the driver's interrupt does not break SQLite's busy-retry sleep;
  - `daemon.New` refuses to start unless this window is strictly below its read deadline (`ErrUnorderedTimeouts`);
  - at 2000 ms against 10 s, a lock-blocked daemon read ends at the busy window with `SQLITE_BUSY`, which the daemon reports as a sanitized 500.
- **`host/daemon/handlers.go:283-309`.** Keep the `LIMITATION(w-daemon-late-read-503)` tag and residual (i) word for word. Rewrite residual (ii) from "*an ORDERING nothing in this code asserts, not a guarantee*" to say three things:
  - the ordering is now **asserted at startup** by `checkReadOrdering` and pinned by `TestReadDeadlineDominatesProductionBusyTimeout`;
  - the lock-wait regime still exists: the deadline still does not govern it, and at production settings it surfaces as **500 Internal at ~busy_timeout** (P3), not 503;
  - tests that shrink `d.readDeadline` *after* `New` bypass the check by construction (§7).

### D3 — What does NOT change

- No store change. `busyTimeoutMillis`, `withBusyTimeout`, `BusyTimeout()` and all store read methods are unchanged. Only the store comment changes.
- No constant changes. `readDeadline` stays 10 s and `busy_timeout` stays 2000 ms. The D7 pin and the store pins stay.
- No status-contract change. A lock-blocked read still answers 500 at the busy window (P3).
- **Z3 / pure core:** no `.ail` file changes, and no contract is added or removed. This is a host-boundary configuration check. S1 does not apply, and S2 is respected because the check sits in the effectful constructor.
- **S3 "why not a package?"** This is daemon startup configuration validation, which is host code by definition.
- **S7:** no new usage surface. The refusal is a `StartupError` whose Detail and wrapped error name both durations. `docs/QUICKSTART.md` does not mention `busy_timeout` (V15), so no quickstart step changes.

## §4 Files

| File | Change | Prototype LOC |
|---|---|---|
| `host/daemon/daemon.go` | `ErrUnorderedTimeouts`, `checkReadOrdering`, call in `New` | +23 |
| `host/daemon/lock_wait_order_test.go` (new) | 3 tests (§5) | +97 |
| `host/daemon/handlers.go` | residual (ii) comment rewrite (D2) | ~±15 (comment only) |
| `host/store/writer_lock.go` | `busyTimeoutMillis` comment rewrite (D2) | ~±6 (comment only) |

## §5 Acceptance criteria

- **AC1. Reorder test.** `TestReadDeadlineDominatesProductionBusyTimeout` calls `New` on a plain temp path. It must succeed, and `d.store.BusyTimeout()` must be `> 0` (non-vacuity guard) and `< d.readDeadline`. It **must go red when the two constants are retuned into the wrong order**, even when each constant's own literal pin is updated to match (R1, R2, R3).
- **AC2. The refusal is driven by the opened store's value.** `TestNewRefusesBusyTimeoutAtOrAboveReadDeadline` sets the window through the DSN at `readDeadline − 1ms` (accepted), `== readDeadline` (refused) and `+ 1ms` (refused).
  - The accepted arm asserts that `BusyTimeout()` equals the DSN value. This makes sure the arm is really driving the window.
  - Each refused arm asserts all of the following:
    - a `*StartupError` with `Stage == StageStoreOpen`;
    - `errors.Is(err, ErrUnorderedTimeouts)`;
    - writer authority is released: `store.Open` on the same DSN succeeds;
    - no registry head was written, with the accepted arm as the positive control.
- **AC3. Helper branches.** `TestCheckReadOrderingRefusesUnknownWindow`: a window of `-1` is refused with `ErrUnorderedTimeouts`, and `0` is accepted.
- **AC4. Honesty edits land (D2).** The `writer_lock.go` sentence "the context always wins" is gone. Residual (ii) in `handlers.go` names `checkReadOrdering` and the measured 500-at-busy-window behaviour. The `LIMITATION(w-daemon-late-read-503)` tag and residual (i) are unchanged. This AC is a text AC. A grep is its instrument-health control and is not claimed as load-bearing (S6).
- **AC5. Gates.** `go vet ./...` is clean. `go test ./... -count=1` is green with `AILANG_BIN` set to the pinned v0.41.0. `./scripts/verify_ail.sh` is unaffected because no `.ail` file changes.

**No wall-clock assertions.** No new test measures elapsed time: every AC is a construction-time refusal or a value comparison. The only timings in this doc are probe measurements (P-A, P-B, P-C) that justify the decision, and none is asserted by a test. That follows the repo's measured history of 1 s bounds flaking 3/3 under full-suite load.

## §6 Test plan and mutation matrix — every row was RUN against the prototype

Harness: `~/.ailang/state/iter189-design/.mut/run.sh`. It applies a Python edit, runs `go test -run 'ReadDeadlineDominates|BusyTimeoutAtOrAbove|CheckReadOrdering' -count=1`, and restores the files. R1/R3 also run `./host/store/`. The B rows run whole packages with no `-run` filter.

| # | Mutation | Anchored to | Killed by | Result |
|---|---|---|---|---|
| R1 | `busyTimeoutMillis` → 15000 **and** `wantBusyTimeoutMS` → 15000 (retune with the pin updated) | the row's own scenario | AC1 `TestReadDeadlineDominatesProductionBusyTimeout` | **KILLED** |
| R2 | `readDeadline` → 1 s **and** the D7 pin → 1 s (retune with the pin updated) | the row's scenario, other direction | AC1 | **KILLED** |
| R3 | `busyTimeoutMillis` → 10000 = `readDeadline` (equality retune, pin updated) | boundary | AC1 | **KILLED** |
| M3 | drop the startup check (`false && err != nil`) | the fix | AC2 `/equal`, `/above` | **KILLED** |
| M4 | `window >= deadline` → `window > deadline` | the fix (boundary) | AC2 `/equal` | **KILLED** |
| M5 | read a constant `2*time.Second` instead of `s.BusyTimeout()` | the fix (source of truth) | AC2 `/equal`, `/above` | **KILLED** |
| M6 | omit `_ = s.Close()` on refusal | **the diff** (new early return after `store.Open`) | AC2 re-open assertion | **KILLED** |
| M7 | refuse with `Stage: StageConfig` | the diff (stage choice) | AC2 stage assertion | **KILLED** |
| M8 | `%w` → `%v` (sentinel not wrapped) | the diff (new sentinel) | AC2 `errors.Is` | **KILLED** |
| M9 | drop the `window < 0` branch | the diff (helper branch) | AC3 | **KILLED** |
| M10 | compare against `writeTimeout` instead of `readDeadline` | the diff: the neighbouring constant in the same `const` block, which D2's comment also names | AC2 `/equal`, `/above` | **KILLED** |
| M11 | arguments swapped: `checkReadOrdering(readDeadline, s.BusyTimeout())` | the diff (helper signature) | AC1 and AC2 `/just_below`, `/above` | **KILLED** |
| M12 | `window < 0` → `window <= 0` (zero refused) | the diff (in-memory stores) | AC3 | **KILLED** |
| M13 | move the check after `registry.Bootstrap`/projection, still closing the store | the diff (placement) | AC2 no-registry-head assertion | **KILLED** |

**Tally: 14 mutations run on the prototype, 14 killed, 0 survivors, 0 equivalent.** Adding the two baseline rows below gives **16 runs in total**, plus one isolation rerun of B1.

**Pre-fix baseline.** These rows use HEAD's `daemon.go` with the new test file removed, and show what the gates catch today:

| # | Mutation at HEAD | Result | Detail |
|---|---|---|---|
| B1 | R1 at HEAD, running `./host/daemon/ ./host/store/ ./host/evidence/` | red, **incidentally** | A dedicated rerun gives `host/daemon` **ok**. The reds are `host/store` `TestBusyTimeoutCachesEffectiveDSNAndDoesNotBlock/default` ("BusyTimeout = 15s; want 2s", a second literal pin) and `host/evidence` `TestOversizeProofReportIsRefused` plus `TestConstructorPinsBusyTimeoutBelowObjectReadTimeout` ("ObjectReadTimeout 3s must exceed BusyTimeout 15s", a test fixture). None references `readDeadline`. The first combined run also showed daemon failures, and the isolated rerun is green. I attribute them to load in the parallel run and do **not** count them as detection. |
| B2 | R2 at HEAD, running `./host/daemon/` | **SURVIVED (green)** | Nothing today asserts the order in this direction. |

## §7 What is NOT fixed (declared residuals)

- **The lock-wait regime still ignores the deadline.** Arm B guarantees only that the busy window ends first. A lock-blocked read is bounded by `busy_timeout`, currently 2 s, and in production it answers **500 Internal**, not 503 (P3).
- **Short test deadlines still overrun.** `d.readDeadline` is a field that tests set after `New` (`daemon.go:296-302`), so the ordering is not checked for them. A lock-blocked read under such a deadline still overruns it by up to the busy window. The existing read-deadline tests avoid this by using an already-expired stimulus (`expiredReadDeadline`), not a lock.
- **Residual (i) is untouched.** A read that completes before cancellation still answers 200. That is `LIMITATION(w-daemon-late-read-503)`, and its tag and text are kept.
- **The `window < 0` branch cannot be reached through `New`** (D1). It is killed only at the helper level.

## §8 Conflict surface

- **Row 23 (`w-store-deadline-free-residue-owner`): no collision.** Row 23 threads caller contexts through `host/broker/approve.go` (8), `host/registry/registry.go` (2) and `host/replay/replay.go` (1), and shrinks `deadlineFreeReadPins` in `host/store/context_read_test.go:370-374` (V16). This doc touches none of those files. Its only `host/store` edit is a comment in `writer_lock.go`. **Arm A, if revived, would collide.** It rewrites the ctx-taking store methods that row 23's callers thread through, which is another reason to sequence A after row 23.
- **Row 26 (`w-bounded-z3-report-producer`): no collision.** Its scope is `host/evidence/proof_producer.go` and envelope production. This doc does not touch `host/evidence`. It copies the `ErrUnorderedTimeouts` *pattern* but declares its own daemon-local sentinel, because importing `evidence` into `daemon` for a sentinel would couple unrelated packages. B1 shows `host/evidence`'s fixture ordering is not a daemon guard, and this doc leaves it alone.
- **`w-daemon-late-read-503` (residual (i)): no collision.** There is no row or doc to collide with (§11 O1). This doc edits the same `handlers.go` comment block, and it changes only residual (ii)'s paragraph.
- **The same package, `host/daemon`:** `daemon_test.go:226`'s D7 literal pin is read, not modified. `read_deadline_test.go` is not modified. The full `./...` suite stayed green on the prototype (V17).

## §9 Non-goals

- No per-request busy capping (arm A, deferred, §2).
- No change to the status for a lock-blocked read (500 → 503). That is O1's successor.
- No change to either constant's value.
- No driver change and no upstream ask. The driver's behaviour matches documented SQLite semantics: `sqlite3_interrupt` does not end a busy-handler sleep.
- No edit to the charter. Findings for the controller are listed in §11.

## §10 Milestones (≈0.3 d total)

- **M1 — startup ordering check + tests (~120 LOC).** Add D1 to `daemon.go` and create `lock_wait_order_test.go`.
  - *Acceptance:* AC1–AC3 green, and every mutation from R1 to M13 is re-run and killed, with the tally recorded in the sprint's Verification Log.
- **M2 — honesty edits (~25 LOC, comments only).** Make the D2 rewrites in `writer_lock.go` and `handlers.go`.
  - *Acceptance:* AC4, plus AC5 on the merged tree: `go vet ./...` and `AILANG_BIN=$HOME/.pinned-ailang/ailang go test ./... -count=1` both green.

## §11 Open items and findings for the controller (not acted on here)

- **O1. Residual (i)'s named owner does not exist, which is row 23's defect class again.**
  - `LIMITATION(w-daemon-late-read-503)` names a successor item.
  - `grep -nE '^[0-9]+\. \*\*w-daemon-late-read' design_docs/world-mission.md` → **no rows**, while the same pattern shape finds row 22 (control, V18).
  - The charter mentions the name only twice, both in prose, and there is **no** design doc for it (V18).
  - Suggestion: queue a row that owns residual (i) *and* the 500-vs-503 classification of a lock-blocked read (P3). Arm A (§2) would be one candidate mechanism there.
- **O2. The row text cites stale line numbers.** It gives `writer_lock.go:179` and `daemon.go:128`; at `a5090f0` they are `:181` and `:133` (V1, V2).

## §12 Verification Log

All commands were run at base `a5090f0` in `/Users/voightkampff/dev/sunholo-data/ailang-world`, unless marked as run in the prototype worktree. Toolchain `go1.26.6 darwin/arm64`; driver `modernc.org/sqlite v1.54.0`; `AILANG v0.41.0` at `$HOME/.pinned-ailang/ailang`.

| # | Command | Observed |
|---|---|---|
| V1 | `grep -n 'const busyTimeoutMillis\|context always wins' host/store/writer_lock.go` | `178:// sits well below the daemon's 10 s read deadline so the context always wins.` · `181:const busyTimeoutMillis = 2000` |
| V2 | `grep -n 'readDeadline = 10' host/daemon/daemon.go; grep -n '{"readDeadline", readDeadline' host/daemon/daemon_test.go` | `133:	readDeadline = 10 * time.Second` · `226:			{"readDeadline", readDeadline, 10 * time.Second},` |
| V3 | `grep -rln busyTimeout host \| xargs grep -ln readDeadline`; control `grep -rln readDeadline host` | check: empty, rc=1. Control fires on **6** files: `archive/archive_test.go`, `archive/archive.go`, `daemon/read_deadline_test.go`, `daemon/handlers.go`, `daemon/daemon.go`, `daemon/daemon_test.go` |
| V4 | `grep -n 'LIMITATION(w-daemon-late-read-503)\|ORDERING nothing in this code asserts' host/daemon/handlers.go` | `283:// LIMITATION(w-daemon-late-read-503): …` · `301://	     request deadline — an ORDERING nothing in this code asserts, not a` |
| V5 | `grep -n 'window := reader.BusyTimeout()\|ErrUnorderedTimeouts, cfg' host/evidence/validator.go`; `grep -n '^func TestConstructorPins…\|^func TestRealStoreBlocked…' host/evidence/*_test.go`; `grep -n '^func New(cfg Config)\|cfg: cfg, store: s, reads: s\|readDeadline: readDeadline' host/daemon/daemon.go` | `85:	window := reader.BusyTimeout()` · `90: … ErrUnorderedTimeouts …` · `realstore_test.go:144` and `:247` · `431:func New(cfg Config)` · `456: cfg: cfg, store: s, reads: s, …` · `457: readDeadline: readDeadline, …`. Also `store.go:293 func (s *Store) BusyTimeout() time.Duration` (read directly). |
| V6 | `grep -rn 'NewValidator(' host --include='*.go' \| grep -v _test.go`; control `grep -rln 'NewValidator(' host \| grep _test` | check: only the definition, `validator.go:63`, so **no production caller**. Control: `authority_test.go`, `realstore_test.go` |
| V7 | `grep -n 'MaxWait: readDeadline' host/daemon/daemon.go` | `511:		MaxWait: readDeadline,` |
| V8 | `grep -n 'SetMaxOpenConns(1)' host/store/store.go`; `grep -rn journal_mode host/store --include='*.go' \| grep -v _test`; control `grep -rln journal_mode host/store` | `305:	db.SetMaxOpenConns(1)`. Non-test `journal_mode`: none. Control fires on `context_read_test.go` |
| V9 | `grep -n 'wantBusyTimeoutMS = 2000' host/store/context_read_test.go; grep -n 'want 2s' host/store/read_object_test.go` | `163:const wantBusyTimeoutMS = 2000` · `153: t.Fatalf("BusyTimeout = %v; want 2s", got)` |
| V10 | `grep -n 'modernc.org/sqlite' go.mod`; `grep -n 'func RegisterConnectionHook\|c.interrupt(c.db)' $GOMODCACHE/modernc.org/sqlite@v1.54.0/sqlite.go`; `grep -c 'BusyHandler\|busy_handler' …/sqlite.go` | `7: modernc.org/sqlite v1.54.0` · `101: c.interrupt(c.db)` · `645:func RegisterConnectionHook` · busy-handler count **0** |
| V11 (P-A) | prototype worktree: `go run ./probe189` (source `~/.ailang/state/iter189-design/probe189/main.go`) | `control unlocked read: 128.833µs` · `per-request PRAGMA exec: 1.159µs; per unlocked read: 5.38µs` · `deadline 300ms, today (busy 2000ms): elapsed 2.045763333s, err=context deadline exceeded` · `deadline 300ms, arm A (busy capped to 299ms): elapsed 315.730416ms overrun 15.730416ms, err=database is locked (5) (SQLITE_BUSY) …; restore-with-req-ctx err=context deadline exceeded, restore-with-bg err=<nil>` · `deadline 700ms, today: elapsed 2.045831834s, err=database is locked (5)` · `deadline 700ms, arm A: elapsed 721.658834ms overrun 21.658834ms` · `pool busy_timeout after arm A: 2000` |
| V12 (P-B) | prototype worktree: `go run ./probe189b` | `(1) pool busy_timeout after a skipped restore: 37 (configured 2000)` · `(2) writer store BusyTimeout()=2s` · `(2a) writer-store read while a mode=ro reader holds SHARED: 79.209µs err=<nil>` · `(2b) writer-store read while a raw handle holds EXCLUSIVE: 2.044035375s err=store: read selected head: context deadline exceeded` |
| V13 (P-C) | prototype worktree: scratch `zz_probe189_test.go`, `go test ./host/daemon/ -run ZZProbe -v` (deleted after the run) | `store.Open(":memory:?_pragma=%zz") err=store: enable foreign keys: invalid URL escape "%zz"` · `control store.Open(:memory:) BusyTimeout=0s` · `control unlocked GET /v1/head: 404 in 145.041µs` · `lock-blocked GET /v1/head (deadline 10s, busy 2s): 500 in 2.046774458s body={"error":{"class":"Internal","message":"internal store failure"}}` |
| V14 | `grep -n '^func (s \*Store) [A-Za-z]*(ctx context.Context' host/store/*.go \| grep -v _test \| wc -l`; `grep -c 'd\.reads\.' host/daemon/handlers.go host/daemon/daemon.go` | `8` · `handlers.go:5`, `daemon.go:1` |
| V15 | `grep -n 'busy_timeout\|busy' docs/QUICKSTART.md` | no hits. Control: the same `grep -n` shape fires on `writer_lock.go` (V1) |
| V16 | `grep -n 'deadlineFreeReadPins = ' -A5 host/store/context_read_test.go` | `370:var deadlineFreeReadPins = map[string]int{` · `approve.go: 8` · `registry.go: 2` · `replay.go: 1` |
| V17 | prototype worktree: `go vet ./...`; `go test ./host/daemon/ ./host/store/ -count=1`; `AILANG_BIN=$HOME/.pinned-ailang/ailang go test ./... -count=1` | `VET-OK` · `ok host/daemon 2.632s` · `ok host/store 6.435s` · full suite: **21 `ok`, 0 `FAIL`**. Without `AILANG_BIN`, `host/verifygate` fails loudly by design ("AILANG_BIN is unset …"), as expected. |
| V18 | `grep -nE '^[0-9]+\. \*\*w-daemon-late-read' design_docs/world-mission.md`; control `grep -nE '^[0-9]+\. \*\*w-daemon-lock-wait' …`; `grep -c 'w-daemon-late-read-503' …`; `ls design_docs/planned design_docs/implemented \| grep -c late-read` | check: no rows, rc=1. Control: `4219:22. **w-daemon-lock-wait-not-deadline-bound**` · 2 prose hits · 0 docs |
| V19 | prototype worktree: mutation harness, 14 fix rows plus 2 baseline rows plus the B1 isolation rerun | as tabulated in §6. The raw lines are in this iteration's transcript, and the harness is at `~/.ailang/state/iter189-design/.mut/run.sh`. |

**Controller claims checked:**
- F1: TRUE, and stronger. The comment is not only unchecked but **false as written** (P3).
- F2: TRUE.
- F3: the conclusion is TRUE, but the **stated control output is FALSE**. It listed 3 files, and the control fires on 6 (V3).
- F4: TRUE, at slightly different lines, 283-309.
- F5: TRUE.
- F6: TRUE, reproduced at 2.0458 s (V11).
