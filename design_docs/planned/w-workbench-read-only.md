# w-workbench-read-only — a truthful localhost lens over the immutable world

- **Status**: Planned — awaiting quorum
- **Item**: queue item 14, `w-workbench-read-only`, clause-5
- **Filed**: 2026-08-11, human attended
- **Designed**: 2026-08-13, iteration 82
- **Estimated**: 1.5–2 days
- **Measurement base**: `9491a10`
- **Files changed by implementation**: enumerated in §8; this document changes only itself
- **Design result**: one server-rendered, read-only localhost route; no SSE, no write verbs,
  no copied grade policy, and no claim that the present log is already a complete provenance graph

Every present-tense codebase statement is backed by a command and observed output in §11.
Controller-supplied facts are labelled there and were re-run where they carry this design.

---

## 1. Problem

The binding human surface says the World, not a transcript, is the workspace. Its reference
chrome is a world browser, timeline scrubber, and provenance explorer. Scenario 3 makes that
promise concrete: begin at a deployment, walk backward through the deploying transition,
proposal, agent, contemporaneous evidence, and originating goal, then compare now with the
world immediately before the surprising transition.

The daemon is the right process boundary for a local workbench, but its measured surface is
narrower than the charter prose suggests. It has seven GET patterns and one POST pattern, not
six or seven total (V1). None serves HTML, SSE, embedded assets, or a file server (V2). The
existing comment above the route table says “seven patterns” although eight follow; that is a
non-blocking repository defect, not a premise to preserve (V1).

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
only in AILANG. No Go source names `gradeOf`, `EvidenceGrade`, or `gradeCode` (V5). Copying its five
arms into a Go switch would create a second, unproved policy and the exact grade-laundering surface
the UX forbids. The stale queue instruction to print `UNSUPPORTED` for unmapped variants cannot be
followed: the exported type has four constructors and no `UNSUPPORTED` constructor (V6).

Finally, `TestReport(ref, false)` canonically has grade `TESTED`, exactly like the true case (V6).
That grade describes the kind of evidence, not its verdict. A badge saying only `TESTED` over a
failed test would be materially misleading. Grade and verdict must be separate channels.

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
  owns: route registration, query parsing, status codes, cardinality-bounded store reads
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
- A full page request has bounded work cardinality and response size; it does not have a bounded
  elapsed wait because the store API has no context plumbing or production SQLite busy timeout
  (V16–V19).
- Progressive enhancement may later add local JavaScript, but it is not an acceptance dependency.

Reject “JSON API plus static assets” for this item. The daemon already has the necessary JSON
reads. A new aggregate JSON API would become another public contract before the typed projection
exists; static assets would add caching and content-type gates without improving the first
read-only grammar.

Reject SSE. The server’s `WriteTimeout` configures a 30-second response-write deadline (V8); it is
not a handler or store-call cancellation mechanism (V16–V19). Keeping a long-lived response open
would require deliberate deadline handling through `http.ResponseController`, which this item does
not need because the log is immutable and explicit refresh is sufficient. This design neither uses
`ResponseController` nor changes any timeout.

### 2.4 Route and version posture

`/v1/*` remains the machine API. `/workbench` is an operator renderer whose HTML may evolve; it is
not inserted into the frozen v1 table and it does not change the semantics of the existing eight
patterns. The route table becomes nine patterns: eight existing patterns plus `GET /workbench`.
The stale “seven patterns” comment is corrected to state the exact machine/API and renderer split.

Unknown workbench query parameters are ignored only if they do not select data. Malformed values
for `world`, `object`, `from`, or `entry` return `400` with a small HTML error page. Well-formed but
absent references return `404`. Store failures return `500`. Those branches preserve the API’s
malformed/absent distinction without exposing raw SQL or stack traces.

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

The adapter calls the same store methods the existing JSON handlers call. It does not call the JSON
handlers through an in-process HTTP client; that would add encoding work and blur error ownership.
The adapter does not expose `d.store` to the presentation package.

Timeline reads stop after 100 entries or the first missing index. Query `from` is non-negative and
addition is overflow-checked. Payload preview is off by default and capped at 64 KiB with a visible
“truncated” marker. The workbench cannot request the daemon API’s maximum 500-entry page in one HTML
render; rendering work and response size remain independently bounded. These caps do not bound
elapsed store wait time.

