# w-world-live-surface — rows 148–151

- Status: **DRAFT — authored attended 2026-10-06 (Mark + Claude), routed for v1.0.0 by
  `D-WORLD-64`.** Not yet quorum-reviewed. The loop's designer takes this as its input: re-measure
  §2 at the pick-time base, run the pick-time quorum, then plan. Nothing here waives a gate.
- **Row 148 is SUPERSEDED here by [`../implemented/w-worldd-agui-event-stream.md`](../implemented/w-worldd-agui-event-stream.md)** (iteration 243, landed `38db793`): bounded 18 s runs + resume (no 30-min connection, frozen D7 `writeTimeout`), resume cursor from AG-UI `state.lastIndex` or `Last-Event-ID`, keyset tail, no session filtering. Rows 149 and 151 consume that interface.
- **Row 149's parts (D6, D7, §4 row 149, §6 149(a)–(e)) are SUPERSEDED by [`w-workbench-live-and-polish.md`](../implemented/w-workbench-live-and-polish.md)** (iteration 244/245 pick-time design, base `b702d73`). Rows 150–151 stay here.
- Measurement base: `2d53d72` (`origin/dev`, iteration 235 record). All measurements are listed in §10.
- Rows: **148** `w-worldd-agui-event-stream` · **149** `w-workbench-live-and-polish` ·
  **150** `w-marketing-capture` · **151** `w-agent-generated-projections`. One design, four rows,
  landed in that order. Each row is an independently mergeable sprint.
- Binding inputs: [HUMAN-SURFACE.md](../HUMAN-SURFACE.md) (P1–P7, §2.5, §3.1, §5 anti-patterns,
  §6 build strategy, all ratified), [DESIGN.md](../DESIGN.md) §3.7 / §11 / open question 9,
  [coding-standards.md](../coding-standards.md) S2, S3, S6, S7.

## 1. Problem and clause mapping

World is about to ship 1.0 and has no way to show itself. Three things are missing.

1. **Nothing is live.** worldd serves the workbench as script-free HTML (CSP
   `default-src 'none'`) and has no event stream (V1, V2). A human has to reload to see a
   commit. P5 (attention choreography, "quiet is health") cannot be rendered by a page that
   doesn't know when the world changed.
2. **Nothing is presentable.** The workbench was built to the ratified grammar (grades,
   checked links, closed query grammar, items 14 and 98–105). It is correct, but visually it is
   an inspector. There is no seeded demo world and no reproducible way to produce release images.
3. **Agents cannot make UI.** DESIGN.md §11.2 says domain UIs are *generated projections* emitted
   in an open agent-UI protocol (M6, open question 9). Mark has named this a headline feature:
   an agent working through World can build the view a human needs, and the view carries
   World's grades and provenance like every other object.

**Clause mapping, stated honestly.** None of these rows closes an UNMET bar clause on its own.
Row 148 strengthens clause 6 (protocol-native boundary: one more open protocol, no new wire
protocol). Rows 149 and 151 strengthen clause 5 (the human surface works). Row 150 moves no
clause. It is release work. They are on the 1.0 roadmap because Mark ruled the v1.0.0
release ships with them (`D-WORLD-64`), not because the bar requires them.

**Godot is ruled out** for World's surface (`D-WORLD-64`). Reason: P1–P3 are text and typed
values, and HUMAN-SURFACE §6 binds the reference renderer to a browser-opened, zero-install
surface served by worldd. A game engine does that badly. Stapledon's Voyage keeps Godot.

## 2. Measured current state (base `2d53d72`)

1. **Routes.** `Daemon.Handler` (`host/daemon/daemon.go:800-821`) registers the ten frozen
   `/v1` routes, `GET /workbench`, the A2A card, `POST /a2a/` and `POST /mcp/`. There is no
   event or stream route (V1).
