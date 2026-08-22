# w-workbench-read-only — a truthful localhost lens over the immutable world

- **Status**: Planned — unparked; the ratified defer-behind-item-18 condition is satisfied by
  measurement (V16), ready for the mission loop's next routing decision
- **Item**: queue item 14, `w-workbench-read-only`, clause-5
- **Filed**: 2026-08-11, human attended
- **Designed**: 2026-08-13, iteration 82
- **Quorum history**: two rounds, BLOCKED both times, all reviewer slots present (no N−1
  degrade); parked `needs-human-review` 2026-08-13
- **Human ruling**: 2026-08-14, Mark Edmondson, attended — **Option B, RATIFIED**: defer behind
  item 18 `w-daemon-read-cancellation`; item 14 does NOT grow into `host/store` context plumbing;
  item 18 gives all eight routes one elapsed-time contract and item 14 ships its renderer onto
  that base; residual WB-R1 is discharged by item 18, not by this item. That park is now
  discharged: item 18 is complete and its merge `d21754f` is an ancestor of `origin/dev` (V16)
- **Revised**: 2026-08-22, iteration 109 — every codebase claim re-measured on the post-item-18
  tree; both round-2 blocking objections answered (§12.1); changes enumerated in §12.2
- **Estimated**: 1.5–2 days
- **Measurement base**: `93e1ba5` (replaces `9491a10`, which had gone 65 non-`design_docs` files
  stale, including every daemon, store, and gate-script file this document cites)
- **Files changed by implementation**: enumerated in §8; this document changes only itself
- **Design result**: one server-rendered, read-only localhost route; no SSE, no write verbs,
  no copied grade policy, and no claim that the present log is already a complete provenance graph

Every present-tense codebase statement is backed by a command and observed output in §11, re-run
first-party at `93e1ba5` on 2026-08-22. Rows of the 2026-08-13 log that are no longer true are
replaced, not kept; §11 marks each and §12.2 summarizes.

---

## 1. Problem

The binding human surface says the World, not a transcript, is the workspace. Its reference
chrome is a world browser, timeline scrubber, and provenance explorer. Scenario 3 makes that
promise concrete: begin at a deployment, walk backward through the deploying transition,
proposal, agent, contemporaneous evidence, and originating goal, then compare now with the
world immediately before the surprising transition.

The daemon is the right process boundary for a local workbench, but its measured surface is
narrower than the charter prose suggests. It has seven GET patterns and one POST pattern
registered in `host/daemon/daemon.go` — eight, not the charter row's six (V1). None serves HTML,
SSE, embedded assets, or a file server (V2). The route-table comment miscount ("seven patterns")
that the 2026-08-13 revision recorded as a non-blocking defect has since been fixed by item 18's
work: at `93e1ba5` the comment reads "The eight patterns below are the complete frozen v1 table"
and eight registrations follow it (V1).

The seven GETs can supply:

1. daemon health;
2. selected head;
3. a world by content reference;
4. an object envelope, with payload only when explicitly requested;
5. one log entry by index;
6. a bounded log page; and
7. one registry head by semantic name (V3).

Those reads are enough for a useful first workbench: select the head or a named world, page the
timeline, select an entry, inspect its transition reference, and follow any reference that is
actually present. They are not a typed deployment → proposal → agent → evidence → goal index.
The log entry exposes a transition object reference, but its HTTP envelope has no proposal,
agent, evidence, goal, predecessor-world, or successor-world fields (V3). The store types likewise
record world, log, and object envelopes, rather than a reverse-edge provenance index (V4).

This distinction is the honesty boundary. Queue item 14 must build the grammar against real
objects, but it must not manufacture Scenario 3’s absent edges from strings inside an opaque
payload. The first release therefore provides a **partial, explicit provenance walk**: every
stored edge is clickable; every absent typed edge is a named `UNAVAILABLE` stop with the missing
relation. Full Scenario 3 becomes satisfiable only when a producer stores typed relations or a
future read projection exposes them. That missing substrate is not smuggled into this 1.5–2 day
renderer item.

Trust rendering has a second boundary. The canonical `Evidence → EvidenceGrade` function exists
only in AILANG. No Go source names `gradeOf`, `EvidenceGrade`, or `gradeCode` (V5). Copying its six
arms into a Go switch — `Evidence` gained a sixth constructor, `ProofReceipt`, since the previous
measurement base — would create a second, unproved policy and the exact grade-laundering surface
the UX forbids. The stale queue instruction to print `UNSUPPORTED` for unmapped variants cannot be
followed: the exported grade type in `world/types.ail` is exactly four constructors,
`PROVEN | TESTED | ATTESTED | CLAIMED`, and no `UNSUPPORTED` token exists anywhere under `world/`
(V6).

Finally, `TestReport(ref, false)` canonically has grade `TESTED`, exactly like the true case: the
mapping arm `TestReport(_, _) => TESTED` (`world/types.ail:45` in the contract, `:54` in the body)
never reads the pass/fail bool, and the committed test vectors grade the `true` and `false`
reports the same code (V6). That grade describes the kind of evidence, not its verdict, so at the
source a failing test report is badge-indistinguishable from a passing one — a rendering hazard
the workbench must surface, not paper over. A badge saying only `TESTED` over a failed test would
be materially misleading. Grade and verdict must be separate channels.

## 2. The design question, settled

### 2.1 Question

What is the smallest persistent localhost surface that makes real worlds, time, and stored
provenance navigable without widening authority, duplicating epistemic policy, or pretending the
current storage model contains relations it does not?

### 2.2 Decision

Add one route outside frozen `/v1`:

```text
GET /workbench
GET /workbench?world=<HashRef>
GET /workbench?from=<non-negative-index>&entry=<non-negative-index>
GET /workbench?object=<HashRef>&payload=0|1
```

These are query states of one renderer route, not four route patterns. Links and forms use GET.
There is no form with a mutating method and no handler path to `Store.Commit`.

The response is server-rendered HTML from a new transport-free package, `host/workbench`.
`host/daemon` owns the HTTP adapter and reads from the already-open store. The package split is:

```text
host/daemon
  owns: route registration, query parsing, status codes, deadline- and cardinality-bounded
        store reads (via the existing readCtx helper; §3.2, §3.3)
  imports: net/http, host/store, host/workbench

host/workbench
  owns: view model, escaping template, navigation grammar, unavailable states
  imports: html/template, io; no net/http, no store, no registry, no cloud SDK
```

The new package is deliberate. `host/boundary` is pinned to exactly one Go file, so no renderer
file lands there and its pin is not relaxed (V7). `host/workbench` is presentation logic rather
than a semantic boundary. Keeping it free of `net/http` makes the transport authority remain in
the daemon closure instead of creating a second HTTP owner.

### 2.3 Why server-rendered HTML

Choose server-rendered HTML with inline CSS and ordinary links/forms.

- It adds one representation route, not a second JSON protocol plus a static-asset protocol.
- It works with JavaScript disabled and makes every navigation state bookmarkable.
- `html/template` escapes object labels, provenance strings, semantic IDs, and payload previews.
- There is no asset build, embedded filesystem, CDN, package manager, or cloud dependency.
- A full page request has bounded work cardinality, bounded response size, and — since item 18
  landed — a bounded elapsed store wait: every getter the workbench calls is context-aware and
  the daemon derives a 10 s request-scoped deadline from `r.Context()` (V16–V18).
- Progressive enhancement may later add local JavaScript, but it is not an acceptance dependency.

Reject “JSON API plus static assets” for this item. The daemon already has the necessary JSON
reads. A new aggregate JSON API would become another public contract before the typed projection
exists; static assets would add caching and content-type gates without improving the first
read-only grammar.

Reject SSE. The server’s `WriteTimeout` configures a 30-second response-write deadline (V8); it is
not a handler or store-call cancellation mechanism — that job now belongs to the request-scoped
read deadline (V16–V18), which a long-lived response would fight. Keeping a long-lived response
open would require deliberate deadline handling through `http.ResponseController`, which this item
does not need because the log is immutable and explicit refresh is sufficient. This design neither
uses `ResponseController` nor changes any timeout or deadline constant.

### 2.4 Route and version posture

`/v1/*` remains the machine API. `/workbench` is an operator renderer whose HTML may evolve; it is
not inserted into the frozen v1 table and it does not change the semantics of the existing eight
patterns. The route table becomes nine registrations in `host/daemon/daemon.go`: eight existing
`/v1` patterns plus `GET /workbench`. The route-table comment already counts eight correctly (V1);
implementation keeps it truthful about the eight-frozen/one-renderer split when the ninth
registration lands (§3.5).

