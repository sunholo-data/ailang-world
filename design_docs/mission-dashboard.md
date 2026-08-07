# AILANG World — mission dashboard

*Snapshot, overwritten every Gate 4. History lives in `world-mission.md` STATUS + `world-mission-log.md`.*

**As of** 2026-08-07, iteration 61 · dev @ `c6a14c0` · GitHub status API **All Systems Operational,
0 incidents** — so unlike iterations 58–60 a green here is attributable, not a sample.

## In flight

- **Item 10 `w-boundary-gate-tree-mutation` — `BG.C` LANDED, and the ITEM IS COMPLETE.** All three
  milestones shipped: `BG.A` `278f102` (iter-58) · `BG.B` `39130ec` (iter-60) · `BG.C` PR #49
  (iter-61, evaluator `sonnet` **94/100 r1, zero blocking**). Doc → `implemented/`.
- **`[NEXT]` is item 8 `w-self-mod-vertical`, milestone `SM.B2a`** — gated on nothing. `8/OD-1`
  (the attended publish authorization) was ANSWERED 2026-08-05 and is now registered.
- **Item 5 `w-mcp-projection` — still BLOCKED** on one prerequisite (the transition registry).

## Latest — a probe that must prove it can fire, on the filesystem it certifies

`BG.C` adds the runtime backstop: five observables of the live target captured before/after
`check()`, four asserted unconditionally, `mtime_ns` asserted **only** after a 20-trial probe fires
**20/20**. `<20/20` FAILS loudly — `verify_go.sh` runs `go test` without `-v`, so a passing test
prints nothing and the *assertion*, not the log line, is the only channel a measurement travels.

**The controller finding.** The plan RECORDS `st_dev` for `t.TempDir()` and `repoRoot` but never
COMPARES them — so a 20/20 probe on a fine-grained tmpfs would license an mtime assertion about a
file on a possibly coarse-grained repo volume. On this host the two devs are equal (`16777230`),
which is exactly why the gap is invisible without an assertion. Now asserted *before* the gate.

| mutation | observed |
|---|---|
| `M7` `cp` write+restore (off the AST deny-list) | `live-target nanosecond ModTime changed` |
| `M10` **new** — `mv` restore, content byte-identical | `live-target inode changed` |
| `M8` **new** — forced cross-filesystem probe | `probe measured a different filesystem than the live target` |
| `M9` **new** — probe write pair removed | `backstop is not armed…: probe fired 0/20` |

AST guard stayed **GREEN** in all four arms — so the backstop is live *independently* of the guard.

**`AC6′`**: 8 interleaved same-session pairs — median A `1.1275 s`, B `1.4700 s`, **ratio 1.3038 ≤
1.50**, absolute ≤ 3.0 s. The doc's ORIGINAL `AC6` (`≤2× 0.435 s`) would have failed **both** arms,
including arm A, which is unmodified base code. iter-57's replacement was not pedantry.

## Loop · cost · asks

- launchd `mission-world`; controller `claude-opus-5`. Executor **`pi:deepseek-v4-flash-0731`**
  (codex bucket dry until Aug 8 11:24); evaluator **`sonnet`** (distinct provider ⇒ generator≠judge).
  No designer/planner fired. Verify profile `ailang-code`; AILANG pinned **v0.30.0**. Issue **#32**.
- **`metered=$0.024`** vs the $5 ceiling — the pi executor run; every other role on a quota bucket.
- **Process fix landed:** the `OD-<n>` registry's own enumeration instrument returned **0** against a
  firing control, so **six** ODs were allocated unregistered — including `8/OD-1`, which is where
  Mark's attended publish approval landed. Instrument replaced, all six registered.
- **Parked on Mark: NONE.**
