# Independent PRODUCT evaluation — iteration220 D-WORLD-49A

**FAIL — 59/100 (exact 58.5454545), mandatory full-test hard failure.** Immutable candidate `cdb714e19811cb067d6fa0f00c135a1216876633`, compared with `origin/dev` `c5c40f3e05dad27b7591ca494bfeea2e7f1920bd`. This is the one authorized D49A follow-up, not a restart of the historical three-round loop. FAIL49/54/59 remain historical failures.

Judge: native `gpt-6.1-sol`, fresh separate context, FLAG `same-model-fresh-context` under amended D48. Preferred Sonnet returned Unknown model according to controller routing; no Astra used. Native token usage UNKNOWN/null. Judge did not generate code or diagnosis. Own isolated tree `.wt-world-iter220-eval`; no production/test fixes, commits, pushes, main/fleet writes, document moves, or subagents. Absolute workdir pinned on each scoped command. Runtime started only after controller terminal authorization.

Read sprint-evaluator skill/full100-point rubric, full mission verification protocol, World CLAUDE/charter/coding standards and DESIGN §§1/14, entire production diff, original fourteen ACs, original five-milestone sprint/plan and three historical judgments, new designer/plan/sprint/diagnosis and raw metadata. Go code/module/gate bytes have no diff between final historical candidate `741e3c3` and current candidate; evidence normalization/authoritative-main merge does not repair product behavior.

## Blocking measurements

Executor exact sequential AIL→Go: AIL0/5.747s; Go1/296.061s. Build,37-test evidence manifest, all plain packages passed. Full race reached: `TestQueryInterfaceReturnsWhileADescendantHoldsStdout`, `host/pkgproj/iface_test.go:232`, missing pid artifact after3.00s. The fixture's3s timeout returned before the required pid was observed; why it was absent is UNPROVED. Untouched source and old comment about1s startup failures cannot prove this failure's cause.

Controller separate exact AIL→Go: AIL0; Go1/325.175s. All plain packages passed. Full race reached two named failures:

- `TestCLIRealSubprocessEpisode`, `cmd/ailang-worldd/cli_test.go:211`: daemon announcement timed out with empty stderr. This is startup, unlike historical round3's line266 replanned commit503. Do not conflate signatures.
- `TestOrdinaryOverflowCarriesNoESRCH`, `host/broker/cleanup_test.go:200`: joined overflow plus `overflow kill: operation not permitted`, instead of exactly `*HandlerOutputOverflowError`. Round2 recorded the same assertion/signature, but no controlled cause or permissible correction is established.

Both executor/controller full Go measurements are terminal rc1 at the identical immutable candidate, normal Go1.26.6/darwin-arm64 environment/default caches and unchanged scripts, with both interpreter variables pinned to released v0.41.0. Expected race-detector control warnings are instrument-positive evidence, not product failures. Raw logs/metadata and controller lossless base64 hashes match. Own exact measurements: AIL0/23.727s then Go1/159.666s. Build and37-test evidence manifest passed; **plain** full suite failed the same descendant-pid assertion (`iface_test.go:232`,3.00s). Race **NOTRUN** under the unchanged fail-fast script. This corroborates the executor's named signature in another suite leg; it does not establish a unique cause or a race-only property. All own processes terminal; no cap fired. Own `go vet ./...` rc0 with empty output and changed-Go `gofmt -l` empty.

Eight paired four-test commands (base first, candidate second) all rc0. They establish no recurrence in that selected execution shape. Only capsule's historical exact elapsed assertion is demonstrably present on an unchanged prior base. Coordinator family failures differ in classification/elapsed signature; CLI commit and daemon ordinary-B failures lack exact named prior-base recurrence. These greens neither waive full reds nor prove a unique environmental cause. Historical round3's “inheritance”, “only under saturation” and “not MCP regressions” conclusions exceed its measurements; retain its verdict, reject that causal overreach.

**Accepted full profiles have not been established without causal correction.** A later green would document that run, not erase the three unresolved full-profile reds or prove a repair. D49 allows verified fixture correction, but this evidence demonstrates no source-proved correction preserving every assertion. Do not enlarge clocks, serialize verify_go, skip tests, adjust caches/GODEBUG/toolchains, remove pid/error-classification assertions or manufacture a fixture patch. Any required production repair must return to normal design. One bounded follow-up does not authorize speculative repair/review loops.

## Original contract assessment

All fourteen ACs remain binding. Status here means functional source/banked-control support, never product acceptance; final full gate/CI remain failed or unmeasured.

