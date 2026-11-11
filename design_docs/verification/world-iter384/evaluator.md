# Evaluator report r1 - iteration 384, row 201 A3a sprint 1 (base 0a25c0d6, head 9bfd0e2b)

## Verdict: PASS - 91/100. No BLOCKING findings.

| Category | Score |
|---|---|
| Correctness vs design and plan §9 | 28/30 |
| Tests (non-vacuous, kill mutants, cover operands) | 24/30 |
| Security / fail-closed | 19/20 |
| Coding standards | 10/10 |
| Scope and size | 10/10 |

## Review
- Schema 7 DDL in credentials_schema.sql and in surfaceApprovalsSQL are the same, column for column. The upgrade is `if version < 7` inside `upgradeCredentials`; a session open does not upgrade. Mutant U1 and the schema-7 test cover it.
- RecordSurfaceApproval does S1, S2, S3 and S4 in one tx. Every refusal returns before the INSERT, so 0 rows are written. S4 uses `now >= expAt` and `now < createdAt`, with the right sense. S5 is upsert, only reachable for an unconsumed, unexpired request. `received_at` is set to `now`. The prune DELETE is in the same tx as PutMintRequest and is guarded by `version >= 7`. A schema-6 store still works.
- judge.go is a pure lookup plus the A1b/A2 law. Digest and challenge are re-decoded and refused unless 64 hex. The enrolled credential comes from the read-only handle. The row's `approver` column is never read by the minter.
- The minter (`pollSurfaceApproval`, `awaitSurfaceApproval`) re-reads the request and re-judges the artifact on every poll. Admission (used, expired) is checked before returning `signed`. The deadline is `min(expiry, now+awaitFor)`. A refusal is printed once per `received_at`. Every error path returns a refusal and nothing mints. The flag checks happen before any store read, and `--await` is the sixth exclusive mode. No terminal is opened.
- §1.4 deviations (P1 to P6): P4 (the judge's `request` and `store` steps) is sound: it is fail-closed and keeps the refusal text byte-identical. P3 is one new `context.Background()` root, pinned in `context_roots_test.go`. P1 and P2 are mechanical test edits. I found nothing unsound.
- Scope: files match §5 plus the two plan files. Nothing under tools/launchd, world/, packages/, scripts/, world-publish, host/daemon, host/surfacepage or host/approval changed. There is no new exec site, route or module. Exactly one new production context root.

## Non-blocking findings
1. **Mutants NOW and NOW-ssh SURVIVED** (`cmd/ailang-worldd/approve_await.go` pollSurfaceApproval, lines 82 to 87). The minter can judge at `p.row.CreatedAt` instead of its poll `now` and every await test stays green. The plan's MUT-WP-JUDGE-NOW only covers the approve wrappers. Production code is correct (it passes `now`). The gap is test-only: no await test enrols or revokes an approver between request creation and the poll. Add an await analogue of `TestApproveByPasskeyJudgesTheApproverAtNow`. Reproduce: `bash mut.sh NOW cmd/ailang-worldd/approve_await.go 's/Artifact, now\)/Artifact, p.row.CreatedAt)/' ./cmd/ailang-worldd/ 'Await|JudgesTheApproverAtNow'` gives rc=0 (SURVIVED). The ledger has both lines.
2. Mutant ROW (removing `|| !p.recorded` after `SurfaceApproval`) survived. It is equivalent: an unrecorded row judges an empty artifact, the refusal is not shown because `recorded` is false, and nothing mints.
3. W-die (poll interval 0) and W-die2 (1ns) are KILLED with rc=1 but no `--- FAIL` line was captured (likely a timing or deadline assertion). They are not claimed as strong kills.
4. The F1 mutant made a test write a stray file `cmd/ailang-worldd/m` (a relative `--write-message m`). I removed it. This shows a flag test that would touch the cwd if exclusivity broke. It is harmless in a green run.

