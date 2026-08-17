# AILANG World — mission dashboard

*Snapshot, overwritten every iteration. History: charter STATUS + `world-mission-log.md`.*
**As of** 2026-08-17 (**attended session, Mark**) · all four ledger decisions RESOLVED ·
`mission_decisions.sh --open` → **0 rows** · **THE QUEUE IS UNBLOCKED.**

## ✅ The four rulings (full text: the decision ledger, which is the source of truth)

- **D-WORLD-5 → A.** Import upstream `serveapi` pinned at **v0.33.1**. Not a new direction — the
  doc's own Decision 2 says *"pin the first released/tagged upstream revision that contains this
  seam"*, and Decision 1 calls the local alternative "the forbidden reinvention". Revision round
  preconditions: run P6.A's frozen conformance fixture against v0.33.1; audit the dependency
  closure against `TestDaemonDependencyAllowlist` BEFORE any `go.mod` change. Surviving P6.B
  prerequisite: the VERIFIED commit-boundary contract (`world/*.ail` work).
- **D-WORLD-17 → A.** Seals bind to their minting validator: `Validator.Resolve(sealed)` as a
  METHOD, free `GradeOfValidated` dropped, tranche 1 library-only and explicitly non-production.
  The cross-validator refusal arm is REQUIRED as a named RED mutation — that arm is what makes
  "bind" non-vacuous. Self-minting is accepted (no library stops a caller lying to itself); the
  enforced property is that no caller makes SOMEONE ELSE'S validator resolve their seal.
- **D-WORLD-18 → A.** Ship the scoped item; unparks DIRECTLY to **sprint-planner** (doc §13 — no
  designer round). Residue enforced, not hidden: `TestNoNewDeadlineFreeStoreReads` pins
  approve 8 / registry 2 / replay 1, so follow-on progress is mechanically observable 11 → 0.
- **D-WORLD-DRIVER-1 → B, with teeth.** The driver stays FLEET-owned and the sync is a COMMIT,
  never working-tree dirt. The 2026-08-15 routing bundle lands as the first fleet-authored commit
  (attended). NEW: `verify_go.sh`'s **driver drift gate** reds while `tools/launchd/` or
  `mission_decisions.sh` diverges from HEAD (liveness control ≥5 tracked files). A drift red means
  "the fleet must commit" — never "absorb it"; CLAUDE.md now says so.

## The guardrail worked on its own exception

The harness permission layer blocked the controller from staging `tools/launchd/*` — frozen-core
enforced by a mechanism that cannot see an in-conversation ratification. Correct, and consistent
with the ruling: the FLEET'S HUMAN stages and commits driver files, so the bundle commit is
authored via Mark's own invocation in this session.

## Routable now (recommended order)

1. **Item 5** — revision round (designer): conformance fixture vs `serveapi` v0.33.1 + closure audit.
2. **Item 17** — revision round (designer): arm-A spec is in the queue row; add the refusal arm.
3. **Item 18** — sprint plan (sprint-planner; no designer needed).

**Routing constraint:** codex exhausted fleet-wide until **2026-08-20 05:34** — designer rotation
falls to the NEXT entry, never `$MODEL`. Gates on the final tree: `test_mission_routing.sh` 9/9 ·
ledger valid 5 rows / 0 open · `verify_ail.sh` + `verify_go.sh` green (see the two landing commits).
