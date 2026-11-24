# Evaluator r2 - iter 429, row 243 V4 S4 (head 976f4028)

## Verdict: PASS. Score 81/100. No BLOCKING findings.

| Category | r1 | r2 |
|---|---|---|
| Correctness vs design / plan 9 | 24 | 27 |
| Tests (non-vacuity, mutants) | 18 | 23 |
| Security / fail-closed | 15 | 16 |
| Coding standards | 8 | 8 |
| Scope and size | 7 | 7 |
| Total | 72 | 81 |

## What changed (verified)
`git diff --stat 22473391 976f4028`: exactly two test files, +110/-0 (operator_service_test.go +94, host/store/operator_runner_test.go +16). No production file changed. Nothing else found.

## Scope trade
Size is now 1,510 against the 1,400 cap. I agree with the trade: the overage is entirely tests pinning an authority gate (attended-cancel confirmation) and the home-upgrade fence that r1 proved was shipped undefended; design 8 says never cut an authority gate to make the count. Disclosed, test-only, no production growth. Scope scored 7/10 (cap exceeded, justified).

## B1: CLOSED
Whole-package re-runs at 976f4028 (ledger `$HOME/.ailang/state/world-iter429/judge/mutants.tsv`, 46 lines):
| id | whole-package result |
|---|---|
| R2-J-CANCEL-CONFIRM (drop `verb == "cancel"` from typed-y set) | KILLED by TestOperatorServiceCancelConfirmation |
| R2-J-UPGRADE-SUPERVISOR (drop supervisor-unloaded check) | KILLED by TestOperatorServiceUpgradePairFence/registered |
| R2-J-UPGRADE-RUNNER-PLIST (drop Plist half) | KILLED by TestOperatorServiceUpgradePairFence/plist-only |
| R2-J-UPGRADE-RECORD (drop Record half; extra) | KILLED by .../installed |
| R2-J-UPGRADE-LOCK (lifecycle lock -> plain os.Open; extra) | KILLED by .../locked |
| R2-J-ST-CANCEL-CONFIRMED (whole host/store) | KILLED by TestOperatorCancelConfirmed |
All kills are runtime assertions. Each of the four upgrade arms fails with only its own guard removed (installed, plist-only, registered, locked each killed by a distinct mutant), so the arms are non-vacuous. The new tests assert both refusal and unchanged state: exit 1, expected message, zero kicks, binary hash unchanged, zero observed writes (upgrade); row `requested` after `n`, `cancelled` after `y` via a read-only handle (cancel); generation+1 and terminal_at readback (store). The previously untested strings "still installed", "lifecycle is busy" and "cannot verify that the operator runner is unloaded" are now asserted.

## r1 findings by name
- B1 (J-CANCEL-CONFIRM, J-UPGRADE-SUPERVISOR, J-UPGRADE-RUNNER-PLIST, untested still-installed / lifecycle-busy refusals): CLOSED.
- Non-blocking 1 (subset survivors J-OWN-UID, J-WRITABLE-OTHERS, J-LOCKFILE-MODE, J-LOCKFILE-NOFOLLOW, J-CONFIG-NOFOLLOW, J-PRESET-SE, J-DISTINCT-STORES, J-HEALTH-STATUS, J-LOOPBACK-IP, J-ADDR-GIVEN): OPEN (not targeted; not re-run). Non-blocking: install-time defense-in-depth checks behind other validators or untestable as non-root, nothing dispatches an effect.
- Non-blocking 2: J-ST-CANCEL-CONFIRMED CLOSED. J-ST-CANCEL-OWNER-ISNULL OPEN (plan calls it equivalent, unreachable state). J-ST-HB-ACCESS OPEN (heartbeat read deadline unpinned).
- Non-blocking 3 (compatible rollback covered only by pre-existing test), 4 (refusal leaves candidate running, record drift), 5 (credentials_schema is self-reported, omitempty fail-closed): UNCHANGED.
- r1 unmeasured items: per-commit green of M1 still unmeasured; r1 verify_go failures (host/daemon, host/verifygate, untouched packages) are gone: round-2 sweep verify_go rc=0.

## Test runs
- Controller round-2 sweep (`gates-r2/legs.txt`): all legs rc=0 (pkg cmd 225 s, verify_ail, verify_go, fence, scope, observation); `GATES_DONE 670s FAILS: none`.
- Mine: `go test -count=1 ./host/store/` ok rc=0; `go test -count=1 -run 'Operator|HomeService|HomeVersion|DoctorService|Quickstart|Help|Keyguard|CredentialWriters' ./cmd/ailang-worldd/` ok rc=0. Mutant runs above used `go vet` before test (compiled) and whole `./cmd/ailang-worldd/` / `./host/store/`.
- Eval worktree clean at the end; throwaway worktrees removed.
