# Evaluator r1: row 201 A1b sprint 4 (base 1fce6999, head 5e85f8bb)

VERDICT: PASS, 93/100 (correctness 28/30, tests 28/30, security 19/20, standards 9/10, scope 9/10)

## Blocking findings
None.

## Non-blocking findings
1. cmd/ailang-worldd/approver.go (K8 read-only open failure, "revoked, but cannot read the credentials store ... exitFatal"): untested. Mutant K8-ERRCODE (exitFatal -> exitOK on that branch) SURVIVED the whole TestApproverRevoke* set. The sibling list-error branch is covered (K8-ROOT killed).
2. host/sshsign Sign: G1 (failed run) is killed only from cmd (TestSignWithSSHKeygenRefusesBadOutput/failed), not by a host/sshsign unit test. Package-level coverage is fine; unit-level is thin.
3. signWithSSHKeygen resolves "ssh-keygen" through PATH, which the scrub does not narrow. Documented, and safe because approveArmored verifies the output against the enrolled key over World's own message (SIGN-NOVERIFY mutant would not compile, so it could not be exercised, but SIGN-SOURCE was killed).
4. QUICKSTART §15 adds the --sign-with / --list --decided lines; they are consistent with the code but I did not execute them (attended steps).

## Real-tool leg (WORLD_REQUIRE_SSH_KEYGEN=1)
PASS, not SKIP: TestSSHSIGCrossCheckWithSSHKeygen (+ed25519, ecdsa-p256), TestTestSignerAgreesWithSSHKeygen, TestSignWithSSHKeygenEndToEnd (+ed25519, ecdsa-p256). The CI step's PASS-grep list matches these names. Touched packages (cmd/ailang-worldd, host/approval/..., store, authority, sshsign, broker incl. AC10 census, childenv) all ok at head.

## Mutants (run against the named single test; survivors re-run narrowly)
Controls: CTRL-KILL (drop Scrubbed line) KILLED by TestSignHandsTheSignerAScrubbedEnvironment; CTRL-SURVIVE (len(id)+0) SURVIVED. Harness: perl literal replace, non-empty git diff asserted, KILLED needs a real "--- FAIL".
KILLED (test): DM-FILTER, DM-ORDER, DM-TIE, DM-APPROVEDBY(+SQL), DM-CAP, DM-LIMIT-HIGH, DM-LIMIT-LOW, DM-VERSION (v3 arm), DM-DEADLINE (no-deadline arm) (all TestDecidedMintRequests); F3 and F2/F1-COUNT/F1-NOSIGNWITH (TestApproveFlagRules); F3-LIST, DECIDED-BRANCH, DECIDED-FILTER, F15, F16-CONSUMED, F16-BY, F16-CRED17 (TestApproveDecidedList); F7, F8, SIGN-KEY, SIGN-SOURCE (TestApproveSignWithFakeSigner); F7A, F7A-OFF (rem<0 variant also killed), DEADLINE-UNIT (TestApproveSignWithDeadline); REAL-ENV-SIGN (TestRealSessionEnvSignsWithSSHKeygen); SIGNER-BIN, G1, G2, G3, CAP-MAX, ARGS-KEY, STDIN, STDERR, WAITDELAY, CMD-CTX (TestSignWithSSHKeygenRefusesBadOutput); CAP-WRITE (TestCappedBufferKeepsItsBound); ENV-SCRUB and ENV-EMPTY (sshsign unit), ENV-SCRUB-CENSUS (broker census, subtest host/sshsign/sshsign.go); K8-FILTER, K8-ROOT, K9, K8-LINE, K8-AT, K8-COUNT (TestApproverRevokeListsWhatItApproved).
SURVIVED: K8-ERRCODE (finding 1). First attempts for DM-VERSION/DM-DEADLINE/F3-LIST read SURVIVED only because my pattern hit PendingMintRequests or I named the wrong test; re-anchored, all three KILLED.
Not mutants (did not compile): SIGN-NOCTX, SIGN-NOVERIFY, ARGS-NS (my sed form). Not counted.

## Security properties confirmed
- Child env is childenv.Scrubbed(os.Environ()) (registry credential absent; SSH_AUTH_SOCK/HOME/PATH pass); unit test plus AC10 census drives sshsign.
- argv is fixed ("-Y sign -n world-approval@sunholo -f <key>"), key a separate argv element (no shell, no option injection); message on stdin; nothing written beside the key.
- Output capped at 16 KiB (capped writer returns len(p), no SIGPIPE games), empty and failed runs refused; output is never trusted: approveArmored parses, looks up the enrolled unrevoked key, and the law verifies over World's own message before any mint.
- Deadline = remaining request lifetime; expired-since-admission refused before the signer runs (F7a); WaitDelay 1s bounds a pipe-holding grandchild (killed in ~2s vs ~20s without).
- approver revoke: write root closed before K8; K8 uses its own read-only store and root; it revokes no session; only prints the `session revoke` lines.
- --list/--decided/--signature/--sign-with/--write-message exclusivity (F1/F2/F3) holds.
- D-WORLD-109 = A holds: only temp stores, fake signers, test-generated keys with SSH_AUTH_SOCK emptied and temp HOME; grep of added tests finds no ~/.ssh, /dev/tty, ~/.ailang/world or UserHomeDir use.
- coding-standards: effects at host boundary (sshsign), context_roots pins updated for the new Background roots.

## Worktree
/Users/voightkampff/dev/sunholo-data/.wt-world-iter378-judge: `git status --short` empty after every mutant (checked). Sprint worktree and main checkout untouched. Scratch: ~/.ailang/state/world-iter378/judge/.
