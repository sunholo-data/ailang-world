# Floor verdict: HOLDS

| agent | eligible | Δ (point) | Δ test | O (median o_t) | O test |
|---|---|---|---|---|---|
| claude | yes | 0 (+0.0000) | pass | -2246/12525 (-0.1793) | pass |
| codex | yes | 0 (+0.0000) | pass | -984/17329 (-0.0568) | pass |

Informative only (never the gate): bootstrap 95% CI for Δ, pooled-median and total-time ratios — see verdict.json `informative`.

- drift claude: median per-task paired shell time +0.084 vs Phase 1 over 23 tasks (within 10%; informative, the paired gate value is unchanged)
- drift codex: median per-task paired shell time -0.035 vs Phase 1 over 23 tasks (within 10%; informative, the paired gate value is unchanged)
