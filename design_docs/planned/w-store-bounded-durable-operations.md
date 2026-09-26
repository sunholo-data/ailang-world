# w-store-bounded-durable-operations — row 23's policy tranche: finite budgets, a cancellable `Commit`, and the strict guard last

**Status**: Planned — **REVISION 1** after quorum round 1 (BLOCKED 2/2; see §12 "Quorum log"). Design + prototype (iteration 197, designer `claude:claude-opus-5-5`). **The bound table (§1) needs Mark's ratification before any milestone merges** (§9). The prototype proves M1's mechanism. It is uncommitted residue in the design worktree; the controller decides what lands.
**Item**: queue row 23 of `design_docs/world-mission.md` (`w-store-deadline-free-residue-owner`), **policy tranche**. The plumbing tranche landed at `b9ecabe` (iteration 194); its doc is `design_docs/implemented/w-store-deadline-free-residue-owner.md` (below, "the plumbing doc"). This tranche implements that doc's §5 answer A.
**Clauses**: clause-2 (bounded, honest daemon behaviour). It is on the clause-6 critical path because it is the named prerequisite of row 106's M5 (`R-106-11`, "Round 2 carve-out" in `design_docs/planned/w-transition-invocation-coordinator.md`).
**Ruling**: `D-WORLD-37` = A, attended by Mark 2026-09-26, verbatim: *"A: finite active operations. Keep the strict store-boundary guard (reject nil or deadline-free contexts) as the closing move; require finite budgets for daemon startup/bootstrap, approval active I/O, publish validation and credential lookup, with human think time between polls consuming no budget; `Commit` accepts a ctx and may be cancelled before the durable commit, all-or-nothing, and an uncertain outcome is reconciled, never blindly retried. The policy tranche routes as loop work under row 23; its design brings the concrete bound table to Mark before anything ships."* This doc does not re-open it.
**Estimate**: ~2.5d (revision 1: whole-surface migration). Nine milestones, each ≤150 production LOC (§5). **First landable slice = M1**: exactly what row 106's M5 needs, and it contains **no duration value**. The prototype is M1, plus one test-only hook added in revision 1 (§11).
**Verified against**: dev `d8db0eb` (worktree `.design-wt-iter197`, detached). Go `go1.26.6 darwin/arm64`, `modernc.org/sqlite v1.54.0`, pinned interpreter `$HOME/.pinned-ailang/ailang` (v0.41.0). Line numbers cited as `file:N@d8db0eb` are from `git show d8db0eb:<file>`. Unsuffixed prototype line numbers refer to the worktree after the prototype. Every "the codebase does X" claim has a §8 row.

---

## 0. Problem, measured

The plumbing tranche threaded caller contexts through 11 store **reads**. It explicitly left four things (plumbing doc §4/§9):

1. **The durable tail takes no context.** Three store methods sit on row 106's durable tail: `Commit(c)` (`store.go:874@d8db0eb`), `AppendIntent(id, intent)` (`journal.go:411@d8db0eb`) and `GetReceipt(id)` (`journal.go:813@d8db0eb`). The coordinator calls exactly these three and no other context-free store method (V2). `Commit` and `AppendIntent` open `s.db.Begin()` (`store.go:909`, `journal.go:420@d8db0eb`). `GetReceipt` issues two `QueryRow`s on `s.db` (`journal.go:819,823@d8db0eb`). Commit's head compare runs under a literal `context.Background()` (`store.go:934@d8db0eb`) (V1).
2. **The wait for the one pooled connection is unbounded.** The pool is `SetMaxOpenConns(1)` (`store.go:305@d8db0eb`). `busy_timeout` = 2000 ms (`writer_lock.go:187@d8db0eb`) bounds SQLite **lock** retries only. A context-free `Begin()` waits behind every in-process holder of that connection (V1, V3).
3. **Four production roots are deadline-free** (plumbing doc §4 guard probe): `runApprove→MintAttendedApproval`, `runPublish→InvokeAttendedPublish→validate`, `runServe→Run→New→Bootstrap`, and `commit middleware→Resolve→ResolveSession`. The 5th durable caller is `/v1/commit` itself (`handlers.go:600@d8db0eb`, `d.store.Commit(commit)`, no ctx) (V2, V4, V5).
4. **The whole exported Store surface outside those reads is context-free** (revision 1, objection 2). At `d8db0eb` there are **19 exported context-free I/O methods** on `*Store`. **16** have call sites outside `host/store`, and they hang off **seven** production roots: the four above plus `runSessionMint`, `runReconcile` and `runTransitions`. **3** have no production call site at all (`PutWorld`, `SelectHead`, `AppendOutcome`). §1.1 inventories every one, with its root and budget row (V6, V16–V18). With a pool of one, any context-free method left on the boundary makes the guard bypassable, and its own wait for the connection is unbounded. A *bounded* caller behind it still times out through `BeginTx(ctx)`, so it is not "permanent starvation", but that caller loses its whole budget to an unbounded holder. **This tranche therefore migrates the entire surface** (M5a/M5b) and guards all of it (M7).

## 1. The bound table (for Mark's ratification)

**Margin rule, applied uniformly.** `B = max(10 × max_observed, F)`, rounded **up** to a whole second.
- `max_observed` is the largest sample over **three independent runs of N=50** real operations on a real file store in `t.TempDir()` on this machine (Apple silicon, macOS). The instruments are `host/store/durable_measure_test.go` and `host/daemon/startup_measure_test.go` (gated `WORLD_BOUND_MEASURE=1`). Raw output: V7, V8.
- `F = busy_timeout + 1 s = 3 s` is the floor for any operation that can take a SQLite lock. The measured stimulus for F is the predecessor's lock-blocked read, which ended at **~2.05 s** with `SQLITE_BUSY` against a 2000 ms `busy_timeout` (`writer_lock.go:176-186@d8db0eb` comment, row 111's measurement). A budget below F would turn one benign lock-retry window into a spurious timeout.
- The factor 10 covers machine variance (slow disk, CI runner, a loaded laptop). Every happy-path store op measured here is ≥ 2,800× below F, so **F, not the multiplier, decides every store-level row**. That is stated, not hidden: the store rows are a lock-window floor, not a latency fit.

