# Sprint plan — w-transition-invocation-coordinator (queue row 106, iteration 196)

**Base:** detached `dev 32f44c6` plus the designer's uncommitted prototype. **Authority:** [design](w-transition-invocation-coordinator.md), especially “Round 2 carve-out — the reviewers' fixes, applied verbatim”; that section overrides earlier milestone and bounded-wait claims. **Scope:** M1, M2, M3, M4a, M4b, M4c only. The controller commits; this planner performs no git write operation.

M1–M4c land as an offline host capability, with **no production caller**. Nothing in `host/daemon` or `host/projection` constructs a `Coordinator`; authorized `/a2a/` `tasks/send` continues to return its constant `-32603`. **M5 and M6 are blocked, not planned.** Their prerequisite is queue row 23's `D-WORLD-37` policy tranche: context-aware receipt lookup, connection acquisition, and pre-durable transaction operations, plus an enforced bound for remaining non-cancellable durable work. M5 must test a sole connection held past the deadline and stalls in each durable operation, showing bounded completion without accumulating blocked workers. An uncertain outcome must remain unconfirmed and reconcile without re-execution; a confirmed commit must not become a failure after a later context check.

## 1. Gates and measurements

Run every command with `export PATH=/opt/homebrew/bin:$PATH; export AILANG_BIN=$HOME/.pinned-ailang/ailang` (`AILANG v0.41.0`). The boundary simulation copied the worktree under `.plan-scratch/boundary`, restored only that copy's `capsule.go` from `HEAD`, and cumulatively applied the prototype hunks. M3 carried the six pure-law tests cut into a standalone `plan_test.go`; M4a carried a standalone construction test; M4b replaced them with the full prototype `coordinator_test.go`. M4c added the proposed `sync.Map` guard and `InFlightError` **in scratch only**. The first M3 test cut was malformed; after correction the entire sequence was rerun. These are the corrected results:

| Boundary | Added production code, approximately | `go vet ./...` rc | `go test ./... -run '^$'` rc |
|---|---:|---:|---:|
| M1 | capsule +27/−4 diff lines | 0 | 0 |
| M2 | binder 7 non-comment LOC | 0 | 0 |
| M3 | plan 94 + errors 40 = 134 | 0 | 0 |
| M4a | coordinator declarations/helpers through `parseOutput`, ~110 | 0 | 0 |
| M4b | `committed` + full `Dispatch`, ~121 additional | 0 | 0 |
| M4c | `sync.Map` claim/release and error, ~10 additional | 0 | 0 |

This re-partitions the design's M4c reconciliation into **M4b**, because the complete prototype `Dispatch` and its refusal suite already depend on `GetReceipt`, `committed`, R15 and R16. M4c is exclusively the quorum's missing concurrency fix. No boundary exceeds roughly 150 production LOC. The M3 pure-law test extraction and M4a construction test are milestone-local file states; do not land duplicate test definitions when M4b brings in the complete prototype test file.

**Re-measured here:** six compile boundaries (both commands, rc above); TR.C `go test ./host/broker -run '^TestRegistryDispatchBindingBoundary$' -count=1` in the M4c scratch tree, rc 0 with the new `binder.go`; prototype file and hunk inventory; M3 and M4a test-file compile dependencies; the design's step-1a omission from the prototype. **UNMEASURED here:** the designer's 32/32 mutation kills, real-interpreter echo/replay/incompatibility behavior, full package/full-repo test verdicts, the new R4 mutation kill, and all M5/M6 gates. The executor must perform those gates. The scratch guard passed `TestInFlightGuard` (all three subtests); its check and two release mutants each produced a targeted red (rc 1). The executor must transfer that scratch test into the landed milestone and rerun it.

The designer's V34 pre-existing macOS flake is `host/capsule TestOrdinaryOverflowCarriesNoESRCH`: EPERM from process-group kill occurred 3/320 with prototype and 3/320 with pristine `HEAD`. A recurrence is **not attributable to this sprint**. Any loopback-socket or signal/process failure in this sandbox is **UNINFORMATIVE UNDER SANDBOX**; the controller reruns those gates outside it.

## 2. Milestones and acceptance

Each row lists exact files and prototype hunks. Commands below inherit the pinned environment above. At **every** milestone repeat `go vet ./...` and `go test ./... -run '^$'`; `go build` alone misses `_test.go` compile failures.

