## Harness

`TestProtectedV24RealWorktree` runs three always-on acceptance controls, then
measures construction and the production integrity checker at each worktree root.
It uses the V23 project pattern sets through the real renderer, mirror and parser.
Each project gets one discarded warm-up, 20 warm samples and 20 cold samples;
per-leg timings, census, verdicts, bounds and the ignored post-snapshot creation
and removal controls are written to `WORLD_PROTECTED_V24_OUT` when set.
An unset input root skips measurement and never satisfies V24.

`V24_SCRATCH` contains disposable source clones at `src/sunholo-platform` and
`src/ailang`, their installed warm worktrees at `sunholo-platform` and `ailang`,
and matching `<project>.deps` manifests (one relative entry per line). Cold
worktrees temporarily receive these dependency entries by rename; deferred
cleanup restores them before removing the cold worktree.

Cold means first Git access on a freshly added worktree with a fresh index and
no untracked cache. Construction runs before the phase. OS page cache is not
purged: this is never an OS-cold claim. Dependencies precede the snapshot in this
harness; M4-wide install ordering remains unproven (R-156-6).

With `WORLD_TOOL_AILANG_BIN` unset and `TMPDIR`/`GOTMPDIR` set to `$PWD/.gotmp`:

```sh
WORLD_PROTECTED_V24_ROOT="$V24_SCRATCH" WORLD_PROTECTED_V24_OUT=<abs path> go test -count=1 ./host/broker/ -run '^TestProtectedV24RealWorktree$' -timeout 900s -v
```

Timing or socket observations inside the sandbox are UNINFORMATIVE UNDER SANDBOX;
the controller measures outside it.

## Results

Controller-measured outside the sandbox, iteration 290 (2026-10-20), on base `17c28324`.
Two runs, rig load average 14–19 during both (`run1.test.log`, `run2.test.log`; both
`result.json` bodies under `runs`). Both runs: `--- PASS: TestProtectedV24RealWorktree/measurement`.

Scratch: `git clone --shared --no-checkout` copies of sunholo-platform `fe80714f` and
ailang `f0f7df204`, one `git worktree add --detach … HEAD` each, and the real checkouts'
ignored entries APFS-clone-copied in (`cp -cRp`), excluding agent and mission state
(`.claude/`, `.ailang/state/`, `.eval_workspace/`, `.dev-logs/`, `.env` files,
`settings.local.json`) and one untracked-not-ignored directory. 51 and 206 entries;
93 732 and 107 711 ignored files.

| Project / class | Phase p95 / max ms (run 1; run 2) | Construction p95 / max ms (run 1; run 2) |
|---|---|---|
| sunholo-platform warm | 653 / 859; 484 / 573 | 468 / 787; 486 / 490 |
| sunholo-platform cold | 755 / 776; 662 / 686 | 918 / 1 007; 709 / 726 |
| ailang warm | 1 130 / 1 133; 819 / 881 | 1 167 / 1 216; 834 / 909 |
| ailang cold | 1 611 / **1 958**; 988 / 1 066 | 1 244 / **1 583**; 935 / 1 331 |

- All 160 samples: production `checkProtected` intact `true`, null 0, error 0, changed 0.
- The ignored leg is 70–80 % of each phase. On ailang it walks `docs/node_modules` (2.5 GB).
- Chosen bounds (max over both runs × 1.5, rounded up, never below the provisional floor):
  phase **2 937 ms**; construction **3 200 ms** (2 375 computed, floor kept).
- Post-snapshot dependency arm, both projects, both runs: a new ignored protected file
  reports `intact: false` with status `ignored`; removing it reports `intact: true`.
- Install order: documented, not enforced (QUICKSTART installs before the session; the
  handler is built on the episode's first exec call). A later install reports CHANGED,
  never a false intact (R-156-6).
- "Cold" is a fresh worktree's first Git access; the OS page cache was not purged.
