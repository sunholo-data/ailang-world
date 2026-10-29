# Launch presets (M0/M1)

A preset is private operator configuration. It pins the repository's local base
branch, daemon paths, grants, agent and checks before the terminal confirmation.
No check runs at start. No fetch or agent launch occurs. Finish runs the checks
pinned in the start record, even if the preset file later changes:
`ailang-worldd episode finish <id> --preset my-repo --dry-run` previews the local
commit and report body; omit --dry-run to keep them and print manual PR steps.

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
launch and checks. Unknown fields refuse. launch is validated for future use.
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
