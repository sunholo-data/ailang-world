# w-real-project-trials — real projects through World, for clause 5's real questions — r1

**Status:** r1 — **RATIFIED by Mark Edmondson, attended 2026-10-16 ("yep gret - please continue")**. Mark chose the two projects:
**`sunholo-data/ailang-parse`** (AILANG) and **`sunholo-data/platform-source`** (Python, the ADK 2 upgrade
planned in its PR #69). Ratification is Mark's, attended.
**Authority:** D-WORLD-58 (2026-10-03): "dogfood now on real work … feeding row 114's corpus of real 'why'
questions"; row 140 (generic toolchain effect) LANDED; row 182 (`world/task-outcome/v1`) in PR #280.

## §1 Why

Clause 5 needs ≥ 3 REAL "why did X happen" questions, each answered by a provenance walk and verified in
≤ 5 minutes. Row 114 Q1 failed because the store held no task, outcome or grading context (row 182 fixes
that). A capability check on 2026-10-16 (`verification/world-row182/e2e/`) showed a walk answering a Q1-type
question in 20.8 s, verified in 32.8 s — but that question was posed by the controller. Real questions only
come from real operation: real work, done through World, captured, where something surprising happens.

## §2 Measured facts

| id | measurement |
|---|---|
| F1 | Claude Code runs on Mark's subscription: `claude auth status` → `authMethod: oauth_token`, `apiProvider: firstParty`; no `ANTHROPIC_API_KEY` in the environment; all 216 Claude transcripts of the 2026-10-15 FINAL report `apiKeySource: none`. codex runs on its subscription. Neither agent bills per call; both draw on plan limits. |
| F2 | World's exec tool is `workspace-exec` (`packages/se-tools/se_tools/exec.ail`), input `{command, args?}` only; the operator's `--exec-profile` (`world/exec-profile/v1`) fixes program, argv prefix, root, env and timeout; argv[0] sha256-pinned; a seven-arm probe gates `serve`. |
| F3 | Exec confinement: srt 0.0.78, **no network** (`allowedDomains: []`), writes only to the worktree and the episode cache, `$HOME` fenced, `env -i`; **timeout ≤ 9000 ms per call**; no memory/CPU limit (row 158); the agent can edit the tests that judge it (row 156). |
| F4 | Row 140 M4 (attended 2026-10-06) ran Python (`pytest -k bq tests/test_cli_bq.py`, exit 0 in 523 ms after `uv sync --frozen` and a `read_roots` fix for uv's interpreter alias), Go and TypeScript through `workspace-exec`; first calls on a cold worktree can time out — warm once outside the sandbox (F4 there). |
| F5 | `session new --repo <dir> --branch <name>` creates the episode worktree on an external repo; `--grant Workspace.Exec=worktree:N` grants N exec calls. |
| F6 | The floor's agent argv (`scripts/floor/arms.py`) allowlists 8 World tools and NOT `workspace-exec`; the trial needs its own argv with it. The floor's config digest is untouched by this design. |
| F7 | ailang-parse requires **`ailang >= 0.50.0`** (`ailang.toml`; CHANGELOG: serve-api `_headers` + `@optional`). World has TWO hard pins: the transition interpreter `--ailang-bin` **v0.41.0** (epoch 1; a restart on v0.54.11 refuses at startup: `registry: existing head … diverges from bootstrap revision`, measured 2026-10-16) and the agent tool binary `--tool-ailang-bin` **v0.52.1** (`ToolBinaryRelease`, `host/daemon/workspace.go:44`; digests in `cmd/ailang-worldd/pins.go`). v0.52.1 satisfies ailang-parse's range on paper; neither pin can move routinely today (row 139). Its dev loop is `./bin/docparse --check`/`--test`/`--prove` and `uv run benchmarks/...`. (r1 said the agent tools ran v0.41.0; corrected.) |
| F8 | platform-source pins `google-adk==1.31.1`, `ag-ui-adk>=0.6,<1`, `a2a-sdk<1` (`backend/pyproject.toml`); `make adk-conformance` → `pytest -m adk_contract` in `backend/`. PR #69 (merged) plans the ADK 2.12 upgrade; Phase 0 (`template-adk-2-regression-gate.md`) is 8 new hermetic `adk_contract` tests (stub `BaseLlm` + real `Runner`), run on 1.31.1 and once on a 2.12 venv, plus dev-Agent-Engine steps. |

## §3 Design

**D1 — two projects, real backlog only.**
- **T-PY (platform-source):** Phase 0's 8 contract tests. Episode A (Claude Code `claude-opus-5-5`): tests 1–4
  (`transfer_declaration_shape`, `before_tool_callback_contextvars`, `error_event_finality`,
  `hallucinated_tool_name`). Episode B (codex `gpt-6.1-sol`): tests 5–8 (`instruction_template_literals`,
  `after_model_replacement_metadata`, `inmemory_append_unknown_session`, `resilient_llm_interrupt_passthrough`).
  Each test must pass on 1.31.1; each is run once on 2.12 and the failures recorded (the doc's exit criterion).
  The dev-Agent-Engine round trip and the journey/eval baselines stay OUTSIDE World (they need deployed
  credentials) as attended steps for Mark.
- **T-AIL (ailang-parse):** open real bugs, one per episode, alternating agents. Proposed first two: **#67**
  (`--deep` prints a stale skipped warning) and **#20** (`outputPath` silently ignored). Further issues are
  added only from the open backlog, never invented.

**D2 — exec profiles (operator-authored, one per project).** Targeted commands only (F3's 9 s cap):
- platform-source: `adk-contract-file` = `.venv/bin/python -m pytest -q -p no:cacheprovider -m adk_contract`
  + positional `testfile` (`.py`) + flags `-k` (regex), `-x` (bool); `adk212-contract-file` = the same under
  `.venv-adk212/bin/python`; `ruff-file` = `.venv/bin/ruff check` + `relpath`.
- ailang-parse: **after row 139 lands** (Mark, attended 2026-10-16: "world needs to be able to cope with ailangs versioning - we release very often"; sequencing "widen 139, do it next"). The episode then uses World's own `ailang-check`/`ailang-run` on an interpreter matching the project's declared range (`ailang = ">=0.50.0"`), not an exec-pinned side binary. Only `quick-check` (`bash benchmarks/quick_check.sh`) is profiled, if it fits 9 s.

**D3 — prep (operator, outside the sandbox, before any mint).** Fresh rig `~/.ailang/state/trial-<date>`:
fetch both repos at `origin/dev`; create the episode worktrees through `session new --repo … --branch
trial/<issue>` at mint; per worktree `uv sync --frozen` (platform-source `backend/`), a second venv
`.venv-adk212` with `google-adk==2.12.0` resolved (recorded lock), ailang-parse's `uv sync` for its benchmark
scripts; warm each profile command once unsandboxed (F4); run the seven-arm probe via `serve`. Every digest
(profile, venv lockfiles, binaries) goes into the trial's `pins`.

**D4 — the agent run.** A small trial runner (`scripts/trial/run.py`, new; reuses
`scripts/floor/run_task.run_agent` for spawning and streaming, with a trial argv that adds
`mcp__world__workspace-exec`). Prompt = the issue or sub-doc text verbatim + a short channel note (World tools
only; exec commands available; finish by summarising the change). Deadline 45 min per episode. Transcripts
kept on the rig and hashed.

**D5 — outcomes (row 182 object, trial grader).** After the agent stops, the runner reads the episode's entry
range (head before/after) and runs the project's **grader** outside the sandbox in the episode worktree:
platform-source `make adk-conformance` (both venvs, separately), ailang-parse `./bin/docparse --check` +
`--test` + `bash benchmarks/quick_check.sh`. It commits one `world/task-outcome/v1` per grader with
`task` = the issue/sub-doc id, `arm` = `world`, `entries` = the episode range, `grading.rule` = "the operator
runs <command> unsandboxed in the episode worktree after the agent stops", `grading.argv` = that command,
`graderVersion` = `trial-grader/1`, class PASS/FAIL from the exit code, stdout/stderr heads. Same verb
(`ailang-worldd outcome record`), same session.

**D6 — PRs (outward-facing; drafts).** For each episode with changes: the operator pushes the episode branch
and opens a **draft** PR on the project with label `ailang-world` and the `World-Provenance:` trailer
(D-WORLD-60). Mark reviews and merges as usual; CI results on the PR are further real outcomes, recorded the
same way when they arrive.

**D7 — harvesting questions (row 114 rules, unchanged).** Anything surprising during or after an episode
becomes a candidate in row 114's harvest log with an event ref — e.g. a contract test passing on 1.31.1 but
failing on 2.12 for an unexpected reason, a grader failure the agent didn't see, a CI failure, a refused exec
call mid-task, an agent declaring done while a test fails. **Only questions someone actually wants answered
are taken**; none are invented. Each runs row 114's protocol: frozen bytes, E2 covering every arm/side the
question names (Q1 finding F-2), baseline B first on raw logs + repo, walk W second on a served copy, fresh
verifiers. The trial stores are declared as new cohorts (K4 platform-source, K5 ailang-parse) by sha256 before
any reveal.

## §4 Acceptance criteria

| AC | criterion |
|---|---|
| AC1 | Both projects run ≥ 1 episode each through World only (transcripts show no non-World tool use; exec calls all profiled). |
| AC2 | Every episode has ≥ 1 `world/task-outcome/v1` with its entry range; `outcome_record_errors` 0. |
| AC3 | Draft PRs opened with the label and trailer; nothing merged by the controller. |
| AC4 | The harvest log lists every candidate question with its event ref, taken or rejected with a reason. |
| AC5 | Each taken question runs row 114's protocol end to end; verdicts recorded, failures preserved. |
| AC6 | Zero changes to the floor harness config/digest, the kernel, or `tools/launchd/*`. |

## §5 Risks

| risk | mitigation |
|---|---|
| 9 s exec cap too short for ADK test files (imports) | measure at prep; split per test with `-k`; if one file cannot fit, record it and narrow the profile — never raise the cap in this trial |
| Agent edits the tests that judge it (row 156) | the grader runs outside; the PR diff shows test edits; harvested as a question if it matters |
| No real surprise happens | then no question is taken; the trial still yields captured stores and PRs; more issues follow from the backlog |
| AILANG version mismatch (F7) | T-AIL waits for row 139 (routine per-release upgrade + project version ranges) |
| Outward-facing PRs | drafts only; Mark merges |
| Platform-source secrets | worktree only; `$HOME` fenced (F3); no deployed-env commands profiled |

## §6 Order of work

1. Row 182 LANDED (`e879e1e`). 2. Ratify this doc. 3. Build `scripts/trial/run.py` (spawn + outcome grader +
draft PR helper) with tests. 4. Prep the rig and profiles; measure the 9 s fit. 5. Mark mints T-PY episode A.
6. Run, record, open drafts, harvest; then episode B. 7. T-AIL (#67, #20) after row 139 lands.

## §7 Unknown

- Whether ADK test files run inside 9 s from a warm venv (measured at prep).
- Whether `quick_check.sh` fits 9 s.
- How many real questions two to four episodes produce.
