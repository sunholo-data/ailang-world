# Independent PRODUCT evaluation — iteration 217, round 2

**FAIL — 54/100 (exact 53.5454545), hard test failure.** Tested immutable candidate `59870e58c4fc393aab62d2cf7ff8c83641ad6850`. This is an independent product judgment, not controller acceptance. The prior round 1 FAIL 49/100 is preserved.

Judge: native GPT6Astra per Agent assignment; generator/controller/executor native GPT6.1Sol. D-WORLD-48 same-vendor different-model fallback **FLAG**; preferred Sonnet was unavailable according to controller routing evidence. Native token usage is unknown, not zero. No subagents, production fixes, remote messages, commits, or harness changes were made.

Own detached worktree: `/Users/voightkampff/dev/sunholo-data/.wt-world-iter217-product-eval-r2`. Worktree creation completed successfully and exact HEAD was verified before review. Runtime started only after controller's full Go process terminated. All evaluator runtime processes are terminal; temporary mutations were restored byte-for-byte in `finally`. Final tracked diff is empty; only this report directory is untracked. Evaluation follows the absolute sprint-evaluator skill/rubric and mission verification protocol, with World scripts in place of compiler-repository make targets.

## Blocking evidence

Controller ran the exact unchanged profiles sequentially in `.check-world-iter217` at the candidate SHA, released interpreter pinned first in PATH, GODEBUG unset and default Go cache. Evaluator independently read both complete raw logs, metadata and lossless base64 copies and verified byte identity. These are controller measurements, not falsely labelled evaluator reruns:

- `bash scripts/verify_ail.sh`: rc 0, 67.857 seconds, external cap 600 seconds.
- `bash scripts/verify_go.sh`: rc 1, 459.398 seconds, external cap 1200 seconds. Plain tests passed all 24 packages. Race tests failed six functions across four packages.

| Failing race test | Location | Actual failure |
|---|---|---|
| TestOrdinaryOverflowCarriesNoESRCH | host/broker/cleanup_test.go:200 | Joined overflow/EPERM error instead of exact typed overflow |
| TestHandlerTimeoutKillsTheWholeProcessGroup | host/broker/handlers_test.go:801 | exec_started=false and forked=false at timeout |
| TestF5WallClockTimeoutHasElapsedBound | host/capsule/capsule_test.go:230 | 4.511691083s exceeds 2s |
| TestDispatchDurableDeadline/uncertain_AppendIntent | host/coordinator/coordinator_test.go:629 | Receipt load deadline before expected uncertain R16 result |
| TestA2ADurableDeadlineWire/uncertain_AppendIntent | host/projection/a2a_wire_test.go:467 | Invocation deadline instead of uncertain reconciliation error |
| TestMCPInvokeDeadlineFreesSlot | host/projection/mcp_bounds_test.go:65 | Recovery id 2 gives HTTP 200, JSON -32603 `host callback timed out`, instead of SSE |

**BLOCK-R2-FULL:** Get the unchanged full AIL→Go profiles green on a repaired immutable candidate. Diagnose each failure with controlled reproductions; repair fixture or product behavior according to evidence. Byte-identical broker/capsule/coordinator source and unchanged A2A uncertain arm establish change scope, not runtime baseline reproduction or an inherited-failure waiver.

**BLOCK-R2-MCP:** The MCP recovery failure is in the new test and cannot be called inherited. Own isolated `-race -count=10` passes and `GOMAXPROCS=1 -race -count=20` passes do not erase the full-profile failure. An independent test-only calibration adds 60ms resolver delay to the existing 30ms resolver delay under the 80ms aggregate budget. It compiles, reproduces the same recovery timeout while explicitly observing the first callback returned and one capsule run, then restores green. This proves an alternative causal path to the symptom; it does not prove which phase exhausted the original full-profile run. Recommended repair: isolate the artificial deadline stimulus from recovery, keep the same handler, assert first callback return and phase-entry counters, and give recovery its own bounded live budget. Re-kill callback-context and aggregate controls, then rerun the full profile. Do not blindly enlarge the first stimulus or change production 3/10/20/30 constants.

## Literal rubric

| Category | Score | Evidence/deduction |
|---|---:|---|
| Tests | 0/20 | Full Go rc1, mandatory hard failure |
| Lint | 10/10 | Own `go vet ./...` rc0; gofmt list of changed Go files empty |
| Acceptance | 24.545/30 | A1 5 + A2 6 + B 3 + C 4 = 18 of 22 criteria; D passes=null contributes zero |
| Code quality | 0/15 | Touched daemon.go 966, daemon_test.go 1343 and projection_test.go 1553 lines; three -5 deductions floor score at zero |
| Documentation | 10/15 | Active HOST_CHANGELOG Unreleased entry and executed host examples/three guide payloads; design status loses 5 |
| Design fidelity | 9/10 | Thin upstream adapter and World admission architecture match; usage/surface refusal coverage is narrower than complete product claims |

