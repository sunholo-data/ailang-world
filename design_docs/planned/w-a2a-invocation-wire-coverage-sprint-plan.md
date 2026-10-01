# Sprint plan — w-a2a-invocation-wire-coverage (iteration 213)

**Authority:** World mission queue row 123 (clause-6 residual of row 106, test-only). **Source finding:** the iteration-208 independent judge's report, §0 items 1–2 and §4.2 F1 (V13). **Parent design:** [w-transition-invocation-coordinator.md](w-transition-invocation-coordinator.md) and its [sprint plan](w-transition-invocation-coordinator-sprint-plan.md), M5 step 7 (V11). **Scope:** one milestone, about 0.25 day. **Landing:** one commit, one new test file plus this plan. **Production code: unchanged** (V15). The planner prototyped the full change in the worktree; every gate and mutation below ran on that prototype and nothing was committed.

Run this prefix before every command. It pins ailang v0.41.0 and keeps TMPDIR outside the repo:

```sh
source /Users/voightkampff/.ailang/state/mission-world-iter213-evidence/env.sh
```

## 1. Problem

On a dispatch error the `/a2a/` handler (base `131c517`, `host/projection/projection.go:339-350`) does three things, in this order (V2):

1. If the invocation context has ended (`ctx.Err() != nil`, `:345`), it sets `Connection: close` (`:346`).
2. It maps the error with `dispatchError(err)` (`:348`, defined at `:391`).
3. It writes the refusal with `protocol.A2AError(w, req.ID, code, msg)` (`:349`).

`TestA2ADispatch` (`projection_test.go:802`) tests step 2 only. Its 17 refusal rows call `dispatchError(fmt.Errorf("wrapped: %w", row.err))` directly (`:834`). Its three HTTP subtests (`:803-805`, V5) never assert the JSON-RPC id of an error. `registry_moved_after_admission` is the only one of them that reaches the dispatch-error write, and it discards the id (`code, msg, _ := a2aErr(…)`, `:1012-1013`). No test in `host/` reads the `Connection` response header (count 0, V6). Steps 1 and 3 therefore had no recorder-level test, and the measurement in §4 bears this out. These three mutants survive **all 432 pre-existing tests** in `host/projection`, `host/daemon` and `host/coordinator`:

- echoing a `null` id instead of the request id;
- always sending `Connection: close`;
- keying `Connection: close` on the error's type instead of on the context.

A fourth mutant, dropping the header entirely, is caught only by the loopback-socket test `TestA2AInvokeDeadlineClosesSocket`. That test is the one the iteration-208 judge could not run in its sandbox (F2).

The second gap is the reviewer's round-2 request. Parent plan M5 step 7 asks for tests "through the invocation path: hold the only connection beyond the request deadline, then stall `GetReceipt`, `AppendIntent`, and `Commit` separately … asserting bounded request completion, no accumulating blocked workers across repeated calls, definite R11 before cutoff and R16 for uncertain cutoff" (V11). The iteration-209 planner narrowed it to "a daemon/projection HTTP assertion with repeated calls". That exists today only at the `Dispatch` level (`TestDispatchDurableDeadline`, `coordinator_test.go:575`, 3 calls per op, V10) and at the store level (`TestDurableOpsHonourHeldConnection`, `durable_test.go:74`). Neither one runs an HTTP request.

### Interpretation of "repeated HTTP held-connection deadline proof" (measured)

The request originates in parent plan M5 step 7 (V11). The iteration-209 planner, the executor's admission and the judge's §0 item 1 all quote it (V11). Read together, the proof means: on **one handler**, issue **three** sequential `/a2a/` requests while one durable operation (`GetReceipt`, `AppendIntent` or `Commit`, each separately, plus the two uncertain variants) is held past the invocation deadline. Each request must:

- complete within the invoke bound;
- answer the exact R11 refusal (definite) or the exact R16 refusal (uncertain), with `Connection: close`;
- leave no blocked goroutine behind, and commit nothing.

After the three requests, a normal call on the same handler must succeed. This is the positive control that no in-flight, uncertain-intent or connection state leaked.

The *real* sole SQLite connection cannot be held from outside `host/store`. The only holder is `s.db.Conn` in the package-internal test (V9), and no exported `Store` method takes a callback (count 0 of 34, V9). The HTTP-level test therefore holds the operation at the coordinator's existing `coordinator.Store` seam (V4). This is the same seam and stall shape as `TestDispatchDurableDeadline`'s `stalledDurableStore` (V10). See OPEN-1.

