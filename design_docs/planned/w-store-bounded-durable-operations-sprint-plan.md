# Sprint plan — w-store-bounded-durable-operations (queue row 23, iteration 197)

**Base:** detached `dev e33378a` plus the designer's uncommitted prototype. **Authority:** [design](w-store-bounded-durable-operations.md), especially §1.2, §2, §5, §6.1, §9, §11 and the round-2 carve-out in §12. **Scope:** M1 only, split below into M1a–M1d. This planner makes no git write and leaves production code untouched. **Merge gate:** nothing in this tranche merges before Mark ratifies the §1 bound table (§9). `D-WORLD-37` = A is the policy direction, not ratification of the proposed numbers.

M1 supplies exactly row 106 M5's context-aware durable tail and the store-side `Close` invariant. It installs **no duration value**, daemon budget, guard, quarantine, HTTP mapping, or coordinator caller migration. M2–M7 await Mark's bound-table ratification. Keep the three context-free wrappers temporarily; remove them in **M6b**, not M4 as the prototype comments say. The three new census pins in `context_roots_test.go` must say “removed by M6b”; `GetEffectReceipt`'s local Background root is likewise removed by M6b's migration. Fix the `UncertainError` comment to say a bare commit must reconcile by its log-index witness (§2.5), never `SelectedHead`.

## 1. Gates and measurements

For **every** Go command: `export PATH=/opt/homebrew/bin:$PATH; export AILANG_BIN=$HOME/.pinned-ailang/ailang` (pinned AILANG v0.41.0). The compile fence is **both** `go vet ./...` and `go test ./... -run '^$'`; `go build` misses test compilation. I copied the worktree into `.plan-scratch/`, assembled cumulative states, and ran both commands at each boundary. “Production LOC” counts added nonblank, noncomment source lines, approximately; it excludes tests and does not subtract replaced old lines.

| Boundary | Cumulative assembled state / added production LOC at this boundary | `go vet ./...` rc | `go test ./... -run '^$'` rc |
|---|---|---:|---:|
| M1a | Entire prototype: `durable.go` 68 + `store.go` 20 + `journal.go` 23 ≈ **111** | 0 | 0 |
| M1b | M1a + `commitBodyHook(ctx)` declaration and calls before head read and world insert; ≈ **3** | 0 | 0 |
| M1c | M1b + `appendIntentBodyHook(ctx)` declaration and mid-body call; ≈ **2** | 0 | 0 |
| M1d | M1c + test-only wrapping `database/sql` driver and source-context assertion; **0** | 0 | 0 |

The M1b and M1c scratch states included their new hook tests. M1d's targeted positive gate (`TestCommitBodyUsesCallerContext|TestCommitBodyHookPhases|TestAppendIntentCancelledMidBodyIsNotCommitted|TestDriverCommitErrorAfterRealCommitIsUncertainAndReconciles`) returned rc 0. I applied each of the four formerly surviving mutations to the scratch tree and restored the bytes after each run: each targeted test returned rc 1. The other 14 kills are the designer's §6.1 evidence, not remeasured here. No socket or `httptest` gate is treated as a verdict in this sandbox; any such result is **UNINFORMATIVE UNDER SANDBOX**. The full-suite `host/pkgproj/TestQueryInterfaceReturnsWhileADescendantHoldsStdout` load flake is outside this store slice and is not an acceptance gate.

Land both gated instruments in **M1a**: `host/store/durable_measure_test.go` and `host/daemon/startup_measure_test.go`. They skip unless `WORLD_BOUND_MEASURE=1`, assert no latency threshold, preserve the reproducible evidence for Mark's bound table, and introduce no runtime budget. The designer's §8 V7/V8/V19 measurements are prior evidence; this planner did **not** rerun three independent N=50 campaigns or measure proposed budget values. Run their explicit measurement commands only when refreshing the table, not as the ordinary M1 pass gate.

## 2. Sub-milestones in compile order

At **each** boundary, run the two compile-fence commands above. Run targeted commands with the same exported environment. `go test ./host/store -run '^TestProductionContextRoots$' -count=1` is the census gate. No `.ail` file changes.

### M1a — durable core, three context methods, and worker lifetime