| AC | Current support and limits |
|---|---|
| DEP | v0.47.2 narrow protocol dependency; actual historical closure54→56 adds only hostcall/mcphttp, no removals; old-pin compile refusal distinct from runtime kills. |
| MAP | Canonical reversible escaping, invalid escape rejection,64-byte whole-surface refusal and same-mount recovery. Exact complete refusal-envelope assertion and a live registry-valid escape-heavy case remain incomplete support, previously nonblocking. |
| SCHEMA | RawMessage preserves9007199254740993/properties/constraints; existing object bytes unchanged; nonobject/annotation refusal retained. |
| CARRIER | Shared bearer resolver, constant401 denial, no denied registry reads plus populated authorized control; middleware exclusion unchanged. |
| ADMIT | Fresh Allowed set per session/item, dynamic registry/error refusal and inner read bounds supported. |
| ABSENCE | Only live-context typed confirmed absence is empty; absent precheck plus genuine same-text/store/context faults propagate. MCP list/invoke and A2A compatibility controls observe injected stage/sink/no mint/no dispatch. |
| INVOKE | Object input and fresh descriptor membership/context check precede random task mint and real coordinator. Cooperative recovery fixture retains same handler/one slot and actual follow-up success. |
| ITEM | Repeated RPC IDs still yield distinct64hex tasks, K+1 snapshots, receipts/journal entries, sequential overlap1; revocation stops later item. |
| BATCH | Upstream owns old-version arrays/new-version400, notifications202, item continuation and whole-host-error abort; first commit retained/third unentered. Decorators are honest wire-calibration tests. |
| BOUNDS | Frozen3/10/20/30 and slots8 injected; actual single/batch production controls retain stage/task/journal counts and successful cooperative follow-up. Named aggregate mutations record32.61s/34.62s failures with actual entered stages and restoration green. |
| WIRE | Thin context-only World wrapper delegates JSON-RPC/MCP/SSE to released seam; generated upstream golden, transport/body/error controls retained. Slow-body30s truncation is a residual, no universal20s delivery claim. |
| ROUTE | Executable AST asserts exactly POST /mcp/; methodless mutant kills this assertion while independent GET405 still passes. |
| CROSS | One daemon, unequal populated sessions, exact encoded/decoded tool/card cardinality including dot/slash examples. |
| FIXTURE | Both surfaces have nonempty controls and recorded ambient union absent; actual three guide payloads execute against real daemon/pinned interpreter. Setup/headers reconstructed, initialize recipe absent. |

All619 executor original artifact hashes verified;29 display-normalized logs recovered from lossless originals, zero unresolved mismatch. Original43 named mutation records include42 compiled runtime/source/decorator controls and1 dependency compile refusal. Shared aggregate/snapshot mutant shapes are not falsely43 independent production behaviors. Prior stage-specific API-removal compile refusals and semantic reversions remain banked; no own mutation claimed this follow-up. D has no production behavior diff; its tests pin B/C and are not credited with an imaginary D-production revert.

## Literal rubric

| Category | Points | Reason |
|---|---:|---|
| Tests |0/20| Exact mandatory full Go red; HARD FAIL. |
| Lint |10/10| Own changed-Go gofmt list empty; own full-repo go vet ./... rc0 with empty output. |
| Acceptance |24.5454545/30| Original sprint A1(5)+A2(6)+B(3)+C(4)=18/22; D(4) remains null. New diagnostic milestone false is honestly unresolved, not completion. |
| Code quality |0/15| Literal three -5 oversized touched files: daemon.go966, daemon_test.go1343, projection_test.go1553. Already oversized at base; no unrelated splitting requested. |
| Documentation |15/15| Unreleased host changelog, applicable runnable Go host examples/guide payloads, truthful unaccepted design status. |
| Design fidelity |9/10| Thin released-seam/host boundary matches original design; full usage/refusal support narrower than strongest wording. |

Total58.5454545 rounds59. Threshold70; hard test failure independently overrides score. No compiler regression-surface bonus triggered; production clocks are correctness bounds rather than optimization goals, so performance bonus not triggered. Correct diagnostic caution is not penalized or promoted into product success.

## Disposition

B1/B2/B3 above block acceptance. Action is a source-proved, measured correction within D49's remaining authorized bounds or an attended scope decision; no default new repair round is licensed. Required candidate-SHA remote CI is UNMEASURED in this judge and receives no acceptance credit. Clause6 remains UNMET; row108 stays unaccepted/unmerged, so dependent rows114/93 gain no claimed movement.

Nonblocking limits remain: row132's callback-capacity and nil-coordinator supporting-hunk survivors, precise long-surface test support, guide command/setup scope and stale Handler comments. No production defect inferred from those survivors, no widening. SSE golden's intentional terminal blank event delimiter explains diff-check blank-line warning; do not rewrite frozen wire bytes for whitespace cosmetics.

EVALUATION_RESULT: fail
EVALUATION_SCORE: 59/100 (exact58.5454545)
EVALUATION_ROUND: D49A bounded follow-up1; historical rounds1/2/3 retained
EVALUATION_REPORT_PATH: design_docs/verification/world-iter220/product-evaluation.json
FEEDBACK_SUMMARY: Exact full Go still red in executor/controller/independent judge; no causal correction or waiver, clause6 UNMET.
ALL_OWNED_RUNTIME: TERMINAL
