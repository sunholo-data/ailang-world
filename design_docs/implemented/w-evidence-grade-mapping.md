# w-evidence-grade-mapping — total mapping for the ratified evidence representation

- **Status**: Planned — **CLEARED for sprint-planner**: quorum round 2 = `gemini-3-1-pro` PASS,
  `gpt5-6-sol` REJECT on one non-directional point, resolved by the controller under the
  narrow-refinement carve-out with the reviewer's verbatim fix (§12)
- **Item**: queue item 13, `w-evidence-grade-mapping`, clause-5
- **Filed**: 2026-08-11, human attended
- **Revised**: 2026-08-13 — designer revision after a both-reject round 1, then a bounded
  controller carve-out revision after round 2 (measured before applied, V25/V26)
- **Estimated**: 0.65 day; within the queued ~0.5–1 day band (§9)
- **Measurement base**: `2ef2271`, 2026-08-13
- **Instrument**: `/tmp/ailang-v0300/ailang`, AILANG v0.30.0 (`e37b370`)
- **Files changed by implementation**: `world/types.ail`,
  `packages/world-core/world/types.ail`, `scripts/verify_ail.sh`, and
  `scripts/world_package_ready_packet.golden.json`
- **Design result**: total, proved grading for the five existing `Evidence` constructors;
  `PROVEN` deliberately remains unreachable.

This revision chooses **Resolution B — representation-only**. It withdraws the earlier proposal
to add `ProofReport` and `ReplayReport`. A top-grade carrier that an agent can mint from an
unchecked `HashRef` would turn the current representation gap into a grade-laundering authority
gap. The validation boundary and the first real producer integration are a separately priced
follow-on, not prose presumed by this item.

Every present-tense codebase statement is backed by a command and observed output in §11.
Controller-supplied measurements are attributed there. The verifier probes were isolated under
`.probe/`; they are measurements, not implementation artefacts.

---

## 1. Problem

The human surface binds displayed facts to `PROVEN > TESTED > ATTESTED > CLAIMED` and requires
ratification to produce a total mapping from kernel evidence variants (V1). The kernel exposes
five `Evidence` constructors and no grading function: `CompilerOutput`, `TestReport`,
`HumanApproval`, `AiReview`, and `RecordedEffect` (V2–V4). The ratified text settles three arms
and leaves `CompilerOutput` and `HumanApproval` unmapped (V1).

This item closes that total-mapping defect only. It does not add proof/replay carriers. Therefore
`PROVEN` remains unreachable from `Evidence` after this item, explicitly and intentionally.

That restriction is necessary because `Proposal.evidence` is agent-authored, yet current
proposal verification does not read it; all production AILANG proposal constructors supply an
empty list; production Go has no `Evidence` constructor or decoder; and the repository has no Z3
proof-report producer (V20–V22). Replay machinery exists and detects divergence, but no path turns
its result into `Evidence` (V8, V23). The earlier design's statement that producers “must” mint a
top-grade variant only after success had no enforcing boundary. That statement is withdrawn.

The risk is prospective rather than live: no current producer writes non-empty `Evidence`, no
current consumer reads it, and item 14's renderer is unbuilt (V20–V22). That does not make an
unvalidated top-grade constructor safe to ratify. A type is the contract future producers will
inherit, and a freely constructible `ProofReport(HashRef)` would let any proposal author claim
the maximum grade before any renderer or validation boundary existed.

## 2. The design question, settled

### 2.1 Question

What atomic change makes the mapping total and gate-enforced without granting unvalidated mint
authority for `PROVEN` or exceeding the queued 0.5–1 day band?

### 2.2 Candidate answers

#### A. Add trust grades until every constructor fits

Rejected. New `COMPILED` or `APPROVED` grades would change the ratified cross-surface vocabulary
and introduce ordering questions not answered by HUMAN-SURFACE. The existing four grades are
sufficient to classify all five current constructors.

#### B. Add proof/replay carriers and leave mapping implicit

Rejected. Carriers alone do not create one canonical mapping, and v0.30.0 accepts a
non-exhaustive ADT match (V11). More importantly, a public constructor containing only `HashRef`
does not establish who may mint it or whether its referent exists, matches its hash, decodes as
the expected report, records success, or represents non-divergent replay.

#### C. Keep five variants and use a lowest-grade default

Rejected. A default obscures policy and lets later constructors inherit a grade silently.
`HumanApproval` is a recorded sovereign-human action, not an agent claim; `CompilerOutput` is a
recorded tool result, not merely prose. Both are `ATTESTED`, and every constructor must have an
explicit arm.

#### D. Add proof/replay carriers and prove a seven-arm function

Rejected in revision round 1. This was the original decision. The proof would establish only
constructor-directed grading. It would not establish report existence, integrity, type, success,
or producer authority. Existing `TestReport → TESTED`, `RecordedEffect → ATTESTED`, and
`AiReview → CLAIMED` semantics are also constructor-directed (V1), but symmetry is not enough:
forging `AiReview` gains only `CLAIMED`, whereas forging a proposed report carrier gains
`PROVEN`, the maximum grade. The asymmetry of consequence makes unchecked top-grade minting the
cardinal grade-laundering case even if weaker constructors retain referent obligations.

The requested validation and producer scope also cannot be absorbed honestly. There is no
production Go `Evidence` boundary and no Z3 report producer to wire (V21, V22). Adding both
carriers would therefore require building the repository's first evidence encode/decode path,
bounded object loading, hash verification, typed report codecs, explicit proof/replay success
results, and two integrations—not merely attaching adapters to two existing producers.

#### E. Keep five variants and prove an explicit five-arm function

**Decision: E (Resolution B).** Add `EvidenceGrade` and a pure `gradeOf(Evidence)` with an exact
five-arm contract. Execute six inline cases through a private integer adapter: one per constructor
plus both boolean shapes of `TestReport`. Pin the verified function and every emitted test
identity in the existing primary gate.

