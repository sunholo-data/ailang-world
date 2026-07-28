# w-mcp-projection — Session-Scoped MCP Projection and A2A Agent Card

**Status**: Planned — **BLOCKED on the upstream serving seam and the landed transition registry**
**Item**: `w-mcp-projection`  
**Clause**: clause-6 (protocol-native boundary; not DESIGN.md M6)  
**Estimate**: ~1.2 World days after prerequisites; upstream work is not included  
**Author**: DESIGN-DOC-CREATOR  
**Verified against**: **AILANG v0.30.0**, commit `e37b370` (clean), at
`/tmp/ailang-v0300/ailang`  
**Date**: 2026-07-28

> **Scheduling truth.** This is sprint-ready as a conditional plan, but it is not executable
> against v0.30.0. The shipped binary cannot provide an exact per-session tool surface, and this
> repo does not yet contain the clause-3 transition registry or its session authority. Milestone
> P6.A may land independently now; P6.B starts only after both prerequisites are real and pinned.
> Calling the current static sketch exports a dynamic transition registry would be a defect.

---

## Motivation

Clause 6 requires the available World transitions to leave the daemon through existing protocols:
MCP for callable tools and A2A for discovery. It also requires the MCP surface to be filtered by
the capabilities of each session. This is an authority boundary, not documentation: an
unauthorized transition must be absent from `tools/list`, absent from the A2A card, and rejected
if invoked by name.

AILANG already owns the MCP and A2A wire implementation. World must reuse that machinery, not
write JSON-RPC, SSE, MCP, or A2A codecs. Live v0.30.0 evidence, however, rules out treating the
CLI flags as the completed boundary:

- `--caps` is process-wide and does not filter discovery;
- one static API key is not a World session;
- `--routes-only` removes the eight embedded `std/io` exports but still publishes the built-in
  `submit_feedback` tool;
- module loading is cwd-sensitive; and
- static `.ail` exports are not the dynamic, worldd-backed transition registry.

The design therefore takes charter reuse path **(c)**: request a narrow public serving seam from
the existing `internal/apiserver`, then let worldd supply session-scoped discovery and invocation
callbacks. The upstream package remains the protocol owner. World remains the state, session,
capability, and transaction owner.

## Premises (hard constraints)

- **P1 — no new wire protocol.** MCP Streamable HTTP remains at `/mcp/`; A2A remains at
  `/.well-known/agent.json` and `/a2a/`. Responses keep upstream framing, including SSE-framed
  MCP HTTP responses. World does not decode and re-encode protocol traffic in a reverse proxy.
- **P2 — exact per-session authority.** A session resolves to an immutable capability snapshot
  for one request. The same predicate filters MCP discovery, A2A skills, and invocation. A
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
  MCP/A2A are additive protocol endpoints and do not translate failures into the REST
  `daemon.APIError` envelope. Existing REST callers see byte-compatible behavior.
- **P7 — local-first and zero-cloud remain enforced.** Prefer zero new dependencies. Path (c)
  necessarily requires a public upstream serving package; its eventual module/version and full
  transitive graph must be pinned and admitted by `TestDaemonDependencyAllowlist` only after a
  first-party `go list -deps` audit proves that graph contains no cloud SDK.
- **P8 — no new AILANG law is needed.** This item adds a protocol adapter over laws owned by the
  transition registry and broker. Its invariants are cross-request/session behavior and belong in
  Go integration/conformance tests. There is no `.ail` sketch, no new Z3 obligation, and no
  `scripts/verify_ail.sh` manifest change.

## High-Impact Decisions

| Decision | Why High Impact | Chosen By | Deadline | Change Cost |
|---|---|---|---|---|
| Choose reuse path (c), not (a) or (b) | v0.30.0 cannot express exact per-session discovery/invocation | live evidence | before sprint | high |
| Upstream owns MCP/A2A codecs; World supplies callbacks | satisfies DESIGN §3.7 without moving World state upstream | clause 6 | upstream API | high |
| One authorization predicate drives list/card/call | prevents discovery and invocation drift | clause 3/6 | compile | high |
| Projection is additive inside worldd | preserves single writer and frozen REST v1 | shipped M2 | route wiring | medium |
| No new `.ail` sketch | the new law is host-boundary session behavior, not pure kernel logic | S1/S2/S3 | design | low |

