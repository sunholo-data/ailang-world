# Floor verdict: FAILS

| agent | eligible | Δ (point) | Δ test | O (median o_t) | O test |
|---|---|---|---|---|---|
| claude | yes | 0 (+0.0000) | pass | -1291/5980 (-0.2159) | pass |
| codex | yes | -2/69 (-0.0290) | FAIL | 5543/22622 (+0.2450) | pass |

Informative only (never the gate): bootstrap 95% CI for Δ, pooled-median and total-time ratios — see verdict.json `informative`.

- drift claude: median per-task paired shell time +0.007 vs Phase 1 over 23 tasks (within 10%; informative, the paired gate value is unchanged)
- drift codex: median per-task paired shell time -0.006 vs Phase 1 over 23 tasks (within 10%; informative, the paired gate value is unchanged)
