# w-workbench-grade-projection — the workbench object page projects real grades from stored, subject-bound proof envelopes through a daemon-held validator

- **Status**: **implemented** (iteration 254: landed `9d16a07` via PR #248; judge sonnet PASS 88/100, zero blocking; CI PR 2/2 green) — design history: designed iteration 253 (orphan, `pi:ollama/glm-5.3:cloud`), quorum r1 BLOCKED → revision, r2 BLOCKED → narrow-refinement carve-out with the reviewers' verbatim fixes (§0.4, §0.5)
- **Queue row**: 105, `w-workbench-grade-projection` · clause-5 · **v1.0.0 release prerequisite** — ruled routable by attended `D-WORLD-69` = **A** (2026-10-09, exception path 26 → 105 → 172 → 150; the loop MAY take it when no critical-path row is routable). PARKED behind item 13's validated-boundary follow-on — that follow-on (row 26, `w-bounded-z3-report-producer`) landed 2026-10-12 as `0390aa5`, so the park is lifted. Filed by the rows-98+102 designer at iter-186 (charter line 1912, V33): "EVERY OBJECT PAGE NOW SHOWS ONE CONSTANT, TRUE GRADE REASON ('no canonical host projection: this workbench handler does not resolve subject-bound evidence into an object grade'); A REAL GRADE NEEDS STORED, DECODED EVIDENCE BOUND TO A SUBJECT AND A DAEMON-HELD VALIDATOR. Filed so the constant reads as a waypoint, not the destination."
- **Filed**: 2026-10-12, iteration 253 (designer role, unattended mission stage)
- **Measurement base**: code base `07687c2` for V1-V37; this revision (V2) re-verifies at the same code base — worktree HEAD `1945fb2` is `07687c2` plus this document's round-1 commit alone (`git diff 07687c2..1945fb2 --stat`: 1 file, the doc), so every new row V38-V42 observes identical code
- **Instrument**: `$HOME/.pinned-ailang/ailang`, AILANG **v0.41.0** (first `--version` line `AILANG v0.41.0`, V7/V41), sha256 `1a67b0146858450182f48082299956ea5b08cdf30131979b188b339fcedb5b9f` — the same binary the row-26 design pinned. No `.ail` file is touched by this design (V24), so no `ailang check` run is required by the fluency protocol; the projection is host-boundary Go.
- **Estimate**: **~1.5 d** (3 milestones; ~220 production LOC). The queue row carries no day figure (it was filed as a waypoint record); this design prices the wiring it forces plus the projection itself.
- **Scope**: `host/daemon` (validator wiring + the workbench grade projection + its tests) and `cmd/ailang-worldd` (two serve flags + usage). **Zero changes** to `host/evidence`, `host/store`, `host/workbench`, `scripts/`, any `.ail` file, `host/boundary`, `host/verifygate`, or `.github/workflows/ci.yml` (all by observation, §9/AC7).

---

## §0 The inherited iter-186 claims, re-measured at base (do not trust, verify — verified here)

Row 105 was filed with two inherited measurements taken at `6127ff3` (iter-186). Both are re-derived at `07687c2` before this design builds on them, as the controller's brief requires.

### 0.1 "0 non-test importers of `host/evidence`" — STILL TRUE, and stronger

At base, the only non-test Go files mentioning `host/evidence` are the package's own `validator.go` and `types.go` (V3). Since then row 26 landed the producer inside that package, and the stronger fact now holds: **0 non-test callers** of `NewValidator` and of `GenerateProof` anywhere in `host/` or `cmd/` (V28, with 14 firing test-call controls). No production composition root exists on either the minting or the validating side. The workbench handler still renders the constant unconditionally (V2).

### 0.2 "Validator needs a key, compiler config and required identities the daemon does not hold" — PARTIALLY STALE, corrected by measurement

Re-derived at base, `evidence.NewValidator(key [32]byte, reader ObjectReader, cfg CompilerConfig, requiredIdentities []string)` needs four things (V4), and the daemon's holdings have changed since iter-186:

| Need | Who holds it at base | Measured fact |
|---|---|---|
| `[32]byte` MAC key | **Nobody.** `daemon.Config` has no key field at all — `Key|Secret|Mac|Hmac` reads **0** in `host/daemon/daemon.go`, control `AilangBin` reads 10 (V5) | This design adds it (D6) |
| `Compiler hashref.HashRef` | **The daemon holds it** when `--ailang-bin` is set: `New` archives the interpreter and stores `d.interpreterHash` = the sha256 content-address of the executable's exact bytes (V6) — the same kind of ref the producer's `CheckerConfig.Compiler` records and the validator compares | Held; no new deliverable |
| `CompilerVersion string` | **Half-held.** `d.interpreterVersion` is the manifest's VERBATIM `--version` output (7 lines on the pinned binary, V7); the validator needs the exact line the producer records — the RAW first line with a trailing `\r` stripped (`proof_producer.go:328`), which differs from the daemon's `releaseFromVersion` spelling (TrimSpace, V7) | This design adds the producer-mirroring derivation (D6) |
| `requiredIdentities []string` | **Nobody.** No Config field exists (V5) | This design adds it (D6) |
| `ObjectReader` (store seam) | **The daemon holds the store**: the same `*store.Store` satisfies `evidence.ObjectReader` (`ReadObject` + `BusyTimeout()`) unchanged (V16, V29) | Held; no new deliverable |

The iter-186 filing predates none of these — it is simply narrower than today's truth: the daemon NOW holds the compiler identity (row 26's producer landed against the same archive discipline), but still holds no key and no identities, and `NewValidator`'s ordering pin (`ObjectReadTimeout > BusyTimeout`) is implied by the daemon's existing `checkReadOrdering` startup pin (V16). **The daemon still lacks two of the four needs, so the inherited claim's conclusion stands and the design must add them as explicit deliverables — it does (D6).**

### 0.3 The named successor that never filed — this row absorbs the minimum

Item 17's §9 decomposition names `w-proven-evidence-production-key-wiring` as tranche 3 ("provision/load durable MAC key material; retain the keyed validator before serving; abort startup closed"). That successor has **no queue row**: the token appears in the charter only inside item 17's landed record and the `D-WORLD-17` ruling (V34) — three mentions, none a numbered row. Row 105's own text demands "a daemon-held validator", and the controller's brief is explicit: what the projection needs and the daemon lacks must be an explicit deliverable, not an assumption. **This design therefore carries the minimal validator-side wiring (key file + identities + fail-closed construction) as in-scope deliverables; producer-side composition (who calls `GenerateProof` in production) stays OUT of scope and unowned** (V28: 0 non-test callers; Non-Goals).

### 0.4 Quorum round 1 → revision