**Files:** `host/store/durable.go`, `store.go`, `journal.go`, `durable_test.go`, `context_roots_test.go`, `durable_measure_test.go`; `host/daemon/startup_measure_test.go`.

**Take prototype hunks:** all of `durable.go`; `Store.workers sync.WaitGroup` and `Close` waiting before `db.Close` and `lock.release`; `Commit` wrapper plus `CommitContext`, `beginDurable`/`finishDurable`, `classifyCommit`, `notCommitted`, `UncertainError`/`IsUncertain`; `AppendIntent` and `GetReceipt` wrappers plus their `*Context` versions; `journalRowFor(ctx, …)` and `bindCommitIntentTx(ctx, …)`; `GetEffectReceipt`'s local compatibility root; the three census pins; all prototype durable tests including the two-arm log-index witness and Close writer-lock test; both gated instruments. Preserve the single-connection pool and buffered worker result channel. Correct the two comments identified in the introduction. Do not add a pre-Commit `ctx.Err()` check: §2.1 found it redundant and unkillable.

**Accept:** `go test ./host/store -run 'TestDurableOpsHonourHeldConnection|TestCommitCancelledMidBodyIsAllOrNothing|TestDurableStallAtCutoffIsUncertainAndReconciles|TestUncertainBareCommitReconcilesByLogIndexNotHead|TestCloseKeepsWriterLockUntilDurableWorkerSettles|TestProductionContextRoots' -count=1`. The held-connection test checks each of the three methods returns a definite context error within 100 ms + 250 ms, leaves no worker, and leaves state untouched; the cutoff test checks both writes return uncertain and reconcile. The log-index test has separate landed-then-B and rolled-back-then-B arms. The Close test checks refusal of a second writer while a post-COMMIT worker remains, then cleanup and re-Open after release. **Mutations:** A01–A14 in §3, one for each refusal or lifecycle branch.

### M1b — Commit body stall points

**Files:** `host/store/store.go`, `durable_test.go`.

**New code:** a production no-op `var commitBodyHook = func(context.Context) {}`; invoke it immediately before `selectedHeadTx(ctx, tx)` and immediately before the `worlds` `ExecContext`. In `TestCommitBodyHookPhases`, arm the hook on the second visit, cancel its exact caller ctx, yield 50 ms for `database/sql`'s asynchronous rollback, and assert a definite `context.Canceled`, no uncertainty, and an indeterminate receipt. The first visit pins the head-read seam; the second pins a body write. Keep it test-serial: hooks are package globals.

**Accept:** `go test ./host/store -run '^TestCommitBodyHookPhases$' -count=1`. The source-context check added in M1d makes the point-statement mutants observably red even when `database/sql` wins its asynchronous rollback race. **Mutations:** B01–B02 in §3; run them at M1d when that check exists.

### M1c — AppendIntent mid-body cancellation

**Files:** `host/store/journal.go`, `durable_test.go`.

**New code:** a production no-op `var appendIntentBodyHook = func(context.Context) {}`; call it after `nextJournalSeqTx(tx)` and before `insertJournalObjectTx(tx, object)` in `AppendIntentContext` only. `TestAppendIntentCancelledMidBodyIsNotCommitted` cancels there, yields 50 ms, requires `errors.Is(err, context.Canceled) && !IsUncertain(err)`, and reads a `ReceiptNotStarted` with no intent. The yield forces the `sql.ErrTxDone` normalization arm, matching the Commit prototype test.

**Accept:** `go test ./host/store -run '^TestAppendIntentCancelledMidBodyIsNotCommitted$' -count=1`. **Mutation:** C01 in §3.

### M1d — post-COMMIT driver fault and context-use pin

**Files:** new `host/store/durable_driver_fault_test.go`; no production change.

**New test-only code:** register a unique `database/sql` driver name once. Its `driver.Conn` wraps `modernc.org/sqlite.Driver.Open`, delegates `Prepare`, `Close`, `Begin`, and `BeginTx`, and wraps `driver.Tx`. One armed `Tx.Commit` calls the **real** `Commit` first, returns its real error if any, then returns a sentinel error. Open a real file store, seed its invocation intent, swap only its `*sql.DB` to the wrapping driver's single-connection DB for the same file while retaining its writer lock, arm the fault, call `CommitContext`, and assert `*UncertainError` with the sentinel in `Cause`. A subsequent `GetReceiptContext` must be resolved: **uncertain → reconciles committed**. The test has a live control: the real COMMIT was reached, so a pre-commit injected error cannot satisfy it. Ensure cleanup closes the replacement DB and releases the original lock.

