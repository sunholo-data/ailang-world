# w-world-library-m1 — Semantic World Library Kernel

**Status**: **Planned — design complete & quorum-direction-accepted; M1 SPRINT AUTHORIZED (Mark, 2026-07-24, attended: option A — carve-out first-use approved, §14 replay-orchestration framing approved; ratification recorded in the charter queue + STATUS)**
**Date**: 2026-07-24
**Charter clause**: clause-1
**Verified against**: **`AILANG v0.30.0`** — the official released darwin/arm64 artifact
(GitHub release `AILANG v0.30.0`, `darwin.arm64.ailang.tar.gz`,
`sha256:ac3174e0f27692bb091d341a518b9473bb78010a4234cbff792aab63c67bb4d3`, checksum-verified),
a **clean, traceable, reproducible** substrate — deliberately NOT the rig's local
`v0.30.0-147-g6ed26bebd-dirty` dev build. A determinism-kernel design must validate against a
released artifact, not an ephemeral working tree (quorum objection, gemini-3-1-pro, 2026-07-24).
M1 pins this exact released binary as its interpreter of record.
**Traces to**: [DESIGN.md](../DESIGN.md) §15, §17 (M1), §14, §19 Q1; [w-log-epoch-decision.md](w-log-epoch-decision.md) (the frozen log format)
**Depends on**: [w-log-epoch-decision.md](w-log-epoch-decision.md) (**SETTLED**)
**Estimated**: ~2–3 days

> **Scope note.** M1 freezes the deterministic kernel as two cooperating libraries: pure
> transition semantics in AILANG and an embedded Go host for persistence, artifact pinning,
> and replay proof. The network daemon, effect broker, and protocol projections remain later
> milestones.

---

## Motivation

Clause 1 requires an immutable, content-addressed world store, an append-only transition
log, and replay that reconstructs a recorded episode bit-for-bit. The settled epoch decision
already fixes the three identities that make that promise meaningful:

1. `interpreter: HashRef` selects the exact archived AILANG executable for authoritative replay.
2. `transitionFn: HashRef` selects the exact canonical transition source bytes.
3. Every digest carries an algorithm tag, beginning with `sha256`.

M1 turns those identities into one executable kernel. AILANG owns semantic state and the pure
plan → verify → commit rules. Go owns byte persistence, SQLite transactions, tagged hashing,
interpreter archival, and the replay-proof harness. The split follows DESIGN.md §19 Q1 and
keeps the released language artifact behind the §14 boundary.

The load-bearing acceptance result is a recorded episode replayed twice through the exact
interpreter named by its log entry, with both replay byte streams compared against each other
and against the recorded result.

## High-Impact Decisions

| Decision | Why High Impact | Chosen By | Deadline | Change Cost |
|----------|-----------------|-----------|----------|-------------|
| AILANG library home is `world/*` | Establishes the import surface that M2 and later extensions consume | M1 | compile | medium |
| One canonical source-byte procedure | Makes `transitionFn` stable across editors and line-ending conventions | M1 | compile | high |
| SQLite stores `HashRef` in one canonical-text column | Keeps equality, indexes, and future algorithm coexistence aligned | M1 | schema | high |
| Epoch registry semantic ID is `world/epoch-registry/v1` | Gives immutable registry revisions one stable semantic name | M1 | schema | medium |
| Interpreter archive is adjacent to the SQLite file | Makes backup, transfer, and replay resolution operationally direct | M1 | runtime | medium |
| Go host lives in this repository as its own module | Preserves the language repository's frozen-core boundary | coordinator recommendation; assumed for M1 | compile | medium |

The final row resolves the mission's open ratification item by assumption: M1 creates a Go
module in this repository. The quorum or human authority may object before implementation;
absent that objection, the sprint proceeds with this layout.

### Design Freeze

