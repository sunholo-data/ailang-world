# AILANG World — mission dashboard

*Snapshot, overwritten every Gate 4. History lives in `world-mission.md` STATUS + `world-mission-log.md`.*

**As of** 2026-08-06, iteration 58 · dev @ `e3808c0` · CI **green BOTH jobs at HEAD** (SHA-addressed);
`ailang-code verify gate` was red on a **declared GitHub Actions outage** at the BG.A merge SHA and is
**green at HEAD** on `e3808c0`, a descendant carrying the same code — all 11 steps, step-log verified.

## In flight

- **Item 10 `w-boundary-gate-tree-mutation` — `BG.A` LANDED** (PR #47 → squash `278f102`, evaluator
  `sonnet` **89/100 r1, zero blocking**). The gate no longer writes the tree it guards: the mutant is
  **declared** via `go list -overlay` + an overlay-aware read; the `defer`-based restore is deleted.
  **AC2/AC3/AC4/AC5 discharged.**
- **`[NEXT]` is milestone `BG.B`** (`AC1a` · `M3`, `M6`), gated on nothing — **but apply the
  three-write-site correction first** (below) or its own AST guard reds on `BG.A`'s landed code.
- **Item 8 `w-self-mod-vertical`** — `SM.B2a` (~780 LOC, first irreversible-publish-capable code)
  queued behind item 10; unchanged.
- **Item 5 `w-mcp-projection` — still BLOCKED** on one prerequisite. Unchanged.

## Latest — a checker that cannot read the tree finds no forbidden imports

Met three times this iteration: inside the fix, inside the doc+plan that specified it, and inside my
own control.

- **`M5` (the SIGKILL harness) passed 4/4 arms with a firing negative control.** Marker awaited →
  artifacts verified → overlay verified to map target→mutant → process verified **alive** → `SIGKILL`
  (`rc=-9`): **0** changed shas, **0** porcelain lines. The **same kill on the base harness**:
  `RESIDUE=YES`, ` M host/store/store.go`. Outcomes differ, so the green measures the mechanism.
- **The plan's `BG.B` write-site count is wrong by one, and the missing one reds `BG.B`'s own guard.**
  Plan says two writes; measured **3** `os.WriteFile` (`:383` marker, `:428` mutant, `:439` overlay),
  **0** `OpenFile/Create/Rename`, KP `os.ReadFile` = 4. Route the AC4 marker through `confinedWrite`
  too — it is *required* to live outside `repoRoot`, which is exactly what the writer permits.
- **The doc and plan both specify a latent bug.** `go/parser` tests `src != nil` on the **interface**,
  so a typed nil `[]byte` parses as an **empty source** — every unreplaced file becomes
  `expected 'package', found 'EOF'`. Isolated in `parseSrc`.
- **My own known-positive control was invalid** (armed=0, unset=0 — an instrument that cannot see a
  positive). Cause: no `-v`, i.e. this sprint's own `V16c`, one iteration after it was measured.

## Loop · cost · asks

- launchd `mission-world`; controller `claude-opus-5`. Executor **`opus`** (env pin; the plan assumed
  codex, so its no-git-writes/snapshot rule did not apply). Evaluator **`sonnet`**. Designer/planner
  not fired; rotation pointer unchanged at `claude:claude-fable-5`. Verify profile `ailang-code`;
  AILANG pinned **v0.30.0**. Issue **#32**.
- **`metered=$0.00`** vs the $5 ceiling — every role on a quota bucket, no quorum round.
- **Parked on Mark: NONE.** `8/OD-2`, `10/OD-1`, `10/OD-2` open, all non-blocking with controller
  defaults recorded. FYI not blocking: item 9's human-gated half (pin CI job 1 vs track `latest`).
