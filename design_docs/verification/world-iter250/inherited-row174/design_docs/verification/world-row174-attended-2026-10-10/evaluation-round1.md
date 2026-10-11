# M1 independent evaluation — row 174

Candidate: `604cee1`, parent: `5ca111547fdd80480a8d42578fb59942ee8a0b63`. Fresh isolated checkout: `/private/tmp/ailang-world-row174-evaluator-20261010`. Evaluator role fallback is **same-model-fresh-context** under amended D-WORLD-48; this is not cross-vendor independence. Scope is M1 only.

## Current verdict

**M1 evaluation FAIL — 68/100; normal-gate requirement unmet.** Technical descriptor implementation review passes, but required normal verification does not. First full gate completed rc=1 (`verify-go-first-red.log`): broker `TestRunBoundedHeadTailTimeoutKeepsPartialOutput` failed with empty partial stdout; all other plain packages passed. Untouched parent reproduces the unchanged timing-sensitive assertion: 13 PASS / 7 FAIL across 20 runs (`baseline-partial-output-control.log`). One unchanged full retry passed the plain broker package (119.249s), then failed the same test in the race leg (170.087s; `verify-go.log`, lines 97–100). Remaining race packages were still running when root requested this round be finalized; the observed failure already makes the full gate nonzero. No failed gate is waived. This is a gate failure, not evidence of a descriptor semantic defect. A separate fixture repair and new committed-candidate evaluation are required before M1 approval.

## Independent evidence

- Immutable fixture bytes equal `git show 5ca1115:packages/se-tools/transitions.json`; SHA-256 `240840ca9a0f8f4c60ac159dd78fce26a6c681200e8a10c64fe0e08c0a228868`. Full decoded entries compare equal after removing only the eight permitted top-level descriptions; workspace-exec compares complete. All 28 tracked `.ail` files compare byte-for-byte to parent.
- Independently ran three discovered tests with pinned interpreter v0.41.0 and tool v0.52.1: compaction, transition source equality, real HTTP exact-eight tools/list. All PASS; zero skips. Canonical payload **13,947→10,301 B**; descriptions **7,430→3,794 UTF-8 B**; live SSE **10,368 B**. Wire and canonical sizes are distinct. Independently parsed and SHA-checked both banked exact-eight wire captures against their committed manifests: parent 14,054 B, candidate 10,368 B; complete names/descriptions/input/output equality holds in each.
- Own isolated mutation restored all eight old descriptions: test exits 1 at `candidate byte budget: 13947/7430 exceeds 10500/3800`, plus all eight unchanged-description assertions. Restored exact original bytes and reran all three tests: GREEN, zero skips. Product tree remains unchanged.
- Independently checked all 14 executor mutant logs against intended failing assertion and nonzero return; executor bank records exact restoration and pristine GREEN. Census deletion/duplication/tenth-entry; effects/schema/output/source/unknown-field/workspace changes; old prose; stale projection; and both archived-binary verification bypasses are detected.
- Controller scoped log contains 29 top-level and 52 subtest PASS identities, zero SKIP: real capsule matrix, grant/budget/replay and binary tamper controls plus session exact-set control. These are observed controls, not name-only claims. Existing optional full-suite skip conventions are separate from this zero-skip matrix.
- verify_ail bank: 16 required identities, 40 named tests, 444 se-tools named tests, 9/9 package steps; no `.ail` changed. New Go files are gofmt clean and 141/71 lines.

## Semantics and design review

Descriptions retain exact IO/FS defaults, Env/Net exclusion and separate grants, operator enablement, loopback Net/redirect restrictions, stdin/argv and args_json, 8-second timeout, declared-effects/policy admission and FS budget-zero behavior, denied non-execution, decision.missing_from_policy action, and exit-zero inner-refusal warning. Edit exact-once, structured check diagnostics, builtin degraded inventory/invalid-or-truncated refusal, complete CLI allowlist/write refusal, path constraints, unknown-key refusal, per-effect cost and World provenance remain. Complete property prose and contracts are unchanged. QUICKSTART already supplies operator setup, cap/session/net guidance, degraded builtin shape and payload examples; no new surface or example module is needed.

Implementation confines production changes to package metadata; added host code is verification only. No parser/compiler/effect-row infrastructure changes trigger the conditional compilation regression category. Metadata budget margin is only 6 description bytes; this is a nonblocking maintenance observation.

## Rubric (M1 milestone scope)

| Category | Score | Evidence/qualification |
|---|---:|---|
| Tests | 0/20 | HARD FAIL: scoped GREEN, full normal retry race leg red |
| Lint | 10/10 | Added Go files gofmt clean; no production Go change |
| Acceptance | 25/30 | Five M1 criteria met; sixth full gate/review criterion failed |
| Code quality | 13/15 | Small focused tests; generic artifact completion deduction limited to 2 because M2 legitimately pending |
| Documentation | 10/15 | Existing examples adequate; design/planner status still pending; no changelog change |
| Design fidelity | 10/10 | Description-only package repair exactly scoped |

Final M1 score **68/100 — FAIL**, with the Tests hard-fail condition independently controlling rejection. A future clean normal gate would permit reassessment; it is not credited now. Performance verification is deferred to M2, not awarded or claimed PASS; historical profiling and metadata bytes do not prove latency efficacy. Generic whole-sprint completion cannot be declared because M2 remains pending.

## Blocking and pending work

Blocking for M1 final approval: repair the timing-sensitive verification fixture in a separately reviewed change and complete clean normal build/plain/race verification on the new candidate. Both red attempts remain banked; parent reproduction does not waive the gate. Nonblocking: update M1 completion artifacts after accepted final verification while preserving M2 pending.

M2 remains **pending**, requiring the human-typed QUICKSTART publish/mint fences and capacity for fixed Haiku4.5 / gpt-6-luna smoke models. No role fallback changes those models. Fresh fixed six-task baseline/candidate/candidate-repeat, four cells, N=1, same current CLI/config, all attempts and hashes banked, unchanged grader/thresholds/capacity policy are required. No reuse-from, pooling substitution or favorable-run selection. Codex headroom <=+20%, Claude <=+25%, each delta>=-0.02 and pass-rate preservation remain unmeasured. Row 174 and clause 4 stay unmet; no FINAL authorization or performance PASS follows from this report.


EVALUATION_RESULT: fail
EVALUATION_SCORE: 68/100
EVALUATION_ROUND: 1
EVALUATION_REPORT_PATH: design_docs/verification/world-row174-attended-2026-10-10/evaluation.md
FEEDBACK_SUMMARY: M1 metadata controls pass, but required full Go retry race gate fails unchanged partial-output timeout fixture; repair separately, rerun full gate, then evaluate new candidate. M2 efficacy remains pending.
