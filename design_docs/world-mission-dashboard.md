# Mission Dashboard — Ailang World

*Snapshot, overwritten every iteration. History: `world-mission.md` (STATUS) + `world-mission-log.md`.
Last written: iteration 110, 2026-08-22.*

**Where we are** — `dev` == `origin/dev` @ `3e0c34c`, CI green (`checks=2`). Last landed: item 17
`w-validated-proven-evidence-boundary`, COMPLETE (iter-108, PR #82 → `189299b`).

**In flight** — queue item **14** `w-workbench-read-only`: doc carve-out-complete (iter-109),
**sprint plan landed this iteration** (`WB.A`–`WB.K`, ≈1.5 d, 32 mutations, AC1–AC8). Not started.

**Next** — sprint-executor on **`WB.A`** (`codex:gpt-5.6-sol` → `opus`), worktree a sibling of the
repo, never `/tmp`. Then row 32 `w-wallclock-ceilings-not-derived` (reds the Go gate on commits that
cannot have caused it, so it costs a headless loop whole iterations) → item 22
`w-daemon-lock-wait-not-deadline-bound` → row 31 `w-stale-planned-doc-citations`.

**Blocked** — row 5 on upstream `sunholo-data/ailang#764` (predicate re-run this iteration: `OPEN`,
0 comments, unchanged since 08-17 — not a human ask). Row 8 `SM.D` is **attended-only** by ratified
design, never headless. Rows 6/7 wait on internal items.

**Routing** — controller `opus` · planner `opus` (lane `fail-closed:env-pin`) · executor
`codex:gpt-5.6-sol` → `opus` · evaluator `sonnet` (generator ≠ judge) · designer rotation **unspent**.
pi/deepseek lane SUSPENDED (`D-WORLD-20`). Gates need **both** `AILANG_BIN=/tmp/ailang-v0300/ailang`
and `GOTOOLCHAIN=go1.25.6`.

**Parked on Mark** — nothing. Ledger: 11 rows, **0 OPEN**. Zero open asks, 4th iteration running.

**Cost** — this iteration `metered=$0.00` of $5 · quota `opus` ×2 · Fable unspent 2nd running.

**The one thing to know this week.** `AC7` — *"only priced files changed"* — was **vacuous by
construction**, and not because the design was wrong: `git diff` cannot see untracked files, and the
sandboxed executor is forbidden from committing, so every file a sprint adds stays untracked for the
sprint's whole life. Two quorum rounds and a restored third reviewer could not have caught it; the
defect lives in the seam between a correct criterion and the controller's own recipe. Wider and
unswept: **`go test -run` exits 0 on an empty match set**, so every AC here shaped *"`go test -run
'TestA|TestB'` passes"* is green before either test exists.
