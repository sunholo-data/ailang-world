# Iteration 429 — controller readings (row 243 V4 S4)

All readings taken by the controller (`claude:claude-opus-5-5`) on the Studio, outside any sandbox, 2026-11-24.

## Rig experiments (design §11), before planning
- **E1 PASS.** `launchctl print gui/501/com.sunholo.ailang-world.worldd` rc 0, `state = running`, `program = ~/.ailang-world/bin/ailang-worldd`; `id -u` = 501; `ps -p <pid> -o uid=,comm=` → `501 ~/.ailang-world/bin/ailang-worldd`. Positive control: the `dev.ailang.mission-world` job pid resolved by the same awk.
- **E3 PASS.** The design's E3 block extracted verbatim and run under `/bin/sh`: rc 0 — synthetic `org.example.v4-fixture.501.<pid>` job in a `mktemp -d` home ran twice with two process identities and exactly the explicit HOME/PATH/UID (the script's own awk assertion under `set -eu`); the synthetic World files absent before and after; `launchctl list | grep -c org.example.v4-fixture` = 0 afterwards. The real World plist mtime unchanged (Nov 21).

## Plan
- Planner round trip (`ctl_roundtrip.sh`, loop M1 M2 M3, `docs` added to the comparison per plan F9): every `diff -r` of `cmd`, `host`, `docs` against `st-M*` rc 0 on a fresh `git archive f074dd49`. The two `.orig` files under `design_docs/verification/world-iter362/` are committed at base.
- Planner mutant ledger (`plan-scratch/mutants/ledger.txt`, 176 rows): K|NM|NM 60, S|NM|NM 73, K|K|K 32, S|S|S 8, S|S|E 3, control 1 — i.e. 73 named-test survivors whose whole-package verdict the sandbox could not take.

## Commits (numstat `git diff -M --numstat f074dd49 HEAD -- cmd host docs app`)
M1 `36bc27c1` 114 · M2 `af1a36bb` 1,232 · M3 `22473391` 1,400 · M4 `976f4028` 1,510 (test-only, judge r1 B1) · M5 `4fc443cd` 1,539 (Linux symlink order fix + test). Each committed by `commit_ms.sh` from the executor's snapshot after gofmt/build/vet/acceptance; snapshot M3 vs planner `st-M3` 18/18 `cmp`-equal; committed tree vs snapshot 0 differ.

## Gates (plan §6 `gates.sh`, outside the sandbox)
- On `22473391` (1,239 s): every leg rc 0 except `verify_go` (62 ok, 2 FAIL): `host/daemon` `TestSeToolsRunReplayIsByteEqual` and `host/verifygate` `TestMissionConfigDefaultFlowIntegration`; both packages byte-identical to base; both pass alone (`rerun_replay.txt` 4.88 s, `rerun_verifygate.txt` 8.37 s). Observation legs (real `~/Library/LaunchAgents`, `~/.ailang-world/bin`, `~/.ailang-world/*.json`, launchctl count) identical before/after.
- On `976f4028` with `GATES_SKIP_VERIFY_GO=1` (670 s): `FAILS: none`; `verify_go` leg reads `SKIPPED` (rc 0, 0 s) — not a run.

## Judge r1 B1, controller-reproduced
At `22473391`, with the supervisor check in `service_upgrade.go` and `|| verb == "cancel"` in `operator_service.go` both disabled: `go test -count=1 ./cmd/ailang-worldd/` rc 0 (232.7 s). At `976f4028`: the same edits, one at a time → `TestOperatorServiceUpgradePairFence/registered` FAIL and `TestOperatorServiceCancelConfirmation` FAIL (`declined cancel: 0 … action state = "cancelled", want "requested"`).

## PR #651 CI
- `976f4028`: build+test and race FAIL — `TestOperatorServiceOwnership/symlink_tool` on ubuntu-24.04 in both (`refusal "wrong type": 1 runner refused`; Linux lstat reports a symlink as 0777, so the writable-by-others check fired first). Race also: `host/broker` `TestWorktreeSnapshotRefusesAConfiguredFilter`, `host/daemon` 8 m timeout (dev's own run at `212a87ba` hit the identical timeout).
- `4fc443cd` attempt 1: symlink arm fixed; remaining reds all in untouched packages — `scripts/marketingcheck` `TestMarketingCaptureFakeChromeSixShots/single-pass` 120 s timeout; `TestFinishConfirmOnSurfaceDispatchesTheNamedDecision` (also red on branch `world/row228-search-h1` today; identical hashes on both sides of its message); `host/daemon` 8 m timeout. `gh run rerun --failed` once → attempt 2 **3/3 success**.
- Merge `ee8ef1c0` (`--match-head-commit 4fc443cd`; `origin/dev` had moved to `af7a610c`, 0 files shared): CI **3/3 success**.
