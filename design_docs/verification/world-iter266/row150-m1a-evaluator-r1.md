# Row 150 (design r4 + plan + M1a) — independent evaluator, iteration 266

Agent tool, `sonnet`, fresh foreground context (56,640 tok, 7 tool calls, 72 s). Judged branch
`design/w-marketing-capture-r3` at `f59a76e` (PR #264) against the plan's M1a section and design D3.

**VERDICT: PASS 88/100, zero blocking.** Correctness 27/30 · tests/non-vacuity 22/25 · coding standards 17/20 ·
honesty 12/15 · scope 10/10.

Non-blocking:
1. Stale V15 row still said the repeat-GET probe was PENDING after V24 ran it — FIXED by the controller in `7dfc76f`.
2. V24/V25 were measured on a test-helper store, not the seeded demo store (the doc says so; M7 re-measures, fallback clause present).
3. The `ensures` clauses in `evidence.ail` largely restate the bodies; only the range bounds are independent. Z3 verifies
   them, and the Go drift test (not the contracts) guards the law — confirmed by the M1a-law mutant.
4. The drift test genuinely executes `world/types.gradeOf` under v0.41.0 (staged scratch dir, generated probe), no hand-copied table.
5. Empty/missing ai-check results cannot pass: empty results, skips, errors, counterexamples, `!Available` all refused; missing binary is t.Fatal.
6. All three quorum-r3 objections are answered in r4 (glm alt-1 premise refuted by V28; kimi V24 + SHAs; sonnet V25 + M7 fallback).

Mutation ledger (all compiled; tree clean after each): M1a-law killed (`RecordedEffect grade differs from executable kernel law:
package=1 kernel=2`); M1a-empty killed (`null report must be refused: empty verify results accepted`); own: authored-claim → 2
killed (inline test fails); own: rosterKindAllowed admits AiReview killed (inline test fails).

Gate ledger (judge, outside sandbox): pin v0.41.0; both focused tests PASS; check/test (31 passed)/ai-check (3 verified) rc 0;
`verify_ail.sh` rc 0 (16 identities, 40 named tests, 444 se-tools tests); `go vet ./host/coordinator` clean.
