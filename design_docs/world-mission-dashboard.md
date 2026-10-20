# World mission dashboard — iteration 287 (2026-10-20)

- **This fire:** row 156 M5a built, PR #326. A protected set now reaches every AILANG policy World renders (base, Declassify, Net, Env), appended after the built-in denies, with safe TOML escaping; absent/[] keeps the old bytes exactly.
- **Scope:** policy rendering and handler plumbing only; the daemon still passes nil until M5b wires the episode profile. No World symlink/hard-link pre-check yet (M5c/M5d).
- **Judge:** Agent sonnet (cross-vendor), own worktree: r1 PASS 88, one survivor (raw DEL, because the test decoded with `encoding/json`) → test-only fix `3852dbc` → r2 PASS 91, zero blocking.
- **Verification:** outside the sandbox, pinned `AILANG_BIN`: vet, named tests, broker and daemon `-race`, `./cmd/...` rc 0; 5 controller mutants killed.
- **Release:** `v1.0.0-rc.1` pre-release. Proof/compiler pin 0.41.0; tool pin 0.54.15.
- **Clauses:** 7/7 MET; goal unmoved. Before v1.0.0 (D-WORLD-81): row 156 (14/21), row 150 (9/17), then row 151.
- **Next:** row 156 M5b (episode profile selection, module-root rebasing), then M5c/M5d (re-measure V7 on v0.54.15 first); M6-pre (V24 timing) before M6a.
- **Parked on Mark:** none (ledger zero OPEN).
- **Routing/cadence:** executor codex gpt-6.1-sol (1 run), evaluator Agent sonnet (2 rounds); designer and planner not needed (doc current, plan ready).
- **Cost:** metered $0.00; codex ≈96k tok; Agent subagents ≈138k tok.
- **Retro:** harness share 0/20; no drift alarm. A test that round-trips through a laxer stand-in decoder must also assert bytes for every character the real format forbids.
- **Record:** world-mission-log.md entry 287; design_docs/verification/world-iter287/.
