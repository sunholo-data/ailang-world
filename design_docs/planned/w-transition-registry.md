# w-transition-registry — Stable, Snapshot-Readable World Transition Catalog

**Status**: Planned
**Item**: clause-3 transition registry; prerequisite of `w-mcp-projection` P6.B
**Estimate**: 3.5 World days in three independently mergeable milestones
**Verified against**: repository HEAD `b0f323a`; AILANG v0.30.0 at
`/tmp/ailang-v0300/ailang`
**Date**: 2026-08-11

> **Scope truth.** This document designs the registry that projection reads. It does not design,
> mount, or encode MCP/A2A, resolve session credentials, or dispatch a transition through the
> coordinator. The minimum result is an immutable registry snapshot API, stable transition IDs,
> typed schema metadata, capability/effect declarations, and the broker-side snapshots and guard
> needed to consume those declarations without inventing another policy engine. TR.B supplies the
> descriptor-bound invocation mechanism; declaration honesty holds only where execution uses it,
> and becomes an enforced prerequisite for P6.B only when TR.C's structural binding gate is green.

## Motivation

`DESIGN.md` §11.1 requires discovery of the transitions available to a particular session as
typed schemas, projected through MCP. The quorum-cleared `w-mcp-projection` design cannot implement
that requirement because no transition catalog exists under `host/`, `world/`, or `cmd/`.

The repository already has three adjacent concepts which must remain distinct:

1. `host/registry.Registry` is the interpreter-epoch registry.
2. `broker.Registry` maps exact effect names to executable effect handlers.
3. Transition source and records are addressed by `HashRef`.

None supplies the missing stable transition name, input/output schema, or declared authority.
This item adds that index without changing the frozen transition/log identity: a stable ID selects
one descriptor in a registry revision, while the descriptor pins the content-addressed transition
source and interpreter that execution and replay already understand.

Eight supplied premises need correction. The quoted discovery text is under current §11.1, not
§11.2 (which is the human UI section). At this HEAD the AILANG gate checks **11**, not 9, `.ail`
modules; its exact semantic totals are still 4 verified `world/` identities and 14 named tests.
Also, the broker implementation is present and its focused tests pass, although
`design_docs/implemented/w-effect-broker-m3.md` still has the stale header `Status: Planned`.
The projection charter's claim that item 4 is parked is therefore stale; the transition registry
is the remaining local prerequisite addressed here. Finally, `rg` is not an installed binary on
this rig or in CI; acceptance commands use the repository's installed `grep` convention (V16–V17).
Controller measurements also found three production `Session.Invoke` calls, all inside
`host/broker/publish_op.go`; no production caller of either exported session constructor; and no
current coordinator/daemon-to-broker dispatch path (V25–V27).

## Premises (hard constraints)

- **P1 — three registries, three names.** The new Go package is `host/transitionreg`, its public
  value is `Snapshot`, and its semantic ID is `world/transition-registry/v1`. It does not add a
  package-scope type named `Registry`.
- **P2 — identity is a pair.** `Descriptor.ID` is the stable logical identity used by projections;
  `Descriptor.TransitionFn` is the immutable content identity used by execution and replay. An
  update may retain the ID only by creating a new registry revision that pins the new hash.
- **P3 — schemas are metadata.** Canonical JSON input/output schemas describe the invocation
  surface. They neither authorize execution nor replace AILANG checking, transition verification,
  a content hash, or a broker decision.
- **P4 — the broker remains the policy owner.** Registry records declare an invocation
  requirement and an effect manifest. Visibility calls the landed `broker.Decide` law over an
  immutable capability snapshot. TR.B provides a descriptor-bound invoker that delegates every
  accepted effect to `broker.Session.Invoke`; TR.C structurally gates the future P6.B execution
  path so registry-mediated dispatch cannot bypass that invoker.
- **P5 — immutable revisions, one mutable head.** Registry revisions are canonical,
  content-addressed store objects. The existing generic registry-head table selects one revision.
  A compare-and-swap update prevents lost updates; no schema migration is required.
- **P6 — one read produces one value.** `ReadSnapshot` reads one head once, resolves that immutable
  object, fully validates it, deep-copies schema bytes and entries, and returns `{Head, Revision,
  Entries}` or an error. It never returns a partial snapshot or a lazy store-backed iterator.
- **P7 — replay authority is unchanged.** A descriptor pins `TransitionFn`, `Interpreter`, and
  `SemanticsEpoch`; invocation copies those pins into the existing commit path. Replay continues
  to resolve the source and interpreter from the recorded hashes, never from the current registry.
- **P8 — no `.ail` addition.** This item adds no pure World law. Transition programs and their
  types remain AILANG package artifacts; the new work is persistence, canonical wire data,
  snapshot isolation, and host effect confinement. Therefore `scripts/verify_ail.sh` remains 4
  verified identities, 11 swept modules, and 14 named tests.
- **P9 — no registration surface.** This item adds no REST route, CLI verb, or network endpoint.
  A later package-install/coordinator path may publish a validated next revision after its normal
  propose → verify → commit authorization; direct projection code receives only a reader.
- **P10 — both CI jobs constrain every milestone.** Each milestone must pass the AILANG gate and
  the Go build/test gate. The full Go test gate cannot run green in this sandbox because socket
  creation is denied; focused non-socket tests and `go build ./...` remain available here.

## High-Impact Decisions

| Decision | Chosen alternative | Rejected alternatives | Evidence |
|---|---|---|---|
| 1 | Stable ID + content hash pair | hash alone; mutable name alone | V3, V4 |
| 2 | `host/transitionreg` + existing object/head store | `host/registry`; broker map; daemon state; new `world/*.ail` | V1, V2, V6 |
| 3 | Eager immutable `Snapshot` | live iterator; repeated head reads; startup cache | consumer contract, V6 |
| 4 | Canonical JSON schemas as metadata | schemas as authority; Go reflection at discovery time | DESIGN §11.1; V4 |
| 5 | One access requirement + complete effect manifest; landed broker law | registry-owned policy engine; `Proposal.requiredCaps` trusted after the fact | V4, V5 |
| 6 | CAS-published immutable revisions, bytewise ID ordering | in-place rows; blind head overwrite; insertion order | V6 |
| 7 | Pin interpreter pair in descriptor; replay recorded pair | current epoch nomination at replay | V7 |
| 8 | Host adapter only; no new kernel module | everything in Go; everything in AILANG | S1–S3, DESIGN §14 |

### Design Freeze

- [ ] Do not name a new package-scope type `Registry`.
- [ ] Do not use a transition source hash as the MCP/A2A tool name.
- [ ] Do not let display text or JSON schemas grant authority.
- [ ] Do not compute schemas from mutable process state during discovery.
- [ ] Do not read the registry head more than once in one `ReadSnapshot` call.
- [ ] Do not let the parsed-snapshot cache become a startup snapshot: every `ReadSnapshot` reads
  the current head first, and every cache hit returns a deep copy keyed by that immutable head hash.
- [ ] Do not return a lazy iterator, mutable slice alias, or partial snapshot.
- [ ] Do not update a registry object in place or publish a head without expected-head CAS.
- [ ] Do not preserve a stable ID across a content change without a new registry revision.
- [ ] Do not implement capability matching in `host/transitionreg` or `host/projection`; call the
  landed broker decision law.
- [ ] Do not allow a registry-mediated transition to request an effect absent from its declared
  effect manifest when using the bound invoker, even when the session independently holds that
  effect capability; TR.C must keep the P6.B dispatch path on that invoker.
- [ ] Do not consult the current transition registry or epoch registry during replay.
- [ ] Do not add a `.ail` file, store table, REST route, CLI verb, MCP/A2A codec, or projection
  package in this item.
- [ ] Do not expose the registry writer to `host/projection`; projection receives `Reader` only.
- [ ] Do not change the v1 ID grammar, `TR-CJSON-1` rules, or any byte/count limit below without a
  new semantic ID and interface hash.
- [ ] Derive v1 golden revision bytes and the semantic/interface hash only from the frozen grammar,
  codec rules, field layout, and bounds below; never hand-edit either fixture or hash.

## Decision 1 — Stable Identity Is an ID/Hash Pair

One registered transition is a `Descriptor`, conceptually:

```text
Descriptor {
  ID                 string
  TransitionFn       HashRef
  Interpreter        HashRef
  SemanticsEpoch     int64
  InputSchema        canonical JSON
  OutputSchema       canonical JSON
  Access             EffectRequirement
  DeclaredEffects    []EffectRequirement
  Title, Description string
}
EffectRequirement { Effect string; Scope string; Cost int64 }
```

