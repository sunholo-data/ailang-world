# w-world-os-surface (proposed re-scope of queue row 188, clause 5) — r2

**Status:** r2 draft, attended with Mark on 2026-10-20. r1 cut r0's five views to one, and r2 cuts its three
distances to two (§3), at Mark's direction: *"each section must be useful and we dont duplicate … elegance, form and function.
where there is nothing left to take away"*. **Nothing here is ratified.** It
proposes re-scoping row 188 (`w-desktop-shell`) and filing one new row; §8 lists the rulings
it needs. Until they are made, row 188 stands as written. Mockups: §9.

Builds on [HUMAN-SURFACE.md](../HUMAN-SURFACE.md) (ratified v0.3), especially P4 (time is
navigation), P5 (quiet must read as health), P6 (five zoom levels), P7 (speculation as a
gesture) and §2.5 (one typed layer under every renderer, the CLI included).

## §0 What Mark asked (attended, 2026-10-20)

Shown the row 149 workbench over Tailscale: *"this is technical but useless to me as a UI
really."* Then, in order (typos corrected):

1. *"think operating system — an OS — we are changing the way people interact with computers."*
2. *"maybe we still drag and drop but its agents working on tasks? we interact with
   combinations of nodes/elements we can play, replay, rewind and stitch together"*
3. *"editors for prompts where we can drag in various file types … repos, folders, projects
   … authority configuration … decision approvals, audit trails … analytics … destination,
   outputs, artifacts that may be inputs to further work, downstream — node graphs?"*
4. *"and we can rewind the node graph to a certain point as well"*
5. *"we perhaps also want a CLI/text way to also work with the agents/graph — it looks like it
   can get very complex looking quickly … "branes" that advance through agents like windows —
   the time machine effect? flip between windows easily, but also have an overview in a system
   manager somewhere"*
6. *"we can represent branches in branes where they have two or more ancestors after a branch,
   upstream and downstream showing the state graph power. I like the 4-d time element, its the
   block universe idea peeking through"*
7. *"the branches aren't just display — they are places we can inspect the state at the time,
   and change it. reviewing all inputs/outputs/decisions etc"*

## §1 Problem (measured)

