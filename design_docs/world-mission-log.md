# Ailang World Mission — iteration log

Append-only. One entry per outer-loop iteration. Newest at the bottom.

---

> **Older entries are ARCHIVED.** This file holds the newest 20. The full record of every
> iteration is in `world-mission-log-archive.md`, and a one-line index of ALL of them —
> the thing to grep before picking work, so the loop never repeats itself — is in
> `world-mission-index.md`.

## 378 — 2026-11-10 — row 201 A1b sprint 4 built: the decided list, the revoke listing, and `approve --sign-with` with its real-tool tests, the signer moved to a scrubbed `host/sshsign` after the sweep caught an AC10 red; judge PASS 93 [PRODUCT]

**Pick:** row 201 `w-world-approvals-by-key`, A1b sprint 4, the route iteration 377 named: the sprint-3 plan's Appendix C, with no re-design. D-WORLD-109 = A applies: fakes, temp stores and test keys only.
- **Re-measured at pick** (`origin/dev` `1fce6999`):
  - No open PR from this account.
  - No `DecidedMintRequests` or `--decided` anywhere in `cmd/`/`host/` Go source. The only `"decided"` hits were broker status strings.
  - No sprint-4 worktree.
  - `mc_quorum_gate2 w-world-approvals-by-key-a1b` rc 0, so no pick-time quorum was owed.
- **Ledger:** valid (100 rows). OPEN: D-WORLD-105 only.
- **Clause map:** 1.0 clauses 1–7 MET, none UNMET. Release rows are not routable:
  - Row 150's attended A1: 0 subjects since 2026-11-07 matching `row ?150|A1\b` (control `A1b`: 16).
  - Row 151's row-188 N2d/N4a: 0 since 2026-10-01 (control N2c: 1).

**Gate 0:**
- Stamped first. `gh` is `sunholo-voight-kampff`. Billing CLEAN. Kill switch absent.
- Directives: 0 allowlisted on #524, out of 10 comments. The watermark is absent, so all comments were read.
- Self-notices: rc 2, an instrument failure and not a verdict. No `-last-seen` file exists for 524 or 449, the same as 371–377. A hand read of #524 found no `FAILED to complete` notice.
- The `mission-world` inbox has 0 unread.
- Skill drift: none. SKILL.md and all 13 resource files at the resolved symlink target match the driver clone's `origin/dev`.
- Dev at `1fce6999`: `checks=3`, all success. Local dev was 15 behind and was fast-forwarded (only untracked dirs in the tree).

**Plan:** Agent `opus` (`fail-closed:env-pin`), 38 min, 397,156 tok. Committed `d3a3ffca`. Same form as iteration 377: every sub-milestone is a compiled unified diff in Appendix P, with an expected numstat measured after each commit.
- **Committed part:** T1 (the sprint-3 judge gaps SIG-NOW, SIG-FORMAT and A4-EXP, tests only), D1 (`DecidedMintRequests`), D2 (`approve --list --decided`), D3 (K8/K9) = 359 lines.
- **Gated tail:** S1–S3, the `--sign-with` seam, fake-signer tests and the production signer with XS1/XS2, the CI step and QUICKSTART. The gate was ≤ 1,463 after D3; it measured 359.
- **Mutants:** 68 run in scratch; 66 killed, 2 named equivalents.
- **Premises found false at HEAD:**
  - **P2:** `DecidedMintRequests(ctx, limit)` cannot serve K8, so it became `(ctx, approvedBy, limit)`.
  - **P4:** the signer's root pin moved to `approve_sig.go`.
  - **P6:** the CI step timeout went from 3 to 8 min.
  - **P8:** `TestApproverRevoke/yes` gains the K9 line.
  - **Item 4:** the revoke SQL's `revoked_at IS NULL` was already killed at base by `TestRevokeApprover`.

**Executor:** env `recipe claude:claude-haiku-5-5`, spawned as Agent `haiku`. Five runs: T1+D1, D2+D3, S1+S2, S3, S4. No hand edits. The controller rebuilt each sub-milestone from its snapshot with `commit_ms.sh`:
- T1 `f267e1b8` → 33
- D1 `387e3b07` → 181
- D2 `1bfe8c88` → 278
- D3 `e9d70925` → 359 (size gate PASS)
- S1 `969626c4` → 461
- S2 `4b5f9761` → 615
- S3 `784c9413` → 909
- S4 `5e85f8bb` → **1,014**

Every value equals the plan's value, except S3's (896), for the correction below.
- **Executor finding, adjudicated by measurement (rule 3h):**
  - **The gap:** §3 names `TestCappedBufferKeepsItsBound`, but the S3 patch did not contain it. Its `-run` regex matched nothing and still exited 0.
  - **The check:** the controller compared all 17 sprint files with the planner's scratch tree. 16 were byte-identical; the 17th differed only by that 13-line function.
  - **The fix:** the function was copied from the planner's tree, not rewritten. MUT-CAP-WRITE was re-run against it (killed), and the source was restored byte-identical. Recorded as plan §7.1 (`3a6135c2`).
- **Run-5 slip:** the executor disclosed running `git add -N -n`, a dry run. The controller checked the index: `git diff --cached` was empty.

**Controller gates:** plan §6 in one bash shell, run BEFORE the judge.
- **Sweep 1 at `784c9413`: rc 1.** Two reds, both in `host/broker`:
  - **(a) A REAL regression from S3: `TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential`.** It found 13 exec sites but only 12 drivers. The new `exec.CommandContext("ssh-keygen", …)` sat in `package main` with the full inherited environment, so the registry credential reached the child, and no census driver could import it. The planner's §6 list ran only the touched packages, and this census enumerates exec sites across the whole repo.
  - **(b) `TestWorktreeSnapshotRefusesAConfiguredFilter`:** the known load-sensitive control arm (earlier instances at iterations 300, 334 and a row-257 sprint). It passed 5/5 alone, and the sprint touches neither package.
- **Fix:** the planner was resumed by SendMessage (13.5 min, 452,035 tok) and compiled **S4** (`7a61cff0` plan, `5e85f8bb` code):
  - A new `host/sshsign.Sign` is the single exec site, with `cmd.Env = childenv.Scrubbed(os.Environ())`. The 16 KiB cap, `WaitDelay` and the G1–G3 texts are unchanged.
  - `signWithSSHKeygen` became a one-line call.
  - The census gained `driveSSHSign`.
  - **Mutant MUT-ENV-SCRUB:** dropping the scrub fails the census ("the child observed AILANG_REGISTRY_API_KEY") and the new unit test `TestSignHandsTheSignerAScrubbedEnvironment`.
  - S4's own mutants: 18 run, 17 killed, 1 named equivalent.
- **Sweep 2 at `5e85f8bb`: rc 0.** This covers:
  - `verify_ail.sh`, gofmt, build+vet and the compile fence;
  - the touched packages, `-race`, and the real-tool leg (PASS ×4, no SKIP);
  - the no-SKIP `-v` run;
  - `VERIFY_GO_LEGS=plain verify_go.sh` (114 ok lines, 0 FAIL; `host/broker` 344 s).

  The `DATA RACE` lines are the race control.
- **Unattributed single red:** the first `commit_ms.sh S1` run had `cmd/ailang-worldd` FAIL at 175.6 s. Its log was lost to a `tail`. On the byte-identical tree the package was ok at 112.6 s, and it was ok in every later run. Details are in `verification/world-iter378/controller-gates.md`.

**Judge:** sonnet via the Agent tool, on its own worktree at `5e85f8bb` (generator ≠ judge). **PASS 93/100, zero blocking.**
- **Sub-scores:** correctness 28, tests 28, security 19, standards 9, scope 9.
- **Controls:** the must-die control was killed and the must-live control survived.
- **Mutants:** about 50 were killed by their named tests, covering `DecidedMintRequests`, `--decided`, `--sign-with` F1/F2/F7/F7a/F8, `host/sshsign` (G1–G3, cap, argv, stdin, `WaitDelay`, scrub), and K8/K9. 3 did not compile and were not counted.
- **Survivor K8-ERRCODE:** in `approver.go:449`, revoke's K8 read-only reopen failure returns `exitFatal`; changing it to `exitOK` survives `TestApprover*`. The **controller reproduced it** and restored the file clean. It is non-exploitable: the revoke has already committed, and only the follow-up listing's exit code is unpinned. Carried forward.
- **Other non-blocking findings:**
  - NB2: `host/sshsign` has no unit test for G1 (a failed run); only the cmd test covers it.
  - NB3: PATH is not narrowed. This is documented, and the output is verified against the enrolled key.
  - NB4: QUICKSTART §15's new lines were not executed; they are attended.
- **Security properties confirmed:**
  - The registry credential is absent from the signer while `SSH_AUTH_SOCK`/`HOME`/`PATH` pass through.
  - argv is fixed, there is no shell, and the message goes on stdin.
  - Output is capped and verified.
  - The deadline equals the remaining lifetime, and an expired request is refused before the signer runs.
  - Revoke's write root closes before K8 lists, and it revokes no session.
  - D-WORLD-109 holds: no `~/.ssh`, `/dev/tty` or `~/.ailang/world`.

**Gate 3b:** PR #543, head `5e85f8bb`.
- `pull_request` run 42880342619: 3/3 success.
- The CI job log shows the real-tool leg PASS on Linux: `TestSSHSIGCrossCheckWithSSHKeygen` ×2, `TestTestSignerAgreesWithSSHKeygen`, and `TestSignWithSSHKeygenEndToEnd` ×2.
- The auto-close keyword scan read 0 (control 1).
- Merged by merge commit **`9429e798`** (`--match-head-commit`).
- Push checks on the merge `9429e798`: 3/3 success (all three jobs, the same set).
- Row 201 A1b is **code-complete**. What remains of A1b is Mark's attended exit run (QUICKSTART §15).

**Routing evidence:** base=1fce6999ebf875299254a1d77003e9e3351a3347@2026-11-10T16:50:36Z.
- Controller: `claude-opus-5-5` (tok: not reported).
- Designer: not spawned. Quorum: not run (the doc is reviewed).
- Planner: `opus` via Agent (397,156 + 452,035 tok).
- Executor: `claude:claude-haiku-5-5` via Agent `haiku` (38,466 + 41,522 + 42,852 + 47,709 + 45,582 tok).
- Evaluator: `sonnet` via Agent (79,233 tok).
- The driver noted codex unusable at preflight (probe rc 75); no role needed it.
- Metered total: $0.00. Details in `verification/world-iter378/routing.md`.

**Progress:** 1.0 clauses 1–7 MET; goal unmoved. Approval by key is now usable without leaving the CLI:
- `approve <id> --sign-with <key>` signs World's exact message with `ssh-keygen` (agent, file or `sk` key), verifies the result against the enrolled key, and approves.
- `approve --list --decided` shows what was approved and by whom.
- `approver revoke` lists the sessions the key approved, so the operator can revoke them.

The A1b design is fully built. A terminal `y` is no longer the only way to approve a mint, and nothing here needs a real key until Mark's attended exit run.

**Ruled out:**
- **Treating the AC10 census red as environmental:** it named the new site, and the mechanism (unscrubbed env, undriven site) was read in the code.
- **Fixing the AC10 red by adding the cmd file to the driver map:** `package main` cannot be imported, and a driver-less entry would be a vacuous pass.
- **Attributing `TestWorktreeSnapshotRefusesAConfiguredFilter` to the sprint:** it passed 5/5 alone, its packages are untouched, and it has prior instances.
- **Dismissing K8-ERRCODE on the judge's label alone:** it was reproduced first.
- **Treating the executor's missing-test report as executor error:** it was correct, and the planner's own tree proved it.

**Retro:**
- **A planner's gate list scoped to the touched packages misses repo-wide inventories.** This repo has several: the AC10 exec-site census, the import allowlist, context roots, and the verifygate subprocess sink. The controller's full sweep caught it, so the sweep-before-judge order earned its keep. Next directive: ask the planner to run every repo-wide inventory test, found by grepping for tests that walk `cmd/` or `host/`, whenever the sprint adds an exec site, a package, an import or a context root.
- **A patch cut from stages can drift from the tree the mutants ran on.** The planner added a test after generating its patches. Next directive: regenerate the patches last, and `cmp` the round-tripped tree against the mutation tree, not only against the stages.
- **Never pipe a gate's only copy of its output through `tail`.** The S1 red's log is gone. `commit_ms.sh` callers now `tee` to a file.

**Next:** row 201's NEXT moves to **A2–A4**, starting with A2 (WebAuthn/passkeys). A2 needs a new child design from the rotation designer, then a quorum. The route line lists row 243 after row 201, and the queue governs.
- Carry the judge's K8-ERRCODE and NB2 test gaps into the next row-201 sprint.
- Mark's attended A1b exit run (QUICKSTART §15 with `--sign-with` and `--list --decided`) stays attended.

## 379 — 2026-11-10 — row 201 A2 sprint 1 built: the pure WebAuthn verification law (CBOR, COSE, authData, clientData, origin, registration law, 128-bit enrolment code), from a new child design passed by quorum r2 and the carve-out; judge PASS 93 [PRODUCT]

**Pick:** row 201 `w-world-approvals-by-key`, parent milestone A2 (WebAuthn/passkeys), the route iteration 378 named. D-WORLD-109 = A applies: test keys, temp stores and Chrome's virtual authenticator only.
- **Re-measured at pick** (`origin/dev` `98374909`):
  - No open PR from this account; no `iter379` worktree; the index has no WebAuthn/passkey entry.
  - No child design for A2 existed (`design_docs/planned/` held only the a0a1 and a1b children).
  - `mc_quorum_gate2 w-world-approvals-by-key-a2` rc 1 (no artifact), so the quorum was owed.
- **Ledger:** valid (100 rows). OPEN: D-WORLD-105 only.
- **Clause map:** 1.0 clauses 1–7 MET, none UNMET. Release rows are not routable:
  - Row 150's attended A1: 0 subjects since 2026-11-07 matching `row ?150|A1\b` (control `A1b`: 20).
  - Row 151's row-188 N2d/N4a: 0 since 2026-10-01 (control N2c: 1).

**Gate 0:**
- Stamped first. `gh` is `sunholo-voight-kampff`. Billing CLEAN. Kill switch absent.
- Directives: 0 allowlisted on #524, out of 11 comments. The issue-scoped watermark is absent, so all comments were read.
- Self-notices: rc 2, an instrument failure and not a verdict; no `-last-seen` file exists for 524 or 449 (same as 371–378).
- The `mission-world` inbox has 0 unread.
- Skill drift: none. SKILL.md and all 13 resource files at the resolved symlink target match the driver clone's `origin/dev`.
- Dev at `98374909`: `checks=3`, all success. Local dev equal to origin.
- PR #545 (`fix/json-int-accessors`, se-tools `asNumber`) merged after record 378, outside the loop; it does not touch row 201's files.

**Design** (`w-world-approvals-by-key-a2.md`, a child of the attended parent). Resolver `recipe claude:claude-opus-5-5 declared:provider-pin`; spawned as Agent `opus` (precedent). Rotation file reads `claude:claude-opus-5-5`; the resolver's answer was followed.
- **r0** (22.6 min, 1,328 lines): three sprints. Challenge = SHA-256 of A0's 90-byte message (`approval.MessageHash`, already pinned); rpID = the origin's host, only `localhost` or `*.ts.net`; ES256 and EdDSA only; credentials schema 6 `webauthn_credentials`; the browser result reaches the CLI as a file, so A2 adds no write route; one static `/surface/passkey` page. The designer drove a headless Chrome 157 virtual authenticator (test key) and verified both algorithms with the standard library. It measured K8-ERRCODE wider than reported (`approver.go:449` and `:455`).
- **Quorum r1 BLOCKED 3/3** ($0.4386; author Anthropic benched; seats oc-kimi-k3, gpt6-1-sol, gemini-4-argon; absent none). The controller measured each premise first (rule 3f); all three real:
  - gpt: the 40-bit enrolment code derives only from the registration file (doc lines 511–521), so a read-then-substitute attacker can grind a collision; "guess blind" was false.
  - gemini: the `blob:` save under `surfaceCSP` was UNVERIFIED with a fallback (line 833).
  - kimi: no Verification Log row read the parent's or the ledger's ruling text; the 16-byte credential-ID floor was UNVERIFIED.
- **r1** (4.8 min, designer resumed by SendMessage): a 128-bit public fingerprint (26 base32 characters, read-then-substitute attacker modelled, CW5 rewritten); V31 measured the `blob:` download in Chrome under the exact CSP (completed, byte-equal, zero violations, with a violation-detecting control); V32–V34 quote the parent, the ledger and L3 §4 (the 16-byte floor is an authenticator rule; floor dropped to 1–1023).
- **Quorum r2 BLOCKED 2/3** ($0.5997; gemini PASS). Both premises measured real:
  - kimi: explicit default ports were accepted and stored, but `location.origin` elides them. Measured V36: Chrome 157 and Node's WHATWG URL both serialise `https://x.ts.net:443/` → `https://x.ts.net` (control `:7644` kept).
  - gpt: only WebAuthn's `timeout` hint bounded the browser wait; no AbortController.
  - Per round: r1 = enrolment code + CSP save + authority/floor; r2 = origin default port + ceremony abort. Different surfaces each round, no SPLIT signal.
- **r2, controller narrow-refinement carve-out** (ratified; precedents 372, 373, 375): both fixes applied verbatim (default-port rule + OR1 arms; the AbortController paragraph + XB2). One named adaptation: V36 measures Chrome's URL serializer via `new URL(...).origin`, not `location.origin` on a served `:443` page (needs TLS on 443).
- Fable/designer diet: one doc, r0 + one protocol-mandated revision. Held.

