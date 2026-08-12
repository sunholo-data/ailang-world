# AILANG World — mission dashboard

*Snapshot, overwritten every iteration. History: `world-mission.md` (STATUS) + `world-mission-log.md`.
Last written: iteration 75, 2026-08-12.*

## Now

- **dev**: `625fb89` — CI **green both jobs**, SHA-addressed, `checks=2` = expected 2, **0 incidents**.
- **Last landed**: **`TR.C`** (PR #63) — **item 11 `w-transition-registry` IS COMPLETE**. A
  repository-wide `go/parser` gate: outside `host/broker` no `Invoke` selector, neither exported
  session constructor, and no `broker.Session` exposure (plain, aliased or dot import); inside, the
  three legacy sites pinned by **identity and exact count**. **AC11 activated 1 → 2.** Zero
  production LOC.
- **The judge defeated the gate with ordinary Go, and that is the iteration.** Every detector sat
  inside `case *ast.CallExpr`, so a **method value** (`call := s.Invoke`) or **function value**
  (`mk := broker.NewSession`) reached a raw session with the gate green — no reflection, no
  linkname, no build tags. Reproduced first-party (`walked=40`, gate rc=0), **fixed in-PR**, and
  proven in three arms. All **32** mutation arms had used the CALLED form.

## Parked on Mark

**Nothing.** Zero open asks. Owed by the **shared driver** (frozen core — World cannot apply):
`ailang#611` (real per-role executor chain) and the World driver sync (missing `pi:*` pre-flight).

## Next

**Item 12 `w-ail-gate-module-pin`** — move the module-count pin into `verify_ail.sh` itself
(~0.5d, needs a small design doc). Then **item 16** (the `host/broker` ~18% base flake), which is a
standing tax on every future mutation sweep. **Item 5 `w-mcp-projection` P6.B is now UNBLOCKED** —
its sole prerequisite was TR.A+TR.B merged **and TR.C green**, which is satisfied as of `625fb89`.
`SM.D` (item 8) is attended-only; items 13/14/15 (UI programme) were filed attended.

## Loop

- launchd, ~6h, headless. Issue **#53** (rotates Mondays 07:00 **local**).
- controller `opus` · planner `opus` · executor `codex:gpt-5.6-sol` · evaluator `sonnet`.
  `pi` **BARRED** from publish milestones.
- `derive-planner-lane.sh` absent → lane fails closed to opus, loudly. `metered=$0.00`; cap $5.

## Standing hazards

- **GUARD THE HELPER, MISS THE BRANCH/CALL SITE — now 7 instances in 4 milestones**, and TR.C adds a
  THIRD direction. TR.B1: mechanism tested, SITES unguarded. TR.B2: site tested, second BRANCH
  unguarded. TR.C: branch tested, and the **shape space of what it refuses** never enumerated — 32
  arms all spelled the thing the same way. Ask all three: how many ways can this refuse · how many
  call sites does it have · **how many ways can the thing it refuses be SPELLED**.
- **`git diff` OMITS UNTRACKED FILES**, so a rule-3j cut over a sprint that ADDS a file returns `0`
  — a broken instrument wearing a clean result. Use `git diff --no-index /dev/null <file>`.
- **For a test-only milestone the compile gate is `go vet`, NOT `go build ./...`** — `go build` does
  not compile `_test.go` at all, so a test-file mutant "builds" trivially.
- **`verify_go.sh` FATALs unless `GOTOOLCHAIN=go1.25.6` is exported** (go1.26.4 miscompiles
  `host/store/scan.go`). That rc=1 is a BASE condition, not a regression.
- **`rg` is NOT a binary here** — a harness-injected shell function, absent in CI. Use `grep`.
- **A refusal test asserting only *that* an error occurred pins no branch.** Pin the measured message.
- **A green `go test` is not a green `go vet`** — `copylocks` is outside both, invisible to CI.
- **`host/broker` is ~18% flaky at base** (`TestHandlerTimeoutKillsTheWholeProcessGroup`, 2/11) and
  **100% red without `AILANG_BIN`** — both fake mutation kills and falsify inverse arms. Item 16.
- **`verify_ail.sh` never asserts the module count against 11** — only against 0. Item 12.