`ID` is the stable surface identity. V1 accepts exactly the anchored regular language
`^[a-z0-9](?:[a-z0-9_-]{0,30}[a-z0-9])?(?:[./][a-z0-9](?:[a-z0-9_-]{0,30}[a-z0-9])?)*$`.
The total ID is 1–128 bytes inclusive and every `.`- or `/`-delimited segment is 1–32 bytes
inclusive. Because the alphabet is ASCII, byte length equals character length. Leading/trailing
`_` or `-`, empty segments, uppercase, whitespace, non-ASCII, and every other byte are invalid.
Ordering is lexicographic comparison of the unsigned bytes from left to right; at the first
difference the smaller byte sorts first, and if one ID is an exact prefix the shorter sorts first.

`TransitionFn` is the executable content identity. Re-registering changed source under the same ID
creates revision N+1 with a different `TransitionFn`; revision N remains addressable and unchanged.
Deleting an ID means omitting it from N+1, never deleting its earlier object. Renaming creates a new
ID and removes the old ID in the next revision. Aliases and deprecation windows are deferred.

**Alternatives.** Hash-only identity makes every source edit rename a tool and violates the
consumer's stable-ID contract. Name-only identity cannot prove which bytes ran. A pair retains a
diffable surface while keeping execution content-addressed. V3 shows hashes are already the replay
identity; this decision adds rather than replaces the missing logical identity.

## Decision 2 — Host Package, Immutable Object, Generic Head

The package is `host/transitionreg`; `host/registry` remains exclusively the epoch registry and
`broker.Registry` remains exclusively the effect-handler map. The semantic object ID is
`world/transition-registry/v1`; its interface hash is the hash of the frozen v1 wire-schema
identifier, not a transition's input/output schema.

The exact interface-hash preimage is the following ASCII byte string with no newline:
`world/transition-registry/v1|TR-CJSON-1|revision:{entries,interfaceHash,parent,revision,semanticID}|descriptor:{access,declaredEffects,description,id,inputSchema,interpreter,outputSchema,semanticsEpoch,title,transitionFn}|id:lower-ascii-1..128-segments-1..32|schema:raw262144-canonical65536|entries:1024|revision:raw16777216-canonical8388608`.
`hashref.SumSHA256` produces
`sha256:743f39f470bf354ebab0ab196598b5ba72db80463d833325cb7672249d4734ac`.
Revision and descriptor keys are exactly the lower-camel-case names enumerated in that preimage;
unknown or missing keys are invalid. `revision` and `semanticsEpoch` are non-negative signed 64-bit
decimal integers; refs and hashes use the existing `HashRef` JSON string representation;
requirements encode exactly `{effect,scope,cost}`, with non-negative signed 64-bit `cost`.

A revision encodes `{semanticID, interfaceHash, revision, parent, entries}` as deterministic JSON
under `TR-CJSON-1` (Decision 4). Entries are stored already sorted by `ID`; schema objects have
recursively sorted keys. Encode → Decode → Encode must reproduce identical bytes. Golden revision
bytes and the v1 interface hash are derived artifacts of this frozen specification. The object is
written through `PutObject`; the existing generic head names it.

`TestCodecGoldenRoundTrip` must assert byte equality against a committed literal golden (and a
literal interface-hash string), in addition to round-trip idempotence. A golden derived at test
time from the encoder or its semantic-ID constant cannot detect indentation, key-order, or tag
mutations because it changes alongside the mechanism it is meant to freeze.

`host/store` gains `CompareAndSetRegistryHead(name, expected, next)`. In one SQLite transaction it
requires the current head to equal `expected` (zero means absent), requires `next` to name an
existing object, and changes exactly that named row. Conflict returns a typed error carrying both
heads. The blind `SetRegistryHead` remains for compatibility but transition-registry publication
must not call it. No DDL changes.

**Alternatives.** Reusing `host/registry` would conflate interpreter compatibility with callable
World behavior. Reusing `broker.Registry` would conflate authority-neutral descriptors with live
effect handlers. Daemon memory would not survive restart. A new `world/*.ail` kernel module would
make host persistence an effectful pure-core concern and grow the frozen kernel for an adapter.

## Decision 3 — Snapshot Is an Eager, Self-Contained Value

`Reader` exposes only `ReadSnapshot(context.Context) (Snapshot, error)`. One call:

1. reads `world/transition-registry/v1`'s head exactly once;
2. loads that head's immutable object;
3. verifies object hash, semantic ID, interface hash, parent/revision rules, canonical encoding,
   ID grammar/order/uniqueness, all non-zero refs, schemas, and requirements;
4. deep-copies entries and schema bytes; and
5. returns `Snapshot{Head, Revision, Entries}`.

Every failure returns a zero `Snapshot` plus a typed error. The returned entries cannot observe a
later head update. `Lookup(ID)` binary-searches the frozen sorted slice; `List()` returns a copy.

This makes the projection rule enforceable: a request stores one registry `Snapshot` value and one
broker `CapabilitySnapshot` value in its request state. It may not call either provider again.
The two snapshots have independent epoch labels (`Head` and `CapabilitySnapshot.Epoch`); they need
not be captured in one database transaction because neither is re-read within that request.

The host adapter may maintain a thread-safe cache of the parsed, validated structure keyed by its immutable object hash. `ReadSnapshot` must still read the head dynamically on every call, but may return a deep copy from this cache if the head's hash is already known, bypassing store I/O and re-parsing for unchanged heads.

**Alternatives.** A live iterator can mix revisions. Re-reading the head for each schema can mix
epochs. A startup cache violates P6.B's requirement that the next request observe a newer head.

## Decision 4 — Typed Schemas Are Canonical Metadata

Input and output schemas are canonical JSON objects embedded in the immutable registry revision.
Registration receives them from the checked AILANG package/interface artifact produced before the
coordinator proposes the registry update; it does not reflect on a live Go function. V1 validation
requires `TR-CJSON-1` and the bounds below, but deliberately does not implement a second AILANG
typechecker or a general JSON-Schema validator.

`TR-CJSON-1` is the named v1 canonical JSON profile. It is implemented dependency-free with Go's
standard-library `encoding/json.Decoder` token stream using `UseNumber`, plus registry-owned
validation and encoding; plain `json.Unmarshal`/`json.Marshal` is forbidden for canonicalization.
The validator first requires UTF-8 input, scans tokens while retaining object-member names, rejects
duplicate names at every depth, rejects invalid escapes and any escaped or literal Unicode
surrogate code point, and requires each schema root to be an object. Strings are Unicode scalar
sequences and are **not normalized**: canonically equivalent NFC/NFD strings remain distinct.
Encoding emits UTF-8 directly, escapes only `"`, `\\`, and U+0000–U+001F (lowercase `\u00xx`, with
the standard short escapes for backspace, tab, LF, FF, and CR), and does not HTML-escape `<`, `>`,
or `&`. Object keys sort lexicographically by their unescaped UTF-8 bytes; arrays retain order.
The literals are exactly `true`, `false`, and `null`.

Numbers are parsed from the RFC 8259 JSON number grammar without binary floating point and emitted
as one normalized arbitrary-precision decimal: no leading `+`, no exponent, no leading integer
zero, no trailing fractional zero, and `-0` becomes `0`; the decimal point appears only when a
non-empty fractional part remains. Thus `1`, `1.0`, and `1e0` all encode as `1`. Before expansion,
the coefficient is limited to 1,024 decimal digits and the exponent magnitude to 10,000; after
normalization a number token is limited to 16,384 bytes. Exceeding any limit is invalid.

Each raw schema input is at most 262,144 bytes **before** parsing and its `TR-CJSON-1` form is at
most 65,536 bytes **after** canonicalization. A revision has at most 1,024 entries. Its raw encoded
input is at most 16,777,216 bytes before decoding and its complete canonical encoding is at most
8,388,608 bytes after canonicalization. Both input and output schemas count independently toward
the complete revision bound. All limits are inclusive; validation occurs before hashing or store
write, and decode enforces both raw and canonical limits before accepting an object.

The hashes of the registry revision and transition source make metadata changes auditable. They do
not make metadata authoritative. Invocation authority is `ID → descriptor in captured snapshot`,
then broker authorization, then normal propose → verify → commit. A forged or stale schema cannot
authorize a missing ID, substitute a different `TransitionFn`, or bypass verification.

**Alternatives.** Treating schemas as authority contradicts the consumer design. Generating them
from reflection at request time is neither stable nor package-language-neutral. Storing only schema
hashes would force projection to perform extra reads; embedding bounded canonical values makes one
snapshot self-contained.