This satisfies HUMAN-SURFACE §7.2's actual requirement that ratification produce a total mapping
while preserving the safety fact that no current `Evidence` value yields `PROVEN`. The follow-on
in §12 owns mint authority, validation, producer wiring, and the six required failure modes.

### 2.3 Exact semantic decisions

The trust vocabulary remains:

```text
PROVEN > TESTED > ATTESTED > CLAIMED
```

The result type is:

```ailang
export type EvidenceGrade
  = PROVEN
  | TESTED
  | ATTESTED
  | CLAIMED
```

`gradeOf` is total over decoded `Evidence` and returns one of the four ratified grades. This item
defines no behaviour for absent, unreadable, or malformed inputs because it performs no loading or
decoding. The validated-boundary follow-on must represent those outcomes explicitly, for example as
`GradeReadResult = Graded(EvidenceGrade) | Unsupported(UnsupportedReason)`, without silently
converting them to a trust grade. An unreadable reference is not an `EvidenceGrade`; it is the
result of attempting to obtain validated evidence, and it belongs to that follow-on or to the
renderer extension (V1, R2-CARVE-OUT).

The total mapping is:

| Existing constructor | Grade | Reason |
|---|---|---|
| `TestReport(HashRef, bool)` | `TESTED` | Ratified mapping. The boolean records outcome; grade records method, so true and false have the same grade. |
| `CompilerOutput(HashRef)` | `ATTESTED` | A recorded deterministic tool outcome, below a named test and not a proof. |
| `HumanApproval(HashRef)` | `ATTESTED` | A recorded sovereign-human act; authority is not mathematical proof. |
| `RecordedEffect(HashRef)` | `ATTESTED` | Ratified mapping for a record of an effect-boundary event. |
| `AiReview(HashRef, float)` | `CLAIMED` | Ratified mapping for an agent assessment and confidence. |

There is no `PROVEN` arm. This is a deliberate non-reachability property, not a missing case.
The exact contract pins all five choices and totality; it does not validate any `HashRef`.

### 2.4 Mint authority and the deferred boundary

Mint authority for `PROVEN` is unresolved and no constructor purporting to grant it lands here.
The follow-on must answer both who may mint validated proof/replay evidence and what executable
boundary enforces that authority. An agent-authored `Proposal.evidence` value is not sufficient
authority merely because it decodes as an ADT constructor.

Until that follow-on lands, a renderer **MUST NOT display `PROVEN` for any kernel `Evidence`**.
It may display the mapped grades for successfully decoded existing values. It has no kernel grade
to display when a value cannot be obtained at all, and MUST NOT substitute one. Item 14 must not infer `PROVEN` from a raw hash,
replay presence, solver-related text, or any unratified host object.

### 2.5 Proposal confidence

`Proposal.confidence: float` remains out of scope. It has no adjacent evidence reference, so the
surface must omit it or label it `UNSUPPORTED` (V1, V2). `gradeOf` does not accept bare floats.

## 3. Proposed change

### 3.1 Placement: `world/types.ail`

Add the grade ADT and `gradeOf` immediately after the existing five-constructor `Evidence` type.
Do not change `Evidence`. `world/types.ail` currently contains declarations but no functions,
contracts, or inline tests; the implemented library assigns it ownership of the `Evidence`
semantic surface, while `world/contracts.ail` owns proposal/commit predicates (V4, V15).

The mapping is the elimination rule for the kernel ADT, not transition policy. Co-location keeps
the type and its total interpretation in one exported interface.

Alternatives cost more or weaken ownership:

- `world/types.ail` changes no module identity or package export, though its interface and package
  hashes change (V5–V7).
- A new `world/grades.ail` would add a twelfth Leg-1 module, fifth package module/export, sixth
  package `.ail` file, seventh tar member, and manifest/build allowlist changes (V5–V7).
- `world/contracts.ail` would mix a type eliminator with transition predicates and still move
  proof/test totals and package hashes (V5, V6, V15).
- An optional package could be absent or version-skewed while kernel `Evidence` still exists.

**Why is this not a package?** Every consumer must share one interpretation of the kernel's
closed `Evidence` ADT. Permitting packages to redefine `Evidence → EvidenceGrade` would let the
same value acquire a different epistemic meaning by package version. This is a kernel elimination
rule, not a domain policy, projection, tool, or effect. Co-location is the smallest interface
change and introduces no new kernel module.

### 3.2 Exact AILANG code

The existing `Evidence` declaration remains byte-for-byte in shape:

```ailang
export type Evidence
  = CompilerOutput(HashRef)
  | TestReport(HashRef, bool)
  | HumanApproval(HashRef)
  | AiReview(HashRef, float)
  | RecordedEffect(HashRef)
```

Add:

```ailang
-- The four ratified trust grades. gradeOf is total over decoded Evidence and
-- returns one of these four; absent/unreadable/malformed input is NOT a grade
-- and belongs to the deferred validated-boundary follow-on.
export type EvidenceGrade
  = PROVEN
  | TESTED
  | ATTESTED
  | CLAIMED

-- Canonical grading for the ratified five-constructor representation.
-- No current Evidence constructor has authority to mint PROVEN.
export func gradeOf(e: Evidence) -> EvidenceGrade ! {}
ensures { result == match e {
  CompilerOutput(_) => ATTESTED,
  TestReport(_, _) => TESTED,
  HumanApproval(_) => ATTESTED,
  AiReview(_, _) => CLAIMED,
  RecordedEffect(_) => ATTESTED
} }
{
  match e {
    CompilerOutput(_) => ATTESTED,
    TestReport(_, _) => TESTED,
    HumanApproval(_) => ATTESTED,
    AiReview(_, _) => CLAIMED,
    RecordedEffect(_) => ATTESTED
  }
}

-- Private adapter because v0.30.0 inline-test expected values cannot be ADT
-- identifiers. Integer codes are test plumbing, not an exported ordering.
func gradeCode(e: Evidence) -> int
tests [
  (CompilerOutput({ algo: "sha256", digest: "compiler" }), 2),
  (TestReport({ algo: "sha256", digest: "tests-pass" }, true), 3),
  (TestReport({ algo: "sha256", digest: "tests-fail" }, false), 3),
  (HumanApproval({ algo: "sha256", digest: "approval" }), 2),
  (AiReview({ algo: "sha256", digest: "review" }, 0.8), 1),
  (RecordedEffect({ algo: "sha256", digest: "effect" }), 2)
]
{
  match gradeOf(e) {
    PROVEN => 4, TESTED => 3, ATTESTED => 2,
    CLAIMED => 1
  }
}
```