### Design Freeze

- [ ] Do not implement MCP, SSE, JSON-RPC, or A2A codecs in this repository.
- [ ] Do not graft the World store or scheduler onto `ailang serve-api`.
- [ ] Do not use `--caps` as evidence of per-session filtering.
- [ ] Do not use a single `--api-key-*` value as a World session model.
- [ ] Do not call the static sketch export list a transition registry.
- [ ] Do not expose `submit_feedback`, `std/io`, introspection helpers, or any tool not present in
  the session's exact authorized transition set.
- [ ] Do not add a direct MCP/A2A-to-`store.Commit` or MCP/A2A-to-REST-commit path.
- [ ] Do not start P6.B before the transition registry/broker session API, the **verified
  commit-boundary contract** (atomic not-started-versus-committed, stable invocation/idempotency ID,
  queryable durable receipt), and the upstream serving seam have landed at pinned revisions.
- [ ] Do not claim that cancellation prevents a durable mutation once the coordinator has ACCEPTED a
  commit, and never return a definitive "not committed" result when the outcome is unknown.
- [ ] Do not treat logical `context` cancellation as proof that an SSE connection closed — the
  socket itself must be observed closed.
- [ ] Do not relax or change the frozen D7 daemon timeout constants; any SSE deadline handling is
  route-local to `/mcp/`.
- [ ] Do not leave any projection request, dependency wait, or streaming connection lifetime
  unbounded.

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

## Decision 2 — The Upstream Serving Seam

The upstream request is behavioral, not a speculative Go signature. The public seam must:

1. mount upstream's MCP HTTP and A2A handlers on a caller-owned `http.ServeMux` (or return
   `http.Handler`s);
2. ask the caller for the request's principal/session before discovery or invocation;
3. ask the caller for the exact tool descriptors visible to that session;
4. route a named invocation back to the caller with the same resolved session;
5. generate both MCP tool descriptors and A2A skills from that caller-supplied set;
6. expose no built-in tool unless the caller explicitly supplies it; and
7. preserve upstream MCP/A2A conformance and SSE framing.

The seam must not accept a World store, broker, or scheduler. It is a protocol-serving package,
not a World feature in the language repository.

**Version discipline.** The design does not invent a nonexistent release number. P6.B pins the
first released/tagged upstream revision that contains this seam and records its commit in
`go.mod`/`go.sum` and the implementation log. Until then, the dependency is **UNVERIFIED** and
P6.B is blocked.

## Decision 3 — Session-Scoped Projection in worldd

The additive `host/projection` adapter has four responsibilities:

1. resolve the configured session credential through the landed broker/session API;
2. read one transition-registry snapshot and one capability snapshot;
3. derive an exact ordered tool set using the registry's stable transition identity; and
4. dispatch an authorized invocation into the normal propose → verify → commit coordinator.

The authorization function is conceptually:

`visible(session, transition) = session is live AND transition capability ∈ session capabilities`.

That expression is explanatory prose, not AILANG source. The implementation reuses the landed
capability predicate; it does not create a second policy engine in `host/projection`.

**One snapshot per request.** Discovery/card generation may observe a newer registry on the next
request, but a single request never mixes registry or capability epochs. Invocation re-resolves
the session and transition at dispatch time; possession of a previously listed tool name conveys
no authority.

**Fail-closed cases.** Unknown/expired session, absent transition, capability denial, registry
read error, or broker unavailability returns the upstream protocol's structured error form. No
case falls back to REST or direct store access.

## Decision 4 — Surface Identity, A2A Card, and Compatibility

Tool identity is the transition registry's stable transition ID. Display text and schemas are
projection metadata; they do not become authority. Ordering is deterministic by stable ID so
cards and tool lists are diffable.

For a given session snapshot:

- `tools/list` names equal the authorized registry ID set exactly;
- A2A `skills[].id` equals the same set exactly;
- every listed tool is invocable through the same session;
- every unlisted tool is rejected even when invoked by a guessed/stale name; and
- `std/io.*`, bare IO names (`exit`, `writeBytes`, etc.), and `submit_feedback` are absent unless
  World someday registers and authorizes transitions with those exact IDs. This item registers
  none.

The agent card uses the upstream A2A representation. The test must not assume a hand-authored JSON
shape beyond protocol-required fields and the exact skill-ID set.

## Decision 5 — Process, Endpoint, and Package Layout (S2/S3)

| Path | Change |
|---|---|
| `host/projection/` | New host-boundary adapter and tests |
| `host/daemon/daemon.go` | Additive mounting of upstream MCP/A2A handlers |
| `host/daemon/*_test.go` | Existing REST regression, dependency allowlist, live protocol tests |
| `cmd/ailang-worldd/` | Minimal flags/config for enabling projection and session credential policy |
| `go.mod`, `go.sum` | Pin the released upstream serving module only after dependency audit |

**Why is this not an AILANG package?** Protocol serialization, HTTP request authentication, and
connection lifecycle are host effects (S2). The transition definitions and policy remain
package/registry data; only their OS/network projection is Go.

**Why is this not kernel growth?** `host/projection` is an adapter at DESIGN §3.7's system
boundary. It owns no transition semantics, store schema, registry update, or capability law.
Nothing is added to `world/`, `host/store`, or `host/registry`.

**Dependencies.** Zero new cloud dependencies. One new direct dependency is expected only because
path (c) exports upstream's existing server. Its module path, version, and transitive packages are
**UNVERIFIED until upstream publishes it**. The executor must not weaken the allowlist broadly;
it adds only individually audited local/compiler/protocol modules, or stops and escalates.

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
above. **P6.B is additionally BLOCKED on a verified commit-boundary contract** (see the premise row
`Commit-boundary contract`): the acceptance tests below depend on semantics the landed code does
not yet expose.

Long-lived `/mcp/` SSE responses use a second configured positive, finite maximum stream lifetime.
The connection closes when the earliest of client cancellation/deadline, the request-operation
deadline while an operation is active, or that stream-lifetime deadline occurs. The stream cap is
not a retry budget and cannot extend an operation's bounded context.

## Milestones (each independently CI-green and mergeable; ~1.2 World days)

### P6.A — Upstream repro + executable conformance fixture (~0.2d)

- File the Upstream Finding below against exact v0.30.0 with the stdio reproduction.
- Add no production endpoint and no dependency.
- Freeze a Go test fixture/expected sets for two synthetic sessions, ready to run against the
  eventual public seam. If the seam is not yet available, keep this fixture in the design/issue,
  not as a skipped or expected-failing CI test.
- Merge criterion: documentation-only World change; both existing CI jobs remain green.

### P6.B — Session projection + MCP/A2A conformance (~1.0d after prerequisites)

**Blocked on THREE prerequisites**, all outside this milestone: (1) the upstream serving seam of
Decision 2, released and tagged; (2) the clause-3 transition registry + broker session API; and
(3) the **verified commit-boundary contract** of Decision 6 (round-2 objection).

- Pin the released upstream serving package after full dependency audit.
- Add `host/projection`, mount the additive handlers, and connect only to the landed transition
  registry/session/broker interfaces.
- Prove exact two-session discovery/card/invocation behavior, stale-name denial, SSE framing,
  bounded cancellation through every dependency, route-local SSE lifetime handling, REST
  deadline regression, single-writer preservation, and dependency allowlist.
- Merge criterion: both CI jobs green, no skips, every acceptance mutation demonstrated RED then
  reverted GREEN.

**Estimate honesty.** World work moves from ~1.0d to ~1.2d: P6.B moves from ~0.8d to ~1.0d for
deadline propagation, cancellation cleanup, route-local SSE deadline control, and the associated
fault-injection/regression tests. Upstream design/review/release latency and the
clause-3 prerequisite are outside that estimate and may dominate calendar time. If the upstream
seam requires more than a narrow export/callback refactor, cut P6.B from this item and keep only
P6.A; do not implement a local codec as schedule compensation.

