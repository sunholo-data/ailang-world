# w-agent-generated-projections (queue row 151: the agent's view) — r2

**Status:** r1, **direction approved by Mark, attended 2026-10-26** (*"yes lets update it"*),
recorded as **D-WORLD-92**. It supersedes row 151's parts of
[`w-world-live-surface.md`](w-world-live-surface.md): D9, §4 row 151 and §6 151(a)–(e). D1
(AG-UI carries events and A2UI is the payload, `D-WORLD-64`) still stands. Row 151 is still a
v1.0.0 release row. This is not quorum-reviewed: the loop's pick-time designer re-measures and
runs the quorum as usual. Nothing here waives a gate.

**r2 (attended 2026-10-26, D-WORLD-93):** adds §2a, named agents with a set output and a view
that comes with the agent, after Mark reviewed the
[front card mockup](https://claude.ai/artifact/UV1WBzA8Zjwn4hGdbr5go9) (private to Mark until
shared). It adds a thumbnail to every view (§5), acceptance items 11–14 (§7), and a separate
build row (§8). Nothing in r1 changes for an agent that isn't named.

## §0 What it is for (Mark, attended 2026-10-26)

> *"The purpose is really to give the AI agent a way to communicate what it's doing in a better
> way than just text, if it needs to."*

Three things follow from that sentence. Each is a rule in this design.

1. **It is communication.** A view tells the human what the agent is doing, has done or needs:
   its plan and where it is in it, a comparison, a number that matters, a result to look at, a
   question with options. It is not an app the agent builds, and it is not decoration.
2. **It is about the agent's work.** Every view belongs to one episode, and it is read in the
   front card of that episode's states, beside the inputs, what happened and the outputs that
   it describes.
3. **It is optional ("if it needs to").** Text stays the default. Most episodes will make no
   view, and that is quiet, not a gap: the front card shows one line, not an empty box.

## §1 What changes from D9

| | D9 (2026-10-06) | r1 |
|---|---|---|
| purpose | "agents can build the view a human needs" | an agent explains its work when text isn't enough (§0) |
| where it shows | the workbench, `/workbench?view=<ref>` | the World app's front card: "From the agent" at the top of Outputs, opening as a tab ([outputs doc](w-world-surface-outputs.md) §5); `/surface/static` for no-script |
| catalog | 8 components | 15, grouped by what the agent is trying to say (§3) |
| media | none | `media` shows one of the same state's outputs through the outputs-only blob route |
| numbers | none | `chart`, with data bound to objects or marked CLAIMED |
| updates | one committed view | a view is updated as the work moves: each update is a new version of the same object, so time navigation shows the view as it was (§4) |
| A2UI version | v1.0 RC, "whichever is stable at pick time" | unchanged rule; on 2026-10-26 the production line is still v0.9.1 and v1.0 is still a candidate (§8) |

## §2 When an agent should make a view

This guidance ships in World's own plugin (`plugins/ailang-world/`), where agents read it.

**Make one when the text would be worse:**
- **A plan with live state:** more than a few steps, each one done, working, waiting or blocked.
- **A comparison:** before and after, option A against option B.
- **Numbers that matter:** a timing, a size, a pass rate across rounds.
- **A result to look at:** an image, a render, a video the agent produced.
- **A question with options:** something the human must choose, with the choices laid out.

**Don't make one** to restate the summary, the files list or the checks (the front card already
shows them), to decorate, or to report progress the steps already show.

**One view per purpose.** An episode can name a few views, such as `plan` and `perf`. Updating a
view replaces what it shows; it never makes a second copy.

## §2a Named agents: a set output, and a view that comes with the agent (r2)

> *"we have very structured agents, for example sprint planners, executors and evaluators, that
> should have good structure to report in these cards."* … *"we don't run every task through
> sprint plan, executor etc. but it's nice if we do use a named agent it has a set output, much
> like our other AI protocol platform skills, and I guess new agents used could come with their
> own custom workbench/A2UI UIs."* (Mark, attended 2026-10-26, typos corrected.)

§2 covers an agent that decides for itself to make a view. A **named agent** is different: it
is a definition someone wrote and reuses, such as the mission loop's designer, planner, executor
and evaluator, a quorum, or a new agent someone installs. Its report should look the same every
time it runs, the way a skill on the platform always returns the same shape.

**The rules:**

1. **Most tasks use no named agent, and nothing changes for them.** Their card has §2's
   optional view or none.
2. **A named agent declares two things beside its prompt:**
   - an **output schema**: the record it must produce at the end of its run;
   - a **view template**: an A2UI composition from §3's catalog whose data model is bound to
     that record, plus a thumbnail (§5).
3. **The agent writes the record, World draws it.** World validates the record against the
   schema and renders it through the template. The agent doesn't compose UI at run time, so the
   same role reads the same way across runs and across models.
4. **A custom UI is a custom composition, never custom code.** A template uses §3's catalog,
   caps and rules, and nothing else: unbound content renders CLAIMED, links must resolve,
   `choice` and `action-button` only propose. A new agent can bring its own layout; it cannot
   bring a script, a style sheet or a URL.