Add `host/daemon/workbench_test.go`. Tests use `httptest.NewRecorder` and the constructed handler,
not `httptest.NewServer`, so the route tests are hermetic in restricted environments. They seed real
store objects and assert status, headers, escaping, links, page bounds, opt-in payload, and every
refusal branch.

### 3.3 Named residual: daemon reads have no elapsed-time bound

This item accepts **WB-R1: a workbench handler can remain blocked in a store read after the HTTP
write deadline**. `http.Server.WriteTimeout` bounds response writes, not handler execution. The
store has no `context.Context` plumbing, every relevant getter is context-free, daemon handlers do
not consult the request context, and the production store has no SQLite `busy_timeout` (V16–V19).
The 100-entry and 64 KiB caps bound work and response size only; they never bound wait time.

This is a pre-existing property of the daemon's JSON read path as well as the new renderer. Adding
context-aware store operations or a production lock-acquisition timeout is cross-cutting daemon/store
work and is not absorbed into this read-only 1.5–2 day item. The required follow-on is a separately
designed daemon read-cancellation item that selects a store/driver mechanism, gives existing JSON
reads and `/workbench` the same elapsed-time contract, and supplies a blocking-store test that reds
when cancellation or the driver timeout is removed. No acceptance criterion in this item treats
`WriteTimeout`, the page cap, or the payload cap as an elapsed-time observable.

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

### 3.5 Documentation correction

Correct the daemon route-table comment while adding the ninth route. The comment must say there are
eight frozen `/v1` machine patterns (seven GET, one POST) plus one unversioned read-only workbench
pattern. This is a comment correction, not a route deletion or v1 redesign.

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
It also does not buy elapsed-time cancellation for store reads; WB-R1 and its daemon-wide follow-on
are explicit in §3.3.

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
| payload is opt-in and capped | `TestWorkbenchPayloadPreviewBound` | small payload renders in full |
| hostile values are escaped | `TestRenderEscapesAllObjectText` | safe literal remains visible |
| links stay local and relative | `TestRenderEmitsOnlyLocalLinks` | expected `/workbench?...` link exists |
| missing edge is explicit | `TestRenderUnavailableProvenanceEdge` | stored edge renders as a link |
| grade is never synthesized | `TestGradeViewRejectsUnavailableOrInvalidInput` | each of four valid labels is accepted |
| failed test cannot hide behind TESTED | `TestGradeViewRequiresTestVerdict` | TESTED+PASS and TESTED+FAIL both render both channels |
| renderer has no transport authority | `TestWorkbenchPackageRemainsTransportFree` in boundary gate | daemon closure still contains exactly one `net/http` |
| loopback refusal remains two-sided | existing `TestNewRefusesNonLoopbackBind` | its accepted loopback table remains |
| server timeout constants remain configured | existing `TestBoundedWaitsAndBodyLimit` | exact non-zero constants remain; this does not prove handler/store cancellation |
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

M25 requires the boundary test to exercise its own detector with an overlay, following the file’s
existing mutation-test idiom. Merely asserting the dependency count once is not enough.

## 7. Acceptance criteria

Every command below was baselined on the unmodified `9491a10` tree. “Fail trigger” states how the
criterion can red. “Producible” states why the named mechanism can emit the observable.

### AC1 — package and route behavior

Command:

```sh
go test ./host/workbench ./host/daemon -count=1
```

Baseline: `host/workbench` does not exist, so the future exact command is not runnable at HEAD;
the existing daemon package was separately run and failed only in socket-binding tests because
this design environment forbids `listen(2)` (V12). Implementation tests named in §5 must use
recorders and buffers, so their observables are producible without sockets.

Pass: all new renderer and handler tests pass. Fail trigger: any mutation M1–M23 or a compile
failure reds a named test/package. Producible: store fixtures, `httptest.ResponseRecorder`, and
`bytes.Buffer` directly generate every asserted status, header, link, and body.

### AC2 — boundary authority remains explicit

Command:

```sh
go test ./host/boundary -run 'Test(BareNetHTTPExemptionIsPerGroup|WorkbenchPackageRemainsTransportFree)' -count=1
```

Baseline: existing `TestBareNetHTTPExemptionIsPerGroup` passes; the new named test is absent at
HEAD, so the exact combined selector presently exercises only the existing arm and is not accepted
as implementation evidence (V12). Pass requires test enumeration/log output proving both names ran.

