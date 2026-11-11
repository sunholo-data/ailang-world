MISSION-ROLE: planner

You are the sprint-planner for Ailang World mission iteration 380 (unattended run; no human is present).
Write the sprint plan for **row 201 `w-world-approvals-by-key`, milestone A2, SPRINT 2** — "the assertion law, the store and the binding".

## Where you work
- Sprint worktree (already created, branch `sprint/world-iter380-row201-a2s2`, base = origin/dev
  `3b77c294000ca40c2fe77e312ae0e8ce0f71a030`, which contains A2 sprint 1 merged as `d5d90e49`):
  `/Users/voightkampff/dev/sunholo-data/.wt-world-iter380-row201-a2s2`
- Scratch (compile your stages, run mutants) OUTSIDE the repo: `$HOME/.ailang/state/world-iter380/plan-scratch/`
  (never under /tmp — location-sensitive tests go red there).
- Write exactly two files in the worktree, and DO NOT git add/commit them (the controller commits):
  - `design_docs/planned/w-world-approvals-by-key-a2-sprint2-plan.md`
  - `design_docs/planned/sprint_w-world-approvals-by-key-a2-s2.json` (same shape as `sprint_w-world-approvals-by-key-a2-s1.json`)
- Leave every code file in the worktree untouched (the executor applies your patches there). No git writes of any kind in the worktree.

## Read first
1. `design_docs/planned/w-world-approvals-by-key-a2.md` (r2, ratified via quorum + carve-out; NOT reopened): §0, §2.9 (assertion law), §2.10 (store, schema 6), §2.11 (authority), §2.15, §3, §6 "Sprint 2", §7 (repo-wide inventories), §8, §9 (milestones, gates, **sprint 2 fallback cut**).
2. `design_docs/planned/w-world-approvals-by-key-a2-sprint1-plan.md` — the format to reproduce (it worked: haiku executor applied 4 patches with zero deviations). Especially **Appendix C** (M3b, carried to the head of THIS sprint), §1.4 premises (P2, P4 bind M3b), §2 patch protocol, §6 gates, §7 executor notes.
3. `CLAUDE.md`, `design_docs/coding-standards.md` (binding).

## Scope and order
- **M3b first** (Appendix C exactly: `passkey_assert.go` with `ApproveByAssertion`, `checkApprover` shared with `ApproveBySignature`, `signed.go`'s `credentialID`/`signCount`/`PasskeyUse()`, `approvaltest/passkey.go` in an external test package per P4, PA1–PA4, RC1 edit, TP1, XW1, XW2's assertion part).
- **Then M4** (schema 6, `host/store/passkeys.go`, PS1–PS5, RC6, edited version pins) and **M5** (consume C13–C18, CS7–CS9, CS2/CS3 arms, RC3, then authority A6 + AU3).
- **Size cap: ≤ 2,000 changed lines**, measured as `git diff --numstat <base> HEAD -- cmd host` (state this exact path filter in the plan; `.gitattributes`/docs are not counted). Sprint 1 estimated M3b ≈ 450 and sprint 2 ≈ 1,600, so **apply the design's §9 sprint-2 fallback early**: AU3 + authority's A6 move to the head of sprint 3 if needed. If still over 2,000 after that, cut at a milestone boundary that leaves every committed state safe (e.g. M5 consume to sprint 3), say which, and list the moved items in an Appendix C for sprint 3 exactly as sprint 1 did. Never exceed 2,000.
- Split into sub-milestones of ≤ ~250 changed lines each where possible (M3b may need M3b1/M3b2, etc.), each independently green.

