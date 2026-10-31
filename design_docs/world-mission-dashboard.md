# World mission dashboard — iteration 328 (2026-10-31)

- **This fire:** row 231 built, the first pick under **D-WORLD-102 = A**. The P0/P3 position predicate and the phone predicate are now CI-gated `marketingcheck position` / `marketingcheck phone`.
  - The banked programs are a parity oracle that CI builds and compares against.
  - Fixtures are generated: 34 JSON files and 10 phone skeletons.
  - A bounded subprocess runner; exit codes carry the verdict.
  - PR #440; judge sonnet PASS 90, zero blocking.
- **Dev:** green. Row 231 LANDED: #440 merged as `6641df63`, merge check-runs 3/3 success on attempt 1.
- **Clauses:** 1–7 MET; goal unmoved. Release order (D-WORLD-81): 150 → 151 → v1.0.0.
- **Row 150:** 23/23 milestones built; waits on the **attended A1**. P0/P3 now run `marketingcheck position`/`phone` (orchestration §4.4). They need the host-local phone corpus at `~/.ailang/state/world-iter328-phone-corpus/`; if it is missing, that is BLOCKED.
- **Row 151:** needs row 188 N2d and N4a, after the attended rows 219 → 217 → 218.
- **Loop's next rows (D-WORLD-102 = A):** 230 → 233, stopping when a release row is routable. Rows 235, 237 and new 238–240 (predicate tolerances, still-probe commit, checker tolerance pins) are hygiene outside that scope.
- **Roles this fire:** all foreground Agent-tool spawns. Designer opus: quorum r1 and r2 BLOCKED, then the carve-out. Planner opus: 7 milestones, 36/36 mutants. Executor haiku. Judge sonnet.
- **Parked on Mark:** no open decision (0 OPEN ledger rows). **Attended actions:** schedule row 150's A1; the 219 → 217 → 218 track unblocks 151.
- **Watch:** `host/broker` `TestExecProtectedCheckTimeout/reserve-after-before` reds locally under load (3 s bound); CI is the arbiter.
- **Cost:** $0.4042 metered (two quorum rounds); everything else on Anthropic quota.
- **Record and evidence:** world-mission-log.md entry 328; verification/world-iter328/.
