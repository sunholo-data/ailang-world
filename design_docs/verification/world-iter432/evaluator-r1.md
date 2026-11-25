# Evaluator report, iteration 432, row 243 V4 S5b (head 6e354749) — round 1

VERDICT: PASS. Score 84/100. No BLOCKING findings.

| Category | Score |
|---|---|
| Correctness vs design and plan §9 | 26/30 |
| Tests | 24/30 |
| Security / fail-closed | 18/20 |
| Coding standards | 9/10 |
| Scope and size | 7/10 |

## Measurements
- Numstat (`git diff -M --numstat f8d8e3ad HEAD -- cmd host docs app`) sums to exactly 1,400. The rows are: operator_mint.go 126, operator_mint_test.go 341, operator_provision_test.go 84, operator_service.go 11, service.go 2, session.go 11, QUICKSTART 13, operator_signed.go 46, authority/mint.go 23, authority/operator_mint_test.go 383, operator_store_test.go 5, credentials.go 23, credentials_test.go 3, store/operator_mint.go 324, episode_lock_test 5. The plan's unit is honest. The cap is met with zero headroom, which is why the resolver's defence-in-depth checks are thinly tested (see N1).
- Production callers (non-test): `MintOperatorStart` 0, `CommitOperatorMint` 1 (only from MintOperatorStart), `AbortOperatorMint` and `ReadOperatorMint` 1 each (only from the attended resolver), `MintProof` 1. `ClaimOperatorAction` has a definition only and no production caller. `ReserveOperatorStep` and `ObserveOperatorStep` do not exist under those names. The transaction is inert, and binding gate 2 is respected.
- `.ail`, app, `daemon.go` and `surface_page.go` are unchanged. The controller's `routes.txt` and `scope.txt` legs are rc=0.
- Controller sweep: `observation-before` and `observation-after` are identical. `verify_go` rc=1 comes only from `host/daemon` `TestSeToolsRunBenchmarkTasksEndToEnd`, which is on the known load-sensitive list. I re-ran it alone and it passed (ok 6.9s). Every other leg is rc=0, including the 347 s `-race` of `cmd/ailang-worldd`.
- My test runs:
  - `go test -count=1 ./host/store ./host/authority ./host/approval` rc=0.
  - The cmd broad regex `Operator|Episode|Preset|...|Runner` rc=0, 70.9 s.
  - `-race` on `OperatorSessionMint|OperatorMint|OperatorProvision` across cmd, authority and store rc=0, 6 s.
  - `go vet ./cmd/... ./host/...` rc=0.
  - `go test -run '^$' ./...` compiled clean.
  - The new tests run in about 1.5 s for cmd and about 2 s for authority.
- CI cost: negligible. There are no `Sleep` or `time.After` calls in the new tests (grep is empty).
- Linux portability: the new tests use `syscall.O_NOFOLLOW`, `Stat_t.Uid` and `Nlink`, `os.Link` and `os.Symlink` in temp dirs, and mode 0600. A symlink open under `O_NOFOLLOW` fails with ELOOP on Linux and is handled as a generic open error. I found nothing macOS-specific.

## Answers to the questions
- **M1 order.** The order is: arg and store checks, then the tty fence, then `terminalEpisodeLock` (only when the id matches `workspaceEpisodePattern`), then the prompt, then store open and `authority.Mint`, then the token write. `defer lock.Close()` covers every return. The refusal comes before the prompt, is non-blocking, names the action id (`expectEpisodeRefusal`) and asserts unchanged credentials, worktree, refs and files. The `durably_owned` subtest closes the descriptor first and still refuses. Opaque ids such as `legacy/opaque` skip the lock and mint as before (`free_and_opaque`). Callers of `terminalEpisodeLock` are now `episode`, `finish`, `launch`, `preset`, `resume`, `session_new` and `session mint`. `approve.go:398` calls `provisionEpisode(req, nil, ...)`; this is S5's signed-approve path, which I did not re-audit.
- **M2 atomicity.** One `BeginTx` covers the session row, the start row and the `operator_action_steps` receipt UPDATE, then one `Commit`. The token is generated before the transaction, and no plaintext enters the store. The receipt holds the credential hash plus the start digest. A premature-commit mutant is killed (`TestOperatorMintReceiptWindow/transaction_separation`). The collision test blocks a second mint after a lost handoff.
- **Sealed proof.** `OperatorMintProof{signed}` has an unexported field. `MintProof()` returns nil for a confirmation proof (`s.confirmed`) or an unverified one. `Bind` re-checks the artifact hash and `Binds`. The only production construction path is `MintOperatorStart` after `RecheckClaimedOperatorAction`. In `package approval` a zero-value `OperatorMintProof{}` has `signed == nil` and is refused.
- **Binding.** The table has about 55 subtests, each with its own refusal class or `BindingError` field and an unchanged-snapshot read-back. They cover ttl, base, grants, project, profile, catalogue, runner generation and owner, action owner and phase, approver identity, key, revoked, enrolment, passkey credential, counter, rp, origin and backup, expiry boundary, clock behind, and the token and task paths. The clock is read inside the transaction (`now := clock()` after `BeginTx`).
- **Resolver safety.** It needs a typed `y` through `e.confirm`, which requires a tty, and the flow is:
  1. `installedRunner`, so it refuses on a machine with no installed runner, before any write.
  2. Read the receipt.
  3. Take a non-blocking episode lock.
  4. Re-read under the lock.
  5. Inspect the token. A matching token retains the owner and is never revoked.
  6. Prompt.
  7. Inside `runnerWrite`, re-inspect the token. If it appeared after the prompt, refuse.
  8. `AbortOperatorMint` in one transaction: re-read, `current != reviewed` gives conflict, delete only that credential id, fail the action, and clear the owner. All of this is gated on the exact owner, generation, digest and step version.
  Tests cover declined, no-tty, stale owner, stale step, token appeared, later effect, expired, a rollback trigger, and an unrelated session, worktree, task and MCP files unchanged. The remedy text is exercised verbatim: the test cuts the command out of the refusal and re-runs that exact string. QUICKSTART is checked by `TestQuickstartRunnerMintRecovery`.