### M1 — capsule argument and cancellation boundary

**Files:** `host/capsule/capsule.go`, `host/capsule/runcontext_test.go`. Take the entire prototype `capsule.go` diff: `Entry.Args`, `argsFile`, `Run` delegation to `RunContext`, `--args-file` staging, caller-derived timeout and `AILANG_RELAX_MODULES=1`; add the three prototype tests. The fixed staging path must accept the module headers that row 107's publish check accepts. Keep the legacy `Run` tests green. **Accept:** `go test ./host/capsule -run 'TestRunContext(Args|Cancel|RelaxedModule)$' -count=1`; `go test ./host/capsule -count=1` (V34 caveat). **Mutations:** CAP-ARGS, CAP-CTX, CAP-RELAX (§3).

### M2 — broker's narrow binder

**Files:** `host/broker/binder.go`, `host/broker/binder_test.go`, `host/broker/invoke_boundary_test.go`. Take both prototype binder files verbatim. Add the design's missing `NEG-session-binder` detector control to the `detector_controls` table in `invoke_boundary_test.go`: a synthetic outside-broker `*broker.SessionBinder` expression must produce zero findings. The existing `POS-session-type` control already detects raw `*broker.Session`. Do not expose the raw session. TR.C's inside exemption remains exactly 3. **Accept:** `go test ./host/broker -run 'TestOpenBinder|TestRegistryDispatchBindingBoundary' -count=1`. **Mutations:** BINDER-RAW, BINDER-DECLARED, TRC-SCOPE (§3). The existing TR.C test was re-measured green with `binder.go` in scratch.

### M3 — pure plan and refusal vocabulary

**Files:** `host/coordinator/plan.go`, `host/coordinator/errors.go`, `host/coordinator/plan_test.go` (milestone-local). Take the complete prototype `plan.go` and `errors.go`, including R15/R16 error types required by M4b. Cut `lawPlan`, `TestPlanRecordSkillID` and all five `TestPlanLaw*` functions out of prototype `coordinator_test.go` into `plan_test.go` with only their required imports. Pin the canonical record, input/output objects, observed head, revision, state root, log head, previous entry and journal binding. The pure Go transcription follows the already verified `world/transitions.ail` law; no `.ail` edit. **Accept:** `go test ./host/coordinator -run 'TestPlan(Law|RecordSkillID)' -count=1`. **Mutations:** PLAN-REV/STATE/LOG/PREV/OBSERVED/SKILL, R13-JOURNAL (§3).

### M4a — coordinator types and helpers

**Files:** `host/coordinator/coordinator.go`, `host/coordinator/new_test.go` (milestone-local). Take prototype `coordinator.go` from package/imports through `parseOutput` (immediately before `// committed`): narrow `Store`/`Runner`/`BinderFor`, `Config`/`New`, `Call`/`Result`, `InvocationID`, task/input validation, execution classification, JSON-object output validation. Add a construction test that rejects missing seams and non-positive caps; test helper behavior without a dispatcher. `plan_test.go` stays. **Accept:** `go test ./host/coordinator -run 'TestNew|TestPlan' -count=1`. **Mutations:** CFG-NIL, CFG-CAP, R1-TASK, R1-INPUT, R9-CLASSIFY, R12-OBJECT, ID-NS (§3); dispatcher-level killers become runnable in M4b.

### M4b — compose and reconcile, without a production caller

**Files:** `host/coordinator/coordinator.go`, `host/coordinator/coordinator_test.go`; remove the milestone-local `plan_test.go` and `new_test.go` as the full prototype test file replaces them. Take the rest of prototype `coordinator.go`, from `// committed` through the end: step 1b `GetReceipt`/committed read-back, bind/check/effect/source/world/execute/output/plan, pre-durable `ctx.Err`, intent then CAS commit, R14 bare conflict, R16 unconfirmed, and no post-success context check. Take complete prototype `coordinator_test.go` (pure-law, fake-runner refusal, real pinned interpreter, replay, boundary tests). Add `TestDispatchRefuses/R1_empty_episode`, `/R1_input_over_cap`, `/R1_unmarshalable_input`, `/R12_output_over_cap` and the missing R4 branch. Add `TestDispatchRefuses/R4_bind_error` using a fake `BinderFor` whose `Bind` returns an error; assert the typed/wrapped error and no execution or durable mutation. R2 and R3 must propagate `transitionreg.Bind` refusals without replacing the captured request's capability snapshot. **Accept:** `go test ./host/coordinator -run 'TestPlan|TestNew|TestDispatch|TestCommitBoundary|TestReplayCommittedInvocation' -count=1`; `go test ./host/broker -run '^TestRegistryDispatchBindingBoundary$' -count=1`. **Mutations:** R1–R16 and composition guards in §3. This milestone includes R13/R15/R16 reconciliation because the prototype `Dispatch` is one dependent unit.

