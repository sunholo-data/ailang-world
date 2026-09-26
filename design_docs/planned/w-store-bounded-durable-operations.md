# w-store-bounded-durable-operations — row 23's policy tranche: finite budgets, a cancellable `Commit`, and the strict guard last

**Status**: Planned — design + prototype (iteration 197, designer `claude:claude-opus-5-5`). **The bound table (§1) needs Mark's ratification before any milestone merges** (§9). The prototype proves M1's mechanism. It is uncommitted residue in the design worktree; the controller decides what lands.
**Item**: queue row 23 of `design_docs/world-mission.md` (`w-store-deadline-free-residue-owner`), **policy tranche**. The plumbing tranche landed at `b9ecabe` (iteration 194); its doc is `design_docs/implemented/w-store-deadline-free-residue-owner.md` (below, "the plumbing doc"). This tranche implements that doc's §5 answer A.
**Clauses**: clause-2 (bounded, honest daemon behaviour). It is on the clause-6 critical path because it is the named prerequisite of row 106's M5 (`R-106-11`, "Round 2 carve-out" in `design_docs/planned/w-transition-invocation-coordinator.md`).
**Ruling**: `D-WORLD-37` = A, attended by Mark 2026-09-26, verbatim: *"A: finite active operations. Keep the strict store-boundary guard (reject nil or deadline-free contexts) as the closing move; require finite budgets for daemon startup/bootstrap, approval active I/O, publish validation and credential lookup, with human think time between polls consuming no budget; `Commit` accepts a ctx and may be cancelled before the durable commit, all-or-nothing, and an uncertain outcome is reconciled, never blindly retried. The policy tranche routes as loop work under row 23; its design brings the concrete bound table to Mark before anything ships."* This doc does not re-open it.
**Estimate**: ~1.5d. Six milestones, each ≤150 production LOC (§5). **First landable slice = M1**: exactly what row 106's M5 needs, and it contains **no duration value**. The prototype is M1.
**Verified against**: dev `d8db0eb` (worktree `.design-wt-iter197`, detached). Go `go1.26.6 darwin/arm64`, `modernc.org/sqlite v1.54.0`, pinned interpreter `$HOME/.pinned-ailang/ailang` (v0.41.0). Line numbers cited as `file:N@d8db0eb` are from `git show d8db0eb:<file>`. Unsuffixed prototype line numbers refer to the worktree after the prototype. Every "the codebase does X" claim has a §8 row.

---

## 0. Problem, measured

The plumbing tranche threaded caller contexts through 11 store **reads**. It explicitly left four things (plumbing doc §4/§9):

1. **The durable tail takes no context.** Three store methods sit on row 106's durable tail: `Commit(c)` (`store.go:874@d8db0eb`), `AppendIntent(id, intent)` (`journal.go:411@d8db0eb`) and `GetReceipt(id)` (`journal.go:813@d8db0eb`). The coordinator calls exactly these three and no other context-free store method (V2). `Commit` and `AppendIntent` open `s.db.Begin()` (`store.go:909`, `journal.go:420@d8db0eb`). `GetReceipt` issues two `QueryRow`s on `s.db` (`journal.go:819,823@d8db0eb`). Commit's head compare runs under a literal `context.Background()` (`store.go:934@d8db0eb`) (V1).
2. **The wait for the one pooled connection is unbounded.** The pool is `SetMaxOpenConns(1)` (`store.go:305@d8db0eb`). `busy_timeout` = 2000 ms (`writer_lock.go:187@d8db0eb`) bounds SQLite **lock** retries only. A context-free `Begin()` waits behind every in-process holder of that connection (V1, V3).
3. **Four production roots are deadline-free** (plumbing doc §4 guard probe): `runApprove→MintAttendedApproval`, `runPublish→InvokeAttendedPublish→validate`, `runServe→Run→New→Bootstrap`, and `commit middleware→Resolve→ResolveSession`. The 5th durable caller is `/v1/commit` itself (`handlers.go:600@d8db0eb`, `d.store.Commit(commit)`, no ctx) (V2, V4, V5).
4. **Context-free writes on those root paths**: `PutObject` and `SetRegistryHead` (approve, registry bootstrap), `AppendClaimedEffectIntent`, `AppendNextEffectIntent` and `AppendEffectOutcome` (publish invoke), and `MintSession` (authority mint) (V6).

## 1. The bound table (for Mark's ratification)

**Margin rule, applied uniformly.** `B = max(10 × max_observed, F)`, rounded **up** to a whole second.
- `max_observed` is the largest sample over **three independent runs of N=50** real operations on a real file store in `t.TempDir()` on this machine (Apple silicon, macOS). The instruments are `host/store/durable_measure_test.go` and `host/daemon/startup_measure_test.go` (gated `WORLD_BOUND_MEASURE=1`). Raw output: V7, V8.
- `F = busy_timeout + 1 s = 3 s` is the floor for any operation that can take a SQLite lock. The measured stimulus for F is the predecessor's lock-blocked read, which ended at **~2.05 s** with `SQLITE_BUSY` against a 2000 ms `busy_timeout` (`writer_lock.go:176-186@d8db0eb` comment, row 111's measurement). A budget below F would turn one benign lock-retry window into a spurious timeout.
- The factor 10 covers machine variance (slow disk, CI runner, a loaded laptop). Every happy-path store op measured here is ≥ 2,800× below F, so **F, not the multiplier, decides every store-level row**. That is stated, not hidden: the store rows are a lock-window floor, not a latency fit.