2. **SSE.** `text/event-stream` appears in this repo only in MCP tests and the CLI's MCP client
   (`cmd/ailang-worldd/mcpclient.go`). The MCP server's SSE lives upstream in
   `ailang/serveapi/protocol/mcphttp` and scopes one response to one request. worldd has no
   long-lived push (V2).
3. **Workbench CSP.** `workbenchCSP` (`host/daemon/workbench.go:20`) is
   `default-src 'none'; style-src 'unsafe-inline'; base-uri 'none'; form-action 'self'; frame-ancestors 'none'`.
   It allows no scripts and no `connect-src`. Making the page live requires a deliberate CSP
   change, and `TestWorkbenchSecurityHeaders` pins this constant (V3).
4. **Read authority.** `isProtected` (`daemon.go:828-830`) protects only `POST /v1/commit`. The
   GET routes are unauthenticated under declared residual R1 (V1).
5. **Log reads.** `Store.GetLogEntry(ctx, index)` (`host/store/store.go:687`) reads one entry by
   index. `Store.SelectedHead` (`:988`) reads the head. `LogHeader` carries `EntryIndex`,
   `SemanticsEpoch`, `TransitionFn`, `Interpreter`, `PrevEntryHash` and `WrittenBy`
   (`:123-130`). A stream can be built entirely from existing reads. No store change is
   needed (V4).
6. **No demo seed.** `ailang-worldd` subcommands are serve, health, head, world, object, log,
   registry, commit, tools, call, why, provenance, doctor, setup, session and help
   (`cmd/ailang-worldd/main.go:238-302`). None seeds a demo world (V5).
7. **Dependencies.** `go.mod` requires only `ailang` and `modernc.org/sqlite`, and boundary
   allowlist tests exist (`host/boundary/allowlist_world_test.go`). Any Go dependency these rows
   add is a slim-kernel (S3) decision. The design adds none (V6).
8. **Mockups.** `design_docs/mockups/approval-inbox.html` and `grounded-prose.html` define the
   visual language: tokens on `:root`, light and dark themes, cards, grade pills. They are the
   style reference for row 149 (V7).
9. **Decision packets are typed, not yet produced.** `DecisionPacket`
   (`world/decision-packet/v1`) is defined in `world/types.ail:155-159`. Whether any host path
   commits one at pick time must be re-measured. If none does, row 148 ships the
   `world.decision.*` mapping with a fixture test, and the live events appear when the
   approval-inbox path commits packets (V10).

**External facts (fetched 2026-10-06; re-check at pick time, both specs are moving):**

