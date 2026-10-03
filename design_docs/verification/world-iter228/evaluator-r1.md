# World iteration 228 — independent evaluator verdict (round 1)

- **Subject:** `sprint/w-verifygate-etxtbsy` at `483ac8e` (PR #194; squash-merged as `382fd7e`).
- **Evaluator:** Agent tool, `opus`. It ran in a fresh context with a read-only checkout and a scratch worktree for its drills, which it removed afterwards. The executor ran on `sonnet` (Agent tool), so the generator and the judge are different models. The resolver's evaluator lane was `pi:openrouter/minimax/minimax-m3` (reroute generator-equals-judge). That lane was unusable this fire because the `openrouter` bucket is in `MISSION_OVER_RATION`. The next lane in the declared chain is `claude:claude-sonnet-4-6`, which would have been the same Sonnet family as the executor, so the controller used the chain's last entry, `opus`.
- **Verdict:** **PASS 88/100, zero blocking.**

## Evaluator mutation table (scratch worktree)

| # | Mutation | Result |
|---|---|---|
| AC1 | delete RLock/RUnlock in `forkLockedDo` | KILLED |
| AC1b | skip RUnlock | KILLED (5 s bound) |
| AC3 | drop CreateTemp from the banned set | KILLED |
| AC8 | disable the aliased-import ban | KILLED |
| E7 | `writeExecutableAt` back to `os.WriteFile` | KILLED (names `toolchain_pin_gate_test.go:484`) |
| E4 | run the kernel test on darwin (no ETXTBSY) | KILLED, so the assertion is live |
| E6a | raw write in a package-level initializer | KILLED |
| E1 | `createTempForkLocked` does its writes without the lock | SURVIVED |
| E2 | `forkLockedWrite` skips the lock when the mode is executable | SURVIVED |
| E3 | RUnlock before Close | SURVIVED |
| E5 | function value `wf := os.WriteFile` at a site | KILLED, but only by the floor |
| E5b | function value in a NEW helper | SURVIVED |
| E8b | `syscall.Open` + `os.NewFile` writes in a new helper | SURVIVED |
| E6b | scan drops its package-level visit branch | SURVIVED (no fixture covers it) |

## Non-blocking findings (filed as queue row 142)

1. The allowlisted wrapper bodies are unguarded (E1, E2, E3). The controller reproduced this by inspection: `rawWriteAllowlist` exempts `forkLockedWrite` and `createTempForkLocked` entirely, and test 1 drives only `forkLockedWrite` at mode 0o644.
2. The scan sees only direct calls (E5b, E8b).
3. No package-level initializer fixture covers the scan branch for package-level declarations (E6b).
4. The CI verbose step does not fail on SKIP. AC6 is read by a person, not machine-checked.
5. A future kernel that drops the exec write-deny rule would red CI through the kernel control. That failure is correct and loud, but it reads like a regression.
6. Wrappers and the allowlist are matched by bare name.

## Gates (evaluator, at 483ac8e)

- `go vet ./host/verifygate/`: clean.
- `go test ./host/verifygate/ -count=1` with `AILANG_BIN` set to the v0.41.0 pin: ok in 146.4 s.
- `verify_ail.sh`: PASSED.

CI on the PR head (run 37137525867) was 2/2 success. Verifygate took 238.1 s under race, against the pre-fix baseline of 215 s (+10.7%). The verbose step printed `--- PASS: TestKernelRefusesExecOfWriterOpenFile` on Linux.

## Controller verification (Gate 3b, merge commit 382fd7e)

CI run 37138466945 was 2/2 success. Verifygate took 130.5 s plain and 155.1 s under race. The verbose step shows `--- PASS` for TestForkLockedWriteBlocksConcurrentFork, TestKernelRefusesExecOfWriterOpenFile and TestVerifygateTestWritesAreForkLocked.
