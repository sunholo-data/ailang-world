# w-self-mod-vertical — one public World package through the brokered self-modification lane

**Status**: Planned
**Date**: 2026-08-05
**Queue item**: 8, `w-self-mod-vertical`, clause-7
**Estimated**: 4–5 days total (SM.B is 2–2.5 d after the round-1 revision; see the Quorum
verification log's standing caveat (2) — whether that wants splitting is a planner question, not a
settled one)
**Designer**: `codex:gpt-5.6-sol` (rotation slot 2). **Quorum**: 3 rounds — BLOCKED, BLOCKED (N−1,
`gpt5-6-sol` absent: budget), then controller carve-out applying the reviewer's own prescribed fix.
**Not yet sprint-planned.**
**Compiler boundary**: every command in this design uses the released
`/tmp/ailang-v0300/ailang` (`AILANG v0.30.0`, commit `e37b370`), never PATH `ailang`.
**Upstream source boundary**: AILANG source claims are measured at
`e37b370d1d7a9c4e7136b319e38bec4d5f2bd9a0`.

> **Thesis:** World will prepare, authorize, publish, receipt, and hash-verify exactly
> `world/core@0.1.0`; the registry POST is an outward broker effect, never a core dependency or an
> ambiently authorized shell command.

---

## The finding in one paragraph

The public registry can carry the intended artifact, but it cannot confer the authority the phrase
“`world/` vendor namespace” suggests. In v0.30.0 there is one optional shared API key, no accounts,
no vendor registration, no scopes, and a server comment explicitly deferring namespace auth. A live
controlled four-arm dry-run accepts `world/probe`, `someoneelse/probe`, and `sunholo/probe` alike,
while rejecting a non-vendor-shaped name. Thus **`world/` is a string World writes, not a namespace
World holds**. This design does not launder that convention into ownership. It makes World's own
publisher reject every non-`world/` manifest, routes real namespace authorization upstream, and
places the irreversible POST behind the landed capability/budget/receipt broker plus an attended
human stamp. The core stays local: it produces a deterministic package projection and receipt; only
the handler knows the credential or network, and consumers install an exact version whose three
hashes match the committed publish receipt.

## Motivation

Clause 7 is not satisfied by placing four files in a tarball or by running a registry command by
hand. The purpose of the vertical is to demonstrate that World can change its own distributable
surface through the same proposal → verify → authorize → effect → receipt chain it demands of every
other mutation.

The first artifact is deliberately modest. `world/core@0.1.0` packages the already-landed pure
semantic library:

- `world/types`
- `world/contracts`
- `world/transitions`
- `world/logepoch`

That choice makes this a distribution and governance extension, not a new kernel behavior. It is
still useful: a remote consumer can resolve the exact public package while local World continues to
execute the canonical checked-in modules. Future behavioral extensions remain separate packages;
this vertical proves their lane without inventing one in the same sprint.

The difficult part is the irreversible boundary. A request can reach the registry, succeed, and
lose its response. Since the version is immutable, a blind retry then returns 409. The landed
three-state receipt law already names this truth: an effect intent without an outcome is
`indeterminate`. The design must preserve that state, probe public metadata, and only then resolve
the receipt. Treating a timeout as an ordinary failed handler result would falsely authorize a retry.

## Premises

- **P1 — naming is compatibility, not ownership.** `world/` is unclaimed in the current public
  index, but v0.30.0 has no vendor registration or per-vendor authorization.
- **P2 — the package can keep current module names.** Manifest validation accepts exported modules
  beginning with either the package name or a single-segment `module_prefix`; therefore package
  `world/core` can set `module_prefix = "world"` and export the four modules unchanged.
- **P3 — a root manifest is too broad.** v0.30.0 tarball creation recursively includes every
  `*.ail`, so making the repository root the package would ship checked design sketches as well as
  the four intended modules.
- **P4 — the package is a projection, not a second source of truth.** Canonical files remain
  `world/*.ail`; `packages/world-core/world/*.ail` are byte-identical publish inputs checked by a
  hard gate.
- **P5 — the broker already supplies serialized decision/debit/dispatch/record flow.** A live
  allowed invocation durably appends an effect intent before handler dispatch and an effect outcome
  afterward; replay never dispatches handlers.
- **P6 — ordinary handler failure is the wrong representation for an ambiguous POST.** The current
  broker writes a resolved `failed` outcome for every returned handler error. Publish needs a typed
  indeterminate return that deliberately leaves the durable intent outcome-free.
- **P7 — the store's three states are load-bearing.** `not-started` has no intent,
  `indeterminate` has an intent and no outcome, and `resolved` has both.
- **P8 — local gates have different scopes.** `verify_ail.sh` checks the canonical `world/` modules
  and design sketches; `verify_go.sh` builds/tests all Go packages and pins replay to v0.30.0. A new
  package projection needs its own manifest/export/drift/smoke gate rather than an assumption that
  either existing script sees it semantically.
- **P9 — public consumption must not enter the kernel.** Registry resolution and install belong to
  an acceptance fixture or extension workspace, never imports in canonical `world/*.ail`,
  `host/store`, `host/replay`, or `cmd/ailang-worldd`.
- **P10 — a public write is human-sovereign.** The headless loop may make the package completely
  ready and emit a dry-run evidence bundle; it may not mint the attended authorization stamp or
  call the live handler.

### Design Freeze

- The one artifact is **`world/core@0.1.0`**. The version is immutable once published.
- Its manifest uses `module_prefix = "world"`; module declarations and imports remain unchanged.
- The canonical source remains `world/*.ail`. The package projection must be byte-identical.
- The outward effect is exactly `Registry.Publish`; cost is one irreversible publish attempt.
- Only manifests whose package vendor is exactly `world` can pass World's publish gate.
- This is convention enforcement on World, **not** registry namespace ownership.
- A public POST requires one attended, content-bound approval stamp. Dry-run does not.
- No automatic retry follows an indeterminate publish.
- Recovery is probe-then-resolve, never retry-then-see.
- No registry URL, client, credential, package cache, or installed registry dependency enters the
  core runtime.
- No live publish occurs during implementation or CI. The first live invocation is an attended
  operational step after all milestones land green.

## Decision 1 — the artifact is `world/core@0.1.0`

Create `packages/world-core/ailang.toml` with this semantic content:

```toml
[package]
name = "world/core"
version = "0.1.0"
edition = "1"
ailang = ">=0.30.0"
module_prefix = "world"
description = "AILANG World's pure semantic core"

[exports]
modules = ["world/types", "world/contracts", "world/transitions", "world/logepoch"]

[effects]
max = []
```

**`edition` was left as an implementer TODO by the designer and is RESOLVED HERE BY MEASUREMENT
`(controller iter-51, 2026-08-05)`.** The draft carried `edition = "2025"` with the honest caveat
that validation had not been checked. Measured at the pinned commit: `internal/pkg/manifest.go:245`
validates `edition` for **non-emptiness only** — `if m.Package.Edition == "" { return
"[package].edition is required" }` — there is no enum and no allowlist (same-file known-positive
control: `ModulePrefix`, which IS structurally validated, appears **12** times). So `"2025"` would
have passed the parser while diverging from the toolchain's own template: `internal/pkg/manifest.go:424`
writes `edition = "1"`, and `ailang init package` produced exactly `edition = "1"` in the
controller's live V-K probe. **The field is therefore a convention the validator cannot police, which
makes copying the template value the only defensible choice** — and it is a small standing example of
this document's central theme: a check that passes everything confers no authority. Set to `"1"`.
The implementer still confirms against a published-package fixture before landing, but it is now a
confirmation rather than an open question.

The narrow package directory contains:

```text
packages/world-core/
  ailang.toml
  _smoke.ail
  world/types.ail
  world/contracts.ail
  world/transitions.ail
  world/logepoch.ail
```

`_smoke.ail` imports at least one exported symbol from every module and exercises the pure
plan→verify→commit flow. Its syntax is copied from the pinned, already-checked replay fixture and
then checked with the pinned binary; no AILANG syntax is invented in prose.

Why copy rather than move: the canonical source root and replay fixture already resolve
`world/*` from repository root. Moving the files would turn a packaging vertical into a module-root
migration and force replay-golden churn for no semantic gain. Why copy rather than use a root
manifest: the tarball walker recursively includes all `.ail`, including checked design sketches.

The copy is generated deterministically by `scripts/build_world_package.sh`. That script starts
from an empty staging directory under `packages/world-core/`, copies exactly the four allowlisted
canonical modules, and refuses any fifth `.ail` other than `_smoke.ail`. The committed projection
allows review of the actual public bytes; the generator and drift gate prevent dual-source drift.

## Decision 2 — `world/` remains a convention until upstream adds namespace auth

Choose a combination of option (a) and option (b): enforce the vendor convention in World's own
gate now, and route a concrete upstream requirement for real namespace authorization. Reject option
(c), a claim-marker package, because an immutable marker proves chronology but grants no authority
and consumes an unrecoverable public name while offering no exclusion.

The local gate parses `ailang.toml` structurally and requires all of:

1. `[package].name == "world/core"` for this vertical.
2. Splitting the name once on `/` yields vendor exactly `world`.
3. Version equals the attended target `0.1.0`.
4. `module_prefix == "world"`.
5. Export set equals the four-item frozen manifest, with no missing or extra module.
6. Effects max is empty.
7. Projection files are byte-identical to canonical sources.

This protects against World accidentally or maliciously publishing `sunholo/*`,
`someoneelse/*`, or a different `world/*` artifact through the World handler. It does **not** stop
another holder of the shared registry key from publishing under `world/`, reserve the string, prove
publisher identity, or provide package signing.

The upstream ask is: per-vendor principals, vendor registration/claim, scoped credentials, and a
server-side authorization check binding authenticated principal → vendor before immutability check.
Until that exists, receipts say `vendorEnforcement: "world-local-convention/v1"`; they never say
“owned”, “registered”, or “authenticated as world”. The upstream ask does not block this ratified
vertical because the human accepted the current public flow; it blocks any later claim of namespace
security.

## Decision 3 — the publish broker contract

### Effect, scope, cost, and grant

| Field | Frozen value / grammar |
|---|---|
| Effect | `Registry.Publish` |
| Scope | `registry:<registry-origin>/package:world/core/version:0.1.0` |
| Cost | `1` |
| Cost unit | one live irreversible HTTP publish attempt, charged before dispatch |
| Capability | exact effect and exact scope; `Budget >= 1`; logical expiry checked against caller-supplied `Now` |
| Budget | attended grant starts at `1`; success, definite failure, and indeterminate dispatch all consume it |