- **Token readback.** `O_NOFOLLOW|O_NONBLOCK`, an fstat regular-file check, uid, Nlink==1, mode==0600, `LimitReader(65)`, an exact 64 bytes and a sha256 match to the credential id. This is correct and bounded.
- **S5b/S5c split.** Appendix C carries S5c (worktree reservation and probe, files, live-input reconstruction), the S5c/S6 paired-process consumer and drain with four named process pins, `session revoke`, and the stop-before-activation gate. I did not re-verify that every S5 Appendix C share was carried: the two parked parity survivors and J-OWN-UID were not located in the grep I ran (non-blocking, N3).
- **Pre-existing edits.** The pre-existing edits do not weaken anything. Specifically:
  - `credentials.go`: a pure extraction into `insertEpisodeStart`, behaviour-identical.
  - `operator_signed.go`: `confirmation *` becomes `binding` plus `confirmed`; `ConfirmationFor` keeps the same checks.
  - `authority.Mint` is unchanged, since `mint.go` only adds a function.
  - `operator_service.go` and `service.go`: one verb and help-string edit.
  - Test edits only add helpers and arms.

## Mutation results

Ledger: `$HOME/.ailang/state/world-iter432/judge/mutants.tsv` (76 lines). The harness is in `judge/scr/`.

