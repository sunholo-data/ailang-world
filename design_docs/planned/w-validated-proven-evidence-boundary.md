# w-validated-proven-evidence-boundary — authority-bearing proof evidence

- **Status**: **PARKED `needs-human-review`** after TWO quorum rounds (iteration 84) — DESIGNED, not landed
- **Park reason**: round-2 `gpt5-6-sol` raises a DIRECTION objection the controller may not settle —
  a content-addressed, self-asserted report is not an *authenticated* one, so `ValidateProof` as
  designed mints from bytes an attacker with store-write access could fabricate (V28 prices the
  premise). The three candidate directions are mutually exclusive architecture choices, so this is a
  human A/B/C, not a completeness nit. See §10.1.
- **Item**: queue item 17, `w-validated-proven-evidence-boundary`
- **Filed**: 2026-08-14, iteration 84
- **Measurement base**: `4557262`
- **Instrument**: `/tmp/ailang-v0300/ailang`, AILANG v0.30.0 (`e37b370`)
- **This tranche estimate**: **3.5 days**
- **Decomposition**: **yes** — three ordered sprint-sized items, **8.5 days total** (§9)
- **Design result**: serialized proof references remain untrusted and grade `CLAIMED`; only a host
  validator may mint the sealed value from which the host-resolved API returns authority-bearing
  `PROVEN`.

This document is the first tranche of the larger queue item. A fresh inventory found no production
Go Evidence boundary, no Z3 report producer, no renderer, and no named Go-test manifest (V3, V4,
V14, V15). Implementing proof, replay, and rendering together would exceed the 3–4 day sprint
guardrail. Section 9 therefore decomposes the work without weakening the authority rule.

Every present-tense repository claim is tied to the numbered Verification Log in §11. The revised
AILANG syntax and its one-sided negative control were checked with the pinned release in `/tmp`
(V25); V18/V22 are retained measurements of the rejected round-1 arm.

---

## 1. Problem

`Evidence` has five public constructors and `gradeOf` maps none to `PROVEN`; the six mapping tests
and the proof identity are already pinned by the AILANG gate (V5, V6). This is intentional. An
AILANG proposal contains `evidence: list[Evidence]`, while the current production AILANG
constructors write an empty list (V7, V8). If the kernel merely added `ProofReport(HashRef) =>
PROVEN`, a proposal author could write that constructor with an arbitrary hash. The constructor
would confuse possession of a reference with authority to make the top-grade claim.

The missing control is an executable trust boundary. Production Go currently contains neither an
Evidence constructor/decoder nor a Z3 report producer (V3, V4). The object store verifies content
addresses on insertion and exposes object payloads on lookup, but lookup itself returns the full
payload and has no Evidence-specific size/type/success validation (V10). The transition registry
contains a useful strict, bounded JSON-codec pattern, but it is not an Evidence codec (V11).

Replay cannot fill the gap by implication. It has structured divergence at transition-source,
result-byte, and world-hash comparisons and returns success only after those comparisons (V9), but
it emits no Evidence value (V3). `RecordedEffect` remains `ATTESTED`; it is a record of an effect,
not evidence that an episode replayed without divergence.

## 2. The design question, settled

### 2.1 Question

How can the kernel carry a proof receipt while bounded I/O, hash verification, typed decoding,
solver execution, and the only authority-bearing `PROVEN` result remain at the effectful Go
boundary—and
how does the system prevent an agent from routing around that boundary by authoring identical
bytes?

### 2.2 Decision

Adopt the reviewer's direction fix. Add a public kernel receipt constructor named
`ProofReceipt(HashRef)`, explicitly treat it as untrusted, and map it to `CLAIMED`. Do **not** add
an agent-constructible kernel arm that returns `PROVEN`. AILANG v0.30.0 has no module-private
constructor or host-issued nominal capability that survives serialization, and the published
`world/types` module exports every `Evidence` constructor and `gradeOf` (V23). A foreign module can
therefore import and execute any such arm; the rejected round-1 design was measured to return
`PROVEN` from a made-up digest with no Go boundary at all (V24).

Authority is instead a Go value with unexported state and a host-only resolved grade:

1. `host/evidence.Validator` is the **only mint authority**. Its successful `ValidateProof`
   return is a `ValidatedEvidence` whose fields and construction function are unexported.
2. `host/evidence.DecodeProposal` decodes `ProofReceipt` as untrusted `ClaimedEvidence`; possession
   of the receipt neither seals it nor changes its kernel grade.
3. `host/evidence.GradeOfValidated(ValidatedEvidence) ResolvedGrade` is the only API that can
   return host `ResolvedGradeProven`. It accepts neither raw `Evidence`, decoded claims, nor a
   `HashRef`. `ResolvedGrade` is represented in Go, not serialized in the kernel ADT.
4. On reload or process restart, a `ProofReceipt` is revalidated against its report and expected
   subject before a new sealed value can exist. No serialized value regains authority by decode.
5. Any renderer introduced later accepts `ValidatedEvidence` (or the read-only resolved view from
   it), never raw `Evidence`, raw JSON, kernel `EvidenceGrade`, or receipt bytes. Until tranche 3,
   no renderer may show `PROVEN`.

The safety statement is now exact at both language boundaries: an agent can spell
`ProofReceipt`, but direct `gradeOf` execution returns `CLAIMED`; an agent cannot construct the Go
sealed value required by `GradeOfValidated`. Kernel `EvidenceGrade.PROVEN` remains a public,
agent-spellable enum value with no `Evidence -> EvidenceGrade` kernel producer in tranche 1. Its
mere spelling carries no authority, just as caller-written Go grade data carries none. That
reserved result is acceptable: removing it would break the already-pinned grade vocabulary,
while assigning it to a public Evidence constructor would recreate grade laundering.

### 2.3 Current and future grade-consuming ingress

Not every current or future grade consumer is forced through the resolver. That round-1 claim was
false: published pure AILANG consumers may call `gradeOf` directly (V23, V24), and may spell the
public `EvidenceGrade.PROVEN` vocabulary value directly. The design instead ensures that an
Evidence-to-grade bypass cannot yield `PROVEN`: exhaustive policy test
`gradeCode_test_7` pins `ProofReceipt => CLAIMED`, and the exact `gradeOf` contract pins the same
arm. A named mutation changes only that arm to `PROVEN` and must red the AILANG policy leg.

Host display/attachment ingress is narrower. The external-package API freeze test and exact named
Go-test manifest require every API capable of returning/displaying `ResolvedGradeProven` to take
`ValidatedEvidence`; adding a raw-Evidence, raw-hash, kernel-grade, or receipt overload reds
`TestPublicAuthoritySurfaceIsFrozen`. Tranche 3 adds the same rule at the first renderer.