### M4c — step 1a in-flight guard (quorum fix)

**Files:** `host/coordinator/coordinator.go`, `host/coordinator/errors.go`, **new** `host/coordinator/inflight_test.go`. Add a per-`Coordinator` `sync.Map` keyed by `InvocationID`. After validation and ID formation, **before `GetReceipt` (1b)**, claim the ID atomically with `LoadOrStore`; if present, return typed R17 `*InFlightError`. On a successful claim, `defer Delete(id)` immediately, before any later call or return. This includes normal success, every error path and a panic unwinding through `Dispatch`. Keep the guard per Coordinator and the current single-process writer premise explicit. `InFlightError` is a coordinator refusal type, not a wire mapping in this slice.

**Choice of the reviewer's two branches:** use **immediate refusal**. It gives a concurrent retry a prompt typed R17 without parking a second worker behind an execution or an unbounded durable store operation. The client can resend the same ID after the first call finishes; step 1b then returns the resolved receipt without re-execution or reports its unresolved intent. This choice preserves the carve-out's M5 block: it does not pretend that `GetReceipt`, `AppendIntent`, or `Commit` are bounded by the request context.

`TestInFlightGuard/R17_concurrent_same_id` uses a Runner double blocked on a channel after counting entry. While first `Dispatch` is inside the runner, a same-ID retry gets R17 and the count remains one. After release and first completion, another same-ID retry reaches 1b, returns `Reconciled: true`, and still has count one. `TestInFlightGuard/release_after_error` has the first runner call fail, then verifies a second same-ID call reaches the runner (not R17). `TestInFlightGuard/release_after_panic` recovers the first runner panic and verifies the ID was released. Also cover early returns after claim (for example, receipt error) with the same follow-up shape. Use channel synchronization, no timing sleep. **Accept:** `go test ./host/coordinator -run '^TestInFlightGuard$' -count=1`; `go test ./host/coordinator -count=1`; `go test ./host/broker -run '^TestRegistryDispatchBindingBoundary$' -count=1`. **Mutations:** R17-CHECK and R17-RELEASE, each with the sole killing subtest named below.

## 3. Mutation drill, anchored to the diff

For every row, mutate **only** the named file/change in a scratch copy or isolated executor edit, run **only** its named killer, record a behavioral `--- FAIL`, then restore byte-identical contents. A compile error, unused import or a passing test does not count as a kill. The designer reported 32/32 prototype mutations killed; this planner did not rerun them. R17-CHECK, R17-RELEASE and R17-PANIC-RELEASE were killed in scratch by the named selected subtests (rc 1); the R4 binder-error and new control mutations remain unmeasured obligations. The table names the **sole intended killing test** at subtest granularity; a row with two distinct refusal branches has two mutants.

