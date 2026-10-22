# w-world-os-surface — S1 sprint plan (row 188, D-WORLD-87)

**Status:** M0a–M4a landed or in review (2026-10-20/21). **Re-cut 2026-10-21 by D-WORLD-88:** the rest of S1 is §"S1 re-cut" below. M4b and M5 in the table are superseded by N2 and N5.
**Design:** [w-world-os-surface.md](w-world-os-surface.md) r2 (§3 the surface, §5 the gaps, §7 S1).
**Mockup:** [mockups/world-os-surface.dc.html](../mockups/world-os-surface.dc.html).
**Base:** ailang-world `dev` `adffe4c1`.
**How it is built:** attended, as World episodes, the way row 189 was. Each milestone is one
episode on the trial rig (`~/.ailang/state/trial-2026-10-16`), minted by Mark with `episode
start`, graded by a rig grader, then a PR, merged on Mark's word.

**Size:** eight milestones (M0 split into M0a–M0c). Each is at most ~30 minutes of agent work and ~250 lines.

## What S1 delivers

The two distances, read-only, over data that exists today:

- **The world timeline:** branes as lanes, merges, the now edge, and *needs you*.
- **One brane:** the earlier and later stacks, and the front card holding inputs, outputs,
  decisions and checks.

Both are served by `worldd` at a new `GET /surface`, with the three themes, the keyboard map and
the address bar.

**Not in S1:**
- live inputs and "change from here" (S5; needs the fork API, R2 on hold);
- the agent band (row 151);
- diffs and images in Outputs (S2);
- brane data inside World's store (S3).

## Constraints the build must meet (measured 2026-10-20 on `dev`)

1. **One script, strict CSP.** The workbench's CSP is pinned as a literal by
   `TestWorkbenchSecurityHeaders`: no inline script, no `on*=` handlers, same-origin links only.
   `/surface` copies that header and serves one script of its own, `/surface/surface.js`
   (`go:embed`, as `host/workbench/live.js` is).
2. **No new routes in the frozen table.** `/v1` is pinned by `TestFrozenV1RouteTableMatchesMux`.
   `/surface` and its two companions register beside `/workbench`, outside the frozen table, and
   are GET-only.
3. **No new Go dependencies.** `TestDaemonDependencyAllowlist` admits only the standard library,
   this repo, and the store's sqlite. The page is `html/template` plus pure layout functions, the
   same way `host/workbench/graph.go` lays out SVG.
4. **Its own tokens.** The workbench's 11-token, two-theme test (`TestTokensBothThemes`) is not
   touched. `/surface` has its own token set (Paper, Slate, Night) with its own test.
5. **The daemon has no git, charter or rig input today.** Its only filesystem inputs are the
   store, the credentials store, the workspace root and exec profiles. S1 adds **one read-only
   input**: a generated surface index (M0a, M0c, M1). S3 later replaces it with World objects.
6. **Loopback only.** GET routes are unauthenticated (row 199, residual R1), so `/surface` is
   for the local machine until row 199 lands.

## The surface index (`world/surface-index/v1`)

A JSON document, bounded at 4 MiB, built outside the daemon by `ailang-worldd surface-index`
and read by it. M0a builds `branes[]` and `states[]` (brane kinds `attended`, `loop`, `mixed`;
state kinds `shared`, `step`); M0c adds the rest, taking each step's reads from the episode
report M0b extended.

| field | from | example |
|---|---|---|
| `branes[]`: key, title, row, kind (`attended` · `loop` · `shared`) | merge branch names (`trial/row189-…`, `sprint/world-iter288-row156-…`, `design/row189-…`) and the charter's queue row titles | `onboarding · row 189` |
| `states[]`: id, brane, kind (`step` · `shared` · `open`), merge hash, PR, parents[], landed time | `git log --first-parent --merges` on `dev`, and the merge-base of each PR | `9003555`, parents `7fb99da` and `f36428a` |
| `steps[]`: episode, task file, agent, rounds, graders x/y, agent minutes, calls, log range, base | the rig's `evidence/<project>/episodes/<ep>.json` and `prs.json` | w189e, 8/8, 3.5 min, 42 calls |
| `ran_on{}`: step → shared state | **none before S3**: no record says which daemon build an episode ran on, so S1 leaves it absent | (absent) |
| `reads{}`: step → states | the episode report's `reads` (World's invocation records), each path mapped to its last change on the step's base's first-parent line (`source.reads_method` = `last-change-at-base/v1`; M0c) | w189e → design, w189b, w189c, w189d, w198, row 138 |
| `files{}`: state → [path, new/edit] | `git diff --name-status` per merge | w189e: 7 files |
| `decisions[]`: actor, text | the episode report (mint, refusals), the merge's PR, ledger review rounds | You · Merged #329 |