The honest limitation is that no existing language mechanism prevents a future authorized source
change from adding a new public kernel producer or a new Go bypass. The gates make those changes
explicit and red-on-removal; they do not make repository maintainers powerless. Likewise, a
foreign module can define—or directly spell—a non-authoritative “proven” value, but it cannot
obtain host `ResolvedGradeProven` or make canonical `gradeOf(ProofReceipt(...))` return `PROVEN`
without a gated kernel change. Any future security-sensitive consumer that accepts a bare kernel
`EvidenceGrade` is outside this authority API and must be rejected by the API-freeze gate.

### 2.4 Why not signatures, secret tags, or an allowlisted hash?

- A signature would merely move mint authority to key custody and introduce key rotation,
  persistence, and recovery. No cross-host proof receipt is required here.
- A secret marker serialized beside the report becomes bearer authority and can leak or be copied.
- An allowlisted hash is mutable ambient state and still does not prove type, subject, or success.
- An exported Go struct with `Grade: PROVEN` is forgeable by any caller and is rejected.

The sealed value is process-local authority. Durable authority is reconstructed only from report
validation, not from trusting the serialized receipt.

### 2.5 Explicit failure semantics

Every validation call returns exactly one of:

```text
Validated(ValidatedEvidence) | Unsupported(UnsupportedReason)
```

The reasons are stable identifiers: `invalid_ref`, `missing`, `oversize`, `hash_mismatch`,
`wrong_semantic_id`, `wrong_interface`, `malformed`, `subject_mismatch`, `tool_mismatch`,
`proof_failed`, and `proof_incomplete`. Store/I/O cancellation is an
explicit operational error, not a grade. Neither `Unsupported` nor error contains an
`EvidenceGrade`; no failure falls back to `CLAIMED`, `ATTESTED`, or `TESTED`.

## 3. Proposed change — tranche 1

### 3.1 Kernel representation in `world/types.ail`

Append one constructor:

```ailang
  | ProofReceipt(HashRef)
```

Extend both the `gradeOf` postcondition and body with:

```ailang
  ProofReceipt(_) => CLAIMED
```

Add a seventh integer expectation to `gradeCode` for that arm, expecting the existing `CLAIMED`
code. The `HashRef` identifies the canonical typed proof report. It is a receipt pointer, not a
capability, and its kernel grade never depends on whether the referenced report happens to exist.

Do not add `ProofReport` fields to AILANG. Report parsing, byte limits, solver metadata, and object
loading are effect-boundary concerns. Do not add an aggregate function over `Proposal`: the pinned
verifier cannot encode a record containing `list[ADT]`, while a bare ADT parameter and ADT result
are supported. The revised bare-ADT seven-arm function verifies non-vacuously under v0.30.0 (V25).

**Why is this in the published package?** `Evidence`, `EvidenceGrade`, and their canonical
elimination rule already live in `world/types`, one of the four published package modules (V23).
That fact forbids placing authority there; it does not forbid placing an explicitly untrusted wire
receipt there. Validation and the only reachable `PROVEN` result remain in a new slim host package
because they require effects and an unforgeable process-local value.

### 3.2 Go representation in new `host/evidence`

Add these conceptual surfaces (names are binding; field layout is implementation detail):

```text
ClaimedEvidence              // decoded untrusted constructors, including ProofReceipt
ValidatedEvidence            // exported type, all fields unexported
ValidationResult             // accessors expose Validated or Unsupported, not writable fields
Validator.ValidateProof(ctx, reportRef, expectedSubject)
DecodeProposal(raw)          // ProofReceipt remains an untrusted claim
GradeOfValidated(sealed) ResolvedGrade
ResolvedGrade                // host enum; ResolvedGradeProven is not serialized
```

No `NewValidatedEvidence`, struct literal, `SetGrade`, raw-hash grade resolver, receipt resolver,
or exported unseal method exists. Put external-package tests in `host/evidence/authority_test.go`
so they compile with only the public API.

The new package depends on `host/hashref` and a minimal object-reader interface, not on daemon,
replay, or renderer. That avoids a cycle and makes the validator testable with a bounded fake.

### 3.3 Canonical `ProofReportV1`

The producer encodes strict canonical JSON with these fields in this order:

```text
schema             = "world/proof-report/v1"
subject            HashRef of the exact checked AILANG module bytes
compiler           HashRef of the exact executable bytes
compilerVersion    = "AILANG v0.30.0"
verified           sorted, non-empty list of required function identities
errors             = 0
counterexamples    = 0
checkPassed        = true
proofSucceeded     = true
```

`SemanticID` is `world/proof-report/v1`; a fixed `InterfaceHashV1` identifies this schema. Unknown,
duplicate, missing, non-canonical, trailing, invalid-UTF-8, and over-limit input is rejected. Raw
report payload is capped at **256 KiB** before decode; the decoded verified list is capped at 256
unique identities and every string at 1 KiB. The implementation must reject payload length before
allocating a second full copy. The chosen raw cap follows an existing 256 KiB strict-codec bound
already used for registry schema input (V11); it is a new Evidence invariant, not inherited store
behaviour.

Validation order is fixed:

1. validate canonical `HashRef`;
2. load exactly one object through `ObjectReader`;
3. reject absent and payload over 256 KiB;
4. recompute the payload hash and compare algorithm plus digest to the requested ref;
5. require exact semantic ID and interface hash;
6. strict-decode and canonical re-encode byte equality;
7. require report subject equals the expected subject;
8. require compiler hash/version equals the configured pinned executable;
9. require non-empty verified identities, `checkPassed`, `proofSucceeded`, zero errors, and zero
   counterexamples;
10. only then construct `ValidatedEvidence`.

The validator never accepts a report merely because `store.PutObject` once checked its hash. It
recomputes at consumption so alternate readers, corruption, and test fakes cannot bypass the
boundary.

### 3.4 First Z3 report producer

Add `host/evidence/proof_producer.go`. It runs only the configured absolute `AILANG_BIN`, verifies
the executable bytes and exact `AILANG v0.30.0` token, invokes `ai-check` for the subject, bounds
wall time, stdout, and stderr, and parses JSON rather than trusting process rc. It emits and stores
a report only when `check.passed=true`, `verify.verified>0`, `verify.errors=0`,
`verify.counterexample=0`, and the configured required identity set is present.

