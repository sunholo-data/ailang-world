# w-transition-registry-production-publisher — A Production Path That Publishes the Transition Registry (row 107)

**Status**: Planned — **REVISION 2** (quorum round 1: three rejects; see "Quorum log". Every objection is answered by measurement plus prototype changes listed there, and all new claims carry Verification Log rows V22–V36). Design doc complete; prototype executed end-to-end in the designer worktree — see "Prototype" below and `.design/proto.patch`. All claims measured at `dev 0aa53e6e41ac2354d3993833a0335e00eb1806e6` in the detached worktree `.design-wt-iter195`; every "the codebase does X" claim has a Verification Log row with the command and its observed output. **No `.ail` file is touched by this row** (measured: `git status --short` lists only `.go` files — V20).
**Item**: queue row 107 of `design_docs/world-mission.md` (clause-6), regroom position **1** (REGROOMED 2026-09-26: "Without it the A2A card lists zero skills and clause 6's 'the transition registry is served' has no content").
**Clause**: clause-6 — *"the transition registry is served over MCP (capability-filtered per session) and an A2A agent card is published"*. Row 40 landed the SERVING half (card, capability filter, absent-head-zero-skills); this row lands the CONTENT half: a production path that populates the registry.
**Estimate**: ~1d (matches the row's estimate; prototype production LOC = 104 + 83 + 116 + 83 + 78 + 256 code lines, i.e. milestones of ≤ ~150 code LOC each — the round-1 three-file shape was SPLIT in revision 2 so the new checks keep every milestone inside the rule; measured V32).
**Prototype files — new (round 1, revised in place)**: `host/transitionreg/publish.go` (104), `cmd/world-publish/transitions.go` (256), `host/transitionreg/publish_test.go`, `cmd/world-publish/transitions_test.go`, `host/daemon/registry_publisher_test.go`.
**Prototype files — new (revision 2)**: `host/transitionreg/collide.go` (83 — genesis construction + the same-ID retry rule), `host/transitionreg/epoch.go` (116 — epoch derivation/validation), `host/transitionreg/verify.go` (83 — publisher constructor + source verification), `host/transitionreg/verify_test.go`, `host/archive/check.go` (78 — the bounded `<archived interpreter> check` subprocess).
**Prototype files — modified, landed**: `cmd/world-publish/main.go` (+12: flags/verb/usage, round 1), `host/store/context_roots_test.go` (+1 pin, round 1), `host/transitionreg/transitionreg.go` (+4: the StoreReader `arch` field — the publisher's archive handle), `host/broker/registry_publish_test.go` (+13: the AC10 subprocess-site driver for `host/archive/check.go`, inventory moved WITH the change — V31).

## Problem

`transitionreg.StoreReader.Publish` and `BuildNext` have **zero** non-test callers (F1), so every
production daemon's registry head is absent and its A2A card at `GET /.well-known/agent.json`
lists **zero skills** — honest (row 40 tested that), but clause 6's "the transition registry is
served" then serves an empty catalogue forever. The only code that ever populates a head is a
TEST helper (`seedTransitionRegistry` in `host/daemon/daemon_test.go`, V12) that hand-encodes a
revision and calls `store.CompareAndSetRegistryHead` directly — bypassing `BuildNext`/`Publish`
entirely. There is no operator-facing path that publishes a revision through the landed CAS
machinery, and no daemon-level test proving the card then lists published content.

**What "honesty-checked" means here, stated exactly (revision 2, objection B):** a published
descriptor is verified for (a) source **presence** (`verifySources`: the `TransitionFn` object
exists in the store — round 1); (b) source **loadability**: the canonical source bytes pass
`<archived interpreter> check <temp file>` — the archived interpreter binary the descriptor
itself pins, run bounded through `procbound` — BEFORE `PutObject` (CLI) and again inside
`PublishSet` before any head move, refusing with a typed `TransitionSourceInvalidError` on
non-zero; and (c) **epoch derivation** (objection A): `SemanticsEpoch` is derived from (or
validated against) the epoch registry under `world/epoch-registry/v1`, refusing
unknown/mismatched interpreter–epoch pairs, with **no default** — enforced inside `PublishSet`,
not just the CLI. `canon.Source` (host/canon/source.go:47-80, F13) enforces **normalisation only**
— it rejects a BOM, invalid UTF-8, and embedded NUL, splits lines (LF/CRLF/CR), trims trailing
space/tab, and performs **no parse and no type-check** — so canonicalisation alone proves
nothing about loadability; the interpreter check in (b) is what does. The loadability check's
honest scope is the hermetic capsule shape (F17): a source is checked STANDALONE in a scratch
root, so a module whose imports cannot resolve there (e.g. a copy of `world/transitions.ail` without its `world/*` siblings) is refused with the interpreter's own LDR001 output; per F18 the capsule execution engine stages only the entry source under `AILANG_FS_SANDBOX`, so the same source would fail identically at execution — the publish-time refusal is a necessary condition for executability, not a limitation. Publication verifies source presence and successful standalone checking under the pinned interpreter, plus advisory epoch nomination. It does not establish invocation compatibility, descriptor/schema correspondence, or successful execution.

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
| **F13** (r2, objection B / M-B1) | `canon.Source` normalises ONLY: rejects a leading BOM, invalid UTF-8, embedded NUL; splits lines on LF/CRLF/CR; trims trailing space/tab; drops trailing empty lines. **NO parse, NO type-check** — a syntactically garbage file canonicalises cleanly if it is UTF-8 without NUL | `sed -n '30,80p' host/canon/source.go` — the eight-step normaliser, no parser call anywhere in the file; control (positive proof it proves nothing): `canon.Source([]byte("not AILANG at all"))` returns those bytes canonicalised without error (V27) |
| **F14** (r2, objection A / M-A1) | The epoch registry exists as landed store state: head `world/epoch-registry/v1` = `store.EpochRegistryV1` (store.go:84); a revision decodes (registry.go:75 `Decode`) to `Registry{Epochs []EpochRecord}`, each `EpochRecord{Epoch int64, Candidates []string}` where Candidates are nominated interpreter RELEASE strings (Candidates[0] primary; the file comment calls them "advisory compatibility metadata"); `registry.Bootstrap(ctx, s, releaseString)` creates epoch 1 with `Candidates: [releaseString]`, idempotently | `sed -n '33,56p' + '101,173p' host/registry/registry.go`; `grep -n EpochRegistryV1 host/store/store.go` → :84; `grep -n 'registry.Bootstrap' host/daemon/daemon.go` → :512 (V22, V23) |
| **F15** (r2, objection A / M-A1 precision) | The daemon bootstraps the epoch registry with `release = releaseFromVersion(m.Version)` — the FIRST NON-BLANK LINE of the archived interpreter manifest's version, trimmed — NOT the verbatim multi-line probed string; with no `cfg.AilangBin`, `release = "unpinned"`. The manifest field holding the version is `Manifest.Version` (`json:"version"`), the verbatim `ailang --version` stdout captured at archival (archive.go:155-156, probeVersion :467). Any derivation MUST apply the same reduction or it false-mismatches | `sed -n '486,515p' + '613,624p' host/daemon/daemon.go` (`release = unpinnedRelease`; `release = releaseFromVersion(m.Version)`; the reduction loop); `ailang --version` observed multi-line: `AILANG v0.41.0 / Commit: 24ee108 / Full: …` (V24) |
| **F16** (r2, objection B / M-B2) | `host/procbound` is the repo's bounded-subprocess helper: `Admit()` reserves one of 8 process-wide slots (refuses `ErrCleanupBacklog` non-blocking), `Wait(wait, d, release)` bounds the child's reap by `d` and hands an unreaped child to a background waiter (`ErrCleanupIncomplete`); `host/archive.probeVersion` (archive.go:467) is the measured in-repo pattern for running the archived binary: `exec.CommandContext` + a wall-clock timeout + `childenv.Scrubbed(os.Environ())` + separated stdout/stderr buffers | `read host/procbound/procbound.go` (Admit :42, Wait :60, MaxOutstanding = 8); `sed -n '450,495p' host/archive/archive.go`; control: the archived copy of the pinned binary executes on this rig (V25: `cp $HOME/.pinned-ailang/ailang <artifacts>/copy-test && … --version` → exit 0, `AILANG v0.41.0`) |
| **F17** (r2, objection B) | The pinned interpreter's `check` subcommand is a real parse+type+effect check with usable exit codes, and the archived copy runs it: garbage bytes → **exit 1** (`Error: type error … undefined variable`); a module-declared standalone source in a scratch dir → **exit 0** (with a benign MOD010 temp-path warning, auto-relaxed); a module-LESS source → **exit 1** (MOD014 "no 'module' declaration"); `world/transitions.ail` copied WITHOUT its `world/*` imports into an empty dir → **exit 1** (`LDR001: module not found: world/logepoch`) | Observed transcripts in V26; run against `$HOME/.pinned-ailang/ailang` (v0.41.0) and its archived copy |
| **F18** (r2, objections A+B) | (a) `BuildNext` REPLACES by ID (`byID[change.ID] = …`, transitionreg.go:190) — so the round-1 CAS-retry rebuild would let the LOSER's bytes silently clobber the winner's fresh entry for a shared ID; (b) the capsule (row 106's execution engine) stages ONLY the entry source in its sandbox root (`AILANG_FS_SANDBOX=root`, no world lib copied) — so a source whose imports cannot resolve hermetically cannot be EXECUTED either; refusing it at publish time is the executability contract, not a limitation | `sed -n '164,196p' host/transitionreg/transitionreg.go`; `sed -n '180,215p' host/capsule/capsule.go` (staging writes only `host/capsule/main.ail`); contrast `host/replay/replay.go:319` which stages a world lib for its own recorded fixtures (V28) |

## Quorum log

**Round 1 — REJECTED 3/3 (all present).**

| Reviewer | Surface | Objection (abridged) | Answer in this revision |
|---|---|---|---|
| gpt6-astra **and** gemini-3-1-pro (same surface) | **A — epoch semantics**: `SemanticsEpoch` is operator-typed and only checked non-zero (Residual 2 defers the real check); the epoch must be DERIVED from / validated against the epoch registry, refusing unknown or mismatched interpreter–epoch pairs, no default-to-1, enforced in `PublishSet` (not just the CLI), with tests proving refusal leaves head and card unchanged | **Answered by measurement + prototype.** The derivation inputs are landed and measured (F14/F15): the epoch registry decodes to epoch→nominated-release-strings, and the interpreter's release string is the manifest's `Version` field reduced exactly as the daemon's bootstrap reduces it (first non-blank line, trimmed). `PublishSet` now REQUIRES the archive (typed refusal when constructed via bare `NewReader` — `NewPublisher(db, archive.New(dbPath))` is the only publishing constructor), and enforces per descriptor: read the epoch-registry head, decode, read the pinned interpreter's manifest, derive the release, and require `descriptor.SemanticsEpoch` ∈ {epochs whose Candidates name that release}. Zero matches → typed `*InterpreterEpochMismatchError` (unknown pair); epoch present in the registry but not nominating this release → the same typed error (mismatched pair); **no epoch registry head → typed `*EpochRegistryAbsentError`** (no default; the daemon owns epoch-1 bootstrapping). The CLI manifest's `semanticsEpoch` becomes OPTIONAL: omitted + exactly one matching epoch → derived and reported; omitted + several matching epochs → refusal naming them (the operator must state one); present → must be in the match set. Mutations MUT-9 (skip the check), MUT-10 (accept a mismatched epoch), MUT-11 (default-to-1) all KILLED; `TestPublishSetRefusesEpochNotNominatingTheInterpreter` + `TestEpochRefusalKeepsCardUnchanged` prove refusal leaves head and card unchanged. Residual 2 is RESOLVED (removed). |
| oc-kimi-k3 | **B — source validity**: the "honesty-checked" claim covers PRESENCE only; nothing verifies the transition source is loadable under the pinned interpreter; the doc never states what `canon.Source` enforces | **Answered by measurement + prototype (the strong arm, not the fallback).** F13 states `canon.Source`'s exact scope (normalisation only, no parse/type-check — V27 shows garbage canonicalises cleanly). The loadability check is REAL and infeasibility was NOT invoked: the verb already holds the archived interpreter; the NEW `host/archive/check.go` (`Archive.CheckSource`) runs `<archived interpreter> check <canonical source staged in a scratch root>` bounded via `procbound.Admit/Wait` plus the `probeVersion` wall-clock pattern (F16), and the publisher refuses with a typed `*TransitionSourceInvalidError` on any refusal, BEFORE `PutObject` (CLI, `transitionFnFile` path, via `transitionreg.EnsureSourceLoadable`) and again inside `PublishSet` (unskippable by non-CLI callers) before any head move. Measured on (i) a real module — `world/transitions.ail` copied alone → exit 1 `LDR001` (import-resolution refusal is the HONEST scope: the capsule engine stages no world lib, F18, so such a source is not executable either; pinned as an AILANG_BIN-gated CLI test) and (ii) garbage bytes → exit 1 (V26). MUT-8 (skip the source check) KILLED ×2. The Problem statement and Q2 now carry the precise "what is verified" wording; "honesty-checked" is no longer the vague claim kimi rejected. |
| oc-kimi-k3 (secondary) | **B — same-ID CAS-retry merge**: Q3 never said what happens when the CAS-retry winner and loser both set the same ID with different descriptor bytes; it must not silently drop, and the CLI output must name the ID | **Answered by measurement + prototype.** Measured: `BuildNext` replaces by ID, so the round-1 retry WOULD have silently clobbered (F18a) — kimi's instinct was right. Q3 now specifies exactly: identical bytes for a shared ID → idempotent merge (`Unchanged` at the winner's revision, entry intact); DIFFERENT bytes for a shared ID → typed `*SameIDConflictError{ID,…}` refusing the retry — the loser writes NOTHING, the head stays at the winner's revision, and both the error and the CLI output NAME THE ID. Pinned by `TestPublishSetSameIDConflictOnCASRetryRefuses` (different bytes → refusal, winner intact, ID in error) and `TestPublishSetCASRetrySameIDIdenticalBytesIsNoOp` (identical bytes → `Unchanged`); the CLI rendering is pinned by `TestPublishErrorLineNamesConflictingID`. MUT-12 (silent clobber) KILLED. |

**Round 2 — BLOCKED 2/1 (all present; gemini-3-1-pro PASS).** Surfaces: **A′ — epoch-check
strength** (gpt6-astra: release-name nomination is advisory, so the doc overclaimed an
interpreter–epoch binding) and **B′ — executability wording + unstated consequence** (oc-kimi-k3:
the Problem statement said the capsule "could resolve" import-bearing sources while F18 says it
stages only the entry source; and no landed `.ail` source can be published yet). Neither disputes
the direction and both carry concrete fixes, so the controller applied the **narrow-refinement
carve-out**: kimi's two fixes VERBATIM (Problem sentence; Residual 7, backed by the new V37);
astra's claim-weakening sentence VERBATIM (Q2) and its two-interpreter test adopted as
AC-EPOCH-TWIN in the form that pins the true behaviour. astra's further proposal — an authorized
HashRef↔epoch binding record — was **not** applied, because it would reverse ratified D1
(epoch = compatibility metadata, HashRef = authoritative pin); the refutation is in Q2. The controller then re-ran astra alone against that revision ($0.22): it
**no longer contested the epoch point** and rejected on a NEW overclaim, "guaranteed
capsule-executable" (introduced by kimi's verbatim Residual-7 text), with a verbatim
claim-weakening fix, applied in Q2, the Problem statement and Residual 7. Its request to prove
invocation through the production capsule path is Residual 8 (row 106 owns the invocation
contract, which does not exist yet). Objections localised onto claim STRENGTH across rounds 2–3
while gemini passed; no further round was run.

## Design

### Q1 — WHERE does publishing live?

**Decision: (a) an operator CLI verb — `world-publish transitions` — over a host-level publisher
core (`transitionreg.PublishSet`); NOT (b) the package pipeline, NOT a REST route.**

Measured comparison:

| Option | Authority gate today | Honest descriptor inputs | LOC | Verdict |
|---|---|---|---|---|
| **(a) operator CLI verb** (on `cmd/world-publish`, the only store-opening attended operator binary — F5) | `requireAttendedOperator` — the SAME human-in-the-loop fence `approve` uses (F9): CI-env tripwire, controlling-TTY probe, typed phrase. Prototype adds it measured (MUT-5 killed by the fence arms). | `--ailang-bin` runs the daemon's own archival call (F10); `transitionFnFile` canonicalises via `canon.Source`, is CHECKED loadable under the archived interpreter BEFORE it is stored (revision 2, Q2), and then stores an object; schemas via the codec's canonicaliser. | 104+83+116+83+78 (core) + 256 (verb shell) | **CHOSEN** |
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
- **SemanticsEpoch is DERIVED / validated against the epoch registry (revision 2, objection A)**:
  the manifest's `semanticsEpoch` is OPTIONAL — omitted, the verb derives it: read the epoch
  registry head (`world/epoch-registry/v1`) → decode → read the pinned interpreter's archived
  manifest → reduce `Manifest.Version` with the daemon's own rule (first non-blank line, trimmed,
  F15) → the epochs whose `Candidates` name that release. Exactly one match → that epoch, and the
  success output says so; several matches → refusal naming them (the operator states one);
  zero matches → refusal (no epoch nominates this interpreter release). **There is NO default —
  in particular no default-to-1.** The enforcement lives in `PublishSet` itself: a publisher is
  constructed with `transitionreg.NewPublisher(db, archive.New(dbPath))`, and `PublishSet`
  refuses (typed `*PublisherArchiveRequiredError`) on a bare `NewReader` store — so a non-CLI
  caller cannot skip the epoch check any more than it can skip `verifySources`. `PublishSet`
  re-validates every descriptor's epoch against the same registry (typed
  `*InterpreterEpochMismatchError` / `*EpochRegistryAbsentError`), so a caller that hand-fills an
  epoch is still bound by the registry. Refusal leaves the head and the card unchanged — pinned
  at unit level (`TestPublishSetRefusesEpochNotNominatingTheInterpreter`) and at daemon level
  (`TestEpochRefusalKeepsCardUnchanged`).
  **What this check does and does not establish (quorum r2, gpt6-astra; narrow-refinement
  carve-out).** Candidates matching establishes advisory release nomination only; it does not establish compatibility of a pinned interpreter HashRef with an epoch. That is the
  strength the ratified log-epoch decision assigns it: **D1 (RATIFIED 2026-07-24, Mark,
  attended — A+B-metadata)** makes the exact-binary `interpreter` HashRef the *authoritative*
  pin and `semanticsEpoch` *compatibility metadata* (`design_docs/planned/w-log-epoch-decision.md`
  lines 26–30, 109, 117; charter line 943). The authoritative half of every published descriptor
  is therefore its `Interpreter` HashRef, which the publisher requires to be an archived,
  manifest-backed binary (above) and which execution and replay resolve exactly (P7). The
  reviewer's further proposal — an explicitly authorized HashRef↔epoch binding record — would
  promote the epoch from metadata to authority, i.e. reverse a ratified human ruling; this row
  does not do that, and if it is wanted it is a new decision for Mark, not a revision of this
  doc. The epoch check is kept because it is exactly as strong as the daemon's own bootstrap
  nomination and it stops an operator typing an epoch the registry does not nominate. The
  reviewer's test is adopted in the form that pins the TRUE behaviour: two archived
  interpreters with different hashes and an identical first version line are BOTH
  epoch-eligible under the same nomination, AND each published descriptor carries its own
  distinct `Interpreter` HashRef (so the card, invocation and replay never conflate them);
  refusal of a non-nominated release still leaves head and card unchanged (acceptance row
  AC-EPOCH-TWIN below).