Fail trigger: adding transport to `host/workbench`, losing daemon `net/http`, adding a second
exception, or disarming the overlay mutant. Producible: `go list -deps -json` and the existing
overlay/read helpers expose package closures and import substitutions in this test package.

### AC3 — full Go gate

Command:

```sh
AILANG_BIN=/tmp/ailang-v0300/ailang scripts/verify_go.sh
```

Baseline: not run in this restricted design environment because it includes tests that bind
loopback sockets and a race leg; direct daemon testing demonstrates the sandbox bind failure (V12).
This is an environment limitation, not a green baseline. The controller/CI environment must first
record a green unmodified-tree baseline, then run the same command after implementation.

Fail trigger: build, plain tests, race tests, binary hygiene, toolchain canary, or AILANG pin fails.
Producible: the committed script invokes each observable and is the repository’s existing Go CI
gate (V11). A sprint may not waive this AC on the strength of AC1.

### AC4 — AILANG gate totals do not move

Command:

```sh
AILANG_BIN=/tmp/ailang-v0300/ailang scripts/verify_ail.sh
```

Baseline at HEAD: exit 0; 11 allowlisted modules, 5/5 required verified identities, 20 required
named tests, and the nine-step package gate pass (V13).

Pass: the exact same totals and package gate pass after implementation. Fail trigger: any AILANG
module, contract identity, required test, package projection, manifest, tar member, or golden moves.
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

Baseline: `host/workbench` is absent, so this exact command fails at the root assertion and cannot
pass vacuously; same-scope daemon control finds `http.Server`/`net/http` in V2/V8. Pass after
implementation is root assertion success, zero forbidden-pattern hits, and a positive HTTP control
in the same call.

Fail trigger: SSE, deadline relaxation, streaming/file-server, or embedded-asset machinery appears.
Producible: the chosen HTML renderer needs none of those symbols. Persistence comes from imports
and behavior tests; this grep is a review assertion, not the sole gate.

### AC7 — only priced files changed

Command:

```sh
git diff --name-only 9491a10 --
```

Baseline: empty at the start of design work (V14). For implementation, pass is a subset of the
implementation files enumerated in §8, excluding this already-landed design doc from sprint edits.
Fail trigger: any unpriced path appears. Producible: the implementation is fully contained by the
listed Go source/tests and the existing boundary gate file.

This snapshot does not persist after merge; code ownership and ordinary review enforce scope.

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
| `host/daemon/daemon.go` | one GET registration; correct eight-pattern comment to eight v1 plus workbench |
| `host/boundary/allowlist_world_test.go` | transport-free package guard and daemon positive arm; retain one-file pin and sole exemption |

### 8.3 Files inspected but not moved

| File/pin | Required disposition |
|---|---|
| `host/boundary/allowlist_world_test.go` `wantFileCount = 1` | remains exactly 1; no new boundary Go file |
| `protectedGoGroups` four-row table | remains four rows; `cmd/ailang-worldd` remains sole bare-http exception |
| `scripts/verify_ail.sh` `LEG1_MODULES` | remains the 11-entry set |
| `scripts/verify_ail.sh` `EXACT_TOTAL_VERIFIED=5` | remains 5 |
| `scripts/verify_ail.sh` `REQUIRED_TESTS` | remains its current 20 identities |
| `scripts/verify_ail.sh` `EXACT_TOTAL_TESTS = 20` | remains 20; note Python spacing |
| `scripts/verify_go.sh` | no edit; it already discovers all Go packages/tests |
| `world/types.ail` | no edit; remains canonical grade policy |
| `host/store/*` | no edit; no reverse index or query added |
| `design_docs/HUMAN-SURFACE.md` | binding input, not edited |
| `design_docs/SCENARIOS.md` | binding Scenario 3 input, not edited |

The likely collision is queue item 7 adding approval-inbox views to `host/workbench` and another
route in `host/daemon`. It must compose the view model rather than replace grade/refusal rules.
Transition/provenance work may later add a typed projection; it must not teach this renderer to
parse opaque payloads heuristically.

## 9. Scope and pricing

The 1.5–2 day price covers:

- 0.35 day: view model, template, accessibility, empty/unavailable states;
- 0.35 day: daemon query adapter and cardinality-bounded reads;
- 0.35 day: route/refusal/security/payload tests;
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
- choose or implement timeout policies;
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