## Carry-forwards (from iteration 379's judge — fold in where they belong, test-only where possible)
- `Passkey.Key` / `Passkey.CredentialID` aliasing not pinned (copies today) — add a pin that a caller mutating the input slices after `NewPasskeyApprover` does not change the approver (and the reverse for any getter).
- `host/approval/testdata/webauthn/SOURCE.txt` names `es256.json`/`eddsa.json` while the files are `chrome-*.json` — fix the text (XW0's SHA pin covers the fixtures; check whether SOURCE.txt is pinned too and update consistently).
- `RegistrationInput.Challenge` zero / `Now` unchecked — decide (design says caller's job): either a typed refusal + test, or a documented decision in §1.3. Measure first.

## Retro directives (iterations 374–379; each was paid for)
1. **Ship compiled patches**: every line of code and test as a unified diff in Appendix P (```` ````diff <M> ```` fences, as sprint 1), generated LAST from compiled, vetted, gofmt-clean, green scratch stages; then extract them FROM THE PLAN FILE with §2's awk, apply to a fresh export of the base, and `diff -r` against your mutant tree — byte-identical or regenerate. Round-trip must include tests added late (iteration 378's S3 dropped one).
2. **Expected numstat after each sub-milestone's commit** (measured after commit, `-- cmd host`), as sprint 1 §7.
3. **Mutants anchored to the DIFF**: one mutant per new production branch/refusal, enumerated from the diff, not from the test plan; run each against the named test row AND the suite; record killed/survived; name equivalents with reasons. Base controls where a test pins pre-existing behaviour.
4. **Flag-state coverage (379's judge finding)**: for every refusal that tests a combination of independent bits/flags (UP, UV, BE, BS, AT, ED…, counter zero/non-zero, revoked/active), give an arm with each other flag set, so `!a && !b`-style mutants die.
5. **Repo-wide inventories**: run every design §7 inventory by name with `-v` against your final scratch tree (AC10 exec-site census, dependency allowlists, context roots, `TestNoNewDeadlineFreeStoreReads`, `TestReadGettersHonorContext`, `TestCredentialsMethodsRequireADeadline`, reflective RC tests, `TestApprovalImportsOnlyTheStandardLibrary`, `TestApprovalTestIsTestOnly`, `TestFrozenV1RouteTableMatchesMux`, tracked-binary gate via verify_go.sh's flags). A new store method must take a deadline context and appear where those inventories demand.
6. The §6 gate script for the controller (one bash shell, `tee` every leg to `$HOME/.ailang/state/world-iter380/gates/`, never only through `tail`), including `-race` on `host/approval/...`, `host/store` (touched tests), `host/authority`, and the inventories, plus `VERIFY_GO_LEGS=plain ./scripts/verify_go.sh` after commit. Measure each gate green at the base first (rule: a gate red at base measures the repo).
7. Executor notes as sprint 1 §7: no git writes; patch -p1 is the only writer; runs of ~2 sub-milestones; cumulative snapshots to `$HOME/.ailang/state/world-iter380/snap/<M>/` with MANIFEST.sha256 (exclude the two plan files); stop conditions; numstat checks. Paths must name iter380 and the worktree above.
8. Acceptance criteria (§9-style) the judge can score against, numbered.

## Hard constraints
- D-WORLD-109: test keys, temp stores and in-memory/test authenticators only. No real key, no `~/.ssh`, no `~/.ailang/world`, no real `/dev/tty`, no mint or publish, no browser/Chrome/Node (sprint 3). Never run bare `ssh-keygen` (it once prompted for ~/.ssh/id_ed25519).
- `host/approval` stays pure and standard-library-only. Effects only at the store/CLI boundary.
- Do not touch: `tools/launchd/*`, any skill, `world/*.ail`, `packages/`, `scripts/` (except reading), `cmd/world-publish/*`, the typed-`y` path, `~/.ailang/state/mission-v1*`.
- Export `AILANG_BIN=$HOME/.pinned-ailang/ailang` in every shell that runs Go tests (`host/daemon` tests red by design without it).
- The controller shell is zsh: run loops/globs/`$VAR` work under `bash -c` or `bash <<'EOF'`.
- Known load-sensitive / flaky tests (re-run once alone, do not chase): `TestSeedRestartIgnoresASurvivingOldCopy/survivor` (host/broker, row 259 — 3 CI instances, 0 local reds in 60), local `-race` of host/broker, host/daemon, host/verifygate, `TestWorktreeSnapshotRefusesAConfiguredFilter`, `TestMissionConfigDefaultFlowIntegration/invalid_ledger`.
- Premises: any sentence in your plan that says "the code currently does X" needs a §1.2 row with the command and its observed output; an empty result needs a known-positive control in the same call. Where the design and HEAD disagree, measure, route minimally, and record it in §1.4.

## Waiting discipline (binding on you)
If you background anything (a long test run, a mutant loop), poll its ARTIFACT (an output/rc file it writes on completion) inside a bounded `date +%s` loop in the same turn. A Monitor is an event stream, not a wait; nothing will re-invoke you when it fires. Never end your turn announcing a wait, a pause, or "once it completes" — end only when both files are written and the round-trip has passed.

## Your final message
Report: the plan path, the cut you chose (what moved to sprint 3 and why, with measured numbers), sub-milestone list with expected numstat after each, mutant counts (total / killed / named equivalents), the round-trip result, the base-gate results, every premise you found false (§1.4), and anything the controller must decide.