The handler refuses wildcard origins, redirects to a different origin, non-HTTPS live origins,
credentials embedded in URLs, alternate package/version payloads, and `--allow-dotted-tool-names`.
Testing may inject a loopback fake validator under a distinct test-only constructor; production
construction accepts only the compiled/approved public validator origin.

The `EffectRequest` is:

```text
Effect = Registry.Publish
Scope  = registry:<origin>/package:world/core/version:0.1.0
Cost   = 1
Now    = caller-supplied logical authorization time
```

The canonical request payload is fixed-field-order JSON containing:

```text
schema, vendor, name, version, registryOrigin,
manifestRef, approvalRef, tarballSHA256, contentHash, interfaceHash,
exports, effects, compilerVersion, compilerSHA256
```

`compilerVersion` is `AILANG v0.30.0` and `compilerSHA256` is the full V-A digest. The handler
recomputes every hash immediately before dispatch and rejects any mismatch with the approval stamp.

### What the EffectRecord binds

The landed `EffectRecord` directly binds effect, scope, cost, budget before/after, request ref,
and result ref. The request object transitively binds vendor, name, version, manifest ref, attended
approval ref, tarball SHA-256, content hash, interface hash, exports, effects, compiler identity, and
registry origin. A success result object contains the same identity and hashes plus HTTP status,
server metadata URL, observed immutable metadata digest, and reconciliation mode. Therefore the
EffectRecord binds all requested publish observables through content-addressed `RequestRef` and
`ResultRef` without widening the frozen record codec.

The dry-run parser must not trust truncated display prefixes. It **independently recomputes** the
three hashes in `host/pkgproj` and uses the CLI only as a cross-check. A receipt containing only
V-J's displayed prefixes is invalid.

> **`DD-1` (SM.A, landed) — the "library extraction" branch this paragraph used to offer does not
> exist, and three of this document's own citations are the proof.** The sentence above originally
> read *"either obtains full hashes from a small library extraction of v0.30.0 package hashing
> logic or independently recomputes them"*. The first branch is impossible: World is
> `module github.com/sunholo-data/ailang-world`, upstream is `module github.com/sunholo-data/ailang`,
> and the hashing lives in `internal/pkg/{hasher,tarball}.go` — which Go's internal rule forbids
> importing across module paths. The Premise Verification Log cites those exact files three times
> *as evidence that the mechanics work*. They do work; the word `internal` in the path is
> simultaneously the reason the plan built on them cannot. **When you cite a path as evidence,
> read it once more as a policy.**
>
> The CLI is no fallback either: `e37b370:cmd/ailang/pkg_publish.go:110-112` prints
> `hash[:24]+"..."` (`"sha256:"` is 7 chars, so 17 hex nibbles = 68 bits survive), and the tarball
> bytes are never persisted — the only `.tar.gz` reference in that file is the multipart form on
> the **upload** path, which `--dry-run` returns before reaching.
>
> `AC6` is therefore satisfied by a **re-implementation** in `host/pkgproj` carrying a mandatory,
> hard-failing cross-check against those 24-char prefixes. **The risk this created has been
> measured and is CLOSED:** the tarball hash rides `compress/gzip`'s DEFLATE output and
> `archive/tar`'s format selection, and World builds at `go 1.25.6` while upstream's module
> declares `go 1.26.5`. Measured 2026-08-05 (iter-53) against the pinned binary, all three arms
> agree, tarball length `5472 = 5472` — and the instrument was shown able to red:
> `MUT-SM-PKGPROJ-CONTENT-SEPARATOR` (`file:%s\n` → `file:%s`) reds the **content** arm alone and
> `MUT-SM-PKGPROJ-TAR-MODE` (`Mode: 0644` → `0600`) reds the **tarball** arm alone, each naming
> both values, with the control green before and after and `pkgproj.go` byte-identical on restore.
> Decomposing the comparison into three arms is load-bearing rather than cosmetic: content and
> interface are pure sha256 over bytes and **cannot** diverge by toolchain, so a content/interface
> mismatch means the re-implementation is wrong while a tarball-only mismatch is a genuine
> cross-toolchain finding. The two have opposite remedies, and one merged verdict destroys exactly
> the information that tells them apart.

### Definite and ambiguous outcomes

| Observation | Receipt transition | Retry? |
|---|---|---|
| local validation fails before POST | resolved `failed-before-dispatch`; budget policy records no live attempt | allowed only after a new verified proposal |
| server returns non-2xx other than matching immutable 409 | resolved `failed`; attempt consumed | no automatic retry |
| server returns success and metadata verifies | resolved `succeeded` | never |
| timeout/reset/cancel after request body may have left process | remain `indeterminate` | forbidden until reconciliation |
| process crashes after durable intent and before outcome | remain `indeterminate` | forbidden until reconciliation |
| reconciliation finds exact version with all hashes equal | resolved `succeeded-reconciled` | never |
| reconciliation gets 404/absence after bounded consistency window | resolved `not-published`; new attended grant required for retry | never under old grant |
| reconciliation finds same version with any hash mismatch | resolved `conflict`; human/operator incident | never |

The broker gains a typed `IndeterminateEffectError`. `Session.Invoke` recognizes only that type
after a durable intent, does not create a failed EffectRecord, does not append an outcome, and
returns the effect invocation ID. Every other handler error retains landed behavior. Budget remains
debited in the in-memory session; after restart, retry refusal comes from the pending receipt plus
the consumed one-shot authorization, not from reconstructing mutable grant state.

### Probe-then-resolve

**REVISED under the quorum round-2 carve-out `(controller iter-51, 2026-08-05)`.** Reviewer
`gemini-3-1-pro` blocked round 2 on this exact paragraph: *"The premise that the upstream AILANG
registry serves package metadata at the exact HTTP path
`packages/world/core/0.1.0/metadata.json` is unverified… Attempting to reconcile an indeterminate
publish state by hardcoding a hallucinated endpoint will result in false 404s, incorrectly resolving
the state to 'not-published' and completely breaking the recovery boundary."* The controller ran the
reviewer's own prescribed measurement (row **V-N** below) before revising. The verdict is more
interesting than either "right" or "wrong": **the path string is correct, but its NATURE was wrong,
and the nature is what the recovery boundary rests on.**

`packages/{vendor}/{name}/{version}/metadata.json` is a **GCS bucket object key**, not a validator
HTTP route. The validator's router registers exactly `/publish`, `/unpublish`, `/rebuild-index`,
`/health`, `/version`, `/api/packages`, `/api/packages/`, `/api/stats`
(`cmd/registry-validator/main.go:58-67` @ `e37b370…`) — **no `/packages/…` route exists**, and every
in-tree use of that key (`main.go:168`, `cache.go:159`, `unpublish.go:190`) is a server-side
`bucket.Object(...)` call. The probe therefore targets the **read-only public bucket** under
`$AILANG_REGISTRY`, i.e.
`https://storage.googleapis.com/ailang-registry/packages/world/core/0.1.0/metadata.json`, and it
must never be described as, pointed at, or failed over to the validator service.

`ReconcileRegistryPublish` accepts an existing effect invocation ID and performs a bounded,
read-only GET of that bucket object. It never invokes the publish handler. It verifies
vendor/name/version and tarball/content/interface hashes against the durable request object — all
three are present in the served document (V-N). An exact match writes a success result object and
resolved outcome. A mismatch writes a conflict result and resolved outcome.

**Absence is the dangerous branch, and a bare 404 may NOT be read as absence.** Resolving
`not-published` is what re-authorizes an irreversible POST, so a false 404 is the one failure in
this design that can double-publish. Measured: a 404 from this bucket returns **GCS XML**
(`<Code>NoSuchKey</Code>`), not JSON — so a reconciler that parses JSON on the error path gets a
parse failure, not a clean signal — and that same 404 is indistinguishable from a typo'd object key,
a wrong or unset `$AILANG_REGISTRY`, a captive portal, a bucket-permission change, or a DNS
interception. Therefore, applying this mission's own instrument discipline to a **runtime** path
rather than only to controller commands: **every absence sample must carry a known-positive control
fetched in the same pass** — a package/version measured to exist (the design's fixture control is
`sunholo/auth@0.4.1`, V-N arm 1) or the registry `index.json`. If the control does not come back
`200` with well-formed JSON, the sample is **UNINFORMATIVE**, not absent: it is discarded, it does
not count toward the bounded window, and exhausting the window on uninformative samples resolves to
a named `probe-unavailable` state that requires a human — never `not-published`. Only a window whose
samples are ALL absent-with-a-firing-control may resolve `not-published`. An empty result from an
instrument that was never shown to work is a claim, not a measurement; here that claim costs a
duplicate immutable publish.

A returned 409 is evidence only after the same metadata probe. It is not automatically success:
someone else may have won the immutable version with different bytes. This rule is especially
important because `world/` has no registry-level owner.

The generic `Recover` function remains fail-closed and non-dispatching. It reports the pending
publish; explicit reconciliation is a separate operation with its own read-only network capability,
`Registry.ReadMetadata`, scope equal to the exact metadata URL, and bounded budget.

## Decision 4 — de-ambient the registry credential

The mission-loop environment and ordinary subprocess environment must have
`AILANG_REGISTRY_API_KEY` unset. The secret moves to a mode-0600 file outside the worktree and is
opened only by an injected `RegistryCredentialProvider` used by the publish handler. The provider
returns bytes in memory; it never logs them, stores them in an object, includes them in a request
hash, or exposes them to dry-run.

The handler launches the pinned binary with an explicit minimal environment. It strips every
registry-related variable, then adds `AILANG_REGISTRY_API_KEY` only to that child for the live
publish call. Dry-run always runs under `env -u AILANG_REGISTRY_API_KEY` and has no credential
provider. All other World-launched processes inherit the stripped environment.

Production startup fails if the ambient variable is set. A separate attended launcher may migrate
the existing ambient value into the secret file without printing it, unset the variable, verify file
mode and ownership, then exec World. That migration is operational and human-attended; it is not CI.

The honest residual risk: process-level isolation is not a same-UID security boundary. The Go daemon
that owns the provider can read the file, its own memory, and child environment; another same-UID
debugger may also do so. This gate buys removal of accidental inheritance and ensures ordinary
agents/shell commands cannot publish merely because they inherited the mission shell. Stronger
secret isolation requires a separate credential broker/OS principal and is deferred.

## Decision 5 — the attended authorization gate

The headless loop may perform all of these steps autonomously:

