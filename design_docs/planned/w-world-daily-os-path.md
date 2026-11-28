# w-world-daily-os-path: the shortest path to using World as your OS for agent work — r0 (regroom proposal)

**Status:** r0, **attended 2026-11-28**, a regroom proposal for Mark. It files no new work row.
It orders rows that already exist, names one milestone, and asks two rulings (D-WORLD-130,
D-WORLD-131, §5). Nothing moves until Mark rules. The loop may not reorder an attended groom
(Gate 2 rule (c)).

## §0 What Mark asked

> *"ok so does this give me ability to use ailang world as a real OS?"* (2026-11-28, after row 274
> M1 and M1b landed)

The honest answer was not yet. Then: *"yes contineu with those"*, to drafting this proposal.

## §1 Two different "OS" goals

| | **Daily OS v1** (this proposal) | **OS from boot** (later) |
|---|---|---|
| What it means | World is where you start, watch, talk to, approve and review **all agent work**, full screen, without herdr or a shell for that work | You boot a machine into World. Network, devices, installs and your own edits are World effects |
| Runs on | the Mac you use today | a NixOS image (row 211) |
| Rows | 243 V4, 274, 188 N3–N5 + the outputs card, 222, 206, 207, 145 | 205, 208, 209, 210, 211, 212, 213–216 |
| Status | every row is designed; most are NEXT or in sprint | all PARKED post-1.x ([`w-world-os-from-boot`](w-world-os-from-boot.md)) |

DESIGN §2's boundary holds for both: Linux or macOS stays the kernel, and World is the semantic
layer that becomes the session.

## §2 Where each daily-OS piece stands (measured at `origin/dev`, 2026-11-28)

