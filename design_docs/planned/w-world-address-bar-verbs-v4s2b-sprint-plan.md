# Sprint plan — w-world-address-bar-verbs-v4, sprint S2b "the host/authority operator judges" (iteration 409, row 243)

> **THE SPRINT — row 243 phase V4, sprint S2b (design §3.1 [r1] "NEW typed verified-action proof can originate only in these
> judges", §5.2 ConfirmBothKeys and the store half of RunnerRechecks, §6 placement "adapters in `host/authority`"; S2 plan
> Appendix C "Moved out (S2-D1)"; judge-408 notes 1–4); D-WORLD-119 and D-WORLD-120 OPEN — S2b is answer-independent (S2
> Appendix C: the adapters "need no schema"; D-WORLD-119's default "build S1 and S2 against fixtures only"):**
> **M1** — NEW `host/authority/judge_operator.go`: the parked adapters, landed and tightened — `JudgeOperatorAction` (the daemon's
> confirmation judge), `RecheckOperatorAction` (the runner's recheck, with the recorded-approver comparison `ErrOperatorApprover`),
> `OperatorJudgeInput` (with an `admit()` mirroring `VerbJudgeInput.admit`), and the judge-only proof `VerifiedOperatorAction`
> (unexported fields; `Binds`, `ApprovedBy`, `PasskeyUse`, `Rechecked`). **Neither judge takes an approver**: each reads the
> enrolled approver the artifact names from the read-only store (V3's `sshApprover`/`passkeyApprover`, reused unchanged), so
> the law judges the store's key, format, enrolment, revocation and counter. Tests: `TestOperatorJudgeBothKeys`,
> **`TestOperatorJudgeReadsTheStore`** (NEW: 19 rows, each a store edit or a row edit, each asserting its `*JudgeError` step, its
> error class and **no proof**), `TestOperatorRunnerRechecks` (rewritten: 2 controls + 9 refusals), **`TestOperatorUnjudgedProofBindsNothing`**
> (NEW) and `TestVerifiedOperatorActionOnlyFromTheJudges` (AST + reflect, now also "no judge takes an approver" with a positive
> control). **M2** — judge-408 notes 2 and 3 in host/approval (test only): `TestOperatorPasskeyUseIsACopy` kills **`X-PK-ALIAS`**
> (reproduced SURVIVED at the base, whole package) and its accessor twin `PK-USE-ALIAS`; `TestOperatorRunnerRecheckBoundsExpiryNotLifetime`
> pins the runner recheck's expiry-not-lifetime behaviour with design §3.1's reason and kills `R1-LIFETIME-ADDED`. **M3** — the
> two judges write nothing (`TestOperatorJudgesWriteNothing`, judge-407's rule; its mutant control `WN-LOG-INJECT` killed).
> **No schema change, no daemon route, no `.ail`, `surfaceCSP` and `app.js` unchanged, no served or CLI behaviour change; the
> judges have no production caller; no pre-existing file changes.**
> - **Measured: 526 changed lines** — `git diff -M --numstat 4b1c39e2 HEAD -- cmd host docs app` (insertions + deletions, renames on;
>   no rename occurs; **the same 526 with no path filter**; nothing lands outside `host/`), after a commit of each of the three
>   scratch stages: **116 production** (`judge_operator.go`), **410 test**. **Cap 1,400: margin 874. No cut.**
> - **Three sub-milestones, each green at its commit** (§1.2 C14): **M1** the judges (401) → **M2** the judge-408 pins (469) →
>   **M3** writes-nothing (526). The controller's `ctl_roundtrip.sh` loops `M1 M2 M3 M4`: **there is no M4** — edit its loop to
>   `M1 M2 M3` (it prints `EMPTY M4` and exits 1 otherwise, after M1–M3 have applied and diffed).
> - **Every line ships as a compiled patch** (Appendix P), extracted from this file and round-tripped onto a fresh
>   `git archive 4b1c39e2` export, byte-identical to the stage trees (§1.2 C16). The executor applies, never retypes.
> - **Pre-existing files changed: 0; deleted: 0** (production or test). Every S2b path is added.
> - **Mutants: 61** (enumerated from the diff, plus the store-sourced checks the judges reach in `judge_verb.go`): **59 killed by
>   their named test and by the whole package, 2 named equivalents, 0 survivors, 0 NM**; no kill rests on a static test alone.
>   The parked patch as iteration 408 shipped it, under the same mutants (`park0`): **13 whole-package survivors** (§1.4 F3).

---

## 0. Binding constraints (each is restated where it governs)

- **D-WORLD-119 and D-WORLD-120 OPEN** (controller, iteration 409: `mission_decisions.sh --open` → both OPEN). The judges encode
  no answer: they judge an `approval.OperatorAction` of either kind, as S2's law does; which may be admitted stays S1's pure
  `OperatorPolicy`. Nothing in S2b is reachable from a route, a store writer or the CLI.
- **No daemon route, no schema change, no `.ail`, `surfaceCSP` and `app.js` unchanged, no served-behaviour change**: the gates'
  `UNTOUCHED_VIOLATIONS` (every changed file other than a `_test.go` or the one NEW production file
  `host/authority/judge_operator.go`, plus any non-test file under `cmd`, `docs`, `host/daemon`, `host/store`, `host/approval`,
  `host/surfacepage`, `world`, `packages`, `scripts`, `tools`, `plugins` or any `.ail`) must say `none`; `DAEMON_ROUTES` base =
  final (47), `SURFACE_POSTS: 3`; `UNCHANGED:` app.js, `surface_page.go` (the CSP literal), `daemon.go` (the route table),
  `host/store/credentials.go` (the schema, `credentialsSchemaVersion = 9` at `:45`); `PREEXISTING_CHANGED: none`;
  `JUDGE_PRODUCTION_CALLERS: none` and `ENVELOPE_PRODUCTION_CALLERS_OUTSIDE_THE_JUDGES: none`, each with a positive control.
- **D-WORLD-109**: every store a test opens is `jgRig`'s schema-current store under `t.TempDir()` (opened read-only by the
  judges); every key is approvaltest's seeded test signer or passkey (`jgSigner` = `NewSigner(62)`, `jgPasskey` =
  `NewPasskey(AlgEdDSA, 61)`, NEW `ojStranger` = `NewSigner(91)`, NEW `ojOtherKey` = `NewPasskey(AlgEdDSA, 92)`, all package-level
  `var`s generated once; S2's `opSigner`/`opPasskey` in host/approval), whose SSHSIG form is byte-compared with real `ssh-keygen`
  by the pre-existing `TestTestSignerAgreesWithSSHKeygen`, run in the gates' `inv` leg with `WORLD_REQUIRE_SSH_KEYGEN=1`. No
  `~/.ssh`, no `~/.ailang/world`, no `/dev/tty`, no `gh`, no network, no browser.
- Coding standards: the law stays pure in host/approval (unchanged); the adapters are the effect boundary (a read-only store
  lookup, a caller-supplied deadline context, no root context, no write); host/approval still imports the standard library only
  (`TestApprovalImportsOnlyTheStandardLibrary`, unedited, green).
- **Not touched**: `tools/launchd/*`, any skill, `world/*.ail`, `packages/`, `scripts/`, `cmd/` (all of it), `plugins/`, `app.js`,
  `surface_page.go`, `daemon.go`, every `host/store`, `host/daemon` file, every production file of `host/approval`, every
  pre-existing file anywhere (`judge_verb.go`, `judge.go`, `judge_test.go` included: reused, not edited).

---

## 1. Pre-flight

### 1.1 Baseline at `4b1c39e2` (measured by the planner, 2026-11-18)

`tools/gates.sh <base-tree> logs/base-gates base 4b1c39e2` (the §6 script, byte-equal to the shipped one but for the `NEW_AU`
line, which `base` mode empties; on a clone of the repo checked out at `4b1c39e2`; `GATES_SKIP_VERIFY_GO=1`):

| leg | rc | time | note |
|---|---|---|---|
| `TOUCHED_PACKAGES` | — | — | *(empty at the base)* |
| `status`, `gofmt`, `build`, `vet` | 0 | 0 / 0 / 5 / 2 s | |
| `verify_ail` | 0 | 41 s | `✓ verify gate PASSED: 16 required identities verified, 40 named tests pass, 503 se-tools named tests pass` — the final must match |
| `fence` | 0 | 24 s | `go test -run '^$' ./...` is the test-file compile fence |
| `pkgs` | 0 | 2 s | `host/approval`, `host/approval/approvaltest`, `host/authority` whole |
| `race` | 0 | 13 s | the same three whole under `-race` (host/authority 5.3 s) |
| `inv` | 0 | 8 s | **36 named `--- PASS`**, 0 SKIP, 0 FAIL (`WORLD_REQUIRE_SSH_KEYGEN=1`: the two `ssh-keygen` cross-checks ran) |
| `DAEMON_ROUTES` / `SURFACE_POSTS` | — | — | **47** (base 47) / **3** |
| `UNTOUCHED_VIOLATIONS` / `PREEXISTING_CHANGED` | — | — | none / none |
| `JUDGE_PRODUCTION_CALLERS` / `ENVELOPE_…_OUTSIDE_THE_JUDGES` | — | — | none (control: `authority.RecheckVerb` production callers 1) / none (control: `approval.CheckConfirmation` production callers 1) |
| `UNCHANGED` | — | — | app.js, surface_page.go, daemon.go, host/store/credentials.go |
| **`GATES_DONE 96s FAILS: none`** | | | `verify_go` is the controller's |

### 1.2 Codebase claims checked (at `4b1c39e2` unless "work" — the planner's stage tree, tag `M3`)

| # | claim | command | observed |
|---|---|---|---|
| C1 | **Input 1 — every cite the patch and design §3.1/§5.2 name, re-measured at `4b1c39e2`** | `grep -n '^func sshApprover\|^func passkeyApprover\|…' host/authority/judge_verb.go host/authority/judge.go host/approval/*.go host/store/*.go` | `judge_verb.go:18` `VerbJudgeInput`, `:24` its `admit()`, **`:30` `sshApprover`, `:44` `passkeyApprover`, `:63` `JudgeVerbSSH`, `:82` `JudgeVerbAssertion`, `:99` `ErrUnknownArtifact`, `:107` `RecheckVerb`** (design §3.1 [r1] cites `:30,44,63,82,107` — **identical**); `judge.go:22–26` `JudgeStepParse/Store/UnknownKey/StoredKey/Law`, `:30` `ErrNotEnrolled`, `:34` `JudgeError`; host/approval `sshsig.go:182` `ParseSSHSIG`, `:38` `ErrNamespace`, `passkey_json.go:77` `ParseAssertionResponse`, `cose.go:91` `ParseStoredCOSEKey`, `signed.go:75` `checkApprover` (P2–P5), `operator_signed.go:17` `SignedOperatorAction`, `:29` `OperatorConfirmInput`, `:40` `OperatorCheckInput` (no creation time), `:58` `ConfirmOperatorAction`, `:69` C2, `:79` `CheckOperatorAction`, `:84` R1, `:91` `operatorProof`, **`:126` `bytes.Clone(p.CredentialID)`** (judge-408's `X-PK-ALIAS` site: "line ~126" stands), `:163` `PasskeyUse`'s clone; host/store `approvers.go:140` `ApproverByKey` (`public_key = ?`, **no format filter**), `passkeys.go:191` `PasskeyByCredentialID` (`w.credential_id = ?`, no format filter), `credentials.go:262` `OpenCredentialsReadOnly`, `:45` `credentialsSchemaVersion = 9`; tests `judge_test.go:21` `jgOrigin/jgRPID`, `:24–25` `jgPasskey`/`jgSigner`, `:30` `jgRig`, `:52` `jgExec`; `m7b_deadline_test.go:10` `boundedTestContext`. `git diff --stat cee5174c 4b1c39e2 -- host cmd`: empty (the record merge touched no code) |
| C2 | **Input 2 — the closest landed analogue, mirrored** | `sed -n` of `judge_verb.go`, `judge_verb_test.go`, `judge_decision.go` | **Mirrored**: V3's `judge_verb.go` — the input struct with `admit()` (`VerbJudgeInput` → `OperatorJudgeInput`), the approver lookups reused unchanged, the artifact dispatch of `RecheckVerb` (`:107–129`: SSHSIG first, else an assertion, else `JudgeStepParse` with `errors.Join(ErrUnknownArtifact, sigErr, respErr)`) factored into `operatorApprover`, every law refusal wrapped `JudgeStepLaw`; V3's `judge_verb_test.go` rig (`jgRig`/`jgExec`, store edits by SQL, read-only reopen). **Differs**: (1) one dispatching judge per role instead of one per key kind (`JudgeVerbSSH`/`JudgeVerbAssertion`), because S2's law takes a `Format`; (2) the judges return the authority-owned proof `*VerifiedOperatorAction` wrapping S2's `*approval.SignedOperatorAction`, where V3 returns the law's `*approval.SignedVerb` directly — design §3.1 [r1] says the proof "can originate only in these judges"; (3) the recheck compares the **recorded approver identity** (`ErrOperatorApprover`) after the law, which `RecheckVerb` does not (design §3.1 "compares the approver identity before any effect"); (4) the decision judge (`judge_decision.go`) was read, not mirrored: it judges a decision row, not an envelope |
| C3 | the parked patch, re-established first-party at `4b1c39e2` | `patch -p1 -E --dry-run`, then apply on a clone, `gofmt -l`, `go vet ./host/authority/`, `go test -count=1 -v` its three tests, `go test ./host/authority/` | applies cleanly; gofmt-clean; vet clean; `TestOperatorJudgeBothKeys`, `TestOperatorRunnerRechecks`, `TestVerifiedOperatorActionOnlyFromTheJudges` `--- PASS`; package `ok`. **272 file lines = 260 changed lines** (108 + 152 numstat; the 12 are diff headers) — both inherited numbers hold |
| C4 | **the parked patch's tests are NOT behaviourally discriminating** (never mutant-tested) | `tools/mut.py park-mut park0 <the M1/M3 ids>` (the parked patch on `4b1c39e2`; whole-package verdicts) | 43 ids: 5 `NM-text` (the parked source inlines the `AdmitInput`), 25 killed, **13 SURVIVED the whole host/authority package**: `PARSE-SENTINEL`, `ADMIT-NOW`, `ADMIT-STORED`, `ADMIT-CONSUMED`, `ADMIT-BYTES`, `R-APPROVER-ERR`, `R-IDENTITY-STEP`, `R-LAW-STEP`, `BINDS-DELEGATE`, `PASSKEYUSE-DELEGATE`, `PASSKEYUSE-NIL`, `RECHECKED-SIGNED`, `WN-LOG-INJECT` (no writes-nothing test). Routed: §1.3 S2b-D2 |
| C5 | every type, field and function the patches name | `grep -n` at HEAD | **existing (reused, unedited)**: `sshApprover`, `passkeyApprover`, `ErrUnknownArtifact`, `ErrNotEnrolled`, `JudgeError{Step}`, `JudgeStepParse/Store/UnknownKey/StoredKey/Law`, `jgRig`, `jgExec`, `jgSigner`, `jgPasskey`, `jgOrigin`, `jgRPID`, `boundedTestContext` (host/authority); `approval.OperatorAction`, `OperatorOptions`, `OperatorResolved`, `OperatorActionSchema`, `OperatorActionNamespace`, `OperatorArgv`, `OperatorSignatureDigest`, `OperatorActionMessage`, `OperatorActionMessageHash`, `ConfirmOperatorAction`, `OperatorConfirmInput{Action, Admit, Challenge, Artifact, Format, Approver}`, `CheckOperatorAction`, `OperatorCheckInput{Action, Challenge, Expiry, Now, Artifact, Format, Approver}`, `SignedOperatorAction` and its `Binds/ApprovedBy/PasskeyUse`, `ErrOperatorLifetime`, `ErrOperatorCounter`, `AdmitInput`, `Verdict`, `RefuseConsumed/Expired/DigestMismatch`, `NotAdmittedError`, `BindingError{Field}`, `ErrNotVerified`, `ErrApproverRevoked`, `ErrApproverClockBehind`, `ErrApproverFormat`, `ErrApproverNotPasskey`, `ErrSignCountRegression`, `ErrNamespace`, `ErrBadSignature`, `VerbNamespace`, `FormatSSH`, `FormatWebAuthn`, `ParseSSHSIG`, `ParseAssertionResponse`, `Approver`, `AlgEdDSA` (host/approval); `opIn`, `opAction`, `opSSH`, `opAssert`, `opSigner`, `opPasskey`, `opC`, `opCreated`, `opExp`, `opNow`, `paApprover` (host/approval tests); `store.Credentials`, `OpenCredentialsReadOnly`, `(*Credentials).Close` (host/store); `approvaltest.NewSigner/NewPasskey/Sign/Assert/PublicKey/CredentialID`. **NEW** (host/authority): `VerifiedOperatorAction{signed, rechecked}`, `verifiedOperatorAction`, `ErrOperatorApprover`, `OperatorJudgeInput{Action, Digest, Challenge, BytesDigest, CreatedAt, ExpiresAt, Now, Consumed}`, `(OperatorJudgeInput).admit`, `operatorApprover`, `JudgeOperatorAction`, `RecheckOperatorAction`, `(*VerifiedOperatorAction).Binds/ApprovedBy/PasskeyUse/Rechecked`; tests `ojStranger`, `ojOtherKey`, `ojAction`, `ojIn`, `ojSSH`, `ojAssert`, `ojOpen`, `ojRefused`, `ojSilent` and every test named in §3 |
| C6 | **Input 3 — fixtures from the real tool** | the test files | every SSHSIG is `approvaltest.Signer.Sign` (the repo's test signing path, byte-compared with `ssh-keygen -Y sign` by `TestTestSignerAgreesWithSSHKeygen`, `approvaltest/sign_crosscheck_test.go:48`); every WebAuthn assertion is `approvaltest.Passkey.Assert` (Chrome's layout), the passkey enrolled through the real `store.EnrolPasskey` by `jgRig` (`judge_test.go:30`), and in host/approval through the real `NewPasskeyApprover` by `paApprover` (`passkey_law_test.go:18`); every message and hash is the real `OperatorActionMessage[Hash]` over the real `OperatorSignatureDigest`; every store is the real schema in `t.TempDir()`, edited by SQL through `jgExec` and reopened by `OpenCredentialsReadOnly`; the one inserted approver's key is `ojStranger.PublicKey().Blob` formatted as a SQL blob literal. No hand-typed bytes, keys or digests |
| C7 | **Input 4 — tests that enumerate packages, imports, symbols, routes, writer sites, schema, context roots, CLI** | `grep -rn` by name; each run `-v` on the final tree (C17 `inv`) | **unchanged and green**: `TestApprovalImportsOnlyTheStandardLibrary`, `TestApprovalTestIsTestOnly`, `TestDaemonDependencyAllowlist` (host/daemon imports host/authority; S2b adds no import to host/authority: `judge_operator.go` imports `context errors host/approval host/store`, all already imported by `judge_verb.go`), `TestWorldBoundaryDependencyAllowlist`, `TestBareNetHTTPExemptionIsPerGroup`, `TestEveryCommandDirectoryIsAProtectedGroup`, `TestBoundaryASTWriteGuard`, **`TestProductionContextRoots`** (no new `context.Background`/`TODO`/`WithoutCancel` in production: the judges take the caller's context; the pin table is unchanged), `TestProductionGoSurface` (the new file is under `host/`), `TestStoreSurfaceTakesContext`, `TestCredentialsSchemaMatchesCanonicalDDL`, `TestCredentialsSchema9VerbRequests`, `TestCredentialsNewerThanNineIsRefused`, `TestFrozenV1RouteTableMatchesMux`, `TestSurfaceAppRoutes`, `TestApprovalRoutesRegistered`, `TestDecisionRoutesRegistered`, `TestConfirmRoutesRegistered`, `TestDaemonCredentialWriterSites`, `TestSurfaceVerbRouteRegistered`, `TestSurfaceVerbIsReadOnly`, `TestSurfaceVerbOutcomes`, `TestEveryVerbHasHelp`, `TestVerbTableMatchesTheCLISwitch`, `TestCLICredentialWriterSites`, `TestOperatorFlagParity`. **Changed: none.** No test enumerates host/authority's exported symbols. **New inventory**: `TestVerifiedOperatorActionOnlyFromTheJudges` (AST + reflect; static, so no mutant's verdict rests on it, §4) |
| C8 | **no route, no writer site, no schema change, no production caller — by test or caller search with a positive control** | the gates (C17) | `DAEMON_ROUTES: 47 (base 47) SURFACE_POSTS: 3`; `TestDaemonCredentialWriterSites`, `TestCLICredentialWriterSites` green unedited (no new writer: the judges call only `ApproverByKey` and `PasskeyByCredentialID` on a read-only handle); `UNCHANGED: host/store/credentials.go`; `JUDGE_PRODUCTION_CALLERS: none (control: authority.RecheckVerb production callers 1)` over all seven NEW identifiers (exported and unexported) outside `judge_operator.go`; `ENVELOPE_PRODUCTION_CALLERS_OUTSIDE_THE_JUDGES: none (control: approval.CheckConfirmation production callers 1)` — S2's 21 identifiers now have exactly one production caller file, the judges' |
| C9 | **judge-408's `X-PK-ALIAS` reproduced on the base first** | `tools/mut.py base-mut base-whole X-PK-ALIAS PK-USE-ALIAS R1-LIFETIME-ADDED` (a `git archive 4b1c39e2` export) | `X-PK-ALIAS` (`bytes.Clone(p.CredentialID)` → `p.CredentialID` at `operator_signed.go:126`): **SURVIVED the whole host/approval package** — the controller's iteration-408 reproduction holds. So did its accessor twin `PK-USE-ALIAS` (`:163`) and `R1-LIFETIME-ADDED` (judge-408 note 3 unpinned). All three **KILLED** on the final tree (§4) |
| C10 | **Input 5 — no asynchronous wait; time is an input** | the test files | no S2b test sleeps, ticks, polls or reads a clock: every judge takes `Now`/`CreatedAt`/`ExpiresAt`/`now` as integers. **The S2b production code has no time comparison of its own** (it wires times into S2's `Admit`, C2 and R1, each of whose operators S2's ledger already mutated: `P1-ADMIT`, `C2-*`, `R1-*`); S2b's time mutants are the **wiring**: `ADMIT-NOW`, `ADMIT-CREATED`, `ADMIT-EXPIRES`, `R-NOW-WIRE`, `R-EXPIRY-WIRE`, each killed at the strict boundary (`now == expiry` refused; lifetime 601 refused; recheck at `expiry − 1` admitted, at `expiry` refused), and `R1-LIFETIME-ADDED` (the runner's comparison) |
| C11 | the new tests, work | `go test -count=1 -v` by name (`accept.sh`) | all `--- PASS`, none SKIP |
| C12 | **judge-408 note 1 — does a judge ever see a caller's approver?** | `go doc`-level read of both signatures; `TestVerifiedOperatorActionOnlyFromTheJudges` | no: `JudgeOperatorAction(ctx, ro, OperatorJudgeInput, artifact)` and `RecheckOperatorAction(ctx, ro, action, challenge, expiry, now, artifact, approvedBy)` take no `approval.Approver` and no struct holding one (reflect over both function types; positive control: `approval.OperatorConfirmInput` holds one and is found). Behaviourally: every approver fact the law reads moves with a **store** edit in `TestOperatorJudgeReadsTheStore` (revoked, enrolled after now, format, counter, stored key, unknown key, unknown credential, unreadable store) and the identity with `TestOperatorJudgeBothKeys`'s `ops-ssh` row; each is its own mutant (§4 `REUSED-*`) |
| C13 | the snapshot block (§7) on every sub-milestone | `tools/snap_rehearse.sh <this plan>` (a clone at `4b1c39e2`, the patches in order, §7's block verbatim) | `M1 2/0/2`, `M2 3/0/3`, `M3 4/0/4` (`STATUS` lines / `DELETED` lines / MANIFEST lines), each `MANIFEST_OK`, rc 0; every stage has only `host/` (the `find` guard skips the absent `cmd/` and `docs/`); the block extracted from this file is §7's verbatim but for `M`, `W` and the snapshot root (`BLOCK_VERBATIM_BUT_W_M_D`) — `logs/snap.txt`, re-run on this final plan file |
| C14 | **every sub-milestone is green when committed** | `tools/accept.sh` on each stage tree (`st-M1`…`st-M3`) | `M1 fast rc=0 ACCEPT_OK` (15 named top-level `--- PASS`), `M2` (17), `M3` (18), each 5–7 s warm, 30–35 s cold — every part far inside one 600 s call (`logs/accept/*.txt`) |
| C15 | the pinned binaries | `$HOME/.pinned-ailang/ailang --version`; `go version`; `which ssh-keygen` | `AILANG v0.41.0`; `go1.26.6 darwin/arm64`; `/usr/bin/ssh-keygen` (the `inv` leg's `TestSSHSIGCrossCheckWithSSHKeygen` and `TestTestSignerAgreesWithSSHKeygen` run, not skip) |
| C16 | **the round trip** | `tools/roundtrip.sh <this plan>` | `CMP_EQUAL` for all three (each block extracted from this file with §2's awk equals the generated patch byte for byte); each applied with `/usr/bin/patch -p1 -E` after a `--dry-run` to a fresh `git archive 4b1c39e2` export, no `.orig`/`.rej`, committed and measured: `NUMSTAT M1 401`, `M2 469`, `M3 526` (each `unfiltered=` the same); after each, `diff -r` of `cmd` and `host` against `st-<M>` rc 0; then `diff -r -x .git` of the whole export against `st-M3` → **`ROUNDTRIP_IDENTICAL_TO_STAGE_TREE`**; `SCRIPT_EQUAL accept.sh`, `SCRIPT_EQUAL gates.sh`, `SAVED_GATES_EQUAL` (`$HOME/.ailang/state/world-iter409/gates.sh` is §6's) — `logs/roundtrip.txt`, re-run on this final plan file. The controller's `ctl_roundtrip.sh` was **not** run in place (it refuses a reused export dir); a copy with only its dir and its loop (`M1 M2 M3`) changed: `M1 applied (413 lines)`, `M2 (74)`, `M3 (63)`, every `diff -r cmd`/`host` vs `st-<M>` rc 0, `DONE`; the two `.orig` files it lists (`design_docs/verification/world-iter362/evaluator-r1/mutations/{doctor,exec_probe}.go.orig`) are **tracked at the base** (`git ls-files`), not patch residue — `logs/ctl_sim.txt` |
| C17 | **the gates on the final tree** | `tools/gates.sh <final-tree> logs/final-gates final 4243349a` | (`final-tree` = the planner's clone at the stage tag `M3`, the three stage commits on `4b1c39e2`, equal to the round-trip export, C16; `GATES_SKIP_VERIFY_GO=1`): **`TOUCHED_PACKAGES: host/approval host/authority`** (= the plan's list); `status` 0; **`verify_ail` rc 0** (the base's PASSED line, byte for byte); `gofmt`, `build`, `vet`, `fence` rc 0; **`pkgs` rc 0**; **`race` rc 0** (host/approval 2.0 s, approvaltest 1.5 s, host/authority 5.7 s); **`inv` rc 0: 44 named `--- PASS`** (the base's 36 + the eight new), 0 SKIP, 0 FAIL, no `NOT PASSED`; **`DAEMON_ROUTES: 47 (base 47)`, `SURFACE_POSTS: 3`, `UNTOUCHED_VIOLATIONS: none`, `PREEXISTING_CHANGED: none`, `JUDGE_PRODUCTION_CALLERS: none (control 1)`, `ENVELOPE_PRODUCTION_CALLERS_OUTSIDE_THE_JUDGES: none (control 1)`, `UNCHANGED:` app.js, surface_page.go, daemon.go, credentials.go**; **`GATES_DONE 61s FAILS: none`**. `verify_go` is the controller's |
| C18 | **"wrong format" is a store fact the law can see** | `sed -n 140,159p host/store/approvers.go`; `TestOperatorJudgeReadsTheStore` rows | `ApproverByKey` matches the key blob only and `PasskeyByCredentialID` the credential ID only, so an SSH key stored on a `webauthn`-format row (`UPDATE approvers SET format = 'webauthn'`) reaches the law, which refuses it at P3 (`ErrApproverFormat`), and a passkey on an `ssh`-format row at P3w (`ErrApproverNotPasskey`); mutants `REUSED-SSH-FORMAT`/`REUSED-PK-FORMAT` (the adapter hard-coding the format instead of reading it) are killed by those rows |

### 1.3 Planner decisions where the design or the directive leaves a choice

- **S2b-D1 — no cut, three sub-milestones, by measurement.** Compiled: M1 401, M2 +68, M3 +57 = **526** (cap 1,400). M1 lands the
  production file with every test that kills its mutants (production never lands untested); M2 is test-only in host/approval
  (judge-408 notes 2 and 3 are about S2's law, not the adapters, so they are independent of M1 and could land in either order;
  they go second so M1 is the reviewable unit); M3 is test-only, the writes-nothing rule.
- **S2b-D2 — what changed in the parked patch, and why** (C4 measured the parked tests: 13 whole-package survivors):
  - *production (`judge_operator.go`, 108 → 116 lines)*: (1) `OperatorJudgeInput.admit()` replaces the inline `AdmitInput`
    literal — V3's `VerbJudgeInput.admit()` shape (`judge_verb.go:24`), and each field's wiring becomes one mutant line;
    (2) the doc comments say what judge-408 note 1 asked: neither judge takes an approver; the runner's bound is expiry and
    revocation, the lifetime bound is the daemon's. **No behaviour changed**: the three parked tests, applied unedited beside the new
    file (`st-M1` with the parked test file swapped in), `--- PASS` — superseded, not weakened.
  - *tests (`judge_operator_test.go`, 152 → 285 lines)*: (1) **NEW `TestOperatorJudgeReadsTheStore`** — 19 rows, every one a store
    edit (unknown SSH key, unknown passkey credential, unreadable store ×2, revoked ×2, enrolled after now ×2, wrong stored format
    ×2, stored counter at the assertion's, unparsable stored key) or a row edit (consumed, now at expiry, lifetime 601, stored
    bytes differ, another action's digest column) or an artifact (typed y, a verb namespace), each asserting **step, class and
    no proof** through `ojRefused` (directive retro 4); the parked test asserted only the step on three rows; (2)
    `TestOperatorJudgeBothKeys` adds `Binds` of another action (kills `BINDS-DELEGATE`), the passkey use (kills
    `PASSKEYUSE-DELEGATE`) and an `ops-ssh` key enrolled under its own name (kills `REUSED-SSH-NAME`: the approver name is the
    store's); (3) `TestOperatorRunnerRechecks` checks both controls' approver and binding, and its refusals add the step (kills
    `R-LAW-STEP`, `R-IDENTITY-STEP`), an unknown key (kills `R-APPROVER-ERR`), typed y, another challenge and an empty recorded
    approver; (4) the nil/zero-value assertions move from the AST test into **NEW `TestOperatorUnjudgedProofBindsNothing`** — run 1
    showed the five nil-guard mutants were killed only by the static test (`w/o static: SURVIVED`), and a row with
    `{rechecked: true}` and no proof (kills `RECHECKED-SIGNED`); (5) the AST test also proves no judge takes an approver (C12);
    (6) fixtures sign through `ojSSH`/`ojAssert` helpers with every error checked (the parked `ojIn` dropped three errors).
- **S2b-D3 — judge-408 note 3, by measurement: PIN, do not add the bound.** Design §3.1 line 80 puts the 600 s bound on the
  *signed request* ("signed request expiry 600 s"), and S2's `ConfirmOperatorAction` enforces it at C2 (`operator_signed.go:69`)
  where the row's creation time is an input; design §3.1 line 83 gives the runner "Recheck expiry and revocation again at each
  new effect boundary", and the last bullet "its signed expiry limits **new** effects". `OperatorCheckInput` (`:40`) carries no
  creation time, so a lifetime check at the runner would need an input the design does not give it; the expiry it does check is
  inside the signed message, so it cannot be lengthened without a new signature. `TestOperatorRunnerRecheckBoundsExpiryNotLifetime`
  states this reason in its comment, shows the confirmation refusing a 3,600 s request (`ErrOperatorLifetime`) and the recheck
  admitting the same recorded expiry at `now` and `expiry − 1` for both key kinds; `R1-LIFETIME-ADDED` (a remaining-time bound
  added to R1) is killed by it, and SURVIVED at the base (C9).
- **S2b-D4 — judge-408 note 2: kill it in host/approval, with its twin.** `TestOperatorPasskeyUseIsACopy` confirms and rechecks with
  a passkey, then flips a byte of the caller's `Approver.Passkey.CredentialID` and requires both proofs' `PasskeyUse` unchanged
  (kills `X-PK-ALIAS`), then flips a byte of a returned slice and requires the next call unchanged (kills `PK-USE-ALIAS`, the
  accessor's clone at `:163`, equally unpinned at the base). It is a NEW file: no S2 test is edited.
- **S2b-D5 — judge-408 note 4: not acted on, carried (Appendix C).** The field drops it names (episode, kind, argv, options) are
  **binding-redundant by construction**: Validate (O9) requires `argv == OperatorArgv(kind, episode, options)`, so dropping one of
  kind/episode/options from the encoding leaves it bound through `argv`, and dropping `argv` leaves it bound through the other
  three — a binding test cannot distinguish the mutant from the original; the round-trip tests (`TestOperatorStrictJSON`,
  `TestOperatorRequestLimits`, behavioural, golden excluded per judge-408) already kill them. Not cheap and not in a file S2b
  touches; the plan does not overclaim it.
- **S2b-D6 — the two equivalents are the lookups' own equality.** `REUSED-SSH-KEY` and `REUSED-PK-CREDID` replace a stored value with
  the artifact's own bytes; the store selected the row **by** those bytes (`public_key = ?`, `credential_id = ?`), so they are
  equal on every row returned. Named, not chased.
- **S2b-D7 — no Chrome, no node, no docs, no cmd**: S2b adds no page, route, command, flag or output line. The gates have no
  `chrome` or `node` leg.

### 1.4 Premises found false or incomplete at HEAD (each routed around minimally)

| # | premise | measured | routing |
|---|---|---|---|
| F1 | directive: the design's and the patch's cites may have moved since `cee5174c`/`41331b4d` | C1: none moved (`git diff cee5174c 4b1c39e2 -- host cmd` empty) | none; cites stand |
| F2 | directive: the parked patch "compiled, vetted, gofmt-clean, green" (UNVERIFIED, inherited) | C3: **true** at `4b1c39e2` | none |
| F3 | S2 plan Appendix C / directive: the parked tests are adequate ("lands it ... with its own mutant pass"; never run) | C4: **false** — 13 whole-package survivors of 38 compiled ids (12 excluding the missing writes-nothing test) | S2b-D2: tests rewritten and extended; run 2 has 0 survivors |
| F4 | judge-408 note 3: the runner might be expected to re-bound the lifetime | design §3.1 lines 80/83: the bound is the signed request's, judged at confirmation; the runner rechecks expiry and revocation; `OperatorCheckInput` has no creation time | S2b-D3: pinned with a named test and a mutant, no production change |
| F5 | judge-408 note 4 / S2 §9 item 4: `TestOperatorConfirmBinding` could kill each field drop by name | the four fields are mutually bound through Validate's O9 (S2b-D5) | carried, with the reason; not a test gap a binding row can close |
| F6 | directive (controller tooling): `ctl_roundtrip.sh` loops `M1 M2 M3 M4` | S2b ships three sub-milestones | the controller edits its loop to `M1 M2 M3` |

**Files outside design §7**: none — §7 names `host/{approval,authority,store,verbs,daemon,surfacepage}`; S2b adds files to
`host/authority` and one test to `host/approval`.

### 1.5 The size and the race-mode cost

| stage | production | test | changed lines | cumulative (numstat after commit) |
|---|---|---|---|---|
| M1 | 116 (`judge_operator.go`) | 285 (`judge_operator_test.go`) | 401 | **401** |
| M2 | 0 | 68 (`host/approval/operator_pins_test.go`) | 68 | **469** |
| M3 | 0 | 57 (`judge_operator_silent_test.go`) | 57 | **526** |

S2's planner measured the parked adapters at 260; S2b lands them at 401 (+141: the store matrix, the split nil test and the
stricter controls), and the whole sprint at 526 against S2 Appendix C's "~260 + pins".

**Race-mode wall time of the new tests** (directive 8): work tree, `go test -race -count=1 -v` by name (`logs/race-new.txt`):
host/authority — `TestOperatorJudgeReadsTheStore` 0.48 s (19 temp stores), `TestOperatorRunnerRechecks` 0.28 s,
`TestOperatorJudgesWriteNothing` 0.10 s, `TestOperatorJudgeBothKeys` 0.08 s, the other two < 0.01 s — **≈ 0.95 s added**; the package
under `-race` 5.6 s twice against 4.6–5.3 s at the base (gates: 5.7 s vs 5.3 s). host/approval — the two pins < 0.01 s each; the
package 2.1 s vs 2.2 s at the base (noise). **≈ 1 s added in all**, in packages far under CI's 8-minute per-package bound; host/daemon
and host/broker gain nothing (no S2b test there). No key generation at run time beyond the four package-level approvaltest `var`s
(seeded, generated once per package).

---

## 2. How the executor applies a sub-milestone (the patch protocol)

Each sub-milestone `<M>` is the block fenced by a line `` ````diff <M> `` (four backticks) and a line of four backticks,
in Appendix P. Paths are `a/<repo path>` and `b/<repo path>`; every S2b file is new (`--- /dev/null`); **no sub-milestone
deletes or changes a pre-existing file**. Apply it **only** this way, in one `bash` shell (the controller's shell is zsh):

```bash
bash <<'EOF2'
set -euo pipefail
M=M1   # <- the sub-milestone, exactly as in the Appendix P heading, in order: M1 M2 M3
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter409-row243-v4s2b
PLAN="$W/design_docs/planned/w-world-address-bar-verbs-v4s2b-sprint-plan.md"
T="$HOME/.ailang/state/world-iter409/exec-patches"; mkdir -p "$T"
awk -v m="$M" '$0 == "````diff " m {f=1; next} f && $0 == "````" {exit} f' "$PLAN" > "$T/$M.patch"
test -s "$T/$M.patch"
cd "$W"
/usr/bin/patch -p1 -E --dry-run < "$T/$M.patch"
/usr/bin/patch -p1 -E < "$T/$M.patch"
if find cmd host \( -name '*.orig' -o -name '*.rej' \) 2>/dev/null | grep .; then echo "STOP: patch left .orig/.rej"; exit 1; fi
git status --short --untracked-files=all
EOF2
```

- If `--dry-run` fails, **stop and report** (no hand edits, no fuzz options). It can only fail if an earlier
  sub-milestone was skipped, reordered or changed.
- Then run that sub-milestone's **Acceptance** (§3): extract `accept.sh` from this plan (`awk '$0 == "````bash
  accept.sh" {f=1; next} f && $0 == "````" {exit} f' "$PLAN" > "$T/accept.sh"`) and run its one part:
  `bash "$T/accept.sh" "$W" <M> fast 2>&1 | tee "$HOME/.ailang/state/world-iter409/exec-logs/accept_<M>_fast.txt"` (never
  only through `tail`). **Never `nohup`, never background a test leg; pass a 600000 ms timeout to every long Bash call.**
  No part launches Chrome or regenerates a golden.
- Then snapshot (§7). The controller commits and compares the numstat with §3's expected value.

---

## 3. Milestones

**Acceptance.** Every part runs the same `[base]` checks first (gofmt over `cmd host`, build, vet of host/approval and
host/authority, the test-file compile fence `go test -run '^$' ./...`), then the named tests with `-v` (every name `--- PASS`,
none SKIP or FAIL) beside V3's verb-judge tests and S2's law tests as regression witnesses, then host/approval, approvaltest
and host/authority whole, then three host/daemon inventories. One part, `fast`, per sub-milestone (C14 measured each, well
inside one 600 s call).

````bash accept.sh
#!/bin/bash
# accept.sh <tree> <M> [fast]: exactly plan §3's acceptance for sub-milestone M (one bash shell; rc 0 =
# ACCEPT_OK). One part per sub-milestone, each well inside one 600 s call. No part launches Chrome (S2b has
# no Chrome test) and no part regenerates a golden.
T=$1; M=$2; P=${3:-fast}
cd "$T" || exit 2
export AILANG_BIN=$HOME/.pinned-ailang/ailang
export WORLD_EXEC_NODE=/opt/homebrew/bin/node
export WORLD_REQUIRE_SSH_KEYGEN=1
unset WORLD_UPDATE_OPERATOR_ACTION_GOLDEN WORLD_UPDATE_VERB_GOLDEN WORLD_UPDATE_CONFIRM_GOLDEN WORLD_UPDATE_CONFIRMATIONS_GOLDEN WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER
set -o pipefail
chk() { tee /dev/stderr | { out=$(cat); echo "$out" | grep -qE -- '--- (FAIL|SKIP)|^FAIL' && exit 1; echo "$out" | grep -q -- '--- PASS' || exit 1; }; }
v() { go test -count=1 -v "$@" 2>&1 | chk; }   # -v: every named test PASS, none SKIP/FAIL
p() { go test -count=1 "$@"; }
step() { echo "== $*"; "$@" || { echo "ACCEPT_RED $M $P: $*"; exit 1; }; }
# [base]
step test -z "$(gofmt -l cmd host)"
step go build ./...
step go vet ./host/approval/ ./host/authority/
if go test -run '^$' ./... 2>&1 | grep -E '^(FAIL|--- FAIL)'; then echo "ACCEPT_RED $M $P: FENCE_RED"; exit 1; fi
A1='TestOperatorJudgeBothKeys|TestOperatorJudgeReadsTheStore|TestOperatorRunnerRechecks|TestOperatorUnjudgedProofBindsNothing|TestVerifiedOperatorActionOnlyFromTheJudges'
A2='TestOperatorPasskeyUseIsACopy|TestOperatorRunnerRecheckBoundsExpiryNotLifetime'
A3='TestOperatorJudgesWriteNothing'
INV_AU='TestJudgeVerbSSH|TestRecheckVerb|TestReadOnlyAtFailuresAreErrors'
INV_A='TestApprovalImportsOnlyTheStandardLibrary|TestOperatorConfirmBothKeys|TestOperatorRunnerRecheckLaw|TestOperatorPureEntryPointsWriteNothing'
case "$M:$P" in
M1:fast) step v ./host/authority/ -run "^($A1|$INV_AU)\$"
         step v ./host/approval/ -run "^($INV_A)\$" ;;
M2:fast) step v ./host/authority/ -run "^($A1|$INV_AU)\$"
         step v ./host/approval/ -run "^($A2|$INV_A)\$" ;;
M3:fast) step v ./host/authority/ -run "^($A1|$A3|$INV_AU)\$"
         step v ./host/approval/ -run "^($A2|$INV_A)\$" ;;
*) echo "unknown $M $P"; exit 2 ;;
esac
step p ./host/approval/ ./host/approval/approvaltest/ ./host/authority/
step v ./host/daemon/ -run '^(TestDaemonDependencyAllowlist|TestConfirmRoutesRegistered|TestDaemonCredentialWriterSites)$'
echo "ACCEPT_OK $M $P"
````

### M1 — the host/authority operator judges and the judge-only proof · 116 prod + 285 test · part: `fast`

- `host/authority/judge_operator.go` (NEW, S2b-D2): `VerifiedOperatorAction{signed, rechecked}` and its only constructor
  `verifiedOperatorAction`; `ErrOperatorApprover`; `OperatorJudgeInput` and `admit()`; `operatorApprover` (an SSHSIG → `sshApprover`,
  else an assertion → `passkeyApprover`, else `JudgeStepParse` joining `ErrUnknownArtifact`); `JudgeOperatorAction` (the approver
  from ro, then `approval.ConfirmOperatorAction`, law refusals → `JudgeStepLaw`, proof not rechecked); `RecheckOperatorAction` (the
  approver from ro at now, then `approval.CheckOperatorAction` over the runner's own action, challenge and expiry, then the
  recorded identity → `ErrOperatorApprover`, proof rechecked); the four nil-safe accessors.
- `host/authority/judge_operator_test.go` (NEW): the rig helpers (`ojAction`, `ojIn` — created 1000, expires 1600, now 1010; mark-ssh's
  SSHSIG and mark-mac's assertion at counter 2 — `ojSSH`, `ojAssert`, `ojOpen`, `ojRefused`) and **`TestOperatorJudgeBothKeys`**
  (both kinds: approver, binding, another action's binding refused, passkey use, not rechecked; `ops-ssh` judged under its stored
  name), **`TestOperatorJudgeReadsTheStore`** (the 19 rows of S2b-D2), **`TestOperatorRunnerRechecks`** (controls for both kinds, the
  passkey's counter recorded; a forged row's action, another challenge, another / no recorded approver, revoked since, expired at
  now, counter never recorded, unknown key, typed y — each step, class and no proof), **`TestOperatorUnjudgedProofBindsNothing`**
  (nil, zero, `{rechecked: true}`), **`TestVerifiedOperatorActionOnlyFromTheJudges`** (every field of both proof types unexported;
  no judge parameter holds an `approval.Approver`, control: `OperatorConfirmInput` does; one constructor literal, called by exactly
  the two judges, the positive control in the same count).
- Expected numstat after commit **401**.

### M2 — judge-408 pins: the passkey use is a copy; the runner bounds expiry, not lifetime · 0 prod + 68 test · part: `fast`

- `host/approval/operator_pins_test.go` (NEW, S2b-D3/D4): **`TestOperatorPasskeyUseIsACopy`** (confirmed and rechecked proofs; the
  caller's credential ID edited after the call; a returned slice edited) and **`TestOperatorRunnerRecheckBoundsExpiryNotLifetime`**
  (both kinds: a 3,600 s recorded expiry refused by the confirmation with `ErrOperatorLifetime`, admitted and bound by the recheck
  at now and `expiry − 1`).
- Expected numstat after commit **469**.

### M3 — the operator judges write nothing · 0 prod + 57 test · part: `fast`

- `host/authority/judge_operator_silent_test.go` (NEW): `ojSilent` (os.Stdout, os.Stderr and the log on one pipe, read after the
  call returns) and **`TestOperatorJudgesWriteNothing`** (control: a write to each is seen; then both judges, against the rig, a
  revoked store and a closed store, on mark-ssh's SSHSIG, mark-mac's assertion, typed y and an unknown key, with two recorded
  approvers → nothing written; `WN-LOG-INJECT` killed).
- Expected numstat after commit **526**.

### 3.1 Every design §5.2 row and judge-408 note S2b drives, traced (directive retro 4)

| design row / note | branch | test (and the no-proof assertion) |
|---|---|---|
| §5.2 Good passkey and good enrolled SSHSIG ("skips proof or discards recorded approver") | both judges, the whole law | `TestOperatorJudgeBothKeys`, `TestOperatorRunnerRechecks` controls (approver = the store's name) |
| §5.2 Plain y, malformed proof | `operatorApprover`'s parse refusal | `…ReadsTheStore` / `…RunnerRechecks` "typed y" (`JudgeStepParse`, `ErrUnknownArtifact`, no proof) |
| §5.2 Forged, unknown, wrong SSH namespace | unknown key / credential lookups; the law | `…ReadsTheStore` "unknown ssh key", "unknown passkey credential" (`JudgeStepUnknownKey`, `ErrNotEnrolled`), "a verb namespace" (`ErrNamespace`) |
| §5.2 Revoked key or replayed/unchanged nonzero counter (the store half) | `RevokedAt`, `SignCount` read from ro | `…ReadsTheStore` "revoked" ×2, "stored counter at the assertion's" (`ErrSignCountRegression`); `…RunnerRechecks` "revoked since", "counter never recorded" |
| §5.2 Expiry equal to now, outside lifetime (the wiring) | `admit()` | `…ReadsTheStore` "now at expiry", "lifetime 601", "consumed", "stored bytes differ", "another action's digest column" |
| §5.2 Forged DB row / mismatch after claim ("trusts row digest, stored approver or daemon status") | `RecheckOperatorAction` | `…RunnerRechecks` (own action re-signed, current key, recorded identity) |
| §3.1 [r1] proof originates only in these judges (judge-408 note 1) | `verifiedOperatorAction`; no approver parameter | `TestVerifiedOperatorActionOnlyFromTheJudges` (static) + `TestOperatorUnjudgedProofBindsNothing` and every store row (behavioural) |
| judge-408 note 2 (`X-PK-ALIAS`) | `operator_signed.go:126`, `:163` | `TestOperatorPasskeyUseIsACopy` |
| judge-408 note 3 (runner lifetime) | `operator_signed.go:84` R1 | `TestOperatorRunnerRecheckBoundsExpiryNotLifetime` |
| judge-407 rule: server entry points write nothing | both judges | `TestOperatorJudgesWriteNothing` |

**Not S2b** (Appendix C names each owner): every CAS/claim/crash/replay/race/recovery/status/prune/cancel row (S3, S4, S6); every
route row (S7); every page row (S8); judge-408 note 4 (carried, S2b-D5).

---

## 4. Mutation table

Enumerated **from the diff** (`tools/mut.py`; the ids below are the ledger's): every branch of `judge_operator.go` — the format
dispatch both ways and both labels, the parse refusal's step and sentinel, each judge's approver-refusal pass-through, every field
of the law's input wired from the judge's arguments (action, challenge, format, approver; `admit()`'s now, created, expires,
consumed, stored digest, bytes digest; the recheck's expiry and now), each law-refusal step, the recorded-identity comparison, its
sentinel and its step, both `rechecked` flags, the constructor's two fields, every accessor's nil guard and delegation, the three
conjuncts of `Rechecked`; an injected log line (the writes-nothing harness's mutant control); **every store-sourced fact the judges
reach** in the reused `judge_verb.go` (store error and unknown row for both lookups, the SSH name, format, key, enrolment and
revocation; the passkey's stored key parse, format, enrolment, revocation, counter and credential ID); and **judge-408's
host/approval sites** (`X-PK-ALIAS`, its accessor twin, and the runner-lifetime pin's mutant). S2b's production has no time
comparison of its own (C10). Each is one or more exact replacements (text not found exactly once → `NM-text`; a mutant that
does not compile → `NM`, then rewritten to compile and re-run). Each runs against its **named** test, then the **WHOLE** package
of the mutated file (`-run '.'`: every test, static ones included); a mutant whose only whole-package killer is the static
`TestVerifiedOperatorActionOnlyFromTheJudges` is re-run with it skipped and the leg recorded (`w/o static`).
Ledger: `$HOME/.ailang/state/world-iter409/plan-scratch/mutants/ledger.txt` (one line per run, keyed by tag: `final1` — run 1 on the
pre-split M3 draft, whose six `NM` were rewritten to compile and re-run, and whose five nil-guard kills were static-only
(`w/o static: SURVIVED`), which drove `TestOperatorUnjudgedProofBindsNothing`; `base-whole` — C9; `park0` — C4, the parked patch;
`final2` — the table below, run 2 on the final stage tree).

**Base reproduction (C9; tag `base-whole`, a `git archive 4b1c39e2` export)**: `X-PK-ALIAS`, `PK-USE-ALIAS`, `R1-LIFETIME-ADDED` —
each named SURVIVED (no such test) and **whole host/approval SURVIVED**.

**Final (run 2, tag `final2`, the M3 stage tree)**:

| id | stage · file | mutation | named test | named | whole package | verdict |
|---|---|---|---|---|---|---|
| `FMT-SSH-OFF` | M1 · `judge_operator.go` | `` approval.ParseSSHSIG(artifact); sigErr == nil { `` → `` approval.ParseSSHSIG(artifact); sigErr == nil && false { `` | `TestOperatorJudgeBothKeys` | killed (TestOperatorJudgeBothKeys) | killed (TestOperatorJudgeBothKeys,TestOperatorJudgeReadsTheStore,TestOperatorRunnerRechecks) | K |
| `FMT-PK-OFF` | M1 · `judge_operator.go` | `` approval.ParseAssertionResponse(artifact); respErr == nil { `` → `` approval.ParseAssertionResponse(artifact); respErr == nil && false { `` | `TestOperatorJudgeBothKeys` | killed (TestOperatorJudgeBothKeys) | killed (TestOperatorJudgeBothKeys,TestOperatorJudgeReadsTheStore,TestOperatorRunnerRechecks) | K |
| `FMT-SSH-LABEL` | M1 · `judge_operator.go` | `` return approval.FormatSSH, a, err `` → `` return approval.FormatWebAuthn, a, err `` | `TestOperatorJudgeBothKeys` | killed (TestOperatorJudgeBothKeys) | killed (TestOperatorJudgeBothKeys,TestOperatorJudgeReadsTheStore,TestOperatorRunnerRechecks) | K |
| `FMT-PK-LABEL` | M1 · `judge_operator.go` | `` return approval.FormatWebAuthn, a, err `` → `` return approval.FormatSSH, a, err `` | `TestOperatorJudgeBothKeys` | killed (TestOperatorJudgeBothKeys) | killed (TestOperatorJudgeBothKeys,TestOperatorJudgeReadsTheStore,TestOperatorRunnerRechecks) | K |
| `PARSE-STEP` | M1 · `judge_operator.go` | `` &JudgeError{Step: JudgeStepParse, Err: errors.Join( `` → `` &JudgeError{Step: JudgeStepLaw, Err: errors.Join( `` | `TestOperatorJudgeReadsTheStore` | killed (TestOperatorJudgeReadsTheStore) | killed (TestOperatorJudgeReadsTheStore,TestOperatorRunnerRechecks) | K |
| `PARSE-SENTINEL` | M1 · `judge_operator.go` | `` errors.Join(ErrUnknownArtifact, sigErr, respErr) `` → `` errors.Join(sigErr, respErr) `` | `TestOperatorJudgeReadsTheStore` | killed (TestOperatorJudgeReadsTheStore) | killed (TestOperatorJudgeReadsTheStore,TestOperatorRunnerRechecks) | K |
| `J-APPROVER-ERR` | M1 · `judge_operator.go` | `` if err != nil {⏎return nil, err⏎}⏎s, err := approval.ConfirmOperatorA… `` → `` if false {⏎return nil, err⏎}⏎s, err := approval.ConfirmOperatorAction `` | `TestOperatorJudgeReadsTheStore` | killed (TestOperatorJudgeReadsTheStore) | killed (TestOperatorJudgeReadsTheStore) | K |
| `J-ACTION-WIRE` | M1 · `judge_operator.go` | `` OperatorConfirmInput{Action: in.Action, `` → `` OperatorConfirmInput{Action: approval.OperatorAction{}, `` | `TestOperatorJudgeBothKeys` | killed (TestOperatorJudgeBothKeys) | killed (TestOperatorJudgeBothKeys,TestOperatorJudgeReadsTheStore) | K |
| `J-CHALLENGE-WIRE` | M1 · `judge_operator.go` | `` Admit: in.admit(), Challenge: in.Challenge, `` → `` Admit: in.admit(), Challenge: [32]byte{}, `` | `TestOperatorJudgeBothKeys` | killed (TestOperatorJudgeBothKeys) | killed (TestOperatorJudgeBothKeys,TestOperatorJudgeReadsTheStore) | K |
| `J-FORMAT-WIRE` | M1 · `judge_operator.go` | `` in.Challenge,⏎Artifact: artifact, Format: format, `` → `` in.Challenge,⏎Artifact: artifact, Format: format[:0] + approval.Forma… `` | `TestOperatorJudgeBothKeys` | killed (TestOperatorJudgeBothKeys) | killed (TestOperatorJudgeBothKeys,TestOperatorJudgeReadsTheStore) | K |
| `J-APPROVER-WIRE` | M1 · `judge_operator.go` | `` in.Challenge,⏎Artifact: artifact, Format: format, Approver: a} `` → `` in.Challenge,⏎Artifact: artifact, Format: format, Approver: approval.… `` | `TestOperatorJudgeBothKeys` | killed (TestOperatorJudgeBothKeys) | killed (TestOperatorJudgeBothKeys,TestOperatorJudgeReadsTheStore) | K |
| `ADMIT-NOW` | M1 · `judge_operator.go` | `` approval.AdmitInput{Now: in.Now, `` → `` approval.AdmitInput{Now: in.CreatedAt, `` | `TestOperatorJudgeReadsTheStore` | killed (TestOperatorJudgeReadsTheStore) | killed (TestOperatorJudgeReadsTheStore) | K |
| `ADMIT-CREATED` | M1 · `judge_operator.go` | `` CreatedAt: in.CreatedAt, ExpiresAt: in.ExpiresAt, Consumed `` → `` CreatedAt: in.Now, ExpiresAt: in.ExpiresAt, Consumed `` | `TestOperatorJudgeReadsTheStore` | killed (TestOperatorJudgeReadsTheStore) | killed (TestOperatorJudgeReadsTheStore) | K |
| `ADMIT-EXPIRES` | M1 · `judge_operator.go` | `` ExpiresAt: in.ExpiresAt, Consumed: in.Consumed, `` → `` ExpiresAt: in.Now + 1, Consumed: in.Consumed, `` | `TestOperatorJudgeBothKeys` | killed (TestOperatorJudgeBothKeys) | killed (TestOperatorJudgeBothKeys,TestOperatorJudgeReadsTheStore) | K |
| `ADMIT-CONSUMED` | M1 · `judge_operator.go` | `` Consumed: in.Consumed, `` → `` Consumed: false, `` | `TestOperatorJudgeReadsTheStore` | killed (TestOperatorJudgeReadsTheStore) | killed (TestOperatorJudgeReadsTheStore) | K |
| `ADMIT-STORED` | M1 · `judge_operator.go` | `` StoredDigest: in.Digest, `` → `` StoredDigest: in.BytesDigest, `` | `TestOperatorJudgeReadsTheStore` | killed (TestOperatorJudgeReadsTheStore) | killed (TestOperatorJudgeReadsTheStore) | K |
| `ADMIT-BYTES` | M1 · `judge_operator.go` | `` BytesDigest: in.BytesDigest} `` → `` BytesDigest: in.Digest} `` | `TestOperatorJudgeReadsTheStore` | killed (TestOperatorJudgeReadsTheStore) | killed (TestOperatorJudgeReadsTheStore) | K |
| `J-LAW-STEP` | M1 · `judge_operator.go` | `` return nil, &JudgeError{Step: JudgeStepLaw, Err: err}⏎}⏎return verifi… `` → `` return nil, &JudgeError{Step: JudgeStepParse, Err: err}⏎}⏎return veri… `` | `TestOperatorJudgeReadsTheStore` | killed (TestOperatorJudgeReadsTheStore) | killed (TestOperatorJudgeReadsTheStore) | K |
| `J-RECHECKED` | M1 · `judge_operator.go` | `` verifiedOperatorAction(s, false) `` → `` verifiedOperatorAction(s, true) `` | `TestOperatorJudgeBothKeys` | killed (TestOperatorJudgeBothKeys) | killed (TestOperatorJudgeBothKeys) | K |
| `R-APPROVER-ERR` | M1 · `judge_operator.go` | `` if err != nil {⏎return nil, err⏎}⏎s, err := approval.CheckOperatorAction `` → `` if false {⏎return nil, err⏎}⏎s, err := approval.CheckOperatorAction `` | `TestOperatorRunnerRechecks` | killed (TestOperatorRunnerRechecks) | killed (TestOperatorRunnerRechecks) | K |
| `R-ACTION-WIRE` | M1 · `judge_operator.go` | `` OperatorCheckInput{Action: action, `` → `` OperatorCheckInput{Action: approval.OperatorAction{}, `` | `TestOperatorRunnerRechecks` | killed (TestOperatorRunnerRechecks) | killed (TestOperatorRunnerRechecks) | K |
| `R-CHALLENGE-WIRE` | M1 · `judge_operator.go` | `` Action: action, Challenge: challenge, `` → `` Action: action, Challenge: [32]byte{}, `` | `TestOperatorRunnerRechecks` | killed (TestOperatorRunnerRechecks) | killed (TestOperatorRunnerRechecks) | K |
| `R-EXPIRY-WIRE` | M1 · `judge_operator.go` | `` Expiry: expiry, Now: now, `` → `` Expiry: now + 1, Now: now, `` | `TestOperatorRunnerRechecks` | killed (TestOperatorRunnerRechecks) | killed (TestOperatorRunnerRechecks) | K |
| `R-NOW-WIRE` | M1 · `judge_operator.go` | `` Expiry: expiry, Now: now, `` → `` Expiry: expiry, Now: 0, `` | `TestOperatorRunnerRechecks` | killed (TestOperatorRunnerRechecks) | killed (TestOperatorRunnerRechecks) | K |
| `R-FORMAT-WIRE` | M1 · `judge_operator.go` | `` Artifact: artifact, Format: format, Approver: a})⏎if err != nil {⏎ret… `` → `` Artifact: artifact, Format: format[:0] + approval.FormatSSH, Approver… `` | `TestOperatorRunnerRechecks` | killed (TestOperatorRunnerRechecks) | killed (TestOperatorRunnerRechecks) | K |
| `R-APPROVER-WIRE` | M1 · `judge_operator.go` | `` Artifact: artifact, Format: format, Approver: a})⏎if err != nil {⏎ret… `` → `` Artifact: artifact, Format: format, Approver: approval.Approver{Name:… `` | `TestOperatorRunnerRechecks` | killed (TestOperatorRunnerRechecks) | killed (TestOperatorRunnerRechecks) | K |
| `R-LAW-STEP` | M1 · `judge_operator.go` | `` return nil, &JudgeError{Step: JudgeStepLaw, Err: err}⏎}⏎if s.Approved… `` → `` return nil, &JudgeError{Step: JudgeStepParse, Err: err}⏎}⏎if s.Approv… `` | `TestOperatorRunnerRechecks` | killed (TestOperatorRunnerRechecks) | killed (TestOperatorRunnerRechecks) | K |
| `R-IDENTITY-OFF` | M1 · `judge_operator.go` | `` if s.ApprovedBy() != approvedBy { `` → `` if false { `` | `TestOperatorRunnerRechecks` | killed (TestOperatorRunnerRechecks) | killed (TestOperatorRunnerRechecks) | K |
| `R-IDENTITY-SENTINEL` | M1 · `judge_operator.go` | `` Err: ErrOperatorApprover} `` → `` Err: approval.ErrBadSignature} `` | `TestOperatorRunnerRechecks` | killed (TestOperatorRunnerRechecks) | killed (TestOperatorRunnerRechecks) | K |
| `R-IDENTITY-STEP` | M1 · `judge_operator.go` | `` &JudgeError{Step: JudgeStepLaw, Err: ErrOperatorApprover} `` → `` &JudgeError{Step: JudgeStepParse, Err: ErrOperatorApprover} `` | `TestOperatorRunnerRechecks` | killed (TestOperatorRunnerRechecks) | killed (TestOperatorRunnerRechecks) | K |
| `R-RECHECKED` | M1 · `judge_operator.go` | `` verifiedOperatorAction(s, true) `` → `` verifiedOperatorAction(s, false) `` | `TestOperatorRunnerRechecks` | killed (TestOperatorRunnerRechecks) | killed (TestOperatorRunnerRechecks) | K |
| `CTOR-SIGNED` | M1 · `judge_operator.go` | `` &VerifiedOperatorAction{signed: s, rechecked: rechecked} `` → `` &VerifiedOperatorAction{signed: nil, rechecked: rechecked} `` | `TestOperatorJudgeBothKeys` | killed (TestOperatorJudgeBothKeys) | killed (TestOperatorJudgeBothKeys,TestOperatorRunnerRechecks) | K |
| `CTOR-FLAG` | M1 · `judge_operator.go` | `` &VerifiedOperatorAction{signed: s, rechecked: rechecked} `` → `` &VerifiedOperatorAction{signed: s, rechecked: true} `` | `TestOperatorJudgeBothKeys` | killed (TestOperatorJudgeBothKeys) | killed (TestOperatorJudgeBothKeys) | K |
| `BINDS-NIL` | M1 · `judge_operator.go` | `` if v == nil {⏎return approval.ErrNotVerified `` → `` if false {⏎return approval.ErrNotVerified `` | `TestOperatorUnjudgedProofBindsNothing` | killed (TestOperatorUnjudgedProofBindsNothing) | killed (TestOperatorUnjudgedProofBindsNothing) | K |
| `BINDS-DELEGATE` | M1 · `judge_operator.go` | `` return v.signed.Binds(a, challenge, expiry) `` → `` return nil `` | `TestOperatorJudgeBothKeys` | killed (TestOperatorJudgeBothKeys) | killed (TestOperatorJudgeBothKeys,TestOperatorUnjudgedProofBindsNothing) | K |
| `APPROVEDBY-NIL` | M1 · `judge_operator.go` | `` if v == nil {⏎return "" `` → `` if false {⏎return "" `` | `TestOperatorUnjudgedProofBindsNothing` | killed (TestOperatorUnjudgedProofBindsNothing) | killed (TestOperatorUnjudgedProofBindsNothing) | K |
| `APPROVEDBY-DELEGATE` | M1 · `judge_operator.go` | `` return v.signed.ApprovedBy() `` → `` return "ssh:mark-ssh" `` | `TestOperatorJudgeBothKeys` | killed (TestOperatorJudgeBothKeys) | killed (TestOperatorJudgeBothKeys,TestOperatorRunnerRechecks,TestOperatorUnjudgedProofBindsNothing) | K |
| `PASSKEYUSE-NIL` | M1 · `judge_operator.go` | `` if v == nil {⏎return nil, 0, false `` → `` if false {⏎return nil, 0, false `` | `TestOperatorUnjudgedProofBindsNothing` | killed (TestOperatorUnjudgedProofBindsNothing) | killed (TestOperatorUnjudgedProofBindsNothing) | K |
| `PASSKEYUSE-DELEGATE` | M1 · `judge_operator.go` | `` return v.signed.PasskeyUse() `` → `` return nil, 0, false `` | `TestOperatorJudgeBothKeys` | killed (TestOperatorJudgeBothKeys) | killed (TestOperatorJudgeBothKeys) | K |
| `RECHECKED-NIL` | M1 · `judge_operator.go` | `` v != nil && v.signed != nil && v.rechecked `` → `` v.signed != nil && v.rechecked `` | `TestOperatorUnjudgedProofBindsNothing` | killed (TestOperatorUnjudgedProofBindsNothing) | killed (TestOperatorUnjudgedProofBindsNothing) | K |
| `RECHECKED-SIGNED` | M1 · `judge_operator.go` | `` v != nil && v.signed != nil && v.rechecked `` → `` v != nil && v.rechecked `` | `TestOperatorUnjudgedProofBindsNothing` | killed (TestOperatorUnjudgedProofBindsNothing) | killed (TestOperatorUnjudgedProofBindsNothing) | K |
| `RECHECKED-FLAG` | M1 · `judge_operator.go` | `` v != nil && v.signed != nil && v.rechecked `` → `` v != nil && v.signed != nil `` | `TestOperatorJudgeBothKeys` | killed (TestOperatorJudgeBothKeys) | killed (TestOperatorJudgeBothKeys) | K |
| `WN-LOG-INJECT` | M3 · `judge_operator.go` | `` "errors"⏎ `` → `` "errors"⏎"log"⏎ `` + `` return verifiedOperatorAction(s, false), nil `` → `` log.Print("judged")⏎return verifiedOperatorAction(s, false), nil `` | `TestOperatorJudgesWriteNothing` | killed (TestOperatorJudgesWriteNothing) | killed (TestOperatorJudgesWriteNothing) | K |
| `REUSED-SSH-STORE-ERR` | — · `judge_verb.go` (reused) | `` ro.ApproverByKey(ctx, sig.PublicKey.Blob)⏎if err != nil { `` → `` ro.ApproverByKey(ctx, sig.PublicKey.Blob)⏎if false { `` | `TestOperatorJudgeReadsTheStore` | killed (TestOperatorJudgeReadsTheStore) | killed (TestJudgeVerbSSH,TestOperatorJudgeReadsTheStore,TestRecheckVerb) | K |
| `REUSED-SSH-UNKNOWN` | — · `judge_verb.go` (reused) | `` if !ok {⏎return approval.Approver{}, &JudgeError{Step: JudgeStepUnkno… `` → `` if !ok && false {⏎return approval.Approver{}, &JudgeError{Step: Judge… `` | `TestOperatorJudgeReadsTheStore` | killed (TestOperatorJudgeReadsTheStore) | killed (TestJudgeVerbSSH,TestOperatorJudgeReadsTheStore,TestOperatorRunnerRechecks,TestRecheckVerb) | K |
| `REUSED-SSH-NAME` | — · `judge_verb.go` (reused) | `` return approval.Approver{Name: enrolled.Name, `` → `` return approval.Approver{Name: "mark-ssh", `` | `TestOperatorJudgeBothKeys` | killed (TestOperatorJudgeBothKeys) | killed (TestOperatorJudgeBothKeys) | K |
| `REUSED-SSH-FORMAT` | — · `judge_verb.go` (reused) | `` Format: enrolled.Format, `` → `` Format: approval.FormatSSH, `` | `TestOperatorJudgeReadsTheStore` | killed (TestOperatorJudgeReadsTheStore) | killed (TestOperatorJudgeReadsTheStore) | K |
| `REUSED-SSH-KEY` | — · `judge_verb.go` (reused) | `` Key: approval.SSHPublicKey{Blob: enrolled.PublicKey} `` → `` Key: approval.SSHPublicKey{Blob: sig.PublicKey.Blob} `` | `TestOperatorJudgeBothKeys` | survived | survived — **equivalent**: `ApproverByKey` selects `public_key = ?` with exactly `sig.PublicKey.Blob`, so the stored key and the signature's key are the same bytes on every row it returns | E |
| `REUSED-SSH-ENROLLED` | — · `judge_verb.go` (reused) | `` EnrolledAt: enrolled.EnrolledAt, `` → `` EnrolledAt: 0, `` | `TestOperatorJudgeReadsTheStore` | killed (TestOperatorJudgeReadsTheStore) | killed (TestOperatorJudgeReadsTheStore) | K |
| `REUSED-SSH-REVOKED` | — · `judge_verb.go` (reused) | `` RevokedAt: enrolled.RevokedAt} `` → `` RevokedAt: 0} `` | `TestOperatorJudgeReadsTheStore` | killed (TestOperatorJudgeReadsTheStore) | killed (TestJudgeVerbSSH,TestOperatorJudgeReadsTheStore,TestOperatorRunnerRechecks,TestRecheckVerb) | K |
| `REUSED-PK-STORE-ERR` | — · `judge_verb.go` (reused) | `` ro.PasskeyByCredentialID(ctx, resp.RawID)⏎if err != nil { `` → `` ro.PasskeyByCredentialID(ctx, resp.RawID)⏎if false { `` | `TestOperatorJudgeReadsTheStore` | killed (TestOperatorJudgeReadsTheStore) | killed (TestJudgeVerbAssertion,TestOperatorJudgeReadsTheStore,TestRecheckVerb) | K |
| `REUSED-PK-UNKNOWN` | — · `judge_verb.go` (reused) | `` if !ok {⏎return approval.Approver{}, &JudgeError{Step: JudgeStepUnkno… `` → `` if !ok && false {⏎return approval.Approver{}, &JudgeError{Step: Judge… `` | `TestOperatorJudgeReadsTheStore` | killed (TestOperatorJudgeReadsTheStore) | killed (TestJudgeVerbAssertion,TestOperatorJudgeReadsTheStore,TestRecheckVerb) | K |
| `REUSED-PK-STOREDKEY` | — · `judge_verb.go` (reused) | `` approval.ParseStoredCOSEKey(a.PublicKey)⏎if err != nil { `` → `` approval.ParseStoredCOSEKey(a.PublicKey)⏎if false { `` | `TestOperatorJudgeReadsTheStore` | killed (TestOperatorJudgeReadsTheStore) | killed (TestJudgeVerbAssertion,TestOperatorJudgeReadsTheStore,TestRecheckVerb) | K |
| `REUSED-PK-FORMAT` | — · `judge_verb.go` (reused) | `` return approval.Approver{Name: a.Name, Format: a.Format, `` → `` return approval.Approver{Name: a.Name, Format: approval.FormatWebAuthn, `` | `TestOperatorJudgeReadsTheStore` | killed (TestOperatorJudgeReadsTheStore) | killed (TestOperatorJudgeReadsTheStore) | K |
| `REUSED-PK-ENROLLED` | — · `judge_verb.go` (reused) | `` Format: a.Format, EnrolledAt: a.EnrolledAt, `` → `` Format: a.Format, EnrolledAt: 0, `` | `TestOperatorJudgeReadsTheStore` | killed (TestOperatorJudgeReadsTheStore) | killed (TestOperatorJudgeReadsTheStore) | K |
| `REUSED-PK-REVOKED` | — · `judge_verb.go` (reused) | `` RevokedAt: a.RevokedAt,⏎ `` → `` RevokedAt: 0,⏎ `` | `TestOperatorJudgeReadsTheStore` | killed (TestOperatorJudgeReadsTheStore) | killed (TestOperatorJudgeReadsTheStore,TestRecheckVerb) | K |
| `REUSED-PK-COUNTER` | — · `judge_verb.go` (reused) | `` SignCount: uint32(p.SignCount)} `` → `` SignCount: 0} `` | `TestOperatorJudgeReadsTheStore` | killed (TestOperatorJudgeReadsTheStore) | killed (TestJudgeVerbAssertion,TestOperatorJudgeReadsTheStore,TestOperatorRunnerRechecks) | K |
| `REUSED-PK-CREDID` | — · `judge_verb.go` (reused) | `` CredentialID: p.CredentialID, Key: key, `` → `` CredentialID: resp.RawID, Key: key, `` | `TestOperatorJudgeBothKeys` | survived | survived — **equivalent**: `PasskeyByCredentialID` selects `credential_id = ?` with exactly `resp.RawID`, so the stored credential ID and the assertion's raw ID are the same bytes on every row it returns | E |
| `X-PK-ALIAS` | M2 → S2 · `operator_signed.go` | `` s.credentialID, s.signCount = bytes.Clone(p.CredentialID), ad.SignCount `` → `` s.credentialID, s.signCount = p.CredentialID, ad.SignCount `` | `TestOperatorPasskeyUseIsACopy` | killed (TestOperatorPasskeyUseIsACopy) | killed (TestOperatorPasskeyUseIsACopy) | K |
| `PK-USE-ALIAS` | M2 → S2 · `operator_signed.go` | `` return bytes.Clone(s.credentialID), s.signCount, true `` → `` return s.credentialID, s.signCount, true `` | `TestOperatorPasskeyUseIsACopy` | killed (TestOperatorPasskeyUseIsACopy) | killed (TestOperatorPasskeyUseIsACopy) | K |
| `R1-LIFETIME-ADDED` | M2 → S2 · `operator_signed.go` | `` if in.Now >= in.Expiry { // R1 `` → `` if in.Now >= in.Expiry \|\| in.Expiry-in.Now > int64(OperatorActionLi… `` | `TestOperatorRunnerRecheckBoundsExpiryNotLifetime` | killed (TestOperatorRunnerRecheckBoundsExpiryNotLifetime) | killed (TestOperatorRunnerRecheckBoundsExpiryNotLifetime) | K |

**Totals: 61 mutants (0 NM, 0 NM-text in run 2), 59 killed by their named test and by the whole package, 2 named equivalents, 0
survivors, 0 pending.** No kill rests on the static test alone (the `w/o static` leg never triggered in run 2). Every
store-sourced check the directive lists is its own row: unknown key (`REUSED-SSH-UNKNOWN`), unknown passkey credential
(`REUSED-PK-UNKNOWN`), revoked (`REUSED-SSH-REVOKED`, `REUSED-PK-REVOKED`), enrolment (`REUSED-SSH-ENROLLED`, `REUSED-PK-ENROLLED`),
format dispatch (`FMT-*`) and stored format (`REUSED-SSH-FORMAT`, `REUSED-PK-FORMAT`), the counter (`REUSED-PK-COUNTER`), the
approver-identity comparison (`R-IDENTITY-*`), the `rechecked` flag (`J-RECHECKED`, `R-RECHECKED`, `CTOR-FLAG`, `RECHECKED-*`), each
`Binds` delegation (`BINDS-*`; S2's `Binds` fields are S2's rows) and each accessor's nil/zero guard (`*-NIL`).

---

## 5. Files (final `git diff --name-status 4b1c39e2 HEAD`)

```
A	host/approval/operator_pins_test.go
A	host/authority/judge_operator.go
A	host/authority/judge_operator_silent_test.go
A	host/authority/judge_operator_test.go
```

One new production file; no pre-existing file of any kind; nothing deleted; no `cmd`, `docs`, `app.js`, `daemon.go`,
`surface_page.go`, `host/store`, host/approval production, `.ail`, `world`, `packages`, `scripts`, `tools` or `plugins` file.
**Pre-existing test files changed: none. Deleted: none.**

---

## 6. Sprint gates (controller, after the last sub-milestone's commit; each measured at the base, §1.1)

`bash $HOME/.ailang/state/world-iter409/gates.sh <worktree> $HOME/.ailang/state/world-iter409/gates final <HEAD sha>` (saved there
byte-equal to this block, C16; one bash shell; `GATES_SKIP_VERIFY_GO=1` skips `verify_go` when the controller runs it separately —
`VERIFY_GO_LEGS=plain AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_go.sh` after commit). Each leg is its own subshell with
its own rc, `tee`'d to `<logdir>/<leg>.txt`; no leg is chained to another; `ALL_DONE.<sha>` holds `FAILS: …`. `TOUCHED_PACKAGES` is
derived from `git diff --name-only` and must equal `host/approval host/authority`. `-race` runs on both touched packages whole
and on approvaltest. S2b adds no production context root, and the `inv` leg runs `go test ./host/store/ -run
'^TestProductionContextRoots$'` (with the rest of `INV_STORE`) regardless. S2b has no Chrome or node test: no `chrome` or `node`
leg. Never start it under `nohup`.

````bash gates.sh
#!/bin/bash
# Row 243 V4 S2b gates (plan §6). usage: gates.sh <tree> <logdir> <base|final> [<sha>]
# The controller: bash gates.sh <worktree> $HOME/.ailang/state/world-iter409/gates final <HEAD sha>
# macOS /bin/bash 3.2: no declare -A. Every leg records its own rc; every leg runs (no set -e, no leg
# chained to another with &&); every leg is tee'd to <logdir>/<leg>.txt. GATES_SKIP_VERIFY_GO=1 skips
# verify_go (the executor sets it; the controller runs it). Never start this under nohup. Writes
# <logdir>/ALL_DONE.<sha> with FAILS: when every leg has finished. S2b has no Chrome test: no chrome leg.
T=$1; L=$2; MODE=$3
mkdir -p "$L"; cd "$T" || exit 2
SHA=${4:-$(git rev-parse HEAD 2>/dev/null || echo unknown)}
B=4b1c39e2a32fa129f2b340f4f9c2b6a312c63e0f
export AILANG_BIN=$HOME/.pinned-ailang/ailang
export WORLD_EXEC_NODE=/opt/homebrew/bin/node
export WORLD_TOOL_AILANG_BIN=$HOME/.pinned-ailang-tools/v0.54.15/ailang
unset WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER WORLD_PASSKEY_RECORD_DIR WORLD_UPDATE_OPERATOR_ACTION_GOLDEN WORLD_UPDATE_VERB_GOLDEN WORLD_UPDATE_WHY_GOLDEN WORLD_UPDATE_FINISH_PROMPT_GOLDEN WORLD_UPDATE_CONFIRM_GOLDEN WORLD_UPDATE_CONFIRMATIONS_GOLDEN
export WORLD_REQUIRE_SSH_KEYGEN=1
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
# The touched Go packages, derived from the diff (printed so the controller can diff it). The plan's list is
# host/approval host/authority; a difference FAILS.
TOUCHED=$(git diff --name-only "$B" HEAD -- cmd host 2>/dev/null | grep '\.go$' | xargs -n1 dirname 2>/dev/null | sort -u | tr '\n' ' ')
echo "TOUCHED_PACKAGES: $TOUCHED" | tee -a "$L/legs.txt"
if [ "$MODE" = final ] && [ "$TOUCHED" != "host/approval host/authority " ]; then FAILS="$FAILS touched_packages_differ"; fi
leg status 'git status --porcelain | wc -l'
leg verify_ail './scripts/verify_ail.sh'
leg gofmt 'out=$(gofmt -l cmd host); echo "$out"; test -z "$out"'
leg build 'go build ./...'
leg vet 'go vet ./...'
leg fence "go test -run '^\$' ./... > '$L/fence.raw' 2>&1; rc=\$?; cat '$L/fence.raw'; if grep -E '^(FAIL|--- FAIL)' '$L/fence.raw' >/dev/null; then exit 1; fi; exit \$rc"
# the two touched packages whole, and approvaltest (the keys every new test uses), each its own rc
leg pkgs "r=0; for p in ./host/approval/ ./host/approval/approvaltest/ ./host/authority/; do go test -count=1 -timeout 30m \$p; r=\$(( r | \$? )); done; exit \$r"
# -race on every touched package whole (and approvaltest)
leg race "r=0; for p in ./host/approval/ ./host/approval/approvaltest/ ./host/authority/; do go test -race -count=1 -timeout 30m \$p; r=\$(( r | \$? )); done; exit \$r"
NEW_AU="TestOperatorJudgeBothKeys TestOperatorJudgeReadsTheStore TestOperatorRunnerRechecks TestOperatorUnjudgedProofBindsNothing TestVerifiedOperatorActionOnlyFromTheJudges TestOperatorJudgesWriteNothing"
NEW_A="TestOperatorPasskeyUseIsACopy TestOperatorRunnerRecheckBoundsExpiryNotLifetime"
if [ "$MODE" = base ]; then NEW_AU=""; NEW_A=""; fi
INV_APPROVAL="TestApprovalImportsOnlyTheStandardLibrary TestSSHSIGCrossCheckWithSSHKeygen TestVerbDomainSeparation TestOperatorConfirmBothKeys TestOperatorRunnerRecheckLaw TestOperatorPureEntryPointsWriteNothing"
INV_APPROVALTEST="TestApprovalTestIsTestOnly TestTestSignerAgreesWithSSHKeygen"
INV_AUTHORITY="TestJudgeVerbSSH TestJudgeVerbAssertion TestRecheckVerb TestReadOnlyAtFailuresAreErrors"
INV_DAEMON="TestDaemonDependencyAllowlist TestFrozenV1RouteTableMatchesMux TestSurfaceAppRoutes TestApprovalRoutesRegistered TestDecisionRoutesRegistered TestConfirmRoutesRegistered TestDaemonCredentialWriterSites TestSurfaceVerbRouteRegistered TestSurfaceVerbIsReadOnly TestSurfaceVerbOutcomes"
INV_BOUNDARY="TestWorldBoundaryDependencyAllowlist TestBareNetHTTPExemptionIsPerGroup TestEveryCommandDirectoryIsAProtectedGroup TestBoundaryASTWriteGuard"
INV_STORE="TestProductionContextRoots TestProductionGoSurface TestStoreSurfaceTakesContext TestCredentialsSchemaMatchesCanonicalDDL TestCredentialsSchema9VerbRequests TestCredentialsNewerThanNineIsRefused"
INV_CMD="TestEveryVerbHasHelp TestVerbTableMatchesTheCLISwitch TestCLICredentialWriterSites TestOperatorFlagParity"
leg inv "set -o pipefail; r=0
go test -v -count=1 ./host/approval/ -run '^($(bar $INV_APPROVAL)$(alt $NEW_A))\$'; r=\$(( r | \$? ))
go test -v -count=1 ./host/approval/approvaltest/ -run '^($(bar $INV_APPROVALTEST))\$'; r=\$(( r | \$? ))
go test -v -count=1 ./host/authority/ -run '^($(bar $INV_AUTHORITY)$(alt $NEW_AU))\$'; r=\$(( r | \$? ))
go test -v -count=1 ./host/daemon/ -run '^($(bar $INV_DAEMON))\$'; r=\$(( r | \$? ))
go test -v -count=1 ./host/boundary/ -run '^($(bar $INV_BOUNDARY))\$'; r=\$(( r | \$? ))
go test -v -count=1 ./host/store/ -run '^($(bar $INV_STORE))\$'; r=\$(( r | \$? ))
go test -v -count=1 ./cmd/ailang-worldd/ -run '^($(bar $INV_CMD))\$'; r=\$(( r | \$? ))
exit \$r"
need inv $INV_APPROVAL $NEW_A $INV_APPROVALTEST $INV_AUTHORITY $NEW_AU $INV_DAEMON $INV_BOUNDARY $INV_STORE $INV_CMD
if grep -E '^\s*--- (SKIP|FAIL)' "$L/inv.txt"; then FAILS="$FAILS inv_skip_or_fail"; fi
echo "named PASS: $(grep -cE '^--- PASS' "$L/inv.txt")" | tee -a "$L/legs.txt"
# No new daemon route: every mux registration in daemon.go, counted (base = final), and the /surface POSTs (3)
ROUTES=$(grep -cE 'mux\.(HandleFunc|Handle)\(' host/daemon/daemon.go); POSTS=$(grep -cE 'mux\.HandleFunc\("POST /surface/' host/daemon/daemon.go)
BROUTES=$(git show "$B:host/daemon/daemon.go" | grep -cE 'mux\.(HandleFunc|Handle)\(')
echo "DAEMON_ROUTES: $ROUTES (base $BROUTES) SURFACE_POSTS: $POSTS" | tee -a "$L/legs.txt"
if [ "$ROUTES" != "$BROUTES" ] || [ "$POSTS" != 3 ]; then FAILS="$FAILS routes_changed"; fi
# Production bytes: S2b's only production file is host/authority/judge_operator.go (NEW). Every other changed
# file is a _test.go or a plan file; no non-test file under cmd, docs, host/daemon, host/store, host/approval,
# host/surfacepage, world, packages, scripts, tools, plugins; no .ail.
NOPROD=$(git diff --name-only "$B" HEAD | grep -vE '(_test\.go$|^host/authority/judge_operator\.go$|^design_docs/planned/(w-world-address-bar-verbs-v4s2b-sprint-plan\.md|sprint_w-world-address-bar-verbs-v4s2b\.json)$)' | tr '\n' ' ')
NOPROD="$NOPROD$(git diff --name-only "$B" HEAD -- cmd docs host/daemon host/store host/approval host/surfacepage world packages scripts tools plugins '*.ail' | grep -v '_test\.go$' | tr '\n' ' ')"
echo "UNTOUCHED_VIOLATIONS: ${NOPROD:-none}" | tee -a "$L/legs.txt"
if [ "$MODE" = final ] && [ -n "$NOPROD" ]; then FAILS="$FAILS untouched_violated"; fi
# Pre-existing files changed or deleted: none (every S2b path is added).
PRE=$(git diff --name-status "$B" HEAD | awk '$1 != "A"' | tr '\n' ' ')
echo "PREEXISTING_CHANGED: ${PRE:-none}" | tee -a "$L/legs.txt"
if [ "$MODE" = final ] && [ -n "$PRE" ]; then FAILS="$FAILS preexisting_changed"; fi
# The judges have no production caller, and the S2 envelope/law's only production caller is the judges' file
# (caller searches with a positive control in the same call).
JCALLERS=$(grep -rnE '\b(JudgeOperatorAction|RecheckOperatorAction|VerifiedOperatorAction|OperatorJudgeInput|ErrOperatorApprover|verifiedOperatorAction|operatorApprover)\b' --include='*.go' cmd host | grep -v '_test\.go:' | grep -vE '^host/authority/judge_operator\.go:' | tr '\n' ' ')
JCONTROL=$(grep -rnE '\bauthority\.RecheckVerb\(' --include='*.go' cmd host | grep -v '_test\.go:' | wc -l | tr -d ' ')
echo "JUDGE_PRODUCTION_CALLERS: ${JCALLERS:-none} (control: authority.RecheckVerb production callers $JCONTROL)" | tee -a "$L/legs.txt"
if [ "$MODE" = final ] && { [ -n "$JCALLERS" ] || [ "$JCONTROL" -lt 1 ]; }; then FAILS="$FAILS judges_have_a_caller"; fi
CALLERS=$(grep -rnE '\b(OperatorAction|OperatorActionSchema|OperatorActionPrefix|OperatorActionNamespace|OperatorActionMessage|OperatorActionMessageHash|OperatorReviewDigest|OperatorSignatureDigest|DecodeOperatorAction|CanonicalOperatorAction|OperatorArgv|OperatorOptions|OperatorResolved|ConfirmOperatorAction|CheckOperatorAction|SignedOperatorAction|OperatorConfirmInput|OperatorCheckInput|MaxOperatorActionBytes|MaxOperatorTaskBytes|OperatorActionLifetime)\b' --include='*.go' cmd host | grep -v '_test\.go:' | grep -vE '^host/approval/operator_(action|signed)\.go:|^host/authority/judge_operator\.go:' | tr '\n' ' ')
CONTROL=$(grep -rnE '\bapproval\.CheckConfirmation\(' --include='*.go' cmd host | grep -v '_test\.go:' | wc -l | tr -d ' ')
echo "ENVELOPE_PRODUCTION_CALLERS_OUTSIDE_THE_JUDGES: ${CALLERS:-none} (control: approval.CheckConfirmation production callers $CONTROL)" | tee -a "$L/legs.txt"
if [ "$MODE" = final ] && { [ -n "$CALLERS" ] || [ "$CONTROL" -lt 1 ]; }; then FAILS="$FAILS envelope_has_a_caller"; fi
# Served bytes unchanged: app.js, the CSP literal's file, the route table and the store's schema
for f in host/surfacepage/app/app.js host/daemon/surface_page.go host/daemon/daemon.go host/store/credentials.go; do
  if git diff --quiet "$B" HEAD -- "$f"; then echo "UNCHANGED: $f" | tee -a "$L/legs.txt"; else echo "CHANGED: $f" | tee -a "$L/legs.txt"; FAILS="$FAILS served_bytes_changed"; fi
done
if [ "${GATES_SKIP_VERIFY_GO:-0}" = 1 ]; then echo "LEG verify_go SKIPPED (GATES_SKIP_VERIFY_GO=1)" | tee -a "$L/legs.txt"
else leg verify_go 'VERIFY_GO_LEGS=plain AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_go.sh'; fi
echo "GATES_DONE $(( $(date +%s) - s0 ))s FAILS:${FAILS:- none}" | tee -a "$L/legs.txt"
echo "GATES_DONE $(( $(date +%s) - s0 ))s FAILS:${FAILS:- none}" > "$L/ALL_DONE.$SHA"
````

---

## 7. Executor directive notes

- **No git writes of any kind** (no add, rm, mv, commit, stash, checkout, reset, branch). **`patch -p1 -E` (§2) is the
  only writer** of the tree; no hand edit, no `gofmt -w`, no golden regeneration, no fuzz. A failed `--dry-run` or a red
  acceptance → **stop and report** with the log; do not repair.
- **Never `nohup`, never background a test leg; pass a 600000 ms timeout to every long Bash call.** Each acceptance part is
  one call (C14 measured each under a minute). **The executor never launches Chrome** and never sets `WORLD_CHROME`.
- **Runs**: **one run, M1 → M2 → M3** (each patch assumes the previous tree; the controller may commit between
  sub-milestones). One part each: `fast`. `accept.sh` exports `AILANG_BIN=$HOME/.pinned-ailang/ailang`,
  `WORLD_EXEC_NODE=/opt/homebrew/bin/node` and `WORLD_REQUIRE_SSH_KEYGEN=1`, and unsets `WORLD_CHROME`,
  `WORLD_REQUIRE_PASSKEY_BROWSER` and every golden-update variable.
- **Commit-time test commands** (`commit_ms.sh <M> <msgfile> <test-cmd…>`; each carries its env through `accept.sh`, which exports
  `AILANG_BIN` itself; with `T=$HOME/.ailang/state/world-iter409/exec-patches W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter409-row243-v4s2b`):
  **M1** `bash $T/accept.sh $W M1 fast`; **M2** `bash $T/accept.sh $W M2 fast`; **M3** `bash $T/accept.sh $W M3 fast`.
  `commit_ms.sh` vets the directories of the changed `.go` files only (`host/authority`, and `host/approval` from M2).
- After each sub-milestone's acceptance, snapshot **cumulatively against the base** to
  `$HOME/.ailang/state/world-iter409/snap/<M>/`: every file added or changed since `4b1c39e2` (repo-relative paths), a
  **`DELETED`** file (one repo path per line; **empty for every S2b stage**), a **`STATUS`** file, and `MANIFEST.sha256`
  (the two plan files, `STATUS` and `DELETED` are not in it). **Never an `rsync` of whole directories.** The `find` takes
  only the top-level dirs that exist (every S2b stage has only `host/`; no stage has `cmd/` or `docs/`), rehearsed on all three
  (C13):
  ```bash
  bash <<'EOF2'
  set -euo pipefail
  M=M1   # the sub-milestone just accepted
  B=4b1c39e2a32fa129f2b340f4f9c2b6a312c63e0f
  W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter409-row243-v4s2b
  D="$HOME/.ailang/state/world-iter409/snap/$M"; rm -rf "${D:?}"; mkdir -p "$D"
  cd "$W"
  { git diff --name-status --no-renames "$B" -- cmd host docs; git ls-files --others --exclude-standard -- cmd host docs | sed 's/^/A\t/'; } | sort -u -t "$(printf '\t')" -k2 > "$D/STATUS"
  awk -F'\t' '$1 != "D" {print $2}' "$D/STATUS" | while read -r p; do mkdir -p "$D/$(dirname "$p")"; cp -p "$p" "$D/$p"; done
  awk -F'\t' '$1 == "D" {print $2}' "$D/STATUS" > "$D/DELETED"
  (cd "$D" && find $(for d in cmd host docs; do [ -d "$d" ] && echo "$d"; done) -type f -print0 | sort -z | xargs -0 shasum -a 256) > "$D/MANIFEST.sha256"
  cat "$D/STATUS"; echo "DELETED: $(wc -l < "$D/DELETED")"; echo "MANIFEST: $(wc -l < "$D/MANIFEST.sha256")"
  EOF2
  ```
  Expected (`STATUS` lines / `DELETED` / MANIFEST): **M1** 2 / 0 / 2; **M2** 3 / 0 / 3; **M3** 4 / 0 / 4. The controller copies each
  snapshot over the worktree, commits it as one commit, and checks the numstat (`git diff -M --numstat 4b1c39e2 HEAD -- cmd host
  docs app | awk '{s+=$1+$2} END {print s}'`): **M1 401, M2 469, M3 526**. A different number → stop.
- **Stop conditions**: a dry-run failure; any `--- FAIL`, `--- SKIP` or `FAIL` line in an acceptance (`ACCEPT_RED`); gofmt
  output; a numstat mismatch; any file outside §5 changed (`git status --short` after each patch must equal the cumulative list
  of §3's paragraphs); a numstat above 1,400.
- **Known load-sensitive tests** (the directive's list) are not this sprint's; none is in an S2b acceptance part (the packages
  run whole are host/approval, approvaltest and host/authority, each a few seconds).
- D-WORLD-109: nothing touches `~/.ailang/world`, `~/.ssh`, a real `/dev/tty` or a real key; every store is a test's `t.TempDir()`.

---

## 8. Risks

- **R1 — the judges are unregistered.** Nothing calls them in production (C8), so a regression in them is invisible to every served
  path until S3/S7 wire them. The mutant table, not a caller, is their guard; S7 must call `JudgeOperatorAction` and the runner
  `RecheckOperatorAction`, never the S2 law directly (Appendix C).
- **R2 — `ErrUnknownArtifact` is shared with `RecheckVerb`.** Both judges classify an unparseable artifact with V3's sentinel; a
  caller distinguishing V3 from V4 refusals must use the step and the route, not the sentinel. Intentional reuse (design §3.1 [r1]
  "Reuse existing key lookup").
- **R3 — `RecheckOperatorAction` compares `ApprovedBy()` strings** (`format:name`). A renamed approver (a new enrolment under the old
  key) would refuse at recheck with `ErrOperatorApprover` — fail-closed. S3 records exactly `VerifiedOperatorAction.ApprovedBy()`.
- **R4 — the writes-nothing test swaps process-global `os.Stdout`/`os.Stderr`**, as S2 M4's does. It is not `t.Parallel`, so no other
  test of the package runs concurrently; under `-race` the package passed (§1.5).

---

## 9. Acceptance criteria (the judge scores against these)

1. **No new daemon route, no schema change, no `.ail`, `surfaceCSP` and `app.js` unchanged, no served-behaviour change**:
   `DAEMON_ROUTES: 47 (base 47)`, `SURFACE_POSTS: 3`, `UNTOUCHED_VIOLATIONS: none`, `UNCHANGED:` app.js / surface_page.go /
   daemon.go / host/store/credentials.go, `verify_ail` identical to the base's; `TestFrozenV1RouteTableMatchesMux`,
   `TestConfirmRoutesRegistered`, `TestDaemonCredentialWriterSites`, `TestCLICredentialWriterSites`, `TestProductionContextRoots`,
   `TestDaemonDependencyAllowlist`, `TestApprovalImportsOnlyTheStandardLibrary` green unedited.
2. **The judges have no production caller**: `JUDGE_PRODUCTION_CALLERS: none` over the seven NEW identifiers, with the positive
   control ≥ 1; and S2's envelope/law has no production caller outside `judge_operator.go` (`ENVELOPE_…_OUTSIDE_THE_JUDGES: none`,
   control ≥ 1).
3. **`VerifiedOperatorAction` cannot be constructed outside `host/authority`**: every field unexported (`TestVerifiedOperatorActionOnlyFromTheJudges`,
   reflect — a keyed or unkeyed literal naming them does not compile outside the package); inside the package exactly one
   composite literal, in `verifiedOperatorAction`, called by exactly `JudgeOperatorAction` and `RecheckOperatorAction` (AST, with
   the positive control in the same count); nil, the zero value and an unsigned `{rechecked: true}` bind nothing
   (`TestOperatorUnjudgedProofBindsNothing`; `BINDS-NIL`, `APPROVEDBY-NIL`, `PASSKEYUSE-NIL`, `RECHECKED-NIL`, `RECHECKED-SIGNED`
   killed behaviourally).
4. **The judges read the approver from the store, never a caller** (judge-408 note 1): no judge parameter holds an
   `approval.Approver` (static, positive control); each store-sourced check is its own mutant, killed by a named test and by the
   whole host/authority package: `REUSED-SSH-UNKNOWN`, `REUSED-PK-UNKNOWN`, `REUSED-SSH-STORE-ERR`, `REUSED-PK-STORE-ERR`,
   `REUSED-SSH-REVOKED`, `REUSED-PK-REVOKED`, `REUSED-SSH-ENROLLED`, `REUSED-PK-ENROLLED`, `REUSED-SSH-FORMAT`, `REUSED-PK-FORMAT`,
   `REUSED-PK-COUNTER`, `REUSED-PK-STOREDKEY`, `REUSED-SSH-NAME`; the two lookup-equality equivalents named with reasons.
5. **Every refusal asserts its `*JudgeError` step, its class and that no proof is produced** (`ojRefused`: 19 confirm rows, 9 recheck
   rows).
6. **`X-PK-ALIAS` killed** by `TestOperatorPasskeyUseIsACopy` (named and whole host/approval), after reproducing SURVIVED at the base
   (C9); `PK-USE-ALIAS` likewise.
7. **Judge-408 note 3 disposed by a named pin**: `TestOperatorRunnerRecheckBoundsExpiryNotLifetime`, citing design §3.1 lines 80 and
   83; `R1-LIFETIME-ADDED` killed (SURVIVED at the base).
8. **Judge-408 note 4 carried with its reason** (S2b-D5, Appendix C); not claimed.
9. **The judges write nothing** (`TestOperatorJudgesWriteNothing`, control in-test; `WN-LOG-INJECT` killed).
10. **Size**: `git diff -M --numstat 4b1c39e2 HEAD -- cmd host docs app` sums to **526** (≤ 1,400; the same unfiltered); each
    sub-milestone's numstat as §7 (401 / 469 / 526).
11. **Exactly these files change: the four added in §5. Pre-existing files changed or deleted: none — production or test**
    (`PREEXISTING_CHANGED: none`). In particular `judge_verb.go`, `judge.go`, `judge_test.go` and every S2 test are unedited.
12. **Mutants**: every `K` row killed by its named test and by the whole package, none by the static test alone; the `E` rows named
    with reasons; the plan's ids are the ledger's (`final2`).
13. **Gates**: §6's script `FAILS: none` (or only a directive-listed load-sensitive test, green alone), `verify_go` rc 0.

---

## Appendix C — after this sprint (S2b's close; what S3–S12 own)

- **Closed here**: S2's "Moved out (S2-D1)" — the `host/authority` adapters, with their first mutant pass (C4 measured the parked
  draft's gaps; run 2 has none); design §3.1 [r1] "NEW typed verified-action proof can originate only in these judges" (the proof
  is `authority.VerifiedOperatorAction`); the store half of §5.2 `TestOperatorConfirmBothKeys`/`…Authenticity`/`…RevocationCounter`
  through the judges; the pure + store halves of `TestOperatorRunnerRechecks` except the claimed-row and daemon-status parts
  (below); judge-408 notes 1, 2 and 3.
- **What S3 inherits from S2b**: the judges' inputs are S3's row: `OperatorJudgeInput{Action: DecodeOperatorAction(row bytes), Digest:
  the digest column, BytesDigest: SHA-256 of the stored bytes, Challenge, CreatedAt, ExpiresAt, Consumed}` at the CAS's fresh `Now`;
  the confirm CAS binds `VerifiedOperatorAction.PasskeyUse()` (a copy, now pinned) and records `ApprovedBy()` as the action's
  approver identity, which the runner passes back as `approvedBy`; the runner calls `RecheckOperatorAction` with its **own** frozen
  action, recorded challenge and expiry, never the row's digest. A second CAS on the same action must not call the judge's counter
  rule again (design §3.1 "never advanced a second time"): R2 is the runner's, W15 the confirmation's. S3 also needs D-WORLD-119 A/B
  (design §8) and stays start-only unless D-WORLD-120 is answered (the audit residual).
- **§5.1/§5.2 rows S3–S12 own** (by design §8; S2's list minus what S2b closed): **S3** schema 10 `operator_actions`,
  `operator_runner_state`, `operator_action_steps`, the receipts and CAS — `TestOperatorAdmissionBoundaries`,
  `TestOperatorDuplicateTap` (keyed on `OperatorReviewDigest`), `TestOperatorReconnect` (durable dedup), `TestOperatorConfirmReplay`,
  `TestOperatorConfirmCrash`, `TestOperatorClaimBeforeExpiry`, `TestOperatorClaimRace`, `TestOperatorTerminalExclusion`,
  `TestOperatorPrune`, and the claimed-row half of `TestOperatorRunnerRechecks` ("trusts … daemon status": the runner reading a
  claimed row and calling `RecheckOperatorAction`); **S4** `TestOperatorRunnerUnavailable`, `TestOperatorCatalogueAdmission`,
  `TestOperatorRecoverUntouched`/`…Reserved` (runner side), `TestOperatorStatusDisconnected`, `TestOperatorCancelInflight`, the
  service rows (§5.5); **S5–S6** §5.3; **S7** the routes (calling `JudgeOperatorAction`, never the S2 law directly) —
  `TestOperatorMisdirected`, `…CrossOrigin`, `…FetchSite`, `…MediaType`, `…BodyBound`, `…UnsupportedMethod`, `…RequestLimits` (body and
  task at the route), `…LiteralTask` (rendering), `…BarStartReview`, `…BarFinishReview`, `…OtherVerbParity` at the route; **S8**
  `…PanelNetwork`, `…PanelLatest`, `…LinkAllowlist`, `…WireTypes` (page), S1's F8 (`null` elements in `app.js` arrays); **S9–S11**
  §5.4; **S12** the herdr hand-off.
- **Carried (not acted on)**: judge-408 note 4 — field drops of episode/kind/argv/options are binding-redundant through Validate's O9
  (S2b-D5); killed behaviourally by S2's round-trip tests, not by `TestOperatorConfirmBinding` by name. Judge-407 note (3) — numeric
  CLI flags mirrored as strings refuse `--ttl x` as terminal not usage; the envelope carries no numeric option (S7's route renders
  the class). Judge-407 note (2) / S1 F8 — `null` array elements throw in `app.js`: S8's. Iteration 407 retro — an attended merge may
  land mid-sprint: S2b touches no shared file.
- **Audit residuals (iteration 406, unchanged)**: an uncertain action has no attended resolver; a headless child outliving the
  runner is unowned; S3's schema 10 start-only unless D-WORLD-120 is answered; the second credentials writer vs `LOCK_NB` and the
  writer-site inventories move together; the herdr plugin hard-codes `codex=gpt-6.1-sol` (`world-config-init:96`).
- **Controller decisions requested**: (1) edit `ctl_roundtrip.sh`'s loop to `M1 M2 M3` (S2b has no M4). (2) None on the ledger:
  D-WORLD-119/120 remain the ledger's; S2b is answer-independent.
- The controller records; this plan moves no design doc and edits no charter.

---

## Appendix P — the patches (apply with §2, in order; generated from the compiled scratch stages; do not edit)

### Patch M1 — the host/authority operator judges and the judge-only proof

````diff M1
diff --git a/host/authority/judge_operator.go b/host/authority/judge_operator.go
new file mode 100644
index 00000000..1d474b16
--- /dev/null
+++ b/host/authority/judge_operator.go
@@ -0,0 +1,116 @@
+package authority
+
+// Row 243 V4 sprint 2b (design §3.1 [r1], §5.2, §6; parked from sprint 2 by its size cut, S2-D1): the
+// operator-action judges beside the verb judges, on a read-only handle with judge.go's step names.
+// Neither judge takes an approver: each reads the enrolled approver the artifact names from ro
+// (sshApprover, passkeyApprover), so the law judges the store's key, format, enrolment, revocation and
+// counter, never a caller's. VerifiedOperatorAction is the typed verified-action proof: its fields
+// are unexported and only JudgeOperatorAction and RecheckOperatorAction construct one.
+
+import (
+	"context"
+	"errors"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/store"
+)
+
+// VerifiedOperatorAction is a judged operator-action proof; the zero value binds nothing.
+type VerifiedOperatorAction struct {
+	signed    *approval.SignedOperatorAction
+	rechecked bool
+}
+
+func verifiedOperatorAction(s *approval.SignedOperatorAction, rechecked bool) *VerifiedOperatorAction {
+	return &VerifiedOperatorAction{signed: s, rechecked: rechecked}
+}
+
+// ErrOperatorApprover is a recheck whose artifact's approver is not the one the action recorded.
+var ErrOperatorApprover = errors.New("authority: the action's recorded approver did not sign it")
+
+// OperatorJudgeInput is a stored operator action at now (S3's row): its strictly decoded envelope, the
+// digest column, the challenge, the SHA-256 of the stored bytes, its times and whether it is consumed.
+type OperatorJudgeInput struct {
+	Action                         approval.OperatorAction
+	Digest, Challenge, BytesDigest [32]byte
+	CreatedAt, ExpiresAt, Now      int64
+	Consumed                       bool
+}
+
+func (in OperatorJudgeInput) admit() approval.AdmitInput {
+	return approval.AdmitInput{Now: in.Now, CreatedAt: in.CreatedAt, ExpiresAt: in.ExpiresAt, Consumed: in.Consumed,
+		StoredDigest: in.Digest, BytesDigest: in.BytesDigest}
+}
+
+// operatorApprover is the artifact's format (an armored SSHSIG, else a passkey assertion) and the
+// enrolled approver it names on ro, or a *JudgeError.
+func operatorApprover(ctx context.Context, ro *store.Credentials, artifact []byte) (string, approval.Approver, error) {
+	if sig, sigErr := approval.ParseSSHSIG(artifact); sigErr == nil {
+		a, err := sshApprover(ctx, ro, sig)
+		return approval.FormatSSH, a, err
+	} else if resp, respErr := approval.ParseAssertionResponse(artifact); respErr == nil {
+		a, err := passkeyApprover(ctx, ro, resp)
+		return approval.FormatWebAuthn, a, err
+	} else {
+		return "", approval.Approver{}, &JudgeError{Step: JudgeStepParse, Err: errors.Join(ErrUnknownArtifact, sigErr, respErr)}
+	}
+}
+
+// JudgeOperatorAction is the daemon's confirmation judge (S7): the approver the artifact names on ro,
+// then approval.ConfirmOperatorAction over in. It returns the proof or a *JudgeError.
+func JudgeOperatorAction(ctx context.Context, ro *store.Credentials, in OperatorJudgeInput, artifact []byte) (*VerifiedOperatorAction, error) {
+	format, a, err := operatorApprover(ctx, ro, artifact)
+	if err != nil {
+		return nil, err
+	}
+	s, err := approval.ConfirmOperatorAction(approval.OperatorConfirmInput{Action: in.Action, Admit: in.admit(), Challenge: in.Challenge,
+		Artifact: artifact, Format: format, Approver: a})
+	if err != nil {
+		return nil, &JudgeError{Step: JudgeStepLaw, Err: err}
+	}
+	return verifiedOperatorAction(s, false), nil
+}
+
+// RecheckOperatorAction is the runner's recheck before an effect (design §3.1): its own frozen action,
+// the recorded challenge and expiry, the current key on ro at now, then the recorded approver identity.
+// The lifetime bound is the daemon's (C2 at confirmation); the runner rechecks expiry and revocation.
+func RecheckOperatorAction(ctx context.Context, ro *store.Credentials, action approval.OperatorAction, challenge [32]byte, expiry, now int64,
+	artifact []byte, approvedBy string) (*VerifiedOperatorAction, error) {
+	format, a, err := operatorApprover(ctx, ro, artifact)
+	if err != nil {
+		return nil, err
+	}
+	s, err := approval.CheckOperatorAction(approval.OperatorCheckInput{Action: action, Challenge: challenge, Expiry: expiry, Now: now,
+		Artifact: artifact, Format: format, Approver: a})
+	if err != nil {
+		return nil, &JudgeError{Step: JudgeStepLaw, Err: err}
+	}
+	if s.ApprovedBy() != approvedBy {
+		return nil, &JudgeError{Step: JudgeStepLaw, Err: ErrOperatorApprover}
+	}
+	return verifiedOperatorAction(s, true), nil
+}
+
+// Binds, ApprovedBy and PasskeyUse are the law's, nil-safe; Rechecked says a runner recheck made it.
+func (v *VerifiedOperatorAction) Binds(a approval.OperatorAction, challenge [32]byte, expiry int64) error {
+	if v == nil {
+		return approval.ErrNotVerified
+	}
+	return v.signed.Binds(a, challenge, expiry)
+}
+
+func (v *VerifiedOperatorAction) ApprovedBy() string {
+	if v == nil {
+		return ""
+	}
+	return v.signed.ApprovedBy()
+}
+
+func (v *VerifiedOperatorAction) PasskeyUse() ([]byte, uint32, bool) {
+	if v == nil {
+		return nil, 0, false
+	}
+	return v.signed.PasskeyUse()
+}
+
+func (v *VerifiedOperatorAction) Rechecked() bool { return v != nil && v.signed != nil && v.rechecked }
diff --git a/host/authority/judge_operator_test.go b/host/authority/judge_operator_test.go
new file mode 100644
index 00000000..c5fd0d8b
--- /dev/null
+++ b/host/authority/judge_operator_test.go
@@ -0,0 +1,285 @@
+package authority
+
+// Row 243 V4 sprint 2b (design §3.1 [r1], §5.2 ConfirmBothKeys and RunnerRechecks; judge-408 note 1): the
+// operator judges on judge_test.go's rig (mark-mac's passkey, counter 1, and mark-ssh, both enrolled at
+// 900) with approvaltest keys (D-WORLD-109). Every approver fact the law judges is the store's: each row
+// edits the store, not an argument, and wants its refusal step, its error class and no proof.
+
+import (
+	"errors"
+	"fmt"
+	"go/ast"
+	"go/parser"
+	"go/token"
+	"io/fs"
+	"reflect"
+	"strings"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
+	"github.com/sunholo-data/ailang-world/host/store"
+)
+
+var (
+	ojStranger = approvaltest.NewSigner(91)
+	ojOtherKey = approvaltest.NewPasskey(approval.AlgEdDSA, 92)
+)
+
+func ojAction() approval.OperatorAction {
+	a := approval.OperatorAction{Schema: approval.OperatorActionSchema, ID: "oa-0123456789abcdef0123456789abcdef", Kind: "start", Episode: "fixture-ep",
+		Store: "/ws/state/credentials.db", Options: approval.OperatorOptions{Preset: "sol", Headless: true}, Task: "t", Review: "r",
+		Resolution: []approval.OperatorResolved{{Name: "agent", Value: "codex"}}}
+	a.Argv = approval.OperatorArgv(a.Kind, a.Episode, a.Options)
+	return a
+}
+
+// ojIn is a requested action created at 1000, expiring at 1600, judged at 1010, with its SSHSIG by
+// mark-ssh and its assertion (counter 2) by mark-mac.
+func ojIn(t *testing.T) (OperatorJudgeInput, []byte, []byte) {
+	t.Helper()
+	a := ojAction()
+	d, err := approval.OperatorSignatureDigest(a)
+	if err != nil {
+		t.Fatal(err)
+	}
+	in := OperatorJudgeInput{Action: a, Digest: d, BytesDigest: d, Challenge: [32]byte{0x6f}, CreatedAt: 1000, ExpiresAt: 1600, Now: 1010}
+	return in, ojSSH(t, jgSigner, approval.OperatorActionNamespace, in), ojAssert(t, jgPasskey, in, 2)
+}
+
+func ojSSH(t *testing.T, s *approvaltest.Signer, ns string, in OperatorJudgeInput) []byte {
+	t.Helper()
+	m, err := approval.OperatorActionMessage(in.Digest, in.Challenge, in.ExpiresAt)
+	if err != nil {
+		t.Fatal(err)
+	}
+	return s.Sign(ns, m)
+}
+
+func ojAssert(t *testing.T, p *approvaltest.Passkey, in OperatorJudgeInput, count uint32) []byte {
+	t.Helper()
+	h, err := approval.OperatorActionMessageHash(in.Digest, in.Challenge, in.ExpiresAt)
+	if err != nil {
+		t.Fatal(err)
+	}
+	return p.Assert(jgOrigin, jgRPID, h, count)
+}
+
+// ojOpen is the rig after setup, read-only.
+func ojOpen(t *testing.T, setup ...string) *store.Credentials {
+	t.Helper()
+	path, _ := jgRig(t)
+	jgExec(t, path, setup...)
+	ro, err := store.OpenCredentialsReadOnly(path)
+	if err != nil {
+		t.Fatal(err)
+	}
+	t.Cleanup(func() { _ = ro.Close() })
+	return ro
+}
+
+// ojRefused wants no proof, a *JudgeError at step, and want: a binding field, an error class or an
+// admission verdict.
+func ojRefused(t *testing.T, name string, v *VerifiedOperatorAction, err error, step string, want any) {
+	t.Helper()
+	var je *JudgeError
+	var na *approval.NotAdmittedError
+	ok := v == nil && errors.As(err, &je) && je.Step == step
+	var be *approval.BindingError
+	switch w := want.(type) {
+	case *approval.BindingError:
+		ok = ok && errors.As(err, &be) && be.Field == w.Field
+	case error:
+		ok = ok && errors.Is(err, w)
+	case approval.Verdict:
+		ok = ok && errors.As(err, &na) && na.Verdict == w
+	}
+	if !ok {
+		t.Errorf("%s: %v, %v; want step %s, %v and no proof", name, v, err, step, want)
+	}
+}
+
+func TestOperatorJudgeBothKeys(t *testing.T) {
+	in, ssh, pk := ojIn(t)
+	for by, art := range map[string][]byte{"ssh:mark-ssh": ssh, "webauthn:mark-mac": pk} {
+		v, err := JudgeOperatorAction(boundedTestContext(t), ojOpen(t), in, art)
+		if err != nil || v.ApprovedBy() != by || v.Binds(in.Action, in.Challenge, in.ExpiresAt) != nil || v.Rechecked() {
+			t.Fatalf("%s: %v", by, err)
+		}
+		other := in.Action
+		other.Task += "!"
+		var be *approval.BindingError
+		if err := v.Binds(other, in.Challenge, in.ExpiresAt); !errors.As(err, &be) || be.Field != "digest" {
+			t.Errorf("%s: binds another action: %v", by, err)
+		}
+		cid, n, ok := v.PasskeyUse()
+		if (by == "webauthn:mark-mac") != ok || (ok && (n != 2 || string(cid) != string(jgPasskey.CredentialID()))) {
+			t.Errorf("%s: passkey use %x %d %v", by, cid, n, ok)
+		}
+	}
+	// another enrolled key is judged under its own stored name
+	ops := fmt.Sprintf(`INSERT INTO approvers (name, format, public_key, enrolled_at) VALUES ('ops-ssh', 'ssh', x'%x', 900)`, ojStranger.PublicKey().Blob)
+	if v, err := JudgeOperatorAction(boundedTestContext(t), ojOpen(t, ops), in, ojSSH(t, ojStranger, approval.OperatorActionNamespace, in)); err != nil || v.ApprovedBy() != "ssh:ops-ssh" {
+		t.Fatalf("ops-ssh: %v", err)
+	}
+}
+
+// TestOperatorJudgeReadsTheStore: every approver fact is the store's (judge-408 note 1), and every
+// admission input is the row's; each row is its own refusal.
+func TestOperatorJudgeReadsTheStore(t *testing.T) {
+	in, ssh, pk := ojIn(t)
+	closed := ojOpen(t)
+	_ = closed.Close()
+	for _, c := range []struct {
+		name  string
+		setup []string
+		ro    *store.Credentials
+		edit  func(*OperatorJudgeInput)
+		art   []byte
+		step  string
+		want  any
+	}{
+		{"typed y", nil, nil, nil, []byte("y\n"), JudgeStepParse, ErrUnknownArtifact},
+		{"unknown ssh key", nil, nil, nil, ojSSH(t, ojStranger, approval.OperatorActionNamespace, in), JudgeStepUnknownKey, ErrNotEnrolled},
+		{"unknown passkey credential", nil, nil, nil, ojAssert(t, ojOtherKey, in, 2), JudgeStepUnknownKey, ErrNotEnrolled},
+		{"ssh store unreadable", nil, closed, nil, ssh, JudgeStepStore, nil},
+		{"passkey store unreadable", nil, closed, nil, pk, JudgeStepStore, nil},
+		{"ssh revoked", []string{`UPDATE approvers SET revoked_at = 1000 WHERE name = 'mark-ssh'`}, nil, nil, ssh, JudgeStepLaw, approval.ErrApproverRevoked},
+		{"passkey revoked", []string{`UPDATE approvers SET revoked_at = 1000 WHERE name = 'mark-mac'`}, nil, nil, pk, JudgeStepLaw, approval.ErrApproverRevoked},
+		{"ssh enrolled after now", []string{`UPDATE approvers SET enrolled_at = 1011 WHERE name = 'mark-ssh'`}, nil, nil, ssh, JudgeStepLaw, approval.ErrApproverClockBehind},
+		{"passkey enrolled after now", []string{`UPDATE approvers SET enrolled_at = 1011 WHERE name = 'mark-mac'`}, nil, nil, pk, JudgeStepLaw, approval.ErrApproverClockBehind},
+		{"ssh key stored as webauthn", []string{`UPDATE approvers SET format = 'webauthn' WHERE name = 'mark-ssh'`}, nil, nil, ssh, JudgeStepLaw, approval.ErrApproverFormat},
+		{"passkey stored as ssh", []string{`UPDATE approvers SET format = 'ssh' WHERE name = 'mark-mac'`}, nil, nil, pk, JudgeStepLaw, approval.ErrApproverNotPasskey},
+		{"stored counter at the assertion's", []string{`UPDATE webauthn_credentials SET sign_count = 2 WHERE name = 'mark-mac'`}, nil, nil, pk, JudgeStepLaw, approval.ErrSignCountRegression},
+		{"stored key unparsable", []string{`UPDATE approvers SET public_key = x'00' WHERE name = 'mark-mac'`}, nil, nil, pk, JudgeStepStoredKey, nil},
+		{"a verb namespace", nil, nil, nil, ojSSH(t, jgSigner, approval.VerbNamespace, in), JudgeStepLaw, approval.ErrNamespace},
+		{"consumed", nil, nil, func(i *OperatorJudgeInput) { i.Consumed = true }, ssh, JudgeStepLaw, approval.RefuseConsumed},
+		{"now at expiry", nil, nil, func(i *OperatorJudgeInput) { i.Now = 1600 }, ssh, JudgeStepLaw, approval.RefuseExpired},
+		{"lifetime 601", nil, nil, func(i *OperatorJudgeInput) { i.CreatedAt = 999 }, ssh, JudgeStepLaw, approval.ErrOperatorLifetime},
+		{"stored bytes differ", nil, nil, func(i *OperatorJudgeInput) { i.BytesDigest[0] ^= 1 }, ssh, JudgeStepLaw, approval.RefuseDigestMismatch},
+		{"another action's digest column", nil, nil, func(i *OperatorJudgeInput) { i.Digest[0] ^= 1; i.BytesDigest = i.Digest }, ssh, JudgeStepLaw, &approval.BindingError{Field: "digest"}},
+	} {
+		ro, row := c.ro, in
+		if ro == nil {
+			ro = ojOpen(t, c.setup...)
+		}
+		if c.edit != nil {
+			c.edit(&row)
+		}
+		v, err := JudgeOperatorAction(boundedTestContext(t), ro, row, c.art)
+		ojRefused(t, c.name, v, err, c.step, c.want)
+	}
+}
+
+// TestOperatorRunnerRechecks: the runner trusts no row digest, stored approver or daemon status.
+func TestOperatorRunnerRechecks(t *testing.T) {
+	in, ssh, pk := ojIn(t)
+	moved := `UPDATE webauthn_credentials SET sign_count = 2 WHERE name = 'mark-mac'`
+	for by, c := range map[string]struct {
+		setup []string
+		art   []byte
+	}{"ssh:mark-ssh": {nil, ssh}, "webauthn:mark-mac": {[]string{moved}, pk}} {
+		v, err := RecheckOperatorAction(boundedTestContext(t), ojOpen(t, c.setup...), in.Action, in.Challenge, in.ExpiresAt, 1599, c.art, by)
+		if err != nil || !v.Rechecked() || v.ApprovedBy() != by || v.Binds(in.Action, in.Challenge, in.ExpiresAt) != nil {
+			t.Fatalf("%s control: %v", by, err)
+		}
+	}
+	forged := in.Action
+	forged.Task = "another task"
+	for name, c := range map[string]struct {
+		setup     []string
+		action    approval.OperatorAction
+		challenge [32]byte
+		art       []byte
+		by        string
+		now       int64
+		step      string
+		want      any
+	}{
+		"a forged row's action":     {nil, forged, in.Challenge, ssh, "ssh:mark-ssh", 1500, JudgeStepLaw, approval.ErrBadSignature},
+		"another challenge":         {nil, in.Action, [32]byte{0x70}, ssh, "ssh:mark-ssh", 1500, JudgeStepLaw, approval.ErrBadSignature},
+		"another recorded approver": {nil, in.Action, in.Challenge, ssh, "webauthn:mark-mac", 1500, JudgeStepLaw, ErrOperatorApprover},
+		"no recorded approver":      {nil, in.Action, in.Challenge, ssh, "", 1500, JudgeStepLaw, ErrOperatorApprover},
+		"revoked since":             {[]string{`UPDATE approvers SET revoked_at = 1400 WHERE name = 'mark-ssh'`}, in.Action, in.Challenge, ssh, "ssh:mark-ssh", 1500, JudgeStepLaw, approval.ErrApproverRevoked},
+		"expired at now":            {nil, in.Action, in.Challenge, ssh, "ssh:mark-ssh", 1600, JudgeStepLaw, approval.RefuseExpired},
+		"counter never recorded":    {nil, in.Action, in.Challenge, pk, "webauthn:mark-mac", 1500, JudgeStepLaw, approval.ErrOperatorCounter},
+		"unknown key":               {nil, in.Action, in.Challenge, ojSSH(t, ojStranger, approval.OperatorActionNamespace, in), "ssh:mark-ssh", 1500, JudgeStepUnknownKey, ErrNotEnrolled},
+		"typed y":                   {nil, in.Action, in.Challenge, []byte("y\n"), "ssh:mark-ssh", 1500, JudgeStepParse, ErrUnknownArtifact},
+	} {
+		v, err := RecheckOperatorAction(boundedTestContext(t), ojOpen(t, c.setup...), c.action, c.challenge, in.ExpiresAt, c.now, c.art, c.by)
+		ojRefused(t, name, v, err, c.step, c.want)
+	}
+}
+
+// TestOperatorUnjudgedProofBindsNothing: what a refusal returns (nil) and a proof no judge built bind nothing.
+func TestOperatorUnjudgedProofBindsNothing(t *testing.T) {
+	for name, v := range map[string]*VerifiedOperatorAction{"nil": nil, "zero": {}, "unsigned, flagged": {rechecked: true}} {
+		_, _, ok := v.PasskeyUse()
+		if ok || !errors.Is(v.Binds(ojAction(), [32]byte{}, 0), approval.ErrNotVerified) || v.ApprovedBy() != "" || v.Rechecked() {
+			t.Errorf("%s: an unjudged proof binds something", name)
+		}
+	}
+}
+
+// The proof can be constructed only by verifiedOperatorAction, called only by the two judges, and
+// neither judge accepts an approver from its caller (the law's own input does: the positive control).
+func TestVerifiedOperatorActionOnlyFromTheJudges(t *testing.T) {
+	for _, typ := range []reflect.Type{reflect.TypeOf(VerifiedOperatorAction{}), reflect.TypeOf(approval.SignedOperatorAction{})} {
+		for i := 0; i < typ.NumField(); i++ {
+			if typ.Field(i).IsExported() {
+				t.Errorf("%s.%s is exported", typ.Name(), typ.Field(i).Name)
+			}
+		}
+	}
+	approverT := reflect.TypeOf(approval.Approver{})
+	holds := func(typ reflect.Type) bool {
+		for i := 0; typ.Kind() == reflect.Struct && i < typ.NumField(); i++ {
+			if typ.Field(i).Type == approverT {
+				return true
+			}
+		}
+		return typ == approverT
+	}
+	if !holds(reflect.TypeOf(approval.OperatorConfirmInput{})) {
+		t.Fatal("control: the law's input holds no approver")
+	}
+	for _, f := range []any{JudgeOperatorAction, RecheckOperatorAction} {
+		for ft, i := reflect.TypeOf(f), 0; i < ft.NumIn(); i++ {
+			if holds(ft.In(i)) {
+				t.Errorf("a judge takes an approver from its caller: %v", ft)
+			}
+		}
+	}
+	pkgs, err := parser.ParseDir(token.NewFileSet(), ".", func(fi fs.FileInfo) bool { return !strings.HasSuffix(fi.Name(), "_test.go") }, 0)
+	if err != nil {
+		t.Fatal(err)
+	}
+	lits, calls := 0, map[string]int{}
+	for _, f := range pkgs["authority"].Files {
+		for _, d := range f.Decls {
+			fd, ok := d.(*ast.FuncDecl)
+			if !ok {
+				continue
+			}
+			ast.Inspect(fd, func(n ast.Node) bool {
+				if cl, ok := n.(*ast.CompositeLit); ok {
+					if id, ok := cl.Type.(*ast.Ident); ok && id.Name == "VerifiedOperatorAction" {
+						lits++
+						if fd.Name.Name != "verifiedOperatorAction" {
+							t.Errorf("%s builds a VerifiedOperatorAction", fd.Name.Name)
+						}
+					}
+				}
+				if ce, ok := n.(*ast.CallExpr); ok {
+					if id, ok := ce.Fun.(*ast.Ident); ok && id.Name == "verifiedOperatorAction" {
+						calls[fd.Name.Name]++
+					}
+				}
+				return true
+			})
+		}
+	}
+	if lits != 1 || len(calls) != 2 || calls["JudgeOperatorAction"] != 1 || calls["RecheckOperatorAction"] != 1 {
+		t.Fatalf("constructor literals %d, callers %v; want 1 and the two judges (the positive control)", lits, calls)
+	}
+}
````

### Patch M2 — judge-408 pins: the passkey use is a copy; the runner bounds expiry, not lifetime

````diff M2
diff --git a/host/approval/operator_pins_test.go b/host/approval/operator_pins_test.go
new file mode 100644
index 00000000..1d1d03fd
--- /dev/null
+++ b/host/approval/operator_pins_test.go
@@ -0,0 +1,68 @@
+package approval_test
+
+// Row 243 V4 sprint 2b (judge-408 notes 2 and 3): the proof's passkey use is its own copy (X-PK-ALIAS), and
+// the runner's recheck bounds the recorded expiry, not the request's lifetime (design §3.1).
+
+import (
+	"bytes"
+	"errors"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+)
+
+// TestOperatorPasskeyUseIsACopy: neither the caller's approver nor a returned slice, edited after the
+// call, changes the credential ID a proof binds (S3's confirm CAS keys on it).
+func TestOperatorPasskeyUseIsACopy(t *testing.T) {
+	want := opPasskey.CredentialID()
+	in := opIn(t, approval.FormatWebAuthn)
+	confirmed, err := approval.ConfirmOperatorAction(in)
+	if err != nil {
+		t.Fatal(err)
+	}
+	in.Approver.Passkey.SignCount = 2
+	rechecked, err := approval.CheckOperatorAction(approval.OperatorCheckInput{Action: in.Action, Challenge: opC, Expiry: opExp, Now: opNow,
+		Artifact: in.Artifact, Format: approval.FormatWebAuthn, Approver: in.Approver})
+	if err != nil {
+		t.Fatal(err)
+	}
+	in.Approver.Passkey.CredentialID[0] ^= 0xff
+	for name, s := range map[string]*approval.SignedOperatorAction{"confirmed": confirmed, "rechecked": rechecked} {
+		cid, _, ok := s.PasskeyUse()
+		if !ok || !bytes.Equal(cid, want) {
+			t.Fatalf("%s: the approver's slice, edited after the call, moved the proof's credential ID to %x", name, cid)
+		}
+		cid[1] ^= 0xff
+		if again, _, _ := s.PasskeyUse(); !bytes.Equal(again, want) {
+			t.Fatalf("%s: a returned slice, edited, moved the proof's credential ID to %x", name, again)
+		}
+	}
+}
+
+// TestOperatorRunnerRecheckBoundsExpiryNotLifetime pins judge-408 note 3. Design §3.1: the 600 s bound is
+// the signed request's ("signed request expiry 600 s"), judged by the daemon at confirmation (C2), where
+// the row's creation time is known; the runner "recheck[s] expiry and revocation again at each new effect
+// boundary" from its own recorded values, and OperatorCheckInput carries no creation time. So a recorded
+// expiry an hour away is the runner's to admit while unexpired, and the confirmation refuses it.
+func TestOperatorRunnerRecheckBoundsExpiryNotLifetime(t *testing.T) {
+	const exp = opCreated + 3600
+	a := opAction("start")
+	for format, art := range map[string][]byte{approval.FormatSSH: opSSH(t, opSigner, approval.OperatorActionNamespace, a, exp),
+		approval.FormatWebAuthn: opAssert(t, a, exp, 2)} {
+		in := opIn(t, format)
+		in.Artifact, in.Admit.ExpiresAt = art, exp
+		if s, err := approval.ConfirmOperatorAction(in); s != nil || !errors.Is(err, approval.ErrOperatorLifetime) {
+			t.Errorf("%s: the confirmation admitted a 3600 s request: %v", format, err)
+		}
+		if in.Approver.Passkey != nil {
+			in.Approver.Passkey.SignCount = 2
+		}
+		for _, now := range []int64{opNow, exp - 1} {
+			s, err := approval.CheckOperatorAction(approval.OperatorCheckInput{Action: a, Challenge: opC, Expiry: exp, Now: now, Artifact: art,
+				Format: format, Approver: in.Approver})
+			if err != nil || s.Binds(a, opC, exp) != nil {
+				t.Errorf("%s at %d: the recheck refused an unexpired recorded expiry: %v", format, now, err)
+			}
+		}
+	}
+}
````

### Patch M3 — the operator judges write nothing

````diff M3
diff --git a/host/authority/judge_operator_silent_test.go b/host/authority/judge_operator_silent_test.go
new file mode 100644
index 00000000..1908edd7
--- /dev/null
+++ b/host/authority/judge_operator_silent_test.go
@@ -0,0 +1,57 @@
+package authority
+
+// Row 243 V4 sprint 2b (judge-407's retro rule, as sprint 2 M4): the two operator judges S7 and the
+// runner will call write nothing to stdout, stderr or the log, accepting or refusing.
+
+import (
+	"io"
+	"log"
+	"os"
+	"strings"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/store"
+)
+
+// ojSilent runs f with os.Stdout, os.Stderr and the log on one pipe and returns what was written.
+func ojSilent(t *testing.T, f func()) string {
+	t.Helper()
+	r, w, err := os.Pipe()
+	if err != nil {
+		t.Fatal(err)
+	}
+	out, errOut := os.Stdout, os.Stderr
+	os.Stdout, os.Stderr = w, w
+	log.SetOutput(w)
+	f()
+	os.Stdout, os.Stderr = out, errOut
+	log.SetOutput(os.Stderr)
+	_ = w.Close()
+	b, _ := io.ReadAll(r)
+	return string(b)
+}
+
+func TestOperatorJudgesWriteNothing(t *testing.T) {
+	if got := ojSilent(t, func() { _, _ = os.Stderr.WriteString("x"); _, _ = os.Stdout.WriteString("y"); log.Print("z") }); !strings.HasPrefix(got, "xy") || !strings.HasSuffix(got, " z\n") {
+		t.Fatalf("control: the harness saw %q, want x, y and the log line", got)
+	}
+	in, ssh, pk := ojIn(t)
+	ro, revoked, closed := ojOpen(t), ojOpen(t, `UPDATE approvers SET revoked_at = 1000`), ojOpen(t)
+	_ = closed.Close()
+	arts := [][]byte{ssh, pk, []byte("y\n"), ojSSH(t, ojStranger, approval.OperatorActionNamespace, in)}
+	ctx := boundedTestContext(t)
+	got := ojSilent(t, func() {
+		for _, h := range []*store.Credentials{ro, revoked, closed} {
+			for _, art := range arts {
+				_, _ = JudgeOperatorAction(ctx, h, in, art)
+				for _, by := range []string{"ssh:mark-ssh", "webauthn:mark-mac"} {
+					_, _ = RecheckOperatorAction(ctx, h, in.Action, in.Challenge, in.ExpiresAt, 1500, art, by)
+				}
+			}
+		}
+	})
+	if got != "" {
+		t.Fatalf("an operator judge wrote %q", got)
+	}
+}
````