**The workbench query grammar is closed. The only accepted keys are `world`, `object`, `from`,
`entry`, and `payload`. Any unknown key, duplicate scalar key, or unsupported parameter
combination returns `400 Bad Request` with a constant HTML error message; no parameter is ignored
and no precedence fallback is applied.** (Round-3 carve-out, `gpt5-6-sol`'s own replacement text,
applied verbatim — the prior rule “unknown parameters are ignored only if they do not select
data” was a SILENT FALLBACK and the mission axiom forbids one: `?paylod=1` would have rendered a
different view instead of refusing. Pinned by `TestWorkbenchRefusalBranches/unknown-parameter`
and `/duplicate-parameter`, mutations M31/M32.) Malformed values
for `world`, `object`, `from`, or `entry` return `400` with a small HTML error page. Well-formed but
absent references return `404`. A store read that exceeds the request deadline returns `503`, the
same class the JSON routes emit (V18). Other store failures return `500` carrying only a constant
message. Those branches preserve a distinction the existing JSON handlers measurably model, not an
assumed one: the six `handle` functions in `host/daemon/handlers.go` spread 18 `writeAPIError`
sites across `400`/`404`/`503`/`500` (plus `413`/`409` on the write path), with malformed input
parsed at the boundary returning `400` and well-formed-but-absent references returning `404`
(V20). The no-raw-error posture is likewise measured: exactly one `Internal` site exists and it
passes the constant `internalErrorMessage`; every remaining `err.Error()` site is a
`400 BadRequest` branch (V21). The workbench inherits both postures rather than claiming them.

### 2.5 Authority and bind posture

The workbench inherits the daemon’s startup bind decision. `New` accepts only `127.0.0.1`, `::1`,
or `localhost`, and tests exercise both accepted and refused arms (V9). There is no workbench
override, proxy, callback URL, remote font, analytics endpoint, fetch to an arbitrary origin, or
link that grants an agent ambient egress.

Loopback is locality, not authentication. The workbench therefore exposes only facts already
available through the daemon’s read paths and never displays secrets by default. Object payloads
remain opt-in. A payload preview is capped before conversion to a string and is visibly marked as
raw bytes, not interpreted HTML. The renderer sets:

```text
Content-Type: text/html; charset=utf-8
Cache-Control: no-store
Content-Security-Policy: default-src 'none'; style-src 'unsafe-inline'; base-uri 'none';
                         form-action 'self'; frame-ancestors 'none'
X-Content-Type-Options: nosniff
Referrer-Policy: no-referrer
```

The inline-style exception is narrow and avoids a second asset route. No script source is allowed.
All outbound-looking hash links are relative `/workbench?...` links.

### 2.6 Grade decision

The renderer never computes an evidence grade in Go.

For the real objects available in this item, no canonical decoded evidence/grade projection is
present. The UI renders `GRADE UNAVAILABLE — no canonical host projection` beside an evidence
slot only when the underlying typed object says evidence should exist but no grade value is
available. It does **not** render the nonexistent grade `UNSUPPORTED`, and it does not downgrade
the object to `CLAIMED`. Absence of a grade is a rendering state, not an `EvidenceGrade`.

If a future canonical projection supplies both a four-value grade and evidence details, the view
model accepts them as display data; it still does not map constructors. The view type uses a
closed display enum validated at construction (`PROVEN`, `TESTED`, `ATTESTED`, `CLAIMED`) plus a
separate availability/result type. Invalid labels are refused, not printed as badges.

For `TestReport`, the view model requires a separate verdict field. It renders, for example:

```text
TESTED   verdict: FAIL
```

The FAIL marker uses text, icon/shape, and accessible label; it is not color alone. A test report
without a verdict is `GRADE UNAVAILABLE — test verdict missing`, even if a caller supplies
`TESTED`. This deliberate fail-closed rule prevents a method badge from concealing the result.
No current constructor may produce `PROVEN`; the workbench never infers it from “Z3”, “proof”, a
replay result, a semantic ID, or a hash (V6).

### 2.7 The Scenario 3 walk, precisely

The timeline pane pages `GET /v1/log`-equivalent store reads in ascending entry order, at most
100 entries per workbench request. Selecting entry N shows the measured fields below (V3b):

```text
entry N
  entryHash
  prevEntryHash
  semanticsEpoch
  transitionFn ───────────────> object inspector if object exists
  interpreter ────────────────> object inspector if object exists
  transitionRef ──────────────> object inspector if object exists
  writtenBy (literal identity string, not an agent-object link)
```

`entry N` is the frozen header's `EntryIndex int64`, rendered once as the timeline selection label;
it is not repeated as a second “entry index” detail. The other frozen header fields are
`PrevEntryHash`, `SemanticsEpoch`, `TransitionFn`, `Interpreter`, and `WrittenBy` (V3b). `EntryHash`
and `TransitionRef` are fields of `LogEntry` outside the frozen `LogHeader`; in particular,
`TransitionRef` is not presented as a header field (V3b).

The selected world pane shows its revision, state root, and log head; state root is an object link.
The current schema supplies no world-by-revision/index lookup. Therefore “world before entry N” is
shown only if an exact world reference is already present in a typed edge supplied to the view.
It is otherwise `UNAVAILABLE: predecessor world relation is not stored by this projection`.

Likewise, the walk stops explicitly at the transition object when typed proposal, agent, evidence,
deployment, policy, or goal relations are absent. `writtenBy` is not promoted to an agent link,
object `provenance` text is not parsed as a reference, and payload bytes are not searched for
hash-shaped substrings. This serves Scenario 3’s navigation grammar and exposes the exact missing
route/projection: **a bounded typed provenance projection from a root object or log entry, including
edge kind, target reference, and historical world reference**. Adding that projection is separate
work because the present seven GET routes cannot produce it from their envelopes (V3–V4).

## 3. Proposed change

### 3.1 New package `host/workbench`

Add `host/workbench/render.go` containing:

- `Page`, `WorldView`, `TimelineView`, `EntryView`, `ObjectView`, `EdgeView`, and `GradeView`;
- constructors that validate bounds and grade/verdict combinations;
- `Render(io.Writer, Page) error` using one parsed `html/template`;
- constants for the renderer’s page limit and payload-preview byte cap; and
- no network, database, filesystem, subprocess, or AILANG imports.

Add `host/workbench/render_test.go` with hostile-string escaping, all empty/unavailable states,
timeline navigation, relative-link checks, grade-label validation, the failed-test dual-channel
case, and size-cap tests. Template tests render to `bytes.Buffer`; they do not bind a socket.

The HTML uses landmarks and headings: header/status, world browser, timeline, inspector, and
provenance walk. Keyboard focus follows the document order. Hashes use full values in copyable
text, with a shortened visual label only through CSS/markup; the accessible name remains exact.

### 3.2 Daemon adapter

Add `host/daemon/workbench.go` containing `handleWorkbench`, query decoding, cardinality-bounded
data loading, security headers, and mapping from store envelopes to the presentation view model.
Register exactly one pattern in `Daemon.Handler`:

```go
mux.HandleFunc("GET /workbench", d.handleWorkbench)
```

The adapter calls the same store methods the existing JSON handlers call, through the same
five-method `readStore` seam (V26), and derives its context exactly the way they do: one call to the
existing `d.readCtx(r)` helper per request, `defer cancel()` immediately, and the derived context
threaded to every getter (V16–V17). It invents no second deadline-derivation helper and spells no
`context.Background()` store read — the committed ratchet `TestNoNewDeadlineFreeStoreReads` reds
on any new deadline-free call site under `host/` or `cmd/` (V19) — and it maps deadline expiry to
the same `503`/`Timeout` posture the JSON routes use (V18). It does not call the JSON handlers
through an in-process HTTP client; that would add encoding work and blur error ownership. The
adapter does not expose `d.store` to the presentation package.

Timeline reads stop after 100 entries or the first missing index. Query `from` is non-negative and
addition is overflow-checked. Payload preview is off by default and capped at 64 KiB with a visible
“truncated” marker. The workbench cannot request the daemon API’s maximum 500-entry page in one HTML
render; rendering work and response size remain independently bounded. Elapsed store wait is
bounded by the inherited request-scoped deadline (§3.3), not by these caps.