## Decision 5 — Access Requirement and Effect Manifest Have Different Jobs

Each descriptor has exactly one `Access` requirement used for discovery/invocation admission, plus
an ordered, duplicate-free `DeclaredEffects` manifest covering every effect the transition may ask
the broker to perform. A requirement contains effect, exact scope, and non-negative cost; expiry and
budget belong to session grants, never registry metadata.

`broker.Session.CapabilitySnapshot(now)` copies the session's remaining grants under its mutex and
returns an immutable, monotonically incremented epoch. `broker.Allows(snapshot, requirement)` feeds
each grant and a corresponding `EffectRequest` to the existing `broker.Decide`; it does not restate
effect, scope, expiry, or budget comparisons. Thus access refusal has the landed four branches:
effect name, scope, expiry, and budget. An absent matching grant is denial, not an error.

`transitionreg.Bind(snapshot, ID, session)` returns a restricted effect invoker containing only
that descriptor's manifest. When execution uses this mechanism, it refuses an undeclared request
before calling `Session.Invoke`; declared requests still pass through `Session.Invoke`, which
performs the real grant/debit/record pipeline. TR.B alone does not force a caller through the
mechanism. TR.C makes the future P6.B path be born bound by rejecting any raw session exposure,
construction, or `Session.Invoke` call outside `host/broker`.

The proposal created for invocation must copy the descriptor's `TransitionFn`, expected effect
names, and requirements. The bound invoker provides proposal/descriptor agreement checking before
handler dispatch and commit, but that refusal is guaranteed only for execution routed through it.
The pre-dispatch manifest guard is the mechanism that keeps the declaration honest;
`Proposal.expectedEffects` and `requiredCaps` alone are descriptive fields and cannot stop a later
undeclared effect. TR.C is therefore required to close the declaration-honesty claim.

**Alternatives.** Putting expiry/budget in the registry would freeze session state into global
metadata. Trusting a declaration without a pre-dispatch guard is unenforced prose. Reimplementing
matching in this package creates the second policy engine P6.B explicitly forbids.

## Decision 6 — Publication Is CAS, Append-Only, and Deterministically Ordered

`BuildNext(current, changes)` is deterministic and side-effect free Go code: validate all candidate
descriptors, apply replacements/removals by stable ID, sort bytewise by ID, set revision N+1 and
parent to the captured head, then encode. `Publish(expectedHead, next)` puts the immutable object
and CASes only the transition-registry head. A failed CAS leaves the unselected immutable object in
the object store; that harmless orphan is content-addressed and may be retried or collected later.

There is no public registration endpoint in this milestone. Tests and bootstrap code can call the
writer. Production registration belongs to the package-install/coordinator lane and must occur only
after normal authorization and verification; a projection dependency is typed as `Reader`, so it
cannot publish. Registry publication itself does not execute a transition.

**Alternatives.** Blind head updates lose concurrent changes. Mutable descriptor rows erase
history. Insertion order makes cards nondeterministic. Granting projection a writer expands a
read-only protocol adapter into authority.

## Decision 7 — Interpreter Epoch Changes Affect New Invocations, Not Replay

A descriptor pins all three execution selectors: transition source hash, interpreter hash, and
semantics epoch. A package update that needs a new interpreter publishes a new descriptor revision,
even if the stable ID and transition source remain unchanged. The epoch registry may advise which
interpreter to nominate while building that revision, but `ReadSnapshot` never silently rewrites a
descriptor from the current nomination.

At invocation, the coordinator copies the captured pins into the existing log/commit data. Replay
uses the entry's `TransitionFn` and `Interpreter`; `SemanticsEpoch` remains compatibility metadata.
Consequently a later registry or epoch-head move cannot redirect an old episode. Registry-mediated
replay is bit-for-bit to exactly the same extent as existing hash-pinned replay; reconstructing the
historical discovery card is deferred and would use the recorded registry head.

**Alternatives.** Resolving the current epoch candidate during replay violates the landed replay
contract. Storing only an epoch number is insufficient because the epoch is not the authoritative
artifact identity.

## Decision 8 — Pure/Host Boundary and Minimum P6.B Contract

This is not “everything in Go.” Transition implementations, their typechecking, and their package
interfaces remain AILANG artifacts, and the broker decision law remains sourced from the checked
AILANG sketch. Go owns only the host adapter: immutable-object encoding, SQLite head selection,
request snapshots, schema transport, and effect dispatch confinement.

This is not “everything in AILANG.” Store I/O, mutex snapshots, CAS, JSON transport, and broker
dispatch are effects at DESIGN §14's host boundary. Adding them to `world/` would violate S2 and
S3. The small deterministic functions in `host/transitionreg`—validation, canonical encoding,
next-revision construction, ordering, and lookup—receive exhaustive/property-style Go tests and
golden bytes. S1 does not require a new `world/` function when no new pure World invariant is being
introduced.

The minimum registry that unblocks P6.B is named `world/transition-registry/v1` and supplies:

- stable ID plus pinned transition/interpreter identity;
- embedded typed input/output schema metadata;
- access and effect declarations;
- deterministic ordered immutable snapshots and dynamic next-request head reads;
- broker capability snapshots using the landed predicate; and
- a descriptor-bound undeclared-effect guard before the broker, structurally enforced for the
  future P6.B dispatch path by TR.C.

Everything beyond that list is Deferred Scope.

### P6.B consumer-contract trace

| Consumer requirement | Satisfied by this design |
|---|---|
| read one registry and one capability snapshot | Decision 3 `Snapshot`; Decision 5 `CapabilitySnapshot` |
| exact ordered set by stable identity | Decisions 1 and 6; bytewise `ID` order |
| live session and required capability | session liveness remains resolver-owned; Decision 5 applies landed `broker.Decide` to captured grants |
| never mix epochs in one request | Decision 3 prohibits provider re-read and returns explicit head/epoch labels |
| ID is authority; text/schema metadata only | Decisions 1 and 4 |
| fail closed on absent transition/read error/denial/broker failure | Decisions 3 and 5 return zero snapshot or typed refusal; projection maps it to protocol form |
| no second projection policy engine | Decision 5 centralizes matching in `broker.Decide`; TR.B supplies the bound invoker and TR.C prevents the future P6.B path from receiving or constructing a raw session or calling `Session.Invoke` directly |

Unknown/expired session and upstream structured-error encoding remain P6.B responsibilities.
`host/transitionreg` returns typed Go errors; it does not import protocol packages.

## Milestones

### TR.A — Immutable descriptor, snapshot, and CAS publication (2 days)

- Add canonical descriptor/revision codec, validation, stable ordering, golden bytes, and every
  decode/publication refusal test.
- Add the generic store CAS method and concurrency/rollback tests; add no table.
- Add `Reader`, eager `Snapshot`, lookup/copy-isolation tests, and dynamic-head test.
- Merge criterion: both CI jobs green; no broker or daemon change; P8 totals unchanged. Before the
  merge, change AC1–AC4 and AC10 from their documented base-tolerant count arms to require exactly
  3, 3, 4, 2, and 1 listed tests respectively; record `MUT-DELETE-TR-A-TEST` RED for a deleted or
  renamed required test. This zero-tolerance activation is a named release-checklist item.

### TR.B — Capability snapshot and declared-effect confinement (1 day)

- Add broker capability snapshots with epoch increments on debit.
- Add `broker.Allows` as a thin caller of `Decide`, with all four refusal arms pinned.
- Add restricted registry-bound effect invocation and one test/mutation per undeclared-effect and
  broker refusal branch.
- Add the exact synthetic two-session fixture consumed later by P6.B, without MCP/A2A code.
- Merge criterion: both CI jobs green; projection can depend only on `Reader` plus immutable
  capability snapshots; all named mutations demonstrated RED and restored GREEN. Before the merge,
  change AC5–AC7 from their documented base-tolerant count arms to require exactly 2, 3, and 3
  listed tests; record `MUT-DELETE-TR-B-TEST` RED for a deleted or renamed required test. This
  zero-tolerance activation is a named release-checklist item.

### TR.C — the binding gate (0.5 day)

- Add `TestRegistryDispatchBindingBoundary`, a repository-wide standard-library `go/parser`/`ast`
  assertion over production `.go` files located from the module root, not a hand-maintained package
  list. Outside `host/broker`, it rejects any broker `Session` type exposure or construction,
  either exported session constructor, and any `Invoke` selector call; this makes a future
  coordinator, daemon, projection, or newly added `cmd/` package red unless it uses the
  descriptor-bound invoker.
