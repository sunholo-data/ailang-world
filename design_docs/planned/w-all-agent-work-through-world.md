# w-all-agent-work-through-world: the roadmap to running every agent task through World

**Status:** goal and roadmap, **set by Mark, attended 2026-10-26**, recorded as **D-WORLD-94**
(the goal and rows 217–222) and **D-WORLD-95** (OPEN: un-parking rows 206 and 207, and where
the sequence sits beside row 188). This is a roadmap, not a build design: each row gets its own
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
2. **219 → 217 → 218.** These change how every later episode runs, so they pay back first.
   Each is an attended episode of its own, run through the rig as today, and from 217 on with
   presets.
3. **Back to row 188** (N2d, N4a, N3, N4, N5), now run through 217 and 218, with **220** beside
   N2d, since both are data work for the app.
4. **221 with 199**, then **201** (already NEXT after row 188's style pass), then **222**.
5. **206 and 207,** if D-WORLD-95 un-parks them. Then **196**, when its ruling and the fleet
   change exist.

**Done when** a week of Mark's agent work, across his repos and both headless and interactive,
shows in World with no task run outside it. Each task that is run outside World is a friction
entry that names the missing piece.

## §4 Rulings

- **D-WORLD-94 (resolved, attended 2026-10-26):** the goal above, and rows 217–222 filed NEXT
  in §3's order.
- **D-WORLD-95 (OPEN):**
  - **R1:** un-park rows 206 and 207 into this sequence after 222. They are PARKED post-1.x
    today, and the goal needs them. Default: yes.
  - **R2:** run 219, 217 and 218 before row 188's N2d. Default: yes; the alternative is after
    N4a.

## Sources

- Episode w188l's mint and launch (attended 2026-10-26); the rig's `daemon-cmdline.txt`.
- [`w-world-os-from-boot`](w-world-os-from-boot.md) §7 (herdr) and §10 (rows 204–212).
- Friction #13 (stale rig build) and #23 (exec node pinned at startup) in
  [`friction.md`](../verification/world-dogfood/friction.md).
- Rows 189 (episode start, profile init, the plugin), 196, 198 (episode report), 199, 201, 206
  and 207 in the charter.
