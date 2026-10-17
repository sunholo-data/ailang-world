# World mission dashboard — iteration 280 (2026-10-17)

- **This fire:** row 156 M1b (PR #307). The protected set now renders for all four consumers, as pure Go with no wiring yet: AILANG `fs_deny_write` re-based to the module-root sandbox, srt `denyWrite` for darwin (dirs plus `**/G`) and linux (dirs only), and root-anchored git pathspecs (`top,glob,icase`, bare `D` and `D/**`).
- **Measured against the real engine:** a pinned test runs AILANG v0.52.1's policy-tool over the V7 fixture table, empty-deny controls and two module-root rows (23 PASS); it skips without the pin.
- **Judge:** independent sonnet r1 **PASS 80** → r2 **PASS 84**, zero blocking. Its two surviving rebase mutants are now killed by three test rows. The diff is ~393 lines against the plan's 250 cap (recorded).
- **Latest GitHub release:** `v1.0.0-rc.1` (pre-release, 2026-10-16). Proof/compiler pin AILANG 0.41.0; tool pin 0.52.1.
- **Clauses:** **7 of 7 MET.** Before `v1.0.0` (D-WORLD-81): finish row 156 (4/19), then row 150 (9/17), then row 151.
- **Parked on Mark:** none open (ledger 70 rows, zero OPEN).
- **Next:** row 156 M1c (real-git pathspec equivalence, incl. the r3 rows).
- **Routing/cadence:** executor Agent haiku (pin `claude:claude-haiku-5-5`; codex over daily ration), judge Agent sonnet; no designer/planner (ready-made milestone). 4h throttle.
- **Metered:** $0. No shared skill, harness or routing-policy edit.
- **Evidence:** log entry 280; judge reports in `~/.ailang/state/world-iter280/`.
