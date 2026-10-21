# w-world-os-surface — S1 sprint plan (row 188, D-WORLD-87)

**Status:** draft plan, attended 2026-10-20, for Mark's review before M0 starts.
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
| **M0b** | w188b | Reads in the episode report, from World's invocation records: `EpisodeReport.reads` = `{count, files[{path, count, first_entry}], truncated, unlisted}`. A read is a committed `ailang-read` call with a v2 record whose effects were allowed and did not fail and whose handler result is not `ok:false`. Its path is the record input's `path`; a ranged read counts. Only `filepath.IsLocal` paths are listed, the rest count in `unlisted`. Files are ordered by first entry, capped at 1000 (`truncated`). The text form adds a `reads:` line | `go test` cases for the read rules (counts and order, ranged, handler- and plan-refused, non-local, the cap, an empty `files: []`, the text line); every other report field unchanged against the running daemon on six episodes; a live check of real episodes' `reads` against independently measured read sets |
| **M0c** ✅ | w188c | **Done.** The episode facts, added to the same index only when it is built with a rig (without one the output is byte-identical to M0a's). `host/surface/episodes.go` (pure): `Step`, `Graders`, `FileChanges`/`Change`, `Read`, `Decision`, `RigSource` (`source.rig`: episodes, attached, unattached) and `source.reads_method`; `ParseEpisodeRecord` (`trial-episode/1`, path fields never decoded), `GroupRounds` (one record per round, ordered by `started_at`), `RoundDir` (`<project>`, `<project>-rN`), `ParsePRs` (only `"<project>#N"`), `ParseReport`, `ParseNameStatus`, `ParseTree`, `ClassifyLastChange` and `AttachRig` over a `Git` seam. `steps[]`: one per episode whose PR is a step state (task, agent, model, base from round 1; graders from the last round; minutes summed over rounds; calls and entries from the report, absent on a 404). `files{}`: `git diff --name-status --no-renames -z <base> <step>`, capped at 1000. `reads{}`: each report path as `new` (not at base), `step` (the last change at base is a merge of a step in the index), `commit` or `outside`, via `git ls-tree` once per step and `git log --first-parent -1` per path. `decisions[]`: each round's grade, the report's refusals, the merge (`you`, or `loop` on a loop brane). `Validate` gains the rig rules. `ailang-worldd surface-index` gains `--rig`, `--project`, `--addr` (loopback only; bounded GET, 4 MiB; records ≤ 1 MiB, ≤ 2000). `ran_on{}` stays absent (no source before S3). **Reads use the last change at the step's base, not the creating commit**: that is what reproduces the design's M8, where w189e's reads come from six earlier sources on three branes; a creator-only mapping would send every long-lived file to its first author | `go test ./host/surface/...`: a fixture rig of 2 episodes (one with two rounds and a refusal) against the golden `testdata/index-rig.json`, the read cases, the 1000 cap, and one `Validate` mutant per rig rule; a real build reproduces §4 (w189e has 2 lines behind it) |
| **M1** ✅ | w188d | **Done.** `worldd serve --surface-index <file>`; startup refuses files or directories under the workspace root, resolving symlinks. `GET /surface/index.json` re-reads and validates each request with the shared `surface.MaxBytes` 4 MiB bound; serves exact bytes with JSON, no-store and nosniff headers. Unconfigured or absent → positive empty state; read/decode/validation failures → 503 without the configured path. | Handler recorder tests: exact bytes and headers, both empty states, rebuild, oversize, decode and duplicate-state refusal, read errors, path redaction, POST 405, and startup refusal directly and through a symlinked directory; frozen-table and dependency guards |
| **M2** ✅ | w188e | **Done.** `GET /surface`: the world timeline, server-rendered — lanes, glyphs, merges, the now edge with health and *needs you*, ticks and legend; Paper, Slate and Night tokens. `host/surfacepage` (standard library and this repo): `LayoutTimeline` (pure; shared lane first, then branes newest first; six glyphs then slate; x by rank of `at`, 28 px; one UTC tick per day), `RelativeTime`, `RenderTimeline` (`html/template`, the AILANG mark inlined, no script). The daemon route shares M1's read-and-validate (`readSurfaceIndex`), uses the workbench's headers, and answers any query with an inert 400. The now edge's health line is the daemon's build and the index's validation, and its sessions are the credentials store's live rows, read-only. **Needs you renders "Nothing needs you" in S1: mints and approvals waiting on Mark are not recorded anywhere the daemon can read (S2/S3 add them), so AC7 waiting is not shown and nothing is invented in its place.** | layout goldens (`LayoutTimeline` is deterministic and has no map-order dependence); one `<script` and no `on*=`; local links only; the tokens test; the row 189 fixture shows w189d's badge "3" and AC7 waiting |
| **M3** | (not yet minted) | `GET /surface?brane=&at=`: one brane — earlier stack above, later below (one column per brane, three deep), the open card's four columns, relative and exact times; a closed query grammar `{brane, at}` | stack truth table on the fixture: w189e has 2 columns behind and 2 ahead; w189d 3 behind; `0ee1945` 3 ahead. Unknown key → inert 400. Every field in the view types renders (field census). |
| **M4** | (not yet minted) | `/surface/surface.js`: the keyboard map, address bar, zoom, Appearance (stored in this browser only), and live refresh from `/agui/` by re-fetch and a fixed region swap | a node harness, as with `live_script_test`: keys move as §3.3a says; typing in the address bar is ignored by the shortcuts; unknown addresses keep the page; the region list is pinned; `STATE_DELTA` paths outside `/lastIndex` and `/logHead` are refused |
| **M5** | (not yet minted) | `docs/QUICKSTART.md` "See your work" section; the plugin's `/world-watch` points at `/surface`; `ailang-worldd surface-index` documented | help-check for every command on the page (as `TestUseWorldPageCommandsMatchTheCLIHelp`); the docs link and stale-phrase checks |

**S1 is done when** Mark opens `/surface` on his machine, built from the real index on `dev`,
and the row 189 brane matches the mockup's facts:

- the two-ancestor merge;
- w189d's three lines behind it;
- `0ee1945`'s three columns ahead;
- AC7 waiting.

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
