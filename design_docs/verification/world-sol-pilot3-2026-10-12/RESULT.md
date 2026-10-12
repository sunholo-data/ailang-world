# Sol pilot 3, ledger iteration 13 (2026-10-12, attended): the "fixed World overhead" reversed sign

Pre-registered in `preregistration.json` (`de55d8e`). Config is byte-identical to iteration 12 (digest
`430748e2…`, `row93-direct`). The only change is the harness's per-event timing (PR #247, `58c4fd5`). Codex only,
gpt-6.1-sol, 23 core tasks, N=1. Runner output: `iter-13-runner.log`. Rows: `world-floor-m4-2026-10-06/iter-13/`.

**Instrument check (pre-registered): PASS.** 46/46 rows carry `timing`; startup + turn + teardown equals wall
within 1% on every row.

| | iter 11 | iter 12 | **iter 13** |
|---|---|---|---|
| O (median per-task overhead) | +20.02% | +23.21% | **−14.4%** |
| Δ pass rate | 0 | −4.35pp | −4.35pp (`config_file_parser`, World) |
| median wall, shell arm | 30.4 s | 29.0 s | **37.4 s** |
| median wall, World arm | 35.1 s | 36.3 s | **29.3 s** |
| shell window (CEST) | 12:30–12:43 | 13:42–13:53 | 16:00–16:15 |
| World window (CEST) | 12:44–12:57 | 13:54–14:08 | 16:16–16:28 |

## Attribution (median World − shell, per task)

| component | Δ |
|---|---|
| wall | −4.2 s |
| startup (spawn → first event) | +0.0 s |
| model turns (turn time minus tool items) | **−4.9 s** |
| tool calls | +1.0 s (World ≈ 250 ms per call × ~4 calls) |
| teardown | −0.2 s |
| tool calls per task | 4 vs 4 |

- **H_startup: false.** MCP connection adds nothing measurable before the first model request.
- **H_turns: false.** The arms make the same number of tool calls.
- **H_tool: false as the main cause.** World tool calls cost a real, stable ≈ +1.0 s per task, which matches the
  zero-spend probe (+0.8 s).
- **H_latency: the remaining delta is model-turn time, and its sign follows *when* the arm ran, not the arm.** Per-call model time
  was 6.4 s for shell and 5.7 s for World in this run. Shell was 8 s *slower* than in iterations 11/12 and World was 7 s *faster*,
  with identical config.

Caveat on the instrument: codex reports shell `command_execution` items only once they finish (≈15 ms per item), so shell tool time
is counted inside its model time. Per the zero-spend probe that is ≈ 50–70 ms per call, too small to change the picture.

## What this means for row 93

The harness runs **every shell task first, then every World task** (here and in the FINAL, §4.11 order). The two arms sit in
different ~15-minute windows of provider latency. Iterations 11–13 put the same arm-to-arm gap at +5.9 s, +7.9 s and −4.2 s with only
a wording change between them. So the "+6–8 s fixed World overhead" is **not established as a World cost**. The only
World-attributable cost measured is the ≈ +1 s per task of MCP tool latency, about +3% at a 30 s median task.

Proposed for the next FINAL preregistration (needs Mark's ruling; it changes §4.11's phase order): interleave the arms, i.e. alternate
shell/World per task with the order counterbalanced, so provider drift falls on both arms equally. Otherwise O measures the clock.

Spend: $0 Claude. Codex tokens on the existing session.