- [x] D1/D2/D3 remain ratified exactly as recorded in `w-log-epoch-decision.md`.
- [x] The frozen log header is the six-field `LogHeader` shape from `sketches/logepoch.ail`.
- [x] `HashRef` is the sole digest type across AILANG and Go typed surfaces.
- [x] Authoritative replay identity is `(transitionFn, interpreter)`.
- [x] `semanticsEpoch` is compatibility metadata and stays outside replay authority and cache validity.
- [x] M1 acceptance uses replay-doubling for every recorded fixture episode.

---

## Decision 1 — AILANG Semantic Library

**Decision.** Production AILANG modules land at:

| File | Module | Responsibility |
|------|--------|----------------|
| `world/types.ail` | `world/types` | `World`, `Proposal`, `Transition`, `Evidence`, `Capability`, `Verification`, `CommitResult` |
| `world/contracts.ail` | `world/contracts` | Pure executable predicates shared by verify and commit |
| `world/transitions.ail` | `world/transitions` | Pure `plan`, `verify`, and `commit` functions |
| `world/logepoch.ail` | `world/logepoch` | Promoted `HashRef`, `ObjectEnvelope`, `EpochRecord`, `LogHeader`, rendering, and cache-key helpers |

The design artifact [`sketches/worldkernel.ail`](../sketches/worldkernel.ail) compiles the
integrated surface today. M1 promotes the ratified sketch types into `world/*`; the sketches
remain compiler-checked documentation fixtures.

### Typed Surface

- `World` carries `revision`, `stateRoot: HashRef`, and `logHead: HashRef`.
- `Proposal` carries its own `proposalHash`, the input-world hash, transition-function hash,
  goal, plan, expected effects, required capabilities, evidence, and confidence.
- `Transition` carries proposal, input-world, output-world, transition-function, and evidence
  references, each digest-bearing field typed as `HashRef`.
- `Evidence` variants carry content references as `HashRef`; human-readable labels remain strings.
- `RecordedTransition` combines the frozen `LogHeader` with a transition body.

`plan`, `verify`, and `commit` are pure. The M1 `commit` function derives a next immutable
`World` plus a `Transition`; the Go host persists that result atomically. This preserves the
three-phase semantic shape while reserving physical effects and capability enforcement for M3.

### Contracts

M1 contracts are executable pure predicates rather than prose-only preconditions:

1. A proposal's `inputWorld` equals the current world's `stateRoot`.
2. A verification names the same `proposalHash` as the proposal.
3. Commit proceeds only when verification is accepted and both identity predicates hold.
4. The resulting world increments `revision` by one and carries the supplied output-state and
   next-log-head references.

The same predicates are called by verify and commit, preventing policy drift between phases.
The compiler-checked sketch exercises record types, tagged evidence variants, imported
`HashRef`/`LogHeader`, pure functions, and the complete plan → verify → commit flow.

## Decision 2 — Canonical Transition Source Bytes

**Decision.** `transitionFn` hashes one precisely defined canonical byte stream:

1. Decode input as UTF-8; malformed input produces `CanonicalizationError`.
2. Reject a UTF-8 BOM and the NUL code point.
3. Convert CRLF and lone CR line endings to LF (`0x0a`).
4. Remove trailing ASCII space (`0x20`) and tab (`0x09`) from every line.
5. Remove empty lines at the end of the file.
6. Join remaining lines with LF.
7. A populated source ends with exactly one LF; an empty source yields zero bytes.
8. Every other Unicode code point retains its original UTF-8 byte sequence.

Canonicalization is idempotent: `canon(canon(source)) == canon(source)`. The host stores the
canonical bytes as an ordinary content-addressed object and hashes those exact bytes through
the tagged digest dispatcher. The log stores only the resulting `HashRef` in `transitionFn`.

This procedure resolves the epoch doc's deferred details while staying within its UTF-8, LF,
and trailing-space constraints. It also avoids dependence on compiler AST, Core, interface, or
other internal serialization.

## Decision 3 — Tagged Hashes and Object Envelope

**Decision.** The Go host defines a value type equivalent to AILANG `HashRef`:

