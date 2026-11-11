# Evaluator report r1 — iteration 381, row 201 A2 sprint 3 (base ab4b2402, head 64a239f8)

## Verdict: PASS — 92 / 100. No BLOCKING findings.

| Category | Score |
|---|---|
| Correctness vs design and plan §9 | 28 / 30 |
| Tests (non-vacuous, kill mutants, operand/mode coverage) | 26 / 30 |
| Security / fail-closed | 19 / 20 |
| Coding standards | 9 / 10 |
| Scope and size | 10 / 10 |
| **Total** | **92** (PASS threshold 70) |

## Review (full `git diff ab4b2402 HEAD -- cmd host` read)
Confirmed in code, not just in tests:
- Challenge from `env.random`; nil, error or all-zero refused (WN5) before `openApproverTerminal` (approver_passkey.go ~L78-L82). Link carries exactly challenge, rp, name, origin.
- Code compare: `subtle.ConstantTimeCompare` over the whole normalised string against `EnrolmentCode(...)` of the file's own parsed bytes, one attempt (WN10). The registration is read once after the code prompt and the same bytes feed `NewPasskeyApprover` and the code, so a swap before the read is caught by WN10 and there is no second read to race.
- Precheck handles (`passkeyEnrolPrecheck`) are opened and returned before any terminal; WN11 really re-reads (second call with `&a`; mutant passing nil is killed). Writer lock (`openCredentialsForSession`) is taken only after `y`.
- Terminal closed by `defer` on every path after open; M6a split keeps typed-`y` path messages and exit codes (A1b tests unedited, mutants MUT-CONFIRM-* killed).
- `approve --webauthn-assertion`: WF2-WF6 all run before step 2/4 (no daemon/network call) and before `provisionEpisode`; the `signed != nil` path skips the terminal; request stays pending on refusal. `approved_by` is derived by the authority layer from the signed approval; no caller field.
- F1: all five modes in one slice, `modes > 1`; `--list` + any mode refused by F2.
- `ParseStoredCOSEKey`: `decodeCBORExact` + `ParseCOSEKey` (rejects trailing bytes and non-canonical input).
- Scope: `git diff --name-status` equals plan §5 (plus plan doc and sprint json); nothing under tools/launchd, world/, packages/, scripts/, world-publish, skills. `host/approval` imports: stdlib only (go list shows no non-stdlib or module imports).
- Size: `git diff --numstat ab4b2402 HEAD -- cmd host` sums to **1686** (insertions+deletions) <= 2000.

## The cut (M7 page, QUICKSTART §16 deferred)
Printed links are `<origin>/surface/passkey#<base64url JSON>`. `host/daemon/daemon.go` registers only exact `/surface`, `/surface/index.json`, `/surface/static`, `/surface/surface.js`, `/surface/app.js`, `/surface/app.css`, `/surface/model.json`, etc., with no subtree or catch-all pattern (grep; not a live request), so `/surface/passkey` is a 404 at HEAD. Safety of the committed state: an operator following this CLI alone gets a 404; no registration file appears, so enrol waits on the code prompt and refuses after 300 s (or on EOF) enrolling nothing, and `--webauthn-link` is followed by an assertion step that cannot be satisfied without an enrolled passkey (WF1 refuses when none). The origin is validated by `ParseCeremonyOrigin` (localhost or *.ts.net, no IP) and chosen by the operator, and the fragment (challenge) never reaches a server. I found no way this misleads into an unsafe action or opens someone else's page beyond what the operator's own `--origin` names. Only cosmetic: help text says "the daemon's page /surface/passkey" as if it exists; follow-up M7 resolves it.

## Findings (all non-blocking)
1. Test gap: the `then approve it:` line's `--db`/`--credentials-db` operand order is untested. OWN-LINK-NEXTSTEP-DB (swap them) SURVIVED; test only asserts the prefix `... --db ` (approve_passkey_test.go:148). Effect if shipped: operator is told wrong args, gets a refusal (fail closed).
2. Resource-handle gap: removing `defer ro.Close()` in `passkeyEnrolPrecheck` (OWN-PRECHECK-RO-CLOSE) SURVIVED; no test observes a leaked read-only handle. Read-only handles do not hold the writer lock, so no safety effect, but "no handle open while the prompt waits" is not enforced by a test for the RO handle.
3. Normaliser looseness untested: accepting `_` as an extra separator (OWN-NORM-UNDERSCORE) SURVIVED. Not exploitable (compare still needs the genuine 26 chars); design says only `-`, space, tab are stripped (the code also strips `\n`, needed for the line).
4. OWN-ALLZERO-ONLY (`challenge[0]==0` instead of all-zero) survives: strictly stricter (fail closed), arguably equivalent for tests.
5. Plan table's three "named equivalent" survivors reproduced (MUT-ASK-ERR-EQUIV, MUT-AS-CRED-EQUIV; MUT-ASK-NL-EQUIV was not run).
6. `--webauthn-assertion ""` (empty value) falls back to the typed-y path, same as pre-existing `--signature ""`. Consistent, not new.
7. Plan §4 multi-line mutants (rows marked with a newline sign) cannot be applied by a single-line harness; I rewrote the relevant ones as OWN-* with real newlines (see below).

