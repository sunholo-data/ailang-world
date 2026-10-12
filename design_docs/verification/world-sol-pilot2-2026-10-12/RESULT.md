# Sol pilot 2, ledger iteration 12 (2026-10-12, attended): MARGINAL — the mechanism moved, the time did not

Pre-registered in `preregistration.json` (`b5a5af8`). The only change from iteration 11 is
`world_tools_wording` `row135-rw` → `row93-direct` (PR #244). Codex only, gpt-6.1-sol, 23 core tasks, N=1.

| | iteration 11 (row135-rw) | iteration 12 (row93-direct) |
|---|---|---|
| O (median per-task overhead) | +20.02% | **+23.21%** |
| Δ pass rate | 0 (22/22) | −4.35pp (shell 23/23, World 22/23) |
| harness faults / capacity re-runs | 0 / 0 | 0 / 1 (`effect_composition`, §4.6) |
| World/shell input tokens (median ratio) | 1.71 | **1.00** |
| median absolute World − shell wall | +5.9 s | +7.9 s |
| tasks above +25% | 11 | 11 |

**Verdict under the fixed rule: MARGINAL** (O ≤ +25%). Δ −4.35pp is below the −2pp line on its own: it is the one
task `config_file_parser`, which passed shell this run and fails World as in every run. That is N=1 noise on
Codex's one flaky task, not a wording effect, but it is recorded as measured.

**Mechanism check (pre-registered):** the input-token ratio fell from 1.71 to **1.00**, so the wording did remove the hidden
discovery round trips. The absolute delta did **not** fall. **The discovery round trips were real but were not the
time cost.**

## Codex-side tool latency (zero-spend probe, `probes/run_timing.py`)

The mock makes codex run JavaScript that times 5 cycles of read → write → check → run inside its own `exec`, with no model:

| per call (median) | read | write | check | run | cycle |
|---|---|---|---|---|---|
| World (MCP via codex) | 216 ms | 216 ms | 344 ms | 248 ms | **1,024 ms** |
| shell (`exec_command`) | 46 ms | 47 ms | 70 ms | 62 ms | **225 ms** |

So tool execution costs World ≈ +0.8 s per task, about 10–15% of the gap. The remaining ~5–7 s is in the **model turns
themselves**: request latency or reasoning time per turn. `codex exec --json` carries no timestamps, so
these banked runs cannot attribute it further.

## Next instrument (no spend to build)

Timestamp every codex stdout event in the runner (`run_task.py`, harness-only) so the next measured run splits wall
time into model-turn vs tool-call intervals per task. Until then, any further tuning would be guesswork.
