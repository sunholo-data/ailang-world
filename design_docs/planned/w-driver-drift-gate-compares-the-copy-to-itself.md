# w-driver-drift-gate-compares-the-copy-to-itself — design for a fleet-source comparison arm in `scripts/verify_go.sh`

- **Status:** planned (design only; nothing implemented)
- **Date:** 2026-09-02 (revised 2026-09-02, quorum round 1 — full-strength, all three reviewers rejected on the enumeration domain and the honesty of the success line; direction Option B kept)
- **Base commit measured at:** `7d3aa00ef880657cc8a769f6726bb700aa54b5e0` (== `origin/dev`), worktree `sprint/w-iter148-driver-drift`
- **Owning queue row:** `design_docs/world-mission.md` row 54 — `w-driver-copy-stale-and-the-drift-gate-compares-it-to-itself` (clause-2, ~0.3d, gated on nothing)
- **Fleet source measured at:** `sunholo-data/ailang` HEAD `722e19c77bd7c5c4b7a8a47c9a31909e5387cc9d`

---

## 1. Problem

`launchd` runs **this repo's** copy of the fleet driver (`~/Library/LaunchAgents/dev.ailang.mission-world.plist` names `<repo>/tools/launchd/mission-control.sh` by absolute path), and that copy is **11 commits and 705 differing lines behind the fleet** (measured today; the row's own 8/430 figures have grown). The gate whose name is "driver drift" — `scripts/verify_go.sh:205` — implements `D-WORLD-DRIVER-1` as `git status --porcelain -- tools/launchd/`, i.e. **working tree vs THIS repo's HEAD**. A stale-but-*committed* copy is invisible to it *by construction*: both of its legs (`driver_tracked` path-liveness and `driver_drift` working-tree diff) compare the working tree against the local HEAD, and neither can see a stale commit. The gate is **green right now** on a copy 705 lines and 11 commits stale, and its success line reads `✓ driver drift gate: 6 tracked driver files, working tree matches HEAD` — a *true* statement (the working tree does match HEAD) that a reader parses as the *false* claim "the driver is current." This is the same family as the mission's standing finding that *a control only certifies the axis you varied*: the axis varied here is "working tree vs local HEAD," and the axis that matters — "committed copy vs fleet source" — is never measured.

---

## 2. Verification Log

Base commit for every row: `7d3aa00ef880657cc8a769f6726bb700aa54b5e0`. Fleet source: `~/dev/sunholo-data/ailang` at `722e19c77bd7c5c4b7a8a47c9a31909e5387cc9d`.

