# Row174 Codex-only smoke comparison — 2026-10-11

**Both candidate runs meet the Codex diagnostic target. Full M2 remains pending.**

| Run | Shell pass | World pass | Median per-task World overhead | Delta | Harness faults |
|---|---:|---:|---:|---:|---:|
| Original baseline (iter8) | 6/6 | 6/6 | +9.68% | +0.00pp | 0 |
| Candidate (iter9) | 6/6 | 6/6 | +8.49% | +0.00pp | 0 |
| Unchanged candidate repeat (iter10) | 6/6 | 6/6 | +10.59% | +0.00pp | 0 |

Each run uses unchanged gpt-6-luna, fixed six-task smoke, N1, fizzbuzz canaries,
teaching/grader/config/thresholds/current clients; both arms are fresh with no row reuse.
O is the median of paired per-task ratios, independently recomputed from raw rows;
it is not the ratio of the cell medians. Every recorded attempt is retained.

Canonical exact-eight descriptors shrank13947→10301 bytes; prose7430→3794,
wire14054→10368. Candidate/repeat metadata is identical. Actual daemon/publisher/
interpreter/tool/runner/config hashes and CLI versions match all three runs.

The original baseline already meets the cheap-model target. Shell timings vary,
so a favorable difference does not identify descriptor prose as the cause of faster
runtime. N1 results do not predict the historical Sol FINAL or establish statistical
improvement. The historical Sol+26.6% failed FINAL remains unchanged.

Iteration6 aborted on Claude weekly quota; iteration7 retained six passing shell
tasks and a failed World canary due port mismatch. Both aborted rows remain in the
ledger and are excluded from this complete paired comparison, never discarded.

Claude regression and the design’s fresh all-four-cell comparison/repeat remain
pending. Product remains unmerged; implementation PASS86 remains scoped to M1.
New FINAL separately needs product landing, a new committed preregistration and
human-typed sessions. The daemon is stopped after each runner invocation.

## Per-task overhead

| Task | Baseline | Candidate | Unchanged repeat |
|---|---:|---:|---:|
| api_call_json | -17.01% | -56.84% | -4.10% |
| cli_args | -36.06% | +0.27% | -24.37% |
| effect_tracking_io_fs | +80.73% | +43.55% | +25.29% |
| higher_order_functions | +22.15% | +85.08% | +29.05% |
| pipeline | +56.63% | +9.72% | -14.52% |
| prompt_injection | -2.78% | +7.25% | +57.74% |

Earlier ADMIN handoff ba79d1e exact-head remote CI is now SUCCESS, both required jobs
green (receipt handoff-ci-result.json). This receipt is not candidate/repeat CI or a
product merge. The Codex runtime evidence is committed locally and preserved in a
persistent Git bundle.