The over-800-line files are inherited (base 963/1341/1541); the rubric still deducts literally. No unrelated splits are requested. Documentation example credit uses the applicable host-only equivalent: executable Go examples and real-daemon guide payloads, not an invented `.ail` feature. The canonical design header still says Planned/planning authorized/acceptance NOT RUN while appended text says implemented locally; that contradiction prevents status credit. D's final-CI/evaluator criterion is pending and must not be marked true to inflate this score. Compiler regression and performance bonus categories do not apply: no shared compiler paths change, and the clocks are correctness bounds rather than an optimization goal.

## Independent repair validation

`TestA2AResendAfterWriteTimeout` passed normal count 3 (10.840s command), race count 3 (13.297s), and the unchanged real-capsule sibling `TestA2AInvokesPublishedTransition` passed (1.682s). Source review confirms the repair retains actual registry, published archive, store and coordinator; the cooperative deterministic echo runner removes irrelevant interpreter startup from this socket/receipt fixture. InvokeWait is 3s, socket WriteTimeout remains 2s, actual post-commit hold is 2.1s, first response fails, retry uses the same task, entry index remains 1/no second entry and exactly one execution is asserted. Idempotent release and joining the first request precede server close.

Two own test-only calibrations compiled and failed the intended assertion: changing retry task produced `same task executed 2 times`; forcing fatal after the commit signal terminated cleanly without blocked server close or test timeout. Both restored green. The real-capsule execution/replay sibling's function bytes are unchanged. This repairs the round 1 resend blocker's tested scope; it does not pass the entire candidate.

## Product contract audit

All 14 ACs were reviewed in source and against raw evidence, plus the evaluator's own relevant MCP/A2A/daemon baseline (rc0, 86.491s). Scope is darwin/arm64, Go 1.26.6, pinned AILANG binary v0.41.0 commit 24ee1088776e21cd06a3781ed18e77f40be06db3; module dependency v0.47.2 is a separate identity.

| AC | Reviewed result and limits |
|---|---|
| DEP | Own dependency closure has 56 imports versus existing baseline evidence 54, adding exactly protocol/hostcall and protocol/mcphttp, removing none. Four frozen protocol files byte-identical across module versions. |
| MAP | Canonical underscore/dot/slash escaping, decoding and 64-byte refusal checked. Whole-surface refusal/recovery tested; exact refusal envelope assertion remains narrower than ideal. |
| SCHEMA | Object schema normalization, refusal of non-object type, exact 9007199254740993 preservation via RawMessage verified. |
| CARRIER | Only Authorization bearer header enters shared resolver; refusal does not dispatch. |
| ADMIT | Fresh per-item admission and K+1 reads, revocation/readmission and no-dispatch counters checked. |
| ABSENCE | Typed confirmed absence differs from store/integrity/context failures and raced heads; own A1/A2 reversals kill intended assertions. |
| INVOKE | World coordinator owns commit/output; absent coordinator refuses; supporting refusal hunk currently lacks a killing test (see survivors). |
| ITEM | Repeated RPC IDs mint distinct random 32-byte/64-hex task IDs; RPC ID is not durable idempotency key. |
| BATCH | Supported default 2025-03-26 array handling, newer-version HTTP400; per-item commits are non-atomic and earlier success survives later refusal. |
| BOUNDS | Aggregate request deadline differs from per-callback budget; Runner capacity8 lasts through callback return, separate subprocess capacity8 lasts through reap. New tiny-budget recovery fixture is blocking under full race load. |
| WIRE | Frozen upstream owns parsing/envelopes/SSE; no World duplicate codec. Golden is regenerated by running upstream and compared. |
| ROUTE | Exact POST /mcp/ registration checked by source AST; own generic-mount mutant is killed while GET405 alone survives. |
| CROSS | Shared admission/card/A2A laws checked, but full A2A uncertain-deadline race test remains red. |
| FIXTURE | Real-daemon pinned interpreter usage and three guide payloads execute; no claim that every shell guide command/initialize recipe was executed. |

Two own real production-clock legs measured restored single write 20.003757208s and batch write 20.001360583s, before the 30s socket cutoff. Resolver stage holds ~2.8s and tools stage ~9.8s. Single N/K/J=1/1/1, admissions2, commits0; batch 3/2/2, admissions3, commits1, no third dispatch; follow-up succeeds. Removing only aggregate context makes both fail at writes 32.6096s/34.6140s with EOF. These prove those staged cases, not universal delivery by20s; the separate late-body 30s cutoff/truncation case remains an explicit residual limit.

