# Independent product review — iteration 217, round 1

Candidate: `be286c1dbf29eb56c922f5015c0bc1a52055f087`. Product base: `54e0fb0a5fb0c988fb8d9f9e9bcda28096a1183e`.

**Disposition: FAIL — round1, 49/100; pristine test failure reproduced at both candidate and exact base.** This is not an attribution of the failure to MCP. No product fix has been made by this evaluator. Full machine-readable commands, outputs, durations, source hashes and actual red-test sets are in `product-evaluator-independent.json` beside this report.

## Independence and scope

Fresh native evaluator context and persistent isolated detached worktree, creation terminal rc0 and exact candidate HEAD confirmed before reads. Same-vendor/different-model routing is the explicitly flagged D-WORLD-48 fallback. The controller reported that the original spawn prompt omitted the required leading `MISSION-ROLE: evaluator`; a later corrective label does not erase that omission. Native token usage is unavailable, not zero. No descendants, external messages, pushes, production fixes or shared-file writes. All own temporary product mutations were restored byte-identically.

Read the absolute sprint-evaluator skill, full scoring rubric, mission-control skill and full verification protocol, World repository guidance, charter, relevant architecture, frozen full design, five-stage plan and JSON. The review started at 21:00 UTC with a 21:45 UTC bound. Root profiles completed before independent production-duration timing work. Controller verdicts were not used as acceptance evidence.

## Blocking finding

`host/daemon/invoke_e2e_test.go:289`: `TestA2AResendAfterWriteTimeout` failed on the unchanged candidate after every product mutation was restored. Exact command with the released pin and normal Go environment:

```sh
AILANG_BIN=/Users/voightkampff/.pinned-ailang/ailang go test ./host/daemon -run '^TestA2AResendAfterWriteTimeout$' -count=3 -v -timeout=90s
```

Externally bounded at 120 seconds; terminal rc1. Repetitions were PASS 9.75s, FAIL 13.73s (`commit did not land`, then `httptest.Server blocked in Close` with an active connection), PASS 8.43s. The same test failed the broad aggregate-mutated daemon run; it is not counted as an MCP aggregate killer. The independent pristine failure cannot be erased by the earlier green full profile.

The fixture uses real capsule execution with injected InvokeWait 3s, server WriteTimeout 2s, a 5s commit watchdog and a 2.1s post-commit hold. Those timing constants predate this sprint. A concrete repair candidate is to use a deterministic cooperative runner for this socket/receipt law while retaining the separate real interpreter/replay test; preserve the observed durable-commit event, actual socket timeout, failed first response, same-task entry 1 reconciliation and absence of entry 2. Also release the commit seam and terminate workers on every failure path. This is a proposed fixture repair, not a proven cause or permission to relax production constants. Exact-base comparison was completed before regression attribution: clean detached base54e0fb0, same binary/environment/command, terminal rc1 in44.884s; PASS4.12s, PASS4.10s, FAIL10.42s at the same line289 and `commit did not land` message. This supports an inherited fixture failure; it does not waive the candidate test gate. The base run emitted no Close warning, so that additional candidate symptom is not falsely called reproduced on base.

## Acceptance evidence

All 14 named MCP product behaviors were examined; the supporting evidence limits below remain explicit. This behavioral assessment does not override the pristine-suite blocker or claim final CI/row108 completion.

