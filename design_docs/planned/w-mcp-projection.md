# w-mcp-projection — Session-Scoped A2A Projection and Agent Card (SPLIT parent; MCP dispatch carved out)

**Status**: Planned — **SPLIT APPLIED 2026-08-25 (iteration 125), THEN RE-QUORUMED AND BLOCKED
AGAIN AT ROUND 4 ON A SINGLE, DIFFERENT, NEWLY-LOCALISED SURFACE — SESSION AUTHORITY. A SECOND
SPLIT IS OWED BEFORE THIS DOC ROUTES TO A PLANNER; DO NOT SPRINT IT AS IT STANDS.** Round 4 is the
first round `gemini-3-1-pro` has **passed**, and `gpt5-6-sol`'s sole reject is confined to
**`P6.B-A2A`**: the repo has no inbound credential->session resolution at all — measured
first-party, `Bearer` **0** and session-lookup functions **0** across `host/`, with the
same-scope known-positive control at **181**. `D-WORLD-26` settled the credential *envelope*;
the *contents* question (who mints a session credential, where credential->(episode, grants)
lives, what expires it) was never built and is a **pre-existing gap at HEAD**, so it is filed as
its own queue row rather than absorbed here. **`P6.T`, `P6.D` and `P6.V` have drawn zero
objections in four rounds and have no session-authority dependence** — they are what the second
split leaves behind. Full evidence and disposition: `## Quorum verification log` -> `### Round 4`.

**Prior status text follows, and it remains accurate about round 3.** Quorum round 3 (2026-08-25,
both reviewers present, `absent_reviewers` empty) blocked the previous revision on two objections,
and the controller's routing call was SPLIT:

- **Objection 1 (`gpt5-6-sol`, DIRECTION-level, confirmed first-party): `serveapi/protocol` has no
  MCP JSON-RPC dispatch, and both routes around that are closed by this repo's own guardrails.**
  That objection — and with it the whole MCP half: the `/mcp/` handler,
  `initialize`/`tools/list`/`tools/call` dispatch, dispatch-bound envelope framing, SSE stream
  lifetime, and the MCP-specific acceptance criteria and mutations — is carried VERBATIM in the
  child doc [`w-mcp-dispatch-projection.md`](w-mcp-dispatch-projection.md), which is honestly
  BLOCKED on the upstream ask
  [`sunholo-data/ailang#885`](https://github.com/sunholo-data/ailang/issues/885) (re-measured this
  round: OPEN, 0 comments; control `#764` reads CLOSED with 6 comments — premise row N21).
- **Objection 2 (`gemini-3-1-pro`, narrow): an executable design must close its own forks.**
  CLOSED by the human: **D-WORLD-26 = ARM A** (Mark, attended, issue `#89`, 2026-08-25T19:06:41Z)
  — the session credential rides the standard `Authorization: Bearer` header. See Open Decisions
  and Decision 3; no open fork remains in this doc.

What this parent retains is executable end to end with nothing depending on MCP dispatch:
**P6.A** (done, record), **P6.T** (toolchain floor), **P6.D** (pinned dependency + one narrow
allowlist line), **P6.V** (verified commit-boundary law), and **P6.B-A2A** (the A2A half of the
old P6.B, re-cut). The existence proof that the A2A half needs no MCP SDK is upstream's own
`serveapi/a2a_handler.go`: **180 lines over stdlib + `protocol`, 0 SDK imports** (re-derived this
round via `gh api` at tag `v0.33.2` — premise row N23).

**Clause-6 scope statement — read this before treating any milestone as charter discharge.**
Charter clause 6 names *"project the transition registry over MCP + publish the A2A agent card."*
After the split, THIS DOC DELIVERS THE A2A HALF AND THE THREE ENABLING MILESTONES ONLY. The MCP
half — "project the transition registry over MCP" — is carried by the child doc and is blocked on
the upstream dispatch seam (`ailang#885`; the blocking predicate, as a runnable command with
controls, lives in the child). Clause 6 is NOT satisfied by this parent alone, and this doc does
not quietly narrow the clause: it partitions it, with the blocked half named and tracked.

