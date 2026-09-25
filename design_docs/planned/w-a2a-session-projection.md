# w-a2a-session-projection — Session-Scoped A2A Agent Card + Fail-Closed `/a2a/` (SPLIT child of `w-mcp-projection`)

**Status**: Planned — **UNBLOCKED (row 39 `w-session-authority` landed `a036062`, iter-181);
REVISED iter-187 against row 39's ACTUAL landed contract; DESCOPED by controller SPLIT to the
buildable now — a session-scoped A2A agent card plus a fail-closed `/a2a/` endpoint on row 39's
`host/authority.Resolver`; transition INVOCATION is moved verbatim to a "Deferred to invocation"
section gated on a coordinator that does not exist; NOT yet quorum-cleared; QUORUM ROUND 1 (iter-187 r1): three REJECTions
(gpt6-astra, gemini-3-1-pro, oc-glm-5-2) MEASURED and each answered — see "Quorum round 1".** The A2A projection
design substance (the parent's `P6.B-A2A`, reviewer-refined across four quorum rounds) is
carried where it still applies and preserved verbatim where it must wait for invocation.
**Item**: split child #2 of `w-mcp-projection` (charter clause 6, queue row 5; the long-time
blocker was queue row 39 `w-session-authority`)  
**Clause**: clause-6 — *"the transition registry is served over MCP (capability-filtered per
session) and an A2A agent card is published; no new wire protocols"*. The A2A half of clause 6
that ships HERE is **publishing the card** (F12: clause 6's A2A half requires publishing the
card; it does not by itself require A2A task invocation).  
**Estimate**: **~0.65d = P6.B-A2A-CARD card ~0.4d + P6.A-CTX bounded-resolution prerequisite ~0.1d + P6.D dependency admission
~0.15d** (descoped from the inherited ~0.7d invocation-bearing P6.B-A2A line)  
**Author**: rotation designer, split #2 round, 2026-08-26 (iteration 126); this revision is the
revision mandated by the doc's own blocking predicate, executed iter-187  
**Verified against**: World `dev` at **`a8e12fd`** (this worktree; `go.mod` is `go 1.26.6`)
**Date**: 2026-08-26 (authored); **Revised**: 2026-09-25 (iteration 187)

## Why this doc exists — the objection, verbatim (historical record, preserved)

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
verified equal to the artifact's fields in the split-round session).

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

The round-4 `proposed_fix`'s second half is literally the disposition the split executed: the
gap ("no inbound credential→session resolution") was BLOCKED-era and is **now closed — row 39
`w-session-authority` landed `a036062` (iter-181) and supplies exactly the resolver the fix
demanded** (measured below, F1–F3). The `catch`'s demanded overlap analysis is now answered in
this revision's Conflict Surface. Everything the objection blocked is re-scoped against the
landed contract; invocation (the coordinator half) is deferred, not lost.

## Revision record — iter-187 controller SPLIT (descope to buildable-now)

The blocking predicate named three conditions for UNBLOCK and a revision against row 39's
ACTUAL contract. All three hold at `a8e12fd`, measured first-party this session:

1. **Row 39 landed** (`a036062`, iter-181, judged PASS ×3), giving `host/authority` its
   `Resolver` (F1), the daemon's `SessionMiddleware` (F2), and the wired resolver construction
   (F3).
2. **Parent enablers green**: `go.mod` is `go 1.26.6` (P6.T landed); P6.V is the parent's
   Z3-verified commit-boundary law and is consumed only by the invocation that is DESCOPED here.
3. **P6.D has NOT landed** (`git grep 'serveapi/protocol' -- go.mod host/daemon/daemon_test.go`
   rc=1) so dependency admission ships HERE (F10).

The controller's iter-187 routing decision, applied without re-litigation, is a **SPLIT** of this
doc's remaining scope into what is **buildable now** versus what is **deferred to invocation**:

- **(A) Milestone P6.D** — dependency admission of the pinned `serveapi/protocol` package,
  carried with updated closure numbers (F10).
- **(B) Milestone P6.B-A2A-CARD** — a session-scoped A2A **agent card** + a fail-closed `/a2a/`
  endpoint, both over row 39's resolver. Buildable now because neither requires a server-side
  invocation coordinator (F7).
- **(C) Deferred to invocation** (new queue row, gated on a coordinator that does not exist) —
  every invocation-only responsibility, moved verbatim: Decision 3 responsibility 4, P4's
  dispatch half, `protocol.Invoker`, AC13, the OS-level socket-closure obligation, and all
  invocation-success AC text.

F7 is the load-bearing fact: **there is no server-side invocation coordinator.** A2A task
invocation would have to be invented end-to-end (`propose → verify → commit` dispatch has no
landed entrypoint and the transition registry is never populated in production, F8). Publishing
the card and rejecting `/a2a/` invocation fail-closed is honest, additive, and executable now.

## Verification Log — VERIFIED BY CONTROLLER iter-187, re-run first-party at `a8e12fd`

Every fact below was measured by the controller today at `dev a8e12fd` and is cited as a
**VERIFIED BY CONTROLLER iter-187** row; rows marked **[RR]** were additionally re-run
first-party in this session at the same worktree with the recorded output. Empty/zero results
are paired with a known-positive control in the same row so a zero is a measurement, not a
broken instrument.

### Row 39's landed contract — the resolver this doc consumes (F1–F3)

| # [RR] | Claim | Command | Observed OUTPUT |
|---|---|---|---|
| **F1** | `host/authority/resolver.go: interface Resolver { Resolve(header string, now int64) ResolveOutcome }`; `authority.New(st CredentialStore) Resolver`; `ResolveOutcome{Success *SessionBinding; Denied *DenialKind}` (one set); `SessionBinding{EpisodeID string; Caps []broker.Capability; ExpiresAt, CreatedAt int64}`; `DenialKind` = DenialAbsent/Malformed/Unknown/Expired with `String()` = `SessionAbsent/InvalidSession/SessionUnknown/SessionExpired`; takes the WHOLE `Authorization` header value | `grep -n 'type Resolver interface\|func New(\|ResolveOutcome\|type SessionBinding' host/authority/resolver.go` + `read` | **confirmed**: `Resolve(header string, now int64) ResolveOutcome` (`:93`); `New(st CredentialStore) Resolver` (`:89`); `ResolveOutcome` at `:77`; `SessionBinding` at `:67`; `DenialKind` consts at `:36-42`; String() `:` maps to SessionAbsent/InvalidSession/SessionUnknown/SessionExpired |
| **F2** | `daemon/middleware.go: SessionMiddleware.Wrap(protected func(*http.Request) bool, mux)` enforces the boundary only where `d.isProtected(r)` is true; today ONLY `POST /v1/commit` (`daemon.go:583`); success carries the binding via `authority.WithBinding`, handlers read it with `authority.FromContext`; absent/unknown/expired→401, malformed→400, REST APIError envelope (`writeAPIError`) | `grep -n 'func (m \*SessionMiddleware) Wrap\|isProtected\|"/v1/commit"\|writeAPIError' host/daemon/middleware.go host/daemon/daemon.go` + `read middleware.go` | **confirmed**: `Wrap(protected func(*http.Request) bool, mux http.Handler) http.Handler` (`middleware.go:42`); `isProtected` = `POST` && `/v1/commit` (`daemon.go:582-583`); `case DenialAbsent, DenialUnknown, DenialExpired → StatusUnauthorized`, `case DenialMalformed → StatusBadRequest` via `writeAPIError` (`middleware.go:52-62`); `WithBinding`/`FromContext` pair (`:63`) |
| **F3** | the daemon constructs the resolver at `daemon.go:447` | `grep -n 'resolver: authority.New' host/daemon/daemon.go` | **confirmed**: `447: scanTimeBudget: integrityScanTimeBudget, resolver: authority.New(s),` |

### Upstream `github.com/sunholo-data/ailang/serveapi/protocol` (F4–F5)

| # [RR] | Claim | Command | Observed OUTPUT |
|---|---|---|---|
| **F4** | the five protocol files (`a2a_wire.go, descriptor.go, descriptor_test.go, envelope.go, interfaces.go`) have BYTE-IDENTICAL blob SHAs at tags `v0.33.2` and `v0.42.0` — pinning `v0.33.2` stays valid, nothing new shipped | `gh api '…/contents/serveapi/protocol/<file>?ref=v0.33.2\|v0.42.0' --jq '.sha'` for all five | **confirmed**: identical SHAs at both tags, e.g. `interfaces.go` `dbe401d…86`, `a2a_wire.go` `48cedfb…81d` — equal at v0.33.2 and v0.42.0 (all five re-run) |
| **F5** | `interfaces.go`: `type Session any`; `SessionResolver{ResolveSession(ctx,*http.Request)(Session,error)}`; `ToolSource{Tools(ctx,Session)([]ToolDescriptor,error)}`; `Invoker{Invoke(ctx,Session,Invocation)(InvocationResult,error)}`; `AgentInfo{Name,Description,Version string}`; `AuthorizationError{Status int; Err error}` with `HTTPStatus()`. NO AGENT-CARD TYPE; upstream `serveapi/a2a_handler.go` builds the card as a `map[string]any` with keys `name,description,url,version,capabilities{streaming:false,pushNotifications:false,stateTransitionHistory:false},defaultInputModes,defaultOutputModes,skills[{id,name,description,tags,examples}]` | `gh api '…/interfaces.go?ref=v0.33.2' --jq '.content' \| base64 -d`; `…/a2a_handler.go?ref=v0.33.2' … \| grep 'map\[string\]any\|"name"\|"capabilities"'` | **confirmed**: both signatures and the exact card keys from `a2a_handler.go:79-84`; card is a `map[string]any` literal — no agent-card type anywhere |