There are six cases because `TestReport` has two policy-relevant boolean shapes. The exact
contract specifies policy; the private adapter only makes each arm executable under the pinned
runner. The runner derives identity from the function carrying `tests`, so emitted identities
are `gradeCode_test_1` through `gradeCode_test_6`, subject to confirming exact JSON output before
committing the manifest (V19, V24).

### 3.3 Gate pins moved in the same commit

In `scripts/verify_ail.sh`, set:

```python
"world/types.ail": {"gradeOf"},
```

Move `EXACT_TOTAL_VERIFIED` from 4 to 5. Add all six emitted `gradeCode` test identities to the
Leg-2 manifest and move `EXACT_TOTAL_TESTS` from 14 to the observed total, expected to be 20.
Named identities are the acceptance oracle; counts are secondary drift guards (S1).

`LEG1_MODULES` remains unchanged because no module is added (V5). A private function's inline
tests are collected, and an exported contracted `gradeOf` still verifies when called by that
private adapter (V24). Those facts make the Leg-2 and Leg-1 layers compatible.

### 3.4 Package projection and ready packet

Run `./scripts/build_world_package.sh` after editing canonical `world/types.ail`. It copies the
four allowlisted canonical modules through a fresh staging directory and replaces the package
projection wholesale (V6). Require `packages/world-core/world/types.ail` to be byte-identical to
the canonical file.

Regenerate `scripts/world_package_ready_packet.golden.json` with
`pkgproj.RecomputeReadyPacket` and `pkgproj.EncodeReadyPacket`; do not hand-edit derived fields.
The repository has a verifier but no dedicated regeneration command (V16). The grade type and
function change `contentHash`, `tarballSHA256`, and ordinarily `tarballBytes`. **`interfaceHash`
does NOT move — CORRECTED 2026-08-14 (iter-85), measured, after this stale sentence propagated a
false premise into a second item's acceptance criteria.** `InterfaceHash`
(`host/pkgproj/pkgproj.go:86`) hashes only `manifest.Package.Name`, `Edition`, `AILANG`, the
sorted `Exports.Modules` list and `Effects.Max`; it never reads a source file, so **no** change to
a module's *contents* can move it — only a change to the manifest's exported module inventory can,
and this item changes no inventory. (The adjacent `ContentHash` *does* read files, which is what
makes the original claim plausible.) Item 15's design doc inherited this sentence as precedent and
built `AC9` on it, requiring a hash move the hash function cannot produce; item 13's own sprint
plan had already diagnosed it and the correction was never written back here, which is why it
recurred. The correct assertion is the *invariance*: `interfaceHash` unchanged is itself the
check.
The packet remains four exports and six tar entries because neither `Evidence` constructors nor
module inventory changes (V6, V7).

## 4. What the proof proves — and does not prove

The exact postcondition proves the selected result for every encoded value of the current
five-constructor `Evidence` ADT. Consequently `gradeOf` never returns `PROVEN`, and after the
round-2 carve-out there is no `UNSUPPORTED` constructor for it to return. A positive user-ADT probe
verifies, and a false postcondition produces a real constructor counterexample, so the proof is
non-vacuous (V9, V10).

The proof does not establish that any referenced report exists, is authentic, or records
success. That limitation is inherited by the existing constructors, but it cannot justify a new
unchecked top-grade carrier. The six runtime cases execute every current arm and both
`TestReport` booleans; they do not validate referents.

The contract is necessary because v0.30.0 does not reject a non-exhaustive ADT match (V11).
With the exact match-equality contract present, a missing constructor arm produces a verifier error;
the exhaustive control verifies (V12). `verify_ail.sh` parses JSON errors and counterexamples,
so a future constructor added without a mapping makes Leg 1 red even when process rc alone would
be misleading (V5, V13).

No aggregate proposal-grade function lands. A verified function taking `Proposal` reaches its
`list[Evidence]` and encounters the measured record-containing-ADT encoder limitation; the bare
ADT function does not (V13, V14). Aggregation also needs rules for empty evidence and for an
unobtainable value, which belong to a separately ratified decision surface.

## 5. Persistent non-vacuity

The persistent proof has three layers:

1. `gradeOf` is named in `REQUIRED_VERIFIED["world/types.ail"]`; removing its contract or identity
   reds Leg 1.
2. The exact verified total moves to five as a secondary guard.
3. Six emitted `gradeCode_test_N` identities are pinned in Leg 2; deleting or changing an arm
   removes or fails a named identity.

The implementation must record deterministic mutations and controls. No random input, timing
threshold, solver exit code, or probabilistic gate is an oracle (S6).

## 6. Mutation table

