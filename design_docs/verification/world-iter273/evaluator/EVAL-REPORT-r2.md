# EVAL REPORT r2 — AILANG World row 150 M3b (head 330bb0b, base 3e22bd7)

Verdict: PASS. Score 93/100. Blocking findings: NONE.

## Score
- Contract / allowlist (20): 20. `git diff --stat 3e22bd7 330bb0b` touches only demo_proof.go, demo_proof_test.go, demo_seed.go, demo_test.go plus design_docs/verification/world-iter273/**.
- Correctness / gates (25): 25. vet rc 0; gofmt clean; `-run ^TestDemo` rc 0 with both pins: named `--- PASS` for TestDemoSeedProofStage, TestDemoProofHTTPProducerRoundTrip, SourceHashRefusal, AdapterGuards, WrongKeyCannotProve, NoVerifiedSetRefuses, UncertainCommitStops, Bounds, KeyIORefusals, KeyDisclosure, TestDemoSeedRequiredIdentities; zero FAIL/SKIP. With WORLD_TOOL_AILANG_BIN unset: rc 0, only TestDemoSeedReadTestCheck and TestDemoSeedProofStage SKIP; all TestDemoProof* and TestDemoSeedRequiredIdentities PASS. host/evidence TestProducerEndToEndMintsEnvelopeValidatorAccepts PASS. `-race` on pin-free S4 tests ok.
- Test strength / non-vacuity (30): 27. All r1 survivors and plan mutants now killed (table). Deduction: two seed-level guards are enforced only by the tool-pin-gated stage test (see NB1).
- Design / hygiene (15): 14. Unused key write removed; seam is minimal.
- Evidence honesty (10): 7. Claims verified; ledger matches my reruns. Scored 7 for CI-coverage framing (NB1).

## Mutation table (all restored; git checkout)
| Mutant | Result | Asserting line |
|---|---|---|
| r1#1 drop servesEntry at seed call (hardcoded `[]string{"sameRef"}`) | KILLED (pin) | demo_test.go:970 seed-required-identities |
| r1#1b drop servesEntry in demoSeedRequiredIdentities | KILLED (pin-free via TestDemoSeedRequiredIdentities, demo_test.go:1090; also stage) | seed-required-identities |
| r1#2 seed uses newClient (no 10s bounds) | KILLED (pin) | demo_test.go:981 seed-s4-http-bounds, :999 seed-s4-request-context |
| r1#2b demoProofClient drops http.Timeout | KILLED (pin-free via TestDemoProofBounds demo_proof_test.go:834; also stage :981) | |
| r1#3 hex(key) appended to seed S4 error | KILLED (pin) | demo_test.go:1015 seed-key-disclosure: stderr contains hex-lower (safe-refusal, lost-response arms) |
| r1#3b hex(raw) in short-key load error | KILLED (pin-free) | demo_proof_test.go:949 TestDemoProofKeyDisclosure/31 |
| M3b-key (different validator key in positive arm) | KILLED | demo_proof_test.go:705 |
| M3b-null-verified (no-contract source in positive arm) | KILLED | demo_proof_test.go:731 ErrVacuousCheck |
| M3b-retry (second POST in PutObject) | KILLED | demo_proof_test.go:498 AdapterGuards/writer-transport-error-uncertain-no-retry, :796 UncertainCommitStops |
Survived: none. Masked: none.

## Controller skip edit
Acceptable. It matches the cliSeBins gating pattern; a configured wrong version still fails via requireBinaryRelease; unset is a skip with explicit reason. CI keeps pin-free coverage of the producer, adapters, bounds (client level), key IO refusals, key disclosure across raw/hex/base64 forms, and the seed's required-identity configuration. The pin-gated stage test remains the only gate on the seed's actual call wiring.

## Seam `demoSeedProve = demoProve`
No runtime change (same function, no wrapper). Tests in this package use no t.Parallel; the only reassignment is in TestDemoSeedProofStage arms, restored by defer, so no race (`-race` clean on the pin-free subset). Hazard is theoretical: a future parallel test touching it.

## Non-blocking findings
- NB1: Mutants "hardcode servesEntry-less list at the call", "seed builds client without bounds" and "leak key in the seed S4 error" are caught only by TestDemoSeedProofStage, which SKIPs where the tool pin is absent (CI). Local/pinned runs are the enforcement; the CI gate would not see a regression of seed wiring. Consider provisioning the pin in CI or a pin-free seed S4 harness.
- NB2: Restoration uses defer, not t.Cleanup; equivalent here (sequential arms), t.Cleanup would be more robust.
- NB3: TestDemoSeedReadTestCheck now also skips without the pin (pre-existing gating, not new).
- NB4: package-level mutable var in production code purely for test observation.