**Plan rows I re-measured (17, all from the controller's UNMEASURED CLI list; every one KILLED by a behavioural Go test):**

| Plan id | Killing test |
|---|---|
| L-FENCE | `TestOperatorSessionMintEpisodeExclusion/durably_owned` |
| L-OPAQUE | `TestOperatorSessionMintEpisodeExclusion/free_and_opaque` |
| L-CLOSE | `TestOperatorSessionMintEpisodeExclusion/declined` |
| NV-M1 | `.../durably_owned` |
| L-TTY-ORDER | `.../tty` |
| R-CLOSE | `TestOperatorMintRecovery/revocation` |
| R-LOCK-CREATE | `TestOperatorMintRecovery/missing_lock` |
| R-TTY | `TestOperatorMintRecovery/revocation` |
| R-TOKEN-RECHECK | `TestOperatorMintRecovery/appeared` |
| R-TOKEN-FIRST | `TestOperatorMintRecovery/adopt` |
| R-EXPIRED | `TestOperatorMintRecovery/expired` |
| R-REVOKE-ID | `TestOperatorMintRecovery/revocation` |
| R-ABORT-ERROR | `TestOperatorMintRecovery/abort_rollback` |
| F-NOFOLLOW | `TestOperatorMintTokenReadback/symlink_valid` |
| F-READ-CAP | `TestOperatorMintTokenReadback/length` |
| F-EXACT-LENGTH | `TestOperatorMintTokenReadback/short` |
| R-LOCK-NB | `TestOperatorEpisodeLockExclusive` |

**My own mutants.** Ids are prefixed E-, E2-, E3- and E4-, so none are in the plan's id list; the `comm -23` against the plan ids is trivially my whole set. Some E- lines in the ledger ran against the wrong package (host/store has no mint tests; the tests live in authority and cmd). Those are superseded by the E2- reruns and should be read as mis-targeted, not as results.
- KILLED, with the killing test:
  - E2-TX-HEADLESS: `OperatorMintBinding/start_options`
  - E2-TX-APPROVEDBY: `/approver_identity`
  - E2-TX-NOW-LT: `/clock_behind`
  - E2-TX-COLLISION: `OperatorMintCollision/credential`
  - E2-TX-STEP-PHASE: `/step_phase`
  - E2-TX-RECEIPT-START: `OperatorMintReceiptWindow` and `OperatorMintRecovery` arms
  - E2-TX-PASSKEY-COUNT: `/passkey_counter`
  - E2-TX-TTL-ZERO: `/zero_ttl`
  - E2-ABORT-OWNER-KEPT: `OperatorMintRecovery/revocation`
  - E2-ABORT-GEN: `/revocation`
  - E2-ABORT-NOCAS: `/stale_owner` and `/stale_step`
  - E2-ABORT-EPISODE-COL (a delete-by-episode-only variant): `/revocation`
  - E2-READ-LATER-LAUNCH: `/later_effect`
  - E2-READ-STATE-UNCERTAIN: `/revocation`
  - E2-SIG-MINTPROOF-CONFIRMED: `OperatorMintBinding/proof`
  - E2-SIG-CONFIRMFOR: `OperatorRunnerRechecks/claimed_store_row`
  - E2-SIG-BIND-ARTIFACT: `OperatorMintBinding/artifact`
  - E-MINT-THIN: `/thin_challenge`
  - E-RES-NLINK: `TokenReadback/link`
  - E-RES-MODE: `/mode`
  - E-RES-RETAIN-NOW: `/expired`
  - The M1-hunk revert (base `session.go` restored) fails `TestOperatorSessionMintEpisodeExclusion` (non-vacuity for M1 holds; the ledger line for it is a manual note without a recorded rc).
- Compile errors (not mutants, fixed and rerun): E-TX-RECEIPT-START, E-RES-SHA. E2-RES-SHA also failed to compile (a bad escape sequence in my mutation spec, not a result) and I did not rerun it, so the token sha256 comparison has no mutation evidence from me.
- Controls: must-die is L-FENCE (killed). Must-live is E4-CTRL-LIVE (`0600` to `0o600`), which survived as expected.
- SURVIVED (whole-package run where it was affordable):
  - E3-TX-STORE-ID-whole (remove `a.Store != c.lock.dbPath`): host packages whole, rc=0.
  - E3-TX-KIND-whole (remove `row.Kind != "start"` in the transaction): host packages whole, rc=0.
  - E4-RES-STORE-ID-whole (remove `r.Store != c.Credentials` in resolver inspect): whole cmd package, rc=0.
  - Four more survived the broad regex only and were not run whole, so they are not strictly SURVIVED verdicts: E3-RES-TOKENPATH-BASE-broad, E3-RES-EPISODE-CONFLICT-broad, E3-RES-REREAD-broad and E3-RES-REGULAR-broad. The last, `IsRegular`, is practically redundant because the later mode check rejects a directory. The re-read ones are covered by the transaction CAS.

**M2 per-hunk non-vacuity.** Each M2 hunk I mutated (transaction guards, abort SQL, proof sealing, readback) was killed by a named M2 test; the three survivors above are the exceptions.

## Findings (all non-blocking, ranked)
- **N1.** Three defence-in-depth guards have no test: the transaction's store-identity check (`a.Store != c.lock.dbPath`), the transaction's `Kind == "start"` check, and the resolver's token-path and store-identity checks in `inspectRunnerMintToken` (basename, `r.Store`). They are partly redundant with the recheck and the action digest, but a regression would pass CI. Fix: add one binding case per guard (`operator_mint.go:97-104`, `host/store/operator_mint.go`).
- **N2.** The test for `reconcile-mint` on a machine with no installed runner is not explicit. `installedRunner` plainly refuses before any write, and the refusal paths are covered indirectly, but a one-line test would pin the "today's real install" claim.
- **N3.** I did not exhaustively verify that Appendix C re-carries the parked parity survivors and J-OWN-UID from S5. The S5c items are carried; check these two named items before S5c planning.
- **N4.** The resolver holds the episode lock through the human prompt. This is intended (the lock is non-blocking for others, which refuse), but an operator who leaves the prompt open blocks terminal episode operations on that episode.

## Verdict rationale
M1 is correct and non-vacuous. M2 is atomic, sealed, bound field by field with state-unchanged read-back, and inert (zero production callers of the mint). The resolver is safe: a typed `y` at a tty, an exact owner/step CAS re-read after the prompt, and a scope limited to the orphan credential. The measurement is honest at 1,400. The only weaknesses are the untested redundant guards (N1).
