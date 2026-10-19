# Judge report r2, iter 284 row 156 M3 (HEAD 3b5a076, test-only delta over 7f6c762)
Score: 93/100 — PASS. No blocking findings; r1 NB1 resolved.
- Delta: handlers_exec_m2_test.go only (+9/-1); expected suffix now literals (darwin <wt>/**/*_test.go then <wt>/tests; linux <wt>/tests; none for protected=false).
- Clean run green: TestExecProtectedSettingsThroughHandler ok; TestExecSettings* ok; go vet ./host/broker ok.
- Mutant p.Protected=nil after load (exec_profile.go after ProtectedDigest line): FAIL /true; restored.
- Mutant drop `Protected:` in handlers_exec.go: FAIL /true; restored. git status clean afterwards.
- New findings: none blocking. Note: my r2 smoke of TestExecSrtProtectedWrites ran without WORLD_EXEC_SRT_NODE_MODULES (skipped); delta is test-only in a different file, r1 18/18 real-srt result stands. Linux still UNMEASURED-LOCAL.