| **F5b** | `a2a_wire.go` (verified v0.33.2): `type A2ARequest struct { JSONRPC string `json:"jsonrpc"`; Method string `json:"method"`; ID json.RawMessage `json:"id"`; Params json.RawMessage `json:"params"` }`; `A2ATaskSendParams{ ID string; Message A2AMessage; Metadata map[string]any }` (tags id/message/metadata); `A2AMessage{ Role string; Parts []A2AContent }`; `A2AContent{ Type, Text string; Data map[string]any }`. `func A2AError(w http.ResponseWriter, id json.RawMessage, code int, msg string)` and `func A2AResult(w http.ResponseWriter, id json.RawMessage, result any)` BOTH write HTTP **200**. Upstream `a2a_handler.go` reads the skill from `params.Metadata["skill_id"].(string)` — this doc does the same and says so | `gh api '…/contents/serveapi/protocol/a2a_wire.go?ref=v0.33.2' --jq .content \| base64 -d` | **VERIFIED BY CONTROLLER iter-187 r1**: exact struct/signature text above from `a2a_wire.go` at `v0.33.2`; `A2AError`/`A2AResult` both `w.WriteHeader(http.StatusOK)` then encode — so a `/a2a/` denial is a JSON-RPC error body at HTTP **200** |

### NAME-GRAMMAR conflict (F6)