`TestCommitBodyUsesCallerContext` pins the exact `selectedHeadTx(ctx, tx)` and world-insert `tx.ExecContext(ctx, …)` call forms. Run it with `TestCommitBodyHookPhases`: the runtime test proves cancellation/rollback and the source check detects a point call switched to `Background` or `Exec`. This narrow source check is intentional: after `BeginTx(ctx)` cancellation, `database/sql` can itself return `ErrTxDone` before the mutated point statement runs, which is why the prototype's two point-statement mutants survived the previous runtime tests.

**Accept:** `go test ./host/store -run 'TestDriverCommitErrorAfterRealCommitIsUncertainAndReconciles|TestCommitBodyUsesCallerContext|TestCommitBodyHookPhases|TestAppendIntentCancelledMidBodyIsNotCommitted' -count=1`, then the M1a durable selector again. **Mutations:** D01 plus B01–B02 in §3. In the assembled scratch tree, D01, B01, B02 and C01 each returned rc 1 under their named tests.

## 3. Runnable mutation ledger

Each row is a literal source edit to make **after** its test lands. Run the named test with `-count=1` (except A06, with `-timeout=5s`), observe red, and restore the file byte-identically. Edits are scoped to the specified function or statement so identical text elsewhere remains intact. The M1a rows were killed in the designer's §6.1 run; the four formerly surviving rows were killed in this planner's scratch run. A source-pin kill is called out as such.

