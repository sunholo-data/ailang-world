# AILANG World — mission dashboard

*Snapshot, overwritten every iteration. History: `world-mission.md` (STATUS) + `world-mission-log.md`.
Last written: iteration 74, 2026-08-12.*

## Now

- **dev**: `88eb850` — CI **green both jobs**, SHA-addressed, `checks=2` = expected 2, **0 incidents**.
- **Last landed**: **`TR.B2`** (PR #62) — **`TR.B` IS COMPLETE**. `host/transitionreg/bind.go`:
  `Bind` (3 ordered refusals, broker's own denial label verbatim), `Check` (all 5 authority pins),
  the confined `Bound.Request`, and the single-read `Request`/`Allowed()` fixture.
  **AC6 + AC7 activated** to exactly 3 each, tolerant arms deleted. Evaluator `sonnet` **96/100**,
  **zero blocking**.
- **The controller sweep found the one branch 21 executor arms could not see**: `equalRequirements`
  has TWO refusals (length, element-wise); only the length one was observed. Fixed as a subtest,
  proven by the inverse arm.

## Parked on Mark

**Nothing.** Zero open asks. Owed by the **shared driver** (frozen core — World cannot apply):
`ailang#611` (real per-role executor chain) and the World driver sync (missing `pi:*` pre-flight).

## Next

**`TR.C`** — the binding gate, and the LAST milestone of item 11. `TR.A`+`TR.B` deliver the
mechanism; without `TR.C` the undeclared-effect guard is an unenforced helper and item 5 `P6.B`'s
prerequisite is NOT satisfied. Then item 12, and item 16 (the `host/broker` ~18% base flake).
`SM.D` (item 8) is attended-only; items 13/14/15 (UI programme) were filed attended.

## Loop

- launchd, ~6h, headless. Issue **#53** (rotates Mondays 07:00 **local**).
- controller `opus` · designer rotation (last `codex:gpt-5.6-sol`) · executor `codex:gpt-5.6-sol` ·
  evaluator `sonnet`. No planner this iteration — TR.B's plan already scoped T4–T7b.
  `pi` **BARRED** from publish milestones.
- `derive-planner-lane.sh` absent → lane fails closed to opus, loudly. `metered=$0.00`; cap $5.

## Standing hazards

- **GUARD THE HELPER, MISS THE BRANCH/CALL SITE** — **6 instances in 3 milestones** (3 in `TR.A2`,
  2 in `TR.B1`, 1 in `TR.B2`), and TR.B2's is the **mirror** of TR.B1's: there the mechanism was
  tested and the SITES were not; here the site was tested and the mechanism's second BRANCH was not.
  Instrument that covers both: per helper/mechanism, ask **how many ways can this refuse, and how
  many does a test observe?** — never how many tests name it.
- **`git diff` OMITS UNTRACKED FILES**, so a rule-3j cut over a sprint that ADDS a file returns `0`
  — a broken instrument wearing a clean result. Use `git diff --no-index /dev/null <file>`.
- **`verify_go.sh` FATALs unless `GOTOOLCHAIN=go1.25.6` is exported** (go1.26.4 miscompiles
  `host/store/scan.go`). That rc=1 is a BASE condition, not a regression.
- **`rg` is NOT a binary here** — a harness-injected shell function, absent in CI. Use `grep`.
- **A refusal test asserting only *that* an error occurred pins no branch.** Pin the measured message.
- **A green `go test` is not a green `go vet`** — `copylocks` is outside both, invisible to CI.
- **`host/broker` is ~18% flaky at base** (`TestHandlerTimeoutKillsTheWholeProcessGroup`, 2/11) and
  **100% red without `AILANG_BIN`** — both fake mutation kills and falsify inverse arms. Item 16.
- **`verify_ail.sh` never asserts the module count against 11** — only against 0. Item 12.