The producer takes required identities from its caller; an empty set is refused. The first
integration is a library API, not a daemon route or an agent tool. Therefore an agent cannot ask
the host to prove arbitrary source and immediately attach the result through a public network
surface. A later transition may call it only after separately ratifying that authority.

### 3.5 What crosses the AILANG/Go boundary

Across serialization, the kernel carries `ProofReceipt(reportRef)`, which always remains an
untrusted `ClaimedEvidence`. Trusted validation separately produces
`ValidatedEvidence{reportRef, subject, ...}` with unexported fields.
`GradeOfValidated(ValidatedEvidence)` is the sole bridge to host `ResolvedGradeProven`; there is no
bridge that serializes authority back into an AILANG constructor.

The pure kernel proves only the untrusted receipt-to-`CLAIMED` mapping. Go performs bounded loading, digest
verification, typed decode, solver-success checks, and provenance enforcement. No contract takes
`Proposal`, so the `record containing list[ADT]` verifier failure is avoided rather than hidden.

### 3.6 Gate changes

AILANG changes move all five coupled surfaces in one implementation:

- `world/types.ail` and byte-identical `packages/world-core/world/types.ail`;
- `EXACT_TOTAL_TESTS` 20 → observed 21 and `REQUIRED_TESTS` gains the emitted seventh
  `gradeCode` identity;
- `EXACT_TOTAL_VERIFIED` remains 5, but `gradeOf` remains named in `REQUIRED_VERIFIED`;
- the frozen four-module export manifest remains four exports while its interface/content hashes
  change;
- `scripts/world_package_ready_packet.golden.json` is regenerated canonically.

Add a named Go-test leg to `scripts/verify_go.sh` before the broad plain/race runs. It executes
`go test -json ./host/evidence -count=1`, parses only terminal `Action=pass` events for `Test...`
identities, requires the exact set declared as `REQUIRED_EVIDENCE_TESTS`, requires that set and the
observed test set are both non-empty, and pins `EXACT_EVIDENCE_TESTS` to the observed count. It
fails on missing, skipped, failed, duplicate, or extra tests. The broad tests remain defense in
depth; the named manifest is the persistent authority gate.

## 4. What the proof proves — and does not prove

The exact `gradeOf` postcondition proves that every encoded constructor maps to the specified
grade, including `ProofReceipt → CLAIMED`, and that the seven-arm match is total over the current
ADT. The seventh runtime integer case makes the new policy executable under the pinned runner.
The isolated revised shape produced `check.passed=true`, one verified function, zero verifier
errors, and zero counterexamples; changing only its body arm to `PROVEN` produced a named
counterexample (V25).

The proof does **not** prove report existence, payload size, hash integrity, schema, subject,
compiler identity, solver success, producer identity, decode provenance, or host refusal. Z3 sees
only a `HashRef`. Those properties are enforced by the Go validator and its named mutation tests.

Most importantly, the proof cannot detect a **consistent lie**: changing both the contract and
body to return `PROVEN` for an unvalidated constructor leaves Leg 1 green; this was measured in
iteration 81, while the hand-authored integer expectations turned Leg 2 red (V19, inherited).
Therefore the exact contract is necessary for totality but is not the authority oracle. The Go
boundary tests and pinned policy cases are independent statements of intended trust.

The AILANG constructor is not opaque. Claiming otherwise would be false. The unforgeable value is
the in-process Go `ValidatedEvidence`; serialized receipts never carry authority, although they
may be inputs to revalidation.

## 5. Persistent non-vacuity

The implementation has four persistent layers:

1. `gradeOf` stays named under `REQUIRED_VERIFIED["world/types.ail"]`; exact verified total stays
   5. Removing its contract or identity reds Leg 1.
2. `gradeCode_test_7` (use the observed emitted name) is added to `REQUIRED_TESTS`, and
   `EXACT_TOTAL_TESTS` moves to 21. Changing the new arm or deleting the case reds Leg 2.
3. `scripts/verify_go.sh` adds non-empty `REQUIRED_EVIDENCE_TESTS` plus
   `EXACT_EVIDENCE_TESTS`. Removing a guard, validator branch, authority-surface check, or named test reds
   the focused leg before broad tests.
4. `host/evidence/gate_mutation_test.go` runs the deterministic Go mutants against downstream
   validation/resolved-grade observables, while the AILANG mutation changes only
   `ProofReceipt(_) => CLAIMED` to `PROVEN`. Go mutants use neutering (`if false && condition`) rather than
   deletion, and each control proves the fixture reaches the success path.

The focused Go leg's anti-vacuity floor is **one discovered package, a non-empty required set, and
at least one terminal named-test pass**; implementation pins the exact non-zero count. A shell
grep over source is not acceptance evidence.

## 6. Mutation table

Rows M1–M5 land in this tranche. M6 is specified now but lands with ordered tranche 2; until then
replay has no route to `PROVEN`, which is the stronger default. Every observable is returned by or
after the mutated mechanism, never a sibling value assigned beside it.

