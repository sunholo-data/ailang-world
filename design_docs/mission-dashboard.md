# AILANG World — mission dashboard

*Snapshot, overwritten every iteration. History: `world-mission.md` (STATUS) + `world-mission-log.md`.
Last written: iteration 78, 2026-08-13.*

## Now

- **dev**: `b2c3f89` — CI **green both jobs**, SHA-addressed (`checks=2` = expected 2). Nothing
  landed to code this iteration; the deliverable is a design doc plus a measured park.
- **Item 16 `w-broker-base-flake` is DESIGNED and PARKED `needs-human-review`.** Doc
  `design_docs/planned/w-broker-base-flake.md` (586 lines), two quorum rounds, **both BLOCKED**,
  both external reviewers **present** in both (`absent_reviewers` empty). R2's objections differ in
  kind: gemini's is a concrete `sync.Mutex` fix (carve-out material), gpt5's **reverses the doc's
  central architectural choice**, so the carve-out does not apply and Standing rule 2 binds.
- **THE QUEUE ROW'S NUMBERS AND ITS MECHANISM ARE BOTH DEAD.** The flake is real — reproduced once,
  `FAIL (5.28s)` against a 100 ms bound — but the rate is **1 in 132 runs = 0.76%**, not the row's
  ~18%. So the row's own `-count=20` acceptance proof reds only **~14%** of the time: a coin-flip
  gate, which **S6** forbids. Second consecutive item where the measurement killed the prescription.
- **THE SPINE — MY OWN 1,987-RUN EXONERATION WAS SCOPED TO A REGIME THE TEST NEVER ENTERS.** My
  probe mirrored `runBounded` line for line and re-used **one warm fixture inode**; the real test
  writes a **fresh** one every run, and on darwin that is **103 ms median vs 3 ms** (24/25 fresh
  execs exceed 100 ms) — the same order as the deadline under test. A faithful mock of a mechanism
  is not a faithful mock of its **inputs**. Corollary the designer drew out: the committed test is
  **partially vacuous per-run on darwin**, because the deadline usually expires before the
  grandchild exists at all.

## Parked on Mark

**ONE open ask — a one-word answer, `A` or `B`, on item 16:**
- **A** (as designed) — bring M2's `killGroup` seam forward into M1: `var killGroup = func(pgid int)
  error` in `host/broker/handlers.go`, +5/−1 behaviour-identical lines, so the test records the
  kill's count, monotonic offset, pgid and errno and makes the ~0.76% event **decisive** the one
  time CI catches it. Costs M1 its "zero production bytes" property.
- **B** — keep M1 production-free, narrow its claim to "localise the stall to the `Execute` window",
  defer mechanism selection to an external tracer or a separate doc. Costs a second rare-event wait.

Measured for the decision: gpt5's "frozen core" framing is **wrong** (`CLAUDE.md:25` scopes it to
`tools/launchd/*` + skills, not `host/`), but its **catch is right** — **S3** does bind `host/`
("why is this not a package?") and the doc's §10 dismisses rather than answers it. Exactly **one**
precedent exists, same shape, landed through this loop: `host/store/store.go:863`
`var commitBeforeOutcomeHook = func() {}` (`6811604`, PR #19).

Also owed by the **shared driver** (frozen core — World cannot apply): `ailang#611` (real per-role
executor chain) and the World driver sync (missing `pi:*` pre-flight).

## Next

**Item 16** on Mark's one word; failing that **item 13** `w-evidence-grade-mapping` (cheapest
high-leverage UI item; `PROVEN` is currently unreachable). **Item 5 `P6.B` remains UNBLOCKED.**
`SM.D` (item 8) is attended-only; 13/14/15 attended.

## Loop

- launchd, ~6h, headless. Issue **#53** (rotates Mondays 07:00 **local**; not due — created Mon
  07:37 local, 16 comments). Running skill **== origin** for the first time in 4 iterations.
- controller `opus` · designer `claude:claude-fable-5` (rotation; the rotation file is **rig-shared,
  not mission-namespaced**) · planner `opus` (fail-closed, `missing-script`) · executor
  `codex:gpt-5.6-sol` · evaluator `sonnet`. `pi` **BARRED** from publish milestones.
- `metered=$0.202435` this iteration (quorum only, both caps raised pre-emptively). Cap $5.

## Standing hazards

- **SCOPE, NOT MECHANISM, IS WHAT BREAKS A CONTROL.** Three instances in one iteration: a probe
  faithful in every syscall but not in its fixture lifecycle; a `-run` naming a **nonexistent test**
  (`no tests to run` → `PASS` → **rc=0**); and `grep -rl … | head` making the pipeline's status
  `head`'s, so the `||` fallback could never fire. Pair every check with a control **in the same
  call and the same scope**, and read `grep`'s exit code (1 = no match, 2 = no such file).
- **GUARD THE HELPER, MISS THE BRANCH/CALL SITE — now five directions** (mechanism/sites ·
  site/second branch · branch/shape-space · shape-space/spelling · **recogniser/enumerator**).
- A queue row's **prescription** rots faster than its diagnosis — re-measure both at pick time.
- `go build ./...` does not compile `_test.go`; use `go vet`. `git diff` omits untracked files.
- `verify_go.sh` needs `GOTOOLCHAIN=go1.25.6` and an `AILANG_BIN` reporting exactly `v0.30.0`;
  without `AILANG_BIN`, `host/broker` is red 100% of the time (correct behaviour, not the flake).