1. Build the package projection from the four allowlisted canonical modules.
2. Run canonical AILANG and Go gates.
3. Run package drift, manifest, export, smoke, tarball-content, and hash gates.
4. Run `/tmp/ailang-v0300/ailang publish --dry-run` with the key explicitly unset.
5. Create a ready-to-publish packet binding exact package/version, registry origin, manifest ref,
   compiler identity, full tarball/content/interface hashes, export/effect lists, proposed scope,
   and expiry.
6. Persist the proposal and evidence, then stop in `READY_AWAITING_HUMAN_PUBLISH`.

An attended human must review that packet and mint a `HumanApproval` object whose content hash is
`approvalRef`. Silence, timeout, or a changed byte is denial. The approval is single-use, exact-scope,
exact-hash, version-specific, expires at a logical time, and authorizes only one cost unit.

### Reuse of landed approval machinery and durable consumption

`HumanApproval` is the role name above, not a new wire object. The implementation reuses
`EffectHumanApprove` and `EffectHumanPollApproval`, `HumanHandler`, `DecideApproval`,
`ApprovalRequestV1`, `ApprovalDecisionV1`, `approvalRequestWire`, and `approvalDecisionWire`
unchanged. `EffectHumanApprove` persists the existing immutable request, with the publish effect,
exact hash/version scope, cost, requester, and logical time carried by the existing request wire;
the operator uses the existing non-effect `DecideApproval` entry point to create exactly an
`ApprovalDecisionV1` of `approve` or `deny`; polling remains the existing observation path.
`approvalRef` is the content hash of that landed decision object. No parallel approval codec or
effect is introduced. The extension is consumption state, not either approval wire: before use,
the broker verifies that the decision is `approve`, follows its `RequestRef` to the corresponding
`ApprovalRequestV1`, and checks the request's effect, canonical scope (including packet ref,
package/version/hashes, and expiry), cost, and request time against the publish packet and current
logical time. Thus none of the landed wire types is extended; the publish binding is a canonical
value of the existing `Scope` field, and old non-publish request/decision bytes remain valid.

The store adds
`AppendClaimedEffectIntent(episodeID string, intent EffectIntent, approvalRef, requestRef hashref.HashRef) (string, int64, error)`.
In one SQLite transaction it derives the normal `invocationID`, inserts the canonical effect-intent
object and journal row exactly as `AppendNextEffectIntent` does, and executes exactly
`INSERT INTO approval_claims (approval_ref, request_ref, invocation_id) VALUES (?, ?, ?)`.
The exact DDL is
`CREATE TABLE IF NOT EXISTS approval_claims (approval_ref TEXT PRIMARY KEY, request_ref TEXT NOT NULL, invocation_id TEXT NOT NULL UNIQUE);`.
Therefore the durable claim tuple is
`(approvalRef, requestRef, invocationID)`, and a uniqueness failure is returned as the named
`ErrApprovalAlreadyConsumed`. Transaction rollback leaves neither claim nor intent. Transaction
commit makes both visible atomically across goroutines, processes, sessions, restart, replay, and
indeterminate recovery. This deliberately changes `host/store/schema.sql`; the milestone must also
advance `currentSchemaVersion` and its independent ledger from 1 to 2 and update the exact-DDL
fixture required by the landed `w-ddl-gate-teeth` gate, which must red on an unratified schema edit.

For `Registry.Publish`, `Session.Invoke` uses that combined operation instead of
`AppendNextEffectIntent`: validate grant and approval first; atomically claim approval and append
the durable effect intent; debit the already-decided one-shot budget; only then load the credential
and call POST. A crash after the transaction commits but before dispatch burns the approval: it is
never refunded, deleted, or moved to another invocation. `Recover` reports the committed intent as
indeterminate; `retryAllowed(indeterminate, reconciled)` remains false until read-only
reconciliation resolves it, and a fresh-session attempt with the stamp fails
`ErrApprovalAlreadyConsumed` before credential load or dispatch. This deliberately prefers a
possibly unused burned approval to a duplicate immutable POST.

The blind, last-writer-wins `appendApprovalHead` race is a pre-existing defect surfaced by this
design. The claim table makes publish consumption safe without relying on that head as a lock;
repairing concurrent approval-head append itself is a separate queue row, not part of this sprint.

Only after that stamp may the controller construct the one-shot capability and call
`Session.Invoke`. The handler checks the approval object again. A headless loop may reconcile an
already-authorized indeterminate attempt using the separate read-only capability, but it may not
mint another live-publish grant. A `not-published` resolution returns to a new attended approval;
it does not resurrect the old stamp.

The first `world/` publish is therefore operationally ready at milestone close but deliberately not
performed by CI or the headless implementation agent. Clause 7 completes only when the attended
live invocation produces a resolved, hash-verified receipt and a clean-room pinned install verifies
the same hashes.

## Decision 6 — local-first consumption and the clause-2 boundary

Canonical World execution continues importing repository-local `world/*`. No core package imports
a registry client, performs HTTP, reads a package cache, or resolves `latest`. The public package is
an output projection and an optional input to external extension workspaces.

The clean-room acceptance fixture creates a scratch consumer outside the repository, declares the
exact dependency `world/core = "0.1.0"`, installs it with the pinned binary, and asserts the
installed metadata/tarball/content/interface hashes equal the resolved publish receipt before
checking and testing the consumer. The fixture never uses `latest`, a range, or trust-on-first-use.

The invariant would be violated by registry/package-cache imports or resolution code in any of:

- `world/types.ail`, `world/contracts.ail`, `world/transitions.ail`, `world/logepoch.ail`
- `host/store/**`
- `host/replay/**`
- `cmd/ailang-worldd/**`
- `go.mod` / `go.sum` if a cloud registry SDK is added

A dependency allowlist test scans those paths and fails on HTTP/cloud/registry client imports,
registry environment variables, `latest`, or package-cache lookup. `host/broker` is the intentional
boundary exception, limited to the publish and metadata-read handlers.

## Decision 7 — verification and evidence assembly

Add `scripts/verify_world_package.sh`, called explicitly by `scripts/verify_ail.sh` after its
existing two legs. The package verifier:

1. Fails if the package directory, manifest, smoke, or any allowlisted module is absent.
2. Fails if any unexpected `.ail` file exists.
3. Compares SHA-256 of each projection with its canonical source.
4. Parses the manifest and asserts the exact frozen fields and export set.
5. Checks and tests from the package source root with the pinned binary.
6. Runs `_smoke.ail` with a hard wall-clock bound.
7. Executes dry-run with `env -u AILANG_REGISTRY_API_KEY` and asserts the package identity, exports,
   empty effects, and three full recomputed hashes.
8. Lists tar entries and asserts the exact allowlist, proving design sketches cannot leak.
9. Writes a deterministic ready packet to the store during the operational workflow; CI compares
   canonical JSON golden bytes but does not persist mission state.

> **`DD-4` (SM.A, landed) — a THIRD LEG, never a new `ROOTS` entry.** The obvious way to make
> `verify_ail.sh` see `packages/` is to add `".|packages"` to its `ROOTS` array. That reds the
> repo's primary gate for a reason unrelated to the code under test: `verify_ail.sh:160`
> (`EXACT_TOTAL_VERIFIED=4`) and `:190` (`EXACT_TOTAL_TESTS = 14`) are **exact equalities, not
> floors**, and the four projections are byte-identical copies of the canonical modules — so they
> re-verify the same four contracts and the total becomes 8. Landed as an explicit third leg
> invoking `scripts/verify_world_package.sh`; measured after landing, the gate still reports
> `4/4` identities and `14` named tests.
>
> **`DD-5` (SM.A, landed) — the `.ailang` cache is inside the tarball's walk.** `CreateTarball`
> skips only directories named `.git`, `tests`, `test` — **not `.ailang`** — while step 5 of this
> very gate (check + test from the package root) is what *creates*
> `packages/world-core/.ailang/cache/**`. A stray `*.ail` there would silently enter the tarball
> and move its hash. Step 8 therefore asserts zero cached `*.ail` **with a firing control** in the
> same run (measured: control fired on 33 non-`.ail` files, zero `.ail` observed).
>
> **`DD-7` (SM.A, landed; found by the controller at landing, named by no prior reviewer) — the
> compiler pin is PLATFORM-SPECIFIC, and CI job 1 does not have the pinned compiler at all.**
> Two independent facts, both measured 2026-08-05:
> **(a)** the rig is darwin/arm64 (Mach-O) and CI is linux/amd64 (ELF), so a single
> `COMPILER_SHA256` constant is a gate that can only ever pass on one of the two. It is now a
> per-platform table with a loud unknown-platform failure — darwin/arm64
> `e9746fef8570bc42…`, linux/amd64 `1e594d158dffa688…`, the latter measured by downloading
> `releases/download/v0.30.0` and verifying its published `.sha256` (`OK`) as the control.
> Non-vacuity: a byte-flipped copy of the pinned binary reporting an identical
> `AILANG v0.30.0` string is rejected, naming both SHAs.
> **(b)** CI job 1 installs `releases/latest`, and the step log at `af0c3b4` (run `30993399332`)
> prints **`AILANG v0.33.0`** for job 1 against **`AILANG v0.30.0`** for job 2 in the same run.
> `latest` moved on 2026-08-04, so **queue item 9's "latent, not active" assessment expired**.
> v0.33.0 additionally fails this gate's own step 5 (measured on the rig: *"5 properties never ran
> (no generator)"*). CI job 1 therefore installs the pin to a separate path and passes it to this
> leg alone via `WORLD_PKG_AILANG_BIN`, leaving what legs 1-2 verify against untouched — that
> broader change is item 9's and is human-gated.
> **Consequence for the ready packet:** `compilerSHA256` is provenance about the *machine*, not
> identity of the *artifact*, and it is the only field in the packet that is platform-dependent
> (content and interface are sha256 over bytes; the tarball was measured to reproduce
> byte-identically across toolchains). It is asserted against the platform table and deliberately
> kept **out** of the byte-compared golden — including it would make the golden pass on exactly
> one platform, i.e. a gate that cannot run in CI.
>
> **S3 answer — "why is `host/pkgproj` not a package?"** It is HOST code, not kernel: it adds
> nothing to `world/**` and introduces no new semantics, sitting at the outward projection and
> verification boundary in the same class as `host/replay`. It computes the hashes that
> **authorize** package publication, so it cannot itself be a published package without
> circularity.

The Go gate covers broker changes through `go test ./...`, but acceptance requires named tests for
allowed success, namespace rejection, missing approval, changed hash, missing credential, definite
failure, ambiguous timeout, crash-window recovery, matching-409 reconciliation, mismatching-409
conflict, absence resolution, replay non-dispatch, and secret redaction.

