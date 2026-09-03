# Mission Dashboard — AILANG World

_Snapshot, overwritten every iteration. History: `world-mission.md` (STATUS), `world-mission-log.md`._

**As of:** 2026-09-03 · iter 151 · `dev` = [`a7b58dd`](https://github.com/sunholo-data/ailang-world/commit/a7b58dd) · CI green (3/3) · local verify gate **restored to green**

## Last iteration
**No PR — bookkeeping repair + a precondition restore.** **HARNESS** · metered **$0.00** (controller-authored; no sub-agent spawned).

**Two consecutive slots died.** Slot 1 (unnumbered, 2026-09-02) merged row 56 as PR [#113](https://github.com/sunholo-data/ailang-world/pull/113) → `725ad5a` and left **zero** trace in all four mission docs — iter-150 then read that SHA as a CI negative control without asking what it was. Slot 2 (iter-150) stamped `gate-4` at epoch `1788394743`; `kern.boottime` = `1788395029`, so the rig rebooted **286 s later** and macOS wiped `/private/tmp`. Both records now landed, row 56 tagged with retroactive Gate-3b evidence, stale worktree pruned.

**The reboot had also deleted the toolchain pin** `/tmp/ailang-v0300/ailang` — the binary every gate here runs on. It fails *closed*, so nothing was falsely green; but the loop could not verify anything and nothing said so. Restored to **`~/.pinned-ailang/ailang`** (CI's own path; `$HOME` survives a boot), checksum- and version-verified; `verify_ail.sh` step 9/9 now prints `compiler pinned by exact bytes: AILANG v0.30.0 on Darwin/arm64`.

## Goal distance
**59 of 72 rows closed — carried, not measured.** Two re-derivations read 36/72 and 53/70 against the carried 58/70, and the better one still over-counts (it classes the open row 57 as closed on prose). Now row 72. Row 50 parked on `D-WORLD-31`.

## Next picks
1. **Row 57** `w-approvals-spine-prints-a-green-no-pending` — queue head, ungated.
2. **Row 71** `w-mission-critical-state-lives-in-a-directory-the-os-wipes-on-boot` — pin half closed; the driver's **only** crash log and sprint worktrees are still in `/tmp`. Fleet-owned.
3. **Row 69** `w-heartbeat-script-absent` — fleet-owned port. Then 58–66, 68, 70, 72, then 39.

## Loop + routing
Controller `claude:claude-opus-5` · designer ROTATION, last used **`claude:claude-fable-5`** (next: `pi:ollama/deepseek-v4-flash:0731-cloud`) · planner `opus` · executor `codex:gpt-5.6-sol` · evaluator `sonnet` (generator≠judge). Iter-151 spawned **no** role: the deliverable was verification and bookkeeping over existing artifacts, and routing a sprint before restoring the pin would have produced a sprint whose gate could not run.

## Parked on Mark
**`D-WORLD-31`** — one word. Ship `D-WORLD-29`'s rule A as ratified, or hold row 50 for the fixture migration. Unchanged, re-asked. Nothing else is blocked; the queue advances either way.

## Standing reds — owned elsewhere, none is a World failure
- **`verify_go.sh` RED on the rig:** the fleet arm — World's driver copy is behind fleet HEAD, i.e. *the fleet must commit*. **CI unaffected** (that arm loud-skips there). Both charter hard-rule legs are **GREEN** this iteration: `verify_ail.sh` rc=0, and `go build ./... && go test ./...` rc=0 (19 ok / 0 FAIL) — row 58's flake did not fire.
- **The running `mission-control` skill is byte-identical to `origin/dev`** (`cmp` against the resolved symlink target).
- **`tools/launchd/mission-heartbeat.sh` is absent here** (row 69) — this iteration lost its `gate-0` stamp to it before switching to V1's absolute path.