| # [RR] | Claim | Command | Observed OUTPUT |
|---|---|---|---|
| **F6** | `protocol.CallerSurface` validates each `ToolDescriptor.Name` against `^[a-zA-Z0-9_-]{1,64}$`; World stable transition IDs (`host/transitionreg/transitionreg.go validateID :314`) allow `.`/`/` and 128 bytes. Results (controller's throwaway Go test): `"effect-paged"`→accepted, `"tools.echo"`→REJECTED, `"world/recovery-transition/v1"`→REJECTED. So "exclusively through CallerSurface" (old AC1) and Decision 4's "skills[].id equals the authorized registry ID set exactly" CANNOT both hold for real IDs | `gh api '…/descriptor.go?ref=v0.33.2' \| grep mcpToolNameRegex`; `grep -n 'func validateID' host/transitionreg/transitionreg.go` + `read :314` | **confirmed**: protocol `var mcpToolNameRegex = regexp.MustCompile("^[a-zA-Z0-9_-]{1,64}$")`; transitionreg `validateID` allows `.` and `/` separators, segments 1..32, id 1..128 bytes. The two grammars are mutually incompatible → **the card must NOT route IDs through `CallerSurface`** (AMENDED AC1); the MCP dispatch child inherits this and must solve naming itself |

### World premises (F7–F9)

| # [RR] | Claim | Command | Observed OUTPUT |
|---|---|---|---|
| **F7** | NO server-side invocation coordinator: `grep -rniE "func .*(propose\|coordinat)" host cmd` (non-test) → 0; `host/transitionreg` imported by NO production package (control: `host/broker/invoke_boundary_test.go`); `host/capsule`/`host/replay` have no production callers; the daemon's only mutation is `POST /v1/commit` (F2) | `grep -rniE "func .*(propose|coordinat)" host cmd --include='*.go' \| grep -v '_test.go'` (output count 0); `grep -rln 'host/transitionreg' host cmd --include='*.go'` | **confirmed**: non-test coordinator funcs = **0**; `host/transitionreg` importers = `host/transitionreg/transitionreg.go` only (non-test) + the single test `host/broker/invoke_boundary_test.go` |
| **F8** | registry NEVER populated in production: `transitionreg.StoreReader.Publish` / `BuildNext` have 0 non-test callers (control: 23 test refs to `store.TransitionRegistryV1`); `ReadSnapshot` returns plain `errors.New("read transition registry: head is absent")` when no head (`transitionreg.go:~78`) — not typed | `grep -rn 'BuildNext' host cmd --include='*.go' \| grep -v '_test.go'`; `grep -rc 'TransitionRegistryV1' host cmd --include='*_test.go'`; `read transitionreg.go` | **confirmed**: `BuildNext` has only its definition, 0 callers; `cmd/world-publish` uses `broker.EffectRegistryPublish` (different type), does NOT import `transitionreg`; test-ref total = **23**; `ReadSnapshot` head-absent error at `:78` is `errors.New(...)` (plain) |
| **F9** | capability filtering available: `transitionreg.NewRequest(ctx, reader, caps CapabilitySource, now)` then `Request.Allowed()` returns descriptors whose `Access` `broker.Allows` admits, in registry bytewise order; `broker.NewSession(store,episodeID,grants,registry)` is a pure constructor; `(*broker.Session).CapabilitySnapshot(now)` satisfies `CapabilitySource` | `grep -n 'func NewRequest\|func (q Request) Allowed\|func NewSession\|func (s \*Session) CapabilitySnapshot\|type CapabilitySource' host/transitionreg/bind.go host/broker/broker.go` | **confirmed**: `NewRequest` `bind.go:118`, `Allowed` `bind.go:128`, `CapabilitySource` `bind.go:17`; `NewSession` `broker.go:87`, `CapabilitySnapshot` `broker.go:67` — F9 exactly |

### P6.D / P6.T / P6.V status (F10) and the MCP sibling (F11), clause 6 (F12)

| # [RR] | Claim | Command | Observed OUTPUT |
|---|---|---|---|
| **F10** | P6.D admission re-measured: `go get github.com/sunholo-data/ailang@v0.33.2` succeeds under go1.26.6 (`go.mod` is now `go 1.26.6`, P6.T landed); go.mod gains `+ailang v0.33.2`, bumps `go-isatty v0.0.20→v0.0.22`, `x/sys v0.46.0→v0.47.0`; closure over `./host/daemon/... ./cmd/ailang-worldd/...` moves **250 → 251** (ADD 1: `serveapi/protocol`, removed set empty); probing import + unchanged allowlist → `TestDaemonDependencyAllowlist` FAILS naming exactly that package; adding ONE package-path line passes. P6.D lands HERE: `git grep -n 'serveapi/protocol' -- go.mod host/daemon/daemon_test.go` → no match | `git grep -n 'serveapi/protocol' -- go.mod host/daemon/daemon_test.go`; `go list -deps ./host/daemon/... ./cmd/ailang-worldd/... wc -l`; `go test ./host/daemon/ -run TestDaemonDependencyAllowlist` | **confirmed**: `git grep` rc=1 (P6.D not landed); current closure = **250 packages** (`250` vs the doc's stale `249→250` — the landed baseline is 250, admission is 250→251); allowlist test currently **ok** (0.4s). Controller also re-measured the post-admission 251 and the RED-names-one-intruder probe |
| **F11** | sibling MCP dispatch child: `ailang#885` was CLOSED 2026-09-13 as "Delivered" (the unread mission-control message, iter-187), but F4 shows protocol unchanged since v0.33.2 → no MCP dispatch shipped; the MCP child stays BLOCKED | `gh api 'repos/sunholo-data/ailang/issues/885' --jq '.state,.title,.closed_at'` | **confirmed**: `closed`, title `serveapi/protocol has no MCP dispatch…`, `closed_at 2026-09-13T15:52:56Z`; protocol blob SHAs unchanged (F4). **Sentence: the MCP dispatch child stays blocked** — `ailang#885` closed without the MCP dispatch seam shipping |
| **F12** | charter clause 6 (`world-mission.md:1001`): "the transition registry is served over MCP (capability-filtered per session) and an A2A agent card is published; no new wire protocols". Its A2A half REQUIRES publishing the card; it does not by itself require A2A task invocation | `grep -n 'transition registry is served over MCP' design_docs/world-mission.md` + `read :1001` | **confirmed** at `:1001`; the card is the mandatory half, task invocation is not stated |
| row 39 | `w-session-authority` LANDED `a036062` (PR #141 squash, iter-181, judged PASS ×3) | `grep -n 'w-session-authority' design_docs/world-mission.md` | **confirmed**: queue row 39 record shows **[LANDED 2026-09-24 (iter-181) — PR #141 squash `a036062`**…**]** |

### Row 39's contract re-verified iter-187 r1 (VERIFIED BY CONTROLLER iter-187 r1)

Round 1 surfaced three objections the controller measured (rows **R1-CTX / R1-HEAD / R1-WIRE** +
**F5b**); each fix below is grounded in the cited measurement. The P6.A-CTX row is the load-bearing
one for the bounded-wait answer; B3 + AC-ABSENT-HEAD consume R1-HEAD; F5b consumes R1-WIRE.

| # | Claim | Command | Observed OUTPUT (VERIFIED BY CONTROLLER iter-187 r1) |
|---|---|---|---|
| **R1-CTX** | `Resolver.Resolve(header, now)` takes NO context (resolver.go:93), so it ignores any request deadline; `host/store/store.go:305` `db.SetMaxOpenConns(1)` = ONE connection, so a resolve can wait behind any in-flight query/commit; `(*Store).ResolveSession` (store.go:1101) already uses `QueryRowContext` so a ctx WOULD be honoured; resolver.go:130-134 maps a store error to `DenialUnknown` (401 "unknown credential") — a naive ctx fix would mis-report a deadline as a denial. REQUIRED FIX = milestone **P6.A-CTX** (`Resolver.ResolveContext(ctx, header, now)` returns a non-nil `error` on store failure, never a denial) | `read host/authority/resolver.go:129` + `read host/store/store.go:305` + `read host/store/store.go:1101` + `read host/authority/resolver.go:130-134` | **confirmed**: `row, ok, err := r.st.ResolveSession(context.Background(), credentialID)` (resolver.go:129); `db.SetMaxOpenConns(1)` (store.go:305); `s.db.QueryRowContext(ctx, ...)` (store.go:1101); `DenialUnknown` mapping (resolver.go:130-134). `Resolve` still lands on `/v1/commit` unchanged (F2); residual declared |
| **R1-HEAD** | GetRegistryHead seam: `store.go:636` `(s *Store) GetRegistryHead(ctx, name) (hashref.HashRef, bool, error)` — `sql.ErrNoRows`→`(zero,false,nil)` (:641-642), other scan error→wrapped error (:643-644), bad ref→error (:647-648); `daemon.go:331-337` `readStore` already includes `GetRegistryHead`; registry name `store.TransitionRegistryV1 = "world/transition-registry/v1"` (store.go:88) is the same name `transitionreg.ReadSnapshot` reads (transitionreg.go:74) | `read host/store/store.go:636-648` + `read host/daemon/daemon.go:331-337` + `grep -n 'TransitionRegistryV1' host/store/store.go host/transitionreg/transitionreg.go` | **confirmed** at all cited lines. FIX: verification row added here; B3 + AC-ABSENT-HEAD cite the lines and specify the head race (publish between check and `NewRequest`; see B3); a test covers "check says absent" |
| **R1-WIRE** | A2A wire types were unverified until row **F5b** (added): the struct/method signatures in `a2a_wire.go` and the HTTP-200 behaviour of `A2AError`/`A2AResult` | see F5b | see F5b (VERIFIED BY CONTROLLER iter-187 r1). FIX: JSON-RPC denial codes on `/a2a/` specified (-32001 / -32600, see Quorum round 1) + AC/MUT rows added |

## Row 39's landed contract — what this doc CONSUMES (this supersedes the blocked-era "Blocking predicate" section)

The blocked-era section measured what was ABSENT. Row 39 landed, so the predicate's signal is
now **unblocked** and the doc is revised against the landed shape. This doc still does NOT
design the session-authority boundary; it CONSUMES the landed contract verbatim:

- **F1** — `authority.New(st CredentialStore) Resolver`; `Resolver.Resolve(header string, now
  int64) ResolveOutcome`; exactly one of `Success *SessionBinding` / `Denied *DenialKind` set.
  `SessionBinding` carries `EpisodeID`, `Caps []broker.Capability`, `ExpiresAt`, `CreatedAt`;
  `DenialKind` = `DenialAbsent/Malformed/Unknown/Expired` (String `SessionAbsent/InvalidSession/
  SessionUnknown/SessionExpired`). It takes the **WHOLE** `Authorization` header value; there is
  no alternate-header and no API-key path inside the resolver by construction.
- **F2** — the daemon wires a `SessionMiddleware.Wrap(protected, mux)` enforcing the boundary
  only where `d.isProtected(r)` is true (today only `POST /v1/commit`); denials map
  absent/unknown/expired→**401**, malformed→**400**, rendered through the REST `writeAPIError`
  envelope; success carries the binding via `authority.WithBinding`/`FromContext`.
- **F3** — the daemon holds one resolver instance: `resolver: authority.New(s)` (`daemon.go:447`).
  The projection handler REUSES this instance (mount-time injection); it adds **no second
  credential path**, does **not** read `X-World-Session`, and **never accepts the serve-api key**
  (D-WORLD-26 = ARM A, both constraints, carried unchanged).

Because resolution is delegated to the landed resolver, the round-4 requirement (a premise row
"maps an opaque credential to episode ID, grants, expiry; absent/malformed/unknown/expired all
return typed denial") is satisfied as **landed code**, not a designed dependency.

## Scope carried from the parent (moved at split #2; substance unchanged; re-scoped iter-187)

Parent identifiers are **RETAINED** wherever they still apply; **AMENDED iter-187** where the
contract changed them; **MOVED → Deferred to invocation** where they wait on a coordinator (F7).
The parent keeps P1/P7/P8, AC10, AC15, AC17; the numbering gaps on both sides are deliberate and
declared.

### Premises (hard constraints; parent identifiers preserved)

- **P2 — exact per-session authority.** A session resolves to an immutable capability snapshot
  for one request. The same predicate filters A2A discovery (card skills) and the `/a2a/`
  admission check. A missing, expired, or unknown session fails closed. **(AMENDED iter-187:**
  the "invocation" half of the old predicate is renamed the `/a2a/` admission check, because
  invocation does not exist (F7).)
- **P3 — transition truth stays in World.** The projected entries come from the landed transition
  registry, not a startup copy of `.ail` exports.
- **P4 — propose → verify → commit remains mandatory.** **MOVED → Deferred to invocation** — the
  dispatch half of P4 names a coordinator that does not exist (F7). Its constraint that the
  projection receives no `store.Commit` / `SetRegistryHead` / REST `/v1/commit` authority is
  trivially satisfied by the descoped shape (the card route is read-only). The full P4 text is
  moved verbatim to the Deferred section.
- **P5 — worldd stays the sole writer.** The projection is hosted in the existing worldd process
  and uses its already-open store/broker handles. It never calls `store.Open`, starts a second
  writer, or weakens `WriterAlreadyActive`.
- **P6 — the landed REST v1 contract is frozen.** No route is renamed, removed, or repurposed.
  The A2A endpoints are additive and are NOT in `d.isProtected`. Existing REST callers see
  byte-compatible behavior. For the card route's denial body this doc deliberately **reuses the
  daemon's REST `daemon.APIError` envelope** (Decisions, B1) — reuse of the shipped writer on a
  NEW endpoint is not a change to any `/v1/` route, so P6 is respected.

### Decisions — the descoped card + fail-closed `/a2a/` projection (Decision 3 re-bound)

Row 39 landed `host/authority.Resolver`. The additive `host/projection` adapter is a **read-only
projection surface** with two responsibilities that are buildable now, mounted additively (B5):

1. **Resolve the session** — for BOTH routes, call the daemon's existing resolver instance:
   `header := r.Header.Get("Authorization"); out, err := resolver.ResolveContext(ctx, header, now)`
   (P6.A-CTX; `ctx` is the ONE bounded request context) where `now = time.Now().Unix()`.
   **D-WORLD-26 = ARM A** carried unchanged: **(i)** never an API key, **(ii)**
   fail closed on absent/malformed/unknown/expired, never degrading to an unauthenticated
   surface; Arm B `X-World-Session` is REJECTED and must not be read, even as a fallback.
   *(r2 carve-out, gemini-3-1-pro verbatim:)* A denial returns the F2 status mapping (401/400) on
   the *card route*, but emits protocol.A2AError at HTTP 200 with code -32001/-32600 on the
   */a2a/* route. See the **Denial matrix** below, which is authoritative.

   **Denial matrix (authoritative; r2 carve-out, gpt6-astra verbatim).** Card route uses the
   injected daemon APIError writer with HTTP 401 for absent/unknown/expired and HTTP 400 for
   malformed; /a2a/ uses protocol.A2AError at HTTP 200, with -32001 for absent/unknown/expired and
   -32600 for malformed, constant messages, and no REST envelope. -32001 is an
   implementation-defined JSON-RPC server-error code (the JSON-RPC 2.0 reserved range
   -32000..-32099). A denied request on either route never acquires a registry snapshot and never
   reaches skill admission.

   | Denial kind | Card route (`/.well-known/agent.json`) | `/a2a/` route |
   |---|---|---|
   | absent | HTTP 401, APIError `SessionAbsent` | HTTP 200, `A2AError` -32001, constant message |
   | unknown | HTTP 401, APIError `SessionUnknown` | HTTP 200, `A2AError` -32001, constant message |
   | expired | HTTP 401, APIError `SessionExpired` | HTTP 200, `A2AError` -32001, constant message |
   | malformed | HTTP 400, APIError `InvalidSession` | HTTP 200, `A2AError` -32600, constant message |
   **Body decision (B1):** deny on the **card route** (`/.well-known/agent.json`) is emitted with
   the daemon's **REST `daemon.APIError` envelope** by REUSING the daemon's `writeAPIError`
   writer, which the daemon injects into the projection config at mount. Justification against
   P6: P6 freezes the `/v1/` route table and its bytes; it does not forbid **reusing** the
   shipped envelope struct on a NEW non-REST endpoint. The card route is a plain-HTTP-JSON fetch
   (no A2ARequest framing), its denial vocabulary and status mapping are byte-identical to the
   already-proven middleware on `/v1/commit`, and reuse keeps one error shape (`SessionAbsent/
   InvalidSession/SessionUnknown/SessionExpired`) with NO parallel envelope in `host/projection`
   (AC1 — the APIError envelope is daemon-owned; projection never hand-formats it). The `/a2a/`
   route is the JSON-RPC surface and emits `protocol.A2AError`, never `APIError` — P6's "do not
   translate into APIError" applies to that JSON-RPC surface, not to a card fetch.
2. **Build the per-session card** (B2): on resolution success, construct a broker session from
   the binding (`broker.NewSession(store, b.EpisodeID, b.Caps, registry)`, F9 — pure constructor,
   no writes), capture **ONE** registry snapshot + **ONE** capability snapshot via
   `transitionreg.NewRequest(ctx, reader, caps, now)` (the broker session's
   `CapabilitySnapshot(now)` is the `CapabilitySource`), and take `Request.Allowed()`. The card's
   `skills[] =` those descriptors in **registry bytewise order** with `id` = the World stable ID
   **VERBATIM** (F6: never route the card through `protocol.CallerSurface` — its MCP name regex
   rejects real IDs; record: the MCP dispatch child inherits this incompatibility and must solve
   naming itself). Build the card as a `map[string]any` literal with the SAME keys upstream uses
   (F5), with `name`/`description`/`version` from a `protocol.AgentInfo` value. No parallel card
   struct type.
3. **Handle ABSENT registry head distinctly (B3):** an absent head (F8) is a legitimate state,
   distinct from a read failure. The projection checks `GetRegistryHead` (the daemon's reads
   seam) BEFORE building the request: `ok==false` → card with **ZERO skills** (HTTP 200, empty
   `skills` array, still authenticated against the session); `err != nil` or a later
   `NewRequest` failure → **5xx fail-closed**. This distinguishes the two without modifying the
   landed `transitionreg` (minimal-frozen-core: prefer a host-boundary injection over a change
   to landed core). A test must cover both. **Head race (specified, per R1-HEAD):** the head can be published between the `GetRegistryHead` check and `NewRequest`. If the check saw absent and `NewRequest` then succeeds, use `NewRequest`'s result (the head arrived — succeed). If the check saw a head and `NewRequest` fails, answer 5xx (a present-but-unreadable registry is a failure, not an absent state). A test covers "check says absent" at minimum.
4. **Serve `/a2a/` POST fail-closed (B4):** resolve the session exactly as (1); parse with
   `protocol.A2ARequest`; non-`"2.0"` → **-32600**; `method != "tasks/send"` → **-32601**; a
   `skill_id` not in this session's `Allowed` set → **-32602** `"not authorized"` (stale/guessed
   names get the same refusal); an AUTHORIZED `skill_id` → a structured A2A error with one
   **CONSTANT** message "transition invocation is not available in this daemon" — code
   **-32603** internal error, justified below — never a fake success, never a REST body, never a
   direct store write. Emit via `protocol.A2AError(w, id, code, msg)`.

**JSON-RPC code for an authorized-but-not-invocable skill (B4 decision):** **-32603 (Internal
error).** Justification: the other three standard JSON-RPC codes are already claimed for the
framing/method/params errors (-32600/-32601/-32602), and -32603 is the standard code for a
request that reaches the handler and fails at execution rather than at parse/params — exactly
this case. The MESSAGE does the real work and is a **single constant string** (never interpolated
with a skill name or request content), so the code's "internal" connotation is bounded and stable;
the code is spec-fixed (not a bespoke implementation-defined value that would be A2A-spec
novelty). It is emitted ONLY via `protocol.A2AError`; there is no success path for `/a2a/`
invocation in this daemon.

The authorization function is conceptually:

`visible(session, transition) = session is live AND transition capability ∈ session capabilities`.

That expression is explanatory prose, not AILANG source. The implementation reuses the landed
capability predicate (`transitionreg`/`broker.Allows`); it does not create a second policy
engine in `host/projection`.

**One snapshot per request.** Card generation may observe a newer registry on the next request,
but a single request never mixes registry or capability epochs. `/a2a/` builds its OWN fresh
per-request request-set at admission time (re-resolves the session and re-reads the snapshot);
possession of a previously listed skill name conveys no authority.

**Fail-closed cases.** Unknown/expired session, absent head (zero skills, still 200) vs. read
error (5xx), and an authorized-but-unsupported invocation (constant -32603) each have a
distinct, tested outcome. No case falls back to REST `/v1/` translation on the `/a2a/` JSON-RPC
surface, and no case does a direct store write from the projection.

### Decision 4 — Surface Identity, A2A Card, and Compatibility — carried, AMENDED iter-187

Tool identity is the transition registry's stable transition ID, emitted **verbatim** (F6). Display
text and schemas are projection metadata; they do not become authority. Ordering is deterministic
by stable ID so cards are diffable.

For a given session snapshot:

- A2A `skills[].id` equals the authorized registry ID set exactly, **verbatim** (AMENDED:
  NOT via `CallerSurface`; its MCP name regex rejects real World IDs — F6);
- every listed skill is discoverable through the same session, and `/a2a/` admits only listed
  IDs — an authorized-but-unsupported ID receives the constant not-available error, an unlisted/
  guessed/stale name receives `-32602`; and
- `std/io.*`, bare IO names (`exit`, `writeBytes`, etc.), and `submit_feedback` are absent unless
  World someday registers and authorizes transitions with those exact IDs. This item registers
  none.

The cross-surface obligation the pre-split AC3 stated — that the eventual MCP `tools/list` name
set must exactly equal the A2A `skills[].id` set per session — is NOT dropped: it travels with
the MCP dispatch child (`w-mcp-dispatch-projection.md`) and binds the MCP surface when it
becomes buildable. **It is NOT satisfiable through `CallerSurface` for real IDs (F6); the MCP
child inherits the grammar conflict and must solve naming itself** (recorded here and there).

The agent card is emitted by World's handler as a `map[string]any` literal with the SAME keys
upstream uses (F5), `name`/`description`/`version` from a `protocol.AgentInfo` value; the `/a2a/`
route constructs/parses exclusively through `protocol`'s `A2ARequest`/`A2AError` helpers. The
test must not assume a hand-authored JSON shape beyond protocol-required fields and the exact
skill-ID set.

### Bounded wait — carried from the parent's Decision 6, DESCOPED to resolution + snapshot reads

Every projection request derives one bounded context whose deadline is the earlier of the client
cancellation/deadline and a configured finite, positive server maximum (validated at startup;
zero/negative/omitted ≠ "unlimited"), passed WITHOUT replacement through both session resolution
and registry/capability snapshot acquisition. Resolution uses **P6.A-CTX**'s `ResolveContext(ctx,
header, now)` — NOT the context-less `Resolve(header, now)` the `/v1/commit` middleware still
uses; the store holds ONE connection (`SetMaxOpenConns(1)`, store.go:305), so the bounded context
bounds a wait behind any in-flight query/commit. `ResolveContext` returns a store error (incl.
`context.Canceled`/`DeadlineExceeded`) as a non-nil `error`, never a `DenialUnknown` 401. The
**card route** answers a resolution error **503** (store / unavailable upstream) or **504**
(deadline hit waiting on the single pooled connection); the **`/a2a/` route** answers `A2AError`
**-32603** with the constant message (F5b; `A2AError` always writes HTTP 200). **Declared
residual:** `/v1/commit`'s `SessionMiddleware` STILL resolves without a request deadline (`Resolve`
→ `ResolveContext(context.Background(), …)`, error mapped to `DenialUnknown` exactly as today);
this doc does not change it. Client disconnect cancels the bounded context; every dependency
observes cancellation and returns promptly. Resolution runs on the calling goroutine, no worker
spawned (`MUT-DROP-DEADLINE-PROJ`). The `/v1/commit`-bound and OS socket-closure material is
invocation-only, **MOVED verbatim** to the Deferred section.

### Files (rows moved from the parent's table; iter-187 descope)

| File | Milestone | Purpose |
|---|---|---|
| `host/projection/projection.go` | P6.B-A2A-CARD | Session-scoped adapter; World-owned card handler (`/.well-known/agent.json`) + fail-closed `/a2a/` handler over row 39's `authority.Resolver` and `protocol`'s A2A helpers; injected `writeAPIError` (daemon) + resolver + `GetRegistryHead` + registry reader |
| `host/projection/projection_test.go` | P6.B-A2A-CARD | Exact-set, 4-kind denial, absent/read-head, JSON-RPC codes, constant not-available, no-`X-World-Session`, no-`CallerSurface`, bounded-wait, wire-form tests |
| `host/daemon/daemon.go` | P6.B-A2A-CARD | Additive mounting of the two World-owned routes; inject the existing resolver + writes seam into the projection config at mount |
| `host/daemon/daemon_test.go` | P6.D + card | P6.D allowlist narrowness line + test; REST single-writer/route regression half |
| `go.mod` / `go.sum` | P6.D | `+github.com/sunholo-data/ailang v0.33.2`; `go-isatty`→0.0.22, `x/sys`→0.47.0 |

## Milestone P6.A-CTX — bounded, context-aware session resolution in `host/authority` (~0.1d; PREREQUISITE, FIRST)

SEQUENCED FIRST, before P6.D / P6.B-A2A-CARD. Adds one `host/authority` entry point that keeps row
39's credential policy byte-for-byte (R1-CTX) so both projection routes resolve within a bounded
deadline (Decision 6) without waiting unboundedly on the store's single connection and without
misclassifying a deadline as an unknown credential.
**Surface** — add to the `Resolver` interface alongside `Resolve`:
`func (r *resolver) ResolveContext(ctx context.Context, header string, now int64) (ResolveOutcome,
error)`. Header parsing, hashing, the indexed PK lookup, expiry and grant decoding stay IDENTICAL to
`Resolve`. Exactly two differences: (1) the store call receives `ctx` (`ResolveSession(ctx,id)`
uses `QueryRowContext` — database/sql honours it while waiting for the single pooled connection);
(2) a store error (incl. `context.Canceled`/`DeadlineExceeded`) is returned as a non-nil `error`,
NOT a denial, so the caller answers 503/504 (or A2A -32603) instead of a false 401.
**Unchanged / residual:** `Resolve(header, now)` keeps today's behaviour — delegates to
`ResolveContext(context.Background(), …)` and maps a non-nil error to `DenialUnknown` — so
`/v1/commit`'s `SessionMiddleware` is byte-identical and STILL resolves without a request deadline
(declared residual; this doc does not change it). Both projection routes call `ResolveContext` with
the ONE bounded context from Decision 6.
**AC (P6.A-CTX):** policy unchanged (every existing `host/authority` test still passes); a deadline
test HOLDS the store's single connection (e.g. an open transaction in the test), calls
`ResolveContext` with a short deadline, and asserts an error wrapping `context.DeadlineExceeded`
within the bound and no goroutine left running (resolution runs on the calling goroutine; assert at
source level: no `go ` statement in the resolve path).
**Mutations (RED-gate the ACs / policy-equivalence):** `MUT-RESOLVE-BACKGROUND` (passes
`context.Background()` → deadline test REDs); `MUT-CTX-ERR-AS-UNKNOWN` (maps the store error to
`DenialUnknown` → non-nil-error test REDs); `MUT-RESOLVE-POLICY-DRIFT` (change one policy branch
only, e.g. expiry `>`→`>=` → an equivalence test running `Resolve` and `ResolveContext` over the
same fixture table REDs).
**LOC / estimate:** ~40-60 LOC in `host/authority/resolver.go` (interface + method, no new file) +
~10 LOC of tests; +~0.1d, in the top-line bump.

## Milestone P6.D — Dependency admission, ATOMIC WITH THIS DOC'S FIRST REAL CONSUMER (~0.15d)

**Why here.** Quorum round 5 rejected pre-landing this dependency behind a dead anchor ("speculative
core growth"); its prescribed fix — *"Move the pinned dependency and narrow allowlist change into
whichever child first becomes unblocked, where a real handler import provides the compile-visible
use"* — is what puts it here. This doc's `host/projection` A2A handlers are that real consumer.

**Ordering.** The sibling MCP child is **still blocked** (F11: `ailang#885` closed without the
seam shipping, protocol unchanged at v0.33.2). So P6.D ships HERE. The predicate is:
`git grep -n 'serveapi/protocol' -- go.mod host/daemon/daemon_test.go` → **rc=1 at this revision
(F10) → P6.D lands here**; if the sibling ever unblocks first it inherits an already-admitted
dependency and this milestone records a skip.

**Spec** (substance carried; numbers corrected iter-187 per F10): `go get
github.com/sunholo-data/ailang@v0.33.2`, then add **exactly one** `allowedDepModules` entry: the
**PACKAGE path** `github.com/sunholo-data/ailang/serveapi/protocol` — never the module root. The
matcher (`disallowedDeps`, `daemon_test.go:790-812` — re-read this session) treats entries as path
prefixes, so the package-path entry admits exactly that package (+ future subpackages) while the
module root would admit `internal/apiserver`'s measured 476 disallowed packages. Update the
`allowedDepModules` doc comment ("module roots") to "module roots or package paths".

- **P6.T already green at this revision** (`go.mod` = `go 1.26.6`) so the go-version floor is met;
  admission is **250 → 251** (closure over `./host/daemon/... ./cmd/ailang-worldd/...`, measured
  **250 today**, F10; the doc's stale blocked-era "249→250" is superseded).
- Include a **narrowness test**: with the new entry in place,
  `disallowedDeps(["github.com/sunholo-data/ailang/internal/apiserver"])` (+ a representative
  cloud path) is non-empty — proving the entry did not widen the gate.
- The compile-visible use is the real handler import; `host/daemon/protocol_use.go` (the parent's
  one-iteration-old dead anchor) is **withdrawn and is not created by any doc**.
- Named risk (AC16): the pin also bumps `go-isatty → v0.0.22` and `x/sys → v0.47.0`
  (already-allowlisted roots) and `go.sum` moves; no other graph movement is permitted.
- Merge criterion: both CI jobs green; `TestDaemonDependencyAllowlist` green with the probe
  demonstration (allowlist minus the new line → REDs naming exactly one intruder).

Files: `go.mod`, `go.sum`, `host/daemon/daemon_test.go` (ONE package-path line + narrowness test).

- [ ] **AC16 — pinned dependency, narrow admission (P6.D):** `go.mod` requires
  `github.com/sunholo-data/ailang v0.33.2`; `allowedDepModules` gains EXACTLY ONE entry, the
  package path `github.com/sunholo-data/ailang/serveapi/protocol`; the closure over both gated
  patterns moves **250 → 251** (ADD 1, removed set empty); the narrowness test proves
  `internal/apiserver` (+ a cloud path) is still refused; the ONLY other graph movement is the two
  measured bumps. Any third movement fails the row.

| Gate | Named RED mutation (concrete edit) | Required red observation |
|---|---|---|
| AC16 `MUT-ALLOWLIST-ROOT` | replace the package-path entry with the module root `github.com/sunholo-data/ailang` | the narrowness test REDs: `internal/apiserver` is no longer refused by `disallowedDeps` |
| AC16 `MUT-FACADE-IMPORT` | import `github.com/sunholo-data/ailang/serveapi` (the facade) at `host/projection`'s protocol-import site | `TestDaemonDependencyAllowlist` REDs naming the intruding packages by name (scale 476 across 86 roots) |

## Milestone P6.B-A2A-CARD — Session-Authority-backed A2A agent card + fail-closed `/a2a/` (~0.4d)

**Starts when:** row 39 landed (DONE — `a036062`), P6.T landed (DONE — go1.26.6), and P6.D above
has landed (here or in the sibling, whichever unblocked first; at this revision P6.D lands here,
F10). Built **exclusively** on `protocol` types/helpers plus row 39's `authority.Resolver` over
the landed `transitionreg` interfaces. **It does not require a server-side invocation coordinator
(F7) because it never invokes a transition.**

- `GET /.well-known/agent.json`: resolve via the injected `authority.Resolver.Resolve(header, now)`
  (B1, B2); deny → the F2 status mapping with the daemon's REST APIError envelope (reused
  `writeAPIError`, injected at mount); success → ONE `transitionreg.NewRequest` request-set,
  `Request.Allowed()` → card skills **verbatim** (F6), `map[string]any` literal with upstream keys
  (F5), `name`/`description`/`version` from `protocol.AgentInfo`. Absent head → zero skills
  (200); read error → 5xx (B3).
- `POST /a2a/`: resolve exactly as the card route; parse `protocol.A2ARequest`; JSON-RPC codes
  -32600/-32601/-32602/-32603 as per Decision 3 (B4); authorized skill → constant not-available
  error via `protocol.A2AError`; never success, never REST, never a store write.
- **Mount** (B5): both routes mounted ADDITIVELY in `host/daemon`; no REST route renamed/removed/
  repurposed (P6 frozen); they are NOT in `d.isProtected` — the projection handler does its own
  resolution so the `/a2a/` route can use A2A-form errors.
- Bounded wait: ONE request context through resolution + snapshot reads (Decision 6, descoped).

**Cards are diffable**: skills sorted in registry bytewise order (F9 `Request.Allowed()` order);
two sessions with unequal caps get unequal exact skill sets.

- Merge criterion: both CI jobs green, no skips (zero-skip clause from parent AC10 travels here
  with the tests it governs), every acceptance mutation demonstrated RED then reverted GREEN.

## Acceptance Criteria (parent IDs RETAINED; **AMENDED iter-187** where the contract changed;
**MOVED → Deferred** items marked; added rows tagged NEW-iter187)

- [ ] **AC1 — wire-contract reuse (AMENDED iter-187 for F6):** World's **`/a2a/`** route constructs
  and parses wire forms EXCLUSIVELY through the pinned `serveapi/protocol` types/helpers
  (`A2ARequest`/`A2AError`/`A2AResult`, `AuthorizationError`); the **card** is built as a
  `map[string]any` literal with the SAME keys upstream uses (F5) with `name`/`description`/
  `version` from a `protocol.AgentInfo` value. Production World code declares NO parallel wire
  struct and hand-formats NO `A2ARequest`/`A2AError`/`A2AResult` bytes. The card ID is emitted
  VERBATIM and is **NOT** routed through `protocol.CallerSurface` or `ValidateMCPName` (F6).
  (The daemon's own REST `APIError` envelope may be reused for the card-route denial body via the
  injected `writeAPIError`, but that envelope is daemon-owned, never reproduced in
  `host/projection`.)
- [ ] **AC2 — exact surface (AMENDED iter-187 — discovery only):** for two non-empty sessions
  with unequal capability sets, the A2A `skills[].id` set equals each session's authorized
  transition-ID set exactly, in registry order, emitted VERBATIM (F6); no extras.
- [ ] **AC3 — card tracks the session:** a session's A2A skill-ID set changes when its capability
  snapshot changes; two concurrent sessions with different snapshots observe different cards from
  the same daemon.
- [ ] **AC4 — ambient exports absent:** `exit`, `writeBytes`, all eight observed `std/io` exports,
  and `submit_feedback` are absent from the card and skill set.
- [ ] **AC5 — invocation enforcement (AMENDED iter-187 → `/a2a/` admission):** an unauthorized,
  stale, or guessed `skill_id` is rejected with **-32602** before any other processing and
  produces no store/log change; an authorized `skill_id` receives the constant not-available
  error (**-32603**), never a success.
- [ ] **AC6 — session failure, carrier fixed (D-WORLD-26 = A):** the resolver reads the session
  credential from the standard `Authorization: Bearer` header ONLY; absent, malformed, unknown,
  or expired credentials fail closed on both routes, encoded per route as the **Denial matrix**
  specifies (r2 carve-out, glm/gemini verbatim): the CARD route uses the F2 HTTP status mapping
  (401/400) with the APIError envelope; the `/a2a/` route uses JSON-RPC error codes at HTTP 200
  via `protocol.A2AError` — -32001 for absent/unknown/expired, -32600 for malformed; no
  default/global capability set is used; the static serve-api key
  is never accepted as a session credential (constraint i); no alternative header (including
  rejected Arm-B `X-World-Session`) is read, even as a fallback.
- [ ] **AC7 — one-snapshot consistency:** a request never mixes registry/capability epochs;
  concurrent registry/session change is observed only on a subsequent request.
- [ ] **AC8 — A2A protocol conformance (AMENDED iter-187 — discovery + `/a2a/` admission):** the
  card route emits the upstream `map[string]any` keyset (F5) with the exact skill-ID set; the
  `/a2a/` route emits ONLY `protocol.A2AError` bodies (always HTTP 200, F5b) with the JSON-RPC
  codes -32001 (session denial, per the Denial matrix), -32600, -32601, -32602 and -32603; no
  success path exists for invocation. *(r2 carve-out: the four-code restriction is removed.)*
- [ ] **AC9 — landed behavior preserved:** REST v1 route/body regressions and cross-process
  writer-lock tests remain green; the two new routes are NOT in `d.isProtected` (adding them to
  the predicate is a mutation that must RED); projection never opens a store.
- [ ] **AC11 — dependency floor:** the printed transitive graph contains no disallowed package and
  `TestDaemonDependencyAllowlist` covers the new package/cmd paths.
- [ ] **AC12 — dynamic source:** changing the transition-registry head changes the next authorized
  card without restart or `.ail` file edits.
- [ ] **AC13 — bounded projection waits → MOVED → Deferred to invocation.** The bounded ONE-context
  requirement for resolution + snapshot reads on the two new routes is retained in the milestone
  (Decision 6) with its own MUT (`MUT-DROP-DEADLINE-PROJ`); the commit-boundary and socket-closure
  halves are invocation-only and moved verbatim.
- [ ] **AC14 — no deadline tampering:** no code in `host/projection` (or in this doc's diff) calls
  `ResponseController.SetWriteDeadline`/`SetReadDeadline` or otherwise relaxes the frozen D7
  deadlines — asserted at source level — and REST `/v1/*` still uses the unchanged D7 values.
  (The route-local SSE relaxation this row forbids is the MCP dispatch child's.)
- [ ] **AC-CARD-DENIALS (NEW-iter187):** on the card route, each of the four denial kinds
  (absent/unknown/expired and malformed) yields its F2 HTTP status (401/401/401 and 400) and the
  daemon's REST APIError envelope `class` (`SessionAbsent`/`SessionUnknown`/`SessionExpired` and
  `InvalidSession`), over the injected `writeAPIError`; each message is a constant per kind.
- [ ] **AC-ABSENT-HEAD (NEW-iter187):** absent registry head → card with **ZERO skills** at HTTP
  200 (still authenticated); a real store/read error (or a later `NewRequest` failure) → **5xx
  fail-closed**. The two are told apart by checking `GetRegistryHead` first (B3). A test covers
  both.
- [ ] **AC-A2A-CODES (NEW-iter187; r2: -32001 added):** `/a2a/` returns `-32001` for an
  absent/unknown/expired session and `-32600` for a malformed one (Denial matrix), `-32600` for
  non-`"2.0"`, `-32601` for
  method ≠ `tasks/send`, `-32602` `"not authorized"` for unlisted/guessed/stale `skill_id`, and
  a constant-message `-32603` for an authorized `skill_id`; every body is `protocol.A2AError`,
  never REST `APIError`, never success, never a store write.
- [ ] **AC-DENIAL-JSONRPC (r2 carve-out, gpt6-astra + oc-glm-5-2 verbatim):** all four denial
  kinds are tested on BOTH routes. On `/a2a/`: HTTP 200 carrier, `protocol.A2AError` body, -32001
  for absent/unknown/expired and -32600 for malformed, constant messages, no REST envelope. On the
  card route: HTTP 401/401/401/400 with the APIError envelope. For every denied request on either
  route, the test asserts that no registry snapshot was acquired and skill admission was never
  reached (e.g. a counting registry reader / `GetRegistryHead` observer reads **0** calls).

## Non-Vacuity — Named RED Mutation for Every Gate

| Gate | Named RED mutation (concrete edit) | Required red observation |
|---|---|---|
| AC1 `MUT-PROTO-OWNER` | declare a local parallel A2A wire struct in `host/projection` and hand-format an `A2AError`/`A2ARequest` body | wire-ownership source test REDs on the World-declared wire struct/envelope bytes |
| AC1 `MUT-CARDSURFACE` | build the card by routing skill IDs through `protocol.CallerSurface` (the F6 path) | `tools.echo`-shaped / `world/recovery-transition/v1`-shaped real ID makes `CallerSurface` error → card test REDs; verbatim-ID assertion REDs |
| AC2 `MUT-SESSION-UNION` | return the union of all transitions for every session | low-capability session exact-set test REDs |
| AC3 `MUT-CARD-GLOBAL` | generate the card from the unfiltered registry | per-session card set-equality/session-change test REDs |
| AC4 `MUT-UNFILTERED-PROJECTION` | re-enable raw v0.30.0 projection | surface test REDs on `std.io.writeBytes`/`writeBytes`, `exit`, or `submit_feedback` |
| AC5 `MUT-A2A-STALE` | admit a `skill_id` not in the session's `Allowed` set (drop the -32602 check) | guessed-name `/a2a/` call gets success/store change and the -32602 test REDs |
| AC5 `MUT-A2A-FAKE-SUCCESS` | emit `A2AResult` fake success for an authorized `skill_id` instead of the constant -32603 error | the constant not-available test REDs (must not be success) |
| AC6 `MUT-DEFAULT-CAPS` | map an unknown session to a process default | unknown-session card/call test REDs |
| AC6 `MUT-ALT-HEADER` | make the resolver fall back to reading `X-World-Session` when `Authorization` is absent | carrier test REDs: a request bearing ONLY the rejected Arm-B header must fail closed, and under the mutation it resolves |
| AC6 `MUT-KEY-AS-SESSION` | accept the static serve-api API-key value as a `Bearer` session credential | constraint-(i) test REDs: the key resolves to a session under the mutation while the test demands fail-closed |
| AC7 `MUT-SPLIT-SNAPSHOT` | re-read capabilities after reading registry within one request | barrier-controlled epoch-consistency test REDs |
| AC8 `MUT-A2A-CODEFLIP` | swap one `/a2a/` code mapping (e.g. `-32602` ⇄ `-32603`) | the JSON-RPC-code test REDs on the swapped branch |
| AC9 `MUT-SECOND-OPEN` | make projection call `store.Open` during mount | writer-lock/live-daemon test REDs with `WriterAlreadyActive` |
| AC9 `MUT-ISPROTECTED-EXPAND` | add `/.well-known/agent.json` (or `/a2a/`) to `d.isProtected` | REST route regression / unauthenticated-health REDs: a public read route now demands a session |
| AC14 `MUT-DEADLINE-RELAX` | add a `ResponseController.SetWriteDeadline(time.Time{})` call to a projection handler | the source-level assertion REDs naming the call site; the REST deadline regression is the same-run control that D7 constants did not move |
| AC-CARD-DENIALS `MUT-CARD-DENIAL-TABLE` | return the wrong HTTP status for one denial kind (e.g. malformed → 401, or all → 500) | card-route denial test REDs on the kind whose mapping the mutation broke (all four kinds asserted) |
| AC-CARD-DENIALS `MUT-CARD-ENVELOPE` | emit a plain JSON body (not the daemon's APIError envelope) for a card-route denial | APIError-envelope assertion REDs on the card route |
| AC-ABSENT-HEAD `MUT-ABSENT-HEAD` | treat absent head (`ok==false`) as a real read error → 5xx | the zero-skills-at-200 test REDs; the read-error-5xx half stays green (control) |
| AC-A2A-CODES `MUT-A2A-MESSAGE-INTERP` | interpolate the `skill_id` (or request content) into the not-available message instead of a constant | the constant-message assertion REDs on `/a2a/` |
| AC-DENIAL-JSONRPC `MUT-DENIAL-401` | *(r2 carve-out, oc-glm-5-2 verbatim)* emit HTTP 401 instead of HTTP 200 + JSON-RPC error for a session denial on `/a2a/` | the new AC REDs on status code |
| AC-DENIAL-JSONRPC `MUT-DENIAL-SNAPSHOT-FIRST` | move the registry-snapshot acquisition (or `GetRegistryHead`) ahead of session resolution on either route | the zero-snapshot-on-denial assertion REDs (counting reader observes ≥1 call on a denied request) |
| AC13/Decision-6 (proj) `MUT-DROP-DEADLINE-PROJ` | replace the propagated request context with `context.Background()` before session resolution or snapshot reads | bounded-wait fault-injection test exceeds the configured bound and REDs |
| zero-skip `MUT-SKIP-SOCKET` | add `t.Skip` on listen failure | zero-skip CI assertion/source check REDs |
| AC11 `MUT-CLOUD-DEP` | add `cloud.google.com/go/storage` to projection imports | `TestDaemonDependencyAllowlist` reports it by name |
| AC12 `MUT-STARTUP-CACHE` | cache transition descriptors once at daemon startup | registry-head-change card returns stale set and REDs |
| AC16 `MUT-ALLOWLIST-ROOT` | (P6.D) package-path → module-root entry | narrowness test REDs: `internal/apiserver` no longer refused (see P6.D table) |
| AC16 `MUT-FACADE-IMPORT` | (P6.D) import the `serveapi` facade at the projection protocol-import site | `TestDaemonDependencyAllowlist` REDs naming intruders (see P6.D table) |

## Conflict Surface

- **SessionMiddleware vs the two new routes.** The middleware enforces only `d.isProtected`
  (today `POST /v1/commit`); the two new routes are NOT protected — they resolve themselves (B5)
  so the `/a2a/` route can emit A2A-form errors. A card fetch with a bad credential therefore gets
  the projection's own denial (F2 mapping + APIError envelope), distinct from the middleware's
  handling of `/v1/commit`, but byte-compatible in status and class.
- **Card-route denial envelope vs P6's "no APIError translation".** P6 freezes the `/v1/` route
  table/bytes; it does not forbid reusing the shipped `writeAPIError` envelope on a NEW non-REST
  endpoint. The card fetch is plain JSON (not JSON-RPC), so `protocol.A2AError` (a JSON-RPC
  response) is the wrong shape there; reusing the daemon's envelope avoids a second wire shape for
  the identical denial semantics and involves no parallel struct in `host/projection` (the daemon
  injects the writer). The `/a2a/` route is the JSON-RPC surface and emits `A2AError` — never
  `APIError`.
- **Row 39's resolver vs existing broker/session/daemon machinery** (the round-4 `catch`'s
  overlap analysis). `bind.go`/`broker.NewSession` (`broker.go:87`, F9) is a pure constructor of
  sessions from resolver output; `host/broker/credential.go` is OUTBOUND-only registry-publish
  credentials (a grep re-check this session: `RegistryCredentialProvider`/`assertNoAmbientRegistry`);
  `host/evidence/`'s `Authenticate*` is evidence-envelope signing; `host/authority` owns the
  INBOUND resolver and is the single credential path. The projection adds no policy engine, no
  store, and no credential store; it consumes the daemon's one resolver instance (F3).
- **worldd single-writer flock.** Projection uses the daemon's handle; no `store.Open`, sidecar
  writer, or lock change. Existing `WriterAlreadyActive` tests remain green.
- **Frozen REST v1 route table.** health, head, world, object, log-entry, log-range, registry
  wildcard, commit all unchanged; `/.well-known/agent.json` and `/a2a/` are additive.
- **`host/registry` vs `host/transitionreg`.** The interpreter epoch registry (`host/registry`,
  `world/epoch-registry/v1`) and the transition registry (`host/transitionreg`) are different
  landed subsystems. The projection reads only `transitionreg` snapshots and neither renames nor
  overloads either.
- **Absent registry head vs read failure.** Distinct via `GetRegistryHead` first (B3, AC-ABSENT-
  HEAD); absent → zero skills at 200, read error → 5xx.
- **NAME grammar (F6).** Card IDs are emitted verbatim and never pass through `CallerSurface`;
  the MCP dispatch child inherits the `^[a-zA-Z0-9_-]{1,64}$` (protocol) vs `validateID` (World)
  conflict and must solve naming itself.
- **Performance baseline.** No existing benchmark removed; P6.B-A2A-CARD adds bounded round-trip
  benchmarks only if measurement shows a material new hot path; no rewrite of `bench/BASELINE.md`.

## Axiom Compliance (re-scored iter-187 for the descoped shape)

| Axiom | Score | Justification |
|---|---:|---|
| A1 Determinism | +1 | stable-ID verbatim ordering (registry bytewise) and one snapshot per request |
| A2 Replayability | 0 | descope removes invocation → the read-only card//a2a/ never records a transaction (was +1 under invocation) |
| A3 Effect Legibility | +1 | protocol I/O stays at the host boundary; card is a read-only projection with zero effects |
| A4 Explicit Authority | +2 | exact per-session card predicate + `/a2a/` admission over the landed row-39 resolver; carrier fixed by D-WORLD-26; fail-closed with no API-key/alternate-header path |
| A5 Bounded Verification | +1 | the (deferred) verify pipeline is not reached; admission reuses landed `broker.Allows` |
| A6 Safe Concurrency | +1 | daemon retains sole handle; immutable request snapshot |
| A7 Machines First | +2 | upstream A2A card keyset + structured JSON-RPC errors (-32600/-32601/-32602/-32603) |
| A8 Minimal Syntax | 0 | no language syntax change |
| A9 Cost Visibility | 0 | no new budget claim; invocation (the expensive half) is deferred |
| A10 Composability | +2 | standard A2A wire; the MCP≡A2A equality obligation is preserved in the MCP dispatch child (with the F6 grammar caveat noted) |
| A11 Structured Failure | +2 | absent/unknown/expired/malformed (401/401/401/400), absent-head vs read-error (200-zero vs 5xx), JSON-RPC codes on `/a2a/` — each distinct, all over the landed resolver (was +1 under invocation's unbuilt coordinator) |
| A12 System Boundary | +2 | projection is an adapter, never kernel state; no store write path |

**Net: +15; hard axioms A1/A3/A4/A7 non-negative.**

## Deferred to invocation (new queue row, gated on a coordinator that does not exist)

Everything moved verbatim from the iter-187 descope. The two prerequisites it is gated on:
**(i) a server-side transition-invocation coordinator** (F7 — none exists), and **(ii) a
production path that publishes the transition registry** (F8 — `StoreReader.Publish`/`BuildNext`
have zero production callers today). Until both exist, the following is NOT buildable and is
re-opened as a new queue row when F7 lands.

- **Decision 3 responsibility 4 (verbatim):** "dispatch an authorized invocation into the normal
  propose → verify → commit coordinator — implementing `protocol.Invoker`, with
  `transitionreg/bind.go`'s typed refusals (`TransitionAbsentError`, `AccessDeniedError`,
  `ProposalMismatchError`) mapped through `protocol.AuthorizationError`/`AuthorizationStatus`."
  (Prereq (i) names the "coordinator" this responsibility requires — it does not exist.)
- **Premise P4 (verbatim):** "**P4 — propose → verify → commit remains mandatory.** Protocol
  invocation enters the same broker/session/transaction path as every other agent action. The
  projection receives no direct `store.Commit`, `SetRegistryHead`, or REST `/v1/commit`
  authority." (Its dispatch half is deferred; its "no store authority" half is satisfied by the
  read-only card shape.)
- **`protocol.Invoker`** implementation and the whole invocation success path (no A2A success for
  tasks/send in this daemon until it lands).
- **AC13 — bounded projection waits (verbatim), invocation half:** "with the configured finite
  server bound, a barrier-blocked session resolver, registry snapshot, authorization provider,
  broker, proposer, verifier, or response write, and a disconnected client each terminate within
  that bound… Cancellation is tested on BOTH SIDES of the commit boundary… cancelling immediately
  BEFORE the coordinator accepts the commit asserts **no durable mutation**; cancelling
  immediately AFTER acceptance asserts **exactly one recoverable, queryable/replayable receipt**…
  The fault-injection test must also assert OS-LEVEL socket closure… a disconnected or
  deadline-expired transport must be OBSERVED closed (`http.Server.ConnState` tracking or a
  client read error), never merely context-cancelled… A logical Go `context` cancellation alone
  does not satisfy this criterion."
- **The commit-boundary / AC13 / socket-closure obligation** (the `MUT-COMMIT-BOUNDARY-LIE`,
  `MUT-LEAK-CONN`, `MUT-SKIP-SOCKET` mutations) — `MUT-COMMIT-BOUNDARY-LIE` and `MUT-LEAK-CONN`
  apply invocation-only and are re-activated when F7 lands; the A2A route reuses the D7 server
  deadlines with no SSE (a disconnected A2A client is handled by the D7 write deadline, which is
  already observed — the OS-level socket-closure *obligation as a test contract* is deferred).
- **AC5's invocation-dispatch half** (dispatch through propose → verify → commit) and **AC8's
  `A2AResult`/invocation-result leg** (emitting a task result for an accepted call).

## Estimate honesty

**~0.55d** total (card P6.B-A2A-CARD ~0.4d + P6.D ~0.15d) — descoped from the inherited ~0.7d
invocation-bearing P6.B-A2A line. This **starts now** (row 39 landed, P6.T landed); it does not
wait on a coordinator because it never invokes. The deferred invocation half is re-priced when
F7/F8 land (the card and `/a2a/` admission surface already lay the foundation).

## Quorum status

**This doc is NOT quorum-cleared.** It was authored at split #2 (iteration 126), blocked, and is
now **revised iter-187** against row 39's landed contract as its own blocking predicate demanded.
At pick time the doc MUST go through the full design quorum (`ailang design-quorum`,
reject-by-default synthesis). Nothing in this doc pre-authorizes a sprint.

#### Quorum round 1 (iter-187 r1) — three REJECTions, Measured and Answered

Round 1: **gpt6-astra / gemini-3-1-pro / oc-glm-5-2 all REJECT** (oc re-run alone after malformed first reply); the controller MEASURED every objection (rows **R1-CTX / R1-HEAD / R1-WIRE** + **F5b**, VERIFIED BY CONTROLLER iter-187 r1). Each answered here:
1. **gpt6-astra — bounded wait cannot reach resolution.** Measured: no ctx (resolver.go:129), one connection (store.go:305), naive fix→false 401 (resolver.go:130-134). **Answered:** prerequisite **P6.A-CTX** (FIRST) — `ResolveContext(ctx,…)` returns non-nil `error` on store failure, never a denial; both routes use one bounded ctx; `Resolve`/`/v1/commit` unchanged (residual). Card 503/504; `/a2a/` -32603 constant.
2. **gemini-3-1-pro — GetRegistryHead never verified.** Measured: store.go:636-648 (`sql.ErrNoRows`→`(zero,false,nil)`, scan errors wrapped, bad ref→error), daemon.go:331-337 reads seam, store.go:88 / transitionreg.go:74 same name. **Answered:** verification row added; B3 + AC-ABSENT-HEAD cite the lines; head race in B3 (absent-check + `NewRequest` success → use its result; head + fail → 5xx); test covers "check says absent".
3. **oc-glm-5-2 — A2A wire types never verified.** Measured (gh api @ v0.33.2): struct + `A2AError`/`A2AResult` sigs; `A2AError` ALWAYS HTTP 200; handler reads `params.Metadata["skill_id"]`. **Answered:** row **F5b** added; /a2a/ DENIAL = JSON-RPC error body at HTTP 200 — **-32001** absent/unknown/expired, **-32600** malformed. AC `AC-DENIAL-JSONRPC`; MUT `MUT-DENIAL-401`.

#### Quorum round 2 (iter-187 r2) — BLOCKED 3/3 on ONE surface; closed under the narrow-refinement carve-out

Round 2 (`absent_reviewers` empty): **gpt6-astra, gemini-3-1-pro, oc-glm-5-2 all REJECT, all on the same surface**. The r1 answer promised `AC-DENIAL-JSONRPC` and `MUT-DENIAL-401`, but neither reached the AC or mutation tables, so AC6 and Decision 3 still demanded HTTP 401/400 on `/a2a/`, where `protocol.A2AError` always writes HTTP 200 (F5b). The controller confirmed the gap by grep (each ID occurred once, only in the r1 history prose). Every objection carried a concrete reviewer-authored `proposed_fix` and none disputed the design direction, so the controller applied the fixes **verbatim** (no reviewer overridden, nothing controller-invented): the authoritative **Denial matrix** in Decision 3 (astra), the route-split wording of Decision 3 responsibility 1 (gemini) and AC6 (glm/gemini), AC8's code set widened to include -32001 (astra), -32001 added to AC-A2A-CODES (glm), plus **AC-DENIAL-JSONRPC** and **MUT-DENIAL-401** (astra/glm). Two sentences are the **controller's own consistency edits**, labelled as such: Decision 3 responsibility 1 now names `ResolveContext` (the r1 P6.A-CTX fix had left the old `Resolve` call in that sentence), and `MUT-DENIAL-SNAPSHOT-FIRST` gives astra's "denied requests never acquire registry snapshots" clause a named mutation. Surface count across rounds: r1 = three surfaces (resolution deadline, head seam, wire types); r2 = one (denial encoding, introduced by r1's own answer). No split is owed.

## Relationship to the parent and to charter clause 6

Charter clause 6 (F12) is partitioned across THREE docs:

- the parent, [`w-mcp-projection.md`](w-mcp-projection.md) — the enabling milestones P6.T/P6.V
  (P6.D is carried by whichever child first unblocks);
- child #1, [`w-mcp-dispatch-projection.md`](w-mcp-dispatch-projection.md) — the MCP dispatch
  half, **still BLOCKED** (F11: `ailang#885` closed without the seam shipping); carries the
  MCP≡A2A equality obligation and inherits the F6 name-grammar conflict;
- child #2, THIS DOC — the A2A card + fail-closed `/a2a/` half, now unblocked on row 39 and
  descoped to what is buildable now; invocation deferred to a new queue row gated on F7/F8.

## Related Documents

- [w-mcp-projection.md](w-mcp-projection.md) — the SPLIT parent: P6.T/P6.V enabler milestones,
  executable now; its quorum log records both split dispositions
- [w-mcp-dispatch-projection.md](w-mcp-dispatch-projection.md) — split child #1: the MCP dispatch
  half, blocked on `ailang#885`; carries the cross-surface MCP≡A2A equality obligation and the F6
  name-grammar problem
- [world-mission.md](../world-mission.md) — clause 6 (`:1001`); **queue row 39
  `w-session-authority`** (LANDED `a036062`, iter-181, the resolver this doc consumes); D-WORLD-5,
  D-WORLD-26 (the carrier ruling binding AC6)
- [coding-standards.md](../coding-standards.md) — S1–S6
- [DESIGN.md](../DESIGN.md) — §3.7 protocol-native boundary, §14, §17
- [w-worldd-m2.md](../implemented/w-worldd-m2.md) — shipped daemon, REST freeze, D7 constants
- [w-effect-broker-m3.md](../implemented/w-effect-broker-m3.md) — landed session/capability/broker
  (the in-process model row 39 extended to the HTTP-facing boundary)
- [w-store-durability.md](../implemented/w-store-durability.md) — the commit-boundary Go surface
  whose proof is consumed only by the deferred invocation half
