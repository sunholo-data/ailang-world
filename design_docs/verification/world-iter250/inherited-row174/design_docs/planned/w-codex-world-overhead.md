# w-codex-world-overhead — concise se-tools descriptors

**Status:** Planned; historical profiling complete, product implementation and smoke validation pending.
**Created:** 2026-10-10, attended session; not a scheduled mission iteration.
**Queue / authority:** row 174; D-WORLD-70 = A. Clause 4 remains UNMET.
**Priority:** P0 critical path before a new row-93 FINAL.
**Target:** next World product change; no version or release-policy change.
**Estimate:** 1 day including metadata review, equivalence controls, and verification; smoke depends on available providers and attended scratch sessions.
**Dependencies:** existing row-93 M4 smoke rig and tuning ledger; fresh independent evaluation before landing; new committed preregistration and attended mint before FINAL.

## Problem and measured scope

The immutable 2026-10-09 FINAL failed the Codex overhead criterion: median **per-task** overhead `4220/15879 = +26.57598%` exceeds +25%. Codex passed the separate pass-rate criterion at −1/69; Claude passed both, with overhead −19.65%. These are the existing verdict, not recomputed thresholds. D-WORLD-70 accepts this FAILS and authorizes a World product repair followed by a new FINAL.

The smallest defensible first repair candidate is **shortening the eight AILANG tool descriptions in `packages/se-tools/transitions.json`**. This removes a measured package-owned prompt payload while keeping the tool contracts and execution path identical. Historical evidence does **not** establish that descriptor size caused the failed overhead. This proposal must earn that claim through fresh smoke evidence; byte savings alone do not satisfy row 174.

Offline profiling is reproducible with:

```sh
python3 design_docs/verification/world-row174-attended-2026-10-10/profile_final.py
```

The profiler reads only the banked FINAL and the current package manifest. It asserts 69 distinct task/trial rows per Codex arm, 23 tasks, one completed usage record per transcript, and every transcript's banked SHA-256. Its output is [final-profile.json](../verification/world-row174-attended-2026-10-10/final-profile.json).

| Observation | Shell | World | Interpretation limit |
|---|---:|---:|---|
| Task-run rows | 69 | 69 | Three trials on the fixed 23 tasks |
| Pooled median wall, context only | 28,525 ms | 35,940 ms | The gate uses median per-task ratios, not this pooled ratio |
| Median input tokens | 126,334 | 261,841 | Aggregate usage, not initial prompt size |
| Median cached input tokens | 93,184 | 214,784 | Provider-accounted cached usage; no latency attribution |
| Median uncached input tokens | 34,605 | 47,034 | Median of each row's input minus cached usage |
| Median output tokens | 717 | 658 | Similar output volume |
| Recorded `turn.started` events per row | 1 | 1 | Codex session-level event; internal model request count is unobserved |
| Completed command/file or MCP items | 223 commands + 67 file changes | 245 MCP calls | Distinct tool surfaces; not interchangeable latency or model-round measures |

World calls comprise 69 reads, 76 writes, 70 runs, 20 CLI calls, 6 checks, and 4 builtin searches. All 245 saved result objects carry equal JSON twice: text content and `structured_content`, 283,896 bytes in each representation. Canonically encoded complete results total 629,416 bytes; median 1,131 bytes. The upstream MCP transport owns that representation, so changing it is outside this patch. CLI results account for 258,366 bytes (41.0% of that total), principally whole-module `iface` responses. `pipeline` has no CLI call and still has +93.51% per-task overhead; selectively trimming `iface` cannot explain or repair every slow case by itself.

The current eight granted descriptors reconstruct to **13,947 canonical JSON bytes**, including **7,430 UTF-8 description bytes**, 2,824 input-schema bytes, and 3,106 output-schema bytes. This is a reconstruction through the same fields `mcpDescriptors` emits, **not a banked tools/list response from FINAL**. The FINAL teaching prompt is 97,530 bytes and stays unchanged. Byte counts are neither token counts nor predictions of wall time.

