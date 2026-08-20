# w-validated-proven-evidence-boundary — authority-bearing proof evidence

- **Status**: **DIRECTION RATIFIED (`D-WORLD-17`); round-7 revision applied (`D-WORLD-19` arm A); round-7b carve-out revision applied; round-9 revision applied (`D-WORLD-21` arm A, attended 2026-08-19); round-9b carve-out revision applied; round-11 revision applied (`D-WORLD-22` arm B via the `D-WORLD-23` arm-A standing rule, Mark attended 2026-08-20, §10.13); round-11b carve-out revision applied (both round-11 objections, reviewers' verbatim fixes, §10.14)** — **round-10 quorum: `gemini-3-1-pro` PASS (2nd consecutive), `gpt5-6-sol` reject on the DECLARED residual; the round-10 park is ANSWERED — tranche 1 keeps scope, the wait-bound claim narrows to exactly what is proven, OPEN queue row 22 keeps the lock-wait residual, and `busy_timeout < ObjectReadTimeout` is pinned at construction (AC18)** — **round-11 quorum: BLOCKED, both present, on two NEW completeness defects in this round's own remedy; the narrow-refinement CARVE-OUT applied for the first time on this document (§10.14)** — **round-12 confirming quorum: BLOCKED, both present, and for the THIRD consecutive round the objections land on the PREVIOUS round's own fix; PARKED `needs-human-review` on `D-WORLD-24`, a one-word A/B on whether tranche 1 SHEDS the producer (§10.15)** — DESIGNED, not landed
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
  §10.8. Round 9 applies the attended `D-WORLD-21` arm-A ruling of 2026-08-19: the round-7b
  `OpenObject(ctx, ref, maxBytes)` streaming spelling is RETIRED for the complete-read
  `ReadObject(ctx, ref, maxBytes) (ObjectMeta, []byte, error)` — the store enforces `maxBytes`
  BEFORE materialization and performs the whole read under the supplied context (round 11
  narrows round 9's "so cancellation is enforceable" to the three-part claim below) — M22 is
  REWRITTEN against the store's actual cancellation mechanism rather
  than re-run, and `gpt5-6-sol`'s owed real-store integration test lands in AC16 with mutation
  rows M22/M23. See §10.10. Round 9b applies both round-9 objections under the
  narrow-refinement carve-out: `ReadObject`'s probe and payload statements now run on ONE
  reserved database connection inside ONE read transaction/snapshot, with the previously
  unverified insert-only premise VERIFIED by positive enumeration (V47, including its
  instrument finding) and the concurrent-mutation integration test landing as AC17 with
  mutation row M25; and §3.2 gains the producer's explicit Go API (`NewProducer`,
  `Producer.GenerateProof`) reconciled with §3.4's bounds. See §10.11. Round 11 applies
  `D-WORLD-22` arm B, resolved by the `D-WORLD-23` arm-A standing rule (Mark, attended,
  2026-08-20T08:01:31Z, `#68` comment `A`): the wait-bound claim narrows to exactly what is
  proven — (1) every wait this tranche's OWN code performs is bounded by `ObjectReadTimeout`;
  (2) a LOCK-contended wait is bounded by `busy_timeout` (2000 ms,
  `host/store/writer_lock.go:179`), which this tranche does not change; (3) the composition is
  safe ONLY WHILE `busy_timeout < ObjectReadTimeout` — the residual is owned by OPEN queue row
  22 `w-daemon-lock-wait-not-deadline-bound` (verified open by command, §10.13), and the
  ordering is PINNED at construction (`ErrUnorderedTimeouts`, AC18, M26, V49). The owed
  round-10 nesting-depth note is measured and discharged with the same remedy (AC19, M27,
  V50). See §10.13.
- **Item**: queue item 17, `w-validated-proven-evidence-boundary`
- **Filed**: 2026-08-14, iteration 84
- **Measurement base**: `bef0153` (rounds ≤ 4); `03c7892` (round 5); `52bc9ec` (round 7, V36–V38); `7806cac` (round 7b, V40–V42); `35fd875` (rounds 9 and 9b, V43–V48); `516836f` (round 11, V49–V50)
- **Instrument**: `/tmp/ailang-v0300/ailang`, AILANG v0.30.0 (`e37b370`)
- **This tranche estimate**: **5.35 days** (round 7 adds the ratified bounded store read; round
  7b adds the `ObjectReadTimeout` deadline machinery; round 9 nets +0.25 day — the streaming
  reader's removal against the real-store integration test's addition; round 9b adds 0.25 day
  for the one-snapshot read transaction and its concurrent-mutation test; round 11 adds 0.40
  day for the timeout-ordering and nesting-depth pins and round 11b a further 0.25 day for the
  cached-accessor amendment and the bounded envelope write — each priced, not absorbed, §9)
- **Decomposition**: **yes** — four ordered sprint-sized items, **11.35 days total** (§9)
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
ObjectReader                 // minimal read seam the validator consumes; its one method is the ruled D-WORLD-21 arm-A signature (§8.2, §10.10): ReadObject(ctx context.Context, ref hashref.HashRef, maxBytes int64) (ObjectMeta, []byte, error)
NewValidator(key [32]byte, reader, compilerConfig)   // configuration REQUIRES a positive ObjectReadTimeout; zero or negative is a constructor refusal (AC16); and when the supplied reader reports its lock-retry window — the real store does, via the additive store.BusyTimeout() accessor, a NONBLOCKING CACHED PROPERTY populated at Open from the resolved DSN and never queried through the pool (§8.2; gpt5-6-sol's round-11 alternative fix, applied verbatim under the carve-out, §10.14) — a POSITIVE ObjectReadTimeout <= that window is ALSO a constructor refusal, the dedicated ErrUnorderedTimeouts (AC18, round 11): the busy_timeout < ObjectReadTimeout composition condition is pinned against the two values actually used at run time, never two literals
Validator.ValidateProof(ctx, reportRef, expectedSubject)   // derives context.WithTimeout(ctx, ObjectReadTimeout) before the object read, even under context.Background() (§3.3 step 2, AC16)
DecodeProposal(raw)          // ProofReceipt remains an untrusted claim; raw capped at 256 KiB before any parse (§3.3)
Validator.Resolve(sealed ValidatedEvidence) ResolutionResult   // method; sole ResolvedGradeProven source; refuses zero/unminted identities with ErrUnmintedAuthority, foreign seals with ErrForeignSeal
ResolutionResult             // unexported fields; mutually exclusive Proven() (ResolvedGrade, bool) and Err() error; zero value is refusal
ResolvedGrade                // host enum; ResolvedGradeProven is not serialized
ObjectWriter                 // minimal write seam the producer stores envelopes through; its one method TAKES A CONTEXT — WriteObject(ctx context.Context, o store.Object) error (gemini-3-1-pro's round-11 fix, applied verbatim under the carve-out, §10.14): the round-11 spelling was "one PutObject-shaped method", and the real store's PutObject takes NO context (V43, re-verified V51), so the producer's envelope write was unbounded against the same one-connection pool AC16's read bound is built around. §3.4's "stores an authenticated envelope" requires a writer, and the validator keeps the read-only ObjectReader
NewProducer(key [32]byte, checker, writer, execTimeout time.Duration, maxOutputBytes int64)   // §3.4's producer, configured at construction (gemini-3-1-pro's round-9 fix; departures from the verbatim spelling are enumerated in §10.11): key is the same caller-owned 32-byte MAC key NewValidator takes; checker pins the absolute AILANG_BIN path plus the exact executable hash and "AILANG v0.30.0" token §3.4 verifies before any run; writer is the ObjectWriter above; execTimeout bounds wall time via context.WithTimeout feeding exec.CommandContext (§3.4); maxOutputBytes is the PER-STREAM cap, applied independently to stdout and stderr through §3.4's limit+1 capped readers; non-positive execTimeout or maxOutputBytes is a constructor refusal (AC5); Producer.GenerateProof PASSES ITS OWN CONTEXT to the writer's WriteObject, so the envelope write is wait-bounded by the caller's deadline rather than by nothing (round 11b, AC20, M28) — the execution bounds are explicit and checkable at the API, per the bounded-waits axiom, not ambient configuration
Producer.GenerateProof(ctx context.Context, sourceRef HashRef, requiredIdentities []string) (HashRef, error)   // the reviewer's round-9 signature, verbatim: resolves sourceRef through the ObjectReader above (ReadObject under the same 256 KiB bound — an over-bound subject module is refused, not proven), materializes the module bytes to a private temp file for the pinned checker, runs it under the configured bounds, and builds the report's verified set only from verify.results[].function entries with status == "verified" (V27); empty requiredIdentities is a refusal (§3.4); returns the stored authenticated envelope's HashRef only on full success — deadline overrun and either-stream overflow are refusals that emit and store nothing
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

The ordering pin's seam is explicit, and the cross-package decision is written down rather than
left implicit (round 11, `D-WORLD-23` obligation (ii)). `ObjectReader` keeps its ONE required
method; the pin rides an optional capability interface, unexported in `host/evidence`:
`busyWindowReporter`, one method, `BusyTimeout() time.Duration`. `NewValidator` type-asserts
the supplied reader against it. A reader that reports no lock-retry window — a bounded test
fake with no lock layer — skips the check: there is no window to order against, and a fake
reporting one it does not have would supply the PROPERTY under test, §6.1's inadmissible
class. The REAL `host/store.Store` does report one, through the additive exported accessor
`BusyTimeout() time.Duration` (§8.2). Round 11 specified that accessor as a LIVE
`PRAGMA busy_timeout` read, and `gpt5-6-sol` rejected it for the reason this whole round
exists to close: with `db.SetMaxOpenConns(1)` (`store.go:298`, V51) a context-free pragma query
can wait unboundedly for the sole pooled connection — the round's own remedy reintroducing the
round's own defect. Its alternative fix is applied verbatim (§10.14): `BusyTimeout()` is a
**genuinely nonblocking cached property**, resolved once during `Open` from the DSN the store
was opened with and thereafter a field read — no connection acquisition, no query, no error
channel. The "LIVE PRAGMA" claim of round 11 is WITHDRAWN as unsupported, and nothing is lost
by withdrawing it: the window is IMMUTABLE for the store's lifetime, because it is applied per
physical connection from the `_pragma` DSN parameter at open time (`writer_lock.go:196`) and
non-test `host/`+`cmd/` code issues **zero** runtime `PRAGMA busy_timeout` writes (V51, with a
firing control on the DSN site). A caller-supplied DSN that set its own busy window — which
`withBusyTimeout` deliberately never overrides — is therefore reported at its true runtime
value by the cached read exactly as it would have been by a live one, so the pin still binds
the window SQLite will actually apply to a lock wait against the timeout the validator will
actually arm. This is strictly stronger than bounding the query would have been: the wait is
REMOVED rather than capped.
`context_read_test.go:209` already pins the constant's VALUE at 2000; nothing anywhere pins
the ORDERING against any consumer's deadline (V49), and `host/daemon/handlers.go:299–302`
names that gap in the codebase's own words. The refusal therefore lives in `host/evidence`'s
constructor — the consumer that owns a deadline owns the refusal; the store cannot know any
consumer's deadline and contributes only the reporter. AC18 makes the refusal non-vacuous and
M26 must red it.

The new package depends on `host/hashref` and the minimal object reader and writer seams above
— the writer existing solely for the producer's envelope store — not on daemon, replay, or
renderer. That avoids a cycle and makes the validator testable with a bounded fake.
`NewValidator` copies an exactly 32-byte MAC key supplied directly by its caller; tranche 1 accepts
no key path and neither creates nor loads key files. The proof producer receives the same
caller-owned key material through its library constructor, `NewProducer` above. Consequently tranche 1 is explicitly
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
point of I/O: step 2 below reads the object through the store's `ReadObject(readCtx, ref,
262144)`, which refuses any payload whose PROBED length exceeds the bound BEFORE
materialization — the probe statement's select list does not contain the payload column, so no
payload byte crosses the driver into Go before the guard (§8.2) — and an attacker-supplied
`HashRef` to a multi-gigabyte object therefore never materializes at all. Round 6's objection
6b stays closed with a stronger property than the round-7b stream gave it: there is no
`io.ReadCloser`, no caller-side `io.LimitReader`, no detection byte, and no reader lifetime to
close on the envelope path (the streaming spelling is retired under `D-WORLD-21` arm A,
§10.10). The decoded report bytes and verified list are separately capped at 256 KiB and 256
unique identities, and every string at 1 KiB. `DecodeProposal` applies the same 256 KiB cap to
its raw input before any parse, refusing oversize input undecoded — the iteration-87
non-blocking note, discharged. The chosen raw cap follows an existing 256 KiB strict-codec
bound already used for registry schema input (V11); it is a new Evidence invariant, not
inherited store behaviour.

Nesting depth within the byte cap is MEASURED, not assumed (round 11, discharging
`gemini-3-1-pro`'s round-10 non-blocking note — V50: four arms, both toolchains go1.25.6 and
go1.26.4, payloads all inside the 256 KiB cap). REFUTED: no stack-overflow panic and no
pathological CPU — a pure `[[[[…` bomb at the deepest depth that fits the cap (131,071 levels,
262,142 bytes) is REFUSED by Go's JSON scanner in under a millisecond, and the deepest
ACCEPTED payload (9,999 levels) costs single-digit milliseconds, bounded by construction
because depth cannot exceed half the byte bound. SUSTAINED, and it is the half that matters:
the refusing limit is `maxNestingDepth = 10000` at `encoding/json/scanner.go:148` of the
pinned toolchain's GOROOT — an UNEXPORTED stdlib implementation constant, not a documented
`encoding/json` API guarantee — so the resilience just measured rests on an internal nothing
asserts: the SAME defect shape as the `busy_timeout < ObjectReadTimeout` ordering, and this
round's spine (§10.13). Same remedy: no depth guard of our own is added — inside a 256 KiB cap
the stdlib limit is sufficient in VALUE — but the behaviour is PINNED. AC19's named test
`TestNestingDepthBombWithinByteCapIsRefused` feeds `DecodeProposal` a depth-131,071 payload
inside the byte cap and requires the typed decode refusal and no `ClaimedEvidence`, with a
shallow-depth control in the same test that must decode — so a future toolchain that raises or
removes the internal limit reds the pin instead of silently invalidating the measurement. M27
must red it.

Validation order is fixed:

1. validate canonical `HashRef`;
2. derive the read deadline — `readCtx` from `context.WithTimeout(ctx, ObjectReadTimeout)` —
   even when the caller supplies no deadline, `context.Background()` included (§3.2, AC16);
   then read exactly one object through `ObjectReader`'s
   `ReadObject(readCtx, ref, maxBytes)` (§8.2): the store probes `length(payload)` and refuses
   an over-bound object with the typed `*store.ObjectTooLargeError` BEFORE any payload
   statement runs, and performs the complete read under `readCtx`; deadline expiry or
   cancellation is an explicit operational error that mints no seal (§2.5). The wait-bound
   claim is scoped to exactly what is proven (`D-WORLD-22` arm B, §10.13): `readCtx` bounds
   every wait THIS TRANCHE'S OWN code performs — the connection wait and the in-flight
   statements; a LOCK-contended wait inside the store is bounded by `busy_timeout` (2000 ms,
   `host/store/writer_lock.go:179`), which this tranche does not change, so the composition
   is safe ONLY WHILE `busy_timeout < ObjectReadTimeout` — an ordering `NewValidator` refuses
   to violate (AC18, M26) and whose completion is owned by OPEN queue row 22
   `w-daemon-lock-wait-not-deadline-bound`;
3. reject absent objects as `missing`, and map the store's typed over-bound refusal to
   `oversize` — the envelope byte bound has exactly ONE owner, `ReadObject`'s `maxBytes`; the
   validator performs no second envelope-length check, because a duplicate check would make the
   store guard's mutation (M4) observable-identical and therefore unkillable;
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
bytes; the `ObjectReadTimeout`-derived deadline in step 2 closes the waits this tranche's OWN
code performs at its OWN call site, inherited from no one — the LOCK-contended wait stays
`busy_timeout`-governed and row-22-owned, safe only under AC18's pinned ordering (§10.13). The
11 DR-2 sites are item 18's follow-on obligation and are
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
The store step is itself wait-bounded (round 11b, `gemini-3-1-pro`'s fix applied verbatim,
§10.14): `GenerateProof` passes its own context to `ObjectWriter.WriteObject(ctx, o)`, whose
real implementation performs the connection reservation and the insert under that context —
analogous to `ReadObject` (§8.2). Round 11 bounded the read side and the subprocess and left
the write side inheriting nothing, against the same one-connection pool (`store.go:298`, V51);
the existing `Store.PutObject(o Object) error` takes no context at all (V43, re-verified V51)
and is left untouched for its **8** existing call sites outside `host/store` (V51; the
controller's first draft said 10, transcribed from an earlier epoch rather than measured — the
same-pipeline `GetObject` control returns 13, reproducing V43 exactly, so the 8 is a
measurement). AC20 pins the bound and M28 kills it.

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
| M4 **oversize** | `host/store/store.go`: change `ReadObject`'s pre-materialization guard `if n > maxBytes` (n is the probed `length(payload)`) to `if false && n > maxBytes` | `TestOversizeProofReportIsRefused`, feeding an OTHERWISE-VALID > 256 KiB envelope through the REAL store — base64's 4/3 inflation makes one constructible: a ~250 KiB correctly-MAC'd report yields a ~333 KiB envelope, over the read bound while the decoded report stays inside its own 256 KiB cap | Correct code: the probe refuses before any payload statement runs; reason `oversize`. Under the mutant the payload statement materializes the full envelope, step 4's recomputed hash MATCHES (the object is genuinely content-addressed), decode and MAC succeed, and the oversize envelope SEALS; failure: `oversize envelope sealed; want oversize`. The kill needs the real store: the mutated line is in `host/store`, which any fake replaces. |
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
| M22 **read detached from context** | `host/store/store.go`: in `ReadObject`, replace the supplied `ctx` with `context.Background()` at EVERY context-taking call — the connection reservation, the transaction begin, and both statement sites (round 9b: the one-snapshot spelling moved the connection WAIT — the very wait AC16's stimulus blocks on — out of `QueryRowContext` and into the reservation; a mutant detaching only the two statement sites would leave that wait context-bounded and the row would survive VACUOUSLY, §10.11); one mechanism (the read detaches from the caller's deadline, exactly the caller-inherited deadline-free state of item 18's DR-2 residue, V41) | `TestRealStoreBlockedObjectReadReturnsWithinObjectReadTimeout` (AC16) | Correct code: the blocked attempt's connection wait observes `readCtx`, `ValidateProof` returns the explicit operational error within the bound, and no seal exists. Under the mutant the wait's context is never done, so the attempt outlives the decoy hold, completes, and seals; failure: `blocked read sealed after exceeding ObjectReadTimeout; want operational timeout error and no seal` (the test-side 20× watchdog reds the hang mode). The kill NEEDS the real store — the mutated lines are in `host/store`, code any fake replaces — which is the round-8 fixture vacuity inverted into the fix. |
| M23 **deadline arming** | `host/evidence/validator.go`: replace the read-deadline derivation `context.WithTimeout(ctx, cfg.ObjectReadTimeout)` with `context.WithCancel(ctx)` (the retired round-7b M22 edit, now killed honestly: the mutant compiles, keeps a cancel, and never arms a deadline — the caller-inherited state AC16 forbids) | `TestRealStoreBlockedObjectReadReturnsWithinObjectReadTimeout` (AC16) | Under the mutant `readCtx` carries no deadline, so nothing expires the real blocked read's connection wait; the attempt outlives the decoy hold, completes, and seals — the same observable as M22 through a DIFFERENT mutated line, in a different package. Failure text as M22. No fake-reader test participates in this kill; the round-7b fake killer is deleted (§6.1, §10.10). |
| M24 **oversize translation** | `host/evidence/validator.go`: change the typed-refusal mapping `if errors.As(err, &tooLarge)` to `if false && errors.As(err, &tooLarge)` | `TestOversizeProofReportIsRefused` | The store's `*store.ObjectTooLargeError` falls through to the generic operational-error branch instead of the `oversize` reason; failure: `got operational store-read error; want oversize`. Paired with M4: M4 proves the store guard is load-bearing, M24 proves the validator's reason mapping is — the byte-slice seam separates what the round-7b stream unified, so §5's every-refusal-branch rule demands both rows. |
| M25 **one-snapshot probe/payload** | `host/store/store.go`: remove the read transaction from `ReadObject` — run the probe and the payload statement as two bare autocommit statements on the same reserved connection (the pre-round-9b spelling, restored as a mutant) | `TestConcurrentMutationCannotDesyncProbeAndPayload` (AC17) | Correct code: the transaction's snapshot pins the probed row, so the returned payload matches the probed length and re-hashes to the requested ref whichever way the injected writer lands (committed after the snapshot, or refused busy — AC17 accepts both). Under the mutant the hook-injected `UPDATE` lands between the two autocommit statements and the payload statement reads the NEW row; failure: `payload diverged from probed row: probe=100 bytes, payload=300 bytes; want one snapshot`. The kill is DETERMINISTIC — the scheduling hook forces the interleaving, so no retry classification participates — and NEEDS the real store: the mutated lines are `host/store` code any fake replaces. |
| M26 **timeout ordering** | `host/evidence/validator.go`: change `NewValidator`'s ordering refusal `if window > 0 && cfg.ObjectReadTimeout <= window` (window is the reader's reported `BusyTimeout()`; exact spelling fixed at implementation — the neutered comparison is the ordering pin itself) to `if false && (window > 0 && cfg.ObjectReadTimeout <= window)` | `TestConstructorPinsBusyTimeoutBelowObjectReadTimeout` (AC18) | Correct code: constructing a validator over the REAL file-backed store with a POSITIVE `ObjectReadTimeout` of 1 s — merely <= the store's live 2000 ms busy window — is refused with the dedicated `ErrUnorderedTimeouts` naming both values. Under the mutant the unordered validator constructs; AC16's non-positive constructor refusal is a DIFFERENT, WEAKER guard and stays green on this stimulus (1 s is positive), which is exactly what makes the two arms distinguishable rather than one check wearing two names. Failure: `NewValidator accepted ObjectReadTimeout 1s <= busy_timeout 2s; want ErrUnorderedTimeouts`. The kill NEEDS the real store: only it reports a live busy window (§6.1). |
| M28 **envelope-write bound** | `host/evidence/proof_producer.go`: in `GenerateProof`, replace the context passed to the writer — `writer.WriteObject(ctx, o)` becomes `writer.WriteObject(context.Background(), o)` (the round-11 state, restored as a mutant: a write inheriting no deadline; the mutant compiles and keeps the write) | `TestProducerEnvelopeWriteIsBoundedByCallerContext` (AC20) | Correct code: with the store's sole pooled connection occupied by a decoy holder and `GenerateProof` called under a short deadline, the producer returns a bounded operational error and stores NO envelope. Under the mutant the write's context is never done, so the call blocks on the connection reservation until the decoy releases and then stores the envelope; failure: `envelope written after the caller's deadline expired; want a bounded operational error and no stored object` (a test-side 20x watchdog reds the hang mode). The kill NEEDS the real store — the reservation wait it exercises is `database/sql` behaviour against `SetMaxOpenConns(1)`, which any in-memory fake writer replaces (§6.1). |
| M27 **depth-bomb refusal** | `host/evidence/proposal_codec.go` (exact file fixed at implementation; the neutered branch is `DecodeProposal`'s own strict-decode refusal, not M9's report-codec branch): change the decode-error refusal `if err != nil` to `if false && err != nil`, returning the zero claim beside the swallowed error | `TestNestingDepthBombWithinByteCapIsRefused` (AC19) | Correct code: a depth-131,071 `[[[[…` payload inside the 256 KiB cap is refused typed — the stdlib scanner's `exceeded max depth` error surfaces as `DecodeProposal`'s decode refusal in well under a millisecond (V50). Under the mutant the depth bomb's error is swallowed and a zero `ClaimedEvidence` returns as success; failure: `depth bomb decoded as ClaimedEvidence; want typed decode refusal`. The same test's shallow depth-10 control keeps a refuse-everything mutant from passing the arm vacuously. |

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
resolver that refuses everything. M22's and M23's shared control lives inside
`TestRealStoreBlockedObjectReadReturnsWithinObjectReadTimeout` itself: with no decoy read in
flight, the same validator configuration validates one good authenticated report through the
same REAL store under the same `context.Background()` and resolves it to
`ResolvedGradeProven`, proving both that the derived deadline does not refuse the fast path
and that the fixture reaches minting, before the blocked arm asserts the timeout refusal.
M24 carries the M2–M14 control shape. M25's control is in-test and two-armed: with the hook
installed but performing no write, the same `ReadObject` returns the original payload matching
its probe, proving the fixture reaches the success path THROUGH the transaction; and the
mutating arm asserts the hook FIRED and records the writer's observed outcome — committed or
busy-refused — so a green can never come from a writer that never ran. M26's
control is inside AC18's own test: over the SAME real store, the ordered arm
(`ObjectReadTimeout` 3 s > the reported 2000 ms window) constructs, validates one good
authenticated report, and resolves `ResolvedGradeProven` — proving the ordering refusal does
not refuse the ordered path and the fixture reaches minting — before the unordered arm asserts
the refusal. M27's control is in-test: the same `DecodeProposal` accepts a depth-10 control
payload before the depth-bomb arm asserts the typed refusal.

### 6.1 Fake audit — what each prescribed fake makes true (round 9)

Round 8's reject was a FIXTURE defect: the round-7b M22's prescribed fake observed the
context, so the mutant died for a property the real `host/store` reader had never been shown
to have. The general rule, now applied to every row above: a prescribed fake is admissible
where it supplies INPUT the mutated mechanism consumes; it is inadmissible where it supplies
the PROPERTY the mutation is supposed to expose. Row by row:

- **M1** — AILANG leg; no fake.
- **M2** — the fake is never reached: the `invalid_ref` refusal precedes the read. Nothing the
  fake makes true is load-bearing.
- **M3** — the fake makes "an absent ref reports absent" true; production has that property at
  the same shape (`GetObject`'s `sql.ErrNoRows → ok=false` branch, and `ReadObject` specifies
  the identical absent branch), and §8.2 admits the real Store as an integration control.
- **M4** — no fake: the killer feeds the real store, because the mutated line IS store code.
- **M5, M7–M9, M11–M14** — the fake supplies attacker-shaped input BYTES; the mutated
  mechanism and its observable live wholly in the validator, so no fake behavior stands in for
  an unshown production property.
- **M10** — same class as M5: hand-authored envelope bytes are the INPUT; the MAC branch under
  mutation is validator code.
- **M15** — the fixture checker emits `verify.errors=1` JSON; the parse-and-refuse path under
  mutation is the producer's own code consuming that input.
- **M16–M19** — no fakes (API-inventory, gate, projection, and golden mutations).
- **M20, M21** — the refusals under test are pure Go value semantics; readers appear only as
  data suppliers in the mint-path controls.
- **M22, M23** — NO fake participates in the kill, by construction: both mutants are killed
  against the real store (AC16). The retired round-7b killer,
  `TestBlockingObjectReadReturnsWithinObjectReadTimeout` with its context-observing fake, is
  DELETED, not retained beside the honest test.
- **M24** — the input is the real store's typed refusal; the mapping under mutation is
  validator code.
- **M25** — no fake: the killer runs the real store, and the test-side scheduling hook WRAPS
  the interval between the two statements, replacing nothing — it supplies a scheduling INPUT
  (when the concurrent write is attempted), while the property under mutation, snapshot
  isolation, is supplied by the real SQLite transaction. The mutating `UPDATE` is
  test-authored by necessity: V47 verifies production has no mutation path on `objects`,
  which is exactly why the test must inject one.
- **M26** — no fake participates in the kill: the killer constructs the validator over the
  REAL store, because only the real store reports a live busy window. A fake reporting a
  window it does not have would supply the PROPERTY under test — the exact inadmissible class
  this section exists to name — so fakes simply do not implement the reporter and skip the
  check by design (§3.2).
- **M28** — no fake writer participates in the kill: the killer writes through the REAL
  store, because the wait it mutates is `database/sql`'s connection reservation against
  `SetMaxOpenConns(1)`. An in-memory fake writer never blocks, so it would supply the
  PROPERTY under test (a write that always completes promptly) — the same inadmissible class
  as M26's reporter, and the reason AC20 is a real-store criterion.
- **M27** — the depth bomb is test-authored INPUT bytes; the refusing mechanism (the stdlib
  scanner's depth limit surfacing through `DecodeProposal`'s own strict-decode refusal) is the
  code under mutation. No fake.

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
3. **Bounded authenticated validation.** The envelope is bounded at I/O: the store's
   `ReadObject` refuses an over-bound payload BEFORE materialization (length probe first, the
   payload column absent from the probe's select list — §8.2) and returns at most `maxBytes`
   bytes; there is no stream and no caller-side cap. Readings that make this criterion
   falsifiable are scoped to CODE, not to this document's prose (which already quotes these
   tokens): base reading at `35fd875`, 0 files / 0 occurrences of `ReadObject` and 0 / 0 of
   `ObjectMeta` in `*.go` under `host/`, `cmd/`, and `scripts/` (V44; the retired `OpenObject`
   also reads 0 / 0 there); head reading, ≥ 1 `func (s *Store) ReadObject(` in non-test
   `host/store` and ≥ 1 `.ReadObject(` call in non-test `host/evidence` — scopes in which base
   and head DIFFER by construction. The same-scope, same-instrument control counts a DIFFERENT
   token — `GetObject`, 23 files / 80 occurrences (V44) — so it fires while every check reads
   zero; a control that can only fire on the check's own hits is not a control (the
   iteration-95 lesson). The round-7 alternative spelling, a `maxBytes` parameter on
   `GetObject`, is DELETED (V40, §10.8), and the round-7b streaming spelling is RETIRED
   (§10.10). Decoded reports are bounded before strict
   decode, hash is recomputed, semantic/interface identities match, envelope and report canonical
   bytes round-trip, HMAC-SHA-256 is compared in constant time, subject and compiler match, the
   `verify.results[]`-derived verified set is non-empty, and success/error/counterexample fields
   agree.
4. **No fallback.** Every validation failure yields its exact `UnsupportedReason` or an explicit
   operational error. No failure result carries any grade.
5. **Producer.** The pinned executable is byte/version checked; execution is time/output bounded
   by the tranche-1-owned runner inside `host/evidence` (§3.4, V35), with deadline overrun and
   output overflow each a refusal that emits no report, and stdout/stderr captured separately;
   the bounds are configured at the API, not ambient — `NewProducer` refuses a non-positive
   `execTimeout` or `maxOutputBytes` (§3.2, round 9b);
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
9. **Required mutations.** Every tranche-1 row M1–M5, M7–M19, and M20–M28 reds with its
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
16. **Deadline-bounded validation (wait bound), proven against the REAL store.** A context
    parameter is not a bound; a deadline is — 11 ratified production store reads already pass
    `context.Background()` (V41), so the validator may inherit no deadline from anyone.
    Validator configuration carries a required positive `ObjectReadTimeout`; `NewValidator`
    refuses zero or negative. `ValidateProof` derives
    `context.WithTimeout(ctx, ObjectReadTimeout)` before reading the envelope object, even when
    the caller supplies no deadline — `context.Background()` included. Expiry or cancellation
    is the explicit operational error of §2.5, mints no seal, and carries no grade. Named
    integration test `TestRealStoreBlockedObjectReadReturnsWithinObjectReadTimeout` proves it
    against the ACTUAL `host/store` — no fake, because round 8 measured that a
    context-observing fake makes the mutant die for a property the real reader has never been
    shown to have (§10.9). Stimulus: a decoy goroutine holds the store's single pooled
    connection (`store.go:298` pins `db.SetMaxOpenConns(1)`, V45) with one in-flight
    `GetObject` of a decoy object sized at implementation so a single read outlasts 20×
    `ObjectReadTimeout` — the hold duration is ASSERTED in-test, not assumed; `ValidateProof`
    then runs under `context.Background()` against a file-backed store, its `ReadObject` waits
    for the pooled connection, and `database/sql`'s connection wait honors the supplied context
    by stdlib contract, so the attempt must return the operational error within the bound,
    carrying no seal. The classification is race-honest: a seal with wall-clock BELOW the bound
    means the attempt won the freed connection — the stimulus retries, bounded at 5, and retry
    exhaustion is a loud instrument failure, never a pass (rule 3a); a seal at ≥ 2× the bound
    is the mutant signature and reds immediately; a test-side 20× watchdog makes the hang mode
    a red, not a hang. The same test's no-decoy control seals and resolves to
    `ResolvedGradeProven` first. The blocking MECHANISM is chosen deliberately (V45): lock
    contention is REJECTED because iteration 94 measured a lock-blocked read returning at
    `busy_timeout` (2.043 s under a 300 ms deadline) — bounded by the wrong constant, the
    queue-row-22 composition this tranche must not absorb; interrupting the target read's own
    blob transfer is REJECTED as the stimulus because iteration 94's in-flight interruption
    measurement was a many-opcode query, and asserting the same of ONE blob-column
    materialization would generalize past the evidence — the exact defect round 8 rejected in
    iteration 94's doc; the pool wait depends only on the stdlib contract plus the measured
    single-connection pool, independent of SQLite's interrupt granularity and of
    `busy_timeout`. Residual, stated: under real LOCK contention, `ReadObject`'s wait is
    bounded by `busy_timeout` (2000 ms, `writer_lock.go:179`, V45), not by
    `ObjectReadTimeout`, whenever the former exceeds the latter — that composition is queue
    row 22's subject and is deliberately not absorbed here; round 11 re-verified row 22 OPEN
    by command (§10.13) and pinned the ordering that keeps the composition safe at
    construction (AC18, M26). Readings that make this criterion
    falsifiable are scoped to CODE, not to this document's prose (which quotes the token
    throughout): base reading, 0 files / 0 occurrences of `ObjectReadTimeout` in `*.go` under
    `host/`, `cmd/`, and `scripts/` at `35fd875` (V44); head reading, ≥ 1 in non-test
    `host/evidence` plus ≥ 1 `context.WithTimeout` in non-test `host/evidence` — scopes in
    which base and head DIFFER by construction (the base `context.WithTimeout` control count
    of 8 sits entirely outside `host/evidence`, which does not exist at base, so the control
    fires where the check cannot). M22 and M23 must red this arm.
17. **One-snapshot object read (probe/payload atomicity), proven against the REAL store.**
    `ReadObject`'s probe and payload statements run on one reserved database connection inside
    one read transaction/snapshot (§8.2; `gpt5-6-sol`'s round-9 fix, verbatim). Named
    integration test `TestConcurrentMutationCannotDesyncProbeAndPayload` proves the payload
    returned cannot differ from the row whose length was checked: against a file-backed real
    store, a package-private scheduling hook in `host/store` (nil in production; it WRAPS the
    interval between the two statements and replaces nothing — admissible under §6.1's rule
    because it supplies a scheduling INPUT while the property under test, snapshot isolation,
    is supplied by the real SQLite transaction) fires between the probe and the payload
    statement and attempts `UPDATE objects SET payload = <different length and content> WHERE
    hash_ref = ?` through a SECOND `sql.DB` handle on the same database file. The mutating SQL
    is TEST-AUTHORED by necessity: V47 verifies production has no mutation path on `objects`,
    which is exactly why the test must inject one. The assertion is the read-side invariant:
    the returned payload's length equals the probed `ObjectMeta` length and its recomputed
    hash equals the requested ref. The test accepts EITHER writer outcome as valid environment
    behavior, because the journal mode is deliberately not assumed (unmeasured): under a
    rollback journal the writer blocks on the read transaction's held shared lock and returns
    busy at `busy_timeout` (2000 ms, V45); under WAL it commits while the reader's snapshot
    stays fixed — the invariant holds under both, and the hook must be asserted to have FIRED
    with the writer's outcome recorded, so a green can never come from a writer that never
    ran. The kill is DETERMINISTIC — the hook forces the interleaving, so no retry
    classification and no bespoke watchdog participate — and NEEDS the real store: the mutated
    lines are `host/store` code any fake replaces. Base/head falsifiability rides on AC3's
    scopes: the enclosing `ReadObject` reads 0 files / 0 occurrences at base (V44), so every
    AC17 token is new by construction. M25 must red this arm.
18. **Timeout ordering pinned at construction (`D-WORLD-23` obligation (ii)).** The narrowed
    wait-bound claim (§3.3 step 2) depends on the composition condition
    `busy_timeout < ObjectReadTimeout` — an ordering NOTHING in the tree asserts today (V49):
    the constant is unexported (`host/store/writer_lock.go:179`); its only existing pin is a
    VALUE pin (`context_read_test.go:209` pins 2000 against no consumer's deadline);
    `ObjectReadTimeout` occurs in zero `.go` files because it is this document's proposed
    constant; and `host/daemon/handlers.go:299–302` names the gap in the codebase's own
    words — "an ORDERING nothing in this code asserts, not a guarantee". The pin: `host/store`
    gains the additive exported accessor `BusyTimeout() time.Duration`, read from the LIVE
    `PRAGMA busy_timeout` on the store's own connection (never a re-export of the constant, so
    a caller-overridden DSN reports its true window); `NewValidator` type-asserts its reader
    for that optional capability and REFUSES a positive `ObjectReadTimeout` less than or equal
    to the reported window with the dedicated `ErrUnorderedTimeouts`, naming both values. This
    binds the two values actually used at run time — the caller's configured timeout and the
    window SQLite will actually apply — never two literals: an assertion comparing two
    constants is a tautology that stays green while the configured timeout changes. Named test
    `TestConstructorPinsBusyTimeoutBelowObjectReadTimeout`, against the REAL file-backed
    store: the ordered arm (`ObjectReadTimeout` = 3 s > the reported 2000 ms) constructs,
    validates one good authenticated report, and resolves `ResolvedGradeProven`; the pin arm
    (a POSITIVE 1 s <= 2000 ms) must be refused at construction. AC16's existing non-positive
    refusal is a WEAKER check and stays green on the pin arm's stimulus — 1 s is positive — so
    the two criteria are distinguishable by construction; the same test also cross-checks the
    accessor against a direct `PRAGMA busy_timeout` query on the same store, so the accessor
    cannot drift into a stale literal. A reader reporting NO window (a bounded fake with no
    lock layer) skips the check: the refusal is scoped to compositions where a lock window
    exists to order against (§3.2, §6.1).
    **Round-11b amendment (`gpt5-6-sol`'s alternative fix, verbatim, §10.14).** The accessor is
    a NONBLOCKING CACHED property resolved once during `Open` from the DSN the store was opened
    with — not a live `PRAGMA` query, which with `SetMaxOpenConns(1)` could wait unboundedly for
    the sole pooled connection and had no error channel to report it. The withdrawn live-read
    claim costs nothing: the window is immutable after `Open` (V51). Two consequences for this
    criterion, and the second is the direct answer to the reviewer's catch. (1) The "cross-check
    against a direct `PRAGMA busy_timeout` query" moves from the CONSTRUCTION path to the TEST:
    AC18's test still issues that query itself, on a store whose connection is free, and
    requires it to equal `BusyTimeout()` — so the cached value cannot drift into a stale literal,
    while no production path performs the query. (2) A new REQUIRED arm: with the store's sole
    pooled connection OCCUPIED by a decoy holder, `BusyTimeout()` must return promptly and
    `NewValidator` must complete — the criterion asserts an upper bound on elapsed time far below
    the decoy's hold, so a regression to a pooled read reds here rather than hanging silently.
    That arm is what makes "nonblocking" a measurement instead of an adjective.
    Base/head falsifiability, scoped to CODE:
    `ErrUnorderedTimeouts` reads 0 occurrences in `*.go` at `516836f` and the accessor
    spelling `func (s *Store) BusyTimeout(` reads 0 in `host/` (same-scope control
    `busyTimeoutMillis` fires with 5 lines, V49); at head each reads ≥ 1 in non-test
    `host/evidence` and non-test `host/store` respectively. M26 must red this arm.
19. **Nesting-depth pin within the byte cap (round-10 note, discharged by measurement).**
    `DecodeProposal`'s resilience to `[[[[…` inside its 256 KiB cap is measured, in both
    directions (V50): the deepest cap-fitting bomb refuses in under a millisecond with no
    panic, and the deepest ACCEPTED payload costs single-digit milliseconds — but the refusing
    limit is the UNEXPORTED stdlib constant `maxNestingDepth = 10000`
    (`encoding/json/scanner.go:148`), not a documented guarantee. Named test
    `TestNestingDepthBombWithinByteCapIsRefused`: a depth-131,071 payload of 262,142 bytes —
    inside the byte cap — through `DecodeProposal` must produce the typed decode refusal and
    no `ClaimedEvidence`; a depth-10 shallow control in the SAME test must decode, so the arm
    cannot pass vacuously. No depth guard of our own is added; the pin exists precisely
    because an unasserted dependency on an unexported internal is the defect shape this round
    closes twice (§10.13). Falsifiability rides on AC16's scope note: `host/evidence` does not
    exist at base, so every AC19 token is new by construction. M27 must red this arm.
20. **The producer's envelope write is wait-bounded (round-11b, `gemini-3-1-pro`'s fix,
    verbatim).** Round 11 bounded the READ side (AC16), the subprocess (AC5) and, in 11b, the
    construction-time reporter — and left the producer's WRITE inheriting nothing, against the
    same one-connection pool (`store.go:298`, V51). The existing `Store.PutObject(o Object)
    error` takes no `context.Context` at all (V43, re-verified V51), so a `PutObject`-shaped
    `ObjectWriter` is unbounded by construction. The fix: `ObjectWriter`'s one method becomes
    `WriteObject(ctx context.Context, o store.Object) error` (§3.2); `host/store` gains that
    additive method, performing the connection reservation and the insert under the supplied
    context, analogous to `ReadObject` (§8.2); and `Producer.GenerateProof` passes its own
    context to it (§3.4). Named test
    `TestProducerEnvelopeWriteIsBoundedByCallerContext`, against the REAL file-backed store: a
    decoy holds the sole pooled connection, `GenerateProof` runs under a short deadline, and the
    criterion requires a bounded operational error and NO stored object — with a same-test
    control in which the decoy is installed but releases immediately, proving the fixture reaches
    the success path and stores the envelope, so a refuse-everything mutant cannot pass the arm
    vacuously. `PutObject` itself is NOT changed and its **8** existing call sites outside
    `host/store` are untouched (V51; same-pipeline control `GetObject` = 13, reproducing V43), for the same additive reason as `ReadObject`. Base/head falsifiability, scoped to
    CODE: `func (s *Store) WriteObject(` reads 0 occurrences in `host/` at `516836f` (same-scope
    control `func (s *Store) PutObject(` fires with 1, V51); at head ≥ 1 in non-test `host/store`.
    M28 must red this arm.

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
  exactly THREE additive surfaces, none changing an existing signature. FIRST, the bounded
  object read, the `D-WORLD-21` arm-A ruled spelling
  (§10.10):
  `ReadObject(ctx context.Context, ref hashref.HashRef, maxBytes int64) (ObjectMeta, []byte, error)`.
  It enforces `maxBytes` BEFORE materialization and performs the complete read under the
  supplied context, by a stated mechanism: two ordered statements on ONE reserved database
  connection inside ONE read transaction/snapshot (`gpt5-6-sol`'s round-9 fix, in the
  reviewer's own words — §10.11) — first a probe, `SELECT interface_hash_ref, semantic_id,
  provenance, length(payload) FROM objects WHERE hash_ref = ?`, whose select list OMITS the
  payload column, so no payload byte crosses the driver into Go before the guard; then a typed
  `*ObjectTooLargeError` refusal (carrying the probed size, returning no payload) when the
  probed length exceeds `maxBytes`; and only under the guard, `SELECT payload FROM objects
  WHERE hash_ref = ?` — BOTH statements issued through the transaction's `QueryRowContext`
  with the SUPPLIED context, and the connection reservation and transaction begin taking that
  same context. The probe is the transaction's first statement and fixes the snapshot, so the
  payload the second statement materializes is the row whose length the probe checked — a
  guarantee carried by the TRANSACTION. Two adjacent facts are corroboration, not the
  mechanism, and neither may be read as satisfying it: V45's `db.SetMaxOpenConns(1)` means no
  other statement FROM THIS PROCESS can interleave between the two, but a pool bound is a pool
  property, not a snapshot — it says nothing about a second process or a second `sql.DB`
  handle on the same file (exactly what AC17's test opens), and it is a store-constructor
  detail a future change could move without touching `ReadObject`. And the objects table
  itself is immutable, content-addressed, and insert-only — a claim this document may now
  state because V47 VERIFIES it by positive enumeration: NINE statements name `objects` in
  non-test Go/SQL under `host/` and `cmd/` — three `INSERT OR IGNORE` writes and six SELECTs,
  with zero UPDATE/DELETE/DROP/ALTER, zero triggers, and zero cascade paths — so today the
  transaction defends the probe/payload pair against a write surface measurement says does not
  exist. The transaction is required anyway: the immutability fact is repository state, not an
  API contract, no gate pins it, and V47's instrument finding shows that the negative reading
  which would "confirm" it is vacuous by construction. The pinned
  driver additionally carries SQLite's length-only column optimization (`OPFLAG_LENGTHARG` in
  the darwin build file, V46), so the probe need not traverse blob overflow pages — recorded
  as corroboration; the select-list guarantee is the load-bearing one. `ReadObject` is
  ADDITIVE, and that is the fact that dissolves the round-7 scope objection (V43): it changes
  NO existing signature, so the **13** non-test out-of-tranche `.GetObject(` CALL sites under
  `host/` and `cmd/` (6/2/2/1/1/1 across `host/broker/approve.go`, `host/broker/broker.go`,
  `host/transitionreg`, `host/daemon/handlers.go`, `host/registry`, `host/replay`) and the
  **4** interface-method DECLARATIONS (`approve.go:470`, `transitionreg.go:41`, `broker.go:39`,
  `daemon.go:324`) are ALL untouched — the two counts are taken separately, because a
  declaration read as a call is this mission's recurring enumeration error — and the frozen
  packages below do not change. `GetObject` is NOT modified; it already takes a
  `context.Context` (item 18 landed that, `store.go:468`), and the round-7 alternative arm, a
  `maxBytes` bound on `GetObject`, was unsatisfiable by construction precisely because it was
  NOT additive (V40) and is DELETED, not retained as an option.
  `ObjectMeta` carries the non-payload columns `GetObject` already returns (`InterfaceHash`,
  `SemanticID`, `Provenance` — V36's refresh at `store.go:468`) that §3.3 steps 4–5 check,
  plus the probed payload length; a bytes-only return would strip data the validator consumes.
  No schema change; the objects
  table is untouched. The binding invariant: no object read for an untrusted ref may
  materialize payload bytes unless the probed length is within `maxBytes`; the refusal carries
  the probed size and no payload. Sequencing with item 18
  (`w-daemon-read-cancellation`): item 18 LANDED first and bounds the SIGNATURE — every store
  read now takes a `context.Context` — but a context parameter is not a wait bound: 11
  production store reads pass `context.Background()` by item 18's own ratified deferral DR-2,
  pinned by `TestNoNewDeadlineFreeStoreReads` (V41). Those 11 sites are item 18's follow-on
  obligation (the pin's own comment names the 11 → 0 path) and MUST NOT be absorbed here; this
  tranche imposes its own deadline at its own call site via `ObjectReadTimeout` (§3.3 step 2,
  AC16). Tests may use Store as an integration control, and AC16's killer test REQUIRES the
  real store (§6.1). SECOND (round 11, `D-WORLD-23` obligation (ii); amended round 11b): the
  exported read-only accessor `BusyTimeout() time.Duration` — a NONBLOCKING CACHED property
  resolved once during `Open` from the DSN the store was opened with, then a field read. Round
  11 specified it as a LIVE `PRAGMA busy_timeout` query and `gpt5-6-sol` rejected that: with
  `SetMaxOpenConns(1)` a context-free pragma query can wait unboundedly for the sole pooled
  connection, and it has no error channel — the round's own remedy carrying the round's own
  defect. Its alternative fix is applied verbatim (§10.14). The withdrawn "live" claim costs
  nothing, because the window is IMMUTABLE after `Open`: it is applied per physical connection
  from the `_pragma` DSN parameter (`writer_lock.go:196`) and non-test `host/`+`cmd/` code
  issues ZERO runtime `PRAGMA busy_timeout` writes (V51). A caller-overridden DSN is therefore
  still reported at its true runtime value, and no production path acquires a connection to
  learn it. Consumed by `host/evidence`'s constructor ordering refusal (§3.2, AC18, M26).
  **THIRD (round 11b, `gemini-3-1-pro`'s fix, verbatim):** the additive
  `WriteObject(ctx context.Context, o Object) error`, which performs the connection reservation
  and the insert under the supplied context, analogous to `ReadObject` — the producer's envelope
  write is otherwise unbounded against the same one-connection pool, because the existing
  `PutObject(o Object) error` takes no context at all (V43, re-verified V51). `PutObject` is
  NOT changed and its **8** existing call sites outside `host/store` are untouched (V51;
  same-pipeline control `GetObject` = 13, reproducing V43) (§3.2, §3.4, AC20, M28). None of the
  three changes the schema or an existing signature, and no new statement runs against
  `objects` beyond the additive read and the additive insert; the 13 out-of-tranche
  `.GetObject(` call sites and 4 interface declarations above are untouched, for the same
  additive reason as `ReadObject`.
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
| **1. This document, `w-validated-proven-evidence-boundary`** | Library-only proof report schema/producer and authenticated envelope; caller-supplied 32-byte MAC key; Evidence codec; bounded/hash/type/success validation; bounded `host/store` object read (`D-WORLD-19` arm A; `D-WORLD-21` arm A `ReadObject` spelling with pre-materialization bound and one-snapshot probe/payload read transaction, round 9b); validator-imposed `ObjectReadTimeout` wait bound proven by the real-store blocked-read integration test; sealed mint authority; untrusted `ProofReceipt → CLAIMED`; host-only resolved `PROVEN`; refusal-branch mutations; persistent named gates; constructor-pinned `busy_timeout < ObjectReadTimeout` ordering (`D-WORLD-23` obligation (ii), AC18) and stdlib nesting-depth pin (AC19), round 11; nonblocking cached busy-window accessor and context-bounded envelope write (AC20), round 11b. Explicitly non-production. | **5.35 d** |
| **2. `w-proven-evidence-production-key-wiring`** | Verify and name the actual production composition root and configuration/state mechanisms; provision/load durable MAC key material; retain the keyed validator before serving; abort startup closed. | **1.0 d** |
| **3. `w-validated-replay-evidence-boundary`** | Typed replay report and untrusted replay receipt; integrate only successful full-episode replay; bind episode/log head and interpreter set; make missing/failed/divergent replay explicitly unsupported; M6 persistent mutation. `RecordedEffect` stays `ATTESTED`. | **3.0 d** |
| **4. `w-proven-evidence-renderer-consumption`** | Renderer/read API that accepts only sealed or freshly revalidated evidence; display `PROVEN` only from that value; explicit `UNSUPPORTED` for every validation failure; end-to-end agent-forgery and restart/revalidation tests. | **2.0 d** |
| **Total** | All reviewer obligations | **11.35 d** |

Tranche 1 arithmetic:

| Work | Time |
|---|---:|
| Strict report/envelope/Evidence codecs and object integration | 0.75 d |
| Bounded store read (`ReadObject` with `ObjectMeta`, pre-materialization length probe, typed refusal, one-snapshot read transaction), store-side tests including the concurrent-mutation test and its scheduling hook, M25 | 0.75 d |
| `ObjectReadTimeout` configuration, `WithTimeout` derivation, real-store blocked-read integration test with its decoy rig, M22/M23/M24 | 0.50 d |
| Bounded pinned proof producer, caller-key MAC integration, and fixtures | 0.60 d |
| Validator, sealed authority, receipt containment, resolved-grade API | 0.75 d |
| Kernel mapping, projection, golden, AILANG pins | 0.35 d |
| Named Go-test manifest and self-mutation gate | 0.45 d |
| Mutations, full pinned gates, review contingency | 0.60 d |
| Timeout-ordering pin: store `BusyTimeout()` cached accessor, constructor refusal, real-store AC18 test incl. its occupied-connection nonblocking arm, M26 (rounds 11 + 11b) | 0.30 d |
| Nesting-depth pin: AC19 depth-bomb test with shallow control, M27 (round 11) | 0.15 d |
| Bounded envelope write: additive store `WriteObject(ctx, …)`, `ObjectWriter`/`GenerateProof` context threading, real-store AC20 decoy test, M28 (round 11b) | 0.20 d |
| **Total** | **5.35 d** |

The tranche-1 estimate removes 0.5 day previously assigned to an unwired key lifecycle, and
round 7 adds 0.50 day for the `D-WORLD-19` arm-A bounded store read — priced as its own row
rather than absorbed, because `host/store` entering the tranche is new work, not a re-labelling.
Round 7b adds 0.25 day for the deadline machinery (`ObjectReadTimeout` plumbing through
configuration and constructor refusal), again priced rather than absorbed into the contingency
row. Round 9 is priced in BOTH directions rather than rounded away: retiring the streaming
reader REMOVES work (no reader lifetime or close path, no caller-side `io.LimitReader`, no
fake blocking-reader test), and the ruling ADDS work (the store-side pre-materialization probe
with its typed refusal, and the real-store blocked-read integration test with its
measured-duration decoy rig and race-honest retry classification, plus mutation rows
M23/M24) — net +0.25 day on the deadline row. Round 9b adds 0.25 day on the bounded-store-read
row for the one-snapshot mechanics — the reserved-connection read transaction, the scheduling
hook, `TestConcurrentMutationCannotDesyncProbeAndPayload`, and mutation row M25 — priced rather
than absorbed; the round-9b producer-API fix (§3.2) is specification of already-priced
construction and moves no row. Round 11 prices both pins rather than absorbing them — +0.25
day for the timeout-ordering pin (the store accessor, the constructor refusal, the real-store
AC18 test, M26) and +0.15 day for the nesting-depth pin (AC19 and M27); round 11b adds +0.05
day to the ordering row (the accessor becomes a cached property and AC18 gains its
occupied-connection nonblocking arm) and +0.20 day for the bounded envelope write (the additive
store `WriteObject`, the context threading, AC20's real-store decoy test and M28). At
**5.35 days** the tranche EXCEEDS the 3–4 day sprint guardrail by 1.35 days. That is stated rather than
hidden — round 9's own standard refuses to disguise the number by leaving 4.0 because it was
4.0 — and the re-scope question is answered here rather than deferred, because `D-WORLD-23`
makes shedding scope the mirror of the standing rule, so the option was examined with
arithmetic. DECOMPOSITION IS REJECTED. The only coherent piece separable without breaking the
authority boundary is the producer (§3.4): the validator never calls it, and the validator's
tests hand-author their envelopes (AC13 already does). Shedding it — the 0.60-day producer
row, the `ObjectWriter` seam and envelope-store share of the codec row (~0.10 d), and
M15/AC5's share of the mutation row (~0.10 d) — removes about 0.80 day and leaves **4.35
days, still over the guardrail**: the decomposition fails its own purpose, while costing
tranche 1 its end-to-end demonstration (no in-repo path from an `ai-check` run to an
authenticated envelope; every validator fixture hand-authored) and deleting
`NewProducer`/`Producer.GenerateProof`, the API surface a reviewer's round-9 fix forced INTO
§3.2 one round ago. Every other candidate is unsheddable by authority or by ruling: the
kernel arm is the item's design result; the validator/seal/`Resolve` chain and the
named-manifest gate ARE the authority boundary; the bounded store read and both real-store
integration tests are `D-WORLD-19`/`D-WORLD-21` arm-A obligations; and the two pins are
`D-WORLD-23` obligation (ii) and the owed round-10 note, added by the same ruling that
unparked the document. So the true number stands at 5.15, not rounded, and the 1.15-day
excess buys, by name: (a) the ratified real-store bounded read with its blocked-read and
concurrent-mutation integration tests; (b) the ordering pin that makes the narrowed claim
more than a quieter claim; (c) the depth pin — together the property §10.12 records no
reviewer had yet been shown: every load-bearing bound this tranche depends on asserted in the
tree, not argued in prose.
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

## 10.10 Round 9 revision (D-WORLD-21 arm A applied)

- **The ruling.** Mark Edmondson answered `D-WORLD-21` attended on 2026-08-19 (decision ledger,
  `design_docs/world-mission.md`): **ARM A — `ReadObject(ctx, ref, maxBytes)`.** The store
  enforces `maxBytes` BEFORE materialization and performs the complete read under the supplied
  context, so cancellation becomes ENFORCEABLE and the streaming reader is RETIRED. Chosen as
  the arm already consistent with the rulings on record — it is what objection 6b demanded,
  what `D-WORLD-19` arm A ratified, and what `gemini-3-1-pro` has PASSED — whereas arm B's
  burden was naming a concrete termination mechanism in `modernc.org/sqlite` that has never
  been shown to exist. `gpt5-6-sol`'s round-8 objection is UPHELD on its merits:
  `io.ReadCloser.Read` takes no context and is under no obligation to unblock on `ctx.Done()`,
  so a `context.WithTimeout` bounds the OPEN call only. Owed unconditionally under either arm:
  the real-store integration test. M22: REWRITTEN, not re-run.
- **(1) The seam, replaced at every normative site.** The round-7b
  `OpenObject(ctx, ref, maxBytes) (ObjectMeta, io.ReadCloser, error)` becomes
  `ReadObject(ctx context.Context, ref hashref.HashRef, maxBytes int64) (ObjectMeta, []byte, error)`
  in the header, §3.2's seam line, §3.3's cap paragraph and steps 2–3, §6 (M4, M22), AC3,
  AC16, §8.2, and §9. There is no `io.ReadCloser`, no caller-side `io.LimitReader`, no
  detection byte, and no reader lifetime to close on the envelope path — §3.4's producer keeps
  its own limit+1 capped readers, a different mechanism on a different input. The
  pre-materialization mechanism is STATED, not waved at: a `length(payload)` probe statement
  whose select list omits the payload column, a typed refusal, and only then the payload
  statement, both under the supplied context (§8.2); V46 corroborates that the pinned driver
  carries SQLite's length-only column optimization.
- **(2) Arm A is ADDITIVE, and the document now says so.** `ReadObject` ADDS a method and
  changes NO existing signature: the 13 non-test out-of-tranche `.GetObject(` call sites and
  the 4 interface-method declarations are untouched — counted separately, because a
  declaration read as a call is this mission's recurring enumeration error — so §8.2's frozen
  packages do not change (V43). This is the fact the document had never stated, and it
  dissolves the round-7 scope objection: the `maxBytes`-on-`GetObject` arm was unsatisfiable
  precisely because it was NOT additive.
- **(3) M22 rewritten, not re-run — and every other prescribed fake audited.** The round-7b
  M22 mutated the `WithTimeout` derivation and was killed by a fake wired to observe the
  context: vacuous by construction, the mutant dying for a property the real `host/store`
  reader has never been shown to have. The rewritten M22 mutates the ACTUAL cancellation
  mechanism — the supplied `ctx` at both `QueryRowContext` call sites inside `host/store`'s
  `ReadObject` — and its killer NEEDS the real store, because the mutated lines are code any
  fake replaces. The old arming edit survives as M23, killed by the same real-store test
  rather than by a fake. New §6.1 audits every remaining row: a prescribed fake is admissible
  where it supplies INPUT the mutated mechanism consumes, inadmissible where it supplies the
  PROPERTY the mutation is supposed to expose. The retired fake killer,
  `TestBlockingObjectReadReturnsWithinObjectReadTimeout`, is DELETED, not kept beside the
  honest test.
- **(4) The owed real-store test lands.**
  `TestRealStoreBlockedObjectReadReturnsWithinObjectReadTimeout` (AC16): a blocked read
  against the actual `host/store`, under `context.Background()`, relying only on
  `ObjectReadTimeout`, returning the §2.5 operational error within the bound with no seal —
  killing M22 and M23. Its blocking mechanism is chosen deliberately and the rejected
  mechanisms are named with their measurements: the decoy read holds the store's measured
  single pooled connection (V45) and the attempt blocks in `database/sql`'s context-aware
  connection wait — NOT lock contention (measured bounded by `busy_timeout`, 2.043 s under a
  300 ms deadline; queue row 22's subject, not absorbed), NOT single-blob mid-transfer
  interruption (iteration 94's interruptibility measurement was a many-opcode query; asserting
  it of one blob column would generalize past the evidence). The lock-contention residual is
  STATED in AC16 rather than absorbed.
- **(5) Rebased onto item 18; every grep re-derived at `35fd875`.** `GetObject` takes a
  context at HEAD (item 18 landed); every criterion that grepped `OpenObject` now greps
  `ReadObject`; all base readings were re-measured rather than inherited (V44), each with a
  same-scope control that fires while every check reads zero — the iteration-95
  non-discriminating-control lesson applied.
- **Deleted.** The envelope-path `io.LimitReader` streaming and its 256 KiB + 1 detection
  byte; the caller-side reader lifetime; the fake blocking-reader test and the round-7b form
  of M22; the `OpenObject` spelling at every normative site (it survives only in §§10.6–10.9's
  history and in V-rows recording the retired designs).
- **Round-8 obligations discharged.** Both of §10.9's unconditional obligations are now in the
  document: the real-store integration test (change 4) and the M22 rewrite (change 3).
- **What did not change.** `D-WORLD-17` arm A, the MAC seam (Mark's option B), `D-WORLD-19`
  arm A (`host/store` IS in the tranche), the non-production tranche-1 framing, the §9
  three-successor decomposition, the 11 DR-2 deadline-free sites (item 18's follow-on, V41),
  and the queue-row-22 lock-wait composition (named as a residual, not absorbed).
- **Pricing.** Removals (streaming reader, reader lifetime, caller-side cap, fake blocking
  test) and additions (store-side probe and typed refusal, the real-store test and its decoy
  rig, rows M23/M24) net to +0.25 day: tranche 1 moves 4.25 → 4.50 days, the decomposition
  total 10.25 → 10.50, and the tranche now exceeds the 3–4 day sprint guardrail by HALF a day,
  stated in §9 in those words rather than rounded away.

## 10.11 Round 9b revision (carve-out: both round-9 objections, reviewers' verbatim fixes)

- **Eligibility.** Both round-9 reviewers rejected; both objections carry concrete
  reviewer-authored fixes; neither disputes the design direction. This is the
  narrow-refinement carve-out, applied for the second time in this document's history (§10.8
  was the first).
- **(1) `gpt5-6-sol` — the unverified immutability premise, fixed in the reviewer's own three
  parts.** The objection is TRUE AS STATED: §8.2 asserted "immutable, content-addressed, and
  insert-only" and no Verification Log row established it. (a) §8.2 now requires BOTH
  statements to run on ONE reserved database connection inside ONE read transaction/snapshot
  — the reviewer's words — and states plainly that V45's `SetMaxOpenConns(1)` is
  corroboration, not the mechanism: a pool bound is a pool property, not a snapshot, and says
  nothing about a second process or a second `sql.DB` handle on the same file; the explicit
  read transaction is what carries the guarantee, and the pool fact may not be read as
  satisfying this fix. (b) The concurrent-mutation integration test lands as
  `TestConcurrentMutationCannotDesyncProbeAndPayload` (new AC17) with mutation row M25; the
  kill is deterministic via a package-private scheduling hook whose §6.1 admissibility is
  argued in place (it wraps, replaces nothing, and supplies only a scheduling input), and it
  is not unkillable — removing the transaction reds it by construction. (c) New V47
  enumerates every non-test statement naming `objects`, with exact commands, observed output,
  and two firing controls — and the enumeration came out LARGER than this round's own
  briefing: NINE statements, not five, because four `JOIN objects` reads in
  `host/store/journal.go` (lines 744, 792, 918, 966) are invisible to a `FROM/INTO objects`
  adjacency pattern. The conclusion HOLDS — three `INSERT OR IGNORE` writes, six SELECTs,
  zero UPDATE/DELETE/DROP/ALTER, zero triggers, zero cascades, and no dynamic SQL building a
  statement against the table — and the instrument finding is recorded twice over in V47:
  the negative grep for the UPDATE-statement form reads 0 repo-wide while five
  `ON CONFLICT ... DO UPDATE SET` upserts sit in the same file (the upsert spelling puts no
  table name between `UPDATE` and `SET`), and even a POSITIVE pattern is only true inside
  its scope (the JOIN sites). (d) The `immutable, content-addressed, and insert-only` claim
  STANDS in §8.2, now citing V47 at the claim — exactly the condition the reviewer set — and
  is demoted from mechanism to corroboration.
- **(2) `gemini-3-1-pro` — the absent producer API, added with the reviewer's signatures as
  the base.** §3.2's conceptual-surfaces block gains `ObjectWriter`, `NewProducer(...)`, and
  `Producer.GenerateProof(ctx context.Context, sourceRef HashRef, requiredIdentities
  []string) (HashRef, error)` — the latter the reviewer's spelling unchanged. Departures
  from the verbatim constructor, each stated with its reason because a silent departure from
  a verbatim fix is what this document has been refuted for twice: (i) `NewProducer` gains a
  `checker` parameter, because §3.4 mandates byte/version verification of the pinned
  executable and a producer constructed without knowing the binary cannot perform it; (ii)
  it gains a `writer` parameter (the new `ObjectWriter` seam), because §3.4 mandates that
  the producer STORES the authenticated envelope whose `HashRef` `GenerateProof` returns,
  and the package's only seam was read-only; (iii) `maxOutputBytes` is defined as the
  PER-STREAM cap applied independently to stdout and stderr, because a single shared cap
  would re-open the round-8 dual-stream objection this document already closed verbatim;
  (iv) constructor refusal of non-positive bounds is added to AC5, mirroring AC16's
  `ObjectReadTimeout` refusal, so the bounded-waits axiom is checkable at the API rather
  than asserted. `key [32]byte`, `execTimeout`, `requiredIdentities`, and the whole
  `GenerateProof` signature are verbatim; the verified set remains sourced from
  `verify.results[].function` filtered on `status == "verified"` (V27). V48 records the
  doc-scoped before/after readings with a firing same-file control.
- **(3) A consequence neither objection stated, applied because leaving it would rot a row:
  M22 is respelled.** The one-snapshot fix moves `ReadObject`'s connection wait out of
  `QueryRowContext` and into the explicit connection reservation. M22's round-9 edit —
  detach the two statement sites only — would leave that wait context-bounded, AC16's test
  would still pass under the mutant, and the row would survive VACUOUSLY. M22 now detaches
  EVERY context-taking call inside `ReadObject` (reservation, transaction begin, both
  statements): the same one mechanism, honest coverage.
- **Pricing.** +0.25 day on the bounded-store-read row (0.50 → 0.75): the
  reserved-connection transaction, the scheduling hook, the concurrent-mutation test, and
  M25. Tranche 1 moves 4.50 → 4.75 days — EXCEEDING the 3–4 day guardrail by three-quarters
  of a day, stated in §9 in those words — and the decomposition total moves 10.50 → 10.75.
  The producer-API fix is specification of already-priced construction and moves no row.
- **What did not change.** §§10.1–10.10 are byte-identical (verified by `cmp` against a
  pre-edit extraction); the `D-WORLD-17`/`D-WORLD-19`/`D-WORLD-21` rulings; the MAC seam;
  AC16's blocked-read stimulus and its rejected alternatives; the 11 DR-2 sites; queue row
  22's `busy_timeout` residual (AC17's writer-outcome tolerance cites the same 2000 ms
  constant but asserts nothing new about it); the four-tranche decomposition.

## 10.12 Round 10 quorum — gemini PASSES again, gpt5 rejects on the DECLARED residual, and the document PARKS on `D-WORLD-22`

Round 10 was the confirming re-quorum after the round-9b carve-out revision, run for the same
document-specific reason recorded in §10.8: on this document a verbatim-applied reviewer fix has
been refuted by its own author twice, so "applied verbatim" is not self-certifying here.

`absent_reviewers` empty, both slots present, `metered=$0.334800`.

**`gemini-3-1-pro` → PASS.** Second consecutive pass, and the first time its summary reads as an
endorsement of the premise work rather than a concession: *"The design thoroughly bounds I/O and
process lifetimes, and exhaustively verifies its premises against the repository."* One
NON-BLOCKING note, recorded here as an OWED obligation rather than applied, because the document
is parking and an unreviewed edit made after the round that parked it is exactly the change nobody
would then be reviewing:

> `DecodeProposal` caps raw untrusted envelope bytes at 256 KiB before parsing but states no
> maximum JSON **nesting depth**. A 256 KiB payload of purely nested structures (`[[[[…`) could
> still cause pathological CPU consumption, or rely implicitly on Go's internal 10,000-deep limit
> to avoid a stack-overflow panic during unmarshaling. Catch: *"Ensure the strict JSON decoder
> used for `DecodeProposal` is resilient to highly nested structures within the 256 KiB byte
> bound."*

§10.6 records what happens to a non-blocking note that is not forwarded: it returns one round
later as a BLOCKING objection. **Whoever resumes applies this in the SAME revision that applies
the `D-WORLD-22` ruling.**

**`gpt5-6-sol` → REJECT, and it is aimed at the residual the document DECLARED.**

> *"The design knowingly violates the bounded-waits axiom: AC16 admits that real `ReadObject` lock
> contention is governed by SQLite's 2000 ms `busy_timeout`, not the validator's
> `ObjectReadTimeout`, with a measured 300 ms deadline returning only after 2.043 s. This directly
> contradicts the normative claims that the complete read runs under the supplied context and that
> cancellation is enforceable."*

Its catch: *"The real-store test exercises only `database/sql` connection-pool waiting and
explicitly excludes the lock-contention path already shown to ignore the requested deadline.
Deferring that path to queue row 22 leaves tranche 1's authority-bearing validation operation
without a complete wait bound."*

**The objection is substantively CORRECT and the controller does not dispute it.** Arm A makes
cancellation enforceable for the pool wait and for an in-flight query; it does not make it
enforceable for a lock-contended wait, which returns via `busy_timeout` (2000 ms,
`writer_lock.go:179`, V45). Iteration 94 measured exactly that — 2.043 s under a 300 ms deadline —
and filed it as queue row 22.

**WHY THIS PARKS RATHER THAN TAKING A THIRD REVISION.** The carve-out's two conditions are (a) a
concrete reviewer-authored `proposed_fix` and (b) no dispute about DIRECTION. (a) holds. (b) does
not, and it fails in the way the carve-out exists to catch. The fix's own words are *"Make deadline
enforcement for lock contention part of this tranche … Remove the residual deferral"* — i.e. fold
**queue row 22, a separately filed and separately owned item**, into item 17's tranche, and delete
the deferral the document declares. That is a scope call, and this mission has now foreclosed the
carve-out on the scope axis five consecutive times before this one; the immediately preceding
instance, `D-WORLD-19`, was a scope question about widening into `host/store` that the human had to
answer.

**And it is narrower than a scope call alone, which is what makes it a decision rather than a
judgement:** `D-WORLD-21`'s ratified text chose arm A on the ground that under it *"cancellation
becomes ENFORCEABLE."* `gpt5-6-sol` has measured a path where it does not. So the reviewer is
disputing a predicate of the ruling on record, and the controller may not resolve that by editing
either the ruling or the document's claim to fit.

**Note the shape, because it is worth more than this document.** The reject landed on a DISCLOSURE.
Round 8's applied lesson was to state residuals rather than absorb them, round 9 stated this one in
AC16 in the reviewer's plain sight, and that honesty is precisely the surface the objection
attached to. Absorbing the residual instead would have drawn no objection and would have been
strictly worse — it is the failure iteration 94 filed row 22 to prevent. **A document is not
penalised for the defects it hides.** The remedy is not to say less; it is for the ask to name the
composition explicitly, which `D-WORLD-22` does.

**THE ASK IS ONE WORD** (`D-WORLD-22`). Both arms are faithful to something already ratified, which
is why the controller cannot pick:

- **Arm A — widen tranche 1 to bound the lock wait**, per `gpt5-6-sol`'s fix verbatim: every wait
  in `ReadObject`, including SQLite busy/lock retry, capped by the earlier of the caller's context
  deadline and the configured store bound (a context-aware busy handler, or a connection-local
  `busy_timeout` derived from the remaining deadline); a real-store acceptance test holding a
  conflicting lock and requiring an operational timeout with no seal within the bound; a mutation
  restoring the fixed 2000 ms wait, which must red; the residual deferral deleted. This makes
  `D-WORLD-21`'s "enforceable" true without qualification — and it absorbs queue row 22 into item
  17, against standing rule 1, re-pricing a tranche already half a day over the guardrail.
- **Arm B — keep the tranche as ruled and weaken the CLAIM to exactly what is proven.** The
  document stops saying cancellation is enforceable full stop, and says instead that every wait
  this tranche's own code performs is bounded by `ObjectReadTimeout`, that a lock-contended wait is
  bounded by `busy_timeout`, that the composition is safe only while `busy_timeout` (2 s) is less
  than `ObjectReadTimeout`, and — the part nothing in the tree asserts today — adds an assertion
  pinning that ordering, so the composition cannot silently invert. Row 22 keeps ownership. This
  costs the strength of the claim and buys the ordering guarantee the codebase currently lacks.

Under BOTH arms, two obligations are owed and neither is contingent on the answer: gemini's
JSON-nesting-depth note above, and — unchanged from §10.9 and still true — that no reviewer has yet
seen a tranche-1 document whose every load-bearing wait bound is asserted in the tree rather than
argued in prose.

## 10.13 Round 11 revision (`D-WORLD-22` arm B applied)

**The ruling and its provenance.** `D-WORLD-22` resolved as **arm B** on 2026-08-20 as the first
consequence of the `D-WORLD-23` arm-A STANDING RULE (Mark, attended, 2026-08-20T08:01:31Z, `#68`
comment `A`): when a quorum objection's `proposed_fix` would fold separately-owned work into the
tranche in front of it, the tranche KEEPS SCOPE, WEAKENS ITS CLAIM TO EXACTLY WHAT IT PROVES,
and RECORDS THE RESIDUAL WITH A NAMED OPEN OWNER — applied by the controller without a new ask.
`gpt5-6-sol`'s round-10 objection stays SUSTAINED as substantively correct; it is answered by
narrowing, never by disputing it. Nothing here re-litigates `D-WORLD-19`/`21`/`22`/`23`, and no
decision is open on this item. The rule licenses this SCOPE call only: it settled no premise and
no design direction, and none was in dispute this round.

**Obligation (i) — the residual's owner is OPEN, asserted by command at the moment of writing.**
Queue row 22 `w-daemon-lock-wait-not-deadline-bound` heads at `design_docs/world-mission.md:3583`.
Stripping struck-through spans and counting completion tokens over the row —
`sed -n '3583,3608p' design_docs/world-mission.md | perl -0pe 's/~~.*?~~//gs' | grep -cE
"LANDED|ITEM COMPLETE"` — returns **0**, while the SAME instrument over row 21, known LANDED,
returns **1**: the zero is a measurement, not a dead pattern. `git log origin/dev --oneline
--grep='lock-wait'` returns exactly ONE commit, `912009d`, and it is the `D-WORLD-23` ledger
commit whose own body says "row 22 keeps ownership of the lock-wait bound" — the branch's one
mention of the residual ASSIGNS it to row 22 rather than landing it. Row 22 may own the
residual. (The controller measured the same facts this iteration; re-run first-party here.)

**D1 — what narrowed, before and after.** Before (the header bullet, formerly lines 24–25): "the
store enforces `maxBytes` BEFORE materialization and performs the whole read under the supplied
context, **so cancellation is enforceable**". After, at every normative site (the header bullet,
§3.3 step 2, §8.2, with AC16 already carrying the disclosure): (1) every wait this tranche's OWN
code performs is bounded by `ObjectReadTimeout`; (2) a LOCK-contended wait is bounded by
`busy_timeout` (2000 ms, `host/store/writer_lock.go:179`), which this tranche does not change;
(3) the composition is safe ONLY WHILE `busy_timeout < ObjectReadTimeout` — residual owned by
OPEN queue row 22. AC16's residual disclosure is UNTOUCHED except to gain the row-22-open and
AC18 cross-references: round 10 proved the disclosure was right and the claim above it wrong, so
this round edits the claim and keeps the honesty. The full grep audit (`enforceable`, `under the
supplied context`, `bounded-wait`) found every remaining hit inside §§10.7–10.12 — append-only
records of rounds that happened, including §10.9's account of the round-8 arm-1 fix ("This makes
cancellation enforceable", formerly line 1358) — and those are left BYTE-IDENTICAL under the
history rule; this section records their supersession instead of rewriting them.

**D2 — obligation (ii): the pin, and where it lives.** A claim narrowed without a pin is a claim
merely made quieter. The composition condition is real and NOTHING asserted it (V49):
`busyTimeoutMillis` is unexported (`writer_lock.go:179`); its only existing pin is a VALUE pin
(`context_read_test.go:209` pins 2000 against no consumer's deadline); `ObjectReadTimeout`
occurs in ZERO `.go` files — it is this document's proposed constant; and
`host/daemon/handlers.go:299–302` already names the gap verbatim: "an ORDERING nothing in this
code asserts, not a guarantee". The cross-package seam decision, written down: the ordering
refusal lives in `host/evidence`'s `NewValidator` — the consumer that owns a deadline owns the
refusal — and `host/store` contributes only the additive exported `BusyTimeout()` accessor
reading the LIVE `PRAGMA busy_timeout` (§3.2, §8.2, priced in §9), so the pin binds the caller's
configured timeout against the window SQLite will actually apply, never two literals. AC18 kills
through the real store; M26 neuters the refusal and is distinguishable from AC16's weaker
non-positive refusal because its stimulus is a POSITIVE 1 s timeout that is merely <= the
2000 ms window.

**D3 — the owed nesting note, measured, and it is the same defect shape.** Four arms, both
toolchains (V50): the panic/CPU half is REFUTED — the deepest cap-fitting bomb (131,071 deep,
262,142 bytes) refuses in ~0.5 ms and the deepest ACCEPTED depth (9,999) costs ~7 ms, bounded by
construction — and the unasserted-internal half is SUSTAINED: the refusing limit is the
unexported `maxNestingDepth = 10000` (`encoding/json/scanner.go:148`), a stdlib internal, not a
guarantee. That is the round's spine, not a coincidence: TWICE this round the design depended on
a bound nothing asserts — SQLite's busy window ordered under the validator's deadline, and the
stdlib's depth limit under the byte cap — and both are discharged by the same remedy: state the
condition, then pin it (AC18/M26; AC19/M27).

**D4 — the §9 delta and its direction.** UP, both pins priced: tranche 1 moves 4.75 → **5.15
days** (+0.25 ordering pin, +0.15 depth pin); the decomposition total moves 10.75 → 11.15.
Decomposition was examined and REJECTED with arithmetic: the only boundary-safe separable piece,
the producer, sheds ~0.80 day and leaves 4.35 days — still over the guardrail — so shedding it
buys no compliance while costing the end-to-end demonstration; everything else is authority
boundary or ratified obligation. §9 states the true number and names what the excess buys.

## 10.14 Round 11b revision (carve-out: both round-11 objections, reviewers' verbatim fixes)

Round 11's quorum BLOCKED, both slots present (`absent_reviewers` empty — no degraded quorum, so
both rejects are considered positions), `metered=$0.3991`, artifact
`w-validated-proven-evidence-boundary-2026-08-20T15-03-08Z.json`. **This is the first round in this
document's history where the narrow-refinement carve-out is AVAILABLE**, and the reason is worth
recording: the carve-out has been foreclosed on the SCOPE axis seven consecutive times, because
every previous blocking fix wanted to fold separately-owned work in. Neither round-11 objection
does. Both are completeness defects on surfaces this tranche already owns, both carry a concrete
reviewer-authored `proposed_fix`, and neither disputes the design DIRECTION — the carve-out's two
conditions, met. Applying them is therefore not a judgement call, and parking would have
manufactured a decision the human does not have.

**Both objections landed on THIS ROUND'S OWN REMEDY, which is the finding.** Round 11 existed to
close "a bound nothing asserts". `gpt5-6-sol` observed that the fix introduced a NEW unbounded
wait: `BusyTimeout()` as a LIVE `PRAGMA` read takes no context, has no error channel, and with
`db.SetMaxOpenConns(1)` (`store.go:298`, V51) can wait indefinitely for the sole pooled
connection — the exact wait mechanism AC16 is built around. `gemini-3-1-pro` observed that round
11 bounded the read side, the subprocess and (in 11b) the reporter, and left the producer's WRITE
side inheriting nothing, because `ObjectWriter` was specified as "one `PutObject`-shaped method"
and `Store.PutObject(o Object) error` takes no context at all (V43, re-verified V51). Item 18
threaded the read side and not the write side; the design copied the asymmetry without noticing
it. Both are the mission's *guard the helper, miss the call site* shape, aimed at the guard this
document had just written.

**`gpt5-6-sol`'s fix offered TWO arms and the alternative was chosen BY MEASUREMENT, not by
preference** — this document has been burned once by picking an arm on judgement (§10.2: arm 2
taken verbatim, then correctly rejected by its own author one round later), so the choice is made
with §10.4's discriminator, *is what the fix REMOVES load-bearing for the doc's claim?* Arm 1
makes construction context-bearing and BOUNDS the query. The alternative makes `BusyTimeout()` a
nonblocking cached property and REMOVES it, at the cost of the "live" claim. V51(c) settles which:
non-test `host/`+`cmd/` code issues **zero** runtime `PRAGMA busy_timeout` operations, the window
being applied per physical connection from the `_pragma` DSN parameter at `writer_lock.go:196`
(control fires), so the window is IMMUTABLE for the store's lifetime and the removed property is
not load-bearing — a cached read and a live read return the same value by construction, including
for a caller-overridden DSN. The alternative is therefore strictly stronger against the
bounded-waits axiom, since eliminating a wait dominates capping one, and it leaves
`NewValidator`'s signature alone. Applied verbatim, including its tail: the "LIVE PRAGMA" claim is
WITHDRAWN as unsupported (§3.2, §8.2, AC18).

**The reviewer's CATCH is answered directly rather than argued away.** It asked how the accessor
obtains a pooled connection, how that acquisition terminates when the sole connection is occupied,
and how query failure is represented — and noted that AC18 tested only ordered and unordered
values. Under the applied fix the first three questions dissolve (no acquisition, no query, no
failure channel), and the fourth is closed by a NEW REQUIRED ARM: with the sole pooled connection
held by a decoy, `BusyTimeout()` must return and `NewValidator` must complete far inside the
decoy's hold. That arm is what makes "nonblocking" a measurement rather than an adjective, and a
regression to a pooled read reds it instead of hanging silently. The live-vs-cached cross-check
moves from the construction path into AC18's test, so the cached value still cannot drift into a
stale literal.

**`gemini-3-1-pro`'s fix is applied verbatim in all three of its parts** — §3.2 redefines
`ObjectWriter`'s one method as `WriteObject(ctx context.Context, o store.Object) error`; §8.2 adds
the additive `WriteObject(ctx, o) error` to `host/store`, performing the reservation and the
insert under the supplied context, analogous to `ReadObject`; §3.4 requires `GenerateProof` to
pass its own context to it. `PutObject` is untouched, so its **8** call sites outside `host/store`
(V51) are untouched, for the same additive reason as `ReadObject`. AC20 and M28 supply the
acceptance criterion and the kill — the fix named neither, and this document's §5 requires every
new guard to have both, so adding them is the standing discipline rather than a controller-invented
resolution of the objection.

**A controller-authored count was WRONG and is corrected here rather than quietly.** The first
draft of these edits said `PutObject` has "10 existing call sites" — transcribed from the queue
row's `2ef2271`-epoch measurement, not measured at HEAD. Measured in V43's scope it is **8**, and
the same-pipeline `GetObject` control returns **13**, reproducing V43 exactly: the pipeline was
sound and the number was inherited. That is rule 3b(v)(b) committed by the controller inside the
very revision that applies a rule about not laundering claims.

**§9 delta:** +0.05 d on the ordering row (cached accessor plus AC18's nonblocking arm) and
+0.20 d for the bounded envelope write, so tranche 1 moves **5.15 → 5.35 d** and the decomposition
total 11.15 → 11.35 d. Priced, not absorbed. The decomposition analysis of §10.13 is unchanged:
the only boundary-safe shed remains the producer, and shedding it still leaves the tranche over
the guardrail.

**Nothing is open for the human on this item.** No premise was disputed this round, no direction
was disputed, and no scope call was made — both fixes stay inside surfaces the tranche already
owns. Per §10.8's document-specific practice — on this document a verbatim-applied reviewer fix
has been refuted by its own author twice — a CONFIRMING re-quorum follows rather than a direct
route to planning.

## 10.15 Round 12 confirming quorum — BLOCKED, and the PATTERN is the finding

`absent_reviewers` empty, both slots present, `metered=$0.4352`, artifact
`w-validated-proven-evidence-boundary-2026-08-20T15-09-51Z.json`. Iteration metered total for
both rounds: `$0.8343` of the $5 ceiling.

**`gpt5-6-sol` → REJECT.** *"`Producer.GenerateProof` must read `sourceRef` through
`ObjectReader`, but `NewProducer(key, checker, writer, execTimeout, maxOutputBytes)` receives no
reader. Even if a reader were added implicitly, both the source read and `WriteObject` use the
caller context without requiring or deriving a deadline, so `context.Background()` can block
indefinitely on the store's single-connection pool. AC20 only proves propagation of an
already-deadlined context, not a bounded operation."*

**`gemini-3-1-pro` → REJECT.** *"Round 11b's additive `WriteObject(ctx, o)` introduces an unpinned
lock-contended wait. Because `NewProducer` configures no store timeout, it cannot pin the ordering
against the store's busy window at construction. If the caller's context to `GenerateProof` has a
deadline shorter than `busy_timeout` (2000 ms), SQLite lock contention will ignore it and overrun
the deadline. AC20 vacuously masks this by testing only connection-pool waiting with a decoy,
reproducing the exact fixture defect rejected in Round 10 for the read side."*

**Both are substantively CORRECT and the controller disputes neither.** `gemini`'s is the sharper
one and it is self-aware about the shape: AC20's decoy exercises the connection-pool wait, which
IS context-aware, and says nothing about the lock wait, which is not — the same fixture-vacuity
class the read side was rejected for in round 10, reproduced on the write side by the fix that
bounded the write side.

**THE PATTERN, AND IT IS WORTH MORE THAN ANY ONE OBJECTION.** Three consecutive rounds have
blocked on a defect *introduced by the previous round's fix*, each time one surface over:

| Round | The fix | What the next round found in it |
|---|---|---|
| 11 | pin `busy_timeout < ObjectReadTimeout` via a live-`PRAGMA` `BusyTimeout()` | the pragma query itself is an unbounded wait on the 1-connection pool (`gpt5-6-sol`) — and the producer's WRITE was never bounded at all (`gemini-3-1-pro`) |
| 11b | cache the accessor (wait removed); add `WriteObject(ctx, o)` and thread the producer's context | the producer has no reader and derives no deadline of its own (`gpt5-6-sol`); the new write inherits the SAME unpinned lock-contended wait the read has (`gemini-3-1-pro`) |
| 12 | — | (this round) |

This is not reviewer churn. Every objection has been in-tranche, concrete, and right. The mechanism
is structural: **item 18 threaded the store's READ signatures and deliberately left the waits
unbounded (DR-2, 11 sites, V41), `PutObject` was never threaded at all (V51), and the busy window
is ordered against nothing (V49)** — so every store surface tranche 1 newly touches arrives with
the same three holes, and bounding one surface only reveals the next. The document is chasing a
boundary that moves as fast as it is fixed.

**WHY THIS PARKS.** The revision budget is spent honestly: one designer revision (round 11), one
quorum, one bounded carve-out revision applying both reviewers' verbatim fixes (round 11b), one
confirming re-quorum (round 12). A third revision would be unbounded re-litigation. And the
`D-WORLD-23` standing rule does not reach this: it licenses keeping scope and narrowing a claim
when an objection would fold **separately-owned** work in, and round 12's objections are almost
entirely about work this tranche already owns — obligation (iii) is explicit that the rule is a
SCOPE licence only.

**THE ASK IS ONE WORD** (`D-WORLD-24`), and it is the mirror of `D-WORLD-22`/`23` rather than a
duplicate of them: those asked whether the tranche ABSORBS separately-owned work (answer: no, keep
scope and narrow the claim). This asks whether it SHEDS work it does own — a scope CUT of a named
core deliverable, which no ratified rule licenses the controller to make.

The observation that makes it live: **both round-12 objections are producer-side.** §10.13 examined
decomposition and rejected shedding the producer on the arithmetic then available — ~0.80 d off
5.15 d leaves 4.35 d, still over the guardrail, so it bought no compliance. That arithmetic did not
and could not include round 12's objections, because they did not exist. They do now, and shedding
the producer dissolves BOTH at once rather than costing another round to fix.

- **Arm A — shed the producer (§3.4) into its own ordered document.** Tranche 1 becomes the
  validator, the seal, the boundary and their gates; both round-12 objections dissolve unfixed
  because both are producer-side; the tranche drops the ~0.80 d §10.13 priced plus round 12's
  unpriced fixes, landing materially closer to the 3–4 d guardrail. **Costs** the end-to-end
  demonstration §10.13 named — no in-repo path from an `ai-check` run to an authenticated
  envelope, every validator fixture hand-authored (AC13 already is) — and deletes
  `NewProducer`/`Producer.GenerateProof`, the API surface `gemini-3-1-pro`'s own round-9 fix
  forced INTO §3.2 three rounds ago.
- **Arm B — keep the producer and apply round 12's fixes.** `NewProducer` gains the reader
  parameter and its own derived deadline (an `ObjectReadTimeout` twin, so the producer stops
  trusting the caller's context); the write's lock-contended residual weakens to the SAME
  composition claim the read side already carries under `D-WORLD-22` arm B, owned by OPEN queue
  row 22, with AC20's decoy arm re-scoped to say what it actually proves. **Costs** another
  revision + re-quorum round and further days on a tranche already 1.35 d over the guardrail —
  and, on this document's record, carries a real chance the fix reveals the next surface.

**WHOEVER RESUMES, under BOTH arms:** AC20's decoy arm must state that it exercises the
connection-pool wait and NOT the lock wait — `gemini`'s vacuity finding is true of the criterion as
written and is owed regardless of which arm is chosen. And `gpt5-6-sol`'s reader-parameter gap
(`GenerateProof` reads `sourceRef` through a reader `NewProducer` never receives) is a plain API
incompleteness that holds under arm B and disappears with the producer under arm A.

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

V43–V46 were run for round 9 on 2026-08-19 from repository root at `35fd875`. V39 and V42
carry in-row round-9 notes: their measurements stand as taken; the AC scopes they describe
were re-based by this round (V44), and the streaming spelling they reference is retired
(§10.10). V46 additionally records an instrument-scope catch made during this round's own
measurement: the per-target darwin driver file is a constants shim on which BOTH the check and
the control read zero, and only the control's zero exposed the wrong scope.

V47–V48 were run for round 9b on 2026-08-19 from repository root at `35fd875` (working tree:
this document modified, uncommitted; no code file differs from HEAD, so code measurements are
HEAD measurements). V47's statement enumeration is deliberately token-scoped — every
occurrence of the word `objects`, classified line by line — because an adjacency pattern
(`FROM objects|INTO objects|UPDATE objects|DELETE FROM objects`) reads five statements and is
blind to the four `JOIN objects` reads in the same package; both readings appear in the row.

V49–V50 were run for round 11 on 2026-08-20 from repository root at `516836f` (working tree:
this document modified, uncommitted; no code file differs from HEAD, so code measurements are
HEAD measurements). The controller measured the same facts earlier in the iteration and
labelled them for re-verification; both rows were RE-RUN first-party by this revision before
being cited. V50's runtime probe deliberately ran on BOTH `GOTOOLCHAIN=go1.25.6` (the pin) and
the ambient go1.26.4, because the behaviour it pins is stdlib, not repository code.

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
| V39 | **CONTROLLER-MEASURED at `52bc9ec`, round 7, after the designer run — PRIOR ART FOR THE BOUNDED-READ IDIOM EXISTS AND THE DOC DID NOT KNOW IT.** Non-test `host/` already contains **two** `io.LimitReader` sites, and one of them (`host/capsule/capsule.go:238`) already uses the **`limit+1` detection-byte** shape §3.3 step 3 now specifies. Neither is an exported bounded-read helper — both are inline stdlib idioms at their own call sites — so there is nothing to REUSE and §8.2's new `host/store` surface is not duplicating an existing API. Recorded because §3.4 applies exactly this discipline to `runBounded` (V35: reuse REJECTED with a reason rather than deferred), and applying it to `runBounded` but not to the idiom this round introduces would be the same inconsistency round 6 blocked on. AC3's scoping is unaffected: it counts `io.LimitReader` in non-test `host/evidence`, a package that does not exist at base, so these two sites cannot satisfy it. **Round-9 note: AC3 no longer counts `io.LimitReader` at all — the envelope path has no stream under `D-WORLD-21` arm A (§10.10); §3.4's producer keeps the limit+1 capped-reader idiom this row surveys, so the row stands as the reuse survey it was.** | `grep -rn 'io.LimitReader' host/ --include='*.go' \| grep -v _test` ; same-path control `grep -rn 'io.ReadAll' host/ --include='*.go' \| grep -v _test \| wc -l` | Two hits: `host/broker/registry_reconcile.go:481` `io.ReadAll(io.LimitReader(resp.Body, cfg.MaxBodyBytes))` and `host/capsule/capsule.go:238` `io.ReadAll(io.LimitReader(pipe, limit+1))`; `host/evidence` does not exist (`test -d host/evidence` rc=1). |
| V40 | At `7806cac`, `GetObject` cannot gain a `maxBytes` parameter without changing out-of-tranche callers: exactly **13** non-test `.GetObject(` call sites exist in `*.go` under `host/` and `cmd/` outside `host/store` — `host/broker/approve.go` 6, `host/broker/broker.go` 2, `host/transitionreg/transitionreg.go` 2, `host/registry/registry.go` 1, `host/daemon/handlers.go` 1, `host/replay/replay.go` 1 — and the last two are in packages §8.2 declares frozen, so the round-7 `maxBytes`-on-`GetObject` arm was unsatisfiable by construction and is deleted (§10.8). | `git rev-parse --short HEAD; grep -rn '\.GetObject(' --include='*.go' host/ cmd/ \| grep -v '^host/store/' \| grep -v _test \| wc -l`; same-call, same-scope control: `grep -rn '\.PutObject(' --include='*.go' host/ cmd/ \| grep -v '^host/store/' \| grep -v _test \| wc -l`; per-file enumeration via `grep -rln` then `grep -c` per file | `7806cac`; count `13`; control `8` — non-zero, and equal to V28's independently measured eight non-test `PutObject` sites, so the instrument reads a known quantity correctly. Per-file enumeration exactly as stated in the claim. |
| V41 | At `7806cac`, a context parameter is not a wait bound — item 18's OWN ratchet says so: `host/store/context_read_test.go` declares **11** production store reads that pass `context.Background()`, left deadline-free "by ratified deferral DR-2", with the comment stating the set "may SHRINK … but it may never GROW" and naming the "11 -> 0" reduction as the FOLLOW-ON item's mechanically observable path. So item 18 bounded the SIGNATURE, and this tranche's validator must arm its own deadline (`ObjectReadTimeout`, AC16) rather than inherit one; the 11 sites belong to item 18's follow-on, not to this tranche (§8.2). | `grep -n 'deadlineFreeReadPins\|TestNoNewDeadlineFreeStoreReads\|DR-2' host/store/context_read_test.go; sed -n '361,379p' host/store/context_read_test.go` | Pin map at `context_read_test.go:370`: `host/broker/approve.go` 8, `host/registry/registry.go` 2, `host/replay/replay.go` 1 — sum 11; matcher regexp `\.(GetObject\|GetWorld\|GetLogEntry\|GetRegistryHead\|SelectedHead)\(\s*context\.Background\(\)`; `TestNoNewDeadlineFreeStoreReads` at `:387` scans non-test `host/` and `cmd/` and reds if any file's count exceeds its pin. The comment reads verbatim "the EXACT set of production store reads this item leaves deadline-free, by ratified deferral DR-2". |
| V42 | AC16's and AC3's base readings at `7806cac`: the token `ObjectReadTimeout` occurs in **zero** `*.go` files under `host/`, `cmd/`, and `scripts/`, and `OpenObject` in **zero** `*.go` files under `host/` and `cmd/` (tests included), so both head readings differ from base by construction; the same-call positive control fires in the same scope. **Round-9 note: `OpenObject` is now the RETIRED spelling; the live base readings for AC3/AC16 (`ReadObject`, `ObjectMeta`, `ObjectReadTimeout`) were re-derived at `35fd875` in V44. This row stands as the `7806cac` measurement it was.** | `printf 'objectreadtimeout_code='; grep -rn 'ObjectReadTimeout' --include='*.go' host/ cmd/ scripts/ \| wc -l; printf 'withtimeout_control='; grep -rn 'context.WithTimeout' --include='*.go' host/ cmd/ \| grep -v _test \| wc -l; printf 'openobject_code='; grep -rn 'OpenObject' --include='*.go' host/ cmd/ \| wc -l` | `objectreadtimeout_code=0`; `openobject_code=0`; same-scope control `withtimeout_control=8` non-test `context.WithTimeout` sites (all outside the not-yet-existing `host/evidence`), so each zero is a measurement, not a broken pattern. |
| V43 | At `35fd875`, `ReadObject` is ADDITIVE — the fact that dissolves the round-7 scope objection: `host/store` exposes exactly TWO exported Object methods (`PutObject` at `store.go:444`; `GetObject` at `store.go:468`, already context-threaded by item 18), and outside `host/store` there are exactly **13** non-test `.GetObject(` CALL sites (`host/broker/approve.go` 6, `host/broker/broker.go` 2, `host/transitionreg/transitionreg.go` 2, `host/daemon/handlers.go` 1, `host/registry/registry.go` 1, `host/replay/replay.go` 1 — reproducing V40 exactly) plus **4** interface-method DECLARATIONS (`approve.go:470`, `transitionreg.go:41`, `broker.go:39`, `daemon/daemon.go:324`), counted SEPARATELY because a declaration is not a call. An ADDED method changes none of the 2, none of the 13, and none of the 4, so §8.2's frozen packages are untouched by construction. | `grep -n 'func (s \*Store) PutObject\|func (s \*Store) GetObject' host/store/store.go`; calls: `grep -rn '\.GetObject(' --include='*.go' host/ cmd/ \| grep -v '^host/store/' \| grep -v _test \| awk -F: '{print $1}' \| sort \| uniq -c` and the same pipeline through `wc -l`; declarations: `grep -rn 'GetObject(' --include='*.go' host/ cmd/ \| grep -v '^host/store/' \| grep -v _test \| grep -v '\.GetObject('`; same-scope control: `grep -rn '\.PutObject(' --include='*.go' host/ cmd/ \| grep -v '^host/store/' \| grep -v _test \| wc -l` | Methods at `store.go:444`/`store.go:468`; call total **13** with the per-file split exactly as claimed; declaration enumeration exactly the four cited lines; control `.PutObject(` = **8**, independently matching V28's and V40's eight non-test put sites, so the instrument reads a known quantity correctly. |
| V44 | Round-9 base readings at `35fd875`, scope `*.go` under `host/`, `cmd/`, and `scripts/`, tests INCLUDED: `ReadObject` **0** files / **0** occurrences; `ObjectMeta` **0** / **0**; `OpenObject` **0** / **0** (the retired spelling never landed in code); `ObjectReadTimeout` **0** / **0**; `host/evidence` ABSENT; and non-test `context.WithTimeout` under `host/`+`cmd/` = **8**, all outside `host/evidence`. The control token DIFFERS from every check token in the same scope and instrument — `GetObject`, **23** files / **80** occurrences — so the control fires while all four checks read zero: a control that can only fire on the check's own hits is not a control (the iteration-95 lesson, applied at design time rather than found at review). | `for tok in ReadObject ObjectMeta OpenObject ObjectReadTimeout; do grep -rl "$tok" --include='*.go' host/ cmd/ scripts/ \| wc -l; grep -rn "$tok" --include='*.go' host/ cmd/ scripts/ \| wc -l; done`; control: `grep -rl 'GetObject' --include='*.go' host/ cmd/ scripts/ \| wc -l; grep -rn 'GetObject' --include='*.go' host/ cmd/ scripts/ \| wc -l`; `grep -rn 'context.WithTimeout' --include='*.go' host/ cmd/ \| grep -v _test \| wc -l`; `ls -d host/evidence \|\| echo ABSENT` | All four check tokens `files=0 occurrences=0`; control `GetObject` `files=23 occurrences=80`; `withtimeout=8` (enumerated: broker/handlers.go:90, broker/registry_reconcile.go:464, capsule.go:152, daemon/handlers.go:270, daemon.go:109 comment + :678, replay.go:323, cmd/ailang-worldd/cli.go:39); `host/evidence` prints ABSENT. |
| V45 | The blocking-mechanism constraints AC16's real-store test relies on, measured at `35fd875`: the production store pool is exactly ONE connection — `store.go:298` `db.SetMaxOpenConns(1)` — so a single in-flight decoy read serializes any concurrent `ReadObject` into `database/sql`'s context-aware connection wait; and the lock-retry window is `busyTimeoutMillis = 2000` (`writer_lock.go:179`), the constant iteration 94 measured OVERRUNNING a 300 ms context deadline 6.8× (a lock-blocked read returned at 2.043 s), which is why lock contention is the REJECTED stimulus and remains queue row 22's subject, not this tranche's. | `grep -n 'db.SetMaxOpenConns(1)' host/store/store.go; grep -n 'const busyTimeoutMillis' host/store/writer_lock.go`; same-file control: `grep -c 'func (s \*Store)' host/store/store.go` (non-zero, so the instrument reads the file) | `298:	db.SetMaxOpenConns(1)`; `179:const busyTimeoutMillis = 2000`; control non-zero. The 2.043 s / 302.13 ms iteration-94 measurements are cited from the charter's queue-row-22 record, not re-run here (re-running them requires the lock rig that item measured with). |
| V46 | The pinned driver carries SQLite's length-only column optimization on THIS platform, corroborating §8.2's probe mechanism: `modernc.org/sqlite v1.54.0`'s shared darwin build file `lib/sqlite_g_0000000000000003.go` is build-tagged `(darwin && amd64) \|\| (darwin && arm64)` and contains the `OP_Column` implementation with **3** `OPFLAG_LENGTHARG` references, including the source comment "The content of large blobs is not loaded, thus saving CPU cycles." CORROBORATION ONLY: the load-bearing guarantee is structural — the probe's select list omits the payload column, so no payload byte crosses the driver regardless. Instrument-scope catch, kept: the per-target `sqlite_darwin_arm64.go` reads **0** for BOTH the check (`OPFLAG_LENGTHARG`) and the control (`OP_Column`) — it is a 5.7 KB constants shim, and only the control's zero exposed that the first scope was wrong (a count is only true inside the scope it was taken in). | `grep -n 'go:build' ~/go/pkg/mod/modernc.org/sqlite@v1.54.0/lib/sqlite_g_0000000000000003.go; grep -c 'OPFLAG_LENGTHARG' <same file>; sed -n '82650,82665p' <same file>`; wrong-scope record + control: `grep -c 'OPFLAG_LENGTHARG' lib/sqlite_darwin_arm64.go; grep -c 'OP_Column' lib/sqlite_darwin_arm64.go; wc -c lib/sqlite_darwin_arm64.go` | Build tag as claimed; count **3** (lines 63847, 82658, 87366); the quoted comment verbatim at the `OP_Column` opcode; shim reads `0` and `0` at 5,710 bytes. `go.mod` pins `modernc.org/sqlite v1.54.0` (line 5). |
| V47 | The `objects` table is INSERT-ONLY, verified by POSITIVE enumeration at `35fd875`. Scope: every occurrence of the token `objects` (case-insensitive, whole-word) in `*.go` and `*.sql` under `host/` and `cmd/`, tests EXCLUDED — 41 lines, classified one by one: **NINE SQL statements**, one `CREATE TABLE` (DDL, `schema.sql:13`), and 31 Go identifiers/comments/JSON tags/route strings. The nine: WRITES (3, all `INSERT OR IGNORE INTO objects` — a form that never modifies an existing row): `host/store/journal.go:401`, `host/store/store.go:455`, `host/store/store.go:942`. READS (6): `store.go:477` (`SELECT interface_hash_ref, semantic_id, provenance, payload FROM objects WHERE hash_ref = ?`), `store.go:692` (`SELECT 1 FROM objects WHERE hash_ref = ?`), and `journal.go:744/792/918/966` (each `SELECT ... FROM journal j JOIN objects o ON o.hash_ref = j.object_ref`, inspected as SELECTs). ZERO UPDATE, DELETE, DROP, or ALTER names the table. Schema: `objects(hash_ref TEXT PRIMARY KEY, interface_hash_ref/semantic_id/provenance TEXT NOT NULL, payload BLOB NOT NULL)`; `CREATE TRIGGER` in non-test `host/`+`cmd/` = **0** (the pattern's only two hits are in `host/store/journal_test.go`, on OTHER tables — `journal`, `approval_claims` — a test-scope firing that doubles as the instrument control); `FOREIGN KEY\|ON DELETE\|ON UPDATE` in `schema.sql` = **0** (same-file control `NOT NULL` = **26**); dynamic SQL: **15** non-test `fmt.Sprintf(` sites in `host/store` (per-file: journal 3, scan 1, store 8, writer_lock 3), every one an `Error() string` body except `writer_lock.go:196` — the `busy_timeout(2000)` DSN pragma — and none builds a statement against `objects`. **THE INSTRUMENT FINDING, twice over. First, the reading that establishes immutability must be a positive enumeration read one by one, never a negative grep for the mutation keywords: `grep -rniE '\bUPDATE[[:space:]]+[a-z_]+[[:space:]]+SET\b'` over the same trees returns 0 INCLUDING tests, while `host/store/store.go` holds FIVE `ON CONFLICT(...) DO UPDATE SET` upserts (lines 618, 709, 759, 836, 978 — on `registry_name`, `transition_fn_ref/interpreter_ref`, and `head_key`; none on `objects`). The upsert spelling puts NO table name between `UPDATE` and `SET`, so a grep aimed at the UPDATE-statement form is blind to a real mutation path by construction: its zero on a genuinely immutable table and its zero on a five-upsert file are the IDENTICAL zero. (Control that the pattern itself works: the same grep against a synthetic `/tmp` file containing `UPDATE objects SET payload = ?;` and `DELETE FROM objects WHERE hash_ref = ?;` matched both lines.) Second, even a POSITIVE pattern is only true inside its scope: the adjacency pattern `FROM objects\|INTO objects\|UPDATE objects\|DELETE FROM objects` reads FIVE statements and misses the four `JOIN objects` reads — only the bare-token sweep with per-line classification is complete.** | `grep -rnwi 'objects' --include='*.go' --include='*.sql' host/ cmd/ \| grep -v '_test'` (41 lines, classified; the adjacency comparison re-run beside it); `sed -n '13,19p' host/store/schema.sql`; `grep -rni 'CREATE TRIGGER' --include='*.go' --include='*.sql' host/ cmd/`; `grep -ciE 'FOREIGN KEY\|ON DELETE\|ON UPDATE' host/store/schema.sql; grep -ci 'NOT NULL' host/store/schema.sql`; `grep -rc 'fmt\.Sprintf(' --include='*.go' host/store/*.go \| grep -v '_test'`; `grep -rniE '\bUPDATE[[:space:]]+[a-z_]+[[:space:]]+SET\b' --include='*.go' --include='*.sql' host/ cmd/ \| wc -l`; `grep -n 'DO UPDATE SET' host/store/store.go`; synthetic control: `printf 'UPDATE objects SET payload = ?;\nDELETE FROM objects WHERE hash_ref = ?;\n' > /tmp/round9b_synth.sql; grep -nE '\bUPDATE[[:space:]]+[a-z_]+[[:space:]]+SET\b\|\bDELETE[[:space:]]+FROM[[:space:]]+objects\b' /tmp/round9b_synth.sql` | Token sweep: 41 lines, the nine statements and one DDL exactly as enumerated in the claim (statement sites inspected in place: `sed -n '470,485p' host/store/store.go`, `sed -n '738,748p;786,796p;912,922p;960,970p' host/store/journal.go`); adjacency pattern: **5** (blind to the four JOINs); triggers: only `journal_test.go:247` and `journal_test.go:561`; `fk_cascade=0`, `notnull_control=26`; Sprintf per-file `3/1/8/3` = 15 non-test (the round-9b briefing said 16 — the measured count is 15, recorded as measured); `update_set_form=0`; upserts at the five cited lines; synthetic control matched lines 1 and 2, `rc=0`. |
| V48 | Round-9b document readings for the producer-API fix, scope: this document only. BEFORE the §3.2 edit, ZERO occurrences of `NewProducer` or `GenerateProof` — the objection's premise, confirmed by measurement; AFTER, both tokens present in §3.2's conceptual-surfaces block and §10.11's record. The same-file, same-instrument control `NewValidator` fires in BOTH runs while the check moves 0 → non-zero, so the before-zero was a measurement, not a broken pattern. | `grep -c 'NewProducer\|GenerateProof' design_docs/planned/w-validated-proven-evidence-boundary.md; grep -c 'NewValidator' <same file>` — run before the round-9b edits and re-run after | Before: check `0`, control `15`. After: check `12` matching lines, control `17` (every added line is the §3.2 surface itself, its dependency paragraph, AC5's constructor-refusal clause, the header and §10.11 records of adding it, or this row's instrument; the control moved 15 → 17 because §3.2's new producer lines and §10.11 name `NewValidator` beside it). |
| V49 | The `busy_timeout < ObjectReadTimeout` composition condition is REAL and NOTHING in the tree asserts the ORDERING, at `516836f`. (a) `const busyTimeoutMillis = 2000` at `host/store/writer_lock.go:179`, UNEXPORTED, feeding the DSN pragma at `:196`; its own comment states the composition informally ("2000 ms sits well below the daemon's 10 s read deadline so the context always wins") — prose, not an assertion. (b) The only existing pin is a VALUE pin: `host/store/context_read_test.go:209` asserts `busyTimeoutMillis` equals the pinned 2000 — it orders the constant against no consumer's deadline. (c) `ObjectReadTimeout` occurs in ZERO `.go` files repo-wide (it is this document's proposed constant, so no code CAN assert an ordering against it yet — which is why AC18 places the assertion in the constructor where both values first coexist); same-scope control counting a different token, `busyTimeoutMillis`, fires with 2 files. AC18's own base tokens also read zero: `ErrUnorderedTimeouts` 0 in `*.go`, accessor spelling `func (s *Store) BusyTimeout(` 0 in `host/` (control `busyTimeoutMillis` 5 lines in `host/`). (d) `host/daemon/handlers.go:299–302` states the gap verbatim: "Today that composition is safe only because busy_timeout (2s) is shorter than the 10s request deadline — an ORDERING nothing in this code asserts, not a guarantee." (e) Obligation (i): queue row 22's head at `design_docs/world-mission.md:3583` carries **0** `LANDED\|ITEM COMPLETE` tokens after `~~…~~` spans are stripped, while the identical instrument on LANDED row 21 returns **1**; and the single `--grep='lock-wait'` commit on `origin/dev`, `912009d`, is the `D-WORLD-23` ledger commit whose body assigns ownership TO row 22 ("row 22 keeps ownership of the lock-wait bound") — corroborating openness, not contradicting it. | `sed -n '175,183p' host/store/writer_lock.go; grep -n 'busyTimeoutMillis' host/store/writer_lock.go host/store/context_read_test.go`; `sed -n '200,215p' host/store/context_read_test.go`; `grep -rn 'ObjectReadTimeout' --include='*.go' . \| wc -l` with control `grep -rln 'busyTimeoutMillis' --include='*.go' . \| wc -l`; `grep -rn 'ErrUnorderedTimeouts' --include='*.go' . \| wc -l; grep -rn 'func (s \*Store) BusyTimeout(' --include='*.go' host/ \| wc -l; grep -rn 'busyTimeoutMillis' --include='*.go' host/ \| wc -l`; `sed -n '292,305p' host/daemon/handlers.go`; `sed -n '3583,3608p' design_docs/world-mission.md \| perl -0pe 's/~~.*?~~//gs' \| grep -cE "LANDED\|ITEM COMPLETE"` beside the row-21 control; `git log origin/dev --oneline --grep='lock-wait'` | Constant and pragma at the cited lines; value pin fires at `:209`; `ObjectReadTimeout` **0** (control 2 files); `ErrUnorderedTimeouts` **0** and accessor spelling **0** (control 5 lines); the handlers comment verbatim as quoted; row-22 tokens **0** vs row-21 control **1**; exactly one `lock-wait` commit, `912009d`, the ledger entry. |
| V50 | JSON nesting inside `DecodeProposal`'s 256 KiB cap: the panic/CPU half of the round-10 note is REFUTED and the unasserted-stdlib-internal half SUSTAINED. Probe: pure `[…]` nesting decoded into `interface{}` by BOTH `json.Decoder.Decode` and `json.Unmarshal`, under BOTH `GOTOOLCHAIN=go1.25.6` (pinned) and ambient go1.26.4 — classification IDENTICAL on every arm. go1.25.6 readings (Decode / Unmarshal): CONTROL-shallow, depth 10, 20 B → both `err=nil`, 375 µs / 6.5 µs; deepest ACCEPTED, depth 9,999, 19,998 B → both `err=nil`, 7.59 ms / 1.83 ms; just-over, depth 10,001, 20,002 B → both refuse `exceeded max depth`, 162 µs / 129 µs; deepest cap-fitting, depth 131,071, 262,142 B → both refuse, 506 µs / 133 µs. No panic on any arm; the worst case INSIDE the byte cap is a sub-millisecond refusal and the worst ACCEPTED case is single-digit milliseconds, bounded by construction (depth ≤ bytes/2). The refusing limit is `const maxNestingDepth = 10000` at `encoding/json/scanner.go:148` of the pinned toolchain's GOROOT — unexported, an implementation detail, not documented API. Controller ran the same four arms earlier this iteration with matching classifications (its timings: 7.03 ms worst-accepted, 418 µs worst-refused — timing noise only). | Probe source `/tmp/depth_probe_iter101/main.go` (four arms, both codecs, timed); `GOTOOLCHAIN=go1.25.6 go run .` and `GOTOOLCHAIN=go1.26.4 go run .`; `grep -n 'maxNestingDepth' "$(GOTOOLCHAIN=go1.25.6 go env GOROOT)/src/encoding/json/scanner.go"` | Four-arm table as in the claim, identical on both toolchains; `scanner.go:148: const maxNestingDepth = 10000`. |
| V51 | Round-11b premises, all four measured first-party at `516836f` before either reviewer fix was applied. **(a) `gpt5-6-sol`'s premise — the single-connection pool — TRUE:** `db.SetMaxOpenConns(1)` at `host/store/store.go:298`. Instrument-scope note, recorded rather than smoothed: the same grep over non-test `host/store/` returns **2** lines, and the second (`journal.go:639`) is a COMMENT quoting the call — one code site, two matching lines, and a bare count would have said two. **(b) `gemini-3-1-pro`'s premise — the write side takes no context — TRUE:** `func (s *Store) PutObject(o Object) error` at `store.go:444`, against the same-file control `func (s *Store) GetObject(ctx context.Context, …)` at `:468`, already context-threaded by item 18. So the read side was threaded and the write side was not; the asymmetry is real, not a doc omission. **(c) The withdrawn LIVE-PRAGMA claim costs nothing — the busy window is IMMUTABLE after `Open`:** non-test `host/`+`cmd/` code issues **0** runtime `PRAGMA busy_timeout` reads or writes; the window is applied per physical connection from the `_pragma` DSN parameter built at `writer_lock.go:196`, which the same-call control finds (**1**), so the zero is a measurement. A property that cannot change after construction is exactly a cacheable one. **(d) AC20's base tokens and the call-site count:** `func (s *Store) WriteObject(` reads **0** in `host/` against the same-shape control `func (s *Store) PutObject(` = **1**; and `.PutObject(` outside `host/store`, non-test, is **8** (`broker/approve.go` 3, `broker/broker.go` 3, `registry` 1, `transitionreg` 1) — **not the 10 the controller's first draft transcribed from an earlier epoch** — with the same-pipeline `GetObject` control returning **13**, reproducing V43's figure exactly and proving the pipeline, not the number, is what was shared. | `grep -rn 'SetMaxOpenConns(1)' --include='*.go' host/store/ \| grep -v _test`; `grep -n 'func (s \*Store) PutObject(\|func (s \*Store) GetObject(' host/store/store.go`; `grep -rn 'PRAGMA busy_timeout' --include='*.go' host/ cmd/ \| grep -v _test \| wc -l` with control `grep -c 'busy_timeout(%d)' host/store/writer_lock.go`; `grep -rn 'func (s \*Store) WriteObject(' --include='*.go' host/ \| wc -l` with control `grep -rn 'func (s \*Store) PutObject(' --include='*.go' host/ \| wc -l`; `grep -rn '\.PutObject(' --include='*.go' host/ cmd/ \| grep -v '^host/store/' \| grep -v _test \| awk -F: '{print $1}' \| sort \| uniq -c` and the same pipeline for `.GetObject(` through `wc -l` | (a) `store.go:298` code + `journal.go:639` comment; (b) both signatures at the cited lines, write context-free, read context-bearing; (c) runtime pragma **0**, DSN control **1**; (d) `WriteObject` **0**, control **1**; `.PutObject(` **8** with the per-file breakdown as claimed, `.GetObject(` control **13**. |

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