| ID | Mutation | Expected RED | Discrimination/control |
|---|---|---|---|
| M1 | Add `FutureEvidence(HashRef)` without a `gradeOf` arm. | JSON reports `verify.errors > 0`; Leg 1 fails. | Add an explicit `ATTESTED` arm; same file verifies (V12 shape). |
| M2 | Change `CompilerOutput → ATTESTED` to `TESTED`. | Its `gradeCode` identity fails. | Both `TestReport` identities stay green. |
| M3 | Change `HumanApproval → ATTESTED` to `CLAIMED`. | Approval identity fails. | `AiReview` stays green. |
| M4 | Change `TestReport → TESTED` to `ATTESTED`. | Both boolean identities fail. | Four other constructor families stay green. |
| M5 | Change `AiReview → CLAIMED` to `ATTESTED`. | AI-review identity fails. | Human approval stays green. |
| M6 | Change `RecordedEffect → ATTESTED` to `CLAIMED`. | Effect identity fails. | Compiler output stays green. |
| M7 | Add any current arm returning `PROVEN`. | Exact contract or its matching named policy test fails. | Restore the specified arm; proof and all six tests pass. |
| M8 | Remove `gradeOf` from `REQUIRED_VERIFIED`. | Manifest inspection fails AC3; required identity disappears from Leg 1 output. | Restore literal identity; Leg 1 reports it verified. |
| M9 | Edit only canonical `world/types.ail`. | Package step 3/9 fails SHA equality. | Rebuild projection; equality returns. |
| M10 | Rebuild projection but keep old golden. | Package step 9/9 fails byte comparison. | Recompute canonical packet; 9/9 passes. |
| M11 | Use a non-pinned AILANG. | Package gate refuses the version. | Pinned v0.30.0 gives measured 9/9 baseline (V7). |

M8 is an explicit review obligation because the current script has no separate mutation test for
its own required-identity dictionary. The literal identity and named Leg-1 output, never a count
alone, are the acceptance evidence.

## 7. Acceptance criteria

1. **Representation-only shape.** `Evidence` remains exactly the existing five constructors.
   Add `EvidenceGrade` with exactly the **four** ratified result constructors in §3.2
   (`PROVEN | TESTED | ATTESTED | CLAIMED` — the round-2 carve-out removed `UNSUPPORTED`; §12).
   Fail if `ProofReport`, `ReplayReport`, any other evidence carrier, or an `UNSUPPORTED`
   grade constructor is added.

2. **Canonical total function.** `gradeOf` is pure, contains the exact five-arm mapping and the
   exact postcondition, `EvidenceGrade` has exactly the four ratified constructors, and no arm
   returns `PROVEN`. Fail on a default arm, missing arm, altered grade, effect, or referent lookup.

3. **Named proof pin.** Add exactly `gradeOf` to
   `REQUIRED_VERIFIED["world/types.ail"]`; require `gradeOf=verified`, zero errors and zero
   counterexamples. Move exact verified total 4 → 5 only as a secondary guard.

4. **Named runtime pins.** Confirm runner JSON and add all six emitted `gradeCode` test identities
   to Leg 2. Expected names are `gradeCode_test_1`…`gradeCode_test_6`; observed output is
   authoritative. Move total 14 → observed 20. Fail on a missing/non-pass identity or total drift.

5. **Future-variant RED/control.** On an isolated copy add `FutureEvidence(HashRef)` without an
   arm and require `verify.errors > 0`; add the explicit arm and require `gradeOf=verified`,
   `errors=0`. Do not use process rc as verdict.

6. **Policy RED/control.** Change only `CompilerOutput → ATTESTED` to `TESTED`; require only its
   named `gradeCode` case to fail, then restore and require all six cases to pass.

7. **No top-grade mint.** Search the implementation diff and resulting `Evidence` declaration;
   fail if it adds proof/replay constructors, a raw-hash-to-`PROVEN` path, or a `PROVEN` arm.

8. **Projection.** Run `./scripts/build_world_package.sh`; require canonical/projected
   `types.ail` byte equality and unchanged four-module allowlist/count.

9. **Ready packet.** Recompute and canonically encode the packet through `host/pkgproj`; require
   content/interface/tar hashes to differ from the old golden, four exports, and six tar entries.
   Fail on hand-authored JSON, unchanged interface hash, or packet drift.

10. **Full pinned gate.** Run
    `AILANG_BIN=/tmp/ailang-v0300/ailang ./scripts/verify_ail.sh`; require rc=0, the five named
    verified identities including `gradeOf`, 20 named tests including six `gradeCode` identities,
    and package 9/9 with non-zero work.

11. **Go compatibility.** Run the pinned `./scripts/verify_go.sh`; fail on the existing
    `RecordedEffect` fixture or any Go regression. No production Go source change is expected
    because no `Evidence` carrier or host boundary is added (V17, V21). **Base condition, not a
    regression:** the rig's local `go` is newer than the module's pin, so this script FATALs
    unless invoked with `GOTOOLCHAIN=go1.25.6`. A FATAL without that variable is the instrument,
    not the change — set it before reading the result.

12. **Scope.** The implementation changes only the four metadata files unless measurement finds
    the existing test manifest elsewhere. Fail on renderer, host codec, replay/proof execution,
    package manifest, proposal confidence, new module, or new evidence constructor changes.

## 8. Conflict surface

### 8.1 Primary AILANG gate

- `scripts/verify_ail.sh:135-147`: eleven Leg-1 module identities remain unchanged.
- `scripts/verify_ail.sh:262-267`: `world/types.ail` moves from an empty required-proof set to
  `{"gradeOf"}`.
- `scripts/verify_ail.sh:278-282`: JSON verifier errors/counterexamples remain authoritative.
- `scripts/verify_ail.sh:310-314`: exact verified total moves 4 → 5.
- `scripts/verify_ail.sh:317-340`: six `gradeCode` identities are added and total moves 14 → 20
  after observing actual runner JSON (V5, V24).

### 8.2 Package projection and golden

- `scripts/build_world_package.sh:7-17,53-71`: the four-module allowlist/count and wholesale
  staging projection remain unchanged; the script regenerates the `types.ail` copy.
- `packages/world-core/world/types.ail`: gains the grade ADT/function/tests but retains the
  five-constructor `Evidence` shape and stays byte-identical to canonical.
