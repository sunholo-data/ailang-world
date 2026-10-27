# w-all-agent-work-through-world: the roadmap to running every agent task through World

**Status:** goal and roadmap, **set by Mark, attended 2026-10-26**, recorded as **D-WORLD-94**
(the goal and rows 217–222) and **D-WORLD-95** (OPEN: un-parking rows 206 and 207, and where
the sequence sits beside row 188). **r1 (attended 2026-10-27):** adds §5, the agents you
already use, and rows 223–225, with **D-WORLD-96** (OPEN: the standing allowance's bounds). This is a roadmap, not a build design: each row gets its own
design at pick time, and nothing here waives a gate.

## §0 What Mark said

> *"I can view these runs in AILANG World? Are we close to where I just put all agent work through
> AILANG World now? It seems difficult to mint and use still though."* … *"Put it on the roadmap
> what we need. I want all my agent work to be on AILANG World."* (Mark, attended 2026-10-26,
> while Codex ran episode w188l; typos corrected.)

**The goal:** every agent task Mark runs goes through World, whether headless or interactive, on
any of his repos. World is the normal way, not the trial rig, and nothing else is used alongside
it. Today World runs well-specified headless tasks (w188a–w188l; measured Codex overhead −5.7%,
`91ada43`). Everything else still happens outside it.

## §1 What it took to run w188l (measured 2026-10-26)

1. **Minting:** one command of about 1,000 characters (repo, branch, profile, model, TTL, the
   store, workspace and sessions paths, the daemon address and 10 `--grant` flags), typed at a
   terminal, then `y`.
