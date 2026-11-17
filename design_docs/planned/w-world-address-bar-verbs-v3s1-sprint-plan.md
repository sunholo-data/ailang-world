# Sprint plan — w-world-address-bar-verbs-v3, sprint S1 "the law" (iteration 400, row 243)

> **THE SPRINT — row 243 phase V3, sprint S1 (design §0, §1, §2.1–§2.3, §2.7 (the parts `CheckConfirmation` and
> `FinishPromptText` serve), §2.9, §3–§8, §11, §12, §14–§16):** the pure law of a key-confirmed `episode finish`, with
> **no route, no store read or write, no schema, no daemon change, no `.ail`, and the CLI's output byte-identical**.
> `host/approval` gains the verb request **`world/verb-request/v1`** (`VerbRequest`, `Validate`, `CanonicalVerb`,
> `DecodeVerb`) and its **86-byte message `world-verb/v1`** (`VerbMessage`, `VerbMessageHash`, `VerbPrefix`,
> `VerbNamespace = "world-verb@sunholo"`), the laws **`ConfirmBySignature`** (P1–P6) and **`ConfirmByAssertion`**
> (P1–P5, P3w, W1–W15 in full), **`SignedVerb`** (`Binds`, `ApprovedBy`, `Use`, `PasskeyUse`) and
> **`CheckConfirmation`** (the CLI's re-check: the law minus P1 and W15). `host/broker` gains **`FinishPromptText`** and
> **`DisplayCommitter`**, and the CLI's `finishPrompt`/`displayCommitter` delegate to them — **byte-identical**, pinned
> by a golden the **base tree's own `finishPrompt` generated** (VT1).
> - **Measured: 1,300 changed lines** — `git diff -M --numstat d3b0a813 HEAD -- cmd host docs` (insertions + deletions,
>   git's rename detection on; no rename occurs), after a commit of each of the four scratch stages: **422 production**
>   (`verb.go` 172, `verb_signed.go` 187, `finish_prompt.go` 38, `episode_finish_operator.go` 15, `committer.go` 10),
>   **878 test** (incl. the 67-line golden), **0 docs**. **Cap 1,400: margin 100.**
> - **The cut fired (S1-D2):** the design's S1 compiled to **1,626** lines (M1 1,034 + M2 326 + M3 266) — over the cap.
>   **M2 (the judges `JudgeVerbSSH`, `JudgeVerbAssertion`, `RecheckVerb`, VJ1) moves to a named sprint S1b, measured
>   326 lines** (compiled, green, its patch kept in scratch, Appendix C); M1's law is fully proven here (VL1–VL8), and
>   M3 stays because sprint 2's card (XC1) and sprint 3's F5 need `FinishPromptText` first.
> - **Four sub-milestones, each green at its commit** (§1.2 C11): **M1a** the request and the message (VL1–VL3) →
>   **M1b** the laws and `SignedVerb` (VL4–VL7) → **M1c** `CheckConfirmation` (VL8, the recorded-Chrome half of VL4) →
>   **M3** the moved prompt text (VT1).
> - **VT2 (`FinishConfirmRefusal` F1–F5) is sprint 3's**, not this sprint's (S1-D3: design §2.3 [r1] and §8's table are
>   the later text; §6's "Sprint 1" list is the stale r0 placement). **D-WORLD-118 = A** (attended, `72007b37`): the
>   B-only branches (F6, `DraftOnly`, "B +≈ 15") are not built in any sprint (S1-D4).
> - **Every line of code, test and golden is a unified diff in Appendix P**, generated from compiled, vetted,
>   gofmt-clean, green scratch stages; extracted **from this file** with §2's `awk`, applied with `patch -p1 -E` to a
>   fresh `git archive d3b0a813` export and compared byte for byte with the stage tree (§1.2 C12). The executor applies;
>   it writes nothing.

- **Design**: [w-world-address-bar-verbs-v3.md](w-world-address-bar-verbs-v3.md), **r2** (quorum r1/r2 BLOCKED;
  controller carve-out r2 applying the reviewers' own fixes; independent audit sonnet PASS 87, 0 blocking; iteration
  399; **not reopened**). Where this plan and the design differ, this plan was measured at `d3b0a813` and wins; each
  difference is a decision (§1.3) or a premise found false (§1.4).
- **Previous sprint of the row**: V2b (iteration 398, PR #587, merge `c251132f`):
  [w-world-address-bar-verbs-v2b-sprint-plan.md](w-world-address-bar-verbs-v2b-sprint-plan.md) (this plan's format).
  **The law analogue**: row 201 A4 sprint 1 ([w-world-approvals-by-key-a4-sprint1-plan.md](w-world-approvals-by-key-a4-sprint1-plan.md)),
  whose `decision.go`/`decision_signed.go`/`judge_decision.go` this sprint mirrors (§1.2 C4–C6).
- **Authority**: D-WORLD-107 = A; D-WORLD-109 = A (test keys only: `host/approval/approvaltest` and A2's recorded Chrome
  fixtures; no store, no browser, no tty); D-WORLD-89 R1–R5 untouched (no POST, no key path reachable from anything);
  **D-WORLD-118 = A** (Mark's attended ruling 2026-11-17, flip `72007b37`, controller-verified
  `scripts/mission_ledger_provenance.sh --id D-WORLD-118` → `verdict=attended`).
- **Planner base**: `d3b0a8139bdecc36f8d083f4b701aef822340b27` (= origin/dev, PR #590's merge). Worktree
  `/Users/voightkampff/dev/sunholo-data/.wt-world-iter400-row243-v3s1`, branch `sprint/world-iter400-row243-v3s1`.
  `git diff --stat 994ba123 d3b0a813 -- cmd host docs` is **empty** (§1.2 C2): the design's measurement base is the code.
- **Planner scratch** (outside the repo, never under `/tmp`): `$HOME/.ailang/state/world-iter400/plan-scratch/` (`base/`
  the pristine `git archive d3b0a813` export, a throwaway repo tagged `BASE`; `base-gen/` the base clone that generated
  VT1's golden; `work/` the stage repo — branch `s1` with tags `M1a M1b M1c M3`, branch `s1b` with tag `S1b-M2`;
  `final-tree/` its clone at `M3` for the final gates; `rt/` the round trip; `snaptest/` the snapshot rehearsal;
  `patches/` (incl. `S1b-M2.patch`); `tools/` (`gates.sh`, `accept.sh`, `accept_all.sh`, `mkpatches.sh`,
  `roundtrip.sh`, `snap_snippet.sh`, `mutants.py`, `mutate.py`, `assemble.py`); `logs/` (`base_gates/`,
  `final_gates/`, `accept_<M>.log`, `accept_all.txt`, `measure.txt`, `mutant_ledger.tsv`, `mut/`, `roundtrip.log`)).
- **Executor**: as routed by the controller. **One run**: M1a → M1b → M1c → M3. Each sub-milestone is one
  `patch -p1 -E` (§2) plus its acceptance (§3). **Apply, never retype.**
- **Touched** (§5): new `host/approval/{verb.go, verb_signed.go, verb_test.go, verb_signed_test.go, verb_check_test.go,
  verb_chrome_test.go}`, `host/broker/{finish_prompt.go, finish_prompt_test.go, testdata/finish-prompt.golden.json}`,
  `cmd/ailang-worldd/finish_prompt_golden_test.go`; edited `cmd/ailang-worldd/{committer.go,
  episode_finish_operator.go}`. **No pre-existing test file changes.** **Outside the design's §7 lists** (each in §1.4):
  `host/broker/testdata/finish-prompt.golden.json` and `cmd/ailang-worldd/finish_prompt_golden_test.go` (VT1's golden and
  its generator/pin, S1-D7); the test split into four `host/approval` files (S1-D12).
- **Not touched**: `host/authority/*` (the judges are S1b's), `host/store/*`, `host/daemon/*`, `host/surfacepage/*`,
  `host/verbs/*`, `host/broker/approve.go` and `decision_law.go` (A4's L1 and `reservedReason` unchanged, AM6),
  `cmd/ailang-worldd/episode_finish.go` (its call `finishPrompt(tty, scope, actualTitle, rendered, protected,
  pinned.Committer)` at `:410` is unchanged), every `episode_finish*_test.go`, `world/*.ail`, `packages/**`,
  `scripts/**`, `tools/launchd/*`, `cmd/world-publish/*`, any skill, `docs/`, `~/.ssh`, `~/.ailang/world`,
  `~/.ailang/state/mission-v1*`.

---

## 0. Binding constraints (each is restated where it governs)

1. **Test keys only (D-WORLD-109 = A).** Every SSHSIG is `approvaltest.NewSigner(seed).Sign(…)`; every passkey assertion
   is `approvaltest.NewPasskey(alg, seed).Assert(…)` or A2's recorded Chrome fixture (`host/approval/testdata/webauthn/
   chrome-{es256,eddsa}.json`, read through the existing `chromeAssertionInput`). No store, no daemon, no `/dev/tty`
   (the prompt tests write to a `bytes.Buffer`), no browser, no `~/.ssh`, no `~/.ailang/world`, no mint, no publish.
2. **Pure (design §0.2, §2.4 option (ii), §3).** No route (`TestFrozenV1RouteTableMatchesMux`, `TestApprovalRoutesRegistered`,
   `TestDecisionRoutesRegistered` unedited and green), no store read or write (`host/store` untouched; nothing new imports
   it), **no credentials schema change** (schema 9 is sprint 3's), no daemon process change, no `.ail`. `host/approval`
   stays stdlib-only (`TestApprovalImportsOnlyTheStandardLibrary`, unedited). `host/broker/finish_prompt.go` imports
   `fmt` and `strings` only.
3. **The CLI's output is unchanged.** `finishPrompt` writes the bytes the base tree's `finishPrompt` wrote (VT1's golden,
   generated at the base, §1.2 C9); every `episode_finish*_test.go` and `preset_committer_test.go` passes **unedited**
   (§6 `prompt` leg). The typed `y` path is unchanged (`episode_finish.go` untouched).
4. **One mutant per new or changed production branch, anchored to the diff** (§4); the named test and the behavioural
   suites with every static test `-skip`ped; every run appends a ledger line; the table's ids are the ledger's.
5. **No `t.Parallel()`** in any new test (`TestNoParallelWriteForkPackagesOutsideVerifygate`, `TestBrokerTestsDoNotCallParallel`).
6. **No git writes by the executor** (§7). The controller commits from snapshots.

---

## 1. Pre-flight

### 1.1 Baseline at `d3b0a813` (measured by the planner, 2026-11-17)

| Leg (`tools/gates.sh <base-export> <logs> base`, §6's script in base mode) | rc | Time |
|---|---|---|
| `status` | **0** | 1s |
| `verify_ail` | **0** | 41s |
| `gofmt` | **0** | 0s |
| `build` | **0** | 4s |
| `vet` | **0** | 2s |
| `fence` | **0** | 27s |
| `pkgs` | **1** | 246s |
| `race` | **0** | 58s |
| `prompt` | **0** | 15s |
| `inv` (named PASS: 33) | **0** | 15s |
| `verify_go` | **1** | 569s |
| `GATES_DONE 979s FAILS: pkgs verify_go` | | |

Every command ran with `export AILANG_BIN=$HOME/.pinned-ailang/ailang` (AILANG v0.41.0), `go1.26.6 darwin/arm64`, in a
`git archive d3b0a813` export at `plan-scratch/base-gates` (a throwaway git repo; `diff -r cmd host docs` against the
worktree: identical). Logs: `plan-scratch/logs/base_gates/`.

- `verify_ail`: `✓ verify gate PASSED: 16 required identities verified, 40 named tests pass, 503 se-tools named tests pass`.
- `pkgs` rc 1: `ok` `host/approval`, `host/approval/approvaltest`, `host/authority`, `cmd/ailang-worldd` (131.7 s); red
  only `host/broker` **`TestWorktreeSnapshotRefusesAConfiguredFilter`** (the directive's known list).
- `race`, `prompt` (the 13 existing finish/preset prompt tests by name), `inv` (33 named `--- PASS`, none missing) rc 0.
- `verify_go` rc 1: **62** packages `ok`; red only `host/broker` **`TestPublishApprovalRefusalSetWithALandedPositiveControl`**
  and `host/verifygate` **`TestMissionConfigDefaultFlowIntegration`** (both on the directive's known list).

**Rule (carried):** a red at the base measures the repo, not the sprint. The controller re-runs a red untouched test
alone once and does not chase it. **Never start a gate under `nohup`.**

### 1.2 Codebase claims checked (at `d3b0a813` unless "work" — the planner's final stage tree `work` at `M3`)

Commands: `plan-scratch/logs/measure.txt` (one batch, verbatim output) unless stated.

| # | Claim | Command → observed output |
|---|---|---|
| C1 | The base and the tools | `git rev-parse HEAD` (worktree) → `d3b0a8139bdecc36f8d083f4b701aef822340b27`; `git status --porcelain \| wc -l` → `0` (before the two plan files); `$HOME/.pinned-ailang/ailang --version` → `AILANG v0.41.0` (`Commit: 24ee108`); `go version` → `go1.26.6 darwin/arm64`; `/usr/bin/patch --version` → `patch 2.0-12u11-Apple` |
| C2 | **Directive input 1: no code drift since the design's base** (negative, with control) | `git diff --stat 994ba123 d3b0a813 -- cmd host docs \| tail -1` → *(empty)*; control `… -- design_docs \| tail -1` → `7 files changed, 1189 insertions(+), 106 deletions(-)` |
| C3 | **Directive input 1: every file:line the design cites for M1–M3, re-measured — none moved** (C2 makes it so; each re-read) | `admit.go:44` `func Admit(in AdmitInput) Verdict {` (`:37` `type AdmitInput struct {`); `signed.go:75` `func checkApprover(admit AdmitInput, a Approver, format string) error {`, `:58` `ApproveBySignature`, `:17` `SignedApproval`, `:39` `NotAdmittedError`, `:46` `ErrNotVerified`, `:49` `BindingError`; `passkey_assert.go:46` `… && ad.SignCount <= p.SignCount { // W15`, `:60` `func assertionSteps(response []byte, a Approver, h [32]byte, hErr error) (AuthData, error) {`; `decision_signed.go:79-84` `func DecideByAssertion(…)` … `checkApprover(AdmitInput{Now: in.Now}, in.Approver, FormatWebAuthn) … // P2-P5, P3w` … `h, hErr := DecisionMessageHash(…)` … `ad, err := assertionSteps(in.Response, in.Approver, h, hErr) // W1-W14` (W60 holds); `judge_decision.go:31` `enrolled, ok, err := ro.ApproverByKey(ctx, sig.PublicKey.Blob)`, `:58` `a, p, ok, err := ro.PasskeyByCredentialID(ctx, resp.RawID)`; `judge.go:21-26` the six `JudgeStep*`, `:30` `ErrNotEnrolled`, `:34` `type JudgeError struct {`; `episode_finish.go:270` `fmt.Fprint(tty, "Run this check? [y/N] ")`, `:410` `finishPrompt(tty, scope, actualTitle, rendered, protected, pinned.Committer)`, `:413` `return fail("pull request declined", "opening a PR requires typed y at the terminal; commit kept at "+sha, steps(sha))`; `episode_finish_operator.go:154` `func finishPrompt(tty io.Writer, scope broker.FinishApprovalScope, title, body string, protected []string, committer ...string) {`, `:170` the format string ending `protected paths changed: %s\nOpen this pull request? [y/N] `; `committer.go:3` `func displayCommitter(value string) string {` (the `""` → `(repository git config)` rule); `approve.go:689` `type FinishApprovalScope struct {`, `:698` `FormatFinishApprovalScope`, `:704` `ParseFinishApprovalScope`, `:722` `RecordOperatorApproval`; `episode_markdown.go:11` `const reportMarkdownMax = 60000` |
| C4 | **Directive input 2: the analogue's shape at HEAD** (what M1 mirrors) | `decision.go:15` `const DecisionPrefix = "world-decision/v1"`, `:18` `const DecisionNamespace = "world-decision@sunholo"`, `:21` `const DecisionMessageLen = len(DecisionPrefix) + 1 + 1 + 32 + 32 + 8`; `decision_signed.go`: `SignedDecision` (unexported fields), `DecideBySignature` = `checkApprover(AdmitInput{Now: in.Now}, …, FormatSSH)` then `VerifyDecisionSSHSIG`; `DecideByAssertion` = W1–W14 **without** W15 (deferred to the writer); `Binds` a `switch` over `!s.ok`/fields; `DecidedBy`, `Use`, `PasskeyUse` guarded by `s == nil \|\| !s.ok`. **Where V3 differs (design §2.3):** `ConfirmBySignature`/`ConfirmByAssertion` run **P1** (`Admit`, as the mint laws do: the verb request is a stored row) and **W15 in full** (no deferred-W15 two-path judge: the counter moves in sprint 3's consume, after the world write); only `CheckConfirmation` is A4-shaped (no P1, no W15). S1-D8 |
| C5 | **Directive input 4: the namespace and prefix constants at HEAD; the new ones differ** | `request.go:26` `const MessagePrefix = "world-approval/v1"`, `sshsig.go:19` `const Namespace = "world-approval@sunholo"`, `decision.go:15`/`:18` (C4). Work: `VerbPrefix = "world-verb/v1"` (13 bytes), `VerbNamespace = "world-verb@sunholo"`; the three prefixes share `world-` and differ at byte 6 (`a`/`d`/`v`), and VL3 (`TestVerbMessageLayout`) asserts it, and that neither of the other two messages prefixes the verb message nor the reverse; the three namespaces are pairwise distinct (VL3) and inside the signed blob (`ErrNamespace`, VL4) |
| C6 | **The names the patches use exist at HEAD or are new** (retro 10) | at HEAD: `approval.Admit`/`AdmitInput`/`Admitted`/`Verdict`/`RefuseConsumed`… (`admit.go`), `checkApprover`/`NotAdmittedError`/`ErrNotVerified`/`BindingError` (`signed.go`), `assertionSteps`/`ErrSignCountRegression`/`ErrApproverNotPasskey`/`ErrUserHandle`/`ErrBackupEligibility` (`passkey_assert.go`), `ParseSSHSIG`/`VerifySSHSIG`/`SSHPublicKey`/`ErrNamespace`/`ErrArmor`/`ErrKeyMismatch`/`ErrBadSignature`/`Namespace` (`sshsig.go`), `FormatSSH`/`FormatWebAuthn`/`Approver`/`ErrApproverName`/`ErrApproverFormat`/`ErrApproverRevoked`/`ErrApproverClockBehind` (`approver.go`), `Passkey`/`COSEKey.Canonical` (`passkey.go`, `cose.go`), `episodePattern`/`stringLeaves`/`ErrSchema`/`ErrString`/`ErrEpisode`/`ErrDecode`/`ErrNotCanonical`/`ErrNegativeExpiry`/`SignedMessage`/`MessageHash`/`MessagePrefix` (`request.go`), `DecisionMessage`/`DecisionMessageHash`/`DecideBySignature`/`DecideByAssertion` (`decision*.go`); test helpers `paApprover` (`passkey_law_test.go:18`), `paOrigin`/`paRPID` (`passkey_testauth_test.go:21`), `chromeAssertionInput` (`passkey_assert_test.go:17`), `approvaltest.NewSigner`/`Sign`/`PublicKey` (`signer.go:21,46,35`), `approvaltest.NewPasskey`/`Assert`/`CredentialID`/`Edit` (`passkey.go:30,96,50,65`); `broker.FinishApprovalScope` (`approve.go:689`); cmd `finishPrompt` (`episode_finish_operator.go:154`), `displayCommitter` (`committer.go:3`). **New**: `approval.VerbRequestSchema`, `VerbRequestLifetime`, `MaxVerbRequestBytes`, `VerbFinish`, `VerbPrefix`, `VerbNamespace`, `VerbMessageLen`, `VerbRequest`, `ErrVerbRequest`, `(VerbRequest).Validate`, `CanonicalVerb`, `DecodeVerb`, `VerbMessage`, `VerbMessageHash` (+ unexported `verbBodyMax`, `verbArgvMax`), `SignedVerb` (+ `Binds`, `ApprovedBy`, `Use`, `PasskeyUse`), `VerbSignatureInput`, `VerbAssertionInput`, `VerbCheckInput`, `ConfirmBySignature`, `ConfirmByAssertion`, `CheckConfirmation`, `ErrVerbFormat` (+ unexported `verifyVerbSSHSIG`); `broker.DisplayCommitter`, `broker.FinishPromptText`; tests `vrRequest`, `vrEscaped`, `vrLen`, `vAdmit`, `vMsg`, `vHash`, `vsInput`, `vaInput`, `vcSSH`, `vcPasskey`, `vPasskey`, `finishPromptFixture`, `finishPromptFixtures`, `finishPromptGolden` |
| C7 | **No verb symbol exists yet** (negative, with control) | `grep -rn 'VerbRequest\|world-verb\|FinishPromptText\|DisplayCommitter\|JudgeVerb\|RecheckVerb\|CheckConfirmation' cmd host --include='*.go' \| wc -l` → `0`; control `grep -rn 'DecisionNamespace' host --include='*.go' \| wc -l` → `23` |
| C8 | **The finish prompt's bytes were only partly pinned at the base** (why VT1 adds a golden) | `grep -rln 'Open this pull request' cmd host --include='*_test.go'` → `cmd/ailang-worldd/episode_finish_m3_test.go` (a `strings.Contains` of the question); `episode_finish_m2_test.go:310` `"committer: Test Operator <operator@example.invalid>"` (one line); nothing else pins the prompt's lines. `displayCommitter(` outside tests: the definition (`committer.go:3`) and **two** callers (`episode_finish_operator.go:159`, `episode_preset.go:255`) |
| C9 | **Directive input 5: VT1's fixtures generated by the real `finishPrompt` at the base** | `cp -cR base base-gen; cp gen/finish_prompt_golden_test.go base-gen/cmd/ailang-worldd/; cd base-gen; AILANG_BIN=$HOME/.pinned-ailang/ailang WORLD_UPDATE_FINISH_PROMPT_GOLDEN=1 go test -count=1 ./cmd/ailang-worldd/ -run '^TestFinishPromptIsTheGolden$' -v` → `--- PASS: TestFinishPromptIsTheGolden`, writing `host/broker/testdata/finish-prompt.golden.json` (sha256 `443911b0…e88594`); the same test without the variable on the base tree → `--- PASS` (the golden is the **base** `finishPrompt`'s bytes, by its own `fmt.Fprintf`s). The golden shipped in M3 is that file, byte for byte (C12); the test that generated it is shipped in M3 (plus three lines pinning `displayCommitter("")`, §4 `CMD-DC`). At `M3`, the same test passes against the delegating `finishPrompt`, and the broker's `TestFinishPromptText` reads the same file |
| C10 | **Every sub-milestone is committed and measured** | `tools/mkpatches.sh` + `git diff -M --numstat BASE <tag> -- cmd host docs` in `work`: `M1a 387`, `M1b 849`, `M1c 1034`, `M3 1300` (`-- cmd host`: the same; no `docs/`). The cut stage `S1b-M2` (`M2` over `M3`): **326** |
| C11 | **Every sub-milestone is green when committed** | `tools/accept_all.sh M1a M1b M1c M3` (an APFS clone of `work` at each tag, §3's `accept.sh`; `logs/accept_<M>.log`) → `M1a rc=0 ACCEPT_OK M1a`; `M1b rc=0 ACCEPT_OK M1b`; `M1c rc=0 ACCEPT_OK M1c`; `M3 rc=0 ACCEPT_OK M3` (09:57–10:02Z, ≈ 4.5 minutes for the four; each arm well inside one 600 s call). Re-run on the **final** stages (after M1c gained VL8's `ssh a negative expiry` row for mutant `S-MSG-ERR`, and M3 the `displayCommitter("")` pin for `CMD-DC`): `M1a rc=0 ACCEPT_OK M1a`; `M1b rc=0 ACCEPT_OK M1b`; `M1c rc=0 ACCEPT_OK M1c`; `M3 rc=0 ACCEPT_OK M3` (10:33–10:37Z, load ≈ 20, `logs/accept_all.txt`) |
| C12 | **The patches reproduce the stage tree** | `tools/roundtrip.sh <this plan>`: a fresh `git archive d3b0a813` export (from the worktree) + `patch -p1 -E` (after `--dry-run`) of Appendix P's four patches **extracted from this file with §2's `awk`**, in order (each `cmp`-equal to the generated patch; no `.orig`/`.rej`), each committed and measured → `CMP_EQUAL M1a`; `NUMSTAT M1a 387`; `CMP_EQUAL M1b`; `NUMSTAT M1b 849`; `CMP_EQUAL M1c`; `NUMSTAT M1c 1034`; `CMP_EQUAL M3`; `NUMSTAT M3 1300`; `ROUNDTRIP_IDENTICAL_TO_STAGE_TREE` (`diff -r -x .git` against `work` at `M3`; `logs/roundtrip.log`). Run on the checkpointed plan (10:09Z, before the mutant runs) and again on this final file, last (after M1c and M3 gained the two mutant-driven test lines and the mutant table was filled) |
| C13 | **Every mutant; the table's ids are the ledger's** | 115 rows (3 `CTL`, 3 `E`, 105 `K`, 3 `N`, 1 `S`); 118 ledger lines (runs R1, R2); 115 of 115 rows run; the last line of each row behaves as expected for **115 of 115**. Table ids = ledger ids (`tools/mktable.py` fills §4 from the ledger's last line per id); ledger `plan-scratch/logs/mutant_ledger.tsv` |
| C14 | **The gates, on the final tree** | `tools/gates.sh <final-tree> logs/final_gates final` (`final-tree` = `work` at `M3` `fd45eb58`, a git repo with every stage committed, byte-identical to the round-trip tree; load 13–25, the planner's mutant runs beside it): `status` 0; `verify_ail` rc 0 (no `.ail` changed; `✓ verify gate PASSED: 16 required identities verified, 40 named tests pass, 503 se-tools named tests pass`, as the base); `gofmt`, `build`, **`vet`**, `fence` rc 0; **`pkgs` rc 1** — `ok` `host/approval` (1.5 s), `host/approval/approvaltest`, `host/authority`, `cmd/ailang-worldd` (289.8 s); red only `host/broker` **`TestExecProbeNet6ForeignDialIsRetriedNotRefused`** (an exec-probe timing test, untouched, the directive's "host/broker exec-probe timeouts") — **green re-run alone** (`--- PASS: TestExecProbeNet6ForeignDialIsRetriedNotRefused (6.45s)`, `logs/final_broker_red_alone.txt`) and green in this run's own `verify_go` (`host/broker` `ok`); the base's `TestWorktreeSnapshotRefusesAConfiguredFilter` passed here; **`race` rc 0**; **`prompt` rc 0** (the 13 existing finish/preset prompt tests, unedited, + `TestFinishPromptIsTheGolden`: every one `--- PASS`); **`inv` rc 0: 47 named `--- PASS`**, 0 SKIP, 0 FAIL, none missing (the base's 33 + the 14 new: the ten `host/approval` tests, `TestFinishPromptText`, `TestDisplayCommitter`, and the two authority CLI-judge inventories added to the list); **`verify_go` rc 0: 64 packages `ok`** (the base's two `verify_go` reds passed here). **`GATES_DONE 1456s FAILS: pkgs`** — the one red an untouched load-timing test, green alone; `ALL_DONE.fd45eb584c720f405ca9076d16c1d2421ebb023b` written (`logs/final_gates/`) |
| C15 | **Exactly these pre-existing test files change: none** | work: `git diff --name-status BASE M3 -- '*_test.go'` → `A cmd/ailang-worldd/finish_prompt_golden_test.go`, `A host/approval/verb_check_test.go`, `A host/approval/verb_chrome_test.go`, `A host/approval/verb_signed_test.go`, `A host/approval/verb_test.go`, `A host/broker/finish_prompt_test.go` — six `A`, **zero `M`** |
| C16 | **No route, no store, no schema, no daemon** | work: `git diff --name-only BASE M3 -- host/daemon host/store host/surfacepage host/verbs host/authority world packages scripts docs` → *(empty)*; `grep -rn 'host/store\|net/http\|database/sql' host/approval/verb.go host/approval/verb_signed.go host/broker/finish_prompt.go` → *(empty)*; control `grep -c 'host/store' host/authority/judge.go` → `1` |
| C17 | **The snapshot line works on a sub-milestone with no `docs/` and no `cmd/`** | `tools/snap_snippet.sh` (§7's block, parameterised) on a `git clone --shared` of the worktree checked out at `d3b0a813` with the patches applied (uncommitted, as the executor's tree): at **M1a** (`host/` only — no `cmd/`, no `docs/`) → `STATUS` 2 lines (`A host/approval/verb.go`, `A host/approval/verb_test.go`), `DELETED: 0`, `MANIFEST: 2`, no error; at **M3** → `STATUS` 12 lines, `DELETED: 0`, `MANIFEST: 12`; `shasum -a 256 -c MANIFEST.sha256` → 0 lines not `OK` |
| C18 | **Directive input 6: the size** | design §8 S1 ≈ **935** raw; compiled **1,626** (M1 1,034 + M2 326 + M3 266) = 1.74× the design's raw, over the 1,400 cap (and over 935 × 1.35 = 1,262). The over-run is test (the design's ≈ 530 test lines measured 1,210 for the three milestones: every Validate rule, P1–P5, W1–W15 and both directions of domain separation have their own row, retro 3) and the law's production (≈ 405 designed, 554 measured). Cut: S1-D2 |

### 1.3 Planner decisions where the design leaves a choice

- **S1-D1 — Four sub-milestones, in dependency order.** **M1a** `verb.go` (request, canonical bytes, strict decode,
  `Validate`, the message) + VL1–VL3; **M1b** `verb_signed.go`'s laws and `SignedVerb` + VL4–VL7; **M1c**
  `CheckConfirmation` (+ `VerbCheckInput`, `ErrVerbFormat`) + VL8 and the recorded-Chrome half of VL4; **M3**
  `FinishPromptText`/`DisplayCommitter` and the cmd delegations + VT1. Each is green alone (C11). M1b is 462 lines (its
  test file carries the whole P1–P6 / W1–W15 matrix); splitting it further would split one law's table across commits.
- **S1-D2 — The cut: M2 → sprint S1b (directive "Size rule").** Compiled, S1 = 1,626 > 1,400. The options measured:
  M1 + M2 = 1,360 (margin 40, and M3 then waits); **M1 + M3 = 1,300 (margin 100)**; trimming tests (−225) would drop
  behavioural rows. **Chosen: M1 + M3.** Reasons: (1) M1's law is fully proven either way (VL1–VL8); (2) M3 is a
  prerequisite of sprint 2 (XC1: the card's prompt equals `FinishPromptText` of the fixture) and of sprint 3 (F5), while
  the judges' first caller is sprint 4's C11 (and sprint 5's CF6), so M2 has the most slack; (3) M2 is self-contained
  (`host/authority/judge_verb.go` + its test, no file shared with S1), compiled and green on top of `M3` (`S1b-M2`,
  326 lines: `judge_verb.go` 129, `judge_verb_test.go` 197; `TestJudgeVerbSSH`, `TestJudgeVerbAssertion`,
  `TestRecheckVerb` PASS), its patch at `plan-scratch/patches/S1b-M2.patch` (sha256 prefix `775b927ad56aa4c0`).
  **S1b** = that patch + VJ1's mutants, ≈ 326 lines; the controller routes it (it may ride with sprint 2 only if
  sprint 2 measures ≤ 1,074). Appendix C.
- **S1-D3 — VT2 is sprint 3's (directive: "resolve by reading").** Design §6 "Sprint 1" lists **VT2
  `FinishConfirmRefusal` F1–F5 (and F6 with `DraftOnly`)**; design §2.3 says `broker.FinishConfirmRefusal` is "built in
  sprint 3 **[r1]**, beside its only caller's store"; §8's table puts "M8 `FinishConfirmRefusal` (≈ 40 + 50)" in **S3**,
  and §8's S1 row lists only M1–M3 (≈ 935, which has no room for M8). The `[r1]` marker and §8 are the r1 text (r1 moved
  the store and everything beside it to sprint 3, §15 row 1); §6's VT2 line carries no `[r1]` and is the r0 placement
  left behind. **The directive's reading is confirmed; followed: VT2 is not built here** (Appendix C: S3 inherits it,
  without F6/`DraftOnly`, S1-D4).
- **S1-D4 — D-WORLD-118 = A: the B-only branches are not built** (ledger `72007b37`, `verdict=attended`). Design §2.3's
  **F6** ("D-WORLD-118 = B: a ready pull request is confirmed with a typed y"), VT2's "F6 with `DraftOnly`", §2.6 C10's
  "F6 under B", §2.7's `--ready` refusal and §8's "B +≈ 15" (S4, S5) are **dropped** from every sprint. Sprints 3–6 are
  now routable (later iterations); this sprint is unchanged by the answer (it touches no store, no route, no schema).
- **S1-D5 — Validate's rules, order and classes** (design §2.2, one comment-labelled branch each, V1–V10): V1 schema
  (`ErrSchema`), V2 every string leaf valid UTF-8 (`ErrString`, via the existing `stringLeaves`, naming the JSON path),
  V3 verb = `VerbFinish` (`ErrVerbRequest`), V4 episode `^[a-z0-9][a-z0-9-]{0,63}$` (the existing `episodePattern`,
  `ErrEpisode`), V5 `argv[0:4] == ailang-worldd episode finish <episode>`, V6 ≤ 64 words, V7 no byte `< 0x20` or
  `0x7f` in any word, V8 prompt, scope, title non-empty, V9 body ≤ 60,000 bytes, V10 canonical bytes ≤
  `MaxVerbRequestBytes` (131,072). New exported names: `VerbFinish` (the one admissible verb: sprint 5's CLI builds
  `Verb: approval.VerbFinish`) and `ErrVerbRequest` (the class of V3, V5–V10); the bounds 64 and 60,000 stay unexported
  (`verbArgvMax`, `verbBodyMax`), each pinned at its boundary by VL2. **As designed, V2 checks UTF-8 only** (the mint
  rule also refuses NUL; the design's verb rule does not, and V7 refuses NUL in argv): §1.4 P6.
- **S1-D6 — `FinishPromptText`/`DisplayCommitter` in a new `host/broker/finish_prompt.go`** (design §7 allows "or a new
  `finish_prompt.go`"), so `approve.go` — A4's file — is untouched. The text is moved **verbatim** (one `Sprintf` for
  the committer line, one for the rest, the same verbs and arguments); `finishPrompt` keeps its variadic signature and
  its one caller (`episode_finish.go:410`) is untouched; `displayCommitter` becomes a one-line delegation (its other
  caller, `episode_preset.go:255`, is unchanged).
- **S1-D7 — VT1's golden: generated by the base's `finishPrompt`, read by both packages.** The generator is the shipped
  `TestFinishPromptIsTheGolden` (cmd), which under `WORLD_UPDATE_FINISH_PROMPT_GOLDEN=1` writes
  `host/broker/testdata/finish-prompt.golden.json` from `finishPrompt`'s bytes, and otherwise compares; it ran on the
  base tree to make the golden (C9). The three fixtures cover draft and ready, an existing and a created label, no /
  one / two protected paths, a pinned and an empty committer, a body over six lines (the seventh hidden), under six,
  and empty, a title with `<`, `&`, `"` and non-ASCII. `TestFinishPromptText` (broker) asserts
  `FinishPromptText(…) + "Open this pull request? [y/N] "` equals each golden prompt and that the golden still covers
  those classes. The executor never sets the variable (§7).
- **S1-D8 — The law's shape** (directive input 2). `SignedVerb` mirrors `SignedApproval`/`SignedDecision` (unexported
  fields; `Binds` B1–B4 as a `switch`; `ApprovedBy`/`Use`/`PasskeyUse` guarded). `ConfirmBySignature` = P1 → P2–P5 →
  P6 (`verifyVerbSSHSIG`, unexported, under `VerbNamespace`); `ConfirmByAssertion` = P1 → P2–P5, P3w → W1–W14
  (`assertionSteps(…, VerbMessageHash(…), hErr)`) → W15 (the mint law's text and arithmetic, `passkey_assert.go:46`).
  `CheckConfirmation(VerbCheckInput{Digest, Challenge, Expiry, Now, Artifact, Format, Approver})` dispatches on
  `Format` (`ssh` → P2–P5 + P6; `webauthn` → P2–P5, P3w + W1–W14; anything else → `ErrVerbFormat`) with
  `checkApprover(AdmitInput{Now: in.Now}, …)` (A4's form, `decision_signed.go:80`) — **no P1, no W15**, and the
  message built **only from the caller's own** digest, challenge and expiry.
- **S1-D9 — VL4 both directions, with both key kinds** (directive input 4). SSH: a mint SSHSIG, a decision SSHSIG, the
  verb message under the mint and the decision namespaces → `ErrNamespace`; the mint and decision messages under the
  verb namespace and a verb message missing its `0x00` → `ErrBadSignature`; a verb SSHSIG → `ApproveBySignature` and
  `DecideBySignature` refuse with `ErrNamespace`. Passkey: a mint and a decision assertion → `ConfirmByAssertion`
  `ErrChallenge` (W6); a verb assertion → `ApproveByAssertion` and `DecideByAssertion` `ErrChallenge` (W6); and A2's
  **recorded Chrome assertions** (a mint request's hash) → `ConfirmByAssertion` and `CheckConfirmation` `ErrChallenge`,
  with `ApproveByAssertion` of the same input as the positive control (`TestChromeMintAssertionIsNotAVerb`).
- **S1-D10 — VL8 pins audit note 5 exactly as designed** (directive input 3): `CheckConfirmation` judges P2–P5 at
  `in.Now`, so an approver revoked after the confirmation (`RevokedAt = expiry + 1`) is **refused**
  (`ErrApproverRevoked`), for both formats. The ≈ 1 s window this creates for the CLI (fail-safe: commit kept, a
  re-run asks afresh) is sprint 5's to test (CF6); Appendix C carries it.
- **S1-D11 — The re-check is proven against a moved counter, both ways**: stored counter = the artifact's (2) and past
  it (7) both verify (kills "W15 left in", §4 `CC-W15-IN`); a counter-less authenticator (0 over 0) passes
  `ConfirmByAssertion`'s W15 (kills the zero exemption's deletion, `A-W15-ZERO`).
- **S1-D12 — New tests in new files** (no pre-existing test file edited, C15): `verb_test.go` (M1a), `verb_signed_test.go`
  (M1b), `verb_check_test.go` and `verb_chrome_test.go` (M1c; the second is `package approval` because
  `chromeAssertionInput` is internal), `finish_prompt_test.go` and `finish_prompt_golden_test.go` (M3).
- **S1-D13 — S1b's RecheckVerb (recorded for S1b, not built here):** `RecheckVerb(ctx, ro, in, artifact)` keeps the
  design's signature by detecting the artifact's kind — an armored SSHSIG (`ParseSSHSIG` succeeds), else a passkey
  assertion (`ParseAssertionResponse` succeeds), else `JudgeStepParse` with `ErrUnknownArtifact` — then looks up the
  approver as `judge_decision.go:31`/`:58` do and calls `CheckConfirmation` with `Expiry = in.ExpiresAt`; it never
  reads `in.CreatedAt`, `BytesDigest` or `Consumed`. Sprint 5 compares `ApprovedBy()` with the row's `approved_by`.

### 1.4 Premises found false or incomplete at HEAD (each routed around minimally)

- **P1 — Design §8 "S1 ≈ 935 raw" (and 935 × 1.35 ≈ 1,262).** Compiled 1,626 (C18). Routed: the cut (S1-D2), 1,300.
- **P2 — Design §6 "Sprint 1 … VT2 `FinishConfirmRefusal` F1–F5".** Stale r0 text (S1-D3). Routed: VT2 to sprint 3.
- **P3 — Design W45 / §1 item 1 "`displayCommitter` … 3 callers".** Measured: the grep's 3 lines are the definition and
  **two** callers (C8). Routed: nothing changes (both callers keep `displayCommitter`, which delegates).
- **P4 — Design §2.2 / W43 "the existing `episode_finish_m2/m3` prompt assertions" as the byte-identity proof.** They
  pin only the question (`strings.Contains`) and one `committer:` line (C8); a moved line, a dropped preview line or a
  changed label would pass them. Routed: VT1's golden (S1-D7) is the byte-identity proof; m2/m3 stay unedited and green
  as the behavioural check (§6 `prompt`).
- **P5 — Design VL4 "a mint assertion fails `ConfirmByAssertion` at W6 and back" with "the A2 recorded authenticator
  fixtures".** The recorded Chrome assertions signed a **mint** message hash (`passkey_assert_test.go:17`), so they can
  only drive the mint→verb direction. Routed: the reverse (a verb assertion refused by `ApproveByAssertion` and
  `DecideByAssertion`) uses `approvaltest.NewPasskey` (S1-D9).
- **P6 — Design §2.2 Validate "every string valid UTF-8".** The mint request's rule also refuses NUL
  (`request.go`'s `ErrString`: "not valid UTF-8 or holds NUL"); the designed verb rule does not. Built as designed (a NUL
  in a title or body is valid UTF-8 and is rendered by the card as text; argv refuses it under V7). Recorded for S2's
  card (Appendix C), not changed here (no re-design).
- **P7 — Design §2.3 `SignedVerb`'s comment lists `ok, digest, challenge, expiry, name, format, key, credentialID,
  signCount`.** Built exactly so; `ConfirmBySignature` sets `format` to `FormatSSH` (the approver's format is `ssh`
  after P3), not `in.Approver.Format` as `ApproveBySignature` does — identical after P3, chosen so the mutant
  `SIG-FORMAT` is a real branch (§4).

### 1.5 The size

Lines are insertions + deletions: **`git diff -M --numstat d3b0a813 HEAD -- cmd host docs`** (git's rename detection on;
nothing is renamed). **Every row is measured** from the scratch stages, each committed in `work`, and re-measured on
the round-trip tree (C10, C12).

| Sub-milestone (in order) | Production | Test | Own | Cumulative vs `d3b0a813` (after its commit) | Run |
|---|---|---|---|---|---|
| M1a — `verb.go`; VL1–VL3 | 172 | 215 | 387 | **387** | 1 |
| M1b — `verb_signed.go` laws + `SignedVerb`; VL4–VL7 | 140 | 322 | 462 | **849** | 1 |
| M1c — `CheckConfirmation`; VL8, VL4 (recorded Chrome) | 47 (+50 −3) | 138 | 185 | **1,034** | 1 |
| M3 — `FinishPromptText`, `DisplayCommitter`, delegations; VT1 + golden | 63 (38 + 10 + 15) | 203 (49 + 87 + 67) | 266 | **1,300** | 1 |
| **The sprint** | **422** | **878** | | **1,300 ≤ 1,400** (margin 100) | |
| *cut: S1b — M2 `judge_verb.go` + VJ1 (over `M3`)* | *129* | *197* | *326* | *(1,626 had it stayed)* | — |

Per file, cumulative (+/−): `verb_signed_test.go` +322; `verb_test.go` +215; `verb_signed.go` +187; `verb.go` +172;
`verb_check_test.go` +110; `finish_prompt_golden_test.go` +87; `finish-prompt.golden.json` +67; `finish_prompt_test.go`
+49; `finish_prompt.go` +38; `verb_chrome_test.go` +28; `episode_finish_operator.go` +1 −14; `committer.go` +4 −6.

---

## 2. How the executor applies a sub-milestone (the patch protocol)

Each sub-milestone `<M>` is the block fenced by a line `` ````diff <M> `` (four backticks) and a line of four backticks,
in Appendix P. Paths are `a/<repo path>` and `b/<repo path>`; a new file's patch has `--- /dev/null` and an
`@@ -0,0 …` hunk. S1 deletes no file; `-E` is kept (the protocol is V2a's). Apply it **only** this way, in one `bash`
shell (the controller's shell is zsh):

```bash
bash <<'EOF2'
set -euo pipefail
M=M1a   # <- the sub-milestone, exactly as in the Appendix P heading, in order: M1a M1b M1c M3
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter400-row243-v3s1
PLAN="$W/design_docs/planned/w-world-address-bar-verbs-v3s1-sprint-plan.md"
T="$HOME/.ailang/state/world-iter400/exec-patches"; mkdir -p "$T"
awk -v m="$M" '$0 == "````diff " m {f=1; next} f && $0 == "````" {exit} f' "$PLAN" > "$T/$M.patch"
test -s "$T/$M.patch"
cd "$W"
/usr/bin/patch -p1 -E --dry-run < "$T/$M.patch"
/usr/bin/patch -p1 -E < "$T/$M.patch"
if find cmd host docs \( -name '*.orig' -o -name '*.rej' \) | grep .; then echo "STOP: patch left .orig/.rej"; exit 1; fi
git status --short --untracked-files=all
EOF2
```

- If `--dry-run` fails, **stop and report** (no hand edits, no fuzz options). It can only fail if an earlier
  sub-milestone was skipped, reordered or changed.
- Then run that sub-milestone's **Acceptance** (§3): extract `accept.sh` from this plan and run it, in one bash shell,
  `tee`'d to a log (never only through `tail`). **Never `nohup`, never background a test leg; pass a 600000 ms
  timeout to every long Bash call.**
- Then snapshot (§7). The controller commits and compares the numstat with §3's expected value.

---

## 3. Milestones

**Acceptance.** Every sub-milestone's acceptance is the arm of this script named for it, preceded by the same
**[base]** lines (gofmt, `go build ./...`, `go vet` of the touched packages, the `_test.go` compile fence). Each `v`
line must show `--- PASS` for every test it names and **no `--- FAIL`, no `--- SKIP`**; each `p` line must end `ok`.
The script ends `ACCEPT_OK <M>` or prints `ACCEPT_RED <M>: <step>` and exits 1. Each arm fits a 600 s call (measured:
the planner's runs, C11). Extract and run it (`M` and `W` as in §2):

```bash
bash <<'EOF2'
set -uo pipefail
M=M1a
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter400-row243-v3s1
PLAN="$W/design_docs/planned/w-world-address-bar-verbs-v3s1-sprint-plan.md"
L="$HOME/.ailang/state/world-iter400/exec-logs"; mkdir -p "$L"
awk '$0 == "````bash accept.sh" {f=1; next} f && $0 == "````" {exit} f' "$PLAN" > "$L/accept.sh"
/bin/bash "$L/accept.sh" "$W" "$M" 2>&1 | tee "$L/accept_$M.log"
EOF2
```

````bash accept.sh
#!/bin/bash
# accept.sh <tree> <M>: exactly plan §3's acceptance for sub-milestone M (one bash shell; rc 0 = ACCEPT_OK).
T=$1; M=$2
cd "$T" || exit 2
export AILANG_BIN=$HOME/.pinned-ailang/ailang
unset WORLD_UPDATE_FINISH_PROMPT_GOLDEN WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER
set -o pipefail
chk() { tee /dev/stderr | { out=$(cat); echo "$out" | grep -qE -- '--- (FAIL|SKIP)|^FAIL' && exit 1; echo "$out" | grep -q -- '--- PASS' || exit 1; }; }
v() { go test -count=1 -v "$@" 2>&1 | chk; }   # -v: every named test PASS, none SKIP/FAIL
p() { go test -count=1 "$@"; }
step() { echo "== $*"; "$@" || { echo "ACCEPT_RED $M: $*"; exit 1; }; }
# [base]
step test -z "$(gofmt -l cmd host)"
step go build ./...
step go vet ./host/approval/... ./host/authority/ ./host/broker/ ./cmd/ailang-worldd/
if go test -run '^$' ./... 2>&1 | grep -E '^(FAIL|--- FAIL)'; then echo "ACCEPT_RED $M: FENCE_RED"; exit 1; fi
A0='TestApprovalImportsOnlyTheStandardLibrary|TestApprovalTestIsTestOnly|TestMintSignatureIsNotADecision|TestDecideByAssertionDefersTheCounter|TestApproveByAssertionOrder|TestApproveByAssertionAcceptsChrome'
VL13='TestVerbRequestRoundTrip|TestVerbRequestDecodeRefusals|TestVerbRequestValidate|TestVerbMessageLayout'
VL47='TestVerbDomainSeparation|TestConfirmBySignatureLaw|TestConfirmByAssertionLaw|TestSignedVerbBinds'
VL8='TestCheckConfirmation|TestChromeMintAssertionIsNotAVerb'
PROMPT='TestFinishPromptIsTheGolden|TestEpisodeFinishDeclineAC7|TestEpisodeFinishBytesAC8|TestEpisodeFinishDispatchAC9AC11AC12|TestEpisodeFinishOutcomesAC10|TestEpisodeFinishProtected|TestEpisodeFinishCredentialsMismatch|TestEpisodeFinishResumeAC13|TestEpisodeFinishRerunAC10|TestEpisodeFinishApprovalReuse|TestEpisodeFinishFixupRound|TestEpisodeFinishFinishedRerun|TestEpisodeFinishResumeGuards|TestPresetCommitterValidationAndPrompt'
case $M in
M1a) step v ./host/approval/ -run "^($VL13|$A0)\$"
     step p ./host/approval/... ;;
M1b) step v ./host/approval/ -run "^($VL13|$VL47|$A0)\$"
     step p ./host/approval/... ./host/authority/ ;;
M1c) step v ./host/approval/ -run "^($VL13|$VL47|$VL8|$A0)\$"
     step v ./host/authority/ -run '^(TestJudgeDecisionSSH|TestJudgeDecisionAssertion)$'
     step p ./host/approval/... ./host/authority/ ;;
M3)  step v ./host/broker/ -run '^(TestFinishPromptText|TestDisplayCommitter|TestBrokerTestsDoNotCallParallel|TestDecideApprovalCallerCensus)$'
     step v ./cmd/ailang-worldd/ -run "^($PROMPT|TestEveryVerbHasHelp|TestVerbTableMatchesTheCLISwitch)\$"
     step v ./host/approval/ -run "^($VL13|$VL47|$VL8)\$"
     step p ./host/broker/ -run '^(TestFinish|TestDisplayCommitter|TestRecordOperator|TestDecideApproval|TestParseFinish|TestFormatFinish)'
     step p ./cmd/ailang-worldd/ -run '^(TestEpisodeFinish|TestFinish|TestPreset|TestEpisodeStart)' ;;
*) echo "unknown $M"; exit 2 ;;
esac
echo "ACCEPT_OK $M"
````

### M1a — the verb request and the message (design §2.2) · 172 prod + 215 test · Run 1

**Patch** `M1a`: **`host/approval/verb.go`** (NEW) — the constants (`VerbRequestSchema`, `VerbRequestLifetime` 10 min,
`MaxVerbRequestBytes` 128 KiB, `VerbFinish`, `VerbPrefix`, `VerbNamespace`, `VerbMessageLen` = 86, unexported
`verbBodyMax` 60,000 and `verbArgvMax` 64); `VerbRequest` (design §2.2's ten fields in order, json tags);
`ErrVerbRequest`; `Validate` (V1–V10, S1-D5); `CanonicalVerb` (`Canonical`'s encoder settings); `DecodeVerb`
(`Decode`'s strictness); `VerbMessage` (prefix ‖ 0x00 ‖ digest ‖ challenge ‖ u64be(expiry), a negative expiry
refused) and `VerbMessageHash`. **`host/approval/verb_test.go`** (NEW): **VL1** `TestVerbRequestRoundTrip` (the exact
canonical bytes, `<`/`&` unescaped, no trailing newline; decode inverts; a `nil` `protected` is `null`) and
`TestVerbRequestDecodeRefusals` (unknown fields and non-JSON → `ErrDecode`; trailing data, a second object, reordered
keys, whitespace, an HTML-escaped spelling → `ErrNotCanonical`); **VL2** `TestVerbRequestValidate` (one row per rule,
each bound on both sides: 4/64/65 words, `0x1f`/`0x20`/`0x7f`, 60,000/60,001 body bytes, 131,072/131,073 canonical
bytes); **VL3** `TestVerbMessageLayout` (86 bytes, every field's offset, big-endian expiry, the hash, negative expiry
for both functions, the prefixes' byte 6 and the three namespaces).
**Expected numstat after commit: 387.** Files: `?? host/approval/verb.go`, `?? host/approval/verb_test.go`.

### M1b — the laws and `SignedVerb` (design §2.3) · 140 prod + 322 test · Run 1

**Patch** `M1b`: **`host/approval/verb_signed.go`** (NEW) — `SignedVerb` (unexported fields), `VerbSignatureInput`,
`VerbAssertionInput`, `verifyVerbSSHSIG`, `ConfirmBySignature` (P1 → P2–P5 → P6), `ConfirmByAssertion` (P1 → P2–P5, P3w
→ W1–W14 → W15), `Binds`, `ApprovedBy`, `Use`, `PasskeyUse`. **`host/approval/verb_signed_test.go`** (NEW): **VL4**
`TestVerbDomainSeparation` (S1-D9, both directions, both key kinds), **VL5** `TestConfirmBySignatureLaw` (P1 × 4
verdicts, P1 before P4, P2, P3, P4, P4 before P6, P5, P6 not armored / another key / another digest / challenge /
expiry; the control's `Binds`, `ApprovedBy`, `Use`, no `PasskeyUse`), **VL6** `TestConfirmByAssertionLaw` (P1 × 4, P2,
P3w × 2, P4, P5, W1, W2, W3, W5, W6 × 3, W7–W14, W15 at the stored counter and at 0 over 1; the control's `PasskeyUse`
(credential ID, counter 2) and `Use`; a counter-less 0 over 0 passes), **VL7** `TestSignedVerbBinds` (each field, the
nil and zero values, `ApprovedBy`/`Use`/`PasskeyUse` empty unless verified, `Use` hands out a copy).
**Expected numstat after commit: 849.** Files: + `?? host/approval/verb_signed.go`, `?? host/approval/verb_signed_test.go`.

### M1c — `CheckConfirmation` (design §2.3 [r1]) · 47 prod + 138 test · Run 1

**Patch** `M1c`: **`verb_signed.go`** gains `VerbCheckInput`, `ErrVerbFormat`, `CheckConfirmation` (S1-D8) and the
header/`SignedVerb` comments name it. **`host/approval/verb_check_test.go`** (NEW): **VL8** `TestCheckConfirmation` —
SSH: revoked since (audit note 5, S1-D10), name, enrolled after now, a webauthn approver, another digest / challenge /
expiry, a negative expiry, a mint- and a decision-namespace SSHSIG, an assertion as an SSHSIG, another key; WebAuthn: revoked since, an
ssh approver, another digest / challenge / expiry, a negative expiry, a mint assertion, a bad signature; no and an
unknown format; the verified re-checks past the expiry (no P1) with the stored counter at and past the artifact's (no
W15), `ApprovedBy`, `PasskeyUse`, `Use`. **`host/approval/verb_chrome_test.go`** (NEW, `package approval`):
`TestChromeMintAssertionIsNotAVerb` (S1-D9, recorded fixtures).
**Expected numstat after commit: 1,034.** Files: + `M host/approval/verb_signed.go` (now cumulative `??`),
`?? host/approval/verb_check_test.go`, `?? host/approval/verb_chrome_test.go`.

### M3 — the finish prompt's text, moved (design §2.2) · 63 prod + 203 test · Run 1

**Patch** `M3`: **`host/broker/finish_prompt.go`** (NEW) — `DisplayCommitter`, `FinishPromptText` (S1-D6).
**`cmd/ailang-worldd/episode_finish_operator.go`** — `finishPrompt` prints `broker.FinishPromptText(…) + "Open this pull
request? [y/N] "` (and drops its `strings` import); **`cmd/ailang-worldd/committer.go`** — `displayCommitter` delegates.
**`host/broker/testdata/finish-prompt.golden.json`** (NEW, generated at the base, C9), **`host/broker/finish_prompt_test.go`**
(NEW: **VT1** `TestFinishPromptText`, `TestDisplayCommitter`), **`cmd/ailang-worldd/finish_prompt_golden_test.go`** (NEW:
**VT1** `TestFinishPromptIsTheGolden` — the delegating `finishPrompt` writes the golden's bytes; the omitted variadic
committer; `displayCommitter("")`). The existing finish/preset prompt tests run unedited (accept `M3`).
**Expected numstat after commit: 1,300.** Files: + `M cmd/ailang-worldd/committer.go`, `M cmd/ailang-worldd/episode_finish_operator.go`,
`?? cmd/ailang-worldd/finish_prompt_golden_test.go`, `?? host/broker/finish_prompt.go`, `?? host/broker/finish_prompt_test.go`,
`?? host/broker/testdata/finish-prompt.golden.json`.

### 3.1 Every input and failure class the law decides, traced (design §2.9; directive retro 4)

Sprint 1 adds no entry point: its code is reached by sprint 4's POST (C8, C11, C12 — through S1b's judges) and sprint
5's CLI (§2.7 step 6(a), CF6). Each design §2.9 / §2.6 / §4 row whose **decision** S1's code makes is driven here,
directly on the law:

| design row (class) | S1 branch | test that drives it (→ asserted) |
|---|---|---|
| §2.6 C8 consumed (T2, a replay on the same request) | P1 `Admit` → `RefuseConsumed` | VL5 `P1 consumed`, `P1 consumed before P4`; VL6 `P1 consumed` → `NotAdmittedError{RefuseConsumed}` |
| C8 clock behind / expired (T4) / stored digest ≠ SHA-256(bytes) | P1 | VL5/VL6 `P1 clock behind`, `P1 expired` (now = expiry and expiry + 1), `P1 digest` |
| C11 unknown key | S1b's judge (lookup); the law's P6 with another key | VL5 `P6 another key` → `ErrKeyMismatch` (S1b: `VJ-unknown-key`) |
| C11 revoked / wrong format / enrolled after now | P2–P5 (`checkApprover`) | VL5 `P2`–`P5`, `P4 revoked before P6`; VL6 `P2`, `P3w` × 2, `P4`, `P5` |
| C11 bad signature / namespace (T1, T5) | P6 `verifyVerbSSHSIG` | VL5 `P6 not armored` (`ErrArmor`), `another digest/challenge/expiry` (`ErrBadSignature`); VL4 (`ErrNamespace`, both directions) |
| C11 origin, RP hash, UV, UP, BE, user handle, type, cross-origin, envelope | W1–W14 (`assertionSteps`) | VL6 `W1`–`W14` (one row each, as A2's matrix) |
| C11 counter (T11's law half) | W15 | VL6 `W15 the stored counter`, `W15 a zero counter over a stored 1`; the 0-over-0 control |
| C12 / T3 the same artifact on another request | `Binds` B1–B4 + the signed digest/challenge/expiry | VL7 (each field, nil, zero); VL5/VL6 `another digest/challenge/expiry` |
| T5 a mint or decision signature confirms; a verb signature mints or decides | prefix byte 6, three namespaces | VL4 (SSH and passkey, both directions); `TestChromeMintAssertionIsNotAVerb` (recorded); VL3 (bytes) |
| C10 `DecodeVerb` + `Validate` (sprint 4 calls them) | V1–V10, `DecodeVerb`'s refusals | VL1, VL2 (each rule, each bound) |
| CF6 the re-check: digest/challenge/expiry not this process's | `CheckConfirmation` builds the message from its input only | VL8 `ssh/webauthn another digest/challenge/expiry`, `a negative expiry` |
| CF6 the artifact does not verify; another namespace; an assertion as an SSHSIG | P6 / W1–W14 | VL8 `bad signature`, `mint-namespace`, `decision-namespace`, `mint assertion`, `an assertion as an SSHSIG` (`ErrArmor`), `another key` |
| CF6 approver revoked since (audit note 5) | P2–P5 at now | VL8 `ssh revoked approver`, `webauthn revoked approver` → `ErrApproverRevoked` (fail-safe) |
| CF6 a consumed row past its expiry; a moved counter | no P1, no W15 | VL8's verified rows (now = expiry + 5; stored counter 2 and 7) |
| CF6 a format the re-check does not know | `ErrVerbFormat` | VL8 `no format`, `an unknown format` |
| CF6 `approved_by` differs (sprint 5 compares) | `ApprovedBy` = format ":" name | VL5/VL6/VL8 controls (`ssh:mark-ssh`, `webauthn:mark-mac`); VL7 zero → `""` |
| C10 F5 / T8 / CF13 "the prompt is the terminal's rendering" | `FinishPromptText` | VT1 (`TestFinishPromptText`, `TestFinishPromptIsTheGolden`: byte-equal to the base `finishPrompt`) |
| the typed path (no flag) | unchanged | every existing finish/preset test, unedited (§6 `prompt`) |

---

## 4. Mutation table

**Procedure** (`tools/mutate.py <run> <workers> [ids…]`, rows in `tools/mutants.py`, enumerated from the diff `BASE..M3`:
every new or changed production branch — each `Validate` rule and its bound and constant, each `DecodeVerb` refusal,
`CanonicalVerb`'s two encoder settings, each `VerbMessage` field position, separator, byte order and the hash, the
prefix and namespace constants, P1 and P2–P5 in each law, P6, the W1–W14 call, W15 (deletion, boundary, the zero
exemption), every `SignedVerb` field the laws set, each `Binds` comparison, the zero-value and `ok` guards, `ApprovedBy`'s
text, `Use`'s copy, `PasskeyUse`'s format guard, `CheckConfirmation`'s two cases, its approver check, its use of the
caller's inputs, its omission of P1 and of W15 (each **inserted** as a mutant: `CC-SSH-P1-IN`, `CC-W15-IN`), its
default, every line, argument, bound and branch of `FinishPromptText`, `DisplayCommitter`'s `""` rule and text, and the
three cmd delegations): for each row, an APFS clone of `work` at `M3` (`cp -cR`; byte-identical to the round-trip tree,
C12), one exact-string replacement whose anchor occurs exactly once (checked before the run: 0 bad anchors), then
`go test -count=1` of the **named test**, then the **behavioural suites** with every static test `-skip`ped —
`host/approval` rows: `./host/approval/...` minus `TestApprovalImportsOnlyTheStandardLibrary`, `TestApprovalTestIsTestOnly`
(`suite`), and again minus `TestVerbMessageLayout` too (`suite-nolit`: VL3 pins the constants' literals besides the
layout, so a constant mutant's kill is shown not to rest on it); `host/broker`/cmd rows: `./host/broker/ -run
'^(TestFinishPromptText|TestDisplayCommitter|TestRecordOperator|TestParseFinish|TestFormatFinish)'` (`suite`) and
`./cmd/ailang-worldd/` over VT1's `TestFinishPromptIsTheGolden` **plus the 13 existing finish/preset prompt tests**
(`suite-cmd`; the class column says whether the existing tests alone kill it). Go rows ran with `GOFLAGS=-vet=off`,
three workers, at load 13–25 (the final gates beside them). A build failure is `NOCOMPILE` and not a mutant (its
rewrite is the next row). Every run appends one line per row to `plan-scratch/logs/mutant_ledger.tsv` (time, run, id,
file, expected, named result, named failures, suites, named regex, seconds); per-row output in `logs/mut/<run>-<id>.log`.
**The ids below are the ledger's ids**; the table shows each row's last line. `K` killed by the named test **and** by
the suites; `S` killed only by a static/literal pin; `E` a named equivalent (reason given); `CTL` unmutated control.

**Runs.** **R1**: every row (112) on the final stage tree. **R2**: the controls and every row whose R1 result was not
as expected, on the same tree: `V2-UTF8`, `D-CANON`, `CMD-DC` did not compile (an import or variable left unused) and are
rewritten as `V2-UTF8b`, `D-CANONb`, `CMD-DCb` (killed); `S-MSG-ERR` (the SSH message's error ignored) was named
against VL5, where every expiry is the row's own, and **survived its named test** while the suite killed it through VL8's
`ssh a negative expiry` row — R2 names VL8 (`TestCheckConfirmation`), whose row was added for exactly this branch before R1.

**Note on the layout rows** (`M-ZERO`, `M-SWAP`, `M-ENDIAN`, `M-HASH`, `M-HASH-ERR`): each is killed by VL3
(`TestVerbMessageLayout` **runs** `VerbMessage`/`VerbMessageHash` and checks every offset — behavioural, not a
literal pin) and by nothing else, as expected: every other test signs and verifies through the same `VerbMessage`, so a
self-consistent change of the layout is invisible to them (`suite-nolit` = PASS). The cross-implementation pin is
sprint 2's CP4 (the page's message hash equals Go's for a golden written by a Go test). `M-PREFIX`, `M-NS-MINT`,
`M-NS-DEC` and `S-NS` are killed with VL3 excluded (VL4's domain separation), so the namespace/prefix kills are
behavioural.

**Summary:** 115 rows (3 `CTL`, 3 `E`, 105 `K`, 3 `N`, 1 `S`); 118 ledger lines (runs R1, R2); 115 of 115 rows run; the last line of each row behaves as expected for **115 of 115**. Table ids = ledger ids (`tools/mktable.py` fills §4 from the ledger's last line per id); ledger `plan-scratch/logs/mutant_ledger.tsv`

| id | file | from → to | named test → result | behavioural suites (static skipped) | expected / class | as expected |
|---|---|---|---|---|---|---|
| `CTL-A` | `verb.go` | control (no change) | `(TestVerb|TestConfirm|TestSignedVerb|TestCheckConfirmation)` → **PASS** | suite=PASS;suite-nolit=PASS | CTL | yes |
| `V1-SCHEMA` | `verb.go` | `if r.Schema != VerbRequestSchema { // V1` → `if false { // V1` | `TestVerbRequestValidate` → **KILLED** | suite=KILLED[TestVerbRequestValidate];suite-nolit=KILLED[TestVerbRequestValidate] | K | yes |
| `V2-UTF8` | `verb.go` | `if !utf8.ValidString(s) {` → `if false {` | `TestVerbRequestValidate` → **NOCOMPILE** |  | N — does not compile; not a mutant (rewritten as `V2-UTF8b`) | yes |
| `V3-VERB` | `verb.go` | `if r.Verb != VerbFinish { // V3` → `if false { // V3` | `TestVerbRequestValidate` → **KILLED** | suite=KILLED[TestVerbRequestValidate];suite-nolit=KILLED[TestVerbRequestValidate] | K | yes |
| `CTL-B` | `finish_prompt.go` | control (no change) | `(TestFinishPromptText|TestDisplayCommitter)` → **PASS** | suite=PASS;suite-cmd=PASS | CTL | yes |
| `CTL-C` | `episode_finish_operator.go` | control (no change) | `TestFinishPromptIsTheGolden` → **PASS** | suite=PASS;suite-cmd=PASS | CTL | yes |
| `V4-EPISODE` | `verb.go` | `if !episodePattern.MatchString(r.Episode) { // V4` → `if false { // V4` | `TestVerbRequestValidate` → **KILLED** | suite=KILLED[TestVerbRequestValidate];suite-nolit=KILLED[TestVerbRequestValidate] | K | yes |
| `V5-LEN` | `verb.go` | `if len(r.Argv) < 4 \|\| r.Argv[0]` → `if len(r.Argv) < 3 \|\| r.Argv[0]` | `TestVerbRequestValidate` → **KILLED** | suite=KILLED[TestVerbRequestValidate];suite-nolit=KILLED[TestVerbRequestValidate] | K | yes |
| `V5-ARGV0` | `verb.go` | `r.Argv[0] != "ailang-worldd" \|\|` → `false \|\|` | `TestVerbRequestValidate` → **KILLED** | suite=KILLED[TestVerbRequestValidate];suite-nolit=KILLED[TestVerbRequestValidate] | K | yes |
| `V5-ARGV1` | `verb.go` | `r.Argv[1] != "episode" \|\|` → `false \|\|` | `TestVerbRequestValidate` → **KILLED** | suite=KILLED[TestVerbRequestValidate];suite-nolit=KILLED[TestVerbRequestValidate] | K | yes |
| `V5-ARGV2` | `verb.go` | `r.Argv[2] != "finish" \|\|` → `false \|\|` | `TestVerbRequestValidate` → **KILLED** | suite=KILLED[TestVerbRequestValidate];suite-nolit=KILLED[TestVerbRequestValidate] | K | yes |
| `V5-ARGV3` | `verb.go` | `r.Argv[3] != r.Episode { // V5` → `false { // V5` | `TestVerbRequestValidate` → **KILLED** | suite=KILLED[TestVerbRequestValidate];suite-nolit=KILLED[TestVerbRequestValidate] | K | yes |
| `V6-MAX` | `verb.go` | `if len(r.Argv) > verbArgvMax { // V6` → `if false { // V6` | `TestVerbRequestValidate` → **KILLED** | suite=KILLED[TestVerbRequestValidate];suite-nolit=KILLED[TestVerbRequestValidate] | K | yes |
| `V6-CONST` | `verb.go` | `const verbArgvMax = 64` → `const verbArgvMax = 65` | `TestVerbRequestValidate` → **KILLED** | suite=KILLED[TestVerbRequestValidate];suite-nolit=KILLED[TestVerbRequestValidate] | K | yes |
| `V6-BOUND` | `verb.go` | `if len(r.Argv) > verbArgvMax { // V6` → `if len(r.Argv) >= verbArgvMax { // V6` | `TestVerbRequestValidate` → **KILLED** | suite=KILLED[TestVerbRequestValidate];suite-nolit=KILLED[TestVerbRequestValidate] | K | yes |
| `V7-LT` | `verb.go` | `if w[j] < 0x20 \|\| w[j] == 0x7f {` → `if w[j] < 0x1f \|\| w[j] == 0x7f {` | `TestVerbRequestValidate` → **KILLED** | suite=KILLED[TestVerbRequestValidate];suite-nolit=KILLED[TestVerbRequestValidate] | K | yes |
| `V7-LE` | `verb.go` | `if w[j] < 0x20 \|\| w[j] == 0x7f {` → `if w[j] <= 0x20 \|\| w[j] == 0x7f {` | `TestVerbRequestValidate` → **KILLED** | suite=KILLED[TestVerbRequestValidate];suite-nolit=KILLED[TestVerbRequestValidate] | K | yes |
| `V7-CTL` | `verb.go` | `if w[j] < 0x20 \|\| w[j] == 0x7f {` → `if w[j] == 0x7f {` | `TestVerbRequestValidate` → **KILLED** | suite=KILLED[TestVerbRequestValidate];suite-nolit=KILLED[TestVerbRequestValidate] | K | yes |
| `V8-PROMPT` | `verb.go` | `{{"prompt", r.Prompt}, {"scope"` → `{{"scope"` | `TestVerbRequestValidate` → **KILLED** | suite=KILLED[TestVerbRequestValidate];suite-nolit=KILLED[TestVerbRequestValidate] | K | yes |
| `V8-SCOPE` | `verb.go` | `{"scope", r.Scope}, {"title"` → `{"title"` | `TestVerbRequestValidate` → **KILLED** | suite=KILLED[TestVerbRequestValidate];suite-nolit=KILLED[TestVerbRequestValidate] | K | yes |
| `V7-DEL` | `verb.go` | `if w[j] < 0x20 \|\| w[j] == 0x7f {` → `if w[j] < 0x20 {` | `TestVerbRequestValidate` → **KILLED** | suite=KILLED[TestVerbRequestValidate];suite-nolit=KILLED[TestVerbRequestValidate] | K | yes |
| `V9-BODY` | `verb.go` | `if len(r.Body) > verbBodyMax { // V9` → `if false { // V9` | `TestVerbRequestValidate` → **KILLED** | suite=KILLED[TestVerbRequestValidate];suite-nolit=KILLED[TestVerbRequestValidate] | K | yes |
| `V8-TITLE` | `verb.go` | `, {"title", r.Title}} { // V8` → `} { // V8` | `TestVerbRequestValidate` → **KILLED** | suite=KILLED[TestVerbRequestValidate];suite-nolit=KILLED[TestVerbRequestValidate] | K | yes |
| `V9-BOUND` | `verb.go` | `if len(r.Body) > verbBodyMax { // V9` → `if len(r.Body) >= verbBodyMax { // V9` | `TestVerbRequestValidate` → **KILLED** | suite=KILLED[TestVerbRequestValidate];suite-nolit=KILLED[TestVerbRequestValidate] | K | yes |
| `V10-SIZE` | `verb.go` | `if len(c) > MaxVerbRequestBytes { // V10` → `if false { // V10` | `TestVerbRequestValidate` → **KILLED** | suite=KILLED[TestVerbRequestValidate];suite-nolit=KILLED[TestVerbRequestValidate] | K | yes |
| `V9-CONST` | `verb.go` | `const verbBodyMax = 60000` → `const verbBodyMax = 60001` | `TestVerbRequestValidate` → **KILLED** | suite=KILLED[TestVerbRequestValidate];suite-nolit=KILLED[TestVerbRequestValidate] | K | yes |
| `V10-BOUND` | `verb.go` | `if len(c) > MaxVerbRequestBytes { // V10` → `if len(c) >= MaxVerbRequestBytes { // V10` | `TestVerbRequestValidate` → **KILLED** | suite=KILLED[TestVerbRequestValidate];suite-nolit=KILLED[TestVerbRequestValidate] | K | yes |
| `V10-CONST` | `verb.go` | `const MaxVerbRequestBytes = 128 << 10` → `const MaxVerbRequestBytes = 129 << 10` | `TestVerbRequestValidate` → **KILLED** | suite=KILLED[TestVerbMessageLayout,TestVerbRequestValidate];suite-nolit=KILLED[TestVerbRequestValidate] | K | yes |
| `V-CANON-ERR` | `verb.go` | `c, err := CanonicalVerb(r)⏎if err != nil {⏎return err⏎}` → `c, _ := CanonicalVerb(r)` | `TestVerbRequestValidate` → **PASS** | suite=PASS;suite-nolit=PASS | E — `CanonicalVerb` of a `VerbRequest` (strings and string slices only) cannot fail; the check is defensive | yes |
| `C-ESCAPE` | `verb.go` | `enc.SetEscapeHTML(false)` → `enc.SetEscapeHTML(true)` | `(TestVerbRequestRoundTrip|TestVerbRequestDecodeRefusals)` → **KILLED** (TestVerbRequestDecodeRefusals,TestVerbRequestRoundTrip) | suite=KILLED[TestVerbRequestDecodeRefusals,TestVerbRequestRoundTrip];suite-nolit=KILLED[TestVerbRequestDecodeRefusals,T… | K | yes |
| `D-CANON` | `verb.go` | `if !bytes.Equal(c, b) {` → `if false {` | `(TestVerbRequestRoundTrip|TestVerbRequestDecodeRefusals)` → **NOCOMPILE** |  | N — does not compile; not a mutant (rewritten as `D-CANONb`) | yes |
| `C-TRIM` | `verb.go` | `return bytes.TrimSuffix(buf.Bytes(), []byte("\n")), nil` → `return buf.Bytes(), nil` | `(TestVerbRequestRoundTrip|TestVerbRequestDecodeRefusals)` → **KILLED** (TestVerbRequestDecodeRefusals,TestVerbRequestRoundTrip) | suite=KILLED[TestVerbRequestDecodeRefusals,TestVerbRequestRoundTrip];suite-nolit=KILLED[TestVerbRequestDecodeRefusals,T… | K | yes |
| `D-UNKNOWN` | `verb.go` | `dec.DisallowUnknownFields()` → `` | `(TestVerbRequestRoundTrip|TestVerbRequestDecodeRefusals)` → **KILLED** (TestVerbRequestDecodeRefusals) | suite=KILLED[TestVerbRequestDecodeRefusals];suite-nolit=KILLED[TestVerbRequestDecodeRefusals] | K | yes |
| `D-ERR` | `verb.go` | `return VerbRequest{}, fmt.Errorf("%w: %v", ErrDecode, err)` → `return VerbRequest{}, ErrNotCanonical` | `(TestVerbRequestRoundTrip|TestVerbRequestDecodeRefusals)` → **KILLED** (TestVerbRequestDecodeRefusals) | suite=KILLED[TestVerbRequestDecodeRefusals];suite-nolit=KILLED[TestVerbRequestDecodeRefusals] | K | yes |
| `M-PREFIX` | `verb.go` | `const VerbPrefix = "world-verb/v1"` → `const VerbPrefix = "world-approval/v1"` | `TestVerbDomainSeparation` → **KILLED** | suite=KILLED[TestCheckConfirmation,TestChromeMintAssertionIsNotAVerb,TestVerbDomainSeparation,TestVerbMessageLayout];su… | K | yes |
| `M-NS-MINT` | `verb.go` | `const VerbNamespace = "world-verb@sunholo"` → `const VerbNamespace = "world-approval@sunholo"` | `TestVerbDomainSeparation` → **KILLED** | suite=KILLED[TestCheckConfirmation,TestVerbDomainSeparation,TestVerbMessageLayout];suite-nolit=KILLED[TestCheckConfirma… | K | yes |
| `M-NS-DEC` | `verb.go` | `const VerbNamespace = "world-verb@sunholo"` → `const VerbNamespace = "world-decision@sunholo"` | `TestVerbDomainSeparation` → **KILLED** | suite=KILLED[TestCheckConfirmation,TestVerbDomainSeparation,TestVerbMessageLayout];suite-nolit=KILLED[TestCheckConfirma… | K | yes |
| `M-ZERO` | `verb.go` | `m = append(m, 0)⏎m = append(m, digest[:]...)` → `m = append(m, 1)⏎m = append(m, digest[:]...)` | `TestVerbMessageLayout` → **KILLED** | suite=KILLED[TestVerbMessageLayout];suite-nolit=PASS | K (VL3: the layout; no self-consistent sign/verify test can see it — see note) | yes |
| `M-SWAP` | `verb.go` | `m = append(m, digest[:]...)⏎m = append(m, challenge[:]...)` → `m = append(m, challenge[:]...)⏎m = append(m, digest[:]...)` | `TestVerbMessageLayout` → **KILLED** | suite=KILLED[TestVerbMessageLayout];suite-nolit=PASS | K (VL3: the layout; no self-consistent sign/verify test can see it — see note) | yes |
| `M-ENDIAN` | `verb.go` | `m = binary.BigEndian.AppendUint64(m, uint64(expiry))` → `m = binary.LittleEndian.AppendUint64(m, uint64(expiry))` | `TestVerbMessageLayout` → **KILLED** | suite=KILLED[TestVerbMessageLayout];suite-nolit=PASS | K (VL3: the layout; no self-consistent sign/verify test can see it — see note) | yes |
| `M-NEG` | `verb.go` | `if expiry < 0 {⏎return nil, ErrNegativeExpiry` → `if expiry < -1 {⏎return nil, ErrNegativeExpiry` | `TestVerbMessageLayout` → **KILLED** | suite=KILLED[TestCheckConfirmation,TestVerbMessageLayout];suite-nolit=KILLED[TestCheckConfirmation] | K | yes |
| `M-HASH` | `verb.go` | `return sha256.Sum256(m), nil` → `return sha256.Sum256(m[1:]), nil` | `TestVerbMessageLayout` → **KILLED** | suite=KILLED[TestVerbMessageLayout];suite-nolit=PASS | K (VL3: the layout; no self-consistent sign/verify test can see it — see note) | yes |
| `M-HASH-ERR` | `verb.go` | `m, err := VerbMessage(digest, challenge, expiry)⏎if err !…` → `m, _ := VerbMessage(digest, challenge, expiry)` | `TestVerbMessageLayout` → **KILLED** | suite=KILLED[TestVerbMessageLayout];suite-nolit=PASS | K (VL3: the layout; no self-consistent sign/verify test can see it — see note) | yes |
| `M-LEN` | `verb.go` | `const VerbMessageLen = len(VerbPrefix) + 1 + 32 + 32 + 8` → `const VerbMessageLen = len(VerbPrefix) + 1 + 32 + 32 + 9` | `TestVerbMessageLayout` → **KILLED** | suite=KILLED[TestVerbMessageLayout];suite-nolit=PASS | S — `VerbMessageLen` is only `make`'s capacity; behaviour unchanged; killed only by VL3's literal pin (expected, as designed) | yes |
| `S-P1` | `verb_signed.go` | `if v := Admit(in.Admit); v != Admitted { // P1⏎return nil…` → `if err := checkApprover(in.Admit, in.Approver, FormatSSH)` | `TestConfirmBySignatureLaw` → **KILLED** | suite=KILLED[TestConfirmBySignatureLaw];suite-nolit=KILLED[TestConfirmBySignatureLaw] | K | yes |
| `S-APPROVER` | `verb_signed.go` | `if err := checkApprover(in.Admit, in.Approver, FormatSSH)…` → `if err := checkApprover(in.Admit, in.Approver, FormatSSH)…` | `TestConfirmBySignatureLaw` → **KILLED** | suite=KILLED[TestConfirmBySignatureLaw];suite-nolit=KILLED[TestConfirmBySignatureLaw] | K | yes |
| `S-P6` | `verb_signed.go` | `in.Admit.ExpiresAt, in.Approver.Key); err != nil { // P6` → `in.Admit.ExpiresAt, in.Approver.Key); false && err != nil…` | `TestConfirmBySignatureLaw` → **KILLED** | suite=KILLED[TestConfirmBySignatureLaw,TestVerbDomainSeparation];suite-nolit=KILLED[TestConfirmBySignatureLaw,TestVerbD… | K | yes |
| `S-P6-CHAL` | `verb_signed.go` | `verifyVerbSSHSIG(in.Signature, in.Admit.StoredDigest, in.…` → `verifyVerbSSHSIG(in.Signature, in.Admit.StoredDigest, in.…` | `TestConfirmBySignatureLaw` → **KILLED** | suite=KILLED[TestConfirmBySignatureLaw,TestSignedVerbBinds,TestVerbDomainSeparation];suite-nolit=KILLED[TestConfirmBySi… | K | yes |
| `S-P6-DIG` | `verb_signed.go` | `verifyVerbSSHSIG(in.Signature, in.Admit.StoredDigest, in.…` → `verifyVerbSSHSIG(in.Signature, in.Admit.BytesDigest, in.C…` | `TestConfirmBySignatureLaw` → **PASS** | suite=PASS;suite-nolit=PASS | E — after P1 admits, `BytesDigest == StoredDigest`, so verifying over either is the same message | yes |
| `S-NS` | `verb_signed.go` | `return VerifySSHSIG(s, VerbNamespace, msg, key)` → `return VerifySSHSIG(s, Namespace, msg, key)` | `TestVerbDomainSeparation` → **KILLED** | suite=KILLED[TestCheckConfirmation,TestConfirmBySignatureLaw,TestSignedVerbBinds,TestVerbDomainSeparation];suite-nolit=… | K | yes |
| `S-PARSE` | `verb_signed.go` | `s, err := ParseSSHSIG(armored)⏎if err != nil {⏎return err…` → `s, _ := ParseSSHSIG(armored)⏎if s == nil {⏎return ErrBadS…` | `TestConfirmBySignatureLaw` → **KILLED** | suite=KILLED[TestCheckConfirmation,TestConfirmBySignatureLaw];suite-nolit=KILLED[TestCheckConfirmation,TestConfirmBySig… | K | yes |
| `S-MSG-ERR` | `verb_signed.go` | `msg, err := VerbMessage(digest, challenge, expiry)⏎if err…` → `msg, _ := VerbMessage(digest, challenge, expiry)` | `TestCheckConfirmation` → **KILLED** | suite=KILLED[TestCheckConfirmation];suite-nolit=KILLED[TestCheckConfirmation] | K | yes |
| `S-FORMAT` | `verb_signed.go` | `name: in.Approver.Name, format: FormatSSH, key: bytes.Clo…` → `name: in.Approver.Name, format: FormatWebAuthn, key: byte…` | `TestConfirmBySignatureLaw` → **KILLED** | suite=KILLED[TestConfirmBySignatureLaw,TestVerbDomainSeparation];suite-nolit=KILLED[TestConfirmBySignatureLaw,TestVerbD… | K | yes |
| `S-KEY` | `verb_signed.go` | `name: in.Approver.Name, format: FormatSSH, key: bytes.Clo…` → `name: in.Approver.Name, format: FormatSSH, key: nil}, nil…` | `TestConfirmBySignatureLaw` → **KILLED** | suite=KILLED[TestConfirmBySignatureLaw,TestSignedVerbBinds];suite-nolit=KILLED[TestConfirmBySignatureLaw,TestSignedVerb… | K | yes |
| `S-BIND` | `verb_signed.go` | `return &SignedVerb{ok: true, digest: in.Admit.StoredDiges…` → `return &SignedVerb{ok: true, digest: in.Admit.StoredDiges…` | `TestConfirmBySignatureLaw` → **KILLED** | suite=KILLED[TestConfirmBySignatureLaw,TestSignedVerbBinds];suite-nolit=KILLED[TestConfirmBySignatureLaw,TestSignedVerb… | K | yes |
| `A-P1` | `verb_signed.go` | `if v := Admit(in.Admit); v != Admitted { // P1⏎return nil…` → `if err := checkApprover(in.Admit, in.Approver, FormatWebA…` | `TestConfirmByAssertionLaw` → **KILLED** | suite=KILLED[TestConfirmByAssertionLaw];suite-nolit=KILLED[TestConfirmByAssertionLaw] | K | yes |
| `A-APPROVER` | `verb_signed.go` | `if err := checkApprover(in.Admit, in.Approver, FormatWebA…` → `if err := checkApprover(in.Admit, in.Approver, FormatWebA…` | `TestConfirmByAssertionLaw` → **KILLED** | suite=KILLED[TestConfirmByAssertionLaw];suite-nolit=KILLED[TestConfirmByAssertionLaw] | K | yes |
| `A-HASH-MINT` | `verb_signed.go` | `h, hErr := VerbMessageHash(in.Admit.StoredDigest, in.Chal…` → `h, hErr := MessageHash(in.Admit.StoredDigest, in.Challeng…` | `TestConfirmByAssertionLaw` → **KILLED** | suite=KILLED[TestChromeMintAssertionIsNotAVerb,TestConfirmByAssertionLaw,TestVerbDomainSeparation];suite-nolit=KILLED[T… | K | yes |
| `A-STEPS` | `verb_signed.go` | `ad, err := assertionSteps(in.Response, in.Approver, h, hE…` → `ad, _ := assertionSteps(in.Response, in.Approver, h, hErr…` | `TestConfirmByAssertionLaw` → **KILLED** | suite=KILLED[TestChromeMintAssertionIsNotAVerb,TestConfirmByAssertionLaw,TestVerbDomainSeparation];suite-nolit=KILLED[T… | K | yes |
| `A-W15-DEL` | `verb_signed.go` | `if (ad.SignCount != 0 \|\| p.SignCount != 0) && ad.SignCo…` → `if false { // W15` | `TestConfirmByAssertionLaw` → **KILLED** | suite=KILLED[TestConfirmByAssertionLaw];suite-nolit=KILLED[TestConfirmByAssertionLaw] | K | yes |
| `A-W15-LT` | `verb_signed.go` | `if (ad.SignCount != 0 \|\| p.SignCount != 0) && ad.SignCo…` → `if (ad.SignCount != 0 \|\| p.SignCount != 0) && ad.SignCo…` | `TestConfirmByAssertionLaw` → **KILLED** | suite=KILLED[TestConfirmByAssertionLaw];suite-nolit=KILLED[TestConfirmByAssertionLaw] | K | yes |
| `A-W15-ZERO` | `verb_signed.go` | `if (ad.SignCount != 0 \|\| p.SignCount != 0) && ad.SignCo…` → `if ad.SignCount <= p.SignCount { // W15` | `TestConfirmByAssertionLaw` → **KILLED** | suite=KILLED[TestConfirmByAssertionLaw];suite-nolit=KILLED[TestConfirmByAssertionLaw] | K | yes |
| `A-FORMAT` | `verb_signed.go` | `name: in.Approver.Name, format: FormatWebAuthn, key: p.Ke…` → `name: in.Approver.Name, format: FormatSSH, key: p.Key.Can…` | `TestConfirmByAssertionLaw` → **KILLED** | suite=KILLED[TestConfirmByAssertionLaw,TestVerbDomainSeparation];suite-nolit=KILLED[TestConfirmByAssertionLaw,TestVerbD… | K | yes |
| `A-CRED` | `verb_signed.go` | `credentialID: bytes.Clone(p.CredentialID), signCount: ad.…` → `credentialID: nil, signCount: ad.SignCount}, nil⏎}⏎⏎// Er…` | `TestConfirmByAssertionLaw` → **KILLED** | suite=KILLED[TestConfirmByAssertionLaw];suite-nolit=KILLED[TestConfirmByAssertionLaw] | K | yes |
| `A-COUNT` | `verb_signed.go` | `credentialID: bytes.Clone(p.CredentialID), signCount: ad.…` → `credentialID: bytes.Clone(p.CredentialID), signCount: p.S…` | `TestConfirmByAssertionLaw` → **KILLED** | suite=KILLED[TestConfirmByAssertionLaw];suite-nolit=KILLED[TestConfirmByAssertionLaw] | K | yes |
| `CC-SSH-CASE` | `verb_signed.go` | `switch in.Format {⏎case FormatSSH:` → `switch in.Format {⏎case "ssh-x":` | `TestCheckConfirmation` → **KILLED** | suite=KILLED[TestCheckConfirmation];suite-nolit=KILLED[TestCheckConfirmation] | K | yes |
| `CC-SSH-APPROVER` | `verb_signed.go` | `if err := checkApprover(AdmitInput{Now: in.Now}, in.Appro…` → `if err := checkApprover(AdmitInput{Now: in.Now}, in.Appro…` | `TestCheckConfirmation` → **KILLED** | suite=KILLED[TestCheckConfirmation];suite-nolit=KILLED[TestCheckConfirmation] | K | yes |
| `CC-SSH-P6` | `verb_signed.go` | `in.Expiry, in.Approver.Key); err != nil { // P6` → `in.Expiry, in.Approver.Key); false && err != nil { // P6` | `TestCheckConfirmation` → **KILLED** | suite=KILLED[TestCheckConfirmation];suite-nolit=KILLED[TestCheckConfirmation] | K | yes |
| `CC-SSH-OWN` | `verb_signed.go` | `verifyVerbSSHSIG(in.Artifact, in.Digest, in.Challenge, in…` → `verifyVerbSSHSIG(in.Artifact, in.Challenge, in.Digest, in…` | `TestCheckConfirmation` → **KILLED** | suite=KILLED[TestCheckConfirmation];suite-nolit=KILLED[TestCheckConfirmation] | K | yes |
| `CC-SSH-P1-IN` | `verb_signed.go` | `if err := checkApprover(AdmitInput{Now: in.Now}, in.Appro…` → `if in.Now >= in.Expiry {⏎return nil, &NotAdmittedError{Ve…` | `TestCheckConfirmation` → **KILLED** | suite=KILLED[TestCheckConfirmation];suite-nolit=KILLED[TestCheckConfirmation] | K | yes |
| `CC-WA-CASE` | `verb_signed.go` | `case FormatWebAuthn:⏎if err := checkApprover(AdmitInput{N…` → `case "webauthn-x":⏎if err := checkApprover(AdmitInput{Now…` | `TestCheckConfirmation` → **KILLED** | suite=KILLED[TestCheckConfirmation,TestChromeMintAssertionIsNotAVerb];suite-nolit=KILLED[TestCheckConfirmation,TestChro… | K | yes |
| `CC-WA-APPROVER` | `verb_signed.go` | `if err := checkApprover(AdmitInput{Now: in.Now}, in.Appro…` → `if err := checkApprover(AdmitInput{Now: in.Now}, in.Appro…` | `TestCheckConfirmation` → **KILLED** | suite=KILLED[TestCheckConfirmation];suite-nolit=KILLED[TestCheckConfirmation] | K | yes |
| `CC-WA-HASH` | `verb_signed.go` | `h, hErr := VerbMessageHash(in.Digest, in.Challenge, in.Ex…` → `h, hErr := MessageHash(in.Digest, in.Challenge, in.Expiry)` | `TestCheckConfirmation` → **KILLED** | suite=KILLED[TestCheckConfirmation,TestChromeMintAssertionIsNotAVerb];suite-nolit=KILLED[TestCheckConfirmation,TestChro… | K | yes |
| `CC-WA-OWN` | `verb_signed.go` | `h, hErr := VerbMessageHash(in.Digest, in.Challenge, in.Ex…` → `h, hErr := VerbMessageHash(in.Digest, in.Digest, in.Expir…` | `TestCheckConfirmation` → **KILLED** | suite=KILLED[TestCheckConfirmation];suite-nolit=KILLED[TestCheckConfirmation] | K | yes |
| `CC-WA-STEPS` | `verb_signed.go` | `ad, err := assertionSteps(in.Artifact, in.Approver, h, hE…` → `ad, _ := assertionSteps(in.Artifact, in.Approver, h, hErr…` | `TestCheckConfirmation` → **KILLED** | suite=KILLED[TestCheckConfirmation,TestChromeMintAssertionIsNotAVerb];suite-nolit=KILLED[TestCheckConfirmation,TestChro… | K | yes |
| `CC-W15-IN` | `verb_signed.go` | `p := in.Approver.Passkey⏎return &SignedVerb{ok: true, dig…` → `p := in.Approver.Passkey⏎if ad.SignCount <= p.SignCount {…` | `TestCheckConfirmation` → **KILLED** | suite=KILLED[TestCheckConfirmation];suite-nolit=KILLED[TestCheckConfirmation] | K | yes |
| `CC-DEFAULT` | `verb_signed.go` | `return nil, fmt.Errorf("%w: %q", ErrVerbFormat, in.Format)` → `return &SignedVerb{ok: true, format: in.Format, name: in.…` | `TestCheckConfirmation` → **KILLED** | suite=KILLED[TestCheckConfirmation];suite-nolit=KILLED[TestCheckConfirmation] | K | yes |
| `CC-BIND` | `verb_signed.go` | `return &SignedVerb{ok: true, digest: in.Digest, challenge…` → `return &SignedVerb{ok: true, digest: in.Digest, challenge…` | `TestCheckConfirmation` → **KILLED** | suite=KILLED[TestCheckConfirmation];suite-nolit=KILLED[TestCheckConfirmation] | K | yes |
| `B-OK` | `verb_signed.go` | `case s == nil \|\| !s.ok:⏎return ErrNotVerified` → `case s == nil:⏎return ErrNotVerified` | `TestSignedVerbBinds` → **KILLED** | suite=KILLED[TestSignedVerbBinds];suite-nolit=KILLED[TestSignedVerbBinds] | K | yes |
| `B-DIGEST` | `verb_signed.go` | `case s.digest != digest:` → `case false:` | `TestSignedVerbBinds` → **KILLED** | suite=KILLED[TestSignedVerbBinds];suite-nolit=KILLED[TestSignedVerbBinds] | K | yes |
| `B-CHALLENGE` | `verb_signed.go` | `case s.challenge != challenge:` → `case false:` | `TestSignedVerbBinds` → **KILLED** | suite=KILLED[TestSignedVerbBinds];suite-nolit=KILLED[TestSignedVerbBinds] | K | yes |
| `B-EXPIRY` | `verb_signed.go` | `case s.expiry != expiry:` → `case false:` | `TestSignedVerbBinds` → **KILLED** | suite=KILLED[TestSignedVerbBinds];suite-nolit=KILLED[TestSignedVerbBinds] | K | yes |
| `AB-GUARD` | `verb_signed.go` | `if s == nil \|\| !s.ok {⏎return ""⏎}⏎return s.format + ":…` → `if s == nil {⏎return ""⏎}⏎return s.format + ":" + s.name` | `TestSignedVerbBinds` → **KILLED** | suite=KILLED[TestSignedVerbBinds];suite-nolit=KILLED[TestSignedVerbBinds] | K | yes |
| `AB-TEXT` | `verb_signed.go` | `return s.format + ":" + s.name` → `return s.name` | `TestConfirmBySignatureLaw` → **KILLED** | suite=KILLED[TestCheckConfirmation,TestConfirmByAssertionLaw,TestConfirmBySignatureLaw,TestVerbDomainSeparation];suite-… | K | yes |
| `U-GUARD` | `verb_signed.go` | `if s == nil \|\| !s.ok {⏎return "", "", nil, false⏎}` → `if s == nil {⏎return "", "", nil, false⏎}` | `TestSignedVerbBinds` → **KILLED** | suite=KILLED[TestSignedVerbBinds];suite-nolit=KILLED[TestSignedVerbBinds] | K | yes |
| `U-CLONE` | `verb_signed.go` | `return s.name, s.format, bytes.Clone(s.key), true` → `return s.name, s.format, s.key, true` | `TestSignedVerbBinds` → **KILLED** | suite=KILLED[TestSignedVerbBinds];suite-nolit=KILLED[TestSignedVerbBinds] | K | yes |
| `PU-FORMAT` | `verb_signed.go` | `if s == nil \|\| !s.ok \|\| s.format != FormatWebAuthn {` → `if s == nil \|\| !s.ok {` | `TestConfirmBySignatureLaw` → **KILLED** | suite=KILLED[TestConfirmBySignatureLaw];suite-nolit=KILLED[TestConfirmBySignatureLaw] | K | yes |
| `PU-OK` | `verb_signed.go` | `if s == nil \|\| !s.ok \|\| s.format != FormatWebAuthn {` → `if s == nil \|\| s.format != FormatWebAuthn {` | `TestSignedVerbBinds` → **PASS** | suite=PASS;suite-nolit=PASS | E — the zero value's `format` is `""`, so `s.format != FormatWebAuthn` refuses it anyway (A4's `PasskeyUse` has the same redundancy) | yes |
| `FP-REPO` | `finish_prompt.go` | `"repository: %s\nbranch` → `"repo: %s\nbranch` | `TestFinishPromptText` → **KILLED** | suite=KILLED[TestFinishPromptText];suite-cmd=KILLED[TestFinishPromptIsTheGolden] | K; existing cmd tests alone: **survive** (only VT1 sees it) | yes |
| `FP-COMMITTER` | `finish_prompt.go` | `"committer: %s\n"` → `"committer %s\n"` | `TestFinishPromptText` → **KILLED** | suite=KILLED[TestFinishPromptText];suite-cmd=KILLED[TestEpisodeFinishProtected,TestFinishPromptIsTheGolden] | K; existing cmd tests alone: killed (TestEpisodeFinishProtected) | yes |
| `FP-BRANCH` | `finish_prompt.go` | `\nbranch: %s\n` → `\nhead: %s\n` | `TestFinishPromptText` → **KILLED** | suite=KILLED[TestFinishPromptText];suite-cmd=KILLED[TestFinishPromptIsTheGolden] | K; existing cmd tests alone: **survive** (only VT1 sees it) | yes |
| `FP-BASE` | `finish_prompt.go` | `\nbase: %s\n` → `\nbase %s\n` | `TestFinishPromptText` → **KILLED** | suite=KILLED[TestFinishPromptText];suite-cmd=KILLED[TestFinishPromptIsTheGolden] | K; existing cmd tests alone: **survive** (only VT1 sees it) | yes |
| `FP-SHA` | `finish_prompt.go` | `\nsha: %s\n` → `\ncommit: %s\n` | `TestFinishPromptText` → **KILLED** | suite=KILLED[TestFinishPromptText];suite-cmd=KILLED[TestFinishPromptIsTheGolden] | K; existing cmd tests alone: **survive** (only VT1 sees it) | yes |
| `FP-TITLE` | `finish_prompt.go` | `\ntitle: %s\n` → `\ntitle %s\n` | `TestFinishPromptText` → **KILLED** | suite=KILLED[TestFinishPromptText];suite-cmd=KILLED[TestFinishPromptIsTheGolden] | K; existing cmd tests alone: **survive** (only VT1 sees it) | yes |
| `FP-BODY` | `finish_prompt.go` | `body (%d bytes), first lines:` → `body (%d bytes):` | `TestFinishPromptText` → **KILLED** | suite=KILLED[TestFinishPromptText];suite-cmd=KILLED[TestFinishPromptIsTheGolden] | K; existing cmd tests alone: **survive** (only VT1 sees it) | yes |
| `FP-PROT` | `finish_prompt.go` | `protected paths changed: %s\n"` → `protected paths: %s\n"` | `TestFinishPromptText` → **KILLED** | suite=KILLED[TestFinishPromptText];suite-cmd=KILLED[TestEpisodeFinishProtected,TestFinishPromptIsTheGolden] | K; existing cmd tests alone: killed (TestEpisodeFinishProtected) | yes |
| `FP-LABEL` | `finish_prompt.go` | `\nlabel: %s\n` → `\nlabel %s\n` | `TestFinishPromptText` → **KILLED** | suite=KILLED[TestFinishPromptText];suite-cmd=KILLED[TestFinishPromptIsTheGolden] | K; existing cmd tests alone: **survive** (only VT1 sees it) | yes |
| `FP-SIX` | `finish_prompt.go` | `lines[:min(6, len(lines))]` → `lines[:min(7, len(lines))]` | `TestFinishPromptText` → **KILLED** | suite=KILLED[TestFinishPromptText];suite-cmd=KILLED[TestFinishPromptIsTheGolden] | K; existing cmd tests alone: **survive** (only VT1 sees it) | yes |
| `FP-ARGS` | `finish_prompt.go` | `s.Repository, s.Head, s.SHA, s.Base,` → `s.Repository, s.Base, s.SHA, s.Head,` | `TestFinishPromptText` → **KILLED** | suite=KILLED[TestFinishPromptText];suite-cmd=KILLED[TestFinishPromptIsTheGolden] | K; existing cmd tests alone: **survive** (only VT1 sees it) | yes |
| `FP-END-NL` | `finish_prompt.go` | `protected paths changed: %s\n"` → `protected paths changed: %s"` | `TestFinishPromptText` → **KILLED** | suite=KILLED[TestFinishPromptText];suite-cmd=KILLED[TestFinishPromptIsTheGolden] | K; existing cmd tests alone: **survive** (only VT1 sees it) | yes |
| `FP-STATE` | `finish_prompt.go` | `if !s.Draft {` → `if s.Draft {` | `TestFinishPromptText` → **KILLED** | suite=KILLED[TestFinishPromptText];suite-cmd=KILLED[TestFinishPromptIsTheGolden] | K; existing cmd tests alone: **survive** (only VT1 sees it) | yes |
| `FP-READY` | `finish_prompt.go` | `state = "ready"` → `state = "draft"` | `TestFinishPromptText` → **KILLED** | suite=KILLED[TestFinishPromptText];suite-cmd=KILLED[TestFinishPromptIsTheGolden] | K; existing cmd tests alone: **survive** (only VT1 sees it) | yes |
| `FP-LEN` | `finish_prompt.go` | `title, len(body), preview` → `title, len(preview), preview` | `TestFinishPromptText` → **KILLED** | suite=KILLED[TestFinishPromptText];suite-cmd=KILLED[TestFinishPromptIsTheGolden] | K; existing cmd tests alone: **survive** (only VT1 sees it) | yes |
| `FP-CREATE` | `finish_prompt.go` | `if s.CreateLabel {` → `if !s.CreateLabel {` | `TestFinishPromptText` → **KILLED** | suite=KILLED[TestFinishPromptText];suite-cmd=KILLED[TestFinishPromptIsTheGolden] | K; existing cmd tests alone: **survive** (only VT1 sees it) | yes |
| `FP-DC` | `finish_prompt.go` | `DisplayCommitter(committer)) +` → `committer) +` | `TestFinishPromptText` → **KILLED** | suite=KILLED[TestFinishPromptText];suite-cmd=KILLED[TestFinishPromptIsTheGolden] | K; existing cmd tests alone: **survive** (only VT1 sees it) | yes |
| `FP-JOIN` | `finish_prompt.go` | `strings.Join(protected, ", ")` → `strings.Join(protected, ",")` | `TestFinishPromptText` → **KILLED** | suite=KILLED[TestFinishPromptText];suite-cmd=KILLED[TestFinishPromptIsTheGolden] | K; existing cmd tests alone: **survive** (only VT1 sees it) | yes |
| `DC-TEXT` | `finish_prompt.go` | `return "(repository git config)"` → `return "(git config)"` | `TestDisplayCommitter` → **KILLED** | suite=KILLED[TestDisplayCommitter,TestFinishPromptText];suite-cmd=KILLED[TestFinishPromptIsTheGolden] | K; existing cmd tests alone: **survive** (only VT1 sees it) | yes |
| `DC-RULE` | `finish_prompt.go` | `if v == "" {⏎return "(repository git config)"` → `if false {⏎return "(repository git config)"` | `TestDisplayCommitter` → **KILLED** | suite=KILLED[TestDisplayCommitter,TestFinishPromptText];suite-cmd=KILLED[TestFinishPromptIsTheGolden] | K; existing cmd tests alone: **survive** (only VT1 sees it) | yes |
| `CMD-Q` | `episode_finish_operator.go` | `+"Open this pull request? [y/N] ")` → `+"Open this pull request? [y/n] ")` | `TestFinishPromptIsTheGolden` → **KILLED** | suite=PASS;suite-cmd=KILLED[TestFinishPromptIsTheGolden] | K; existing cmd tests alone: **survive** (only VT1 sees it) | yes |
| `CMD-DC` | `committer.go` | `{ return broker.DisplayCommitter(value) }` → `{ return value }` | `TestFinishPromptIsTheGolden` → **NOCOMPILE** |  | N — does not compile; not a mutant (rewritten as `CMD-DCb`) | yes |
| `CMD-IDENTITY` | `episode_finish_operator.go` | `identity = committer[0]` → `identity = ""` | `TestFinishPromptIsTheGolden` → **KILLED** | suite=PASS;suite-cmd=KILLED[TestEpisodeFinishProtected,TestFinishPromptIsTheGolden] | K; existing cmd tests alone: killed (TestEpisodeFinishProtected) | yes |
| `CMD-TEXT` | `episode_finish_operator.go` | `broker.FinishPromptText(scope, title, body, protected, id…` → `broker.FinishPromptText(scope, title, title, protected, i…` | `TestFinishPromptIsTheGolden` → **KILLED** | suite=PASS;suite-cmd=KILLED[TestFinishPromptIsTheGolden] | K; existing cmd tests alone: **survive** (only VT1 sees it) | yes |
| `V2-UTF8b` | `verb.go` | `if !utf8.ValidString(s) {` → `if !utf8.ValidString(s) && false {` | `TestVerbRequestValidate` → **KILLED** | suite=KILLED[TestVerbRequestValidate];suite-nolit=KILLED[TestVerbRequestValidate] | K | yes |
| `D-CANONb` | `verb.go` | `if !bytes.Equal(c, b) {` → `if !bytes.Equal(c, b) && false {` | `(TestVerbRequestRoundTrip|TestVerbRequestDecodeRefusals)` → **KILLED** (TestVerbRequestDecodeRefusals) | suite=KILLED[TestVerbRequestDecodeRefusals];suite-nolit=KILLED[TestVerbRequestDecodeRefusals] | K | yes |
| `CMD-DCb` | `committer.go` | `{ return broker.DisplayCommitter(value) }` → `{ _ = broker.DisplayCommitter; return value }` | `TestFinishPromptIsTheGolden` → **KILLED** | suite=PASS;suite-cmd=KILLED[TestFinishPromptIsTheGolden] | K; existing cmd tests alone: **survive** (only VT1 sees it) | yes |

---

## 5. Files (final `git status --short --untracked-files=all` after M3, against the base)

```
 M cmd/ailang-worldd/committer.go
 M cmd/ailang-worldd/episode_finish_operator.go
?? cmd/ailang-worldd/finish_prompt_golden_test.go
?? host/approval/verb.go
?? host/approval/verb_chrome_test.go
?? host/approval/verb_check_test.go
?? host/approval/verb_signed.go
?? host/approval/verb_signed_test.go
?? host/approval/verb_test.go
?? host/broker/finish_prompt.go
?? host/broker/finish_prompt_test.go
?? host/broker/testdata/finish-prompt.golden.json
```
Plus the two plan files, which the controller commits. After each sub-milestone, exactly the files its §3 paragraph
names (cumulatively). Nothing under `world/`, `packages/`, `scripts/`, `tools/`, `docs/`, `host/authority/`,
`host/store/`, `host/daemon/`, `host/surfacepage/`, `host/verbs/`, nor `host/broker/approve.go`,
`cmd/ailang-worldd/episode_finish.go`, any `*_test.go` that exists at the base.

---

## 6. Sprint gates (controller, after the last sub-milestone's commit; each measured at the base, §1.1)

Run it in **one bash shell, in the foreground, never under `nohup`** (macOS `/bin/bash` 3.2; every leg records its own rc
and every leg runs — no `set -e`; every leg is `tee`'d to `gates/<leg>.txt`), extracting the script from this plan. The
executor sets `GATES_SKIP_VERIFY_GO=1`; the controller runs it without, after the last commit (`verify_go` reads
committed blobs):

```bash
bash <<'EOF2'
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter400-row243-v3s1
G="$HOME/.ailang/state/world-iter400/gates"; mkdir -p "$G"
awk '$0 == "````bash gates.sh" {f=1; next} f && $0 == "````" {exit} f' "$W/design_docs/planned/w-world-address-bar-verbs-v3s1-sprint-plan.md" > "$G/gates.sh"
/bin/bash "$G/gates.sh" "$W" "$G" final "$(git -C "$W" rev-parse HEAD)"
EOF2
```

It writes `gates/legs.txt`, one `<leg>.txt` per leg, and **`gates/ALL_DONE.<sha>`** (`GATES_DONE <s> FAILS: …`) when every
leg has finished (the judge waits on it). The legs: `status`, `verify_ail` (no `.ail` changed: the base's result),
`gofmt`, `build`, `vet` (`go vet ./...`), `fence` (`go test -run '^$' ./...`, the test-file compile fence), `pkgs`
(`host/approval/...`, `host/authority`, `host/broker`, `cmd/ailang-worldd` whole), `race` (`-race` on `host/approval/...`
and `host/authority` whole, and on the touched tests of `host/broker` and `cmd`), **`prompt`** (the CLI's finish prompt:
the 13 existing `episode_finish_m2/m3` and `preset_committer` tests **by name, unedited**, plus VT1's
`TestFinishPromptIsTheGolden`), `inv` (every named inventory `-v` + every new test: none missing, none SKIP — incl.
`TestApprovalImportsOnlyTheStandardLibrary`, `TestApprovalTestIsTestOnly`, `TestDaemonDependencyAllowlist`,
`TestFrozenV1RouteTableMatchesMux`, `TestDaemonCredentialWriterSites` (unchanged: three sites), `TestApprovalRoutesRegistered`,
`TestDecisionRoutesRegistered`, `TestWorldBoundaryDependencyAllowlist`, `TestProductionContextRoots`, `TestProductionGoSurface`,
`TestStoreSurfaceTakesContext`, `TestNoParallelWriteForkPackagesOutsideVerifygate`, `TestBrokerTestsDoNotCallParallel`,
`TestDecideApprovalCallerCensus`, `TestEveryVerbHasHelp`, `TestVerbTableMatchesTheCLISwitch`, and A4's
`TestMintSignatureIsNotADecision`/`TestDecideByAssertionDefersTheCounter`, the judges' `…IsTheCLIJudge`/`TestJudgeDecision*`),
`verify_go` (`VERIFY_GO_LEGS=plain AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_go.sh`). No live Chrome and
no real-stack leg: S1 has no page and no daemon path. No production context root is added (S1 adds no
`context.Background`/`TODO`), so `TestProductionContextRoots` is an inventory, not a changed gate. The script, verbatim:

````bash gates.sh
#!/bin/bash
# Row 243 V3 S1 gates (plan §6). usage: gates.sh <tree> <logdir> <base|final> [<sha>]
# The controller: bash gates.sh <worktree> $HOME/.ailang/state/world-iter400/gates final <HEAD sha>
# macOS /bin/bash 3.2: no declare -A. Every leg records its own rc; every leg runs (no set -e); every leg is tee'd
# to <logdir>/<leg>.txt. GATES_SKIP_VERIFY_GO=1 skips verify_go (the executor sets it; the controller runs it).
# Never start this under nohup. Writes <logdir>/ALL_DONE.<sha> with FAILS: when every leg has finished.
T=$1; L=$2; MODE=$3
mkdir -p "$L"; cd "$T" || exit 2
SHA=${4:-$(git rev-parse HEAD 2>/dev/null || echo unknown)}
export AILANG_BIN=$HOME/.pinned-ailang/ailang
export WORLD_EXEC_NODE=/opt/homebrew/bin/node
unset WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER WORLD_PASSKEY_RECORD_DIR WORLD_UPDATE_VERB_GOLDEN WORLD_UPDATE_WHY_GOLDEN WORLD_TOOL_AILANG_BIN WORLD_UPDATE_FINISH_PROMPT_GOLDEN
: > "$L/legs.txt"
FAILS=""
s0=$(date +%s)
leg() { # leg <name> <command string>: tee'd to $L/<name>.txt, rc recorded in legs.txt
  local name=$1 cmd=$2 s rc
  s=$(date +%s)
  bash -c "$cmd" 2>&1 | tee "$L/$name.txt"
  rc=${PIPESTATUS[0]}
  echo "LEG $name rc=$rc $(( $(date +%s) - s ))s" | tee -a "$L/legs.txt"
  [ "$rc" -eq 0 ] || FAILS="$FAILS $name"
}
need() { # need <leg> <test id>...: each must show --- PASS in <leg>.txt
  local lg=$1 t; shift
  for t in "$@"; do grep -qE -- "^ *--- PASS: $t " "$L/$lg.txt" || { echo "NOT PASSED ($lg): $t" | tee -a "$L/legs.txt"; FAILS="$FAILS ${lg}_missing"; }; done
}
alt() { local x="$*"; x="${x// /|}"; echo "${x:+|$x}"; }
bar() { echo "$*" | tr ' ' '|'; }
leg status 'git status --porcelain | wc -l'
leg verify_ail './scripts/verify_ail.sh'
leg gofmt 'out=$(gofmt -l cmd host); echo "$out"; test -z "$out"'
leg build 'go build ./...'
leg vet 'go vet ./...'
leg fence "go test -run '^\$' ./... > '$L/fence.raw' 2>&1; rc=\$?; cat '$L/fence.raw'; if grep -E '^(FAIL|--- FAIL)' '$L/fence.raw' >/dev/null; then exit 1; fi; exit \$rc"
leg pkgs "go test -count=1 -timeout 30m ./host/approval/... ./host/authority/ ./host/broker/ ./cmd/ailang-worldd/"
leg race "go test -race -count=1 ./host/approval/... ./host/authority/ && go test -race -count=1 ./host/broker/ -run '^(TestFinish|TestDisplayCommitter|TestRecordOperator|TestDecideApproval)' && go test -race -count=1 ./cmd/ailang-worldd/ -run '^(TestEpisodeFinish|TestFinish|TestPresetCommitter)'"
# the CLI's finish prompt: the existing prompt tests, unedited, by name
PROMPT="TestEpisodeFinishDeclineAC7 TestEpisodeFinishBytesAC8 TestEpisodeFinishDispatchAC9AC11AC12 TestEpisodeFinishOutcomesAC10 TestEpisodeFinishProtected TestEpisodeFinishCredentialsMismatch TestEpisodeFinishResumeAC13 TestEpisodeFinishRerunAC10 TestEpisodeFinishApprovalReuse TestEpisodeFinishFixupRound TestEpisodeFinishFinishedRerun TestEpisodeFinishResumeGuards TestPresetCommitterValidationAndPrompt"
NEW_PROMPT="TestFinishPromptIsTheGolden"; [ "$MODE" = base ] && NEW_PROMPT=""
leg prompt "set -o pipefail; go test -v -count=1 ./cmd/ailang-worldd/ -run '^($(bar $PROMPT)$(alt $NEW_PROMPT))\$'"
need prompt $PROMPT $NEW_PROMPT
if grep -E '^\s*--- (SKIP|FAIL)' "$L/prompt.txt"; then FAILS="$FAILS prompt_skip_or_fail"; fi
NEW_A="TestVerbRequestRoundTrip TestVerbRequestDecodeRefusals TestVerbRequestValidate TestVerbMessageLayout TestVerbDomainSeparation TestConfirmBySignatureLaw TestConfirmByAssertionLaw TestSignedVerbBinds TestCheckConfirmation TestChromeMintAssertionIsNotAVerb"
NEW_B="TestFinishPromptText TestDisplayCommitter"
if [ "$MODE" = base ]; then NEW_A=""; NEW_B=""; fi
INV_APPROVAL="TestApprovalImportsOnlyTheStandardLibrary TestApprovalTestIsTestOnly TestMintSignatureIsNotADecision TestDecideByAssertionDefersTheCounter"
INV_DAEMON="TestDaemonDependencyAllowlist TestFrozenV1RouteTableMatchesMux TestDaemonCredentialWriterSites TestApprovalRoutesRegistered TestDecisionRoutesRegistered TestSurfaceVerbRouteRegistered TestSurfaceVerbIsReadOnly TestSurfaceVerbWhyIsReadOnly"
INV_BOUNDARY="TestWorldBoundaryDependencyAllowlist TestBareNetHTTPExemptionIsPerGroup TestEveryCommandDirectoryIsAProtectedGroup TestBoundaryASTWriteGuard TestWorkbenchPackageRemainsTransportFree TestPinfetchStaysOutOfStoreReplayAndPublish"
INV_STORE="TestProductionContextRoots TestProductionGoSurface TestStoreSurfaceTakesContext TestNoNewDeadlineFreeStoreReads TestReadGettersHonorContext"
INV_VG="TestNoParallelWriteForkPackagesOutsideVerifygate TestSubprocessSinkGate TestVerifygateTestWritesAreForkLocked"
INV_BROKER="TestBrokerTestsDoNotCallParallel TestDecideApprovalCallerCensus TestRegistryDispatchBindingBoundary"
INV_CMD="TestEveryVerbHasHelp TestVerbTableMatchesTheCLISwitch"
leg inv "set -o pipefail
go test -v -count=1 ./host/approval/... -run '^($(bar $INV_APPROVAL)$(alt $NEW_A))\$' &&
go test -v -count=1 ./host/authority/ -run '^(TestJudgeDecisionSSH|TestJudgeDecisionAssertion|TestJudgeAssertionIsTheCLIJudge|TestJudgeSSHSignatureIsTheCLIJudge)\$' &&
go test -v -count=1 ./host/broker/ -run '^($(bar $INV_BROKER)$(alt $NEW_B))\$' &&
go test -v -count=1 ./host/daemon/ -run '^($(bar $INV_DAEMON))\$' &&
go test -v -count=1 ./host/boundary/ -run '^($(bar $INV_BOUNDARY))\$' &&
go test -v -count=1 ./host/store/ -run '^($(bar $INV_STORE))\$' &&
go test -v -count=1 ./host/verifygate/ -run '^($(bar $INV_VG))\$' &&
go test -v -count=1 ./cmd/ailang-worldd/ -run '^($(bar $INV_CMD))\$'"
need inv $INV_APPROVAL $NEW_A TestJudgeDecisionSSH TestJudgeDecisionAssertion TestJudgeAssertionIsTheCLIJudge TestJudgeSSHSignatureIsTheCLIJudge $INV_BROKER $NEW_B $INV_DAEMON $INV_BOUNDARY $INV_STORE $INV_VG $INV_CMD
if grep -E '^\s*--- (SKIP|FAIL)' "$L/inv.txt"; then FAILS="$FAILS inv_skip_or_fail"; fi
echo "named PASS: $(grep -cE '^--- PASS' "$L/inv.txt")" | tee -a "$L/legs.txt"
if [ "${GATES_SKIP_VERIFY_GO:-0}" = 1 ]; then echo "LEG verify_go SKIPPED (GATES_SKIP_VERIFY_GO=1)" | tee -a "$L/legs.txt"
else leg verify_go 'VERIFY_GO_LEGS=plain AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_go.sh'; fi
echo "GATES_DONE $(( $(date +%s) - s0 ))s FAILS:${FAILS:- none}" | tee -a "$L/legs.txt"
echo "GATES_DONE $(( $(date +%s) - s0 ))s FAILS:${FAILS:- none}" > "$L/ALL_DONE.$SHA"
````

**Expected:** every leg rc 0 except the base's known reds (§1.1: `pkgs` `TestWorktreeSnapshotRefusesAConfiguredFilter`;
`verify_go` `TestPublishApprovalRefusalSetWithALandedPositiveControl`, `TestMissionConfigDefaultFlowIntegration`), each
re-run alone once; `inv` and `prompt` print **named PASS** for every test. The planner's run on the final tree: §1.2 C14.

---

## 7. Executor directive notes

- **No git writes of any kind** (no add, rm, mv, commit, stash, checkout, reset, branch). **`patch -p1 -E` (§2) is the
  only writer** of the tree; no hand edit, no `gofmt -w`, no fuzz, **never `WORLD_UPDATE_FINISH_PROMPT_GOLDEN=1`** (the
  golden ships in M3's patch; regenerating it on a moved tree would launder a changed prompt). A failed `--dry-run` or a
  red acceptance → **stop and report** with the log; do not repair.
- **Never `nohup`, never background a test leg; pass a 600000 ms timeout to every long Bash call.** Each acceptance arm
  fits one call (C11: the four together took ≈ 4.5 minutes).
- **One run**: **M1a → M1b → M1c → M3. The order is fixed** (each patch assumes the previous tree). After each
  sub-milestone's acceptance, snapshot **cumulatively against the base** to `$HOME/.ailang/state/world-iter400/snap/<M>/`:
  every file added or changed since `d3b0a813` (repo-relative paths), a **`DELETED`** file (one repo path per line; S1
  deletes nothing, so it is empty), a **`STATUS`** file, and `MANIFEST.sha256` (the two plan files, `STATUS` and
  `DELETED` are not in it). **Never an `rsync` of whole directories.** The `find` takes only the top-level dirs that
  exist (M1a–M1c have only `host/`; no sub-milestone has `docs/`), rehearsed in scratch on M1a and M3 (C17):
  ```bash
  bash <<'EOF2'
  set -euo pipefail
  M=M1a   # the sub-milestone just accepted
  B=d3b0a8139bdecc36f8d083f4b701aef822340b27
  W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter400-row243-v3s1
  D="$HOME/.ailang/state/world-iter400/snap/$M"; rm -rf "$D"; mkdir -p "$D"
  cd "$W"
  { git diff --name-status --no-renames "$B" -- cmd host docs; git ls-files --others --exclude-standard -- cmd host docs | sed 's/^/A\t/'; } | sort -u -t "$(printf '\t')" -k2 > "$D/STATUS"
  awk -F'\t' '$1 != "D" {print $2}' "$D/STATUS" | while read -r p; do mkdir -p "$D/$(dirname "$p")"; cp -p "$p" "$D/$p"; done
  awk -F'\t' '$1 == "D" {print $2}' "$D/STATUS" > "$D/DELETED"
  (cd "$D" && find $(for d in cmd host docs; do [ -d "$d" ] && echo "$d"; done) -type f -print0 | sort -z | xargs -0 shasum -a 256) > "$D/MANIFEST.sha256"
  cat "$D/STATUS"; echo "DELETED: $(wc -l < "$D/DELETED")"; echo "MANIFEST: $(wc -l < "$D/MANIFEST.sha256")"
  EOF2
  ```
  Expected (`STATUS` lines / `DELETED` / MANIFEST): **M1a** 2 / 0 / 2; **M1b** 4 / 0 / 4; **M1c** 6 / 0 / 6; **M3**
  12 / 0 / 12. The controller copies the snapshot over the worktree, `git rm`s each `DELETED` path (none), commits each
  snapshot as one commit, and checks the numstat (`git diff -M --numstat d3b0a813 HEAD -- cmd host docs | awk '{s+=$1+$2}
  END {print s}'`): **M1a 387, M1b 849, M1c 1,034, M3 1,300**. A different number → stop.
- `export AILANG_BIN=$HOME/.pinned-ailang/ailang` in every shell that runs Go tests (`accept.sh` and `gates.sh` set it).
  The controller's shell is zsh: run each block under `bash <<'EOF2' … EOF2`.
- **Stop conditions**: a dry-run failure; any `--- FAIL`, `--- SKIP` or `FAIL` line in an acceptance (`ACCEPT_RED`); gofmt
  output; a numstat mismatch; any file outside §5 changed (`git status --short` after each patch must equal the
  cumulative list of §3's paragraphs).
- **Known load-sensitive tests** (the directive's list) are not this sprint's: a red there is re-run once alone.
- D-WORLD-109: nothing touches `~/.ailang/world`, `~/.ssh`, a real `/dev/tty`, a browser, or a real daemon; S1's tests
  open no store at all.

---

## 8. Risks

- **R1 — Size margin 100.** The patches are exact; any executor-side growth is a stop (§7).
- **R2 — The golden is shared across packages.** `cmd`'s test reads `../../host/broker/testdata/finish-prompt.golden.json`
  (a repo-relative path, as other cmd tests read their testdata). A future change to the prompt must change
  `FinishPromptText` and regenerate the golden with `WORLD_UPDATE_FINISH_PROMPT_GOLDEN=1` — deliberately, in review.
- **R3 — The cut leaves the judges for S1b.** Sprint 4 (C11) and sprint 5 (CF6) need them; S1b is compiled and green on
  this sprint's tip (S1-D2), so its risk is only scheduling.
- **R4 — Load.** `host/broker` and `cmd/ailang-worldd` carry the directive's load-sensitive tests; a red outside the touched
  tests is re-run alone once (§1.1).

---

## 9. Acceptance criteria (the judge scores against these)

1. **The request and the message (design §2.2)**: `CanonicalVerb` emits the exact field-order, unescaped, newline-free
   bytes (VL1's literal); `DecodeVerb` refuses unknown fields (`ErrDecode`), trailing data, reordered keys, whitespace and
   any non-canonical spelling (`ErrNotCanonical`); `Validate` refuses each of V1–V10 with its class and field, each bound
   tested on both sides (VL2); `VerbMessage` is 86 bytes in the designed layout, `VerbMessageHash` its SHA-256, a
   negative expiry refused (VL3).
2. **Domain separation, both directions (VL4)**: a mint and a decision SSHSIG never confirm (`ErrNamespace`); a verb
   SSHSIG never approves (`ApproveBySignature`) nor decides (`DecideBySignature`) (`ErrNamespace`); a mint and a decision
   assertion never confirm, a verb assertion never mints nor decides (W6 `ErrChallenge`); A2's recorded Chrome assertions
   are refused by `ConfirmByAssertion` and `CheckConfirmation` (W6) while the mint law accepts them; the new prefix and
   namespace differ from both existing ones (C5).
3. **The laws (VL5, VL6)**: `ConfirmBySignature` runs P1 (each verdict) → P2–P5 → P6, with the order pinned by pairs;
   `ConfirmByAssertion` runs P1 → P2–P5, P3w → W1–W14 → **W15 in full** (no deferred-W15 judge, design §2.3), and a
   counter-less 0-over-0 passes.
4. **`SignedVerb` (VL7)**: `Binds` refuses each differing field, the nil and the zero value; `ApprovedBy`/`Use`/`PasskeyUse`
   are empty unless verified; `Use` returns a copy; `PasskeyUse` only for webauthn.
5. **`CheckConfirmation` (VL8)** = the law minus P1 and W15: verifies past the expiry and with the counter at or past the
   artifact's; refuses a revoked approver (audit note 5, as designed), another digest / challenge / expiry, a mint- or
   decision-namespace SSHSIG, a mint assertion, an unknown format; the message is built from the caller's inputs only.
6. **The CLI's finish prompt is byte-identical (VT1)**: `FinishPromptText(…) + "Open this pull request? [y/N] "` equals the
   golden the **base** `finishPrompt` generated (C9) for three fixtures; the delegating `finishPrompt` writes the same
   bytes; `DisplayCommitter("")` is `(repository git config)`, and the CLI's `displayCommitter` agrees; every existing
   `episode_finish*_test.go` and `preset_committer_test.go` passes **unedited** (§6 `prompt`).
7. **Pure**: no route, no store read or write, no schema change, no daemon change, no `.ail` (C16; `host/approval`
   stdlib-only; `TestFrozenV1RouteTableMatchesMux`, `TestApprovalRoutesRegistered`, `TestDecisionRoutesRegistered`,
   `TestDaemonCredentialWriterSites` unedited and green).
8. **Exactly these pre-existing test files change: none** (C15). New test files: `verb_test.go`, `verb_signed_test.go`,
   `verb_check_test.go`, `verb_chrome_test.go`, `finish_prompt_test.go`, `finish_prompt_golden_test.go`; new testdata:
   `host/broker/testdata/finish-prompt.golden.json`.
9. **VT2 placement and D-WORLD-118 = A recorded** (S1-D3, S1-D4); the cut recorded with S1b's measured size (S1-D2,
   Appendix C).
10. **Size** ≤ 1,400 (`git diff -M --numstat d3b0a813 HEAD -- cmd host docs`): 1,300 measured; each sub-milestone's numstat
    as §3 (387, 849, 1,034, 1,300).
11. **Every mutant in §4 behaves as expected**, its id in the ledger; no suite kill rests on a static test.
12. **The round trip is identical** (C12) and the gates (§6) pass — `vet`, `fence`, `prompt`, `inv`, and `verify_go`
    after commit included; `./scripts/verify_ail.sh` gives the base's result.

**What an operator can do after this sprint:** nothing new at the terminal or in the app (by design: S1 is the law). The
finish prompt they see is byte-for-byte the same; the code that will judge a key's confirmation of it exists and is
proven, so the card (sprint 2) and the store (sprint 3) build on settled types.

---

## Appendix C — after this sprint (S1's close; what S1b, S2, S3 and S5 inherit)

- **Moved out (S1-D2): M2 → sprint S1b, 326 lines** (`host/authority/judge_verb.go` 129: `VerbJudgeInput`,
  `JudgeVerbSSH`, `JudgeVerbAssertion`, `RecheckVerb` + `ErrUnknownArtifact`, two lookup helpers; `judge_verb_test.go`
  197: VJ1 — every `JudgeError.Step` for the three judges, the step order pinned by pairs (`parse before store`,
  `unknown-key before law`, `stored-key before law`), RecheckVerb's no-P1/no-W15 controls and its revoked refusal).
  Compiled and green on top of `M3` (`work` branch `s1b`, tag `S1b-M2`); patch `plan-scratch/patches/S1b-M2.patch`
  (sha256 prefix `775b927ad56aa4c0`, 338 lines). S1b's planner must re-measure it on the merged tip, enumerate VJ1's
  mutants (step labels, the lookup branches, `RecheckVerb`'s kind detection and its use of `in.ExpiresAt`) and add
  `TestJudgeVerbSSH|TestJudgeVerbAssertion|TestRecheckVerb` to the gates' `inv`. S1-D13 records its one design choice.
- **Moved out (S1-D3): VT2** `broker.FinishConfirmRefusal` F1–F5 → sprint 3's M8, as design §2.3 [r1] and §8 say;
  **F6 and `DraftOnly` are not built** (D-WORLD-118 = A, S1-D4), nor S4's/S5's "B +≈ 15".
- **Audit notes (log 399) routed**: note 4 → **S5** (CF12 must say the `resume-dispatch` reuse path after a
  surface-confirmed approval whose push failed skips CF6/CF7 — its test is S5's); note 5 / gemini r2 → pinned here as
  designed (VL8 refuses an approver revoked since, S1-D10) and **carried to S5**: the ≈ 1 s window after C15 in which a
  revocation makes the CLI refuse a validly recorded approval (fail-safe: commit kept, re-run asks afresh) is CF6's to
  test; notes 1–3 → **S3** (its "four GETs" are three; P6's CR1–CR4 vs CR1–CR2 and the CX15→CX17 numbering gap) and
  **S3/S4** (the raced supersession row).
- **Residuals carried**: design RV1–RV9 unchanged (S1 builds none of their mechanisms); **P6** — a verb request's title or
  body may carry NUL (valid UTF-8; V2 as designed): S2's card renders it with `textContent` (CP3) and must not truncate at
  it; the planner of S2 should add a NUL row to CP3.
- **What S2 (the card page against fixtures) inherits**: `approval.VerbRequest`'s ten json tags in order (CP2's
  `CONFIRM_FIELDS` must equal them — a Go test can read them by reflection), `approval.CanonicalVerb` (the fixture bytes
  the card decodes), `approval.VerbMessage`/`VerbMessageHash`/`VerbMessageLen` = 86/`VerbNamespace` (CP4's golden and
  CP6's *Save the message* + signing line), `VerbRequestLifetime`, `MaxVerbRequestBytes`, and
  `broker.FinishPromptText` (XC1: the card's prompt equals it for the fixture). `ConfirmationInput`/`ConfirmationCard`
  use primitives (design §2.5), so S2 needs no store type. S2's M0 (VB8, VP12) is unaffected by S1.
- **What S3 inherits**: `DecodeVerb` + `Validate` for C10 and `requests.json`'s decode-or-refuse rows; `FinishPromptText`
  for F5; `Admit`'s verdict order for `ConsumeVerbRequest`; `SignedVerb.Use`/`PasskeyUse` for the consume's
  `ApproverUse`.
- **What S4/S5 inherit**: from S1b, the judges; from S1, `SignedVerb.Binds` (C12), `ApprovedBy` (C14's `decidedBy`, CF6's
  comparison), `CheckConfirmation` (CF6).

---

## Appendix P — the patches (apply with §2, in order; generated from the compiled scratch stages; do not edit)

### Patch M1a — the verb request and the message (design §2.2); VL1–VL3

````diff M1a
diff --git a/host/approval/verb.go b/host/approval/verb.go
new file mode 100644
index 0000000..bdd9a91
--- /dev/null
+++ b/host/approval/verb.go
@@ -0,0 +1,172 @@
+package approval
+
+// Row 243 V3 (design_docs/planned/w-world-address-bar-verbs-v3.md §2.2): the verb request a
+// confirm-on-surface `episode finish` stores, and the 86-byte message an enrolled key signs to
+// confirm it. A third kind beside the mint message (SignedMessage) and the decision message
+// (DecisionMessage): its own prefix and SSHSIG namespace, so no signature of one kind is accepted
+// as another (design §2.2 "Domain separation"). Pure: no I/O, no clock, no randomness.
+
+import (
+	"bytes"
+	"crypto/sha256"
+	"encoding/binary"
+	"encoding/json"
+	"errors"
+	"fmt"
+	"reflect"
+	"time"
+	"unicode/utf8"
+)
+
+// VerbRequestSchema is the schema of every stored verb request.
+const VerbRequestSchema = "world/verb-request/v1"
+
+// VerbRequestLifetime is a verb request's lifetime, as a mint request's (RequestLifetime).
+const VerbRequestLifetime = 10 * time.Minute
+
+// MaxVerbRequestBytes bounds a verb request's canonical bytes: a body of at most 60,000 bytes
+// (cmd/ailang-worldd reportMarkdownMax) plus the prompt, scope, title and argv.
+const MaxVerbRequestBytes = 128 << 10
+
+// VerbFinish is the only admissible verb.
+const VerbFinish = "episode finish"
+
+// VerbPrefix starts every signed verb message (design §2.2).
+const VerbPrefix = "world-verb/v1"
+
+// VerbNamespace is the SSHSIG namespace every verb confirmation is signed under.
+const VerbNamespace = "world-verb@sunholo"
+
+// VerbMessageLen is the length of VerbMessage: 13 + 1 + 32 + 32 + 8.
+const VerbMessageLen = len(VerbPrefix) + 1 + 32 + 32 + 8
+
+// verbBodyMax is the pull request body's cap (cmd/ailang-worldd episode_markdown.go reportMarkdownMax).
+const verbBodyMax = 60000
+
+// verbArgvMax bounds the canonical argv's words.
+const verbArgvMax = 64
+
+// VerbRequest is a stored, single-use confirmation request. Its canonical JSON (CanonicalVerb), in
+// this field order, is what is digested, stored and shown on the card.
+type VerbRequest struct {
+	Schema    string   `json:"schema"`    // VerbRequestSchema
+	Verb      string   `json:"verb"`      // VerbFinish, the only admissible value
+	Argv      []string `json:"argv"`      // the canonical argv: ailang-worldd episode finish <episode> …
+	Episode   string   `json:"episode"`   // ^[a-z0-9][a-z0-9-]{0,63}$
+	Prompt    string   `json:"prompt"`    // the terminal's lines (broker.FinishPromptText), without the question
+	Scope     string   `json:"scope"`     // the finish approval scope the daemon records
+	Title     string   `json:"title"`     // the pull request's title
+	Body      string   `json:"body"`      // the full pull request body
+	Committer string   `json:"committer"` // the pinned committer, raw ("" allowed)
+	Protected []string `json:"protected"` // the protected paths the commit changed
+}
+
+// ErrVerbRequest is a verb request that Validate refuses; the text names the field.
+var ErrVerbRequest = errors.New("approval: bad verb request")
+
+// Validate refuses a verb request the daemon or the CLI must never act on (design §2.2). Each
+// refusal wraps ErrSchema, ErrString, ErrEpisode or ErrVerbRequest and names its JSON field.
+func (r VerbRequest) Validate() error {
+	if r.Schema != VerbRequestSchema { // V1
+		return fmt.Errorf("%w: schema is %q, want %q", ErrSchema, r.Schema, VerbRequestSchema)
+	}
+	if err := stringLeaves(reflect.ValueOf(r), "", func(path, s string) error { // V2
+		if !utf8.ValidString(s) {
+			return fmt.Errorf("%w: %s", ErrString, path)
+		}
+		return nil
+	}); err != nil {
+		return err
+	}
+	if r.Verb != VerbFinish { // V3
+		return fmt.Errorf("%w: verb is %q, want %q", ErrVerbRequest, r.Verb, VerbFinish)
+	}
+	if !episodePattern.MatchString(r.Episode) { // V4
+		return fmt.Errorf("%w: episode %q", ErrEpisode, r.Episode)
+	}
+	if len(r.Argv) < 4 || r.Argv[0] != "ailang-worldd" || r.Argv[1] != "episode" || r.Argv[2] != "finish" || r.Argv[3] != r.Episode { // V5
+		return fmt.Errorf("%w: argv does not start ailang-worldd episode finish %s", ErrVerbRequest, r.Episode)
+	}
+	if len(r.Argv) > verbArgvMax { // V6
+		return fmt.Errorf("%w: argv has %d words, at most %d", ErrVerbRequest, len(r.Argv), verbArgvMax)
+	}
+	for i, w := range r.Argv { // V7
+		for j := 0; j < len(w); j++ {
+			if w[j] < 0x20 || w[j] == 0x7f {
+				return fmt.Errorf("%w: argv[%d] holds a control byte", ErrVerbRequest, i)
+			}
+		}
+	}
+	for _, f := range []struct{ name, val string }{{"prompt", r.Prompt}, {"scope", r.Scope}, {"title", r.Title}} { // V8
+		if f.val == "" {
+			return fmt.Errorf("%w: %s is empty", ErrVerbRequest, f.name)
+		}
+	}
+	if len(r.Body) > verbBodyMax { // V9
+		return fmt.Errorf("%w: body is %d bytes, at most %d", ErrVerbRequest, len(r.Body), verbBodyMax)
+	}
+	c, err := CanonicalVerb(r)
+	if err != nil {
+		return err
+	}
+	if len(c) > MaxVerbRequestBytes { // V10
+		return fmt.Errorf("%w: the request is %d bytes, at most %d", ErrVerbRequest, len(c), MaxVerbRequestBytes)
+	}
+	return nil
+}
+
+// CanonicalVerb is the compact JSON of r in struct field order, with HTML escaping off and no
+// trailing newline (Canonical's encoder settings). It does not call Validate.
+func CanonicalVerb(r VerbRequest) ([]byte, error) {
+	var buf bytes.Buffer
+	enc := json.NewEncoder(&buf)
+	enc.SetEscapeHTML(false)
+	if err := enc.Encode(r); err != nil {
+		return nil, err
+	}
+	return bytes.TrimSuffix(buf.Bytes(), []byte("\n")), nil
+}
+
+// DecodeVerb is the strict inverse of CanonicalVerb (Decode's strictness): unknown fields, trailing
+// data and any input that does not re-encode to the same bytes are refused. It does not call
+// Validate; callers do.
+func DecodeVerb(b []byte) (VerbRequest, error) {
+	dec := json.NewDecoder(bytes.NewReader(b))
+	dec.DisallowUnknownFields()
+	var r VerbRequest
+	if err := dec.Decode(&r); err != nil {
+		return VerbRequest{}, fmt.Errorf("%w: %v", ErrDecode, err)
+	}
+	c, err := CanonicalVerb(r)
+	if err != nil {
+		return VerbRequest{}, err
+	}
+	if !bytes.Equal(c, b) {
+		return VerbRequest{}, ErrNotCanonical
+	}
+	return r, nil
+}
+
+// VerbMessage is the 86-byte message a confirmation signs (design §2.2):
+// VerbPrefix ‖ 0x00 ‖ digest ‖ challenge ‖ expiry (uint64, big-endian).
+func VerbMessage(digest, challenge [32]byte, expiry int64) ([]byte, error) {
+	if expiry < 0 {
+		return nil, ErrNegativeExpiry
+	}
+	m := make([]byte, 0, VerbMessageLen)
+	m = append(m, VerbPrefix...)
+	m = append(m, 0)
+	m = append(m, digest[:]...)
+	m = append(m, challenge[:]...)
+	m = binary.BigEndian.AppendUint64(m, uint64(expiry))
+	return m, nil
+}
+
+// VerbMessageHash is the SHA-256 of VerbMessage: a passkey's WebAuthn challenge.
+func VerbMessageHash(digest, challenge [32]byte, expiry int64) ([32]byte, error) {
+	m, err := VerbMessage(digest, challenge, expiry)
+	if err != nil {
+		return [32]byte{}, err
+	}
+	return sha256.Sum256(m), nil
+}
diff --git a/host/approval/verb_test.go b/host/approval/verb_test.go
new file mode 100644
index 0000000..09b17b6
--- /dev/null
+++ b/host/approval/verb_test.go
@@ -0,0 +1,215 @@
+package approval_test
+
+// Row 243 V3 sprint 1 M1a (design §2.2, §6 VL1-VL3): the verb request's canonical bytes, its strict
+// decoder, every Validate rule, and the 86-byte message's layout.
+
+import (
+	"bytes"
+	"crypto/sha256"
+	"encoding/binary"
+	"errors"
+	"strings"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+)
+
+// vrRequest is a valid verb request: a draft finish of w1 whose title and body need no escaping
+// but hold the bytes CanonicalVerb must not HTML-escape.
+func vrRequest() approval.VerbRequest {
+	return approval.VerbRequest{
+		Schema:    approval.VerbRequestSchema,
+		Verb:      approval.VerbFinish,
+		Argv:      []string{"ailang-worldd", "episode", "finish", "w1", "--preset", "p", "--confirm-on-surface"},
+		Episode:   "w1",
+		Prompt:    "committer: (repository git config)\nrepository: o/r\n",
+		Scope:     "github-pr:o/r#remote=origin&head=world/w1",
+		Title:     "Fix <b> & tidy",
+		Body:      "line 1\nline <2> & \"3\"\n",
+		Committer: "",
+		Protected: []string{"a/b", "c"},
+	}
+}
+
+// VL1: the canonical bytes are compact, in field order, with HTML escaping off and no trailing
+// newline; DecodeVerb inverts them exactly.
+func TestVerbRequestRoundTrip(t *testing.T) {
+	r := vrRequest()
+	c, err := approval.CanonicalVerb(r)
+	if err != nil {
+		t.Fatal(err)
+	}
+	want := `{"schema":"world/verb-request/v1","verb":"episode finish","argv":["ailang-worldd","episode","finish","w1","--preset","p","--confirm-on-surface"],` +
+		`"episode":"w1","prompt":"committer: (repository git config)\nrepository: o/r\n","scope":"github-pr:o/r#remote=origin&head=world/w1",` +
+		`"title":"Fix <b> & tidy","body":"line 1\nline <2> & \"3\"\n","committer":"","protected":["a/b","c"]}`
+	if string(c) != want {
+		t.Fatalf("CanonicalVerb =\n%s\nwant\n%s", c, want)
+	}
+	got, err := approval.DecodeVerb(c)
+	if err != nil {
+		t.Fatal(err)
+	}
+	c2, _ := approval.CanonicalVerb(got)
+	if !bytes.Equal(c, c2) || got.Title != r.Title || len(got.Argv) != len(r.Argv) || got.Protected[1] != "c" {
+		t.Fatalf("round trip: %+v", got)
+	}
+	if err := got.Validate(); err != nil {
+		t.Fatalf("the valid request: %v", err)
+	}
+	nilProtected := vrRequest()
+	nilProtected.Protected = nil
+	c, _ = approval.CanonicalVerb(nilProtected)
+	if back, err := approval.DecodeVerb(c); err != nil || back.Protected != nil || !strings.HasSuffix(string(c), `"protected":null}`) {
+		t.Fatalf("nil protected: %v %q", err, c)
+	}
+}
+
+// vrEscaped is "<b>" as an HTML-escaping encoder spells it (backslash u003c b backslash u003e).
+var vrEscaped = string([]byte{'\\', 'u', '0', '0', '3', 'c', 'b', '\\', 'u', '0', '0', '3', 'e'})
+
+// VL1: unknown fields are ErrDecode; trailing data, reordered keys, whitespace and any other
+// non-canonical spelling are ErrNotCanonical.
+func TestVerbRequestDecodeRefusals(t *testing.T) {
+	c, err := approval.CanonicalVerb(vrRequest())
+	if err != nil {
+		t.Fatal(err)
+	}
+	s := string(c)
+	for _, tc := range []struct {
+		name string
+		in   string
+		want error
+	}{
+		{"unknown field", strings.Replace(s, `{"schema"`, `{"typed":"y","schema"`, 1), approval.ErrDecode},
+		{"unknown field at the end", strings.TrimSuffix(s, "}") + `,"y":true}`, approval.ErrDecode},
+		{"not JSON", "y", approval.ErrDecode},
+		{"empty", "", approval.ErrDecode},
+		{"trailing data", s + "x", approval.ErrNotCanonical},
+		{"trailing newline", s + "\n", approval.ErrNotCanonical},
+		{"a second object", s + s, approval.ErrNotCanonical},
+		{"reordered keys", `{"verb":"episode finish","schema":"world/verb-request/v1"` + s[len(`{"schema":"world/verb-request/v1","verb":"episode finish"`):], approval.ErrNotCanonical},
+		{"whitespace", strings.Replace(s, `"schema":`, `"schema": `, 1), approval.ErrNotCanonical},
+		{"html-escaped", strings.Replace(s, "<b>", vrEscaped, 1), approval.ErrNotCanonical},
+	} {
+		if r, err := approval.DecodeVerb([]byte(tc.in)); !errors.Is(err, tc.want) || r.Schema != "" {
+			t.Errorf("%s: (%q, %v), want %v", tc.name, r.Schema, err, tc.want)
+		}
+	}
+}
+
+// VL2: one case per Validate rule (design §2.2), each with its boundary where it has one.
+func TestVerbRequestValidate(t *testing.T) {
+	argv := func(n int) []string {
+		a := []string{"ailang-worldd", "episode", "finish", "w1"}
+		for len(a) < n {
+			a = append(a, "--allow-failing")
+		}
+		return a
+	}
+	for _, tc := range []struct {
+		name string
+		edit func(r *approval.VerbRequest)
+		want error // nil: admissible
+		text string
+	}{
+		{"V1 schema", func(r *approval.VerbRequest) { r.Schema = "world/mint-request/v1" }, approval.ErrSchema, "schema"},
+		{"V2 utf-8 title", func(r *approval.VerbRequest) { r.Title = "a\xffb" }, approval.ErrString, ": title"},
+		{"V2 utf-8 protected", func(r *approval.VerbRequest) { r.Protected = []string{"ok", "\xc3"} }, approval.ErrString, "protected[1]"},
+		{"V3 verb", func(r *approval.VerbRequest) { r.Verb = "episode start" }, approval.ErrVerbRequest, "verb"},
+		{"V4 episode", func(r *approval.VerbRequest) { r.Episode, r.Argv[3] = "W1", "W1" }, approval.ErrEpisode, "episode"},
+		{"V5 argv short", func(r *approval.VerbRequest) { r.Argv = r.Argv[:3] }, approval.ErrVerbRequest, "argv does not start"},
+		{"V5 argv[0]", func(r *approval.VerbRequest) { r.Argv[0] = "world" }, approval.ErrVerbRequest, "argv does not start"},
+		{"V5 argv[1]", func(r *approval.VerbRequest) { r.Argv[1] = "session" }, approval.ErrVerbRequest, "argv does not start"},
+		{"V5 argv[2]", func(r *approval.VerbRequest) { r.Argv[2] = "start" }, approval.ErrVerbRequest, "argv does not start"},
+		{"V5 argv[3]", func(r *approval.VerbRequest) { r.Argv[3] = "w2" }, approval.ErrVerbRequest, "argv does not start"},
+		{"V5 argv of exactly four", func(r *approval.VerbRequest) { r.Argv = r.Argv[:4] }, nil, ""},
+		{"V6 64 words", func(r *approval.VerbRequest) { r.Argv = argv(64) }, nil, ""},
+		{"V6 65 words", func(r *approval.VerbRequest) { r.Argv = argv(65) }, approval.ErrVerbRequest, "65 words"},
+		{"V7 newline", func(r *approval.VerbRequest) { r.Argv[5] = "p\nq" }, approval.ErrVerbRequest, "argv[5] holds a control byte"},
+		{"V7 0x1f", func(r *approval.VerbRequest) { r.Argv[4] = "\x1f" }, approval.ErrVerbRequest, "argv[4] holds a control byte"},
+		{"V7 DEL", func(r *approval.VerbRequest) { r.Argv[6] = "a\x7f" }, approval.ErrVerbRequest, "argv[6] holds a control byte"},
+		{"V7 space and tilde", func(r *approval.VerbRequest) { r.Argv[5] = " ~p" }, nil, ""},
+		{"V8 prompt", func(r *approval.VerbRequest) { r.Prompt = "" }, approval.ErrVerbRequest, "prompt is empty"},
+		{"V8 scope", func(r *approval.VerbRequest) { r.Scope = "" }, approval.ErrVerbRequest, "scope is empty"},
+		{"V8 title", func(r *approval.VerbRequest) { r.Title = "" }, approval.ErrVerbRequest, "title is empty"},
+		{"V8 committer may be empty", func(r *approval.VerbRequest) { r.Committer, r.Body, r.Protected = "", "", nil }, nil, ""},
+		{"V9 body 60000", func(r *approval.VerbRequest) { r.Body = strings.Repeat("b", 60000) }, nil, ""},
+		{"V9 body 60001", func(r *approval.VerbRequest) { r.Body = strings.Repeat("b", 60001) }, approval.ErrVerbRequest, "body is 60001 bytes"},
+		{"V10 canonical 131072", func(r *approval.VerbRequest) { r.Prompt = strings.Repeat("p", 131072-vrLen(t)+1) }, nil, ""},
+		{"V10 canonical 131073", func(r *approval.VerbRequest) { r.Prompt = strings.Repeat("p", 131073-vrLen(t)+1) }, approval.ErrVerbRequest, "131073 bytes"},
+	} {
+		r := vrRequest()
+		tc.edit(&r)
+		err := r.Validate()
+		if tc.want == nil {
+			if err != nil {
+				t.Errorf("%s: %v, want admissible", tc.name, err)
+			}
+			continue
+		}
+		if !errors.Is(err, tc.want) || !strings.Contains(err.Error(), tc.text) {
+			t.Errorf("%s: %v, want %v naming %q", tc.name, err, tc.want, tc.text)
+		}
+	}
+}
+
+// vrLen is the canonical length of vrRequest() with the one-byte prompt "p".
+func vrLen(t *testing.T) int {
+	t.Helper()
+	r := vrRequest()
+	r.Prompt = "p"
+	c, err := approval.CanonicalVerb(r)
+	if err != nil {
+		t.Fatal(err)
+	}
+	return len(c)
+}
+
+// VL3: the message is 86 bytes in design §2.2's layout; its hash is the SHA-256 of those bytes; a
+// negative expiry has no message; the three prefixes differ at byte 6 and none prefixes another.
+func TestVerbMessageLayout(t *testing.T) {
+	if approval.VerbPrefix != "world-verb/v1" || approval.VerbNamespace != "world-verb@sunholo" || approval.VerbMessageLen != 86 ||
+		approval.VerbRequestSchema != "world/verb-request/v1" || approval.MaxVerbRequestBytes != 131072 || approval.VerbRequestLifetime.Seconds() != 600 {
+		t.Fatalf("constants: %q %q %d %q %d %v", approval.VerbPrefix, approval.VerbNamespace, approval.VerbMessageLen,
+			approval.VerbRequestSchema, approval.MaxVerbRequestBytes, approval.VerbRequestLifetime)
+	}
+	var d, c [32]byte
+	for i := range d {
+		d[i], c[i] = byte(i), byte(0x80+i)
+	}
+	const expiry = int64(0x0102030405060708)
+	m, err := approval.VerbMessage(d, c, expiry)
+	if err != nil {
+		t.Fatal(err)
+	}
+	if len(m) != 86 || string(m[:13]) != "world-verb/v1" || m[13] != 0 || !bytes.Equal(m[14:46], d[:]) || !bytes.Equal(m[46:78], c[:]) ||
+		binary.BigEndian.Uint64(m[78:]) != uint64(expiry) || m[78] != 0x01 || m[85] != 0x08 {
+		t.Fatalf("VerbMessage = %x", m)
+	}
+	h, err := approval.VerbMessageHash(d, c, expiry)
+	if err != nil || h != sha256.Sum256(m) {
+		t.Fatalf("VerbMessageHash = %x, %v", h, err)
+	}
+	if m, err := approval.VerbMessage(d, c, -1); m != nil || !errors.Is(err, approval.ErrNegativeExpiry) {
+		t.Fatalf("negative expiry: %x, %v", m, err)
+	}
+	if h, err := approval.VerbMessageHash(d, c, -1); h != ([32]byte{}) || !errors.Is(err, approval.ErrNegativeExpiry) {
+		t.Fatalf("negative expiry hash: %x, %v", h, err)
+	}
+	if m, err := approval.VerbMessage(d, c, 0); err != nil || len(m) != 86 {
+		t.Fatalf("expiry 0: %v", err)
+	}
+	mint, _ := approval.SignedMessage(d, c, expiry)
+	decision, _ := approval.DecisionMessage(approval.VerbApprove, d, c, expiry)
+	if approval.VerbPrefix[:6] != approval.MessagePrefix[:6] || approval.VerbPrefix[6] == approval.MessagePrefix[6] || approval.VerbPrefix[6] == approval.DecisionPrefix[6] {
+		t.Fatalf("prefixes: %q %q %q", approval.VerbPrefix, approval.MessagePrefix, approval.DecisionPrefix)
+	}
+	for name, other := range map[string][]byte{"mint": mint, "decision": decision} {
+		if bytes.HasPrefix(m, other[:len(other)-40]) || bytes.HasPrefix(other, m[:20]) || bytes.Equal(m, other) {
+			t.Errorf("the verb message and the %s message share a prefix", name)
+		}
+	}
+	if approval.VerbNamespace == approval.Namespace || approval.VerbNamespace == approval.DecisionNamespace {
+		t.Fatalf("namespaces: %q %q %q", approval.VerbNamespace, approval.Namespace, approval.DecisionNamespace)
+	}
+}
````

### Patch M1b — the laws and `SignedVerb` (design §2.3); VL4–VL7

````diff M1b
diff --git a/host/approval/verb_signed.go b/host/approval/verb_signed.go
new file mode 100644
index 0000000..44cce54
--- /dev/null
+++ b/host/approval/verb_signed.go
@@ -0,0 +1,140 @@
+package approval
+
+// Row 243 V3 sprint 1 (design §2.3): the pure law of a verb confirmation by an enrolled key, over the
+// 86-byte verb message (verb.go). ConfirmBySignature and ConfirmByAssertion are the mint laws'
+// shape (signed.go, passkey_assert.go) with the verb message: P1 admission, P2-P5 checkApprover,
+// then the signature, and for a passkey W1-W14 (assertionSteps) and W15 in full. No I/O; now is an
+// input.
+
+import (
+	"bytes"
+	"fmt"
+)
+
+// SignedVerb is a verified verb confirmation. Its fields are unexported, so only ConfirmBySignature,
+// and ConfirmByAssertion can return one with ok set; the zero value binds nothing.
+type SignedVerb struct {
+	ok           bool
+	digest       [32]byte
+	challenge    [32]byte
+	expiry       int64
+	name, format string
+	key          []byte
+	credentialID []byte // webauthn: the passkey that signed
+	signCount    uint32 // webauthn: the assertion's counter
+}
+
+// VerbSignatureInput is what ConfirmBySignature judges: the stored request's admission input, its
+// challenge, the armored SSHSIG and the enrolled approver its key names.
+type VerbSignatureInput struct {
+	Admit     AdmitInput // row times, consumed, the stored digest and the SHA-256 of the stored bytes
+	Challenge [32]byte   // the request row's challenge
+	Signature []byte     // the armored SSHSIG
+	Approver  Approver
+}
+
+// VerbAssertionInput is what ConfirmByAssertion judges: the stored request's admission input, its
+// challenge, the assertion file and the enrolled approver its credential ID names, Passkey set.
+type VerbAssertionInput struct {
+	Admit     AdmitInput
+	Challenge [32]byte
+	Response  []byte
+	Approver  Approver
+}
+
+// verifyVerbSSHSIG parses an armored SSHSIG and verifies it over VerbMessage(digest, challenge,
+// expiry) under VerbNamespace with key.
+func verifyVerbSSHSIG(armored []byte, digest, challenge [32]byte, expiry int64, key SSHPublicKey) error {
+	s, err := ParseSSHSIG(armored)
+	if err != nil {
+		return err
+	}
+	msg, err := VerbMessage(digest, challenge, expiry)
+	if err != nil {
+		return err
+	}
+	return VerifySSHSIG(s, VerbNamespace, msg, key)
+}
+
+// ConfirmBySignature is the verb signature law (P1-P6, in this order): admission first, so a replay
+// says "already used"; the approver's state before the cryptography; then the SSHSIG over the verb
+// message of the stored digest, the row's challenge and the row's expiry.
+func ConfirmBySignature(in VerbSignatureInput) (*SignedVerb, error) {
+	if v := Admit(in.Admit); v != Admitted { // P1
+		return nil, &NotAdmittedError{Verdict: v}
+	}
+	if err := checkApprover(in.Admit, in.Approver, FormatSSH); err != nil { // P2-P5
+		return nil, err
+	}
+	if err := verifyVerbSSHSIG(in.Signature, in.Admit.StoredDigest, in.Challenge, in.Admit.ExpiresAt, in.Approver.Key); err != nil { // P6
+		return nil, err
+	}
+	return &SignedVerb{ok: true, digest: in.Admit.StoredDigest, challenge: in.Challenge, expiry: in.Admit.ExpiresAt,
+		name: in.Approver.Name, format: FormatSSH, key: bytes.Clone(in.Approver.Key.Blob)}, nil
+}
+
+// ConfirmByAssertion is the verb assertion law (P1-P5, P3w, W1-W15, in this order), as
+// ApproveByAssertion with the challenge VerbMessageHash. W15 runs in full: the counter moves in the
+// consume, which runs after the world write, so a retry of the same artifact still passes (design
+// §2.3).
+func ConfirmByAssertion(in VerbAssertionInput) (*SignedVerb, error) {
+	if v := Admit(in.Admit); v != Admitted { // P1
+		return nil, &NotAdmittedError{Verdict: v}
+	}
+	if err := checkApprover(in.Admit, in.Approver, FormatWebAuthn); err != nil { // P2-P5, P3w
+		return nil, err
+	}
+	h, hErr := VerbMessageHash(in.Admit.StoredDigest, in.Challenge, in.Admit.ExpiresAt)
+	ad, err := assertionSteps(in.Response, in.Approver, h, hErr) // W1-W14
+	if err != nil {
+		return nil, err
+	}
+	p := in.Approver.Passkey
+	if (ad.SignCount != 0 || p.SignCount != 0) && ad.SignCount <= p.SignCount { // W15
+		return nil, fmt.Errorf("%w: enrolled %d, this assertion %d; if you did not cause this, run: ailang-worldd approver revoke %s",
+			ErrSignCountRegression, p.SignCount, ad.SignCount, in.Approver.Name)
+	}
+	return &SignedVerb{ok: true, digest: in.Admit.StoredDigest, challenge: in.Challenge, expiry: in.Admit.ExpiresAt,
+		name: in.Approver.Name, format: FormatWebAuthn, key: p.Key.Canonical(), credentialID: bytes.Clone(p.CredentialID), signCount: ad.SignCount}, nil
+}
+
+// Binds reports whether s is a verified confirmation of exactly this digest, challenge and expiry.
+func (s *SignedVerb) Binds(digest, challenge [32]byte, expiry int64) error {
+	switch {
+	case s == nil || !s.ok:
+		return ErrNotVerified
+	case s.digest != digest:
+		return &BindingError{Field: "digest"}
+	case s.challenge != challenge:
+		return &BindingError{Field: "challenge"}
+	case s.expiry != expiry:
+		return &BindingError{Field: "expiry"}
+	}
+	return nil
+}
+
+// ApprovedBy is the confirmation's approved_by and the finish approval's decidedBy, format ":" name
+// ("webauthn:mark-mac"); "" unless verified.
+func (s *SignedVerb) ApprovedBy() string {
+	if s == nil || !s.ok {
+		return ""
+	}
+	return s.format + ":" + s.name
+}
+
+// Use is what the consume binds in SQL: the approver's name, format and key (a copy).
+func (s *SignedVerb) Use() (name, format string, key []byte, ok bool) {
+	if s == nil || !s.ok {
+		return "", "", nil, false
+	}
+	return s.name, s.format, bytes.Clone(s.key), true
+}
+
+// PasskeyUse is what the consume binds for a passkey: the credential ID (a copy) and the
+// assertion's counter; ok only for a verified webauthn confirmation.
+func (s *SignedVerb) PasskeyUse() (credentialID []byte, signCount uint32, ok bool) {
+	if s == nil || !s.ok || s.format != FormatWebAuthn {
+		return nil, 0, false
+	}
+	return bytes.Clone(s.credentialID), s.signCount, true
+}
diff --git a/host/approval/verb_signed_test.go b/host/approval/verb_signed_test.go
new file mode 100644
index 0000000..eed528e
--- /dev/null
+++ b/host/approval/verb_signed_test.go
@@ -0,0 +1,322 @@
+package approval_test
+
+// Row 243 V3 sprint 1 M1b (design §2.3, §6 VL4-VL7): the verb law with approvaltest keys and
+// passkeys (D-WORLD-109): domain separation both directions, the signature law's P1-P6, the
+// assertion law's P1-P5, P3w, W1-W15, and what a SignedVerb binds and hands its consume.
+
+import (
+	"bytes"
+	"errors"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
+)
+
+var vD, vC = [32]byte{0x76, 1}, [32]byte{0x76, 2}
+
+const vCreated, vExpiry, vNow = int64(1000), int64(1600), int64(1010)
+
+var vPasskey = approvaltest.NewPasskey(approval.AlgES256, 71)
+
+// vAdmit is an admitted request: created 1000, expires 1600, now 1010, bytes = digest vD.
+func vAdmit() approval.AdmitInput {
+	return approval.AdmitInput{Now: vNow, CreatedAt: vCreated, ExpiresAt: vExpiry, StoredDigest: vD, BytesDigest: vD}
+}
+
+// vMsg and vHash are what a confirmation of (digest, challenge, expiry) signs.
+func vMsg(t *testing.T, d, c [32]byte, e int64) []byte {
+	t.Helper()
+	m, err := approval.VerbMessage(d, c, e)
+	if err != nil {
+		t.Fatal(err)
+	}
+	return m
+}
+
+func vHash(t *testing.T, d, c [32]byte, e int64) [32]byte {
+	t.Helper()
+	h, err := approval.VerbMessageHash(d, c, e)
+	if err != nil {
+		t.Fatal(err)
+	}
+	return h
+}
+
+// vsInput is a valid signature confirmation: s's SSHSIG over the request's verb message, s
+// enrolled as mark-ssh at 900.
+func vsInput(t *testing.T, s *approvaltest.Signer) approval.VerbSignatureInput {
+	t.Helper()
+	return approval.VerbSignatureInput{Admit: vAdmit(), Challenge: vC, Signature: s.Sign(approval.VerbNamespace, vMsg(t, vD, vC, vExpiry)),
+		Approver: approval.Approver{Name: "mark-ssh", Format: approval.FormatSSH, Key: s.PublicKey(), EnrolledAt: 900}}
+}
+
+// vaInput is a valid assertion confirmation: vPasskey (enrolled as mark-mac at 900, stored counter
+// 1) asserts the request's verb message hash with counter count.
+func vaInput(t *testing.T, count uint32, e ...approvaltest.Edit) approval.VerbAssertionInput {
+	t.Helper()
+	return approval.VerbAssertionInput{Admit: vAdmit(), Challenge: vC,
+		Response: vPasskey.Assert(paOrigin, paRPID, vHash(t, vD, vC, vExpiry), count, e...), Approver: paApprover(t, vPasskey)}
+}
+
+// VL4: domain separation, both directions, for both key kinds. A mint or decision SSHSIG never
+// confirms a verb (wrong namespace); a verb SSHSIG never approves a mint nor decides (wrong
+// namespace); a mint or decision assertion never confirms (W6), nor a verb assertion mints or
+// decides (W6).
+func TestVerbDomainSeparation(t *testing.T) {
+	s := approvaltest.NewSigner(72)
+	mintMsg, _ := approval.SignedMessage(vD, vC, vExpiry)
+	decMsg, _ := approval.DecisionMessage(approval.VerbApprove, vD, vC, vExpiry)
+	for _, tc := range []struct {
+		name string
+		sig  []byte
+		want error
+	}{
+		{"a mint SSHSIG", s.Sign(approval.Namespace, mintMsg), approval.ErrNamespace},
+		{"a decision SSHSIG", s.Sign(approval.DecisionNamespace, decMsg), approval.ErrNamespace},
+		{"the verb message, mint namespace", s.Sign(approval.Namespace, vMsg(t, vD, vC, vExpiry)), approval.ErrNamespace},
+		{"the verb message, decision namespace", s.Sign(approval.DecisionNamespace, vMsg(t, vD, vC, vExpiry)), approval.ErrNamespace},
+		{"the mint message, verb namespace", s.Sign(approval.VerbNamespace, mintMsg), approval.ErrBadSignature},
+		{"the decision message, verb namespace", s.Sign(approval.VerbNamespace, decMsg), approval.ErrBadSignature},
+		{"the verb message without its zero byte", s.Sign(approval.VerbNamespace, append([]byte(approval.VerbPrefix), vMsg(t, vD, vC, vExpiry)[14:]...)), approval.ErrBadSignature},
+	} {
+		in := vsInput(t, s)
+		in.Signature = tc.sig
+		if v, err := approval.ConfirmBySignature(in); v != nil || !errors.Is(err, tc.want) {
+			t.Errorf("%s: (%v, %v), want %v", tc.name, v, err, tc.want)
+		}
+	}
+	verb := vsInput(t, s).Signature
+	mint := approval.SignatureInput{Admit: vAdmit(), Challenge: vC, Signature: verb,
+		Approver: approval.Approver{Name: "mark-ssh", Format: approval.FormatSSH, Key: s.PublicKey(), EnrolledAt: 900}}
+	if a, err := approval.ApproveBySignature(mint); a != nil || !errors.Is(err, approval.ErrNamespace) {
+		t.Errorf("a verb SSHSIG minted: (%v, %v)", a, err)
+	}
+	dec := approval.DecisionSignatureInput{Verb: approval.VerbApprove, Packet: vD, Request: vC, Expiry: vExpiry, Now: vNow, Signature: verb, Approver: mint.Approver}
+	if d, err := approval.DecideBySignature(dec); d != nil || !errors.Is(err, approval.ErrNamespace) {
+		t.Errorf("a verb SSHSIG decided: (%v, %v)", d, err)
+	}
+	if v, err := approval.ConfirmBySignature(vsInput(t, s)); err != nil || v.ApprovedBy() != "ssh:mark-ssh" {
+		t.Fatalf("control: (%v, %v)", v, err)
+	}
+
+	mintHash, _ := approval.MessageHash(vD, vC, vExpiry)
+	decHash, _ := approval.DecisionMessageHash(approval.VerbApprove, vD, vC, vExpiry)
+	for name, h := range map[string][32]byte{"a mint assertion": mintHash, "a decision assertion": decHash} {
+		in := vaInput(t, 2)
+		in.Response = vPasskey.Assert(paOrigin, paRPID, h, 2)
+		if v, err := approval.ConfirmByAssertion(in); v != nil || !errors.Is(err, approval.ErrChallenge) {
+			t.Errorf("%s confirmed: (%v, %v), want ErrChallenge (W6)", name, v, err)
+		}
+	}
+	va := vaInput(t, 2)
+	if a, err := approval.ApproveByAssertion(approval.AssertionInput{Admit: va.Admit, Challenge: vC, Response: va.Response, Approver: va.Approver}); a != nil || !errors.Is(err, approval.ErrChallenge) {
+		t.Errorf("a verb assertion minted: (%v, %v)", a, err)
+	}
+	if d, err := approval.DecideByAssertion(approval.DecisionAssertionInput{Verb: approval.VerbApprove, Packet: vD, Request: vC, Expiry: vExpiry, Now: vNow,
+		Response: va.Response, Approver: va.Approver}); d != nil || !errors.Is(err, approval.ErrChallenge) {
+		t.Errorf("a verb assertion decided: (%v, %v)", d, err)
+	}
+	if v, err := approval.ConfirmByAssertion(va); err != nil || v.ApprovedBy() != "webauthn:mark-mac" {
+		t.Fatalf("control: (%v, %v)", v, err)
+	}
+}
+
+// VL5: the signature law's order: P1 each verdict (before the approver), P2-P5 (before the
+// cryptography), then P6: a bad signature, a signature over another digest, challenge or expiry.
+func TestConfirmBySignatureLaw(t *testing.T) {
+	s, other := approvaltest.NewSigner(73), approvaltest.NewSigner(74)
+	for _, tc := range []struct {
+		name    string
+		edit    func(in *approval.VerbSignatureInput)
+		want    error
+		verdict approval.Verdict
+	}{
+		{"P1 consumed", func(in *approval.VerbSignatureInput) { in.Admit.Consumed = true }, nil, approval.RefuseConsumed},
+		{"P1 consumed before P4", func(in *approval.VerbSignatureInput) { in.Admit.Consumed, in.Approver.RevokedAt = true, 950 }, nil, approval.RefuseConsumed},
+		{"P1 clock behind", func(in *approval.VerbSignatureInput) { in.Admit.Now = vCreated - 1 }, nil, approval.RefuseClockBehind},
+		{"P1 expired", func(in *approval.VerbSignatureInput) { in.Admit.Now = vExpiry }, nil, approval.RefuseExpired},
+		{"P1 digest", func(in *approval.VerbSignatureInput) { in.Admit.BytesDigest[0] ^= 1 }, nil, approval.RefuseDigestMismatch},
+		{"P2 name", func(in *approval.VerbSignatureInput) { in.Approver.Name = "Mark" }, approval.ErrApproverName, 0},
+		{"P3 format", func(in *approval.VerbSignatureInput) { in.Approver.Format = approval.FormatWebAuthn }, approval.ErrApproverFormat, 0},
+		{"P4 revoked", func(in *approval.VerbSignatureInput) { in.Approver.RevokedAt = 950 }, approval.ErrApproverRevoked, 0},
+		{"P4 revoked before P6", func(in *approval.VerbSignatureInput) { in.Approver.RevokedAt, in.Signature = 950, []byte("x") }, approval.ErrApproverRevoked, 0},
+		{"P5 clock behind", func(in *approval.VerbSignatureInput) { in.Approver.EnrolledAt = vNow + 1 }, approval.ErrApproverClockBehind, 0},
+		{"P6 not armored", func(in *approval.VerbSignatureInput) { in.Signature = []byte("y\n") }, approval.ErrArmor, 0},
+		{"P6 another key", func(in *approval.VerbSignatureInput) { in.Approver.Key = other.PublicKey() }, approval.ErrKeyMismatch, 0},
+		{"P6 another digest", func(in *approval.VerbSignatureInput) {
+			d := vD
+			d[1] ^= 1
+			in.Admit.StoredDigest, in.Admit.BytesDigest = d, d
+		}, approval.ErrBadSignature, 0},
+		{"P6 another challenge", func(in *approval.VerbSignatureInput) { in.Challenge[1] ^= 1 }, approval.ErrBadSignature, 0},
+		{"P6 another expiry", func(in *approval.VerbSignatureInput) { in.Admit.ExpiresAt++ }, approval.ErrBadSignature, 0},
+	} {
+		in := vsInput(t, s)
+		tc.edit(&in)
+		v, err := approval.ConfirmBySignature(in)
+		var na *approval.NotAdmittedError
+		switch {
+		case v != nil:
+			t.Errorf("%s: confirmed", tc.name)
+		case tc.verdict != 0 && (!errors.As(err, &na) || na.Verdict != tc.verdict):
+			t.Errorf("%s: %v, want verdict %v", tc.name, err, tc.verdict)
+		case tc.want != nil && !errors.Is(err, tc.want):
+			t.Errorf("%s: %v, want %v", tc.name, err, tc.want)
+		}
+	}
+	v, err := approval.ConfirmBySignature(vsInput(t, s))
+	if err != nil || v.Binds(vD, vC, vExpiry) != nil || v.ApprovedBy() != "ssh:mark-ssh" {
+		t.Fatalf("control: (%v, %v)", v, err)
+	}
+	if name, format, key, ok := v.Use(); !ok || name != "mark-ssh" || format != approval.FormatSSH || !bytes.Equal(key, s.PublicKey().Blob) {
+		t.Fatalf("Use = %q %q %x %v", name, format, key, ok)
+	}
+	if _, _, ok := v.PasskeyUse(); ok {
+		t.Fatal("an ssh confirmation has a passkey use")
+	}
+}
+
+// VL6: the assertion law's order: P1, P2-P5 and P3w, W1-W14 over the verb message hash, then W15 in
+// full (a counter that does not advance is refused here); the control hands the consume the
+// credential ID and the assertion's counter.
+func TestConfirmByAssertionLaw(t *testing.T) {
+	for _, tc := range []struct {
+		name    string
+		edit    func(t *testing.T, in *approval.VerbAssertionInput)
+		want    error
+		verdict approval.Verdict
+	}{
+		{"P1 consumed", func(_ *testing.T, in *approval.VerbAssertionInput) { in.Admit.Consumed = true }, nil, approval.RefuseConsumed},
+		{"P1 clock behind", func(_ *testing.T, in *approval.VerbAssertionInput) { in.Admit.Now = vCreated - 1 }, nil, approval.RefuseClockBehind},
+		{"P1 expired", func(_ *testing.T, in *approval.VerbAssertionInput) { in.Admit.Now = vExpiry + 1 }, nil, approval.RefuseExpired},
+		{"P1 digest", func(_ *testing.T, in *approval.VerbAssertionInput) { in.Admit.BytesDigest[0] ^= 1 }, nil, approval.RefuseDigestMismatch},
+		{"P2 name", func(_ *testing.T, in *approval.VerbAssertionInput) { in.Approver.Name = "Mark" }, approval.ErrApproverName, 0},
+		{"P3w format", func(_ *testing.T, in *approval.VerbAssertionInput) { in.Approver.Format = approval.FormatSSH }, approval.ErrApproverNotPasskey, 0},
+		{"P3w no passkey", func(_ *testing.T, in *approval.VerbAssertionInput) { in.Approver.Passkey = nil }, approval.ErrApproverNotPasskey, 0},
+		{"P4 revoked", func(_ *testing.T, in *approval.VerbAssertionInput) { in.Approver.RevokedAt = 950 }, approval.ErrApproverRevoked, 0},
+		{"P5 clock behind", func(_ *testing.T, in *approval.VerbAssertionInput) { in.Approver.EnrolledAt = vNow + 1 }, approval.ErrApproverClockBehind, 0},
+		{"W1 envelope", func(_ *testing.T, in *approval.VerbAssertionInput) { in.Response = []byte("{}") }, approval.ErrPasskeyResponse, 0},
+		{"W2 raw id", func(t *testing.T, in *approval.VerbAssertionInput) {
+			in.Response = vPasskey.Assert(paOrigin, paRPID, vHash(t, vD, vC, vExpiry), 2, approvaltest.Edit{RawID: []byte("another")})
+		}, approval.ErrCredentialIDMismatch, 0},
+		{"W3 user handle", func(t *testing.T, in *approval.VerbAssertionInput) {
+			in.Response = vPasskey.Assert(paOrigin, paRPID, vHash(t, vD, vC, vExpiry), 2, approvaltest.Edit{UserHandle: []byte("mallory")})
+		}, approval.ErrUserHandle, 0},
+		{"W5 type", func(t *testing.T, in *approval.VerbAssertionInput) {
+			in.Response = vPasskey.Assert(paOrigin, paRPID, vHash(t, vD, vC, vExpiry), 2, approvaltest.Edit{Type: "webauthn.create"})
+		}, approval.ErrClientDataType, 0},
+		{"W6 another digest", func(_ *testing.T, in *approval.VerbAssertionInput) {
+			d := vD
+			d[1] ^= 1
+			in.Admit.StoredDigest, in.Admit.BytesDigest = d, d
+		}, approval.ErrChallenge, 0},
+		{"W6 another challenge", func(_ *testing.T, in *approval.VerbAssertionInput) { in.Challenge[1] ^= 1 }, approval.ErrChallenge, 0},
+		{"W6 another expiry", func(_ *testing.T, in *approval.VerbAssertionInput) { in.Admit.ExpiresAt++ }, approval.ErrChallenge, 0},
+		{"W7 origin", func(t *testing.T, in *approval.VerbAssertionInput) {
+			in.Response = vPasskey.Assert("https://y.ts.net", paRPID, vHash(t, vD, vC, vExpiry), 2)
+		}, approval.ErrOriginMismatch, 0},
+		{"W8 cross origin", func(t *testing.T, in *approval.VerbAssertionInput) {
+			in.Response = vPasskey.Assert(paOrigin, paRPID, vHash(t, vD, vC, vExpiry), 2, approvaltest.Edit{CrossOrigin: true})
+		}, approval.ErrCrossOrigin, 0},
+		{"W9 flags", func(t *testing.T, in *approval.VerbAssertionInput) {
+			in.Response = vPasskey.Assert(paOrigin, paRPID, vHash(t, vD, vC, vExpiry), 2, approvaltest.Edit{Set: 0x40})
+		}, approval.ErrAuthDataFlags, 0},
+		{"W10 rp id", func(t *testing.T, in *approval.VerbAssertionInput) {
+			in.Response = vPasskey.Assert(paOrigin, "y.ts.net", vHash(t, vD, vC, vExpiry), 2)
+		}, approval.ErrRPIDHash, 0},
+		{"W11 user presence", func(t *testing.T, in *approval.VerbAssertionInput) {
+			in.Response = vPasskey.Assert(paOrigin, paRPID, vHash(t, vD, vC, vExpiry), 2, approvaltest.Edit{Clear: 0x01})
+		}, approval.ErrUserPresence, 0},
+		{"W12 user verification", func(t *testing.T, in *approval.VerbAssertionInput) {
+			in.Response = vPasskey.Assert(paOrigin, paRPID, vHash(t, vD, vC, vExpiry), 2, approvaltest.Edit{Clear: 0x04})
+		}, approval.ErrUserVerification, 0},
+		{"W13 backup eligibility", func(t *testing.T, in *approval.VerbAssertionInput) {
+			in.Response = vPasskey.Assert(paOrigin, paRPID, vHash(t, vD, vC, vExpiry), 2, approvaltest.Edit{Set: 0x08})
+		}, approval.ErrBackupEligibility, 0},
+		{"W14 bad signature", func(t *testing.T, in *approval.VerbAssertionInput) {
+			in.Response = vPasskey.Assert(paOrigin, paRPID, vHash(t, vD, vC, vExpiry), 2, approvaltest.Edit{BadSignature: true})
+		}, approval.ErrBadSignature, 0},
+		{"W15 the stored counter", func(t *testing.T, in *approval.VerbAssertionInput) {
+			in.Response = vPasskey.Assert(paOrigin, paRPID, vHash(t, vD, vC, vExpiry), 1)
+		}, approval.ErrSignCountRegression, 0},
+		{"W15 a zero counter over a stored 1", func(t *testing.T, in *approval.VerbAssertionInput) {
+			in.Response = vPasskey.Assert(paOrigin, paRPID, vHash(t, vD, vC, vExpiry), 0)
+		}, approval.ErrSignCountRegression, 0},
+	} {
+		in := vaInput(t, 2)
+		tc.edit(t, &in)
+		v, err := approval.ConfirmByAssertion(in)
+		var na *approval.NotAdmittedError
+		switch {
+		case v != nil:
+			t.Errorf("%s: confirmed", tc.name)
+		case tc.verdict != 0 && (!errors.As(err, &na) || na.Verdict != tc.verdict):
+			t.Errorf("%s: %v, want verdict %v", tc.name, err, tc.verdict)
+		case tc.want != nil && !errors.Is(err, tc.want):
+			t.Errorf("%s: %v, want %v", tc.name, err, tc.want)
+		}
+	}
+	in := vaInput(t, 2)
+	v, err := approval.ConfirmByAssertion(in)
+	if err != nil || v.Binds(vD, vC, vExpiry) != nil || v.ApprovedBy() != "webauthn:mark-mac" {
+		t.Fatalf("control: (%v, %v)", v, err)
+	}
+	id, n, ok := v.PasskeyUse()
+	if !ok || n != 2 || !bytes.Equal(id, vPasskey.CredentialID()) {
+		t.Fatalf("PasskeyUse = %x %d %v", id, n, ok)
+	}
+	if _, format, key, ok := v.Use(); !ok || format != approval.FormatWebAuthn || !bytes.Equal(key, in.Approver.Passkey.Key.Canonical()) {
+		t.Fatalf("Use = %s %x %v", format, key, ok)
+	}
+	zero := vaInput(t, 0) // a counter-less authenticator (stored 0, asserted 0) passes W15
+	zero.Approver.Passkey.SignCount = 0
+	if _, err := approval.ConfirmByAssertion(zero); err != nil {
+		t.Fatalf("a counter-less authenticator: %v", err)
+	}
+}
+
+// VL7: Binds refuses each field that differs, the nil and the zero value; ApprovedBy, Use and
+// PasskeyUse are empty unless verified.
+func TestSignedVerbBinds(t *testing.T) {
+	v, err := approval.ConfirmBySignature(vsInput(t, approvaltest.NewSigner(75)))
+	if err != nil {
+		t.Fatal(err)
+	}
+	d, c := vD, vC
+	d[2], c[2] = 9, 9
+	for field, err := range map[string]error{
+		"digest":    v.Binds(d, vC, vExpiry),
+		"challenge": v.Binds(vD, c, vExpiry),
+		"expiry":    v.Binds(vD, vC, vExpiry+1),
+	} {
+		var be *approval.BindingError
+		if !errors.As(err, &be) || be.Field != field {
+			t.Errorf("%s: %v", field, err)
+		}
+	}
+	if err := v.Binds(vD, vC, vExpiry); err != nil {
+		t.Fatalf("control: %v", err)
+	}
+	var none *approval.SignedVerb
+	zero := &approval.SignedVerb{}
+	for name, s := range map[string]*approval.SignedVerb{"nil": none, "zero": zero} {
+		if err := s.Binds([32]byte{}, [32]byte{}, 0); !errors.Is(err, approval.ErrNotVerified) {
+			t.Errorf("%s: Binds = %v", name, err)
+		}
+		if s.ApprovedBy() != "" {
+			t.Errorf("%s: ApprovedBy = %q", name, s.ApprovedBy())
+		}
+		if _, _, _, ok := s.Use(); ok {
+			t.Errorf("%s: Use ok", name)
+		}
+		if _, _, ok := s.PasskeyUse(); ok {
+			t.Errorf("%s: PasskeyUse ok", name)
+		}
+	}
+	_, _, key, _ := v.Use()
+	key[0] ^= 1
+	if _, _, again, _ := v.Use(); again[0] == key[0] {
+		t.Fatal("Use hands out its own key")
+	}
+}
````

### Patch M1c — `CheckConfirmation` (design §2.3 [r1]); VL8, VL4 (recorded Chrome)

````diff M1c
diff --git a/host/approval/verb_check_test.go b/host/approval/verb_check_test.go
new file mode 100644
index 0000000..7511a9f
--- /dev/null
+++ b/host/approval/verb_check_test.go
@@ -0,0 +1,110 @@
+package approval_test
+
+// Row 243 V3 sprint 1 M1c (design §2.3 [r1], §6 VL8): CheckConfirmation, the CLI's re-check of a
+// recorded confirmation. It verifies an artifact over the caller's own digest, challenge and expiry
+// with no P1 (the row is consumed and may be past its expiry) and no W15 (the daemon moved the
+// counter to this very assertion's); it judges the approver at now, so a key revoked since the
+// confirmation is refused (fail-safe; audit note 5, carried to sprint 5).
+
+import (
+	"bytes"
+	"errors"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
+)
+
+// vcSSH is a recorded SSHSIG confirmation of (vD, vC, vExpiry) re-checked at now: after the
+// expiry, as the CLI's poll may run (the 10 s grace).
+func vcSSH(t *testing.T, s *approvaltest.Signer) approval.VerbCheckInput {
+	t.Helper()
+	return approval.VerbCheckInput{Digest: vD, Challenge: vC, Expiry: vExpiry, Now: vExpiry + 5, Format: approval.FormatSSH,
+		Artifact: s.Sign(approval.VerbNamespace, vMsg(t, vD, vC, vExpiry)),
+		Approver: approval.Approver{Name: "mark-ssh", Format: approval.FormatSSH, Key: s.PublicKey(), EnrolledAt: 900}}
+}
+
+// vcPasskey is a recorded assertion (counter 2) re-checked after the daemon moved the stored
+// counter to 2.
+func vcPasskey(t *testing.T) approval.VerbCheckInput {
+	t.Helper()
+	a := paApprover(t, vPasskey)
+	a.Passkey.SignCount = 2
+	return approval.VerbCheckInput{Digest: vD, Challenge: vC, Expiry: vExpiry, Now: vExpiry + 5, Format: approval.FormatWebAuthn,
+		Artifact: vPasskey.Assert(paOrigin, paRPID, vHash(t, vD, vC, vExpiry), 2), Approver: a}
+}
+
+func TestCheckConfirmation(t *testing.T) {
+	s := approvaltest.NewSigner(76)
+	mintMsg, _ := approval.SignedMessage(vD, vC, vExpiry)
+	mintHash, _ := approval.MessageHash(vD, vC, vExpiry)
+	other := vD
+	other[3] ^= 1
+	for _, tc := range []struct {
+		name string
+		in   func(t *testing.T) approval.VerbCheckInput
+		edit func(in *approval.VerbCheckInput)
+		want error
+	}{
+		{"ssh revoked approver", func(t *testing.T) approval.VerbCheckInput { return vcSSH(t, s) }, func(in *approval.VerbCheckInput) { in.Approver.RevokedAt = vExpiry + 1 }, approval.ErrApproverRevoked},
+		{"ssh approver name", func(t *testing.T) approval.VerbCheckInput { return vcSSH(t, s) }, func(in *approval.VerbCheckInput) { in.Approver.Name = "Mark" }, approval.ErrApproverName},
+		{"ssh enrolled after now", func(t *testing.T) approval.VerbCheckInput { return vcSSH(t, s) }, func(in *approval.VerbCheckInput) { in.Approver.EnrolledAt = vExpiry + 6 }, approval.ErrApproverClockBehind},
+		{"ssh a webauthn approver", func(t *testing.T) approval.VerbCheckInput { return vcSSH(t, s) }, func(in *approval.VerbCheckInput) { in.Approver.Format = approval.FormatWebAuthn }, approval.ErrApproverFormat},
+		{"ssh another digest", func(t *testing.T) approval.VerbCheckInput { return vcSSH(t, s) }, func(in *approval.VerbCheckInput) { in.Digest = other }, approval.ErrBadSignature},
+		{"ssh another challenge", func(t *testing.T) approval.VerbCheckInput { return vcSSH(t, s) }, func(in *approval.VerbCheckInput) { in.Challenge = other }, approval.ErrBadSignature},
+		{"ssh another expiry", func(t *testing.T) approval.VerbCheckInput { return vcSSH(t, s) }, func(in *approval.VerbCheckInput) { in.Expiry-- }, approval.ErrBadSignature},
+		{"ssh a mint-namespace SSHSIG", func(t *testing.T) approval.VerbCheckInput { return vcSSH(t, s) }, func(in *approval.VerbCheckInput) { in.Artifact = s.Sign(approval.Namespace, mintMsg) }, approval.ErrNamespace},
+		{"ssh a decision-namespace SSHSIG", func(t *testing.T) approval.VerbCheckInput { return vcSSH(t, s) }, func(in *approval.VerbCheckInput) {
+			in.Artifact = s.Sign(approval.DecisionNamespace, vMsg(t, vD, vC, vExpiry))
+		}, approval.ErrNamespace},
+		{"ssh a negative expiry", func(t *testing.T) approval.VerbCheckInput { return vcSSH(t, s) }, func(in *approval.VerbCheckInput) { in.Expiry = -1 }, approval.ErrNegativeExpiry},
+		{"ssh an assertion as an SSHSIG", func(t *testing.T) approval.VerbCheckInput { return vcSSH(t, s) }, func(in *approval.VerbCheckInput) { in.Artifact = vcPasskey(t).Artifact }, approval.ErrArmor},
+		{"ssh another key", func(t *testing.T) approval.VerbCheckInput { return vcSSH(t, s) }, func(in *approval.VerbCheckInput) { in.Approver.Key = approvaltest.NewSigner(77).PublicKey() }, approval.ErrKeyMismatch},
+		{"webauthn revoked approver", vcPasskey, func(in *approval.VerbCheckInput) { in.Approver.RevokedAt = vExpiry + 1 }, approval.ErrApproverRevoked},
+		{"webauthn an ssh approver", vcPasskey, func(in *approval.VerbCheckInput) { in.Approver.Format = approval.FormatSSH }, approval.ErrApproverNotPasskey},
+		{"webauthn another digest", vcPasskey, func(in *approval.VerbCheckInput) { in.Digest = other }, approval.ErrChallenge},
+		{"webauthn another challenge", vcPasskey, func(in *approval.VerbCheckInput) { in.Challenge = other }, approval.ErrChallenge},
+		{"webauthn another expiry", vcPasskey, func(in *approval.VerbCheckInput) { in.Expiry++ }, approval.ErrChallenge},
+		{"webauthn a negative expiry", vcPasskey, func(in *approval.VerbCheckInput) { in.Expiry = -1 }, approval.ErrChallenge},
+		{"webauthn a mint assertion", vcPasskey, func(in *approval.VerbCheckInput) { in.Artifact = vPasskey.Assert(paOrigin, paRPID, mintHash, 2) }, approval.ErrChallenge},
+		{"webauthn bad signature", vcPasskey, func(in *approval.VerbCheckInput) {
+			in.Artifact = vPasskey.Assert(paOrigin, paRPID, vHash(t, vD, vC, vExpiry), 2, approvaltest.Edit{BadSignature: true})
+		}, approval.ErrBadSignature},
+		{"no format", vcPasskey, func(in *approval.VerbCheckInput) { in.Format = "" }, approval.ErrVerbFormat},
+		{"an unknown format", func(t *testing.T) approval.VerbCheckInput { return vcSSH(t, s) }, func(in *approval.VerbCheckInput) { in.Format = "typed" }, approval.ErrVerbFormat},
+	} {
+		in := tc.in(t)
+		tc.edit(&in)
+		if v, err := approval.CheckConfirmation(in); v != nil || !errors.Is(err, tc.want) {
+			t.Errorf("%s: (%v, %v), want %v", tc.name, v, err, tc.want)
+		}
+	}
+	// The verified re-checks: no P1 (past the expiry, a consumed row has no admission input at all)
+	// and no W15 (the stored counter is already this assertion's, and a counter-less 0 over 0).
+	for name, in := range map[string]approval.VerbCheckInput{"ssh": vcSSH(t, s), "webauthn, counter moved to 2": vcPasskey(t)} {
+		v, err := approval.CheckConfirmation(in)
+		if err != nil || v.Binds(vD, vC, vExpiry) != nil {
+			t.Fatalf("%s: (%v, %v)", name, v, err)
+		}
+		if want := in.Format + ":" + in.Approver.Name; v.ApprovedBy() != want {
+			t.Fatalf("%s: ApprovedBy = %q, want %q", name, v.ApprovedBy(), want)
+		}
+	}
+	ahead := vcPasskey(t)
+	ahead.Approver.Passkey.SignCount = 7 // a later assertion moved the counter on; the recorded one still verifies
+	v, err := approval.CheckConfirmation(ahead)
+	if err != nil {
+		t.Fatalf("a counter already past the artifact's: %v", err)
+	}
+	id, n, ok := v.PasskeyUse()
+	if !ok || n != 2 || !bytes.Equal(id, vPasskey.CredentialID()) {
+		t.Fatalf("PasskeyUse = %x %d %v", id, n, ok)
+	}
+	if _, format, key, ok := v.Use(); !ok || format != approval.FormatWebAuthn || !bytes.Equal(key, ahead.Approver.Passkey.Key.Canonical()) {
+		t.Fatalf("Use = %s %x %v", format, key, ok)
+	}
+	ssh, _ := approval.CheckConfirmation(vcSSH(t, s))
+	if _, _, key, ok := ssh.Use(); !ok || !bytes.Equal(key, s.PublicKey().Blob) {
+		t.Fatalf("ssh Use = %x %v", key, ok)
+	}
+}
diff --git a/host/approval/verb_chrome_test.go b/host/approval/verb_chrome_test.go
new file mode 100644
index 0000000..f0d3322
--- /dev/null
+++ b/host/approval/verb_chrome_test.go
@@ -0,0 +1,28 @@
+package approval
+
+// Row 243 V3 sprint 1 M1c (design §6 VL4, the recorded half): Chrome's recorded assertions (A2's
+// fixtures, testdata/webauthn) signed a mint request's message hash; the verb laws refuse each at
+// W6, while the mint law (the control) approves it.
+
+import (
+	"errors"
+	"testing"
+)
+
+func TestChromeMintAssertionIsNotAVerb(t *testing.T) {
+	for _, alg := range []string{"es256", "eddsa"} {
+		in := chromeAssertionInput(t, alg)
+		if _, err := ApproveByAssertion(in); err != nil {
+			t.Fatalf("%s: control: the recorded assertion does not mint: %v", alg, err)
+		}
+		v, err := ConfirmByAssertion(VerbAssertionInput{Admit: in.Admit, Challenge: in.Challenge, Response: in.Response, Approver: in.Approver})
+		if v != nil || !errors.Is(err, ErrChallenge) {
+			t.Fatalf("%s: ConfirmByAssertion = (%v, %v), want ErrChallenge", alg, v, err)
+		}
+		v, err = CheckConfirmation(VerbCheckInput{Digest: in.Admit.StoredDigest, Challenge: in.Challenge, Expiry: in.Admit.ExpiresAt, Now: in.Admit.Now,
+			Artifact: in.Response, Format: FormatWebAuthn, Approver: in.Approver})
+		if v != nil || !errors.Is(err, ErrChallenge) {
+			t.Fatalf("%s: CheckConfirmation = (%v, %v), want ErrChallenge", alg, v, err)
+		}
+	}
+}
diff --git a/host/approval/verb_signed.go b/host/approval/verb_signed.go
index 44cce54..6d9d7c2 100644
--- a/host/approval/verb_signed.go
+++ b/host/approval/verb_signed.go
@@ -3,16 +3,18 @@ package approval
 // Row 243 V3 sprint 1 (design §2.3): the pure law of a verb confirmation by an enrolled key, over the
 // 86-byte verb message (verb.go). ConfirmBySignature and ConfirmByAssertion are the mint laws'
 // shape (signed.go, passkey_assert.go) with the verb message: P1 admission, P2-P5 checkApprover,
-// then the signature, and for a passkey W1-W14 (assertionSteps) and W15 in full. No I/O; now is an
-// input.
+// then the signature, and for a passkey W1-W14 (assertionSteps) and W15 in full. CheckConfirmation
+// is the CLI's re-check of a recorded confirmation: the same law without P1 and W15. No I/O; now is
+// an input.
 
 import (
 	"bytes"
+	"errors"
 	"fmt"
 )
 
 // SignedVerb is a verified verb confirmation. Its fields are unexported, so only ConfirmBySignature,
-// and ConfirmByAssertion can return one with ok set; the zero value binds nothing.
+// ConfirmByAssertion and CheckConfirmation can return one with ok set; the zero value binds nothing.
 type SignedVerb struct {
 	ok           bool
 	digest       [32]byte
@@ -42,6 +44,16 @@ type VerbAssertionInput struct {
 	Approver  Approver
 }
 
+// VerbCheckInput is what CheckConfirmation re-checks: the CLI's own digest, challenge and expiry,
+// the clock, the recorded artifact, its format and the enrolled approver it names.
+type VerbCheckInput struct {
+	Digest, Challenge [32]byte
+	Expiry, Now       int64
+	Artifact          []byte
+	Format            string // FormatSSH or FormatWebAuthn
+	Approver          Approver
+}
+
 // verifyVerbSSHSIG parses an armored SSHSIG and verifies it over VerbMessage(digest, challenge,
 // expiry) under VerbNamespace with key.
 func verifyVerbSSHSIG(armored []byte, digest, challenge [32]byte, expiry int64, key SSHPublicKey) error {
@@ -98,6 +110,41 @@ func ConfirmByAssertion(in VerbAssertionInput) (*SignedVerb, error) {
 		name: in.Approver.Name, format: FormatWebAuthn, key: p.Key.Canonical(), credentialID: bytes.Clone(p.CredentialID), signCount: ad.SignCount}, nil
 }
 
+// ErrVerbFormat is a CheckConfirmation format other than FormatSSH and FormatWebAuthn.
+var ErrVerbFormat = errors.New("approval: a confirmation's format is ssh or webauthn")
+
+// CheckConfirmation is the CLI's re-check of a recorded confirmation (design §2.3 [r1], §2.7 step
+// 6(a)): the verb law without P1 (the row is consumed by then) and without W15 (the counter was
+// judged and moved by the daemon), over the caller's own digest, challenge and expiry, as
+// DecideByAssertion runs W1-W14 without W15. The approver's state (P2-P5) is judged at now, so a
+// key revoked since the confirmation is refused (fail-safe).
+func CheckConfirmation(in VerbCheckInput) (*SignedVerb, error) {
+	switch in.Format {
+	case FormatSSH:
+		if err := checkApprover(AdmitInput{Now: in.Now}, in.Approver, FormatSSH); err != nil { // P2-P5
+			return nil, err
+		}
+		if err := verifyVerbSSHSIG(in.Artifact, in.Digest, in.Challenge, in.Expiry, in.Approver.Key); err != nil { // P6
+			return nil, err
+		}
+		return &SignedVerb{ok: true, digest: in.Digest, challenge: in.Challenge, expiry: in.Expiry,
+			name: in.Approver.Name, format: FormatSSH, key: bytes.Clone(in.Approver.Key.Blob)}, nil
+	case FormatWebAuthn:
+		if err := checkApprover(AdmitInput{Now: in.Now}, in.Approver, FormatWebAuthn); err != nil { // P2-P5, P3w
+			return nil, err
+		}
+		h, hErr := VerbMessageHash(in.Digest, in.Challenge, in.Expiry)
+		ad, err := assertionSteps(in.Artifact, in.Approver, h, hErr) // W1-W14; no W15
+		if err != nil {
+			return nil, err
+		}
+		p := in.Approver.Passkey
+		return &SignedVerb{ok: true, digest: in.Digest, challenge: in.Challenge, expiry: in.Expiry,
+			name: in.Approver.Name, format: FormatWebAuthn, key: p.Key.Canonical(), credentialID: bytes.Clone(p.CredentialID), signCount: ad.SignCount}, nil
+	}
+	return nil, fmt.Errorf("%w: %q", ErrVerbFormat, in.Format)
+}
+
 // Binds reports whether s is a verified confirmation of exactly this digest, challenge and expiry.
 func (s *SignedVerb) Binds(digest, challenge [32]byte, expiry int64) error {
 	switch {
````

### Patch M3 — the finish prompt's text, moved (design §2.2); VT1 + the base-generated golden

````diff M3
diff --git a/cmd/ailang-worldd/committer.go b/cmd/ailang-worldd/committer.go
index 27ab50b..9f82021 100644
--- a/cmd/ailang-worldd/committer.go
+++ b/cmd/ailang-worldd/committer.go
@@ -1,8 +1,6 @@
 package main
 
-func displayCommitter(value string) string {
-	if value == "" {
-		return "(repository git config)"
-	}
-	return value
-}
+import "github.com/sunholo-data/ailang-world/host/broker"
+
+// displayCommitter is broker.DisplayCommitter (row 243 V3: one rule for the terminal and the card).
+func displayCommitter(value string) string { return broker.DisplayCommitter(value) }
diff --git a/cmd/ailang-worldd/episode_finish_operator.go b/cmd/ailang-worldd/episode_finish_operator.go
index 1fcfc93..f5fcb13 100644
--- a/cmd/ailang-worldd/episode_finish_operator.go
+++ b/cmd/ailang-worldd/episode_finish_operator.go
@@ -12,7 +12,6 @@ import (
 	"net/http"
 	"os"
 	"path/filepath"
-	"strings"
 	"time"
 
 	"github.com/sunholo-data/ailang-world/host/broker"
@@ -156,17 +155,5 @@ func finishPrompt(tty io.Writer, scope broker.FinishApprovalScope, title, body s
 	if len(committer) > 0 {
 		identity = committer[0]
 	}
-	fmt.Fprintf(tty, "committer: %s\n", displayCommitter(identity))
-	lines := strings.Split(body, "\n")
-	preview := strings.Join(lines[:min(6, len(lines))], "\n")
-	state := "draft"
-	if !scope.Draft {
-		state = "ready"
-	}
-	label := "ailang-world"
-	if scope.CreateLabel {
-		label += " (will be created)"
-	}
-	fmt.Fprintf(tty, "repository: %s\nbranch: %s\nsha: %s\nbase: %s\ntitle: %s\nbody (%d bytes), first lines:\n%s\nlabel: %s\n%s\nprotected paths changed: %s\nOpen this pull request? [y/N] ",
-		scope.Repository, scope.Head, scope.SHA, scope.Base, title, len(body), preview, label, state, strings.Join(protected, ", "))
+	fmt.Fprint(tty, broker.FinishPromptText(scope, title, body, protected, identity)+"Open this pull request? [y/N] ")
 }
diff --git a/cmd/ailang-worldd/finish_prompt_golden_test.go b/cmd/ailang-worldd/finish_prompt_golden_test.go
new file mode 100644
index 0000000..245391f
--- /dev/null
+++ b/cmd/ailang-worldd/finish_prompt_golden_test.go
@@ -0,0 +1,87 @@
+package main
+
+// Row 243 V3 sprint 1 M3 (design §2.2, §6 VT1): the bytes finishPrompt writes at the terminal, for
+// three fixtures (draft and ready, create-label, protected paths, an empty committer), pinned to the
+// golden that the base tree's finishPrompt generated before the text moved to
+// broker.FinishPromptText (WORLD_UPDATE_FINISH_PROMPT_GOLDEN=1 on the base tree, plan §1.2). The
+// broker's TestFinishPromptText reads the same golden; displayCommitter (episode start's preset
+// prompt) is pinned to the same "" rule.
+
+import (
+	"bytes"
+	"encoding/json"
+	"os"
+	"reflect"
+	"strings"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/broker"
+)
+
+const finishPromptGolden = "../../host/broker/testdata/finish-prompt.golden.json"
+
+// finishPromptFixture is one golden row: finishPrompt's inputs and the bytes it wrote.
+type finishPromptFixture struct {
+	Name      string
+	Scope     broker.FinishApprovalScope
+	Title     string
+	Body      string
+	Protected []string
+	Committer string
+	Prompt    string
+}
+
+func finishPromptFixtures() []finishPromptFixture {
+	scope := broker.FinishApprovalScope{Repository: "sunholo-data/ailang-world", Remote: "origin", Head: "world/w1",
+		SHA: strings.Repeat("a1", 20), Base: "dev", TitleSHA256: strings.Repeat("b2", 32), BodySHA256: strings.Repeat("c3", 32),
+		Draft: true, ExpiresAt: 1700000900}
+	ready := scope
+	ready.Draft, ready.CreateLabel, ready.Head, ready.Base = false, true, "world/w2", "main"
+	label := scope
+	label.CreateLabel = true
+	return []finishPromptFixture{
+		{Name: "draft, existing label, no protected path, repository committer", Scope: scope, Title: "Row 243: the law",
+			Body: "## Episode w1\n\nline 3\nline 4\nline 5\nline 6\nline 7 (not shown)\nline 8\n"},
+		{Name: "ready, label created, two protected paths, pinned committer", Scope: ready, Title: "Tidy <b> & \"quote\"",
+			Body: "one\ntwo\nthree", Protected: []string{"cmd/ailang-worldd/main.go", "go.mod"}, Committer: "World Bot <bot@example.invalid>"},
+		{Name: "draft, label created, one protected path, empty body", Scope: label, Title: "Café — naïve",
+			Body: "", Protected: []string{"docs/QUICKSTART.md"}, Committer: ""},
+	}
+}
+
+func TestFinishPromptIsTheGolden(t *testing.T) {
+	got := finishPromptFixtures()
+	for i := range got {
+		var tty bytes.Buffer
+		finishPrompt(&tty, got[i].Scope, got[i].Title, got[i].Body, got[i].Protected, got[i].Committer)
+		got[i].Prompt = tty.String()
+	}
+	if os.Getenv("WORLD_UPDATE_FINISH_PROMPT_GOLDEN") == "1" {
+		b, err := json.MarshalIndent(got, "", "  ")
+		if err != nil {
+			t.Fatal(err)
+		}
+		if err := os.WriteFile(finishPromptGolden, append(b, '\n'), 0o644); err != nil {
+			t.Fatal(err)
+		}
+	}
+	raw, err := os.ReadFile(finishPromptGolden)
+	if err != nil {
+		t.Fatal(err)
+	}
+	var want []finishPromptFixture
+	if err := json.Unmarshal(raw, &want); err != nil {
+		t.Fatal(err)
+	}
+	if len(want) != 3 || !reflect.DeepEqual(got, want) {
+		t.Fatalf("finishPrompt's bytes moved from the base's:\n got %+v\nwant %+v", got, want)
+	}
+	var none bytes.Buffer // the variadic committer omitted is the empty committer
+	finishPrompt(&none, want[0].Scope, want[0].Title, want[0].Body, want[0].Protected)
+	if none.String() != want[0].Prompt {
+		t.Fatalf("no committer argument: %q", none.String())
+	}
+	if got := displayCommitter(""); got != "(repository git config)" { // episode start's preset prompt shares the rule
+		t.Fatalf("displayCommitter(\"\") = %q", got)
+	}
+}
diff --git a/host/broker/finish_prompt.go b/host/broker/finish_prompt.go
new file mode 100644
index 0000000..a22f888
--- /dev/null
+++ b/host/broker/finish_prompt.go
@@ -0,0 +1,38 @@
+package broker
+
+// Row 243 V3 sprint 1 M3 (design §2.2): the finish prompt's text as a pure function of the approval
+// scope, the title, the body, the protected paths and the committer, moved verbatim from
+// cmd/ailang-worldd finishPrompt (which now prints FinishPromptText and its question), so the
+// terminal and a stored verb request carry the same bytes.
+
+import (
+	"fmt"
+	"strings"
+)
+
+// DisplayCommitter is how a prompt shows a pinned committer: "" is the repository's git config.
+func DisplayCommitter(v string) string {
+	if v == "" {
+		return "(repository git config)"
+	}
+	return v
+}
+
+// FinishPromptText is the text `episode finish` prints before "Open this pull request? [y/N] ":
+// the committer, the scope's repository, branch, sha and base, the title, the body's length and
+// first six lines, the label, draft or ready, and the protected paths changed.
+func FinishPromptText(s FinishApprovalScope, title, body string, protected []string, committer string) string {
+	lines := strings.Split(body, "\n")
+	preview := strings.Join(lines[:min(6, len(lines))], "\n")
+	state := "draft"
+	if !s.Draft {
+		state = "ready"
+	}
+	label := "ailang-world"
+	if s.CreateLabel {
+		label += " (will be created)"
+	}
+	return fmt.Sprintf("committer: %s\n", DisplayCommitter(committer)) +
+		fmt.Sprintf("repository: %s\nbranch: %s\nsha: %s\nbase: %s\ntitle: %s\nbody (%d bytes), first lines:\n%s\nlabel: %s\n%s\nprotected paths changed: %s\n",
+			s.Repository, s.Head, s.SHA, s.Base, title, len(body), preview, label, state, strings.Join(protected, ", "))
+}
diff --git a/host/broker/finish_prompt_test.go b/host/broker/finish_prompt_test.go
new file mode 100644
index 0000000..0132f26
--- /dev/null
+++ b/host/broker/finish_prompt_test.go
@@ -0,0 +1,49 @@
+package broker
+
+// Row 243 V3 sprint 1 M3 (design §2.2, §6 VT1): FinishPromptText plus the question is, byte for
+// byte, what the base tree's cmd/ailang-worldd finishPrompt wrote for the three golden fixtures
+// (testdata/finish-prompt.golden.json, generated by that finishPrompt; cmd's
+// TestFinishPromptIsTheGolden pins the delegation to the same file).
+
+import (
+	"encoding/json"
+	"os"
+	"path/filepath"
+	"testing"
+)
+
+func TestFinishPromptText(t *testing.T) {
+	raw, err := os.ReadFile(filepath.Join("testdata", "finish-prompt.golden.json"))
+	if err != nil {
+		t.Fatal(err)
+	}
+	var fixtures []struct {
+		Name              string
+		Scope             FinishApprovalScope
+		Title, Body       string
+		Protected         []string
+		Committer, Prompt string
+	}
+	if err := json.Unmarshal(raw, &fixtures); err != nil {
+		t.Fatal(err)
+	}
+	var draft, ready, created, protected, repoCommitter bool
+	for _, f := range fixtures {
+		if got := FinishPromptText(f.Scope, f.Title, f.Body, f.Protected, f.Committer) + "Open this pull request? [y/N] "; got != f.Prompt {
+			t.Errorf("%s:\n got %q\nwant %q", f.Name, got, f.Prompt)
+		}
+		draft, ready = draft || f.Scope.Draft, ready || !f.Scope.Draft
+		created, protected, repoCommitter = created || f.Scope.CreateLabel, protected || len(f.Protected) > 1, repoCommitter || f.Committer == ""
+	}
+	if len(fixtures) != 3 || !draft || !ready || !created || !protected || !repoCommitter {
+		t.Fatalf("the golden no longer covers draft, ready, a created label, protected paths and the repository committer: %d fixtures", len(fixtures))
+	}
+}
+
+func TestDisplayCommitter(t *testing.T) {
+	for in, want := range map[string]string{"": "(repository git config)", "World Bot <bot@example.invalid>": "World Bot <bot@example.invalid>", " ": " "} {
+		if got := DisplayCommitter(in); got != want {
+			t.Errorf("DisplayCommitter(%q) = %q, want %q", in, got, want)
+		}
+	}
+}
diff --git a/host/broker/testdata/finish-prompt.golden.json b/host/broker/testdata/finish-prompt.golden.json
new file mode 100644
index 0000000..ed862d9
--- /dev/null
+++ b/host/broker/testdata/finish-prompt.golden.json
@@ -0,0 +1,67 @@
+[
+  {
+    "Name": "draft, existing label, no protected path, repository committer",
+    "Scope": {
+      "Repository": "sunholo-data/ailang-world",
+      "Remote": "origin",
+      "Head": "world/w1",
+      "SHA": "a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1",
+      "Base": "dev",
+      "TitleSHA256": "b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2",
+      "BodySHA256": "c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3",
+      "Draft": true,
+      "CreateLabel": false,
+      "ExpiresAt": 1700000900
+    },
+    "Title": "Row 243: the law",
+    "Body": "## Episode w1\n\nline 3\nline 4\nline 5\nline 6\nline 7 (not shown)\nline 8\n",
+    "Protected": null,
+    "Committer": "",
+    "Prompt": "committer: (repository git config)\nrepository: sunholo-data/ailang-world\nbranch: world/w1\nsha: a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1\nbase: dev\ntitle: Row 243: the law\nbody (69 bytes), first lines:\n## Episode w1\n\nline 3\nline 4\nline 5\nline 6\nlabel: ailang-world\ndraft\nprotected paths changed: \nOpen this pull request? [y/N] "
+  },
+  {
+    "Name": "ready, label created, two protected paths, pinned committer",
+    "Scope": {
+      "Repository": "sunholo-data/ailang-world",
+      "Remote": "origin",
+      "Head": "world/w2",
+      "SHA": "a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1",
+      "Base": "main",
+      "TitleSHA256": "b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2",
+      "BodySHA256": "c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3",
+      "Draft": false,
+      "CreateLabel": true,
+      "ExpiresAt": 1700000900
+    },
+    "Title": "Tidy \u003cb\u003e \u0026 \"quote\"",
+    "Body": "one\ntwo\nthree",
+    "Protected": [
+      "cmd/ailang-worldd/main.go",
+      "go.mod"
+    ],
+    "Committer": "World Bot \u003cbot@example.invalid\u003e",
+    "Prompt": "committer: World Bot \u003cbot@example.invalid\u003e\nrepository: sunholo-data/ailang-world\nbranch: world/w2\nsha: a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1\nbase: main\ntitle: Tidy \u003cb\u003e \u0026 \"quote\"\nbody (13 bytes), first lines:\none\ntwo\nthree\nlabel: ailang-world (will be created)\nready\nprotected paths changed: cmd/ailang-worldd/main.go, go.mod\nOpen this pull request? [y/N] "
+  },
+  {
+    "Name": "draft, label created, one protected path, empty body",
+    "Scope": {
+      "Repository": "sunholo-data/ailang-world",
+      "Remote": "origin",
+      "Head": "world/w1",
+      "SHA": "a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1",
+      "Base": "dev",
+      "TitleSHA256": "b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2b2",
+      "BodySHA256": "c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3c3",
+      "Draft": true,
+      "CreateLabel": true,
+      "ExpiresAt": 1700000900
+    },
+    "Title": "Café — naïve",
+    "Body": "",
+    "Protected": [
+      "docs/QUICKSTART.md"
+    ],
+    "Committer": "",
+    "Prompt": "committer: (repository git config)\nrepository: sunholo-data/ailang-world\nbranch: world/w1\nsha: a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1a1\nbase: dev\ntitle: Café — naïve\nbody (0 bytes), first lines:\n\nlabel: ailang-world (will be created)\ndraft\nprotected paths changed: docs/QUICKSTART.md\nOpen this pull request? [y/N] "
+  }
+]
````

