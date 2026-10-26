# w-world-surface-outputs (row 188 follow-on: the front card's outputs, and work in other programs) — r0

**Status:** r0, drafted with Mark attended on 2026-10-26; **not yet approved.** Builds on
[w-world-os-surface.md](w-world-os-surface.md), whose §3.2 says "Outputs open in place". This
design delivers that and extends it: every file type, a tabbed front card, and a roadmap for
work done in other programs. The sequence slots into row 188 S1
([sprint plan](w-world-os-surface-sprint-plan.md)) after N2b (#367): **N2d** (data) and
**N4a** (the front card), ahead of N4's style pass. Rulings R1–R5 are in §9.

## §0 What Mark asked (attended, 2026-10-26)

This came from trying the N2b app over the tailnet. In order (typos corrected):

1. *"in the mockups we included being able to view the file outputs not just the file names."*
2. *"would it make sense to have inputs on left, outputs on right, what happened in middle?"*
3. *"and the A2UI space looks cramped for any real outputs the agents may use it for?"*
4. *"we should support all file types, video as well. perhaps even file types that launch the
   subsequent programs — but that can go on roadmap. trying to think of workflows such as using
   blender and blender files, godot etc — how ailang world makes it seamless"*
5. *"we probably won't be able to fit it all in these panes so perhaps tabbed UI will be needed
   for larger components"*

## §1 Problem (measured on the N2b app, #367)

- **Outputs are names only.** N1's `stack.json` gives each file as `{path, tag}`, where the tag
  is `edit` or `new`, so w189e's front card lists seven file names and nothing behind them. The
  artifact board opened its real diffs in place (README +1, QUICKSTART +3, the new page and its
  test). The v3 working copy dropped that, because its frozen data had no diffs, and N2b ported v3.
- **Two columns tell one story.** The card has four columns (Inputs 1.5fr; Outputs, Decisions
  and Checks 1fr each). w189e shows "round 1: 8/8 graders" under Decisions and "Passed all 8
  checks in one round" under Checks: one event, told twice, in two places.
- **The agent's view has no room.** "From the agent" is a single dashed line between the title
  and the columns. That suits the empty state, which is every state today, since row 151 is not
  built. It cannot hold a real A2UI view: a table, a chart, or a decision with buttons.
- **Only text is imaginable.** There is no notion of images, video, audio or binary outputs.
  World's own repo tracks no binaries (the tracked-binary gate), but the repos World is meant for,
  such as a Blender or Godot project, are mostly binary.

## §2 The front card: in → what happened → out

Three columns, read left to right, matching the timeline, where left is earlier:

| Inputs (1fr) | What happened (1fr) | Outputs (1.25fr) |
|---|---|---|
| the agent, the task, the grants; *came from* | one ordered list: who, what, when | the agent's view first (§5), then files, media and other outputs; *led to*; the PR |

**What happened** merges Decisions and Checks into the episode's story, oldest first. Each entry
has an actor (You, World, the agent, CI), a sentence and a time. Totals (graders, rounds,
agent-minutes, World calls) shrink from four tiles to one footer line. Here is w188k, measured
on 26 Oct:

| who | what |
|---|---|
| You | minted for 12 h, 10 grants · 08:14 |
| World | round 1 · 5/8 graders · checks could not run (the exec node was broken, friction #23) |
| You | restarted the daemon; asked for round 2 |
| World | round 2 · 7/8 · the one red fails on dev too |
| You | "get back to root by clicking on icon or world" |
| World | round 3 · 7/8 · 44/44 app checks |
| CI | running on ubuntu-24.04 |

## §3 Tabs: the card is a workspace

The panes can't hold a video, a long diff or an agent's own UI (§0.5). The front card therefore
becomes a tabbed workspace under its title bar:

- **Summary** is always the first tab and holds the three columns (§2).
- **Opening an output** (click or Enter on a file, a medium or the agent's view) adds a tab
  that uses the whole card body. Tabs carry the file name and a × to close them. A tab that is
  already open is focused, not opened twice.
- **The Esc ladder:** a panel (keys or look) closes, then the active tab, then the brane. Each
  press undoes the most recent opening.
- **Maximise** (`f`, or a button on the tab): the open tab takes the whole stage, and the stack
  bars and the rail step aside. This is the room for video and agent views. `f` or Esc restores.
- **The URL holds the open tab:** `/surface?brane=&at=&open=<path>` (or `open=view` for the
  agent's view), so a link lands on that diff or video. This extends `/surface`'s query grammar
  (R3). `/surface/static` keeps its own.
- **Moving to another state** (↑/↓, the rail, a chip) closes that state's tabs. Whether a state
  should remember its tabs is open (§11).

## §4 Every file type

Each output opens in the viewer for its kind. The kind comes from the extension, checked against
the bytes' magic number; on a disagreement the file is treated as binary.

| kind | examples | the tab shows |
|---|---|---|
| text, code | `.go`, `.js`, `.gd`, `.tscn`, `.tres`, `.json`, `.yaml` | the diff (line numbers, +/−), with a toggle to the whole file; long files show the first 2,000 lines and then "N more" |
| markdown | `.md` | rendered, with a toggle to the diff; no raw HTML passthrough |
| table, tree | `.csv`, `.tsv`; `.json` | a table, or a collapsible tree; a toggle to the diff |
| image | `.png`, `.jpg`, `.gif`, `.webp`, `.avif`, `.svg` | before and after side by side, with dimensions and sizes; SVG only ever as an `<img>` |
| video | `.mp4`, `.webm`, `.mov` (H.264) | `<video controls>` with seeking; before and after when both exist |
| audio | `.wav`, `.mp3`, `.ogg` | `<audio controls>` |
| PDF | `.pdf` | the first page as a preview (Tier 2, §6.1); download |
| binary, unknown | `.blend`, `.glb`, `.fbx`, `.zip`, `.pck` | a card: type, size before and after, sha256 before and after, the step and World call that wrote it, download |
| deleted | any | the before, marked deleted |

Godot's scenes and resources (`.tscn`, `.tres`) and its scripts (`.gd`) are text, so a Godot
project gets real diffs from the first row of this table.

**Where the bytes come from:**
- **A state's outputs:** its commit (`<at>:<path>`). The *before* is the first parent's copy. A
  step's commit is its episode's merged commit; a shared state's is the merge.
- **Reads:** through `host/worktree.Output`, the scrubbed git path that AC10 requires.
- **Outputs that never reached git** (an episode's uncommitted writes): the store holds their
  sha256 but not their bytes, so the tab says **"bytes not kept"**. Whether World should keep
  them is R5.

**The routes (N2d):**
- `GET /surface/file.json?brane=&at=&path=` returns:
  - `path`, `tag`, `kind` and `mime`;
  - `size` and `sha256`, each `{before, after}`;
  - for text, the diff hunks, capped, with `truncated`;
  - for kinds the browser cannot render, the reason.
- `GET /surface/blob?brane=&at=&path=&side=before|after` returns the bytes:
  - with the exact allowlisted `Content-Type`, `nosniff` and `Content-Security-Policy:
    default-src 'none'; sandbox`;
  - `inline` for allowlisted media, otherwise `attachment`;
  - Range support (`http.ServeContent`) so video seeks;
  - a size cap.
- **Only paths listed in that state's outputs are served.** Anything else is a 404. The blob
  route is not a way to browse the repo (§7).

## §5 The agent's view (A2UI) is an output

- **None yet** (every state today): one quiet line at the top of Outputs, "From the agent · no
  view yet · row 151". The band between the title and the columns goes.
- **Present:** a row at the top of Outputs with the view's title and kind. It opens as a tab
  and can be maximised, which gives a real table, chart or decision the whole card or stage.
- **The rules are unchanged** from w-world-os-surface §3.2: the view is checked against World's
  closed catalog, and its buttons only propose. Your answer still goes through the decision path.

## §6 Work in other programs (roadmap)

The workflow to design for: *the agent animates a walk cycle in `walk.blend` → you watch the MP4
in the front card → open the file in Blender → fix the hand → save → the agent continues from
your version.* Every hop is attributed and on one brane. Godot is the same with a level instead
of a walk cycle.

### §6.1 Preview derivers (Tier 2)

Kinds the browser can't render get a preview made by a World tool. The tool runs under an exec
profile, as `go test` does today, and writes the preview as a recorded effect with provenance.
Previews are cached by (sha256, deriver version). Candidates:
- **Blender:** the thumbnail that Blender already stores inside a `.blend`; otherwise a frame
  rendered with `blender -b <file> -o … -f 1`.
- **Godot:** a `godot --headless` screenshot or web export.
- **glTF/GLB:** a model viewer whose script is served from `'self'`.
- **PDF:** the first page rendered to an image.
- **HEVC `.mov`:** transcoded to H.264.

### §6.2 Agents drive the programs headless

Exec profiles per program, confined the way the `ailang-world` profile confines `go`: declared
binaries, read roots, caches, and no network.
- **`blender`:** `blender -b --python <script>`.
- **`godot`:** `--headless`, exports and test runners.

The agent's work in Blender is then World calls like any other, so it has budgets, refusals and
provenance.

### §6.3 Open in the program (the desktop shell)

The browser can't launch programs; N3's Tauri shell can.
- "Open in Blender" on an output launches a native program, which is an effect. It needs a
  grant (proposed: `Desktop.Open`, scoped to the worktree) and an approval like a mint's (row
  201's approvals by key would cover it).
- The browser shows the action disabled, with "in the desktop app".

### §6.4 The round trip: your edits are steps

While a file is open in its program, World watches that episode's worktree (or snapshots it
when you press "Done editing").
- **Your save becomes a step:** "You · edited walk.blend", with the before and after sha256,
  attributed to your approval key rather than to the agent.
- **The next agent round starts from that state.** Nothing is copied between tools, and nothing
  goes missing from the history.

## §7 Invariants

- **Read-only surface.** None of these routes writes. Opening in a program is the desktop
  shell's effect, behind a grant and an approval.
- **Outputs only.** The file and blob routes serve exactly the paths a state lists as outputs.
- **Absent stays absent.** Bytes that weren't kept, and kinds with no viewer, say so; nothing
  is invented.
- **The CSP stays strict:**
  - no inline script, and scripts only from `'self'`;
  - media only from `'self'`, and only through the blob route;
  - blobs served under `default-src 'none'; sandbox`.
- **World's own repo stays binary-free.** Media is served from the git objects of the repo being
  worked on, and the tracked-binary gate is untouched.

## §8 Proposed sequence

1. **N2d (Codex, data):** `file.json` and `blob`. It covers kinds and the magic check, diff
   hunks with caps, outputs-only paths, Range, size caps and every refusal (`400`/`404`/`405`/
   `413`). The CSP is per R2. Graded by Go tests plus a browser check that an image and a video
   load under the new CSP.
2. **N4a (Claude, the front card):** §2's three columns and the What-happened list; §3's tabs,
   with the Esc ladder, maximise and `open=`; §4's viewers; §5's agent view as an output. Graded
   by the browser grader, with the real w189e diffs and a media fixture.
3. **N4 (Mark-led style pass):** as planned, over the new card.
4. **Roadmap rows (on R4):** preview derivers (§6.1); program exec profiles (§6.2); open in
   the program (§6.3, after N3 and row 201); the human-edit round trip (§6.4).

## §9 Rulings for Mark

- **R1, layout:** Inputs | What happened | Outputs, with Decisions and Checks merged. Default:
  yes.
- **R2, CSP:** add `img-src 'self'; media-src 'self'` to `surfaceCSP` for the blob route. No PDF
  framing; PDFs go through §6.1 instead. Default: yes.
- **R3, tabs:** a tabbed front card, and an `open` parameter in `/surface`'s query grammar.
  Default: yes.
- **R4, roadmap:** file §6's four rows, with positions left to the loop's groom. Default: yes.
- **R5, bytes:** keep the bytes of outputs that never reached git (a content-addressed blob
  store in World), or show "bytes not kept". Default: not kept, until a row decides.

## §10 Sketch

A static sketch of the tabbed front card, built from real data: w188k's round-3 diffs, the
grader's timeline screenshot before and after round 3, and a video tab using the World intro
film as an example only (it is not one of w188k's outputs). It is served on the tailnet for
review and committed as a text-only fixture once reviewed, with its media referenced, not
committed.

## §11 Open questions

1. **Tab memory.** Should a state remember its tabs when you come back to it, and how many
   tabs before the oldest closes?
2. **Big binaries.** LFS pointers, multi-GB `.blend` files: does the blob route follow LFS, and
   where is the cap?
3. **Diffing binaries.** An image diff (onion skin, swipe) beyond side by side? A `.blend`
   diff at the object level (which objects changed), from a deriver?
4. **Watching the worktree** (§6.4). Is a file-system watch or an explicit "Done editing"
   snapshot the human-edit step? A watch is seamless but noisy; a snapshot is one deliberate act.
