# World mission dashboard — iteration 440 (2026-11-26)

- **Done this fire:** row 243 **V4 sprint 5f2 LANDED** (PR #685, merge `150d0d96`). The judge gave **FAIL 67** in round 1 and **PASS 77** in round 2, with 0 blocking findings.
  - **The activation prerequisites are in.** The worktree and files helpers now borrow and verify the runner's own episode lock instead of opening their own. Judge-439's NULL-catalogue (J01) and handle-close (J05) pins are in. N04/N12/N15 pin the publication windows.
  - **Live change on every install:** the daemon's auto-GC and `episode gc` now take that per-episode lock (non-blocking) and re-check before removing a worktree. When unsure they skip with a reason, and they never remove. Side effects: GC now leaves a small `.operator-episode-<ep>.lock` sidecar per episode, and `--credentials-db` for a non-default workspace skips every worktree (judge notes N3, N2).
  - 1,476 lines, 76 over the cap from the round-2 test fix. Nothing dispatches yet: no production caller.
- **Why the judge failed round 1:** M2's test helpers grabbed the lock before the code under test, so three existing refusal tests went red or dead. My whole-package fail-set comparison showed the same thing. The codex round-2 fix was test-only (80 lines), and each arm now kills a named production mutation.
- **Why activation moved again (to S5f3):** the planner measured the full caller at 1,984 cumulative lines. Even its core plus these prerequisites is 1,702, over the 1,650 ceiling. The upgrade path still cannot drain a paired runner. The judge measured about 584 lines over this head, so S5f3 fits one sprint. This is the third deferral in a row.
- **My own slip:** the first executor run stopped correctly on an `ACCEPT_OK` marker I had carried over from iteration 439's directive, which this plan's script does not print. I corrected it and re-ran.
- **DECISIONS FOR MARK:** **D-WORLD-121** — take row 267 (CI race-bound headroom) before S5f3? Recommended **A**; default B (groomed order), which applied this fire. This sprint's PR and merge race jobs were both green at attempt 1 (`cmd/ailang-worldd` 273.2 s and 351.0 s, `host/daemon` 360.3 s and 409.3 s, against a 480 s bound).
- **Next loop work:** **V4 S5f3, activation.** It brings the first production caller with:
  - the attended paired drain and rollback, and the reconcile matrix
  - pre-digest signed-start compatibility, terminal `session revoke` exclusion and J-OWN-UID
  - the four process pins
- **Bar:** 1.0 clauses 1–7 MET; no UNMET clause. Goal unmoved; row 243 serves D-WORLD-94.
- **Routing this fire:** planner and executor `codex:gpt-6.2-sol` (914k planner; 75k + 104k executor tok). The judge was Agent `sonnet` (128k + 64k tok), a different vendor. Metered $0.
- **Open PRs:** #617 is attended UI work that the loop does not touch.
- **Bookkeeping thread:** #635 (next rotation Monday 2026-11-30 07:00 local).
