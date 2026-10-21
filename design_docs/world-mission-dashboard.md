# World mission dashboard — iteration 292 (2026-10-21)

- **This fire:** row 156 M6b built, PR #339. A protected exec result now carries the `protected` block: the episode baseline plus a before- and an after-command integrity verdict. The child's budget reserves `protectedCheckTimeout`, so the after-check always fits the call deadline. The after-check is detached from cancellation but never outlives the original deadline. A failed baseline skips both checks; a check error or a CHANGED verdict never fails the call (D-156-1). Unprotected profiles emit no key and run no new Git.
- **Controller catches:** two wiring mutants survived the executor's suite (budget read before the before-phase; a cancellable after-phase) and are now killed. `verify_go.sh` was red twice where the plan's gates could not see: an unpinned `WithoutCancel` root (`TestProductionContextRoots`, the second iteration running), and a load-fragile `/after-cap` row (now anchored to an absolute deadline target). `verify_go.sh` rc 0 on the fixed tree, race leg included.
- **Judge:** Agent sonnet, own worktree: r1 PASS 80 → fixes → r2 **PASS 90**, zero blocking throughout. It found the stale-`now` survivor on its own.
- **Release:** `v1.0.0-rc.1` pre-release. Proof/compiler pin 0.41.0; tool pin 0.54.15.
- **Clauses:** 7/7 MET; goal unmoved. Before v1.0.0 (D-WORLD-81): row 156 (19/24, M5c/M5d held), row 150 (9/17), then row 151.
- **Next:** row 156 M6c (tamper, ignored-file, flag and moved-HEAD mutants; pin replay), then M7a/M7b.
- **Carried:** a budget spent by a slow before-phase still returns "budget is spent" rather than a report-only result; `/after-cap` uses `/usr/bin/perl`; arm 8's startup probe budget under load (dogfood friction #22).
- **Parked on Mark:** D-WORLD-86: drop M5c/M5d (A, recommended) or build them anyway (B). Default: hold.
- **Routing/cadence:** executor codex gpt-6.1-sol, evaluator Agent sonnet. No designer or planner: design and plan ready.
- **Cost:** metered $0.00. codex executor 112,329 tok; Agent judge subagent ≈160k tok over 2 rounds.
- **Retro:** harness share 0/20; no drift alarm. Process fix: a directive whose diff adds a production context root names `TestProductionContextRoots` as a gate (charter Repo Profile).
- **Record:** world-mission-log.md entry 292 (log rotated: 20 entries live).