- `scripts/verify_world_package.sh:33-34,86-128`: four modules/exports, exact `.ail` allowlist,
  SHA equality, frozen manifest, and empty effects remain fixed.
- `scripts/verify_world_package.sh:190-245`: six tar entries remain fixed; the canonical ready
  packet and golden byte comparison derive new hashes and size.
- `scripts/world_package_ready_packet.golden.json`: regenerated in the same implementation; the
  design does not guess its new byte length (V6, V7).

A separate `world/grades.ail` would move all fixed inventory above. Co-location avoids that
projection/manifest expansion while still changing interface-derived hashes.

### 8.3 Evidence construction and consumption

AILANG `Evidence` is currently write-only and never written non-empty. `verify(w,p)` calls only
`proposalMatchesWorld`; production proposal constructors pass `[]`; contract fixtures also use
`[]`. A same-scope `.stateRoot` control fires (V20).

Production Go contains no `Evidence` constructor or decoder. The same search finds 13 test hits,
including the test-local `RecordedEffect` encoding, proving the instrument reaches the intended
scope (V21). Consequently this item changes no host codec and wires no producer.

These absences do not license future unchecked minting. They explain why validation plus producer
integration is a separate item rather than a hidden adapter task. The follow-on must inventory
the then-current paths again; this measurement is a baseline, not permanent authority.

### 8.4 Proof and replay producers

There is no production Z3 proof-report producer: production Z3 references are comments about Go
predicates mirroring proved AILANG predicates, while the same production scope has 424 `hashref`
hits (V22). Thus “wire the existing proof producer” has no target at this revision.

Replay is real. `host/replay/replay.go` defines `DivergenceError`, compares produced data, and
emits `archive.KindHashMismatch` (V8, V23). What is absent is an evidence/report boundary after a
successful result. This item neither promotes `RecordedEffect` nor treats replay machinery's
existence as proof evidence.

### 8.5 Renderer and prose consequences

Item 14 remains unbuilt. Before the validated producer boundary lands, a renderer may show the
mapped grades, but must not display `PROVEN` for kernel evidence and has no kernel
constructor for an unobtainable value.

HUMAN-SURFACE's “no total mapping” statement becomes stale after implementation, while its lack
of a `PROVEN` carrier remains true (V1). Mission prose describing the five-variant mismatch must
be updated narrowly: the mapping gap closes, the carrier gap does not (V18). Historical
implemented design documents remain historical and are not silently rewritten.

## 9. Scope and pricing

Representation-only implementation arithmetic:

| Work | Time |
|---|---:|
| Add grade ADT, exact contract/function, and six cases | 0.15 d |
| Move named proof/test pins; run RED/control mutations | 0.20 d |
| Rebuild package projection and canonical golden | 0.15 d |
| Run pinned AILANG/Go gates and inspect implementation | 0.15 d |
| **Total** | **0.65 d** |

This fits the queued ~0.5–1 day band. The original carrier design's 0.75-day estimate omitted the
first evidence boundary and first Z3 report producer. A boundary with bounded loading, hash
verification, typed codecs, success results, replay integration, proof production, and six
negative classes is plausibly a separate multi-day item and must be measured before scheduling.
Silently including it here would make the queue estimate false.

The representation-only change is atomic: grade type, total proof, named runtime cases, package
projection, and golden move together. The producer boundary is also logically atomic with any
future top-grade carriers; neither carrier may land ahead of that enforcement.

## 10. What this item is NOT doing

- It does not add `ProofReport`, `ReplayReport`, or any route to `PROVEN`.
- It does not build or change item 14's read-only workbench.
- It does not grade `Proposal.confidence`.
- It does not aggregate `Proposal.evidence` or define empty-list/minimum-grade policy.
- It does not load, validate, decode, or repair a `HashRef`.
- It does not create an evidence codec or production Go constructor.
- It does not create a Z3 proof-report producer or replay-report object.
- It does not alter replay, proof execution, proposal verification, commit policy, or effects.
- It does not add a module, package export, tar member, Go API, renderer API, or effect.
- It does not publish `world/core` or change its version.
- It does not claim constructor-directed grading validates referents.

## 11. Verification Log

Commands V1–V19 were run from repository root at `2ef2271` on 2026-08-13. V20–V24 are controller
measurements at the same HEAD; they are attributed rather than represented as this author's runs.
Probe commands use `AILANG_RELAX_MODULES=1` only for `.probe/` paths. Semantic verdicts come from
JSON rather than process rc.

