# w-worldd-home-service (queue row 219): World as an always-on user service — r0

**Status:** r0, written attended 2026-10-27 through World itself (episode `h2610271506` on the
trial rig; Mark: *"design row 219 (w-worldd-home-service). No code yet."*). It is a design, not a
build: no code, no plist and no service exist yet. It is not quorum-reviewed. The pick-time
designer re-measures every premise marked **verify** in §1 on the Mac and runs the quorum as
usual. Roadmap: [`w-all-agent-work-through-world`](w-all-agent-work-through-world.md) §2 row 1,
§3 step 2 (`223a + 227 → 223b → 226, then 219 → 217 → 218`). Rulings it needs: §8 R1–R4.

## §0 What is wrong today

Row 219's own words: **"THE DAEMON IS A HAND-STARTED PROCESS ON A TRIAL RIG."** Measured over
w188k, w188l and this week's attended work:

- **It dies with whatever started it.** The rig daemon ran as a child of the operator's Claude
  session; a host crash and reboot on 26 Oct ended both and nothing restarted it. The next launch
  failed as `pre-flight: /v1/health -> HTTP 0`, naming neither cause nor fix (friction **24**).
- **It is stale and nothing says so.** During w188l it was a 6-day-old build without `/surface`
  (roadmap §1.3); before that, a pre-row-139 build serving the old tool pin (friction 13).
- **It pins exec tools at startup.** `--exec-node /opt/homebrew/bin/node` is resolved once to
  `Cellar/node/26.0.0`; a `brew upgrade` removed that, and every `workspace-exec` died in dyld
  with exit 134 while the daemon kept granting exec (friction 23).
- **Replacing its binary is a ritual.** `cp … .new && mv`, because macOS kills a binary copied over
  a running one (roadmap §1.6), and the restart is by hand from `daemon-cmdline.txt`.
- **The CLI and the daemon drift apart.** Mark, attended 2026-10-27: a CLI-only swap made
  `episode start` refuse, because it requires the daemon to be *exactly its own build*
  (`checkDaemon`, §1 P3). Swapping one without the other is always a refusal.
- **Generated scripts call a path that moves.** After `hosts init --replace`, `world-context`
  reported "daemon down" to a live session: the scripts called a bare `ailang-worldd`, not on the
  agent's PATH, and exit 127 fell into the daemon-down branch. The fix now writes the absolute
  path of the binary that ran `hosts init` (friction **32**), which is only as stable as wherever
  that binary was built.

