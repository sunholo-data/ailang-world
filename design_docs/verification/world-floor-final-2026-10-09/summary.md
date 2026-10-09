# Floor verdict: FAILS

| agent | eligible | Δ (point) | Δ test | O (median o_t) | O test |
|---|---|---|---|---|---|
| claude | yes | 0 (+0.0000) | pass | -958/4875 (-0.1965) | pass |
| codex | yes | -1/69 (-0.0145) | pass | 4220/15879 (+0.2658) | FAIL |

Informative only (never the gate): bootstrap 95% CI for Δ, pooled-median and total-time ratios — see verdict.json `informative`.

- drift probe claude: median per-task shell time -0.030 vs Phase 1 over 23 tasks (within 10%; informative, the gate value is unchanged)
- drift probe codex: median per-task shell time +0.035 vs Phase 1 over 23 tasks (within 10%; informative, the gate value is unchanged)
