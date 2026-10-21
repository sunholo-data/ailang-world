# World mission dashboard — iteration 293 (2026-10-21)

- **This fire:** row 156 M6c built, PR #340, test-only. Every tamper class the design names is now pinned through the production handler: creation (including in a new directory), deletion, rename, replacement, `.gitignore`/`info/exclude` ignored creations, assume-unchanged/skip-worktree index flags, and a moved HEAD. The episode-start snapshot exempts standing `.venv` files. Recorded protected results replay byte-for-byte without running anything. Real srt refuses a linked worktree's index/HEAD writes.
- **Controller catches:** the judge's survivor (`--untracked-files=normal`) came with a suggested fixture that did not kill it, because Git recurses into a new directory under a file-name pathspec in both modes. Re-aimed at a `tests/**` directory pattern, and now killed. A race-leg red under judge co-load was re-run isolated (green), then the full unloaded `verify_go.sh` gave rc 0.
- **Judge:** Agent sonnet, own worktree: r1 PASS 80 → fix → r2 **PASS 87**, zero blocking throughout.
- **Release:** `v1.0.0-rc.1` pre-release. Proof/compiler pin 0.41.0; tool pin 0.54.15.
- **Clauses:** 7/7 MET; goal unmoved. Before v1.0.0 (D-WORLD-81): row 156 (20/24, M5c/M5d held), row 150 (9/17), then row 151.
- **Next:** row 156 M7a (`why` shows intact / changed / unchecked / unprotected runs), then M7b (QUICKSTART guidance).
- **Carried:**
  - the parser's unmerged `u` branch is unpinned but fails closed;
  - `exec_protected_record_test.go` is 1 094 lines (split at M7);
  - "budget is spent" after a slow before-phase;
  - `/after-cap` needs `/usr/bin/perl`;
  - arm 8's probe budget under load (dogfood friction #22).
- **Parked on Mark:** D-WORLD-86: drop M5c/M5d (A, recommended) or build them anyway (B). Default: hold.
- **Routing/cadence:** executor codex gpt-6.1-sol, evaluator Agent sonnet. No designer or planner: the design and plan are ready.
- **Cost:** metered $0.00. codex executor 116,969 tok; Agent judge subagent ≈136k tok over 2 rounds.
- **Retro:** harness share 0/20; no drift alarm. Process fix: never run the judge's mutation drills while `verify_go.sh` runs. Also, re-run the mutant against a judge's suggested fixture before committing it (charter Repo Profile).
- **Record:** world-mission-log.md entry 293 (log rotated: 20 entries live).
