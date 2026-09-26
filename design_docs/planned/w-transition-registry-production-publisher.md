# w-transition-registry-production-publisher — A Production Path That Publishes the Transition Registry (row 107)

**Status**: Planned (design doc complete; prototype executed end-to-end in the designer worktree —
see "Prototype" below and `.design/proto.patch`). All claims measured at `dev 0aa53e6e41ac2354d3993833a0335e00eb1806e6` in the detached worktree `.design-wt-iter195`; every "the codebase does X" claim has a Verification Log row with the command and its observed output. **No `.ail` file is touched by this row** (measured: `git status --short` lists only `.go` files — V20).
**Item**: queue row 107 of `design_docs/world-mission.md` (clause-6), regroom position **1** (REGROOMED 2026-09-26: "Without it the A2A card lists zero skills and clause 6's 'the transition registry is served' has no content").
**Clause**: clause-6 — *"the transition registry is served over MCP (capability-filtered per session) and an A2A agent card is published"*. Row 40 landed the SERVING half (card, capability filter, absent-head-zero-skills); this row lands the CONTENT half: a production path that populates the registry.
**Estimate**: ~1d (matches the row's estimate; prototype production LOC = 145 + 206 + 13 modified-line insertions, i.e. three milestones of ≤ ~150 code LOC each).
**Prototype files changed** (all in `.design/proto.patch`): NEW `host/transitionreg/publish.go` (145 code LOC), NEW `cmd/world-publish/transitions.go` (206 code LOC), MODIFIED `cmd/world-publish/main.go` (+12: flags/verb/usage), MODIFIED `host/store/context_roots_test.go` (+1 pin), plus three NEW test files (`host/transitionreg/publish_test.go`, `cmd/world-publish/transitions_test.go`, `host/daemon/registry_publisher_test.go`).

## Problem

`transitionreg.StoreReader.Publish` and `BuildNext` have **zero** non-test callers (F1), so every
production daemon's registry head is absent and its A2A card at `GET /.well-known/agent.json`
lists **zero skills** — honest (row 40 tested that), but clause 6's "the transition registry is
served" then serves an empty catalogue forever. The only code that ever populates a head is a
TEST helper (`seedTransitionRegistry` in `host/daemon/daemon_test.go`, V12) that hand-encodes a
revision and calls `store.CompareAndSetRegistryHead` directly — bypassing `BuildNext`/`Publish`
entirely. There is no operator-facing, honesty-checked path that publishes a revision through the
landed CAS machinery, and no daemon-level test proving the card then lists published content.

## Findings (F-table)

Controller facts re-verified at `0aa53e6`; every row below names the command and observed output
(full transcripts in the Verification Log).

| # | Fact | Verification (command → observed) |
|---|---|---|
| **F1** | `Publish`/`BuildNext` have **0** non-test callers outside `host/transitionreg`; **26** test references to `TransitionRegistryV1` | `grep -rnE '\.(Publish\|BuildNext)\(' --include='*.go' . \| grep -v _test.go \| grep -v '^./host/transitionreg/'` → **EMPTY (exit 1)**; control in the same breath `grep -rn TransitionRegistryV1 --include='*_test.go' . \| wc -l` → **26**. Controller's F1 numbers CONFIRMED exactly. |
| **F2** | `BuildNext` (transitionreg.go:164) calls `current.Validate()` first. `Validate` accepts a synthetic `Revision{SemanticID: SemanticIDV1, InterfaceHash: InterfaceHashV1, Revision: 0}` with no entries, **BUT** `BuildNext` sets `next.Parent = SumSHA256(EncodeRevision(current))` — a NON-ZERO parent — while `Publish` with zero `expectedHead` (transitionreg.go:221) **requires `next.Parent.IsZero()`**. There is **no genesis constructor**: the first revision MUST be constructed directly as `Revision 1, Parent zero`. | Executed as test `TestGenesisCannotUseBuildNextWithZeroExpectedHead`: `BuildNext(genesis, changes)` → `next.Parent.IsZero() == false`; `Publish(zero, next)` → error **"parent is not captured head (absent)"**. Mutation MUT-7 proves the direct construction is load-bearing (2 tests red under it). |
| **F3** | Descriptor fields exactly as stated: `ID, TransitionFn, Interpreter (hashref.HashRef), SemanticsEpoch, InputSchema, OutputSchema ([]byte canonical), Access (EffectRequirement), DeclaredEffects, Title, Description`; `Descriptor.Validate()` at transitionreg.go:266 | `read host/transitionreg/transitionreg.go` (struct at top; Validate at :266). Schema canonicality is enforced by `Validate` (byte-equal re-canonicalisation) and again by `EncodeRevision → Validate` inside `Publish`. |
| **F4** | The only production reader: `host/daemon/daemon.go:532 Reader: transitionreg.NewReader(d.store)`; `AgentCard` at `host/projection/projection.go:218` filters by `allowedDescriptors` (one snapshot + one capability snapshot, broker.Allows) and an absent head is a legitimate zero-skills 200 | `grep -n` both lines → `:532`, `:218` (controller said ~212; the `AgentCard` doc block begins at :212, the func at :218 — re-measured). |
| **F5** | `cmd/world-publish` verbs packet/approve/publish/reconcile publish PACKAGES through `broker.EffectRegistryPublish` and do **NOT** import `host/transitionreg` (**0** matches); `cmd/ailang-worldd/cli.go` is a pure HTTP **client** (`object get` → `GET /v1/objects/<ref>`; `commit` → `POST /v1/commit`) — it opens no store | `grep -c transitionreg cmd/world-publish/main.go cmd/ailang-worldd/cli.go` → **0 / 0**; `sed -n '155,253p' cmd/ailang-worldd/cli.go` (all verbs are `client.get`/`client.do` HTTP calls). |
| **F6** | `world/transitions.ail` exports exactly 4: `plan` (:23), `verify` (:45), `applyRevision` (:57), `commit` (:76); 95 lines | `grep -n '^export func' world/transitions.ail` → 4 rows as listed. |
| **F7** | P9 says the REGISTRY item adds no registration surface and anticipates "a later package-install/coordinator path may publish a validated next revision after its normal propose → verify → commit authorization; direct projection code receives only a reader" | `sed -n '49,85p' design_docs/planned/w-transition-registry.md` — P9 verbatim. **P9 constrains the registry item, not this row**: row 107's own text offers "a CLI verb" as the production path. The prototype keeps P9's reader-only invariant: the daemon still constructs only `transitionreg.NewReader`. |
| **F8** (new) | `store.Open` acquires an exclusive single-writer lock; a CLI cannot open the store while the daemon serves it | `sed -n '238,266p' host/store/store.go` — `acquireWriterLock(canonical)` before SQLite opens; `OpenReadOnly` "succeeds while another process holds writer authority". Publishing therefore happens **with the daemon down**, then the daemon's next `ReadSnapshot` sees the head (no startup cache — registry freeze: "every ReadSnapshot reads the current head first"). |
| **F9** (new) | The world-publish publish path is gated TODAY by: one-shot `--approval-ref` (minted by the ATTENDED `approve` verb via `broker.MintAttendedApproval`, scoped `Effect: EffectRegistryPublish`), `--credential-file`, `requireAttendedOperator` (CI-env refusal + controlling-TTY probe + typed phrase), and `--live` | `sed -n 'func runPublish'` block + `cmd/world-publish/fences.go:182` (refuseAutomationEnvironment), `:226` (requireAttendedOperator), `host/broker/approve.go:398` (`Effect: EffectRegistryPublish`), `:562` (scope check). The `approve` verb calls the SAME `requireAttendedOperator` before minting. |
| **F10** (new) | Honest descriptor inputs exist and are measurable: interpreter pin = `archive.New(dbPath).Archive(ailangBin)` — the SAME call the daemon runs at startup (daemon.go:494-505, root `<store>.artifacts`, `probeVersion` runs `<bin> --version`); transition source = `canon.Source(bytes)` + `store.PutObject`; schema canonicalisation = `canonicalSchema` (codec.go:296, unexported — prototype exports `CanonicalSchema`); epoch value = the interpreter-epoch registry's `firstEpoch = 1` (host/registry/registry.go:34) | All read in place; the prototype uses each (see Design Q2). |
| **F11** (new) | The daemon test harness is fully in-process: `newHandlerDaemon` (temp store via `daemon.New`), `requestRecorderAuth` (httptest recorder, **no sockets**), `mintSessionGrants` (mints a session with given effect grants). The existing card tests hand-seed a head via `seedTransitionRegistry` — the test-only copy this row replaces | `sed -n '22,60p' host/daemon/handlers_test.go`; `sed -n '1003,1033p' host/daemon/daemon_test.go`. |
| **F12** (new) | Two inventory gates FIRE on this row and must move with it (S8-class coupling): (a) `cmd/world-publish`'s FROZEN flag surface `flagNames` with exact-set AC24(b); (b) `host/store`'s `contextRootPins` production-context-roots inventory | (a) `sed -n` flagNames block — "ADDING a flag reds even if it is never passed"; (b) observed: `go test ./host/store` FAILED with `production context roots changed: got … cmd/world-publish/transitions.go|runTransitions|Background:1 …` until the pin was updated (V15). |

## Design

### Q1 — WHERE does publishing live?

**Decision: (a) an operator CLI verb — `world-publish transitions` — over a host-level publisher
core (`transitionreg.PublishSet`); NOT (b) the package pipeline, NOT a REST route.**

Measured comparison:

| Option | Authority gate today | Honest descriptor inputs | LOC | Verdict |
|---|---|---|---|---|
| **(a) operator CLI verb** (on `cmd/world-publish`, the only store-opening attended operator binary — F5) | `requireAttendedOperator` — the SAME human-in-the-loop fence `approve` uses (F9): CI-env tripwire, controlling-TTY probe, typed phrase. Prototype adds it measured (MUT-5 killed by the fence arms). | `--ailang-bin` runs the daemon's own archival call (F10); `transitionFnFile` canonicalises via `canon.Source` and stores an object; schemas via the codec's canonicaliser. | 145 (core) + 206 (verb shell) | **CHOSEN** |
| (b) a step of the world-publish PACKAGE pipeline | The package publish is a NETWORK write to a remote registry origin gated by a one-shot `EffectRegistryPublish` approval + credential + `--live` (F9). Coupling LOCAL transition visibility to a REMOTE package publish is the wrong dependency direction: a daemon operator with no package to publish could never populate their store. Measured: `publish` requires `--registry-origin` (a URL) and refuses without a credential file. | A package ready-packet carries PACKAGE-level digests (`contentHash`, `tarballSHA256`) — no per-transition ID/schema/capability fields exist today (F6: `world/transitions.ail` exports 4 functions, no descriptor manifest). Filling descriptors from a package needs a package-manifest format that does not exist. | Would require inventing a package-borne descriptor manifest + remote gating for a local write | REJECTED for this row; retained as the LATER evolution (P9's "package-install path") that calls the SAME `PublishSet` — the seam is designed for it |
| (c) daemon REST route (e.g. a `worldd transitions` verb) | Would need session-middleware + a broker decision per write — the machinery of row 106, which does not exist. P9 (F7): the registry item adds no REST route; a write endpoint any local process can reach is exactly the ambient-authority shape clause 3 forbids. `cmd/ailang-worldd` is a pure HTTP client (F5) — a verb there would have nothing to call. | — | — | REJECTED: opens a network write surface |
| (d) startup registration from `world/*.ail` exports | None — and it is the shape P3/the a2a doc line 197 explicitly rejects ("registry, not a startup copy of `.ail` exports"); `.ail` exports cannot carry interpreter hashes or capability declarations. | — | — | REJECTED (row text itself flags it) |

**Does a local operator CLI open an ambient-authority path (clause 3)?** No, by measurement of
what the write CAN grant: a registry entry grants **nothing** to any session — visibility is
still `broker.Allows`-filtered per session (F4; the prototype's daemon test negative arm
proves it), and future execution is confined to `transitionreg.Bind` → broker session (P4/P7).
Clause 3 is about AGENTS reaching the outside world; a local operator CLI is outside the agent
trust boundary — the same authority as running the daemon, deleting the DB, or running
`world-publish approve` (which also opens the store directly). What DOES matter, and is
enforced: (i) the verb is NOT drivable by an unattended loop — `requireAttendedOperator` (MUT-5);
(ii) it stays OFF the daemon's network surface — no REST route exists (F7); (iii) single-writer
(F8) means it cannot run under the daemon's feet.

### Q2 — Descriptor honesty: what must the publisher VERIFY before publishing?

Measured inventory of existing verification helpers the publisher can lean on:

- **TransitionFn exists in the store**: `store.GetObject` (existence by ref). `PublishSet.verifySources`
  refuses with the NEW typed `TransitionSourceAbsentError` before any write. Killed by MUT-2
  (two tests). A descriptor whose source object is absent can never be resolved by execution or
  replay, so the card must never list it — the prototype's `TestPublisherRefusalKeepsCardEmpty`
  pins this at daemon level.
- **Interpreter is an archived, manifest-backed interpreter**: interpreters are NOT store objects
  (they live under `<store>.artifacts/interpreters/`), so this is verified in the verb:
  `--ailang-bin` performs the daemon's own `archive.Archive` + `ReadManifest` (F10); `--interpreter-ref`
  must `Resolve` + `ReadManifest` or the verb refuses (MUT-6 killed). The manifest field is
  deliberately NOT in the descriptor file — the operator cannot type an unverified pin.
- **Schemas canonicalise**: `Descriptor.Validate` (F3) enforces byte-canonical schemas at
  `EncodeRevision` time inside `Publish`; the verb canonicalises inputs up front via the newly
  exported `transitionreg.CanonicalSchema` (the codec's own canonicalizer, unexported until this
  row — 3-line export, no second canonical form).
- **SemanticsEpoch**: the manifest must supply a non-zero epoch; the interpreter-epoch registry's
  first epoch is `1` (F10). The CLI refuses a zero epoch with that fact in the message. (A
  registry-derived epoch is a residual — see Residuals.)
- **ID grammar, effect requirement sanity, entry ordering**: `Descriptor.Validate` +
  `SortedDescriptors` inside `BuildNext`/genesis construction — landed code, unchanged.

### Q3 — CAS conflicts and idempotence

Measured behaviour (prototype, `host/transitionreg/publish_test.go`):

- **Two publishers race**: `Publish` CAS-fails the loser with the typed
  `*store.RegistryCASConflict` (detected via `store.IsRegistryCASConflict`). `PublishSet` performs
  **ONE bounded retry** — re-read the head, rebuild over the winner, re-publish — then surfaces
  a typed error if the second CAS also conflicts (no retry storm: `TestPublishSetSecondConflictSurfacesTypedError`
  asserts exactly 2 CAS calls). The retry MERGES rather than clobbers:
  `TestPublishSetCASRetryMergesWinner` proves the winner's entry survives in the merged revision
  (MUT-1 — skipping the captured-head CAS — is killed by the growth test).
- **Re-publishing the identical descriptor set is a NO-OP**: `PublishSet` compares canonical
  entry bytes (`sameEntries`, via `EncodeRevision` of synthetic revisions — the same bytes the
  store will hold) and returns `{Unchanged: true, Revision: N}` WITHOUT writing. A revision that
  adds no entries would be registry noise. Killed by MUT-3 at BOTH levels (unit + CLI stdout
  "UNCHANGED at revision 1").

### Q4 — The daemon-level acceptance test

`host/daemon/registry_publisher_test.go::TestPublishedTransitionsAppearOnAgentCard` (prototype,
green): starts a REAL daemon via the existing harness (`newHandlerDaemon` — `daemon.New` on a
temp store, F11), publishes through **the production path** — `transitionreg.NewReader(d.store).PublishSet(...)`,
the exact function the CLI verb calls after its fences — then:

- **Positive arm**: a session minted with the descriptor's Access capability
  (`mintSessionGrants(t, d, "ep", "world.apply")`) GETs `/.well-known/agent.json` → HTTP 200,
  skills = exactly the published transition's ID + title, verbatim.
- **Negative arm**: a session minted with an unrelated capability (`fs.read`) gets a legitimate
  **zero-skills** card at 200 — the capability filter applies to published entries, not just
  seeded fixtures.
- Companion negative path: `TestPublisherRefusalKeepsCardEmpty` — a descriptor with an absent
  source object is refused and the card stays empty.

No wall-clock assertions anywhere; the recorder harness binds no sockets (F11), and the full
`host/daemon` package is green in this sandbox (V13). The CLI's own fence/idempotence arms live
in `cmd/world-publish/transitions_test.go` (happy path publishes a REAL archived fake
interpreter — the house shell-script pattern from `host/archive/archive_test.go` — then asserts
the no-op republish).

### Q5 — Human policy decisions?

One genuine one, with a landing tranche that does not need it:

**Decision for Mark — A/B (one word):** the attended confirmation phrase for
`world-publish transitions` is the EXISTING shared phrase
`publish world/core@0.1.0 irreversibly` (fences.go:66) — accurate for the package publish,
odd for a local registry write.
- **A (prototype; recommended): keep the shared phrase.** It is deliberately awkward to type, it
  is one more fence shared with a vetted path, and a transitions-specific phrase adds new fence
  surface for zero authority gain (the gate's strength is the TTY probe + CI tripwire, not the
  phrase's semantics).
- **B: a transitions-specific phrase** (`publish transition registry irreversibly`), threaded as a
  per-verb phrase constant.

Everything in this row lands under A without Mark's input; B is a 5-line follow-up. See
"Decision for Mark".

## Milestones (each ≤ ~150 production code LOC; sizes are `grep -vE '^\s*(//|$)'` code lines)

| M | Deliverable | LOC | Acceptance criteria | Gate commands |
|---|---|---|---|---|
| **M1** | `host/transitionreg/publish.go`: `PublishSet` + `CurrentRevision` + `CanonicalSchema` export + typed errors (`TransitionSourceAbsentError`, `EmptyGenesisError`) + genesis direct-construction + bounded CAS retry + idempotent no-op | 145 | Unit suite green: genesis publishes rev 1 (F2's asymmetry pinned by `TestGenesisCannotUseBuildNextWithZeroExpectedHead`); growth → rev 2 via `BuildNext`; identical republish = `Unchanged`, head unmoved; absent source refused typed; CAS retry merges winner in exactly 2 CAS calls; second conflict surfaces `IsRegistryCASConflict` | `go vet ./host/transitionreg/ && go test ./host/transitionreg/ -count=1` |
| **M2** | `cmd/world-publish/transitions.go`: the `transitions` verb — refuse `--live`, `--store` fence, store-open-before-gate (WriterAlreadyActive surfaced as "stop the daemon first"), `requireAttendedOperator`, interpreter pin verification, manifest parse, `transitionFnFile` → `canon.Source` → `PutObject`, schema canonicalisation, `PublishSet` + typed result lines. Moves the FROZEN `flagNames` (+3) and `contextRootPins` (+1) inventories — SAME commit (F12) | 206 + 13 modified | CLI fence table green (live refused, CI STOP, wrong phrase STOP, missing pin exit 1, unarchived interpreter refused, missing manifest usage); happy path publishes revision 1 and republish prints `UNCHANGED at revision 1`; AC24(b) flag-set equality still green; `TestProductionContextRoots` green with the new pin | `go vet ./cmd/world-publish/ && go test ./cmd/world-publish/ -count=1 -run 'TestTransitions'` (plus full package for AC24) |
| **M3** | `host/daemon/registry_publisher_test.go`: the daemon-level acceptance test (Q4) — publish through `PublishSet` on the live daemon's store, card lists the skill for an allowed session, zero skills for a disallowed session, refusal keeps the card empty | test-only (0 production LOC) | `TestPublishedTransitionsAppearOnAgentCard` + `TestPublisherRefusalKeepsCardEmpty` green; no wall-clock assertions; no sockets | `go test ./host/daemon/ -count=1 -run 'TestPublish|TestCard'` |
| **M4** (follow-up, no new code) | Docs: QUICKSTART row for the verb (S7 — usage surfaces ship usage docs: `--help` text is in the flagset; the quickstart happy path INCLUDING payload construction is the manifest example from `transitions_test.go`) | ~15 doc LOC | Quickstart executed-verbatim: run the verb against a scratch store, restart daemon, `curl /.well-known/agent.json` lists the skill | manual + `go test ./... -count=1` full gate |

## Acceptance + mutation table

All seven mutations were **EXECUTED** on the prototype (mutate → targeted `go test` → revert);
every one was killed by the named assertion, none by the compiler (the one vet-only build break
in MUT-7's first cut was fixed so the TEST kills it, per S6).

| # | Mutation (the defect it models) | Killed by (assertion that fires) | Result |
|---|---|---|---|
| MUT-1 | Publisher skips the captured-head CAS (`expected` always zero → blind overwrite) | `TestPublishSetGenesisThenGrowth` (second publish errors; the CAS conflict the landed store returns for a stale expected head is not swallowed) | **KILLED** |
| MUT-2 | Publisher accepts a descriptor whose TransitionFn object is absent (`verifySources` call removed) | `TestPublishSetRefusesAbsentTransitionSource` (wants `*TransitionSourceAbsentError`) AND daemon-level `TestPublisherRefusalKeepsCardEmpty` (card must stay empty) | **KILLED ×2** |
| MUT-3 | Idempotent republish creates a new revision (unchanged early-return removed) | `TestPublishSetIdempotentRepublishIsNoOp` (revision must stay 1, head unmoved) AND CLI `TestTransitionsVerbHappyPathAndIdempotence` (stdout must say `UNCHANGED at revision 1`) | **KILLED ×2** |
| MUT-4 | Publisher writes Revision+2 (genesis `Revision: 2`) | `TestPublishSetGenesisThenGrowth` — `Publish`'s own "revision 2 is not expected 1" refusal fires through the test's success assertion | **KILLED** |
| MUT-5 | CLI skips the attended-operator gate | `TestTransitionsVerbFences/ci_environment_stops` and `/wrong_phrase_stops` (expect exit 3 STOP; got success) | **KILLED ×2** |
| MUT-6 | CLI accepts an unverified `--interpreter-ref` (Resolve/ReadManifest removed) | `TestTransitionsVerbFences/unverifiable_interpreter_ref_refused` (expect exit 1 naming "is not archived next to the store") | **KILLED** |
| MUT-7 | Genesis built via `BuildNext` over a synthetic revision 0 instead of the direct revision-1 construction | `TestPublishSetGenesisThenGrowth` + `TestPublishSetIdempotentRepublishIsNoOp` — `Publish` refuses the non-zero parent ("parent is not captured head (absent)"): the F2 asymmetry is load-bearing | **KILLED ×2** |

**Tally: 7/7 killed.**

The row's example mutations map as follows: "publisher skips CompareAndSet" = MUT-1;
"publisher accepts a descriptor whose TransitionFn object is absent" = MUT-2; "card ignores the
capability filter" — that mutant lives in LANDED row-40 code (`host/projection`), and this row's
daemon-level negative arm is its regression killer for published (non-seeded) entries:
`TestPublishedTransitionsAppearOnAgentCard`'s second session asserts zero skills, which fails if
the filter is bypassed for real published content; "publisher writes Revision+2" = MUT-4;
"idempotent republish creates a new revision" = MUT-3.

## Conflict Surface

| Surface | Collision? | Measured basis |
|---|---|---|
| **Row 106** (invocation coordinator, sibling, NOT in scope) | **No code collision; one convention hand-off.** Row 106 composes `transitionreg.Bind` + capsule + broker session + `store.Commit` and wires `/a2a/` dispatch. It will need to resolve `TransitionFn` source objects; this row introduces the object form `world/transition-source/v1` (journalObject convention: `InterfaceHash = SumSHA256(semanticID)`) in the CLI. Row 106's design should either adopt it or name its own — the prototype's constant lives in `cmd/world-publish/transitions.go`. | Row 106 text (world-mission.md:1868) names Bind/capsule/commit — none of which this row touches. `git status` shows no overlap with any capsule/replay/broker file. |
| **Row 108** (MCP dispatch, blocked upstream on ailang#885) | **No.** This row adds no MCP surface, no `serveapi` import, no `.ail` change (V20). | regroom table: row 108 is blocked on `serveapi/protocol` blobs; nothing here touches them. |
| **Row 93** (clause-4 floor) | **No.** The floor run re-enters after 106+107+108; a populated registry gives it content but changes no floor measurement path. | regroom table row 4. |
| **Clause 7 world-publish pipeline** | **Yes — deliberate, measured, and additive.** `cmd/world-publish/main.go` gains 12 lines (3 flags, 1 switch case, 1 usage line) and the FROZEN `flagNames` list moves (+3, the exact-set AC24(b) gate stays green); `host/store/context_roots_test.go` pin moves (+1, F12). The existing verbs' behaviour is untouched (full package suite green except the pre-existing AILANG_BIN-gated test). | V15, V16, V13. |
| **Row 40 / host/projection** | **No.** Zero projection files touched; the card's absent-head zero-skills behaviour stays as landed. | `git status --short` (V20). |

## Decision for Mark

**One word: A or B.** The `world-publish transitions` verb reuses the shared attended phrase
`publish world/core@0.1.0 irreversibly` (A — prototype default) or gets its own
(`publish transition registry irreversibly`, B — 5-line follow-up). **Recommendation: A** — the
gate's strength is the TTY probe + CI tripwire + deliberate friction, not the phrase's wording;
B adds new fence surface for no authority gain. **The row lands complete without this decision**
(prototype ships A); B is cosmetic follow-up.

## Residuals / named owners

1. **`world/transition-source/v1` object convention vs. row 106** — owner: row 106's designer.
   Adopt or rename before invocation lands, so capsule/replay and the publisher agree on the
   source-object form. (Prototype: `cmd/world-publish/transitions.go`.)
2. **Registry-derived `SemanticsEpoch`** — owner: this row's executor or a follow-up. The CLI
   requires a non-zero epoch and points at `firstEpoch = 1`; deriving it from the epoch-registry
   head (host/registry) is ~20 LOC and removes the last operator-typed field.
3. **P9's package-install path** — owner: a later row. `PublishSet` is the seam: a package
   pipeline step that extracts descriptors from a package manifest calls it under the package's
   own propose → verify → commit authorization. This row deliberately does not build that
   (Q1-b measurement).
4. **MUT-dry-run for the verb** (`--dry-run` prints the would-be revision without writing) —
   owner: executor's judgement (~15 LOC); the other verbs have it, `transitions` does not yet.
5. **QUICKSTART row (M4)** — owner: executor, same commit as M2 (S7).

## Verification Log

All commands run in `.design-wt-iter195` at HEAD `0aa53e6e41ac2354d3993833a0335e00eb1806e6`.

- **V1 (F1)** `grep -rnE '\.(Publish|BuildNext)\(' --include='*.go' . | grep -v _test.go | grep -v '^./host/transitionreg/'` → **empty, exit 1**. Control: `grep -rn TransitionRegistryV1 --include='*_test.go' . | wc -l` → **26**. `git rev-parse HEAD` → `0aa53e6e41ac2354d3993833a0335e00eb1806e6`.
- **V2 (F2)** `TestGenesisCannotUseBuildNextWithZeroExpectedHead` (executed, green): `BuildNext(Revision{SemanticID: SemanticIDV1, InterfaceHash: InterfaceHashV1, Revision: 0}, …)` → `next.Parent.IsZero() == false`; `Publish(zero, next)` → **"publish transition registry: parent is not captured head (absent)"**. Also read: `BuildNext` sets `Parent: hashref.SumSHA256(currentBytes)`; `Publish` zero-expected branch requires `next.Parent.IsZero()`.
- **V3 (F3)** `grep -n 'func (d Descriptor) Validate' host/transitionreg/transitionreg.go` → **:266**; struct fields read in the same file.
- **V4 (F4)** `grep -n 'Reader:.*transitionreg.NewReader' host/daemon/daemon.go` → **:532**; `grep -n 'func (h \*Handler) AgentCard' host/projection/projection.go` → **:218** (doc block from :212).
- **V5 (F5)** `grep -c transitionreg cmd/world-publish/main.go cmd/ailang-worldd/cli.go` → **0, 0**; `sed -n '155,253p' cmd/ailang-worldd/cli.go` → all verbs are HTTP client calls (`/v1/objects/`, `/v1/log`, `/v1/commit`).
- **V6 (F6)** `grep -n '^export func' world/transitions.ail` → plan :23, verify :45, applyRevision :57, commit :76; `wc -l` → 95.
- **V7 (F7/P9)** `sed -n '49,85p' design_docs/planned/w-transition-registry.md` → P9 verbatim; row 107 text offers "a CLI verb".
- **V8 (F8)** `sed -n '238,266p' host/store/store.go` → `acquireWriterLock(canonical)` in `Open`; `OpenReadOnly` doc: "succeeds while another process holds writer authority". `host/daemon/daemon.go` New(): "another process already holds writer authority for this database (single-writer is enforced, not conventional)".
- **V9 (F9)** `sed -n` `runPublish`/`runApprove` blocks: publish requires `requireApprovalRef` + `requireCredentialFile` under `--live`; both verbs call `requireAttendedOperator`; `host/broker/approve.go:398` `Effect: EffectRegistryPublish`, `:562` scope check; fences.go:66 phrase, :182 CI tripwire, :226 attended stack.
- **V10 (F10)** `sed -n '440,505p' host/daemon/daemon.go` → `archive.New(cfg.DBPath)` + `a.Archive(cfg.AilangBin)` + `ReadManifest` at startup; `archive.go:52` `.artifacts` suffix, `:184-206` New/Root/dirFor/pathFor, `probeVersion` runs `execPath --version` under a scrubbed env; `host/canon/source.go:51` `func Source`; `host/registry/registry.go:34` `const firstEpoch int64 = 1`; codec.go:296 `canonicalSchema` (unexported).
- **V11 (honesty helpers)** `grep -rn 'PutVerifyResult' host cmd --include='*.go' | grep -v _test` → only `host/replay/replay.go` (and store itself) — the verification cache has no production writer today, so this row does NOT gate on it (row 106 may).
- **V12 (test-only copy)** `sed -n '1003,1033p' host/daemon/daemon_test.go` → `seedTransitionRegistry` hand-encodes a revision and calls `store.CompareAndSetRegistryHead` directly — the pattern this row's production path replaces at daemon level.
- **V13 (gates)** `go vet ./...` → **exit 0**. `go test ./... -count=1` → **19 packages ok**; 4 FAIL, all on the pinned-binary precondition: cmd/world-publish `TestWorldCoreManifestMatchesTheCommittedGolden` ("AILANG_BIN unset: interfaceHashV2 needs the pinned binary; never skip"), host/broker `TestEpisodeLiveReplayThreeArmsAndEvidence` ("AILANG_BIN must name the pinned released interpreter"), host/pkgproj 3× ("AILANG_BIN unset: never skip"), host/verifygate 4× ("AILANG_BIN is unset … verify_go.sh already refuses to run without it"). **Environmental, pre-existing**: this rig's only `ailang` is `v0.43.1-8-ga2256b1c5-dirty` (measured: `ailang --version`), which those gates refuse by design; the failures name AILANG_BIN, not any file this row touched. **No socket-binding failures** (the daemon tests use httptest recorders); host/daemon, host/transitionreg, host/store, cmd/world-publish's new tests all green.
- **V14 (mutations)** MUT-1…MUT-7 executed as described in the mutation table; each run's observed `--- FAIL: <killer>` captured in-session (e.g. MUT-1 → `--- FAIL: TestPublishSetGenesisThenGrowth`; MUT-5 → `--- FAIL: TestTransitionsVerbFences/ci_environment_stops` and `/wrong_phrase_stops`; MUT-6 → `--- FAIL: …/unverifiable_interpreter_ref_refused`; MUT-7 → `--- FAIL: TestPublishSetGenesisThenGrowth` + `TestPublishSetIdempotentRepublishIsNoOp`). **7/7 killed.**
- **V15 (context-roots gate)** observed before the pin update: `TestProductionContextRoots` FAIL — "production context roots changed: got map[… cmd/world-publish/transitions.go|runTransitions|Background:1 …]"; after adding the pin: `go test ./host/store/ -count=1 -run TestProductionContextRoots` → **ok**.
- **V16 (frozen flag surface)** flagNames comment: "AC24(b) compares this list with the FlagSet as an exact set, so ADDING a flag reds even if it is never passed"; with +3 flags the full `go test ./cmd/world-publish/` suite (minus the AILANG_BIN-gated row) is green — the inventory moved WITH the change in the same diff.
- **V17 (writer-lock operation)** the verb's message for a daemon-held store is exercised by design: `store.IsWriterAlreadyActive` branch in `runTransitions` (code read at V8; the in-process tests use fresh temp stores so the branch's failure mode is the V8 lock, not a stub).
- **V18 (archive fake)** `host/archive/archive_test.go:47-54` — the house shell-script fake interpreter (`#!/bin/sh`, 0o755) is the pattern the CLI test reuses; first attempt with `/bin/echo` was SIGKILLed by the archived-copy probe (observed: "signal: killed") — macOS code-signing on copied binaries; the shell-script fake passes.
- **V19 (protocol messages)** happy path stdout (observed): "published transition registry revision 1 (head sha256:…); the daemon's A2A card will list these skills on its next read"; republish: "transition registry UNCHANGED at revision 1 (head sha256:…); nothing written".
- **V20 (no .ail touched)** `git status --short` → only `cmd/world-publish/main.go`, `host/store/context_roots_test.go` modified and 5 new `.go` files (3 test). **No `.ail` file, no `scripts/verify_ail.sh` constants move** — the S8 ledger is untouched.
- **V21 (patch)** `git diff > .design/proto.patch` (+ `git diff --no-index` for the 5 new files) → 7 file diffs, 2085 lines; regenerated cleanly to remove a duplicated-append artefact; `grep -c '^diff --git'` → **7**.

## Controller-fact audit (F1–F7 as supplied)

F1 **CONFIRMED** (empty + control 26). F2 **CONFIRMED with a sharper finding**: no genesis
constructor exists, and `BuildNext`-over-genesis is *structurally incompatible* with
`Publish`'s absent-head branch (V2, MUT-7) — the first revision must be constructed directly.
F3 **CONFIRMED** (fields + Validate :266). F4 **CONFIRMED** (daemon.go:532; AgentCard at :218,
doc block at :212). F5 **CONFIRMED** (verbs; 0 transitionreg imports in both cmds; worldd CLI
is an HTTP client). F6 **CONFIRMED** (4 exports, :23/:45/:57/:76). F7 **CONFIRMED** (P9 verbatim;
this row is P9's anticipated "later path" and adds no REST route — the verb is a local,
attended CLI). **None of F1–F7 was found FALSE.**