Add `host/daemon/workbench_test.go`. Tests use `httptest.NewRecorder` and the constructed handler,
not `httptest.NewServer`, so the route tests are hermetic in restricted environments. They seed real
store objects and assert status, headers, escaping, links, page bounds, opt-in payload, and every
refusal branch.

### 3.3 The elapsed-time base: WB-R1 is discharged by item 18

The 2026-08-13 revision of this section accepted **WB-R1: a workbench handler can remain blocked
in a store read after the HTTP write deadline**, because at base `9491a10` the store had no
`context.Context` plumbing, no daemon handler consulted the request context, and the production
DSN set no SQLite `busy_timeout`. Mark's attended 2026-08-14 ruling (Option B, ratified) answered
the objection by sequencing: item 18 `w-daemon-read-cancellation` lands first, and this item ships
its renderer onto that base. Item 18 has landed — its merge `d21754f` is an ancestor of
`origin/dev` (V16) — and every clause of the old residual is now false by measurement, so
**WB-R1 is DISCHARGED by item 18**:

- the daemon derives a request-scoped deadline from the request: the single `readCtx` helper
  returns `context.WithTimeout(r.Context(), d.readDeadline)` at `host/daemon/handlers.go:270`
  (V16);
- `readDeadline` is a real constant, `10 * time.Second` at `host/daemon/daemon.go:128`, held as a
  `Daemon` field and wired by `New` (V18);
- every store method the workbench will call is context-aware: `GetObject`, `GetWorld`,
  `GetLogEntry`, `GetRegistryHead`, and `SelectedHead` in `host/store/store.go`, and `ReadObject`
  in `host/store/read_object.go`, all take `ctx context.Context` first (V17); and
- expiry is an explicit wire status, not a hang or a silent fallback: `writeReadTimeout` emits
  `503` with class `Timeout` (`host/daemon/handlers.go:324`), and the committed
  `TestDaemonReadDeadline/blocking-store` and `TestReadCtxCancelledAfterHandler` arms red when
  cancellation or context propagation is removed (V25).

The workbench's obligation on that base is inheritance, not reinvention. Its handler MUST derive
its context through the same `d.readCtx(r)` helper the JSON routes use — never a second helper,
never `context.Background()` — and `GET /workbench` must join the cancelled-after-handler and
blocking-store test tables so that removing the propagation reds a named test (AC8; mutations M29,
M30). A workbench that met every rendering criterion but derived its own deadline would be a
second timeout policy, which this item forbids exactly as it forbids a second grade policy. The
100-entry and 64 KiB caps still bound work and response size only.

Two residuals remain, and both are external to this item:

- **Lock waits are bounded by SQLite `busy_timeout`, not by the request deadline.** The
  production DSN sets `busy_timeout(2000)` — `busyTimeoutMillis = 2000` at
  `host/store/writer_lock.go:181` — and no production code links that 2 s window to the 10 s
  `readDeadline`; the composition is safe today only because 2 s < 10 s, an ordering nothing
  asserts (V22). That is queue item 22, `w-daemon-lock-wait-not-deadline-bound`, a separate row
  this item names and does not absorb or fix.
- **A read that completes after the deadline but before cancellation is observed still answers
  200.** The timeout classifier runs only on the error path; the code names this limitation
  first-party (`LIMITATION(w-daemon-late-read-503)` in `host/daemon/handlers.go`), and it belongs
  to that successor item, not to this renderer (V22).

### 3.4 Boundary gate adjustment

Do not add a file to `host/boundary`; keep its exact file-count pin at one. Modify the existing
`allowlist_world_test.go` only.

The existing protected group `cmd/ailang-worldd` includes the daemon in its dependency closure and
is the sole group whose extra bare-`net/http` deny list is empty (V10). Re-assert its positive arm
rather than accepting inherited green:

- keep `cmd/ailang-worldd` the only per-group bare-`net/http` exception;
- add a named assertion that its dependency closure contains exactly one `net/http` package after
  `host/workbench` is added;
- assert `host/workbench` itself has zero `net/http`, `net/url`, cloud, registry, and broker imports;
- assert `cmd/ailang-worldd` reaches `host/workbench` exactly once in its package closure; and
- add a same-test known-positive transport control using the existing daemon closure.

This changes neither the protected group count nor the `wantFileCount = 1` pin. It makes removal
of the pure-package boundary red, instead of relying on a green inherited exception.

### 3.5 Documentation accuracy at the ninth registration

The old "seven patterns" miscount is already gone: at `93e1ba5` the comment above the route table
reads "The eight patterns below are the complete frozen v1 table" and eight registrations follow
(V1). When the ninth registration lands, implementation must keep that comment truthful — eight
frozen `/v1` machine patterns (seven GET, one POST) plus one unversioned read-only workbench
pattern — rather than letting the renderer route blur into the frozen table. This is comment
maintenance, not a route deletion or v1 redesign.

## 4. What this buys — and what it does not

This item buys a zero-install browser-opened surface over the actual local daemon; a reusable visual
grammar for objects, time, links, unavailable edges, trust availability, and failure verdicts; a
work- and response-size-bounded route that cannot commit; and executable seams the approval inbox
can later compose.

It de-risks the medium decision cheaply. If server-rendered navigation proves too limited, the
typed view model and read adapter remain useful; no public aggregate JSON schema or streaming
contract has been frozen.

It does not buy complete Scenario 3 data. The operator can walk every stored edge and see precisely
where the chain stops, but deployment/proposal/agent/evidence/goal relations require a canonical
typed producer/projection. It does not buy world diff: the current world envelope contains roots,
not a state-tree diff protocol. It does not buy live updates, replay execution, or a grade decoder.
Elapsed-time cancellation for store reads is not bought here either — item 18 already bought it,
and this item inherits it (§3.3). What remains unbounded-by-deadline is the lock-wait window,
which is item 22's row, not this one (V22).

The success standard is therefore truthful navigation, not a staged mock that always reaches the
goal. A visible unavailable edge is successful behavior when the substrate lacks the relation.

## 5. Persistent non-vacuity

Every property intended to survive implementation is attached to a committed Go test or the
existing Go gate. One-shot greps appear only as baselines and review aids.

| Persistent property | Gate | Positive/control arm |
|---|---|---|
| route is GET-only | `TestWorkbenchRouteIsReadOnly` | GET returns 200; POST/PUT/DELETE return 405 |
| route renders real store data | `TestWorkbenchRendersSeededWorldAndTimeline` | two distinct seeded objects render distinct refs |
| malformed differs from absent | `TestWorkbenchRefusalBranches` | valid seeded ref returns 200 |
| query grammar is closed — nothing is ignored | `TestWorkbenchRefusalBranches/unknown-parameter`, `/duplicate-parameter` | the five accepted keys, each supplied once, return 200 |
| payload is opt-in and capped | `TestWorkbenchPayloadPreviewBound` | small payload renders in full |
| hostile values are escaped | `TestRenderEscapesAllObjectText` | safe literal remains visible |
| links stay local and relative | `TestRenderEmitsOnlyLocalLinks` | expected `/workbench?...` link exists |
| missing edge is explicit | `TestRenderUnavailableProvenanceEdge` | stored edge renders as a link |
| grade is never synthesized | `TestGradeViewRejectsUnavailableOrInvalidInput` | each of four valid labels is accepted |
| failed test cannot hide behind TESTED | `TestGradeViewRequiresTestVerdict` | TESTED+PASS and TESTED+FAIL both render both channels |
| renderer has no transport authority | `TestWorkbenchPackageRemainsTransportFree` in boundary gate | daemon closure still contains exactly one `net/http` |
| loopback refusal remains two-sided | existing `TestNewRefusesNonLoopbackBind` | its accepted loopback table remains |
| server timeout constants remain configured | existing `TestBoundedWaitsAndBodyLimit` | exact non-zero constants remain; store cancellation is gated separately by the two rows below |
| workbench reads carry the request deadline | `TestWorkbenchReadDeadline` + `GET /workbench` in the `TestReadCtxCancelledAfterHandler` route table | existing `/v1` arms keep passing; blocking-store arm answers `503` class `Timeout` |
| no new deadline-free store read | existing `TestNoNewDeadlineFreeStoreReads` | pins {approve 8, registry 2, replay 1}; a workbench `context.Background()` read site reds it (V19) |
| boundary file pin does not move | existing boundary AST/file-count test | it enumerates the one real Go file |