| ID | Claim | Command | Observed output |
|----|-------|---------|-----------------|
| V1 | launchd runs the repo copy, not the fleet's | `grep -n "mission-control" ~/Library/LaunchAgents/dev.ailang.mission-world.plist` | `14:		<string>/Users/voightkampff/dev/sunholo-data/ailang-world/tools/launchd/mission-control.sh</string>` |
| V2 | Local driver line count | `wc -l tools/launchd/mission-control.sh` | `550 tools/launchd/mission-control.sh` |
| V3 | Fleet driver line count | `wc -l ~/dev/sunholo-data/ailang/tools/launchd/mission-control.sh` | `1075 /Users/voightkampff/dev/sunholo-data/ailang/tools/launchd/mission-control.sh` |
| V4 | Lines present only in the fleet copy | `diff ~/dev/sunholo-data/ailang/tools/launchd/mission-control.sh tools/launchd/mission-control.sh \| grep -c '^<'` | `615` |
| V5 | Lines present only in the local copy | `diff ~/dev/sunholo-data/ailang/tools/launchd/mission-control.sh tools/launchd/mission-control.sh \| grep -c '^>'` | `90` |
| V6 | Commits the local copy does not carry (since 2026-08-15) | `git -C ~/dev/sunholo-data/ailang log --oneline --since=2026-08-15 -- tools/launchd/mission-control.sh \| wc -l` | `11` |
| V7 | Newest of those commits | `git -C ~/dev/sunholo-data/ailang log --oneline --since=2026-08-15 -- tools/launchd/mission-control.sh \| head -1` | `aca6908bd feat(mission-driver): controller fallback is now a chain — codex, then flat-rate GLM-5.3, then its OpenRouter twin (Mark directive 2026-08-31)` |
| V8 | `driver_tracked` (path-liveness leg) | `git ls-files tools/launchd/ scripts/mission_decisions.sh \| wc -l` | `6` |
| V9 | `driver_drift` (working-tree leg) is empty → gate green | `git status --porcelain -- tools/launchd/ scripts/mission_decisions.sh` | *(empty)* |
| V10 | Local `tools/launchd/lib/` does not exist | `ls tools/launchd/lib/` | `ls: tools/launchd/lib/: No such file or directory` |
| V11 | Local driver has no `pin_root` reference | `grep -c 'pin_root' tools/launchd/mission-control.sh` | `0` |
| V12 | Fleet driver has a `pin_root` reference (known-positive control, same file-set) | `grep -c 'pin_root' ~/dev/sunholo-data/ailang/tools/launchd/mission-control.sh` | `1` |
| V13 | `~/.ailang-driver-pin/world` is a worktree of `ailang` (not `ailang-world`), detached | `git -C ~/.ailang-driver-pin/world rev-parse HEAD` and `git -C ~/.ailang-driver-pin/world log -1 --format='%h %ad' --date=short` | `da96b98a5e6fe4ba0a96c1cb76b4b6cd31211f9a` / `da96b98a5 2026-08-26` |
| V14 | Control: `~/.ailang-driver-pin/v1` is current | `git -C ~/.ailang-driver-pin/v1 rev-parse HEAD` and `git -C ~/.ailang-driver-pin/v1 log -1 --format='%h %ad' --date=short` | `f930119aa4adb6181438522cce5c1badf7f1f275` / `f930119aa 2026-09-02` |
| V15 | Nothing in this repo references the pin | `grep -rn 'ailang-driver-pin\|pin-root' tools/ scripts/ .github/` | *(0 hits)* |
| V16 | The env override for the dropped executor id is commented out | `grep -n 'MISSION_EXECUTOR_FALLBACK' ~/.config/ailang/mission-world.env` | `105:# MISSION_EXECUTOR_FALLBACK="${MISSION_EXECUTOR_FALLBACK:-opus}"` |
| V17 | CI runs the gate | `sed -n '160,170p' .github/workflows/ci.yml` | `166:        run: ./scripts/verify_go.sh` |
| V18 | Both repos are public (API fetch is possible) | `gh api repos/sunholo-data/ailang --jq .visibility` ; `gh api repos/sunholo-data/ailang-world --jq .visibility` | `public` / `public` |
| V19 | Fleet HEAD | `git -C ~/dev/sunholo-data/ailang rev-parse HEAD` | `722e19c77bd7c5c4b7a8a47c9a31909e5387cc9d` |
| V20 | Intersection blob comparison: local committed copy vs fleet HEAD, per tracked driver path | `for f in $(git ls-files tools/launchd/ scripts/mission_decisions.sh); do lb=$(git rev-parse HEAD:"$f"); fb=$(git -C ~/dev/sunholo-data/ailang rev-parse HEAD:"$f" 2>/dev/null); [ -n "$fb" ] && { [ "$lb" = "$fb" ] && echo "MATCH  $f" || echo "DIFFER  $f"; }; done` | `MATCH  scripts/mission_decisions.sh`<br>`DIFFER  tools/launchd/derive-planner-lane.sh`<br>`DIFFER  tools/launchd/mission-control.sh`<br>`MATCH  tools/launchd/mission-template.plist`<br>`DIFFER  tools/launchd/test_mission_routing.sh`<br>`MATCH  tools/launchd/testdata/planner-lane/n-no-backtick-bullet.md` |
| V21 | Local blob of `mission-control.sh` | `git rev-parse HEAD:tools/launchd/mission-control.sh` | `c8c92a5c569a51bb649b2a833475641c1e70e5ee` |
| V22 | Fleet blob of `mission-control.sh` | `git -C ~/dev/sunholo-data/ailang rev-parse HEAD:tools/launchd/mission-control.sh` | `253154ad54548b46b521ac20dd23a302c959c077` |
| V23 | The gate is green right now (M4 re-derive) | `driver_tracked=$(git ls-files tools/launchd/ scripts/mission_decisions.sh \| wc -l \| tr -d ' '); echo "$driver_tracked"; driver_drift=$(git status --porcelain -- tools/launchd/ scripts/mission_decisions.sh); echo "[$driver_drift]"` | `6` / `[]` |
| V24 | Gate line numbers in `verify_go.sh` | `grep -n "DRIVER DRIFT GATE\|driver_tracked\|driver_drift\|working tree matches HEAD" scripts/verify_go.sh` | `190:# DRIVER DRIFT GATE — D-WORLD-DRIVER-1, RESOLVED B (Mark, attended 2026-08-17).`<br>`200:driver_tracked=$(git ls-files ...`<br>`201:if [ "$driver_tracked" -lt 5 ]; then`<br>`205:driver_drift=$(git status --porcelain -- tools/launchd/ scripts/mission_decisions.sh)`<br>`206:if [ -n "$driver_drift" ]; then`<br>`212:echo "   ✓ driver drift gate: $driver_tracked tracked driver files, working tree matches HEAD"` |
| V25 | CI does not set `CI` explicitly; GitHub Actions sets `CI=true` by default | `grep -n "CI=" .github/workflows/ci.yml` | *(no explicit `CI=`; GitHub Actions sets `CI=true` for every job by default)* |
| V26 | Sibling-path derivation fails in a worktree (justifies an env-var default, not a derived path) | `FLEET_REPO="$(cd "$(dirname "$(dirname "$0")")/../ailang" 2>/dev/null && pwd)"; echo "$FLEET_REPO"` | *(empty — this worktree is `.wt-world-iter148`, not `ailang-world`, so `../ailang` does not resolve)* |
| V27 | Every local driver path also exists in the fleet (intersection == full local driver set) | `for f in $(git ls-files tools/launchd/ scripts/mission_decisions.sh); do git -C ~/dev/sunholo-data/ailang cat-file -e HEAD:"$f" 2>/dev/null || echo "NO-FLEET: $f"; done` | *(no `NO-FLEET` lines — all 6 local driver paths exist in the fleet)* |
| V28 | Fleet `lib/pin-root.sh` blob (reference for the addition-shaped mutant) | `git -C ~/dev/sunholo-data/ailang rev-parse HEAD:tools/launchd/lib/pin-root.sh` | `5902a60c049f0fac0301eacf56ee578e893e3bc5` |
| V29 | Fleet driver sources the pin helper (known-positive control for V11) | `sed -n '394,396p' ~/dev/sunholo-data/ailang/tools/launchd/mission-control.sh` | `if [ -f "$REPO/tools/launchd/lib/pin-root.sh" ]; then`<br>`. "$REPO/tools/launchd/lib/pin-root.sh"`<br>`pin_root_to_committed_ref "$@"` |
| V30 | Local driver derives `REPO` from `$0` (no pin) | `sed -n '40p' tools/launchd/mission-control.sh` | `REPO="${MISSION_WORKDIR:-$(cd "$(dirname "$0")/../.." && pwd)}"` |
| V31 | `MISSION_WORKDIR` is empty in this fire's environment | `echo "MISSION_WORKDIR=${MISSION_WORKDIR:-<unset/empty>}"` | `MISSION_WORKDIR=<unset/empty>` |
| V32 | World driver path count (the arm's intersection domain) | `git ls-tree -r --name-only HEAD -- tools/launchd scripts/mission_decisions.sh \| sort \| wc -l` | `6` |
| V33 | Fleet driver path count (the authoritative tree) | `git -C ~/dev/sunholo-data/ailang ls-tree -r --name-only HEAD -- tools/launchd scripts/mission_decisions.sh \| sort \| wc -l` | `48` |
| V34 | MISSING LOCALLY (fleet paths absent locally) | `comm -13 <(git ls-tree -r --name-only HEAD -- tools/launchd scripts/mission_decisions.sh \| sort) <(git -C ~/dev/sunholo-data/ailang ls-tree -r --name-only HEAD -- tools/launchd scripts/mission_decisions.sh \| sort) \| wc -l` | `42` |
| V35 | MISSING IN FLEET (local paths absent in fleet) | `comm -23 <(git ls-tree -r --name-only HEAD -- tools/launchd scripts/mission_decisions.sh \| sort) <(git -C ~/dev/sunholo-data/ailang ls-tree -r --name-only HEAD -- tools/launchd scripts/mission_decisions.sh \| sort) \| wc -l` | `0` |
| V36 | Union of both trees | `cat <(git ls-tree -r --name-only HEAD -- tools/launchd scripts/mission_decisions.sh \| sort) <(git -C ~/dev/sunholo-data/ailang ls-tree -r --name-only HEAD -- tools/launchd scripts/mission_decisions.sh \| sort) \| sort -u \| wc -l` | `48` |
| V37 | The 42 MISSING LOCALLY are overwhelmingly files World must NOT carry | `comm -13 <(git ls-tree -r --name-only HEAD -- tools/launchd scripts/mission_decisions.sh \| sort) <(git -C ~/dev/sunholo-data/ailang ls-tree -r --name-only HEAD -- tools/launchd scripts/mission_decisions.sh \| sort)` | other missions' plists (`dev.ailang.mission-motoko.plist`, `dev.ailang.mission-docs.plist`, `dev.ailang.mission-recovery*.plist`, `dev.ailang.nightly-eval.plist`, `dev.ailang.rig-watchdog.plist`, `dev.ailang.server.plist`, `dev.ollama.serve.plist`, `dev.ailang.os-rotation-filler.plist`, `dev.ailang.coordinator.plist.template`), other missions' env files (`mission-env/mission-v1.env`, `mission-docs.env`, `mission-motoko.env`), fleet-only scripts (`nightly-eval.sh`, `nightly-lang-eval.sh`, `rig-watchdog.sh`, `rig-lock.sh`, `os-rotation-filler.sh`, `mission-recovery.sh`, `install_coordinator.sh`, five `test_*.sh`), 12 planner-lane testdata fixtures |
| V38 | `lib/pin-root.sh` is among the 42 MISSING LOCALLY (the one World genuinely SHOULD carry) | `comm -13 <(git ls-tree -r --name-only HEAD -- tools/launchd scripts/mission_decisions.sh \| sort) <(git -C ~/dev/sunholo-data/ailang ls-tree -r --name-only HEAD -- tools/launchd scripts/mission_decisions.sh \| sort) \| grep pin-root` | `tools/launchd/lib/pin-root.sh` |
| V39 | CB1: a bare function call under `set -euo pipefail` exits immediately on non-zero, so `fleet_rc=$?` is never reached | `set -euo pipefail; f(){ return 2; }; echo before; f; rc=$?; echo "REACHED rc=$rc"` | `before` then `script exit=2` — `REACHED rc=2` is never printed |
| V40 | CB1 control: with `set -uo pipefail` the return code is captured | `set -uo pipefail; f(){ return 2; }; echo before; f; rc=$?; echo "REACHED rc=$rc"` | `before` / `REACHED rc=2` / `script exit=0` |
| V41 | `--evidence-manifest-check` mode precedes the `AILANG_BIN` gate, so an isolated mode needs no `AILANG_BIN` | `grep -n "evidence-manifest-check\|AILANG_BIN is unset" scripts/verify_go.sh` | `111:if [ "${1:-}" = "--evidence-manifest-check" ]` / `120:if [ -z "${AILANG_BIN:-}" ]` |
| V42 | CB1-bis: the ISOLATED-MODE block exactly as round 2 drafted it (bare call) loses the CI loud skip | harness: `set -euo pipefail`; `check_driver_fleet(){ echo SKIPPED; return 2; }`; `check_driver_fleet; rc=$?; [ "$rc" -eq 2 ] && rc=0; echo "REACHED-MAP: rc=$rc"; exit "$rc"` — status captured without a pipe | `SKIPPED` printed, **`REACHED-MAP` never printed**, script `rc=2`. The 2→0 mapping is dead code. |
| V43 | CB1-bis control/fix arm: the `if` form captures it | same harness, `rc=0; if check_driver_fleet; then :; else rc=$?; fi; [ "$rc" -eq 2 ] && rc=0; echo "REACHED-MAP: rc=$rc"; exit "$rc"` | `SKIPPED` / `REACHED-MAP: rc=0` / script `rc=0`. Asserted `rc_bare != rc_if` explicitly → `2 != 0`, DIFFER-OK, so this is not a false symmetry. |
| V44 | REFUTED controller hypothesis: a top-level `[ "$rc" -eq 2 ] && rc=0` does NOT trip `set -e` when the test is false | `set -euo pipefail; rc=1; [ "$rc" -eq 2 ] && rc=0; echo "REACHED-AFTER-AND: rc=$rc"; exit "$rc"` | `REACHED-AFTER-AND: rc=1`, script `rc=1`. There is ONE `set -e` defect in the isolated mode, not two — recorded so no later round re-litigates it. |

**AC baselines (run on the UNMODIFIED tree at `7d3aa00`):**

| AC | Baseline result on unmodified tree |
|----|-----------------------------------|
| AC1 (rig staleness detection) | Gate exits 0, prints `working tree matches HEAD` (V23/V24); no `DRIVER DRIFT vs FLEET` — the defect. AC is red at base by design (the fix is absent). |
| AC2 (CI loud skip) | Gate exits 0, prints `working tree matches HEAD` (V23/V24); no `SKIPPED` line. AC is red at base by design. |
| AC3 (rig typed refusal when fleet absent) | Gate exits 0 (V23); `AILANG_FLEET_REPO` has no effect because the arm does not exist. AC is red at base by design. |
| AC4 (path-liveness control) | Gate exits 0 (V23); no "0 comparable files" refusal. AC is red at base by design. |
| AC5 (REQUIRED path absent locally) | Gate exits 0 (V23); no `MISSING LOCALLY` line. AC is red at base by design. |
| AC6 (MISSING IN FLEET, synthetic) | Gate exits 0 (V23); no `MISSING IN FLEET` line. AC is red at base by design. |
| AC7 (CI branch exits ZERO) | Gate exits 0 (V23); the arm does not exist. AC is red at base by design. |
| AC8 (regression guard) | Existing working-tree arm fires on dirt (code present at V24); AC passes at base and must still pass after the fix. |

---

## 3. Options considered

The central decision (C3): `verify_go.sh` runs in CI (V17) where **the fleet checkout does not exist** (CI checks out `ailang-world` only). A gate that hard-fails when the fleet source is unreachable turns CI permanently red; a gate that silently passes when the fleet source is unreachable is exactly the vacuous pass this row exists to close.

**A measured constraint that shaped the choice (quorum round 1):** the fleet tree is far larger than World's. At HEAD `7d3aa00` vs fleet `722e19c7`, World tracks **6** driver paths, the fleet carries **48**, the union is **48**, MISSING LOCALLY is **42**, and MISSING IN FLEET is **0** (V32–V36). The 42 MISSING LOCALLY are overwhelmingly files World must NOT carry — other missions' plists and env files, fleet-only scripts, and 12 planner-lane testdata fixtures (V37). Only one of the 42, `tools/launchd/lib/pin-root.sh`, is a file World genuinely SHOULD carry (V38; it is the pin mechanism, Finding 1). This rules out any design that requires "both blobs exist and match over the union": that would red permanently on 42 correctly-absent files, and World cannot fix it (frozen core, C1) — the exact permanent-red failure mode Option C is rejected for. The design must therefore compare the intersection for content drift, add an explicit REQUIRED set for the few fleet paths World must carry, and treat the remaining fleet-only paths as a loud, counted, non-fatal residual.

### Option A — Committed expected-hash manifest (embedded in `verify_go.sh`)
- **Mechanism:** hardcode the fleet HEAD blob hash for each frozen-core path in `verify_go.sh`; CI compares the local committed blob against the recorded hash.
- **Consequence:** CI gets a non-vacuous check without the fleet checkout, but the manifest is a *snapshot* that goes stale the moment the fleet advances, and a stale manifest is the same vacuous pass this row exists to close — it just moves the self-comparison to a snapshot-comparison.
- **Cost:** must be re-derived and re-committed on every fleet commit (editing `verify_go.sh`); ~0.05d per update, plus a new staleness failure mode.
- **What it does NOT catch:** a fleet advance not yet recorded in the manifest.

### Option B — Rig-only live arm, loud when skipped, typed when it cannot run
- **Mechanism:** an arm that compares the local committed blob against the fleet checkout's HEAD blob for the same path, where the fleet checkout exists (the rig). When the fleet checkout is absent: loud-skip in CI (`CI=true`), typed refusal on the rig (not CI).
- **Consequence:** catches the real exposure where the driver actually runs; honest in CI (a loud skip, never a false "current" claim); no network dependency, no manifest maintenance, no permanent red in CI.
- **Cost:** ~0.2d. Depends on the fleet checkout being present and current on the rig.
- **What it does NOT catch:** CI cannot certify driver currency (the arm is loud-skipped there); a fleet advance between rig runs is caught at the next rig run.

### Option C — Fetch the fleet blob over the public API
- **Mechanism:** `gh api repos/sunholo-data/ailang/contents/...` (V18: public) to fetch the live fleet blob and compare in CI and on the rig.
- **Consequence:** genuinely live everywhere, but reds *permanently* whenever the fleet advances — and World cannot fix it (frozen core, C1); the red would mean "the fleet must commit," but it would fire on every fleet commit even before the fleet lands the update in World, and it depends on network/API availability.
- **Cost:** ~0.15d; network dependency, rate limits, and a permanent-red failure mode.
- **What it does NOT catch:** nothing about staleness — but it is unusable as a hard gate because it reds on every fleet advance.

### Option D — Combination (rig arm + manifest)
- **Mechanism:** Option B for the real check on the rig, plus Option A's manifest for CI.
- **Consequence:** best coverage, but adds manifest maintenance and reintroduces the manifest-staleness residual in CI.
- **Cost:** ~0.3d (at the top of scope).
- **What it does NOT catch:** a fleet advance not yet recorded in the manifest (CI leg).

### Choice: **Option B — rig-only live arm, loud when skipped, typed when it cannot run.**

Reason: it is the row's own suggested mechanism ("assert the copy's content-hash equals the fleet checkout's HEAD blob for the same path, with a loud, typed refusal when the fleet checkout is absent"), it runs where the driver runs (the rig), it is honest in CI (a loud skip, never a false "current" claim), it does not reintroduce a staleness failure mode (unlike the manifest), and it does not red permanently on fleet advance in CI (unlike the API). The manifest (A) is rejected because it reintroduces the exact vacuous-pass failure it claims to close; the API (C) is rejected because it reds permanently on fleet advance and World cannot fix it (frozen core). The union measurement (V32–V38) is the reason the arm's domain is the intersection **plus** an explicit REQUIRED set **plus** a non-fatal unclassified report, rather than the full union — see §4.1.

**Residual of the choice:** CI cannot certify driver currency — the fleet arm is loud-skipped there, so CI's green is honest but does not assert the driver is current; a fleet advance between rig runs is caught at the next rig run; the arm depends on the fleet checkout being present and current on the rig; and unclassified fleet-only paths (41 today, §7) are reported but not certified.

---

## 4. Design

The change is confined to `scripts/verify_go.sh` (C2). The existing working-tree-vs-HEAD arm (V24, lines 200–212) and its `driver_tracked` path-liveness control are **kept unchanged** (C4); the new arm is **in addition** to them.

### 4.1 New function `check_driver_fleet()`

Add a function near the top of the script (after `check_evidence_manifest`), and a `--driver-fleet-check` isolated mode beside the existing `--evidence-manifest-check` mode (V41: that mode is at line 111, before the `AILANG_BIN` gate at line 120, so an isolated mode needs no `AILANG_BIN`):

```bash
# FLEET-COMPARISON ARM — D-WORLD-DRIVER-1, iter-148 round 2. The working-tree-vs-HEAD
# arm cannot see a stale-but-COMMITTED copy (it compares the copy to itself). This arm
# compares the committed copy against the FLEET source, which is where the driver
# actually lives. The driver is FLEET-owned; World detects and reports, the fleet
# commits. This arm is World's own file and is in scope.
#
# DOMAIN (quorum round 1): the arm compares (a) every path World tracks under
# tools/launchd/ and scripts/mission_decisions.sh, and (b) every path in the EXPLICIT
# REQUIRED_FLEET_PATHS set below. It does NOT require every fleet path to exist
# locally: the fleet carries 48 paths under this prefix, World legitimately carries 6,
# and 42 of the fleet's paths are other missions' plists/envs, fleet-only scripts, and
# testdata fixtures World must NOT carry (V32–V37). Requiring all 48 would red
# permanently and World cannot fix it (frozen core, C1). Instead, a fleet-only path
# that is neither tracked by World nor in REQUIRED_FLEET_PATHS is UNCLASSIFIED:
# reported loudly and counted, but non-fatal. A path becomes REQUIRED by deliberate,
# reviewed addition to the set below — the classification gate. Until classified, a
# fleet-only path is a counted residual, never a silent pass and never a permanent red.
#
# The fleet source is a checkout of sunholo-data/ailang. On the rig it is a sibling
# of this repo; in CI it does not exist (CI checks out ailang-world only). The path
# is configurable via AILANG_FLEET_REPO and defaults to the known rig path.
AILANG_FLEET_REPO="${AILANG_FLEET_REPO:-$HOME/dev/sunholo-data/ailang}"

# REQUIRED_FLEET_PATHS: fleet paths World is REQUIRED to carry. Each entry is a
# deliberate, reviewed commitment; a required path absent locally is a typed
# MISSING_LOCALLY failure. First member: tools/launchd/lib/pin-root.sh — the pin
# mechanism the fleet driver sources (V29) and World's copy lacks (V10–V12, Finding 1).
# Add a new fleet path here only when World genuinely must carry it; until then it is
# an unclassified fleet-only path (counted, non-fatal).
REQUIRED_FLEET_PATHS=(
  "tools/launchd/lib/pin-root.sh"
)

check_driver_fleet() {
  # Returns 0 (green), 1 (FATAL/typed refusal), or 2 (loud skip in CI).
  if [ -d "$AILANG_FLEET_REPO" ] && git -C "$AILANG_FLEET_REPO" rev-parse --git-dir >/dev/null 2>&1; then
    fleet_head="$(git -C "$AILANG_FLEET_REPO" rev-parse HEAD)"
    compared=0
    differing=""
    missing_in_fleet=""
    missing_locally=""
    unclassified=0

    # Phase 1 — every path World tracks under the driver prefix, EXCEPT REQUIRED paths
    # (those are owned by Phase 2, so a REQUIRED path is never double-counted). A local
    # path the fleet lacks is a typed MISSING_IN_FLEET failure (never a silent continue).
    while IFS= read -r path; do
      [ -z "$path" ] && continue
      required=0
      for rp in "${REQUIRED_FLEET_PATHS[@]}"; do
        [ "$rp" = "$path" ] && required=1 && break
      done
      [ "$required" -eq 1 ] && continue   # REQUIRED paths are owned by Phase 2
      local_blob="$(git rev-parse --verify "HEAD:$path" 2>/dev/null || true)"
      fleet_blob="$(git -C "$AILANG_FLEET_REPO" rev-parse --verify "HEAD:$path" 2>/dev/null || true)"
      if [ -z "$fleet_blob" ]; then
        missing_in_fleet="$missing_in_fleet
  $path (tracked by World, absent in fleet)"
        continue
      fi
      compared=$((compared + 1))
      if [ "$local_blob" != "$fleet_blob" ]; then
        differing="$differing
  $path (local $local_blob != fleet $fleet_blob)"
      fi
    done < <(git ls-files tools/launchd/ scripts/mission_decisions.sh)

    # Phase 2 — REQUIRED_FLEET_PATHS. A required path absent locally is a typed
    # MISSING_LOCALLY failure; absent in fleet is MISSING_IN_FLEET; present-but-
    # differing is a DIFFERING failure.
    for path in "${REQUIRED_FLEET_PATHS[@]}"; do
      if ! git cat-file -e "HEAD:$path" 2>/dev/null; then
        missing_locally="$missing_locally
  $path (REQUIRED by World, absent locally)"
        continue
      fi
      local_blob="$(git rev-parse --verify "HEAD:$path" 2>/dev/null || true)"
      fleet_blob="$(git -C "$AILANG_FLEET_REPO" rev-parse --verify "HEAD:$path" 2>/dev/null || true)"
      if [ -z "$fleet_blob" ]; then
        missing_in_fleet="$missing_in_fleet
  $path (REQUIRED by World, absent in fleet)"
        continue
      fi
      compared=$((compared + 1))
      if [ "$local_blob" != "$fleet_blob" ]; then
        differing="$differing
  $path (local $local_blob != fleet $fleet_blob)"
      fi
    done

    # Phase 3 — unclassified fleet-only paths: loud, counted, non-fatal. These are
    # fleet paths World neither tracks nor requires (other missions' files, fleet-only
    # scripts, testdata). They are reported so the residual is visible, but they do not
    # red the gate — requiring them would red permanently on 42 correctly-absent files.
    while IFS= read -r path; do
      [ -z "$path" ] && continue
      if git cat-file -e "HEAD:$path" 2>/dev/null; then
        continue   # World tracks it (Phase 1)
      fi
      required=0
      for rp in "${REQUIRED_FLEET_PATHS[@]}"; do
        [ "$rp" = "$path" ] && required=1 && break
      done
      [ "$required" -eq 1 ] && continue   # REQUIRED (Phase 2)
      unclassified=$((unclassified + 1))
      echo "   ⚠ unclassified fleet-only path (not tracked, not required): $path" >&2
    done < <(git -C "$AILANG_FLEET_REPO" ls-tree -r --name-only HEAD -- tools/launchd scripts/mission_decisions.sh)

    # Refusal order (quorum round 1): path-liveness first (an empty comparison must
    # never read as green), then DIFFERING, then MISSING_IN_FLEET, then MISSING_LOCALLY.
    # MISSING_LOCALLY is deliberately LAST so the branch-specific ACs (AC1/AC4/AC6) can
    # reach their branches against the current World tree, which lacks lib/pin-root.sh
    # (V10) and would otherwise trip MISSING_LOCALLY first and mask the branch under test.
    # At base the real gate reports DIFFERING (the stale driver, the row's core) first;
    # MISSING_LOCALLY (pin-root, Finding 1) surfaces on the next run after the driver is
    # landed. Any of the four is a hard red, so the order only affects which message leads.
    if [ "$compared" -eq 0 ]; then
      echo "verify_go.sh: FATAL: fleet-comparison arm enumerated 0 comparable driver files against $AILANG_FLEET_REPO; the instrument is broken, so every 'matches fleet' result is void" >&2
      return 1
    fi

    if [ -n "$differing" ]; then
      echo "verify_go.sh: FATAL: DRIVER DRIFT vs FLEET (D-WORLD-DRIVER-1) — the committed copy differs from fleet HEAD $fleet_head:" >&2
      printf '%s\n' "$differing" >&2
      echo "  The driver is fleet-owned; land the current driver as a fleet-authored commit. World's controller must not edit or absorb it." >&2
      return 1
    fi
    if [ -n "$missing_in_fleet" ]; then
      echo "verify_go.sh: FATAL: DRIVER DRIFT vs FLEET (D-WORLD-DRIVER-1) — World-tracked paths MISSING IN FLEET:" >&2
      printf '%s\n' "$missing_in_fleet" >&2
      echo "  A World-tracked driver path is absent from the fleet; reconcile which tree owns it." >&2
      return 1
    fi
    if [ -n "$missing_locally" ]; then
      echo "verify_go.sh: FATAL: DRIVER DRIFT vs FLEET (D-WORLD-DRIVER-1) — REQUIRED fleet paths MISSING LOCALLY:" >&2
      printf '%s\n' "$missing_locally" >&2
      echo "  The driver is fleet-owned; land the required file as a fleet-authored commit. World's controller must not edit or absorb it." >&2
      return 1
    fi
    echo "   ✓ fleet-comparison arm: $compared tracked frozen-core files match fleet HEAD $fleet_head — tracked copy is current (untracked fleet additions not certified)"
    if [ "$unclassified" -gt 0 ]; then
      echo "   ⚠ $unclassified unclassified fleet-only paths not certified (see above)" >&2
    fi
    return 0
  fi
  if [ -n "${CI:-}" ]; then
    # CI: no fleet checkout. Loud skip — this is NOT a certification.
    echo "   ⚠ fleet-comparison arm SKIPPED (fleet checkout absent at $AILANG_FLEET_REPO) — driver currency NOT certified here"
    return 2
  fi
  if [ -z "${CI:-}" ]; then
    # Rig, fleet source absent: an unreachable source must never read as agreement.
    echo "verify_go.sh: FATAL: DRIVER DRIFT (D-WORLD-DRIVER-1) — fleet source $AILANG_FLEET_REPO is absent; the fleet-comparison arm cannot run, so driver currency is NOT certified" >&2
    return 1
  fi
  return 0
}

# CB1 APPLIES HERE TOO — this is the SIBLING call site, and quorum round 2 rejected
# the bare form unanimously (gpt5-6-sol, gemini-3-1-pro, oc-glm-5-2). Under
# `set -euo pipefail` a bare `check_driver_fleet` returning 2 exits the script before
# `rc=$?` runs, so the 2 -> 0 mapping below is dead code and the CI loud skip becomes a
# hard rc=2. Measured, two arms: bare form -> rc=2, mapping line never reached; if-form
# -> rc=0, mapping reached (V42/V43). The `if` form is exempt from `set -e`.
# The trailing `[ ... ] && rc=0` is safe at top level: measured, rc=1 does NOT trip
# `set -e` there (V44).
if [ "${1:-}" = "--driver-fleet-check" ]; then
  rc=0
  if check_driver_fleet; then
    :
  else
    rc=$?
  fi
  # A loud skip (2) is a non-failure by design (CI must stay green); only a typed
  # refusal (1) is a hard failure. Map 2 -> 0 so the isolated mode mirrors the main
  # flow's non-fatal handling of the CI branch (CB1).
  [ "$rc" -eq 2 ] && rc=0
  exit "$rc"
fi
```

### 4.2 Call site in the main flow

Immediately after the existing driver drift gate's success line (V24, line 212), invoke the new arm and propagate a hard failure. **The call site must be `set -e`-safe (CB1):** `verify_go.sh` runs under `set -euo pipefail` from line 16, and `set +e` appears only at lines 247/263/279 — all after this call site. A *bare* `check_driver_fleet` call that returns non-zero would exit the script immediately (V39), so `fleet_rc=$?` would never run and the CI loud-skip (return 2) would turn CI permanently red — the exact outcome Option B was chosen over Option C to avoid. The `if` form is exempt from `set -e`, so the return code is captured and only a typed refusal (1) is a hard failure:

```bash
fleet_rc=0
if check_driver_fleet; then
  :
else
  fleet_rc=$?
fi
if [ "$fleet_rc" -eq 1 ]; then
  exit 1
fi
```

### 4.3 Success-line honesty

Change the existing success line (V24, line 212) so it cannot be read as a stronger claim than it makes. The current line's defect is precisely that "working tree matches HEAD" reads as "the driver is current." The new line labels the working-tree arm and lets the fleet arm's result stand beside it:

```bash
echo "   ✓ driver drift gate: $driver_tracked tracked driver files, working tree matches HEAD (working-tree arm)"
```

The fleet arm's success line is adopted **verbatim** from the quorum (oc-glm-5-2) and never claims more than the arm measured — it scopes "current" to the *tracked* copy and explicitly disclaims untracked fleet additions:

```bash
echo "   ✓ fleet-comparison arm: $compared tracked frozen-core files match fleet HEAD $fleet_head — tracked copy is current (untracked fleet additions not certified)"
```

The full driver-drift output is then, on the rig (fleet present, all match — the green state after the fleet lands the current driver including `lib/pin-root.sh`):

```
   ✓ driver drift gate: 6 tracked driver files, working tree matches HEAD (working-tree arm)
   ✓ fleet-comparison arm: 7 tracked frozen-core files match fleet HEAD 722e19c7 — tracked copy is current (untracked fleet additions not certified)
   ⚠ 41 unclassified fleet-only paths not certified (see above)
```

and in CI (fleet absent):

```
   ✓ driver drift gate: 6 tracked driver files, working tree matches HEAD (working-tree arm)
   ⚠ fleet-comparison arm SKIPPED (fleet checkout absent at /home/runner/dev/sunholo-data/ailang) — driver currency NOT certified here
```

The CI line cannot be parsed as "the driver is current" — it says the opposite. The rig line ties "current" to a specific fleet HEAD sha and to the *tracked* copy only, so it is scoped to the measured fleet state and reds the moment the fleet advances or a required path is missing. The bare phrase "driver is current" appears nowhere in the output or in §4.4's table.

### 4.4 Refusal branches and typed exit behaviour

| Branch | Condition | Behaviour | Exit |
|--------|-----------|-----------|------|
| Fleet present, 0 comparable files | `[ "$compared" -eq 0 ]` (checked first) | `FATAL: ... enumerated 0 comparable driver files ... instrument is broken` | 1 |
| Fleet present, any blob differs | `[ "$local_blob" != "$fleet_blob" ]` (Phase 1 or Phase 2) | `FATAL: DRIVER DRIFT vs FLEET` naming each differing path | 1 |
| Fleet present, World-tracked path absent in fleet | `[ -z "$fleet_blob" ]` (Phase 1 or Phase 2) | `FATAL: ... MISSING IN FLEET` naming each | 1 |
| Fleet present, REQUIRED path absent locally | `! git cat-file -e "HEAD:$path"` (Phase 2) | `FATAL: ... REQUIRED fleet paths MISSING LOCALLY` naming each | 1 |
| Fleet present, all match | all blobs equal, no missing | `✓ fleet-comparison arm: N tracked frozen-core files match fleet HEAD <sha> — tracked copy is current (untracked fleet additions not certified)` | 0 |
| Fleet present, unclassified fleet-only paths | fleet path not tracked, not required (Phase 3) | `⚠ unclassified fleet-only path: <path>` (counted, non-fatal) | 0 (with success) |
| Fleet absent, CI set | `[ -n "${CI:-}" ]` | `⚠ fleet-comparison arm SKIPPED ... NOT certified here` | 2 (not a failure) |
| Fleet absent, not CI (rig) | `[ -z "${CI:-}" ]` | `FATAL: DRIVER DRIFT ... fleet source ... absent ... NOT certified` | 1 |

---

## 5. Acceptance criteria

Each AC is a command that can fail. Baselines on the unmodified tree are recorded in §2. AC1–AC7 are red at base by design (the fix is absent); AC8 is a regression guard. **AC1–AC7 drive the isolated `--driver-fleet-check` mode** (CB2) so they do not pay the full `go build`/`go test`/`-race` cost of `./scripts/verify_go.sh` (lines 272–297); the mode is before the `AILANG_BIN` gate (V41), so no `AILANG_BIN` is needed. AC8 exercises the existing working-tree arm, which lives only in the main flow, so it uses the full script.

1. **AC1 — content drift (DIFFERING).** A controlled fleet repo where a World-tracked file differs. Command:
   ```bash
   tmp=$(mktemp -d); git -C "$tmp" init -q
   for f in $(git ls-files tools/launchd/ scripts/mission_decisions.sh); do
     mkdir -p "$tmp/$(dirname "$f")"; git show "HEAD:$f" > "$tmp/$f"
   done
   echo "# fleet-only drift" >> "$tmp/tools/launchd/mission-control.sh"
   git -C "$tmp" add -A && git -C "$tmp" commit -qm fleet
   AILANG_FLEET_REPO="$tmp" ./scripts/verify_go.sh --driver-fleet-check
   ```
   Assert: `exit != 0` **and** output contains `DRIVER DRIFT vs FLEET` naming `tools/launchd/mission-control.sh`. Baseline: gate exits 0, prints `working tree matches HEAD` (V23/V24). *What would this still pass under if the claim were false?* A gate that compared the copy to itself (the current defect) — which is exactly the baseline.

2. **AC2 — CI loud skip, honest.** In CI (no fleet checkout), the arm exits 0 but prints `SKIPPED ... NOT certified` and does **not** print `driver is current`. Command: `CI=true AILANG_FLEET_REPO=/nonexistent ./scripts/verify_go.sh --driver-fleet-check`. Assert: `exit == 0` **and** output contains `SKIPPED` **and** output does not contain `driver is current`. Baseline: gate exits 0, prints `working tree matches HEAD` with no `SKIPPED` (V23/V24). *What would this still pass under if the claim were false?* A gate that silently passed when the fleet source is unreachable — which is the vacuous pass this row closes; the AC fails at base because the `SKIPPED` line is absent.

3. **AC3 — rig typed refusal when fleet source absent.** On the rig (not CI) with the fleet source absent, the arm exits non-zero with a typed refusal, not a pass. Command: `AILANG_FLEET_REPO=/nonexistent ./scripts/verify_go.sh --driver-fleet-check` (CI unset). Assert: `exit != 0` **and** output contains `fleet source ... absent`. Baseline: gate exits 0 (V23); `AILANG_FLEET_REPO` has no effect because the arm does not exist. *What would this still pass under if the claim were false?* A gate that treated an unreachable source as agreement — the exact failure the row names; the AC fails at base because no refusal is emitted.

4. **AC4 — path-liveness control.** The arm refuses when it enumerates zero comparable files. Command:
   ```bash
   tmp=$(mktemp -d); git -C "$tmp" init -q && git -C "$tmp" commit -qm empty --allow-empty
   AILANG_FLEET_REPO="$tmp" ./scripts/verify_go.sh --driver-fleet-check
   ```
   Assert: `exit != 0` **and** output contains `0 comparable driver files`. Baseline: gate exits 0 (V23); the arm does not exist. *What would this still pass under if the claim were false?* An enumerator that sees nothing and reports green — the non-vacuity rule; the AC fails at base because no refusal is emitted.

5. **AC5 — REQUIRED path absent locally (MISSING_LOCALLY).** A controlled fleet repo that carries `tools/launchd/lib/pin-root.sh` (REQUIRED) while the current World tree lacks it (V10). Command:
   ```bash
   tmp=$(mktemp -d); git -C "$tmp" init -q
   for f in $(git ls-files tools/launchd/ scripts/mission_decisions.sh); do
     mkdir -p "$tmp/$(dirname "$f")"; git show "HEAD:$f" > "$tmp/$f"
   done
   mkdir -p "$tmp/tools/launchd/lib"; echo "# pin-root" > "$tmp/tools/launchd/lib/pin-root.sh"
   git -C "$tmp" add -A && git -C "$tmp" commit -qm fleet
   AILANG_FLEET_REPO="$tmp" ./scripts/verify_go.sh --driver-fleet-check
   ```
   Assert: `exit != 0` **and** output contains `MISSING LOCALLY` naming `tools/launchd/lib/pin-root.sh`. Baseline: gate exits 0 (V23); no `MISSING LOCALLY` line. *What would this still pass under if the claim were false?* An intersection-only arm that never looks at fleet-only paths — the exact gap gpt5-6-sol flagged; the AC fails at base because no `MISSING LOCALLY` is emitted.

6. **AC6 — MISSING IN FLEET (synthetic).** A controlled fleet repo missing one World-tracked path. Synthetic because MISSING IN FLEET is **0** today (V35), so the branch is empty in the wild and must be exercised by construction. Command:
   ```bash
   tmp=$(mktemp -d); git -C "$tmp" init -q
   for f in $(git ls-files tools/launchd/ scripts/mission_decisions.sh); do
     [ "$f" = "tools/launchd/test_mission_routing.sh" ] && continue   # drop one tracked path
     mkdir -p "$tmp/$(dirname "$f")"; git show "HEAD:$f" > "$tmp/$f"
   done
   git -C "$tmp" add -A && git -C "$tmp" commit -qm fleet
   AILANG_FLEET_REPO="$tmp" ./scripts/verify_go.sh --driver-fleet-check
   ```
   Assert: `exit != 0` **and** output contains `MISSING IN FLEET` naming `tools/launchd/test_mission_routing.sh`. Baseline: gate exits 0 (V23); no `MISSING IN FLEET` line. *What would this still pass under if the claim were false?* A silent `[ -z "$fleet_blob" ] && continue` that skips a local path absent from the fleet — the exact silent-fallback gemini-3-1-pro flagged; the AC fails at base because no `MISSING IN FLEET` is emitted.

7. **AC7 — CI branch exits ZERO (CB1 guard).** The CI loud-skip branch (return 2) must never turn the gate red. Command: `CI=true AILANG_FLEET_REPO=/nonexistent ./scripts/verify_go.sh --driver-fleet-check`. Assert: `exit == 0`. Baseline: gate exits 0 (V23); the arm does not exist. *What would this still pass under if the claim were false?* A call site that propagates the CI skip as a hard failure (e.g. a bare `check_driver_fleet` under `set -e`, or `exit $?` in the isolated mode) — which would turn CI permanently red; the AC fails if the CI branch ever exits non-zero.

8. **AC8 — regression guard (existing arm preserved).** The existing working-tree-vs-HEAD arm still fires on working-tree dirt. Command: `touch tools/launchd/mission-control.sh && AILANG_BIN=<pinned> ./scripts/verify_go.sh` (then restore). Assert: `exit != 0` **and** output contains the existing `DRIVER DRIFT` message. Baseline: the existing arm fires (V24). This AC passes at base and must still pass after the fix; it is a regression guard, not a test of the new arm, and it uses the full script because the working-tree arm lives only in the main flow.

---

## 6. Test and mutation plan

For every refusal branch, one neutering mutation. Neuter with `if false && <cond>`, never by deleting the block. Each row states which write the assertion reads. The addition-shaped mutant is re-shaped per gpt5-6-sol's catch: the round-1 mutant added the file **locally first**, so it never tested detection of a **fleet-only** path; the round-2 mutant adds the path to the **fleet** and asserts the arm detects it.

| Mutation | File | Test that kills it | Why the observable is DOWNSTREAM of the mechanism |
|----------|------|--------------------|--------------------------------------------------|
| Neuter DIFFERING (Phase 1): `if false && [ "$local_blob" != "$fleet_blob" ]` | `scripts/verify_go.sh` | AC1 (controlled fleet repo with a differing tracked file): assert `exit != 0` and output contains `DRIVER DRIFT vs FLEET`. With the mutation the FATAL never fires, the arm reports green, and the test fails. | The assertion reads the DIFFERING write; the write is produced only by the `!=` branch. Neutering the branch removes the write, so the assertion cannot pass. |
| Neuter MISSING_IN_FLEET (Phase 1): `if false && [ -z "$fleet_blob" ]` | `scripts/verify_go.sh` | AC6 (controlled fleet repo missing a tracked path): assert `exit != 0` and output contains `MISSING IN FLEET`. With the mutation the branch never fires, the path is silently skipped, the arm reports green, and the test fails. | The assertion reads the MISSING_IN_FLEET write; the write is produced only by the `-z` branch. Neutering it removes the write. |
| Neuter MISSING_LOCALLY (Phase 2): `if false && ! git cat-file -e "HEAD:$path"` | `scripts/verify_go.sh` | AC5 (controlled fleet repo with `lib/pin-root.sh`, World lacks it): assert `exit != 0` and output contains `MISSING LOCALLY` naming `pin-root.sh`. With the mutation the branch never fires, the arm reports green, and the test fails. | The assertion reads the MISSING_LOCALLY write; the write is produced only by the required-absent branch. Neutering it removes the write. |
| Neuter MISSING_IN_FLEET (Phase 2): `if false && [ -z "$fleet_blob" ]` | `scripts/verify_go.sh` | AC6 variant where a REQUIRED path is absent in fleet: assert `exit != 0` and output contains `MISSING IN FLEET`. With the mutation the branch never fires, the arm reports green, and the test fails. | The assertion reads the MISSING_IN_FLEET write; produced only by the `-z` branch in Phase 2. |
| Neuter DIFFERING (Phase 2): `if false && [ "$local_blob" != "$fleet_blob" ]` | `scripts/verify_go.sh` | AC1 variant where a REQUIRED path differs: assert `exit != 0` and output contains `DRIVER DRIFT vs FLEET`. With the mutation the branch never fires, the arm reports green, and the test fails. | The assertion reads the DIFFERING write; produced only by the `!=` branch in Phase 2. |
| Neuter the path-liveness refusal: `if false && [ "$compared" -eq 0 ]` | `scripts/verify_go.sh` | AC4 (empty fleet repo): assert `exit != 0` and output contains `0 comparable driver files`. With the mutation the refusal never fires, the arm reports green on an empty comparison, and the test fails. | The assertion reads the path-liveness write; the write is produced only by the `compared -eq 0` branch. Neutering it removes the write. |
| Neuter the CI loud skip: `if false && [ -n "${CI:-}" ]` | `scripts/verify_go.sh` | AC2 (`CI=true`, fleet absent): assert `exit == 0` and output contains `SKIPPED`. With the mutation the skip never fires, the arm falls through to the rig refusal (exit 1), and the test fails. | The assertion reads the skip write; the write is produced only by the absent-and-CI branch. Neutering it removes the write. |
| Neuter the rig typed refusal: `if false && [ -z "${CI:-}" ]` | `scripts/verify_go.sh` | AC3 (fleet absent, not CI): assert `exit != 0` and output contains `fleet source ... absent`. With the mutation the refusal never fires, the function returns 0, and the test fails. | The assertion reads the rig-refusal write; the write is produced only by the absent-and-not-CI branch. Neutering it removes the write. |
| **Addition-shaped mutant (fleet-only path detection)** — add `tools/launchd/lib/new-helper.sh` to the **fleet** (absent locally, unclassified), not to World | fleet repo (test-only, reverted after) | Run the arm against that fleet repo; assert output contains `unclassified fleet-only path: tools/launchd/lib/new-helper.sh`. If Phase 3's enumerator misses it, the report is absent and the test fails. | The assertion reads the unclassified report; the report is produced only by Phase 3 enumerating the **fleet** tree. This proves the arm LOOKS at the fleet tree (an addition), not merely that it FIRES on a removal — the round-1 mutant's gap. |

---

## 7. Declared residuals

1. CI cannot certify driver currency — the fleet arm is loud-skipped there, so CI's green is honest but does not assert the driver is current.
2. A fleet advance between rig runs is caught only at the next rig run of `verify_go.sh`.
3. The arm depends on the fleet checkout being present and current on the rig; a stale *fleet* checkout would make the comparison compare against a stale fleet HEAD.
4. **Unclassified fleet-only paths are not certified.** The fleet carries 48 paths under the driver prefix (V33); World tracks 6 (V32) and requires 1 (`lib/pin-root.sh`, V38), so **41** fleet-only paths are unclassified today. They are reported loudly and counted (Phase 3), but do not red the gate — requiring them would red permanently on 42 correctly-absent files (other missions' plists/envs, fleet-only scripts, testdata, V37), and World cannot fix that (frozen core, C1). A fleet-only path becomes certified only when it is deliberately added to `REQUIRED_FLEET_PATHS`; until then it is a counted residual. This is the honest residual the quorum demanded: the success line says "tracked copy is current (untracked fleet additions not certified)".
5. The arm does not detect a *working-tree* change to the driver that is also present in the fleet (that is the existing working-tree arm's job, and it is preserved).
6. The `CI` env-var heuristic (V25) is GitHub-Actions-specific; a non-CI run that exports `CI=true` would loud-skip the arm even with a fleet checkout present (the comparison still runs when the fleet checkout exists, so this only affects the absent-fleet case).

---

## 8. Conflict surface

- **Existing working-tree-vs-HEAD arm (V24, lines 200–212):** the new arm is in addition; only the success line text (line 212) is relabelled to `(working-tree arm)`. The arm's logic and its `driver_tracked` control are untouched (C4). No behavioural conflict.
- **CI job (V17, line 166):** runs `./scripts/verify_go.sh` on `ubuntu-latest` in a checkout of `ailang-world` only. The fleet arm loud-skips (`CI=true`, no fleet checkout), so CI stays green. No conflict.
- **`--evidence-manifest-check` mode:** the new `--driver-fleet-check` mode is a separate branch; no overlap. Both sit before the `AILANG_BIN` gate (V41), so neither needs `AILANG_BIN`.
- **`AILANG_BIN` gate:** the main-flow call site is after the driver drift gate, which is after the `AILANG_BIN` check; the arm does not depend on `AILANG_BIN`. The isolated `--driver-fleet-check` mode is before the gate (V41) and needs no `AILANG_BIN`. No conflict.
- **`set -euo pipefail` (line 16):** the call site uses the `if` form (CB1, §4.2), which is exempt from `set -e`, so the CI loud-skip (return 2) is captured and never exits the script. No conflict.
- **Rig go gate:** on the rig, the fleet arm reds while the local copy is stale (the current state, V20) and while `lib/pin-root.sh` is missing (V10, Finding 1). This reds the rig's go gate until the fleet lands the current driver as a fleet-authored commit. That red is the point — it means "the fleet must commit," never "absorb it into your change" (per CLAUDE.md). It is a new, intended red, not a collision.
- **New env var `AILANG_FLEET_REPO`:** no existing consumer; default `$HOME/dev/sunholo-data/ailang` matches the known rig path (V3/V19).

---

## 9. Findings for the queue

1. **World's driver is missing the pin mechanism entirely** (V10–V12, V29–V31): the fleet driver sources `tools/launchd/lib/pin-root.sh` and re-execs into a pinned commit, while World's copy has no `lib/`, no `pin_root` reference, and derives `REPO` from `$0` with `MISSION_WORKDIR` empty. **Round 1 change:** this is no longer invisible to the arm — `lib/pin-root.sh` is the first member of `REQUIRED_FLEET_PATHS` (§4.1), so the arm now reds with a typed `MISSING LOCALLY` until the fleet lands it. It remains a distinct defect from the stale-copy row and still deserves its own row, but the arm now *detects* it rather than passing over it.
2. **The stale driver's `${VAR:-default}` governs effective config** (V16, V31): `MISSION_EXECUTOR_FALLBACK` resolves to the dropped `:floor` id because the env override is commented out. This is a *consequence* of the stale copy (the row's own live proof, M7), not a separate defect, but it means the fleet-comparison arm's green is necessary, not sufficient, for correct effective config.
3. **The row cites the gate at `verify_go.sh:205`; the gate spans lines 200–212** (V24). The `driver_drift` working-tree leg is at 205; the `driver_tracked` control is at 200–203 and the success line at 212. Cosmetic, no action needed.

---

## 10. Quorum round 1 — objections and disposition

Full-strength quorum: all three external reviewers present, `absent_reviewers` empty, cross-checked. All three rejected on the same surface — the arm's enumeration domain and the honesty of its success line. Direction (Option B) was not disputed. The controller additionally measured the union (V32–V38) and two first-party defects (CB1, CB2). Disposition per objection:

| # | Objection (verbatim core) | Disposition |
|---|----------------------------|-------------|
| 1 | gpt5-6-sol: "The proposed fleet arm can falsely certify 'driver is current' because it enumerates only paths already tracked by World and silently ignores fleet-only paths... `tools/launchd/lib/pin-root.sh` is present in the fleet but absent locally." Proposed fix: enumerate the union, require both blobs to exist and match, report MISSING LOCALLY / MISSING IN FLEET / DIFFERING, add an AC leaving `pin-root.sh` absent locally. | **Adopted the principle, not the literal union.** The controller measured that the literal union fix produces a gate that can never be green: at HEAD `7d3aa00` vs fleet `722e19c7`, World tracks 6, fleet carries 48, union 48, MISSING LOCALLY 42, MISSING IN FLEET 0 (V32–V36). The 42 MISSING LOCALLY are overwhelmingly files World must NOT carry — other missions' plists/envs, fleet-only scripts, 12 testdata fixtures (V37); only `lib/pin-root.sh` is genuinely required (V38). Requiring all 48 would red permanently on 42 correctly-absent files, and World cannot fix it (frozen core, C1) — the permanent-red mode Option C was rejected for. Instead: intersection comparison for DIFFERING (kept); MISSING IN FLEET as a typed failure (gemini's fix, below); an explicit committed `REQUIRED_FLEET_PATHS` set whose first member is `lib/pin-root.sh`, with a required-but-absent path a typed MISSING_LOCALLY failure (AC5); and unclassified fleet-only paths as a loud, counted, non-fatal report (Phase 3, §7 residual 4). The success line is scoped to the tracked copy (oc-glm's line, below). |
| 2 | gemini-3-1-pro: "The script silently ignores rogue local driver files that the fleet does not have. By using `[ -z \"$fleet_blob\" ] && continue`, any tracked file in World that is absent from the fleet will be silently skipped." Proposed fix: remove the `[ -z "$fleet_blob" ] && continue` line. | **Adopted.** The silent `continue` is removed; a World-tracked path absent from the fleet is now a typed `MISSING IN FLEET` failure (Phase 1, §4.1). Because MISSING IN FLEET is **0** today (V35), the branch is empty in the wild, so its test arm is **synthetic** (AC6 constructs a fleet repo missing one tracked path). |
| 3 | oc-glm-5-2: "The fleet arm's success line prints 'driver is current', but by the doc's own Residual 4 and Finding 1, the arm only compares files World *tracks*... the arm prints 'driver is current' when the driver is factually NOT current." Proposed fix: replace the success line with `echo "   ✓ fleet-comparison arm: $compared tracked frozen-core files match fleet HEAD $fleet_head — tracked copy is current (untracked fleet additions not certified)"` and remove the bare "driver is current" phrase from §4.3 and §4.4. | **Adopted verbatim.** The success line is exactly oc-glm's text (§4.3); the bare "driver is current" phrase is removed from §4.3's example output and §4.4's table. The success line now scopes "current" to the tracked copy and disclaims untracked fleet additions. |
| CB1 | Controller: §4.2's call site is broken under `set -euo pipefail` (line 16); a bare `check_driver_fleet` call that returns non-zero exits immediately, so `fleet_rc=$?` is never reached. Measured: with `set -euo pipefail` the script exits rc=2 and never prints REACHED; with `set -uo pipefail` it prints `REACHED: rc=2` and exits 0 (V39/V40). The CI branch (return 2) would turn CI permanently red. | **Fixed.** The call site uses the `if` form — `if check_driver_fleet; then :; else fleet_rc=$?; fi` — which is exempt from `set -e`, so the return code is captured and only a typed refusal (1) is a hard failure (§4.2). AC7 asserts the CI branch exits ZERO. |
| CB2 | Controller: §4.1 adds a `--driver-fleet-check` isolated mode, but every AC1–AC5 invokes the full `./scripts/verify_go.sh`, which runs `go build`, `go test -count=1` and a `-race` leg — an internal §4-vs-§5 inconsistency costing minutes per arm. Confirmed viable: `--evidence-manifest-check` is at line 111, before the `AILANG_BIN` gate at line 120 (V41). | **Fixed.** AC1–AC7 drive `--driver-fleet-check` (§5), which needs no `AILANG_BIN` (V41). AC8 (regression guard for the existing working-tree arm) necessarily uses the full script because that arm lives only in the main flow. |

| CB1-bis | **Quorum round 2 — ALL THREE reviewers reject on ONE surface, unanimously and independently:** the `--driver-fleet-check` isolated-mode block in §4.1 repeats the exact `set -e` defect CB1 was opened to fix. gemini-3-1-pro: *"the author left a bare `check_driver_fleet` call in the `--driver-fleet-check` isolated mode block... the script immediately exits with 2... turning a CI skip into a hard failure."* oc-glm-5-2: *"CB1's fix (the `if` form) was applied to the main-flow call site (§4.2) but NOT to the isolated mode. AC7... would fail against this code."* gpt5-6-sol (restored, see below): *"AC2 and AC7 fail, and CI callers of the isolated mode receive a hard failure—the exact CB1 failure the document claims to have fixed."* | **Fixed under Gate 2's NARROW-REFINEMENT CARVE-OUT** (bounded controller 2nd revision after the one re-quorum; the carve-out is already ratified for this mission, first used iter-13). Both preconditions hold and were checked: every remaining objection carries a concrete reviewer-authored `proposed_fix`, and none disputes the design DIRECTION — all three dispute one code shape. The applied text is the reviewers' own: gpt5-6-sol and gemini-3-1-pro proposed the **identical** `if` form, so the majority-verbatim form was taken (oc-glm-5-2's `\|\| rc=$?` variant is equivalent and was not needed). It also matches §4.2's already-accepted pattern, so no new shape was invented. **Non-vacuity shown, not asserted** (iter-147's durable lesson: a carve-out fix applied verbatim is still a fix whose PLACEMENT the controller authored): the guarding criterion AC7 **fails** on the pre-fix block and **passes** on the post-fix block, measured as two arms with the exit codes asserted to differ — V42 `rc=2`, V43 `rc=0` (`2 != 0`). |
| ABSENT | **`gpt5-6-sol` dropped out of round 2 on `budget`** — `estimated cost $0.1006 (doc ~13981 input tok) exceeds cap $0.1000 (pre-flight refusal, zero spend)`, i.e. refused over **$0.0006**. `.synthesis.absent_reviewers` = `[{"model":"gpt5-6-sol","reason":"budget"}]`, cross-checked against `[.reviewers[]\|select(.present==false)\|.model]` = `["gpt5-6-sol"]`, control `has("synthesis")` = true. | **Restored, not waived.** This is the self-selecting trigger the shared skill names: the reviewer drops out on budget precisely because the revision GREW the doc, and it was gpt5-6-sol whose round-1 `proposed_fix` the controller had **declined** (the literal union), so its opinion was the most load-bearing available. Re-run alone at `--max-cost-usd 0.25` for **$0.087215**: verdict **reject**, on the same CB1-bis defect, with the same `if`-form fix — and it did **not** re-raise the union objection, which is the evidence that the 42/0/48 disposition was accepted rather than merely unreviewed. Round 2 is therefore recorded as **full strength, 3 of 3 present**, never as "proceed at N−1". |

The 42/0/48 measurement (V32–V36) is the load-bearing constraint: it is why the arm's domain is the intersection **plus** an explicit REQUIRED set **plus** a non-fatal unclassified report, rather than the full union, and why the success line is scoped to the tracked copy.

---

## 11. Evaluator round 1 (sonnet, own worktree) — score 86/100 PASS, one BLOCKING finding, fixed in-sprint

**BLOCKING #1 — `compared` was unpinned, so the arm could under-certify while printing currency.**
The judge found, and the controller reproduced first-party, that a one-line regression in Phase 1
(a stray `continue`) silently drops a tracked driver path from the comparison and **nothing catches
it**: the path is not `differing` (never compared), not `missing_in_fleet` (that branch is never
reached for it), and not `unclassified` (Phase 3 skips it because World still tracks it). The arm
prints `✓ ... tracked copy is current` over a quietly smaller set. That is **this row's own defect
one level up** — a success claim wider than the axis actually measured — so shipping it would have
undercut the row. Reproduced on an 8-path green fixture: mutant LANDED (sha256), PARSES
(`bash -n` rc=0), `rc=0`, count **8 → 7**, currency claim still printed.

**Fix (controller, in-sprint): a two-counter Phase-1 accounting invariant.** `dispositioned` must
account for every path the loop SAW, and `expected_enumerated` — computed by a **separate**
`git ls-tree ... | wc -l` call, so a skip placed *before* the in-loop increment cannot hide from it
— must equal the number the loop saw. Mismatch is a typed FATAL, never a green.

**Non-vacuity shown, not asserted** (the standing lesson: a fix the controller places is a fix whose
placement the controller must show can fail). Two mutants, both LANDED by sha256 and both parsing:

| Mutant | Skip placed | Result | Which counter caught it |
|---|---|---|---|
| MUT-A (the judge's) | AFTER the counter | rc=1, `ACCOUNTING BROKEN`, no currency claim — `offered 8, saw 8, verdict on 7` | `dispositioned` |
| MUT-B (harder) | BEFORE the counter | rc=1, `ACCOUNTING BROKEN`, no currency claim — `offered 8, saw 7, verdict on 7` | `expected_enumerated` |

Control: the unmutated fixture is `rc=0` at 8 files. **Both counters are load-bearing** — each
catches a shape the other misses, so neither is decorative.

**NON-BLOCKING #2, accepted as a declared residual (see §7).** Phase 3's loud-residual net is scoped
to the same fixed pathspec as Phase 1, so a fleet-only driver file *outside* `tools/launchd/` and the
literal `scripts/mission_decisions.sh` is invisible — not even counted as unclassified. The judge
demonstrated it with an addition-shaped mutant (`scripts/mission_decisions_v2.sh`) and then measured
that it does **not** currently materialise against the real fleet. Latent, not live; filed to the
queue rather than fixed here (a pre-existing scope question is a queue row, not a sprint widening).

**Judge verdicts on the named targets:** T1 (the controller's own post-executor revert) PASS —
no shipped hunk perturbs any anchor `host/verifygate` matches on, whole package green; T2 (`set -e`
call sites) PASS; T3 (`set -u` + empty array, whitespace/glob paths) PASS, the suspicion **refuted**;
T4 no reproducible bug; T5 disclosed, not new; T6 real but latent (#2 above); T7 one vacuity gap
(#1 above).

