# Row 93 M4: smoke tuning (2026-10-06)

The attended executor ran this on darwin arm64 with `claude` 2.1.291 and `codex-cli` 0.159.2. Models were `claude-haiku-4-5` and `gpt-6-luna` only. The store was the scratch tuning store `~/.ailang/state/floor-tune`, using sessions `fl-tune-cc` and `fl-tune-cx` (minted by Mark; 8 explicit grants each). Design: `design_docs/planned/w-resident-agent-non-inferiority-floor-run.md` §4.9, §4.10 and §6 M4, plus V-M4-1..4.

**Runner:** `scripts/floor/run.py iterate`. **Config:** `scripts/floor/floor_config.json`. **Ledger** (AC4.1, one row per iteration): `design_docs/verification/world-floor-tuning-ledger.jsonl`.

## Smoke set

The design defines a `smoke` **tier**: 23 tasks, all `[IO]`. That tier has no stdin, argv, Env, Net or Declassify task. The brief instead asked for a fixed core subset that covers the P4 tasks. These 6 core tasks were used:

- `pipeline` (stdin)
- `cli_args` (argv, `[IO,FS,Env]`, input file)
- `api_call_json` (`[Net,IO]` and the loopback mock on 127.0.0.1:7655)
- `prompt_injection` (`[Declassify,IO]`)
- `effect_tracking_io_fs` (`[IO,FS]`)
- `higher_order_functions` (`[IO]`)

The phase canary (§4.9) is the smoke-tier task `fizzbuzz`. It runs once per (agent, arm) before that cell's tasks. It is graded but not counted. It must pass attestation, and since iteration 1 it must also PASS the task, or the phase stops.

Run settings:
- N=1 per arm.
- A 600 s deadline per task-run.
- `--max-budget-usd 1.00` per Claude task-run (no run came near it; the maximum was $0.21).

**Cell reuse.** A cell (agent, arm) whose arm digest is unchanged from an earlier iteration is not re-run; its rows are copied. The digest covers the argv template, the teaching and template digests, the task hashes, the tool sha and the World tools block. `iteration.json` (`rows_reused_from`) records which cells were reused, and the ledger row records the same.

## Iterations

Δ = (World passes − shell passes) / 6. O = the median over tasks of `wall_W / wall_S − 1`. A task with a harness fault in either arm is excluded from O. At N=1 every number is a single draw: one task moves Δ by 16.7 pp.

| it | knob (side) | commit | claude S/W, Δ, O | codex S/W, Δ, O | Claude $ | wall |
|---|---|---|---|---|---|---|
| 0 | baseline: M2 config as merged | `c2dc2db` | 6/6, 0, **+0.342** | 6/**0**, **−1.000**, −0.504 | 1.637 | 674 s |
| 1 | `codex_mcp_approval=approve` (harness, codex World) | `dba991a` | (reused 0) 6/6, 0, +0.342 | 6/6, 0, −0.022 | 0.000 | 167 s |
| 2 | World tools wording `m1`→`row135` (World) | `60d1408` | 6/6, 0, +0.180 | 6/5, −0.167, −0.095 | 0.780 | 407 s |
| 3 | `claude_isolation=true` (harness, both Claude arms) | `64cbf71` | 6/6, 0, +0.075 | (reused 0/2) 6/5, −0.167, −0.095 | 1.577 | 367 s |
| 4 | wording `row135`→`row135-rw` (World) | `a48950c` | 6/6, 0, +0.062 | 6/6, 0, +0.027 | 0.848 | 524 s |
| 5 | none: fresh repeat of all 4 cells on the it-4 config | `a48950c`+evidence | 6/6, 0, +0.153 | **4**/6 (2 shell harness faults), +0.333, +1.302 (4 tasks) | 1.070 | 684 s |

What each knob did:

- **Iteration 0 (baseline).**
  - **codex:** World passed 0/6. Every World MCP call failed with `MCP tool call requires approval, but approval policy is never`. The exhibit is `exhibits/iter0-codex-world-pipeline.jsonl`; the codex World canary also failed. M2's codex World canary made no MCP call (it asked for a shell), so it could not see this. The codex World arm had never been functional.
  - **claude:** O = +0.34, driven by `cli_args` (+1.73) and `prompt_injection` (+1.39).