## Mutants run (ledger: $HOME/.ailang/state/world-iter384/judge/mutants.tsv, 58 lines, one per run incl. re-runs)
Tally of ledger lines: 44 KILLED, 5 SURVIVED, 7 COMPILEERR, 2 NOEDIT. The COMPILEERR and NOEDIT lines were re-run with a corrected edit, which is a further line.

Controls: C-die KILLED (S4 sentinel swap, TestRecordSurfaceApprovalGuards/S4-expired). C-live SURVIVED (comment edit). W-live SURVIVED (comment edit).

Store: S4a (`>=` to `>`), S4b (`<` to `<=`), S4c (consumed case off), S4d (digest case off), S2a (approver-name length), S2b and S2b2 (artifact `>` to `>=`), S2c (`<1` to `<0`), S2d (hex check off), S2e (format `&&` operand), S5 (ON CONFLICT DO NOTHING), S5r (received_at = row's), S1 (version guard off), S3 (unknown-request off), P1 (prune version guard off), P2 (prune deletes all), U1 (upgrade `<7` to `<6`): all KILLED. Tests: TestRecordSurfaceApprovalGuards, TestCredentialsSchema7SurfaceApprovals, TestSurfaceApprovalPrunedWithItsRequest.

Judge: J1 (now to CreatedAt), J2 (Consumed false), J3 (challenge length operand), J4 (RevokedAt 0), J5 (zero challenge), J6 (BytesDigest of nil), J7 (ssh unknown-key off), J8 (assertion unknown-key off), J-nolookup (passkey lookup with nil credential), J-nolookup-ssh (ssh lookup with nil key): all KILLED by TestJudgeAssertionIsTheCLIJudge and TestJudgeSSHSignatureIsTheCLIJudge.

Minter and flags: M6 (`>=` to `>`), M6b (awaitFor ignored), M6c (`awaitFor>0` operand), M4 (print every poll), M4b (shown reset), M5 (signed ignored), M2 (admission check off), NF (deleted request not named), F-ssh (format switch inverted), AW7 (no re-judge: JudgeSSHSignature replaced by nil, nil), F1 (await removed from the exclusive-mode list), F2 (`!*await` operand), F3 (`<1s` to `<0`), W-die and W-die2: all KILLED. NOW, NOW-ssh and ROW SURVIVED (see findings 1 and 2).

Plan §4 sample: the §4 rows overlap my diff-anchored set on S1, S2, S3, S4, S5, the prune, the judge lookups, M2 to M6, the flag rules and AW7. All were killed by the tests the plan names, except the await-time judging gap in finding 1. My own mutants number 47 distinct edits against the sample's 15 minimum.

## Non-vacuity per milestone
- M1b: S1, S2 (a-e), S3, S4 (a-d) and S5 guards each reverted alone, and the guard tests (SA2) fail each time.
- M1c: P1 and P2 revert the prune, and TestSurfaceApprovalPrunedWithItsRequest (SA3) fails.
- M2: J-nolookup and J-nolookup-ssh skip the credential lookup, and JG1 and JG2 (the two TestJudge*IsTheCLIJudge tests) fail.
- M3: AW7 mints from the artifact without re-judging, and TestApproveAwaitNeverTrustsTheRow and TestApproveAwaitSSHSIG fail.

## Size
`git diff --numstat 0a25c0d6 HEAD -- cmd host`: 1543 (insertions + deletions), under the 2,000 limit.

## Test runs
- `go test -count=1 -timeout 20m ./host/store/ ./host/authority/ ./cmd/ailang-worldd/` rc=0 (store 13.8s, authority 0.9s, cmd 281s). Output: judge/t1.out.
- `go test -race -count=1 ./host/authority/ ./host/store/ -run 'Credentials|SurfaceApproval|MintRequest|Judge'` ok, both packages. Output: judge/race.out.
- `go vet ./host/... ./cmd/...` rc=0.
- `go test -run '^$' ./...` rc=0 (compile fence).
- I did not run verify_go.sh or the whole `go test ./...`, as directed.
- `git status --short` in the judge worktree is clean at the end.