2. **Launching:** a second step. `episode start` prints the agent's command, and the operator
   runs it (here through the rig's `scripts/trial/trial.py`).
3. **Watching:** the rig daemon was the 20 Oct build (`d7c647d`), so `/surface` was a 404. The
   live log was only at `/workbench` and `provenance --episode`. A running episode isn't a state
   on the timeline until it merges and someone rebuilds the surface index by hand.
4. **Checks:** the sandbox has no loopback, no Chrome and no `cargo`, so the task's Chrome test
   can only run at grading, by the operator.
5. **Finishing:** commit, PR body, push and PR are all operator steps.
6. **The daemon:** started by hand from `daemon-cmdline.txt`, it pins `node` at startup. One
   `brew upgrade` broke every exec call silently (friction #23). Replacing its binary needs the
   `cp … .new && mv` dance, because macOS kills a binary copied over a running one.
7. **This conversation:** an attended session like this one doesn't go through World at all.

## §2 The roadmap

| # | gap (§1) | row | the item, briefly | exit criterion |
|---|---|---|---|---|
| 1 | 6, 3 | **219 w-worldd-home-service** | World as an always-on user service: a home store (not the trial rig), a user LaunchAgent or systemd unit, `ailang-worldd upgrade` that swaps the binary safely and refuses while sessions run, exec tools resolved per call, and a stale-build warning in the now card. A migration from the trial rig's store. | the daemon survives a reboot and a `brew upgrade` with no manual step, and an upgrade lands with one command |
| 2 | 1, 2 | **217 w-episode-one-command** | launch presets: repo, profile, grants, TTL, agent and model in one named file per repo, so a mint is `ailang-worldd episode start --preset world-codex --task t.md` and then `y`. Episode id and branch are derived when omitted. `--launch` starts the agent (Claude or Codex) from an empty directory and records its transcript. | Mark starts a Codex episode and a Claude episode on two different repos with one short command each |
| 3 | 5 | **218 w-episode-finish** | `episode finish <ep>`: runs the preset's declared checks, commits the worktree, pushes and opens the PR with the episode report as its body (an outward effect that Mark approves), records the PR on the episode and closes the session. | an episode goes from start to merged PR without `trial.py` or a hand-written PR |
| 4 | 3 | **220 w-surface-index-follow** | the daemon keeps the surface index current: it rebuilds on new merges in the configured repos, shows running episodes as open states, and has one surface for several repos. | a merged PR appears in the app within a minute with no manual command, and a running episode appears as an open state |
| 5 | 4 | **221 w-exec-profile-dev-tools** (with row 199's loopback) | profile-declared, confined dev tools the agent can run itself: headless Chrome for browser checks, `cargo`, and node test runners. | an episode like w188l where the agent runs its own Chrome test and the operator's re-run only confirms it |
| 6 | 1 | **201 w-world-approvals-by-key** (existing, NEXT) | approve a mint or an outward effect with a passkey or SSH key, from the app or a phone, instead of typing at a terminal. | its own |
| 7 | 7 | **222 w-attended-session-as-episode** | an interactive Claude Code or Codex session becomes an episode. The plugin's `/world-start` asks Mark to approve a mint (a typed `!` command, or row 201's key), the session's own file and shell tools are switched off in favour of World's, and `episode finish` closes it. | an attended session like this one runs end to end through World |
| 8 | 3, 7 | **206 w-world-live-agent-tab** and **207 w-world-agent-states** (existing, PARKED) | watch and talk to a running agent from the app; know at a glance which agent is working, blocked or done. | their own |
| 2b | 1, 2 | **223 w-agent-host-adapters** (r1; design [r0](w-agent-host-adapters.md)) | the agents Mark already runs get into World with no flags: a Claude Code profile installed once machine-wide (the plugin, World's MCP server, its own file and shell tools denied, a session-start hook that attaches a session); a Codex `world` profile; a **herdr plugin** that starts each agent pane as an episode and labels the pane with it; Claude Desktop and IDE agents through World's MCP server. The app shows each agent's level (§5). | launching Claude Code, Codex or a herdr agent pane needs no World flags, and each shows in the app as **confined**; Claude Desktop shows as **observed** |
| 6b | 1 | **224 w-standing-session-allowances** (r1) | Mark approves a preset once with his key: bounded sessions a day, fixed grants, TTL, repos, an expiry. Host adapters draw sessions from it without a prompt; each draw is recorded against the approval; it is revocable at any time. Bounds are D-WORLD-96. | a week of launches with no per-session prompt, every session traceable to its allowance and the approval |
| 6c | 4, 7 | **225 w-world-run-wrapper** (r1) | `ailang-worldd run --preset <p> -- <agent> …`: any CLI agent with no adapter runs under World's sandbox, with the episode worktree as its only write root and no network but the daemon and the model API. Its writes land as worktree changes: effect-level provenance, not typed calls. | an agent with no adapter (pi) runs an episode through the wrapper, and its writes outside the worktree are refused |
| 9 | (the loop) | **196 w-loop-executors-through-world** (existing, GATED) | the mission loop's executors work through a pinned World. It is still gated on Mark's delegated-minting ruling and a fleet driver change. | its own |

**What doesn't change:**
- **A mint stays Mark's act** (typed `y`, or row 201's key). The roadmap makes it short, not
  automatic. Delegated minting stays row 196's separate ruling.
- **Every outward effect needs Mark's approval:** push, PR, publish.
- **The frozen core is untouched.** Row 219's service is World's own user agent;
  `tools/launchd/*` (the fleet driver) is never touched.

## §3 Order (recommended)

1. **Now:** finish row 188's N2c (w188l, running). Then swap the trial rig's daemon to a current
   build with `--surface-index`, so the app shows the next episodes. That swap is an interim
   step until row 219.
2. **219 → 217 → 223 → 218.** These change how every later episode runs, so they pay back
   first. Each is an attended episode of its own, run through the rig as today, and from 217
   on with presets. 223 is where the herdr plugin lands.
3. **Back to row 188** (N2d, N4a, N3, N4, N5), now run through 217 and 218, with **220** beside
   N2d, since both are data work for the app.
4. **221 with 199**, then **201** (already NEXT after row 188's style pass), then **224**
   (the allowance, which needs 201's keys), **225**, then **222**.
5. **206 and 207,** if D-WORLD-95 un-parks them. Then **196**, when its ruling and the fleet
   change exist.

**Done when** a week of Mark's agent work, across his repos and both headless and interactive,
shows in World with no task run outside it. Each task that is run outside World is a friction
entry that names the missing piece.

## §5 The agents Mark already uses (r1, attended 2026-10-27)

> *"When using AILANG World as an OS, will it be possible to launch IDEs, Claude Desktop, OpenAI
> Codex, herdr etc and have those within AILANG World when they run agents? Either
> automatically, or with a quick setting? I'm looking for user friction."* … *"Let's start some
> work strands to do this, yes. A herdr plugin is a good idea, but yes, the full experience in
> AILANG World OS gives us full control, yes?"* (Mark, attended 2026-10-27, typos corrected.)

**Two levels, shown per agent in the app:**
- **Confined:** the agent's only way to touch files, the shell and the network is World's
  tools, under a session Mark authorised. World is the boundary, and provenance is typed calls.
- **Observed:** the host keeps its own tools; World sees only what goes through it. It must
  never be mistaken for confined.

| host | on macOS today (row 223) | in World OS (rows 209–211) |
|---|---|---|
| Claude Code (CLI, IDE extensions, herdr panes) | confined, through a machine-wide profile | confined |
| Codex CLI | confined, through a `world` profile (w188l's exact flags) | confined |
| herdr | inherits its panes' agents; the plugin makes each pane an episode | confined |
| any CLI agent without an adapter | confined at the effect level through row 225's wrapper | confined |
| Claude Desktop, Cursor and other IDE agents | observed, through World's MCP server | **confined at the effect level** |

**World OS gives full control of every effect on the machine:** the filesystem, processes and
network are World's, so any program runs confined, and its writes land in an episode's
worktree as recorded changes. **Two limits stay:**
1. **Provenance is coarser for programs that don't speak World.** World sees their file writes,
   processes and network flows, not their intent. Typed, explainable calls (`why`) come only
   from agents that use World's tools.
2. **Off-machine actions aren't World's.** What a cloud service does on its own servers, such as
   a hosted connector editing a Drive file, happens outside the machine. World controls the
   request leaving the machine (network as a granted effect, row 208), not the action at the
   other end.

**Friction, the target:** you open Claude Code, Codex or a herdr pane as today. The first launch
of the day asks for one tap (row 201), or nothing under a standing allowance (row 224). Every
run shows in the World app with its level.

## §4 Rulings

- **D-WORLD-94 (resolved, attended 2026-10-26):** the goal above, and rows 217–222 filed NEXT
  in §3's order.
- **D-WORLD-95 (OPEN):**
  - **R1:** un-park rows 206 and 207 into this sequence after 222. They are PARKED post-1.x
    today, and the goal needs them. Default: yes.
  - **R2:** run 219, 217 and 218 before row 188's N2d. Default: yes; the alternative is after
    N4a.

- **D-WORLD-96 (OPEN, r1): the standing allowance's bounds (row 224).** Default:
  - per preset, at most 10 sessions a day, with the preset's grants and a 12 h TTL;
  - the repos listed in the preset;
  - the approval lasts 30 days, revocable at any time from the app or the CLI;
  - going over the bounds asks for a tap;
  - the mission loop (row 196) is excluded unless ruled separately.

## Sources

- Episode w188l's mint and launch (attended 2026-10-26); the rig's `daemon-cmdline.txt`.
- [`w-world-os-from-boot`](w-world-os-from-boot.md) §7 (herdr) and §10 (rows 204–212).
- Friction #13 (stale rig build) and #23 (exec node pinned at startup) in
  [`friction.md`](../verification/world-dogfood/friction.md).
- Rows 189 (episode start, profile init, the plugin), 196, 198 (episode report), 199, 201, 206
  and 207 in the charter.
