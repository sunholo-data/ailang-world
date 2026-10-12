# Sol pilot, ledger iteration 11 (2026-10-12, attended): MARGINAL

Pre-registered in `preregistration.json` (`4a0feb4`) before any inference. Codex only, `gpt-6.1-sol`,
all 23 core tasks, N=1, codex-cli 0.162.0, row-174 product (dev `054c564`, manifest `d171b695`).

| metric | value | rule |
|---|---|---|
| O (median per-task World overhead) | **+20.02%** | go ≤ +15%; marginal (+15%, +25%]; no-go > +25% |
| Δ pass rate | 0 (22/23 shell, 22/23 World) | ≥ −2pp |
| harness faults / capacity re-runs | 0 / 0 | 0 |

**Verdict under the fixed rule: MARGINAL** → reported to Mark; no FINAL is pre-registered by this result alone.

- The only failure in either arm is `config_file_parser` (logic_error), the same task Codex failed in the
  2026-10-09 FINAL.
- Per-task overheads: 11 of 23 sit above +25% (`ast_patch_roundtrip` +52.8%, `state_machine_elevator` +50.9%,
  `prompt_injection` +48.7%, `cli_args` +45.2%, `effect_tracking_io_fs` +43.6% …). The median (12th of 23,
  `merge_sort`) is +20.0%. The same 11-above-the-line shape produced +26.6% in the FINAL.
- Comparison with the FINAL (+26.6%, client 0.159.2, pre-174 product, N=3) is directional only: N=1,
  different client, different product. Neither the client nor row 174 can be credited separately.
- Informative only (tuning ledger), never pooled into a FINAL.

Evidence: `world-floor-m4-2026-10-06/iter-11/`, ledger row 11, `iter-11-runner.log`, `runtime-pre.json`.
Token scan over the evidence against the pilot session token: 0 hits (positive control on the session file: 1 hit).
