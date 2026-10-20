# World mission dashboard — iteration 286 (2026-10-20)

- **This fire:** row 156 M4c built, PR #325. A profile declaring `protected` now runs startup probe arm 8 (`arm8-protected-write`): the sandbox must refuse the protected writes with host bytes unchanged (darwin also refuses creation), and on linux the base-name write is caught by `checkProtected` (intact=false).
- **Scope:** startup only; no per-run integrity record yet (M6). Absent/[] `protected` keeps the exact seven-arm probe.
- **Judge:** Agent sonnet (cross-vendor), own worktree: r1 PASS 84 with 11 own survivors (5 spec'd) → test-only fix `9edf116` → r2 PASS 91, zero blocking.
- **Verification:** outside the sandbox, pinned `AILANG_BIN`: vet, named tests, broker `-race` rc 0; real srt arm 8 PASS on darwin; controller mutants (unwired arm, no darwin creation) survived run 1 and are killed after a test-only fix. Linux arm 8 runs only in CI.
- **Release:** `v1.0.0-rc.1` pre-release. Proof/compiler pin 0.41.0; tool pin 0.54.15.
- **Clauses:** 7/7 MET; goal unmoved. Before v1.0.0 (D-WORLD-81): row 156 (13/21), row 150 (9/17), then row 151.
- **Next:** row 156 M5a (protected AILANG policies), then M5b; M6-pre (V24 real-worktree timing) before M6a.
- **Parked on Mark:** none (ledger zero OPEN).
- **Routing/cadence:** executor codex gpt-6.1-sol (3 runs), evaluator Agent sonnet (2 rounds); designer and planner not needed (doc current, plan ready).
- **Cost:** metered $0.00; codex ≈226k tok; Agent subagents ≈142k tok.
- **Retro:** harness share 0/20; no drift alarm. Executor directives should demand a production-entry-point test for any milestone that wires something, plus a negative control for every leg the design spells out.
- **Record:** world-mission-log.md entry 286; design_docs/verification/world-iter286/.