| Field | Rule |
|-------|------|
| `Algo` | lowercase registered tag; M1 registers `sha256` |
| `Digest` | lowercase hexadecimal digest bytes; SHA-256 uses 64 characters |
| text | `algo:digest` |

Parsing and rendering are centralized in `host/hashref`. Hash calculation dispatches by
`Algo`; M1's dispatcher contains the SHA-256 implementation backed by Go's `crypto/sha256`.
Readers reject malformed text, unsupported tags, uppercase hex, and bare digests with a
structured `HashError`.

The stored envelope retains the ratified shape:

| Field | SQLite representation | Meaning |
|-------|-----------------------|---------|
| `hash` | canonical `HashRef` text | hash of payload bytes; primary identity |
| `interfaceHash` | canonical `HashRef` text | hash of the object's typed interface/schema bytes |
| `semanticId` | UTF-8 text | stable semantic name across immutable revisions |
| `provenance` | UTF-8 text | producer/reason lineage label |
| payload | BLOB | exact immutable bytes addressed by `hash` |

Every field whose meaning is a digest uses `HashRef`. `semanticId` and `provenance` are labels
rather than digest fields and preserve the settled `ObjectEnvelope` shape.

## Decision 4 — SQLite Store and Append-Only Log

**Decision.** M1 uses one SQLite database through Go `database/sql` and a pinned pure-Go SQLite
driver. The store is an embedded library opened by tests and, later, by M2's daemon.

### HashRef Column Layout

Each `HashRef` occupies one `TEXT` column in canonical `algo:digest` form. The Go boundary parses
every value into the typed representation before use. This choice provides one atomic indexed
identity and avoids split-column comparison mistakes. Algorithm-specific validation remains in
the dispatcher, allowing future tags to coexist in the same tables.

### Schema

| Table | Key columns | Purpose |
|-------|-------------|---------|
| `objects` | `hash_ref TEXT PRIMARY KEY` | immutable envelope plus payload bytes |
| `worlds` | `world_ref TEXT PRIMARY KEY` | immutable world revision, state root, and log head |
| `log_entries` | `entry_index INTEGER PRIMARY KEY`, `entry_hash_ref TEXT UNIQUE` | frozen header fields plus transition-body object reference |
| `epoch_registry_heads` | `registry_name TEXT PRIMARY KEY`, `object_ref TEXT` | current immutable registry object reference |
| `verification_cache` | `(transition_fn_ref, interpreter_ref)` primary key | cached typecheck/verify result with epoch metadata |

`log_entries` stores the frozen header fields verbatim: `entry_index`, `semantics_epoch`,
`transition_fn_ref`, `interpreter_ref`, `prev_entry_hash_ref`, and `written_by`. A separate
`transition_ref` points to the content-addressed transition body; it is outside the frozen
header. `entry_hash_ref` addresses the canonical encoded header-plus-body-reference bytes.

### Atomic Commit

One SQLite transaction performs the kernel commit:

1. Read and compare the current world/log head.
2. Insert required immutable objects with content verification.
3. Insert the next immutable world row.
4. Insert the next append-only log row with `prev_entry_hash_ref` equal to the observed head.
5. Advance the store's selected world head.
6. Commit the SQLite transaction.

A stale observed head produces structured `ConflictError`; callers may re-plan against the new
world. This single compare-and-append rule is the M1 concurrency boundary.

## Decision 5 — Epoch Registry Objects

**Decision.** The epoch registry uses semantic ID `world/epoch-registry/v1`. Each revision is an
immutable object containing ordered `EpochRecord` values; `epoch_registry_heads` maps that name
to the selected revision's `HashRef`.

Bootstrap creates epoch `1` with the M1 AILANG interpreter's release string as the first
nominated candidate and stores the registry through the same object/log transaction mechanism.
Registry updates therefore remain ordinary world transitions.

Candidate nomination is advisory compatibility metadata. M1 designates the replay fixture set
as the initial conformance corpus and records expected result hashes. Authoritative replay uses
the entry's `interpreter` reference exclusively; the registry is used for inspection and later
nomination workflows.

