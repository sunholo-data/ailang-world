# Evaluator report r1 — iteration 385, row 201 A3a sprint 2 "the routes" (base a8c66655, head d52b3ae3)

VERDICT: PASS — 92/100. No BLOCKING findings.

## Scores
- Correctness vs design / plan §9 (30): 28
- Tests: non-vacuous, kill mutants, cover operands (30): 27
- Security / fail-closed (20): 19
- Coding standards (10): 9
- Scope and size (10): 9
Total 92.

## Review summary (production diff read in full)
- Every new route is wrapped by `approvalRoute`: Host is checked against O (loopback at the arriving LocalAddr port + `--surface-origin` + active passkeys' origins) and refused 421 before any handler. The only store access ahead of the Host check is the read-only enrolled-origin lookup that O needs.
- POST order is R1, R3 (Origin exact in O; missing or "null" refused), R4 (Sec-Fetch-Site absent or same-origin), R5 (media type), R6 (MaxBytesReader + strict grammar), R7, R8, R9, R10, R11 (sshsig cut, then `authority.JudgeAssertion` at the daemon `now`), R12 (`Binds`), R13 (non-waiting `OpenCredentials`, closed immediately), R14. Every refusal precedes `RecordSurfaceApproval`, and the single record site is after the judge. It never mints and reads no Cookie or Authorization header. The error path of the judge refuses with 422, and a non-store judge error never records.
- Only the constant-text OPTIONS check exists in tests. No CORS header is set anywhere in production, and a preflight is 405 (the mux has no OPTIONS route).
- The read JSON carries no challenge, artifact, token or message hash. The list carries only hash16 and the card carries the raw nonce and stored bytes, as designed.
- Intermediate commits: M4d/M4e register read routes only, and M5a adds the POST with the guards in the same commit, so no route writes before its guards exist.
- Scope: `git diff --name-status` matches plan §5 plus the two plan files. Nothing is touched under tools/launchd, world/, packages/, scripts/, host/approval, host/store or host/authority. The only `isProtected`/`passkey.js`/`app.js` production-diff lines are comments and test assertions. There is no `surfaceCSP` diff. There is no new exec site, `/v1` route, CORS header or `t.Parallel` (the grep hits are test assertions on the absence of Access-Control).
- Numstat `git diff --numstat a8c66655 HEAD -- cmd host`: 1596 changed lines (ins+del; prod ~ 20/2 main.go, 18 daemon.go, 291+91+267 new prod files), under the 2000 cap.
- The fallback cut is honestly bounded: an sshsig POST returns 422 naming the terminal path (R11 cut, killed by J17), and RT6, RT13 and AW6's route half are in Appendix C.

## BLOCKING findings
None.

## Non-blocking findings
1. J36 SURVIVED: `handleSurfaceRequestCard` (surface_approve.go, `if err != nil && !isOlderStore(err)` after ActivePasskeys). Replacing it with `if false` leaves every test green, so a non-older ActivePasskeys error is untested. The card would then silently show zero passkeys instead of a 500. Read-only and fail-soft; low severity. A test needing an injectable store error would close it.
2. R12 `Binds` is a documented equivalent (plan). The plan's MUT-R12 was not re-run by me; it is covered by design analysis only.
3. The `*.ts.net` fix-text arms depend on `ParseCeremonyOrigin`. J11b (the plan's MUT-CROSSORIGIN-TSNET) killed.
4. approvalOrigins opens a read-only store handle on every route call, including foreign-Host requests. This is necessary for enrolled origins (design) and costs only a read; noted as a cheap unauthenticated-request amplification surface bounded by loopback bind.

## Mutants (ledger: /Users/voightkampff/.ailang/state/world-iter385/judge/mutants.tsv, 64 lines; driver scripts in .../judge/bin/)
Commands: `go test -vet=off -count=1 -failfast <pkg> -run <named regex>`; each restored via `git checkout --`, worktree clean at the end.
- 56 KILLED, 2 SURVIVED, 5 COMPILEERR (J11, J31, N01, N08 first attempts, J31b: unused import/var, redone as J11b/N01b/N08b; J31b not redone), 1 no-diff.
- Plan §4 rows reproduced (roughly 25, M4a to M5c; each killed by its named test or a sibling in the same package): HOST-PREFIX-OF (J01), O-127 (J05), O-ENROLLED (J06), O-CONFIGURED (J07), R3 (J12), R4 (J13), R5 (J14), R9 (J15), R8 (J16), SSHSIG-CUT (J17), WRITER-KEPT (J18), R14-RACE (J19), RETRY-AFTER (J20), STOREERR-BUSY (J21), NOSTORE (J22), NOSNIFF (J23), CONFIGURED-ORIGINS (J25), LOCALADDR (J26), R6-SIZE (J27), HOST-CHECK (J28), HOST-CASE (J02), SFS-ABSENT/SAMESITE (J03/J04), NEW-ORIGINS (N25), CLI-GRAMMAR (N23), THIRD-WRITER-style/RECORD-BEFORE-LAW (K01), carried AW-JUDGE-AT-CREATION ssh/pk (N24a/N24b: both KILLED by TestApproveAwaitJudgesTheApproverAtNow/ssh and /webauthn).
- My own mutants not in the table (>=15), all KILLED: J08 (empty Origin accepted), J29 (wrong approver recorded), J30 (nil artifact recorded), J32 (record under another id), J33 (unavailable not set), J34 (decided cap 1000), J35 (card 404 -> 200), J36 (SURVIVED, see finding 1), N01b (Contains for Host), N02 (drop localhost origin), N05 (credential16 -> 8), N07 (413 -> 400), N08b (trailing-token check off), N09 (consumed text uses Created), N12 (421 -> 400), N16 (R8 409 -> 400), N17 (judge refusal status 200: fail-open), N18 (CrossOrigin status 200), N19 (Retry-After 0), N20 (Recorded dropped), N21 (card passkeys dropped), N22 (success status 202), J11b, K02 (any judge error -> 503), LEAK1 (challenge leaked into list), LEAK2 (request bytes leaked into decided/card).
- Must-die control K01 (record before judge): KILLED by TestDaemonCredentialWriterSites (and plan's named recording-before-law test class). Must-live controls: L02 (comment-only edit) SURVIVED as expected; N24c (identity edit) produced no diff and is recorded as COMPILEERR(no diff), not credited. Note L01 (plan treats `now` vs `d.approveNow()` at record time as killed) was KILLED by R14-race, confirming the plan row rather than being a live control.

## Per-milestone non-vacuity
- M4a/M4b model: N05, N08b, N09, N20, LEAK2 killed in surfacepage tests.
- M4c guards: J01..J10, N01b, N02 killed by TestApprovalGuards; N25 (New's origin check), N23 (serve flag) killed.
- M4d/M4e read routes: J28 (Host check removed -> TestApprovalRoutesRefuseForeignHost), J33..J35, LEAK1 killed by TestRequestsJSON/TestRequestCardJSON.
- M5a/M5b/M5c POST: J12-J15 (a guard call removed), K01 (record before judging) killed, J18 (kept writer) killed, J20/J21 busy/Retry-After killed; carried AW-NOW N24a/N24b killed. The E2E seam TestSurfaceApproveThenAwaitMints passes in the baseline run.

## Test runs (AILANG_BIN pinned; Chrome/passkey env unset)
- `go test -count=1 -timeout 20m ./host/surfacepage/ ./host/daemon/ ./cmd/ailang-worldd/`: rc=0 (surfacepage 0.6s, daemon 272s, cmd 107s), output at judge/base.txt.
- `go test -race -count=1 ./host/surfacepage/`: rc=0. `go test -race -count=1 ./host/daemon/ -run '^(TestApprov|TestRequest|TestDaemonCredentialWriterSites|TestSurfaceAppRoutes)'`: rc=0 (judge/race_vet.txt).
- `go vet ./host/... ./cmd/...`: rc=0. `go test -run '^$' ./...`: compiles, no failures.
- No flakes seen. `git status --short` in the judge worktree is clean.
