# w-mcp-projection — Clause-6 Enabling Milestones (SPLIT parent; MCP dispatch and A2A projection carved out)

**Status**: Planned — **CLEARED TO SPRINT. Round 5 (2026-08-26, iteration 126; both reviewers
PRESENT, `absent_reviewers` EMPTY, `metered=$0.219099`) BLOCKED on two narrow objections, and BOTH
HAVE BEEN APPLIED VERBATIM under the narrow-refinement carve-out — one of them after the
controller re-measured it first-party and found it MATERIAL.** (1) `gemini-3-1-pro` refuted
premise row **N7** as a narrowed search; the repo-wide sweep it prescribed found **two live pin
sites P6.T never named** — `actions/setup-go@v5`'s `go-version: '1.25.6'` at `ci.yml:28` and
`:109` — so P6.T now moves **four** `ci.yml` sites plus `go.mod`, and N7 carries the full sweep
with its controls. (2) `gpt5-6-sol` rejected P6.D as speculative core growth: split #2 had moved
its only real consumer to the A2A child one iteration earlier, leaving a dead compile anchor.
**P6.D, `host/daemon/protocol_use.go`, AC16 and its two mutations are removed from this parent** —
dependency admission is now **atomic with its first real consumer** and travels with whichever
child unblocks first, specification intact. **WHAT REMAINS: `P6.A` (record), `P6.T` (~0.1d) and
`P6.V` (~0.3d) — ~0.4d, no dependency admission, no handler code, nothing external on the critical
path.** Clause 6 is partitioned across THREE docs (see the clause-6 scope statement below). Full
disposition, per-round surface table and round count: `## Quorum verification log` →
`### Round 5`.

