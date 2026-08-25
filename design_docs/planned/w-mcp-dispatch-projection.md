# w-mcp-dispatch-projection — Session-Scoped MCP Projection (SPLIT child of `w-mcp-projection`)

**Status**: Planned — **BLOCKED UPSTREAM on
[`sunholo-data/ailang#885`](https://github.com/sunholo-data/ailang/issues/885); NOT
QUORUM-CLEARED; deliberately NOT sprint-ready.** This doc exists to carry a confirmed
quorum objection and its blocking predicate honestly, not to specify work that cannot yet be
designed: the shape of the fix depends on what `#885` delivers, and a full spec written now
would be designed against an unknown seam.  
**Item**: split child of `w-mcp-projection` (charter clause 6, queue row 5)  
**Clause**: clause-6 — the *"project the transition registry over MCP"* half  
**Estimate**: deliberately NONE — estimating against an undelivered seam would be fiction; the
estimate is written when the design is (at unblock + pick time)  
**Author**: rotation designer, split round, 2026-08-25 (iteration 125)  
**Verified against**: upstream `github.com/sunholo-data/ailang` tag `v0.33.2` (`63e7909f`) via the
GitHub API (no local upstream checkout); World `dev` at `2e44e3e`  
**Date**: 2026-08-25

## Why this doc exists — the objection, verbatim

The parent, `design_docs/planned/w-mcp-projection.md`, was blocked at quorum round 3
(2026-08-25, both reviewers present, `absent_reviewers` empty) by a DIRECTION-level objection
from `gpt5-6-sol`, **confirmed first-party by the controller with firing controls**. The
disposition was SPLIT — a controller routing call, explicitly not `needs-human-review`: the
reviewer-clean A2A remainder stayed in the parent; this doc carries the objection and the scope
it blocks. The objection, verbatim:

> "The selected upstream seam is insufficient for the design's own no-codec rule … no MCP request
> parser or handler for JSON-RPC method dispatch, initialization, `tools/list`, or `tools/call`. A
> World-authored `/mcp/` handler therefore appears forced to implement MCP/JSON-RPC parsing and
> dispatch locally, directly contradicting P1, the Design Freeze, and AC1."
>
> — `gpt5-6-sol`, quorum round 3 on `w-mcp-projection`, 2026-08-25 (artifact
> `.ailang/state/mission-quorum/w-mcp-projection-2026-08-25T16-33-38Z.json`)

The objection is correct, and it is not a defect in the parent's measurements — it is a scope
finding: `serveapi/protocol@v0.33.2` delivered the whole A2A wire surface and the MCP *envelope*
helpers, but no MCP JSON-RPC dispatch. The MCP half of clause 6 therefore has exactly two local
routes, and both are closed by this repo's own guardrails.

## The measured evidence — both routes are closed

Rows marked **[R]** were re-derived first-party in the split session (2026-08-25, commands as
shown); rows marked **[I]** are inherited from the controller's iteration-124 measurements
(charter STATUS stamp + the parent's premise rows), labelled as such.

| # | Claim | Command / evidence | Result |
|---|---|---|---|
| E1 [R] | `protocol` ships MCP envelope helpers but the MCP HANDLER lives outside it and delegates dispatch to the SDK | `gh api 'repos/sunholo-data/ailang/contents/serveapi/mcp_handler.go?ref=v0.33.2'`, decoded: **187 lines**, `grep -c modelcontextprotocol` → **1**; same-call control `a2a_handler.go`: **180 lines**, SDK import count **0**, import block = `context`, `encoding/json`, `fmt`, `log`, `net/http` + `serveapi/protocol` | Confirmed — the A2A handler is the existence proof of the SDK-free shape; the MCP handler is not |
| E2 [R] | `mcp_handler.go` does not string-match MCP method names AT ALL | `grep -c 'tools/list\|jsonrpc'` in `mcp_handler.go` → **0**; same grep in `a2a_handler.go` (its own method strings) → **2** | Confirmed — the 0 is not a refutation of the reviewer but the STRONGER form of the claim: the entire parse/dispatch lives in `github.com/modelcontextprotocol/go-sdk`, so there is nothing dispatch-shaped in upstream's zero-cloud-importable surface |
| E3 [I] | Importing the MCP SDK moves the gated closure 249 → 283: **+34 packages across 5 new module roots**, and `TestDaemonDependencyAllowlist` reds on **28** disallowed packages including `golang.org/x/oauth2`, `go-sdk/auth`, `go-sdk/oauthex` | controller-measured iteration 124 over both gated patterns (`./host/daemon/... ./cmd/ailang-worldd/...`), with a sentinel control; recorded in the charter's iteration-124 STATUS stamp and the parent's status header | Inherited — an outbound-credential stack in the daemon core; breaches charter clause 2 AND clause 3 |
| E4 [I] | The `protocol` arm, by contrast, is clean: closure 249 → 250, the single added package IS `serveapi/protocol`, removed set EMPTY (sentinel control fired), stdlib-only across all four files | parent premise rows N3/N10/N11 (revision round, first-party with controls) | Inherited — the contrast that makes route (ii)'s cost a property of the SDK, not of importing upstream per se |
| E5 [R] | `#885` is OPEN with 0 comments | `gh issue view 885 --repo sunholo-data/ailang --json state,title,createdAt,comments` → `state=OPEN`, `comments=0`, created `2026-08-25T16:37:44Z`; same-call control `gh issue view 764 …` → `CLOSED`, 6 comments | Confirmed — the instrument can see closure and comments, and the ask is unanswered |
| E6 [R] | Upstream latest release is still `v0.33.2`; no `v0.34.*` tag exists | `gh api repos/sunholo-data/ailang/releases/latest --jq .tag_name` → `v0.33.2`; `gh api …/git/matching-refs/tags/v0.34` → empty (rc=0); control `…/tags/v0.33` → `v0.33.0`, `v0.33.1`, `v0.33.2` | Confirmed — nothing newer has shipped that could contain a dispatch seam |

## Why neither route can be taken

- **Route (i) — World writes its own MCP/JSON-RPC parsing and dispatch.** Forbidden by the
  parent's P1 and Design Freeze, its AC1, and `DESIGN.md` §3.7's protocol-native rule ("World
  invents no wire protocol where an open one exists" — and reimplementing an existing protocol's
  codec is the same reinvention). This is the exact contradiction the objection names.
- **Route (ii) — import `github.com/modelcontextprotocol/go-sdk`.** The measured closure (E3) puts
  an OAuth/outbound-credential stack (`golang.org/x/oauth2`, `go-sdk/auth`, `go-sdk/oauthex`)
  inside the daemon core: 28 allowlist violations, breaching charter clause 2 (zero-cloud) and
  clause 3 (no ambient authority). D-WORLD-5 (Mark, attended, 2026-08-17) prescribes the route
  out of a disallowed graph: **ask upstream for a narrow module, never a broad relaxation.** That
  default executed as `#885` — the same route that previously produced `#764` → `v0.33.2`.

## Blocking predicate — runnable, not prose

A future iteration decides "still blocked?" by RUNNING these, never by transcribing this doc.
Each probe carries a same-call control so an empty/negative read is a measurement, not a guess.

```bash
# 1. Is the upstream ask still open / unanswered?
gh issue view 885 --repo sunholo-data/ailang --json state,comments \
  --jq '{state: .state, comments: (.comments|length)}'
# measured 2026-08-25: {"state":"OPEN","comments":0}
# control (proves the instrument can see closure and comments):
gh issue view 764 --repo sunholo-data/ailang --json state,comments \
  --jq '{state: .state, comments: (.comments|length)}'
# measured 2026-08-25: {"state":"CLOSED","comments":6}

# 2. Has anything newer than v0.33.2 shipped?
gh api repos/sunholo-data/ailang/releases/latest --jq .tag_name
# measured 2026-08-25: v0.33.2
# control (proves the tag instrument sees tags at all):
gh api 'repos/sunholo-data/ailang/git/matching-refs/tags/v0.33' --jq '.[].ref'
# measured 2026-08-25: refs/tags/v0.33.0, v0.33.1, v0.33.2
```

**UNBLOCKED means ALL of:** `#885` is CLOSED as delivered (or a release tag newer than `v0.33.2`
exists whose diff shows a dispatch seam), AND the delivered seam re-measures as stdlib-only in
closure by the parent's own method (import-block read of every file in the seam package;
`go list -deps` delta over both gated patterns; the unmodified `TestDaemonDependencyAllowlist`
naming exactly the expected intruder(s)). A `#885` closed as WONTFIX also unblocks — the
DECISION, not the design: it routes back to D-WORLD-5's arms for a human ruling, never to a local
workaround (both local routes remain closed by the guardrails above).

## What upstream must deliver for this design to become writable

An SDK-free MCP JSON-RPC dispatch seam importable by zero-cloud consumers: parsing and method
dispatch for `initialize`, `notifications/initialized`, `tools/list`, and `tools/call` over MCP
Streamable HTTP, driven by the same callback interfaces `protocol` already ships
(`SessionResolver`, `ToolSource`, `Invoker`, `CallerSurface`), living in `serveapi/protocol` or
an equivalently stdlib-only package whose admission moves World's gated closure by ~+1 package —
not +34 across 5 module roots. `#885` asks for exactly this, citing `a2a_handler.go` (E1) as the
existence proof that upstream can ship handler-shaped code over `protocol` alone.

## Scope this doc carries (moved from the parent at the split)

- the `/mcp/` endpoint and the World-authored MCP HTTP handler;
- MCP JSON-RPC dispatch: initialization, `tools/list`, `tools/call`;
- dispatch-bound envelope framing (`WriteMCPEnvelope`/`RequestID` used from a served handler);
- SSE stream lifetime: the `/mcp/`-route-local `ResponseController` deadline relaxation, the
  finite stream-lifetime maximum, and OS-level closure on expiry (quorum r1/r2 obligations,
  `gemini-3-1-pro`);
- the cross-surface criterion (the parent's old AC3): per session, the MCP `tools/list` name set
  must exactly equal the A2A `skills[].id` set;
- the SSE-framing conformance half of the parent's old AC8 (`event: message` + `data:` JSON,
  proven against the frozen P6.A fixture);
- the parent's old AC14 (SSE/REST deadline separation — the relaxation is route-local to `/mcp/`,
  and the frozen D7 constants stay byte-unchanged); and
