# AILANG World — mission dashboard

*Snapshot, overwritten every Gate 4. History: `world-mission.md` STATUS + `world-mission-log.md`.*

**As of** 2026-08-10, **iteration 68** · dev @ `9789b87` · CI green both jobs · issue **#53**.

## In flight

- **`VL.A` LANDED** (PR #56 → `9789b87`) — item 9's three headless pieces all shipped. The `.ail`
  gate now **refuses a dev build**; the proving shim is **committed**; the always-firing warning is
  gone. **Item 9 is COMPLETE** — its doc-less charter row is the spec and every piece is closed.
- **`[NEXT]`: item 5 `w-mcp-projection` needs its ONE prereq (the transition registry) written, or
  item 6b. Item 8 has no headless milestone left (`SM.D` is attended-only).**
- **Parked on Mark: nothing blocking. One scope ask below.**

## Latest — the defect, and two greens that answered the wrong question

**The defect, measured at base:** a shim reporting `AILANG v0.33.0-105-…-dirty` drove the real gate
to `rc=0` / `verify gate PASSED`. The primary `.ail` gate passed exactly the build CLAUDE.md forbids.

**Landed:** five refusal branches each with a stable reason **CODE** (they *funnel* —
`NOT_A_RELEASE` is a catch-all, so an `rc`-keyed table scores **3 of 8** mutations falsely SURVIVED).
The release shape admits a pre-release identifier by design: upstream published `v0.24.1-rc1` with
`isPrerelease: false`, so strict `^vX.Y.Z$` had a measured 1-in-64 chance of redding CI on a release.

**Two local greens were both wrong, in the same way.** (1) The executor's equality-shaped notice was
quiet in CI and printed `moved from 'v0.33.0' to 'v0.30.0'` on **every local run** — always-firing
again, and false. Fixed to membership over the two lanes' releases. (2) Both gates were rc=0 locally,
re-run outside the codex sandbox on purpose — and CI red **twice**: a hardcoded rig path, then the
go-verify job's missing Z3. **Spine: a green proves the tree passes where you ran it; only CI proves
it passes where it must — and re-running "outside the sandbox" answers a different question.**

## Loop · cost · asks

- Planner `opus` · executor `codex:gpt-5.6-sol` · repair+adjudication `opus` · judge `sonnet`
  **38/100 FAIL r1 — the judge was RIGHT and its blocking finding is the iteration's spine**;
  `metered=$0.00`. 10 mutations, 10 killed. **v0.30.0** pinned.
- **ASK `9/OD-11`** (one word): may a milestone add a **Z3 install step to CI job 2**? Blocked
  headless because `9/OD-10` clause (a) scopes item 9 to zero `.github/` edits. Cost of not doing
  it: `host/verifygate` asserts the version block's contract, never `verify gate PASSED`.