## Files to Create/Modify

| File | Milestone | Purpose |
|---|---|---|
| `design_docs/planned/w-mcp-projection.md` | P6.A | This design and verified repro |
| upstream issue in `sunholo-data/ailang` | P6.A | Public serving seam request |
| `host/projection/projection.go` | P6.B | Session-scoped adapter |
| `host/projection/projection_test.go` | P6.B | Exact-set and denial tests |
| `host/daemon/daemon.go` | P6.B | Additive upstream handler mount |
| `host/daemon/daemon_test.go` | P6.B | REST/single-writer/dependency regression |
| `cmd/ailang-worldd/main.go` and tests | P6.B | Minimal enablement/config |
| `go.mod`, `go.sum` | P6.B | Exact released upstream pin |

No `.ail`, schema, benchmark baseline, REST route, or CI workflow file changes.

## Conflict Surface

- **HTTP Server Timeouts vs SSE.** The existing daemon has frozen D7
  `ReadHeaderTimeout` 5s, `ReadTimeout` 30s, `WriteTimeout` 30s, and `IdleTimeout` 120s
  (`host/daemon/daemon.go:77-91`, wired at `daemon.go:409-414`). The `/mcp/` streaming handler
  alone uses `http.NewResponseController(w).SetWriteDeadline(time.Time{})` and
  `SetReadDeadline(time.Time{})` to relax the global read/write deadlines for an established SSE
  stream. This is a per-connection relaxation on one route, not a change to the server defaults:
  the D7 constants and every REST `/v1/*` path remain byte-for-byte unchanged. Because
  `IdleTimeout` governs idle time between requests rather than an active SSE handler, it is not
  the stream-lifetime bound. Decision 6's explicit finite stream maximum closes the active
  connection, while its bounded request context limits each operation and client cancellation
  ends both. A controller error fails the stream closed in upstream structured protocol form; it
  does not continue with an accidentally inherited or unbounded deadline.
- **worldd single-writer flock.** Projection uses the daemon's handle; no `store.Open`, sidecar
  writer, or lock change. Existing cross-process `WriterAlreadyActive` tests remain green.
- **Frozen REST v1 route table.** The shipped mux patterns remain unchanged:
  health, head, world, object, log-entry, log-range, registry wildcard, and commit. MCP/A2A are
  additive at upstream-standard paths; REST response bytes remain regression-tested.
- **Shared JSON error envelope.** `/v1/*` continues using `daemon.APIError`. MCP/A2A failures use
  upstream protocol errors; mapping them into `APIError` would conflate protocols.
- **`host/registry` and `world/epoch-registry/v1`.** The landed registry is an interpreter epoch
  registry, not the transition registry. This item neither renames nor overloads it. P6.B waits
  for a separately landed transition-registry interface.
- **`scripts/verify_ail.sh`.** Exact current floor is **4/4 required verified identities across 9
  modules and 14/14 named tests**. No sketch means no manifest or count change.
- **`scripts/verify_go.sh`.** Continues to run `go build ./...` and unskipped `go test ./...`
  against pinned v0.30.0. Protocol socket tests must run on CI; local sandbox denial is not a skip
  condition.
- **CI workflow.** The single `CI` workflow has exactly two jobs:
  `ailang-code verify gate` and `go host build + test gate`; the latter also runs benchmark smoke.
  No new workflow/job is required.
- **Dependency allowlist.** `TestDaemonDependencyAllowlist` enforces zero-cloud over the real
  transitive daemon/cmd build graph. The upstream package must be enumerated and pinned, and every
  non-stdlib transitive module audited. A blanket `github.com/sunholo-data/ailang/...` exception
  without printed graph evidence is forbidden.