| Mutation / refusal or guard | File: exact change | Sole killing test |
|---|---|---|
| CAP-ARGS | `host/capsule/capsule.go`: omit `--args-file` append while retaining the file write | `TestRunContextArgs` |
| CAP-CTX | `capsule.go`: derive timeout from `context.Background()` instead of `parent` | `TestRunContextCancel` |
| CAP-RELAX | `capsule.go`: remove `AILANG_RELAX_MODULES=1` from child env | `TestRunContextRelaxedModule` |
| BINDER-RAW / TR.C | `host/coordinator/coordinator.go`: introduce a `*broker.Session` type use outside broker | `TestRegistryDispatchBindingBoundary/outside_broker_is_clean` |
| BINDER-DECLARED | `host/broker/binder.go`: `b.s.Bind(m)` → `b.s.Bind(Manifest{})`, so the negative-cost manifest is accepted | `TestOpenBinder` |
| TRC-SCOPE | `host/broker/invoke_boundary_test.go`: classify `SessionBinder` as a forbidden `session-type` selector | `TestRegistryDispatchBindingBoundary/detector_controls` (`NEG-session-binder`) |
| PLAN-REV | `host/coordinator/plan.go`: `Revision: w.Revision + 1` → `+ 2` | `TestPlanLawRevision` |
| PLAN-STATE | `plan.go`: `StateRoot: out.Hash` → `in.Hash` | `TestPlanLawStateRoot` |
| PLAN-LOG | `plan.go`: `LogHead: entryHash` → `rec.Hash` | `TestPlanLawLogHead` |
| PLAN-PREV | `plan.go`: `PrevEntryHash: w.LogHead` → `w.Ref` | `TestPlanLawPrevEntry` |
| PLAN-OBSERVED | `plan.go`: `ObservedHead: w.Ref` → zero ref in `Commit` | `TestPlanLawObservedHead` |
| PLAN-SKILL | `plan.go`: `SkillID: d.ID` → empty string | `TestPlanRecordSkillID` |
| CFG-NIL | `coordinator.go`: remove missing-seam refusal in `New` | `TestNewRefusesMissingSeamsAndCaps` |
| CFG-CAP | `coordinator.go`: accept `MaxInput <= 0` or `MaxOutput <= 0` | `TestNewRefusesMissingSeamsAndCaps` |
| R1-TASK | `coordinator.go`: guard `!validTaskID` with `false &&` | `TestDispatchRefuses/R1_bad_task_id` |
| R1-INPUT | `coordinator.go`: replace nil-input guard with `if false` | `TestDispatchRefuses/R1_nil_input` |
| R1-EPISODE | `coordinator.go`: `if call.EpisodeID == ""` → `if false` | `TestDispatchRefuses/R1_empty_episode` (new) |
| R1-CAP | `coordinator.go`: replace `len(in) > c.cfg.MaxInput` with `false && len(in) > c.cfg.MaxInput` | `TestDispatchRefuses/R1_input_over_cap` (new) |
| R1-MARSHAL | `coordinator.go`: replace `err != nil` in the marshal refusal with `false && err != nil` | `TestDispatchRefuses/R1_unmarshalable_input` (new) |
| R2-ABSENT | `coordinator.go`: swallow `transitionreg.Bind`'s `TransitionAbsentError` and return generic error | `TestDispatchRefuses/R2_absent` |
| R3-CAPS | `coordinator.go`: pass ambient Admin caps in place of `call.Request.Caps` to `transitionreg.Bind` | `TestDispatchRefuses/R3_access_denied` |
| R4-BIND | `coordinator.go`: swallow a fake binder's `Bind` error | `TestDispatchRefuses/R4_bind_error` (new) |
| R5-PIN | `coordinator.go`: ignore `call.PinnedFn` (`_ = call.PinnedFn` in its arm) | `TestDispatchRefuses/R5_pin_mismatch` |
| R6-ABSENT | `coordinator.go`: source `if !ok` → `if false` | `TestDispatchRefuses/R6_absent_source` |
| R6-HASH | `coordinator.go`: `sum != d.TransitionFn` → `false && sum != d.TransitionFn` | `TestDispatchRefuses/R6_corrupt_source` |
| R7-WORLD | `coordinator.go`: selected-head `if !ok` → `if false` | `TestDispatchRefuses/R7_no_world` |
| R8-EFFECT | `coordinator.go`: declared-effects guard → `if false` | `TestDispatchRefuses/R8_effects_declared` |
| R9-INCOMPAT | `coordinator.go`: classify interpreter mismatch as `ExecutionError` | `TestDispatchIncompatibleRealInterpreter` |
| R10-EXEC | `coordinator.go`: replace `classifyExec(err)` return with an untyped error | `TestDispatchRefuses/R10_exec_failed` |
| R11-EXEC | `coordinator.go`: disable `ctx.Err()` guard on runner error (`false && cerr != nil`) | `TestDispatchRefuses/R11_cancelled_during_exec` |
| R11-BOUNDARY | `coordinator.go`: disable `ctx.Err()` guard immediately before `AppendIntent` | `TestCommitBoundary/cancel_before_commit` |
| R12-OUTPUT | `coordinator.go`: decode into `any` instead of requiring non-nil `map[string]any` | `TestDispatchRefuses/R12_output_not_object` |
| R12-JSON | `coordinator.go`: after `json.Unmarshal(out, &obj)`, replace `err != nil || obj == nil` with `false && err != nil` | `TestDispatchRefuses/R12_output_not_json` |
| R12-CAP | `coordinator.go`: `if len(out) > c.cfg.MaxOutput` → `if false` | `TestDispatchRefuses/R12_output_over_cap` (new) |
| R13-RECONCILE | `coordinator.go`: `if seen` → `if false && seen` | `TestDispatchRefuses/R13_resend_reconciles_committed` |
| R13-JOURNAL | `plan.go`: set `Commit.InvocationID` to empty while leaving intent ID | `TestCommitBoundary/receipt` |
| R14-CAS | `coordinator.go`: re-read the head and silently re-plan before commit | `TestDispatchRefuses/R14_head_moved` |
| R14-TYPED | `coordinator.go`: `if store.IsConflict(err)` → `if false` | `TestDispatchRefuses/R14_head_moved` |
| R15-INTENT | `coordinator.go`: `if rc.State == store.ReceiptResolved` → `if true` | `TestDispatchRefuses/R15_resend_not_committed` |
| R16-UNCONFIRMED | `coordinator.go`: return bare `err` instead of `UnconfirmedError` on non-conflict commit error | `TestDispatchRefuses/R16_unconfirmed_then_reconciled` |
| POSTCOMMIT-SUCCESS | `coordinator.go`: add `if err := ctx.Err(); err != nil { return Result{}, err }` after successful `Commit` | `TestCommitBoundary/deadline_during_commit_reports_success` |
| ID-NS | `coordinator.go`: `"a2a:" + episodeID` → `"effect:" + episodeID` | `TestDispatchEchoRealInterpreter` |
| R17-CHECK | `coordinator.go`: `if _, loaded := c.inFlight.LoadOrStore(id, …); loaded` → `if false && loaded` (keep atomic store and use `loaded`) | `TestInFlightGuard/R17_concurrent_same_id` |
| R17-RELEASE | `coordinator.go`: change `defer c.inFlight.Delete(id)` to a no-op `defer func() {}()` | `TestInFlightGuard/release_after_error` |
| R17-PANIC-RELEASE | `coordinator.go`: move `Delete(id)` to normal-return code rather than a defer | `TestInFlightGuard/release_after_panic` |