5. **Scores are attributed, never grades.** An evaluator's "PASS 81/100" is that role's
   verdict, shown with the role's name on it. It never becomes an object's grade; `grade-badge`
   still shows only real grades. A line is TESTED only when a recorded re-run backs it (for
   example, the controller re-ran the executor's gate); the role's own word is CLAIMED.
6. **The header is World's, not the agent's.** Role, model, tokens, time, worktree and receipt
   come from the routing and session records, so an agent can't misreport who it was or what it
   cost.
7. **Absent stays absent.** A named agent that ends without a valid record shows "no report"
   and the validator's reason. Nothing is filled in.
8. **Templates are versioned with the agent.** The template ships in the agent's plugin. The
   episode records its hash, so an old state renders with the template it ran with.

**The first templates are the mission loop's roles**, built from what they already write (the
mockup's iteration 305 card shows all five, with real data):

| role | record it already writes | the view |
|---|---|---|
| designer | the design doc and its V-rows | revision strip; each V-row with its bound output file |
| planner | the sprint JSON and `planner-receipt.md` | milestones and their state; next milestone's readiness; flags left for a decision |
| executor | its report and the controller's re-runs | changes with diff thumbnails; gates; named mutants and where each went red |
| evaluator | `evaluator-r<N>.md` | score and rubric; mutation results; each finding and what happened to it |
| quorum | `quorum-r<N>.json` | one row per seat, one column per round, the objection in the seat's words; the disposition |

The loop's driver and its shared skill are fleet-owned (D-WORLD-DRIVER-1). World adds adapters
that read these files; any change to what the roles write is **proposed** to the fleet, never
applied from this repo.

**Where they show:** in Outputs as **Role reports**, one thumbnail each, opening as a tab named
for the role. What happened uses the roles as its actors. A file that a report owns (the judge
report, a quorum JSON) opens that report rather than its raw text.

## §3 The catalog

World's catalog is a closed subset of A2UI. A view that uses anything else is refused at
verify time and never rendered. The groups follow §0's purposes.

| group | component | what it shows | World's rule |
|---|---|---|---|
| say | `text` | a sentence or paragraph | grounded prose (HUMAN-SURFACE P2): spans bound to objects are links; unbound text renders **CLAIMED** |
| arrange | `row`, `column` | layout only | no `modal` and no `tabs`: the card already has tabs, and a view never takes over the screen |
| structure | `list`, `table`, `card` | items, rows, a titled group | cells follow the `text` rule |
| progress | `checklist` | a plan: steps with a state (done, working, waiting, blocked) | a state bound to a real object (a check, a run, a decision) shows that object's state; unbound states render CLAIMED |
| progress | `timeline` | events in order | each event is bound to a log entry, or is CLAIMED |
| numbers | `chart` | bar or line, at most 2 series and 200 points | data bound to objects (grader results, exec timings, sizes) or rendered CLAIMED; no free-form styling |
| results | `media` | an image, video or audio output | must name an output **of the same state**; served by the blob route (outputs doc §4); never a URL, never inline bytes |
| results | `output-link` | opens an output's tab (its diff or viewer) | must name an output of the same state |
| link | `object-link` | a link to any world object | must resolve at verify time (the `checkedEdge` contract) |
| link | `grade-badge` | an object's grade | renders the object's **real** grade; a grade value in the payload is unrepresentable |
| ask | `choice` | a question with options | becomes a decision packet; the answer goes through the decision path with its own authority check |
| ask | `action-button` | one proposed action | yields a transition proposal, never an effect |

**Not in the catalog, by design:**
- HTML, markdown with raw HTML, scripts, styles, external URLs, `data:` URIs and images by URL.
- A `state-graph` (w-world-os-surface R6, default no): the timeline is World's own.

**Limits** (pick-time designer may tune): 64 KB of payload and 200 components per view, 8 named
views per episode, and one update per view per 5 s (later ones coalesce).

## §4 Lifecycle: a view moves with the work

- **Propose:** the agent calls `propose_view(name, a2ui_messages)` on World's MCP projection. The
  view is a typed world object, `world/view/a2ui/<episode>/<name>`, committed through the normal
  propose → verify → commit path. It is a step with a trace, attributed to the agent's session,
  and inside the session's budget.
- **Update:** calling `propose_view` again with the same name commits a new version
  (A2UI's `updateComponents` and `updateDataModel` are applied, then the whole view is verified
  again). The front card at a state shows the view as it was at that state, so scrubbing back
  shows the plan as it stood then.
- **Live:** each commit emits AG-UI `CUSTOM world.view.a2ui` on the existing stream (row 148), and
  the app redraws the open tab without a reload.
- **Remove:** `deleteSurface` commits a tombstone; earlier versions stay in history.

## §5 Where it shows

- **No view** (the normal case): one quiet line at the top of Outputs, "From the agent · no view".
- **A view:** a row at the top of Outputs with its title, its kind (plan, comparison, numbers,
  result, question) and "updated hh:mm". Clicking it opens a tab; `f` maximises it
  (outputs doc §3). The URL is `open=view:<name>`.
- **A question** (`choice`): the row is marked "needs you", and the decision also appears
  wherever decisions do. When row 207 lands, the agent's state is "blocked" until it's answered.
- **A thumbnail** (r2): every view, and every named agent's report, has a small picture in
  Outputs drawn from its own data: a checklist's ticks, a chart's bars, a score ring, a
  quorum's verdict grid. The template says which component the thumbnail summarises; an ad-hoc
  view's thumbnail is its first component.
- **No script:** `/surface/static` renders the same view server-side from the same verified
  model, and `choice` and `action-button` are shown disabled with "answer in the app".

## §6 Build shape

- **`.ail`:** the view object type, the catalog and the rules, as a pure validator with Z3
  contracts where encodable (S1). The validator owns "unknown component", "grade in payload",
  "media names a non-output", the caps, and the CLAIMED marking.
- **Go (S2):** resolves the bindings (object links, grades, check and run states, outputs), then
  produces a **render model**: verified JSON with every grade and state already resolved.
  `GET /surface/view.json?brane=&at=&name=` serves it.
- **The app** renders the render model with its own small renderer in `app.js` (no framework,
  no stock A2UI renderer), so World's rules are enforced before anything reaches the page. A
  stock renderer can still draw the raw A2UI; that is the interoperability test, not the product
  path.
- **MCP:** `propose_view` (D9 named it). Refusals return the validator's reason, so the agent
  can fix and retry.
- **The plugin:** §2's guidance and two worked examples.

## §7 Acceptance (replaces 151(a)–(e))

1. A view with a dangling `object-link` is refused.
2. A payload carrying a grade value is refused.
3. Unbound text, an unbound checklist state and unbound chart data each render CLAIMED.
4. `action-button` produces a proposal and `choice` produces a decision packet; neither produces
   an effect.
5. `media` naming a path that isn't an output of the same state is refused; one that is plays
   through the blob route under the strict CSP.
6. A view updated three times has three versions, and the front card at each state shows the
   version current at that state.
7. A component outside the catalog, and each cap exceeded, is refused with the validator's reason.
8. The same view renders in the app's tab and in `/surface/static`, and parses in a stock A2UI
   renderer at the pinned spec version (a pinned fixture).
9. An episode with no view shows the one quiet line and no empty box.
10. **An attended episode:** a real agent, given §2's guidance, makes a plan view on a task with
    several steps and keeps it current, and makes no view on a one-line fix.

11. **(r2)** A named agent's record that fails its schema is refused with the validator's
    reason, and the card shows "no report" for that role.
12. **(r2)** A template that uses a component outside the catalog, a script or a URL is
    refused when the agent is installed, not when it runs.
13. **(r2)** An old state renders with the template version recorded for it, after the
    template has changed.
14. **(r2)** A role's score shows attributed to the role, and no object's grade changes because
    of it.

## §8 Order and dependencies

- **Needs** row 188's N2d (the blob route, for `media` and `output-link`) and N4a (the tabbed
  front card). Both come before N4. Row 151 stays a v1.0.0 release row (D-WORLD-64), so the
  1.0 path now runs through those two steps.