| ID | Claim | Exact command | Observed output |
|---|---|---|---|
| V1 | Ratified surface defines four ordered grades, `UNSUPPORTED`, partial mappings, absent proof carriers, and confidence/unreadable-link rulings. | `nl -ba design_docs/HUMAN-SURFACE.md \| sed -n '67,78p;193,204p;268,307p;316p'` | `PROVEN > TESTED > ATTESTED > CLAIMED`; no total mapping; five variants; three mappings; no `PROVEN` carrier; confidence unsupported; preserve unreadable state. |
| V2 | Canonical `Evidence` has five constructors; `Proposal` contains `list[Evidence]` and bare confidence. | `nl -ba world/types.ail \| sed -n '20,52p'` | `23-28` five constructors; `50 evidence: list[Evidence]`; `51 confidence: float`. |
| V3 | Grade words have no non-comment implementation; `Evidence` control fires. | `test -d world && test -d host; printf 'grade='; rg -n 'PROVEN\|TESTED\|ATTESTED\|CLAIMED' world host --glob '*.{ail,go}' \| rg -v '^.*:[0-9]+:\s*(--\|//)' \| wc -l; printf 'control='; rg -n 'Evidence' world host --glob '*.{ail,go}' \| wc -l` | roots rc=0; `grade=0`; `control=20`. |
| V4 | `world/types.ail` has no functions/tests/contracts; declaration and sibling-behaviour controls fire. | `printf 'types_behaviour='; rg -n 'export func\|^tests \[\|requires \{|ensures \{' world/types.ail \| wc -l; printf 'types_decl_control='; rg -n '^export type' world/types.ail \| wc -l; printf 'world_behaviour_control='; rg -n 'export func\|^tests \[\|requires \{|ensures \{' world --glob '*.ail' \| wc -l` | `types_behaviour=0`; `types_decl_control=7`; `world_behaviour_control=26`. |
| V5 | Gate pins eleven modules, four proofs, and fourteen named tests and parses verifier failures. | `nl -ba scripts/verify_ail.sh \| sed -n '132,147p;258,315p;317,345p'` | existing `world/types.ail`; empty required set; JSON failure checks; totals 4 and 14. |
| V6 | Package build/verify pins four modules/exports, projection equality, six tar entries, and golden comparison. | `nl -ba scripts/build_world_package.sh \| sed -n '5,72p'; nl -ba scripts/verify_world_package.sh \| sed -n '32,35p;86,128p;190,245p'` | four-module staged copy; exact allowlists; SHA equality; four exports; six entries; golden compare. |
| V7 | Canonical/projected types match, sibling control differs, and pinned baseline gate is green. | `sha256sum world/types.ail packages/world-core/world/types.ail world/contracts.ail packages/world-core/world/contracts.ail; cmp -s world/types.ail packages/world-core/world/types.ail; echo types=$?; cmp -s world/types.ail world/contracts.ail; echo control=$?; /tmp/ailang-v0300/ailang --version; AILANG_BIN=/tmp/ailang-v0300/ailang ./scripts/verify_ail.sh \| tail -8` | types hashes equal; contracts pair equal; `types=0`, `control=1`; v0.30.0; package `9/9`; 4 proofs, 14 tests. |
| V8 | Replay exists and deliberate divergence fails. | `nl -ba host/replay/replay.go \| sed -n '1,29p;62,80p'; nl -ba host/replay/replay_test.go \| sed -n '168,264p'` | replay compares bytes/world hashes; named match and divergence tests require equality/typed errors. |
| V9 | A bare-ADT grade function with an exact match-equality contract verifies. | `AILANG_RELAX_MODULES=1 /tmp/ailang-v0300/ailang ai-check .probe/adt_grade.ail > .probe/adt_grade.ail.json; python3 -c 'import json; d=json.load(open(".probe/adt_grade.ail.json")); print(d["check"],d["verify"])'` | check true; verified 1; counterexample 0; errors 0; `gradeOf` verified. |
| V10 | False postcondition gives a constructor counterexample. | `AILANG_RELAX_MODULES=1 /tmp/ailang-v0300/ailang ai-check .probe/adt_false.ail > .probe/adt_false.ail.json; python3 -c 'import json; d=json.load(open(".probe/adt_false.ail.json")); print(d["check"],d["verify"])'` | counterexample 1; model contains `CompilerOutput`; non-vacuous result expression. |
| V11 | v0.30.0 checker accepts a non-exhaustive ADT match. | `AILANG_RELAX_MODULES=1 /tmp/ailang-v0300/ailang ai-check .probe/non_exhaustive.ail > .probe/non_exhaustive.ail.json; python3 -c 'import json; d=json.load(open(".probe/non_exhaustive.ail.json")); print(d["check"],d["verify"])'` | three constructors, two arms; check true; zero type errors. |
| V12 | Contract makes missing arm red; exhaustive control verifies. | `for f in .probe/non_exhaustive_contract.ail .probe/exhaustive_contract.ail; do AILANG_RELAX_MODULES=1 /tmp/ailang-v0300/ailang ai-check "$f" > "$f.json"; python3 -c 'import json,sys; d=json.load(open(sys.argv[1])); print(sys.argv[1],d["check"],d["verify"])' "$f.json"; done` | missing arm: verifier error `non-exhaustive pattern match`, model `C`; exhaustive: verified 1, errors 0. |
| V13 | Record containing `list[ADT]` errors while flat-record control verifies. | `AILANG_RELAX_MODULES=1 /tmp/ailang-v0300/ailang ai-check .probe/record_adt.ail > .probe/record_adt.ail.json; python3 -c 'import json; d=json.load(open(".probe/record_adt.ail.json")); print(d["check"],d["verify"])'` | `bagLabelled` unknown sort; `flatLabelled` verified; errors 1. |
| V14 | Repository comments record Proposal/ADT verifier limitation and silent rc. | `nl -ba world/contracts.ail \| sed -n '10,16p;41,45p;67,72p;102,120p'` | Proposal evidence yields unknown sort while rc=0; record-free predicate verifies. |
| V15 | Library assigns evidence to `world/types`, predicates to `world/contracts`, and uses `HashRef`. | `nl -ba design_docs/implemented/w-world-library-m1.md \| sed -n '69,110p'` | responsibility table and typed surface state those facts. |
| V16 | Golden verifier exists but dedicated regeneration command does not; same-scope control fires. | `test -d scripts; printf 'regen='; rg -n 'regenerat.*world_package_ready_packet\|write.*world_package_ready_packet' scripts \| wc -l; printf 'control='; rg -n 'world_package_ready_packet' scripts \| wc -l` | `regen=0`; `control=2`. |
| V17 | Constructor consumers are declarations and one Go test fixture; non-test Go consumer is absent with test control. | `test -d world && test -d host && test -d cmd && test -d packages; rg -n 'CompilerOutput\|TestReport\|HumanApproval\|AiReview\|RecordedEffect' world host cmd packages --glob '*.{ail,go}'; printf 'non_test='; rg -l 'CompilerOutput\|TestReport\|HumanApproval\|AiReview\|RecordedEffect' host cmd --glob '*.go' --glob '!*_test.go' \| wc -l; printf 'test_control='; rg -l 'CompilerOutput\|TestReport\|HumanApproval\|AiReview\|RecordedEffect' host cmd --glob '*_test.go' \| wc -l` | declarations plus `episode_test.go`; `non_test=0`; `test_control=1`. |
| V18 | Mission prose repeats five-variant and unreachable-`PROVEN` premise. | `rg -n 'five.*variants\|PROVEN.*unreachable\|no.*Evidence.*carrier' design_docs/world-mission.md` | hits at `2247-2253`, `2781-2786`. |
| V19 | Direct ADT expectations fail while private integer adapter works; emitted identity uses adapter name. | `AILANG_RELAX_MODULES=1 /tmp/ailang-v0300/ailang ai-check .probe/exact_grade.ail > .probe/exact_grade.json; AILANG_RELAX_MODULES=1 /tmp/ailang-v0300/ailang test --format json .probe/exact_grade.ail > .probe/exact_test.json; AILANG_RELAX_MODULES=1 /tmp/ailang-v0300/ailang ai-check .probe/exact_contract.ail > .probe/exact_contract.json; AILANG_RELAX_MODULES=1 /tmp/ailang-v0300/ailang ai-check .probe/code_tests.ail > .probe/code_ai.json; AILANG_RELAX_MODULES=1 /tmp/ailang-v0300/ailang test --format json .probe/code_tests.ail > .probe/code_test.json` | exact contract verified; direct ADT expected identifiers rejected; private adapter cases emitted as `gradeCode_test_*` and passed. |
| V20 | **Controller measurement:** AILANG evidence is write-only and all proposal construction is empty. | `rg -n '\.evidence\|evidence:' world packages design_docs/sketches --glob '*.ail'; printf 'control='; rg -n '\.stateRoot' world packages design_docs/sketches --glob '*.ail' \| wc -l; nl -ba world/transitions.ail \| sed -n '33,53p;87,94p'; rg -n 'evidence: \[\]' world/contracts.ail` | only declarations, unrelated `list[HashRef]`, empty `evidence: []` constructors/tests; `verify` does not read evidence; `.stateRoot` control 6 hits. |
| V21 | **Controller measurement:** production Go has no `Evidence` constructor/decoder. | `printf 'production='; grep -rn Evidence --include='*.go' --exclude='*_test.go' host/ cmd/ \| wc -l; printf 'test_control='; grep -rn Evidence --include='*_test.go' host/ cmd/ \| wc -l` | `production=0`; same-scope `test_control=13`, including test-local `RecordedEffect`. |
| V22 | **Controller measurement:** no Z3 proof-report producer exists. | `printf 'z3='; rg -n 'Z3' host cmd --glob '*.go' --glob '!*_test.go'; printf 'hashref_control='; rg -ni 'hashref' host cmd --glob '*.go' --glob '!*_test.go' \| wc -l` | every Z3 hit is a comment about mirrored predicates; `hashref_control=424`. |
| V23 | **Controller measurement:** replay has typed divergence/hash-mismatch behaviour but no Evidence output path. | `nl -ba host/replay/replay.go \| sed -n '105,117p;279,289p'; printf 'evidence_control='; rg -n 'Evidence' host/replay --glob '*.go' \| wc -l; printf 'replay_control='; rg -n 'DivergenceError\|KindHashMismatch' host/replay --glob '*.go' \| wc -l` | `DivergenceError` at 111; emits `archive.KindHashMismatch` at 285; no replay path constructs `Evidence`; replay control fires. |
| V24 | **Controller measurement:** private inline tests are collected and exported contract still verifies with private caller. | Pinned v0.30.0 module with exported contracted `gradeOf`, private tested `gradeCode`; `ailang test --format json <dir>/` plus verifier JSON. | `gradeCode_test_1`, `gradeCode_test_2` pass; `passed=2 failed=0 success=True`, rc=0; `gradeOf verified`, `verify.errors=0`. |
| V25 | **Controller measurement — §3.2's LITERAL code run on the real tree, not a mock, both before and after the round-2 carve-out.** An isolated copy of `world/*.ail` (never the live checkout) with the real `world/logepoch (HashRef)` import; baseline control taken first. | `cp world/*.ail <iso>/world/`; `ailang test --format json world/` (baseline); apply §3.2 verbatim; `ailang ai-check -timeout 5s world/types.ail`; `ailang test --format json world/`; then the carve-out edit; then mutation M1; `shasum -a 256` before/after each. | **Baseline control: `passed=14 failed=0`, reproducing the live gate's `EXACT_TOTAL_TESTS=14`.** After §3.2: `gradeOf` **`verified`**, `verify.verified=1`, `errors=0` (so `EXACT_TOTAL_VERIFIED` 4→5 is measured); `passed=20 failed=0`, identities exactly `gradeCode_test_1`–`gradeCode_test_6` (so 14→20 is measured, not guessed). **M1** (sixth `Evidence` variant, no arm) LANDED (sha `6ad0c8e7…`→`d894b46c…`) → `verify.errors=1`, `gradeOf status=error`, `non-exhaustive pattern match`, while **`check.passed` stayed `true` and rc stayed `0`**; restored byte-identical. After the carve-out (sha `6ad0c8e7…`→`2cf5b004…`): `gradeOf verified`/`errors=0`, `passed=20 failed=0`, and **M1 still reds** (`errors=1`); restored byte-identical. |
| V26 | **Controller measurement:** a probe of this fix that did NOT land still printed a plausible green, and was discarded rather than banked. | First carve-out attempt asserted `'UNSUPPORTED' not in src`, which fired on a COMMENT, so the file was never written — yet the subsequent Leg-1/Leg-2 commands printed `gradeOf verified` and `passed=20`, indistinguishable from a real post-fix result. | The green was the PRE-fix file measured twice. Only a sha-based landed-proof (`6ad0c8e7…`→`2cf5b004…`) makes V25's carve-out arm evidence. Recorded because a mutation that never ran and a mutation that did not red share an exit code. |

