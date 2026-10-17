# Trial runner (row 183, `w-real-project-trials`)

Runs one real-project episode through World and records how it turned out
(design `design_docs/planned/w-real-project-trials.md` §3 D3–D6). Python 3, stdlib only; it
reuses `scripts/floor/run_task.py` (spawn, transcript, private MCP config, agent env) and
`scripts/floor/arms.py` (argv) without changing them.

## 1. Prepare (operator, outside the sandbox — D3)

Rig layout under `--root` (e.g. `~/.ailang/state/trial-<date>`):

    bin/ailang-worldd          the rig's daemon/CLI binary
    private/                   mode 0700, outside every worktree
    private/<ep>.session       mode 0600, the episode's 64-hex session token
    void/<ep>/                 empty cwd for the agent (created if absent)

Per episode: mint the session with `session new --repo <project> --branch trial/<issue>` and the
grants the tools need **including `--grant Workspace.Exec=worktree:N`** (pre-flight refuses a
session whose `tools/list` lacks `workspace-exec`); prepare the venvs; warm every exec-profile
command once unsandboxed; record every digest in the trial's `pins`.

Trial spec (`trial.json`; relative paths resolve against the spec's directory):

```json
{"project": "platform-source", "repo": "/abs/platform-source",
 "exec_commands": ["adk-contract-file", "adk212-contract-file", "ruff-file"],
 "episodes": [
  {"id": "tps-a", "agent": "claude", "model": "claude-opus-5-5", "task": "adk2-phase0-tests-1-4",
   "task_text_file": "tasks/tps-a.md", "worktree": "/abs/worktree/of/tps-a", "deadline_s": 2700,
   "graders": [
     {"name": "adk-131", "argv": ["make", "adk-conformance"], "cwd": "backend", "timeout_s": 900}]}]}
```

`graders`/`exec_commands`/`deadline_s` may sit at the top level as defaults. Optional per episode:
`void_dir`, `max_budget_usd` (omitted = no `--max-budget-usd`; the agents run on subscriptions, F1).

## 2. Run

    python3 scripts/trial/trial.py run --trial trial.json --episode tps-a \
        --root ~/.ailang/state/trial-2026-10-16 --addr 127.0.0.1:7644 [--evidence DIR]

It records `log_from = head + 1`, spawns the agent (World MCP tools only, plus
`mcp__world__workspace-exec`; prompt = task text verbatim + a fixed channel note), records
`log_to = head`, hashes the worktree change, then runs every grader argv **unsandboxed** in the
worktree (wall cap, default 900 s) and commits one `world/task-outcome/v1` per grader with
`ailang-worldd outcome record --session <rig>/private/<ep>.session`. Finally it captures the
D-WORLD-60 trailer with `ailang-worldd provenance --since <log_from> --episode <ep>`.

Evidence (default `<root>/evidence/<project>/`): `episodes/<ep>.json` (agent rc, wall, argv
redacted, transcript path + sha256, log range, solution hash, each grader's rc/heads/full-output
paths/outcome ref or error, `outcome_record_errors`, trailer), `transcripts/<ep>.jsonl`,
`graders/<ep>.<name>.{stdout,stderr}`. A failed record never raises; it is counted. Exit: 0 done,
1 done with record errors (AC2 wants 0), 2 refused or pre-flight failed (nothing graded).

## 2b. The paired shell arm (row 187 evidence loop)

An episode with `"arm": "shell"` (and `"pair_of": "<world episode id>"`) runs the SAME task text with
the same agent in a plain git worktree of the same base commit. It uses the floor's shell argv
(`arms.build_argv(agent, 'shell')`): native tools, a sandboxed shell, no web tools, Codex
`-s workspace-write`. It needs no World session, mint or daemon, and records no outcome or trailer; the
same graders run after it. Write a paired task's text channel-neutral (no World- or shell-only
instructions): only the channel note differs per arm. Prepare the shell worktree yourself
(`git worktree add -b <branch> <dir> <base>`) with the same venvs.

## 3. Open the draft PR (operator — D6)

    python3 scripts/trial/trial.py pr-body --episode-row <evidence>/episodes/tps-a.json \
        [--evidence-link design_docs/verification/<dir>/] > body.md

It prints the body (agent summary, graders with outcome refs, the `World-Provenance:` trailer as
the last line) and never pushes or calls `gh`. It refuses when no trailer was recorded, and when
anything token-shaped (a bare 64-hex run, `Bearer`) would be printed. The operator pushes the
episode branch and opens a **draft** PR with label `ailang-world`; Mark reviews and merges.

Tests: `cd scripts/trial && python3 -m unittest -v` (builds `./cmd/ailang-worldd` to check payloads
against the real validator; needs `go`).