| Piece | Row | State |
|---|---|---|
| Start and finish an episode from the browser | 243 V4 | in sprint (S9c landed iter 450; S9d next; S10–S12 remain; D-WORLD-122, the Sol preset, is open before S11) |
| Apps: a format, checked, installable | 274 M1, M1b | M1 merged (#712); M1b-1 merged (#715); M1b-2 in PR (#718); M1b-3 (verbs read the registry) remains |
| **Launch an app from the surface** | 274 M2 | designed (published apps only, D-WORLD-129); waits for 243 V4 |
| The app's checked outputs in the report | 274 M3 | designed (D-WORLD-125 = A); with M2 |
| The World app in the browser | 188 N2a–N2c | landed (October; `/surface`, live over AG-UI) |
| **The desktop app, full screen** | 188 N3 | **not started** (no `desktop/` in the repo). Tauri, chosen by D-WORLD-88. It needs the operator to build it, because the agent sandbox has no `cargo` |
| Style pass, docs | 188 N4, N5 | after N3 |
| **Outputs you can open in the card** (files, images, diffs) | outputs doc N2d, N4a | partly built: `GET /surface/file.json` is routed; `/surface/blob` (images, video, binaries) is not, and the card's viewers (N4a) are not built |
| Your attended session is itself an episode | 222 | NEXT (attended track) |
| **Watch and type to a running agent** | 206 | NEXT, after 222 |
| Which agent is waiting on you | 207 | NEXT, after 222 |
| **Builds and tests longer than 9 s** | 145 | design ruled (D-WORLD-84); an evidence row "after the v1.0.0 tag" (D-WORLD-83), and only rc.1 exists |

**The four gaps that make World feel like a daily OS are in bold:** launching apps, a
full-screen app, live agents, and long-running commands. All four are designed. None is far
from being built.

## §3 Proposed order

Two tracks run side by side, as they do today. The loop keeps its queue, and the attended track
is Mark plus an attended session.

**Loop track** (unchanged first position):
1. **243 V4** to its end (S9d–S12). This is still the critical path.
2. **274 M1b-3**, small: `app list|show` and `/surface/apps.json` read the registry first.
3. **274 M2 and M3**: launch a published app from the surface start flow, and the report's
   checked outputs. **Proposed change:** ahead of rows 272 and 273. M2 needs only 243 V4 (since
   D-WORLD-123 moved device apps out, it no longer builds on 272's grants).
4. **272, then 273**, as filed.

**Attended track** (proposed new order, each item started attended):
1. **188 N3, the desktop app.** Small: a Tauri window that loads `/surface` from the local
   daemon, plus `make-app.sh` to build `World.app`. The operator builds it. Then N4 (style, with
   Mark) and N5 (docs).
2. **Outputs N2d and N4a.** Finish N2d with `/surface/blob` and its refusals (`file.json` is routed already), then
   the front card's viewers. This is what makes an app's results visible in place (M3 records
   them, and N4a shows them).
3. **222, then 206 and 207.** Attended sessions become episodes, then the live agent tab and
   agent states. This is the step that lets World replace herdr for agent work.
4. **145, the long-run lane**, re-gated from "after the v1.0.0 tag" to now (D-WORLD-131). A
   daily OS cannot cap every build and test at 9 s. World's own `go test ./...` takes about
   155 s.

**Stays parked (OS from boot):** 205, 208, 209, 210, 211, 212 and 213–216, until daily OS v1 is
met. Rows 211 and 215 already name 188 N3 as a dependency, so step 1 of the attended track also
unblocks the later path.

## §4 The milestone: "daily OS v1"

Met when Mark, on the Studio, can do all of these in one attended session without a shell or
herdr for the agent work:

1. open `World.app` full screen (188 N3);
2. launch the `docparse-convert` app from the surface, approve its start, and see it running
   (274 M2);
3. watch the agent live and type to it (206), and see that it is the one waiting on him (207);
4. let it run a check longer than 9 s inside World (145);
5. open its outputs in the card when it finishes (N2d, N4a) and the checked outputs in the report
   (274 M3);
6. finish it to a PR from the surface (243 V4).

Each step is a recorded World fact, so the milestone's evidence is the episode's own report.

## §5 Rulings for Mark

- **D-WORLD-130: adopt "daily OS v1" (§4) and the two-track order in §3.** This moves 274 M2/M3
  ahead of 272/273 in the loop track. It starts the attended track with 188 N3, then the outputs
  card, then 222 → 206/207. **A:** adopt as written. **B:** adopt the attended track only, and
  leave the loop order as filed (274 after 272). **Recommendation:** A. **Default while
  unanswered:** hold (the current order stands).
- **D-WORLD-131: re-gate row 145 (the long-run lane) from "after the v1.0.0 tag" to "after
  274 M2"**, so long builds run inside World before daily OS v1 is claimed. **A:** yes. **B:**
  keep the v1.0.0-tag gate. **Recommendation:** A. A daily OS whose builds must leave World is
  not one, and the design (D-WORLD-84) is already ruled. **Default while unanswered:** B.

## §6 Sizes (estimates, from the designs; the loop re-measures at pick time)

| Item | Basis | Estimate |
|---|---|---|
| 274 M1b-3 | the read paths for `app list|show` and `/surface/apps.json` | ~300 lines |
| 274 M2 + M3 | template + repo profile → preset → V4 start; report v1 fields | 1–2 loop sprints |
| 188 N3 | a Tauri window and a build script | about a day attended |
| N2d + N4a | the outputs doc §8 items 1–2 | 2 attended episodes |
| 222 → 206/207 | `w-all-agent-work-through-world` §3 steps 4–5 | several sprints |
| 145 | D-WORLD-84's lane: a start effect, a result transition, caps | 1–2 loop sprints |

## §7 Measurements (2026-11-28, `origin/dev`)

| # | Measurement | Result |
|---|---|---|
| P-1 | Row 188 progress | `git log --grep='row 188 N'`: N2b (#367), N2c (#374) merged; `w-world-os-surface-sprint-plan.md` §"S1 re-cut": N3 = the Tauri shell in `desktop/`; `git ls-tree origin/dev desktop` → 0 entries |
| P-2 | Rows 206/207 | charter: `[NEXT — un-parked by D-WORLD-95 R1 … after row 222]` |
| P-3 | Row 145 | charter: `[evidence row — after the v1.0.0 tag, picked by measured cost (D-WORLD-83)]`; `git tag` → `v1.0.0-rc.1` only |
| P-4 | Rows 205, 208–216 | charter: `[PARKED — post-1.x roadmap …]` |
| P-5 | Outputs N2d/N4a | `w-world-surface-outputs.md` §8 items 1–2; `git grep '"GET /surface/(file.json\|blob)' host/`: `file.json` routed, `blob` absent |
| P-6 | Row 274 | #712 merged `4f5c6bea`, #715 merged `d9c1f627`, #718 open |
