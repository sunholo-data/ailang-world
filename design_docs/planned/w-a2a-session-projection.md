# w-a2a-session-projection — Session-Scoped A2A Projection and Agent Card (SPLIT child of `w-mcp-projection`)

**Status**: Planned — **BLOCKED on charter queue row 39 `w-session-authority` — a LOCAL
prerequisite design, not an upstream issue (that is the other child,
[`w-mcp-dispatch-projection.md`](w-mcp-dispatch-projection.md), blocked on `ailang#885`) and not
a human decision (nothing here is parked on Mark; the split is a controller routing call); NOT
QUORUM-CLEARED; deliberately NOT sprint-ready.** The A2A projection design below is
substantially written — it was the parent's `P6.B-A2A`, reviewer-refined across four quorum
rounds — but its session-resolution responsibility rests on an API that does not exist at HEAD:
the repo has **no inbound credential→session resolution at all** (measured below, re-run
first-party this session). Charter row 39 designs that boundary; this doc CONSUMES it.  
**Item**: split child #2 of `w-mcp-projection` (charter clause 6, queue row 5; the blocker is
queue row 39 `w-session-authority`)  
**Clause**: clause-6 — the *"publish the A2A agent card"* half, plus session-scoped A2A
invocation through propose → verify → commit  
**Estimate**: **~0.7d** (the parent's P6.B-A2A line, inherited unchanged by the split) — and it
does **not start** until row 39's session-authority design lands AND the parent's
`P6.T`/`P6.D`/`P6.V` are green (see Estimate honesty)  
**Author**: rotation designer, split #2 round, 2026-08-26 (iteration 126)  
**Verified against**: World `dev` at `fcf18fa` — the session-authority probes were re-run
first-party in this session and are identical to the controller's iteration-125 values at
`2e44e3e` (`git diff 2e44e3e..fcf18fa -- host/ cmd/` is empty, measured this session); upstream
`github.com/sunholo-data/ailang` measurements are inherited from the parent at tag `v0.33.2`
(`63e7909f`) and labelled as inherited  
**Date**: 2026-08-26

## Why this doc exists — the objection, verbatim

The parent, `design_docs/planned/w-mcp-projection.md`, was re-quorumed at round 4 (2026-08-25,
both reviewers present, `absent_reviewers` empty) after split #1 carved out the MCP dispatch
half. `gemini-3-1-pro` **PASSED** — the doc's first pass in four rounds — while `gpt5-6-sol`
**REJECTED** on exactly one surface: session authority, a property of **`P6.B-A2A` alone**.
Objections localising onto one surface while another reviewer flips to pass is the skill's
decomposition signal firing as written: the doc's SCOPE was wrong, not its content. The
disposition was **SPLIT #2** — a controller routing call, explicitly not `needs-human-review` —
and this doc carries the objection and the scope it blocks. All three review fields follow,
character-for-character from the round-4 artifact
(`.ailang/state/mission-quorum/w-mcp-projection-2026-08-25T21-25-44Z.json`; the quoted text was
verified equal to the artifact's fields in this session).

**`strongest_objection` (`gpt5-6-sol`, quorum round 4 on `w-mcp-projection`, 2026-08-25):**

> The session-authority boundary is not executable: the doc equates `host/broker.NewSession(store,
> episodeID, grants, registry)` with an API that resolves an opaque Bearer credential, but N15
> verifies only a session constructor. No landed credential lookup, credential-to-episode/grants
> mapping, expiry/revocation source, or authentication API is identified. D-WORLD-26 selects the
> HTTP carrier only; it does not supply these missing semantics. Therefore the claims that
> unknown/expired credentials fail closed and that all retained milestones are executable now rest
> on an unverified codebase premise.

**`catch` (same reviewer, same round):**

> Verify whether a credential resolver/store already exists, including exact file, symbol, lookup
> key, expiry behavior, and immutable capability-snapshot semantics. Also analyze its overlap with
> existing broker/session and daemon authentication machinery. If none exists, P6.B-A2A has an
> unacknowledged prerequisite and must not be presented as immediately executable.

**`proposed_fix` (same reviewer, same round):**

> Add a premise row such as: `N27 [R] | Bearer credential resolution exists | <commands locating
> and testing exact API> | <symbol maps an opaque credential to episode ID, grants, expiry and
> registry epoch; absent/malformed/unknown/expired all return typed denial>`. Then replace
> Decision 3 responsibility 1 with that exact API and add its conflict-surface entry. If the row
> cannot be verified, change status to BLOCKED and add an independently mergeable prerequisite
> milestone defining a host-extension credential resolver with bounded lookup, explicit
> expiry/revocation, constant-time credential comparison, no API-key or alternate-header fallback,
> and tests for unknown/expired credentials; P6.B-A2A must depend on it.

The `catch` was executed first-party by the controller in iteration 125 and re-run in this
session (the measured table below): no credential resolver/store exists — under any name — so
`P6.B-A2A` had an unacknowledged prerequisite, exactly as the reviewer said. The
`proposed_fix`'s second half is literally the disposition this split executes: the status of
this scope IS **BLOCKED**, and the "independently mergeable prerequisite milestone defining a
host-extension credential resolver" IS **charter queue row 39 `w-session-authority`** — filed as
a charter queue row rather than a milestone of this doc because the gap is **pre-existing at
HEAD**: clause 3 landed the session/capability model deliberately in-process and no HTTP-facing
session authority was ever built, so this is a defect the parent *fails to fix*, not one it
*introduces*, and a pre-existing defect surfaced by a reviewer is a queue row, not a revision.
`P6.B-A2A` depends on it, exactly as the fix demands.

## The measured evidence — no inbound credential→session resolution exists

Rows marked **[R]** were re-run first-party in this split-#2 session (2026-08-26, iteration 126,
at `fcf18fa`; `host/` holds **96** Go files, scope asserted with `test -d`). Every value is
identical to the controller's iteration-125 measurement at `2e44e3e` — and
`git diff 2e44e3e..fcf18fa -- host/ cmd/` is empty, which is why. Counts are of matching LINES.

| # | Probe | Command | Result | Reading |
|---|---|---|---|---|
| S1 [R] | inbound bearer handling | `grep -rn --include='*.go' 'Bearer' host/` | **0** | nothing reads a bearer credential |
| S2 [R] | any authorization header path | `grep -rn --include='*.go' 'Authorization' host/` | **1** | a prose comment at `host/broker/approve.go:51`, not a code path |
| S3 [R] | credential machinery | `grep -rn --include='*.go' 'Credential' host/` | **128** | ALL OUTBOUND — `RegistryCredentialProvider` / `FileRegistryCredentialProvider` / `AssertNoAmbientRegistryCredential` in `host/broker/credential.go`: a credential World *presents* to an upstream registry |
| S4 [R] | authentication machinery | `grep -rn --include='*.go' 'Authenticate' host/` | **29** | ALL evidence-envelope — `AuthenticatedEnvelope` + codec in `host/evidence/`: signed evidence, not HTTP session auth |
| S5 [R] | session lookup by string | `grep -rnE --include='*.go' 'func .*(GetSession\|LookupSession\|ResolveSession\|SessionByID\|FindSession)' host/ cmd/` | **0** | **nothing in this repo resolves a session by string** |
| S6 [R] | known-positive control, SAME scope and instrument shape | `grep -rn --include='*.go' 'Session' host/` | **181** | the instrument sees the surface it is searching |
| S7 [R] | negative control, fresh absent literal, same scope | `grep -rn --include='*.go' '<fresh literal>' host/` | **0** | the instrument's zeros are measurements |
| S8 [R] | the session "API" the parent cited is a CONSTRUCTOR | `grep -n 'func NewSession' host/broker/broker.go` | `:87` — `func NewSession(s *store.Store, episodeID string, grants []Capability, registry Registry) *Session` | **the grants are an ARGUMENT** — nothing decides what grants a credential carries |

So `D-WORLD-26` (Mark, attended, issue `#89`, 2026-08-25T19:06:41Z, answer `A`) settled the
credential **envelope** — the session credential rides the standard `Authorization: Bearer`
header — and the **contents** were never built: who mints a session credential, where the
credential → (episode, grants) mapping lives, and what expires or revokes it. The reject is
correct, and it is a strictly larger finding than "the doc under-specifies a resolver".

## Blocking predicate — runnable, not prose

A future iteration decides "still blocked?" by RUNNING these, never by transcribing this doc.
Each probe carries a same-call control so a zero from a broken instrument is distinguishable
from a true zero.

```bash
# 1. Does an inbound credential→session resolver exist at HEAD?
grep -rnE --include='*.go' \
  'func .*(GetSession|LookupSession|ResolveSession|SessionByID|FindSession)' host/ cmd/ | wc -l
# measured 2026-08-26 (iteration 126) at fcf18fa: 0
# known-positive control — SAME scope, same instrument shape (proves the instrument sees the surface):
grep -rn --include='*.go' 'Session' host/ | wc -l
# measured 2026-08-26 at fcf18fa: 181
# negative control — a freshly invented absent literal reads 0, so a zero is a measurement:
grep -rn --include='*.go' 'ZZQX_ABSENT_LITERAL_9' host/ | wc -l
# measured 2026-08-26 at fcf18fa: 0

# 2. Has charter row 39 produced its session-authority design?
grep -n 'w-session-authority' design_docs/world-mission.md | head -3
ls design_docs/planned/ | grep -i 'session-authority'
# measured 2026-08-26: the row is present (world-mission.md:4245) and no design doc exists yet
# → still blocked

# 3. Are the parent's enabler milestones green? (P6.T's own observable; P6.D/P6.V per their
#    merge criteria in the parent)
grep -n 'GOTOOLCHAIN' .github/workflows/ci.yml
# measured 2026-08-26: go1.25.6 at :21 and :102 → P6.T unlanded → still blocked
```

**UNBLOCKED means ALL of:** row 39's session-authority design has landed and names the actual
resolver API (predicate 2 is the AUTHORITATIVE signal — predicate 1's function-name grep is the
instrument that measured the absence, and a resolver landed under a name outside that shape
satisfies row 39, not the grep; re-write predicate 1 against the landed symbol at that point),
AND the parent's `P6.T`, `P6.D` and `P6.V` are green (predicate 3 plus their merge criteria).
Only then is this doc revised against row 39's ACTUAL contract and quorumed at pick time.

## The dependency this doc consumes — row 39's resolver (the shape is row 39's to settle)

**This doc does NOT design the session-authority boundary.** That is row 39's job: who mints a
session credential and through what surface; where the credential → (episode, grants) mapping
lives, given that clause 3 put `broker.Session` in-process on purpose; what expires or revokes
it; lookup bounds; and comparison semantics (the `proposed_fix` names constant-time credential
comparison). What this doc states is the interface it will CONSUME, as a dependency:

- an opaque `Authorization: Bearer` credential resolves to **`(episodeID, grants, expiry)`**;
- **typed denial** on absent / malformed / unknown / expired — each distinguishable, each
  fail-closed;
- **no API-key fallback** (D-WORLD-26 constraint (i): a `Bearer` value here is a SESSION
  credential and never the static `serve-api` key, which iteration 24 measured process-wide)
  and no alternate-header fallback (Arm B `X-World-Session` is REJECTED and must not be read,
  even as a fallback);
- **no unauthenticated degrade path** (D-WORLD-26 constraint (ii); clause 3's inviolable half).

If row 39's landed contract diverges from this statement, THIS doc is revised at pick time,
before its quorum — the dependency direction is: row 39 settles the shape, this doc adapts.

## Scope carried from the parent (moved at split #2; substance unchanged)

Parent identifiers are **RETAINED** throughout — premises P2–P6, Decisions 3 and 4, acceptance
criteria AC1–AC9 and AC11–AC14, and every named mutation — so four rounds of quorum history
stay traceable. The parent keeps P1/P7/P8, AC10, and AC15–AC17; the numbering gaps on both
sides are deliberate and declared.

### Premises (hard constraints; parent identifiers preserved)

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

(The parent retains **P1** — no new wire protocol; A2A remains at `/.well-known/agent.json` and
`/a2a/` with upstream framing, no reverse-proxy re-encoding — as the clause-wide guardrail both
children cite; it binds this doc's handlers in full. **P7** (zero-cloud, measured) and **P8**
(P6.V adds exactly one law) also stay with the parent's enabler milestones, which this doc
consumes.)

### Decision 3 — Session-Scoped Projection in worldd (A2A) — carried; responsibility 1 re-bound per the round-4 `proposed_fix`

The additive `host/projection` adapter has four responsibilities. Responsibilities 2–4 name
LANDED interfaces (anchors inherited from the parent's premise row N25, re-derived 2026-08-25 at
`2e44e3e`; `host/ cmd/` unchanged since — measured this session). Responsibility 1 names the
interface row 39 must land — that re-binding is exactly the round-4 fix's demand, executed:

1. resolve the session credential — **through row 39's inbound credential resolver, which does
   not yet exist at HEAD** (S5/S6 above; this is the doc's blocker) — implementing
   `protocol.SessionResolver`. `ResolveSession(ctx, r)` reads the standard
   `Authorization: Bearer <session-credential>` header — and nothing else (**D-WORLD-26 = ARM
   A**, with both constraints carried unchanged: **(i)** never an API key, **(ii)** fail closed
   on absent/malformed/unknown/expired, never degrading to an unauthenticated surface; Arm B
   `X-World-Session` REJECTED even as a fallback) — hands the opaque credential to the row-39
   resolver, and receives `(episodeID, grants, expiry)` or a typed denial. Only then does the
   landed constructor come in: `host/broker.NewSession(store, episodeID, grants, registry)`
   (`broker.go:87`) CONSTRUCTS the session from the resolver's output. It is not, and never
   was, a credential resolver — the parent's use of it in this position is precisely what the
   round-4 objection caught, and this rewrite is the correction.
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

### Decision 4 — Surface Identity, A2A Card, and Compatibility — carried unchanged

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

The cross-surface obligation the pre-split AC3 stated — that the eventual MCP `tools/list` name
set must exactly equal the A2A `skills[].id` set per session — is NOT dropped: it travels with
the MCP dispatch child (`w-mcp-dispatch-projection.md`) and binds the MCP surface when it
becomes buildable.

The agent card is emitted by World's A2A handler using `protocol`'s wire representation
(`AgentInfo`, `A2ARequest`/`A2AResult`/`A2AError` and the descriptor set from `CallerSurface`).
The test must not assume a hand-authored JSON shape beyond protocol-required fields and the exact
skill-ID set.

### Bounded waits and the commit boundary — carried from the parent's Decision 6

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
It is never mapped to `daemon.APIError`, retried, or used as a reason to fall back to REST,
direct store access, or another capability/session source.

The reviewer-authored commit-boundary CONTRACT (quorum round 2, `gpt5-6-sol`, adopted verbatim)
and its grounding in landed public APIs (`JournalIntent`/`bindCommitIntentTx`/`InvocationID`/
`GetReceipt`/`GetEffectReceipt` — parent premise row N14) STAY in the parent, because the
parent's P6.V is what proves it as a Z3-verified law. This doc CONSUMES that verified law:
AC13 below tests both sides of the boundary against it. There is NO SSE stream in this doc —
A2A is single bounded request/response within the frozen D7 deadlines — but the OS-level kernel
of quorum r2's socket objection survives here: a disconnected or deadline-expired transport must
be OBSERVED closed (`http.Server.ConnState` tracking or a client read error), never merely
context-cancelled (AC13). The stream-lifetime maximum and route-local deadline relaxation are
the MCP dispatch child's.

### Files (the rows moved from the parent's table)

| File | Milestone | Purpose |
|---|---|---|
| `host/projection/projection.go` | P6.B-A2A | Session-scoped adapter; World-owned A2A handlers + agent card over `protocol`; session resolver binding (`Authorization: Bearer` → row-39 resolver) |
| `host/projection/projection_test.go` | P6.B-A2A | Exact-set, denial, wire-form, bounded-wait, socket-closure tests |
| `host/daemon/daemon.go` | P6.B-A2A | Additive mounting of the World-owned A2A handlers |
| `host/daemon/daemon_test.go` | P6.B-A2A | REST/single-writer/dependency regression half (the P6.D allowlist-line half stays with the parent) |
| `cmd/ailang-worldd/main.go` and tests | P6.B-A2A | Minimal flags/config for enabling projection and session credential policy |

## Milestone P6.D — Dependency admission, ATOMIC WITH THIS DOC'S FIRST REAL CONSUMER (~0.15d) — INHERITED FROM THE PARENT AT ITS ROUND-5 CARVE-OUT (2026-08-26, iteration 126)

**Why this milestone is here and not in the parent.** Quorum round 5 rejected the parent for
pre-landing this dependency behind a dead compile anchor: *"P6.D deliberately lands an otherwise
unused direct dependency and a production file whose only purpose is a dummy symbol reference …
This is speculative core growth and directly violates the minimal-frozen-core and
route-to-extension axioms"* (`gpt5-6-sol`, verbatim). Its prescribed fix — *"Move the pinned
dependency and narrow package-path allowlist change into whichever child first becomes unblocked,
where a real handler or adapter import provides the compile-visible use"* — is what puts it here.
**This doc's `host/projection` A2A handler import IS that real consumer.**

**Ordering caveat, stated because it is a race and not a certainty.** Two children can unblock:
this one (on charter queue row 39 `w-session-authority`) and
[`w-mcp-dispatch-projection.md`](w-mcp-dispatch-projection.md) (on
[`ailang#885`](https://github.com/sunholo-data/ailang/issues/885)). **Whichever unblocks FIRST
carries this milestone; the other inherits an already-admitted dependency and adds nothing.**
Before starting P6.D here, check whether it has already landed —
`git grep -n 'serveapi/protocol' -- go.mod host/daemon/daemon_test.go` — and if it has, skip this
milestone and record the skip.

The specification, carried from the parent unchanged in substance:

`go get github.com/sunholo-data/ailang@v0.33.2`, then add **exactly one** `allowedDepModules`
entry: the **PACKAGE path** `github.com/sunholo-data/ailang/serveapi/protocol` — never the
module root. The matcher (`disallowedDeps`, `daemon_test.go:801`:
`d == m || strings.HasPrefix(d, m+"/")`) treats entries as path prefixes, so the package-path
entry admits exactly that package (and any future subpackage beneath it), while the module root
would admit `internal/apiserver`'s measured 476 disallowed packages. Update the
`allowedDepModules` doc comment ("module roots") to "module roots or package paths".

- **Requires the parent's `P6.T` to be green first**: `v0.33.2` declares `go 1.26.6` and this
  repo pins `go1.25.6` until P6.T lands (two-arm measured 2026-08-25: `go get ...@v0.33.2` rc=1
  under `GOTOOLCHAIN=go1.25.6` with the exact floor message, rc=0 under `go1.26.6`;
  known-positive control `go get github.com/google/uuid@v1.6.0` rc=0 under `go1.25.6`).
- Include a **narrowness test**: with the new entry in place,
  `disallowedDeps(["github.com/sunholo-data/ailang/internal/apiserver"])` (plus a
  representative cloud path) is non-empty — proving the entry did not widen the gate.
- **The compile-visible use is the real handler import**, not an anchor file: `host/projection`'s
  first `protocol` import supplies it. `host/daemon/protocol_use.go` — the parent's
  one-iteration-old dead anchor — is **withdrawn and is not created by any doc.** Closure over
  both gated patterns moves **249 → 250**, addition = the package itself, removed set empty
  (inherited measurement, parent premise rows, iteration 120/124).
- Named risk with its own acceptance row (AC16): the pin also bumps `go-isatty → v0.0.22` and
  `x/sys → v0.47.0` (already-allowlisted roots) and `go.sum` moves; no other graph movement is
  permitted.
- Merge criterion: both CI jobs green; `TestDaemonDependencyAllowlist` green with the probe
  demonstration recorded (allowlist minus the new line → REDs naming exactly one intruder).

Files: `go.mod` (`require`), `go.sum` (pin + the two measured indirect bumps),
`host/daemon/daemon_test.go` (ONE package-path allowlist line + the narrowness test).

- [ ] **AC16 — pinned dependency, narrow admission (P6.D; identifier retained from the parent):**
  `go.mod` requires `github.com/sunholo-data/ailang v0.33.2`; `allowedDepModules` gains EXACTLY
  ONE entry, the package path `github.com/sunholo-data/ailang/serveapi/protocol`; the closure
  over both gated patterns moves by exactly +1 package (the package itself, removed set empty);
  the narrowness test proves `internal/apiserver` (and a representative cloud path) is still
  refused; and the ONLY other module-graph movement is the two measured indirect bumps
  (`go-isatty v0.0.20 → v0.0.22`, `x/sys v0.46.0 → v0.47.0`) — any third movement fails the row.

| Gate | Named RED mutation (concrete edit) | Required red observation |
|---|---|---|
| AC16 `MUT-ALLOWLIST-ROOT` | replace the package-path entry with the module root `github.com/sunholo-data/ailang` | the narrowness test REDs: `internal/apiserver` is no longer refused by `disallowedDeps` |
| AC16 `MUT-FACADE-IMPORT` | import `github.com/sunholo-data/ailang/serveapi` (the facade) at `host/projection`'s protocol-import site | `TestDaemonDependencyAllowlist` REDs naming the intruding packages by name (measured scale: 476 across 86 roots) |

## Milestone P6.B-A2A — Session projection + A2A conformance, World-owned handlers (~0.7d)

The A2A half of the pre-split P6.B, carried from the parent unchanged in substance. **Starts
only after row 39's session-authority design lands, the parent's P6.T and P6.V are green, and
P6.D above has landed (here or in the sibling child, whichever unblocked first).** Add `host/projection` with a **World-authored** A2A handler and agent-card endpoint
built exclusively on `protocol` types/helpers, mounted additively in the daemon; connect only to
the landed `transitionreg`/coordinator interfaces plus row 39's resolver. The session resolver
reads `Authorization: Bearer` per D-WORLD-26 = A and fails closed. Upstream's own
`a2a_handler.go` (**180 lines, stdlib + `protocol`, 0 SDK imports** — inherited, parent premise
row N23, re-derived 2026-08-25 via `gh api` at tag `v0.33.2`) is the existence proof of exactly
this shape.

- Prove exact two-session discovery/card/invocation behavior, stale-name denial, A2A wire-form
  conformance against the frozen P6.A fixture's A2A leg (the fixture is frozen in the parent's
  P6.A record; its MCP leg travels with the MCP dispatch child), bounded cancellation through
  every dependency, both sides of the commit boundary (AC13), OS-level socket closure on
  disconnect/expiry, REST deadline regression, single-writer preservation, and the dependency
  allowlist (AC11).
- The MCP half — `/mcp/` handler, JSON-RPC dispatch, SSE stream lifetime — is NOT here and does
  not gate this milestone; it is the MCP dispatch child's scope, blocked on `#885`.
- Merge criterion: both CI jobs green, no skips — **zero skipped protocol tests** (this is the
  protocol-test clause of the parent's AC10, travelling here with the tests it governs;
  `MUT-SKIP-SOCKET` is its mutation) — and every acceptance mutation demonstrated RED then
  reverted GREEN.

## Acceptance Criteria (parent identifiers RETAINED; AC10, AC15 and AC17 remain in the parent; **AC16 moved HERE with P6.D at the parent's round-5 carve-out** and is stated in the P6.D milestone above — the gaps are deliberate)

- [ ] **AC1 — wire-contract reuse:** World's A2A handler and agent card construct and parse wire
  forms exclusively through the pinned `serveapi/protocol` types and helpers (`ToolDescriptor`,
  `CallerSurface`/`AuthorizedSurface`, `ValidateMCPName`, `AgentInfo`,
  `A2ARequest`/`A2AError`/`A2AResult`, `AuthorizationError`/`AuthorizationStatus`); production
  World code declares no parallel wire struct and hand-formats no envelope bytes. (World OWNS the
  `http.Handler`s — that is the parent's Decision 2 scope, not a violation of this row.)
- [ ] **AC2 — exact surface:** for two non-empty sessions with unequal capability sets, the A2A
  `skills[].id` set equals each session's authorized transition-ID set exactly; no extras.
- [ ] **AC3 — card tracks the session:** a session's A2A skill-ID set changes when its capability
  snapshot changes; two concurrent sessions with different snapshots observe different cards from
  the same daemon. (The cross-surface MCP-equality half of the pre-split AC3 travels with the
  MCP dispatch child.)
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
  against a re-derived local shape. (The MCP SSE-framing half of the pre-split AC8 travels with
  the MCP dispatch child.)
- [ ] **AC9 — landed behavior preserved:** REST v1 route/body regressions and cross-process
  writer-lock tests remain green; projection never opens a store.
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
  values (constants `daemon.go:78-97`, wired in `newServer` at `:619-628` — anchors inherited
  from parent row N25 at `2e44e3e`; `host/` unchanged since, measured this session), including
  its 30s read/write deadlines and 120s idle timeout. (The route-local SSE relaxation this row
  FORBIDS here is exactly what the MCP dispatch child must design for `/mcp/` when it unblocks.)

## Non-Vacuity — Named RED Mutation for Every Gate (carried; parent identifiers retained)

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
| zero-skip (from parent AC10's protocol-test clause) `MUT-SKIP-SOCKET` | add `t.Skip` on listen failure | zero-skip CI assertion/source check REDs |
| AC11 `MUT-CLOUD-DEP` | add `cloud.google.com/go/storage` to projection imports | `TestDaemonDependencyAllowlist` reports it by name |
| AC12 `MUT-STARTUP-CACHE` | cache transition descriptors once at daemon startup | registry-head-change test returns stale set and REDs |
| AC13 `MUT-DROP-DEADLINE` | replace the propagated request context with `context.Background()` before one injected blocking dependency | bounded-wait fault-injection test exceeds the configured bound and REDs; mutation cleanup is terminated by the test harness |
| AC13 `MUT-COMMIT-BOUNDARY-LIE` | move the cancellation check from *before* coordinator acceptance to *after* it, so a cancelled-mid-commit request reports a definitive "not committed" | the after-boundary test REDs: a durable receipt exists under the invocation/idempotency ID while the caller was told the commit did not happen |
| AC13 `MUT-LEAK-CONN` | on client disconnect or deadline expiry with a blocked response write, cancel the Go context but never close the `http` connection | the `ConnState`/client-read-error assertion REDs — proving the test observes OS-level socket closure, not just logical cancellation |
| AC14 `MUT-DEADLINE-RELAX` | add a `ResponseController.SetWriteDeadline(time.Time{})` call to the A2A handler | the source-level assertion REDs naming the call site; the REST deadline regression is the same-run control that the D7 constants did not move |

## Conflict Surface (the entries moved from the parent)

- **HTTP Server Timeouts vs the A2A route.** The existing daemon has frozen D7
  `ReadHeaderTimeout` 5s, `ReadTimeout` 30s, `WriteTimeout` 30s, and `IdleTimeout` 120s
  (constants at `host/daemon/daemon.go:78-97`, wired in `newServer` at `daemon.go:619-628` —
  anchors inherited from parent row N25 at `2e44e3e`; `host/` unchanged since, measured this
  session). This doc adds ONLY plain bounded request/response endpoints
  (`/.well-known/agent.json`, `/a2a/`), which run entirely WITHIN those deadlines: no route in
  this doc's scope relaxes a deadline, and `host/projection` contains no
  `ResponseController.SetWriteDeadline`/`SetReadDeadline` call (AC14's source assertion). The
  long-lived-SSE-vs-30s-write-deadline conflict is REAL but is the MCP dispatch child's conflict
  surface, recorded there. The D7 constants and every REST `/v1/*` path remain byte-for-byte
  unchanged.
- **Row 39's resolver vs existing broker/session and daemon machinery.** The round-4 `catch`
  demands an overlap analysis between any credential resolver and the existing broker/session
  and daemon authentication machinery. That analysis belongs to row 39 (it designs the
  resolver); what is measured HERE is the seam this doc will consume: `broker.NewSession`
  (`broker.go:87`) constructs sessions from resolver output, `host/broker/credential.go` is
  OUTBOUND-only registry-publish credentials (S3), and `host/evidence/`'s `Authenticate*` is
  evidence-envelope signing (S4) — neither is, nor must be conflated with, inbound HTTP session
  auth. This doc adds no second policy engine and no credential store.
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
  neither renames nor overloads either.
- **Broker/propose-verify-commit.** Direct REST commit and store commit are explicitly outside the
  projection. If the landed broker lacks a single session-aware dispatch entrypoint, P6.B-A2A
  stops; it does not synthesize one in the adapter.
- **Protocol socket tests on CI.** Protocol socket tests must run on CI; local sandbox denial is
  not a skip condition (zero-skip merge criterion, `MUT-SKIP-SOCKET`).
- **Performance baseline.** No existing benchmark is removed. P6.B-A2A adds bounded protocol
  round-trip benchmarks only if measurement shows a material new hot path; this item does not
  rewrite `bench/BASELINE.md`.

## Axiom Compliance (inherited from the parent's pre-split table — the projection substance is unchanged; re-scored at unblock if row 39's contract shifts it)

| Axiom | Score | Justification |
|---|---:|---|
| A1 Determinism | +1 | stable-ID ordering and one snapshot per request |
| A2 Replayability | +1 | calls use the normal recorded transaction path |
| A3 Effect Legibility | +1 | protocol I/O stays at the host boundary |
| A4 Explicit Authority | +2 | exact per-session card/call predicate; carrier fixed by D-WORLD-26, fail closed; resolution delegated to row 39's designed boundary, never improvised here |
| A5 Bounded Verification | +1 | existing verify pipeline; no protocol bypass |
| A6 Safe Concurrency | +1 | daemon retains sole handle; immutable request snapshot |
| A7 Machines First | +2 | A2A schemas and structured upstream errors |
| A8 Minimal Syntax | 0 | no language syntax change |
| A9 Cost Visibility | 0 | no new budget claim; benchmark scope deferred honestly |
| A10 Composability | +2 | standard A2A, upstream-owned wire contract; the MCP≡A2A equality obligation is preserved in the MCP dispatch child |
| A11 Structured Failure | +1 | missing/expired/denied sessions fail in protocol form, with row 39's typed denials underneath |
| A12 System Boundary | +2 | projection is an adapter, never kernel state |

**Net: +14; hard axioms A1/A3/A4/A7 are non-negative.**

## Premise rows carried from the parent (inherited, with their original dates and iterations)

Rows below keep their parent identifiers. **[I]** = inherited exactly as measured in the
parent's 2026-08-25 logs (revision round and split round, iteration 124/125); the
session-authority surface they touch was ADDITIONALLY re-measured first-party in this session
(the S-table above, 2026-08-26 at `fcf18fa`).

| # | Premise | Original evidence | Status here |
|---|---|---|---|
| N15 [I] | Transition registry + broker session API landed | `ls host/transitionreg/` → `transitionreg.go`, `bind.go`, tests; `grep -n 'func NewSession' host/broker/broker.go` → `:87` (parent revision round, 2026-08-25) | Inherited; S8 re-read the signature this session. NOTE the round-4 sharpening: N15 verifies a session CONSTRUCTOR, not a credential resolver — that gap is row 39 |
| N16 [I] | No HTTP session-credential carrier exists in `host/` | `grep -rn 'Bearer\|Authorization' host/` → only a prose comment at `host/broker/approve.go:51` (parent revision round, 2026-08-25) | Inherited; re-run this session as S1/S2 with identical values. The CARRIER was since chosen (D-WORLD-26 = A); the RESOLVER remains unbuilt (row 39) |
| N17 [I] | D7 constants/wiring anchors | constants `host/daemon/daemon.go:78-97`, wiring `newServer` at `:619-628` (parent revision round, re-derived split round at `2e44e3e`) | Inherited; `host/` unchanged `2e44e3e..fcf18fa` (measured this session) |
| N23 [I] | A2A existence proof at `v0.33.2`: `a2a_handler.go` is 180 lines, stdlib + `protocol`, 0 SDK imports | `gh api 'repos/sunholo-data/ailang/contents/serveapi/a2a_handler.go?ref=v0.33.2'`, decoded (parent split round, 2026-08-25, iteration 125) | Inherited, UNVERIFIED this session (would be settled by re-running the same `gh api` call); also retained in the parent as Decision 2's measurement |
| N26 [I] | D-WORLD-26 RESOLVED = ARM A with both constraints | charter ledger row read first-party in the parent's split round (`design_docs/world-mission.md:708`, 2026-08-25): answered `A` by Mark, attended, issue `#89`, `2026-08-25T19:06:41Z` | Inherited; its operative force (resolver contract, AC6, `MUT-ALT-HEADER`/`MUT-KEY-AS-SESSION`) binds THIS doc |

## Estimate honesty

**~0.7d**, the parent's split-round P6.B-A2A line carried unchanged — the split moved the scope,
it did not re-price it. Plainly: **this work does not start now.** It starts only when row 39's
session-authority design has landed AND the parent's `P6.T`/`P6.D`/`P6.V` are green, and the
estimate is contingent on row 39's landed contract matching the dependency stated above — if it
diverges, the pick-time revision re-prices before quorum. No part of this 0.7d is on the
parent's critical path (the parent's remainder is ~0.55d and waits on nothing).

## Quorum status

**This doc is NOT quorum-cleared.** It was authored at split #2 (iteration 126) and no reviewer
has passed it — the round-4 quorum blocked the PARENT, and this doc is the carrier of that
block's surviving objection, not a resolution of it. When the blocking predicate above reads
unblocked and this item is picked, this doc is revised against row 39's ACTUAL landed contract
and MUST then go through the full design quorum (`ailang design-quorum`, reject-by-default
synthesis) at pick time. Nothing in this doc pre-authorizes a sprint.

## Relationship to the parent and to charter clause 6

Charter clause 6 names *"project the transition registry over MCP + publish the A2A agent
card."* After split #2 the clause is partitioned across THREE docs:

- the parent, [`w-mcp-projection.md`](w-mcp-projection.md) — the three enabling milestones
  (`P6.T` toolchain floor, `P6.D` pinned `serveapi/protocol` dependency + one narrow allowlist
  line, `P6.V` verified commit-boundary law), objection-free across all four quorum rounds,
  executable now, and consumed by BOTH children;
- child #1, [`w-mcp-dispatch-projection.md`](w-mcp-dispatch-projection.md) — the MCP dispatch
  half, blocked UPSTREAM on
  [`sunholo-data/ailang#885`](https://github.com/sunholo-data/ailang/issues/885);
- child #2, THIS DOC — the A2A projection half, blocked LOCALLY on charter queue row 39
  `w-session-authority`.

**Clause 6 is satisfied only when all three land.** No doc narrows the clause; the splits
partition it, each blocked half named with a runnable predicate.

## Related Documents

- [w-mcp-projection.md](w-mcp-projection.md) — the SPLIT parent: enabler milestones, executable
  now; its quorum log records both split dispositions
- [w-mcp-dispatch-projection.md](w-mcp-dispatch-projection.md) — split child #1: the MCP
  dispatch half, blocked on `ailang#885`; carries the cross-surface MCP≡A2A equality obligation
- [world-mission.md](../world-mission.md) — clause 6; **queue row 39 `w-session-authority`**
  (this doc's blocker, filed iteration 125 on the controller's first-party evidence);
  **D-WORLD-26** (the carrier ruling whose operative force binds this doc); D-WORLD-5
- [coding-standards.md](../coding-standards.md) — S1–S6
- [DESIGN.md](../DESIGN.md) — §3.7 protocol-native boundary, §14, §17
- [w-worldd-m2.md](../implemented/w-worldd-m2.md) — shipped daemon, REST freeze, D7 constants
- [w-effect-broker-m3.md](../implemented/w-effect-broker-m3.md) — landed session/capability/
  broker (the in-process model row 39 must extend to an HTTP-facing boundary)
- [w-store-durability.md](../implemented/w-store-durability.md) — the commit-boundary Go surface
  whose Z3 proof (parent P6.V) this doc consumes