Every field names its source in the index, so a value the builder could not measure is absent,
never invented, and the card renders it as missing.

## Milestones

**M0 is split in two** (w188a, 2026-10-20). The episode sandbox runs Go only (`go-test`,
`go-vet`, `gofmt-l`; no git, no Python), so the builder is pure Go plus a thin git subcommand,
not a Python script. The pure functions are tested from fixtures inside the sandbox, and the
grader runs the subcommand against the real repository:

- **M0a**: the index from git history alone (the shared world, the steps, their branes).
- **M0b** (w188b, 2026-10-20): the files each episode read, in the episode report's `reads`,
  from World's own invocation records rather than agent transcripts, so the source is the
  same for every agent.
- **M0c**: the episode facts from the rig (steps' episodes, files, decisions) and `reads{}`
  mapped to states, built from the rig's episode records and those reports. `ran_on{}` has no
  measured source before S3, so S1 leaves it absent.

| M | Episode | Delivers | Grader checks (a rig script, written before the run) |
|---|---|---|---|
| **M0a** ✅ | w188a | **Done.** The index from git history alone. `host/surface` (standard library only): the type `Index`; `ParseMerges`, `BraneOf`, `ChartRowTitles`, `Build` (deterministic) and `Validate`. Plus `ailang-worldd surface-index --repo --charter --out [--rev dev] [--limit 200]`, the thin git side (git log, then git merge-base per merge), which writes atomically and never inside the repo, and `surface-index --validate <file>`. Fixtures: 24 real `dev` merges ending at `adffe4c1`, their merge-bases, and charter rows 101, 156, 188, 189 and 198. The golden has 24 shared states, 24 steps and 6 branes. | `go test ./host/surface/...`: parse errors by line, `BraneOf` for every branch form, row slugs (188's tags, 101's none), the golden, and `Validate` refusing 18 mutants. A real `surface-index` on `dev`: `9003555` is shared with parents `7fb99da` and `f36428a`; `f36428a` is on `row-189` (`w-world-onboarding`) with base `0ee1945`; #328's step is on `row-156` (loop), with base `0ee1945`; #317's step has no brane; no path or email appears in the file. `--validate` accepts the index and refuses the grader's mutants |
| **M0b** ✅ | w188b | **Done** (#336). Reads in the episode report, from World's invocation records: `EpisodeReport.reads` = `{count, files[{path, count, first_entry}], truncated, unlisted}`. A read is a committed `ailang-read` call with a v2 record whose effects were allowed and did not fail and whose handler result is not `ok:false`. Its path is the record input's `path`; a ranged read counts. Only `filepath.IsLocal` paths are listed, the rest count in `unlisted`. Files are ordered by first entry, capped at 1000 (`truncated`). The text form adds a `reads:` line | `go test` cases for the read rules (counts and order, ranged, handler- and plan-refused, non-local, the cap, an empty `files: []`, the text line); every other report field unchanged against the running daemon on six episodes; a live check of real episodes' `reads` against independently measured read sets |
| **M0c** ✅ | w188c | **Done.** The episode facts, added to the same index only when it is built with a rig (without one the output is byte-identical to M0a's). `host/surface/episodes.go` (pure): `Step`, `Graders`, `FileChanges`/`Change`, `Read`, `Decision`, `RigSource` (`source.rig`: episodes, attached, unattached) and `source.reads_method`; `ParseEpisodeRecord` (`trial-episode/1`, path fields never decoded), `GroupRounds` (one record per round, ordered by `started_at`), `RoundDir` (`<project>`, `<project>-rN`), `ParsePRs` (only `"<project>#N"`), `ParseReport`, `ParseNameStatus`, `ParseTree`, `ClassifyLastChange` and `AttachRig` over a `Git` seam. `steps[]`: one per episode whose PR is a step state (task, agent, model, base from round 1; graders from the last round; minutes summed over rounds; calls and entries from the report, absent on a 404). `files{}`: `git diff --name-status --no-renames -z <base> <step>`, capped at 1000. `reads{}`: each report path as `new` (not at base), `step` (the last change at base is a merge of a step in the index), `commit` or `outside`, via `git ls-tree` once per step and `git log --first-parent -1` per path. `decisions[]`: each round's grade, the report's refusals, the merge (`you`, or `loop` on a loop brane). `Validate` gains the rig rules. `ailang-worldd surface-index` gains `--rig`, `--project`, `--addr` (loopback only; bounded GET, 4 MiB; records ≤ 1 MiB, ≤ 2000). `ran_on{}` stays absent (no source before S3). **Reads use the last change at the step's base, not the creating commit**: that is what reproduces the design's M8, where w189e's reads come from six earlier sources on three branes; a creator-only mapping would send every long-lived file to its first author | `go test ./host/surface/...`: a fixture rig of 2 episodes (one with two rounds and a refusal) against the golden `testdata/index-rig.json`, the read cases, the 1000 cap, and one `Validate` mutant per rig rule; a real build reproduces §4 (w189e has 2 lines behind it) |
| **M1** ✅ | w188d | **Done.** `worldd serve --surface-index <file>`; startup refuses files or directories under the workspace root, resolving symlinks. `GET /surface/index.json` re-reads and validates each request with the shared `surface.MaxBytes` 4 MiB bound; serves exact bytes with JSON, no-store and nosniff headers. Unconfigured or absent → positive empty state; read/decode/validation failures → 503 without the configured path. | Handler recorder tests: exact bytes and headers, both empty states, rebuild, oversize, decode and duplicate-state refusal, read errors, path redaction, POST 405, and startup refusal directly and through a symlinked directory; frozen-table and dependency guards |
| **M2** ✅ | w188e | **Done.** `GET /surface`: the world timeline, server-rendered — lanes, glyphs, merges, the now edge with health and *needs you*, ticks and legend; Paper, Slate and Night tokens. `host/surfacepage` (standard library and this repo): `LayoutTimeline` (pure; shared lane first, then branes newest first; six glyphs then slate; x by rank of `at`, 28 px; one UTC tick per day), `RelativeTime`, `RenderTimeline` (`html/template`, the AILANG mark inlined, no script). The daemon route shares M1's read-and-validate (`readSurfaceIndex`), uses the workbench's headers, and answers any query with an inert 400. The now edge's health line is the daemon's build and the index's validation, and its sessions are the credentials store's live rows, read-only. **Needs you renders "Nothing needs you" in S1: mints and approvals waiting on Mark are not recorded anywhere the daemon can read (S2/S3 add them), so AC7 waiting is not shown and nothing is invented in its place.** | layout goldens (`LayoutTimeline` is deterministic and has no map-order dependence); one `<script` and no `on*=`; local links only; the tokens test; the row 189 fixture shows w189d's badge "3" and AC7 waiting |
| **M3** ✅ | w188f | **Done.** `GET /surface?brane=KEY&at=ID`: deterministic `LayoutStack`, mirrored earlier/later bars, four-column front card and whole-brane rail. Own lane first (three before/after); shared states collect merging/leaving steps; steps collect other branes across the previous-base/current-base shared window (base alone for a first or missing previous base), and lead through their merge to shared and leaving lanes. Columns preserve first appearance; cards use time then id. Index reads, files, decisions and episode checks render without invented values; closed grammar accepts exactly one `brane` and `at`, malformed queries return 400 and missing selections framed 404. | Five fixture truth tables cover merge, first step, cross-brane window, downstream work and shared departure; determinism and every-state rendering check markers, local links and no script. Route tests cover stack 200, malformed/repeated/partial query 400, unknown/wrong-brane 404 and POST 405. |
| **M4a** ✅ | w188g | **Done.** `GET /surface/surface.js` (`host/surfacepage/surface.js`, `go:embed` as `Script`; `/workbench/live.js`'s headers), the one `<script src defer>` on both pages; without it the pages work as before. Pure seams, exported for node: `laneStates` (M2/M3's lane rules, by `at` then id), `parseAddress` (`world`, `back`, `world NAME[@N\|@HEX]`; NAME is a key, a title or a hyphen-joined title suffix, `shared`/`dev` the shared lane; ambiguous, unknown, out-of-range and non-unique-prefix addresses are errors that never navigate), `keyAction` (§3.3a for both views over `buildModel`'s fixed model shape), `shouldHandle` (no shortcut while typing), `buildModel`, `loadLook`/`saveLook`/`resolveLook` (`world.appearance`; Paper, Slate, Night, Match system; Shape and colour or Shape; Full, Reduced or the system's). The wiring: one `keydown` listener, focus in `--you`, the address bar in the header (server-rendered address, the index fetched once), the keys overlay, *needs you*, and the Appearance panel (gear or `,`). Two layout fixes: the timeline's label column is sized from the lane names (ellipsized past 480 px, the full name in the lane's `<title>`); a stack decision's pill, text and time share one flex line, and only the text wraps | a node harness (`TestSurfaceScriptPure`, skipping only without node; fails when `WORLD_EXEC_NODE` is set and it cannot run): every address form and error, every key in both views, `shouldHandle`, throwing and garbage storage, `resolveLook`, `buildModel` on the rendered rig pages against Go's own layouts, every stack's address round-trip, and the wiring in a fake DOM. Go: the route's headers, exactly one script and no inline script on both pages, the seams in the asset |
| **M4b** | (not yet minted) | live refresh from `/agui/` by re-fetch and a fixed region swap, in `/surface/surface.js` | a node harness, as with `live_script_test`: the region list is pinned; `STATE_DELTA` paths outside `/lastIndex` and `/logHead` are refused |
| **M5** | (not yet minted) | `docs/QUICKSTART.md` "See your work" section; the plugin's `/world-watch` points at `/surface`; `ailang-worldd surface-index` documented | help-check for every command on the page (as `TestUseWorldPageCommandsMatchTheCLIHelp`); the docs link and stale-phrase checks |

## S1 re-cut: the Tauri app on the artifact's design (D-WORLD-88, attended 2026-10-21)

**What happened.** M2 to M4a built the surface as server-rendered pages. They are correct, and
graded state by state, but Mark found them "functionally the same but not great". He asked for "an
operating system full screen … somewhere I want to live in". A same-data spike compared Tauri v2
(Claude) with SwiftUI (Codex), and D-WORLD-88 chose Tauri:
- one HTML interface for both the native window and the browser on his tailnet;
- self-hosted fonts under `/surface`'s own CSP;
- the approved artifact as the visual target, ported from its own source.

Two working copies followed, judged by Mark on the tailnet:
- **v2:** the artifact's markup, styles and layout maths with real data;
- **v3:** timeline pan, zoom and an overview strip, and a human layer: titles, the agents' own summaries, brane goals and progress, and a Today line.

Mark: "starting to look good … we will iterate more im sure but this is good to build on top of".

**The working copy** lives outside the repo, in the session scratchpad (`spike/tauri-v3`):
- `dist/`: plain HTML, CSS and JS, with fonts and the logo;
- the Rust shell;
- `make-app.sh`;
- `spike-data.json`: frozen real data, recording the source of every field in a `provenance` block.

N2 brings the interface into the repo; N3 brings the shell.

**What stays from M0–M4a.** The index (M0a–M0c), the serving route (M1) and the stack rules (M3's
`LayoutStack`) stay. So do M4a's address grammar and Appearance seams, which the app reuses. The
server-rendered timeline and stack pages stay as the no-script fallback and the graders' structural
fixture, until N2's model route carries the same rules.

**Principle: human first.** Every field a person reads must come from the records. Nothing is
invented, and each field's source is named. A missing field shows as an honest gap, which becomes a
data milestone. Identifiers (hashes, PR numbers, episode ids) appear on hover and in the detail
view, never as the main label.

| N | agent | delivers | grader checks |
|---|---|---|---|
| **N0** ✅ | Codex | **The index's human layer.** Step `title` from the merge commit's body line (GitHub puts the PR title there; git alone, no network). Step `summary`: the agent's own round-1 summary from the rig record, first sentence plus the full text, bounded. Brane `name` and `goal`, the charter row's bolded problem statement. The **delivered verdict**: the last round's graders, plus the operator's recorded re-runs (`graders/<ep>.<grader>.regrade.json`, a rig convention made explicit). `Validate` rules for each. | real build: w189e's title "The Use World on your repo page" and summary "I've written the …"; row 189's name and goal; w188d/w188e/w188f delivered 8/8 with their re-runs; no path or email; titles never invented (a merge without a PR has none) |
| **N1** | Claude | **`GET /surface/model.json`:** the view model in Go, from the index plus the daemon's live facts. Lanes in order, with glyphs and progress; states with labels; links; every state's stack from `LayoutStack`; the now edge (health, build, live sessions, running episodes, Today). Bounded, deterministic, loopback, outside the frozen `/v1` table. The JS renders; the rules stay in tested Go. | every state's stack equals `row188_stack_rules.py` on the real index; Today matches the episode records; an index refusal comes back as a reason |
| **N2a** | Claude | **Fonts and the surface's CSP (D-WORLD-90).** A narrow second allowance in `scripts/tracked_binary_gate.sh`: font files only under `host/surfacepage/fonts/`, each with a valid font signature (`00 01 00 00`, `OTTO`, `wOF2`), under a per-file and a directory cap, and listed by exact sha256 in the gate. Montserrat and IBM Plex Mono (OFL) are committed there with their licence texts, at the pinned google/fonts commits the fixture's `fetch-fonts.sh` uses. `GET /surface/fonts/<name>` serves them, embedded. `/surface`'s own CSP is the workbench's plus `font-src 'self'`. | the gate's own test suite gains arms for the fonts (accepted), a renamed binary in the fonts directory (bad signature), an unlisted font (not allowlisted), an over-cap font and a font elsewhere in the tree; the real tree passes; the font route serves the exact bytes with the right type; the CSP literal |
| **N2b** | Codex | **The app shell at `/surface`.** The working copy's `dist/` is embedded (`go:embed`) and reads `/surface/model.json`. It uses N2a's CSP and fonts. Navigation (pan, zoom, overview), M4a's address grammar and keys, Appearance, and **live refresh over AG-UI** (M4b's rules: deltas only for `/lastIndex` and `/logHead`, then a re-fetch of the model). | a node harness on the pure seams; headless Chrome screenshots of the timeline, a hover card and w189e's brane at 1600×1000 and 1600×932, reviewed by the operator against the working copy; one script, no `on*=`, local links; the CSP literal and font route; live refresh driven by a real commit on a store copy |
| **N3** | operator, then Claude | **The Tauri shell in the repo** (`desktop/`): a small Rust app whose window loads `http://127.0.0.1:<port>/surface` from the local daemon, or a tailnet address given at launch. `make-app.sh` builds `World.app`. The World sandbox has no `cargo`, so the operator builds it until the `ailang-world` exec profile gains a build command (Mark's yes). | the app launches and stays up; the window shows `/surface` from the daemon; no bundled data |
| **N4** | Mark-led | **The style pass:** iterated with Mark against the artifact rendered locally, compared screenshot to screenshot. | Mark's eye |
| **N5** | Claude | **Docs:** QUICKSTART's "See your work" (the app, the tailnet via `tailscale serve`, never Funnel); the plugin's `/world-watch` points at `/surface`. | help-check for every command on the page |

Operator re-runs are recorded beside the last round’s grader logs as `evidence/<dir>/graders/<ep>.<grader>.regrade.json` (`trial-regrade/1`, matching episode and grader, `passed: true`, `rc: 0`, at most 64 KiB); the index copies only the resulting verdict and grader name.

**Data milestones after S1** (each fills a gap the working copy shows):
- ran-on recording (S3);
- waiting and open steps (the AC7 kind);
- mutants caught (from judge records);
- agent views (row 151);
- approvals and mints as records, which also gives *needs you* a source (row 201, D-WORLD-89).

**Order:** N0 → N1 → N2a → N2b → N3, with N4 continuous from N2 on, then N5. Agents alternate between Claude and Codex (Mark, 2026-10-21).

**S1 is done when** Mark opens the World app on his machine, and his tailnet link, on the real index and live daemon, and lives in it. The row 189 brane must still match the mockup's facts:

- the two-ancestor merge;
- w189d's three lines behind it;
- `0ee1945`'s three columns ahead;
- w189e's six read sources on three branes.

AC7 waiting is shown only once open steps are recorded (a data milestone); until then the gap is shown honestly.

He then rules S2 on.

## Order and parallelism

M0a → M0b → M0c → M1 → (M2 ∥ M3 after M1) → M4 → M5. M1 needs only M0a's index to start, so
it can run beside M0b and M0c. M2 and M3 touch different templates and can run as two
episodes at once, since minting no longer stops the daemon (row 189 M1). Each PR's CI is about
25 minutes, so the critical path is about 4 PR cycles plus Mark's mints and reviews.

## Risks

- **Index drift.** The index is a projection. If it is rebuilt late, it shows an old world.
  M1 serves the index's own build time on the now edge, so staleness is visible, never silent.
- **Reads mapping.** "Read" edges come from the episode report's `reads` (World's invocation
  records, M0b), mapped to each file's last change on the step's base's first-parent line. A
  file edited by many steps maps only to the latest of them before the step began. M0c records
  the method in the index (`source.reads_method` = `last-change-at-base/v1`), and the card
  labels these edges as *read*, never as a cause.
- **Branch-name parsing.** Branes from branch names miss merges that follow no convention
  (`docs/…`). Those land on the shared lane with no brane, and the M0 fixture covers one.