| ID / class | Exact file and neutered edit | Named check that fires | Downstream observable and predicted failure text |
|---|---|---|---|
| M1 **arbitrary** | `world/types.ail`: change only body arm `ProofReceipt(_) => CLAIMED` to `ProofReceipt(_) => PROVEN` | `gradeCode_test_7` plus `gradeOf` verification | Agent-authored receipt reaches canonical kernel `PROVEN`; runtime fails `got 4, want 1`, while the unchanged contract yields a counterexample. |
| M2 **missing** | `host/evidence/validator.go`: change `if !ok` to `if false && !ok` | `TestMissingProofReportCannotResolveProven` | The fake reader returns absent and the full resolver is called; failure: `missing report resolved PROVEN; want unsupported missing`. |
| M3 **malformed** | `host/evidence/report_codec.go`: change strict-decode error guard to `if false && err != nil` while returning a zero report | `TestMalformedProofReportCannotResolveProven` | Malformed/trailing JSON reaches validator and unexpectedly seals; failure: `malformed report resolved PROVEN; want unsupported malformed`. |
| M4 **mismatched** | `host/evidence/validator.go`: change subject comparison to `if false && report.Subject != expectedSubject` | `TestMismatchedProofSubjectCannotResolveProven` | A valid successful report for subject A is resolved for B; failure: `mismatched report resolved PROVEN; want unsupported subject_mismatch`. |
| M5 **failed** | `host/evidence/validator.go`: change success guard to `if false && (!report.ProofSucceeded || report.Errors != 0 || report.Counterexamples != 0)` | `TestFailedProofReportCannotResolveProven` | A typed report with `proofSucceeded=false` passes the whole resolver; failure: `failed proof resolved PROVEN; want unsupported proof_failed`. |
| M6 **divergent** (tranche 2) | `host/replay/evidence.go`: change `if err != nil` after `ReplayEpisode` to `if false && err != nil`, allowing `*DivergenceError` to reach report creation | `TestDivergentReplayCannotResolveProven` | The test corrupts recorded result bytes, calls replay-evidence production then the common resolver, and unexpectedly gets `PROVEN`; failure: `divergent replay resolved PROVEN; want unsupported replay_divergent`. This is downstream of replay comparison and report validation. |
| M7 hash integrity | `host/evidence/validator.go`: neuter recomputed-hash comparison | `TestPayloadHashMismatchCannotResolveProven` | Corrupt object from a fake reader resolves; `hash-mismatched payload resolved PROVEN; want unsupported hash_mismatch`. |
| M8 wrong type | `host/evidence/validator.go`: neuter semantic/interface checks | `TestWrongReportTypeCannotResolveProven` | A different typed object resolves; `wrong report type resolved PROVEN; want unsupported wrong_semantic_id`. |
| M9 producer false-green | `host/evidence/proof_producer.go`: neuter `verify.errors != 0` guard | `TestProofProducerRefusesVerifierErrors` | Producer stores a report despite JSON errors; `producer emitted report with verify.errors=1`. |
| M10 seal bypass | `host/evidence/grade.go`: add a grade resolver accepting `HashRef`, `ProofReceipt`, raw `Evidence`, or kernel `EvidenceGrade` | `TestPublicAuthoritySurfaceIsFrozen` | External-package API inventory gains a forbidden resolver; `public authority surface exposes non-sealed PROVEN ingress`. |
| M11 named-manifest removal | `scripts/verify_go.sh`: remove one literal required name, leaving the test present | `TestEvidenceNamedManifestRejectsUnpinnedTest` in `host/verifygate/evidence_manifest_gate_test.go` | Isolated gate sees an extra observed test; `evidence test set differs from REQUIRED_EVIDENCE_TESTS`. |
| M12 projection drift | edit only `world/types.ail` | existing world-package step 3/9 | `projection hash mismatch: world/types.ail` (exact wording confirmed during implementation). |
| M13 stale ready packet | rebuild projection but retain old golden | existing world-package step 9/9 | `ready packet differs byte-for-byte from golden`. |

For M2–M9, the control first validates one good report for the same expected subject and observes a
sealed result whose `GradeOfValidated` result is `ResolvedGradeProven`. Thus a mutant cannot pass
merely because the test never reached minting. M1's control executes an agent-authored receipt and
observes `CLAIMED` before applying the one-sided mutation. For M6, the non-divergent control must
produce a replay report and resolve host `PROVEN` before the corrupted record is introduced.

## 7. Acceptance criteria

1. **Authority surface.** `ValidatedEvidence` has no exported fields or public constructor.
   `Validator.ValidateProof` is the sole mint. External package tests prove no authority-bearing
   grade API accepts a raw `HashRef`, decoded proposal evidence, receipt, or caller-written
   `EvidenceGrade`.
2. **Agent containment.** `DecodeProposal` may decode `ProofReceipt`, but it remains
   `ClaimedEvidence`; canonical `gradeOf` returns `CLAIMED`, and no host grade API accepts it.
3. **Bounded validation.** Proof reports are capped at 256 KiB before strict decode, hash is
   recomputed, semantic/interface identities match, canonical bytes round-trip, subject and
   compiler match, verified set is non-empty, and success/error/counterexample fields agree.
4. **No fallback.** Every validation failure yields its exact `UnsupportedReason` or an explicit
   operational error. No failure result carries any grade.
5. **Producer.** The pinned executable is byte/version checked; execution is time/output bounded;
   JSON fields—not rc—decide success; an empty required identity set is refused; only successful
   reports are stored.
6. **Kernel mapping.** Add only `ProofReceipt(HashRef)` in tranche 1 and map it to `CLAIMED`; extend
   exact contract/body and integer policy case. No kernel Evidence constructor produces `PROVEN`.
7. **AILANG pins.** Keep `gradeOf` named; add observed seventh test identity; move only
   `EXACT_TOTAL_TESTS` 20 → 21. Do not invent `EXACT_TOTAL_MODULES`.
8. **Go persistent gate.** Add the exact non-empty named-test manifest and exact count to
   `scripts/verify_go.sh`; add an isolated self-mutation test in `host/verifygate`.
9. **Required mutations.** M1–M5 and M7–M11 red with the named messages; controls green. M7 is the
   sole payload-hash mutation owner. M6 is an
   acceptance criterion of tranche 2 and replay remains unable to yield `PROVEN` before then.
10. **Projection/golden.** Canonical/projected `types.ail` are byte-identical, the frozen four
    exports and six tar entries remain unchanged, and the canonical ready packet golden is
    regenerated.
11. **Pinned baseline and final gates.** Before mutation, require the measured base gate from V17.
    After implementation run `AILANG_BIN=/tmp/ailang-v0300/ailang ./scripts/verify_ail.sh` and
    `AILANG_BIN=/tmp/ailang-v0300/ailang GOTOOLCHAIN=go1.25.6 ./scripts/verify_go.sh`; both must be
    green. A red base is repository failure, not change evidence.
12. **No premature display.** Tranche 1 exposes no renderer/public daemon route and no surface may
    display `PROVEN`. Tranche 3 must accept the sealed/revalidated value only.

## 8. Conflict surface

### 8.1 Five coupled AILANG moves

- `world/types.ail`: adds one untrusted constructor, `CLAIMED` contract/body arm, and integer case.
- `packages/world-core/world/types.ail`: byte-identical projection changes in the same commit.
- `scripts/verify_ail.sh`: `REQUIRED_TESTS` and Python `EXACT_TOTAL_TESTS = 20` move; shell
  `EXACT_TOTAL_VERIFIED=5` and `REQUIRED_VERIFIED["world/types.ail"]={"gradeOf"}` remain pinned.
- `scripts/verify_world_package.sh` step 4 retains the frozen four-module export manifest; step 3 sees
  new projected bytes; step 9 sees new content/interface/tarball hashes.
- `scripts/world_package_ready_packet.golden.json`: canonical byte golden changes.

`LEG1_MODULES` remains the same set and retains its anti-vacuity floor. There is no
`EXACT_TOTAL_MODULES` pin (V5, V13).

### 8.2 Go packages and gate

- New `host/evidence`: codec, validator, producer, resolved-grade API, focused tests, mutation
  harness.
