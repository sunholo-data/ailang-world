# AILANG World — mission dashboard

*Snapshot, overwritten every iteration. History: `world-mission.md` (STATUS) + `world-mission-log.md`.
Last written: iteration 73, 2026-08-11.*

## Now

- **dev**: `6e207ca` — CI **green both jobs**, SHA-addressed, `checks=2` = expected 2, **0 incidents**.
- **Last landed**: **`TR.B1`** (PR #61) — broker `CapabilitySnapshot` (epoch on debit, one
  `debitGrant` mechanism), `Allows` delegating to the landed `Decide` via a single `decideOver`,
  and the confined `BoundInvoker` seam. **AC5 activated** (exactly 2). Evaluator `sonnet`
  **84/100**; its one blocking finding reproduced first-party and **FIXED in-PR**.
- **`Invoke` → unexported `invoke`**, `Invoke` a one-line wrapper: the bound invoker would have been
  a **4th** production `Invoke` selector call, and TR.C pins exactly **3**. New gate `AC-INVOKE3`.

## Parked on Mark

**Nothing.** Zero open asks. Owed by the **shared driver** (frozen core — World cannot apply):
`ailang#611` (real per-role executor chain) and the World driver sync (missing `pi:*` pre-flight).

## Next

**`TR.B2`** (descriptor-bound confinement + two-session fixture, AC6/AC7 — plan §3 T4–T5 already
written), then **`TR.C`**, the binding gate — P6.B's prerequisite is satisfied only when `TR.C` is
green. Then item 12, and new **item 16** (the `host/broker` ~18% base flake). `SM.D` (item 8) is
attended-only; items 13/14/15 (UI programme) were filed attended.

## Loop

- launchd, ~6h, headless. Issue **#53** (rotates Mondays 07:00 **local**).
- controller/planner `opus` · designer rotation (last `codex:gpt-5.6-sol`) · executor
  `codex:gpt-5.6-sol` · evaluator `sonnet`. `pi` **BARRED** from publish milestones.
- `derive-planner-lane.sh` absent → lane fails closed to opus, loudly. `metered=$0.00`; cap $5.

## Standing hazards

- **GUARD THE HELPER, MISS THE CALL SITE** — now this repo's most reliably recurring defect class:
  3 instances in `TR.A2`, **2 more in `TR.B1`** (an un-copied slice on `Bind`'s INPUT side while its
  output accessor was pinned; a 3rd `debitGrant` call site pinned by nothing). Unifying N call sites
  into one mechanism makes you test **the mechanism** and stop testing **the sites**.
- **`rg` is NOT a binary here** — a harness-injected shell function, absent in CI. Use `grep`.
- **A refusal test asserting only *that* an error occurred pins no branch.** Pin the measured message.
- **A rule-3j audit anchored to a DECISION LIST cannot contain branches the sprint itself writes.**
- **A green `go test` is not a green `go vet`** — `copylocks` is outside both, invisible to CI.
- **`host/broker` is ~18% flaky at base** (`TestHandlerTimeoutKillsTheWholeProcessGroup`, 2/11) and
  **100% red without `AILANG_BIN`** — both fake mutation kills and falsify inverse arms. Item 16.
- **`verify_ail.sh` never asserts the module count against 11** — only against 0. Item 12.