The controller's banked current pinned, single-client daemon probe completed 18 calls (six each of read/write/run): reads 135.6–172.4 ms, writes 136.5–147.2 ms, runs 174.1–500.3 ms. Standalone full binary read+SHA-256 measured interpreter 39.7–44.7 ms and tool 41.2–71.1 ms. These are current idle measurements in [daemon-probe.log](../verification/world-row174-attended-2026-10-10/daemon-probe.log), not historical per-phase timings and not a loaded benchmark. They identify small fixed costs; they do not plausibly attribute the full historical 7.4-second pooled-median gap to daemon work. The separately captured 16,991-byte HTTP tools/list contains nine tools, including `workspace-exec`, and must not be reported as the exact floor eight-tool wire payload.

## Root-cause candidates and disposition

| Candidate | Evidence | Decision |
|---|---|---|
| Repeated tool descriptions/schema payload | Eight descriptors are 13,947 B, with 7,430 B descriptions; package owns these descriptions | First controlled candidate: compact descriptions only; efficacy pending |
| Large/duplicated tool results | Exact duplication in all 245 saved results; CLI contributes 41.0% | Preserve results in this patch; duplication belongs to released protocol machinery, CLI output has its own semantics |
| Interpreter/tool re-hashing (row 165) | Both paths verify every execution; current standalone costs about 38–45 ms each | Preserve integrity semantics; no verify-once cache or TOCTOU policy in row 174 |
| Building undeclared handler families (row 168) | Binder opens episode registry before effectful dispatch; episode AILANG handlers are cached after success | Separate product item; all floor tools require the AILANG family, and FINAL uses per-run episodes; no evidence of repeated construction failure here |
| Extra agent/model rounds | World has fewer completed tool items; single `turn.started` in both arms | Unproven: transcripts lack internal request timestamps/counts; do not infer from total input tokens |

This is a shared descriptor repetition pattern across the eight tools, not a `pipeline` special case. Treat all eight descriptions consistently. `workspace-exec` is outside the floor grant set and keeps its current description in this change.

## Goals and solution

Reduce the canonical eight-tool descriptor list to **at most 10,500 bytes** and total top-level description text to **at most 3,800 UTF-8 bytes**, from 13,947 and 7,430 respectively. Keep all nine manifest entries and every field other than the eight designated top-level descriptions equal to the pre-change baseline.

Rewrite each description around the operation, how to select arguments, decision-bearing exceptions, failure interpretation, budget, and World provenance. Retain information needed to call correctly rather than repeating an implementation walkthrough. The unchanged input/output schemas already describe argument types, defaults/limits in property descriptions, and return fields. Complete operator/confined-execution explanations stay discoverable in QUICKSTART; add a concise tool-use table if review finds a removed detail absent there.

For example, the read description can state:

> Read an episode-worktree file. Use a non-empty relative path without `..` segments or a leading `-`; policy confines symlinks and protected paths. Unknown keys refuse. Returns file content or a refusal plus World plan/effect refs; costs one Workspace.Read and commits one log entry.

The run description must retain the default exact IO/FS capability set, Env/Net mutual exclusion and grant/operator enablement, loopback Net restriction, `stdin`/`argv` behavior, timeout, denied-program non-execution, and the warning that an inner refusal can exit zero. Preserve the decision/action guidance (`decision.missing_from_policy`, inspect stdout/stderr). Avoid suggesting that requesting caps can create authority. The edit description must retain exact-once replacement; check must retain structured diagnostic behavior; builtin search must retain its degraded inventory shape and refusal on invalid/truncated inventory; CLI must retain its allowlisted operation names and write refusal. Property descriptions remain byte-for-byte unchanged.

**Why is this a package?** The change belongs entirely to published se-tools metadata. No kernel code or World transition law grows. The host remains the effect boundary and the package remains the machine-facing tool specification. Publication uses the existing transitions publisher and attended fence; this design grants no new publishing/session authority.

## Conflict surface and preserved behavior

This patch touches no parser, typechecker, compiler, interpreter, coordinator, binder, capsule, MCP transport, grader, budget, or result schema. It reuses the existing projection: `mcpDescriptors` forwards `Description`, normalizes the existing input schema, and copies output schema. Both MCP and A2A consumers read registry metadata, so changes are not client-specific. MCP fresh admission and the broker remain authoritative.

