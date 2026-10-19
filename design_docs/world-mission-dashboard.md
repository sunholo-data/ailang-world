# World mission dashboard — iteration 283 (2026-10-19)

- **This fire:** row 156 M2a + M2b built, PR #314. A profile may declare `protected`; it loads as a canonical list plus a set digest, and a bad pattern refuses daemon startup before any sandbox, probe or tool subprocess.
- **Compatibility:** absent / `[]` keep `world/exec-profile/v1` and the raw-byte profile digest; `null`, non-arrays and non-string elements are refused naming `protected` / `protected[i]`. No consumer reads the field yet.
- **Judge:** Agent sonnet (cross-vendor), own worktree: r1 PASS 90, zero blocking; NB1 (a field-name check that the temp path satisfied) reproduced and fixed doc-only → r2 PASS 92.
- **Verification:** outside the sandbox with the pinned `AILANG_BIN`: vet, focused tests, the 51-row refusal table, broker + daemon `-race` all rc 0; commits rebuilt from snapshots, sha256-identical; PR CI 3/3 required before merge.
- **Release:** `v1.0.0-rc.1` pre-release. Proof/compiler pin 0.41.0; tool pin 0.54.15.
- **Clauses:** 7/7 MET; goal unmoved. Before v1.0.0 (D-WORLD-81): row 156 (8/21), row 150 (9/17), then row 151.
- **Next:** row 156 M3a (protected sandbox settings), M3b; M6-pre (V24 real-worktree timing) before M6a.
- **Parked on Mark:** none (ledger zero OPEN).
- **Routing/cadence:** executor codex gpt-6.1-sol, evaluator Agent sonnet; designer and planner not needed (doc current, plan ready).
- **Cost:** metered $0.00; codex 151k tok; Agent subagents ≈136k tok.
- **Retro:** harness share 0/20; no drift alarm (all clauses MET). Read a baseline red's first FAIL line before calling the base red (it was a missing `AILANG_BIN`).
- **Record:** world-mission-log.md entry 283; design_docs/verification/world-iter283/.
