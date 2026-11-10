# Evaluator r1 - iteration 374, row 201 sprint 2 (86aa2f9c)

Verdict: PASS, 90/100. No blocking findings.

Score: Tests 20/20, Lint 10/10 (gofmt -l empty, go vet clean), Acceptance 27/30, Quality 14/15, Docs 11/15
(QUICKSTART s14 is 12 lines vs ~45 planned; design doc still in planned/ by design), Fidelity 9/10. Conditional categories not triggered.

Gates run in the judge worktree: go build ./..., go vet (4 pkgs), go test of the 4 pkgs -count=1 (all ok; cmd/ailang-worldd 380 s),
go test -race TestConcurrentApprovalsMintOnce -count=3 ok. No t.Parallel, no /dev/tty, ssh dir or ~/.ailang use in new tests; no binary blobs; gofmt clean.

Mutation: 58 applied, 56 killed, 2 non-counted (CR-TOOMANY survived = real gap; ST-INSERTOUTSIDE was an equivalent mutant of my own making).
First-run compile-error kills (AP-DIGEST, AP-CLOCKREUSE, BM-EXPIRES) were re-done as clean mutants and are killed by behaviour.
KILLED (test): ST-ROWS, ST-CLOCK2, ST-CLOCKPRE, ST-NOCHECK, ST-EXPGE, ST-DIGESTW, ST-CONSUMEDW, ST-CREATEDLE, ST-NOTX, ST-CLASSIFY,
ST-OLDER-CONSUME, ST-CAPGE, ST-PRUNECONS, ST-PENDFILTER, ST-PENDEXP, ST-ID, ST-UPG-NOMR, ST-CHECKNIL, ST-CLOCKNIL, ST-OLDER-PUT (store tests);
AU-PRECONSUME, AU-NOCHECK, AU-REQCREATED, AU-TOKENHASH (TestMintForRequestChecksTheRowInsideTheTransaction subtests);
AP-ADMIT, AP-REFUSALIF (Replay/Expired), AP-NODERIVE (StaleWorld), AP-SAMESTORE, AP-GRAMMAR, AP-AFTERNIL, AP-ADDR (BadIDsAndMisuse),
AP-DIGEST (TamperedRow), AP-MCPEXISTS, AP-CLOCKREUSE (ExpiryDuringWorktreePreparation + Expired/in-the-transaction), AP-APPROVEDBY,
AP-ANSWERBY, AP-LISTNOW, PE-NOCONSUME, PE-PLAINREAD, PE-BOUNDFROMLIFE, PE-PROMPTCHECK, PE-PROMPTBOUNDARY (<= to <), PE-ALWAYSBOUND (ExistingCallersReadUnbounded),
FM-CONSUME, FM-GIT, FR-CREDSDB, FR-PROFILESDIR, FR-OUTDIRCREATE, BM-REPOABS, BM-EXPIRES, CR-EXPIRY, CR-RANDOMNIL, CR-CHALLENGE-SWAP,
CR-LAUNCH, CR-HEADLESS, CR-RESUME, CR-PRESET, CR-PROVISION.
Per-milestone non-vacuity: each milestone's headline test was killed by that milestone's own mutants (M4a UPG-NOMR; M4b PENDFILTER/CAPGE/PRUNECONS/ID;
M4c ROWS/CLOCKPRE/NOTX/DIGESTW; M5a FM-*/FR-*/BM-*; M5b CR-*; M6a AU-*/PE-*; M6b AP-*; M6c CLOCKREUSE/PROMPTCHECK/NODERIVE).

Non-blocking findings:
1. cmd/ailang-worldd/mint_request.go createMintRequest: the ErrTooManyPendingRequests -> exitUsage branch (vs exitFatal) is untested at the
   cmd level (mutant `false && errors.Is(err, store.ErrTooManyPendingRequests)` survives; not in the plan's 83). Store-level cap is tested.
2. QUICKSTART s14 is 12 lines (plan ~45): no expected-output blocks or --list sample; marked pending first verbatim run, so acceptable.
3. M5a CL14 positive control asserts presence of the "unmapped field git"/"stale key" lines, not the exact list; weaker but the real
   completeness check plus FM-* mutants kill all map drift.

Security review (by reading + mutants): consume and session insert share one tx with guarded UPDATE + RowsAffected==1; Clock read once after BeginTx
and passed to Check, SET and WHERE; digest, consumed_at IS NULL, created_at<=at, expires_at>at all in the WHERE; MintMatches checks all seven columns
with RequestCreatedAt; approve decodes+admits before any terminal, same-store check, prepareMint(fromRequest) + byte-equal re-derivation with approval.Diff
naming the field; answerBy only set by approve (existing callers pass 0, unbounded, after never called); request creation opens no terminal.

Deviations (a)-(i): (a) ST7 small-store arm justified, PENDFILTER/PENDEXP killed only with it in the CapsAndPrunes test: ACCEPT.
(b) unwrapped error: insertSessionRow already wraps "store: mint session:": ACCEPT. (c) see NB3: ACCEPT. (d) inline v3 DDL: ACCEPT (independent of production constant).
(e) realSessionEnv().random in the rig only draws unique ids; harmless: ACCEPT. (f) 1610 = now+ttl as provisionEpisode sets; end-to-end test passes: ACCEPT.
(g),(h) did not disturb the killed AP-ADDR / stale-world arms: ACCEPT. (i) 13 arms >= 12 fields; QUICKSTART brevity NB2: ACCEPT.
F1-F8: exit codes applied consistently (usage=1 refusals, fatal=2 at/after write; AP/CR arms assert them).

Worktree restored: git status --porcelain empty; production files cmp-identical to the backups.
EVALUATION_VERDICT: PASS score=90