**Plan:** `derive-planner-lane.sh` → `opus fail-closed:env-pin`; Agent `opus`, 58 min, 445,890 tok. Committed `3eaf2910`.
- **Sprint 1 = M0, M1, M2, M3a = 1,958 lines** (cap 2,000). The design's §9 fallback cut fired on both clauses: the fixtures are the iteration-379 Chrome recording with `generate.sh` stubbed until M7, and the assertion half (M3b: `ApproveByAssertion`, `PasskeyUse`, `approvaltest.Passkey`, PA1–PA4, RC1, TP1, XW1/XW2) moves to the head of sprint 2 (Appendix C).
- Mutants: 197 on the planner's final tree, 184 killed, 13 named equivalents; base controls MUT-K8A/K8B/SG1 survive at base (confirms V3/V4).
- Round trip: the patches extracted from the plan file, applied to a fresh export of `03d78b6a`, reproduce the mutant tree byte for byte.
- **Premises found false or incomplete:** P2 the probe kept only `assertion2`'s authData; P3 PE5's rename turns `approver_enroll_test.go:230` red in sprint 1 (edited in M3a); P4 AD1 cannot import `approvaltest` (cycle), so it edits Chrome's bytes; P5 1,958 vs the design's ~1,800; P7 a failed `json.Decoder` keeps returning `More()` true (a draft hung; every token error is now refused and pinned).
- A new load-sensitive test was noted: `TestMissionConfigDefaultFlowIntegration/invalid_ledger` (10 s fixture deadline, 7.3 s alone).

**Executor:** env `recipe claude:claude-haiku-5-5`, spawned as Agent `haiku`. Two runs (M0+M1, M2+M3a; 1.7 and 0.7 min). No hand edits, no git writes. The controller rebuilt each sub-milestone from its snapshot with `commit_ms.sh`: M0 `897e6ce0` → 75, M1 `720824bc` → 905, M2 `86cedeae` → 1,510, M3a `ca2605d7` → 1,958. Every value equals the plan's, and the final `cmd`/`host` tree is `diff -r`-identical to the planner's round-trip tree.

**Controller gates:** plan §6 as one bash script at `ca2605d7`, BEFORE the judge: **rc 0.** `verify_ail.sh` (16 identities, 40 named tests, 503 se-tools tests); gofmt; build + vet; the compile fence; touched packages; `-race` on `host/approval/...`, `host/sshsign` and the cmd revoke/enrol tests (0 DATA RACE); the 16 named repo-wide inventories with `-v` (16 PASS, 0 SKIP/FAIL); `VERIFY_GO_LEGS=plain verify_go.sh` (62 ok, 0 FAIL). Details in `verification/world-iter379/controller-gates.md`.