The baseline comparison includes `id`, `title`, `transitionFnFile`, complete input/output schema objects including property descriptions, `access`, and `declaredEffects`; it also includes any future/unknown field rather than comparing only a hand-picked field subset. It compares the full entry with only the approved top-level `description` key removed. `workspace-exec` compares as a complete entry. No transition `.ail` bytes change; existing source/effect-plan hashes must remain equal, while a descriptor registry revision may change because its prose changed.

The existing `TestF1PinnedInterpreterHashMismatchRefusedBeforeExec` mutates archived interpreter bytes and asserts a hash-mismatch error. `TestWorkspaceToolBinaryIsArchivedAndVerified` mutates the archived tool after a successful call, asserts refusal, and checks dispatch count unchanged. Preserve both; their bodies were read, not inferred from names. The existing se-tools real-capsule matrix remains the execution regression gate.

## Implementation plan and files

1. Freeze the exact baseline manifest from reviewed parent `5ca1115` in a test fixture. Implement metadata equivalence and compact-size acceptance before editing prose. Capture an exact eight-grant tools/list from the scratch daemon to distinguish live wire bytes from canonical reconstruction.
2. Rewrite the eight descriptions and review retained semantics against the current manifest and QUICKSTART. Run the baseline equivalence, live projection and source/plan regression checks. Add usage documentation only for detail removed from the descriptor that lacks an existing home.
3. Run mutation controls, normal gates, and fresh independent evaluation. Then run one-knob M4 smoke comparison/repeat with cheap tuning models when available; append evidence to the existing tuning ledger. Run all four smoke cells fresh: the current `live_config`/`arm_digest` omits product manifest descriptions, so `--reuse-from` cannot certify this candidate. Separately bank the manifest SHA, served exact-eight metadata SHA, daemon binary SHA and client CLI versions beside the rig's config digest. Current CLI differs from the historical FINAL (0.159.2→0.162.0); compare baseline/candidate using the same current CLI and do not attribute cross-version differences to this patch. Stop before FINAL until this product change lands and the newly committed preregistration plus attended mint exist.

**Product scope:**

- `packages/se-tools/transitions.json` (~−3.5 KB prose): eight top-level descriptions only.
- `host/coordinator/setools_descriptor_test.go` (~100 LOC): immutable baseline comparison and exact tool census/byte limits; real manifest, not a generated self-baseline.
- `host/coordinator/testdata/se-tools-descriptor-baseline.json` (~16 KB): parent manifest fixture, with the source commit and SHA independently recorded in the test.
- `docs/QUICKSTART.md` (0–40 LOC): usage guidance removed from descriptions only if currently absent; review links and run existing relevant walkthrough.

**Evidence only:**

- `design_docs/verification/world-row174-attended-2026-10-10/profile_final.py`: offline historical profiler, already present.
- `design_docs/verification/world-row174-attended-2026-10-10/final-profile.json`: exact offline output, already present.

No mission charter, ledger, scheduler, shared skill, installed runtime, historical FINAL artifact, floor grader, task/model selection, threshold, or N changes belong in this product patch.

## Acceptance and non-vacuous mutations