| # | Operation / root | Measured stimulus (3 runs × N=50: p50 / p99 / **max**) | Proposed budget | Installed at (file / function) | On expiry |
|---|---|---|---|---|---|
| B1 | **Daemon startup/bootstrap**: `daemon.New` (store open + interpreter archive + `registry.Bootstrap` + the startup integrity scans `ScanUnreadableLog`/`ScanUnreadableWorlds`), `runServe→Run` | `New(fresh DB, archive)`: 841–844 ms / 849–864 ms / **873 ms**. `New(existing DB, archive)`: 58 ms / 64–67 ms / 71 ms (V8) | **9 s** (= ⌈10 × 0.873⌉). The measured `New` includes the integrity scan (`daemon.go:517`). One 64-row scan page takes ≤ 0.2 ms (V19). The scan keeps its own 2 s page-boundary budget (`daemon.go:168,558`) inside B1 | `host/daemon/daemon.go` `Run`: `nctx, cancel := context.WithTimeout(ctx, startupBudget)` around `New(nctx, cfg)` **only**. The serve lifetime keeps the signal context | `*StartupError{Stage: StageRegistry, Err: context.DeadlineExceeded}` → `runServe` prints `ailang-worldd: <StartupError>` (proposed detail: `startup exceeded its 9s budget`) and returns `exitFatal` (the existing path, `cmd/ailang-worldd/main.go:207-209`). The writer lock is released by `d.abort` |
| B2 | **Approval active I/O**, per active call: `runApprove→MintAttendedApproval` (request → decide → poll) | Composite. Worst fixture = 9 reads (plumbing doc §4) + ≤ 6 writes (V6: request `PutObject`, decision `PutObject`, 2 × head `PutObject`+`SetRegistryHead`). Read max ≤ 0.088 ms, write max ≤ 1.07 ms (V7), so 10 × composite ≈ 50 ms, below F | **5 s** — **policy value** (see note) | `cmd/world-publish/main.go` `runApprove`, replacing the `context.Background()` at `:269`, **after** `requireAttendedOperator` (`:258`) | `broker: … context deadline exceeded` → `world-publish: mint approval: …`, exit `exitError` (1). The mint is all-or-nothing per write. A partial mint (request written, no decision) is the existing pending state, and the operator re-runs |
| B3 | **Publish validation**: `invoke→validatePublishApproval` (2 object reads) | `GetObject` max 0.023 ms (V7) | **3 s** (= F) | `host/broker/broker.go` `invoke`, `:245`: `vctx, cancel := context.WithTimeout(ctx, validationBudget)` around the validation call **only**. The network publish keeps `handlerExecTimeout` (30 s, `handlers.go:16`) | Validation error before dispatch → CLI `world-publish: publish: … context deadline exceeded`, exit 1, approval **not consumed** (validation precedes the claim) |
| B3′ | **Publish root** (derived, needed by the guard): `runPublish→InvokeAttendedPublish` | Derived: 30 s exec + B3 + B6 | **36 s** (= 30 + 3 + 3) | `cmd/world-publish/main.go` `runPublish`, replacing `context.Background()` at `:382`, after `requireAttendedOperator` (`:352`) | Same exit path. An expiry after the claim is the existing broker "indeterminate" outcome, reported and reconciled by `world-publish reconcile`, never re-published |
| B4 | **Credential lookup**: session middleware → `ResolveContext` → `ResolveSession` | `ResolveSession` max 0.029 ms (V7) | **3 s** (= F) | `host/daemon/middleware.go` `Wrap`, `:49`: `ResolveContext(WithTimeout(r.Context(), credentialBudget), …)` replaces `Resolve`. The Background root in `authority/resolver.go` `Resolve` is deleted | **HTTP 503**, class `Timeout` (the existing `writeReadTimeout` envelope). Today a lookup error collapses to **401 SessionUnknown**, which tells a valid client its token is bad. That changes |
| B5 | **Durable tail on `/a2a/`** (row 106 M5): `GetReceiptContext` + `AppendIntentContext` + `CommitContext` together | Max: `GetReceipt` 0.088 ms, `AppendIntent` 1.07 ms, `Commit` 0.94 ms (V7). Sum 2.1 ms | **3 s** for the three together (= F) | Installed by **row 106's M5**, not here: `tail := WithTimeout(ctx, durableTailBudget)`. Recommendation to M5, not absorbed: run the capsule under `InvokeWait − durableTailBudget` so a slow transition cannot starve the tail | Pre-cutoff → R11 "not committed". Post-cutoff → `*store.UncertainError` → R16 "not confirmed; resend the same task id" |
| B6 | **`/v1/commit`**: `handleCommit → CommitContext` | `Commit` max 0.94 ms (V7) | **3 s** (= F) | `host/daemon/handlers.go` `handleCommit`, `:600`: `CommitContext(WithTimeout(r.Context(), commitBudget), commit)` | Pre-cutoff: **503** `Timeout`, "commit not committed; safe to resend". Uncertain: **503** class `CommitUncertain`, "outcome unknown; reconcile with `GET /v1/log/{entryIndex}` (or `GET /v1/receipts/{invocationId}` when one was sent); a different head does **not** mean it failed. Do not resend until reconciled." The reconciliation rule is §2.5 (revision 1, objection 1). Conflict stays **409** |
| B7 | **Pooled-connection acquisition** (inside every row above) | Uncontended `BeginTx`+`Rollback`: p99 4–7.6 µs, max 75 µs (V7). Contended: the prototype test holds the connection and the op returns at the 100 ms deadline, with margin (V10) | **No separate timer.** It consumes the enclosing operation's budget | `store.beginDurable` → `db.BeginTx(ctx, nil)` (database/sql waits on `ctx.Done()`, V3) | `store: begin <op>: context deadline exceeded`, i.e. definite not committed |
| B9 | **Reconcile scan** (new, revision 1): `runReconcile → PendingEffectIntents(200)` | `PendingEffectIntents(page 200)` over 200 pending effect intents: p50 1.61–1.65 ms / p99 1.86–1.97 ms / **max 2.06 ms** (V19) | **3 s** (= F) | `cmd/world-publish/main.go` `runReconcile`: `sctx := WithTimeout(Background, reconcileScanBudget)` for the scan at `:464`. The HTTP reconcile at `:483` keeps its own root; that root performs no store read (plumbing doc §4) | `world-publish: scan pending effect intents: … context deadline exceeded`, exit 1. Read-only, so nothing is left behind |
| B10 | **Transition publish CLI** (new, revision 1): `runTransitions → buildChanges → PublishSet` (store reads, `PutObject`, `CompareAndSetRegistryHead`) | Store ops: `PutObject` max 0.52 ms, `CompareAndSetRegistryHead` max 0.55 ms (V19), both below F. The same ctx also carries one `ailang check` subprocess per change, bounded by the existing `checkTimeout` = 10 s, which honours ctx (`archive/check.go:20,75`) | **`len(entries) × 10 s + 3 s`**. Derived from the existing per-check timer plus F, **not a new fit** | `cmd/world-publish/transitions.go` `runTransitions`: the root at `:90` moves after `readManifest` (entry count known) and becomes `WithTimeout(Background, n×checkTimeout + F)` | `world-publish transitions: … context deadline exceeded`, exit 1. The registry advance is one CAS: before it, nothing is published. Stored source objects are content-addressed and inert |
| B11 | **Session mint/revoke CLI** (new, revision 1): `runSessionMint → authority.Mint → MintSession`; `runSessionRevoke → authority.Revoke → RevokeSession` | `MintSession` max 0.43 ms, `RevokeSession` max 0.99 ms (V19) | **3 s** (= F) | `cmd/ailang-worldd/session.go` `:158` and `:218`: `WithTimeout(Background, sessionBudget)` replaces `context.Background()` | Non-zero exit, with a `mint session: … context deadline exceeded` / `revoke: …` line. Each write goes through `beginDurable`/`finishDurable` (M5b), so the §2.3 contract applies. An uncertain mint prints `outcome unknown: run session list …` rather than a token |
| B12 | **No production root** (new, revision 1): `broker.Recover` (→ `PendingIntents`, `GetReceipt`, `PendingEffectIntents`, `GetEffectReceipt`), `replay` (→ `PutVerifyResult`), and `PutWorld`, `SelectHead`, `AppendOutcome` (test fixtures only) | Max: `PendingIntents(page 1000)` over 1000 intents **6.73 ms**; `PutVerifyResult` 1.40 ms; `PutWorld` 0.41 ms; `SelectHead` 0.036 ms; `AppendOutcome` 0.67 ms; `GetEffectReceipt` 0.16 ms (V19) | **None installed.** There is no root to install it at (V17). Under the rule, a future caller's budget is F = 3 s | Nowhere. After M7 the guard refuses any deadline-free call, so a future production caller **must** bring a deadline | `store.ErrNoDeadline` at the boundary (M7) |
| B8 | **Non-cancellable durable step** (the driver `COMMIT` after the cutoff) | Whole `Commit` max 0.94 ms (V7). Worker side: SQLite lock ≤ `busy_timeout` 2 s, then fsync | **Caller side: enforced** = the op's remaining budget. **Worker side: not enforceable in Go**; one lock window (row 111), then fsync | `store.finishDurable` selects on `ctx.Done()` and returns `*UncertainError`. The COMMIT goroutine holds the sole connection, so **at most one** exists per store (goroutine-count tested, V10) | `*store.UncertainError`, which is never `errors.Is(…, context.DeadlineExceeded)` |

