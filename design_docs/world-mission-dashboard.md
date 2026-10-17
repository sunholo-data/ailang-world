# World mission dashboard — iteration 279 (2026-10-17)

- **This fire:** row 156 M1a (PR #300). The protected-path grammar core, pure Go with no wiring yet: parse and refuse per design §4.1, canonical form, digest, and a matcher that folds exactly like AILANG's fileguard (NFC plus Unicode case orbits, so `ſ` and the Kelvin sign fold too).
- **Judge findings fixed:** padded entries are now refused (fileguard trims them); the digest no longer depends on input order (one spelling per fold key, the byte-smallest); a `.git` segment is refused at any depth. Design §4.1 says so.
- **Carry-over closed:** the M0b judge's four untested validator guards now have self-test cases, and the contract test fails rather than skips without python3.
- **Judge:** independent sonnet r1 **PASS 88** → r2 **PASS 94**, zero blocking.
- **Latest GitHub release:** `v1.0.0-rc.1` (pre-release, 2026-10-16). Proof/compiler pin AILANG 0.41.0; tool pin 0.52.1.
- **Clauses:** **7 of 7 MET.** Before `v1.0.0` (D-WORLD-81): finish row 156 (3/19), then row 150 (9/17), then row 151.
- **Parked on Mark:** **D-WORLD-82** (row 156 design decisions D-156-1..4; defaults A, A, A, B; the plan follows the defaults).
- **Next:** row 156 M1b (render the four consumers, rebase module roots).
- **Routing/cadence:** executor Agent haiku (pin `claude:claude-haiku-5-5`; codex over daily ration), judge Agent sonnet; no designer/planner (ready-made milestone). 4h throttle.
- **Metered:** $0. No shared skill, harness or routing-policy edit.
- **Evidence:** log entry 279; role directives in `~/.ailang/state/world-iter279/`.
