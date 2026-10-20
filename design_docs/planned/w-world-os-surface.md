# w-world-os-surface (proposed re-scope of queue row 188, clause 5) — r0

**Status:** r0 draft, written attended with Mark on 2026-10-20. **Nothing here is ratified.** It
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
| M7 | **The read cone is real and crosses branes.** w189e's transcript (`evidence/ailang-world/transcripts/w189e.jsonl` on the rig) has 22 `ailang-read` calls over 17 distinct files. By the commit that created each file, they come from: row 189's design doc (`28f90b5`), w189b (`episode.go`), w189c (`profile_init.go`), w189d (`plugin_help_test.go`, `host/plugindoc`, the plugin's README and skills), row 198's brane (`episode_report.go`), and row 138, two and a half weeks older (`session_new.go`, 3 Oct). That is six earlier steps on three branes, out of 931 commits in its ancestry (`git rev-list --count 0ffbb9c`). |
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
at that moment, with everything that went into it and came out of it (§3.0). Changing anything
there (an input, a grant, a decision) makes a new branch from that state. The state itself is
never edited, so every view of the past stays true.

Three kinds of edge connect states. Each one must be labelled with its kind wherever it is
drawn, because they answer different questions:

| edge | meaning | source today | example |
|---|---|---|---|
| **contains** | B's history includes A (ancestry) | git parents; World's `PrevEntryHash` on one strand | `9003555` contains `7fb99da` and `f36428a` |
| **read** | a step in B looked at an object A wrote | invocation records (M6) | w189e read `episode_report.go`, which row 198's brane wrote (M7) |
| **runs on** | a step in B ran on a tool built from A | the rig keeps binaries by commit (M3); nothing records it as a World fact yet | w189c ran on `fd7b325` |

**Upstream and downstream cones.** Pick any state and World can light up:

- its **upstream cone**: everything that could have shaped it (contains), and, narrower, what
  actually did (read and runs on);
- its **downstream cone**: everything built on it since, across every brane, including steps
  that only ran on it.

Ancestry alone says everything before w189e is upstream of it. The read and runs-on edges say
which few states mattered: for w189e, six steps on three branes (M7). Only World can draw that narrower cone, because only World
records the reads, and that is the "state graph power" the views exist to show.

## §3 The surface — four views over one model

Every view reads the same records (HUMAN-SURFACE §2.5: one typed layer underneath). None of
them owns state, so none can disagree with another.

### §3.0 Inside a state — the place every view opens onto

Every view below shows states, and every state opens the same way. This is P6's **inspect**
and **walk** zoom, applied to a moment in time:

- **Inputs:** the task or prompt, the files and objects the step read (M7), the world it
  started from, and the build it ran on (M3).
- **Outputs:** files written and edited, objects committed, the PR, and artifacts that feed
  later steps.
- **Decisions:** who minted the session and with which grants and TTL, what the agent asked the
  human and the answer, refusals by plan and in handlers, review findings and rounds, and the
  merge.
- **Checks:** graders, mutants caught, CI, judge scores, calls and minutes.

Two actions sit on every state:

- **Open the world as of here:** browse its files, objects and log, read-only. This is what
  replay (P4) gives a human.
- **Change from here:** edit any input, grant or decision and run again. This always starts a
  new brane at that state (needs the fork API, S5), then offers to re-run the downstream steps
  on the new brane, so the two futures can be compared side by side (P7). The original brane
  is untouched.

Worked example, w189e (on the v5 board):

- **Inputs:** `tasks/w189e.md`; 22 reads over 17 files; daemon `0ee1945`.
- **Outputs:** the new page and its help-check test; links from three files; PR #329, 7 files.
- **Decisions:** Mark's typed mint (12 h); grants of read and write in the worktree and exec ×600,
  of which 6 were used; 0 refusals; "merge 329"; the session revoked.
- **Checks:** 8/8 in 1 round; 3 help-check mutants caught.

### §3.1 System manager — the overview

What is true right now, across all branes. Quiet must read as health (P5), so this view always
shows a positive, timestamped "World healthy, checked HH:MM", never an empty list.

- **Vitals:** World health and daemon build, agents running, how many things wait for the
  human.
- **Branes table:** each brane with its current state, step pips, residents, state and what it
  needs from you. Clicking a brane flips to it (§3.2).
- **Residents:** agents and loops resident on this World, with what they are doing.
- **Authority:** live sessions, expired sessions kept for audit, revocations, active repo
  profiles, network grants.
- **Work:** episodes, PRs merged, agent minutes and new friction notes over a window.

### §3.2 Branes — Time Machine

One brane as a stack of windows, one per state, receding into depth. The newest is at the
front, unless the brane has an open step waiting for you, in which case that step leads.

- Each window is the state at P6's **decide** zoom: who made it, what changed, checks, rounds,
  minutes, calls, and **"World could now…"**. The last one is the operator-facing meaning of
  the step ("mint while another episode runs").
- **Flip** back and forward with keys, scroll or the timeline rail. **Tabs** switch branes.
- **Branch from this state** (P7) starts a new brane. It needs the fork API (§5).
- At a merge, the window shows **both** ancestors as tabs on its edge, so you can flip up into
  the brane the other ancestor came from. That is how a brane view stays a single strand and
  still shows the graph: the graph is reached through the merges, not drawn all at once.

### §3.3 Block — the 4-D view

Zoomed out from the branes: the whole block as a spacetime diagram. Time runs along one axis,
branes run as worldlines along it, and branches and merges are where worldlines split and
join.

- **Now** is a plane you can drag through the block. Everything left of it is fixed history;
  to its right are open steps and speculative branches (P7), drawn dashed and labelled
  speculative.
- Selecting any state lights up its **upstream and downstream cones** (§2), with a toggle for
  edge kinds: ancestry only, or read and runs-on.
- Selecting two states gives the **world-diff** between them (P4).

This is the view that makes the block-universe idea visible. It is an overview to navigate
from; it is not where work is edited.

### §3.4 Graph — composing one step

When a step needs rewiring, open it as a node graph. **Sources** are repos, folders and files.
**Prompts** are editors you drop files into. **Authority** nodes carry grants, and a disallowed
wire refuses to connect. **Workers** are agents. **Gates** are checks and approvals.
**Outputs** are PRs and artifacts, which can feed other nodes downstream. **Lenses** are
analytics and the audit trail.

Running the graph proposes a step (P1: a decision packet), and the graph's own edit history is
scrubbable and branchable like any brane. The graph is deliberately one level down: it gets
complex fast (Mark, §0.5), so the brane is the default view and the graph is opened per step.

### §3.5 Console — the same thing as text

A World shell that addresses branes, states and graphs by name. No capability exists in the
views that the console lacks, and the reverse holds too (§2.5). Today's commands keep working
unchanged; the brane verbs are new and only read, except `rewind` and `run`.

| proposed verb | what it does | today's nearest equivalent |
|---|---|---|
| `world branes` | the System manager's brane table | none (charter rows and `episodes.md` by hand) |
| `world brane <b> log` | a brane's steps, newest first | `ledger.py` rows on the rig |
| `world brane <b> show <n>` | one step at decide zoom | `ailang-worldd episode report <ep>` (one episode) |
| `world brane <b> cone <n> [--down] [--edges contains,read,runs-on]` | upstream or downstream cone | `ailang-worldd why <entry>` (one entry, upstream only) |
| `world brane <b> diff <n> <m>` | world-diff between two states | `git diff` between merge commits |
| `world brane <b> rewind <n> --branch <name>` | new brane from state n; the original is untouched | none (needs the fork API) |
| `world brane <b> run --agent <a> --prompt @file` | propose the next step as a decision packet | `ailang-worldd episode start …`, then the typed `y` |

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

State 5's other ancestor is on the release brane (iteration 288, #328). State 4 reaches the later
states by all three edges: by **contains** (state 5 includes it), by **runs on** (states 5
and 6 ran on `0ee1945`, which includes it) and by **read** (state 6 read `profile_init.go`). That is the shape the Block
view (§3.3) exists to show.

## §5 What exists, and what has to be built

| need | exists today | to build |
|---|---|---|
| States, immutable and replayable | `worlds` and the append-only log (M5); replay proven | — |
| Steps with cost, checks, PR | `episode report`, outcome records (row 182), ledger | a projection that joins them per step |
| **Inside a state: inputs** (§3.0) | reads and effects are World records (M6). The task prompt is a file on the trial rig, not a World object. | commit the task and prompt as an object at `episode start`, so "change the input" edits a World object |
| **Inside a state: decisions** | grants and TTL at mint (credentials store); refusals (`episode report`); review rounds in the ledger. Merges happen on GitHub, outside World. | record the merge decision and the review rounds as World facts, linked to the state they produced |
| **The brane object** | nothing names a brane (M2) | `world/brane/v1`: a name, its states, its steps (episode ids, ruling ids, PRs); appended as states advance, never edited |
| Read edges | invocation records (M6) | an index from object to the invocations that read it |
| Runs-on edges | binaries kept by commit on the rig, by convention only (M3) | record the daemon's build (`vcs.revision`, `vcs.modified`) on every invocation record, or on the episode's mint |
| **Fork: "branch from here"** | none: one `selected_world_head` (M5) | **M5's fork API**: more than one head, named by brane; a fork is a new head pointing at an existing world, which is O(1) through structure sharing (DESIGN.md §13.5) |
| **Merge: two-ancestor states** | git only (M4) | a merge record object naming every ancestor world, referenced from the merging entry's transition. The frozen six-field `LogHeader` keeps one `PrevEntryHash`, so ancestors cannot go in the header. Conflicts need a conflict view of their own (§10.1). |
| Views (§3) | the workbench and AG-UI stream served by `worldd` | System, Branes and Block as read-only projections first; then the Graph |
| Console verbs (§3.5) | `episode report`, `why`, `provenance` | the read verbs first; `rewind` and `run` after the fork API |
| Approval as a gesture | the typed `y` at a TTY (the human-typed rule) | only on a ruling (§8 R3) that keeps the human-typed guarantee: the gesture must be unforgeable by an agent in the same way |

### §5.1 Fit with the protocol-led approach (D-WORLD-64)

D-WORLD-64 (attended, 2026-10-06) ruled three things. The reference renderer stays the
`worldd`-served workbench. AG-UI is the event transport and A2UI is the payload for
agent-generated views ("both, layered"). Godot is out. This design keeps all three, and
introduces one gap.

- **Transport: AG-UI, unchanged.** Row 148 (landed, `38db793`) streams committed log entries
  over SSE from `POST /agui/`, with resume by `state.lastIndex`. The encoder (`host/agui`)
  already emits `RUN_STARTED`, `RUN_FINISHED`, `RUN_ERROR`, `STATE_SNAPSHOT`, `STATE_DELTA` and
  `CUSTOM`. A brane advancing *is* new committed entries, so every view in §3 stays live over
  the existing stream with no second transport. How a brane projection rides `STATE_*` events
  is S1 design work.
- **Payload: A2UI fits the System manager and the Branes view.** Row 151 (not landed) gives
  agents a closed catalog: text, list, table, card, object-link, grade-badge, timeline and
  action-button. The System manager is tables and cards. A brane window is a card with
  grade-badges, and the brane rail is a timeline. So an agent could generate its own brane or
  system view, with D9's rules intact: links must resolve, grades come from the referenced
  object, and ungrounded text renders CLAIMED.
- **Actions: already right.** "Branch from here", `run` and approve-from-a-view are
  action-buttons, and D9 maps an action-button to a **transition proposal**, never an effect.
  That is §6's authority rule in protocol form. Approval as a gesture (R3) would still go
  through the decision path; it is never a button that approves on its own.
- **The gap: the Block view and the Graph are not in the catalog.** A spacetime diagram with
  cones and a node editor are neither a table nor a timeline. There are two ways forward:
  - (a) draw them as **bespoke views in the reference renderer**, as the workbench already does
    for entry edges (`host/workbench/graph.go`, `LayoutGraph`), so agents cannot generate
    them;
  - (b) **extend the catalog** with a read-only `state-graph` component whose every node must be
    an object-link, so the same grade and link rules apply.

  Recommendation: (a) for S4 and S6, and (b) only if agents turn out to need to make these
  views (R6).
- **Renderer and shell.** R5's default (browser first) is D-WORLD-64's reference renderer.
  Row 188's desktop shell wraps the same pages and adds native notifications. The Console is
  the CLI over the same REST routes, the "one typed layer" of HUMAN-SURFACE §2.5.

## §6 Invariants (unchanged by this design)

- **Nothing is ever deleted.** Rewind is a branch, never a reset. The original brane keeps its
  later states, and "rewound" history stays in the block.
- **Authority does not move.** Minting stays human-typed. No view, verb or gesture lets an agent
  approve, mint, publish or push. A view can only *propose* a step (P1).
- **No new network exposure** without Mark's ruling. GET routes are unauthenticated today
  (row 199 R1), so serving these views beyond loopback waits on that row.
- **One typed layer.** No view holds state of its own. Everything shown reads from the store,
  so the console and every view agree by construction.

## §7 Proposed sequence

Each stage is useful on its own and adds no authority until S6.

1. **S1 — read-only projections over today's data.** System and Branes, built from charter
   rows, episode reports, outcome records and git merges. No store change. This answers M1
   directly.
2. **S2 — read console verbs.** `branes`, `log`, `show`, `cone` (contains and read edges), and
   `diff`.
3. **S3 — the brane object and runs-on recording.** Brane membership becomes a World fact
   rather than a projection over git.
4. **S4 — the Block view**, over S3's data, with cones by edge kind.
5. **S5 — the fork API (M5 pulled forward)** and merge records: "branch from here", `rewind`,
   and real two-ancestor World states.
6. **S6 — the Graph and `run`**: composing and proposing a step. Approval as a gesture only if
   R3 is ruled yes.

## §8 Rulings for Mark

Proposed. When answered, each one is filed OPEN in the charter and recorded through
`mission_answer.sh`, as usual.

- **R1 — re-scope row 188** from "a Tauri shell around the workbench" to this surface,
  sequenced as §7. Default if unanswered: hold (row 188 stays as written).
- **R2 — file a new row for the fork API and merge records** (S5), pulling that part of M5
  forward as the dependency of "branch from here". Default: hold.
- **R3 — approval as a gesture** (hold to approve) as well as the typed `y`, provided the
  gesture is human-only in the same sense as the typed `y`. Default: no (the typed `y`
  stays the only path).
- **R4 — the word "brane"** in the product, CLI included, or a plainer word (strand, line).
  Default: "brane" in the design, revisited before S2 ships.
- **R5 — the first renderer**: the browser via `worldd` (as the workbench is today), or the
  desktop shell row 188 named. Default: the browser first; the shell wraps the same pages
  later.
- **R6 — a `state-graph` component in World's A2UI catalog** (§5.1 b), so agents can generate
  Block and Graph views as well as people. Default: no, and the two views stay bespoke in the
  reference renderer.

## §9 Mockups

Private design canvas (attended, 2026-10-20), real data unless a board says otherwise:
`https://claude.ai/artifact/4z7uXX5GgRj9ZzdAjxwD6L`

- v1 — an app about episodes (Home, Episode story, Start request);
- v2 — World as the operating system (the desk, intent in, open a thing);
- v3 — everything is a node (Graph, §3.4);
- v4 — branes: Brane (§3.2), Console (§3.5), System (§3.1);
- v5 — Block (§3.3 and §3.0): the 4-D view with cones and a draggable "now"; any state opens
  to its inputs, outputs, decisions and checks.

Illustrative, not real: the w189e activity lines on v1, and the shortened `why --json`
output on the Console board.

## §10 Open questions

1. **Merge semantics for World objects.** Git merges text. A World merge of two typed states
   needs a rule per object kind, plus the conflict surface. What is the smallest kind set that
   covers the merges actually seen (M4 was the charter file)?
2. **What starts a brane?** A charter row, a ruling, an episode with no parent brane? S1 can
   infer branes from charter rows; S3 has to decide.
3. **Cone size.** Upstream by ancestry grows without bound. Should the default cone be read
   and runs-on, with ancestry behind a toggle?
4. **The mission loop as a resident.** The loop's iterations form a brane. Should its states
   show at the same zoom as attended episodes, or fold into one step per row?