| # | Operation / root | Measured stimulus (3 runs × N=50: p50 / p99 / **max**) | Proposed budget | Installed at (file / function) | On expiry |
|---|---|---|---|---|---|
| B1 | **Daemon startup/bootstrap**: `daemon.New` (store open + interpreter archive + `registry.Bootstrap`), `runServe→Run` | `New(fresh DB, archive)`: 841–844 ms / 849–864 ms / **873 ms**. `New(existing DB, archive)`: 58 ms / 64–67 ms / 71 ms (V8) | **9 s** (= ⌈10 × 0.873⌉) | `host/daemon/daemon.go` `Run`: `nctx, cancel := context.WithTimeout(ctx, startupBudget)` around `New(nctx, cfg)` **only**. The serve lifetime keeps the signal context | `*StartupError{Stage: StageRegistry, Err: context.DeadlineExceeded}` → `runServe` prints `ailang-worldd: <StartupError>` (proposed detail: `startup exceeded its 9s budget`) and returns `exitFatal` (the existing path, `cmd/ailang-worldd/main.go:207-209`). The writer lock is released by `d.abort` |
| B2 | **Approval active I/O**, per active call: `runApprove→MintAttendedApproval` (request → decide → poll) | Composite. Worst fixture = 9 reads (plumbing doc §4) + ≤ 6 writes (V6: request `PutObject`, decision `PutObject`, 2 × head `PutObject`+`SetRegistryHead`). Read max ≤ 0.088 ms, write max ≤ 1.07 ms (V7), so 10 × composite ≈ 50 ms, below F | **5 s** — **policy value** (see note) | `cmd/world-publish/main.go` `runApprove`, replacing the `context.Background()` at `:269`, **after** `requireAttendedOperator` (`:258`) | `broker: … context deadline exceeded` → `world-publish: mint approval: …`, exit `exitError` (1). The mint is all-or-nothing per write. A partial mint (request written, no decision) is the existing pending state, and the operator re-runs |
| B3 | **Publish validation**: `invoke→validatePublishApproval` (2 object reads) | `GetObject` max 0.023 ms (V7) | **3 s** (= F) | `host/broker/broker.go` `invoke`, `:245`: `vctx, cancel := context.WithTimeout(ctx, validationBudget)` around the validation call **only**. The network publish keeps `handlerExecTimeout` (30 s, `handlers.go:16`) | Validation error before dispatch → CLI `world-publish: publish: … context deadline exceeded`, exit 1, approval **not consumed** (validation precedes the claim) |
| B3′ | **Publish root** (derived, needed by the guard): `runPublish→InvokeAttendedPublish` | Derived: 30 s exec + B3 + B6 | **36 s** (= 30 + 3 + 3) | `cmd/world-publish/main.go` `runPublish`, replacing `context.Background()` at `:382`, after `requireAttendedOperator` (`:352`) | Same exit path. An expiry after the claim is the existing broker "indeterminate" outcome, reported and reconciled by `world-publish reconcile`, never re-published |
| B4 | **Credential lookup**: session middleware → `ResolveContext` → `ResolveSession` | `ResolveSession` max 0.029 ms (V7) | **3 s** (= F) | `host/daemon/middleware.go` `Wrap`, `:49`: `ResolveContext(WithTimeout(r.Context(), credentialBudget), …)` replaces `Resolve`. The Background root in `authority/resolver.go` `Resolve` is deleted | **HTTP 503**, class `Timeout` (the existing `writeReadTimeout` envelope). Today a lookup error collapses to **401 SessionUnknown**, which tells a valid client its token is bad. That changes |
| B5 | **Durable tail on `/a2a/`** (row 106 M5): `GetReceiptContext` + `AppendIntentContext` + `CommitContext` together | Max: `GetReceipt` 0.088 ms, `AppendIntent` 1.07 ms, `Commit` 0.94 ms (V7). Sum 2.1 ms | **3 s** for the three together (= F) | Installed by **row 106's M5**, not here: `tail := WithTimeout(ctx, durableTailBudget)`. Recommendation to M5, not absorbed: run the capsule under `InvokeWait − durableTailBudget` so a slow transition cannot starve the tail | Pre-cutoff → R11 "not committed". Post-cutoff → `*store.UncertainError` → R16 "not confirmed; resend the same task id" |
| B6 | **`/v1/commit`**: `handleCommit → CommitContext` | `Commit` max 0.94 ms (V7) | **3 s** (= F) | `host/daemon/handlers.go` `handleCommit`, `:600`: `CommitContext(WithTimeout(r.Context(), commitBudget), commit)` | Pre-cutoff: **503** `Timeout`, "commit not committed; safe to resend". Uncertain: **503** class `CommitUncertain`, "outcome not confirmed; read `/v1/head` before resending". Conflict stays **409** |
| B7 | **Pooled-connection acquisition** (inside every row above) | Uncontended `BeginTx`+`Rollback`: p99 4–7.6 µs, max 75 µs (V7). Contended: the prototype test holds the connection and the op returns at the 100 ms deadline, with margin (V10) | **No separate timer.** It consumes the enclosing operation's budget | `store.beginDurable` → `db.BeginTx(ctx, nil)` (database/sql waits on `ctx.Done()`, V3) | `store: begin <op>: context deadline exceeded`, i.e. definite not committed |
| B8 | **Non-cancellable durable step** (the driver `COMMIT` after the cutoff) | Whole `Commit` max 0.94 ms (V7). Worker side: SQLite lock ≤ `busy_timeout` 2 s, then fsync | **Caller side: enforced** = the op's remaining budget. **Worker side: not enforceable in Go**; one lock window (row 111), then fsync | `store.finishDurable` selects on `ctx.Done()` and returns `*UncertainError`. The COMMIT goroutine holds the sole connection, so **at most one** exists per store (goroutine-count tested, V10) | `*store.UncertainError`, which is never `errors.Is(…, context.DeadlineExceeded)` |

