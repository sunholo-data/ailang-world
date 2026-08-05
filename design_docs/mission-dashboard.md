# AILANG World — mission dashboard

*Snapshot, overwritten every Gate 4. History lives in `world-mission.md` STATUS + `world-mission-log.md`.*

**As of** 2026-08-05, iteration 52 · dev @ `c0ca1df` · CI green both jobs (SHA-addressed)

## In flight

- **Item 8 `w-self-mod-vertical` — SPRINT-PLANNED** (plan + handoff in `.ailang/state/sprints/`,
  gitignored as all 18 prior sprint artifacts are; planner `opus`). **Re-scoped 4 → 6 milestones**:
  `SM.A · SM.B1 · SM.B2a · SM.B2b · SM.C · SM.D`. Stays **one queue item** — the split is internal.
- **`[NEXT]` is `SM.A`** — package projection, drift/export/tar gate, smoke, boundary guard.
  Gated on nothing, no kernel touch, ~620 LOC. It builds the ready packet `8/OD-1` authorizes.
- **Item 5 `w-mcp-projection` — still BLOCKED** on one prerequisite (transition registry absent at
  HEAD, measured iter-50, control fired). Unchanged this iteration.

## Latest

- **Mark approved the `world/` publish** (`#32`, 08:25). `8/OD-1` **RATIFIED as policy** — but not
  the exact-bytes stamp SM.D describes, and it cannot be: the packet doesn't exist until SM.A
  builds it. **An authorization is not an attendance.** SM.D stays attended-only, never in CI.
- **The planner refuted the design's central reuse claim.** Decision 3's "extract v0.30.0
  package-hashing logic" is **impossible** — it lives in upstream's `internal/pkg/` and World is a
  different module. `AC6` needs a re-implementation (`host/pkgproj`) + a 24-char cross-check.
  The doc cited that path three times *as evidence for* the plan the path forbids.
- **`DD-3`**: bumping the schema version makes `store.go:354`'s bare `return nil` reachable — a v1
  store would open fine and **never run `schemaSQL`**. Answered from ratified text, non-blocking.
- **DDL blast radius ~3× the doc's Conflict Surface** (second fixture in `journal_test.go`,
  `frozenFutureSchemaVersion = 2` collision, literal `PRAGMA user_version = 1`) — all into SM.B1's
  single commit. **5 ACs judged vacuous and replaced**; 36 mutations total.

## Loop · cost · asks

- launchd `mission-world`; controller `claude-opus-5`. Planner **`opus`** (lane
  `opus fail-closed:env-pin`, verbatim). Designer/executor/evaluator **not fired** — a planning
  iteration has nothing to execute or judge; designer rotation stays `codex:gpt-5.6-sol`.
- Verify profile `ailang-code`; AILANG pinned **v0.30.0** at `/tmp/ailang-v0300/ailang`; upstream
  read only at `e37b370d…`. Bookkeeping issue **#32** (week of 2026-08-03).
- **`metered=$0.00`** vs the $5 ceiling — all roles on quota buckets. First planning iter at zero.
- **Parked on Mark: NONE.** `8/OD-1` answered today; `8/OD-3` from ratified charter text; `8/OD-2`
  (upstream namespace auth) open but **non-blocking by design**. Next free OD: **`OD-9`**.