`scripts/verify_go.sh` already runs `go build ./...`, plain tests, and race tests (V11). The new
package and tests are therefore in the persistent CI path without adding another shell total.

## 6. Mutation table

Each mutation is a compiling one-line neuter. Refusal branches are separate rows.

| ID | Guard/refusal | One-line mutation that must red | Named test |
|---|---|---|---|
| M1 | method pattern is GET-only | change pattern to `"/workbench"` | `TestWorkbenchRouteIsReadOnly` |
| M2 | malformed world ref → 400 | `if false && err != nil {` | `TestWorkbenchRefusalBranches/malformed-world` |
| M3 | absent world → 404 | `if false && !ok {` | `TestWorkbenchRefusalBranches/absent-world` |
| M4 | malformed object ref → 400 | `if false && err != nil {` | `TestWorkbenchRefusalBranches/malformed-object` |
| M5 | absent object → 404 | `if false && !ok {` | `TestWorkbenchRefusalBranches/absent-object` |
| M6 | negative `from` → 400 | `if false && from < 0 {` | `TestWorkbenchRefusalBranches/negative-from` |
| M7 | malformed `entry` → 400 | `if false && err != nil {` | `TestWorkbenchRefusalBranches/malformed-entry` |
| M8 | absent entry → 404 | `if false && !ok {` | `TestWorkbenchRefusalBranches/absent-entry` |
| M9 | store error → 500 | `if false && err != nil {` | `TestWorkbenchRefusalBranches/store-error` |
| M10 | payload remains opt-in | change `showPayload` default to `true` | `TestWorkbenchPayloadPreviewBound/default-off` |
| M11 | payload cap | change `preview = preview[:maxPayloadPreview]` to `preview = payload` | `TestWorkbenchPayloadPreviewBound/oversize` |
| M12 | timeline page cap | change `workbenchPageLimit` from 100 to 101 | `TestWorkbenchTimelineBound` |
| M13 | overflow refusal | `if false && from > math.MaxInt64-int64(limit) {` | `TestWorkbenchRefusalBranches/from-overflow` |
| M14 | HTML escaping | wrap hostile label in `template.HTML(...)` | `TestRenderEscapesAllObjectText` |
| M15 | local-only links | prefix object link with `https://example.com` | `TestRenderEmitsOnlyLocalLinks` |
| M16 | absent edge stays visible | `if false && !edge.Available {` | `TestRenderUnavailableProvenanceEdge` |
| M17 | invalid grade refused | `if false && !validGrade(label) {` | `TestGradeViewRejectsUnavailableOrInvalidInput` |
| M18 | unavailable grade not CLAIMED | return `GradeView{Label:"CLAIMED"}` on missing input | `TestGradeViewRejectsUnavailableOrInvalidInput` |
| M19 | TestReport requires verdict | `if false && kind == TestReport && verdict == nil {` | `TestGradeViewRequiresTestVerdict/missing` |
| M20 | failed verdict shown | render verdict as constant `PASS` | `TestGradeViewRequiresTestVerdict/fail` |
| M21 | no inferred PROVEN | change unavailable label to `PROVEN` | `TestGradeViewRejectsUnavailableOrInvalidInput/no-proven-inference` |
| M22 | CSP forbids scripts | delete `default-src 'none'` token | `TestWorkbenchSecurityHeaders` |
| M23 | no-store | set `Cache-Control` to `public` | `TestWorkbenchSecurityHeaders` |
| M24 | transport-free renderer | add blank import `_ "net/http"` to `host/workbench` | `TestWorkbenchPackageRemainsTransportFree` |
| M25 | daemon positive transport arm | neuter exact `net/http == 1` assertion | boundary gate’s mutation subtest for daemon transport control |
| M26 | loopback refusal | `if false && !isLoopbackHost(cfg.BindHost) {` | existing `TestNewRefusesNonLoopbackBind` |
| M27 | write-timeout configuration remains unchanged | set `WriteTimeout: 0` | existing `TestBoundedWaitsAndBodyLimit` (connection-write configuration only; not an elapsed store-wait guard) |
| M28 | boundary file-count pin | change `wantFileCount` to 2 | existing boundary AST/file-count test |
| M29 | deadline derived via `readCtx` | replace `d.readCtx(r)` with `context.WithCancel(r.Context())` | `TestWorkbenchReadDeadline/blocking-store` |
| M30 | store timeout renders 503, not 500 | `if false && timedOut(ctx, err) {` in the workbench adapter | `TestWorkbenchReadDeadline/blocking-store` (status + class assertion) |
| M31 | unknown query key → 400 (closed grammar) | `if false && !isAcceptedKey(k) {` in the workbench query parser | `TestWorkbenchRefusalBranches/unknown-parameter` |
| M32 | duplicate scalar query key → 400 (closed grammar) | `if false && len(vs) > 1 {` in the workbench query parser | `TestWorkbenchRefusalBranches/duplicate-parameter` |

M25 requires the boundary test to exercise its own detector with an overlay, following the file’s
existing mutation-test idiom. Merely asserting the dependency count once is not enough.

M29/M30's observable is the HTTP status and error class on `GET /workbench` under a blocking
store fake, and it is downstream of the mechanism and of nothing else: the fake is not SQLite, so
no `busy_timeout` exists in that test to expire and produce the same value (V22), and
`http.Server.WriteTimeout` closes the connection rather than writing a `503` body (V8) — only the
`readCtx`-derived deadline plus the `timedOut` classifier can put `503`/`Timeout` on that
recorder. Under M29 the fake blocks with no deadline to expire, so the test's own wait bound
trips instead of a 503 arriving; under M30 the expiry is misclassified as `500`. Both mutations
compile.

## 7. Acceptance criteria

Every command below was baselined on the unmodified `93e1ba5` tree on 2026-08-22 (this revision's
environment, which — unlike the 2026-08-13 one — permits `listen(2)`; V12). “Fail trigger” states
how the criterion can red. “Producible” states why the named mechanism can emit the observable.
AC3 and AC4 are green at base by construction; they are regression pins whose red condition is
sprint-induced breakage, and neither is citable as evidence any workbench behavior exists —
AC1/AC2/AC5/AC8 carry the feature evidence.

### AC1 — package and route behavior

Command:

```sh
go test ./host/workbench ./host/daemon -count=1
```

Baseline at `93e1ba5`: `host/workbench` does not exist, so the exact command fails at base —
`go test` cannot match the package — and the criterion cannot pass vacuously. The daemon package
alone passes at base (`ok`, rc=0; V12). Implementation tests named in §5 still use recorders and
buffers so they stay hermetic in environments where `listen(2)` is refused.

Pass: all new renderer and handler tests pass. Fail trigger: any mutation M1–M23 or M31–M32 (the closed-grammar refusals) or a compile
failure reds a named test/package. Producible: store fixtures, `httptest.ResponseRecorder`, and
`bytes.Buffer` directly generate every asserted status, header, link, and body.

### AC2 — boundary authority remains explicit

Command:

```sh
go test ./host/boundary -run 'Test(BareNetHTTPExemptionIsPerGroup|WorkbenchPackageRemainsTransportFree)' -count=1
```

Baseline at `93e1ba5`: existing `TestBareNetHTTPExemptionIsPerGroup`
(`host/boundary/allowlist_world_test.go:875`, V10) passes within the green boundary package
(V12); the new named test is absent at base, so the exact combined selector presently exercises
only the existing arm and is not accepted as implementation evidence. Pass requires test
enumeration/log output proving both names ran.

Fail trigger: adding transport to `host/workbench`, losing daemon `net/http`, adding a second
exception, or disarming the overlay mutant. Producible: `go list -deps -json` and the existing
overlay/read helpers expose package closures and import substitutions in this test package.

### AC3 — full Go gate

Command:

```sh
AILANG_BIN=/tmp/ailang-v0300/ailang scripts/verify_go.sh
```

Baseline at `93e1ba5`: run in full for this revision — `✓ go gate PASSED: build clean, plain and
race tests pass with pinned AILANG_BIN (AILANG v0.30.0 …)`, rc=0 (V24). Green at base by
construction: this AC is a regression pin, and its power to red is the implementation breaking a
build, test, race, hygiene, or pin leg — it is never evidence the workbench works.

