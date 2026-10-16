# World mission dashboard — iteration 275 (2026-10-16)

- **This fire:** row 150 **M4b** LANDED `acb3f54` (PR #289). Demo seed stage S6 runs the real CLI why chain for the read, test, check and approval entries, and refuses unless coordinator, ok and every named link hold. It requires exactly one packet anchored to S5's entry, then assembles the public manifest in memory, checks every role and ref over REST, and scrubs secrets, paths and timestamps. The seed still refuses success and writes no manifest.
- **Independent judge:** sonnet r1 **FAIL 66**: the manifest's contents were never asserted, and the controller reproduced it. After a codex fix round, r2 **PASS 86**, zero blocking, with all five r1 survivors killed with the tool pin unset.
- **CI:** one red, and it was real. A table of about 120 test controls shared one 30 s deadline, which ran out under `-race` on the 2-core runner but never on the Mac. The controller fixed it with per-subtest deadlines, and a 10 s drill confirms it (49 failures before, 0 after). Judge r3 reviewed the fix: **PASS 86**. Rebased over attended row 184 (`why.go` changed) and re-tested: 32/32.
- **Latest GitHub release:** none; World 1.0 is unreleased. The proof/compiler pin is the released AILANG 0.41.0; the tool pin is 0.52.1.
- **Clauses:** **6 of 7 MET.** Clause 5 is UNMET. Row 139 (routine AILANG upgrades) is attended (M1 and M2 landed, M3 next). Row 183 (real-project trials) is IN-SPRINT ATTENDED. Row 114 waits on 183. Row 184 landed attended this afternoon, which opens row 185.
- **Row 150:** 8 of 17 milestones done (M1a, M1b, M2a, M2b, M3a, M3b, M4a, M4b). Next is M5a: bounded anchored roster, ATTESTED and authored CLAIMED.
- **Parked on Mark:**
  - **D-WORLD-79:** looks moot. Row 181's fix landed attended (`38e7859`, race leg in its own job), and this PR's race job was green on 3 of 4 heads. Close it if you agree.
  - **D-WORLD-75:** unchanged (default A).
- **Next:** row 185 (`w-episode-requires-git`) if no attended session holds it, since it was ordered "NEXT after 184"; otherwise row 150 M5a.
