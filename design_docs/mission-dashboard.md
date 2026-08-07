# AILANG World — mission dashboard

*Snapshot, overwritten every Gate 4. History: `world-mission.md` STATUS + `world-mission-log.md`.*

**As of** 2026-08-08, iteration 63 · dev @ `abb3a3d` · status API **All Systems Operational, 0
incidents** — this green is attributable, licensing an infrastructure inference as well as a code one.

## In flight

- **Item 8 `w-self-mod-vertical` — `SM.B2b` LANDED** (PR #51 → `abb3a3d`, evaluator `sonnet`
  **74/100, one blocking — closed and re-verified**): an attended approval now binds **bytes, not a
  name**, and is spent **exactly once, durably**. `AC8`/`AC9`/`AC9a`/`AC9b`/`AC9c` discharged;
  SM.B1's carried `NB-2` closed.
- **`[NEXT]` is `SM.C`** — probe-then-resolve reconciliation, replay evidence, clean-room fixture,
  attended runbook. Gated on nothing.
- **Item 5 `w-mcp-projection` — still BLOCKED** on one prerequisite (the transition registry).

## Latest — a guard is not a gate until something reds when you remove it

The milestone was **inherited from a dead iteration**: Gate 2's rule (c) found an orphan worktree
holding five uncommitted production files, no commit, no branch, zero charter rows. It built and
vetted clean — and **redded five landed SM.B2a tests**, with a pristine-base control green. Then a
three-stage cascade, each stage a different role, each measured. The **controller** mutated
the executor's own new gate: neutering `scope.Effect != EffectRegistryPublish` left the **entire
`host/broker` package green** — the neighbouring arm rejects at the parser and never reaches the
term. The **evaluator**, given that as a named target, found **six more**; the **executor** audited
the function and found **twelve**. `AC9` now carries **20 negative arms, one per refusal branch**;
all seven policy branches re-verified RED by the controller.

**Rule 3d, bought and caught in one breath:** an expiry mutation redded in exactly the predicted
direction — and the FAIL was a pre-existing load flake. *A red is not evidence until you can name
the test that produced it.* Also: the doc **contradicts itself in one paragraph**.

## Loop · cost · asks

- Controller `claude-opus-5`. Executor **`opus`** — codex measured quota-dry (resets Aug 8 11:24)
  **and** `pi` barred for publish-capable code (Mark, attended 2026-08-06); FLAGGED. Evaluator
  **`sonnet`** ⇒ generator≠judge. AILANG pinned **v0.30.0**. Issue **#32**. **`metered=$0.00`**.
- **Safety:** `ailang publish` never invoked; every publisher is the re-exec'd test binary on
  loopback; no secret printed. **Parked on Mark: NONE.**
