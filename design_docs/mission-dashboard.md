# AILANG World — mission dashboard

*Snapshot, overwritten every Gate 4. History lives in `world-mission.md` STATUS + `world-mission-log.md`.*

**As of** 2026-08-05, iteration 54 · dev @ `1856bfb` · CI green both jobs (SHA-addressed on the merge commit)

## In flight

- **Item 8 `w-self-mod-vertical` — milestone `SM.B1` LANDED** (PR #43 → squash `1856bfb`, evaluator
  sonnet **91/100 zero blocking**). `approval_claims`, schema **1→2**, atomic
  `AppendClaimedEffectIntent`, and `DD-3`'s silent hole closed loudly — a store left at
  `user_version = 1` now fails with `*LegacySchemaVersionError` instead of opening successfully and
  never running `schemaSQL`.
- **`[NEXT]` is `SM.B2a`** — the brokered publish handler, the de-ambient credential, the typed
  indeterminate (~780 LOC). Gated on nothing. **`AC12`'s "network confined to `host/broker`" control
  stops being vacuous exactly here** — re-assert it when the first `net/http` dep lands.
- **Item 5 `w-mcp-projection` — still BLOCKED** on one prerequisite (transition registry absent at
  HEAD). Unchanged this iteration.

## Latest

- **A 15.7 MB compiled binary had been sitting in the repo since SM.A, and five independent checks
  passed it.** `ailang-worldd` (darwin/arm64 Mach-O) was added by `13315da` and survived the codex
  executor, a sonnet evaluator at 87/100 **zero blocking**, the controller's four-gate re-run and
  both CI jobs — because nothing enumerated tracked file *types*. Removed, and `verify_go.sh` now
  reds on any tracked binary blob using git's own classification (portable darwin↔ubuntu, no
  allowlist). Mutation-proven both ways, including a blinded-detector arm.
- **SM.B1's milestone-gating ledger check was VACUOUS as delivered.**
  `TestSchemaVersionLedgerIsIndependent` greps its own source; its two *negative* needles were split
  so they wouldn't match their own check-lines, but the *positive* needle was one literal that did.
  `var schemaV2SQL = string(schemaSQL)` — the ledger becoming the file it exists to attest — passed
  with `ok 0.290s`. Repaired by anchoring to `^` plus a semantic backstop.
- **The executor's own mutation could not have caught it**: it used the bare form its negative needle
  was written to catch. **A mutation shaped to the check tests the check, not the threat.**

## Loop · cost · asks

- launchd `mission-world`; controller `claude-opus-5`. Executor **`codex:gpt-5.6-sol`** (one bounded
  30-min run, rc=0); evaluator **`sonnet`** (generator≠judge, cross-provider). Designer/planner
  **not fired** — doc and plan both landed; designer rotation unchanged at `codex:gpt-5.6-sol`.
- Verify profile `ailang-code`; AILANG pinned **v0.30.0** at `/tmp/ailang-v0300/ailang`. Issue **#32**.
- **`metered=$0.00`** vs the $5 ceiling — every role on a quota bucket.
- **Parked on Mark: NONE.** `8/OD-2` open but non-blocking. Worth his attention, not blocking:
  item 9's human-gated half (pin CI job 1 vs keep tracking `latest`) — the rig's PATH `ailang` has
  drifted again, now `v0.33.0-23-g78f30e053-dirty`.
