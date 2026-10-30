# Launch presets

A preset is private operator configuration. It pins the repository's local base
branch, daemon paths, grants, agent and checks before the terminal confirmation.
No check or fetch runs at start. Add --launch --headless to launch after mint. Finish runs the checks
pinned in the start record, even if the preset file later changes:
`ailang-worldd episode finish <id> --preset my-repo --dry-run` previews the local
commit and report body. Omit --dry-run to approve the exact PR bytes at the terminal:
finish pushes without force, creates the ailang-world label if approved, opens a draft
PR (or use --ready), records the outcome, then revokes the episode credentials.
Approval and records use a one-time operator proof beside the credentials store;
the daemon refuses session credentials on this route. Protected changes are named
in the commit, body and prompt; --exec-profiles-dir locates the pinned profile.

Create and inspect a preset (replace the repository and branch with yours):

```sh
ailang-worldd preset init my-repo --repo /absolute/path/to/repo --base dev
ailang-worldd preset list
ailang-worldd preset show my-repo --json
ailang-worldd preset check my-repo
ailang-worldd episode start --preset my-repo --task /absolute/path/to/task.md
```

The last command requires a controlling terminal and your typed y. The id is
my-repo-<yymmddhhmm> and the new branch is world/<id>. The worktree starts at the
local dev commit shown on the prompt, even when the main checkout has another
branch checked out. --branch can adopt an existing free branch. --agent, --model,
--ttl and --branch overrides are all shown and recorded. --grant is refused
with a launch preset. --preset se-tools retains its existing grant-set meaning.

init creates ~/.ailang-world/presets (0700) and my-repo.json (0600), once only.
Edit that private JSON to choose grants, checks, model, TTL and daemon. The
default daemon is the home service's service.json. For a separate daemon use:

```json
"daemon": {
  "addr": "127.0.0.1:7644",
  "db": "/private/world/world.db",
  "credentials_db": "/private/world/credentials.db",
  "workspace_root": "/private/worktrees",
  "sessions_dir": "/private/world/sessions"
}
```

The complete generated file has schema world/preset/v1, name, repo, remote,
base, agent, ttl_s and grants. Optional fields are profile, model, daemon,
launch and checks. Unknown fields refuse. launch selects the mode and wall cap.
Grants use {"effect":"Workspace.Read","scope":"worktree","budget":5000}, or
{"set":"se-tools","budget":50}. Workspace.Exec requires profile.

Checks have name, argv, cwd (default "."), timeout_s (default 900), and optional
unsandboxed. For example, add a profile name and:

```json
"checks": [
  {"name":"vet","argv":["go","vet","./..."],"timeout_s":9},
  {"name":"suite","argv":["go","test","./..."],"unsandboxed":true}
]
```

A sandboxed check must match the named profile's fixed argv, argument grammar,
working directory and timeout. Full suites generally need unsandboxed: true
until the long-run lane exists. This class is shown at approval; finish will
ask again before executing host code. check --no-daemon validates locally and
prints that the daemon match was skipped. --presets-dir and --exec-profiles-dir
select fixture or alternate configuration directories.

episode report <id> --json reads the immutable start block, including the exact
preset bytes, base sha and pinned checks. Session self exposes a smaller start
summary. Revoke ends the credential and retains start history. A schema-2
credentials store must first be upgraded by serve; the CLI never upgrades it.
The fixture tests in preset_test.go exercise this workflow in temporary dirs.

## Headless launch (M2)

Mint and launch together at the attended terminal:

```sh
ailang-worldd episode start --preset my-repo --task /absolute/path/to/task.md --launch --headless
```

Or launch a previously minted, live episode:

```sh
ailang-worldd episode launch my-episode --preset my-repo --headless
ailang-worldd episode report my-episode --json
ailang-worldd episode finish my-episode --preset my-repo
```

Launch reads the preset name, agent, model, task digest and launch settings pinned
at mint. The current preset supplies the daemon and paths and must name the same
store. Without --headless, the pinned launch.mode selects the mode; its default
is interactive, which arrives in M3. For headless defaults use
`"launch": {"mode": "headless", "deadline_s": 3600}`; a missing or zero deadline
also means 3600 seconds.