**Human think time consumes no budget, by construction.** Every CLI budget (B2, B3′) is created **after** `requireAttendedOperator` returns (`main.go:258` → `:269`; `:352` → `:382`, V4). The human's typing happens before the timer exists. **No production poll loop waits on a human**: `MintAttendedApproval` performs request → decide → poll in one call. The human attends the CLI and does not answer a remote poll (plumbing doc F2, V4). If a later row adds a human-waiting poll loop, each poll step gets its own `WithTimeout(parent, B2)` inside the loop, and the sleep between steps sits outside every timer.

**B2 is a policy value, not a fit.** The measured work is ~50 ms after the multiplier. The budget is human-facing, and `walkApprovalHead` is O(approval-head objects) (`approve.go:230`), which grows without a measured ceiling. **Recommended: 5 s.** Alternative: **15 s**, if Mark prefers "never time out an operator on a large approvals history" over "fail fast".

### 1.1 Inventory of every exported context-free Store method at `d8db0eb` (revision 1, objection 2)

Production call sites: V16. Roots: V17. Test-only fixture use: V18. "In place" means the signature gains a leading `ctx` in that milestone, and every call site, test sites included, migrates mechanically in the same milestone. No wrapper is kept.

| Method (`file:line@d8db0eb`) | Production call sites (outside `host/store`) | Production root(s) | Budget row | Milestone |
|---|---|---|---|---|
| `Commit` (`store.go:874`) | `coordinator.go:285`; `daemon/handlers.go:600` | row 106 M5 (`/a2a/`); `POST /v1/commit` | B5; B6 | M1 (`CommitContext` + wrapper), M2, wrapper removed M6 |
| `AppendIntent` (`journal.go:411`) | `coordinator.go:282` | row 106 M5 | B5 | M1, wrapper removed M6 |
| `GetReceipt` (`journal.go:813`) | `coordinator.go:208`; `broker/recover.go:157` | row 106 M5; `broker.Recover` has **no** production caller | B5; B12 | M1, wrapper removed M6 |
| `PutObject` (`store.go:451`) | `approve.go:115,186,208`; `broker.go:217,301,358`; `registry.go:165`; `transitionreg.go:263`; `cmd/world-publish/transitions.go:308` | runApprove, runPublish, runServe (bootstrap), runTransitions | B2, B3′, B1, B10 | M5a |
| `SetRegistryHead` (`store.go:618`) | `approve.go:211`; `registry.go:168` | runApprove; runServe | B2; B1 | M5a |
| `AppendNextEffectIntent` (`journal.go:458`) | `broker.go:252` | runApprove (mint request invoke), runPublish | B2; B3′ | M5a |
| `AppendClaimedEffectIntent` (`journal.go:553`) | `broker.go:249` | runPublish | B3′ | M5a |
| `AppendEffectOutcome` (`journal.go:718`) | `broker.go:291,313` | runApprove, runPublish | B2; B3′ | M5a |
| `GetEffectReceipt` (`journal.go:852`) | `broker/recover.go:232` | none | B12 | M5a |
| `CompareAndSetRegistryHead` (`store.go:682`) | `transitionreg.go:266` | runTransitions | B10 | M5b |
| `MintSession` (`store.go:1083`) | `authority/mint.go:54` | runSessionMint (`session.go:158`) | B11 | M5b |
| `PendingEffectIntents` (`journal.go:954`) | `broker/recover.go:214,216`; `cmd/world-publish/main.go:464` | runReconcile; Recover (none) | B9; B12 | M5b |
| `PendingIntents` (`journal.go:906`) | `broker/recover.go:137,139` | none | B12 | M5b |
| `ScanUnreadableLog` (`scan.go:53`) | `daemon.go:566` | runServe (`New → scanIntegrity`) | B1 | M5b |
| `ScanUnreadableWorlds` (`scan.go:93`) | `daemon.go:588` | runServe | B1 | M5b |
| `PutVerifyResult` (`store.go:751`) | `replay/replay.go:204` | none (`host/replay` has no production importer) | B12 | M5b |
| `PutWorld` (`store.go:508`) | **none** | — | B12 | M5b |
| `SelectHead` (`store.go:838`) | **none** | — | B12 | M5b |
| `AppendOutcome` (`journal.go:663`) | **none** | — | B12 | M5b |

**Why the three uncalled methods are converted, not removed.** The reviewer's fix allows removal when a method has no production caller. These three, though, are cross-package **test fixtures**: `PutWorld` has 7 sites and `SelectHead` 7 in `host/daemon` and `host/broker` tests, and `AppendOutcome` 1 in `broker/recover_test.go` (V18). Those packages cannot reach `s.db` to plant rows another way. Converted and guarded, they cost a test deadline and cannot bypass the boundary. Removal would push raw-SQL fixture helpers into an exported test-support surface, which is a larger change.

**Every write in M5a/M5b goes through `beginDurable`/`finishDurable`,** single-statement ones included. The §2.3 three-outcome contract is then uniform across the surface. Whether `ExecContext` on an autocommit statement can report a ctx error after a durable write is not re-derived per method.

**Outside the method guard, by design:**
- `BusyTimeout` (`store.go:293`) returns a cached value and does no I/O.
- `Close` (`store.go:409`) is lifecycle: it waits for in-flight operations, and after M7 each of those is itself deadline-bounded.
- The package functions `Open`/`OpenReadOnly` are not `*Store` methods. Each builds a **fresh** pool, so its I/O cannot wait on another holder's connection; lock waits are bounded by `busy_timeout`. `Open` measured max 5.23 ms (V7), and inside B1 for the daemon.