All commands ran from repository root at `9491a10` on 2026-08-13. Rows marked controller-supplied
were re-run before use. Empty-result rows include a root assertion and same-scope positive control.

| ID | Codebase claim | Command | Observed output |
|---|---|---|---|
| V1 | Handler has 8 patterns: 7 GET, 1 POST; comment miscounts 7 | `grep -c 'mux.HandleFunc' host/daemon/daemon.go; sed -n '456,470p' host/daemon/daemon.go` | `8`; seven shown `GET` registrations and one `POST /v1/commit`; comment says “seven patterns” (controller-supplied, re-verified) |
| V2 | no current production Go HTML/SSE/file/embed route in `host/`+`cmd/`; instrument sees HTTP in same scope | `test -d host && test -d cmd && find host cmd -name '*.go' ! -name '*_test.go' -type f -print0 | xargs -0 grep -lE 'text/html|text/event-stream|http.FileServer|embed.FS' | wc -l; grep -c 'http' host/daemon/daemon.go` | roots exist; `0`; positive control `26` (controller-supplied, re-verified) |
| V3 | exact response fields and bounded log composition | `sed -n '30,70p;194,360p' host/daemon/handlers.go` | world fields `ref,revision,stateRoot,logHead`; object fields `hash,interfaceHash,semanticId,provenance,payload`; log fields header, `entryHash,transitionRef`; clamp default 100/max 500 |
| V3b | `LogHeader` directly contains all six frozen fields; `LogEntry` contains `EntryHash` and `TransitionRef` outside that header | `grep -n "type LogHeader" -A 20 host/store/store.go` | `LogHeader`: `EntryIndex int64`, `SemanticsEpoch int64`, `TransitionFn hashref.HashRef`, `Interpreter hashref.HashRef`, `PrevEntryHash hashref.HashRef`, `WrittenBy string`; `LogEntry`: `Header LogHeader`, `EntryHash hashref.HashRef`, `TransitionRef hashref.HashRef`; intervening comment says the header fields are frozen and `LogEntry` may grow |
| V4 | store records object/world/log envelopes and only keyed getters for these reads | `rg -n 'type (Object|World|LogEntry)|func \(s \*Store\) (GetObject|GetWorld|GetLogEntry)' host/store/store.go` | types at 90,100,123; getters at 467,522,551 |
| V5 | no Go grade vocabulary; same-scope Evidence control fires | `test -d . && find . -name '*.go' -type f -print0 | xargs -0 grep -lE 'gradeOf|EvidenceGrade|gradeCode' | wc -l; find . -name '*.go' -type f -print0 | xargs -0 grep -l 'Evidence' | wc -l` | root exists; grade hits `0`; Evidence control `1` (controller-supplied, re-verified) |
| V6 | four grades, total five-arm mapping, false TestReport is TESTED, no UNSUPPORTED | `test -d world && grep -rn 'UNSUPPORTED' world/ | wc -l; grep -c 'CLAIMED' world/types.ail; sed -n '27,68p' world/types.ail` | `0`; control `4`; four-grade ADT; five arms; true and false TestReport tests both code 3 (controller-supplied, re-verified) |
| V7 | boundary has exactly one Go file and pins it | `test -d host/boundary && find host/boundary -maxdepth 1 -name '*.go' -type f | wc -l; sed -n '1158,1168p' host/boundary/allowlist_world_test.go` | `1`; `const wantFileCount = 1` and fatal mismatch (controller-supplied, re-verified) |
| V8 | exact server timeouts and wiring | `sed -n '73,94p;514,522p' host/daemon/daemon.go` | 5s header, 30s read, 30s write, 120s idle; all four wired into `http.Server` |
| V9 | loopback predicate/refusal has both arms in committed tests | `rg -n 'func isLoopbackHost|TestIsLoopbackHostMirrorsSketchPredicate|TestNewRefusesNonLoopbackBind' host/daemon/daemon.go host/daemon/daemon_test.go` | predicate at daemon.go:171; tests at daemon_test.go:88 and :130 |
| V10 | per-group HTTP exception is only `cmd/ailang-worldd` | `sed -n '30,65p;870,905p' host/boundary/allowlist_world_test.go` | four groups; store/replay/world-publish forbid bare HTTP; daemon command has nil extra list; test pins asymmetry (controller-supplied, re-verified) |
| V11 | Go gate builds and runs plain/race tests | `grep -nE 'go build ./\.\.\.|go test ./\.\.\.' scripts/verify_go.sh` | build, plain test, and bounded race-test commands present |
| V12 | design-environment baseline for target packages | `go test ./host/daemon ./host/boundary -count=1` | daemon fails when socket tests get `bind: operation not permitted`; boundary passes `ok` |
| V13 | AILANG pins and baseline | `AILANG_BIN=/tmp/ailang-v0300/ailang scripts/verify_ail.sh` | exit 0; 11 modules; `5/5`; 20 tests; package gate 9/9 |
| V14 | worktree initially had no changes | `git status --short` | empty output; positive repository/path controls are V1 and V7 |
| V15 | exact AILANG gate pins | `sed -n '135,153p;306,346p' scripts/verify_ail.sh` | 11-entry `LEG1_MODULES`; shell `EXACT_TOTAL_VERIFIED=5`; Python `REQUIRED_TESTS`; `EXACT_TOTAL_TESTS = 20` (controller-supplied, re-verified) |
| V16 (controller M1) | store has no context plumbing; same-file method control fires | `grep -c "context.Context" host/store/store.go; grep -c "func (s \*Store)" host/store/store.go` | `0`; same-file known-positive control `14` (controller-supplied measurement) |
| V17 (controller M2) | all workbench-relevant read getters are context-free | `grep -nE 'func \(s \*Store\) (GetObject|GetWorld|GetLogEntry|SelectedHead)' host/store/store.go` | `467 GetObject(ref hashref.HashRef)`, `522 GetWorld(ref hashref.HashRef)`, `551 GetLogEntry(index int64)`, `802 SelectedHead()`; none accepts `context.Context` (controller-supplied measurement) |
| V18 (controller M3) | daemon handlers do not consult request context | `grep -n "r.Context()\|context\." host/daemon/handlers.go` | zero hits (controller-supplied measurement) |
| V19 (controller M4) | production store code has no SQLite busy timeout | `grep -rn "busy_timeout" host/store/*.go \| grep -v _test` | zero hits; known-positive control exists only in `host/store/writer_lock_test.go:609` test DSN (controller-supplied measurement) |