**Item**: `w-mcp-projection`  
**Clause**: clause-6 (protocol-native boundary; not DESIGN.md M6) — the A2A half + enablers; the
MCP half is the child's  
**Estimate**: ~1.25 World days (P6.T ~0.1 + P6.D ~0.15 + P6.V ~0.3 + P6.B-A2A ~0.7); no upstream
wait remains IN THIS DOC (the child's wait is `#885`, and it is the child's to carry)  
**Author**: DESIGN-DOC-CREATOR (original 2026-07-28, rotation designer codex/gpt-5.6-sol;
revision round 2026-08-25; split applied 2026-08-25, iteration 125)  
**Verified against**: pinned **AILANG v0.30.0**, commit `e37b370`, at `/tmp/ailang-v0300/ailang`
(the `.ail` verifier axis) **and** upstream **`github.com/sunholo-data/ailang` tag `v0.33.2`
(`63e7909f`)** (the Go serving-seam axis — D-WORLD-5's two independent version axes); World `dev`
at `2e44e3e` (split-round anchors re-derived there — premise row N25)  
**Date**: 2026-07-28; **revised in place 2026-08-25**; **split applied 2026-08-25 (iteration 125)**

> **Scheduling truth (2026-08-25, split round).** Every milestone in THIS doc is executable NOW:
> `serveapi/protocol` shipped in upstream tag `v0.33.2` (measured — premise rows N1–N3, latest
> release re-confirmed this round in row N22), this repo has landed the transition registry
> (`host/transitionreg/`, item 11, iter-75) and the broker session API (`host/broker.Session`,
> `NewSession(store, episodeID, grants, registry)` at `broker.go:87`, re-derived this round), and
> the last open fork (session credential carrier) is CLOSED by D-WORLD-26 = A. The single internal
> residual is milestone P6.V (the commit-boundary law's Z3 proof), a milestone of this doc. The
> FIRST milestone is a toolchain move — `v0.33.2`'s `go.mod` declares `go 1.26.6` while this repo
> pins `GOTOOLCHAIN: go1.25.6` at `ci.yml:21,102` and in `go.mod` (re-derived this round) — and it
> is independently mergeable before any dependency lands. Nothing here waits on `#885`.

---

## Motivation

Clause 6 requires the available World transitions to leave the daemon through existing protocols:
MCP for callable tools and A2A for discovery, with the projected surface filtered by the
capabilities of each session. This is an authority boundary, not documentation: an unauthorized
transition must be absent from discovery, absent from the A2A card, and rejected if invoked by
name. **This doc delivers that boundary over A2A** — the agent card at `/.well-known/agent.json`,
the A2A endpoint at `/a2a/`, and session-scoped invocation through propose → verify → commit —
plus the three enabling milestones (toolchain floor, pinned dependency, verified commit-boundary
law) that the MCP half will ALSO consume when it unblocks. The MCP surface itself lives in the
child doc, blocked on the upstream dispatch seam.

AILANG owns the A2A wire implementation. World reuses that machinery — the pinned
`serveapi/protocol` wire types and helpers — and writes no JSON-RPC, SSE, MCP, or A2A codec. Live
v0.30.0 evidence, however, rules out treating the CLI flags as the completed boundary:

- `--caps` is process-wide and does not filter discovery;
- one static API key is not a World session;
- `--routes-only` removes the eight embedded `std/io` exports but still publishes the built-in
  `submit_feedback` tool;
- module loading is cwd-sensitive; and
- static `.ail` exports are not the dynamic, worldd-backed transition registry.

The July draft therefore took charter reuse path **(c)**: request a narrow public serving seam
from the existing `internal/apiserver`, then let worldd supply session-scoped discovery and
invocation callbacks. **That request has been DELIVERED — for the A2A half** (upstream `#498`
lane B in v0.33.0/v0.33.1, then the protocol-only module ask `#764` in **v0.33.2**). Quorum
round 3 then established, confirmed first-party, that the delivery covers the A2A half ONLY:
`protocol` carries the whole A2A wire surface plus MCP *envelope* helpers, but no MCP JSON-RPC
dispatch. The design here therefore IS: **import `serveapi/protocol` pinned at `v0.33.2`; World
implements its own A2A HTTP handlers and callback-bounding over those types** (Decision 2 — the
upstream handlers deliberately live outside `protocol` and cannot be imported without failing the
zero-cloud gate). Upstream remains the wire-contract owner. World remains the state, session,
capability, and transaction owner — and, by measured necessity, the A2A-handler owner.

## Premises (hard constraints)

- **P1 — no new wire protocol.** A2A remains at `/.well-known/agent.json` and `/a2a/`; responses
  keep upstream framing. World does not decode and re-encode protocol traffic in a reverse proxy,
  and this doc mounts NO `/mcp/` endpoint at all — the MCP surface, including its SSE-framed HTTP
  responses, travels with the child doc.
- **P2 — exact per-session authority.** A session resolves to an immutable capability snapshot
  for one request. The same predicate filters A2A discovery (card skills) and invocation. A
  missing, expired, or unknown session fails closed.
- **P3 — transition truth stays in World.** The projected entries come from the landed
  transition registry, not a startup copy of `.ail` exports. v0.30.0 static serving proves only
  protocol projection, not dynamic registry integration.
- **P4 — propose → verify → commit remains mandatory.** Protocol invocation enters the same
  broker/session/transaction path as every other agent action. The projection receives no direct
  `store.Commit`, `SetRegistryHead`, or REST `/v1/commit` authority.
- **P5 — worldd stays the sole writer.** The projection is hosted in the existing worldd process
  and uses its already-open store/broker handles. It never calls `store.Open`, starts a second
  writer, or weakens `WriterAlreadyActive`.
- **P6 — the landed REST v1 contract is frozen.** No route is renamed, removed, or repurposed.
  The A2A endpoints are additive and do not translate failures into the REST `daemon.APIError`
  envelope. Existing REST callers see byte-compatible behavior.
- **P7 — local-first and zero-cloud remain enforced.** The path-(c) import is now measured, not
  hypothetical: `serveapi/protocol` at `v0.33.2` imports **only the standard library** across its
  four source files, and admitting it moves the daemon-core dependency closure by **exactly one
  package** (249 → 250 over both gated patterns). It is admitted by **one allowlist line naming
  the PACKAGE path** `github.com/sunholo-data/ailang/serveapi/protocol` — never the module root,
  which would silently admit `internal/apiserver` and its measured 476 disallowed packages.
- **P8 — the projection adds no new AILANG law; milestone P6.V adds exactly one.** The protocol
  adapter's invariants remain cross-request/session behavior in Go integration/conformance tests
  (unchanged from July). But prerequisite 3's residual is precisely a missing pure-core law: none
  of the **10** identities pinned in `REQUIRED_VERIFIED` (`scripts/verify_ail.sh:274-279`) is a
  commit-boundary law. P6.V therefore adds a commit-boundary identity in `world/*.ail` and pins
  it in the manifest — this item now DOES touch `world/` and the `verify_ail.sh` manifest,
  reversing the July premise. P6.V must design around the recorded v0.30.0 limitation: a contract
  on any `Proposal`-taking predicate Z3-errors `unknown sort 'Proposal'` (ADT-bearing sorts)
  while `ai-check` exits 0 **silently** — the gate's JSON `verify.errors` check is the detector,
  and the law must be stated over encodable (record/scalar) types, with the
  `w-m1-ailang-hardening` test-only fallback if it cannot be.

## High-Impact Decisions

| Decision | Why High Impact | Chosen By | Deadline | Change Cost |
|---|---|---|---|---|
| Import `serveapi/protocol` pinned at **v0.33.2**; World writes the A2A handlers | the requested seam shipped for the A2A half; D-WORLD-5 (Mark, attended, 2026-08-17) rules the pinned import the sanctioned path; the handler-bearing packages would fail the zero-cloud gate | D-WORLD-5 + 2026-08-25 measurement | resolved | high |
| **MCP half SPLIT into `w-mcp-dispatch-projection.md`, blocked on `#885`** | `protocol` has no MCP JSON-RPC dispatch and both local routes are closed by this repo's own guardrails (round-3 finding, confirmed first-party) | quorum round 3 + controller routing call (iteration 125) | resolved (split) | high |
| **Session credential carrier = `Authorization: Bearer` (Arm A)** | session authority surface in a repo where clause 3 is inviolable; never an API key; fail closed | **D-WORLD-26** (Mark, attended, issue `#89`, 2026-08-25T19:06:41Z) | resolved | high |
| Toolchain floor moves FIRST (`go1.25.6 → go1.26.6`) | `v0.33.2` requires `go >= 1.26.6` (two-arm measured); independently mergeable before any dependency | live evidence | P6.T | medium |
| Allowlist admits the PACKAGE path, never the module root | the matcher is prefix-based; the root would admit `internal/apiserver`'s 476 disallowed packages | measured matcher semantics | P6.D | high |
| One authorization predicate drives card/call | prevents discovery and invocation drift | clause 3/6 | compile | high |
| Projection is additive inside worldd | preserves single writer and frozen REST v1 | shipped M2 | route wiring | medium |
| Commit-boundary law is pure-core `world/*.ail` work, scoped HERE as P6.V | "verified" in this repo means Z3-proven and pinned in `REQUIRED_VERIFIED`; the Go surface alone does not discharge it | charter prereq-3 residual | P6.V | medium |

### Design Freeze

- [ ] Do not implement MCP, SSE, JSON-RPC, or A2A codecs in this repository: wire types,
  envelope helpers, and name validation come from the pinned `serveapi/protocol`; World-owned
  handlers do not re-declare parallel wire structs or hand-format envelope bytes.
- [ ] Do not implement ANY MCP surface in this doc's scope — no `/mcp/` route, no JSON-RPC
  dispatch, no `tools/list`/`tools/call`, no MCP SDK import. That entire half is the child doc's,
  and it is blocked on `#885`. (This is the round-3 objection made structural: the two closed
  routes stay closed.)
- [ ] Do not graft the World store or scheduler onto `ailang serve-api`.
- [ ] Do not use `--caps` as evidence of per-session filtering.
- [ ] Do not use a single `--api-key-*` value as a World session model — and, per D-WORLD-26's
  constraint (i), never treat the `Authorization: Bearer` value as an API key: it is a SESSION
  credential, resolved and authorized per request.
- [ ] Do not call the static sketch export list a transition registry.
- [ ] Do not expose `submit_feedback`, `std/io`, introspection helpers, or any tool not present in
  the session's exact authorized transition set.
- [ ] Do not add a direct A2A-to-`store.Commit` or A2A-to-REST-commit path.
- [ ] Do not start P6.B-A2A before P6.T (toolchain), P6.D (pinned dependency + narrow allowlist),
  and P6.V (the **verified** commit-boundary law) have each landed CI-green. The July wording's
  three external prerequisites are discharged: the seam is `serveapi/protocol@v0.33.2`, the
  registry is `host/transitionreg/`, the session API is `host/broker.Session`; the Go half of the
  commit-boundary contract (atomic not-started-versus-committed via `JournalIntent` +
  `bindCommitIntentTx`, stable `InvocationID`, queryable `GetReceipt`/`GetEffectReceipt`) is
  landed — only the Z3 proof (P6.V) remains.
- [ ] Do not import `serveapi` proper, its embedded MCP/A2A handlers, or `CallbackRunner` — their
  closure spans the MCP SDK's 9 module roots and fails `TestDaemonDependencyAllowlist` by
  construction. The allowlist entry is the package path, never the module root.
- [ ] Do not claim that cancellation prevents a durable mutation once the coordinator has ACCEPTED a
  commit, and never return a definitive "not committed" result when the outcome is unknown.
- [ ] Do not treat logical `context` cancellation as proof that a connection closed — the socket
  itself must be observed closed.
- [ ] Do not relax or change the frozen D7 daemon timeout constants; `host/projection` performs NO
  `ResponseController` deadline changes at all (the SSE-relaxation question is MCP-bound and
  travels with the child doc).
- [ ] Do not leave any projection request or dependency wait unbounded.

## Decision 1 — Reject Reuse Paths (a) and (b) on Evidence

**Path (a), modules served directly, is insufficient.** On the pinned binary, `--routes-only`
does remove ordinary exports and the eight embedded `std/io` functions, but it leaves
`submit_feedback`. With no `--routes-only`, the sketch directory publishes 27 tools, including
`exit` and `writeBytes`. `--caps ''` publishes the same 27 names. More fundamentally, the flags
hold process-wide state and the source directory is static; neither accepts a World session or a
dynamic transition-registry callback.

**Path (b), a sidecar, does not repair the semantic mismatch.** A single sidecar still has one
capability set and key. A sidecar per session would turn process lifetime into the session model,
but would still expose `submit_feedback`, would still serve static modules, would need a new
worldd-backed invocation bridge, and would multiply cwd, port, cleanup, and crash ownership.
Reverse-proxy filtering would make World parse and rewrite MCP/A2A, which is the forbidden
reinvention. Therefore path (b) adds failure modes without satisfying P2–P4.

**Chosen fallback: path (c).** File an upstream request to export the already-existing serving
machinery behind a public, callback-driven Go API. This conclusion is limited to v0.30.0. If an
existing public API is found before implementation, the sprint may use it only if the same
conformance suite proves P1–P4; the design direction does not change.

**Post-delivery status (revision round, 2026-08-25).** The falsification is scoped precisely:

- **SURVIVED**: the rejections of paths (a) and (b). Their evidence was measured against the
  pinned v0.30.0 binary, which is unchanged; the CLI flags are still process-wide, the sketch
  directory is still static, and a sidecar still cannot model World sessions. Nothing shipped
  upstream alters that evidence.
- **SURVIVED**: path (c)'s direction — upstream owns the wire contract, World supplies
  session-scoped discovery and invocation.
- **FALSIFIED**: the premise "upstream exposes nothing public," and with it the July reading of
  path (c) as *"request a seam over `internal/apiserver` and wait."* The request was filed
  (`#498`, then `#764` for a protocol-only module) and **delivered**: `serveapi/protocol` at tag
  `v0.33.2`. Path (c) is no longer a request; it is an import (Decision 2).

## Decision 2 — The Upstream Serving Seam: DELIVERED as `serveapi/protocol@v0.33.2` for the A2A half; World owns the A2A handlers

**What shipped (measured 2026-08-25 by `git show 'v0.33.2:serveapi/protocol/<f>'`).** Four source
files (`interfaces.go`, `descriptor.go`, `envelope.go`, `a2a_wire.go`, plus `descriptor_test.go`),
importing **only the standard library** (`context`, `encoding/json`, `errors`, `fmt`, `net/http`,
`regexp`, `sort`, `strings`; zero non-stdlib). Exported surface:

- `interfaces.go` — `Session any`; `SessionResolver.ResolveSession(ctx, *http.Request)`;
  `ToolSource.Tools(ctx, Session) ([]ToolDescriptor, error)`;
  `Invoker.Invoke(ctx, Session, Invocation) (InvocationResult, error)`; `Invocation`,
  `InvocationResult`, `AgentInfo`, `AuthorizationError` (with `HTTPStatus()`).
- `descriptor.go` — `ToolDescriptor`, `AuthorizedSurface`,
  `CallerSurface([]ToolDescriptor)` (validates, dedups, and deterministically orders the
  caller-supplied set), `ValidateMCPName(string) error`.
- `envelope.go` — `RequestID([]byte)`, `CallbackMessage(error)`, `WriteMCPEnvelope(...)`,
  `AuthorizationStatus(error) int`, `ErrCallbackCapacity`.
- `a2a_wire.go` — `A2ARequest`, `A2ATaskSendParams`, `A2AMessage`, `A2AContent`,
  `A2AError(...)`, `A2AResult(...)`.

Against the July draft's seven behavioral requirements: (2) principal-before-discovery, (3)
caller-supplied exact descriptors, (4) invocation routed back with the resolved session, (5) one
caller-supplied set feeding both surfaces, and (6) no built-in tool are satisfied **as types and
helpers** — `Session`/`SessionResolver`/`ToolSource`/`Invoker`/`CallerSurface` are the seam, and
there is no built-in tool because there is no handler to embed one. Requirements (1) mount-on-mux
and (7) served conformance/SSE framing are **deliberately NOT in `protocol`**: `CallbackRunner`,
the embedded A2A `http.Handler`, and the MCP handler live in `serveapi/` proper, whose dependency
closure spans the MCP SDK's 9 module roots (charter iter-90 measured the facade seam at **476
disallowed packages across 86 module roots**) — importing it would fail the exact zero-cloud gate
this design must pass.

**The round-3 sharpening (confirmed first-party; re-derived again this round, premise row N23):
the missing piece is asymmetric.** For A2A, everything needed is importable and the handler shape
is proven: upstream's own `serveapi/a2a_handler.go` is **180 lines over stdlib + `protocol` with
0 MCP-SDK imports** (re-derived this round via
`gh api 'repos/sunholo-data/ailang/contents/serveapi/a2a_handler.go?ref=v0.33.2'`; control in the
same call: `mcp_handler.go` carries **1** SDK import in 187 lines). For MCP, `protocol` carries
envelope HELPERS (`WriteMCPEnvelope`, `RequestID`, `ValidateMCPName`, `AuthorizationStatus`) but
**no JSON-RPC method dispatch** — `mcp_handler.go` does not even string-match method names
(`grep -c 'tools/list\|jsonrpc'` → **0**, while the same grep finds A2A method strings in
`a2a_handler.go` → **2**), because it hands the entire dispatch to
`github.com/modelcontextprotocol/go-sdk`. Therefore **World writes its own A2A HTTP handler and
callback-bounding here**, and the MCP handler question moves to the child doc until `#885`
delivers a dispatch seam. That is D-WORLD-5 (ARM A: import upstream, pinned) executing as written
on the half it can reach, not a new human ask.

**Version discipline — resolved by measurement.** D-WORLD-5 named `v0.33.1`; its own condition is
"pin the first released/tagged upstream revision that contains this seam." Measured 2026-08-25:
`serveapi/protocol` is **absent at `v0.33.1`** (`git ls-tree v0.33.1 serveapi/protocol/` → 0
entries; same-call control: `serveapi/` itself lists 2 entries) and **present at `v0.33.2`**
(`63e7909f`), which is the **latest tag**. Re-confirmed in the split round without a local
checkout (premise row N22): the latest upstream RELEASE is still `v0.33.2` and no `v0.34.*` tag
exists (`matching-refs/tags/v0.34` → empty; control `tags/v0.33` → 3 refs). The pin is therefore
**`v0.33.2`**, recorded in `go.mod`/`go.sum` at P6.D. The upstream stdlib-only guarantee is
CI-enforced upstream (`scripts/check_protocol_closure.sh` ships in the tag — charter iter-120,
inherited), and re-enforced here first-party by the allowlist gate.

## Decision 3 — Session-Scoped Projection in worldd (A2A)

The additive `host/projection` adapter has four responsibilities, each naming a **landed**
interface (anchors re-derived at `2e44e3e` this round, premise row N25):

1. resolve the session credential through the landed broker session API
   (`host/broker.Session`, `NewSession(store, episodeID, grants, registry)` at `broker.go:87`) —
   implementing `protocol.SessionResolver`. **Carrier: CLOSED by D-WORLD-26 = ARM A.**
   `ResolveSession(ctx, r)` reads the standard `Authorization: Bearer <session-credential>`
   header — and nothing else. Two constraints travel with the ruling and bind the resolver's
   contract: **(i)** the static `serve-api` API key remains forbidden as a session model
   (iteration 24 measured it process-wide, so it cannot represent a session) — a `Bearer` value
   here is a SESSION credential and never an API key; **(ii)** clause 3 still binds — the
   resolver fails CLOSED on an absent, malformed, or unknown credential, never degrading to an
   unauthenticated surface. Arm B (`X-World-Session`) is REJECTED and must not be read, even as
   a fallback.
2. read one transition-registry snapshot (`host/transitionreg.StoreReader.ReadSnapshot(ctx)` →
   `Snapshot`, with `Snapshot.Lookup(id)`/`Snapshot.List()`) and one capability snapshot —
   implementing `protocol.ToolSource` via `protocol.CallerSurface`;
3. derive an exact ordered tool set using the registry's stable transition identity
   (`transitionreg.Descriptor`, `SortedDescriptors`); and
4. dispatch an authorized invocation into the normal propose → verify → commit coordinator —
   implementing `protocol.Invoker`, with `transitionreg/bind.go`'s typed refusals
   (`TransitionAbsentError`, `AccessDeniedError`, `ProposalMismatchError`) mapped through
   `protocol.AuthorizationError`/`AuthorizationStatus`.

The authorization function is conceptually:

`visible(session, transition) = session is live AND transition capability ∈ session capabilities`.

That expression is explanatory prose, not AILANG source. The implementation reuses the landed
capability predicate; it does not create a second policy engine in `host/projection`.

**One snapshot per request.** Card generation may observe a newer registry on the next request,
but a single request never mixes registry or capability epochs. Invocation re-resolves the
session and transition at dispatch time; possession of a previously listed skill name conveys no
authority.

**Fail-closed cases.** Unknown/expired session, absent transition, capability denial, registry
read error, or broker unavailability returns the upstream protocol's structured error form. No
case falls back to REST or direct store access.

## Decision 4 — Surface Identity, A2A Card, and Compatibility

Tool identity is the transition registry's stable transition ID. Display text and schemas are
projection metadata; they do not become authority. Ordering is deterministic by stable ID so
cards are diffable.

For a given session snapshot:

- A2A `skills[].id` equals the authorized registry ID set exactly;
- every listed skill is invocable through the same session;
- every unlisted transition is rejected even when invoked by a guessed/stale name; and
- `std/io.*`, bare IO names (`exit`, `writeBytes`, etc.), and `submit_feedback` are absent unless
  World someday registers and authorizes transitions with those exact IDs. This item registers
  none.

The cross-surface obligation the old AC3 stated — that the eventual MCP `tools/list` name set
must exactly equal the A2A `skills[].id` set per session — is NOT dropped by the split: it
travels with the child doc and binds the MCP surface when it becomes buildable.

The agent card is emitted by World's A2A handler using `protocol`'s wire representation
(`AgentInfo`, `A2ARequest`/`A2AResult`/`A2AError` and the descriptor set from `CallerSurface`).
The test must not assume a hand-authored JSON shape beyond protocol-required fields and the exact
skill-ID set.

## Decision 5 — Process, Endpoint, and Package Layout (S2/S3)

| Path | Change |
|---|---|
| `.github/workflows/ci.yml` | P6.T: `GOTOOLCHAIN: go1.25.6 → go1.26.6` at BOTH sites (`:21`, `:102`) — the only workflow change |
| `go.mod` | P6.T: `go 1.25.6 → 1.26.6`; P6.D: `require github.com/sunholo-data/ailang v0.33.2` |
| `go.sum` | P6.D: pin movement, including the measured indirect bumps (`go-isatty v0.0.22`, `x/sys v0.47.0`) |
| `host/daemon/daemon_test.go` | P6.D: ONE `allowedDepModules` line (the package path) + narrowness test; comment updated to "module roots or package paths" |
| `world/*.ail` + `scripts/verify_ail.sh` | P6.V: commit-boundary law + its `REQUIRED_VERIFIED` pin |
| `host/projection/` | P6.B-A2A: World-owned A2A handlers + agent card over `protocol` types, session resolver (`Authorization: Bearer`), adapter, and tests — no `/mcp/` route |
| `host/daemon/daemon.go` | P6.B-A2A: additive mounting of the World-owned A2A handlers |
| `cmd/ailang-worldd/` | P6.B-A2A: minimal flags/config for enabling projection and session credential policy |

**Why is this not an AILANG package?** Protocol serialization, HTTP request authentication, and
connection lifecycle are host effects (S2). The transition definitions and policy remain
package/registry data; only their OS/network projection is Go.

**Why is this not kernel growth?** `host/projection` is an adapter at DESIGN §3.7's system
boundary. It owns no transition semantics, store schema, registry update, or capability law.
The projection milestones (P6.D/P6.B-A2A) add nothing to `world/`, `host/store`, or
`host/registry`. The one `world/` addition in this item is P6.V's commit-boundary LAW — a proof
over the already-landed store/broker semantics, which is the opposite of kernel growth: it pins
existing behavior, adds none.

**Dependencies — measured, no longer hypothetical.** One new direct dependency:
`github.com/sunholo-data/ailang v0.33.2`, whose admitted surface is the single stdlib-only
package `serveapi/protocol`. Closure delta over both gated patterns
(`./host/daemon/... ./cmd/ailang-worldd/...`): **249 → 250 packages, the one addition being the
package itself; removed set empty** (controller-measured with a sentinel control at iter-124;
the pristine 249 re-derived first-party in the revision round). The upgrade is **not
transitively free**: `go get ...@v0.33.2` also bumps two already-allowlisted indirect roots —
`github.com/mattn/go-isatty v0.0.20 → v0.0.22` and `golang.org/x/sys v0.46.0 → v0.47.0` — and
rewrites `go.mod`'s `go` directive to `1.26.6` (re-derived first-party on a copy of this repo's
`go.mod`). Both bumped roots are already in `allowedDepModules`, so the allowlist is unaffected,
but `go.sum` moves; AC16 names this as its own acceptance row. No lockfile/module-pin gate exists
in `verify_go.sh` to collide with (`grep 'go\.sum|go mod verify|go mod tidy'` → 0 hits; same-call
control on the same file fired). The executor must not weaken the allowlist beyond the ONE
package-path line.