- Inside `host/broker`, enumerate only the pre-registry calls in `publish_op.go`: two in
  `mintAttendedApproval` and one in `invokeAttendedPublish`. Pin a separate expected exemption count
  of exactly **3** and assert both the identities and count; a fourth call, an unknown location, or
  silently raising the exemption count is RED. This is an exact legacy carve-out, not a package
  wildcard.
- Use ASTs because a text scanner cannot distinguish code from prose and can match its own needle
  list. Any detector tokens used by the test are assembled rather than present as literal scan
  needles. The test assertion fails, rather than merely logging, because `scripts/verify_go.sh`
  runs `go test` without `-v`; normal `go test ./...` therefore enforces the boundary in CI.
- The production-file inventory is a filesystem walk that skips nested Go modules and asserts a
  non-zero floor, required anchors, test-file exclusion, and a `go list` superset cross-check.
  Detector identifiers are exact matches; in particular `InvokeAttendedPublish` is not an
  `Invoke` selector call. Hermetic positive and negative controls cover every detector mechanism.
- TR.C may land before or after TR.B, but TR.B does not close the declaration-honesty claim by
  itself. The `w-mcp-projection` P6.B prerequisite is satisfied only when TR.A and TR.B are merged
  and TR.C is green; its charter must record that exact condition before implementation starts.

TR.C changes no production dispatch because V25–V27 show that no coordinator/daemon broker path
exists yet. There are therefore no concrete coordinator or daemon files to modify: P6.B will build
the first registry-mediated path under the structural boundary.

All three milestones are independently mergeable and at most 2 days. TR.A alone has useful
persistence and read semantics; TR.B depends on TR.A and adds authority integration without
changing its wire format; TR.C is a structural CI boundary and may land on either side of TR.B.

## Files to Create/Modify

| File | Milestone | Purpose |
|---|---|---|
| `host/transitionreg/transitionreg.go` | TR.A | descriptor, requirements, reader/writer, snapshot |
| `host/transitionreg/codec.go` | TR.A | canonical v1 encoding and validation |
| `host/transitionreg/transitionreg_test.go` | TR.A/TR.B | identity, snapshot, publication, guard tests |
| `host/store/store.go` | TR.A | generic compare-and-set registry-head primitive |
| `host/store/store_test.go` | TR.A | CAS conflict, absent-head, rollback, concurrency tests |
| `host/broker/broker.go` | TR.B | immutable capability snapshot and restricted integration seam |
| `host/broker/decide.go` | TR.B | exported thin `Allows` adapter over `Decide` |
| `host/broker/broker_test.go` | TR.B | snapshot epoch and four denial-arm tests |
| `host/broker/invoke_boundary_test.go` | TR.C | repository-wide AST binding assertion and exact three-call legacy exemption |
| `design_docs/planned/w-transition-registry.md` | design/TR.C | this document and AC11 zero-tolerance activation |
| `design_docs/verification/w-transition-registry/trc-mutations.md` | TR.C | 23-mutation non-vacuity transcript |

No change to `world/`, `host/registry`, `host/store/schema.sql`, `host/replay`, daemon/CLI code,
protocol code, `scripts/verify_ail.sh`, or CI workflow.

## Conflict Surface

- **Registry naming.** `host/registry.Registry`, `broker.Registry`, and the new transition catalog
  can be confused in imports and reviews. The package/type/semantic-ID freeze makes each explicit.
- **Generic head writes.** A new CAS method touches `host/store`, a kernel-adjacent package. It must
  preserve `SetRegistryHead` behavior for the epoch registry and must never CAS the wrong name.
- **Concurrent publication.** Put-before-CAS can leave an unselected immutable object; blind retry
  must not overwrite a winner. Tests distinguish acceptable orphan data from selected state.
- **Canonical JSON.** Changing field order, escaping, schema normalization, or ID sort changes the
  registry hash and downstream card diffs. `TR-CJSON-1` uses only standard-library
  `encoding/json` tokenization plus local validation/encoding, so it adds no module dependency and
  does not touch the daemon dependency allowlist. Golden bytes and round-trip idempotence freeze v1.
- **Aliasing.** Returning internal schema byte slices would let a caller mutate a supposedly frozen
  snapshot. Copy-isolation tests cover construction, `List`, and `Lookup`.
- **Broker locking.** Capability snapshots share the session mutex with debit/dispatch. Copying must
  not deadlock by calling a lock-taking path while locked, and the epoch must advance exactly when
  remaining authority changes.
- **Authorization drift.** `Access` controls discovery while `DeclaredEffects` confines execution.
  TR.B's mechanism tests prove an authorized access capability does not authorize an undeclared
  effect and that a declared effect still fails without its own live grant; TR.C proves the future
  registry-mediated path cannot bypass that mechanism.
- **Proposal drift.** The later coordinator adapter must compare proposal pins/declarations with the
  captured descriptor. This item supplies the values and guard, not coordinator wiring.
- **Replay.** No current replay file is modified. A future attempt to resolve current registry state
  during replay would violate Decision 7 and the existing replay tests.
- **Schema size.** Embedded schemas can inflate every snapshot. V1 caps raw/canonical schemas at
  262,144/65,536 bytes, entries at 1,024, and raw/canonical revisions at
  16,777,216/8,388,608 bytes; external schema refs are deferred until measurement requires them.
- **P6.B boundary.** This item cannot map errors to MCP/A2A, assess session liveness, or mount routes.
  Claiming those behaviors here would hide work in the wrong milestone.

## Systemic-Issue Audit

| Pattern | Audit result |
|---|---|
| Same concept duplicated | prevented: stable IDs live only in transition snapshots; hashes remain store/replay identity |
| Authority duplicated | prevented: matching delegates to `broker.Decide`; registry schemas are metadata |
| Negative-result instrumentation | V1 uses a known-positive `host/registry` control in the same command |
| Vacuous new-package gates | count gates tolerate only 0 or the complete set at base; each milestone merge checklist permanently removes its 0 arm and proves deletion/rename RED |
| Refusal-set under-testing | all decoder, reader, publisher, access, and effect-manifest branches are enumerated below |
| Kernel gravity | no `world/` or DDL growth; one generic CAS primitive is the only store change |
| Unbounded work | schema/revision sizes are bounded; no network/subprocess/wait is introduced |
| User-facing usage gap | no user-facing surface ships; P6.B owns protocol self-description and usage docs |

## Deferred Scope

- MCP/A2A projection, error encoding, routes, and cards (`w-mcp-projection` P6.B).
- Session credential resolution and session-liveness lifecycle.
- A REST/CLI registration API or operator UI.
- Automatic compiler schema extraction if the upstream package artifact lacks canonical schemas.
- Multiple access requirements, conditional effects, wildcard scopes, policy expressions, aliases,
  deprecation periods, and localization.
- Registry replication, garbage collection of unselected objects, pagination, subscriptions, and
  push invalidation.
- Historical card reconstruction and recording the registry head in the frozen log format.
- Changes to AILANG kernel types or the required-check manifest.

## Acceptance Criteria

Every command below was run at unmodified HEAD `b0f323a`. A count gate accepts only 0 (the whole
milestone test set is not implemented) or the exact complete count (then runs it). It never keys a
later milestone to a directory created by an earlier one. The TR.A/TR.B/TR.C merge checklists above
require removal of the corresponding 0 arm, so test deletion cannot regress into the base-tolerant
state. Every delivered gate also has a compiling RED mutation below.

> **Controller repair (iter-70, measured — see V16).** As drafted, AC1–AC8 counted test names with
> `rg -c '^Test'`. There is **no `rg` binary on this rig or in CI**: `rg` here is a *shell function
> injected by the agent harness's shell snapshot*, so it exists only inside the measuring
> environment and is absent from `env -i` and from `ubuntu-latest`. The repository itself uses `rg`
> **nowhere** (0 occurrences across `ci.yml` and all six `scripts/*.sh`, which use `grep`
> throughout). The pre-feature arm masked this completely: with `host/transitionreg` absent, every
> AC returns rc=0 **without ever executing the `rg` branch**, so the recorded base measurements are
> true and say nothing about the arm that matters. Once the package exists the ACs would have gone
> RED for a missing tool rather than a missing test — the right colour for the wrong reason. All
> eight are rewritten to `grep -c '^Test'`, which is exact-equivalent on `go test -list` output
> (measured both arms: known-present name → **1**, absent name → **0**). The Premise Verification
> Log below is a *record of commands actually run* and is deliberately left as recorded.

