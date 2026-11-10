# Evaluator r1: row 201 A1b sprint 3 (iteration 377), judge worktree at 9f014e0d (base 545fdaae)

VERDICT: PASS, 92/100. Correctness 28/30, tests 26/30, security 19/20, standards 9/10, scope 10/10.

## Automated checks (mine)
- `go test ./cmd/ailang-worldd/ ./host/approval/... ./host/store/ ./host/authority/ -count=1` all ok (cmd 134.9s, store 13.7s, approval 0.7s, authority 1.5s).
- `go vet ./...` clean; `go test -run '^$' ./...` compiles everywhere (no output). AILANG_BIN=$HOME/.pinned-ailang/ailang. Did not run verify_go.sh or go test ./... (per instruction).
- Constraints: no ssh-keygen invocation (the only hit is the expected stdout string in a test); no `~/.ssh`, `/dev/tty`, exec.Command in the new tests; `readRegularFileCapped` is O_NONBLOCK then fstat, regular-file only, 16 KiB cap, re-capped on read; `git diff --stat` on `.github` and `tools` is empty; TestProductionContextRoots untouched; QUICKSTART section 15 is marked attended-pending.

## Mutation results (I re-derived; the planner's "102 killed" was not trusted)
Harness: /Users/voightkampff/.ailang/state/world-iter377/mut.sh and muts.sh. Each mutant asserts a non-empty `git diff` of the file, runs the named killer test only, then `git checkout`. Judge worktree ends clean (`git status --short | grep -v '^??'` empty).

PROCESS NOTE: my first full pass reported 100% KILLED and was INVALID: macOS has no `timeout`, so every run was "command not found". I caught it because a manual control (SIG-FORMAT) passed. I fixed the harness (KILLED now requires `--- FAIL`/panic/FAIL in the output; otherwise UNKNOWN) and re-ran everything. Three mutants did not compile (C11-KEY, A3-NOAPPR, CONSUME-SIGNED) and two of mine were mis-anchored onto sprint 2's revoke SQL (C8-REVOKED, C8-ENR); all five were rewritten and re-run. Only the re-runs count.

Totals: about 79 distinct mutants run (all 12 signed.go P/B/clone rows plus P6-DIGEST and APPROVEDBY-SEP; the store C1/C3-C8 and the approver SQL rows; the authority A1-A5 rows; the CLI F1-F14, SIG-*, MSG-*, V1/V2 rows; the four G1 store-flag rows), plus these production branches NOT in section 4: P6-DIGEST, APPROVEDBY-SEP, A4-CHAL, A4-EXP, SIG-FORMAT, SIG-NOW, MSG-ORDER (SignedMessage args swapped), TTY-AUTHORITY (typed path records `ssh:x`), C8-AT. Killed 77; equivalent 0 re-confirmed; survived 2 (distinct) plus 1 named-test-only survivor.
- SURVIVOR SIG-FORMAT (`approve_sig.go`: `Format: enrolled.Format` -> `approval.FormatSSH`): TestApproveBySignatureRefusalsBeforeAnyEffect still ok. Control: P3 (law level) and C8-FORMAT (store) both KILLED. Defence in depth means a non-ssh row is still refused at the store (C11 `other-format`), so this is not exploitable; it is only an untested CLI-level arm. Non-blocking.
- SURVIVOR SIG-NOW (`approve_sig.go`: `Admit.Now: now` -> `row.CreatedAt`): ok. Control: SIG-ENROLLEDAT KILLED. The mutant would wrongly refuse (ErrApproverClockBehind) a key enrolled AFTER the request was created but before approval, which is a realistic flow (request, then enroll, then approve). No test pins that the key may be enrolled between creation and approval. Gap in tests, not a defect in the code. Non-blocking.
- A4-EXP (`Binds(..., rm.ExpiresAt)` -> `rm.RequestCreatedAt`) survives the single test named for it in my own pairing (TestMintForRequestRefusesAnUnboundApproval) but is KILLED by the whole authority package (`Test.*`) and by the CLI test TestApproveBySignatureMintsWithNoTerminal. Control satisfied; not a survivor.
- Pre-existing, outside sprint 3: removing `AND revoked_at IS NULL` from sprint 2's REVOKE SQL (approvers.go:204) survived TestSignedConsumeChecksTheApproverInTheTransaction. That test does not cover that SQL, so this proves nothing about it by itself; I did not run it against the sprint-2 killers. Mentioned only for completeness.
- Every named security branch confirmed killed with a real `--- FAIL` (spot-checked output for C8: subtests active, revoked-in-the-window, unknown FAIL).

## Security properties confirmed (code read plus mutants)
1. A signature cannot approve another request: the message binds stored digest, row challenge and row expiry. P6 (challenge/expiry/digest each killed), Binds B2-B4 and authority A4 (A4-CHAL, A4-EXP, A4-DIGEST killed) re-check at mint time, and the store UPDATE binds `challenge = ?` (C6/C7 killed), so a swapped challenge cannot consume.
2. A revoked, replaced or unknown key cannot approve: P4 before cryptography, and again INSIDE the consume transaction (C8 UPDATE with `revoked_at IS NULL`, `format`, `public_key`, `enrolled_at <= at`; C8 revoked/format/key/enrolled/N/ERR all killed). The revoke-race test (a revoke on a second connection inside Clock) refuses and leaves request pending with no session.
3. No write on refusal paths: the approver mark and the request update share one transaction, so approver refusal rolls back; the CLI refuses all signature checks before the network (SIG-ORDER row not re-run by me, but SIG-REFUSAL and F10-F14 killed); `--write-message` is O_EXCL, never overwrites, writes no store row and returns before any terminal (F5, F5-EXCL, F6-RETURN killed).
4. `approved_by` cannot be forged: the authority derives it from `rm.Signed` and refuses a caller-set ApprovedBy or Approver (A1, A1-APPR killed); the store independently requires consistency (C3/C4 killed); `SignedApproval` fields are unexported with a fail-closed zero value (B1, B5, B6 killed).
5. Typed-`y` path unchanged: terminal still opened when `Signed == nil` (V1, V1-CONSUME killed), records `tty` (A5, TTY-AUTHORITY killed), touches no approver (C8-ALWAYS killed, CS4).
6. File input: FIFO/device refused via fstat after O_NONBLOCK; 16 KiB bound (F9-BOUND killed).
7. G1 closed: approver verbs and `approve --list` honour `--credentials-db` and the default `--db` (all four G1 mutants killed).

## Findings
Blocking: none.

Non-blocking:
1. Test gap, SIG-NOW: add a CLI test where the approver is enrolled after the request's created_at but before now, and approves successfully.
2. Test gap, SIG-FORMAT: a CLI arm with a non-ssh-format approver row would pin the CLI-level refusal (the store backstop currently carries it).
3. The plan's section 4 killer for MUT-A4-EXP-type edits is not isolated to a single named test (only whole-package kills); tighten if sprint 4 touches A4.
4. Judge process lesson: a mutation harness must be checked for a surviving control before trusting 100% KILLED (macOS lacks `timeout`).
5. Plan's "102 mutants" is a planner claim; my independent count is about 79 distinct, with 77 killed and 2 survivors above (neither exploitable).

Files: /Users/voightkampff/dev/sunholo-data/.wt-world-iter377-judge (restored clean); harness in /Users/voightkampff/.ailang/state/world-iter377/.