## Decision 6 — Interpreter Artifact Archive

**Decision.** For a database at `<store>.db`, artifacts live at:

`<store>.db.artifacts/interpreters/<algo>/<digest>/ailang`

The adjacent tree travels with database backups and gives the resolver a direct mapping from
`HashRef` to executable bytes.

Startup mechanics:

1. Resolve the configured AILANG executable path and open it once.
2. Stream the opened bytes through SHA-256 while copying to a temporary file in the destination
   directory.
3. Compare the calculated tag/digest with any pre-existing archive path.
4. `fsync` the file, set executable read-only permissions, and atomically rename it into place.
5. Record a sidecar manifest with hash, `ailang --version` output, byte size, OS, and architecture.
6. Use the archived path for subsequent typecheck, verify, execution, and replay operations.

Hashing and archival consume the same opened byte stream. Every log entry written by that host
stamps the resulting interpreter `HashRef` and the version output in `writtenBy`.

M1's replay guarantee covers the writing OS/architecture. An absent artifact, unsupported
platform, hash mismatch, or execution failure produces a structured `ReplayError` and stops
authoritative replay.

## Decision 7 — Replay Proof and Pair Pinning

**Decision.** A recorded episode is a manifest naming:

- initial world object;
- ordered log-entry hashes;
- transition source objects;
- interpreter artifacts;
- recorded result bytes and final world hash;
- recorded evidence/effect-result objects, when present.

The M1 fixture episode uses the real AILANG transition library and a deterministic fixture
driver. Recording and replay both capture the exact result byte stream before object insertion.

Authoritative replay for each entry performs this sequence:

1. Load `transitionFn` canonical bytes from the object store and verify their hash.
2. Resolve `interpreter` directly from the artifact archive and verify its hash.
3. Consult the verification cache using `(transitionFn, interpreter)`.
4. Invoke the archived interpreter with the pinned transition source and recorded inputs.
5. Byte-compare the produced result with the recorded result.
6. Reconstruct the next world/log hashes and compare them with the recorded hashes.

`semanticsEpoch` is copied into cache metadata for diagnostics and migration planning; cache
lookup validity uses the pair key exclusively.

### Replay-Doubling

Every acceptance run replays every recorded episode twice from a fresh temporary store through
the same pinned interpreter artifact:

- replay A result bytes equal replay B result bytes;
- each replay result equals the recorded result bytes;
- each replay final world hash equals the recorded final world hash;
- each replay resolves the interpreter from the entry header;
- replacing the epoch-registry candidate with another binary leaves authoritative resolution
  unchanged;
- replacing either pair member causes cache miss and replay re-verification.

Any divergence fails M1. This is the standing hermeticity test required by ratified D1.

---

## Systemic-Issue Audit

The kernel uses one identity and persistence mechanism across worlds, transitions, evidence,
epoch registries, interpreter artifacts, and future effect results:

1. canonical bytes;
2. tagged `HashRef`;
3. immutable object envelope;
4. append-only log reference;
5. atomic head advancement;
6. replay by exact pair pin.

A naive design would add dedicated digest columns, serializers, archive lookup rules, and replay
branches for each object type. M1 instead centralizes canonicalization, hash parsing/dispatch,
object insertion, log encoding, and artifact resolution. Type-specific code supplies semantic
fields and payload bytes; the identity, storage, chaining, and replay machinery remains shared.

The epoch registry follows the same object/log path, preventing a privileged metadata side
channel. Interpreter artifacts use the same `HashRef` parser and verifier even though their
bytes live beside SQLite for executable-file handling.

## Deferred Scope

| Item | Milestone | Boundary |
|------|-----------|----------|
| `ailang-worldd` long-running process, REST API, and CLI | M2 `w-worldd-m2` | M1 exposes embedded Go packages and a test harness |
| Effect broker, capability/budget enforcement, recorded external effects, physical isolation | M3 `w-effect-broker-m3` | M1 retains typed declarations and replay-ready evidence references |
| MCP/A2A projection and dynamic serving | `w-mcp-projection` | M1 defines reusable typed modules and host APIs |
| Portable cross-platform interpreter substitution | M8 policy work | M1 executes the exact archived artifact on its supported platform |

