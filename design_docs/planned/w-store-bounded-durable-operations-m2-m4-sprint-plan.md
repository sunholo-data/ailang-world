# Sprint plan — w-store-bounded-durable-operations, milestones M2–M4 (queue row 23, iteration 206)

**Base:** detached `b0ce897` = dev `58f6022` plus the rebased M1 (PR #153, being merged this iteration). **Authority:** [design](w-store-bounded-durable-operations.md) §1 (B4, B6), §2.3, §2.5, §5 rows M2/M3/M4, §6.2 rows AC5/AC5b/AC6, §7, §10, and the ratification **`D-WORLD-40` = A** (Mark, attended 2026-09-28): the §1 bound table (B6 = 3 s, B4 = 3 s), the three-outcome contract, and the §2.5 reconciliation rule (log index, never head; opt-in receipt identity) are binding. A credential-lookup timeout answers **503, not 401**. **Scope:** M2, M3, M4 only. M5a–M7 are not planned here. This planner made no git write in the repo. The only file written in the worktree is this plan.

**Prototype:** all three milestones were built in a scratch copy (`~/.ailang/state/world-iter206/planner-scratch/proto`, outside the repo) and saved as `~/.ailang/state/world-iter206/planner-scratch/prototype_m2_m4.diff` (1421 lines, sha256 `d85a2f88…4468`; `git apply --check` against `b0ce897` rc 0). **28 mutations (30 executions: M2-03 and M2-12 were run in both their M2-era and final-tree forms). All 28 are killed, 0 survive, 0 are unmeasured** (§3). The prototype is evidence that the plan is runnable. It is not a substitute for re-running the ledger on the landed files.

## 0. Premises re-measured against this tree

The design was written at `d8db0eb`. Every M2–M4 premise was re-measured at `b0ce897` (§5, V1–V16). Where a premise is false, the plan below has been adjusted to match.

| # | Design premise | Measured at `b0ce897` | Consequence for this plan |
|---|---|---|---|
| P1 | `/v1/commit` calls `d.store.Commit(commit)` at `handlers.go:600` | **`:669`**. The handler starts at `:615`, and the 200 body at `:678` (V1) | Line moved. The premise holds |
| P2 | The middleware calls `Resolve` at `middleware.go:49`; `resolver.go` `Resolve` is a Background root | `:49` holds. `resolver.go:128-129` `Resolve` → `ResolveContext(context.Background(), …)`. Census pin `host/authority/resolver.go\|Resolve\|Background` is at `context_roots_test.go:29` (V2, V6) | Holds |
| P3 | `GET /v1/log/{entryIndex}` exists | **Yes**: `daemon.go:644` `GET /v1/log/{index}` → `handleLogEntry` (V3) | M2 adds no route |
| P4 | `GET /v1/receipts/{id}` is new | **Absent**: 0 non-test mentions in `host/daemon` (V3) | M3 adds it, and that **changes the frozen /v1 table** (P5) |
| P5 | *(not stated in the design)* | The /v1 table is frozen by `TestFrozenV1RouteTableMatchesMux` (`route_table_test.go:51`: `len(mux) != 9 \|\| len(sketch) != 9`). That test compares the mux against `routes()` in `design_docs/sketches/worlddapi.ail:71-83`, which is Leg-1 `ai-check`ed by `verify_ail.sh` (V4) | M3 must edit **three** places in one commit: the sketch's `routes()`, the literal `9 → 10` in the route test, and the "nine /v1 patterns" doc comment in `daemon.go`. The sketch edit keeps `verify_ail.sh` rc 0 (V13). This is **OPEN-1** |
| P6 | 503 uses "the existing `writeReadTimeout` envelope" | `writeReadTimeout` (`handlers.go:338`) hard-codes the text `read deadline (%s) exceeded` (V5) | Reusing it would mislabel a commit or credential timeout as a read. The plan adds `writeCommitTimeout` and a middleware-local message. Both keep class `Timeout` and status 503 |
| P7 | *(not stated)* The API classes mirror the sketch's `ApiError`/`httpStatus` "exactly" (`handlers.go:21-23`) | `TestTimeoutStatusMirrorsSketch` **parses** the sketch's arms (`read_deadline_test.go:569-631`). Session classes are already absent from the sketch, so "exactly" is already not literally true | M2 adds `CommitUncertain(int) => 503` to the sketch (it is a /v1/commit class), plus a parsing mirror test. Sketch tests go 20 → 21. `verify_ail.sh` rc 0 with pins unmoved (16 identities, 40 named tests: Leg 2 only tests `world/`) (V13) |
| P8 | AC5 arm (iii) (B takes A's index) kills MUT-M2-ABSENT-AS-UNKNOWN | **False.** In (iii) the row at A's index **exists** (it is B's) and differs, so arm (iii) exercises the *differs* branch. The *absent* branch needs an arm where nothing lands at A's index | Added arm **(iii-b)**, "A rolled back, nothing else". MUT-M2-ABSENT-AS-UNKNOWN is killed there, over HTTP and at store level |
| P9 | AC6's mutations are named `MUT-M3-COLLAPSE` / `MUT-M3-NOCTX` | AC6 belongs to **M4** (§5) | Renamed `MUT-M4-*`. "Keep `Resolve`" is not a mutation once `Resolve` is deleted, so it is replaced by MUT-M4-NOCTX / MUT-M4-NOBUDGET / MUT-M4-RESTORE-BG-ROOT |
| P10 | M3 namespace `rest:` | The coordinator mints `a2a:<episode>:<task>` (`coordinator.go:100`). The store refuses `effect:` (`journal.go:325,837`) (V9) | `/v1/commit` accepts **only** `rest:<non-empty>`, which rules out `a2a:` collisions as well as `effect:`. `GET /v1/receipts/{id}` accepts the same namespace only |
| P11 | *(M7, recorded, not planned)* §3 counts 28 guarded methods | `b0ce897` already has **36** exported `*Store` methods, including `ObjectCommits`, `ObjectReferences` and `ObjectsBySemanticID` (added after `d8db0eb`) and M1's 3 `…Context` variants. M2's `CommitLanded` makes **37** (V12) | Not in scope. M7's planner must re-derive the count. The reflect-generated table absorbs it, and `CommitLanded` already takes ctx first |
| P12 | *(not stated)* The daemon cannot reach `store`'s package-global hooks or `s.db` from `host/daemon` tests | Confirmed. The unexported `durableAfterCommitHook` and `s.db` are invisible to `package daemon`. The existing pattern is an embed-and-override seam (`d.reads readStore`, `daemon.go:365`) | M2 adds a parallel **`d.commits durableStore`** seam. AC5 over HTTP uses a wrapping fake that calls the **real** `CommitContext`. The real store-side mechanism is M1's already-landed tests |

**Two defects found while prototyping. Both are fixed in the plan.**

1. **The fake has to model an ended context.** In the first draft of the AC5 fake, the uncertain arms returned `*UncertainError` at once. At that point the handler's ctx is still live, so a mutant that consults `timedOut` before `IsUncertain` (MUT-M2-TIMEOUT-BEFORE-UNCERTAIN) **survived** (rc 0). Every real uncertain outcome carries an ended ctx, because the caller stopped waiting. The fake now waits for `ctx.Done()` before returning, and the mutant is killed (rc 1). This is the same class of defect as iter-92's E7: a fake that hands the surviving branch exactly what it needs.
2. **Test fakes that embed `*store.Store` silently inherit a new `readStore` method.** Adding `GetReceiptContext` to `readStore` compiles, but `failingStore` then passes the receipt read through to the real store, and `TestInternalErrorsAreSanitized/read-routes` went red (`status = 200, want 500`). M3 adds explicit `GetReceiptContext` overrides to `blockingStore`, `recordingStore` and `failingStore`, and a `receipt` row to `seedReadRoutes`. The receipt route therefore inherits every read-deadline, cancel-release and sanitisation property (`MUT-M3-RECEIPT-UNBOUNDED` is killed by `TestDaemonReadDeadline`).

**Found, not fixed (outside M2–M4):** QUICKSTART §2 runs `/tmp/ailang-worldd commit --file /tmp/genesis.json` **without `--session`** (`QUICKSTART.md:56,76`), but `POST /v1/commit` is session-gated (`daemon.go:664`). By construction it therefore answers 401 `SessionAbsent`. This is an existing S7 drift. It was established by reading the code, not by running the daemon. See OPEN-5.

## 1. Gates and measurements

For **every** Go command: `export PATH=/opt/homebrew/bin:$PATH AILANG_BIN=$HOME/.pinned-ailang/ailang GOCACHE=$HOME/.cache/go-build TMPDIR=$HOME/.ailang/state/world-iter206/<role>-tmp; mkdir -p $TMPDIR; unset AILANG_REGISTRY_API_KEY`. The pinned binary is AILANG v0.41.0. The compile fence is **both** `go vet ./...` **and** `go test ./... -run '^$'`. CI runs `-race`, so every accept list below includes a `-race` run on the touched packages.

"Production code lines" counts added and deleted non-blank, non-comment lines in non-test `.go` files, per snapshot diff.

| Boundary | Cumulative state / production code lines | `go vet ./...` rc | `go test ./... -run '^$'` rc | `verify_ail.sh` rc |
|---|---|---:|---:|---:|
| base `b0ce897` | — | 0 | 0 | 0 (M1 CI) |
| M2 | `durable.go` `CommitLanded`; daemon seam, B6, two 503s; sketch arm; QUICKSTART. **+36 / −5** | 0 | 0 | 0 |
| M3 | opt-in `invocationId`, intent append, `GET /v1/receipts/{id}`; sketch route; QUICKSTART. **+86 / −13** | 0 | 0 | 0 |
| M4 | middleware B4 + `ResolveContext`; `Resolve` deleted; census pin removed. **+22 / −14** (plus `resolver.go` −37 comment/code lines) | 0 | 0 | n/a (no `.ail`) |

The design estimates were ~110, ~90 and ~40. All three milestones are well under 150.

**Full suite at M4:** `go test ./... -count=1` gave **23 ok, 1 FAIL**. The failure is `host/pkgproj` `TestQueryInterfaceReturnsWhileADescendantHoldsStdout` (`iface_test.go:232: instrument: descendant pid not recorded`). Re-run alone it passes (`ok 12.965s`). `go list -deps ./host/pkgproj` contains none of `host/store`, `host/daemon` or `host/authority` (the grep rc is 1, and the `ailang-world` deps count is 2). This is the known row-113 instrument flake, the same one as design V15, and it is not an acceptance gate.

**`-race`:** `go test -race ./host/store ./host/daemon ./host/authority ./host/projection ./cmd/ailang-worldd -count=1` gave 5/5 ok. The targeted selector ran `-race -count=5` ok, with **120 `=== RUN`** lines (= 24 × 5: 9 store + 15 daemon).

**Sandbox:** every new daemon test uses `httptest.NewRecorder` and `d.Handler().ServeHTTP`. **No socket is opened**, so the sandbox does not make these results uninformative. No `httptest.Server` or listening test is part of these accept lists.

**Package-global hooks:** M2–M4 add **none**. The daemon fakes are instance fields (`d.commits`, `d.resolver`, `d.reads`, `d.commitBudget`, `d.credentialBudget`) on a per-test daemon. The only edits to hook-touching tests are M2's two call-site swaps in `durable_test.go` (`logWitness` → `CommitLanded`). Those tests keep M1's restore-after-join ordering unchanged: `releaseOnce(release); s.workers.Wait(); hook = prev`. The executor must not reorder those defers.

## 2. Sub-milestones in compile order

At **each** boundary run the compile fence, then the accept list. `.ail` changes in M2 and M3 also need `./scripts/verify_ail.sh`. Pre-existing tests stay green, including the full `host/daemon` package: its tests exercise every read route through `seedReadRoutes`.

### M2 — `/v1/commit` bounded by B6; two distinct 503s; INV-LOG as `Store.CommitLanded` (AC5)

**Files:** `host/store/durable.go`, `host/store/durable_test.go`, new `host/store/commit_landed_test.go`; `host/daemon/daemon.go`, `handlers.go`, `daemon_test.go`, new `host/daemon/commit_budget_test.go`; `design_docs/sketches/worlddapi.ail`; `docs/QUICKSTART.md`.

**Store (`durable.go`, append):** this promotes the M1 test helper `logWitness` verbatim in behaviour. `LogEntry` is comparable (every field is an `int64`, a `string` or a `hashref.HashRef{algo, digest string}`), so `==` compares every stored field.

```go
// CommitLanded is INV-LOG (design §2.5, ratified D-WORLD-40) ... [doc as in the patch]
func (s *Store) CommitLanded(ctx context.Context, c Commit) (bool, error) {
	got, ok, err := s.GetLogEntry(ctx, c.Entry.Header.EntryIndex)
	if err != nil {
		return false, err // unknown
	}
	if !ok {
		return false, nil // absent: not committed
	}
	return got == c.Entry, nil // equal: entry landed; different: not committed
}
```

In `durable_test.go`, delete `logWitness` and its `reflect` import, and replace its two call sites with `s.CommitLanded(bctx, a)` and `again.CommitLanded(ctx2s(t), …)`.

**Daemon (`daemon.go`):**
- Constant `commitBudget = 3 * time.Second` (B6). Its doc names D-WORLD-40.
- Fields `commits durableStore`, `commitBudget time.Duration`.
- `type durableStore interface { CommitContext(ctx context.Context, c store.Commit) error }`.
- `New` wires `commits: s, commitBudget: commitBudget`.

**`handlers.go`:**
- Replace `:669-677` with:
  ```go
  ctx, cancel := context.WithTimeout(r.Context(), d.commitBudget)
  defer cancel()
  if err := d.commits.CommitContext(ctx, commit); err != nil { d.writeCommitError(w, r, ctx, commit, err); return }
  ```
- Add `writeCommitError`, whose switch order is: `ConflictError` → 409; `store.IsUncertain` → `writeCommitUncertain` (**before** `timedOut`, because an uncertain error also carries an ended ctx); `timedOut(ctx, err)` → `writeCommitTimeout`; default → `writeInternalError`.
- Add `writeCommitTimeout(w, budget)`: 503 `Timeout`, `commit deadline (%s) exceeded before the durable step: not committed; safe to resend`.
- Add `writeCommitUncertain(w, c store.Commit)`: 503 `CommitUncertain`, `commit outcome unknown: reconcile with GET /v1/log/%d before any resend: a row equal to your entry means it landed; absent or different means it did not; a different head does not mean it failed`.

Take `writeCommitError` and the `store.Commit`-typed `writeCommitUncertain` from the patch's final form now, so M3 only adds branches. Leave the 200 body unchanged.

**Sketch:** add `| CommitUncertain(int)` to `ApiError`, `((CommitUncertain(1)), 503)` to the `httpStatus` tests, and the arm `CommitUncertain(_) => 503`.

**QUICKSTART §2:** add one paragraph with the command `/tmp/ailang-worldd log get 0`. It says: a commit is bounded (3 s); `Timeout` means not committed and safe to resend; `CommitUncertain` means do not resend and reconcile by reading `log get <index>` and comparing it with the sent `entry` (equal means landed; 404 or different means not landed); the head is not evidence. The prototype wording is in the patch.

**Tests.** `TestCommitLandedReadsTheLogRowAtItsIndex` covers four arms, each a distinct branch:
- `absent_row_is_not_committed`;
- `equal_row_is_entry_landed`;
- `same_hash_different_header_is_not_committed`: A′ = A with `WrittenBy` changed, planted through `CommitContext`, which is possible because the store never recomputes `EntryHash`;
- `read_error_is_unknown`: an expired ctx.

`TestCommitBudgetAndUncertainReconcile` uses the wrapping fake `commitSeam{real *store.Store; mode}`:
- `block`: wait for `ctx.Done()`, then the real `CommitContext(ctx, …)`. That is a definite not-committed from the real store. There is a 2 s escape into the real commit, so a ctx-ignoring mutant turns red on status and elapsed time rather than hanging.
- `land-then-uncertain`: a real commit under its own 1 s ctx, then wait for `ctx.Done()`, then return `&store.UncertainError{…}`.
- `uncertain-not-landed`: wait for `ctx.Done()`, then return `UncertainError`.

Its arms:
- (i) held → 503 `Timeout` within budget + 500 ms; the message contains "not committed"; the head is unchanged.
- (ii) A lands and reports uncertain → 503 `CommitUncertain` naming `/v1/log/1`. B commits on A's world. The head ≠ A's world, `CommitLanded(A)` = true, and the wire form `GET /v1/log/1` equals `logJSON(A.Entry)`.
- (iii) A rolls back; B takes index 1 from genesis → false, and the wire row differs.
- (iii-b) A rolls back and nothing else happens → false, and the wire answer is 404.
- (iv) `d.readDeadline = expiredReadDeadline`; `GET /v1/log/1` → 503 `Timeout` whose message does not contain "resend".

`TestCommitUncertainStatusMirrorsSketch` checks that `writeCommitUncertain`'s status equals the **parsed** sketch arm. `TestBoundedWaitsAndBodyLimit` gains a constants row `{"commitBudget", commitBudget, 3 * time.Second}` and the wiring assertions `d.commitBudget == commitBudget` and `d.commits == durableStore(d.store)`.

**Accept:**
```
go test ./host/store -run '^(TestCommitLandedReadsTheLogRowAtItsIndex|TestUncertainBareCommitReconcilesByLogIndexNotHead|TestCloseKeepsWriterLockUntilDurableWorkerSettles|TestProductionContextRoots)$' -count=1 -v   # expect 9 store === RUN + census
go test ./host/daemon -run '^(TestCommitBudgetAndUncertainReconcile|TestCommitUncertainStatusMirrorsSketch|TestTimeoutStatusMirrorsSketch|TestBoundedWaitsAndBodyLimit|TestInternalErrorsAreSanitized|TestFrozenV1RouteTableMatchesMux)$' -count=1 -v
go test -race ./host/store ./host/daemon -count=1
./scripts/verify_ail.sh        # rc 0; sketch tests 20 -> 21; pins unmoved (16 / 40)
```
Count the `=== RUN` lines; do not trust rc alone. **Mutations:** M2-01 … M2-12 in §3.

### M3 — opt-in `invocationId` on `/v1/commit`; `GET /v1/receipts/{id}` (AC5b)

**Files:** `host/daemon/daemon.go`, `handlers.go`, `route_table_test.go`, `read_deadline_test.go`, `commit_budget_test.go`; `design_docs/sketches/worlddapi.ail`; `docs/QUICKSTART.md`.

**`daemon.go`:**
- `durableStore` gains `AppendIntentContext(ctx, id string, intent store.JournalIntent) (int64, hashref.HashRef, error)`.
- `readStore` gains `GetReceiptContext(ctx, id string) (store.Receipt, bool, error)`. `*store.Store` satisfies both.
- Register `mux.HandleFunc("GET /v1/receipts/{id}", d.handleReceipt)`.
- Update the "nine /v1 patterns" comment to "ten … (nine GET, one POST; GET /v1/receipts/{id} added by row 23 M3, D-WORLD-40)".

**`handlers.go`** (exact code in the patch):
- `commitRequest` gains `InvocationID string \`json:"invocationId,omitempty"\``. With `omitempty`, existing request bodies and `encodeCommit` are byte-unchanged.
- `const restInvocationPrefix = "rest:"`; `restInvocationID(id)` returns `HasPrefix(id, "rest:") && len(id) > 5`.
- `commitIntent(c)` derives the `JournalIntent` from the commit, with `LogicalTime = c.Entry.Header.EntryIndex` (design §2.5).
- In `handleCommit`, after `decodeCommit`: a non-empty id outside `rest:` gets 400 `BadRequest` "invocationId must be rest:<id>". Otherwise set `commit.InvocationID`.
- Under the **same** B6 ctx, when an id is present, first run `d.commits.AppendIntentContext(ctx, id, commitIntent(commit))`, then `CommitContext`. Both errors go to `writeCommitError`.
- `writeCommitError` adds a case: `*store.DuplicateInvocationError` or `*store.InvocationMismatchError` → 400 `BadRequest` "invocationId already names a different commit: …". This is the client's own input, so it is not sanitised.
- `writeCommitUncertain` adds a first branch: when `c.InvocationID != ""`, it says "reconcile with GET /v1/receipts/<id> …: resolved means it landed; not-started or indeterminate means it did not".
- `handleReceipt`: refuse non-`rest:` ids (400). Otherwise read under `d.readCtx(r)`, `GetReceiptContext`, and map `timedOut` → `writeReadTimeout`, error → `writeInternalError`. Return 200 `receiptResponse{invocationId, state, resultRef?}`, where `resultRef` = `Outcome.ResultRef` when resolved.

**Sketch:** append `{ method: "GET", path: "/v1/receipts/{id}" }` to `routes()`. **Route test:** change `9` → `10` in both length checks and in the message.

**`read_deadline_test.go`:**
- Add `{"receipt", "/v1/receipts/rest:" + label, "GetReceiptContext"}` to `seedReadRoutes`.
- Add `GetReceiptContext` overrides: `blockingStore` blocks on `b.block(ctx)`, `recordingStore` notes and delegates, `failingStore` returns `errSentinelInternal`.

**QUICKSTART §2:** add the `invocationId` paragraph and `curl -s http://127.0.0.1:7644/v1/receipts/rest:<your-id>`. The CLI has no `receipt` verb (OPEN-4).

**Tests.** `TestCommitInvocationReceiptIdentity`, with the fake extended to pass `AppendIntentContext` through (`block` mode holds it too):
- `A_landed_uncertain_then_B_on_top`: 503 `CommitUncertain` naming `/v1/receipts/rest:A`. B lands. The receipt is `resolved` and `resultRef` = A's world.
- `A_rolled_back_uncertain`: the receipt is `indeterminate`.
- `held_connection_on_the_intent_step_is_503_timeout`: 503 `Timeout` within budget + 500 ms, and the receipt is `not-started`.
- `foreign_namespaces_are_refused`: for `effect:ep:1`, `a2a:ep:task` and `rest:`, both the POST and the GET answer 400 with a message naming `rest:<id>`, and the head does not move.

**Accept:**
```
go test ./host/daemon -run '^(TestCommitInvocationReceiptIdentity|TestCommitBudgetAndUncertainReconcile|TestFrozenV1RouteTableMatchesMux|TestDaemonReadDeadline|TestReadCtxCancelledAfterHandler|TestInternalErrorsAreSanitized|TestCommitUncertainStatusMirrorsSketch)$' -count=1 -v
go test ./host/daemon -count=1 && go test -race ./host/daemon ./host/store -count=1
./scripts/verify_ail.sh        # rc 0
```
**Mutations:** M3-01 … M3-09 in §3.

### M4 — credential lookup bounded by B4; lookup expiry is 503, never 401; `Resolve` deleted (AC6)

**Files:** `host/daemon/middleware.go`, `daemon.go`, `daemon_test.go`, new `host/daemon/credential_budget_test.go`; `host/authority/resolver.go`, `resolver_test.go`, `mint_test.go`, `resolvecontext_test.go`; `host/projection/projection_test.go`; `cmd/ailang-worldd/session_test.go`; `host/store/context_roots_test.go`.

**`middleware.go`:**
- `SessionMiddleware` gains `budget time.Duration` and `fail func(http.ResponseWriter, *http.Request, error)`.
- `NewSessionMiddleware(resolver, budget, fail)`. Its only caller is `daemon.go` `Handler`, measured (V14).
- `Wrap` becomes:
  ```go
  lctx, cancel := context.WithTimeout(r.Context(), m.budget)
  out, err := m.resolver.ResolveContext(lctx, header, time.Now().Unix())
  timeout := err != nil && timedOut(lctx, err)
  cancel()
  if err != nil {
      if timeout { writeAPIError(w, "Timeout", fmt.Sprintf("credential lookup deadline (%s) exceeded", m.budget), http.StatusServiceUnavailable); return }
      m.fail(w, r, err); return
  }
  ```
  The binding still goes into `r.Context()`, not `lctx`. `lctx` is cancelled before the handler runs, so a commit never inherits the lookup budget.

**`daemon.go`:** constant `credentialBudget = 3 * time.Second` (B4); a field; `New` wires it; `Handler` returns `NewSessionMiddleware(d.resolver, d.credentialBudget, d.writeInternalError).Wrap(…)`.

**`authority/resolver.go`:**
- Delete `Resolve` from the `Resolver` interface and from `*resolver`.
- Move the header/now contract text onto `ResolveContext`'s interface doc.
- Reword three comments that say "Resolve".
- `ResolveContext`'s behaviour is unchanged.

**Census:** delete `"host/authority/resolver.go|Resolve|Background": 1`. The middleware introduces no root, because its parent is `r.Context()`.

**Test migration (mechanical, 24 call sites, V7).** Add a `resolveNow(t, r, header, now)` helper in `authority/resolver_test.go` and in `cmd/ailang-worldd/session_test.go`. It calls `ResolveContext` under a 5 s deadline and fails on error. Using a deadline now keeps these tests valid under M7's guard.
- Replace every `res.Resolve(`, `New(st).Resolve(` and `authority.New(st).Resolve(` in those files, plus `resolvecontext_test.go`: 15 + 6 + 2 + 1 sites.
- In `TestResolveContext_StoreErrorIsErrorNotDenial`, delete the block asserting that `Resolve` collapses a store error to `DenialUnknown`. That is the behaviour this milestone removes.
- In `projection_test.go`, delete `blockingResolver.Resolve`. It is dead once the interface loses the method.

**Tests.** `TestCredentialLookupExpiryIsNotADenial` uses `heldCredentials{real *store.Store}`, a `CredentialStore` that waits for `ctx.Done()`, with a 2 s escape into the real store, where a valid token resolves:
- `held_lookup_is_503_timeout_within_budget`: a **valid** minted token with a 100 ms budget → 503 `Timeout` "credential lookup …" within budget + 500 ms; the head is unchanged.
- `failed_lookup_is_sanitized_500`: `heldCredentials{err: "private lookup detail"}` → 500 `Internal`; the body omits the detail and the error log carries it.

`TestBoundedWaitsAndBodyLimit` gains `{"credentialBudget", credentialBudget, 3 * time.Second}` and a wiring assertion. The existing 401/400 denial tests (`session_middleware_test.go`) must stay green unchanged. They are the negative control: a responsive store still answers 401 `SessionUnknown`.

**Accept:**
```
go test ./host/daemon -run '^(TestCredentialLookupExpiryIsNotADenial|TestBoundedWaitsAndBodyLimit|TestSessionMiddleware_.*)$' -count=1 -v
go test ./host/authority ./host/projection ./cmd/ailang-worldd -count=1
go test ./host/store -run '^TestProductionContextRoots$' -count=1 -v
go test -race ./host/daemon ./host/authority ./host/projection ./cmd/ailang-worldd ./host/store -count=1
```
**Mutations:** M4-01 … M4-07 in §3.

## 3. Runnable mutation ledger (28 rows, all executed on the prototype)

Harness: `~/.ailang/state/world-iter206/planner-scratch/mutate.py <ledger.json> [IDs]`. For each row it requires the exact old text to occur exactly once, applies the edit, runs `go test <pkg> -run '<selector>' -count=1 -v -timeout 60s`, records rc and the `=== RUN` count, restores the file bytes, and asserts a sha256 match. Transcripts are in `…/planner-scratch/mut/<ID>.txt`. The ledgers are `m2_muts.json`, `m3_muts.json`, `m4_muts.json` and `m2_final_forms.json`. The final full re-run of all ledgers on the M4 tree is `final_mut_run.txt`. Every rc below is from that re-run, or from `m2_final_forms.json` for the two rows marked @M3. Each row's edit is shown in its **final-tree form**. **rc 1 with `=== RUN` ≥ 1 is a kill.** No row was killed by a compile error: MUT-M3-ANY-NS first produced an unused-import compile failure (`RUN=0`, so vacuous) and was rewritten until it compiled and killed.

| ID | File: exact old → new | Red test (selector) | rc | RUN |
|---|---|---|---:|---:|
| M2-01 MUT-M2-NOCTX | `handlers.go`: `if err := d.commits.CommitContext(ctx, commit); err != nil {` → `if err := d.store.Commit(commit); err != nil {` | `^TestCommitBudgetAndUncertainReconcile$/^i_held` (200, want 503) | 1 | 2 |
| M2-02 MUT-M2-NOBUDGET | `handlers.go`: `ctx, cancel := context.WithTimeout(r.Context(), d.commitBudget)` → `ctx, cancel := context.WithCancel(r.Context())` | same (200 after the 2 s escape) | 1 | 2 |
| M2-03 MUT-M2-UNCERTAIN-AS-TIMEOUT @M3 | `handlers.go` `writeCommitError`: `\t\twriteCommitUncertain(w, c)\n` → `\t\twriteCommitTimeout(w, d.commitBudget)\n` | `…/^ii_` (class Timeout, want CommitUncertain) | 1 | 2 |
| M2-04 MUT-M2-TIMEOUT-BEFORE-UNCERTAIN | `handlers.go`: `case store.IsUncertain(err): // before timedOut: …` → `case store.IsUncertain(err) && !timedOut(ctx, err):` | `…/^ii_` | 1 | 2 |
| M2-05 MUT-M2-HEAD-WITNESS (HTTP) | `durable.go` `CommitLanded` body → `head, _, err := s.SelectedHead(ctx); if err != nil { return false, err }; return head == c.NextWorld.Ref, nil` | `…/^ii_` (`CommitLanded(A) = false; want true`) | 1 | 2 |
| M2-06 MUT-M2-HEAD-WITNESS (store) | same edit | `./host/store` `^TestUncertainBareCommitReconcilesByLogIndexNotHead$` | 1 | 3 |
| M2-07 MUT-M2-HASH-ONLY | `durable.go`: `\treturn got == c.Entry, nil` → `\treturn got.EntryHash == c.Entry.EntryHash, nil` | `./host/store` `^TestCommitLandedReadsTheLogRowAtItsIndex$` (same_hash arm) | 1 | 5 |
| M2-08 MUT-M2-ABSENT-AS-UNKNOWN (store) | `durable.go`: `if !ok { return false, nil }` → `if !ok { return false, fmt.Errorf("store: log entry %d absent: outcome unknown", c.Entry.Header.EntryIndex) }` | `./host/store` `^TestCommitLandedReadsTheLogRowAtItsIndex$` (absent arm) | 1 | 5 |
| M2-09 MUT-M2-ABSENT-AS-UNKNOWN (HTTP) | same edit | `./host/daemon` `^TestCommitBudgetAndUncertainReconcile$/^iii_b` | 1 | 2 |
| M2-10 MUT-M2-ERR-AS-NOTCOMMITTED | `durable.go`: in the `err != nil` arm, `return false, err` → `return false, nil` | `./host/store` `^TestCommitLandedReadsTheLogRowAtItsIndex$` (read_error arm) | 1 | 5 |
| M2-11 MUT-M2-BUDGET-UNWIRED | `daemon.go` `New`: `commits: s, commitBudget: commitBudget,` → `commits: s, commitBudget: readDeadline,` | `^TestBoundedWaitsAndBodyLimit$` | 1 | 6 |
| M2-12 MUT-M2-SKETCH-DRIFT @M3 | `handlers.go`: `"a different head does not mean it failed", c.Entry.Header.EntryIndex), http.StatusServiceUnavailable)` → `…, http.StatusGatewayTimeout)` | `^TestCommitUncertainStatusMirrorsSketch$` | 1 | 1 |
| M3-01 MUT-M3-NO-INTENT | `handlers.go`: delete `\t\tcommit.InvocationID = request.InvocationID\n` | `^TestCommitInvocationReceiptIdentity$/^A_landed` | 1 | 2 |
| M3-02 MUT-M3-SKIP-APPEND | `handlers.go`: `\tif commit.InvocationID != "" {\n\t\t// Both durable steps share B6.` → `\tif false && commit.InvocationID != "" {…` | same (400: the store's bind refuses an unjournalled id) | 1 | 2 |
| M3-03 MUT-M3-EFFECT-NS | `handlers.go` `restInvocationID`: `return strings.HasPrefix(id, restInvocationPrefix) && len(id) > len(restInvocationPrefix)` → `return (strings.HasPrefix(id, restInvocationPrefix) \|\| strings.HasPrefix(id, "effect:")) && len(id) > len(restInvocationPrefix)` | `…/^foreign` (message no longer names `rest:<id>`) | 1 | 2 |
| M3-04 MUT-M3-ANY-NS | same function: → `return strings.TrimSpace(id) != ""` | `…/^foreign` | 1 | 2 |
| M3-05 MUT-M3-APPEND-NOBUDGET | `handlers.go`: `d.commits.AppendIntentContext(ctx, commit.InvocationID` → `d.commits.AppendIntentContext(r.Context(), commit.InvocationID` | `…/^held` (answered after 2.0 s) | 1 | 2 |
| M3-06 MUT-M3-UNCERTAIN-LOG-ROUTE | `handlers.go` `writeCommitUncertain`: `\tif c.InvocationID != "" {\n\t\twriteAPIError(w, "CommitUncertain"` → `\tif false && c.InvocationID != "" {…` | `…/^A_landed` | 1 | 2 |
| M3-07 MUT-M3-RECEIPT-RESULTREF | `handlers.go`: `\t\tbody.ResultRef = rc.Outcome.ResultRef.String()` → `\t\tbody.ResultRef = rc.IntentRef.String()` | `…/^A_landed` | 1 | 2 |
| M3-08 MUT-M3-RECEIPT-UNBOUNDED | `handlers.go` `handleReceipt`: `\tctx, cancel := d.readCtx(r)\n\tdefer cancel()\n\trc, _, err := d.reads.GetReceiptContext(ctx, id)` → `\tctx, cancel := context.WithCancel(r.Context())\n…` | `^TestDaemonReadDeadline$` (receipt route 200, want 503) | 1 | 4 |
| M3-09 MUT-M3-ROUTE-SKETCH-DRIFT | `worlddapi.ail`: delete `,\n    { method: "GET", path: "/v1/receipts/{id}" }` | `^TestFrozenV1RouteTableMatchesMux$` | 1 | 1 |
| M4-01 MUT-M4-COLLAPSE | `middleware.go`: `\t\tif err != nil {\n\t\t\t// A lookup that did not finish is not a denial: never 401.` → insert `writeSessionDenial(w, authority.DenialUnknown); return` as the first statement | `^TestCredentialLookupExpiryIsNotADenial$/^held` (401, want 503) | 1 | 2 |
| M4-02 MUT-M4-COLLAPSE-500 | `middleware.go`: `\t\t\tm.fail(w, r, err)\n` → `\t\t\twriteSessionDenial(w, authority.DenialUnknown)\n` | `…/^failed` (401, want 500) | 1 | 2 |
| M4-03 MUT-M4-NOCTX | `middleware.go`: `lctx, cancel := context.WithTimeout(r.Context(), m.budget)` → `lctx, cancel := context.WithCancel(context.Background())` | `…/^held` (200 after the escape) | 1 | 2 |
| M4-04 MUT-M4-NOBUDGET | same line → `lctx, cancel := context.WithCancel(r.Context())` | `…/^held` | 1 | 2 |
| M4-05 MUT-M4-TIMEOUT-AS-500 | `middleware.go`: `timeout := err != nil && timedOut(lctx, err)` → `timeout := false` | `…/^held` (500, want 503) | 1 | 2 |
| M4-06 MUT-M4-BUDGET-UNWIRED | `daemon.go` `New`: `credentialBudget: credentialBudget,` → `credentialBudget: readDeadline,` | `^TestBoundedWaitsAndBodyLimit$` | 1 | 6 |
| M4-07 MUT-M4-RESTORE-BG-ROOT | `resolver.go`: before `// ResolveContext implements Resolver …`, insert `func (r *resolver) resolveBG(h string, n int64) (ResolveOutcome, error) { return r.ResolveContext(context.Background(), h, n) }` | `./host/store` `^TestProductionContextRoots$` (Background=22) | 1 | 1 |

Rows M2-03 and M2-12 were also killed at the M2 snapshot in their M2-era forms. Those forms had the inline switch and `writeCommitUncertain(w, entryIndex int64)`. The plan asks the executor to land the `store.Commit`-typed helpers in M2, so the final-tree forms above apply from M2 on. M2-01's red is `status = 200, want 503` because the mutant bypasses the seam and the real store commits at once. M4-03 and M4-04 are red on status (200) **and** elapsed time (2 s escape), never a hang.

## 4. Execution notes

- **Land the fakes exactly as prototyped.** The uncertain arms of `commitSeam` **must** wait for `ctx.Done()` (defect 1 in §0). All three store-embedding read fakes **must** override `GetReceiptContext` (defect 2 in §0). Do not replace the wrapping fakes with fakes that never call the real store: every "block" arm ends in a real `CommitContext`, `AppendIntentContext` or `ResolveSession` call, which is what keeps the not-committed classification real.
- **Why a seam and not the real hook.** `host/daemon` cannot set `store`'s unexported `durableAfterCommitHook` or hold `s.db`. Adding exported test hooks to `store` would put test-only API on the surface M7 guards. The real post-COMMIT mechanism is already proven at store level by M1's `TestUncertainBareCommitReconcilesByLogIndexNotHead` and `TestDriverCommitErrorAfterRealCommitIsUncertainAndReconciles`. M2's HTTP test proves the mapping and INV-LOG over the wire.
- **Budgets compose.** A protected commit can spend B4 (3 s) on the lookup and then B6 (3 s) on the commit. That is ≤ 6 s against `writeTimeout` = 30 s, so the 503 is always writable. `lctx` is cancelled before the handler runs.
- **Sketch edits** are Leg-1 `ai-check` only. Leg 2 tests only `world/`, so `EXACT_TOTAL_TESTS` = 40 and `EXACT_TOTAL_VERIFIED` = 16 are unmoved. There is no packaged-module or floor-raise coupling: `worlddapi.ail` is not packaged. Run `verify_ail.sh` at the M2 and M3 boundaries anyway.
- **QUICKSTART (S7).** The new reconcile commands (`log get 0`, the `curl` receipt read) are ordinary reads and can be run verbatim. The 503 paths cannot be produced on demand against a live daemon. Their prose is justified by the AC5/AC5b tests, not by a verbatim run, and the evaluator should read it that way.
- **Conflict surface (§7), re-checked.** Row 111's 503 read-path mapping touches `handlers.go` `timedOut`/`writeReadTimeout`; M2 **reuses** `timedOut` and does not modify it. Row 106 M5 touches `daemon.go` to construct the coordinator; M2–M4 add fields and a `New` line there, which is a rebase at most.
- **Mutation discipline.** Re-run the §3 ledger on the landed files. Restore bytes by copy and sha256, never with `git checkout --`. Count `=== RUN`. A `-run 'A\|B'` pattern runs zero tests; use `'^(A|B)$'`.

## 5. Verification Log

All commands were run at `/Users/voightkampff/dev/sunholo-data/.planner-wt-iter206` (`b0ce897`) unless marked *proto*, which means the scratch copy `~/.ailang/state/world-iter206/planner-scratch/proto` (a `git archive b0ce897` export plus the prototype).

| # | Claim | Command | Observed |
|---|---|---|---|
| V1 | Commit call and 200 body lines | `grep -n 'd.store.Commit(commit)\|func (d \*Daemon) handleCommit\|writeJSON(w, http.StatusOK, commitResponse' host/daemon/handlers.go` | `615` handler, `669` `d.store.Commit(commit)`, `678` 200 body = `commit.NextWorld.Ref` |
| V2 | Middleware and resolver root | `grep -n 'm.resolver.Resolve(' host/daemon/middleware.go`; `grep -n 'func (r \*resolver) Resolve(\|ResolveContext(context.Background()' host/authority/resolver.go` | `middleware.go:49`; `resolver.go:128`, `:129` |
| V3 | Log route exists; receipts route absent (positive control in the same probe) | `grep -n 'mux.HandleFunc("GET /v1/log/{index}"\|receipts' host/daemon/daemon.go`; `grep -rn receipts host/daemon/*.go \| grep -v _test \| wc -l` | `daemon.go:644` present; receipts **0** |
| V4 | Frozen table pin | `grep -n 'len(mux) != 9' host/daemon/route_table_test.go`; `grep -n 'method: "POST", path: "/v1/commit"' design_docs/sketches/worlddapi.ail` | `route_table_test.go:51`; `worlddapi.ail:81` (last entry of `routes()`) |
| V5 | Read-timeout text is read-specific | `grep -n 'func writeReadTimeout\|"Timeout", fmt.Sprintf' host/daemon/handlers.go` | `:338-339` `read deadline (%s) exceeded` |
| V6 | Census pin M4 removes | `grep -n 'authority/resolver.go\|Resolve\|Background' host/store/context_roots_test.go` | `:29` |
| V7 | `Resolve` test call sites | `grep -rn '\.Resolve(' --include='*_test.go' host/authority cmd/ailang-worldd \| wc -l`; `grep -n 'func (blockingResolver) Resolve' host/projection/projection_test.go` | **25** (24 migrated + 1 deleted with its collapse block); `projection_test.go:194` |
| V8 | QUICKSTART commit runs without a session; `/v1/commit` is protected | `grep -n 'commit --file /tmp/genesis.json' docs/QUICKSTART.md`; `grep -n 'r.URL.Path == "/v1/commit"' host/daemon/daemon.go` | `:56`, `:76` (no `--session`); `daemon.go:664`. **Code reading, not executed** |
| V9 | Namespaces | `grep -n 'func InvocationID' host/coordinator/coordinator.go`; `grep -n 'strings.HasPrefix(id, "effect:")' host/store/journal.go` | `a2a:` at `:100`; `effect:` refused at `:325`, `:837` |
| V10 | A resolved invocation is an idempotent no-op before the head compare | `grep -n 'A fully matching resolved invocation is an idempotent no-op' host/store/store.go` | `:1003` (feeds OPEN-3) |
| V11 | Three read fakes embed `*store.Store` | `grep -n 'type blockingStore struct\|type recordingStore struct\|type failingStore struct' -A1 host/daemon/read_deadline_test.go \| grep -c store.Store` | **3**. *proto*: before the override, `TestInternalErrorsAreSanitized/read-routes` → `read_deadline_test.go:811: status = 200, want 500` |
| V12 | Exported `*Store` method count | `ls host/store/*.go \| grep -v _test \| xargs grep -h '^func (s \*Store) [A-Z]' \| wc -l` (worktree, then *proto*) | **36** → **37** |
| V13 | Sketch edits pass the pinned gate | *proto*: `./scripts/verify_ail.sh` after M2 and after M3; `cd design_docs && $AILANG_BIN test sketches/worlddapi.ail` against the base export | rc 0 both times, `16 required identities verified, 40 named tests pass`; worlddapi tests **20 → 21** (`httpStatus_test_7`) |
| V14 | `NewSessionMiddleware` callers | `grep -rn 'NewSessionMiddleware' --include='*.go' .` | the definition plus `daemon.go:655` only |
| V15 | Boundary gates | *proto*: per snapshot, `go vet ./...` and `go test ./... -run '^$'` | M2, M3, M4 all `vet=0 compile=0` |
| V16 | Production size | *proto*: `git diff <a> <b> -- '*.go' ':!*_test.go'`, counting `^[+-]` non-blank, non-comment lines | M2 +36/−5; M3 +86/−13; M4 +22/−14 |
| V17 | Race and stability | *proto*: `go test -race ./host/store ./host/daemon ./host/authority ./host/projection ./cmd/ailang-worldd -count=1`; targeted `-race -count=5 -v \| grep -c '^=== RUN'` | 5/5 ok; **120** RUN, ok |
| V18 | Full suite | *proto*: `go test ./... -count=1`; then `go test ./host/pkgproj -count=1`; `go list -deps ./host/pkgproj \| grep -E 'host/(store\|daemon\|authority)$'` | 23 ok, 1 FAIL (`pkgproj` `iface_test.go:232`), `ok` alone, deps-grep rc 1 |
| V19 | Mutations | *proto*: `python3 mutate.py all_muts.json`; `python3 mutate.py m2_final_forms.json` | 28 rows plus the 2 final-tree re-forms, all rc 1 with RUN ≥ 1 (§3). First-draft survivor MUT-M2-TIMEOUT-BEFORE-UNCERTAIN rc **0** before the fake fix, rc 1 after |
| V20 | Patch applies to base | `git apply --check ~/.ailang/state/world-iter206/planner-scratch/prototype_m2_m4.diff` | rc 0; worktree `git status --porcelain` empty before this plan was written |

## 6. OPEN questions

- **OPEN-1 — the frozen /v1 table grows from 9 to 10 routes.** §2.5 names `GET /v1/receipts/{id}`, and D-WORLD-40 A ratified "opt-in receipt identity". The design never says that this edits the frozen `worlddapi.ail` `routes()` and the route-table pin. The plan reads the ratification as authorising the edit, because no alternative route exists. **The controller should confirm that reading before M3 routes.** If it is refused, M3 cannot be delivered as designed. M2 and M4 are unaffected.
- **OPEN-2 — `CommitLanded` has no production caller.** The design names `store.CommitLanded(ctx, c)`. The daemon cannot use it without waiting out an uncertain outcome, and clients reconcile over `GET /v1/log/{index}`. It lands as the store's statement of INV-LOG: the in-process API that the AC5 tests and the wire rule are checked against, and that a future in-process caller would use. It is exported and ctx-first, so it is inside M7's guard. If the controller prefers no uncalled exported API, the alternative is to keep it as a test helper. That would move MUT-M2-05…10 onto test code, which weakens them.
- **OPEN-3 — replaying a resolved `rest:` commit returns a misleading 200.** A resend of an already-resolved `invocationId` is a store-level idempotent no-op (`store.go:1003`, before the head compare). The handler then answers 200 `{"selectedHead": <A's world>}` even when B has since moved the head, because the 200 body is `commit.NextWorld.Ref`, not a read. This predates this tranche but only becomes reachable in M3. Options: (a) leave it as a residual, since the design says not to resend until reconciled; (b) a follow-up that answers such a replay with the receipt rather than a head. The design does not decide, and this plan does not invent policy. The default is (a).
- **OPEN-4 — CLI parity.** The sketch comment says CLI subcommands map 1:1 onto routes, but no test enforces it. The design adds no CLI verb. The plan's QUICKSTART uses `curl` for `/v1/receipts`, and `invocationId` is set by editing the commit JSON. A `receipt get` verb is left out unless the controller wants it (~15 LOC in `cmd/ailang-worldd/cli.go`).
- **OPEN-5 — existing QUICKSTART drift (S7).** §2's `commit --file` omits `--session` against a session-gated `/v1/commit` (V8). M2 edits the same section but does not repair this. Fixing it requires re-ordering the walkthrough (minting a session needs the daemon stopped), which is out of scope. **A new queue row is needed** unless the controller rides it on M2.
- **OPEN-6 — non-timeout lookup errors.** B4 only specifies expiry → 503. The plan maps every other lookup error to the sanitised **500 `Internal`**, not to 401, on the ratified principle that a store failure is not an authentication failure. This is a planner reading. The design does not state it.

**Residuals (named owners, unchanged from design §10 unless new):**
- A `rest:` intent whose commit then conflicts (409) or times out pre-cutoff stays pending forever: its receipt is `indeterminate`. A retry needs a **new** id: the same id with a re-planned commit gets 400 `already names a different commit`. `broker.Recover` (no production caller) would list such intents as indeterminate. Owner: whoever wires `Recover`. **Not a new row now.**
- M7's method count (P11). Owner: row 23 M7.

## 7. Controller disposition of the OPEN questions (iteration 206)

None of these is a new policy question for Mark, so none is filed as a decision. Each follows either the quorum-passed design or a landed precedent.

- **OPEN-1 → proceed.** `GET /v1/receipts/{id}` is named in the design's M3 row (§5) and AC5b (§6.2), which passed the quorum and whose contracts D-WORLD-40 = A ratified. Precedent: row 96 grew the frozen table from 8 to 9 routes as ordinary loop work (iteration 201, `e859501`). M3 updates the sketch `routes()`, the pinned count and the `daemon.go` comment in one commit.
- **OPEN-2 → as designed.** `store.CommitLanded(ctx, c)` lands as the exported statement of INV-LOG, as §2.5 and the M2 row name it.
- **OPEN-3 → residual, new queue row.** The 200 reply to a resend of an already-resolved `rest:` commit predates this tranche, and fixing it needs a design choice (read the head, or answer with the receipt). It is filed as its own queue row at Gate 4, and this tranche does not change it.
- **OPEN-4 → no CLI verb.** The design adds none, and QUICKSTART uses `curl`.
- **OPEN-5 → new queue row.** The QUICKSTART `commit` step without `--session` is pre-existing S7 drift. It needs the walkthrough reordered, so it is filed as its own row at Gate 4. M2 must not make it worse.
- **OPEN-6 → accept the planner's reading.** The design's AC6 mutation `MUT-M3-COLLAPSE` (lookup error → `DenialUnknown`) must go red, so the design already treats the 401 collapse as the defect. A sanitised 500 still fails closed, because no request is served. This is reported to Mark as a default, not asked.
