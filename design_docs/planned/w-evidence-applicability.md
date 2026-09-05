# W-EVIDENCE-APPLICABILITY: Evidence About the Current Requirement

**Status:** Proposed design — queue admission requested by Mark, 2026-09-05; implementation not approved; quorum not run.
**Author:** Astra.
**Target:** World domain-package experiment; supports clause 5 without adding a release requirement.
**Priority:** P1 candidate, after review and existing authority/proof prerequisites.
**Estimated:** 3–5 engineering days, provisional; production integration sized separately.
**Parents:** [HUMAN-SURFACE](../HUMAN-SURFACE.md), [proof authority](../implemented/w-validated-proven-evidence-boundary.md), [coding standards](../coding-standards.md).
**Cross-project origin:** [Astra vision and World correction](https://github.com/sunholo-data/ailang/blob/e1b23db03/design_docs/planned/m-astra-vision.md).

## Problem and coverage decision

An authentic proof can remain valid about its original subject while no longer answering the human's current question. World already distinguishes evidence claims from authenticated proof authority. This proposal adds a domain-level relationship between a versioned requirement, its assumptions, and evidence applicability. It does not create a new proof verifier, global evidence grade or decision-packet schema.

Observed surfaces at World `fe1e411`: `world/types.ail` stores `Proposal.goal` and `plan` as strings; `ProofReportV1` in `host/evidence/types.go` binds subject/compiler and verified identities. `Validator.ValidateProof` checks authenticity, expected subject, compiler, outcome and completeness. These observations motivate an explicit requirement binding in a domain package; they do NOT establish the absence of all equivalent mechanisms across World.

## Goals and non-goals

Answer, deterministically and with reasons: “Does this evidence apply to the selected requirement revision under the selected assumptions and observation policy?” Preserve historical evidence and show the difference between historical validity and current applicability.

Non-goals: new `Evidence` variants or grades; modifying `DecisionPacket/v1`; rewriting history; resolving natural-language entailment automatically; treating applicability as approval; reimplementing proof authentication; a global dependency graph; invalidating evidence by mutating the original evidence object.

## Proposed domain objects

Proposed package namespace: `packages/world-requirements/` (publication name and canonical encoding freeze at design approval). All names here describe proposed data, not currently available AILANG syntax.

| Object | Content |
|---|---|
| Requirement revision | Immutable ID; stable logical requirement key; predecessor revision(s); owner identity reference; original request reference; human-readable criterion; check-spec reference; assumption-set reference |
| Evidence binding | Requirement revision; evidence reference; subject digest; exact dependency-manifest digest; toolchain/policy references; observation freshness rule; originating transition reference |
| Applicability request | Selected requirement revision and current identities, plus explicitly recorded logical evaluation time and policy |
| Applicability result | `applicable`, `stale`, `unavailable` or `unsupported`, with ordered reason codes and references |

Use World's existing content-addressed object envelope/canonicalization conventions. The binding is itself a claim until an authorized validator verifies its links and required evidence; an agent cannot assert an `applicable` or `PROVEN` flag into authority. A read-only report may display a claimed binding as claimed, without upgrading its evidence kind.

## Decision rules

1. Resolve all required references through bounded reads. Missing, malformed or unauthenticated authoritative input produces `unavailable`/`unsupported`, never an empty successful set.
2. Authenticate proof receipts through the existing host validator when a rule requires proof. Keep the domain applicability result separate from the proof's resolved grade. A library minted seal is not serializable production authority; restart revalidates from receipts under the production owner's rules.
3. Compare the binding's exact requirement revision, subject and dependency manifest, assumptions and policy with the current applicability request. Any known mismatch is `stale`, preserving the prior binding unchanged.
4. For observations, evaluate a declared freshness policy at the supplied recorded logical time. No ambient clock in the pure decision. “Captured at T” does not establish that the observation is true. Expiry is separate from authenticity.
5. Only fully resolved matching inputs whose required evidence meets the rule produce `applicable`. When a changed requirement happens to be logically equivalent, still require a new authorized binding; text similarity never transfers evidence.
6. Forked successor revisions are explicit competing choices. Do not pick “latest” by timestamp. The current requirement is selected by an authorized World transition referencing an exact revision; presentation cannot make that selection.

Initial dependency scope is deliberately conservative: a complete declared manifest for the bounded task, with exact hashes and a verified capture procedure. Unknown closure completeness prevents `applicable`. Do not claim a minimal causal slice. Changing one input may force revalidation of the whole task before finer-grained reuse is justified.

## Evidence presentation contract

Retain the ratified kind labels `PROVEN`, `TESTED`, `ATTESTED`, `CLAIMED`, subject to their existing minting rules. Add a companion view with separate fields: kind, verdict, proposition/scope, assumptions, applicability, observed/evaluated time and provenance.

Example: “Proof succeeded for requirement R1, source S1; STALE for selected R2: requirement changed.” A failed test remains TESTED with FAIL, not an upgrade or downgrade chosen by prose. An approved action is permission under scope, not proof of the outcome. `PROVEN` remains unavailable until the production proof-to-renderer path required by the existing proof design is satisfied.

This is a proposed additive refinement of HUMAN-SURFACE P3, not a unilateral rewrite of its ratified grammar. The parent must be reconciled at approval.

## Architecture and why this is a package

Pure applicability decisions belong in a domain package, with contracts written first per S1. A thin host adapter resolves object references and supplies authenticated evidence results; it cannot mint its own proof authority. Store old and new objects through existing World transitions. No new table or frozen kernel export is proposed.

The package returns data. Production execution authority stays with existing broker and proposal enforcement. This phase is a read-only applicability/reporting capability: `applicable` does not authorize dispatch, and no consumer may treat it as a reusable authorization token. The [vertical child](w-requirement-change-vertical.md) specifies the separate current-state execution guard.

Proposed adapter placement is a package-specific host module, only if the existing extension boundary cannot host it. Before sprint planning, read that extension boundary and document “why not a package” for every proposed host change. No claim is made that the generic package loader is already sufficient.

## Plan and file scope

| Milestone | Deliverable | Provisional scope |
|---|---|---|
| EA1 | Object/binding specification; pure decision laws; canonical fixture vectors | New domain package, ~200–350 LOC plus named tests/contracts |
| EA2 | Bounded resolver adapter and read-only applicability view | ~200–350 host/adapter LOC plus tests; no changes to proof minting |
| EA3 | Requirement revision, dependency edit and restart examples | Fixtures/runbook; required verification manifests updated through existing scripts |

Do not hand-edit generated `packages/world-core` mirrors. If a frozen core change becomes necessary, stop and split it into a separately justified design rather than expanding this one.

## Acceptance criteria and mutation controls

| ID | Acceptance | Named mutation that must fail |
|---|---|---|
| EA1 | R1 evidence matches R1 but is stale for R2; historical query still reports the R1 result | `EA-IGNORE-REQUIREMENT`: remove revision comparison |
| EA2 | Dependency-body edit with same public interface makes the binding stale | `EA-INTERFACE-ONLY`: replace source-manifest identity with interface identity |
| EA3 | Unresolved link, incomplete manifest and forged proof receipt never produce applicable/proven | `EA-MISSING-AS-EMPTY`, `EA-TRUST-CLAIM`: separately neuter resolution and proof validation |
| EA4 | Changed assumptions or expired observations produce named, deterministic reasons | `EA-IGNORE-ASSUMPTIONS`, `EA-IGNORE-EXPIRY`: each has its own control |
| EA5 | Serialization/restart does not resurrect an in-memory proof seal | `EA-DECODE-AS-AUTHORITY`: bypass revalidation |
| EA6 | Renderer shows kind, verdict and applicability independently | `EA-HIDE-STALE`, `EA-HIDE-FAIL`: stale proof and failed-test fixtures discriminate |
| EA7 | Read-only report dispatches zero effect handlers and does not mutate historical objects | `EA-REPORT-DISPATCH`: counting handler; compare history identities before/after |
| EA8 | Null input and budget exhaustion are loud; a real matching case passes | `EA-ZERO-SUCCESS`, `EA-UNBOUND-READ`: null and bounded-reader adversarial fixtures |

Controls are proposed implementation tests, not tests already present. Contract proof identities and named test floors must be installed through World's existing non-vacuity gates. Use the pinned released AILANG binary for any new `.ail`; this document contains no new `.ail` snippet or unverified syntax claim.

## High-impact decisions / design freeze

- [ ] Mark approves package-first scope and additive evidence presentation.
- [ ] Canonical identity and owner-selection protocol frozen; no implicit “latest.”
- [ ] Dependency-capture completeness and freshness policy specified for one domain.
- [ ] Production proof adapter ownership identified; unavailable remains the default until satisfied.
- [ ] Host integration inventory confirms no frozen-core expansion.

Designer may choose layout and internal algorithms within these invariants. New freshness policy members or wire fields require schema versioning, not reinterpretation of stored objects.

## Verification log, overlap and risks

Read `world/types.ail`, `host/evidence/{types,validator}.go`, `host/store/schema.sql`, HUMAN-SURFACE and the implemented proof/decision-lifecycle designs. These establish the positive reuse surfaces above. No language limitation or absence claim is used to justify a new kernel feature. Manual topic search found the proof authority, evidence-grade and frozen lifecycle parents; this child owns CURRENT applicability, not those mechanisms. Semantic search in the AILANG checkout had fallback embeddings, so no neural-score completeness claim is made for World.

Prior review ran fresh evidence/workbench tests and one broker recovery test; those are background controls, not acceptance of this unimplemented feature. Risks: incomplete dependency manifests, over-invalidation, misleading grade hierarchies and authenticity confused with relevance. Prefer conservative unavailability to fabricated applicability. The experiment can fail on utility without weakening authority.

## Axiom compliance

Directional design assessment, not implementation approval. No hard violation proposed on A1/A3/A4/A7.

| Axiom | Score | Design constraint |
|---|---:|---|
| A1 Determinism | +1 | Identity-bound inputs; deterministic mechanical results |
| A2 Replayability | +1 | Preserve inputs and outcomes needed to reproduce the decision |
| A3 Effect legibility | 0 | Existing effect semantics unchanged |
| A4 Explicit authority | +1 | Metadata and model output never mint execution authority |
| A5 Bounded verification | +1 | Explicit limits and checks with named refusal outcomes |
| A6 Safe concurrency | 0 | Sequential first version; snapshot identity checked |
| A7 Machines first | +1 | Structured, versioned artifacts with explicit unavailable states |
| A8 Minimal syntax | +1 | No new language syntax |
| A9 Cost visibility | +1 | Record tool/model costs and failure overhead |
| A10 Composability | +1 | Reuse existing compiler, evidence and protocol boundaries |
| A11 Structured failure | +1 | Unknown, incomplete and stale cannot masquerade as success |
| A12 System boundary | +1 | Separate claims, verification, permission and action |

**Net +10.** Re-score if implementation scope changes.