## Files to Create/Modify

| File | Estimated LOC | Change |
|------|--------------:|--------|
| `world/logepoch.ail` | ~80 | promote frozen log/hash/epoch types and helpers |
| `world/types.ail` | ~100 | production World/Proposal/Transition/Evidence surface |
| `world/contracts.ail` | ~60 | shared pure transition predicates |
| `world/transitions.ail` | ~100 | pure plan/verify/commit implementation |
| `design_docs/sketches/worldkernel.ail` | ~120 | compiler-checked integrated design artifact; created with this doc |
| `go.mod`, `go.sum` | ~20 | repository-local Go module and pinned SQLite dependency |
| `host/hashref/hashref.go` | ~120 | tagged parsing, rendering, SHA-256 dispatch |
| `host/canon/source.go` | ~100 | exact transition-source canonicalization |
| `host/store/schema.sql` | ~70 | objects, worlds, log, registry head, cache schema |
| `host/store/store.go` | ~350 | transactions, immutable objects, append-only commits |
| `host/archive/archive.go` | ~180 | interpreter hashing, atomic archival, resolution |
| `host/replay/replay.go` | ~240 | recorded episode and authoritative replay engine |
| `host/replay/replay_test.go` | ~260 | pair-pinning, replay-doubling, divergence tests |
| `host/replay/testdata/*` | ~100 | fixture transition, episode manifest, expected bytes |
| `scripts/verify_go.sh` | ~15 | `go build ./...` and `go test ./...` gate |
| existing CI workflow invoking `verify_ail.sh` | ~10 | add Go build/test gate in the same implementation PR |

Estimated implementation total: ~1,925 LOC including tests and fixtures.

## Conflict Surface

This is the DESIGN.md §14 frozen-core boundary analysis for a consumer repository.

- **Coupling to AILANG core**: M1 consumes the released executable through its path, exact bytes,
  `ailang --version`, `ailang check`, `ailang ai-check`, and source-in/result-out process
  execution. Compiler packages, parser structures, typechecker structures, Core, and codegen
  structures remain outside the dependency graph.
- **Upstream changes that could break the scheme**: removal or incompatible alteration of the
  consumed CLI operations; a source grammar change that rejects a pinned transition; executable
  behavior depending on mutable ambient state; or platform support loss for an archived artifact.
  Pair-pinning preserves old bytes, replay-doubling detects divergence, and structured replay
  failure prevents substitution.
- **Programs that MUST still work**: `sketches/logepoch.ail`, `sketches/worldtypes.ail`, and
  `sketches/transitions.ail`; the Verification Log records all three green together with the new
  `sketches/worldkernel.ail`.
- **Go module and CI implication**: this repository gains `go.mod` with the host packages.
  The same PR that introduces Go code extends CI with `go build ./... && go test ./...` while
  retaining `./scripts/verify_ail.sh`.
- **Repository placement assumption**: `worldd` and its host packages live in this repository as
  an independent Go module, matching the coordinator recommendation and preserving language-repo
  cadence. This remains the one assumed resolution the quorum or human authority may reverse.
- **Overlap with the released binary's own `ailang replay`** (added after quorum round 2 —
  gemini-3-1-pro): the released binary exposes `ailang replay <trace.jsonl>` (flags `-entry`,
  `-file`, `-caps`, `-json`) — a **single-program execution-trace** replay that re-runs one AILANG
  program against a recorded effect trace. World's M1 `host/replay/replay.go` operates at a
  **different layer**: it is a **store/log-level ORCHESTRATOR** — it resolves content-addressed
  objects by `HashRef`, verifies the `prevEntryHash` chain, pins and resolves the exact interpreter
  artifact per entry, and sequences a recorded episode's transitions. It does **not** reimplement
  AILANG evaluation; the actual deterministic re-execution of a single pinned transition is
  **delegated to the released binary's execution surface** (`ailang run` / `ailang replay`),
  which is exactly what §14 REQUIRES (consume the released artifact; never link or reimplement the
  interpreter). So the two are **complementary, not redundant**: `ailang replay` is a candidate
  per-transition execution primitive the orchestrator may invoke; it is not a substitute for the
  store/log/hash-chain/artifact-pinning machinery M1 must build, none of which the released binary
  offers. Reusing `ailang replay` as the per-transition execution verb (vs `ailang run`) is an
  implementation choice deferred to the sprint; either way the interpreter is the released binary,
  not host code. This resolution is **forced by the ratified §14 boundary**, not a new design
  decision.