- **serve-api cwd sensitivity.** Any P6.A reproduction that loads `sketches/` runs with
  `design_docs/` as cwd. The charter's 2026-07-23 live-test row does not record its cwd, so it is
  insufficient evidence for launch configuration. P6.B embeds handlers and therefore has no
  module-loading cwd.
- **Broker/propose-verify-commit.** Direct REST commit and store commit are explicitly outside the
  projection. If the landed broker lacks a single session-aware dispatch entrypoint, P6.B stops;
  it does not synthesize one in the adapter.
- **Performance baseline.** No existing benchmark is removed. P6.B adds bounded protocol
  round-trip benchmarks only if measurement shows a material new hot path; this ~1d item does not
  rewrite `bench/BASELINE.md`.

## Systemic-Issue Audit

| Question | Finding | Response |
|---|---|---|
| Is the mismatch local to one endpoint? | No; discovery, A2A card, authentication, and invocation share it | Require one upstream callback seam and one World predicate |
| Can configuration solve it? | No; `--caps` and API key are process-wide | Do not encode sessions as flags |
| Can a sidecar solve it safely? | No; static exports, built-in feedback, and lifecycle remain | Reject path (b) |
| Is a new protocol required? | No | Reuse upstream codecs/handlers |
| Does this reveal a missing landed subsystem? | Yes; only the epoch registry is shipped | Make transition registry/broker session API an explicit prerequisite |
| Could a gate pass vacuously? | Yes; empty registry, one identical session, list-only tests, or skipped sockets | Require non-empty unequal session sets, calls, denials, and zero skips |

## Deferred Scope

- Designing or landing the clause-3 capability law, broker, session store, or transition registry.
- Any A2UI/AG-UI work (DESIGN.md M6 Generated UI).
- Remote/public bind, TLS termination, OAuth, multi-tenant federation, or cloud deployment.
- Changes to the REST v1 contract or shared error envelope.
- Upstream implementation details beyond the minimum public serving seam.
- Hot reload of transition definitions; registry snapshots already provide request-to-request
  dynamism.
- Benchmark thresholds for the clause-4 non-inferiority experiment.

## Acceptance Criteria

- [ ] **AC1 — upstream reuse:** production World code contains no MCP/A2A/SSE/JSON-RPC codec;
  requests and responses pass through the pinned upstream serving package.
- [ ] **AC2 — exact surface:** for two non-empty sessions with unequal capability sets,
  `tools/list` equals each authorized transition-ID set exactly; no extras.
- [ ] **AC3 — A2A agreement:** each session's A2A skill-ID set exactly equals its MCP tool-name
  set and changes when the session capability snapshot changes.
- [ ] **AC4 — ambient exports absent:** `exit`, `writeBytes`, all eight observed `std/io`
  exports, and `submit_feedback` are absent from both surfaces.
- [ ] **AC5 — invocation enforcement:** every listed tool dispatches through propose → verify →
  commit; an unauthorized, stale, or guessed tool name is rejected before broker dispatch and
  produces no store/log change.
- [ ] **AC6 — session failure:** missing, unknown, and expired sessions fail closed for list,
  card, and call; no default/global capability set is used.
- [ ] **AC7 — one-snapshot consistency:** a request never mixes registry/capability epochs;
  concurrent registry/session change is observed only on a subsequent request.
- [ ] **AC8 — protocol conformance:** MCP HTTP response handling asserts upstream SSE framing
  (`event: message` plus `data:` JSON), and A2A uses upstream's card/task handlers.
- [ ] **AC9 — landed behavior preserved:** REST v1 route/body regressions and cross-process
  writer-lock tests remain green; projection never opens a store.
- [ ] **AC10 — honest gates:** `verify_ail.sh` remains exactly 4 identities / 9 modules / 14
  tests; `verify_go.sh`, benchmark smoke, and both CI jobs pass with zero skipped protocol tests.
- [ ] **AC11 — dependency floor:** the printed transitive graph contains no disallowed package
  and `TestDaemonDependencyAllowlist` covers the new package/cmd paths.
- [ ] **AC12 — dynamic source:** changing the transition-registry head changes the next
  authorized list/card without restart or `.ail` file edits.