| ID / criterion | File and exact old → new edit | Red test |
|---|---|---|
| A01 held connection, writes | `durable.go`, `beginDurable`: `s.db.BeginTx(ctx, nil)` → `s.db.Begin()` | `TestDurableOpsHonourHeldConnection/CommitContext` and `/AppendIntentContext` |
| A02 held connection, receipt | `journal.go`, `GetReceiptContext`: both `journalRowFor(ctx, s.db, id, …)` → `journalRowFor(context.Background(), s.db, id, …)` | `TestDurableOpsHonourHeldConnection/GetReceiptContext` |
| A03 cutoff returns at deadline | `durable.go`, `finishDurable`: replace the entire `select { case err := <-done: …; case <-ctx.Done(): … }` with `return classifyCommit(ctx, op, <-done)` | `TestDurableStallAtCutoffIsUncertainAndReconciles/CommitContext` |
| A04 uncertain cannot unwrap to ctx | `durable.go`: no `Unwrap` method → add `func (e *UncertainError) Unwrap() error { return e.Cause }` | `TestDurableStallAtCutoffIsUncertainAndReconciles/CommitContext` |
| A05 post-cutoff is uncertain | `durable.go`, `finishDurable` default arm: `return &UncertainError{Op: op, Cause: ctx.Err()}` → `return fmt.Errorf("store: %s not committed: %w", op, ctx.Err())` | `TestDurableStallAtCutoffIsUncertainAndReconciles/CommitContext` |
| A06 worker must exit | `durable.go`, `finishDurable`: `make(chan error, 1)` → `make(chan error)` | `TestDurableStallAtCutoffIsUncertainAndReconciles/CommitContext` (`-timeout=5s`, expected red or timeout dump) |
| A07 Commit uses cutoff | `store.go`, end of `CommitContext`: `return s.finishDurable(ctx, "commit", tx)` → `return tx.Commit()` | `TestDurableStallAtCutoffIsUncertainAndReconciles/CommitContext` |
| A08 AppendIntent uses cutoff | `journal.go`, `AppendIntentContext`: `if err := s.finishDurable(ctx, "append intent", tx); err != nil {` → `if err := tx.Commit(); err != nil {` | `TestDurableStallAtCutoffIsUncertainAndReconciles/AppendIntentContext` |
| A09 pre-cutoff Commit cause | `store.go`, `CommitContext`: `defer func() { err = notCommitted(ctx, "commit", err) }()` → `defer func() {}()` | `TestCommitCancelledMidBodyIsAllOrNothing` |
| A10 bare landed-then-B witness | `durable.go`, `finishDurable`: `err := tx.Commit(); durableAfterCommitHook()` → `durableAfterCommitHook(); err := tx.Commit()` (preserve the surrounding goroutine) | `TestUncertainBareCommitReconcilesByLogIndexNotHead/A_landed_then_B_on_top` |
| A11 bare post-cutoff branch | `durable.go`, `finishDurable` default arm: `return &UncertainError{Op: op, Cause: ctx.Err()}` → `return fmt.Errorf("store: %s not committed: %w", op, ctx.Err())` | `TestUncertainBareCommitReconcilesByLogIndexNotHead/A_landed_then_B_on_top` |
| A12 lock retained until worker settles | `store.go`, `Close`: delete `s.workers.Wait()` before `s.db.Close()` | `TestCloseKeepsWriterLockUntilDurableWorkerSettles` |
| A13 Close eventually finishes | `durable.go`, worker goroutine: `defer s.workers.Done()` → `defer func() {}()` | `TestCloseKeepsWriterLockUntilDurableWorkerSettles` |
| A14 wrapper census refusal | `store.go`, `Commit` wrapper: `s.CommitContext(context.Background(), c)` → `s.CommitContext(context.TODO(), c)` | `TestProductionContextRoots` |
| B01 head read uses caller ctx | `store.go`, `CommitContext`: `selectedHeadTx(ctx, tx)` → `selectedHeadTx(context.Background(), tx)` | `TestCommitBodyUsesCallerContext` (source pin), with `TestCommitBodyHookPhases` as runtime control |
| B02 world insert uses caller ctx | `store.go`, `CommitContext` Step 3: `tx.ExecContext(ctx,` → `tx.Exec(` (only the `INSERT OR IGNORE INTO worlds` call) | `TestCommitBodyUsesCallerContext` (source pin), with `TestCommitBodyHookPhases` as runtime control |
| C01 AppendIntent pre-cutoff cause | `journal.go`, `AppendIntentContext`: `defer func() { err = notCommitted(ctx, "append intent", err) }()` → `defer func() {}()` | `TestAppendIntentCancelledMidBodyIsNotCommitted` |
| D01 driver COMMIT error is uncertain | `durable.go`, `classifyCommit` default arm: `return &UncertainError{Op: op, Cause: err}` → `return fmt.Errorf("store: %s not committed: %w", op, err)` | `TestDriverCommitErrorAfterRealCommitIsUncertainAndReconciles` |

A05 and A11 deliberately mutate the same refusal branch for **different** scenarios: the invocation receipt and the bare log-index witness. The driver fault D01 is a different branch: `tx.Commit()` itself returned a driver error after real durability. A06's mutant can leave a worker parked while `Close` waits, so its timeout dump is the red verdict; never kill a process to clear it.

For A03, the complete old block in `finishDurable` is:

```go
select {
case err := <-done:
	return classifyCommit(ctx, op, err)
case <-ctx.Done():
	select {
	case err := <-done:
		return classifyCommit(ctx, op, err)
	default:
		return &UncertainError{Op: op, Cause: ctx.Err()}
	}
}
```

Replace that block with `return classifyCommit(ctx, op, <-done)`. For A10, replace the consecutive lines `err := tx.Commit()` then `durableAfterCommitHook()` with `durableAfterCommitHook()` then `err := tx.Commit()` inside the `finishDurable` worker. These exact replacements avoid touching the same names elsewhere.

## 4. Execution notes

The prototype's `durableAfterCommitHook` is a test-only seam implemented as a production no-op; it already covers a landed commit followed by a caller deadline. Keep all package-global hook tests serial and restore hooks in `defer`. Use a short watchdog on every stall so a ctx-ignoring mutant reports red instead of leaving the gate indefinitely. Re-run the four new mutations after transferring the scratch tests; the scratch transcript proves the design is runnable but is not a substitute for mutation evidence on the landed files. After M1, row 106 M5 can adopt the three context methods, subject to the same §9 merge gate. M2–M7 remain unplanned here.
