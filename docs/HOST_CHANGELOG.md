# Host changelog

## v1.0.0-rc.1 — the 1.0 release candidate

`v1.0.0-rc.1` is a GitHub pre-release. All seven clauses of the 1.0 bar
(`design_docs/world-mission.md`, "The bar") are met; `v1.0.0` follows after queue
rows 156 (protected test paths), 150 (marketing capture) and 151 (agent-generated
projections) (D-WORLD-81). The transition interpreter stays at AILANG v0.41.0
(epoch 1); the workspace tool binary stays at v0.52.1.

- **`workspace-exec`** (row 140): a ninth se-tools transition and the
  `Workspace.Exec` effect. It runs only the commands an operator exec profile
  lists (`serve --exec-profile`, with `--exec-sandbox`), under the srt sandbox:
  writes confined to the worktree and the episode's caches, no network, and a
  startup probe each profile must pass. Without a profile every call is refused
  `no exec profile configured`. The `se-tools` session preset is nine grants.
- **Outcome capture** (row 182): `ailang-worldd outcome record --file
  <outcome.json|->` commits one graded task's outcome as a
  `world/task-outcome/v1` object under the episode's session, so a provenance
  walk reaches the task, the grading rule and the entry range. Exit 1 refused
  before sending, 2 commit refused, 3 outcome uncertain.
- **`worktree_delta`** (row 184): `ailang-run`, `ailang-check`, `ailang-cli`
  and `workspace-exec` results list the files the call created, modified or
  deleted (at most 256, on the `git-status` basis in a git worktree). `why`
  prints a `worktree: N created, M modified, K deleted` line.
- **Degraded episodes** (row 185): `session new --allow-non-git` adopts an
  existing non-git episode directory; `session list` marks it `DEGRADED` and
  `/v1/health` reports it under `degraded_episodes`.
- **Release manifest and epoch nomination** (row 139 M1/M2):
  `pins/ailang-releases.json` is the single source of the AILANG releases World
  uses, compiled into both CLIs. `world-publish epoch nominate --store S
  --release vX.Y.Z (--candidate | --new-epoch) --gate-verdict V --binary B`
  appends one epoch-registry revision; `--candidate` is headless, `--new-epoch`
  attended.
- **`why` input fields** (PR #292): when an input is elided, its short
  top-level fields (`path`, `command`, `op`, …) print first.
- **Clause 4, the floor** (D-WORLD-77): the 2026-10-15 FINAL holds. Claude Code
  Δ 0 (69/69 vs 69/69), overhead −17.9%; codex Δ 0 (69/69), overhead −5.7%;
  414/414 graded cells passed
  (`design_docs/verification/world-floor-final-2026-10-15/`).
- **Clause 5, provenance teeth** (D-WORLD-80): three unseeded real questions
  answered by walks verified against sealed ground truth (Q4 54.6 s, Q6 25.3 s,
  Q8 39.0 s); the live baseline also verified, more slowly (67.0 s, 87.5 s,
  116.8 s)
  (`design_docs/verification/w-provenance-teeth-unseeded-baselined-walk/index.md`).

## Session-filtered MCP tools — accepted in v1.0.0-rc.1

The host adds session-filtered MCP tools at `POST /mcp/`, alongside the agent card
and A2A surface. Each invocation item re-admits against the current registry and
receives a fresh task ID. Older-version batches run sequentially with independent
commits: a later failure can hide earlier committed results, and retrying can repeat
an effect. Encoded MCP names are reversible; A2A IDs remain verbatim.

Credential resolution, registry reads, aggregate requests and transport writes use
the frozen 3/10/20/30-second budgets. The shared callback runner uses 20 seconds and
eight slots; the distinct subprocess reservation limit is also eight. MCP wire and
SSE handling come from upstream module v0.47.2; the archived compiler remains
v0.41.0. Slow bodies and noncooperative callbacks retain the documented residuals.

The resend socket fixture uses deterministic cooperative execution to test durable
commit followed by a lost response and same-task reconciliation. The separate real
capsule execution/replay test remains unchanged. This entry was first recorded as an
unreleased candidate with acceptance pending; it ships in `v1.0.0-rc.1`.
