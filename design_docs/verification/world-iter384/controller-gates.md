# Iteration 384 — controller gates (row 201 A3a sprint 1, head `9bfd0e2b`)

Base `0a25c0d6`. Sprint branch `sprint/world-iter384-row201-a3a-s1`, PR #554, merge `b44b5f7d`.

## V32 base-drift trigger (design §14), at `0a25c0d6`
- V1: all 12 sha256 prefixes equal the design's (`daemon.go` `9a9357770a4b0360` … `passkey_cdp.mjs` `2a336e5e810f8ed6`); `git diff --quiet 0140bd7f HEAD -- cmd host` → identical; `go1.26.6 darwin/arm64`.
- V8: `OpenCredentialsReadOnly` sites in `host/daemon` (non-test) = 5; writer opens (`OpenCredentials(`) = 0.
- V18: `go test ./host/daemon/ -run '^(TestFrozenV1RouteTableMatchesMux|TestSurfaceAppRoutes|TestPasskeyPageRoutes|TestDaemonDependencyAllowlist)$'` rc 0.
- V15: `TestPasskeyPageInChrome` (required, headless Chrome, temp profile, virtual authenticator) rc 0, 3 subtests PASS.
- Verdict: no drift → no r3.

## Round trip (controller, independent)
`git archive 0a25c0d6` → the 7 Appendix-P patches extracted from the plan file with §2's awk → `patch -p1`:
`diff -rq` against the planner's `plan-scratch/work/` tree empty for `cmd` and `host`. After the last
sub-milestone commit the sprint worktree's `cmd`/`host` are `diff -rq`-identical to the round-trip tree.

## Sub-milestone commits (commit_ms.sh; numstat `git diff --numstat 0a25c0d6 HEAD -- cmd host`)
| M | commit | numstat | plan |
|---|---|---|---|
| M1a | 38a6d903 | 70 | 70 |
| M1b | 2912eeea | 404 | 404 |
| M1c | 4fb902b4 | 463 | 463 |
| M2a | e34a9e61 | 839 | 839 |
| M2b | 6f75243d | 971 | 971 |
| plan | b506ce93 | 971 | — |
| M3a | c6d4fb6c | 1313 | 1313 |
| M3b | 9bfd0e2b | 1543 | 1543 |

**Controller error, recovered:** `commit_ms.sh` restores every untracked path that is not in the
snapshot, so its first run DELETED the two uncommitted plan files (earlier iterations committed the
plan first). Regenerated from the planner's own assembly parts (`plan-scratch/plan/*` +
`tools/assemble.py`, `tools/mkjson.py`): 2,687 lines as reported; all 7 extracted patches `cmp`-identical
to the copies extracted before the deletion. Committed as `b506ce93` before M3.

## Plan §6 pre-commit script (gates.sh) at 9bfd0e2b
rc 0, `PRE_COMMIT_GATES_GREEN 408s`: verify_ail, gofmt, build, vet, compile fence, touched packages,
`-race` legs, 33 named inventory/schema/help/quickstart/new tests `--- PASS` (0 SKIP/FAIL), AW leg 8 `--- PASS`.

## VERIFY_GO_LEGS=plain verify_go.sh at 9bfd0e2b
rc 1 (load 15–19, run alone, before the judge): `✓ 0 other binary blobs among 6753 tracked files`, toolchain
floor ✓, `✓ all 64 required top-level evidence tests passed exactly once`; two reds in packages the sprint
does not touch (`git diff --stat 0a25c0d6 HEAD -- host/broker host/verifygate` empty):
`host/broker` `TestExecProtectedCheckTimeout/reserve` and `host/verifygate`
`TestMissionConfigDefaultFlowIntegration/invalid_ledger` (both on the plan's §6 load-sensitive list).
Re-run alone: both rc 0 (`TestExecProtectedCheckTimeout` 8/8 subtests, `invalid_ledger` PASS).
CI on the PR head 3/3 is the arbiter (below).

## Judge survivor reproduced (controller)
`sed 's/Artifact, now)/Artifact, p.row.CreatedAt)/' cmd/ailang-worldd/approve_await.go` (2 lines changed)
→ `go test ./cmd/ailang-worldd/ -run 'Await|JudgesTheApproverAtNow'` rc 0: SURVIVED. Test-only gap carried
to sprint 2.

## CI (Gate 3b)
- PR #554 head `9bfd0e2be5bb36e72a27dc3d92340c63691727e8`: check-runs total=3 completed=3 success=3 (run 43111055834), first attempt.
- Merged `gh pr merge 554 --merge --match-head-commit 9bfd0e2b…` → `b44b5f7dd8860fb044ee9ba323d3a07164862226` at 22:37:23Z.
- Push CI on the merge `b44b5f7d` (run 43113317174): 3/3 success at 22:49Z. LANDED.
