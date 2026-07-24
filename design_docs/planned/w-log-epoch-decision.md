# w-log-epoch-decision — Log Epochs, Content-Addressed Transitions, Versioned Hashes

**Status**: **SETTLED — ready for M1** (D1 RATIFIED by Mark 2026-07-24, attended: **Option A + B-as-metadata hybrid** — see "Ratified Decision" below; D2 & D3 stood from 2026-07-23). The transition-log format is the frozen kernel; the charter guardrail required explicit human ratification — exercised.
**Date**: 2026-07-23
**Charter clause**: clause-1
**Verified against**: shipped `ailang` binary `v0.30.0-144-g07fbb29c5` (`ailang --version`)
**Traces to**: [DESIGN.md](../DESIGN.md) §19 open question 11; [REFERENCES.md](../REFERENCES.md) design-deltas 1, 2, 5
**Blocks**: `w-world-library-m1` (M1 freezes the transition-log format — these are the hardest things to retrofit)
**Estimated**: 0.5 day (this doc); M1 implications sized in M1's own plan

> **Revision note.** Revised after quorum round 1 (BLOCKED): authoritative replay now
> pins the content-addressed interpreter artifact (`interpreter: HashRef` in every log
> header); corpus attestation demoted to non-authoritative nomination — closes the
> corpus-equivalence soundness gap (gpt5-6-sol + gemini-3-1-pro). A finite conformance
> corpus cannot establish semantic equivalence over all untested/future programs, so
> epoch membership never licenses binary substitution during authoritative replay, and
> verify/typecheck caches are keyed to the exact interpreter artifact, not the epoch.

---

## Ratified Decision (Mark, 2026-07-24, attended — resolves the fork below)

**D1 = Option A + Option B's identity recorded as metadata.** Every log-entry header carries:

1. `interpreter: HashRef` — content hash of the **exact released `ailang` binary bytes**: the
   **authoritative replay pin**. Bit-for-bit determinism on the writing platform; matches 1.0's
   single-machine local-first scope (clause 2), where platform-lock costs nothing today.
2. `release: string` — the platform-independent source identity (`ailang version` tag + commit,
   verified available in ailang#471) — **compatibility metadata, NOT authoritative**.
3. `semanticsEpoch: int` — as drafted (epoch registry an object in the store).

**The M8 door, held open without a format migration:** when multi-node arrives, promoting the
portable identity (release + conformance-corpus gate) to replay-authoritative is a **ratified
policy change over fields that already exist in every header** — never a log rewrite. Until
then, cross-platform replay is explicitly out of guarantee.
**Hermeticity obligation → M1 acceptance:** replay-doubling — every M1 acceptance run replays
each recorded episode twice through the pinned binary and byte-compares; divergence fails M1
(answers gpt5-6-sol's unproven-hermeticity objection with a standing test instead of an
assumption).
**Archive obligation:** the pinned binary artifact must be retained (D1 carries the obligation;
location is an M1 operational detail, per the blast-radius table).

## Open Decision (parked for Mark — quorum round 2, 2026-07-23) — RESOLVED ABOVE, kept for the record

**D2 (content-addressed transition functions) and D3 (SHA-256 + tagged `HashRef`) are settled**
— they drew no objection across two quorum rounds. **D1's replay-pin identity is NOT settled.**
Round 1 closed the corpus-equivalence gap (any-attested-binary replay is non-deterministic) by
pinning the *exact* interpreter artifact. Round 2's two reviewers then split, in opposite
directions, on what "the interpreter" should be — the classic **determinism ⇄ portability** fork
of content-addressed replay. This is a frozen-kernel-format choice, so it goes to the human gate.

**The fork (three coherent options):**

| Option | Replay pin | Determinism | Portability | Cost of the objection it answers |
|---|---|---|---|---|
| **A. Exact-binary pin** (current draft D1) | content hash of the exact released `ailang` executable | strong *if* the binary is hermetic | **platform-locked** — a macOS-ARM64 log can't replay on Linux-AMD64; substitution forbidden (gemini's objection) | assumes binary bytes ⇒ semantics; unproven hermeticity (gpt5-6-sol's objection) |
| **B. Release/semantics-identity pin + corpus-gated local substitution** (gemini) | upstream release/commit (or a semantics-version) identity; host runs its own arch's build of that release | weaker — substitution is corpus-gated (**partly re-opens round 1's gap**) | **cross-platform** (needed for clause-2 local-first fleets / M8 multi-node) | portability |
| **C. Content-addressed replay-runtime closure** (gpt5-6-sol) | interpreter + ABI + shared libs + deterministic env (locale/tz) + fuel/timeout, all hashed; sanitized replay sandbox | strongest | still platform-specific (doesn't solve B's portability) | hermeticity |

