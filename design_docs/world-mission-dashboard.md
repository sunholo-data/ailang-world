# World mission dashboard — iteration 291 (2026-10-20)

- **This fire:** row 156 M6a built, PR #338. A protected episode's exec handler captures its baseline once, at construction. That covers the committed protected-tree digest, the base commit and the fold-keyed ignored snapshot, under one 4 848 ms bound. A failure is kept on the handler and never fails construction. The bounded `protected` record shape exists but is not emitted until M6b.
- **Bounds:** design V24's 3 751 ms phase / 4 848 ms construction are applied and pinned by a test.
- **Controller catches:** one existing timeout test was sized for the old 2 500 ms bound and turned marginal at 3 751 ms; it is now derived from the bound. `verify_go.sh` caught an unpinned `context.Background()` root (`TestProductionContextRoots`) that the plan's gate list could not see.
- **Judge:** Agent sonnet, own worktree: r1 PASS 79 → fixes → r2 PASS 88 → one-bound fix → r3 **PASS 90**, zero blocking throughout.
- **Release:** `v1.0.0-rc.1` pre-release. Proof/compiler pin 0.41.0; tool pin 0.54.15.
- **Clauses:** 7/7 MET; goal unmoved. Before v1.0.0 (D-WORLD-81): row 156 (18/24, M5c/M5d held), row 150 (9/17), then row 151.
- **Next:** row 156 M6b (attach before/after integrity to exec results, reserving the check bound from the child budget), then M6c.
- **Carried:** row 156's arm-8 startup probe times out at 8 s under a loaded full suite (dogfood friction #22). Candidate: size its budget for load.
- **Parked on Mark:** D-WORLD-86: drop M5c/M5d (A, recommended) or build them anyway (B). Default: hold.
- **Routing/cadence:** executor codex gpt-6.1-sol, evaluator Agent sonnet. No designer or planner: design and plan ready.
- **Cost:** metered $0.00. codex executor 121,065 tok; Agent judge subagent ≈256k+ tok over 3 rounds.
- **Retro:** harness share 0/20; no drift alarm. When a bound constant changes, re-size every test stimulus sized against it. Run `verify_go.sh`, not just the plan's named gates.
- **Record:** world-mission-log.md entry 291.