- **A2UI version:** the rule is unchanged: target whatever is stable at pick time, recorded in
  D1's version constant. On 2026-10-26 v0.9.1 is the production release and v1.0 is still a
  candidate ([google/A2UI](https://github.com/google/A2UI)). The catalog uses only components
  both have, or maps to them, so a version bump is one mapping file.
- **Named agents (§2a) are a separate build row,** filed when Mark places it. The
  recommendation is right after row 151 and off the 1.0 path. Row 151 builds the catalog, the
  validator and the renderer that §2a reuses, so §2a adds the schema and template layer, the
  loop's five adapters and the install-time check.
- **Row 150's demo seed** may gain one agent view (D9's idea), if row 150's split rows want it.

## §9 Open questions for the pick-time design

1. **The caps** in §3: start there, or measure on the first attended episodes?
2. **Who may see a view:** everyone who can see the episode (the default), or narrower?
3. **`choice` with free text:** allow an "other" answer, or options only?
4. **Old versions:** keep every version (the default; they are small), or only versions at
   states someone opened?
5. **(r2) Where a template lives:** beside the agent definition in its plugin, or in a
   World registry of templates that agents name?
6. **(r2) Who may install a named agent with a template:** anyone who can mint, or only an
   approval like a grant's?

## Sources

- [google/A2UI](https://github.com/google/A2UI) (fetched 2026-10-26: v0.9.1 production, v1.0
  release candidate)
- [A2UI v1.0 basic catalog implementation guide (candidate)](https://a2ui.org/specification/v1.0-basic-catalog-implementation-guide/)
- D9 and D1 in [`w-world-live-surface.md`](w-world-live-surface.md); P2 and P3 in
  [HUMAN-SURFACE.md](../HUMAN-SURFACE.md)
