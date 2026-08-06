# AILANG World — mission dashboard

*Snapshot, overwritten every Gate 4. History lives in `world-mission.md` STATUS + `world-mission-log.md`.*

**As of** 2026-08-06, iteration 56 · dev @ `deeb804` + PR #46 · CI green both jobs (SHA-addressed)

## In flight

- **Item 10 `w-boundary-gate-tree-mutation` — DESIGN DOC LANDED** (PR #46). Promoted ahead of
  `SM.B2a` this iteration on a new measurement, not on the queue order.
- **`[NEXT]` is item 10's sprint-planner run**, gated on nothing. Then `SM.B2a`.
- **Item 8 `w-self-mod-vertical`** — `SM.B2a` (~780 LOC, brokered publish, the first
  irreversible-publish-capable code) still the next milestone; unchanged, deliberately not started.
- **Item 5 `w-mcp-projection` — still BLOCKED** on one prerequisite. Unchanged.

## Latest — the gate that guards the tree can poison it, and the build cannot see it

- **The teeth-proof writes live production sources and restores with a `defer`.** That `defer`
  survives return and `t.Fatal` and **nothing else**: SIGKILL mid-mutation → `rc=137`, residue
  permanent; Go's own `-test.timeout` panic → `rc=2`, same (60s control completes clean). Both kill
  paths are in the repo's own gate — `verify_go.sh` runs `-race -timeout 8m` inside an
  `os.killpg(SIGKILL)` at 600s.
- **The residue is invisible to the build**: `go build` **rc=0**, `go vet` **rc=0** with all three
  mutants applied. It reds *the boundary gate itself*, accusing an innocent file of a
  network-boundary violation — during the one sprint whose job is to add network code.
- **Correction to mission records**: the mutants are NOT "deliberately non-compiling" (charter,
  queue row and this dashboard all said so). They compile. The harm model is worse, not milder.
- **Fix**: `go list -overlay` — the overlay closure is diff-identical to a physically poisoned tree
  (control fired, 69-package difference) and the tree stays 0-dirty.

## Loop · cost · asks

- launchd `mission-world`; controller `claude-opus-5`. Designer **`claude:claude-fable-5`**
  (rotation slot 1; pointer advanced). Planner/executor/evaluator **not fired** — a design-doc
  iteration. Verify profile `ailang-code`; AILANG pinned **v0.30.0**. Issue **#32**.
- **`metered=$0.154`** vs the $5 ceiling — two quorum rounds, all four reviewer slots present.
- **Parked on Mark: NONE.** `8/OD-2`, `10/OD-1`, `10/OD-2` open, all non-blocking with controller
  defaults recorded. FYI not blocking: item 9's human-gated half (pin CI job 1 vs track `latest`).