- the RED mutations `MUT-PLAIN-JSON`, `MUT-LEAK-SSE-CONN`, `MUT-SSE-REST-DEADLINE`.

These are **obligations, not a spec**: when `#885` delivers, the full design — decisions,
milestones, acceptance criteria with named RED mutations, premise verification log with firing
controls — must be written against the ACTUAL delivered seam, and every obligation above must
reappear in that draft's acceptance table. What is already known to carry over unchanged from the
parent: the session carrier ruling (`D-WORLD-26` = ARM A: `Authorization: Bearer
<session-credential>`, fail closed on absent/malformed/unknown, never an API key), the propose →
verify → commit invocation path, the one-snapshot-per-request rule, and P6.V's verified
commit-boundary law (which the parent lands and this half consumes).

## Non-negotiables that survive into any future draft

- No local JSON-RPC/MCP dispatch implementation (parent P1 / Design Freeze / `DESIGN.md` §3.7).
- No MCP SDK import into the gated daemon graph (charter clauses 2 + 3; the measured
  28-violation closure, E3).
- The narrow PACKAGE-path allowlist discipline: admit the delivered seam by package path, never
  a module root (the parent's measured matcher semantics).
- Session authority per D-WORLD-26 = ARM A, with both of its constraints.

## Quorum status

**This doc is NOT quorum-cleared.** It was authored at the split (iteration 125) and no reviewer
has passed it — the round-3 quorum blocked the PARENT, and this doc is the carrier of that
block's surviving objection, not a resolution of it. When the blocking predicate above reads
unblocked and this item is picked, the real design is written against the delivered seam and
MUST then go through the full design quorum (`ailang design-quorum`, reject-by-default
synthesis) at pick time. Nothing in this doc pre-authorizes a sprint.

## Relationship to the parent and to charter clause 6

Charter clause 6 names *"project the transition registry over MCP + publish the A2A agent
card."* The parent (`w-mcp-projection.md`) delivers the A2A half plus the enabling milestones —
P6.T (toolchain floor), P6.D (pinned `serveapi/protocol` dependency + one narrow allowlist
line), P6.V (verified commit-boundary law) — all of which this half will ALSO consume when it
unblocks. **Clause 6 is satisfied only when BOTH docs land.** Neither doc narrows the clause;
they partition it, and this one carries the blocked half with its predicate named and runnable.

## Related Documents

- [w-mcp-projection.md](w-mcp-projection.md) — the SPLIT parent: A2A half + enablers, executable
  now; its split-round quorum entry records this disposition
- [world-mission.md](../world-mission.md) — clause 6; D-WORLD-5 (the ask-upstream default this
  doc executes); D-WORLD-26 (the carrier ruling that survives into any future draft); the
  iteration-124 STATUS stamp holding the controller's first-party closure measurements this doc
  inherits (E3/E4)
- [DESIGN.md](../DESIGN.md) — §3.7 protocol-native: the reinvention ban that closes route (i)
- upstream [`sunholo-data/ailang#885`](https://github.com/sunholo-data/ailang/issues/885) — the
  blocking ask (SDK-free MCP dispatch seam)
- upstream [`sunholo-data/ailang#764`](https://github.com/sunholo-data/ailang/issues/764) →
  `serveapi/protocol@v0.33.2` — the precedent: the same ask-upstream route, previously delivered