- **Transition source is LOADABLE under the pinned interpreter (revision 2, objection B)**:
  `canon.Source` (F13) is normalisation only — it proves nothing about parse or types. The verb
  therefore runs the REAL check: the new `host/archive/check.go` (`Archive.CheckSource`) stages
  the canonical bytes in a scratch root and runs `<archived interpreter binary> check <file>`
  bounded via `procbound.Admit/Wait` plus a wall-clock context and the scrubbed child
  environment (the `probeVersion` pattern, F16), and the publisher refuses with a typed
  `*TransitionSourceInvalidError` (carrying the interpreter's captured output) on any refusal —
  BEFORE `PutObject` for the `transitionFnFile` path (the CLI calls
  `transitionreg.EnsureSourceLoadable`), and again inside `PublishSet` (which reads the stored
  payload) before any head move, so the `transitionFn` ref path and every non-CLI caller are
  covered too. Honest scope (F17/F18): the check is hermetic — a source whose imports cannot
  resolve in the scratch root is refused with the interpreter's own LDR001 output; per F18 the
  capsule execution engine (row 106's `Bind` → capsule) stages only the entry source under
  `AILANG_FS_SANDBOX`, so the same source would fail identically at execution. A refusal here therefore predicts a refusal at
  execution; an ACCEPTANCE predicts nothing stronger than the check itself. Publication verifies source presence and successful standalone checking under the pinned interpreter, plus advisory epoch nomination. It does not establish invocation compatibility, descriptor/schema correspondence, or successful execution. (quorum r2 astra re-run, verbatim; see Residual 8.) The check adds no new authority
  surface: `check` parses and type-checks, it never runs the module. Two S8-class inventories
  FIRE on this addition and move with it (V31): the AC10 subprocess-site census
  (`host/broker/registry_publish_test.go` gained the `host/archive/check.go` driver) and the
  `host/verifygate` subprocess-sink gate (the check's output sinks are the mandated
  concurrency-safe buffers, not bare `bytes.Buffer`).
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
- **Same ID on the CAS-retry merge (revision 2, kimi's secondary; exact semantics)**: before the
  retry rebuilds, `PublishSet` compares each change's canonical descriptor bytes with the
  WINNER's head entry for that ID. **Identical bytes → idempotent merge**: the shared ID's entry
  is already exactly what we were going to write, the retry converges to `Unchanged` at the
  winner's revision, nothing is dropped and nothing is rewritten
  (`TestPublishSetCASRetrySameIDIdenticalBytesIsNoOp`). **Different bytes → typed
  `*SameIDConflictError{ID, WinnerEpoch…}`**: the retry is REFUSED before any write — the loser
  writes nothing, the head stays at the winner's revision with the winner's bytes intact, and
  the error (and the CLI's output for it) NAMES THE ID: "both this publish and revision N's
  winner set `tools.echo` with different bytes — re-run with the winner's descriptor, or remove
  the entry". There is no silent drop and no silent clobber (`TestPublishSetSameIDConflictOnCASRetryRefuses`;
  MUT-12 kills the clobber). A same-ID REPLACE outside a race (no CAS conflict) is unaffected:
  it is the ordinary `Change` update semantics — the collision rule applies only to the retry,
  where two publishers both believed they were writing fresh bytes.
- **Re-publishing the identical descriptor set is a NO-OP**: `PublishSet` compares canonical
  entry bytes (`sameEntries`, via `EncodeRevision` of synthetic revisions — the same bytes the
  store will hold) and returns `{Unchanged: true, Revision: N}` WITHOUT writing. A revision that
  adds no entries would be registry noise. Killed by MUT-3 at BOTH levels (unit + CLI stdout
  "UNCHANGED at revision 1").

### Q4 — The daemon-level acceptance test

`host/daemon/registry_publisher_test.go::TestPublishedTransitionsAppearOnAgentCard` (prototype,
green): starts a REAL daemon via the existing harness shape (`daemon.New` on a temp store, F11) —
**revision 2: with `AilangBin` set to the house shell-script fake, so the daemon itself archives
the interpreter and bootstraps the epoch registry with its release** (the production shape:
daemon.go:494-512) — publishes through **the production path** —
`transitionreg.NewPublisher(d.store, archive.New(cfg.DBPath)).PublishSet(…)` pinning the
DAEMON'S OWN archived interpreter ref (`d.interpreterRef`) with `SemanticsEpoch` derived from
the bootstrapped registry — then:

- **Positive arm**: a session minted with the descriptor's Access capability
  (`mintSessionGrants(t, d, "ep", "world.apply")`) GETs `/.well-known/agent.json` → HTTP 200,
  skills = exactly the published transition's ID + title, verbatim.
- **Negative arm**: a session minted with an unrelated capability (`fs.read`) gets a legitimate
  **zero-skills** card at 200 — the capability filter applies to published entries, not just
  seeded fixtures.
- Companion negative path: `TestPublisherRefusalKeepsCardEmpty` — a descriptor with an absent
  source object is refused and the card stays empty.
- **Revision 2 (objection A): `TestEpochRefusalKeepsCardUnchanged`** — after one honest
  publish populates the card, a second descriptor with a registry-UNVERIFIED epoch
  (`SemanticsEpoch: 2` while the epoch registry holds only epoch 1 for that release) is refused
  with the typed `*InterpreterEpochMismatchError`, and the card still lists EXACTLY the first
  skill: refusal leaves both the head and the card unchanged.

No wall-clock assertions anywhere; the recorder harness binds no sockets (F11), and the full
`host/daemon` package is green in this sandbox (V13). The CLI's own fence/idempotence arms live
in `cmd/world-publish/transitions_test.go` (happy path publishes a REAL archived fake
interpreter — the house shell-script pattern from `host/archive/archive_test.go` — then asserts
the no-op republish; revision 2 adds arms for epoch derivation output, epoch-mismatch refusal,
garbage-source refusal before `PutObject`, and the AILANG_BIN-gated real-module LDR001 refusal).

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

## Milestones (each ≤ ~150 production code LOC; sizes are `grep -vE '^\s*(//|$)'` code lines, measured V32; the round-1 three-file shape was SPLIT in revision 2 so the new checks keep every milestone inside the rule)

| M | Deliverable | LOC | Acceptance criteria | Gate commands |
|---|---|---|---|---|
| **M1** | `host/transitionreg/publish.go`: `PublishSet` flow + `SetResult` + `CurrentRevision` + the publisher archive gate (typed `*PublisherArchiveRequiredError` on a bare `NewReader`) + the bounded CAS retry + the idempotent no-op | 104 | Bare-`NewReader` publish refused typed; genesis publishes rev 1; growth → rev 2 via `BuildNext`; identical republish = `Unchanged`, head unmoved; second conflict surfaces `IsRegistryCASConflict` in exactly 2 CAS calls | `go vet ./host/transitionreg/ && go test ./host/transitionreg/ -count=1` |
| **M2** (r1, carried) | `host/transitionreg/collide.go`: genesis direct-construction (`EmptyGenesisError`, the F2 asymmetry) **+ the same-ID retry rule** (`*SameIDConflictError`, `refuseConflictingSameIDs`, `descriptorBytes`) | 83 | `TestGenesisCannotUseBuildNextWithZeroExpectedHead` green; same-ID identical bytes → `Unchanged` at the winner's revision in exactly 1 CAS call; same-ID different bytes → the typed error naming the ID, winner's bytes intact; removal-only/empty genesis refused | `go test ./host/transitionreg/ -count=1 -run 'TestGenesis|TestPublishSetCASRetry|TestPublishSetSameID|TestPublishSetRefusesEmpty'` |
| **M3** (r2) | `host/transitionreg/epoch.go`: epoch derivation and validation — `EpochsForInterpreter` (registry head read + `registry.Decode` + manifest release reduction), `epochRegistry`, `verifyEpochs`, `epochIn`, `EpochRegistryAbsentError`, `InterpreterEpochMismatchError` | 116 | Derived epochs = the registry's nominations for the manifest-derived release; mismatched/unknown epoch refused typed with head unmoved; absent epoch registry refused typed (no default); epoch 1 refused when epoch 1 does not nominate the release (no default-to-1) | `go test ./host/transitionreg/ -count=1 -run 'TestEpochs|TestRelease|TestPublishSetRefusesEpoch|TestPublishSetRefusesDefault|TestPublishSetRefusesAbsent'` |
| **M4** (r2) | `host/transitionreg/verify.go`: `NewPublisher(db, archive)` + `CanonicalSchema` + `verifySources`/`TransitionSourceAbsentError` + `verifyLoadable` + `EnsureSourceLoadable`/`TransitionSourceInvalidError` (the typed wrap over `archive.CheckSource`) | 83 | Absent source refused typed; unloadable source refused typed BEFORE any head move with the entry ID named; the loadable positive control publishes | `go test ./host/transitionreg/ -count=1 -run 'TestPublishSetRefusesAbsent|TestPublishSetRefusesUnloadable|TestEnsure|TestCanonicalSchema|TestPublishSetRefusesWithoutArchive'` |
| **M5** (r2) | `host/archive/check.go`: the bounded check subprocess — `CheckSource` stages the canonical source in a scratch root and runs `<archived interpreter> check <file>` under `procbound.Admit/Wait` + a wall-clock context + `childenv.Scrubbed`, with concurrency-safe output sinks | 78 | A real archived interpreter accepts a standalone module-declared source and refuses an import-bearing one (LDR001) and garbage; the subprocess-sink gate (verifygate) is green; the AC10 subprocess-site inventory names this file and its driver scrubs the registry credential | `go vet ./host/archive/ && go test ./host/archive/ ./host/verifygate/ -count=1 -run 'TestCheck|TestSubprocessSinkGate' && go test ./host/broker/ -count=1 -run TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential` |
| **M6** | `cmd/world-publish/transitions.go` — verb shell: fences (refuse `--live`, `--store`, store-open-before-gate, `requireAttendedOperator`), `pinnedInterpreter` verification, `PublishSet` invocation + `publishSetErrorLine` typed rendering (**naming the same-ID conflict's ID**) | ~130 of 256 | CLI fence table green (live refused, CI STOP, wrong phrase STOP, missing pin exit 1, unarchived interpreter refused, missing manifest usage); same-ID conflict line names the ID | `go test ./cmd/world-publish/ -count=1 -run 'TestTransitionsVerbFences|TestPublishErrorLine'` |
| **M7** | `cmd/world-publish/transitions.go` — materialisation: `readManifest`, `buildChanges` **with epoch derivation (optional `semanticsEpoch`: derive when unique, refuse naming the epochs when ambiguous or unmatched, validate when present)**, `pinnedTransitionFn` (`canon.Source` → **`EnsureSourceLoadable` BEFORE `PutObject`** → object), schemas, requirement helpers | ~126 of 256 | Happy path publishes revision 1 with a DERIVED epoch and the output says so; explicit epoch must be in the registry's match set or exit 1; garbage source refused BEFORE the object is stored; AILANG_BIN-gated arm: an import-bearing module refused with LDR001 output, no object stored; republish prints `UNCHANGED at revision 1`; AC24(b) flag-set equality green (no new flags); `TestProductionContextRoots` green with the existing pin | `go test ./cmd/world-publish/ -count=1 -run 'TestTransitions'` (plus full package for AC24 and context roots) |
| **M8** | `host/daemon/registry_publisher_test.go`: the daemon-level acceptance tests (Q4) — publish through `NewPublisher(…).PublishSet` on the live daemon's store (with the daemon's own archived interpreter and bootstrapped epoch registry), card lists the skill for an allowed session, zero skills for a disallowed session, refusal keeps the card empty, **epoch-refusal keeps head and card unchanged** | test-only (0 production LOC) | `TestPublishedTransitionsAppearOnAgentCard` + `TestPublisherRefusalKeepsCardEmpty` + `TestEpochRefusalKeepsCardUnchanged` green; no wall-clock assertions; no sockets | `go test ./host/daemon/ -count=1 -run 'TestPublish|TestCard|TestEpochRefusal'` |
| **M9** (follow-up, no new code) | Docs: QUICKSTART row for the verb (S7 — usage surfaces ship usage docs: `--help` text is in the flagset; the quickstart happy path INCLUDING payload construction is the manifest example from `transitions_test.go`) | ~15 doc LOC | Quickstart executed-verbatim: run the verb against a scratch store, restart daemon, `curl /.well-known/agent.json` lists the skill | manual + `go test ./... -count=1` full gate |

## Acceptance + mutation table

All twelve mutations were **EXECUTED** on the revision-2 prototype (mutate → targeted `go test` →
revert); every one was killed by the named assertion, none by the compiler. MUT-1…MUT-7 are
unchanged from round 1 and were RE-RUN against the revised prototype — after the revision-2 file
split (publish/collide/epoch/verify) all twelve were run once more against the FINAL files so the
tally reflects the code the executor will receive (V29); two mutation cuts needed adjusting to stay
honest (MUT-2 removes BOTH presence arms because `verifyLoadable` double-covers `verifySources`;
MUT-7's replacing form also drops `collide.go`'s then-unused `hashref` import so the TEST kills it,
not the compiler — S6).

| # | Mutation (the defect it models) | Killed by (assertion that fires) | Result |
|---|---|---|---|
| MUT-1 | Publisher skips the captured-head CAS (`expected` always zero → blind overwrite) | `TestPublishSetGenesisThenGrowth` (second publish errors; the CAS conflict the landed store returns for a stale expected head is not swallowed) | **KILLED** |
| MUT-2 | Publisher accepts a descriptor whose TransitionFn object is absent (BOTH presence arms removed: the `verifySources` call AND `verifyLoadable`'s absent-object branch — the revision-2 core double-covers presence, so the faithful defect-mutation removes both) | `TestPublishSetRefusesAbsentTransitionSource` (wants `*TransitionSourceAbsentError`) AND daemon-level `TestPublisherRefusalKeepsCardEmpty` (card must stay empty) | **KILLED ×2** |
| MUT-3 | Idempotent republish creates a new revision (unchanged early-return removed) | `TestPublishSetIdempotentRepublishIsNoOp` (revision must stay 1, head unmoved) AND CLI `TestTransitionsVerbHappyPathAndIdempotence` (stdout must say `UNCHANGED at revision 1`) | **KILLED ×2** |
| MUT-4 | Publisher writes Revision+2 (genesis `Revision: 2`) | `TestPublishSetGenesisThenGrowth` — `Publish`'s own "revision 2 is not expected 1" refusal fires through the test's success assertion | **KILLED** |
| MUT-5 | CLI skips the attended-operator gate | `TestTransitionsVerbFences/ci_environment_stops` and `/wrong_phrase_stops` (expect exit 3 STOP; got success) | **KILLED ×2** |
| MUT-6 | CLI accepts an unverified `--interpreter-ref` (Resolve/ReadManifest removed) | `TestTransitionsVerbFences/unverifiable_interpreter_ref_refused` (expect exit 1 naming "is not archived next to the store") | **KILLED** |
| MUT-7 | Genesis built via `BuildNext` over a synthetic revision 0 instead of the direct revision-1 construction | `TestPublishSetGenesisThenGrowth` + `TestPublishSetIdempotentRepublishIsNoOp` — `Publish` refuses the non-zero parent ("parent is not captured head (absent)"): the F2 asymmetry is load-bearing | **KILLED ×2** |
| **MUT-8** (r2, objection B) | Publisher/CLI skip the source-loadability check (`EnsureSourceLoadable` call in `pinnedTransitionFn` and the `verifyLoadable` call in `PublishSet` removed — garbage bytes publish) | `TestPublishSetRefusesUnloadableTransitionSource` (wants `*TransitionSourceInvalidError`; head must stay absent) AND CLI `TestTransitionsVerbRefusesGarbageSourceBeforePutObject` (exit 1; the source object must NOT be in the store) | **KILLED ×2** |
| **MUT-9** (r2, objection A) | `PublishSet` skips the epoch check (`verifyEpochs` call removed) | `TestPublishSetRefusesEpochNotNominatingTheInterpreter` AND `TestPublishSetRefusesDefaultEpochOne` (both wanted typed refusals; the mismatched epoch publishes instead) AND daemon-level `TestEpochRefusalKeepsCardUnchanged` (the card must stay unchanged) | **KILLED ×3** |
| **MUT-10** (r2, objection A) | The derivation ignores the release: `EpochsForInterpreter` returns every epoch the registry holds, so any existing epoch is accepted for any interpreter | `TestPublishSetRefusesDefaultEpochOne` (epoch 1 EXISTS in the registry but does not nominate this release → must refuse) AND CLI `TestTransitionsVerbEpochArms` (the no-nomination arm expected exit 1, got a published revision) | **KILLED ×2** |
| **MUT-11** (r2, objection A) | Epoch defaults to 1 when unverifiable/absent (absent epoch-registry head → a synthetic epoch-1 registry; a non-nominating epoch → coerced to 1) | `TestPublishSetRefusesDefaultEpochOne` AND `TestPublishSetRefusesEpochNotNominatingTheInterpreter` AND `TestPublishSetRefusesAbsentEpochRegistry` (all three wanted refusals) AND CLI `TestTransitionsVerbEpochArms` (both refusal arms published instead) | **KILLED ×4** |
| **MUT-12** (r2, kimi secondary) | Same-ID CAS-retry silently clobbers: the `refuseConflictingSameIDs` call removed, so the loser's differing bytes replace the winner's | `TestPublishSetSameIDConflictOnCASRetryRefuses` — both of its arms fire: the typed `*SameIDConflictError` is absent (the publish SUCCEEDS) and the winner's `title-winner` bytes are gone | **KILLED** |
| **AC-EPOCH-TWIN** (r2, astra; carve-out, to be built by the executor — NOT yet run) | (a) two archived interpreters with different hashes and an identical first `--version` line: both derive the same nominated epoch, and each published descriptor carries its own distinct `Interpreter` HashRef; (b) a release the registry does not nominate is refused and head + card are unchanged | Mutation to kill: the publisher substitutes the nominated release's FIRST archived interpreter for the descriptor's own pin (the two descriptors then share one HashRef → (a) reds). Planner sizes it; the tally above stays 12/12 until it runs. |

**Tally: 12/12 killed.**

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
| **Row 106** (invocation coordinator, sibling, NOT in scope) | **No code collision; one convention hand-off.** Row 106 composes `transitionreg.Bind` + capsule + broker session + `store.Commit` and wires `/a2a/` dispatch. It will need to resolve `TransitionFn` source objects; this row introduces the object form `world/transition-source/v1` (journalObject convention: `InterfaceHash = SumSHA256(semanticID)`) in the CLI. Row 106's design should either adopt it or name its own — the prototype's constant lives in `cmd/world-publish/transitions.go`. **Revision 2 note for its designer (objection B):** the publish-time loadability check is hermetic (F17/F18) — it proves the source loads under the pinned archived interpreter in the capsule shape; if row 106's execution engine later stages a world library into the capsule root, the check's scope should widen with it (same helper, `transitionreg.EnsureSourceLoadable`). Published card content carries presence + hermetic-loadability + registry-epoch guarantees as of this row. | Row 106 text (world-mission.md:1868) names Bind/capsule/commit — none of which this row touches. `git status` shows no overlap with any capsule/replay/broker file. |
| **Row 108** (MCP dispatch, blocked upstream on ailang#885) | **No.** This row adds no MCP surface, no `serveapi` import, no `.ail` change (V20). | regroom table: row 108 is blocked on `serveapi/protocol` blobs; nothing here touches them. |
| **Row 93** (clause-4 floor) | **No.** The floor run re-enters after 106+107+108; a populated registry gives it content but changes no floor measurement path. | regroom table row 4. |
| **Clause 7 world-publish pipeline** | **Yes — deliberate, measured, and additive.** `cmd/world-publish/main.go` gains 12 lines (3 flags, 1 switch case, 1 usage line) and the FROZEN `flagNames` list moves (+3, the exact-set AC24(b) gate stays green); `host/store/context_roots_test.go` pin moves (+1, F12). The existing verbs' behaviour is untouched (full package suite green except the pre-existing AILANG_BIN-gated test). | V15, V16, V13. |
| **Row 40 / host/projection** | **No.** Zero projection files touched; the card's absent-head zero-skills behaviour stays as landed. | `git status --short` (V20). |
| **AC10 subprocess-site census + subprocess-sink gate (r2)** | **Yes — both fired during revision 2 and moved WITH the change (V31).** The new bounded check subprocess (`host/archive/check.go`) joins the AC10 census — `host/broker/registry_publish_test.go` gained its driver, and the child env is observed scrubbed like every other World-launched subprocess; `host/verifygate`'s sink gate forced the concurrency-safe buffer pattern. The check subprocess deliberately lives in `host/archive` (next to `probeVersion`): placing it in `host/transitionreg` measurably cycles (`registry_publish_test.go imports host/transitionreg … imports host/broker from bind.go`). | Observed FAIL → green transcripts in V31; both inventories green in the full run (V30). |

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
2. **~~Registry-derived `SemanticsEpoch`~~ — RESOLVED in revision 2** (objection A): the epoch
   is now derived from / validated against the epoch registry inside `PublishSet`
   (`EpochsForInterpreter` + `verifyEpochs`, typed refusals, no default); the CLI manifest field
   is optional and derived when unique. No residual remains on this point.
3. **P9's package-install path** — owner: a later row. `PublishSet` is the seam: a package
   pipeline step that extracts descriptors from a package manifest calls it under the package's
   own propose → verify → commit authorization (constructing `NewPublisher` so the epoch and
   loadability checks cannot be skipped). This row deliberately does not build that
   (Q1-b measurement).
4. **MUT-dry-run for the verb** (`--dry-run` prints the would-be revision without writing) —
   owner: executor's judgement (~15 LOC); the other verbs have it, `transitions` does not yet.
5. **QUICKSTART row (M7)** — owner: executor; landed in sprint milestone M7 (`docs/QUICKSTART.md` §6, after the verb exists — the sprint plan's renumbering), marked *attended — pending first verbatim run* until Mark executes it (S7). The executor ran only the non-TTY steps (V40); the TTY fences were NOT driven through a pty wrapper (MUT-5).
6. **Loadability-check scope vs. a future world-library-staging capsule (r2)** — owner: row 106's
   designer. The check is hermetic by measurement (F17/F18); if row 106's engine stages a world
   library for executed sources, the publish check should widen to the same shape
   (`EnsureSourceLoadable` is the seam).
7. **Operational consequence (quorum r2, oc-kimi-k3; verbatim).** Measured consequence: the repo's only landed transition module, `world/transitions.ail` (F6), imports `world/*` and is refused by this check (V26(b)); therefore after this row lands, the production path exists and everything it publishes passes standalone `check` under the pinned interpreter (not a guarantee of capsule executability — see Residual 8), but no existing `.ail` source can be published through it. First real card content requires either a new hermetically self-contained transition module or row 106 widening capsule staging plus `EnsureSourceLoadable` (Residual 6). Backed by V37.
8. **Invocation contract (quorum r2 astra re-run; claim weakened verbatim above).** Publication verifies source presence and successful standalone checking under the pinned interpreter, plus advisory epoch nomination. It does not establish invocation compatibility, descriptor/schema correspondence, or successful execution. Owner: **row 106's designer**. Row 106 defines the source-object and entry-point/calling convention; it must (a) reuse/replace this row's `world/transition-source/v1` convention (Residual 1), (b) add a verification row that publishes a self-contained transition through this path with a real pinned interpreter and invokes that exact stored source through the production capsule path, and (c) decide, with a negative fixture that passes standalone `check` but lacks the required callable interface, whether publication rejects it (widening `EnsureSourceLoadable`) or invocation reports a typed incompatibility. Until then the card can list a skill that checks but cannot be invoked — the same state every card skill is in today, since `/a2a/` refuses every invocation (row 106's problem statement).

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
- **V22 (F14 / controller M-A1 re-verify)** `sed -n '33,56p' host/registry/registry.go` → `type EpochRecord struct { Epoch int64 …; Candidates []string … }` with the comment "Candidates is the ordered list of nominated interpreter release strings. Candidates[0] is the first (primary) nomination. This is advisory metadata: authoritative replay uses the log entry's interpreter HashRef exclusively"; `grep -n 'func Decode' host/registry/registry.go` → **:75** `Decode(payload []byte) (Registry, error)`; `grep -n 'EpochRegistryV1' host/store/store.go` → **:82-84** `const EpochRegistryV1 = "world/epoch-registry/v1"`; `sed -n '126,133p' host/registry/registry.go` → `Bootstrap(… releaseString …)` builds `Epochs: []EpochRecord{{Epoch: firstEpoch, Candidates: []string{releaseString}}}`. **M-A1 CONFIRMED.**
- **V23 (F14 / M-A1)** `grep -n 'registry.Bootstrap' host/daemon/daemon.go` → **:512** `if _, _, err := registry.Bootstrap(ctx, s, release); err != nil {` (StageRegistry abort). **CONFIRMED.**
- **V24 (F15 / M-A1 precision + M-A2)** `sed -n '486,515p' host/daemon/daemon.go` → `release := unpinnedRelease` (:491, `const unpinnedRelease = "unpinned"` :179); with `cfg.AilangBin != ""`: `release = releaseFromVersion(m.Version)` (:505); `sed -n '613,624p'` → `releaseFromVersion` returns the first non-blank line, `strings.TrimSpace`d, else `"unpinned"`. `$HOME/.pinned-ailang/ailang --version` → observed multi-line (`AILANG v0.41.0` / `Commit: 24ee108` / `Full: …` / `Built: …`), so the registry's Candidates hold `"AILANG v0.41.0"`, NOT the verbatim probed string — **M-A1's "release is the archived interpreter's probed version string" is imprecise (not false): it is that string reduced by releaseFromVersion**; the prototype applies the same reduction (`releaseFromManifest`). **M-A2 CONFIRMED**: `sed -n '150,160p' host/archive/archive.go` → `Manifest.Version string \`json:"version"\` — "the verbatim `ailang --version` output captured at archival" (probeVersion :467; `Version: version` :374).
- **V25 (F16 / M-B2)** `read host/procbound/procbound.go` → `Admit() (release func(), err error)` (:42, MaxOutstanding = 8, non-blocking CAS reservation), `Wait(wait func() error, d time.Duration, release func()) error` (:60, `ErrCleanupIncomplete` hands the reap to a background waiter); `sed -n '450,495p' host/archive/archive.go` → probeVersion pattern: `exec.CommandContext` + `a.probeTimeout` (10s, :75) + `childenv.Scrubbed(os.Environ())` (:476) + separated buffers. Control that the archived copy is executable on this rig: `cp $HOME/.pinned-ailang/ailang /tmp/r2probe/artifacts/interpreters/copy-test && chmod 755 … && …/copy-test --version` → **exit 0, `AILANG v0.41.0`** (contrast V18: Apple-signed `/bin/echo` was SIGKILLed; the ad-hoc-signed pinned binary survives archival). **M-B2 CONFIRMED.**
- **V26 (F17 / objection B, pinned v0.41.0 + its archived copy)** all observed in-session: (a) `printf 'this is not AILANG source at all' > garbage.ail && ailang check garbage.ail` → **exit 1** (`Error: type error … undefined variable: this`); (b) `cp world/transitions.ail /tmp/r2probe/alone.ail && cd /tmp/r2probe && ailang check alone.ail` → **exit 1** (`Error: module loading error: failed to load world/logepoch: LDR001: module not found: world/logepoch`) — a standalone check of an import-bearing module fails for IMPORT-RESOLUTION reasons; from the repo root (imports resolvable) the same file checks **exit 0** (`✓ No errors found!`); (c) module-declared standalone source in a scratch dir → **exit 0** with `WARNING MOD010 (temp-path) … Auto-relaxed for temporary directory`; (d) module-LESS source (`export func main() -> string { … }`) → **exit 1** (`MOD014: no 'module' declaration`); (e) hyphenated module path → **exit 1** (`PAR_HYPHEN_IN_MODULE`). The honest scope chosen from (b)+(F18): hermetic check, import-bearing sources refused (they are not executable under the capsule either).
- **V27 (F13 / M-B1 re-verify)** `sed -n '30,80p' host/canon/source.go` → `Source` rejects BOM (`input begins with a UTF-8 BOM`), invalid UTF-8, embedded NUL; `splitLines` on LF/CRLF/CR; `trimTrailingSpaceTab`; trailing empty lines dropped — **no parser, no type-checker anywhere in the file**. **M-B1 CONFIRMED.** Known-positive control that canonicalisation alone proves nothing: the CLI fixture `TestTransitionsVerbRefusesGarbageSourceBeforePutObject` feeds bytes that `canon.Source` ACCEPTS (valid UTF-8, no NUL) and the publish is refused ONLY by the interpreter check (`Error: … PAR_/type error` captured in the typed refusal) — the same bytes sail through `canon.Source` to `PutObject` under MUT-8.
- **V28 (F18 / kimi secondary premise)** `sed -n '164,196p' host/transitionreg/transitionreg.go` → `byID[change.ID] = cloneDescriptor(*change.Descriptor)` — `BuildNext` REPLACES by ID, so a naive retry WOULD clobber the winner's same-ID entry (the defect kimi suspected is real in round-1's code); `sed -n '180,215p' host/capsule/capsule.go` → capsule stages ONLY the entry source (`entry.Source` at `host/capsule/main.ail`, `AILANG_FS_SANDBOX=root`); contrast `host/replay/replay.go:319` (`copyDir(worldLibDir, filepath.Join(root, "world"))`) which stages a world lib for its own recorded fixtures — the capsule engine (row 106's `Bind`) has no such staging.
- **V29 (r2 mutations re-run)** all twelve mutations were executed against the FINAL (post-split) prototype — mutate → targeted `go test` → revert — with these observed kills: MUT-1 → `--- FAIL: TestPublishSetGenesisThenGrowth`; MUT-2 (both presence arms removed) → `--- FAIL: TestPublishSetRefusesAbsentTransitionSource` AND `--- FAIL: TestPublisherRefusalKeepsCardEmpty`; MUT-3 → `--- FAIL: TestPublishSetIdempotentRepublishIsNoOp` AND `--- FAIL: TestTransitionsVerbHappyPathAndIdempotence`; MUT-4 → `--- FAIL: TestPublishSetGenesisThenGrowth`; MUT-5 → `--- FAIL: TestTransitionsVerbFences/ci_environment_stops` AND `/wrong_phrase_stops`; MUT-6 → `--- FAIL: TestTransitionsVerbFences/unverifiable_interpreter_ref_refused`; MUT-7 → `--- FAIL: TestPublishSetGenesisThenGrowth` AND `TestPublishSetIdempotentRepublishIsNoOp`; MUT-8 → `--- FAIL: TestPublishSetRefusesUnloadableTransitionSource` AND `--- FAIL: TestTransitionsVerbRefusesGarbageSourceBeforePutObject`; MUT-9 → `--- FAIL: TestPublishSetRefusesEpochNotNominatingTheInterpreter` AND `TestPublishSetRefusesDefaultEpochOne` AND `TestEpochRefusalKeepsCardUnchanged`; MUT-10 → `--- FAIL: TestPublishSetRefusesDefaultEpochOne` AND `TestTransitionsVerbEpochArms`; MUT-11 → the previous three unit refusals AND `TestPublishSetRefusesAbsentEpochRegistry` AND `TestTransitionsVerbEpochArms`; MUT-12 → `--- FAIL: TestPublishSetSameIDConflictOnCASRetryRefuses`. **12/12 killed**; every revert verified by `go build` + the targeted suite returning to green.
- **V30 (r2 full gate, post-split)** `go vet ./...` → **exit 0**. `AILANG_BIN=$HOME/.pinned-ailang/ailang go test ./... -count=1` → **exit 0, 23 packages ok, 0 FAIL** (`grep -cE '^ok  '` → 23; `grep -cE '^FAIL'` → 0) — the same 23 ok / 0 FAIL the controller measured for round 1's prototype with the pinned v0.41.0 binary. Without AILANG_BIN the only failures are the AILANG_BIN-gated rows (round 1's `TestWorldCoreManifestMatchesTheCommittedGolden` plus this row's new `TestTransitionsVerbRefusesImportBearingSourceUnderPinnedInterpreter`, which names the requirement in its message) — house pattern, never skipped.
- **V31 (r2 S8 couplings, both observed firing then green)** (a) AC10 subprocess-site census: `go test ./host/broker -run TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential` first FAILED — `subprocess site file "host/transitionreg/verify.go" has no driver` — until the driver was added; the check subprocess then landed in `host/archive/check.go` (a cycle measurement forced the move: `registry_publish_test.go imports host/transitionreg … imports host/broker from bind.go: import cycle not allowed in test` — observed verbatim), and the census went green with the new `host/archive/check.go` driver and its child env observed WITHOUT the registry credential sentinel. (b) `go test ./host/verifygate -run TestSubprocessSinkGate` first FAILED — `host/transitionreg/verify.go:250: subprocess sink "stdout"… is a bare bytes.Buffer … use the syncBuffer pattern` — fixed with the mandated concurrency-safe buffers; gate green. Both inventories moved WITH the change in the same diff (F12 discipline).
- **V32 (r2 LOC)** `grep -vE '^\s*(//|$)' <file> | wc -l` → `host/transitionreg/publish.go` **104**, `collide.go` **83**, `epoch.go` **116**, `verify.go` **83**, `host/archive/check.go` **78**, `cmd/world-publish/transitions.go` **256** (split across M6/M7); every milestone ≤ ~150.
- **V33 (r2 patch)** `git diff > .design/proto.patch` plus `git diff --no-index /dev/null <file>` for the ten new files → **14 file diffs** (`grep -c '^diff --git'` → 14), 2369 lines: 4 modified landed files (`cmd/world-publish/main.go`, `host/store/context_roots_test.go`, `host/broker/registry_publish_test.go`, `host/transitionreg/transitionreg.go`) + 10 new.
- **V34 (r2 no-inventory-move check)** the FROZEN flag surface is untouched by revision 2 (no new flags: `grep -c 'flagNames'` unchanged; AC24(b) green in the full run, V30) and `TestProductionContextRoots` is green with the ROUND-1 pin (`cmd/world-publish/transitions.go|runTransitions|Background` still exactly 1 — the revision-2 code threads the one `ctx := context.Background()` root through the new helpers instead of adding roots).
- **V35 (objection-A prototype, epoch)** observed green: `go test ./host/transitionreg/ -run 'TestEpochs|TestRelease|TestPublishSetRefusesEpoch|TestPublishSetRefusesDefault|TestPublishSetRefusesAbsent|TestPublishSetRefusesWithoutArchive'` → **ok**; daemon arm `go test ./host/daemon/ -run 'TestEpochRefusal|TestPublished|TestPublisherRefusal'` → **ok**; CLI arms `go test ./cmd/world-publish/ -run 'TestTransitionsVerbEpochArms|TestTransitionsVerbHappyPath'` → **ok** (happy-path stdout includes `semantics epoch 1 derived from world/epoch-registry/v1 for interpreter release "test-interpreter-version"`).
- **V36 (objection-B prototype, loadability + same-ID)** observed green: `go test ./host/transitionreg/ -run 'TestEnsure|TestPublishSetRefusesUnloadable|TestPublishSetSameIDConflict|TestPublishSetCASRetry'` → **ok**; the refusing fake's output is carried verbatim in the typed error (asserted); `go test ./cmd/world-publish/ -run 'TestTransitionsVerbRefusesGarbageSource|TestTransitionsVerbRefusesImportBearing|TestPublishErrorLine'` → **ok** (the import-bearing arm refused with `LDR001` in the message under the pinned binary).
- **V37 (publishable-content census, controller-measured at 0759ee7, quorum r2 kimi)** `grep -rn "^import" world/*.ail` → `world/transitions.ail:3 import world/logepoch`, `:4 import world/types`, `:8 import world/contracts`; `world/contracts.ail:3,4` import `world/logepoch`, `world/types`; `world/types.ail:3` imports `world/logepoch`. Control: `ls world/*.ail | wc -l` → **4** (the fourth, `world/logepoch.ail`, imports nothing and exports no transition). Cross-referenced with V26(b): every landed transition-source candidate imports `world/*` and is refused hermetically.
- **V38 (executor iter-195, sprint §4 finding: `CheckSource` verdict was host-path-dependent)** Fresh-dir probe, pinned `$HOME/.pinned-ailang/ailang` (AILANG v0.41.0, Darwin 25.6.0), each probe in a NEW `mktemp -d` entered via `/private$TMPDIR` (as `exec.Cmd.Dir` resolves it), file staged as `entry.ail`: strict (no relax) `module check/entry` → **rc=1** `Error MOD010: module 'check/entry' doesn't match file path 'entry'`; strict `module host/capsule/main` → **rc=1** MOD010; strict `module entry` → **rc=0** `No errors found`; `AILANG_RELAX_MODULES=1` `module check/entry` → **rc=0** (`WARNING MOD010 (relaxed)`); relax `module quickstart/echo` (the QUICKSTART §6 source) → **rc=0**; relax import-bearing (`import world/types (World)`) → **rc=1** `LDR001: module not found: world/types` — the hermetic import check keeps its teeth. Fix landed in M1 (`cmd.Env = append(childenv.Scrubbed(os.Environ()), "AILANG_RELAX_MODULES=1")`); removing it reds `TestCheckSourceUnderPinnedInterpreter` (M1-MUT-c, V41). This corrects F17/V26(c)'s "module-declared → exit 0", which held only from a shell whose `$PWD` kept the `/var/folders` alias. Transcript banked at `~/.ailang/state/world-iter195/exec-s4-probe.txt`.
- **V39 (executor iter-195, AC-EPOCH-TWIN)** `TestPublishSetEpochTwinInterpretersKeepDistinctPins` (M4, `host/transitionreg/epoch_test.go`) and `TestEpochTwinInterpretersBothListedOnCard` (M5, `host/daemon/registry_publisher_test.go`) pass on the landed code. TWIN-MUT (in `verifyEpochs`, overwrite `d.Interpreter` with the first interpreter seen for the same derived release) → `--- FAIL: TestPublishSetEpochTwinInterpretersKeepDistinctPins` AND `--- FAIL: TestEpochTwinInterpretersBothListedOnCard`; the same mutation with `go test ./host/transitionreg/ ./host/daemon/ ./cmd/world-publish/ -skip Twin` → `ok` ×3 (survives everything else), so the twin tests are load-bearing. Reverted; `git diff --stat HEAD` empty.
- **V40 (executor iter-195, full gate on the final sprint tree, `sprint/w-transition-registry-production-publisher` @ M7)** `go vet ./...` → **rc=0**; `AILANG_BIN=$HOME/.pinned-ailang/ailang go test ./... -count=1` → **rc=0, 23 ok / 0 FAIL** (no `host/pkgproj` flake this run); `AILANG_BIN=… ./scripts/verify_ail.sh` → **rc=0** (`world package gate PASSED: 9/9 steps`; `verify gate PASSED: 11 required identities verified, 40 named tests pass`) — unchanged, no `.ail` touched. Every milestone boundary M1–M6 passed its §2 gate (vet rc=0, `go test ./... -run '^$'` rc=0, targeted tests ok); M6's boundary also ran the full suite (23 ok) and verify_ail (rc=0). QUICKSTART §6 non-TTY steps: `go build ./cmd/world-publish` rc=0; the verb run with stdin `/dev/null` → rc=3 `STOP fence=tty reason=no-controlling-terminal` (the attended fence, observed; the store opens before the fence by design, as in `runPublish`).
- **V41 (executor iter-195, mutation drill on the final tree)** Every sprint-plan §6 row executed (mutate → `go vet` compile fence → named killers → `git checkout` → `git diff --stat HEAD` empty): M1-MUT-a → `TestCheckSourceReportsTheArchivedInterpretersVerdict`, `TestCheckSourceUnderPinnedInterpreter`, `TestEnsureSourceLoadableStagesHermetically`, `TestPublishSetRefusesUnloadableTransitionSource`; M1-MUT-b → `TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential/host/archive/check.go`; M1-MUT-c → `TestCheckSourceUnderPinnedInterpreter`; MUT-1 → `TestPublishSetGenesisThenGrowth`; MUT-2 → `TestPublishSetRefusesAbsentTransitionSource` + `TestPublisherRefusalKeepsCardEmpty`; MUT-3 → `TestPublishSetIdempotentRepublishIsNoOp` + `TestTransitionsVerbHappyPathAndIdempotence`; MUT-4 → `TestPublishSetGenesisThenGrowth`; MUT-5 → `TestTransitionsVerbFences/ci_environment_stops` + `/wrong_phrase_stops`; MUT-6 → `TestTransitionsVerbFences/unverifiable_interpreter_ref_refused`; MUT-7 → `TestPublishSetGenesisThenGrowth` + `TestPublishSetIdempotentRepublishIsNoOp`; MUT-8 (core) → `TestEnsureSourceLoadableStagesHermetically`; MUT-8 → `TestPublishSetRefusesUnloadableTransitionSource` + `TestTransitionsVerbRefusesGarbageSourceBeforePutObject`; MUT-9 → `TestPublishSetRefusesEpochNotNominatingTheInterpreter` + `TestPublishSetRefusesDefaultEpochOne` + `TestEpochRefusalKeepsCardUnchanged`; MUT-10 → `TestEpochsForInterpreterRefusesAbsentAndUnnominated` + `TestPublishSetRefusesDefaultEpochOne` + `TestTransitionsVerbEpochArms/omitted_epoch_with_no_nomination_refused_(no_default)`; MUT-11 absent arm → `TestEpochsForInterpreterRefusesAbsentAndUnnominated`; MUT-11 (all arms) → `TestPublishSetRefusesEpochNotNominatingTheInterpreter` + `…DefaultEpochOne` + `…AbsentEpochRegistry` + `TestTransitionsVerbEpochArms/{omitted_epoch_with_no_nomination_refused_(no_default),absent_epoch_registry_refused_typed}`; MUT-12 → `TestPublishSetSameIDConflictOnCASRetryRefuses`; TWIN-MUT → V39; row example "card ignores the capability filter" (`req.Allowed()` → `req.Registry.List()` in `host/projection`) → `TestPublishedTransitionsAppearOnAgentCard`. **19/19 KILLED, none by the compiler.** Three first-run authoring errors were corrected and re-run, not counted: MUT-5's first form was compiler-killed (`*stopError` typed nil needed); M1-MUT-a's first `host/archive` arm used an anchored regex that ran zero tests; MUT-11's first form lacked the "non-nominating epoch coerced to 1" arm. Transcript banked at `~/.ailang/state/world-iter195/exec-mutations.txt`.

## Controller-fact audit (F1–F7 as supplied)

F1 **CONFIRMED** (empty + control 26). F2 **CONFIRMED with a sharper finding**: no genesis
constructor exists, and `BuildNext`-over-genesis is *structurally incompatible* with
`Publish`'s absent-head branch (V2, MUT-7) — the first revision must be constructed directly.
F3 **CONFIRMED** (fields + Validate :266). F4 **CONFIRMED** (daemon.go:532; AgentCard at :218,
doc block at :212). F5 **CONFIRMED** (verbs; 0 transitionreg imports in both cmds; worldd CLI
is an HTTP client). F6 **CONFIRMED** (4 exports, :23/:45/:57/:76). F7 **CONFIRMED** (P9 verbatim;
this row is P9's anticipated "later path" and adds no REST route — the verb is a local,
attended CLI). **None of F1–F7 was found FALSE.**

## Controller-measurement audit (revision 2: M-A1, M-A2, M-B1, M-B2 as supplied)

**M-A1 CONFIRMED, with one precision correction (not FALSE):** the epoch registry exists exactly
as measured — `store.EpochRegistryV1` (store.go:84), `Registry{Epochs []EpochRecord}` with
`Candidates []string` (advisory release strings, Candidates[0] primary), `registry.Decode`
(registry.go:75), daemon bootstrap at daemon.go:512 (V22/V23). The precision: the bootstrapped
`release` is NOT the verbatim probed version string — it is `releaseFromVersion(m.Version)`,
the FIRST NON-BLANK LINE, TRIMMED (daemon.go:505, :613-624), and `"unpinned"` when no
interpreter is configured (V24). The design consequence is load-bearing and measured in:
the prototype's `releaseFromManifest` applies the same reduction, or every multi-line
`ailang --version` output would false-mismatch the registry's Candidates.
**M-A2 CONFIRMED:** the manifest field is `Manifest.Version` (`json:"version"`, archive.go:155-156)
— "the verbatim `ailang --version` output captured at archival", written by `probeVersion`
(archive.go:467, `Version: version` at :374) (V24).
**M-B1 CONFIRMED:** `canon.Source` (source.go:47-80) is the eight-step normaliser only — BOM,
UTF-8, NUL rejection; LF/CRLF/CR line splitting; trailing space/tab trim — no parse, no
type-check (V27, with the known-positive control that non-AILANG clean-UTF-8 bytes sail
through canonicalisation and are refused only by the interpreter check).
**M-B2 CONFIRMED:** `procbound.Admit/Wait` (procbound.go:42/:60, MaxOutstanding 8) is the
bounded-subprocess helper, and `host/archive`'s `probeVersion` (10s `CommandContext`, scrubbed
env, separated buffers) is the in-repo pattern — the prototype's check uses both, and the
archived copy of the pinned binary executes on this rig (V25). **None of the four controller
measurements was found FALSE** (M-A1 imprecise on the release-string reduction, as above).