**Judge:** sonnet via the Agent tool, on its own worktree at `ca2605d7` (generator ≠ judge). **Round 1: PASS 91/100, zero blocking** (correctness 28, tests 26, security 18, standards 9, scope 10).
- 77 mutants: 44 from plan §4 (all killed) and its own anchored to the diff; must-die control killed, must-live survived. Per-milestone non-vacuity held for M0, M1, M2 and M3a.
- **Survivor UV-OR-BE** (`if !ad.UV && !ad.BE`, `passkey.go:96`): the only no-UV arm cleared UV on Chrome's BE=0 flags. **The controller reproduced it** (mutant rc 0), then pinned UV and UP in every backup state with test-only commit `5a4b1224` (three arms) and proved the arms kill `!UV && !BE`, `!UV && !BS` and `!UP && !BE`, production restored byte-identical. UV-required is D-WORLD-89 R2, so this was closed now rather than carried.
- **Round 2** (judge resumed by SendMessage, at `5a4b1224`): confirmed, six UV/UP variants killed, **PASS 93/100**.
- Non-blocking, carried to sprint 2: `Passkey.Key`/`CredentialID` aliasing not pinned (copies today); `SOURCE.txt` names `es256.json`/`eddsa.json` while the files are `chrome-*.json`; `RegistrationInput.Challenge` zero/`Now` unchecked (caller's job). No Ed25519 point check (design-optional).

**Gate 3b:** PR #546, head `5a4b1224`.
- `pull_request` run 42917909921 attempt 1: `go host build + test gate` **FAIL** on `TestSeedRestartIgnoresASurvivingOldCopy/survivor` (`host/broker/exec_seed_test.go:656`, "the old staging dir survived the restart"). The sprint does not touch `host/broker` (`git diff --quiet 98374909 HEAD -- host/broker`, control: 18 approval files changed). It passed **10/10** locally at the sprint head, and it has 0 mentions in the charter, log and archive: **instance 1** of a CI flake, logged here and not filed as a row (D-WORLD-109 (2)). The job was re-run once (attempt 2): 3/3 success.
- **Instance 2, same iteration:** the record PR #547 (docs only, code identical to the green merge `d5d90e49`) went red in the `go race gate` job on the same `TestSeedRestartIgnoresASurvivingOldCopy/survivor` (`host/broker` FAIL at 298 s under `-race`; the two `DATA RACE` lines in that log are the known-positive race control `w-race-gate-blindspot/racecontrol`). Two CI instances in one evening, 0 local reds in 10.
- **Instance 3:** record PR #547 head `ed3700ee`, `go host build + test gate`, the same subtest. Locally at `d5d90e49`: 20/20 under `-race` and 30/30 at `-cpu 1,2` (0 reds in 60 local runs in all). At the third instance it becomes a row (D-WORLD-109 (2)): **row 259 `w-broker-seed-restart-survivor-ci-flake`**, filed below the release rows and row 201 under rule (d), with the three instances named and a mechanism hypothesis marked unmeasured. `host/broker`'s test binary links `host/approval` and `host/sshsign` (the AC10 census), but the seed test exercises neither, and the same code was green on the merge.
- The auto-close keyword scan of the PR body read 0 (control 1).
- Merged by merge commit **`d5d90e49`** (`--match-head-commit`). Push checks on the merge: 3/3 success.

**Routing evidence:** base=98374909268b9eab35382c5bf1779a3784edbf5b@2026-11-10T20:50:02Z.
- Controller: `claude-opus-5-5` (tok: not reported).
- Designer: `claude:claude-opus-5-5` via Agent `opus` (320,153 + 365,888 tok); r2 by controller carve-out.
- Quorum: r1 $0.4386 (108,047 in / 7,393 out), r2 $0.5997 (124,095 in / 14,972 out); absent_reviewers [] both rounds.
- Planner: `opus` via Agent (445,890 tok).
- Executor: `claude:claude-haiku-5-5` via Agent `haiku` (39,424 + 37,756 tok).
- Evaluator: `sonnet` via Agent (102,582 + 106,877 tok).
- The driver noted codex unusable at preflight (probe rc 75); no role needed it.
- Metered total: $1.04 (quorum only). Details in `verification/world-iter379/routing.md`.

**Progress:** 1.0 clauses 1–7 MET; goal unmoved. World can now verify a passkey registration with the standard library alone: a Chrome passkey's registration is parsed (bounded CBOR, COSE ES256/EdDSA), bound to the challenge, origin and rpID, required to carry user presence and user verification, and turned into an approver with a 128-bit enrolment code the operator will compare at the terminal. Nothing is reachable from a store, CLI or route yet; that is sprints 2 and 3.

**Ruled out:**
- **Forwarding the quorum's objections unmeasured:** each of the five was measured first; all were real.
- **Treating the 40-bit code as adequate because comparison is constant-time:** constant time does not stop an offline second-preimage search over a public input.
- **Carrying UV-OR-BE forward as a non-blocking note:** UV-required is a ruled property, and the pin cost one test-only commit plus a judge resume.
- **Attributing the CI red to the sprint:** untouched package, 10/10 locally, first instance; a re-run of the same job went green.
- **Accepting the judge's 906/1,511 numstat as a plan mismatch:** it counted `.gitattributes`; the plan's unit (`-- cmd host`) reads 905/1,510.

**Retro:**
- **Two quorum rounds of real objections on a fresh security design is the norm, not an outlier** (373, 375 and now 379 all went r0 → r1 → carve-out r2). The designer directive could pre-empt the recurring classes: secrecy claims about values derived from attacker-readable inputs, and browser/CSP premises left UNVERIFIED when a probe harness is already running.
- **The judge's best finding was a flag-state coverage gap** — a refusal tested in one combination of independent bits. Next planner directive: for every boolean-flag refusal, an arm with each other flag set.
- **A judge's numstat can differ by unit, not by error.** State the numstat path filter in the evaluator directive.

**Next:** row 201 A2 **sprint 2**: M3b (the assertion law, plan Appendix C) at its head, then M4 (credentials schema 6, `webauthn_credentials`) and M5 (the consume counter and authority), from the A2 doc without re-design. The planner must apply §9's sprint-2 fallback early (it grows by about 450 lines). Carry the judge's aliasing pins, the `SOURCE.txt` names and the zero-challenge check. Then A2 sprint 3 (M6 CLI, M7 ceremony page + QUICKSTART §16; Mark's real-passkey exit run is attended, §10), then A3–A4, then row 243.

## 380 — 2026-11-11 — row 201 A2 sprint 2 built: the WebAuthn assertion law, credentials schema 6 with the passkey store, the consume that binds the passkey and advances its counter in the mint's transaction, and authority's `webauthn:<name>`; judge PASS 93 [PRODUCT]

**Pick:** row 201 `w-world-approvals-by-key`, parent milestone A2, sprint 2, the route iteration 379 named. D-WORLD-109 = A applies: test keys and temp stores only.
- **Re-measured at pick** (`origin/dev` `3b77c294`):
  - No open PR from this account; no `iter380` worktree; no `A2 sprint 2` commit on `origin/dev`.
  - The A2 child design r2 has a quorum artifact (`mc_quorum_gate2 w-world-approvals-by-key-a2` rc 0), so no pick-time quorum was owed and no designer ran.
- **Ledger:** valid (100 rows). OPEN: D-WORLD-105 only.
- **Clause map:** 1.0 clauses 1–7 MET, none UNMET. Release rows are not routable:
  - Row 150's attended A1: 0 subjects since 2026-11-07 matching `row ?150|A1\b` (control `A1b`: 20).
  - Row 151's row-188 N2d/N4a: 0 since 2026-10-01 (control N2c: 1).

**Gate 0:**
- Stamped first. `gh` is `sunholo-voight-kampff`. Billing CLEAN. Kill switch absent.
- Directives: 0 allowlisted on #524, out of 12 comments. The issue-scoped watermark is absent, so all comments were read. A hand read of #524 found no `FAILED to complete` notice.
- Self-notices: rc 2, an instrument failure and not a verdict; no `-last-seen` file exists for 524 or 449 (same as 371–379).
- The `mission-world` inbox has 0 unread.
- Skill drift: none. SKILL.md and all 13 resource files at the resolved symlink target match the driver clone's `origin/dev`.
- Dev at `3b77c294` (record 379): `checks=3`, all success. Local dev equal to origin.

**Plan:** `derive-planner-lane.sh` → `opus fail-closed:env-pin`; Agent `opus`, 93.8 min, 478,824 tok. Committed `75cb07e9`.
- **Sprint 2 = M3b1–M3b3, M4a–M4c, M5a–M5b = 1,843 lines** (`-- cmd host`, cap 2,000). The design's §9 sprint-2 fallback was priced as a measurement and did not need to fire: M3b measured 890 (sprint 1 estimated 450), M4 648, the consume 239, and A6 + AU3 only 66, so they stay. Appendix C moves nothing to sprint 3.
- Mutants: 139, each run against its named test and its package; 135 killed, 4 named equivalents (decoded inputs that cannot alias). Two tests strengthened after the first pass found gaps.
- Round trip: the patches extracted from the plan file, applied to a fresh export of `3b77c294`, reproduce the mutant tree byte for byte. **The controller re-did it** (its own export and awk): `diff -rq cmd host` empty.
- **Premises found false or incomplete:** P1 the size estimates; P2 schema 6 broke four `cmd/ailang-worldd` test fixtures that rebuild a v4 store by dropping `approvers` only (they now also drop `webauthn_credentials`; test-only, in M4a), though the design said sprint 2 touches no `cmd/` file; P3 `EnrolPasskey`'s validation is ten checks, not nine (S1–S10); P4 the `user_version 6` check compared the store with the constant itself and could never fail (now the literal 6); P5 the test authenticator's `name` parameter dropped (the recordings carry no user).
- **Carry-forwards from 379's judge:** aliasing pinned (`TestNewPasskeyApproverDoesNotAliasItsInput`, PA4); `SOURCE.txt` names the `chrome-*.json` files and XW0 checks the names; zero challenge / `Now` measured accepted by the law and left as the caller's job (`EnrolPasskey` refuses `EnrolledAt <= 0`; sprint 3's random challenge owns the challenge).
- A harness trap the planner found: any gate started under `nohup` reds `scripts/marketingcheck`'s `shell-hup` test (3/3). The plan forbids `nohup`.

**Executor:** env `recipe claude:claude-haiku-5-5`, spawned as Agent `haiku`. Four runs of two sub-milestones (1.5, 5.0, 1.3 and 1.4 min). No hand edits, no git writes, no stop condition. The controller rebuilt each sub-milestone from its snapshot with `commit_ms.sh`: M3b1 `8c68004c` → 345, M3b2 `5a1c9383` → 672, M3b3 `6a2984c2` → 890, M4a `0f0319d9` → 1,027, M4b `5ca2f4dc` → 1,300, M4c `17039568` → 1,538, M5a `19039a40` → 1,777, M5b `1fe7c86d` → 1,843. Every value equals the plan's, and the final `cmd`/`host` tree is `diff -r`-identical to the planner's mutant tree.

**Controller gates:** plan §6 as one bash script at `1fe7c86d`, BEFORE the judge: **rc 0.** `verify_ail.sh` (16 identities, 40 named tests, 503 se-tools tests); gofmt; build + vet; the compile fence; touched packages (6 ok); `-race` on `host/approval/...`, `host/store`, `host/authority` (0 DATA RACE); 22 named inventories and reflective tests with `-v` (22 PASS, 0 SKIP/FAIL); `VERIFY_GO_LEGS=plain verify_go.sh` (62 ok, 0 FAIL; its two DATA RACE lines are the race control). Details in `verification/world-iter380/controller-gates.md`.

**Judge:** sonnet via the Agent tool, on its own worktree at `1fe7c86d` (generator ≠ judge). **Round 1: PASS 93/100, zero blocking** (correctness 28, tests 29, security 19, standards 9, scope 8).
- 66 of its own diff-anchored mutants (every conjunct of the W7/W8/W12/W13 flag refusals, the counter law's four forms, S1–S10, the lookups, C13–C18, `usePasskey`, A6), must-die control killed, must-live control survived; 48 plan rows sampled across all 8 sub-milestones, no unexpected survivor.
- Three of its mutants (P1, C17, A6-count) did not compile on the first run; it re-ran them compiling and reported them killed, but the re-runs are not in its `mutants.txt`. **The controller reproduced all three** with `-vet=off`: killed by `TestApproveByAssertionOrder`, `TestPasskeyConsumeChecksTheRowInTheTransaction` and `TestMintForRequestDerivesPasskeyUse`.
- Non-vacuity: M3b1, M4b and M5a's new tests do not exist at the parent commits.
- Non-blocking: N2, `ApproveByAssertion` does not guard a zero-valued `Passkey.Key`. The controller read `COSEKey.verify`: its default arm returns `ErrCOSEAlg`, so it fails closed. Logged, not filed (D-WORLD-109 (2)).

**Gate 3b:** PR #548, head `1fe7c86d`: 3/3 success on the first run. Mergeable CLEAN. The auto-close keyword scan of the PR body read 0 (control 1). Merged by merge commit **`3e55451f`** (`--match-head-commit`). Push checks on the merge: 3/3 success (run 42952060607).

**Routing evidence:** base=3b77c294000ca40c2fe77e312ae0e8ce0f71a030@2026-11-11T00:50:05Z.
- Controller: `claude-opus-5-5` (tok: not reported).
- Designer: not spawned (the A2 r2 doc covers sprint 2). Quorum: not re-run (artifact present, doc unchanged).
- Planner: `opus` via Agent (478,824 tok).
- Executor: `claude:claude-haiku-5-5` via Agent `haiku` (46,214 + 49,619 + 40,775 + 44,161 tok).
- Evaluator: `sonnet` via Agent (73,582 tok).
- The driver noted codex unusable at preflight (probe rc 75); no role needed it.
- Metered total: $0.00. Details in `verification/world-iter380/routing.md`.

**Progress:** 1.0 clauses 1–7 MET; goal unmoved. World can now verify a passkey *assertion* (an approval) as well as a registration, store enrolled passkeys in credentials schema 6, and consume a mint approved by a passkey: the consume re-reads the passkey row inside the mint's transaction, binds the credential ID, refuses a counter that does not strictly advance (both zero allowed), and records `approved_by = webauthn:<name>`, derived by authority from the verified approval. No CLI or route builds a passkey approval yet; that is sprint 3.

**Ruled out:**
- **Cutting AU3/A6 to sprint 3 pre-emptively** (sprint 1's plan said to apply the fallback early): measured at 66 lines with 157 lines of margin; moving them would leave the store accepting a passkey consume no caller builds.
- **Banking the judge's "re-ran and killed" for its three compile-error mutants:** the re-runs were not on disk; the controller reproduced them instead.
- **Treating N2 as a gap:** a zero `COSEKey` reaches `verify`'s default arm and is refused.

**Retro:**
- **The planner-ships-patches route held for an eight-patch sprint across three packages:** zero executor deviations, numstat equal at every step, and the controller's own round trip matched. Fourth consecutive sprint on this route.
- **A judge's narrative claim about a re-run is not evidence until it is on disk.** Its `mutants.txt` recorded `COMPILEERR` for three mutants its report called killed. Next evaluator directive: every mutant re-run must append its own line to the mutant file.
- **A design's "touches no file in X" is a claim the planner should measure:** schema 6 reached `cmd/` through test fixtures that rebuild an older schema by dropping tables.

**Next:** row 201 A2 **sprint 3**: M6 (the CLI: `approver enroll --format webauthn`, `approve --webauthn-link`/`--webauthn-assertion`, list/revoke show passkeys) and M7 (the static `/surface/passkey` ceremony page, the Chrome driver and `generate.sh`, QUICKSTART §16), from the A2 doc without re-design; the design's sprint-3 fallback makes M7 its own fourth sprint if needed. Consider sprint 3 refusing an all-zero challenge in its random-challenge step. Then Mark's attended real-passkey exit run (§10), A3–A4, then row 243.

## 381 — 2026-11-11 — row 201 A2 sprint 3 built: the passkey operator CLI (`approver enroll --format webauthn` with the enrolment-code fence, list/revoke labels, `approve --webauthn-link`/`--webauthn-assertion`); the design's sprint-3 fallback fired, so the ceremony page is sprint 4; judge PASS 92 [PRODUCT]

**Pick:** row 201 `w-world-approvals-by-key`, parent milestone A2, sprint 3, the route iteration 380 named. D-WORLD-109 = A applies: test keys and temp stores only.
- **Re-measured at pick** (`origin/dev` `ab4b2402`):
  - No open PR from this account; no `iter381` worktree; no `A2 sprint 3` commit on `origin/dev`.
  - The A2 child design r2 has a quorum artifact (`mc_quorum_gate2 w-world-approvals-by-key-a2` rc 0), so no pick-time quorum was owed and no designer ran.
- **Ledger:** valid (100 rows). OPEN: D-WORLD-105 only.
- **Clause map:** 1.0 clauses 1–7 MET, none UNMET. Release rows 150/151 not routable (unchanged since 380; no attended A1 or row-188 N2d/N4a work landed in between).

**Gate 0:**
- Stamped first. `gh` is `sunholo-voight-kampff`. Billing CLEAN. Kill switch absent.
- Directives: 0 allowlisted on #524, out of 13 comments. The issue-scoped watermark is absent, so all comments were read.
- Self-notices: rc 2, an instrument failure and not a verdict; no `-last-seen` file exists for 524 or 449 (same as 371–380).
- Inbox: no unread `mission-world` message addressed to this loop (the unread list holds other missions' controlplane fyi notes).
- Skill drift: none. SKILL.md and all 13 resource files at the resolved symlink target match the driver clone's `origin/dev` (`de918c252`).
- Dev at `ab4b2402` (record 380): `checks=3`, all success. Local `dev` was behind origin (record 380 merged via GitHub); state read from origin.
- Weekly issue sweep: not owed (not the first fire after Monday's rotation).

**Plan:** `derive-planner-lane.sh` → `opus fail-closed:env-pin`; Agent `opus`, 181 min, 536,858 tok. Committed `3cff9784`.
- **The design's §9 sprint-3 fallback fired, on measurement.** The planner built M6 in full and measured it at **1,686 lines** (`-- cmd host`: 467 production, 1,219 test) against the design's 1,220 estimate (1.38×). 314 lines of margin cannot hold M7, whose three browser probes alone measured 202 lines before any page or Go code. M7 (`/surface/passkey` + `passkey.js`, PG1, XB1, XB2, `generate.sh`) and QUICKSTART §16 move whole to a fourth sprint, listed in the plan's Appendix C. The printed links point at `/surface/passkey`, a 404 until then.
- Eight sub-milestones M6a–M6h. Mutants: 124 rows, 121 killed against the named test and the package, 3 named equivalents; 3 base controls; a 146-line per-run ledger (17 first drafts that did not compile were rewritten and re-run, each with its own line).
- Round trip: the planner's eight patches, extracted from the plan file, rebuild its mutant tree byte for byte. **The controller re-did it** (its own `git archive ab4b2402` and awk): `diff -rq cmd host` empty.
- For sprint 4: XB2's virtual-time approach measured — `Emulation.setVirtualTimePolicy` advanced the page's 240000 ms abort timer to exactly 240000 ms of page time in 53–55 ms real time, 3 of 3 runs, under the exact `surfaceCSP` (headless Chrome, temp profile, no authenticator).
- **Premises found false or incomplete:** P1 no exported function turns a stored COSE key back into a key (added `host/approval.ParseStoredCOSEKey`, pure, with a test; outside the design's conflict surface); P2 the size estimate; P3 `host/store/context_roots_test.go` (+1 root) is also outside the conflict-surface table; P4 the enrol-refusals arm the design expected to edit stays green unedited; P6 a first CW4 draft flaked ~1 in 6 and was fixed before shipping; P7 WF5 also needs a refusal for a stored key that will not parse; P8 QUICKSTART §16 belongs with M7.
- **Carry-forward decided:** WN5 refuses a nil, failed **or all-zero** challenge with `exitFatal` before any terminal (the registration law measured accepting a zero challenge; the caller now owns it).

**Executor:** env `recipe claude:claude-haiku-5-5`, spawned as Agent `haiku`. Four runs of two sub-milestones (8.6, 5.2, 2.6 and 13.6 min). No hand edits, no git writes, no stop condition; one timing exception used in run 1 (`TestDemoSeedPreflightRefusals` hit its 10 s `ailang --version` deadline under load in the whole-package run and passed alone, 38/38 subtests). The controller rebuilt each sub-milestone from its snapshot with `commit_ms.sh`: M6a `67f4d9af` → 38, M6b `41f5cb9a` → 237, M6c `e7241f61` → 701, M6d `48d75b44` → 945, M6e `62f57011` → 1,190, M6f `0ba3a76d` → 1,397, M6g `f4655266` → 1,538, M6h `64a239f8` → 1,686. Every value equals the plan's, and the final `cmd`/`host` tree is `diff -rq`-identical to the planner's mutant tree.

**Controller gates:** plan §6 as one bash script at `64a239f8`, BEFORE the judge.
- Pre-commit legs green (`verify_ail.sh`, gofmt, build, vet, the compile fence, touched packages, `-race`, all 30 named inventories with `-v`, 0 SKIP/FAIL).
- `VERIFY_GO_LEGS=plain verify_go.sh`: **rc 1 under machine load ~32** on four timing tests — `TestHostsLaunchersExecHost` (`cmd/ailang-worldd`, `codex: signal: killed`) and three `host/broker` exec-probe timeouts. **All four pass alone at the head**; `host/broker` imports no touched package; the planner measured the same class red at the base under the same load. Details in `verification/world-iter381/controller-gates.md`.

**Judge:** sonnet via the Agent tool, on its own worktree at `64a239f8` (generator ≠ judge). **Round 1: PASS 92/100, zero blocking** (correctness 28, tests 26, security 19, standards 9, scope 10).
- A 67-line mutant ledger on disk (`verification/world-iter381/judge-mutants.tsv`): 35 plan rows sampled across M6a/b/c/f/g, all killed except the two named equivalents; 23 of its own diff-anchored mutants; must-die controls killed, must-live control survived. Its one COMPILEERR mutant was rewritten and re-run **with its own ledger line** — iteration 380's retro directive held.
- Non-vacuity per milestone: M6c/M6f/M6g's new tests do not compile or fail at the parent; deleting M6d/M6e's test files lets five mutants survive (WN10 invert and prefix, WN11, N12, the y/N gate), so those test-only milestones are what pin CW2/CW4/CW5/CW6.
- The cut judged safe: the daemon serves only exact `/surface…` paths, so the link is a 404; enrol then refuses after 300 s and enrols nothing.
- Non-blocking, all fail closed: the next-step line's `--db`/`--credentials-db` order is unpinned; `defer ro.Close()` in the pre-check is unobserved; `_` as an extra code separator survives (the compare still needs the genuine 26 characters); `challenge[0]==0` in place of the all-zero check survives (stricter). Logged, not filed (D-WORLD-109 (2)).
- **The controller reproduced one kill and one survivor:** dropping WN5's all-zero check fails `TestPasskeyEnrollRaceAndRefusals/WN5`; swapping the next-step `--db` pair leaves `TestApproveByPasskey` green.

**Gate 3b:** PR #550, head `64a239f8`. Attempt 1: 2/3 — the race job red on `scripts/marketingcheck` `TestMarketingSupervisorFakeDaemonCleanup/shell-int` (trap `signal: killed`, cleanup 29.93 s over its bound; run 42993721120). Known row-254 class (folded into row 258), a new subtest (`shell-int`; earlier instances were `shell-wait` and `shell-hup`); this sprint touches no `scripts/` file; its two `DATA RACE` lines are the gate's positive control. `gh run rerun --failed` → **3/3 success**. Mergeable CLEAN. The auto-close keyword scan of the PR body read 0 (control 1). Merged by merge commit **`0ce117c5`** (`--match-head-commit`). Push checks on the merge: 3/3 success (run 42997956004).

**Routing evidence:** base=0ce117c56f44a413e5206143413d269e8ede3c4d@2026-11-11T09:46:27Z.
- Controller: `claude-opus-5-5` (tok: not reported).
- Designer: not spawned (the A2 r2 doc covers sprint 3). Quorum: not re-run (artifact present, doc unchanged).
- Planner: `opus` via Agent (536,858 tok).
- Executor: `claude:claude-haiku-5-5` via Agent `haiku` (49,537 + 45,854 + 40,458 + 50,472 tok).
- Evaluator: `sonnet` via Agent (105,937 tok).
- The driver noted codex unusable at preflight (probe rc 75); no role needed it.
- Metered total: $0.00. Details in `verification/world-iter381/routing.md`.

**Progress:** 1.0 clauses 1–7 MET; goal unmoved. World now has the operator half of passkey approvals: an operator can enrol a passkey (`approver enroll --format webauthn`, fenced by a 128-bit code compared in constant time and a typed `y`), see passkeys in `approver list`/`revoke`, print a ceremony link for a pending request, and mint from a saved assertion file with no World prompt. What is missing is the page that performs the browser ceremony; until sprint 4 the links 404, so nothing can be enrolled or approved by passkey yet.

**Ruled out:**
- **Keeping M7 in sprint 3:** M6 alone measured 1,686 of 2,000; the design's own fallback names this cut.
- **Treating the four `verify_go` timing reds as regressions:** each passes alone at the head, three are in a package that imports nothing touched, and the base showed the same class under the same load.
- **Treating the PR's race red as this sprint's:** it is in `scripts/marketingcheck`, untouched; the row-254 class; green on re-run.
- **Filing the four fail-closed survivors as rows:** D-WORLD-109 (2) — logged here.

**Retro:**
- **The planner-ships-patches route held a fifth consecutive sprint:** zero executor deviations, numstat equal at every step, the controller's own round trip matched.
- **The judge's ledger directive worked:** its one recompiled mutant has two lines on disk, and the controller could reproduce a kill and a survivor straight from the ledger rows.
- **A design's size estimate keeps undershooting CLI work by ~1.4×** (A2 sprint 2 M3b 2.0×, sprint 3 M6 1.38×); the §9 fallback-as-measurement rule is what kept both sprints under the cap without a re-plan.
- **The rig was loaded all fire (load 12–38 from other sessions):** the planner spent 181 min, and both the planner and the controller saw timing reds in untouched packages. A gate verdict under shared load needs the alone re-run and a base control; both were done.

**Next:** row 201 A2 **sprint 4**: M7 from the sprint-3 plan's Appendix C (the static `/surface/passkey` page + `passkey.js` with the [r2] AbortController/240000 ms timer, PG1, live-Chrome XB1/XB2 on a CDP virtual authenticator, `generate.sh`, QUICKSTART §16 tagged attended-pending), from the A2 doc without re-design. Consider the judge's one-line early `--name` check (Appendix C) and pinning the next-step `--db` order. Then Mark's attended real-passkey exit run (§10), A3–A4, then row 243.

## 382 — 2026-11-11 — row 201 A2 sprint 4 built: the passkey ceremony page under the unchanged surfaceCSP, live-Chrome XB1/XB2 on a CDP virtual authenticator, `generate.sh` made real, QUICKSTART §16 (attended — pending); judge PASS 91; NOT LANDED — the CI host job is red on a runner apt-mirror DNS outage that a `dev` control reproduces [PRODUCT]

**Pick:** row 201 `w-world-approvals-by-key`, parent milestone A2, sprint 4 — M7 from the sprint-3 plan's Appendix C, the route iteration 381 named. D-WORLD-109 = A applies: test keys, temp stores and a CDP virtual authenticator only.
- **Re-measured at pick** (`origin/dev` `27ce163e`): no open PR from this account; no `iter382` worktree; no `A2 sprint 4`/M7 commit on `origin/dev`; `/surface/passkey` not served at HEAD (`generate.sh` still the sprint-1 stub).
- The A2 child design r2 has a quorum artifact (`mc_quorum_gate2 w-world-approvals-by-key-a2` rc 0), so no pick-time quorum was owed and no designer ran.
- **Ledger:** valid (100 rows). OPEN: D-WORLD-105 only.
- **Clause map:** 1.0 clauses 1–7 MET, none UNMET. Release rows 150/151 not routable (unchanged since 381).
- **Orphans:** iteration 381's fire exited rc=143 after its record landed (the driver's own notice on #524 says so; record #551 is on `origin/dev`), so nothing was orphaned. **Iteration 381's record PR #551 was red once on `TestMarketingCaptureFakeChromeSixShots/inventory` (row 247 class, 120 s ctx) and green on one `rerun --failed`**; that is not in log 381 and is recorded here as 381 asked.

**Gate 0:**
- Stamped first. `gh` is `sunholo-voight-kampff`. Billing CLEAN. Kill switch absent.
- Directives: 0 allowlisted on #524, out of 15 comments. The issue-scoped watermark is absent, so all comments were read.
- Self-notices: rc 2, an instrument failure and not a verdict; no `-last-seen` file exists for 524 or 449 (same as 371–381). The one driver notice since 381 is its `rc=143 after landing its record` note.
- Inbox: no unread `mission-world` message addressed to this loop (the two newest are this mission's own 381 fyi and rc=143 notes).
- Skill drift: none. SKILL.md and all 13 resource files at the resolved symlink target match the driver clone's `origin/dev`.
- Dev at `27ce163e` (record 381): `checks=3`, all success. Local `dev` was 12 behind origin (records merged via GitHub); state read from origin.
- Weekly issue sweep: not owed (not the first fire after Monday's rotation).

**Plan:** `derive-planner-lane.sh` → `opus fail-closed:env-pin`; Agent `opus`, 218 min, 479,457 tok. Committed `d688d22d` on the sprint branch.
- Seven sub-milestones M7a–M7g, **1,054 lines** `-- cmd host` (production ≈ 206, the rest tests and the CDP driver); the 2,000 cap held with 946 to spare, nothing narrowed. M7e is 327 lines (over the ~250 aim) because the routes may not land without both browser tests.
- Both optional carry-forwards taken in M7a: WN2a refuses a `--name` outside `approval.ApproverNamePattern` before any link, and CW7 pins the next-step line's `--db … --credentials-db …` order.
- Mutants: 69 (64 killed by the named test and the package's passkey tests, 5 named equivalents); a 251-line per-run ledger.
- **Premises found false:** P1 the design puts XB1/XB2 in `host/daemon`, but `ceremonyLink` is in `package main`, so the browser tests and driver live in `cmd/ailang-worldd` while the page is still served through `daemon.New(...).Handler()`; P2 `generate.sh` cannot reproduce the committed fixtures byte-for-byte (new keys each run), so it writes only into an empty directory it is given and its output passes the fixture law; P3 Chrome's virtual authenticator always picks ES256, so EdDSA through the page is a named equivalent. Six browser races were fixed in the driver/test, none by loosening an assertion (e.g. virtual time is now paused before the click so the abort lands at exactly 240000 ms under load; a UV-false authenticator replaces `setUserVerified(false)`).
- **Controller decisions on the planner's asks:** (1) page/JS mutants run against the package's passkey tests rather than the whole `cmd/ailang-worldd` — accepted on the planner's C22 grep and the judge's independent check that only those tests reach the page bytes; (2) `-timeout 20m` on whole-package runs — accepted; (3) `TestWorldContextAndStatusline` added to the known load-sensitive list — accepted.
- Round trip: **the controller re-did it** (its own `git archive 27ce163e` → the 7 patches awk-extracted from the plan → `patch -p1`): `diff -rq` against the planner's mutant tree empty for `cmd`, `host`, `docs`.

**Executor:** env `recipe claude:claude-haiku-5-5`, spawned as Agent `haiku`. Two runs (M7a–M7d 7.5 min, M7e–M7g 7.8 min), merging the plan's four runs to save wall clock. No hand edits, no git writes, no stop condition, no timing re-run. It reported the pre-existing env-gated skips in whole-package `-v` runs (14 in `cmd/ailang-worldd`, ~30 in `host/daemon`) rather than treating them as M7 skips; each is unchanged at the base. The controller rebuilt each sub-milestone from its snapshot with `commit_ms.sh`: M7a `e4734707` → 11, M7b `a7f45088` → 195, M7c `30c5debe` → 367, M7d `c1483836` → 597, M7e `17c21dae` → 924, M7f `e8077ba6` → 1,054, M7g `e07b7bad` → 1,054. Every value equals the plan's, and the final `cmd`/`host`/`docs` tree is `diff -rq`-identical to the planner's.

**Controller gates:** plan §6 as one bash script at `e07b7bad`.
- Pre-commit script **green** (`PRE_COMMIT_GATES_GREEN`): verify_ail, gofmt, build, vet, compile fence, touched packages, `-race` legs, 29 named inventories/help/quickstart tests `--- PASS` (0 SKIP/FAIL), and **XB1+XB2 required ×3 in a row: 9, 13, 15 s** (headless Chrome 157, temp profile, CDP virtual authenticator).
- `VERIFY_GO_LEGS=plain verify_go.sh`: rc 1, **run concurrently with the judge's mutants** (the slot clock forced the overlap), on three package-level `test timed out after 10m0s` panics (host/broker, host/daemon, host/verifygate). Re-run alone with `-timeout 20m`: host/daemon rc 0 (381 s), host/broker rc 0 (256 s), host/verifygate red once on `TestMissionConfigDefaultFlowIntegration/valid_ledger` (plan's known list) and green alone. Details in `verification/world-iter382/controller-gates.md`.

**Judge:** sonnet via the Agent tool, on its own worktree at `e07b7bad` (generator ≠ judge). **Round 1: PASS 91/100, zero blocking** (correctness 28, tests 27, security 19, standards 9, scope 8).
- A 46-line mutant ledger (`verification/world-iter382/judge-mutants.tsv`): 38 of its own diff-anchored mutants across page checks (version, each op, origin, rp, base64url, JSON, timer, late result, expiry, abort signal, button disable, save-after-failure, code label/order, UV), routes/headers, WN2a and `generate.sh`; 34 killed, must-die control killed, must-live control survived. XB1/XB2 required twice (XB1 6.4/3.3 s, XB2 4.1/4.0 s).
- Non-vacuity per milestone: removing M7e's route registrations fails PG1 and XB1/XB2; page mutants kill XB1/XB2; removing WN2a fails CW1's arm; each `generate.sh` refusal removal fails `TestPasskeyFixtureGeneratorRefuses`. It confirmed XB1 drives the CLI's own printed links into the real enrol and `--webauthn-assertion` paths.
- Non-blocking, all fail closed or cosmetic: `J-allow-empty` survives (the virtual authenticator holds one credential, so the allow list is unexercised — a second enrolled credential in XB1 would kill it); **`R-csp` survives** (PG1 compares the header with the `surfaceCSP` identifier, so no test pins the CSP's value); `J-blob-type` (saved MIME type unasserted); `H-go-enabled` (redundant initial `disabled`). Logged, not filed (D-WORLD-109 (2)).
- **The controller reproduced one kill and one survivor:** `if false && …` on WN2a fails `TestPasskeyEnrollFlagRules`; adding `blob:` to `surfaceCSP`'s `style-src` leaves `TestPasskeyPageRoutes` and every host/daemon CSP/surface test green.

**Gate 3b:** PR #552, head `e07b7bad`, MERGEABLE/CLEAN.
- Attempt 1: 1/3 — the race job red on `scripts/marketingcheck` `TestMarketingSupervisorFakeDaemonCleanup/shell-hup` (row 254/258 class; this sprint touches no `scripts/`), and the host job red at `apt-get install -y -qq bubblewrap socat ripgrep` (`Temporary failure resolving 'azure.archive.ubuntu.com'`, exit 100) before any Go test.
- Attempts 2, 3 and 4 (`gh run rerun --failed`): the race job green from attempt 2; the host job red every time at the same apt step, 7 DNS-failure lines each.
- **Negative control:** a `workflow_dispatch` of CI on `dev` `27ce163e` (all 3 green at 10:34Z) failed its host job at the same step (`Row 140 srt confinement matrix`, apt DNS failures, the step timed out at 8 min). GitHub's status API read "All Systems Operational". So the red is a runner-mirror outage, inherited, not this sprint.
- **Not merged** (`dev` has no branch protection, so the merge button works; merging over a red is forcing a guardrail). This record is committed on PR #552's branch so the next fire finds it there (Gate 1's blocked-PR rule).

**Routing evidence:** base=27ce163eb1ae8af55151ea6548591c15a3126835@2026-11-11T10:54:18Z.
- Controller: `claude-opus-5-5` (tok: not reported).
- Designer: not spawned (the A2 r2 doc covers sprint 4). Quorum: not re-run (artifact present, doc unchanged).
- Planner: `opus` via Agent (479,457 tok).
- Executor: `claude:claude-haiku-5-5` via Agent `haiku` (69,540 + 51,236 tok).
- Evaluator: `sonnet` via Agent (77,943 tok).
- The driver noted codex unusable at preflight (probe rc 75); no role needed it.
- Metered total: $0.00.

**Progress:** 1.0 clauses 1–7 MET; goal unmoved. Once PR #552 merges, A2 is code-complete: an operator can enrol a passkey from the printed link in any browser on `localhost`/`*.ts.net`, and approve a pending request with it, with no World terminal prompt. The page itself is now proven end to end on a virtual authenticator; Mark's attended run on a real passkey (§10, QUICKSTART §16) is what remains of A2.

**Ruled out:**
- **Treating the host-job red as this sprint's:** it fails at an apt install before any Go test, 4 of 4 times, and the same step fails on `dev` `27ce163e`, which was green this morning.
- **Merging over it:** `dev` is unprotected, but an item is not LANDED until remote CI passes on its merge.
- **Treating the `verify_go` timeouts as regressions:** all three were package-level 10-minute timeouts under concurrent judge load; each package is green alone (verifygate after one known-list re-run); two of the three import nothing touched.
- **Filing the judge's four survivors as rows:** D-WORLD-109 (2) — logged here.

**Retro:**
- **The planner-ships-patches route held a sixth consecutive sprint:** zero executor deviations, numstat equal at every step, the controller's round trip matched. Two executor runs of 3–4 sub-milestones each were as clean as four runs of two.
- **The planner over-ran its stated budget by 1.8×** (218 min against ~120). It spent the time repeating mutant runs after load-driven test fixes and a Go build cache another session trimmed mid-run; with a 6 h slot that left ~2 h for execute, gate, judge, CI and record, which forced the gate sweep to overlap the judge and produced three false timeouts. If the next sprint carries browser mutants, give the planner a hard mid-run checkpoint (e.g. ship the patches first, mutants second).
- **An outage that hits a pre-test install step is invisible to every local gate**, and repeated `rerun --failed` cannot clear it. The `dev` dispatch control is what turned "CI is flaky" into "CI is down at this step for every commit" in one run.

**Next:** land PR #552: `gh run rerun --failed` on its newest head until all 3 checks pass (the predicate is the host job's apt step completing), then merge `--match-head-commit` and confirm 3/3 push checks on the merge. Then Mark's attended real-passkey exit run (A2 §10, QUICKSTART §16). Then row 201 A3 (the surface flow) needs a new child design: the rotation file's last-used designer is `claude:claude-opus-5-5`, so the next entry is `codex:gpt-6.1-sol` if its lane is usable. Consider pinning `surfaceCSP`'s value in a test (judge R-csp) in that sprint.

## 383 — 2026-11-11 — row 201 A3 child design r2 (`w-world-approvals-by-key-a3.md`): quorum blocked twice, controller carve-out, independent audit VALID 84; D-WORLD-114 filed on its two amendments to the parent; row 201 A2 sprint 4 LANDED (PR #552, merge `88ff1565`) once the CI runner apt outage cleared on attempt 8 [PRODUCT]

**Pick:** carry-over first: land PR #552 by iteration 382's predicate. While the outage persisted, ci-health's rule for a provider-side red applied: do not fix-forward or revert, and pick routable work that needs no landing gate. That was row 201 A3's child design, which the charter's Next names after Mark's attended A2 exit run (attended, not loop work). The outage cleared on attempt 8, and #552 landed in the same fire.
- **Re-measured at pick** (`origin/dev` `27ce163e`, record 381; PR #552's branch carries record 382, read as Gate-1 input): no `a3` doc, branch, PR or worktree; `mc_quorum_gate2 w-world-approvals-by-key-a3` rc 1.
- **Ledger:** valid (100 rows before, 101 after). OPEN: D-WORLD-105; D-WORLD-114 filed this iteration.
- **Clause map:** 1.0 clauses 1–7 MET, none UNMET. Release rows 150/151 not routable (unchanged).
- **Orphans:** none. Iteration 382's fire ended normally; its fyi arrived 3 minutes before this fire.

**Gate 0:**
- Stamped first. `gh` is `sunholo-voight-kampff`. Billing CLEAN. Kill switch absent.
- Directives: 0 allowlisted on #524, of 16 comments. The watermark file is absent, so all comments were read.
- Self-notices: rc 2, an instrument failure and not a verdict; no `-last-seen` file exists for 524 or 449 (same as 371–382).
- Inbox: nothing unread addressed to `mission-world` (`-inbox mission-world` empty).
- Skill drift: none. SKILL.md and all 13 resource files at the resolved symlink target match the driver clone's `origin/dev`.
- Weekly issue sweep: not owed.

**Gate 1 / 3b — PR #552 (head `0140bd7f`):** `rerun --failed` attempts 2–8 this fire (16:06–17:00Z). The host job was red at `apt-get install -y -qq bubblewrap socat ripgrep` (exit 100) on attempts 2–7, with 6–7 `Temporary failure resolving 'azure.archive.ubuntu.com'` lines each, before any Go test ran. **Attempt 8** (17:00Z): the apt step passed and every later step ran; the host job was green, so all **3/3** checks passed on `0140bd7f` (`commits/<sha>/check-runs` checks=3, all success). `ailang-code verify` and `go race` were green throughout.
- **Landed:** `gh pr merge 552 --merge --match-head-commit 0140bd7f…` → merge **`88ff1565`** at 17:11:12Z (no closing keywords in the PR text). Push CI on the merge (run 43065674075): 3/3 success (run 43065674075, 17:26Z).
- **Controls:**
  - the GitHub status API: `none` / All Systems Operational, no unresolved incidents;
  - from the rig, `azure.archive.ubuntu.com` resolves (a traffic-manager CNAME to a Sweden Central load balancer);
  - the 382 `dev` dispatch control already failed at the same step.
- **Reading:** the runner-side Ubuntu mirror is broken, which makes this a provider-side red. ci-health: re-run, record the diagnosis, do not fix-forward, do not revert.
- **Not done:** pointing the apt step at another mirror would have been a fix-forward of a provider outage; it cleared on its own about 3 hours after it was first seen.

**Design** (`w-world-approvals-by-key-a3.md`, a child of the attended parent). Resolver: `recipe claude:claude-opus-5-5 declared:provider-pin`; spawned as Agent `opus`, following the 379 precedent. The rotation file reads `claude:claude-opus-5-5`; codex was unusable at preflight (rc 75). Worktree `.wt-world-iter383-a3`, branch `design/world-iter383-row201-a3`, stacked on #552's head.
- **r0** (12.2 min, 798 lines). Split:
  - **A3a**, the web flow, in three sprints:
    - credentials schema 7 `surface_approvals`;
    - one shared judge;
    - the CLI minter `approve <id> --await`;
    - `GET /surface/requests.json`, `GET /surface/requests/{id}`, the card page and `POST /surface/approve/{id}`, with Host, Origin, `Sec-Fetch-Site` and JSON checks;
    - *Needs you*, the decision list and an opt-in notification.
  - **A3b**, the Tauri notification and Touch ID: constraints only, because no Tauri app exists.
  - The designer measured that minting is CLI-side (`session_new.go:404`) and that the daemon never keeps the credentials writer (`host/daemon/credentials.go:50-53`). **The controller re-read both, plus the absence of `desktop/`, first-party.** It also probed Chrome's POST headers and the cross-origin preflight (V17).
- **Quorum round 1 BLOCKED 3/3** ($0.385). Author Anthropic benched; seated oc-glm-5-3, oc-kimi-k3, gpt6-1-sol; absent none. The controller checked each objection against r0, and all three were real:
  - R3 blocked SSHSIG behind a tailnet origin;
  - the `crypto.subtle` premise was cited as measured but was not;
  - the card omitted signed fields.
- **r1** (designer resumed by SendMessage, 4.7 min, 928 lines):
  - a surface origin set with `serve --surface-origin` (RT13, CF1);
  - V23–V26 measured in Chrome 157 at `http://localhost` under `surfaceCSP` (digest, the 90-byte message and `MessageHash` byte-equal to Go on two golden requests; `Notification` behaviour);
  - a `FIELDS` table rendering every `MintRequest` key path, with PG3.
- **Quorum round 2 BLOCKED 3/3** ($0.426, same seats). Every objection carried a reviewer-authored `proposed_fix` and none disputed the direction:
  - glm: the missing-path refusal assumes no `omitempty`;
  - kimi: two amendments to the attended parent have no record;
  - sol: the card's GET and POST have no network deadline.
- **r2 by controller narrow-refinement carve-out:** the fixes applied in the reviewers' own text.
  - glm's branch (a) holds, measured: `omitempty` count 0 (V28).
  - V29 (`env.after`) and V30 (`DecidedMintRequests("")`) measured.
  - kimi's *Amendment posture*, V31 (the parent's R1–R5 fix no mint location or A3 vehicle) and V32 (base-drift trigger), plus the §10 mirror and the §15 checklist.
  - sol's 10-second AbortController paragraph and XA5.
- **Independent audit** (evaluator sonnet via the Agent tool; generator ≠ judge): **CARVE-OUT VALID, 84/100, zero blocking.**
  - It re-ran V28–V30, read V31 and D-WORLD-88 R1 (the deferral sequences A3's native add-ons and does not contradict R1), and judged kimi's objection a missing record, not a direction dispute.
  - Its consistency findings were applied (`fad08b71`): PG3b/XA5 added to §9's gates and the `-run` pattern (XA5 would otherwise not have run), PG3b renamed, V29/V30 cited.
  - Noted, not changed: glm's third golden case is a test obligation rather than a measured row, and the D-WORLD-114 sentence is a controller addition.
- **Designer diet:** one doc, r0 plus one protocol-mandated revision. Held; r2 cost no designer run.

**Not done this fire, on purpose: the A3a sprint-1 planner.** Sprint 1 builds the CLI minter, which is amendment (b), the thing D-WORLD-114 lets Mark veto, and the remaining slot after a design fire is not enough for a planner, executor and judge round (382's planner alone took 218 min). So the default (A) waits one report cycle before any A3a code.

**Routing evidence:** base=27ce163eb1ae8af55151ea6548591c15a3126835@2026-11-11T16:06:23Z.
- Controller: `claude-opus-5-5` (tok: not reported).
- Designer: `claude:claude-opus-5-5` via Agent `opus` (r0 197,650 tok; r1 resume 251,486 tok as reported); r2 by controller carve-out.
- Quorum: `ailang design-quorum --author claude:claude-opus-5-5`, rounds 1+2 $0.811 metered (artifacts `mission-quorum-world/w-world-approvals-by-key-a3-2026-11-11T16-31-03Z.json`, `…T16-47-14Z.json`).
- Planner and executor: not spawned (design-only fire; see above).
- Evaluator: `sonnet` via Agent (44,668 tok), auditing the r2 carve-out.
- The driver noted codex unusable at preflight (probe rc 75).
- Metered total: $0.81 (quorum only).

**Progress:** 1.0 clauses 1–7 MET; the goal is unmoved. A2 is code-complete on `dev`: an operator can enrol a passkey from the printed link and approve a pending request with it, with no World terminal prompt. Mark's real-passkey run (A2 §10, QUICKSTART §16) is now runnable. Row 201 also has a reviewed design for its last human-surface milestone: an operator approves a mint from the surface's *Needs you* card with a passkey or SSH signature.

**Ruled out:**
- **Merging #552 over the red (attempts 2–7):** an item is not LANDED until remote CI passes; it merged only after attempt 8 went 3/3.
- **Fix-forward of the apt step:** ci-health forbids it for a provider-side red.
- **Parking A3 as `needs-human-review` after round 2:**
  - every remaining objection had a concrete fix and none disputed the direction, so the carve-out applied;
  - an independent judge confirmed it;
  - the human question that remains (the two amendments) is asked as D-WORLD-114 with a default, not as a park.
- **Codex for the designer:** its lane was unusable at preflight (rc 75), and the resolver named opus.

**Retro:**
- **A design premise about browser behaviour gets probed at design time, not cited.** r0 measured the POST headers but asserted `crypto.subtle` byte-equality, and round 1 caught it. The designer measured it in five minutes once asked. Next design directive: list every browser-capability premise and require a V-row per premise.
- **Amending an attended parent needs a record in r0, not a disclosure.** r0 labelled the departure *DEPARTED-FROM-PARENT*. Round 1 let it pass; round 2 blocked on it. Next directive: any departure from an attended parent gets an *Amendment posture* paragraph and a ledger ask in r0.
- **The outage lasted about 3 hours** (first seen before 14:53Z, still failing at 16:54Z, clear at 17:00Z). Re-running and doing other work while waiting cost one design fire and no lost work; the same handling applies next time.

**Next:**
1. If this record's own PR is still open, land it first (docs only: the A3 design, record 383).
2. Then the A3a sprint-1 planner from `w-world-approvals-by-key-a3.md` r2 (no re-design). D-WORLD-114 default A applies from iteration 384 if unanswered. Run V32's base-drift check first.
3. Mark's A2 exit run (§10, QUICKSTART §16) is runnable and stays attended.

## 384 — 2026-11-11 — row 201 A3a sprint 1 LANDED (PR #554, merge `b44b5f7d`): credentials schema 7 `surface_approvals`, one shared judge `host/authority/judge.go`, and the CLI minter `approve <id> --await`; judge PASS 91; D-WORLD-114 default A applied [PRODUCT]

**Pick:** the A3a sprint-1 planner, the charter's Next and record 383's Next 2. Record 383 landed (#553 → `0a25c0d6`), so there was no carry-over PR. D-WORLD-114 was unanswered; its default A ("keep both amendments") applies from this iteration, because the ask was posted once in report 383 (18:14Z).
- **Re-measured at pick** (`origin/dev` `0a25c0d6`, record 383): no `a3a` plan, sprint branch, PR or worktree. The A3 doc has its quorum artifacts from iteration 383 (rounds 1 and 2, then the carve-out), so no quorum was owed at pick.
- **V32 base-drift trigger** (the design's own rule before sprint 1), at `0a25c0d6`:
  - V1: all 12 sha256 prefixes match, and `cmd/` and `host/` are byte-identical to the design's base `0140bd7f`;
  - V8: 5 `OpenCredentialsReadOnly` sites and no writer open in `host/daemon`;
  - V18: `--- PASS` ×4;
  - V15: `TestPasskeyPageInChrome` PASS.
  - No drift, so no r3.
- **Ledger:** valid. OPEN: D-WORLD-105 (default C) and D-WORLD-114 (default A, now applied).
- **Clause map:** 1.0 clauses 1–7 MET, none UNMET. Release rows 150/151 not routable (unchanged).
- **Orphans:** none. Iteration 383's fire ended normally.

**Gate 0:**
- Stamped first. `gh` is `sunholo-voight-kampff`. Billing CLEAN. Kill switch absent.
- Directives: 0 allowlisted on #524, of 17 comments. The watermark file is absent, so all comments were read.
- Self-notices: the script printed `watermark channel unavailable` (no `-last-seen` file for 524 or 449). This is an instrument failure, not a verdict, as in 371–383.
- Inbox: nothing unread addressed to `mission-world`.
- Skill drift: none. SKILL.md and all 13 resource files at the resolved symlink target match the driver clone's `origin/dev`.
- Weekly issue sweep: not owed.

**Planner** (Agent `opus`, 69 min, 452,011 tok): plan `w-world-approvals-by-key-a3a-sprint1-plan.md` with 7 compiled patches in Appendix P, 1,543 lines (579 production, 964 test) against the design's ≈ 1,650. The fallback cut did not fire. It ran 104 diff-anchored mutants (100 killed by both the named test and the suite, 4 named equivalents) and did its own round trip.
- **Premises found false (plan §1.4):**
  - P1: six older-store test setups must drop the new table, so "every `TestApprove*` unedited" holds only for the assertions;
  - P2: three schema pins edited;
  - P3: the context-root pin is `approve_await.go|pollSurfaceApproval`;
  - P4: the judge needs two more steps (`request`, `store`) to carry WF4/F12's texts and refuse a malformed row without a panic;
  - P5: a schema-6 store also refuses profile-bound mints until `serve` restarts, as every earlier bump did;
  - P6: JG1's text agreement is pinned by the CLI tests, because the wrappers live in `package main`.
- The planner added PK-NOW (`TestApproveByPasskeyJudgesTheApproverAtNow`) after its own first mutant run found `MUT-WP-JUDGE-NOW` surviving.

**Controller round trip:** `git archive 0a25c0d6`, the 7 patches extracted from the plan file with §2's awk and applied → `diff -rq` against the planner's `work/` tree is empty for `cmd` and `host`.

**Executor** (env `recipe claude:claude-haiku-5-5`, spawned as Agent `haiku`, as 374–382): two runs. Run A was M1a–M2b (7.9 min) and Run B was M3a–M3b (5.8 min). There were no hand edits, no git writes, no stop condition and no timing exception. The controller rebuilt each sub-milestone from its snapshot with `commit_ms.sh`: M1a `38a6d903` → 70, M1b `2912eeea` → 404, M1c `4fb902b4` → 463, M2a `e34a9e61` → 839, M2b `6f75243d` → 971, M3a `c6d4fb6c` → 1,313, M3b `9bfd0e2b` → 1,543. That equals the plan at every step, and the final tree equals the controller's round trip.
- **Controller error, recovered:** the first `commit_ms.sh` run deleted the two uncommitted plan files. The helper restores every untracked path that is not in the snapshot, and earlier iterations had committed the plan first. The planner's assembly parts survived in its scratch directory, so `assemble.py` and `mkjson.py` regenerated both files: 2,687 lines, with all 7 extracted patches `cmp`-identical to the copies taken before the deletion. They were committed as `b506ce93` before M3.

**Gates** (controller, before the judge, not concurrent with it):
- Plan §6 script: rc 0, `PRE_COMMIT_GATES_GREEN 408s`, with 33 named inventory, schema, help, quickstart and new tests `--- PASS` and the AW leg 8/8.
- `VERIFY_GO_LEGS=plain verify_go.sh`: rc 1 at load 15–19. The binary, toolchain and 64-evidence-test legs are ✓. Two reds are in packages the sprint does not touch, both on §6's load-sensitive list: `host/broker` `TestExecProtectedCheckTimeout/reserve` and `host/verifygate` `TestMissionConfigDefaultFlowIntegration/invalid_ledger`. Re-run alone, both are rc 0.

**Judge** (evaluator `sonnet` via Agent, 31 min, 77,541 tok; own worktree; generator ≠ judge): **PASS 91/100, zero blocking** (correctness 28, tests 24, security 19, standards 10, scope 10).
- 47 distinct mutants of its own (58 ledger lines). The must-die control was killed and the must-live controls survived. Every S-guard, the prune, the judge lookups, the minter's flag/poll/consume/AW7 branches were killed. Per-milestone non-vacuity held: M1b guards → SA2, M1c prune → SA3, judge lookup → JG1/JG2, mint-from-row → AW7.
- It judged P1–P6 sound.
- **Survivor:** `NOW`/`NOW-ssh` — the minter judging at `p.row.CreatedAt` instead of its poll `now` (`approve_await.go:91,94`) leaves every await test green. The controller reproduced it: `sed 's/Artifact, now)/Artifact, p.row.CreatedAt)/'` changed 2 lines, and `go test ./cmd/ailang-worldd/ -run 'Await|JudgesTheApproverAtNow'` returned rc 0. The production code is right; this is a test-only gap, carried to sprint 2 as an await analogue of PK-NOW. `ROW` (dropping `|| !p.recorded`) is equivalent.

**Gate 3b:**
- PR #554 head `9bfd0e2b`: 3/3 success on the first attempt (run 43111055834).
- `gh pr merge 554 --merge --match-head-commit 9bfd0e2b…` → merge **`b44b5f7d`** at 22:37:23Z. No closing keywords.
- Push CI on the merge `b44b5f7d` (run 43113317174): **3/3 success** (22:49Z). LANDED.

**Routing evidence:** base=0a25c0d6c6c9af6cdf81da75270f41af41a44bad@2026-11-11T20:06:31Z.
- Controller: `claude-opus-5-5` (tok: not reported).
- Designer: not spawned (A3 r2 covers A3a sprint 1).
- Planner: `opus` via Agent (452,011 tok), `agent-tool opus fail-closed:env-pin`.
- Executor: `claude:claude-haiku-5-5` via Agent `haiku` (66,873 + 39,719 tok).
- Evaluator: `sonnet` via Agent (77,541 tok), `agent-tool sonnet declared:alias-pin`.
- The driver noted codex unusable at preflight (probe rc 75).
- Metered total: $0.00.

**Progress:** 1.0 clauses 1–7 MET; the goal is unmoved. Row 201 A3a has its record and its minter on `dev`: credentials schema 7 stores a signed surface approval, one judge serves the daemon and the CLI, and `approve <id> --await` mints from a recorded approval with no terminal, re-judging the artifact rather than trusting the row. Nothing writes the table in production until sprint 2's POST route, so `--await` waits until expiry today.

**Ruled out:**
- **A design r3:** V32 found no drift.
- **The design's fallback cut** (M3 to sprint 2): the plan measured 1,543 lines against the 2,000 cap.
- **Treating the two `verify_go` reds as regressions:** untouched packages, both on the known load-sensitive list, both green alone, and CI 3/3.
- **Adding the NOW test in this fire:** the controller would be writing unreviewed code after the judge. It is a sprint-2 test obligation instead.

**Retro:**
- **`commit_ms.sh` deletes anything untracked that is not in the snapshot, the plan included.** Commit the plan files before the first sub-milestone commit. The banked copy (`verification/world-iter384/commit_ms.sh`) now refuses (exit 2) instead of deleting an untracked `design_docs/` path.
- **"Ship patches before mutants, checkpoint the plan" worked:** 69 min against 218 at 382 for a plan of similar size. Keep it.
- **Merging the plan's three runs into two executor runs** cost nothing (13.7 min in total).

**Next:**
1. Row 201 A3a sprint 2, "the routes" (M4 the read routes, the origin set and the guards; M5 the POST), from `w-world-approvals-by-key-a3.md` r2 §9 (no re-design). Carry the judge's NOW gap: an await test that enrols or revokes between the request and the poll.
2. Mark's A2 exit run (§10, QUICKSTART §16) stays attended and runnable. D-WORLD-114 stays OPEN for Mark to confirm or veto; default A is being built.

## 385 — 2026-11-12 — row 201 A3a sprint 2 LANDED (PR #556, merge `d1027477`): the read routes, the surface origin set and its guards, and `POST /surface/approve/{id}` (judge, then record; SSHSIG arm cut to sprint 3); judge PASS 92 [PRODUCT]

**Pick:** the A3a sprint-2 planner, which was record 384's Next 1 and the charter's Next. Record 384 had landed (#555 → `a8c66655`), so there was no carry-over PR. D-WORLD-114 is still unanswered, and its default A (build as designed) still applies.
- **Re-measured at pick** (`origin/dev` `a8c66655`, record 384): there is no `a3a` sprint-2 plan, sprint branch, open PR or worktree.
  - Sprint 1's symbols are on HEAD: `host/store/surface_approvals.go`, `JudgeAssertion`/`JudgeSSHSignature` in `host/authority/judge.go:65,102`, and `approve_await.go`.
  - Production code has 0 hits for `requests.json`/`"/surface/approve` and 0 for `SurfaceOrigins|surface-origin`. The control, `/surface/passkey`, has 11.
- **Quorum reader:** rc 0. The A3 doc has its two iteration-383 artifacts (both `blocked`, seats glm-5-3/kimi-k3/gpt6-1-sol, `absent_reviewers` `[]`), followed by the controller carve-out and the audit (VALID 84). No quorum was owed at pick.
- **Ledger:** valid, 101 rows. OPEN: D-WORLD-105 (default C) and D-WORLD-114 (default A, applied).
- **Clause map:** 1.0 clauses 1–7 MET, none UNMET. Release rows 150/151 are not routable (unchanged).
- **Orphans:** none. The only driver notice on #524 since 383 is iteration 381's `rc=143 after landing its record` (2026-11-11T10:51Z), already credited in log 382. All iteration-384 worktrees are clean.

**Gate 0:**
- Stamped first. `gh` is `sunholo-voight-kampff`. Billing CLEAN. Kill switch absent.
- Directives: 0 allowlisted on #524 (18 comments) and 0 on #449 (47). Neither watermark file exists, so every comment was read.
- Self-notices: rc 2, `watermark channel unavailable`. This is an instrument failure, not a verdict, as in 371–384.
- Inbox: nothing unread is addressed to `mission-world`. The newest messages are this mission's own 383/384 reports.
- Skill drift: none. `SKILL.md` and all 12 `resources/*.md` match the driver clone's `origin/dev`, at both the resolved symlink target and the pin worktree (inode 271495741 vs 272502315, two separate files). The mismatch control fired.
- Weekly issue sweep: not owed (#524 was opened Monday).

**Planner** (Agent `opus`, 77 min, 466,125 tok): plan `w-world-approvals-by-key-a3a-sprint2-plan.md` with 8 compiled patches in Appendix P.
- **Size:** 1,995 lines (689 production, 1,306 test) against the 2,000 cap. **The design's fallback cut fired:** with the SSHSIG arm the planner measured about 2,100 lines. R11's `sshsig` arm, RT6, RT13 and AW6's route half move to sprint 3 (Appendix C). An `sshsig` POST is a 422 naming `approve <id> --signature`, so nothing is silently accepted. The guards, the origin set `O` and `--surface-origin`/CF1 are not cut, as the design requires.
- **Mutants:** 115 diff-anchored, 114 killed by both the named test and the suite, 1 named equivalent (`MUT-R12`: `Binds` cannot fail after the law verified the same row). Run 1 exposed two test gaps, fixed before run 2:
  - `MUT-HEX32-LEN`: a 31-byte challenge case was added;
  - `MUT-R11-NOW`: the test rig now enrols the passkey after the requests are made.
- **Premises found false (plan §1.4):**
  - P1: `ensureCredentialsStore` opens its writer through `store.UpgradeCredentials`, not `OpenCredentials`, so RT12's inventory pins those two calls;
  - P2: the design's "five routes" includes sprint 3's card page and `approve.js`, so sprint 2 registers three;
  - P3: the CLI's admission refusal text lives in `package main`, so the model repeats it and a test pins the equality;
  - P4: the page's retries belong to sprint 3, and only the daemon side (503, `Retry-After: 1`, no kept writer) is here;
  - P5: "below schema 7 → 503" is caught by `RecordSurfaceApproval` after the writer opens, while a store below 6 is caught earlier by the judge's lookup;
  - P6: RT11's refusals are the handler's own (422/400), with `isProtected` asserted unchanged;
  - P7: the store refuses a 17th pending request, so only the decided cap of 20 is pinned.
- The carried NOW obligation is in M4a as `TestApproveAwaitJudgesTheApproverAtNow`. With iteration 384's `sed` applied (2 lines), both subtests fail.

**Controller round trip:** I ran `git archive a8c66655`, extracted the 8 patches from the plan file with §2's awk, and applied them. `diff -rq` against the planner's `work/` tree is empty for `cmd` and `host`, and the appended-byte control was detected. The plan was committed (`0d88d820`) before any sub-milestone commit (retro 384).

**Executor** (env `recipe claude:claude-haiku-5-5`, spawned as Agent `haiku`, as 374–384), two runs:
- Run A, M4a–M4e: 5.9 min, 52,970 tok.
- Run B, M5a–M5c: 7.3 min, 45,017 tok.
- No hand edits, no git writes, no stop condition.
- **One deviation:** Run B started M5a's whole-`host/daemon` leg under `nohup`, which §7 forbids. Its log reads `ok host/daemon 250.244s` and its rc file `0`, and the controller's §6 gate re-ran the package green. The deviation is procedural only.

The controller rebuilt each sub-milestone from its snapshot with `commit_ms.sh`. Cumulative numstat after each commit:

| Sub-milestone | Commit | Numstat |
|---|---|---|
| M4a | `fb4e2ae2` | 350 |
| M4b | `b9025b20` | 517 |
| M4c | `627d825e` | 793 |
| M4d | `f76c6bea` | 1,257 |
| M4e | `6a3f07ae` | 1,396 |
| M5a | `76aa7dcd` | 1,741 |
| M5b | `456fb26c` | 1,912 |
| M5c | `d52b3ae3` | 1,995 |

That equals the plan at every step. The final `cmd` and `host` equal the round-trip tree, and the M4e/M5c snapshots equal the committed files.

**Gates** (controller, before the judge, not concurrent with it):
- Plan §6 script: rc 0, `PRE_COMMIT_GATES_GREEN 355s`. Every inventory and all 20 new tests `--- PASS` by name, the three `-race` legs, and the AW leg 9/9.
- After commit, `VERIFY_GO_LEGS=plain verify_go.sh`: **rc 0** ("all 64 required top-level evidence tests passed exactly once"). There were no load-sensitive reds this time.

**Judge** (evaluator `sonnet` via Agent, 8.6 min, 84,797 tok; own worktree; generator ≠ judge): **PASS 92/100, zero blocking** (correctness 28, tests 27, security 19, standards 9, scope 9).
- Its ledger has 64 lines: 56 killed, 2 survived (J36 and the must-live control), 5 non-compiling edits and 1 no-diff edit (neither credited). It sampled about 25 of the plan's §4 rows and added more than 25 of its own, including fail-open statuses, wrong approver/artifact, and injected challenge and request-bytes leaks, which were caught.
- **Survivor J36, controller-reproduced:** in `handleSurfaceRequestCard`, replacing the `ActivePasskeys` error check with `if false` leaves the approval tests green (rc 0). The path is read-only and fail-soft (the card shows no passkeys instead of a 500), so it was not a blocker. It is carried to sprint 3 as a test obligation.
- **Must-die control, re-measured:** the judge's K01 (record before the judge) was killed by the writer *inventory* under `-failfast`, which does not show that a behavioural test catches it. The controller re-ran it with the inventory excluded: a record inserted before `JudgeAssertion` makes `TestApprovePOSTRefusalArms` fail on four R11 arms' 0-row checks (rc 1). Unmutated, the same command is rc 0. The planner's `MUT-RECORD-BEFORE-LAW` row agrees.
- **The judge's numstat was wrong:** it reported 1,596. The stated unit (insertions + deletions) reads 1,995, and insertions alone 1,991. Both are under the cap, so the verdict is unaffected.

**Gate 3b:**
- PR #556 head `d52b3ae3`, run 43138995090 attempt 1: `go host build + test gate` red on `host/broker` `TestSeedRestartIgnoresASurvivingOldCopy/survivor` (`exec_seed_test.go:656`). The sprint touches 0 `host/broker` files. This is **row 259, instance 4**.
- `gh run rerun --failed` → attempt 2: **3/3 success**.
- `gh pr merge 556 --merge --match-head-commit d52b3ae3…` → merge **`d1027477`** at 02:36:43Z. No closing keywords.
- Push CI on the merge (run 43142623795): **3/3 success**. LANDED.

**Routing evidence:** base=a8c666554f284651de33d2f6a2c85a16f5b72c02@2026-11-12T00:06:27Z.
- Controller: `claude-opus-5-5` (tok: not reported).
- Designer: not spawned (A3 r2 covers A3a sprint 2).
- Planner: `opus` via Agent (466,125 tok), `agent-tool opus fail-closed:env-pin`.
- Executor: `claude:claude-haiku-5-5` via Agent `haiku` (52,970 + 45,017 tok).
- Evaluator: `sonnet` via Agent (84,797 tok), `agent-tool sonnet declared:alias-pin`.
- The driver noted codex unusable at preflight (probe rc 75).
- Metered total: $0.00.

**Progress:** 1.0 clauses 1–7 MET; the goal is unmoved. Row 201 A3a now has its daemon half on `dev`. The surface can list pending requests, serve a request's card data, and accept a passkey approval that the shared law verifies before one row is recorded. `approve <id> --await` mints from that row with no terminal, and the end-to-end test proves the seam in process. No page posts yet; that is sprint 3.

**Ruled out:**
- **A design r3:** sprint 2 builds §2.3/§2.4 as written. Every design-vs-HEAD drift the planner found is recorded as P1–P7 and routed minimally.
- **Shipping the SSHSIG POST arm in this sprint:** about 2,100 lines against the 2,000 cap. The design names this exact cut.
- **Treating the PR's `host/broker` red as a regression:** the package is untouched, the failure is row 259's known CI-only flake (now 4 instances), and a re-run went 3/3.
- **Crediting the judge's K01 as behavioural proof of judge-before-record:** it was killed by an inventory. The controller re-measured it against the behavioural tests alone (above).
- **Adding the J36 test in this fire:** the controller would be writing unreviewed code after the judge. It becomes a sprint-3 obligation.

**Retro:**
- **Commit the plan before `commit_ms.sh` worked:** the plan was `0d88d820` before M4a, and nothing was deleted (the helper's refusal guard never fired).
- **A judge's control can be killed by the wrong test.** Under `-failfast`, an inventory that pins call sites kills any mutant that adds a call, so a "record before judge" control proves only the inventory. That was one friction this fire. Directive tweak for the next evaluator: run must-die controls with the inventory excluded, and name which test killed them.
- **The executor broke a `nohup` rule it was given verbatim**, after 374–384 obeyed it. The result was verifiable from its own artifact. Keep the rule in the directive; no further action until it recurs.
- **Planner time:** 77 min for 1,995 lines and 115 mutants, against 69 min for 1,543 lines at 384. The "ship patches first, checkpoint the plan" practice still holds.

**Next:**
1. Row 201 A3a sprint 3, "the page and the app", from `w-world-approvals-by-key-a3.md` r2 §9 (no re-design):
   - M6 the card page and `approve.js` (PG2/PG3/PG3b, network deadlines);
   - M7 *Needs you*, Decided, and the opt-in notification;
   - M8 live-Chrome XA1–XA5 and QUICKSTART §17 (attended — pending);
   - plus Appendix C's SSHSIG arm (R11 `sshsig`, RT6, RT13, AW6's route half) and the J36 test (an injected `ActivePasskeys` error on the card route → 500).
   - The design's own estimate is about 2,200 lines before Appendix C, so expect its fallback cut (M7's notification and Decided move to a fourth sprint).
2. Mark's A2 exit run (§10, QUICKSTART §16) stays attended and runnable. D-WORLD-114 stays OPEN for Mark to confirm or veto; default A is being built.

## 386 — 2026-11-12 — row 201 A3a sprint 3 LANDED (PR #558, merge `70823879`): the card page and `approve.js`, *Needs you*, the live-Chrome XA1/XA2/XA3-list/XA5 tests, QUICKSTART §17, and the carried SSHSIG POST arm and J36 test; notification and Decided cut to sprint 4; judge PASS 90 [PRODUCT]

**Pick:** the A3a sprint-3 planner. This was record 385's Next 1 and the charter's Next. Record 385 had landed (#557 → `82fb9afa`), so there was no carry-over PR. D-WORLD-114 is still unanswered, and its default A (build as designed) still applies.
- **Re-measured at pick** (`origin/dev` `82fb9afa`):
  - There is no sprint-3 plan, branch, open PR or worktree.
  - Production Go has 0 hits for `approve.js` and 0 for a `/surface/approve/` page route. The sprint-2 SSHSIG refusal is at `host/daemon/surface_approve.go:255`. Controls: `/surface/passkey` has 11 hits and `requests.json` has 5.
  - All three iteration-385 worktrees are clean.
- **Quorum reader:** rc 0, the iteration-383 artifact. No quorum was owed at pick.
- **Ledger:** valid, 101 rows. OPEN: D-WORLD-105 (default C) and D-WORLD-114 (default A, applied).
- **Clause map:** 1.0 clauses 1–7 MET, none UNMET. Release rows 150/151 are not routable (unchanged).
- **Orphans:** none. #524 has no `FAILED to complete` notice (0 of 20 comments).

**Gate 0:**
- Stamped first. `gh` is `sunholo-voight-kampff`. Billing CLEAN. Kill switch absent. Dev CI 3/3 at `82fb9afa` (attempt 2).
- Directives: 0 allowlisted on #524 (20 comments) and 0 on #449 (47). Neither watermark file exists, so every comment was read.
- Self-notices: rc 2, `watermark channel unavailable`. This is an instrument failure, not a verdict, as in 371–385.
- Inbox: 0 unread for `mission-world`; the control read 110 already-read messages.
- Skill drift: none. `SKILL.md` and all 12 `resources/*.md` at the resolved symlink target match the driver clone's `origin/dev`.
- Weekly issue sweep: not owed (#524 was opened Monday).

**Carried from 385:** the record merge `82fb9afa` was red once on the race gate. `host/verifygate` hit its 8-minute `-race` bound in `TestSeToolsGateArmsAreLoadBearing/plan-bytes`, which is row 249's flake, folded into row 258. Attempt 2 (run 43146351741) went 3/3. It was posted as a follow-up on #524 and is recorded here as row 258's next instance.

**Planner** (Agent `opus`, 97 min, 496,816 tok): plan `w-world-approvals-by-key-a3a-sprint3-plan.md` with 9 compiled patches in Appendix P.
- **Size:** 1,843 lines (463 production, 1,380 test) by `git diff --numstat 82fb9afa HEAD -- cmd host`, plus 18 QUICKSTART lines.
- **The design's first fallback step fired:** with M7's notification, the Decided list and the `N` key, the planner measured about 2,143 lines. XA4, XA3's Decided half and the `N` key move to sprint 4 (plan Appendix C). The carried SSHSIG arm fits and stays. Nothing on the design's never-cut list (PG3, PG3b, XA2 arm c, XA5) was cut.
- **Mutants:** 88, Go and JS: 87 killed by both the named test and the package suite, and 1 named equivalent (`JS-PAGEHIDE`: the browser already aborts a page's requests on navigation). No kill depends on an inventory test. Run 1 exposed `JS-WALK-TYPE`; XA2's type arm was changed to use `reuse`, and run 2 killed it.
- **J36 proof:** at the base, `if false` in place of the `ActivePasskeys` error check leaves 29 tests green. With M6a's `TestRequestCardPasskeyReadError`, it fails (a 200 card with `"passkeys":[]` instead of a 500). Dropping the passkey table did not reach the read, because the store's open refuses first. The test renames a column instead.
- **Live Chrome** (157.0.8103.51, headless): every XA test, the two planner-added page tests, and the A2 tests `--- PASS`.
- **Premises found false (plan §1.4):**
  1. XA2 arm (b) is a 409, not a 422, because the digest check (R9) runs before the signature (R11).
  2. The card data has no "recorded" field, so the page shows the POST's own answer.
  3. `grantNotifications`/`notify` move with XA4, and `click` becomes `press`.
  4. The design's Chrome `-run` pattern matches the fixture recorder, which skips by design.
  5. The J36 table-drop injection does not reach the passkey read (above).
  6. A store error in the SSH judge is a 500, not a 503.
  7. The `N` key lives in `surface.js`, not `app.js`.
  8. The `daemon.go` line numbers have moved.
- Part of the planning (M6a, M5d, the J36 proof) was done by a forked agent, which the planner reviewed and stacked unchanged.

**Controller round trip:** I ran `git archive 82fb9afa`, extracted the 9 patches from the plan file with §2's awk, and applied them. `diff -rq` against the planner's `work/` tree is empty for `cmd`, `host` and `docs`, and the appended-byte control was detected. The plan was committed (`888d415c`) before any sub-milestone commit.

**Executor** (env `recipe claude:claude-haiku-5-5`, spawned as Agent `haiku`, as 374–385), two runs:
- Run A, M6a–M7a: 10.4 min, 50,422 tok.
- Run B, M8a–M8d (incl. the live-Chrome legs): 4.8 min, 69,352 tok.
- No hand edits, no git writes, no `nohup`, no stop condition, and no timing exception.

The controller rebuilt each sub-milestone from its snapshot with `commit_ms.sh`. Cumulative numstat after each commit:

| Sub-milestone | Commit | Numstat |
|---|---|---|
| M6a | `854c0a82` | 29 |
| M5d | `3a023498` | 280 |
| M6b | `fea1889a` | 788 |
| M6c | `693f01c7` | 943 |
| M7a | `34aaa64b` | 976 |
| M8a | `9424444d` | 1,383 |
| M8b | `4d161376` | 1,502 |
| M8c | `ea34fd73` | 1,788 |
| M8d | `2936964e` | 1,843 |

That equals the plan at every step (docs 18). The final committed tree equals the round-trip tree.

**Gates** (controller, before the judge, not concurrent with it):
- Plan §6 script: rc 0, `PRE_COMMIT_GATES_GREEN 470s`. It covered every inventory and new test `--- PASS` by name, the three `-race` legs, the AW leg 9/9, and the live-Chrome leg (8 top-level `--- PASS`, 0 FAIL/SKIP).
- After commit, `VERIFY_GO_LEGS=plain verify_go.sh`: **rc 0** ("all 64 required top-level evidence tests passed exactly once") on the first run.

**Judge** (evaluator `sonnet` via Agent, 24.1 min, 114,178 tok; own worktree; generator ≠ judge): **PASS 90/100, zero blocking** (correctness 28, tests 25, security 19, standards 10, scope 8).
- Its ledger has 71 lines: 5 non-compiling edits and 1 run that executed no tests (neither credited). Every kill names its test, and none was by an inventory alone. The `innerHTML` mutant is killed only by PG2's static banned-token test (C03b), which survives when that test is excluded (C03c).
- Its review confirmed the following:
  - The SSHSIG arm judges on the read-only handle at `now` before any writer opens, and every refusal returns before the write.
  - Both page routes go through `approvalRoute`.
  - `approve.js` hashes the stored bytes itself and refuses unlisted, mistyped or missing fields. It writes only through `textContent`, retries only on 503 `Busy`, and shows an unanswered POST as "outcome unknown".
  - The 409 deviation and the sprint-4 cut are sound.
- **Survivors (non-blocking):**
  - The Needs-you list's `esc()` on `agent`, the `id` href and `recorded` (A02/A03/A05) is unasserted.
  - No Chrome test renders a hostile field value (C03/C04).
  - The expiry display, label text and unknown-outcome styling are unasserted.
- **A02 reproduced by the controller:**
  - Removing `esc()` on `q.agent` leaves `TestSurfaceAppPure` green (rc 0) with `WORLD_EXEC_NODE` set.
  - I first checked that the judge's command (which omitted that variable) still runs the test: it does, so the survivor is real.
  - The code escapes correctly, so this is a test gap only. It carries to sprint 4.

**Gate 3b:**
- PR #558 head `2936964e`, run 43170254467 attempt 1: **3/3 success**.
- `gh pr merge 558 --merge --match-head-commit 2936964e…` → merge **`70823879`** at 07:00:43Z. No closing keywords.
- Push CI on the merge (run 43172124328) attempt 1: the race gate was red on `host/daemon` `TestMCPPostBudgetProductionConstants/batch`, at `mcp_budget_test.go:208` ("timed callback did not return within finite teardown").
  - That is a 1 s wall-clock bound under `-race`, and the package took 430 s.
  - The merge tree is byte-identical to the PR head tree (`bcc9100e`), whose race job passed. The sprint touches no MCP-budget code.
  - This is the test's first recorded instance (0 hits in the charter, log, archive and index), so it is filed as a row-258-class CI flake.
  - `gh run rerun --failed` → attempt 2: **3/3 success**. LANDED.

**Routing evidence:** base=82fb9afa9086eed2f781bb0d4fe54d4bae3b70d7@2026-11-12T04:07:19Z.
- Controller: `claude-opus-5-5` (tok: not reported).
- Designer: not spawned (A3 r2 covers A3a sprint 3).
- Planner: `opus` via Agent (496,816 tok), `agent-tool opus fail-closed:env-pin`.
- Executor: `claude:claude-haiku-5-5` via Agent `haiku` (50,422 + 69,352 tok).
- Evaluator: `sonnet` via Agent (114,178 tok), `agent-tool sonnet declared:alias-pin`.
- The driver noted codex unusable at preflight (probe rc 75).
- Metered total: $0.00.

**Progress:** 1.0 clauses 1–7 MET; the goal is unmoved. Row 201 A3a now has its page on `dev`:
- A pending request is listed under *Needs you* in the surface app and opens a card page.
- The card renders every signed field and verifies the digest itself, then lets the approver sign with a passkey or an SSH key and posts it.
- The daemon judges either signature before recording it, and `approve <id> --await` mints.
- A swapped, extended or truncated request is refused before any `credentials.get`.
- What remains for A3a is sprint 4 (the opt-in notification, the Decided list, the `N` key) and Mark's attended exit run.

**Ruled out:**
- **A design r3:** sprint 3 builds §2.7/§2.8/§2.10 as written. Every design-vs-HEAD drift is recorded in plan §1.4 and routed minimally.
- **Cutting the SSHSIG arm a second time:** it fits under the cap (1,843) once M7's notification and Decided list move, which is the design's first fallback step.
- **Treating the merge race red as a regression:** the identical tree passed the same job at the PR head, the test is a 1 s wall-clock bound in untouched code, and a re-run was used to confirm.
- **Adding the A02 test in this fire:** the controller would be writing unreviewed code after the judge. It becomes a sprint-4 obligation.

**Retro:**
- **The "inventory kills are not behavioural" directive worked.** The judge named the killing test for every mutant and re-ran the `innerHTML` kill without the static test (C03c survived). That separated a static guard from a behavioural one with no controller re-measurement. Keep it in the evaluator directive.
- **A judge's survivor command can drop an environment variable the plan sets.** A02's command omitted `WORLD_EXEC_NODE`. It made no difference here, because the test finds node on `PATH`, but the controller checked before banking the survivor. For the next directive: require mutant commands to copy the plan's acceptance environment verbatim.
- **The executor obeyed `nohup`/background rules** this time, after 385's single deviation, with the line repeated in the directive's rules.
- **Planner time:** 97 min for 1,843 lines and 88 mutants with live Chrome, against 77 min at 385 (no browser). Within the 100-minute aim.

**Next:**
1. Row 201 A3a sprint 4, from this plan's Appendix C (no re-design):
   - the opt-in notification (XA4, design §2.9);
   - the Decided list (XA3's other half) and the `N` key;
   - the judge's A02/A03/A05 test gap (assert the Needs-you list escapes `agent`, the `id` href and `recorded`, e.g. by a hostile-value case in `TestSurfaceAppPure`);
   - optionally, a Chrome case that renders a hostile field value on the card.
2. Mark's A2 exit run (§10, QUICKSTART §16) and the A3a exit run (QUICKSTART §17, after sprint 4) stay attended. D-WORLD-114 stays OPEN for Mark to confirm or veto; default A is being built.

## 387 — 2026-11-12 — row 201 A3a sprint 4 LANDED (PR #560, merge `4134ce89`): the Decided list, the opt-in browser notification (XA4), the `N` key, the Needs-you `esc()` tests and a hostile-value card render; A3a code-complete except the attended exit run; judge PASS 90 [PRODUCT]

**Pick:** the A3a sprint-4 planner. This was record 386's Next 1 and the charter's Next loop row. Record 386 had landed (#559 → `bed9d135`), so there was no carry-over PR. D-WORLD-114 is still unanswered, and its default A (build as designed) still applies. The notification built here is the web one of design §2.9 (A3a), not the native one that D-WORLD-114 (a) defers.
- **Re-measured at pick** (`origin/dev` `bed9d135`):
  - There is no sprint-4 plan, branch, open PR or worktree. `gh pr list --state open` was empty.
  - `app.js` has 0 `Notification` hits, and `grantNotifications` has no match in `cmd`/`host`. Control in the same call: `getRecorder|fetchRewrite` matched 4 files.
  - `needsYouHTML` renders `pending` only. `requests.json`'s `decided` exists (`approve_model.go:68`). The `N` key exists only in `surface.js:156`.
- **Quorum reader:** rc 0, the iteration-383 artifact. No quorum was owed at pick.
- **Ledger:** valid, 101 rows. OPEN: D-WORLD-105 (default C) and D-WORLD-114 (default A, applied).
- **Clause map:** 1.0 clauses 1–7 MET, none UNMET. Row 150 waits on attended A1, and row 151 sits after the critical path (D-WORLD-64); neither is routable.
- **Orphans:** none. #524 has no `FAILED to complete` notice.

**Gate 0:**
- Stamped first. `gh` is `sunholo-voight-kampff`. Billing CLEAN. Kill switch absent. Dev CI 3/3 at `bed9d135` (all checks green).
- Directives: 0 allowlisted on #524 (22 comments) and 0 on #449 (47). Neither watermark file exists, so every comment was read.
- Self-notices: rc 2, `watermark channel unavailable`. This is an instrument failure, not a verdict, as in 371–386.
- Inbox: 0 unread for `mission-world`.
- Skill drift: none. `SKILL.md` and all 13 `resources/*.md` at the resolved symlink target match the driver clone's `origin/dev`.
- Weekly issue sweep: not owed (#524 was opened Monday).

**Carried from 386:** the record merge `bed9d135` (docs only, 0 Go files) was red once on the race gate. `host/surfacefollow` `TestBackoff` hit a 1 s test-synchronization wait under `-race` (`scheduler_test.go:131`). It is the test's first recorded instance, in the row-258 class. Attempt 2 (run 43177817156) went 3/3. It was posted as a follow-up on #524 and is recorded here.

**Planner** (Agent `opus`, 84.6 min, 395,119 tok): plan `w-world-approvals-by-key-a3a-sprint4-plan.md` with 6 compiled patches in Appendix P.
- **Size:** 389 lines (354 + 35) by `git diff --numstat bed9d135 HEAD -- cmd host`, plus 4 QUICKSTART lines. Nothing was cut, including the optional hostile-value Chrome test.
- **Production:** `host/surfacepage/app/app.js` only. No Go production file changes, no route is added, and `surfaceCSP` is byte-unchanged.
- **Mutants:** 60: 59 killed by both the named test and the package suite, and 1 named equivalent (`NP-NOAPI`: a guard for browsers without the Notifications API, unreachable in Chrome and node). No kill depends on an inventory or static test.
- **A02/A03/A05 proof:** each of `needsYouHTML`'s 8 `esc()` calls was removed one at a time. At the base, 6 survive `TestSurfaceAppPure` (agent, the link `id=`, its `href`, recorded, the refused id, the unavailable text). On M9a and on the final tree, all 8 are killed.
- **Live Chrome:** every test and subtest `--- PASS` on the final tree, including the new XA4 (`TestNeedsYouNotifiesOnlyWhenAskedOnce`), XA3's Decided half and `TestApprovePageShowsAHostileValueAsText`.
- **Premises found false (plan §1.4):**
  1. The design's `N` key is the static pages' (`surface.js`), whose key handler never runs on the app shell. `app.js` gains its own.
  2. The timeline's Enter swallowed Enter on now-card links.
  3. The app polls every 15 s only in follow mode. In v1 mode a new request appears, and notifies, only on the next log commit or reload. "No new timer" was kept as designed and noted in Appendix C.
  4. The test driver does not await Promises, and a poll during a refresh is dropped, so the tests wait for each refresh to finish.
  5. The design's line numbers have drifted.
  6. The notification recorder never calls the real constructor.
  7. `World.open` answers before the brane has opened.
- **Planner decisions** (§1.3): the page's first read raises no notification (the design says "on a refresh"). Times are UTC `hh:mm`, as elsewhere in the app.
- **Gate flake flagged:** its whole-gate `verify_go.sh` failed once on `host/broker` `TestPublishApprovalRefusalSetWithALandedPositiveControl` (`POST=2`, want 1), an untouched package, green alone at the base and on the final tree. This is a new instance, added to the plan's known list. The controller's `verify_go` run did not reproduce it (below).

**Controller round trip:** I ran `git archive bed9d135`, extracted the 6 patches from the plan file with §2's awk, and applied them. `diff -rq` against the planner's `work/` tree is empty for `cmd`, `host` and `docs`, and the appended-byte control was detected. The plan was committed (`b4f0a2fb`) before any sub-milestone commit.

**Executor** (env `recipe claude:claude-haiku-5-5`, spawned as Agent `haiku`, as 374–386), two runs:
- Run A, M9a–M9c: 1.6 min, 45,509 tok.
- Run B, M10a–M10c: 3.3 min, 45,004 tok.
- No hand edits, no git writes, no `nohup`, no stop condition, and no timing exception.

The controller rebuilt each sub-milestone from its snapshot with `commit_ms.sh`. Cumulative numstat after each commit:

| Sub-milestone | Commit | Numstat |
|---|---|---|
| M9a | `fce449c4` | 6 |
| M9b | `0b5ee70d` | 36 |
| M9c | `c7e473de` | 93 |
| M10a | `b4df4261` | 239 |
| M10b | `5d5fdeb8` | 357 |
| M10c | `6723807a` | 389 |

That equals the plan at every step (docs 4). The final committed tree equals the round-trip tree.

**Gates** (controller, before the judge, not concurrent with it):
- Plan §6 script: rc 0, `PRE_COMMIT_GATES_GREEN 577s`. All seven packages `ok`, every inventory and new test `--- PASS` by name, the three `-race` legs, the AW leg, and the live-Chrome leg including XA4 and the hostile-value test.
- After commit, `VERIFY_GO_LEGS=plain verify_go.sh`: **rc 0** ("all 64 required top-level evidence tests passed exactly once") on the first run.

**Judge** (evaluator `sonnet` via Agent, 8.4 min, 71,557 tok; own worktree; generator ≠ judge): **PASS 90/100, zero blocking** (correctness 27, tests 27, security 19, standards 9, scope 8).
- Its ledger has 36 lines, every kill naming its test: `TestSurfaceAppPure`, `TestNeedsYouNotifiesOnlyWhenAskedOnce`, `TestNeedsYouListsPendingRequests` or `TestApprovePageShowsAHostileValueAsText`. None rests on an inventory or static test. The must-die control was killed, and three must-live controls survived.
- Its review confirmed the following:
  - `requestPermission()` is reached only from the Notify-me click.
  - A notification needs `granted`, never fires on the first read, fires once per id, and skips refused and recorded requests.
  - The Decided list escapes every field.
  - `N`/Enter do not fire inside inputs.
  - Deviations (i)–(v) are sound.
- **Survivor (non-blocking):** `DEC-GUARD-ARR`, `Array.isArray(r.decided)` → truthiness in `decidedHTML`. The daemon always sends an array. **Controller-reproduced:** `TestSurfaceAppPure` rc 0 with the edit and `WORLD_EXEC_NODE` set (`ctl-DEC-GUARD-ARR-survivor.log`).
- **Judge claim corrected by the controller:** it calls the notification's unescaped `href` id safe "by store schema". The table's CHECK is only `GLOB 'mr-[0-9a-f]*'`, so it constrains the fourth character alone. The real guard is `PutMintRequest`'s `MintRequestIDPattern` (`^mr-[0-9a-f]{16}$`, `host/store/mint_requests.go:105`), plus the fixed same-origin `/surface/approve/` prefix. The conclusion holds; its stated reason was the weaker one.
- **Judge shortfall:** the directive asked for at least 12 mutants of its own beyond the plan table. Only 11 of its ids are absent from the table, and some are renamed plan rows (`CARD-INNERHTML` is the plan's `CARD-TEXTCONTENT`). Its plan-row sample (24+) and its controls are real. Noted for the next directive, and not scored against the sprint.

**Gate 3b:**
- PR #560 head `6723807a`, run 43200210085 attempt 1: **3/3 success**. Mergeable CLEAN. Closing-keyword scan: 0 hits over the PR body, title and the 7 commit messages (control fired).
- `gh pr merge 560 --merge --match-head-commit 6723807a…` → merge **`4134ce89`** at 10:41:58Z.
- Push CI on the merge (run 43202662787) attempt 1: **3/3 success**. LANDED.

**Routing evidence:** base=4134ce898548ec7b9c16c89f952544cdb5ea5d7a@2026-11-12T10:58:15Z.
- Controller: `claude-opus-5-5` (tok: not reported).
- Designer: not spawned (A3 r2 covers A3a sprint 4).
- Planner: `opus` via Agent (395,119 tok), `agent-tool opus fail-closed:env-pin`; derive-planner-lane `opus fail-closed:env-pin`.
- Executor: `claude:claude-haiku-5-5` via Agent `haiku` (45,509 + 45,004 tok), resolver `recipe claude:claude-haiku-5-5 declared:provider-pin`.
- Evaluator: `sonnet` via Agent (71,557 tok), `agent-tool sonnet declared:alias-pin`.
- The driver noted codex unusable at preflight (probe rc 75, quota unreadable).
- Judge independence: cross-model (opus planner, haiku executor, sonnet judge). Metered total: $0.00.

**Progress:** 1.0 clauses 1–7 MET; the goal is unmoved. Row 201 A3a is now code-complete on `dev`:
- An operator sees *Needs you* and *Decided* in the surface app, can opt in to a browser notification for each new request, and can jump to it with `N`.
- The card renders every signed field as text and verifies the digest itself. The operator signs with a passkey or an SSH key, the daemon judges before it records, and `approve <id> --await` mints.
- What remains for A3a is Mark's attended exit run (QUICKSTART §17). A3b (native) waits for row 188 N3. Next is A4.

**Ruled out:**
- **A design r3:** sprint 4 builds §2.8/§2.9 as written. Every design-vs-HEAD drift is recorded in plan §1.4 and routed minimally.
- **Treating the planner's `host/broker` POST-count red as a regression:** `host/broker` is untouched, the test is green alone at the base and on the final tree, and the controller's `verify_go` run passed it.
- **Adding the `DEC-GUARD-ARR` test or `href` hardening in this fire:** the controller would be writing unreviewed code after the judge. Both are optional A4-era polish (Appendix C), not A3a obligations.

**Retro:**
- **A small sprint is cheap end to end:** planner 85 min, executor 4.9 min, gates 10 min, judge 8 min, CI about 15 min per run. The size rule ("prefer a small, fully proven sprint") produced 389 lines with nothing cut.
- **The judge's "own mutants" count needs an instrument.** A count claimed in prose was not checked against the plan table. Next directive: require the judge to list its own-mutant ids and the command that shows they are absent from plan §4 (a `comm` against the table's ids).
- **Judge reasons can be weaker than its conclusions.** "Safe by the store schema" pointed at the weak CHECK, not the Go writer guard. Reproduce the reason, not only the verdict (gate-2 rule (a), aimed at a judge's justification).
- **Protocol shortcut, stated:** the controller did not re-read `resources/verification-protocol.md` (136 KB) in full this fire; it worked from gate-2's recall list and the rules quoted in the gates it read in full (0, 1, 2, 3 route, 4). Nothing in this sprint needed a rule beyond those. Say so rather than imply a full read.
- **The planner flagged a new flake instead of chasing it.** `TestPublishApprovalRefusalSetWithALandedPositiveControl` is recorded as a first instance, in the row-258 class. Record it as instance 2 if it recurs.

**Next:**
1. Row 201 **A4**: `Human.Approve` decision packets (row 172) approved the same way (parent `w-world-approvals-by-key.md` §A4). It needs a NEW child design by the rotation designer: the last used was `claude:claude-opus-5-5`, so next is `codex:gpt-6.1-sol` (re-probe it; rc 75 this fire), falling to `pi:ollama/glm-5.3:cloud`. Then the quorum with `--author`, then the plan.
2. Mark's A2 exit run (§10, QUICKSTART §16) and the A3a exit run (QUICKSTART §17) are both runnable and attended. D-WORLD-114 stays OPEN for Mark to confirm or veto; default A is built.

## 388 — 2026-11-12 — row 201 A4 child design r2 (`w-world-approvals-by-key-a4.md`): an operator approves or denies a `Human.Approve` decision packet by typed `y`, SSH key or passkey; quorum blocked twice, controller carve-out, independent audit VALID 84; D-WORLD-115 filed on its one amendment to the parent [PRODUCT]

**Pick:** row 201 A4, the charter's Next loop row and record 387's Next 1. Record 387 had landed (#561 → `b0fbcd3e`), so there was no carry-over PR.
- **Re-measured at pick** (`origin/dev` `b0fbcd3e`): there was no `a4` doc, branch, PR or worktree. `mc_quorum_gate2 w-world-approvals-by-key-a4` returned rc 1; the A3 control returned rc 0.
- **The gap, measured by the controller before the designer ran:** the exported operator entry point `broker.DecideApproval` has **no** production caller (control: 5 `_test.go` files call it). The only production decisions are `RecordOperatorApproval` (`host/daemon/episode_finish.go:173`, which creates and decides its own request) and `MintAttendedApproval` (`cmd/world-publish/main.go:310`, publish). A decision packet produced in a daemon episode therefore stays `approvals chain: pending`, and nothing can decide it.
- **Ledger:** valid, 101 rows before and 102 after. OPEN: D-WORLD-105 (default C), D-WORLD-114 (default A, applied), and D-WORLD-115, filed this iteration.
- **Clause map:** 1.0 clauses 1–7 MET, none UNMET. Release rows 150/151 are not routable (unchanged).
- **Orphans:** none.

**Gate 0:**
- Stamped first. `gh` is `sunholo-voight-kampff`. Billing CLEAN. Kill switch absent. Dev CI 3/3 at `b0fbcd3e`.
- Directives: 0 allowlisted on #524 (24 comments). The watermark file is absent, so every comment was read.
- Self-notices: rc 2, `watermark channel unavailable`, an instrument failure and not a verdict (as in 371–387).
- Inbox: 0 unread for `mission-world`.
- Skill drift: none. `SKILL.md` and all 13 `resources/*.md` at the resolved symlink target match the driver clone's `origin/dev`.
- Weekly issue sweep: not owed.

**Carried from 387:** the record PR #561 (docs only; its code tree equals merge `4134ce89`) was red once on its head's race gate. `scripts/marketingcheck` `TestMarketingSupervisorFakeDaemonCleanup/shell-wait` hit its 30 s shell-trap cleanup bound (row 254's class, folded into row 258). Run 43205372503 attempt 2 went 3/3. It was posted as a follow-up on #524 and is recorded here.

**Design** (`w-world-approvals-by-key-a4.md`, a child of the attended parent). Resolver: `recipe claude:claude-opus-5-5 declared:provider-pin`; spawned as Agent `opus`, following the 372–383 precedent. The rotation file reads `claude:claude-opus-5-5`. Record 387's Next named `codex:gpt-6.1-sol` as the rotation's next entry, but codex was over ration (`MISSION_OVER_RATION=codex openrouter`, preflight probe rc 75), and the resolver's answer was followed. Worktree `.wt-world-iter388-a4`, branch `design/world-iter388-row201-a4`.
- **r0** (16.6 min, 853 lines, 31 V-rows):
  - **The shape:** the daemon decides, because it is the only writer of the world store. The design adds one route, `POST /surface/decide/{ref}`, with three arms: `typed` (the operator-proof file of `episode finish`), `sshsig` and `webauthn`, behind A3a's guards. It also adds a subject route `GET /surface/decisions/{ref}` and the CLI form `ailang-worldd approve <sha256:request-ref> [--deny] [key flags]`.
  - **`decidedBy` is derived by the daemon** (`operator:terminal`, `ssh:<name>`, `webauthn:<name>`), never supplied by the caller.
  - **A pure refusal ladder** D1–D10 lives in `host/broker/decision_law.go`. The broker now refuses a second decision of the same request.
  - **The signed message** is `world-decision/v1 ‖ verb ‖ packet digest ‖ request digest ‖ expiry`, 91 bytes, under the SSHSIG namespace `world-decision@sunholo`. The verb is signed, and the domain is separate from the mint message.
  - **R4:** `validatePublishApproval` never reads `DecidedBy` (V11). A4 therefore fences publishing twice:
    - L1 refuses publish-shaped and finish-shaped scopes for every format;
    - L2 has the publish consumer refuse `ssh:`/`webauthn:` decisions.
  - **The controller re-measured the R4 premise first-party:** `DecidedBy` occurs 0 times in the body of `validatePublishApproval` (control: 2 occurrences in the file).
  - **Designer-found facts (V9, a throwaway probe):** identical requests share a ref; a request can be decided twice and the newest decision wins; `Human.PollApproval` ignores scope.
- **Quorum round 1 BLOCKED 2/2** ($0.167). The author's vendor (Anthropic) was benched. gemini-4-argon and gpt6-1-sol were present; oc-glm-5-3 and oc-kimi-k3 were absent as `unreachable` (ollama 429, weekly usage limit).
  - **sol:** D4's "exactly one packet" rested on a 101-packet scan ordered by hash, with no branch for completeness. The controller confirmed `packetScanLimit = 100` (`decisionpacket_validate.go:34`).
  - **gemini:** a crash between `RecordKeyDecision` and `DecideApproval` strands the passkey counter, and two refs lacked a `sha256:*` CHECK.
  - Both were real, each carried a fix, and neither disputed the direction.
  - A first launch was wrapped in `timeout`, which is absent on macOS. It exited rc 127, ran no reviewer and wrote no artifact; it was re-run without the wrapper.
- **r1** (designer resumed by SendMessage, 4.5 min, 931 lines, 34 V-rows):
  - a complete request-wide packet discovery (`Discovery{Complete,…}`; D4a refuses an incomplete scan);
  - both DDL CHECKs, verbatim;
  - K0, an idempotent retry on the exact artifact of an unlanded row;
  - V32: a counter-0 assertion passes the counter check every time, so a second use is refused by D6 and the primary key instead. The controller re-read `usePasskey`'s SQL (`passkeys.go:225`).
- **Quorum round 2 BLOCKED 2/2** ($0.194, the same seats, glm and kimi still 429). Both reviewers independently found ONE defect on ONE surface: r1's K0 sits after the read-only judge's W15 counter check, so a counter-N retry dies before K0 can see it. Each carried a concrete fix, and neither disputed the direction.
- **r2 by the controller's narrow-refinement carve-out.** §2.7's K0 paragraph is replaced by gpt6-1-sol's two-path judge, in its own words:
  - W15 is deferred to the writer transaction;
  - an exact match bypasses only the counter;
  - revocation is re-checked on both paths;
  - the mint judge is unchanged.
  The two fixes differ on one point: gemini's also skips `useApprover`. sol's stricter re-check was kept. RD14 is added with both reviewers' test text verbatim.
- **Independent audit** (evaluator `sonnet` via Agent, 2.0 min, 75,259 tok; generator ≠ judge): **CARVE-OUT VALID, 84/100, zero blocking.**
  - It re-ran V3–V6, V11, V32 and V33, all matching.
  - It traced five retry paths (counter-7 retry, a different counter-7 assertion, counter-0 replay, a revoked approver, SSHSIG), and all close.
  - It confirmed R4's L1 and L2 each stop a key-decided publish alone, with must-die and must-live rows.
  - **Its N1 was real:** an end-to-end revoked-approver arm cannot kill its mutant, because the read-only judge refuses a revoked approver first. It moved to a store-level KD9.
  - N2, N3, N5 and N6 were applied in the same consistency commit (`0722d533`). N4 (folding K0 into the Judge and handler-order bullets) is left to the planner.
- **Designer diet:** one doc, r0 plus one protocol-mandated revision. Held; r2 cost no designer run.
- **Size:** 3 sprints, about 5,250 changed lines (≈1,950 / 1,800 / 1,500), each with a priced fallback cut.

**D-WORLD-115 filed** (the design's §0 amendment posture and §12 ask, filed in r0 as iteration 383's retro asked). The parent §2.2 specifies the mint message with a stored 32-byte daemon challenge. A4 instead signs the decision message with the packet digest in the challenge's slot: freshness comes from an expiry of at most 10 minutes, and single use comes from the approvals chain. **A** (default): keep it. **B**: add the stored challenge (≈ +550 lines in sprint 2).

**Not done this fire, on purpose: the A4 sprint-1 planner.** D-WORLD-115 amends the attended parent and is new in this report. The design says sprint 1 does not depend on the answer, but sprint 1 builds the decision message the ask is about, so its default waits one report cycle (as D-WORLD-114 did at 383–384).

**Routing evidence:** base=b0fbcd3eee7b5b29190e678e296635a3dd3a0d0c@2026-11-12T12:59:39Z.
- Controller: `claude-opus-5-5` (tok: not reported).
- Designer: `claude:claude-opus-5-5` via Agent `opus` (r0 227,163 tok; r1 resume 260,698 tok as reported); r2 by controller carve-out.
- Quorum: `ailang design-quorum --author claude:claude-opus-5-5`, rounds 1+2 $0.361 metered (gemini-4-argon $0.250, gpt6-1-sol $0.111); oc-glm-5-3 and oc-kimi-k3 absent both rounds (ollama 429). Quoted as BLOCKED at N−2, not "quorum passed".
- Planner and executor: not spawned (design-only fire; see above). Resolver: planner `agent-tool opus fail-closed:env-pin`, executor `recipe claude:claude-haiku-5-5 declared:provider-pin`.
- Evaluator: `sonnet` via Agent (75,259 tok), `agent-tool sonnet declared:alias-pin`, auditing the r2 carve-out.
- The driver noted codex unusable at preflight (probe rc 75, quota unreadable).
- Metered total: $0.36 (quorum only).

**Progress:** 1.0 clauses 1–7 MET; the goal is unmoved. Row 201 now has a reviewed design for its last milestone: an operator approves or denies an agent's `Human.Approve` request (the decision packet row 172 produces) with a typed `y`, an enrolled SSH key or a passkey, from the CLI or the surface, and a key decision can never authorize a publish.

**Ruled out:**
- **Parking A4 as `needs-human-review` after round 2:** the only remaining objection was one defect, raised by both reviewers, each with a concrete fix and no dispute of the direction. The carve-out applied, and an independent judge confirmed it.
- **Gemini's variant of the round-2 fix** (also skip `useApprover` on the exact-match path): sol's transactional revocation re-check is strictly safer. KD9 pins it.
- **Codex for the designer:** it was over ration with probe rc 75, and the resolver named opus.
- **Treating glm/kimi's absence as a lane park:** the quorum degrades by design and records the absence by name. The doc's own designer ran, and the absent seats were reviewers, not a role.

**Retro:**
- **Two-round convergence on one surface is the carve-out's best case.** In round 2 both reviewers found the identical ordering defect independently, which made the objection unambiguous. A fix that adds a new path (K0) needs its reachability traced through every earlier step. Next designer directive: for any new branch, require a trace from the request's entry point to that branch for each input class.
- **Tooling:** macOS has no `timeout`. The first quorum launch exited rc 127, ran nothing and wrote no artifact, and only the rc file showed it. Bound background commands with the harness's own timeout instead (rule 6), and read the rc before the artifact.
- **The ollama quorum seats were both 429 on their weekly limit.** The quorum ran at 2 of 3 seats both rounds. If it recurs next week, it is a capacity pattern to report, not an item fact (standing rule 8e).
- **Protocol shortcut, stated:** the controller did not re-read `resources/verification-protocol.md` (136 KB) in full this fire. It worked from gate 2's recall list and the gates it read (0, 1, 2, 3 route, 4, 5).

**Next:**
1. If this record's own PR is still open, land it first (docs only: the A4 design and record 388).
2. Then the **A4 sprint-1 planner** ("the law and the typed root": the decision message and law, the broker's no-redecide rule, the R4 fences L1/L2, the complete discovery, `GET /surface/decisions/{ref}`, the typed POST arm, `approve <request-ref> [--deny]` typed) from `w-world-approvals-by-key-a4.md` r2, with no re-design. D-WORLD-115's default A applies from iteration 389 if Mark has not answered. Fold the audit's N4 into the plan.
3. Mark's A2 exit run (QUICKSTART §16) and A3a exit run (QUICKSTART §17) are runnable and attended. D-WORLD-114 and D-WORLD-115 stay OPEN for Mark to confirm or veto.

## 389 — 2026-11-12 — row 201 A4 sprint 1 LANDED (PR #563, merge `338ff821`): the decision message, the pure law D1–D10, complete packet discovery, first-decision-final, the R4 fences L1/L2, and `GET /surface/decisions/{ref}` + `POST /surface/decide/{ref}` (typed arm); the typed CLI moved to sprint 2; judge PASS 88 [PRODUCT]

**Pick:** the A4 sprint-1 planner. This was record 388's Next 2 and the charter's Next loop row. Record 388 had landed (#562 → `9a6b959d`), so there was no carry-over PR. D-WORLD-115 is unanswered, and its default A (keep the decision message with a derived challenge) applies from this iteration, as filed.
- **Re-measured at pick** (`origin/dev` `9a6b959d`):
  - There was no A4 sprint-1 plan, branch or worktree, and `gh pr list --author @me --state open` was empty.
  - `grep -rl "DecideApproval(" cmd host --include="*.go"` → `host/broker/approve.go` (the definition) and 5 `_test.go` files: there is still no production caller (design V3).
  - `packetScanLimit = 100` (`decisionpacket_validate.go:34`) and `validatePublishApproval` (`approve.go:517`) are as the design measured.
- **Quorum reader:** rc 0, two iteration-388 artifacts. No quorum was owed at pick.
- **Ledger:** valid, 102 rows before and after. OPEN: D-WORLD-105 (default C), D-WORLD-114 (default A, applied), D-WORLD-115 (default A, applied from 389).
- **Clause map:** 1.0 clauses 1–7 MET, none UNMET. Release rows 150/151 are not routable (unchanged).
- **Orphans:** none.

**Gate 0:**
- Stamped first. `gh` is `sunholo-voight-kampff`. Billing CLEAN. Kill switch absent. Dev CI 3/3 at `9a6b959d`.
- Directives: 0 allowlisted on #524 (25 comments). The watermark file is absent, so every comment was read.
- Self-notices: rc 2, `watermark channel unavailable`. This is an instrument failure, not a verdict, as in 371–388.
- Inbox: 0 unread for `mission-world`.
- Skill drift: none. `SKILL.md` and all 13 `resources/*.md` at the resolved symlink target match the driver clone's `origin/dev`.
- Weekly issue sweep: not owed (#524 was opened Monday).

**Planner** (Agent `opus`, 105 min, 435,275 tok): plan `w-world-approvals-by-key-a4-sprint1-plan.md` with 8 compiled patches in Appendix P.
- **Size:** 1,985 lines (587 production, 1,398 test) by `git diff --numstat 9a6b959d HEAD -- cmd host`, 15 under the cap.
- **Cut:** the typed CLI (CD1–CD4, ≈ 400 lines at the design's price) moved to sprint 2's head by design §9's priced fallback. With it the sprint would have been ≈ 2,385. So after this sprint nothing shipped writes the operator proof the POST needs: the route is guarded and has no production caller yet.
- **Mutants:** 96 rows. 88 were killed by both the named test and the package suite, 3 are named equivalents and 5 are named survivors with reasons. The main survivor class is the `humanMu` lock, which tests that post one request at a time cannot observe. No kill rests on an inventory or static test. Ledger: 104 lines.
- **Premises found false (plan §1.4):**
  - **P2:** the design says first-decision-final leaves the existing deciders unaffected. The unedited `TestEpisodeFinishOperatorProofExpiryBoundary` approves the same finish scope twice in the same second and got a 500. The plan's fix is that a byte-identical re-decision returns the existing ref and writes nothing, while any different second decision is refused.
  - **P3:** `context_thread_test.go`'s store-read census moves from 5/4 to 6/5 (the new decision walk). It is the only edited existing test.
  - **P8:** the broker's no-redecide check is unreachable through the POST, because the law's D6 refuses first. It has unit coverage only.
  - The rest: `TestSurfaceAppRoutes` needs no edit (P1); `Discovery` lives in `host/broker` because the coordinator imports the broker (P4); the subject's `passkeys` field and R4d's must-die half move to sprint 2 (P6, P7); an unrelated `DecisionRefusal` already exists in `host/surface` (P9).
- **N4 / K0 handler order (decision S1-D9):** Host → session/cookie refusal → Origin → `Sec-Fetch-Site` → media type → body → lock → read facts → law D1–D10 → **arm slot** → `DecideApproval` → 200. Sprint 2's key arm fills the arm slot in this order: read-only judge (W15 deferred) → writer → `RecordKeyDecision` with K0 → `DecideApproval` → `LandKeyDecision`. No sprint-1 line moves.

**Controller round trip:** I ran `git archive 9a6b959d`, extracted the 8 patches from the committed plan with §2's awk, and applied them. `diff -rq` against the planner's `work/` tree is empty for `cmd` and `host`; the control (the same diff against the unpatched base) shows 17 differences. The plan was committed (`9ac369fa`) before any sub-milestone commit.

**Executor** (env `recipe claude:claude-haiku-5-5`, spawned as Agent `haiku`, as 374–388), two runs:
- **Run A (M1a–M2), 3.6 min + 18.3 min resumed, 50,487 + 67,802 tok.** Its first turn hit the Bash tool's default 120 s timeout on M1a's compile fence, moved the acceptance shell into the background, and **ended its turn announcing a wait** (standing rule 7's sub-agent case). The controller checked the artifacts: no log had grown, and the background shell was lost. It then resumed the executor by `SendMessage` with explicit instructions: a 600 s tool timeout, foreground only, no Monitors. The resumed run re-ran every M1a leg in the foreground and finished cleanly.
- **Run B (M3a–M3d), 8.4 min, 51,374 tok.** Its directive carried the timeout rule from the start, and it ran clean. The whole `host/daemon` leg was `ok` in 343 s.
- No hand edits, no git writes, no stop condition, and no timing exception.

The controller rebuilt each sub-milestone from its snapshot with `commit_ms.sh`. Cumulative numstat after each commit:

| Sub-milestone | Commit | Numstat |
|---|---|---|
| M1a | `678f95e5` | 162 |
| M1b | `1eb10754` | 508 |
| M1c | `c929b9b8` | 733 |
| M2 | `5a3cd887` | 1004 |
| M3a | `bc0e925c` | 1356 |
| M3b | `74496cc8` | 1628 |
| M3c | `b7d512f4` | 1901 |
| M3d | `4c2cac17` | 1985 |

That equals the plan at every step, and the final committed tree equals the round-trip tree.
- **Controller instrument slip:** my Run-A commit loop used `declare -A` for per-milestone test commands. macOS bash 3.2 has no associative arrays, so all four commits ran M2's test command as their test gate. Each commit was still gofmt/build/vet-gated, and the executor had run each sub-milestone's exact acceptance. Run B's loop passed each command directly.

**Gates** (controller; the plan §6 script with every leg recording its own rc so one red does not stop the rest; 2,828 s under rig load 40–210):
- `verify_ail.sh`, gofmt, build, vet and the compile fence: rc 0. `-race` on the touched packages: no `FAIL` and no `DATA RACE`. All **68** named inventory and new tests print `--- PASS`. approval, approvaltest, coordinator and surfacepage are `ok`.
- **`host/daemon`:** 5 reds. `TestMarketingCaptureInventoryMutation`, `TestSessionSelfGit`, `TestServeStartsWhenTheProbeSpendsTheStartupClock` and `TestRegistryPhaseHasItsOwnClock` are on the known list. `TestWorkspaceProtectedPolicyWiring` is new to the list. All except `TestServeStarts…` (red at the base, per the planner) passed alone (`ok` 11.7 s).
- **`host/broker`:** `TestProtectedChecker/index-flag` was green alone. `TestExecProtectedCheckTimeout/reserve-after-before` ("returned at/after the original deadline") was **red once alone at the head**, then **green at the base**, then **green twice alone at the head**. It is the known `/reserve` timing family, and the sprint touches no exec code. Recorded as a load flake with all four runs cited.
- **After commit, `VERIFY_GO_LEGS=plain verify_go.sh`: rc 1 under load.** The evidence-manifest gate passed (64 of 64). `go test ./...` hit **five 10-minute package-timeout panics** (`cmd/ailang-worldd`, `host/broker`, `host/coordinator`, `host/daemon`, `host/verifygate`), and reds in untouched `host/capsule`, `cmd/world-publish`, `host/pkgproj`, `host/store` and `scripts/marketingcheck`. The two reds not on any known list (`host/store` `TestConvertedStoreMethodsM5b`, `host/pkgproj` ×2) were `ok` alone. Clean-runner CI was taken as the gate, as in 381–384.

**Judge** (evaluator `sonnet` via Agent, 49.8 min, 76,719 tok; own worktree; generator ≠ judge): **PASS 88/100, zero blocking** (correctness 27, tests 24, security 19, standards 10, scope 8).
- It confirmed the 91-byte layout, the law order including D4a/D4b, discovery's bound, sentinel and cursor, L1 and L2, and the POST's guard order, `humanMu` span, derived `decidedBy` and verb-bound proof.
- **P2 is sound:** the byte-identical comparison covers `requestRef`, verb, `decidedBy` and `now`, so a different decision can never pass.
- **The CLI cut is honest:** the route is unused but guarded, with no reachable attack surface.
- **Mutants:** its ledger has 46 lines. 36 compile-checked mutants (+2 COMPILEERR) ran: 29 killed, each by a named behavioural test, and 7 survived. `comm -12` against the plan's ids found 0 shared. The must-die controls (law skipped, proof consumption skipped) were killed, and the must-live control J06 survived.
- **Non-blocking test gaps** (carried to sprint 2):
  - **J07:** the L1 effect-substring arm (`Contains` → `HasSuffix` survives).
  - **J31/J32:** the daemon tests do not pin D1's regex anchor or hex case. Production is correct, but a regression would let `hashref.MustParse` panic on a 65-hex ref.
  - **J33:** the GET's D4a mapping is untested.
  - J06, J11 and J38 are equivalent mutants.
- **Controller-reproduced:** J31 (`{64}$` → `{64}` in `surface_decide.go:29`) gives `go test ./host/daemon/ -run '^(TestDecide|TestDecision)'` rc 0. Restored, worktree clean.

**Gate 3b:**
- PR #563 head `4c2cac17`, run 43296861559 attempt 1: **3/3 success**. Mergeable CLEAN. Closing-keyword scan: 0 hits over the PR body, title and the 9 commit messages (control fired).
- `gh pr merge 563 --merge --match-head-commit 4c2cac17…` → merge **`338ff821`** at 20:20:41Z.
- Push CI on the merge (run 43299366793) attempt 1: **3/3 success**. LANDED.

**Routing evidence:** base=338ff821c49ad10c7b6707df1d6398ec0c5cbe9a@2026-11-12T20:21:03Z.
- Controller: `claude-opus-5-5` (tok: not reported).
- Designer: not spawned (A4 r2 covers sprint 1).
- Planner: `opus` via Agent (435,275 tok), resolver `agent-tool opus fail-closed:env-pin`; derive-planner-lane `opus fail-closed:env-pin`.
- Executor: `claude:claude-haiku-5-5` via Agent `haiku` (50,487 + 67,802 resumed + 51,374 tok), resolver `recipe claude:claude-haiku-5-5 declared:provider-pin`.
- Evaluator: `sonnet` via Agent (76,719 tok), `agent-tool sonnet declared:alias-pin`.
- The driver noted codex unusable at preflight (probe rc 75, quota unreadable).
- Judge independence: cross-model (opus planner, haiku executor, sonnet judge). Metered total: $0.00.

**Progress:** 1.0 clauses 1–7 MET; the goal is unmoved. Row 201 A4 now has its law on `dev`:
- A `Human.Approve` decision packet has a pure, reviewed refusal ladder and a complete discovery.
- A request can be decided at most once.
- A key-made decision can never authorize a publish.
- The daemon serves the decision subject and accepts a typed decision behind the A3a guards.

An operator cannot yet decide a packet from a terminal: that is sprint 2's first milestone.

**Ruled out:**
- **Treating the run-A executor's one-line "I am waiting" return as a result:** the artifacts showed it had not progressed, so it was resumed rather than banked (standing rule 7, sub-agent corollary (a)).
- **Treating the head's `TestExecProtectedCheckTimeout` red as a regression:** it went green twice alone at the head, the base control passed, and no exec code is touched.
- **Blocking the landing on the loaded rig's `verify_go` rc 1:** every red was a timeout panic, on a known list, or green alone, and clean-runner CI went 3/3 on head and merge.
- **Re-design:** none was needed. The planner's premise deviations (P2, P3, P8) were judged sound by the evaluator.

**Retro:**
- **Sub-agent timeout → background → ended turn:** this is the second time an executor ended its turn on a wait (iteration 385's was a `nohup`). The cause here is mechanical: the Bash tool's default 120 s timeout, not a choice to background. Run B's directive carried "every Go Bash call passes `timeout: 600000`, foreground", and it ran clean. Keep that line in every executor and judge directive (now in the iteration-389 templates).
- **Controller tooling:** macOS `/bin/bash` is 3.2, with no `declare -A`. Use a function per milestone, or pass commands positionally, as Run B's loop did.
- **Gates under load:** the plan §6 script's `set -e` stops at the first red package. A per-leg rc variant (`gates_run.sh`, every leg recording its own rc) ran everything in one pass. Keep it.
- **Protocol shortcut, stated:** the controller did not re-read `resources/verification-protocol.md` (136 KB) in full this fire. It worked from gate 2's recall list and the gates it read (0, 1, 2, 3 route + role-spawn-routing, 3b, 4, 5).

**Next:**
1. If this record's own PR is still open, land it first (docs only).
2. Then the **A4 sprint-2 planner** ("the keys"), headed by the typed CLI CD1–CD4 moved from sprint 1, from `w-world-approvals-by-key-a4.md` r2 and the sprint-1 plan's Appendix C, with no re-design. Expect design §9's webauthn fallback cut (sprint 2 would be ≈ 2,200 otherwise). Carry the judge gaps J07, J31/J32 and J33, and the sprint-1 handler order S1-D9 for K0.
3. Mark's A2 exit run (QUICKSTART §16) and A3a exit run (QUICKSTART §17) are runnable and attended. D-WORLD-114 and D-WORLD-115 stay OPEN for Mark to confirm or veto; both defaults are applied.

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
