# AILANG World — mission dashboard

*Snapshot, overwritten every Gate 4. History lives in `world-mission.md` STATUS + `world-mission-log.md`.*

**As of** 2026-08-07, iteration 62 · dev @ `3fd889f` · status API **All Systems Operational, 0
incidents** — this green is attributable, and licenses an infrastructure inference as well as a code one.

## In flight

- **Item 8 `w-self-mod-vertical` — `SM.B2a` LANDED** (PR #50 → `3fd889f`, evaluator `sonnet`
  **98/100 r1, zero blocking**). Brokered publish handler, de-ambient credential, typed
  indeterminate outcome. `AC7` · `AC10` · `AC11` discharged.
- **`[NEXT]` is `SM.B2b`** — `AC8` (dispatch half) + `AC9`/`AC9a`/`AC9b`/`AC9c`: attended-stamp
  binding and single-use approval consumption. Gated on nothing. `SM.B2a` wired
  `AppendClaimedEffectIntent` but does **no approval validation** — that is `SM.B2b`'s whole job.
- **Item 5 `w-mcp-projection` — still BLOCKED** on one prerequisite (the transition registry).

## Latest — the credential was already leaking, and a replaced AC is what found it

The doc's original `AC10` ("all non-publish subprocesses observe it unset") is satisfiable by
launching **zero** subprocesses. Iter-52's planner replaced it with one that must re-derive the site
count **by command in-run**, drive **every** site, and `t.Fatal` on an empty enumeration. Executing
that literally measured **two of five production subprocess sites leaking a live,
irreversible-publish credential** (verified first-party at base `0c47667`):

| site | defect at base |
|---|---|
| `host/archive/archive.go` `probeVersion` | bare `exec.Command(...).CombinedOutput()` — **no `cmd.Env` at all** |
| `host/replay/replay.go` `runPinnedTransition` | sets `Dir`/`Stdout`/`Stderr`, **not `Env`** |

Both fixed; `host/childenv` holds the variable list once so four packages cannot drift.

**Judge `NB-1`, fixed not carried** — the *direction* is why: `Scrubbed` returned **nil** for
degenerate inputs and `exec` reads a nil `cmd.Env` as **INHERIT**. Fail-OPEN, in the one package
written to prevent it. Now always non-nil; guard proven by a **compiling** mutant (`f9e2e40`).

**Twice a mutation redded for the wrong reason** — one never landed (sha256 unchanged), one failed to
build. *"Vacuous"*, *"never ran"* and *"doesn't build"* are three facts wearing one exit code.

## Loop · cost · asks

- Controller `claude-opus-5`. Executor **`opus`** — codex quota-dry to Aug 8 11:24 **and** `pi` barred
  for publish-capable code (Mark, attended 2026-08-06); documented fallback, FLAGGED. Evaluator
  **`sonnet`** ⇒ generator≠judge. AILANG pinned **v0.30.0**. Issue **#32**.
- **`metered=$0.00`** vs the $5 ceiling — all roles on quota buckets.
- **Safety:** `ailang publish` never invoked in any form; no non-loopback request; no secret printed.
- **Parked on Mark: NONE.**
