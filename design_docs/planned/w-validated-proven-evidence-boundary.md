# w-validated-proven-evidence-boundary — authority-bearing proof evidence

- **Status**: **DIRECTION RATIFIED (`D-WORLD-17`); round-7 revision applied (`D-WORLD-19` arm A, attended 2026-08-19); round-7b carve-out revision applied (both round-7 fixes, converged), pending re-quorum** — DESIGNED, not landed
- **Park reason**: the iteration-84 park is answered. Mark Edmondson attended and ratified option B
  on 2026-08-14: authenticate canonical proof reports with a single-host, host-held MAC key. This
  revision applies that direction. Round 4 further narrows tranche 1 to an explicitly
  non-production, caller-keyed library boundary and measures the verifier's record-with-ADT-field
  encoding failure. See §10.4. Round 5 applies the attended `D-WORLD-17` ratification of
  2026-08-17: every seal is bound to its minting validator, and the free `GradeOfValidated` is
  replaced by the `Validator.Resolve` method. See §10.5. Revision 2 closes round 5's two
  sustained completeness objections: the zero-value forge (mint validity, `ErrUnmintedAuthority`,
  M21) and the §3.4/§9 bounded-runner contradiction, resolved (a) — tranche 1 owns the runner
  (V35). See the §10.5 revision-2 note. Round 7 applies the attended `D-WORLD-19` arm-A ruling
  of 2026-08-19: both round-6 objections are closed with the reviewers' own fixes — the
  sum-style `ResolutionResult` refusal channel (6a) and the bounded `host/store` object read
  (6b) — and the iteration-87 `DecodeProposal` byte-bound note is discharged. See §10.7.
  Round 7b applies both round-7 objections under the narrow-refinement carve-out — the two
  reviewer fixes CONVERGE on one corrected store signature (`OpenObject(ctx, ref, maxBytes)`),
  the `GetObject`-mutation arm is deleted as unsatisfiable (V40), and the false "item 18 bounds
  the WAIT" sentence is corrected: a context parameter is not a bound; a deadline is (V41). See
  §10.8.
- **Item**: queue item 17, `w-validated-proven-evidence-boundary`
- **Filed**: 2026-08-14, iteration 84
- **Measurement base**: `bef0153` (rounds ≤ 4); `03c7892` (round 5); `52bc9ec` (round 7, V36–V38); `7806cac` (round 7b, V40–V42)
- **Instrument**: `/tmp/ailang-v0300/ailang`, AILANG v0.30.0 (`e37b370`)
- **This tranche estimate**: **4.25 days** (round 7 adds the ratified bounded store read; round
  7b adds the `ObjectReadTimeout` deadline machinery — priced, not absorbed, §9)
- **Decomposition**: **yes** — four ordered sprint-sized items, **10.25 days total** (§9)
- **Design result**: serialized proof references remain untrusted and grade `CLAIMED`; only a host
  validator may mint the sealed value from which only the minting validator's `Resolve` method
  returns authority-bearing `PROVEN`.

This document is the first tranche of the larger queue item. The measured inventory found no production
Go Evidence boundary, no Z3 report producer, no renderer, and no named Go-test manifest (V3, V4,
V14, V15). Implementing proof, replay, and rendering together would exceed the 3–4 day sprint
guardrail. Section 9 therefore decomposes the work without weakening the authority rule.

Every present-tense repository claim is tied to the numbered Verification Log in §11. The revised
AILANG syntax and its one-sided negative control were checked with the pinned release in `/tmp`
(V25); V18/V19/V22/V24/V26 remain explicitly labelled history of rejected designs or earlier
comparative probes.

---

## 1. Problem

`Evidence` has five public constructors and `gradeOf` maps none to `PROVEN` (V30); the six mapping
tests and the proof identity are already pinned by the AILANG gate (V5, V6). This is intentional. An
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
   return is a `ValidatedEvidence` whose fields and construction function are unexported, and it
   can succeed only after authenticating the report with the host-held key.
2. `host/evidence.DecodeProposal` decodes `ProofReceipt` as untrusted `ClaimedEvidence`; possession
   of the receipt neither seals it nor changes its kernel grade.
3. `Validator.Resolve(sealed ValidatedEvidence) ResolutionResult` — a METHOD on `Validator`, per
   the attended `D-WORLD-17` ratification, with the round-6 sum-style result applied under the
   attended `D-WORLD-19` arm A — is the only API that can produce host `ResolvedGradeProven`.
   `ResolutionResult` has unexported fields and exposes mutually exclusive
   `Proven() (ResolvedGrade, bool)` and `Err() error` accessors: success carries exactly
   `ResolvedGradeProven`; refusal carries exactly one of `ErrUnmintedAuthority` or
   `ErrForeignSeal` and no grade; and the zero value is an explicit refusal state, never
   success. There is no free resolver: the former package-level `GradeOfValidated` is DROPPED,
   with no deprecated alias, package-level helper, or convenience wrapper, because a resolver
   detached from a validator instance is the round-4 defect itself. `Resolve` accepts neither
   raw `Evidence`, decoded claims, nor a `HashRef`, and it refuses any `ValidatedEvidence`
   sealed under a different mint identity (`ErrForeignSeal`) — and any zero or unminted
   identity outright, before comparison (`ErrUnmintedAuthority`, §3.2). `ResolvedGrade` is
   represented in Go, not serialized in the kernel ADT.
4. Within the library API, a `ProofReceipt` is revalidated against its authenticated report,
   expected subject, and the caller-supplied key before a new sealed value can exist. No serialized
   value regains authority by decode. Production key loading and restart wiring are tranche 2.
5. Any renderer introduced later accepts `ValidatedEvidence` (or the read-only resolved view from
   it), never raw `Evidence`, raw JSON, kernel `EvidenceGrade`, or receipt bytes. Until tranche 4,
   no renderer may show `PROVEN`.

The safety statement is now exact at both language boundaries: an agent can spell
`ProofReceipt`, but direct `gradeOf` execution returns `CLAIMED`; an agent can construct only a
ZERO-VALUE Go `ValidatedEvidence` — exported types make `var s ValidatedEvidence` unpreventable —
which `Resolve` refuses by name before any comparison, and even holding a genuinely minted seal,
cannot make any validator carrying a different mint identity resolve it. Kernel `EvidenceGrade.PROVEN` remains a public,
agent-spellable enum value with no `Evidence -> EvidenceGrade` kernel producer in tranche 1. Its
mere spelling carries no authority, just as caller-written Go grade data carries none. That
reserved result is acceptable: removing it would break the already-pinned grade vocabulary,
while assigning it to a public Evidence constructor would recreate grade laundering.

The binding property is stated exactly, neither more nor less. What is enforced: a zero-value
`Validator` or `ValidatedEvidence` never resolves to any grade (`ErrUnmintedAuthority`), and you
CANNOT make a validator carrying someone else's mint identity resolve your seal — `Resolve`
refuses any `ValidatedEvidence` whose mint identity differs from the receiver's
(`ErrForeignSeal`). The guarantee is per-identity, not per-Go-variable: a value copy of a
validator carries the same key and the same identity and IS the same mint authority, deliberately
(§3.2). What is not enforced: you CAN still construct your own `Validator` with your own key via
`NewValidator` and self-mint into it, and that is accepted, because no library can stop a caller
lying to itself. Tranche 1 enforces provenance relative to a minting identity; it does not, and
cannot, make a caller-constructed validator's output mean anything to anyone else.

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
`TestPublicAuthoritySurfaceIsFrozen`. Tranche 4 adds the same rule at the first renderer.

The honest limitation is that no existing language mechanism prevents a future authorized source
change from adding a new public kernel producer or a new Go bypass. The gates make those changes
explicit and red-on-removal; they do not make repository maintainers powerless. Likewise, a
foreign module can define—or directly spell—a non-authoritative “proven” value, but it cannot
obtain host `ResolvedGradeProven` or make canonical `gradeOf(ProofReceipt(...))` return `PROVEN`
without a gated kernel change. Any future security-sensitive consumer that accepts a bare kernel
`EvidenceGrade` is outside this authority API and must be rejected by the API-freeze gate.

### 2.4 Why a host-held MAC, but not a bearer marker or an allowlisted hash?

The pre-ratification design rejected signatures because they move authority to key custody. The
human decision deliberately changes that position: the missing property is provenance, and a MAC
is the appropriate provenance primitive for this single-host boundary. `ValidateProof` verifies a
host-issued tag before minting. HMAC-SHA-256 needs neither a cross-host receipt nor a public-key
deployment, while re-executing the compiler and solver on every grade resolution would add an
effectful, timeout-prone critical-path operation. Secure production key custody is required before
the boundary is deployable, but is deliberately assigned to tranche 2 rather than implied by the
tranche-1 library API.

- A secret bearer marker is still rejected: a fixed marker serialized beside the report can be
  copied to unrelated bytes. The MAC instead binds a tag to the exact canonical report bytes.
- An allowlisted hash is mutable ambient state and still does not prove type, subject, or success.
- An exported Go struct with `Grade: PROVEN` is forgeable by any caller and is rejected.

The sealed value remains process-local authority. Durable authority is reconstructed only by full
validation, including MAC verification; the serialized receipt and envelope remain untrusted
inputs. Supplying a different key makes every report tagged under the prior key unvalidatable.
There is no correctness-critical rotation protocol: affected reports must be regenerated by
rerunning the bounded producer, which is the same cost option A would pay on every resolution.
Production loss, replacement, file custody, and restart behavior are tranche-2 concerns.

### 2.5 Explicit failure semantics

Every validation call returns exactly one of:

```text
Validated(ValidatedEvidence) | Unsupported(UnsupportedReason)
```

The reasons are stable identifiers: `invalid_ref`, `missing`, `oversize`, `hash_mismatch`,
`wrong_semantic_id`, `wrong_interface`, `malformed`, `unauthenticated_report`,
`subject_mismatch`, `tool_mismatch`, `proof_failed`, and `proof_incomplete`. Store/I/O cancellation
or deadline expiry — including expiry of the validator's own derived `ObjectReadTimeout` deadline
(§3.2, AC16) — is an explicit operational error that mints no seal, not a grade. Neither
`Unsupported` nor error contains an `EvidenceGrade`; no failure falls back to `CLAIMED`,
`ATTESTED`, or `TESTED`.

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
code. The `HashRef` identifies the canonical authenticated proof-report envelope defined in §3.3.
It is a receipt pointer, not a capability, and its kernel grade never depends on whether the
referenced envelope happens to exist.