1. **AC1 — identity, codec, and schema validation.** Command:
   `export PATH=/opt/homebrew/bin:$PATH; count=$(GOTOOLCHAIN=go1.25.6 go test ./host/transitionreg -list 'Test(DescriptorIdentityAndContentUpdate|CodecGoldenRoundTrip|DescriptorValidationRefusals)$' 2>/dev/null | grep -c '^Test' || true); test "$count" -eq 3 && GOTOOLCHAIN=go1.25.6 go test ./host/transitionreg -run 'Test(DescriptorIdentityAndContentUpdate|CodecGoldenRoundTrip|DescriptorValidationRefusals)$' -count=1`.
   Base observed: rc=0 because count=0 (package absent), not because of a directory gate. TR.A
   activation requires count=3; three named tests PASS.
2. **AC2 — eager one-head snapshot and copy isolation.** Command:
   `export PATH=/opt/homebrew/bin:$PATH; count=$(GOTOOLCHAIN=go1.25.6 go test ./host/transitionreg -list 'Test(ReadSnapshotReadsHeadOnce|SnapshotIsEagerAndCopyIsolated|ReadSnapshotRefusals)$' 2>/dev/null | grep -c '^Test' || true); test "$count" -eq 3 && GOTOOLCHAIN=go1.25.6 go test ./host/transitionreg -run 'Test(ReadSnapshotReadsHeadOnce|SnapshotIsEagerAndCopyIsolated|ReadSnapshotRefusals)$' -count=1`.
   Base observed: rc=0 because count=0 (package absent). TR.A activation requires count=3.
3. **AC3 — CAS publication and deterministic ordering.** Command:
   `export PATH=/opt/homebrew/bin:$PATH; count=$(GOTOOLCHAIN=go1.25.6 go test ./host/transitionreg -list 'Test(PublishCASConflictPreservesWinner|ConcurrentPublishHasOneWinner|StableIDByteOrder|PublishRefusals)$' 2>/dev/null | grep -c '^Test' || true); test "$count" -eq 4 && GOTOOLCHAIN=go1.25.6 go test ./host/transitionreg -run 'Test(PublishCASConflictPreservesWinner|ConcurrentPublishHasOneWinner|StableIDByteOrder|PublishRefusals)$' -count=1`.
   Base observed: rc=0 because count=0 (package absent). TR.A activation requires count=4.
4. **AC4 — generic store CAS preserves the epoch registry.** Command:
   `export PATH=/opt/homebrew/bin:$PATH; count=$(GOTOOLCHAIN=go1.25.6 go test ./host/store -list 'Test(RegistryHeadRoundTrip|CompareAndSetRegistryHead)$' 2>/dev/null | grep -c '^Test' || true); test "$count" -eq 2 && GOTOOLCHAIN=go1.25.6 go test ./host/store -run 'Test(RegistryHeadRoundTrip|CompareAndSetRegistryHead)$' -count=1`.
   Base observed: rc=0 because count=1 is the known existing control. TR.A activation removes the
   count=1 arm, requires count=2, and runs both tests.
5. **AC5 — capability snapshot and all landed denial arms.** Command:
   `export PATH=/opt/homebrew/bin:$PATH; count=$(GOTOOLCHAIN=go1.25.6 go test ./host/broker -list 'Test(CapabilitySnapshotEpochAndIsolation|AllowsUsesDecideAllFourDenials)$' 2>/dev/null | grep -c '^Test' || true); test "$count" -eq 2 && GOTOOLCHAIN=go1.25.6 go test ./host/broker -run 'Test(CapabilitySnapshotEpochAndIsolation|AllowsUsesDecideAllFourDenials)$' -count=1`.
   Base observed: rc=0 because count=0 (TR.B tests absent), independently of TR.A's directory.
   TR.B activation requires count=2. Delivered: both named
   tests PASS and enumerate effect-name, scope, expiry, and budget denials.
6. **AC6 — descriptor-bound declaration-honesty mechanism.** Command:
   `export PATH=/opt/homebrew/bin:$PATH; count=$(GOTOOLCHAIN=go1.25.6 go test ./host/transitionreg -list 'Test(GuardedSessionRefusesUndeclaredEffect|GuardedSessionStillRequiresBrokerGrant|ProposalDescriptorAgreementRefusals)$' 2>/dev/null | grep -c '^Test' || true); test "$count" -eq 3 && GOTOOLCHAIN=go1.25.6 go test ./host/transitionreg -run 'Test(GuardedSessionRefusesUndeclaredEffect|GuardedSessionStillRequiresBrokerGrant|ProposalDescriptorAgreementRefusals)$' -count=1`.
   Base observed: rc=0 because count=0 (package absent). TR.B activation requires count=3; the
   three named tests PASS with zero handler calls
   on every refusal.
7. **AC7 — exact two-session consumer fixture and dynamic next-request source.** Command:
   `export PATH=/opt/homebrew/bin:$PATH; count=$(GOTOOLCHAIN=go1.25.6 go test ./host/transitionreg -list 'Test(TwoSessionExactOrderedSets|NextReadObservesNewHeadWithoutRestart|SingleRequestKeepsCapturedEpochs)$' 2>/dev/null | grep -c '^Test' || true); test "$count" -eq 3 && GOTOOLCHAIN=go1.25.6 go test ./host/transitionreg -run 'Test(TwoSessionExactOrderedSets|NextReadObservesNewHeadWithoutRestart|SingleRequestKeepsCapturedEpochs)$' -count=1`.
   Base observed: rc=0 because count=0 (package absent). TR.B activation requires count=3.
8. **AC8 — replay remains hash-pinned.** Command:
   `export PATH=/opt/homebrew/bin:$PATH; test "$(GOTOOLCHAIN=go1.25.6 go test ./host/replay -list 'Test(FixtureEpisodeReplaysBitForBit|EpochRegistryCandidateCannotRedirect)$' | grep -c '^Test')" -eq 2 && GOTOOLCHAIN=go1.25.6 AILANG_BIN=/tmp/ailang-v0300/ailang go test ./host/replay -run 'Test(FixtureEpisodeReplaysBitForBit|EpochRegistryCandidateCannotRedirect)$' -count=1`.
   Base observed: rc=0, package PASS in 2.197s. Delivered result remains rc=0 without importing
   `host/transitionreg` in production replay code.
9. **AC9 — AILANG gate totals do not move.** Command:
   `export PATH=/opt/homebrew/bin:$PATH; out=$(AILANG_BIN=/tmp/ailang-v0300/ailang ./scripts/verify_ail.sh 2>&1); rc=$?; m=$(printf '%s\n' "$out" | grep -c '4/4 required world/ identities verified across 11 module(s)'); t=$(printf '%s\n' "$out" | grep -c 'all 14 required named tests pass'); s=$(printf '%s\n' "$out" | grep -c 'world package gate PASSED: 9/9 steps'); [ "$rc" -eq 0 ] && [ "$m" -eq 1 ] && [ "$t" -eq 1 ] && [ "$s" -eq 1 ]`.
   Base observed: rc=0; 4/4 verified identities across 11 modules, 14 named tests, all 9 package
   steps, and final PASS. Delivered output has the same 4/11/14 totals.
10. **AC10 — build plus focused Go packages.** Command:
    `export PATH=/opt/homebrew/bin:$PATH; GOTOOLCHAIN=go1.25.6 AILANG_BIN=/tmp/ailang-v0300/ailang go build ./... && GOTOOLCHAIN=go1.25.6 go test ./host/store -run 'TestRegistryHeadRoundTrip$' -count=1 && GOTOOLCHAIN=go1.25.6 go test ./host/broker -run 'Test(ReplayReturnsRecordedBytesWithoutDispatch|ReplayGapNeverFallsBackToLive)$' -count=1 && count=$(GOTOOLCHAIN=go1.25.6 go test ./host/transitionreg -list 'TestCodecGoldenRoundTrip$' 2>/dev/null | grep -c '^Test' || true) && test "$count" -eq 1 && GOTOOLCHAIN=go1.25.6 go test ./host/transitionreg -count=1`.
    Base observed: rc=0; build, store, and broker focused tests PASS; transition-registry count=0
    because the package is absent. TR.A activation removes the zero arm, requires count=1, and runs
    the full package. Full `verify_go.sh` was also
    run at base with `GOTOOLCHAIN=go1.25.6`; it reached build successfully but test returned rc=1
    because this sandbox denies loopback listeners (`bind: operation not permitted`). CI must run
    `./scripts/verify_go.sh` green before any milestone merges.