R13 is a reconciliation success, rather than a refusal; its mutation proves no second execution. R14 is a bare store conflict and R16 is an unconfirmed outcome. `*store.DuplicateInvocationError` remains a defensive error, **not** unreachable: R17 protects concurrent callers sharing one Coordinator, while the store retains the last guard. Avoid an R15 inference across multiple processes or a widened connection pool; row 23 must revisit that premise before M5.

## 4. Execution notes and risks

1. Apply milestone files in the measured order. Keep `plan_test.go`/`new_test.go` only for their boundaries, then merge their tests into the prototype `coordinator_test.go` at M4b. Copy milestone states only under worktree-local `.plan-scratch/` while rehearsing, and remove those scratch copies when done; do not modify prototype files to simulate boundaries.
2. At each boundary run the two compile fences and the row's targeted command with the pinned interpreter. At M4c run `go vet ./...`, `go test ./... -run '^$'`, `go test ./host/coordinator -count=1`, `go test ./host/broker -run '^TestRegistryDispatchBindingBoundary$' -count=1`, then the non-socket full gate `go test ./... -count=1` and `./scripts/verify_ail.sh`. The controller reruns sandbox-sensitive gates outside the sandbox. Do not classify socket or signal-process failures here as product verdicts.
3. Execute the mutation drill after the tests exist. Record red test identities and restore each edit byte-identically. A full green run without mutations does not discharge S6's load-bearing claims.
4. No `world/`, `.ail`, `tools/launchd/*`, `.claude/*`, `.agents/*`, daemon or projection change belongs in this slice. The invocation source fixture stays inline in Go; the production skill and usage walkthrough wait for M5/M6, after row 23. The major residual risk is context-free `GetReceipt`, `AppendIntent` and `Commit`: keep the coordinator uncalled in production until the prerequisite lands.