## Milestones

### SM.A — deterministic package projection and non-vacuous dry-run gate (~1 day)

- Land `packages/world-core/**`, generator, exact drift/export/tar-content gate, and smoke.
- Extend `verify_ail.sh` to invoke the package verifier with the pinned binary.
- Produce the deterministic ready-packet schema and golden.
- Owns AC1, AC2, AC3, AC4, AC5, AC6.
- Independently landable and CI-green; no Open Decision blocks it.

### SM.B — brokered publish handler, attended stamp, and de-ambient credential (~2–2.5 days)

- Land `Registry.Publish`, one-shot scope/cost law, approval validation, minimal child environment,
  credential provider, typed indeterminate outcome, and all named handler tests.
- Land the atomic durable approval-claim/effect-intent transaction and ratify its schema/DDL fixture.
- Preserve existing handler failure semantics for every non-indeterminate error.
- Add the clause-2 dependency allowlist.
- Owns AC7, AC8, AC9, AC10, AC11, AC12.
- Independently landable and CI-green against a local fake validator; no live write.
- No Open Decision blocks it.

### SM.C — reconciliation, replay evidence, and attended operational runbook (~0.75–1 day)

- Land read-only exact-metadata reconciliation and three-state receipt transitions.
- Land clean-room exact-version/hash verification fixture against a local immutable registry fake.
- Land the attended runbook that stops at readiness by default.
- Route the namespace-auth ask upstream without claiming it is implemented.
- Owns AC13, AC14, AC15, AC16, AC17.
- Independently landable and CI-green; no Open Decision blocks it.

### SM.D — first public publish and hash-verified consumption (~0.25 day, attended)

- Human reviews the ready packet and mints the exact one-shot approval.
- Invoke the live broker effect once; reconcile if and only if indeterminate.
- Verify the public metadata and a clean-room exact-version install against the receipt.
- Owns AC18, AC19, AC20.
- **Blocked on `8/OD-1`**; never run headless or in CI.

Total implementation estimate: 4–5 days. SM.A–SM.C may land without performing a public write.

## Files to Create/Modify

| Path | Change |
|---|---|
| `packages/world-core/ailang.toml` | New exact package manifest |
| `packages/world-core/_smoke.ail` | New four-module smoke |
| `packages/world-core/world/{types,contracts,transitions,logepoch}.ail` | New byte-identical publish projection |
| `scripts/build_world_package.sh` | New deterministic allowlisted projection builder |
| `scripts/verify_world_package.sh` | New manifest/drift/smoke/tar/hash/dry-run gate |
| `scripts/verify_ail.sh` | Invoke the new package gate explicitly |
| `host/broker/registry_publish.go` | New live publish and read-only reconciliation handlers |
| `host/broker/registry_publish_test.go` | New contract, ambiguity, reconciliation, and redaction tests |
| `host/broker/approve.go` | Reuse existing request/decision wires; validate publish scope and route decision refs into single-use claim |
| `host/broker/approve_test.go` | Prove landed approval request/decision/poll compatibility with publish-bound decisions |
| `host/broker/broker.go` | Route publish through atomic approval-claim/effect-intent append; preserve typed indeterminate result without appending outcome |
| `host/broker/record.go` | No wire change expected; helper/schema constants only if required |
| `host/broker/recover.go` | Surface publish findings to explicit reconciler; never dispatch |
| `host/broker/dependency_test.go` | Extend core/handler dependency allowlists |
| `host/store/journal.go` | Add the claim-plus-effect-intent transaction beside the existing receipt API without changing intent/outcome wires |
| `host/store/store.go` | Add `ErrApprovalAlreadyConsumed`; advance schema version 1→2 and its version-handling path |
| `host/store/schema.sql` | Add durable `approval_claims` table with unique approval and invocation refs |
| `host/store/schema_version_test.go` | Update the landed independent exact-DDL/version fixture and ledger from version 1 to 2 |
| `host/store/journal_test.go` | Test atomic claim-plus-intent rollback, durability, and uniqueness |
| `docs/SELF_MOD_PUBLISH.md` | Attended ready/approve/invoke/reconcile/install runbook |
| `.github/workflows/ci.yml` | No live network step; only invoke repository gates if not already reached through scripts |

## Conflict Surface

| Landed behavior / path | Collision | Resolution |
|---|---|---|
| Broker decision/debit ledger, `host/broker/broker.go`, `decide.go` | One-shot publish must debit exactly once; special ambiguity must not refund or double-dispatch | Cost `1`, budget `1`; serialize through existing Session; retry prohibited while receipt pending |
| EffectRecord fixed codec, `host/broker/record.go` | Adding publish fields directly would change golden bytes and replay contract | Bind publish payload through existing content-addressed request/result refs; no record-wire widening |
| Ordinary handler failures, `host/broker/broker.go`, `handlers.go` | A broad special case could leave normal failures pending | Only typed `IndeterminateEffectError` suppresses outcome; negative-control test proves ordinary timeout remains resolved failed |
| Generic recovery, `host/broker/recover.go` | Recovery intentionally never dispatches or auto-resolves | Keep it fail-closed; explicit read-only reconciler handles only `Registry.Publish` |
| Landed human approval, `host/broker/approve.go` | Existing decisions are immutable but infinitely reusable; its linked registry head is also a blind last-writer-wins update | Reuse `HumanHandler`, both human effects, `DecideApproval`, and both V1 wires unchanged; bind publish to the decision ref and consume it through the atomic store claim, never through the head as a lock; queue the pre-existing head race separately |
| Append-only journal and receipt law, `host/store/journal.go` | Rewriting/deleting intent or outcome would violate frozen history | Append one outcome after probe; never mutate rows; retain three states |
| Store transaction API, `host/store/store.go` | No claim-if-unused primitive exists, so separate claim and intent writes permit restart/concurrency reuse | Add `AppendClaimedEffectIntent` and `ErrApprovalAlreadyConsumed`; one transaction inserts the unique claim tuple and the ordinary effect intent, with all-or-nothing visibility |
| Store schema, `host/store/schema.sql` | Durable single-use consumption requires state not represented by immutable approval objects or the journal | Add `approval_claims(approval_ref PRIMARY KEY, request_ref NOT NULL, invocation_id NOT NULL UNIQUE)`; no update/delete/refund path |
| Landed DDL gate, `design_docs/implemented/w-ddl-gate-teeth.md` / `host/store/schema_version_test.go` | Its independent embedded/compile-time fixture intentionally reds on every `schema.sql` edit | Ratify the added table in the same SM.B change by updating the independent exact-DDL/version fixture and its table count; retain the RED-on-unratified-edit control |
| `scripts/verify_ail.sh` | Current roots do not semantically assert package projection identity | Explicitly invoke new exact package gate; keep canonical required identities and 14 tests unchanged |
| `scripts/verify_go.sh` | Replay silently skips without pinned `AILANG_BIN`; broker tests need no live registry | Retain pin/anti-skip and use local fake validator only |
| `.github/workflows/ci.yml` | CI `ailang-verify` currently installs `latest`; live credential must never enter CI | Package gate uses exported pinned `AILANG_BIN`; no secrets/live publish; item 9 owns broader latest→pin correction |
| Replay goldens, `host/replay/**` | Module renames or source-byte changes alter v0.30.0-scoped fixtures/goldens | Keep module paths unchanged; projection copies do not replace fixture imports; no golden re-record |
| Canonical modules, `world/*.ail` | A second editable copy can drift; package prefix rule could tempt rename | Canonical files win; byte-identity gate; use manifest `module_prefix = "world"` |
| Coding standards S2/S3, `design_docs/coding-standards.md` | Network in kernel or presenting packaged kernel as new behavior violates boundaries | Handler-only network; artifact is a distribution projection; future behavior stays package-first |
| Thesis local-first and self-mod boundaries, `design_docs/DESIGN.md` §1/§14 | Runtime registry dependency or compiler modification would cross hard boundary | Registry only at outward handler/install fixture; no compiler or kernel-runtime dependency |
| Package immutability | `0.1.0` cannot be replaced or publisher-unpublished | Exact attended preflight; reconciliation; fixes roll forward to a new version |
| Public multi-vendor assumptions | World becomes first observed second vendor | Label convention honestly; verify public metadata; route upstream auth ask |

## Systemic-Issue Audit

This item exposes one systemic security shape: a shared secret in an ambient shell plus a client that
does not require it locally makes irreversible authority look like an ordinary command. The registry
also validates name shape while deferring name authority, so a green publish preflight can be
mistaken for namespace authorization. The durable correction is a separation of proofs:

- syntactic name validity;
- World-local vendor policy;
- attended authorization of exact bytes;
- registry authentication with a shared secret;
- registry-level namespace authorization (**absent**).

No one proof may be reported as another. The receipt records all five, including the absent last
one. This applies beyond packages: every outward handler should distinguish transport
authentication, resource authorization, human authorization, and artifact integrity rather than
compress them into “authorized”.

The gate-vacuity audit finds three likely false greens and closes each: a copy step that copied zero
files, a dry-run that never checked full hashes, and a 409 treated as success without comparing
metadata. Each has a named red mutation below.

## Deferred Scope

- Implementing registry-side vendor registration, per-vendor principals/scopes, or signing in the
  upstream AILANG repository. This design routes the requirement; World does not patch the compiler.
- Publishing any second World package or any version after `world/core@0.1.0`.
- Turning the packaged semantic core into a runtime dependency of World.
- Package cascade updates, automatic version selection, `latest`, or roll-forward automation.
- OS-principal, enclave, keychain, or separate-process secret isolation.
- Public unpublish. Publisher unpublish is unavailable and never part of recovery.
- Renaming canonical `world/*` modules into a package-name path.
- Changing AILANG's tarball include rules or extracting its hashing API upstream.

## Acceptance Criteria

- **AC1 (SM.A)** The exact four canonical modules have package projections, and deleting any one
  makes `scripts/verify_world_package.sh` fail with its path named.
- **AC2 (SM.A)** Changing one byte in any projection makes the gate fail with canonical/projection
  SHA-256 values; unchanged copies pass.
- **AC3 (SM.A)** A manifest named `sunholo/core`, `someoneelse/core`, or `world/other` fails the
  World-local policy, while exact `world/core` passes; the v0.30.0 generic manifest parser is also
  shown accepting at least one non-world vendor control so the local gate carries the difference.
- **AC4 (SM.A)** Adding a fifth package `.ail` or a design sketch path makes the tar-entry allowlist
  fail; the unmodified tar lists only manifest, smoke, and four modules.