## 12. Quorum and revision-round decision record

Revision round 1 was **BLOCKED by both reviewers**.

- `gemini-3-1-pro` correctly objected that tests attached to private `gradeCode` emit
  `gradeCode_test_N`, not `gradeOf_test_N`. This revision changes every manifest/acceptance
  reference to `gradeCode` identities and records the private-test collection and verifier
  compatibility measurements. It also corrects V19: the controller's confirming probe had two
  cases; the load-bearing result is the identity pattern, not an invented “three of three” count.
- `gpt5-6-sol` correctly objected that the proposed top-grade carriers were forgeable and that
  “producers must construct them only after success” was unenforced prose. This revision withdraws
  the claim that this item makes `PROVEN` honestly reachable, withdraws both carriers, and chooses
  representation-only total mapping. The objection cannot be answered by symmetry with weaker
  constructors because unchecked access to `PROVEN` has categorically higher consequence.

The controller's inventory also corrects the reviewer's assumed remedy: there are not two real
evidence producers waiting to be wired. There is no production Evidence boundary and no Z3 report
producer; replay exists but emits no Evidence (V20–V23). Therefore the validation/producer work is
not silently expanded into this 0.65-day item.

### Round 2: `gemini-3-1-pro` PASS, `gpt5-6-sol` REJECT — narrow-refinement carve-out APPLIED