- **AG-UI** ([docs.ag-ui.com](https://docs.ag-ui.com/concepts/architecture)) is an
  event protocol carried over SSE. The client POSTs `RunAgentInput` (`threadId`, `runId`, state,
  messages, tools), and the server streams typed events: lifecycle (`RUN_STARTED`,
  `RUN_FINISHED`, `RUN_ERROR`, `STEP_*`), text, tool calls, state (`STATE_SNAPSHOT`,
  `STATE_DELTA` as RFC 6902 JSON Patch, `MESSAGES_SNAPSHOT`), `RAW` and namespaced `CUSTOM`.
- **A2UI** ([a2ui.org v1.0 RC](https://a2ui.org/specification/v1.0-a2ui/);
  [v0.9 announcement](https://developers.googleblog.com/a2ui-v0-9-generative-ui/)) is a
  declarative UI format. An agent sends `createSurface`, `updateComponents`, `updateDataModel`
  and `deleteSurface`. Components are restricted to a **renderer-declared catalog**
  (JSON Schema), so no executable code crosses the wire. Data binds by JSON Pointer, and user
  actions return as `action` messages. v0.9.1 is the production line, and v1.0 is a release
  candidate with stable targeted for Q4 2026. A2UI is documented as running over AG-UI.

## 3. Decisions, with alternatives

### D1 — Open question 9 is answered "both, layered": AG-UI is the transport, A2UI is the payload

The two protocols are not competitors. AG-UI carries *what happened*, and A2UI describes
*what to draw*. World emits AG-UI for every live surface (row 148) and carries A2UI documents
inside it for agent-generated views (row 151).

*Alternatives:* AG-UI alone leaves agents no declarative UI format, so they would invent one,
which violates "no new wire protocols". A2UI alone has no world-event transport. A homegrown
event JSON is ruled out by §3.7.
*Cost:* both specs are pre-1.0. Mitigation: pin the spec versions in a constant, and keep each
mapping in one file so a spec bump is a one-file change.

### D2 — The event stream is a projection of the log, and it resumes by log index

Each event is derived from a committed log entry and its referenced objects. Nothing is
emitted that isn't in the log. The SSE `id:` field is the entry index, so a client that
reconnects with `Last-Event-ID: n` (or `?from=n`) receives exactly the events for entries
`n+1…head`. The same entries always produce the same events, so **the stream is replayable,
and a test can assert the stream bytes for a fixture log**. This is World's determinism
carried into the UI.

The handler **tails the log by index**: it reads `SelectedHead` and pages with `GetLogEntry`,
polling on a short bounded tick (start at 250 ms; the plan measures it). It does not use an
in-process pub/sub hooked into each commit path.

*Alternatives:* an in-process broadcaster on `Store.Commit` would have to hook every commit path
(REST, MCP, A2A, coordinator) and misses writes by a second process on the same SQLite file.
SQLite has no cross-process change notification. Polling one indexed read is cheap, correct
across processes, and trivially bounded.
*Cost:* up to one tick of latency, which is invisible at human speed.

### D3 — The route is `POST /agui/`, AG-UI-conformant, with World semantics carried in CUSTOM events

- **Route.** `POST /agui/` accepts `RunAgentInput`. `threadId` names the world (the store's
  selected head chain), and `runId` is client-chosen. The response is `text/event-stream`. This
  is the AG-UI standard path, so stock AG-UI clients work unmodified. The workbench reads it
  with `fetch` and a streamed body (EventSource is GET-only). It is additive and outside the
  frozen `/v1` table, the same pattern as `/a2a/` and `/mcp/`.
- **Run lifecycle.** The handler emits `RUN_STARTED`, then one `STATE_SNAPSHOT`, then per-entry
  events, until the client disconnects or the per-connection budget (D5) ends it with
  `RUN_FINISHED`. Errors become `RUN_ERROR` with the daemon's existing error classes.
- **Event mapping (v1, frozen per D1's version constant):**

  | World fact | AG-UI event |
  |---|---|
  | head, revision, log head, entry count, pending-decision count | `STATE_SNAPSHOT` at start; `STATE_DELTA` (JSON Patch) on change |
  | one committed log entry | `CUSTOM` `world.entry.committed` `{index, entryHash, transitionFn, writtenBy, epoch}` |
  | a decision packet created, decided or timed out (`world/decision-packet/v1`) | `CUSTOM` `world.decision.*` carrying the packet ref and its grades |
  | an effect intent or outcome recorded | `CUSTOM` `world.effect.*` |
  | an agent-generated view (row 151) | `CUSTOM` `world.view.a2ui` carrying A2UI messages |

  `CUSTOM` names are namespaced `world.` as AG-UI requires. Text and tool-call events are
  **not** emitted in v1: World's narrative output goes through grounded prose (P2), and that
  binding is row 151's follow-up, not this row's.

*Alternatives:* `GET /v1/events` as plain SSE is simpler but is a new dialect (a §3.7
violation) and no AG-UI client could consume it. WebSocket adds a dependency and a
second transport for no gain.

### D4 — Authority: same read posture as `/v1` GET today, filtered by session when a credential is presented

Without a credential, the stream carries exactly what the unauthenticated GET routes already
expose (residual R1). It adds nothing new to residual R1. With a session credential, the stream
is filtered to what the session may read, using the projection's existing session resolution.
**Decision packets are never actionable from the stream.** A decision is a transition, and it
goes through the existing authenticated commit path (HUMAN-SURFACE §5: no admin backdoor).

When the follow-up row closes residual R1, the stream inherits the change by construction, because it
calls the same `isProtected` predicate.

### D5 — Bounds

The stream enforces a per-connection budget, a global connection cap and a per-tick read cap.
The daemon's read deadline applies to every store read inside the loop. Initial numbers are
30 min per connection, 16 connections, 100 entries per tick. These are placeholders: the
plan measures them and the frozen D7 handler deadlines are respected per read. A slow client
is dropped rather than buffered without limit.

### D6 — The workbench goes live with one vendored, self-hosted script, and the CSP widens by exactly two directives

The CSP becomes `… script-src 'self'; connect-src 'self'` and nothing else changes. No inline
script, no `eval`, no CDN. The script (`workbench.js`, hand-written, **no framework, no build
step**) is embedded with `go:embed` and served by worldd at a hashed path. It consumes
`/agui/`, applies JSON Patch to a small local state, and swaps server-rendered fragments. Every
fragment still comes from the same `html/template` code, so the grades, checked links and closed
grammar stay server-side and keep their tests.

**The page must work with the script disabled.** Without JS it remains today's correct static
page. This is a test, not an aspiration.

*Alternatives:* htmx or Datastar are vendored JS too, and they add a dependency whose behaviour
the gates don't cover. They stay an option if the plan measures the hand-written script growing
past about 300 lines. React or Svelte bring a build chain into a Go repo that has none. That is
the user's fallback, not the default.

### D7 — Visual polish follows the mockups. The graph view is SVG laid out server-side

Row 149 ports the mockups' token system (`:root` variables, light and dark) into the workbench
templates. It adds the **decision-packet stream** and a **live timeline** with a "quiet is
health" footer (P5). It adds a **world graph view**: an SVG of the recent transition → object →
evidence neighbourhood, laid out deterministically in Go (a layered layout, not a physics
simulation), so it is snapshot-testable and identical on every capture. The live script only
animates insertion.

*Alternative:* a client-side force graph (d3-force or three.js) looks more dynamic but is
non-deterministic. That breaks reproducible marketing captures and pixel tests, and it adds
roughly 100 KB or more of vendored JS. Revisit after 1.0 if wanted.

### D8 — Marketing captures are a reproducible build artifact, not hand-made screenshots

Row 150 adds `ailang-worldd demo seed --db <path>`. It builds a deterministic demo world by
committing real transitions through the normal path: a goal, proposals, a decision packet with
evidence of each grade, an effect outcome and a provenance chain. Real data, not a fixture.
It also adds `scripts/capture_marketing.sh`, which serves the demo store and drives **system
Chrome headless** (`--headless --screenshot --window-size`) over a fixed shot list (workbench
light and dark, the decision packet, the graph, a provenance walk, a phone width). Output goes
to `website/static/img/v1/`.

The checks an agent can run itself: capture twice and assert byte-identical PNGs. Assert that every
shot is non-blank and that its size is in range.

*Alternatives:* `chromedp` is a Go dependency, which S3 rules out for a release script. Hand
screenshots go stale the first time the UI changes.

### D9 — Agent-generated views are typed world objects, rendered through World's own catalog

An agent proposes an A2UI surface as a world object (semantic id
`world/view/a2ui/<name>`, payload = A2UI messages pinned to the D1 spec version) through the
**normal propose → verify → commit pipeline**. It is a transition with authority and a trace,
like everything else. Verification checks the payload against **World's catalog**: a small,
closed set of components (text, list, table, card, object-link, grade-badge, timeline,
action-button) that the workbench can render.

The safety properties come from World, not from A2UI:

- `object-link` must reference a stored object (checked at verify time, using the same
  `checkedEdge` contract as the workbench). A view cannot link to something that doesn't exist.
- `grade-badge` renders the referenced object's **actual** grade. An agent cannot supply a grade
  value: §5's grade laundering is made unrepresentable, not just policed.
- Text from the agent renders as **CLAIMED** unless it is bound to an object (P2 grounded prose:
  ungrounded spans are visibly ungrounded).
- `action-button` maps to a **transition proposal**, never a direct effect. Pressing it opens
  the normal decision path with its own authority check.

The workbench renders committed views server-side (`/workbench?view=<ref>`, a new closed-grammar
key) and live via `world.view.a2ui`. Because the payload is standard A2UI, any A2UI renderer
(Lit, React, Flutter) can also draw it. That is the "no proprietary UI schema" promise of
DESIGN §11.2.

*Alternatives:* letting agents ship HTML or JS runs arbitrary code on the operator's origin and
is ruled out. An off-the-shelf A2UI renderer in the workbench would need a framework and could not
enforce World's grade and link rules. Server-side rendering through World's catalog keeps
those rules in tested Go.

## 4. Design sketch per row

**148 `w-worldd-agui-event-stream`** (`host/agui/`, new package; `host/daemon` registers the route)
- `agui.Encode(entry, objs) []Event` is pure and deterministic: the mapping table in D3.
  `agui.Tail(ctx, store, from, budget)` is the polling loop (D2, D5).
- The handler adds `POST /agui/` additively, plus the route-table doc comment update.
- QUICKSTART gets a new section, "watch the world live" (`curl -N` example). This is S7.

**149 `w-workbench-live-and-polish`** (`host/workbench`, `host/daemon/workbench.go`)
- The CSP changes per D6, with `TestWorkbenchSecurityHeaders` updated to the new literal and not to
  the symbol (the item-14 tautological-oracle lesson).
- Adds `workbench.js` (embedded), the fragment endpoints in the closed grammar, the token
  stylesheet, the packet stream, the live timeline, the quiet-is-health footer and the SVG graph view.

**150 `w-marketing-capture`** (`cmd/ailang-worldd` `demo seed`; `scripts/capture_marketing.sh`)
- Demo seed through real commits (D8), the capture script, a shot list file, and the
  `website/static/img/v1/` output. The README and website use the images.

**151 `w-agent-generated-projections`** (`world/view` AILANG types + verifier; `host/workbench` renderer)
- The A2UI view object type and the catalog live in `.ail`, as a pure validator with Z3 contracts
  where encodable (S1). Host rendering lives in Go (S2).
- The demo seed gains one agent-built view (for example "the decisions waiting on you this
  week") so row 150's shot list can show the feature.
- The MCP projection gains a `propose_view` tool, so a resident agent can actually do it.

## 5. Failure modes

| Failure | Behaviour |
|---|---|
| Client disconnects mid-stream | The handler returns on `ctx.Done()`. There are no goroutine leaks (`-race` plus a leak test) |
| Store read error or deadline | `RUN_ERROR` with the daemon's error class, then the stream closes. Never a silent stall |
| Client is too slow | Dropped at the write deadline. No unbounded buffering |
| Reconnect with a stale `Last-Event-ID` beyond head | `RUN_ERROR` "unknown index". Never a silent empty stream |
| Script disabled or blocked | The static page is unchanged and fully usable (tested) |
| View references a missing object | Refused at verify time. Never rendered as a dead link |
| View claims a grade | Unrepresentable in the catalog, so the payload is refused |
| A2UI spec bump | One-file mapping change. The pinned version constant fails a test until it is updated |

## 6. Milestones and acceptance (per row; the planner refines)

- **148:** (a) a fixture log of N entries produces byte-identical streams on two runs; (b) resume
  from `k` yields exactly entries `k+1…N`; (c) a commit by a *second process* on the same db
  appears within two ticks; (d) the connection cap and budget refuse or close as specified;
  (e) a stock AG-UI client parses the stream (a conformance fixture from the AG-UI repo,
  pinned).
- **149:** (a) the CSP equals the literal string; (b) with JS off, every existing workbench test passes
  unchanged; (c) with JS on, a commit appears in the timeline without a reload (headless Chrome
  check); (d) the graph SVG is golden-file stable; (e) grade rendering is never weakened (every
  existing grade test passes).
- **150:** (a) `demo seed` commits through the normal path, and `ailang-worldd why` walks its
  provenance chain; (b) capture runs twice and produces byte-identical PNGs; (c) the shot list
  is complete, with no blank images.
- **151:** (a) a view with a dangling `object-link` is refused; (b) a payload carrying a grade
  value is refused; (c) unbound text renders CLAIMED; (d) `action-button` produces a
  proposal, never an effect; (e) the same view renders in the workbench and parses in a stock
  A2UI renderer (pinned fixture).

## 7. Test plan and mutation rows (seed list; S6 requires the planner to finish it)

| Mutant | Must be killed by |
|---|---|
| Skip one entry when paging | 148(a) and 148(b) |
| Emit events for uncommitted state | 148(a): the fixture log has none |
| Drop the connection cap | 148(d) |
| CSP adds `'unsafe-inline'` to `script-src` | 149(a) literal check |
| Template renders a grade from view payload | 151(b) and the grade tests |
| `action-button` calls the broker directly | 151(d) |
| Seed writes the store directly, bypassing commit | 150(a) via `why` |

## 8. Non-goals

- No Godot, native app or Tauri shell in 1.0. Tauri stays the sanctioned post-1.0 path (HUMAN-SURFACE §6).
- No goal composer (surface 1) in these rows.
- No text-message or tool-call AG-UI events. No chat UI. The workbench is the workspace (§2.5).
- No closing of residual R1 here. The stream inherits whatever the R1 row does.
- No explainer film. Films release after 1.0, per the standing rule.

## 9. Conflict surface

- **Rule (d) and the critical path.** These rows move no UNMET clause. Placement and the loop
  exception are recorded in `D-WORLD-64`. They never displace a routable critical-path row.
- **CSP widening (D6)** loosens a pinned security header. It is limited to `'self'`, carries a
  literal-string test, and keeps the no-JS page fully functional.
- **A2UI is pre-1.0.** If v1.0 stable changes the message shapes before row 151 lands, row 151
  targets whichever version is stable at pick time, and D1's constant records it.
- **Read posture.** The stream widens no data exposure beyond the existing GET routes (D4). A
  reviewer should check that claim against the event mapping.

## 10. Verification log (base `2d53d72`, attended, 2026-10-06)

| # | Claim | Command / source |
|---|---|---|
| V1 | Routes and `isProtected` | `sed -n 790,831p host/daemon/daemon.go` |
| V2 | No worldd SSE of its own | `grep -rln text/event-stream host cmd` → MCP tests + `cmd/ailang-worldd/mcpclient.go` only |
| V3 | Workbench CSP literal | `host/daemon/workbench.go:20` |
| V4 | Log and head readers, header fields | `host/store/store.go:123-140, 687, 988` |
| V5 | Subcommand list, no demo seed | `grep -n 'case "' cmd/ailang-worldd/main.go` (lines 238–302) |
| V6 | Go deps | `go.mod` (two direct requires) |
| V7 | Mockup tokens | `design_docs/mockups/approval-inbox.html` `:root` block |
| V8 | AG-UI events/transport | docs.ag-ui.com/concepts/architecture, fetched 2026-10-06 |
| V9 | A2UI messages/catalog/status | a2ui.org/specification/v1.0-a2ui, fetched 2026-10-06 |
| V10 | `DecisionPacket` type exists | `grep -n DecisionPacket world/types.ail` → `:159` (semantic id comment `:155`) |

## 11. Quorum verification log

Not yet run. The loop's pick-time quorum records here.