- **Iteration 1, knob (c), V-M4-3.** `-c mcp_servers.world.default_tools_approval_mode=approve` on the codex World arm only. The key is enum-validated: `bogus` → `expected one of auto, prompt, writes, approve`. codex World went 0/6 → 6/6, with O −0.02.
- **Iteration 2, knob (a).** The World tools block was rewritten to describe row 135's tool as it ships. `ailang-run` takes `stdin`, `argv` and `caps`, and the run step names the task's caps. M1's text said it had none of those.
  - **claude:** O +0.34 → +0.18.
  - **codex:** one miss, `higher_order_functions`. codex answered "this workspace is explicitly read-only, and the available AILANG tools cannot override that restriction" and never wrote.
- **Iteration 3, knob (b), V-M4-1.** Both Claude arms got `--setting-sources ""` plus `--settings {"enabledPlugins": {<3 built-ins>: false}, "autoMemoryEnabled": false}`, merged into the shell arm's sandbox settings object.
  - Every Claude init event in iteration 3 and later showed `plugins: []` and 0 hook events. The runner attests this per run.
  - The shell confinement canary under these flags still refused the Bash and Write escapes (`canary-shell-isolation/`).
  - Cost per Claude run did not drop measurably, because the ~97.5 KB teaching prompt dominates.
  - claude O fell to +0.08.
- **Iteration 4, knob (a2).** `row135-rw` adds one sentence: the read-only sandbox binds native tools only, and `ailang-write`/`ailang-edit` write the worktree. codex World went 6/6, and both agents had Δ 0 with O ≤ +0.06.
- **Iteration 5 (no-change repeat, noise).**
  - **claude:** reproduced Δ 0 and O +0.15.
  - **codex:** the shell arm hit **2 harness faults**: `Selected model is at capacity. Please try a different model.` on `api_call_json` and `cli_args` (a provider capacity refusal, `api_error`). codex World passed 6/6. With the two faulted tasks excluded, the codex O of +1.30 rests on 4 tasks: `pipeline` took 87 s in World against 12.7 s in shell.

**Stopped after iteration 5.** The last knob (it 4) and the repeat (it 5) left both agents' smoke Δ at 0. The remaining movement in O is draw-to-draw noise at N=1: Claude O ranged +0.06 to +0.18 on two identical World configs. No further allowed knob had a measured cause behind it.

**Not tuned: teaching prompt full vs `--compact` (V-M4-2).** `ailang prompt --compact` fails on the pinned v0.52.1 tool with `Error: "v0.16.6-compact" is not a known prompt version`, and `--list` has only `v0.7.4-compact`. The knob is unavailable without an older prompt, which would not be the same teaching content, so `teaching: full` stays.

## Tuned config

| | |
|---|---|
| file | `scripts/floor/floor_config.json` (sha256 `a0bf41172eda31bfd7c1e409475416da1851775d426ac3b31ce9013fd995eebc`) |
| knobs | `world_tools_wording=row135-rw`, `claude_isolation=true`, `codex_mcp_approval=approve`, `teaching=full` |
| live smoke config digest | `ce87a7e23524dd7fa8d0b072cf8e6e6c55b3aead146a91ed011a495b006d23cf` (`run.py digest --mode smoke`, after D-WORLD-67). Iterations 4 and 5 ran under `0e2e6d78…e55b`: the argv was the same, but the digest did not yet carry the pre-registered rules |
| live FINAL config digest | **`aaa9b2b6c99db7d723f5fe0806de152eeab2b2ae019821feb1174b2dc1e7714a`** (`preregistration.draft.json`, regenerated after D-WORLD-67). The first draft's `49ea2005…8464` had no `rules` block |

**Prereg draft:** `preregistration.draft.json`.
- Models `claude-opus-5-5` and `gpt-6.1-sol`; N=3; the 23 core tasks with YAML hashes at `76a5aef`.
- Thresholds −1/50, 1/20 and 1/4; the statistics and protocol as ruled.
- Pinned CLI versions; the argv templates; `--max-budget-usd 5.00` per Claude task-run (ruled by D-WORLD-67).
- The §4.6 capacity rule (`capacity_rule`, and `rules` in the live config).
- The ledger head.

M5 commits a copy as `preregistration.json`. `run.py final` refuses unless that file is committed and its `config_digest` equals the live FINAL config. A CLI upgrade, a knob, the teaching prompt or a task YAML all move the live config digest.

## Rulings (D-WORLD-67, Mark Edmondson, attended 2026-10-06)