## 2. The change (M1, prototyped): one new file, `host/projection/a2a_wire_test.go` (+414/−0)

Every refusal is produced by the **real `coordinator.Coordinator`** (built by `coordinator.New`) behind the **real `Handler.A2A`**, using the real resolver, a minted credential, the real registry reader and an in-memory store. A failure is injected only through the coordinator's own construction seams, `Store`, `Runner` and `BinderFor` (V4). `projection.Config.Coordinator` is a concrete `*coordinator.Coordinator` (V3), so a fake dispatcher would have required a production interface seam, and none was added.

Shared fixture: `newWireRig` (`:40`) commits a genesis world, stores canonical echo source and publishes `tools.echo` with no declared effects. It is the same shape as the existing `testA2ADispatchSuccess`. `assertWireRefusal` (`:100`) asserts, from the recorded response:

- HTTP **200** and `Content-Type: application/json`;
- `jsonrpc:"2.0"`, no `result` member, and an error object with **only** `code` and `message` (through the existing strict `a2aErr`, V16);
- the JSON-RPC **id echoed verbatim** (each row sends a distinct string id, for example `"wire-R11"`);
- the **exact** code and message;
- **no detail interpolation**: the row's hidden inner-error text plus the internal prefixes `coordinator:`, `store:`, `transition registry:` and `context deadline exceeded` must not appear anywhere in the body;
- the **Connection header**: `close` when the invocation deadline ended the request, **absent** otherwise. This pins the `ctx.Err()` condition from both sides.

Each wire row also asserts that its injection point was **reached exactly once**. Admission shares some messages: `not authorized` with an unlisted skill, and the not-available constant with the nil-coordinator branch. Without the reach check, a refusal raised before dispatch could pass as the row.

### Tests

**`TestA2ADispatchWire`** (`:193`) checks representative typed refusals on the wire:

| Subtest | Real failure source | Wire expectation | Connection |
|---|---|---|---|
| `R2_R3` | `BinderFor` returns a binder whose `Bind` fails with `*transitionreg.AccessDeniedError{Label: secret}`. `transitionreg.Bind` wraps it with `%w` (`bind.go:78`). | −32602 `not authorized` | absent |
| `R11` | the Runner blocks until the invocation context ends (`InvokeWait` 300 ms); the coordinator returns `coordinator: execute: <ctx err>` (`coordinator.go:304`) | −32603 `invocation exceeded its deadline` | **close** |
| `R14` | a **real store conflict**: the Runner moves the world head mid-run (`PutWorld` + `SelectHead`), so `store.Commit`'s compare-and-append returns `*store.ConflictError` (`store.go:1148`). The test confirms the head moved. | −32603 `world head moved during invocation; not committed; send a new task id` (no head hashes, no `stale observed`) | absent |
| `R16` | the Store's `Commit` returns `*store.UncertainError{Cause: secret}`; the coordinator wraps it in `UnconfirmedError` (`coordinator.go:328`) | −32603 `invocation outcome is not confirmed; resend the same task id` | absent |
| `R4_R6_default` | the Store's `GetObject` (source load) returns an untyped error carrying a secret, as `coordinator: load source: %w` (`:276`) | −32603 `transition invocation is not available in this daemon` | absent |

**`TestA2ADurableDeadlineWire`** (`:348`) is the repeated held-operation proof. A `stallStore` holds one named operation until the invocation context ends. A `release` channel frees it if the context never ends, so a context-dropping mutant fails instead of hanging. `InvokeWait` is 300 ms and the completion bound is `InvokeWait` + 3 s (the executor widened both from the prototype's 100 ms / +400 ms for loaded `-race` CI runners; the mutation verdicts below were measured at the prototype's bounds, and the executor re-ran A1, A2, B1, B2, B3, C-R14, C-R16 and D1 at the widened bounds: all KILLED).

| Subtest | Held op | Expected per call (×3 on one handler) |
|---|---|---|
| `GetReceipt` | GetReceipt | R11 + close |
| `AppendIntent` | AppendIntent | R11 + close |
| `Commit` | Commit | R11 + close |
| `uncertain_AppendIntent` | AppendIntent → `UncertainError` | R16 + close |
| `uncertain_Commit` | Commit → `UncertainError` | R16 + close |

