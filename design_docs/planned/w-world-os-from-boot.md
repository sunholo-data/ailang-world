# w-world-os-from-boot (roadmap notes: World as the thing you boot into) — r0

**Status:** r0, **exploration notes, not a plan.** Recorded at Mark's request, attended
2026-10-26 (*"yes please record this in roadmap potential notes"*). No ruling is asked for; it
is a direction for after 1.x. Mark then asked for the gaps to go on the roadmap (*"yep good
write them up and put in to roadmap"*). They are queue rows **204–212**, all **PARKED
post-1.x**, so the loop does not pick them up until Mark un-parks one. The mapping is in §10.

## §0 What Mark asked (attended, 2026-10-26)

These came up while trying the N2b app (#367). Typos are corrected.

1. *"Could this be built in as its own OS from boot? Like perhaps building on a Linux
   distribution etc. Do we have enough to have a user/agent drive all compute tasks through
   AILANG World? And how does this impact disk space with our traces and being able to roll
   back all the time, for system requirements vs other operating systems?"*
2. *"Have you seen Omarchy? It is a new Linux distro that is AI-led. Is that a suitable target?
   I see AILANG World as possibly replacing herdr, like I'm writing in now. Is that accurate,
   or does it run alongside, or what?"*

## §1 The boundary this keeps

[DESIGN.md §2](../DESIGN.md) says: *"AILANG World is not a replacement for Linux or macOS. It
is a semantic operating environment that runs on top of existing operating systems."* §1 names
the **OS gravity well** as a risk: everything wants to move into the kernel.

An "OS from boot" keeps that boundary:

- **Linux stays the kernel.** It keeps processes, memory, devices and drivers.
- **World is the semantic layer, and it becomes the session.** It holds goals, typed state,
  effects, grants, decisions and history.
- **Nothing here moves into World's kernel.** This is packaging and a shell, not a kernel
  change.

## §2 The shape

| layer | choice | why |
|---|---|---|
| base | a minimal **immutable** Linux, with NixOS first choice | its system config is declarative and rolls back, so it can become a World object |
| session | a Wayland kiosk compositor (e.g. `cage`) running the N3 Tauri shell full screen | you boot straight into World, with no desktop behind it |
| authority | `worldd` as a system service that owns the store, the credentials and the grants | launching programs becomes a granted, recorded effect ([outputs doc](w-world-surface-outputs.md) §6.3) |
| filesystem | Btrfs or ZFS, snapshotted per World state | rewinding big binaries costs only the changed blocks (§5) |

## §3 What World already mediates

These work today, on macOS and Linux:

- agent file I/O in worktrees;
- checks and runs;
- confined commands (srt: Seatbelt on macOS, bubblewrap on Linux);
- sessions, grants and budgets;
- refusals, and provenance for every step.

## §4 Gaps before "all compute goes through World"

Each line could become a candidate row. None is filed.

| gap | today |
|---|---|
| long-running and interactive processes | exec is capped at 9 s; row 145 (long-run lane, D-WORLD-84) is the start of this |
| GUI programs, and what they change | not mediated; outputs doc §6.3/§6.4 (open in a program, human edits as steps) |
| the human's own actions as steps | only typed `y` and decisions are recorded; file edits made by hand are not |
| network | agents never get it. A human OS needs it as a granted, recorded effect |
| installing software and changing the system | outside World; on NixOS the system config could be a World object |
| devices, GPU, audio | outside World |
| identity and approvals | the typed `y` only; row 201 (approvals by key) is designed, not built |
| retention and GC | none yet; outputs doc R5 is the first piece (§5) |
| a base that doesn't change underneath | friction #23: a Homebrew upgrade replaced node under the daemon, and every exec exited 134. An immutable base removes this class of failure |

## §5 Disk space and rollback, measured

These numbers come from the trial rig on 2026-10-26: 10 days and about 33 attended episodes.

| what | size | notes |
|---|---|---|
| the trace (`world.db`) | **69.7 MB** | 3,461 log entries; 33,257 typed objects totalling 42.9 MB (average 1.3 KB, largest 0.4 MB), so about **20 KB per call, all in** |
| exec caches | 8.1 GB | 25 episodes at 0.3–0.8 GB each, mostly per-episode Go build caches; **can be rebuilt** |
| episode worktrees | 8.1 GB | 31 full checkouts; platform-source ones are 1.1–2.4 GB with dependencies |
| pinned interpreters and the capsule cache | 0.3 GB | |
| graders' evidence | 55 MB | |

Findings:

- **History is cheap.** The trace is under 0.5% of the rig; caches and checkouts are over 95%.
  A heavy agent day of about 5,000 calls is roughly 100 MB of trace before compression, which
  comes to a few GB a year.
- **The cost is copies, and copies can be reclaimed.** Two levers are a shared,
  content-addressed build cache (instead of one per episode) and dropping merged worktrees.
- **World's rollback is not an undo.** Time Machine, APFS/Btrfs/ZFS snapshots and NixOS
  generations restore **blocks or configs**. They are cheap because copy-on-write keeps only
  changed blocks, but they can't say *why* anything changed. World **branches from any past
  state** ("change from here"), keeps every later state, and records who decided what.
- **Storage policy:**
  - facts are kept forever and are small;
  - bytes are kept by reference (git objects today; CoW snapshots plus content-addressed dedup
    on an OS build);
  - caches and merged worktrees are disposable.
- **The real design question is retention.** What must a past state still be able to *show*,
  and what must it be able to *re-run*? That is outputs-doc R5 widened from one store of
  files to a policy for the whole system.

## §6 Omarchy: a place to ship World, not the base

Facts as of 2026-10-26; sources are at the end.
[Omarchy](https://itsfoss.com/news/omarchy-ai-agent-focus/) is DHH's distro: Arch Linux with
the Hyprland window manager. Version 4 calls itself "the malleable OS for the age of agents".
It ships with Claude Code, Codex, OpenCode and others wired in, has agent skills for
troubleshooting and themes, and installs local models (LM Studio, Ollama).

- **Good as a place to ship.** Its users already run the agents World's episodes launch. World
  could be an Omarchy package with a menu entry, opening the app as a Hyprland window.
- **Poor as the from-boot base.** It is rolling and mutable, the opposite of §2's immutable
  base. Friction #23 is what a rolling base does to a daemon pinned to a toolchain.
- **Opposite philosophies that complement each other.** Omarchy gives agents the whole machine
  with the user's rights, and nothing records what they did. World is the missing layer:
  grants, a sandbox, a typed record of each step, and "change from here". Put simply,
  Omarchy wires the agents in; World makes them accountable.

## §7 herdr: alongside World today, replaced only if World grows a live agent tab

[herdr](https://www.bitdoze.com/herdr-agent-multiplexer/) is a terminal multiplexer like tmux,
built for agents. It has panes and tabs, a sidebar that marks each agent blocked, working,
done or idle, and a persistent server, so sessions survive detaching. Mark runs the attended
operator sessions in it.

| | herdr | World today |
|---|---|---|
| watch several agents | ✅ live terminal panes | ✅ timeline and now card, after the fact or near-live |
| type to an agent | ✅ the whole point | ❌ |
| agent state | read from the terminal output | typed: runs, decisions, checks |
| what the agent did, why, and undoing it | ❌ scrollback only | ✅ the whole point |
| limits on what an agent can do | ❌ | ✅ grants, sandbox, budgets |

**Today they run side by side.** herdr is where Mark talks to agents; World is where their work
is limited, recorded and reviewed. An attended operator session (a Claude Code pane in herdr)
is itself outside World. Only the episodes it launches through `trial.py` are World episodes,
and the operator's merges and pushes are not World state.

**What World would need to replace herdr** (candidate rows, none filed):

1. **A live agent tab:** see the agent's terminal and type into it. This builds on outputs-doc
   §5 (the agent's view as an output) and adds input.
2. **Real agent states:** "blocked" comes from World's own decisions and approvals (row 201),
   not from reading the terminal.
3. **Operator sessions as episodes:** the session that drives episodes is itself one. This is
   row 196 territory, gated on delegated minting.

**What would say it's time:** Mark opens World to *talk* to agents, not only to review them.

## §8 A first spike, when wanted

A NixOS VM image that boots into the World app (cage + the N3 shell) over a local `worldd`,
with the store on Btrfs. It would measure four things:

- boot to first screen;
- the disk taken per episode, with a shared build cache versus one per episode;
- snapshot cost per state on a 2 GB binary file edited in place;
- what breaks when the base updates (the friction #23 question, on an immutable base).

It needs N3 (Tauri) first, and Mark's go-ahead. It is not on the 1.x path.

## §9 Open questions

1. Should the immutable base be NixOS, or Fedora Silverblue / another image-based system?
   That depends on whether the system config should be a World object (NixOS) or an image
   World only points at.
2. Where does the human's own work become steps: an editor plugin, a file watcher, or only
   saves made through World?
3. Network as a human effect: is each request recorded, each grant, or only each session?
4. Retention defaults: how long does a past state stay re-runnable, as opposed to only
   viewable?
5. Shipping on Omarchy: does World want to depend on their packaging, or only offer an
   installer?

## §10 Where each gap lives on the roadmap

| gap | row | status |
|---|---|---|
| long-running and interactive exec | 145 | design ruled (D-WORLD-84); an evidence row after v1.0.0 |
| agent-made UI (part of the agent's view) | 151 | v1.0.0 release row |
| the full-screen app (Tauri shell) | 188 N3 | in sprint, after N2c |
| approvals by key | 201 | NEXT after row 188's style pass |
| previews, program profiles, open in the program, human edits as steps | **213–216** ([outputs doc](w-world-surface-outputs.md) §6) | PARKED post-1.x (D-WORLD-91 R4) |
| kept output bytes | outputs doc R5 → row 204 | not kept until row 204 decides (D-WORLD-91) |
| operator sessions as episodes | 196 | GATED (delegated minting, fleet driver) |
| forge and package-docs network | 146, 162 | PARKED post-1.0 |
| retention and GC | **204** | PARKED post-1.x |
| shared build cache | **205** | PARKED post-1.x |
| live agent tab (see and type) | **206** | PARKED post-1.x |
| agent states from World's records | **207** | PARKED post-1.x |
| network as a human effect | **208** | PARKED post-1.x (design first) |
| system config as a World object; toolchain check | **209** | PARKED post-1.x |
| devices, GPU, audio | **210** | PARKED post-1.x (design first) |
| the NixOS boot image spike (§8) | **211** | PARKED post-1.x, needs 188 N3 |
| an Omarchy package | **212** | PARKED post-1.x, needs 188 N3 and 209 |

## Sources (fetched 2026-10-26)

- [Omarchy bets its future on AI agents (It's FOSS)](https://itsfoss.com/news/omarchy-ai-agent-focus/)
- [Omarchy agentic AI (tbreak)](https://tbreak.com/omarchy-linux-distro-agentic-ai/)
- [What is Omarchy Linux (MindStudio)](https://www.mindstudio.ai/blog/what-is-omarchy-linux/)
- [Herdr review (bitdoze)](https://www.bitdoze.com/herdr-agent-multiplexer/)
- [Herdr guide (Better Stack)](https://betterstack.com/community/guides/ai/herdr-ai-agent/)
- [herdr (Tecmint)](https://tecmint.com/herdr-run-ai-coding-agents-in-linux-terminal)