Round 2 had both external reviewers present (`absent_reviewers` empty, no N−1 degrade; `$0.087085`).
`gemini-3-1-pro` **passed**, naming as its closest non-blocking concern the private adapter's
implicit integer scale and specifically the unreachable `UNSUPPORTED => 0` arm, with the fix "simply
remove `UNSUPPORTED => 0`". `gpt5-6-sol` **rejected** on the same element, more strongly: adding
`UNSUPPORTED` to the exported kernel `EvidenceGrade` is prospective API design in a frozen core when
`gradeOf` can never return it and no consumer exists, and it conflates grading with
evidence-acquisition failure.

The controller applied the **narrow-refinement carve-out** (both limbs met: the objection carries a
concrete reviewer-authored `proposed_fix`, and it disputes no design direction — the total mapping,
the Z3-proven contract, the placement, the five arms and their grades, `PROVEN`'s deferral, the
private adapter and every gate-pin move are all accepted; the cut is of a non-core sentinel
constructor). Two independent reviewers converging on the same element strengthened it. `gpt5-6-sol`'s
fix was taken **verbatim**: `EvidenceGrade` is now exactly `PROVEN | TESTED | ATTESTED | CLAIMED`,
the `UNSUPPORTED => 0` adapter arm is gone, and its prescribed replacement text (including the
`GradeReadResult = Graded(EvidenceGrade) | Unsupported(UnsupportedReason)` illustration) is in §2.3.
Acceptance criteria, renderer guidance (§2.4, §8.5) and the §4 proof discussion were updated with it,
and `gemini-3-1-pro`'s dead-arm concern is closed by the same edit.

**The carve-out fix was MEASURED before it was applied, not after** (V25). It costs nothing and
breaks nothing: on an isolated copy of the real `world/` tree, with the fix landed (sha
`6ad0c8e7…` → `2cf5b004…`), `gradeOf` still reports `verified` with `verify.errors=0`, Leg 2 still
reports `passed=20 failed=0` with identities `gradeCode_test_1`–`gradeCode_test_6`, and mutation
**M1 still reds** (`verify.errors=1`, `non-exhaustive pattern match`, `check.passed` still `true`
and rc still `0`), restored byte-identical. So the totality mechanism survives the carve-out.

### Follow-on proposal: `w-validated-proven-evidence-boundary`

Price only after a fresh inventory; treat it as a multi-day producer-and-boundary item, not a
0.5–1 day type ratification. It must land atomically with any proof/replay carriers and:

1. define explicit mint authority and an opaque/validated value unavailable to proposal authors;
2. bounded-load the `HashRef`, verify its hash, decode a typed report, and require an explicit
   successful proof or replay result;
3. build the first production Evidence encode/decode boundary and the missing Z3 report producer;
4. integrate successful replay without treating `RecordedEffect` as replay proof;
5. return explicit error/`UNSUPPORTED` on every validation failure, with no fallback grade;
6. prove by named mutations that **arbitrary, missing, malformed, mismatched, failed, and
   divergent** reports cannot yield `PROVEN`;
7. permit renderer display of `PROVEN` only for the validated value produced by that boundary.

## Related

- **Proposed queue row: `w-validated-proven-evidence-boundary`.** Multi-day, re-measure before
  scheduling: define mint authority; add bounded loading, hash verification, typed report decode,
  explicit proof/replay success, the first production Evidence boundary and producer wiring; and
  prove that arbitrary, missing, malformed, mismatched, failed, and divergent reports cannot
  yield `PROVEN`. Top-grade carriers and renderer permission land only with this boundary.
- `design_docs/HUMAN-SURFACE.md` — binding trust gradient and total-mapping ratification point.
- `design_docs/coding-standards.md` — S1 identities, S3 package placement, S6 honest gates.
- `design_docs/DESIGN.md` §1 and §14 — kernel and controlled self-modification boundary.
- `design_docs/implemented/w-world-library-m1.md` — module responsibility split.
- `design_docs/planned/w-ail-gate-module-pin.md` — Leg-1 identity policy.
- `world/types.ail` — canonical `Evidence` and `Proposal` types.
- `world/transitions.ail` — current proposal construction and verification.
- `host/replay/replay.go` — replay result and divergence machinery.
- `scripts/verify_ail.sh` — primary proof/test gate.
- `scripts/build_world_package.sh`, `scripts/verify_world_package.sh` — projection and packet gate.
