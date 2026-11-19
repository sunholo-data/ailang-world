# Ailang World Mission — iteration log

Append-only. One entry per outer-loop iteration. Newest at the bottom.

---

> **Older entries are ARCHIVED.** This file holds the newest 20. The full record of every
> iteration is in `world-mission-log-archive.md`, and a one-line index of ALL of them —
> the thing to grep before picking work, so the loop never repeats itself — is in
> `world-mission-index.md`.

## 390 — 2026-11-12 — row 201 A4 sprint 2 LANDED (PR #566, merge `9ea7b45f`): the typed CLI `approve <request-ref> [--deny]`, credentials schema 8 `key_decisions`, the decision SSHSIG law and judge, the POST's `sshsig` arm and `--sign-with` on a request ref; the webauthn arm and CD6–CD7 moved to sprint 3; judge PASS 88 [PRODUCT]

**Pick:** the A4 sprint-2 planner. It was record 389's Next 2 and the charter's Next loop row. Record 389 had landed (#564 → `8338eebb`), so there was no carry-over PR. D-WORLD-115 is unanswered and its default A still applies.
- **Re-measured at pick** (`origin/dev` `8338eebb`):
  - There was no A4 sprint-2 plan, branch or worktree, and `gh pr list --author @me --state open` was empty.
  - `credentialsSchemaVersion = 7` (`host/store/credentials.go:43`).
  - `approve.go:199` refuses any non-`mr-` id.
  - `surface_decide.go` admits only the typed arm.
- **Quorum:** none owed (the A4 r2 design covers sprint 2).
- **Ledger:** valid, 102 rows before and after. OPEN: D-WORLD-105 (default C), D-WORLD-114 (default A, applied), D-WORLD-115 (default A, applied).
- **Clause map:** 1.0 clauses 1–7 MET, none UNMET. Release rows 150/151 are not routable (unchanged).
- **Orphans:** none.

**Gate 0:**
- Stamped first. `gh` is `sunholo-voight-kampff`. Billing CLEAN. Kill switch absent. Dev CI 3/3 at `8338eebb`.
- Directives: 0 allowlisted on #524 (26 comments), read through both the issue-scoped watermark (absent) and `mission-world-last-seen` (LEGACY, 2026-10-27).
- Self-notices: rc 2, `watermark channel unavailable`. This is an instrument failure, not a verdict, as in 371–389.
- Inbox: none addressed to `mission-world`. The global unread list (3,060) is other agents' traffic.
- Skill drift: none. `SKILL.md` and all 13 `resources/*.md` match the driver clone's `origin/dev` ref, at both the resolved symlink target and the pin worktree (two different inodes, both clean).
- Weekly issue sweep: not owed (#524 was opened Monday).

**Planner** (Agent `opus`, 77 min, 536,608 tok): plan `w-world-approvals-by-key-a4-sprint2-plan.md` with 10 compiled patches in Appendix P.
- **Size:** 1,983 lines (1,929 insertions + 54 deletions; 738 production, 1,245 test) by `git diff --numstat 8338eebb HEAD -- cmd host`, 17 under the cap.
- **Cut (S2-D1):** the design's priced fallback fired. The whole webauthn arm (DM7, the passkey halves of RD10/K0/KD7/KD8/RD13/RD14, KD2, CD8–CD9, `DecisionSubject.Passkeys`) moved to sprint 3. The build was still 2,019 with CD6–CD7, so `--write-message` and `--signature --expiry` on a request ref (63 lines) moved too. Never cut: CD1–CD4, R4e, RD13 (a), KD5/KD6, schema 8, the judge, and the J-gap tests.
- **Mutants:** 125 rows. 109 were killed by both the named test and the package suite, 8 are named equivalents and 6 are named survivors. One row (`S-K1-DEADLINE`) is killed by an inventory test only, and both results are recorded. Ledger: 135 lines.
- **Premises found false (plan §1.4):**
  - **P1:** `sshsign.Sign` hard-codes the mint namespace, so `--sign-with` could never sign a decision. The real-`ssh-keygen` sibling test found it; the fake signer could not. The fix is `sshsign.SignIn` (namespace parameter; `Sign` delegates) plus a second seam, `sessionEnv.signDecision`.
  - **P2:** schema 8 breaks the downgrade fixtures inside tests the design lists as "must stay green, unedited". Each site gets one `DROP TABLE key_decisions` token, as A3a sprint 1 did for schema 7 (`38a6d903`), and four store tests pin the new table set and version. **Controller decision: accepted.** Schema 8 cannot land otherwise. The judge audited every edit.
  - **P3:** `approve` refuses the global `--addr`, so the decision mode takes `--daemon`.
  - **P4:** `DecideBody.Typed` needs `omitempty`.
  - **P5–P8:**
    - KD2 is passkey-only, so it moves with the webauthn arm.
    - RD11's forged-`decidedBy` 400 already holds in sprint 1.
    - KD9's HTTP mapping is unreachable end to end.
    - No route changes this sprint.
- **The CLI finds the daemon (S2-D2):** through `--daemon <host:port>`, defaulting to the installed service's address. Its `/v1/health` must name the same `--db` and credentials store, the same check `episode finish` F6 makes.
- **Schema 8 is forward-compatible (S2-D3):** the format is `ssh` or `webauthn`, and the passkey counter stays in `webauthn_credentials.sign_count`. Sprint 3 adds Go only; there is no schema 9.

**Controller round trip:** I ran `git archive 8338eebb`, extracted the 10 patches from the plan with §2's awk, and applied them. `diff -rq` against the planner's `work/` tree is empty for `cmd` and `host`; the control (the unpatched base against the same tree) shows 42 differences. The plan was committed (`f6652a40`) before any sub-milestone commit. The final committed tree equals the round-trip tree.

**Executor** (env `recipe claude:claude-haiku-5-5`, spawned as Agent `haiku`, as 374–389), three runs:
- Run 1 (M4a–M4c): 1.8 min, 46,553 tok.
- Run 2 (M4d–M4g): 2.9 min, 55,361 tok.
- Run 3 (M5a–M5c): 9.7 min, 52,940 tok.
- Every directive carried "every Go Bash call passes `timeout: 600000`, foreground, no Monitor" from the start. All three ran clean in one pass: no resume, no hand edits, no git writes, no stop condition, no timing exception.
- The whole-package legs were `ok`: `host/daemon` in 287 s and `cmd/ailang-worldd` in 164 s.

The controller rebuilt each sub-milestone from its snapshot with `commit_ms.sh`; every commit was gofmt/build/vet-gated and ran its own sub-milestone's test command. Cumulative numstat after each commit:

| Sub-milestone | Commit | Numstat |
|---|---|---|
| M4a | `c85af51f` | 55 |
| M4b | `640289ca` | 450 |
| M4c | `307dda55` | 559 |
| M4d | `1e7dcca8` | 777 |
| M4e | `3a4026e9` | 1142 |
| M4f | `7fb0bc84` | 1258 |
| M4g | `447cfe0c` | 1380 |
| M5a | `1a5342db` | 1658 |
| M5b | `4a7f0f04` | 1804 |
| M5c | `e4510375` | 1983 |

That equals the plan at every step.

**Gates** (controller; the plan §6 `gates.sh`, every leg recording its own rc; 788 s at load 16–30):
- `verify_ail`, gofmt, build, vet and the fence: rc 0. `pkgs` (approval, authority, store, sshsign, surfacepage, daemon, cmd): rc 0.
- `broker_named` rc 0. `broker_whole_reported` rc 0: the base's known-list `TestExecProtectedCheckTimeout` red did not recur.
- `race` rc 0.
- `inv` rc 0: **73** named tests print `--- PASS`, none missing, none skipped. That is 49 carried, including `TestFrozenV1RouteTableMatchesMux` (13), `TestDaemonCredentialWriterSites` (two sites), `TestProductionContextRoots` and `TestCredentialsMethodsRequireADeadline`, plus 24 new.
- `GATES_DONE FAILS: none`.
- **After commit, `VERIFY_GO_LEGS=plain verify_go.sh`: rc 0 on the first run.** All 64 required evidence tests passed once each, 62 packages were `ok`, and the tracked-binary gate passed 24 of 24.

**Judge** (evaluator `sonnet` via Agent, 18.1 min, 114,138 tok; own worktree; generator ≠ judge): **PASS 88/100, zero blocking** (correctness 27, tests 24, security 19, standards 10, scope 8).
- **Confirmed:**
  - The handler order is sprint 1's S1-D9 slot: law under `humanMu`, then the read-only judge, `Binds`, the non-waiting writer, `RecordKeyDecision`, `DecideApproval` and `LandKeyDecision`.
  - `decidedBy` is derived as `ssh:<name>`.
  - Unknown-field and `webauthn` bodies are refused with 400 and reach no arm, so the cut is honestly bounded.
  - The verb is bound inside the 91-byte message.
  - Mint and decision signatures cannot substitute for each other.
  - The K0 exact match covers the artifact, packet ref, verb, format, approver and expiry, with `useApprover` run on both paths.
  - The CLI's hash and `requestRef` checks run before any prompt or proof.
  - Mint signing is unchanged, pinned behaviourally by a real-`ssh-keygen` test.
- **P2 audit:** one `DROP TABLE key_decisions` token at each of 14 sites; the table-set and version pins are additive. No assertion is weakened, except that `surface_approvals_test.go` now uses `credentialsSchemaVersion` instead of the literal 7, a non-blocking change in pin form (`credentials_test.go` still pins 8 literally).
- **Mutants:** its ledger has 68 lines. It re-ran 17 plan rows and ran 43 of its own (`comm -23`). The must-die control (`ssh:evil`) was killed. The must-live controls (JM13 `found &&`, JM14 the `DecidedBy` literal) survived as expected. One wrong-package first run is marked INVALID and was re-run.
- **Non-blocking test gaps** (carried to sprint 3):
  - **JM22:** `LandKeyDecision`'s `request_ref` filter is unpinned (tests hold one pending row).
  - K0's approver and format operands are unpinned at the store API (equivalent through the daemon).
  - **JM01–JM04:** the four `isApproverRace` operands map to 422, and that mapping is unpinned. A concurrent `approver revoke` can reach it, and it still fails closed.
  - **JM27:** `LastUsedAt` pass-through is unpinned.
  - **JM46:** the CLI's default daemon address is unpinned.
  - `--daemon` is not restricted to loopback. The operator chooses the address.
- **Controller-reproduced:** JM22 (`WHERE request_ref = ? AND decision_ref IS NULL` → `WHERE decision_ref IS NULL` in `host/store/key_decisions.go:143`) gives `go test ./host/store/ -run KeyDecision` rc 0. The file was restored and the worktree is clean.

**Gate 3b:**
- PR #566 head `e4510375`, run 43326730776 attempt 1: **3/3 success**. Mergeable CLEAN. Closing-keyword scan: 0 hits over the PR body, title and 11 commit messages (control fired).
- `gh pr merge 566 --merge --match-head-commit e4510375…` → merge **`9ea7b45f`** at 23:49:51Z.
- Push CI on the merge: **run 43332117627 attempt 1: 2/3, with `go host build + test gate` red only on `TestWorktreeSnapshotRefusesAConfiguredFilter` ("control: plain git status did not run the clean filter"). That is the known queue row `w-filter-refusal-control-flake` (filed at iteration 285), in a broker test this sprint does not touch: the sprint's only broker change is one test row in `decision_law_test.go`. It was green 5/5 alone on the merged tree locally. `gh run rerun --failed`, attempt 2: **3/3 success**. LANDED**.

- **Carried from record 389 (owed by its note):** record PR #564 (docs only, code tree = `338ff821`), run 43302235014 attempt 1: the race job was red once on host/daemon `TestCommitInvocationReceiptIdentity/A_landed_uncertain_then_B_on_top` (`commit_budget_test.go:323`: commit B got 409 `HeadConflict`). The other two jobs were green. This is instance 2 of that test (iteration 226 was instance 1), in the row-258 class. `rerun --failed`, attempt 2: 3/3 success. The job's two `DATA RACE` warnings come from the gate's own positive control (`design_docs/verification/w-race-gate-blindspot/racecontrol/main.go`), not from a failure. Verified first-party this iteration from the job log and the run's attempt history.

**Routing evidence:** base=9ea7b45fbb8dd17e522e238ff2436d098a90981d.
- Controller: `claude-opus-5-5` (tok: not reported).
- Designer: not spawned (A4 r2 covers sprint 2).
- Planner: `opus` via Agent (536,608 tok), resolver `agent-tool opus fail-closed:env-pin`; derive-planner-lane `opus fail-closed:env-pin`.
- Executor: `claude:claude-haiku-5-5` via Agent `haiku` (46,553 + 55,361 + 52,940 tok), resolver `recipe claude:claude-haiku-5-5 declared:provider-pin`.
- Evaluator: `sonnet` via Agent (114,138 tok), `agent-tool sonnet declared:alias-pin`.
- The driver noted codex unusable at preflight (probe rc 75, quota unreadable).
- Judge independence: cross-model (opus planner, haiku executor, sonnet judge). Metered total: $0.00.

**Progress:** 1.0 clauses 1–7 MET; the goal is unmoved. Row 201 A4 now has keys:
- An operator can approve or deny a `Human.Approve` decision packet at a terminal: `ailang-worldd approve <sha256:request-ref> [--deny]` checks the subject, takes `y`, writes the verb-bound proof and posts it.
- An operator can also decide with an enrolled SSH key: `--sign-with`, signed under `world-decision@sunholo`, judged and recorded in `key_decisions`, then landed after the decision.
- A key decision still cannot authorize a publish (L1 + L2, R4e).

A passkey cannot decide yet, and there is no out-of-band `--write-message`/`--signature` path: both are sprint 3's head.

**Ruled out:**
- **Treating "must stay unedited" literally for P2:** schema 8 makes a fixture that drops only `surface_approvals` a non-credentials store. The edit is one token per site, has precedent at `38a6d903`, and the judge audited it.
- **Running the judge's mutation drills while `verify_go.sh` ran:** they were sequenced instead (charter process fix, iteration 293).
- **Waiting for merge CI before opening the PR:** PR CI ran on clean runners in parallel with the local gates, at no cost to the rig.

**Retro:**
- **The timeout line works:** all three executor runs and the judge carried "every Go Bash call passes `timeout: 600000`, foreground, no Monitor", and none ended its turn on a wait. Iteration 389 lost a whole run cycle to that failure. Keep the line in every directive.
- **A real-tool sibling paid for itself again:** P1 (the mint namespace hard-coded in `sshsign`) was invisible to the fake signer and caught by `ssh-keygen` itself. This is the charter's iteration-299 rule, now at a second instance.
- **Overlapping remote CI with local gates:** opening the PR as soon as the commits exist overlaps remote CI with the local gates and the judge.
- **Protocol shortcut, stated:** the controller did not re-read `resources/verification-protocol.md` (136 KB) in full this fire. It worked from gate 2's recall list and the gates it read (0, 1, 2, 3 route + role-spawn-routing, 3b outline, 4, 5).

**Next:**
1. If this record's own PR is still open, land it first (docs only).
2. Then the **A4 sprint-3 planner**, from the design (r2), this plan's Appendix C and the sprint-3 table. Its head is the webauthn arm and CD6–CD7 moved from sprint 2; then comes the surface (M6–M7: the `decisions` listing, the decision card, live Chrome XD1–XD4, QUICKSTART §18). Sprint 3 is ≈ 2,060 lines by the design's price, so expect its own fallback cut (the design's: the card page becomes a fourth sprint). Run planner obligation P1 (the 91-byte message measured with `crypto.subtle`) first. Carry the judge gaps JM22, the K0 approver/format operands at the store, JM01–JM04, JM27 and JM46.
3. Mark's A2 exit run (QUICKSTART §16) and A3a exit run (QUICKSTART §17) are runnable and attended. D-WORLD-114 and D-WORLD-115 stay OPEN for Mark to confirm or veto; both defaults are applied.

## 391 — 2026-11-13 — row 201 A4 sprint 3 LANDED (PR #569, merge `e6761873`): the passkey decision arm (`DecideByAssertion`, `JudgeDecisionAssertion`, the store counter, the POST's `webauthn` arm), CD6–CD9 on a request ref, a live-Chrome passkey decision and QUICKSTART §18; the whole surface moved to sprint 4; judge PASS 91 [PRODUCT]

**Pick:** the A4 sprint-3 planner. It was record 390's Next 2 and the charter's route line. Record 390 and its addendum had landed (#567 → `a50841f8`, #568 → `7004f591`), so there was no carry-over PR.
- **Re-measured at pick** (`origin/dev` `7004f591`):
  - There was no A4 sprint-3 plan, branch or worktree.
  - `credentialsSchemaVersion = 8` (`host/store/credentials.go:44`); `sshsign.SignIn` exists (`host/sshsign/sshsign.go:50`).
  - `grep -n -i webauthn host/daemon/surface_decide.go` → no hits: sprint 2 refuses the member as an unknown field.
  - `docs/QUICKSTART.md` ended its numbered sections at §17.
- **Open PRs by the fleet account:** one, #565 `fix/contracts-typecheck` (derive `Eq` for `EvidenceGrade`/`TimeoutOutcome` for ailang#2283). It has no worktree in this clone and no World iteration token, so it is **not attributed to this mission** and was left alone (Gate 2 PR-attribution rule (c)).
- **Quorum:** none owed (the A4 r2 design covers sprint 3).
- **Ledger:** valid, 102 rows. OPEN: D-WORLD-105 (default C), D-WORLD-114 (default A, applied), D-WORLD-115 (default A, applied).
- **Clause map:** 1.0 clauses 1–7 MET, none UNMET. Release rows 150/151 are not routable (unchanged).
- **Orphans:** none (no iter-391 worktree, branch or PR existed; 390's worktrees are its own landed work).

**Gate 0:**
- Stamped first. `gh` is `sunholo-voight-kampff`. Billing CLEAN. Kill switch absent (`mission-world.disabled`). Dev CI 3/3 at `7004f591`.
- Directives: 0 allowlisted on #524 (27 comments) and 0 on #449 (47 comments); both issue-scoped watermarks are absent.
- Self-notices: rc 2, `watermark channel unavailable`. This is an instrument failure, not a verdict, as in 371–390.
- Inbox: none addressed to `mission-world`. The global unread list (3,071) is other agents' traffic.
- Skill drift: none. `SKILL.md` and all 13 `resources/*.md` match the driver clone's `origin/dev`, at both the resolved symlink target and the pin worktree.
- Weekly issue sweep: not owed (#524 was opened Monday).

**Planner** (Agent `opus`, 56.6 min, 454,705 tok): plan `w-world-approvals-by-key-a4-sprint3-plan.md` with 7 compiled patches in Appendix P.
- **Size:** 1,493 lines by `git diff --numstat 7004f591 HEAD -- cmd host docs` (376 production, 1,088 test, 29 docs); `-- cmd host` = 1,464.
- **P1 (planner obligation):** Chrome 157.0.8103.51 headless, under exactly `surfaceCSP`, rebuilt both golden 91-byte decision messages and their hashes byte-equal to Go.
- **Cut (S3-D1):** the WHOLE surface moved to a fourth sprint (≈ 1,250): the `requests.json` `decisions` listing, *Needs you*, the decision card and `decide.js`, `approve --list` decisions and `DecisionSubject.Passkeys`. This goes **further than the design's priced fallback**, which kept the listing and the *Needs you* line. The planner's reasons: the time cap, and that *Needs you* should link to the card it ships with rather than to a CLI command re-pointed a sprint later. **Controller decision: accepted.** It is a sequencing change inside the loop's own A4 child design, no attended text changes, every shipped arm is reachable by an operator, and the judge scored it honestly bounded (scope 8/10).
- **Mutants:** 76 rows: 65 killed by the named test and the package, 5 by a unit test only (JM01–JM04 and `KA-COUNTER-CHANGED`, not reachable through the POST), 4 named equivalents (incl. JM27), 2 named survivors (`KA-RACE-TEXT`, `C-LINK-ERR`, both unreachable). Ledger: 80 lines, runs R1–R3.
- **Premises found false (plan §1.4):**
  - **P1:** the A2 passkey page signs whatever challenge its link carries and always reads "Approve request …", even for a deny. Routed around through the CLI; `passkey.js` unchanged.
  - **P2:** two sprint-2 tests pinned sprint 2's own cut (the refusal of `--signature`/`--write-message` on a request ref) and had to be edited.
  - **P3–P6:** the `webauthn` member as `json.RawMessage` breaks a sprint-2 struct comparison (a pointer fixes it); a counter refusal carries the store's text; RD13 (b)'s fresh assertion must be N+2 in the test because the failed write already moved the counter; identical requests share a ref, so the daemon tests need distinct requesters.
  - **P7 — a D-WORLD-109 slip, stated plainly:** the planner's first draft of the JM46 test wrote an incomplete temp service record, the CLI fell back to `127.0.0.1:7644`, and **one read-only `/v1/health` GET reached the machine's real World service**. Nothing was written. The landed test refuses to run unless the temp record resolves to the fake daemon, and the judge mutated the record to incomplete and confirmed it fails at `approve_decide_jgaps_test.go:31` before any dial.
  - **P8:** `TestSurfaceAppRoutes` gets no new routes this sprint.

**Controller round trip:** I ran `git archive 7004f591`, extracted the 7 patches from the plan with §2's awk, and applied them. `diff -rq` against the planner's final tree is empty for `cmd`, `host` and `docs`; the control (the unpatched base against the same tree) shows 22 differences. The plan was committed (`bb8d5685`) before any sub-milestone commit. The final committed tree equals the planner's tree.

**Executor** (env `recipe claude:claude-haiku-5-5`, spawned as Agent `haiku`, as 374–390), three runs:
- Run 1 (M6a–M6c): 2.6 min, 55,314 tok.
- Run 2 (M6d–M6e): 1.5 min, 50,103 tok. Its report claimed my directive named a wrong plan filename; `grep -c` for that name in the directive is 0, so the claim was false and harmless.
- Run 3 (M6f–M6g): 3.1 min, 45,708 tok, including the live-Chrome leg (5 PASS, 0 SKIP).
- All three ran clean in one pass: no resume, no hand edits, no git writes, no stop condition.
- I gave the executor a porcelain-based snapshot command (changed files only, as sprint 2's plan had) in place of the sprint-3 plan §7's `rsync` wording, because `commit_ms.sh` commits every file in the snapshot.

The controller rebuilt each sub-milestone from its snapshot with `commit_ms.sh`; every commit was gofmt/build/vet-gated and ran its own sub-milestone's test command. Cumulative numstat (`-- cmd host docs`) after each commit:

| Sub-milestone | Commit | Numstat |
|---|---|---|
| M6a | `aac4e5b0` | 124 |
| M6b | `0362239a` | 381 |
| M6c | `6962cbea` | 628 |
| M6d | `c607244c` | 963 |
| M6e | `c0d0d17d` | 1328 |
| M6f | `9b2bf00a` | 1411 |
| M6g | `e3907bbf` | 1493 |

That equals the plan at every step.

**Gates** (controller; the plan §6 `gates.sh`, every leg recording its own rc; 760 s at load ≈ 10–20):
- `verify_ail`, gofmt, build, vet and the fence: rc 0. `pkgs` rc 0 (305 s). `broker_named` rc 0. `broker_whole_reported` rc 0. `race` rc 0. `inv` rc 0: **67** named tests `--- PASS`, none missing or skipped. `inv2` rc 0.
- **`chrome` rc 1:** `TestPasskeyPageAbortsAPendingCeremony/get-late`, an A2 test whose file, page driver and `host/surfacepage/passkey` this sprint does not touch (`git diff --quiet 7004f591 HEAD --` on all three: clean). The failure: `AbortedAfter` was `240000.1` against "exactly 240000 ms" (250 s of virtual page time in 28 ms). The whole Chrome leg (7 tests, incl. the new `TestApproveDecideWebAuthnLinkInChrome`) then passed **7/7 twice** on the controller, and the judge saw it pass alone. First instance of this test going red; filed here as a load-sensitive Chrome timing flake (row-258 class), not a regression.
- **After commit, `VERIFY_GO_LEGS=plain verify_go.sh`: rc 0 on the first run.**

**Judge** (evaluator `sonnet` via Agent, 37.3 min incl. its wait for the gate sweep, 123,613 tok; own worktree; generator ≠ judge): **PASS 91/100, zero blocking** (correctness 28, tests 26, security 19, standards 10, scope 8).
- **Sequencing:** the judge read and planned while the gates ran and started no Go command until the gates' done-marker existed (charter process fix, iteration 293).
- **Confirmed:** numstat 1,493 / 1,464; changed files = plan §5 + the two plan files; no new `/v1` route, CORS header, exec site, clock read in a law or `t.Parallel`; an approve assertion posted as deny (or the reverse) is 422 with nothing decided; letting the exact-retry path skip `useApprover` is killed (`TestRecordKeyDecisionRevokedSince`); dropping the verb from the CLI-built message is killed (`TestApproveDecideWebAuthnLinkAndAssertion`).
- **Mutants:** a 67-line ledger: 22 plan rows re-run, 41 of its own (`comm -12` against the plan's ids is empty), a must-die control (W14 disabled → killed) and a must-live control (comment edit → survived).
- **Pre-existing test edits:** exactly two sprint-2 files lose the refusal rows for the flags this sprint enables; the behaviours are re-asserted in the new artifact tests.
- **Deny on the passkey page:** not blocking; named for sprint 4's card. The CLI prints the verb in capitals, the `decide-deny-<16 hex>` request text and the message hash; the unchanged A2 page reads "Approve request decide-deny-…". The signature binds verb, packet, request and expiry, and QUICKSTART §18 tells the operator to check the printed hash.
- **Non-blocking:** `hErr` placement unpinned (only matters for inputs no entry point produces); the packet the CLI posts with `--signature`/`--webauthn-assertion` is not asserted (the daemon refuses a wrong one); the two file-size read caps are unpinned.
- **One judge finding refuted by the controller:** the judge reported the W15 mint operand (`|| p.SignCount != 0`, `host/approval/passkey_assert.go:46`) as SURVIVED. Its run used `-run '^(TestDecideByAssertion|TestApproveByAssertion)'`. I applied the same edit in its worktree and ran the whole package: `go test ./host/approval/... -count=1` → **FAIL**, killed by `TestSignCountLaw` (`passkey_law_test.go:206`, "stored 5, new 0 … want ErrSignCountRegression"). File restored, worktree clean. The A2 mint law's counter operand is pinned; the judge's scope was narrower than its sentence (rule 3b).

**Gate 3b:**
- PR #569 head `e3907bbf`, run 43354034113: **3/3 success**. Mergeable CLEAN. Closing-keyword scan: 0 hits over the PR body, title and 8 commit messages (control: 8 `Co-Authored-By` lines found).
- `gh pr merge 569 --merge --match-head-commit e3907bbf…` → merge **`e6761873`**.
- Push CI on the merge: run 43358584140 attempt 1: **3/3 success** (`ailang-code verify gate`, `go host build + test gate`, `go race gate`). **LANDED**.

**Routing evidence:** base=e67618734b2d3a48ccc13b6014e0f2d5b3475ebc.
- Controller: `claude-opus-5-5` (tok: not reported).
- Designer: not spawned (A4 r2 covers sprint 3).
- Planner: `opus` via Agent (454,705 tok), resolver `agent-tool opus fail-closed:env-pin`; derive-planner-lane `opus fail-closed:env-pin`.
- Executor: `claude:claude-haiku-5-5` via Agent `haiku` (55,314 + 50,103 + 45,708 tok), resolver `recipe claude:claude-haiku-5-5 declared:provider-pin`.
- Evaluator: `sonnet` via Agent (123,613 tok), `agent-tool sonnet declared:alias-pin`.
- The driver noted codex unusable at preflight for planner and executor (probe rc 75, quota unreadable); neither role was routed to codex this fire.
- Judge independence: cross-model (opus planner, haiku executor, sonnet judge). Metered total: $0.00.

**Progress:** 1.0 clauses 1–7 MET; goal unmoved. Row 201 A4 now has the passkey:
- An operator can decide a `Human.Approve` packet with a passkey: `approve <request-ref> [--deny] --webauthn-link` prints a link to the existing passkey page and the message hash; the operator checks the hash, uses the device, saves the file, then runs `approve <request-ref> [--deny] --webauthn-assertion <file> --expiry <printed>`. The daemon judges the assertion, checks the counter inside the credentials writer's transaction, decides, and lands the `key_decisions` row.
- An operator can also decide out of band with an SSH key: `--write-message` writes the 91-byte message, and `--signature <file> --expiry <printed>` posts it.
- A key decision still cannot authorize a publish (L1 + L2).

The surface lists nothing new yet: pending decisions are not on *Needs you* until sprint 4.

**Ruled out:**
- **Holding to the design's fallback (listing + *Needs you* in sprint 3):** the planner measured 1,493 with the whole surface out and judged the half-surface incoherent; the judge agreed it is honestly bounded. Accepted, not overridden.
- **Treating the `chrome` leg red as the sprint's:** the failing test's file, its driver and the passkey page are untouched (`git diff --quiet` on all three), and the leg passed 7/7 twice on re-run.
- **Acting on the judge's W15 survivor:** refuted at package scope (above).

**Retro:**
- **A narrow `-run` turns a killed mutant into a "survivor":** the judge's W15 finding was the third instance of a judge reading its filtered run as package truth (sprint-2 plan rows were recorded both ways for the same reason). Next judge directive: "a SURVIVED verdict needs the whole-package run of that package, not only the named tests".
- **The P7 slip is the reason the D-WORLD-109 rule names "no real daemon":** a fallback address in production code is a real network target for any test that forgets to pin it. The landed test now proves the pin before running; future CLI tests that touch the service record should do the same.
- **Overlapping the judge's reading with the gate sweep** (judge spawned at once, Go runs gated on a done-marker) saved about 15 minutes with no shared-load runs.
- **Protocol shortcut, stated:** the controller did not re-read `resources/verification-protocol.md` (136 KB) or all of `gate-3-route.md` / `gate-3b-ci-green.md` in full this fire. It worked from gate 2's recall list and the gates it read (0, 1, 2, 3 route head + role-spawn-routing, 3b outline, 4 head, 5).

**Next:**
1. If this record's own PR is still open, land it first (docs only).
2. Then the **A4 sprint-4 planner**, from the design (r2) and this plan's Appendix C: the `requests.json` `decisions` listing (DMD1–DMD3), the decision card `GET /surface/decide/{ref}` + `decide.js` (PD1, PD2; `BigInt` for the expiry per P1), *Needs you*, live Chrome XD1–XD4, `approve --list` decisions (design Q2), `DecisionSubject.Passkeys`, and the card's own verb text so a deny no longer reads "Approve". About 1,250 lines.
3. Mark's A2 exit run (QUICKSTART §16) and A3a exit run (QUICKSTART §17) are runnable and attended; design §10's A4 steps 1–3 and 5 are runnable after this sprint, step 4 (the card on a phone) waits for sprint 4. D-WORLD-114 and D-WORLD-115 stay OPEN for Mark to confirm or veto; both defaults are applied.

## 392 — 2026-11-13 — row 201 A4 sprint 4 MERGED (PR #571, merge `842b9331`): the surface — `requests.json` decisions, `approve --list` decisions, *Needs you* + notification, `DecisionSubject.Passkeys`, the decision card + `decide.js`, live Chrome XD1/XD2/XD4; A4 code-complete; judge PASS 88; merge CI's verify gate red on the AILANG v0.60.4 release (row 260, D-WORLD-116) [PRODUCT]

**Pick:** the A4 sprint-4 planner, record 391's Next 2 and the charter's route line. Record 391 had landed (#570 → `799660ac`), so there was no carry-over PR.
- **Re-measured at pick** (`origin/dev` `799660ac`):
  - There was no A4 sprint-4 plan, branch or worktree.
  - `host/surfacepage/decide_model.go` `type DecisionSubject struct` had no `Passkeys` field.
  - `host/daemon/daemon.go:1138` `GET /surface/requests.json`, `:1143` `GET /surface/decisions/{ref}`, `:1144` `POST /surface/decide/{ref}`; `docs/QUICKSTART.md:1861` §18 exists.
- **Open PRs by the fleet account:** one, #565 `fix/contracts-typecheck`. It has no worktree in this clone or the ailang clone and no mission token, and the V1, fleet and triage logs on `origin/dev` do not mention it, so it is **not attributed to this mission** and was left alone (Gate 2 PR-attribution rules (c), (d)). It became relevant at Gate 3b (below).
- **Quorum:** none owed (the A4 r2 design covers sprint 4).
- **Ledger:** valid, 102 rows at pick. OPEN: D-WORLD-105 (default C), D-WORLD-114 (default A, applied), D-WORLD-115 (default A, applied).
- **Clause map:** 1.0 clauses 1–7 MET, none UNMET. Release rows 150/151 are not routable (unchanged).
- **Orphans:** none (no iter-392 worktree, branch or PR existed).

**Gate 0:**
- Stamped first. `gh` is `sunholo-voight-kampff`. Billing CLEAN. Kill switch absent (`mission-world.disabled`). Dev CI 3/3 at `799660ac`.
- Directives: 0 allowlisted on #524 (28 comments); the issue-scoped watermark is absent, so every comment was read. #449 was not re-read this fire.
- Self-notices: rc 2, `watermark channel unavailable`. This is an instrument failure, not a verdict, as in 371–391.
- Inbox: none addressed to `mission-world` (`ailang messages list --unread --inbox mission-world` → none). The global unread list (3,073) is other agents' traffic.
- Skill drift: none. `SKILL.md` and all 13 `resources/*.md` at the resolved symlink target match the driver clone's `origin/dev`. All rules were read from that resolved copy; the pin worktree's copy was not diffed separately.

**Planner** (Agent `opus`, 92.6 min, 351,500 tok): plan `w-world-approvals-by-key-a4-sprint4-plan.md` with 9 compiled patches in Appendix P.
- **Size:** 1,947 lines by `git diff --numstat 799660ac HEAD -- cmd host docs` (655 production: 286 Go, 369 page/script; 1,284 test; 8 docs); `-- cmd host` = 1,939. Nothing was cut. Appendix C's ≈ 1,250 estimate ran over mainly on the three Chrome tests and their rig (528 lines).
- **Mutants:** 126 rows: 113 killed, 7 named equivalents, 6 named survivors (two listing error arms and two `Passkeys` store-error arms that cannot be reached, the card's 240-second abort, `credentials: 'omit'`). Ledger 148 lines; its ids equal the plan table's. Sprint 3's five named survivors carried.
- **Premises found false (plan §1.4):** P1/P2 `approve --list` never contacts a daemon and refuses `--daemon`, so it reads the world store offline through `daemon.PendingDecisions`; P3 a listing row's refusal costs a `ValidatePacket` scan each; P4 the sprint-3 judge's W15 survivor was a narrow-run artifact (as 391 found); P5 `passkey.js` is SHA-256-pinned by an A3 test; P6 the notification lives outside `needsYouHTML`; P7/P8 two sprint-1 tests pin the subject's key list and the route list (one line each); P9 PD1 compares real stored payloads because the wire structs are unexported.
- **Verb text (directive scope item 7):** the card shows its own verb on its buttons, status and file names, and XD1 asserts both verbs. The CLI `--webauthn-link` path keeps the A2 page's wording as **RS-VERB**, because the design freezes `passkey.js` (§2.8, §13), an A3 test pins its bytes and mints share its link format; changing it needs a design amendment.
- **RS-RIG:** the Chrome card tests post to a rig using the daemon's real grammar, law and judges, not the daemon's own POST over a seeded world store (≈ 150 more lines). The first real card-to-daemon POST is the attended run's step 4.

**Controller round trip:** I ran `git archive 799660ac`, extracted the 9 patches with §2's awk and applied them. `diff -rq` against the planner's final tree is empty for `cmd`, `host` and `docs`; the control (the unpatched base) shows 29 differences. The plan was committed (`4edc3276`) before any sub-milestone commit. The final committed tree equals the planner's tree.

**Executor** (env `recipe claude:claude-haiku-5-5`, spawned as Agent `haiku`, as 374–391), four runs:
- Run 1 (M7a): 1.4 min, 46,037 tok. **Stopped** after M7a's acceptance passed: the plan §7 snapshot command aborts under `set -euo pipefail` because its MANIFEST line runs `find cmd host docs` and M7a has no `docs/`. The copied files were right. I committed M7a and added one correction to every executor directive (`find $(for d in cmd host docs; do [ -d "$d" ] && echo "$d"; done)`).
- Run 1b (M7b–M7c): 2.2 min, 44,979 tok. Run 2 (M7d–M7f): 1.8 min, 54,366 tok. Run 3 (M7g–M7i, live Chrome, real `ssh-keygen`): 1.6 min, 48,020 tok. All clean: no resume, no hand edits, no git writes.

The controller rebuilt each sub-milestone from its snapshot with `commit_ms.sh` (now excluding the snapshot's `STATUS`/`DELETED` files); every commit was gofmt/build/vet-gated and ran its own test command. Cumulative numstat (`-- cmd host docs`):

| Sub-milestone | Commit | Numstat |
|---|---|---|
| M7a | `f9552d05` | 106 |
| M7b | `d45e525a` | 437 |
| M7c | `c22a4aa7` | 587 |
| M7d | `1bdd5773` | 719 |
| M7e | `a20bfc73` | 810 |
| M7f | `6e0ab1d2` | 1350 |
| M7g | `145d3336` | 1635 |
| M7h | `b7e0ba5b` | 1878 |
| M7i | `cbb3783a` | 1947 |

That equals the plan at every step.

**Gates** (controller; the plan §6 script, extracted from the plan and `diff`-equal to the planner's `gates_final.sh`; 794 s at load ≈ 10–11): every leg rc 0 — `verify_ail`, gofmt, build, vet, fence, `pkgs` (293 s), `broker_named`, `race`, `inv` (69 named tests `--- PASS`, none missing or skipped), `chrome` (all ten page tests `--- PASS`), and `VERIFY_GO_LEGS=plain verify_go.sh` rc 0 on the first run (375 s).

**Judge** (evaluator `sonnet` via Agent, 67.4 min incl. its wait for the gates, 141,030 tok; own worktree; generator ≠ judge): **PASS 88/100, zero blocking** (correctness 28, tests 23, security 19, standards 9, scope 9).
- **Confirmed:** numstat 1,947 / 1,939; 34 changed paths = plan §5 + the two plan files; no edit to `passkey.js`, `approve.js` or the A2/A3a page dirs; two new `GET` routes, both behind `approvalRoute`; no `/v1`, CSP or CORS change. The card hashes the decoded bytes itself, checks `requestRef`, builds the message from the clicked verb and refuses a bad subject before any `credentials.get`; the POST re-runs the law. *Needs you* escapes every value. `approve --list` imports no HTTP or dial code.
- **Pre-existing test edits:** three files, none weakened beyond the change (the card routes added to two route lists; `"passkeys"` added to the subject's key list).
- **Mutants:** 17 plan rows re-run (all killed by behavioural tests); 44 of its own (`comm -12` against the plan ids empty), incl. 12 `app.js` and 8 `decide.js` rows; must-die killed, two must-live survived. Ledger 73 lines.
- **Residuals:** RS-VERB and RS-RIG judged honestly described and acceptable as non-blocking.
- **Non-blocking (logged here, not filed, per D-WORLD-109's rule):** (1) no test lists decisions when there are zero mint requests (`JD-GO-LIST-EARLYRET`). **Controller-reproduced:** restoring the early `return exitOK` in `approve.go` leaves `go test ./cmd/ailang-worldd/ -run '^(TestApproveList|TestApprove)'` rc 0. The code is correct; the gap is the test. (2) `decide.js`'s pre-sign re-read that finds the request no longer ready is unexercised (the POST law refuses it). (3) The displayed deny hash is unasserted; a `null` body or a `toISOString` RangeError leaves the page on "loading" with no controls (fails closed). (4) The listing's 40-head bound and sentinel slice are unpinned. (5) `approve --list` exits 0 with a stderr line on an unreadable store (S4-D3, documented). (6) Five `decide.js` operand mutants were not run against the whole package.
- **Judge ledger hygiene:** four `-WHOLE` rows read "KILLED" by `TestMarketingChromeSixShots`, which is red unmutated without `WORLD_TOOL_AILANG_BIN`; its `-RECLASS` rows carry the real verdicts.

**Gate 3b:**
- PR #571 head `cbb3783a`: **3/3 success**. Mergeable CLEAN. Closing-keyword scan: 0 hits over the PR title, body and 10 commit messages (control: 10 `Co-Authored-By` lines found).
- `gh pr merge 571 --merge --match-head-commit cbb3783a…` → merge **`842b9331`**.
- Push CI on the merge, run 43393173068:
  - Attempt 1: `go host build + test gate` success; `go race gate` failure on `TestBackoff` (`host/surfacefollow/scheduler_test.go:131`, "timed out waiting for test synchronization"; package untouched by this sprint — **instance 2**, iteration 386 was instance 1; the two `DATA RACE` warnings are the gate's own race-detector positive control); **`ailang-code verify gate` failure**.
  - Attempt 2 (`gh run rerun --failed`): race **success**; verify gate **failure** again.
- **The verify-gate red is a compiler release, not this sprint:**
  - The job's legs 1–2 install `releases/latest` (`.github/workflows/ci.yml:43`). AILANG **v0.60.4** was released at 2026-11-13T08:23:45Z; the merge's run started 08:26. The log says `UNRECOGNISED RELEASE: legs 1-2 resolved 'v0.60.4'` and fails `✗ world/contracts.ail: check.passed != true`.
  - v0.60.4 carries ailang#2283 (merged 2026-11-12T23:46Z), which type-checks `requires`/`ensures` (`TC_CONTRACT_001`).
  - Reproduced with the released v0.60.4 darwin binary (sha256 OK): `ai-check world/contracts.ail` → `TC_CONTRACT_001`, "No instance for Eq[EvidenceGrade]" at `world/types.ail:43`. This sprint changed no `.ail`, `world/`, `packages/`, `scripts/` or `.github/` file (`git diff --stat 799660ac 842b9331 -- …` empty), so the parent fails identically.
  - Per-file sweep under v0.60.4: `TC_CONTRACT_001` in `world/{types,contracts,transitions}.ail`, their `packages/world-core` copies, `packages/world-core/_smoke.ail` and `host/replay/testdata/transition_fixture.ail`. (Other files report errors only because the sweep ran outside their module roots; not findings.)
  - **PR #565** (unattributed) adds `deriving (Eq)` to `EvidenceGrade`/`TimeoutOutcome`: under v0.60.4 `contracts.ail` passes, but `verify_ail.sh` on its head fails at **`world/transitions.ail`** — `plan`'s ensures compares the anonymous record `{ algo: string, digest: string }` with `==` (line 29). So #565 is not a complete fix.
- **Disposition:** the sprint is **MERGED**, both Go gates green on the merge; it is not called LANDED while the merge's verify gate is red. It becomes LANDED when row 260 turns that gate green on `dev`. I did not touch #565 (not attributed), did not fix forward in this fire (the change moves the projected `world/core` package bytes and needs a 0.1.2 release, which is attended), and did not pin the compiler. Filed **row 260** `w-ailang-v0604-contract-eq` and **D-WORLD-116** (A: the loop builds its own complete fix-forward and leaves #565 alone; B: land #565 first; C: pin CI to v0.60.3; default A from iteration 393).

**Routing evidence:** base=799660ac4521919e2efe9c864fa8f5a0d7666ac8.
- Controller: `claude-opus-5-5` (tok: not reported).
- Designer: not spawned (A4 r2 covers sprint 4).
- Planner: `opus` via Agent (351,500 tok), resolver `agent-tool opus fail-closed:env-pin`.
- Executor: `claude:claude-haiku-5-5` via Agent `haiku` (46,037 + 44,979 + 54,366 + 48,020 tok), resolver `recipe claude:claude-haiku-5-5 declared:provider-pin`.
- Evaluator: `sonnet` via Agent (141,030 tok), resolver `agent-tool sonnet declared:alias-pin`.
- The driver noted codex unusable at preflight for planner and executor (probe rc 75, quota unreadable); neither role was routed to codex.
- Judge independence: cross-model (opus planner, haiku executor, sonnet judge). Metered total: $0.00.

**Progress:** 1.0 clauses 1–7 MET; goal unmoved. Row 201 A4 is code-complete. After this sprint an operator can open *Needs you*, see each pending agent decision, open its card, read every field of both signed objects, and approve or deny with a passkey or by pasting an SSH signature; `approve --list` shows the same pending decisions. What is left on row 201 is attended: Mark's exit runs (A1b §15, A2 §16, A3a §17, A4 §10 incl. the card on a phone), and A3b waits on row 188 N3.

**Ruled out:**
- **Calling the verify-gate red this sprint's:** no `.ail` changed; the red reproduces on identical `.ail` bytes with the released v0.60.4, and the PR head was green on v0.60.3 three minutes before the release.
- **Calling #565 the fix:** measured; it leaves `world/transitions.ail` red under v0.60.4.
- **Merging, rebasing or commenting on #565:** not attributed to this mission (Gate 2 rule (d)).
- **A quick fix-forward in this fire:** it changes the projected package bytes and needs the attended 0.1.2 release decision and an `interfaceHash` measurement; a judged sprint is the right size, and the decision is filed.
- **Re-running the verify gate again:** deterministic on the same compiler and bytes.

**Retro:**
- **A plan's snapshot command was never run on a sub-milestone without `docs/`.** The planner tested acceptance per sub-milestone but not the §7 snapshot line, which aborts under `pipefail` when a `find` root is missing. Next planner directive: "run §7's snapshot command yourself on the first sub-milestone and on one with no docs change".
- **A head green is not a merge green when a job tracks `releases/latest`.** The PR was 3/3 on v0.60.3; the merge ran three minutes after v0.60.4 shipped. Gate 3b's merge poll caught it. Nothing to change in the loop; row 260 is the fix.
- **Judge directive:** ask the judge to run each package once unmutated and list baseline-red tests before reading any whole-package mutant run (its `TestMarketingChromeSixShots` rows needed a reclassification pass).
- **Protocol shortcut, stated:** the controller did not re-read `resources/verification-protocol.md` (136 KB) or all of `gate-3-route.md` / `gate-3b-ci-green.md` this fire. It worked from gate 2's recall list and the gates it read (0, 1, 2, 3 head + role-spawn-routing, 3b head, 4 head and rotation rules, 5, ci-health).

**Next:**
1. **This record's PR stays OPEN, unmerged** while `dev`'s verify gate is red (CLAUDE.md: nothing lands red; precedent 382). The next fire reads the charter, log and ledger from its branch `mission/world-iter392-record` (Gate 1's blocked-record rule), stacks row 260's work on that branch, and lands both once the verify gate is green.
2. Then **row 260** (dev's red outranks the queue): D-WORLD-116 default A unless Mark has answered — a design/plan for the complete contract fix-forward (two enum derivings, a named record type with `Eq` for `plan`, the package regeneration, identity goldens, `interfaceHash`/`interfaceHashV2` before and after), green under v0.60.4 and the pinned v0.41.0, stopping at a ready-to-release `world/core` 0.1.2. Do not touch #565.
3. Then row 243 (address-bar verbs): a new design r0 from the rotation designer, then quorum.
4. Mark's attended exit runs (QUICKSTART §15–§18, A4 design §10) are all runnable; D-WORLD-114, D-WORLD-115 and D-WORLD-116 are OPEN with defaults.

## 393 — 2026-11-16 — row 260 LANDED (PR #573, merge `10e6b5dc`): the kernel's contracts type-check under AILANG v0.60.4 — four `TC_CONTRACT_001` sites fixed (two enum `deriving (Eq)`, credited to #565; `plan`/`verify` restated field-wise), a `world/core` 0.1.2 candidate (not published), v0.60.4 recorded as observed; dev's verify gate green again; judge PASS 94 [PRODUCT]

**Pick:** row 260 — dev's `ailang-code verify gate` was red (a red on the owning repo outranks the queue), D-WORLD-116 default **A** applied from this iteration. The record of 392 sat on PR #572 (blocked while dev was red); per Gate 1's blocked-record rule this iteration read the charter, log and ledger from `mission/world-iter392-record` and stacked the sprint on it, so #573's merge also landed #572.
- **Re-measured at pick:** released v0.60.4 darwin binary (sha256 OK) at `/tmp/ailang-v0604-i393/ailang`; `verify_ail.sh` at `d0c9c4ce` → rc 1, `✗ world/contracts.ail: check.passed != true` (log banked as `world-iter393/baseline-verify_ail-v0604.log`). v0.60.4 still `releases/latest`.
- **Orphans — three dead attempts at this iteration, credited here:** the 2026-11-13 10:31 fire (CRASHED at gate-3, rc 1, "Failed to authenticate. API Error: 403 … Cloudflare"), the 14:31 fire (CRASHED at gate-3, rc 1, "You've hit your weekly limit · resets Nov 16 at 6am"), and the 2026-11-16 02:31 pi fire (PAUSED-NO-CAPACITY at gate-3, controller rung demoted). They left two empty worktrees (`.wt-world-iter393-row260`, `.claude/worktrees/w393-row260`) and scratch under `/tmp`; no commits, PRs or edits. This fire reused the empty `.wt-world-iter393-row260`. Between 2026-11-13 18:32 and 2026-11-15 22:32 every fire (14) paused for no capacity (Anthropic weekly limit, codex quota unreadable, pi demoted) — a three-day gap, not a loop fault.
- **Open PRs by the fleet account:** #572 (this mission's 392 record — adopted as above) and #565 (`fix/contracts-typecheck`; not attributed; left untouched).
- **Quorum:** a new design doc, so the pick-time quorum ran (below).
- **Ledger:** valid, 103 rows at pick. OPEN: D-WORLD-105 (default C), D-WORLD-114 (A, applied), D-WORLD-115 (A, applied), D-WORLD-116 (A, applied this iteration).
- **Clause map:** 1.0 clauses 1–7 MET, none UNMET. Release rows 150/151 not routable. Row 260 moves no clause (it restores clause 1's gate under the current compiler).

**Gate 0:**
- Stamped first. `gh` is `sunholo-voight-kampff`. Billing CLEAN. Kill switch absent (`mission-world.disabled`).
- Directives: 0 allowlisted on #524 (45 comments; issue-scoped watermark absent, so all were read) and 0 on #449.
- Self-notices: rc 2, `watermark channel unavailable` (no readable watermark) — instrument failure, not a verdict, as in 371–392.
- Inbox: none addressed to `mission-world`.
- Skill drift: none — `SKILL.md` + 13 `resources/*.md` (14 files) match the driver clone's `origin/dev` at both the resolved symlink target and the pin worktree.

**Designer** (env `recipe claude:claude-opus-5-5 declared:provider-pin`, spawned as Agent `opus` per the operator's standing request; the rotation file's last-used value is `claude:claude-opus-5-5` and its next entries are codex (quota unreadable) and pi glm (demoted today), so the resolver's answer was followed, as at 314): doc `design_docs/planned/w-ailang-v0604-contract-eq.md`.
- r0 (22.4 min, 178,658 tok): **four** sites, not three — v0.60.4 reports only the first failing clause per module, so `verify`'s `result.proposalHash == p.proposalHash` (transitions:48) was hidden behind `plan`. 8 of 32 gated/packaged `.ail` files fail, all from the four sites; with them fixed every leg runs to completion on v0.60.4. Chosen: enum `deriving (Eq)`; `HashRef` equalities restated field-wise on one line (a multi-line `ensures` breaks `test_kernel_contracts.py`'s line-based control, V8). Alternatives B (`HashRef deriving (Eq)`) and C (a nominal type) measured and rejected. v0.30.0 found unused by any CI lane (the "replay goldens are v0.30.0-scoped" comment is stale).
- r1 (3.4 min, 216,667 tok cumulative transcript): V15 (field equality ⇔ derived record equality, both directions, both CI compilers; A and B kill the same 7 body mutants), V16 (contentHash invariant under CHANGELOG/manifest text; a module-byte control moves it), V17 (premises banked), V18 (the two new `Eq` instances are real and invisible to both interface hashes → kind `feature`).

**Quorum** (`--author claude:claude-opus-5-5`; seats gemini-4-argon, oc-glm-5-3, oc-kimi-k3; reserve gpt6-1-sol; no absent seat either round):
- **r1 BLOCKED** ($0.266): gemini PASS; glm REJECT (equivalence and hash invariance asserted, premises uncited); kimi REJECT (release-kind inconsistency, premises uncited). → the one designer revision.
- **r2 BLOCKED** ($0.253): gemini PASS; glm REJECT (V17a's v0.54.15 install undisposed, AC8 "both jobs" unenumerated, `sameRef` reuse, V14 not banked); kimi REJECT (M2's hand-appended JSON never measured; S2 line 173 vs 174).
- Every r2 objection carried a reviewer-authored fix and none disputed the direction → **narrow-refinement carve-out** (r2 = the controller applying the reviewers' fixes, each measured first-party, rule 3f): **V19** v0.54.15 is the floor grader's `$TOOL`, never run on `world/`, and verifies 16/16 at base and candidate anyway; **V20** the exact hand-appended M2 bytes in a fresh worktree → `ok host/pins`, verifygate subset `ok`, notice absent; **V21** `sameRef(…)` in `plan`'s `ensures` is verified on v0.60.4 but **`skipped`** on v0.41.0 ("SMT declaration closure could not resolve a required sort") — so it cannot replace the field-wise form. No re-quorum. Quorum spend $0.519 metered. Artifacts in `world-iter393/quorum/` and `~/.ailang/state/mission-quorum-world/`.

**Planner** (Agent `opus`, `opus fail-closed:env-pin`, 32.0 min, 142,919 tok): plan `w-ailang-v0604-contract-eq-sprint-plan.md` (230 lines), compiled `M1.patch`/`M2.patch`, `gates.sh`.
- M1 14 files +74/−42; M2 4 files +19/−1. Both dry-run fixtures recaptured from real v0.41.0 and v0.54.11 runs (v0.41.0 twice, byte-identical). contentHash `sha256:359a82fd…4a72c` = the design's value; interface hashes unchanged; tarball 11030 bytes.
- **Premises found false:** design §3 item 14's "commit `563c6c9`" field does not exist in observed-only rows (the verbatim snippet was shipped); `scripts/marketingcheck` also flakes in the full Go gate; M1 larger than the design's +55/−30.
- Mutants MUT-13..20 added (all killed); **MUT-21 survives**: a *consistent* revert to version 0.1.1 is caught only by `gates.sh`'s AC6 assert — no repo test refuses the already-published 0.1.1.

**Executor** (env `recipe claude:claude-haiku-5-5`, spawned as Agent `haiku`, as 374–392), two runs, no hand edits, no git writes:
- M1 (13.8 min, 32,910 tok): checksum OK, numstat = plan, `gates.sh m1` all 11 legs rc 0 (799 s). Controller: `git diff | git patch-id` = M1.patch's → committed `8adbca8b`.
- M2 (33.4 min, 42,586 tok): numstat = plan; its first gate run was cut by the background-task limit (no verdict); the second: every leg rc 0 except `verify_go_plain` on `cmd/ailang-worldd` `TestCLIRealSubprocessEpisode` ("daemon announcement timed out", a 5 s bound, under load avg ~9). The executor stopped as told. Controller: the test passed `-count=5` in isolation on the M2 tree, and a full `VERIFY_GO_LEGS=plain verify_go.sh` re-run was rc 0 (62 `ok`, `✓ go gate PASSED (plain)`). patch-id = M2.patch's → committed `9e64b719`.

**Judge** (evaluator `sonnet` via Agent, `agent-tool sonnet declared:alias-pin`, 17.3 min, 77,794 tok; own worktree at `9e64b719`; generator ≠ judge): **PASS 94/100, zero blocking** (correctness 38/40, tests 18/20, security 10/10, standards 14/15, scope 14/15).
- Re-derived: `verify_ail.sh` rc 0 on v0.60.4 and on v0.41.0 (16 / 40 / 503 / 9/9); `test_kernel_contracts.py` rc 0; 8 Go packages ok. `build_world_package.sh` reproduces the package copies byte for byte; fresh dry-runs equal both fixtures; registry `world/core` 0.1.2 → 404 (unpublished), 0.1.1 → 200. Scope = design §3 items 1–16 exactly. #565 untouched.
- 33 mutants: 27 killed, 4 equivalent/unpinned (m23 equivalent, m28 comment, m26/m27 doc prose), 2 controls survived as required. Baseline: no deterministic red; `host/verifygate` `TestMissionConfigDefaultFlowIntegration` timed out once under 8-way parallel load, passes alone.
- **Non-blocking (logged here, not filed, per D-WORLD-109's rule):** (1) `docs/SELF_MOD_PUBLISH.md`'s confirmation phrase and fresh-store name are unpinned prose — **controller-reproduced**: reverting `publish world/core@0.1.2 irreversibly` to `0.1.1` at line 233 leaves `host/runbook` and `cmd/world-publish` `ok` (fails safe: the tool refuses a mismatched phrase). (2) `test_kernel_contracts.py` misses a consistent weakening of both copies (the golden catches it). (3) the verifygate fixture's 10 s deadline flakes under parallel load.

**Gate 3b:**
- PR #573 head `20b06a2d`: **3/3 success** first try (the verify gate installed `releases/latest` = v0.60.4). Mergeable CLEAN. Closing-keyword scan of the body: 0 (title says "supersedes #565", not a closing keyword).
- `gh pr merge 573 --merge --match-head-commit 20b06a2d…` → merge **`10e6b5dc`**; #572 shows MERGED (its commit is in the branch).
- Push CI on the merge, run 43829772485: attempt 1 — `ailang-code verify gate` **success**, race success, `go host build + test gate` **failure** on `host/broker` `TestSeedRestartIgnoresASurvivingOldCopy/survivor` (`exec_seed_test.go:656`) — **row 259, instance 5**; the merge tree is byte-identical to the PR head that passed (`git diff 20b06a2d 10e6b5dc` empty). `gh run rerun --failed` → attempt 2: attempt 2 — the Go job failed again, this time at `w-race-gate-blindspot/run.sh`: `INSTRUMENT FAILURE: 3 of 4 KNOWN-BAD probes completed`, because the runner could not download `go1.26.5` (`go: download go1.26.5: golang.or…`; the instrument correctly treats a skip as a refusal). The same step ran all four probes on the PR run 43828016452 an hour earlier. `gh run rerun --failed` → attempt 3: **3/3 success**. **LANDED.**

**Routing evidence:** base=10e6b5dc5a02e1142e966c43ab4eb30e418962e8@2026-11-16T08:31:26Z (Gate 4 record; Gate 1's was `842b9331…@2026-11-16T05:31:49Z`).
- Controller: `claude-opus-5-5` (tok: not reported).
- Designer: `claude:claude-opus-5-5` via Agent `opus` (178,658 + 216,667 tok), resolver `recipe claude:claude-opus-5-5 declared:provider-pin`; quorum r1/r2 blocked → designer revision r1 → controller carve-out r2.
- Planner: `opus` via Agent (142,919 tok), `opus fail-closed:env-pin`.
- Executor: `claude:claude-haiku-5-5` via Agent `haiku` (32,910 + 42,586 tok), `recipe claude:claude-haiku-5-5 declared:provider-pin`.
- Evaluator: `sonnet` via Agent (77,794 tok), `agent-tool sonnet declared:alias-pin`.
- Driver note: codex unusable at preflight for planner and executor (probe rc 75, quota unreadable) → opus / haiku. No role routed to codex.
- Judge independence: cross-model (opus designer/planner, haiku executor, sonnet judge). Metered total: $0.519 (quorum only).

**Progress:** 1.0 clauses 1–7 MET; goal unmoved. Dev's verify gate is green again under the current compiler, which also makes row 201 A4 sprint 4 (merge `842b9331`) LANDED. A `world/core` 0.1.2 candidate is ready for Mark to publish (design §7); the published 0.1.1's modules fail `TC_CONTRACT_001` under v0.60.4 (design V1, package copies), so a consumer on the current compiler needs 0.1.2.

**Ruled out:**
- **Pinning CI to v0.60.3 (D-WORLD-116 C):** not needed; adapted.
- **Landing on #565:** measured incomplete (no `transitions.ail` fix; version 0.1.1, unpublishable).
- **`HashRef deriving (Eq)` and `sameRef` in `ensures`:** the first adds a third exported instance and churns `logepoch.ail` (V6); the second is `skipped` on the frozen v0.41.0 (V21).
- **Calling the M2 `TestCLIRealSubprocessEpisode` red or the merge's `TestSeedRestart…` red this sprint's:** the first passes 5/5 in isolation and in a full re-run; the second is row 259's known flake on a tree identical to a green head.

**Retro:**
- **A compiler that reports one failure per module hides the rest.** 392's three sites were really four; the designer found the fourth only by fixing and re-running. Measuring the full failure set by iterated fix-and-rerun is now in this doc's V-log; worth the same habit at the next compiler bump.
- **The quorum blocked twice on measurements, never on direction.** r1 and r2 objections were each answerable by a controller command (V19–V21 took minutes). Instance N of the carve-out; no rule change.
- **The executor's first M2 gate run was killed by the background-task limit and a second run overlapped the load.** Next executor directive: run `gates.sh` with `GATES_SKIP_VERIFY_GO=1` and let the controller run `verify_go.sh` itself once.
- **Row 259 reached instance 5.** It stays below the release rows under rule (d); the count is in the row.

**Gate 5 bookkeeping:** the weekly thread rotated — #574 opened with the weekly report, #524 closed with a pointer, `mission-world-gh-issue` = 574 and `-prev` = 524. The weekly external-issue sweep ran: 0 orphans of 1 open issue enumerated (the bookkeeping thread; negative control fired).

**Next:**
1. Mark: publish `world/core` 0.1.2 (D-WORLD-117; design §7 runbook) — attended.
2. Row 243 (address-bar verbs): a new design r0 from the rotation designer, then quorum.
3. Mark's attended exit runs (QUICKSTART §15–§18, A4 design §10) remain runnable; D-WORLD-114/115/116 OPEN with defaults applied.

## 394 — 2026-11-16 — row 243 design r2 DONE (`w-world-address-bar-verbs.md`): operator verbs in the World app's address bar through one Go grammar shared with the CLI; parent V1–V5, V1 specified (read-only: `host/verbs`, `GET /surface/verb.json`, `episode report` in the bar); quorum blocked twice, controller carve-out r2, independent audit VALID 84 [PRODUCT]

**Pick:** row 243 `w-world-address-bar-verbs`, record 393's Next and the charter's next loop row (D-WORLD-109 = A: after row 201, design → quorum → build; D-WORLD-107 = A). Record 393 had landed (#575 → `48c92a52`), so there was no carry-over PR.
- **Re-measured at pick** (`origin/dev` `48c92a52`): no `address-bar-verbs` doc, worktree or PR (the two `mission/*243*` branches are the attended filing #456 and iteration 243's old record). `mc_quorum_gate2 w-world-address-bar-verbs` → rc 1; control `w-world-approvals-by-key-a4` → rc 0.
- **Open PRs by the fleet account:** #565 (`fix/contracts-typecheck`, unattributed, superseded by #573 at 393) — left untouched.
- **Clause map:** 1.0 clauses 1–7 MET, none UNMET. Release rows 150/151 not routable (unchanged since 393). Row 243 moves no clause (clause 5, operator surface, MET); it serves D-WORLD-94.
- **Ledger:** valid, 104 rows. OPEN: D-WORLD-105 (default C), 114 (A, applied), 115 (A, applied), 116 (A, built at 393), 117 (publish `world/core` 0.1.2?, pending Mark). No row changed.

**Gate 0:**
- Stamped first. `gh` is `sunholo-voight-kampff`. Billing CLEAN. Kill switch absent (`mission-world.disabled`).
- Directives: 0 allowlisted on #574 (2 comments) and 0 on #524 (47 comments; rotation-week catch). Both issue-scoped watermarks are absent, so all comments were read; nothing to `--mark`.
- Self-notices: rc 2, `watermark channel unavailable` — instrument failure, not a verdict, as in 371–393.
- Inbox: a controlplane transcript-scan notice for fire `world-1794807071-53892` (the 393 fire): `verdict=LEAK`, 2 hits, key name `OPENROUTER_API_KEY`, both in session files under another project's directory (`-Users-voightkampff-dev-sunholo-data-ailang`). No value or path recorded here. Nothing else addressed to `mission-world`.
- Skill drift: none — `SKILL.md` + 13 `resources/*.md` match the driver clone's `origin/dev` at both the resolved symlink target and the pin worktree.
- Main checkout: local `dev` 48 behind `origin/dev`, clean except two untracked worktree dirs; read state from origin, worked in `.wt-world-iter394-row243`.

**Designer** (env `recipe claude:claude-opus-5-5 declared:provider-pin`, spawned as Agent `opus` per the operator's standing request, as 372–393; rotation file `claude:claude-opus-5-5`, next entry codex over ration — resolver's answer followed):
- r0 (10.1 min, 156,558 tok; `ad144327`): 540 lines, 30 V-rows. Parent with five phases; only V1 specified. Premises found false vs the row text: P1 the row's example grammar (`start world <task>`, `finish <episode>`) is not the CLI's (`episode start|launch|finish|report`; `world` is both a CLI verb and the bar's navigation keyword); P2 the bar's parser is JavaScript (`surface.js` `parseAddress`), so "one parser" must be Go reached through a route; P3 `episode finish` commits (`episode_finish.go:379`) before its prompt (`:409`), so "the prompt before anything runs" is not literally true of the CLI; P4 half of a no-terminal start exists (`episode start --request` + card + `approve --await`) but both ends are CLI processes. Controller spot-checked P2, P3 and the `parseInterleaved` call sites first-party before the quorum.
- r1 (SendMessage resume, 5.0 min, 204,294 tok reported; `15070038`): 702 lines, 35 V-rows; applied the three round-1 fixes the controller adopted and nothing else.
- r2: no designer run (controller carve-out, below). **Designer diet:** one doc, r0 + one protocol-mandated revision. Held.

**Quorum** (`--author claude:claude-opus-5-5`; seats oc-glm-5-3, oc-kimi-k3, gpt6-1-sol; reserve gemini-4-argon; all present both rounds, `absent_reviewers` = `[]`):
- **r1 BLOCKED 3/3** ($0.327): gpt6-1-sol and glm — silent fallback (a failed verb request reinterpreted as navigation: `why` becomes `world why` after a 404/network error; a 400/500 swallowed); kimi — "navigation never regresses" falsified by the doc's own table (27 CLI words shadow brane names, only 3 disclosed; `ailang-worldd world get <ref>` re-read as navigation). Controller adopted sol's stricter fix (it subsumes glm's) and kimi's, and measured glm's premises before routing: `grep -rniE "Access-Control|cors" host/daemon` → 3 comment-only lines (no CORS header set), CSP control 5; the 28 switch arms; `serve` not in `ownHelpVerbs`.
- **r2 BLOCKED 3/3** ($0.348): kimi and glm independently — bare/unknown `episode` falls through step 4 to navigation, contradicting RK1/T8/VB3; gpt6-1-sol (glm's catch too) — the page's fetch has no deadline. Objections localised onto two surfaces round 1's fixes touched; round 1's surfaces drew nothing. Every objection had a reviewer-authored fix, none disputed the direction → **narrow-refinement carve-out r2** (`c062bf69`): step 4b in kimi's text; sol's AbortController budget (`VERB_ROUTE_BUDGET_MS = 15000`, above the measured 10 s `d.readCtx`, V37) with glm's message; V36 (`cmd`/`host` unchanged base→r1, control firing), V37 (`readCtx` = `WithTimeout(…, 10 s)`, all report reads take it), V38 (`maxReportEpisodeBytes = 256`).

**Independent audit** (evaluator `sonnet` via Agent, 1.5 min, 59,806 tok; read-only; generator ≠ judge): **CARVE-OUT VALID, 84/100** (carve-out 34/40, consistency 17/25, evidence 19/20, buildability 14/15). Re-ran V5, V12, V18, V20, V31, V33–V38; all reproduce except a one-line citation (`main.go:374` → 375).
- **Blocking B1, controller-reproduced and fixed** (`99c51381`): step 3b still said a three-word `world …` "falls to step 5" while 4b seized it, and 4b had no refusal text for `world` (no step-4 row). Now: refused with today's address help text (`surface.js:11`), no `cli`; §2.4 and VB3 rows pin `world a b c`.
- Non-blocking applied: the `refused` mapping in §2.2; the citation; V36's control labelled as the r1 value; the controller's own "a second submit aborts the first" sentence **withdrawn** as beyond the reviewers' text (left to the planner, with how a superseded abort differs from the timeout message). Left to the planner: the 15 s budget as a tunable; VP2's controlled-clock mechanics.
- Audit report and both quorum artifacts banked under `design_docs/verification/world-iter394/`.

**Routing evidence:** base=48c92a521fcbf46a0f1120609de80d656326ba5b@2026-11-16T10:11:57Z (Gate 4 record; Gate 1's was the same SHA @2026-11-16T09:40:33Z).
- Controller: `claude-opus-5-5` (tok: not reported).
- Designer: `claude:claude-opus-5-5` via Agent `opus` (r0 156,558 tok; r1 resume 204,294 tok as reported), resolver `recipe claude:claude-opus-5-5 declared:provider-pin`; r2 by controller carve-out.
- Quorum: rounds 1+2 $0.675 metered (glm, kimi, sol present both rounds).
- Planner and executor: not spawned (design-only fire, 388 precedent). Resolver: planner `agent-tool opus fail-closed:env-pin`; executor env `claude:claude-haiku-5-5`.
- Evaluator: `sonnet` via Agent (59,806 tok), `agent-tool sonnet declared:alias-pin`, auditing the r2 carve-out.
- Driver note: codex unusable at preflight for planner and executor (probe rc 75, quota unreadable) → opus / haiku. No role routed to codex.
- Judge independence: cross-model (opus designer, sonnet judge) and cross-vendor quorum (Z-AI, Moonshot, OpenAI). Metered total: $0.675 (quorum only).

**Progress:** 1.0 clauses 1–7 MET; goal unmoved. Row 243, the no-terminal operator path Mark asked for (D-WORLD-94/107), now has a reviewed design and a buildable first phase (V1, one sprint ≈ 1,750 lines). V1 is read-only by construction; mutating verbs from the bar wait for V3's key-approval law and V4's host runner.

**Ruled out:**
- **Parking after round 2:** every remaining objection had a reviewer fix and none disputed direction; the carve-out applied and the audit confirmed it.
- **glm's round-1 fallback-on-404 fix:** superseded by sol's stricter no-auto-fallback rule, which answers glm's objection too; V32 shows the page is `//go:embed`ded in the daemon, so a 404 is not "an older daemon".
- **Running the planner in the same fire:** design-only fires stop at the audit (388 precedent); the plan starts from a reviewed r2 next fire.
- **A ledger ask now:** V1 amends no settled sibling. V3 will need one (a key-confirmed `episode finish` vs A4 §2.6 L1's stated reason); filed when V3's child design is written, not before.

**Retro:**
- **Round 2 objected to two surfaces round 1's fixes created** (a new classification step with no head-match, a new fetch path with no deadline). Same shape as 388's K0: a fix that adds a branch needs its own input-class trace. Next designer directive (V3+): for every new step or path, trace each input class AND each failure/timeout class to its outcome.
- **The controller's r2 added one sentence no reviewer asked for, and the audit caught it.** The carve-out applies reviewers' text only; withdrawn rather than defended.
- **Watermark files for #574 and #524 are absent**, so every fire reads every comment (harmless while there are none from Mark), and the self-notices instrument stays rc 2. Unchanged since 371; not this loop's harness to fix.

**Gate 5 bookkeeping:** rotation not due (thread #574 opened today, after the Monday 07:00 boundary). Weekly external-issue sweep already run this week at 393.

**Next:**
1. Row 243 V1 sprint planner from `w-world-address-bar-verbs.md` r2 (one sprint ≈ 1,750; decide the superseded-abort behaviour and the hang-test clock).
2. Mark: publish `world/core` 0.1.2 (D-WORLD-117) — attended.
3. Mark's attended exit runs (QUICKSTART §15–§18) remain runnable; D-WORLD-114/115/116 OPEN with defaults applied.

## 395 — 2026-11-16 — row 243 V1 sprint LANDED (PR #582, merge `fb4d8b89`): operator verbs in the World app's address bar — `host/verbs` (one grammar for the CLI and the page), read-only `GET /surface/verb.json`, `episode report` in the bar with the CLI's own bytes, `approve` navigating to row 201's cards, the CLI line for every other verb; no silent fallback; judge PASS 91 [PRODUCT]

**Pick:** row 243 V1 sprint planner, record 394's Next 1 and the charter's next loop row (D-WORLD-109 = A). Record 394 had landed (#576 → `a51684af`), so there was no carry-over PR.
- **Re-measured at pick** (`origin/dev` `db2b7423`): no V1 plan, branch or worktree. `mc_quorum_gate2 w-world-address-bar-verbs` → rc 0 (`…-2026-11-16T09-53-06Z.json`); control (a fresh absent id) → rc 1. No quorum owed at pick.
- **Dev since the design (`48c92a52`):** four attended PRs (#577–#580) touched the surface V1 edits (`daemon.go` +1 route, `surface_home_summary.go`, `home_summary.go`, `app.css`, `app_test.go`, `approve.js`, QUICKSTART). The planner was told and re-measured every premise at `db2b7423`.
- **Clause map:** 1.0 clauses 1–7 MET, none UNMET. Release rows 150/151 not routable (unchanged). Row 243 moves no clause (clause 5, operator surface, MET); it serves D-WORLD-94.
- **Ledger at pick:** valid, 104 rows; OPEN D-WORLD-105, 114, 115, 116, 117 (defaults applied).

**Gate 0:**
- Stamped first. `gh` is `sunholo-voight-kampff`. Billing CLEAN. Kill switch absent (`mission-world.disabled`).
- Directives: 0 allowlisted on #574 (3 comments). The issue-scoped watermark file is absent, so every comment was read; nothing to `--mark`.
- Self-notices: rc 2, `watermark channel unavailable` — instrument failure, not a verdict, as in 371–394.
- Inbox: nothing addressed to `mission-world` (the newest unread are fleet `harness-resolved` replies to other missions).
- Skill drift: none — `SKILL.md` + 13 `resources/*.md` at the resolved symlink target match the driver clone's `origin/dev`, per file.
- Main checkout: local `dev` = `origin/dev` at Gate 1 (`db2b7423`), clean except two untracked worktree dirs.

**Attended rulings that landed mid-iteration** (`#581`, merge `c98eb536`, 16:25 local; read at Gate 3b): **D-WORLD-105 = C** (row 233 closed for 1.0; a mux collision keeps failing closed), **D-WORLD-114 = A** and **D-WORLD-115 = A** (the A3/A4 amendments ratified; the attended exit runs passed, PR #577), **D-WORLD-116 = A** (row 260's fix-forward ratified), **D-WORLD-117 = A** (publish `world/core` 0.1.2 — Mark published it, attended). `scripts/mission_ledger_provenance.sh` → `verdict=attended` for all five (flips `299b3aa4`, `f5ccbaf2`, `0061962f`, `ca707c16`, `e34813fa`). The loop actioned nothing on them beyond this record: every default had already been applied. The ledger now has **no OPEN row**. Row 256 is marked LANDED (PR #529) in the same attended PR.

**Planner** (Agent `opus`, 76.5 min, 433,801 tok; resolver `agent-tool opus fail-closed:env-pin`): `w-world-address-bar-verbs-v1-sprint-plan.md` (3,374 lines) + `sprint_w-world-address-bar-verbs-v1.json`, 8 compiled patches in Appendix P.
- **Size:** 1,909 (`git diff --numstat db2b7423 HEAD -- cmd host docs`), 91 under the cap; M1–M3 = 1,491, under the design's 1,500 fallback threshold, so nothing moved to V2 (XV1 and the `approve` navigate arm stay).
- **Mutants:** 120 rows, 224 ledger lines; 106 killed by both the named test and the behaviour suite (no inventory/static kills); 7 non-compiling rewritten and re-run; 2 named survivors (`E-ASb` — a start/finish read failing after the deadline is unreachable because the builder fails first; `JW-GO-EMPTY` — an empty bar asks nothing and goes home); 5 no-mutation controls green.
- **Premises found false (§1.4):** P1 `episode_markdown.go` reads the printer's unexported fields (+4/−4); P2 `TestSurfaceAppPure` pins `app.js`'s export list (+2/−1, the only existing test edited); P3 the Enter handler cleared the bar before `go` ran; P4 `renderCrumbs` overwrites the placeholder; P5 no `cmd` rig produced a real report (the VR3 test seeds its own journal intents); P6 the `/v1` report handler maps only the builder's error to 503 (kept via a wrapped error type); P7 QUICKSTART §18 tests read to EOF, so §19 carries no bold label; P12 (the planner's own error) the after-fetch "superseded" check is reachable — an answer can beat the newer submit's abort — restored, with a mutant.
- **Decisions the r2 audit left to the plan (§1.3):** V1-D1 a second submit aborts the pending ask and the aborted one is dropped silently (no message, no navigation, input untouched); a timeout says "the verb route did not answer in time". V1-D2 `askVerb` takes fetch/AbortController/timers as an injected `io`; the node test fires the budget timer on demand; `VERB_ROUTE_BUDGET_MS = 15000` stays a constant, not a tunable.
- **Controller round trip:** extracted the 8 patches from the committed plan with §2's awk, applied them to `git archive db2b7423`: `diff -rq` against the planner's final tree is empty for `cmd`, `host`, `docs`; control (unpatched base vs the same tree) 26 differences. Plan committed `71c84a08` before any sub-milestone.

**Executor** (env `recipe claude:claude-haiku-5-5`, spawned as Agent `haiku`, as 374–394), three runs, all foreground with 600 s tool timeouts:
- Run 1 (M1a–M1c) 3.8 min, 47,723 tok; Run 2 (M2–M3b) 3.8 min, 46,491 tok; Run 3 (M4a–M4b) 1.6 min, 41,688 tok. Every `--dry-run` clean, every acceptance `ACCEPT_OK`, no timing exception, no stop, no hand edit, no git write.
- The controller rebuilt each sub-milestone from its snapshot with `commit_ms.sh` (gofmt/build/vet/test before each commit). Cumulative numstat after each, `-- cmd host docs`, equals the plan's to the line: M1a `0b7e1a19` 203 · M1b `64899fc2` 510 · M1c `9e24bb2b` 697 · M2 `be5fc8c9` 905 · M3a `dc88b6b1` 1,248 · M3b `45987472` 1,491 · M4a `dc4aa6ff` 1,762 · M4b `4962fbea` 1,909. Final tree byte-identical to the planner's for `cmd`, `host`, `docs`; worktree clean; M4b MANIFEST all OK.

**Controller gates** (plan §6 script on `4962fbea`, per-SHA dir `gates-ctl-4962fbea`, 951 s): status, verify_ail (16 / 40 / 503; no `.ail` changed), gofmt, build, vet, fence (39 `ok`), pkgs (`host/verbs`, `host/surfacepage`, `host/daemon` 313 s, `cmd/ailang-worldd` 199 s), race, inv (65 named PASS, none missing/SKIP), chrome (15 PASS incl. XV1), `verify_go` (pinned v0.41.0) — all rc 0, `FAILS: none`.

**Evaluator** (Agent `sonnet`, own worktree `.eval-world-iter395-row243-v1`, 55.5 min, 153,112 tok; resolver `agent-tool sonnet declared:alias-pin`; generator ≠ judge — opus planner, haiku executor, sonnet judge): **PASS 91/100**, 0 blocking (correctness 28/30, tests 25/30, security 19/20, standards 9/10, scope 10/10).
- Measured 1,759 + 150 = 1,909; 34 paths all in plan §5 (+ the two plan files); `ParseInterleaved` 10 call sites at base and head; row 201 files and the static `surface.js` untouched; the only pre-existing test edit is `app_test.go` +2/−1 (five export names appended; nothing weakened).
- Own test runs: all four touched packages, `-race` on verbs/surfacepage, `vet`, the compile fence, the live-Chrome leg (3/3) — green.
- Mutants: ledger 95 lines (`judge/mutants.tsv`); 17 plan rows re-run as the plan says; ~60 own ids not in the plan's table, incl. `app.js` branches; must-die killed, must-live survived; every kill behavioural.
- **No silent fallback confirmed** on its own reading: `verbOutcome` returns `{address:true}` only for 200 + JSON + schema `world/surface-verb/v1` + `v===1` + an own-key outcome with non-empty members; every other path says why and leaves view and input.
- **Non-blocking, logged here and handed to V2 (not filed, per D-WORLD-109):** (1) a never-verb's refusal text can lose its CLI line unchecked; (2) **an empty submit does not supersede a pending ask** — controller-reproduced from source: `go('')` returns `out()` before `askVerb`, so a late `address` answer can still navigate (to what the operator typed, within the loopback answer window); the base `go()` already had this shape for index-backed addresses (`await loadIndex()` before `openState`), V1 widens it to every line; plus `X-JS-FINALLY-NULL` (a third submit's supersession untested); (3) `MaxWords` not pinned apart from its constant; (4) a `result` answer without `cli` is accepted by the page; (5) Escape-closes-panel and Enter refocus untested; (6) the moved printer's `start:` and `uncommitted intents:` lines untested (pre-existing gap); (7) the exact 256-byte episode boundary untested; (8) refusal text is a header `say` that fades in 2.6 s (named residual RS-SAY).

**Gate 3b:**
- Dev moved during the sprint to `c98eb536` (the attended rulings PR #581, docs only). `db2b7423`'s push CI had gone red on `go host build + test gate`: `TestSeedRestartIgnoresASurvivingOldCopy/survivor` (`host/broker`) — **row 259, instance 6** (log 189,192 B, 175 `--- PASS` lines as the control); `c98eb536`'s CI is 3/3 green. Sprint diff vs the new dev commits: 0 shared files (controls 4 and 34).
- PR #582 head `4962fbea`: **3/3 success** first try (verify gate, Go host gate, race gate; per-SHA poll `ci_pr582_head4962fbea.verdict`). `mergeable` MERGEABLE/CLEAN on that head. Closing-keyword scan of title, body and the 9 commit messages: 0.
- `gh pr merge 582 --merge --match-head-commit 4962fbea…` → merge **`fb4d8b89`**.
- Push CI on the merge, run 43907035880: **3/3 success** first try (`ailang-code verify gate`, `go host build + test gate`, `go race gate`). **LANDED.**

**Routing evidence:** base=db2b7423327d5223cd1a37cf40e0794c2019026a@2026-11-16T13:40:16Z (Gate 1; worktree snap at 13:41:59Z, same SHA).
- Controller: `claude-opus-5-5` (tok: not reported).
- Designer: not spawned (doc reviewed at 394).
- Planner: `opus` via Agent (433,801 tok), resolver `agent-tool opus fail-closed:env-pin` (driver note: codex `gpt-6.1-sol` unusable at preflight, probe rc 75, quota unreadable → opus).
- Executor: `claude:claude-haiku-5-5` env pin, spawned as Agent `haiku` (47,723 + 46,491 + 41,688 tok), resolver `recipe claude:claude-haiku-5-5 declared:provider-pin` (driver note: codex unusable → haiku).
- Evaluator: `sonnet` via Agent (153,112 tok), `agent-tool sonnet declared:alias-pin`.
- Quorum: none (doc reviewed at 394). Metered total: $0.
- Judge independence: cross-model (opus plan, haiku execution, sonnet judge).

**Progress:** 1.0 clauses 1–7 MET; goal unmoved. Row 243 V1 is LANDED: Mark can type `episode report <e>` into the World app's address bar and read the CLI's own text, `approve <id|ref>` opens row 201's card, and every other verb shows its exact CLI line. Mutating verbs from the bar wait for V3 (key-approval law) and V4 (host runner); V2 (`why`) is next.

**Ruled out:**
- **Rebasing the sprint onto `c98eb536` before the PR:** the new dev commits share no file with the sprint (0 of 4 vs 34), and the PR's checks run on the merge ref.
- **Treating `db2b7423`'s red as the sprint's or as dev-red:** it is the row-259 flake on a commit before the sprint, and dev's tip is green.
- **Filing the judge's eight notes as queue rows:** D-WORLD-109 logs non-blocking notes and hands them to the next phase; the empty-submit gap is a V2 design input, not a V1 blocker (it navigates only to the operator's own line, and the base had the same race for index-backed addresses).
- **Moving XV1 / the `approve` arm to V2:** M1–M3 measured 1,491, under the design's 1,500 threshold.

**Retro:**
- **The attended rulings landed while the sprint ran.** Gate 0 read the ledger at `db2b7423`; the five rows were resolved at 15:48 local and merged at 16:25. Nothing depended on them (every default was already applied), and Gate 3b's fetch caught them. No rule change needed: Gate 0 (e) already says re-read after the fetch and rebase the record.
- **A plan with compiled patches + diff-anchored mutants keeps the execution cheap:** three haiku runs totalled 9.2 min with zero stops; the planner (76.5 min) and judge (55.5 min) carried the weight. Same shape as 389–392.
- **Watermark files for #574 and #524 remain absent**, so the self-notices instrument stays rc 2 (unchanged since 371; not this loop's harness).

**Gate 5 bookkeeping:** rotation not due (#574 opened 2026-11-16 after the Monday 07:00 local boundary). Weekly external-issue sweep already run this week at 393.

**Next:**
1. Row 243 **V2** child design (`why` in the bar), from design §0.3, with the judge's V1 notes as inputs (empty-submit supersession, the never-verb CLI line, `result` without `cli`).
2. Mark's attended no-terminal exit run for V1 (design §10) when convenient.

## 396 — 2026-11-16 — row 243 V2 design r1 DONE (`w-world-address-bar-verbs-v2.md`): `why <target>` in the World app's address bar — the walk moves to `host/why` behind a `Getter`, the CLI and the daemon call one implementation, served in-process by V1's `GET /surface/verb.json`; read-only; split V2a (the seam) / V2b (the bar); quorum r1 blocked 3/3, r2 PROCEED 3/3; independent audit PASS 84 [PRODUCT]

**Pick:** row 243 V2 child design. This was record 395's Next 1 and the charter's current route (D-WORLD-109 = A: design → quorum → build). Record 395 had landed (#583 → `200ea5d5`), so there was no carry-over PR.
- **Re-measured at pick** (`origin/dev` `200ea5d5`):
  - No V2 doc, plan, branch or worktree existed. `origin/mission/attended-address-bar-verbs` is the attended D-WORLD-107 filing, already merged.
  - No open PR from the fleet account.
  - The parent's §0.3 V2 row reads "child design r0".
- **Clause map:** 1.0 clauses 1–7 MET, none UNMET. Release rows 150/151 are not routable (unchanged). Row 243 moves no clause (clause 5, operator surface, is MET); it serves D-WORLD-94.
- **Ledger at pick:** valid, 104 rows, no OPEN row.

**Gate 0:**
- Stamped first. `gh` is `sunholo-voight-kampff`. Billing CLEAN. Kill switch absent (`mission-world.disabled`).
- Directives: 0 allowlisted on #574 (4 comments). The issue-scoped watermark file is absent, so every comment was read; nothing to `--mark`.
- Self-notices: rc 2, `watermark channel unavailable`. That is an instrument failure, not a verdict, as in 371–395.
- Inbox: nothing addressed to `mission-world`. The newest unread items are triage digests and this mission's own 395 notice.
- Skill drift: none. `SKILL.md` and all 13 `resources/*.md` (14 files) match the driver clone's `origin/dev`, both at the resolved symlink target and in the pin worktree.
- Main checkout: local `dev` was 12 behind and clean except two untracked worktree dirs. I fast-forwarded it (`--ff-only`) to `200ea5d5` = `origin/dev`.
- Dev CI at `200ea5d5`: 3/3 success (`checks=3`; control `fb4d8b89` also 3).

**Designer** (resolver `recipe claude:claude-opus-5-5 declared:provider-pin`, spawned as Agent `opus` per the operator's standing request, as 372–395). The rotation file reads `claude:claude-opus-5-5`; the next entry, codex, is over ration (`MISSION_OVER_RATION=codex`), so I followed the resolver's answer. Work happened in `.wt-world-iter396-row243-v2` on `design/world-iter396-row243-v2`.
- **r0** (13.9 min, 217,488 tok; `2562348c`): 661 lines, 38 W-rows. The designer found five of the parent row's premises false or incomplete at HEAD:
  - P1: `why.go` is 1,050 lines, but only lines 424–1050 are the walk. Lines 44–253 are the CLI's shared read layer, used by 8 other files.
  - P2: there is no reader interface on the CLI side; the daemon already has the `readStore` seam.
  - P3: the CLI's 60 s `whyBudget` differs from the daemon's 10 s `readCtx`.
  - P4: the CLI prints a deadline as "not found" (`why.go:399-400`).
  - P5: the CLI shows read failures as ✗ broken links.
  - I spot-checked P3, P4 and the judge-note-3 premise (`app.js:485` `VERB_NEEDS.result = ['cli']`) first-party before the quorum.
- **r1** (SendMessage resume, 3.0 min, 251,434 tok reported; `39bb1ff7`): 755 lines, 44 W-rows. It applied the quorum round-1 fixes I listed and nothing else.
- **Designer diet:** one doc, r0 plus one protocol-mandated revision. Held.

**Quorum** (`--author claude:claude-opus-5-5`; seats gemini-4-argon, oc-glm-5-3, oc-kimi-k3; reserve gpt6-1-sol; Claude benched; `absent_reviewers` = `[]` in both rounds):
- **r1 BLOCKED 3/3** ($0.414). Each objection had a reviewer-authored fix, and none disputed the direction. I measured each premise at HEAD before routing:
  - **gemini-4-argon:** the 1 MiB cap measured `Stdout`, a field the deferred writer never sees. Measured: `verbAnswer` has `Text`/`JSON`/`Stderr`/`Exit` and no `Stdout` (`surface_verb.go:18-32`).
  - **oc-glm-5-3:** `/v1/head` was forked into a verbatim copy of `headState`. Measured: `mcpclient.go:259`, with 9 non-test call sites. Glm also raised three gaps:
    - the activity/session middleware bypass was untraced (`activity.go:78-90`, `daemon.go:1157`);
    - the typed reason strings were unchecked;
    - `cli.go:81` was misquoted (its format has two `%s`).
  - **oc-kimi-k3:** the 7-route census was asserted, not traced, and `HealthCache` was unexplained. **Kimi's premise was partly false:** `buildChain` (`why.go:804`) reads `/v1/health` on every walk, not only for `--expect-store`. I handed the designer that measurement, so r1 took kimi's second branch: document which health members reach the answer. Kimi's first branch would have deleted a route the walk actually reads. Kimi also caught `-3` in the wrong §2.3 row.
- **r2 PROCEED 3/3** ($0.400). There were two non-blocking notes:
  - glm: the page's 15 s budget and abort carry no W-row. The audit later found it true.
  - kimi: VS2's pinned `d.*` set omits `whyGetter`.
- No carve-out was needed.

**Independent audit** (evaluator `sonnet` via Agent, own detached worktree `.eval-world-iter396-row243-v2` at `39bb1ff7`, read-only, 2.8 min, 102,419 tok; generator ≠ judge): **PASS 84/100**, 0 blocking (evidence 28/30, consistency 19/25, soundness 23/25, buildability 14/20).
- It re-ran 37 W-rows, and all reproduce; every negative row carries a firing control.
- It confirmed glm's note (a) as true: `app.js:484` sets 15000, `:512-519` abort through `{signal}`, and the page says "did not answer in time". It confirmed kimi's note (b) as a real inconsistency that fails loudly.
- It confirmed read-only, bounds and no silent fallback, found that A1–A5 are each needed against the parent's text, and judged the §9 dispositions reasonable.
- **Seven non-blocking planner inputs.** I reproduced 1, 2, 3 and 5 first-party, then recorded all seven in the doc's §13.2 as a controller record, not a design change:
  1. `sum`/`whyTarget` are used by `demo_seed.go:658,730…`, outside the V2a adapter.
  2. The moving test files depend on cmd's rig, and `elide`, so they are not "package line only".
  3. `c.get` collides with `cli.go:94`'s `client.get`.
  4. VS2's `d.*` set is incomplete.
  5. `lookup.go:99`'s "cannot run it yet" reason would ride on a `result`.
  6. The 1 MiB cap is on raw bytes, while the JSON body can be larger.
  7. Smaller citation and size notes.
- The audit report and both quorum artifacts are banked under `design_docs/verification/world-iter396/`.

**Routing evidence:** base=200ea5d5b7c2944e561feeff16a4b2b3042333d1@2026-11-16T17:39:40Z (Gate 1).
- Controller: `claude-opus-5-5` (tok: not reported).
- Designer: `claude:claude-opus-5-5` via Agent `opus` (r0 217,488 tok; r1 resume 251,434 tok as reported). Resolver `recipe claude:claude-opus-5-5 declared:provider-pin`.
- Quorum: rounds 1+2 $0.814 metered (gemini-4-argon, oc-glm-5-3, oc-kimi-k3 present in both rounds).
- Planner and executor: not spawned (design-only fire, 388/394 precedent). Driver note: codex unusable at preflight for planner and executor (probe rc 75, quota unreadable) → opus / haiku.
- Evaluator: `sonnet` via Agent (102,419 tok), `agent-tool sonnet declared:alias-pin`, auditing r1. **Spawned as required; no role failed to spawn.**
- Judge independence: cross-model (opus designer, sonnet judge) and a cross-vendor quorum (Google, Z-AI, Moonshot). Metered total: $0.814 (quorum only).

**Progress:** 1.0 clauses 1–7 MET; goal unmoved. Row 243 V2 (`why` in the bar) now has a reviewed design and two buildable sprints:
- **V2a:** move the walk behind a `Getter` in `host/why`, with the CLI repointed and no visible change.
- **V2b:** the bar runs `why` with the CLI's own bytes and refuses every failure with the CLI line.

V2b also fixes the V1 judge's empty-submit supersession gap (A4).

**Ruled out:**
- **Kimi's first fix branch** (drop `/v1/health` from the route set, delete `HealthCache`): its premise that only `--expect-store` reads health is false. `buildChain` reads it on every walk (`why.go:804`).
- **A carve-out or a third designer run:** round 2 passed 3/3 with no absentee.
- **Editing the design for the audit's findings:** all are non-blocking buildability notes that the planner resolves by measurement. They are recorded verbatim in §13.2, not resolved by the controller (394's lesson: a carve-out applies reviewers' text only).
- **Running the planner in the same fire:** design-only fires stop at the audit (388/394 precedent).
- **A ledger ask:** V2 amends only the parent (A1–A5, within D-WORLD-107 = A); no settled sibling is touched. Q1–Q4 carry loop defaults.

**Retro:**
- **394's directive held.** The designer brief demanded an input-class *and* failure/timeout-class trace for every new path, and no round-1 objection landed on §2.7's failure trace. All three landed on internal consistency between the doc's code sketch and the tree (fields, a duplicated function, a call census). The next designer directive should add: "for every type, field or function your sketch names, cite where it exists at HEAD or mark it new".
- **Measuring an objection's premise changed the fix, not just its confidence.** Kimi's reject was right that the census was untraced, but wrong about what the trace would show. Forwarding the objection unmeasured would most likely have produced the deletion branch. This is rule 3f paying for itself.
- **The judge and the quorum find different classes.** Reviewers read only the doc and blocked on doc-vs-doc consistency; the judge ran the code and found four compile-time planner traps the quorum could not see. Both are needed.
- **Watermark files for #574 and #524 remain absent.** The self-notices instrument stays at rc 2 (unchanged since 371; not this loop's harness).

**Gate 5 bookkeeping:** rotation not due (#574 opened 2026-11-16, after the Monday 07:00 local boundary). The weekly external-issue sweep already ran this week at 393.

**Next:**
1. Row 243 **V2a** sprint planner from `w-world-address-bar-verbs-v2.md` r1, starting with §13.2's audit inputs 1–3 (the V2a adapter set, splitting the test files, `c.get`). Decide Q1 (one sprint or two) on the measured patch size.
2. Mark's attended no-terminal exit run for V1 (parent design §10), when convenient.

## 397 — 2026-11-17 — row 243 V2a sprint LANDED (PR #585, merge `2d530cbe`): the `why` walk moves to `host/why` behind a `Getter` (`Run`, `ParseTarget`, `Head`, `Elide`, `Sum`); the CLI keeps a thin front and a compatibility layer; one `/v1/head` reader; behaviour-neutral for the CLI; judge PASS 90 [PRODUCT]

**Pick:** row 243 V2a sprint planner. This was record 396's Next and the queue head (D-WORLD-109 = A). Record 396 had landed (#584 → `18d30a17`), so there was no carry-over PR.
- **Re-measured at pick** (`origin/dev` `18d30a17`):
  - No V2a plan, branch or worktree existed. `mc_quorum_gate2 w-world-address-bar-verbs-v2` → rc 0 (r2 artifact `…T18-02-04Z.json`: `proceed`, `absent_reviewers` `[]`), so no quorum was owed at pick.
  - `git diff --stat 200ea5d5 18d30a17 -- cmd host docs` was empty, so the code the design measured was unchanged.
  - No open PR from the fleet account. The last four slot verdicts were all `COMPLETED rc=0`, so there was no orphaned fire.
- **Clause map:** 1.0 clauses 1–7 MET, none UNMET. Release rows 150/151 are not routable (unchanged). Row 243 moves no clause (clause 5, operator surface, is MET); it serves D-WORLD-94.
- **Ledger at pick:** valid, 104 rows, no OPEN row.

**Gate 0:**
- Stamped first. `gh` is `sunholo-voight-kampff`. Billing CLEAN. Kill switch absent (`mission-world.disabled`).
- Directives: 0 allowlisted on #574 (5 comments) and on #524 (47 comments). Both issue-scoped watermark files are absent, so every comment was read; nothing to `--mark`.
- Self-notices: rc 2, `watermark channel unavailable`. That is an instrument failure, not a verdict, as in 371–396.
- Inbox: nothing new addressed to `mission-world` (one 2026-09-14 coordinator item and one 2026-10-17 CLI note, both old).
- Skill drift: none. `SKILL.md` and all 13 `resources/*.md` match the driver clone's `origin/dev`, both at the resolved symlink target and in the pin worktree.
- Main checkout: local `dev` = `origin/dev` = `18d30a17`, clean except two untracked worktree dirs.
- Dev CI at `18d30a17`: 3/3 success (`checks=3`).

**Planner** (Agent `opus`, 88.1 min, 397,975 tok; resolver `agent-tool opus fail-closed:env-pin`; `derive-planner-lane.sh` → `opus fail-closed:env-pin`): `w-world-address-bar-verbs-v2a-sprint-plan.md` (3,888 lines) + `sprint_w-world-address-bar-verbs-v2a.json`, 3 compiled patches in Appendix P.
- **Size:** 1,299 (`git diff -M --numstat 18d30a17 HEAD -- cmd host docs`), 1 under the directive's 1,300 cap. The design estimated ≈ 970; the move measures `233 285` against the rehearsal's `56 209` (§1.4 P7).
- **The cut (V2a-D2):** design §2.2 (`WhyFlags` and friends, the CLI's use of them, VA1) moves whole to V2b. A compiled probe measured it at **111** lines, so V2a + §2.2 ≈ 1,410. The design's literal cut (keep M3's tests, drop the repointing) would have landed `WhyFlags` with no caller. V2b becomes ≈ 1,041.
- **Sub-milestones:** M0 (VC1 characterisation pin: 11 CLI paths no test drove at the base, golden generated at the base) · M12 (design M1 + M2 as one unit, because neither half compiles alone) · M3t (the seam's tests).
- **Mutants:** 72 rows, 86 ledger lines, all as expected. 63 killed by both the named test and the behavioural suite (static tests excluded); 3 killed only by VW4's constant pin (`W-LOGPAGE`, `W-MAXRESP`, `W-SCANMAX`, each named with its V2b behavioural kill); 2 named equivalents (`C-HEALTHCACHE`, `C-READER-CTX`); 1 non-compiling form rewritten; 2 controls.
- **Premises found false (§1.4):**
  - P1: with type aliases, the positional `entryWire`/`worldWire` literals compile but `go vet` rejects them, so they became defined types.
  - P6: macOS `patch` leaves a deleted file as 0 bytes without `-E`.
  - P2: audit input 2 missed `e2e-tampered`, which also uses cmd's rig.
  - P3: 13 delegations, not 15.
  - P4: `elideAt` is removed, not aliased.
  - P8: VW2/VW3 need cmd's walk rig, so they live in cmd.
- **Audit inputs:**
  - 1: `sum` delegates to `why.Sum` and `whyTarget` converts inside `locate`, so `demo_seed.go` is unedited.
  - 2: `why_input_test.go` moves with 3 lines changed; `why_protected_test.go` splits (18 pure subtests and `TestWhyProtectedJSON` move, both e2e subtests stay).
  - 3: a `clientGetter(c)` closure, so `cli.go` is unedited.
  - 7: 8 `headState` call lines; VS2a pins direct imports only.
  - 4–6 carry to V2b (Appendix C).
- **Controller round trip:** I extracted the 3 patches from the plan with §2's awk and applied them with `patch -p1 -E` to `git archive 18d30a17` in a scratch repo. Numstat after each was 77 / 1,039 / 1,299 with both rename lines. `diff -rq` against the planner's `work2` was empty for `cmd`, `host`, `docs`; the control (the unpatched sprint worktree vs `work2`) differed. Plan committed `f53b9273` before any sub-milestone.

**Executor** (env `recipe claude:claude-haiku-5-5`, spawned as Agent `haiku`, as 374–396), one run, 6.1 min, 73,064 tok, foreground:
- M0, M12, M3t: every `--dry-run` clean, every acceptance `ACCEPT_OK` (M12 and M3t split into foreground calls). No timing exception, no stop, no hand edit, no git write. Snapshot counts matched §7 (STATUS 2/10/12, DELETED the two moved paths, MANIFEST 2/8/10).
- It flagged one judgement call: 31 `--- SKIP` lines in a diagnostic `-v` run of the whole cmd package. I checked them: all are pre-existing environment gates (`WORLD_CHROME`, `WORLD_TOOL_AILANG_BIN`, real-stack, PTY) in untouched files, and the acceptance's own `p` arm ends `ok`. CI does not set `WORLD_TOOL_AILANG_BIN` either, so the real-stack why chains went into the judge's directive (below).
- The controller committed each snapshot with `commit_ms.sh` (extended this iteration to `git rm` the snapshot's `DELETED` paths; gofmt/build/vet/test before each commit). Cumulative numstat: M0 `7b4586d6` 77 · M12 `347a39f5` 1,039 (`233 285 {cmd/ailang-worldd => host/why}/why.go`, `3 3 …/why_input_test.go`) · M3t `5d71dc69` 1,299. Final tree byte-identical to the planner's; worktree clean.

**Controller gates** (plan §6 script on `5d71dc69`, per-SHA dir `gates-ctl-5d71dc69`, 705 s): status, verify_ail, gofmt, build, vet (`./...`), fence, pkgs (`host/why`, `host/verbs`, `cmd/ailang-worldd`), race, inv (53 named PASS, none missing/SKIP), `verify_go` (pinned v0.41.0) — all rc 0, `FAILS: none`.

**Evaluator** (Agent `sonnet`, own worktree `.eval-world-iter397-row243-v2a`, 52.6 min, 94,072 tok; resolver `agent-tool sonnet declared:alias-pin`; generator ≠ judge — opus plan, haiku execution, sonnet judge): **PASS 90/100**, 0 blocking (correctness 28/30, tests 25/30, security 20/20, standards 9/10, scope 8/10).
- Read the rename diff as a mechanical move: exports, `net/http`/`os`/`flag`/`time` gone, `r.c.do(GET)` → the injected `r.Get`; no semantic change in `GetJSON`, `Head`, `ScanEntries`, `Locate`, `BuildChain`, `Elide`. Stream split and exits 0/1/3 hold on every path it mutated.
- **VC1's golden is a real base capture:** it copied the test and golden into a throwaway base worktree and it passed there.
- **The real-stack why chains RAN (not SKIP) and passed** with `WORLD_TOOL_AILANG_BIN` = the pinned v0.54.15: `TestDemoSeedWhyNamedLinks/real-stack`, `TestCLIAgainstSeToolsDaemon` (CI never runs these).
- Pre-existing tests: `go test -v -run TestWhy` gives 33 names at base and 40 at head; `comm -23 base head` is empty, so no test or subtest was dropped.
- Own runs: `host/why` (+ `-race`), the whole `cmd/ailang-worldd` (117.9 s), `vet`, the compile fence — all rc 0.
- Mutants: 28 of its own (ledger `judge-mutants.tsv`, 28 lines; `comm` vs the plan's ids: no overlap) incl. a must-die and a must-live control.
- **Non-blocking, logged here and handed to V2b (not filed, per D-WORLD-109):** (1) six base lines no test pins, unchanged by the move (`--scan 0`, `why head` with no head, the `rest:` receipt state operand, the fallback-plan operand, the head-absent error, a mid-rune `Elide` cut); (2) the `--json` trailing newline and the `--result` error text rest on VC1's golden alone, and `TestWhyRunReplaysTheCLI` is partly circular on them; (3) V2a sits at 1,299 of 1,300, so V2b has no headroom for V2a follow-ups; (4) VW5 cannot tell a closure that still calls `why.Head` from the real delegation (behaviourally equal).
- **Process lapse, self-reported:** it started one mutant batch under `nohup` (forbidden), killed it, restored its worktree, and re-ran those mutants in the foreground; nothing from the aborted batch is in the ledger. I checked: its worktree is clean, the throwaway base worktree is gone, and no mutant process is left. It started after the gates' `ALL_DONE` marker, so the gate legs were unaffected.

**Gate 3b:**
- Dev did not move during the sprint (`origin/dev` `18d30a17` at push).
- PR #585 head `5d71dc69`: **3/3 success** first try (verify gate, Go host gate, race gate; per-SHA poll `ci_pr585_head5d71dc69.verdict`). `mergeable` MERGEABLE/CLEAN on that head. Closing-keyword scan of title, body and the 4 commit messages: 0 (control 1).
- `gh pr merge 585 --merge --match-head-commit 5d71dc69…` → merge **`2d530cbe`**.
- Push CI on the merge, run 43988503883: **3/3 success** first try (`ailang-code verify gate`, `go host build + test gate`, `go race gate`; per-SHA poll `ci_merge2d530cbe.verdict`). **LANDED.**

**Routing evidence:** base=18d30a17a058e9adad71c42b6078b538374378d7@2026-11-16T21:40:09Z (Gate 1).
- Controller: `claude-opus-5-5` (tok: not reported).
- Designer: not spawned (doc reviewed at 396).
- Planner: `opus` via Agent (397,975 tok), resolver `agent-tool opus fail-closed:env-pin` (driver note: codex `gpt-6.1-sol` unusable at preflight, probe rc 75, quota unreadable → opus).
- Executor: `claude:claude-haiku-5-5` env pin, spawned as Agent `haiku` (73,064 tok), resolver `recipe claude:claude-haiku-5-5 declared:provider-pin` (driver note: codex unusable → haiku).
- Evaluator: `sonnet` via Agent (94,072 tok), `agent-tool sonnet declared:alias-pin`. **Spawned as required; no role failed to spawn.**
- Quorum: none (doc reviewed at 396). Metered total: $0.
- Judge independence: cross-model (opus plan, haiku execution, sonnet judge).

**Progress:** 1.0 clauses 1–7 MET; goal unmoved. Row 243 V2a is LANDED: the `why` walk is one function, `why.Run`, over an injected reader, which V2b's daemon will call in-process. `ailang-worldd why` prints the same bytes with the same exits; the World app's bar still answers `why` with its CLI line.

**Ruled out:**
- **The design's literal §8 cut** (keep M3's tests, drop the CLI repointing): it would have landed `WhyFlags` with no caller. §2.2 moved whole instead (measured 111).
- **Folding V2b into this sprint** (design Q1): V2a alone measured 1,299 against a 1,300 cap.
- **Treating the executor's 31 `-v` SKIPs as a stop:** all are pre-existing environment gates in untouched files; the acceptance `p` arm is not verbose and ended `ok`.
- **Filing the judge's notes as queue rows:** D-WORLD-109 logs non-blocking notes and hands them to the next phase.

**Retro:**
- **The real-stack why chains are invisible to CI.** `TestDemoSeedWhyNamedLinks/real-stack` and `TestCLIAgainstSeToolsDaemon` skip without `WORLD_TOOL_AILANG_BIN`, and CI does not set it. For a behaviour-neutral refactor of the walk they are the strongest end-to-end check, so this iteration had the judge run them with the pinned v0.54.15. Instance 1; the next sprint touching `why` (V2b) should name them in its plan's gates. Not a rule change yet.
- **The snapshot contract needed a `DELETED` list** for the first move-shaped sprint. The planner specified it and the controller's `commit_ms.sh` gained a `git rm` step; both are banked under `verification/world-iter397/`.
- **Watermark files for #574 and #524 remain absent**, so the self-notices instrument stays rc 2 (unchanged since 371; not this loop's harness).

**Gate 5 bookkeeping:** rotation not due (#574 opened 2026-11-16 10:07 local, after the Monday 07:00 boundary). Weekly external-issue sweep already run this week at 393.

**Next:**
1. Row 243 **V2b** sprint planner (the bar: `AdmitWhy`, the daemon's in-process getter, the route's `why` branch and the A2 cap, the page A2–A4, XV2, QUICKSTART §19), with §2.2 (`WhyFlags`, VA1; 111 measured) folded in, audit inputs 4–6, and the V2a judge's notes as inputs. Budget ≈ 1,041.

## 398 — 2026-11-17 — row 243 V2b sprint LANDED (PR #587, merge `c251132f`): `why <target>` answered in the World app's address bar — in-process through a private mux over the seven `/v1` GET handlers, read-only, every failure a refusal with the CLI line, the bar's answer byte-equal to the CLI's; live Chrome XV2; judge PASS 91 [PRODUCT]

**Pick:** row 243 V2b sprint planner. This was record 397's Next and the queue head (D-WORLD-109 = A). Record 397 had landed (#586 → `c09e0cdb`), so there was no carry-over PR.
- **Re-measured at pick** (`origin/dev` `c09e0cdb`):
  - No V2b plan, branch or worktree existed. `mc_quorum_gate2 w-world-address-bar-verbs-v2` → rc 0 (r2 artifact `…T18-02-04Z.json`: `proceed`, `absent_reviewers` `[]`; the r1 artifact reads `blocked`, as recorded at 396), so no quorum was owed at pick.
  - `git diff --stat 2d530cbe c09e0cdb -- cmd host docs` was empty (control: the whole-repo diff showed 18 record files), so the tree V2a left was unchanged.
  - No open PR from the fleet account. The last three slot verdicts were all `COMPLETED rc=0`, so there was no orphaned fire.
- **Clause map:** 1.0 clauses 1–7 MET, none UNMET. Release rows 150/151 are not routable (unchanged). Row 243 moves no clause (clause 5, operator surface, is MET); it serves D-WORLD-94.
- **Ledger at pick:** valid, 104 rows, no OPEN row.

**Gate 0:**
- Stamped first. `gh` is `sunholo-voight-kampff`. Billing CLEAN. Kill switch absent (`mission-world.disabled`).
- Directives: 0 allowlisted on #574 (6 comments, all the loop's own). The issue-scoped watermark file is absent, so every comment was read; nothing to `--mark`.
- Self-notices: rc 2, `watermark channel unavailable` (instrument failure, not a verdict, as 371–397). I read #574's comments by hand instead: no driver crash notice.
- Inbox: nothing new addressed to `mission-world` (the 2026-09-14 coordinator item and the 2026-10-17 CLI note, both old).
- Skill drift: none. `SKILL.md` and all 13 `resources/*.md` at the resolved symlink target `cmp` equal to the driver clone's `origin/dev`.
- Main checkout: local `dev` = `origin/dev` = `c09e0cdb`, clean except untracked worktree dirs. Dev CI at `c09e0cdb`: 3/3 success (`checks=3`).

**Planner** (Agent `opus`, 60.9 min, 442,485 tok; resolver `agent-tool opus fail-closed:env-pin`; `derive-planner-lane.sh` → `opus fail-closed:env-pin`): `w-world-address-bar-verbs-v2b-sprint-plan.md` (2,646 lines) + `sprint_w-world-address-bar-verbs-v2b.json`, 5 compiled patches in Appendix P.
- **Size:** 1,255 (`git diff -M --numstat c09e0cdb HEAD -- cmd host docs`; 343 production, 908 test, 4 docs) against a 1,400 cap; no cut. Over the ≈ 1,135 calibrated budget by tests: the two new `host/daemon` test files are 509 lines.
- **Sub-milestones:** M4 (§2.2 + `AdmitWhy`) 255 · M6p (the page, A2–A4) 382 · M5g (the getter, `verbWhy`, the cap) 781 · M5r (the route's branch, precedence, golden) 1,106 · M6e (VA3, VE1/VE3, XV2, QUICKSTART §19) 1,255.
- **Mutants:** 99 table rows (ledger 110 lines): 92 killed by the named test AND the behavioural suites with static tests skipped; 2 killed only by a static pin (a flag usage string, V2a's `W-LOGPAGE`); 1 killed only by live Chrome (DOM wiring); 3 controls; 1 non-compiling form rewritten. One first-run survivor exposed a weak deadline-before-failure test, which was fixed.
- **Inputs resolved:**
  - Audit 4: the private mux is built per walk, so there is no new `Daemon` field and no `sync.Once`; the pinned set is the seven handlers plus `whyRoutes`, `whyMux`, `whyGetter`, `readCtx`, `readDeadline`, `writeInternalErrorLog`.
  - Audit 5: Lookup's `why` row no longer carries the stale reason; VR1 asserts an empty reason on every result.
  - Audit 6: the cap is on raw bytes; the encoded body is bounded at 6× the cap, tested.
  - V2a judge notes: `--scan 0` and the `rest:` state words pinned; RS-REST closed (VE1 runs a resolved `rest:` receipt through the route); the bar's `--json` pinned against VC1's base-recorded golden; the head-absent lines stay residuals; VW5 unchanged.
  - V2a survivors: `W-MAXRESP` now killed by VR5, `W-SCANMAX` by VA3; `W-LOGPAGE` still survives.
- **Premises found false (§1.4), 11:** among them, five more pre-existing test files need edits than the design listed; a pending receipt reads `indeterminate`/`not-started`, not `pending`; VR8's control is `approvalOrigins`; the walk rig only makes a ✗ through its HTTP proxy, so exit-3 parity is shown on the report seed; the Chrome driver has no Escape key, so XV2 dispatches it in-page; precedence row (a) needed a test because a cancelled walk otherwise came back as a normal result.
- **Deviations (planner's, accepted by the judge):** no separate `verb-why.golden.json` (the `why` rows live in V1's route golden, generated by the route test); the daemon tests split in two files so the getter milestone lands its own tests.
- **Controller round trip:** I extracted the 5 patches from the plan with §2's awk and applied them with `patch -p1 -E` to `git archive c09e0cdb` in a scratch repo. Numstat after each: 255 / 382 / 781 / 1,106 / 1,255. `diff -rq` against the planner's `final-tree` and `work2` was empty for `cmd`, `host`, `docs`; the control (unpatched worktree vs `final-tree`) differed. Plan committed `1a244fee` before any sub-milestone.

**Executor** (env `recipe claude:claude-haiku-5-5`, spawned as Agent `haiku`, as 374–397), two runs, foreground:
- Run 1 (M4, M6p, M5g): 2.6 min, 54,776 tok. Run 2 (M5r, M6e): 2.1 min, 70,912 tok.
- Every `--dry-run` clean, every acceptance `ACCEPT_OK` (I re-grepped each log: PASS 45/16/26/37/40, 0 `--- FAIL`, 0 `--- SKIP`, 0 `^FAIL`). Live Chrome passed first time in M6p, M5r and M6e; the real-stack leg passed in M6e. No exception re-run, no stop, no hand edit, no git write. Snapshots matched §7 (STATUS/MANIFEST 5, 10, 13, 19, 22; DELETED empty).
- The controller committed each snapshot with `commit_ms.sh` (gofmt/build/vet/test before each commit). Commits: M4 `b3458dee` 255 · M6p `66305ae7` 382 · M5g `55d7c197` 781 · M5r `d4daf620` 1,106 · M6e `b613e447` 1,255, each equal to the plan. Final tree byte-identical to the planner's; worktree clean.

**Controller gates** (plan §6 script on `b613e447`, per-SHA dir `gates-ctl-b613e447`, 829 s, in the background while the judge read): status, verify_ail, gofmt, build, vet, fence, pkgs (`host/verbs`, `host/why`, `host/surfacepage`, `host/daemon`, `cmd/ailang-worldd`), race, inv (**71** named PASS, none missing/SKIP), **realstack** (`TestDemoSeedWhyNamedLinks/real-stack`, `TestCLIAgainstSeToolsDaemon` PASS, 0 SKIP), **chrome** (23 PASS incl. `TestVerbBarWhyInChrome`, 0 SKIP), `verify_go` (pinned v0.41.0) — all rc 0, `FAILS: none`.

**Evaluator** (Agent `sonnet`, own worktree `.eval-world-iter398-row243-v2b`, 23.9 min, 97,807 tok; resolver `agent-tool sonnet declared:alias-pin`; generator ≠ judge — opus plan, haiku execution, sonnet judge): **PASS 91/100**, 0 blocking (correctness 28/30, tests 27/30, security 19/20, standards 9/10, scope 10/10; 2 withheld for the survivor groups). It waited for the gates' `ALL_DONE` before running anything.
- **Read-only:** an eighth handler or a `POST /v1/head` on the private mux is killed by `TestSurfaceVerbWhyIsReadOnly` and `TestWhyGetterTransportParity`.
- **Precedence:** mapping a cancelled walk to a normal result is killed by `TestSurfaceVerbWhyDeadline` (non-vacuous); the 6× encoded bound is tested and true; restoring Lookup's stale reason is killed by `TestLookup` and the daemon outcome tests.
- **Bar = CLI, independently:** a one-byte change to the bar's text or JSON is killed by `TestVerbRouteWhyIsTheCLIs`, which compares against the real CLI path and the base-recorded golden, not `why.Run`. VC1's golden (`testdata/`) has no diff against the base.
- Its own runs: the four touched host packages (`host/daemon` 258.6 s), the touched cmd tests, `-race` on verbs/why, `vet`, the compile fence, the real-stack leg (RAN, not SKIP) and both live-Chrome bar tests — all rc 0. Scope: 25 paths, all in plan §5 + the two plan files; no `.ail`, `t.Parallel`, exec site or context root.
- Mutants: ledger `judge/mutants.tsv`, 72 lines (58 KILLED, 10 SURVIVED, 4 INVALID); 53 ids not in the plan's table; must-live and must-die controls behaved; a 13-id plan sample reproduced.
- Pre-existing tests: none weakened, several tightened; no `Test*` name disappeared (verbs 6→10, daemon 382→392, surfacepage 56→56, cmd 425→428, why 10→10).
- **Non-blocking, logged here (not filed, per D-WORLD-109):** (1) `AdmitWhy`'s `Quote` in the refusal's CLI line is unpinned — a plain join survives the whole `host/verbs` package; **reproduced by the controller** (mutant `host/verbs/flags.go:137`, `go test ./host/verbs/` rc 0, restored); the effect is a mis-quoted reason line only. (2) three cosmetic page strings (the placeholder suffix, the keys-row `why T`, `index.html`'s placeholder) survive `host/surfacepage`. (3) `W-LOGPAGE` and the flag-usage string are killed by static pins only. (4) its daemon-side survivors were confirmed in the narrow `-run` only, not the whole 258 s package. It did not re-run the M5g commit on its own.

**Gate 3b:**
- Dev did not move during the sprint (`origin/dev` `c09e0cdb` at push and at merge).
- PR #587 head `b613e447`: **3/3 success** first try (verify gate, Go host gate, race gate; per-SHA poll `ci_pr587_headb613e447.verdict`). `mergeable` MERGEABLE/CLEAN on that head. Closing-keyword scan of the PR body and the 6 commit messages: 0 (control 1).
- `gh pr merge 587 --merge --match-head-commit b613e447…` → merge **`c251132f`**.
- Push CI on the merge, run 44010591675: **3/3 success** first try (`ailang-code verify gate`, `go host build + test gate`, `go race gate`; per-SHA poll `ci_mergec251132f.verdict`). **LANDED.**

**Routing evidence:** base=c09e0cdb3c8a3eb25c208dbacedddfb690fe8ca0@2026-11-17T01:39:53Z (Gate 1).
- Controller: `claude-opus-5-5` (tok: not reported).
- Designer: not spawned (doc reviewed at 396).
- Planner: `opus` via Agent (442,485 tok), resolver `agent-tool opus fail-closed:env-pin` (driver note: codex `gpt-6.1-sol` unusable at preflight, probe rc 75, quota unreadable → opus).
- Executor: `claude:claude-haiku-5-5` env pin, spawned as Agent `haiku` (54,776 + 70,912 tok), resolver `recipe claude:claude-haiku-5-5 declared:provider-pin` (driver note: codex unusable → haiku).
- Evaluator: `sonnet` via Agent (97,807 tok), `agent-tool sonnet declared:alias-pin`. **Spawned as required; no role failed to spawn.**
- Quorum: none (doc reviewed at 396). Metered total: $0.
- Judge independence: cross-model (opus plan, haiku execution, sonnet judge).

**Progress:** 1.0 clauses 1–7 MET; goal unmoved. Row 243 V2 is complete: an operator can type `why <target>` in the World app's address bar and read the same answer the CLI prints, served in-process and read-only; every failure says why and carries the CLI line. Row 243's next phase is V3 (the verb-approval law), which needs a child design and a ledger ask.

**Ruled out:**
- **A V2c cut:** the plan measured 1,255 against the 1,400 cap, so nothing moved out.
- **Treating the judge's `Quote` survivor as blocking:** the refusal still carries a CLI line and refuses; only an argument needing quotes would print unquoted. Reproduced and logged for V3's planner.
- **Re-running the whole `host/daemon` package for the judge's narrow-run survivors:** the controller's `pkgs` leg ran the whole package green on the same SHA; the survivors are test-strength notes, not regressions.

**Retro:**
- **The real-stack leg is now a named plan gate** (397's instance 1 → applied here): the planner put `TestDemoSeedWhyNamedLinks/real-stack` and `TestCLIAgainstSeToolsDaemon` in §6 as `realstack`, measured at the base, and both the gates and the judge saw them RUN. CI still skips them; that stays a known gap, not a rule change.
- **Fastest V-sprint shape so far:** planner 61 min → executor 2 runs, 4.7 min total → gates 829 s concurrent with the judge (24 min). The judge's sequencing wait cost nothing because the gates finished first.
- **Test-heavy plans exceed the design's size estimate** (1,255 vs ≈ 1,135): the overrun was all tests. Instance 1 for "test lines dominate the overrun"; no rule change.
- **Watermark files for #574/#524 remain absent**, so the self-notices instrument stays rc 2 (unchanged since 371; not this loop's harness).

**Gate 5 bookkeeping:** rotation not due (#574 opened 2026-11-16 after the Monday 07:00 boundary; the next rotation is Monday 2026-11-23). Weekly external-issue sweep already run this week at 393.

**Next:**
1. Row 243 **V3** child design (the verb-approval law: key-confirmed mutating verbs in the bar; it will need a ledger ask, e.g. `episode finish` vs A4 §2.6 L1). Designer = rotation (file last-used `claude:claude-opus-5-5`, so next is `codex:gpt-6.1-sol` — check the resolver; codex probed rc 75 this fire), then quorum `--author`. Inputs: the V2b judge's non-blocking notes (the `Quote` survivor, the cosmetic page strings).

## 399 — 2026-11-17 — row 243 V3 design DONE (`w-world-address-bar-verbs-v3.md` r2): `episode finish --confirm-on-surface` — an enrolled key on `/surface/confirm/<id>` replaces the PR prompt's typed `y`; commit-then-approve; six sprints, schema 9 gated behind the new ask D-WORLD-118; quorum r1/r2 blocked, controller carve-out r2, independent audit PASS 87 [PRODUCT]

**Pick:** row 243 V3 child design. This was record 398's Next and the queue head (D-WORLD-109 = A). Record 398 had landed (#588 → `994ba123`), so there was no carry-over PR.
- **Re-measured at pick** (`origin/dev` `994ba123`): no `-v3` doc, branch or worktree existed (`ls design_docs/planned | grep verbs` shows v1/v2 only). No open PR from the fleet account. The last five slot verdicts were all `COMPLETED rc=0`, so there was no orphaned fire.
- **Clause map:** 1.0 clauses 1–7 MET, none UNMET. Release rows 150/151 are not routable (unchanged). Row 243 moves no clause (clause 5, operator surface, is MET); it serves D-WORLD-94.
- **Ledger at pick:** valid, 104 rows, no OPEN row.

**Gate 0:**
- Stamped first. `gh` is `sunholo-voight-kampff`. Billing CLEAN. Kill switch absent (`mission-world.disabled`).
- Directives: 0 allowlisted on #574 (7 comments). The issue-scoped watermark file is absent, so every comment was read; nothing to `--mark`.
- Self-notices: rc 2, `watermark channel unavailable` (instrument failure, not a verdict, as 371–398).
- Inbox: `mission-world` has no unread message.
- Skill drift: none. `SKILL.md` and all 13 `resources/*.md`, at both the resolved symlink target and the pin worktree, `cmp` equal to the driver clone's `origin/dev` (28 of 28).
- Main checkout: local `dev` = `origin/dev` = `994ba123`, clean except untracked worktree dirs. Dev CI at `994ba123`: 3/3 success (`checks=3`).

**Designer** (resolver `recipe claude:claude-opus-5-5 declared:provider-pin`, spawned as Agent `opus` per the operator's standing request, as 372–398). The rotation file reads `claude:claude-opus-5-5`; the next entry, codex, is over ration (`MISSION_OVER_RATION=codex`), so I followed the resolver. Worktree `.wt-world-iter399-row243-v3`, branch `design/world-iter399-row243-v3`. The directive carried 394's and 396's retro rules (trace every input and failure class; cite every named symbol at HEAD or mark it NEW).
- **r0** (20.9 min, 345,943 tok; `674fd735`): 862 lines, 52 W-rows. Seven of the parent's V3 premises measured; four false or incomplete:
  - `episode start` does not need the mode: mints are already key-approvable (`--request`, the A3a card, `approve --await`, R1); the gap is launch presets (`episode.go:249`), which is V4's.
  - Finish has a second `y` per unsandboxed check (`episode_finish.go:270`); surface mode refuses before any check runs when the preset has one.
  - "A4's judges" are bound to A4's decision message; V3 takes A3a's stored-request shape under a new message.
  - "The exact prompt bytes" show only six body lines; the subject also binds scope, title, full body, committer and protected paths.
  - Q3 decided: **commit-then-approve.** A preview can never equal the real prompt: the commit sha depends on the clock, the body embeds it, and `--dry-run` itself records outcomes and commits (W7–W9).
  - Holds: credentials schema 8 → 9 (`credentials.go:44`; I reproduced it before the quorum, with the `key_decisions` control).
- **r1** (SendMessage resume, 8.8 min; `52b7d85a`): 985 lines, W53–W61. Answered round 1 with the controller's measurements: schema 9 and every store read moved behind D-WORLD-118; the CLI re-verifies the signature over its own bytes before it reads the decision back; the two-store write reordered (world first, then one credentials transaction).
- **Designer diet:** one doc, r0 plus one protocol-mandated revision. Held.

**Quorum** (`--author claude:claude-opus-5-5`; seated `codex:gpt-6.1-sol`, gemini-4-argon, oc-glm-5-3; reserve oc-kimi-k3; claude-sonnet-5 sits out):
- **r1 BLOCKED** ($0.4355; artifact `…v3-2026-11-17T06-03-31Z.json`). `codex:gpt-6.1-sol` ABSENT (quota) → the reserve took the seat.
  - gemini-4-argon **pass** (concern: the two-store write was not atomic).
  - oc-glm-5-3 **reject:** schema 9 landed before the answer, and the C branch had no disposition for stores already at v9. **Measured first:** the newer-store refusal already exists (`credentials.go:341-342`), so reverting the constant would strand a real store. The objection was real.
  - oc-kimi-k3 **reject:** the read-back (`GET /v1/objects/{ref}?payload=true`) and A4's writer hold were cited without W-rows. **Measured first:** route `daemon.go:1113`, payload branch `handlers.go:508-510`, `why.Reader.Object` returns raw `Payload []byte` (`why.go:66-70`, `:142`); A4 holds its writer across the world write (`surface_decide_key.go:66→90`, `surface_decide.go:211`). Premise true, rows missing.
- **r2 BLOCKED** ($0.5409; `…T06-19-11Z.json`). Sol absent (quota) again.
  - gemini-4-argon **pass** (concern: a revocation in the ≈ 1 s after C15 makes the CLI refuse a validly recorded approval; fail-safe).
  - oc-glm-5-3 **reject:** "the orphaned C14 approval is inert" was unmeasured; it feared a re-run lands in `resume-dispatch` and dispatches under it. **Measured: false.** The approve phase writes only the request and decision (`host/daemon/episode_finish.go:166-178`), `resume-dispatch` needs an episode-finish record at HEAD with a failed-dispatch outcome (`episode_finish_resume.go:19-42`), and reuse reads `record.Approval` and refuses `ExpiresAt < now+300` (`:122`).
  - oc-kimi-k3 **reject:** the read-back seam's types were measured piecewise, not joined; T11 cited C14 for C15; `decision_ref`'s CHECK does not enforce hex. **Measured:** `worldReader.c` is `*client` (`why_cli.go:181-190`), `finishClient = newClient` → `*client` (`episode_finish.go:52`, `cli.go:31`), `episode report` uses `newClient` (`:84`). The seam holds.
- **Controller carve-out r2** (`e81921c6`, `a4a9b2c1`): both remaining rejects kept the direction and carried concrete fixes. I applied them, choosing glm's "rewrite to the traced truth" branch by measurement, and added W62–W65 and §16. The audit caught one missed reviewer fix (gemini's "T2 and T11 C14→C15"; I had done T11 only), applied as `a4a9b2c1`.
- Objection surfaces by round: r1 on the store schema and an uncited read path; r2 on the r1 fixes' own new claims (the orphan, the seam). The 388/394 shape again: a fix that adds a path draws the next round's objections. Total $0.9764.

**Evaluator** (Agent `sonnet`, 1.8 min, 107,865 tok; resolver `agent-tool sonnet declared:alias-pin`; generator ≠ judge — opus designer, sonnet judge, cross-vendor quorum): **PASS 87/100**, 0 blocking (premises 29/30, carve-out 15/20, ask safety 18/20, soundness 17/20, plannability 8/10).
- Re-ran 34 W-rows (W53–W65 and 21 older ones): all reproduced. Three un-rowed §2 citations checked true.
- Carve-out hunk by hunk: every hunk traced to a reviewer's text except "pinned by CS-tests" (harmless controller words) and the missed T2 fix (now applied). The traced truth for glm holds at HEAD.
- The default while D-WORLD-118 is open is safe: sprints 1–2 register no route, read no store and change no schema.
- **Non-blocking, logged here (not filed, per D-WORLD-109), for the S1/S3 planners:** (1) §8 S3 says "the four GET registrations" but §2.5 defines three (the fourth is S4's POST). (2) §2.9 P6 says CR1–CR4; §6 names CR1–CR2, and the test list jumps CX15 → CX17. (3) A raced supersession between C14 and C15 is covered by "C15 fails" but has no window-table row (inert per W62). (4) The `resume-dispatch` reuse path after a surface-confirmed approval whose push failed skips CF6/CF7, matching the typed path's authority; CF12 should say so. (5) gemini r2's revocation-window concern.

**Gate 3b:** see the record PR below (design + record in one docs PR, as 394/396).

**Routing evidence:** base=994ba12358dc50246db2862d78f66d64e03b29f4@2026-11-17T06:27:09Z (Gate 4).
- Controller: `claude-opus-5-5` (tok: not reported).
- Designer: `claude:claude-opus-5-5` via Agent `opus` (r0 345,943 tok; r1 resume 442,286 tok as reported). Resolver `recipe claude:claude-opus-5-5 declared:provider-pin`.
- Planner and executor: not spawned (design-only fire, 388/394/396 precedent). Driver note: codex unusable at preflight for planner and executor (probe rc 75, quota unreadable) → opus / haiku.
- Evaluator: `sonnet` via Agent (107,865 tok), `agent-tool sonnet declared:alias-pin`. **Spawned as required; no role failed to spawn.**
- Quorum: gemini-4-argon, oc-glm-5-3, oc-kimi-k3 (reserve); `codex:gpt-6.1-sol` absent on quota both rounds. Metered total: $0.9764 (quorum only).
- Judge independence: cross-model (opus designer, sonnet judge) and a cross-vendor quorum (Google, Z-AI, Moonshot).

**Progress:** 1.0 clauses 1–7 MET; goal unmoved. Row 243 V3 is designed: how a finish's PR prompt can be confirmed with an enrolled key on the World app instead of a typed `y`, and the ask that decides whether it may be.

**Ruled out:**
- **Preview-bound approval (the parent's Q3 default):** measured impossible — the sha depends on the clock and `--dry-run` commits (W7–W9).
- **`episode start` in V3:** mints are already key-approvable under R1; its gap is V4's launch-preset runner.
- **glm r2's feared reuse of an orphaned approval:** measured false (W62); no residual needed.
- **A third designer run or a park:** both r2 rejects kept the direction with concrete fixes, so the carve-out applied.
- **Re-running the absent `codex:gpt-6.1-sol` seat:** its absence was quota (capacity), not budget; the reserve seat filled it, so the quorum ran at three present reviewers both rounds.

**Retro:**
- **Instance 2 of "the controller's carve-out misses a reviewer fix the audit catches"** (394's audit caught an added sentence; here a missed T2). The audit is doing its job. No rule change; next carve-out: grep the doc for every identifier each `proposed_fix` names before committing.
- **Instance 3 of "a fix that adds a path draws the next round"** (388, 394, 399). 394's directive rule (trace every new path's input and failure classes) held for r0, but r1's new paths (the CLI re-check, the reorder) brought new unmeasured claims. Next designer revision directive: every r1 claim about existing behaviour gets a W-row before the doc is returned, with the controller's own measurements as the template.
- **Pass-concerns are worth reading:** gemini's r1 concern (the non-atomic write) was real and cheap to fix in the same revision.
- **`design-quorum` now requires `--controller-note` with `--controller-verdict`** (rc 2, first call, $0). Instance 1; noted for the next fire.

**Gate 5 bookkeeping:** rotation not due (#574 opened 2026-11-16, after the Monday 07:00 boundary; next rotation Monday 2026-11-23). The weekly external-issue sweep already ran this week at 393.

**Next:**
1. Row 243 **V3 sprint 1** planner (the pure law: the verb request, the `world-verb/v1` message, the two key laws, `CheckConfirmation`/`RecheckVerb`, `FinishPromptText`; ≈ 935 raw). Sprints 1–2 need no answer; S3–S6 wait for **D-WORLD-118**. Feed the planner the audit's five notes and design §6's planner obligations.

## 400 — 2026-11-17 — row 243 V3 sprint 1 LANDED (PR #591, merge `b0817474`): the verb-confirmation law — `world/verb-request/v1`, the `world-verb/v1` message, `ConfirmBySignature`/`ConfirmByAssertion`/`SignedVerb`, `CheckConfirmation`, and the finish prompt's text moved to `host/broker` (byte-identical); judges cut to S1b; judge PASS 94 [PRODUCT]

**Pick:** row 243 V3 sprint 1, the planner step named as Next in record 399. That was also the queue head.
- **Attended ruling at Gate 0:** **D-WORLD-118 = A** (Mark, attended 2026-11-17, `72007b37`, merged as #590 → `d3b0a813`). Provenance: `mission_ledger_provenance.sh --id D-WORLD-118` → `verdict=attended`, flip `72007b37`, creation `4ef4abf8`. A key may confirm any finish, draft or ready, and V3 is built as designed, credentials schema 9 included. **V3 sprints 3–6 are now loop-routable.** The B-only branches (F6, `DraftOnly`) are not built.
- **Re-measured at pick** (`origin/dev` `d3b0a813`):
  - None of the sprint-1 symbols existed (`world-verb/v1`, `CheckConfirmation`, `RecheckVerb`, `FinishPromptText`, `verb-request/v1`: 0 hits in `host cmd internal`). Control: `AdmitWhy` was found.
  - No open fleet PR. The checked quorum reader returned rc 0 (reviewed), so no pick-time quorum was owed.
- **Clause map:** 1.0 clauses 1–7 MET, none UNMET. Release rows 150/151 are still not routable. Row 243 serves D-WORLD-94 and moves no clause.
- **Ledger at pick:** valid, 105 rows, no OPEN row.

**Gate 0:**
- Stamped first. `gh` is `sunholo-voight-kampff`. Billing CLEAN. Kill switch absent.
- Directives: 0 allowlisted on #574 (8 comments). The watermark file is absent, so every comment was read, and there was nothing to `--mark`.
- Self-notices: rc 2, `watermark channel unavailable`. This is an instrument failure, not a verdict, as in 371–399.
- Inbox: no unread message for `mission-world` (`messages list -unread -inbox mission-world`; control: the read list returns rows).
- Skill drift: none. `SKILL.md` and all 13 `resources/*.md` at the resolved symlink target `cmp` equal to the driver clone's `origin/dev`.
- Dev CI at `d3b0a813` was pending at Gate 1 and was 3/3 success by Gate 4.

**Planner** (Agent `opus`, 56.7 min, 373,966 tok; resolver `agent-tool opus fail-closed:env-pin`; `derive-planner-lane` → `opus fail-closed:env-pin`). It wrote `w-world-address-bar-verbs-v3s1-sprint-plan.md` (2,325 lines) and its JSON, committed as `a71f2296`.
- **Cut, by measurement:** S1 compiled to **1,626** lines (design ≈ 935) against the 1,400 cap, so it kept **M1 + M3 = 1,300**. The M2 judges (`JudgeVerbSSH`, `JudgeVerbAssertion`, `RecheckVerb`, VJ1) move to a new **S1b, 326 lines, compiled and green on top of M3** (`plan-scratch/patches/S1b-M2.patch`).
  - M3 stays in S1 because S2's card test (XC1) and S3's F5 need `FinishPromptText`.
  - The judges' first caller is S4.
- **Seven §1.4 premises found false or incomplete:**
  - The size.
  - §6 still lists VT2 under S1. It is leftover r0 text: §2.3 [r1] and §8 put `FinishConfirmRefusal` in S3. Not built here.
  - `displayCommitter` "3 callers": it is the definition plus 2 callers.
  - The existing m2/m3 tests prove only the question and one committer line. A base-generated golden was added.
  - The recorded Chrome fixtures cover only mint→verb; the reverse direction uses `approvaltest` keys.
  - `Validate` checks UTF-8 but not NUL, as designed; noted for S2's card.
  - A format-constant choice with identical behaviour.
- **Mutants:** 115 ledger rows (109 mutants, 3 controls, 3 rewritten non-compiling rows); 105 killed, 3 named equivalents.
- **Controller round trip:** byte-identical. I extracted the four patches from the plan, applied them to a fresh `git archive d3b0a813`, and `diff -r` against the planner's `work` and `final-tree` was silent. Control: 12 files differ from base, which matches the plan's 12.

**Executor** (env `recipe claude:claude-haiku-5-5`, spawned as Agent `haiku` as in 374–399; one foreground run, 3.1 min, 59,154 tok).
- M1a → M1b → M1c → M3: every `--dry-run` was clean and every acceptance `ACCEPT_OK`. I re-grepped each log: PASS 121/125/129/74, 0 `--- FAIL`, 0 `--- SKIP`, 0 `^FAIL`.
- Snapshots matched §7 (STATUS/MANIFEST 2/4/6/12, DELETED empty).
- I committed each snapshot with `commit_ms.sh`, which runs gofmt/build/vet/test before each commit:
  - M1a `b1573c64` 387
  - M1b `1e182629` 849
  - M1c `011f6716` 1,034
  - M3 `b4033af7` 1,300
- Each numstat equals the plan's. The final tree passes the M3 `MANIFEST.sha256` (0 mismatches) and the worktree is clean.

**Controller gates** (plan §6 `gates.sh` on `b4033af7`, 965 s, run while the judge read):
- rc 0: status, `verify_ail`, gofmt, build, vet, fence, pkgs (253 s), race, **prompt** (the 13 existing finish/preset prompt tests by name, unedited, plus VT1) and **inv** (47 named passes).
- **`verify_go` rc 1** on two untouched-area tests:
  - `host/daemon` `TestWorkspaceExecConfiguredBindsTheEpisodeHandler`: a 5 s exec timeout on `/bin/sh -c 'pwd -P'`.
  - `host/verifygate` `TestMissionConfigDefaultFlowIntegration/invalid_ledger`, which is on the known list.
  - Load average was about 16 with the judge running at the same time. Each **passed alone on the same head** (0.86 s; 18 s). The sprint touches neither package.

**Evaluator** (Agent `sonnet`, 29.8 min, 97,223 tok; `agent-tool sonnet declared:alias-pin`; its own worktree `.wt-world-iter400-row243-v3s1-eval`): **PASS 94/100**, 0 blocking.
- Scores: correctness 28/30, tests 27/30, security 19/20, standards 10/10, scope 10/10.
- **Its own mutants:** a 132-row ledger, all `J-` ids, 0 shared with the plan's.
  - The must-die control (`VerbNamespace` set to the mint namespace) was killed; the must-live comment-only control survived.
  - Static-only kills were re-run with the static test excluded.
  - Two survivors: `J-W15-OR`, an equivalent redundant operand copied from the mint law's W15; and `J-CC-DEFAULT`, an error-message format.
- **Other checks:**
  - Law order was checked against `signed.go`/`passkey_assert.go`/`decision_signed.go`. `CheckConfirmation` = law − P1 − W15, and its message is built only from the caller's inputs.
  - Domain separation holds both ways. The namespaces `world-approval@`/`world-decision@`/`world-verb@sunholo` and the prefixes are disjoint.
  - The golden was regenerated from the base's own `finishPrompt` in a throwaway worktree.
  - No pre-existing `_test.go` was modified (`--diff-filter=M` empty).
  - M1b is green on its own.
  - It judged the M2 cut sound.
- **Reproduced first-party before banking:** the must-die mutant in the judge's worktree is killed by `TestVerbDomainSeparation` and `TestCheckConfirmation` with `TestVerbMessageLayout` skipped. After restoring, the tree is clean and the tests pass.
- **Non-blocking:**
  1. Prompt-line byte-identity rests on the base-generated golden alone. The existing finish tests catch only the question, the committer and the protected lines.
  2. `VerbPrefix` and the 0x00 separator are pinned only by the layout test, which is inherent because sign and verify share one function.
  3. The two survivors above.
  4. The `verify_go` timing reds above.
  5. `Validate` does not refuse NUL and does not constrain `Committer`/`Protected`, as designed.

**Gate 3b:** PR #591 head `b4033af7`: CI 3/3 success (verify, go host build + test, race). Merged with `--match-head-commit` → **`b0817474`** (tree byte-identical to the head). **Merge push CI attempt 1: 2/3.** The race job (run 44074508663) went red on store-deadline and timeout tests in untouched code: `cmd/ailang-worldd` `TestApproveRefusesReplay` (`put mint request: begin: context deadline exceeded`) and `TestDemoProofWrongKeyCannotProve` (proof commit uncertain), `host/daemon` `TestHumanApproveIsSerializedAcrossEpisodes/together`, and an 8 m package timeout panic. Its two `DATA RACE` lines are the gate's positive control. Negative control: the same tree's race job was green on the PR head, and the last 11 dev merges' race gates were all success. `gh run rerun --failed` (the row-254 timing-class precedent, 382/398) → **attempt 2: 3/3 success** (`ci_merge_b0817474_attempt2.verdict`).

**Routing evidence:** base=b0817474f853cd7ce859974e1abff49bf210a77b@2026-11-17T11:29:54Z (Gate 4). The Gate-1 base was d3b0a8139bdecc36f8d083f4b701aef822340b27@2026-11-17T09:39:57Z, and the worktree base snap agreed with it (no drift).
- Controller: `claude-opus-5-5` (tok: not reported).
- Designer: not spawned (doc reviewed at 399).
- Planner: `opus` via Agent (373,966 tok). Driver note: codex `gpt-6.1-sol` unusable at preflight (probe rc 75, quota unreadable) → opus.
- Executor: `claude:claude-haiku-5-5` env pin, spawned as Agent `haiku` (59,154 tok); resolver `recipe claude:claude-haiku-5-5 declared:provider-pin`. Driver note: codex unusable → haiku.
- Evaluator: `sonnet` via Agent (97,223 tok), `agent-tool sonnet declared:alias-pin`. **Spawned as required; no role failed to spawn.**
- Quorum: none (the doc was reviewed at 399). Metered total: $0.
- Judge independence: cross-model (opus plan, haiku execution, sonnet judge).

**Progress:** 1.0 clauses 1–7 MET; goal unmoved. Row 243 V3 has its law: the code that will judge a key's confirmation of a finish prompt exists and is proven, and the prompt itself now lives where the card and the store can render it byte-identically.

**Ruled out:**
- **Building M2 (the judges) in S1:** S1 compiled to 1,626 > 1,400. M1 + M2 (1,360, margin 40) would have delayed M3, which S2 and S3 need sooner.
- **VT2 (`FinishConfirmRefusal`) in S1:** §6's listing is stale r0 text; §2.3 [r1] and §8 place it in S3.
- **The B-only branches (F6, `DraftOnly`):** D-WORLD-118 = A.
- **Treating the `verify_go` reds as sprint regressions:** neither package is touched, and both tests passed alone on the same head.

**Retro:**
- **Instance 2 of "test lines dominate the size overrun"** (398: 1,255 vs ≈ 1,135; 400: 1,626 vs ≈ 935 compiled). The design's raw estimates undercount tests by about 1.7× on law sprints. No rule change yet. If a third instance comes, the designer directive should size tests from the analogous landed sprint's measured test/prod ratio.
- **Instance 2 of "an untouched-area exec-timeout test reds verify_go while the judge runs concurrently"**: `TestWorkspaceExecConfiguredBindsTheEpisodeHandler` is new to the flake list. It passed alone in 0.86 s. Added to the next directives' known list.
- The 398 fire shape held: planner 57 min → round trip → executor 3 min → gates concurrent with the judge (30 min).

**Gate 5 bookkeeping:** rotation not due (#574 was opened 2026-11-16, after the Monday 07:00 boundary; the next rotation is Monday 2026-11-23). The weekly external-issue sweep already ran this week at 393.

**Next:**
1. Row 243 **V3 sprint 2** planner: the card page against fixtures (≈ 825 raw: M0 V2b pins, M4 `ConfirmationInput`/`ConfirmationCard`, M5 page + `confirm.js`, M6 XC1). Fold S1b's compiled judges (326 lines, `world-iter400/plan-scratch/patches/S1b-M2.patch`) in only if the compiled total is ≤ 1,400; otherwise S1b runs first as its own small sprint. Carry S1's Appendix C (NUL in `Validate` → the card; audit notes 1–5 to S3/S5).
2. Then S3 (credentials schema 9 + read routes), which is now routable under D-WORLD-118 = A.

## 401 — 2026-11-17 — row 243 V3 sprint 2 LANDED with S1b folded in (PR #593, merge `2b5cc609`): the confirmation card — `ConfirmationCard`, the embedded `/surface/confirm` page and `confirm.js` (digest-checked, `CONFIRM_FIELDS`, text-only, the 86-byte `world-verb/v1` message built in the page), XC1 in live Chrome, the V2b pins, and the verb judges; CI red once on S2's own node harness (a tick-count wait), fixed as M7; judge PASS 92 [PRODUCT]

**Pick:** row 243 V3 sprint 2, the planner step named as Next in record 400. That was also the queue head.
- **Re-measured at pick** (`origin/dev` `07d7c90e`):
  - None of S2's symbols existed (`ConfirmationCard`, `CONFIRM_FIELDS`, `confirm.js`, `JudgeVerbSSH`: 0 hits in `host cmd`). Control: `FinishPromptText` was found.
  - No open fleet PR, and no stale worktree for row 243 V3 S2. The design was quorum-reviewed at 399, so no pick-time quorum was owed.
- **Clause map:** 1.0 clauses 1–7 MET, none UNMET. Release rows 150/151 still wait on attended steps (A1; row 188 N2d/N4a), unchanged since 400 (no directive or attended commit since). Row 243 serves D-WORLD-94 and moves no clause.
- **Ledger at pick:** valid, 105 rows, no OPEN row.

**Gate 0:**
- Stamped first. `gh` is `sunholo-voight-kampff`. Billing CLEAN. Kill switch absent.
- Directives: 0 allowlisted on #574 (9 comments). The watermark file is absent, so every comment was read, and there was nothing to `--mark`.
- Self-notices: the script ran with rc 0, but printed `watermark channel unavailable`: both watermark files (#574, #524) are absent, so this is the same instrument gap as 371–400, not a verdict. #574's comments, read by hand, carry no driver crash notice.
- Inbox: no unread message for `mission-world` (`messages list -unread -inbox mission-world`).
- Skill drift: none. `SKILL.md` and all 13 `resources/*.md` at the resolved symlink target `cmp` equal to the driver clone's `origin/dev` (control: a nonexistent path fired).
- Dev CI at `07d7c90e`: 3/3 success (`check-runs` total 3).

**Planner** (Agent `opus`, 86 min, 396,179 tok; resolver `agent-tool opus fail-closed:env-pin`; `derive-planner-lane` → `opus fail-closed:env-pin`). It wrote `w-world-address-bar-verbs-v3s2-sprint-plan.md` and its JSON, committed as `21abef59`.
- **S1b folded by measurement:** S2 alone compiled to 1,051 (design ≈ 825, 1.27×). Iteration 400's `S1b-M2.patch` (`775b927ad56aa4c0`) applied cleanly to `07d7c90e` at 326 lines. The mutant run found two S1b survivors (`JV-ADMIT-BYTES`, `JV-RCK-NOW`), and 3 test lines kill them, so S1b came to 329 and the total to 1,380.
- Sub-milestones: M0 68 · M4 293 · M5a 711 · M5b 942 · M6 1,051 · S1b 1,380.
- **Mutants:** 145 rows (3 controls, 138 killed, 3 string-pin survivors by design — VP12's browser-only `app.js` strings — and 1 killed only by the static-asset test and live Chrome).
- **Six §1.4 premises found false or incomplete:**
  - F1: A4's field-table test reads the served wire, but S2 has no route, so it reflects over `VerbRequest`.
  - F2: A4's page is tested only in Chrome, so S2 adds a node stand-in DOM that runs in CI.
  - F3: the type is `Confirmation` because the function is `ConfirmationCard`.
  - F4: VP12's `app.js` strings are pinned as source strings.
  - F5: no S2 route applies a CSP, so XC1 copies the daemon's real headers at run time.
  - F6: five new test files sit outside design §7.
- **Controller round trip:** byte-identical. I extracted the six patches from the plan, applied them to a fresh `git archive 07d7c90e` and ran `diff -r` against the planner's `work` tree: only `.git` differed. Control: 13 paths differ from base, matching plan §5.

**Executor** (env `recipe claude:claude-haiku-5-5`, spawned as Agent `haiku` as in 374–400; one foreground run, 3.2 min, 54,579 tok).
- M0 → M4 → M5a → M5b → M6 → S1b: every `--dry-run` was clean and every acceptance `ACCEPT_OK`. I re-grepped each log: PASS 7/3/21/10/15/15, 0 `--- FAIL`, 0 `--- SKIP`, 0 `^FAIL`. Snapshots matched §7 (2/4/10/11/12/14, DELETED empty).
- Commits by `commit_ms.sh` (gofmt/build/vet/test before each): M0 `650de9cf` 68 · M4 `38329114` 293 · M5a `ef15ea19` 711 · M5b `5bd02067` 942 · M6 `dc3ac827` 1,051 · S1b `1b20f08a` 1,380.
- **My own instrument slip, caught and undone:** M6's commit-time test omitted `WORLD_CHROME`, so XC1 failed with "browser required" and `commit_ms.sh` refused. The cumulative S1b snapshot then committed M6's file inside S1b. I reset the unpushed branch to M5b and re-committed M6 (XC1 PASS 3.29 s) and S1b separately. The final tree is identical to the merged-commit tree (`tree-equal: yes`).

**Controller gates** (plan §6 `gates.sh` on `1b20f08a`, 977 s, run while the judge read):
- rc 0: status, `verify_ail`, gofmt, build, vet, fence, pkgs (309 s), race, **node** (9 named), **inv** (43 named passes), **chrome** (17 incl. XC1).
- **Plan gap:** §6's script has no leg for `host/authority` (S1b's package), though its text says "S1b's judges". I added an `authority` leg: 102 `--- PASS`, plus `-race` on `host/authority` and `host/approval`, rc 0.
- **`verify_go` rc 1** on one untouched-area test: `host/broker` `TestExecProbeForeignListenerOnTheReversePortDoesNotRefuse/chatty`, at load ≈ 16. It **passed alone on the same head** (5.74 s). The sprint touches no `host/broker` file.

**Evaluator r1** (Agent `sonnet`, 24 min, 118,716 tok; `agent-tool sonnet declared:alias-pin`; its own worktree `.wt-world-iter401-row243-v3s2-eval`): **PASS 92/100**, 0 blocking.
- Scores: correctness 28/30, tests 25/30, security 19/20, standards 10/10, scope 10/10.
- **Its own mutants:** 49 `E-` ids, 0 shared with the plan's 146.
  - Must-die control `E-JS-DIGESTCHK` (the digest check removed) was killed by `TestConfirmPageFailsClosed`.
  - Must-live controls survived.
  - Two constant mutants were killed only by the golden, and it recorded that.
- **Survivors, non-blocking:**
  - `null` request bytes leave the page on "Reading the request…", still with no controls.
  - The summary/pending text is unpinned.
  - Five S1b judge-input mappings are unpinned: `EnrolledAt` ×2, a judge-level clock-behind case, `LastUsedAt`, `BackupEligible`.
- **Other checks:**
  - Fail-closed on every class.
  - `CONFIRM_FIELDS` is read both sides.
  - No HTML sink anywhere.
  - The golden is Go-written.
  - No route: proven by `TestConfirmRoutesNotRegistered` plus a killed register-route mutant.
  - S1b step order matches design §2.3.
  - Only `why_admit_test.go` among pre-existing tests was edited (+1 row).
  - Intermediate commits M5a and M6 are green alone.
- **Reproduced first-party before banking:** in the judge's worktree, `confirm.js:141`'s digest check replaced by `if (false)` is killed by `TestConfirmPageFailsClosed` with the static and golden tests skipped (rc 1). After restoring, the tree is clean and the test passes.

**Gate 3b, attempt at `1b20f08a` — RED on the sprint's own test:** PR #593, CI 2/3. The host job failed `TestConfirmPageMessage` (`'loading' !== 'ready'`) and `TestConfirmPageStates` (`'loading'`, expected `closed`); `…FailsClosed`/`…RendersText` passed (run 44116934040).
- **Cause, found by the planner (resumed by SendMessage; it authored the harness):** `open()` waited a fixed **400 `setImmediate` ticks** for the page to leave `loading`. The page needs **two** `crypto.subtle` digests, which resolve off the event loop (node's thread pool), so a slower runner ran out of ticks.
  - **Reproduced both ways** by delaying the page's digest inside the harness's vm context: at 3 ms the base fails exactly CI's pair; at 50/200 ms all four tests fail; with the fix, 0/50/200 ms all pass.
- **M7, test-only:** `open()` now polls node's own `setTimeout(…, 2)` against a 10 s wall-clock deadline and fires the page's stubbed timers once the page has fetched and set its abort timer. No assertion and no production file changed.
  - 19 of 19 re-checked mutants still die.
  - Plan addendum `c268c997`, commit `57778071`.
  - M5b acceptance re-run: 10 PASS. `go test -count=3 ./host/surfacepage/` is green.
- **Size:** 1,382 by the stated unit. The plan section said 1,390 because it added M7's gross +6/−4. The file is new since base, so its net cumulative change is +2. The judge confirmed 1,382.

**Evaluator r2** (the same judge, resumed by SendMessage, 2.0 min, 130,092 tok cumulative): **PASS 92**, 0 blocking.
- **Own instrument** (`judge/r2/repro.txt`): old `open()` fails `states` at 3 ms and all four sections at 50 ms; M7 passes at 0/3/50 ms.
- **Re-ran 8 of its own mutants on M7**, including the must-die control and two abort mutants (`E-JS-ABORT-TIMER-M7`, `E-JS-ABORT-NOSIGNAL-M7`). All 8 were killed with the static and golden tests in and out.
- `refused()` still asserts the exact state, so a stuck `loading` cannot pass vacuously. The one new cost: a genuinely hung page now takes a bounded 10 s to fail.

**Gate 3b, attempt at `57778071`:**
- **Attempt 1: 2/3.** The host job was red on `scripts/marketingcheck` `TestMarketingSupervisorFakeDaemonCleanup/shell-int` (`signal: killed`, cleanup 29.93 s over its bound; run 44120752559). This is the row-254/258 supervisor-cleanup class, with three earlier instances in this log.
  - The sprint touches no `scripts/` file.
  - Negative control: the same test was `ok` (35.06 s) on the `1b20f08a` run, which has the same `scripts/` tree.
  - `surfacepage` and `authority` were `ok` in this run.
- `gh run rerun --failed` → **attempt 2: 3/3 success**, polled on `actions/runs/44120752559` attempt 2 (`ci_pr593_run44120752559_att2_s2.verdict`).
- `mergeable` MERGEABLE/CLEAN on the head. Merged with `--match-head-commit` → **`2b5cc609`**, tree byte-identical to the head.
- Merge push CI: **3/3 success on attempt 1** (`ci_merge_2b5cc609_s2.verdict`).

**Routing evidence:** base=2b5cc60904e82a1baa525eb95490022563e78b28@2026-11-17T16:37:19Z (Gate 4). The Gate-1 base was 07d7c90e5979893154d2244037bedeffd42386fe@2026-11-17T13:40:07Z, and the sprint worktree was cut from it (no drift).
- Controller: `claude-opus-5-5` (tok: not reported).
- Designer: not spawned (doc reviewed at 399).
- Planner: `opus` via Agent (396,179 tok + the M7 resume; the harness reports cumulative 418,718). Driver note: codex `gpt-6.1-sol` unusable at preflight (probe rc 75, quota unreadable) → opus.
- Executor: `claude:claude-haiku-5-5` env pin, spawned as Agent `haiku` (54,579 tok); resolver `recipe claude:claude-haiku-5-5 declared:provider-pin`. Driver note: codex unusable → haiku. M7 was applied by the controller with the plan's own `patch -p1 -E` heredoc, the executor's protocol, because it was one test-file patch.
- Evaluator: `sonnet` via Agent (r1 118,716 tok; r2 cumulative 130,092), `agent-tool sonnet declared:alias-pin`. **Spawned as required; no role failed to spawn.**
- Quorum: none. Metered total: $0.
- Judge independence: cross-model (opus plan and fix, haiku execution, sonnet judge on both rounds).

**Progress:** 1.0 clauses 1–7 MET; goal unmoved. Row 243 V3 is at 2 of 6 sprints, and the judges S4/S5 need are built. The confirmation card exists and is proven against fixtures, in node and in live Chrome. It is not reachable until sprint 3 registers its routes.

**Ruled out:**
- **Leaving S1b for its own iteration:** S2 measured 1,051 and S1b 329, so 1,380 fits the 1,400 cap.
- **Treating the `1b20f08a` CI red as a flake:** it was this sprint's own test. It was reproduced by an injected delay and fixed.
- **Treating the `57778071` red as a regression:** it was in an untouched `scripts/` test that was green on the previous head with the same tree, and it is a known class.
- **Treating the `verify_go` red as a regression:** untouched `host/broker`; it passed alone.

**Retro:**
- **New, instance 1: "a node harness bounded by a tick count is a timing assumption."** A page whose state change waits on `crypto.subtle` (thread pool) cannot be bounded by `setImmediate` ticks. Local node v26 passed, and the CI runner did not. The planner's directive should say: *harness waits on async page work use a wall-clock deadline, and are proven with an injected delay*. Two pre-existing tick waits (`script_test.go:417`, `app_verb_test.go:169`) wait only on in-loop promise chains and were left alone. Recorded; the rule is added to the next planner directive.
- **New, instance 1: "the plan's gate script omits a package the sprint touches"** (no earlier instance found by grep of log + archive): here, `host/authority`. The controller should diff the gate script's package list against `git diff --name-only` before running it. Applied this iteration by hand.
- **Controller instrument slip:** a commit-time test for a Chrome-gated test must carry the plan's exact env (`WORLD_CHROME`). `commit_ms.sh` refused correctly. The danger was the next cumulative snapshot silently absorbing the refused milestone. The rule: after any `commit_ms.sh` refusal, stop the chain and do not commit the next snapshot.
- **Size accounting:** compute a fix's numstat after the commit, by the base-relative unit, never as the sum of the fix's own +/−. Here that gave 1,382, not 1,390.

## 402 — 2026-11-17 — row 243 V3 sprint 3 LANDED (PR #595, merge `c9e997a9`): credentials schema 9 `verb_requests` (put/supersede/cap/prune, one-transaction consume whose refusals write nothing), `broker.FinishConfirmRefusal` F1–F5, and the three confirm GETs (page, script, card JSON) plus `requests.json` `confirmations`; no POST, no CLI change; 1,391 lines; judge PASS 92 [PRODUCT]

**Pick:** row 243 V3 sprint 3, the planner step named as Next in record 401. That was also the queue head.
- **Re-measured at pick** (`origin/dev` `9733813c`): no open fleet PR and no stale worktree for row 243 V3 S3; `mission_decisions.sh --open` → no OPEN row (D-WORLD-118 = A since 400). The design was quorum-reviewed at 399, so no pick-time quorum was owed.
- **Clause map:** 1.0 clauses 1–7 MET, none UNMET. Release rows 150/151 still wait on attended steps, unchanged since 401 (no directive or attended commit since). Row 243 serves D-WORLD-94 and moves no clause.

**Gate 0:**
- Stamped first. `gh` is `sunholo-voight-kampff`. Billing CLEAN. Kill switch absent.
- Directives: 0 allowlisted on #574 (10 comments). The watermark file is absent, so every comment was read, and there was nothing to `--mark`.
- Self-notices: **rc 2** (`watermark channel unavailable`: both watermark files, #574 and #524, are absent) — an instrument failure, not a verdict; same gap as 371–401.
- Inbox: no unread message for `mission-world` (`messages list -unread -inbox mission-world`).
- Skill drift: none. `SKILL.md` and all 13 `resources/*.md` at the resolved symlink target `cmp` equal to the driver clone's `origin/dev` (per-file verdicts printed).
- Dev CI at `9733813c`: 3/3 success (`check-runs` total 3).

**Planner** (Agent `opus`, 92 min, 455,957 tok; resolver `agent-tool opus fail-closed:env-pin`; `derive-planner-lane` → `opus fail-closed:env-pin`). It wrote `w-world-address-bar-verbs-v3s3-sprint-plan.md` and its JSON, committed as `685bc668`.
- **No cut:** M7 compiled first at 836; the first full compile was 1,408; trimming comments and one duplicate test brought it to 1,387, and two tests added after the mutant runs to 1,391 (margin 9). The measured alternative (S3a = M7 + M8 at 913, S3b = M9 at 478) was not needed.
- Sub-milestones: M7a 377 · M7b 805 · M8 913 · M9a 1,235 · M9b 1,391.
- **Mutants:** 176 rows (166 killed, 2 route-table, 3 controls, 5 named equivalents), 276 ledger lines over runs R1–R4. R1 found a redundant digest-hex check (removed); R2's survivor (prune pending rows only) got a 2-line test.
- **§1.4 premises found false:** F1 the design's DDL does not refuse a `decision_ref` on a `superseded` row (kept verbatim; inert — only consume writes it; the actual behaviour is pinned); F2 CS7's "constant set to 8" cannot be built in one binary, so a v10 store under the v9 build exercises the same branch; F3 the bump forces 19 fixture edits in 14 files plus 5 store-test edits outside design §7; F4–F5 audit notes 1–2 (three GETs, not four; CR1–CR4); F6–F8 a new `ErrBadVerbRequest`, a new `finish_confirm.go`, a new assets test.
- **Audit note 3 (raced supersession):** S3 pins the store half both ways (put-then-consume refuses `ErrVerbRequestSuperseded` and writes nothing; consume-then-put leaves the approved row). The 409 mapping and the window-table row are S4's.
- **Controller round trip:** byte-identical. I extracted the five patches from the plan, applied them to a fresh `git archive 9733813c` and ran `diff -r` against the planner's final tree: only `.git` differed. The plan's §6 touched-package list equals my own derivation from the diff (6 packages) — iteration 401's gate-script gap did not recur.

**Executor** (env `recipe claude:claude-haiku-5-5`, spawned as Agent `haiku` as in 374–401; one foreground run, 18 min, 53,490 tok).
- M7a → M7b → M8 → M9a → M9b: every `--dry-run` clean, all nine acceptance parts `ACCEPT_OK`. I re-grepped each log: 0 `--- FAIL`, 0 `--- SKIP`, 0 `^FAIL`. Snapshots matched §7 (21/22/24/28/30 STATUS lines; DELETED 1 from M9a). The executor's tree was byte-identical to the planner's final tree.
- Commits by `commit_ms.sh` (gofmt/build/vet + the part's `fast` acceptance before each; no refusal): M7a `a4d9f852` 377 · M7b `4424ec73` 805 · M8 `cd0d91d4` 913 · M9a `d614fe3d` 1,235 · M9b `21fa7d54` 1,391 — each equal to the plan.

**Controller gates** (plan §6 `gates.sh` on `21fa7d54`, 887 s, concurrent with the judge's reading phase):
- rc 0: status, `verify_ail`, gofmt, build, vet, fence, pkgs (281 s), race, **race3** (`-race -count=3` over the new store tests), node, **inv** (65 named passes; `TestConfirmRoutesNotRegistered` gone), chrome.
- **`verify_go` rc 1** on one untouched-area test: `host/verifygate` `TestMissionConfigDefaultFlowIntegration/valid_ledger` (fixture timeout, on the known load-sensitive list). It **passed alone on the same head** (controller 7.8 s; judge 8.0 s).

**Evaluator** (Agent `sonnet`, 26 min, 107,905 tok; resolver `agent-tool sonnet declared:alias-pin`; its own worktree `.wt-world-iter402-row243-v3s3-eval`). **PASS 92/100**, 0 blocking (correctness 28, tests 27, security 19, standards 9, scope 9).
- 55 ledger runs: 48 killed (all by behavioural tests), 6 survived, 1 compile error (re-written). 41 of its ids are not in the plan's ledger (`comm -23`). Must-die `J-CTL-DIE-F` killed; must-live `J-LIVE-UPGRADE-BOUND` and `J-LIVE-CONSUME-AND-NULL` survived (equivalent). 12 plan rows reproduced; `P-LIST-LIVE` killed in the daemon package as the plan says.
- Pre-existing test edits: all 18 are what the bump forces; no assertion weakened. M7a and M9a green on their own.
- **Non-blocking (logged, not filed — D-WORLD-109):** (1) no test pins the exact-expiry-second boundary of the pending list (`J-LIST-EXPIRY-BOUND`, `>`→`>=`) or of the cap count (`J-PUT-COUNT-BOUND`); the code is consistent with consume's `now >= ExpiresAt`. **The controller reproduced the list survivor** (store rc 0 and daemon `^(TestConfirm|TestRequestsJSON)` rc 0 with the mutant applied). (2) A day-expired request of the same episode is superseded rather than pruned on the next put (harmless, one row per episode, undocumented in §2.4). (3) F1 confirmed and accepted.

**Gate 3b:** PR [#595](https://github.com/sunholo-data/ailang-world/pull/595) head `21fa7d54`: **3/3 success on attempt 1** (`ci_pr595_head21fa7d54_s2.verdict`, the second 540 s segment). `mergeable` MERGEABLE/CLEAN on the head. Merged with `--match-head-commit` → **`c9e997a9`**, tree byte-identical to the head. Merge push CI: **3/3 success on attempt 1** (`ci_merge_c9e997a9_s2.verdict`, second segment).

**Routing evidence:** base=c9e997a9c3bf8197356d6af70f6e03c345a9fdf3@2026-11-17T20:20:05Z (Gate 4). The Gate-1 base was 9733813cd8ab41313c158f83d3a0d2bc7c7d7d7f@2026-11-17T17:40:24Z, and the sprint worktree was cut from it (no drift).
- Controller: `claude-opus-5-5` (tok: not reported).
- Designer: not spawned (doc reviewed at 399).
- Planner: `opus` via Agent (455,957 tok). Driver note: codex `gpt-6.1-sol` unusable at preflight (probe rc 75, quota unreadable) → opus.
- Executor: `claude:claude-haiku-5-5` env pin, spawned as Agent `haiku` (53,490 tok); resolver `recipe claude:claude-haiku-5-5 declared:provider-pin`. Driver note: codex unusable → haiku.
- Evaluator: `sonnet` via Agent (107,905 tok), `agent-tool sonnet declared:alias-pin`. **Spawned as required; no role failed to spawn.**
- Quorum: none. Metered total: $0.
- Judge independence: cross-model (opus plan, haiku execution, sonnet judge).

**Progress:** 1.0 clauses 1–7 MET; goal unmoved. Row 243 V3 is at 3 of 6 sprints. A confirmation request can now be stored and its card read on the surface; nothing can approve one until sprint 4's POST.

**Ruled out:**
- **Cutting S3 into S3a/S3b:** measured 1,391 ≤ 1,400 with every part proven; the cut (913 / 478) was costed and not needed.
- **Treating the `verify_go` red as a regression:** untouched `host/verifygate`, listed load-sensitive, green alone on the same head (two independent re-runs).
- **Fixing the expiry-second boundary in this PR:** a 9-line margin and a non-blocking, display-only one-second skew; carried to S4's planner, whose consume/POST work touches the same rows.

**Retro:**
- **Held, instance 2 of the planner-directive rules earned at 401:** the gate script derived its touched-package list from the diff and FAILED on a mismatch (it matched); the commit-time commands carried their env, and `commit_ms.sh` refused nothing. Both 401 slips were closed by putting the rule into the directive, not into the controller's memory.
- **New, instance 1: "a boundary on a time comparison is a mutant nobody enumerates."** Both survivors are `>` vs `>=` on `expires_at` in SQL strings. The planner's diff-anchored enumeration took each SQL predicate as one mutant (removing it), not its comparison operator. For S4: enumerate every time comparison's operator (`>`/`>=`, `<`/`<=`) as its own mutant, with an exact-boundary test row. Recorded; added to the next planner directive.
- **Shape that worked (third time):** planner opus → controller round trip → executor haiku → commit_ms per snapshot → gates in background concurrent with a judge that waits on `ALL_DONE.<sha>`. 2 h 40 min to merge.

## 403 — 2026-11-18 — row 243 V3 sprint 4 LANDED (PR #597, merge `00fa924f`): `POST /surface/confirm/{id}` (C0–C15 in the design's order; nothing written before C14; the writer held across `RecordOperatorApproval` and the one-transaction consume), `surfacepage.Confirmed`, the window tests CX14–CX21 with audit note 3's raced supersession, the exact-expiry-second pins, and `confirm.js`'s submit half; no CLI, schema or `.ail` change; 1,097 lines; judge PASS 91 [PRODUCT]

**Pick:** row 243 V3 sprint 4, the planner step named as Next in record 402. That was also the queue head.
- **Re-measured at pick** (`origin/dev` `430668f8`): no open PR from this mission and no stale worktree for row 243 V3 S4 (`git log origin/dev --grep "V3 S4"` empty; the three `iter402` worktrees are merged leftovers); `mission_decisions.sh --open` → no OPEN row; `--check` → ledger valid, 105 rows. The design was quorum-reviewed at 399, so no pick-time quorum was owed.
- **Clause map:** 1.0 clauses 1–7 MET, none UNMET. Release rows 150/151 still wait on attended steps, unchanged since 402 (no directive or attended commit since). Row 243 serves D-WORLD-94 and moves no clause.
- **The previous fire:** 402's slot verdict reads `REAPED at=gate-5` (21:34Z, rc 0, 14,190 s). Nothing was lost: its #574 report (21:22Z), its controlplane fyi and its memory update all landed, and its record PR #596 had merged (`430668f8`, after 3 race-job attempts — instance 1 of the CI race headroom note). Not an orphan; no log entry is missing.

**Gate 0:**
- Stamped first. `gh` is `sunholo-voight-kampff`. Billing CLEAN. Kill switch absent.
- Directives: 0 allowlisted on #574 (12 comments). The watermark file is absent, so every comment was read, and there was nothing to `--mark`.
- Self-notices: **rc 2** (`watermark channel unavailable`: both watermark files, #574 and #524, are absent) — an instrument failure, not a verdict; same gap as 371–402. #574's newest comments, read by hand: the 401 and 402 reports and the driver's 402 `REAPED at=gate-5` notice.
- Inbox: no unread message addressed to `mission-world` beyond 402's own controlplane fyi.
- Skill drift: none. `SKILL.md` and all 13 `resources/*.md` at the resolved symlink target `cmp` equal to the driver clone's `origin/dev` (per-file verdicts printed).
- Dev CI at `430668f8`: 3/3 success (`check-runs` total 3). The main checkout was 9 behind origin and clean apart from untracked worktree dirs; fast-forwarded.

**Planner** (Agent `opus`, 65 min, 487,331 tok; resolver `agent-tool opus fail-closed:env-pin`; `derive-planner-lane` → `opus fail-closed:env-pin`). It wrote `w-world-address-bar-verbs-v3s4-sprint-plan.md` and its JSON, committed as `4a00462c`.
- **No cut:** M10 compiled first at 916; M11 added 181; total **1,097** (311 production, 786 test). The design's "M11 → S6" cut was measured and not needed.
- Sub-milestones: M10a 433 (handler unregistered, `Confirmed`, the fourth writer site, CX0–CX6) · M10b 605 (CX7–CX13, the boundary pins) · M10c 917 (CX14–CX21, audit note 3, C15's mapping) · M10d 916 (the POST registered; tests through the mux) · M11 1,097 (submit half, CP8–CP10).
- **Mutants:** 150 rows (141 killed, 3 controls, 4 named equivalents, 2 killed only by the route-table test), 155 ledger lines over R1–R3. **Every time comparison's operator was its own mutant** (the 402 retro rule), including S3's two survivors, both now killed by `TestVerbRequestExpiryBoundary`.
- **§1.4 premises found false:** F1 the design's CX19 file holds an operand table; A4's concurrency proof is in the retry test file (both shapes used); F2 below v9 is refused at C7 (read-only lookup), not C13; F3 two cites moved (`ParseApproveBody` `:214`, routes `daemon.go:1141-1153`); F4 CX21's CLI half has no CLI before S5 (routed there); F5 a failed C14 still writes the request object and moves the approvals head (CX20 asserts 0 decisions, the meaning the design intends).
- **Decisions the judge was asked about:** S4-D6 (the C14/C15 500 carries the design's fixed text, not `internalErrorMessage`) and S4-D7 (three package-variable test seams in production).
- **Base gates:** the planner's `pkgs` and `race` legs were red at the base under its own load (`TestAGUISlowReaderCutIsResumable`, `TestConvertedStoreMethodsM5b`); its draft `race` leg chained with `&&`, so the daemon subset never ran at the base — it caught this itself, fixed the script, and ran that subset alone at the base (green, 65 PASS).
- **Race-mode cost:** the new daemon tests add ≈ 2.1 s under `-race` (directive item 8, the 402 headroom note).
- **Controller round trip:** byte-identical. I extracted the five patches from the plan, applied them to a fresh `git archive 430668f8` and ran `diff -r` against the planner's final tree: only `.ailang` build caches differed. The touched-package list (daemon, store, surfacepage) equals my own derivation from the patch headers.

**Executor** (env `recipe claude:claude-haiku-5-5`, spawned as Agent `haiku` as in 374–402; one foreground run, 12 min, 50,645 tok).
- M10a → M10b → M10c → M10d → M11: every `--dry-run` clean, all eight acceptance parts `ACCEPT_OK` (the two `daemon` parts 258 s and 269 s, each one call). I re-grepped each log: 0 `--- FAIL`, 0 `--- SKIP`, 0 `^FAIL`. Snapshots matched §7 (5/6/7/9/12 STATUS lines, DELETED 0). The executor's tree was byte-identical to my round-trip tree.
- Commits by `commit_ms.sh` (gofmt/build/vet + the part's `fast` acceptance, and M11's `chrome`, before each; no refusal): M10a `206335ee` 433 · M10b `29bb39a5` 605 · M10c `af892b6f` 917 · M10d `47d0a1df` 916 · M11 `eff19d24` 1,097 — each equal to the plan.

**Controller gates** (plan §6 `gates.sh` on `eff19d24`, 926 s, concurrent with the judge's reading phase): **all 13 legs rc 0** — status, `verify_ail`, gofmt, build, vet, fence, pkgs (255 s), race, race3, node, inv (70 named passes), chrome, and **`verify_go` (493 s, green — first time in four sprints with no load-sensitive red)**. Touched packages derived from the diff = the plan's. `/surface` POST registrations 2 → 3.

**Evaluator** (Agent `sonnet`, 30 min, 97,129 tok; resolver `agent-tool sonnet declared:alias-pin`; its own worktree `.wt-world-iter403-row243-v3s4-eval`). **PASS 91/100**, 0 blocking (correctness 28, tests 26, security 19, standards 9, scope 9).
- 50 ledger runs, all its own ids (`comm -12` vs the plan's ids = 0). Must-die controls `J-C9-EQ`, `J-C15-BEFORE-C14` (a real C14/C15 swap, killed by four tests) and `J-STORE-CAP-GE` killed; must-live `J-CTL-LIVE` and `J-JS-CTL` survived. Intermediate commits `206335ee` (0 confirm POST registrations) and `47d0a1df` checked in a throwaway worktree.
- S4-D6 acceptable (replacing the constant with `err.Error()` is killed); S4-D7 acceptable (seams assigned only in tests, restored by `t.Cleanup`; `-race -count=3` clean), with a hermeticity note.
- **Non-blocking (logged, not filed — D-WORLD-109):** (1) C14's `now` is not asserted — `J-NOW-C14` (`now+1`) survives the whole package (258 s). **The controller reproduced it** (narrow rc 0 with the mutant; must-die `J-C9-EQ` rc 1 in the same call). Carried to S5's planner. (2) `confirm.js` accepts `confirmed:true` without also checking status 200 (defence in depth). (3) The busy 503's reason text is unasserted. (4) A second POST route is caught only by the static route-table test (an unwrapped one cannot compile).

**Gate 3b:** PR [#597](https://github.com/sunholo-data/ailang-world/pull/597) head `eff19d24`: **3/3 success on attempt 1** (`ci_pr597_headeff19d24_s2.verdict`, the second 540 s segment). `mergeable` MERGEABLE/CLEAN on the head. Merged with `--match-head-commit` → **`00fa924f`**, tree byte-identical to the head. Merge push CI: **3/3 success on attempt 1** (`ci_merge_00fa924f_s2.verdict`, second segment).

**Routing evidence:** base=00fa924fdefb93f0f12e2037cd09f0f8949e6136@2026-11-18T00:02:03Z (Gate 4). The Gate-1 base was 430668f86d4bdfb226f92ec39de99440abbbb397@2026-11-17T21:40:07Z, and the sprint worktree was cut from it (no drift at Gate 3).
- Controller: `claude-opus-5-5` (tok: not reported).
- Designer: not spawned (doc reviewed at 399).
- Planner: `opus` via Agent (487,331 tok). Driver note: codex `gpt-6.2-sol` unusable at preflight (probe rc 75, quota unreadable) → opus.
- Executor: `claude:claude-haiku-5-5` env pin, spawned as Agent `haiku` (50,645 tok); resolver `recipe claude:claude-haiku-5-5 declared:provider-pin`. Driver note: codex unusable → haiku.
- Evaluator: `sonnet` via Agent (97,129 tok), `agent-tool sonnet declared:alias-pin`. **Spawned as required; no role failed to spawn.**
- Quorum: none. Metered total: $0.
- Judge independence: cross-model (opus plan, haiku execution, sonnet judge).

**Progress:** 1.0 clauses 1–7 MET; goal unmoved. Row 243 V3 is at 4 of 6 sprints. A key can now confirm a stored finish request on the surface (the world decision is written and the request consumed); nothing creates a request outside tests until sprint 5's `episode finish --confirm-on-surface`.

**Ruled out:**
- **Cutting M11 to S6:** measured 1,097 ≤ 1,400 with every part proven; the cut was not needed.
- **Treating the planner's red base `pkgs`/`race` legs as a regression:** both are untouched-area tests at the base under the planner's own parallel load; the final-tree gates were all green on the controller's run.
- **Fixing J-NOW-C14 in this PR:** non-blocking, a timestamp no consumer reads in S4; S5's CLI read-back is where the decision's time is first consumed, so its planner pins it.

**Retro:**
- **Held, instance 2 of the 402 rule "every time comparison's operator is its own mutant":** the planner enumerated each `>`/`>=` in Go and SQL; S3's two survivors are now killed by an exact-boundary test. Closed by putting the rule in the directive.
- **Held, instance 3 of the planner-directive rules from 401:** touched-package list derived from the diff and checked (matched); commit-time commands carried their env (`commit_ms.sh` refused nothing).
- **New, instance 1: "an unasserted argument passed to a writer survives".** The judge's only whole-package survivor is the `now` passed to C14's world write; every other C14 argument (requester, scope, approver) was pinned by CX17. For S5's planner: enumerate every argument of every store/world write call as its own mutant, not only the call's presence.
- **Shape that worked (fourth time):** planner opus → controller round trip → executor haiku → commit_ms per snapshot → gates in background concurrent with a judge that waits on `ALL_DONE.<sha>`. About 2 h 10 min from Gate 1 to merge.

**Gate 5 bookkeeping:** rotation not due (#574 opened 2026-11-16, after the Monday 07:00 boundary; next rotation Monday 2026-11-23). The weekly external-issue sweep already ran this week at 393.

## 404 — 2026-11-18 — row 243 V3 sprint 5 LANDED (PR #599, merge `48b8b6b6`): `episode finish --confirm-on-surface` (refusals, the request put, the tty-free wait, the re-check over the process's own message, the decision read-back), CF1–CF14 incl. CX21's CLI half, C14's `now` pinned, the bar's `approve vr-` and finish reason, QUICKSTART §20 + CQ1; no route, schema or `.ail` change; 1,398 lines; judge PASS 91 [PRODUCT]

**Pick:** row 243 V3 sprint 5, the planner step named as Next in record 403. That was also the queue head.
- **Re-measured at pick** (`origin/dev` `18c3d132`): no open PR from this mission and no stale worktree for row 243 V3 S5 (`git log origin/dev --grep "V3 S5"` empty; the three `iter403` worktrees are clean merged leftovers); `mission_decisions.sh --open` → no OPEN row; `--check` → ledger valid, 105 rows. The design was quorum-reviewed at 399, so no pick-time quorum was owed.
- **Clause map:** 1.0 clauses 1–7 MET, none UNMET. Release rows 150/151 still wait on attended steps (no directive or attended commit since 403). Row 243 serves D-WORLD-94 and moves no clause.
- **The previous fire:** 403's slot verdict reads `REAPED_at=gate-5` (00:43Z, rc 0, 11,128 s) — the second in a row after 402's. Nothing was lost: 403's #574 report (00:19Z) posted and its record PR #598 merged (`18c3d132`). Not an orphan; no log entry is missing.

**Gate 0:**
- Stamped first. `gh` is `sunholo-voight-kampff`. Billing CLEAN. Kill switch absent.
- Directives: 0 allowlisted on #574 (13 comments) and 0 on #524 (47 comments, rotation `-prev`). Both watermark files are absent, so every comment was read; nothing to `--mark`.
- Self-notices: **rc 2** (`watermark channel unavailable`; the #574 and #524 files are absent) — an instrument failure, not a verdict; same gap as 371–403.
- Inbox: no unread message addressed to `mission-world`.
- Skill drift: none. `SKILL.md` and every `resources/*.md` at the resolved symlink target and in the pin (`022f68f1a`) `cmp` equal to the driver clone's `origin/dev`.
- Dev CI at `18c3d132`: 3/3 success (`check-runs` total 3). The main checkout equalled origin/dev.

**Planner** (Agent `opus`, 90 min, 502,955 tok; resolver `agent-tool opus fail-closed:env-pin`; `derive-planner-lane` → `opus fail-closed:env-pin`). It wrote `w-world-address-bar-verbs-v3s5-sprint-plan.md` and its JSON, committed as `000aa4c6`.
- **No cut (S5-D1):** M12 compiled first at 1,277 (1.75× its ≈ 730 raw); M13 59; M14 62; total **1,398** (323 production, 12 docs, 1,063 test). The design's "M13 → S6" cut was measured and not needed; "CQ1 stays (S7)" read as the charter's QUICKSTART guardrail (S5-D3).
- Sub-milestones: M12a 540 (the surface code with no caller, unit tests, two context-root pins) · M12b 815 (wait / re-check / read-back tests, the C14 `now` daemon pin) · M12c 1,277 (the flag, its refusals, the prompt hook, the real-daemon rig; the flag is user-visible only from here) · M13 1,336 · M14 1,398.
- **Mutants:** 118 run (111 killed by the named test and the package's no-static suite, 7 named equivalents, 0 survivors, 1 control); ledger `plan-scratch/mutants/ledger.txt`. Every argument of the put and every time comparison was its own mutant (the 402 and 403 retro rules).
- **§1.4 premises found false:** F1 the finish rig has no daemon (`finishClient = memoryDaemon(f.serve)`), so the rig adds a real `daemon.New` for `/v1/objects` and the card's POST on a loopback request with no listener; F3 the daemon's clock cannot be injected from `cmd/`, so one test waits for the next wall second (`nextWallSecond`, ≤ 1 s, 3 s cap); F4 `GET /v1/objects/{ref}` moved `daemon.go:1113` → `:1115`; F6 the golden's generator had no `vr-` query (two pre-existing test files outside design §7 gained one line each); F7 the 128 KiB check is already `Validate`'s; F8 no test enumerates the CLI's credential writer sites.
- **Carried items:** CX21's CLI half → `TestFinishConfirmOnSurfaceDispatchesTheNamedDecision`; judge-403 note 1 (`J-NOW-C14`) reproduced surviving S4's tests, now killed by `TestConfirmPOSTRecordsAtItsClock` and by CF3; S1 note 4 (CF12: the resume reuse skips the re-check and read-back — safe, S5-D9); S1 note 5 (the revocation window — fail-safe, S5-D10); judge-403 notes 2–4 not taken (2-line margin), carried in Appendix C.
- **Base gates:** `FAILS: none` (782 s). Its final gates: every leg green except two untouched `host/daemon` timing tests under load 13–19.
- **Controller round trip:** byte-identical. I extracted the five patches from the plan, applied them to a fresh `git archive 18c3d132` and ran `diff -r` against the planner's `work` and `final` trees: no difference. Numstat 1,398.

**Executor** (env `recipe claude:claude-haiku-5-5`, spawned as Agent `haiku` as in 374–403; one foreground run, 23 min, 73,626 tok).
- M12a → M12b → M12c → M13 → M14: every `--dry-run` clean; all ten acceptance parts `ACCEPT_OK` (the two `daemon` parts ≈ 339 s each, the three `cmd` parts ≈ 230 s, each one call). I re-grepped each log: 0 `--- FAIL`, 0 `--- SKIP`, 0 `^FAIL`. Snapshots matched §7 (3/4/6/14/16 STATUS lines, DELETED 0). The executor's tree was byte-identical to my round-trip tree.
- Commits by `commit_ms.sh` (gofmt/build/vet + the part's `fast` acceptance before each; no refusal): M12a `2f0539f2` 540 · M12b `34e97cd9` 815 · M12c `ba783c32` 1,277 · M13 `af6545d9` 1,336 · M14 `0b48604e` 1,398 — each equal to the plan.

**Controller gates** (plan §6 `gates.sh` on `0b48604e`, 1,499 s, concurrent with the judge's reading phase): touched packages derived from the diff = the plan's five; status, `verify_ail`, gofmt, build, vet, fence, race, race3, node, inv (71 named passes), chrome all rc 0; `DAEMON_ROUTES` 47 = base, `SURFACE_POSTS` 3, `UNTOUCHED_VIOLATIONS` none. **`pkgs` and `verify_go` red, only in untouched areas:** `TestMCPPostBudgetProductionConstants/batch` (host/daemon) and `TestMissionConfigDefaultFlowIntegration/invalid_ledger` (host/verifygate) each pass alone at the head; `TestSeToolsRunBenchmarkTasksEndToEnd` (host/daemon; its mock saw a stray `"GET  "` beside the one POST) fails alone at the head 1 of 3 **and alone at the base `18c3d132` 1 of 3 with the identical stray-GET signature** — pre-existing, not this sprint (S5 touches no setools code). Logs `gates/final-0b48604e/alone_*.txt`.

**Evaluator** (Agent `sonnet`, 52 min, 121,547 tok; resolver `agent-tool sonnet declared:alias-pin`; its own worktree `.wt-world-iter404-row243-v3s5-eval`; waited for `ALL_DONE.0b48604e` before running anything). **PASS 91/100**, 0 blocking (correctness 28, tests 25, security 19, standards 9, scope 10).
- 51 ledger lines, 47 own ids not in the plan's 124 (`comm -12` empty). Must-die `CTL-DIE` (the dispatch with no decision ref) killed by three named tests; must-live `CTL-LIVE-COMMENT` survived. Intermediate commits `34e97cd9` (flag absent) and `ba783c32` (flag present) checked green in a throwaway worktree.
- Re-check (a) driven by raw-SQL rows with another digest, challenge, expiry, namespace, an unenrolled key, a garbage artifact, a revoked approver and a mismatched `approved_by` — each refused "does not verify". S5-D4, S5-D8, S5-D9, S5-D10 and F3 judged sound (`-race -count=3` over the wait tests green, 15.2 s); S5-D15 (the planner's no-static mutant suite) judged weaker than whole-package — the judge re-ran its own survivors whole.
- **Non-blocking (logged, not filed — D-WORLD-109), carried to S6's planner:** (1) `E-LK-VR-LEN` — `ApproveTarget`'s `vr-` arm is not pinned to one argument (`approve vr-… --deny` would navigate; navigation only); (2) `E-JS-INCLUDES` / `E-JS-NODASH` — the `app.js` allowlist's `startsWith` and its trailing `-` are unpinned (widening to `/surface/confirm/` is killed); (3) `E-RC-INPUT-NOW` — `RecheckVerb`'s `Now` is near-equivalent (enrolment-time checks only); (4) `E-AWAIT-NOW-STALE` (a wait that never advances) is killed only by the package timeout, not a named test.

**Gate 3b:** PR [#599](https://github.com/sunholo-data/ailang-world/pull/599) head `0b48604e`: **3/3 success on attempt 1** (`ci_pr599_head0b48604e_s2.verdict`, the second 540 s segment). `mergeable` MERGEABLE/CLEAN on the head. Merged with `--match-head-commit` → **`48b8b6b6`**, tree byte-identical to the head. Merge push CI: **3/3 success on attempt 1** (`ci_merge_48b8b6b6_s2.verdict`, second segment).

**Routing evidence:** base=48b8b6b6e1270432d7d53d198c559d8903c56987@Gate-4 (record worktree cut from it). The Gate-1 base was 18c3d13231b39c7e2b0b437a964bce660d6739ec@2026-11-18T01:39:54Z, and the sprint worktree was cut from it (no drift at Gate 3).
- Controller: `claude-opus-5-5` (tok: not reported).
- Designer: not spawned (doc reviewed at 399).
- Planner: `opus` via Agent (502,955 tok). Driver note: codex `gpt-6.2-sol` unusable at preflight (probe rc 75, quota unreadable) → opus.
- Executor: `claude:claude-haiku-5-5` env pin, spawned as Agent `haiku` (73,626 tok); resolver `recipe claude:claude-haiku-5-5 declared:provider-pin`. Driver note: codex unusable → haiku.
- Evaluator: `sonnet` via Agent (121,547 tok), `agent-tool sonnet declared:alias-pin`. **Spawned as required; no role failed to spawn.**
- Quorum: none. Metered total: $0.
- Judge independence: cross-model (opus plan, haiku execution, sonnet judge).

**Progress:** 1.0 clauses 1–7 MET; goal unmoved. Row 243 V3 is at 5 of 6 sprints. From a terminal, `episode finish --confirm-on-surface` now waits for a key on the World card instead of a typed `y`, and verifies that confirmation itself before it opens the pull request. Sprint 6 adds the *Needs you* link and the live-browser end-to-end runs.

**Ruled out:**
- **Attributing `TestSeToolsRunBenchmarkTasksEndToEnd` to this sprint:** it fails alone at the base with the same stray-GET signature (1 of 3), and S5 touches no setools code.
- **Taking judge-403 notes 2–4 in this PR:** the planner measured a 2-line margin under the cap; they stay in S5's Appendix C for S6.
- **Fixing the judge-404 survivors in this PR:** navigation-only / near-equivalent, non-blocking; S6 touches the same bar and page code.

**Retro:**
- **Held, instance 3 of "every time comparison's operator is its own mutant" and instance 2 of "every argument of a write call is its own mutant":** the planner enumerated both; the judge's own argument mutants on the put were all killed.
- **New, instance 1: "a planner's mutant 'package' leg that excludes static tests is not the whole package"** (S5-D15). The planner's 0-survivor claim held for its 118 mutants, but the judge found 4 whole-package survivors among its own 47, all in the small M13 surface (bar arity, allowlist shape). For S6's planner: mutants on allowlists/grammars need a row per character class of the boundary (arity, prefix vs substring, trailing delimiter), and SURVIVED/KILLED is read from the whole package.
- **New, instance 1: the flaky `TestSeToolsRunBenchmarkTasksEndToEnd`** (a stray empty `GET` reaching its mock) now reds 3 of 3 loaded head runs and 1 of 3 alone runs at base and head. Not filed (D-WORLD-109's hygiene cap); added to the known-flaky list in the directives. If it reds CI it becomes a row.
- **Shape that worked (fifth time):** planner opus → controller round trip → executor haiku → commit_ms per snapshot → gates in background concurrent with a judge that waits on `ALL_DONE.<sha>`. About 3 h 5 min from Gate 1 to merge (the planner 90 min of it).

**Gate 5 bookkeeping:** rotation not due (#574 opened 2026-11-16, after the Monday 07:00 boundary; next rotation Monday 2026-11-23). The weekly external-issue sweep already ran this week at 393.

## 405 — 2026-11-18 — row 243 V3 sprint 6 LANDED, ROW 243 V3 COMPLETE (PR #601, merge `0077a7b0`): *Needs you* lists the finish confirmations (`app.js`, the only production change), CP7 over the real handler's golden, XC2 and XC3 in live Chrome, judge-404's survivors pinned, `TestCLICredentialWriterSites`; 425 lines; judge PASS 91 [PRODUCT]

**Pick:** row 243 V3 sprint 6, the planner step named as Next in record 404. That was also the queue head.
- **Re-measured at pick** (`origin/dev` `539cecd0`): no open PR from this mission and no stale worktree for row 243 V3 S6 (`git log origin/dev --grep "V3 S6"` empty; the three `iter404` worktrees are clean merged leftovers); `mission_decisions.sh --open` → no OPEN row; `--check` → ledger valid, 105 rows. The design was quorum-reviewed at 399, so no pick-time quorum was owed.
- **Clause map:** 1.0 clauses 1–7 MET, none UNMET. Release rows 150/151 still wait on attended steps (no directive or attended commit since 404). Row 243 serves D-WORLD-94 and moves no clause.
- **The previous fire:** 404's slot verdict reads `COMPLETED` (05:26Z, rc 0, 13,726 s) — `stamp complete` ran last, so the two-in-a-row `REAPED_at=gate-5` pattern (402, 403) did not recur.

**Gate 0:**
- Stamped first. `gh` is `sunholo-voight-kampff`. Billing CLEAN. Kill switch absent.
- Directives: 0 allowlisted on #574 (14 comments). The #574 watermark file is absent, so every comment was read; nothing to `--mark`.
- Self-notices: `watermark channel unavailable` (the #574 and #524 files are absent) — an instrument failure, not a verdict; same gap as 371–404.
- Inbox: no unread message addressed to `mission-world` (the newest is 404's own controlplane fyi).
- Skill drift: none. `SKILL.md` and all 13 `resources/*.md` at the resolved symlink target `cmp` equal to the driver clone's `origin/dev`; the pin is still `022f68f1a`, the same as 404's.
- Dev CI at `539cecd0`: `go host build + test gate` was pending at Gate 1 (the record merge was 13 min old); it read 3/3 success later in the fire.

**Carried from 404 (unlogged there; logged here as 404's retro asked):** 403's record merge `18c3d132` race job went red on attempt 1 on `host/verifygate` `TestReleaseChangeNotice` (ArmLocalPin got one `UNRECOGNISED RELEASE` notice for `v0.41.0` while printing `v0.41.0` in its own expected list — self-contradicting; the first failure of this test). Attempt 2 was green. Candidate mechanism, UNCONFIRMED: `verify_ail.sh:141` `printf | grep -qxF` under `set -o pipefail` (v0.41.0 is line 1 of the list, so `printf` can take SIGPIPE → rc 141) — 0/3000 locally on macOS, Linux untested. Log `~/.ailang/state/world-iter403/ci_merge18c3_race.txt`. No recurrence on 404's four CI runs nor on this fire's four (PR #601 head + merge, base `539cecd0`, all attempt 1). Instance 1 of 3 (D-WORLD-109).

**Planner** (Agent `opus`, 56 min, 371,701 tok; resolver `agent-tool opus fail-closed:env-pin`; `derive-planner-lane` → `opus fail-closed:env-pin`). It wrote `w-world-address-bar-verbs-v3s6-sprint-plan.md` and its JSON, committed as `2e2262dc`.
- **No cut (S6-D1):** design ≈ 355 raw → 425 compiled (1.20×). Every carried item taken (order 1 → 3 → 2).
- Sub-milestones: M15 157 (*Needs you* + CP7 + the handler-written golden) · M16 290 (XC2, XC3) · M17 425 (the carried items).
- **Mutants:** 49 run, whole-package verdicts with static tests included (the 404 retro rule): 47 killed (one, `W-SITE-POLL`, by the writer-site inventory only, as expected for an inventory), 2 named equivalents (`esc()` on an id the regex already gated), 0 survivors, 2 controls. `E-LK-VR-ARGS1` survived its first draft and was killed by adding `--deny vr-…` to `TestApproveTarget`. Ledger `plan-scratch/mutants/ledger.txt`.
- **§1.4 premises found false:** F1 the design's refused-row rule did not cover a summarised row whose id is not a confirmation id (now listed refused, no link); F2 no harness read a daemon-made `requests.json` (CP7 uses a golden, the verb-route pattern); F3 the Chrome driver's card helper looks for ids the confirm card lacks (XC2/XC3 use `openCard raw`, `waitFor`, `tap`); F4 `storeOwn`'s 2030 scope is refused by C10 against the wall clock (XC2 sets now + 900); F5 `E-AWAIT-NOW-STALE` is killed by S5's named test run alone — the package-timeout hang comes from S5's rig tests running first; F6 `confirm.js` already checks status 200 (the missing half was the test); F7 notes 2–3 land outside design §7; F8 XC2 checks the decision through the CLI's own read-back.
- **Base gates** `539cecd0`: `FAILS: none` (1,368 s). **Final gates** (planner's stage tree): `FAILS: none`.
- **Controller round trip:** byte-identical. I extracted the three patches from the plan, applied them to a fresh `git archive 539cecd0` and ran `diff -r` against the planner's `work` and `st-M17` trees: no difference. Numstat 425.

**Executor** (env `recipe claude:claude-haiku-5-5`, spawned as Agent `haiku` as in 374–404; one foreground run, 12 min, 52,484 tok).
- M15 → M16 → M17: every `--dry-run` clean; every part `ACCEPT_OK` (M17 `cmd1` 206 s, `cmd2` 43 s, `daemon1` 149 s, `daemon2` 191 s — the two big packages whole, split by test name). I re-grepped each log: 0 `--- FAIL`, 0 `^FAIL`; the only `--- SKIP` lines are the two expected counted skips inside the `skipped` step. Snapshots matched §7 (4/5/11 STATUS lines, DELETED 0, manifests verified). The executor's tree was byte-identical to my round-trip tree.
- Commits by `commit_ms.sh` (gofmt/build/vet + the part's acceptance before each; numstat filter widened to `-- cmd host docs app`; no refusal): M15 `a760bd9f` 157 · M16 `169c8070` 290 (commit-time test included the `chrome` part: XC2 `--- PASS` 1.26 s, XC3 1.39 s, plus XC1, the needs-you pair, XD1 and the bar's report test) · M17 `983abaf1` 425 — each equal to the plan.

**Controller gates** (plan §6 `gates.sh` on `983abaf1`, 1,150 s, concurrent with the judge's reading phase): touched packages derived from the diff = the plan's four; status, `verify_ail` (the base's line), gofmt, build, vet, fence, `pkgs` (four packages whole, 577 s), race, race3, node, inv (52 named passes), **chrome** (17 `--- PASS` incl. XC2 and XC3), **`verify_go`** (432 s) — **all rc 0, `FAILS: none`**; `DAEMON_ROUTES` 47 = base, `SURFACE_POSTS` 3, `UNTOUCHED_VIOLATIONS` none. First clean `verify_go` leg of V3 S5–S6 (404's was red only on untouched flakes).

**Evaluator** (Agent `sonnet`, 29 min, 87,662 tok; resolver `agent-tool sonnet declared:alias-pin`; its own worktree `.wt-world-iter405-row243-v3s6-eval`; waited for `ALL_DONE.983abaf1` before running anything). **PASS 91/100**, 0 blocking (correctness 28, tests 26, security 19, standards 9, scope 9).
- 44 ledger lines, 36 own ids not in the plan's table. Every kill by a behavioural test (node harness on rendered HTML, or a direct call), none by an inventory alone. Judge-404's survivors `E-LK-VR-LEN`, `E-JS-INCLUDES`, `E-JS-NODASH` and `E-RC-INPUT-NOW` (190 s whole `cmd`) reproduced SURVIVING at the base `539cecd0` and KILLED at the head. `TestCLICredentialWriterSites` non-vacuous both ways. Corrupting the digest in the page's 86-byte message fails XC2 and XC3. The CP7 golden regenerated at the head in a throwaway worktree came out byte-identical (it is the real handler's output). XC2/XC3 skip counted without `WORLD_CHROME` and fail under `WORLD_REQUIRE_PASSKEY_BROWSER=1`.
- **Non-blocking (logged, not filed — D-WORLD-109):** (1) `J-JS-RE-NL` — adding the `m` flag to `confirmLineHTML`'s `/^vr-[0-9a-f]{16}$/` survives the whole `host/surfacepage` package (no test feeds an id with an embedded newline). **Controller-reproduced:** mutant applied in the eval worktree, `go test -count=1 ./host/surfacepage/` rc 0 (`world-iter405/ctl_mut_re_nl.txt`), restored, tree clean. Production is correct today; a test gap only. (2) `J-JS-CS-NOGUARD` — dropping `Array.isArray` on `confirmations` survives (no non-array fixture; the guard is right by inspection). (3) `J-JS-ALLOW-DECIDE` — the pre-S6 `/surface/decide/sha256:` allowlist prefix is unpinned against `includes` (the same class as `E-JS-INCLUDES`, on the older line). (4) XC3 is a happy-path proof; a removed digest re-check survives the two Chrome tests narrowly but is killed by `TestPollVerbRequestRechecks`. (5) The judge did not re-run M15/M16 alone (the controller's commit-time acceptance ran each) nor `E-AWAIT-NOW-STALE` at the base.

**Gate 3b:** PR [#601](https://github.com/sunholo-data/ailang-world/pull/601) head `983abaf1`: **3/3 success on attempt 1** (`ci_pr601_head983abaf1_s2.verdict`, the second 540 s segment). `mergeable` MERGEABLE/CLEAN. Merged with `--match-head-commit` → **`0077a7b0`**, tree byte-identical to the head. Merge push CI: **3/3 success on attempt 1** (`ci_merge_0077a7b0_s2.verdict`, second segment).

**Routing evidence:** base=0077a7b03321dd333a570e20342a4581dd4dc6c3@Gate-4 (record worktree cut from it). The Gate-1 base was 539cecd0d37930fb7ab0021bc4e529d794c653b2@2026-11-18T05:39:44Z, and the sprint worktree was cut from it (no drift at Gate 3).
- Controller: `claude-opus-5-5` (tok: not reported).
- Designer: not spawned (doc reviewed at 399).
- Planner: `opus` via Agent (371,701 tok). Driver note: codex `gpt-6.2-sol` unusable at preflight (probe rc 75, quota unreadable) → opus.
- Executor: `claude:claude-haiku-5-5` env pin, spawned as Agent `haiku` (52,484 tok); resolver `recipe claude:claude-haiku-5-5 declared:provider-pin`. Driver note: codex unusable → haiku.
- Evaluator: `sonnet` via Agent (87,662 tok), `agent-tool sonnet declared:alias-pin`. **Spawned as required; no role failed to spawn.**
- Quorum: none. Metered total: $0.
- Judge independence: cross-model (opus plan, haiku execution, sonnet judge).

**Progress:** 1.0 clauses 1–7 MET; goal unmoved. **Row 243 V3 is complete (6 of 6 sprints, S1–S6 at iterations 400–405).** The World surface now shows a pending `episode finish --confirm-on-surface` request under *Needs you* with a link to its card, and a live-Chrome test drives the whole path: the real CLI waits, a passkey confirms on the card, and the CLI opens the pull request without reading the terminal. What stays attended: Mark's own no-terminal exit run with his real key (parent §10 step 4; QUICKSTART §14–§20 pending their first verbatim run). V4 (design §0.4) inherits who runs the finish with no terminal (F1 is kept), launching from the bar, `episode start` with a launch preset, the herdr hand-off, per-check confirmation of unsandboxed checks (RV3), and whether a pending confirmation should notify.

**Ruled out:**
- **Filing the judge's new survivors as rows:** non-blocking test gaps on correct production code; D-WORLD-109 logs them here. They are carried as V4 planner inputs.
- **Moving the V3 design to `implemented/`:** the parent row 243 is not done (V4 remains), and the V2 design stayed in `planned/` on the same grounds.
- **Attributing the 403 `TestReleaseChangeNotice` red to S5/S6:** neither touches `host/verifygate` or `verify_ail.sh`, and eight CI runs since are green on attempt 1.

**Retro:**
- **Held, instance 1 → 2 of "grammars and allowlists get a mutant per boundary class, and SURVIVED/KILLED is read from the whole package" (404's rule):** the planner ran every mutant whole-package with statics, and its boundary rows (arity, prefix vs substring, trailing dash, case, 15/17 hex) caught `E-LK-VR-ARGS1` in its own draft. The judge's three new survivors are boundary classes the planner did not enumerate (`m` flag on an anchored regex, non-array input, an older allowlist line) — the rule is right but its class list is not closed; add "regex flags" and "input type (non-array/non-string)" to the directive.
- **Held, sixth time:** planner opus → controller round trip → executor haiku → commit_ms per snapshot → gates in background concurrent with a judge that waits on `ALL_DONE.<sha>`. About 2 h 10 min from Gate 1 to merge (planner 56 min).
- **The `stamp complete`-last fix held:** 404's slot read `COMPLETED`.

**Gate 5 bookkeeping:** rotation not due (#574 opened 2026-11-16, after the Monday 07:00 boundary; next rotation Monday 2026-11-23). The weekly external-issue sweep already ran this week at 393.

## 406 — 2026-11-18 — row 243 V4 design DONE (`w-world-address-bar-verbs-v4.md` r1): start and finish an episode from the World app through an opt-in user LaunchAgent runner, a signed operator-action law and per-effect receipts; create-and-launch before finish; 12 sprints, S1–S2 answer-independent; quorum r1 blocked 1/3, r2 PROCEED 3/3; independent audit PASS 81; new asks D-WORLD-119/120 [PRODUCT]

**Pick:** row 243 V4 design r0, the Next named in record 405 and the queue head.
- **Re-measured at pick** (`origin/dev` `41331b4d`): no V4 doc, branch or worktree (`git log origin/dev --grep V4` → none for row 243; `ls design_docs/planned | grep v4` empty; the `iter404`/`iter405` worktrees are clean merged leftovers). `mission_decisions.sh --check` → valid, 105 rows; `--open` → none.
- **Clause map:** 1.0 clauses 1–7 MET, none UNMET. Release rows 150/151 still wait on attended steps. Row 243 serves D-WORLD-94 (D-WORLD-109 = A) and moves no clause.
- **Attended context found at pick (not a directive):** two untracked briefs Mark wrote this morning in the main checkout, `design_docs/planned/w-phone-episode-workflow.md` and `w-human-agent-surface.md` (~09:15Z, after the first attended Sol episode through World), plus that episode's PR #604 (`world/mobile-needs-you`, `app.js`/`app.css`/`index.html`). The phone brief says "Row 243 V4 owns the runner design" and asks for create-and-launch first, a runner choice backed by a rig experiment, and one headless Sol from a browser. I copied both briefs to `/tmp/w406-designer/inputs/` (sha256 pinned in the doc's W31) and handed them to the designer as UNRATIFIED product context, never as authority. I did not touch, commit or move them. Claim sent to `controlplane` (`inbox_1794995082755_a8143879`). #604 merged attended at 09:48:30Z mid-fire (`65108353`), so the record is cut from that base.

**Gate 0:**
- Stamped first. `gh` is `sunholo-voight-kampff`. Billing CLEAN. Kill switch absent.
- Directives: 0 allowlisted on #574 (15 comments) and 0 on #524 (47), read against both the absent issue-scoped watermark and the LEGACY `mission-world-last-seen` (`2026-10-27T09:03:42Z`). Nothing to `--mark`.
- Self-notices: rc 2 `watermark channel unavailable` (the #574 and #524 files are absent) — an instrument failure, not a verdict; same gap as 371–405.
- Inbox: the only new `mission-world` traffic is this fire's own driver lane-degraded note.
- Skill drift: none — `SKILL.md` and all 13 `resources/*.md` at the resolved symlink target and the pin `cmp` equal to the driver clone's `origin/dev` (pin `9828cda4a`).
- Dev CI at `41331b4d`: 3/3 check-runs, none non-green.

**Designer** (rotation: last used `claude:claude-opus-5-5` → next **`codex:gpt-6.2-sol`**; resolver said `recipe claude:claude-opus-5-5 declared:provider-pin`, the env seed). The driver had marked codex unusable at preflight (probe rc 1, planner/executor handed to opus/haiku). My own bounded probe from `/tmp` failed rc 1 on `Not inside a trusted directory` (my instrument, not the lane); re-run from the repo it returned `ok`, rc 0, and `MISSION_OVER_RATION` was empty. So I followed the rotation (399 skipped codex only because it was over ration), spawned through the codex recipe (`codex exec --sandbox workspace-write`, `< /dev/null`, directive delivery asserted, 3,000 s cap, rc file), not the Agent tool — the routing table gives a `codex:` designer to the recipe. Worktree `.wt-world-iter406-row243-v4`, branch `design/world-iter406-row243-v4`, base `41331b4d`. The directive carried the iter-399 designer rules (W-row per claim with control, file:line or NEW per symbol, every input and failure class traced, threat table, one mutant per branch and per boundary class incl. regex flags and input type with whole-package verdicts — the 405 retro), the judge-405 notes, the #604 collision, and the two sandbox rules (sandbox-denied probes are UNINFORMATIVE; guards outside the worktree are observations under a synthetic HOME).
- **r0** `5d5a50d6` (29 min, 242,512 tok, 805 lines, 32 W-rows, five named planner experiments E1–E5).
- **r1** `c7810dac` (one revision run, 11 min, 67,485 tok, 880 lines, 42 `[r1]` marks) — the Fable-diet unit is one doc plus one protocol revision; this was that revision.

**What the doc decides:** Q4 = **B**, a separate opt-in user LaunchAgent (`service runner …`) running as the operator with an explicit PATH/HOME and non-interactive Git/GitHub, polling a NEW credentials **schema 10** `operator_actions` table (+ `operator_runner_state`, `operator_action_steps`); a NEW signed `world/operator-action/v1` law reusing the V3 judges; single claim by CAS + kernel lock; a per-effect write-ahead receipt and reconcile rule for worktree, mint, files, launch, checks, commit, push, label, PR and record (§3.3); the bar gains a `review` answer, a guarded inert POST and a card; **create-and-launch first** (first browser slice S8: one headless start), then finish through a private verified F1 entry with two key gestures (local prep, then V3's real PR prompt); RV3 and notifications out; merge out. **12 sprints ≈ 9,420 raw / 12,719 calibrated** (13,664 with the optional herdr S12) against the parent's ≈ 2,000 — the doc attributes the growth to admission, ownership, recovery and their test matrix. S1 (judge-405 pins + the pure admitted grammar) and S2 (the action law against fixtures) need no answer.

**Quorum** (`--author codex:gpt-6.2-sol`; OpenAI benched; seats claude-sonnet-5@claude-p, gemini-4-argon, oc-glm-5-3; reserve oc-kimi-k3; `absent_reviewers` = `[]` in both rounds):
- **r1 on r0: BLOCKED 1/3** ($0.229). sonnet REJECT: W07's `gh auth status` rc 1 is a negative signal on the finish half's credential premise, and E2 as scripted may not exercise the Keychain helper. gemini PASS (pin the /tmp briefs; measure #604's hunks). glm PASS (bound the terminal side of the per-episode lock; three uncited anchors; §4.3 numbering 5,7,6).
- **Rule 3f — I measured the objection instead of forwarding it.** W07 ran inside the designer's sandbox (restricted egress) — UNINFORMATIVE. I ran §11's E2 script verbatim outside any sandbox from the process tree of the GUI LaunchAgent `dev.ailang.mission-world` (probe ← claude ← `mission-control.sh` PID 49385 ← launchd 1; `launchctl print gui/501` lists that PID), under `env -i`: installed-job PATH `/usr/bin:/bin:/usr/sbin:/sbin` → **rc 2** (`gh` missing; confirms the PATH finding); Homebrew PATH → `ADMIN`, `rig-auth-read-ok`, **rc 0** (`gh auth status`, `gh label list`, `viewerPermission`, and `git -c credential.interactive=never ls-remote` over the HTTPS origin via `osxkeychain`). Scope: read-level, one rig, one remote, a launchd-descended process that is not the runner's own plist (E3 stays the plist test). I also measured #604 (head `b9791084`, merge-base `41331b4d`; `app.js` hunks clear of S1's pin sites `:434`, `:456–458`, `:507`). Both went to the designer as measurements (W33, W34), with every reviewer fix.
- **r2 on r1: PROCEED 3/3** ($0.231). Non-blocking notes, carried to the planner: sonnet — §2.1 should say W33 ran in the driver's tree, not the runner's; gemini — "non-blocking with a ≤5 s bound" should read `LOCK_NB` with the 5 s bounding only the store read for the refusal message; glm — pin `runPresetEpisodeStart`, `launchEpisode`, `pollVerbRequest`, `RecheckVerb`, `systemServiceController`, `serviceEnv.confirm` by name. **Controller-measured:** all six exist as named (`episode_preset.go:53`, `episode_launch.go:258`, `episode_finish_surface.go:214`, `host/authority/judge_verb.go:107`, `service_launchctl.go:24`, `service.go:167`). No carve-out was needed or made.

**Evaluator** (Agent `sonnet`, foreground, 3.7 min, 135,013 tok; resolver `agent-tool sonnet declared:alias-pin`; directive `~/.ailang/state/world-iter406/evaluator-directive.md`): **PASS 81/100**, 0 blocking (premises 26/30, revision discipline 13/15, ask safety 15/20, soundness 20/25, plannability 7/10). It re-ran 28 W-rows (all reproduce except W31/W34, which drifted in time) and three uncited citations; every r1 hunk traced to a reviewer or to W33/W34, nothing unrequested added. Non-blocking, ranked: (1) an `uncertain` action has no attended resolver verb, and the terminal refusal would then block Mark's own escape hatch; (2) a headless child outliving the runner is unowned (`agentrun.Run` waits on its own child; the launch receipt is appended by the launcher after exit) and E3's fake child never outlives its parent; (3) S3's schema 10 carries finish structures while only D-WORLD-119 gates it — a one-way migration; (4) the second long-lived credentials writer vs the ratified non-waiting writer lock and the pinned writer-site inventories is unaddressed; (5) the asks omit decision-relevant consequences (schema 10 one-way; the runner holds ADMIN-level GitHub credentials, so "PR only" is code, not credential; herdr rides on unmeasured E4); (6) the proof types must live in `host/authority` with unexported constructors or `package main` can forge them; (7) W31's human-agent brief changed after r1 (now `3b0f4535…`, 74 lines, edited ~09:57Z — Mark is still writing it); (8) #604 is merged; (9) the herdr plugin hard-codes `codex=gpt-6.1-sol` (`world-config-init:96`); (10) the three r2 notes; (11) P1 is not really a false premise. **Recommendation: split at the plan level** — plan S1–S2 now; S3–S8 as "V4a start/launch" after D-WORLD-119 with schema 10 start-only; S9–S11 as a "V4b finish" addendum after D-WORLD-120 absorbing (1), (2), (4), (6); S12 as a herdr child after E4.
- **Controller-reproduced before acting:** (8) `gh pr view 604` → MERGED `2026-11-18T09:48:30Z`, merge `65108353`; (4) `host/store/writer_lock_unix.go:41` is `Flock(LOCK_EX|LOCK_NB)` returning `WriterAlreadyActive` (`:34`), and the doc names the writer lock only in one §5.2 row (`:212`); (7) `shasum` of the main-checkout brief → `3b0f45352a25…`, 74 lines.
- **Disposition:** not a third designer run (the diet is spent and nothing blocks). (5) is folded into the filed D-WORLD-119/120 rows as a controller filing note in their evidence cells, quoting the audit; the asks' A/B/C and defaults are the designer's text unchanged. (1)–(4), (6), (9), (10) and the split are V4 planner inputs (D-WORLD-109: logged here, not filed as rows).

**Ledger:** filed **D-WORLD-119** (surface start/launch through the runner; recommend A; default: S1–S2 only, fixtures) and **D-WORLD-120** (no-terminal finish, two gestures; recommend A; default: F1 kept, no runner finish). `mission_decisions.sh --check` → valid, 107 rows; `--open` → D-WORLD-119, D-WORLD-120.

**Gate 3b:** the design and this record land in one docs PR (as 394/396/399) — see the PR below.

**Routing evidence:** base=65108353df0ebd2597a81f75372dfdd9eeb4c543@2026-11-18T10:30:58Z (Gate 4; the record worktree is cut from it). The Gate-1 base was 41331b4d880e4ea1a4c19a2f1e65a74204580bc4@2026-11-18T09:40:04Z and the design worktree was cut from it; `drift gate1` → DRIFT (2 commits, #604's attended merge), benign — the design commits cherry-picked onto `65108353` with no conflict and #604 touched no file the doc measures except as W34 records.
- Controller: `claude-opus-5-5` (tok: not reported).
- Designer: `codex:gpt-6.2-sol` via the codex recipe — r0 242,512 tok, r1 67,485 tok. Rotation pointer advanced to `codex:gpt-6.2-sol`. Resolver answer `recipe claude:claude-opus-5-5 declared:provider-pin` (the env seed) recorded beside it; the rotation is owned by the skill, and the driver comment says so.
- Planner: not spawned (design fire). Executor: not spawned.
- Evaluator: `sonnet` via Agent (135,013 tok), `agent-tool sonnet declared:alias-pin`. **Spawned as required; no role failed to spawn.**
- Quorum: claude-sonnet-5@claude-p ($0, subscription), gemini-4-argon, oc-glm-5-3; r1 $0.229 + r2 $0.231. **Metered total $0.460.**
- Judge independence: cross-vendor (OpenAI designer, Anthropic judge; OpenAI benched from the quorum).

**Progress:** 1.0 clauses 1–7 MET; goal unmoved. Row 243 V4 is designed: the route to Mark's no-terminal exit run (start → launch → report → finish → PR from the World app) now has a measured runner choice and a 12-sprint plan, of which S1–S2 are loop-routable now and S3–S12 wait for D-WORLD-119/120.

**Ruled out:**
- **Following the resolver's `claude:claude-opus-5-5` for the designer:** that is the env seed; the rotation's next entry was codex, which probed rc 0 and was in ration.
- **Treating the driver's preflight "codex unusable (rc 1)" as current:** re-probed rc 0 from the repo; my own first probe's rc 1 was the trusted-directory check, an instrument fault.
- **Forwarding sonnet's auth objection to the designer as a claim:** measured it (E2 rc 0 with the Homebrew PATH) and handed over the measurement.
- **A controller carve-out over r1 for the r2 notes:** r2 proceeded; the notes are planner inputs, and glm's symbol list was measured instead.
- **Treating Mark's untracked briefs as directives or committing them:** they are attended drafts in progress (one changed again during this fire); they shaped order and scope only, as the doc says line by line.

**Retro:**
- **Instance 2 of "a sandboxed designer's negative probe reaches the quorum as a negative premise"** (instance 1: iteration 399's designer ran no runtime probe at all and said so). Here W07 was labelled UNVERIFIED in the body but read as evidence by a reviewer; the controller's out-of-sandbox measurement cleared it in one round. Process note for the next codex designer directive: any probe touching the network, launchd or `ps` is the CONTROLLER's to run before quorum, not the designer's — run E-experiments that are read-only before the first quorum round.
- **Held:** codex designer through the recipe, one revision, one re-quorum, independent sonnet audit, all inside ~1 h 40 min from Gate 1.
- **Watch:** the doc is 12 sprints. The audit's plan-level split is the planner's first decision for V4 S1.

**Gate 5 bookkeeping:** rotation not due (#574 opened 2026-11-16 09:06Z, after the Monday 07:00 local boundary; next Monday 2026-11-23). The weekly external-issue sweep already ran this week at 393.

## 407 — 2026-11-18 — row 243 V4 sprint 1 LANDED (PR #606, merge `d0f22aa7`): the judge-405 `app.js` pins and `host/verbs.AdmitOperator`, the pure, unregistered admitted operator subset joined to the CLI's real flag sets; D-WORLD-119/120 a policy input whose zero value admits nothing; 976 lines; judge PASS 93 [PRODUCT]

**Pick:** row 243 V4 sprint 1, the Next named in record 406 and the queue head.
- **Re-measured at pick** (`origin/dev` `d92983d5`): no V4 plan, sprint branch or worktree (`git log origin/dev --grep V4` → only the design commits; `gh pr list --state open` → none). Quorum reader rc 0 for `w-world-address-bar-verbs-v4`: r1 `blocked` (sonnet reject), r2 `proceed` 3/3, `absent_reviewers` `[]` in both. `mission_decisions.sh --check` → valid, 107 rows; `--open` → D-WORLD-119, D-WORLD-120. S1 needs neither answer (design §8 and the D-WORLD-119 default).
- **Clause map:** 1.0 clauses 1–7 MET, none UNMET. Release rows 150/151 still wait on attended steps. Row 243 serves D-WORLD-94 (D-WORLD-109 = A) and moves no clause.
- Claim sent to `controlplane` before routing (Mark was attended this morning).

**Gate 0:**
- Stamped first. `gh` is `sunholo-voight-kampff`. Billing CLEAN. Kill switch absent (`mission-world.disabled`).
- Directives: 0 allowlisted on #574 (16 comments), read from the epoch because the issue-scoped watermark is absent; the LEGACY `mission-world-last-seen` (`2026-10-27`) is older, so nothing was missed. Nothing to `--mark`.
- Self-notices: rc 2 `watermark channel unavailable` (the #574 and #524 files are absent) — an instrument failure, not a verdict; same gap as 371–406.
- Inbox: no unread message addressed to `mission-world`; the only `mission-world` traffic is this mission's own outbound reports.
- Skill drift: none — `SKILL.md` and all 13 `resources/*.md` at the resolved symlink target `cmp` equal to the driver clone's `origin/dev` (pin `55d336df`, drift 0).
- Dev CI at `d92983d5`: 3/3 check-runs, none non-green.

**Planner** (Agent `opus`, foreground, 59 min, 364,238 tok; resolver `agent-tool opus fail-closed:env-pin`, `derive-planner-lane.sh` → `opus fail-closed:env-pin`; directive `~/.ailang/state/world-iter407/planner-directive.md`, the 405 directive adapted, with the 404/405 boundary-class rule extended to every grammar the parser adds). Plan `w-world-address-bar-verbs-v4s1-sprint-plan.md` + JSON, committed `b672f64e`. **No cut:** 976 of 1,400 lines (design 620 raw → 1.57×; pins 127 for 210 raw, the subset 849 for 410 — the boundary-class matrix). Four sub-milestones: M1 the three pins (127), M2 `host/verbs/operator.go` + six grammar rows (570), M3 policy/never/identifier/other-verb rows (778), M4 the CLI join (976). 153 mutants with whole-package verdicts: 148 killed, 5 named equivalents (three regex `s`-flag mutants, `EQ-NOTADMIT-GT0`, `SEAM-AFTER-PARSE`), 0 survived; judge-405's three survivors reproduced SURVIVED on the base first.
- **Premises found false (§1.4):** F2 — there was nothing to join to: the CLI's `episode start`/`finish` flag sets are locals of two long functions and `--help` is static text, so M4 adds a nil-by-default 10-line seam `episodeFlagSetSeen` (the one production change outside `operator.go`; precedent `finishAfterStep`); F3 — `TestSurfaceVerbIsReadOnly` refuses `strconv` in host/verbs production (the first draft went red), so bool words are a table checked against `strconv.ParseBool` in tests; F4 — the "1–64 ASCII" episode grammar is `^[a-z0-9][a-z0-9-]{0,63}$`; F5 — a wrong-case verb never reaches the subset (`Lookup` answers it); F6 — `--headless` without `--launch` is the CLI's own refusal (usage); F7 — the CLI accepts a repeated flag (last wins), the subset refuses it; F8 — a `null` element in `confirmations`/`decisions`/`pending` throws in `app.js` (the daemon cannot send one; carried to S8); F10 — `Lookup` answers `world-publish …` as an address, the subset classes it `never`.
- **Controller round trip (independent):** extracted all four patches from the plan file, applied them with `patch -p1 -E` to a fresh `git archive d92983d5`, `diff -r` against the planner's `st-M4` → byte-identical; a one-line perturbation control made it differ (rc 1) and restoring it matched again.

**Executor** (Agent `haiku`, foreground, 5.7 min, 64,665 tok; resolver `recipe claude:claude-haiku-5-5 declared:provider-pin` — spawned through the Agent tool's `haiku` pin as in 402–406, no hook denial; directive `~/.ailang/state/world-iter407/executor-directive.md`). M1→M4 in one run, every acceptance part `ACCEPT_OK` (M4 incl. `cmd/ailang-worldd` whole in three parts, 151/31/42 s), snapshots 1/0/1, 3/0/3, 4/0/4, 7/0/7. Controller checks: each snapshot `MANIFEST` `shasum -c` ok; the executor's final tree `cmp`-equal to the planner's stage tree on all 7 files; main checkout untouched (only Mark's untracked drafts). `commit_ms.sh` per snapshot, 0 refusals, numstat 127 / 570 / 778 / 976 = plan: `9fd474a6`, `6f502e92`, `e563bca9`, `407b5ace`.

**Gates** (plan §6 `gates.sh` on `407b5ace`, background, 1,230 s): status, verify_ail, gofmt, build, vet, fence, pkgs, race, node, inv (54 named PASS), chrome (5 `--- PASS`), **verify_go** — all rc 0, `FAILS: none`. `TOUCHED_PACKAGES` = the plan's three; `DAEMON_ROUTES` 47 = base, `SURFACE_POSTS` 3; `UNTOUCHED_VIOLATIONS` none; `CMD_SEAM_LINES` `10 0`; `PARSER_PRODUCTION_CALLERS` none (control `verbs.AdmitWhy` 1); `app.js`, `surface_page.go`, `host/store/credentials.go` unchanged. `TestSeToolsRunBenchmarkTasksEndToEnd` passed in `pkgs` this time.

**Evaluator** (Agent `sonnet`, foreground, 30 min, 122,000 tok; resolver `agent-tool sonnet declared:alias-pin`; own worktree `.wt-world-iter407-row243-v4s1-eval`; directive asked open questions, no controller findings supplied): **PASS 93/100**, 0 blocking (correctness 28/30, tests 28/30, security 20/20, standards 9/10, scope 8/10). Its own ledger `~/.ailang/state/world-iter407/judge/mutants.tsv` (112 lines, 54 ids not in the plan's): every policy operand, each grammar's bounds/anchors/case, the bool table, the duplicate counter, the canonical argv, the mirror's membership and kinds, the seam call sites and CLI-side flag drift (add/rename both ways) all killed by named behavioural tests; a 60,000-line differential fuzz (deleted afterwards) admitted 142 lines, each parsed by the CLI's real flag set to the same episode and preset with no flag outside the five admitted. Judge-405's survivors reproduced SURVIVED at the base and KILLED at the head, whole package. Each commit M1–M4 green alone in a throwaway worktree. Non-blocking, ranked: (1) dropping `fs.SetOutput(io.Discard)` in `AdmitOperator` survives the whole `host/verbs` package — no test asserts the pure parser writes nothing, so once S7 calls it a malformed bar line would print flag usage to the daemon's stderr; (2) F8 still unpinned; (3) numeric CLI flags are mirrored as strings, so `--ttl x` is refused `terminal` rather than `usage` (both refuse); (4) it ran the cmd mutants narrow (`^TestOperator`) for time; (5) its ledger has ragged rows from multi-line edits.
- **Controller-reproduced before acting:** (1) in the eval worktree, `fs.SetOutput(io.Discard)` → `_ = io.Discard` (`git diff --stat` 1 file, `go vet` ok), `go test -count=1 ./host/verbs/` → `ok` rc 0 — the mutant survives; restored, tree clean. Disposition: logged here and carried to the S7 planner (the first production caller), not filed as a row (D-WORLD-109).

**Gate 3b:** PR #606 head `407b5ace` CI 3/3 green (second 9-minute poll segment). Merged at `d0f22aa7`. **`origin/dev` had moved while the sprint ran**: PR #607 (`ui/episode-composer-review`, attended — Mark's episode composer: `app.js` +51/−28, `composer.js` NEW, `index.html`, `app.css`, `host/daemon/daemon.go` +1 route, `surface_page.go`, `host/surfacepage/app.go`, `cmd/ailang-worldd/surface_draft_catalog.go`) merged at `390af570`, so the merge commit combines both. Because CI's node harness tests may skip, I re-ran on `d0f22aa7` locally with `WORLD_EXEC_NODE` set: `go build`, `go vet` of the four packages rc 0; `host/surfacepage` whole rc 0, 86 PASS / 0 SKIP / 0 FAIL incl. all three pins; `host/verbs` ok; `cmd/ailang-worldd -run '^TestOperator'` 3 PASS; host/daemon route/writer/dependency tests ok. The three pinned `app.js` sites are untouched by #607 (now `app.js:435`, `:458`, `:508`, one line down). Merge commit `d0f22aa7` CI **3/3 green** (second segment). **LANDED.**

**Routing evidence:** base=d92983d59f887905251df5b99e825c299136760f@2026-11-18T13:39:54Z (Gate 1; sprint worktree re-measured SAME at 13:41:23Z). Gate 4 base d0f22aa7b95c2c35f82b7ac4e8ebfce6cf36408b@2026-11-18T15:48:21Z (DRIFT from Gate 1: #607 attended + #606 this sprint; benign — re-checked as above).
- Controller: `claude-opus-5-5` (tok: not reported).
- Designer: not spawned (plan fire; the doc is reviewed). Rotation pointer unchanged at `codex:gpt-6.2-sol` → next design gets `pi:ollama/glm-5.3:cloud`.
- Planner: `opus` via Agent (364,238 tok). Executor: `haiku` via Agent for the `claude:claude-haiku-5-5` pin (64,665 tok). Evaluator: `sonnet` via Agent (122,000 tok). **Every required role spawned; none failed to spawn.**
- Judge independence: same vendor, different model and fresh context (opus planner, haiku executor, sonnet judge), as 400–405.
- Metered: **$0.00** (no quorum, no metered lane).

**Progress:** 1.0 clauses 1–7 MET; goal unmoved. Row 243 V4 is at 1 of 12 sprints: the bar now has a pure, tested decision of which operator lines the app may later offer to run, tied to the CLI's real flags, with every answer-dependent mode behind a policy that admits nothing until Mark answers.

**Ruled out:**
- **Treating the CI green on the merge as proof the pins still hold after #607:** CI may skip node tests; re-ran them locally with `WORLD_EXEC_NODE` on `d0f22aa7` (86 PASS, 0 SKIP).
- **Spawning the haiku executor through `claude-sub`:** the Agent tool's `haiku` pin is the same model, accepted by the spawn-pin hook in 402–407; recorded as before.
- **Fixing judge finding (1) in this sprint:** it is real (reproduced) but harmless until a production caller exists; S7 owns the first caller and gets it as a planner input.
- **A re-quorum at pick:** the doc's r2 artifact is PROCEED 3/3 with no absent seat; nothing changed since.

**Retro:**
- **Held, a 7th time:** planner opus → controller round trip → executor haiku → commit_ms per snapshot → background gates beside an Agent judge in its own worktree. Gate 1 → merge CI green ≈ 2 h 10.
- **New (instance 1): an attended merge to `dev` mid-sprint touched the same page file as the sprint's tests.** The merge was conflict-free and green, but CI could not prove the node pins against the new `app.js`. A local re-run of the touched packages on the merge commit closed it in under 2 minutes. Planner/controller note for S7/S8, which change `app.js`: re-measure against `origin/dev` immediately before the sprint PR merges, not only at pick.
- **Planner-directive rule carried:** for any pure function a later sprint will call from a server, a test that it writes nothing to stdout/stderr (judge-407 (1)).

**Next:** row 243 V4 sprint 2 planner — the canonical action envelope and the typed SSHSIG/passkey proof and recheck against fakes (design §8 S2; answer-independent). Carry S1's Appendix C (the `OperatorIntent` S2 canonicalizes; `OperatorActionIDPattern`; the policy stays an input) and judge-407 notes (1) and (3). S3–S8 wait for D-WORLD-119; S9–S11 for D-WORLD-120.

**Gate 5 bookkeeping:** rotation not due (#574 opened 2026-11-16 09:06Z; next Monday 2026-11-23). The weekly external-issue sweep already ran this week at 393.

## 408 — 2026-11-18 — row 243 V4 sprint 2 LANDED (PR #609, merge `cee5174c`): the `world/operator-action/v1` envelope, its strict decoder, two digests and the 97-byte signed message, and the key law `host/approval.ConfirmOperatorAction`/`CheckOperatorAction` over SSHSIG and passkey proofs; unregistered; the `host/authority` adapters (260 lines) parked as a compiled patch for S2b; 1,394 lines; judge PASS 92 [PRODUCT]

**Pick:** row 243 V4 sprint 2, the Next named in record 407 and the queue head.
- **Re-measured at pick** (`origin/dev` `8db22168`): no S2 plan, branch or worktree (`git log origin/dev --grep V4` → S1 and design commits only; `gh pr list --state open` → none; `git worktree list` → 406/407 trees only). Quorum reader rc 0 for `w-world-address-bar-verbs-v4`: r1 `blocked`, r2 `proceed`, `absent_reviewers` `[]` in both. `mission_decisions.sh --check` → valid, 107 rows; `--open` → D-WORLD-119, D-WORLD-120. S2 needs neither (design §8 and the D-WORLD-119 default).
- **Clause map:** 1.0 clauses 1–7 MET, none UNMET. Release rows 150/151 still wait on attended steps. Row 243 serves D-WORLD-94 (D-WORLD-109 = A) and moves no clause.
- Claim sent to `controlplane` before routing (`inbox_1795023668044_38a93001`).

**Gate 0:**
- Stamped first. `gh` is `sunholo-voight-kampff`. Billing CLEAN. Kill switch absent (`mission-world.disabled`).
- Directives: 0 allowlisted on #574 (17 comments) and #524 (47), each read from the epoch because the issue-scoped watermarks are absent. Nothing to `--mark`.
- Self-notices: rc 2 `watermark channel unavailable` (the #574 and #524 files are absent) — an instrument failure, not a verdict; same gap as 371–407.
- Inbox: no unread message addressed to `mission-world`; the `mission-world` traffic is this mission's own outbound reports.
- Skill drift: none — `SKILL.md` and all 13 `resources/*.md` at the resolved symlink target `cmp` equal to the driver clone's `origin/dev` (pin `8bfca4d03`, drift 0). My first loop printed DIFF for all 13: it built the `git show` path without the `resources/` prefix. Re-run with the path fixed: 14 SAME, and a nonexistent path gave DIFF as the negative control.
- Dev CI at `8db22168`: 3/3 check-runs green.
- The main checkout was 19 commits behind `origin/dev` with only untracked files (Mark's drafts and stale scratch); `git merge --ff-only origin/dev` brought it level, and the drafts were untouched.

**Planner** (Agent `opus`, foreground, 37 min, 371,145 tok; resolver `agent-tool opus fail-closed:env-pin`, `derive-planner-lane.sh` → `opus fail-closed:env-pin`; directive `~/.ailang/state/world-iter408/planner-directive.md`, the 407 directive adapted to S2, carrying judge-407's writes-nothing rule). Plan `w-world-address-bar-verbs-v4s2-sprint-plan.md` + JSON, committed `b8545962`. **Cut by measurement:** everything compiled came to 1,654 lines against the 1,400 cap, so the `host/authority` adapters (`JudgeOperatorAction`, `RecheckOperatorAction`, `authority.VerifiedOperatorAction` and its construction test; 260 lines) are parked, compiled, vetted and green on M4, at `~/.ailang/state/world-iter408/plan-scratch/patches/PARK-authority-judges.patch`. Landed 1,394 (design 840 raw → 1.66× for the full scope). Four sub-milestones: M1 the envelope and strict decoder (530), M2 limits/validation/digests/the bar join (773), M3 the key law (1,304), M4 writes-nothing (1,394). 150 mutants with whole-package verdicts: 146 killed, 4 named equivalents (`ID-DOTALL`, `D4-SHAPE-ERR`, `W15-ZERO-P`, `PASSKEYUSE-OK`), 0 survived; judge-407's `fs.SetOutput(io.Discard)` survivor reproduced SURVIVED on the base, then killed by `TestAdmitOperatorWritesNothing`.
- **Premises found false (§1.4):** F2 — the proof cannot live in `host/authority` as design §3.1 says, because the SSHSIG/WebAuthn primitives it needs are unexported in `host/approval` (V3's `SignedVerb` lives there for the same reason); only the store adapter can live in `host/authority`. F3 — `encoding/json` refuses none of duplicate keys, trailing data or invalid UTF-8, and matches keys case-insensitively, so the decoder walks the bytes first. F4 — the challenge belongs in the signed message and the row, not the envelope. F5 — the reused W15 counter condition has a redundant clause (inherited from V3). F6 — the repo's test signer (compared byte for byte with real `ssh-keygen` by an existing test) replaces per-test `ssh-keygen`.
- **Controller round trip (independent):** `ctl_roundtrip.sh` extracted all four patches from the plan file and applied them with `patch -p1 -E` to a fresh `git archive 8db22168`. My first run printed `rc=1` for every `diff -r`: the script's own redirect target `diff_$M_$d` hit an unbound variable under `set -u` (the same script ran at 407). Fixed to `${M}_${d}` and re-run on a fresh export: all four stages byte-identical to the planner's `st-M1`…`st-M4`; a one-line perturbation made it differ (rc 1) and restoring it matched again.

**Executor** (Agent `haiku`, foreground, 1.7 min, 51,503 tok; resolver `recipe claude:claude-haiku-5-5 declared:provider-pin` — spawned through the Agent tool's `haiku` pin as in 402–407, no hook denial; directive `~/.ailang/state/world-iter408/executor-directive.md`). M1→M4 in one run, every part `ACCEPT_OK` (PASS 15/19/27/31, 0 FAIL, 0 SKIP), snapshots 3/0/3, 4/0/4, 6/0/6, 8/0/8. Controller checks: each snapshot `MANIFEST` `shasum -c` ok; the executor's final tree `cmp`-equal to the planner's stage tree on all 8 files; main checkout untouched. `commit_ms.sh` per snapshot, 0 refusals, numstat 530 / 773 / 1,304 / 1,394 = plan: `d6ae2544`, `0b2a9b21`, `51a6e0e9`, `f8fd7c99`.

**Gates** (plan §6 `gates.sh` on `f8fd7c99`, background, 630 s): status, verify_ail, gofmt, build, vet, fence, pkgs, race, inv (51 named PASS) rc 0. `TOUCHED_PACKAGES` = `host/approval host/verbs` as planned; `DAEMON_ROUTES` 47 = base; `UNTOUCHED_VIOLATIONS` none; `PREEXISTING_CHANGED` none; `ENVELOPE_PRODUCTION_CALLERS` none (control `approval.CheckConfirmation` 1); `app.js`, `surface_page.go`, `credentials.go` unchanged. **`verify_go` rc 1**: the only failing package was the untouched `host/verifygate`, where `TestMissionConfigDefaultFlowIntegration` hit its 10 s fixture deadline (`instrument failure: context deadline exceeded`) while the judge ran in parallel. Re-run alone at `f8fd7c99`: PASS (4.8 s and 5.0 s per subtest). It is on the known load-sensitive list.

**Evaluator** (Agent `sonnet`, foreground, 15 min, 109,469 tok; resolver `agent-tool sonnet declared:alias-pin`; own worktree `.wt-world-iter408-row243-v4s2-eval`; directive asked open questions, no controller findings supplied): **PASS 92/100**, 0 blocking (correctness 28/30, tests 27/30, security 19/20, standards 10/10, scope 8/10). Its own ledger `~/.ailang/state/world-iter408/judge/mutants.tsv` (68 lines, ~45 mutants, ids `X-*`): each decoder refusal D1–D9, Validate rules, the message byte ranges, both digests' inputs and each key-law check (UV, UP, RP id, origin, challenge, revocation, enrollment) killed by named behavioural tests, goldens excluded and re-run; must-die and must-live controls behaved. 27 hostile decoder inputs of its own (duplicate keys at three depths, case-folded keys, raw and escaped NUL, overlong UTF-8, surrogate halves, BOM, `1e999999`, 100k-deep nesting) all refused with the zero value; a scratch package outside `host/approval` fails to compile `approval.SignedOperatorAction{ok: true}` (positive control compiles). M1 and M3 green alone. Judge-407's survivor reproduced SURVIVED at the base and KILLED at the head. F2–F5 confirmed. Non-blocking, ranked: (1) design §3.1 says the proof originates only in the `host/authority` judges; S2 guarantees only `host/approval`, and `ConfirmOperatorAction` trusts the `Approver` it is handed (as V3's `ConfirmBySignature` does) — the parked adapter must land before any caller; (2) `X-PK-ALIAS`: dropping the `bytes.Clone` of the credential id in the proof survives the whole package; (3) `CheckOperatorAction` does not re-bound the lifetime (matches the design; unpinned); (4) dropping `episode`/`kind`/`argv`/`options` from the encoding is killed by strict-JSON/limits/golden tests, not by `TestOperatorConfirmBinding` as the plan says (they are bound redundantly through `argv`); (5) M1 and M3 are ~530 lines each, over the per-milestone guide; (6) the `verify_go` red above (it did not re-run it).
- **Controller-reproduced before acting:** (2) in the eval worktree, line 126 `bytes.Clone(p.CredentialID)` → `p.CredentialID` (`git diff --stat` 1 file, `go vet` ok), `go test -count=1 ./host/approval/` → `ok` rc 0 — the mutant survives; restored, tree clean. Disposition: carried to S2b's planner (a one-line aliasing pin beside the parked adapters), not filed as a row (D-WORLD-109).

**Gate 3b:** PR #609 head `f8fd7c99`: first poll **`go host build + test gate: failure`** — `host/broker` `TestSeedRestartIgnoresASurvivingOldCopy/survivor` (`the old staging dir survived the restart: <nil>`); the job log (190,666 B, 175 `--- PASS`) shows `host/approval`, `host/authority` and the rest `ok`. `host/broker` is untouched (0 files in the diff), and the test is on the known-flaky list (2 earlier mentions in this log). Re-ran the failed job, confirmed the re-run was live (`pending=1`) before polling, then CI 3/3 green. `origin/dev` re-measured just before merging: still `8db22168`, PR `MERGEABLE CLEAN`. Merged at `cee5174c`; merge commit CI **3/3 green**. **LANDED.**

**Routing evidence:** base=8db22168689f31e19349b075e0199c8d9c940e61@2026-11-18T17:40:13Z (Gate 1; sprint worktree re-measured SAME at 17:41:23Z). Gate 4 base cee5174cf30d532574a17233780cf98b390003fa@2026-11-18T19:17:20Z (drift from Gate 1 = this sprint's merge only).
- Controller: `claude-opus-5-5` (tok: not reported).
- Designer: not spawned (plan fire; the doc is reviewed). Rotation pointer unchanged at `codex:gpt-6.2-sol` → next design gets `pi:ollama/glm-5.3:cloud`.
- Planner: `opus` via Agent (371,145 tok). Executor: `haiku` via Agent for the `claude:claude-haiku-5-5` pin (51,503 tok). Evaluator: `sonnet` via Agent (109,469 tok). **Every required role spawned; none failed to spawn.**
- Judge independence: same vendor, different model and fresh context (opus planner, haiku executor, sonnet judge), as 400–407.
- Metered: **$0.00** (no quorum, no metered lane).

**Progress:** 1.0 clauses 1–7 MET; goal unmoved. Row 243 V4 is at 2 of 12 sprints, with S2's adapter tail parked. The repo now defines exactly what a key holder signs to approve an operator action, and how a daemon and a runner judge that signature. Nothing calls it yet, and nothing answer-dependent is built.

**Ruled out:**
- **Treating the `verify_go` red as the sprint's:** the only failing package was untouched `host/verifygate`; the test passes alone with half its deadline to spare.
- **Treating the PR's first CI red as the sprint's:** `host/broker` is untouched and the failing test is on the flaky list; the same job shows the sprint's packages `ok`, and the re-run was green.
- **Folding the parked adapters into S3:** S3 waits on D-WORLD-119; the adapters do not, and the judge's finding (1) says they should land before any caller exists. Scheduled as S2b, next.
- **A re-quorum at pick:** the doc's r2 artifact is PROCEED 3/3 with no absent seat; nothing changed since.

**Retro:**
- **Held, an 8th time:** planner opus → controller round trip → executor haiku → commit_ms per snapshot → background gates beside an Agent judge in its own worktree. Gate 1 → merge CI green ≈ 1 h 35.
- **Instance 1: my own copied round-trip script carried a latent bug.** The `diff_$M_$d` redirect fails under `set -u`, so every `diff -r` read `rc=1`. It read as a mismatch, not as an error. Iteration 407 used the same script. Fixed in place and re-run, with a perturbation control. Lesson, already in standing rule 3a: a uniform reading across every case is a reason to check the instrument first.
- **Instance 2 of a background gate leg racing the judge** (`verify_go` deadline under load): the 407 run was green, this one was not. Record only; the re-run-alone discipline handled it.

**Next:** row 243 V4 sprint 2b planner: apply and prove the parked `host/authority` adapters (`PARK-authority-judges.patch`: `JudgeOperatorAction`, `RecheckOperatorAction`, `VerifiedOperatorAction` with an unexported constructor), and pin judge-408 notes (2) `X-PK-ALIAS` and (3) the runner's lifetime. It is answer-independent. S3–S8 wait for D-WORLD-119; S9–S11 for D-WORLD-120.

**Gate 5 bookkeeping:** rotation not due (#574 opened 2026-11-16 09:06Z; next Monday 2026-11-23). The weekly external-issue sweep already ran this week at 393.

## 409 — 2026-11-18 — row 243 V4 sprint 2b LANDED (PR #611, merge `07b30a5e`): the `host/authority` operator-action judges `JudgeOperatorAction`/`RecheckOperatorAction` read the approver from the store and return an unexported-field `VerifiedOperatorAction`; judge-408's `X-PK-ALIAS` killed and the runner lifetime pinned; unregistered; 526 lines; judge PASS 92; V4 now waits on D-WORLD-119 [PRODUCT]

**Pick:** row 243 V4 sprint 2b, the Next named in record 408 and the queue head.
- **Re-measured at pick** (`origin/dev` `4b1c39e2`): no S2b plan, branch or worktree (`git log origin/dev --grep S2b` → none; `gh pr list --author @me --state open` → none; `git worktree list` → no 409 tree). The parked patch `PARK-authority-judges.patch` was on disk (11,732 B) and `patch -p1 -E --dry-run` applied it cleanly to a fresh worktree at `4b1c39e2`. Quorum reader rc 0 for `w-world-address-bar-verbs-v4`: r1 `blocked` (reject/pass/pass), r2 `proceed` (pass ×3), `absent_reviewers` `[]` in both. `mission_decisions.sh --open` → D-WORLD-119, D-WORLD-120. S2b needs neither (S2 Appendix C: "needs no schema").
- **Clause map:** 1.0 clauses 1–7 MET, none UNMET. Release rows 150/151 still wait on attended steps. Row 243 serves D-WORLD-94 (D-WORLD-109 = A) and moves no clause.
- Claim sent to `controlplane` before routing (`inbox_1795038059946_78e064e1`).

**Gate 0:**
- Stamped first. `gh` is `sunholo-voight-kampff`. Billing CLEAN. Kill switch absent (`mission-world.disabled`).
- Directives: 0 allowlisted on #574 (18 comments), read from the epoch because the issue-scoped watermark is absent (the mission-scoped `mission-world-last-seen` reads `2026-10-27T09:03:42Z`). Nothing to `--mark`.
- Self-notices: rc 2 `watermark channel unavailable` (the #574 and #524 files are absent) — an instrument failure, not a verdict; same gap as 371–408.
- Inbox: no unread message addressed to `mission-world`.
- Skill drift: none — `SKILL.md` and all 13 `resources/*.md` at the resolved symlink target `cmp` equal to the driver clone's `origin/dev` (pin `d2f2ec238`). My first loop printed DIFF for all 13 because `ls resources/` yields bare names and the `git show` path lacked the prefix — the same slip 408 recorded; re-run with the prefix: all SAME.
- Dev CI at `4b1c39e2`: 3/3 check-runs green. Main checkout level with `origin/dev`; only untracked files (Mark's drafts, stale scratch), untouched.

**Planner** (Agent `opus`, foreground, 27.6 min, 283,161 tok; resolver `agent-tool opus fail-closed:env-pin`, `derive-planner-lane.sh` → `opus fail-closed:env-pin`; directive `~/.ailang/state/world-iter409/planner-directive.md`, the 408 directive adapted to S2b, handing judge-408 notes (1)–(4) as questions). Plan `w-world-address-bar-verbs-v4s2b-sprint-plan.md` + JSON, committed `1f083c2f`. No cut: 526 lines (116 production, 410 test) in three sub-milestones — M1 the judges (401), M2 judge-408 pins in `host/approval` (469), M3 writes-nothing (526). 61 mutants with whole-package verdicts: 59 killed, 2 named equivalents (`REUSED-SSH-KEY`, `REUSED-PK-CREDID` — the store selects the row by the artifact's own bytes), 0 survived, none resting on the AST test alone.
- **Premises found false (§1.4):** the parked patch's tests were NOT discriminating — its first mutant pass left **13 whole-package survivors**, which the new `TestOperatorJudgeReadsTheStore` (19 rows, each asserting step, error class and no proof) and `TestOperatorUnjudgedProofBindsNothing` kill. Judge-408 note (4) overstated what a binding row can kill (episode/kind/argv/options bind each other through `Validate`'s O9) — carried to Appendix C. The controller's round-trip script looped M1–M4 while S2b has three milestones.
- **Judge-408 notes:** (1) closed — neither judge takes an approver; each store-sourced check (unknown key, unknown credential, revoked, enrolled after now, stored format, counter, unreadable store, unparsable stored key, approver name) is its own mutant with a named killer. (2) `X-PK-ALIAS` reproduced SURVIVED at the base, killed by `TestOperatorPasskeyUseIsACopy` (with its accessor twin `PK-USE-ALIAS`). (3) pinned, not added: design §3.1 puts the 600 s bound on the signed request (checked at confirmation) and gives the runner expiry and revocation only; `TestOperatorRunnerRecheckBoundsExpiryNotLifetime` kills `R1-LIFETIME-ADDED`.
- **Controller round trip (independent):** `ctl_roundtrip.sh` with the loop edited to `M1 M2 M3`; all three stages `diff -r` rc 0 against the planner's `st-M*`; a one-line perturbation made it differ (rc 1), restoring matched (rc 0). The two `.orig` files the export lists are tracked at the base (`design_docs/verification/world-iter362/…`), not patch residue.

**Executor** (Agent `haiku`, foreground, 1.3 min, 43,260 tok; resolver `recipe claude:claude-haiku-5-5 declared:provider-pin` — spawned through the Agent tool's `haiku` pin as in 402–408, no hook denial). M1→M3 in one run, every part `ACCEPT_OK` (PASS 19/21/22, 0 FAIL, 0 SKIP), snapshots 2/0/2, 3/0/3, 4/0/4. Controller checks: each `MANIFEST` `shasum -c` ok; the final tree `cmp`-equal to the planner's `st-M3` on all 4 files. `commit_ms.sh` per snapshot, 0 refusals, numstat 401 / 469 / 526 = plan: `e582b846`, `7f490bd7`, `0e9e8d35`.

**Gates** (plan §6 `gates.sh` on `0e9e8d35`, background, 610 s, run BEFORE the judge touched Go — the judge's directive held it on the `ALL_DONE` file): status, verify_ail, gofmt, build, vet, fence, pkgs, race, inv (44 named PASS) rc 0. `TOUCHED_PACKAGES` = `host/approval host/authority` as planned; routes 47 = base; `UNTOUCHED_VIOLATIONS` none; `PREEXISTING_CHANGED` none; `JUDGE_PRODUCTION_CALLERS` none (control `authority.RecheckVerb` 1); envelope callers outside the judges none (control `approval.CheckConfirmation` 1); `app.js`, `surface_page.go`, `daemon.go`, `credentials.go` unchanged. **`verify_go` rc 1**: the only failing package was `host/daemon` (byte-identical to base, `git diff --quiet`), `TestSeToolsRunBenchmarkTasksEndToEnd` — the mock saw a stray `"GET  "` beside the one POST, the same signature as 337 and 404 — at load average 9.96 → 13.17 (sibling missions). Re-run alone at `0e9e8d35`: PASS (5.5 s).

**Evaluator** (Agent `sonnet`, foreground, 13.3 min, 81,956 tok; resolver `agent-tool sonnet declared:alias-pin`; own worktree `.wt-world-iter409-row243-v4s2b-eval`; directive asked open questions, no controller findings supplied): **PASS 92/100**, 0 blocking (correctness 29/30, tests 27/30, security 19/20, standards 9/10, scope 10/10). Its ledger `~/.ailang/state/world-iter409/judge/mutants.tsv` (34 lines; 31 ids not in the plan's ledger by `comm -23`): 27 killed by named behavioural tests, 2 survived (`ADMIT-DIGEST-SWAP` equivalent — the admit comparison is symmetric; a builtin `println` in a judge), 2 must-live controls survived as intended, must-die `NIL-BINDS` killed. A scratch package outside `host/authority` fails to compile `authority.VerifiedOperatorAction{rechecked: true}` (positive control compiles). `X-PK-ALIAS` survives at the base and is killed at the head. M1 and M2 green alone. Plan's lifetime reading of design §3.1 confirmed. Non-blocking: (1) the writes-nothing helper swaps `os.Stdout`/`os.Stderr` and `log` only, so a builtin `println` (fd 2 direct) survives; (2) the construction test is AST/reflect, but every mutant it kills is also killed behaviourally; (3) S3 must persist the recorded approver exactly as `format:name`. It said itself that it did not sample the plan's §4 rows.
- **Controller-reproduced before disposing:** (1) in the eval worktree, `println("leak")` inserted at the top of `JudgeOperatorAction` (`git diff --stat` 1 file, `go vet` ok), `go test -count=1 ./host/authority/` → `ok` — the mutant survives; restored, tree clean. Disposition: logged here and carried to S3's planner (an fd-level capture), not filed (D-WORLD-109). **Judge gap closed by the controller:** plan rows `J-RECHECKED`, `R-RECHECKED`, `RECHECKED-SIGNED` re-run in the eval worktree — each KILLED by its named test and by the whole `host/authority` package; a comment-only control SURVIVED; tree clean.

**Gate 3b:** PR #611 head `0e9e8d35`: CI **3/3 green** first time (two 9-minute poll segments). `origin/dev` re-measured just before merging: still `4b1c39e2`; PR `MERGEABLE CLEAN`. Merged with `--match-head-commit` at `07b30a5e`; merge commit CI **3/3 green**. **LANDED.**

**Routing evidence:** base=4b1c39e2a32fa129f2b340f4f9c2b6a312c63e0f@2026-11-18T21:40:24Z (Gate 1; sprint worktree SAME at 21:41:27Z). Gate 4 base 07b30a5ee1f6f45b26a87d9fa16ccae913313192@2026-11-18T22:53:48Z (drift from Gate 1 = this sprint's merge only).
- Controller: `claude-opus-5-5` (tok: not reported).
- Designer: not spawned (plan fire; the doc is reviewed). Rotation pointer unchanged at `codex:gpt-6.2-sol` → next design gets `pi:ollama/glm-5.3:cloud`.
- Planner: `opus` via Agent (283,161 tok). Executor: `haiku` via Agent for the `claude:claude-haiku-5-5` pin (43,260 tok). Evaluator: `sonnet` via Agent (81,956 tok). **Every required role spawned; none failed to spawn.**
- Judge independence: same vendor, different model and fresh context (opus planner, haiku executor, sonnet judge), as 400–408.
- Metered: **$0.00** (no quorum, no metered lane).

**Progress:** 1.0 clauses 1–7 MET; goal unmoved. Row 243 V4 is at 2 of 12 sprints with S2's tail now landed: a verified operator-action proof can come only from the two `host/authority` judges, which read the approver from the store. **Every remaining V4 sprint needs an answer** (design §8: S3–S8 wait for D-WORLD-119 A/B, S9–S11 also for D-WORLD-120, S12 for 119 A), so the loop has no answer-independent V4 work left.

**Ruled out:**
- **Treating the `verify_go` red as the sprint's:** `host/daemon` is byte-identical to the base; the test passes alone.
- **Re-running the parked patch as-is:** its own tests left 13 survivors; it landed only after the planner's mutant pass.
- **Adding the lifetime bound to the runner recheck (judge-408 note 3):** design §3.1 assigns it to confirmation; the runner's input has no creation time. Pinned instead.
- **A re-quorum at pick:** r2 PROCEED 3/3, no absent seat; nothing changed since.

**Retro:**
- **Held, a 9th time:** planner opus → controller round trip → executor haiku → commit_ms per snapshot → gates before the judge's drills (this time strictly sequenced via the `ALL_DONE` file) → judge sonnet in its own worktree. Gate 1 → merge CI green ≈ 1 h 20 (21:40Z → 22:58Z).
- **Instance 3 (third recorded failing iteration: 337, 404, 409) of `TestSeToolsRunBenchmarkTasksEndToEnd`'s stray-GET flake.** Under D-WORLD-109 (2) it becomes its own row: filed as row 261, a 258-class hygiene row under rule (d), as row 259 was.
- **Instance 2 of the controller's skill-drift loop printing DIFF for every resource file** (408, 409): `ls resources/` returns bare names. The gate's own snippet (gate-1 rule (d)) globs `"$base"/resources/*.md` and is correct; the slip is in my hand-rolled form. Use the gate's snippet verbatim.
- **A worktree helper returned rc 0 before the checkout existed** (standing rule 7's launcher-vs-work clause): the first `cd` failed; a re-check moments later found it. Polled the artifact thereafter.

**Next:** no answer-independent row 243 work remains. Until Mark answers D-WORLD-119, the next fire has no routable build row (release rows 150/151 are attended; rows 258/259/261 sit under rule (d)) and records a rule-(d) block report. With 119 = A or B, the next loop work is the V4 sprint 3 planner (schema 10 `operator_actions`, receipts and CAS; carrying judge-409's `println` fd-capture note and the `format:name` approver persistence).

**Gate 5 bookkeeping:** rotation not due (#574 opened 2026-11-16 09:06Z; next Monday 2026-11-23). The weekly external-issue sweep already ran this week at 393.

## 410 — 2026-11-19 — no pick: rule (d) block report; row 243 V4 S3–S12 wait on D-WORLD-119 (S9–S11 also 120); Gate 0's self-notices rc 2 added to the fleet ticket `gate0:self-notices-rc2-no-watermark-file` (row 262); dev green at `f267fac4` [ADMIN]

**Kind:** no-pick iteration. Nothing is loop-routable without an answer from Mark, so standing rule (d) applies. There is no design, plan, code or product acceptance. The record is the STATUS stamp (407 moved to the archive), the Current route line, queue row 262, this entry, the regenerated index, the dashboard and `verification/world-iter410/gate1-2-readings.txt`. An independent judge reviews the record before merge.

**Pick and why:** none. Clause map at `f267fac4`: 1.0 clauses 1–7 MET, so there is no UNMET clause and no critical-path row.
- **What moved since iteration 409's Gate-1 base (`4b1c39e2`):** only 409's own work. 7 commits: plan `1f083c2f`, M1–M3 `e582b846`/`7f490bd7`/`0e9e8d35`, merge `07b30a5e`, record `c6214e24`, merge `f267fac4`. 0 `| D-WORLD` lines changed.
- **Ledger:** `--check` valid, 107 rows; `--open` → D-WORLD-119, D-WORLD-120, both unchanged since iteration 406 filed them.
- **Row 243 (D-WORLD-109 granted scope):** V4 design §8, re-read this fire, gates S3–S8 on 119 A/B, S9–S11 also on 120 A/B, and S12 on 119 A. S1, S2 and S2b are landed. No answer-independent V4 sprint is left.
- **Row 150** waits on the attended A1 (descriptor publish, `session mint`, capture, promotion). Subjects in range matching `row ?150|A1|marketing` → **0**; control `row 243|record|merge` → **7**.
- **Row 151** needs row 188's N2d and N4a. Subjects since 2026-10-01 matching `N2d|N4a` → **0**; control `N2c` → **1**. Non-test `.go` files naming N2d's `file.json` route → **0**; control `surface` → **92**.
- **Rows 217–219** are on the attended track beside row 188 (D-WORLD-94). **Rows 258, 259, 261** sit under rule (d), below the release rows (D-WORLD-109 (2)).
- **External predicates:** no row in the routable set waits on an upstream party this fire; the only outward predicate is the new fleet ticket below.

**Gate 0:**
- Stamped first (rc 0). Kill switch armed (`mission-world.disabled` absent). `gh` is `sunholo-voight-kampff`. Billing CLEAN. The env was read by name, then specific non-secret values.
- Directives: 0 allowlisted on #574 (19 comments, read from the epoch because `mission-574-last-seen` is absent), 0 on #574 and #524 since `2026-10-27T09:03:42Z` (the mission-scoped watermark, per the charter's two-watermark rule), 0 ever on issue #1. Nothing to `--mark`; no watermark moved.
- **Self-notices (step 6a):** the prescribed form (two issue-scoped watermarks) → rc 2 `watermark channel unavailable`, as on every fire since 371. Re-run with `--watermark-file mission-world-last-seen` → **rc 1**, `A FIRE DIED`: 1 crash notice (#524, `2026-11-13T09:48:25Z`, rc 1) and 1 transcript notice (#524, fire `world-1794807071-53892`, 2 hits), control 107:4 ok. Both were already credited: the crash is the 2026-11-13 10:31 CET fire (driver log `slot-verdict: CRASHED at=gate-3 rc=1`, 10:48:12 local), listed in entry 393's "Orphans — three dead attempts"; the transcript notice is in entry 394's Gate 0 inbox line. So nothing new died. The died-mid-flight traces ran anyway: 0 open PRs from the fleet account, the three iteration-409 worktrees have 0 porcelain lines each, no `iter410` branch existed on origin, and the last two slot verdicts are `COMPLETED rc=0`.
- **Harness ticket:** the fleet already had `gate0:self-notices-rc2-no-watermark-file` open (filed by stapledon 2026-11-18). World had recorded the same rc 2 in about 39 entries without filing it. Filed World's occurrence under the same signature (`--blocking none`, `--workaround none`; message `inbox_1795052586078_70d26102`); `ticket open` now reads count 2, missions `stapledon,world`. Queue row 262 records it. One error in the filed evidence: it says the mission-scoped re-run "returns rc=0"; it exits rc 1. My first reading went through `| tail`, and that was tail's code. The verdict text in the ticket is correct.
- Inbox: `messages list --unread --inbox mission-world` → none. Skill drift: `SKILL.md` and all 13 `resources/*.md` at the resolved symlink target `cmp` equal to the driver clone's `origin/dev`, and equal to the pin's copies (gate-1 rule (d) form, with the `resources/` prefix this time).
- Dev CI at `f267fac4`: 3/3 check-runs success. Main checkout level with `origin/dev`; only the known untracked residue (Mark's drafts, `.wt-world-iter306-r0`, `.claude/worktrees/`), untouched.

**Gate 1:** `git rev-parse dev origin/dev` both `f267fac4…`. Gate-1 base `f267fac413182976765b697648740f1b145c5044@2026-11-19T01:39:50Z`. 0 open PRs from the fleet account. Record written in its own worktree (`mission/world-iter410-record` from `f267fac4`, `mission-worktree.sh add` then `wait` → ready).

**Designer / planner / executor:** not spawned. A no-pick iteration has no doc, plan or code, so the routing table has nothing to send them. No spawn was attempted, so none failed. The driver's routing note says the codex lane probed rc 1 this fire (planner → opus, executor → `claude:claude-haiku-5-5`); neither role was needed.

**Independent evaluator (record):** REQUIRED. `resolve-role-spawn.sh evaluator` → `agent-tool sonnet declared:alias-pin`. Agent `sonnet`, fresh foreground read-only context, judging the record commit before merge. Result in the follow-through below.

**Routing evidence:** base=f267fac413182976765b697648740f1b145c5044@2026-11-19T01:39:50Z (Gate 1).
- Controller `claude:claude-opus-5-5` (tok: not reported). Designer, planner, executor unused.
- Evaluator: Agent `sonnet` (tokens in the follow-through). Generator ≠ judge: opus controller, sonnet judge.
- Quorum: none (no doc). Metered $0.

**Progress:** 1.0 clauses 1–7 MET; goal unmoved. Row 243 V4 is at 2 of 12 sprints (S2b landed in 409). Every remaining V4 sprint waits on D-WORLD-119.

**Ruled out:**
- Picking 150, 151, 217–219, 243, 258, 259 or 261: measured reasons above.
- Treating the mission-scoped self-notices verdict (`A FIRE DIED`) as a new death: both notices are credited (393, 394), and the traces are clean.
- Hand-writing `mission-574-last-seen` to make the prescribed read work: Gate 0 step 8 forbids writing a watermark by hand, and the fix is the fleet's.
- Promoting the idle fingerprint: Gate 5 forbids it when the fire degraded, and this fire's codex probe failed (driver routing note).
- Re-asking D-WORLD-119/120: both asks are already in the digest, unchanged.
- Rotating the log: 21 live entries with this one, below ~40. The index is regenerated with `rotate-log --index-only`.

**Retro:**
1. **A "still" with a known cause, carried 39 times.** Each entry since 371 wrote "rc 2, an instrument failure, not a verdict" and moved on. That was correct per entry, but the gap was never escalated, and the gate could not see a death on its own for that whole stretch. The fix for the loop is the ticket, now filed, not a local workaround.
2. **Instance 3 of an exit code read through a pipe** (verification rule 3): my first self-notices re-run read `rc=0` from `| tail -15`. The script exits 1. Caught when I re-captured the readings under `bash` with the code taken directly. The ticket evidence carries the wrong code (noted above).
3. **zsh again** (memory `controller-shell-is-zsh`): `set -- $spec` did not word-split, and `PIPESTATUS` is absent, in my first readings capture. Re-run as a `bash` script; nothing was written from the bad run.
- Harness share: this entry files 1 ticket (an escalation, the correct outcome) and fixes nothing. No skill edit.

**Next:** with D-WORLD-119 = A or B, row 243 V4 sprint 3 planner (schema 10 `operator_actions`, receipts and CAS; carrying judge-409's `println` fd-capture note and the `format:name` approver persistence). Unanswered, another rule-(d) block report. A dev red outranks either way. The fleet's `harness-resolved` reply for `gate0:self-notices-rc2-no-watermark-file` unparks row 262.

**Gate 5 bookkeeping:** rotation not due (#574 opened 2026-11-16 09:06Z; next boundary Monday 2026-11-23 07:00 local). The weekly external-issue sweep ran this week at 393.

**Judge follow-through (iteration 410):** Agent `sonnet`, fresh foreground read-only context, judged record commit `8c851a1b`. Verdict **PASS 93/100, zero blocking**. It used 61,775 subagent tokens, 14 tool uses and 48 s. Generator ≠ judge: opus controller, sonnet judge. The directive gave it the files and seven open questions, not the controller's numbers.
- **What it re-measured:**
  - `--open` lists D-WORLD-119/120 only. V4 §8 gates every sprint after S2 on 119 or 120. D-WORLD-109 grants rows 201 and 243 only.
  - No non-test `.go` file names `file.json`. No subject since 2026-10-01 names N2d or N4a.
  - The 7 commits in `4b1c39e2..origin/dev`.
  - Both self-notices forms, with exit codes taken directly: rc 2 and rc 1, control 107:4 ok.
  - Its crediting check: entry 393's orphan line and the driver's `CRASHED at=gate-3 rc=1` at 10:48:12 match the crash notice. Entry 394 matches the transcript notice.
  - The ticket listing (count 2, `stapledon,world`). `--check` valid at 107 rows. Heading formats. STATUS 407 at the archive's end.
  - It found no closing keyword in the commit message or the diff.
- **Non-blocking notes and their disposition:**
  1. The ticket's evidence says rc 0, not rc 1. Already disclosed above. I did not file a correction: a second filing under the signature would count as a second occurrence.
  2. The 2026-11-13 14:31 fire (CRASHED at gate-3, 15:20:16 local) is not in the reconciliation. It produced no second notice because the driver's crash comment suppresses identical failures until the rc changes, and entry 393 credits it.
  3. "Since 371" is inherited from the carried entries; entry 393 confirms it. It did not count the entries itself.
  4. Its run omitted `--driver-src`, so it printed `signature-source: unchecked`. The verdict was the same.
- **What it did not check:** dev CI check-runs, skill-drift `cmp`, the codex probe note, the ticket message id, the inbox read, and every other "iteration 409" reference in the charter.

## 411 — 2026-11-19 — no pick: row 262 RESOLVED on the fleet's fix (ailang #2421, `cf540337e`), checked first-party; then a rule (d) block report, row 243 V4 S3–S12 still wait on D-WORLD-119 (S9–S11 also 120); dev green at `2c9bb491` [ADMIN]

**Kind:** no-pick iteration with one bookkeeping unpark. Nothing is loop-routable without an answer from Mark, so standing rule (d) applies. The one change in state since iteration 410 is the fleet's `harness-resolved` reply for row 262. It is checked and recorded here. There is no design, plan, code or product acceptance. The record is the STATUS stamp (408 moved to the archive), the Current route line, queue row 262, this entry, the regenerated index, the dashboard and `verification/world-iter411/gate1-2-readings.txt`. An independent judge reviews the record before merge.

**Pick and why:** row 262's unpark (bookkeeping), then nothing. Clause map at `2c9bb491`: 1.0 clauses 1–7 MET, so there is no UNMET clause and no critical-path row.
- **What moved since iteration 410's Gate-1 base (`f267fac4`):** only 410's own record. 3 commits: `8c851a1b`, `beee5376`, merge `2c9bb491`. 0 `| D-WORLD` table rows changed.
- **Ledger:** `--check` valid, 107 rows; `--open` → D-WORLD-119, D-WORLD-120, both unchanged since iteration 406 filed them.
- **Row 243 (D-WORLD-109 granted scope):** V4 design §8 gates S3–S8 on 119 A/B, S9–S11 also on 120 A/B, and S12 on 119 A. S1, S2 and S2b are landed. No answer-independent V4 sprint is left. Nothing in the 3 commits touches the design.
- **Row 150** waits on the attended A1. Subjects in `f267fac4..2c9bb491` matching `row ?150|A1|marketing` → **0**; control `record|iteration 410` → **3**.
- **Row 151** needs row 188's N2d and N4a. Subjects since 2026-10-01 matching `N2d|N4a` → **0**; control `N2c` → **1**. Non-test `.go` files naming `file.json` → **0**; control `surface` → **109**. (Iteration 410 read 92 for the same control from the working tree; this fire reads `git grep` over `origin/dev`, a different scope. It is a control, not a quantity this record relies on.)
- **Rows 217–219** are on the attended track beside row 188 (D-WORLD-94). **Rows 258, 259, 261** sit under rule (d), below the release rows (D-WORLD-109 (2)).
- **External predicates:** row 262's (the fleet's reply) flipped this fire, see below. No other routable row waits on an outside party.

**Row 262 — the fleet's reply, checked:**
- Reply `inbox_1795055723042_70effb63` from `mission-fleet`, category `harness-resolved`, signature `gate0:self-notices-rc2-no-watermark-file`: ailang #2421, merge `cf540337ede442443011db75c4cd9bf2776f7d33`, fleet iteration 123, its judge sonnet PASS 91. It says step 6a now reads the whole of each issue, labelled DEGRADED, when every watermark file is absent, and exits 0 or 1 by hits. A present but unreadable, empty or invalid file is still rc 2.
- `cf540337e` is an ancestor of `origin/dev` in the driver clone (rc 0) and is this fire's pin (`AILANG_DRIVER_PINNED`, `git -C $MISSION_DRIVER_ROOT rev-parse HEAD`).
- Same inputs, two scripts (both issue-scoped watermark files absent, as on every fire since 371): the pinned script exits **rc 1**, `crash=1 transcript=1` on #524, control 107:4 ok. The fix's parent (`git show cf540337e^1:scripts/mission_gate0_self_notices.sh`) exits **rc 2**, `watermark channel unavailable`. That is the negative control: the change in reading comes from the fix, not from the inputs.
- `ailang mission ticket open` no longer lists the signature (0 rows; control 18 open rows).
- The two notices the read surfaces were credited before: the crash (#524, 2026-11-13T09:48:25Z, rc 1) by entry 393's orphan line, matched by time to the driver's `CRASHED at=gate-3 rc=1` (iteration 410 and its judge both checked this; the URL text itself appears 0 times in the log), and the transcript notice (fire `world-1794807071-53892`, 2 hits) in entry 394. So nothing new died.
- Row 262 is now `RESOLVED (iter-411)` in place, with the resolution and a later-fires check. The fleet message is acknowledged in the report below.

**Gate 0:**
- Stamped first (rc 0). Kill switch armed (`mission-world.disabled` absent). `gh` is `sunholo-voight-kampff`. Billing CLEAN. The env was read by name, then specific non-secret values.
- Directives: 0 allowlisted on #574 (20 comments, read from the epoch because `mission-574-last-seen` is absent). Nothing to `--mark`; no watermark moved.
- Self-notices (step 6a), prescribed form: rc 1, `A FIRE DIED` on the two credited #524 notices (above). Died-mid-flight traces ran anyway: 0 open PRs on the repo, the iteration-409/410 worktrees have 0 porcelain lines each, no `iter411` branch existed on origin before this record, and the last two slot verdicts are `COMPLETED rc=0`. The transcript notice count is 1 (fire `world-1794807071-53892`, 2 hits), the same one entry 394 reported.
- Inbox: one unread on `mission-world`, the `harness-resolved` reply above.
- Skill drift: `SKILL.md` and all 13 `resources/*.md` at the resolved symlink target `cmp` equal to the driver clone's `origin/dev` (14 compared, 0 drift).
- Dev CI at `2c9bb491`: 3/3 check-runs success. The main checkout was 3 behind (`2c9bb491` not yet pulled) with only the known untracked residue; fast-forwarded with `--ff-only`, nothing else touched.

**Gate 1:** `git rev-parse dev origin/dev` both `2c9bb491…`. Gate-1 base `2c9bb491f52bf42acff7ddd9c4ec57659ec7c42e@2026-11-19T05:40:36Z`. 0 open PRs. Record written in its own worktree (`mission/world-iter411-record` from `2c9bb491`, `mission-worktree.sh add` then `wait` → ready).

**Designer / planner / executor:** not spawned. A no-pick iteration has no doc, plan or code, so the routing table has nothing to send them. No spawn was attempted, so none failed. The driver's routing note says the codex lane probed rc 1 this fire (planner → opus, executor → `claude:claude-haiku-5-5`); neither role was needed.

**Independent evaluator (record):** REQUIRED. `resolve-role-spawn.sh evaluator` → `agent-tool sonnet declared:alias-pin`. Agent `sonnet`, fresh foreground read-only context, judging the record commit before merge. Result in the follow-through below.

**Routing evidence:** base=2c9bb491f52bf42acff7ddd9c4ec57659ec7c42e@2026-11-19T05:42:06Z (Gate 4).
- Controller `claude:claude-opus-5-5` (tok: not reported). Designer, planner, executor unused.
- Evaluator: Agent `sonnet` (tokens in the follow-through). Generator ≠ judge: opus controller, sonnet judge.
- Quorum: none (no doc). Metered $0.

**Progress:** 1.0 clauses 1–7 MET; goal unmoved. Row 243 V4 is at 2 of 12 sprints. Every remaining V4 sprint waits on D-WORLD-119. Row 262 closed.

**Ruled out:**
- Picking 150, 151, 217–219, 243, 258, 259 or 261: measured reasons above.
- Treating step 6a's `A FIRE DIED` as a new death: both notices are credited (393, 394), and the traces are clean.
- Taking the fleet's word for the fix: checked against the parent script on the same inputs instead.
- Correcting iteration 410's filed evidence (it said rc 0 where the script exits 1): the ticket is resolved, and the error is already disclosed in entry 410.
- Promoting the idle fingerprint: Gate 5 forbids it when the fire degraded, and this fire's codex probe failed (driver routing note).
- Re-asking D-WORLD-119/120: both asks are already in the digest, unchanged.
- Rotating the log: 22 live entries with this one, below ~40. The index is regenerated with `rotate-log --index-only`.

**Retro:**
1. **The escalation path worked, in one fire.** Iteration 410 filed World's occurrence at about 01:43Z; the fleet's fix merged and its reply arrived by 02:35Z; this fire checked it. A ticket was the right outcome and a local workaround would have been wasted.
2. **A predicate you stop reading is a predicate you cannot see flip.** The reply was the only change in state since 410. It showed up only because Gate 0 reads the inbox and Gate 2 re-reads external predicates; the queue row itself said nothing.
3. Shell: the first env print used zsh `${(P)…}` and was refused by the harness; re-run under `bash -c` (memory `controller-shell-is-zsh`). Nothing was written from the refused command.
- Harness share: this entry closes 1 ticket's row and files none. No skill edit.

**Next:** with D-WORLD-119 = A or B, row 243 V4 sprint 3 planner (schema 10 `operator_actions`, receipts and CAS; carrying judge-409's `println` fd-capture note and the `format:name` approver persistence). Unanswered, another rule-(d) block report. A dev red outranks either way.

**Gate 5 bookkeeping:** rotation not due (#574 opened 2026-11-16 09:06Z; next boundary Monday 2026-11-23 07:00 local). The weekly external-issue sweep ran this week at 393.

**Judge follow-through (iteration 411):** Agent `sonnet`, fresh foreground read-only context, judged record commit `6ff2af6a`. Verdict **PASS 90/100, zero blocking**. It used 51,886 subagent tokens, 10 tool uses and 39 s. Generator ≠ judge: opus controller, sonnet judge. The directive gave it the files and seven open questions, not the controller's numbers.
- **What it re-measured:**
  - The reply message, its commit, and that `cf540337e` is the pinned driver's HEAD.
  - Both step-6a scripts on the prescribed inputs, exit codes taken directly: pinned rc 1 (DEGRADED whole-issue read, `crash=1`, `transcript=1`, control 107:4 ok), parent rc 2 `watermark channel unavailable`. Both issue-scoped watermark files absent.
  - The credits in entries 393 (crash, matched by time) and 394 (transcript fire).
  - `--check` valid at 107 rows, `--open` = 119 and 120, 0 changed `| D-WORLD` rows in `f267fac4..2c9bb491`.
  - V4 §8's sprint gates, rows 150 and 151's stated blockers, and that the D-WORLD-64 exception needs a routable row.
  - The rotation: 3 stamps left, 408 in the archive, queue rows unchanged, the canonical heading regex.
  - The commit message (no closing keyword), `check-ignore` on the readings file, dev CI 3/3 and 0 open PRs.
- **Non-blocking notes and their disposition:**
  1. The crash credit is by time, not URL, and the time gap is not spelled out here. Entry 410 and its judge matched the notice (09:48:25Z) to the driver's `CRASHED at=gate-3 rc=1` (10:48:12 local), which is the same instant. No change.
  2. The "until #524 rotates out on 2026-11-23" date. #574 opened 2026-11-16 and the next Monday-07:00 rotation is 2026-11-23, at which point #574 becomes `-prev` and #524 drops out of the read. This is a prediction for later fires, not a measurement.
  3. "Acknowledged in the report below" was a future action when judged. The fleet message is acked after merge, and the Gate 5 report names it.
- **What it did not check:** the row 150/151 grep counts, the skill-drift `cmp`, the slot verdicts and worktree porcelain, the ticket-open counts, the full index and dashboard diffs, and the middle of the long stamp text.

**Gate 3b (record PR #614):** head `4fcdfb5a` run 44467497039 attempt 1 read 2/3. `ailang-code verify gate` and `go host build + test gate` passed. `go race gate` failed on `host/broker` `TestSeedRestartIgnoresASurvivingOldCopy/survivor` (`exec_seed_test.go:656`, "the old staging dir survived the restart: <nil>"). The job's two `DATA RACE` warnings are its own known-positive control (`racecontrol/main.go`), as expected. Attribution: the PR changes only `design_docs/` (0 other paths), and the parent `2c9bb491` was 3/3 on the same gate. This is row 259's flake, instance 6, now added to that row. `gh run rerun --failed` → attempt 2 3/3. The instance note is a new head, so the PR is re-polled on it before merge.