| ID | Required assertion | Mutant that must make it red |
|---|---|---|
| AC1 | Independently pinned exact nine-ID manifest census and exact eight-ID floor set; unique IDs | Delete `ailang-read`, duplicate its ID, or add a tenth entry; exact-set assertion fails |
| AC2 | All non-description fields equal to immutable parent fixture; complete `workspace-exec` equal; eight non-empty descriptions actually changed | Add Net to a declared effect, set `additionalProperties:true`, remove `world` from output schema, or alter transitionFnFile; field-equivalence assertion fails |
| AC3 | Canonical eight-tool list ≤10,500 B and eight descriptions ≤3,800 B, each present; baseline reconstruction asserts 13,947/7,430 B | Restore all eight old descriptions; size assertion fails (baseline 13,947/7,430 exceeds both) |
| AC4 | Live exact-eight-grant tools/list emits the manifest's rewritten descriptions and unchanged names/input/output contracts | Make projection return the old descriptions; live description-equality assertion fails |
| AC5 | Existing real-capsule package matrix and archived-binary tamper controls pass; transition `.ail` sources equal parent | Edit a source byte; explicit source-byte-equivalence assertion fails. Disable binary verification; existing tamper assertions fail before exec |
| AC6 | M4 tuning comparison uses the fixed six tasks, same cheap models/config except this one package-description change; all attempts banked; no harness faults; World pass rate preserved for both agents | Supply zero task rows or remove one task/trial; existing rig's census/eligibility assertions fail rather than publishing a favorable empty report |
| AC7 | Fresh smoke and no-change repeat show Codex O ≤+0.20 (headroom target) and Claude O ≤+0.25 with Δ≥−0.02 for each; report same-config noise and no pooling substitution | Replace median per-task O with pooled median ratio; existing report's outlier fixture assertion fails. Feed a known >+0.25 fixture; report must fail |
| AC8 | Required normal verification and independent review pass; documentation updated and new scratch evidence/ledger entry banked | A failed check or pending efficacy is a failure to complete, not an administrative pass |

AC6/AC7 are **pending and not executed in this designer session**: provider/model calls were explicitly prohibited. Scratch publication and session mint additionally require the user's typed fence under QUICKSTART §0; the prepared scratch rig is not a minted session. The headroom target +20% is an internal tuning acceptance goal, not a change to the ratified +25% FINAL threshold. A byte-size pass plus missing/failing smoke cannot close row 174. At the tuning rig's fixed N=1, a repeat is useful evidence of noise, not statistical proof. Record all no-change repeats and remain within existing capacity-fault rules; no selecting a favorable run or changing tasks/models/N to rescue this candidate.

Full gates: `./scripts/verify_ail.sh`, `go build ./... && go test ./...`, plus required repo Go/race gate for the execution milestone, under pinned released binaries. Tests must report discovered controls and zero skips in the relevant package matrix. Historical FINAL remains immutable. The new FINAL is a separately preregistered, attended execution after landing, using the unchanged gate models, task set, grader, thresholds, N and sealing protocol.

## Verification log

| ID | Evidence read/run on 2026-10-10 | Result and scope |
|---|---|---|
| V1 | CLAUDE.md; coding-standards.md; DESIGN §1/§14; mission row174/D-WORLD-70 | Confirmed scope, package-first rule, mutations, immutable benchmark and new FINAL requirement |
| V2 | `profile_final.py` over banked `runs/codex/{shell,world}/r{1,2,3}.jsonl` and referenced transcripts | Passed all row/usage/SHA assertions; statistics banked in final-profile.json |
| V3 | Current `packages/se-tools/transitions.json`; `host/projection/mcpname.go` `mcpDescriptors` | Confirmed schema/description forwarding and byte reconstruction; actual historical tools/list remains unmeasured |
| V4 | All completed historical Codex MCP items, exact text JSON vs structured JSON comparison | 245/245 equal; no attribution to internal model behavior |
| V5 | `host/daemon/workspace.go` `registry`, `ailangHandler`, `verifiedTool.Execute`, `verifyArchivedTool` | Handler cached by episode/sandbox after success; archived tool re-hashed each call. Requested `host/broker/workspace.go` absent; actual workspace implementation is daemon-owned |
| V6 | `host/capsule/capsule.go` `RunContext`, `verifyExecutable`; `cachetemplate.go` `seedCache` | Confirmed full interpreter read/hash then private template cache copy and subprocess; no cache/verification edits proposed |
| V7 | `host/coordinator/effectful.go` `dispatchEffectful`, `composeEffectOutput`; package read/run headers | Plan executes once before effects; existing single-effect, no-finish se-tools return handler output plus World block; no finish cost invented |
| V8 | Bodies of `TestF1PinnedInterpreterHashMismatchRefusedBeforeExec`, `TestWorkspaceToolBinaryIsArchivedAndVerified`; `setools_package_test.go` fixture loader | Confirmed concrete tamper refusals and real-capsule matrix; no assertion about test behavior inferred from names |
| V9 | Mission rows165/168; capsule-output-cap design; git history of capsule/daemon/projection | Related integrity/binder optimizations already queued; this metadata-only change is distinct and does not absorb them |
| V10 | Controller's 18-call pinned daemon probe and standalone read/hash probe, banked in `daemon-probe.log` with reproducible `daemon_probe_test.go.txt` in the row174 evidence directory | Current single-client timings received above. Not per-phase historical trace and not loaded attribution |
| V11 | FINAL preregistration and M4 README | Full teaching 97,530 B; fixed M4 six-task smoke and pinned models; prior same-config O noise explicitly recorded |
| V12 | Related-doc lexical search/read (rows165/168, row93 floor design, se-domain design, capsule integrity design) | No existing concise-descriptor repair found in those candidate designs. Neural/provider search and scaffold auto-search intentionally not run under session's no-provider-call instruction; no similarity score fabricated |
| V13 | `scripts/floor/run.py` `arm_digest` and `live_config`; QUICKSTART §0 | Confirmed config digest lacks manifest descriptions; fresh four-cell comparison and separate product hashes required. Publish/mint requires a human-typed terminal fence |
| P1 | Fresh live exact-eight-tools wire capture | Pending implementation measurement; nine-tool probe is not substituted |
| P2 | Metadata candidate smoke and repeat | Pending availability; byte reduction alone proves no latency improvement |