## Mutants (ledger: `$HOME/.ailang/state/world-iter381/judge/mutants.tsv`, 67 lines; harness mut.py/own*.py alongside)
Totals in ledger: 50 KILLED, 12 SURVIVED, 1 COMPILEERR, 3 SKIP (anchor not found, multi-line plan rows, redone as OWN-*), 1 control line.
Plan-table sample (single named test + compile check each), 35 executed, spread over M6a/b/c/f/g:
- M6a: MUT-OPEN-AFTER, MUT-CONFIRM-COMPLETE, MUT-CONFIRM-BUDGET all KILLED; MUT-ASK-ERR-EQUIV SURVIVED (named equivalent).
- M6b: MUT-L-FORMAT, MUT-L5, MUT-L4-ALG, MUT-CRED-LEN KILLED (MUT-PSK-DECODE/WRAP were skipped as multi-line; replaced by OWN-PSK-ERR and OWN-PSK-ERRWRAP, KILLED).
- M6c: MUT-WN1-PUBKEY, WN1-AND, WN2-REG, WN4, WN5-ZERO, WN5-ERR, WN6-CONST, WN10, WN10-PREFIX, WN10-UPPER, WN10-TAB, WN11-KEY, WN11-CRED, N12-CRED, LINK-CHALLENGE, ROW-KEY all KILLED.
- M6f: MUT-F1-LINK, WF1, LINK-ALLOW, HASH-EXPIRY KILLED.
- M6g: MUT-F1-ASSERT, WF-ASSERT-REFUSAL, WF2-CAP, WF5-KEY, AS-REVOKED, AS-COUNT KILLED; MUT-AS-CRED-EQUIV SURVIVED (named equivalent).
Own mutants from the diff (22 + 1 rerun): OWN-PSK-ERR, OWN-PSK-ERRWRAP, OWN-PRE1-REFUSAL, OWN-CODE-INVERT (compare `!=1` to `==1`), OWN-LINK-PATH, OWN-YN-GATE (write without y gate), OWN-F1-GT2, OWN-WN1-ORIGIN, OWN-WN1-REG, OWN-WN1-FIRST (each WN1 operand separately), OWN-WN11-NILA (second precheck without the registration), OWN-TERM-CLOSE (terminal not closed), OWN-BYTESDIGEST, OWN-STOREDDIGEST-RERUN (first attempt OWN-STOREDDIGEST was COMPILEERR, unused var; rewritten and re-run, KILLED, both lines in ledger), OWN-REGPATH-STAT (Lstat to Stat), OWN-LINK-ORIGIN-FIELD, OWN-CTRL-MUSTDIE-BUDGET (300 s to 120 s) all KILLED. SURVIVED: findings 1-4 above.
Controls: must-die OWN-CTRL-MUSTDIE-BUDGET and MUT-WN10 KILLED; must-live OWN-CTRL-MUSTLIVE-NOOP (`300*time.Second + 0`) SURVIVED as required; unmutated named tests pass (CTRL-UNMUTATED rc=0).
Worktree restored after every mutant; final `git status --short` empty, HEAD 64a239f8.

## Per-milestone non-vacuity (output: judge/nonvacuity.txt)
- M6c: parent 41f5cb9a prod + M6c tests: build FAILED (`undefined: ceremonyLink`).
- M6f: parent M6e prod + M6f test file: `TestApproveByPasskey` FAIL (`flag provided but not defined: -webauthn-link`).
- M6g: parent M6f prod + M6g test file: `TestApproveByPasskey` FAIL (`-webauthn-assertion` not defined).
- M6d/M6e (test-only): with those two test files deleted at HEAD, NOMDE-WN10-INVERT, NOMDE-WN10-PREFIX, NOMDE-WN11-NILA, NOMDE-N12-CRED, NOMDE-YN-GATE all SURVIVE the remaining suite (`TestPasskeyEnroll.*|TestApprover.*`); with the files restored the same mutants are KILLED (MUT-WN10, -PREFIX, WN11-KEY/CRED, N12-CRED, OWN-YN-GATE, OWN-WN11-NILA, OWN-CODE-INVERT). So M6d/M6e are what pin CW2/CW4/CW5/CW6.

## Test runs (AILANG_BIN=$HOME/.pinned-ailang/ailang; logs in judge/run_*.txt)
- `go test -count=1 ./host/approval/... ./host/store/ ./cmd/ailang-worldd/` rc=0 (cmd 129 s)
- `go test -race -count=1 ./host/approval/... ./host/store/` rc=0
- `go test -race -count=1 ./cmd/ailang-worldd/` rc=0 (241 s)
- `go vet ./...` rc=0; `go test -run '^$' ./...` rc=0
- Not run by design: verify_go.sh, whole `go test ./...`.
