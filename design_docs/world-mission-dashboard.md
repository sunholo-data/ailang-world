# World mission dashboard — iteration 276 (2026-10-16)

- **This fire:** row 150 **M5a** LANDED `bcd2d7c` (PR #295). The validated decision packet now carries an internal evidence roster, read only from its carrying invocation record's input. The roster is exactly four kind/ref pairs in a strict shape, its reads are bounded under one shared 10 s deadline, and every refusal happens before any label. ATTESTED needs a consistent stored broker record. CLAIMED resolves only the demo-authored claim, never AiReview. TESTED, PROVEN and rendering are M5b.
- **Independent judge:** sonnet r1 **PASS 90**, zero blocking. All three plan mutants were killed, plus about 40 of the judge's own. The controller reproduced one kill with the tool pin unset.
- **CI:** PR 3/3 green on `a9c0771`. The merge's code tree is identical to that head.
- **Mid-fire attended landings:** row 114 LANDED, so **clause 5 is MET and all 7 clauses are MET** (D-WORLD-80 = A). **`v1.0.0-rc.1`** is tagged as a GitHub pre-release (D-WORLD-81 = A). D-WORLD-75 = A, D-WORLD-79 is closed moot, row 185 landed (PR #291), and rows 186–189 were filed.
- **Latest GitHub release:** `v1.0.0-rc.1` (pre-release, 2026-10-16). The proof/compiler pin is the released AILANG 0.41.0; the tool pin is 0.52.1.
- **Clauses:** **7 of 7 MET.** What remains before `v1.0.0` is D-WORLD-81's order: row 156, then finish row 150, then row 151.
- **Row 150:** 9 of 17 milestones done (M1a–M5a). Next on the row is M5b: TESTED/PROVEN resolution and rendering, which must also settle the plan's "source subject sibling" phrase.
- **Parked on Mark:** none. The decision ledger has 68 rows and zero OPEN.
- **Next:** row 156 (`w-exec-protected-test-paths`), design first. It is NEXT per D-WORLD-81.
- **Routing/cadence:** executor codex gpt-6.1-sol (recipe), judge Agent sonnet. No designer or planner was needed. 4h throttle.
- **Metered:** $0. No shared skill, harness or routing-policy edit.
- **Evidence:** log entry 276; `design_docs/verification/world-iter276/`; gate logs in `~/.ailang/state/world-iter276/`.