Round 1 was BLOCKED (`gemini-3-1-pro` REJECT; `oc-kimi-k3` REJECT; `claude-sonnet-5@claude-p` ABSENT — quota, the anthropic daily ration, so the quorum degraded to N-1; controller PASS). Neither reject disputed the design direction; each named ONE narrow, fully specified surface, and the controller measured both first-party. This bounded revision re-verified each measurement at the same code base before folding it in (V38-V42 record command AND observed output; the controller's brief cited `daemon.go:687` for the assignment and `validator.go:207-208` for the resolver pair — both re-measured here as `:688` and `:208-209`, recorded corrections to the citations, not to the facts). **Quorum note:** a first attempt accidentally ran on the pinned v0.41.0 binary, whose model registry lacks `oc-kimi-k3` and `claude-sonnet-5@claude-p` (both read `unknown-model`) — recorded as a pin-version gap, not a degrade; the quorum re-ran on the PATH v0.54.5 binary. This revision applies exactly:

- **F1 (`oc-kimi-k3`) — the D6 CompilerVersion premise was asserted, not verified, and AC3 sourced the envelope's version from the daemon side (a self-consistent tautology).** The premise is now VERIFIED end-to-end, link by link: the Manifest field's own doc comment defines `Version` as "the verbatim `ailang --version` output captured at archival" (`archive.go:155`, V38); `probeVersion` returns `stdout.String()` VERBATIM — trailing whitespace preserved, stderr never joined, and `TrimSpace` reads **0** anywhere in `host/archive` with a firing same-file control (V39); the daemon assigns it UNCHANGED (`d.interpreterVersion = m.Version`, `daemon.go:688`, V40); and the producer's recorded value in the real row-26 tests is exactly the pinned binary's first `--version` line — all four `CompilerVersion` pins equal `AILANG v0.41.0`, which is what `$HOME/.pinned-ailang/ailang --version | head -1` prints (V41). **The real gap is FIXED, not just documented:** AC3 and the M2/M3 test specification no longer source the envelope's `CompilerVersion` from any daemon-side value — the test derives it INDEPENDENTLY (run the fake's `--version`, apply the producer's exact first-line rule to THAT raw stdout) and a NAMED cross-check test (`TestWorkbenchGradeCompilerVersionDerivationsAgree`, AC3/M11) reds loudly if the daemon-side and producer-side derivations of the same fake run ever diverge, with the side of each value stated in the AC. The hermetic cross-check is the load-bearing arm; the pinned-binary arm is an optional skipped-when-absent OBSERVATION.
- **F2 (`gemini-3-1-pro`) — D4's pseudocode called an API surface no V-row verified.** Added V42: `ValidationResult.Validated()` at `validator.go:110`, `.Unsupported()` at `:111`, `.Err()` at `:114`, `ValidateProof(ctx, reportRef, expectedSubject)` at `:120`, `ResolutionResult.Proven()` at `:208`, `.Err()` at `:209`, `Resolve(sealed)` at `:211` — every method D4's state machine calls exists with the called shape, and D4 now cites the row.

### 0.5 Quorum round 2 → narrow-refinement carve-out (iteration 254)

Round 2 (`2026-10-12T20-58-08Z`; author `pi:ollama/glm-5.3:cloud`, glm seat benched; seats `oc-kimi-k3`, `claude-sonnet-5@claude-p`, `gpt6-1-sol` ABSENT unreachable → reserve `gemini-3-1-pro`) was BLOCKED by three rejects, none disputing the design direction, each carrying a concrete `proposed_fix`. Two (`oc-kimi-k3`, `gemini-3-1-pro`) found that §0.4 cited V38-V42 while §10 ended at V37 — orphan iteration 253's revision was cut off before its log rows were written. `oc-kimi-k3` also found the exactly-100 truncation boundary defect; `claude-sonnet-5@claude-p` found the cross-check fixture unspecified. The iteration-254 controller applied the reviewers' fixes verbatim, re-running every command first-party at base `edb1b63`: V38-V42 (kimi's exact commands), V43 (sonnet's exact fixture bytes, plus a second discriminating fixture and a non-discriminating control), V44 (the page cap that admits kimi's 101-row fetch), the `+1` fetch/`truncated` rule in D4/D5, and the exactly-100 / 101 boundary arms in M3, AC4 and M8. Per surface: round 2's objections land on (a) verification-log completeness, (b) the truncation boundary, and (c) fixture specification; none on the projection direction.

## Axiom Compliance

| Axiom | Score | Justification |
|-------|-------|---------------|
| A1: Determinism | +1 | The grade is a pure function of (stored candidate bytes, daemon key, pinned interpreter identity, configured identities); candidates enumerate in the index's `hash_ref` ascending total order and the first validating candidate wins — the same store state always renders the same page (V13) |
| A2: Replayability | 0 | No replay machinery touched; the receipt→`CLAIMED` mapping and replay paths are unchanged |
| A3: Effect Legibility | +1 | Every effectful operation is host-boundary Go in the daemon; the AILANG surface is untouched (no `.ail` change, V24) |
| A4: Explicit Authority | +2 | `PROVEN` is minted on the page only through `ValidateProof` → `ValidatedEvidence` → `Resolve` → `Proven()` — the validator stays the sole mint, the projection holds no free resolver, the key is daemon-held and fail-closed, and the unconfigured state renders the honest constant (V4, V11) |
| A5: Bounded Verification | +1 | The projection performs at most one indexed page query and at most `gradeCandidateScanLimit` (100) bounded validator reads, each capped at `evidence.MaxBytes` (256 KiB) by the store's pre-materialization probe, all under the request's `readDeadline` (10 s); truncation is disclosed on the page, never hidden (V13-V16, D5) |
| A6: Safe Concurrency | 0 | No new shared production state; the validator is immutable after construction; scans are per-request |
| A7: Machines First | +1 | Every non-grade state is one stable sentence constant mapped from the validator's typed `UnsupportedReason` vocabulary; operational failures keep the existing 503/500 classification, never a fabricated reason (D4, V35) |
| A8: Minimal Syntax | 0 | No language surface at all |
| A9: Cost Visibility | +1 | The per-page cost is visible and bounded: the unconfigured page adds zero reads (AC1); the configured page adds one indexed metadata query plus at most 100 bounded point reads; today's production stores hold zero envelopes, so the steady state is one empty page (V28) |
| A10: Composability | +2 | Reuses the store's semantic-ID lookup index, the landed validator, its codecs and `Resolve`, the existing render `GradeView`, the existing `readStore` seam and `writeWorkbenchStoreError` classifier — zero changes to `host/evidence`, `host/store`, `host/workbench` (V13, V16, V25, V29) |
| A11: Structured Failure | +1 | Every refusal is a named reason with a named test; the mapping is total over the 12-token `UnsupportedReason` vocabulary plus an honest fallback; each reason's premise is test-constructed (AC4) |
| A12: System Boundary | +2 | The projection sits exactly where the workbench origin drew the boundary: the grade vocabulary stays AILANG-side, no Go transcription of `gradeOf` exists or is added, and the kernel keeps only the untrusted claim (V12, V25) |

**Net Score: +12** → **Decision: move forward.** Hard axioms A1/A3/A4/A7 all ≥ 0.

---

## Problem Statement

**Current state (measured at `07687c2`):**

- Every workbench object page renders one constant grade reason, unconditionally, at exactly one construction site (`workbench.go:477`); the constant is declared at `workbench.go:43` with the comment "Named reasons for what the object inspector cannot show" (V2). The constant is TRUE today — no daemon-held validator exists (V3, V28) — but it is also the ONLY state the page can express: an envelope in the store, a mismatched subject, an unauthenticated report, and an absent store are indistinguishable to the operator.
- Row 26 landed the bounded producer: real `ai-check` output can now be minted into MAC-tagged canonical envelope objects in the store, content-addressed, carrying `SemanticID = "world/proof-report/v1"` and `InterfaceHash = InterfaceHashV1` (V8). But nothing production-side can READ them back into a grade: `NewValidator` has 0 non-test callers (V28).
- The ratified grade vocabulary is four grades, `PROVEN | TESTED | ATTESTED | CLAIMED` (`world/types.ail:33-38`), and the render already supports all four labels plus a separate verdict channel (`host/workbench/render.go`, V25) — but no Go mechanism mints any of them: the only grade resolver in Go is the validator's `Resolve`, which mints exactly `ResolvedGradeProven` and nothing else (V11). The AILANG-side `gradeOf` maps `TestReport(_, _)` to `TESTED` regardless of its pass/fail boolean — a recorded rendering hazard (V12).
- The daemon holds the interpreter's content-address and archived manifest when `--ailang-bin` is set (V6) but no key material and no identity requirements (V5), so the validator the projection needs cannot be constructed from today's `Config`.

**Impact:** without this row, the v1.0.0 workbench cannot show a single real grade even for a store that holds real, producer-minted, subject-bound proof envelopes — the constant would be FALSE in the very world it was written to describe as a waypoint, and the honest states around a real grade (absent, foreign, unauthenticated, failed) stay unrenderable.

## Goals

**Primary goal:** the workbench object page projects `PROVEN` for an object exactly when the store holds a decoded, subject-bound, MAC-authenticated proof envelope that the daemon-held validator accepts — and otherwise renders one measured, per-object, TRUE reason; the constant becomes the true statement of the unconfigured daemon, not dead text.

**Success metrics:**
1. End-to-end (AC3): configured daemon + stored envelope minted through the exported codecs under the daemon's key → the subject's object page renders `<span>PROVEN</span>` exactly once, through `ValidateProof` → `Validated` → `Resolve` → `Proven()` — with the envelope's `CompilerVersion` derived PRODUCER-side (independently of the daemon) and the two version derivations cross-checked by a named test, so the green path cannot be a daemon-sourced tautology (round 1, `oc-kimi-k3`).
2. Every reachable state has its own reason and its own named test, and the reason is TRUE by the state the test constructs (AC4).
3. The unconfigured daemon renders the constant verbatim and performs zero evidence reads (AC1); the existing `TestWorkbenchObjectGrade` passes unchanged (V18).
4. No gate surface changes: no `host/evidence` test is added (the 64-member manifest stays), no `.ail` moves, no new subprocess site, no `host/boundary` file, no new exported declaration in `host/evidence` (AC7, §9).

## High-Impact Decisions

| Decision | Why High Impact | Chosen By | Deadline | Change Cost |
|----------|-----------------|-----------|----------|-------------|
| **PROVEN-only projection.** The page mints `PROVEN` and nothing else; `TESTED`/`ATTESTED`/`CLAIMED` are not projected because no daemon-held mechanism mints them (V11), and the projection deliberately does NOT transcribe `gradeOf` | The charter records the hazard: `gradeOf` grades `TestReport(_, false)` `TESTED`, badge-indistinguishable from a passing one (V12); the item-13 doc ratified that as AILANG-side semantics — "the boolean records outcome; grade records method" — and the workbench origin ruled "grade and verdict must be separate channels" and "the renderer never computes an evidence grade in Go" (V25). A Go transcription would create the exact grade-laundering surface both docs forbid | agent (measured, V11/V12/V25) | design | low |
| **The validator is a daemon field behind a tiny unexported interface seam (`proofValidator`), production-wired to `*evidence.Validator`** | It mirrors the repo's `readStore`/`durableStore` pattern (V16): production always wires the real thing, tests wrap. Without the seam, the Resolve-refusal branch and the reason-mapping arms for reasons a real store cannot produce (`invalid_ref`, `missing`, `hash_mismatch`, `wrong_semantic_id`) would be unkillable — "a branch nothing can pin is not a guard" is this mission's own measured lesson (`host/evidence/types.go:60-68` comment, V9) | agent (mirrors V16's pattern) | design | low |
| **Candidate discovery is one `ObjectsBySemanticID` page at `"world/proof-report/v1"`, bounded at 100, with truncation disclosed on the page** — no subject index is added | The reference index records entry/world fields, not payload subjects (row 103's own scope); a subject index is a store schema change this row must not make (row 26's no-new-store-surface discipline). The semantic-ID keyset page never reads payloads (V13), and the producer writes exactly that semantic ID (V8) | agent (measured, V8/V13) | design | low |
| **The unconfigured state performs ZERO evidence reads and keeps the constant verbatim** | The constant must remain TRUE and reachable (controller's brief); the unconfigured daemon's page must not grow a scan; `TestWorkbenchObjectGrade`'s existing assertions stay the contract (V2, V18) | agent | design | low |
| **Key custody is an operator-supplied exactly-32-byte key file (`--proof-key-file`), never auto-provisioned** | An auto-generated key would be a validator nobody's producer can satisfy — it would look configured while every envelope fails MAC, a false authority. The producer needs the SAME key on its side; the file is the minimal operator affordance that keeps the daemon and a future producer caller in agreement | agent | design | low |
| **The validator's `CompilerVersion` is derived from the archived manifest by MIRRORING the producer's derivation (raw first line, trailing `\r` stripped), not `releaseFromVersion`'s TrimSpace — and every TEST derives the envelope's version producer-side, never from the daemon** | Measured: the manifest carries the verbatim `--version` output (the field's own contract, V38; `probeVersion` returns stdout verbatim with `TrimSpace` reading 0 on the archive path, V39; the daemon assigns it unchanged, V40); the producer records `strings.TrimRight(strings.SplitN(out, "\n", 2)[0], "\r")` (V7), and the row-26 tests pin exactly the pinned binary's real first line (V41). A TrimSpace-side derivation would make every envelope `tool_mismatch` the day a pinned binary's first line carries trailing spaces — silent divergence between the two halves of the same discipline — so the named cross-check reds on a discriminating fixture (AC3/M11) (V43: exact bytes `"AILANG v0.41.0 \r\nbuild abc\n"`, on which the producer derivation yields `"AILANG v0.41.0 "` and a TrimSpace derivation yields `"AILANG v0.41.0"` — round 2, `claude-sonnet-5@claude-p`) and no envelope's version ever comes from the daemon side | agent (fixture-measured, V7, V38-V41) | design | low |
| **Reason strings are frozen sentence constants, asserted exactly by tests** | The constant today is already a frozen sentence with a byte-exact test (V2); new reasons join the same discipline — a wording change reds tests deliberately | agent | design | low |

### Design Freeze

All seven decisions above are made on measured evidence in this document; none requires human input before implementation.

- [x] PROVEN-only, no `gradeOf` transcription (frozen; V11/V12/V25)
- [x] The `proofValidator` seam, production-wired to the real validator (frozen; V16 pattern)
- [x] One semantic-ID page, scan bound 100, truncation disclosed (frozen; V8/V13)
- [x] Unconfigured ⇒ constant + zero reads (frozen; V2/V18)
- [x] Operator key file, exactly 32 bytes, never auto-provisioned (frozen)
- [x] Producer-mirroring version-line derivation, with the envelope's version derived producer-side and the two halves cross-checked by a named test (frozen; V7, V38-V41; AC3/M11 — round 1's tautology fix)
- [x] Frozen reason sentences (frozen)

## Solution Design

### Overview

One new production file in `host/daemon` (`workbench_grade.go`), small wiring in `daemon.go`, a one-line construction-site change in `workbench.go`, and two serve flags in `cmd/ailang-worldd`. The projection is a pure read path: it discovers candidate envelope objects through the store's existing semantic-ID index, hands each to the daemon-held validator against the page's object as the expected subject, and renders the first validating candidate as `PROVEN` through the validator's own `Resolve`. Every other outcome renders one measured reason sentence. Nothing in `host/evidence`, `host/store`, or `host/workbench` changes.

The four design questions this document was asked to answer map to: **D1** (where evidence lives, what binds it to a subject), **D2** (how the daemon obtains a validator, what happens when it cannot), **D3** (what grades an authenticable envelope projects, and the pass/fail hazard), **D4** (the honest per-object states and their reasons).

### D1 — Where stored, subject-bound evidence lives, and what binds it to a subject (question 1)

- **Table.** Envelope objects live in the store's `objects` table (`hash_ref, interface_hash_ref, semantic_id, provenance, payload` — the `INSERT` SQL the row-23 store uses, V15). The producer writes them with `Hash = SumSHA256(envelope)`, `InterfaceHash = InterfaceHashV1`, `SemanticID = "world/proof-report/v1"`, `Provenance = "bounded-ai-check-producer"` (V8) — so every envelope is discoverable by the store's existing `objects_by_semantic_id` index (row 103; keyset page over `(semantic_id, hash_ref)`, never reads payloads, `MaxSemanticIDPage = 500`, V13).
- **Decoding.** The payload is the canonical two-member authenticated envelope (`{"report": <base64url>, "mac": <base64url>}`); `DecodeAuthenticatedEnvelope` strictly decodes it and decodes the inner canonical nine-member report exactly once (V9). The validator performs this decode inside `ValidateProof`; the projection never re-decodes (no second decode — the package's own measured lesson, V9).
- **Subject binding.** The binding is the report's `subject` member: a canonical `HashRef` inside the MAC-covered canonical report bytes (`ProofReportV1.Subject`, `types.go:41-51`, the second member). `ValidateProof` computes `HMAC-SHA-256(key, envelope.Report)` and requires constant-time equality BEFORE comparing `report.Subject != expectedSubject` (V9, V10) — so an envelope is bound to a subject BY THE MAC: an attacker cannot rebind a report without the key. The caller supplies the subject per call: the projection passes the object page's ref as `expectedSubject`, which is what makes the grade the OBJECT's grade, not the envelope's.
- **Production reality.** At base, zero production-written envelopes exist: `GenerateProof` has 0 non-test callers (V28). The empty-store reason is therefore the expected live state after this row lands — a fact, not a defect, and the reason says exactly that.

### D2 — How the daemon obtains a validator, and what happens when it cannot (question 2)

**What `NewValidator` needs (V4):** a `[32]byte` key (copied to a fresh heap allocation), an `ObjectReader` (the store satisfies it: `ReadObject(ctx, ref, maxBytes)` + `BusyTimeout()`), a `CompilerConfig{Compiler hashref.HashRef, CompilerVersion string, ObjectReadTimeout time.Duration}` (compiler non-zero, version non-empty, timeout strictly positive and, when `BusyTimeout() > 0`, strictly above the busy window), and a non-empty `requiredIdentities` slice.

**Who holds each at daemon start, after this design (D6):** the key comes from the operator's `--proof-key-file` (exactly 32 bytes, read once at startup, never logged); the reader is the daemon's own store handle; `Compiler` is `d.interpreterHash` — held whenever `--ailang-bin` is set (V6); `CompilerVersion` is derived from the archived manifest by mirroring the producer's first-line derivation (V7); `requiredIdentities` come from repeatable `--proof-required-identity` flags. `ObjectReadTimeout` reuses the daemon's `readDeadline` (10 s), whose ordering against the store's busy window is ALREADY pinned at startup by `checkReadOrdering` (V16) — the daemon's existing pin implies the validator's, and `NewValidator` re-checks it anyway.

**When they are absent:** the daemon is *unconfigured* — `d.grader == nil` — and every object page renders the constant reason, verbatim, performing **zero** evidence reads (AC1). This is the state every daemon is in today (no key file exists as a concept, V5), so the constant remains TRUE and reachable, not dead text: it is the true statement of a daemon that holds no validator. When the flags are HALF-set, startup fails closed with named `StartupError`s at `StageConfig` (AC2): a key file without `--ailang-bin`; a key file without identities; a missing, non-regular, short, or long key file. The daemon never degrades to a half-configured validator.

**The required-identities semantics, stated honestly:** the daemon's identity list is daemon-wide. A subject whose envelope proves fewer identities than the daemon requires renders the proof-incomplete reason (its own authenticated report is the witness, V10). Per-subject identity configuration is deferred (Future Work) — this row does not invent a config surface the queue never asked for.

### D3 — What the page projects, and the pass/fail hazard (question 3)

- **An authenticable, subject-bound envelope projects exactly one grade: `PROVEN`.** The chain is `ValidateProof(ctx, candidateRef, pageRef)` → `ValidatedEvidence` → `Resolve(sealed).Proven()` → `workbench.GradeView{Available: true, Label: GradePROVEN}`. `Resolve` is the sole grade resolver and mints only `ResolvedGradeProven` (V11); the projection calls it and surfaces a refusal as an operational error (never a silent fallback) — the seam makes that branch killable (M5).
- **The mapping to the ratified vocabulary is: PROVEN from a validated seal; nothing else.** `TESTED`, `ATTESTED`, `CLAIMED` are not projected because no daemon-held mechanism can mint them: the only Go-side resolver resolves `PROVEN` (V11), and the AILANG-side `gradeOf` lives in `world/types.ail` with its own ratified semantics (V12). The projection adds no Go transcription of `gradeOf` — measured at base, no Go source names `gradeOf`/`EvidenceGrade`/`gradeCode` (V25, re-measured by the workbench origin's V5 at its own base and unchanged in kind here: the projection file adds none).
- **The pass/fail hazard, handled explicitly.** The charter records: "`gradeOf`'s mapping sends `TestReport(_, _)` to `TESTED` regardless of its pass/fail boolean, so a FAILING test report carries the same badge as a passing one" (V12). The item-13 grade doc ratified that as deliberate AILANG-side semantics — grade records method, verdict records outcome — and the workbench origin ruled the rendering consequence: "Grade and verdict must be separate channels" (V25). This design keeps the hazard out of the workbench by NOT keeping any such mapping: there is no `TESTED` badge to wear, so a failing test report has no badge at all. The render-side guard that would catch a future `TESTED` wiring already exists and is already mutation-drilled: `NewGradeView(GradeTESTED, KindTestReport, nil)` refuses with `ErrMissingVerdict` (V25), and row 34 drilled the grade/verdict template line (V25). The projection's `PROVEN` view carries no verdict — grade and verdict stay separate channels (M10 pins it).

### D4 — The projection state machine and the honest reason strings (question 4)

`objectGradeView(ctx, ref) (workbench.GradeView, error)` in the new file, called from the object-page construction site (workbench.go:477's replacement):

```text
if d.grader == nil:
    return NewGradeUnavailable(objectGradeUnavailableReason)          # the constant, verbatim, zero reads
candidates, err := d.reads.ObjectsBySemanticID(ctx, evidence.ProofSemanticID, "", gradeCandidateScanLimit+1)   # one extra row: proves truncation (V44)
if err != nil: return GradeView{}, err                                # operational: 503/500, never a reason
truncated := len(candidates) == gradeCandidateScanLimit+1             # only a REAL 101st row truncates (round 2, oc-kimi-k3)
if truncated: candidates = candidates[:gradeCandidateScanLimit]      # iterate only the first 100
if len(candidates) == 0:
    return NewGradeUnavailable(gradeReasonNoStoredEvidence)
firstFailure := ""
for cand in candidates (hash_ref ascending, the index's total order):
    res := d.grader.ValidateProof(ctx, cand.Hash, ref)                # ref is the expected SUBJECT
    if res.Err() != nil: return GradeView{}, wrap(res.Err())          # operational
    if sealed, ok := res.Validated(); ok:
        if _, proven := d.grader.Resolve(sealed).Proven(); !proven:
            return GradeView{}, errors.New("...refused its own seal")  # operational; seam-killable (M5)
        return GradeView{Available: true, Label: workbench.GradePROVEN}
    if reason, unsupported := res.Unsupported(); unsupported && reason != UnsupportedSubjectMismatch:
        if firstFailure == "": firstFailure = reason                   # first non-mismatch failure, deterministic
if firstFailure != "":
    return NewGradeUnavailable(reasonSentence(firstFailure) + qualifier(truncated))
return NewGradeUnavailable(gradeReasonBoundElsewhere + qualifier(truncated))
```

**The API surface the pseudocode calls is verified, not assumed (round 1, `gemini-3-1-pro`):** `func (r ValidationResult) Validated() (ValidatedEvidence, bool)` (`validator.go:110`), `.Unsupported() (UnsupportedReason, bool)` (`:111`), `.Err() error` (`:114`), `func (v *Validator) ValidateProof(ctx context.Context, reportRef, expectedSubject hashref.HashRef) ValidationResult` (`:120`), `func (r ResolutionResult) Proven() (ResolvedGrade, bool)` (`:208`), `.Err() error` (`:209`), and `func (v *Validator) Resolve(sealed ValidatedEvidence) ResolutionResult` (`:211`) — every method and both constructors the state machine calls exist with the called shape (V42).

**Why this classification is TRUE per state.** `ValidateProof`'s check order (V10) is what makes each sentence honest: hash → semantic ID → interface → decode → MAC → **subject** → tool → proof-failed → incomplete. Everything after the MAC check carries exact subject knowledge (`tool_mismatch`, `proof_failed`, `proof_incomplete` are BY CONSTRUCTION subject-bound and authenticated when the validator reports them); everything before it does not, so those sentences never claim binding. `subject_mismatch` proves a MAC-valid envelope bound elsewhere. The first-non-mismatch rule is deterministic (candidates enumerate in `hash_ref` ascending order, V13) and every sentence states only what the deciding candidate's result established.

**The reason vocabulary (frozen sentences; the mapping is total over the validator's 12 tokens, V10):**

| Validator result(s) | Frozen reason sentence |
|---|---|
| (no grader) | `no canonical host projection: this workbench handler does not resolve subject-bound evidence into an object grade` — the constant, byte-identical to today's (V2) |
| 0 candidates | `no stored evidence: the store holds no object with semantic ID world/proof-report/v1` |
| all candidates `subject_mismatch` | `stored proof envelopes exist; none is bound to this subject` |
| `tool_mismatch` | `subject-bound evidence authenticates under the daemon's key but names a compiler other than the daemon's pinned interpreter` |
| `proof_failed` | `subject-bound evidence authenticates under the daemon's key but records a failed proof outcome` |
| `proof_incomplete` | `subject-bound evidence authenticates under the daemon's key but does not verify every identity the daemon requires` |
| `unauthenticated_report` | `a stored proof envelope does not authenticate under the daemon's key` — deliberately claims no subject binding: the MAC is checked before the subject (V10), so the failing envelope's subject is not established |
| any of `invalid_ref, missing, oversize, hash_mismatch, wrong_semantic_id, wrong_interface, malformed` | `a stored candidate object does not satisfy the proof envelope's identity, content or schema checks` — one folded sentence for the seven structural classes (rationale below) |
| any future/unknown reason token | `stored proof envelopes exist and none validates for this object` — the invariant of the branch, true for any mix |
| truncated scan (any non-PROVEN terminal state) | the sentence plus ` (scan bounded to the first 100 stored candidates)` |

**Why the seven structural classes fold into one sentence.** Three of them cannot be produced through this path at all in a sole-writer store: the scan itself parses every candidate ref (`invalid_ref`) and predicates the semantic ID (`wrong_semantic_id`, V13/V15), and both insert paths verify content addresses (`hash_mismatch`, V15). `missing` requires the row to vanish between the index query and the read inside one request — impossible under the daemon's sole writer, reachable only through the seam. Per-class sentences for classes nothing real can produce would be four unpinnable arms — the mission's "a branch nothing can pin is not a guard" lesson (V9). The fold keeps every ARM pinnable: `malformed` and `wrong_interface` are store-reachable (a mislabelled object with the proof semantic ID), so the folded sentence is killed by real-store arms, and the seam pins the remaining tokens through the mapping-table test (AC4). `oversize` is likewise store-reachable (a >256 KiB candidate), giving a third real arm.

**Operational failures never become reasons.** The scan error, the validator's read error, and a Resolve refusal return `(GradeView{}, err)`; the handler routes them through the existing `writeWorkbenchStoreError` — timeout → 503, `ReferenceIndexUnavailableError` → 503, else sanitized 500 (V35). A grade reason is a fact about evidence, never a report about infrastructure. One declared residue: `LookupIndexUnavailableError` from the scan answers the workbench's sanitized 500 while the `/v1` semantic-ID route classifies the same typed error as 503 (V35) — the shared classifier is not widened by this row; the residue is declared here, not absorbed.

### D5 — The scan bound and its cost

`gradeCandidateScanLimit = 100` (a named constant in the new file; the magnitude mirrors `WorkbenchPageLimit`). One indexed metadata query (payloads never read, V13) plus at most 100 validator validations, each one bounded point read (`ReadObject`'s pre-materialization length probe refuses anything over `evidence.MaxBytes` = 256 KiB, V14) plus one HMAC and one strict decode, all under the request's `readDeadline` (10 s, V16). Worst case per object page: 100 × ≤256 KiB = ≤25 MiB of bounded reads; today's production steady state is one empty page (V28). The query fetches `gradeCandidateScanLimit+1` = 101 rows (within `MaxSemanticIDPage` = 500, V44) and iterates only the first 100; the scan is declared truncated ONLY when the 101st row exists, and then every non-PROVEN terminal sentence carries the truncation qualifier. Exactly 100 stored candidates is an exhaustive scan and carries NO qualifier — the page neither claims exhaustiveness it did not perform nor claims a truncation it did not perform (round 2, `oc-kimi-k3`, verbatim fix).

### D6 — The wiring (Config, startup, CLI)

**`daemon.Config` gains two fields** (additive; no existing Config consumer breaks — no reflection over Config exists, V27):

```go
ProofKeyFile           string   // --proof-key-file: a 32-byte MAC key file; empty = unconfigured (the constant state)
ProofRequiredIdentities []string // --proof-required-identity (repeatable): identities every validated proof must carry
```

**`Daemon` gains one field:** `grader proofValidator` (nil when unconfigured), where `proofValidator` is the small unexported seam (§High-Impact Decisions). **`New` wires it** immediately after the `--ailang-bin` block (which has already set `d.interpreterHash`/`d.interpreterVersion`, V6):

```go
if cfg.ProofKeyFile != "" {
    if d.interpreterHash.IsZero() {                       // --ailang-bin unset: no compiler identity
        return nil, d.abort(StageConfig, "the workbench grade validator needs the pinned interpreter; --ailang-bin is unset", nil)
    }
    key, err := loadProofKey(cfg.ProofKeyFile)            // exactly 32 bytes, regular file, bounded read (33-byte limit reader)
    if err != nil { return nil, d.abort(StageConfig, "cannot load the proof key file", err) }
    v, err := evidence.NewValidator(key, s, evidence.CompilerConfig{
        Compiler:          d.interpreterHash,              // the archived bytes' content address (V6)
        CompilerVersion:   proofCompilerVersionLine(d.interpreterVersion), // producer-mirroring first line (V7)
        ObjectReadTimeout: readDeadline,                   // already pinned above the busy window (V16)
    }, cfg.ProofRequiredIdentities)
    if err != nil { return nil, d.abort(StageConfig, "cannot construct the workbench grade validator", err) }
    d.grader = v
}
```

`loadProofKey` opens the path, stats the OPEN file (no path TOCTOU), refuses non-regular files, reads at most 33 bytes, and requires exactly 32 — a growing or shrunken key file is refused by name, and the read is bounded regardless of races. `proofCompilerVersionLine` is `strings.TrimRight(strings.SplitN(verbatim, "\n", 2)[0], "\r")` — the producer's own spelling (`proof_producer.go:328`, V7), deliberately NOT `releaseFromVersion`'s TrimSpace. **The premise this derivation rests on — that `d.interpreterVersion` carries the VERBATIM `--version` stdout — is verified link-by-link, not asserted (round 1, `oc-kimi-k3`):** the `Manifest` field's own doc comment defines `Version` as "the verbatim `ailang --version` output captured at archival" (`archive.go:155`, V38); `probeVersion` returns `stdout.String()` verbatim — trailing whitespace preserved, stderr never joined — and `TrimSpace` reads **0** anywhere in `host/archive` with a firing same-file control (V39); `New` assigns it into the daemon unchanged (`d.interpreterVersion = m.Version`, `daemon.go:688`, V40). The discipline's other half — that the PRODUCER records this same line — is pinned by the row-26 tests themselves: all four `CompilerVersion` pins in `producer_realstore_test.go` equal the pinned binary's real first `--version` line (V41). The two halves are held together by a NAMED cross-check test, not by trust: every test derives the envelope's `CompilerVersion` producer-side (independently of the daemon), and `TestWorkbenchGradeCompilerVersionDerivationsAgree` reds if the daemon-side and producer-side derivations of the same fake run ever diverge (AC3, M11). The key bytes are never logged (the existing `writeInternalErrorLog` discipline; V35).

**The CLI gains two serve flags** (`--proof-key-file PATH`, repeatable `--proof-required-identity NAME` — the `run-allow-caps` flag pattern, V27) plus their lines in the `usage` constant's three spots (header comment, usage block, flag docs, V27). The help tests are `Contains`-based and resolve text from the same constants, so they stay green (V27).

### D7 — What each failure does (refusal, never partial trust)

| Failure | Observable | Page state after |
|---|---|---|
| Unconfigured (`grader == nil`) | the constant, verbatim | unchanged page, zero added reads |
| Half-set flags at startup | named `StartupError` at `StageConfig`; daemon never starts | no serving |
| Key file missing / non-regular / ≠32 bytes | `cannot load the proof key file: …` naming the arm | no serving |
| Scan operational error (timeout / lookup index unavailable / store error) | 503 Timeout / sanitized 500 through `writeWorkbenchStoreError` (V35) | no fabricated reason |
| Candidate read error (`res.Err() != nil`) | wrapped operational error → sanitized 500 | no fabricated reason |
| Resolve refuses its own seal | operational error → sanitized 500 (seam-killable, M5) | never a PROVEN badge without the resolver's word |
| 0 candidates | the empty-store reason | honest |
| All candidates foreign | the bound-elsewhere reason (+ qualifier if truncated) | honest |
| First non-mismatch failure | that class's frozen sentence (+ qualifier if truncated) | honest |
| Envelope validates | `<span>PROVEN</span>` exactly once | the real grade |

No path renders a grade without the validator's seal, and no path renders a reason whose premise the store state does not establish.

### Implementation Plan (milestones)

- **M1 — wiring + CLI (~0.4 d; ~70 prod LOC).** Config fields, `loadProofKey`, `proofCompilerVersionLine`, the `New` wiring block, the `proofValidator` seam, the two CLI flags + usage lines. Tests: `TestProofKeyWiringRefusesEveryHalfConfiguration` (six arms: missing file, directory path, 31 bytes, 33 bytes, key-without-`--ailang-bin`, key-without-identities), `TestServeProofFlags` (usage contains both flags; `--proof-key-file` with a missing interpreter reaches the named stage errors). Exit gate: `go build ./... && go test ./host/daemon -run 'TestProofKeyWiring' -count=1 && go test ./cmd/ailang-worldd -count=1` green.
- **M2 — the projection (~0.4 d; ~90 prod LOC).** `workbench_grade.go`: the scan-limit constant, the reason sentences + map + fallback, `objectGradeView`, the construction-site change in `workbench.go`. Tests: `TestWorkbenchGradeUnconfiguredServesTheConstantWithoutScanning` (zero-calls wrapper + the constant), `TestWorkbenchGradeProvenFromSubjectBoundEnvelope` (real validator, real store, envelope minted through the exported codecs under the daemon's key with `Subject` = the page object, `Compiler = SumSHA256(fake interpreter bytes)`, `CompilerVersion` = the PRODUCER-SIDE derivation — the test runs the fake's `--version` itself and applies the producer's first-line rule (`proof_producer.go:328`, V7) to that RAW stdout, never a daemon-side value (AC3) — `Verified` ⊇ required identities, on AC3's discriminating fixture (V43: exact bytes `"AILANG v0.41.0 \r\nbuild abc\n"`, on which the producer derivation yields `"AILANG v0.41.0 "` and a TrimSpace derivation yields `"AILANG v0.41.0"` — round 2, `claude-sonnet-5@claude-p`) — the full `ValidateProof → Validated → Resolve → Proven` chain). Exit gate: both green; `TestWorkbenchObjectGrade` STILL green and `git diff --stat -- host/daemon/workbench_test.go` empty.
- **M3 — the state table + gates (~0.7 d; ~60 prod LOC + the remaining test LOC).** All reason-state tests (each constructs its premise and asserts the exact sentence), the 12-token + unknown-token + fallback mapping-table test (seam-driven), the truncation boundary test (exactly 100 foreign envelopes renders the bound-elsewhere sentence WITHOUT the qualifier; 101 renders it WITH — round 2, `oc-kimi-k3`), the three sanitized-operational tests, the verdict/label sweep, the CLI flag test, the version cross-check `TestWorkbenchGradeCompilerVersionDerivationsAgree` (AC3: the daemon-side `proofCompilerVersionLine(d.interpreterVersion)` vs the producer-side derivation of the test's own `--version` run, on the discriminating fixture (V43: exact bytes `"AILANG v0.41.0 \r\nbuild abc\n"`, on which the producer derivation yields `"AILANG v0.41.0 "` and a TrimSpace derivation yields `"AILANG v0.41.0"` — round 2, `claude-sonnet-5@claude-p`); the pinned-binary arm an optional skip-when-absent OBSERVATION), and the full mutation re-drill M1-M11. Exit gates: `go test ./host/daemon ./cmd/ailang-worldd -count=1` green; `./scripts/verify_go.sh` green (UNINFORMATIVE UNDER SANDBOX here — the controller re-runs; the socket-bind finding V30); `AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_ail.sh` green (no `.ail` change, V24).

### Files to Modify/Create

**New files:**
- `host/daemon/workbench_grade.go` — the projection: the scan-limit constant, the frozen reason sentences and the `UnsupportedReason → sentence` map with fallback, the `proofValidator` seam, `objectGradeView`, `loadProofKey`, `proofCompilerVersionLine` (**~150 prod LOC**)
- `host/daemon/workbench_grade_test.go` — the named tests of AC1-AC6, including the version cross-check `TestWorkbenchGradeCompilerVersionDerivationsAgree` (the producer-side derivation vs `proofCompilerVersionLine` over the daemon's archived verbatim `d.interpreterVersion`, both sides named; V38-V41) (~650 test LOC; the envelope minting mirrors `authority_test.go`'s `envelopeFor` recipe, V32; the configured daemon uses the `fakeInterpreter` helper, V31)

**Modified files:**
- `host/daemon/daemon.go` — the two Config fields + comments, the `grader` field, the `New` wiring block (**~45 prod LOC**)
- `host/daemon/workbench.go` — the construction site: `Grade: workbench.NewGradeUnavailable(objectGradeUnavailableReason)` becomes `grade, err := d.objectGradeView(ctx, ref)` + the error branch + `Grade: grade`; the constant's declaration is untouched (**~5 prod LOC**)
- `cmd/ailang-worldd/main.go` — two flags + three `usage` spots (**~22 prod LOC**)
- `cmd/ailang-worldd/main_test.go` — `TestServeProofFlags` (**~45 test LOC**)

**Unchanged (with receipts):** `host/evidence/**` (no new exported declaration — `TestPublicAuthoritySurfaceIsFrozen`'s exact manifest untouched, V23; no new test — the 64-member `REQUIRED_EVIDENCE_TESTS` and the `:329` count literal stay, V20), `host/store/**` (consumed: `ObjectsBySemanticID`, `ReadObject`, `PutObject`'s insert identities — no store file changes, V13-V15), `host/workbench/render.go` + `live.js` (the render already supports PROVEN and the verdict channel; live.js swaps no inspector region, V25/V26), `scripts/verify_go.sh` + `scripts/verify_ail.sh` (V20/V24), every `.ail` (V24), `host/boundary` (V22), `host/verifygate` (V20), `.github/workflows/ci.yml` (both gates already run in CI).

## Examples

### Example 1 — the configured daemon renders a real grade (AC3's demonstration)

**Before (today):** every object page renders `GRADE UNAVAILABLE — no canonical host projection: this workbench handler does not resolve subject-bound evidence into an object grade` (workbench.go:477, V2) — for a store that holds a producer-minted envelope exactly as for an empty one.

**After (illustrative):** with `ailang-worldd serve --db world.db --ailang-bin <pinned> --proof-key-file key.mac --proof-required-identity applyRevision`, and a store holding the subject module object plus the producer's envelope object (`SemanticID world/proof-report/v1`, `Subject` = the module's ref, `Verified ⊇ [applyRevision]`), the module's object page renders `<p><span>PROVEN</span></p>` — minted by the validator's seal, resolvable by nothing else.

### Example 2 — the honest states around a real grade

| Store state (all under a configured daemon) | Page renders |
|---|---|
| no envelope object exists | `GRADE UNAVAILABLE — no stored evidence: the store holds no object with semantic ID world/proof-report/v1` |
| one envelope, valid MAC, bound to another subject | `GRADE UNAVAILABLE — stored proof envelopes exist; none is bound to this subject` |
| one envelope, `Subject` = this object, MAC under a foreign key | `GRADE UNAVAILABLE — a stored proof envelope does not authenticate under the daemon's key` |
| one envelope, authenticated, `Compiler` ≠ the daemon's interpreter | `GRADE UNAVAILABLE — subject-bound evidence authenticates under the daemon's key but names a compiler other than the daemon's pinned interpreter` |
| 150 foreign envelopes | the bound-elsewhere sentence + ` (scan bounded to the first 100 stored candidates)` |

## Success Criteria

See §7 (AC1-AC7) and §8 (M1-M10). Summary gates: both pinned verify gates green at HEAD; the mutation matrix killed with named tests; `TestWorkbenchObjectGrade` green and its file untouched (observation, not a seeded write); the evidence manifest still 64, untouched.

## Testing Strategy

**Unit/state tests (`workbench_grade_test.go`):** every reason state constructs its PREMISE first (the store state and, where relevant, a direct validator call proving the classification) and then asserts the exact sentence — a reason is never asserted without its premise. The mapping-table test drives `objectGradeView` through a fake `proofValidator` for all 12 `UnsupportedReason` tokens plus one unknown token plus the `Validated`-and-`Resolve`-refused arms, asserting each sentence exactly and that no state leaks another's text. The unconfigured test wraps `d.reads` (the `refFailingStore` embed pattern, V18) and counts `ObjectsBySemanticID` calls.

**Integration (real store + real validator, no fakes on the golden path):** AC3's end-to-end uses `newHandlerDaemon` + `fakeInterpreter` (V31) + a real 32-byte key file + a real envelope minted through `EncodeProofReportV1` + HMAC-SHA-256 + `EncodeAuthenticatedEnvelope` (the `envelopeFor` recipe, V32) + `d.store.PutObject` with the producer's identities (V8) — no hand-authored bytes: every byte is produced by the exported codecs.

**Manual/CI:** no manual steps; CI's existing `go test ./... -count=1` covers both changed packages; the verify gates are unchanged.

**Binary independence:** the projection tests use the fake interpreter, never the 100 MB pinned binary — the daemon tests stay hermetic (V31), and no test in this design joins the binary-dependent `host/evidence` set.

## Deferred Decisions

- The key file's permission bits (e.g., refusing group/world-readable keys) — deferred: platform-dependent behavior (Windows), no in-repo precedent; named as a risk, owned by the unfiled production-custody successor (§0.3).
- Per-subject required identities — deferred (D2's honest narrowness; Future Work).
- Surfacing WHICH envelope validated on the PROVEN page (a link from the grade line to the envelope object's page) — the render has no slot; Future Work.
- The `LookupIndexUnavailableError` 500-vs-503 classification residue (D4) — deferred to a workbench error-table owner; declared, not absorbed.

## Non-Goals

- **No `host/evidence` change** — no new exported declaration, no new test, no codec or validator edit; the frozen public surface and the 64-member manifest stay (V20, V23).
- **No store change** — no subject index, no new read surface; the scan consumes the existing keyset page (V13).
- **No render change** — `render.go` and the drilled grade/verdict line are untouched (V25).
- **No producer composition root** — who calls `GenerateProof` in production stays unowned (V28); this row wires the CONSUMER only.
- **No full key-custody regime** — auto-provisioning, rotation, KMS: the unfiled tranche-3 successor's scope; this row carries only the minimal load-and-hold the projection needs.
- **No TESTED/ATTESTED/CLAIMED projection and no `gradeOf` transcription** — the hazard defense is the absence of the mapping (D3).
- **No `/v1` JSON grade surface** — the queue row is the workbench object page; the `/v1` routes are untouched.
- **No `.ail` change, no subprocess, no gate-script edit** (V21, V24).

## Timeline

| Milestone | Estimate | Production LOC |
|---|---:|---:|
| M1 wiring + CLI flags + fail-closed tests | 0.4 d | ~70 |
| M2 the projection + the constant state + the PROVEN end-to-end | 0.4 d | ~90 |
| M3 the state table, truncation, sanitized paths, mutation re-drill, gates | 0.7 d | ~60 |
| **Total** | **~1.5 d** | **~220** |

The queue row carries no day figure (a waypoint record); this estimate prices the forced wiring (M1), the projection (M2), and the state-table test mass (M3) honestly.

## Risks & Mitigations

| Risk | Impact | Mitigation |
|---|---|---|
| A store with >100 envelopes truncates the scan | Low | The qualifier discloses the bound on the page; raising it is a one-constant change; the steady state today is zero envelopes (V28) |
| The daemon-wide identity list renders coarse `proof_incomplete` reasons for subjects proving fewer identities | Med | The reason is TRUE (the report is its own authenticated witness); per-subject config is named Future Work, not silently assumed |
| Reason-wording churn reds tests | Low | Deliberate: the sentences are frozen like the constant (its test already byte-asserts, V2) |
| A world-readable key file is accepted | Med | Declared risk + deferred decision (§Deferred); the key never appears in logs; a leaked key's worst case is a forged PROVEN badge on a loopback, operator-owned workbench — not authority elsewhere |
| The version-line derivation diverges from the producer's on a future pin | Med | The derivation mirrors the producer's exact spelling (V7) and M3's tool-mismatch test fails loudly if the two halves drift; a pin move regenerates the producer fixtures in the same commit (row 26's discipline) |
| The seam lets a test fake the validator and mask a production wiring bug | Low | The golden path (AC3) uses a REAL validator on a real store; fakes cover only mapping arms no real store can produce; production always wires `*evidence.Validator` in `New` (D6) |
| Sandbox: gates unrunnable in this worktree (socket bind refused, V30) | Low | The controller re-runs both verify gates outside the sandbox; the in-process tests that run here are green (V18) |

## Related Documents

- `design_docs/implemented/w-bounded-z3-report-producer.md` — row 26, the direct relative: the producer whose envelopes this row projects; the store seams (`ObjectReader`, `PutObject`, `BusyTimeout`), the deadline discipline (`requireDeadline`, `ErrNoDeadline`), and the fixture/provenance conventions this design composes with unchanged
- `design_docs/implemented/w-validated-proven-evidence-boundary.md` — item 17, the parent: the validator, the seal, the `D-WORLD-17` ruling; §9's decomposition whose unfiled tranche-3 successor this row absorbs minimally (§0.3)
- `design_docs/implemented/w-workbench-object-provenance-and-grade.md` — rows 98+102, the iter-186 designer that filed this row's facts: §2b the constant's rationale, V10/V12/V20 the inherited measurements, §9 P-c this row's filing
- `design_docs/implemented/w-workbench-read-only.md` — the workbench origin: §2.6 the grade decision ("the renderer never computes an evidence grade in Go"; "grade and verdict must be separate channels")
- `design_docs/implemented/w-evidence-grade-mapping.md` — item 13: the ratified `gradeOf` semantics this design refuses to transcribe ("the boolean records outcome; grade records method")
- `design_docs/implemented/w-workbench-unpinned-render-hunks.md` — row 34: the mutation-drilled grade/verdict template line this design consumes without touching
- `design_docs/implemented/w-store-object-reference-index.md` — row 103: the `objects_by_semantic_id` index and the keyset page the scan consumes
- `ailang docs search` is refused in this deployment (AILANG_AGENT_POLICY unset — recorded in the session log); all related documents are named by path above

## References

- Charter queue row 105: `design_docs/world-mission.md:1912` (verbatim scope, the waypoint framing, the iter-186 filing, V33)
- Rulings: `D-WORLD-69` = A (the exception path 26 → 105 → 172 → 150); `D-WORLD-17` = A (the seal-binding ruling whose tranche-3 residue §0.3 absorbs)
- `design_docs/coding-standards.md` — effects at the host boundary; the AILANG fluency protocol (no `.ail` written by this design, V24)

## Future Work

- Production producer composition: who calls `GenerateProof` with the daemon's key — the unfiled producer-side successor (V28's 0 callers is the receipt).
- Full key custody (provisioning, rotation, permission hardening) — the unfiled tranche-3 scope beyond the minimal load-and-hold here.
- Per-subject required-identity configuration (a subject→identities map) — deferred with the coarse-but-true reason as the honest interim.
- A link from the PROVEN grade line to the validating envelope's object page (a render slot — row 34's drilled line would need a deliberate extension).
- TESTED/ATTESTED/CLAIMED projection, each behind its own validated mechanism that binds a verdict (for TESTED) — the render's `ErrMissingVerdict` guard is the waiting half of that work (V25).

---

## §7 Acceptance criteria

Fresh row: no inherited ACs. Each criterion carries its killing mutation (§8) and each test constructs its premise before asserting its observable.

**AC1 — the unconfigured daemon serves the constant, scans nothing, and the existing grade test is untouched.** A daemon constructed without `--proof-key-file` (the state every daemon is in today, V5) renders, for every object page, the byte-identical constant reason (V2) and performs ZERO `ObjectsBySemanticID` calls (asserted through a counting wrapper on `d.reads`, the `refFailingStore` embed pattern, V18). `host/daemon/workbench_test.go` is not modified by this row — established by observation (`git diff --stat` empty for the file; `TestWorkbenchObjectGrade` green at base, V18, and green after). The constant's declaration and text are byte-identical before and after (observed: `sed -n '40,44p' host/daemon/workbench.go` on both sides). Named test: `TestWorkbenchGradeUnconfiguredServesTheConstantWithoutScanning`. M1 and M9 kill it.

**AC2 — validator wiring is fail-closed with named refusals.** `--proof-key-file` constructs the validator only when the interpreter identity exists (`--ailang-bin` set) and identities are named; every half-configuration refuses at `StageConfig` with its own observable detail: key file absent; key path a directory (non-regular); key file 31 bytes; 33 bytes; key without `--ailang-bin` ("the workbench grade validator needs the pinned interpreter; --ailang-bin is unset"); key without identities (the validator's own refusal surfaced through the daemon's stage). `loadProofKey` stats the OPEN file and reads at most 33 bytes (bounded regardless of races); the key is never logged. The wiring's positive observable is delegated to AC3 — the only way to see the validator's pinned identity is through what it accepts and refuses (AC4's tool-mismatch arm is the negative side). Named tests: `TestProofKeyWiringRefusesEveryHalfConfiguration` (six arms), `TestServeProofFlags` (CLI: usage contains both flags; the flags reach `Config` and their refusals fire in order). M2 and M3 kill it.

**AC3 — PROVEN end-to-end from stored, decoded, subject-bound evidence through the daemon-held validator.** A configured daemon (real key file, `fakeInterpreter` for `--ailang-bin`, V31) whose store holds a subject object and one envelope object minted through the EXPORTED codecs under the daemon's key (`Subject` = the page object's ref, `Compiler = SumSHA256(fake interpreter bytes)`, `CompilerVersion` = the fake's first version line, `Verified ⊇` the required identities, clean proof outcome, stored via `PutObject` with the producer's `SemanticID`/`InterfaceHash` identities, V8/V32) renders, on the SUBJECT's object page, `<p><span>PROVEN</span></p>` exactly once — through `ValidateProof` → `Validated` → `Resolve` → `Proven()`, with no `GRADE UNAVAILABLE` line and no verdict span. Named test: `TestWorkbenchGradeProvenFromSubjectBoundEnvelope`. M4 (subject binding), M5 (Resolve bypass), M6 (label substitution), M10 (verdict channel) kill it.

**AC4 — every measured state renders its own TRUE reason; the mapping is total.** Each named test constructs the premise and asserts the exact sentence (D4's table): `TestWorkbenchGradeReasonNoStoredEvidence` (a direct store query returns 0 candidates — the premise), `TestWorkbenchGradeReasonBoundToOtherSubjects` (a MAC-valid envelope for subject B; page for A; premise proven by a direct `ValidateProof` call returning `subject_mismatch`), `TestWorkbenchGradeReasonUnauthenticatedEnvelope` (subject-bound envelope under a foreign key), `TestWorkbenchGradeReasonToolMismatch`, `TestWorkbenchGradeReasonProofFailed` (`CheckPassed=false`), `TestWorkbenchGradeReasonProofIncomplete` (identity missing), `TestWorkbenchGradeReasonStructuralCandidate` (two real-store arms: a proof-semantic-ID object with the wrong interface hash; one with a non-decodable payload), `TestWorkbenchGradeReasonScanTruncationIsDisclosed` (boundary arms, round 2 `oc-kimi-k3`: exactly 100 foreign envelopes → the bound-elsewhere sentence WITHOUT the qualifier; 101 foreign envelopes → the bound-elsewhere sentence carries the truncation qualifier), and `TestWorkbenchGradeReasonMapsEveryUnsupportedReason` — the seam-driven table over all 12 `UnsupportedReason` tokens (V10's count) plus one unknown token (the invariant fallback sentence) plus the `Validated`-with-refused-`Resolve` arm, asserting each sentence exactly and that no state leaks another's. M7 (a map-entry swap) and M8 (the dropped qualifier) kill it.

**AC5 — operational failures never masquerade as grade reasons.** The scan's operational error (wrapper returning a store error) answers the workbench's existing 503/500 classification (timeout → 503 "Timeout"; `ReferenceIndexUnavailableError` → 503; else sanitized 500, V35) — never a 200 with a reason; a candidate read error (`res.Err()`) and a Resolve refusal likewise return the sanitized operational path. Named tests: `TestWorkbenchGradeScanFailureIsSanitized`, `TestWorkbenchGradeValidationReadErrorIsSanitized`, `TestWorkbenchGradeValidatorRefusalIsSanitized`. M5 and M9 kill it.

**AC6 — the TESTED pass/fail hazard stays out: PROVEN-only, verdict channel untouched.** The projection contains exactly one grade-label literal (`GradePROVEN`); no state renders `TESTED`, `ATTESTED`, or `CLAIMED`, and the PROVEN view carries no verdict; a sweep over every reason state (driving the AC4 table) asserts no `<span>` label appears in any unavailable state and no `verdict-` span appears in the PROVEN state. The defense for the recorded hazard (V12) is the absence of any mapping that could wear it — stated in D3 and pinned here. Named test: `TestWorkbenchGradeNeverRendersVerdictOrUnmintedLabels`. M6 and M10 kill it.

**AC7 — the conflict surface is exactly the six files, and every enumerated gate is green or explicitly untouched.** Established by observation at merge time, never by a seeded write: `git diff --stat` names exactly `host/daemon/{workbench_grade.go,workbench_grade_test.go,daemon.go,workbench.go}` and `cmd/ailang-worldd/{main.go,main_test.go}`; `git diff --quiet` is empty for `host/evidence`, `host/store`, `host/workbench`, `scripts/`, `world/`, `host/boundary`, `host/verifygate`, `.github/`; `go test ./host/evidence -count=1` stays 64/64 with `AILANG_BIN` set (the manifest gate's own scope, V20); `./scripts/verify_go.sh` and `AILANG_BIN=… ./scripts/verify_ail.sh` both green (controller-run outside the sandbox; the socket-bind finding V30 recorded as UNINFORMATIVE UNDER SANDBOX).

## §8 Mutation table

Every mutant is a compiling, behaviour-changing edit; a build failure is not a kill. Every observable is the named test's assertion, never a sibling value.

| ID | Exact file and neutered edit | Named check that fires | Downstream observable and predicted failure text |
|---|---|---|---|
| **M1 unconfigured-scan** | `host/daemon/workbench_grade.go`: hoist the `ObjectsBySemanticID` call above the `d.grader == nil` return (the scan runs unconditionally) | `TestWorkbenchGradeUnconfiguredServesTheConstantWithoutScanning` (AC1) | the counting wrapper observes 1 scan call on an unconfigured daemon; failure: `ObjectsBySemanticID calls = 1, want 0` — the unconfigured page must stay O(1) |
| **M2 key-size acceptance** | `host/daemon/workbench_grade.go`: `loadProofKey`'s `len(raw) != 32` → `len(raw) < 32` (and the `<=` twin) | `TestProofKeyWiringRefusesEveryHalfConfiguration` (AC2, the 31- and 33-byte arms) | a 33-byte file loads (or a 31-byte file loads); failure: `New accepted a %d-byte key file, want the 32-byte refusal` |
| **M3 unset-interpreter acceptance** | `host/daemon/daemon.go`: `if d.interpreterHash.IsZero()` → `if false && d.interpreterHash.IsZero()` | `TestProofKeyWiringRefusesEveryHalfConfiguration` (AC2, the key-without-`--ailang-bin` arm) | the refusal detail changes to the validator's `compiler identity is unset`; failure: `startup detail = "cannot construct the workbench grade validator: …", want "the workbench grade validator needs the pinned interpreter; --ailang-bin is unset"` |
| **M4 subject binding** | `host/daemon/workbench_grade.go`: `d.grader.ValidateProof(ctx, cand.Hash, ref)` → `d.grader.ValidateProof(ctx, cand.Hash, cand.Hash)` | `TestWorkbenchGradeProvenFromSubjectBoundEnvelope` (AC3) | every candidate reports `subject_mismatch` against its own ref; failure: `body contains the bound-elsewhere reason; want <span>PROVEN</span>` — the grade must be the PAGE OBJECT's, not the envelope's |
| **M5 Resolve bypass** | `host/daemon/workbench_grade.go`: render PROVEN on `res.Validated()` alone, dropping the `Resolve(sealed).Proven()` gate | `TestWorkbenchGradeValidatorRefusalIsSanitized` (AC5) and the AC4 table's refused-`Resolve` arm | a seal the resolver refuses still wears PROVEN; failure: `status = 200 with <span>PROVEN</span>; want the sanitized operational failure` — the resolver stays the sole mint |
| **M6 label substitution** | `host/daemon/workbench_grade.go`: `workbench.GradePROVEN` → `workbench.GradeTESTED` | `TestWorkbenchGradeProvenFromSubjectBoundEnvelope` + `TestWorkbenchGradeNeverRendersVerdictOrUnmintedLabels` (AC3/AC6) | failure: `body lacks <span>PROVEN</span> / contains <span>TESTED</span>` — the hazard badge can never appear |
| **M7 reason-map swap** | `host/daemon/workbench_grade.go`: the map's `unauthenticated_report` entry → the bound-elsewhere sentence | `TestWorkbenchGradeReasonUnauthenticatedEnvelope` + `TestWorkbenchGradeReasonMapsEveryUnsupportedReason` (AC4) | failure: `reason = "stored proof envelopes exist; none is bound to this subject"; want "a stored proof envelope does not authenticate under the daemon's key"` — a MAC failure never claims subject knowledge |
| **M8 truncation qualifier dropped** | `host/daemon/workbench_grade.go`: the `qualifier(truncated)` append removed | `TestWorkbenchGradeReasonScanTruncationIsDisclosed` (AC4) | failure: `reason lacks " (scan bounded to the first 100 stored candidates)"` — a bounded scan never claims exhaustiveness. Its boundary arm (exactly 100 foreign envelopes → no qualifier) also kills the mutant `truncated := len(candidates) >= gradeCandidateScanLimit` (round 2) |
| **M9 operational→reason laundering** | `host/daemon/workbench_grade.go`: the scan-error return becomes `NewGradeUnavailable(gradeReasonNoStoredEvidence)` | `TestWorkbenchGradeScanFailureIsSanitized` + AC1's zero-calls arm | failure: `status = 200 with the empty-store reason; want 503/500` — infrastructure facts never become evidence facts |
| **M10 verdict channel** | `host/daemon/workbench_grade.go`: the PROVEN view gains `HasVerdict: true, Verdict: VerdictPass` | `TestWorkbenchGradeProvenFromSubjectBoundEnvelope` + `TestWorkbenchGradeNeverRendersVerdictOrUnmintedLabels` (AC3/AC6) | the grade line gains `✓ verdict: PASS`; failure: `grade line != "<p><span>PROVEN</span></p>"` — grade and verdict stay separate channels |

The first-match determinism of the candidate loop is structural (the loop breaks at the first `Validated`; continuing past it changes only reads, no observable), so no mutation row is claimed for it — stated, not hidden.

## §9 Conflict Surface

What already lives in the positions this design extends, and what must still pass:

| Neighbor / surface | Collision and ordering |
|---|---|
| **The evidence package's authority invariants** (`host/evidence`, V3/V4/V9-V11/V23) | **No collision — the design's central constraint.** No `host/evidence` file changes: no new exported declaration (`TestPublicAuthoritySurfaceIsFrozen`'s exact `frozenPublicSurface` manifest stays green unedited, V23), no new test (the 64-member `REQUIRED_EVIDENCE_TESTS` and the `verify_go.sh:329` count literal `64` stay — the manifest's scope is `go test -json ./host/evidence`, V20, and this design adds no `host/evidence` test), no validator or codec edit. The projection mints nothing: `PROVEN` still requires the validator's seal and `Resolve`'s word |
| **`scripts/verify_go.sh`** (V20) | **Untouched.** The named-manifest gate is scoped to `host/evidence` (its call site builds evidence JSON from `go test -json ./host/evidence`, V20); the new tests live in `host/daemon` and `cmd/ailang-worldd`, covered by the script's global `go test ./... -count=1`. No count literal moves |
| **The evidence-manifest self-mutation gate** (`host/verifygate/evidence_manifest_gate_test.go`, V20) | Untouched by construction: it executes the live script against synthetic sets; an additive `host/daemon` test file changes nothing it observes |
| **The module manifest gate** (`host/verifygate/module_manifest_gate_test.go`, V24) | **Untouched** — no `.ail` file is created, moved, or edited (V24: `find` enumerates the same set before and after); the gate's module inventory is unchanged |
| **The AC10 subprocess-site census** (`host/broker/registry_publish_test.go:1135`, V21) | **No new site.** The census enumerates every production `exec.Command*` under `host/`+`cmd/` and requires a driver per site-file; this design adds zero launch sites (the projection is pure reads, HMAC, and strict JSON decode). The `host/verifygate` subprocess-sink gate is likewise untouched (no `cmd.Stderr` sink added) |
| **`host/boundary`'s file-count pin** (`allowlist_world_test.go:1349`, V22) | **No new file there.** `wantFileCount = 1` stays; the new file lives in `host/daemon`, never in `host/boundary` |
| **The pinned-binary verify gate** (`scripts/verify_ail.sh`, V24) | **Untouched** — no `.ail` change; the gate's `ROOTS` sweeps (`design_docs`, `world/`, `packages/se-tools`) see an identical set |
| **`host/store`** (rows 22/23/103, V13-V16) | **No file changes.** `ObjectsBySemanticID` (the row-103 keyset page), `ReadObject` (the probe-bounded read), `PutObject`'s verified inserts — all consumed as-is; `requireDeadline`/`ErrNoDeadline` is the projection's defense-in-depth, and the request's `readDeadline` derivation satisfies it by construction (V16) |
| **`host/workbench` render + live** (V25/V26) | **No file changes.** The render already supports `PROVEN` and the verdict channel; row 34's mutation-drilled grade/verdict template line is consumed untouched; `live.js` swaps no inspector region (its `REGIONS` list excludes it, V26), so the live path needs nothing |
| **`cmd/ailang-worldd`** (V27) | Two additive flags + usage lines; the help tests are `Contains`-based over the same single-source constants; `runDaemon` stays the test seam; no verb, route, or guard changes |
| **The daemon's route table and error classifier** (V35) | **No route added.** The projection lives inside the existing `/workbench` object branch; its operational failures route through the UNCHANGED `writeWorkbenchStoreError` (the `LookupIndexUnavailableError` 500-vs-503 residue is declared in D4, not fixed here) |
| **The queue's ordering** | Row 105's park is lifted (row 26 landed `0390aa5`); rows 172/150 follow on the `D-WORLD-69` path — nothing here touches their scope; the unfiled key-wiring successor's scope is absorbed only to the minimum D6 states |

## §10 Verification Log

All commands were run from this worktree root at `07687c2` on 2026-10-12 with `/usr/bin/grep` (this rig's default `grep` is ugrep; recorded commands must be reproducible). Negative-existence rows carry a KNOWN-POSITIVE CONTROL in the same call — a pattern that fires in the same scope, proving the instrument works. No verification artifacts are written (the controller's one-deliverable rule); every measurement is reproduced by the command shown.

| ID | Claim | Exact command and same-call control | Observed output |
|---|---|---|---|
| V1 | Measurement base; the worktree carries only this document | `git rev-parse --short HEAD; git status --short` | `07687c2`; empty status before this document was written |
| V2 | The constant exists, is served at exactly one construction site, and its test asserts it byte-exactly | `/usr/bin/grep -rn "no canonical host projection" host/`; `/usr/bin/grep -rn "objectGradeUnavailableReason" host/ cmd/` | first call: `host/daemon/workbench.go:43` (the declaration) + three `host/workbench/render_test.go` lines (arbitrary reason TEXT in render tests, not the constant); second: `workbench.go:43`, `workbench.go:477` (`Grade: workbench.NewGradeUnavailable(objectGradeUnavailableReason)` — the single construction site), `workbench_test.go:988` (the byte-exact `want` line) |
| V3 | No non-test Go file outside `host/evidence` imports it | `/usr/bin/grep -rln "host/evidence" --include="*.go" host/ cmd/ \| grep -v _test.go`; control: the same grep WITHOUT the `grep -v _test` filter | check: `host/evidence/validator.go host/evidence/types.go` (in-package only); control fires with 9 files incl. `host/verifygate/evidence_manifest_gate_test.go` — the instrument sees package mentions |
| V4 | `NewValidator`'s four needs and its ordering pin | `sed -n '63,93p' host/evidence/validator.go` | `NewValidator(key [32]byte, reader ObjectReader, cfg CompilerConfig, requiredIdentities []string)`: refuses nil reader, non-positive `ObjectReadTimeout`, zero compiler, empty version, empty identities (each naming its own field); `window < 0` refused, `window > 0 && ObjectReadTimeout <= window` refused with `ErrUnorderedTimeouts`; key copied to a fresh heap allocation |
| V5 | The daemon's Config holds no key material | `printf 'key: '; /usr/bin/grep -c "Key\|Secret\|Mac\|Hmac" host/daemon/daemon.go; printf 'control: '; /usr/bin/grep -c "AilangBin" host/daemon/daemon.go` | `key: 0`; `control: 10` (the Config's interpreter fields fire — the instrument reads the file) |
| V6 | The daemon holds the interpreter's content-address and manifest version when `--ailang-bin` is set | `sed -n '676,689p' host/daemon/daemon.go` | the `if cfg.AilangBin != ""` block: `ref, err := a.Archive(cfg.AilangBin)`; `d.arch, d.interpreterHash = a, ref` (:687); `d.interpreterVersion = m.Version` (:688) — `Archive` content-addresses the executable's exact bytes (sha256, `archive.go:277-283`) |
| V7 | The version-line divergence this design must mirror | `$HOME/.pinned-ailang/ailang --version \| head -1` + `\| wc -l`; `/usr/bin/grep -n "TrimRight(strings.SplitN" host/evidence/proof_producer.go`; `sed -n '848,859p' host/daemon/daemon.go` | first line `AILANG v0.41.0`; the output is 7 lines (the manifest carries all of it verbatim); the producer's derivation is `proof_producer.go:328: line := strings.TrimRight(strings.SplitN(string(out), "\n", 2)[0], "\r")` (raw first line, no TrimSpace); `releaseFromVersion` (:852-859) uses `strings.TrimSpace` — a DIFFERENT spelling |
| V8 | The producer's envelope objects carry the exact identities the scan relies on | `sed -n '495,497p' host/evidence/proof_producer.go` | `p.writer.PutObject(writeCtx, store.Object{Hash: ref, InterfaceHash: InterfaceHashV1, SemanticID: ProofSemanticID, Provenance: "bounded-ai-check-producer", Payload: env, …})` |
| V9 | The subject binding is the MAC-covered `subject` member; the codec decodes exactly once | `sed -n '41,51p' host/evidence/types.go`; `sed -n '156,164p' host/evidence/validator.go`; `sed -n '61,63p' host/evidence/types.go` | `ProofReportV1` struct (:41-51) with `Subject hashref.HashRef` at :43 (the decoded report's second member); validator: MAC over `envelope.Report` with `hmac.Equal` (:156-159) BEFORE `if report.Subject != expectedSubject` (:162-163); the `decoded`-field comment (:61) carries the no-second-decode lesson |
| V10 | `ValidateProof`'s check order (what makes each reason sentence TRUE) and the 12-token vocabulary | `sed -n '124,184p' host/evidence/validator.go`; `/usr/bin/grep -n "Unsupported.*UnsupportedReason = " host/evidence/validator.go` | order: parse ref → read (oversize/missing) → hash → semantic ID → interface → decode → MAC → subject → tool → proof-failed → incomplete → validated; 12 `Unsupported*` tokens (`invalid_ref, missing, oversize, hash_mismatch, wrong_semantic_id, wrong_interface, malformed, unauthenticated_report, subject_mismatch, tool_mismatch, proof_failed, proof_incomplete`) |
| V11 | `Resolve` is the sole grade resolver and mints only PROVEN | `/usr/bin/grep -n "ResolvedGradeProven\|func (v \*Validator) Resolve" host/evidence/validator.go` | `:200 const ResolvedGradeProven ResolvedGrade = 1` (the ONLY grade constant); `:211` `Resolve(sealed)`; `:221` returns `ResolutionResult{grade: ResolvedGradeProven, ok: true}` only for its own seal, else `ErrUnmintedAuthority`/`ErrForeignSeal` |
| V12 | The ratified vocabulary and the pass/fail hazard | `sed -n '31,59p' world/types.ail`; `/usr/bin/grep -n "regardless of its pass" design_docs/world-mission.md` | `world/types.ail:33-38`: `EvidenceGrade = PROVEN \| TESTED \| ATTESTED \| CLAIMED`; `gradeOf`'s contract and body both map `TestReport(_, _) => TESTED` — the boolean is never read; charter line 3863: "gradeOf's mapping sends TestReport(_, _) to TESTED regardless of its pass/fail boolean, so a FAILING test report carries the same badge as a passing one" |
| V13 | Candidate discovery: the semantic-ID keyset page never reads payloads and is bounded | `sed -n '9,64p' host/store/objects_by_semantic_id.go` | `MaxSemanticIDPage = 500`; the SQL selects `hash_ref, interface_hash_ref, semantic_id, provenance` with `WHERE semantic_id = ? AND hash_ref > ? ORDER BY hash_ref LIMIT ?` ("never reads payload"; `hash_ref` is the PRIMARY KEY — a unique total order, so the scan is deterministic); `lookupIndexAvailable` gate returns typed `LookupIndexUnavailableError` |
| V14 | Bounded vs unbounded object reads (the projection must not use `GetObject` on candidates) | `/usr/bin/grep -n "func (s \*Store) GetObject" -A 4 host/store/store.go`; `sed -n '40,83p' host/store/read_object.go` | `GetObject` scans the FULL payload with no length probe (only `requireDeadline`); `ReadObject` probes `length(payload)` first and refuses `*ObjectTooLargeError` before materialization — the validator reads candidates through `ReadObject(ctx, ref, MaxBytes)` |
| V15 | Both insert paths verify content addresses, so `hash_mismatch` is a corruption sentinel | `/usr/bin/grep -n "func verifyObject\|verifyObject(o)" host/store/*.go \| /usr/bin/grep -v _test` | `store.go:514` `verifyObject` ("computed digest must equal o.Hash"); called from `PutObject` (:557) and `Commit` (:1113) — a stored object's payload hashes to its ref unless the DB was tampered with externally |
| V16 | The daemon's read seam, deadline, and ordering pin (what the projection runs under) | `sed -n '484,496p' host/daemon/daemon.go`; `/usr/bin/grep -n "readDeadline = \|checkReadOrdering(s.BusyTimeout" host/daemon/daemon.go`; `/usr/bin/grep -n "func (d \*Daemon) readCtx" -A 2 host/daemon/handlers.go` | `readStore` includes `ObjectsBySemanticID` (:492) — the scan joins the existing seam; `readDeadline = 10 * time.Second` (:144); `checkReadOrdering` refuses `window >= deadline` at startup (:643) — implying `NewValidator`'s pin for `ObjectReadTimeout = readDeadline`; `readCtx` = `context.WithTimeout(r.Context(), d.readDeadline)` (handlers.go:287) |
| V17 | The workbench object handler's construction site and error classifier | `/usr/bin/grep -n "Grade: workbench.NewGradeUnavailable" host/daemon/workbench.go`; `sed -n '113,124p' host/daemon/workbench.go` | `:477` the single `Grade:` construction inside the `?object=` branch; `writeWorkbenchStoreError`: timeout → 503 "Timeout", `*store.ReferenceIndexUnavailableError` → 503, else sanitized 500 |
| V18 | The existing grade test passes at base, unconfigured, with a proof-labelled object present | `go test ./host/daemon -count=1 -run 'TestWorkbenchObjectGrade' -v 2>&1 \| tail -3`; the test's own body: `sed -n '974,1010p' host/daemon/workbench_test.go` | `--- PASS: TestWorkbenchObjectGrade (0.01s)` with all six subtests PASS; the test (`:974`) builds a `world/proof-report/v1`-labelled object with the CORRECT interface hash and asserts the constant line once and NO label spans — the unconfigured state renders the constant for proof-labelled objects too, which AC1 preserves |
| V19 | The unconfigured test daemon holds no key | `sed -n '21,33p' host/daemon/handlers_test.go` | `newHandlerDaemon` builds `Config{DBPath: …, BindHost: …}` — no key, no interpreter: the unconfigured state is what every existing daemon test exercises |
| V20 | The evidence-manifest gate's scope, size, and call site (why this design leaves it untouched) | `python3 -c "import re; raw=open('scripts/verify_go.sh').read(); m=re.search(r'REQUIRED_EVIDENCE_TESTS = \{(.*?)\}', raw, re.S); print(len(re.findall(r'\"(Test\w+)\"', m.group(1))))"`; `/usr/bin/grep -n "check_evidence_manifest \"\$evidence_json\"\|exact top-level\|go test -json ./host/evidence" scripts/verify_go.sh` | 64 members; `:326 go test -json ./host/evidence -count=1` (the gate's scope is the evidence package ALONE); `:322` banner "64 exact top-level tests"; `:329 check_evidence_manifest "$evidence_json" 64` — a design adding no `host/evidence` test changes none of it |
| V21 | The AC10 census enumerates every production launch site (this design adds none) | `/usr/bin/grep -n "func enumerateSubprocessSites\|TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential\|has no driver" host/broker/registry_publish_test.go \| head -3`; `/usr/bin/grep -rn "exec\.Command" host/ cmd/ --include="*.go" \| /usr/bin/grep -v _test \| /usr/bin/grep -v verifygate \| wc -l` | the census at `:1077`/`:1135` walks `host/`+`cmd/` non-test Go and fails on any site-file without a driver (`:1219`); 12 textual `exec.Command` matches in 10 files (11 launch sites + 1 comment in `capsulecache.go` naming its absence) — the projection adds zero |
| V22 | `host/boundary`'s file-count pin | `sed -n '1349,1352p' host/boundary/allowlist_world_test.go` | `const wantFileCount = 1` with `t.Fatalf` on any other count — the new file must NOT live there (it lives in `host/daemon`) |
| V23 | The frozen public-authority surface is an exact manifest | `sed -n '270,289p' host/evidence/authority_test.go`; `/usr/bin/grep -n "func TestPublicAuthoritySurfaceIsFrozen" host/evidence/authority_test.go` | `frozenPublicSurface` is the EXACT inventory ("a MANIFEST, not a pattern: any new exported symbol … reds this test"); `:427` the test — this design adds no `host/evidence` declaration, so the manifest stays green unedited |
| V24 | No `.ail` file is touched, and the AIL gate's sweeps are unchanged | `find world host design_docs -name "*.ail" \| wc -l; find world -name "*.ail" \| wc -l` | 20 total under the swept scopes+`host/` (4 in `world/`, 7 row-26 fixtures + 1 replay fixture under `host/*/testdata/`, 8 sketches); this design creates/edits none — `verify_ail.sh`'s `ROOTS` (`design_docs`, `./world`, `packages/se-tools/se_tools`) see an identical set |
| V25 | The render already supports the projection's outputs, and the two-channel discipline is already drilled | `sed -n '14,33p' host/workbench/render.go`; `sed -n '252,273p' host/workbench/render.go`; `/usr/bin/grep -n "NewGradeView(GradeTESTED, KindTestReport, nil)" -A 2 host/workbench/render_test.go`; `sed -n '71,90p' design_docs/implemented/w-workbench-read-only.md` | `GradeView{Available, Label, HasVerdict, Verdict, Unavailable}` with all four labels; `NewGradeUnavailable(reason)`, `NewGradeView` refuses invalid labels and refuses `KindTestReport` without a verdict (`ErrMissingVerdict`); render_test drills the TESTED-without-verdict refusal; the workbench origin: "the renderer never computes an evidence grade in Go… Grade and verdict must be separate channels" |
| V26 | The live path never re-renders the inspector (no live-surface change needed) | `/usr/bin/grep -n "REGIONS =" host/workbench/live.js` | `REGIONS = ["nav…","section[aria-label=\"timeline\"]","section[aria-label=\"live\"]","section[aria-label=\"world graph\"]","section[aria-label=\"decisions\"]"]` — the inspector section (where the grade renders) is not swapped |
| V27 | The CLI's flag/usage surfaces this design extends | `sed -n '193,218p' cmd/ailang-worldd/main_test.go`; `sed -n '45,62p' cmd/ailang-worldd/help_test.go`; `/usr/bin/grep -n "fs.Func(\"run-allow-caps\"" cmd/ailang-worldd/main.go`; `/usr/bin/grep -n "reflect.TypeOf(Config\|NumField" host/daemon/*_test.go cmd/ailang-worldd/*_test.go` | `TestServeWorkspaceFlags` observes serve flags through failure messages (the pattern AC2's CLI test mirrors); the help tests are `Contains`-based over the single-source `usage` constant; `run-allow-caps` is the repeatable-func flag pattern; NO test reflects over `Config` fields — additive fields are safe |
| V28 | No production caller mints or validates envelopes today (the empty state is the live state) | `printf 'check: '; /usr/bin/grep -rn "GenerateProof(" --include="*.go" host/ cmd/ \| /usr/bin/grep -v _test \| /usr/bin/grep -v "func (p \*Producer)" \| wc -l; printf 'control: '; /usr/bin/grep -rn "GenerateProof(" --include="*_test.go" host/evidence/ \| wc -l; printf 'NewValidator-non-test: '; /usr/bin/grep -rn "NewValidator(" --include="*.go" host/ cmd/ \| /usr/bin/grep -v _test \| /usr/bin/grep -v "func NewValidator" \| wc -l` | `check: 0`; `control: 14` (the producer's own tests fire — the instrument sees call sites); `NewValidator-non-test: 0` — both halves of the evidence chain are library-only at base |
| V29 | The store satisfies `evidence.ObjectReader` unchanged | `sed -n '50,61p' host/evidence/validator.go`; `/usr/bin/grep -n "func (s \*Store) ReadObject\|func (s \*Store) BusyTimeout" host/store/read_object.go host/store/store.go` | `ObjectReader` = `ReadObject(ctx, ref, maxBytes) (ObjectMeta, []byte, error)` + `BusyTimeout() time.Duration`; `*store.Store` provides both — the daemon's own handle wires the validator with no adapter |
| V30 | SANDBOX FINDING (recorded, not worked around): the socket-binding daemon tests cannot run in this worktree | `go test ./host/daemon -count=1 -run 'TestWorkbenchObjectGrade\|TestWorkbenchLiveCursorIsResumable' 2>&1 \| tail -6` | `TestWorkbenchObjectGrade` PASS in-process; `TestWorkbenchLiveCursorIsResumable` panics `httptest: failed to listen on a port: listen tcp6 [::1]:0: bind: operation not permitted` — **UNINFORMATIVE UNDER SANDBOX** (a fence refusal, not a base defect); the controller re-runs the full gates outside the sandbox |
| V31 | The configured-daemon test shape exists (no 100 MB binary needed) | `sed -n '42,52p' host/daemon/daemon_test.go` | `fakeInterpreter(t, dir, version)` writes a `#!/bin/sh` script answering `--version` and returns its exact bytes — the same helper `TestHealthAndHeadRoundTrip` uses to exercise the archive hermetically |
| V32 | The envelope-minting recipe for daemon tests is the exported codecs (no hand-authored bytes) | `sed -n '48,66p' host/evidence/authority_test.go` | `envelopeFor`: `EncodeProofReportV1(report)` → `HMAC-SHA-256(key, raw)` → `EncodeAuthenticatedEnvelope{Report, MAC, MACValid}` — every byte produced by the exported codecs; the daemon test mirrors this against the daemon's key |
| V33 | The queue row's text and filing record | `awk 'NR==1912' design_docs/world-mission.md` | row 105 verbatim: the D-WORLD-69 = A exception path, the PARKED-behind-item-13 marker, the constant quoted, "A REAL GRADE NEEDS STORED, DECODED EVIDENCE BOUND TO A SUBJECT AND A DAEMON-HELD VALIDATOR", "filed so the constant reads as a waypoint, not the destination", the iter-186 attribution |
| V34 | The key-wiring successor is named but has NO queue row (why this design carries the minimum) | `/usr/bin/grep -n "w-proven-evidence-production-key-wiring" design_docs/world-mission.md \| cut -c1-80` | three hits, all inside item 17's landed record (`:4061`, `:4078`) and the `D-WORLD-17` ruling (`:901`) — a named successor, never filed as a numbered row |
| V35 | The operational-failure classification the projection routes through (and the declared residue) | `sed -n '523,575p' host/daemon/handlers.go \| /usr/bin/grep -n "LookupIndexUnavailable\|503\|StatusServiceUnavailable"`; `sed -n '113,124p' host/daemon/workbench.go` | the `/v1` semantic-ID route classifies `*store.LookupIndexUnavailableError` → 503; the workbench classifier handles timeout + `ReferenceIndexUnavailableError` only — a scan-side lookup-index failure answers the workbench's sanitized 500: a declared residue (D4), not silently absorbed |
| V36 | The helpers and stages the wiring uses | `/usr/bin/grep -n "func (d \*Daemon) abort" -A 3 host/daemon/daemon.go`; `/usr/bin/grep -n "StageConfig     =" host/daemon/daemon.go`; `/usr/bin/grep -n "func boundedTestContext" -A 4 host/daemon/m7b_deadline_test.go`; `/usr/bin/grep -n "func (h HashRef) IsZero" host/hashref/hashref.go` | `abort(stage, detail, err)` closes the store and returns the `StartupError` (:831-834); `StageConfig = "config"` (:307); `boundedTestContext` = 30 s test ctx (m7b_deadline_test.go:10-15); `HashRef.IsZero()` exists (hashref.go:75) — the unset-interpreter check is expressible |
| V37 | Every proposed new identifier is unallocated | `for tok in gradeReasonNoStoredEvidence gradeReasonBoundElsewhere gradeReasonUnauthenticated gradeReasonToolMismatch gradeReasonProofFailed gradeReasonProofIncomplete gradeReasonStructuralCandidate gradeReasonNoneValidated gradeScanTruncation gradeCandidateScanLimit proofValidator loadProofKey proofCompilerVersionLine ProofKeyFile ProofRequiredIdentities; do printf '%s=%s ' "$tok" "$(/usr/bin/grep -rn "$tok" host/ cmd/ 2>/dev/null \| wc -l \| tr -d ' ')"; done; printf 'control NewGradeUnavailable=%s' "$(/usr/bin/grep -rn "NewGradeUnavailable" host/ cmd/ \| wc -l \| tr -d ' ')"` | all 15 tokens read `0`; control `NewGradeUnavailable=7` (the render constructor fires — the instrument reads Go). `objectGradeView` is likewise free (the only `objectGrade` matches are substrings of `objectGradeUnavailableReason`, V2) |
| V38 | `Manifest.Version` is defined as the verbatim `--version` output (round 1 F1, `oc-kimi-k3`) | `sed -n '150,160p' host/archive/archive.go` (iter-254 controller, base `edb1b63`) | `// Version is the verbatim \`ailang --version\` output captured at archival.` / `Version string \`json:"version"\`` |
| V39 | `probeVersion` returns stdout verbatim; no TrimSpace on the archive path | `sed -n '467,488p' host/archive/archive.go`; `printf 'check: '; /usr/bin/grep -c 'TrimSpace' host/archive/archive.go; printf 'control: '; cat host/archive/{archive,capsulecache,check,environment}.go \| /usr/bin/grep -c 'stdout.String()'` | `return stdout.String(), nil` (:487; stderr only in error text); check **0** in `archive.go`; control **6**. Package-wide non-test TrimSpace = **1**, at `check.go:134` inside `outputHead` (error-message formatting for captured check output — not on the version path) |
| V40 | The daemon assigns the manifest version unchanged | `sed -n '684,690p' host/daemon/daemon.go` | `:688 d.interpreterVersion = m.Version`; `:689 release = releaseFromVersion(m.Version)` (the TrimSpace spelling is confined to `release`, `daemon.go:852-858`) |
| V41 | The producer-side tests pin the pinned binary's first `--version` line | `/usr/bin/grep -n 'CompilerVersion' host/evidence/producer_realstore_test.go`; `$HOME/.pinned-ailang/ailang --version \| head -1 \| od -c`; `... --version \| wc -l` | four pins `"AILANG v0.41.0"` at :39, :140, :243, :349; first line bytes `A I L A N G ' ' v 0 . 4 1 . 0 \n` (no trailing space or `\r`); full output **7** lines |
| V42 | Every validator method D4 calls exists with the called shape (round 1 F2, `gemini-3-1-pro`) | `/usr/bin/grep -n 'func (r ValidationResult) Validated\|func (r ValidationResult) Unsupported\|func (r ValidationResult) Err\|func (v \*Validator) ValidateProof\|func (r ResolutionResult) Proven\|func (r ResolutionResult) Err\|func (v \*Validator) Resolve' host/evidence/validator.go` | 7 hits: :110 `Validated() (ValidatedEvidence, bool)`, :111 `Unsupported() (UnsupportedReason, bool)`, :114 `Err() error`, :120 `ValidateProof(ctx context.Context, reportRef, expectedSubject hashref.HashRef) ValidationResult`, :208 `Proven() (ResolvedGrade, bool)`, :209 `Err() error`, :211 `Resolve(sealed ValidatedEvidence) ResolutionResult` |
| V43 | The cross-check fixture actually discriminates the two derivations (round 2, `claude-sonnet-5@claude-p`) | scratch `go run` (outside the repo) of `strings.TrimRight(strings.SplitN(out, "\n", 2)[0], "\r")` (producer, `proof_producer.go:328`) vs the `releaseFromVersion` loop (`daemon.go:852-858`) over three fixtures | `"AILANG v0.41.0 \r\nbuild abc\n"` → producer `"AILANG v0.41.0 "` vs TrimSpace `"AILANG v0.41.0"` **differ=true**; `"\nAILANG v0.41.0\n"` → `""` vs `"AILANG v0.41.0"` **differ=true**; control `"AILANG v0.41.0\n"` → both `"AILANG v0.41.0"` **differ=false**. The first fixture is THE discriminating fixture AC3/M2/M3/M11 cite, by exact bytes |
| V44 | The 101-row fetch is within the store's page cap | `/usr/bin/grep -n 'MaxSemanticIDPage' host/store/objects_by_semantic_id.go` | `const MaxSemanticIDPage = 500` (:13); `limit < 1 \|\| limit > MaxSemanticIDPage` refuses (:34) — 101 is admitted |

**Sweep note (the directive's own rule):** every sentence in this document claiming the codebase does or does not do X cites a row above. Negative claims (no non-test importers: V3; no key material: V5; no production caller: V28; no new subprocess site: V21; no `.ail` change: V24; unallocated identifiers: V37) each carry a same-call firing control. No sentence claims a validator property without V4/V9-V11, a store property without V13-V16, a render property without V25/V26, or a gate property without V20-V24.

---

*Filed by the designer role, iteration 253, unattended mission stage. The controller reads the worktree diff; no git write operations were performed by this role.*