### Non-blocking repository findings

1. `host/daemon/daemon.go` calls eight registrations “seven patterns.” Implementation should fix
   the comment while accurately separating eight frozen `/v1` patterns from the new renderer.
2. The charter row’s six-route count and `UNSUPPORTED` instruction are stale. Neither is treated
   as an implementation requirement because measured code and the ratified grade type supersede it.

## 12. Quorum and revision-round decision record

| Round | Reviewer | Verdict | Blocking findings | Resolution |
|---|---|---|---|---|
| Initial quorum → revision 1 | `gpt5-6-sol` | REJECT / BLOCKED | The design falsely treated `WriteTimeout` and cardinality caps as an elapsed-time bound on context-free store reads. Controller measurements V16–V19 confirmed the issue and showed there is no production driver busy timeout either. | Revision 1 deletes the false elapsed-time claims, preserves the caps only as work/response-size bounds, and names accepted residual WB-R1 plus a separately scoped daemon-wide cancellation follow-on in §3.3. No store plumbing or timeout policy is added to this item. |
| Initial quorum → revision 1 | `gemini-3-1-pro` | REJECT / BLOCKED | §2.7 named `LogHeader` fields without direct premise evidence and omitted `EntryIndex`; it could also leave the location of `TransitionRef` ambiguous. | Revision 1 records direct measurement V3b. The premise was substantively confirmed: all five named fields exist, so the EntryView design is retained rather than weakened. `EntryIndex` is rendered once as `entry N`, and the text now states that `EntryHash` and `TransitionRef` sit outside the frozen header. |

## Related

- [The Human Surface](../HUMAN-SURFACE.md) — binding workbench grammar and anti-patterns
- [Human Interaction Scenarios](../SCENARIOS.md) — Scenario 3 provenance walk
- [AILANG World Design](../DESIGN.md) — §1 thesis and §14 hard boundaries
- [Evidence grade mapping](../implemented/w-evidence-grade-mapping.md) — canonical mapping and
  deliberate absence of `PROVEN` mint authority
- [worldd M2](../implemented/w-worldd-m2.md) — daemon transport and bounded-wait decisions