- `host/store`: no schema change is required; use a narrow reader interface and existing Object.
  Tests may use Store as an integration control.
- `host/hashref`: reused unchanged.
- `host/verifygate`: gains isolated mutation coverage for the named Go manifest.
- `scripts/verify_go.sh`: gains a focused JSON-parsed exact named-test leg; broad build/plain/race
  legs remain.
- `.github/workflows/ci.yml`: no new job is required because CI already invokes both scripts;
  confirm the Go job exports the pinned binary and Go toolchain during implementation.
- `host/daemon`, `cmd/**`, `host/replay`, and renderer surfaces do not change in tranche 1.

### 8.3 Package and wire compatibility

The `world/core` interface hash changes because a public ADT changes. Content and tarball hashes
and normally byte length change. Export count remains four and tar entry count remains six (V13).
Consumers matching `Evidence` must add the new arm; the exact kernel proof catches canonical
mapping totality, but downstream package compilation is the compatibility oracle.

Proof-report objects introduce semantic ID `world/proof-report/v1` and its fixed interface hash.
They use the existing objects table; no database migration or registry head is added. The report
schema is versioned rather than decoded heuristically.

## 9. Scope, pricing, and decomposition

The queue's ~1.5–2 day price is not credible after the fresh inventory. Item 13's “multi-day
producer-and-boundary” warning is correct, but the complete obligation list is larger than one
sprint.

| Ordered document | Closes | Estimate |
|---|---|---:|
| **1. This document, `w-validated-proven-evidence-boundary`** | Proof report schema/producer; first production Evidence codec; bounded/hash/type/success validation; sealed mint authority; untrusted `ProofReceipt → CLAIMED`; host-only resolved `PROVEN`; arbitrary/missing/malformed/mismatched/failed mutations; persistent named gates. | **3.5 d** |
| **2. `w-validated-replay-evidence-boundary`** | Typed replay report and untrusted replay receipt; integrate only successful full-episode replay; bind episode/log head and interpreter set; make missing/failed/divergent replay explicitly unsupported; M6 persistent mutation. `RecordedEffect` stays `ATTESTED`. | **3.0 d** |
| **3. `w-proven-evidence-renderer-consumption`** | Renderer/read API that accepts only sealed or freshly revalidated evidence; display `PROVEN` only from that value; explicit `UNSUPPORTED` for every validation failure; end-to-end agent-forgery and restart/revalidation tests. | **2.0 d** |
| **Total** | All seven reviewer obligations | **8.5 d** |

Tranche 1 arithmetic:

| Work | Time |
|---|---:|
| Strict report/Evidence codecs and object integration | 0.75 d |
| Bounded pinned proof producer and fixtures | 0.75 d |
| Validator, sealed authority, receipt containment, resolved-grade API | 0.75 d |
| Kernel mapping, projection, golden, AILANG pins | 0.35 d |
| Named Go-test manifest and self-mutation gate | 0.45 d |
| Mutations, full pinned gates, review contingency | 0.45 d |
| **Total** | **3.5 d** |

Ordering is binding. Tranche 2 cannot reuse raw proof receipts as replay receipts. Tranche 3
cannot infer trust from either serialized constructor. Until tranche 2 lands, replay evidence
cannot yield `PROVEN`; until tranche 3 lands, no renderer may display `PROVEN` at all.

## 10. What this is NOT doing

- It does not claim an AILANG constructor is opaque or unforgeable, or that every grade consumer
  passes through Go.
- It does not accept a report because its `HashRef` is well formed or because Store once inserted it.
- It does not add replay evidence in tranche 1 or promote `RecordedEffect` above `ATTESTED`.
- It does not treat a replay cache hit, execution success, or absence of an error as replay proof.
- It does not add a daemon route, CLI verb, agent tool, renderer, or public proof-as-a-service.
- It does not define aggregate grade ordering for `list[Evidence]`, empty evidence, or proposal
  confidence.
- It does not put I/O, decoding, hashing, or solver execution in `world/`.
- It does not attach a contract to `Proposal` and does not widen the measured verifier limitation
  from “record containing `list[ADT]`” to “all ADTs.”
- It does not change store schema, world commit policy, replay semantics, effect broker, package
  version, or registry publication.
- It does not allow `PROVEN` rendering before ordered tranche 3.

## 10.1 Round 2 revision

- **Objection A — ADOPTED (option A).** The round-1 public `ValidatedProof => PROVEN` arm was a
  grade-laundering route. The replacement is public `ProofReceipt => CLAIMED`; authority-bearing
  `PROVEN` is represented only as host `ResolvedGradeProven` returned from sealed `ValidatedEvidence`.
  V23 establishes publication and V24 establishes the executable foreign-consumer bypass that
  invalidated the old direction. Section 2.3 states the future-ingress enforcement and limitation.
- **Objection B — ADOPTED verbatim.** The sentence assigning payload-hash comparison to an M4
  table-driven arm was deleted. M7 is the sole payload-hash mutation owner. The remaining rows
  each have one owner, and AC9 enumerates M7 once.

## 10.2 Round 3 quorum — BLOCKED, and why this document PARKS

Round 2's revision cleared round 1 and drew two **new** objections, both from present reviewers
(no N−1 degrade; `absent_reviewers` empty in both rounds). The controller measured the empirical
half of each rather than forwarding it, and the two land in different lanes.