- **AC5 (SM.A)** `_smoke.ail` imports all four exports and runs under the pinned binary; deleting one
  import makes the smoke-coverage manifest fail even if compilation stays green.
- **AC6 (SM.A)** Dry-run runs with the credential unset, emits exact package/version/exports/effects,
  and its recomputed full tarball/content/interface hashes equal the ready packet.
- **AC7 (SM.B)** A request with no exact `Registry.Publish` grant, wrong scope, expired logical time,
  or budget zero is denied before handler dispatch, with a persisted denial record.
- **AC8 (SM.B)** A valid grant with budget one atomically persists the approval claim and effect
  intent before dispatch and leaves budget zero; after constructing a fresh session with a fresh
  budget, the same approval still fails `ErrApprovalAlreadyConsumed` before credential load and
  POST, proving this is not satisfied by in-memory budget alone.
- **AC9 (SM.B)** Missing, expired, already-consumed, wrong-scope, wrong-hash, denied, or malformed
  landed `ApprovalDecisionV1`/`ApprovalRequestV1` fails before credential load and POST; an approved
  decision traversed through `DecideApproval` and `EffectHumanPollApproval` is the positive control.
- **AC9a (SM.B)** After one invocation consumes an approval, close and reopen the store and create
  a new process/session/budget; reuse returns `ErrApprovalAlreadyConsumed`, and the shared fake's
  total POST count remains exactly one.
- **AC9b (SM.B)** Two concurrent sessions with distinct fresh budgets race the same approval behind
  a start barrier; exactly one obtains the durable claim and dispatches exactly one POST, while the
  other returns `ErrApprovalAlreadyConsumed` before credential load.
- **AC9c (SM.B)** A first invocation dispatches once and returns typed indeterminate; after store
  reopen, a fresh session/budget retries with the same approval and receives
  `ErrApprovalAlreadyConsumed`, not budget denial, while the total POST count remains exactly one
  and recovery performs only read-only reconciliation.
- **AC10 (SM.B)** Production startup fails when `AILANG_REGISTRY_API_KEY` is ambient; dry-run and all
  non-publish subprocesses observe it unset; logs/objects/errors contain no sentinel secret.
- **AC11 (SM.B)** A definite handler failure appends a resolved failed outcome, while a typed
  ambiguous transport result leaves the same durable intent indeterminate; the two arms differ.
- **AC12 (SM.B)** Beyond landed `TestBrokerDependencyAllowlist` and its null-case control, a new
  World-boundary gate enumerates `world/**`, `host/store/**`, `host/replay/**`, and
  `cmd/ailang-worldd/**`; inserting a compiling registry/HTTP/cloud import into each protected Go
  group, or a registry/package-cache lookup marker into each protected AILANG group, fails with the
  exact path, while the unmodified enumeration is non-empty and network code confined to
  `host/broker` passes.
- **AC13 (SM.C)** Recovery reports an indeterminate publish without dispatching any handler.
- **AC14 (SM.C)** Exact public metadata resolves an indeterminate receipt to
  `succeeded-reconciled`, binding all three hashes.
- **AC15 (SM.C)** A 409 with mismatching metadata resolves to conflict, never success; matching
  metadata is the negative control and resolves success.
- **AC16 (SM.C)** Bounded repeated absence resolves `not-published` without a POST; a later live
  retry still requires a new attended approval/grant.
- **AC16a (SM.C, added under the round-2 carve-out)** An absence sample whose same-pass
  known-positive control does NOT return `200` with well-formed JSON is recorded **UNINFORMATIVE**,
  does not decrement the bounded window, and can never contribute to a `not-published` resolution.
  Three arms, all required and all distinguishable in the receipt: (i) control `200` + target `404`
  → counts as absent; (ii) control `200` + target `200` → resolves success/conflict by hash; (iii)
  control non-`200` (fixture: unreachable host, wrong `$AILANG_REGISTRY`, and a 403) + target `404`
  → **`probe-unavailable`**, human-required, and the assertion is that the receipt does NOT read
  `not-published`. Arm (iii) is the one that can fail today; arms (i)/(ii) are its controls.
- **AC16b (SM.C, added under the round-2 carve-out)** The reconciler treats a non-JSON error body as
  an error path, not a parse crash: a fixture serving the measured GCS XML
  `<Code>NoSuchKey</Code>` (V-N arm 2) yields a clean `absent` classification with the control
  firing, and a truncated/garbage body yields `probe-unavailable`, never `absent`.
- **AC17 (SM.C)** Replay returns the recorded publish result and performs zero network calls.
- **AC18 (SM.D)** The attended approval object hash equals the approval ref in the durable request,
  and the effect scope/package/version/hashes equal the ready packet byte-for-byte.
- **AC19 (SM.D)** The public registry metadata for `world/core@0.1.0` equals the resolved receipt on
  vendor, name, version, tarball SHA-256, content hash, and interface hash.
- **AC20 (SM.D)** A clean-room consumer declares parser-valid exact `world/core = "0.1.0"` and
  installs with the pinned binary, but compile/test is instrumented to be unreachable until the
  installed tarball/content/interface hashes equal the receipt; changing one receipt hash fails at
  that comparison while the exact-version manifest still parses, proving added World behavior
  beyond the landed parser's existing `latest`/range rejection.

## Non-Vacuity — named RED mutation for every added gate

Before trusting any mutation arm: take a backup first; assert the target anchor count is exactly
one (or the enumerated file count is exact); record baseline SHA-256; apply the mutation; assert the
post-mutation SHA-256 differs; print the diff; confirm the mutant compiles when compilation is not
the intended gate; then run the named gate. Restore from the backup and assert byte-identical
SHA-256. If the mutation did not land, its result is void.

| Gate | RED mutation and predicted failure |
|---|---|
| Projection presence | **MUT-SM-PROJECTION-MISSING**: remove staged `world/contracts.ail` after proving 4→3 allowlisted module count; named missing-file failure |
| Projection byte identity | **MUT-SM-PROJECTION-DRIFT**: change one existing comment byte after anchor-count=1 and differing SHA; drift failure names both hashes |
| Exact vendor/name | **MUT-SM-VENDOR**: replace exactly one `world/core` with `sunholo/core`; generic v0.30.0 dry-run remains green but World gate reds |
| Exact export set | **MUT-SM-EXPORT-DROP**: remove exactly `world/logepoch`; manifest still parses, export-set gate reds |
| Empty effects | **MUT-SM-EFFECT-ADD**: add one effect after proving empty list anchor; policy reds |
| Tar allowlist | **MUT-SM-TAR-LEAK**: add one valid compiling `leak.ail`; tar creation succeeds, allowlist reds on exact unexpected entry |
| Smoke coverage | **MUT-SM-SMOKE-IMPORT-DROP**: remove exactly one module import while leaving smoke compilable; coverage identity gate reds |
| Full-hash readiness | **MUT-SM-HASH-FLIP**: flip one hex nibble in ready packet after exact anchor assertion; comparison reds before approval |
| Pinned compiler | **MUT-SM-BINARY-DIRTY**: substitute a fake version output containing `-dirty`; binary identity gate reds before dry-run |
| Broker capability | **MUT-SM-GRANT-SCOPE**: change exact version scope to `0.1.1`; denial record appears and dispatch counter remains zero |
| One-shot budget | **MUT-SM-COST-ZERO**: change request cost 1→0 after anchor assertion; contract test reds because irreversible attempt would be free |
| Human approval | **MUT-SM-APPROVAL-HASH**: flip one approved tarball digest nibble; handler rejects before credential-provider and POST counters increment |
| Atomic approval claim | **MUT-SM-CLAIM-NONATOMIC**: split the claim insert and effect-intent append into independently committed transactions (equivalently remove the unique expected-unused insert); AC9b reds because both racing sessions can pass the claim boundary or the claim/intent atomicity assertion fails |
| Durable approval claim | **MUT-SM-CLAIM-MEMORY**: replace the table insert with an in-memory used-ref map; AC9a reds after store/process reopen when the second session reaches dispatch |
| Burn on crash | **MUT-SM-CLAIM-REFUND**: delete/refund the claim on the typed-indeterminate/crash arm; AC9c reds when the fresh session reuses it or the POST counter reaches two |
| Ambient-secret startup | **MUT-SM-AMBIENT-KEY**: launch with sentinel ambient key; production constructor must red; negative control with variable absent succeeds |
| Secret redaction | **MUT-SM-SECRET-ERROR**: fake validator echoes sentinel key in error body; sanitizer test reds if sentinel reaches error/log/object |
| Typed ambiguity | **MUT-SM-INDETERMINATE-AS-FAILED**: map typed ambiguous error through ordinary failure arm; receipt-state assertion reds (`resolved` vs `indeterminate`) |
| Ordinary failure preservation | **MUT-SM-ALL-ERRORS-PENDING**: classify ordinary exit error as indeterminate; negative-control test reds because landed resolved-failed behavior changed |
| Core dependency allowlist | **MUT-SM-CORE-HTTP**: insert a compiling `net/http` import/use in a protected Go core test fixture; allowlist reds specifically on import |
| Recovery non-dispatch | **MUT-SM-RECOVERY-DISPATCH**: call fake publish handler from recovery; dispatch counter reds at 1 |
| Matching reconciliation | **MUT-SM-RECON-HASH**: flip metadata interface hash; success expectation reds and receipt resolves conflict |
| 409 discrimination | **MUT-SM-409-AUTO-SUCCESS**: short-circuit 409 to success; mismatching-metadata arm reds while matching arm remains success |
| Absence bound | **MUT-SM-ABSENCE-UNBOUNDED**: remove sample bound under a fake always-404 registry; wall-clock/process-group bound reds |
| Probe instrument validity (AC16a) | **MUT-SM-PROBE-NO-CONTROL**: after asserting the control-fetch call site appears exactly once and the post-edit SHA-256 differs, delete it so absence is believed on the target 404 alone. AC16a arm (iii) must red — the receipt resolves `not-published` under a broken probe. **Negative control, mandatory (rule 3d): run arm (iii) UNPATCHED in the same session and require `probe-unavailable`.** Same outcome in both arms means the fixture's control was never wired and the mutation proved nothing |
| Probe instrument validity, silent-degrade variant (AC16a) | **MUT-SM-PROBE-CONTROL-ALWAYS-OK**: leave the control call site in place but hardcode its verdict to pass. This is the more realistic regression — the code still *looks* controlled. Arm (iii) must red identically to `MUT-SM-PROBE-NO-CONTROL`; if it does not, the control is decorative |
| Error-body handling (AC16b) | **MUT-SM-XML-AS-ABSENT**: make the classifier treat ANY non-200 as `absent` after asserting the branch anchor count is exactly one; the garbage-body arm of AC16b must red while the `NoSuchKey`-XML arm stays green — the two arms differing is the whole point, and a mutation that reds both means the fixture cannot tell them apart |
| Replay non-dispatch | **MUT-SM-REPLAY-NETWORK**: register a counting metadata handler and consult it during replay; zero-call assertion reds |
| Version pin | **MUT-SM-LATEST**: replace exact dependency once with `latest`; clean-room manifest gate reds before install |
| Hash-before-compile | **MUT-SM-HASH-GATE-BYPASS**: bypass the receipt-hash comparison while keeping exact `0.1.0`; AC20's mismatched-hash arm reds if compile/test becomes reachable |
| Receipt/public equality | **MUT-SM-PUBLIC-DIGEST**: serve valid metadata with one differing tarball nibble; AC19 comparator reds |

