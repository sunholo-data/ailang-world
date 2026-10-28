# w-episode-one-command (queue row 217): one short command starts an episode, and its agent — r0

**Status:** r0, written attended 2026-10-28 through World itself (episode `h2610281352` on the
trial rig; Mark: *"Write the design doc for row 217 (w-episode-one-command) … per
w-all-agent-work-through-world §3 step 2. Cover: the preset file (repo, base branch, remote,
profile, grants, TTL, agent, model, store and daemon address) and where presets live; the checks
section row 218 reads (D-WORLD-100 R1: trial.py's grader shape); episode start --preset <name>
--task <file> with the episode id and branch derived when omitted; --launch starting Claude or
Codex from an empty directory with the transcript recorded beside the session; and the mint
staying Mark's typed y. Read design_docs/planned/w-episode-finish.md first; 217 must give it what
§2 and §3.2 expect. List open rulings as R1…Rn with defaults. Log friction as you go, starting at
43."*). This is a design only. No code exists yet, and it has not been quorum-reviewed. Whoever picks
up the row re-measures every premise marked **verify** in §1 (they concern the `claude` and
`codex` CLIs, which this session can't run) and runs the quorum as usual. Roadmap:
[`w-all-agent-work-through-world`](w-all-agent-work-through-world.md) §2 row 2, built in §3 step
2's order (`219 → 217 → 218`). Rulings it needs: §6 R1–R10. Friction logged while writing:
**43** and **44**.

## §0 What is wrong today

Row 217's own words: **"A MINT IS A 1,000-CHARACTER COMMAND, AND LAUNCHING IS A SECOND STEP."**

- **The mint names everything every time.** w188l's mint named the repo, branch, profile, model,
  TTL, the store, workspace and sessions paths, the daemon address and 10 `--grant` flags
  (roadmap §1). The herdr plugin hides this behind `world.env`, but only for herdr, and only one
  set of grants for every repo.
- **Launching is a second step.** `episode start` prints the agent's command line. Mark runs it,
  or `trial.py` does. The transcript is recorded only when `trial.py` runs the agent.
- **Nothing records how the episode started.** The credential holds the episode, grants,
  expiry and exec profile binding. It doesn't hold the repo, the base commit, the agent, the
  model or the task. Row 218 needs all of them (§2).
- **Resuming an episode is a second full mint** (friction 25, owned by this row).

**Exit criterion (row 217):** Mark starts a Codex episode and a Claude episode on two different
repos with one short command each, and each agent runs.

## §1 Premise verification log

Each premise is **measured** (read in this worktree's code on 2026-10-28, base `9ec31fc`),
**quoted** (from a charter row, ruling, design or friction entry, not re-measured), or **verify**
(the behaviour of a host CLI that this session can't run).

| # | premise | source | status |
|---|---|---|---|
| P1 | Row 217's item: launch presets, one named file per repo (repo, profile, grants, TTL, agent, model, store and daemon), so a mint is `episode start --preset <name> --task <file>` then `y`. The episode id and branch are derived when omitted. `--launch` starts Claude or Codex from an empty directory with the transcript recorded beside the session. The mint stays Mark's typed `y` (or row 201's key). | `design_docs/world-mission.md` row 217 | measured |
| P2 | **`--preset` already exists, and means something else.** On `episode start` and `session new` it names a grant set, and the only one is `se-tools` (`unknown --preset %q (known: se-tools)`). `sePresetGrants` is the union of se-tools' declared (effect, scope) pairs, each given `--budget`. | `cmd/ailang-worldd/session_new.go` (`sePresetGrants`, `sessionNewGrants`) | measured |
| P3 | `episode start` today: flags `--repo --branch --profile --task --agent claude\|codex --model --ttl (3600) --preset --grant --budget --db --credentials-db --workspace-root --sessions-dir --exec-profiles-dir --addr (127.0.0.1:47330)`. It validates everything, checks that the daemon on `--addr` is its own build and pin and reads its credentials store and profile (`checkDaemon`), then runs `provisionEpisode`: the `/dev/tty` fence, the y/N, the worktree, the mint. It writes `<sessions>/<ep>.session` and, for claude, `<ep>.mcp.json` (both O_EXCL, 0600), then **prints** the launch line. "It never runs the agent." The `--task` path is only printed. | `cmd/ailang-worldd/episode.go` | measured |
| P4 | **The worktree is cut from the main checkout's `HEAD`.** `host/worktree.Add` runs `worktree add -b <branch> <dir>` with no start point, and nothing records the base commit. | `host/worktree/worktree.go` (`Add`); [`w-episode-finish`](w-episode-finish.md) P16 | measured |
| P5 | **The credential is all that a mint records.** `SessionRow` holds the credential id, episode, grants, expiry, creation time, and the exec profile's project and digest. `/v1/session/self` serves the episode, worktree, expiry, grants with use, entries, address and surface, plus a `git` block (base, lag). Nothing records the repo, remote, base branch, agent, model or task. | `host/store/credentials.go` (`SessionRow`); `host/authority/mint.go`; `host/daemon/session_self.go` | measured |
| P6 | `session revoke --episode <ep>` deletes every credential row of the episode. Anything kept on those rows is gone after revoke. | [`w-episode-finish`](w-episode-finish.md) P15 | quoted |
| P7 | `prepareMint` already **adopts** an existing git worktree at `<workspace-root>/<ep>` (`r.reuse`). `episode start` still refuses when `<ep>.session` or `<ep>.mcp.json` exists, so continuing an episode whose session expired meant moving files aside by hand and a second long mint (friction 25). | `session_new.go` (`prepareMint`); `episode.go`; friction 25 | measured / quoted |
| P8 | **The herdr plugin's `world.env` is the nearest thing to a preset.** It holds `WORLDD`, `WORLD_ADDR`, `WORLD_DB`, `WORLD_WS`, `WORLD_SESSIONS`, `WORLD_HOSTS`, `WORLD_PROFILE`, `WORLD_TTL`, one `WORLD_GRANTS` string of 10 `--grant` flags for every repo, `WORLD_REPOS` (`<workspace>=<repo>[:<profile>]`) and `WORLD_MODELS` (`claude=… codex=… pi=`). The popup's default episode id is `h$(date +%y%m%d%H%M)`. It mints with `episode start`, then starts the agent from `$HERDR_PLUGIN_STATE_DIR/empty/<ep>` (0700, must be empty) through `hosts init`'s launcher. It records no transcript. | `integrations/herdr/ailang-world/bin/world-config-init`, `bin/world-agent-popup` | measured |
| P9 | `hosts init` writes `bin/claude-world` (`claude --settings <deny Bash, Edit, Write, Read, Glob, Grep, NotebookEdit, WebFetch, WebSearch> --strict-mcp-config --mcp-config <bearer ${WORLD_SESSION}>`) and `bin/codex-world` (`env CODEX_HOME=<hosts>/codex codex`, read-only sandbox, shell tool off). Each refuses with a three-part message when `WORLD_SESSION` is unset. The token never reaches a file other than the session file. | `cmd/ailang-worldd/hosts.go` | measured |
| P10 | **`trial.py`'s grader shape:** `{name, argv, cwd, timeout_s}`. `argv` is a non-empty list of non-empty strings. `name` defaults to `basename(argv[0])` and must be unique. `cwd` defaults to `.` and must be relative and inside the worktree. `timeout_s` has a default. | `scripts/trial/trial.py` (spec load, `run_grader`) | measured |
| P11 | `trial.py`'s headless launch: the agent starts from an empty `void_dir` (refused if not empty), with a 0600 MCP config deleted after the run, and the World context block from `GET /v1/session/self` (`--append-system-prompt` for Claude, `-c developer_instructions=` for Codex). The raw stream goes to `<evidence>/transcripts/<ep>.jsonl`. The prompt is the task text verbatim plus a fixed channel note. The argv comes from `arms.build_argv`: `claude -p <prompt> --model <m> --output-format stream-json --verbose --no-session-persistence …` or `codex exec --json --ephemeral --ignore-user-config …`. | `scripts/trial/trial.py` (`run_trial_agent`, `world_context`, `build_prompt`); `scripts/floor/arms.py` (`build_argv`) | measured |
| P12 | **An exec profile's commands are ids, not argv.** `workspace-exec` takes a command id and args matched against that command's grammar. The profile fixes the program, its leading arguments, cwd and env. **`timeout_ms` is one value for the whole profile, at most 9,000** (D-140-4). | `host/broker/exec_profile.go` (`ExecMaxTimeoutMS`, `ExecProfile`) | measured |
| P13 | D-WORLD-100 (resolved, all defaults): R1 a `checks` list in row 217's preset, in `trial.py`'s grader shape; R2 checks run in the episode's exec sandbox as recorded `Workspace.Exec`, with a per-check `unsandboxed: true` named and confirmed before it runs; R6 an operator route for finish's records. | `design_docs/world-mission.md` D-WORLD-100 | measured |
| P14 | D-WORLD-98 R1 (resolved): the home service is on `127.0.0.1:7644`, and the defaults of `episode start`, `hosts init` and the herdr plugin move to it. `~/.ailang-world/service.json` (0600) records the address, store, credentials and workspace paths, and row 219's design names "row 217's presets" as one of its readers. | D-WORLD-98; [`w-worldd-home-service`](w-worldd-home-service.md) §3.2 | measured / quoted |
| P15 | D-WORLD-99 R3: a worktree's lag is measured against the repo's local `dev`, "with a per-preset override when row 217 lands". | D-WORLD-99 | measured |
| P16 | D-WORLD-96: row 224's standing allowance is approved **per preset** (the preset's grants and repos, 12 h TTL, 10 sessions a day). An allowance has to bind *which* preset bytes it approved. | D-WORLD-96 | measured |
| P17 | `claude --session-id <uuid>` fixes the session id, and an interactive session's transcript is written to `~/.claude/projects/<cwd with / as ->/<uuid>.jsonl`. A positional prompt starts an interactive session with that first message. Interactive `codex` writes `rollout-*.jsonl` under `$CODEX_HOME/sessions/YYYY/MM/DD/`, and `codex "<prompt>"` starts interactively with that prompt. | `claude --help`, `codex --help`, their docs | **verify** |
| P18 | Episode ids match `^[a-z0-9][a-z0-9-]{0,63}$`. | `cmd/ailang-worldd/doctor.go` (`workspaceEpisodePattern`) | measured |

**P12 changes row 218's example.** [`w-episode-finish`](w-episode-finish.md) §3.2 shows
`{"name": "go-test", "argv": ["go", "test", "./..."], "timeout_s": 900}` as a sandboxed check.
Under P12 it can't be one: no profile command runs past 9 s, and a sandboxed check must be one of
the profile's commands. §3.2 below catches this when the episode starts, not when it finishes
(R5). Until row 145's long-run lane exists, a full test suite is an `unsandboxed: true` check.

## §2 Conflict surface, and what row 218 gets

**Row 218 (finish) is the main reader.** Each thing it expects, and where 217 provides it:

| 218 needs | where in 218 | 217 gives it |
|---|---|---|
| the preset's `checks` | §2, §3.2, R1 | `checks` in the preset (§3.1), validated at start (§3.2) and **pinned** at mint (§3.4), so editing the preset file mid-episode can't change what finish runs |
| repo, base branch, remote | §2, §3.6–3.7, §7 Q4 | required preset fields. At mint the remote is resolved to its URL and recorded (§3.4) |
| the base commit | §3.3 (attribution against "its base"), §3.8 (`"base": "<sha>"`), P16 | the worktree is cut from `<base>` at a recorded sha (§3.3, R3) |
| the branch | §3.1 step 2, §3.9 | recorded. A fix-up round starts on an existing branch with `--branch` (§3.3), so finish checks the recorded branch, not `world/<ep>` |
| the task title | §3.4 (commit title) | the task's first `# ` heading, else its first line, recorded with the task's sha256 (§3.4) |
| agent and model | §3.5 (body line 1) | recorded (§3.4), including any flag override |
| checks that can actually run sandboxed | §3.2, R2 | each sandboxed check is matched to a profile command and its timeout at start (R5) |
| a record the agent can't forge | R6 | the start record is written with the credential, in the credentials store, under its lock (R2). No session token can write it |
| renewing an expired session | F2, §7 Q3 | `episode start --resume <ep>` (§3.6, R10). 218's `--renew` can call the same path |

`episode report <ep> --json` gains a `start` block (§3.4). Finish reads it there, and its
`writes` and `finish` additions sit beside it.

**Other rows:**
- **Row 219 (home service).** Presets read `service.json` for the daemon, store, credentials,
  workspace and sessions paths (P14) unless a preset names them. Defaults follow D-WORLD-98 R1
  (`7644`). The rig stays reachable with a preset that names `47330` and the rig's paths.
- **Row 223 (host adapters).** Stage B, "the popup uses presets", is this row's. `world.env`'s
  `WORLD_REPOS` lines become `<workspace>=<preset>`, and the popup runs `episode start --preset
  <p>` and then `episode launch <ep>` in its new pane (§3.5), so herdr and the CLI record
  transcripts the same way. `preset import-herdr` (§3.7) converts an existing `world.env` once.
- **Row 224 (standing allowances).** An allowance approves one preset **digest** (P16). Editing
  the preset ends the allowance's cover, and the next start asks for a `y`. 217 makes the digest
  stable (§3.1) and shows it on the prompt.
- **Row 228 (base and lag).** The context block's base becomes the recorded base, not a guess, and
  lag is measured against the preset's `base` (D-WORLD-99 R3's per-preset override).
- **Row 201 (keys)** replaces the typed `y` for this mint as for any other. 217 doesn't change the
  fence.
- **Row 225 (run wrapper)** takes `--preset` with the same meaning and file.
- **Row 226 (scorecard)** reads the start record (agent, model, preset) and the launch record
  (wall time, exit).
- **Row 222 (attended session as episode)** mints through the same preset path from `/world-start`.
- **`trial.py`** keeps its own spec for paired and graded trials. Its grader validation and
  `world_context` are the shapes this row ports to Go (P10, P11), with shared test vectors.
- **Row 196 (the loop)** stays gated. `--launch` is attended: it runs after a typed `y` at a
  terminal.

## §3 The design

### §3.1 The preset file

One JSON file per preset. A repo can have several presets (`world-codex`, `world-claude`).
Schema `world/preset/v1`:

```json
{
  "schema": "world/preset/v1",
  "name": "world-codex",
  "repo": "~/dev/ailang-world",
  "remote": "origin",
  "base": "dev",
  "profile": "ailang-world",
  "agent": "codex",
  "model": "gpt-6.1-sol",
  "ttl_s": 43200,
  "grants": [
    {"effect": "Workspace.Read",  "scope": "worktree", "budget": 5000},
    {"effect": "Workspace.Write", "scope": "worktree", "budget": 5000},
    {"effect": "Workspace.Exec",  "scope": "worktree", "budget": 600}
  ],
  "daemon": "service",
  "launch": {"mode": "interactive", "deadline_s": 3600},
  "checks": [
    {"name": "go-vet", "argv": ["go", "vet", "./cmd/..."], "cwd": ".", "timeout_s": 9},
    {"name": "go-test-all", "argv": ["go", "test", "./..."], "cwd": ".", "timeout_s": 900,
     "unsandboxed": true}
  ]
}
```

| field | required | meaning |
|---|---|---|
| `name` | yes | must equal the file's basename and match the episode grammar (P18). `se-tools` is reserved (R8) |
| `repo` | yes | the main checkout. `~` is expanded. Must be a git repository, not a linked worktree |
| `remote` | yes | a remote **name** in `repo`. Its URL is resolved and recorded at mint |
| `base` | yes | the branch the worktree is cut from and the PR targets (218 §7 Q4). Also the lag reference (P15) |
| `profile` | no | the exec profile's project (`--profile`). Required when a grant is `Workspace.Exec` or any check is sandboxed |
| `agent` | yes | `claude` or `codex` (`pi` once `episode start` prints a pi line) |
| `model` | no | passed to the agent. Empty: the host's default |
| `ttl_s` | yes | session lifetime in seconds |
| `grants` | yes | an explicit list. `{"set": "se-tools", "budget": N}` expands to today's preset (P2), so a short preset stays possible. The prompt always lists every grant |
| `daemon` | no | `"service"` (default): address, store, credentials, workspace root and sessions dir from `service.json` (P14). Otherwise an object with `addr`, `db`, `credentials_db`, `workspace_root`, `sessions_dir`, each falling back to the D-CLI-4 defaults. The rig is a preset with `"addr": "127.0.0.1:47330"` and the rig's paths |
| `launch` | no | `mode` (`interactive` or `headless`, R6) and `deadline_s` (headless only) |
| `checks` | no | §3.2 |

Unknown keys are refused, as exec profiles refuse them. The **digest** is
`sha256` of the file's exact bytes, as the exec profile's is (P12). It is shown on the prompt,
pinned at mint and bound by row 224's allowances.

**Where presets live (R1).** `~/.ailang-world/presets/<name>.json`, beside `hosts/` and
`service.json`. The directory is 0700 and each file 0600. A symlink, a group- or world-writable
file, or a directory inside the workspace root is refused, with the same layout rule
`prepareMint` applies to the credentials store. **Grants and checks are authority.** A preset in
the repo would be writable by the agent through its own worktree, and once merged it could widen
the next episode's grants or weaken the checks finish runs. "One named file per repo" (P1) is met
by naming: presets are kept per repo, operator-side.

### §3.2 The checks section (row 218 reads it)

Exactly `trial.py`'s grader shape (P10, D-WORLD-100 R1) plus 218's one flag:
`{name, argv, cwd, timeout_s, unsandboxed}`. The same validation as `trial.py` applies: argv
non-empty, names unique, cwd relative and inside the worktree, and a default `timeout_s` of 600.
`test_trial.py`'s grader-spec vectors are reused as Go test vectors so the two can't drift.

217 adds one check at load and again at start, because under D-WORLD-100 R2 a sandboxed check runs
as a `workspace-exec` call (P12):

- **A sandboxed check** (no `unsandboxed`) must **match one command of the preset's profile**: its
  `argv` begins with that command's fixed argv, the rest is admitted by the command's grammar,
  and its `cwd` is the command's. Its `timeout_s` must be at most the profile's `timeout_ms`
  (≤ 9 s). The match is recorded as `{"exec": "<command id>", "args": [...]}` beside the check in
  the pinned copy (§3.4), so finish calls exactly that. **No match is a refusal at start**, which
  names the check, the nearest command, and the fix (`"unsandboxed": true`, or a profile command).
  Without this, a bad check would surface only at finish, after the work is done (R5).
- **An unsandboxed check** is shown on the mint prompt as `UNSANDBOXED`, so Mark sees at start
  that finish will run agent-written code on the host. Finish still asks its own `y` before it
  runs (218 §3.2). The prompt at start doesn't replace that one.

217 never runs a check. Running checks is 218's job.

### §3.3 `episode start --preset <name> --task <file>`

```
ailang-worldd episode start [<episode>] --preset <name> [--task <file>] [--launch [--headless]]
    [--branch <name>] [--agent claude|codex] [--model <m>] [--ttl <s>] [--resume <episode>]
```

**Derived when omitted:**
- **Episode id (R4):** `<preset>-<yymmddhhmm>` from the local clock, e.g.
  `world-codex-2610281352`. If `<workspace-root>/<id>` or `<sessions>/<id>.session` already
  exists, `-2` … `-9` is appended, and a tenth collision is refused. The result must match P18,
  and a preset name of more than 53 characters is refused at load for that reason.
- **Branch:** `world/<episode>`, a new branch, as today. With `--branch <existing>` (a fix-up
  round, 218 §3.9 and §7 Q1), the worktree checks out the existing branch with no `-b`. It is
  refused if that branch is checked out in another worktree, and the prompt says
  `existing branch`.
- **Base (R3):** the worktree is cut from the repo's **local** `<base>` ref, never `HEAD`. This needs
  a start point in `host/worktree.Add` (P4), which gains one argument. The mint does no fetch and no
  network. The prompt shows the base sha and how far it is behind `<remote>/<base>` as last
  fetched, so a stale base is visible before the `y`, not found by the agent (friction 35).

**Flags over the preset (R9).** `--agent`, `--model`, `--ttl` and `--branch` override their field.
Each override is named on the prompt (`model gpt-6.1 (preset: gpt-6.1-sol)`) and recorded. Grants
and checks are not overridable with `--preset`: `--grant` beside `--preset` is refused, and the
remedy is to edit the preset (a new digest) or use another one. Without `--preset`, every flag of
today's `episode start` works unchanged.

**Order.** `episode start` keeps its order (P3). Everything is validated first: the preset file and
its layout, the checks (§3.2), the repo, the base ref, the remote, the task file, and the derived
id and its files. Then the daemon is checked (`checkDaemon`, plus "its store is the preset's").
Then the fence. **The fence and the `y` are unchanged:** `provisionEpisode` is called as today,
and there is no `--yes`, no stdin answer, and no environment variable that answers it. The prompt
gains the preset block:

```
episode world-codex-2610281352 from preset world-codex (sha256:3f9a…, ~/.ailang-world/presets/world-codex.json)
  repo      ~/dev/ailang-world, remote origin = git@github.com:sunholo-data/ailang-world.git
  base      dev @ 9ec31fc (origin/dev as last fetched: same commit)
  worktree  ~/.ailang/world-ws/world-codex-2610281352 on new branch world/world-codex-2610281352
  agent     codex gpt-6.1-sol, launched interactive from ~/.ailang/world/launch/world-codex-2610281352
  task      t.md "Fix the surface index follow" (sha256:7c1e…)
  profile   ailang-world (sha256:e938…)
  checks    go-vet (sandboxed: go-vet), go-test-all (UNSANDBOXED: finish asks again)
  grants    Workspace.Read=worktree:5000, Workspace.Write=worktree:5000, Workspace.Exec=worktree:600
  session   12h0m
Create the worktree and mint a session for episode world-codex-2610281352? [y/N]
```

### §3.4 What the mint records (the start record)

In the same credentials-store transaction as the credential (R2), one row in a new table
`episode_starts`, keyed by episode, is written. **It is not deleted by `session revoke`** (P6):
revoke ends the credential, and the start stays as history. The row holds:

```json
{"schema": "world/episode-start/v1", "episode": "<ep>", "started_at": "<RFC 3339>",
 "preset": {"name": "world-codex", "path": "…", "sha256": "…", "bytes": "<the exact file>"},
 "repo": "<abs path>", "remote": {"name": "origin", "url": "<resolved>"},
 "base": {"branch": "dev", "sha": "<40 hex>"}, "branch": "world/<ep>", "branch_existed": false,
 "agent": "codex", "model": "gpt-6.1-sol", "overrides": ["model"],
 "task": {"path": "<abs>", "sha256": "…", "title": "…"},
 "checks": [{"name": "go-vet", "argv": […], "cwd": ".", "timeout_s": 9,
             "exec": "go-vet", "args": []}, …],
 "ttl_s": 43200, "resumed_from": null,
 "launches": [{"mode": "interactive", "started_at": "…", "ended_at": "…", "exit": 0,
               "transcript": "<sessions>/<ep>.transcript.jsonl", "transcript_sha256": "…",
               "host": "codex 0.x", "argv_redacted": […]}]}
```

The preset's bytes are kept, so the pinned checks are exactly what Mark approved even if the file
changes or is deleted. The task's bytes are copied to `<sessions>/<ep>.task.md` (0600) for the
same reason. Its sha256 is in the row.

The daemon reads it through the credentials store it already opens. `episode report <ep>` shows it
as `start`, and `/v1/session/self` gains `base`, `branch`, `preset`, `agent` and `model`, so the
agent's context block states its base from a record instead of computing it (row 228).

### §3.5 `--launch`: the agent starts from an empty directory

After the `y` and the mint, `--launch` runs the same steps as a new verb,
`ailang-worldd episode launch <ep>`. The verb exists so that a crashed agent can be restarted and
herdr can launch in its own pane without minting again. `episode launch` asks no `y`: it starts an
agent on a session Mark already approved, and it refuses when that session has expired or been
revoked, or when another launch of the episode is running (a 0600 lock file
`<sessions>/<ep>.launch.lock` holding the pid).

1. **The empty directory.** `<sessions>/../launch/<ep>/` (default `~/.ailang/world/launch/<ep>`) is
   created 0700. It is refused if it isn't empty, or if it lies inside the workspace root or the
   worktree. This is the plugin's rule (P8) and `trial.py`'s (P11), now in one place.
2. **The context block.** `GET /v1/session/self` is rendered as `world_context` does (P11), as a
   Go port pinned to the same text by a shared test. It is refused if the token appears in it.
3. **The argv** by agent and mode (R6):
   - **interactive** (default): `hosts init`'s launcher, `<hosts>/bin/<agent>-world`, with
     `WORLD_SESSION` set from the session file in the child's environment only, `--model`, the
     context block (`--append-system-prompt`, or Codex `developer_instructions`), and the task
     text plus the channel note as the first prompt (P17). For Claude, `--session-id <uuid>` is
     added so the transcript's path is known in advance. No `.mcp.json` is written: the launcher's
     config reads the token from the environment (P9), so `--launch` writes one fewer copy of the
     token.
   - **headless** (`--headless`, which needs `--task`): `arms.build_argv`'s World argv (P11), with
     the per-episode `.mcp.json` for Claude, deleted when the agent exits, as `trial.py` does.
4. **The run.** The agent runs as a child in its own process group. Interactive mode inherits the
   terminal. Headless mode sends stdout to the transcript and stderr to `<ep>.agent.err`, with
   `launch.deadline_s` as the wall cap. SIGINT and SIGTERM are forwarded. `episode launch` waits
   for the child.
5. **The transcript, beside the session (R7):** `<sessions>/<ep>.transcript.jsonl`, 0600.
   - **headless:** the stream itself (Claude stream-json, Codex `--json`), as `trial.py` writes it;
   - **interactive:** the host's own session file, copied when the agent exits. For Claude that is
     the file named by the `--session-id` it was given. For Codex it is the one rollout under the
     launcher's `CODEX_HOME` whose recorded cwd is this launch directory (P17). If none is found,
     the launch record says `transcript: missing` with the reason. It is never silently absent.

   Each later launch appends to the transcript, separated by a `{"world_launch": n}` line. The
   transcript is scanned for the token before it is kept (`trial.py`'s rule). A hit truncates the
   file to a single refusal line and records `transcript: scrubbed`.
6. **After exit.** A `launches[]` entry is added to the start record (§3.4), and the lock and any
   `.mcp.json` are removed. The next step is printed: `ailang-worldd episode finish <ep>` (row
   218) and `episode report <ep>`. **The session is not revoked when the agent exits.** Finish
   closes it (218 §3.8), and Mark may relaunch before then.

### §3.6 Resume: `--resume <ep>` (friction 25)

`episode start --resume <ep>` re-mints a session for an episode whose session has **expired**.
It is refused while one is live (the remedy is `episode launch <ep>`). It reads the pinned preset
from the start record, not the file, so the grants and checks stay those of the original start. If
the file's digest has changed since, the prompt says so and lets Mark choose: `--resume` keeps the
pinned bytes, and a new `episode start` uses the new ones. It adopts the worktree (P7). The expired
`<ep>.session` and `.mcp.json` are renamed `<ep>.session.expired-<unix>` (still 0600), never
deleted. The `y` is fresh, and the start record gets `resumed_from` with the old credential id.
218's F2 `--renew` can be this path with a finish-only TTL.

### §3.7 The `preset` verbs

```
ailang-worldd preset list                 # name, repo, agent, digest, last used
ailang-worldd preset show <name> [--json] # the file, its digest, and every derived value
ailang-worldd preset check <name>         # step 1 of episode start, with no fence and no mint
ailang-worldd preset init <name> --repo <dir> --base <branch> [--agent …] [--profile …]
ailang-worldd preset import-herdr [--config-dir <dir>]   # one preset per WORLD_REPOS line
```

`preset check` is §3.3's whole validation: the file and its layout, the checks against the
profile (R5), the base ref, the remote, and the daemon match. It is read-only, so herdr's popup
and row 227's `preflight` can run it before asking for anything. `init` and `import-herdr` write a
new file once (O_EXCL) and print it. They never overwrite a file, which is `world-config-init`'s
rule.

## §4 Failure modes (three-part, row 227)

| # | condition | first line → `next:` |
|---|---|---|
| F1 | no preset by that name | `no preset world-codx in ~/.ailang-world/presets (have: world-codex, sunholo-claude)` → `ailang-worldd preset list` |
| F2 | preset file unsafe | `preset world-codex is group-writable (0640)` / a symlink / inside the workspace root → `chmod 600 <file>` or move it |
| F3 | bad field | `preset world-codex: base "deve" is not a branch of ~/dev/ailang-world` → fix the field, then `preset check world-codex` |
| F4 | sandboxed check can't run sandboxed | `check go-test-all: no command of profile ailang-world runs "go test ./..." within 9 s` → add `"unsandboxed": true`, or a profile command (row 145 for long runs) |
| F5 | id collisions exhausted | `world-codex-2610281352 … -9 all exist` → pass an episode id |
| F6 | daemon mismatch | `checkDaemon`'s messages, plus `the daemon serves store X; preset world-codex names Y` → the preset's `daemon`, or `service status` |
| F7 | `--grant` with `--preset` | `grants come from preset world-codex` → edit the preset, or drop `--preset` |
| F8 | launch dir not empty | `~/.ailang/world/launch/<ep> is not empty` → remove its contents; the session is untouched |
| F9 | session expired at launch | `episode <ep>'s session expired at 04:34` → `ailang-worldd episode start --resume <ep>` |
| F10 | host CLI missing or launcher stale | `doctor --host`'s message (row 223a) → `ailang-worldd hosts init --replace` |
| F11 | interactive transcript not found | the launch record says `transcript: missing (no rollout under <CODEX_HOME> for <dir>)` → not fatal; recorded as friction for the scorecard |
| F12 | resume of a live session | `episode <ep> has a live session (3h left)` → `ailang-worldd episode launch <ep>` |

## §5 What this design does not do

It doesn't answer the `y`, delay it, or batch it: the mint stays Mark's typed act until row 201
or 224 replaces the typing, not the approval. It runs no checks (218 does), commits nothing, and
pushes nothing. It doesn't fetch. It doesn't start a daemon (`episode start` never has). It
doesn't make presets World objects: a preset is operator-local configuration, like an exec
profile, and only its digest and bytes travel into the start record.

## §6 Rulings (proposed; recorded as D-WORLD-NN when made)

- **R1, where presets live.** **(A, default)** operator-local `~/.ailang-world/presets/<name>.json`,
  0600 in a 0700 directory, refused inside any workspace root. **(B)** in each repo under
  `.ailang-world/presets/`, read from the main checkout's `<base>` commit (`git show`), never
  from a worktree. B travels with the repo, but every merged agent PR could then change the
  grants and checks of the next episode.
- **R2, where the start record is written.** **(A, default)** an `episode_starts` table in the
  credentials store, in the mint's own transaction, kept across revoke. **(B)** a World object
  through an operator route, which builds 218 R6's route one row early. **(C)** a 0600 file
  beside the session. B puts the start record in the log with everything else, but it needs that
  route before 218, and C is invisible to the daemon.
- **R3, the base.** **(A, default)** the repo's local `<base>` ref, with no fetch. The prompt
  shows its sha and its lag behind `<remote>/<base>` as last fetched. **(B)** `git fetch <remote>
  <base>` before the prompt, which is always current but needs network and Mark's credential at
  every mint.
- **R4, the derived episode id.** **(A, default)** `<preset>-<yymmddhhmm>`, `-2`…`-9` on
  collision. **(B)** herdr's `h<yymmddhhmm>` (P8), the same everywhere. **(C)** a slug of the
  task's title, more readable, but two runs of one task collide by design.
- **R5, sandboxed checks in the grader shape.** **(A, default)** keep `argv` (D-WORLD-100 R1). A
  check without `unsandboxed` must match a profile command and fit the profile's timeout, and
  is refused at start otherwise. The matched command id is pinned. **(B)** a sandboxed check
  names the command (`"exec": "go-vet", "args": […]`) instead of `argv`, which is more direct but
  is not `trial.py`'s shape.
- **R6, how `--launch` runs the agent.** **(A, default)** interactive in the current terminal
  through `hosts init`'s launchers, with `--headless` for unattended runs that need `--task`.
  **(B)** always headless, as `trial.py` runs it.
- **R7, the transcript.** **(A, default)** structured: the headless stream as written, and the
  interactive host's own session file copied after exit. **(B)** a pty recording (`script -q`)
  for every mode, which needs no host knowledge but is raw terminal bytes that no parser in the
  repo reads.
- **R8, the meaning of `--preset`.** **(A, default)** `--preset` names a launch preset. `se-tools`
  is reserved and keeps today's meaning (the grant set) with a one-line note, and `session new
  --preset se-tools` is unchanged. **(B)** a new flag (`--use <name>`) for launch presets,
  leaving `--preset` as it is, at the cost of the charter's wording and two flags that mean
  "preset".
- **R9, flags over a preset.** **(A, default)** `--agent`, `--model`, `--ttl` and `--branch`
  override, and each override is named on the prompt and recorded. Grants and checks are closed.
  **(B)** a preset is fully closed, and any flag beside `--preset` other than the id, `--task`,
  `--launch` and `--resume` is refused.
- **R10, resuming an expired episode.** **(A, default)** `episode start --resume <ep>` (§3.6):
  the pinned preset, the adopted worktree, the old session moved aside, a fresh `y`. **(B)**
  none in 217, leaving it to 218's finish-only `--renew`, but friction 25 was a mid-work resume,
  not a finish.

## §7 Open questions

1. **A default preset per repo.** Should `episode start --task t.md` with no `--preset` pick
   one from the current directory's repo? Default: no. The agent's directory is empty by design,
   and inferring the repo from the operator's shell is the ambient state World avoids. Herdr maps
   workspace to preset explicitly.
2. **TTL headroom for finish (218 §7 Q3).** Should `ttl_s` include finish time? Default: no, and
   `--resume` or 218's `--renew` covers an expired session.
3. **Model names.** Nothing checks that `model` exists without a model call. Should `preset check
   --online` run the host's no-model `--help` probe (friction 31's lesson)? Default: yes, in
   `--online` only.
4. **pi.** `episode start` has no pi launch line (P8 mints pi as codex). Should `agent: "pi"` be
   allowed in a preset, launched through `pi-world`? Default: after 223's pi certification.

## §8 Milestones

| M | what | gate |
|---|---|---|
| **M0** | preset schema, loader and layout checks; the checks section and its profile match (R5); `preset list/show/check/init` | Go tests over fixture presets and profiles; `test_trial.py`'s grader vectors pass in Go; ACs 4, 5, 13 |
| **M1** | `episode start --preset`: derivation, `Add` with a start point, the prompt block, the start record (R2), `start` in `episode report` and `/v1/session/self` | Go tests with a fixture repo, the fake terminal and a fixture clock; ACs 1–3, 6–8, 12 |
| **M2** | `episode launch` and `--launch --headless`: the empty directory, the context block, the transcript, the launch record, the lock | Go tests with a fake agent binary on `PATH`; AC9 |
| **M3** | interactive launch and the host-transcript copy; `--resume` | fake host binaries writing to the P17 paths; ACs 10, 11; P17 verified on the Mac |
| **M4** | herdr stage B (`WORLD_REPOS` → presets, `import-herdr`, the pane runs `episode launch`); the exit run | the plugin's `test/run.sh`; AC14 |

## §9 Acceptance criteria (no model call; an operator or CI can check each)

1. **One short command.** With fixture preset `p` (repo whose `HEAD` is on `feature`, `base:
   dev`), `episode start --preset p --task t.md` under the fake terminal answering `y` creates the
   worktree at `dev`'s sha (not `HEAD`'s), on `world/p-<yymmddhhmm>`, with the session file and a
   start record whose `base.sha` is `dev`'s. The command line is under 60 characters.
2. **The prompt names it all.** The prompt shows the preset digest, the base sha, every grant,
   every check with its sandbox class, and every override. Answering `n` leaves the worktree list,
   the credentials store, the sessions directory and the launch directory byte-identical.
3. **The fence is unchanged.** Under `setsid` with stdin `</dev/null`, `episode start --preset p`
   exits 1 with nothing written. `episode start --help` lists no flag that answers the prompt.
4. **Preset layout.** A preset that is group-writable, a symlink, inside the workspace root, has an
   unknown key, or is named `se-tools` is refused with its F-line, before the daemon check.
5. **Checks validated at start.** A sandboxed check that matches no profile command, or whose
   `timeout_s` exceeds the profile's `timeout_ms`, is refused with F4 at `preset check` and at
   `episode start`, before the fence. Duplicate names, an absolute `cwd` and `..` are refused
   with `trial.py`'s messages.
6. **Pinned.** After start, editing the preset's checks leaves `episode report --json | jq
   .start.checks` and `.start.preset.sha256` unchanged.
7. **Not forgeable by the agent.** No route that accepts a session credential writes
   `episode_starts`. A test calls every session-gated route with the episode's token and finds
   the table unchanged. `session revoke --episode` leaves the row in place.
8. **Derivation.** With a fixture clock, the id is `p-2610281352`. With that worktree present, it is
   `p-2610281352-2`. With `-2`…`-9` present, F5. `--branch world/p-old` (existing) checks it out
   without `-b`, and the prompt says `existing branch`.
9. **Headless launch.** With a fake `claude` on `PATH` that records its cwd, its environment and
   the listing of its cwd, and prints a known stream: the cwd is the launch directory and it was
   empty, no environment variable holds the token, the `.mcp.json` is gone after exit,
   `<ep>.transcript.jsonl` is 0600 and equals the stream, and the launch record holds the exit
   code and the transcript's sha256.
10. **Interactive transcript.** A fake host that writes its session file at P17's path gets it
    copied after exit. A fake that writes the token into it gets `transcript: scrubbed` and a
    file without the token.
11. **Resume.** With the session expired, `--resume ep` renames the old session file (still 0600),
    adopts the worktree, asks a fresh `y` and records `resumed_from`. With the session live, it
    refuses with F12.
12. **The daemon from the service.** With `daemon: "service"` and a fixture `service.json`, the
    address and store come from it. A daemon serving another store refuses with F6.
13. **`--preset se-tools`** still mints today's grant set, byte-for-byte the same as before.
14. **The exit run** (attended, M4). Mark runs `episode start --preset <codex-preset> --task
    a.md --launch` on one repo and `episode start --preset <claude-preset> --task b.md
    --launch` on another. Each is under 80 characters plus one `y`, each agent makes at least one
    World call, and each `episode report` shows its `start` block and a transcript.

## Sources

- Charter row 217 (`design_docs/world-mission.md` line 6282), and rows 145, 196, 201, 218, 219,
  222–228; ledger D-WORLD-94, -96, -98, -99, -100.
- [`w-all-agent-work-through-world`](w-all-agent-work-through-world.md) §1, §2 (row 2) and §3 step 2;
  [`w-episode-finish`](w-episode-finish.md) (all of it; §2, §3.2, §3.9, §7 Q1, Q3 and Q4 and R6
  especially); [`w-worldd-home-service`](w-worldd-home-service.md) §3.2;
  [`w-agent-host-adapters`](w-agent-host-adapters.md) §2.3 and §3.
- [`friction.md`](../verification/world-dogfood/friction.md) entries 25, 31, 35, and 43–44 (logged
  while writing this).
- Code: `cmd/ailang-worldd/{episode,session_new,hosts,service,doctor}.go`;
  `host/worktree/worktree.go`; `host/store/credentials.go`; `host/authority/mint.go`;
  `host/daemon/session_self.go`; `host/broker/exec_profile.go`; `scripts/trial/trial.py`;
  `scripts/floor/arms.py`; `integrations/herdr/ailang-world/bin/{world-config-init,world-agent-popup}`.