- [ ] **AC13 — bounded projection waits:** with the configured finite server bound, a
  barrier-blocked session resolver, registry snapshot, authorization provider, broker,
  proposer, verifier, or upstream response write, and a disconnected client each
  terminate within that bound. The same request context reaches every dependency; a writable
  transport receives the upstream structured timeout/cancellation error and a disconnected
  transport closes; no retry and no REST/store fallback occurs.
  **Cancellation is tested on BOTH SIDES of the commit boundary** (round-2 fix, `gpt5-6-sol`):
  cancelling immediately BEFORE the coordinator accepts the commit asserts **no durable
  mutation**; cancelling immediately AFTER acceptance asserts **exactly one recoverable,
  queryable/replayable receipt** under the stable invocation/idempotency ID — and in neither case
  does the caller receive a definitive "not committed" answer while the outcome is unknown.
  **The fault-injection test must also assert OS-LEVEL socket closure** (round-2 fix,
  `gemini-3-1-pro`): on stream-deadline expiry the underlying TCP connection is actively closed —
  observed via `http.Server.ConnState` tracking or a client read error — proving no silent
  connection leak once the D7 write deadline is relaxed. A logical Go `context` cancellation alone
  does not satisfy this criterion.
- [ ] **AC14 — SSE/REST deadline separation:** an `/mcp/` SSE connection can remain active past
  the frozen 30s global write deadline but closes at the earlier client/request/finite-stream
  bound. Tests and source assertions prove that only the streaming handler calls
  `ResponseController.SetWriteDeadline`/`SetReadDeadline`, and that REST `/v1/*` still uses the
  unchanged D7 values at `daemon.go:77-91` and `daemon.go:409-414`, including its 30s read/write
  deadlines and 120s idle timeout.

## Non-Vacuity — Named RED Mutation for Every Gate

| Gate | Named RED mutation (concrete edit) | Required red observation |
|---|---|---|
| AC1 `MUT-PROTO-OWNER` | Replace upstream handler with a local hand-written `tools/list` JSON response | source/ownership test rejects local protocol codec |
| AC2 `MUT-SESSION-UNION` | Return the union of all transitions for every session | low-capability session exact-set test REDs |
| AC3 `MUT-CARD-GLOBAL` | Generate the A2A card from the unfiltered registry | MCP/card set-equality test REDs |
| AC4 `MUT-UNFILTERED-PROJECTION` | Re-enable raw v0.30.0 projection | surface test REDs on `std.io.writeBytes`/`writeBytes`, `exit`, or `submit_feedback` |
| AC5 `MUT-CALL-BYPASS` | Dispatch by tool name without re-running authorization | guessed-name call changes dispatch counter/store and denial test REDs |
| AC6 `MUT-DEFAULT-CAPS` | Map an unknown session to a process default | unknown-session list/call test REDs |
| AC7 `MUT-SPLIT-SNAPSHOT` | Re-read capabilities after reading registry within one request | barrier-controlled epoch-consistency test REDs |
| AC8 `MUT-PLAIN-JSON` | Strip SSE framing and return only JSON | framing assertion REDs on missing `event:`/`data:` |
| AC9 `MUT-SECOND-OPEN` | Make projection call `store.Open` during mount | writer-lock/live-daemon test REDs with `WriterAlreadyActive` |
| AC10 `MUT-SKIP-SOCKET` | add `t.Skip` on listen failure | zero-skip CI assertion/source check REDs |
| AC11 `MUT-CLOUD-DEP` | add `cloud.google.com/go/storage` to projection imports | `TestDaemonDependencyAllowlist` reports it by name |
| AC12 `MUT-STARTUP-CACHE` | cache transition descriptors once at daemon startup | registry-head-change test returns stale set and REDs |
| AC13 `MUT-DROP-DEADLINE` | replace the propagated request context with `context.Background()` before one injected blocking dependency | bounded-wait fault-injection test exceeds the configured bound and REDs; mutation cleanup is terminated by the test harness |
| AC13 `MUT-COMMIT-BOUNDARY-LIE` | move the cancellation check from *before* coordinator acceptance to *after* it, so a cancelled-mid-commit request reports a definitive "not committed" | the after-boundary test REDs: a durable receipt exists under the invocation/idempotency ID while the caller was told the commit did not happen |
| AC13 `MUT-LEAK-SSE-CONN` | on stream-deadline expiry, cancel the Go context but never close the `http` connection | the `ConnState`/client-read-error assertion REDs — proving the test observes OS-level socket closure, not just logical cancellation |
| AC14 `MUT-SSE-REST-DEADLINE` | apply the SSE `ResponseController` deadline relaxation to the parent mux/REST paths instead of only the `/mcp/` streaming handler | route-scope/source assertion and REST deadline regression RED while the frozen D7 constants remain expected |