## Open Decisions

- **8/OD-1 — attended authorization for the irreversible first publish (OPEN, blocks SM.D only).**
  Does Mark approve the exact ready packet for `world/core@0.1.0` and mint the single-use stamp?
  Controller default: **do not publish**; remain `READY_AWAITING_HUMAN_PUBLISH`. Alternative: approve
  the exact packet and execute SM.D attended. Silence is not approval.
- **8/OD-2 — upstream namespace authorization disposition (OPEN, non-blocking).** Will the AILANG
  registry add vendor principals/registration before a later World release? Controller default:
  file/route the requirement and label `world/` convention-only. Alternative: if upstream lands and
  is measured, tighten the handler and receipt to server-enforced vendor authorization. This does
  not block `0.1.0` under Mark's ratified current-registry decision.

## Quorum verification log

Reviewers: `gpt5-6-sol`, `gemini-3-1-pro`, plus the controller's in-session verdict. Artifacts in
`.ailang/state/mission-quorum/`. Every round's outcome, including the degraded one, is recorded here
by name — a reviewer that did not run is never silently a pass.

**Round 1 — `w-self-mod-vertical-2026-08-05T04-05-25Z.json` — BLOCKED, 2/2 reject, `$0.102967`.**
Both reviewers converged on ONE defect, from opposite ends: `gemini-3-1-pro` — *"You successfully
identified the existing `approve.go` machinery in your PV log, but dropped the finding completely;
it never made it into the actual design specs, file lists, or conflict analysis."* `gpt5-6-sol` —
the same gap stated as its consequence: *"An indeterminate attempt could therefore be followed by a
fresh session reusing the same approval and dispatching a duplicate immutable POST."* **The
provenance of that defect is worth recording, because it is not the designer's:** the `approve.go`
row was inserted into the Verification Log by the CONTROLLER *after* the body was written, so the
contradiction between evidence and design was planted by construction. A post-hoc evidence row that
the body has never seen is not corroboration — it is a second author writing into the same
document. Both objections were ACCEPTED (each carried a concrete reviewer-authored `proposed_fix`;
neither disputed the design direction) and routed to the designer with the reviewers' verbatim text
plus controller measurements A1–A6 — in particular that `SetRegistryHead`
(`host/store/store.go:601`) is a blind upsert with no expected-previous and no transaction, so the
reviewer's *"if the current store lacks atomic claim-if-unused support, explicitly add the required
store API/schema"* branch fires and "no schema change" was not available.

**Round 2 — `w-self-mod-vertical-2026-08-05T04-12-48Z.json` — BLOCKED, `$0.03663`, DEGRADED to N−1:
`gpt5-6-sol` ABSENT (reason recorded by the tool: `budget`).** The revision satisfied both round-1
objections (`approval_claims` + atomic `AppendClaimedEffectIntent`, AC9a/b/c, four new mutations,
estimate honestly moved 3–4 → 4–5 d). `gemini-3-1-pro` then found a NEW and better objection on the
recovery path — quoted in full in the revised Probe-then-resolve section.