**Prior head text follows (the split-#2 head; accurate when written — its "zero objections across
all four rounds" was true of rounds 1–4 and round 5 then found two narrow ones).** Planned — **SPLIT #2 APPLIED 2026-08-26 (iteration 126): ROUND 4'S BLOCKING
OBJECTION IS DISCHARGED BY SPLIT, NOT ANSWERED — `P6.B-A2A`, the whole A2A projection surface,
the session-authority objection it drew, and `D-WORLD-26`'s operative force are carved out into
the child [`w-a2a-session-projection.md`](w-a2a-session-projection.md), honestly BLOCKED on
charter queue row 39 `w-session-authority`. WHAT REMAINS HERE — `P6.A` (record), `P6.T`, `P6.D`,
`P6.V` — HAS DRAWN ZERO OBJECTIONS ACROSS ALL FOUR QUORUM ROUNDS AND WAITS ON NOTHING.** Clause 6
is now partitioned across THREE docs (see the clause-6 scope statement below). Full disposition
and the moved-scope inventory: `## Quorum verification log` → `### Round 4 disposition — SPLIT #2`.

**Prior head text follows (the round-4 head; accurate when written — its "a second split is
owed" is the split executed above).** Planned — **SPLIT APPLIED 2026-08-25 (iteration 125), THEN RE-QUORUMED AND BLOCKED
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

What this parent retains after split #2 is executable end to end with nothing depending on MCP
dispatch OR on session authority: **P6.A** (done, record), **P6.T** (toolchain floor), **P6.D**
(pinned dependency + one narrow allowlist line), and **P6.V** (verified commit-boundary law).
**P6.B-A2A** (the A2A half of the old P6.B) moved to the A2A child at split #2, together with
premises P2–P6, Decisions 3 and 4, the projection half of Decision 6, acceptance criteria
AC1–AC9/AC11–AC14 and their mutations — identifiers preserved there. (The `a2a_handler.go`
existence proof — **180 lines over stdlib + `protocol`, 0 SDK imports**, re-derived at the split
round via `gh api` at tag `v0.33.2`, premise row N23 — stays cited HERE because it grounds
Decision 2's measured asymmetry, and is inherited by the child as its handler-shape proof.)

**Clause-6 scope statement — read this before treating any milestone as charter discharge.**
Charter clause 6 names *"project the transition registry over MCP + publish the A2A agent card."*
After split #2, THIS DOC DELIVERS THE THREE ENABLING MILESTONES ONLY. Clause 6 is now
partitioned across THREE docs: this parent (`P6.T`/`P6.D`/`P6.V` — the enablers BOTH children
consume, blocked on nothing); child #1
[`w-mcp-dispatch-projection.md`](w-mcp-dispatch-projection.md) — "project the transition
registry over MCP", blocked UPSTREAM on `ailang#885`; and child #2
[`w-a2a-session-projection.md`](w-a2a-session-projection.md) — "publish the A2A agent card" plus
session-scoped A2A invocation, blocked LOCALLY on charter queue row 39 `w-session-authority`.
Each child carries its blocking predicate as a runnable command with controls. Clause 6 is NOT
satisfied by this parent alone — nor by any two of the three docs — and this doc does not
quietly narrow the clause: it partitions it, with both blocked halves named and tracked.

**Item**: `w-mcp-projection`  
**Clause**: clause-6 (protocol-native boundary; not DESIGN.md M6) — the enabling milestones only; the
MCP half is child #1's, the A2A half is child #2's  
**Estimate**: ~0.55 World days (P6.T ~0.1 + P6.D ~0.15 + P6.V ~0.3) — a re-cut, not a discount:
the removed ~0.7d is P6.B-A2A, carried unchanged by the A2A child (see Estimate honesty); no
wait of any kind remains IN THIS DOC (child #1 waits on `#885`, child #2 on charter row 39, and
each wait is its carrier's to run)  
**Author**: DESIGN-DOC-CREATOR (original 2026-07-28, rotation designer codex/gpt-5.6-sol;
revision round 2026-08-25; split applied 2026-08-25, iteration 125; split #2 applied
2026-08-26, iteration 126)  
**Verified against**: pinned **AILANG v0.30.0**, commit `e37b370`, at `/tmp/ailang-v0300/ailang`
(the `.ail` verifier axis) **and** upstream **`github.com/sunholo-data/ailang` tag `v0.33.2`
(`63e7909f`)** (the Go serving-seam axis — D-WORLD-5's two independent version axes); World `dev`
at `2e44e3e` (split-round anchors re-derived there — premise row N25)  
**Date**: 2026-07-28; **revised in place 2026-08-25**; **split applied 2026-08-25 (iteration
125)**; **split #2 applied 2026-08-26 (iteration 126)**

> **Scheduling truth (2026-08-26, split-#2 round).** Every milestone in THIS doc is executable
> NOW: `serveapi/protocol` shipped in upstream tag `v0.33.2` (measured — premise rows N1–N3,
> latest release re-confirmed in row N22), the toolchain floor is a self-contained move, and the
> single internal residual is milestone P6.V (the commit-boundary law's Z3 proof), a milestone of
> this doc. The FIRST milestone is a toolchain move — `v0.33.2`'s `go.mod` declares `go 1.26.6`
> while this repo pins `GOTOOLCHAIN: go1.25.6` at `ci.yml:21,102` and in `go.mod` (re-verified
> unlanded 2026-08-26 at `fcf18fa`: `go1.25.6` at both `ci.yml` sites and in `go.mod`) — and it
> is independently mergeable before any dependency lands. Nothing here waits on `#885`, on
> charter row 39, or on any human decision. (The split-round text this supersedes also cited the
> landed transition registry, the broker session constructor, and D-WORLD-26 = A: all still true,
> all now consumed by the A2A child, not by any milestone here.)

---

## Motivation

Clause 6 requires the available World transitions to leave the daemon through existing protocols:
MCP for callable tools and A2A for discovery, with the projected surface filtered by the
capabilities of each session. This is an authority boundary, not documentation: an unauthorized
transition must be absent from discovery, absent from the A2A card, and rejected if invoked by
name. That boundary is now delivered by the two split children — the A2A surface (the agent card
at `/.well-known/agent.json`, the A2A endpoint at `/a2a/`, and session-scoped invocation through
propose → verify → commit) by [`w-a2a-session-projection.md`](w-a2a-session-projection.md),
blocked on charter row 39, and the MCP surface by
[`w-mcp-dispatch-projection.md`](w-mcp-dispatch-projection.md), blocked on the upstream dispatch
seam (`ailang#885`). **This doc delivers the three enabling milestones** — toolchain floor,
pinned dependency, verified commit-boundary law — that BOTH halves consume when they unblock.

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
zero-cloud gate; since split #2 the handler-authoring work itself lives in the A2A child, while
this parent lands the pin). Upstream remains the wire-contract owner. World remains the state,
session, capability, and transaction owner — and, by measured necessity, the A2A-handler owner
(that authorship now scoped in the A2A child).

## Premises (hard constraints)

- **P1 — no new wire protocol.** A2A remains at `/.well-known/agent.json` and `/a2a/`; responses
  keep upstream framing. World does not decode and re-encode protocol traffic in a reverse
  proxy. Post-split this doc mounts NO protocol endpoint at all — the MCP surface (including its
  SSE-framed HTTP responses) travels with child #1 and the A2A surface with child #2. P1 is
  retained HERE as the clause-wide guardrail both children cite and must satisfy.
- *(**P2–P6** — exact per-session authority, transition truth, mandatory propose → verify →
  commit, single writer, frozen REST v1 — travel with the A2A child at split #2, identifiers
  preserved there: they are properties of the projection surface, and no milestone in
  P6.T/P6.D/P6.V touches them.)*
- **P7 — local-first and zero-cloud remain enforced.** The path-(c) import is now measured, not
  hypothetical: `serveapi/protocol` at `v0.33.2` imports **only the standard library** across its
  four source files, and admitting it moves the daemon-core dependency closure by **exactly one
  package** (249 → 250 over both gated patterns). It is admitted by **one allowlist line naming
  the PACKAGE path** `github.com/sunholo-data/ailang/serveapi/protocol` — never the module root,
  which would silently admit `internal/apiserver` and its measured 476 disallowed packages.
- **P8 — the projection adds no new AILANG law; milestone P6.V adds exactly one.** The protocol
  adapter's invariants remain cross-request/session behavior in Go integration/conformance tests
  (unchanged from July; those adapter tests now live in the A2A child's scope). But prerequisite 3's residual is precisely a missing pure-core law: none
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
| **Session credential carrier = `Authorization: Bearer` (Arm A)** — record; operative force carried by the A2A child since split #2 | session authority surface in a repo where clause 3 is inviolable; never an API key; fail closed | **D-WORLD-26** (Mark, attended, issue `#89`, 2026-08-25T19:06:41Z) | resolved | high |
| **A2A half SPLIT into `w-a2a-session-projection.md`, blocked on charter row 39 `w-session-authority`** | the repo has no inbound credential→session resolution at all (round-4 finding, confirmed first-party); a pre-existing gap is a queue row, not a revision | quorum round 4 + controller routing call (iteration 126) | resolved (split) | high |
| Toolchain floor moves FIRST (`go1.25.6 → go1.26.6`) | `v0.33.2` requires `go >= 1.26.6` (two-arm measured); independently mergeable before any dependency | live evidence | P6.T | medium |
| Allowlist admits the PACKAGE path, never the module root | the matcher is prefix-based; the root would admit `internal/apiserver`'s 476 disallowed packages | measured matcher semantics | P6.D | high |
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
- [ ] Do not call the static sketch export list a transition registry.
- [ ] Do not start the A2A child's milestone (P6.B-A2A, carried by
  `w-a2a-session-projection.md` since split #2) before P6.T (toolchain), P6.D (pinned dependency
  + narrow allowlist), and P6.V (the **verified** commit-boundary law) have each landed CI-green
  — recorded here because the enablers are this doc's deliverables; the child's start condition
  additionally requires charter row 39's session-authority design, and the condition binds the
  child. The July wording's
  three external prerequisites are discharged: the seam is `serveapi/protocol@v0.33.2`, the
  registry is `host/transitionreg/`, the session API is `host/broker.Session` (a session CONSTRUCTOR — round 4 established the
  inbound credential RESOLVER was never built; that is charter row 39, consumed by the A2A
  child); the Go half of the
  commit-boundary contract (atomic not-started-versus-committed via `JournalIntent` +
  `bindCommitIntentTx`, stable `InvocationID`, queryable `GetReceipt`/`GetEffectReceipt`) is
  landed — only the Z3 proof (P6.V) remains.
- [ ] Do not import `serveapi` proper, its embedded MCP/A2A handlers, or `CallbackRunner` — their
  closure spans the MCP SDK's 9 module roots and fails `TestDaemonDependencyAllowlist` by
  construction. The allowlist entry is the package path, never the module root.
- [ ] Do not claim that cancellation prevents a durable mutation once the coordinator has ACCEPTED a
  commit, and never return a definitive "not committed" result when the outcome is unknown.
- *(Five freeze rows travel with the A2A child at split #2, verbatim there: the
  API-key-as-session ban with D-WORLD-26's constraint (i), the ambient-export exposure ban, the
  direct A2A-to-commit ban, the socket-closure-observation rule, and the D7-deadline /
  unbounded-wait rules for projection code — they constrain projection code, of which this doc
  now has none.)*

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
callback-bounding — work that since split #2 lives in the A2A child,
[`w-a2a-session-projection.md`](w-a2a-session-projection.md)** — and the MCP handler question
moves to child #1 until `#885` delivers a dispatch seam. That is D-WORLD-5 (ARM A: import
upstream, pinned) executing as written on the half it can reach, not a new human ask.

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

## Decision 3 — Session-Scoped Projection in worldd (A2A) — MOVED TO THE A2A CHILD (split #2)

Carried in full, identifier preserved, by
[`w-a2a-session-projection.md`](w-a2a-session-projection.md) — with its responsibility 1
re-bound to the row-39 resolver interface per the round-4 `proposed_fix`: this doc's former text
equated `host/broker.NewSession` (a session CONSTRUCTOR whose grants are an argument) with a
credential resolver, which is exactly what round 4 caught. Nothing in P6.T/P6.D/P6.V consumes
Decision 3.

## Decision 4 — Surface Identity, A2A Card, and Compatibility — MOVED TO THE A2A CHILD (split #2)

Carried unchanged, identifier preserved, by
[`w-a2a-session-projection.md`](w-a2a-session-projection.md). (The cross-surface MCP≡A2A
equality obligation it recorded remains with the MCP dispatch child, as before.)

## Decision 5 — Process, Endpoint, and Package Layout (S2/S3)

| Path | Change |
|---|---|
| `.github/workflows/ci.yml` | P6.T: **FOUR sites** → `go1.26.6` — `GOTOOLCHAIN` at `:21`/`:102` and `setup-go`'s `go-version` at `:28`/`:109` (the two `setup-go` sites were added at the round-5 carve-out; the only workflow change) |
| `go.mod` | P6.T: `go 1.25.6 → 1.26.6`. **The `require github.com/sunholo-data/ailang v0.33.2` line was DEFERRED with P6.D at the round-5 carve-out** and lands with its first real consumer, in whichever child unblocks first. |
| `world/*.ail` + `scripts/verify_ail.sh` | P6.V: commit-boundary law + its `REQUIRED_VERIFIED` pin |

*(The `host/projection/`, `host/daemon/daemon.go` mount, and `cmd/ailang-worldd/` rows moved to
the A2A child's files table at split #2. The `go.sum`, `host/daemon/daemon_test.go` allowlist and
`host/daemon/protocol_use.go` rows moved out at the round-5 carve-out with P6.D — the last of
those is withdrawn outright and is created by no doc.)*

**Why is this not kernel growth?** After split #2 this doc adds nothing protocol-facing at all:
P6.T moves toolchain pins, P6.D adds one pinned stdlib-only dependency plus one narrow allowlist
line, and the one `world/` addition is P6.V's commit-boundary LAW — a proof over the
already-landed store/broker semantics, which is the opposite of kernel growth: it pins existing
behavior, adds none. (The "why is the projection an adapter, not an AILANG package and not
kernel growth" rationale — host effects at DESIGN §3.7's boundary, no transition semantics or
capability law owned — travels with the children, which own the adapter code.)

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

## Decision 6 — The commit-boundary contract (what P6.V proves); bounded waits MOVED TO THE A2A CHILD

**Bounded waits across the projection path (Standing rule 6) — MOVED at split #2.** The
per-request bounded-context contract — one bounded context per projection request with a
mandatory finite server maximum, propagation without replacement through every dependency,
prompt cancellation and cleanup, structured protocol timeout errors, no REST/store fallback —
governs projection request handling and travels verbatim with the A2A child (its AC13 tests it;
the child's carried Decision-6 section holds the full text). What REMAINS here is the half that
P6.V proves:

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
(`host/broker/recover.go:126`). What the A2A child's P6.B-A2A still needs from P6.V is the word
the prerequisite hinged on: a Z3-**verified** commit-boundary law pinned in `REQUIRED_VERIFIED` — the Go surface
exists; the proof does not.

**SSE lifecycle — MOVED WITH THE MCP SURFACE (split round).** Quorum r2's `gemini-3-1-pro`
objection concerned zombie TCP connections on the long-lived `/mcp/` SSE stream, and the r3
revision made that stream World-authored. Under the split there is NO SSE stream in this doc:
A2A is single bounded request/response, served entirely within the frozen D7 deadlines, and
`host/projection` performs no `ResponseController` deadline relaxation (the A2A child's AC14
asserts this at source level; both the criterion and the code moved there at split #2). What
SURVIVES from the objection is its OS-level kernel: a disconnected or deadline-expired transport
must be OBSERVED closed (`http.Server.ConnState` tracking or a client read error), never merely
context-cancelled (the child's AC13). The stream-lifetime maximum, the
`/mcp/`-route-local deadline relaxation, and mutations
`MUT-LEAK-SSE-CONN`/`MUT-SSE-REST-DEADLINE` travel with the child doc, where the streaming
handler will exist.

## Milestones (each independently CI-green and mergeable — hard repo convention; ~0.4 World days)

### P6.A — Upstream repro + conformance fixture — **DONE (iter-24, 2026-07-28)**

Retained as record: the upstream finding was filed (`sunholo-data/ailang#764`, following `#498`)
and this design landed. That finding is now **RESOLVED BY DELIVERY** — see Upstream Findings.
The frozen two-session conformance fixture carries forward unchanged: its A2A leg into the A2A
child ([`w-a2a-session-projection.md`](w-a2a-session-projection.md), since split #2), its MCP
leg with the MCP dispatch child.

### P6.T — Toolchain floor `go1.25.6 → go1.26.6` (~0.1d) — **FIRST, and independently mergeable**

No dependency, no new code. **FOUR pin sites in `ci.yml`, not two** (corrected at the round-5
carve-out, 2026-08-26, applying `gemini-3-1-pro`'s proposed fix — the prior premise row N7 had
narrowed its search to a single `grep` and made two of them invisible):

| site | today | after P6.T |
|---|---|---|
| `ci.yml:21` (`ailang-code verify gate` job env) | `GOTOOLCHAIN: go1.25.6` | `go1.26.6` |
| `ci.yml:28` (`actions/setup-go@v5`, same job) | `go-version: '1.25.6'` | `'1.26.6'` |
| `ci.yml:102` (`go host build + test gate` job env) | `GOTOOLCHAIN: go1.25.6` | `go1.26.6` |
| `ci.yml:109` (`actions/setup-go@v5`, same job) | `go-version: '1.25.6'` | `'1.26.6'` |
| `go.mod:3` | `go 1.25.6` | `go 1.26.6` |

Missing `:28`/`:109` is not cosmetic: `setup-go` would install 1.25.6 while `GOTOOLCHAIN` demanded
1.26.6, so every CI job would fetch a second toolchain at build time — the version skew the
reviewer predicted. The `verify_go.sh` deny-list (`:214-224`, enumerating exactly
`go1.26.0`–`go1.26.5`) is **untouched** — `go1.26.6` is not in it, and the committed canary is the
version-agnostic detector.

- **Independent justification, stated because the round-5 carve-out deferred the dependency this
  milestone was originally sequenced for** (`gpt5-6-sol`'s own clause: *"keep only independently
  useful changes here"*). P6.T stands on its own: the rig's `go` **binary** is `go1.26.4`, which
  IS deny-listed, so under `GOTOOLCHAIN=auto` today's `go 1.25.6` module resolves to a refused
  toolchain and **every local gate run in this repo must carry an explicit
  `GOTOOLCHAIN=go1.25.6`** — a standing tax on every sprint here, and a live foot-gun, both
  independent of any upstream pin. P6.T removes it. It adds no dependency, no allowlist entry and
  no production code, so it is untouched by the minimal-frozen-core objection that deferred P6.D.
- Why it is ALSO the dependency's prerequisite (no longer this doc's reason to land it):
  `v0.33.2` requires `go >= 1.26.6` (two-arm measured: `go get ...@v0.33.2` rc=1
  under `GOTOOLCHAIN=go1.25.6` with the exact floor message, rc=0 under `go1.26.6`;
  known-positive control `go get github.com/google/uuid@v1.6.0` rc=0 under `go1.25.6`). That
  sequencing now serves whichever child admits the dependency, not this doc.
- Rig nuance (measured 2026-08-25): the rig's `go` **binary** is `go1.26.4` (deny-listed);
  under `GOTOOLCHAIN=auto`, toolchain selection follows `go.mod` — today (`go 1.25.6`) it
  resolves to `go1.26.4` and the deny-list correctly refuses it, and in a `go 1.26.6` module it
  resolves to `go1.26.6` (measured). So until P6.T lands, local gate runs must pin
  `GOTOOLCHAIN=go1.25.6` explicitly; after it lands, `auto` is safe.
- Acceptance: canary re-run under the new toolchain — repro prints `OK` under `go1.26.6` AND the
  known-bad control prints `BUG` under `go1.26.5` (both re-run first-party 2026-08-25); full
  `verify_go.sh` rc=0 (deny-list, driver drift gate, armed race control, build+test); both CI
  jobs green with zero other diffs.

### P6.D — DEFERRED AT THE ROUND-5 CARVE-OUT (2026-08-26, iteration 126) — dependency admission is now ATOMIC WITH ITS FIRST REAL CONSUMER

Quorum round 5's blocking objection (`gpt5-6-sol`, present, `absent_reviewers` empty) is applied
here as its author wrote it. Verbatim:

> P6.D deliberately lands an otherwise unused direct dependency and a production file whose only
> purpose is a dummy symbol reference (`host/daemon/protocol_use.go: var _ =
> protocol.ValidateMCPName`). Both actual consumers are blocked child designs. This is speculative
> core growth and directly violates the minimal-frozen-core and route-to-extension axioms; making
> the dependency graph non-vacuous does not make the dependency operationally necessary.

and its proposed fix, also verbatim, which is what has been done:

> Remove P6.D, `host/daemon/protocol_use.go`, AC16, and its mutations from this parent. Move the
> pinned dependency and narrow package-path allowlist change into whichever child first becomes
> unblocked, where a real handler or adapter import provides the compile-visible use. Keep only
> independently useful changes here, and revise the scheduling text to state that dependency
> admission occurs atomically with its first real consumer rather than being pre-landed.

**Note the objection is a defect this doc INTRODUCED, one iteration ago, and not a pre-existing
one** — so it is answered here rather than routed to the charter queue (the split rule's clause
(c) discriminates exactly this way). Before split #2, P6.D's compile-visible use was P6.B-A2A's
first handler import: a real consumer. Split #2 moved that consumer to the A2A child and the
re-anchor replaced it with a dead symbol reference — trading a scoping problem for an axiom
violation. The reviewer caught it on the first round after it appeared.

**SCHEDULING RULE (the operative half).** Dependency admission — the `v0.33.2` pin in `go.mod`,
the single `allowedDepModules` **package-path** entry
`github.com/sunholo-data/ailang/serveapi/protocol`, its narrowness test, the `go.sum` movement,
and the 249 → 250 closure assertion — **occurs atomically with its first real consumer and is
never pre-landed.** It is carried by whichever child unblocks first:

- [`w-a2a-session-projection.md`](w-a2a-session-projection.md) — blocked on charter queue row 39
  `w-session-authority`; carries the admission as its milestone step 1 (its A2A handler import is
  the compile-visible use).
- [`w-mcp-dispatch-projection.md`](w-mcp-dispatch-projection.md) — blocked on
  [`ailang#885`](https://github.com/sunholo-data/ailang/issues/885); carries it instead if it
  unblocks first.

Whichever lands it, the other child inherits an already-admitted dependency and adds nothing.
The full P6.D specification (matcher semantics at `daemon_test.go:801`, the module-root-vs-
package-path measurement, the two named indirect bumps, the narrowness test, AC16 and mutations
`MUT-ALLOWLIST-ROOT`/`MUT-FACADE-IMPORT`) travels with it and is preserved verbatim in the A2A
child. Nothing was discarded; it was re-sequenced.

**`host/daemon/protocol_use.go` is not created by any doc.** It existed only in this milestone's
one-iteration-old re-anchor and is withdrawn.

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

### P6.B-A2A — MOVED TO THE A2A CHILD AT SPLIT #2 (2026-08-26, iteration 126)

The session projection + A2A conformance milestone (~0.7d) is carried, unchanged in substance
and under its own identifier, by
[`w-a2a-session-projection.md`](w-a2a-session-projection.md) — honestly BLOCKED on charter
queue row 39 `w-session-authority`, with the round-4 objection travelling verbatim in its
opening problem statement. It starts only after row 39's session-authority design lands AND
P6.T/P6.D/P6.V are green. Nothing in this doc's three milestones depends on it.

**Estimate honesty.** July's "~1.2d after prerequisites" became ~1.55d all-in at the revision
round; split #1 re-cut it to ~1.25d; split #2 re-cuts it to **~0.55d in this doc** (P6.T ~0.1 +
P6.D ~0.15 + P6.V ~0.3). Each re-cut is a re-cut, not a discount, and the removed effort is
accounted for by name: the MCP share of the old P6.B (~0.3d plus whatever the dispatch seam
actually demands) moved at split #1 to the MCP dispatch child and is deliberately UNESTIMATED
there (estimating against an unknown seam would be fiction), and **the ~0.7d P6.B-A2A moved at
split #2 to the A2A child at the SAME figure, unchanged** — it starts only when charter row 39's
design lands and this doc's three milestones are green. Nothing external remains on THIS doc's
critical path. If P6.V's encodable statement proves harder than one law, split P6.V out rather
than absorbing the overrun silently.

## Files to Create/Modify

| File | Milestone | Purpose |
|---|---|---|
| `design_docs/planned/w-mcp-projection.md` | P6.A (done) | This design; revised in place 2026-08-25; SPLIT applied 2026-08-25 (iteration 125) |
| `design_docs/planned/w-mcp-dispatch-projection.md` | split child #1 | Carries the MCP half + the round-3 objection verbatim; blocked on `#885`; NOT quorum-cleared |
| `design_docs/planned/w-a2a-session-projection.md` | split child #2 | Carries the A2A half + the round-4 objection verbatim; blocked on charter queue row 39 `w-session-authority`; NOT quorum-cleared |
| upstream issue in `sunholo-data/ailang` | P6.A (done) | `#764` — protocol-only module; RESOLVED by `v0.33.2` (`#885` — dispatch seam — is the CHILD's blocker, tracked there) |
| `.github/workflows/ci.yml` | P6.T | FOUR sites → `go1.26.6`: `GOTOOLCHAIN` at `:21`/`:102` **and `setup-go`'s `go-version` at `:28`/`:109`** (corrected at the round-5 carve-out; the only change) |
| `go.mod` | P6.T | `go 1.26.6` directive ONLY. The `require github.com/sunholo-data/ailang v0.33.2` line was DEFERRED with P6.D at the round-5 carve-out and lands with its first real consumer. |
| `world/*.ail` | P6.V | Commit-boundary law (encodable shape; fluency protocol first) |
| `scripts/verify_ail.sh` | P6.V | `REQUIRED_VERIFIED` gains the commit-boundary identity |

*(The `host/projection/projection.go`, `host/projection/projection_test.go`,
`host/daemon/daemon.go` mount, and `cmd/ailang-worldd/main.go` rows moved to the A2A child at
split #2. The `go.sum`, `host/daemon/daemon_test.go` allowlist and `host/daemon/protocol_use.go`
rows moved out at the round-5 carve-out with P6.D — dependency admission lands with its first
real consumer, never here.)*

No schema, benchmark baseline, or REST route changes. **No `/mcp/` route, no A2A route, and no
handler file of any kind exists anywhere in this doc's scope** — that is the two splits, not an
omission. The July claim "no
`.ail` or CI workflow file changes" remains **withdrawn**: P6.T touches `ci.yml` (toolchain pins
only) and P6.V touches `world/` plus the `verify_ail.sh` manifest — both deliberate, both named
above.

## Conflict Surface

*(Seven entries moved to the A2A child at split #2, verbatim there: HTTP server timeouts vs the
A2A route, the worldd single-writer flock, the frozen REST v1 route table, the shared JSON error
envelope, `host/registry` vs `host/transitionreg`, broker/propose-verify-commit dispatch, and
the performance baseline — they are conflict surfaces of projection code, of which this doc now
has none. The child also adds a row-39-overlap entry the round-4 `catch` demanded.)*

- **`scripts/verify_ail.sh`.** Floor re-measured in the split round (2026-08-25, full gate run,
  PASSED — premise row N24): **10 required identities verified, 40 named tests pass, 9/9
  world-package steps**. P6.V deliberately RAISES the identity floor by pinning the
  commit-boundary law in `REQUIRED_VERIFIED` (`verify_ail.sh:274-279`) — a manifest change this
  doc names, not an accident.
- **`scripts/verify_go.sh`.** Runs the driver-drift gate, the go1.26.x deny-list (`:214-224`,
  exactly `go1.26.0`–`go1.26.5`), the armed race-detector control, and `go build`/unskipped
  `go test`. P6.T interplay measured 2026-08-25: the rig `go` binary is `go1.26.4`
  (deny-listed), so under `GOTOOLCHAIN=auto` the gate correctly FATALs today and passes once
  `go.mod` says `1.26.6` (selection then resolves to `go1.26.6`, measured). (The
  "protocol socket tests must run on CI; local sandbox denial is not a skip condition" rule
  travels with the A2A child's tests.)
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
- **serve-api cwd sensitivity.** Historical (P6.A record) only: no milestone in this doc runs an
  upstream server, and the A2A child's World-owned handlers are embedded in worldd — so
  module-loading cwd is moot.

## Systemic-Issue Audit

| Question | Finding | Response |
|---|---|---|
| Is the mismatch local to one endpoint? | No; discovery, A2A card, authentication, and invocation share it | Require one upstream callback seam and one World predicate |
| Can configuration solve it? | No; `--caps` and API key are process-wide | Do not encode sessions as flags |
| Can a sidecar solve it safely? | No; static exports, built-in feedback, and lifecycle remain | Reject path (b) |
| Is a new protocol required? | No | Reuse upstream codecs/wire types |
| Does this reveal a missing landed subsystem? | It did in July (only the epoch registry existed); as of 2026-08-25 the transition registry (`host/transitionreg`) and broker session API are landed, and the remaining gap is one missing Z3 law | Consume the landed subsystems; scope the law as P6.V |
| Is the toolchain move a one-off? | No — go1.26.6 also unblocks item 4e's parked remediation | Land P6.T narrowly here; leave 4e to the queue (Deferred Scope) |
| Could a gate pass vacuously? | Yes; empty registry, one identical session, list-only tests, or skipped sockets | Require non-empty unequal session sets, calls, denials, and zero skips (the session-set/denial/socket halves bind the A2A child since split #2; zero skips and the named RED mutations bind every gate that remains here) |
| Did one doc bundle surfaces with different readiness? | Yes — round 3 surfaced it (MCP blocked on a missing dispatch seam vs A2A), and round 4 surfaced it AGAIN: session authority (a property of P6.B-A2A alone; no inbound credential→session resolution exists) vs the enablers, which have never drawn an objection | SPLIT #1 (iteration 125): the MCP half → `w-mcp-dispatch-projection.md`, blocked on `#885`; SPLIT #2 (iteration 126): the A2A half → `w-a2a-session-projection.md`, blocked on charter queue row 39; the enablers stay |

## Deferred Scope

- **The MCP half of clause 6 — SPLIT OUT, not silently deferred.** The `/mcp/` endpoint, MCP
  JSON-RPC dispatch (`initialize`, `tools/list`, `tools/call`), dispatch-bound envelope framing,
  SSE stream lifetime and the route-local `ResponseController` relaxation, the cross-surface
  MCP≡A2A equality criterion (old AC3), old AC8's SSE-framing assertion, old AC14, and mutations
  `MUT-PLAIN-JSON`/`MUT-LEAK-SSE-CONN`/`MUT-SSE-REST-DEADLINE` are carried by
  [`w-mcp-dispatch-projection.md`](w-mcp-dispatch-projection.md), blocked on `ailang#885`.
- **The A2A half of clause 6 — SPLIT OUT at split #2, not silently deferred.** P6.B-A2A,
  premises P2–P6, Decisions 3 and 4, the bounded-wait contract, acceptance criteria
  AC1–AC9/AC11–AC14 with their named mutations, and the projection conflict-surface entries are
  carried by [`w-a2a-session-projection.md`](w-a2a-session-projection.md), blocked on charter
  queue row 39 `w-session-authority` (no inbound credential→session resolution exists — round-4
  finding, confirmed first-party).
- Designing or landing the clause-3 capability law, broker, in-process session model, or
  transition registry (all landed; this item only consumes them) — and designing the HTTP-facing
  session-authority boundary, which is charter row 39's to design and the A2A child's to
  consume.
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

Split #1 moved three criteria to the MCP dispatch child (old AC3's cross-surface MCP≡A2A
equality, old AC8's MCP SSE-framing assertion, old AC14's SSE/REST deadline separation; the
freed numbers were re-used for A2A-scoped criteria — that mapping is in the split-round quorum
entry). Split #2 now moves the A2A-scoped criteria — **AC1–AC9 and AC11–AC14, identifiers
preserved** — to the A2A child, `w-a2a-session-projection.md`, together with their named
mutations. What remains below is the enabler set: **AC10, AC15, AC17** (AC16 was deferred with P6.D at
the round-5 carve-out — see below), retained under
their original identifiers so four rounds of quorum history stay traceable. The numbering gaps
are deliberate on both sides, and the split-#2 inventory is recorded in the
`Round 4 disposition — SPLIT #2` quorum entry.

- [ ] **AC10 — honest gates:** `verify_ail.sh` holds the floor re-measured 2026-08-25 in the
  split round (10 required identities, 40 named tests, 9/9 world-package steps — premise row
  N24) plus EXACTLY the P6.V additions — the commit-boundary identity(ies) named in the P6.V
  implementation report, nothing else moves; `verify_go.sh`, benchmark smoke, and both CI jobs
  pass with zero skips. (The zero-skipped-PROTOCOL-tests clause travels to the A2A child with
  the protocol tests it governs, along with `MUT-SKIP-SOCKET`.)
- [ ] **AC15 — toolchain floor (P6.T):** `GOTOOLCHAIN: go1.26.6` at BOTH `ci.yml` sites and
  `go 1.26.6` in `go.mod`; the `verify_go.sh` deny-list still enumerates exactly
  `go1.26.0`–`go1.26.5`; the committed canary prints `OK` under `go1.26.6` while its known-bad
  control prints `BUG` under a deny-listed toolchain in the same run; full `verify_go.sh` rc=0.
- [ ] ~~**AC16 — pinned dependency, narrow admission (P6.D)**~~ — **DEFERRED at the round-5
  carve-out (2026-08-26)** together with P6.D itself, per `gpt5-6-sol`'s verbatim fix. The row is
  not weakened and not dropped: it travels intact to whichever child admits the dependency (see
  the P6.D section above) and is preserved in `w-a2a-session-projection.md`. The identifier stays
  reserved so five rounds of quorum history remain traceable.
- [ ] **AC17 — verified commit-boundary law (P6.V):** a named commit-boundary identity is
  Z3-verified on the pinned v0.30.0 binary and pinned in `REQUIRED_VERIFIED`; the gate's JSON
  checks (`verify.errors == 0`, `counterexample == 0`, required identity present and `verified`)
  all hold — guarding the silent `unknown sort` exit-0 trap; if the encodable-shape fallback was
  taken, the limitation is recorded in the doc and the named test-only law is in the 40+ named
  tests.

## Non-Vacuity — Named RED Mutation for Every Gate

| Gate | Named RED mutation (concrete edit) | Required red observation |
|---|---|---|
| AC15 `MUT-TOOLCHAIN-REGRESS` | set ONE `ci.yml` site (or `go.mod`) back to `go1.25.6` | dependency resolution REDs with the measured floor message `requires go >= 1.26.6 (running go 1.25.6; GOTOOLCHAIN=go1.25.6)` |
| AC15 `MUT-CANARY-BLIND` | run the committed canary under deny-listed `go1.26.5` | repro prints `BUG: Field="" want "stateRoot"` (re-run 2026-08-25) — proving the detector still SEES the defect class, so `OK` under go1.26.6 is a measurement |
| AC17 `MUT-LAW-BREAK` | mutate the commit-boundary law body to permit a receipt without its journal intent (or a second receipt under one invocation ID) | `ai-check` reports `counterexample > 0` and `verify_ail.sh` REDs |
| AC17 `MUT-SORT-SILENT` | retype one law parameter to an ADT-bearing (`Proposal`-class) sort | `verify.errors > 0` in the JSON and `verify_ail.sh` REDs — proving the silent-exit-0 trap is guarded, not relied on |

Moved with split #1: `MUT-PLAIN-JSON` (SSE framing), `MUT-LEAK-SSE-CONN` (stream-deadline
zombie connection), and `MUT-SSE-REST-DEADLINE` (route-local relaxation scope) — recorded in the
MCP dispatch child as obligations that must reappear in its acceptance table when it becomes
writable. Moved with split #2, carried in the A2A child's acceptance table under their preserved
names: the eighteen projection mutations `MUT-PROTO-OWNER`, `MUT-SESSION-UNION`,
`MUT-CARD-GLOBAL`, `MUT-UNFILTERED-PROJECTION`, `MUT-CALL-BYPASS`, `MUT-DEFAULT-CAPS`,
`MUT-ALT-HEADER`, `MUT-KEY-AS-SESSION`, `MUT-SPLIT-SNAPSHOT`, `MUT-A2A-SHAPE`,
`MUT-SECOND-OPEN`, `MUT-SKIP-SOCKET`, `MUT-CLOUD-DEP`, `MUT-STARTUP-CACHE`,
`MUT-DROP-DEADLINE`, `MUT-COMMIT-BOUNDARY-LIE`, `MUT-LEAK-CONN`, and `MUT-DEADLINE-RELAX`.

## Axiom Compliance

| Axiom | Score | Justification |
|---|---:|---|
| A1 Determinism | +1 | pinned toolchain, pinned dependency version, and a Z3-proven law are deterministic, re-runnable gates |
| A2 Replayability | 0 | no transaction-path change in this doc (the projection's replay axioms travel with the children) |
| A3 Effect Legibility | +1 | no new effects; P6.V's law is pure core, and the one new dependency is measured stdlib-only |
| A4 Explicit Authority | +1 | the allowlist admits ONE package path, never a module root — authority to import is explicit, narrow, and gate-enforced |
| A5 Bounded Verification | +2 | P6.V adds a Z3-verified commit-boundary law pinned in `REQUIRED_VERIFIED`; every gate carries a named RED mutation |
| A6 Safe Concurrency | 0 | no concurrency surface changes in this doc |
| A7 Machines First | +1 | machine-checkable gates (allowlist test, toolchain canary, JSON verify checks) over prose claims |
| A8 Minimal Syntax | 0 | no language syntax change |
| A9 Cost Visibility | 0 | no new budget claim |
| A10 Composability | +1 | upstream-owned wire contract admitted at a pinned version through a narrow package-path seam |
| A11 Structured Failure | +1 | gate failures are named and specific (measured floor message, intruder-by-name, `verify.errors` JSON) |
| A12 System Boundary | +2 | zero kernel growth: P6.V pins existing behavior; nothing protocol-facing lands in this doc |

**Net: +10; hard axioms A1/A3/A4/A7 are non-negative.** (The pre-split projection scoring — net
+14 — travels with the A2A child, which inherits that table and re-scores at unblock if row 39's
contract shifts it.)

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
| N7 [R] | **CORRECTED at the round-5 carve-out (2026-08-26, iteration 126) — the previous form of this row narrowed its search to one file and thereby hid a live pin site.** Repo-wide, the toolchain is pinned in exactly ONE tracked file: `.github/workflows/ci.yml`, at **FOUR** sites — `GOTOOLCHAIN: go1.25.6` at `:21` and `:102`, **and `go-version: '1.25.6'` for `actions/setup-go@v5` at `:28` and `:109`** — plus `go.mod`'s `go 1.25.6` directive. No other tracked file carries a live pin: `scripts/verify_go.sh:222` is an advice STRING inside an error message, `host/store/toolchain_canary_test.go:7` and `scripts/verify_world_package.sh:226` are comments, `bench/BASELINE.md` and `docs/SELF_MOD_PUBLISH.md:120` are recorded historical measurements, and the four `sprint_*.json` files are archived sprint records. `.golangci.yml`, `.golangci.yaml`, `Makefile`, `Dockerfile` and `.tool-versions` — the candidates the reviewer named explicitly — are all **ABSENT** from the repo (`test -e` on each). | `git grep -n 'GOTOOLCHAIN' -- .` -> **5** tracked files; `git grep -n '1\.25\.6' -- .` -> **10** tracked files; `git grep -n 'go1\.26' -- .` ; `for f in .golangci.yml .golangci.yaml Makefile Dockerfile .tool-versions; do test -e $f; done`. Known-positive control in the same shape, `git grep -c 'modernc.org/sqlite'` -> **12** files; fresh-literal negative control -> **0**. | Confirmed; **and the correction is material** — `ci.yml:28`/`:109` were absent from P6.T's modification list and are now added to it |
| N8 [R] | `verify_go.sh` deny-list + canary | `sed -n '214,224p' scripts/verify_go.sh` → case arms exactly `go1.26.0\|...\|go1.26.5`; canary re-run: `GOTOOLCHAIN=go1.26.6 go run .` → `OK`; known-bad control `GOTOOLCHAIN=go1.26.5` → `BUG: Field="" want "stateRoot"` | Confirmed; go1.26.6 is NOT deny-listed and the detector still sees the defect |
| N9 [R] | **DISAGREEMENT with the directive's "the rig's default `go` is already go1.26.6"** | `which -a go` → `/opt/homebrew/bin/go` only; `go version` → **go1.26.4** darwin/arm64; `go env GOTOOLCHAIN` → `auto`; in-repo (`go 1.25.6` module) `go env GOVERSION` → **go1.26.4** (deny-listed!); in a `go 1.26.6` module → **go1.26.6** | The BINARY is go1.26.4; `auto` selection reaches go1.26.6 only once `go.mod` says so. Consequence recorded in P6.T: pin `GOTOOLCHAIN` explicitly for local gate runs until P6.T lands |
| N10 [R] | Pristine dependency closure over BOTH gated patterns | `go list -deps ./host/daemon/... ./cmd/ailang-worldd/... \| sort -u \| wc -l` → **249**; same-call control: `grep -c github.com/google/uuid` → 1 | Confirmed 249 (the charter's iter-120 "pristine 244" is the older baseline; intervening landed work moved it — both agree post-pin is 250 and the non-stdlib delta is exactly 1) |
| N11 [C] | Probe import + unchanged allowlist REDs the gate naming exactly one intruder | `go test ./host/daemon/ -run TestDaemonDependencyAllowlist` with a `serveapi/protocol` probe import | `1 package(s) outside the zero-cloud allowlist ... github.com/sunholo-data/ailang/serveapi/protocol`; with the ONE package-path line added: green at **250**; sentinel control proved the diff can see an addition |
| N12 [R] | Allowlist shape + matcher semantics | `allowedDepModules` at `host/daemon/daemon_test.go:747` (11 module-root entries); matcher `d == m \|\| strings.HasPrefix(d, m+"/")` (~`:800-806`) | Confirmed: entries are path prefixes, so a PACKAGE-path entry admits only that package (+ subpackages beneath it) and refuses `internal/apiserver` — the matcher DOES distinguish; no extra matcher milestone needed, only the narrowness test |
| N13 [R] | `verify_ail.sh` floor | full gate run, `AILANG_BIN=/tmp/ailang-v0300/ailang ./scripts/verify_ail.sh` → PASSED | **10 required identities verified, 40 named tests pass, 9/9 world-package steps**; `REQUIRED_VERIFIED` at `verify_ail.sh:274-279`, 4 module entries, none a commit-boundary law (July's "4/9/14" floor is long stale) |
| N14 [R] | Commit-boundary Go surface is landed and public | `sed -n '24,30p' host/store/journal.go`; `grep -n bindCommitIntentTx host/store/store.go`; `grep -n 'GetReceipt\|GetEffectReceipt' host/store/journal.go`; `grep -n recoverCommitPending host/broker/recover.go` | `JournalIntent` at `journal.go:28`; **`bindCommitIntentTx` at `store.go:1025` — NOT `:1015` as the directive stated (line drift; one-line disagreement, re-derived value adopted)**; `GetReceipt :813`; `GetEffectReceipt :852`; `recoverCommitPending :126` |
| N18 [R] | No lockfile/module-pin gate collides with the `go.sum` movement | `grep -n 'go\.sum\|go mod verify\|go mod tidy' scripts/verify_go.sh` → 0 hits; same-call control `grep -c GOVERSION` on the same file → 1 | Confirmed (the planned `w-ail-gate-module-pin` pins the **`.ail`** module set, a different axis) |
| N19 [I] | Upstream closure guarantee is CI-enforced upstream (`check-protocol-closure`, 5-arm refusal self-test) | charter iter-120 row | Inherited, unverified here; first-party enforcement is N11's gate either way |
| N20 [I] | `Proposal`-sort Z3 limitation (`unknown sort`, silent exit 0) | charter (recorded v0.30.0 limitation; also `w-m1-ailang-hardening`) | Inherited; P6.V designs around it and MUT-SORT-SILENT proves the JSON guard |

Rows **N15, N16 and N17** (transition registry + broker session CONSTRUCTOR; no HTTP
session-credential carrier in `host/`; D7 constants/wiring drift) moved to the A2A child at
split #2, identifiers preserved and marked inherited there — they existed only to serve
`P6.B-A2A`. Surviving references to them in this doc are historical records (the round-4 quorum
quote names N15; the Open Decisions record cites N16), and each resolves in the child, which
carries the rows alongside the objection that names them.

### Split-round rows (2026-08-25, iteration 125 — all [R], re-derived first-party in the split session)

| # | Premise | Command / evidence | Result |
|---|---|---|---|
| N21 [R] | `#885` is OPEN with 0 comments (the child's blocking predicate holds) | `gh issue view 885 --repo sunholo-data/ailang --json state,title,createdAt,comments` → `state=OPEN`, `comments=0`, created `2026-08-25T16:37:44Z`, title "serveapi/protocol has no MCP dispatch — zero-cloud consumers can reach A2A but not MCP (follow-on to #764)"; same-call control `gh issue view 764 …` → `CLOSED`, 6 comments | Confirmed — the instrument can see closure and comments; the MCP half stays blocked |
| N22 [R] | Upstream latest release is still `v0.33.2`; no `v0.34.*` tag exists | `gh api repos/sunholo-data/ailang/releases/latest --jq .tag_name` → `v0.33.2`; `gh api …/matching-refs/tags/v0.34` → empty (rc=0); same-call control `…/tags/v0.33` → `v0.33.0`, `v0.33.1`, `v0.33.2` | Confirmed — the pin target is unchanged and nothing newer shipped |
| N23 [R] | The A2A existence proof + the dispatch gap, at tag `v0.33.2`, without any local upstream checkout | `gh api 'repos/sunholo-data/ailang/contents/serveapi/a2a_handler.go?ref=v0.33.2'` (and `mcp_handler.go`), base64-decoded: `a2a_handler.go` **180 lines**, import block = `context`, `encoding/json`, `fmt`, `log`, `net/http` + `serveapi/protocol`, `grep -c modelcontextprotocol` → **0**; `mcp_handler.go` **187 lines**, SDK import count **1**; `grep -c 'tools/list\|jsonrpc'` in `mcp_handler.go` → **0** while the same grep in `a2a_handler.go` → **2** | Confirmed — the A2A handler is buildable over stdlib+`protocol`; MCP dispatch is entirely SDK-side (the stronger form of the r3 objection: `mcp_handler.go` never string-matches method names at all) |
| N24 [R] | `verify_ail.sh` floor unchanged | `AILANG_BIN=/tmp/ailang-v0300/ailang ./scripts/verify_ail.sh` → `✓ verify gate PASSED: 10 required identities verified, 40 named tests pass`; `✓ world package gate PASSED: 9/9 steps performed non-zero work`; binary `--version` → v0.30.0, commit `e37b370` | Confirmed — AC10's floor is a live split-round measurement, not carried prose |
| N25 [R] | Every local line-number anchor this doc cites, at `2e44e3e` | `grep -n GOTOOLCHAIN .github/workflows/ci.yml` → `:21`,`:102`; `head go.mod` → `go 1.25.6`; `grep -n 'func NewSession' host/broker/broker.go` → `:87`; `grep -n 'const ('` → `daemon.go:78`, `grep -n 'func newServer'` → `:619`; `grep -n REQUIRED_VERIFIED scripts/verify_ail.sh` → `:274` (4 module entries; 10 identities by count); `grep -n allowedDepModules host/daemon/daemon_test.go` → `:747` (var), matcher `d == m \|\| strings.HasPrefix(d, m+"/")` at `:801`; `JournalIntent` → `journal.go:28`, `bindCommitIntentTx` → `store.go:1025`, `GetReceipt` → `:813`, `GetEffectReceipt` → `:852`, `recoverCommitPending` → `recover.go:126`; deny-list case arms exactly `go1.26.0\|…\|go1.26.5` | Confirmed, zero drift since the revision round (one refinement: the matcher's exact line is `:801`; the revision round cited `~:800-806`) |
| N26 [R] | D-WORLD-26 is RESOLVED = ARM A with both constraints | charter Human Decision Ledger row read first-party this session (`design_docs/world-mission.md:708`): answered `A` by Mark, attended, issue `#89`, `2026-08-25T19:06:41Z`, verbatim one-character comment, allowlist-enforced directive read; constraints (i) API-key-never-a-session and (ii) fail-closed travel with the answer | Confirmed — Open Decision 2 is CLOSED; zero forks remain in this doc (since split #2 the ruling's operative force — resolver contract, AC6, `MUT-ALT-HEADER`/`MUT-KEY-AS-SESSION` — binds the A2A child) |

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
   resolver's contract (all three carried by the A2A child since split #2 — and round 4 then
   established that this answer settled the credential ENVELOPE while the resolver CONTENTS were
   never built: charter queue row 39): **(i)** the static `serve-api` API key remains forbidden as a session
   model (iteration 24 measured it process-wide, so it cannot represent a session) — a `Bearer`
   value here is a SESSION credential and never an API key; **(ii)** clause 3 still binds — the
   resolver fails closed on an absent, malformed, or unknown credential, never degrading to an
   unauthenticated surface.
3. **Transition descriptor schema — ANSWERED BY LANDING.** `host/transitionreg.Descriptor`
   (stable ID, schema, description, effect requirements) is used verbatim, projected through
   `protocol.ToolDescriptor`/`CallerSurface`; no projection-owned schema. This is now Decision 3
   text (carried by the A2A child since split #2), not a decision.
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

### Round 4 disposition — SPLIT #2 (2026-08-26, iteration 126; no quorum was run this round)

**Round count: FOUR quorum rounds plus one carve-out revision.** Per the skill's rule (a) — from
round 3 on, record which surface each round's objections name — the per-round surfaces are: **R1**
— three objections on three surfaces (bounded waits, `gpt5-6-sol`; SSE-vs-REST timeout conflict,
`gemini-3-1-pro`; upstream appetite, controller, non-blocking). **R2** — commit boundary
(`gpt5-6-sol`) + SSE socket lifetime (`gemini-3-1-pro`, pass-with-objection). **R3** — MCP
dispatch seam (`gpt5-6-sol`) + status/fork self-consistency (`gemini-3-1-pro`). **R4** —
**session authority ONLY** (`gpt5-6-sol`), while **`gemini-3-1-pro` PASSED for the first time in
four rounds** (its single non-blocking P6.V wording nit was applied verbatim by the controller in
iteration 125 and is NOT re-applied here). Monotone convergence: three surfaces, then two, then
one-plus-a-nit — and each departing surface left as its own doc or charter row, never as a
deletion.

**Disposition: SPLIT, per the skill's decomposition rule** — objections localising onto ONE
surface while another reviewer flips to pass mean the doc's SCOPE is wrong, not its content.
**Round 4's blocking objection is therefore DISCHARGED BY SPLIT, NOT ANSWERED**: all three of its
fields (`strongest_objection`, `catch`, `proposed_fix`) are carried VERBATIM — verified
character-for-character against the round-4 artifact
(`.ailang/state/mission-quorum/w-mcp-projection-2026-08-25T21-25-44Z.json`), with a mutated-quote
negative control — into the opening problem statement of the new child,
[`w-a2a-session-projection.md`](w-a2a-session-projection.md). The pre-existing gap the objection
surfaced (no inbound credential→session resolution anywhere in `host/ cmd/` — the controller's
first-party iteration-125 measurement, re-run at `fcf18fa` in the split-#2 session with identical
values) is **charter queue row 39 `w-session-authority`**: a defect this doc *fails to fix*, not
one it *introduces* — a queue row, not a third revision. The `proposed_fix`'s own second half
prescribed exactly this ("change status to BLOCKED and add an independently mergeable
prerequisite milestone … P6.B-A2A must depend on it").

**Moved to the A2A child (identifiers preserved):** milestone P6.B-A2A (~0.7d, unchanged);
premises P2–P6; Decisions 3 (responsibility 1 re-bound to the row-39 resolver interface, per the
`proposed_fix`) and 4; the bounded-wait half of Decision 6; acceptance criteria AC1–AC9 and
AC11–AC14 plus AC10's zero-skipped-protocol-tests clause; the eighteen projection mutations
(named in the Non-Vacuity note above); seven conflict-surface entries; premise rows N15/N16/N17;
the `host/projection/*`, daemon-mount, and `cmd/ailang-worldd` files rows; and D-WORLD-26's
operative force. **What remains here:** P6.A (record), P6.T, P6.D (compile-visible use
re-anchored onto P6.D's own `host/daemon/protocol_use.go`, since the former anchor — "the
P6.B-A2A skeleton's first import" — departed), P6.V; AC10 and AC15–AC17; six mutations (with
`MUT-FACADE-IMPORT` re-anchored to the same P6.D use site); premises P1/P7/P8. **Everything that
remains has drawn ZERO objections across all four rounds.** The surfaces left for the re-quorum
to probe are exactly the enabler ones: the toolchain floor, the pinned-dependency admission, and
the commit-boundary law.

### Round 5 — BLOCKED, then NARROW-REFINEMENT CARVE-OUT (2026-08-26, iteration 126; both reviewers PRESENT, `absent_reviewers` EMPTY; `metered=$0.219099`)

Artifact: `.ailang/state/mission-quorum/w-mcp-projection-2026-08-26T02-12-17Z.json`
(`--max-cost-usd 0.35`, controller verdict `pass`). Reviewers `gpt5-6-sol` (`reject`, $0.152755)
and `gemini-3-1-pro` (`reject`, $0.066344); `total_tokens_in` 60,429 / `out` 549.
**`absent_reviewers` is EMPTY and both reviewers are EXTERNAL** — the `presentCount` was not
satisfied by the controller's own verdict, so this is a real two-reviewer reading.

**Round count: FIVE**, plus two splits and now one carve-out revision. Said out loud because the
skill requires it: a doc past round 4 is data about this loop's SCOPING, not about the document.
Surfaced to Mark in the iteration-126 report for exactly that reason.

**Per-round surface tracking (the decomposition rule's clause (a), continued).**

| round | `gpt5-6-sol` surface | `gemini-3-1-pro` surface | disposition |
|---|---|---|---|
| 1 | multiple | multiple | revise |
| 2 | multiple | multiple | revise |
| 3 | MCP dispatch seam (DIRECTION) | open fork (`D-WORLD-26`) | **SPLIT #1** + human decision |
| 4 | session authority (ONE surface) | **PASS** (narrow P6.V nit) | **SPLIT #2** + charter row 39 |
| 5 | P6.D dependency-admission timing | P6.T premise completeness (N7) | **CARVE-OUT** (both fixes applied verbatim) |

**Round 5 is NOT a decomposition signal and was not treated as one.** The objections did not
localise onto one surface — they landed on two different milestones — and no reviewer flipped to
pass. Both are narrow, both carry a concrete reviewer-authored `proposed_fix`, and neither
disputes the design DIRECTION (`D-WORLD-5`, Mark, attended, already settled that World imports
`serveapi/protocol` pinned behind one narrow allowlist line; `gpt5-6-sol` disputes only the
*timing* of admission). That is the narrow-refinement carve-out's exact predicate, so the
controller made a bounded 2nd revision applying the reviewers' own text. It is not force-passing:
each fix REMOVES or CORRECTS the objected-to thing rather than proceeding over it.

**`gemini-3-1-pro`'s objection was MEASURED first-party, not forwarded (rule 3f), and the
measurement CONFIRMED it and found a live defect.** The repo-wide sweep it asked for
(`git grep -n 'GOTOOLCHAIN' -- .` -> 5 tracked files; `git grep -n '1\.25\.6' -- .` -> 10 tracked
files; known-positive control `modernc.org/sqlite` -> 12 files; fresh-literal negative control ->
0) found **two pin sites P6.T's modification list did not name**: `actions/setup-go@v5`'s
`go-version: '1.25.6'` at `ci.yml:28` and `:109`, one in each job. Bumping `GOTOOLCHAIN` alone
would have left `setup-go` installing 1.25.6 against a 1.26.6 demand — the version skew the
reviewer predicted, in the milestone this iteration was about to sprint. The premise row that
existed to rule this out (N7) had narrowed its own search to `grep -n GOTOOLCHAIN
.github/workflows/ci.yml`, which cannot see a `go-version` key. N7 is rewritten with the
repo-wide command, its counts, its controls, and the negative result for every candidate the
reviewer named (`.golangci.yml`, `.golangci.yaml`, `Makefile`, `Dockerfile`, `.tool-versions` —
all ABSENT, `test -e` on each). **This is the quorum paying for itself**: a narrowed instrument
reporting a clean premise is this loop's signature defect, and an external reviewer caught it.

**`gpt5-6-sol`'s objection is a defect split #2 INTRODUCED one iteration ago**, so clause (c) of
the decomposition rule routes it here rather than to the charter queue: before split #2 P6.D's
compile-visible use was P6.B-A2A's first handler import — a real consumer — and moving that
consumer to the A2A child left a dead symbol reference (`host/daemon/protocol_use.go: var _ =
protocol.ValidateMCPName`) whose only job was to make the dependency graph non-vacuous. The fix
is applied as written: P6.D, `host/daemon/protocol_use.go`, AC16 and mutations
`MUT-ALLOWLIST-ROOT`/`MUT-FACADE-IMPORT` are removed from this parent; dependency admission moves
to whichever child unblocks first, atomically with its first real consumer; and the scheduling
text now says so. Nothing was discarded — the full specification is preserved in
`w-a2a-session-projection.md`.

**Consequence for P6.T, stated rather than hidden.** `gpt5-6-sol`'s fix contains the clause *"keep
only independently useful changes here"*, and P6.T's original justification was *"`v0.33.2`
requires `go >= 1.26.6`"* — i.e. the dependency that just deferred. Rather than let the same
objection land on P6.T next round, its independent justification is now stated and measured in
the milestone body: the rig's `go` binary is `go1.26.4`, which `verify_go.sh:214-224` deny-lists,
so every local gate run in this repo currently requires an explicit `GOTOOLCHAIN=go1.25.6`. P6.T
removes that standing tax and adds no dependency, no allowlist entry and no production code, so
the minimal-frozen-core objection does not reach it. A prior justification — *"it also unblocks
queue item 4e's parked remediation"*, recorded in the charter at iteration 90 — is **retired as
stale**: item `4e w-race-gate-blindspot` is COMPLETE (iteration 46, 2026-08-04, PR #36 → `f19acac`).

**What the parent now contains:** `P6.A` (record), **`P6.T`** (~0.1d, toolchain floor, four pin
sites) and **`P6.V`** (~0.3d, verified commit-boundary law in `world/*.ail`) — ~0.4d. Both are
free of dependency admission, both are independently CI-green and mergeable, and `P6.V` has drawn
no objection in five rounds.

## Related Documents

- [w-mcp-dispatch-projection.md](w-mcp-dispatch-projection.md) — SPLIT child #1: carries the
  MCP dispatch half and the round-3 objection verbatim; blocked on `ailang#885`; NOT
  quorum-cleared
- [w-a2a-session-projection.md](w-a2a-session-projection.md) — SPLIT child #2: carries the A2A
  projection half and the round-4 objection verbatim; blocked on charter queue row 39
  `w-session-authority`; NOT quorum-cleared
- [world-mission.md](../world-mission.md) — clause 6, queue row 5, **queue row 39
  `w-session-authority`** (child #2's blocker), D-WORLD-5/D-WORLD-25/**D-WORLD-26** (the carrier
  ruling whose operative force is child #2's), premise measurements this doc consumes
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