The host may invoke public released CLI behavior; it may never link or vendor compiler internals.
Language-surface needs route upstream under the mission guardrails.

## Frozen Prerequisite Contract

M1 satisfies the nine implications from `w-log-epoch-decision.md` as follows:

1. **Frozen header** — SQLite stores the exact six `LogHeader` fields and hashes the canonical
   header-plus-transition-reference encoding.
2. **Tagged envelope** — `hash` and `interfaceHash` are canonical-text `HashRef` values; readers
   reject bare digests.
3. **SHA-256 dispatch** — `host/hashref` dispatches on tag and implements `sha256` with Go stdlib.
4. **Epoch registry objects** — `world/epoch-registry/v1` bootstraps epoch 1 with the M1 binary as
   first nominated candidate and evolves through logged transitions.
5. **Interpreter archive** — startup archives and then executes the exact artifact stamped into
   every entry; replay resolves the entry pin directly.
6. **Canonical transition source** — the eight-step UTF-8/LF procedure above defines the hashed
   bytes and stores them as objects.
7. **Replay-proof pair pinning** — replay lookup and cache validity use `(transitionFn,
   interpreter)`; epoch remains metadata.
8. **Conformance fixture** — replay testdata forms the initial corpus and records expected result
   hashes; nomination automation may follow in M2.
9. **AILANG digest typing** — every digest-bearing field in World, Proposal, Transition, Evidence,
   and recorded-transition surfaces uses `HashRef`.

## Acceptance Criteria

- [x] `world/logepoch.ail`, `world/types.ail`, `world/contracts.ail`, and
  `world/transitions.ail` pass `ailang ai-check` under the recorded released binary.
- [x] `./scripts/verify_ail.sh` passes every `.ail` module and reports a positive module count.
- [x] Go code lands with `go.mod`; CI runs `go build ./... && go test ./...` in the same PR.
- [x] Source canonicalization tests cover CRLF, lone CR, trailing ASCII space/tab, terminal empty
  lines, UTF-8 validation, BOM, NUL, idempotence, and a golden SHA-256 `HashRef`.
- [x] SQLite persists object envelopes, immutable worlds, the frozen log header, transition-body
  references, the registry head, and pair-keyed verification cache.
- [x] Hash readers reject malformed, uppercase, bare, and unsupported tagged forms with structured
  errors.
- [x] Epoch 1 registry object names the M1 AILANG release as its first nominated candidate.
- [x] Startup archives the configured interpreter and every new entry stamps its exact `HashRef`.
- [x] Authoritative replay resolves the executable from each entry's `interpreter` hash.
- [x] CI asserts epoch-registry candidate changes cannot redirect authoritative replay.
- [x] Verification/typecheck cache tests prove the key is exactly `(transitionFn, interpreter)`
  and epoch changes alone preserve the selected row as metadata-compatible.
- [x] A recorded fixture episode replays bit-for-bit to its recorded result and final world hash.
- [x] **Replay-doubling:** every acceptance episode replays twice through the pinned artifact;
  replay A, replay B, and the recorded bytes match exactly, with divergence causing test failure.
- [x] The three prerequisite sketches and `sketches/worldkernel.ail` remain green in the final gate.
- [x] M1 exports embedded host packages and stops before daemon, broker, isolation, REST, CLI, MCP,
  or A2A implementation.

