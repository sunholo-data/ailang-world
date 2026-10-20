# World mission dashboard — iteration 288 (2026-10-20)

- **This fire:** row 156 M5b built, PR #328. Each episode's protected set (selected by binding, episode mapping, or the only profile; rebased to its module root) now reaches its own AILANG policy and run variants. The real v0.54.15 policy-tool refuses protected write/edit.
- **Scope:** policy half only. No exec, no usable profile, or no `protected` keeps today's bytes; a binding read error fails the AILANG family closed. No World symlink/hard-link pre-check yet (M5c/M5d).
- **Judge:** Agent sonnet (cross-vendor), own worktree: r1 PASS 88, survivor M10 (same-set build sharing untested) → test-only fix `c8ab0c3` → r2 PASS 91, zero blocking.
- **Verification:** outside the sandbox, pinned `AILANG_BIN`: vet, named tests, real-tool arm, broker and daemon `-race`, `./cmd/...` rc 0; controller survivor C1 (join ignoring the set) closed with a new test.
- **Release:** `v1.0.0-rc.1` pre-release. Proof/compiler pin 0.41.0; tool pin 0.54.15.
- **Clauses:** 7/7 MET; goal unmoved. Before v1.0.0 (D-WORLD-81): row 156 (15/21), row 150 (9/17), then row 151.
- **Next:** row 156 — re-measure V7 on v0.54.15, then keep or drop M5c/M5d; then M6-pre (V24 timing) before M6a.
- **Parked on Mark:** none (ledger zero OPEN).
- **Routing/cadence:** executor codex gpt-6.1-sol (1 run), evaluator Agent sonnet (2 rounds); designer and planner not needed (doc current, plan ready).
- **Cost:** metered $0.00; codex ≈112k tok; Agent subagents ≈140k tok.
- **Retro:** harness share 0/20; no drift alarm. When a milestone edits a cache or single-flight path, the directive should require both join outcomes as test rows.
- **Record:** world-mission-log.md entry 288; design_docs/verification/world-iter288/.