| id | measurement |
|---|---|
| M1 | The live workbench (rows 105/149) renders World's internals: log entries, objects, hashes. Mark operates **work**: what is running, what needs him, what it cost, whether it passed, the PR. Row 188 as written wraps that same workbench in a Tauri window. |
| M2 | A unit of work bigger than an episode already exists, but nothing names it. Row 189 is one: Mark's rulings D-189-1…5 (#316), then five episodes w189a…w189e (#320, #321, #324, #327, #329), then AC7. Agent minutes per episode fell 37.3 → 26.8 → 18.5 → 6.8 → 3.5. The only place that sequence is written down is `episodes.md` on the trial rig and the charter row's prose. |
| M3 | **Each step ran on the world the step before it made.** The rig's daemon binaries are kept by commit: w189b (19 Oct 19:34–20:37) ran on `ailang-worldd.9527bc5`, the #320 merge (w189a's credentials split). w189c (21:34–22:36) ran on `fd7b325`, the #321 merge (w189b's `episode start`). w189d and w189e (20 Oct) ran on `0ee1945`, which contains #324 (w189c) and #323 (row 198). A state does not only precede the next one; it becomes **the tool** the next one runs on. |
| M4 | **Real two-ancestor states exist, but only in git.** Two strands forked from `0ee1945` (#323 merged, 20 Oct 08:25): the mission loop's iteration 288 (row 156 M5b, 3 commits, #328) and episode w189d (row 189 M5, 2 commits, #327). They merged at 10:22 and 10:31, so `9003555` has two ancestors from two different strands of work. Earlier, row 198's strand forked from `96d4901` (06:06), w189c's #324 landed underneath it, and #323 merged at 08:25 with a conflict in `world-mission.md`, resolved by keeping both sides (row 199's note and row 200). |
| M5 | **World's own store is linear.** `host/store/store.go`: `worlds` holds immutable, content-addressed revisions (state root and log head), but there is one `selected_world_head` row. Each `LogHeader` has one `PrevEntryHash`, and the six-field header is frozen (`world/logepoch.ail`). A stale observed head is a `ConflictError`, not a branch. HUMAN-SURFACE P7 already records the gap: *"one selected head and no fork, branch, or compare-arms API"*. |
| M6 | World does record what each step **read**, which git cannot. Invocation records carry `episodeId`, `skillId`, input, output, plan and effects (row 182 M2), and `episode report w189e` counts 22 `ailang-read` calls. That gives a causal upstream (what this step actually looked at), narrower than ancestry (everything before it). |
| M7 | **What a step actually came from is small, and crosses branes.** w189e's transcript (`evidence/ailang-world/transcripts/w189e.jsonl` on the rig) has 22 `ailang-read` calls over 17 distinct files. By the commit that created each file, they come from: row 189's design doc (`28f90b5`), w189b (`episode.go`), w189c (`profile_init.go`), w189d (`plugin_help_test.go`, `host/plugindoc`, the plugin's README and skills), row 198's brane (`episode_report.go`), and row 138, two and a half weeks older (`session_new.go`, 3 Oct). That is six earlier steps on three branes, out of 931 commits in its ancestry (`git rev-list --count 0ffbb9c`). |
| M8 | Human acts per episode fell from 11 to 6 with row 189 (`docs/USE_WORLD_ON_YOUR_REPO.md`). The mint among them is still a typed `y` at a terminal. |

## §2 The model — a block universe of world states

**Every world state that has ever existed still exists.** States are immutable and
content-addressed (M5), so the full history is a block that can be looked at whole, not a past
that has gone. "Now" is a slice through the block, not a place you are stuck in. That is the
block-universe reading Mark named, and World's storage already makes it literal for one strand
(the append-only log, replay).

On top of that block:

- **State.** One world revision: what was true, who made it, and the checks that held.
- **Brane.** A named strand of work through the block: a sequence of states that advances as
  agents and humans take steps. Row 189 is a brane. So is the mission loop's release work, and
  so is a single one-episode fix. A brane is what Mark flips through (§3.2).
- **Step.** The move from one state to the next on a brane: an episode, a ruling, a merge. It
  has an actor, a cost, checks and a PR.
- **Branch.** A state with two or more successors on different branes. "Branch from here"
  starts a new brane at any past state. It never removes anything: the original brane keeps
  its later states (append-only, §6).
- **Merge.** A state with **two or more ancestors**. `9003555` (M4) is the worked example: one
  ancestor carries row 156 M5b from the loop's brane and the other carries the plugin from
  row 189's brane.

**Every state is a place, not a picture.** Selecting a state *enters* it: the world as it was
at that moment, with everything that went into it and came out of it (§3.2). Changing anything
there (an input, a grant, a decision) makes a new branch from that state. The state itself is
never edited, so every view of the past stays true.

Three kinds of edge connect states. Wherever one is drawn it is labelled with its kind, because
each answers a different question:

| edge | meaning | source today | example |
|---|---|---|---|
| **built on** | B's history includes A (ancestry) | git parents; World's `PrevEntryHash` on one strand | `9003555` is built on `7fb99da` and `f36428a` |
| **read** | a step in B looked at an object A wrote | invocation records (M6) | w189e read `episode_report.go`, which row 198's brane wrote (M7) |
| **ran on** | a step in B ran on a tool built from A | the rig keeps binaries by commit (M3); nothing records it as a World fact yet | w189c ran on `fd7b325` |

**What it came from, and what it led to.** For any state, World can show:

- **what it came from**: everything that could have shaped it (built on) and, more narrowly,
  what actually did (read, ran on);
- **what it led to**: everything built on it since, across every brane, including steps that
  only ran on it.

Ancestry alone says everything before w189e is something it came from. The read and ran-on edges say
which few states mattered: for w189e, six steps on three branes (M7). Only World can show that
narrower set, because only World records the reads, and that is the "state graph power" the
view exists to show.

## §3 The surface — two distances, one action

There is one view. The whole world is the block (§2), and you look at it from one of two
distances. Each answers a question the other does not:

| distance | the question it answers | what you see |
|---|---|---|
| **The world timeline** (far) | *Where is everything, and what needs me now?* | Every brane as a lane through time, with branches and merges where lanes split and join. The right edge is **now**. |
| **One brane** (near) | *What happened here, how did it get here, and what if it had gone differently?* | The brane's states as a Time Machine stack. **The front card is the state, open**: inputs (live controls), outputs, decisions and checks. |

Clicking a lane or a state opens its brane, at that state. The address bar takes you back out
(§3.3). The zoom is spatial: you move into the point you clicked.

**One colour for you.** Amber marks the human, and nothing else is amber:

- on the timeline, a number on each state counts your decisions there (w189d: 3, for the
  mint, the approval it asked you for, and the merge);
- a dashed, pulsing state is waiting for you, and so is everything under *needs you*;
- on a card, every decision carries its actor (You, Agent, World, Review, Loop), and yours are
  filled amber;
- inputs you have set or changed turn amber, and so does the one button that proposes them.

Text is kept to labels and numbers. Checks are tiles, outputs are file names, decisions are a
few words each, and the proposal bar appears only when there is something to propose.

### §3.1 Far: the world timeline, with now as its edge

- **Lanes** are branes, with the shared world (dev) as one lane among them. A step leaves the
  shared world at the state it started from and returns to it at its merge, so a merge state's
  two ancestors are visible as two lines meeting.
- **Pointing at a state shows its connections:** blue for what it came from, amber for what
  it led to, with arcs for what it read (purple) and the build it ran on (teal). These are
  the same links the open card lists (§3.2), drawn in place. There is no toggle and no mode.
- **The now edge** holds everything a separate system overview would have shown, and nothing
  else:
  - a positive, timestamped health line ("World healthy · checked 12:55 · daemon `0ee1945`,
    clean"), so quiet reads as health (P5);
  - what is running and how many sessions are live;
  - **needs you**: open steps and proposed branches, each one click from its card.
- **Open and proposed steps sit on the now edge, dashed.** That covers AC7 today, and any
  branch you have proposed but not yet approved.

### §3.2 Near: one brane, the front card open

The brane's states are stacked in depth like Time Machine windows. Older states recede behind,
and the newest is in front, unless an open step is waiting for you, in which case that step
leads. **The front card is the state itself, fully open**, so there is no "go inside" step:

| | what it holds | w189e, measured |
|---|---|---|
| **Inputs** | *came from*: the states it was built on, read, and ran on; plus the task and the agent | built on `9003555`; read design, w189b, w189c, w189d, w198, row 138; ran on `0ee1945`; `tasks/w189e.md`; 22 reads over 17 files |
| **Outputs** | *led to*: the states that were built on it, read it, or ran on it; plus files, objects, the PR | merged into `0ffbb9c`; the page and its test; links from three files; PR #329, 7 files |
| **Decisions** | the mint, its grants and TTL, questions the agent asked and the answers, refusals, review rounds, the merge | typed mint, 12 h; read and write in the worktree, exec ×600 (6 used); 0 refusals; "merge 329"; session revoked |
| **Checks** | graders, mutants caught, CI, judges, calls and minutes | 8/8 in 1 round; 3 help-check mutants caught; 42 calls, 3.5 agent minutes |

**Moving.** *Back in time* and *toward now* step along **this brane only**. That is why
navigation never has to choose between ancestors. Leaving the brane is always explicit: click
a *came from* or *led to* link, and the stack becomes that state's brane, with that state in
front. At a merge (a shared state with two ancestors), *back* follows the shared world. The
other ancestor is a link, one click away, and so is its brane. A rail lists every state on the
brane, so you can jump to any of them.

**The one action: change from here.** The front card's **inputs are live controls**: the agent,
the task, the grants. Change from here is not a button or a mode; it is what happens when you
change one of them.

- **Unchanged,** they are the record of what ran. The card says so.
- **Changed,** the card lists what changed ("agent Claude Opus 5.5 → Codex, grants −Delete
  ×200"), names the branch it would make, and offers **Propose** and **Discard**. For a step,
  the branch re-runs that step from the world it started in. For a shared state, it starts new
  work from that world.
- **On an open step** (AC7 today), the same controls are how you start it: pick the agent,
  write the task, and propose the start.
- A proposal is a decision packet (P1). It appears on the now edge as *needs you* (and, for a
  branch, as a new dashed lane), and nothing runs until the typed `y`.
- The original brane keeps every later state. Nothing is edited, so every view of the past
  stays true.
- Grants that can never be given, such as network, are shown and locked rather than hidden.
  The default grants are the ones every row 189 episode was minted with (read, write, delete
  ×200, exec ×600, in the worktree).

### §3.3 The address bar — the same space as text

The address bar is both the breadcrumb and the command line. Every place you can click to has
a text form:

| you type | you get |
|---|---|
| `world` | the world timeline |
| `world onboarding` | that brane, newest (or open) state in front |
| `world onboarding@6` | that brane, its 6th state in front |
| `world shared@0ee1945` | the shared world, that state in front, by hash |
| `back` | out to the world timeline |

The same addresses work in a terminal. Today's commands (`episode report`, `why`, `provenance`)
already answer front-card questions, so their output is what the card renders.

### §3.3a Keys — one map for both distances

Time always runs on one axis and lanes on the other, so the same keys mean the same thing at
either distance:

| key | world timeline | one brane |
|---|---|---|
| ↑ ↓ | lane above · below | back in time · toward now |
| ← → | earlier · later on the lane | the brane above · below, at the nearest moment |
| End · Home | most recent · earliest | most recent · earliest |
| Enter | open the brane | — |
| Esc | — | back to the timeline |
| N | jump to what needs you | jump to what needs you |
| / | the address bar | the address bar |
| ? | show these keys | show these keys |

Keys never fire while you are typing in the address bar or the task. Esc there returns to the
board.

### §3.4 What r1 and r2 cut, and why

| earlier view | where its job went |
|---|---|
| System manager (r0) | the now edge (§3.1). Its brane table was the timeline's lanes a second time. |
| Block as a separate screen (r0) | the world timeline, as the home |
| Graph, the node editor (r0) | the front card. Its sources, prompts, grants and outputs *are* the state's inputs, outputs and decisions, and "change from here" edits them. |
| Console as a view (r0) | the address bar (§3.3), and the same addresses in a terminal |
| "Inside a state" as a third distance (r1) | the front card of the brane stack (r2). A separate inside view needed its own back and forward, which is ambiguous at a merge. |
| The "cones" toggle (r1) | always on: pointing at a state shows what it came from and what it led to |
| The "Change from here" button and its panel (r1) | the front card's inputs, which are the controls (r2) |

## §4 Worked example — row 189 as a brane (real data, 2026-10-19/20)

| # | state | who | PR / merge | ran on | agent min | checks |
|---|---|---|---|---|---|---|
| 1 | design: D-189-1…5 = C, A, C, A, A | Mark | #316 → `5386215` | — | — | — |
| 2 | M0+M1 credentials store split (w189a) | Claude | #320 → `9527bc5` | (an earlier build) | 37.3 | 5/5 |
| 3 | M2+M3 `episode start` (w189b) | Claude | #321 → `fd7b325` | `9527bc5` | 26.8 | 8/8 |
| 4 | M4 `profile init` (w189c) | Claude | #324 → `fe2cc96` | `fd7b325` | 18.5 | 8/8 |
| 5 | M5 plugin (w189d) | Claude | #327 → `9003555` (**merge: two ancestors**, M4) | `0ee1945` | 6.8 | 8/8 |
| 6 | M6 docs page (w189e) | Claude | #329 → `0ffbb9c` | `0ee1945` | 3.5 | 8/8 |
| 7 | AC7: a real platform-source episode by the page | Mark | open | — | — | — |

State 5's other ancestor is on the release brane (iteration 288, #328). State 4 reaches the
later states by all three edges:

- **built on:** state 5 includes it;
- **ran on:** states 5 and 6 ran on `0ee1945`, which includes it;
- **read:** state 6 read `profile_init.go`.

## §5 What exists, and what has to be built

| need | exists today | to build |
|---|---|---|
| States, immutable and replayable | `worlds` and the append-only log (M5); replay proven | — |
| Steps with cost, checks, PR | `episode report`, outcome records (row 182), ledger | a projection that joins them per step |
| **Front card: inputs** | reads and effects are World records (M6). The task prompt is a file on the trial rig, not a World object. | commit the task and prompt as an object at `episode start`, so "change from here" edits a World object |
| **Front card: decisions** | grants and TTL at mint (credentials store); refusals (`episode report`); review rounds in the ledger. Merges happen on GitHub, outside World. | record the merge decision and the review rounds as World facts, linked to the state they produced |
| **The brane object** | nothing names a brane (M2) | `world/brane/v1`: a name, its states, its steps (episode ids, ruling ids, PRs); appended as states advance, never edited |
| Read edges | invocation records (M6) | an index from object to the invocations that read it |
| Ran-on edges | binaries kept by commit on the rig, by convention only (M3) | record the daemon's build (`vcs.revision`, `vcs.modified`) on every invocation record, or on the episode's mint |
| **Fork: "change from here"** | none: one `selected_world_head` (M5) | **M5's fork API**: more than one head, named by brane; a fork is a new head pointing at an existing world, which is O(1) through structure sharing (DESIGN.md §13.5) |
| **Merge: two-ancestor states** | git only (M4) | a merge record object naming every ancestor world, referenced from the merging entry's transition. The frozen six-field `LogHeader` keeps one `PrevEntryHash`, so ancestors cannot go in the header. Conflicts need a rule (§10.1). |
| The view (§3) | the workbench and AG-UI stream served by `worldd` | the timeline and the brane stack as read-only projections first; then change |
| Addresses (§3.3) | `episode report`, `why`, `provenance` | address parsing in the CLI and the view, one grammar for both |
| Approval as a gesture | the typed `y` at a TTY (the human-typed rule) | only on a ruling (§8 R3) that keeps the human-typed guarantee: the gesture must be unforgeable by an agent in the same way |

### §5.1 Fit with the protocol-led approach (D-WORLD-64)

D-WORLD-64 (attended, 2026-10-06) ruled three things:
- the reference renderer stays the `worldd`-served workbench;
- AG-UI is the event transport and A2UI is the payload for agent-generated views ("both,
  layered");
- Godot is out.

This design keeps all three.

- **Transport: AG-UI, unchanged.** Row 148 (landed, `38db793`) streams committed log entries
  over SSE from `POST /agui/`, with resume by `state.lastIndex`. The encoder (`host/agui`)
  already emits `RUN_STARTED`, `RUN_FINISHED`, `RUN_ERROR`, `STATE_SNAPSHOT`, `STATE_DELTA` and
  `CUSTOM`. A brane advancing *is* new committed entries, so the view stays live over the
  existing stream. How brane projections ride `STATE_*` events is S1 design work.
- **Payload: A2UI fits the brane stack.** Row 151 (not landed) gives agents a
  closed catalog: text, list, table, card, object-link, grade-badge, timeline and
  action-button. A brane is a timeline of cards, and the front card is four lists whose
  *came from* and *led to* entries are object-links. So an agent could generate it, with D9's rules intact:
  links must resolve, grades come from the referenced object, and ungrounded text renders
  CLAIMED.
- **Actions: already right.** "Change from here" and "propose" are action-buttons, and D9 maps
  an action-button to a **transition proposal**, never an effect. That is §6's authority rule in
  protocol form.
- **The one gap: the world timeline is not in the catalog.** Lanes, merges and connections are
  neither a table nor a timeline. There are two ways forward:
  - (a) draw them as a **bespoke view in the reference renderer**, as the workbench already
    does for entry edges (`host/workbench/graph.go`, `LayoutGraph`);
  - (b) **extend the catalog** with a read-only `state-graph` component whose every node must be
    an object-link.

  Recommendation: (a). Do (b) only if agents turn out to need to draw it (R6).
- **Renderer and shell.** R5's default (browser first) is D-WORLD-64's reference renderer.
  Row 188's desktop shell wraps the same page and adds native notifications.

## §6 Invariants (unchanged by this design)

- **Nothing is ever deleted.** "Change from here" is a branch, never a reset. Every past state
  stays in the block.
- **Authority does not move.** Minting stays human-typed. No view, address or gesture lets an
  agent approve, mint, publish or push. The view can only *propose* (P1).
- **No new network exposure** without Mark's ruling. GET routes are unauthenticated today
  (row 199 R1), so serving the view beyond loopback waits on that row.
- **One typed layer.** The view holds no state of its own. Everything shown reads from the
  store, so the view and the CLI agree by construction.

## §7 Proposed sequence

Each stage is useful on its own and adds no authority until S5.

1. **S1 — the timeline and the brane stack, read-only, over today's data.** Branes from charter rows, states from
   git merges, episode reports and outcome records. No store change. This answers M1
   directly.
2. **S2 — the open front card and addresses.** Inputs, outputs, decisions and checks; *came
   from* and *led to* over the read edges that exist (M6); one address grammar in the view and
   the CLI.
3. **S3 — the brane object, ran-on recording, and the task as a World object.** Brane
   membership and inputs become World facts rather than projections over git and the rig.
4. **S4 — the fork API (M5 pulled forward) and merge records.** Real two-ancestor World
   states, and the store side of "change from here".
5. **S5 — change from here.** Proposing a branch as a decision packet; approval stays the typed
   `y` unless R3 is ruled yes.

## §8 Rulings for Mark

Proposed. When answered, each one is filed OPEN in the charter and recorded through
`mission_answer.sh`, as usual.

- **R1 — re-scope row 188** from "a Tauri shell around the workbench" to this view, sequenced
  as §7. Default if unanswered: hold (row 188 stays as written).
- **R2 — file a new row for the fork API and merge records** (S4), pulling that part of M5
  forward as the dependency of "change from here". Default: hold.
- **R3 — approval as a gesture** as well as the typed `y`, provided the gesture is human-only
  in the same sense as the typed `y`. Default: no (the typed `y` stays the only path).
- **R4 — the word "brane"** in the product, addresses included, or a plainer word (strand,
  line). Default: "brane" in the design, revisited before S1 ships.
- **R5 — the first renderer**: the browser via `worldd` (as the workbench is today), or the
  desktop shell row 188 named. Default: the browser first; the shell wraps the same page later.
- **R6 — a `state-graph` component in World's A2UI catalog** (§5.1 b), so agents can draw the
  world timeline too. Default: no.

## §9 Mockup

One board, real data from 19–20 October 2026, on the private design canvas:
`https://claude.ai/artifact/4z7uXX5GgRj9ZzdAjxwD6L`

- It opens on the world timeline, showing what w189e came from and led to.
- Point at a state to see its connections, click a lane or a state to open its brane, and
  type addresses in the bar.
- "Change from here" proposes a branch, which then appears in the block and under *needs you*.
- Proposals made on the board are the only content that is not real.

r0's eleven boards (v1–v5) are superseded and are kept in the canvas's version history.

## §10 Open questions

1. **Merge semantics for World objects.** Git merges text. A World merge of two typed states
   needs a rule per object kind, plus somewhere to resolve conflicts. What is the smallest set
   of kinds that covers the merges actually seen (M4 was the charter file)?
2. **What starts a brane?** A charter row, a ruling, an episode with no parent brane? S1 can
   infer branes from charter rows; S3 has to decide.
3. **The mission loop as a resident.** The loop's iterations form a brane. Should its states
   show like attended episodes, or fold into one step per row?
4. **Long time.** The board covers two days. Months of history need the world timeline to
   compress quiet stretches. What does it compress, and how does it show that it has?