## Decision 6 — Bounded waits across the projection path (Standing rule 6)

Every projection request derives one bounded context whose deadline is the earlier of the client
cancellation/deadline and a configured finite server maximum. The server maximum is mandatory,
positive, finite, and validated at startup; zero, negative, or an omitted value does not mean
"unlimited." The same context is propagated without replacement through session resolution,
registry snapshot acquisition, authorization, broker dispatch, and the complete propose → verify
→ commit path.

Client disconnect cancels that context and all work started for the request. Each dependency must
observe cancellation, release request-owned resources, and return promptly. Deadline expiry or
client cancellation is returned in the upstream protocol's structured timeout/cancellation error
form when the transport remains writable; a disconnected transport is closed after cancellation.
It is never mapped to `daemon.APIError`, retried, or used as a reason to fall back to REST, direct
store access, or another capability/session source.

**The commit boundary (quorum round 2, `gpt5-6-sol` — the reviewer's own wording, adopted verbatim
under the narrow-refinement carve-out).** An earlier draft of this paragraph claimed that
cancellation after commit *begins* still yields no store/log mutation. That is not achievable: an
HTTP context can expire while an atomic store commit is already in flight, so unless the
prerequisite coordinator/store exposes a verified atomic "not-started versus committed" contract,
the commit may succeed even though the caller observed cancellation. The truthful contract is the
reviewer's:

> Before the coordinator accepts a commit, cancellation guarantees no durable mutation. Commit uses
> a stable invocation/idempotency ID and has a defined point of no return. Once accepted, the
> bounded commit critical section may complete despite client disconnect; its durable receipt is
> recorded and queryable/replayable. The caller must never receive a definitive "not committed"
> result when the outcome is unknown.

Proposals remain non-durable until verification succeeds and the coordinator ACCEPTS the commit. A
dependency that returns success after the deadline cannot revive the request. Request cleanup
cancels child work and discards staged state, so a stalled resolver, registry, authorization
provider, broker, or verifier terminates within the configured bound with no dispatch/store/log
mutation — the accepted commit critical section being the single, explicitly bounded exception
above.

