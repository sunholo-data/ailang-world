# AILANG World — mission dashboard

*Snapshot, overwritten every iteration. History: `world-mission.md` (STATUS) + `world-mission-log.md`.
Last written: iteration 76, 2026-08-12.*

## Now

- **dev**: `d201a1e` — CI green both jobs on the parent `304120b`, SHA-addressed, `checks=2` = expected 2.
- **This iteration produced a DESIGN DOC, not a milestone**: item 12 `w-ail-gate-module-pin` had no
  doc, so the routing step was designer → quorum. `design_docs/planned/w-ail-gate-module-pin.md`
  (557 lines) is quorum-cleared and **ready for sprint-planner** — that is the next iteration's pick.
- **The gate defect is now measured at HEAD, not inherited.** Three arms, tree restored byte-identical
  each time: stray `world/*.ail` → `12 module(s)` **rc=0 PASSED**; delete a leaf sketch →
  `10 module(s)` **rc=0 PASSED**; **both composed → `11 module(s)`, success line BYTE-IDENTICAL to the
  baseline, PASSED**.
- **That third arm killed the queue row's own prescription.** The row asked for
  `EXACT_TOTAL_MODULES=11`; a count pin passes the add-one-delete-one mutant. The doc ports the
  identity allowlist the sibling leg already implements (`verify_world_package.sh:86-96`) and that
  coding-standards **S1** mandates ("never aggregate counts alone").

## Parked on Mark

**Nothing.** Zero open asks. Owed by the **shared driver** (frozen core — World cannot apply):
`ailang#611` (real per-role executor chain) and the World driver sync (missing `pi:*` pre-flight).

## Next

**Item 12 sprint** — plan + execute the doc above (~0.5d, top of band). Then **item 16** (the
`host/broker` ~18% base flake), a standing tax on every mutation sweep. **Item 5 `P6.B` is UNBLOCKED**
(prerequisite discharged at `625fb89`). `SM.D` (item 8) is attended-only; items 13/14/15 attended.

## Loop

- launchd, ~6h, headless. Issue **#53** (rotates Mondays 07:00 **local**; not due — created Mon 07:37 local, 14 comments).
- controller `opus` · designer `claude:claude-fable-5` (rotation) · planner `opus` · executor
  `codex:gpt-5.6-sol` · evaluator `sonnet`. `pi` **BARRED** from publish milestones.
- `metered=$0.169` (quorum R1 $0.0691 + R2 $0.0999); cap $5.

## Standing hazards

- **GUARD THE HELPER, MISS THE BRANCH/CALL SITE — 7 instances in 4 milestones**, three directions:
  mechanism tested/SITES unguarded · site tested/second BRANCH unguarded · branch tested and the
  **shape space of what it refuses** never enumerated. Ask all three.
- **A REVIEWER'S OBJECTION IS A CLAIM TOO — and can be right for the wrong reason.** Iter-76's
  surviving round-2 objection named two concrete exploits; **both were REFUTED** in isolated trees
  with firing controls, while the defect it pointed at was real. Measure before applying, and record
  the correction instead of laundering it into agreement.
- **An isolated-tree test needs a PRISTINE-COPY control** — an incomplete copy reds for the wrong
  reason and looks exactly like a kill. My own first newline probe created **1 file where it needed
  2**; only an asserted creation-count caught it, and its "PIN DEFEATED" line was discarded as vacuous.
- **`git diff` OMITS UNTRACKED FILES** — use `git diff --no-index /dev/null <file>`.
- **For a test-only milestone the compile gate is `go vet`, NOT `go build ./...`.**
- **`verify_go.sh` FATALs unless `GOTOOLCHAIN=go1.25.6`** (go1.26.4 miscompiles `host/store/scan.go`)
  — a BASE condition, not a regression. Confirmed again this iteration (local go is 1.26.4).
- **`go test ./...` runs PACKAGES CONCURRENTLY** (`verify_go.sh:108`, no `-p 1`), and `host/boundary`
  enumerates *and reads* the live `world/` (`allowlist_world_test.go:197`, `:293`) — so no test may
  mutate the live tree. `host/broker`'s AST gate does NOT collide (filters `.go` at `:149`).
- **`host/verifygate` has NO file-count pin** (`host/boundary`'s `wantFileCount = 1` is scoped to
  `host/boundary` only) — a new `_test.go` file there is safe.
- **`rg` is NOT a binary here.** **`host/broker` is ~18% flaky at base** (item 16).
- **`verify_ail.sh` never asserts the module count against 11** — only against 0. Item 12, designed.