**A — `gpt5-6-sol`, DIRECTION. This is the park.** Verbatim: *"The validator mints authority from a
forgeable, self-asserted JSON report. Hash recomputation proves only content integrity; semantic /
interface IDs, compiler hash, verified identities, and success flags are all public values an
attacker can encode into canonical bytes. `ValidateProof` neither reruns the proof nor authenticates
that the trusted producer created the report."* Its `catch` is precisely a demand for a negative
control this document does not contain, which is a fair reading of §11. **V28 prices the premise
without settling it:** the daemon's surface is 7 `GET` routes plus exactly one write, `POST
/v1/commit`, with no object-write route (control: 8 of 8 registrations enumerated) — but
`PutObject` has 10 non-test call sites and `host/broker/broker.go:289` stores bytes derived from an
effect result, so "writable object store" is not excluded by the transport. The reviewer's own
`proposed_fix` offers two mutually exclusive architectures — re-execute the pinned checker inside
`ValidateProof` and treat stored reports as non-authoritative cache, or authenticate reports with a
host-held signing/MAC key. Choosing between those is a design DIRECTION call with real cost
consequences (a validator that re-runs the checker is a validator that executes a compiler on every
grade resolution). Standing rule 2 forbids the controller proceeding over a contested direction and
the narrow-refinement carve-out does not apply, so this parks for a human A/B/C.

**B — `gemini-3-1-pro`, PREMISE. Correct premise, WEAKENING fix — see V27, and do not apply it
verbatim.** The objection is right that `verify.verified` is an integer, so §3.3/§3.4's "list of
required function identities" cannot come from that field. Its fix — replace the identity set with a
`verifiedCount` integer — was measured and is a **downgrade**: `verify.results[]` carries
`function` and `status` per identity (control: the same field reads `status='counterexample'` on a
one-sided mutant), so identity-level validation is available and strictly stronger than a count.
The repair is to re-point §3.3/§3.4 at `verify.results[]`. That is the DESIGNER's call, recorded
here and deliberately **not applied** — a controller-invented resolution is forbidden even when the
measurement is clean, and the carve-out's verbatim-application safeguard would have shipped the
weakening.

**Whoever resumes this document** starts from the human's A/B/C on objection A, applies the V27
repair to §3.3/§3.4, and adds the negative control objection A's `catch` asks for — a test that
hand-authors otherwise-perfect canonical `ProofReportV1` bytes and requires an explicit
`unauthenticated_report` result rather than a seal.

## 11. Verification Log

Unless labelled inherited, every command was run from repository root at `4557262` on 2026-08-14.
Negative measurements assert their roots first and include a same-scope positive control in the
same call. Glob-shaped arguments are quoted. V18, V22, V24, and V25 use only scratch files under
`/tmp` when they write.

| ID | Claim | Exact command and same-call control | Observed output |
|---|---|---|---|
| V1 | Measurement base and clean initial tree. | `git rev-parse --short HEAD; git status --short` | `4557262`; no status lines. |
| V2 | Required inputs exist in this worktree. | `test -f design_docs/implemented/w-evidence-grade-mapping.md && test -f world/types.ail && test -f scripts/verify_ail.sh && test -f host/replay/replay.go; wc -l ...` | All tests rc=0; representative counts: prior doc 620, `types.ail` 132, gate 378, replay 396. |
| V3 | Production Go has no Evidence constructor/decoder; test scope is visible. | `test -d host && test -d cmd; printf prod=; grep -rn "Evidence" host/ cmd/ --include='*.go' \| grep -v '_test\.go' \| wc -l; printf control=; grep -rn "Evidence" host/ cmd/ --include='*_test.go' \| wc -l` | roots yes; `prod=0`; same-scope test control `13`. |
| V4 | There is no production Z3 report producer; same production scope contains hashref code. | `test -d host && test -d cmd; grep -rinE 'Z3\|z3' host/ cmd/ --include='*.go' \| grep -v '_test\.go' \| wc -l; grep -rin 'hashref' host/ cmd/ --include='*.go' \| grep -v '_test\.go' \| wc -l` | `9` Z3 hits (inspection: comments); same-scope hashref control `424`. |
| V5 | Existing proof and all six integer tests are pinned; four pin mechanisms are distinct. | `nl -ba scripts/verify_ail.sh \| sed -n '130,345p' \| rg 'LEG1_MODULES\|world/types.ail\|gradeCode_test\|EXACT_TOTAL_VERIFIED\|REQUIRED_TESTS\|EXACT_TOTAL_TESTS'` | `LEG1_MODULES`; `world/types.ail:{"gradeOf"}`; six `gradeCode_test` names; shell total `5`; Python total `20`. |
| V6 | No additional PROVEN grep exists; the gate file is searchable. | `{ rg -n 'PROVEN' scripts/verify_ail.sh .github/workflows --glob '*.yml' || true; } \| wc -l; grep -c 'EXACT_TOTAL_VERIFIED' scripts/verify_ail.sh` | PROVEN lines `0`; same-file control `4`. Thus prohibition is the six pinned expectations and nothing else. |
| V7 | Evidence and grade use are confined to types/projection; proposal predicate control fires. | `test -d world; grep -rc Evidence world/*.ail; rg -n gradeOf world packages/world-core/world --glob '*.ail'; rg -n proposalMatchesWorld world packages/world-core/world --glob '*.ail' \| wc -l` | only `world/types.ail:8`; `gradeOf` only declaration/body adapter plus projection; control `10`. |
| V8 | AILANG Proposal contains `list[Evidence]`; production constructors emit empty lists. | `nl -ba world/types.ail \| sed -n '87,100p'; rg -n 'evidence: \[\]' world/transitions.ail world/contracts.ail` | field at line 97; empty-list writes/fixtures in both files. Same-scope non-empty field declaration is the control. |
| V9 | Replay has real divergence and a success result only after comparisons. | `nl -ba host/replay/replay.go \| sed -n '90,250p' \| rg 'DivergenceError\|return ReplayResult'; ... \| rg ReplayResult \| wc -l` | typed divergence definition; returns at lines 169, 216, 234; success at 242; control count `15`. |
| V10 | Store Object carries payload, verifies on put, and returns full payload on get. | `nl -ba host/store/store.go \| sed -n '88,96p;440,492p' \| rg 'type Object\|Payload\|verifyObject\|GetObject\|SELECT interface_hash_ref'; rg -n 'func \(s \*Store\) (PutObject\|GetObject)' host/store/store.go \| wc -l` | payload field; `verifyObject` call; SQL payload lookup; two-method control. No Evidence-specific bound appears in this scope. |
| V11 | A strict bounded JSON pattern exists in transition registry. | `test -d host/transitionreg; rg -n 'DisallowUnknownFields\|multiple JSON values\|trailing JSON\|maxRevisionRaw\|maxSchemaRaw' host/transitionreg/codec.go; rg -n 'func parseJSON\|func DecodeRevision' host/transitionreg/codec.go \| wc -l` | raw caps include 262144 and 16777216; trailing/multiple checks; two function controls. |
| V12 | Go authority-relevant Proposal is a different narrow struct and has no evidence field. | `nl -ba host/transitionreg/bind.go \| sed -n '20,30p'; nl -ba world/types.ail \| sed -n '87,100p'` | Go fields are transition/interpreter/epoch/caps/effects; AILANG control contains evidence at line 97. |
| V13 | Projection and fixed package surfaces are live. | `cmp -s world/types.ail packages/world-core/world/types.ail; echo types_cmp=$?; rg -n 'EXPECTED_MODULE_COUNT\|types.ail\|exact package fields\|tar contains exactly\|ready-packet golden' scripts/build_world_package.sh scripts/verify_world_package.sh` | `types_cmp=0`; four-module build; four exports; six-entry and golden gate lines. |
| V14 | Production host/cmd has no grade/PROVEN renderer consumer; daemon scope is visible. | `test -d host && test -d cmd; { rg -n 'EvidenceGrade\|gradeOf\|PROVEN' host cmd --glob '*.go' --glob '!**/*_test.go' || true; } \| wc -l; rg -n 'handleObject\|handleWorld' host/daemon --glob '*.go' \| wc -l` | target `0`; same-scope daemon handler control `6`. |
| V15 | Go gate runs broad tests but has no named-test manifest today. | `{ rg -n 'REQUIRED.*TEST\|EXACT.*TEST\|test2json\|go test -json' scripts/verify_go.sh || true; } \| wc -l; rg -n 'go test ./\.\.\.' scripts/verify_go.sh \| wc -l` | named-pin target `0`; broad-test control `4`. |
| V16 | Go gate requires pinned AILANG and active Go is separately guarded. | `nl -ba scripts/verify_go.sh \| sed -n '19,41p;79,89p;102,126p'` | exact `v0.30.0` token check; go1.26.0–1.26.5 denylist; build/plain/race legs. |
| V17 | Pinned AILANG baseline is green and non-vacuous. | `export AILANG_BIN=/tmp/ailang-v0300/ailang; /tmp/ailang-v0300/ailang --version; ./scripts/verify_ail.sh` | v0.30.0 commit e37b370; rc=0; 5/5 identities across 11 modules; 20 named tests; package 9/9; terminal PASS. |
| V18 | The now-rejected round-1 six-constructor bare-ADT mapping checked and verified with pinned release. | `apply_patch` created `/tmp/iter84_proven_probe.ail`; `export AILANG_BIN=/tmp/ailang-v0300/ailang; AILANG_RELAX_MODULES=1 /tmp/ailang-v0300/ailang ai-check /tmp/iter84_proven_probe.ail > /tmp/iter84_proven_probe.json; python3 -c 'import json; d=json.load(open(...)); print(d["check"]["passed"],d["verify"]["verified"],d["verify"]["errors"],d["verify"]["counterexample"])'; /tmp/ailang-v0300/ailang test /tmp/iter84_proven_probe.ail` | `True 1 0 0`; `gradeCode_test_1` and `_2` pass. Temporary-path warning only. This is retained history, not evidence for the revised mapping. |
| V19 | A consistent contract/body lie evades Z3 but runtime policy tests red. | **Inherited controller measurement, iteration 81**: add the same `=> PROVEN` arm to contract and body, run pinned gate; compare Leg 1 and Leg 2. | Leg 1 `verify.verified=1`, errors/counterexamples `0`; six integer expectations make Leg 2 red. Not re-run because mutating live sources is outside this design-only task. |
| V20 | HashRef is nominal in Go with unexported fields and strict constructors. | `nl -ba host/hashref/hashref.go \| sed -n '43,125p'; rg -n 'func (New\|Parse\|Sum)' host/hashref/hashref.go` | `HashRef` fields `algo`, `digest` unexported; malformed/unsupported cases return `HashError`; constructor controls present. |
| V21 | No `EXACT_TOTAL_MODULES` variable exists; module inventory still has a positive control. | `{ rg -n 'EXACT_TOTAL_MODULES' scripts/verify_ail.sh || true; } \| wc -l; rg -n 'LEG1_MODULES' scripts/verify_ail.sh \| wc -l` | target `0`; same-file `LEG1_MODULES` control is non-zero (multiple lines). |
| V22 | **CONTROLLER-RUN (iteration 84), and it supplies the arm V18 lacks: V18's `verified=1` is NON-VACUOUS — the contract genuinely constrains the new `ValidatedProof` arm.** V18 ran only the positive arm, and a green there is equally consistent with "the contract binds" and "the contract is decorative". | Reproduced V18 first-party, then mutated **only the body** arm (`ValidatedProof(_) => CLAIMED`) while leaving the contract at `PROVEN`; asserted the mutant LANDED by sha256 (`3d7558d4…` → `3b125da3…`, `diff` = the single line 29); re-ran `ai-check` with the pinned binary, reading stdout and stderr to separate files (a warning on stdout otherwise voids the JSON parse). | Positive arm reproduces **without** `AILANG_RELAX_MODULES=1` as well as with it (`check.passed=true`, `verified=1`, `errors=0`, `counterexample=0`, `skipped=0`) — so V18's relax flag was **not** load-bearing and that narrowing need not travel with the finding. Mutant: **rc=1**, `verified=0`, `counterexample=1`, and Z3 names the witness exactly — `$p_e = (ValidatedProof (mk_HashRef "!0!" "!1!"))`. This does **not** contradict V19: mutating **one** side reds, mutating **both** stays green. Both facts hold, and §4 states both. |
| V23 | `world/types` is one of the four published package-module exports; the same exact set is pinned at every package gate surface. | `readonly EXPORTS=(world/types world/contracts world/transitions world/logepoch); rg -n 'readonly EXPORTS=\(world/types world/contracts world/transitions world/logepoch\)\|"exports": \{"modules":\["world/types", "world/contracts", "world/transitions", "world/logepoch"\]\}\|exact export set\|Modules:\[\]string\{"world/types","world/contracts","world/transitions","world/logepoch"\}\|world/types world/contracts world/transitions world/logepoch' scripts/verify_world_package.sh; printf 'control=%s\n' "${#EXPORTS[@]}"` | Matches at lines 34, 120, 153, 175, and 239; same-scope positive control `control=4`. Thus “four exports” means four exported modules, not four functions. |
| V24 | A foreign pure-AILANG consumer can import the rejected public constructor and execute `gradeOf(ValidatedProof(made-up HashRef)) => PROVEN`; a nonexistent-constructor negative control fails at import. | **CONTROLLER-RUN at `4557262` with pinned v0.30.0**: copied `world/types.ail` and `world/logepoch.ail` to scratch; applied the round-1 constructor/contract/body change and asserted three `ValidatedProof` occurrences; checked and tested `world/consumer.ail` importing `Evidence`, `EvidenceGrade`, `gradeOf`, `ValidatedProof`, and `HashRef`; in the same scratch scope replaced that import/call with `NoSuchCtor` and rechecked. | Positive: `ailang check` rc=0, `No errors found!`; `ailang test` passes `launder_code_test_1`, whose `PROVEN` arm returns 4. Negative control: rc=1, `Error: IMP010: symbol 'NoSuchCtor' not exported by 'world/types'`. No Go boundary, decoder, validator, or sealed value occurs on the positive path. |
| V25 | Revised `ProofReceipt => CLAIMED` syntax checks, verifies, and executes; a same-scope one-sided body mutation to `PROVEN` is non-vacuously rejected. | `apply_patch` created `/tmp/iter84_receipt_probe.ail` and identical `/tmp/iter84_receipt_mutant.ail` except the mutant body arm; `export AILANG_BIN=/tmp/ailang-v0300/ailang; $AILANG_BIN ai-check /tmp/iter84_receipt_probe.ail >/tmp/iter84_receipt_probe.stdout 2>/tmp/iter84_receipt_probe.stderr; $AILANG_BIN test /tmp/iter84_receipt_probe.ail >/tmp/iter84_receipt_test.stdout 2>/tmp/iter84_receipt_test.stderr; $AILANG_BIN ai-check /tmp/iter84_receipt_mutant.ail >/tmp/iter84_receipt_mutant.stdout 2>/tmp/iter84_receipt_mutant.stderr`; `python3` loaded each stdout JSON and printed `d['check']['passed'], d['verify']['verified'], d['verify']['errors'], d['verify']['counterexample']`; `rg 'PASS\|FAIL\|gradeCode\|counterexample\|ProofReceipt' /tmp/iter84_receipt_test.stdout /tmp/iter84_receipt_mutant.stdout /tmp/iter84_receipt_mutant.stderr` inspected runtime and negative-control outputs. | Positive rc=0: `check.passed=True`, `verify.verified=1`, `verify.errors=0`, `verify.counterexample=0`; both runtime tests pass, including made-up receipt → code 1. Mutant rc=1: `check.passed=True`, `verify.verified=0`, `verify.errors=0`, `verify.counterexample=1`; witness is `(ProofReceipt (mk_HashRef "!0!" "!1!"))`. |
| V26 | **CONTROLLER-RUN (iteration 84), and it is the arm that makes the round-2 fix a MEASUREMENT rather than a promise: the revised design DEFEATS the exact attack that defeated round 1, run in BOTH arms with only the design as the variable.** V24 and V25 each measure one design; neither runs the *attack* against the *fix*. | Rebuilt V24's scratch tree from `world/types.ail` + `world/logepoch.ail` at `4557262`, applied the **round-2** change instead (asserted LANDED: 3 `ProofReceipt` occurrences at lines 29/49/58), and re-ran V24's **byte-for-byte identical attack module** — a foreign `world/consumer` importing `gradeOf` and the receipt constructor, minting from the same literal `digest: "i-made-this-up"`, with an inline test asserting the result is `PROVEN` (code 4). | **Outcomes DIFFER, which is the whole evidence.** Round-1 arm: the attack test **PASSES** (`launder_code_test_1`, rc=0) — the foreign module really obtains `PROVEN`. Round-2 arm: `ailang check` still rc=0 `No errors found!` (so the refusal is *semantic*, not a type error the attacker would notice), and the attack test **FAILS** — `✗ attack_code_test_1`, **`test 0: expected 4, got 1`**, rc=1. `1` is `CLAIMED`. Same attacker, same literal, same instrument, same scratch scope; the only variable is the kernel arm. Note what this does NOT prove, per §2.3's own limitation: it bounds the `Evidence -> EvidenceGrade` route, not a future consumer that spells `EvidenceGrade.PROVEN` directly. |
| V27 | **CONTROLLER-RUN (iteration 84), round-2 quorum. `gemini-3-1-pro`'s PREMISE is CORRECT and its `proposed_fix` is an unnecessary WEAKENING — the identities the design needs ARE emitted, one field over from the one the doc named.** The objection: §3.3/§3.4 require a list of required function identities, but `verify.verified` is an integer, so "the Go producer cannot extract or validate identities that the AILANG binary does not emit". | Re-read the pinned binary's own `ai-check` JSON and enumerated the `verify` object rather than the single field the objection names, with the one-sided mutant of V22 as the same-instrument negative control. | `verify.verified` is indeed `int = 1` — the premise holds. But `verify.results` is a **list** whose entries carry keys `['function', 'status', 'duration']`: `function='gradeOf'`, `status='verified'`. Control (mutant, same instrument): `verify.verified = 0` while `results[0]` still reads `function='gradeOf'`, `status='counterexample'` — so the field distinguishes verified from unverified **per identity**. Therefore the producer CAN validate a sorted non-empty identity set, from `verify.results[].function` filtered on `status == "verified"`. Adopting the proposed `verifiedCount` integer would discard identity-level checking (a count of 1 cannot say *which* function was verified) — strictly weaker than what the instrument supports. **The correct repair is to re-point §3.3/§3.4 at `verify.results[]`, not to downgrade the field.** This is a DESIGNER decision, not a controller one, and it is recorded here rather than applied. |
| V28 | **CONTROLLER-RUN (iteration 84). The empirical half of `gpt5-6-sol`'s round-2 objection: is a "writable object store" excluded by the daemon's surface? NO — the read-only HTTP surface does not exclude it, and there is a non-test path that stores bytes derived from an effect result.** This does not resolve the objection; it prices it. | Enumerated the daemon's route registrations and every non-test `PutObject` caller, each with a same-scope control. | Daemon registers **8** routes: seven `GET /v1/*` and exactly one write, `POST /v1/commit` (`host/daemon/daemon.go:461-468`); control — total `HandleFunc` registrations in the same scope = **8**, so the enumeration is complete, and there is **no object-write route**. But `PutObject` has **10** non-test call sites (control: `GetObject` non-test callers = **16**), and `host/broker/broker.go:289` stores `resultObject(result)` — the bytes of an effect invocation's result. So the reviewer's premise is **not** excluded by the transport: at least one non-test, agent-adjacent path writes object bytes whose content derives from an effect the agent requested. **Stated at its true strength: this shows the threat model is live enough to need a decision, NOT that a full forgery has been demonstrated** — no probe here constructs canonical `ProofReportV1` bytes through that path or seals them. |

## 12. Related

- `design_docs/implemented/w-evidence-grade-mapping.md` — predecessor and binding deferred
  authority obligations.
- `design_docs/HUMAN-SURFACE.md` — trust gradient and grade-laundering prohibition.
- `design_docs/coding-standards.md` — pure core, effect boundary, slim kernel, and honest gates.
- `design_docs/DESIGN.md` §§1, 14 — immutable world and controlled-change architecture.
- Planned successor: `w-validated-replay-evidence-boundary`.
- Planned successor: `w-proven-evidence-renderer-consumption`.
