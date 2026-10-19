# World mission dashboard — iteration 284 (2026-10-19)

- **This fire:** row 156 M3a + M3b built, PR #319. A declared `protected` set now reaches the sandbox: srt `denyWrite` gains the directory literal (both OSes) and `<wt>/**/G` (darwin only), appended after World's deny names and srt's defaults.
- **Compatibility:** absent / `[]` render the old settings bytes and digest; both existing goldens unchanged.
- **Real srt:** `TestExecSrtProtectedWrites`, 18 rows through the handler — darwin 18/18 locally; linux via CI bubblewrap. Deleting the append turns `/modify` red (`WORLD-PROBE-WROTE`, host hash changed).
- **Judge:** Agent sonnet (cross-vendor), own worktree: r1 PASS 91, zero blocking; NB1 (handler arm derived its expectation from the profile) measured and fixed test-only → r2 PASS 93.
- **Verification:** outside the sandbox, pinned `AILANG_BIN`: vet, focused tests, broker + daemon `-race` rc 0; commits rebuilt from snapshots, sha256-identical; PR CI 3/3 required before merge.
- **Release:** `v1.0.0-rc.1` pre-release. Proof/compiler pin 0.41.0; tool pin 0.54.15.
- **Clauses:** 7/7 MET; goal unmoved. Before v1.0.0 (D-WORLD-81): row 156 (10/21), row 150 (9/17), then row 151.
- **Next:** row 156 M4a (refuse configured git filters), M4b; M6-pre (V24 real-worktree timing) before M6a.
- **Parked on Mark:** none (ledger zero OPEN). D-WORLD-84/85 resolved attended; rows 145/189/198/199 filed attended.
- **Routing/cadence:** executor codex gpt-6.1-sol, evaluator Agent sonnet; designer and planner not needed (doc current, plan ready).
- **Cost:** metered $0.00; codex 114k tok; Agent subagents ≈113k tok.
- **Retro:** harness share 0/20; no drift alarm (all clauses MET). zsh does not split `$FILES` — use arrays in multi-file commands.
- **Record:** world-mission-log.md entry 284; design_docs/verification/world-iter284/.
