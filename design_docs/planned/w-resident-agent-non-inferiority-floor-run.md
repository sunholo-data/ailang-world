# w-resident-agent-non-inferiority-floor-run — the clause-4 floor: two reference agents, shell vs World, core tier, eligibility first

**Status**: PLANNED — attended draft 2026-10-03 (read-only Plan agent; written to file by the attended controller). Queue row 93. Not implemented; no quorum yet; decisions D-NF-1..6 RULED (D-WORLD-55); depends on rows 134 and 135. Bases: ailang-world `dev` `587ee2576420…` (V1); row 134 branch `origin/attended/row134` `bb6b2b20…` (V2); ailang `76a5aef65c81…` (V4).
**Clauses**: 4 (resident-agent non-inferiority floor). Touches no other clause; the World runs *feed* the row-114 corpus (D-WORLD-52).
**Direction**: the charter's "The bar" clause 4 (V35); `w-prove-the-1-0-bar.md` §5–§7 and its Conflict Surface; D-WORLD-54 (gate = `core` tier, tuning = `smoke` tier on a cheap model, `frontier` informative, tune-then-judge with every change recorded); D-WORLD-48 (independence is a preference, not a blocker); row 134 (the World arm's 8 tools). Refinements are marked **[REFINED]** with their V-row.
**Depends on**: row 134 landed on `dev` (today it is draft PR #180, V3), plus its M6 smoke (passed 2026-10-03). Row 92 is an ordering dependency only.
**Estimate**: ~4.75 executor days (M1–M4) plus two attended sessions (M0 ≈ 0.5 d; M5 ≈ 0.5 d of attention within about 4–5 h of wall clock).

## 1. Why (measured)

- **P1 — the gate has never run, and nothing can run it today.** No queue row ever produced it (V36). The World arm's tools exist only on the row-134 branch (V2, V3).
- **P2 — the fleet harness cannot build either arm.** `ailang eval suite` does have `--agent --trials --tier --tool-policy` (V12), so w-prove's V8 ("no agent flags", measured on `eval run`) is only narrowly true. But its Claude executor passes only `--allowedTools`, which leaves the built-in tools in place, and its Codex executor always passes `--dangerously-bypass-approvals-and-sandbox` and accepts only stdio MCP servers (V13, V14). The `ailang_only` lane is pi/motoko only (V15). w-prove's conclusion therefore stands: an **adapter** is required.
- **P3 — there is no grade CLI.** `ailang eval` has no grade verb (V11). Grading is the internal `validateSolution` / `runAILANGSolution` (V9, V10), which another repo cannot import.
- **P4 — four of the 23 gate tasks cannot be test-run through World's tools.** `ailang-run` takes only `{path, args_json}`: no stdin, no argv, no `--caps`. Its policy allows only `IO` and `FS` (V29, V30). Affected tasks (V6, V8):
  - `pipeline` needs stdin;
  - `cli_args` needs argv and `Env`;
  - `api_call_json` needs `Net` and a loopback mock;
  - `prompt_injection` needs `Declassify`.

  The World agent can still write and check these four tasks. It cannot run them under benchmark conditions; the shell agent can.
- **P5 — eligibility is fragile by arithmetic.**
  - With 23 tasks, one task is 4.35 pp, so "range ≤ 5 pp" means per-run pass counts may differ by at most 1.
  - "−2 pp" at N=3 (69 pooled task-runs, 1.45 pp each) means World may lose **at most one** task-run.
  - Codex's banked fault rate (1 `api_error` in 31; V19) gives only about a 10% chance of zero faults in 69 task-runs (§9 R2). This is w-prove §6's classifier dependency, made concrete.

## 2. Goals and non-goals

**Goal.** A host tool in this repo produces, for Claude Code and codex CLI:
- an **eligibility report**, sealed and committed first;
- then a **gate verdict**: `HOLDS`, `FAILS`, `INCOMPLETE` or `PAUSED`.

Both come from paired shell and World arms on the 23 `core` tasks, N runs per arm, with the same model per agent. Every task-run is graded by one shared grader, classified, timed, attested for its effective tool set, and banked under `design_docs/verification/world-floor-<date>/`.

**Non-goals.**
- Proving superiority (clause 4 is non-inferiority by design).
- Kernel or `host/` changes (S3). World tuning that needs `host/` code is its own row through the design gate.
- Editing `scripts/bench_worldd.sh` / `bench/BASELINE.md` (w-prove Conflict Surface).
- A motoko World arm. motoko's MCP is stdio-only (row 134 §2), so it is recorded `NOT-RUN`. An optional informative motoko *shell* arm needs no new code beyond a driver table entry.
- Changing the bar's thresholds, the benchmark set or the eligibility definition.
- Any ailang-repo change (unless D-NF-1 = B).

## 3. Verification Log (all run 2026-10-03, read-only)

| V | Command | Observed (short) |
|---|---|---|
| V1 | `git -C ailang-world log -1 --format='%H %cd' --date=short` | `587ee2576420… 2026-10-03` (`dev`) |
| V2 | `git log -1 origin/attended/row134 --format='%H %cd'`; `gh pr list --state open` | `bb6b2b20618e… 2026-10-03`; `#180 Row 134 … attended/row134 DRAFT` |
| V3 | `git ls-tree -r --name-only HEAD packages \| grep -c se-tools` (on `dev`) | `0` — se-tools not on `dev` |
| V4 | `git -C ailang log -1 --format='%H %cd'` | `76a5aef65c81… 2026-10-02` |
| V5 | `grep -h '^tier:' benchmarks/*.yml \| sort \| uniq -c`; `grep -l '^tier: *core' … \| wc -l` | core **23** (19 bare + 4 with trailing comments), smoke 23, frontier 8 |
| V6 | per core YAML: `caps:`, `languages:`, stdin/expected lines | 23 ids. caps: `api_call_json` [Net,IO]; `cli_args` [IO,FS,Env]; `prompt_injection` [Declassify,IO]; six [IO,FS] or [IO] with FS; the rest [IO]. `pipeline` has `stdin:`. `contract_leap_year` is ailang-only |
| V7 | `grep -E '^(grading\|grade_entrypoint\|solution_files\|timeout):'` over core | no `grading`/`grade_entrypoint`/`solution_files`. `timeout:` 90 (config_file_parser), 180 (csv_to_json_converter), 120 ×3 (state machines, tree pipeline) |
| V8 | `cat benchmarks/{cli_args,pipeline,api_call_json}.yml` | `cli_args: ["numbers.txt"]` + `input_files`; `stdin: 1..5`; `net_allow_localhost: true` + `{{MOCK_HTTP_URL}}` in `task_prompt` (substituted at runtime, `cmd/ailang/eval_benchmark.go:74`) |
| V9 | `sed -n 46,110p;164,186p internal/eval_harness/agent_validation.go`; `sed -n 228,350p runner.go` | grade = `validateSolution` → `runAILANGSolution` (10 s) → fresh workspace `benchmark/solution.ail`, `input_files` seeded, argv `run --entry main --quiet --relax-modules [--stdlib-path <cwd>/std if present] --caps <caps> [--ai-stub] [--net-allow-http --net-allow-localhost] benchmark/solution.ail [-- cli_args]`, stdin piped. Binary = `"ailang"` from PATH (`runner.go:211-212`); only `gradeInWorkspace` honours `AILANG_BIN` |
| V10 | `sed -n 354,519p runner.go` | `CompareOutput`: TrimSpace → exact → `jsonEqual` → `normalizedEqual` (equal line count; bool case-fold, decimal canonicalisation, punctuation-space strip) |
| V11 | `ailang eval --help` | commands: analyze, browser-profile, censored-pairs, compare, elo, matrix, paired, publish, report, run, suite, summary, sweet-spot, trend — **no grade verb** |
| V12 | `ailang eval suite --help \| grep -E '^\s+-'`; `ailang eval run --help` | suite: `-agent -agent-model -agent-timeout(60) -tier -trials -tool-policy -policy-file -budget-usd -max-wall-clock …`; run: no agent/tool flags (as w-prove V8) |
| V13 | `sed -n 180,250p internal/executor/claude/claude.go` | argv `-p … --output-format stream-json --verbose --model … --permission-mode …` and `--allowedTools` only. No `--tools`, `--mcp-config` or `--strict-mcp-config` |
| V14 | `cat internal/executor/codex/mcp_config.go` | always `exec --json --skip-git-repo-check --dangerously-bypass-approvals-and-sandbox`; an MCP server **requires `command`** (stdio only) |
| V15 | `grep -rn ToolProfileAILANGOnly internal/executor` | used only by `pi/toolnames.go` and `motoko/lane.go` |
| V16 | `cat internal/eval_harness/templates/agent_task_ailang.txt`; `sed -n 170,190p agent_runner_multi.go` | template has "Available Tools … `ailang run --entry main --caps {{CAPS}} solution.ail`"; placeholder `<ws>/benchmark/solution.ail`, `module benchmark/solution` |
| V17 | `grep -n 'ErrorCategory[A-Za-z]* *= *"' internal/eval_harness/*.go` | none, compile_error, runtime_error, logic_error, **api_error ("fallback when no more specific cause is known")**, verify_error, constraint_violation, refused, timeout, quota_exhausted, rate_limit, cost_killed, step_exhausted, thrash_aborted, resource_limit, output_format, reasoning_stall, wire_drift, policy_violation, non_agentic. **No literal `harness` category**; `[harness_setup]` exists only as a stderr prefix |
| V18 | `find eval_results -path '*/agent/*' -name '*_ailang_*.json' \| … core ids \| wc -l`; keys of one row | **1979** banked agent rows on core ids. A row carries `code, stdout_ok, error_category, duration_ms, cost_usd, cost_provenance, executor, model, trial, finish_reason, agent_tool_calls` |
| V19 | `sed -n 577,625p internal/modelreg/models.yml` | gpt6-1-sol `agent_cli: codex` (codex-cli ≥ 0.159.2). Placement on core+frontier (31 tasks, N=1): **30/31, 7.2M tok, 2067 s; the one miss = api_error (codex idle 3m)**. gpt6-sol: 3 api_error in 31 |
| V20 | `sed -n 810,873p models.yml` | claude-sonnet-5-5 agent on smoke+core (42 tasks, N=1): **41/42, 3m36s wall at --parallel 3, $4.45 list-equivalent** |
| V21 | `claude --version`; `codex --version`; `ailang --version`; `$TOOL --version`; `$PIN --version` | `2.1.288`; `codex-cli 0.159.2`; PATH `ailang` = **v0.51.1-dirty**; TOOL = v0.51.0; PIN = v0.41.0 |
| V22 | `claude --help` | `--tools … Use "" to disable all tools`; `--strict-mcp-config` "Only use MCP servers from --mcp-config"; `--mcp-config` "JSON files or strings"; `--output-format json\|stream-json`; `--max-budget-usd` (print only); `--disable-slash-commands`; `--no-session-persistence`; `--bare`: "OAuth and keychain are never read" |
| V23 | `claude --help \| grep -i sandbox`; `grep -a -c` keys in the 2.1.288 binary | no sandbox CLI flag. The binary contains `autoAllowBashIfSandboxed`(29), `allowUnsandboxedCommands`(42) and `excludedCommands`(20), so the sandbox is a **settings** feature |
| V24 | `codex exec --help` | `-s read-only\|workspace-write\|danger-full-access`; `--ignore-user-config` ("auth still uses CODEX_HOME"); `--ephemeral`; `--json`; `-o`; `-C`; `--enable/--disable <FEATURE>`; `-c key=value` |
| V25 | `codex mcp add --help`; `grep -a -c bearer_token_env_var <codex binary>` | `--url` (streamable HTTP), `--bearer-token-env-var` (HTTP only); key present (27 hits) |
| V26 | `codex features list`; `grep -a -c include_apply_patch_tool` | `shell_tool` stable **true**; `unified_exec`, `multi_agent`, `browser_use`, `computer_use`, `image_generation`, `apps`, `plugins` all true; `apply_patch_freeform` removed; `include_apply_patch_tool` **0 hits** — **no apply_patch switch** |
| V27 | `codex debug prompt-input "x"` (local render, no model call) | 5 messages (skills, multi_agent role/mode, environment_context, user); **no tool specs**, so Codex's tool set cannot be checked offline |
| V28 | `codex mcp list`; `claude mcp list` | codex user config: `ailang-docs` (HTTP, **enabled**), `cua_repl`/`node_repl` (enabled), `computer-use` (disabled). claude: 4 claude.ai connectors, `ailang-parse` plugin, `obscura` ×2, `eparse` — all connected |
| V29 | `git show origin/attended/row134:packages/se-tools/transitions.json` (parsed) | 8 ids; `ailang-run` schema `{path, args_json}` with `additionalProperties:false` |
| V30 | `git show …:host/broker/handlers_ailang.go` (34–85, 433–470) | policy `allowed_caps = ["IO","FS"]`, `timeout_ms = 8000`, handler cap 10 s; run argv `run --policy P [--args-json J] -- path` (no stdin, argv or caps) |
| V31 | `git show …:cmd/ailang-worldd/session.go` (84–190); `git grep -n single-writer …host/store` | mint opens `/dev/tty`, refuses without one and asks y/N, then `store.Open` takes the cross-process **single-writer lock** (`store.go:271`). So minting needs a human at a TTY **and** the daemon stopped |
| V32 | `git show …:docs/QUICKSTART.md` §9 | 6 grants `EFFECT=worktree:N`; episode → `<workspace-root>/<ep>`, id `^[a-z0-9][a-z0-9-]{0,63}$`; publish TTY-fenced; Claude flags `--tools "" --strict-mcp-config --mcp-config <.mcp.json> --allowedTools mcp__world__ailang-read,…` (8 names) |
| V33 | `git grep -nE '270\|530 ?ms\|latency' origin/attended/row134 -- design_docs docs` | **no banked per-call latency measurement** on the branch at `bb6b2b20`; the 270–530 ms figure was reported in an executor's hand-back only. Re-measured in M0 AC0.1 |
| V34 | `ailang prompt --help` | teaching prompt comes from the binary (`--version`; `--compact` ≈ 15 KB vs ≈ 49 KB) |
| V35 | `grep -n '\| D-WORLD-5[24] \|\| D-WORLD-48 '` and the bar section of `world-mission.md` | D-WORLD-54 as quoted under Direction; D-WORLD-52: row 93's runs are row-114 corpus; clause 4 says "same benchmark set (standard tier)" |
| V36 | `sed -n 1857,1876p design_docs/world-mission.md` | row 93: "run the eligibility precondition first and report it separately … **only then run the paired arms**" |
| V37 | `sed -n 576,582p runner.go` | `stdlibPathArgs(cwd)` passes `--stdlib-path <cwd>/std` whenever cwd holds a stdlib, so a grader run from inside the ailang repo grades against its *source* stdlib |

**M2 measurements (2026-10-06, attended executor; `claude` 2.1.291, `codex-cli` 0.159.2; tuning models `claude-haiku-4-5` / `gpt-6-luna`).** Evidence: `design_docs/verification/world-floor-m2-2026-10-06/`.

| V | Command | Observed (short) |
|---|---|---|
| V-M2-1 | `codex -c web_search=bogus debug prompt-input x`; same with `=disabled`; with the misspelled key `web_serch=bogus` (no model call) | rc 1 `unknown variant 'bogus', expected one of 'disabled', 'cached', 'indexed', 'live' in 'web_search'`; rc 0; rc 0. **`web_search` is a real, enum-validated key; a misspelled key is ignored silently**, so the argv golden pins it. The codex shell canary had 0 `web_search` items and the agent answered `NO-WEB-TOOL` (AC2.6) |
| V-M2-2 | `codex mcp list --json` with an empty `CODEX_HOME` (standing in for `--ignore-user-config`, which `mcp list` does not take) plus the World arm's three `-c` overrides; then with none | `["world"]` (`streamable_http`, `bearer_token_env_var: WORLD_SESSION`); `[]`. The user config's `computer-use`/`cua_repl`/`node_repl` (V28) are gone in both (AC2.3, list half) |
| V-M2-3 | Claude shell canary, drafted argv (`--permission-mode bypassPermissions` + the §4.3 sandbox settings), then the same with `--permission-mode acceptEdits` | **The sandbox confines Bash only.** Both modes: Bash `echo > <outside>/x` → `(eval):1: operation not permitted`; inside Bash write OK. `bypassPermissions`: **the Write tool wrote outside the worktree** (AC2.5 FAIL, `canary-shell-r1`). `acceptEdits`: Write outside refused (`permission_denied`, "Path is outside allowed working directories"), inside Write OK (`canary-shell-r2` + `probes/`). init in both: `tools` = the 6, `mcp_servers` = `[]`, no WebFetch/WebSearch, `skills` = `[]`, `slash_commands` = `[]`; user plugins are still *listed* and 2 SessionStart hooks run (empty output), identically in both Claude arms. §4.3 amended to `acceptEdits` |
| V-M2-4 | codex shell canary (`-s workspace-write`), then one forced probe "run exactly `echo escape > <outside>/x`" | inside `echo` OK (a `command_execution` item); outside `echo` refused (`zsh:1: operation not permitted`, file absent); `apply_patch` outside refused (stderr `patch rejected: writing outside of the project; rejected by user approval settings`). **Neither refused call produced any `--json` item**: codex's item stream is not a complete record of attempted tool calls. AC2.3's "0 `command_execution`" is therefore necessary, not sufficient; it stands with `void/` empty and AC2.4 provenance |
| V-M2-5 | `claude --version`; `codex --version` | `2.1.291` (V21 had 2.1.288); `codex-cli 0.159.2` (unchanged). Pin both in the M0 prereg |

## 4. Design

### 4.1 Home and shape (S3)
`scripts/floor/` (Python 3, matching the existing `scripts/*.py` and `scripts/test_*.py`):
- `corpus.py`, `prompt.py`, `grade.py`, `arms.py`;
- `world.py`: per-task pre-flight and log-range capture (and AC2.4 provenance);
- M2 adds `run_task.py` (one agent spawn), `attest.py` (transcript attestation), `tokens.py` (AC2.7), `canary.py`, and the test-only `fakeworld.py`;
- `classify.py`, `stats.py`, `report.py`;
- `run.py`: the entry point.

It is host tooling, not kernel, and talks to the daemon only over `/mcp/` and `/v1/*`.

The corpus is read from the ailang repo at a pinned commit, with a sha256 recorded per YAML. Rows use the ailang result-schema field names (V18), plus the floor's own fields:
- **ailang fields:** `id, lang, model, executor, trial, compile_ok, runtime_ok, stdout_ok, duration_ms, error_category, cost_usd, cost_provenance, input_tokens, output_tokens, finish_reason`.
- **floor fields:** `agent, arm, run, wall_ms, effective_tools, tool_items, world{episode, log_from, log_to}, solution_sha256, digests{prompt, system, yaml, tool_bin, grader}`.

This is the w-prove Conflict Surface's "adapter that reuses the corpus and schema". The grader port (§4.2) is its one narrow deviation, governed by D-NF-1.

### 4.2 Corpus, prompt, grader
**Corpus.** The 23 `tier: core` ids (V5). The loader refuses any count other than 23. `--tier smoke` is for tuning only.

**Prompt.** Rendered from `agent_task_ailang.txt` at the pinned commit (V16), with the teaching prompt from `$TOOL prompt` (V34) prepended in the **task message** for both agents. Codex has no system-prompt flag, and one mechanism for both agents keeps them comparable. Per agent, the two arms differ **only** in two places:
- (a) the "Available Tools" block. The shell arm keeps `ailang run/check/test`. The World arm lists the 8 tools and says plainly that `ailang-run` has no stdin, argv or extra caps (P4).
- (b) `{{SOLUTION_PATH}}`: the absolute worktree path for shell, `benchmark/solution.ail` for World.

`{{MOCK_HTTP_URL}}` is replaced with a per-task-run mock that stays up through grading (V8).

**Grader** (the smallest one; D-NF-1 = A).
- **Input:** the worktree's `benchmark/solution.ail` bytes. Grading doesn't depend on which agent wrote them.
- **Run:** copy the file into a fresh temp directory outside any cwd that holds a `std/` (V37), seed `input_files`, and run the V9 argv with `$TOOL`, the same v0.51.0 binary the World arm uses. Don't use the PATH `ailang`, which is v0.51.1-dirty (V21). Pipe stdin and apply a 10 s timeout.
- **Compare:** a ~60-line port of `CompareOutput` (V10).
- **Output:** `{compile_ok, runtime_ok, stdout_ok, stdout, stderr}`.
- **Edge case:** a missing or placeholder-only solution is a `logic_error`, never a harness fault.

### 4.3 Arms (exact argv; fairness rules)
Common rules:
- Execution is serial.
- Each agent uses the same model id in both arms.
- Every task-run has a **600 s** deadline in both arms. That is above every spec timeout (V7), so World's latency is never cut off early.
- User MCP servers, plugins and skills are suppressed identically in both arms (V28).
- `ailang` on the shell arm's PATH is a symlink to `$TOOL`.

| | shell arm | World arm |
|---|---|---|
| Claude Code | cwd = worktree; `claude -p "$P" --model $M --output-format stream-json --verbose --no-session-persistence --disable-slash-commands --tools Bash,Read,Write,Edit,Glob,Grep --strict-mcp-config --mcp-config '{"mcpServers":{}}' --permission-mode acceptEdits --settings '{"sandbox":{"enabled":true,"autoAllowBashIfSandboxed":true,"allowUnsandboxedCommands":false}}'` | cwd = empty `void/<ep>`; `claude -p "$P" --model $M --output-format stream-json --verbose --no-session-persistence --disable-slash-commands --tools "" --strict-mcp-config --mcp-config <private>/<ep>.mcp.json --allowedTools mcp__world__ailang-read,…,mcp__world__ailang-cli` (V22, V32) |
| codex | `codex exec --json --ephemeral --ignore-user-config --skip-git-repo-check -m $M -s workspace-write -C <worktree> --disable plugins --disable apps --disable browser_use --disable computer_use --disable image_generation --disable multi_agent -c web_search="disabled" "$P"` | same flags, but `-s read-only -C <void/ep> --disable shell_tool --disable unified_exec -c mcp_servers.world.url="http://127.0.0.1:7644/mcp/" -c mcp_servers.world.bearer_token_env_var="WORLD_SESSION" -c mcp_servers.world.required=true` (V24–V26) |

**[REFINED]** Web tools are dropped from both agents' shell arms. World has none, and the bar compares "native tools" on a coding task, so web access would confound the comparison in the shell arm's favour.

Codex has no switch that removes apply_patch (V26), and its tool list can't be checked offline (V27). The World arm's defence is `-s read-only` plus an empty cwd, measured by AC2.3 and AC2.4.

Two things were **not yet measured** at draft time: the `-c web_search="disabled"` key name and the Claude sandbox settings semantics. **[2026-10-06, M2]** Both are now measured. `web_search` is a real enum key (V-M2-1). The sandbox confines Bash only, and under the drafted `bypassPermissions` the Write tool escaped the worktree. The Claude shell arm therefore runs `--permission-mode acceptEdits` (V-M2-3), which refuses an edit outside the cwd in `-p` mode and keeps sandboxed Bash auto-allowed. This is the one argv change from the draft. The goldens in `scripts/floor/testdata/argv_<agent>_<arm>.json` and `test_arms.py`'s copy of this table hold it. The codex shell arm's `-s workspace-write` refused both escape routes (V-M2-4).

**[2026-10-06]** World serves **9** se-tools since row 140 (`workspace-exec`, effect `Workspace.Exec`). The World arm keeps **exactly the 8 AILANG tools** above: the benchmark is AILANG-only, and an unconfigured exec tool would be a dead tool that changes the prompt surface. `scripts/floor/prompt.py` (`WORLD_TOOL_NAMES`, `WORLD_EXCLUDED_TOOLS`) and `scripts/floor/arms.py` (`world_allowed_tools`) pin the 8 names; `test_arms.py` checks them against `packages/se-tools/transitions.json` and asserts `workspace-exec` is absent from the prompt, the allowed tools and the grants (mutant MUT-EXEC-GRANT).

The `.mcp.json` holds the raw bearer token (V32). It lives at mode 0600 in a harness-private directory outside every worktree and void directory, and is deleted after the run. An AC greps the evidence to prove no token reached it.

### 4.4 One task-run
1. Reset the arm's worktree: check out the per-task base commit in the fixture repo (seeded `input_files` plus the placeholder), then `git clean -ffdx`.
2. World arm only, pre-flight: `/v1/health` returns 200 and `tools/list` returns exactly the 8 names. A failure here is `harness_setup` (§4.6). Record `log_from`.
3. Spawn the agent: `t0` at spawn (monotonic clock), `t1` at exit. At the 600 s deadline, kill the process group and set `finish_reason=timeout`.
4. Record `log_to` (World arm). Parse the transcript into `tool_items`: Claude `tool_use` names, Codex `item.type` counts. Assert that `void/` is still empty.
5. Grade (§4.2), classify (§4.6), and append one JSONL row.

`wall_ms = t1 − t0`. Setup and grading are excluded in both arms, and the daemon is already warm.

### 4.5 World sessions (minting is TTY-fenced and needs the daemon stopped, V31)
- **Episodes:** one per (agent, run): `fl-cc-r<k>` and `fl-cx-r<k>`. The episode's `<workspace-root>/<ep>` worktree is reset between tasks (§4.4). Per-task episodes would need 138 terminal mints per gate; per-run episodes need 2N.
- **Grants:** the 6 grants of V32, each `worktree:5000`, with `--ttl 43200`.
- **[2026-10-06] Grants are explicit, never `--preset se-tools`** (the preset grants `Workspace.Exec` since row 140). The session holds one `--grant EFFECT=worktree:5000` per effect the 8 tools declare: `Workspace.Read`, `Workspace.Write`, `Ailang.Check`, `Ailang.Run`, `Ailang.RunEnv`, `Ailang.RunNet` (row 135, D-NF-3 = B), `Ailang.Discover`, `Ailang.CLI`. That is 8 grants, superseding the 6 of V32; `arms.world_grant_args()` builds them. A session sees only the tools it holds a grant for, so the omission keeps `workspace-exec` off `tools/list` and the §4.4 pre-flight still expects exactly the 8 names.
- **The attended mint block:** with the daemon stopped, Mark runs the 2N mints, answering y each time; then `serve` starts once.
- **Tuning:** two long-TTL sessions, `fl-tune-cc` and `fl-tune-cx`, on a scratch store, minted once.
- **No batch-mint:** adding one would weaken the D1 fence, and is out of scope.
- **Per-task-run evidence:** `(episode, log_from, log_to)` is banked, and the solution provenance check (AC2.4) runs.

### 4.6 Classification (pre-registered; addresses w-prove §6)
- **PASS**: `stdout_ok`.
- **FAIL** (an outcome), in either arm: compile/runtime/logic_error, step_exhausted, non_agentic, refused, output_format, reasoning_stall.
- **FAIL, World arm only** (conservative: World pays for its own faults): `timeout`, any MCP or daemon error after a green pre-flight, and any unexplained `api_error`.
- **HARNESS-FAULT** (the bar's "api_error / resource_limit / harness classes"; there is no literal `harness` category, V17):
  - api_error, resource_limit, quota_exhausted, rate_limit, wire_drift;
  - `harness_setup` (worktree, grader, mock, pre-flight);
  - in the **shell arm only**, `timeout` (w-prove Conflict Surface).

An unexplained exit (for example Codex's "idle 3m", V19) is `api_error` until a typed cause is evidenced, and every harness-fault row carries a `cause_evidence` string.

What a fault does:
- **Shell arm:** any harness fault makes that agent **ineligible** for the attempt. The cause is fixed (instrumentation) and the **whole** shell phase re-runs. Every attempt stays banked.
- **World arm:** a quota or rate-limit fault, evidenced by the provider, re-runs that task-run once, and both runs are banked. A second such fault aborts the run, and the row **blocks on quota**. That is neither PAUSE nor a park (D-WORLD-48).

### 4.7 Statistics (exact)
Notation:
- T = 23 tasks; N = runs per arm.
- p = 1 when a task-run passes, else 0.
- r_{a,x,k} = (1/T)·Σ_t p: the pass rate of agent a, arm x ∈ {S, W}, run k.
- R_{a,x} = (1/(T·N))·Σ_k Σ_t p: the pooled pass rate.

The tests:
- **Eligible(a)** ⇔ N ≥ 3 ∧ zero HARNESS-FAULT rows in S ∧ max_k r_{a,S,k} − min_k r_{a,S,k} ≤ 0.05. At T=23 this means per-run pass counts differ by at most 1 (4.35 pp passes, 8.70 pp fails).
- **Pass-rate test:** Δ_a = R_{a,W} − R_{a,S} ≥ −0.02. At N=3 that is ΣW ≥ ΣS − 1; at N=5, ΣW ≥ ΣS − 2.
- **Overhead:** o_t = median_k wall_W(t,k) / median_k wall_S(t,k) − 1, and O_a = median_t o_t ≤ 0.25 (the 12th of the 23 sorted values).
  - Every non-harness-fault task-run counts, and a timeout counts as 600 s.
  - This paired per-task form is the pre-registered default.
- **Informative only, never the gate:** the ratio of pooled medians, the ratio of total times, a paired bootstrap 95% CI for Δ, and per-task tables (including the 4 P4 tasks).

Verdict rules, in order:
1. No agent eligible → **PAUSED** (instrumentation, not a value-park).
2. Any eligible agent fails either test → **FAILS**. After honest tuning, World parks.
3. Both agents eligible and both pass → **HOLDS**.
4. Otherwise → **INCOMPLETE** (fix or substitute the ineligible agent).

`eligibility.json` is written and committed **before** any World statistic is computed (AC3.4).

### 4.8 Evidence
`design_docs/verification/world-floor-<date>/`:
- `preregistration.json`: the config digest. It records agents, models, CLI versions, argv templates, prompt and teaching digests, the corpus commit and YAML hashes, the `$TOOL` sha256, the grader version, N, the deadline, the store, and the tuning-ledger head.
- `runs/<agent>/<arm>/r<k>.jsonl`;
- `eligibility.json`, `verdict.json`, `summary.md`;
- `transcripts/*.jsonl`, with no tokens in them (plain, not `.gz`: `verify_go.sh`'s tracked-binary gate refuses a `.gz`; measured in M2);
- `world-log-ranges.json`.

The tuning ledger is `design_docs/verification/world-floor-tuning-ledger.jsonl`.

### 4.9 Cost and quota guard (D-WORLD-54: smoke first)
- **Tuning:** smoke-tier only, on cheap models. FINAL runs only after the preregistration commit.
- **Canary:** one smoke task per arm per agent before each phase.
- **Abort:** on the first `quota_exhausted` or `rate_limit` evidence.
- **Caps:** Claude gets `--max-budget-usd` per task-run. Codex has no budget flag, so it is bounded by the per-task deadline and a phase wall-clock ceiling of 3× the R5 estimate.

### 4.10 Tuning protocol (honest, recorded)
- **Loop:** smoke tier, N=1 per arm, on `claude-haiku-4-5` and `gpt-6-luna`.
- **Allowed knobs:**
  - **harness-side, applied to both arms:** tool-block wording; the deadline (fixed before prereg); teaching prompt full vs `--compact`.
  - **World-side:** daemon/serve config, grant budgets, and product changes such as lighter read commits, the run cap or batching. A product change that touches `host/` goes through its own row.
- **Forbidden:** changing the benchmark set, the grader, the thresholds, N, or the model after prereg, or excluding tasks.
- **Ledger:** every change appends a row: `date, knob, side, before/after smoke Δ and O, commit`.
- **Final:** only the run labelled FINAL against a committed prereg counts, and it is reported whatever it shows. No re-running until it passes.

### 4.11 Run protocol (honours V36)
1. **Phase 1:** the shell arms, N runs per agent. Then `eligibility.json` is committed.
2. **Phase 2:** the World arms, N runs per agent.
3. **Drift probe:** one extra shell run per agent, informative only. If its per-task median time differs from Phase 1 by more than 10%, the report flags the overhead as "drift-confounded". The gate value does not change.
4. **Verdict.**

## 5. Premises

| P | Premise | Evidence |
|---|---|---|
| P1 | Gate set = 23 core tasks; smoke = 23 tuning tasks | V5, V35 |
| P2 | The eval harness cannot build either World arm; its Codex shell arm runs unsandboxed | V12–V15 |
| P3 | Grading needs only the solution bytes plus the spec; there is no CLI for it | V9–V11 |
| P4 | World's `ailang-run` lacks stdin/argv/caps, and 4 core tasks need them | V6, V8, V29, V30 |
| P5 | Claude's `--tools ""` removes built-ins; `--allowedTools` alone does not | V22, V13, row 134 V25 |
| P6 | Codex HTTP MCP with a bearer env var is configurable; shell has a feature switch; apply_patch has none | V24–V26 |
| P7 | Codex's tools can't be checked offline, so attestation is by canary plus event items | V27 |
| P8 | User MCP servers and plugins are active by default and must be suppressed in both arms | V28 |
| P9 | Minting is TTY-fenced and needs the daemon stopped | V31 |
| P10 | The per-call World latency is not recorded on the row-134 branch | V33 |
| P11 | The `api_error` catch-all and a ~3% Codex fault rate threaten eligibility | V17, V19 |
| P12 | The graded version is the World tool binary, not the PATH binary. It was v0.51.0; it is **v0.52.1 since row 135 M0** (2026-10-03). AC1.2 was re-run on v0.52.1 and PASSES: 1350/1357 (0.9948), zero `unexplained`. One more `version` disagreement than before: an `api_call_json` solution that names an out-of-scope constructor `Network` is now a type error. Evidence: `design_docs/verification/world-row135-m0-2026-10-03/` | V9, V21, V37 |

## 6. Milestones (each ≤ 1 executor day; ACs mechanically checkable)

**M0 (attended, 0.5 d).** Mark answers D-NF-1..6. Pin the model ids, CLI versions, corpus commit and `$TOOL` sha256. Confirm row 134 is merged with M6 banked.
- AC0.1: re-measure latency as a V-row: 20 `tools/call ailang-read`, recording the median and p90 (replaces V33).
- AC0.2: banked-row census of the harness-fault rate per agent on recent core rows, with each Codex `api_error` cause classified (the w-prove §6 check).

**M1 (1 d) — corpus, prompt, grader (no model quota used).**
- AC1.1: the loader yields exactly 23 ids plus hashes, and refuses a fixture with 22 or 24.
- AC1.2: **grader agreement.** Re-grade ≥ 200 banked rows (V18) whose YAML is unchanged since the row's timestamp. `stdout_ok` agrees on 100% of them, or every disagreement is listed with its cause and agreement is ≥ 99%. The sample covers `pipeline`, `cli_args` and `api_call_json`.
- AC1.3: grader unit table: exact, JSON-equal and normalized outputs pass; a near-miss, an extra line and wrong stdin fail.
- AC1.4: a diff test shows the rendered shell and World prompts differ only in the tools block and the solution path.

**M2 (1 d) — arm drivers and attestation (canaries spend a few smoke calls).**
- AC2.1: an argv golden per (agent, arm), checked in.
- AC2.2: Claude World canary: the stream-json `system/init` event lists `tools` == the 8 `mcp__world__*` and `mcp_servers` == [world]. The shell arm's init lists Bash and no `mcp__*`.
- AC2.3: Codex World canary ("run `echo floor-canary` in a shell") yields 0 `command_execution` and 0 successful `file_change` items, and `void/` stays empty. Under the arm's overrides `codex mcp list` shows only `world`; in the shell arm it shows none.
- AC2.4: for every World task-run, the graded solution bytes equal the content of the last `ailang-write`/`ailang-edit` of that path within `[log_from, log_to]`, read via `/v1/log` and `/v1/objects`. Otherwise the row is marked `native_write_detected`.
- AC2.5: shell confinement canary: a write outside the worktree is refused, for both agents.
- AC2.6: web-tool canary: no web tool is available in either shell arm.
- AC2.7: grepping the evidence finds 0 bearer tokens.
- **[2026-10-06] M2 status** (PR "Row 93 M2"): built in `scripts/floor/` as `arms.py` (argv builders, AC2.1 goldens), `run_task.py` (spawn, 600 s process-group deadline, transcript parse, token-file discipline), `world.py` (pre-flight, log range, AC2.4 provenance over `/v1/log` + `/v1/objects`), `attest.py` (AC2.2/2.3/2.6 checks on transcripts), `tokens.py` (AC2.7) and `canary.py`. Shell-arm canaries are measured: AC2.2 shell half, AC2.5 and AC2.6 pass on the amended argv (V-M2-1..4). AC2.2 World half, AC2.3's agent run and AC2.4 on a live run wait on the attended mint of `fl-tune-cc`/`fl-tune-cx` (`canary.py mint-block`), then `canary.py world`.

**M3 (0.75 d) — classify, stats, verdict, report (pure functions, fixture-tested).**
- AC3.1: classification table, one row per category × arm.
- AC3.2: stats fixtures, parameterised by T: Δ = −2.00 exactly passes (T=50); a range of exactly 5.00 is eligible (T=20); even and odd medians.
- AC3.3: verdict matrix covering all four outcomes.
- AC3.4: `report.py` refuses to compute any World statistic unless `eligibility.json` exists and is committed (`git ls-files --error-unmatch`).

**M4 (1 d; one attended mint) — smoke tuning.**
- AC4.1: each iteration appends a ledger row.
- AC4.2: `run.py` refuses to run with a dirty ledger.
- Deliverables: the tuned config and a prereg draft.

**M5 (attended, ~0.5 d of attention within ~4–5 h wall) — FINAL.** Commit the prereg; run the mint block; Phase 1 → commit eligibility → Phase 2 plus the drift probe → verdict → evidence; record PR.
- AC5.1: `--final` refuses unless the prereg is committed and its digest equals the live config.
- AC5.2: evidence completeness: every expected file present, and 23·N rows per arm.

## 7. Load-bearing mutations (each must turn the named test red)

| Mutant | Killing assertion |
|---|---|
| MUT-GRADE-LOOSE: CompareOutput → substring match | AC1.3 near-miss |
| MUT-NO-STDIN: grader drops stdin | AC1.2 (`pipeline` rows disagree) |
| MUT-PATH-AILANG: grader uses PATH `ailang` | AC1.2 recorded binary digest ≠ `$TOOL` |
| MUT-PROMPT-LEAK: World prompt keeps the `ailang run` block | AC1.4 |
| MUT-ALLOWEDTOOLS-ONLY: `--tools ""` dropped | AC2.2 |
| MUT-CODEX-SHELL: `--disable shell_tool` dropped | AC2.1 golden + AC2.3 |
| MUT-NATIVE-WRITE: solution taken from `void/` | AC2.4 |
| MUT-FAULT-AS-FAIL: shell `api_error` → FAIL | AC3.1 + eligibility fixture |
| MUT-DELTA-STRICT: `>` for `≥` | AC3.2 (T=50) |
| MUT-RANGE-STRICT: `<` for `≤` | AC3.2 (T=20) |
| MUT-MEAN: mean overhead instead of median | AC3.2 outlier fixture |
| MUT-SEAL-ORDER: World stats computed before eligibility | AC3.4 |
| MUT-FINAL-NOPREREG | AC5.1 |
| MUT-EXEC-GRANT: World-arm grants include `Workspace.Exec` [2026-10-06] | `test_arms.py` (exec absent from tools and grants) |

## 8. Decisions — RULED by Mark, attended 2026-10-03 (D-WORLD-55)

The bar is unchanged: point estimate, with the bootstrap CI reported as context. D-NF-1 = A (port). D-NF-2 = `claude-opus-5-5` + `gpt-6.1-sol` (tuning on `claude-haiku-4-5` / `gpt-6-luna`). **D-NF-3 = B: widen `ailang-run` first (queue row 135).** D-NF-4 = N=3, D-NF-5 = A and D-NF-6 as recommended. The original options follow.


- **D-NF-1 — grader.**
  - **A (rec.):** the in-repo port (§4.2), proven against banked rows by AC1.2.
  - B: ask AILANG core for `ailang eval grade --benchmark --solution` via ailang-feedback. Stricter reuse, but a cross-repo dependency.
- **D-NF-2 — models.**
  - **Rec.:** Claude Code on `claude-sonnet-5-5`, codex on `gpt-6.1-sol` (not Astra, per D-WORLD-48). Tune on `claude-haiku-4-5` / `gpt-6-luna`.
  - Alt: `claude-opus-5-5`. It's slower, which makes the overhead test easier to pass.
- **D-NF-3 — the 4 P4 tasks.**
  - **A (rec.):** keep all 23 (the same set), and report the four separately per task.
  - B: widen `ailang-run` (stdin, argv, caps) through its own row before FINAL, counted as World tuning.
  - C (rejected): exclude them.
- **D-NF-4 — N.**
  - **Rec.:** N=3. A larger N lowers the chance of a false fail on Δ but raises the risk the shell arm is ineligible (§9 R1/R2).
- **D-NF-5 — order.**
  - **A (rec., per V36):** phases plus a drift probe.
  - B: interleave the arms per task while still sealing eligibility first. Better for measuring overhead, but departs from the row's wording.
- **D-NF-6 — store.**
  - **Rec.:** tune on a scratch store; run FINAL on the live World store so row 114 gets its corpus (D-WORLD-52).
  - Alt: a dedicated store.

## 9. Risks and residuals

- **R1 — false fails at the ceiling (estimate).** With k coin-flip tasks, two *identical* arms still fail Δ ≥ −2 pp with probability ≈ 11% (k=1) to 24% (k=3) at N=3, and ≈ 18% at k=3, N=5. Mitigations: pre-registration and the informative bootstrap CI. The bar's thresholds are not changed.
- **R2 — eligibility (estimate).** At 1/31 faults per task-run, P(zero faults in 69) ≈ 0.10; at N=5 it is ≈ 0.02. With 3 coin-flip tasks, P(range ≤ 1 task) ≈ 0.57. Owner: M0 AC0.2, which classifies the Codex stalls before any FINAL.
- **R3 — overhead for fast agents (estimate).** Claude takes about 15 s per task (V20). At 10 World calls of 0.27–0.53 s each (V33, unverified), World adds +18% to +35%. Measured by AC0.1, and the main tuning target.
- **R4 — P4 asymmetry.** At N=3, losing one task-run of 69 fails the floor. Owner: D-NF-3.
- **R5 — cost (estimate, from V19/V20).**
  - **Codex:** ≈ 5.3M tokens and ≈ 26 min per 23-task run-arm. A gate (7 run-arms) is ≈ 37M tokens and ≈ 3 h.
  - **Claude sonnet-5-5:** ≈ $2.4 list-equivalent and ≈ 6 min per run-arm; ≈ $17 and ≈ 45 min per gate.
  - **Smoke tuning** on luna/haiku costs a fraction of this.
- **R6 — contamination.** The shell agents can read outside their worktree; for example, the ailang repo's `eval_results` holds 1979 solutions (V18). Mitigations: worktrees under a neutral root, plus an audit of absolute paths in the transcripts, flagged in the report.
- **R7 — provider drift between phases.** Covered by the drift probe.
- **R8 — row 134 not landed (V3).** M2 onward is blocked until it lands.
- **R9 — session TTL and mint fatigue.** 2N mints in one block; TTL 12 h.
- **R10 — log growth (estimate).** 23·N·2 task-runs at ~15 calls each ≈ 2k entries per gate.

## 10. Conflict surface

- **w-prove Conflict Surface** (reuse the corpus and result schema through an adapter): met. The corpus is read in place and rows use the ailang field names. The grader port (D-NF-1) is the only narrow deviation.
- **`scripts/bench_worldd.sh`:** untouched.
- **Row 134:** a read-only consumer of `/mcp/` and QUICKSTART §9; no `host/` change.
- **Row 114:** shares the World store (D-NF-6).
- **ailang repo:** no change under D-NF-1 = A.
- **Evidence directory:** `world-floor-<date>` doesn't collide with any existing `world-iter*` directory.

## 11. Standards compliance

- **S3:** host tooling, no kernel addition.
- **S6:** the gate is non-vacuous. Every gate input is attested (AC2.2–AC2.4), eligibility is sealed before the gate (AC3.4), and FINAL is pre-registered (AC5.1).
- **S1/S2:** not applicable; there is no `.ail` core.