## Axiom Compliance

| Axiom | Score | Justification |
|---|---:|---|
| A1 Determinism | +1 | stable-ID ordering and one snapshot per request |
| A2 Replayability | +1 | calls use the normal recorded transaction path |
| A3 Effect Legibility | +1 | protocol I/O stays at the host boundary |
| A4 Explicit Authority | +2 | exact per-session list/card/call predicate |
| A5 Bounded Verification | +1 | existing verify pipeline; no protocol bypass |
| A6 Safe Concurrency | +1 | daemon retains sole handle; immutable request snapshot |
| A7 Machines First | +2 | MCP/A2A schemas and structured upstream errors |
| A8 Minimal Syntax | 0 | no language syntax change |
| A9 Cost Visibility | 0 | no new budget claim; benchmark scope deferred honestly |
| A10 Composability | +2 | standard MCP/A2A, upstream-owned |
| A11 Structured Failure | +1 | missing/expired/denied sessions fail in protocol form |
| A12 System Boundary | +2 | projection is an adapter, never kernel state |

**Net: +14; hard axioms A1/A3/A4/A7 are non-negative.**

## Premise Verification Log (live evidence, 2026-07-28)

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
| **Commit-boundary contract** | **VERIFIED — `w-store-durability` SD.B/SD.C, landed** | `store.AppendIntent` durably records the stable invocation ID; `Commit.InvocationID` binds the request and writes its outcome in the same transaction; `GetReceipt` exposes the three-state durable receipt. SD.C's real-process crash tests kill at the named pre-outcome boundary and prove world + receipt atomicity, including the split-transaction RED mutation. |
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

### UF-P6-1 — v0.30.0 cannot expose an exact caller-supplied/session-scoped MCP+A2A surface

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

## Open Decisions

1. **Exact upstream module/version.** **RECOMMENDED DEFAULT:** pin the first tagged release that
   exports the Decision-2 seam; do not use a pseudo-version in the World sprint. P6.B waits.
2. **Session credential carrier.** **RECOMMENDED DEFAULT:** reuse the landed clause-3 session
   resolver's HTTP header convention. If clause 3 lands without one, choose one header in the
   clause-3 API before P6.B; do not reuse the static serve-api key.
3. **Transition descriptor schema.** **RECOMMENDED DEFAULT:** use the landed transition
   registry's stable ID, input schema, description, and required capability verbatim. Do not
   create a projection-owned registry schema.
4. **If upstream export pulls a disallowed graph.** **RECOMMENDED DEFAULT:** ask upstream for a
   small protocol-only module/package. Do not broadly relax `TestDaemonDependencyAllowlist`.

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

## Related Documents

- [world-mission.md](../world-mission.md) — clause 6, reuse paths, guardrails, premise log
- [coding-standards.md](../coding-standards.md) — S1–S6
- [DESIGN.md](../DESIGN.md) — §3.7, §14, §17
- [w-worldd-m2.md](../implemented/w-worldd-m2.md) — shipped daemon and REST freeze
- [w-effect-broker-m3.md](w-effect-broker-m3.md) — prerequisite session/capability/broker design
- [worlddapi.ail](../sketches/worlddapi.ail) — frozen REST-law precedent; unchanged