11. **AC11 — structural registry-dispatch binding boundary.** Command:
    `export PATH=/opt/homebrew/bin:$PATH; count=$(GOTOOLCHAIN=go1.25.6 go test ./host/broker -list 'Test(ReplayReturnsRecordedBytesWithoutDispatch|RegistryDispatchBindingBoundary)$' 2>/dev/null | grep -c '^Test' || true); test "$count" -eq 2 && GOTOOLCHAIN=go1.25.6 AILANG_BIN=/tmp/ailang-v0300/ailang go test ./host/broker -run 'Test(ReplayReturnsRecordedBytesWithoutDispatch|RegistryDispatchBindingBoundary)$' -count=1`.
    Base observed this revision: rc=0, count=1 from the known existing replay control; the binding
    test is absent. TR.C activation removes the count=1 arm, requires count=2, and runs both tests.
    The count gate is independent of whether any future projection/coordinator directory exists.

## Non-Vacuity — Named RED Mutation for Every Gate

| Gate / refusal branch | Named compiling mutation | Required RED observation |
|---|---|---|
| AC1 stable ID | `MUT-ID-ACCEPT-EMPTY`: neuter empty/grammar rejection with `if false &&` | descriptor-refusal subtest names ID |
| AC1 zero transition ref | `MUT-ID-ZERO-FN`: neuter `TransitionFn.IsZero` guard | zero-ref subtest RED |
| AC1 zero interpreter ref | `MUT-ID-ZERO-INTERP`: neuter interpreter guard | zero-ref subtest RED |
| AC1 schema syntax | `MUT-SCHEMA-ANY-BYTES`: neuter JSON-object guard | malformed-schema subtest RED |
| AC1 schema bound | `MUT-SCHEMA-NO-LIMIT`: neuter size guard | oversize-schema subtest RED |
| AC1 canonical bytes | `MUT-CODEC-INDENT`: switch to indented JSON | golden-byte test RED while build stays green |
| AC2 missing head | `MUT-READ-EMPTY-OK`: return empty snapshot on absent head | missing-head subtest RED |
| AC2 store read error | `MUT-READ-SWALLOW`: replace injected read error with empty snapshot | read-error subtest RED |
| AC2 absent object | `MUT-READ-ABSENT-OK`: accept `ok=false` | absent-object subtest RED |
| AC2 object hash | `MUT-READ-NO-REHASH`: neuter payload hash check | corrupt-object subtest RED |
| AC2 semantic/interface IDs | `MUT-READ-ANY-TYPE`: neuter both identity guards, one at a time | corresponding typed-object subtest RED |
| AC2 partial/lazy alias | `MUT-SNAPSHOT-ALIAS`: retain caller/schema slices | copy-isolation test RED |
| AC2 one head read | `MUT-SNAPSHOT-REREAD`: read head once per entry | call-counter assertion RED |
| AC2 cache freshness/copy | `MUT-SNAPSHOT-CACHE-BYPASS`: return a cached snapshot without a fresh head read or without a deep copy, one at a time | head-call counter or copy-isolation assertion RED |
| AC3 parent/revision | `MUT-REVISION-SKIP`: set N+2 or wrong parent | revision-chain subtest RED |
| AC3 duplicate/order | `MUT-ORDER-INSERTION`: skip sort/dedup validation | ordering/duplicate subtests RED |
| AC3 expected-head conflict | `MUT-CAS-BLIND`: replace CAS predicate with unconditional upsert | loser overwrites winner; conflict test RED |
| AC3 next-object existence | `MUT-CAS-DANGLING`: neuter object-existence guard | dangling-head subtest RED |
| AC3 store/CAS error | `MUT-PUBLISH-SWALLOW`: return nil on injected Put/CAS errors, one at a time | matching publication subtest RED |
| AC4 wrong registry name | `MUT-CAS-EPOCH-HEAD`: hardcode epoch registry name | `TestCompareAndSetRegistryHead/epoch_registry_isolation` RED; the epoch round-trip does not call CAS and remains green |
| AC5 capability isolation (session→snapshot) | `MUT-CAPS-ALIAS-SESSION`: retain the session grant slice in `CapabilitySnapshot` | later-debit isolation assertion RED |
| AC5 capability isolation (snapshot→caller) | `MUT-CAPS-ALIAS-CALLER`: return the snapshot grant slice from `Grants` | caller-mutation isolation assertion RED |
| AC5 epoch update | `MUT-CAPS-STATIC-EPOCH`: do not increment epoch after debit | epoch assertion RED |
| AC5 effect-name denial | `MUT-ALLOW-NAME`: force only `denied:effect-name` to allowed | named denial arm RED |
| AC5 scope denial | `MUT-ALLOW-SCOPE`: force only `denied:scope` to allowed | named denial arm RED |
| AC5 expiry denial | `MUT-ALLOW-EXPIRED`: force only `denied:expired` to allowed | named denial arm RED |
| AC5 budget denial | `MUT-ALLOW-BUDGET`: force only `denied:budget` to allowed | named denial arm RED |
| AC6 absent transition | `MUT-BIND-MISSING`: return an empty permit for missing ID | absent-transition subtest RED |
| AC6 undeclared effect | `MUT-EFFECT-UNDECLARED`: neuter manifest-membership guard | handler counter becomes 1; test RED |
| AC6 declared but ungranted | `MUT-EFFECT-BYPASS-BROKER`: invoke handler directly | denial/zero-dispatch test RED |
| AC6 proposal function mismatch | `MUT-PROPOSAL-FN`: neuter function-ref agreement | named agreement arm RED |
| AC6 proposal requirements mismatch | `MUT-PROPOSAL-CAPS`: neuter requirement agreement | named agreement arm RED |
| AC6 proposal effects mismatch | `MUT-PROPOSAL-EFFECTS`: neuter manifest agreement | named agreement arm RED |
| AC7 per-session filtering | `MUT-SESSION-UNION`: return union to both snapshots | exact-set test RED |
| AC7 startup cache | `MUT-STARTUP-CACHE`: memoize first registry snapshot | next-read test RED |
| AC7 split request epoch | `MUT-CAPS-REREAD`: obtain capabilities again after registry barrier | captured-epoch test RED |
| AC8 registry redirects replay | `MUT-REPLAY-CURRENT-REGISTRY`: resolve source/interpreter from current head | existing replay authority test RED |
| AC9 new `.ail` | `MUT-AIL-EMPTY-MODULE`: add an empty `world/transitionregistry.ail` without manifest update | repaired AC9's printed-total assertion has `modules11=0`; the underlying script remains rc=0 |
| AC10 package build | `MUT-GO-CODEC-TAG`: change the v1 semantic ID constant to `world/transition-registry/v2` | full new-package suite RED while all production code still compiles |
| TR.A activated test inventory | `MUT-DELETE-TR-A-TEST`: delete or rename one required AC1–AC4/AC10 test after removing the base arm | exact activated count RED; zero is no longer accepted |
| TR.B activated test inventory | `MUT-DELETE-TR-B-TEST`: delete or rename one required AC5–AC7 test after removing the base arm | exact activated count RED; zero is no longer accepted |
| AC11 fourth raw invoke inside broker | `MUT-BINDING-FOURTH-INVOKE-IN`: add a fourth production `Session.Invoke` inside `host/broker` | exact exemption count RED while production compiles |
| AC11 raw invoke outside broker | `MUT-BINDING-FOURTH-INVOKE-OUT`: add an `Invoke` selector outside `host/broker` | outside-clean assertion RED |
| AC11 moved exemption identity | `MUT-BINDING-MOVE-SITE`: move one exempt call into a new helper while retaining count 3 | identity-set assertion RED while count remains green |
| AC11 live/replay constructors | `MUT-BINDING-CTOR-LIVE`, `MUT-BINDING-CTOR-REPLAY` | matching hermetic control and outside-clean assertion RED |
| AC11 type/alias/dot import | `MUT-BINDING-SESSION-TYPE`, `MUT-BINDING-ALIAS-IMPORT`, `MUT-BINDING-DOT-IMPORT` | matching detector mechanism RED |
| AC11 raised exemption | `MUT-BINDING-RAISE-COUNT`: silently raise the enumerated exemption count | pinned count/identity assertion RED |
| AC11 detector neutered | `MUT-BINDING-IF-FALSE-{INVOKE,CTOR-LIVE,CTOR-REPLAY,SESSION-TYPE,DOT-IMPORT}` | the corresponding hermetic positive control sees zero; assertion RED while mutant compiles |
| TR.C activated test inventory | `MUT-DELETE-TR-C-TEST`: delete or rename the required AC11 test after removing the base arm | exact activated count RED; the known control alone is no longer accepted |