Fail trigger: build, plain tests, race tests, binary hygiene, toolchain canary, or AILANG pin fails.
Producible: the committed script invokes each observable and is the repository’s existing Go CI
gate (V11). A sprint may not waive this AC on the strength of AC1, and may not cite its green as
feature evidence.

### AC4 — AILANG gate totals do not move

Command:

```sh
AILANG_BIN=/tmp/ailang-v0300/ailang scripts/verify_ail.sh
```

Baseline at `93e1ba5`: exit 0; 11 allowlisted modules, 10/10 required verified identities, 40
required named tests, and the nine-step package gate pass (V13 — the totals moved from the old
base's 5/20, which is exactly why a stale baseline could not be kept).

Pass: the exact same totals and package gate pass after implementation. Fail trigger: any AILANG
module, contract identity, required test, package projection, manifest, tar member, or golden
moves. Like AC3, green at base by construction — a regression pin, not feature evidence.
Producible: this is a Go-only item and need not alter any input enumerated by the script.

### AC5 — route cardinality and representation

Command:

```sh
grep -c 'mux.HandleFunc' host/daemon/daemon.go
```

Baseline: `8` (V1). Pass: exact output `9`, with seven `/v1` GETs, one `/v1` POST, and one
`GET /workbench`. Fail trigger: route omitted, duplicated, or another route added. Producible: one
literal `mux.HandleFunc` registration increases this detector by exactly one.

This is an acceptance snapshot, not the persistent gate. `TestWorkbenchRouteIsReadOnly` and the
daemon handler tests persist method/path behavior.

### AC6 — no new long-lived transport

Command:

```sh
test -d host/workbench && \
! grep -rniE 'text/event-stream|ResponseController|http\.Flusher|http\.FileServer|embed\.FS' \
  --include='*.go' host/workbench host/daemon && \
grep -q 'http' host/daemon/daemon.go
```

Baseline at `93e1ba5`: `host/workbench` is absent, so this exact command exits 1 at the root
`test -d` assertion (V23) and cannot pass vacuously; same-scope daemon control finds
`http.Server`/`net/http` in V2/V8. Pass after implementation is root assertion success, zero
forbidden-pattern hits, and a positive HTTP control in the same call.

Fail trigger: SSE, deadline relaxation, streaming/file-server, or embedded-asset machinery appears.
Producible: the chosen HTML renderer needs none of those symbols. Persistence comes from imports
and behavior tests; this grep is a review assertion, not the sole gate.

### AC7 — only priced files changed

Command:

```sh
git diff --name-only 93e1ba5 -- ':!design_docs'
```

Baseline: empty output, rc=0 at revision time (V14). The pathspec excludes `design_docs` because
the checkout carries controller-staged `design_docs` moves and this revised doc, none of which is
implementation scope. For implementation, pass is that every listed path is one of the files
enumerated in §8. Fail trigger: any unpriced non-`design_docs` path appears. Producible: the
implementation is fully contained by the listed Go source/tests and the existing boundary gate
file.

This snapshot does not persist after merge; code ownership and ordinary review enforce scope.

### AC8 — the workbench read path carries the request deadline

Command:

```sh
go test ./host/daemon -run 'TestWorkbenchReadDeadline|TestReadCtxCancelledAfterHandler' -count=1 -v
```

Baseline at `93e1ba5`: `TestReadCtxCancelledAfterHandler` and `TestDaemonReadDeadline` (including
its `blocking-store` arm) pass over the existing routes, and `TestWorkbenchReadDeadline` does not
exist (V25) — so the selector currently proves only the inherited arm and is not accepted as
implementation evidence. Pass requires `=== RUN` enumeration showing both names ran, the
cancelled-after-handler route table covering `GET /workbench`, and a blocking-store arm answering
`503` with class `Timeout` on `/workbench`.

Fail trigger: M29 (deadline derivation removed) or M30 (expiry classified as 500) reds the
blocking-store arm; dropping `/workbench` from the route table reds the enumeration requirement.
Producible: the blocking-store fake and recorder plumbing already exist in `host/daemon` for the
JSON routes (V25); the same stimulus applied to `/workbench` yields the asserted status because
only the derived deadline can expire against a non-SQLite fake (§6, M29/M30 note).

## 8. Conflict surface

### 8.1 Files implementation will add

| File | Reason / collision |
|---|---|
| `host/workbench/render.go` | new pure presentation package; likely approval-inbox composition seam |
| `host/workbench/render_test.go` | renderer non-vacuity and trust-honesty tests |
| `host/daemon/workbench.go` | HTTP adapter and cardinality-bounded store reads |
| `host/daemon/workbench_test.go` | route/refusal/security tests |

### 8.2 Files implementation will modify

| File | Exact change / collision |
|---|---|
| `host/daemon/daemon.go` | one GET registration; keep the (now-correct) eight-pattern comment truthful about the eight-v1/one-renderer split (§3.5) |
| `host/daemon/read_deadline_test.go` | add `GET /workbench` to the cancelled-after-handler route table and a workbench blocking-store arm (AC8) |
| `host/boundary/allowlist_world_test.go` | transport-free package guard and daemon positive arm; retain one-file pin and sole exemption |

### 8.3 Files inspected but not moved

| File/pin | Required disposition |
|---|---|
| `host/boundary/allowlist_world_test.go` `wantFileCount = 1` | remains exactly 1; no new boundary Go file |
| `protectedGoGroups` four-row table | remains four rows; `cmd/ailang-worldd` remains sole bare-http exception |
| `scripts/verify_ail.sh` `LEG1_MODULES` | remains the 11-entry set |
| `scripts/verify_ail.sh` `EXACT_TOTAL_VERIFIED=10` | remains 10 (V15) |
| `scripts/verify_ail.sh` `REQUIRED_TESTS` | remains its current 40 identities |
| `scripts/verify_ail.sh` `EXACT_TOTAL_TESTS=40` | remains 40 (V15) |
| `scripts/verify_go.sh` | no edit; it already discovers all Go packages/tests |
| `world/types.ail` | no edit; remains canonical grade policy |
| `host/store/*` | no edit; its context-first getters are consumed as-is; no reverse index or query added |
| `host/daemon/handlers.go` `readCtx`/`timedOut`/`writeReadTimeout` | consumed by the workbench adapter, never duplicated; no edit to the helpers |
| `host/store/writer_lock.go` `busyTimeoutMillis = 2000` | untouched; its relation to `readDeadline` is item 22's scope (V22) |
| `design_docs/HUMAN-SURFACE.md` | binding input, not edited |
| `design_docs/SCENARIOS.md` | binding Scenario 3 input, not edited |

The likely collision is queue item 7 adding approval-inbox views to `host/workbench` and another
route in `host/daemon`. It must compose the view model rather than replace grade/refusal rules.
Transition/provenance work may later add a typed projection; it must not teach this renderer to
parse opaque payloads heuristically.

## 9. Scope and pricing

The 1.5–2 day price covers:

- 0.35 day: view model, template, accessibility, empty/unavailable states;
- 0.35 day: daemon query adapter — `readCtx`-derived deadline-bounded, cardinality-bounded reads;
- 0.35 day: route/refusal/security/payload tests plus the AC8 deadline-propagation arms;
- 0.25 day: boundary positive/negative mutation arms;
- 0.20 day: full gates, mutation run, and documentation correction;
- 0.20 day contingency for template and sandbox-independent test fixes.

No frontend toolchain, Node dependency, generated asset, API schema, database migration, or AILANG
change is priced. Discovery of a requirement for any of those stops the sprint and returns to
design; it is not absorbed as “small glue.”

## 10. What this item is NOT doing

This item is read-only. It explicitly does **not**:

- build the approval inbox or any approve/reject/defer/attenuate action;
- define or consume a decision-packet schema;
- choose or implement timeout policies — it inherits `readCtx`/`readDeadline` and changes neither;
- touch `busy_timeout`, link it to `readDeadline`, or otherwise absorb item 22 (V22);
- define a grade mapping, copy `gradeOf` into Go, or invent `UNSUPPORTED`;
- mint `PROVEN` or acquire item 17’s proof/replay authority boundary;
- build session-filtered MCP, A2UI, AG-UI, chat, or external projection;
- add SSE, WebSocket, polling, notification, analytics, CDN, or cloud dependency;
- add a shell, subprocess launcher, filesystem browser, or arbitrary URL fetch;
- add a commit, effect, capability, budget, policy-edit, or registry-write control;
- add a database table, reverse-edge index, state-tree diff, or provenance schema;
- claim full Scenario 3 traversal where typed edges are absent;
- reinterpret `writtenBy` or free-form provenance strings as content references;
- render color as the sole grade or verdict channel;
- change the daemon timeout constants or relax one route with `ResponseController`; or
- move any AILANG verification/module/test total.

## 11. Verification Log

All commands ran from repository root at `93e1ba5` on 2026-08-22 (iteration 109), re-run
first-party for this revision — nothing below is inherited from the 2026-08-13 log or transcribed
from the controller's table without re-measurement. Go commands ran with `GOTOOLCHAIN=go1.25.6`.
Empty-result rows include a root assertion and a same-scope positive control. Rows V12, V14, and
V16–V19 REPLACE former rows whose claims are no longer true; V20–V25 are new.

| ID | Codebase claim | Command | Observed output |
|---|---|---|---|
| V1 | Handler registers 8 patterns: 7 GET, 1 POST; the comment now counts eight | `grep -c 'mux.HandleFunc' host/daemon/daemon.go; grep -n 'mux.HandleFunc\|patterns' host/daemon/daemon.go` | `8`; seven `GET /v1/...` registrations plus `POST /v1/commit` at :555–:562; comment at :550 reads "The eight patterns below are the complete frozen v1 table" — the 2026-08-13 log's "seven patterns" miscount has been FIXED since `9491a10` |
| V2 | no production Go HTML/SSE/file/embed route in `host/`+`cmd/`; instrument sees HTTP in same scope | `test -d host && test -d cmd && find host cmd -name '*.go' ! -name '*_test.go' -type f -print0 \| xargs -0 grep -lE 'text/html\|text/event-stream\|http.FileServer\|embed.FS' \| wc -l; grep -c 'http' host/daemon/daemon.go` | roots exist; `0`; same-scope positive control `26` |
| V3 | exact response fields and bounded log composition | `grep -n '"revision"\|"stateRoot"\|"logHead"\|"semanticId"\|"provenance"\|"payload"\|"entryHash"\|"transitionRef"' host/daemon/handlers.go; grep -n 'clampLimit' host/daemon/handlers.go` | world `ref,revision,stateRoot,logHead` (:33–:36, :87–:90); object `hash,interfaceHash,semanticId,provenance[,payload]` (:40–:43, :79–:83); log `entryHash,transitionRef` (:58–:59, :104); `clampLimit` at :239, non-positive defaults to 100 (:238), Z3-proven ceiling 500 (:439) |
| V3b | `LogHeader` directly contains all six frozen fields; `EntryHash`/`TransitionRef` sit outside it | `grep -n "type LogHeader" -A 17 host/store/store.go` | header :112–:119: `EntryIndex int64`, `SemanticsEpoch int64`, `TransitionFn`, `Interpreter`, `PrevEntryHash`, `WrittenBy string`; `LogEntry` :125–:129 carries `EntryHash` and `TransitionRef`, whose own comment says it "is OUTSIDE the frozen header" |
| V4 | store records object/world/log envelopes; keyed getters, now context-first | `grep -nE 'type (Object\|World\|LogEntry) struct' host/store/store.go; grep -nE 'func \(s \*Store\) (GetObject\|GetWorld\|GetLogEntry\|SelectedHead)' host/store/store.go` | types at :92, :102, :125; getters at :475, :530, :559, :810 — each takes `ctx context.Context` first |
| V5 | no Go grade vocabulary; same-scope Evidence control fires | `test -d . && find . -name '*.go' -type f -print0 \| xargs -0 grep -lE 'gradeOf\|EvidenceGrade\|gradeCode' \| wc -l; find . -name '*.go' -type f -print0 \| xargs -0 grep -l 'Evidence' \| wc -l` | root exists; grade hits `0`; Evidence control `8` files (grew from 1 at the old base — the instrument still sees positives in the same scope) |
| V6 | four grades; SIX-arm total mapping; false TestReport is TESTED; no UNSUPPORTED anywhere under `world/` | `test -d world && grep -rc 'UNSUPPORTED' world/; grep -c 'CLAIMED' world/types.ail; sed -n '15,80p' world/types.ail` | `0` in every file under `world/` (all four `.ail` modules and the compile cache); same-file control `6`; `EvidenceGrade = PROVEN \| TESTED \| ATTESTED \| CLAIMED` exactly four constructors (:34–:38); `Evidence` has SIX constructors including `ProofReceipt` (:23–:29); `gradeOf` is total over all six (:42–:59) with `TestReport(_, _) => TESTED` at :45 (contract) and :54 (body), verdict-blind; test vectors grade `true` and `false` reports both code 3 (:67–:68); no arm yields `PROVEN` |
| V7 | boundary has exactly one Go file and pins it | `test -d host/boundary && find host/boundary -maxdepth 1 -name '*.go' -type f \| wc -l; grep -n 'wantFileCount' host/boundary/allowlist_world_test.go` | `1`; `const wantFileCount = 1` at :1163 with fatal mismatch at :1165 |
| V8 | exact server timeouts and wiring | `sed -n '78,95p' host/daemon/daemon.go; grep -n 'Timeout:' host/daemon/daemon.go` | `readHeaderTimeout = 5 * time.Second` (:81), `readTimeout = 30 * time.Second` (:85), `writeTimeout = 30 * time.Second` (:90), `idleTimeout = 120 * time.Second` (:94); all four wired into `http.Server` at :619–:622 |
| V9 | loopback predicate/refusal has both arms in committed tests | `grep -n 'func isLoopbackHost\|TestIsLoopbackHostMirrorsSketchPredicate\|TestNewRefusesNonLoopbackBind' host/daemon/daemon.go host/daemon/daemon_test.go` | predicate at daemon.go:184; tests at daemon_test.go:99 and :142 |
| V10 | per-group HTTP exception is only `cmd/ailang-worldd`, and a test pins the asymmetry | `grep -n 'extraForbidden\|protectedGoGroups\|func TestBareNetHTTPExemptionIsPerGroup' host/boundary/allowlist_world_test.go` | per-group `extraForbidden []string` field (:39); groups table at :42 with the comment "the ONLY group the exception is true of, hence the empty extraForbidden" (:48); `TestBareNetHTTPExemptionIsPerGroup` at :875 checks `forbiddenImport("net/http", …)` per group in both directions (:897–:901) |
| V11 | Go gate builds and runs plain/race tests | `grep -nE 'go build ./\.\.\.\|go test ./\.\.\.' scripts/verify_go.sh` | `go build ./...` :239, `go test ./... -count=1` :258, `go test ./... -count=1 -race -timeout 8m` :262 |
| V12 | this environment's baseline for the target packages: both green, sockets bind | `GOTOOLCHAIN=go1.25.6 go test ./host/daemon ./host/boundary -count=1` | `ok …/host/daemon 2.128s`, `ok …/host/boundary 1.354s`, rc=0 — REPLACES the 2026-08-13 row: that environment refused `listen(2)`; this one does not, so the daemon package passes whole |
| V13 | AILANG gate baseline | `AILANG_BIN=/tmp/ailang-v0300/ailang scripts/verify_ail.sh` | exit 0; "swept .ail module set equals the LEG1_MODULES allowlist (11 modules)"; "✓ verify gate PASSED: 10 required identities verified, 40 named tests pass"; package gate 9/9 — totals MOVED from the old base's 5/20 |
| V14 | checkout state at revision time; scoped diff against base is empty | `git status --short; git diff --name-only 93e1ba5 -- ':!design_docs'` | status: two controller-staged `design_docs` renames (`w-validated-proven-evidence-boundary*` → `implemented/`) plus modified `design_docs/world-mission.md`, none of them this item's work; the non-`design_docs` diff prints nothing, rc=0 — REPLACES the old "clean worktree" row, which is no longer true of this checkout |
| V15 | exact AILANG gate pins | `grep -n 'EXACT_TOTAL_VERIFIED\|EXACT_TOTAL_TESTS' scripts/verify_ail.sh; sed -n '140,152p' scripts/verify_ail.sh` | `EXACT_TOTAL_VERIFIED=10` (:323); `EXACT_TOTAL_TESTS=40` (:349); `LEG1_MODULES` lists exactly 11 entries (:140–:152) — the old 5/20 pins are gone from the script |
| V16 | item 18 landed; the daemon derives a request-scoped deadline from `r.Context()` | `git branch -r --contains d21754f; grep -n 'r\.Context()' host/daemon/handlers.go` | branch list includes `origin/dev`; one code hit — `:270 return context.WithTimeout(r.Context(), d.readDeadline)` inside the single `readCtx` helper (:269), the other hit :257 being its doc comment — REPLACES "store has no context plumbing" |
| V17 | all six workbench-relevant store reads are context-aware | `grep -n 'ctx context.Context' host/store/store.go host/store/read_object.go` | `GetObject` :475, `GetWorld` :530, `GetLogEntry` :559, `GetRegistryHead` :636, `SelectedHead` :810 in `store.go`; `ReadObject` at `read_object.go:43` — all take `ctx` first; REPLACES "getters are context-free" |
| V18 | `readDeadline` is a real constant bounding GET store reads; expiry is an explicit wire status | `grep -n 'readDeadline' host/daemon/daemon.go; grep -n '"Timeout"' host/daemon/handlers.go` | `readDeadline = 10 * time.Second` (:128), `Daemon` field (:290), wired in `New` (:437), rendered via `writeReadTimeout` (:598); `handlers.go:324` emits `writeAPIError(w, "Timeout", …, http.StatusServiceUnavailable)` — REPLACES "handlers do not consult request context" |
| V19 | a committed ratchet pins deadline-free store reads and covers any new workbench file | `sed -n '370,386p' host/store/context_read_test.go` | `deadlineFreeReadPins = {host/broker/approve.go: 8, host/registry/registry.go: 2, host/replay/replay.go: 1}`; regex matches `context.Background()` flowing into the six getters; `TestNoNewDeadlineFreeStoreReads` walks ALL production `.go` under `host/` and `cmd/` and requires zero for every unpinned file — REPLACES "no production busy timeout" (that claim is also false now; see V22) |
| V20 | the JSON handlers measurably model the malformed/absent distinction (objection B's asked-for row) | `grep -oE 'http\.Status[A-Za-z]+' host/daemon/handlers.go \| sort \| uniq -c; grep -c 'writeAPIError(' host/daemon/handlers.go; grep -c 'func (d \*Daemon) handle' host/daemon/handlers.go` | `10 StatusBadRequest`, `1 StatusConflict`, `1 StatusInternalServerError`, `4 StatusNotFound`, `6 StatusOK`, `1 StatusRequestEntityTooLarge`, `1 StatusServiceUnavailable`; `18` `writeAPIError` sites; `6` handle funcs; `handleWorld`'s own comment states the split as deliberate: "parse the ref at the boundary (malformed -> 400) … distinguish absent from malformed (absent -> 404)" |
| V21 | a 500 never carries a raw error string — refuting the charter row 14 sentence carried forward from the old base | `grep -n 'internalErrorMessage\|"Internal"' host/daemon/handlers.go; grep -n 'err\.Error()' host/daemon/handlers.go; grep -c 'noSuchSymbolXyzzy' host/daemon/handlers.go` | `const internalErrorMessage = "internal store failure"` (:132) is documented as "the ONLY message a 500 ever carries on the wire" (:118); exactly ONE `"Internal"` site (:162) and it passes the constant; all five `err.Error()` survivors are `"BadRequest"`/400 sites (:338, :371, :548, :555, :566); negative control: invented symbol → `0`, rc=1. The charter's row-14 sentence "the `Internal` branches passing `err.Error()` verbatim to an unauthenticated localhost client" (`design_docs/world-mission.md:3069`) described base `9491a10` and is REFUTED at `93e1ba5` — item 18's M3 landed the sanitization |
| V22 | lock waits are `busy_timeout`-bounded (2 s), unlinked from the 10 s deadline — item 22's residual, external to this item | `grep -n 'busyTimeoutMillis' host/store/writer_lock.go; grep -rln 'busyTimeoutMillis' --include='*.go' host/ \| grep -v _test; grep -rln 'readDeadline' --include='*.go' host/ \| grep -v _test` | `const busyTimeoutMillis = 2000` (:181), injected as `_pragma busy_timeout(2000)` (:198); in production code the constant appears ONLY in `host/store/writer_lock.go`, while `readDeadline` appears only in `host/daemon/daemon.go` and `host/daemon/handlers.go` (the `host/archive/archive.go:174` hit is a comment analogy) — no shared file, no cross-reference; `handlers.go`'s `LIMITATION(w-daemon-late-read-503)` comment records both residuals first-party, including "safe only because busy_timeout (2s) is shorter than the 10s request deadline — an ORDERING nothing in this code asserts" |
| V23 | AC6's root assertion fails at base, so it cannot pass vacuously | the exact AC6 compound command | rc=1 at `test -d host/workbench`; `host/workbench` does not exist at `93e1ba5` |
| V24 | repo baseline is green at `93e1ba5` | `GOTOOLCHAIN=go1.25.6 go build ./...; GOTOOLCHAIN=go1.25.6 go vet ./...; AILANG_BIN=/tmp/ailang-v0300/ailang GOTOOLCHAIN=go1.25.6 scripts/verify_go.sh` | build rc=0; vet rc=0; `✓ go gate PASSED: build clean, plain and race tests pass with pinned AILANG_BIN (AILANG v0.30.0 …)`, rc=0 |
| V25 | the inherited deadline tests pass at base; the workbench arm does not exist yet | `GOTOOLCHAIN=go1.25.6 go test ./host/daemon -run 'TestReadCtxCancelledAfterHandler\|TestDaemonReadDeadline' -count=1 -v` | 4 `=== RUN` lines including `TestDaemonReadDeadline/blocking-store`; both top-level tests PASS; `ok … 0.561s`, rc=0; `TestWorkbenchReadDeadline` is absent, which is why AC8 cannot be inherited green |
| V26 | the `readStore` seam exists in `host/daemon/` and has EXACTLY five context-aware methods (objection (B)'s asked-for row, round 3) | `test -d host/daemon; grep -rn 'readStore' host/daemon/*.go | grep -v _test; awk '/type readStore interface/,/^}/' host/daemon/daemon.go; grep -rc 'zzzNoSuchSeamXY' host/daemon/*.go` | root asserted present; `type readStore interface` declared at `host/daemon/daemon.go:323`, held as the field `reads readStore` (:275), documented at :308 as "the daemon's request-read surface: EXACTLY the five store"; the five methods are `GetObject`, `GetWorld`, `GetLogEntry`, `GetRegistryHead`, `SelectedHead`, **each taking `ctx context.Context` as its first parameter**; same-scope positive control: 2 non-test files under `host/daemon/` contain `interface`; negative control (invented symbol, same scope) returns 0 on every file. The doc's §3.2 claim was CORRECT and merely uncited — the reviewer was right that nobody had measured it |

### Non-blocking repository findings

1. `world/types.ail:40`'s comment reads "the ratified five-constructor representation" while
   `Evidence` two screens above it has six constructors (`ProofReceipt` was added since); the
   mapping itself is total over all six (V6). A stale code comment, not a grade-policy defect —
   ordinary maintenance outside this item's scope, and this document does not inherit either
   number from a comment.
2. The mission charter's row 14 still carries "the `Internal` branches passing `err.Error()`
   verbatim to an unauthenticated localhost client" — refuted by measurement at `93e1ba5` (V21).
   The charter needs a one-line correction that this document cannot make.
3. The charter's six-route count and `UNSUPPORTED` instruction remain dead: eight registrations
   (V1) and a four-constructor grade type with no `UNSUPPORTED` anywhere under `world/` (V6).

## 12. Quorum and revision-round decision record

| Round | Reviewer | Verdict | Blocking findings | Resolution |
|---|---|---|---|---|
| Initial quorum → revision 1 | `gpt5-6-sol` | REJECT / BLOCKED | The design falsely treated `WriteTimeout` and cardinality caps as an elapsed-time bound on context-free store reads. Controller measurements (the 2026-08-13 log's rows V16–V19 — this revision REPLACES those row IDs with the post-item-18 measurements; §12.2) confirmed the issue and showed there was no production driver busy timeout either. | Revision 1 deletes the false elapsed-time claims, preserves the caps only as work/response-size bounds, and names accepted residual WB-R1 plus a separately scoped daemon-wide cancellation follow-on in §3.3. No store plumbing or timeout policy is added to this item. |
| Initial quorum → revision 1 | `gemini-3-1-pro` | REJECT / BLOCKED | §2.7 named `LogHeader` fields without direct premise evidence and omitted `EntryIndex`; it could also leave the location of `TransitionRef` ambiguous. | Revision 1 records direct measurement V3b. The premise was substantively confirmed: all five named fields exist, so the EntryView design is retained rather than weakened. `EntryIndex` is rendered once as `entry N`, and the text now states that `EntryHash` and `TransitionRef` sit outside the frozen header. |

### 12.1 Round 2, the park, and the 2026-08-14 ratification

Round 2 ran with all reviewer slots present (no N−1 degrade) and BLOCKED again. The two surviving
objections and how this revision answers them:

**(A) `gpt5-6-sol` — REJECT.** The objection: §3.3 knowingly accepted WB-R1 even though bounded
waits are a mission axiom; cardinality and response-size caps do not bound handler duration,
`WriteTimeout` does not cancel context-free store operations, and the workbench must not proceed
while `GetObject`, `GetWorld`, `GetLogEntry`, and `SelectedHead` can block beyond a defined
request deadline. The reviewer's own `proposed_fix` offered two limbs: grow this item to carry
context-aware store reads and a propagation-killing test, **or defer `/workbench` until the
separately proposed daemon read-cancellation item lands**.

The reviewer was right, and the route taken is the reviewer's own second limb, ratified by Mark
(attended, 2026-08-14) as **Option B**: item 14 does not grow into `host/store` context plumbing;
item 18 `w-daemon-read-cancellation` lands first and gives all eight routes one elapsed-time
contract; residual WB-R1 is discharged by item 18, not by this item. That condition is now
satisfied by measurement, not by argument: item 18's merge is an ancestor of `origin/dev` (V16),
the four methods the reviewer named — and two more — take `ctx` first (V17), a 10 s
request-scoped deadline is derived from `r.Context()` with expiry as an explicit `503 Timeout`
(V16, V18), and committed tests red when propagation or cancellation is removed (V25). §3.3 now
describes that base accurately and binds the workbench to inherit it (AC8, M29, M30). Nothing in
this revision re-litigates the ruling.

**(B) `gemini-3-1-pro` — REJECT.** The objection: §2.4 claimed the workbench "preserves the API's
malformed/absent distinction" while the Verification Log carried no evidence that the JSON
handlers model that distinction or prevent raw error exposure — an unverified premise used to
justify the new route's error posture. The asked-for row now exists: V20 measures the status
posture (400/404/503/500 spread over 18 `writeAPIError` sites in 6 handlers, with the split
stated as deliberate in the handlers' own comments) and V21 measures the raw-error posture
(exactly one `Internal` site, emitting a constant; every `err.Error()` survivor on a 400 path),
and §2.4 cites both. V21 is STRONGER than the concern this document previously inherited from the
charter: charter row 14 still says the `Internal` branches pass `err.Error()` verbatim, and
measurement refutes that at `93e1ba5` — item 18's M3 landed the sanitization, so the charter row
needs correcting (outside this file; §11 finding 2).

### 12.2 What the 2026-08-22 revision (iteration 109) changed

- Header: quorum history, park, and Mark's Option B ratification recorded; the park marked
  discharged; measurement base moved `9491a10` → `93e1ba5`.
- §3.3 rewritten: WB-R1 discharged by item 18 with per-clause measurements; the workbench's
  inheritance obligation (`readCtx`, no second helper, AC8/M29/M30) stated; the two external
  residuals — item 22's lock-wait window and the late-read 200 — named without absorbing either.
- §2.3, §3.2, §4: elapsed-wait language updated from "unbounded" to the measured 10 s bound.
- §2.4: objection (B)'s citations added (V20, V21); the workbench's deadline-expiry `503` posture
  stated.
- §1, §2.4, §3.5: the "seven patterns" route-comment defect claim DELETED — the comment now
  counts eight correctly (V1), so the old sentence was no longer true.
- §1: mapping arm count corrected five → six (`Evidence` gained `ProofReceipt`); `gradeOf`'s
  verdict-blind `TESTED` arm cited at `world/types.ail:45`/`:54` as the rendering hazard §2.6's
  separate-verdict channel exists to surface.
- §5/§6/§7: deadline-propagation gate rows, M29/M30 with their observable argument, and AC8
  added; AC1/AC2/AC6 rebaselined at `93e1ba5`; AC3/AC4 baselined green and labelled regression
  pins (AILANG totals moved 5/20 → 10/40); AC7 rebased to `93e1ba5` and scoped to exclude
  `design_docs`.
- §8: `read_deadline_test.go` added to the modified-files table; the `readCtx` helpers and
  `busyTimeoutMillis` added as inspected-not-moved pins; verify_ail pin rows updated to 10/40.
- §11: every row re-run first-party at `93e1ba5`; V12 and V14 replaced (their old claims were
  environment- and worktree-specific and are no longer true); V16–V19 replaced (the store now HAS
  context plumbing, handlers DO consult the request context, and the production DSN DOES set a
  busy timeout); V20–V25 added.


### 12.3 Round 3 (2026-08-22, iteration 109) — re-quorum, the restored reviewer, and the carve-out

The revision above was re-quorumed on 2026-08-22. **Verdict: BLOCKED, and the synthesis was
missing the reviewer it most needed.** `gemini-3-1-pro` was present and rejected;
**`gpt5-6-sol` was ABSENT with reason `budget`** — it refused because the document had GROWN, and
it had grown answering *that reviewer's own* round-2 objection. Acting on the synthesis as printed
would have taken a verdict with a named hole in exactly the shape the mission's rule predicts:
the reviewer drops out precisely when its opinion is most load-bearing. It was therefore re-run
alone at a raised cap (`ailang design-review --reviewer gpt5-6-sol --max-cost-usd 0.60`, rc=0,
`present=true`, $0.089025).

**The restored reviewer did NOT repeat its round-2 objection.** The unbounded-store-wait objection
is gone from its verdict — the sequencing answer worked, and item 18's landed base is what
answered it. Instead it raised a NEW and real defect that the budget-absent synthesis would have
concealed entirely:

> "The query contract explicitly says unknown workbench parameters are ignored, which is a silent
> fallback and directly violates the mission axiom requiring no silent fallbacks. A typo such as
> `?paylod=1` silently produces a different view instead of refusing the unsupported request."

That was true of §2.4 as revised, and it is this mission's own headline discipline turned on this
document.

**Both surviving objections satisfy the narrow-refinement carve-out**: each carries a concrete
reviewer-authored `proposed_fix`, and neither disputes the design DIRECTION (one is completeness /
attribution, the other determinism / refusal-strictness — nothing challenges the renderer route,
the read-only posture, or the sequencing behind item 18). The controller therefore applied the
reviewers' **verbatim** text as a bounded third revision, overriding nothing:

- **`gpt5-6-sol` (determinism).** §2.4's closed-grammar paragraph is the reviewer's replacement
  text applied word-for-word. Its named arms `TestWorkbenchRefusalBranches/unknown-parameter` and
  `/duplicate-parameter` are added to §5, its asked-for mutation rows are M31/M32 in §6, and AC1's
  fail-trigger range is widened to cover them.
- **`gemini-3-1-pro` (completeness).** V26 is the row it asked for, and §3.2 now cites it. The
  measurement CONFIRMS the document: `type readStore interface` is declared at
  `host/daemon/daemon.go:323` with exactly five methods, every one context-first. The reviewer was
  right that nobody had measured it, and wrong that the claim might be false — which is why the
  controller measured the premise rather than forwarding it to another revision round.

Metered spend across round 3: **$0.127437** (quorum $0.038412 + restored reviewer $0.089025) of the
$5 iteration ceiling.

## Related

- [The Human Surface](../HUMAN-SURFACE.md) — binding workbench grammar and anti-patterns
- [Human Interaction Scenarios](../SCENARIOS.md) — Scenario 3 provenance walk
- [AILANG World Design](../DESIGN.md) — §1 thesis and §14 hard boundaries
- [Evidence grade mapping](../implemented/w-evidence-grade-mapping.md) — canonical mapping and
  deliberate absence of `PROVEN` mint authority
- [worldd M2](../implemented/w-worldd-m2.md) — daemon transport and bounded-wait decisions
- [Daemon read cancellation](../implemented/w-daemon-read-cancellation.md) — item 18, the landed
  base that discharges WB-R1: `readCtx`, `readDeadline`, context-first getters, `503 Timeout`