1. **The three argv additions are canonical.** They are Claude isolation, the codex World `default_tools_approval_mode=approve`, and `--max-budget-usd` (5.00 FINAL).
   - They are in §4.3's table, the goldens, `test_arms.py`'s table copy and `arms.build_argv`'s defaults.
   - Mutation results on a copy of `scripts/floor`, each named test FAILED (killed); the unmutated `ArgvGoldens` and `CapacityRule` passed:

     | Mutant | Killing tests |
     |---|---|
     | MUT-ALLOWEDTOOLS-ONLY | `test_builder_matches_checked_in_golden`, `test_claude_world_removes_every_builtin` |
     | MUT-CODEX-SHELL | `test_builder_matches_checked_in_golden`, `test_codex_world_has_no_shell` |
     | **MUT-CODEX-NOAPPROVE** (new) | `test_builder_matches_checked_in_golden`, `test_codex_world_tools_are_approved` |
2. **The 6-task core smoke set is accepted** (§4.10).
3. **Capacity faults: option A, pre-registered** (§4.6). Only a typed `Selected model is at capacity` fault is re-run, up to 3 attempts in all, in both agents and both arms.
   - Every attempt is recorded; the last attempt is the outcome. Three capacity faults stay HARNESS-FAULT.
   - Tests in `test_run.CapacityRule`: capacity then pass → PASS with 2 attempts recorded; 3 capacity faults → HARNESS-FAULT with no 4th attempt; timeout, other api_errors (typed and untyped), a wrong answer and a pass → no retry; eligibility and the summary use the last attempt.
   - **MUT-RETRY-ANY-APIERROR** (retry any `api_error`) turns `test_no_retry_for_timeouts_other_api_errors_or_wrong_answers` red (4 failures).
   - The rule is in the prereg and in the live config digest (`test_rule_is_in_the_digest`).
4. **Stale design figures corrected:** V34 (`--compact` is unusable, and the prompt is about 97.5 KB, not about 49 KB), §4.2 (a), P4, R2, R3, R5 and R10.

Row 153 vs row 93 ordering (Q4) was not answered, so D-WORLD-65's order stands.

## Row 153 (plan-phase deadline)

| | claude | codex |
|---|---|---|
| World task-runs (30 counted + 5 canaries) | 35 | 35 (7 in it 0 made no successful call) |
| World MCP tool calls in counted runs | 140 | 101 (`mcp_tool_call` items, including it 0's failed ones) |
| `host callback timed out` in transcripts | **0** | **0** |
| daemon `context deadline exceeded` lines | **0** | **0** |

Serial load (one agent at a time) on the scratch store did not reproduce row 153. That is evidence at this load only, not that the defect is gone.

## AC2.4 live write provenance

- 84 World writes and edits were committed by the agents' `ailang-write`/`ailang-edit` across the World task-runs and canaries.
- Every World task-run's graded bytes equalled World's log fold: `solution_provenance: world`.
- **0 `native_write_detected`.**
- `void/` stayed empty in every World run.

## Mutation proofs (AC4.2, AC5.1)

Each mutant was applied to a copy of `scripts/floor` and the named test run against it:

| Mutant | Killing test | Result |
|---|---|---|
| MUT-DIRTY-LEDGER (`require_clean_ledger` sees nothing dirty) | `test_run.LedgerGuard.test_refuses_dirty_ledger` | FAILED (killed) |
| same | `test_run.LedgerGuard.test_iterate_refuses_before_any_work` | FAILED (killed) |
| MUT-FINAL-NOPREREG (`check_final` drops `require_committed`) | `test_run.FinalGuard.test_final_refuses_uncommitted_prereg` | FAILED (killed) |
| unmutated controls | the same 3 tests | OK |

Live refusal: `run.py final --prereg …/preregistration.draft.json` → `REFUSED: … is not committed (git ls-files --error-unmatch)`, rc 2.

## Spend

| | |
|---|---|
| Claude `total_cost_usd`, iterations 0–5 | **$5.91** |
| Claude, probes and the isolation canary | about $0.05 |
| codex | 12.0M tokens (input + output, summed from `turn.completed` usage; no budget flag) |
| wall clock in iterations | 2,822 s (47 min) |

## Files

- `iter-<k>/`, one directory per iteration:
  - `iteration.json`: summary, arm digests, reuse, canary, row 153, provenance, spend.
  - `runs/<agent>/<arm>/r1.jsonl`: ailang field names plus floor fields.
  - `canary.jsonl`
- `canary-shell-isolation/`: V-M4-1, AC2.2/2.5/2.6 under the isolation flags.
- `exhibits/iter0-codex-world-pipeline.jsonl`: the approval refusal.
- `preregistration.draft.json`

Full transcripts stay out of the repo. They are about 0.8 MB per iteration, under `~/.ailang/state/floor-tune/m4/transcripts/iter-<k>/`, and every row carries the transcript's `transcript_sha256` and `transcript_bytes`. Every `iter-<k>/` passed the AC2.7 token grep against both session files when it was written.
