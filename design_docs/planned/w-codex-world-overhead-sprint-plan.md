# Row 174 — compact se-tools descriptions sprint plan

Ready for implementation, attended 2026-10-10; benchmark efficacy pending.
Design: [w-codex-world-overhead.md](w-codex-world-overhead.md).
Base: `5ca111547fdd80480a8d42578fb59942ee8a0b63`.

Two milestones separate implementation proof from benchmark efficacy. M1 cannot close row174.
Historical FINAL stays immutable. New FINAL requires landed product, newly committed
preregistration and attended mint; it is outside this implementation handoff.
Planner lane derived `codex:gpt-6.1-sol`, `declared:planner-lane-default-pin`;
native inherited-model fallback FLAGGED under D-WORLD-48, not an independently enforced pin.
Sonnet unavailable until API reset. All four quorum triggers false; no quorum pass claimed.
Fresh implementation evaluation remains required, with any available-lane fallback flagged.

## Scope and reuse

M1 contributes to existing `packages/se-tools` metadata; no new dependency or kernel law.
Released v0.52.1 `ailang pkg search se-tools` attempted; registry DNS unavailable
(`storage.googleapis.com: no such host`). Existing local manifest/design inspected;
no alternative package needed. M2 reuses the existing M4 runner/grader unchanged.
JSON follows the shared sprint-planner schema directly, avoiding the generator's automatic
GitHub-message import. No provider inference/messages, registry send, PR or shared-skill edit.
Editable product: `packages/se-tools/transitions.json`,
`host/coordinator/setools_descriptor_test.go`,
`host/coordinator/testdata/se-tools-descriptor-baseline.json`, missing usage advice only
in `docs/QUICKSTART.md`; focused daemon/projection test if needed for real tools/list proof.
Evidence: `design_docs/verification/world-row174-attended-2026-10-10/`.
Preserve `.ail`, charter, tasks/models/N/thresholds/grader/shared harness/teaching prompt,
result/schema/transport, binary verification and execution authority.

## Milestones

### M1: Description-only compaction and verification (~160 LOC)

Dependencies: none. Estimate: 3–5 active hours plus full-gate elapsed time.
Freeze fixture directly from reviewed parent with `git show`; independently pin commit/SHA:
`240840ca9a0f8f4c60ac159dd78fce26a6c681200e8a10c64fe0e08c0a228868`.
Compare complete raw entries with only approved eight top-level description keys removed;
unknown fields and schema-property prose count. Compare workspace-exec completely.
Review descriptions against design's retained semantics: run default exact IO/FS caps,
Env/Net exclusion, grant/operator requirements, loopback Net, stdin/argv, timeout,
non-execution on denial, exit-zero inner refusal, decision.missing_from_policy and stdout/stderr;
exact-once edit, structured check diagnostics, degraded builtin inventory with invalid/truncated
refusal, CLI operation allowlist/write refusal; preserve argument, budget and World guidance.
Add operator advice only if absent from QUICKSTART; execute relevant existing payload examples.
No new example module or syntax gate: no `.ail` is written.

- [ ] Independently assert unique exact nine IDs: ailang-read/write/edit/check/run, builtins-search, examples-search, ailang-cli, workspace-exec; exact floor excludes only workspace-exec.
- [ ] Complete non-description entries and complete workspace-exec equal parent; eight non-empty descriptions changed; every transition .ail source byte equals parent.
- [ ] Parent canonical descriptor/description bytes =13,947/7,430; candidate <=10,500/3,800 UTF-8 bytes.
- [ ] Live tools/list with exactly eight grants emits exact census/new descriptions/unchanged names and schemas; bank wire bytes separately from canonical bytes. Nine-tool probe is not substituted.
- [ ] Named mutations fail intended assertions; restored controls pass; relevant real-capsule matrix discovers tests with zero skips.
- [ ] Full normal gates and fresh evaluator pass; report metadata proof without efficacy/row-closure claim.

