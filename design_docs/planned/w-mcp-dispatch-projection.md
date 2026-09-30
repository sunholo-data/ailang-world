# w-mcp-dispatch-projection — Session-Scoped MCP Projection over the Delivered `mcphttp` Seam

**Status**: **PARKED — needs-human-review after two blocked quorum rounds (2026-09-30); no sprint authorized**. The blocking predicate the placeholder carried — upstream
[`sunholo-data/ailang#885`](https://github.com/sunholo-data/ailang/issues/885) — is **DELIVERED**:
tag `v0.47.2` ships `serveapi/protocol/mcphttp` (handler.go, methods.go, wire.go) plus
`serveapi/protocol/hostcall`, the issue is CLOSED COMPLETED, and the proxy/sumdb fetchability
caveat has cleared (200/200, controller-measured — V1/V2). This revision is therefore the **full
design against the actual delivered seam**, per the placeholder's own unblock condition: every
carried obligation is adjudicated against what `mcphttp` really does (it is POST-only and
stateless — several SSE-stream obligations REDUCE, explicitly, below), every codebase claim is
logged with its command and observed output, and every acceptance criterion names a RED mutation.
**NOT quorum-cleared**: the controller runs the design quorum at pick time; nothing here
pre-authorizes a sprint.

**Item**: charter queue **row 108**; split child #1 of `w-mcp-projection` (charter clause 6 — the
*"project the transition registry over MCP"* half)
**Clause**: clause-6 — MCP projection of the transition registry, capability-filtered per session
**Estimate**: **~0.9 World days** (M108-1 ~0.1 + M108-2 ~0.15 + M108-3 ~0.3 + M108-4 ~0.35), against
the charter row's `~1d` ceiling; no dependency waits remain (the coordinator landed as row 106, the
session authority as row 39, the toolchain floor as P6.T, the protocol pin as P6.D). Revision 2 adds
the whole-POST budget (D108-10), the resolver adapter's inner bound, and the production-constants
test leg — a ~0.05d bump in each of M108-3/M108-4, inside the ceiling.
**Author**: rotation designer, glm-5.3, iteration 209
**Revision**: **r2 — the one post-quorum revision** (quorum round 1, 2026-09-30: verdict BLOCKED,
2 REJECTs). Both objections are answered by measurement, not argument: (1) *gemini-3-1-pro* — the
`procbound.MaxOutstanding` premise was unverified → **V32** now reads `host/procbound/procbound.go`
first-party and **D108-5 is rewritten** to separate the Runner's slot cap from procbound's
process-wide reservation limit (they are different mechanisms; the value 8 is an alignment choice,
not an identity); (2) *oc-kimi-k3* — the "whole POST answers inside the frozen 30 s `writeTimeout`
by construction" claim was FALSE (three sequential `hostcall.Run`s, each with its own fresh 20 s
timeout, sum to 33 s of inner-bounded callbacks or 60 s without inner bounds, plus body read) →
**D108-10** replaces the claim with a derived whole-POST budget: one aggregate request context
(`invokeDeadline`, the landed A2A route's own whole-route bound, V35), the resolver adapter gains
the frozen `credentialBudget` inner bound (V35), and the honest transport accounting (write clock
starts at header-read; body read eats the window, V34) is recorded with the one client-induced
truncation leg (R108-6). A production-constants test leg + `MUT-AGGREGATE-REMOVED` now pin the
arithmetic (M108-4).
**Verified against**: upstream `github.com/sunholo-data/ailang` tag **`v0.47.2`** read first-party
from the module cache at
`/Users/voightkampff/go/pkg/mod/github.com/sunholo-data/ailang@v0.47.2` (the tag the controller
verified delivered: issue ailang#885 CLOSED COMPLETED — charter row 108, V1) and World `dev` at
**`8a9f961`** (this worktree's HEAD; all [R] rows below were measured there). **Tag commit ≠
delivering-PR squash — distinguished deliberately**: tag `v0.47.2` points at commit `e939cba03`
(first-party in the module cache: `@v/v0.47.2.info` records
`"Hash":"e939cba032c0f38fbecffb47e8261fb58e1a2ca1"`, `"Ref":"refs/tags/v0.47.2"`), while the #885
work merged as PR #1369's SQUASH commit `aa3063c37`; the charter's queue-status record names both
(`design_docs/world-mission.md:1720`: "`e939cba03`; PR #1369 squash `aa3063c37`") — neither SHA is
citable as the other.
**Date**: 2026-09-29

---

## Why this doc exists — the split history and the objection it now discharges

The parent, [`w-mcp-projection.md`](w-mcp-projection.md), was blocked at quorum round 3
(2026-08-25) by a DIRECTION-level objection from `gpt5-6-sol`, confirmed first-party by the
controller. The disposition was SPLIT (iteration 125): the reviewer-clean A2A half moved to
[`w-a2a-session-projection.md`](../implemented/w-a2a-session-projection.md) (landed, iteration
187), the enablers landed with it (P6.T toolchain floor; P6.D the pinned
`github.com/sunholo-data/ailang/serveapi/protocol` package-path admission), row 106 landed the
invocation coordinator (HEAD `8a9f961`), and this doc carried the objection and the MCP half,
honestly blocked on `#885`. The objection, verbatim:

> "The selected upstream seam is insufficient for the design's own no-codec rule … no MCP request
> parser or handler for JSON-RPC method dispatch, initialization, `tools/list`, or `tools/call`. A
> World-authored `/mcp/` handler therefore appears forced to implement MCP/JSON-RPC parsing and
> dispatch locally, directly contradicting P1, the Design Freeze, and AC1."
>
> — `gpt5-6-sol`, quorum round 3 on `w-mcp-projection`, 2026-08-25 (artifact
> `.ailang/state/mission-quorum/w-mcp-projection-2026-08-25T16-33-38Z.json`)

`#885` asked for exactly what the objection demanded — an SDK-free MCP JSON-RPC dispatch seam — and
it DELIVERED in `v0.47.2`: `mcphttp.NewHandler(mcphttp.Config)` returns a complete, wire-owning
`http.Handler` (initialize, ping, tools/list, tools/call, notifications, batching, transport
precondition checks, SSE response framing, the frozen host-failure envelope) over four World-side
callback seams. The objection's forced-local-codec route is therefore closed by delivery, not by
argument: **World writes no JSON-RPC, no SSE bytes, and no envelope — it implements three adapter
interfaces and a name mapping.** This revision replaces the placeholder in place; the placeholder's
obligations are adjudicated one-by-one below (section "Obligations carried from the split").

---

## The delivered seam — what `mcphttp` actually is

Read first-party from the module cache at tag `v0.47.2` (file:line anchors below). Package doc
comment, `mcphttp/handler.go:1-13`:

> "serves the tools subset of MCP over stateless Streamable HTTP using only the standard library
> and serveapi/protocol. It exists so a consumer whose build must stay free of the MCP SDK (and
> its OAuth stack) can still project host tools over MCP without writing its own JSON-RPC codec
> (ailang#885). … Scope: initialize, ping, tools/list, tools/call and notifications, on POST only.
> There are no sessions, no GET stream and no server-initiated messages."

**The exported construction surface** (Config at handler.go:37-44, NewHandler at :48-63):

```go
func NewHandler(config Config) (http.Handler, error)
type Config struct {
    Agent    protocol.AgentInfo      // Name and Version required (TrimSpace-non-empty)
    Resolver protocol.SessionResolver // ResolveSession(ctx, *http.Request) (Session, error)
    Tools    protocol.ToolSource      // Tools(ctx, Session) ([]ToolDescriptor, error)
    Invoker  protocol.Invoker         // Invoke(ctx, Session, Invocation) (InvocationResult, error)
    Runner   *hostcall.Runner         // required; "There is no unbounded default."
}
```

`NewHandler` validates every field and fails closed at construction (handler.go:48-63; the
validation switch at :50-61). The
callback contracts are `protocol.Session = any` plus the three one-method interfaces
(`protocol/interfaces.go:6-21`); `hostcall.New(timeout, maxConcurrent)` requires both positive
(`hostcall/runner.go:23-31`).

**The request order of operations** (handler.go:72-107 ServeHTTP, 109-131 authorize): (1) non-POST → 405
(`writeMethodNotAllowed`, `Allow: POST`, wire.go:97-102); (2) body read capped at
`MaxRequestBodyBytes = 4 MiB` (wire.go:19-21) — over-cap answers the frozen envelope
(`WriteMCPEnvelope`, 200/-32603/"invalid MCP request body"); (3) **authorize**: `ResolveSession`
then `Tools` then `protocol.CallerSurface(descriptors)`, EACH through `hostcall.Run` — i.e. a
denied session learns nothing about the transport, and an admitted POST has already built the
authorized surface **before** any Content-Type/Accept/version check; (4) transport preconditions:
`Content-Type: application/json` and an `Accept` listing BOTH `application/json` and
`text/event-stream`, else HTTP 400 JSON bodies (wire.go:71-77, handler.go:90-105); (5) the
`MCP-Protocol-Version` header: empty → default `2025-03-26`; unsupported → 400 listing the
supported set; `SupportedVersions = ["2025-11-25","2025-06-18","2025-03-26"]` (handler.go:27-30;
`defaultVersion = "2025-03-26"` at :32-34).

**Dispatch** (methods.go:28-43; `initialize` at :45-61): `initialize` (negotiates the client's version if supported, else
`SupportedVersions[0]`; capabilities `{"tools":{"listChanged":false}}`; stateless), `ping`,
`tools/list` (`{"tools":[…]}` from the authorized surface, **sorted by tool NAME** by
`CallerSurface`, descriptor.go:39), `tools/call` (`surface.Lookup` — unknown name → JSON-RPC
`-32602 "unknown tool"` at 200; then `Invoker.Invoke` through the Runner), notifications/responses
accepted and dropped (stateless), and 2025-03-26 batch handling. **A failed host callback answers
the WHOLE POST with the frozen envelope, never a partial stream** (handler.go:182, citing the
seam's own `TestEmbeddedMCPFrozenCallbackEnvelopes`).

**The response wire** (wire.go:79-90): success is ONE SSE event — `Content-Type:
text/event-stream`, `Cache-Control: no-cache`, `X-Content-Type-Options: nosniff`, HTTP 200, body
`event: message\ndata: {json}\n\n` — and the handler RETURNS (the stream ends; stateless). An
all-notification POST answers 202 `application/json` with no body (wire.go:92-95). Transport
errors are HTTP-status JSON bodies. Host-callback failures are `protocol.WriteMCPEnvelope`:
HTTP 200, `-32603`, `CallbackMessage(err)` — "host callback timed out/canceled/capacity
exceeded/failed" (`protocol/envelope.go:21-50`). Authorization failures whose error exposes
`HTTPStatus()` of exactly 401/403 are rendered by the seam as plain `http.Error(w, err, status)`
(handler.go:114-115 + `protocol.AuthorizationStatus`, envelope.go:52-61 — 400 is NOT honored and
falls to the envelope).

**Descriptor validation the seam imposes** (`protocol/descriptor.go:30-66, 143-147`):
`ValidateMCPName` = `^[a-zA-Z0-9_-]{1,64}$` on every tool name; a non-nil `InputSchema` that
unmarshals to a JSON object whose `"type"` is `"object"`; duplicate names refused; the whole
surface fails closed on ANY violation (one bad descriptor refuses everything — the upstream
precedent this design aligns its own refusal rules with).

**Upstream tests the seam itself** (test files in the package): `TestWireContract`,
`TestInvalidJSONIsAParseError`, `TestNewHandlerRefusesMissingFields`, `TestOversizedBodyIsEnvelope`,
plus a differential `TestParityWithSDK`/`TestSDKClientInterop` where the MCP SDK is imported ONLY
in `_test` files, outside the build closure. World's conformance tests pin the OBSERVABLE from
World's mount; they do not duplicate upstream's internals.

**The seam's dependency guarantee, re-measured first-party**: `go list -deps
./serveapi/protocol/mcphttp` inside the v0.47.2 module tree yields exactly **three non-stdlib
packages — `serveapi/protocol`, `serveapi/protocol/hostcall`, `serveapi/protocol/mcphttp`**; the
hostcall closure yields exactly two (`protocol`, `hostcall`) (V8). No MCP SDK, no OAuth stack, no
network egress.

---

## Decisions

### D108-1 — The go.mod pin: **`v0.47.2`**, not `v0.48.0`

Recommend pinning `github.com/sunholo-data/ailang` at **`v0.47.2`**. Reasons: (a) `v0.47.2` is the
tag `#885` delivered in and the exact tree this doc read first-party from the module cache (V3,
V8); (b) both tags are fetchable — proxy.golang.org AND sum.golang.org answered 200 for both
`v0.47.2` and `v0.48.0` (controller-measured iteration 209, V2; the charter's 2026-09-28 404
caveat cleared); (c) the serveapi/protocol tree is **byte-identical at both pins** — the
34-commit compare delta between v0.47.2 and v0.48.0 touches zero serveapi/protocol files
(controller-measured, V3) — and my first-party shasum extends the A2A doc's F4 chain: all five
protocol files are byte-IDENTICAL at `v0.33.2` and `v0.47.2` (V4), so the bump moves the A2A wire
code by zero bytes. **Trade-off, named**: `v0.48.0` is newer but buys World nothing (only
protocol+mcphttp+hostcall enter the daemon graph, byte-identical in both) while widening the
review surface to 34 commits this doc never read; pinning `v0.47.2` does not block a later
fleet-wide bump (the pin moves once per admission, as it did for v0.33.2). The charter's pre-land
rule stays: before the bump commits, `curl` proxy and sumdb and require 200 (V2, U3).

### D108-2 — The transition-ID → tool-name mapping `M` (collision-free, total, reversible)

World stable IDs admit `[a-z0-9_-./]`, 1..128 bytes, dot/slash-separated segments of 1..32
(`host/transitionreg/transitionreg.go:320-345`, V16); `ValidateMCPName` admits
`^[a-zA-Z0-9_-]{1,64}$` (V16). `M` encodes with `_` as the escape character:

| source char | encoding |
|---|---|
| `_` | `_u` (escape the escape) |
| `.` | `_d` |
| `/` | `_s` |
| `[a-z0-9-]` | itself |

So `tools.ghost → tools_dghost`, `world/recovery-transition/v1 → world_srecovery-transition_sv1`.
**Total** on the ID grammar; **injective** (each source char has a unique image and no image is a
prefix-ambiguity of another: an encoded `_` is always followed by exactly `u`/`d`/`s`); therefore
**reversible**: decode is a left-to-right scan — `_u→_`, `_d→.`, `_s→/`, any other `_x` or any
uppercase char is an INVALID encoding and decode fails LOUDLY (fail closed). `decode∘encode =
identity` is proven by a generated corpus unit test (MUT-MAP-COLLIDE/MUT-MAP-DECODE-INVALID
killers). **Collision refusal**: injectivity makes collisions impossible by construction; the
adapter nevertheless refuses the whole surface (fail-closed, matching the seam's own
whole-surface refusal precedent for duplicate names, V17) if two allowed IDs ever encode
identically or if any encoded name exceeds 64 bytes (possible: a 128-byte ID can encode past 64 —
refused, logged by descriptor, never truncated). The refusal surfaces as a `Tools` error, which
the seam renders as the frozen envelope — never a silently lossy projection. **Schema
normalization** rides the same adapter: the registry carries `InputSchema: []byte("{}")` today
(V17 — every fixture), which `CallerSurface` REJECTS (`"type" != "object"`, V17); the adapter
passes a schema that is already a JSON object of `type:"object"` through verbatim and
substitutes the minimal permissive `{"type":"object"}` otherwise — honest because schemas are
carried and never validated anywhere (the landed card's own description says so — daemon.go:588,
V19), and the seam validates only the shape. `OutputSchema` passes verbatim (CallerSurface only requires it be
valid JSON).

### D108-3 — Where the handler mounts: adapters in `host/projection`, route in `host/daemon`

Both candidate homes were read (V18-V20). **The three World adapters — resolver (`authority.Resolver`
→ `protocol.SessionResolver`), tools (allowed set → `[]protocol.ToolDescriptor` through `M`), and
invoker (row-106 coordinator → `protocol.Invoker`, the R-106-6 residual the coordinator doc
assigned to this row, V22) — live in `host/projection`**, in a new file `mcp.go` (+ `mcpname.go`
for the mapping), and `host/projection` is the ONLY importer of `mcphttp` and `hostcall`. The
daemon changes are the mirror of its A2A mount: `projection.Config` gains scalar bounds — r2 grows
the set from two to six, ALL values the daemon already owns by name (`CallbackTimeout` =
`invokeDeadline`, `MaxCallbacks` = `procbound.MaxOutstanding`, `CredentialBudget`, `ReadDeadline`,
`InvokeDeadline`, `WriteWindow` = `writeTimeout`; NO new daemon constant is minted) — and
`projection.New`
constructs `hostcall.New(...)`, `mcphttp.NewHandler(...)`, and the D108-10 aggregate wrapper, and
exposes the WRAPPED result, and
`daemon.Handler()` registers **one additive route: `mux.HandleFunc("POST /mcp/", …)`**
(daemon.go:711-730 shape, V19). Why here and not `host/daemon`: putting the adapters in the
daemon core grows the gated core package with projection-shaped code the A2A precedent already
moved out; putting them in a THIRD new package duplicates the landed admission orchestration
(`bindingCaps`, the B3 absent-head pre-check, `allowedDescriptors`) that `host/projection` already
owns unexported — and the cross-surface criterion (AC-CROSS-108) is testable in-package only if both
surfaces share it. Why not a separate `host/mcpprojection` sibling: same duplication argument; the
package doc comment updates from "the session-scoped A2A projection surface" to the session-scoped
A2A **and** MCP projection surfaces — an additive charter change to one comment, judged at
quorum. `d.isProtected` is untouched (POST /v1/commit only, V19): `/mcp/` resolves the session
itself, exactly like `/a2a/` (landed AC9 mutation discipline mirrored as MUT-MCP-PROTECTED).

### D108-4 — Session carrier: D-WORLD-26 ARM A, rendered through the seam's 401 path

The resolver adapter calls the ONE landed `authority.Resolver.ResolveContext(ctx,
r.Header.Get("Authorization"), now)` — the standard `Authorization: Bearer <session-credential>`
header ONLY, fail closed on absent/malformed/unknown/expired, never an API key, no alternate
header even as fallback (D-WORLD-26 = ARM A constraints, landed and tested on `/a2a/`, V20). A
denial maps to `*protocol.AuthorizationError{Status: http.StatusUnauthorized, Err: …}` with a
constant per-kind message (the A2A vocabulary's wording), which the seam renders as
`http.Error(w, msg, 401)` — the ONLY HTTP-status rendering the seam honors (401/403, V14).
Malformed reads 401 here rather than A2A's 400 because the seam does not honor 400
(`AuthorizationStatus` returns 0 for it → the generic envelope); recorded, not hidden. A
resolution store failure (not a denial) returns the raw error → the seam's frozen envelope — a
failed lookup is not an authentication result, matching the landed resolver contract (V20). The
session value is the `*authority.SessionBinding` itself (`protocol.Session = any`); both adapters
type-assert it and fail closed on mismatch.

### D108-5 — The `hostcall.Runner` bound: per-call timeout `invokeDeadline`; slot cap `8` — TWO mechanisms, honestly separated

The Runner bounds EVERY host call on `/mcp/` — resolution, listing, and invocation alike — with
one per-call timeout and one slot cap (V21): **timeout = the frozen `invokeDeadline` (20 s,
daemon.go:98)**, the same bound `/a2a/` invokes under (row-106's deadline source, V22). **The slot
cap is `8`, and its justification is now corrected and measured (V32, the r2 fix for the gemini
objection): the Runner's `maxConcurrent` and `procbound.MaxOutstanding` are DIFFERENT mechanisms
and neither enforces the other.** `procbound.MaxOutstanding = 8` (procbound.go:31) is the
process-wide limit on **subprocess reservations** — "admitted children not yet reaped, including
those whose reap a background waiter owns" (:29-30); a reservation is held from `Admit()` before
`cmd.Start` (broker/handlers.go:121) until the child is REAPED, which the landed broker can do
AFTER the calling handler has returned, via `procbound.Wait`'s background waiter
(broker/handlers.go:132). The Runner's slot is held only while the in-process callback is running
(runner.go:15-16). So a Runner slot is not a subprocess reservation: eight free Runner slots do
NOT imply procbound headroom (retained waiters may hold all eight reservations), and eight held
Runner slots do not imply eight children. **The value 8 is an alignment CHOICE, not an identity**:
the mount's steady-state concurrent invokes (each of which becomes one broker subprocess holding
one reservation) are sized to the process's unreaped-children headroom, so the Runner does not
admit callbacks whose children procbound would refuse anyway; and the two overflow modes are
distinct and both fail closed with distinguishable frozen envelopes — Runner: `ErrCallbackCapacity`
→ "host callback capacity exceeded" (runner.go:47-51, V21); procbound: `ErrCleanupBacklog` from
`Admit()` → the broker error → the callback returns an error → "host callback failed" (V32).
Two honest consequences recorded: (i) the A2A route invokes through the SAME broker with NO Runner
at all (V18), so the Runner cap bounds only the MCP surface and `procbound` remains the only
process-wide subprocess bound — the design does not claim the Runner caps children; (ii) a
larger Runner would queue callbacks the process cannot run; a smaller one queues for nothing.
Per-request cost: one slot at a time (the seam runs its callbacks sequentially, V33). Beyond the
cap, callers wait on the slot queue until the per-call timeout or the aggregate (D108-10) fires,
then fail closed with `ErrCallbackCapacity` → the seam's frozen "host callback capacity
exceeded" envelope (V21). **Two guarantees, separated honestly (V21, the
Runner's own doc comments)**: the RESPONSE deadline holds ALWAYS — `hostcall.Run`'s select returns
`callCtx.Err()` when the timeout fires (runner.go:61-66), so the POST is answered with the frozen
"host callback timed out" envelope within the EFFECTIVE bound — the Runner timeout clipped by the
D108-10 aggregate's remainder — no matter what the callback does. The SLOT is
different: the upstream Runner cannot forcibly terminate an in-process callback ("In-process
callbacks cannot be forcibly terminated; a non-cooperative callback keeps its slot",
runner.go:38-39; "A slot remains occupied until host code actually returns, even after timeout",
:15-16) — the slot frees ONLY after the callback cooperates with the deadline ctx and RETURNS.
This design's slot-free claim is therefore conditional, and the condition is the landed row-106
ctx cooperation: the invoker dispatches under the callback ctx the seam passed through `Run`, the
coordinator runs the capsule under that same ctx (`RunContext(ctx, …)`, coordinator.go:271) and
returns when it is done (`ctx.Err()` checked, coordinator.go:275), so a blocked invocation
returns at the deadline, its slot frees, and the next request is admitted. Severing the
cooperation — dispatching under a ctx that ignores the deadline — is exactly what MUT-CB-BACKGROUND
mutates, and the deadline test's subsequent-call leg is what reds (the POST itself still answers
within the bound — now the D108-10 aggregate, not the Runner timeout alone). The daemon supplies
the scalars — `invokeDeadline` (Runner timeout), `procbound.MaxOutstanding` (slot-cap value),
`credentialBudget`, `readDeadline`, and `writeTimeout` (for D108-10's derived validation), all
frozen constants injected through `projection.Config`; NO new daemon constant is minted;
`projection.New` validates every one positive and constructs the Runner — a zero/negative bound
is a construction error, never "unbounded" (MUT-RUNNER-UNBOUNDED) — and re-checks the landed
invariant `invokeDeadline < writeTimeout` (daemon.go:555-556, V35) at the projection mount so the
aggregate in D108-10 cannot inherit a bad pair even if the daemon's own check is ever refactored.

**Inner bounds align EVERY callback with the landed D7 discipline — all three adapters, none
silent (the r2 fix for the kimi catch that the resolver had none)**: the resolver adapter applies
`credentialBudget` (3 s, daemon.go:153) — the SAME bound the landed session middleware applies to
every credential lookup (`context.WithTimeout(r.Context(), m.budget)` around `ResolveContext`,
middleware.go:55, wired with `credentialBudget` at daemon.go:730; row-23 bound table B4, V35) — so
a resolution store read no longer rides the bare 20 s Runner timeout; the tools adapter applies
`readDeadline` (10 s, daemon.go:142) — the same bound the A2A card route uses — and the invoker
applies `invokeDeadline` (20 s) inside the Runner's own 20 s. On inner-budget expiry the resolver
returns the raw deadline error, NOT a typed denial — a lookup that did not finish is not an
authentication result (the B4 rule, middleware.go:59-62) — and the seam renders the frozen
"host callback timed out" envelope; the 401s remain reserved for typed denials (D108-4). The
three inner bounds do not compose into a whole-POST bound by themselves (3 + 10 + 20 = 33 s > the
30 s window, V36): that composition is D108-10's aggregate, which is REQUIRED, not optional.
**The Runner is load-bearing for the row-23 store policy**: the seam calls callbacks under
`r.Context()`, which carries NO deadline, and every `*Store` read rejects a deadline-free call
with `ErrNoDeadline` (store.go:236-244, V19); the Runner's `WithTimeout` on every call
(runner.go:42) is what makes every store read under `/mcp/` deadline-carrying, and the inner
bounds + D108-10's aggregate tighten it per leg.

### D108-6 — The route table (additive) and GET behavior (refused, 405)

One registration, `mux.HandleFunc("POST /mcp/", …)` — the same subtree-pattern shape the landed
A2A mount uses at daemon.go:729 — added AFTER the frozen `/v1/*` block, which is byte-untouched
(the ten `/v1` patterns + `GET /workbench`, V19). **GET (and every non-POST) on `/mcp/` is refused
405 with `Allow: POST` by the ServeMux method-pattern mismatch** (Go's own 405 branch,
GOROOT `net/http/server.go:2699-2705`, V23); the seam's internal `writeMethodNotAllowed` is the
same observable and remains as defense-in-depth. This ADJUDICATES the GET-stream obligations: the
delivered seam has NO GET stream (package doc, V9), so there is nothing to relax or keep alive.

### D108-7 — Version negotiation: the seam owns it; World asserts the observable

World has no version policy of its own: `SupportedVersions`, the `2025-03-26` default, the
unsupported-version 400, and `initialize`'s negotiate-or-`SupportedVersions[0]` behavior are all
seam-owned (V10-V12) and pass through untouched. World's conformance test asserts the observable
from World's mount (initialize returns a supported version; an unsupported header → 400).

### D108-8 — The invocation path: the row-106 coordinator through the R-106-6 adapter

`tools/call` arrives as `{name, arguments}`; the invoker **decodes** the name with `M⁻¹`, builds
its OWN fresh admission snapshot (one `transitionreg.NewRequest` over the binding's immutable
caps — B2's rule at invoke time: possession of a listed name conveys nothing), refuses unless
the decoded ID is in that snapshot's `Allowed()` set, unmarshals `arguments` to a JSON object
(the coordinator's own calling convention), and dispatches `coordinator.Dispatch` with a
**server-minted task id** (32 random bytes hex — fits the coordinator's `validTaskID`, V22).
Consequences stated honestly: (a) the seam's callback shape forces **one snapshot per callback**
(D108-9); (b) **MCP `tools/call` has no idempotent-resend channel** — a client retry after a lost
response executes a NEW invocation (a fresh task id, a new journal row); deriving the task id
from the client's JSON-RPC request id was considered and REJECTED (client-controlled, may repeat
across distinct calls — silent reconciliation of two different calls is worse); recorded as
residual R108-2; (c) MCP has no metadata channel, so the A2A route's optional `transition_fn` pin
does not exist here (`PinnedFn` stays nil); (d) **the error taxonomy is the seam's**: any
coordinator refusal — including capability denials, invalid params, conflicts, and deadlines —
renders as the frozen whole-POST envelope (`-32603` + `CallbackMessage`), never per-code JSON-RPC
errors; the typed detail goes to the daemon's one-line error log through the injected log seam
(the A2A `dispatchError` matrix has no MCP analog, by the seam's contract, V13).

### D108-9 — Snapshots: one per callback, each fail-closed against its own epoch

The seam calls `Tools` (at authorize time) and `Invoke` (at dispatch time) as separate callbacks
with no surface hand-off, so the A2A route's one-snapshot-per-REQUEST property cannot hold
literally here: a POST can observe two registry epochs. The design takes the strictly-safer
reading: **each admission is derived fresh and checked against its own snapshot**, so a tool
unpublished between list and call fails closed at call time. The cost is bounded by the landed
`StoreReader` head-keyed snapshot cache (transitionreg.go:63-93, V19): the second `NewRequest`
in a POST is a head read + cache hit, not a re-parse. The A2A doc's AC7 (one-snapshot
consistency) cannot hold verbatim on this surface; the amended reading is: never "one epoch per
POST", always "each admission against its own snapshot, fail-closed" (AC-INVOKE-108,
MUT-INVOKE-NO-READMIT).

### D108-10 — The whole-POST budget: ONE aggregate request context at `invokeDeadline` — REQUIRED, and the honest transport accounting (r2, answers the kimi REJECT)

**The claim r1 carried — "the handler returns inside the frozen D7 `writeTimeout` (30 s) … by
construction" — was FALSE, and this decision replaces it with a derived bound.** The measured
facts (V33): a `tools/call` POST runs up to THREE sequential `hostcall.Run` calls — `ResolveSession`
then `Tools` in `authorize` (handler.go:110-126), then `Invoke` at dispatch (methods.go:83-85) —
and `Run` derives a FRESH `context.WithTimeout(ctx, runner.timeout)` per call (runner.go:42), so
with no outer deadline the POST's server-side work is bounded by the SUM of the three per-call
timeouts: 60 s at a bare 20 s Runner, still 33 s with every inner bound applied (3 + 10 + 20, V36)
— past the frozen 30 s `writeTimeout` at which the transport kills the write. The seam passes
`r.Context()` through unchanged and that context carries no deadline (its cancellation sources are
client disconnect and handler return, not any clock), so nothing upstream composes the legs.
**The fix is the landed A2A discipline transposed**: the A2A route already serves its ENTIRE POST
— resolve + body decode + admit + dispatch — inside ONE bounded context,
`context.WithTimeout(r.Context(), h.invokeWait)` with `invokeWait = invokeDeadline` (20 s,
projection.go:266, daemon.go:592, V35, which is why `/a2a/` always answers inside the window.
The MCP mount cannot put the ctx inside the handlers (the seam owns them), so it wraps the OUTSIDE:
`projection.New` returns the seam's `http.Handler` wrapped in a World-authored middleware whose
entire body is `ctx, cancel := context.WithTimeout(r.Context(), postBudget); defer cancel();
h.ServeHTTP(w, r.WithContext(ctx))` — the **aggregate**. `postBudget` derives to `invokeDeadline`
(the frozen 20 s), so the MCP surface gets exactly the A2A surface's whole-POST bound, and the
three callbacks' effective bounds become the MIN of their inner bound, the Runner timeout, and
the aggregate's remainder (`context.WithTimeout` keeps the EARLIER parent deadline, so the
aggregate can only tighten a leg, never loosen one — no bounded-wait is weakened anywhere).
When the aggregate fires, `Run`'s slot-wait short-circuit returns the parent's `ctx.Err()`
(runner.go:47-49) or its own select returns `callCtx.Err()` (:61-66) — both
`context.DeadlineExceeded` → `CallbackMessage` renders the SAME frozen "host callback timed out"
envelope (envelope.go:27-29), and every later leg of that POST sees an already-expired parent and
fails closed instantly. Fail-closed-early is strictly better than the unaggregated alternative:
without the aggregate the answer composes at 33 s, past the window — the client would see a
TRUNCATED CONNECTION instead of the frozen envelope (this is exactly what `MUT-AGGREGATE-REMOVED`
red-flags in M108-4's production-constants leg).

**The derived budget table** (production constants, `tools/call` POST; every cell cites its read —
V19/V32-V36):

| leg | bound | mechanism (all frozen) | clock |
|---|---|---|---|
| headers | ≤ 5 s | transport `readHeaderTimeout` (daemon.go:89) | from request start t₀ |
| body read (≤ 4 MiB) | client-paced, ≤ 30 s from t₀ | transport `ReadTimeout` (daemon.go:93); armed at t₀, applied to the body read (GOROOT server.go:982-985, 1041) | overlaps the write window |
| resolve | ≤ 3 s | resolver adapter inner `credentialBudget` (daemon.go:153), then min(Runner 20 s, aggregate) | inside the aggregate |
| tools | ≤ 10 s | tools adapter inner `readDeadline` (daemon.go:142) | inside the aggregate |
| surface build | CPU-only, no I/O | `CallerSurface` over the authorized slice | ~ms |
| invoke | ≤ 20 s | invoker adapter inner `invokeDeadline` (daemon.go:98) | inside the aggregate |
| **Σ inner bounds** | **33 s** | arithmetic (V36) — EXCEEDS the window; the inner bounds alone do not compose | — |
| **AGGREGATE** | **20 s** | the wrapper ctx = `invokeDeadline`; clips the Σ to 20 s | fires at handler entry + 20 s |
| response write | ≥ ~10 s margin | transport `WriteTimeout` 30 s (daemon.go:97) | the margin the landed invariant `invokeDeadline < writeTimeout` (daemon.go:555-556) already guarantees |

**The honest transport accounting (V34, measured in GOROOT go1.26.6)**: the write clock does NOT
start when the handler starts — `readRequest` arms `WriteTimeout` in a `defer`
(`c.rwc.SetWriteDeadline(time.Now().Add(d))`, server.go:986-989) that runs when `readRequest`
RETURNS, i.e. right after the request HEADERS are parsed and BEFORE the handler is called; and
the seam reads the body inside `ServeHTTP` (the `io.ReadAll` at handler.go:77), AFTER the wrapper has run. So the
wrapper's entry is ≈ the write-deadline start, and **body-read time is charged against the same
window the callbacks need** — stated plainly, as the controller's re-measure demanded. Three legs
follow, each with its own bound and none asserted:

1. **Body prompt (arrives within the aggregate), callbacks slow — the leg the aggregate fixes**:
   the POST is answered ≤ 20 s after
   handler entry with the frozen envelope or the invocation result, and the response write has
   ≥ ~10 s inside the frozen 30 s window. DELIVERED inside the window, always — this is the
   replacement for the deleted "by construction" sentence, and it is DERIVED (the table), not
   asserted, and MEASURED (the production-constants leg).
2. **Body slow (client-paced, slower than the aggregate)**: the body read eats the aggregate
   first, and a body arriving after the aggregate has fired composes its answer instantly with a
   SHRINKING write margin. A ctx deadline does not
   interrupt `r.Body` reads (they are bounded by the transport's read clock, not the request
   context), so the read proceeds until the bytes arrive or `ReadTimeout` kills it at t₀+30 s;
   when the read completes the first `hostcall.Run` sees an already-expired parent and the POST
   fails closed INSTANTLY with the frozen "host callback timed out" envelope. If the client
   consumed the whole window dripping its own body, that final write may truncate at the socket.
   This is the ONE leg on which the frozen envelope may be undelivered; it is client-induced,
   bounded by the transport clock, and the handler still returns, the callback still returns under
   the row-106 ctx cooperation, and the slot still frees — no zombie stream, no leak (R108-6
   records it; the landed A2A route has the IDENTICAL leg — its body decode rides the same 20 s ctx
   under the same 30 s transport clocks — so the MCP mount adds no new exposure).
3. **No leg hangs**: the answer is composed by min(handler-entry + 20 s, t₀ + 30 s) — the
   aggregate and the transport read clock are the two composition bounds, and the response write
   is bounded by the transport write clock (header-read end + 30 s) behind them; every wait on the
   surface has a named mechanism — per-callback inner bound, Runner timeout, aggregate, transport
   read clock, transport write clock. The bounded-waits axiom is DERIVED from this table, not
   asserted, and the production-constants leg is its empirical control.

**Why this shape and not the alternatives** (the controller's "most robust" mandate): raising the
inner bounds is impossible (the D7 block is frozen byte-untouched, and `invokeDeadline`/
`readDeadline`/`credentialBudget` are its constants); tightening the inner bounds to make 3+10+20
fit 30 s would mint NEW daemon numbers with no landed discipline behind them; and letting the sum
exceed the window (kimi's option (b), truncate-and-disclose) delivers the client a connection
reset on exactly the slow-authorize path the aggregate makes deliverable. The aggregate reuses the
ONE landed whole-POST mechanism (the A2A route's bounded ctx) and the ONE landed margin invariant
(`invokeDeadline < writeTimeout`, validated at New, daemon.go:555-556), so every number in the
derivation is already frozen, already validated, and already proven in production on `/a2a/`.
The wrapper is NOT a codec: it writes no status, no header, no byte (pinned by a source-level
assertion in the same scan family as AC-WIRE-108's no-envelope/no-relaxation legs) — it is
context plumbing only, so the no-local-codec non-negotiable is untouched.

---

## Obligations carried from the split — adjudicated against the actual delivered seam

The placeholder carried eight obligations "not a spec". Each is adjudicated here; reductions are
explicit and evidence-backed.

| # | Carried obligation | Adjudication against the delivered seam |
|---|---|---|
| 1 | The `/mcp/` endpoint and the World-authored MCP HTTP handler | **CARRIES, REDUCED.** The endpoint is `POST /mcp/` (D108-6). The handler is NOT World-authored: `mcphttp.NewHandler` owns HTTP. World authors only the three adapters + the mapping (D108-2/3). AC-ADMIT-108/AC-WIRE-108. |
| 2 | MCP JSON-RPC dispatch: initialization, tools/list, tools/call | **CARRIES, SEAM-OWNED.** The seam also dispatches ping and notifications (package doc, V9). World writes zero dispatch code. Asserted by conformance (AC-WIRE-108). |
| 3 | Dispatch-bound envelope framing (`WriteMCPEnvelope`/`RequestID` used from a served handler) | **REDUCES TO CONFORMANCE.** The seam uses `protocol.RequestID` + `protocol.WriteMCPEnvelope` + `CallbackMessage` itself on every host-failure and over-cap path (V13); World code never calls them. The obligation becomes: (a) assert the frozen envelope observable on a host-callback failure, and (b) a source-level assertion that World's diff hand-forms no JSON-RPC envelope bytes (AC-WIRE-108, MUT-ENVELOPE-HANDROLL). |
| 4 | SSE stream lifetime: `/mcp/`-route-local `ResponseController` relaxation, finite stream-lifetime maximum, OS-level closure on expiry | **REDUCED — the long-lived stream does not exist; r2 replaces the false envelope claim with the derived budget.** The seam is POST-only, stateless, no GET stream, no server-initiated messages (V9); its SSE response is ONE event — no relaxation code and no stream-lifetime constant exist. r1 asserted "the handler returns inside the frozen D7 `writeTimeout` (30 s) by construction"; that was FALSE (three sequential `hostcall.Run`s sum to 33 s of inner-bounded callbacks, V33/V36) and is WITHDRAWN. What survives — now DERIVED (D108-10's budget table) and MEASURED (the production-constants leg), not asserted: the whole POST is aggregate-bounded at `invokeDeadline` (20 s), so a blocked-callback POST is ALWAYS answered with the frozen envelope ≤ 20 s after handler entry and its response write has the ~10 s margin the landed `invokeDeadline < writeTimeout` invariant guarantees — DELIVERED inside the frozen window on every leg where the body arrives within the aggregate; the one client-induced exception (a body dripped past the window) fails closed, truncates at most the client's own window, and still leaks nothing (D108-10 leg 2, R108-6). Conditional on the callback's ctx cooperation (D108-5's separated guarantees; the landed row-106 discipline supplies it), the callback returns, the slot frees, and the next request succeeds (AC-INVOKE-108's deadline test; MUT-CB-BACKGROUND severs exactly the cooperation leg; MUT-AGGREGATE-REMOVED severs the budget leg). `MUT-LEAK-SSE-CONN` and `MUT-SSE-REST-DEADLINE` are **RETIRED** (below), their surviving kernels re-anchored as `MUT-GET-SERVED` + `MUT-DEADLINE-RELAX-MCP`. |
| 5 | Cross-surface criterion: per session, MCP tools/list name set == A2A skills[].id set | **CARRIES, RESTATED MODULO `M`** (the F6 grammar conflict, V16: IDs contain `.`/`/`, `ValidateMCPName` forbids them; the A2A doc itself recorded this residual "travels to the MCP child", V26). Restated: per session, `tools/list` name set == `M`(skills[].id set), with `M` injective on that set — hence equal cardinalities and `decode(M(ids)) == ids`. Order is NOT part of the criterion (A2A emits registry-bytewise order; MCP emits `CallerSurface`'s name-sorted order — both measured, V17). AC-CROSS-108. |
| 6 | SSE-framing conformance (`event: message` + `data:` JSON against the frozen P6.A fixture) | **CARRIES, RE-ANCHORED TO THE POST RESPONSE.** The seam's success response IS the SSE frame (V13). Conformance asserts it from World's mount; the expected bytes are **generated by the seam itself** (a recorded golden produced by a generator run of the pinned seam, checked in as produced — the hand-typed-fixture ban; upstream's own `handler_test.go` is the generator precedent). AC-WIRE-108, MUT-PLAIN-JSON (restated). |
| 7 | AC14 SSE/REST deadline separation (route-local relaxation; frozen D7 constants byte-unchanged) | **REDUCED TO THE 405 + NO-RELAXATION PAIR.** The relaxation this row forbade was for a GET stream the seam does not have; GET is refused 405 (D108-6, V23) and the frozen D7 constants stay byte-unchanged (AC-ROUTE-108; the landed `TestBoundedWaitsAndBodyLimit` pins them). The source-level no-`ResponseController`-relaxation assertion extends to this diff (MUT-DEADLINE-RELAX-MCP). |
| 8 | RED mutations `MUT-PLAIN-JSON`, `MUT-LEAK-SSE-CONN`, `MUT-SSE-REST-DEADLINE` | **RESTATES/REPLACES, ANCHORED TO THIS DIFF.** `MUT-PLAIN-JSON` restated (a World-authored wrapper rewrites the SSE response → framing test reds); `MUT-LEAK-SSE-CONN`/`MUT-SSE-REST-DEADLINE` RETIRED (obligation 4's reductions); their surviving kernels are `MUT-GET-SERVED` (GET must stay refused) and `MUT-DEADLINE-RELAX-MCP` (no relaxation in the diff). The full mutation set below is derived from what the milestones fix. |

---

## Non-negotiables (carried; none weakened)

- **No local JSON-RPC/MCP dispatch implementation** (parent P1 / Design Freeze / `DESIGN.md` §3.7
  protocol-native rule). The delivered seam IS the codec; World's diff contains no JSON-RPC
  parser, no SSE writer, no envelope formatter (asserted at source level, AC-WIRE-108).
- **No MCP SDK import into the gated daemon graph** (charter clauses 2 + 3; the measured
  28-violation OAuth closure was why `#885` was filed). First-party: `mcphttp`'s closure is
  stdlib + `protocol` + `hostcall` + `mcphttp` (V8); the allowlist extension admits exactly those
  package paths (AC-DEP-108).
- **Narrow PACKAGE-path allowlist discipline**: the seam is admitted by package path(s), never the
  module root; the narrowness tests are extended to cover the new paths and keep refusing
  root/facade widening (AC-DEP, MUT-ALLOWLIST-ROOT-MCP, MUT-FACADE-IMPORT-MCP).
- **Session carrier is D-WORLD-26 ARM A**: `Authorization: Bearer <session-credential>`, fail
  closed on absent/malformed/unknown/expired; never an API key; no alternate header
  (D108-4, AC-CARRIER).
- **The invocation path is the row-106 invocation coordinator** — resolver → capability predicate
  (`transitionreg.Request.Allowed` over `broker.Allows`) → coordinator invoke, one bounded context
  per callback (D108-8/9). `host/projection` is NOT a second policy engine; `broker.Allows` stays
  the only capability decider.
- **The row-23 store policy**: every `*Store` read is ctx-first and rejects deadline-free calls
  with `ErrNoDeadline`; the design invents no new store seam — it reuses the landed bounded reads
  (`transitionreg.NewReader`, the daemon's `readStore` heads seam) under the Runner's
  deadline-carrying contexts (D108-5, V19).
- **Bounded waits, twice over — per leg and per POST (r2's answer to the quorum's second REJECT;
  carried as a non-negotiable so it cannot be traded away in review)**: every callback on `/mcp/`
  is bounded by min(its inner frozen scalar, the Runner timeout, the D108-10 aggregate's remainder),
  and the whole POST is bounded by the aggregate (`invokeDeadline`) plus, behind it, the frozen
  transport clocks (`readTimeout`/`writeTimeout`). NO wait on this surface is asserted bounded:
  each is bounded by a named mechanism in D108-10's table (V33-V36), and the production-constants
  test leg plus `MUT-AGGREGATE-REMOVED` + `MUT-RESOLV-NO-BUDGET` + `MUT-TOOLS-NO-DEADLINE` +
  `MUT-RUNNER-UNBOUNDED` are the non-vacuity proof. The D108-10 aggregate is a REQUIRED element
  of the mount (removing it re-opens the 33 s > 30 s window overrun), not an optimization.
- **The aggregate wrapper owns no wire bytes**: it derives a context and delegates — no status,
  no header, no body — pinned by the same source-level scan family as AC-WIRE-108's
  no-envelope/no-relaxation assertions, so the no-local-codec rule above stays intact (D108-10).

---

## Premise Verification Log

Legend: **[R]** re-derived first-party in this session at World `dev` `8a9f961` (or the module
cache, as named); **[C]** controller-measured iteration 209 (cited with the controller's command;
not re-measured here — this sandbox has no network and the directive bounds the premise-7 probe
to reads); **[U]** UNMEASURED-IN-SANDBOX with the exact command the controller must run; never
approximated. Every empty/negative read is paired with a same-call known-positive control.

| # | Claim | Command | Observed output | Control (same call) |
|---|---|---|---|---|
| V1 [C] | `#885` is CLOSED as delivered in `v0.47.2`; `#764` also CLOSED | controller: `gh issue view 885 --repo sunholo-data/ailang --json state,title,comments` and `gh issue view 764 …` (also charter row 108: "issue CLOSED COMPLETED", PR #1369 squash `aa3063c37`) | `885: CLOSED` (delivered); `764: CLOSED`, 6 comments | `#764` known-positive fired in the same sweep |
| V2 [C] | proxy.golang.org AND sum.golang.org answer 200 for BOTH `v0.47.2` and `v0.48.0`; upstream latest release is `v0.48.0` (pushed 2026-09-29T09:52:19Z) | controller: `curl -sI https://proxy.golang.org/github.com/sunholo-data/ailang/@v/v0.47.2.info` (+ `.mod`, `v0.48.0`, and the `sum.golang.org` lookups) | 200/200 for both tags (one transient curl 000 on v0.48.0 re-measured 200 twice) | the 2026-09-28 charter 404 caveat's own re-check (iteration 207: "200/200") |
| V3 [C] | the serveapi/protocol tree is UNCHANGED between `v0.47.2` and `v0.48.0` (34-commit compare touches zero serveapi/protocol files) | controller: the compare between the two tags | zero serveapi/protocol files in the delta | — (paired with V4's first-party identity below) |
| V4 [R] | the five pinned protocol files are byte-identical at `v0.33.2` and `v0.47.2`, so the bump moves the A2A wire code by zero bytes | `cmp -s` each of `a2a_wire.go descriptor.go envelope.go interfaces.go descriptor_test.go` between `…/ailang@v0.33.2/serveapi/protocol/` and `…/ailang@v0.47.2/serveapi/protocol/` | all five `IDENTICAL` (v0.47.2 adds only the `hostcall/` and `mcphttp/` subdirectories, absent at v0.33.2) | same-call `ls` of both directories (5 files at v0.33.2; +2 subdirs at v0.47.2) |
| V5 [R] | World's go.mod pins `ailang v0.33.2`; only two direct deps; `go 1.26.6` | `head -8 go.mod` | `require ( github.com/sunholo-data/ailang v0.33.2; modernc.org/sqlite v1.54.0 )`; `go 1.26.6` | the full file was read; indirect block enumerated |
| V6 [R] | CI pins go1.26.6 at all four sites (P6.T discharged) | `grep -n "GOTOOLCHAIN\|go-version" .github/workflows/ci.yml` | `22: GOTOOLCHAIN: go1.26.6`, `38: go-version: '1.26.6'`, `120: GOTOOLCHAIN: go1.26.6`, `127: go-version: '1.26.6'` | the four sites are exactly P6.T's table |
| V7 [R] | the baseline gated closure at HEAD is **256 packages**; the ONLY upstream package in it is `serveapi/protocol` | `GOTOOLCHAIN=go1.26.6 go list -deps ./host/daemon/... ./cmd/ailang-worldd/... \| wc -l` and the same piped to `grep 'github.com/sunholo-data' \| sort -u` | `256`; upstream hits: exactly `github.com/sunholo-data/ailang/serveapi/protocol` (plus 17 `ailang-world` packages) | the upstream grep's positive control is the protocol line itself; toolchain `go version` = `go1.26.6 darwin/arm64` |
| V8 [R] | `mcphttp`'s import closure is stdlib + exactly `protocol` + `hostcall` + `mcphttp`; `hostcall`'s is stdlib + `protocol` + `hostcall` — the SDK-free claim | `cd …/ailang@v0.47.2 && GOTOOLCHAIN=go1.26.6 go list -deps ./serveapi/protocol/mcphttp` (and `./serveapi/protocol/hostcall`) | non-stdlib packages listed: exactly `…serveapi/protocol`, `…serveapi/protocol/hostcall`, `…serveapi/protocol/mcphttp` (hostcall: the first two) | the same list's stdlib majority (runtime, net/http, …) proves the instrument enumerated |
| V9 [R] | the seam's scope: POST-only, no sessions, no GET stream, no server-initiated messages; `SupportedVersions = 2025-11-25, 2025-06-18, 2025-03-26` | read `mcphttp/handler.go:1-34` | package doc verbatim (quoted in "The delivered seam"); `var SupportedVersions = []string{"2025-11-25", "2025-06-18", "2025-03-26"}` (at :30); `defaultVersion = "2025-03-26"` (at :34) | same read: the `writeMethodNotAllowed` call at handler.go:73-76 + its definition at wire.go:97-102 |
| V10 [R] | every `Config` field is required; no unbounded default | read `mcphttp/handler.go:36-63` (Config :37-44, NewHandler :48-63) | `NewHandler` switch refuses nil Resolver/Tools/Invoker/Runner and blank Agent name/version | same read: `hostcall.New` refuses `timeout <= 0` and `maxConcurrent <= 0` (runner.go:23-31) |
| V11 [R] | authorize runs resolve → tools → surface BEFORE Content-Type/Accept/version checks — denials never leak transport detail | read `mcphttp/handler.go:59-118` | `ServeHTTP`: body-cap → `authorize` → `isJSONContent` → `acceptsBoth` → version → `serveBody` | same read: the body (≤4 MiB) is read before authorize; the JSON-RPC message is NOT parsed before it |
| V12 [R] | transport preconditions: Content-Type `application/json`; Accept must list BOTH media types; unsupported version header → 400 | read `mcphttp/wire.go:106-127`, handler.go:90-105 | exact precondition order and messages as quoted in "The delivered seam" | `acceptsBoth`/`isJSONContent` parse via `mime.ParseMediaType` |
| V13 [R] | success = ONE SSE event then handler returns; 202 for all-notification POSTs; host failure = frozen 200/-32603 envelope for the WHOLE POST; unknown tool = -32602 | read `mcphttp/wire.go:79-102`, handler.go:180-183, methods.go:72-101 | `writeSSE` (text/event-stream, no-cache, nosniff, `event: message\ndata: …\n\n`), `writeAccepted` (202), `WriteMCPEnvelope` (200, -32603, `CallbackMessage`), `errorResponse(msg.ID, codeInvalidParams, "unknown tool %q")` | upstream's own `TestOversizedBodyIsEnvelope`, `TestWireContract` in the same package (read) |
| V14 [R] | the seam honors exactly 401/403 as HTTP statuses; 400 is not honored | read `protocol/envelope.go:52-61` + handler.go:113-118 | `AuthorizationStatus`: `if status == http.StatusUnauthorized \|\| status == http.StatusForbidden { return status }; return 0` | `protocol.AuthorizationError` implements `HTTPStatus()` (interfaces.go:28-40) |
| V15 [R] | the invoker result contract: non-empty, valid JSON → `content[0].text` + `structuredContent`; empty/invalid → -32603 | read `mcphttp/methods.go:83-100` | as quoted (`hostcall.Run` wraps `Invoke` at :83-85; the result contract at :89-100) | same read: `hostcall.Run` wraps `Invoke` |
| V16 [R] | the two name grammars conflict: `ValidateMCPName` = `^[a-zA-Z0-9_-]{1,64}$`; World IDs admit `[a-z0-9_-./]`, 1..128 bytes, segments 1..32 | read `protocol/descriptor.go:143-147` + `host/transitionreg/transitionreg.go:320-345` | both grammars verbatim (`mcpToolNameRegex`; `validateID` charset and bounds) | the A2A doc's F6 row (controller's throwaway test: `"tools.echo"` REJECTED) |
| V17 [R] | `CallerSurface` requires a JSON-object `type:"object"` InputSchema, sorts by name, refuses duplicates; registry fixtures carry `InputSchema: []byte("{}")` today | read `protocol/descriptor.go:30-66` + `grep -n "InputSchema" host/daemon/registry_publisher_test.go host/transitionreg/*_test.go` | `input["type"] != "object"` → refusal; `sort.Slice` by Name; `duplicate tool name %q`; fixtures: `{}` at registry_publisher_test.go:94,171,204,275 and transitionreg_test.go:422 | `codec.go:341` `parseJSON(d.InputSchema, maxSchemaCanonical)` (valid-JSON, ≤64 KiB — but NOT type-checked) |
| V18 [R] | the A2A existence proof's shape: one bounded ctx per route; resolver → allowedDescriptors (head pre-check + `NewRequest` + `Allowed()`) → `coord.Dispatch` | read `host/projection/projection.go:220-223` (AgentCard; `MaxWait` ctx at :221), `:265-268` (A2A; `InvokeWait` ctx at :266, ResolveContext at :268), `:339-346` (`coord.Dispatch(ctx, coordinator.Call{…})` at :339; `Connection: close` on ctx error at :346) | `AgentCard` (`MaxWait` ctx), `A2A` (`InvokeWait` ctx; ResolveContext; JSON-RPC codes; dispatch through the coordinator) | same read: `dispatchError` taxonomy at :391 + `allowedDescriptors`' B3 head pre-check at :376-377 |
| V19 [R] | the daemon's mounts and bounds: the frozen D7 block; the additive A2A mount; `isProtected` = POST /v1/commit only; the store guard | read `host/daemon/daemon.go:85-160, 315-395, 555-610, 711-740` + `host/store/store.go:236-244` | `invokeDeadline = 20s`(:98), `writeTimeout 30s`(:97), `readDeadline 10s`(:142), `credentialBudget 3s`(:153), `maxCommitBytes 8MiB`(:116); `projection.New` at :571-598; routes at :713-729; `isProtected` at :738-740; `ErrNoDeadline` + `requireDeadline` at store.go:236-244 | same reads returned the ten `/v1` patterns and `GET /workbench` |
| V20 [R] | row 39's resolver contract: `ResolveContext` returns typed denials, a store error is never a denial; the A2A route reuses it under D-WORLD-26 | read `host/authority/resolver.go:28-159` + projection.go:269-285 | `ResolveOutcome{Success\|Denied}`, exactly one set, never (nil,nil); DenialAbsent/Malformed/Unknown/Expired; header shape `Bearer <64-hex>` | the A2A denial constants at projection.go:73-88 |
| V21 [R] | the Runner's semantics: per-call `WithTimeout`; a slot is held until host code returns (even after timeout); over-cap → `ErrCallbackCapacity` | read `hostcall/runner.go:15-67` | `Run` verbatim as summarized in D108-5; the Runner's own doc comments: "A slot remains occupied until host code actually returns, even after timeout" (:15-16) and "In-process callbacks cannot be forcibly terminated; a non-cooperative callback keeps its slot" (:38-39) | `protocol.ErrCallbackCapacity` → `CallbackMessage` → "host callback capacity exceeded" (envelope.go) |
| V22 [R] | the row-106 coordinator is LANDED and assigns this row the `protocol.Invoker` adapter | read `host/coordinator/coordinator.go:41-112` (Config/New/Call/`InvocationID`/`validTaskID`) and `:196` (`func (c *Coordinator) Dispatch(ctx context.Context, call Call) (Result, error)`) + `design_docs/planned/w-transition-invocation-coordinator.md:184,426` | `Dispatch(ctx, Call)` verbatim; doc: "Row 108 adds a thin adapter when it lands (R-106-6), and a method declaration named Invoke is not a selector, so TR.C permits it"; residuals line 426: "R-106-6: protocol.Invoker adapter for MCP. Owner: row 108" | HEAD `git log --oneline -3`: `8a9f961 record(208): row 106 landed…`, `27f2574 row 106: invoke published transitions over A2A` |
| V23 [R] | a method-pattern mismatch answers 405 with an `Allow` header (the GET-refusal instrument) | `go doc net/http.ServeMux` + read `$(go env GOROOT)/src/net/http/server.go:2699-2705` | the mux's 405 branch verbatim: `w.Header().Set("Allow", strings.Join(allowedMethods, ", ")); Error(w, StatusText(StatusMethodNotAllowed), StatusMethodNotAllowed)` | same source read shows the pattern-matching doc ("A pattern with the method GET matches both GET and HEAD…") |
| V24 [R] | P6.D is landed: ONE package-path line + narrowness test, both gates green at HEAD | `grep -n "serveapi/protocol\", // w-a2a\|func TestAilangProtocolAdmissionIsNarrow\|func TestDaemonDependencyAllowlist" host/daemon/daemon_test.go` + `GOTOOLCHAIN=go1.26.6 go test ./host/daemon/ -run 'TestDaemonDependencyAllowlist\|TestAilangProtocolAdmissionIsNarrow' -count=1` | line 787 (the pinned package-path entry), 897, 975; test run: `ok github.com/sunholo-data/ailang-world/host/daemon 0.439s` | the allowlist run logs its own package count (non-vacuity leg refused an empty list) |
| V25 [R] | the matcher admits subpackages of a listed path (prefix semantics) — mcphttp/hostcall are already prefix-admitted; explicit manifest lines are still added | read `host/daemon/daemon_test.go:830` (the matcher line) + `:991` (the P6.D control comment) + the landed control leg (`protocol/subpkg` admitted) | `d == m \|\| strings.HasPrefix(d, m+"/")`; the P6.D control comment: "a future subpackage rides the same prefix" | the refused leg (`internal/apiserver`, `serveapi`, cloud path) in the same test |
| V26 [R] | the A2A card emits IDs VERBATIM, never through `ValidateMCPName`, and records this conflict as the MCP child's residual | read `host/projection/projection.go:26-29, 240-242` + `design_docs/implemented/w-a2a-session-projection.md` residuals | "(F6: never through CallerSurface / ValidateMCPName)"; A2A doc: "protocol.CallerSurface's MCP name grammar rejects World stable IDs (travels to the MCP child)" | the card's skills[] loop in the same read |
| V27 [R] | the frozen P6.A fixture's MCP leg: the observed ambient tool names and the CF-D-3 rule ("the fixture must become a real test in the same PR") | read `design_docs/planned/w-mcp-projection.md:873-876` + `design_docs/world-mission-log-archive.md:3153-3155` | recorded probe: unfiltered `tools/list` = 27 tools including `eprintln`, `exit`, `flush`, `print`, `printErr`, `println`, `readLine`, `writeBytes`, `submit_feedback`; `--routes-only` still lists `submit_feedback`; CF-D-3 verbatim | the same record's A2A leg is the landed A2A AC4 (names absent from the card) |
| V28 [C] | the fleet probe (INHERITED, re-measured where readable): pin bump + mcphttp import took a throwaway worktree 254 → 256 packages (mcphttp + hostcall), both dependency gates green unedited | **controller-reported fleet probe** — quoted in the charter's queue-STATUS entry for row 108 (`design_docs/world-mission.md:1720`: "The fleet reports an upstream probe: … took the build from 254 to 256 packages (`mcphttp` + `hostcall`), with both dependency-gate tests green unedited") — NOT the queue row itself (world-mission.md:1878), which carries only the delivery record and never says 254/256; THIS repo's remeasure commands: baseline `GOTOOLCHAIN=go1.26.6 go list -deps ./host/daemon/... ./cmd/ailang-worldd/... \| wc -l` (= **256**, V7, re-run green this iteration) and, for the bump shape, V29's throwaway-worktree probe (`go get github.com/sunholo-data/ailang@v0.47.2 && go mod tidy` + a throwaway `mcphttp` import, then the same list plus the two gate tests) | as reported by the controller (fleet probe, recorded at world-mission.md:1720): "took the build from 254 to 256 packages (mcphttp + hostcall), with both dependency-gate tests green unedited"; the probe's 254 base was a different worktree, so its absolute numbers are NOT this repo's; what carries is the SHAPE: exactly two new packages, both under the admitted prefix (V8 first-party) | V7 + V24 are the first-party baseline and gates-green controls |
| V29 [C] | the POST-BUMP gated closure count and both gates green unedited, at THIS HEAD | controller-measured iteration 209, throwaway detached worktree `.wt-world-iter209-probe` at `8a9f961`: `go get github.com/sunholo-data/ailang@v0.47.2 && go mod tidy`, a throwaway production probe-import of `mcphttp` in `host/daemon`, then `go list -deps ./host/daemon/... ./cmd/ailang-worldd/... \| wc -l` and `go test ./host/daemon/ -run 'TestDaemonDependencyAllowlist\|TestAilangProtocolAdmissionIsNarrow' -count=1` | **256 → 258** [go list -deps]; the additions are EXACTLY `…/serveapi/protocol/mcphttp` and `…/serveapi/protocol/hostcall` (the upstream `grep` shows exactly 3 packages: protocol, hostcall, mcphttp); `go build ./host/daemon/` clean; **both gates green UNEDITED: `ok github.com/sunholo-data/ailang-world/host/daemon 0.442s`** — the landed prefix matcher (V25) admits both subpackages | the pre-bump 256 (V7) re-measured in the same probe session (256 before the bump); the upstream grep's positive control is the protocol line itself |
| V30 [C] | go.sum movement under the bump (indirect roots) | controller-measured iteration 209, same probe worktree: `git diff go.mod go.sum` after the tidy | go.mod: the pin line only (`ailang v0.33.2` → `v0.47.2`); go.sum: 12 lines moved — the ailang pin's two lines plus **`golang.org/x/sys` v0.47.0 → v0.48.0** (the one indirect root the bump drags; already allowlisted). NO new module root enters the graph | any NEW module root outside `allowedDepModules` reds `TestDaemonDependencyAllowlist` by name — that test ran green in V29, so no new root exists (the red would have been the control) |
| V31 [C] | pre-land fetchability re-check (the charter's own rule) | controller-measured iteration 209 (pre-quorum): `curl -s -o /dev/null -w '%{http_code}' https://proxy.golang.org/github.com/sunholo-data/ailang/@v/v0.47.2.info` and `https://sum.golang.org/lookup/github.com/sunholo-data/ailang@v0.47.2` | **200 / 200** (both answers in the same session; one transient curl `000` on the v0.48.0 lookup re-measured 200 twice — a network blip, not a proxy state) | V2's 200/200 is the same-session repeat; the executor re-runs this immediately before the bump commits |
| **V32 [R]** (r2) | `procbound.MaxOutstanding = 8` exists, and is the process-wide limit on **subprocess reservations** (active children + retained background waiters) — NOT the Runner's slot cap; the two are different mechanisms | `sed -n '1,60p' host/procbound/procbound.go` + `grep -n "MaxOutstanding" host/procbound/procbound.go` + `grep -rn "procbound" host/ --include=*.go \| grep -v '^host/procbound' \| grep -v _test` | `const MaxOutstanding = 8` at **:31**; doc comment **:29-30**: "MaxOutstanding is the process-wide limit on reservations: admitted children not yet reaped, including those whose reap a background waiter owns"; package doc **:8-10**: "MaxOutstanding bounds active children plus retained background waiters process-wide … Admit is a nonblocking compare-and-swap reservation, not an observation"; the ONLY production caller is `host/broker/handlers.go:121` (`procbound.Admit()` before `cmd.Start()`) with `:132` `procbound.Wait(cmd.Wait, time.Until(cleanupDeadline), release)` — a reservation is held Admit→REAP, which can OUTLIVE the calling handler via the background waiter, so a Runner slot (callback start→return, runner.go:15-16) and a reservation are different lifetimes | the same sweep's `sed -n '100,170p' host/broker/handlers.go` shows the reservation is taken per subprocess start, not per HTTP request; `ErrCleanupBacklog` (:26) and `ErrCleanupIncomplete` (:22) are the refusal/cleanup errors, distinct from `protocol.ErrCallbackCapacity` |
| **V33 [R]** (r2) | one `tools/call` POST runs up to THREE sequential `hostcall.Run` calls, EACH with a fresh `context.WithTimeout(ctx, runner.timeout)` — so with no outer deadline the POST is bounded by the SUM of per-call timeouts (60 s at a bare 20 s Runner; 33 s with all inner bounds), and `r.Context()` carries NO deadline of its own | `sed -n '60,135p' …/ailang@v0.47.2/serveapi/protocol/mcphttp/handler.go` + `sed -n '60,105p' …/mcphttp/methods.go` + `sed -n '1,70p' …/hostcall/runner.go` + `sed -n '1,70p' …/protocol/envelope.go` | `ServeHTTP` reads the body BEFORE authorize (`io.ReadAll` at :77, `h.authorize` at :86); `authorize` runs `hostcall.Run(r.Context(), …, ResolveSession …)` then `hostcall.Run(r.Context(), …, Tools …)` then `protocol.CallerSurface` (CPU-only, no `Run`); `callTool` runs `hostcall.Run(req.ctx, …, Invoke …)` (methods.go:83-85); `Run` derives `callCtx, cancel := context.WithTimeout(ctx, runner.timeout)` **per call** (runner.go:42) and short-circuits the slot wait on parent expiry (`if ctx.Err() != nil { return zero, ctx.Err() }`, runner.go:47-49) | same reads: the final select returns `callCtx.Err()` (runner.go:61-66); `CallbackMessage(context.DeadlineExceeded)` → "host callback timed out" (envelope.go:27-29) — so an aggregate parent deadline renders as the SAME frozen envelope |
| **V34 [R]** (r2) | the transport clocks and WHEN they start: `WriteTimeout` is armed at the END of header parsing (before the handler runs and before the body is read) — body-read time is charged against the write window; `ReadTimeout` is armed at request start and bounds the body read | `sed -n '975,1000p' $(go env GOROOT)/src/net/http/server.go` + `sed -n '1035,1048p' $(go env GOROOT)/src/net/http/server.go` + `grep -n "wholeReqDeadline\|SetReadDeadline\|SetWriteDeadline" $(go env GOROOT)/src/net/http/server.go` + `go version` | `readRequest`: `wholeReqDeadline = t0.Add(ReadTimeout)` (**:982-985**), applied as the body-read deadline `c.rwc.SetReadDeadline(wholeReqDeadline)` (**:1041**); `WriteTimeout` is set in a **defer** `c.rwc.SetWriteDeadline(time.Now().Add(d))` (**:986-989**) which runs when `readRequest` RETURNS — i.e. the write clock starts at header-read completion, handler entry ≈ write-deadline start | `go version` = `go1.26.6 darwin/arm64` (the pinned toolchain, V6); the same file's doc: "WriteTimeout … is reset whenever a new request's header is read" (server.go:3003-3004) |
| **V35 [R]** (r2) | the landed whole-POST comparator + the B4 credential-budget discipline + the margin invariant the MCP mount reuses: the A2A route serves its ENTIRE POST inside ONE bounded ctx at `invokeDeadline`; the session middleware bounds each credential lookup at `credentialBudget`; `New` validates `invokeDeadline < writeTimeout` | `sed -n '240,300p' host/projection/projection.go` + `sed -n '30,90p' host/daemon/middleware.go` + `sed -n '555,610p' host/daemon/daemon.go` (exact-line `grep -n` for each cited anchor in the same sweep) | A2A: `ctx, cancel := context.WithTimeout(r.Context(), h.invokeWait)` (projection.go:266) with resolve + `MaxBytesReader` body decode + dispatch all inside it; `InvokeWait: invokeDeadline` (daemon.go:592); middleware: `lctx, cancel := context.WithTimeout(r.Context(), m.budget)` around `ResolveContext` (middleware.go:55), `budget = credentialBudget` (daemon.go:730, :153); `New`: `if invokeDeadline <= 0 \|\| invokeDeadline >= writeTimeout { … construction error }` (daemon.go:555-556) | the same middleware read's timeout arm renders 503 "never 401" (middleware.go:59-62) — the B4 rule the resolver adapter's inner bound mirrors; the A2A route's 20 s ctx under the 30 s window is the margin the aggregate inherits |
| **V36 [R]** (r2) | **the derived whole-POST budget** (the derivation row the kimi REJECT demanded): the Σ of the three inner bounds is **33 s**, EXCEEDS the frozen 30 s `writeTimeout`, so the inner bounds alone do not compose into a whole-POST bound; the D108-10 aggregate (`invokeDeadline`, 20 s) is the binding constraint on the SUM, and the response-write margin is `writeTimeout − invokeDeadline` = **10 s** | derived by arithmetic from constants READ in V19/V32/V33/V34/V35 (no new measurement): `credentialBudget` 3 s + `readDeadline` 10 s + `invokeDeadline` 20 s; `writeTimeout` 30 s; GOROOT clock starts per V34 | Σ inner bounds = 3 + 10 + 20 = **33 s > 30 s** (the r1 claim "returns inside the frozen 30 s by construction" is refuted by this row); aggregate = 20 s clips the Σ; margin = 30 − 20 = **10 s**, the same margin the landed `invokeDeadline < writeTimeout` validation (V35) guarantees for `/a2a/` | the empirical control is M108-4's production-constants leg (measured at the production scalars, not shrunk) and its killer `MUT-AGGREGATE-REMOVED`: with the aggregate removed the POST's callbacks run to their full inner bounds and the leg observes the 33 s shape — past the window — and REDs |

---

## Conflict Surface

**Touched:**

- `host/projection/` — two new files (`mcpname.go`, `mcp.go`) + the package doc comment (A2A →
  A2A + MCP) + `projection.Config` (the injected frozen scalars `CallbackTimeout`,
  `MaxCallbacks`, `CredentialBudget`, `ReadDeadline`, `InvokeDeadline`, `WriteWindow` — r2 grows
  this from two fields to six, all values the daemon already owns; NO new daemon constant) +
  `projection.New` (Runner + `mcphttp.Config` construction + the D108-10 aggregate wrapper + the
  `InvokeDeadline < WriteWindow` re-validation) + the new tests (`mcpname_test.go`, `mcp_test.go`,
  `mcp_conformance_test.go`). The landed A2A handlers and their tests are untouched.
- `host/daemon/daemon.go` — the `Config`→`projection.Config` scalar wiring (the frozen constants
  already in the D7 block and `procbound.MaxOutstanding`, by name) and one `mux.HandleFunc`
  registration (+ its comment block). `host/daemon/daemon_test.go` — two allowlist manifest lines
  + the narrowness test's new controls/manifest leg (+ the route-table test's one new entry).
  **The D7 constant block itself is byte-untouched** (the wrapper TIGHTENS a derived context; it
  relaxes nothing — AC-ROUTE-108's controls stay).
- `go.mod` / `go.sum` — the pin bump (V29/V30) and the indirect movement it drags.
- CI — **no workflow change**: the toolchain floor is discharged (V6), the proxy fetchability is
  controller-verified (V2/V31), and the gates that enforce this design are the EXISTING two
  dependency tests plus the new package tests inside the existing `go test ./...` job.

**NOT touched (and asserted untouched):**

- The frozen `/v1/*` route table and the D7 constant block (AC-ROUTE-108; landed
  `route_table_test.go` + `TestBoundedWaitsAndBodyLimit` are the controls). r2 adds the aggregate
  WITHOUT touching this block: every number in D108-10 is an existing frozen constant read by
  name (`invokeDeadline`, `credentialBudget`, `readDeadline`, `readTimeout`, `writeTimeout`);
  the only new identifier is the derived `postBudget` local in `host/projection`, and it derives
  to `invokeDeadline`.
- The store and the broker — no new seam, no schema change, no capability change; reads ride the
  landed `transitionreg.NewReader` + `readStore` heads seam (D108-5, V19).
- The A2A surface's behavior (landed quorum-cleared ACs stay green unedited).
- `world/*.ail` and `scripts/verify_ail.sh` — no new law (protocol/session invariants are
  host-boundary behavior, the parent's own P8 rationale; the mapping is pure Go, unit-tested).

---

## Milestones

Each milestone is independently CI-green and mergeable (hard repo convention); each acceptance
criterion names a RED mutation that must FAIL when applied and PASS when reverted, with the test
that kills it (rule-3i discipline).

### M108-1 — The pin bump + the explicit seam admission (~0.1d)

**Deliverable**: `go.mod` pins `github.com/sunholo-data/ailang v0.47.2` (with the pre-land
fetchability re-check, V31); `allowedDepModules` gains exactly two lines —
`github.com/sunholo-data/ailang/serveapi/protocol/mcphttp` and
`…/serveapi/protocol/hostcall` — each with its justification comment (the manifest is the
by-name review surface even though the landed prefix matcher already admits them, V25);
`TestAilangProtocolAdmissionIsNarrow` is extended with (a) by-name admitted controls for both new
paths and (b) a manifest-shape leg asserting the three exact ailang package-path entries exist;
the refused set stays (root, `serveapi` facade, `internal/apiserver`, cloud).
**Files**: `go.mod`, `go.sum`, `host/daemon/daemon_test.go`.
**Acceptance**: V29's controller re-measure lands in the PR description (closure 256→expected
258, gates green unedited); no new allowlist root. Mutations: **MUT-ALLOWLIST-ROOT-MCP** (replace
the three package-path lines with the module root → the narrowness test REDs: `internal/apiserver`
no longer refused); **MUT-ALLOWLIST-DROP-SEAM-LINES** (delete the two new lines → the
manifest-shape leg REDs); **MUT-FACADE-IMPORT-MCP** (import `github.com/sunholo-data/ailang/serveapi`
anywhere in the daemon graph → `TestDaemonDependencyAllowlist` REDs naming the intruders);
**MUT-PIN-REGRESS** (revert the pin to v0.33.2 → `go build ./...` REDs: the `mcphttp` import is
unresolvable).

### M108-2 — The tool-name mapping + schema normalization (~0.15d)

**Deliverable**: `host/projection/mcpname.go` — `encodeToolName(id) (string, error)`,
`decodeToolName(name) (string, error)` (the D108-2 table; decode fails LOUDLY on any invalid
escape or uppercase char), the ≤64-byte refusal, and the schema normalizer (conforming object
schemas verbatim; else the minimal permissive `{"type":"object"}`).
**Files**: `host/projection/mcpname.go`, `host/projection/mcpname_test.go`.
**Acceptance**: a generated corpus test proves `decode(encode(id)) == id` for every ID-shape
admitted by `validateID` (dot/slash/underscore/segment/length edges), injectivity over all pairs,
and the refusal arms (collision, >64, invalid decode input). Mutations:
**MUT-MAP-IDENTITY** (emit IDs verbatim, no mapping → the injectivity/corpus test REDs on the
dot/slash corpus, and later the conformance suite: `CallerSurface` refuses `tools.echo`);
**MUT-MAP-COLLIDE** (make encode many-to-one, e.g. also map `-`→`_` → the corpus injectivity test
REDS: two distinct IDs encode identically); **MUT-MAP-DECODE-INVALID** (decode accepts an
unmapped `_x` → the round-trip test REDs); **MUT-MAP-NO-LEN-GUARD** (drop the 64-byte refusal →
the refusal test REDs: a long-ID descriptor is silently served instead of refused).

### M108-3 — The three adapters + the mount + the aggregate wrapper (~0.3d)

**Deliverable**: `host/projection/mcp.go` — the resolver adapter (ResolveContext → typed
`*protocol.AuthorizationError{401}` per denial kind, constant messages, log seam for resolution
failures; **the `credentialBudget` inner bound (r2)** — `context.WithTimeout` around
`ResolveContext` exactly as the landed session middleware does, middleware.go:55, V35 — and on
inner-budget expiry the raw deadline error, never a typed denial, so the seam renders the frozen
envelope and the 401s stay reserved for typed denials), the tools adapter (the landed B3
head-pre-check + `NewRequest` + `Allowed()` → encode
→ normalize → `protocol.ToolDescriptor` set; fail-closed typed errors on collision/overflow; the
`readDeadline` inner bound), the invoker adapter (decode → fresh `NewRequest`+`Allowed()` →
object-check on arguments → server-minted task id → `coordinator.Dispatch` → `InvocationResult`;
the `invokeDeadline` inner bound); `projection.Config` + `New` (Runner construction from the
daemon's injected frozen scalars, validation — including the `InvokeDeadline < WriteWindow`
re-check, daemon.go:555-556's landed invariant, V35); **the D108-10 aggregate wrapper**
(`projection.New` returns the seam handler wrapped in the context-deriving middleware — its
entire body is the `WithTimeout` + `defer cancel` + `h.ServeHTTP(w, r.WithContext(ctx))` triple,
no status, no header, no byte written); the daemon's route registration. Tests: the denial
matrix, the per-session exact set, the invoke matrix (listed/revoked/guessed), the
mount-validation test, **and the resolver's bounded-wait fault-injection leg** (a blocked
credential store read exceeds the 3 s `credentialBudget` → typed "timed out" path, never a 401,
never an unbounded wait). **Files**: `host/projection/mcp.go`, `host/projection/projection.go`
(doc comment + Config), `host/projection/mcp_test.go`, `host/daemon/daemon.go`,
`host/daemon/route_table_test.go`.
**Acceptance**: AC-CARRIER-108/AC-ADMIT-108/AC-INVOKE-108/AC-BOUNDS-108/AC-ROUTE-108 below. Mutations:
**MUT-RESOLV-DENY-DROPPED** (the adapter ignores `out.Denied` and returns a session →
`TestMCPDenialMatrix` REDs); **MUT-KEY-AS-SESSION** (the adapter accepts a non-Bearer/static key
as a session → the denial matrix REDs — D-WORLD-26 constraint (i) restated for `/mcp/`);
**MUT-RESOLV-NO-BUDGET** (r2: drop the resolver adapter's inner `credentialBudget` bound so the
lookup rides the bare Runner timeout → the resolver bounded-wait fault-injection leg REDs: a
blocked credential read is answered at ~20 s instead of ≤ 3 s, and the leg's
elapsed-against-`credentialBudget` assertion fails — this is the mutation that pins Kimi's
"resolver has no tighter discipline" catch); **MUT-MCP-PROTECTED** (add `/mcp/` to `isProtected` → the denial-shape test REDs: the REST
envelope replaces the seam's 401); **MUT-TOOLS-UNFILTERED** (Tools returns every descriptor, not
the `Allowed()` set → `TestMCPListExactSetPerSession` REDs); **MUT-SCHEMA-PASSTHRU** (emit the
registry's raw `{}` schema → the surface-conformance test REDs: every POST gets the frozen
envelope because `CallerSurface` refuses); **MUT-TOOLS-NO-DEADLINE** (drop the tools adapter's
inner `readDeadline` bound → the bounded-wait fault-injection test REDs: a blocked head read
exceeds the configured bound); **MUT-INVOKE-NO-READMIT** (the invoker skips the fresh
`Allowed()` check and dispatches the decoded ID regardless → `TestMCPInvokeListedThenRevoked`
REDS: a call succeeds that must fail closed); **MUT-CB-BACKGROUND** (the invoker dispatches
under `context.Background()` instead of the callback ctx → the deadline test REDs: the POST
itself still answers within the aggregate — that guarantee is the Runner's own select under the
D108-10 parent and holds regardless (D108-5's separated guarantees) — but the callback now ignores the deadline,
never returns, and its slot never frees, so the subsequent-call leg fails with capacity rather
than success); **MUT-RUNNER-UNBOUNDED** (construct the
Runner with a zero/negative timeout or slot count → the mount-validation test REDs:
`projection.New`/`hostcall.New` refuse — never "unbounded").

### M108-4 — Conformance + the cross-surface criterion + the frozen fixture + the production-constants leg (~0.35d)

**Deliverable**: `host/projection/mcp_conformance_test.go` — the seam's wire contract observed
from World's mount over a real seeded registry and two sessions with unequal capability sets:
initialize/ping/tools/list/tools/call happy paths; the SSE frame asserted against a golden
GENERATED by the pinned seam (recorded in-repo as produced by a generator run — never
hand-typed); the frozen envelope on a host-callback failure; the transport 400s (Content-Type,
Accept, version header); the 4 MiB cap; GET → 405 with `Allow: POST`; the source-level
no-envelope-handroll and no-`ResponseController`-relaxation assertions over the diff (+ the
same-scan leg asserting the D108-10 wrapper writes no status, header, or byte — the wrapper is
pinned as context plumbing only); the
cross-surface test (`tools/list` name set == `M`(card `skills[].id` set), decode-equality, both
surfaces over ONE daemon); the frozen P6.A MCP leg as a live test (the ambient names of V27
absent from BOTH sessions — CF-D-3 discharged); the deadline/socket test in the landed
row-106 shape — **now TWO legs, mechanism and arithmetic (r2's answer to Kimi's
sidestep catch)**:
- *the mechanism leg* (`TestMCPInvokeDeadlineFreesSlot`, fast, CI-friendly): a blocking fixture
  transition with a **shrunk** Runner timeout and a shrunk aggregate — the frozen "host callback
  timed out" envelope within the bound, ALWAYS, per D108-5's separated guarantees; the capsule
  exits under the callback ctx, the callback returns and the slot frees (the cooperation leg
  MUT-CB-BACKGROUND severs), and a subsequent call succeeds.
- *the production-constants leg* (`TestMCPPostBudgetProductionConstants`, ~21 s wall, the leg the
  r1 draft lacked): the SAME mount with the daemon's REAL production scalars —
  `credentialBudget` 3 s, `readDeadline` 10 s, `invokeDeadline` 20 s, `writeTimeout` 30 s,
  aggregate = 20 s, Runner slot cap 8 — and fault injection that saturates authorize exactly as
  the budget table says (a blocked credential read held to its full 3 s inner bound, a blocked
  head read held to its full 10 s inner bound), then a blocked invoke. Assertions, each against
  D108-10's table: (a) the POST is answered with the frozen `200/-32603 "host callback timed
  out"` envelope; (b) **elapsed ≤ the aggregate (20 s) + a small scheduling tolerance** — the
  measured answer must NOT wait for the invoke's own 20 s inner bound (which would put the POST
  at 33 s), because the aggregate clips the leg at ≈ 7 s into the invoke; (c) the response was
  WRITTEN inside the simulated 30 s write window (the leg records the handler-entry-to-write
  interval and asserts the ≥ ~10 s margin the landed invariant guarantees); (d) the invoke
  callback OBSERVED a ctx cancelled at the aggregate, not at its inner bound (the callback's
  `ctx.Err()` fires while `time.Since(start) < invokeDeadline` — the direct observable that the
  parent deadline is the binding one); (e) a subsequent POST succeeds (the slot freed — the
  cooperation discipline under the aggregate).
**Files**: `host/projection/mcp_conformance_test.go`, the golden fixture file (generated), the
mutation-harness entries.
**Acceptance**: AC-WIRE-108/AC-CROSS-108/AC-FIXTURE-108/AC-ROUTE-108/AC-BOUNDS-108 (the
production-constants leg is AC-BOUNDS-108's empirical control, V36). Mutations: **MUT-PLAIN-JSON** (restated —
a World-authored wrapper around the mounted handler rewrites the SSE response to a plain
`application/json` body → the framing test REDs); **MUT-GET-SERVED** (register `GET /mcp/` on
the MCP handler → `TestMCPGetRefused` REDs: GET must answer 405, not 200); **MUT-DEADLINE-RELAX-MCP**
(add a `ResponseController.SetWriteDeadline(time.Time{})` call anywhere in the diff → the
source-level scan REDs naming the site; the landed D7 constants test is the same-run control);
**MUT-AGGREGATE-REMOVED** (r2, the production-constants killer: mount the bare seam handler with
no aggregate wrapper — or equivalently derive the wrapper ctx from `writeTimeout` instead of
`invokeDeadline` — → the production-constants leg REDs on assertion (b): authorize consumes its
full 3 s + 10 s inner bounds and the invoke then runs to its own 20 s, so the POST answers at
≈ 33 s > 20 s + tolerance, past the frozen window; the leg that r1 lacked is precisely what
makes this arithmetic un-hideable); **MUT-CROSS-SURFACE-DIVERGE** (filter the MCP tools differently from the A2A card, e.g. drop one
descriptor in the tools adapter → the cross-surface test REDs); **MUT-AMBIENT-HARDCODED**
(hard-code an ambient-named tool into the adapter → the P6.A-leg test REDs);
**MUT-ENVELOPE-HANDROLL** (hand-format a JSON-RPC error body in World code → the wire-owner
source assertion REDs). **Retired, with reason**: `MUT-LEAK-SSE-CONN` and `MUT-SSE-REST-DEADLINE`
were written against a long-lived GET SSE stream with a route-local deadline relaxation — the
delivered seam is POST-only and stateless and mounts no such stream (V9), so both mutations have
no production code to mutate; their surviving kernels (transport refused on GET; no relaxation in
the diff; the handler returns within the bound) are anchored as `MUT-GET-SERVED`,
`MUT-DEADLINE-RELAX-MCP`, and the M108-3/M108-4 deadline legs.

---

## Acceptance Criteria

| AC | Criterion (each named mutation must FAIL applied, PASS reverted) | Mutations (killer test) |
|---|---|---|
| **AC-DEP-108** | `go.mod` pins `ailang v0.47.2`; the closure over `./host/daemon/... ./cmd/ailang-worldd/...` gains exactly the two seam packages (controller re-measure V29); the allowlist manifest carries both package paths with justifications; the narrowness test's refused set and controls hold | MUT-ALLOWLIST-ROOT-MCP, MUT-ALLOWLIST-DROP-SEAM-LINES, MUT-FACADE-IMPORT-MCP, MUT-PIN-REGRESS (`TestAilangProtocolAdmissionIsNarrow`, `TestDaemonDependencyAllowlist`, `go build`) |
| **AC-MAP-108** | `M` is total on the ID grammar, injective, reversible with loud decode failure; collision and >64-byte names refuse the whole surface fail-closed; schemas normalize per D108-2 | MUT-MAP-IDENTITY, MUT-MAP-COLLIDE, MUT-MAP-DECODE-INVALID, MUT-MAP-NO-LEN-GUARD (`TestMCPNameRoundTrip`, `TestMCPNameRefusal`, corpus test) |
| **AC-CARRIER-108** | `/mcp/` reads ONLY `Authorization: Bearer <session-credential>` via the ONE landed resolver; absent/malformed/unknown/expired fail closed as the seam's 401 with constant per-kind messages; never an API key, never an alternate header; a denied request acquires NO registry/capability snapshot (counting observer reads 0) and never reaches dispatch; `/mcp/` is NOT in `isProtected` | MUT-RESOLV-DENY-DROPPED, MUT-KEY-AS-SESSION, MUT-MCP-PROTECTED (`TestMCPDenialMatrix`, `TestMCPDeniedNeverTouchesRegistry`) |
| **AC-ADMIT-108** | `tools/list` is exactly this session's `Allowed()` set under `M`, dynamic per request (head change without restart), zero tools at 200 on an absent head, fail-closed on genuine read failures; blocked transitions are absent; every emitted schema passes `CallerSurface` | MUT-TOOLS-UNFILTERED, MUT-SCHEMA-PASSTHRU, MUT-TOOLS-NO-DEADLINE (`TestMCPListExactSetPerSession`, `TestMCPSchemaNormalization`, bounded-wait fault-injection) |
| **AC-INVOKE-108** | `tools/call` decodes the name, re-admits against a FRESH snapshot (a listed-then-revoked name never dispatches), arguments must be a JSON object, the task id is server-minted, the ONLY write path is `coordinator.Dispatch` under the callback ctx; a blocked call's POST is answered within the **D108-10 aggregate** ALWAYS (`Run` returns on ctx Done, runner.go:61-66, under the aggregate parent — the r1 "within the Runner bound" phrasing is withdrawn: the Runner timeout bounds each LEG, the aggregate bounds the POST, V36), and the slot frees once the callback cooperates with the deadline ctx and returns — the landed row-106 ctx discipline is the cooperation, and MUT-CB-BACKGROUND severs it | MUT-INVOKE-NO-READMIT, MUT-CB-BACKGROUND (`TestMCPCallDispatch`, `TestMCPInvokeListedThenRevoked`, `TestMCPInvokeDeadlineFreesSlot`) |
| **AC-BOUNDS-108** | the Runner is constructed from the daemon's injected scalars — per-call timeout `invokeDeadline` (20 s, frozen), slot cap `8` **taken from `procbound.MaxOutstanding` as an alignment VALUE choice, not an identity: the Runner cap counts in-process callback slots, procbound counts process-wide subprocess reservations held to REAP — different mechanisms, honestly separated (V32, D108-5)** — and validated at `New`; **every one of the three adapters carries its own frozen inner bound (`credentialBudget` 3 s resolve / `readDeadline` 10 s tools / `invokeDeadline` 20 s invoke — r2: the resolver's was missing and is added), and the whole POST is bounded by the D108-10 aggregate (`invokeDeadline`), so min(inner, Runner, aggregate) holds on every leg and the Σ is clipped to 20 s inside the 30 s frozen window with the landed 10 s write margin (V36)**; every store read under `/mcp/` carries a deadline (`ErrNoDeadline` never observed on this surface); the derivation is pinned empirically by the production-constants leg at the REAL scalars, not a shrunk bound | MUT-RUNNER-UNBOUNDED, MUT-RESOLV-NO-BUDGET, MUT-AGGREGATE-REMOVED (mount-validation test; the resolver bounded-wait fault-injection leg; `TestMCPPostBudgetProductionConstants`; the `ErrNoDeadline` pin is the same-run control) |
| **AC-WIRE-108** | World's diff contains no JSON-RPC/SSE/envelope formatting and no deadline relaxation; the seam's wire behavior is asserted from the mount (SSE frame vs the seam-generated golden; 202 on notification-only POSTs; 400s on Content-Type/Accept/version; the 4 MiB cap; the frozen envelope on host failure; initialize/ping negotiation) | MUT-PLAIN-JSON, MUT-ENVELOPE-HANDROLL, MUT-DEADLINE-RELAX-MCP (`TestMCPWireConformance`, source-level scans) |
| **AC-ROUTE-108** | GET (and every non-POST) on `/mcp/` answers 405 with `Allow: POST`; the frozen `/v1/*` table, `GET /workbench`, and the D7 constants are byte-unchanged; the A2A routes behave identically to their landed tests | MUT-GET-SERVED (`TestMCPGetRefused`; landed `route_table_test.go` + `TestBoundedWaitsAndBodyLimit` are the frozen-surface controls) |
| **AC-CROSS-108** | per session, the MCP `tools/list` name set equals `M`(the A2A card `skills[].id` set) — the split's cross-surface criterion RESTATED MODULO `M` (obligation 5, the F6 grammar conflict); decode equality holds; cardinalities equal | MUT-CROSS-SURFACE-DIVERGE (`TestMCPToolsListMatchesAgentCard`) |
| **AC-FIXTURE-108** | the frozen P6.A two-session fixture's MCP leg is a LIVE test (CF-D-3): two sessions with unequal capability sets observe different tool sets, and the recorded ambient names (`eprintln`, `exit`, `flush`, `print`, `printErr`, `println`, `readLine`, `writeBytes`, `submit_feedback`) are absent from both | MUT-AMBIENT-HARDCODED (`TestMCPAmbientExportsAbsent`) |

Estimate honesty: the four milestones sum to ~0.9d against the charter row's `~1d`; r2 adds the
aggregate wrapper + resolver inner bound (~0.05d in M108-3) and the production-constants test leg
(~0.05d in M108-4, a ~21 s wall-clock test — the price of pinning the budget arithmetic at the
REAL scalars rather than a shrunk bound); the biggest
single risk (the schema normalization, V17) is quarantined in M108-2 with its own mutation set,
and the seam itself carries upstream tests that World's conformance suite complements rather
than re-derives.

---

## Out of scope / residuals

- **R108-1 — the A2A remainder of clause 6**: landed by
  [`w-a2a-session-projection.md`](../implemented/w-a2a-session-projection.md) (card + `/a2a/`)
  and row 106 (invocation). Nothing here touches it; clause 6 is satisfied only when BOTH halves
  land (the parent's scope statement).
- **R108-2 — MCP has no idempotent-resend channel**: each `tools/call` mints a fresh task id, so
  a client retry after a lost response executes a NEW invocation (D108-8). A per-call idempotency
  key would need either an arguments-key convention or an upstream seam change — neither is
  invented here; owner: a future row if an MCP client needs it. **The natural future mechanism is
  already landed and is named here WITHOUT expanding this row's scope**: same-task-id journal
  reconciliation — the row-106 coordinator already answers a resent task id from the journal and
  never re-executes it (`GetReceipt` hit → return the committed result, coordinator.go:206-215),
  the discipline the A2A route's own comment names ("The client can then retry with the same task
  ID after reconciling an uncertain durable outcome", projection.go:342-344). An MCP surface
  could one day become idempotent by returning the server-minted task id in the invocation
  result and letting a retry carry it back — but `tools/call` has no channel for that today (the
  result carries no task id, and `arguments` is the only client-controlled input), so the
  mechanism stays a named future option, not a design element of this row.
- **R108-3 — a registry schema with an invalid `x-mcp-header` annotation** breaks the MCP surface
  fail-closed (the normalizer fixes only the top-level shape; `CallerSurface` walks properties,
  V17). The refusal is loud (frozen envelope + one-line log naming the descriptor); the fix is
  the publisher's. Owner: the production-publisher row.
- **R108-4 — invocation IDs share the `a2a:` journal namespace**: MCP-initiated invocations are
  indistinguishable from A2A-initiated ones by ID prefix (the coordinator's `InvocationID` scheme
  is landed and frozen here); a per-surface prefix needs a coordinator API change. Owner: none
  unless replay/audit needs the distinction.
- **R108-5 — the `transition_fn` pin affordance is A2A-only** (MCP's `Invocation` carries no
  metadata channel). Owner: none; recorded so the asymmetry is deliberate, not forgotten.
- **R108-6 (r2) — the one leg where the frozen envelope may be UNDELIVERED, recorded honestly
  instead of hidden**: the seam reads the request body INSIDE `ServeHTTP` (the `io.ReadAll` at handler.go:77), and
  Go arms `WriteTimeout` at header-read completion (GOROOT server.go:986-989, V34), so body-read
  time is charged against the same 30 s write window the callbacks need. A client that drips its
  own body slow enough consumes the window itself: the D108-10 aggregate still fires, the read is
  still bounded by the transport's `ReadTimeout` clock (t₀+30 s), and when the read completes the
  POST fails closed instantly with the frozen "host callback timed out" envelope — but that final
  write may truncate at the socket, and the client sees a connection error rather than the
  envelope. This is client-induced (the server cannot deliver an answer to a client that spent the
  answer's window), it is bounded by the frozen transport clocks (the body read unblocks at
  t₀+30 s, the response write at header-read end+30 s; no zombie stream, no
  leaked slot — the callback ctx cooperation still frees the Runner slot), and it is IDENTICAL on
  the landed A2A route (whose body decode rides the same 20 s ctx under the same 30 s transport
  clocks), so `/mcp/` adds no new exposure. The World-authored wrapper cannot fix it without
  becoming a body-reading codec (which would violate the no-local-codec non-negotiable) or
  relaxing the frozen D7 clocks (forbidden). Owner: none — transport-fundamental; recorded so the
  r2 claim is exactly "delivered inside the window on every leg where the body arrives within the
  aggregate", never a universal claim.
- **Quorum-internal notes**: the quorum should scrutinize (a) the mount decision D108-3 (widening
  `host/projection`'s charter vs a third package), (b) the 401-for-malformed choice in D108-4
  (seam-imposed), (c) the retired mutations (obligations 4/7/8), and — added in r2 — (d) the
  D108-10 aggregate and its wrapper (a World-authored middleware around the seam handler: is a
  ~5-line context-deriving wrapper inside the no-local-codec rule? the design's answer is YES —
  it writes no status, no header, no byte, and is pinned by a source-level scan — but the quorum
  should ratify that reading explicitly), and (e) the slot-cap VALUE choice (Runner `8` =
  `procbound.MaxOutstanding` as an alignment, V32 — the quorum may prefer a different number;
  the design's claim is only that the two mechanisms must not be conflated) — each adjudicated
  above with its evidence row.

---

## Quorum verification log

**Round 1 — BLOCKED.** Author `pi:ollama/glm-5.3:cloud`; Gemini and Kimi rejected. Gemini found the unverified `procbound.MaxOutstanding` claim; Kimi found that sequential callback bounds can exceed the frozen 30 s write window. Sonnet was absent on quota and Astra unreachable. The one permitted designer revision added V32–V36, separated the two caps, and added an aggregate request budget plus a production-constants test.

**Round 2 — BLOCKED.** The same two external seats rejected. Gemini found `MUT-GET-SERVED` vacuous: changing the mux registration still routes GET to upstream `mcphttp`, which independently returns 405, so the named test stays green. Kimi found that the default 2025-03-26 protocol permits batch requests, while V33 and D108-5/D108-10 model only one `tools/call` per POST. The controller confirmed in `mcphttp/handler.go:140-198` that authorization runs once and `serveMessages` dispatches batch items sequentially; `methods.go:83-85` invokes once per valid `tools/call` item. A host error aborts the entire POST with the frozen envelope, while an unknown method returns an item error. Thus a batch has two authorization callbacks plus K invocation callbacks, all under the aggregate context if installed. These facts are **review findings, not yet incorporated into the design or acceptance tests**. Sonnet was absent on quota and Astra unreachable again; two external reviewers were present. Round 2 met the skill's one-revision/one-requorum stop rule. The substantive batch omission is outside its narrow-refinement carve-out. The next step requires human review of the park and a further designer revision, followed by full quorum. No planner, executor, or evaluator was routed, and no code was landed.

---

## Related Documents

- [`w-mcp-projection.md`](w-mcp-projection.md) — the split parent: enablers landed (P6.A/P6.T),
  the round-3 objection verbatim, the carried obligations this doc adjudicates, and the P6.A
  fixture record (V27)
- [`../implemented/w-a2a-session-projection.md`](../implemented/w-a2a-session-projection.md) —
  the A2A sibling: the landed existence proof (handlers, P6.D admission, F6 name-grammar
  residual that travels here)
- [`w-transition-invocation-coordinator.md`](w-transition-invocation-coordinator.md) — row 106:
  the invocation contract, the deadline discipline, and R-106-6 (the `protocol.Invoker` adapter
  this design owns)
- [`../world-mission.md`](../world-mission.md) — charter clause 6; queue row 108 (the delivery
  record, the proxy caveat, the re-measure mandate); D-WORLD-26 (ARM A); D-WORLD-5
- upstream [`sunholo-data/ailang#885`](https://github.com/sunholo-data/ailang/issues/885) — the
  delivered ask, closed in tag `v0.47.2`; [`#764`](https://github.com/sunholo-data/ailang/issues/764)
  → `serveapi/protocol` — the precedent admission this design extends