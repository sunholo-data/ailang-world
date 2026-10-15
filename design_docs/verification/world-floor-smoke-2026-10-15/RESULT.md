# Attended known-deltas smoke (2026-10-15, tuning iteration 14): all three checks HOLD

Pre-registered in `preregistration.json` (`ecea876`, cap note `3f45253`; config digest `01454188…`). Run
13:08–13:17 CEST, Mark attending. 7 tasks × 2 arms × 2 agents, N=1, FINAL models (claude-opus-5-5,
gpt-6.1-sol), arms interleaved in the FINAL's `paired_order`, fresh rig built at `8258341`. Two changes vs the
2026-10-14 FINAL: World wording `row93-direct`, and the corpus pinned at sunholo-data/ailang `18dc54b`
(PR #1837, `config_file_parser` says the program itself writes the file). Rows: `world-floor-m4-2026-10-06/iter-14/`;
log: `smoke.log`.

| check (fixed before the run) | result |
|---|---|
| C1 codex World/shell input tokens ≤ 1.2× (FINAL 2.03×) | **0.99×** (128k vs 129k median) — HOLDS |
| C2 `config_file_parser` passes in all 4 cells | **4/4** — HOLDS |
| C3 no regressions: 28/28 cells, 4/4 canaries, 0 harness faults, 0 attestation findings, 0 native writes, token scan clean | HOLDS |

## C4 (informative): codex World minus shell, per task

| task | tokens W/S | Δ wall | Δ model | Δ tool | FINAL tokens W/S | FINAL Δ wall |
|---|---|---|---|---|---|---|
| api_call_json (control) | 0.48 | −13.1 s | −13.2 s | +0.5 s | 0.99 | −12.9 s |
| config_file_parser | 0.97 | +0.7 s | −0.1 s | +0.9 s | 1.74 | +1.0 s |
| effect_composition | 0.75 | −2.0 s | −2.6 s | +0.8 s | 2.05 | +11.0 s |
| error_handling | 1.00 | −1.6 s | −1.9 s | +0.6 s | 1.73 | +11.0 s |
| explicit_dataflow_ssa | 1.00 | −0.6 s | −1.3 s | +0.8 s | 2.08 | +12.7 s |
| json_transform | 0.99 | −4.1 s | −4.4 s | +0.7 s | 2.03 | +20.0 s |
| state_machine_elevator | 1.00 | −5.9 s | −6.2 s | +0.5 s | 2.07 | +13.3 s |
| **median** | | **−2.0 s** | −2.6 s | **+0.7 s** | | |

The five tasks that were +11 to +20 s slower under World in the FINAL are now at or below shell. The token
excess is gone, and with it the model-time gap. What remains is World's real tool cost, ≈ +0.7 s per task.
The control task, which had no token excess in the FINAL either, was World-faster in both runs. That
matches the mechanism: the time followed the hidden discovery tokens, not World.

Small N: one pass over 7 hand-picked tasks. This confirms the mechanism and the task fix; it does not
estimate O over the 23-task corpus. That is the FINAL's job.

Claude: 14/14, O −23.8% on these tasks. Spend $7.04 Claude (Opus at ~$0.43 per run); codex on the subscription.

**Decision rule (pre-registered):** C1–C3 hold → propose the next FINAL with exactly these two changes.