**Why a human, not another headless round:** A and C maximize determinism at the cost of
portability; B restores portability by accepting *probabilistic* (corpus-tested) equivalence
instead of *proof* — which is exactly the tradeoff round 1 rejected. There is no headless-provable
"right" answer; it's a values call about what World's replay guarantee **means** (bit-for-bit on
one machine forever, vs portable-modulo-conformance across a fleet), bounded by what a *released*
`ailang` binary can expose across the §14 frozen-core boundary.

**Upstream dependency (routed to `sunholo-data/ailang`):** Option B is only available if a released
`ailang` exposes a **platform-independent semantics-version identity** (e.g. a spec/conformance
hash decoupled from the platform executable) that replay could pin. Whether that exists or is
feasible is an upstream language/tooling question — filed per the mission's two-channel guardrail.
(Also filed: the designer's round-1 finding that no stable canonical-form/AST serialization is
exposed, which forced D2 to hash raw source bytes.)

**Recommended framing for Mark:** 1.0 is single-machine (DESIGN.md §15), so **A or C** is viable
*now*; but the log format being frozen here must survive to M8 multi-node, so choosing A/C is a
deliberate bet that cross-platform replay is a future *extension* concern (re-pin under a new hash
`algo` tag per D3's migration story) rather than a day-1 format constraint. If the upstream
semantics-identity of Option B is cheaply available, B avoids that bet. This is the decision.

---

## Motivation

World's core guarantee is that the append-only transition log replays **bit-for-bit
forever** (DESIGN.md §3.4, §15). Three things threaten that guarantee if left implicit,
and all three are baked into the log/object *encoding* — the single most expensive thing
to change after M1 ships:

1. A log entry's meaning depends on the **evaluation semantics** it was written under.
   AILANG evolves upstream; World consumes the released binary and never links compiler
   internals (§14). Without pinning, a compiler release silently changes what old entries
   replay to. (Urbit solved this by freezing Nock; Temporal never solved it and pays with
   Worker-Versioning + Patching APIs forever; Event Sourcing calls it the schema-evolution
   tax.)
2. A log entry's meaning depends on **which transition function ran**. Name+version
   references are ambiguous the moment code is edited in place. (Unison's answer:
   content-address the code.)
3. Every reference in (2) and every object in the store depends on a **hash algorithm**.
   Git baked bare SHA-1 into identity; the SHA-256 migration has taken a decade and is
   still incomplete.

This doc decides all three so M1 can freeze the encoding. It is a decision doc: each
section states the decision, the rationale, and the rejected alternatives.

## High-Impact Decisions

| Decision | Why High Impact | Chosen By | Deadline | Change Cost |
|----------|-----------------|-----------|----------|-------------|
| D1: content-addressed interpreter pin (authoritative replay) + semantics epoch (compatibility metadata) in every log header | Wrong choice makes old entries unreplayable — or worse, non-deterministically replayable — after any compiler release | human (this doc) | design (pre-M1) | high |
| D2: transition functions content-addressed by canonical source bytes | Identity scheme is baked into every log entry and cache key | human (this doc) | design (pre-M1) | high |
| D3: SHA-256 with explicit `algo` tag in every hash and object envelope | Bare digests make algorithm migration a decade-long fork of history | human (this doc) | design (pre-M1) | high |
| Epoch-registry granularity (per-store vs global) | Affects multi-node story (M8) but not the M1 encoding | agent at M1 | compile | low |
| Binary-archive mechanics (where archived interpreter artifacts live) | Operational, not encoded in the log (the *obligation* to archive is D1; the *location* is not) | agent at M1 | runtime | low |

### Design Freeze (all resolved by this doc)

- [x] D1: `interpreter: HashRef` (content hash of the exact released binary bytes — the **authoritative replay pin**) and `semanticsEpoch: int` (compatibility metadata) in every log-entry header; epoch registry is an object in the store
- [x] D2: `transitionFn: HashRef` (content hash of canonical source) in every log-entry header
- [x] D3: `HashRef = { algo, digest }`; textual form `"sha256:<hex>"`; `algo` field in the object envelope

---

## Decision 1 — Log-epoch interpreter-semantics versioning

**Decision.** Every log-entry header carries **two interpreter facts with different
authority**:

1. **`interpreter: HashRef`** — the content hash of the **exact released `ailang` binary
   bytes** that define the entry's evaluation semantics. This is the **authoritative
   replay identity**: authoritative replay resolves `E.interpreter` to one exact archived
   artifact and executes with it. A version string does not uniquely identify executable
   bytes; the content hash does.
2. **`semanticsEpoch: int`** — a monotonic counter **decoupled from the binary version
   string**, naming the entry's **compatibility class**. Epochs are operational metadata
   (tooling, migration planning, "which binaries probably replay this"); they carry **no
   replay authority**.

An **epoch registry** (itself a content-addressed object in the store, updated only by
ordinary logged transitions) maps each epoch to the set of released binaries **nominated
as candidates** believed compatible:

```
EpochRecord { epoch, candidateBinaries: ["v0.30.0", "v0.30.1", ...], successor, frozen }
```

- **Authoritative replay**: to replay entry E, resolve `E.interpreter` to the exact
  content-addressed artifact and execute with it — nothing else. The daemon **archives
  the exact artifact** for every interpreter hash that appears in the log. If the
  artifact is unavailable, or cannot run on the host (architecture/OS), replay **FAILS
  EXPLICITLY** with a structured error (Axiom A11) — it never silently substitutes
  another binary, however plausibly compatible.
- **Attestation is non-authoritative nomination.** A candidate binary is nominated into
  epoch N by replaying a **conformance corpus** (a designated slice of the actual log,
  plus fixture transitions) and comparing result hashes against the recorded ones. Pass →
  append the version to `candidateBinaries` via a logged transition. Fail → the binary
  starts epoch N+1. But a finite corpus **cannot establish semantic equivalence over all
  untested or future programs**: two binaries in the same epoch may diverge on an edge
  case the corpus never exercises. Nomination therefore does **not** authorize
  substitution during authoritative replay, and does **not** license verify/typecheck
  cache reuse across distinct interpreter artifacts. What it *is* good for: choosing
  which current binary to write new entries with, planning migrations, and advisory
  "should replay this" tooling — all falsifiable, none load-bearing for determinism.
- **Epoch bump rule** (unchanged in mechanism, re-scoped in meaning): a new epoch is
  minted only when a released binary is observed (or documented upstream) to change
  **evaluation semantics** — not on every release. A bugfix release that replays the
  corpus identically is *nominated into the current epoch* instead of bumping it. The
  boundary this draws is a **migration boundary**, not an interchangeability proof.
- **Migration when semantics MUST change** (Temporal's lesson): history is **never
  patched**. The old epoch is marked `frozen`, its `successor` set; new entries are
  written under the new epoch. Old entries replay forever under an old-epoch binary. If a
  world needs to *move* to new semantics, that is a **rebase transition** — a new, logged
  transition (under the new epoch) recording "state S, verified equivalent to the old
  replay, adopted as the new baseline," with the equivalence evidence attached. This is
  the roll-forward discipline of §14 applied to the interpreter itself: no in-place
  history rewrite, no Temporal-style `patched()` branches inside transition code.
- The binary version string that *wrote* an entry is kept as `writtenBy` — **provenance
  only**, never replay identity. `interpreter` is its exact-bytes counterpart: same
  consuming-the-released-artifact relationship to the binary, but content-addressed and
  authoritative.

**Rationale.**
- **Determinism admits no equivalence assumption.** Bit-for-bit replay (A1, A2) is a
  claim about what one exact executable does with one exact input. Any scheme that lets
  replay choose among "equivalent" binaries stakes the guarantee on an equivalence that a
  finite corpus cannot prove — and that a single untested edge case falsifies. Pinning
  the exact artifact is the only replay identity that is *checkable by construction*
  rather than *believed by induction*. This mirrors Decision 2 exactly: what content
  addressing does for the transition function's source, the interpreter hash does for the
  interpreter's bytes.
- **The epoch still earns its place — as metadata.** Without it, every release would look
  like a semantic boundary and migration planning would degenerate into ad-hoc judgments
  about which binaries are "really" equivalent. The epoch makes the *believed*
  compatibility class explicit and machine-checked (via the corpus) — while the exact pin
  keeps that belief out of the trusted replay path. Archive cost is bounded by binaries
  *actually used to write entries*, not by every release ever.
- World cannot take Urbit's route (freeze the interpreter): the AILANG core is frozen *to
  World* (§14 — consume released binary, route gaps upstream) but the language itself
  evolves upstream on its own authority. Hashing and archiving the released binary's own
  bytes consumes the released artifact as-is — no compiler-internal hooks required.
- Attestation-by-replay turns "is this release *probably* semantics-preserving?" from a
  changelog-reading exercise into a mechanical check World can run itself — an honest
  answer to an operational question, not a soundness proof.

**Rejected alternatives.**
- *Epoch membership as an interchangeability license* (this doc's round-1 draft: "replay
  with any attested binary for that epoch"). Unsound: attestation over a finite corpus
  cannot establish equivalence over all programs, so "any attested binary" makes replay
  non-deterministic on untested edge cases — violating A1 and the bit-for-bit guarantee
  the doc exists to protect. Rejected on quorum objection (gpt5-6-sol + gemini-3-1-pro);
  the epoch survives only as a compatibility class.
- *Pin the exact binary VERSION STRING per entry as replay identity.* A version string
  does not uniquely identify executable bytes (rebuilds, dirty builds, tag reuse) — it
  names a release, not an artifact. The content hash of the binary is the sound form of
  this idea, and is what D1 now adopts. The string is kept only as `writtenBy`
  provenance.
- *Pick one canonical binary per epoch* (the weaker fix). Repairs determinism but leaves
  identity resting on a registry lookup rather than on the entry itself, and still keys
  caches by epoch. Subsumed by the content-addressed pin, which puts the exact artifact
  identity in the header where the rest of the encoding already puts it.
- *No versioning; assume ailang evaluation is stable.* Urbit could promise that by
  freezing Nock; World does not control AILANG and §14 forbids forking it. One upstream
  semantics change would strand the entire log.
- *Temporal-style patching inside transition code* (`if patched("fix-x") ...`). Pushes
  versioning complexity into every transition function forever; Temporal's own docs treat
  it as a pain point to migrate away from. World's immutable-world model makes the
  roll-forward rebase strictly cleaner.
- *Snapshot-and-truncate on every semantics change* (drop replayability of old entries).
  Violates §3.4 replayability outright; rejected on principle, though snapshots remain a
  *performance* tool (Urbit's scaling lesson) orthogonal to this decision.

**Prior art**: Urbit/Arvo (frozen-kernel replay stability — the requirement, not the
mechanism); Temporal Worker Versioning + Patching (what unversioned evolution costs);
Event Sourcing / Fowler (schema-evolution tax; version from day one); Jujutsu op-log
(meta-operations — here, epoch-registry updates — logged through the same mechanism as
domain operations).

## Decision 2 — Content-addressed transition functions

**Decision.** The log references transition functions **by content hash, not by
name+version**. Concretely:

- **What is hashed**: the **canonical source bytes** of the transition function's `.ail`
  module — the exact text stored in the content-addressed object store, normalized only
  at the byte level (UTF-8, `\n` line endings, no trailing whitespace). **Not** the typed
  AST, **not** an interface hash.
- **Identity is content-addressed, not input-addressed** (the Nix distinction taken
  explicitly): a transition function's identity is the hash of what it *is*, not of the
  inputs/toolchain that produced it. The toolchain that *interprets* it is pinned
  separately — by Decision 1's content-addressed interpreter artifact.
- **Relationship to Decision 1** (both are needed): `transitionFn` pins **what code**
  ran; `interpreter` pins **exactly what evaluates it**; `semanticsEpoch` tags the
  compatibility class. Authoritative replay identity of an entry is the pair
  `(transitionFn, interpreter)`; typecheck/`ai-check`/verify results are cached under
  that pair and are valid **forever** — both members are content hashes of immutable
  bytes, so nothing a cache row depends on can change under it. `semanticsEpoch` is
  retained as compatibility metadata on cache rows (useful for queries and migration
  planning) but is **not** part of cache validity: distinct interpreter artifacts never
  share cache rows, even inside one epoch. A new artifact invalidates nothing — it just
  starts new cache rows.
- The store's existing `interface hash` and `semantic id` fields (DESIGN.md §15) are kept
  as **metadata alongside** the source hash: interface hash powers compatibility queries
  ("which transitions still satisfy this signature?"), semantic id powers human/agent
  naming. Neither is replay identity.
- **Unison-pitfall mitigations**: (a) objects in the store *are* plain `.ail` text —
  `ailang check`, grep, diff, and editors work on them unmodified; no UCM-style bespoke
  tooling is required because names live in metadata, not in the identity layer. (b) Hash
  churn on refactor is accepted and harmless: a refactored function is a **new object**;
  old log entries keep pointing at the old hash, which the store retains (append-only).
  Caches keyed by hash are never invalidated by churn, only added to.

**Rationale.**
- Name+version references reintroduce ambiguity ("which v3 — before or after the force-
  push?") that the entire replay guarantee rests on eliminating.
- Canonical source bytes is the **only hashable form available across the §14 boundary**:
  a normalized/typed AST hash would require linking compiler internals or depending on an
  unstable internal serialization — the released binary exposes source in, results out.
  This is a boundary-driven decision, not merely a simplicity preference. (If a future
  released `ailang` ships a stable canonical-form command, adopting it is a new hash
  *algo* tag under Decision 3 — e.g. `ail-cast1:` — a coexistence event, not a redesign.)
- Content-addressing (vs input-addressing) is correct here because transition functions
  are pure values interpreted at run time — there is no build step whose inputs could
  diverge from the artifact. Nix needs input-addressing because building is expensive and
  impure at the edges; World's "build" is Decision 1's interpreter pin.

**Rejected alternatives.**
- *Name + version references.* Ambiguous under edits; makes cache invalidation heuristic;
  the exact failure Unison exists to fix.
- *Hash the typed/normalized AST.* Better churn behavior (comment/whitespace-only edits
  keep the hash — though byte-canonicalization already covers whitespace), but crosses
  the §14 boundary: no stable AST serialization is exposed by the released binary, and
  World must not link internals. Revisit only if upstream ships a stable canonical form.
- *Interface hash as identity.* Two functions with the same signature and different
  bodies would collide — replay would be free to run the wrong code. Kept as metadata.
- *Input-addressed identity (Nix-style).* Ties identity to toolchain inputs World doesn't
  control and doesn't need; loses dedup of identical sources arriving via different
  routes.

**Prior art**: Unison (content-addressed definitions; names as metadata; typecheck cached
forever; the plain-text-tooling and hash-churn pitfalls); Nix/Dolstra (input- vs
content-addressed distinction, taken explicitly); Git object model (content-addressed
store with structural sharing).

## Decision 3 — Explicit, versioned hash algorithm

**Decision.** **SHA-256**, with the algorithm identifier **encoded in every hash
reference and every object envelope** from day one:

- Every hash in the system is a `HashRef { algo: string, digest: string }`; the canonical
  textual/serialized form is **`"<algo>:<hex-digest>"`** — e.g.
  `sha256:9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08`. A bare,
  untagged digest is a **malformed reference** — readers MUST reject it, not guess.
- The object-store envelope carries `hash: HashRef` (and `interfaceHash: HashRef`) — the
  algorithm is an explicit **field of the encoding**, not a property of the database.
- Pleasing consistency with the round-1 revision: the same tagged `HashRef` now **also
  carries the interpreter-artifact hash** (Decision 1). Transition source, object
  envelopes, the log chain, and the released binary's own bytes all share one
  self-describing reference form — pinning the interpreter required no new identity
  machinery, and a future algorithm migration covers interpreter pins for free.
- The hash-chained log (TigerBeetle-style `prevEntryHash`) uses the same tagged form.
- **Migration story** (the load-bearing part): a future algorithm change (e.g. to
  `blake3:`) is a **mixed-algorithm coexistence**, not a rewrite: new objects are written
  under the new tag; old objects are **never re-hashed** — their identity is immutable;
  references across the boundary are legal (a `blake3:`-addressed entry may reference a
  `sha256:`-addressed function); readers dispatch on the tag and fail loudly on unknown
  algorithms. Git's decade of pain came from omitting exactly this tag.

**Rationale.**
- SHA-256 over BLAKE3: it is in the Go standard library (`crypto/sha256` — zero
  dependencies for the M1 Go host), hardware-accelerated on current CPUs, and universally
  verifiable by third-party tooling (`shasum -a 256`, SQLite extensions, git-SHA-256
  ecosystem). BLAKE3 is faster in software, but hashing is nowhere near the bottleneck
  for a laptop-scale SQLite store, and BLAKE3 costs a third-party Go dependency in the
  trusted base. Crucially, **the encoding makes this choice cheap to revisit** — if
  hashing ever shows up in a profile, `blake3:` coexists by design.
- Tag-in-the-reference (not only in the envelope) because references travel without their
  envelopes — log headers, cache keys, evidence, cross-world links. Every place a digest
  appears must be self-describing.
- Textual `algo:` prefix over binary multihash varints: the multihash *principle* (self-
  describing hashes) at full strength, but debuggable with `sqlite3` and grep — consistent
  with the plain-text-tooling stance of Decision 2. The textual form is the canonical one;
  a packed binary encoding may be added later as pure representation (it must round-trip
  to the same canonical text).

**Rejected alternatives.**
- *Bare hex digests* (Git's original sin). The entire motivation for delta 5; rejected.
- *BLAKE3 now.* Faster, but buys nothing at M1 scale and costs a dependency; the tagged
  encoding keeps the door open at near-zero switching cost.
- *Binary multihash (varint codes) as canonical form.* Compact but opaque in a SQLite/
  plain-text world; adopted in spirit (self-describing), rejected as the canonical
  encoding.
- *Algorithm version only in a store-level header* ("this whole DB is sha256"). Makes
  mixed-algorithm coexistence impossible — which is precisely the migration mode Git
  proved you need.

## Type sketch (compiler-checked)

The real source lives at [`sketches/logepoch.ail`](../sketches/logepoch.ail) and passes
`ailang check` **and** `ailang ai-check` (0 errors) with the shipped binary
`v0.30.0-144-g07fbb29c5` — see the Verification Log. Pasted verbatim:

```ailang
module sketches/logepoch

-- Decision-doc sketch for w-log-epoch-decision (clause-1).
-- These types freeze the SHAPE of the three decisions; M1 implements them.
-- Revised after quorum round 1: authoritative replay pins the exact
-- content-addressed interpreter artifact; epoch membership is operational
-- metadata and never licenses binary substitution.

-- Decision 3: every hash carries its algorithm explicitly.
-- Canonical textual form: "${algo}:${digest}" e.g. "sha256:9f86d0...".
export type HashRef = { algo: string, digest: string }

-- Decision 3: the object-store envelope names its hash algorithm as a field —
-- migration is mixed-algorithm coexistence, never a rewrite of history.
export type ObjectEnvelope = {
  hash: HashRef,
  interfaceHash: HashRef,
  semanticId: string,
  provenance: string
}

-- Decision 1: a semantics epoch is decoupled from the binary version string.
-- The registry maps each epoch to CANDIDATE binaries nominated by corpus
-- attestation as believed-compatible. Nomination is NON-authoritative: a
-- finite corpus cannot establish semantic equivalence over all programs, so
-- candidacy is operational metadata (tooling, migration planning, "which
-- binaries probably replay this") — it never authorizes substituting a
-- different binary during authoritative replay and never licenses cache
-- reuse across distinct interpreter artifacts.
export type EpochRecord = {
  epoch: int,
  candidateBinaries: list[string],
  successor: int,
  frozen: bool
}

-- Decisions 1 + 2 meet in the log header:
--   interpreter (content hash of the exact released ailang binary bytes)
--     pins WHICH artifact evaluates the entry,
--   transitionFn (content hash) pins WHAT code ran,
--   semanticsEpoch tags the compatibility class (operational metadata only).
-- AUTHORITATIVE replay identity is the pair (transitionFn, interpreter).
export type LogHeader = {
  entryIndex: int,
  semanticsEpoch: int,
  transitionFn: HashRef,
  interpreter: HashRef,
  prevEntryHash: HashRef,
  writtenBy: string
}

-- Canonical rendering of a hash reference (the on-disk / in-log text form).
export func renderRef(r: HashRef) -> string {
  "${r.algo}:${r.digest}"
}

-- Epoch lookup predicate: a registry record describes a log entry's
-- compatibility class iff the epochs match. This is operational grouping
-- ONLY. Authoritative replay ignores candidateBinaries and executes with the
-- exact archived artifact named by h.interpreter; if that artifact is
-- unavailable or unsupported on the host, replay FAILS EXPLICITLY (A11) —
-- it never silently substitutes another binary.
export func servesEntry(rec: EpochRecord, h: LogHeader) -> bool {
  rec.epoch == h.semanticsEpoch
}

-- The verify/typecheck cache key: results cache forever because both the
-- code (content hash) and the EXACT interpreter artifact that produced them
-- are pinned. semanticsEpoch is deliberately absent: it is compatibility
-- metadata, not validity — distinct interpreter artifacts never share cache
-- rows, even inside one epoch.
export func cacheKey(h: LogHeader) -> string {
  "${renderRef(h.transitionFn)}@${renderRef(h.interpreter)}"
}
```

## Conflict Surface

This doc touches **no** ailang parser/lexer/typechecker/codegen surface — it is a World
(consumer-side) encoding decision. The frozen-core boundary analysis (§14):

- **Coupling to the ailang core**: **released-artifact identity only** — two
  compiler-side facts are consumed, both from the released artifact as shipped: (a) the
  `ailang --version` string (as `writtenBy` provenance and `candidateBinaries` entries),
  and (b) the **content hash of the released binary's own bytes** (as the `interpreter`
  replay pin, with the exact artifact archived by the daemon). Hashing/archiving the
  released binary is *consuming the released artifact* — the same relationship as
  `writtenBy`, made exact — it links no compiler internals and depends on no internal
  serialization (AST, Core, iface), so it stays inside the §14 frozen-core boundary;
  Decision 2 explicitly rejects AST hashing *because* that would cross it.
- **What could break the scheme from upstream**: (a) a release that changes evaluation
  semantics — handled by design: corpus nomination fails, a new epoch is minted, and
  old entries keep replaying under their archived exact artifacts; (b) removal of
  the version string or of source-in/results-out CLI operation — would be a breaking
  change to ailang's public contract, routed upstream like any other need, not worked
  around.
- **Programs that MUST still work**: the existing sketches
  [`sketches/worldtypes.ail`](../sketches/worldtypes.ail) and
  [`sketches/transitions.ail`](../sketches/transitions.ail) — re-verified green under the
  same binary in the Verification Log (this doc adds types beside them, changes nothing).

## Implications for M1 (`w-world-library-m1`)

What M1 must implement and **freeze** as a result of this doc:

1. **Log-header schema** includes `semanticsEpoch: int`, `transitionFn: HashRef`,
   `interpreter: HashRef` (content hash of the exact released binary bytes —
   authoritative replay pin), `prevEntryHash: HashRef` (hash chain), `writtenBy: string`
   — shape per `LogHeader` in the sketch. Frozen at M1.
2. **Object envelope** carries `hash`/`interfaceHash` as tagged `HashRef`s (extends §15's
   field list with the algo tag). All digests serialized as `"<algo>:<hex>"`; readers
   reject untagged digests. Frozen at M1.
3. **SHA-256 via Go stdlib** in the M1 Go host; the digest function is dispatch-on-tag
   from day one even while only one tag exists.
4. **Epoch registry as store objects**, bootstrapped with epoch 1 ← the binary that
   ships M1 as its first nominated candidate; registry updates are ordinary logged
   transitions (Jujutsu discipline: meta-ops through the same log).
5. **Interpreter-artifact archive**: the daemon hashes the running binary at startup,
   stamps that hash into every entry it writes, and **archives the exact artifact** for
   every interpreter hash appearing in the log. Authoritative replay resolves the pin
   from the archive; artifact missing or unrunnable on the host → **structured explicit
   failure** (A11), never substitution. (Archive *location/mechanics* remain
   agent-deferred; the *obligation* is frozen here.)
6. **Transition functions stored as canonical source bytes** (UTF-8, `\n`, no trailing
   whitespace — M1 specifies the exact canonicalization) and referenced from the log only
   by `HashRef`.
7. **Replay-proof CI gate** (already in M1's charter) must assert the pair-pinning:
   authoritative replay resolves the binary from the entry's `interpreter` hash (never
   from the epoch registry), and verify-result caching is keyed by
   `(transitionFn, interpreter)` with `semanticsEpoch` stored as compatibility metadata.
8. **Conformance-corpus mechanism** for candidate nomination (non-authoritative): M1
   designates the fixture-transition set; full nomination automation may land M2, but
   the corpus and the recorded result hashes exist from M1 (they are just the
   replay-proof artifacts, reused). No M1 code path may consume `candidateBinaries` for
   authoritative replay or cache validity.
9. World/Proposal/Transition/Evidence types in AILANG (M1's deliverable) adopt `HashRef`
   wherever a digest appears — no `string` digests anywhere in the typed surface.

**Deferred to implementers** (agent may choose): exact canonicalization procedure details
(within the byte-level constraint above); SQLite column layout for `HashRef` (single text
column in canonical form vs two columns); epoch-registry object naming; binary-archive
location/mechanics.

## Axiom Compliance

| Axiom | Score | Justification |
|-------|-------|---------------|
| A1: Determinism | +2 | Authoritative replay resolves to ONE exact content-addressed interpreter artifact — no corpus-based equivalence assumption stands between an entry and its replay (round-1 draft's "any attested binary" was a genuine A1 hole; closed) |
| A2: Replayability | +2 | Pair-pinning `(transitionFn, interpreter)` — both content hashes of immutable bytes — plus the archive obligation is what makes bit-for-bit replay survivable long-term |
| A3: Effect Legibility | 0 | No effect-system impact |
| A4: Explicit Authority | 0 | No capability changes (epoch-registry updates use ordinary transition authority) |
| A5: Bounded Verification | +1 | Verify/typecheck results cache forever under `(transitionFn, interpreter)` — sound because both members are exact, never shared across distinct artifacts |
| A6: Safe Concurrency | 0 | No concurrency impact |
| A7: Machines First | +1 | Self-describing hashes and machine-checked (advisory) attestation replace changelog-reading judgment |
| A8: Minimal Syntax | 0 | No language surface |
| A9: Cost Visibility | 0 | Negligible (hash tag bytes) |
| A10: Composability | +1 | Mixed-algorithm coexistence; cross-epoch references compose |
| A11: Structured Failure | +1 | Unknown algo / untagged digest / missing-or-unrunnable interpreter artifact all fail loudly with structured errors — replay never silently substitutes |
| A12: System Boundary | +1 | Decision 2 is *shaped by* the §14 boundary (source-bytes hashing, no internals) |

**Net Score: +9** ✅ — no hard violations (A1/A3/A4/A7 all ≥ 0).

## Verification Log

| Claim | Check | Result |
|-------|-------|--------|
| Sketch types compile under shipped binary | `cd design_docs && ailang check sketches/logepoch.ail` | ✓ No errors (v0.30.0-144-g07fbb29c5) |
| Sketch passes ai-check | `ailang ai-check sketches/logepoch.ail` | ✓ `"passed": true, "error_count": 0` (note: JSON is the default output; there is no `--json` flag on this binary) |
| Existing sketches unbroken | `ailang check sketches/worldtypes.ail && ailang check sketches/transitions.ail` | ✓ No errors, both |
| §15 store fields as claimed (hash, interface hash, semantic id, provenance) | Read DESIGN.md §15 | ✓ Confirmed verbatim |
| §14 boundary as claimed (released binary only, no internals, no fork) | Read DESIGN.md §14 hard boundary 2 | ✓ Confirmed |
| Open question 11 wording matches scope of this doc | Read DESIGN.md §19 item 11 | ✓ Confirmed (epoch header + content-addressed fns + decide-before-M1) |
| Deltas 1/2/5 verbatim as cited | Read REFERENCES.md deltas | ✓ Confirmed; delta 3 (input- vs content-addressed) taken inside Decision 2 |
| No existing planned doc covers this | `design_docs/planned/` did not exist before this doc; searched DESIGN/REFERENCES for prior decisions | ✓ First planned doc in this repo |
| Revised sketch (interpreter pin, round-1 fix) passes the CI gate | `./scripts/verify_ail.sh` from repo root | ✓ `"passed": true` for all modules; `checked 3 module(s)` (v0.30.0-144-g07fbb29c5) |
| Revised sketch passes plain check | `cd design_docs && ailang check --relax-modules sketches/logepoch.ail` | ✓ No errors found! |

## References

- [DESIGN.md](../DESIGN.md) — §3.4 replayability, §14 frozen-core boundary, §15 physical architecture, §17 M1, §19 open question 11
- [REFERENCES.md](../REFERENCES.md) — deltas 1, 2, 3, 5; Unison; Nix (Dolstra); Urbit/Arvo; Git object model; Jujutsu op-log; Temporal; TigerBeetle; Event Sourcing (Fowler)
- Prior art URLs live in REFERENCES.md (all verified live 2026-07-23)
