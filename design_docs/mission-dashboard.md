# AILANG World — mission dashboard

*Snapshot, overwritten every Gate 4. History lives in `world-mission.md` STATUS + `world-mission-log.md`.*

**As of** 2026-08-05, iteration 54 · dev @ `bbee51c` · CI green both jobs (SHA-addressed)

## In flight

- **Item 8 `w-self-mod-vertical` — `SM.B1` LANDED** (PR #43 → `1856bfb`, evaluator sonnet **91/100
  zero blocking**). `approval_claims`, schema **1→2**, atomic `AppendClaimedEffectIntent`, and
  `DD-3` closed loudly — a store left at `user_version = 1` now fails with
  `*LegacySchemaVersionError` instead of opening and silently skipping `schemaSQL`.
- **`[NEXT]` is `SM.B2a`** — brokered publish handler, de-ambient credential, typed indeterminate
  (~780 LOC). Gated on nothing. **`AC12`'s "network confined to `host/broker`" control is vacuous
  until exactly this milestone** — re-assert it there, never inherit it as green.
- **Item 5 `w-mcp-projection` — still BLOCKED** on one prerequisite (transition registry absent at
  HEAD). Unchanged.

## Latest — three findings that rhyme

- **SM.B1's milestone-gating ledger check was VACUOUS as delivered.**
  `TestSchemaVersionLedgerIsIndependent` greps its own source; its two *negative* needles were split
  so they wouldn't match their own check-lines, but the *positive* needle was one literal that did.
  `var schemaV2SQL = string(schemaSQL)` — the ledger becoming the file it exists to attest — passed
  with `ok 0.290s`. The executor's own mutation redded only because it used the bare form its
  negative needle was written to catch. **A mutation shaped to the check tests the check, not the
  threat.** Repaired by anchoring to `^` plus a semantic backstop.
- **A 15.7 MB `ailang-worldd` Mach-O had been tracked since SM.A** — past the executor, an evaluator
  at 87/100 *zero blocking*, four controller gates and both CI jobs, because nothing enumerated file
  *types*. Removed; `verify_go.sh` now reds on any tracked binary (PR #42 → `e24a6f0`).
- **Dev went red on this iteration's own docs-only commit** — a SIGPIPE race in SM.A's CI step
  (`--version | grep -q` under `pipefail`). Measured **3/40** vs **0/200** no-pipe. The first fix was
  wrong (capturing still pipes) and the stress arm caught it. Both sites converted — the second was
  safe only by **accident of size** (167 B fits the 64 KiB buffer). PR #44 → `bbee51c`.

## Loop · cost · asks

- launchd `mission-world`; controller `claude-opus-5`. Executor **`codex:gpt-5.6-sol`** (1 bounded
  run, rc=0); evaluator **`sonnet`** (generator≠judge). Designer/planner **not fired**; rotation
  unchanged at `codex:gpt-5.6-sol`. Verify profile `ailang-code`; AILANG pinned **v0.30.0**. Issue **#32**.
- **`metered=$0.00`** vs the $5 ceiling — every role on a quota bucket.
- **Parked on Mark: NONE.** `8/OD-2` open, non-blocking. FYI not blocking: item 9's human-gated half
  (pin CI job 1 vs keep tracking `latest`); the rig's PATH `ailang` has drifted to `v0.33.0-23-…-dirty`.
