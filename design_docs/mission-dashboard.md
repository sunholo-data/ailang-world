# AILANG World — mission dashboard

*Snapshot, overwritten every Gate 4. History: `world-mission.md` STATUS + `world-mission-log.md`.*

**As of** 2026-08-08, iteration 65 · dev @ `0cd00eb` · status API **All Systems Operational, 0
incidents** — so the green is attributable (infrastructure inference, not just a code one).

## In flight

- **Item 8 `w-self-mod-vertical` — `SM.C` LANDED** (PR #52 → `0cd00eb`, evaluator `sonnet`
  **93/100, zero blocking**): probe-then-resolve reconciliation, replay evidence, attended runbook.
  `AC13`–`AC17` discharged; 8 named mutations + 23/23 refusal branches pinned.
- **`SM.D` is BLOCKED on `8/OD-1`** (attended, never headless/CI), so **item 8 has no
  headless-routable milestone left** — the next iteration must pick elsewhere.
- **Item 5 `w-mcp-projection` — still BLOCKED** on one prerequisite (the transition registry).

## Latest — a clean `rc=0` is what a dead iteration looks like

Gate 2 found a second consecutive orphan worktree: 525 lines of untested `SM.C` code, no commit,
no charter row. Rather than just adopt it, the loop **diagnosed why its slots die**. `Background
tasks still running after 600s; terminating.` appears **exactly twice in 67 iterations** — and those
two are **exactly** the two orphans. The controller spawns its executor as a background agent, ends
its turn to wait, and the harness kills it at 600 s; the driver then logs `rc=0` with **no watchdog
firing**. Survived by never ending the turn while the agent ran.

**Three more measured things.** `AC13`'s landed guard was satisfiable by a recovery that *does*
dispatch — it passes the real handler, which refuses before its own counter moves, and stayed GREEN
under mutation. **Rule 3e(b) caught my own contamination**: the executor's `verify_go.sh` was green,
mine was **red** — a `.go` file added to `host/boundary` trips its `wantFileCount = 1` pin (third
bite; fixed by moving the gate, not relaxing the pin). The evaluator's one finding reproduced: an
arm-(iii) fixture is a **bystander, not a guard**. `8/OD-2` → `sunholo-data/ailang#633`, all measured.

## Loop · cost · asks

- Controller `claude-opus-5`. Executor **`opus`** — codex measured quota-dry first-party (`rc=1`,
  resets 11:24) **and** `pi` barred for publish-capable code (Mark, attended 2026-08-06); FLAGGED.
  Evaluator **`sonnet`** ⇒ generator≠judge. AILANG **v0.30.0** · Issue **#32** · **`metered=$0.00`**.
- **Safety:** `ailang publish` never invoked; reconciliation's only verb is `GET`; no secret printed.
- **Parked on Mark:** **`8/OD-1`** (attended approval, blocks `SM.D` and thus all of item 8) · the
  driver env-var proposal (`CLAUDE_CODE_PRINT_BG_WAIT_CEILING_MS=0`; frozen core, World cannot apply).