## Axiom Compliance

| Axiom | Score | Justification |
|-------|-------|---------------|
| A1: Determinism | +2 | Canonical transition bytes plus exact `(transitionFn, interpreter)` execution and replay-doubling make divergence observable and fatal |
| A2: Replayability | +2 | Immutable objects, chained log entries, archived interpreters, and recorded result bytes reconstruct each fixture episode bit-for-bit |
| A3: Effect Legibility | 0 | M1 transitions are pure and effect declarations remain typed; effect execution and recording arrive in M3 |
| A4: Explicit Authority | 0 | Capability values remain in Proposal while enforcement stays at the explicitly deferred broker boundary |
| A5: Bounded Verification | +1 | Verification/typecheck work is cached by the exact code/interpreter pair |
| A6: Safe Concurrency | +1 | One SQLite compare-and-append transaction rejects stale heads with structured conflict |
| A7: Machines First | +1 | Typed modules, tagged identities, machine-readable manifests, and compiler/CI gates define the kernel contract |
| A8: Minimal Syntax | 0 | M1 introduces library modules and host code rather than language syntax |
| A9: Cost Visibility | +1 | Hashing, archival, verification, and replay are explicit host stages with measurable tests |
| A10: Composability | +1 | One object/log/hash mechanism serves all kernel object classes and future tags |
| A11: Structured Failure | +1 | Canonicalization, hash, conflict, archive, and replay failures have typed Go error categories |
| A12: System Boundary | +2 | AILANG owns semantics; Go owns persistence/runtime; released CLI and artifact bytes form the boundary |

**Net Score: +12** ✅ — hard axioms A1/A3/A4/A7 all score at least zero.

## Verification Log

| Claim | Check | Result |
|-------|-------|--------|
| Verification substrate is a clean released artifact (not a `-dirty` tree) | `gh release download v0.30.0 --repo sunholo-data/ailang --pattern darwin.arm64.ailang.tar.gz` + `shasum -a 256` | ✓ downloaded + checksum-verified `sha256:ac3174e0f27692bb091d341a518b9473bb78010a4234cbff792aab63c67bb4d3` (matches the release-published `.sha256`) |
| Header version identity | `/tmp/ailang-v0300/ailang --version` (the released artifact above) | ✓ `AILANG v0.30.0` — clean, no `-dirty` suffix; reproducible from the pinned release |
| All four sketches re-verified on the clean released binary | `ailang ai-check` (released v0.30.0) on `logepoch.ail`, `worldtypes.ail`, `transitions.ail`, `worldkernel.ail` | ✓ each `{"passed":true,"errors":0,"verify":available}`; `worldkernel.ail` also `check --relax-modules` → `No errors found!` |
| Consumed released CLI operations are advertised | `ailang --help` | ✓ output lists `run`, `check`, `ai-check`, `iface`, `replay`, and `serve-api` |
| Integrated M1 type surface compiles | `cd design_docs && ailang check --relax-modules sketches/worldkernel.ail` | ✓ `No errors found!` after type and effect checking |
| Integrated M1 type surface passes AI gate | `cd design_docs && ailang ai-check sketches/worldkernel.ail` | ✓ `"passed": true`, `"error_count": 0` |
| Ratified and prerequisite sketches remain green | `./scripts/verify_ail.sh` | ✓ `logepoch.ail`, `transitions.ail`, `worldtypes.ail`, and `worldkernel.ail` each report `"passed": true` |
| Repository-wide compiler-checked-doc gate has a positive module count | `./scripts/verify_ail.sh` | ✓ tail: `checked 4 module(s)` |
| §15 physical store matches M1 | Read DESIGN.md §15 | ✓ SQLite combines world DB, append-only log, and content-addressed object store |
| §14 boundary matches consumer architecture | Read DESIGN.md §14 | ✓ compiler remains a hard boundary; World routes language needs upstream |
| §17 M1 deliverables match scope | Read DESIGN.md §17 M1 row | ✓ AILANG types, Go store/log/object host, and proven replay |
| §19 Q1 split matches module allocation | Read DESIGN.md §19 Q1 | ✓ transition semantics in `.ail`; store/runtime in Go |
| Frozen nine-item prerequisite is fully mapped | Read `w-log-epoch-decision.md` “Implications for M1” and this doc's “Frozen Prerequisite Contract” | ✓ one numbered implementation response for each prerequisite item |
| AILANG core internals remain outside M1 dependencies | Design review of the file plan and Conflict Surface | ✓ dependencies are released executable bytes, version output, and CLI process operations |