### M2: Fresh smoke comparison and no-change repeat (~40 LOC evidence)

Dependencies: M1, fixed tuning models available, human-typed publish/mint under QUICKSTART §0.
Prepared main-repo `.ailang/state/row174-attended-2026-10-10/rig`, port7646, is not minted.
Agent must not type the fence. Use existing `run.py iterate` and ledger guards; no reuse-from:
live_config/arm_digest omit descriptions. Run baseline, candidate, then unchanged candidate
repeat, all four agent/arm cells fresh; keep same current CLI/config across the comparison.
Fixed models Haiku4.5 (`claude-haiku-4-5`) / `gpt-6-luna`; role-lane fallback cannot alter them.
N=1; deadline600s; Claude cap$1/task; canary fizzbuzz; six tasks pipeline, cli_args,
api_call_json, prompt_injection, effect_tracking_io_fs, higher_order_functions.
Bank all attempts/canaries/faults under existing capacity policy and ledger; separately pin
manifest SHA, served exact-eight metadata SHA, daemon SHA, interpreter/tool/client versions
beside config digest. Commit previous ledger row before another iteration, as existing guard requires.

- [ ] Baseline/candidate/repeat contain all tasks/cells, no harness faults, World pass rate preserved for both agents.
- [ ] Candidate/repeat median per-task O<=+0.20 Codex and <=+0.25 Claude; each delta>=-0.02. Report same-config noise and every run, no favorable selection or pooled substitution.
- [ ] Existing zero/missing-row, pooled-ratio/outlier and failing-threshold controls remain load-bearing.
- [ ] Missing/failing smoke leaves M2 pending/failed; byte reduction cannot close row174 or authorize FINAL.

## Verification and mutation handoff

Use released v0.41.0 for AILANG_BIN/WORLD_PKG_AILANG_BIN; released v0.52.1 for
WORLD_TOOL_AILANG_BIN. Scope new descriptor/live tests, TestMCPListExactSetPerSession,
se-tools daemon tools/list/QUICKSTART checks, coordinator TestSeTools matrix,
TestF1PinnedInterpreterHashMismatchRefusedBeforeExec and TestWorkspaceToolBinaryIsArchivedAndVerified.
Verbose discovered PASS controls, zero relevant SKIPs; test names alone are not live proof.

| Isolated mutant | Assertion that must fail |
|---|---|
| Delete/duplicate ailang-read or add tenth entry | Independent exact manifest/floor census |
| Add Net effect, allow additionalProperties, remove output world, change transitionFnFile | Complete non-description equivalence |
| Alter workspace-exec description or add unknown entry field | Complete workspace-exec/full-entry equivalence |
| Restore all eight old descriptions | Candidate byte limits |
| Projection serves old descriptions | Live tools/list description equality |
| Edit one transition .ail byte | Parent source-byte equivalence |
| Disable archived interpreter/tool verification | Existing tamper refusal and unchanged-dispatch controls |

Bank intended failure identity, restore from backups, rerun pristine controls; no weaker assertions.
Full gates: verify_ail.sh, go build ./... && go test ./..., verify_go.sh including race.
Controller already banked verify_ail PASS:16 required identities,40 named tests,444 se-tools,
9/9 package steps. Executor proves `.ail` equality; controller/evaluator runs full Go gate
on candidate after scoped work. Evidence must identify candidate tree, not infer from parent.
Day1 M1; Day2 M2 when model/human prerequisites exist (availability delay unbounded).
Total200 LOC excludes ~16KB immutable fixture, ~3.5KB prose deletion and logs. Recent79
commits/7d mix records/product/huge FINAL artifacts, so no comparable LOC velocity is claimed.
Acceptance uses named tests/mutations rather than invented coverage percentages. Risks:
wording comprehension, wire/canonical confusion, N=1 noise and model/fence availability.
Controller owns commits, executor handoff and landing; executor returns M1 proof plus honest M2 status.