TR.B rule-3j branch inventory (enumerated from the implementation diff):

| Branch | Named compiling mutation | Required RED observation |
|---|---|---|
| J1 negative manifest cost | `MUT-MANIFEST-NEG-COST-OK` | malformed-manifest negative-cost subtest RED |
| J2 duplicate manifest declaration | `MUT-MANIFEST-DUP-OK` | malformed-manifest duplicate subtest RED |
| J3 declaration scope is exact | `MUT-DECLARED-NAME-ONLY` | undeclared-scope subtest RED |
| J4 declaration cost is exact | `MUT-DECLARED-COST-ANY` | undeclared-cost subtest RED |
| J5 replay debit increments epoch | `MUT-EPOCH-LIVE-ONLY` | replay-debit epoch subtest RED |
| J6 `Allows` uses captured `Now` | `MUT-ALLOW-NOW-ZERO` | captured-now subtest RED |
| J7 bind preserves broker denial label | `MUT-BIND-COLLAPSE-LABEL` | exact-label subtests RED |
| J8 bind refuses zero snapshot | `MUT-BIND-EMPTY-SNAPSHOT-OK` | zero-snapshot subtest RED |
| J9 proposal compares interpreter | `MUT-PROPOSAL-INTERP` | interpreter-mismatch subtest RED |
| J10 proposal compares semantics epoch | `MUT-PROPOSAL-EPOCH` | semantics-epoch-mismatch subtest RED |
| J11 request propagates reader failure | `MUT-REQUEST-SWALLOW` | injected-read-error subtest RED |
| J12 allowed descriptors are detached | `MUT-REQUEST-ALIAS` | returned-descriptors-are-copies subtest RED |
| J13 allowed order is snapshot order | `MUT-REQUEST-REORDER` | order-is-the-snapshot-order subtest RED |
| J14 target bind error is propagated | `MUT-TARGET-BIND-SWALLOW` | target-bind-error subtest RED |

Eleven frozen refusal branches in Decisions 1, 2, and 4 have no named mutation above; the TR.A
sprint plan adds mutations for all eleven. Store errors
are driven by injected interfaces; SQLite CAS conflicts are driven by two handles/racers. Mutation
runs must record the failing test name, not only rc=1.

## Axiom Compliance

| Axiom | Score | Justification |
|---|---:|---|
| A1 Determinism | +2 | canonical bytes, bytewise IDs, immutable snapshots |
| A2 Replayability | +2 | descriptor and log pin the executable pair; replay ignores current heads |
| A3 Effect Legibility | +2 | complete declared-effect manifest plus broker records |
| A4 Explicit Authority | +2 | one access requirement and broker-owned decision; TR.B supplies bound refusal and TR.C structurally requires the future P6.B path to use it |
| A5 Bounded Verification | +1 | bounded schemas/revisions and existing verification path |
| A6 Safe Concurrency | +2 | CAS publication and immutable request snapshots |
| A7 Machines First | +2 | typed schema metadata and stable machine IDs |
| A8 Minimal Syntax | 0 | no language syntax change |
| A9 Cost Visibility | +1 | requirements include broker cost |
| A10 Composability | +1 | protocol-neutral reader consumed by MCP/A2A later |
| A11 Structured Failure | +2 | typed reader, conflict, lookup, and authority errors |
| A12 System Boundary | +2 | pure transition artifacts stay AILANG; host effects stay Go |

**Net: +19; hard axioms A1/A3/A4/A7 are positive.**

## Premise Verification Log

All commands were run on 2026-08-11 in this worktree. Commands shown here include the required PATH
prefix; output is quoted or counted rather than inferred.

| ID | Codebase claim | Command | Observed output |
|---|---|---|---|
| V0 | inspected revision | `export PATH=/opt/homebrew/bin:$PATH; git rev-parse --short HEAD` | `b0f323a` |
| V1 | transition registry absent; search instrument works | `export PATH=/opt/homebrew/bin:$PATH; printf 'transition_registry_hits='; rg -i -l 'transition[ -]?registry' host world cmd \| wc -l; printf 'host_registry_registry_hits='; rg -l 'Registry' host/registry \| wc -l` | `transition_registry_hits=0`; positive control `host_registry_registry_hits=2` |
| V2 | colliding registry/session types | `export PATH=/opt/homebrew/bin:$PATH; rg -n '^type Registry\|Registry maps\|type Session\|func NewSession\|type Capability' host/broker/broker.go host/broker/decide.go host/registry/registry.go` | epoch `Registry` at `registry.go:50`; handler `Registry` at `broker.go:35`; `Session` at :46; `NewSession` at :58; `Capability` at `decide.go:15` |
| V3 | recorded transition identities are hashes | `export PATH=/opt/homebrew/bin:$PATH; rg -n 'TransitionFn\|TransitionRef' host/store/journal.go host/replay/replay.go host/daemon/handlers.go \| head -30` | journal fields at :34-35; replay source ref at :68-69; daemon DTO fields at :48/:57 and :94/:103 |
| V4 | exported World surface | `export PATH=/opt/homebrew/bin:$PATH; rg -n '^export type\|^export func' world/types.ail world/transitions.ail` | 8 exported types; `plan`, `verify`, `applyRevision`, `commit` |
| V5 | proposal declaration fields | `export PATH=/opt/homebrew/bin:$PATH; rg -n 'transitionFn:\|expectedEffects:\|requiredCaps:' world/types.ail world/transitions.ail` | type fields at `types.ail:45,48,49`; constructor values at `transitions.ail:32,35,36` |
| V6 | generic object/head persistence exists | `export PATH=/opt/homebrew/bin:$PATH; rg -n 'epoch_registry_heads\|func \(s \*Store\) (SetRegistryHead\|GetRegistryHead)\|INSERT OR IGNORE INTO objects' host/store/schema.sql host/store/store.go` | head table `schema.sql:50`; object inserts :450/:841; set/get at :606/:624 |
| V7 | replay pins source/interpreter rather than registry nomination | `export PATH=/opt/homebrew/bin:$PATH; rg -n 'authoritative\|TransitionFn addresses\|ENTRY.s own interpreter\|registry candidate' host/replay/replay.go \| head -30` | authoritative pair documented at lines 1,17,63,68,145,178 |
| V8 | AILANG hardcoded totals | `export PATH=/opt/homebrew/bin:$PATH; rg -n 'REQUIRED_VERIFIED\|EXACT_TOTAL_VERIFIED\|REQUIRED_TESTS\|EXACT_TOTAL_TESTS' scripts/verify_ail.sh` | manifest :189; verified total 4 at :238; named tests :261; test total 14 at :268 |
| V9 | current AILANG gate result; supplied module count refuted | `export PATH=/opt/homebrew/bin:$PATH; AILANG_BIN=/tmp/ailang-v0300/ailang ./scripts/verify_ail.sh` | rc=0; `4/4 ... across 11 module(s)`; all 14 named tests; 9/9 package steps; PASS |
| V10 | CI contains exactly two jobs and invokes both scripts | `export PATH=/opt/homebrew/bin:$PATH; rg -n '^  [a-zA-Z0-9_-]+:$\|scripts/verify_(ail\|go)' .github/workflows/ci.yml` | jobs `ailang-verify` at :17 and `go-verify` at :98; scripts at :96/:165 |
| V11 | broker code is landed but implemented doc header is stale | `export PATH=/opt/homebrew/bin:$PATH; rg -n '^\*\*Status\*\*\|^# w-effect-broker' design_docs/implemented/w-effect-broker-m3.md \| head -5` plus V2 and focused test below | file is under `implemented/` but line 3 says `Status: Planned`; broker types/functions exist |
| V12 | focused broker/store/replay baselines | `export PATH=/opt/homebrew/bin:$PATH; GOTOOLCHAIN=go1.25.6 AILANG_BIN=/tmp/ailang-v0300/ailang go test ./host/broker -run 'Test(ReplayReturnsRecordedBytesWithoutDispatch\|ReplayGapNeverFallsBackToLive)$' -count=1; ... go test ./host/store -run 'TestRegistryHeadRoundTrip$' -count=1; ... go test ./host/replay -run 'Test(FixtureEpisodeReplaysBitForBit\|EpochRegistryCandidateCannotRedirect)$' -count=1` | all rc=0; broker 0.381s, store 0.298s, replay 2.197s |
| V13 | full Go gate sandbox limitation | `export PATH=/opt/homebrew/bin:$PATH; GOTOOLCHAIN=go1.25.6 AILANG_BIN=/tmp/ailang-v0300/ailang ./scripts/verify_go.sh` | toolchain control and `go build ./...` passed; `go test ./...` rc=1 with multiple `listen ... bind: operation not permitted`; non-socket packages including store/replay passed |
| V14 | new package and doc absent at baseline | `export PATH=/opt/homebrew/bin:$PATH; test ! -d host/transitionreg; printf ' rc=%s\n' "$?"; test ! -e design_docs/planned/w-transition-registry.md; printf ' rc=%s\n' "$?"` | both rc=0 before this document was created |
| V15 | DESIGN section location and protocol row | `export PATH=/opt/homebrew/bin:$PATH; rg -n '^### 11\.[12]\|Discovery.*what transitions\|Project the .*transition registry' design_docs/DESIGN.md` | protocol row at :146; AI UI heading at :332; discovery at :336; human UI §11.2 at :348 |

