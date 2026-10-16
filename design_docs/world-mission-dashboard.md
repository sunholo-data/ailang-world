# World mission dashboard — iteration 273 (2026-10-16)

- **This fire:** row 150 **M3b** LANDED `c08729b` (PR #277). Demo seed stage S4 now loads the operator proof key and mints the logepoch proof through the production HTTP producer. The proof needs sameRef and servesEntry. Request bounds are 10 s, with 90 s for the whole stage. The seed still refuses completion and writes no manifest. An uncertain commit makes exactly one write and is never retried.
- **Independent judge:** sonnet. r1 **PASS 88**, zero blocking, 3 of 23 mutants survived. A fix round followed, plus a controller CI fix. r2 **PASS 93**, every survivor killed.
- **Found on the way:** the new stage test hard-failed in CI without the v0.52.1 tool pin, which CI never sets. That is the row 160 class. It now skips there like its siblings, so three seed-wiring guards are rig-only until row 160 lands.
- **Latest GitHub release:** none; World 1.0 is unreleased. The proof/compiler pin is the released AILANG 0.41.0; the tool pin is 0.52.1.
- **Clauses:** **6 of 7 MET.** Clause 4 is MET by attended ruling D-WORLD-77 (row 93 LANDED), which landed during this fire. Clause 5 is UNMET (as of charter `c9ff55e`). Row 114 Q1 ran attended, and both arms came back UNVERIFIED. Row 182 (outcome capture) has since LANDED attended. Row 183 (real-project trials) is IN-SPRINT ATTENDED and will supply 114's next questions. Row 139 (routine AILANG upgrades) is NEXT as an attended build.
- **Row 150:** 6 of 17 milestones done (M1a, M1b, M2a, M2b, M3a, M3b). Next is M4a: real approval transition, packet anchoring and budget delta.
- **Harness:** row 178 (worktree helper pending marker) RESOLVED by the fleet (`51f9ecb`).
- **Parked on Mark:**
  - **D-WORLD-75:** unchanged (default A).
- **Next:** re-read the 139 and 183 tags first; take 139 if it is loop-routable. Otherwise row 150 M4a.
- **Attended:** rows 139 and 183; production descriptor publish, session mint and release capture.
- **Routing/cadence:** executor codex gpt-6.1-sol (recipe; codex back in ration), judge Agent sonnet. No designer or planner was needed. 4h throttle.
- **Metered:** $0. No shared skill, harness or routing-policy edit.
- **Evidence:** log entry 273; `design_docs/verification/world-iter273/`.
