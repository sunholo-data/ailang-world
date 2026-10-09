# Row 93 M5 — FINAL floor run, 2026-10-09 (attended, Mark Edmondson)

**Verdict: FAILS** (design §4.7 rule 2: an eligible agent fails a test). Clause 4 stays UNMET.

| agent | eligible | Δ = R_W − R_S (bar ≥ −0.02) | O = median_t o_t (bar ≤ +0.25) |
|---|---|---|---|
| claude (claude-opus-5-5) | yes (69/69 shell, range 0) | **0** — pass (69/69 World) | **−0.197** — pass (World faster) |
| codex (gpt-6.1-sol) | yes (67/69 shell, range 1/23) | **−1/69 = −0.0145** — pass (66/69 World) | **+0.266** — **FAIL** by 1.6 pp |

The miss is codex's time overhead alone. The 12th of the 23 sorted per-task overheads is
`json_transform` +0.266; 11 tasks sit above +0.25 (`pipeline` +0.94, `cli_args`/`explicit_dataflow_ssa` +0.48,
`float_eq` +0.44 …). codex's only failures in either arm are `config_file_parser` (logic_error: 2/3 shell, 3/3 World).

Run under the committed preregistration `preregistration.json` (config digest `9f6f02c0…2076`, AC5.1 checked
at the start of every phase). Order per §4.11: Phase 1 shell r1–r3 → `eligibility.json` sealed and committed
(`11c5049`) before any World row → Phase 2 World r1–r3 → drift probe → `report.py verdict`.

## Phases

| phase-run | status | claude pass | codex pass | harness faults | capacity re-runs | claude spend (CLI estimate; OAuth) | wall |
|---|---|---|---|---|---|---|---|
| final-drift-r4 | complete | 23/23 | 22/23 | {'claude': 0, 'codex': 0} | 0 | $10.69 | 18 min |
| final-shell-r1 | complete | 23/23 | 22/23 | {'claude': 0, 'codex': 0} | 0 | $10.76 | 19 min |
| final-shell-r2 | complete | 23/23 | 22/23 | {'claude': 0, 'codex': 0} | 0 | $10.74 | 18 min |
| final-shell-r3 | complete | 23/23 | 23/23 | {'claude': 0, 'codex': 0} | 0 | $10.81 | 19 min |
| final-world-r1 | complete | 23/23 | 22/23 | {'claude': 0, 'codex': 0} | 0 | $10.19 | 18 min |
| final-world-r2 | complete | 23/23 | 22/23 | {'claude': 0, 'codex': 0} | 0 | $10.13 | 20 min |
| final-world-r3 | complete | 23/23 | 22/23 | {'claude': 0, 'codex': 0} | 0 | $10.11 | 20 min |

Total estimated Claude spend $73.42 (subscription OAuth, `apiKeySource: none`); codex on ChatGPT login.

## Integrity checks
- Zero harness faults and zero capacity re-runs in every phase; every canary passed.
- World arm: all 138 task-runs have `solution_provenance = world` and `native_write_detected = false`
  (`world-log-ranges.json`, AC2.4).
- Drift probe (informative): claude −3.0%, codex +3.5% median per-task shell time vs Phase 1 — within 10%,
  so the overhead is not drift-confounded.
- Token scan (AC2.7) over this directory against the 6 session tokens: 0 findings.
- Informative bootstrap 95% CI for Δ: claude [0, 0], codex [−0.043, 0] (`verdict.json`).

## Store
FINAL ran on the durable store `~/.ailang/state/floor-final/store/world.db` (se-tools registry revision 1,
sessions `fl-cc-r1..3`, `fl-cx-r1..3`), not D-NF-6's live store, which is schema v2 and refused (row 124).

## Files
`preregistration.json`, `eligibility.json`, `verdict.json`, `summary.md`, `runs/<agent>/<arm>/r<k>.jsonl`,
`drift/runs/…/r4.jsonl`, `canary/`, `phases/`, `transcripts/` (plain jsonl), `world-log-ranges.json`.

Per §4.10 the FINAL is reported as it stands; it is not re-run.