| Acceptance ID | Independent assessment and evidence |
|---|---|
| DEP | Own `go list` closure is 54→56, adding only released `protocol/hostcall` and `protocol/mcphttp`, removing none. Four existing protocol production files byte-identical across pins. Module v0.47.2 is distinct from execution binary v0.41.0. Old-pin build refusal is honestly a compile gate. |
| MAP | Own runtime identity/decode mutation kills name corpus and live calls; exact escaping, canonical decode and 64-byte refusal inspected. Whole-surface failure and same-handler recovery run. Exact-envelope and registry-valid escape-expansion fixture support are narrower than the strongest plan wording. |
| SCHEMA | RawMessage-based normalization preserves object properties, typed bytes and integer 9007199254740993. Explicit nonobjects/bad annotation refused. Own passthrough behavior reversion and named float/constraint mutation logs inspected. |
| CARRIER | Authorization-only resolver, constant denial and bounded credential context; denial matrix plus authorized nonzero-read control run. Middleware does not preempt MCP. Resolver errors are upstream envelopes. |
| ADMIT | Unequal session sets, empty real absence, dynamic head and genuine errors covered. Shared helper reads fresh registry rather than caches request authorization. Named unfiltered/no-deadline kills verified. |
| ABSENCE | Own A1 and A2 behavior reversions produce the exact cancellation/classification and genuine-error failures. Direct/real/same-text/deadline cases run for list and fresh invocation, with no dispatch/mint on refusal. A2A law retained by source and targeted controls; separate A2A socket fixture blocker remains. |
| INVOKE | Object input, validated session/name, fresh admission and final context check precede mint/dispatch. Real coordinator only, entropy failure and revocation/cancellation controls run. Nil-coordinator rejection branch is correct but independently found unpinned. |
| ITEM | Repeated RPC ID still creates distinct random 64-hex tasks; K+1 admission reads, sequential overlap 1, receipts and real journal entries measured; removal between items refuses later item. D adds tests pinning C behavior, not D production behavior. |
| BATCH | Frozen upstream owns default/explicit 2025-03-26 arrays, newer-version 400, notification 202, malformed/empty refusal, item-error continuation and whole-host-error abort. Earlier successful item remains committed; third item never invoked. Wire-only decorators are calibration controls, not imaginary product branches. |
| BOUNDS | Injected 3/10/20/30 and slot8 source wiring confirmed. Own actual production single/batch legs green around20s; outer-budget removal compiles and makes both named timing/counter controls red above32s/35s; restored both green. Runner callback-return and subprocess-reap bounds are distinct. Literal800 Runner-capacity survivor is a coverage follow-up. |
| WIRE | World owns callbacks plus context wrapper; released mcphttp owns HTTP/JSON-RPC/SSE. Generated golden source and runtime output audited, body4MiB/headers/version/host-error tests run, source ownership guard calibrated. Late-body socket truncation remains an explicit residual. |
| ROUTE | Exact executable `POST /mcp/` registration inspected. Own methodless-mount mutation fails AST assertion while GET405 remains green, demonstrating why network405 alone is vacuous. Existing A2A route controls retained. |
| CROSS | One daemon, both unequal populated sessions, exact cardinality and decoded MCP/Card equality including dot/slash IDs. Divergence mutant audited. |
| FIXTURE | Both session surfaces have nonempty positive controls; recorded ambient export union and existing A2A bans absent. Hardcoded ambient mutant audited. |

## Source completeness and per-stage nonvacuity

Enumerated every production diff hunk, with complete diffs and hashes in JSON: transition registry reader; projection Config/New/shared admission; new name/schema helpers; new MCP factory/resolver/tools/invoker/task mint; daemon scalars/mount; go.mod/go.sum. No kernel `.ail`, coordinator, fleet, verification script or workflow production change. The x/sys version bump is the released module's required transitive update; it does not add a package path to the measured closure.

* **A1:** Exact own production revert removes the newly exported type and fails compilation, explicitly not a runtime kill. Retaining the public API while reverting emission/post-lookup cancellation compiles, kills the direct/wrapped type assertions and cancellation cases, and restores green. Negative integrity cases and retained API example survive as regression controls. Existing `%w` wrapping is a prerequisite, not a newly introduced A1 hunk.
* **A2:** Exact own classifier revert compiles and fails genuine absent-precheck snapshot-error tests. Own broad projection run lists all actual A2A/MCP classification failures. True absence, present→absence, failed precheck and publication controls that already held before A2 survive and are not misfiled as A2-new-behavior killers.
* **B:** Exact own revert removes new APIs/dependency and compile-refuses. Retained-API identity mapping/schema passthrough reversion compiles, named runtime assertions fail, broad actual red set is recorded, and restoration passes. No compile failure is mislabeled as a semantic kill.
* **C:** Exact own revert compile-refuses because C introduced required APIs/configuration. Scoped own behavioral controls cover daemon scalars/method registration and actual aggregate context removal, with compile fences, runtime reds and exact restore greens. Two additional branches survive the whole projection package and are queued coverage findings below.
* **D:** `git diff C..D` contains no own production behavior. Its conformance controls are filed under B/C ownership; own normal D tests and actual aggregate mutation demonstrate useful coverage. There is no fabricated D-production revert.

## Runtime and artifact audit

Root exact unchanged AIL profile rc0 in23.322s then exact unchanged Go profile rc0 in368.028s, including full normal plain/race packages, were read as raw evidence with immutable candidate metadata; not accepted as a controller verdict. Own focused transitionreg/projection/daemon baseline rc0 in91.570s. Own eleven recorded mutation/reversion experiments include semantic kills, three explicit API/dependency compile refusals and two genuine surviving mutants. Every active test subprocess is terminal; all product hashes restored.