**The reviewer's contract is now GROUNDED in landed public APIs (measured 2026-08-25, anchors
re-derived in the split round — premise row N25)**: the atomic not-started-versus-committed
statement is `JournalIntent` (`host/store/journal.go:28`, "the canonical statement of a planned
commit"), bound inside the commit transaction by `bindCommitIntentTx` (`host/store/store.go:1025`);
the stable invocation/idempotency ID is `InvocationID`, threaded through journal, receipts, and
recovery; the queryable durable receipt is `Store.GetReceipt` (`journal.go:813`) /
`GetEffectReceipt` (`journal.go:852`), consumed by `recoverCommitPending`
(`host/broker/recover.go:126`). What P6.B-A2A still needs from P6.V is the word the prerequisite
hinged on: a Z3-**verified** commit-boundary law pinned in `REQUIRED_VERIFIED` — the Go surface
exists; the proof does not.

**SSE lifecycle — MOVED WITH THE MCP SURFACE (split round).** Quorum r2's `gemini-3-1-pro`
objection concerned zombie TCP connections on the long-lived `/mcp/` SSE stream, and the r3
revision made that stream World-authored. Under the split there is NO SSE stream in this doc:
A2A is single bounded request/response, served entirely within the frozen D7 deadlines, and
`host/projection` performs no `ResponseController` deadline relaxation (AC14 asserts this at
source level). What SURVIVES here from the objection is its OS-level kernel: a disconnected or
deadline-expired transport must be OBSERVED closed (`http.Server.ConnState` tracking or a client
read error), never merely context-cancelled (AC13). The stream-lifetime maximum, the
`/mcp/`-route-local deadline relaxation, and mutations
`MUT-LEAK-SSE-CONN`/`MUT-SSE-REST-DEADLINE` travel with the child doc, where the streaming
handler will exist.

## Milestones (each independently CI-green and mergeable — hard repo convention; ~1.25 World days)

### P6.A — Upstream repro + conformance fixture — **DONE (iter-24, 2026-07-28)**

Retained as record: the upstream finding was filed (`sunholo-data/ailang#764`, following `#498`)
and this design landed. That finding is now **RESOLVED BY DELIVERY** — see Upstream Findings.
The frozen two-session conformance fixture carries forward into P6.B-A2A unchanged (its A2A leg;
the MCP leg of the fixture travels with the child).

### P6.T — Toolchain floor `go1.25.6 → go1.26.6` (~0.1d) — **FIRST, and independently mergeable**

No dependency, no new code. Move `GOTOOLCHAIN: go1.25.6 → go1.26.6` at **both** `ci.yml` sites
(`:21`, `:102`) and `go.mod`'s directive `go 1.25.6 → 1.26.6`. The `verify_go.sh` deny-list
(`:214-224`, enumerating exactly `go1.26.0`–`go1.26.5`) is **untouched** — `go1.26.6` is not in
it, and the committed canary is the version-agnostic detector.

- Why first: `v0.33.2` requires `go >= 1.26.6` (two-arm measured: `go get ...@v0.33.2` rc=1
  under `GOTOOLCHAIN=go1.25.6` with the exact floor message, rc=0 under `go1.26.6`;
  known-positive control `go get github.com/google/uuid@v1.6.0` rc=0 under `go1.25.6`).
- Rig nuance (measured 2026-08-25): the rig's `go` **binary** is `go1.26.4` (deny-listed);
  under `GOTOOLCHAIN=auto`, toolchain selection follows `go.mod` — today (`go 1.25.6`) it
  resolves to `go1.26.4` and the deny-list correctly refuses it, and in a `go 1.26.6` module it
  resolves to `go1.26.6` (measured). So until P6.T lands, local gate runs must pin
  `GOTOOLCHAIN=go1.25.6` explicitly; after it lands, `auto` is safe.
- Acceptance: canary re-run under the new toolchain — repro prints `OK` under `go1.26.6` AND the
  known-bad control prints `BUG` under `go1.26.5` (both re-run first-party 2026-08-25); full
  `verify_go.sh` rc=0 (deny-list, driver drift gate, armed race control, build+test); both CI
  jobs green with zero other diffs.

### P6.D — Pinned dependency + ONE narrow allowlist line (~0.15d)

`go get github.com/sunholo-data/ailang@v0.33.2`, then add **exactly one** `allowedDepModules`
entry: the **PACKAGE path** `github.com/sunholo-data/ailang/serveapi/protocol` — never the
module root. The matcher (`disallowedDeps`, `daemon_test.go:801`:
`d == m || strings.HasPrefix(d, m+"/")`) treats entries as path prefixes, so the package-path
entry admits exactly that package (and any future subpackage beneath it), while the module root
would admit `internal/apiserver`'s measured 476 disallowed packages. Update the
`allowedDepModules` doc comment ("module roots") to "module roots or package paths".

- Include a **narrowness test**: with the new entry in place,
  `disallowedDeps(["github.com/sunholo-data/ailang/internal/apiserver"])` (plus a
  representative cloud path) is non-empty — proving the entry did not widen the gate.
- Include a compile-visible use (or the P6.B-A2A skeleton's first import) so the closure actually
  contains the package: closure over both patterns moves **249 → 250**, addition =
  the package itself, removed set empty.
- Named risk with its own acceptance row (AC16): the pin also bumps `go-isatty → v0.0.22` and
  `x/sys → v0.47.0` (already-allowlisted roots) and `go.sum` moves; no other graph movement is
  permitted.
- Merge criterion: both CI jobs green; `TestDaemonDependencyAllowlist` green with the probe
  demonstration recorded (allowlist minus the new line → REDs naming exactly one intruder).

### P6.V — VERIFIED commit-boundary law (pure-core `world/*.ail`) (~0.3d)

Discharges the one residual of prerequisite 3. Add a commit-boundary identity to `world/*.ail`
stating the reviewer's contract in pure form — commit acceptance has a defined point of no
return; before acceptance nothing durable exists; after acceptance exactly one receipt exists
under the stable invocation ID, and outcome-unknown is never reported as "not committed" — and
pin it in `REQUIRED_VERIFIED` (`scripts/verify_ail.sh:274-279`), raising the floor from 10
identities **on the Z3 path**. **Conditional, not a mandate** (applying `gemini-3-1-pro`'s
round-4 non-blocking observation verbatim): if the encodable-shape fallback below fires, the
identity floor correctly stays at 10 and the milestone closes on the named test-only law plus
its limitation row instead. AC10 and AC17 already state the conditional correctly; this
sentence is what an executor could otherwise misread as a strict obligation to raise the floor
regardless of which path was taken.

- **Load the AILANG language reference via the `ailang-docs` MCP BEFORE writing any `.ail`**
  (charter fluency protocol; binding).
- Design constraint (recorded v0.30.0 limitation, charter): a contract on a `Proposal`-taking
  predicate Z3-errors `unknown sort 'Proposal'` while `ai-check` exits 0 silently; the gate's
  JSON `verify.errors` check catches it. State the law over encodable record/scalar shapes
  (e.g., intent/receipt states keyed by invocation ID), not ADT-bearing ones; if no encodable
  statement exists, fall back to the `w-m1-ailang-hardening` pattern — named test-only law plus
  an explicit limitation row — and say so in the implementation report.
- Merge criterion: `verify_ail.sh` green — with the raised identity floor **if the Z3 path was
  taken**, or at the unchanged floor of 10 with the test-only law and its limitation row **if the
  fallback fired**; no Go changes required; both CI jobs green.

### P6.B-A2A — Session projection + A2A conformance, World-owned handlers (~0.7d)

Re-cut of the old P6.B under the split: the A2A half only. Starts only after P6.T, P6.D, and
P6.V are green. Add `host/projection` with a **World-authored** A2A handler and agent-card
endpoint built exclusively on `protocol` types/helpers, mounted additively in the daemon;
connect only to the landed `transitionreg`/`broker.Session`/coordinator interfaces. The session
resolver reads `Authorization: Bearer` per D-WORLD-26 = A and fails closed. Upstream's own
`a2a_handler.go` (180 lines, stdlib + `protocol`, 0 SDK imports — re-derived this round, premise
row N23) is the existence proof of exactly this shape.

- Prove exact two-session discovery/card/invocation behavior, stale-name denial, A2A wire-form
  conformance against the frozen P6.A fixture's A2A leg, bounded cancellation through every
  dependency, both sides of the commit boundary (AC13), OS-level socket closure on
  disconnect/expiry, REST deadline regression, single-writer preservation, and the dependency
  allowlist.
- The MCP half of the old P6.B — `/mcp/` handler, JSON-RPC dispatch, SSE stream lifetime — is
  NOT here and does not gate this milestone; it is the child doc's scope, blocked on `#885`.
- Merge criterion: both CI jobs green, no skips, every acceptance mutation demonstrated RED then
  reverted GREEN.

**Estimate honesty.** July's "~1.2d after prerequisites" became ~1.55d all-in at the revision
round; the split re-cuts it to **~1.25d in this doc** (P6.T ~0.1 + P6.D ~0.15 + P6.V ~0.3 +
P6.B-A2A ~0.7). The difference is not a discount: the MCP share of the old P6.B (~0.3d, plus
whatever the dispatch seam actually demands) moves to the child and is deliberately UNESTIMATED
there — its shape depends on what `#885` delivers, and estimating against an unknown seam would
be fiction. Nothing external remains on THIS doc's critical path. If P6.V's encodable statement
proves harder than one law, split P6.V out rather than absorbing the overrun silently — but do
not start P6.B-A2A before it lands.

## Files to Create/Modify

| File | Milestone | Purpose |
|---|---|---|
| `design_docs/planned/w-mcp-projection.md` | P6.A (done) | This design; revised in place 2026-08-25; SPLIT applied 2026-08-25 (iteration 125) |
| `design_docs/planned/w-mcp-dispatch-projection.md` | split child | Carries the MCP half + the round-3 objection verbatim; blocked on `#885`; NOT quorum-cleared |
| upstream issue in `sunholo-data/ailang` | P6.A (done) | `#764` — protocol-only module; RESOLVED by `v0.33.2` (`#885` — dispatch seam — is the CHILD's blocker, tracked there) |
| `.github/workflows/ci.yml` | P6.T | `GOTOOLCHAIN` at `:21` and `:102` → `go1.26.6` (only change) |
| `go.mod` | P6.T / P6.D | `go 1.26.6` directive; `require github.com/sunholo-data/ailang v0.33.2` |
| `go.sum` | P6.D | Pin + measured indirect bumps (`go-isatty v0.0.22`, `x/sys v0.47.0`) |
| `host/daemon/daemon_test.go` | P6.D / P6.B-A2A | ONE package-path allowlist line + narrowness test; REST/single-writer/dependency regression |
| `world/*.ail` | P6.V | Commit-boundary law (encodable shape; fluency protocol first) |
| `scripts/verify_ail.sh` | P6.V | `REQUIRED_VERIFIED` gains the commit-boundary identity |
| `host/projection/projection.go` | P6.B-A2A | Session-scoped adapter; World-owned A2A handlers + agent card over `protocol` |
| `host/projection/projection_test.go` | P6.B-A2A | Exact-set, denial, wire-form, bounded-wait, socket-closure tests |
| `host/daemon/daemon.go` | P6.B-A2A | Additive mount of the World-owned A2A handlers |
| `cmd/ailang-worldd/main.go` and tests | P6.B-A2A | Minimal enablement/config |

No schema, benchmark baseline, or REST route changes. **No `/mcp/` route and no MCP handler file
exists anywhere in this doc's scope** — that is the split, not an omission. The July claim "no
`.ail` or CI workflow file changes" remains **withdrawn**: P6.T touches `ci.yml` (toolchain pins
only) and P6.V touches `world/` plus the `verify_ail.sh` manifest — both deliberate, both named
above.

## Conflict Surface

- **HTTP Server Timeouts vs the A2A route.** The existing daemon has frozen D7
  `ReadHeaderTimeout` 5s, `ReadTimeout` 30s, `WriteTimeout` 30s, and `IdleTimeout` 120s
  (constants at `host/daemon/daemon.go:78-97`, wired in `newServer` at `daemon.go:619-628` —
  re-derived at `2e44e3e` this round: `const (` at `:78`, `func newServer` at `:619`). Under the
  split this doc adds ONLY plain bounded request/response endpoints (`/.well-known/agent.json`,
  `/a2a/`), which run entirely WITHIN those deadlines: no route in this doc's scope relaxes a
  deadline, and `host/projection` contains no
  `ResponseController.SetWriteDeadline`/`SetReadDeadline` call (AC14's source assertion). The
  revision-round conflict — long-lived `/mcp/` SSE streams versus the 30s write deadline, and the
  finite stream-lifetime bound that must close them — is REAL but is now the child doc's conflict
  surface, recorded there. The D7 constants and every REST `/v1/*` path remain byte-for-byte
  unchanged.
- **worldd single-writer flock.** Projection uses the daemon's handle; no `store.Open`, sidecar
  writer, or lock change. Existing cross-process `WriterAlreadyActive` tests remain green.
- **Frozen REST v1 route table.** The shipped mux patterns remain unchanged:
  health, head, world, object, log-entry, log-range, registry wildcard, and commit. The A2A
  endpoints are additive at upstream-standard paths; REST response bytes remain
  regression-tested.
- **Shared JSON error envelope.** `/v1/*` continues using `daemon.APIError`. A2A failures use
  upstream protocol errors; mapping them into `APIError` would conflate protocols.
- **`host/registry` vs `host/transitionreg`.** The interpreter epoch registry (`host/registry`,
  `world/epoch-registry/v1`) and the transition registry (`host/transitionreg`, item 11, landed
  iter-75) are DIFFERENT landed subsystems. This item reads only `transitionreg` snapshots and
  neither renames nor overloads either. The July text "P6.B waits for a separately landed
  transition-registry interface" is discharged.
- **`scripts/verify_ail.sh`.** Floor re-measured in the split round (2026-08-25, full gate run,
  PASSED — premise row N24): **10 required identities verified, 40 named tests pass, 9/9
  world-package steps**. P6.V deliberately RAISES the identity floor by pinning the
  commit-boundary law in `REQUIRED_VERIFIED` (`verify_ail.sh:274-279`) — a manifest change this
  doc names, not an accident.
- **`scripts/verify_go.sh`.** Runs the driver-drift gate, the go1.26.x deny-list (`:214-224`,
  exactly `go1.26.0`–`go1.26.5`), the armed race-detector control, and `go build`/unskipped
  `go test`. P6.T interplay measured 2026-08-25: the rig `go` binary is `go1.26.4`
  (deny-listed), so under `GOTOOLCHAIN=auto` the gate correctly FATALs today and passes once
  `go.mod` says `1.26.6` (selection then resolves to `go1.26.6`, measured). Protocol socket
  tests must run on CI; local sandbox denial is not a skip condition.
- **CI workflow.** The single `CI` workflow has exactly two jobs:
  `ailang-code verify gate` and `go host build + test gate`; the latter also runs benchmark
  smoke. No new workflow/job. P6.T's `GOTOOLCHAIN` edit at `:21`/`:102` is the ONLY workflow
  change in this item.
- **Dependency allowlist.** `TestDaemonDependencyAllowlist` enforces zero-cloud over the real
  transitive daemon/cmd build graph. Matcher semantics verified at HEAD (`disallowedDeps`,
  `daemon_test.go:801`: `d == m || strings.HasPrefix(d, m+"/")`): entries are path prefixes, so
  the PACKAGE-path entry is strictly narrower than the module root — it admits
  `serveapi/protocol` and would admit a future subpackage beneath it, but refuses every sibling
  including `internal/apiserver`. A blanket `github.com/sunholo-data/ailang` (module-root) entry
  is forbidden: it would pass the gate while admitting the measured 476-package cloud subtree.
- **serve-api cwd sensitivity.** Historical (P6.A record) only: P6.B-A2A runs no upstream server
  at all — World-owned handlers are embedded in worldd — so module-loading cwd is moot.
- **Broker/propose-verify-commit.** Direct REST commit and store commit are explicitly outside the
  projection. If the landed broker lacks a single session-aware dispatch entrypoint, P6.B-A2A
  stops; it does not synthesize one in the adapter.
- **Performance baseline.** No existing benchmark is removed. P6.B-A2A adds bounded protocol
  round-trip benchmarks only if measurement shows a material new hot path; this item does not
  rewrite `bench/BASELINE.md`.

## Systemic-Issue Audit

| Question | Finding | Response |
|---|---|---|
| Is the mismatch local to one endpoint? | No; discovery, A2A card, authentication, and invocation share it | Require one upstream callback seam and one World predicate |
| Can configuration solve it? | No; `--caps` and API key are process-wide | Do not encode sessions as flags |
| Can a sidecar solve it safely? | No; static exports, built-in feedback, and lifecycle remain | Reject path (b) |
| Is a new protocol required? | No | Reuse upstream codecs/wire types |
| Does this reveal a missing landed subsystem? | It did in July (only the epoch registry existed); as of 2026-08-25 the transition registry (`host/transitionreg`) and broker session API are landed, and the remaining gap is one missing Z3 law | Consume the landed subsystems; scope the law as P6.V |
| Is the toolchain move a one-off? | No — go1.26.6 also unblocks item 4e's parked remediation | Land P6.T narrowly here; leave 4e to the queue (Deferred Scope) |
| Could a gate pass vacuously? | Yes; empty registry, one identical session, list-only tests, or skipped sockets | Require non-empty unequal session sets, calls, denials, and zero skips |
| Did one doc bundle surfaces with different readiness? | Yes — round 3 surfaced it: MCP (blocked on a missing dispatch seam) and A2A (buildable now) shared one doc, the harder half holding the easier hostage | SPLIT (iteration 125): A2A + enablers stay here; the MCP half is carried in `w-mcp-dispatch-projection.md`, blocked on `#885` |

## Deferred Scope

- **The MCP half of clause 6 — SPLIT OUT, not silently deferred.** The `/mcp/` endpoint, MCP
  JSON-RPC dispatch (`initialize`, `tools/list`, `tools/call`), dispatch-bound envelope framing,
  SSE stream lifetime and the route-local `ResponseController` relaxation, the cross-surface
  MCP≡A2A equality criterion (old AC3), old AC8's SSE-framing assertion, old AC14, and mutations
  `MUT-PLAIN-JSON`/`MUT-LEAK-SSE-CONN`/`MUT-SSE-REST-DEADLINE` are carried by
  [`w-mcp-dispatch-projection.md`](w-mcp-dispatch-projection.md), blocked on `ailang#885`.
- Designing or landing the clause-3 capability law, broker, session store, or transition registry
  (all landed; this item only consumes them).
- Adopting the `serveapi` facade, its embedded MCP/A2A handlers, or `CallbackRunner` — measured
  at 476 disallowed packages across 86 module roots; permanently out, not merely deferred, while
  the zero-cloud gate stands.
- Any A2UI/AG-UI work (DESIGN.md M6 Generated UI).
- Remote/public bind, TLS termination, OAuth, multi-tenant federation, or cloud deployment.
- Changes to the REST v1 contract or shared error envelope.
- Item 4e's parked toolchain remediation: `go1.26.6` also unblocks it (charter iter-120), but
  picking it up is queue business, not this item's scope — P6.T changes only this item's named
  files.
- Hot reload of transition definitions; registry snapshots already provide request-to-request
  dynamism.
- Benchmark thresholds for the clause-4 non-inferiority experiment.

## Acceptance Criteria

The split moved three criteria out of this doc — old AC3's cross-surface MCP≡A2A equality, old
AC8's MCP SSE-framing assertion, and old AC14's SSE/REST deadline separation now live with the
child. Their numbers are re-used below for A2A-scoped criteria so the numbering stays dense; the
old-to-new mapping is recorded in the split-round quorum entry.

- [ ] **AC1 — wire-contract reuse:** World's A2A handler and agent card construct and parse wire
  forms exclusively through the pinned `serveapi/protocol` types and helpers (`ToolDescriptor`,
  `CallerSurface`/`AuthorizedSurface`, `ValidateMCPName`, `AgentInfo`,
  `A2ARequest`/`A2AError`/`A2AResult`, `AuthorizationError`/`AuthorizationStatus`); production
  World code declares no parallel wire struct and hand-formats no envelope bytes. (World OWNS the
  `http.Handler`s — that is Decision 2's scope, not a violation of this row.)
- [ ] **AC2 — exact surface:** for two non-empty sessions with unequal capability sets, the A2A
  `skills[].id` set equals each session's authorized transition-ID set exactly; no extras.
- [ ] **AC3 — card tracks the session:** a session's A2A skill-ID set changes when its capability
  snapshot changes; two concurrent sessions with different snapshots observe different cards from
  the same daemon. (The cross-surface MCP-equality half of the old AC3 travels with the child.)
- [ ] **AC4 — ambient exports absent:** `exit`, `writeBytes`, all eight observed `std/io`
  exports, and `submit_feedback` are absent from the card and skill set.
- [ ] **AC5 — invocation enforcement:** every listed skill dispatches through propose → verify →
  commit; an unauthorized, stale, or guessed name is rejected before broker dispatch and produces
  no store/log change.
- [ ] **AC6 — session failure, carrier fixed (D-WORLD-26 = A):** the resolver reads the session
  credential from the standard `Authorization: Bearer` header ONLY; absent, malformed, unknown,
  or expired credentials fail closed for card and call; no default/global capability set is used;
  the static `serve-api` API-key value is never accepted as a session credential (constraint i);
  no alternative header (including the rejected Arm-B `X-World-Session`) is read, even as a
  fallback.
- [ ] **AC7 — one-snapshot consistency:** a request never mixes registry/capability epochs;
  concurrent registry/session change is observed only on a subsequent request.
- [ ] **AC8 — A2A protocol conformance:** World's A2A handler speaks exclusively `protocol`'s
  wire forms (`A2ARequest`/`A2AResult`/`A2AError`, `AgentInfo`), proven against the frozen P6.A
  fixture's A2A leg (agent card 200 + protocol-required fields + the exact skill-ID set), not
  against a re-derived local shape. (The MCP SSE-framing half of the old AC8 travels with the
  child.)
- [ ] **AC9 — landed behavior preserved:** REST v1 route/body regressions and cross-process
  writer-lock tests remain green; projection never opens a store.
- [ ] **AC10 — honest gates:** `verify_ail.sh` holds the floor re-measured 2026-08-25 in the
  split round (10 required identities, 40 named tests, 9/9 world-package steps — premise row
  N24) plus EXACTLY the P6.V additions — the commit-boundary identity(ies) named in the P6.V
  implementation report, nothing else moves; `verify_go.sh`, benchmark smoke, and both CI jobs
  pass with zero skipped protocol tests.
- [ ] **AC11 — dependency floor:** the printed transitive graph contains no disallowed package
  and `TestDaemonDependencyAllowlist` covers the new package/cmd paths.
- [ ] **AC12 — dynamic source:** changing the transition-registry head changes the next
  authorized card without restart or `.ail` file edits.
- [ ] **AC13 — bounded projection waits:** with the configured finite server bound, a
  barrier-blocked session resolver, registry snapshot, authorization provider, broker,
  proposer, verifier, or response write, and a disconnected client each terminate within that
  bound. The same request context reaches every dependency; a writable transport receives the
  upstream structured timeout/cancellation error and a disconnected transport closes; no retry
  and no REST/store fallback occurs.
  **Cancellation is tested on BOTH SIDES of the commit boundary** (round-2 fix, `gpt5-6-sol`):
  cancelling immediately BEFORE the coordinator accepts the commit asserts **no durable
  mutation**; cancelling immediately AFTER acceptance asserts **exactly one recoverable,
  queryable/replayable receipt** under the stable invocation/idempotency ID — and in neither case
  does the caller receive a definitive "not committed" answer while the outcome is unknown.
  **The fault-injection test must also assert OS-LEVEL socket closure** (round-2 fix,
  `gemini-3-1-pro`, retained under the split for the A2A route): on client disconnect or deadline
  expiry the underlying TCP connection is actively closed — observed via `http.Server.ConnState`
  tracking or a client read error — proving no silent connection leak. A logical Go `context`
  cancellation alone does not satisfy this criterion.
- [ ] **AC14 — no deadline tampering:** no code in `host/projection` (or anywhere in this doc's
  diff) calls `ResponseController.SetWriteDeadline`/`SetReadDeadline` or otherwise relaxes the
  frozen D7 deadlines — asserted at source level — and REST `/v1/*` still uses the unchanged D7
  values (constants `daemon.go:78-97`, wired in `newServer` at `:619-628`), including its 30s
  read/write deadlines and 120s idle timeout. (The route-local SSE relaxation this row FORBIDS
  here is exactly what the child doc must design for `/mcp/` when it unblocks.)
- [ ] **AC15 — toolchain floor (P6.T):** `GOTOOLCHAIN: go1.26.6` at BOTH `ci.yml` sites and
  `go 1.26.6` in `go.mod`; the `verify_go.sh` deny-list still enumerates exactly
  `go1.26.0`–`go1.26.5`; the committed canary prints `OK` under `go1.26.6` while its known-bad
  control prints `BUG` under a deny-listed toolchain in the same run; full `verify_go.sh` rc=0.
- [ ] **AC16 — pinned dependency, narrow admission (P6.D):** `go.mod` requires
  `github.com/sunholo-data/ailang v0.33.2`; `allowedDepModules` gains EXACTLY ONE entry, the
  package path `github.com/sunholo-data/ailang/serveapi/protocol`; the closure over both gated
  patterns moves by exactly +1 package (the package itself, removed set empty); the narrowness
  test proves `internal/apiserver` (and a representative cloud path) is still refused; and the
  ONLY other module-graph movement is the two measured indirect bumps
  (`go-isatty v0.0.20 → v0.0.22`, `x/sys v0.46.0 → v0.47.0`) — any third movement fails the row.
- [ ] **AC17 — verified commit-boundary law (P6.V):** a named commit-boundary identity is
  Z3-verified on the pinned v0.30.0 binary and pinned in `REQUIRED_VERIFIED`; the gate's JSON
  checks (`verify.errors == 0`, `counterexample == 0`, required identity present and `verified`)
  all hold — guarding the silent `unknown sort` exit-0 trap; if the encodable-shape fallback was
  taken, the limitation is recorded in the doc and the named test-only law is in the 40+ named
  tests.

## Non-Vacuity — Named RED Mutation for Every Gate

| Gate | Named RED mutation (concrete edit) | Required red observation |
|---|---|---|
| AC1 `MUT-PROTO-OWNER` | Declare a local parallel A2A wire struct in `host/projection` and hand-format the agent-card JSON instead of using `protocol` types | wire-ownership source test REDs on the World-declared wire struct/envelope literal |
| AC2 `MUT-SESSION-UNION` | Return the union of all transitions for every session | low-capability session exact-set test REDs |
| AC3 `MUT-CARD-GLOBAL` | Generate the A2A card from the unfiltered registry | per-session card set-equality/session-change test REDs |
| AC4 `MUT-UNFILTERED-PROJECTION` | Re-enable raw v0.30.0 projection | surface test REDs on `std.io.writeBytes`/`writeBytes`, `exit`, or `submit_feedback` |
| AC5 `MUT-CALL-BYPASS` | Dispatch by skill name without re-running authorization | guessed-name call changes dispatch counter/store and denial test REDs |
| AC6 `MUT-DEFAULT-CAPS` | Map an unknown session to a process default | unknown-session card/call test REDs |
| AC6 `MUT-ALT-HEADER` | Make the resolver fall back to reading `X-World-Session` when `Authorization` is absent | carrier test REDs: a request bearing ONLY the rejected Arm-B header must fail closed, and under the mutation it resolves |
| AC6 `MUT-KEY-AS-SESSION` | Accept the static serve-api API-key value as a `Bearer` session credential | constraint-(i) test REDs: the key resolves to a session under the mutation while the test demands fail-closed |
| AC7 `MUT-SPLIT-SNAPSHOT` | Re-read capabilities after reading registry within one request | barrier-controlled epoch-consistency test REDs |
| AC8 `MUT-A2A-SHAPE` | Emit an invocation result as hand-built JSON that omits a protocol-required field instead of using `A2AResult` | frozen-fixture conformance test REDs on the missing/malformed field |
| AC9 `MUT-SECOND-OPEN` | Make projection call `store.Open` during mount | writer-lock/live-daemon test REDs with `WriterAlreadyActive` |
| AC10 `MUT-SKIP-SOCKET` | add `t.Skip` on listen failure | zero-skip CI assertion/source check REDs |
| AC11 `MUT-CLOUD-DEP` | add `cloud.google.com/go/storage` to projection imports | `TestDaemonDependencyAllowlist` reports it by name |
| AC12 `MUT-STARTUP-CACHE` | cache transition descriptors once at daemon startup | registry-head-change test returns stale set and REDs |
| AC13 `MUT-DROP-DEADLINE` | replace the propagated request context with `context.Background()` before one injected blocking dependency | bounded-wait fault-injection test exceeds the configured bound and REDs; mutation cleanup is terminated by the test harness |
| AC13 `MUT-COMMIT-BOUNDARY-LIE` | move the cancellation check from *before* coordinator acceptance to *after* it, so a cancelled-mid-commit request reports a definitive "not committed" | the after-boundary test REDs: a durable receipt exists under the invocation/idempotency ID while the caller was told the commit did not happen |
| AC13 `MUT-LEAK-CONN` | on client disconnect or deadline expiry with a blocked response write, cancel the Go context but never close the `http` connection | the `ConnState`/client-read-error assertion REDs — proving the test observes OS-level socket closure, not just logical cancellation |
| AC14 `MUT-DEADLINE-RELAX` | add a `ResponseController.SetWriteDeadline(time.Time{})` call to the A2A handler | the source-level assertion REDs naming the call site; the REST deadline regression is the same-run control that the D7 constants did not move |
| AC15 `MUT-TOOLCHAIN-REGRESS` | set ONE `ci.yml` site (or `go.mod`) back to `go1.25.6` | dependency resolution REDs with the measured floor message `requires go >= 1.26.6 (running go 1.25.6; GOTOOLCHAIN=go1.25.6)` |
| AC15 `MUT-CANARY-BLIND` | run the committed canary under deny-listed `go1.26.5` | repro prints `BUG: Field="" want "stateRoot"` (re-run 2026-08-25) — proving the detector still SEES the defect class, so `OK` under go1.26.6 is a measurement |
| AC16 `MUT-ALLOWLIST-ROOT` | replace the package-path entry with the module root `github.com/sunholo-data/ailang` | the narrowness test REDs: `internal/apiserver` is no longer refused by `disallowedDeps` |
| AC16 `MUT-FACADE-IMPORT` | import `github.com/sunholo-data/ailang/serveapi` (the facade) in `host/projection` | `TestDaemonDependencyAllowlist` REDs naming the intruding packages by name (measured scale: 476 across 86 roots) |
| AC17 `MUT-LAW-BREAK` | mutate the commit-boundary law body to permit a receipt without its journal intent (or a second receipt under one invocation ID) | `ai-check` reports `counterexample > 0` and `verify_ail.sh` REDs |
| AC17 `MUT-SORT-SILENT` | retype one law parameter to an ADT-bearing (`Proposal`-class) sort | `verify.errors > 0` in the JSON and `verify_ail.sh` REDs — proving the silent-exit-0 trap is guarded, not relied on |

Moved with the split: `MUT-PLAIN-JSON` (SSE framing), `MUT-LEAK-SSE-CONN` (stream-deadline
zombie connection), and `MUT-SSE-REST-DEADLINE` (route-local relaxation scope) — recorded in the
child doc as obligations that must reappear in its acceptance table when it becomes writable.

## Axiom Compliance

| Axiom | Score | Justification |
|---|---:|---|
| A1 Determinism | +1 | stable-ID ordering and one snapshot per request |
| A2 Replayability | +1 | calls use the normal recorded transaction path |
| A3 Effect Legibility | +1 | protocol I/O stays at the host boundary |
| A4 Explicit Authority | +2 | exact per-session card/call predicate; carrier fixed by D-WORLD-26, fail closed |
| A5 Bounded Verification | +1 | existing verify pipeline; no protocol bypass |
| A6 Safe Concurrency | +1 | daemon retains sole handle; immutable request snapshot |
| A7 Machines First | +2 | A2A schemas and structured upstream errors (the MCP surface follows via the child) |
| A8 Minimal Syntax | 0 | no language syntax change |
| A9 Cost Visibility | 0 | no new budget claim; benchmark scope deferred honestly |
| A10 Composability | +2 | standard A2A, upstream-owned wire contract; the MCP≡A2A equality obligation is preserved in the child |
| A11 Structured Failure | +1 | missing/expired/denied sessions fail in protocol form |
| A12 System Boundary | +2 | projection is an adapter, never kernel state |

**Net: +14; hard axioms A1/A3/A4/A7 are non-negative.**

## Premise Verification Log — revision round (live evidence, 2026-08-25)

Rows marked **[R]** were re-derived first-party in this revision round; rows marked **[C]** were
measured by the controller this same iteration (2026-08-25, each with a firing control) and are
cited with their commands; rows marked **[I]** are labelled inherited. Where a re-derivation
disagreed with the inherited value, the row says so explicitly.

| # | Premise | Command / evidence | Result |
|---|---|---|---|
| N1 [R] | `serveapi/protocol` exists at `v0.33.2`, the latest tag; `v0.34.0` does not exist | `git tag --sort=-v:refname \| head`; `git rev-parse v0.34.0` (fails) vs control `git rev-parse 'v0.33.2^{commit}'` → `63e7909f...` | Confirmed. Files: `interfaces.go`, `descriptor.go`, `envelope.go`, `a2a_wire.go` (+`descriptor_test.go`). **The charter's v0.34.0 mention is wrong** (row 5 itself already records this) |
| N2 [R] | `serveapi/protocol` is absent at `v0.33.1` → the pin resolves to v0.33.2 under D-WORLD-5's own "first tag containing the seam" discipline | `git ls-tree --name-only v0.33.1 serveapi/protocol/` → 0 entries; same-call control `git ls-tree --name-only v0.33.1 serveapi/` → 2 entries | Confirmed — D-WORLD-5's literal "v0.33.1" is superseded by its own condition; **not** an open decision |
| N3 [R] | `serveapi/protocol` imports ONLY the standard library | `git show 'v0.33.2:serveapi/protocol/<f>.go'`, import blocks of all four files | Confirmed: `context`, `encoding/json`, `errors`, `fmt`, `net/http`, `regexp`, `sort`, `strings`; zero non-stdlib |
| N4 [R] | `v0.33.2` `go.mod` declares `go 1.26.6` | `git show 'v0.33.2:go.mod' \| head -4` | Confirmed |
| N5 [R] | Toolchain floor, two-arm + control | ARM1 `GOTOOLCHAIN=go1.25.6 go get github.com/sunholo-data/ailang@v0.33.2` → **rc=1**, `requires go >= 1.26.6 (running go 1.25.6; GOTOOLCHAIN=go1.25.6)`; ARM2 `GOTOOLCHAIN=go1.26.6` → **rc=0**; control `GOTOOLCHAIN=go1.25.6 go get github.com/google/uuid@v1.6.0` → **rc=0** | Confirmed; arms differ as intended, control fired |
| N6 [R] | The `v0.33.2` pin is NOT transitively free | ARM2 on a copy of this repo's `go.mod`/`go.sum`: `go: upgraded go 1.25.6 => 1.26.6`, `upgraded github.com/mattn/go-isatty v0.0.20 => v0.0.22`, `added ...ailang v0.33.2`, `upgraded golang.org/x/sys v0.46.0 => v0.47.0` | Confirmed, verbatim tool output; both bumped roots already allowlisted |
| N7 [R] | Repo pins: `GOTOOLCHAIN: go1.25.6` at `ci.yml:21` and `:102`; `go.mod` declares `go 1.25.6` | `grep -n GOTOOLCHAIN .github/workflows/ci.yml`; `head go.mod` | Confirmed |
| N8 [R] | `verify_go.sh` deny-list + canary | `sed -n '214,224p' scripts/verify_go.sh` → case arms exactly `go1.26.0\|...\|go1.26.5`; canary re-run: `GOTOOLCHAIN=go1.26.6 go run .` → `OK`; known-bad control `GOTOOLCHAIN=go1.26.5` → `BUG: Field="" want "stateRoot"` | Confirmed; go1.26.6 is NOT deny-listed and the detector still sees the defect |
| N9 [R] | **DISAGREEMENT with the directive's "the rig's default `go` is already go1.26.6"** | `which -a go` → `/opt/homebrew/bin/go` only; `go version` → **go1.26.4** darwin/arm64; `go env GOTOOLCHAIN` → `auto`; in-repo (`go 1.25.6` module) `go env GOVERSION` → **go1.26.4** (deny-listed!); in a `go 1.26.6` module → **go1.26.6** | The BINARY is go1.26.4; `auto` selection reaches go1.26.6 only once `go.mod` says so. Consequence recorded in P6.T: pin `GOTOOLCHAIN` explicitly for local gate runs until P6.T lands |
| N10 [R] | Pristine dependency closure over BOTH gated patterns | `go list -deps ./host/daemon/... ./cmd/ailang-worldd/... \| sort -u \| wc -l` → **249**; same-call control: `grep -c github.com/google/uuid` → 1 | Confirmed 249 (the charter's iter-120 "pristine 244" is the older baseline; intervening landed work moved it — both agree post-pin is 250 and the non-stdlib delta is exactly 1) |
| N11 [C] | Probe import + unchanged allowlist REDs the gate naming exactly one intruder | `go test ./host/daemon/ -run TestDaemonDependencyAllowlist` with a `serveapi/protocol` probe import | `1 package(s) outside the zero-cloud allowlist ... github.com/sunholo-data/ailang/serveapi/protocol`; with the ONE package-path line added: green at **250**; sentinel control proved the diff can see an addition |
| N12 [R] | Allowlist shape + matcher semantics | `allowedDepModules` at `host/daemon/daemon_test.go:747` (11 module-root entries); matcher `d == m \|\| strings.HasPrefix(d, m+"/")` (~`:800-806`) | Confirmed: entries are path prefixes, so a PACKAGE-path entry admits only that package (+ subpackages beneath it) and refuses `internal/apiserver` — the matcher DOES distinguish; no extra matcher milestone needed, only the narrowness test |
| N13 [R] | `verify_ail.sh` floor | full gate run, `AILANG_BIN=/tmp/ailang-v0300/ailang ./scripts/verify_ail.sh` → PASSED | **10 required identities verified, 40 named tests pass, 9/9 world-package steps**; `REQUIRED_VERIFIED` at `verify_ail.sh:274-279`, 4 module entries, none a commit-boundary law (July's "4/9/14" floor is long stale) |
| N14 [R] | Commit-boundary Go surface is landed and public | `sed -n '24,30p' host/store/journal.go`; `grep -n bindCommitIntentTx host/store/store.go`; `grep -n 'GetReceipt\|GetEffectReceipt' host/store/journal.go`; `grep -n recoverCommitPending host/broker/recover.go` | `JournalIntent` at `journal.go:28`; **`bindCommitIntentTx` at `store.go:1025` — NOT `:1015` as the directive stated (line drift; one-line disagreement, re-derived value adopted)**; `GetReceipt :813`; `GetEffectReceipt :852`; `recoverCommitPending :126` |
| N15 [R] | Transition registry + broker session API landed | `ls host/transitionreg/` → `transitionreg.go`, `bind.go`, tests; `grep -n 'func NewSession' host/broker/broker.go` → `:87` | Confirmed — July prerequisite (2) discharged |
| N16 [R] | No HTTP session-credential carrier exists in `host/` (Open Decision 2 was live at the revision round) | `grep -rn 'Bearer\|Authorization' host/` → only a prose comment at `host/broker/approve.go:51`; same-call control: `net/http` grep hits 6 files in `host/broker` (reconcile+tests) and `host/daemon/daemon.go` | Confirmed absence at authoring time — the carrier has SINCE been chosen by the human: **D-WORLD-26 = ARM A** (split round; see Open Decisions and row N26) |
| N17 [R] | D7 constants/wiring line drift | `grep -n 'const (' host/daemon/daemon.go` → `:78`; `grep -n 'func newServer'` → `:619` | Constants `:78-97`, wiring `:619-628`; July's `:77-91`/`:409-414` refreshed throughout this doc |
| N18 [R] | No lockfile/module-pin gate collides with the `go.sum` movement | `grep -n 'go\.sum\|go mod verify\|go mod tidy' scripts/verify_go.sh` → 0 hits; same-call control `grep -c GOVERSION` on the same file → 1 | Confirmed (the planned `w-ail-gate-module-pin` pins the **`.ail`** module set, a different axis) |
| N19 [I] | Upstream closure guarantee is CI-enforced upstream (`check-protocol-closure`, 5-arm refusal self-test) | charter iter-120 row | Inherited, unverified here; first-party enforcement is N11's gate either way |
| N20 [I] | `Proposal`-sort Z3 limitation (`unknown sort`, silent exit 0) | charter (recorded v0.30.0 limitation; also `w-m1-ailang-hardening`) | Inherited; P6.V designs around it and MUT-SORT-SILENT proves the JSON guard |

### Split-round rows (2026-08-25, iteration 125 — all [R], re-derived first-party in the split session)

| # | Premise | Command / evidence | Result |
|---|---|---|---|
| N21 [R] | `#885` is OPEN with 0 comments (the child's blocking predicate holds) | `gh issue view 885 --repo sunholo-data/ailang --json state,title,createdAt,comments` → `state=OPEN`, `comments=0`, created `2026-08-25T16:37:44Z`, title "serveapi/protocol has no MCP dispatch — zero-cloud consumers can reach A2A but not MCP (follow-on to #764)"; same-call control `gh issue view 764 …` → `CLOSED`, 6 comments | Confirmed — the instrument can see closure and comments; the MCP half stays blocked |
| N22 [R] | Upstream latest release is still `v0.33.2`; no `v0.34.*` tag exists | `gh api repos/sunholo-data/ailang/releases/latest --jq .tag_name` → `v0.33.2`; `gh api …/matching-refs/tags/v0.34` → empty (rc=0); same-call control `…/tags/v0.33` → `v0.33.0`, `v0.33.1`, `v0.33.2` | Confirmed — the pin target is unchanged and nothing newer shipped |
| N23 [R] | The A2A existence proof + the dispatch gap, at tag `v0.33.2`, without any local upstream checkout | `gh api 'repos/sunholo-data/ailang/contents/serveapi/a2a_handler.go?ref=v0.33.2'` (and `mcp_handler.go`), base64-decoded: `a2a_handler.go` **180 lines**, import block = `context`, `encoding/json`, `fmt`, `log`, `net/http` + `serveapi/protocol`, `grep -c modelcontextprotocol` → **0**; `mcp_handler.go` **187 lines**, SDK import count **1**; `grep -c 'tools/list\|jsonrpc'` in `mcp_handler.go` → **0** while the same grep in `a2a_handler.go` → **2** | Confirmed — the A2A handler is buildable over stdlib+`protocol`; MCP dispatch is entirely SDK-side (the stronger form of the r3 objection: `mcp_handler.go` never string-matches method names at all) |
| N24 [R] | `verify_ail.sh` floor unchanged | `AILANG_BIN=/tmp/ailang-v0300/ailang ./scripts/verify_ail.sh` → `✓ verify gate PASSED: 10 required identities verified, 40 named tests pass`; `✓ world package gate PASSED: 9/9 steps performed non-zero work`; binary `--version` → v0.30.0, commit `e37b370` | Confirmed — AC10's floor is a live split-round measurement, not carried prose |
| N25 [R] | Every local line-number anchor this doc cites, at `2e44e3e` | `grep -n GOTOOLCHAIN .github/workflows/ci.yml` → `:21`,`:102`; `head go.mod` → `go 1.25.6`; `grep -n 'func NewSession' host/broker/broker.go` → `:87`; `grep -n 'const ('` → `daemon.go:78`, `grep -n 'func newServer'` → `:619`; `grep -n REQUIRED_VERIFIED scripts/verify_ail.sh` → `:274` (4 module entries; 10 identities by count); `grep -n allowedDepModules host/daemon/daemon_test.go` → `:747` (var), matcher `d == m \|\| strings.HasPrefix(d, m+"/")` at `:801`; `JournalIntent` → `journal.go:28`, `bindCommitIntentTx` → `store.go:1025`, `GetReceipt` → `:813`, `GetEffectReceipt` → `:852`, `recoverCommitPending` → `recover.go:126`; deny-list case arms exactly `go1.26.0\|…\|go1.26.5` | Confirmed, zero drift since the revision round (one refinement: the matcher's exact line is `:801`; the revision round cited `~:800-806`) |
| N26 [R] | D-WORLD-26 is RESOLVED = ARM A with both constraints | charter Human Decision Ledger row read first-party this session (`design_docs/world-mission.md:708`): answered `A` by Mark, attended, issue `#89`, `2026-08-25T19:06:41Z`, verbatim one-character comment, allowlist-enforced directive read; constraints (i) API-key-never-a-session and (ii) fail-closed travel with the answer | Confirmed — Open Decision 2 is CLOSED; zero forks remain in this doc |

## Premise Verification Log (historical, 2026-07-28 — superseded where the 2026-08-25 log disagrees)

| Premise | Command/evidence | Result |
|---|---|---|
| Pinned binary | `/tmp/ailang-v0300/ailang --version` | v0.30.0, commit `e37b370` |
| Flags | `/tmp/ailang-v0300/ailang serve-api --help` | `--mcp`, `--mcp-http`, `--a2a`, `--caps`, API-key flags, `--routes-only` present |
| S5 reference | `/tmp/ailang-v0300/ailang prompt` | rc=0; 2,535 lines captured |
| S5 builtins | `/tmp/ailang-v0300/ailang builtins list` | rc=0; 338 lines captured |
| Unfiltered sketch MCP | stdio initialize + `tools/list`, cwd `design_docs`, path `sketches/` | 27 tools; eight IO names plus `submit_feedback` observed |
| `--routes-only` | same stdio probe with `--routes-only` | only `submit_feedback`; std/io suppressed |
| `--caps ''` | same stdio probe with empty caps | unchanged 27-tool discovery |
| HTTP live rerun | `serve-api --mcp-http --a2a --port 8231 ...` | **BLOCKED BY SANDBOX**: `Error: listen tcp :8231: bind: operation not permitted` |
| HTTP framing/card | controller E2, 2026-07-28 | established: card 200; MCP HTTP body SSE-framed |
| Registry identity | inspection of `host/registry/registry.go` | only `world/epoch-registry/v1`; it is interpreter nomination metadata |
| REST mux/error | inspection of `host/daemon/daemon.go`, `handlers.go` | frozen routes and shared `APIError` confirmed |
| Daemon timeout constants | VERIFIED BY CONTROLLER (host/daemon/daemon.go:409-414, constants :77-91) | `ReadHeaderTimeout` 5s, `ReadTimeout` 30s, `WriteTimeout` 30s, `IdleTimeout` 120s; D7 bound-constant block is frozen |
| `http.ResponseController` available | VERIFIED BY CONTROLLER (`go.mod:3`) | `go 1.26.4` — well past the Go 1.20 floor for `http.NewResponseController` |
| Transition registry landed? | VERIFIED BY CONTROLLER (repo-wide search for `[Tt]ransition[ -]?[Rr]egistry`) | matches in `design_docs/` ONLY — **zero** hits in `host/`, `world/`, `cmd/`. The clause-3 prerequisite is real, not defensive |
| **Commit-boundary contract** | **"VERIFIED" here meant Go-tested (`w-store-durability` SD.B/SD.C, landed) — NOT Z3-verified in this repo's `REQUIRED_VERIFIED` sense; the 2026-08-25 log (N13/N14) and milestone P6.V carry the correction** | `store.AppendIntent` durably records the stable invocation ID; `Commit.InvocationID` binds the request and writes its outcome in the same transaction; `GetReceipt` exposes the three-state durable receipt. SD.C's real-process crash tests kill at the named pre-outcome boundary and prove world + receipt atomicity, including the split-transaction RED mutation. |
| CI jobs | inspection of `.github/workflows/ci.yml` | exactly `ailang-verify` and `go-verify`; benchmark smoke in `go-verify` |

### Gate output tails

`AILANG_BIN=/tmp/ailang-v0300/ailang ./scripts/verify_ail.sh`:

```text
✓ 4/4 required world/ identities verified across 9 module(s)
✓ all 14 required named tests pass (failed_tests=0)
✓ verify gate PASSED: 4 required identities verified, 14 named tests pass
```

`AILANG_BIN=/tmp/ailang-v0300/ailang ./scripts/verify_go.sh`:

```text
── go build ./...
── go test ./... -count=1
FAIL github.com/sunholo-data/ailang-world/cmd/ailang-worldd
FAIL github.com/sunholo-data/ailang-world/host/daemon
Error: listen tcp 127.0.0.1:0: bind: operation not permitted
```

The non-socket packages, including `host/replay` (13.299s), passed. The full Go gate is
**UNVERIFIED in this sandbox**, not green.

No new `.ail` exists, so new-sketch Z3 obligations = **0** and new-sketch tests = **0**.
The repo sweep independently proves the unchanged existing totals above.

## Upstream Findings

### UF-P6-1 — v0.30.0 cannot expose an exact caller-supplied/session-scoped MCP+A2A surface — **RESOLVED BY DELIVERY (2026-08-25)**

**Resolution:** the request below was executed as `sunholo-data/ailang#498` (lane B, M1–M3 in
v0.33.0/v0.33.1) and, after iter-90's closure audit failed the facade (476 disallowed packages),
as `#764` — a protocol-only, stdlib-only package — delivered in **`serveapi/protocol` at
`v0.33.2`**. The reproduction and hypothesis are retained below as the historical record of what
was asked and why; the observed v0.30.0 behavior itself is unchanged (Decision 1's (a)/(b)
rejections still rest on it).

**Version:** AILANG v0.30.0, commit `e37b370`.

**Minimal reproduction actually run (stdio leg, no socket required):**

1. From repo root, run `ailang prompt` and `ailang builtins list` (completed above).
2. From `design_docs/`, start the pinned binary with `serve-api --mcp sketches/`.
3. Send MCP `initialize`, `notifications/initialized`, and `tools/list`.
4. Repeat with `--routes-only`; repeat with `--caps ''`.

**Observed:** unfiltered and empty-caps runs each list 27 tools, including `eprintln`, `exit`,
`flush`, `print`, `printErr`, `println`, `readLine`, `writeBytes`, and `submit_feedback`.
`--routes-only` removes the module/stdlib exports but still lists `submit_feedback`. The help
surface offers process-wide caps and one static key, with no caller-supplied session/discovery
hook.

**Hypothesis:** the built-in feedback tool is appended outside route filtering, and discovery is
constructed from loaded exports rather than a caller-supplied provider. Source inspection of the
upstream implementation was not available in this worktree, so those causes are hypotheses.

**Requested upstream resolution:** export the existing serving machinery with the callback
behavior in Decision 2. Do not add World persistence upstream.

**Controller action:** reproduce first, then file in `sunholo-data/ailang` and send the required
mission-control message. This design agent did not mutate external issue/message state.

### UF-P6-2 — `serveapi/protocol` has no MCP JSON-RPC dispatch — FILED as `#885` (2026-08-25), carried by the child doc

The round-3 finding (see the split-round quorum entry) was routed upstream on 2026-08-25 as
[`sunholo-data/ailang#885`](https://github.com/sunholo-data/ailang/issues/885) — "serveapi/protocol
has no MCP dispatch — zero-cloud consumers can reach A2A but not MCP (follow-on to #764)" —
carrying the measured closure numbers and the `a2a_handler.go` existence proof. This is
D-WORLD-5's prescribed default executing as written (a disallowed graph asks upstream, never a
broad relaxation), the same route that produced `#764` → `v0.33.2`, and it is not a new human
ask. The full problem statement, evidence table, and runnable blocking predicate live in
[`w-mcp-dispatch-projection.md`](w-mcp-dispatch-projection.md); this parent tracks it only as
out-of-scope context (re-measured this round: OPEN, 0 comments — premise row N21).

## Open Decisions

**NONE — zero forks remain open in this doc**, which is precisely what quorum round 3's second
objection demanded of an executable design. The historical rows and their dispositions:

1. **Exact upstream module/version — ANSWERED BY MEASUREMENT + D-WORLD-5's own condition.**
   `v0.33.2` is the first (and latest) tag containing `serveapi/protocol` (premise rows N1/N2,
   re-confirmed N22); D-WORLD-5's "first tagged release containing the seam" resolves to it
   mechanically. Its literal mention of `v0.33.1` predates the protocol-only delivery; no re-ask
   needed.
2. **Session credential carrier — CLOSED BY THE HUMAN: D-WORLD-26 = ARM A (2026-08-25T19:06:41Z).**
   Clause 3 landed WITHOUT an HTTP credential convention (premise row N16: no
   `Authorization`/`Bearer` convention anywhere in `host/` beyond one prose comment), so the fork
   was real; round 3's `gemini-3-1-pro` objection named it, and it was surfaced to Mark as
   `D-WORLD-26`. Mark answered **A** (attended, issue `#89`, verbatim one-character comment `A`;
   allowlist-enforced directive read; ledger row read first-party this session — premise row
   N26): the carrier is the standard `Authorization: Bearer <session-credential>` header, read by
   `protocol.SessionResolver.ResolveSession(ctx, *http.Request)`. Arm B (`X-World-Session`) is
   REJECTED. Two constraints travel with the answer and are binding on Decision 3, AC6, and the
   resolver's contract: **(i)** the static `serve-api` API key remains forbidden as a session
   model (iteration 24 measured it process-wide, so it cannot represent a session) — a `Bearer`
   value here is a SESSION credential and never an API key; **(ii)** clause 3 still binds — the
   resolver fails closed on an absent, malformed, or unknown credential, never degrading to an
   unauthenticated surface.
3. **Transition descriptor schema — ANSWERED BY LANDING.** `host/transitionreg.Descriptor`
   (stable ID, schema, description, effect requirements) is used verbatim, projected through
   `protocol.ToolDescriptor`/`CallerSurface`; no projection-owned schema. This is now Decision 3
   text, not a decision.
4. **Disallowed upstream graph — ANSWERED BY DELIVERY (for the A2A half) and ROUTED UPSTREAM
   AGAIN (for the MCP half).** The default this row prescribed ("ask upstream for a protocol-only
   module") was executed as `#764` and delivered as `serveapi/protocol`; the measured admitted
   graph is one stdlib-only package (premise rows N3/N10/N11). The same default, applied to the
   round-3 dispatch gap, produced `#885` — the CHILD doc's blocker, not an open decision here.

## Quorum verification log

**Gate**: pick-time / creation-time design quorum (`ailang design-quorum`, reject-by-default
synthesis: any present reviewer or the controller rejects → BLOCKED). Reviewers
`gpt5-6-sol` + `gemini-3-1-pro`, `--max-cost-usd 0.25`, plus the controller's in-session verdict.
Artifacts: `.ailang/state/mission-quorum/w-mcp-projection-2026-07-27T23-40-18Z.json` (r1) and
`…-2026-07-27T23-46-08Z.json` (r2).

**GENERATOR = JUDGE COLLISION, FLAGGED IN BOTH ROUNDS.** This doc (and its r1→r2 revision) was
authored by the rotation designer `codex:gpt-5.6-sol`, so the `gpt5-6-sol` reviewer seat is a
**SELF-review and is not independent**. It was retained rather than excluded because reject-by-default
synthesis means a self-*pass* cannot manufacture a PROCEED, so keeping the seat can only add
objections, never remove them. The independent rejectors were `gemini-3-1-pro` and the controller
(Opus). This is worth recording as evidence, not just as a caveat: **the self-seat did not
rubber-stamp itself in either round** — it produced the strongest objection both times, and in
round 2 it was the *only* reviewer still rejecting.

### Round 1 — BLOCKED (both reviewers present; `metered=$0.063444`)

| Reviewer | Verdict | Objection (summarised) | Disposition |
|---|---|---|---|
| `gpt5-6-sol` (self-seat) | **reject** | No bounded-wait contract anywhere on the projection path: no timeout source, maximum, context-propagation requirement, cleanup rule, protocol error mapping, acceptance test or RED mutation for a stalled resolver/registry/broker/verifier/client. "Broker unavailability returns an error" is insufficient. | **Fix applied in full, in the reviewer's own terms** → new `Decision 6`, `AC13`, and the reviewer's own mutation name `MUT-DROP-DEADLINE`. |
| `gemini-3-1-pro` | **reject** | Conflict-Surface gap: mounting upstream SSE-framed MCP HTTP handlers on the existing worldd daemon ignores that a REST daemon's `ReadTimeout`/`WriteTimeout` abruptly terminate long-lived SSE streams. | **Fix applied** → the requested `HTTP Server Timeouts vs SSE` Conflict-Surface entry, `AC14`, `MUT-SSE-REST-DEADLINE`. |
| controller (Opus) | pass | Design direction sound; recorded three non-blocking concerns (P6.A is documentation-only; the item is transitively blocked behind the parked `w-effect-broker-m3`; the doc does not assess upstream appetite for exporting `internal/apiserver` under PROGRAM.md's frozen-core discipline). | Carried into the report. |

**The controller ran gemini's "verify" step first-party, and it collapsed gemini's two-branch fix to
one branch.** Gemini offered "use `http.ResponseController` … **or** explicitly document that the
current daemon lacks global timeouts". The second branch is **FALSE for this repo**: VERIFIED BY THE
CONTROLLER at `host/daemon/daemon.go:409-414`, wiring constants declared at `:77-91` —
`ReadHeaderTimeout` 5s, `ReadTimeout` 30s, `WriteTimeout` 30s, `IdleTimeout` 120s, the D7
bound-constant block that `w-worldd-m2` milestone A2 ratified and FROZE. So the revision was
directed to the `ResponseController` branch only, scoped to `/mcp/`, with the D7 constants and every
REST `/v1/*` path explicitly byte-unchanged — and with the follow-up the reviewer did not ask but
the freeze demands: *what bounds the connection once its write deadline is relaxed?* (answer: a
second explicit finite stream-lifetime maximum; `IdleTimeout` is correctly NOT that bound, since it
governs idleness between requests rather than an active handler).

### Round 2 — BLOCKED, 1 of 2 reviewers (both present; `metered=$0.072912`)

| Reviewer | Verdict | Objection (summarised) |
|---|---|---|
| `gemini-3-1-pro` | **PASS** | Non-blocking: delegating SSE socket lifecycle to a not-yet-written upstream handler risks zombie TCP connections if upstream mishandles the injected cancellation context. Proposed fix: assert **OS-level** socket closure, not just logical context cancellation. |
| `gpt5-6-sol` (self-seat) | **reject** | Decision 6/AC13 promised an impossible guarantee: that expiry after commit *begins* still yields no store/log mutation. An HTTP context can expire mid-atomic-commit, so without a verified atomic "not-started versus committed" contract the commit may succeed while the caller observed cancellation. |
| controller (Opus) | pass | Both r1 fixes verified as applied in the reviewers' own terms; estimate moved honestly 1.0d → 1.2d rather than being absorbed. |

### Disposition — NARROW-REFINEMENT CARVE-OUT APPLIED (bounded 2nd revision, controller)

Both remaining items pass both limbs of the carve-out: **(a)** each carries a concrete,
reviewer-authored `proposed_fix` (`gpt5-6-sol` supplied verbatim replacement prose for the
cancellation paragraph, the AC13 revision, and the premise row; `gemini-3-1-pro` supplied the
`ConnState`/client-read-error assertion), and **(b)** neither disputes the design DIRECTION — path
(c), upstream-owned codecs, session-scoped projection inside worldd and the bounded-wait decision
all stand. The objection is a **truthfulness-of-claim** defect: an acceptance criterion asserting a
guarantee the system cannot provide. Applied:

1. **Decision 6's commit-boundary paragraph replaced with the reviewer's verbatim contract** (quoted
   as a block quote so the substitution is auditable), with the over-strong claim it replaces
   stated explicitly rather than silently deleted.
2. **AC13 now tests cancellation on BOTH SIDES of the boundary** — before acceptance → no durable
   mutation; after acceptance → exactly one recoverable, queryable/replayable receipt under the
   stable invocation/idempotency ID; never a definitive "not committed" while the outcome is
   unknown. It also carries gemini's OS-level socket-closure assertion.
3. **Two new premise rows**, including `Commit-boundary contract` marked **UNVERIFIED —
   PREREQUISITE**: no landed API exposes these semantics, so the row records the gap instead of
   inventing a mechanism.
4. **P6.B gains the commit-boundary contract as an explicit blocker**, in Decision 6, the Design
   Freeze and the milestone text.
5. **Three new RED mutations**: `MUT-COMMIT-BOUNDARY-LIE`, `MUT-LEAK-SSE-CONN`, and the reviewer's
   own `MUT-DROP-DEADLINE` from r1.

This SATISFIES the objections; it is not force-passing. Standing rule 2 still forbids proceeding
over a contested design DIRECTION, and no reviewer contested one. **No third quorum round was run**
— the carve-out is one bounded revision, not a re-litigation.

**Routing consequence, recorded honestly**: the carve-out normally routes a doc straight to
sprint-planner. It did **not** here, and the reason is the doc's own conclusion rather than the
quorum — P6.B is blocked on an upstream serving seam plus two unlanded local prerequisites, so a
sprint plan would be a plan for work that cannot start. Only **P6.A** is actionable now, and P6.A is
"file the upstream finding + land this record", which the controller executed inline this iteration.

### Revision round — 2026-08-25 (premise falsified by upstream delivery)

**Why this round exists.** The doc's chosen architecture — path (c), "request a narrow public
serving seam over `internal/apiserver`" — rested on "upstream exposes nothing public." Upstream
shipped exactly the requested seam (`serveapi/protocol` at `v0.33.2`, via `#498` → `#764`), so
the clearing of the blocker IS the falsifier of the central decision. Per the standing rot rules
this required re-derivation, not patching. Human rulings executed, not re-opened: **D-WORLD-5**
(ARM A — import upstream, pinned; Mark, attended, 2026-08-17) and **D-WORLD-25** (finish item 14
first; item 14 completed iter-123, so row 5 is live).

**What changed, for the next quorum's delta read:**

1. Status/scheduling: BLOCKED → UNBLOCKED; all July prerequisites discharged except the one-word
   residual "verified," which is now milestone **P6.V** in this doc (pure-core `world/*.ail` law
   + `REQUIRED_VERIFIED` pin) rather than an external wait.
2. Decision 1: (a)/(b) rejections retained (evidence unchanged on the pinned binary); the
   "nothing public upstream" premise marked FALSIFIED; path (c) marked satisfied-by-delivery.
3. Decision 2: rewritten from a behavioral request to the measured delivered surface, pinned at
   `v0.33.2` (`v0.33.1` measured NOT to contain the package; `v0.34.0` does not exist). Scope
   named: the handlers/`CallbackRunner` are deliberately outside `protocol` (facade closure = 476
   disallowed packages, iter-90) — **World writes its own handlers and callback-bounding**; this
   is D-WORLD-5 executing, not new scope.
4. Milestones restructured: P6.A (done) → **P6.T toolchain floor go1.25.6→go1.26.6 first,
   independently mergeable** → P6.D pinned dependency + ONE package-path allowlist line → P6.V
   verified commit-boundary law → P6.B projection with World-owned handlers. Each independently
   CI-green.
5. New AC15–AC17 with named RED mutations (`MUT-TOOLCHAIN-REGRESS`, `MUT-CANARY-BLIND`,
   `MUT-ALLOWLIST-ROOT`, `MUT-FACADE-IMPORT`, `MUT-LAW-BREAK`, `MUT-SORT-SILENT`); AC1/AC10/AC14
   re-derived; a fresh dated premise log (N1–N20) with two loud disagreements recorded
   (`bindCommitIntentTx` at `store.go:1025`, not `:1015`; the rig `go` BINARY is go1.26.4, not
   go1.26.6 — selection semantics reconcile the directive's claim).

**Both carried objections re-verified under the new architecture:**

- `gpt5-6-sol` (r2 reject, carve-out applied): the reviewer's verbatim commit-boundary contract
  is RETAINED unchanged in Decision 6 and AC13 still tests both sides of the boundary. The
  revision *strengthens* its standing: the contract's Go surface is now landed public API
  (`JournalIntent`/`bindCommitIntentTx`/`InvocationID`/`GetReceipt`), and the previously
  unverifiable half is scoped as P6.V with its own gate and RED mutations. The carve-out
  disposition holds.
- `gemini-3-1-pro` (r2 pass with live objection): SSE zombie-connection risk previously landed on
  a hypothetical upstream handler; under the new scope the `/mcp/` handler is World-authored, so
  the objection lands on first-party code and is ADOPTED as a design obligation (Decision 6);
  AC13's OS-level `ConnState`/client-read-error closure assertion and `MUT-LEAK-SSE-CONN` are
  retained verbatim.

### Round 3 — BLOCKED (2026-08-25, both reviewers PRESENT; `metered=$0.1658`)

Artifact: `.ailang/state/mission-quorum/w-mcp-projection-2026-08-25T16-33-38Z.json`.
Cap raised to `--max-cost-usd 0.25` **deliberately and pre-emptively**: the doc grew 641 → 974
lines in this revision, and a reviewer dropping out on budget immediately after a doc grows is the
skill's named self-selecting trap (the eye that closes is the one whose objection drove the
revision). `absent_reviewers` is **EMPTY** — this is a full-strength N=2 block, not an N−1 degrade.
Note the July generator=judge collision does **not** apply to this round: the revision was authored
by the Fable designer, so the `gpt5-6-sol` seat is independent here for the first time.

- **`gpt5-6-sol` → REJECT** ($0.1169): *"The selected upstream seam is insufficient for the design's
  own no-codec rule … no MCP request parser or handler for JSON-RPC method dispatch, initialization,
  `tools/list`, or `tools/call`. A World-authored `/mcp/` handler therefore appears forced to
  implement MCP/JSON-RPC parsing and dispatch locally, directly contradicting P1, the Design Freeze,
  and AC1."*
  **CONFIRMED FIRST-PARTY BY THE CONTROLLER, and it is a DIRECTION-level objection, so the
  narrow-refinement carve-out does not apply.** Measured this iteration (see the SDK arm in the
  status header and premise rows). Note the instrument correction that sharpened it: the
  controller's first probe grepped `serveapi/mcp_handler.go` for `"tools/list"`/`"jsonrpc"` and got
  **0** — which looked like a refutation and was in fact the known-positive control FAILING TO FIRE.
  The right reading is that `mcp_handler.go` does not string-match method names *at all* because it
  delegates the entire dispatch to the SDK, which is a stronger form of the reviewer's claim than
  the one the grep was written to test.
- **`gemini-3-1-pro` → REJECT** ($0.0489): *"The document declares itself 'UNBLOCKED' and
  'executable NOW', yet explicitly leaves Open Decision 2 ('Session credential carrier') as an
  unclosed fork … An executable design must close its own forks before being handed to an
  executor."* **VALID and narrow.** OD2 is a genuine session-authority decision in a repo where
  clause 3 is inviolable, so the controller does not settle it (standing rule 8); it is surfaced to
  the human as `D-WORLD-26`. The status header above is corrected in this iteration so the doc no
  longer claims executability it does not have.
- **controller (in-session) → pass.** Recorded for completeness and **overridden by the two
  rejections** under reject-by-default synthesis. The controller's verdict was cast on the
  premise-work axis (every load-bearing claim re-derived first-party with firing controls, two of
  the controller's own numbers corrected by the designer) and did not weigh the dispatch gap, which
  neither the controller nor the designer had measured at that point. Recorded as a friction: a
  `pass` cast before the hardest surface was probed.

**Objection-surface tracking (skill rule: from round 3 on, record WHICH surface each round's
objections name).** R1 — three objections across three surfaces. R2 — commit boundary
(`gpt5-6-sol`) + SSE socket lifetime (`gemini-3-1-pro`, pass-with-objection). R3 — **MCP dispatch
seam** + **status/fork self-consistency**. The objections have not localised onto one surface, but
the R3 finding is itself a decomposition signal of the kind that rule exists to catch: two
consumers (MCP, A2A) with different readiness bundled under one doc, the harder one holding the
easier one hostage. Disposition is therefore **SPLIT** — a controller routing call, explicitly
*not* `needs-human-review` — with the surviving objection carried verbatim into the MCP doc's
opening problem statement when it is authored.

### Split round — 2026-08-25 (iteration 125) — SPLIT APPLIED; no quorum was run this round

**Disposition executed.** Round 3's controller routing call (SPLIT — explicitly not
`needs-human-review`) is applied in this revision. No reviewer ran against this text: the
controller re-quorums the reduced parent after this revision lands, and the child doc must be
quorumed at pick time when it unblocks (it is authored NOT quorum-cleared, and says so).

**Where each round-3 objection went:**

- `gpt5-6-sol` (DIRECTION-level, confirmed first-party) — **carried VERBATIM into the child**,
  [`w-mcp-dispatch-projection.md`](w-mcp-dispatch-projection.md), as its opening problem
  statement, with the measured both-routes-closed evidence and a runnable blocking predicate
  (the `#885` state/comment read plus the upstream release tag, each with a same-call control —
  re-measured this round: OPEN, 0 comments; latest release still `v0.33.2`; premise rows
  N21/N22). Everything in this parent that depended on MCP dispatch moved with it: the `/mcp/`
  handler and endpoint, `initialize`/`tools/list`/`tools/call` dispatch, dispatch-bound envelope
  framing, SSE stream lifetime and the route-local `ResponseController` relaxation, old AC3's
  cross-surface MCP≡A2A equality, old AC8's SSE-framing assertion, old AC14 (SSE/REST deadline
  separation), and mutations `MUT-PLAIN-JSON`, `MUT-LEAK-SSE-CONN`, `MUT-SSE-REST-DEADLINE`.
- `gemini-3-1-pro` (narrow: "an executable design must close its own forks") — **CLOSED by the
  human**: `D-WORLD-26` = ARM A (Mark, attended, issue `#89`, `2026-08-25T19:06:41Z`; ledger row
  read first-party, premise row N26). Open Decision 2 is closed with the ruling's two constraints
  binding Decision 3 and AC6, plus two new named mutations (`MUT-ALT-HEADER`,
  `MUT-KEY-AS-SESSION`) so the closure is testable, not decorative. This doc now has zero open
  forks.

**What remains here** is the reviewer-clean remainder: P6.A (done), P6.T, P6.D, P6.V, and
P6.B-A2A — with `a2a_handler.go` (180 lines, stdlib + `protocol`, 0 SDK imports; premise row
N23) as the existence proof that the retained scope needs nothing from either closed route. The
AC renumbering under the split: old AC3/AC8/AC14 partially moved to the child (their MCP halves);
the numbers are re-used here for A2A-scoped criteria (card-tracks-session, A2A wire conformance,
no-deadline-tampering respectively), and the moved halves are enumerated in Deferred Scope and
the child.

**Objection-surface tracking (skill rule, from round 3 on):** R1 — bounded waits + SSE/REST
timeout conflict + (controller, non-blocking) upstream appetite: three surfaces. R2 — commit
boundary (`gpt5-6-sol`) + SSE socket lifetime (`gemini-3-1-pro`, pass-with-objection). R3 — MCP
dispatch seam + status/fork self-consistency. Split round — no reviewers ran, so no new
objection surface; of R3's two surfaces, the dispatch seam LEFT this doc (child, blocked on
`#885`) and the fork-consistency surface is CLOSED by a human answer (D-WORLD-26). The surfaces
that remain live in this parent for the next quorum to probe are exactly the retained ones: A2A
wire conformance, session authority (now carrier-fixed), the commit boundary (P6.V), and the
dependency/toolchain gates.

### Round 4 — BLOCKED, but the objections LOCALISED and one reviewer FLIPPED TO PASS (2026-08-25, iteration 125; both reviewers PRESENT, `absent_reviewers` EMPTY; `metered=$0.213298`)

Artifact: `.ailang/state/mission-quorum/w-mcp-projection-2026-08-25T21-25-44Z.json`. Cap raised to
`--max-cost-usd 0.35` because the doc grew again in the split round (974 -> 1204 lines) — the
skill's named self-selecting budget trap, where the reviewer that drops out on budget is
systematically the one whose objection drove the revision. It did not fire: both external
reviewers are present and the block is full-strength N=2, not an N-1 degrade.

- **`gemini-3-1-pro` -> PASS** ($0.063068) — **the first pass this document has received in four
  rounds.** Its only remark is non-blocking and narrow: P6.V's milestone text said
  *"raising the floor from 10 identities"* unconditionally, which an executor could misread as a
  strict mandate even where the encodable-shape fallback fires and the floor correctly stays at
  10. **Applied verbatim by the controller** in this same iteration (P6.V body + merge criterion,
  now explicitly conditional); AC10/AC17 already stated the conditional correctly.
- **`gpt5-6-sol` -> REJECT** ($0.15023) — *"The session-authority boundary is not executable: the
  doc equates `host/broker.NewSession(store, episodeID, grants, registry)` with an API that
  resolves an opaque Bearer credential, but N15 verifies only a session constructor. No landed
  credential lookup, credential-to-episode/grants mapping, expiry/revocation source, or
  authentication API is identified. D-WORLD-26 selects the HTTP carrier only; it does not supply
  these missing semantics."*
- **controller (in-session) -> pass**, overridden by the rejection under reject-by-default.

**THE OBJECTION IS CONFIRMED FIRST-PARTY AND SHARPENED BY THE MEASUREMENT** (skill rule 3f — a
reviewer's objection is a claim too, and a premise objection is the controller's to *measure*, not
to forward). Audited at `2e44e3e` over `host/` (96 Go files; scope asserted with `test -d`), all
in one call, counts are of matching LINES:

| Probe | Result | Reading |
|---|---|---|
| `grep -rn --include='*.go' 'Bearer' host/` | **0** | no inbound bearer-credential handling exists |
| `grep -rn --include='*.go' 'Authorization' host/` | **1** | a prose comment in `approve.go`, not a code path (matches premise row N16) |
| `grep -rn --include='*.go' 'Credential' host/` | **128** | **all OUTBOUND** — `RegistryCredentialProvider` / `FileRegistryCredentialProvider` / `AssertNoAmbientRegistryCredential` (`host/broker/credential.go`), i.e. a credential World *presents* to an upstream registry |
| `grep -rn --include='*.go' 'Authenticate' host/` | **29** | **all evidence-envelope** — `AuthenticatedEnvelope` and its codec (`host/evidence/`), i.e. signed evidence, not HTTP session auth |
| `grep -rnE 'func .*(GetSession\|LookupSession\|ResolveSession\|SessionByID\|FindSession)' host/ cmd/` | **0** | **nothing resolves a session by string** anywhere in the repo |
| known-positive control, SAME scope, same command shape: `'Session' host/` | **181** | the instrument sees the surface it is searching |
| negative control, fresh literal, same scope | **0** | the instrument's zeros are measurements |

`host/broker/broker.go:87` reads `func NewSession(s *store.Store, episodeID string, grants []Capability, registry Registry) *Session` — the grants are an **argument**. Nothing in the repo decides what grants a credential carries. So `D-WORLD-26` settled the **envelope** (`Authorization: Bearer`) and the **contents** question is untouched: who mints a session credential, where the credential -> (episode, grants) mapping is stored, and what expires or revokes it. The reject is correct, and it is a **strictly larger** finding than "the doc under-specifies a resolver".

**DISPOSITION — SPLIT AGAIN, and this is the skill's own decomposition signal firing exactly as
written, not a judgement call.** The rule: *objections that localise onto one surface while
another reviewer starts passing mean the doc's SCOPE is wrong — it bundles surfaces with
different correctness bars, and the hardest one is holding the others hostage.* Round 3 was two
objections on two surfaces; round 4 is `gemini-3-1-pro` **passing** while `gpt5-6-sol` rejects on
**one** surface — session authority — which is a property of **P6.B-A2A alone**. `P6.T`, `P6.D`
and `P6.V` have never drawn an objection in any of the four rounds and have no session-authority
dependence whatever.

And by the rule's clause (c): this defect is one the doc **fails to fix**, not one it
**introduces**. Clause 3 landed the session/capability model deliberately in-process; no
HTTP-facing session authority was ever built, and P6.B-A2A is simply the first thing that needs
one. **A pre-existing defect surfaced by a reviewer is a QUEUE ROW, not a revision** — so it is
filed on the first-party evidence above rather than absorbed by growing this document a third
time. It is explicitly **NOT** `needs-human-review`: a split is a controller routing call, and
filing it as a decision would manufacture an ask Mark does not have (standing rule 8).

**Objection-surface tracking (rounds 3+).** R1 — three objections, three surfaces. R2 — commit
boundary + SSE socket lifetime. R3 — MCP dispatch seam + status/fork self-consistency. **R4 —
session authority ONLY, with the other reviewer passing.** The trend is monotone convergence:
three surfaces, then two, then one-plus-a-nit, and each departing surface left as its own doc or
row rather than as a deletion.

**Round count, said out loud because the skill requires it.** This document is at **four quorum
rounds** plus a carve-out revision. That is data about this loop's *scoping* rather than about
this document: `w-mcp-projection` has been bundling surfaces of very different readiness since
July, and each round has cost real metered dollars to discover one more of them. Only the human
can act on that pattern, so it is surfaced in the iteration report rather than acted on here.


## Related Documents

- [w-mcp-dispatch-projection.md](w-mcp-dispatch-projection.md) — the SPLIT child: carries the
  MCP dispatch half and the round-3 objection verbatim; blocked on `ailang#885`; NOT
  quorum-cleared
- [world-mission.md](../world-mission.md) — clause 6, queue row 5, D-WORLD-5/D-WORLD-25/
  **D-WORLD-26** (the carrier ruling this doc executes), premise measurements this doc consumes
- [coding-standards.md](../coding-standards.md) — S1–S6
- [DESIGN.md](../DESIGN.md) — §3.7, §14, §17
- [w-worldd-m2.md](../implemented/w-worldd-m2.md) — shipped daemon, REST freeze, D7 constants
- [w-effect-broker-m3.md](../implemented/w-effect-broker-m3.md) — landed session/capability/broker
  (moved to implemented/ since the July draft linked it from planned/)
- [w-store-durability.md](../implemented/w-store-durability.md) — the commit-boundary Go surface
  P6.V puts a proof over
- [w-race-gate-blindspot.md](../implemented/w-race-gate-blindspot.md) — the toolchain canary P6.T
  re-runs (committed repro + controls)
- [w-m1-ailang-hardening.md](../implemented/w-m1-ailang-hardening.md) — the encodable-shape /
  test-only fallback pattern P6.V reuses if the `Proposal`-sort limitation bites
- [worlddapi.ail](../sketches/worlddapi.ail) — frozen REST-law precedent; unchanged