**Human think time consumes no budget, by construction.** Every CLI budget (B2, B3′) is created **after** `requireAttendedOperator` returns (`main.go:258` → `:269`; `:352` → `:382`, V4). The human's typing happens before the timer exists. **No production poll loop waits on a human**: `MintAttendedApproval` performs request → decide → poll in one call. The human attends the CLI and does not answer a remote poll (plumbing doc F2, V4). If a later row adds a human-waiting poll loop, each poll step gets its own `WithTimeout(parent, B2)` inside the loop, and the sleep between steps sits outside every timer.

**B2 is a policy value, not a fit.** The measured work is ~50 ms after the multiplier. The budget is human-facing, and `walkApprovalHead` is O(approval-head objects) (`approve.go:230`), which grows without a measured ceiling. **Recommended: 5 s.** Alternative: **15 s**, if Mark prefers "never time out an operator on a large approvals history" over "fail fast".

**Nothing in this table bounds** SQLite lock waits beyond `busy_timeout` (row 111), archive/replay subprocess descendants (rows 112/113), or the archive `--version` probe's own 10 s timer inside B1 (row 112). B1 bounds only the store work on the startup path that accepts `ctx`.

## 2. Cancellation cutoff and the commit-result contract

### 2.1 Where cancellation stops being honoured (prototype line level)

`CommitContext` (`store.go:884`) and `AppendIntentContext` (`journal.go:421`) have the same shape:

1. `s.beginDurable(ctx, op)` (`durable.go:44`) → `db.BeginTx(ctx, nil)`. The connection wait honours ctx (database/sql `DB.conn`: `select { case <-ctx.Done(): … case ret := <-req: }`, `sql.go:1369-1375`, V3).
2. Every statement in the body runs inside a transaction bound to ctx. When ctx ends, database/sql's `Tx.awaitDone` rolls the transaction back **asynchronously** (`sql.go:2209-2221`, V3). Ctx-bound statements (`QueryRowContext`/`ExecContext`) then fail with `ctx.Err()`. Unbound helper statements (`nextJournalSeqTx`, `insertJournalObjectTx`) fail with `sql.ErrTxDone`. `notCommitted` (`durable.go:104`) normalises that case to the ctx cause. This race was **found by the prototype**: 5/30 runs reached `ErrTxDone` (§6, MUT-NOTCOMMITTED-DROP).
3. **The cutoff** is `finishDurable` (`durable.go:62`, called at `store.go:1029` and `journal.go:460`), which calls `tx.Commit()`. Inside database/sql, `Tx.Commit` checks `tx.ctx.Done()` and returns `ctx.Err()` **before** `tx.done.CompareAndSwap(false, true)` (`sql.go:2291-2301`). The rollback path takes the same CAS (`sql.go:2327`), so **exactly one of commit and rollback wins**. **The CompareAndSwap at `sql.go:2299` is the cancellation cutoff.** Before it, cancellation rolls back. After it, the driver runs `exec(context.Background(), "commit")` (`modernc tx.go:36`, V3), which is not interruptible.
4. The caller waits for step 3 only until `ctx.Done()`. After that it gets `*UncertainError` while the goroutine finishes COMMIT (`durable.go:64-79`).

The prototype has **no separate `ctx.Err()` pre-check** before `tx.Commit()`. database/sql already performs that check (point 3). The mutation that removed a pre-check survived, so the code was deleted as redundant (§6).

### 2.2 All-or-nothing

Each durable op is **one SQLite transaction** on the default rollback journal (no `journal_mode` pragma, V9). Every pre-cutoff exit is a rollback: the deferred `tx.Rollback()` or `awaitDone`. The post-cutoff COMMIT is SQLite's atomic commit. A commit-bound `Commit` writes its `committed` journal outcome **in the same transaction** as the objects, world, log and head (`store.go` step 6), so **resolved receipt ⇔ commit landed**. The prototype test `TestCommitCancelledMidBodyIsAllOrNothing` cancels between the head write and the outcome write and asserts: head unchanged, next world absent, receipt still indeterminate (V10).

### 2.3 The three outcomes a caller sees

