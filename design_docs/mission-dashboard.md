# AILANG World — mission dashboard

*Snapshot, overwritten every Gate 4. History lives in `world-mission.md` STATUS + `world-mission-log.md`.*

**As of** 2026-08-07, iteration 60 · dev @ `39130ec` · CI **GREEN at HEAD**, both jobs, SHA-addressed,
step-log verified (`ailang-code verify gate` 11/11 · `go host build + test gate` 13/13, `failed=none`).
Iteration 59's red is **retired by outcome divergence**: the identical tree that was `cancelled/steps=0`
is now `success` with 24 steps executed. Iteration 59 was right — it was the provider.

## In flight

- **Item 10 `w-boundary-gate-tree-mutation` — `BG.B` LANDED** (PR #48 → squash `39130ec`; evaluator
  `sonnet` **88/100 r1, zero blocking**). `AC1a` discharged; `M3`, `M6` and the deny-list control all
  fired with their control arms.
- **`[NEXT]` is milestone `BG.C`** — the runtime backstop (`AC1b`, `M7`), gated on nothing. Carry:
  `C1`'s nanosecond-`ModTime` premise is **APFS-only** (200/200 on darwin, unmeasured on CI's ext4);
  `BG.C` is a fail-loud 20/20 granularity probe with a pre-authorized refutation path (`10/OD-3`).
- **Item 8 `w-self-mod-vertical`** — `SM.B2a` next after item 10. Unchanged.
- **Item 5 `w-mcp-projection` — still BLOCKED** on one prerequisite (the transition registry). Unchanged.

## Latest — the test that proved the writer was blind to the writer

`BG.B` routes every harness write through one `confinedWrite` that rejects `repoRoot` destinations
**before a byte is written**, and installs an **AST** write-guard that must prove it can SEE before its
silence counts. The finding is in the *test*, not the writer:

- The recording-writer test as first delivered **synthesised its own paths** and called `confinedWrite`
  directly instead of driving the real harness. Its "exact count = 2" counted only the writes it had
  itself just made — self-fulfilling.
- **Measured**: with `mutateViaOverlay`'s own writes reverted to bare `os.WriteFile`, it still
  **PASSED 4/4 arms**. Blind to the write path it claimed to cover.
- Repaired to drive the real harness with the sink **teed**; arm mutations factored into shared helpers
  so gate and test cannot drift. **Non-vacuity by outcome divergence**: the identical probe now
  **REDS 4/4** — `the harness recorded ZERO writes through the confined sink`.

| mutation | observed |
|---|---|
| `M3` `confinedWrite` on a live repo path | rejected, message exactly as predicted, target sha256 **UNCHANGED** |
| `M6` the pre-BG.A defect verbatim (direct write + deferred restore) | AST guard **REDS** naming `:503` and `:506` |
| deny-list truncated to zero | `AST deny-list has 0 entries, want 4` |

Gates baselined on a **pristine** tree first (rule 3e) so a red would be attributable: `verify_go.sh`
rc=0, `verify_ail.sh` rc=0, boundary gate rc=0 — before any work, and again after.

## Loop · cost · asks

- launchd `mission-world`; controller `claude-opus-5`. Executor **`pi:deepseek-v4-flash-0731`**
  (codex bucket dry, resets **Aug 8 11:24** — quota-relief policy, probe rc=1 first-party);
  evaluator **`sonnet`** (distinct provider ⇒ generator≠judge). No designer/planner fired.
  Verify profile `ailang-code`; AILANG pinned **v0.30.0**. Issue **#32**.
- **`metered=$0.024`** vs the $5 ceiling (the pi executor run; every other role on a quota bucket).
- **Parked on Mark: NONE.** V1 iter-156 **accepted both** shared-skill proposals — the
  green-during-outage rule is **live** in the shared skill; the case-sensitive stale-charter tell is
  queued as V1 iter-157's edit. Workaround (`grep -ci`) in use here until it lands.