**Round 3 — CONTROLLER CARVE-OUT, no third reviewer round.** The one remaining blocking objection
(a) carried a concrete reviewer-authored `proposed_fix` and (b) disputed no design direction — it
asked for a measurement and for the design to follow wherever that measurement led. The controller
ran the reviewer's own prescribed check (row **V-N**) rather than arguing, and applied the fix the
reviewer specified. **The measurement's verdict was neither "reviewer right" nor "reviewer wrong",
which is why it was worth running:** the path *string* `packages/{vendor}/{name}/{version}/metadata.json`
is correct, but it is a **GCS bucket object key, not a validator HTTP route** — the validator serves
8 routes and none of them is `/packages/…`. The reviewer's stated failure mode (false 404 →
`not-published` → re-authorized irreversible POST) was therefore reachable for a different reason
than the one suspected. Landed under the carve-out: the rewritten Probe-then-resolve, **V-N** with
both live arms, **AC16a**/**AC16b**, and three new RED mutations
(`MUT-SM-PROBE-NO-CONTROL`, `MUT-SM-PROBE-CONTROL-ALWAYS-OK`, `MUT-SM-XML-AS-ABSENT`), each carrying
its own negative control per rule 3d. Also closed by measurement in the same pass, rather than left
as an implementer TODO: `edition` (validated for non-emptiness only → `"1"`, the toolchain's own
template value).

**Standing caveat, recorded rather than resolved.** Round 2 ran one reviewer. The carve-out is a
controller judgement applied to a single reviewer's objection, and `gpt5-6-sol` has not seen the
round-2 or round-3 text. Two questions the controller raised in its own round-2 note went unanswered
by any reviewer and are carried into the sprint rather than treated as cleared: **(1)** whether a
`schema.sql` change is acceptable inside this item given that the landed `w-ddl-gate-teeth` DDL gate
reds on any schema edit by design, and **(2)** whether 4–5 days is still one queue item or wants
splitting at the SM.B boundary. Neither blocks writing the doc; both are the first things a sprint
planner should price.

## Provenance

Mark's attended 2026-08-04 queue-row decision at commit `de80792` is ratified mission state: public
registry, `world/` vendor string, brokered/receipted publish, local-first consumption, exact version
and hash verification. This design resolves the newly measured fact that the registry offers no
vendor claim operation rather than reopening that decision.

Controller-supplied measurements V-A through V-M are first-party evidence from iteration 51 on
2026-08-05. Every row derived from them carries the required marker below. Independent repository
measurements were made in this worktree at base `dd2c173`; upstream code was read only from pinned
commit `e37b370d1d7a9c4e7136b319e38bec4d5f2bd9a0`.

## Premise Verification Log

**World base measured:** `dd2c173` (controller-stated HEAD for V-M; current worktree content used
for independent path reads). **AILANG upstream base measured:**
`e37b370d1d7a9c4e7136b319e38bec4d5f2bd9a0`. **Oldest admitted measurement:** Mark's namespace
census on 2026-08-04, refreshed by controller census on 2026-08-05.

| Premise | Verified | Evidence (command + observed output) |
|---|---|---|
| V-A pinned binary identity | Yes | `(controller iter-51, 2026-08-05)` `/tmp/ailang-v0300/ailang --version` → `AILANG v0.30.0`, commit `e37b370`; `shasum -a 256` → `e9746fef8570bc42b8cc52c0e88b7088468a5d2bd38bb8c42e27e5859b8f3fb5`; PATH binary observed dirty v0.33 dev build |
| V-B publish CLI and endpoints | Yes | `(controller iter-51, 2026-08-05)` `/tmp/ailang-v0300/ailang publish --help` → flags `--dry-run`, `--allow-dotted-tool-names`; `cmd/ailang/pkg_info.go:18-26` → validator env precedence `AILANG_REGISTRY_VALIDATOR`, `AILANG_REGISTRY_API`, default |
| V-C shared-secret auth only | Yes | `(controller iter-51, 2026-08-05)` pinned `cmd/ailang/pkg_publish.go` read → key present sets `X-API-Key`, absent omits header; validator `main.go:54,106-113` → one `REGISTRY_API_KEY`, absent/mismatch HTTP 403 |
| V-D no vendor registration | Yes | `(controller iter-51, 2026-08-05)` pinned validator `main.go:161-177` → vendor/name shape check and verbatim `Namespace auth — deferred (accept all publishers for now)` |
| V-E vendor dry-run discrimination | Yes | `(controller iter-51, 2026-08-05)` four arms under `env -u AILANG_REGISTRY_API_KEY`: `world/probe rc=0`, `someoneelse/probe rc=0`, `sunholo/probe rc=0`, known-positive invalid `novendor rc=1` with `[package].name must be vendor/name format` |
| V-F public census | Yes | `(controller iter-51, 2026-08-05)` read-only GET index → schema `ailang.registry/v1`, 40 packages, histogram `sunholo 40`, `world/` count 0, positive control `sunholo/` count 40 |
| V-G immutable/unrecallable publish | Yes | `(controller iter-51, 2026-08-05)` validator metadata stat → existing version HTTP 409; `unpublish.go:159-160` → server API-key configuration required, operator-only |
| V-H ambient authority | Yes | `(controller iter-51, 2026-08-05)` presence-only probe `[ -n "$AILANG_REGISTRY_API_KEY" ] && echo SET || echo UNSET` → `SET`; `HOME` positive control → `SET`; value never printed |
| V-I publish gates/order | Yes | `(controller iter-51, 2026-08-05)` pinned `pkg_publish.go` → manifest, constraint warning, path rewrite/restore, assets, smoke, tarball, hashes, upload; extension-without-smoke hard-fails |
| V-J dry-run observables | Yes | `(controller iter-51, 2026-08-05)` `world/probe@0.1.0` dry-run → tarball SHA prefix, content-hash prefix, interface-hash prefix, exports `[world/probe/core]`, effects `[]` |
| V-K init ignores requested name | Yes | `(controller iter-51, 2026-08-05)` v0.30.0 init in `.w143-publish-probe` with arg `worldprobe` → `name = "local/.w143-publish-probe"` |
| V-L World has no manifest | Yes | `(controller iter-51, 2026-08-05)` `find . -name ailang.toml -not -path './.git/*'` → empty; same-call positive control `find . -name '*.ail'` → `world/logepoch.ail`, `world/types.ail`, others |
| V-M broker/store machinery | Yes | `(controller iter-51, 2026-08-05)` at `dd2c173`: `host/broker/broker.go` Session/NewSession/Invoke/decide/putRecord/replay; `decide.go` Capability/EffectRequest; `record.go` EffectRecord; `journal.go` GetReceipt/GetEffectReceipt and three states |
| **V-N — the metadata probe endpoint: a BUCKET OBJECT, not a validator route; measured with both controls** | Yes — **added under the round-2 carve-out, running reviewer `gemini-3-1-pro`'s own prescribed measurement** | `(controller iter-51, 2026-08-05)` **Router:** `git show e37b370…:cmd/registry-validator/main.go \| grep -nE 'HandleFunc'` → `:58` `/publish`, `:59` `/unpublish`, `:60` `/rebuild-index`, `:61` `/health`, `:62` `/version`, `:65` `/api/packages`, `:66` `/api/packages/`, `:67` `/api/stats` — **8 routes, none serving `/packages/…`** (this is a complete enumeration, not a `head`-truncated one). **Key construction:** `git grep -n 'metadata.json' e37b370… -- cmd/ internal/` → `main.go:168`, `cache.go:159`, `unpublish.go:190`, all `fmt.Sprintf("packages/%s/%s/%s/metadata.json", …)` passed to `v.bucket.Object(...)` — a GCS key, server-side. **Live, read-only, two arms:** ARM 1 known-positive control `GET https://storage.googleapis.com/ailang-registry/packages/sunholo/auth/0.4.1/metadata.json` → **HTTP 200, 1289 B**, JSON keys `content_hash`, `interface_hash`, `manifest`, `name`, `published_at`, `published_by`, `schema`, `tarball_hash`, `tarball_size_bytes`, `validation`, `version` — so all three hashes AC19 compares are actually served. ARM 2 negative control, the exact object this design will create, `…/packages/world/core/0.1.0/metadata.json` → **HTTP 404, 217 B**, body is **GCS XML** `<Code>NoSuchKey</Code>`, not JSON. Net: the reviewer's path *string* was right and its *nature* was wrong, and the nature is what the recovery boundary rests on — see the revised Probe-then-resolve |
| Module-prefix compatibility | Yes | `git -C /Users/voightkampff/dev/sunholo-data/ailang show e37b370d1d7a9c4e7136b319e38bec4d5f2bd9a0:internal/pkg/manifest.go \| sed -n '220,310p'` → export accepted when it equals/starts with package name **or** configured single-segment `module_prefix`; exact-version dependency validation also shown |
| Tarball breadth | Yes | `git -C /Users/voightkampff/dev/sunholo-data/ailang show e37b370d1d7a9c4e7136b319e38bec4d5f2bd9a0:internal/pkg/tarball.go` → recursive walk includes `ailang.toml`, every `*.ail`, `AGENT.md`, `_smoke.ail`, and `assets/`; excludes git/test dirs/lockfile; sorted zero timestamps |
| **V-O — `internal/pkg` extraction is impossible across modules (`DD-1`)** | Yes | `(controller iter-53, 2026-08-05)` `go.mod:1` here → `module github.com/sunholo-data/ailang-world`; `git show e37b370:go.mod` → `module github.com/sunholo-data/ailang`, `go 1.26.5` (World declares `go 1.25.6`). `git ls-tree --name-only e37b370 internal/pkg/` → 31 entries incl. `hasher.go`, `tarball.go`, `manifest.go`; control `cmd/ailang/` → 207 entries, so the instrument reads the tree. `git show e37b370:cmd/ailang/pkg_publish.go \| grep -nE 'os\.WriteFile\|os\.Create\|\.tar\.gz'` → `:79`/`:240` (toml rewrite + restore) and `:251` (multipart form, **upload** path only); control `grep -c 'Printf'` → 25. Tarball bytes are never persisted; hashes only ever printed at `[:24]` |
| **V-P — the three-arm cross-check AGREES across toolchains, and is proven able to RED (`AC6`)** | Yes — **the finding SM.A existed to surface** | `(controller iter-53, 2026-08-05)` live `CrossCheck` against `packages/world-core` with the pinned binary: `content`, `interface`, `tarball` all AGREE, tarball length `5472 = 5472`, so World's `go1.25.6` `compress/gzip`+`archive/tar` output reproduces the `go1.26.5`-built CLI's byte-for-byte. **Negative controls, each with sha256 proof the mutation applied and byte-identical restore** (`pkgproj.go` baseline `65efe4fb7e59…`): `MUT-SM-PKGPROJ-CONTENT-SEPARATOR` (`file:%s\n`→`file:%s`, sha `c258cdde…`) → RED on the **content** arm only; `MUT-SM-PKGPROJ-TAR-MODE` (`Mode: 0644`→`0600`, sha `5d13faad…`) → RED on the **tarball** arm only; each names both values; control green before and after |
| **V-Q — CI job 1 verifies `.ail` against an UNPINNED compiler, and it went active 2026-08-04 (`DD-7`, queue item 9)** | Yes — **contradicts item 9's recorded "latent, not active"** | `(controller iter-53, 2026-08-05)` `gh release view` → `latest` = **v0.33.0**, published `2026-08-04T12:25:38Z` (control: the `v0.30.0` tag still resolves as a distinct release). Step log for `af0c3b4`, run `30993399332`, SHA-addressed: job `ailang-code verify gate` → **`AILANG v0.33.0`**; job `go host build + test gate` → **`AILANG v0.30.0`** in the *same run*, which is the control proving the difference is real. Rig PATH `ailang` also measured **`v0.33.0-1-gdd68e0741`**, and it fails the package gate's own step 5 (*"5 properties never ran (no generator)"*) |
| **V-R — the pinned compiler's SHA-256 is platform-specific (`DD-7`)** | Yes | `(controller iter-53, 2026-08-05)` `file /tmp/ailang-v0300/ailang` → `Mach-O 64-bit executable arm64`, sha `e9746fef8570bc42…`. Downloaded `releases/download/v0.30.0/linux.x64.ailang.tar.gz` (39,825,888 B), published `.sha256` verified `OK` as the control, extracted → `ELF 64-bit LSB executable, x86-64`, sha `1e594d158dffa688…`. Non-vacuity of the pin: a byte-flipped copy still reporting `AILANG v0.30.0` (sha `74b475bc4715…`) is REJECTED naming both values |
| **V-S — `AC-B1.2`'s ledger gate was VACUOUS as delivered: a self-referential source-grep needle (`SM.B1`)** | Yes — **the finding SM.B1 existed to surface, and the executor's own mutation could not see it** | `(controller iter-54, 2026-08-05)` `TestSchemaVersionLedgerIsIndependent` reads its OWN source. Its two NEGATIVE needles were split (`"var schemaV2SQL = "+"schemaSQL"`) so they would not match the check-lines containing them; its POSITIVE needle was the single literal `"const schemaV2SQL = schemaV1SQL +"`, which **the check's own line satisfies** — so it passed regardless of the declaration. Measured, mutation confirmed landed by sha256 (`76d63695…`→`f277a80b…`) before the result was believed: `var schemaV2SQL = string(schemaSQL)` — the ledger becoming the very file it exists to attest — returned **`ok 0.290s`**; both negative needles were dodged by the `string(...)` conversion and the DDL comparison then compared `schema.sql` to itself. The executor's `MUT-SM-V2-LEDGER-DERIVED` redded ONLY because it used the bare `var schemaV2SQL = schemaSQL` form the negative needle was written to catch. **Repair**: line-anchored `regexp` (`(?m)^const schemaV2SQL = schemaV1SQL \+`), where an indented check-line cannot reach — the same instrument the charter's `^## STATUS 2026` rotation invariant uses — plus a semantic backstop (`schemaV2SQL == schemaSQL` is a hard failure). Post-repair BOTH forms RED, unmutated control green, restore byte-identical. The `sonnet` judge independently added a THIRD form, `const schemaV2SQL = schemaV1SQL + ""`, which also REDs |
| **V-T — `ErrApprovalAlreadyConsumed` is NOT returned on the concurrent-collision path (`AC9b`, SM.B2b)** | Yes — measured by the judge, not inferred | `(evaluator sonnet iter-54, 2026-08-05)` `AppendClaimedEffectIntent` guards with `SELECT EXISTS` then relies on the `approval_claims` PRIMARY KEY. Two callers can both pass the pre-check; the loser then fails on the constraint. Applying `MUT-SM-CLAIM-MEMORY` (removing the `SELECT EXISTS` check) surfaces the collision path directly: `reused approval error = *fmt.wrapError store: claim approval "sha256:…": constraint failed: UNIQUE constraint failed: approval_claims.approval_ref (1555), want ErrApprovalAlreadyConsumed`. **Correctness is preserved** (no double-consumption); only the error CLASSIFICATION differs. SM.B1 deliberately does not fix it — `AC9b` (two sessions racing one approval, under `-race`) is the criterion that must |
| **V-U — a 15.7 MB compiled binary was tracked in the repo, and five independent checks passed it** | Yes | `(controller iter-54, 2026-08-05)` `git ls-tree HEAD ailang-worldd` → `100755 blob`, `file` → `Mach-O 64-bit executable arm64`, 15,740,242 B, `A`-added by SM.A's squash `13315da`. It passed the codex executor, the `sonnet` evaluator (**87/100, zero blocking**), the controller's four-gate re-run and **both CI jobs**, because none enumerates tracked file *types*. Detector chosen: git's own classification, `git diff --numstat <empty-tree> HEAD` → `-`/`-` for binary blobs — **1 of 142** tracked files, against a control asserting the enumeration returned all 142 (a `git ls-files` count agrees at 142). `file(1)` was rejected as non-portable (`Mach-O` vs `ELF`). Gate landed in `scripts/verify_go.sh` (already run by CI job 2 at `ci.yml:121`, so no headless `ci.yml` edit), placed BEFORE the Go toolchain deny-list so it runs even where that check rejects the toolchain. Three arms: control `✓ 0 binary blobs among 141 tracked files`; `MUT-TRACKED-BINARY` → FATAL naming the path, green line absent; `MUT-DETECTOR-BLIND` (`binary_numstat=""`) → FATAL *"the instrument is broken, so every 'no binaries' result it reports is void"* (sha `05b5db9c…`→`8f050f9e…`, restore byte-identical). Confirmed running on the **ubuntu runner** in CI's step log at the same 141 count |
| Current module declarations/imports | Yes | `sed -n '1,220p' world/types.ail; sed -n '1,220p' world/contracts.ail; sed -n '1,260p' world/transitions.ail; sed -n '1,220p' world/logepoch.ail` → declarations `world/types`, `world/contracts`, `world/transitions`, `world/logepoch`; cross-imports use same prefixes; positive presence in every file |
| Replay path collision count | Yes — **number CORRECTED by controller; the instrument's scope is itself load-bearing** | `(controller iter-51, 2026-08-05, re-derived per rule 3b(v-b))` the row's own stated command returns **`18`**, not the `19` first written — a transcribed quantity that did not reproduce. Per-file enumeration in the same call: `world/transitions.ail` 4, `world/contracts.ail` 3, `world/types.ail` 3, `world/logepoch.ail` 1, `host/replay/testdata/transition_fixture.ail` 4, `host/store/store.go` 1, `design_docs/verification/w-m1-ailang-hardening/fixtures/verified_baseline.ail.txt` 2 (= 18). The equivalent `grep -rE … --include='…'` — which does **not** honour `.gitignore` — returns **26 lines across 12 files**, because it also walks the gitignored build cache `world/.ailang/cache/compile/**`. Neither is "the" answer: `18` is the tracked-source surface this design sizes against; the cache is regenerable and never tracked. Material to **AC4**: that cache holds **0** `*.ail` files (same-call control: **79** `*.json`), so it cannot leak through `CreateTarball`'s `*.ail` filter — but `CreateTarball` skips only `.git`, `tests`, `test` (`internal/pkg/tarball.go:43-47` @ `e37b370…`), **not** `.ailang`, so AC4's tar-entry allowlist is doing real work rather than restating an upstream exclusion |
| Broker intent-before-dispatch | Yes | `sed -n '1,360p' host/broker/broker.go` → allowed live flow `PutObject` → `AppendNextEffectIntent` → budget debit → `handler.Execute` → record/result → `AppendEffectOutcome` |
| Ordinary errors resolve failed | Yes | same command → any `handler.Execute` error creates `Allowed:true, Failed:true` EffectRecord then appends `Status:"failed"` outcome |
| **P6 is a real gap, not a re-invention — the typed-ambiguity path does NOT already exist on the dispatch side** | Yes — **controller independently checked for the opposite**, because "this capability is missing" is the claim class most likely to be wrong | `(controller iter-51, 2026-08-05)` `grep -n 'Indeterminate' host/broker/broker.go` → **rc=1, zero hits**, with a same-file known-positive control `grep -c 'putRecord' host/broker/broker.go` → **4**, so the instrument reads the file. `grep -rn 'IndeterminateEffectError' host/ --include='*.go'` (non-test) → **6 hits, ALL in `host/broker/recover.go`** (`:16,19,29,46,155,217`). So `IndeterminateEffectError` is landed but is produced **only by the post-hoc recovery scan**, never returnable by a handler at dispatch time. SM.B's typed-ambiguous return is therefore new work on a real seam, and `recover.go`'s existing `Recover`/`retryAllowed(indeterminate, reconciled)` (`:63,:81`) is the surface SM.C extends rather than replaces |
| **`host/broker/approve.go` and the dependency allowlist already exist — SM.B/AC12 EXTEND, they do not create** | Yes | `(controller iter-51, 2026-08-05)` `grep -nE '^func \|^type ' host/broker/approve.go` → `HumanHandler` (:37), `NewHumanHandler` (:41), `Execute` (:85), `DecideApproval` (:133), plus the approval request/decision/head wire types (:49-80). `host/broker/allowlist_test.go` → `TestBrokerDependencyAllowlist` (:85) **and** `TestBrokerDependencyAllowlistNullCase` (:104), i.e. the gate ships with its own anti-vacuity control. The sprint must bind the publish approval to these, and AC12 must state what it adds beyond the landed test |
| **`ailang.toml` already rejects `latest` and range specifiers at parse time — AC20 must claim more than that** | Yes | `(controller iter-51, 2026-08-05)` pinned `internal/pkg/manifest.go` (~`:298`) → `dep.Version == "latest" \|\| strings.ContainsAny(dep.Version, "^~><=")` → error `requires exact versions`. So AC20's "contains no `latest` or range" half is enforced upstream and would pass without World doing anything; its non-vacuous content is the **hash-verification-before-compile** half |
| A1 — landed human approval effects | Yes | `(controller iter-51, 2026-08-05)` `host/broker/approve.go:85` `HumanHandler.Execute` switches on `req.Effect`: `EffectHumanApprove` decodes `approvalInputWire`, builds `approvalRequestWire{Effect, Scope, Cost, Requester, Now}`, wraps `ApprovalRequestV1`, `store.PutObject`s it, and calls `appendApprovalHead` with a zero decision ref before returning `pendingBytes`; `EffectHumanPollApproval` parses `RequestRef`, returns pending when `findApprovalDecision` misses and `observedDecisionWire{Status:"decided", Decision:<raw>}` when found; the default rejects an unimplemented effect |
| A2 — approval decision is an operator entry point | Yes | `(controller iter-51, 2026-08-05)` `host/broker/approve.go:133` `DecideApproval` is documented verbatim: *"DecideApproval is an operator entry point, not an effect. It creates one immutable decision object and moves only the approvals registry head."* It accepts exactly `approve`/`deny`, requires an existing headed `ApprovalRequestV1`, writes `ApprovalDecisionV1` with `approvalDecisionWire{RequestRef, Decision, DecidedBy, Now}`, and appends the approval head |
| A3 — approvals registry is a linked head chain | Yes | `(controller iter-51, 2026-08-05)` `appendApprovalHead` at `approve.go:178` reads `GetRegistryHead(ApprovalsV1)`, stores `approvalHeadWire{RequestRef, PreviousHead, DecisionRef}`, then calls `SetRegistryHead`; `findApprovalRequest`, `findApprovalDecision`, and `walkApprovalHead` traverse that chain backwards |
| A4 — landed approval has neither consumption nor compare-and-set | Yes | `(controller iter-51, 2026-08-05)` an `ApprovalDecisionV1` remains forever reusable through `findApprovalDecision`; no burn/claim marker exists. `host/store/store.go:601` `SetRegistryHead` is one blind `INSERT ... ON CONFLICT(registry_name) DO UPDATE SET object_ref = excluded.object_ref`, with no expected previous value or enclosing transaction. The controller's non-test store-function search found exactly that one Set/Compare/CAS/Swap-shaped API. Thus concurrent approval-head appends can lose a head and the reviewer's conditional fires: this design adds the durable `approval_claims` schema and atomic `AppendClaimedEffectIntent` API; the lost-update defect in the landed approval head is pre-existing, surfaced here, and queued separately |
| A5 — landed broker allowlist already has anti-vacuity | Yes | `(controller iter-51, 2026-08-05)` `host/broker/allowlist_test.go:85` has `TestBrokerDependencyAllowlist`, `:104` has `TestBrokerDependencyAllowlistNullCase`, with `isStdlibImportPath` (`:32`), `disallowedDeps` (`:37`), and `goListDeps` (`:60`); AC12 therefore adds explicit World core/store/replay/daemon and AILANG marker coverage rather than re-claiming the landed broker gate |
| A6 — recovery ordering surface | Yes | `(controller iter-51, 2026-08-05)` `host/broker/recover.go` defines `IndeterminateEffectError` (`:19`), `IndeterminateEffect` (`:43`), `mayReportNotStarted(hasIntent bool)` (`:59`), `retryAllowed(indeterminate, reconciled bool)` (`:63`), `Recover` (`:81`), `recoverPending` (`:88`), `recoverCommitPending` (`:96`), and `recoverEffectPending` (`:173`); the combined claim/intent commits before credential load/POST, so recovery sees a burned approval with an indeterminate intent and `retryAllowed` remains false until reconciliation |
| EffectRecord binding surface | Yes | `sed -n '1,240p' host/broker/record.go` → fixed fields effect/scope/cost/budgets/allowed/failed/denial/requestRef/resultRef; canonical JSON codec and content-addressed objects |
| Receipt three-state implementation | Yes | `sed -n '1,180p' host/store/journal.go; sed -n '680,770p' host/store/journal.go` → constants `not-started`, `indeterminate`, `resolved`; `GetReceipt`/`GetEffectReceipt` derive state from intent/outcome presence |
| Recovery is non-dispatching | Yes | `rg -n 'indeterminate\|reconcile\|recovery' host/broker host/store --glob='*.go'` → `host/broker/recover.go` comments and tests state registries are never consulted and recovery does not dispatch/auto-resolve; known-positive handler registrations also enumerated |
| Current AILANG gate scope | Yes | `sed -n '1,240p' scripts/verify_ail.sh` → roots `design_docs` and `world`; required canonical identities 4 and named tests 14; package projection semantic assertions absent while positive canonical assertions are present |
| Current Go/replay gate scope | Yes | `sed -n '1,220p' scripts/verify_go.sh` → requires `AILANG_BIN` reporting v0.30.0, runs build/plain/race tests; `rg -n 'transition_fixture\|recorded_result\|golden\|v0.30.0\|AILANG_BIN' host/replay ... \| wc -l` → `18`, with enumeration showing fixture and golden use |
| CI boundary | Yes | `sed -n '1,180p' .github/workflows/ci.yml` → AILANG job installs `latest`; Go job pins v0.30.0 plus SHA verification; both invoke repository scripts; no publish step; known-positive CI jobs/steps enumerated |
| Local-first thesis and self-mod boundary | Yes | `rg -n '^## 14\|^### 14\|self-mod\|extension\|package' design_docs/DESIGN.md; sed -n '610,645p' design_docs/DESIGN.md` → §14 proposal→Verify→Commit extension lane, compiler/kernel/authority hard boundaries, version-pinned roll-forward; §3 says hosted registries are handlers, never core imports |
| Coding standards | Yes | `sed -n '1,220p' design_docs/coding-standards.md` → S2 effects at boundary, S3 slim package-first kernel, S5 pinned language verification, S6 non-vacuous gates, S7 usage docs |
| Landed broker contract | Yes | `sed -n '1,280p' design_docs/implemented/w-effect-broker-m3.md` → broker records every result, receives injected store, caller supplies time, zero-cloud handler boundary, bounded subprocesses, replay evidence contract |

## Related Documents

- `design_docs/world-mission.md` — ratified bar, clauses, and queue item 8
- `design_docs/DESIGN.md` §1, §3, §14 — semantic core, local-first boundary, controlled self-modification
- `design_docs/coding-standards.md` — S2, S3, S5, S6, S7
- `design_docs/implemented/w-effect-broker-m3.md` — landed capability/budget/record/replay contract
- `design_docs/implemented/w-world-library-m1.md` — canonical semantic library and replay
- `design_docs/implemented/w-m1-ailang-hardening.md` — non-vacuous AILANG gate
- `design_docs/implemented/w-race-gate-blindspot.md` — house design-doc structure and mutation discipline