| Outcome | How the caller recognises it | Meaning | Caller action |
|---|---|---|---|
| **committed** | `err == nil` | Durable | Success. Never re-check ctx afterwards (row 106 MUT-POSTCOMMIT-CTX stays) |
| **not committed** | `err != nil && !store.IsUncertain(err)`. When caused by the caller's lifetime, `errors.Is(err, context.Canceled/DeadlineExceeded)` also holds | Rolled back; nothing landed | Safe to resend with the same id. `AppendIntent` is idempotent on identical bytes |
| **uncertain** | `store.IsUncertain(err)`. `*UncertainError` **does not unwrap to ctx errors** (MUT-UNCERTAIN-UNWRAPS kills that) | Committed or not, never partial | **Reconcile, never blindly retry**: invocation → `GetReceiptContext(id)` (resolved ⇒ committed; indeterminate ⇒ not committed, per row 106's R15 argument on the single-connection pool). Bare commit → `SelectedHead` == `NextWorld.Ref`. Because the sole connection is held until COMMIT settles, **the next acquisition observes the settled outcome** |

### 2.4 Is an error from `tx.Commit()` itself uncertain? Yes, by contract (conservative)

- `ctx.Err()` / `sql.ErrTxDone` from `Tx.Commit` are **definite not-committed**. database/sql returns them only when it did **not** win the CAS, i.e. the driver COMMIT never ran (`sql.go:2291-2301`). `classifyCommit` (`durable.go:85`) maps them to "not committed".
- A **driver** error is classified **uncertain**. The evidence goes both ways. On `SQLITE_BUSY`, modernc checks `sqlite3_get_autocommit` and forces a `rollback` (`modernc tx.go:35-53`, V3), and that case is really not-committed. But an I/O error during the journal delete/fsync phase can leave a hot journal, and SQLite resolves it at the next open. The prototype cannot inject that, and this design does not assume it. Calling every driver error uncertain is sound (the receipt read settles it), costs one read, and cannot turn a landed commit into "failed". **Unproven arm**: the prototype cannot reach a driver COMMIT error, so MUT-DRIVERERR-NOTCOMMITTED **survived**. M1 adds the fault-injection seam (§5, §6).

## 3. The strict guard — the closing move, landing LAST (M6)

`requireDeadline(ctx) error` rejects `ctx == nil` or `_, ok := ctx.Deadline(); !ok` with `store.ErrNoDeadline`. It is installed as the first statement of:

- the **8 reads** the plumbing doc named: `GetObject`, `GetWorld`, `GetLogEntry`, `GetRegistryHead`, `GetVerifyResult`, `SelectedHead`, `ReadObject`, `ResolveSession`;
- the **durable ops made ctx-accepting here**: `CommitContext`, `AppendIntentContext`, `GetReceiptContext` (M1), and the M4 writes: `PutObjectContext`, `SetRegistryHeadContext`, `AppendClaimedEffectIntentContext`, `AppendNextEffectIntentContext`, `AppendEffectOutcomeContext`, `GetEffectReceiptContext`;
- `RevokeSession` (already ctx-accepting; a write).

**Census that proves no production root reaches the store without a deadline.** The landed `TestProductionContextRoots` (`context_roots_test.go`) pins every `context.Background/TODO/WithoutCancel` reference by file|function. M6 flips it:
(a) every store-reaching root in `contextRootPins` is removed, or re-pinned **with the `WithTimeout` that wraps it in the same function**. A new AST rule requires each surviving `Background()` in a store-reaching file to be the direct first argument of `context.WithTimeout`/`WithDeadline`.
(b) The four chains of the plumbing doc's guard probe are re-run as tests **with the guard installed**, and all four pass. The same four with the budget removed fail with `ErrNoDeadline`: one mutation per root.
(c) A dynamic table test calls every guarded method with `nil`, `context.Background()` and a deadline ctx, and expects `ErrNoDeadline`, `ErrNoDeadline` and success.

`Commit`/`AppendIntent`/`GetReceipt`/`GetEffectReceipt` wrappers **no longer exist** by M6 (removed in M5). A Background-wrapping store method would make the guard trivially bypassable.

## 4. Why is this not a package? (S3)

This is Go host code **at the effect boundary** (S2): transaction lifetime, a connection pool and HTTP/CLI budgets. It has no AILANG semantics to publish. A package cannot own the store's own transaction. No `.ail` is touched, and the kernel does not grow.

## 5. Milestones (dependency order, each ≤150 production LOC)

| M | Scope | Prod LOC | Unblocks / tests |
|---|---|---|---|
| **M1** (prototype) | `durable.go` (`beginDurable`, `finishDurable`, `classifyCommit`, `notCommitted`, `UncertainError`). `CommitContext`, `AppendIntentContext`, `GetReceiptContext`. `journalRowFor(ctx, …)`. **Thin wrappers kept**: `Commit`, `AppendIntent` and `GetReceipt`, because 65 test call sites and 2 production sites (`handlers.go:600`, `broker/recover.go:157`) would otherwise land in this milestone (V2). The census gains 3 named rows (`journal.go\|AppendIntent`, `\|GetReceipt`, `\|GetEffectReceipt`), each marked "removed by M5". **Additions over the prototype**: a driver fault seam (a `sql.Register`ed wrapping driver used only by tests, whose `Tx.Commit` can error or stall **after** the real COMMIT) to kill MUT-DRIVERERR-NOTCOMMITTED and prove "uncertain → reconciles committed". A `commitBodyHook(ctx)` stall inside the transaction body to kill MUT-HEADREAD-BG/MUT-EXEC-NOCTX. A mid-body cancel for `AppendIntentContext` to kill MUT-APPEND-NOTCOMMITTED-DROP | ~130 (prototype: 63 code lines in `durable.go` + 61 changed lines in store/journal, V11) | **Row 106 M5** (it adopts these three methods per R-106-11). No duration value in M1 |
| **M2** | `/v1/commit` → `CommitContext` with B6. 503 `Timeout` / 503 `CommitUncertain` mapping. `docs/QUICKSTART.md` gains the uncertain-commit reconcile step (S7) | ~45 | AC5 |
| **M3** | Credential lookup B4: middleware → `ResolveContext` with timeout, 503 on lookup expiry. Delete the `authority/resolver.go` `Resolve` Background root | ~40 | AC6 |
| **M4** | Ctx-accepting writes on the root paths (V6): `PutObjectContext`, `SetRegistryHeadContext`, the three effect-journal writes, `GetEffectReceiptContext`. Thread them through `broker/approve.go`, `broker/broker.go`, `registry/registry.go` | ~140 | AC7 |
| **M5** | Root budgets B1, B2, B3, B3′ at their install sites. **Remove the wrappers** `Commit`/`AppendIntent`/`GetReceipt`/`GetEffectReceipt` and migrate the 2 production and ~65 test call sites (the test churn is mechanical) | ~90 | AC8–AC10 |
| **M6** | Strict guard (§3), census flip, four-chain probe as tests | ~60 | AC11–AC13. **Lands last** |

Row 106's M5 needs **M1 only**. M2–M6 complete `D-WORLD-37`.

## 6. Acceptance criteria and non-vacuity

### 6.1 Executed on the prototype (14 mutations: 10 killed, 4 survived)

Harness: `~/.ailang/state/world-iter197/designer-scratch/mutate.py`. It applies each mutation, runs `go test ./host/store/ -run <selector> -count=5`, restores the file in `finally`, and writes transcripts to `…/mut/<NAME>.txt` (V12). Tests have a watchdog (`boundedCall`), so a ctx-ignoring mutant **fails** instead of hanging.

| AC | Mutation | Killing assertion (observed red) | Result |
|---|---|---|---|
| AC1 durable ops honour a held sole connection | MUT-BEGIN-NOCTX: `BeginTx(ctx,nil)` → `Begin()` | `TestDurableOpsHonourHeldConnection`: `durable_test.go:87: still blocked 1.10s after a 100ms deadline` | **KILLED** |
| AC1 (receipt arm) | MUT-RECEIPT-NOCTX: `GetReceiptContext` reads under `Background()` | `…/GetReceiptContext`: `still blocked 1.10s after a 100ms deadline (err after unblock: <nil>)` | **KILLED** |
| AC2 caller returns at deadline past the cutoff | MUT-CUTOFF-NOSELECT: wait on `<-done` only | `TestDurableStallAtCutoffIsUncertainAndReconciles`: `still blocked 1.10s after a 100ms deadline` | **KILLED** |
| AC3 uncertain ≠ ctx error | MUT-UNCERTAIN-UNWRAPS: add `Unwrap() → Cause` | `…:181: err = store: commit outcome is uncertain …; want *UncertainError that is not a ctx error` | **KILLED** |
| AC3 (other direction) | MUT-UNCERTAIN-AS-NOTCOMMITTED: post-cutoff expiry reported as not committed | `…:181: err = store: commit not committed: context deadline exceeded; want *UncertainError …` | **KILLED** |
| AC4 no accumulating blocked worker | MUT-LEAK-UNBUFFERED: unbuffered `done` channel | `…:192: goroutines 5 > baseline 4 after release` | **KILLED** |
| AC2 (Commit wiring) | MUT-COMMIT-DROP-CUTOFF: `return tx.Commit()` | `…:181: err = <nil>; want *UncertainError …` | **KILLED** |
| AC2 (AppendIntent wiring) | MUT-APPEND-DROP-CUTOFF | `…/AppendIntentContext:181: err = <nil>; want *UncertainError …` | **KILLED** |
| AC5′ pre-cutoff cancel reports its cause, all-or-nothing | MUT-NOTCOMMITTED-DROP (Commit) | `TestCommitCancelledMidBodyIsAllOrNothing:132: err = store: next journal seq: sql: transaction has already been committed or rolled back; want definite Canceled` | **KILLED** (after the test was made deterministic, see below) |
| census pins wrapper roots | MUT-WRAPPER-TODO: wrapper uses `context.TODO()` | `TestProductionContextRoots: context_roots_test.go:118: root references: Background=18 TODO=1` | **KILLED** |
| (in-body ctx) | MUT-HEADREAD-BG: head compare under `Background()` | — | **SURVIVED**: no in-body stall seam. M1 obligation (`commitBodyHook`) |
| (in-body ctx) | MUT-EXEC-NOCTX: world insert `ExecContext` → `Exec` | — | **SURVIVED**: same seam. Point statements cannot be stalled without it |
| AC3 driver arm | MUT-DRIVERERR-NOTCOMMITTED: driver COMMIT error → "not committed" | — | **SURVIVED**: no driver fault seam. M1 obligation (wrapping test driver) |
| AC5′ (AppendIntent) | MUT-APPEND-NOTCOMMITTED-DROP | — | **SURVIVED**: no mid-body hook in `AppendIntentContext`. M1 obligation |

Two findings came from running the mutations, not from design. (1) MUT-NOTCOMMITTED-DROP first **survived at 25/30**: `awaitDone`'s rollback is asynchronous, so the `ErrTxDone` arm is reached only when the rollback wins the race. The test hook now cancels and then yields 50 ms, which makes the arm deterministic (killed 5/5 runs). (2) An explicit `ctx.Err()` pre-check before `tx.Commit()` survived its removal, because database/sql performs it (§2.1). It was deleted rather than kept as untestable code (S6).

### 6.2 Acceptance for M2–M6 (mutation named per AC; refusal branches get one mutation each)

| AC | Test (to be written in its milestone) | Mutation that must turn it red |
|---|---|---|
| AC5 `/v1/commit` bounded | Hold the connection. POST `/v1/commit` answers **503 Timeout** within B6 + ε, and the head is unchanged. Stall at the cutoff: **503 CommitUncertain**, then `/v1/head` reconciles | MUT-M2-NOCTX (`Commit` wrapper kept at `:600`); MUT-M2-UNCERTAIN-AS-TIMEOUT (both branches return one class) |
| AC6 lookup expiry is not a denial | Hold the connection. A valid token gets **503**, not 401 | MUT-M3-COLLAPSE (lookup error → `DenialUnknown`); MUT-M3-NOCTX (keep `Resolve`) |
| AC7 root-path writes cancellable | Per write: held connection → ctx error within budget + ε | One MUT-M4-NOCTX-<method> per converted write (6) |
| AC8 startup bounded | `Run` with a held connection / stalled bootstrap write → `StartupError` wrapping `DeadlineExceeded` within B1 + ε, and the writer lock is released (re-Open succeeds) | MUT-M5-STARTUP-SIGCTX (pass the signal ctx to `New`); MUT-M5-SERVE-BOUNDED (the budget wraps the serve lifetime, so a served request fails after 9 s) |
| AC9 approval budget excludes think time | Fake attended prompt blocks 2×B2, then mint succeeds | MUT-M5-TIMER-BEFORE-PROMPT (timer created before `requireAttendedOperator`) |
| AC10 validation-only budget | A stalled validation read errors within B3, and the approval is not consumed. A handler exec taking > B3 still succeeds | MUT-M5-VALIDATION-UNBOUNDED; MUT-M5-BUDGET-WRAPS-DISPATCH |
| AC11 guard rejects deadline-free | Table: every guarded method × {nil, Background, deadline} | One MUT-M6-UNGUARDED-<method> per method (18); MUT-M6-NIL-ONLY (checks nil but not the deadline) |
| AC12 census flip | `TestProductionContextRoots` with the WithTimeout-adjacency rule | MUT-M6-BARE-BACKGROUND (a bare `Background()` passed to a store call in `runApprove`) |
| AC13 four chains pass under the guard | The four plumbing-doc probe chains as tests | Remove each root's budget (4 mutations) → `ErrNoDeadline` |

**Row-106 M5-style test (required, already in the prototype):** `TestDurableOpsHonourHeldConnection` holds the sole pooled connection (`s.db.Conn`) past a 100 ms caller deadline. For each of `CommitContext`/`AppendIntentContext`/`GetReceiptContext` it asserts `errors.Is(err, context.DeadlineExceeded) && !IsUncertain(err)`, `elapsed ≤ 100 ms + 250 ms`, and `runtime.NumGoroutine()` back at baseline. It then asserts the store is untouched and that a released store commits (live control). `TestDurableStallAtCutoffIsUncertainAndReconciles` injects the stall at the cutoff for both writes (V10).

## 7. Conflict Surface

| File (this tranche) | M | Other open rows touching it |
|---|---|---|
| `host/store/durable.go` (new), `store.go`, `journal.go`, `context_roots_test.go` | M1, M4, M6 | **111** arm A (per-request busy_timeout on store reads: same files, and it must keep pool-state restoration). **26** (producer source reads and envelope write: it must adopt `*Context` writes and B-budgets, not a hidden root). **106** M5 (adopts M1's methods; interface in `host/coordinator/coordinator.go:27-34`) |
| `host/daemon/handlers.go`, `middleware.go`, `daemon.go` | M2, M3, M5 | **111** (daemon read-path status mapping: coordinate the 503 envelopes). **106** M5 (daemon constructs the coordinator) |
| `host/authority/resolver.go` | M3 | none measured |
| `host/broker/approve.go`, `broker.go`, `publish_op.go`; `host/registry/registry.go` | M4, M5 | **25** (capsule blocked-child coverage): no file overlap. **112** touches `replay.go`/archive, not these |
| `cmd/world-publish/main.go`, `cmd/ailang-worldd/main.go` | M5 | none measured |

**This tranche does NOT**: bound SQLite lock waits beyond `busy_timeout` (row 111); bound archive/replay/pkgproj subprocess descendants (rows 112/113) or the archive probe inside startup; add blocked-child kill coverage (row 25); design row 26's producer; or widen the pool. If the pool ever widens, row 106's R15 inference and §2.3's "next acquisition observes the settled outcome" both fail, and that is a residual (§10).

## 8. Verification Log

All commands were run in `/Users/voightkampff/dev/sunholo-data/.design-wt-iter197`. `AILANG_BIN=$HOME/.pinned-ailang/ailang` exported. Every negative result is paired with a positive control in the same command.

| # | Claim | Command | Observed |
|---|---|---|---|
| V1 | Commit/Begin/Background/pool lines at HEAD | `git show d8db0eb:host/store/store.go \| grep -n 'SetMaxOpenConns\|func (s \*Store) Commit(\|s.db.Begin()\|selectedHeadTx(context.Background'` | `305: db.SetMaxOpenConns(1)`, `874: func (s *Store) Commit(c Commit) error {`, `909: tx, err := s.db.Begin()`, `934: … selectedHeadTx(context.Background(), tx)` |
| V1b | journal lines at HEAD | `git show d8db0eb:host/store/journal.go \| grep -n 'func (s \*Store) AppendIntent(\|func (s \*Store) GetReceipt(\|journalRowFor(s.db, id'` | `411`, `813`, `819/823` (GetReceipt), `858/862` (GetEffectReceipt) |
| V2 | Production callers of the three durable-tail methods (positive: the coordinator's 3 are found) | `grep -rn '\.Commit(\|\.AppendIntent(\|\.GetReceipt(\|\.AppendOutcome(' --include='*.go' host cmd \| grep -v _test.go \| grep -v 'tx\.Commit()'` | `coordinator.go:208 GetReceipt`, `:282 AppendIntent`, `:285 Commit`, `broker/recover.go:157 GetReceipt`, `daemon/handlers.go:600 Commit`. No production `AppendOutcome` (the negative sits in the same command as 5 positives). Test call sites: ~65 across store/daemon/broker/coordinator (per-file counts from the same grep) |
| V3 | database/sql and modernc behaviour (Go source, not memory) | `grep -n` on `$(go env GOROOT)/src/database/sql/sql.go` and `$(go env GOMODCACHE)/modernc.org/sqlite@v1.54.0/{tx.go,conn.go}` | `sql.go:1369-1375` conn wait `select … case <-ctx.Done()`; `2209-2221` `awaitDone` → `tx.rollback(discardConnection)`; `2287-2301` `Tx.Commit`: ctx check, then `CompareAndSwap(false,true)`; `2327` rollback CAS; `1901-1905` `keepConnOnRollback = hasSessionResetter && hasConnectionValidator`; modernc `conn.go:940,950` implements both (the connection is kept, so `:memory:` survives); `tx.go:35-53` `exec(context.Background(),"commit")`, then forced rollback if still in a transaction |
| V4 | Human prompt precedes the CLI roots | `grep -n 'MintAttendedApproval(context.Background\|InvokeAttendedPublish(context.Background\|requireAttendedOperator(in' cmd/world-publish/main.go` | `258 requireAttendedOperator` < `269 MintAttendedApproval(context.Background()…`; `352` < `382 InvokeAttendedPublish(context.Background()…` |
| V5 | Daemon roots and constants | `grep -n 'signal.NotifyContext\|Run(ctx' cmd/ailang-worldd/main.go`; `grep -n 'readDeadline = \|writeTimeout = ' host/daemon/daemon.go`; `grep -n 'd.store.Commit' host/daemon/handlers.go`; `grep -n 'm.resolver.Resolve' host/daemon/middleware.go` | `203 signal.NotifyContext(context.Background(),…)`, `207 daemon.Run(ctx,…)`; `writeTimeout = 30s` (:94), `readDeadline = 10s` (:138); `handlers.go:600`; `middleware.go:49 out := m.resolver.Resolve(header, …)` |
| V6 | Context-free writes on the root paths | `grep -n '\.PutObject(\|\.SetRegistryHead(\|\.AppendClaimedEffectIntent(\|\.AppendNextEffectIntent(\|\.AppendEffectOutcome(\|\.MintSession(\|\.CompareAndSetRegistryHead(' host/broker/{approve,publish_op,broker}.go host/registry/registry.go host/daemon/daemon.go host/authority/*.go` | approve.go:115,186,208 PutObject; :211 SetRegistryHead; broker.go:217,301,358 PutObject; :249 AppendClaimedEffectIntent; :252 AppendNextEffectIntent; :291,:313 AppendEffectOutcome; registry.go:165 PutObject, :168 SetRegistryHead; authority/mint.go:54 MintSession. `CompareAndSetRegistryHead`: no match on these paths (negative, same command) |
| V7 | Store op timings | `for r in 1 2 3; do WORLD_BOUND_MEASURE=1 go test ./host/store/ -run TestMeasureDurableOps -count=1 -v \| grep 'n=50'; done` (banked `…/designer-scratch/measure_store.txt`) | Max over 3 runs: AppendIntentContext **1.073 ms**; CommitContext **0.939 ms** (p99 ≤ 0.607 ms); GetReceiptContext absent 0.088 / resolved 0.072 ms; GetObject 0.023 ms; SelectedHead 0.097 ms; ResolveSession 0.029 ms; BeginTx+Rollback 0.075 ms (p99 ≤ 7.6 µs); Open fresh 5.23 ms; Open existing 0.584 ms |
| V8 | Startup timings | `for r in 1 2; do WORLD_BOUND_MEASURE=1 go test ./host/daemon/ -run TestMeasureStartup -count=1 -v \| grep 'n=50'; done` (banked `…/measure_startup.txt`) | `New(fresh, archive)` p50 841/844 ms, p99 849/864 ms, **max 851/873 ms**; `New(existing, archive)` p50 58 ms, max 71 ms |
| V9 | No `journal_mode` pragma (positive: `busy_timeout` found) | `git grep -n 'journal_mode\|busy_timeout' d8db0eb -- 'host/store/*.go' ':!*_test.go'` | Only `busy_timeout` lines (`writer_lock.go:189,200,204,208,210`). Zero `journal_mode`, so the default rollback journal |
| V10 | Prototype tests green and stable | `go test ./host/store -run 'Durable\|MidBody' -count=20` | `ok … 10.768s` (before the stall test was parametrised); after: `-count=10` `ok … 6.749s` |
| V11 | Prototype size | `git diff --numstat`; `grep -vc '^\s*//\|^\s*$' host/store/durable.go`; changed non-comment lines in store/journal diff | store.go +25/−19; journal.go +36/−14; context_roots_test.go +19/−15 (gofmt realigned the map); durable.go 109 lines / 63 code; 68 changed code lines in store+journal |
| V12 | Mutation run | `python3 ~/.ailang/state/world-iter197/designer-scratch/mutate.py` | 14 rows as in §6.1 (10 KILLED, 4 SURVIVED); transcripts in `…/mut/`. A 15th mutation (remove an explicit pre-cutoff `ctx.Err()` check) survived and the check was deleted (§2.1) |
| V13 | Census catches the new wrappers | first `go test ./host/store/` after adding wrappers, before the pin update | `TestProductionContextRoots: … Background=20` with `journal.go\|AppendIntent`, `\|GetReceipt`, `\|GetEffectReceipt:2` unexpected (reduced to 1 by a local `ctx`) |
| V14 | Residual owners OPEN | `grep -n '^25\. \|^26\. \|^111\. \|^112\. \|^113\. ' design_docs/world-mission.md` | 25 `w-capsule-blocked-child-kill-coverage`, 26 (gate met, open), 111 `w-daemon-lock-blocked-read-status-and-deadline-enforcement`, 112 `w-archive-replay-grandchild-bound`, 113 `w-pkgproj-quality-descendant-leak`, all present as queue entries |
| V15 | Full suite with the prototype | `AILANG_BIN=… go test ./... ` (banked `…/fullsuite.txt`), then the failing package alone, then `go list -deps ./host/pkgproj \| grep -E '^os/exec$\|ailang-world'` | **23 ok, 1 FAIL**: `host/pkgproj` `TestQueryInterfaceReturnsWhileADescendantHoldsStdout` (`iface_test.go:232: instrument: descendant pid not recorded`). Re-run alone: `ok host/pkgproj 11.425s`. Its deps list `os/exec` (positive) and `host/childenv`, and **not** `host/store`, so the prototype cannot reach it. Row 113's area; attributed as an instrument flake |

## 9. Decision for Mark — one word

**Ratify the bound table in §1 as proposed?**

| Answer | Meaning |
|---|---|
| **A** | Ratify B1–B8 as written: startup 9 s, approval active call 5 s, validation 3 s, publish root 36 s, credential lookup 3 s, durable tail 3 s, `/v1/commit` 3 s, acquisition inside the enclosing budget, and uncertain-on-expiry after the cutoff. Also ratify the margin rule `max(10 × max_observed, busy_timeout + 1 s)` and the three-outcome contract of §2.3. M1 routes first, and row 106 M5 unblocks. |
| **B** | Name the alternative, e.g. "B: approval 15 s" or "B: floor = 5 s". The designer re-derives the table and returns it. |

**Default while unanswered: nothing in this tranche merges** (the ruling's "before anything ships"). Row 106 M5 stays blocked.

## 10. Residuals (named OPEN owners)

| Residual | Owner |
|---|---|
| Lock-wait bound beyond `busy_timeout` (a post-cutoff COMMIT can still wait one 2 s lock window; worker side of B8) | **Row 111** |
| Archive `--version` probe and replay grandchildren inside startup/replay | **Row 112** |
| pkgproj descendants | **Row 113** |
| Blocked-child kill coverage | **Row 25** |
| Producer source-read/envelope-write budgets | **Row 26** (adopts the `*Context` writes and B-rule) |
| Tail reservation (`InvokeWait − durableTailBudget` for the capsule); coordinator adoption of M1 | **Row 106** M5 (recommendation, not absorbed) |
| Pool widening invalidates "next acquisition settles the outcome" and R15 | **Row 106** R-106-11 already owns R15. **No row owns a pool-widening change today**; none is planned, so no new row |
| Remaining context-free store writes **off** the four root paths (`MintSession`, `PutVerifyResult`, `PutWorld`, `SelectHead`, scans) | **Row 23** (this tranche, M4 inventory note). If Mark wants them bounded before the guard lands: **NEW ROW NEEDED**. M6's guard does not cover methods that take no ctx, so they are reachable only through the named CLI roots in the census |

## 11. Prototype manifest and results

**Files changed** (uncommitted, worktree `.design-wt-iter197`):
- `host/store/durable.go`: **new**. `UncertainError`, `IsUncertain`, `durableCommitHook`, `beginDurable`, `finishDurable`, `classifyCommit`, `notCommitted`.
- `host/store/store.go`: `Commit` → wrapper; body → `CommitContext(ctx, c)` using `beginDurable`, `selectedHeadTx(ctx, …)`, `ExecContext`, `finishDurable`; `bindCommitIntentTx(ctx, …)`.
- `host/store/journal.go`: `AppendIntent`/`GetReceipt` → wrappers; `AppendIntentContext`, `GetReceiptContext`; `journalRowFor(ctx, q, …)`; `GetEffectReceipt` keeps one local compatibility root.
- `host/store/context_roots_test.go`: 3 new pins (removed by M5).
- `host/store/durable_test.go`: **new**. `TestDurableOpsHonourHeldConnection`, `TestCommitCancelledMidBodyIsAllOrNothing`, `TestDurableStallAtCutoffIsUncertainAndReconciles` (2 subtests).
- `host/store/durable_measure_test.go`, `host/daemon/startup_measure_test.go`: **new**. Gated measurement instruments (skip unless `WORLD_BOUND_MEASURE=1`).

No `.ail` changed, so no `verify_ail.sh` run is needed.

**Results**: `go vet ./host/... ./cmd/...` clean. `go test -run '^$' ./...` compiles. Targeted `go test ./host/store -run 'Durable|MidBody' -count=10` ok. Full suite: **23 packages ok, 1 FAIL** (`host/pkgproj`, an unrelated instrument flake that passes alone, V15). Mutations: **10 killed, 4 survived**, and each survivor is a named M1 obligation (§6.1).