Own aggregate removal: compile0; single RED32.610s, N/K/J1/1/1, admissions2, commits0; batch RED35.304s, N/K/J3/2/2, admissions3, commits1. Both real resolver/tools holds are about2.802s/9.802s, and counter/commit assertions precede timing. Restored single20.00354s and batch20.00183s, including callback return/follow-up recovery and no third dispatch. The broad mutant run also fails the unrelated A2A fixture as documented above. Fixed axes: local darwin/arm64, default parallelism/caches, no injected CPU or memory load. This does not certify other platforms or stress axes.

All619 executor raw text artifacts independently checked for byte count and SHA256: zero mismatches. All43 named mutation records inspected for selected firing assertion and green restoration. There are42 compiled runtime/source/decorator controls plus1 old-pin dependency refusal, not43 distinct production edits. The two aggregate names share the same outer WithTimeout→WithCancel mutation; single and batch effects are reported separately. Per-POST snapshot/no-readmission rows also share a mutation shape. Three response/version decorators calibrate observable upstream contracts. Their limited scope is not hidden.

## Nonblocking findings and evidence limits

1. `host/projection/mcp.go:19`: replacing cfg.MaxCallbacks with literal800 survives selected bounds tests and the entire projection suite. Candidate code is correct. Add actual shared-handler eight-slot saturation, ninth-entry refusal until callback RETURN, then recovery. No claim that the current daemon is unbounded; no claim that this was tested in whole daemon.
2. `host/projection/mcp.go`: replacing nil-coordinator refusal with successful empty result survives all MCP and whole projection tests. Candidate code correctly refuses; daemon supplies a coordinator. Add a valid admitted library invocation with nil coordinator and exact failure expectation.
3. `TestMCPSurfaceRefusalRecovers` does not pin exact id/message/status; its live65-byte fixture is segmented, while the escape-heavy unit fixture violates the registry segment grammar. Strengthen that specific support, retaining same-handler recovery.
4. `TestMCPQuickstartPayloadsVerbatim` executes three guide payloads verbatim against a real ephemeral daemon and pinned interpreter; setup and headers are reconstructed. It does not execute every guide command verbatim. The guide lacks a runnable initialize command.
5. `projection.go:156` still describes two routes and no spawned worker; MCP adds a third route and uses upstream Runner callbacks.

Per verification3n, the two unpinned supporting branches are queue findings, not a reason to quietly expand this sprint. No production defect was observed in those branches. Complete CI and final administrative design/plan acceptance remain controller obligations.

## Concrete rubric score

**49/100 (48.545 before rounding), FAIL, round1.** The hard failure is the independently reproduced pristine test failure, independent of score. Literal frozen-artifact scoring is deliberately distinct from the14 named behavior assessment.

| Category | Score | Evidence |
|---|---:|---|
| Tests | 0/20 | Own pristine A2A repetition rc1; earlier full green retained, not substituted. |
| Lint | 10/10 | Exact candidate full-repository vet rc0 with empty output and clean formatting. |
| Acceptance | 24.545/30 | Sprint JSON18/22 feature criteria under passes=true; D passes=null, its4criteria not falsely marked accepted. |
| Code quality | 0/15 | Literal rubric deducts5 for each touched file over800lines: daemon.go966, daemon_test.go1343, projection_test.go1553. All three were already over800 at base (963/1341/1541); inherited size, not introduced regression. No unrelated splitting requested. |
| Documentation | 5/15 | Runnable Go API examples/real-daemon guide payloads earn equivalent host-example5. No host changelog entry; frozen design header remains Planned/implementation acceptance NOT RUN. |
| Design fidelity | 9/10 | Thin upstream-owned wire architecture matches; full runbook execution is narrower than plan wording. |

No compiler-infrastructure or optimization bonus category is triggered. Pending final acceptance metadata is preserved honestly; this evaluator neither marks CI green nor mutates frozen artifacts to improve a score.

## Terminal disposition

EVALUATION_RESULT: fail

EVALUATION_SCORE: 49/100

EVALUATION_ROUND: 1

Blocking finding: BLOCK-1, inherited but currently failing pristine A2A socket/receipt fixture. All own test processes terminal, both worktrees retained, candidate production diff empty, only evaluator Markdown/JSON untracked. Review completed before21:45UTC. No production fixes, no descendants, no unreported token estimate. Required repair must preserve the fixture law and be independently reviewed on a new immutable candidate before acceptance.