## Independent nonvacuity and evidence scope

All seven production/module paths and each stage's hunks are preserved in JSON. They and all `.ail` files are byte-identical to round1. A1, A2, B and C have independent evaluator compile-fenced behavioral controls below. D adds conformance tests and has **no own production diff**; its controls belong to B/C, not a fictitious D production revert. Full API-removing reverts can compile-refuse; evaluator runtime reversals retain API scaffolding and do not claim full-hunk coverage.

- **OWN-A1-behavior-revert** (`host/transitionreg/transitionreg.go`): compile rc 0, selected rc 1, restored rc 0. Selected failures: TestRegistryHeadAbsentErrorClassification, TestRegistryHeadAbsentErrorClassification/direct, TestRegistryHeadAbsentErrorClassification/new_request, TestRegistryHeadAbsentAfterLookupCancellation, TestRegistryHeadAbsentAfterLookupCancellation/absent, TestRegistryHeadAbsentAfterLookupCancellation/present, TestRegistryHeadAbsentAfterLookupCancellation/deadline. Broad actual red set: TestRegistryHeadAbsentErrorClassification, TestRegistryHeadAbsentErrorClassification/direct, TestRegistryHeadAbsentErrorClassification/new_request, TestRegistryHeadAbsentAfterLookupCancellation, TestRegistryHeadAbsentAfterLookupCancellation/absent, TestRegistryHeadAbsentAfterLookupCancellation/present, TestRegistryHeadAbsentAfterLookupCancellation/deadline.
- **OWN-A2-classifier-revert** (`host/projection/projection.go`): compile rc 0, selected rc 1, restored rc 0. Selected failures: TestAgentCard_AbsentPrecheckSnapshotFailure, TestAgentCard_AbsentPrecheckSnapshotFailure/store_error, TestAgentCard_AbsentPrecheckSnapshotFailure/real_reader, TestAgentCard_AbsentPrecheckSnapshotFailure/deadline, TestMCPAbsentPrecheckSnapshotFailure, TestMCPAbsentPrecheckSnapshotFailure/list/direct, TestMCPAbsentPrecheckSnapshotFailure/list/real_reader, TestMCPAbsentPrecheckSnapshotFailure/list/same_text, TestMCPAbsentPrecheckSnapshotFailure/list/deadline, TestMCPAbsentPrecheckSnapshotFailure/invoke/direct, TestMCPAbsentPrecheckSnapshotFailure/invoke/real_reader, TestMCPAbsentPrecheckSnapshotFailure/invoke/same_text, TestMCPAbsentPrecheckSnapshotFailure/invoke/deadline. Broad actual red set: TestAgentCard_AbsentPrecheckSnapshotFailure, TestAgentCard_AbsentPrecheckSnapshotFailure/store_error, TestAgentCard_AbsentPrecheckSnapshotFailure/real_reader, TestAgentCard_AbsentPrecheckSnapshotFailure/deadline, TestProjectionAbsenceRaceControls, TestProjectionAbsenceRaceControls/same_text_genuine, TestA2A_AbsentPrecheckSnapshotFailure, TestA2A_AbsentPrecheckSnapshotFailure/direct, TestA2A_AbsentPrecheckSnapshotFailure/real_reader, TestMCPAbsentPrecheckSnapshotFailure, TestMCPAbsentPrecheckSnapshotFailure/list/direct, TestMCPAbsentPrecheckSnapshotFailure/list/real_reader, TestMCPAbsentPrecheckSnapshotFailure/list/same_text, TestMCPAbsentPrecheckSnapshotFailure/list/deadline, TestMCPAbsentPrecheckSnapshotFailure/invoke/direct, TestMCPAbsentPrecheckSnapshotFailure/invoke/real_reader, TestMCPAbsentPrecheckSnapshotFailure/invoke/same_text, TestMCPAbsentPrecheckSnapshotFailure/invoke/deadline.
- **OWN-B-mapping-schema-behavior-revert** (`host/projection/mcpname.go`): compile rc 0, selected rc 1, restored rc 0. Selected failures: TestMCPNameRoundTrip, TestMCPNameRefusal, TestMCPSchemaNormalization. Broad actual red set: TestMCPAbsentPrecheckSnapshotFailure, TestMCPAbsentPrecheckSnapshotFailure/invoke/direct, TestMCPAbsentPrecheckSnapshotFailure/invoke/real_reader, TestMCPAbsentPrecheckSnapshotFailure/invoke/same_text, TestMCPAbsentPrecheckSnapshotFailure/invoke/deadline, TestMCPSurfaceRefusalRecovers, TestMCPListedThenTrulyAbsent, TestMCPInvokeDeadlineFreesSlot, TestMCPBatchAccounting, TestMCPBatchAccounting/repeated_rpc_id, TestMCPBatchAccounting/removed_between_items, TestMCPBatchConformance, TestMCPBatchConformance/old_versions_and_item_errors, TestMCPBatchConformance/whole_host_failure, TestMCPBatchConformance/notifications_and_bad_arrays, TestMCPBatchConformance/newer_versions_refuse, TestMCPTrueAbsenceAndPublication, TestMCPAbsentPrecheckSuccessfulInvocation, TestMCPListExactSetPerSession, TestMCPCallDispatch, TestMCPInvokeListedThenRevoked, TestMCPNameRoundTrip, TestMCPNameRefusal, TestMCPSchemaNormalization, ExampleEncodeMCPName.
- **OWN-C-exact-method-registration** (`host/daemon/daemon.go`): compile rc 0, selected rc 1, restored rc 0. Selected failures: TestMCPMountMethodSource.
- **CAL-resend-distinct-retry** (`host/daemon/invoke_e2e_test.go`): compile rc 0, selected rc 1, restored rc 0. Selected failures: TestA2AResendAfterWriteTimeout.
- **CAL-resend-postcommit-failure-cleanup** (`host/daemon/invoke_e2e_test.go`): compile rc 0, selected rc 1, restored rc 0. Selected failures: TestA2AResendAfterWriteTimeout.
- **CAL-slot-recovery-resolver-deadline** (`host/projection/mcp_bounds_test.go`): compile rc 0, selected rc 1, restored rc 0. Selected failures: TestMCPInvokeDeadlineFreesSlot.
- **OWN-C-aggregate-context-removed** (`host/projection/mcp.go`): compile rc 0, selected rc 1, restored rc 0. Selected failures: TestMCPPostBudgetProductionConstants, TestMCPPostBudgetProductionConstants/single, TestMCPPostBudgetProductionConstants/batch. Broad actual red set: TestMCPWireOwnershipSource.

