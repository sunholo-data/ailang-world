# AILANG World — mission dashboard

*Snapshot, overwritten every iteration. History: `world-mission.md` (STATUS) + `world-mission-log.md`.
Last written: iteration 77, 2026-08-12.*

## Now

- **dev**: `40164ea` (PR #64 squash) — **CI GREEN BOTH JOBS on the merge commit**, SHA-addressed,
  `present=2` = expected 2, `unresolved_incidents=0`, so the green is attributable.
- **Item 12 `w-ail-gate-module-pin` is COMPLETE.** `scripts/verify_ail.sh` now pins the Leg-1 module
  **SET** by identity (`LEG1_MODULES`, 11 repo-relative paths), enumerating once and comparing
  **before** any `ai-check` runs, NUL-delimited end to end. Evaluator `sonnet` **93/100 PASS, zero
  blocking**. Totals **4/11/14 unmoved**; `verify_go.sh` rc=0, 34 `ok`, 0 `FAIL`, 2 healthy races.
- **The 11 is no longer decorative.** Five committed arms in `host/verifygate/module_manifest_gate_test.go`,
  each running a pristine control from its own isolated `t.TempDir()` root first, so an arm cannot
  pass vacuously. Deliberately **not** a count pin: add-one-delete-one prints a success line
  byte-identical to the baseline's (§S1 — never aggregate counts alone).
- **THE JUDGE DEFEATED THE GATE THE MILESTONE JUST LANDED, AND THE REPAIR IS IN-PR.**
  `find -name '*.ail'` is case-**SENSITIVE**, so `world/SNEAKY.AIL` never entered the swept set —
  the gate printed its own new success line, `✓ swept .ail module set equals the LEG1_MODULES
  allowlist (11 modules)`, with an unenumerated module sitting in `world/`. Reproduced first-party
  (control: `-name` 4, `-iname` 5), repaired with `-iname`, pinned by a fifth committed arm.

## Parked on Mark

**Nothing.** Zero open asks. Owed by the **shared driver** (frozen core — World cannot apply):
`ailang#611` (real per-role executor chain) and the World driver sync (missing `pi:*` pre-flight).

## Next

**Item 16** — the `host/broker` ~18% base flake, a standing tax on every mutation sweep. Then
**item 13** `w-evidence-grade-mapping` (cheapest high-leverage UI item; `PROVEN` is currently
unreachable). **Item 5 `P6.B` remains UNBLOCKED.** `SM.D` (item 8) is attended-only; 13/14/15 attended.

## Loop

- launchd, ~6h, headless. Issue **#53** (rotates Mondays 07:00 **local**; not due — created Mon 07:37 local, 15 comments).
- controller `opus` · planner `opus` (lane fail-closed, `missing-script`) · executor
  `codex:gpt-5.6-sol` · evaluator `sonnet`. `pi` **BARRED** from publish milestones.
- `metered=$0.00` — every lane a quota bucket. Cap $5.

## Standing hazards

- **GUARD THE HELPER, MISS THE BRANCH/CALL SITE — 8 instances in 5 milestones**, now FOUR directions:
  mechanism tested/SITES unguarded · site tested/second BRANCH unguarded · branch tested and the
  **shape space** of what it refuses never enumerated · and now **the shape space enumerated only in
  the spelling the author used** (every arm wrote the extension lowercase). Ask all four.
- **A RECOGNISER'S COVERAGE IS A PROPERTY OF ITS INPUT GRAMMAR.** Before trusting any set-compare,
  ask what its ENUMERATOR cannot see — case, symlinks, roots, extensions. `find -name` is
  case-sensitive; `-iname` is not.
- **AN EMPTY jq/`gh` RESULT IN A POLL IS A BROKEN INSTRUMENT UNTIL PROVEN OTHERWISE.** Iter-77's own
  Gate-3b poll printed `checks=0` for 10 minutes on a commit that was **2/2 green**: an unterminated
  jq string emitted nothing and `${1:-0}` rendered it as a confident zero. **And `set -- $out` does
  NOT word-split in zsh** — `$#` is 1. Always pair a poll with a known-positive SHA.
- **An isolated-tree test needs a PRISTINE-COPY control** — an incomplete copy reds for the wrong
  reason and looks exactly like a kill.
- **`git diff` OMITS UNTRACKED FILES** — use `git diff --no-index /dev/null <file>` (measured again
  this iteration: ordinary diff **0** added lines on the new test file, `--no-index` **263**).
- **A MUTATION CAN LAND BY sha AND STILL NOT REACH ITS BRANCH.** Iter-77 injected `>/dev/null`
  mid-command where a trailing `>&2` overrode it; the arm survived and was DISCARDED, not banked.
- **For a test-only milestone the compile gate is `go vet`, NOT `go build ./...`.**
- **`verify_go.sh` FATALs unless `GOTOOLCHAIN=go1.25.6`** (go1.26.4 miscompiles `host/store/scan.go`)
  — a BASE condition, not a regression.
- **`go test ./...` runs PACKAGES CONCURRENTLY** (`verify_go.sh:108`, no `-p 1`), and `host/boundary`
  enumerates *and reads* the live `world/` — so no committed test may mutate the live tree.
- **`rg` is NOT a binary here.** **`host/broker` is ~18% flaky at base** (item 16).
- **`rm -rf design_docs/verification` deletes THREE tracked sibling dirs.** Scope the path.
