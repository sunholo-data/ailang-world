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

Measured outside the sandbox on 2026-10-20 (iteration 290, base `17c28324`). There were five runs:

- **Runs 1–3, controller.** The rig load average was 14–19 from other missions. Runs 1 and 2 used the harness at `0637d1b3`, whose measurement path is identical to the final one. Run 3 used the final harness, with the cold-class guard on.
- **Run A, independent judge (sonnet).** Its own mutation fixtures were running at the same time.
- **Run B, independent judge.** Nothing else was running.

Each run's body is under `runs` in `result.json`, with its log beside it.

The scratch was built from `git clone --shared --no-checkout` copies of sunholo-platform `fe80714f` and ailang `f0f7df204`. Each got one `git worktree add --detach … HEAD`. The real checkouts' ignored entries were then APFS-clone-copied in (`cp -cRp`). Excluded were agent and mission state (`.claude/`, `.ailang/state/`, `.eval_workspace/`, `.dev-logs/`), `.env` files, `settings.local.json` and one untracked-not-ignored directory. That left 51 and 206 entries, holding 93 732 and 107 711 ignored files.

| Run | Verdicts | Slowest phase ms | Slowest construction ms |
|---|---|---|---|
| 1 | 80/80 intact | 1 958 | 1 583 |
| 2 | 80/80 intact | 1 066 | 1 331 |
| 3 | 80/80 intact | 954 | 1 052 |
| A | 1 `null` (check timed out at 2 500 ms; harness FAILED closed) | **2 501** | **3 232** |
| B | 80/80 intact | 1 779 | 1 842 |

Run 1 p95 / max per class, phase then construction, in ms:

| Project | Warm phase | Cold phase | Warm construction | Cold construction |
|---|---|---|---|---|
| sunholo-platform | 653 / 859 | 755 / 776 | 468 / 787 | 918 / 1 007 |
| ailang | 1 130 / 1 133 | 1 611 / 1 958 | 1 167 / 1 216 | 1 244 / 1 583 |

**Bounds.** Each bound is the maximum over all five runs × 1.5, rounded up, and never below the provisional floor.

- Phase: **3 751 ms**.
- Construction: **4 848 ms**.
- Both maxima come from run A, ailang cold. The 2 501 ms phase is an intact sample. The 3 232 ms construction is sample 17, whose production check then timed out (`null`), so the construction bound is sized from a failed sample and errs long.
- Runs 1–3 alone would have given 2 937 / 3 200.
- So the 1.5× margin over one unloaded run did not cover one loaded run. These bounds are not load-proof. An over-bound check records `error` / `null`, never intact.

**The ignored leg** is 60–95 % of a phase. On ailang it walks `docs/node_modules`, which is 2.5 GB.

**Post-snapshot dependency arm** (both projects, every run): a new ignored protected file reports `intact: false` with status `ignored`. Removing it reports `intact: true`.

**Install order** is documented but not enforced. QUICKSTART installs dependencies before the session, and the handler is built on the episode's first exec call. A later install reports CHANGED, never a false intact (R-156-6).

**"Cold"** here means a fresh worktree's first Git access. The OS page cache was not purged.
