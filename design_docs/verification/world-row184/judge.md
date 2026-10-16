# Row 184 judge r1 (independent) -- score 88/100, PASS, no blocking

Gates (all rc=0): go vet; go test ./host/broker/... ./cmd/ailang-worldd/... (real tool + srt env); -race broker subset;
verify_ail.sh. `go test ./host/...` had ONE failure, TestExecHandlerTimeoutKillsTheGroup (calibration: first control took
3.67s) while I ran a concurrent go test; passes 3/3 in isolation -> load flake, not row-184 (but note it runs the delta now).

AC1 worktree_delta.go:131-160 delta, handlers_ailang.go:455-461,619 / test worktree_delta_test.go:278 (literal bytes pin) + tool_test:75
AC2 handlers_exec.go:258,291 / test :441,456 (git stub + real srt)
AC3 worktree_delta.go:178-181 / test :149 (300 -> 256, total 300, git+walk)
AC4 gitSnapshot --ignored=no, basis; walkSnapshot cap / tests :78,:171
AC5 why.go:573-628,806-821,876 / why_test.go:366 (synthetic result objects, not end to end)
AC6 tool_test:31 real v0.52.1: exactly 36 .ailang/cache files, re-run => no delta. PASS.
Overhead on tps-a (55,755 files, linked worktree): ~25 ms per snapshot, ~50 ms per call.
Hostile .git/config: fsmonitor, hooksPath, pager, sshCommand, alias, diff.external NOT executed. filter.<x>.clean (+ .gitattributes
`* filter=x`) IS executed by `git status` (marker file created). Requires write to .git/config, which the sandbox denies
(exec_settings.go:70; .gitattributes also deny-write) -> hardening, non-blocking.