Each call must return after at least `InvokeWait` (the stall was really held) and within the bound. After the three calls:

- the stall was reached exactly 3 times;
- `runtime.NumGoroutine()` settles back to the pre-loop baseline;
- the selected head is still genesis;
- with the stall disarmed, a normal call on the same handler returns a `completed` task with the id echoed and **no** `Connection` header.

The existing `TestA2ADispatch` direct table is **kept unchanged**. It still covers all 17 arms; the wire tests cover the write path and a representative arm set.

## 3. Acceptance criteria

| AC | Criterion | Command |
|---|---|---|
| AC-WIRE | the five wire rows pass, each asserting status, id, exact code and message, no interpolation, reach count and Connection | cmd 1 |
| AC-REPEAT | the five held-op rows pass: 3 bounded calls each, R11 or R16 + close, no goroutine growth, no commit, then a successful control | cmd 1 |
| AC-KILL | every mutation in §4 is KILLED by a new test and the file is sha1-identical after revert; A1, B2 and B3 survive every pre-existing test in the three packages | harness `mutate.py` + `survivors.py` |
| AC-GATES | vet, gofmt, race on projection/daemon/coordinator, the full suite and verify_ail are green | §5 |
| AC-TESTONLY | `git status --short` shows only the new test file and this plan; the `projection.go` sha1 equals HEAD's | V15 |

cmd 1 is kept outside the table so that the alternation is a plain RE2 `|`. A run with zero `=== RUN` lines is not green, and this command must show 12:

```sh
go test -race -count=1 -v -run '^(TestA2ADispatchWire|TestA2ADurableDeadlineWire)$' ./host/projection/ | grep -c '^ *=== RUN'
```

## 4. Mutation table (run on the prototype, against PRODUCTION code)

Harness: `/Users/voightkampff/.ailang/state/mission-world-iter213-evidence/mutate.py`. For each mutation it:

1. takes the sha1 of the file and backs it up to `$TMPDIR` (outside the repo);
2. asserts the anchor occurs exactly once and applies the change;
3. runs **NEW** (cmd 1's regex, `-count=1 -v`), then **OLD** (`-run '^TestA2ADispatch$'`, the pre-existing table with its HTTP subtests);
4. restores from the byte copy and re-checks the sha1.

Log: `…/mutations.log` (run 1, superseded after the R11 runner was bounded, is `…/mutations-run1-superseded.log`; see History). `RUN` is 12 for NEW and 21 for OLD on every row, so neither regex was vacuous. **All 15 restores are sha1-identical** (`projection.go` `48a9fb55…`, `coordinator.go` `8e25ffbe…`; 15× `restored=OK`, 0× `MISMATCH`).

**Survivor control.** `…/survivors.py` re-applies each mutation that the old table survives. It then runs the **entire pre-existing** projection, daemon and coordinator suites with the new tests skipped (`-skip '^(TestA2ADispatchWire|TestA2ADurableDeadlineWire)$'`, 432 `=== RUN`). Log: `…/survivors.log`.

| Mutation | File: change | NEW verdict (first failing assertion) | OLD table | All 432 pre-existing tests |
|---|---|---|---|---|
| **A1-wrong-id** | projection.go:349 `A2AError(w, req.ID, …)` → `A2AError(w, nil, …)` | **KILLED** 12/12 — `JSON-RPC id=null, want … "wire-R2R3"` | **SURVIVED** | **SURVIVED** (rc=0, 0 FAIL) |
| A2-interpolate | :349 message → `msg+": "+err.Error()` | **KILLED** 12/12 — `error=-32602 "not authorized: transition registry: bind …SECRET-R2R3…"` | KILLED (registry_moved) | — |
| A3-http-error | :349 → `http.Error(w, fmt.Sprintf("%d %s", code, msg), 500)` | **KILLED** 12/12 | KILLED (registry_moved) | — |
| A4-other-envelope | :349 → hand-encoded JSON-RPC envelope with `error.data = err.Error()` (the judge's "different error shape") | **KILLED** 12/12 — `error object carries unexpected key "data"` | KILLED (registry_moved) | — |
| **B1-drop-close** | :345 `if ctx.Err() != nil` → `if false` | **KILLED** 8 — R11 + all five held rows: `Connection="", want close` | **SURVIVED** | KILLED only by the loopback `TestA2AInvokeDeadlineClosesSocket` |
| **B2-always-close** | :345 → `if true` | **KILLED** 5 — R2_R3/R14/R16/R4_R6: `Connection="close", want absent` | **SURVIVED** | **SURVIVED** (rc=0, 0 FAIL) |
| **B3-close-by-error-type** | :345 → `if errors.Is(err, DeadlineExceeded) \|\| errors.Is(err, Canceled)` | **KILLED** 3 — uncertain_AppendIntent / uncertain_Commit: `Connection="", want close` (UncertainError does not unwrap, V8) | **SURVIVED** | **SURVIVED** (rc=0, 0 FAIL) |
| C-R2R3 | dispatchError absent/denied arm → `codeInternal, notAvailableMessage` | **KILLED** — `TestA2ADispatchWire/R2_R3` | KILLED | — |
| C-R11 | ctx arm code → `codeInvalidParams` | **KILLED** — R11 + 3 held R11 rows | KILLED | — |
| C-R14 | conflict arm code → `codeInvalidParams` | **KILLED** — R14 | KILLED | — |
| C-R16 | Unconfirmed arm code → `codeInvalidParams` | **KILLED** — R16 + 2 held R16 rows | KILLED | — |
| C-R4R6 | default arm code → `codeInvalidParams` | **KILLED** — R4_R6_default | KILLED | — |
| **D1-invokewait-ignored** | projection.go:266 `WithTimeout(…, h.invokeWait)` → `time.Hour` | **KILLED** 8 — `call 0 did not complete within 500ms (prototype bound; widened to +3 s, D1 still KILLED) … stalled GetReceipt`; R11: `"transition execution failed"` | **SURVIVED** | KILLED by `TestProjection_BoundedWait//a2a/_route…` + loopback test |
| D2-receipt-ctx-dropped | coordinator.go:236 `GetReceipt(ctx, …)` → `GetReceipt(context.WithoutCancel(ctx), …)` | **KILLED** 12/12 — but via the store's `ErrNoDeadline` (`binder reached 0 times`), not via the stall bound; D1 is the stall-bound kill | KILLED | — |
| **D3-uncertain-append-unwrapped** | coordinator.go:317 uncertain AppendIntent → `return Result{}, err` (drops `UnconfirmedError`) | **KILLED** — uncertain_AppendIntent: `"…not available…", want "…not confirmed; resend…"` | **SURVIVED** | KILLED by `TestDispatchDurableDeadline/uncertain_*` (Dispatch level) |

**15/15 KILLED by the new tests.** The row's point is A1, B2 and B3: wire-path mutants that the old table **and every other pre-existing test in the three packages** survive, and that only the new recorder-level tests kill. B1 was previously caught only by the loopback-socket test, which needs a bindable port. It is now also caught at the recorder, with no socket, which answers the judge's F2 sandbox gap for this mutant.

Rule-3n anchoring: the diff adds no production lines. The enumeration is therefore anchored to the production lines the new tests read: projection.go `:266` (D1), `:345-346` (B1–B3), `:348`→`dispatchError` arms (C-*), `:349` (A1–A4), and coordinator.go `:236`/`:317` (D2/D3, the held-op seam the repeated test exercises).

**History.** Run 1 **hung** on D1 at the 600 s harness timeout. The `R11` subtest's runner waited on `<-ctx.Done()` with no other exit, so a context-dropping mutant blocked forever. The runner was bounded (`select` with a 5 s fail-safe that returns an error) and the whole drill was re-run from scratch. The harness's `finally` restore had already returned `projection.go` to `48a9fb55…` (checked by `shasum` before the re-run). Run 1's A3 was a BUILD-FAIL, because `http.Error(w, msg, …)` left `code` unused, so it was not a kill. A3 was reworded to use `code`, and in run 2 it compiles and is KILLED.

## 5. Gates (prototype; logs in the evidence dir)

| Gate | Result | Log |
|---|---|---|
| `go vet ./...` | rc=0 | `gate_vet.log` |
| `gofmt -l host cmd` | empty (0 lines) | `gate_gofmt.log` |
| `go test -race -count=1 -v ./host/projection/ ./host/daemon/ ./host/coordinator/` | rc=0; `=== RUN` 81 / 298 / 65 = **444** (432 pre-existing + 12 new); **0** `--- FAIL` | `gate_race_targeted.log` |
| `go test -count=1 ./...` | rc=0; 24 `ok`, 0 `FAIL`; no load-flake re-runs needed. Re-run with this plan file present (the doc scanners see it): rc=0, 24 `ok`, 0 `FAIL`; verify_ail rc=0 | `gate_full.log`, `gate_full_with_plan.log`, `gate_verify_ail_with_plan.log` |
| `bash scripts/verify_ail.sh` | rc=0; 16 identities, 40 named tests, 9/9 package steps, PUB011 ratchet (`uncontracted_exports=3`) | `gate_verify_ail.log` |
| stability: cmd 1 with `-race -count=20` | rc=0; 240 `=== RUN`, 0 `--- FAIL` (timing bound and goroutine baseline are not flaky under race) | `gate_stability_x20.log` |

## 6. OPEN (not settled here)

- **OPEN-1 (scope, human judgement): real sole-connection hold at HTTP level.** The repeated test holds each durable operation at the `coordinator.Store` seam, not on the store's real pooled connection. The real hold is reachable only inside `host/store` (V9). Driving it from an HTTP test would need a new exported store test hook, which is a production surface, so none was added. The real-connection leg remains `TestDurableOpsHonourHeldConnection` (store level); the HTTP leg is this row. Whether the composition is enough is a call for the reviewer.
- **OPEN-2 (fault model disclosure): R2/R3 at dispatch is defense-in-depth.** The only producer of `AccessDeniedError` is `transitionreg.Bind` (`bind.go:71`), after the same `broker.Allows` check that admission's `Request.Allowed()` already applied to the **same** captured Request (`bind.go:128`, V7). An admitted skill therefore cannot reach this arm through the production binder. The wire row injects the typed error through `BinderFor`; it proves the arm's wire rendering, not a reachable production path.
- **OPEN-3 (semantics pinned, not changed): R16 under an expired deadline also closes the transport.** The uncertain held rows assert `Connection: close`, because `ctx.Err() != nil` holds. That is the production behaviour at `:345`, whose comment says "A timed-out invocation must release its transport". B3 shows a type-keyed refactor would silently drop the close for R16. If the intended contract is "close only on R11", that is a public-behaviour change and needs a ruling. The test makes either choice deliberate.
- **Not done (measured reason):** no repeated loopback-socket variant. `TestA2AInvokeDeadlineClosesSocket` already does one timed-out call plus a normal `task-normal` control on one real listener (`invoke_e2e_test.go:326`, `:371`, V12). The iteration-209 request accepts "daemon/projection" (V11), and the recorder test needs no bindable port.

## 7. Verification Log

The raw transcript is `/Users/voightkampff/.ailang/state/mission-world-iter213-evidence/verification_log_raw.txt`. Production-file claims were read from the **HEAD blob** (`git show HEAD:…`), so a mutation in flight could not skew them.

| V | Claim | Command → observed |
|---|---|---|
| V1 | base and branch | `git rev-parse --short HEAD; git branch --show-current` → `131c517`, `sprint/w-a2a-invocation-wire-coverage` |
| V2 | the wire path | `git show HEAD:host/projection/projection.go \| grep -n 'context.WithTimeout(r.Context(), h.invokeWait)\|h.coord.Dispatch\|if ctx.Err() != nil {\|"Connection", "close"\|dispatchError(err)\|A2AError(w, req.ID, code, msg)\|^func dispatchError'` → `266`, `339`, `345`, `346`, `348`, `349`, `391` |
| V3 | the Coordinator field is a concrete type (no dispatcher interface) | `… \| grep -n 'Coordinator \*coordinator.Coordinator\|coord *\*coordinator.Coordinator'` → `147 Coordinator *coordinator.Coordinator`, `164 coord *coordinator.Coordinator` |
| V4 | coordinator injection seams | `git show HEAD:host/coordinator/coordinator.go \| grep -n '^type Store interface\|^type Runner interface\|^type BinderFor\|^type Config struct'` → `27`, `37`, `42`, `45` |
| V5 | TestA2ADispatch structure | `git show HEAD:host/projection/projection_test.go \| grep -n '^func TestA2ADispatch(\|t.Run("registry_moved…\|…invalid_parts_and_pin\|…success_and_R13_reconciled\|dispatchError(fmt.Errorf'` → `802`, `803`, `804`, `805`, `834`; the registry_moved assertion is `1012 code, msg, _ := a2aErr(…)` / `1013 if code != codeInternal \|\| msg != "transitions that declare effects…"` (id discarded) |
| V6 | no test reads the Connection header | `git grep -n 'Get("Connection")' HEAD -- 'host/*_test.go' \| wc -l` → `0`; control in the same file family: `git show HEAD:host/projection/projection_test.go \| grep -c 'Get("Content-Type")'` → `4`; the socket test uses `ConnState` (`invoke_e2e_test.go:335`) |
| V7 | coordinator/transitionreg error sites the rows travel | `… coordinator.go \| grep -n …` → `236 GetReceipt(ctx, id)`, `276 coordinator: load source: %w`, `304 coordinator: execute: %w`, `322 if store.IsConflict(err)`, `328 &UnconfirmedError{… Err: err} // R16`; `bind.go` → `71 &AccessDeniedError{ID: id, Label: decision.Label}`, `78 fmt.Errorf("transition registry: bind %q: %w", id, err)`, `128 func (q Request) Allowed()`; `git grep -n 'AccessDeniedError{' HEAD -- 'host/*.go' \| grep -v _test` → only `bind.go:71` |
| V8 | UncertainError does not unwrap (so B3 drops the R16 close) | `git show HEAD:host/store/durable.go \| grep -c 'func (e \*UncertainError) Unwrap'` → `0`; control `git show HEAD:host/coordinator/errors.go \| grep -n 'func (e \*UnconfirmedError) Unwrap'` → `78`; `store.go:1148 return &ConflictError{ObservedHead…, SelectedHead: selected}` |
| V9 | the sole-connection hold is store-internal | `git show HEAD:host/store/durable_test.go \| grep -n 'conn, err := s.db.Conn(…)\|^func TestDurableOpsHonourHeldConnection'` → `74`, `94`; exported `*Store` methods whose signature takes a `func(`: `0`, of `34` exported methods (`grep -h '^func (s \*Store) [A-Z]' … \| wc -l`) |
| V10 | coordinator-level repeated stall exists, but not over HTTP | `… coordinator_test.go \| grep -n '^type stalledDurableStore\|^func TestDispatchDurableDeadline\|for i := 0; i < 3; i++ {'` → `531`, `575`, `618` |
| V11 | origin of the request | `grep -n 'no accumulating blocked workers across repeated calls' design_docs/planned/w-transition-invocation-coordinator-sprint-plan.md` → `32: 7. Add projection tests …`; iteration-209 planner (copied to the evidence dir as `world_iter209_planner.md`) `:11 … add a daemon/projection HTTP assertion with repeated calls …`; executor `:33 … The additional repeated HTTP held-connection deadline check was not added.`; judge `eval-report-iter208.md:175` |
| V12 | daemon socket tests | `git show HEAD:host/daemon/invoke_e2e_test.go \| grep -n '^func TestA2AInvokesPublishedTransition\|^func TestA2AResendAfterWriteTimeout\|^func TestA2AInvokeDeadlineClosesSocket\|task-normal'` → `94`, `273`, `326`, `371` |
| V13 | charter row and finding | `grep -n 'w-a2a-invocation-wire-coverage' design_docs/world-mission.md` → `1910: 123. **w-a2a-invocation-wire-coverage** · clause-6 residual …`; finding: `eval-report-iter208.md:16-17` (§0 items 1–2), `:149` (F1) |
| V14 | pinned binary | `$HOME/.pinned-ailang/ailang --version` → `AILANG v0.41.0` |
| V15 | test-only diff | `git status --short` → `?? host/projection/a2a_wire_test.go` (+ this plan); `git show HEAD:host/projection/projection.go \| shasum` = `shasum host/projection/projection.go` = `48a9fb552e16791c6aa20518db99b3f8439285bf` |
| V16 | a2aErr is strict about the envelope | `… projection_test.go \| grep -n '^func a2aErr\|if k != "code" && k != "message"'` → `280`, `301` |
| V17 | the testConfig default InvokeWait (the wire rows override it where they need a deadline) | `… \| grep -n 'InvokeWait: 2 \* time.Second'` → `220` |