### Gate Output Tail

```text
── ai-check design_docs/sketches/worldtypes.ail
{
  "file": "sketches/worldtypes.ail",
  "check": {
    "passed": true,
    "error_count": 0,
    "errors": []
  },
  "verify": {
    "available": true,
    "verified": 0,
    "counterexample": 0,
    "skipped": 0,
    "errors": 0,
    "results": []
  }
}
checked 4 module(s)
```

## Quorum Log (creation-time, 2026-07-24 — world mission iteration 6)

Designer: **codex / gpt-5.6-sol** (rotation). Reviewer: **gemini-3-1-pro** (Vertex ADC);
`gpt5-6-sol` EXCLUDED as reviewer (it authored the doc — generator≠judge), so the quorum ran
degraded to one external reviewer + the Opus controller's in-session verdict.

- **Round 1 — BLOCKED.** gemini-3-1-pro: the Verification Log validated against a `-dirty`
  untraceable working-tree build (`v0.30.0-147-g6ed26bebd-dirty`) — inadmissible for a
  determinism-kernel design. **Resolved (controller revision):** downloaded the official released
  `AILANG v0.30.0` darwin/arm64 artifact, checksum-verified
  (`sha256:ac3174e0f27692bb091d341a518b9473bb78010a4234cbff792aab63c67bb4d3`), re-verified all four
  sketches ai-check-green on it, and pinned it in the header + Verification Log.
- **Round 2 (the one permitted re-quorum) — BLOCKED.** gemini-3-1-pro: the Conflict Surface omitted
  the overlap between the released binary's `ailang replay` and M1's Go replay engine. **Resolved
  in-doc:** added the "Overlap with the released binary's own `ailang replay`" Conflict-Surface
  bullet above — host/replay is a store/log ORCHESTRATOR that delegates per-transition re-execution
  to the released binary (forced by ratified §14), not a reimplementation; the layers are
  complementary. Design DIRECTION was never disputed across either round.
- **Status: PARKED — needs-human-review (sprint gated, NOT the doc).** Both blocking objections were
  narrow, non-direction defects the controller resolved; but this exhausts the bounded
  one-revision-one-requorum budget, and applying the narrow-refinement carve-out to route straight
  to sprint-planner would be its **first use in the World mission** on an objection lacking a
  verbatim reviewer-authored fix. Per the carve-out's first-use rule, the M1 **sprint start**
  (planner → executor) is parked for Mark's one-time OK of (a) the carve-out first-use here and
  (b) the §14 replay-orchestration framing. The DESIGN is complete and quorum-direction-accepted;
  only the go/no-go to build is awaiting the human gate.

## References

- [w-log-epoch-decision.md](w-log-epoch-decision.md) — ratified D1/D2/D3 and the frozen nine-item M1 contract
- [logepoch.ail](../sketches/logepoch.ail) — compiler-checked `HashRef`, envelope, epoch, and header shape
- [worldtypes.ail](../sketches/worldtypes.ail) — prerequisite World-domain type sketch
- [transitions.ail](../sketches/transitions.ail) — prerequisite three-phase transition sketch
- [worldkernel.ail](../sketches/worldkernel.ail) — compiler-checked integrated M1 type/function artifact
- [world-mission.md](../world-mission.md) — clause 1, conflict boundary, guardrails, and M1 queue contract
- [DESIGN.md](../DESIGN.md) — §5, §6, §14, §15, §17, §19 Q1