Do not add `ProofReport` fields to AILANG. Report parsing, byte limits, solver metadata, and object
loading are effect-boundary concerns. Do not add an aggregate function over `Proposal`: the pinned
verifier fails Z3 encoding when a function parameter is a record with an ADT-typed field, whether
the field is a bare ADT or `list[ADT]`. Scalar-only records and records with `list[scalar]` verify,
so neither records nor lists alone are the limitation; a bare ADT parameter also verifies. The
failure occurs even when the contract does not read the ADT-typed field, and `ai-check` still exits
0 with a clean `check` section while the affected `verify.results[].status` is `error` (V31–V32).

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
ObjectReader                 // minimal read seam the validator consumes; its one method is the reconciled round-7b signature (§8.2, §10.8): OpenObject(ctx context.Context, ref hashref.HashRef, maxBytes int64) (ObjectMeta, io.ReadCloser, error)
NewValidator(key [32]byte, reader, compilerConfig)   // configuration REQUIRES a positive ObjectReadTimeout; zero or negative is a constructor refusal (AC16)
Validator.ValidateProof(ctx, reportRef, expectedSubject)   // derives context.WithTimeout(ctx, ObjectReadTimeout) before any open or read, even under context.Background() (§3.3 step 2, AC16)
DecodeProposal(raw)          // ProofReceipt remains an untrusted claim; raw capped at 256 KiB before any parse (§3.3)
Validator.Resolve(sealed ValidatedEvidence) ResolutionResult   // method; sole ResolvedGradeProven source; refuses zero/unminted identities with ErrUnmintedAuthority, foreign seals with ErrForeignSeal
ResolutionResult             // unexported fields; mutually exclusive Proven() (ResolvedGrade, bool) and Err() error; zero value is refusal
ResolvedGrade                // host enum; ResolvedGradeProven is not serialized
```

No `NewValidatedEvidence`, struct literal, `SetGrade`, raw-hash grade resolver, receipt resolver,
exported unseal method, or package-level grade resolver of any name exists — the free
`GradeOfValidated` is dropped, not renamed. Put external-package tests in
`host/evidence/authority_test.go` so they compile with only the public API.

Binding mechanism, stated as an invariant the reader can check, not an assurance. The mint
identity is an unexported pointer to a per-instance, non-zero-size heap allocation made only
inside `NewValidator` — non-zero-size because Go permits two distinct zero-size allocations to
share an address, which would let two identities collide by accident; pointing it at the
validator's own heap-held key copy satisfies this. `ValidateProof` stamps that pointer into the
seal's unexported `mintedBy` field at mint time. `Resolve` performs two ordered checks before
returning any grade:

1. **Mint validity.** If the receiver validator's identity or the seal's `mintedBy` is nil,
   `Resolve` returns a `ResolutionResult` whose `Err()` is the dedicated sentinel error
   `ErrUnmintedAuthority` and whose `Proven()` reports no grade, BEFORE any comparison. A zero
   identity is never valid. This closes the zero-value forge:
   `Validator` and `ValidatedEvidence` are exported types, so a foreign package can always
   construct `var v evidence.Validator; var s evidence.ValidatedEvidence` — unexported fields
   stop writes, not zero values — and a bare equality check would compare zero to zero and pass.
   Here that pair is refused at step 1 and never reaches the comparison.
2. **Binding.** Otherwise, if `sealed.mintedBy` differs from the receiver's identity, `Resolve`
   returns a `ResolutionResult` whose `Err()` is the dedicated sentinel error `ErrForeignSeal`
   and whose `Proven()` reports no grade.

The checkable property: (a) both identity fields are unexported, so no foreign package can
assign them; (b) the only assignments in the package are inside `NewValidator` and
`ValidateProof`; (c) the Go zero value of either type has a nil identity and is refused at step
1; (d) `ResolutionResult`'s own fields are unexported and its zero value reports `Proven()`
false — an explicit refusal — so no code path, including an uninitialized return, can express
success without the package constructing it. Therefore every `ResolutionResult` whose `Proven()`
reports a grade requires two non-nil, equal pointers, and such
a pair exists only in values descended from the same `NewValidator` call. `ErrUnmintedAuthority`
and `ErrForeignSeal` are each produced by exactly one check and nowhere else, and neither is a
member of §2.5's `UnsupportedReason` set, which belongs to `ValidateProof`'s validation
pipeline — so each observable identifies its check uniquely rather than reading a shared error
enum or a bare absence.

Copy semantics are explicit, and the guarantee is worded no stronger than the mechanism. A value
copy of a `Validator` (`v2 := *v1`, assignment, or embedding) carries the same identity pointer
and the same key: a copy IS the same mint authority, deliberately, because the identity names the
`NewValidator` call, not the Go variable. A copied seal likewise carries its `mintedBy` pointer
and remains resolvable by that same authority. What a copy cannot do is the load-bearing part: a
foreign package cannot conjure a value carrying a given identity without copying a value that
already holds it, since the pointer target is allocated only inside `NewValidator` and the fields
are unexported. The enforced guarantee is therefore per-identity: only a validator carrying the
identity minted by the `NewValidator` call that produced a seal can resolve that seal.

The new package depends on `host/hashref` and a minimal object-reader interface, not on daemon,
replay, or renderer. That avoids a cycle and makes the validator testable with a bounded fake.
`NewValidator` copies an exactly 32-byte MAC key supplied directly by its caller; tranche 1 accepts
no key path and neither creates nor loads key files. The proof producer receives the same
caller-owned key material through its library constructor. Consequently tranche 1 is explicitly
**library-only and non-production**: it demonstrates the authenticated authority boundary but does
not wire a production composition root, provision durable key material, or claim startup/restart
behavior. Those obligations belong to the named tranche-2 wiring item in §9.

### 3.3 Canonical `ProofReportV1`

The producer encodes strict canonical JSON `ProofReportV1` bytes with these fields in this order:

```text
schema             = "world/proof-report/v1"
subject            HashRef of the exact checked AILANG module bytes
compiler           HashRef of the exact executable bytes
compilerVersion    = "AILANG v0.30.0"
verified           sorted, non-empty list of verify.results[].function values whose status is "verified"
errors             = 0
counterexamples    = 0
checkPassed        = true
proofSucceeded     = true
```

`SemanticID` is `world/proof-report/v1`; a fixed `InterfaceHashV1` identifies the authenticated
envelope schema. The receipt's `HashRef` addresses canonical envelope bytes containing exactly two
fields in order: `report`, the base64url-no-padding encoding of the canonical `ProofReportV1`
bytes, and `mac`, the base64url-no-padding 32-byte tag. The tag is therefore carried beside—not
inside—the bytes it authenticates. HMAC-SHA-256 with the 32-byte host key covers exactly the
decoded canonical `ProofReportV1` byte sequence, with no whitespace normalization or reconstructed
struct accepted as a substitute.

Unknown, duplicate, missing, non-canonical, trailing, invalid-UTF-8, and over-limit envelope or
report input is rejected. The sole classification exception is the envelope's `mac` member:
missing, malformed, or wrong-length `mac` reaches the authentication step so all tag failures have
the stable `unauthenticated_report` reason. The raw envelope is capped at **256 KiB** at the
point of I/O: step 2 below streams the object payload from the store's bounded read through an
`io.LimitReader` sized 256 KiB + 1 bytes, so the FIRST copy is bounded and an attacker-supplied
`HashRef` to a multi-gigabyte object can never materialize more than the cap plus one detection
byte (round 6's objection 6b, applied verbatim under the attended `D-WORLD-19` arm A — the
retired wording "before allocating a second full copy" conceded exactly the unbounded first
copy). The decoded report bytes and verified list are separately capped at 256 KiB and 256
unique identities, and every string at 1 KiB. `DecodeProposal` applies the same 256 KiB cap to
its raw input before any parse, refusing oversize input undecoded — the iteration-87
non-blocking note, discharged. The chosen raw cap follows an existing 256 KiB strict-codec
bound already used for registry schema input (V11); it is a new Evidence invariant, not
inherited store behaviour.

Validation order is fixed:

1. validate canonical `HashRef`;
2. derive the read deadline — `readCtx` from `context.WithTimeout(ctx, ObjectReadTimeout)` —
   even when the caller supplies no deadline, `context.Background()` included (§3.2, AC16);
   then open exactly one object through `ObjectReader`'s
   `OpenObject(readCtx, ref, maxBytes)` (§8.2) and stream it through an `io.LimitReader` sized
   256 KiB + 1 — the payload is never fully materialized before the bound, and deadline expiry
   or cancellation is an explicit operational error that mints no seal (§2.5);
3. reject absent objects, and reject as `oversize` any read that fills the limit's extra
   detection byte;
4. recompute the envelope payload hash and compare algorithm plus digest to the requested ref;
5. require exact semantic ID and interface hash;
6. strict-decode the envelope and report, then require canonical re-encode byte equality for each,
   deferring only `mac` presence/shape to step 7;
7. recompute HMAC-SHA-256 over the exact decoded canonical report bytes and compare the 32-byte tag
   in constant time; absent, malformed, or unequal tags return `unauthenticated_report`;
8. require report subject equals the expected subject;
9. require compiler hash/version equals the configured pinned executable;
10. require the report's `verified` field to be sorted and non-empty and to contain the configured
   required identity set, plus `checkPassed`, `proofSucceeded`, zero errors, and zero
   counterexamples; §3.4 binds producer emission of that field to `verify.results[]`;
11. only then construct `ValidatedEvidence`.

One landing between rounds 6 and 7 must not be misread as having closed this — in either
dimension. Item 18 (`w-daemon-read-cancellation`) landed and threaded `context.Context` through
`GetObject` (V36). That bounds the SIGNATURE — not the bytes, and not, by itself, the wait. A
context caps nothing unless someone arms it with a deadline, and item 18's own ratified deferral
DR-2 deliberately leaves 11 production store reads passing `context.Background()`, pinned
exactly by `TestNoNewDeadlineFreeStoreReads` (V41). A freshly context-threaded store therefore
still OOMs on an attacker-supplied ref (no byte bound) and still blocks indefinitely for any
caller that supplies no deadline (no inherited wait bound). The bounded read above closes the
bytes; the `ObjectReadTimeout`-derived deadline in step 2 closes the wait at this tranche's OWN
call site, inherited from no one. The 11 DR-2 sites are item 18's follow-on obligation and are
deliberately NOT fixed here (§8.2). Round 7 wrote "item 18 bounds the WAIT" into this paragraph;
that sentence was false and is superseded (§10.8, V41).

The validator never accepts a report merely because `store.PutObject` once checked its hash. It
recomputes at consumption so alternate readers, corruption, and test fakes cannot bypass the
integrity boundary. MAC verification follows strict decode because it needs the uniquely decoded
report byte string, but precedes all attacker-controlled semantic assertions: a correct hash proves
only integrity, while a correct tag proves that the host producer authenticated those exact bytes.

### 3.4 First Z3 report producer

Add `host/evidence/proof_producer.go`. It runs only the configured absolute `AILANG_BIN`, verifies
the executable bytes and exact `AILANG v0.30.0` token, invokes `ai-check` for the subject, bounds
wall time, stdout, and stderr, and parses JSON rather than trusting process rc. It may use
`verify.verified>0` as a cheap precondition, but constructs the report's identity set only from
`verify.results[].function` entries filtered by `status == "verified"`. It emits, MAC-tags, and
stores an authenticated envelope only when `check.passed=true`, `verify.errors=0`,
`verify.counterexample=0`, and that sorted set contains the configured required identity set.

**Tranche 1 owns the bounded runner** (resolution (a) of the round-5 §3.4/§9 contradiction),
because the bounded producer is tranche 1's named deliverable and cannot emit an authenticated
report without executing the pinned checker, so the bounding primitive belongs to the tranche
that executes. No exported bounded-subprocess helper exists for `host/evidence` to import: the
repository's three non-test subprocess sites each bound locally, and the one general-purpose
runner, `runBounded`, is unexported and broker-internal (V35). The producer therefore implements
its own bounding inside `host/evidence`, on the pattern V35 measures in `runBounded`: wall time
is bound by `context.WithTimeout` feeding `exec.CommandContext`, with the child in its own
process group and the whole group killed on cancellation so a forked grandchild cannot outlive
the deadline; both stdout and stderr are independently read through their own capped
readers sized limit+1 bytes (gemini-3-1-pro's round-8 NON-BLOCKING fix, applied verbatim: a
singular "capped reader" for a dual-stream capture leaves stderr's memory bound implicit, so an
attacker-controlled checker could spam stderr and OOM the host while the parsed stream stays
within its cap), so overflow on EITHER stream is detected rather than silently truncated; deadline overrun and output overflow are each producer refusals
that emit and store no report. The producer departs from `runBounded` in exactly one measured
property: `runBounded` merges stderr into stdout, while the producer captures the two streams
SEPARATELY and parses JSON only from stdout, because this producer's output is parsed and a
warning merged into the parsed stream voids the parse (V22 records that failure mode; V27
separates the streams for the same reason). Reusing broker's runner is rejected, not deferred:
it would add a `host/broker` dependency §3.2 excludes, and exporting it is a broker API change
outside this tranche. The tranche-1 arithmetic in §9 already prices this construction inside the
0.60-day producer row; no pricing change follows.

The producer takes required identities from its caller; an empty set is refused. The first
integration is a library API, not a daemon route or an agent tool. Therefore an agent cannot ask
the host to prove arbitrary source and immediately attach the result through a public network
surface. A later transition may call it only after separately ratifying that authority.

### 3.5 What crosses the AILANG/Go boundary

Across serialization, the kernel carries `ProofReceipt(reportRef)`, which always remains an
untrusted `ClaimedEvidence`. Trusted validation separately produces
`ValidatedEvidence{reportRef, subject, ...}` with unexported fields.
`Validator.Resolve(sealed)`, called on a validator carrying the minting identity, is the sole
bridge to host `ResolvedGradeProven` — its `ResolutionResult`'s `Proven()` accessor is the only
place that grade exists; there is no free resolver and no bridge that serializes authority back
into an AILANG constructor.

The pure kernel proves only the untrusted receipt-to-`CLAIMED` mapping. Go performs bounded loading,
digest verification, strict typed decode, MAC authentication, solver-success checks, and
provenance enforcement. The MAC tag crosses storage only in the authenticated envelope and never
becomes kernel authority. No contract takes `Proposal`, so its record parameter with an
ADT-bearing `evidence` field avoids the measured encoding failure rather than assuming that
ignoring the field would evade it (V31–V32).

### 3.6 Gate changes

AILANG changes move all five coupled surfaces in one implementation:

- `world/types.ail` and byte-identical `packages/world-core/world/types.ail`;
- current `EXACT_TOTAL_TESTS` 39 → observed 40 and `REQUIRED_TESTS` gains the emitted seventh
  `gradeCode` identity;
- current `EXACT_TOTAL_VERIFIED` remains 10, while `gradeOf` remains named in
  `REQUIRED_VERIFIED`;
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
   10. Removing its contract or identity reds Leg 1.
2. `gradeCode_test_7` (use the observed emitted name) is added to `REQUIRED_TESTS`, and
   `EXACT_TOTAL_TESTS` moves from 39 to 40. Changing the new arm or deleting the case reds Leg 2.
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

Every tranche-1 refusal branch has its own compiling neutering mutation. M6 is specified now but
lands with ordered tranche 3; until then replay has no route to `PROVEN`, which is the stronger
default. Every observable is returned by or after the mutated mechanism, never a sibling value
assigned beside it.

| ID / class | Exact file and neutered edit | Named check that fires | Downstream observable and predicted failure text |
|---|---|---|---|
| M1 **arbitrary** | `world/types.ail`: change only body arm `ProofReceipt(_) => CLAIMED` to `ProofReceipt(_) => PROVEN` | `gradeCode_test_7` plus `gradeOf` verification | Agent-authored receipt reaches canonical kernel `PROVEN`; runtime fails `got 4, want 1`, while the unchanged contract yields a counterexample. |
| M2 **invalid ref** | `host/evidence/validator.go`: change `if refErr != nil` to `if false && refErr != nil` | `TestInvalidProofRefIsRefused` | The returned reason changes from the unique `invalid_ref`; failure: `got missing; want invalid_ref`. |
| M3 **missing** | `host/evidence/validator.go`: change `if !ok` to `if false && !ok` | `TestMissingProofReportIsRefused` | The returned reason is not `missing`; failure: `got malformed; want missing`. |
| M4 **oversize** | `host/evidence/validator.go`: change the streamed-read overflow guard — the check that the `io.LimitReader` sized 256 KiB + 1 filled its detection byte (§3.3 step 3) — to `if false && overflowed` | `TestOversizeProofReportIsRefused`, feeding a > 256 KiB object through the real bounded store read | The returned reason is not `oversize`; under the mutant the payload proceeds TRUNCATED at the limit, so step 4's recomputed hash cannot match; failure: `got hash_mismatch; want oversize`. |
| M5 **hash integrity** | `host/evidence/validator.go`: change the recomputed-hash guard to `if false && !sameHash` | `TestPayloadHashMismatchIsRefused` | The returned reason is not `hash_mismatch`; failure: `got unauthenticated_report; want hash_mismatch`. |
| M6 **divergent** (tranche 3) | `host/replay/evidence.go`: change `if err != nil` after `ReplayEpisode` to `if false && err != nil`, allowing `*DivergenceError` to reach report creation | `TestDivergentReplayCannotResolveProven` | The test corrupts recorded result bytes, calls replay-evidence production then the common resolver, and unexpectedly gets `PROVEN`; failure: `divergent replay resolved PROVEN; want unsupported replay_divergent`. This is downstream of replay comparison and report validation. |
| M7 **wrong semantic ID** | `host/evidence/validator.go`: change semantic-ID guard to `if false && got != proofSemanticID` | `TestWrongSemanticIDIsRefused` | The unique result changes; failure: `got unauthenticated_report; want wrong_semantic_id`. |
| M8 **wrong interface** | `host/evidence/validator.go`: change interface-hash guard to `if false && got != proofInterfaceHash` | `TestWrongInterfaceIsRefused` | The unique result changes; failure: `got unauthenticated_report; want wrong_interface`. |
| M9 **malformed** | `host/evidence/report_codec.go`: change strict-decode/canonicality guard to `if false && err != nil` while returning a zero report | `TestMalformedProofReportIsRefused` | The returned reason changes from `malformed`; failure: `got unauthenticated_report; want malformed`. |
| M10 **MAC authentication** | `host/evidence/validator.go`: change the single tag guard to `if false && (len(tag) != sha256.Size \|\| !hmac.Equal(want, tag))` | `TestOtherwisePerfectReportWithoutMACIsUnauthenticated` and `TestOtherwisePerfectReportWithWrongMACIsUnauthenticated` | Hand-authored otherwise-perfect absent-tag and wrong-tag envelopes seal; failure: `report resolved PROVEN; want unauthenticated_report`. The mutant compiles and directly neuters the whole constant-time MAC verification branch. |
| M11 **subject mismatch** | `host/evidence/validator.go`: change subject guard to `if false && report.Subject != expectedSubject` | `TestMismatchedProofSubjectIsRefused` | Failure: `got tool_mismatch; want subject_mismatch`. |
| M12 **tool mismatch** | `host/evidence/validator.go`: change compiler hash/version guard to `if false && toolMismatch` | `TestMismatchedProofToolIsRefused` | Failure: `got proof_incomplete; want tool_mismatch`. |
| M13 **proof failed** | `host/evidence/validator.go`: change success/error/counterexample guard to `if false && proofFailed` | `TestFailedProofReportIsRefused` | Failure: `got proof_incomplete; want proof_failed`. |
| M14 **proof incomplete** | `host/evidence/validator.go`: change verified-identity guard to `if false && incomplete` | `TestIncompleteProofReportIsRefused` | An authenticated report missing a required identity seals; failure: `incomplete proof resolved PROVEN; want proof_incomplete`. |
| M15 producer false-green | `host/evidence/proof_producer.go`: change `if verify.Errors != 0` to `if false && verify.Errors != 0` | `TestProofProducerRefusesVerifierErrors` | Producer stores a report despite JSON errors; `producer emitted report with verify.errors=1`. |
| M16 seal bypass | `host/evidence/grade.go`: add a grade resolver accepting `HashRef`, `ProofReceipt`, raw `Evidence`, or kernel `EvidenceGrade` | `TestPublicAuthoritySurfaceIsFrozen` | External-package API inventory gains a forbidden resolver; `public authority surface exposes non-sealed PROVEN ingress`. |
| M17 named-manifest removal | `scripts/verify_go.sh`: remove one literal required name, leaving the test present | `TestEvidenceNamedManifestRejectsUnpinnedTest` in `host/verifygate/evidence_manifest_gate_test.go` | Isolated gate sees an extra observed test; `evidence test set differs from REQUIRED_EVIDENCE_TESTS`. |
| M18 projection drift | edit only `world/types.ail` | existing world-package step 3/9 | `projection hash mismatch: world/types.ail` (exact wording confirmed during implementation). |
| M19 stale ready packet | rebuild projection but retain old golden | existing world-package step 9/9 | `ready packet differs byte-for-byte from golden`. |
| M20 **cross-validator binding** | `host/evidence/validator.go`: change the seal-identity guard in `Resolve` to `if false && sealed.mintedBy != v.id` (exact field names fixed at implementation; the neutered comparison is the binding check itself) | `TestAttackerChosenValidatorCannotMintForHostAuthority` | A seal minted by validator 2 (attacker-constructed, attacker-keyed) is presented to validator 1's `Resolve`, which must return a `ResolutionResult` whose `Err()` is the dedicated `ErrForeignSeal` and whose `Proven()` reports no grade — an observable only the binding comparison produces, not a shared reason enum or a bare absence. Under the mutant the foreign seal resolves; failure: `foreign seal resolved ResolvedGradeProven; want ErrForeignSeal`. |
| M21 **zero-value mint validity** | `host/evidence/validator.go`: change the mint-validity guard in `Resolve` to `if false && (v.id == nil \|\| sealed.mintedBy == nil)` (exact field names fixed at implementation; the neutered check is the nil-identity refusal itself, leaving the M20 comparison intact) | `TestZeroValueForgeryCannotResolve` | The test constructs `var v evidence.Validator; var s evidence.ValidatedEvidence` from the external test package — the two-line forge exported types make unpreventable — and calls `v.Resolve(s)`, requiring a `ResolutionResult` whose `Err()` is the dedicated `ErrUnmintedAuthority` and whose `Proven()` reports no grade; the same test also constructs `var r evidence.ResolutionResult` and requires `Proven()` false, pinning the declared zero-value-is-refusal contract. Under the mutant the two nil identities reach the M20 comparison, compare equal (zero == zero), and the forge resolves; failure: `zero-value seal resolved ResolvedGradeProven; want ErrUnmintedAuthority`. |
| M22 **deadline derivation** | `host/evidence/validator.go`: replace the read-deadline derivation `context.WithTimeout(ctx, cfg.ObjectReadTimeout)` with `context.WithCancel(ctx)` (exact spelling fixed at implementation; the neutered mechanism is the deadline derivation itself — the mutant compiles, keeps a cancel, and simply never arms a deadline, which is exactly the caller-inherited state AC16 forbids) | `TestBlockingObjectReadReturnsWithinObjectReadTimeout` | The test configures `ObjectReadTimeout` at 50 ms, supplies a fake `ObjectReader` whose `OpenObject` returns a reader blocking until ITS ctx is done, calls `ValidateProof` under `context.Background()`, and waits on a done channel guarded by a test-side 20× timer so the mutant FAILS rather than hangs. Correct code: the derived deadline expires, the blocked read unblocks, `ValidateProof` returns the explicit operational error within the bound, and no seal exists. Under the mutant no deadline is ever armed, `context.Background()` is never done, and the guard fires; failure: `validation did not return within 20x ObjectReadTimeout under context.Background(); want explicit operational timeout error and no seal`. |

For M2–M14, the control first validates one good authenticated report for the same expected subject and observes
that the minting validator's `Resolve` returns `ResolvedGradeProven`. Thus a mutant cannot pass
merely because the test never reached minting. M1's control executes an agent-authored receipt and
observes `CLAIMED` before applying the one-sided mutation. For M6, the non-divergent control must
produce a replay report and resolve host `PROVEN` before the corrupted record is introduced. For
M20, the control is same-instance: validator 1 first mints its own seal and observes its own
`Resolve` return `ResolvedGradeProven`, proving the fixture reaches minting and resolution before
the foreign-seal refusal is asserted. M21 carries the same same-instance control inside its own
test: a `NewValidator`-built validator resolves its own minted seal to `ResolvedGradeProven`
before the zero-value pair is asserted, so the refusal cannot be satisfied vacuously by a
resolver that refuses everything. M22's control lives in the same test: the same validator
configuration first validates one good authenticated report through a PROMPT (non-blocking)
fake reader under the same `context.Background()` and resolves it to `ResolvedGradeProven`,
proving the derived deadline does not refuse the fast path before the blocking arm asserts the
timeout refusal.

## 7. Acceptance criteria

1. **Authority surface.** `ValidatedEvidence` has no exported fields or public constructor.
   `Validator.ValidateProof` is the sole mint, and the `Validator.Resolve` method is the sole
   grade resolver — no package-level resolver exists under any name, and
   `TestPublicAuthoritySurfaceIsFrozen` reds if one appears. `Resolve` returns the sum-style
   `ResolutionResult`: unexported fields, mutually exclusive `Proven() (ResolvedGrade, bool)`
   and `Err() error` accessors, zero value an explicit refusal; the frozen-surface test pins
   those two accessors and reds on any added exported field or additional grade-bearing
   accessor. External package tests prove no
   authority-bearing grade API accepts a raw `HashRef`, decoded proposal evidence, receipt, or
   caller-written `EvidenceGrade`.
2. **Agent containment.** `DecodeProposal` may decode `ProofReceipt`, but it remains
   `ClaimedEvidence`; canonical `gradeOf` returns `CLAIMED`, and no host grade API accepts it.
   `DecodeProposal` refuses raw input over 256 KiB before any parse (§3.3; the iteration-87
   note, discharged).
3. **Bounded authenticated validation.** The envelope is bounded at I/O: streamed from the
   store's bounded read through an `io.LimitReader` sized 256 KiB + 1, never fully materialized
   first. Readings that make this criterion falsifiable are scoped to CODE, not to this
   document's prose (which already quotes these tokens): base reading, 0 occurrences of
   `OpenObject` across the six non-test `host/store` files and no `host/evidence` package at
   `52bc9ec` (V36), re-confirmed 0 repo-wide under `host/` and `cmd/` at `7806cac` (V42); head
   reading, ≥ 1 `OpenObject` in non-test `host/store` and ≥ 1 `io.LimitReader` in non-test
   `host/evidence` — scopes in which base and head DIFFER by construction. The round-7
   alternative spelling, a `maxBytes` parameter on `GetObject`, is DELETED (V40, §10.8). Decoded reports are bounded before strict
   decode, hash is recomputed, semantic/interface identities match, envelope and report canonical
   bytes round-trip, HMAC-SHA-256 is compared in constant time, subject and compiler match, the
   `verify.results[]`-derived verified set is non-empty, and success/error/counterexample fields
   agree.
4. **No fallback.** Every validation failure yields its exact `UnsupportedReason` or an explicit
   operational error. No failure result carries any grade.
5. **Producer.** The pinned executable is byte/version checked; execution is time/output bounded
   by the tranche-1-owned runner inside `host/evidence` (§3.4, V35), with deadline overrun and
   output overflow each a refusal that emits no report, and stdout/stderr captured separately;
   JSON fields read from stdout—not rc—decide success; identities come from `verify.results[]`
   entries with `status == "verified"`; an empty required identity set is refused; only
   successfully MAC-tagged envelopes are stored.
6. **Kernel mapping.** Add only `ProofReceipt(HashRef)` in tranche 1 and map it to `CLAIMED`; extend
   exact contract/body and integer policy case. No kernel Evidence constructor produces `PROVEN`.
7. **AILANG pins.** Keep `gradeOf` named and `EXACT_TOTAL_VERIFIED=10`; add the observed seventh
   `gradeCode` test identity and move `EXACT_TOTAL_TESTS` 39 → 40. Do not invent
   `EXACT_TOTAL_MODULES`.
8. **Go persistent gate.** Add the exact non-empty named-test manifest and exact count to
   `scripts/verify_go.sh`; add an isolated self-mutation test in `host/verifygate`.
9. **Required mutations.** Every tranche-1 row M1–M5, M7–M19, M20, M21, and M22 reds with its
   named observable; controls green. M5 is the sole payload-hash mutation owner. M6 is an
   acceptance criterion of tranche 3 and replay remains unable to yield `PROVEN` before then.
10. **Projection/golden.** Canonical/projected `types.ail` are byte-identical, the frozen four
    exports and six tar entries remain unchanged, and the canonical ready packet golden is
    regenerated.
11. **Pinned baseline and final gates.** Before mutation, require the measured base gate from V17.
    After implementation run `AILANG_BIN=/tmp/ailang-v0300/ailang ./scripts/verify_ail.sh` and
    `AILANG_BIN=/tmp/ailang-v0300/ailang GOTOOLCHAIN=go1.25.6 ./scripts/verify_go.sh`; both must be
    green. A red base is repository failure, not change evidence.
12. **No premature display.** Tranche 1 exposes no renderer/public daemon route and no surface may
    display `PROVEN`. Tranche 4 must accept the sealed/revalidated value only.
13. **Forgery negative control.** With a correctly configured host key, hand-author canonical
    `ProofReportV1` bytes whose schema, subject, compiler hash/version, sorted non-empty verified
    identities, success flags, zero errors/counterexamples, envelope metadata, and receipt payload
    hash are otherwise perfect. For an absent tag and separately for a wrong 32-byte tag,
   `ValidateProof` must return the exact stable reason string `unauthenticated_report` and no
    `ValidatedEvidence`; under M10, both assertions must fail. The reason is read from the
    validator result after the MAC branch, not from fixture state written beside it.
14. **Cross-validator binding.** `Validator.Resolve` refuses a `ValidatedEvidence` minted under
    any other mint identity (per-identity, not per-Go-variable: a value copy of a validator
    carries the same identity and key and is the same authority — §3.2):
    `TestAttackerChosenValidatorCannotMintForHostAuthority` constructs validator 2 with its own
    key, mints a seal through it, presents that seal to validator 1's `Resolve`, and requires a
    `ResolutionResult` whose `Err()` is the dedicated `ErrForeignSeal` and whose `Proven()`
    reports no grade — after first resolving a same-instance seal to
    `ResolvedGradeProven` as its control. The enforced property is exactly that a caller cannot
    make someone else's validator resolve their seal; self-minting into a caller-constructed
    validator remains possible and is accepted, because no library can stop a caller lying to
    itself. M20 must red this arm. (This numbered slot was vacated in round 4 when the
    unsatisfiable startup-key AC was removed; this is a new criterion, not that one revived.)
15. **Mint validity.** A zero or unset identity never resolves to any grade:
    `TestZeroValueForgeryCannotResolve` constructs a zero-value `Validator` and a zero-value
    `ValidatedEvidence` from the external test package — the forge unexported fields cannot
    prevent, because exported types always have constructible zero values — calls
    `Resolve` on the pair, and requires a `ResolutionResult` whose `Err()` is the dedicated
    `ErrUnmintedAuthority` and whose `Proven()` reports no grade,
    after its same-test control resolves a `NewValidator`-minted same-instance seal to
    `ResolvedGradeProven`. The same test constructs `var r evidence.ResolutionResult` and
    requires `Proven()` false: the zero value of the result type is itself an explicit refusal,
    never success. The nil-identity refusal executes BEFORE AC14's binding comparison,
    so the zero-zero pair is refused rather than compared equal. `ErrUnmintedAuthority` is
    produced by the mint-validity check and nowhere else and is distinct from `ErrForeignSeal`,
    because "never minted" and "minted by someone else" are different refusals. M21 must red
    this arm.
16. **Deadline-bounded validation (wait bound).** A context parameter is not a bound; a deadline
    is — 11 ratified production store reads already pass `context.Background()` (V41), so the
    validator may inherit no deadline from anyone. Validator configuration carries a required
    positive `ObjectReadTimeout`; `NewValidator` refuses zero or negative. `ValidateProof`
    derives `context.WithTimeout(ctx, ObjectReadTimeout)` before opening or reading the envelope
    object, even when the caller supplies no deadline — `context.Background()` included. Expiry
    or cancellation is the explicit operational error of §2.5, mints no seal, and carries no
    grade. Named test `TestBlockingObjectReadReturnsWithinObjectReadTimeout` proves it with a
    blocking reader: a fake `ObjectReader` whose read completes only when its ctx is done, called
    under `context.Background()`, must return the operational error within the configured bound
    (guarded by a test-side 20× timer so the failure mode is a red, not a hang), after the
    same test's prompt-reader control seals and resolves to `ResolvedGradeProven`. Readings that
    make this criterion falsifiable are scoped to CODE, not to this document's prose (which
    quotes the token throughout): base reading, 0 occurrences of `ObjectReadTimeout` in `*.go`
    under `host/`, `cmd/`, and `scripts/` at `7806cac` (V42); head reading, ≥ 1 in non-test
    `host/evidence` plus ≥ 1 `context.WithTimeout` in non-test `host/evidence` — scopes in which
    base and head DIFFER by construction (the base `context.WithTimeout` control count of 8 sits
    entirely outside `host/evidence`, which does not exist at base). M22 must red this arm.

### 8.1 Five coupled AILANG moves

- `world/types.ail`: adds one untrusted constructor, `CLAIMED` contract/body arm, and integer case.
- `packages/world-core/world/types.ail`: byte-identical projection changes in the same commit.
- `scripts/verify_ail.sh`: `REQUIRED_TESTS` and Python `EXACT_TOTAL_TESTS = 39` move to 40; shell
  `EXACT_TOTAL_VERIFIED=10` and `REQUIRED_VERIFIED["world/types.ail"]={"gradeOf"}` remain pinned.
- `scripts/verify_world_package.sh` step 4 retains the frozen four-module export manifest; step 3 sees
  new projected bytes; step 9 sees new content/interface/tarball hashes.
- `scripts/world_package_ready_packet.golden.json`: canonical byte golden changes.

`LEG1_MODULES` remains the same set and retains its anti-vacuity floor. There is no
`EXACT_TOTAL_MODULES` pin (V5, V13).

### 8.2 Go packages and gate

- New `host/evidence`: codec, validator, producer, resolved-grade API, focused tests, mutation
  harness.
- `host/store`: **IN the tranche**, by the attended `D-WORLD-19` arm-A ruling of 2026-08-19 —
  a ratified widening of the previously frozen package list, not scope creep. The store gains
  exactly ONE additive bounded object-read surface, the round-7b reconciled spelling (§10.8):
  `OpenObject(ctx context.Context, ref hashref.HashRef, maxBytes int64) (ObjectMeta, io.ReadCloser, error)`.
  `GetObject` is NOT modified: 13 non-test out-of-tranche `.GetObject(` call sites exist under
  `host/` and `cmd/` outside `host/store` (V40), two of them in `host/daemon` and `host/replay`,
  which this very list declares frozen — so the round-7 alternative arm, a `maxBytes` bound on
  `GetObject`, was unsatisfiable by construction and is DELETED, not retained as an option.
  `ObjectMeta` carries the non-payload columns `GetObject` already returns (`InterfaceHash`,
  `SemanticID`, `Provenance` — V36's refresh at `store.go:468`) that §3.3 steps 4–5 check; a
  reader-only return would strip data the validator consumes. No schema change; the objects
  table is untouched. The binding invariant: no read opened for an untrusted ref may materialize
  more than 256 KiB + 1 bytes of payload. Sequencing with item 18
  (`w-daemon-read-cancellation`): item 18 LANDED first and bounds the SIGNATURE — every store
  read now takes a `context.Context` — but a context parameter is not a wait bound: 11
  production store reads pass `context.Background()` by item 18's own ratified deferral DR-2,
  pinned by `TestNoNewDeadlineFreeStoreReads` (V41). Those 11 sites are item 18's follow-on
  obligation (the pin's own comment names the 11 → 0 path) and MUST NOT be absorbed here; this
  tranche imposes its own deadline at its own call site via `ObjectReadTimeout` (§3.3 step 2,
  AC16). Tests may use Store as an integration control.
- `host/hashref`: reused unchanged.
- `host/verifygate`: gains isolated mutation coverage for the named Go manifest.
- `scripts/verify_go.sh`: gains a focused JSON-parsed exact named-test leg; broad build/plain/race
  legs remain.
- `.github/workflows/ci.yml`: no new job is required because CI already invokes both scripts;
  confirm the Go job exports the pinned binary and Go toolchain during implementation.
- `host/daemon`, `cmd/**`, `host/replay`, and renderer surfaces do not change in tranche 1.
  Therefore tranche 1 has no production composition root and is not deployable until the
  production MAC-key wiring item in §9 lands.

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
| **1. This document, `w-validated-proven-evidence-boundary`** | Library-only proof report schema/producer and authenticated envelope; caller-supplied 32-byte MAC key; Evidence codec; bounded/hash/type/success validation; bounded `host/store` object read (`D-WORLD-19` arm A, round-7b reconciled `OpenObject` spelling); validator-imposed `ObjectReadTimeout` wait bound; sealed mint authority; untrusted `ProofReceipt → CLAIMED`; host-only resolved `PROVEN`; refusal-branch mutations; persistent named gates. Explicitly non-production. | **4.25 d** |
| **2. `w-proven-evidence-production-key-wiring`** | Verify and name the actual production composition root and configuration/state mechanisms; provision/load durable MAC key material; retain the keyed validator before serving; abort startup closed. | **1.0 d** |
| **3. `w-validated-replay-evidence-boundary`** | Typed replay report and untrusted replay receipt; integrate only successful full-episode replay; bind episode/log head and interpreter set; make missing/failed/divergent replay explicitly unsupported; M6 persistent mutation. `RecordedEffect` stays `ATTESTED`. | **3.0 d** |
| **4. `w-proven-evidence-renderer-consumption`** | Renderer/read API that accepts only sealed or freshly revalidated evidence; display `PROVEN` only from that value; explicit `UNSUPPORTED` for every validation failure; end-to-end agent-forgery and restart/revalidation tests. | **2.0 d** |
| **Total** | All reviewer obligations | **10.25 d** |

Tranche 1 arithmetic:

| Work | Time |
|---|---:|
| Strict report/envelope/Evidence codecs and object integration | 0.75 d |
| Bounded store read (`OpenObject` with `ObjectMeta`), `io.LimitReader` streaming, store-side tests | 0.50 d |
| `ObjectReadTimeout` configuration, `WithTimeout` derivation, blocking-reader test, M22 | 0.25 d |
| Bounded pinned proof producer, caller-key MAC integration, and fixtures | 0.60 d |
| Validator, sealed authority, receipt containment, resolved-grade API | 0.75 d |
| Kernel mapping, projection, golden, AILANG pins | 0.35 d |
| Named Go-test manifest and self-mutation gate | 0.45 d |
| Mutations, full pinned gates, review contingency | 0.60 d |
| **Total** | **4.25 d** |

The tranche-1 estimate removes 0.5 day previously assigned to an unwired key lifecycle, and
round 7 adds 0.50 day for the `D-WORLD-19` arm-A bounded store read — priced as its own row
rather than absorbed, because `host/store` entering the tranche is new work, not a re-labelling.
Round 7b adds 0.25 day for the deadline machinery (`ObjectReadTimeout` plumbing through
configuration and constructor refusal, the blocking-reader test, M22), again priced rather than
absorbed into the contingency row. At 4.25 days the tranche now EXCEEDS the 3–4 day sprint
guardrail by a quarter day. That is stated rather than hidden, and it is the re-quorum's and the
human's to weigh — not the pricing table's to disguise by leaving 4.0 because it was 4.0.
Authenticated-envelope canonicalization, mechanism-specific negative controls, and refusal-branch
mutations remain; production provisioning/loading and composition-root tests move intact to
tranche 2. No tranche adds a rotation protocol or a validator-side compiler execution path.

Tranche 2 has its own acceptance criteria: (a) its design verification log must enumerate the
existing production configuration type and startup function, private-state-directory handling,
secret/key handling, and secure or atomic file-write helpers, each with a command and observed
output before reuse is claimed; (b) it must name the exact `host/daemon` and/or `cmd/**` files
brought into scope; (c) initialization must provision or load exactly 32 bytes, retain the
resulting keyed validator before serving requests, and abort startup on any failure; and (d) an
integration test must start that actual composition root against missing, valid, symlinked,
wrong-permission, and replaced key cases and observe the specified fail-closed behavior. For
those production-wiring surfaces, no existing helper or path is assumed by tranche 1 and reuse
decisions belong to that measured successor design. The bounded subprocess runner is the named
exception and is NOT on tranche 2's list: tranche 1 owns it, measured its reuse decision in this
document's own log (V35: the only general runner is unexported and broker-internal, so nothing
reusable exists), and implements it inside `host/evidence` (§3.4).

Ordering is binding. Tranche 1 is non-production until tranche 2 wires key custody. Tranche 3
cannot reuse raw proof receipts as replay receipts. Tranche 4 cannot infer trust from either
serialized constructor. Until tranche 3 lands, replay evidence cannot yield `PROVEN`; until
tranche 4 lands, no renderer may display `PROVEN` at all.

## 10. What this is NOT doing

- It does not claim an AILANG constructor is opaque or unforgeable, or that every grade consumer
  passes through Go.
- It does not accept a report because its `HashRef` is well formed or because Store once inserted it.
- It does not add replay evidence in tranche 1 or promote `RecordedEffect` above `ATTESTED`.
- It does not treat a replay cache hit, execution success, or absence of an error as replay proof.
- It does not add a daemon route, CLI verb, agent tool, renderer, or public proof-as-a-service.
- It does not provision/load a production key, wire a production composition root, or claim that
  tranche 1 is deployable; those are explicit tranche-2 obligations.
- It does not define aggregate grade ordering for `list[Evidence]`, empty evidence, or proposal
  confidence.
- It does not put I/O, decoding, hashing, or solver execution in `world/`.
- It does not attach a contract to `Proposal`. The measured verifier limitation is confined to
  function parameters whose record type contains an ADT-typed field, bare or inside a list; it
  makes no claim about deeper nesting, result-position records, or tuples (V31–V32).
- It does not change store schema, world commit policy, replay semantics, effect broker, package
  version, or registry publication.
- It does not allow `PROVEN` rendering before ordered tranche 4.

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
At the refreshed base, `PutObject` has 8 non-test call sites and the same
`host/broker/broker.go:289` path stores bytes derived from an
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

## 10.3 Round 3 revision

- **Human direction — RATIFIED option B, attended.** Mark Edmondson ratified a host-held MAC key
  on 2026-08-14. That revision made HMAC-SHA-256 provenance a prerequisite for `ValidateProof`
  and preserved hash recomputation as integrity rather than authentication. Round 4 retains that
  primitive but supersedes the tranche-1 production key-custody claims. Option A is rejected on
  critical-path compiler/solver cost and architecture; option C is rejected as a resting state.
- **`gpt5-6-sol` — objection SUSTAINED and answered.** The validator no longer trusts
  self-asserted canonical bytes. M10 and AC13 hand-author an otherwise-perfect report whose
  receipt hash is correct but whose tag is absent or wrong, and require the unique downstream
  `unauthenticated_report` reason. The tag-guard mutant is the compiling
  `if false && (len(tag) != sha256.Size || !hmac.Equal(want, tag))` form the objection's `catch`
  demanded.
- **`gemini-3-1-pro` — premise repaired, weakening rejected.** `verified` is now explicitly sourced
  from sorted `verify.results[].function` entries filtered by `status == "verified"`; the integer
  `verify.verified` is only an optional cheap precondition. The current pinned run exposes six
  per-identity results and includes the known-positive `gradeOf` entry (V27).
- **Freshness corrections at `bef0153`.** The declared base and required-file counts are refreshed
  (V1–V2); the gate pins are corrected from 5/20 to 10/39 before the proposed seventh test moves
  39 → 40 (V5, V17); the per-identity checker measurement is rerun with separated streams (V27);
  and the live threat-model enumeration is corrected from 10/16 to 8/13 non-test Put/Get call
  sites while retaining the eight-route control (V28). V29 records the seven-file freshness sweep
  that selected the affected rows.

## 10.4 Round 4 revision

- **Objection A — COMPLETENESS/SCOPE, SUSTAINED.** The previous AC14 required daemon-owned
  first-startup key creation while §3.4 and §8.2 excluded every production composition root. That
  AC was unsatisfiable by construction; this objection is accepted, not argued away. This revision
  takes **Arm 2** because the existing decomposition and package boundary already define tranche 1
  as a library integration, while no verified production config/startup/key-file mechanism is
  available in this document. Tranche 1 is now explicitly non-production, accepts a caller-owned
  32-byte key rather than a path, and makes no startup, daemon custody, or restart claim. The named
  successor `w-proven-evidence-production-key-wiring` owns the production inventory, exact
  composition-root changes, fail-closed initialization, and five-case startup integration test.
  This preserves the ratified MAC provenance primitive without pretending it is deployed.
- **Objection B — PREMISE, SUSTAINED and measured.** Pinned same-call probes correct three parts of
  the earlier description. First, the measured failure is not specifically `list[ADT]`: a function
  parameter whose record has either a bare ADT field or a `list[ADT]` field errors with unknown
  record sort. Second, records and lists are not independently the problem: scalar-only records
  and records containing `list[scalar]` verify, as does a bare ADT parameter. Third, the contract
  need not read the ADT-bearing field; the parameter type alone triggers the failure. The scope is
  limited to those measured parameter shapes and does not cover deeper nesting, result-position
  records, or tuples (V31–V32).
- **Silent failure characteristic.** Both pinned `ai-check` calls return rc=0 with
  `check.passed=true` and `check.error_count=0`; only `verify.errors` and the affected
  `verify.results[].status="error"` reveal the Z3 encoding failure. The repository comment at
  `world/contracts.ail:11-13` corroborates the Proposal case and cites upstream issue 477, but the
  first-party probe rows—not that comment—support this design boundary (V31–V32).

## 10.5 Round 5 revision

- **Direction — RATIFIED option A, attended.** Mark Edmondson ratified decision-ledger row
  `D-WORLD-17` in `design_docs/world-mission.md` on 2026-08-17 (recommendation adopted verbatim):
  bind every seal to its minting validator; tranche 1 stays library-only and explicitly
  NON-PRODUCTION. This closes the round-4 blocking catch — with no production root, the key is
  caller-supplied, so a public `NewValidator` plus a free `GradeOfValidated` let any Go caller
  construct a validator, mint a seal, and obtain `ResolvedGradeProven` — by making possession of
  a seal worthless without the minting validator instance. Option B was rejected for chaining the
  item behind an unwritten successor and re-creating the guard-lands-before-the-guarded-thing
  vacuity pattern (the AC12/AC13 precedent).
- **What changed.** The free `GradeOfValidated(sealed) ResolvedGrade` is DROPPED — not deprecated,
  not aliased, not wrapped, because its existence detached from a validator instance is the
  defect — and replaced by the `Validator.Resolve(sealed) ResolvedGrade` method (§2.2 item 3,
  §3.2, §3.5). Each validator carries an unexported per-instance identity stamped into the seal
  at mint; `Resolve` refuses foreign seals with the dedicated `ErrForeignSeal` (§3.2). The
  cross-validator refusal arm lands as named RED mutation M20 with
  `TestAttackerChosenValidatorCannotMintForHostAuthority` (§6, control paragraph included), and
  new AC14 states the enforced property (§7). All five pre-revision references to the dropped
  symbol were moved (V33); the symbol appears in no `.go` or `.ail` file, so this is a
  document-only change (V34).
- **What is and is not enforced.** You cannot make someone else's validator resolve your seal;
  you can still self-mint into a validator you constructed yourself, and that is accepted — no
  library can stop a caller lying to itself. Round 4's explicitly non-production tranche-1
  framing is unchanged and still correct; production key custody stays with
  `w-proven-evidence-production-key-wiring`.
- **Round 5, revision 2.** Quorum round 5 returned BLOCKED with both reviewers present; both
  objections were controller-measured first-party and SUSTAINED, both are COMPLETENESS defects,
  and the ratified `D-WORLD-17` direction is untouched.
  - **`gpt5-6-sol` — the binding was forgeable by Go zero values.** `Validator` and
    `ValidatedEvidence` are exported, so `var v Validator; var s ValidatedEvidence; v.Resolve(s)`
    compared zero identity to zero identity and passed — two lines of Go minting
    `ResolvedGradeProven` with no key and no `NewValidator` call; the round-4 defect one layer
    down (round 4 was a resolver plus a public constructor; this was a resolver plus no
    constructor at all). Fixed with a mint-validity invariant (§3.2): the identity is now an
    unexported pointer to a non-zero-size allocation made only inside `NewValidator`, and
    `Resolve` refuses any nil identity with the dedicated `ErrUnmintedAuthority` BEFORE the
    binding comparison, so a zero identity is never valid. The objection's copy-semantics half
    is answered by statement, not denial: a value copy of a validator carries the same identity
    pointer and key and IS the same authority — the guarantee is re-worded per-identity, not
    per-Go-variable (§2.2, §3.2, AC14), so the claim is no stronger than the mechanism. New RED
    mutation **M21** (`TestZeroValueForgeryCannotResolve`, observable `ErrUnmintedAuthority`)
    lands in §6 with a same-instance control; new AC15 states the criterion; criterion 9's
    required-mutation list gains M21.
  - **`gemini-3-1-pro` — §3.4 and §9 contradicted each other.** §3.4 mandated the tranche-1
    producer bound wall time, stdout, and stderr, while §9 placed the bounded subprocess runner
    on tranche 2's to-be-measured inventory and forbade tranche 1 assuming any helper — an
    acceptance criterion resting on a capability the same document called unmeasured and
    unowned. Resolved by branch **(a)**: tranche 1 owns the bounded runner, because the bounded
    producer is tranche 1's named deliverable and cannot emit an authenticated report without
    executing the pinned checker. New row V35 measures the primitive first-party (exactly three
    non-test `exec.CommandContext` sites; the sole general-purpose runner `runBounded` is
    unexported and broker-internal, so nothing reusable exists and reuse is rejected, not
    deferred); §3.4 now specifies the tranche-1 construction, including overrun and overflow as
    report-emitting-nothing refusals and separate stdout/stderr capture; §9's tranche-2 list
    drops the runner and names the exception; AC5 states the same ownership. Pricing is
    unchanged — the 0.60-day producer row already covered this construction.

## 10.6 Round 6 quorum — BLOCKED, and why this document PARKS again

Round 6 is the single re-quorum Gate 2 permits after a revision pass. Both reviewers were
present (`absent_reviewers: []`, per-reviewer cap raised to `$0.45` *before* the round, on a doc
that had grown 718 → 901 lines), `metered=$0.174209`. Verdict: **blocked**. Two objections
survive, and this document parks a fourth time — but note what changed: the ratified DIRECTION
(D-WORLD-17 arm A) was never in question in either round-5 or round-6 review. Every objection
across both rounds has been a completeness defect in the mechanism, and three of the four were
answered in-loop.

**Objection 6a — `gpt5-6-sol`, API shape. SUSTAINED, and NOT the reason this parks.**
`Validator.Resolve(sealed) ResolvedGrade` is a single-return signature with no refusal channel,
while §§2.2/3.2/3.5, M20/M21 and AC14/AC15 all require it to return `ErrUnmintedAuthority` or
`ErrForeignSeal` *and no grade*. The declared API cannot express what the whole authority
boundary rests on, and the failure mode is a silent default grade. The defect is the
controller's: the round-5 directive prescribed that signature verbatim from the attended
ratification's wording, and neither the designer nor round 5's reviewers questioned it.
Its `proposed_fix` is concrete, uncontested and direction-compatible — a sum-style
`ResolutionResult` with unexported fields and mutually exclusive `Proven() (ResolvedGrade, bool)`
and `Err() error` accessors, whose ZERO VALUE is an explicit refusal rather than success (the
same invariant M21 already pins one layer down). **This fix is pre-agreed and applies under
either arm of the ask below.** It is recorded here rather than applied because the
narrow-refinement carve-out is all-or-nothing: one disqualifying objection forecloses it for the
whole document.

**Objection 6b — `gemini-3-1-pro`, unbounded allocation on untrusted input. SUSTAINED, premise
MEASURED FIRST-PARTY, and this is what parks the document.**
§3.3 step 3 caps the envelope at 256 KiB "before allocating a second full copy" — which concedes
that the FIRST copy is already fully in memory, allocated by the store layer, before any length
check runs. An attacker-supplied `HashRef` to a multi-gigabyte object OOMs the process before the
validator ever looks at a length. The reviewer's premise was not forwarded on trust; the
controller re-derived it at `03c7892` (rule 3f):

- `host/store` exposes exactly **two** exported `Object` methods — `PutObject(o Object) error`
  (`store.go:443`) and `GetObject(ref hashref.HashRef) (Object, bool, error)` (`store.go:467`).
  `GetObject` returns the whole `Object`; there is no streaming or size-bounded form.
- Non-test `host/store` contains **zero** occurrences of `io.Reader`, `io.LimitReader`,
  `maxBytes` or `MaxBytes`. Known-positive control in the same call: **23** exported `Store`
  methods exist (`store.go` 14, `journal.go` 9), so the zero is a measurement, not a broken grep.

So the objection is true, and its `proposed_fix` — "modify §8.2 to permit extending the store API
with a bounded read method (`OpenObject(ref) (io.ReadCloser, error)` or a `maxBytes` parameter)"
— cannot be applied as a narrow refinement. It **widens §8.2's frozen package boundary into
`host/store`**, and `host/store`'s read path is the declared subject of a different ratified queue
item (18, `w-daemon-read-cancellation`, D-WORLD-18 arm A, whose own ratchet
`TestNoNewDeadlineFreeStoreReads` pins that package's residual read sites). Deciding whether
tranche 1 may extend `host/store` — and how that interacts with an item queued to bound the same
package — is a scope call, and Gate 2 forbids a controller-invented resolution of one. **Fourth
consecutive confirmation of the rule that a scope/direction dispute forecloses the carve-out.**

**THE ASK, one word.** Both arms keep D-WORLD-17 arm A intact and both include objection 6a's
pre-agreed fix.

- **A — tranche 1 may extend `host/store` with a bounded object read.** Adopt 6b's fix verbatim:
  `OpenObject(ref) (io.ReadCloser, error)` (or `GetObject` gaining a `maxBytes` bound), §3.3 step
  2 streams through an `io.LimitReader` at 256 KiB, and §8.2's package list widens to include
  `host/store`. Closes the OOM vector inside this tranche. Cost: a second item now writes to
  `host/store` while item 18 is queued to bound it, so the two must be sequenced or merged.
- **B — tranche 1 does not touch `host/store`.** The unbounded first copy is recorded as an
  explicit, named limitation of a tranche that is already declared library-only and
  NON-PRODUCTION, with an acceptance criterion stating the residual and a named successor owning
  the bounded read — either item 18 or the tranche-2 wiring item. Keeps the frozen package
  boundary and the tranche ordering; leaves a real resource-exhaustion vector unfixed in code that
  by construction never runs in production.

Ledger row: `D-WORLD-19`.

## 10.7 Round 7 revision (D-WORLD-19 arm A applied)

- **The ruling.** Mark Edmondson answered `D-WORLD-19` attended on 2026-08-19T04:57:30Z, on
  bookkeeping issue #68, verbatim: *"D world 19 - A yes"*. Arm A adopts both round-6 fixes in
  the reviewers' own words; nothing in this revision re-litigates either objection, and
  `D-WORLD-17` arm A, the MAC seam (option B), the non-production tranche-1 framing, and the
  §9 decomposition are all untouched.
- **6a applied (gpt5-6-sol's fix, verbatim).** The single-return
  `Validator.Resolve(sealed) ResolvedGrade` is replaced at the three normative sites (§2.2
  item 3, the §3.2 signature block, §3.5) by the sum-style
  `Validator.Resolve(sealed ValidatedEvidence) ResolutionResult`: unexported fields, mutually
  exclusive `Proven() (ResolvedGrade, bool)` and `Err() error` accessors, success carrying
  exactly `ResolvedGradeProven`, refusal carrying exactly one of `ErrUnmintedAuthority` or
  `ErrForeignSeal` with no grade, zero value an explicit refusal. M20, M21, AC1, AC14, and
  AC15 are updated in lockstep, and M21/AC15 additionally pin the zero-value
  `ResolutionResult` as refusal. §§10.3–10.6 retain the old spelling, as records must. The
  round-6 record already assigns this defect to the CONTROLLER — the round-5 directive
  transcribed the signature from the attended ratification's wording — so the old spelling is
  a propagated transcription error, not load-bearing precedent.
- **6b applied (gemini-3-1-pro's fix, verbatim).** §8.2's package list widens to include
  `host/store`, which gains exactly one bounded object-read surface
  (`OpenObject(ref) (io.ReadCloser, error)` or a `maxBytes` bound on `GetObject`), and §3.3
  step 2 now streams the payload through an `io.LimitReader` sized 256 KiB + 1, so the first
  copy is bounded and the retired "before allocating a second full copy" wording — which
  conceded the unbounded first copy — is gone. M4 moves to the streaming overflow guard with a
  corrected failure prediction (`got hash_mismatch`, since a truncated payload fails the
  recomputed hash at step 4).
- **Item 18 landed in between — and did NOT close 6b.** Between rounds 6 and 7, item 18
  (`w-daemon-read-cancellation`) landed and threaded `context.Context` through `GetObject`
  (V36). A reader checking "has item 18 already fixed this?" finds a freshly context-threaded
  store and could conclude the OOM vector is closed. It is not: context bounds the WAIT, not
  the BYTES. §3.3 and §8.2 both state this explicitly because it is the single most likely way
  the bounded read gets silently dropped by a future implementer.
- **Iteration-87 note discharged.** `gemini-3-1-pro`'s non-blocking note — `DecodeProposal(raw)`
  stated no byte bound before parsing — was carried by the item head and dropped by two
  consecutive directives, and the same byte-bound concern returned one round later as blocking
  6b. It is now folded in: `DecodeProposal` caps raw input at 256 KiB before any parse (§3.2
  API list, §3.3, AC2), recorded here so it cannot drop again.
- **The vacuous-AC trap, avoided by scope.** The token `ResolutionResult` already appears in
  this document at base — once, in §10.6's own account of round 6 — so any criterion counting
  it in prose reads ≥ 1 at base AND head and is vacuous (the iteration-92/93/94 shape, three
  times running). Every new reading in this revision is therefore scoped to named CODE files
  with a stated base and head reading that differ (AC3, V36).
- **Freshness refresh at `52bc9ec`.** Round 6's recorded control of "23 exported `Store`
  methods" is stale: at HEAD it is 25 across the six non-test `host/store` files (V36). V5's
  cited pin lines moved 311/350 → 323/363 with values unchanged (V37). The sweep from the
  earliest declared base touches 41 non-design files (V38); the rows this revision relies on
  were re-measured at `52bc9ec` (V36–V38, plus in-place refresh notes on V5/V10 and re-run
  confirmations recorded in the §11 preamble).
- **Pricing.** `host/store` entering the tranche adds 0.50 day: tranche 1 moves 3.5 → 4.0 days
  and the decomposition total 9.5 → 10.0 days (§9).

## 10.8 Round 7b revision (carve-out: both round-7 objections, reviewers' verbatim fixes, reconciled)

- **Why a revision and not a park.** Round 7's quorum returned blocked with both reviewers
  present, but both objections AFFIRM the attended `D-WORLD-19` arm-A direction and correct only
  the API shape and one false sentence, each with a concrete reviewer-authored `proposed_fix` —
  and the two fixes CONVERGE on the same corrected signature. That is the narrow-refinement
  carve-out's exact precondition, and the carve-out was used. Convergence, not verbatim
  application, is what justifies proceeding: this document's own history holds two rounds in
  which a reviewer's verbatim fix was applied and then refuted by its own author (round 3 →
  round 4, gpt5-6-sol's arm 2; round 6 → round 7, gemini-3-1-pro's
  `OpenObject(ref) (io.ReadCloser, error)`, which was RATIFIED BY THE HUMAN and still wrong —
  it demanded context-threading in prose while dropping the `context.Context` parameter from
  the signature). "Applied verbatim" is not self-certifying here; two independent providers
  arriving unprompted at the same correction is the stronger warrant.
- **Objection 7a — gemini-3-1-pro, SUSTAINED (against its own round-6 fix).** Both round-7
  spellings were defective. A `maxBytes` bound on `GetObject` changes a signature with 13
  non-test out-of-tranche call sites under `host/` and `cmd/` (V40), two of them in
  `host/daemon` and `host/replay`, which §8.2 itself declares frozen — unsatisfiable by
  construction. `OpenObject(ref) (io.ReadCloser, error)` drops the context and breaks item 18's
  wait-bounding discipline. Fix applied verbatim: the `GetObject` arm is DELETED from §3.3,
  §8.2, and AC3 — not retained as an alternative — and the new surface accepts
  `ctx context.Context` first.
- **Objection 7b — gpt5-6-sol, SUSTAINED, and the controller measured it as WORSE than
  stated.** The round-7 sentence "item 18 added `context.Context` (a WAIT bound)" was FALSE.
  The reviewer argued a caller MAY pass `context.Background()`; at HEAD, 11 production store
  reads ALREADY DO, ratified as item 18's DR-2 deferral and pinned by
  `TestNoNewDeadlineFreeStoreReads` (V41). A context parameter is not a bound; a deadline is.
  Item 18 bounded the SIGNATURE, and the tranche must impose its own deadline at its own call
  site rather than inherit one. Fix applied verbatim: required positive `ObjectReadTimeout` in
  validator configuration; `ValidateProof` derives `context.WithTimeout(ctx, ObjectReadTimeout)`
  before opening or reading, even when the caller supplies no deadline; timeout/cancellation is
  an explicit operational error minting no seal (§2.5); named blocking-reader test; mutation
  M22; new AC16. This requirement is load-bearing, not belt-and-braces: the deadline the
  validator derives at its own call site is the only wait bound standing between it and the
  ratified deadline-free residue. The 11 DR-2 sites themselves are item 18's follow-on and are
  deliberately NOT absorbed here (§8.2; rule stated so no implementer absorbs them).
- **The reconciled signature — a reconciliation, not a third invention.**
  `OpenObject(ctx context.Context, ref hashref.HashRef, maxBytes int64) (ObjectMeta, io.ReadCloser, error)`
  is gpt5-6-sol's fuller spelling, adopted because it satisfies gemini-3-1-pro's fix STRICTLY:
  gemini's own text asked only that the context be present and `GetObject` be left alone, and
  this signature is gemini's corrected signature plus the byte bound and metadata gpt5 requires.
  `ObjectMeta` is kept with a measured reason, not by default: `GetObject` returns
  `InterfaceHash`, `SemanticID`, and `Provenance` beside the payload (V36's refresh at
  `store.go:468`), and §3.3 steps 4–5 check exactly those fields, so a reader-only return would
  strip data the validator consumes.
- **Supersessions.** §10.7's sentence "item 18's `context.Context` bounds the WAIT, not the
  BYTES" — and the copies of that claim that round 7 wrote into §3.3, §8.2, and V36 — is
  superseded by §3.3's corrected paragraph, §8.2's corrected bullet, V36's in-row correction
  note, and V41. §§10.1–10.7 remain byte-identical as records; the false sentence stands in
  §10.7 as history of what round 7 believed. The false sentence's provenance is already
  recorded: the controller wrote it into the round-7 directive and the designer transcribed it.
- **What did not change.** `D-WORLD-17` arm A, the MAC seam (Mark's option B), `D-WORLD-19`
  arm A (`host/store` IS in the tranche), the non-production tranche-1 framing, the §9
  three-successor decomposition, and every §10.7 change not named above.
- **Pricing.** +0.25 day for the deadline machinery: tranche 1 moves 4.0 → 4.25 days, the
  decomposition total 10.0 → 10.25, and the tranche now exceeds the 3–4 day guardrail by a
  quarter day, stated in §9 rather than hidden.

## 10.9 Round 8 quorum — ONE PASS, ONE REJECT, and the document PARKS on `D-WORLD-21`

Round 8 was a CONFIRMING re-quorum, run after the round-7b carve-out revision even though the
carve-out permits routing straight to sprint-planner. The reason is recorded in §10.8 and is
document-specific: a verbatim-applied reviewer fix has been refuted by its own author twice here
(rounds 3→4 and 6→7), so "applied verbatim" is not self-certifying on this doc.

`absent_reviewers` empty, both slots present, `metered=$0.244360`.

**`gemini-3-1-pro` → PASS.** Only the second reviewer flip to `pass` in this item's history (round 4
was the first), and the first on a round where the store surface was in scope. Its one remaining
concern was filed NON-BLOCKING: §3.4 said output is read through *a* capped reader (singular) for a
capture the same paragraph describes as dual-stream, leaving stderr's bound implicit, so an
attacker-controlled checker could spam stderr and OOM the host while the parsed stream stays inside
its cap. Its one-sentence `proposed_fix` is applied VERBATIM in §3.4 in the same commit that records
this round — deliberately, because the round-5 controller miss recorded in §10.6 was exactly a
non-blocking note left unforwarded that returned one round later as a BLOCKING objection.

**`gpt5-6-sol` → REJECT, and it is the same class one layer deeper for the third round running.**

| round | the thing that looked like a bound | why it was not one |
|---|---|---|
| 7 | a `context.Context` PARAMETER on `GetObject` | a parameter binds nothing until something arms it with a deadline; 11 production reads pass `context.Background()` by ratified deferral DR-2 (V41) |
| 7b fix | `context.WithTimeout` derived inside `ValidateProof` | bounds the OPEN call, not the subsequent read |
| 8 | the deadline reaching the read at all | `OpenObject` returns an ordinary `io.ReadCloser`, whose `Read` takes no context and is under no obligation to unblock when `ctx` expires |

The objection verbatim: *"`context.WithTimeout` therefore bounds the open operation but does not
inherently bound the subsequent `io.LimitReader`/`Read` loop. M22 uses a fake reader deliberately
wired to observe the context, so it can pass while the real `host/store` reader blocks
indefinitely."*

The second sentence is the sharper half and it lands on this document's own mutation table.
**M22 is vacuous by construction, and the vacuity is supplied by the FIXTURE rather than by the
production code** — a context-observing fake makes the mutation die for a property the real store
reader has never been shown to have. That is iteration 92's finding arriving in its inverse
direction (there, a doc-prescribed fake returned the exact value a neutered arm needed, so the
mutant PASSED; here, the fake is STRONGER than production, so the test passes where production
would hang), and nothing in this loop audits what a prescribed fake makes true.

**WHY THIS PARKS RATHER THAN TAKING ANOTHER CARVE-OUT.** `gpt5-6-sol`'s `proposed_fix` offers two
arms, and they are not interchangeable:

- **Arm 1** — replace the seam with `ReadObject(ctx, ref, maxBytes) (ObjectMeta, []byte, error)`,
  requiring the store to enforce `maxBytes` before materialization and perform the COMPLETE read
  under the supplied context. This makes cancellation enforceable, and it **returns a byte slice** —
  i.e. it retires the streaming reader that objection 6b demanded, that `D-WORLD-19` arm A ratified,
  and that `gemini-3-1-pro` has just PASSED. Applying it verbatim would undo the round that finally
  satisfied the other reviewer.
- **Arm 2** — keep the reader and "define a context-aware reader contract whose reads are guaranteed
  to terminate on `ctx.Done()` and document the concrete store mechanism that enforces it."

So the two reviewers now want different things at the same seam: streaming to avoid materialization
versus a complete-read-under-context to make cancellation enforceable. That is a DIRECTION dispute
between reviewers, not a completeness defect, and §10.4's discriminator applies — *is what the fix
REMOVES load-bearing for the doc's claim?* Arm 1 removes exactly the property the other reviewer's
objection was about. **Fifth consecutive confirmation that a direction dispute forecloses the
narrow-refinement carve-out**, and the first where the dispute is between the two REVIEWERS rather
than between a reviewer and the design.

Both arms also share one obligation that is not contingent on the answer, and it should be
implemented whichever arm wins: `gpt5-6-sol`'s **real-store integration test** — a blocked or
contended object read under `context.Background()`, relying only on `ObjectReadTimeout`, proving the
read returns within the bound with no seal, and mutating **the actual cancellation mechanism rather
than `WithTimeout`**. That is the arm that would have caught M22's fixture vacuity.

**THE ASK IS ONE WORD** (`D-WORLD-21`).

## 11. Verification Log

Rows V1–V8, V12–V14, V17, V21, and V27–V32 were rerun from repository root at
`bef0153` on 2026-08-15. Other rows are retained old-base measurements whose cited repository
paths are outside V29's seven changed non-design files, or are explicitly labelled inherited/history.
Negative measurements assert their roots first and include a same-scope positive control in the
same call. Glob-shaped arguments are quoted. V18, V22, V24, and V25 use only scratch files under
`/tmp` when they write. V33–V34 were run for the round-5 revision on 2026-08-18 from repository
root at `03c7892`; V35 was run for round-5 revision 2 on 2026-08-18 at the same `03c7892` base.

V36–V38 were run for round 7 on 2026-08-19 from repository root at `52bc9ec`. Because the sweep
from the earliest declared base now touches 41 non-design files (V38), the round-7 pass re-ran
the cheap check of every retained row citing a swept file rather than marking any row fresh by
default: V2 (only `scripts/verify_ail.sh` moved, 386 → 399 lines), V6 (`PROVEN` 0, control 4),
V9 (divergence returns still 169/216/234 with success at 242; the `ReplayResult` control count
moved 15 → 17), V14 (0), V15 (0), V16 (exact-token and toolchain-denylist structure intact),
V21 (0, control 9), V28 (8 routes; 8 put / 13 get non-test sites — unchanged), and V35 (the
same three non-test `CommandContext` sites). V5 and V10 carry in-row round-7 refresh notes
(V37, V36). V17 remains a `bef0153` gate measurement: it is not re-certified here, and AC11
re-establishes the base gate at implementation time. Round-7 re-runs use `grep` (this
environment's `rg` is unavailable); commands shown in V36–V38 are the ones actually executed.

V40–V42 were run for round 7b on 2026-08-19 from repository root at `7806cac`. V36 carries an
in-row round-7b correction note: its measurements stand, but three words of its claim — "a WAIT
bound" — were false and are corrected there rather than silently rewritten (§10.8, V41).

| ID | Claim | Exact command and same-call control | Observed output |
|---|---|---|---|
| V1 | Measurement base and clean initial tree before this revision. | `git rev-parse --short HEAD; git status --short` | `bef0153`; no status lines. |
| V2 | Required inputs exist in this worktree. | `test -f design_docs/implemented/w-evidence-grade-mapping.md && test -f world/types.ail && test -f scripts/verify_ail.sh && test -f host/replay/replay.go; wc -l design_docs/implemented/w-evidence-grade-mapping.md world/types.ail scripts/verify_ail.sh host/replay/replay.go` | All tests rc=0; respective counts `631`, `279`, `386`, `396`. |
| V3 | At `bef0153`, production Go has no Evidence constructor/decoder; test scope is visible. | `test -d host && test -d cmd; printf prod=; grep -rn "Evidence" host/ cmd/ --include='*.go' \| grep -v '_test\.go' \| wc -l; printf control=; grep -rn "Evidence" host/ cmd/ --include='*_test.go' \| wc -l` | roots yes; `prod=0`; same-scope test control `13`. |
| V4 | At `bef0153`, there is no production Z3 report producer; same production scope contains hashref code. | `test -d host && test -d cmd; grep -rinE 'Z3\|z3' host/ cmd/ --include='*.go' \| grep -v '_test\.go' \| wc -l; grep -rin 'hashref' host/ cmd/ --include='*.go' \| grep -v '_test\.go' \| wc -l` | `9` Z3 hits (inspection: comments); same-scope hashref control `424`. |
| V5 | At `bef0153`, the gate pins exactly 10 verified identities and 39 named tests; within `REQUIRED_TESTS`, exactly six names have the `gradeCode_test` prefix, and `gradeOf` remains required. | `grep -nE "^EXACT_TOTAL_VERIFIED=\|^EXACT_TOTAL_TESTS = " scripts/verify_ail.sh; python3 -c "import re; s=open('scripts/verify_ail.sh').read(); m=re.search(r'REQUIRED_TESTS = \{(.*?)\}', s, re.S); names=re.findall(r'\"([a-zA-Z0-9_]+)\"',m.group(1)); print(len(names)); print([x for x in names if x.startswith('gradeCode_test')]); m2=re.search(r'REQUIRED_VERIFIED = \{(.*?)\n\}',s,re.S); print('gradeOf', 'gradeOf' in m2.group(1))"; git log --oneline -1 -S "EXACT_TOTAL_VERIFIED=10" -- scripts/verify_ail.sh` | Lines 311/350: `10` and `39`; Python enumerates `39`, the six `_1`…`_6` names, and `gradeOf True`; introducing commit `aaada20 feat(15): freeze the v1 DecisionPacket…`. **Round-7 refresh at `52bc9ec`: values unchanged, lines moved to 323/363 by item 18's landing (V37).** |
| V6 | No additional PROVEN grep exists; the gate file is searchable. | `{ rg -n 'PROVEN' scripts/verify_ail.sh .github/workflows --glob '*.yml' || true; } \| wc -l; grep -c 'EXACT_TOTAL_VERIFIED' scripts/verify_ail.sh` | PROVEN lines `0`; same-file control `4`. Thus prohibition is the six pinned expectations and nothing else. |
| V7 | In `world/*.ail`, only `world/types.ail` contains `Evidence` (8 lines); `gradeOf` occurs only in canonical/projected types, while the same scoped instrument finds 10 `proposalMatchesWorld` controls. | `grep -rc Evidence world/*.ail; rg -n gradeOf world packages/world-core/world --glob '*.ail'; printf control=; rg -n proposalMatchesWorld world packages/world-core/world --glob '*.ail' \| wc -l` | `world/types.ail:8`, the other three world modules `0`; `gradeOf` lines 30/41/71 in canonical and projection; `control=10`. |
| V8 | AILANG `Proposal` contains `list[Evidence]`; the scoped production fixtures write empty evidence lists. | `rg -n 'export type Proposal\|evidence:' world/types.ail \| head -10; rg -n 'evidence: \[\]' world/transitions.ail world/contracts.ail` | `Proposal` starts at 89 and its Evidence field is line 97; empty-list hits occur twice in transitions and six times in contracts; the non-empty declaration in the same first call is the control. |
| V9 | Replay has real divergence and a success result only after comparisons. | `nl -ba host/replay/replay.go \| sed -n '90,250p' \| rg 'DivergenceError\|return ReplayResult'; ... \| rg ReplayResult \| wc -l` | typed divergence definition; returns at lines 169, 216, 234; success at 242; control count `15`. |
| V10 | Store Object carries payload, verifies on put, and returns full payload on get. | `nl -ba host/store/store.go \| sed -n '88,96p;440,492p' \| rg 'type Object\|Payload\|verifyObject\|GetObject\|SELECT interface_hash_ref'; rg -n 'func \(s \*Store\) (PutObject\|GetObject)' host/store/store.go \| wc -l` | payload field; `verifyObject` call; SQL payload lookup; two-method control. No Evidence-specific bound appears in this scope. **Round-7 refresh at `52bc9ec`: `PutObject` at :444; `GetObject` at :468 is now `(ctx context.Context, ref hashref.HashRef) (Object, bool, error)` — context-threaded by item 18 but still returning the whole `Object` with its full `[]byte` payload (V36).** |
| V11 | A strict bounded JSON pattern exists in transition registry. | `test -d host/transitionreg; rg -n 'DisallowUnknownFields\|multiple JSON values\|trailing JSON\|maxRevisionRaw\|maxSchemaRaw' host/transitionreg/codec.go; rg -n 'func parseJSON\|func DecodeRevision' host/transitionreg/codec.go \| wc -l` | raw caps include 262144 and 16777216; trailing/multiple checks; two function controls. |
| V12 | At `bef0153`, Go's authority-relevant Proposal is a different narrow struct and has no evidence field. | `nl -ba host/transitionreg/bind.go \| sed -n '20,30p'; nl -ba world/types.ail \| sed -n '87,100p'` | Go lines 22–27 contain transition/interpreter/epoch/caps/effects; AILANG's same-call control contains `evidence: list[Evidence]` at line 97. |
| V13 | Projection and fixed package surfaces are live at the new base. | `cmp -s world/types.ail packages/world-core/world/types.ail; echo types_cmp=$?; rg -n 'EXPECTED_MODULE_COUNT\|exact export set\|tar contains exactly\|ready-packet golden' scripts/build_world_package.sh scripts/verify_world_package.sh` | `types_cmp=0`; build pin `EXPECTED_MODULE_COUNT=4`; verifier lines require exact export set, exact tar entries, and the ready-packet golden. |
| V14 | At `bef0153`, production host/cmd has no grade/PROVEN renderer consumer; daemon scope is visible. | `test -d host && test -d cmd; { rg -n 'EvidenceGrade\|gradeOf\|PROVEN' host cmd --glob '*.go' --glob '!**/*_test.go' || true; } \| wc -l; rg -n 'handleObject\|handleWorld' host/daemon --glob '*.go' \| wc -l` | target `0`; same-scope daemon handler control `6`. |
| V15 | Go gate runs broad tests but has no named-test manifest today. | `{ rg -n 'REQUIRED.*TEST\|EXACT.*TEST\|test2json\|go test -json' scripts/verify_go.sh || true; } \| wc -l; rg -n 'go test ./\.\.\.' scripts/verify_go.sh \| wc -l` | named-pin target `0`; broad-test control `4`. |
| V16 | Go gate requires pinned AILANG and active Go is separately guarded. | `nl -ba scripts/verify_go.sh \| sed -n '19,41p;79,89p;102,126p'` | exact `v0.30.0` token check; go1.26.0–1.26.5 denylist; build/plain/race legs. |
| V17 | Pinned AILANG baseline is green and non-vacuous at `bef0153`. | `export AILANG_BIN=/tmp/ailang-v0300/ailang; $AILANG_BIN --version; ./scripts/verify_ail.sh` | v0.30.0 commit `e37b370`; rc=0; 10/10 identities across 11 modules; 39 named tests; world package 9/9 with non-zero work; terminal `verify gate PASSED: 10 required identities verified, 39 named tests pass`. |
| V18 | The now-rejected round-1 six-constructor bare-ADT mapping checked and verified with pinned release. | `apply_patch` created `/tmp/iter84_proven_probe.ail`; `export AILANG_BIN=/tmp/ailang-v0300/ailang; AILANG_RELAX_MODULES=1 /tmp/ailang-v0300/ailang ai-check /tmp/iter84_proven_probe.ail > /tmp/iter84_proven_probe.json; python3 -c 'import json; d=json.load(open(...)); print(d["check"]["passed"],d["verify"]["verified"],d["verify"]["errors"],d["verify"]["counterexample"])'; /tmp/ailang-v0300/ailang test /tmp/iter84_proven_probe.ail` | `True 1 0 0`; `gradeCode_test_1` and `_2` pass. Temporary-path warning only. This is retained history, not evidence for the revised mapping. |
| V19 | A consistent contract/body lie evades Z3 but runtime policy tests red. | **Inherited controller measurement, iteration 81**: add the same `=> PROVEN` arm to contract and body, run pinned gate; compare Leg 1 and Leg 2. | Leg 1 `verify.verified=1`, errors/counterexamples `0`; six integer expectations make Leg 2 red. Not re-run because mutating live sources is outside this design-only task. |
| V20 | HashRef is nominal in Go with unexported fields and strict constructors. | `nl -ba host/hashref/hashref.go \| sed -n '43,125p'; rg -n 'func (New\|Parse\|Sum)' host/hashref/hashref.go` | `HashRef` fields `algo`, `digest` unexported; malformed/unsupported cases return `HashError`; constructor controls present. |
| V21 | No `EXACT_TOTAL_MODULES` variable exists; module inventory still has a positive control. | `{ rg -n 'EXACT_TOTAL_MODULES' scripts/verify_ail.sh || true; } \| wc -l; rg -n 'LEG1_MODULES' scripts/verify_ail.sh \| wc -l` | target `0`; same-file `LEG1_MODULES` control is non-zero (multiple lines). |
| V22 | **CONTROLLER-RUN (iteration 84), and it supplies the arm V18 lacks: V18's `verified=1` is NON-VACUOUS — the contract genuinely constrains the new `ValidatedProof` arm.** V18 ran only the positive arm, and a green there is equally consistent with "the contract binds" and "the contract is decorative". | Reproduced V18 first-party, then mutated **only the body** arm (`ValidatedProof(_) => CLAIMED`) while leaving the contract at `PROVEN`; asserted the mutant LANDED by sha256 (`3d7558d4…` → `3b125da3…`, `diff` = the single line 29); re-ran `ai-check` with the pinned binary, reading stdout and stderr to separate files (a warning on stdout otherwise voids the JSON parse). | Positive arm reproduces **without** `AILANG_RELAX_MODULES=1` as well as with it (`check.passed=true`, `verified=1`, `errors=0`, `counterexample=0`, `skipped=0`) — so V18's relax flag was **not** load-bearing and that narrowing need not travel with the finding. Mutant: **rc=1**, `verified=0`, `counterexample=1`, and Z3 names the witness exactly — `$p_e = (ValidatedProof (mk_HashRef "!0!" "!1!"))`. This does **not** contradict V19: mutating **one** side reds, mutating **both** stays green. Both facts hold, and §4 states both. |
| V23 | `world/types` is one of the four published package-module exports; the same exact set is pinned at every package gate surface. | `readonly EXPORTS=(world/types world/contracts world/transitions world/logepoch); rg -n 'readonly EXPORTS=\(world/types world/contracts world/transitions world/logepoch\)\|"exports": \{"modules":\["world/types", "world/contracts", "world/transitions", "world/logepoch"\]\}\|exact export set\|Modules:\[\]string\{"world/types","world/contracts","world/transitions","world/logepoch"\}\|world/types world/contracts world/transitions world/logepoch' scripts/verify_world_package.sh; printf 'control=%s\n' "${#EXPORTS[@]}"` | Matches at lines 34, 120, 153, 175, and 239; same-scope positive control `control=4`. Thus “four exports” means four exported modules, not four functions. |
| V24 | A foreign pure-AILANG consumer can import the rejected public constructor and execute `gradeOf(ValidatedProof(made-up HashRef)) => PROVEN`; a nonexistent-constructor negative control fails at import. | **CONTROLLER-RUN at `4557262` with pinned v0.30.0**: copied `world/types.ail` and `world/logepoch.ail` to scratch; applied the round-1 constructor/contract/body change and asserted three `ValidatedProof` occurrences; checked and tested `world/consumer.ail` importing `Evidence`, `EvidenceGrade`, `gradeOf`, `ValidatedProof`, and `HashRef`; in the same scratch scope replaced that import/call with `NoSuchCtor` and rechecked. | Positive: `ailang check` rc=0, `No errors found!`; `ailang test` passes `launder_code_test_1`, whose `PROVEN` arm returns 4. Negative control: rc=1, `Error: IMP010: symbol 'NoSuchCtor' not exported by 'world/types'`. No Go boundary, decoder, validator, or sealed value occurs on the positive path. |
| V25 | Revised `ProofReceipt => CLAIMED` syntax checks, verifies, and executes; a same-scope one-sided body mutation to `PROVEN` is non-vacuously rejected. | `apply_patch` created `/tmp/iter84_receipt_probe.ail` and identical `/tmp/iter84_receipt_mutant.ail` except the mutant body arm; `export AILANG_BIN=/tmp/ailang-v0300/ailang; $AILANG_BIN ai-check /tmp/iter84_receipt_probe.ail >/tmp/iter84_receipt_probe.stdout 2>/tmp/iter84_receipt_probe.stderr; $AILANG_BIN test /tmp/iter84_receipt_probe.ail >/tmp/iter84_receipt_test.stdout 2>/tmp/iter84_receipt_test.stderr; $AILANG_BIN ai-check /tmp/iter84_receipt_mutant.ail >/tmp/iter84_receipt_mutant.stdout 2>/tmp/iter84_receipt_mutant.stderr`; `python3` loaded each stdout JSON and printed `d['check']['passed'], d['verify']['verified'], d['verify']['errors'], d['verify']['counterexample']`; `rg 'PASS\|FAIL\|gradeCode\|counterexample\|ProofReceipt' /tmp/iter84_receipt_test.stdout /tmp/iter84_receipt_mutant.stdout /tmp/iter84_receipt_mutant.stderr` inspected runtime and negative-control outputs. | Positive rc=0: `check.passed=True`, `verify.verified=1`, `verify.errors=0`, `verify.counterexample=0`; both runtime tests pass, including made-up receipt → code 1. Mutant rc=1: `check.passed=True`, `verify.verified=0`, `verify.errors=0`, `verify.counterexample=1`; witness is `(ProofReceipt (mk_HashRef "!0!" "!1!"))`. |
| V26 | **CONTROLLER-RUN (iteration 84), and it is the arm that makes the round-2 fix a MEASUREMENT rather than a promise: the revised design DEFEATS the exact attack that defeated round 1, run in BOTH arms with only the design as the variable.** V24 and V25 each measure one design; neither runs the *attack* against the *fix*. | Rebuilt V24's scratch tree from `world/types.ail` + `world/logepoch.ail` at `4557262`, applied the **round-2** change instead (asserted LANDED: 3 `ProofReceipt` occurrences at lines 29/49/58), and re-ran V24's **byte-for-byte identical attack module** — a foreign `world/consumer` importing `gradeOf` and the receipt constructor, minting from the same literal `digest: "i-made-this-up"`, with an inline test asserting the result is `PROVEN` (code 4). | **Outcomes DIFFER, which is the whole evidence.** Round-1 arm: the attack test **PASSES** (`launder_code_test_1`, rc=0) — the foreign module really obtains `PROVEN`. Round-2 arm: `ailang check` still rc=0 `No errors found!` (so the refusal is *semantic*, not a type error the attacker would notice), and the attack test **FAILS** — `✗ attack_code_test_1`, **`test 0: expected 4, got 1`**, rc=1. `1` is `CLAIMED`. Same attacker, same literal, same instrument, same scratch scope; the only variable is the kernel arm. Note what this does NOT prove, per §2.3's own limitation: it bounds the `Evidence -> EvidenceGrade` route, not a future consumer that spells `EvidenceGrade.PROVEN` directly. |
| V27 | At `bef0153`, `verify.verified` is an integer, while `verify.results[]` supplies per-identity `function` and `status`; the known-positive `gradeOf` is present and verified. | `export AILANG_BIN=/tmp/ailang-v0300/ailang; $AILANG_BIN ai-check -timeout 5s world/types.ail > /tmp/iter87-types.json 2>/tmp/iter87-types.err; rc=$?; printf 'ai_check_rc=%s stderr_bytes=' "$rc"; wc -c < /tmp/iter87-types.err; python3 -c "import json; d=json.load(open('/tmp/iter87-types.json')); v=d['verify']; print(sorted(v.keys())); print(type(v['verified']).__name__,repr(v['verified'])); print(len(v['results'])); print(sorted(v['results'][0].keys())); [print(x['function'],x['status']) for x in v['results']]"` | rc=0, stderr 0 bytes; `verified` is `int 6`; results length 6 with keys `duration,function,status`; all six are `verified`, including `gradeOf`. This repairs the premise without weakening identity validation to a count. |
| V28 | At `bef0153`, the daemon has exactly eight registrations in `host/daemon/daemon.go` (seven GET and `POST /v1/commit`), while the scoped non-test `host`/`cmd` inventory has eight `PutObject` and 13 `GetObject` call sites; broker line 289 stores an effect-derived result object. | `rg -n 'HandleFunc\(' host/daemon/daemon.go; printf route_total=; rg -n 'HandleFunc\(' host/daemon/daemon.go \| wc -l; printf put_non_test=; rg -n '\.PutObject\(' host cmd --glob '*.go' --glob '!**/*_test.go' \| wc -l; printf get_non_test_control=; rg -n '\.GetObject\(' host cmd --glob '*.go' --glob '!**/*_test.go' \| wc -l; rg -n '\.PutObject\(' host cmd --glob '*.go' --glob '!**/*_test.go'; nl -ba host/broker/broker.go \| sed -n '280,294p'` | Registrations are lines 461–468, total 8; scoped counts are `put_non_test=8`, `get_non_test_control=13`; all eight put sites enumerate, and broker 288–290 forms `resultObject(result)` then stores it. The route enumeration contains no object-write route, but this same-scope positive store path means transport shape does not exclude writable-object threats. It does not demonstrate a complete report forgery. |
| V29 | Freshness sweep from the old declared base touches seven non-design files; the unexcluded full diff is a positive control with 18 paths. | `git diff --name-only 4557262..HEAD -- ':!design_docs'; printf 'control_total='; git diff --name-only 4557262..HEAD \| wc -l` | Seven paths: `docs/SELF_MOD_PUBLISH.md`, `host/verifygate/module_manifest_gate_test.go`, `packages/world-core/world/types.ail`, `scripts/verify_ail.sh`, `scripts/world_package_ready_packet.golden.json`, `sprint_w-decision-lifecycle-freeze.json`, `world/types.ail`; `control_total=18`. |
| V30 | In the scoped `Evidence` declaration at `bef0153`, exactly five constructor lines occur and none is a `PROVEN` producer; the adjacent grade declaration is the same-scope positive control for `PROVEN`. | `grep -n 'type Evidence' -A 12 world/types.ail` | Lines 24–28 enumerate `CompilerOutput`, `TestReport`, `HumanApproval`, `AiReview`, and `RecordedEffect`; line 34 contains the control `PROVEN`, while no constructor line does. |
| V31 | With pinned v0.30.0, a record parameter containing `list[Evidence]` fails encoding even when its contract reads only a scalar field; a bare `Evidence` parameter is the same-file, same-call positive control. The failure is silent to process rc and `check`. | `sed -n '1,220p' /tmp/iter87_adt_probe.ail`; `set +e; AILANG_RELAX_MODULES=1 /tmp/ailang-v0300/ailang ai-check -timeout 20s /tmp/iter87_adt_probe.ail >/tmp/iter87_adt_probe.rerun.json 2>/tmp/iter87_adt_probe.rerun.err; probe_rc=$?; printf 'rc=%s stderr_bytes=' "$probe_rc"; wc -c </tmp/iter87_adt_probe.rerun.err; python3 -c "import json; d=json.load(open('/tmp/iter87_adt_probe.rerun.json')); print(d['check']['passed'],d['check']['error_count'],d['verify']['verified'],d['verify']['errors']); [print(x['function'],x['status'],x.get('reason','')) for x in d['verify']['results']]"` | Probe defines `arm1(Rec)` with `evidence: list[Evidence]` but its contract/body read only `r.goal`; `arm2(Evidence)` is the control. `rc=0`, stderr `0` bytes, `check.passed=True`, `check.error_count=0`, `verify.verified=1`, `verify.errors=1`; `arm2 verified`, while `arm1 error` reports `Invalid constant declaration: unknown sort 'Rec'` and `unknown constant $p_r`. |
| V32 | The measured boundary is an ADT-typed field in a record parameter, bare or inside a list—not records or lists by themselves. The repository's Proposal comment is corroboration, not the measurement. | `sed -n '1,260p' /tmp/iter87_adt_probe2.ail`; `set +e; AILANG_RELAX_MODULES=1 /tmp/ailang-v0300/ailang ai-check -timeout 20s /tmp/iter87_adt_probe2.ail >/tmp/iter87_adt_probe2.rerun.json 2>/tmp/iter87_adt_probe2.rerun.err; probe_rc=$?; printf 'rc=%s stderr_bytes=' "$probe_rc"; wc -c </tmp/iter87_adt_probe2.rerun.err; python3 -c "import json; d=json.load(open('/tmp/iter87_adt_probe2.rerun.json')); print(d['check']['passed'],d['check']['error_count'],d['verify']['verified'],d['verify']['errors']); [print(x['function'],x['status'],x.get('reason','')) for x in d['verify']['results']]"; nl -ba world/contracts.ail \| sed -n '9,16p'` | Same file/call: scalar-only record `arm3` and record with `list[string]` `arm6` are `verified`; record with bare `Evidence` `arm4` and record with `list[Evidence]` `arm5` are `error`. `rc=0`, stderr `0` bytes, `check.passed=True`, `check.error_count=0`, `verify.verified=2`, `verify.errors=2`; errors name unknown sorts `BareAdtRec` and `ListAdtRec`. Repository lines 11–13 independently note Proposal's unknown-sort failure and upstream `sunholo-data/ailang#477`. No deeper nesting, result-position record, or tuple shape was tested. |
| V33 | Round-5 sweep (R1): before this revision the document referenced the free resolver symbol at exactly **five** lines — 78, 90, 201, 295, 398 — two more than the three the revision brief listed, which is why the sweep and not the list is authoritative. After the revision, every remaining occurrence names the symbol only to record that it is dropped (header, §2.2, §3.2, §10.5, and this log); none specifies it as a live API surface. | `grep -n 'GradeOfValidated' design_docs/planned/w-validated-proven-evidence-boundary.md` run before the edits and again after; counts via `grep -c` on the same pattern and path | Before: `doc_count=5`, lines 78/90/201/295/398 (the line-78/201/295 subset matching the brief's `sed -n '78p;201p;295p'` spot-check). After the R1/R2/R3 edits and before these two rows landed: 5 hits at lines 9/82/221/705/710, each a drop-notice; final re-run after this table row is recorded beneath the table. |
| V34 | At `03c7892`, the dropped symbol appears in zero `.go`/`.ail` files and `host/evidence` does not exist, so R1 is a document-only change with no code to migrate; the same-call positive control fires. | `printf 'code_hits='; { grep -rln 'GradeOfValidated' --include='*.go' --include='*.ail' . \|\| true; } \| wc -l; printf 'control='; grep -rln 'func (s \*Store)' --include='*.go' . \| wc -l; ls -d host/evidence \|\| echo ABSENT; ls host/ \| wc -l` | `code_hits=0`; same-instrument control `4` files define `func (s *Store)` methods; `host/evidence` prints ABSENT while `ls host/` lists 15 packages, so the instrument sees both positives and the directory root. |
| V35 | At `03c7892`, no exported bounded-subprocess helper exists in non-test `host`/`cmd` for `host/evidence` to import; the same-call inspection of unexported `runBounded` is the positive control and the measured pattern §3.4 reimplements (wall-time bound, output cap with overflow detection, overrun kill → refusal). | `git rev-parse --short HEAD; rg -n 'CommandContext' host cmd --glob '*.go' --glob '!**/*_test.go'; rg -n 'func Run[A-Z]\|func .*Bounded' host cmd --glob '*.go' --glob '!**/*_test.go'; nl -ba host/broker/handlers.go \| sed -n '60,140p'` | `03c7892`. Exactly 3 non-test `CommandContext` sites: `host/replay/replay.go:327`, `host/capsule/capsule.go:154`, `host/broker/handlers.go:93`. The exported-helper pattern matches only `func runBounded` at `host/broker/handlers.go:88` — 1 hit, lowercase, so unexported and broker-internal; zero `func Run[A-Z]` subprocess helpers. Inspection shows `context.WithTimeout(ctx, bounds.execTimeout)` feeding `exec.CommandContext`; `Setpgid: true` with `cmd.Cancel` invoking `killGroup` (`syscall.Kill(-pgid, SIGKILL)`) so a forked grandchild dies with the group; `io.LimitedReader{N: bounds.maxOutputBytes + 1}` with overflow returning `HandlerOutputOverflowError` and deadline expiry returning `HandlerTimeoutError`, both refusals returning no output. One property NOT to copy: `cmd.Stderr = cmd.Stdout` merges the streams, which §3.4's producer must not do because it parses stdout as JSON. |
| V36 | At `52bc9ec`, 6b's premise still holds and item 18's landing did NOT close it: non-test `host/store` exposes exactly two exported `Object` methods, `GetObject` returns the whole payload as `[]byte`, and the six non-test files contain zero `io.Reader`/`io.LimitReader`/`maxBytes` and never touch the `io` package — item 18 added `context.Context`, not a byte bound. The stale round-6 control is corrected: **25** exported `Store` methods, not 23. **Round-7b correction: this row's claim originally called the added parameter "(a WAIT bound)". It is not one — a context parameter bounds nothing unless armed with a deadline, and 11 production store reads pass `context.Background()` under item 18's own ratified DR-2 deferral (V41). Item 18 bounded the SIGNATURE; the measurements in this row are unchanged (§10.8).** | `git rev-parse --short HEAD; grep -n 'func (s \*Store) PutObject\|func (s \*Store) GetObject' host/store/store.go; F=(host/store/journal.go host/store/scan.go host/store/store.go host/store/writer_lock.go host/store/writer_lock_other.go host/store/writer_lock_unix.go); grep -nE 'io\.Reader\|io\.LimitReader\|maxBytes' "${F[@]}"; echo "io_rc=$?"; grep -nE '^\s*"io"\|\bio\.' "${F[@]}"; echo "ioimport_rc=$?"; for f in "${F[@]}"; do printf '%s ctx=%s\n' "$f" "$(grep -cE 'context\.Context' "$f")"; done; grep -hnE '^func \([a-z]+ \*Store\) [A-Z]' host/store/*.go \| grep -v _test \| wc -l` | `52bc9ec`; `store.go:444 func (s *Store) PutObject(o Object) error`, `store.go:468 func (s *Store) GetObject(ctx context.Context, ref hashref.HashRef) (Object, bool, error)`; `io_rc=1`, `ioimport_rc=1` (zero hits across all six files); same-path control: `context.Context` reads `store.go 7`, the other five files `0`; **25** exported `Store` methods across the six non-test files (`journal.go` 9, `scan.go` 2, `store.go` 14, the three `writer_lock*` files 0). |
| V37 | At `52bc9ec` the AILANG gate pins are unchanged in VALUE and moved in LINE: V5's cited 311/350 are `bef0153` positions; item 18's landing moved them. The doc's §3.6/§5/§8.1 claims about the pins remain correct. | `grep -nE "^EXACT_TOTAL_VERIFIED=\|^EXACT_TOTAL_TESTS = " scripts/verify_ail.sh` | `323:EXACT_TOTAL_VERIFIED=10`; `363:EXACT_TOTAL_TESTS = 39`. |
| V38 | Round-7 freshness sweep from the EARLIEST declared base: 41 non-design files changed `aaada20..52bc9ec`, and the same-instrument control scoped to one cited file fires. Every retained row citing a swept file was re-run or explicitly labelled unrefreshed (see the preamble); none was marked fresh by default. | `git diff --name-only aaada20..HEAD -- ':!design_docs' \| wc -l; git diff --name-only aaada20..HEAD -- host/store/store.go \| wc -l` | `41`; control `1`. |
| V39 | **CONTROLLER-MEASURED at `52bc9ec`, round 7, after the designer run — PRIOR ART FOR THE BOUNDED-READ IDIOM EXISTS AND THE DOC DID NOT KNOW IT.** Non-test `host/` already contains **two** `io.LimitReader` sites, and one of them (`host/capsule/capsule.go:238`) already uses the **`limit+1` detection-byte** shape §3.3 step 3 now specifies. Neither is an exported bounded-read helper — both are inline stdlib idioms at their own call sites — so there is nothing to REUSE and §8.2's new `host/store` surface is not duplicating an existing API. Recorded because §3.4 applies exactly this discipline to `runBounded` (V35: reuse REJECTED with a reason rather than deferred), and applying it to `runBounded` but not to the idiom this round introduces would be the same inconsistency round 6 blocked on. AC3's scoping is unaffected: it counts `io.LimitReader` in non-test `host/evidence`, a package that does not exist at base, so these two sites cannot satisfy it. | `grep -rn 'io.LimitReader' host/ --include='*.go' \| grep -v _test` ; same-path control `grep -rn 'io.ReadAll' host/ --include='*.go' \| grep -v _test \| wc -l` | Two hits: `host/broker/registry_reconcile.go:481` `io.ReadAll(io.LimitReader(resp.Body, cfg.MaxBodyBytes))` and `host/capsule/capsule.go:238` `io.ReadAll(io.LimitReader(pipe, limit+1))`; `host/evidence` does not exist (`test -d host/evidence` rc=1). |
| V40 | At `7806cac`, `GetObject` cannot gain a `maxBytes` parameter without changing out-of-tranche callers: exactly **13** non-test `.GetObject(` call sites exist in `*.go` under `host/` and `cmd/` outside `host/store` — `host/broker/approve.go` 6, `host/broker/broker.go` 2, `host/transitionreg/transitionreg.go` 2, `host/registry/registry.go` 1, `host/daemon/handlers.go` 1, `host/replay/replay.go` 1 — and the last two are in packages §8.2 declares frozen, so the round-7 `maxBytes`-on-`GetObject` arm was unsatisfiable by construction and is deleted (§10.8). | `git rev-parse --short HEAD; grep -rn '\.GetObject(' --include='*.go' host/ cmd/ \| grep -v '^host/store/' \| grep -v _test \| wc -l`; same-call, same-scope control: `grep -rn '\.PutObject(' --include='*.go' host/ cmd/ \| grep -v '^host/store/' \| grep -v _test \| wc -l`; per-file enumeration via `grep -rln` then `grep -c` per file | `7806cac`; count `13`; control `8` — non-zero, and equal to V28's independently measured eight non-test `PutObject` sites, so the instrument reads a known quantity correctly. Per-file enumeration exactly as stated in the claim. |
| V41 | At `7806cac`, a context parameter is not a wait bound — item 18's OWN ratchet says so: `host/store/context_read_test.go` declares **11** production store reads that pass `context.Background()`, left deadline-free "by ratified deferral DR-2", with the comment stating the set "may SHRINK … but it may never GROW" and naming the "11 -> 0" reduction as the FOLLOW-ON item's mechanically observable path. So item 18 bounded the SIGNATURE, and this tranche's validator must arm its own deadline (`ObjectReadTimeout`, AC16) rather than inherit one; the 11 sites belong to item 18's follow-on, not to this tranche (§8.2). | `grep -n 'deadlineFreeReadPins\|TestNoNewDeadlineFreeStoreReads\|DR-2' host/store/context_read_test.go; sed -n '361,379p' host/store/context_read_test.go` | Pin map at `context_read_test.go:370`: `host/broker/approve.go` 8, `host/registry/registry.go` 2, `host/replay/replay.go` 1 — sum 11; matcher regexp `\.(GetObject\|GetWorld\|GetLogEntry\|GetRegistryHead\|SelectedHead)\(\s*context\.Background\(\)`; `TestNoNewDeadlineFreeStoreReads` at `:387` scans non-test `host/` and `cmd/` and reds if any file's count exceeds its pin. The comment reads verbatim "the EXACT set of production store reads this item leaves deadline-free, by ratified deferral DR-2". |
| V42 | AC16's and AC3's base readings at `7806cac`: the token `ObjectReadTimeout` occurs in **zero** `*.go` files under `host/`, `cmd/`, and `scripts/`, and `OpenObject` in **zero** `*.go` files under `host/` and `cmd/` (tests included), so both head readings differ from base by construction; the same-call positive control fires in the same scope. | `printf 'objectreadtimeout_code='; grep -rn 'ObjectReadTimeout' --include='*.go' host/ cmd/ scripts/ \| wc -l; printf 'withtimeout_control='; grep -rn 'context.WithTimeout' --include='*.go' host/ cmd/ \| grep -v _test \| wc -l; printf 'openobject_code='; grep -rn 'OpenObject' --include='*.go' host/ cmd/ \| wc -l` | `objectreadtimeout_code=0`; `openobject_code=0`; same-scope control `withtimeout_control=8` non-test `context.WithTimeout` sites (all outside the not-yet-existing `host/evidence`), so each zero is a measurement, not a broken pattern. |

Final V33 re-run after all round-5 edits including revision 2: `grep -n 'GradeOfValidated'` on
this file returns **8** hits at lines 9 (header), 85 (§2.2 drop-notice), 231 (§3.2 drop-notice),
786/791 (§10.5 records), and 882/883/886 (the V33/V34 rows and this sentence, which quote the
pattern as an instrument). Zero of the eight is a live API reference; the normative API surfaces
in §2.2, §3.2, and §3.5 name only `Validator.Resolve`.

Round-7 re-run of the same grep after the §10.7 edits: still **8** hits — the round-5 line
numbers above are that round's positions and have since shifted — at lines 9 (header), 93
(§2.2 drop-notice), 240 (§3.2 drop-notice), 850/855 (§10.5 records), 1077/1078 (the V33/V34
rows), and the preceding note's own grep quotation; all remain drop-notices, records, or
instrument quotes. The retired single-return signature survives only at three sites, every one a record —
the §10.5 and §10.6 histories and §10.7's account of replacing it — and the normative surfaces
name only the sum-style `Validator.Resolve(sealed ValidatedEvidence) ResolutionResult`.

## 12. Related

- `design_docs/implemented/w-evidence-grade-mapping.md` — predecessor and binding deferred
  authority obligations.
- `design_docs/HUMAN-SURFACE.md` — trust gradient and grade-laundering prohibition.
- `design_docs/coding-standards.md` — pure core, effect boundary, slim kernel, and honest gates.
- `design_docs/DESIGN.md` §§1, 14 — immutable world and controlled-change architecture.
- Planned successor: `w-proven-evidence-production-key-wiring`.
- Planned successor: `w-validated-replay-evidence-boundary`.
- Planned successor: `w-proven-evidence-renderer-consumption`.