The JSON retains complete raw command output, exact argv, elapsed times, finite external caps, actual failing-test sets, mutant diffs and original/mutant/restored hashes. Every restore hash matches; no external cap expired. Broad mutation runs are scoped exactly, not treated as universal hunk coverage. A1 integrity negative controls surviving a typed-absence reversal are correctly retained as regression controls. B combines mapping/schema behavior; it does not independently establish every possible schema hunk dependency. The old dependency pin's compile refusal is distinct from runtime nonvacuity.

Audited 619 executor-evidence artifact hashes and 23 repair artifact hashes (642 total), with no mismatch. All 43 named mutant raw selections/restores and failure-text assertions were audited. These are 42 compiled runtime/source/test-decorator controls plus one old-pin compile refusal, **not 43 independent production behaviors**: aggregate-removed/per-item controls share the same context change, snapshot controls overlap, and three upstream decorators are calibrations. Executor/round1 evidence is corroboration and is not relabelled as this evaluator's own execution.

Two round1 supporting-hunk survivors remain unchanged, nonblocking and queued132: `mcp.go:19` capacity cfg.MaxCallbacks→800, and `mcp.go:123` nil-coordinator refusal→success. They survived compiled whole-projection tests in round1; round2 corroborates unchanged source, not fresh reruns or invented kills. No claim of numeric line/branch coverage is made.

Other support limits: surface refusal test lacks exact status/ID/message assertions; the unit escape-heavy name is outside registry segment grammar while the live long-name case is valid; guide testing reconstructs setup/headers and executes three payloads rather than all guide commands; projection.go Handler comment still says two routes/no workers. No platform/load matrix, original full-failure unique root cause, final CI or landing is claimed.

## Artifacts and final markers

- `product-evaluator-independent.json`: complete machine-readable judgment, 38 command records, 8 own mutation/calibration records, root profile raw evidence, source hashes, stage hunks and named-mutant audit.
- Existing round1 FAIL report is unchanged and its SHA256 is retained in JSON.
- Original authorization hash remains immutable; separately recorded current design metadata hash is not substituted for it.

EVALUATION_RESULT: fail
EVALUATION_SCORE: 54/100
EVALUATION_ROUND: 2
TESTED_SHA: 59870e58c4fc393aab62d2cf7ff8c83641ad6850
EVALUATION_REPORT_PATH: /Users/voightkampff/dev/sunholo-data/.wt-world-iter217-product-eval-r2/design_docs/verification/world-iter217/evaluator-product-round2/product-evaluator-independent.json
ALL_OWNED_RUNTIME: TERMINAL