| V16 | **`rg` is not a binary on this rig or in CI — it is a harness-injected shell function**, so the drafted AC1–AC8 could not have run outside the measuring environment | `export PATH=/opt/homebrew/bin:$PATH; type rg; whence -p rg \|\| echo NO-BINARY; env -i PATH=/usr/bin:/bin:/opt/homebrew/bin sh -c 'command -v rg \|\| echo ABSENT'; grep -c 'rg ' .github/workflows/ci.yml scripts/*.sh` | `rg is a shell function from …/shell-snapshots/snapshot-zsh-….sh`; `NO-BINARY`; `ABSENT`; **0** occurrences in `ci.yml` and in all six `scripts/*.sh` |
| V17 | `grep -c '^Test'` is exact-equivalent to the drafted `rg -c '^Test'` on `go test -list` output, in **both** directions | `export PATH=/opt/homebrew/bin:$PATH; GOTOOLCHAIN=go1.25.6 go test ./host/store -list 'TestRegistryHeadRoundTrip$' \| grep -c '^Test'; GOTOOLCHAIN=go1.25.6 go test ./host/store -list 'TestNoSuchNameXYZ$' \| grep -c '^Test'` | known-positive **1**; known-negative **0** — so the replacement instrument is proven to see a present name and to refuse an absent one |
| V18 | the sandbox red the designer recorded in V13 is the documented loopback-bind denial, **not** a repo defect — controller re-ran the affected package OUTSIDE the codex sandbox | `export PATH=/opt/homebrew/bin:$PATH; GOTOOLCHAIN=go1.25.6 AILANG_BIN=/tmp/ailang-v0300/ailang go test ./host/daemon/` (main checkout, same commit) | `ok github.com/sunholo-data/ailang-world/host/daemon 1.734s`, **rc=0**; the sandbox log carries **12** `bind: operation not permitted` lines |
| V19 | controller re-derivation of V9 (the refuted module count), independently of the designer | `export PATH=/opt/homebrew/bin:$PATH; AILANG_BIN=/tmp/ailang-v0300/ailang ./scripts/verify_ail.sh` | rc=0; `✓ 4/4 required world/ identities verified across 11 module(s)`; `✓ all 14 required named tests pass`; `✓ world package gate PASSED: 9/9 steps`. The controller's supplied "4/9/14" was a transcription from a stale charter row — **the designer's refutation is confirmed** |
| V20 | Go `encoding/json` behavior on adversarial vectors, measured by `/tmp/tr_json_probe.go` with `Decoder.UseNumber` then `Marshal` | `export PATH=/opt/homebrew/bin:$PATH; go run /tmp/tr_json_probe.go` | duplicate `{"a":1,"a":2}` accepted as `{"a":2}`; `1`, `1.0`, `1e0` preserved as three different spellings; `"é水"` emitted as UTF-8; lone surrogate and invalid UTF-8 each accepted and replaced by `�`; large integer `9007199254740993123456789` preserved with `UseNumber`; non-object `[1,2]` accepted. Therefore `TR-CJSON-1` pre-validates duplicates/UTF-8/surrogates/root and performs its own decimal normalization instead of trusting standard marshal output. |
| V21 | revised count gates at base | `export PATH=/opt/homebrew/bin:$PATH;` run the exact AC1–AC7 commands above, printing each rc/count | AC1 rc=0 count=0; AC2 0/0; AC3 0/0; AC4 0/1 (known existing control); AC5 0/0; AC6 0/0; AC7 0/0. No result depends on the transition-registry directory. |
| V22 | revised aggregate gate at base | `export PATH=/opt/homebrew/bin:$PATH;` run the exact AC10 command above, then print rc/count | build PASS; store PASS; broker PASS; AC10 rc=0 count=0 because the new package is absent. |
| V23 | `grep` exists without the harness shell snapshot | `export PATH=/opt/homebrew/bin:$PATH; env -i PATH=/usr/bin:/bin:/opt/homebrew/bin sh -c 'command -v grep'` | `/usr/bin/grep`, rc=0. |
| V24 | v1 interface hash derived from the frozen ASCII preimage in Decision 2 | `export PATH=/opt/homebrew/bin:$PATH; printf '%s' 'world/transition-registry/v1\|TR-CJSON-1\|revision:{entries,interfaceHash,parent,revision,semanticID}\|descriptor:{access,declaredEffects,description,id,inputSchema,interpreter,outputSchema,semanticsEpoch,title,transitionFn}\|id:lower-ascii-1..128-segments-1..32\|schema:raw262144-canonical65536\|entries:1024\|revision:raw16777216-canonical8388608' \| shasum -a 256` | `743f39f470bf354ebab0ab196598b5ba72db80463d833325cb7672249d4734ac`, rc=0. |
| V25 | controller measurement of every production `Session.Invoke` call | `export PATH=/opt/homebrew/bin:$PATH; grep -rn '\.Invoke(' --include='*.go' host/ cmd/ \| grep -v _test.go` | exactly **3**, all in `host/broker/publish_op.go` at lines 135, 162, and 279; the current `_test.go` known-positive arm returns **90** |
| V26 | controller measurement of production session construction | `export PATH=/opt/homebrew/bin:$PATH; grep -rn 'NewSession(\|NewReplaySession(\|newSession(' --include='*.go' host/ cmd/ \| grep -v _test.go` | exported `broker.NewSession`: **0** production callers; exported `NewReplaySession`: **0** production callers; unexported `newSession`: exactly **3**, all in `host/broker/publish_op.go` at lines 132, 159, and 276, providing the firing positive control |
| V27 | current coordinator/daemon-to-broker dispatch path | controller call-graph synthesis of V25 and V26, tracing the only production construction and invocation sites | **none**: the daemon does not invoke the broker; exported `Session.Invoke` (`broker.go:126`) is reached in production only from inside `host/broker` |

No `.ail` source or syntax sketch is proposed, so S5 pinned-binary source validation is not
applicable. V9 nevertheless runs the pinned binary over every existing `.ail` module.

**V16–V19 and V25–V27 are controller rows (iteration 70), added after the designer run.** V16/V17 justify the
AC repair noted above; V18 and V19 are first-party re-derivations of a designer claim and of a
designer *refutation of the controller*, respectively — a sub-agent's finding is a claim in both
directions, and both survived re-measurement.

## Open Decision

1. **Bootstrap contents.** Recommended default: an empty revision 1 is legal and non-vacuously
   tested; actual package transitions enter through the later coordinator/package-install lane.
   P6.B fixtures publish synthetic descriptors through the writer, never hardcode production tools.
This does not affect the v1 wire format: the same frozen revision encoder represents zero or more
entries, and the empty revision's bytes and hash are golden-tested by TR.A. It only selects which
otherwise-valid descriptors an authorized production bootstrap later publishes.

## Related Documents

**Quorum note.** This narrow carve-out revision applies `gpt5-6-sol`'s option-1 binding fix and
`gemini-3-1-pro`'s non-blocking immutable-hash cache fix; Decisions 1–8 retain their direction.

- [DESIGN.md](../DESIGN.md) — §3.7, §5, §11.1/11.2, §14
- [coding-standards.md](../coding-standards.md) — S1–S7
- [w-mcp-projection.md](w-mcp-projection.md) — blocked consumer contract
- [w-effect-broker-m3.md](../implemented/w-effect-broker-m3.md) — landed broker design
- [world/types.ail](../../world/types.ail) — capability/proposal/transition shapes
- [host/registry/registry.go](../../host/registry/registry.go) — interpreter epoch registry, not reused
- [host/replay/replay.go](../../host/replay/replay.go) — authoritative hash-pinned replay