**Friction 33** (added with this doc's PR) is the CLI-only swap of §0's fifth bullet; the
pick-time designer folds it in.

**Exit criterion (row 219):** the daemon survives a reboot and a `brew upgrade` with no manual
step, and an upgrade lands with one command.

## §1 Premise verification log

Each premise is **measured** (read in this worktree's code or docs on 2026-10-27), **quoted**
(from a charter row or friction entry, not re-measured), or **verify** (a macOS fact this
session could not measure: it has no shell, only World's tools).

| # | premise | source | status |
|---|---|---|---|
| P1 | `/v1/health`'s `build` is computed **lazily**: `CurrentBuild()` is a `sync.OnceValue` first called inside `handleHealth`, and it hashes the file at `os.Executable()`'s *path*. A binary renamed over before the first health request reports the **new** file's sha256 while the **old** code runs. | `host/daemon/build.go` (`currentBuild`, `executableDigest`); `host/daemon/daemon.go` (`handleHealth`) | measured |
| P2 | `SameBuild` matches on executable sha256 first, then on a clean vcs revision; a build in a git worktree carries no revision, so only identical bytes match. | `host/daemon/build.go` | measured |
| P3 | `episode start` refuses a daemon whose build, tool pin, credentials store or profile digest differs from its own, and never starts or restarts one; with none serving it prints a `serve` line. | `cmd/ailang-worldd/episode.go` (`checkDaemon`, `serveLine`) | measured |
| P4 | Two default addresses: `serve --bind` and every client verb's `--addr` default to `127.0.0.1:7644` (`daemon.DefaultBind`/`DefaultAddr`); `episode start --addr` and `hosts init --addr` default to `127.0.0.1:47330`, the rig's port; the herdr plugin's install says `--addr http://127.0.0.1:47330`. | `host/daemon/daemon.go`; `cmd/ailang-worldd/episode.go`; `cmd/ailang-worldd/hosts.go`; `integrations/herdr/ailang-world/README.md` | measured |
| P5 | The home layout already exists as defaults (D-CLI-4): store `~/.ailang/world/world.db`, credentials `credentials.db` beside it, workspace root `~/.ailang/world-ws`, sessions `~/.ailang/world/sessions`, pins `~/.pinned-ailang` and `~/.pinned-ailang-tools`, exec profiles `~/.ailang-world-exec/profiles`, hosts `~/.ailang-world/hosts`. | `cmd/ailang-worldd/setup.go` (`defaultWorldPaths`), `episode.go`, `hosts.go` | measured |
| P6 | `hostsWorlddPath` writes `os.Executable()` **with symlinks resolved** into the generated scripts. | `cmd/ailang-worldd/hosts_awareness.go` | measured |
| P7 | SIGTERM starts a bounded drain: `http.Server.Shutdown` under 10 s, then a hard close and a non-zero exit. The startup budget is 9 s. | `host/daemon/daemon.go` (`shutdownTimeout`, `drain`, `Run`) | measured |
| P8 | `world-publish transitions` needs the **world** writer lock, so the daemon must be stopped for it; `session new/mint/revoke` use the credentials store's own lock and work while it serves. | `docs/QUICKSTART.md` §0; `doctor.go` (`doctorStore`) | measured |
| P9 | The World-Provenance trailer is `store=<db path from /v1/health> episode=<ep> entries=<from>-<to>`: **a path**, and entry indexes are per store. Every merged `ailang-world` PR so far names the rig's db path. | `cmd/ailang-worldd/provenance.go` | measured |
| P10 | `/v1/health` admits additive fields only. | `host/daemon/daemon.go` (`HealthResponse` comment) | measured |
| P11 | `doctor --quick` promises that its only subprocess is the host's `--version`; doctor's classes are daemon 10, session 11, host 12, exec 13. `whoami` exits 10 for daemon unavailable. | `cmd/ailang-worldd/doctor.go`; `docs/QUICKSTART.md` | measured |
| P12 | The rig: `world.db` 69.7 MB for 3,461 entries; `store/exec-cache` 8.1 GB; `ws/` 8.1 GB (31 worktrees), measured 2026-10-26. | charter row 204 | quoted |
| P13 | A LaunchAgent in `~/Library/LaunchAgents` loads into `gui/<uid>` at login; `launchctl bootstrap`/`bootout`/`kickstart -k`/`print` drive it; `KeepAlive` restarts it, rate-limited by `ThrottleInterval` (default 10 s); `ExitTimeOut` (default 20 s) is the SIGTERM→SIGKILL gap; a job's environment is launchd's minimal one (`PATH=/usr/bin:/bin:/usr/sbin:/sbin`) plus `EnvironmentVariables`; launchd refuses a group- or world-writable plist; `XPC_SERVICE_NAME` carries the job's label. | launchd.plist(5), launchctl(1) | **verify** |
| P14 | `rename(2)` over a running binary leaves the running process on the old inode, which is why `cp .new && mv` works and `cp` over the file does not. | roadmap §1.6 (operator-measured), POSIX | quoted / **verify** the kill-on-overwrite half |
| P15 | With FileVault on, no user process runs before the user logs in after a reboot, so "survives a reboot" means *serving shortly after login*. | macOS | **verify** |

**P1 is a latent defect that this design would make live.** Today the `mv` ritual has the same
exposure, but only once per hand restart. A service whose upgrade renames the binary would hit
it every time a daemon goes unasked between its start and the swap. M0 fixes it first.

## §2 Conflict surface

- **The fleet's launchd jobs.** `tools/launchd/*` is the fleet's frozen driver and plist
  (guardrail, charter line 1500; `D-WORLD-DRIVER-1`), and `dev.ailang.mission-world` is the loop's
  job. This design writes **no file under `tools/launchd/`**, uses a label with no
  `dev.ailang.` substring (§3.1), and has no code path that reads, loads, unloads or edits any
  other job. A test pins the write set (AC11).
- **Row 156 (protected test paths, in-sprint)** works in `host/broker`'s exec path, where M3 also
  works. **M3 waits for row 156 to land**, or the two sequence through one controller.
- **Row 223 (built: `doctor --session/--host/--quick`, `hosts init`, the herdr plugin).** 219 changes
  their *defaults and pointers* only (§5). It keeps doctor's classes and `--quick`'s
  no-extra-subprocess promise (P11).
- **Row 227 (error clarity).** 219 owns the daemon-down/stale-build/service messages and uses
  227's three-part form (§6). 227 owns the rest.
- **Row 217 (presets: "store and daemon" per preset).** A preset's daemon defaults to the home
  service, and `service.json` (§3.2) is where 217 reads that default.
- **Row 209 (PARKED: a pinned-toolchain manifest checked at start).** M3 does the minimum the
  exit criterion needs (resolve per call, name the cause), and 209 stays the full form.
- **Row 204 (PARKED: retention and GC).** 219 deletes nothing, rig included (§7). The rig's
  16 GB is 204's to reclaim.
- **Row 196 (GATED: the loop through World).** The loop must not start depending on this
  service by accident. 196 decides that, and 219 does not.
- **Row 220 (surface index follow).** The plist carries `--surface-index` (§3.1). 220 makes the
  index current.
- **Rows 211/212 (PARKED: World OS, Omarchy).** Both name a `worldd` system or user service. §3.5's
  systemd unit is the user-service half of that, and the system service stays theirs.

## §3 The design

### §3.1 One LaunchAgent, World's own

**The label is `com.sunholo.ailang-world.worldd`, and the plist is
`~/Library/LaunchAgents/com.sunholo.ailang-world.worldd.plist`** (mode 0644, owned by the user).
The label shares no prefix with the fleet's `dev.ailang.*` jobs, so a fleet script that greps
`dev.ailang` never sees it. `service install` refuses any label starting with `dev.ailang`.

The plist is rendered from `service install`'s validated flags. Every path in it is absolute,
because launchd does not expand `~`:

| key | value | why |
|---|---|---|
| `ProgramArguments` | `[<home>/.ailang-world/bin/ailang-worldd, serve, --db …, --credentials-db …, --bind 127.0.0.1:7644, --ailang-bin <home>/.pinned-ailang/ailang, --workspace-root <home>/.ailang/world-ws, --tool-ailang-bin <the pinned tool path>, --exec-sandbox <srt node_modules>, --exec-profiles-dir <home>/.ailang-world-exec/profiles, --exec-node /opt/homebrew/bin/node, --examples-dir …, --surface-index …]` | one stable binary path (§3.3); the D-CLI-4 defaults (P5); the port is R1. `--exec-node` names the **symlink**, which M3 follows per call |
| `RunAtLoad` | `true` | starts at login |
| `KeepAlive` | `true` | restarts after a crash, a kill or a startup refusal |
| `ThrottleInterval` | `30` | a startup refusal (port busy, lock held) restarts at most twice a minute, not six times. `service status` names the loop (§6 F3) |
| `ExitTimeOut` | `30` | longer than the 10 s drain plus close (P7), so launchd never SIGKILLs a draining daemon |
| `ProcessType` | `Standard` | it serves interactive agents, so not `Background` |
| `WorkingDirectory` | `<home>/.ailang/world` | never a repo or a worktree |
| `Umask` | `63` (`077`) | every file it creates is private |
| `EnvironmentVariables` | `PATH=/usr/bin:/bin:/usr/sbin:/sbin` only | nothing ambient, and no `AILANG_REGISTRY_API_KEY` (§6 F6) |
| `StandardOutPath` / `StandardErrorPath` | `<home>/Library/Logs/ailang-world/worldd.log` (both) | the listen line, integrity lines and every sanitized-500 cause in one file, in order |

**Loopback only, by construction.** `serve` already refuses a non-loopback `--bind` with no
override (`isLoopbackHost`). `service install` validates the same predicate before rendering, so
it can never write a plist that would crash-loop on it.

**Logs.** launchd opens the log file at each spawn and never rotates it. `service start` and
`service upgrade` rotate it before the restart: over 10 MiB it is renamed to `worldd.log.1`,
and one generation is kept. `service status` prints its size and warns over 100 MiB.
In-flight rotation is §9 Q3.

### §3.2 The home: three directories that already exist, plus one record

- **Store** (the D-CLI-4 defaults, P5): `~/.ailang/world/world.db` and
  `~/.ailang/world/credentials.db`; workspace root `~/.ailang/world-ws`; sessions
  `~/.ailang/world/sessions`. They are reused, not invented, so `doctor`, `setup`, `session new`
  and `episode start` agree with the service with no flags.
- **Binary:** `~/.ailang-world/bin/ailang-worldd` (§3.3).
- **Record:** `~/.ailang-world/service.json`, mode 0600, schema `world/service/v1`. It holds the
  label, the plist path and its sha256, the binary path and its sha256, the previous binary's
  sha256, the address, the store, credentials and workspace paths, the log path, `installed_at`,
  and the last 10 upgrades (from, to, time, outcome). It is the one file `doctor --quick`,
  `whoami`, `hosts init`, the herdr plugin and row 217's presets read to find the service, with
  no subprocess. It is operator-local state, not a World object (§9 Q2).

The binary must lie **outside the workspace root and outside the store directory**: the
workspace root is an agent's only write root, and the exec sandbox writes episode caches under
the store's directory (P12's `store/exec-cache`). `~/.ailang-world/bin` is neither.
`service install` refuses a layout where it would be.

### §3.3 One stable binary path, a real file

`~/.ailang-world/bin/ailang-worldd` is the **only** `ailang-worldd` the service, the herdr plugin,
the generated host scripts and the operator's shell use. `service install` prints the one line to
put it first on `PATH`. It edits no shell rc file.

- **It is a regular file, never a symlink.** P6: `hosts init` resolves symlinks before writing the
  path into scripts, so a `bin/ailang-worldd → versions/<sha>/ailang-worldd` link would bake the
  *versioned* path into every script, and the next upgrade would strand them. A real file at one
  path keeps every script's bytes identical across upgrades (AC14).
- **It is replaced only by `rename(2)`** from a sibling temp file on the same filesystem (P14),
  never by writing into it.
- **One file means the CLI and the daemon cannot drift.** `episode start` run from it matches the
  daemon by bytes (P2), so the CLI-only swap Mark hit becomes impossible through the supported
  path. A binary run from anywhere else that differs gets §6 F14's message.

### §3.4 The verbs

All of them are `ailang-worldd service <verb>`. `install`, `upgrade`, `stop` and `uninstall` are
**attended**: they open `/dev/tty` with session new's fence (`refusing: no controlling terminal`
without one), show exactly what they will change, and act only on a typed `y`. An agent can never
run them: `workspace-exec` runs only profile commands, and it has no TTY. `status` and `start` are
not fenced. `status` is read-only, and `start` loads a plist that was already approved.

Every verb validates everything before its first side effect, and holds
`~/.ailang-world/service.lock` (flock) while it changes anything. `doctor` and `episode start`
report a held lock as "a service change is in progress (pid N)", never as "daemon down".

**`service install [--bind] [--db] [--credentials-db] [--workspace-root] [--interpreter-dir]
[--tools-dir] [--exec-sandbox] [--exec-profiles-dir] [--exec-node] [--examples-dir]
[--surface-index] [--binary <file>] [--replace]`**
1. Validate:
   - macOS (§3.5 for Linux);
   - both pins hash to the compiled-in digests (setup's check, never executed);
   - the store directory and workspace root exist and the root does not contain the store;
   - the binary location is outside both (§3.2);
   - the port is free;
   - the store's writer lock is free;
   - `launchctl getenv AILANG_REGISTRY_API_KEY` is empty;
   - no job with the label is loaded;
   - no plist exists at the path (unless `--replace`, which shows a diff).

   `--binary` defaults to the running executable.
2. Show the full plist text, `<binary> → ~/.ailang-world/bin/ailang-worldd (sha256 …)`, the
   `service.json` it will write and the log path. Warn, without refusing, when the store has no
   transition registry or no world head (doctor's two `store` warnings), and print the publish and
   genesis steps (§7 step 2).
3. On `y`, in order:
   - copy the binary (temp file plus rename);
   - write the plist (temp plus rename, 0644);
   - write `service.json`;
   - `launchctl bootstrap gui/<uid> <plist>`;
   - poll `/v1/health` for up to 30 s until `build.executable_sha256` equals the installed file's
     sha256.
4. If health never matches, `bootout`, keep the files, print the last 20 log lines and exit 2
   (§6). Nothing is left loaded and half-working.

**`service status [--json]`** — read-only. Its subprocess is `launchctl print gui/<uid>/<label>`
(a 2 s timeout). It prints:
- installed (the plist exists, and its sha256 matches `service.json`, or "edited by hand");
- loaded, the pid, uptime, the number of runs and the last exit status;
- the address, and `/v1/health`'s answer;
- the **running** build against the **installed** file's sha256 ("restart pending" when they
  differ);
- the store, the log path, its size, and on failure the log's last 10 lines.

It exits 0 when serving and matching, 10 when not serving (doctor's daemon class), and 1 for a
usage error.

**`service upgrade [--binary <file>] [--pause]`** — the swap and the restart are one act:
1. Validate the candidate (default: the running executable):
   - it is a regular file, not the installed file itself ("already installed", exit 0, when its
     sha256 matches);
   - it runs `<candidate> version --json` (M0) under a 5 s timeout with only `PATH` and `HOME`;
   - that answer names this platform, the **same tool pin** as the installed setup
     (D-WORLD-57; otherwise it refuses with "run `ailang-worldd setup` with the new binary first"),
     and the **same interpreter release** the store's epoch registry admits (otherwise the epoch
     change is attended: row 139's `world-publish epoch nominate`);
   - the service is installed and healthy, or the refusal says to install it first.
2. **Refuse while sessions run** (R2): when `/v1/health`'s `active` (M0) shows a call in flight,
   or any unexpired session committed an entry in the last 5 minutes, it names each episode and
   its last call and exits 1, changing nothing. `--pause` overrides after the `y`, which then lists
   the episodes whose calls will fail for the restart's few seconds.
3. Show `from <rev/sha/age> → to <rev/sha/age>`, the pins (unchanged), the plist (unchanged, or
   the diff when the candidate renders different flags) and the expected downtime.
4. On `y`:
   - copy the candidate to `bin/ailang-worldd.new`;
   - hard-link the current file to `bin/ailang-worldd.prev`;
   - `rename(.new → ailang-worldd)`;
   - rotate the log;
   - `launchctl kickstart -k gui/<uid>/<label>` (SIGTERM, the 10 s drain, then a respawn from the
     new file);
   - poll health for up to 30 s for `build.executable_sha256 == sha256(candidate)`.
5. **On success**, update `service.json` and print the from/to line. **On failure**, roll back:
   rename `.prev` back, `kickstart -k`, wait for the *old* sha256, and exit 2 with the new build's
   last log lines. **If the rollback also fails**, `bootout` (stop the crash loop), keep both files,
   and exit 3 with the exact commands to recover (§6 F13).

**`service stop`** — show the active sessions (as in upgrade step 2), then on `y`, `launchctl
bootout`. The plist stays, so the next login starts it again; `service start` starts it now.
This is the step `world-publish` needs (P8; §6 F9).

**`service start`** — verify the plist's sha256 against `service.json`. A hand-edited plist shows
its diff and asks `y`. Then rotate the log, `bootstrap`, and wait for health as install does.

**`service uninstall`** — validate as stop does, then show what goes and what stays. It removes
the plist and marks `service.json` `uninstalled`. It keeps the binary, `.prev`, the store,
credentials, sessions, worktrees and logs, each named with its path. On `y`, `bootout` and
remove the plist. **It never deletes data**, and the printout ends with the paths to delete by
hand.

**`ailang-worldd version [--json]`** (M0) prints this binary's build identity, its platform, and
the tool and interpreter releases it pins. It needs no daemon.

### §3.5 Linux: a systemd user unit, later

`setup` supports linux/amd64, and row 219 says "LaunchAgent or systemd unit". The same verbs map
one to one:
- `~/.config/systemd/user/ailang-world-worldd.service` (`Restart=always`, `RestartSec=30`,
  `TimeoutStopSec=30`);
- `systemctl --user enable --now`, `restart`, `status`, `disable --now`;
- `loginctl enable-linger` is the operator's own choice, never done by World.

It ships in M6. Until then, `service install` on Linux refuses with §6 F1's message.

### §3.6 Exec tools resolved per call (row 219's item; friction 23)

`--exec-node` and each profile command's program keep the **configured** path (for example
`/opt/homebrew/bin/node`) as their identity, never its startup resolution:
- **Every call** resolves the path (`EvalSymlinks` and a `stat`, microseconds). When the target
  (path, inode, size, mtime) differs from the last one probed, the profile's existing toolchain
  probe arm runs once against the new target before the call proceeds.
- **The startup confinement rules are re-applied** to the new target: a regular file outside the
  workspace root and the state directory. A re-resolution never weakens them.
- **A target that no longer starts** refuses the call with §6 F15's message, never a bare
  exit 134.
- `/v1/health.exec_tools` reports the configured path, the current target and the last probe's
  result.
- Each exec effect record names the resolved target. Its sha256 is computed on change and cached,
  so provenance says which `node` ran.

A profile that names a Cellar path itself (for example a `go env` result) cannot follow an
upgrade. It refuses by name, and the message points at `profile init` to re-propose it. The
profile's digest is bound at mint, so World never rewrites it silently.

### §3.7 What the daemon reports about itself (M0)

These are additive `/v1/health` fields (P10):
- **`build` computed eagerly at startup** (the P1 fix);
- `started_at`;
- `supervisor`: the launchd label from `XPC_SERVICE_NAME` when it equals World's label, else
  `none`;
- `active`: calls in flight, and per unexpired session its episode and last-call time.

## §4 Stale build in the now card

The app's view model (`/surface/model.json`) gains `daemon {build, started_at, supervisor,
installed_sha256, restart_pending}`. `installed_sha256` is the stable file's hash, recomputed only
when its size or mtime changes. The now card shows one line:
- **restart pending** (running ≠ installed): *"World is running build `abc123` (6 days old); the
  installed binary is `def456`. Run `ailang-worldd service upgrade`."* A supported upgrade swaps
  both together, so this state means a hand copy or a failed upgrade, and it is always shown.
- **not supervised** (`supervisor: none`): *"World is not running as a service. It stops when
  its terminal does. Run `ailang-worldd service install`."*
- otherwise the build and its age, in small text.

"Older than `dev`" needs a repo the daemon does not have (§9 Q1).

## §5 Everything that points at World points at the service

| piece | today | after 219 |
|---|---|---|
| **addresses** (P4) | `7644` for serve and client verbs; `47330` for `episode start`, `hosts init` and the herdr plugin | one default, `daemon.DefaultAddr` (R1), read from `service.json` when present |
| **`doctor`** | daemon ✗ "nothing is listening at …; fix: start World: ailang-worldd serve …" | the full pass adds a `service` check (installed, loaded, pid, plist digest, running = installed, log size; `launchctl print` runs only here). `--quick` reads `service.json` and the plist with **no new subprocess** (P11), so the daemon-down fix becomes `ailang-worldd service status` / `service start`. A CLI that is not the installed binary gets `!` with §6 F14 |
| **`whoami`** | exit 10 "daemon unavailable" | the same exit 10, with §6 F2/F3's what, why and next |
| **`episode start`** | refuses a mismatched daemon: "rebuild … stop the daemon and start serve again" | the same check. The fix names `service upgrade` (or "use `~/.ailang-world/bin/ailang-worldd`"), and with no daemon it prints `service start`, not a `serve` line. `--addr` defaults per R1 |
| **`hosts init`** | writes the running binary's resolved path (P6) and `--addr` 47330 | when `service.json` exists, it writes the **installed** path and address, and says so when the running binary is another one. AC14: the generated files are byte-identical before and after an upgrade |
| **herdr plugin** | `world.env` hand-edited (`WORLDD`, `WORLD_ADDR`, `WORLD_DB`, `WORLD_WS`, `WORLD_SESSIONS`, `WORLD_HOSTS`) | `bin/world-config-init` reads its defaults from `service.json` (sh and jq, already required), so a new `world.env` needs no edit. An existing one is never rewritten: `service install` prints the lines to change |
| **QUICKSTART / USE_WORLD_ON_YOUR_REPO** | `serve … &` and `daemon-cmdline.txt` | a "Run World as a service" section ahead of §1 with its command lines bound by tests (coding-standards S7); hand `serve` stays documented for scratch stores |

## §6 Failure modes

Every message has three parts in the launchers' existing form:

```
ailang-worldd service <verb>: <what failed>
  why: <the cause, measured>
  next: <a command, or the exact action>
```

Exit codes:
- **0** ok;
- **1** refused before any change;
- **2** failed after a change began, rolled back cleanly;
- **3** failed and **not** fully rolled back, with the state named;
- **10** for `status` when not serving.

| # | case | what / why / next |
|---|---|---|
| F1 | not macOS | `service install: not supported on linux yet` / the systemd unit is M6 / `run serve by hand (QUICKSTART §1)` |
| F2 | not installed | `no World service is installed` / no `service.json` and no plist at `<path>` / `ailang-worldd service install` |
| F3 | installed, not serving | `the World service is not serving at 127.0.0.1:7644` / `launchctl print`: not loaded (stopped), **or** loaded and restarting every 30 s with last exit N and the log's last line, quoted / `ailang-worldd service start`, or the fix for the quoted refusal |
| F4 | port taken | `127.0.0.1:7644 is in use` / another process listens there (the rig, or another serve) / `lsof -nP -iTCP:7644 -sTCP:LISTEN`, stop it, or `--bind` another port |
| F5 | store writer lock held | `another process holds the writer lock on <db>` / `serve` would fail at store-open / stop that daemon (`doctor` names a running worldd) |
| F6 | registry key in launchd's env | `launchd's environment carries AILANG_REGISTRY_API_KEY` / every World process would inherit publish authority, and serve refuses it / `launchctl unsetenv AILANG_REGISTRY_API_KEY` |
| F7 | pins missing or drifted | `the tool pin is missing at <path>` / serve archives both pins at startup / `ailang-worldd setup` |
| F8 | nothing published | (warning) `the home store has no transition registry` / no tool will serve / §7 step 2 |
| F9 | publish while the service runs | (`world-publish`'s refusal) `the world writer lock is held by the World service` / publish needs the world lock (P8) / `ailang-worldd service stop`, publish, `ailang-worldd service start` |
| F10 | bad upgrade candidate | `<file> is not an ailang-worldd for darwin/arm64` or `it pins tool AILANG v… (installed v…)` or `it serves interpreter v…, which the store's epoch registry does not admit` / nothing was swapped / `ailang-worldd setup` with that binary, or row 139's nominate |
| F11 | sessions active | `2 episodes are working: w223e (last call 40 s ago), w188m (in flight)` / a restart fails their calls for a few seconds / wait, `session revoke --episode …`, or `service upgrade --pause` |
| F12 | new build not healthy | (exit 2) `build def456 did not answer /v1/health within 30 s; rolled back to abc123, which is serving` / the log's last lines, quoted / read `<log>`; fix and re-run `service upgrade --binary <file>` |
| F13 | rollback failed | (exit 3) `rollback failed too: the service is stopped` / both builds kept: `bin/ailang-worldd` (new) and `bin/ailang-worldd.prev` (old) / `mv ~/.ailang-world/bin/ailang-worldd.prev ~/.ailang-world/bin/ailang-worldd && ~/.ailang-world/bin/ailang-worldd service start` (restores the old build), printed with real paths |
| F14 | CLI ≠ daemon | `this ailang-worldd (<path>, sha def456) is not the installed one (abc123)` / `episode start` requires the same build (P2–P3) / run `~/.ailang-world/bin/ailang-worldd`, or install this one: `<path> service upgrade` |
| F15 | exec tool moved and broken | `exec node /opt/homebrew/bin/node → Cellar/node/26.12.1 does not start: <the dyld line>` / the toolchain probe failed after the target changed / `brew reinstall node`, then retry (no restart needed) |
| F16 | plist edited by hand | (status `!`) `the plist differs from the one installed` / its sha256 ≠ `service.json` / `service start` shows the diff, or `service install --replace` |
| F17 | stale rig walk | `why`/`provenance` given an entry from a rig trailer: the indexes overlap across stores (P9) / `--addr` names the home store, not `store=<rig path>` / `why --expect-store <path>` refuses a mismatch (M4); start the rig daemon read-mostly from `daemon-cmdline.txt` on 47330 to walk it |

## §7 The trial rig: freeze it in place, start the home fresh (R4)

**The rig's store never moves and is never deleted.** Every `World-Provenance` trailer already
merged names its `world.db` by path, and indexes are per store (P9). Copying the history into the
home store would make those trailers point at a store that is no longer appended to, and leave two
stores with the same history. So:

1. **Install beside it.** The home service runs on R1's port with its own store, credentials,
   workspace root and sessions, while the rig keeps serving its live episodes (this design's own
   episode, `h2610271506`, runs on the rig) on 47330. The writer locks are per database, so the two
   don't interact.
2. **Bootstrap the home store** as QUICKSTART §9 bootstraps any fresh store: `service stop`, the
   attended `world-publish transitions --store ~/.ailang/world/world.db …`, `service start`, then
   the genesis commit. Row 219 adds no new bootstrap path.
3. **Point new work at home:** `hosts init --replace` from the installed binary; the herdr
   plugin's `world.env` lines that `service install` printed; exec profiles are shared
   (`~/.ailang-world-exec/profiles`), so they need nothing.
4. **Drain the rig.** No new mint there. When the rig's `session list` shows no unexpired
   session, or Mark revokes the rest, stop the rig daemon. It was never a launchd job (friction
   24), so there is nothing to unload.
5. **Freeze it.** The rig directory stays at its path. Walking an old trailer means starting the
   rig daemon by hand from `daemon-cmdline.txt` on 47330, as today (F17). Its `ws/` and
   `exec-cache` (16 GB, P12) wait for row 204, which must keep `world.db` and its archive.
6. **Record the cutover** in `episodes.md`/the friction log: the last rig entry index, the first
   home entry, and the date.

## §8 Rulings (proposed; recorded as D-WORLD-NN when made)

- **R1, the one address.** **(A, default)** home on `127.0.0.1:7644` (`daemon.DefaultAddr`), with
  `episode start`, `hosts init` and the plugin defaults moved to it, so the rig and home can run
  side by side during §7. **(B)** home on 47330 after the rig stops: no default changes, but no
  overlap.
- **R2, upgrade while sessions work.** **(A, default)** refuse, with `--pause` after the `y` as the
  explicit override. **(B)** always refuse.
- **R3, the verb's name.** Row 219 says `ailang-worldd upgrade`. **(A, default)**
  `ailang-worldd service upgrade`, one verb family, with the row text updated at pick. **(B)** both
  spellings.
- **R4, the rig's store.** **(A, default)** freeze in place and start the home fresh (§7). **(B)**
  copy the history (the trailers' paths would then name a dead store).

## §9 Open questions

1. **Stale against what?** The now card can show "running ≠ installed" and build age. "Behind
   `dev`" needs the daemon to read a checkout, or row 220's merge watcher to tell it. Default:
   leave it to 220.
2. **Service events in World's log.** Should an upgrade commit a `world/service-event/v1` object,
   so "why did entry N run on build X?" has a provenance answer (clause 5)? It needs a commit path
   without an agent session. Default: no, `service.json`'s history only, and revisit with row 226.
3. **Log rotation while serving** (launchd never rotates). Default: rotate on restart only, warn
   at 100 MiB.
4. **A reboot with no login** (the studio as a server). A LaunchAgent needs the user's login, and
   FileVault blocks pre-login boot anyway (P15). A root LaunchDaemon is out of scope by design. The
   exit criterion is measured as *serving within 60 s of login*.

## §10 Milestones

| M | what | gate |
|---|---|---|
| **M0** | eager build identity (P1); `ailang-worldd version [--json]`; health's `started_at`, `supervisor`, `active` | Go tests; the AC15 swap-before-first-health test |
| **M1** | `service install/status/stop/start/uninstall` on macOS behind a `launchctl` seam (a fake in tests, `launchctl` in production), the TTY fence, `service.json`, plist rendering and validation, log rotation | Go tests with the fake seam; the operator ACs 1–4, 10–11 on the Mac |
| **M2** | `service upgrade` with the active-session refusal and rollback | Go tests (fake seam, a failing candidate); operator ACs 7–9 |
| **M3** | per-call exec resolution (§3.6), **after row 156** | Go tests with a re-pointed symlink fixture; confinement arms re-run; its own judge (confinement-adjacent) |
| **M4** | §5's pointers: one default address (R1), the doctor `service` check, whoami, episode start's messages, hosts init, the herdr `world-config-init` defaults, `why --expect-store`, the now-card line, QUICKSTART with test-bound command lines | Go tests; `sh integrations/herdr/ailang-world/test/run.sh`; ACs 12–14 |
| **M5** | the cutover (§7), attended | ACs 5–6 on the Mac; the cutover record |
| **M6** | the systemd user unit (§3.5) | Go tests; AC16 on a Linux host |

## §11 Acceptance criteria (no model call; an operator or CI can check each)

`L=com.sunholo.ailang-world.worldd`, `P=~/Library/LaunchAgents/$L.plist`,
`B=~/.ailang-world/bin/ailang-worldd`, `A=http://127.0.0.1:7644` (per R1).

1. **The fence.** `B service install </dev/null` run with no controlling terminal (from a script
   or CI) exits 1 with `refusing: no controlling terminal`, and `test ! -e $P` holds. Answering
   anything but `y` at the prompt exits 1 and also leaves `test ! -e $P` (Go test: no write before
   the answer).
2. **The plist.** After `y`: `plutil -lint $P` passes; `stat -f %Lp $P` is `644`;
   `/usr/libexec/PlistBuddy -c 'Print :ProgramArguments:0' $P` is `$B`; the `--bind` value is a
   loopback host; `KeepAlive` is true; `ExitTimeOut` ≥ 30; and `EnvironmentVariables` holds only
   `PATH`.
3. **The build.** Within 30 s of `y`, `curl -s $A/v1/health | jq -r .build.executable_sha256` is
   `sha256:$(shasum -a 256 $B | cut -d' ' -f1)`, and `.supervisor` is `$L`.
4. **KeepAlive.** `kill -9 <pid from service status>`: within 45 s, `service status` exits 0 with a
   new pid and the same build.
5. **Reboot** (operator, attended). Reboot and log in, typing no World command: health answers
   within 60 s of login.
6. **`brew upgrade`** (operator, attended; CI uses AC6's fixture). After a Homebrew upgrade that
   moves `node`'s Cellar target, the next `workspace-exec` call succeeds with no restart, and
   `/v1/health.exec_tools` shows the new target. Fixture form (Go test): re-point a symlink
   between two node stubs and the call follows it; delete the target and the call refuses with
   F15's three lines, never exit 134.
7. **One-command upgrade.** Build a binary whose bytes differ (`go build -ldflags=-X=main.x=1`),
   then run `<it> service upgrade` and answer `y`. Then:
   - `shasum $B` equals the candidate's, and health's `build` equals it;
   - `$B.prev` holds the old bytes;
   - `$B episode start probe-ep --repo <a git repo>` passes its daemon check and stops at the y/N
     prompt, where the operator answers `n`, so nothing is minted.
8. **The active-session refusal.** Within 5 minutes of a `call ailang-read` on an unexpired
   session, `service upgrade` exits 1, naming that episode, and `shasum $B` is unchanged.
9. **Rollback** (Go test, fake seam). A candidate whose health never matches leaves the old sha at
   `$B`, the old build serving, and exit 2; a fake whose second start also fails ends `bootout`'d
   with exit 3 and both files present.
10. **Uninstall.** After `y`:
    - `test ! -e $P` holds, and `launchctl print gui/$(id -u)/$L` exits non-zero;
    - `world.db`, `credentials.db`, the sessions directory, the worktrees, `$B` and the log are
      all still present, with unchanged sha256s for the two databases.
11. **Nothing else is touched.** Around install, upgrade, stop, start and uninstall:
    - `shasum ~/Library/LaunchAgents/* | grep -v $L` is identical before and after;
    - `launchctl print gui/$(id -u)/dev.ailang.mission-world`'s state line is unchanged;
    - `git status --porcelain tools/launchd` is empty;
    - a Go test pins the service package's write set to {`$P`, `~/.ailang-world/bin/*`,
      `~/.ailang-world/service.json`, `~/.ailang-world/service.lock`, the log directory};
    - `service install` with a `dev.ailang.`-prefixed label refuses.
12. **Three-part messages.** Reproduce F2–F7 and F9–F17 each (the fake seam in CI, real
    conditions where the Mac allows: a stopped service, a held port, `launchctl setenv` of the
    key). Each prints a first line, a `  why:` line and a `  next:` line, and exits with §6's code.
    A grep for `HTTP 0` and for a bare `exit status 134` over their output finds nothing.
13. **Doctor and whoami.** With the service stopped:
    - `doctor --quick` prints daemon ✗ with fix `ailang-worldd service start` and exits 10, and its
      only subprocess is still the host's `--version` (Go test on the exec seam);
    - `whoami` exits 10 with F3's next.
14. **Stable scripts.** After `B hosts init --replace`:
    - `grep -rl "$B" ~/.ailang-world/hosts/bin` lists every generated script that calls World, and
      no other `ailang-worldd` path appears;
    - after AC7's upgrade, `B hosts init` (without `--replace`) succeeds with every file unchanged.
15. **Eager identity** (Go test). Start the daemon from a temp copy, rename another binary over
    that path *before* the first `/v1/health`, then query: `build.executable_sha256` is the
    original's.
16. **Linux** (until M6). `service install` exits 1 with F1's three lines and writes nothing.

## Sources

- Charter row 219 (`design_docs/world-mission.md` line 6281), and rows 196, 204, 209, 211, 212,
  217, 223 and 227; the guardrail at line 1500.
- [`w-all-agent-work-through-world`](w-all-agent-work-through-world.md) §1–§3 and §5;
  [`w-agent-host-adapters`](w-agent-host-adapters.md) §2.5 and §3.
- [`friction.md`](../verification/world-dogfood/friction.md) entries 13, 23, 24 and 32 (no entry 33
  as read).
- [`docs/QUICKSTART.md`](../../docs/QUICKSTART.md): "Run your agents through World", §0, §1 and
  "Start a whole episode in one act", §2.
- Code: `cmd/ailang-worldd/{main,episode,doctor,setup,hosts,hosts_awareness,provenance}.go`;
  `host/daemon/{daemon,build}.go`; `integrations/herdr/ailang-world/README.md`.
