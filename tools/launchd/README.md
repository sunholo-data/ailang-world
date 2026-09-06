# tools/launchd — DE-FORKED 2026-09-06

This directory used to hold **a 1,047-line copy of `mission-control.sh`**, plus copies of
its test suites and this mission's plist. They are gone.

The world mission now runs the **shared driver** from `sunholo-data/ailang`:

```
launchd job  dev.ailang.mission-world
  ProgramArguments  /bin/bash <ailang>/tools/launchd/mission-control.sh
  MISSION_PROFILE   world
  MISSION_WORKDIR   <this repo>          <- the mission still WORKS here
```

## Why the fork existed, and why it does not need to

The driver computed its own root as `${MISSION_WORKDIR:-<script's ../..>}`, so "where the
driver LIVES" and "where the mission WORKS" were the same thing. A mission working in
another repo therefore *needed* its own driver copy. That was not a quirk of this mission
— it is what the model forced.

The cost was real and measured: this fork silently missed every routing fix landed
upstream for weeks. It had no memory gate, no boot stagger, no role-generic pre-flight,
and no lane-degradation ledger — the ledger having been written *because of this
mission's* five silently-demoted iterations (18/19/21/22, 2026-07).

`internal/mission` now separates the two. A mission elsewhere is just a different
`workdir`.

## Changing how this loop runs

Edit the registry entry in the **ailang** repo — not here:

```
missions/world.toml          schedule, workdir, doc
ailang mission doctor        does what is installed match what was reviewed?
ailang mission install world render the artifacts to *.staged
ailang mission apply world   promote them and reload launchd
```

## Known state

**This mission is UNPINNED.** `pin-root.sh` re-execs from a worktree of `$REPO`, which is
this repo, and this repo has no `lib/pin-root.sh` — so every fire logs `DRIVER PIN FAILED`
and runs the working tree. That is **loud, not silent**, and it is no worse than before
(this mission was never pinned). Pinning it needs driver-root and work-dir decoupled
inside `pin-root.sh` itself — see `m-driver-pin-rollout`, which is parked on that same
question.