With sessions_dir /absolute/path/world/sessions, the agent starts in the empty
0700 directory /absolute/path/world/launch/my-episode. A nonempty directory,
symlink or overlap with workspace/session data refuses. A 0600
/absolute/path/world/sessions/my-episode.launch.lock holds the launch pid; remove
a stale lock only after that launch has ended.

The 0600 stream is /absolute/path/world/sessions/my-episode.transcript.jsonl;
stderr is /absolute/path/world/sessions/my-episode.agent.err. Later launches
append a numbered separator. Literal session tokens cause the stream to be
replaced with a refusal line. Temporary output stays in the sessions directory.
Claude's private MCP config is deleted after exit; Codex receives the token only
in its child environment. The report's start.launches records the exit, reason,
transcript digest and redacted argv. The lock is removed, and the session stays
live so another launch or episode finish can use it. Finish is the next step.

## Interactive launch and resume (M3)

Interactive launch is the default when the pinned preset has no launch mode.
Run it from a terminal window; the agent uses that terminal for input and output.
Prepare the confined launchers once, then launch the approved session:

```sh
ailang-worldd hosts init --out-dir /absolute/path/world/hosts --addr http://127.0.0.1:47330
ailang-worldd episode launch my-episode --preset my-repo --hosts-dir /absolute/path/world/hosts
```

Launch checks the launcher digests and daemon address before taking the launch
lock. A missing or stale launcher prints the exact hosts init --replace remedy.
Start and launch together with --launch; interactive mode needs no --task.
Use --headless with a task for unattended runs. The token reaches the interactive
child only as WORLD_SESSION; no per-episode MCP config is written.

After exit, World copies the host's structured transcript into the private 0600
/absolute/path/world/sessions/my-episode.transcript.jsonl. It never moves,
changes or deletes the host's source. Claude uses CLAUDE_CONFIG_DIR, or HOME
plus /.claude, then projects/<slug>/<session-id>.jsonl. The slug replaces every
non-alphanumeric byte in the resolved launch directory with a hyphen; for long
slugs World finds the host's directory by its first 200 characters and a hyphen.
Codex uses /absolute/path/world/hosts/codex/sessions/YYYY/MM/DD/rollout-*.jsonl.
World searches the local dates and times of the run and matches the first
session_meta record's cwd to the launch directory. Zero or multiple matches
record transcript_state: missing with a reason, leave the previous copy alone,
and preserve the agent's exit outcome. Copied streams use the same token scan,
scrubbing and numbered separators as headless launch. Terminal output is separate.
Only regular host files are copied: the open refuses a symlink leaf and checks
the open file again. Files larger than 64 MiB record missing with that reason;
no copy is kept. Private staging files are removed on every outcome. The wall
cap is the session's remaining lifetime; parent cancellation records canceled.

An expired session can be resumed with its pinned grants, checks, agent, model,
TTL and existing worktree:

```sh
ailang-worldd episode start --resume my-episode --preset my-repo
ailang-worldd episode launch my-episode --preset my-repo --hosts-dir /absolute/path/world/hosts
```

Resume accepts only --preset, --presets-dir and --exec-profiles-dir beside
--resume. It asks a fresh terminal y before minting. Old .session and .mcp.json
files, when present, are renamed with .expired-<unix> appended and retain their
permissions. A new 0600 session file holds the fresh token. The start record's
resumes list records the old credential and the current preset file digest;
the original mint fields stay pinned. An edited preset is named in the prompt,
and its new grants are used only by a new episode start. A live credential
refuses with F12 and points to episode launch; revoked or finished episodes
have no credential to resume and need a new start. A credential remains live
at its exact expiry second; resume is allowed only after that second.

The current credential is the last successful resumes entry's credential_id,
or the original start credential when no resume has succeeded. Revoking that
id ends the episode even while older expired rows remain. The mint transaction
rechecks current authority and liveness, so competing resumes cannot both mint.
A failure after mint revokes the new credential, marks its history entry failed,
and restores the old file names. Failed entries do not advance current authority.
The adopted worktree must still have the pinned branch and belong to the pinned
repository, checked using both repositories' resolved git common directories.