M7's reflect test pins this allow-list to exactly `{BusyTimeout, Close}`.

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
| **uncertain** | `store.IsUncertain(err)`. `*UncertainError` **does not unwrap to ctx errors** (MUT-UNCERTAIN-UNWRAPS kills that) | Committed or not, never partial | **Reconcile, never blindly retry**: invocation → `GetReceiptContext(id)` (resolved ⇒ committed; indeterminate ⇒ not committed, per row 106's R15 argument on the single-connection pool). Bare commit → the log-index witness of §2.5, **never** the selected head. Because the sole connection is held until COMMIT settles, **the next acquisition observes the settled outcome** |

### 2.4 Is an error from `tx.Commit()` itself uncertain? Yes, by contract (conservative)

- `ctx.Err()` / `sql.ErrTxDone` from `Tx.Commit` are **definite not-committed**. database/sql returns them only when it did **not** win the CAS, i.e. the driver COMMIT never ran (`sql.go:2291-2301`). `classifyCommit` (`durable.go:85`) maps them to "not committed".
- A **driver** error is classified **uncertain**. The evidence goes both ways. On `SQLITE_BUSY`, modernc checks `sqlite3_get_autocommit` and forces a `rollback` (`modernc tx.go:35-53`, V3), and that case is really not-committed. But an I/O error during the journal delete/fsync phase can leave a hot journal, and SQLite resolves it at the next open. The prototype cannot inject that, and this design does not assume it. Calling every driver error uncertain is sound (the receipt read settles it), costs one read, and cannot turn a landed commit into "failed". **Unproven arm**: the prototype cannot reach a driver COMMIT error, so MUT-DRIVERERR-NOTCOMMITTED **survived**. M1 adds the fault-injection seam (§5, §6).

### 2.5 Reconciling an uncertain *bare* commit (revision 1, objection 1)

**Adopted verbatim from the reviewer:** *"After CommitUncertain, a matching head can confirm success only under verified head-identity invariants; a nonmatching head leaves the outcome unknown and never authorizes resubmission. Definitive reconciliation requires durable evidence tied to the original commit."*

**Head-identity invariants are not verified in this store, so the head is not used at all.** The store checks that `NextWorld.Ref` and `Entry.EntryHash` are well-formed references (`validateRef`). It verifies object payloads (`verifyObject`), but it never recomputes a world ref or an entry hash (V20). A `worlds` row can also be written by `PutWorld` as well as by `Commit` (`INSERT OR IGNORE`, `store.go:519,971`, V21). So neither a matching head nor a present world row is evidence that *this* commit landed. The controller's starting point ("the next World row and the log entry are written") was checked: the world row **fails** as a witness, and the log row **passes**, with one stated limit.

**Durable evidence that already exists, inventoried:**
1. **The intent/receipt journal.** A `Commit` carrying an `InvocationID` writes its `committed` outcome in the same transaction as world, log and head (`store.go` step 6). `journal` has `UNIQUE (invocation_id, kind)` (`schema.sql:81-87`, V22). So a resolved receipt ⇔ that commit landed. **`/v1/commit` does not use it today**: the wire has no invocation id (`handlers.go:72-77`, V22).
2. **The log row at the commit's own index.** `log_entries.entry_index` is the `PRIMARY KEY` and `entry_hash_ref` is `UNIQUE` (`schema.sql:36-38`). The only production writer is `Commit`'s `INSERT INTO log_entries` (`store.go:981`). No `UPDATE`/`DELETE` of `log_entries` exists in `host/` or `cmd/` (V21, with the `INSERT` as positive control). The row at A's index is therefore **A's own slot in immutable history**.

**Rule INV-LOG** (no new outcome machinery; M2 promotes the prototype's test helper `logWitness` to `store.CommitLanded(ctx, c)`). Read the log row at `A.Entry.Header.EntryIndex` once A's outcome has settled:

| Observation | Verdict | Why |
|---|---|---|
| no row at that index | **A not committed** — resubmission allowed | A's transaction inserts exactly this row, and rows are never deleted |
| a row that differs from `A.Entry` in any stored field | **A not committed** — resubmission allowed, and it will get a conflict | The index is the primary key, so A's insert could not also have succeeded |
| a row equal to `A.Entry` in every stored field | **A's log entry is durable** — do **not** resubmit | See the stated limit below |
| the read itself times out (503) | **unknown** — no resubmission advice | The settled outcome has not been observed |

**Stated limit.** Because the store binds neither `NextWorld` to the entry nor `EntryHash` to its header (V20), a *different*, inconsistent commit A′ carrying an identical log row but a different world would also match. The daemon therefore reports that arm as `entryLanded`, a true statement about the log, and not as "A's world is durable". **Settledness**: the reconciling read runs on the same sole pooled connection, so it cannot begin until A's COMMIT goroutine has released it (§2.3).

**Full commit identity: reuse the receipt machinery (M3, opt-in, additive).** `/v1/commit` accepts an optional `invocationId` (namespace `rest:`; `effect:` stays refused). The handler appends the intent derived from the commit (`LogicalTime = EntryIndex`), then commits with that `InvocationID`, and a new `GET /v1/receipts/{id}` exposes the receipt. For those clients a resolved receipt ⇔ *this* commit landed, world included. That is the reviewer's "durable evidence tied to the original commit", reused rather than duplicated. Clients that omit the id get INV-LOG. **No path ever compares the head.**

**Test** (prototype, M1 seam, V10/V12): `TestUncertainBareCommitReconcilesByLogIndexNotHead`.
- Arm 1, `A_landed_then_B_on_top`: A's COMMIT completes, then a post-commit stall (`durableAfterCommitHook`) outlasts A's deadline, so the caller gets `*UncertainError`. B commits on top of A. The test asserts that the head is **not** A's world, which is what the unsafe head rule would have read as "A failed", and that `logWitness(A)` = **landed**.
- Arm 2, `A_rolled_back_then_B_takes_the_index`: A stalls **before** the cutoff, is rolled back and gets `*UncertainError`. B lands at A's index from genesis. `logWitness(A)` = **not committed**.

MUT-AFTERHOOK-BEFORE-COMMIT and MUT-UNCERTAIN-AS-NOTCOMMITTED-R1 are both killed (§6.1).

## 3. The strict guard — the closing move, landing LAST (M7)

`requireDeadline(ctx) error` rejects `ctx == nil` or `_, ok := ctx.Deadline(); !ok` with `store.ErrNoDeadline`. **Revision 1: it covers the entire exported `*Store` I/O surface, 28 methods**, installed as each method's first statement:

- the **9 already ctx-accepting at `d8db0eb`**: `GetObject`, `GetWorld`, `GetLogEntry`, `GetRegistryHead`, `GetVerifyResult`, `SelectedHead`, `ReadObject`, `ResolveSession`, `RevokeSession`;
- the **19 converted here** (§1.1). All carry their original names with a leading `ctx` by the end of M6: M5a/M5b convert in place, and M6 removes M1's `…Context` wrappers and renames them back: `Commit`, `AppendIntent`, `GetReceipt` (M1); `PutObject`, `SetRegistryHead`, `AppendNextEffectIntent`, `AppendClaimedEffectIntent`, `AppendEffectOutcome`, `GetEffectReceipt` (M5a); `CompareAndSetRegistryHead`, `MintSession`, `PendingIntents`, `PendingEffectIntents`, `ScanUnreadableLog`, `ScanUnreadableWorlds`, `PutVerifyResult`, `PutWorld`, `SelectHead`, `AppendOutcome` (M5b).

**Surface test** `TestStoreSurfaceTakesContext` (M7, reflect over `reflect.TypeOf(&Store{})`). Every exported method except the allow-list must have `ctx context.Context` as its first parameter (`Type.In(1)`). A second assertion pins the allow-list to exactly `{BusyTimeout, Close}`, with the reasons given in §1.1. A new exported context-free method, or a growth of the allow-list, turns it red.

**Census that proves no production root reaches the store without a deadline.** The landed `TestProductionContextRoots` (`context_roots_test.go`) pins every `context.Background/TODO/WithoutCancel` reference by file|function. M7 flips it:
(a) every store-reaching root in `contextRootPins` is removed, or re-pinned **with the `WithTimeout` that wraps it in the same function**. A new AST rule requires each surviving `Background()` in a store-reaching file to be the direct first argument of `context.WithTimeout`/`WithDeadline`.
(b) The four chains of the plumbing doc's guard probe are re-run as tests **with the guard installed**, and all four pass. The same four with the budget removed fail with `ErrNoDeadline`: one mutation per root.
(c) A dynamic table test calls **each of the 28** guarded methods with `nil`, `context.Background()` and a deadline ctx, and expects `ErrNoDeadline`, `ErrNoDeadline` and success. The table is generated from the same reflect walk, so a new method cannot be forgotten.

No context-free wrapper exists by M7: M1's three are removed in M6, and M5a/M5b convert in place. A Background-wrapping store method would make the guard trivially bypassable, and the surface test forbids one.

## 4. Why is this not a package? (S3)

This is Go host code **at the effect boundary** (S2): transaction lifetime, a connection pool and HTTP/CLI budgets. It has no AILANG semantics to publish. A package cannot own the store's own transaction. No `.ail` is touched, and the kernel does not grow.

## 5. Milestones (dependency order, each ≤150 production LOC)

| M | Scope | Prod LOC | Unblocks / tests |
|---|---|---|---|
| **M1** (prototype) | `durable.go` (`beginDurable`, `finishDurable`, `classifyCommit`, `notCommitted`, `UncertainError`, and the test-only `durableCommitHook`/`durableAfterCommitHook` seams). `CommitContext`, `AppendIntentContext`, `GetReceiptContext`. `journalRowFor(ctx, …)`. **Thin wrappers kept**: `Commit`, `AppendIntent` and `GetReceipt`, because 65 test call sites and 2 production sites (`handlers.go:600`, `broker/recover.go:157`) would otherwise land in this milestone (V2). The census gains 3 named rows (`journal.go\|AppendIntent`, `\|GetReceipt`, `\|GetEffectReceipt`), each marked "removed by M6". **Additions over the prototype**: a driver fault seam (a `sql.Register`ed wrapping driver used only by tests, whose `Tx.Commit` returns an error) to kill MUT-DRIVERERR-NOTCOMMITTED. "Uncertain → reconciles committed" is already proven by the revision-1 `durableAfterCommitHook` arm (§2.5). A `commitBodyHook(ctx)` stall inside the transaction body to kill MUT-HEADREAD-BG/MUT-EXEC-NOCTX. A mid-body cancel for `AppendIntentContext` to kill MUT-APPEND-NOTCOMMITTED-DROP | ~130 (prototype: 63 code lines in `durable.go` + 61 changed lines in store/journal, V11) | **Row 106 M5** (it adopts these three methods per R-106-11). No duration value in M1 |
| **M2** | `/v1/commit` → `CommitContext` with B6. The 503 `Timeout` / 503 `CommitUncertain` mapping. `store.CommitLanded(ctx, c)`, i.e. INV-LOG promoted from the prototype's `logWitness` (§2.5). `docs/QUICKSTART.md`: the uncertain-commit step reads `GET /v1/log/{entryIndex}`, never the head (S7) | ~110 | AC5 |
| **M3** | Opt-in `invocationId` on `/v1/commit`, reusing the intent/receipt journal (§2.5); `GET /v1/receipts/{id}`; QUICKSTART | ~90 | AC5b |
| **M4** | Credential lookup B4: middleware → `ResolveContext` with timeout, 503 on lookup expiry. Delete the `authority/resolver.go` `Resolve` Background root | ~40 | AC6 |
| **M5a** | In place, ctx first, through `beginDurable`/`finishDurable`: `PutObject`, `SetRegistryHead`, `AppendNextEffectIntent`, `AppendClaimedEffectIntent`, `AppendEffectOutcome`, `GetEffectReceipt`. Callers: `broker/approve.go`, `broker/broker.go`, `broker/recover.go`, `registry/registry.go`, `transitionreg/transitionreg.go`, `cmd/world-publish/transitions.go`. Test sites are migrated mechanically here | ~140 | AC7a |
| **M5b** | In place: `CompareAndSetRegistryHead`, `MintSession`, `PendingIntents`, `PendingEffectIntents`, `ScanUnreadableLog`, `ScanUnreadableWorlds`, `PutVerifyResult`, `PutWorld`, `SelectHead`, `AppendOutcome`. `broker.Recover(ctx, …)`, `Daemon.scanIntegrity(ctx)`, `authority.Mint`/`Revoke` forward ctx. `runReconcile` temporarily gets a named census root for `:464` until M6. Test sites migrated here | ~140 | AC7b |
| **M6** | Root budgets B1, B2, B3, B3′, B9, B10, B11 at their install sites. **Remove M1's wrappers** and rename `…Context` → the original names; migrate the 2 production and ~65 test call sites mechanically | ~120 | AC8–AC10, AC14 |
| **M7** | Strict guard over all 28 methods (§3), `TestStoreSurfaceTakesContext`, census flip, the four-chain probe as tests | ~70 | AC11–AC13. **Lands last** |

Row 106's M5 needs **M1 only** (unchanged by revision 1). M2–M7 complete `D-WORLD-37`.

## 6. Acceptance criteria and non-vacuity

### 6.1 Executed on the prototype (16 mutations: 12 killed, 4 survived; re-run in full after revision 1)

Harness: `~/.ailang/state/world-iter197/designer-scratch/mutate.py`. It applies each mutation, runs `go test ./host/store/ -run <selector> -count=3` (revision 1 re-run; `-count=5` originally), restores the file in `finally`, and writes transcripts to `…/mut/<NAME>.txt` (V12). Tests have a watchdog (`boundedCall`), so a ctx-ignoring mutant **fails** instead of hanging.

| AC | Mutation | Killing assertion (observed red) | Result |
|---|---|---|---|
| AC1 durable ops honour a held sole connection | MUT-BEGIN-NOCTX: `BeginTx(ctx,nil)` → `Begin()` | `TestDurableOpsHonourHeldConnection`: `durable_test.go:90: still blocked 1.10s after a 100ms deadline` | **KILLED** |
| AC1 (receipt arm) | MUT-RECEIPT-NOCTX: `GetReceiptContext` reads under `Background()` | `…/GetReceiptContext`: `still blocked 1.10s after a 100ms deadline (err after unblock: <nil>)` | **KILLED** |
| AC2 caller returns at deadline past the cutoff | MUT-CUTOFF-NOSELECT: wait on `<-done` only | `TestDurableStallAtCutoffIsUncertainAndReconciles`: `still blocked 1.10s after a 100ms deadline` | **KILLED** |
| AC3 uncertain ≠ ctx error | MUT-UNCERTAIN-UNWRAPS: add `Unwrap() → Cause` | `…:181: err = store: commit outcome is uncertain …; want *UncertainError that is not a ctx error` | **KILLED** |
| AC3 (other direction) | MUT-UNCERTAIN-AS-NOTCOMMITTED: post-cutoff expiry reported as not committed | `…:181: err = store: commit not committed: context deadline exceeded; want *UncertainError …` | **KILLED** |
| AC4 no accumulating blocked worker | MUT-LEAK-UNBUFFERED: unbuffered `done` channel | `…:192: goroutines 5 > baseline 4 after release` | **KILLED** |
| AC2 (Commit wiring) | MUT-COMMIT-DROP-CUTOFF: `return tx.Commit()` | `…:181: err = <nil>; want *UncertainError …` | **KILLED** |
| AC2 (AppendIntent wiring) | MUT-APPEND-DROP-CUTOFF | `…/AppendIntentContext:181: err = <nil>; want *UncertainError …` | **KILLED** |
| AC5′ pre-cutoff cancel reports its cause, all-or-nothing | MUT-NOTCOMMITTED-DROP (Commit) | `TestCommitCancelledMidBodyIsAllOrNothing:132: err = store: next journal seq: sql: transaction has already been committed or rolled back; want definite Canceled` | **KILLED** (after the test was made deterministic, see below) |
| AC5 (bare-commit reconciliation, rev. 1) | MUT-AFTERHOOK-BEFORE-COMMIT: the post-commit seam moved before `tx.Commit()`, so A never lands | `TestUncertainBareCommitReconcilesByLogIndexNotHead:281: B: store: stale observed world head …` | **KILLED** |
| AC3 (bare-commit arm, rev. 1) | MUT-UNCERTAIN-AS-NOTCOMMITTED-R1: post-cutoff expiry reported as not committed | `…:264: A: err = store: commit not committed: context deadline exceeded; want *UncertainError` | **KILLED** |
| census pins wrapper roots | MUT-WRAPPER-TODO: wrapper uses `context.TODO()` | `TestProductionContextRoots: context_roots_test.go:118: root references: Background=18 TODO=1` | **KILLED** |
| (in-body ctx) | MUT-HEADREAD-BG: head compare under `Background()` | — | **SURVIVED**: no in-body stall seam. M1 obligation (`commitBodyHook`) |
| (in-body ctx) | MUT-EXEC-NOCTX: world insert `ExecContext` → `Exec` | — | **SURVIVED**: same seam. Point statements cannot be stalled without it |
| AC3 driver arm | MUT-DRIVERERR-NOTCOMMITTED: driver COMMIT error → "not committed" | — | **SURVIVED**: no driver fault seam. M1 obligation (wrapping test driver) |
| AC5′ (AppendIntent) | MUT-APPEND-NOTCOMMITTED-DROP | — | **SURVIVED**: no mid-body hook in `AppendIntentContext`. M1 obligation |

Two findings came from running the mutations, not from design. (1) MUT-NOTCOMMITTED-DROP first **survived at 25/30**: `awaitDone`'s rollback is asynchronous, so the `ErrTxDone` arm is reached only when the rollback wins the race. The test hook now cancels and then yields 50 ms, which makes the arm deterministic (killed 5/5 runs). (2) An explicit `ctx.Err()` pre-check before `tx.Commit()` survived its removal, because database/sql performs it (§2.1). It was deleted rather than kept as untestable code (S6).

### 6.2 Acceptance for M2–M7 (mutation named per AC; refusal branches get one mutation each)

| AC | Test (to be written in its milestone) | Mutation that must turn it red |
|---|---|---|
| AC5 `/v1/commit` bounded; uncertain reconciles by log index, not head | (i) Hold the connection: POST answers **503 Timeout** within B6 + ε and the head is unchanged. (ii) **Objection-1 test over HTTP**: commit A lands, then stalls post-commit (`durableAfterCommitHook`) → **503 CommitUncertain**. Commit B lands on top. `CommitLanded(A)` → `entryLanded` while `/v1/head` ≠ A's world. (iii) A is rolled back pre-cutoff, B takes A's index → `CommitLanded(A)` = not committed. (iv) Reconcile read times out → unknown, and the body carries no resend advice | MUT-M2-NOCTX (`Commit` wrapper kept at `:600`); MUT-M2-UNCERTAIN-AS-TIMEOUT (both branches return one class); **MUT-M2-HEAD-WITNESS** (`CommitLanded` compares `SelectedHead` with `NextWorld.Ref` → arm (ii) red); MUT-M2-HASH-ONLY (compares `EntryHash` only → a same-hash, different-header row is misread; fixture row planted); MUT-M2-ABSENT-AS-UNKNOWN (arm (iii) red) |
| AC5b opt-in receipt identity | With an `invocationId`: A lands then uncertain, B lands; `GET /v1/receipts/{id}` → resolved, `resultRef` = A's world. Rolled-back arm → indeterminate | MUT-M3-NO-INTENT (commit sent without `InvocationID` → receipt not-started); MUT-M3-EFFECT-NS (the `effect:` prefix is accepted) |
| AC6 lookup expiry is not a denial | Hold the connection. A valid token gets **503**, not 401 | MUT-M3-COLLAPSE (lookup error → `DenialUnknown`); MUT-M3-NOCTX (keep `Resolve`) |
| AC7a / AC7b every converted method cancellable | Per method: held connection → definite not-committed ctx error within budget + ε; for writes, a stall at the cutoff → `*UncertainError` | One MUT-M5-NOCTX-<method> per converted method: **6** (M5a) + **10** (M5b) |
| AC8 startup bounded | `Run` with a held connection / stalled bootstrap write → `StartupError` wrapping `DeadlineExceeded` within B1 + ε, and the writer lock is released (re-Open succeeds) | MUT-M6-STARTUP-SIGCTX (pass the signal ctx to `New`); MUT-M6-SERVE-BOUNDED (the budget wraps the serve lifetime, so a served request fails after 9 s) |
| AC9 approval budget excludes think time | Fake attended prompt blocks 2×B2, then mint succeeds | MUT-M6-TIMER-BEFORE-PROMPT (timer created before `requireAttendedOperator`) |
| AC10 validation-only budget | A stalled validation read errors within B3, and the approval is not consumed. A handler exec taking > B3 still succeeds | MUT-M6-VALIDATION-UNBOUNDED; MUT-M6-BUDGET-WRAPS-DISPATCH |
| AC11 guard rejects deadline-free | Reflect-generated table: all **28** methods × {nil, Background, deadline} | One MUT-M7-UNGUARDED-<method> per method (**28**); MUT-M7-NIL-ONLY (checks nil but not the deadline) |
| AC11b no context-free method on the surface | `TestStoreSurfaceTakesContext` | MUT-M7-CONTEXT-FREE-METHOD (add `func (s *Store) Ping() error`); MUT-M7-ALLOWLIST-GROW (add `Ping` to the allow-list → the exact-allow-list assertion fires) |
| AC14 new CLI roots bounded | B9/B10/B11: held connection → CLI exits 1 within budget + ε with the named line | MUT-M6-RECONCILE-BG, MUT-M6-TRANSITIONS-BG, MUT-M6-SESSION-BG (each root reverted to bare `Background()` → the AC12 census and this test both red) |
| AC12 census flip | `TestProductionContextRoots` with the WithTimeout-adjacency rule | MUT-M7-BARE-BACKGROUND (a bare `Background()` passed to a store call in `runApprove`) |
| AC13 four chains pass under the guard | The four plumbing-doc probe chains as tests | Remove each root's budget (4 mutations) → `ErrNoDeadline` |

**Row-106 M5-style test (required, already in the prototype):** `TestDurableOpsHonourHeldConnection` holds the sole pooled connection (`s.db.Conn`) past a 100 ms caller deadline. For each of `CommitContext`/`AppendIntentContext`/`GetReceiptContext` it asserts `errors.Is(err, context.DeadlineExceeded) && !IsUncertain(err)`, `elapsed ≤ 100 ms + 250 ms`, and `runtime.NumGoroutine()` back at baseline. It then asserts the store is untouched and that a released store commits (live control). `TestDurableStallAtCutoffIsUncertainAndReconciles` injects the stall at the cutoff for both writes, and `TestUncertainBareCommitReconcilesByLogIndexNotHead` covers objection 1 (V10).

## 7. Conflict Surface

| File (this tranche) | M | Other open rows touching it |
|---|---|---|
| `host/store/durable.go` (new), `store.go`, `journal.go`, `scan.go`, `context_roots_test.go` | M1, M2, M5a, M5b, M6, M7 | **111** arm A (per-request busy_timeout on store reads: same files, and it must keep pool-state restoration). **26** (producer source reads and envelope write: it must adopt `*Context` writes and B-budgets, not a hidden root). **106** M5 (adopts M1's methods; interface in `host/coordinator/coordinator.go:27-34`) |
| `host/daemon/handlers.go`, `middleware.go`, `daemon.go` | M2, M3, M4, M5b, M6 | **111** (daemon read-path status mapping: coordinate the 503 envelopes). **106** M5 (daemon constructs the coordinator) |
| `host/authority/resolver.go` | M4 | none measured |
| `host/broker/approve.go`, `broker.go`, `publish_op.go`; `host/registry/registry.go` | M5a, M6 | **25** (capsule blocked-child coverage): no file overlap. **112** touches `replay.go`/archive, not these |
| `cmd/world-publish/main.go`, `transitions.go`; `cmd/ailang-worldd/main.go`, `session.go` | M5a, M5b, M6 | none measured |
| `host/transitionreg/transitionreg.go`, `host/authority/mint.go`, `host/broker/recover.go` | M5a, M5b | **106** (the coordinator is a transitionreg consumer; no file overlap) |
| `host/replay/replay.go` (`PutVerifyResult` call at `:204`, inside `ReplayEntry`) | M5b | **112** edits the same file in `runPinnedTransition`, a different function. Sequence either way; rebase only |

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
| V10 | Prototype tests green and stable | `go test ./host/store -run 'Durable\|MidBody' -count=20`; revision 1: `-run 'Durable\|MidBody\|BareCommit' -count=10` | `ok … 10.768s` (before the stall test was parametrised); after: `-count=10` `ok … 6.749s`; revision 1: `ok … 9.081s` |
| V11 | Prototype size | `git diff --numstat`; `grep -vc '^\s*//\|^\s*$' host/store/durable.go`; changed non-comment lines in store/journal diff | store.go +25/−19; journal.go +36/−14; context_roots_test.go +19/−15 (gofmt realigned the map); durable.go 109 lines / 63 code; 68 changed code lines in store+journal |
| V12 | Mutation run | `python3 ~/.ailang/state/world-iter197/designer-scratch/mutate.py` | Revision 1 re-run: 16 rows as in §6.1 (12 KILLED, 4 SURVIVED); transcripts in `…/mut/`. A 15th mutation (remove an explicit pre-cutoff `ctx.Err()` check) survived and the check was deleted (§2.1) |
| V13 | Census catches the new wrappers | first `go test ./host/store/` after adding wrappers, before the pin update | `TestProductionContextRoots: … Background=20` with `journal.go\|AppendIntent`, `\|GetReceipt`, `\|GetEffectReceipt:2` unexpected (reduced to 1 by a local `ctx`) |
| V14 | Residual owners OPEN | `grep -n '^25\. \|^26\. \|^111\. \|^112\. \|^113\. ' design_docs/world-mission.md` | 25 `w-capsule-blocked-child-kill-coverage`, 26 (gate met, open), 111 `w-daemon-lock-blocked-read-status-and-deadline-enforcement`, 112 `w-archive-replay-grandchild-bound`, 113 `w-pkgproj-quality-descendant-leak`, all present as queue entries |
| V16 | Every exported context-free method and its production call sites (positive: `MintSession` → `authority/mint.go:54` in the same loop) | `for f in host/store/*.go (non-test); git show d8db0eb:$f \| grep -n '^func (s \*Store) [A-Z]'`; then `for m in <19 names>; grep -rn "\.$m(" host cmd --include='*.go' \| grep -v _test.go \| grep -v '^host/store/'` | 30 exported methods: 9 take ctx, 2 non-I/O (`BusyTimeout :293`, `Close :409`), **19 context-free**. Call sites exactly as in §1.1. `PutWorld: NONE`, `SelectHead: NONE`, `AppendOutcome: NONE` (negatives in the same loop as 16 positives) |
| V17 | Roots: `broker.Recover` and `host/replay` have no production caller (positive: the `func Recover` declaration itself is found) | `grep -rnw 'Recover\|PendingPublishes' host cmd --include='*.go' \| grep -v _test.go`; `grep -rl '"github.com/sunholo-data/ailang-world/host/replay"' host cmd --include='*.go' \| grep -v _test.go` next to the same grep for `host/store` | Only `recover.go:77,91,98,101,108` (the declarations and comments), and no call site. `host/replay`: **0** production importers, against **23** for `host/store` (positive control). `session.go:158 authority.Mint(context.Background()…)`, `:218 authority.Revoke(context.Background()…)`; `main.go:464 db.PendingEffectIntents(reconcileScanLimit)`; `transitions.go:90 ctx := context.Background()` |
| V18 | Cross-package test fixtures use the uncalled methods | `for m in PutWorld SelectHead AppendOutcome …; grep -rn "\.$m(" host cmd --include='*_test.go' \| grep -v '^host/store/'` | `PutWorld` 7 (broker/recover_test, daemon handlers/integrity/bench/daemon tests); `SelectHead` 7 (same files); `AppendOutcome` 1 (`broker/recover_test.go`); `PutVerifyResult` 0 outside store (control) |
| V19 | Surface timings (revision 1) | `for r in 1 2 3; do WORLD_BOUND_MEASURE=1 go test ./host/store/ -run TestMeasureContextFreeSurface -count=1 -v; done` (banked `…/measure_surface.txt`) | Max over 3 runs: PendingIntents(page 1000) **6.73 ms**; PendingEffectIntents(page 200) **2.06 ms**; AppendEffectOutcome 1.79 ms; PutVerifyResult 1.40 ms; RevokeSession 0.99 ms; AppendNextEffectIntent 0.70 ms; AppendOutcome 0.67 ms; CompareAndSetRegistryHead 0.55 ms; PutObject 0.52 ms; MintSession 0.43 ms; SetRegistryHead 0.43 ms; PutWorld 0.41 ms; ScanUnreadableLog(64) 0.20 ms; GetEffectReceipt 0.16 ms; ScanUnreadableWorlds(64) 0.13 ms; SelectHead 0.036 ms. All ≤ 6.73 ms, so 10× < F, and F decides |
| V20 | The store verifies objects but not world/entry content addresses (positive: `verifyObject` found) | `sed -n '/^func (s \*Store) CommitContext/,/^func bindCommitIntentTx/p' host/store/store.go \| grep -n 'SumSHA256\|verifyObject'`; `sed -n '/^func decodeCommit/,/^}/p' host/daemon/handlers.go \| grep -c 'Sum\|Compute\|verify'` vs `grep -c parseRef` | Commit: one match, `verifyObject`; **no** `SumSHA256`. decodeCommit: 0 hash computations, 10 `parseRef` |
| V21 | Log rows are write-once, written only by Commit; world rows have two writers (positive: the INSERTs) | `grep -rnE 'INSERT (OR IGNORE )?INTO (log_entries\|worlds)\|(UPDATE\|DELETE FROM) (log_entries\|worlds)' host cmd --include='*.go' --include='*.sql' \| grep -v _test`; `grep -n 'entry_index *INTEGER PRIMARY KEY\|entry_hash_ref *TEXT NOT NULL UNIQUE\|world_ref *TEXT PRIMARY KEY' host/store/schema.sql` | `store.go:519` (PutWorld) and `:971` (Commit) `INSERT OR IGNORE INTO worlds`; `:981` `INSERT INTO log_entries`; **zero** UPDATE/DELETE. `schema.sql:24` world_ref PK, `:37` entry_index PK, `:38` entry_hash_ref UNIQUE |
| V22 | The receipt journal binds an invocation; the `/v1/commit` wire has no id | `sed -n 81,87p host/store/schema.sql`; `grep -n 'type commitRequest' -A6 host/daemon/handlers.go` | `UNIQUE (invocation_id, kind)`; `commitRequest{ObservedHead, Objects, NextWorld, Entry}`, with no invocation field |
| V15 | Full suite with the prototype | `AILANG_BIN=… go test ./... ` (banked `…/fullsuite.txt`), then the failing package alone, then `go list -deps ./host/pkgproj \| grep -E '^os/exec$\|ailang-world'` | **23 ok, 1 FAIL**: `host/pkgproj` `TestQueryInterfaceReturnsWhileADescendantHoldsStdout` (`iface_test.go:232: instrument: descendant pid not recorded`). Re-run alone: `ok host/pkgproj 11.425s`. Its deps list `os/exec` (positive) and `host/childenv`, and **not** `host/store`, so the prototype cannot reach it. Row 113's area; attributed as an instrument flake |

## 9. Decision for Mark — one word

**Ratify the bound table in §1 as proposed?**

| Answer | Meaning |
|---|---|
| **A** | Ratify B1–B12 as written: startup 9 s, approval active call 5 s, validation 3 s, publish root 36 s, credential lookup 3 s, durable tail 3 s, `/v1/commit` 3 s, acquisition inside the enclosing budget, uncertain-on-expiry after the cutoff, reconcile scan 3 s, transition publish `n × 10 s + 3 s`, session mint/revoke 3 s, and no budget where there is no production root (the guard refuses deadline-free callers). Also ratify the margin rule `max(10 × max_observed, busy_timeout + 1 s)`, the three-outcome contract of §2.3, and the bare-commit reconciliation rule of §2.5 (log index, never head; opt-in receipt identity). M1 routes first, and row 106 M5 unblocks. |
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
| ~~Context-free methods off the four root paths~~ | **Closed by revision 1.** All 19 are migrated (M5a/M5b) and guarded (M7), so the earlier "NEW ROW NEEDED" line is deleted |
| Store-side verification that `Entry.EntryHash` and `NextWorld.Ref` are content addresses. It would turn INV-LOG's `entryLanded` into full world identity for clients without an `invocationId` (§2.5 stated limit) | **NEW ROW NEEDED** (optional hardening, kernel validation; not required by `D-WORLD-37`). No existing queue row was found for it: a charter grep for entry-hash/world-ref verification found only failure-injection history rows |

## 11. Prototype manifest and results

**Files changed** (uncommitted, worktree `.design-wt-iter197`):
- `host/store/durable.go`: **new**. `UncertainError`, `IsUncertain`, `durableCommitHook`, `durableAfterCommitHook` (**revision 1**: a production no-op variable plus a 2-line call-order change in `finishDurable`, needed for objection 1's landed-then-uncertain test), `beginDurable`, `finishDurable`, `classifyCommit`, `notCommitted`.
- `host/store/store.go`: `Commit` → wrapper; body → `CommitContext(ctx, c)` using `beginDurable`, `selectedHeadTx(ctx, …)`, `ExecContext`, `finishDurable`; `bindCommitIntentTx(ctx, …)`.
- `host/store/journal.go`: `AppendIntent`/`GetReceipt` → wrappers; `AppendIntentContext`, `GetReceiptContext`; `journalRowFor(ctx, q, …)`; `GetEffectReceipt` keeps one local compatibility root.
- `host/store/context_roots_test.go`: 3 new pins (removed by M5).
- `host/store/durable_test.go`: **new**. `TestDurableOpsHonourHeldConnection`, `TestCommitCancelledMidBodyIsAllOrNothing`, `TestDurableStallAtCutoffIsUncertainAndReconciles` (2 subtests), and **revision 1** `TestUncertainBareCommitReconcilesByLogIndexNotHead` (2 arms) with its helpers `bareCommit` and `logWitness`.
- `host/store/durable_measure_test.go`, `host/daemon/startup_measure_test.go`: **new**. Gated measurement instruments (skip unless `WORLD_BOUND_MEASURE=1`). **Revision 1** adds `TestMeasureContextFreeSurface` (V19).

No `.ail` changed, so no `verify_ail.sh` run is needed.

**Revision 1 results**: `go vet ./...` rc 0. `go test ./host/verifygate ./host/runbook ./host/store -count=1` all ok. `go test ./host/store -run 'Durable|MidBody|BareCommit' -count=10` ok (V10). Mutations re-run in full: **16, of which 12 killed and 4 survived** (§6.1); the survivors are unchanged, named M1 obligations. The full suite was **not** re-run for revision 1; its only production change is the no-op `durableAfterCommitHook`.

**Original results**: `go vet ./host/... ./cmd/...` clean. `go test -run '^$' ./...` compiles. Targeted `go test ./host/store -run 'Durable|MidBody' -count=10` ok. Full suite: **23 packages ok, 1 FAIL** (`host/pkgproj`, an unrelated instrument flake that passes alone, V15). Mutations: **10 killed, 4 survived**, and each survivor is a named M1 obligation (§6.1).

## 12. Quorum log

### Round 1 — BLOCKED 2/2 (gpt6-astra, gemini-3-1-pro)

The glm and kimi seats were unreachable on the Ollama weekly limit. Claude was benched as the author's vendor. Neither reviewer disputed the design direction. Both objections were verified true by the controller before routing back, and both proposed fixes are applied.

**Objection 1 — gpt6-astra (bare-commit reconciliation).** Verbatim, short: *"§2.3 treats SelectedHead == NextWorld.Ref as the outcome check, and B6 tells clients to read /v1/head before resending. A commit can land, return UncertainError, and then be followed by another successful commit before reconciliation reads the head. A different current head therefore does not prove the uncertain commit failed."*

*Resolved:*
- **Text.** The reviewer's paragraph is adopted verbatim (§2.5). The head is removed from every reconciliation path: §2.3, B6 and M2's QUICKSTART step.
- **Evidence inventory.** Two candidate witnesses were measured. The world row **fails**: it has two writers (`PutWorld` and `Commit`, both `INSERT OR IGNORE`, V21), and no content-address verification (V20). The log row **passes** (INV-LOG): the primary key sits at the commit's own index, `Commit` is the only writer, and there is no UPDATE/DELETE (V21). The stated limit is that entry identity is not world identity (§2.5).
- **Reuse.** The intent/receipt journal is reused for `/v1/commit` as an opt-in `invocationId` (M3, V22), so no duplicate outcome machinery is added.
- **Test.** The reviewer's A-lands-uncertain-then-B-commits test is in **M1's prototype** (`TestUncertainBareCommitReconcilesByLogIndexNotHead`, both arms, 2 new mutations killed) and in **AC5 (M2)** over HTTP, with MUT-M2-HEAD-WITNESS.
- **Prototype change** (said so): the test-only `durableAfterCommitHook` seam in `finishDurable`, a production no-op.

**Objection 2 — gemini-3-1-pro (unmigrated context-free methods).** Verbatim, short: *"§0.4 identifies MintSession as a context-free write on a production root, but §10 contradicts this by claiming it is off the root paths and leaving it (along with PutWorld, SelectHead, etc.) unmigrated and unguarded. … any call to an unguarded context-free method will issue an unbounded Begin() that bypasses the M6 boundary guard … The strict boundary guard is structurally broken if it leaves unmigrated methods on the interface."*

*Resolved:*
- **The contradiction.** It was real (`MintSession` hangs off `runSessionMint`). §0.4 and §10 now agree, through one inventory (§1.1): all **19** exported context-free methods, with production call sites (V16), roots (V17) and budget rows.
- **New budget rows.** Three are measured (B9 reconcile scan, B10 transition publish, B11 session CLI; V19). B12 is explicit "no production root; the guard refuses" for `broker.Recover`, replay and the three fixture-only methods.
- **Milestones.** M4 is split into **M5a/M5b** (≤140 LOC each), converting in place with test sites migrated in the same milestone. The three uncalled methods are converted rather than removed, because they are cross-package test fixtures (V18).
- **Guard.** M7's guard covers all **28** exported I/O methods (AC11: 28 mutations + NIL-ONLY), plus `TestStoreSurfaceTakesContext`, which pins the allow-list to `{BusyTimeout, Close}` (AC11b).
- **Residuals.** The "NEW ROW NEEDED" line for off-path methods is deleted.
- **Correction.** The controller noted that "permanently starve" overstates the effect for *bounded* callers, which time out through `BeginTx(ctx)`. §0.4 says so, while keeping the substantive point.
- **Scope.** M1 is unchanged in scope; it is still exactly row 106 M5's prerequisite.