## Review status, risks, and axioms

The design-doc-creator skill is applied to this document. All four attended quorum triggers are false for this narrow implementation: no human design-freeze decisions, no override of shared machinery, no cost/KPI or banked-data schema changes, and no load-bearing external-system contract premise. Normal independent implementation evaluation still applies. No provider quorum is run here under the explicit no-provider-call scope; this is not a claimed quorum pass.

Risks: concise wording can reduce tool comprehension; AC6/AC7 reject pass-rate or overhead regression. Canonical bytes can differ from actual MCP wire bytes; AC4/P1 bank both. The dominant source could be model/client prompt assembly or provider latency; this design makes no external causal claim and parks efficacy if smoke cannot establish a product improvement. Descriptor registry hashes change naturally; source bytes, schemas, authority and execution remain pinned by AC2/AC5. Avoid broadening this patch into result filtering, prompt rewriting, transport ownership, or hash caching.

| Axiom | Score | Reason |
|---|---:|---|
| A1 Determinism | 0 | Execution/source unchanged; byte comparisons deterministic |
| A2 Replayability | 0 | Recorded plans/effects unchanged |
| A3 Effect Legibility | 0 | Preserve effect-specific advice and exact declaredEffects |
| A4 Explicit Authority | 0 | Preserve schemas, access and operator/session restrictions |
| A5 Bounded Verification | +1 | Small deterministic census/equivalence/size controls with named mutants |
| A6 Safe Concurrency | 0 | No concurrency changes |
| A7 Machines First | +1 | Less repeated prose in machine-facing discovery, pending efficacy |
| A8 Minimal Syntax | 0 | No syntax change |
| A9 Cost Visibility | +1 | Bank exact payload and usage measures, separate cached usage and unknown attribution |
| A10 Composability | 0 | Existing projections/package publishing reused |
| A11 Structured Failure | 0 | Preserve failure advice and structures |
| A12 System Boundary | 0 | Package metadata only; effects stay at host boundary |

Net +3; no A1/A3/A4/A7 hard violation. Language-support claims and new diagnostic codes are outside this design.

## Related documents

- [Resident-agent floor run](w-resident-agent-non-inferiority-floor-run.md): owns the frozen benchmark, M4 tuning and preregistered FINAL; this change repairs the product metadata, not the benchmark.
- [Software-engineering domain](w-software-engineering-domain.md): defines se-tools' package/effect contracts; this change preserves them.
- [Capsule integrity/load design](w-capsule-output-cap-load-flake.md): explains why verification hoisting needs separate security work; this change leaves it untouched.
- [Historical FINAL README](../verification/world-floor-final-2026-10-09/README.md): immutable failed result.
- [M4 smoke README](../verification/world-floor-m4-2026-10-06/README.md): tuning models, fixed smoke set, ledger location and observed noise.

No clause closure, mission iteration bookkeeping, publication, session mint, provider call or FINAL rerun is represented by this design deliverable.
