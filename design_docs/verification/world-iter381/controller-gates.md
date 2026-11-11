# Controller gates — iteration 381 (sprint head `64a239f8`)

Plan §6 run as one bash script (`gates.sh`, extracted verbatim from the plan plus the after-commit `verify_go.sh` leg) AFTER the last commit and BEFORE the judge.

- **Pre-commit legs: green** (`PRE_COMMIT_GATES_GREEN`). `verify_ail.sh` PASSED (40 named tests); gofmt clean; `go build ./...`; `go vet ./...`; compile fence `go test -run '^$' ./...`; touched packages (`host/approval/...`, `host/store`, `host/authority`, `cmd/ailang-worldd`, `host/daemon`) ok; `-race` on `host/approval/...`, `host/store` and the plan's `cmd/ailang-worldd` set: 0 FAIL; the 30 named inventory, reflective, help and quickstart tests: every one `--- PASS`, 0 SKIP/FAIL.
- **`VERIFY_GO_LEGS=plain ./scripts/verify_go.sh` after commit: rc 1**, under machine load ~32 from other sessions, on four timing tests:
  - `cmd/ailang-worldd` `TestHostsLaunchersExecHost` (`hosts_test.go:223: codex: signal: killed`, a launcher-subprocess timeout);
  - `host/broker` `TestExecProbeNet6BoundBeforeAnyPick`, `TestExecToolFailedProbeRefusesF15` (`handler subprocess timed out after 5s`), `TestExecToolArgv0FollowsARepointedLink` (`TimedOut:true`).
  - **Each re-run alone at the head: PASS** (6.56 s; 8.32 s, 0.75 s, 1.97 s). `go list -deps ./host/broker/` names no touched package (0 matches for `host/approval` or `cmd`). The planner measured the same `host/broker` and `host/verifygate` timing reds at the BASE under the same load (plan §1.1/§6).
- Every other package in the sweep `ok`.

Per-sub-milestone rebuild with `commit_ms.sh` (gofmt, build, vet, the sub-milestone's named tests): every commit rc 0. Numstat (`git diff --numstat ab4b2402 HEAD -- cmd host`) after each commit: 38 / 237 / 701 / 945 / 1,190 / 1,397 / 1,538 / 1,686 = the plan's §7 values. Final `cmd`/`host` tree `diff -rq`-identical to the planner's mutant tree.

Plan round trip re-done by the controller: the 8 patches extracted from the plan file with §2's awk and applied to a fresh `git archive` of `ab4b2402` reproduce the planner's `work/` tree (`diff -rq cmd host` empty, rc 0).

Judge findings reproduced by the controller in the judge worktree (clean afterwards):
- kill: dropping WN5's all-zero check (`|| challenge == [32]byte{}` → `|| false`, vet OK) → `TestPasskeyEnrollRaceAndRefusals/WN5` FAIL (rc 1);
- survivor: swapping `--db`/`--credentials-db` in `--webauthn-link`'s printed next step → `TestApproveByPasskey` ok (rc 0). Fails closed: the operator's pasted command would be refused.

CI: PR #550 head `64a239f8` attempt 1 red in the race job only, `scripts/marketingcheck` `TestMarketingSupervisorFakeDaemonCleanup/shell-int` (`supervisor_test.go:706` trap `signal: killed`, `:710` cleanup 29.93 s over its bound), run 42993721120. Known row-254 class (folded into row 258); this sprint touches no `scripts/` file. `gh run rerun --failed` → 3/3 success.
