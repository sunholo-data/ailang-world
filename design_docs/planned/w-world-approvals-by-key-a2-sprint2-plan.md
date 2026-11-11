# Sprint plan — w-world-approvals-by-key-a2, sprint 2 (iteration 380, row 201)

> **THE CUT — measured, and the design's §9 sprint-2 fallback did NOT have to fire.** Sprint 2 is
> "the assertion law, the store and the binding": **M3b** (carried from sprint 1's Appendix C: the
> assertion law `ApproveByAssertion`, `checkApprover`, `SignedApproval.PasskeyUse`, the test
> authenticator `approvaltest.Passkey`, PA1–PA4, RC1, TP1, XW1, XW2's assertion part), **M4**
> (credentials schema 6, `webauthn_credentials`, `EnrolPasskey` and the three lookups, PS1–PS5, RC6)
> and **M5** (the consume's C13–C18, CS7–CS9, the CS2/CS3 arms, RC3, then authority's A6 and AU3).
> - **Measured: 1,843 changed lines under `cmd/ host/`** — `git diff --numstat <base> HEAD -- cmd host`,
>   insertions + deletions, after a commit of each of the 8 scratch stages (456 production, 1,217
>   test, 170 test support and testdata). **≤ 2,000, 157 lines of margin.** No `.gitattributes` or
>   doc line is counted (this sprint changes neither).
> - **The fallback was applied early, as a measurement, not as a cut.** Sprint 1 estimated M3b at
>   ≈ 450 and design §9 sprint 2 at ≈ 1,600 (≈ 2,050 together). The planner built M3b first and
>   measured it (890), then M4 (648) and M5's consume (239) — 1,777 — and only then priced the
>   fallback's two items, A6 (3 production lines) and AU3 (63 test lines): 1,843 with them. Moving them
>   would save 66 lines that the cap does not need, and would leave the store accepting a passkey
>   consume no caller builds. So **AU3 and A6 stay in sprint 2 (M5b)**, and **Appendix C moves nothing
>   to sprint 3** (it lists what sprint 3 inherits). Had the total exceeded 2,000, the cut point was
>   M5b first, then M5a (the consume) as a whole milestone — each leaves every committed state safe.
>
> **Every line of code and test is a unified diff in Appendix P**, generated LAST from scratch stages
> that were compiled, vetted, gofmt-clean and green at every sub-milestone, after every mutant ran on
> the final tree. The patches were then extracted **from this file** with §2's `awk`, applied to a
> fresh export of `3b77c294`, each sub-milestone's acceptance run, and the result compared
> byte-for-byte (`diff -r cmd host`) with the tree the mutants ran on (§1.2 C30–C31). The executor
> applies each patch and runs its acceptance; it writes no code by hand.

- **Design**: [w-world-approvals-by-key-a2.md](w-world-approvals-by-key-a2.md), **r2** (quorum r2 +
  controller carve-out; **not reopened**). Where this plan and the design differ, this plan was
  measured against `3b77c294` and wins; each difference is a decision (§1.3) or a premise found false
  (§1.4).
- **Starting ground**: sprint 1 (M0, M1, M2, M3a) landed as PR #546 (merge `d5d90e49`), recorded in
  `a8956f3a`; its Appendix C is this sprint's head. The iteration-379 judge's three carry-forwards
  are folded in (§1.3 S2-D9–S2-D11).
- **Planner base**: `3b77c294000ca40c2fe77e312ae0e8ce0f71a030` (= origin/dev). Worktree
  `/Users/voightkampff/dev/sunholo-data/.wt-world-iter380-row201-a2s2`, branch
  `sprint/world-iter380-row201-a2s2`. Planner scratch (outside the repo):
  `$HOME/.ailang/state/world-iter380/plan-scratch/` (`base/` the pristine export, `work/` the tree
  the mutants ran on, `stages/` one commit per sub-milestone, `tools/`, `logs/`).
- **Executor**: as routed by the controller (a small model). **Four runs**: Run 1 (M3b1, M3b2),
  Run 2 (M3b3, M4a), Run 3 (M4b, M4c), Run 4 (M5a, M5b). Each sub-milestone is one `patch -p1` (§2)
  plus its acceptance (§3). **Apply, never retype.**
- **Touched**: `host/approval/{signed.go, signed_test.go, webauthn_fixture_test.go,
  testdata/webauthn/SOURCE.txt}`, new `host/approval/{passkey_assert.go, passkey_assert_test.go,
  passkey_testauth_test.go, passkey_law_test.go, approvaltest/passkey.go}`;
  `host/store/{credentials_schema.sql, credentials.go, mint_requests.go, credentials_test.go,
  approvers_test.go, mint_requests_test.go, consume_binding_test.go}`, new `host/store/{passkeys.go,
  passkeys_test.go, passkey_lookups_test.go, passkey_consume_test.go}`; `host/authority/mint.go`, new
  `host/authority/mint_passkey_test.go`; `cmd/ailang-worldd/{approver_rig_test.go,
  approve_sig_test.go}` (one fixture statement each, §1.4 P2).
- **Not touched**: every other `host/approval` file (`approver.go`, `passkey.go`, `cbor.go`, `cose.go`,
  `authdata.go`, `clientdata.go`, `origin.go`, `passkey_json.go`, `request.go`, `admit.go`, `sshsig*.go`),
  `host/approval/approvaltest/signer*.go`, `host/store/approvers.go` (`useApprover` and
  `EnrolApprover` unchanged), every other `cmd/ailang-worldd` file, `host/daemon/**`,
  `host/surfacepage/**`, `.gitattributes`, `docs/QUICKSTART.md`, `world/*.ail`, `packages/**`,
  `scripts/**`, `tools/launchd/*`, `cmd/world-publish/*`, any skill, the typed-`y` path, `~/.ssh`,
  `~/.ailang/world`, `~/.ailang/state/mission-v1*`.

---

## 0. Binding constraints (each is restated where it governs)

1. **Fakes, temp stores and test keys only (D-WORLD-109 = A).** The only keys are the recorded
   fixtures (Chrome 157's **virtual** authenticator, iteration 379) and `approvaltest.Passkey`, a
   test authenticator whose private key is a repeated seed byte. No browser, Chrome, Node,
   `ssh-keygen` or real authenticator runs (the browser driver is sprint 3's M7). Every store is
   under `t.TempDir()`. No new test starts a subprocess. Nothing opens `~/.ailang/world`, `~/.ssh`, a
   real credentials store, a real daemon or `/dev/tty`; nothing mints into a real store or publishes.
2. **Pure core, effects at the boundary (coding standards S1–S3).** The assertion law is in
   `host/approval`, standard library only, no I/O (`TestApprovalImportsOnlyTheStandardLibrary`: the
   new production file imports only `bytes`, `crypto/sha256`, `errors`, `fmt`, all already
   imported there). `approvaltest` is test support, imported only by `_test.go` files
   (`TestApprovalTestIsTestOnly`). The only effects are SQL in `host/store`, each new method behind
   `requireDeadline`. No new Go module, package, exec site, route, context root or `.ail`.
3. **One mutant per new production branch, anchored to the diff** (§4): every `if`, `case`, operand,
   bound and copied field of Appendix P's production hunks, each run against its single named test
   **and** its package suite; named equivalents with their reasons; base controls where a test pins
   pre-existing behaviour.
4. **Typed errors with `errors.Is`/`errors.As`, and nothing written on a refusal.** Every store
   refusal arm checks the row counts of both `approvers` and `webauthn_credentials`, or the request
   still pending with no session and `last_used_at`/`sign_count` unchanged.
5. **No `t.Parallel()`** in any new test (design §6).
6. **Bisectable.** Each sub-milestone ends gofmt-clean (`gofmt -l cmd host` empty), `go build ./...`
   and `go vet` green, the `_test.go` compile fence (`go test -run '^$' ./...`) green, with its own
   tests and every earlier sub-milestone's tests green (§1.2 C31).
7. **No git writes by the executor** (§7). The controller commits from snapshots.

---

## 1. Pre-flight

### 1.1 Baseline at `3b77c294` (measured by the planner, 2026-11-11)

Every command ran with `export AILANG_BIN=$HOME/.pinned-ailang/ailang` (`AILANG v0.41.0`), `go1.26.6
darwin/arm64`, in a `git archive 3b77c294` export at `$HOME/.ailang/state/world-iter380/plan-scratch/base`
(a throwaway git repo; `diff -r cmd host` against the worktree: identical, `BASE_EQUALS_WORKTREE`; the
worktree itself was only read). Full logs: `plan-scratch/logs/baseline/`.

| Command | rc | Time | Result |
|---|---|---|---|
| `git -C <worktree> status --porcelain \| wc -l` | 0 | — | `0` (clean at `3b77c294`) |
| `gofmt -l cmd host` | **0** | 0 s | empty |
| `go build ./... && go vet ./...` | **0** | 6 s | clean |
| `go test -run '^$' ./...` (compile fence) | **0** | 23 s | no FAIL line |
| `go test -count=1 ./host/approval/... ./host/store/ ./host/authority/` | **0** | 12 s | ok ×4 |
| `go test -race -count=1 ./host/approval/... ./host/store/ ./host/authority/` | **0** | 31 s | ok ×4, no race report |
| `go test -count=1 ./cmd/ailang-worldd/` | **0** | 114 s | ok |
| the §6 inventories and reflective tests by name with `-v` (16 design-§7 tests; RC1, RC4, RC3, RC2; RC5 run alone) | **0** | 12 s | all 21 `--- PASS`, 0 SKIP/FAIL (`inv_all.log`, `rc5.log`); RC6 is new in M4c |
| `./scripts/verify_ail.sh` | **0** | 36 s | `✓ verify gate PASSED` |
| `VERIFY_GO_LEGS=plain ./scripts/verify_go.sh` | **1** | 374 s | **red at the base, `host/broker` only** (see below); every other leg green: `✓ 0 other binary blobs among 6677 tracked files`, `✓ all 64 required top-level evidence tests passed exactly once`, 61 other packages `ok` |

**The base's `verify_go.sh` red is `host/broker` under whole-suite load, and it measures the machine, not
this sprint.** In two base runs `host/broker` failed `TestPublisherErrorRedactsTheSecretAndKeepsTheMarker`
(run 1), then `TestNotPublishedDoesNotReAuthorizeTheConsumedApproval`,
`TestEveryReconcileRequestWasLoopbackAndAGet` and `TestWorktreeSnapshotRefusesAConfiguredFilter` (run 2):
the first three because a loopback fake server **received a stray `GET /`** it never asked for
(`registry_reconcile_test.go:1706: a reconciliation fetched "/"`), the fourth is the known load-sensitive
filter control. All four **PASS alone** on the same base (`broker_alone.log`, 1.2 s), and CI at
`3b77c294` is green (run 42933266511). `host/broker` is not touched by this sprint. Rule: a red at the
base measures the repo; the controller re-runs `host/broker` alone once and does not chase it.

**A harness trap, measured and routed (§7):** the planner's first base run was started under `nohup`,
and `scripts/marketingcheck`'s `TestMarketingSupervisorFakeDaemonCleanup/shell-hup` then failed 3 of 3
times **alone** (`hangup trap must exit 129 … signal: killed`). Run without `nohup` it passes (5.9 s,
`rerun_marketing_nonohup.log`). A shell started with SIGHUP ignored cannot trap it (POSIX: signals
ignored on entry to a non-interactive shell cannot be trapped), so **no gate may be run under `nohup`**.

### 1.2 Codebase claims checked (at `3b77c294` unless "work" — the planner's final scratch tree)

| # | Claim | Command → observed output |
|---|---|---|
| C1 | The base | `git rev-parse HEAD origin/dev` → `3b77c294000ca40c2fe77e312ae0e8ce0f71a030` twice; `git status --porcelain \| wc -l` → `0`; `git log --oneline -1` → `3b77c294 Merge pull request #547 … world-iter379-record` |
| C2 | The signature law's P2–P5 are inline | `grep -n '// P2\|// P3\|// P4\|// P5' host/approval/signed.go` → `59 … // P2`, `62 if in.Approver.Format != FormatSSH { // P3`, `65 … // P4`, `68 … // P5`; `SignedApproval` at `:16` has 7 fields (`ok digest challenge expiry name format key`), RC1's disposition map has 7 keys (`signed_test.go:165-173`) |
| C3 | What sprint 1 left for M3b | `grep -n '^func' host/approval/{passkey_json,clientdata,authdata,cose}.go` → `ParseAssertionResponse` (`passkey_json.go:77`), `ParseClientData`, `ParseAuthData(b, registration)`, `(COSEKey) verify`, `Canonical`; `passkey.go` has `b64Challenge` and `ErrClientDataType … ErrCredentialIDMismatch`; `approver.go:35 Passkey *Passkey`. Control in the same command: `ApproveByAssertion` → **no match** (`grep -rn ApproveByAssertion host cmd` empty; `grep -rn ApproveBySignature host` finds `signed.go:55`) |
| C4 | The store at schema 5 | `credentials.go:42 const credentialsSchemaVersion = 5`; `:73 case version >= 5:` (4 tables); `:358 if version < 5 { … approversSQL }`; `mint_requests.go` C1 `:306`, C2 `:309`, C3 `:312`, C4 `:315`, C5 `:324`, `BeginTx` `:327`, C8 `:348` — **design §2.10's line numbers hold exactly** |
| C5 | The two arms §2.10 names | `consume_binding_test.go:138 {"approver-format-not-ssh", … rc.Approver.Format = "webauthn" }` (expects `ErrInconsistentApproval`), `:288 {"other-format", …` (`Format, ApprovedBy = "webauthn", "webauthn:x"`, expects `ErrApproverKeyChanged`) |
| C6 | Three tests pin the 4-table set | `grep -n requireExactTableNames host/store/{credentials,approvers,mint_requests}_test.go` → `credentials_test.go:67`, `approvers_test.go:113`, `mint_requests_test.go:132`, each `[]string{"approvers", "episode_starts", "mint_requests", "session_credentials"}` |
| C7 | **The schema bump breaks four `cmd/ailang-worldd` arms** (work, before the routing of §1.4 P2) | `go test ./cmd/ailang-worldd/` → `FAIL`: `TestApproveBySignatureRefusalsBeforeAnyEffect/v4-store`, `TestApproverEnrollRefusals/v4`, `TestApproverRevoke/v4`, `TestApproverList`, each `user_version 4 with tables [episode_starts mint_requests session_credentials webauthn_credentials] is not a credentials store`; their fixture builders are `approver_rig_test.go:67` and `approve_sig_test.go:219` (`DROP TABLE approvers`, `PRAGMA user_version = 4` on a fresh, now schema-6, store). `host/daemon` whole package on the same tree: `ok` (282 s) |
| C8 | The version-compare sites outside `host/store` | `grep -rn CredentialsSchemaVersion cmd host` (non-store) → `cmd/ailang-worldd/episode_preset.go:235`, `session_new.go:380`, `host/daemon/credentials.go:61`: each compares with the exported constant, so a schema-5 store now asks for one `serve` restart (design §2.10 "The cost"); their packages stay green (C7, C27) |
| C9 | The deadline inventory | `credentials_test.go:260 func TestCredentialsMethodsRequireADeadline` lists 13 methods; `context_read_test.go` pins only the seven world-store getters by regex, so the four new credentials methods appear in this test only |
| C10 | `ErrApproverFormat`'s text is pinned in `cmd` | `grep -rn 'format must be ssh' cmd` → `approver_enroll_test.go:234`, `approve_sig_test.go:212`: P3w gets its own error (`ErrApproverNotPasskey`, §1.3 S2-D3) |
| C11 | **SOURCE.txt names files that do not exist, and XW0 does not see it** (judge carry-forward) | `grep -n 'es256.json\|eddsa.json' host/approval/testdata/webauthn/SOURCE.txt` → `17:  es256.json  sha256 c037…`, `18:  eddsa.json  sha256 1544…` while `ls testdata/webauthn` → `chrome-eddsa.json chrome-es256.json SOURCE.txt generate.sh`; XW0 (`webauthn_fixture_test.go:99-107`) checks three substrings of SOURCE.txt and pins no hash of it; base control **BC-SOURCE** (§4): renaming that line to `zz-not-a-fixture.json` at the base → `TestPasskeyFixturesAreTheGeneratedBytes` **SURVIVED** |
| C12 | **`Passkey.CredentialID`/`Key` aliasing** (judge carry-forward) | `NewPasskeyApprover` reads `in.Response` only through `json.Unmarshal` + `base64…DecodeString` (`passkey_json.go:38`, `:53`), which allocate; `authdata.go:64` copies the credential ID and `cose.go:78` the coordinates out of that buffer. So the result cannot alias the caller's input: MUT-AD-ALIAS-EQUIV and MUT-COSE-ALIAS-EQUIV (dropping both copies) SURVIVE the new pin, as predicted. The aliasing that **can** happen is the assertion law's: `SignedApproval.credentialID` is copied from the caller's `Approver.Passkey.CredentialID` — MUT-SA-CRED-ALIAS and MUT-PU-CLONE are KILLED by PA4 |
| C13 | **`RegistrationInput.Challenge` zero / `Now` unchecked** (judge carry-forward) — a throwaway probe test (deleted) on work | `NewPasskeyApprover` with `Challenge: [32]byte{}` and a registration made for it, `Now: 0` → `err=<nil> EnrolledAt=0`; `Now: -5` → `err=<nil> EnrolledAt=-5` (`PROBE` lines). Routed by S2-D11 |
| C14 | The import guard's scope | `imports_test.go:28 os.ReadDir(".")`: it reads `host/approval` only; `approvaltest/passkey.go` (test support) may use `crypto/rand` for ES256's random signature nonce; `TestApprovalTestIsTestOnly` keeps it out of production |
| C15 | Go 1.26 APIs the test authenticator uses | `go doc crypto/ecdsa ParseRawPrivateKey` → `func ParseRawPrivateKey(curve elliptic.Curve, data []byte) (*PrivateKey, error)`; `go doc crypto/ecdsa PublicKey.Bytes` → `func (pub *PublicKey) Bytes() ([]byte, error)` |
| C16 | Chrome's assertion envelopes have **no** `userHandle` member | `jq '.assertion.response \| keys'` (via python) on both fixtures → `['authenticatorData', 'clientDataJSON', 'signature']`; registration responses add `publicKey`, `publicKeyAlgorithm`, `transports` (non-string members: XW1 decodes them as raw JSON) |
| C25 | **Every repo-wide inventory, on the final tree** (work, `-v`, tee'd) | the 16 design-§7 tests + `TestSignedApprovalFieldsAreAllBound`, `TestEnrolmentPromptShowsEveryPersistedField`, `TestRequestConsumeCompleteness`, `TestApproverRowMatchesTheTable`, **`TestPasskeyRowMatchesTheTable`**, `TestRequestMintFieldsAreSourced` → all 22 printed `--- PASS`, 0 SKIP/FAIL (`plan-scratch/logs/final_gates/inv.log`, the §6 script run on the committed final tree); no inventory changed: no exec site, module, package, route, context root or `t.Parallel` was added |
| C26 | Whole packages on the final tree | `go test -count=1 ./host/approval/... ./host/store/ ./host/authority/ ./cmd/ailang-worldd/ ./host/daemon/` → all `ok` (approval 0.4 s, approvaltest 0.2 s, store 11.4 s, authority 0.7 s, `cmd/ailang-worldd` 113.6 s, `host/daemon` 283.3 s; `logs/final_gates/pkgs.log`). Other tests that walk `cmd/` or `host/` (the AC10 census, `TestProductionGoSurface`, `TestNoNewDeadlineFreeStoreReads`, `TestApprovalTestIsTestOnly`, the boundary allowlists) are in C25; the whole repository is C27 |
| C27 | **The whole repository on the final tree** (committed in the stage repo) | `VERIFY_GO_LEGS=plain ./scripts/verify_go.sh` on `plan-scratch/stages` at tag `M5b` (committed; `diff -r cmd host` equal to work) → **rc 0**: `✓ 0 other binary blobs among 6687 tracked files` (10 new text files), toolchain floor, `✓ all 64 required top-level evidence tests passed exactly once`, 62 packages `ok` (`host/broker` 320 s, `host/daemon` 355 s, `host/verifygate` 328 s, `cmd/ailang-worldd` 212 s, `host/store` 17 s), `✓ go gate PASSED (plain)`; `./scripts/verify_ail.sh` → rc 0 (`✓ verify gate PASSED`). Whole §6 script: rc 0 in 775 s (`logs/final_gates/`). This run had no `host/broker` red (contrast §1.1) |
| C28 | `go vet ./...`, gofmt and `-race` on the final tree | `go vet ./...` rc 0, empty; `gofmt -l cmd host` empty; `go test -race -count=1 ./host/approval/... ./host/store/ ./host/authority/` rc 0 (store 31.5 s), no race report (`logs/final_gates/{vet,gofmt,race}.log`) |
| C30 | **The patches reproduce the mutant tree** | a fresh `git archive 3b77c294` export + `patch -p1` of Appendix P's eight patches **extracted from this file with §2's `awk`** (each `cmp`-equal to the generated patch; `--dry-run` first; no `.orig`/`.rej`) → `diff -r cmd host` against the tree §4's mutants ran on (`plan-scratch/work`): **identical** (`ROUNDTRIP_IDENTICAL_TO_MUTANT_TREE`, `tools/roundtrip.sh`, `logs/roundtrip.log`); every test added late in the planning (the PS3 read-back arm, the PS1 literal pin, the split test files) is in the patches, because the comparison is byte-for-byte over `cmd/` and `host/`; nothing left uncommitted |
| C31 | **Every sub-milestone is bisectable** | the same fresh export, after **each** patch: `tools/accept.sh <M>` = §3's acceptance verbatim (gofmt empty, build + vet clean, the compile fence, the named `-v` tests with no `--- SKIP`/`--- FAIL` and at least one `--- PASS`, M4a's whole `cmd/ailang-worldd`, then `[pkgs]`) → `ACCEPT <M> rc=0` for all eight; the cumulative numstat after a commit of each = §7's value: **345, 672, 890, 1027, 1300, 1538, 1777, 1843** (`NUMSTAT` lines of `logs/roundtrip.log`) |

### 1.3 Planner decisions where the design leaves a choice

- **S2-D1 — The cut** (§1.5): none needed; the fallback was priced by measurement and not fired.
  Sub-milestones are split at ≤ ~250 lines where a split leaves each green on its own (M3b1 is 345: the
  law, its accessor and RC1 cannot be separated, because RC1 perturbs a verified passkey approval).
- **S2-D2 — `checkApprover(admit, a, format)`** (design §2.9) lives in `signed.go` beside
  `ApproveBySignature`, which now calls it for P2–P5; its order is P2, **P3w**, P3, P4, P5. P3w is
  `format == FormatWebAuthn && (a.Format != format || a.Passkey == nil)`; P3 is unchanged for ssh. So
  an ssh approver is judged exactly as before (`TestApproveBySignatureOrder` unedited), and the
  shared mutants die by both laws' tests (MUT-P3W-SCOPE by SA1, MUT-P2/P4/P5 by PA1).
- **S2-D3 — P3w's error is new, `ErrApproverNotPasskey`**, because `ErrApproverFormat`'s text ("must be
  ssh") is pinned by two `cmd` tests (C10) and would mislead on the passkey path.
- **S2-D4 — The W-errors.** W3 `ErrUserHandle` (a present userHandle, even empty, that is not
  `[]byte(Name)`); W13 `ErrBackupEligibility`; W15 `approval.ErrSignCountRegression`, wrapped with
  "enrolled N, this assertion M; if you did not cause this, run: ailang-worldd approver revoke <name>"
  (design §2.9's "names `approver revoke <name>` as the fix"; PA3 pins the text). W6 refuses when
  `MessageHash` itself fails (a negative expiry) with `ErrChallenge`. W14 reuses `COSEKey.verify`
  (`ErrBadSignature`).
- **S2-D5 — `SignedApproval` from a passkey**: `format "webauthn"`, `key = Passkey.Key.Canonical()`
  (what `approvers.public_key` holds, design §2.4), `credentialID = bytes.Clone(Passkey.CredentialID)`,
  `signCount` = the assertion's. `PasskeyUse()` returns a copy (B7: unverified → false; B8: any
  format but webauthn → false).
- **S2-D6 — `approvaltest.Passkey`'s API** (design §2.15 sketches `Register(origin, rpID, name, …)`):
  `NewPasskey(alg, seed)`, `CredentialID()`, `Key()`, `Register(origin, rpID, challenge, …Edit)`,
  `Assert(origin, rpID, challenge, counter, …Edit)`. No `name` parameter: neither recorded response
  carries the user (C16), so a name would be unused. `Edit` is one struct of knobs (`Type`,
  `CrossOrigin`, `TopOrigin`, `ClientData`, `Set`/`Clear` flag bits, `Extensions`, `UserHandle`,
  `RawID`, `BadSignature`); the zero `Edit` is Chrome's layout (flags `0x45`/`0x05`, registration
  counter 1, AAGUID `0102…08`, a 32-byte credential ID, member order `type, challenge, origin,
  crossOrigin`). XW1 pins that layout against both recordings.
- **S2-D7 — PasskeyRow, EnrolledPasskey and the lookups.** `PasskeyRow` is the design's. The lookups
  share one join (`enrolledPasskeys`): `PasskeyByCredentialID(ctx, id) (ApproverRow, PasskeyRow, bool,
  error)` as designed, `Passkey(ctx, name) (PasskeyRow, bool, error)` (revoked included), and
  `ActivePasskeys(ctx) ([]EnrolledPasskey, error)` with `EnrolledPasskey{Approver ApproverRow; Passkey
  PasskeyRow}` (the design names the content, not the type), filtered `revoked_at IS NULL AND format
  = 'webauthn'`, ordered by name, at most 64.
- **S2-D8 — `EnrolPasskey`'s validation is ten branches, S1–S10** (§1.4 P3), and its three
  pre-checks are one table-driven loop (name, key, credential ID — in that order — each naming the
  holder). Both `INSERT`s are in one transaction.
- **S2-D9 — Carry-forward, aliasing** (C12): PA4 pins that editing the input approver's credential ID
  and key, or the slice `PasskeyUse()` returned, does not change the approval;
  `TestNewPasskeyApproverDoesNotAliasItsInput` pins that overwriting the registration bytes after
  `NewPasskeyApprover` does not change the `Approver` (its two mutants are named equivalents: the
  input is decoded, never retained).
- **S2-D10 — Carry-forward, SOURCE.txt** (C11): the two lines now read `chrome-es256.json  from the
  probe's es256.json, sha256 …` (and EdDSA's), so the committed names and the probe outputs' names
  are both stated; XW0 additionally requires `chrome-es256.json  from` and `chrome-eddsa.json  from`
  (MUT-SOURCE-NAMES, KILLED; base control BC-SOURCE, SURVIVED at the base). The two fixture SHA pins
  are unchanged; SOURCE.txt itself was never hash-pinned and is not now.
- **S2-D11 — Carry-forward, `RegistrationInput.Challenge` zero / `Now` unchecked** (C13): **decided as
  the caller's job, no new branch in the pure law.** The challenge is the enrol process's
  `env.random` (sprint 3, WN5 refuses a failed or nil reader); `Now ≤ 0` cannot be stored: M4b's
  `EnrolPasskey` S4 refuses `EnrolledAt <= 0` (PS2 arms `S4-enrolled-0`, `S4-enrolled--1`;
  MUT-S4, MUT-S4-LT KILLED) and the DDL has no row without it. Sprint 3's planner may make WN5 also
  refuse an all-zero challenge (Appendix C).
- **S2-D12 — The upgrade test is renamed** `TestCredentialsUpgradeToV5` → `TestCredentialsUpgradeToV6`
  (design PS4) and gains a v5 arm (an ssh approver revoked at 950 and a start record kept); it counts
  `webauthn_credentials` rows with SQL, so M4a does not depend on M4c's lookups.
- **S2-D13 — PS1's version pin is the literal 6** (`version != 6`): the old pin compared the store
  with `credentialsSchemaVersion` itself, so MUT-CR-VERSION (6 → 5) survived it (first mutant pass);
  `go vet` refuses the two-constant form (`suspect or`).
- **S2-D14 — `passkeys.go` is written in stage sections** (`// ---- M4a ----` … `// ---- M5a ----`,
  the convention of `approvers.go`), so each sub-milestone's patch appends one section.

### 1.4 Premises found false or incomplete at HEAD (each routed around minimally)

- **P1 — Design §9 / sprint 1 Appendix C: M3b ≈ 450 lines, sprint 2 ≈ 1,600 (≈ 2,050 together).**
  Measured: M3b **890** (157 production, 563 tests, 164 test support, 6 testdata), M4 **648**, M5
  **305**: **1,843**. M3b is twice its estimate (the test authenticator alone is 164 lines, PA1 with
  flag-state coverage 218), M4 and M5 are under theirs. No routing needed: the total is under the cap
  (§1.5).
- **P2 — Design §8: sprint 2 touches no `cmd/` file.** False: the schema bump turns four
  `cmd/ailang-worldd` arms red (C7), because two test fixture builders take a fresh store "back to
  schema 4" by dropping only `approvers`; at schema 6 the fresh store also holds
  `webauthn_credentials`. Routed minimally in M4a: each builder's statement list gains
  `DROP TABLE webauthn_credentials` first (2 lines each, test-only; one comment line). No production
  file outside `host/` changes; the whole `cmd/ailang-worldd` package is in M4a's acceptance.
- **P3 — Design §2.10: "each its own `ErrBadApprover` branch (S1–S9)"** lists ten conditions (format;
  names differ; empty name or key; `EnrolledAt`; `LastUsedAt`/`RevokedAt`; credential ID length;
  `Alg`; AAGUID length; `RPID`/`Origin`; `SignCount`). Routed: numbered **S1–S10**, each its own `case`,
  each compound condition with one arm per operand (MUT-S3-NAME … MUT-S10-MAX).
- **P4 — Design §6 PS1 "user_version 6" was pinned only through the constant** (S2-D13): the existing
  check `version != credentialsSchemaVersion` is tautological for a fresh store. Routed: the literal.
- **P5 — Design §6 PA1 builds every arm "from a test-authenticator knob", and §2.15's
  `Register`/`Assert` take a `name`.** The recorded responses carry no user (C16); the name is used
  only by W3's `userHandle`, which is a knob. Routed by S2-D6 (no `name` parameter).
- **P6 — Design §2.11: "There is deliberately no separate 'webauthn without PasskeyUse' refusal".**
  True, and kept: A6 is one branch; MUT-A6 shows the store's C13 is what refuses when it is skipped
  (AU3 expects success).
- **P7 — The planner's first mutant pass found three test gaps, each closed before the final run**:
  MUT-CR-VERSION survived (P4); MUT-CR-TABLES6 survived its first named test (PS1 opens a fresh file,
  which is initialised, never re-checked) and is now named against `TestCredentialsUpgradeToV6`,
  whose second `UpgradeCredentials` reopens the version-6 table set; four drafts did not compile
  (`h`, `cdh`, `bytes` left unused; `go vet`'s `suspect or`) and were rewritten into compiling forms.
  Only the final run counts (§4).
- **P8 — A gate that ran under `nohup` was red for the harness's reason, not the repo's** (§1.1):
  routed by §7's "never under `nohup`" and §6's re-run rule.

### 1.5 The size, and the cut

Lines are insertions + deletions under `cmd/ host/`: **`git diff --numstat <base> HEAD -- cmd host`**
(`.gitattributes` and docs are outside the filter and this sprint changes neither). **Every row is
measured** from the scratch stages, each committed in a throwaway repo (`plan-scratch/stages`, tags
`M3b1` … `M5b`).

| Sub-milestone | Production | Test | Test support + testdata | Step | Cumulative vs `3b77c294` (after its commit) | Run |
|---|---|---|---|---|---|---|
| M3b1 — the assertion law, `checkApprover`, `PasskeyUse`; PA2, PA3 (Chrome), PA4, RC1, the two carry-forward pins | 157 | 182 | 6 | **345** | 345 | 1 |
| M3b2 — the test authenticator; TP1, XW1 | 0 | 163 | 164 | **327** | 672 | 1 |
| M3b3 — PA1 (order, flag-state coverage), PA3 | 0 | 218 | 0 | **218** | 890 | 2 |
| M4a — schema 6: the table, the upgrade, PS1, PS4; the two `cmd` fixture builders | 59 | 78 | 0 | **137** | 1,027 | 2 |
| M4b — `EnrolPasskey` (S1–S10, three pre-checks, both inserts); PS2, PS5 | 111 | 162 | 0 | **273** | 1,300 | 3 |
| M4c — `PasskeyByCredentialID`, `Passkey`, `ActivePasskeys`; PS3, RC6 | 71 | 167 | 0 | **238** | 1,538 | 3 |
| M5a — the consume: C13–C18; CS7–CS9, CS3's edit, RC3 | 55 | 184 | 0 | **239** | 1,777 | 4 |
| M5b — authority A6; AU3 | 3 | 63 | 0 | **66** | **1,843** | 4 |
| **Sprint** | **456** | **1,217** | **170** | **1,843 ≤ 2,000** | | |
| *the fallback, priced and not taken* — M5b to sprint 3 | — | — | — | *−66* | *1,777* | — |

---

## 2. How the executor applies a sub-milestone (the patch protocol)

Each sub-milestone `<M>` is the block fenced by a line `` ````diff <M> `` (four backticks) and a line of
four backticks, in Appendix P. Paths are `a/<repo path>` and `b/<repo path>`. Apply it **only** this
way, in one `bash` shell (the controller's shell is zsh):

```bash
bash <<'EOF'
set -euo pipefail
M=M3b1   # <- the sub-milestone, exactly as in the Appendix P heading: M3b1 M3b2 M3b3 M4a M4b M4c M5a M5b
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter380-row201-a2s2
PLAN="$W/design_docs/planned/w-world-approvals-by-key-a2-sprint2-plan.md"
T="$HOME/.ailang/state/world-iter380/exec-patches"; mkdir -p "$T"
awk -v m="$M" '$0 == "````diff " m {f=1; next} f && $0 == "````" {exit} f' "$PLAN" > "$T/$M.patch"
test -s "$T/$M.patch"
cd "$W"
patch -p1 --dry-run < "$T/$M.patch"
patch -p1 < "$T/$M.patch"
if find cmd host \( -name '*.orig' -o -name '*.rej' \) | grep .; then echo "STOP: patch left .orig/.rej"; exit 1; fi
git status --short
EOF
```

- If `--dry-run` fails, **stop and report** (no hand edits, no fuzz options). It can only fail if an
  earlier sub-milestone was skipped or changed.
- Then run that sub-milestone's **Acceptance** (§3) in one bash shell with
  `export AILANG_BIN=$HOME/.pinned-ailang/ailang`, each command's output `tee`'d to a log under
  `$HOME/.ailang/state/world-iter380/exec-logs/` (never only through `tail`).
- Then snapshot (§7). The controller commits and compares the numstat with §7's expected value.

---

## 3. Milestones

Every sub-milestone's acceptance starts with the same three lines (abbreviated **[base]** below), and
ends with the package line (abbreviated **[pkgs]**):

```
test -z "$(gofmt -l cmd host)"
go build ./... && go vet ./host/approval/... ./host/store/ ./host/authority/ ./cmd/ailang-worldd/
go test -run '^$' ./...
...
go test ./host/approval/... ./host/store/ ./host/authority/ -count=1          # [pkgs]
```

Each `-v` run below must show `--- PASS` for every test it names and **no `--- SKIP` or `--- FAIL`**.

### M3b1 — the assertion law, `checkApprover`, `PasskeyUse` · 157 prod + 182 test + 6 testdata · Run 1

**Patch** `M3b1`:
- `host/approval/signed.go`: `SignedApproval` gains `credentialID []byte`, `signCount uint32`;
  `checkApprover(admit, a, format)` (P2, **P3w**, P3, P4, P5; S2-D2) and `ApproveBySignature` calls
  it; `PasskeyUse() (credentialID []byte, signCount uint32, ok bool)` (B7, B8; a copy).
- new `host/approval/passkey_assert.go`: `AssertionInput`, `ApproveByAssertion` (P1, P2–P5 via
  `checkApprover`, W1–W15 in the design's order), `ErrApproverNotPasskey`, `ErrUserHandle`,
  `ErrBackupEligibility`, `ErrSignCountRegression` (S2-D3, S2-D4, S2-D5).
- new `host/approval/passkey_assert_test.go` (package `approval`): **PA2**
  `TestApproveByAssertionAcceptsChrome` (both recordings approve the fixture request:
  `Binds`, `ApprovedBy() == "webauthn:mark-mac"`, `Use()`'s canonical key, `PasskeyUse() = (rawId, 2,
  true)`), **PA3/XW2 last arm** `TestChromeAssertionAfterAssertion2IsRefused` (stored counter = the
  recorded `assertion2`'s 3 → `ErrSignCountRegression`), **PA4** `TestSignedApprovalPasskeyUse` (nil,
  zero, unverified webauthn, verified ssh, verified ssh carrying a counter → `(nil, 0, false)`; the
  passkey approval unchanged after editing the returned ID, the input's credential ID and key), and
  `TestNewPasskeyApproverDoesNotAliasItsInput` (S2-D9).
- `host/approval/signed_test.go`: **RC1** classifies `credentialID` and `signCount`, observes
  `PasskeyUse`, and perturbs every field in both a verified ssh and a verified passkey approval.
- `host/approval/testdata/webauthn/SOURCE.txt` and `webauthn_fixture_test.go` (XW0): S2-D10.

**Acceptance**
```
[base]
go test ./host/approval/ -run '^(TestApproveByAssertionAcceptsChrome|TestChromeAssertionAfterAssertion2IsRefused|TestSignedApprovalPasskeyUse|TestNewPasskeyApproverDoesNotAliasItsInput|TestSignedApprovalFieldsAreAllBound|TestApproveBySignatureOrder|TestSignedApprovalAccessorsFailClosed|TestPasskeyFixturesAreTheGeneratedBytes)$' -count=1 -v
[pkgs]
```

**Expected cumulative numstat after the controller's commit of M3b1: 345**.

### M3b2 — the test authenticator · 163 test + 164 test support · Run 1

**Patch** `M3b2`:
- new `host/approval/approvaltest/passkey.go`: `Passkey`, `NewPasskey(alg, seed)` (EdDSA
  `ed25519.NewKeyFromSeed`; ES256 `ecdsa.ParseRawPrivateKey(P256, seed×32)`), `CredentialID()`,
  `Key()`, `Edit`, `Register`, `Assert` (S2-D6). Test support: imported only by `_test.go` files.
- new `host/approval/passkey_testauth_test.go` (package `approval_test`, sprint 1 §1.4 P4): the PA
  request (`paOrigin = https://x.ts.net:8443`, digest `{1}`, challenge `{2}`, expiry 1600, `paHash`);
  **TP1** `TestTestPasskeyVerifies` (both algorithms, a tailnet and a localhost origin: `Register` →
  `NewPasskeyApprover`, `Assert` → `ApproveByAssertion`, the credential ID and canonical key carried);
  **XW1** `TestTestPasskeyMatchesChromeLayout` (envelope members ⊆ Chrome's, `clientDataJSON` member
  order, the attestationObject's 30 bytes outside authData, the registration authData with the
  credential ID and x, y masked, the 37-byte assertion authData, for both recordings).

**Acceptance**
```
[base]
go test ./host/approval/... -run '^(TestTestPasskeyVerifies|TestTestPasskeyMatchesChromeLayout|TestApprovalTestIsTestOnly|TestApprovalImportsOnlyTheStandardLibrary)$' -count=1 -v
[pkgs]
```

**Expected cumulative numstat after the controller's commit of M3b2: 672**.

### M3b3 — the assertion law's refusals and its counter law · 218 test · Run 2

**Patch** `M3b3`:
- new `host/approval/passkey_law_test.go` (package `approval_test`): `paApprover`, `paInput`, `paCase`;
  **PA1** `TestApproveByAssertionOrder`: for **both algorithms**, 34 arms over P1–P5, P3w, W1–W10, W14,
  W15 including the order pairs (P1 before P4, P4 before W14, W5 before W6, W7 before W10, W12 before
  W14, W13 before W14, W14 before W15) and W6's negative-expiry arm; **16 flag-state arms** for W11–W13
  (each refusal with UP, UV, BE, BS, ED in every other state that reaches it, and with the enrolled
  `BackupEligible` both ways); and 6 acceptance controls (clean, enrolled-at-now, `userHandle` = the
  name, an extensions map, BE and BE|BS enrolled). **PA3** `TestSignCountLaw`: stored/new 1/2, 0/0,
  0/7, 0/1, 7/2³²−1, 0/2³²−1 accepted with the new counter carried; 2/2, 3/2, 5/0, 1/0, 2³²−1/0,
  2³²−1/2³²−1 refused, the error naming `approver revoke mark-mac`.

**Acceptance**
```
[base]
go test ./host/approval/ -run '^(TestApproveByAssertionOrder|TestSignCountLaw)$' -count=1 -v
[pkgs]
```

**Expected cumulative numstat after the controller's commit of M3b3: 890**.

### M4a — credentials schema 6 · 59 prod + 78 test · Run 2

**Patch** `M4a`:
- `host/store/credentials_schema.sql`: the `webauthn_credentials` table (design §2.10, verbatim) and
  its comment; the `approvers` comment no longer promises `ADD COLUMN`.
- `host/store/credentials.go`: `credentialsSchemaVersion = 6`; `credentialsTables(6)` = 5 tables,
  `credentialsTables(5)` the old 4; `upgradeCredentials` gains `if version < 6 { passkeysSQL }`.
- new `host/store/passkeys.go` (section M4a): `passkeysTable`, `passkeysSchemaVersion = 6`, `passkeysSQL`.
- `host/store/credentials_test.go`: **PS1** (`canonicalPasskeysDDL`; the five-table set; the schema
  file and `passkeysSQL` both equal it; `user_version` 6, S2-D13).
- `host/store/approvers_test.go`: `makeV5Credentials`; **PS4** `TestCredentialsUpgradeToV6` (renamed;
  v1, v3, v4 and **v5** upgrade once, idempotently, to the five tables with every row kept and no
  passkey row; a session open leaves each version alone).
- `host/store/mint_requests_test.go`: `TestCredentialsUpgradeToV4`'s table list gains the fifth table.
- `cmd/ailang-worldd/approver_rig_test.go`, `approve_sig_test.go`: the v4 fixture builders also drop
  `webauthn_credentials` (§1.4 P2).

**Acceptance**
```
[base]
go test ./host/store/ -run '^(TestCredentialsSchemaMatchesCanonicalDDL|TestCredentialsUpgradeToV6|TestCredentialsUpgradeToV4|TestCredentialsV5WithoutApproversIsRefused|TestCredentialsV4WithoutMintRequestsIsRefused)$' -count=1 -v
go test ./cmd/ailang-worldd/ -run '^(TestApproveBySignatureRefusalsBeforeAnyEffect|TestApproverEnrollRefusals|TestApproverRevoke|TestApproverList)$' -count=1 -v
go test ./cmd/ailang-worldd/ -count=1                                          # ~2 min: the schema bump reaches the CLI
[pkgs]
```

**Expected cumulative numstat after the controller's commit of M4a: 1027**.

### M4b — `EnrolPasskey` · 111 prod + 162 test · Run 3

**Patch** `M4b`:
- `host/store/passkeys.go` (section M4b): `maxSignCount`, `maxActivePasskeys`, `PasskeyRow`,
  `ErrPasskeyCredentialEnrolled`, `ErrSignCountRegression`, `ErrPasskeyChanged`, `olderForPasskeys`
  (`Need: 6, Feature: "passkeys"`), `checkPasskeyRows` (S1–S10), `EnrolPasskey` (deadline; below v6 →
  `*OlderCredentialsStoreError`; S1–S10; in one transaction the name, key and credential-ID
  pre-checks naming the holder, then both `INSERT`s).
- `host/store/credentials_test.go`: `TestCredentialsMethodsRequireADeadline` gains `EnrolPasskey`.
- new `host/store/passkeys_test.go`: `testPasskey`, `rowCounts`, `openV5`, `isOlderPasskeys`; **PS2**
  `TestEnrolPasskeyRefusals` (v5 → `Older{Version 5, Need 6, "passkeys need schema 6"}` and no row; 20
  `ErrBadApprover` arms over S1–S10, every operand, each with `0 approvers, 0 passkeys`; the inclusive
  edges 1- and 1,023-byte credential IDs, counter 0 with EdDSA, counter 2³²−1 enrol; then a name
  taken (by a passkey and by an ssh approver), a key enrolled (both), a credential enrolled, each
  naming the holder with the counts unchanged, again after both holders are revoked); **PS5**
  `TestEnrolApproverStillSSHOnly`.

**Acceptance**
```
[base]
go test ./host/store/ -run '^(TestEnrolPasskeyRefusals|TestEnrolApproverStillSSHOnly|TestCredentialsMethodsRequireADeadline)$' -count=1 -v
[pkgs]
```

**Expected cumulative numstat after the controller's commit of M4b: 1300**.

### M4c — the passkey lookups · 71 prod + 167 test · Run 3

**Patch** `M4c`:
- `host/store/passkeys.go` (section M4c): `EnrolledPasskey`, `passkeyColumns`, `scanPasskey`,
  `enrolledPasskeys` (deadline; below v6 → `Older`; one join on `name`), `PasskeyByCredentialID`,
  `Passkey`, `ActivePasskeys` (S2-D7).
- `host/store/credentials_test.go`: `TestCredentialsMethodsRequireADeadline` gains the three lookups.
- new `host/store/passkey_lookups_test.go`: **PS3** `TestPasskeyLookups` (exact rows by credential ID
  with `last_used_at` and `revoked_at` read back; misses on a prefix, a longer ID and nil; `Passkey`
  of a revoked approver found, of an ssh approver, `""`, `a%`, `zz` missed; `ActivePasskeys` = `a,c`,
  still `a,c` beside a hand-made passkey row of an ssh approver, then 64 of 72 by name; a v5 store
  refuses all three); **RC6** `TestPasskeyRowMatchesTheTable` (fields ↔ map ↔ `PRAGMA table_info`; a
  row of distinct non-zero values read back field by field through all three lookups).

**Acceptance**
```
[base]
go test ./host/store/ -run '^(TestPasskeyLookups|TestPasskeyRowMatchesTheTable|TestCredentialsMethodsRequireADeadline)$' -count=1 -v
[pkgs]
```

**Expected cumulative numstat after the controller's commit of M4c: 1538**.

### M5a — the consume binds the passkey (C13–C18) · 55 prod + 184 test · Run 4

**Patch** `M5a`:
- `host/store/mint_requests.go`: `ApproverUse.Passkey *PasskeyUse`; **C13** (`(Format == "webauthn")
  != (Passkey != nil)` → `ErrIncompleteConsume`), **C14** (credential ID 1–1023 bytes, counter
  0..2³²−1), **C15** (below v6 → the passkeys refusal) — inserted **after C5 and before `BeginTx`**, so
  C4 still refuses the unedited `approver-format-not-ssh` arm (MUT-C13-BEFORE-C4); inside the
  transaction, after C8's `useApprover`, `usePasskey`.
- `host/store/passkeys.go` (section M5a): `PasskeyUse{CredentialID, SignCount}`; `usePasskey` (the
  design's guarded `UPDATE … WHERE name = ? AND credential_id = ? AND (sign_count < ? OR (sign_count =
  0 AND ? = 0))`; when it changes no row, **C16** absent → `ErrApproverUnknown`, **C17** another
  credential → `ErrPasskeyChanged`, **C18** otherwise → `ErrSignCountRegression` naming `approver
  revoke`); any refusal rolls back the request and `last_used_at`.
- `host/store/consume_binding_test.go`: **the one named edit** — CS3 `other-format` gains
  `Passkey: &PasskeyUse{32 bytes, 1}` and still expects `ErrApproverKeyChanged`; **RC3** classifies
  `ApproverUse.Passkey`, `PasskeyUse.CredentialID` (required: zeroed → `ErrIncompleteConsume`) and
  `PasskeyUse.SignCount` (bounded; 0 is a counter). CS2 `approver-format-not-ssh` is **unedited**.
- new `host/store/passkey_consume_test.go`: **CS7** `TestPasskeyConsumeCompleteness` (C13 both ways, C14
  four operands, C15 on a v5 store: each before the clock, request pending, no session, counter and
  `last_used_at` unchanged); **CS8** `TestPasskeyConsumeAdvancesTheCounter` (1/2, 0/0, 0/7,
  1/2³²−1, 1- and 1,023-byte credentials accepted: counter = new, `last_used_at` 1010, `approved_by
  webauthn:x`, 1 session; 2/2, 3/2, 5/0, 2³²−1/2³²−1 refused with nothing changed; the same counter on
  a second request refused); **CS9** `TestPasskeyConsumeChecksTheRowInTheTransaction` (a second
  connection inside the `Clock` hook revokes x, rewrites its credential ID, deletes its passkey row,
  or raises its counter to 9 → `ErrApproverRevoked`, `ErrPasskeyChanged`, `ErrApproverUnknown`,
  `ErrSignCountRegression`, each rolling everything back).

**Acceptance**
```
[base]
go test ./host/store/ -run '^(TestPasskeyConsumeCompleteness|TestPasskeyConsumeAdvancesTheCounter|TestPasskeyConsumeChecksTheRowInTheTransaction|TestRequestConsumeCompleteness|TestSignedConsumeChecksTheApproverInTheTransaction|TestApprovedByMatchesTheApprover|TestTypedConsumeTouchesNoApprover)$' -count=1 -v
[pkgs]
```

**Expected cumulative numstat after the controller's commit of M5a: 1777**.

### M5b — authority derives the passkey use (A6) · 3 prod + 63 test · Run 4

**Patch** `M5b`:
- `host/authority/mint.go`: **A6** — after A3, `if id, count, ok := rm.Signed.PasskeyUse(); ok {
  consume.Approver.Passkey = &store.PasskeyUse{CredentialID: id, SignCount: int64(count)} }`.
- new `host/authority/mint_passkey_test.go`: **AU3** `TestMintForRequestDerivesPasskeyUse`
  (`approvaltest.Passkey` enrolled through `EnrolPasskey` with the canonical key; its assertion of
  the stored request through `ApproveByAssertion`; a zero `SignedApproval` → `ErrNotVerified`, nothing
  minted; the verified one → `approved_by = webauthn:mark-mac`, counter 5, 1 session).

**Acceptance**
```
[base]
go test ./host/authority/ -run '^(TestMintForRequestDerivesPasskeyUse|TestMintForRequestRecordsTheSigner|TestMintForRequestRefusesAnUnboundApproval|TestRequestMintFieldsAreSourced)$' -count=1 -v
[pkgs]
```

**Expected cumulative numstat after the controller's commit of M5b: 1843**.

---

## 4. Mutation table (139 rows run on the planner's final tree: 135 KILLED, 4 named equivalents SURVIVED as predicted; plus 2 base controls)

**Derivation.** Anchored to Appendix P's production hunks (`signed.go`, `passkey_assert.go`, `credentials.go`, `credentials_schema.sql`, `passkeys.go`, `mint_requests.go`, `authority/mint.go`): one row per new `if`, `case`, condition operand, bound, copied field, SQL guard clause and constant a test can tell apart; plus design §6's named mutants (C13–C15 moved before C4; the counter `UPDATE` deleted, widened to `<=`, its both-zero clause dropped; `RowsAffected` unchecked; the classify arms; A6 deleted), the retro's flag-state mutants (`!a && !b` forms over UP, UV, BE, BS for W11–W13; each operand of S3, S5, S6, S9, S10, C14), the carry-forward rows, and one pre-existing branch the new PS5 pins. **Stage** is the sub-milestone whose patch adds the branch; its killing test may land later in the same sprint (PA1, PA3 in M3b3 kill most M3b1 rows). **Procedure**: `plan-scratch/tools/mutate.py <tree> tools/mutants.json` (outside the repo; generated by `tools/mk_mutants.py`, which asserts every anchor occurs exactly once): apply **one** edit, assert the file changed, `go vet` the package (a mutant that does not compile is not a mutant: the first pass had 4 NOCOMPILE drafts, each rewritten into a compiling form, §1.4 P7), run **the single named test** (`go test <pkg> -run '^<test>$' -count=1`) **and then the whole package** (`go test <pkg> -count=1`), restore, assert the restore is byte-identical. A row is KILLED only when both fail. Run with `export AILANG_BIN=$HOME/.pinned-ailang/ailang` on the final tree (`plan-scratch/work`, byte-identical to the round-tripped tree, §1.2 C30); log `plan-scratch/logs/mutants_final.log` (`TOTAL 139, as expected 139, mismatches 0`). **Never commit or snapshot a mutant.** The exact strings are `tools/mutants.json`; in the table below tabs are dropped, a newline prints as ⏎ and a backquote as '. Base controls ran the same harness on the pristine base (`tools/base_controls.json`).

| Id | Stage | Mutant (single edit; ⏎ is a newline) | Named test (arm) | Named | Suite |
|---|---|---|---|---|---|
| MUT-P2 | M3b1 | approval/signed.go: if !ApproverNamePattern.MatchString(a.Name) { // P2 → if false { // P2 | `TestApproveByAssertionOrder` (P2-name) | KILLED | KILLED |
| MUT-P3W | M3b1 | approval/signed.go: if format == FormatWebAuthn && (a.Format != format \|\| a.Passkey == nil) { // P3w → if false { // P3w | `TestApproveByAssertionOrder` (P3w-format-ssh (P3 then says ErrApproverFormat); P3w-no-passkey (panics)) | KILLED | KILLED |
| MUT-P3W-FORMAT | M3b1 | approval/signed.go: (a.Format != format \|\| a.Passkey == nil) { // P3w → (a.Passkey == nil) { // P3w | `TestApproveByAssertionOrder` (P3w-format-ssh) | KILLED | KILLED |
| MUT-P3W-NIL | M3b1 | approval/signed.go: (a.Format != format \|\| a.Passkey == nil) { // P3w → (a.Format != format) { // P3w | `TestApproveByAssertionOrder` (P3w-no-passkey (panics)) | KILLED | KILLED |
| MUT-P3W-SCOPE | M3b1 | approval/signed.go: if format == FormatWebAuthn && (a.Format → if (a.Format | `TestApproveBySignatureOrder` (clean (an ssh approver has no passkey)) | KILLED | KILLED |
| MUT-P3 | M3b1 | approval/signed.go: if a.Format != format { // P3 → if false { // P3 | `TestApproveBySignatureOrder` (format-webauthn) | KILLED | KILLED |
| MUT-P4 | M3b1 | approval/signed.go: if a.RevokedAt != 0 { // P4 → if false { // P4 | `TestApproveByAssertionOrder` (P4-revoked-good-assertion) | KILLED | KILLED |
| MUT-P5 | M3b1 | approval/signed.go: if admit.Now < a.EnrolledAt { // P5 → if false { // P5 | `TestApproveByAssertionOrder` (P5-enrolled-after-now) | KILLED | KILLED |
| MUT-P5-LE | M3b1 | approval/signed.go: if admit.Now < a.EnrolledAt { // P5 → if admit.Now <= a.EnrolledAt { // P5 | `TestApproveByAssertionOrder` (control enrolled-at-now) | KILLED | KILLED |
| MUT-SIG-CHECK | M3b1 | approval/signed.go: if err := checkApprover(in.Admit, in.Approver, FormatSSH); err != nil { // P2-P5 → if err := checkApprover(in.Admit, in.Approver, FormatSSH); false && err != nil { // P2-P5 | `TestApproveBySignatureOrder` (name-upper) | KILLED | KILLED |
| MUT-B7 | M3b1 | approval/signed.go: if s == nil \|\| !s.ok \|\| s.format != FormatWebAuthn { // B7, B8 → if s == nil \|\| s.format != FormatWebAuthn { // B7, B8 | `TestSignedApprovalPasskeyUse` (unverified-webauthn) | KILLED | KILLED |
| MUT-B8 | M3b1 | approval/signed.go: if s == nil \|\| !s.ok \|\| s.format != FormatWebAuthn { // B7, B8 → if s == nil \|\| !s.ok { // B7, B8 | `TestSignedApprovalPasskeyUse` (verified-ssh-with-count) | KILLED | KILLED |
| MUT-B-NIL | M3b1 | approval/signed.go: if s == nil \|\| !s.ok \|\| s.format != FormatWebAuthn { // B7, B8 → if !s.ok \|\| s.format != FormatWebAuthn { // B7, B8 | `TestSignedApprovalPasskeyUse` (nil (panics)) | KILLED | KILLED |
| MUT-PU-CLONE | M3b1 | approval/signed.go: return bytes.Clone(s.credentialID), s.signCount, true → return s.credentialID, s.signCount, true | `TestSignedApprovalPasskeyUse` (the returned ID edited) | KILLED | KILLED |
| MUT-PU-COUNT | M3b1 | approval/signed.go: return bytes.Clone(s.credentialID), s.signCount, true → return bytes.Clone(s.credentialID), 0, true | `TestApproveByAssertionAcceptsChrome` (PasskeyUse counter 2) | KILLED | KILLED |
| MUT-P1 | M3b1 | approval/passkey_assert.go: if v := Admit(in.Admit); v != Admitted { // P1 → if v := Admit(in.Admit); false && v != Admitted { // P1 | `TestApproveByAssertionOrder` (P1-consumed) | KILLED | KILLED |
| MUT-ABA-CHECK | M3b1 | approval/passkey_assert.go: if err := checkApprover(in.Admit, in.Approver, FormatWebAuthn); err != nil { → if err := checkApprover(in.Admit, in.Approver, FormatWebAuthn); false && err != nil { | `TestApproveByAssertionOrder` (P4-revoked-good-assertion) | KILLED | KILLED |
| MUT-W1 | M3b1 | approval/passkey_assert.go: resp, err := ParseAssertionResponse(in.Response) // W1⏎if err != nil { → resp, err := ParseAssertionResponse(in.Response) // W1⏎if false && err != nil { | `TestApproveByAssertionOrder` (W1-envelope) | KILLED | KILLED |
| MUT-W2 | M3b1 | approval/passkey_assert.go: if !bytes.Equal(resp.RawID, p.CredentialID) { // W2 → if false { // W2 | `TestApproveByAssertionOrder` (W2-another-credentials-good-assertion) | KILLED | KILLED |
| MUT-W3 | M3b1 | approval/passkey_assert.go: if resp.HasUserHandle && !bytes.Equal( → if false && resp.HasUserHandle && !bytes.Equal( | `TestApproveByAssertionOrder` (W3-user-handle) | KILLED | KILLED |
| MUT-W3-HAS | M3b1 | approval/passkey_assert.go: if resp.HasUserHandle && !bytes.Equal( → if !bytes.Equal( | `TestApproveByAssertionOrder` (control clean (no userHandle)) | KILLED | KILLED |
| MUT-W3-OR | M3b1 | approval/passkey_assert.go: if resp.HasUserHandle && !bytes.Equal( → if resp.HasUserHandle \|\| !bytes.Equal( | `TestApproveByAssertionOrder` (control clean) | KILLED | KILLED |
| MUT-W4 | M3b1 | approval/passkey_assert.go: cd, err := ParseClientData(resp.ClientDataJSON) // W4⏎if err != nil { → cd, err := ParseClientData(resp.ClientDataJSON) // W4⏎if false && err != nil { | `TestApproveByAssertionOrder` (W4-client-data) | KILLED | KILLED |
| MUT-W5 | M3b1 | approval/passkey_assert.go: if cd.Type != "webauthn.get" { // W5 → if false { // W5 | `TestApproveByAssertionOrder` (W5-type-create) | KILLED | KILLED |
| MUT-W5-CREATE | M3b1 | approval/passkey_assert.go: if cd.Type != "webauthn.get" { // W5 → if cd.Type != "webauthn.create" { // W5 | `TestApproveByAssertionOrder` (control clean) | KILLED | KILLED |
| MUT-W6 | M3b1 | approval/passkey_assert.go: if err != nil \|\| cd.Challenge != b64Challenge(h[:]) { // W6 → if err != nil \|\| false && cd.Challenge != b64Challenge(h[:]) { // W6 | `TestApproveByAssertionOrder` (W6-the-row-challenge-not-the-hash) | KILLED | KILLED |
| MUT-W6-ERR | M3b1 | approval/passkey_assert.go: if err != nil \|\| cd.Challenge != b64Challenge(h[:]) { // W6 → if cd.Challenge != b64Challenge(h[:]) { // W6 | `TestApproveByAssertionOrder` (W6-negative-expiry-has-no-hash) | KILLED | KILLED |
| MUT-W6-EXPIRY | M3b1 | approval/passkey_assert.go: MessageHash(in.Admit.StoredDigest, in.Challenge, in.Admit.ExpiresAt) → MessageHash(in.Admit.StoredDigest, in.Challenge, in.Admit.ExpiresAt+1) | `TestApproveByAssertionOrder` (control clean) | KILLED | KILLED |
| MUT-W7 | M3b1 | approval/passkey_assert.go: if cd.Origin != p.Origin { // W7 → if false { // W7 | `TestApproveByAssertionOrder` (W7-origin) | KILLED | KILLED |
| MUT-W8-CROSS | M3b1 | approval/passkey_assert.go: if cd.CrossOrigin \|\| cd.HasTopOrigin { // W8 → if cd.HasTopOrigin { // W8 | `TestApproveByAssertionOrder` (W8-cross-origin) | KILLED | KILLED |
| MUT-W8-TOP | M3b1 | approval/passkey_assert.go: if cd.CrossOrigin \|\| cd.HasTopOrigin { // W8 → if cd.CrossOrigin { // W8 | `TestApproveByAssertionOrder` (W8-top-origin) | KILLED | KILLED |
| MUT-W9 | M3b1 | approval/passkey_assert.go: ad, err := ParseAuthData(resp.AuthenticatorData, false) // W9⏎if err != nil { → ad, err := ParseAuthData(resp.AuthenticatorData, false) // W9⏎if false && err != nil { | `TestApproveByAssertionOrder` (W9-attested-data) | KILLED | KILLED |
| MUT-W9-REG | M3b1 | approval/passkey_assert.go: ParseAuthData(resp.AuthenticatorData, false) // W9 → ParseAuthData(resp.AuthenticatorData, true) // W9 | `TestApproveByAssertionOrder` (control clean) | KILLED | KILLED |
| MUT-W10 | M3b1 | approval/passkey_assert.go: if ad.RPIDHash != sha256.Sum256([]byte(p.RPID)) { // W10 → if false { // W10 | `TestApproveByAssertionOrder` (W10-rp-id) | KILLED | KILLED |
| MUT-W11 | M3b1 | approval/passkey_assert.go: if !ad.UP { // W11 → if false { // W11 | `TestApproveByAssertionOrder` (W11-no-UP) | KILLED | KILLED |
| MUT-W11-AND-UV | M3b1 | approval/passkey_assert.go: if !ad.UP { // W11 → if !ad.UP && !ad.UV { // W11 | `TestApproveByAssertionOrder` (W11-no-UP (UV set)) | KILLED | KILLED |
| MUT-W11-AND-BE | M3b1 | approval/passkey_assert.go: if !ad.UP { // W11 → if !ad.UP && !ad.BE { // W11 | `TestApproveByAssertionOrder` (W11-no-UP-BE) | KILLED | KILLED |
| MUT-W11-AND-ED | M3b1 | approval/passkey_assert.go: if !ad.UP { // W11 → if !ad.UP && ad.BS == ad.BE { // W11 | `TestApproveByAssertionOrder` (W11-no-UP-BE (BE without BS)) | KILLED | KILLED |
| MUT-W12 | M3b1 | approval/passkey_assert.go: if !ad.UV { // W12 → if false { // W12 | `TestApproveByAssertionOrder` (W12-no-UV) | KILLED | KILLED |
| MUT-W12-AND-UP | M3b1 | approval/passkey_assert.go: if !ad.UV { // W12 → if !ad.UV && !ad.UP { // W12 | `TestApproveByAssertionOrder` (W12-no-UV (UP set)) | KILLED | KILLED |
| MUT-W12-AND-BE | M3b1 | approval/passkey_assert.go: if !ad.UV { // W12 → if !ad.UV && !ad.BE { // W12 | `TestApproveByAssertionOrder` (W12-no-UV-BE) | KILLED | KILLED |
| MUT-W12-AND-BS | M3b1 | approval/passkey_assert.go: if !ad.UV { // W12 → if !ad.UV && !ad.BS { // W12 | `TestApproveByAssertionOrder` (W12-no-UV-BE-BS) | KILLED | KILLED |
| MUT-W12-BE-ONLY | M3b1 | approval/passkey_assert.go: if !ad.UV { // W12 → if !ad.UV && (!ad.BE \|\| p.BackupEligible) { // W12 | `TestApproveByAssertionOrder` (W12-no-UV-BE-not-enrolled) | KILLED | KILLED |
| MUT-W13 | M3b1 | approval/passkey_assert.go: if ad.BE != p.BackupEligible { // W13 → if false { // W13 | `TestApproveByAssertionOrder` (W13-BE-not-enrolled) | KILLED | KILLED |
| MUT-W13-ONEWAY | M3b1 | approval/passkey_assert.go: if ad.BE != p.BackupEligible { // W13 → if ad.BE && !p.BackupEligible { // W13 | `TestApproveByAssertionOrder` (W13-enrolled-BE-absent) | KILLED | KILLED |
| MUT-W13-BS | M3b1 | approval/passkey_assert.go: if ad.BE != p.BackupEligible { // W13 → if ad.BS != p.BackupEligible { // W13 | `TestApproveByAssertionOrder` (W13-BE-not-enrolled (BS clear)) | KILLED | KILLED |
| MUT-W14 | M3b1 | approval/passkey_assert.go: resp.Signature); err != nil { // W14 → resp.Signature); false && err != nil { // W14 | `TestApproveByAssertionOrder` (W14-bad-signature) | KILLED | KILLED |
| MUT-W14-CDJ | M3b1 | approval/passkey_assert.go: p.Key.verify(append(bytes.Clone(resp.AuthenticatorData), cdh[:]...), resp.Signature) → p.Key.verify(append(bytes.Clone(resp.AuthenticatorData), cdh[:0]...), resp.Signature) | `TestApproveByAssertionOrder` (control clean) | KILLED | KILLED |
| MUT-W15 | M3b1 | approval/passkey_assert.go: if (ad.SignCount != 0 \|\| p.SignCount != 0) && ad.SignCount <= p.SignCount { // W15 → if false { // W15 | `TestSignCountLaw` (2/2) | KILLED | KILLED |
| MUT-W15-LT | M3b1 | approval/passkey_assert.go: && ad.SignCount <= p.SignCount { // W15 → && ad.SignCount < p.SignCount { // W15 | `TestSignCountLaw` (2/2) | KILLED | KILLED |
| MUT-W15-BOTHZERO | M3b1 | approval/passkey_assert.go: if (ad.SignCount != 0 \|\| p.SignCount != 0) && ad.SignCount → if ad.SignCount | `TestSignCountLaw` (0/0) | KILLED | KILLED |
| MUT-W15-AND | M3b1 | approval/passkey_assert.go: if (ad.SignCount != 0 \|\| p.SignCount != 0) → if (ad.SignCount != 0 && p.SignCount != 0) | `TestSignCountLaw` (5/0) | KILLED | KILLED |
| MUT-W15-BEFORE-W14 | M3b1 | approval/passkey_assert.go: cdh := sha256.Sum256(resp.ClientDataJSON)⏎if err := p.Key.verify(append(bytes.Clone(resp.AuthenticatorData), cdh[:]...), resp.Signature); err != nil { // W14⏎return nil, err⏎}⏎if (ad.SignCount != 0 \|\| p.SignCount != 0) && ad.SignCount <= p.SignCount { // W15⏎return nil, fmt.Errorf( → if (ad.SignCount != 0 \|\| p.SignCount != 0) && ad.SignCount <= p.SignCount { // W15⏎return nil, fmt.Errorf("%w", ErrSignCountRegression)⏎}⏎cdh := sha256.Sum256(resp.ClientDataJSON)⏎if err := p.Key.verify(append(bytes.Clone(resp.AuthenticatorData), cdh[:]...), resp.Signature); err != nil { // W14⏎return nil, err⏎}⏎if false {⏎return nil, fmt.Errorf( | `TestApproveByAssertionOrder` (W14-before-W15) | KILLED | KILLED |
| MUT-SA-COUNT | M3b1 | approval/passkey_assert.go: signCount: ad.SignCount} → signCount: p.SignCount} | `TestApproveByAssertionAcceptsChrome` (PasskeyUse counter 2) | KILLED | KILLED |
| MUT-SA-CRED-ALIAS | M3b1 | approval/passkey_assert.go: credentialID: bytes.Clone(p.CredentialID) → credentialID: p.CredentialID | `TestSignedApprovalPasskeyUse` (the input credential ID edited) | KILLED | KILLED |
| MUT-SA-KEY | M3b1 | approval/passkey_assert.go: key: p.Key.Canonical() → key: p.Key.X | `TestApproveByAssertionAcceptsChrome` (Use() key is the canonical COSE key) | KILLED | KILLED |
| MUT-SA-FORMAT-EQUIV | M3b1 | approval/passkey_assert.go: format: FormatWebAuthn, key: → format: in.Approver.Format, key: | `TestApproveByAssertionAcceptsChrome` (equivalent: P3w has refused any other format) | SURVIVED | SURVIVED |
| MUT-SA-KEY-ALIAS-EQUIV | M3b1 | approval/passkey_assert.go: key: p.Key.Canonical() → key: append(p.Key.Canonical()[:0:0], p.Key.Canonical()...) | `TestSignedApprovalPasskeyUse` (equivalent: Canonical() always allocates) | SURVIVED | SURVIVED |
| MUT-CR-UPGRADE6 | M4a | store/credentials.go: if version < 6 {⏎statements = append(statements, passkeysSQL) → if false {⏎statements = append(statements, passkeysSQL) | `TestCredentialsUpgradeToV6` (v1, v3, v4, v5) | KILLED | KILLED |
| MUT-CR-TABLES6 | M4a | store/credentials.go: credentialsTable, passkeysTable} → credentialsTable} | `TestCredentialsUpgradeToV6` ((b) the second UpgradeCredentials reopens the version-6 table set) | KILLED | KILLED |
| MUT-CR-TABLES5 | M4a | store/credentials.go: case version == 5:⏎return []string{approversTable, "episode_starts", mintRequestsTable, credentialsTable} → case version == 5:⏎return []string{approversTable, "episode_starts", mintRequestsTable, credentialsTable, passkeysTable} | `TestCredentialsUpgradeToV6` (v5 (refused as not a credentials store)) | KILLED | KILLED |
| MUT-CR-VERSION | M4a | store/credentials.go: const credentialsSchemaVersion = 6 → const credentialsSchemaVersion = 5 | `TestCredentialsSchemaMatchesCanonicalDDL` (user_version 6) | KILLED | KILLED |
| MUT-SQL-DRIFT | M4a | store/credentials_schema.sql: alg             INTEGER NOT NULL CHECK (alg IN (-7, -8)), → alg             INTEGER NOT NULL CHECK (alg IN (-7, -8, -257)), | `TestCredentialsSchemaMatchesCanonicalDDL` (schema file vs canonical) | KILLED | KILLED |
| MUT-PKSQL-DRIFT | M4a | store/passkeys.go: alg             INTEGER NOT NULL CHECK (alg IN (-7, -8)), → alg             INTEGER NOT NULL CHECK (alg IN (-7, -8, -257)), | `TestCredentialsSchemaMatchesCanonicalDDL` (passkeysSQL vs canonical) | KILLED | KILLED |
| MUT-EP-DEADLINE | M4b | store/passkeys.go: func (c *Credentials) EnrolPasskey(ctx context.Context, a ApproverRow, p PasskeyRow) error {⏎if err := requireDeadline(ctx); err != nil { → func (c *Credentials) EnrolPasskey(ctx context.Context, a ApproverRow, p PasskeyRow) error {⏎if err := requireDeadline(ctx); false && err != nil { | `TestCredentialsMethodsRequireADeadline` (EnrolPasskey(Background)) | KILLED | KILLED |
| MUT-EP-OLDER | M4b | store/passkeys.go: if c.version < passkeysSchemaVersion {⏎return c.olderForPasskeys()⏎}⏎if err := checkPasskeyRows → if false {⏎return c.olderForPasskeys()⏎}⏎if err := checkPasskeyRows | `TestEnrolPasskeyRefusals` (v5 store) | KILLED | KILLED |
| MUT-EP-CHECK | M4b | store/passkeys.go: if err := checkPasskeyRows(a, p); err != nil { → if err := checkPasskeyRows(a, p); false && err != nil { | `TestEnrolPasskeyRefusals` (S1-format-ssh) | KILLED | KILLED |
| MUT-S1 | M4b | store/passkeys.go: case a.Format != "webauthn": // S1 → case false: // S1 | `TestEnrolPasskeyRefusals` (S1-format-ssh) | KILLED | KILLED |
| MUT-S2 | M4b | store/passkeys.go: case a.Name != p.Name: // S2 → case false: // S2 | `TestEnrolPasskeyRefusals` (S2-names-differ) | KILLED | KILLED |
| MUT-S3 | M4b | store/passkeys.go: case a.Name == "" \|\| len(a.PublicKey) == 0: // S3 → case false: // S3 | `TestEnrolPasskeyRefusals` (S3-nil-key) | KILLED | KILLED |
| MUT-S4 | M4b | store/passkeys.go: case a.EnrolledAt <= 0: // S4 → case false: // S4 | `TestEnrolPasskeyRefusals` (S4-enrolled-0) | KILLED | KILLED |
| MUT-S5 | M4b | store/passkeys.go: case a.LastUsedAt != 0 \|\| a.RevokedAt != 0: // S5 → case false: // S5 | `TestEnrolPasskeyRefusals` (S5-last-used) | KILLED | KILLED |
| MUT-S6 | M4b | store/passkeys.go: case len(p.CredentialID) < 1 \|\| len(p.CredentialID) > 1023: // S6 → case false: // S6 | `TestEnrolPasskeyRefusals` (S6-cred-1024) | KILLED | KILLED |
| MUT-S7 | M4b | store/passkeys.go: case p.Alg != -7 && p.Alg != -8: // S7 → case false: // S7 | `TestEnrolPasskeyRefusals` (S7-alg-rs256) | KILLED | KILLED |
| MUT-S8 | M4b | store/passkeys.go: case len(p.AAGUID) != 16: // S8 → case false: // S8 | `TestEnrolPasskeyRefusals` (S8-aaguid-15) | KILLED | KILLED |
| MUT-S9 | M4b | store/passkeys.go: case p.RPID == "" \|\| p.Origin == "": // S9 → case false: // S9 | `TestEnrolPasskeyRefusals` (S9-rp-empty) | KILLED | KILLED |
| MUT-S10 | M4b | store/passkeys.go: case p.SignCount < 0 \|\| p.SignCount > maxSignCount: // S10 → case false: // S10 | `TestEnrolPasskeyRefusals` (S10-count--1) | KILLED | KILLED |
| MUT-S3-NAME | M4b | store/passkeys.go: case a.Name == "" \|\| len(a.PublicKey) == 0: → case len(a.PublicKey) == 0: | `TestEnrolPasskeyRefusals` (S3-empty-name) | KILLED | KILLED |
| MUT-S3-KEY | M4b | store/passkeys.go: case a.Name == "" \|\| len(a.PublicKey) == 0: → case a.Name == "": | `TestEnrolPasskeyRefusals` (S3-nil-key) | KILLED | KILLED |
| MUT-S4-LT | M4b | store/passkeys.go: case a.EnrolledAt <= 0: → case a.EnrolledAt < 0: | `TestEnrolPasskeyRefusals` (S4-enrolled-0) | KILLED | KILLED |
| MUT-S5-USED | M4b | store/passkeys.go: case a.LastUsedAt != 0 \|\| a.RevokedAt != 0: → case a.RevokedAt != 0: | `TestEnrolPasskeyRefusals` (S5-last-used) | KILLED | KILLED |
| MUT-S5-REVOKED | M4b | store/passkeys.go: case a.LastUsedAt != 0 \|\| a.RevokedAt != 0: → case a.LastUsedAt != 0: | `TestEnrolPasskeyRefusals` (S5-revoked) | KILLED | KILLED |
| MUT-S6-FLOOR | M4b | store/passkeys.go: case len(p.CredentialID) < 1 \|\| len(p.CredentialID) > 1023: → case len(p.CredentialID) > 1023: | `TestEnrolPasskeyRefusals` (S6-cred-0) | KILLED | KILLED |
| MUT-S6-FLOOR16 | M4b | store/passkeys.go: case len(p.CredentialID) < 1 \|\| len(p.CredentialID) > 1023: → case len(p.CredentialID) < 16 \|\| len(p.CredentialID) > 1023: | `TestEnrolPasskeyRefusals` (edge: a 1-byte credential enrols) | KILLED | KILLED |
| MUT-S6-CEIL | M4b | store/passkeys.go: case len(p.CredentialID) < 1 \|\| len(p.CredentialID) > 1023: → case len(p.CredentialID) < 1 \|\| len(p.CredentialID) >= 1023: | `TestEnrolPasskeyRefusals` (edge: a 1,023-byte credential enrols) | KILLED | KILLED |
| MUT-S7-ES | M4b | store/passkeys.go: case p.Alg != -7 && p.Alg != -8: → case p.Alg != -7: | `TestEnrolPasskeyRefusals` (edge: alg -8 enrols) | KILLED | KILLED |
| MUT-S7-ED | M4b | store/passkeys.go: case p.Alg != -7 && p.Alg != -8: → case p.Alg != -7 && p.Alg != -257: | `TestEnrolPasskeyRefusals` (edge: alg -8 enrols) | KILLED | KILLED |
| MUT-S9-RP | M4b | store/passkeys.go: case p.RPID == "" \|\| p.Origin == "": → case p.Origin == "": | `TestEnrolPasskeyRefusals` (S9-rp-empty) | KILLED | KILLED |
| MUT-S9-ORIGIN | M4b | store/passkeys.go: case p.RPID == "" \|\| p.Origin == "": → case p.RPID == "": | `TestEnrolPasskeyRefusals` (S9-origin-empty) | KILLED | KILLED |
| MUT-S10-NEG | M4b | store/passkeys.go: case p.SignCount < 0 \|\| p.SignCount > maxSignCount: → case p.SignCount > maxSignCount: | `TestEnrolPasskeyRefusals` (S10-count--1) | KILLED | KILLED |
| MUT-S10-MAX | M4b | store/passkeys.go: case p.SignCount < 0 \|\| p.SignCount > maxSignCount: → case p.SignCount < 0: | `TestEnrolPasskeyRefusals` (S10-count-2^63-1 (the DDL CHECK refuses 2^32 with a driver error)) | KILLED | KILLED |
| MUT-S10-GE | M4b | store/passkeys.go: case p.SignCount < 0 \|\| p.SignCount > maxSignCount: → case p.SignCount < 0 \|\| p.SignCount >= maxSignCount: | `TestEnrolPasskeyRefusals` (edge: counter 2^32-1 enrols) | KILLED | KILLED |
| MUT-PRE-NAME | M4b | store/passkeys.go: {'SELECT name FROM approvers WHERE name = ?', a.Name, ErrApproverNameTaken},⏎ → (deleted) | `TestEnrolPasskeyRefusals` (name-taken (the PRIMARY KEY backstop is not ErrApproverNameTaken)) | KILLED | KILLED |
| MUT-PRE-KEY | M4b | store/passkeys.go: {'SELECT name FROM approvers WHERE public_key = ?', a.PublicKey, ErrApproverKeyEnrolled},⏎ → (deleted) | `TestEnrolPasskeyRefusals` (key-enrolled) | KILLED | KILLED |
| MUT-PRE-CRED | M4b | store/passkeys.go: {'SELECT name FROM webauthn_credentials WHERE credential_id = ?', p.CredentialID, ErrPasskeyCredentialEnrolled},⏎ → (deleted) | `TestEnrolPasskeyRefusals` (credential-enrolled) | KILLED | KILLED |
| MUT-PRE-HOLDER | M4b | store/passkeys.go: return fmt.Errorf("%w: approver %q holds it", pre.err, holder) → return fmt.Errorf("%w: approver %q holds it", pre.err, "") | `TestEnrolPasskeyRefusals` (key-enrolled (names the holder)) | KILLED | KILLED |
| MUT-INS-APPROVER | M4b | store/passkeys.go: INSERT INTO approvers (name, format, public_key, enrolled_at) VALUES (?, ?, ?, ?) → SELECT ?, ?, ?, ? | `TestEnrolPasskeyRefusals` (good enrol (6 approvers, 5 passkeys)) | KILLED | KILLED |
| MUT-INS-PASSKEY | M4b | store/passkeys.go: INSERT INTO webauthn_credentials (name, credential_id, rp_id, origin, alg, aaguid, backup_eligible, sign_count) VALUES (?, ?, ?, ?, ?, ?, ?, ?) → SELECT ?, ?, ?, ?, ?, ?, ?, ? | `TestEnrolPasskeyRefusals` (good enrol (6 approvers, 5 passkeys)) | KILLED | KILLED |
| MUT-LK-DEADLINE | M4c | store/passkeys.go: func (c *Credentials) enrolledPasskeys(ctx context.Context, where string, args ...any) ([]EnrolledPasskey, error) {⏎if err := requireDeadline(ctx); err != nil { → func (c *Credentials) enrolledPasskeys(ctx context.Context, where string, args ...any) ([]EnrolledPasskey, error) {⏎if err := requireDeadline(ctx); false && err != nil { | `TestCredentialsMethodsRequireADeadline` (PasskeyByCredentialID, Passkey, ActivePasskeys (Background)) | KILLED | KILLED |
| MUT-LK-OLDER | M4c | store/passkeys.go: if c.version < passkeysSchemaVersion {⏎return nil, c.olderForPasskeys() → if false {⏎return nil, c.olderForPasskeys() | `TestPasskeyLookups` (v5 lookups) | KILLED | KILLED |
| MUT-LK-JOIN | M4c | store/passkeys.go: JOIN webauthn_credentials w ON w.name = a.name → JOIN webauthn_credentials w ON 1 | `TestPasskeyLookups` (ActivePasskeys pairs) | KILLED | KILLED |
| MUT-LK-BYID-MISS | M4c | store/passkeys.go: if err != nil \|\| len(got) == 0 {⏎return ApproverRow{}, PasskeyRow{}, false, err → if err != nil {⏎return ApproverRow{}, PasskeyRow{}, false, err | `TestPasskeyLookups` (PasskeyByCredentialID(prefix) (panics)) | KILLED | KILLED |
| MUT-LK-NAME-MISS | M4c | store/passkeys.go: if err != nil \|\| len(got) == 0 {⏎return PasskeyRow{}, false, err → if err != nil {⏎return PasskeyRow{}, false, err | `TestPasskeyLookups` (Passkey(s) (panics)) | KILLED | KILLED |
| MUT-LK-ACTIVE-REVOKED | M4c | store/passkeys.go: 'a.revoked_at IS NULL AND a.format → 'a.format | `TestPasskeyLookups` (b revoked) | KILLED | KILLED |
| MUT-LK-ACTIVE-FORMAT | M4c | store/passkeys.go: a.revoked_at IS NULL AND a.format = 'webauthn' ORDER BY → a.revoked_at IS NULL ORDER BY | `TestPasskeyLookups` (an ssh row with a hand-made passkey row) | KILLED | KILLED |
| MUT-LK-ACTIVE-ORDER | M4c | store/passkeys.go: ORDER BY a.name LIMIT ? → ORDER BY a.name DESC LIMIT ? | `TestPasskeyLookups` (a,c) | KILLED | KILLED |
| MUT-LK-ACTIVE-CAP | M4c | store/passkeys.go: maxActivePasskeys = 64 → maxActivePasskeys = 65 | `TestPasskeyLookups` (64 of 72) | KILLED | KILLED |
| MUT-LK-REVOKED-READ | M4c | store/passkeys.go: a.LastUsedAt, a.RevokedAt = lastUsed.Int64, revoked.Int64 → a.LastUsedAt, a.RevokedAt = lastUsed.Int64, 0 | `TestPasskeyLookups` (b's revoked_at) | KILLED | KILLED |
| MUT-LK-USED-READ | M4c | store/passkeys.go: a.LastUsedAt, a.RevokedAt = lastUsed.Int64, revoked.Int64 → a.LastUsedAt, a.RevokedAt = 0, revoked.Int64 | `TestPasskeyLookups` (a's last_used_at) | KILLED | KILLED |
| MUT-LK-COLUMNS | M4c | store/passkeys.go: w.rp_id, w.origin, w.alg → w.origin, w.rp_id, w.alg | `TestPasskeyRowMatchesTheTable` (rp_id and origin swapped) | KILLED | KILLED |
| MUT-LK-WHERE-ID | M4c | store/passkeys.go: 'w.credential_id = ?', id) → 'w.credential_id >= ?', id) | `TestPasskeyLookups` (a prefix misses) | KILLED | KILLED |
| MUT-C13 | M5a | store/mint_requests.go: if a := rc.Approver; a != nil && (a.Format == "webauthn") != (a.Passkey != nil) { // C13 → if false { // C13 | `TestPasskeyConsumeCompleteness` (C13-webauthn-without-passkey) | KILLED | KILLED |
| MUT-C13-NIL | M5a | store/mint_requests.go: if a := rc.Approver; a != nil && (a.Format == "webauthn") != (a.Passkey != nil) { // C13 → if a := rc.Approver; (a.Format == "webauthn") != (a.Passkey != nil) { // C13 | `TestTypedConsumeTouchesNoApprover` (a typed consume (panics)) | KILLED | KILLED |
| MUT-C13-ONEWAY | M5a | store/mint_requests.go: a != nil && (a.Format == "webauthn") != (a.Passkey != nil) { // C13 → a != nil && a.Format == "webauthn" && a.Passkey == nil { // C13 | `TestPasskeyConsumeCompleteness` (C13-ssh-with-passkey) | KILLED | KILLED |
| MUT-C14 | M5a | store/mint_requests.go: p.SignCount > maxSignCount { // C14 → p.SignCount > maxSignCount && false { // C14 | `TestPasskeyConsumeCompleteness` (C14-empty-credential) | KILLED | KILLED |
| MUT-C14-FLOOR | M5a | store/mint_requests.go: len(p.CredentialID) < 1 \|\| len(p.CredentialID) > 1023 \|\| → len(p.CredentialID) > 1023 \|\| | `TestPasskeyConsumeCompleteness` (C14-empty-credential) | KILLED | KILLED |
| MUT-C14-CEIL | M5a | store/mint_requests.go: len(p.CredentialID) < 1 \|\| len(p.CredentialID) > 1023 \|\| → len(p.CredentialID) < 1 \|\| | `TestPasskeyConsumeCompleteness` (C14-1024-byte-credential) | KILLED | KILLED |
| MUT-C14-CEIL-GE | M5a | store/mint_requests.go: len(p.CredentialID) > 1023 \|\| → len(p.CredentialID) >= 1023 \|\| | `TestPasskeyConsumeAdvancesTheCounter` (a 1,023-byte credential) | KILLED | KILLED |
| MUT-C14-NEG | M5a | store/mint_requests.go: \|\| p.SignCount < 0 \|\| p.SignCount > maxSignCount { // C14 → \|\| p.SignCount > maxSignCount { // C14 | `TestPasskeyConsumeCompleteness` (C14-counter--1) | KILLED | KILLED |
| MUT-C14-MAX | M5a | store/mint_requests.go: \|\| p.SignCount < 0 \|\| p.SignCount > maxSignCount { // C14 → \|\| p.SignCount < 0 { // C14 | `TestPasskeyConsumeCompleteness` (C14-counter-2^32) | KILLED | KILLED |
| MUT-C14-MAX-GE | M5a | store/mint_requests.go: p.SignCount > maxSignCount { // C14 → p.SignCount >= maxSignCount { // C14 | `TestPasskeyConsumeAdvancesTheCounter` (stored 1, new 2^32-1) | KILLED | KILLED |
| MUT-C15 | M5a | store/mint_requests.go: if c.version < passkeysSchemaVersion { // C15 → if false { // C15 | `TestPasskeyConsumeCompleteness` (C15-v5-store) | KILLED | KILLED |
| MUT-C13-BEFORE-C4 | M5a | store/mint_requests.go: if rc.Approver != nil && rc.ApprovedBy != rc.Approver.Format+":"+rc.Approver.Name { // C4 → if a := rc.Approver; a != nil && (a.Format == "webauthn") != (a.Passkey != nil) { // C13 moved before C4⏎return ErrIncompleteConsume⏎}⏎if rc.Approver != nil && rc.ApprovedBy != rc.Approver.Format+":"+rc.Approver.Name { // C4 | `TestApprovedByMatchesTheApprover` (approver-format-not-ssh (unedited, design CS7 ordering mutant)) | KILLED | KILLED |
| MUT-USE-PASSKEY | M5a | store/mint_requests.go: if p := rc.Approver.Passkey; p != nil { // C16-C18 → if p := rc.Approver.Passkey; false && p != nil { // C16-C18 | `TestPasskeyConsumeAdvancesTheCounter` (2/2 refused) | KILLED | KILLED |
| MUT-UP-CRED | M5a | store/passkeys.go: WHERE name = ? AND credential_id = ? AND → WHERE name = ? AND (credential_id = ? OR 1) AND | `TestPasskeyConsumeChecksTheRowInTheTransaction` (credential-rewritten) | KILLED | KILLED |
| MUT-UP-LE | M5a | store/passkeys.go: AND (sign_count < ? OR → AND (sign_count <= ? OR | `TestPasskeyConsumeAdvancesTheCounter` (2/2) | KILLED | KILLED |
| MUT-UP-BOTHZERO | M5a | store/passkeys.go: OR (sign_count = 0 AND ? = 0)) → OR (0 AND ? = 0)) | `TestPasskeyConsumeAdvancesTheCounter` (0/0) | KILLED | KILLED |
| MUT-UP-SET | M5a | store/passkeys.go: SET sign_count = ? WHERE → SET sign_count = ? * 0 + sign_count WHERE | `TestPasskeyConsumeAdvancesTheCounter` (1 -> 2 (the counter moves)) | KILLED | KILLED |
| MUT-UP-ROWS | M5a | store/passkeys.go: err != nil \|\| n == 1 { → err != nil \|\| n >= 0 { | `TestPasskeyConsumeAdvancesTheCounter` (2/2 refused) | KILLED | KILLED |
| MUT-C16 | M5a | store/passkeys.go: case errors.Is(err, sql.ErrNoRows): // C16 → case false: // C16 | `TestPasskeyConsumeChecksTheRowInTheTransaction` (passkey-deleted) | KILLED | KILLED |
| MUT-C17 | M5a | store/passkeys.go: case !bytes.Equal(id, p.CredentialID): // C17 → case !bytes.Equal(id, p.CredentialID) && false: // C17 | `TestPasskeyConsumeChecksTheRowInTheTransaction` (credential-rewritten) | KILLED | KILLED |
| MUT-C18 | M5a | store/passkeys.go: ErrSignCountRegression, name, count, p.SignCount, name) → ErrPasskeyChanged, name, count, p.SignCount, name) | `TestPasskeyConsumeChecksTheRowInTheTransaction` (counter-raised) | KILLED | KILLED |
| MUT-A6 | M5b | authority/mint.go: if id, count, ok := rm.Signed.PasskeyUse(); ok { // A6 → if id, count, ok := rm.Signed.PasskeyUse(); false && ok { // A6 | `TestMintForRequestDerivesPasskeyUse` (the passkey mint (C13 then refuses)) | KILLED | KILLED |
| MUT-A6-COUNT | M5b | authority/mint.go: SignCount: int64(count)} → SignCount: 1 + 0*int64(count)} | `TestMintForRequestDerivesPasskeyUse` (counter 5) | KILLED | KILLED |
| MUT-A6-CRED | M5b | authority/mint.go: store.PasskeyUse{CredentialID: id, → store.PasskeyUse{CredentialID: append(id, 0), | `TestMintForRequestDerivesPasskeyUse` (the credential that signed (C17)) | KILLED | KILLED |
| MUT-SOURCE-NAMES | M3b1 | approval/testdata/webauthn/SOURCE.txt:   chrome-es256.json  from the probe's es256.json, sha256 →   es256.json  sha256 | `TestPasskeyFixturesAreTheGeneratedBytes` (SOURCE.txt names chrome-es256.json) | KILLED | KILLED |
| MUT-PS5-BASE | M4b | store/approvers.go: if row.Format != "ssh" { → if false { | `TestEnrolApproverStillSSHOnly` (EnrolApprover(webauthn) (approvers.go:89, pre-existing)) | KILLED | KILLED |
| MUT-AD-ALIAS-EQUIV | M3b1 | approval/authdata.go: a.CredentialID = append([]byte(nil), rest[18:18+n]...) → a.CredentialID = rest[18:18+n] | `TestNewPasskeyApproverDoesNotAliasItsInput` (equivalent: rest is base64-decoded from the response, never the caller's bytes) | SURVIVED | SURVIVED |
| MUT-COSE-ALIAS-EQUIV | M3b1 | approval/cose.go: *c.dst = append([]byte(nil), b.b...) → *c.dst = b.b | `TestNewPasskeyApproverDoesNotAliasItsInput` (equivalent: b.b is a slice of the decoded attestationObject, never the caller's bytes) | SURVIVED | SURVIVED |
| BC-SOURCE | base `3b77c294` | approval/testdata/webauthn/SOURCE.txt:   es256.json  sha256 →   zz-not-a-fixture.json  sha256 | `TestPasskeyFixturesAreTheGeneratedBytes` (SOURCE.txt's file names are not pinned at the base (judge carry-forward)) | SURVIVED | SURVIVED |
| BC-CS2-ORDER | base `3b77c294` | store/mint_requests.go: if rc.Approver != nil && rc.ApprovedBy != rc.Approver.Format+":"+rc.Approver.Name { // C4 → if false { // C4 | `TestApprovedByMatchesTheApprover` (the unedited approver-format-not-ssh arm is C4's (it is what MUT-C13-BEFORE-C4 relies on)) | KILLED | KILLED |

**Named equivalents (SURVIVED, as predicted):**
- **MUT-SA-FORMAT-EQUIV** (M3b1): equivalent because P3w has refused any other format.
- **MUT-SA-KEY-ALIAS-EQUIV** (M3b1): equivalent because Canonical() always allocates.
- **MUT-AD-ALIAS-EQUIV** (M3b1): equivalent because rest is base64-decoded from the response, never the caller's bytes.
- **MUT-COSE-ALIAS-EQUIV** (M3b1): equivalent because b.b is a slice of the decoded attestationObject, never the caller's bytes.

**Propagation-only (named, not mutated):** the `err` returns that hand a callee's error back unchanged (`ParseAssertionResponse`, `ParseClientData`, `ParseAuthData`, `COSEKey.verify` inside the assertion law are mutated through their `if err != nil` guards above; their own branches are sprint 1's); `usePasskey`'s and `enrolledPasskeys`' driver-error branches (`ExecContext`, `RowsAffected`, `QueryContext`, `rows.Err` failing), which no test can reach without a fault-injecting driver, as in A1b; `EnrolPasskey`'s `BeginTx`/`Commit` failures. Test support (`approvaltest/passkey.go`) is not production and is not mutated: XW1 pins its layout to Chrome's and TP1 to the production laws.

**Design §6's mutant list (copied), and where each lands:** K1–K16, Q1–Q8, D1–D8, J1–J7, O1–O3, E2–E4, R1–R15 — sprint 1 (landed). **P3w, W1–W15, B7–B8 — this sprint (M3b1).** **S1–S9 (S1–S10 here, §1.4 P3), the three `EnrolPasskey` pre-checks — M4b. C13–C18 — M5a. A6 — M5b.** WN1–WN11, WF1–WF6, the widened F1, L4–L5 — sprint 3.

---

## 5. Files (final `git status --short` after M5b)

| File | Sub-milestone | Change |
|---|---|---|
| `host/approval/signed.go`, `signed_test.go` | M3b1 | `credentialID`, `signCount`, `checkApprover`, `PasskeyUse`; RC1 |
| `host/approval/passkey_assert.go`, `passkey_assert_test.go` | M3b1 | new: the assertion law; PA2, PA3 (Chrome), PA4, the aliasing pin |
| `host/approval/testdata/webauthn/SOURCE.txt`, `webauthn_fixture_test.go` | M3b1 | the carry-forward (S2-D10) |
| `host/approval/approvaltest/passkey.go`, `host/approval/passkey_testauth_test.go` | M3b2 | new: the test authenticator; TP1, XW1 |
| `host/approval/passkey_law_test.go` | M3b3 | new: PA1, PA3 |
| `host/store/credentials_schema.sql`, `credentials.go`, `passkeys.go` (M4a section) | M4a | schema 6 |
| `host/store/credentials_test.go`, `approvers_test.go`, `mint_requests_test.go` | M4a (+M4b, M4c: deadline rows) | PS1, PS4, the table list |
| `cmd/ailang-worldd/approver_rig_test.go`, `approve_sig_test.go` | M4a | the v4 fixture builders (§1.4 P2) |
| `host/store/passkeys.go` (M4b section), `passkeys_test.go` | M4b | new test file: PS2, PS5 |
| `host/store/passkeys.go` (M4c section), `passkey_lookups_test.go` | M4c | new test file: PS3, RC6 |
| `host/store/mint_requests.go`, `passkeys.go` (M5a section), `consume_binding_test.go`, `passkey_consume_test.go` | M5a | C13–C18; CS7–CS9, CS3's edit, RC3 |
| `host/authority/mint.go`, `mint_passkey_test.go` | M5b | A6; AU3 |
| this plan and `sprint_w-world-approvals-by-key-a2-s2.json` | planner | |

---

## 6. Sprint gates (controller, after the last sub-milestone; each green at the base, §1.1)

In **one** `bash` shell, from the worktree, **not under `nohup`** (§1.1). Every leg `tee`s its full
output to a log; a leg's verdict is read from its log, never from `tail` alone:

```
bash <<'EOF'
set -euo pipefail
export AILANG_BIN=$HOME/.pinned-ailang/ailang
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter380-row201-a2s2; cd "$W"
L=$HOME/.ailang/state/world-iter380/gates; mkdir -p "$L"
./scripts/verify_ail.sh 2>&1 | tee "$L/verify_ail.log"                      # no .ail change: exactly green
out=$(gofmt -l cmd host); echo "$out" | tee "$L/gofmt.log"; test -z "$out"
go build ./... 2>&1 | tee "$L/build.log"
go vet ./... 2>&1 | tee "$L/vet.log"
go test -run '^$' ./... 2>&1 | tee "$L/fence.log"                           # compile fence for every _test.go
go test -count=1 ./host/approval/... ./host/store/ ./host/authority/ ./cmd/ailang-worldd/ ./host/daemon/ 2>&1 | tee "$L/pkgs.log"
go test -race -count=1 ./host/approval/... ./host/store/ ./host/authority/ 2>&1 | tee "$L/race.log"
# the repo-wide inventories (design §7) and the reflective tests, by name, -v:
go test -v -count=1 ./host/broker/ -run '^TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential$' 2>&1 | tee "$L/inv.log"
go test -v -count=1 ./host/daemon/ -run '^(TestDaemonDependencyAllowlist|TestFrozenV1RouteTableMatchesMux|TestSurfaceAppRoutes|TestDaemonConstructorCallerInventory)$' 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./host/boundary/ -run '^(TestWorldBoundaryDependencyAllowlist|TestBareNetHTTPExemptionIsPerGroup|TestEveryCommandDirectoryIsAProtectedGroup)$' 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./host/store/ -run '^(TestProductionContextRoots|TestProductionGoSurface|TestNoNewDeadlineFreeStoreReads|TestReadGettersHonorContext|TestCredentialsMethodsRequireADeadline)$' 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./host/verifygate/ -run '^TestNoParallelWriteForkPackagesOutsideVerifygate$' 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./host/approval/... -run '^(TestApprovalImportsOnlyTheStandardLibrary|TestApprovalTestIsTestOnly)$' 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./host/approval/ ./host/store/ ./host/authority/ -run '^(TestSignedApprovalFieldsAreAllBound|TestEnrolmentPromptShowsEveryPersistedField|TestRequestConsumeCompleteness|TestApproverRowMatchesTheTable|TestPasskeyRowMatchesTheTable|TestRequestMintFieldsAreSourced)$' 2>&1 | tee -a "$L/inv.log"
for t in TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential TestDaemonDependencyAllowlist TestFrozenV1RouteTableMatchesMux \
  TestSurfaceAppRoutes TestDaemonConstructorCallerInventory TestWorldBoundaryDependencyAllowlist TestBareNetHTTPExemptionIsPerGroup \
  TestEveryCommandDirectoryIsAProtectedGroup TestProductionContextRoots TestProductionGoSurface TestNoNewDeadlineFreeStoreReads \
  TestReadGettersHonorContext TestCredentialsMethodsRequireADeadline TestNoParallelWriteForkPackagesOutsideVerifygate \
  TestApprovalImportsOnlyTheStandardLibrary TestApprovalTestIsTestOnly TestSignedApprovalFieldsAreAllBound \
  TestEnrolmentPromptShowsEveryPersistedField TestRequestConsumeCompleteness TestApproverRowMatchesTheTable \
  TestPasskeyRowMatchesTheTable TestRequestMintFieldsAreSourced; do
  grep -qE -- "^--- PASS: $t " "$L/inv.log" || { echo "$t did not PASS"; exit 1; }; done
if grep -E '^\s*--- (SKIP|FAIL)' "$L/inv.log"; then exit 1; fi
echo PRE_COMMIT_GATES_GREEN | tee "$L/pre_commit.ok"
EOF
# AFTER the controller commits (the tracked-binary gate reads committed blobs; this leg runs go test ./...):
bash -c 'set -o pipefail; export AILANG_BIN=$HOME/.pinned-ailang/ailang; cd /Users/voightkampff/dev/sunholo-data/.wt-world-iter380-row201-a2s2 && VERIFY_GO_LEGS=plain ./scripts/verify_go.sh 2>&1 | tee $HOME/.ailang/state/world-iter380/gates/verify_go.log; echo "verify_go rc=$?" | tee -a $HOME/.ailang/state/world-iter380/gates/verify_go.log'
git -C /Users/voightkampff/dev/sunholo-data/.wt-world-iter380-row201-a2s2 status --short   # only the two plan files, until the controller commits them
```

**At the base**: every leg above is green at `3b77c294` (§1.1) except the after-commit
`verify_go.sh`, which is red there on `host/broker` alone under whole-suite load (stray loopback
`GET /`; all four tests pass alone; CI green). `TestPasskeyRowMatchesTheTable` is new (RC6) and
absent at the base by design; the other 21 named tests PASS there.
**Known load- or environment-sensitive (re-run once, alone; do not chase)**: `host/broker` (the four
§1.1 tests and the local `-race` leg), `TestSeedRestartIgnoresASurvivingOldCopy/survivor` (row 259),
the local `-race` of `host/daemon` and `host/verifygate`,
`TestMissionConfigDefaultFlowIntegration/invalid_ledger`, and anything run under `nohup`
(`scripts/marketingcheck` shell-hup). None is in a package this sprint touches. Any other red, any
race report, or any red in a touched package is a regression. CI's `go-race-verify` is the race
verdict.

**Must stay green, unedited**: every a0a1, A1b and sprint-1 test, except the named edits: RC1
(`signed_test.go`), XW0's substring list (`webauthn_fixture_test.go`), PS1 and the deadline rows
(`credentials_test.go`), `TestCredentialsUpgradeToV5` → `…ToV6` (`approvers_test.go`), the table list
in `mint_requests_test.go:132`, CS3 `other-format` and RC3 (`consume_binding_test.go`), and the two `cmd`
fixture builders. In particular `TestApproveBySignatureOrder`, `TestSignedApprovalBinds`,
`TestSignedApprovalAccessorsFailClosed`, every `TestNewPasskeyApprover*`, `TestApprovedByMatchesTheApprover`
(its `approver-format-not-ssh` arm unedited), `TestSignedConsumeChecksTheApproverInTheTransaction`,
`TestTypedConsumeTouchesNoApprover`, `TestEnrolApproverRefusals`, every `TestApprover*`/`TestApprove*` of
`cmd/ailang-worldd`, and every §7 inventory.

---

## 7. Executor directive notes

- **No git writes.** Never `git add`, `commit`, `stash`, `checkout`, `reset`, `restore`, `rebase`, `push`,
  `worktree` or `apply`. `patch -p1` (§2) is the only writer; read-only `git status`/`git diff` are fine.
- **Runs.** Run 1 = M3b1, M3b2; Run 2 = M3b3, M4a; Run 3 = M4b, M4c; Run 4 = M5a, M5b. Each run reads
  §0, §2 and its own §3 sections, then applies its patches in order.
- **Expected numstat after each sub-milestone, measured AFTER the controller's commit of it** (an
  uncommitted tree hides the new untracked files from numstat):
  ```
  bash -c 'cd /Users/voightkampff/dev/sunholo-data/.wt-world-iter380-row201-a2s2 && git diff --numstat 3b77c294000ca40c2fe77e312ae0e8ce0f71a030 HEAD -- cmd host | awk "{s+=\$1+\$2} END {print s}"'
  ```
  After the controller's commit of M3b1 **345**, M3b2 **672**, M3b3 **890**, M4a **1027**, M4b **1300**,
  M4c **1538**, M5a **1777**, M5b **1843**. **Any other value means a patch was not applied as written:
  stop and report.**
- **Snapshot after every sub-milestone, cumulatively**:
  ```
  bash -c 'set -euo pipefail; M=M3b1; W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter380-row201-a2s2
    D=$HOME/.ailang/state/world-iter380/snap/$M; rm -rf "$D"; mkdir -p "$D"; cd "$W"
    git status --porcelain --untracked-files=all | while read -r st f; do
      case "$f" in design_docs/planned/w-world-approvals-by-key-a2-sprint2-plan.md|design_docs/planned/sprint_w-world-approvals-by-key-a2-s2.json) continue;; esac
      mkdir -p "$D/$(dirname "$f")"; cp -p "$f" "$D/$f"; done
    (cd "$D" && find . -type f ! -name MANIFEST.sha256 | sort | xargs shasum -a 256) > "$D/MANIFEST.sha256"
    cat "$D/MANIFEST.sha256"'
  ```
  (`M` = `M3b1` … `M5b`.) No file is deleted, so the porcelain status has only `M` and `??`; stop on any
  other status. The controller rebuilds one commit per sub-milestone (iteration 377's `commit_ms.sh`
  with `WT` and the snapshot root pointed here).
- **Stop conditions.** A `patch --dry-run` failure; a numstat other than the expected value; a red
  earlier-milestone test; a `--- SKIP` or `--- FAIL` in any `-v` run; any need to touch a file outside §5;
  running Chrome, Node, `ssh-keygen` or any browser; any reference to `~/.ssh`, `~/.ailang/world` or a
  real `/dev/tty`. Stop and report; do not work around.
- **Timing.** Every new `host/approval` test runs in under 1 s; `host/store`'s whole package ~12 s;
  M4a's whole `cmd/ailang-worldd` ~2 min; the compile fence ~25 s.
- **Shell.** Wrap loops, globs and `$VAR` expansions in `bash -c` or a `bash <<'EOF'` heredoc (zsh does
  not word-split). Export `AILANG_BIN` in the same shell as the Go tests. `tee` every gate's output to a
  file; never pipe a gate's only copy through `tail`. **Never start a gate under `nohup`** (§1.1). Poll
  any background command's artifact in a bounded `date +%s` loop in the same turn; never end a turn
  announcing a wait.

---

## 8. Risks

- **R1 — Size.** 1,843 of 2,000, measured; the per-step expected numstat makes drift visible at once.
- **R2 — A schema-5 store asks for one `serve` restart** (design §2.10 "The cost"): bound, started and
  resumed mints, `approver *` and the signature approval refuse a v5 store from this build until
  `serve` upgrades it. This is the design's accepted cost; the CLI's refusal text names the restart
  (`OlderCredentialsStoreError`), pinned by PS2/PS3/CS7 for `passkeys`.
- **R3 — No operator path reaches the new branches yet** (design §0): `EnrolPasskey`, the lookups and
  the passkey consume are called only by tests until sprint 3's CLI. That is the design's sprint
  boundary, and every branch is mutant-pinned now.
- **R4 — The test authenticator is a model of Chrome.** XW1 pins it to the two recordings' layout;
  sprint 3's XB1 (a live Chrome) is the end-to-end check, and the attended exit run the real one.
- **R5 — Haiku and patches.** The executor transcribes nothing; a skipped patch fails `--dry-run` or the
  numstat; a hand "fix" after a red test is forbidden by §7.

---

## 9. Acceptance criteria (the judge scores against these)

1. **M3b — the assertion law.** `ApproveByAssertion` refuses each of P1–P5, P3w and W1–W15 with its
   typed error and a nil approval, in the design's order (the seven order pairs of PA1), for ES256 and
   EdDSA; W11–W13 are refused in every other flag state that reaches them (PA1's 16 flag-state arms);
   the counter law is "both zero passes, otherwise strictly greater", with the refusal naming `approver
   revoke` (PA3); Chrome's two recorded assertions approve the fixture request and the recorded
   `assertion` is refused after `assertion2`'s counter (PA2, XW2).
2. **`checkApprover` is shared**: `ApproveBySignature`'s behaviour is unchanged (`TestApproveBySignatureOrder`
   unedited and green) and its P2–P5 mutants die by both laws' tests.
3. **`SignedApproval`**: `PasskeyUse()` is ok only for a verified webauthn approval (B7, B8), returns a
   copy, and is unaffected by later edits of the input approver; RC1 classifies and perturbs both new
   fields in an ssh and a passkey approval (PA4, RC1).
4. **The test authenticator** makes registrations and assertions the production laws accept for both
   algorithms (TP1), in Chrome's measured layout (XW1); it is test support only (`TestApprovalTestIsTestOnly`).
5. **Schema 6**: a fresh store has the five tables with `webauthn_credentials` byte-equal (normalised)
   in the schema file, `passkeysSQL` and the canonical DDL, `user_version` 6 (PS1); v1, v3, v4 and v5
   stores upgrade once, idempotently, every row kept, and a session open never upgrades (PS4).
6. **`EnrolPasskey`** refuses each of S1–S10 (every operand) with `ErrBadApprover`, a taken name, an
   enrolled key and an enrolled credential ID with their typed errors naming the holder (also after
   revoke), and a v5 store with the passkeys refusal — writing nothing to either table each time; it
   writes both rows together (PS2); `EnrolApprover` still refuses `webauthn` (PS5).
7. **The lookups** return exact rows, miss without error, exclude revoked and non-webauthn rows from
   `ActivePasskeys`, order by name and cap at 64, refuse a v5 store, and read every `PasskeyRow` field
   back (PS3, RC6); all four new methods require a deadline (`TestCredentialsMethodsRequireADeadline`).
8. **The consume**: C13–C15 refuse before the clock is read (CS7) and sit after C4 (the unedited CS2
   arm still gets `ErrInconsistentApproval`); the counter moves forward only, with `last_used_at`,
   `approved_by = webauthn:x` and the session in one transaction, and every refusal (C16–C18, a revoke)
   rolls all of it back, including when the row changes inside the transaction's window (CS8, CS9); CS3's
   one named edit; RC3 classifies `ApproverUse.Passkey` and `PasskeyUse`.
9. **Authority**: a verified passkey approval mints with `approved_by = webauthn:<name>` and the counter
   advanced; skipping A6 is refused by the store (AU3, MUT-A6).
10. **Carry-forwards**: SOURCE.txt names `chrome-*.json` and XW0 pins it; the aliasing pins hold; the
    zero-challenge/`Now` decision is recorded (S2-D9–S2-D11).
11. Each §4 mutant, applied alone, makes its named test **and** its package suite fail; the named
    equivalents survive for the stated reasons; the base controls behave as stated.
12. Every §6 gate passes, including every inventory and reflective test by name and `VERIFY_GO_LEGS=plain
    ./scripts/verify_go.sh` after the commit (modulo §6's re-run rule for `host/broker`); no test skips;
    `git status --short` lists only §5.
13. The numstat after each sub-milestone's commit equals §7's value; the sprint total is ≤ 2,000.
14. No real key, browser, authenticator, store, `~/.ssh` access, mint or publish anywhere in the sprint
    (D-WORLD-109); `host/approval` stays pure and standard-library-only.

---

## Appendix C — for sprint 3 (nothing moved from this sprint)

The §9 sprint-2 fallback did not fire (§1.5), so **sprint 3 is M6 and M7 as designed**. What sprint 3
inherits from this sprint, and must use as built:
- `approval.ApproveByAssertion(AssertionInput{Admit, Challenge, Response, Approver})` with
  `Approver.Passkey` filled from the store row: `RPID`, `Origin`, `CredentialID`, `Key` (parse the
  stored canonical COSE key), `BackupEligible`, and **`SignCount` = the stored counter** (W15 compares
  against it; the store re-checks in C18).
- `approval.ParseAssertionResponse` (sprint 1) gives the `rawId` for the lookup before any network
  call (WF3, WF5); `store.PasskeyByCredentialID(ctx, id) (ApproverRow, PasskeyRow, bool, error)`.
- `store.ActivePasskeys(ctx) ([]EnrolledPasskey, error)` for `--webauthn-link` (≤ 64, by name).
- `store.EnrolPasskey(ctx, ApproverRow{Format: "webauthn", PublicKey: Key.Canonical()}, PasskeyRow{…})`
  and its typed errors (`ErrApproverNameTaken`, `ErrApproverKeyEnrolled`, `ErrPasskeyCredentialEnrolled`)
  for N12; its S4 refuses `EnrolledAt <= 0`.
- `store.Passkey(ctx, name)` for `approver list`'s L4/L5 (revoked included).
- `approval.ErrApproverNotPasskey` (P3w), `ErrUserHandle`, `ErrBackupEligibility`,
  `ErrSignCountRegression`; `store.ErrSignCountRegression`, `ErrPasskeyChanged`.
- **A decision to revisit (S2-D11):** WN5 may also refuse an all-zero `env.random` challenge; this
  sprint left the pure law without that branch.
- The `approvaltest.Passkey` `Edit` knobs cover every CW arm's file (UV-clear, another origin, a
  counter, `RawID`); its `Register`/`Assert` take no name (S2-D6).

---

## Appendix P — the patches (apply with §2; generated from the compiled scratch stages; do not edit)

### Patch M3b1 — the assertion law, `checkApprover`, `PasskeyUse`

````diff M3b1
--- a/host/approval/passkey_assert.go
+++ b/host/approval/passkey_assert.go
@@ -0,0 +1,95 @@
+package approval
+
+// Row 201 A2 (design §2.9): the assertion law. A mint request is approved by a passkey only when
+// it is admitted, the approver is an active webauthn approver, and its passkey signed, with user
+// verification, this request's message hash at the enrolled origin and relying party ID, with a
+// counter that moved forward. No I/O.
+
+import (
+	"bytes"
+	"crypto/sha256"
+	"errors"
+	"fmt"
+)
+
+// AssertionInput is what ApproveByAssertion judges.
+type AssertionInput struct {
+	Admit     AdmitInput // the existing admission input (row times, consumed, both digests)
+	Challenge [32]byte   // the request row's challenge
+	Response  []byte     // the assertion file
+	Approver  Approver   // the enrolled row the assertion's credential ID names, with Passkey set
+}
+
+var (
+	ErrApproverNotPasskey  = errors.New("approval: the approver is not a webauthn approver with an enrolled passkey")
+	ErrUserHandle          = errors.New("approval: the assertion's userHandle is not the approver's name")
+	ErrBackupEligibility   = errors.New("approval: the assertion's backup eligibility (BE) is not the enrolled passkey's")
+	ErrSignCountRegression = errors.New("approval: the passkey's sign counter did not advance (a cloned authenticator, or an older assertion)")
+)
+
+// ApproveByAssertion is the assertion law (P1-P5, P3w, W1-W15, in this order): admission and the
+// approver's state first, the credential ID before any signed content is parsed, and the counter
+// last, so it is judged only on a genuine signature.
+func ApproveByAssertion(in AssertionInput) (*SignedApproval, error) {
+	if v := Admit(in.Admit); v != Admitted { // P1
+		return nil, &NotAdmittedError{Verdict: v}
+	}
+	if err := checkApprover(in.Admit, in.Approver, FormatWebAuthn); err != nil { // P2-P5, P3w
+		return nil, err
+	}
+	p := in.Approver.Passkey
+	resp, err := ParseAssertionResponse(in.Response) // W1
+	if err != nil {
+		return nil, err
+	}
+	if !bytes.Equal(resp.RawID, p.CredentialID) { // W2
+		return nil, ErrCredentialIDMismatch
+	}
+	if resp.HasUserHandle && !bytes.Equal(resp.UserHandle, []byte(in.Approver.Name)) { // W3
+		return nil, ErrUserHandle
+	}
+	cd, err := ParseClientData(resp.ClientDataJSON) // W4
+	if err != nil {
+		return nil, err
+	}
+	if cd.Type != "webauthn.get" { // W5
+		return nil, ErrClientDataType
+	}
+	h, err := MessageHash(in.Admit.StoredDigest, in.Challenge, in.Admit.ExpiresAt)
+	if err != nil || cd.Challenge != b64Challenge(h[:]) { // W6
+		return nil, ErrChallenge
+	}
+	if cd.Origin != p.Origin { // W7
+		return nil, ErrOriginMismatch
+	}
+	if cd.CrossOrigin || cd.HasTopOrigin { // W8
+		return nil, ErrCrossOrigin
+	}
+	ad, err := ParseAuthData(resp.AuthenticatorData, false) // W9
+	if err != nil {
+		return nil, err
+	}
+	if ad.RPIDHash != sha256.Sum256([]byte(p.RPID)) { // W10
+		return nil, ErrRPIDHash
+	}
+	if !ad.UP { // W11
+		return nil, ErrUserPresence
+	}
+	if !ad.UV { // W12
+		return nil, ErrUserVerification
+	}
+	if ad.BE != p.BackupEligible { // W13
+		return nil, ErrBackupEligibility
+	}
+	cdh := sha256.Sum256(resp.ClientDataJSON)
+	if err := p.Key.verify(append(bytes.Clone(resp.AuthenticatorData), cdh[:]...), resp.Signature); err != nil { // W14
+		return nil, err
+	}
+	if (ad.SignCount != 0 || p.SignCount != 0) && ad.SignCount <= p.SignCount { // W15
+		return nil, fmt.Errorf("%w: enrolled %d, this assertion %d; if you did not cause this, run: ailang-worldd approver revoke %s",
+			ErrSignCountRegression, p.SignCount, ad.SignCount, in.Approver.Name)
+	}
+	return &SignedApproval{ok: true, digest: in.Admit.StoredDigest, challenge: in.Challenge,
+		expiry: in.Admit.ExpiresAt, name: in.Approver.Name, format: FormatWebAuthn, key: p.Key.Canonical(),
+		credentialID: bytes.Clone(p.CredentialID), signCount: ad.SignCount}, nil
+}
--- a/host/approval/passkey_assert_test.go
+++ b/host/approval/passkey_assert_test.go
@@ -0,0 +1,130 @@
+package approval
+
+// Row 201 A2 sprint 2 M3b (design §6, PA2-PA4, XW2): the assertion law over Chrome's recorded
+// assertions, and what a verified passkey approval hands the consume. The refusal arms are
+// passkey_law_test.go's (package approval_test: they need approvaltest, which imports this package).
+
+import (
+	"bytes"
+	"encoding/binary"
+	"errors"
+	"reflect"
+	"testing"
+)
+
+// chromeAssertionInput is the recorded assertion (counter 2) against the passkey the recorded
+// registration enrols (counter 1), admitted for the fixture request whose message hash both signed.
+func chromeAssertionInput(t *testing.T, alg string) AssertionInput {
+	t.Helper()
+	c, reg := chromeRegistration(t, alg)
+	a, err := NewPasskeyApprover(reg)
+	if err != nil {
+		t.Fatal(err)
+	}
+	d := fixtureDigest(t)
+	return AssertionInput{Admit: AdmitInput{Now: fixtureExpiry - 300, CreatedAt: fixtureExpiry - 600, ExpiresAt: fixtureExpiry,
+		StoredDigest: d, BytesDigest: d}, Challenge: fixtureChallenge(), Response: c.Assertion, Approver: a}
+}
+
+// chromeApproval is ApproveByAssertion(chromeAssertionInput(t, "es256")), which must succeed.
+func chromeApproval(t *testing.T) *SignedApproval {
+	t.Helper()
+	s, err := ApproveByAssertion(chromeAssertionInput(t, "es256"))
+	if err != nil || s == nil {
+		t.Fatalf("the recorded assertion: (%v, %v), want a verified approval", s, err)
+	}
+	return s
+}
+
+// PA2 (XW2): both recorded assertions approve the fixture request through the production law.
+func TestApproveByAssertionAcceptsChrome(t *testing.T) {
+	for _, alg := range []string{"es256", "eddsa"} {
+		in := chromeAssertionInput(t, alg)
+		s, err := ApproveByAssertion(in)
+		if err != nil {
+			t.Fatalf("%s: %v", alg, err)
+		}
+		if err := s.Binds(fixtureDigest(t), fixtureChallenge(), fixtureExpiry); err != nil {
+			t.Fatalf("%s: Binds(the fixture request) = %v", alg, err)
+		}
+		if got := s.ApprovedBy(); got != "webauthn:mark-mac" {
+			t.Fatalf("%s: ApprovedBy() = %q", alg, got)
+		}
+		n, f, k, ok := s.Use()
+		if !ok || n != "mark-mac" || f != FormatWebAuthn || !bytes.Equal(k, in.Approver.Passkey.Key.Canonical()) {
+			t.Fatalf("%s: Use() = (%q, %q, %x, %v), want the canonical COSE key", alg, n, f, k, ok)
+		}
+		id, count, ok := s.PasskeyUse()
+		if !ok || !bytes.Equal(id, chromeFixture(t, alg).RegRawID) || count != 2 {
+			t.Fatalf("%s: PasskeyUse() = (%x, %d, %v), want the credential ID, 2, true", alg, id, count, ok)
+		}
+	}
+}
+
+// PA3 (XW2, last arm): the recorded assertion (counter 2) is refused once the stored counter is
+// the second recorded assertion's (3, design §1.4 P2 of sprint 1).
+func TestChromeAssertionAfterAssertion2IsRefused(t *testing.T) {
+	for _, alg := range []string{"es256", "eddsa"} {
+		in := chromeAssertionInput(t, alg)
+		in.Approver.Passkey.SignCount = binary.BigEndian.Uint32(chromeFixture(t, alg).Asr2AuthData[33:37])
+		if in.Approver.Passkey.SignCount != 3 {
+			t.Fatalf("%s: assertion2's counter %d, want 3", alg, in.Approver.Passkey.SignCount)
+		}
+		if s, err := ApproveByAssertion(in); s != nil || !errors.Is(err, ErrSignCountRegression) {
+			t.Fatalf("%s: (%v, %v), want ErrSignCountRegression", alg, s, err)
+		}
+	}
+}
+
+// PA4: PasskeyUse is ok only for a verified webauthn approval (B7, B8), and neither the input
+// approver nor a returned slice can change it afterwards (judge carry-forward, iteration 379).
+func TestSignedApprovalPasskeyUse(t *testing.T) {
+	var nilApproval *SignedApproval
+	for name, s := range map[string]*SignedApproval{
+		"nil":                     nilApproval,
+		"zero":                    {},
+		"unverified-webauthn":     {format: FormatWebAuthn, name: "mark-mac", credentialID: []byte{1}, signCount: 2},
+		"verified-ssh":            verifiedApproval(t),
+		"verified-ssh-with-count": func() *SignedApproval { s := verifiedApproval(t); s.credentialID, s.signCount = []byte{1}, 2; return s }(),
+	} {
+		if id, n, ok := s.PasskeyUse(); ok || id != nil || n != 0 {
+			t.Errorf("%s.PasskeyUse() = (%x, %d, %v), want (nil, 0, false)", name, id, n, ok)
+		}
+	}
+	in := chromeAssertionInput(t, "es256")
+	want := bytes.Clone(in.Approver.Passkey.CredentialID)
+	wantKey := in.Approver.Passkey.Key.Canonical()
+	s, err := ApproveByAssertion(in)
+	if err != nil {
+		t.Fatal(err)
+	}
+	id, _, _ := s.PasskeyUse()
+	id[0] ^= 0xff                               // the caller's copy
+	in.Approver.Passkey.CredentialID[1] ^= 0xff // the input ApproveByAssertion read
+	in.Approver.Passkey.Key.X[0] ^= 0xff
+	if id2, n, ok := s.PasskeyUse(); !ok || n != 2 || !bytes.Equal(id2, want) {
+		t.Fatalf("PasskeyUse() after editing a returned ID and the input = (%x, %d, %v), want (%x, 2, true)", id2, n, ok, want)
+	}
+	if _, _, k, _ := s.Use(); !bytes.Equal(k, wantKey) {
+		t.Fatalf("Use() after editing the input key = %x, want %x", k, wantKey)
+	}
+}
+
+// The registration law's result does not alias its input (judge carry-forward, iteration 379).
+func TestNewPasskeyApproverDoesNotAliasItsInput(t *testing.T) {
+	_, in := chromeRegistration(t, "es256")
+	in.Response = bytes.Clone(in.Response)
+	a, err := NewPasskeyApprover(in)
+	if err != nil {
+		t.Fatal(err)
+	}
+	p := *a.Passkey
+	want := Approver{Name: a.Name, Format: a.Format, EnrolledAt: a.EnrolledAt, Passkey: &p}
+	p.CredentialID, p.Key.X, p.Key.Y = bytes.Clone(p.CredentialID), bytes.Clone(p.Key.X), bytes.Clone(p.Key.Y)
+	for i := range in.Response {
+		in.Response[i] = 'x'
+	}
+	if !reflect.DeepEqual(a, want) {
+		t.Fatalf("the approver changed when its input did: %+v, want %+v", a.Passkey, want.Passkey)
+	}
+}
--- a/host/approval/signed.go
+++ b/host/approval/signed.go
@@ -12,15 +12,18 @@
 )
 
 // SignedApproval is a verified signature approval. Its fields are unexported, so only
-// ApproveBySignature, in this package, can return one with ok set; the zero value binds nothing.
+// ApproveBySignature or ApproveByAssertion, in this package, can return one with ok set; the zero
+// value binds nothing.
 type SignedApproval struct {
-	ok        bool
-	digest    [32]byte
-	challenge [32]byte
-	expiry    int64
-	name      string
-	format    string
-	key       []byte
+	ok           bool
+	digest       [32]byte
+	challenge    [32]byte
+	expiry       int64
+	name         string
+	format       string
+	key          []byte
+	credentialID []byte // webauthn only (row 201 A2): the passkey that signed
+	signCount    uint32 // webauthn only: the assertion's sign counter
 }
 
 // SignatureInput is what ApproveBySignature judges.
@@ -56,18 +59,9 @@
 	if v := Admit(in.Admit); v != Admitted { // P1
 		return nil, &NotAdmittedError{Verdict: v}
 	}
-	if !ApproverNamePattern.MatchString(in.Approver.Name) { // P2
-		return nil, ErrApproverName
+	if err := checkApprover(in.Admit, in.Approver, FormatSSH); err != nil { // P2-P5
+		return nil, err
 	}
-	if in.Approver.Format != FormatSSH { // P3
-		return nil, ErrApproverFormat
-	}
-	if in.Approver.RevokedAt != 0 { // P4
-		return nil, ErrApproverRevoked
-	}
-	if in.Admit.Now < in.Approver.EnrolledAt { // P5
-		return nil, ErrApproverClockBehind
-	}
 	if err := VerifyApproval(in.Signature, in.Admit.StoredDigest, in.Challenge, in.Admit.ExpiresAt, in.Approver.Key); err != nil { // P6
 		return nil, err
 	}
@@ -76,6 +70,27 @@
 		key: bytes.Clone(in.Approver.Key.Blob)}, nil
 }
 
+// checkApprover is P2-P5, shared by both laws (row 201 A2 §2.9): the name, the format (a webauthn
+// approver must carry its passkey, P3w), not revoked, and not enrolled after now.
+func checkApprover(admit AdmitInput, a Approver, format string) error {
+	if !ApproverNamePattern.MatchString(a.Name) { // P2
+		return ErrApproverName
+	}
+	if format == FormatWebAuthn && (a.Format != format || a.Passkey == nil) { // P3w
+		return ErrApproverNotPasskey
+	}
+	if a.Format != format { // P3
+		return ErrApproverFormat
+	}
+	if a.RevokedAt != 0 { // P4
+		return ErrApproverRevoked
+	}
+	if admit.Now < a.EnrolledAt { // P5
+		return ErrApproverClockBehind
+	}
+	return nil
+}
+
 // Binds reports whether s is a verified approval of exactly this digest, challenge and expiry.
 func (s *SignedApproval) Binds(digest, challenge [32]byte, expiry int64) error {
 	if s == nil || !s.ok { // B1
@@ -110,3 +125,12 @@
 	}
 	return s.name, s.format, bytes.Clone(s.key), true
 }
+
+// PasskeyUse is what the consume binds for a passkey (row 201 A2 §2.9): the credential ID (a copy)
+// and the assertion's counter; ok only for a verified webauthn approval.
+func (s *SignedApproval) PasskeyUse() (credentialID []byte, signCount uint32, ok bool) {
+	if s == nil || !s.ok || s.format != FormatWebAuthn { // B7, B8
+		return nil, 0, false
+	}
+	return bytes.Clone(s.credentialID), s.signCount, true
+}
--- a/host/approval/signed_test.go
+++ b/host/approval/signed_test.go
@@ -163,13 +163,15 @@
 
 // signedApprovalFields is the disposition of every SignedApproval field: what reads it.
 var signedApprovalFields = map[string]string{
-	"ok":        "B1/B5/B6: an unverified value binds nothing and has no approver",
-	"digest":    "B2: Binds",
-	"challenge": "B3: Binds",
-	"expiry":    "B4: Binds",
-	"name":      "ApprovedBy and Use",
-	"format":    "ApprovedBy and Use",
-	"key":       "Use",
+	"ok":           "B1/B5/B6: an unverified value binds nothing and has no approver",
+	"digest":       "B2: Binds",
+	"challenge":    "B3: Binds",
+	"expiry":       "B4: Binds",
+	"name":         "ApprovedBy and Use",
+	"format":       "ApprovedBy and Use",
+	"key":          "Use",
+	"credentialID": "PasskeyUse (B7, B8): the consume binds the passkey that signed",
+	"signCount":    "PasskeyUse (B7, B8): the consume moves the counter to it",
 }
 
 // RC1: every SignedApproval field has a disposition, and changing any one of them in a verified
@@ -185,27 +187,35 @@
 		t.Fatalf("SignedApproval has %d fields, signedApprovalFields %d", typ.NumField(), len(signedApprovalFields))
 	}
 	perturb := map[string]func(*SignedApproval){
-		"ok":        func(s *SignedApproval) { s.ok = false },
-		"digest":    func(s *SignedApproval) { s.digest[0] ^= 1 },
-		"challenge": func(s *SignedApproval) { s.challenge[0] ^= 1 },
-		"expiry":    func(s *SignedApproval) { s.expiry++ },
-		"name":      func(s *SignedApproval) { s.name += "x" },
-		"format":    func(s *SignedApproval) { s.format += "x" },
-		"key":       func(s *SignedApproval) { s.key = bytes.Clone(s.key); s.key[0] ^= 0xff },
+		"ok":           func(s *SignedApproval) { s.ok = false },
+		"digest":       func(s *SignedApproval) { s.digest[0] ^= 1 },
+		"challenge":    func(s *SignedApproval) { s.challenge[0] ^= 1 },
+		"expiry":       func(s *SignedApproval) { s.expiry++ },
+		"name":         func(s *SignedApproval) { s.name += "x" },
+		"format":       func(s *SignedApproval) { s.format += "x" },
+		"key":          func(s *SignedApproval) { s.key = bytes.Clone(s.key); s.key[0] ^= 0xff },
+		"credentialID": func(s *SignedApproval) { s.credentialID = append(bytes.Clone(s.credentialID), 0) },
+		"signCount":    func(s *SignedApproval) { s.signCount++ },
 	}
 	if len(perturb) != len(signedApprovalFields) {
 		t.Fatalf("%d perturbations for %d fields", len(perturb), len(signedApprovalFields))
 	}
 	observe := func(s *SignedApproval) string {
 		n, f, k, ok := s.Use()
-		return fmt.Sprintf("%v|%s|%s|%s|%x|%v", s.Binds(fixtureDigest(t), fixtureChallenge(), fixtureExpiry), s.ApprovedBy(), n, f, k, ok)
+		id, count, pok := s.PasskeyUse()
+		return fmt.Sprintf("%v|%s|%s|%s|%x|%v|%x|%d|%v", s.Binds(fixtureDigest(t), fixtureChallenge(), fixtureExpiry), s.ApprovedBy(), n, f, k, ok, id, count, pok)
 	}
-	base := observe(verifiedApproval(t))
+	// Every field is perturbed in a verified ssh approval and in a verified passkey approval
+	// (row 201 A2); each must change what one of them reports.
 	for field, p := range perturb {
-		s := verifiedApproval(t)
-		p(s)
-		if observe(s) == base {
-			t.Errorf("changing %s changed nothing Binds, ApprovedBy or Use report", field)
+		changed := false
+		for _, verified := range []func(*testing.T) *SignedApproval{verifiedApproval, chromeApproval} {
+			s := verified(t)
+			p(s)
+			changed = changed || observe(s) != observe(verified(t))
 		}
+		if !changed {
+			t.Errorf("changing %s changed nothing Binds, ApprovedBy, Use or PasskeyUse report", field)
+		}
 	}
 }
--- a/host/approval/testdata/webauthn/SOURCE.txt
+++ b/host/approval/testdata/webauthn/SOURCE.txt
@@ -13,8 +13,8 @@
 challenge_hex, the registration and the first assertion exactly as PublicKeyCredential.toJSON()
 returned them, and the second assertion's authenticatorData (the probe kept only that part).
 Normalised from the probe outputs by the sprint-1 planner (TestPasskeyFixturesAreTheGeneratedBytes
-pins the bytes):
-  es256.json  sha256 c0370982dbca416a7677b4fa1ff2cd64e9e533cefc608d3dcab42a1000d6caf6
-  eddsa.json  sha256 1544721fc45dd65a553a5baab8a43fcd4df0375b0b8c2c897a199d2526f6f719
+pins the bytes of the two fixtures; the probe's own output files, not committed, were):
+  chrome-es256.json  from the probe's es256.json, sha256 c0370982dbca416a7677b4fa1ff2cd64e9e533cefc608d3dcab42a1000d6caf6
+  chrome-eddsa.json  from the probe's eddsa.json, sha256 1544721fc45dd65a553a5baab8a43fcd4df0375b0b8c2c897a199d2526f6f719
 
 Regeneration lands in sprint 3 (M7, the committed DevTools driver); see generate.sh.
--- a/host/approval/webauthn_fixture_test.go
+++ b/host/approval/webauthn_fixture_test.go
@@ -100,7 +100,7 @@
 	if err != nil {
 		t.Fatal(err)
 	}
-	for _, want := range []string{"Chrome/157.0.8103.50", "virtual authenticator", "D-WORLD-109"} {
+	for _, want := range []string{"Chrome/157.0.8103.50", "virtual authenticator", "D-WORLD-109", "chrome-es256.json  from", "chrome-eddsa.json  from"} {
 		if !strings.Contains(string(src), want) {
 			t.Errorf("SOURCE.txt does not name %q", want)
 		}
````

### Patch M3b2 — the test authenticator

````diff M3b2
--- a/host/approval/approvaltest/passkey.go
+++ b/host/approval/approvaltest/passkey.go
@@ -0,0 +1,164 @@
+package approvaltest
+
+// Row 201 A2 (design §2.15): a test authenticator. It makes registration and assertion files in
+// Chrome 157's measured layout (TestTestPasskeyMatchesChromeLayout) from a fixed seed: no browser,
+// no real authenticator, no randomness in the key (D-WORLD-109). Edit makes each refusal arm.
+
+import (
+	"bytes"
+	"crypto/ecdsa"
+	"crypto/ed25519"
+	"crypto/elliptic"
+	"crypto/rand"
+	"crypto/sha256"
+	"encoding/base64"
+	"encoding/binary"
+	"encoding/json"
+	"strconv"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+)
+
+// Passkey is an ES256 or EdDSA test authenticator with one 32-byte credential.
+type Passkey struct {
+	ed ed25519.PrivateKey
+	ec *ecdsa.PrivateKey
+	id []byte
+}
+
+// NewPasskey returns the authenticator whose 32-byte private key is the byte seed (1-254), repeated.
+func NewPasskey(alg int64, seed byte) *Passkey {
+	raw := bytes.Repeat([]byte{seed}, 32)
+	id := sha256.Sum256(append([]byte("approvaltest-credential"), byte(alg), seed))
+	p := &Passkey{id: id[:]}
+	switch alg {
+	case approval.AlgEdDSA:
+		p.ed = ed25519.NewKeyFromSeed(raw)
+	case approval.AlgES256:
+		k, err := ecdsa.ParseRawPrivateKey(elliptic.P256(), raw)
+		if err != nil {
+			panic(err)
+		}
+		p.ec = k
+	default:
+		panic("approvaltest: NewPasskey needs approval.AlgES256 or approval.AlgEdDSA")
+	}
+	return p
+}
+
+// CredentialID is a copy of the credential ID.
+func (p *Passkey) CredentialID() []byte { return bytes.Clone(p.id) }
+
+// Key is the public key as host/approval parses it.
+func (p *Passkey) Key() approval.COSEKey {
+	if p.ed != nil {
+		return approval.COSEKey{Alg: approval.AlgEdDSA, X: bytes.Clone(p.ed.Public().(ed25519.PublicKey))}
+	}
+	b, err := p.ec.PublicKey.Bytes()
+	if err != nil {
+		panic(err)
+	}
+	return approval.COSEKey{Alg: approval.AlgES256, X: b[1:33], Y: b[33:65]}
+}
+
+// Edit changes one ceremony from Chrome's layout; the zero Edit is Chrome's.
+type Edit struct {
+	Type         string // clientDataJSON type ("" is the ceremony's own)
+	CrossOrigin  bool
+	TopOrigin    string // a topOrigin member when not ""
+	ClientData   []byte // replaces the whole clientDataJSON when not nil
+	Set, Clear   byte   // authData flag bits set, then cleared (UP 0x01, UV 0x04, BE 0x08, BS 0x10, AT 0x40)
+	Extensions   bool   // ED and the extensions map {"credProtect": 2}
+	UserHandle   []byte // an assertion's userHandle member when not nil (Chrome sends none)
+	RawID        []byte // replaces id and rawId when not nil
+	BadSignature bool   // flips the assertion signature's last byte
+}
+
+func one(e []Edit) Edit {
+	if len(e) == 0 {
+		return Edit{}
+	}
+	return e[0]
+}
+
+// Register is a registration file for challenge at origin (fmt "none", flags UP|UV|AT, counter 1).
+func (p *Passkey) Register(origin, rpID string, challenge [32]byte, e ...Edit) []byte {
+	ed := one(e)
+	att := append(bytes.Repeat([]byte{1, 2, 3, 4, 5, 6, 7, 8}, 2), byte(len(p.id)>>8), byte(len(p.id)))
+	att = append(append(att, p.id...), p.Key().Canonical()...)
+	ad := authData(rpID, 0x45, 1, ed, att)
+	obj := append([]byte("\xa3\x63fmt\x64none\x67attStmt\xa0\x68authData\x58"), byte(len(ad)))
+	return envelope(p, ed, map[string]string{"clientDataJSON": b64(clientData("webauthn.create", challenge, origin, ed)),
+		"attestationObject": b64(append(obj, ad...))})
+}
+
+// Assert is an assertion file signing challenge at origin with counter (flags UP|UV).
+func (p *Passkey) Assert(origin, rpID string, challenge [32]byte, counter uint32, e ...Edit) []byte {
+	ed := one(e)
+	ad := authData(rpID, 0x05, counter, ed, nil)
+	cdj := clientData("webauthn.get", challenge, origin, ed)
+	h := sha256.Sum256(cdj)
+	sig := p.sign(append(bytes.Clone(ad), h[:]...))
+	if ed.BadSignature {
+		sig[len(sig)-1] ^= 1
+	}
+	r := map[string]string{"clientDataJSON": b64(cdj), "authenticatorData": b64(ad), "signature": b64(sig)}
+	if ed.UserHandle != nil {
+		r["userHandle"] = b64(ed.UserHandle)
+	}
+	return envelope(p, ed, r)
+}
+
+func (p *Passkey) sign(msg []byte) []byte {
+	if p.ed != nil {
+		return ed25519.Sign(p.ed, msg)
+	}
+	h := sha256.Sum256(msg)
+	sig, err := ecdsa.SignASN1(rand.Reader, p.ec, h[:])
+	if err != nil {
+		panic(err)
+	}
+	return sig
+}
+
+func authData(rpID string, flags byte, counter uint32, ed Edit, attested []byte) []byte {
+	f := (flags | ed.Set) &^ ed.Clear
+	if ed.Extensions {
+		f |= 0x80
+	}
+	h := sha256.Sum256([]byte(rpID))
+	b := append(binary.BigEndian.AppendUint32(append(h[:], f), counter), attested...)
+	if ed.Extensions {
+		b = append(b, "\xa1\x6bcredProtect\x02"...)
+	}
+	return b
+}
+
+func clientData(typ string, challenge [32]byte, origin string, ed Edit) []byte {
+	if ed.ClientData != nil {
+		return ed.ClientData
+	}
+	if ed.Type != "" {
+		typ = ed.Type
+	}
+	q := func(s string) string { b, _ := json.Marshal(s); return string(b) }
+	s := `{"type":` + q(typ) + `,"challenge":` + q(b64(challenge[:])) + `,"origin":` + q(origin) + `,"crossOrigin":` + strconv.FormatBool(ed.CrossOrigin)
+	if ed.TopOrigin != "" {
+		s += `,"topOrigin":` + q(ed.TopOrigin)
+	}
+	return []byte(s + "}")
+}
+
+func b64(b []byte) string { return base64.RawURLEncoding.EncodeToString(b) }
+
+func envelope(p *Passkey, ed Edit, response map[string]string) []byte {
+	id := p.id
+	if ed.RawID != nil {
+		id = ed.RawID
+	}
+	b, err := json.Marshal(map[string]any{"id": b64(id), "rawId": b64(id), "type": "public-key", "response": response})
+	if err != nil {
+		panic(err)
+	}
+	return b
+}
--- a/host/approval/passkey_testauth_test.go
+++ b/host/approval/passkey_testauth_test.go
@@ -0,0 +1,163 @@
+package approval_test
+
+// Row 201 A2 sprint 2 M3b2 (design §6, TP1, XW1): the test authenticator approvaltest.Passkey (a
+// fixed seed, never a real key: D-WORLD-109) makes files the production law accepts, in Chrome's
+// layout. This is an external test package because approvaltest imports host/approval (sprint 1
+// §1.4 P4).
+
+import (
+	"bytes"
+	"encoding/base64"
+	"encoding/json"
+	"os"
+	"path/filepath"
+	"strings"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
+)
+
+const paOrigin, paRPID = "https://x.ts.net:8443", "x.ts.net"
+
+var paDigest, paChallenge, paEnrol = [32]byte{1}, [32]byte{2}, [32]byte{3}
+
+// paHash is the WebAuthn challenge of the PA request: its message hash (design §2.1).
+func paHash(t *testing.T, digest, challenge [32]byte, expiry int64) [32]byte {
+	t.Helper()
+	h, err := approval.MessageHash(digest, challenge, expiry)
+	if err != nil {
+		t.Fatal(err)
+	}
+	return h
+}
+
+// TP1: for both algorithms the test authenticator's registration enrols and its assertion approves,
+// at a tailnet origin and at localhost; the approval names the credential and the canonical key.
+func TestTestPasskeyVerifies(t *testing.T) {
+	for _, alg := range []int64{approval.AlgES256, approval.AlgEdDSA} {
+		for origin, rpID := range map[string]string{paOrigin: paRPID, "http://localhost:7644": "localhost"} {
+			pk := approvaltest.NewPasskey(alg, 5)
+			a, err := approval.NewPasskeyApprover(approval.RegistrationInput{Name: "mark-mac", Origin: origin, Challenge: paEnrol,
+				Response: pk.Register(origin, rpID, paEnrol), Now: 900})
+			if err != nil || !bytes.Equal(a.Passkey.CredentialID, pk.CredentialID()) || !bytes.Equal(a.Passkey.Key.Canonical(), pk.Key().Canonical()) {
+				t.Fatalf("%d at %s: enrol: %v", alg, origin, err)
+			}
+			s, err := approval.ApproveByAssertion(approval.AssertionInput{
+				Admit:     approval.AdmitInput{Now: 1010, CreatedAt: 1000, ExpiresAt: 1600, StoredDigest: paDigest, BytesDigest: paDigest},
+				Challenge: paChallenge, Approver: a,
+				Response: pk.Assert(origin, rpID, paHash(t, paDigest, paChallenge, 1600), 2),
+			})
+			if err != nil {
+				t.Fatalf("%d at %s: assert: %v", alg, origin, err)
+			}
+			_, _, key, _ := s.Use()
+			if id, n, ok := s.PasskeyUse(); !ok || n != 2 || !bytes.Equal(id, pk.CredentialID()) || !bytes.Equal(key, pk.Key().Canonical()) {
+				t.Fatalf("%d at %s: PasskeyUse (%x, %d, %v), key %x", alg, origin, id, n, ok, key)
+			}
+		}
+	}
+}
+
+// chromeParts is one recorded ceremony's envelope members and decoded response fields.
+func chromeParts(t *testing.T, file []byte) (top, response []string, fields map[string][]byte) {
+	t.Helper()
+	var env, resp map[string]json.RawMessage
+	if json.Unmarshal(file, &env) != nil || json.Unmarshal(env["response"], &resp) != nil {
+		t.Fatalf("not an envelope: %s", file)
+	}
+	fields = map[string][]byte{}
+	for k := range env {
+		top = append(top, k)
+	}
+	for k, v := range resp {
+		response = append(response, k)
+		var s string
+		if json.Unmarshal(v, &s) == nil {
+			fields[k], _ = base64.RawURLEncoding.DecodeString(s)
+		}
+	}
+	return top, response, fields
+}
+
+// memberOrder is the member names of a JSON object, in order.
+func memberOrder(t *testing.T, b []byte) []string {
+	t.Helper()
+	d := json.NewDecoder(bytes.NewReader(b))
+	var names []string
+	if _, err := d.Token(); err != nil {
+		t.Fatal(err)
+	}
+	for d.More() {
+		k, err := d.Token()
+		if err != nil {
+			t.Fatal(err)
+		}
+		names = append(names, k.(string))
+		var skip json.RawMessage
+		if err := d.Decode(&skip); err != nil {
+			t.Fatal(err)
+		}
+	}
+	return names
+}
+
+// XW1: the test authenticator's files have Chrome's skeleton: the envelope members it writes are
+// Chrome's, the clientDataJSON member order, the attestationObject's bytes outside authData, the
+// flags, counters and AAGUID, the credential ID length, and the COSE key's bytes with x and y masked.
+func TestTestPasskeyMatchesChromeLayout(t *testing.T) {
+	for alg, name := range map[int64]string{approval.AlgES256: "es256", approval.AlgEdDSA: "eddsa"} {
+		raw, err := os.ReadFile(filepath.Join("testdata", "webauthn", "chrome-"+name+".json"))
+		if err != nil {
+			t.Fatal(err)
+		}
+		var rec struct{ Registration, Assertion json.RawMessage }
+		if err := json.Unmarshal(raw, &rec); err != nil {
+			t.Fatal(err)
+		}
+		pk := approvaltest.NewPasskey(alg, 7)
+		for _, c := range []struct {
+			kind         string
+			chrome, ours []byte
+		}{
+			{"registration", rec.Registration, pk.Register("http://localhost:56547", "localhost", [32]byte{})},
+			{"assertion", rec.Assertion, pk.Assert("http://localhost:56547", "localhost", [32]byte{}, 2)},
+		} {
+			ctop, cresp, cf := chromeParts(t, c.chrome)
+			otop, oresp, of := chromeParts(t, c.ours)
+			for _, sub := range []struct{ ours, chrome []string }{{otop, ctop}, {oresp, cresp}} {
+				for _, m := range sub.ours {
+					if !strings.Contains(" "+strings.Join(sub.chrome, " ")+" ", " "+m+" ") {
+						t.Errorf("%s %s: member %q is not one Chrome writes (%v)", name, c.kind, m, sub.chrome)
+					}
+				}
+			}
+			if got, want := strings.Join(memberOrder(t, of["clientDataJSON"]), ","), strings.Join(memberOrder(t, cf["clientDataJSON"]), ","); got != want {
+				t.Errorf("%s %s: clientDataJSON members %s, Chrome's %s", name, c.kind, got, want)
+			}
+			if c.kind == "assertion" {
+				if a, b := of["authenticatorData"], cf["authenticatorData"]; len(a) != 37 || len(b) != 37 || !bytes.Equal(a[32:], b[32:]) || !bytes.Equal(a[:32], b[:32]) {
+					t.Errorf("%s assertion authData %x, Chrome's %x", name, a, b)
+				}
+				continue
+			}
+			oa, ca := of["attestationObject"], cf["attestationObject"]
+			if len(oa) != len(ca) || !bytes.Equal(oa[:30], ca[:30]) { // fmt "none", attStmt {}, authData's head
+				t.Fatalf("%s attestationObject head %x (%d bytes), Chrome's %x (%d bytes)", name, oa[:30], len(oa), ca[:30], len(ca))
+			}
+			oad, cad := oa[30:], ca[30:]
+			mask := func(ad []byte) []byte {
+				m := bytes.Clone(ad)
+				copy(m[55:87], make([]byte, 32))  // the credential ID's bytes
+				copy(m[97:129], make([]byte, 32)) // x
+				if alg == approval.AlgES256 {
+					copy(m[132:164], make([]byte, 32)) // y
+				}
+				return m
+			}
+			if !bytes.Equal(mask(oad), mask(cad)) {
+				t.Errorf("%s registration authData skeleton:\n ours  %x\n chrome %x", name, mask(oad), mask(cad))
+			}
+		}
+	}
+}
````

### Patch M3b3 — the assertion law's refusals and its counter law

````diff M3b3
--- a/host/approval/passkey_law_test.go
+++ b/host/approval/passkey_law_test.go
@@ -0,0 +1,218 @@
+package approval_test
+
+// Row 201 A2 sprint 2 M3b3 (design §6, PA1, PA3): the assertion law's refusals, in its order, and
+// its counter law, made by the test authenticator (passkey_testauth_test.go).
+
+import (
+	"bytes"
+	"errors"
+	"fmt"
+	"strings"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
+)
+
+// paApprover is pk's registration enrolled as mark-mac at 900 (stored counter 1).
+func paApprover(t *testing.T, pk *approvaltest.Passkey) approval.Approver {
+	t.Helper()
+	a, err := approval.NewPasskeyApprover(approval.RegistrationInput{Name: "mark-mac", Origin: paOrigin, Challenge: paEnrol,
+		Response: pk.Register(paOrigin, paRPID, paEnrol), Now: 900})
+	if err != nil {
+		t.Fatalf("enrol the test passkey: %v", err)
+	}
+	return a
+}
+
+// paInput is a valid AssertionInput: pk's assertion (counter 2) of the PA request (created 1000,
+// expires 1600), admitted at 1010.
+func paInput(t *testing.T, pk *approvaltest.Passkey, e ...approvaltest.Edit) approval.AssertionInput {
+	t.Helper()
+	return approval.AssertionInput{
+		Admit:     approval.AdmitInput{Now: 1010, CreatedAt: 1000, ExpiresAt: 1600, StoredDigest: paDigest, BytesDigest: paDigest},
+		Challenge: paChallenge,
+		Response:  pk.Assert(paOrigin, paRPID, paHash(t, paDigest, paChallenge, 1600), 2, e...),
+		Approver:  paApprover(t, pk),
+	}
+}
+
+// paCase is one PA1 arm: the assertion is signer's (pk; another authenticator when other, under
+// pk's credential ID when underPK) at origin and rpID ("" is paOrigin, paRPID) over hash (nil is
+// the PA request's message hash) with counter count (0 is 2) and Edit e; then edit changes the
+// input; want is the refusal (verdict: admission's).
+type paCase struct {
+	name         string
+	e            approvaltest.Edit
+	origin, rpID string
+	hash         *[32]byte
+	count        uint32
+	other        bool
+	underPK      bool
+	edit         func(*approval.AssertionInput)
+	want         error
+	verdict      approval.Verdict
+}
+
+// PA1: one arm per P1-P5, P3w and W1-W15 on an otherwise valid assertion, in the law's order, for
+// both algorithms; each refusal returns a nil approval. Pairs pin the order.
+func TestApproveByAssertionOrder(t *testing.T) {
+	row, rowHash, zero := paChallenge, paHash(t, paChallenge, paChallenge, 1600), [32]byte{}
+	later := paHash(t, paDigest, paChallenge, 1601)
+	cases := []paCase{
+		{name: "P1-consumed", edit: func(in *approval.AssertionInput) { in.Admit.Consumed = true }, verdict: approval.RefuseConsumed},
+		{name: "P1-consumed-and-revoked", edit: func(in *approval.AssertionInput) { in.Admit.Consumed, in.Approver.RevokedAt = true, 950 }, verdict: approval.RefuseConsumed},
+		{name: "P1-clock-behind", edit: func(in *approval.AssertionInput) { in.Admit.Now = 999 }, verdict: approval.RefuseClockBehind},
+		{name: "P1-expired", edit: func(in *approval.AssertionInput) { in.Admit.Now = 1600 }, verdict: approval.RefuseExpired},
+		{name: "P1-bytes-changed", edit: func(in *approval.AssertionInput) { in.Admit.BytesDigest[0] ^= 1 }, verdict: approval.RefuseDigestMismatch},
+		{name: "P2-name", edit: func(in *approval.AssertionInput) { in.Approver.Name = "Mark" }, want: approval.ErrApproverName},
+		{name: "P3w-format-ssh", edit: func(in *approval.AssertionInput) { in.Approver.Format = approval.FormatSSH }, want: approval.ErrApproverNotPasskey},
+		{name: "P3w-no-passkey", edit: func(in *approval.AssertionInput) { in.Approver.Passkey = nil }, want: approval.ErrApproverNotPasskey},
+		{name: "P4-revoked-good-assertion", edit: func(in *approval.AssertionInput) { in.Approver.RevokedAt = 950 }, want: approval.ErrApproverRevoked},
+		{name: "P4-revoked-bad-assertion", e: approvaltest.Edit{BadSignature: true}, edit: func(in *approval.AssertionInput) { in.Approver.RevokedAt = 950 }, want: approval.ErrApproverRevoked},
+		{name: "P5-enrolled-after-now", edit: func(in *approval.AssertionInput) { in.Approver.EnrolledAt = 1011 }, want: approval.ErrApproverClockBehind},
+		{name: "W1-envelope", edit: func(in *approval.AssertionInput) {
+			in.Response = bytes.Replace(in.Response, []byte(`"public-key"`), []byte(`"x"`), 1)
+		}, want: approval.ErrPasskeyResponse},
+		{name: "W2-another-credentials-good-assertion", other: true, want: approval.ErrCredentialIDMismatch},
+		{name: "W2-raw-id", e: approvaltest.Edit{RawID: []byte("another credential")}, want: approval.ErrCredentialIDMismatch},
+		{name: "W3-user-handle", e: approvaltest.Edit{UserHandle: []byte("mallory")}, want: approval.ErrUserHandle},
+		{name: "W3-empty-user-handle", e: approvaltest.Edit{UserHandle: []byte{}}, want: approval.ErrUserHandle},
+		{name: "W4-client-data", e: approvaltest.Edit{ClientData: []byte(`{"type":"webauthn.get"}`)}, want: approval.ErrClientData},
+		{name: "W5-type-create", e: approvaltest.Edit{Type: "webauthn.create"}, want: approval.ErrClientDataType},
+		{name: "W5-before-W6", e: approvaltest.Edit{Type: "webauthn.create"}, hash: &row, want: approval.ErrClientDataType},
+		{name: "W6-the-row-challenge-not-the-hash", hash: &row, want: approval.ErrChallenge},
+		{name: "W6-another-expiry", hash: &later, want: approval.ErrChallenge},
+		{name: "W6-another-digest", hash: &rowHash, want: approval.ErrChallenge},
+		{name: "W6-negative-expiry-has-no-hash", hash: &zero, edit: func(in *approval.AssertionInput) {
+			in.Admit.Now, in.Admit.CreatedAt, in.Admit.ExpiresAt, in.Approver.EnrolledAt = -10, -20, -1, -30
+		}, want: approval.ErrChallenge},
+		{name: "W7-origin", origin: "https://x.ts.net", want: approval.ErrOriginMismatch},
+		{name: "W7-before-W10", origin: "https://y.ts.net", rpID: "y.ts.net", want: approval.ErrOriginMismatch},
+		{name: "W8-cross-origin", e: approvaltest.Edit{CrossOrigin: true}, want: approval.ErrCrossOrigin},
+		{name: "W8-top-origin", e: approvaltest.Edit{TopOrigin: paOrigin}, want: approval.ErrCrossOrigin},
+		{name: "W9-attested-data", e: approvaltest.Edit{Set: 0x40}, want: approval.ErrAuthDataFlags},
+		{name: "W9-backup-state-without-eligibility", e: approvaltest.Edit{Set: 0x10}, want: approval.ErrAuthDataFlags},
+		{name: "W10-rp-id", rpID: "y.ts.net", want: approval.ErrRPIDHash},
+		{name: "W14-bad-signature", e: approvaltest.Edit{BadSignature: true}, want: approval.ErrBadSignature},
+		{name: "W14-another-key-same-credential-id", other: true, underPK: true, want: approval.ErrBadSignature},
+		{name: "W14-before-W15", count: 1, e: approvaltest.Edit{BadSignature: true}, want: approval.ErrBadSignature},
+		{name: "W15-counter-not-advanced", count: 1, want: approval.ErrSignCountRegression},
+	}
+	// W11-W13 (flag-state coverage, iteration 379's judge): each refusal with every other flag
+	// state that reaches it, so a mutant that also needs another bit dies. BE marks the enrolled
+	// passkey backup-eligible.
+	for _, f := range []struct {
+		name               string
+		set, clear         byte
+		enrolledBE, badSig bool
+		want               error
+	}{
+		{"W11-no-UP", 0, 0x01, false, false, approval.ErrUserPresence},
+		{"W11-no-UP-or-UV", 0, 0x05, false, false, approval.ErrUserPresence},
+		{"W11-no-UP-BE", 0x08, 0x01, true, false, approval.ErrUserPresence},
+		{"W11-no-UP-BE-BS", 0x18, 0x01, true, false, approval.ErrUserPresence},
+		{"W11-no-UP-ED", 0x80, 0x01, false, false, approval.ErrUserPresence},
+		{"W12-no-UV", 0, 0x04, false, false, approval.ErrUserVerification},
+		{"W12-no-UV-BE", 0x08, 0x04, true, false, approval.ErrUserVerification},
+		{"W12-no-UV-BE-BS", 0x18, 0x04, true, false, approval.ErrUserVerification},
+		{"W12-no-UV-BE-not-enrolled", 0x08, 0x04, false, false, approval.ErrUserVerification},
+		{"W12-no-UV-ED", 0x80, 0x04, false, false, approval.ErrUserVerification},
+		{"W12-before-W14", 0, 0x04, false, true, approval.ErrUserVerification},
+		{"W13-BE-not-enrolled", 0x08, 0, false, false, approval.ErrBackupEligibility},
+		{"W13-BE-BS-not-enrolled", 0x18, 0, false, false, approval.ErrBackupEligibility},
+		{"W13-enrolled-BE-absent", 0, 0, true, false, approval.ErrBackupEligibility},
+		{"W13-enrolled-BE-absent-ED", 0x80, 0, true, false, approval.ErrBackupEligibility},
+		{"W13-before-W14", 0x08, 0, false, true, approval.ErrBackupEligibility},
+	} {
+		be := f.enrolledBE
+		cases = append(cases, paCase{name: f.name, want: f.want, edit: func(in *approval.AssertionInput) { in.Approver.Passkey.BackupEligible = be },
+			e: approvaltest.Edit{Set: f.set &^ 0x80, Clear: f.clear, Extensions: f.set&0x80 != 0, BadSignature: f.badSig}})
+	}
+	// Acceptance controls: each knob the arms use, alone and harmless, still verifies.
+	enrolBE := func(in *approval.AssertionInput) { in.Approver.Passkey.BackupEligible = true }
+	controls := []paCase{
+		{name: "clean"},
+		{name: "enrolled-at-now", edit: func(in *approval.AssertionInput) { in.Approver.EnrolledAt = in.Admit.Now }},
+		{name: "user-handle-is-the-name", e: approvaltest.Edit{UserHandle: []byte("mark-mac")}},
+		{name: "extensions", e: approvaltest.Edit{Extensions: true}},
+		{name: "BE-enrolled", e: approvaltest.Edit{Set: 0x08}, edit: enrolBE},
+		{name: "BE-BS-enrolled", e: approvaltest.Edit{Set: 0x18}, edit: enrolBE},
+	}
+	for _, alg := range []int64{approval.AlgES256, approval.AlgEdDSA} {
+		pk, other := approvaltest.NewPasskey(alg, 7), approvaltest.NewPasskey(approval.AlgEdDSA, 9)
+		for i, c := range append(cases, controls...) {
+			t.Run(fmt.Sprint(alg, "/", c.name), func(t *testing.T) {
+				in := paInput(t, pk)
+				signer, origin, rpID, hash, count := pk, paOrigin, paRPID, paHash(t, paDigest, paChallenge, 1600), uint32(2)
+				if c.other {
+					signer = other
+				}
+				if c.underPK {
+					c.e.RawID = pk.CredentialID()
+				}
+				if c.origin != "" {
+					origin = c.origin
+				}
+				if c.rpID != "" {
+					rpID = c.rpID
+				}
+				if c.hash != nil {
+					hash = *c.hash
+				}
+				if c.count != 0 {
+					count = c.count
+				}
+				in.Response = signer.Assert(origin, rpID, hash, count, c.e)
+				if c.edit != nil {
+					c.edit(&in)
+				}
+				s, err := approval.ApproveByAssertion(in)
+				var na *approval.NotAdmittedError
+				switch {
+				case i >= len(cases):
+					if err != nil || s.ApprovedBy() != "webauthn:mark-mac" || s.Binds(paDigest, paChallenge, 1600) != nil {
+						t.Fatalf("control: (%v, %v), want a verified approval of the PA request", s, err)
+					}
+				case s != nil:
+					t.Fatalf("a refused assertion returned %v", s)
+				case c.verdict != 0 && (!errors.As(err, &na) || na.Verdict != c.verdict):
+					t.Fatalf("err %v, want not admitted (%s)", err, c.verdict)
+				case c.verdict == 0 && !errors.Is(err, c.want):
+					t.Fatalf("err %v, want %v", err, c.want)
+				}
+			})
+		}
+	}
+}
+
+// PA3: the counter law (W15): both zero passes; otherwise the new counter must be strictly greater
+// than the enrolled one; the verified approval carries the new counter.
+func TestSignCountLaw(t *testing.T) {
+	pk := approvaltest.NewPasskey(approval.AlgEdDSA, 7)
+	for _, c := range []struct {
+		stored, new uint32
+		ok          bool
+	}{
+		{1, 2, true}, {0, 0, true}, {0, 7, true}, {0, 1, true}, {7, 4294967295, true}, {0, 4294967295, true},
+		{2, 2, false}, {3, 2, false}, {5, 0, false}, {1, 0, false}, {4294967295, 0, false}, {4294967295, 4294967295, false},
+	} {
+		in := paInput(t, pk)
+		in.Approver.Passkey.SignCount = c.stored
+		in.Response = pk.Assert(paOrigin, paRPID, paHash(t, paDigest, paChallenge, 1600), c.new)
+		s, err := approval.ApproveByAssertion(in)
+		if !c.ok {
+			if s != nil || !errors.Is(err, approval.ErrSignCountRegression) || !strings.Contains(err.Error(), "approver revoke mark-mac") {
+				t.Errorf("stored %d, new %d: (%v, %v), want ErrSignCountRegression naming approver revoke mark-mac", c.stored, c.new, s, err)
+			}
+			continue
+		}
+		if err != nil {
+			t.Errorf("stored %d, new %d: %v, want a verified approval", c.stored, c.new, err)
+			continue
+		}
+		if _, n, ok := s.PasskeyUse(); !ok || n != c.new {
+			t.Errorf("stored %d, new %d: PasskeyUse counter %d ok %v, want %d", c.stored, c.new, n, ok, c.new)
+		}
+	}
+}
````

### Patch M4a — credentials schema 6

````diff M4a
--- a/cmd/ailang-worldd/approve_sig_test.go
+++ b/cmd/ailang-worldd/approve_sig_test.go
@@ -216,7 +216,7 @@
 		{name: "v4-store", setup: func(t *testing.T, a approveRig) {
 			raw := rawCredentials(t, a.credPath)
 			defer func() { _ = raw.Close() }()
-			for _, stmt := range []string{`DROP TABLE approvers`, `PRAGMA user_version = 4`} {
+			for _, stmt := range []string{`DROP TABLE webauthn_credentials`, `DROP TABLE approvers`, `PRAGMA user_version = 4`} {
 				if _, err := raw.Exec(stmt); err != nil {
 					t.Fatal(err)
 				}
--- a/cmd/ailang-worldd/approver_rig_test.go
+++ b/cmd/ailang-worldd/approver_rig_test.go
@@ -57,14 +57,14 @@
 	return db, creds
 }
 
-// approverStoreV4 is approverStore taken back to schema 4: the approvers table dropped and
-// user_version 4, which is exactly the version-4 table set.
+// approverStoreV4 is approverStore taken back to schema 4: the webauthn_credentials and approvers
+// tables dropped and user_version 4, which is exactly the version-4 table set.
 func approverStoreV4(t *testing.T) (db, creds string) {
 	t.Helper()
 	db, creds = approverStore(t)
 	raw := rawCredentials(t, creds)
 	defer func() { _ = raw.Close() }()
-	for _, stmt := range []string{`DROP TABLE approvers`, `PRAGMA user_version = 4`} {
+	for _, stmt := range []string{`DROP TABLE webauthn_credentials`, `DROP TABLE approvers`, `PRAGMA user_version = 4`} {
 		if _, err := raw.Exec(stmt); err != nil {
 			t.Fatal(err)
 		}
--- a/host/store/approvers_test.go
+++ b/host/store/approvers_test.go
@@ -23,6 +23,18 @@
 	setPositiveFixtureVersion(t, db, 4)
 }
 
+// makeV5Credentials writes a schema-5 credentials store at path: makeV4Credentials, the approvers
+// table and user_version 5.
+func makeV5Credentials(t *testing.T, path string, rows []SessionRow) {
+	t.Helper()
+	makeV4Credentials(t, path, rows)
+	db := rawDB(t, path)
+	if _, err := db.Exec(approversSQL); err != nil {
+		t.Fatal(err)
+	}
+	setPositiveFixtureVersion(t, db, 5)
+}
+
 // approverCount counts the approvers rows of c's store.
 func approverCount(t *testing.T, c *Credentials) int {
 	t.Helper()
@@ -38,9 +50,9 @@
 	return ApproverRow{Name: name, Format: "ssh", PublicKey: []byte("key-" + name), EnrolledAt: at}
 }
 
-// AS5: the version 1, 3 and 4 stores upgrade to 5 in one step, keeping every session row, the
-// start history and the mint request row, and gaining an empty approvers table.
-func TestCredentialsUpgradeToV5(t *testing.T) {
+// AS5, PS4: the version 1, 3, 4 and 5 stores upgrade to 6 in one step, keeping every session row,
+// the start history, the mint request row and the approver, and gaining the empty tables.
+func TestCredentialsUpgradeToV6(t *testing.T) {
 	const reqID = "mr-00000000000000aa"
 	for _, tc := range []struct {
 		name string
@@ -73,6 +85,15 @@
 			}
 			return rows
 		}},
+		{"v5", 5, func(t *testing.T, path string) []SessionRow {
+			rows := []SessionRow{oneBoundRow()}
+			makeV5Credentials(t, path, rows)
+			if _, err := rawDB(t, path).Exec(`INSERT INTO approvers (name, format, public_key, enrolled_at, revoked_at) VALUES ('kept', 'ssh', x'6b', 900, 950);
+				INSERT INTO episode_starts (episode, start_json) VALUES ('w900', '{"episode":"w900","note":"kept"}')`); err != nil {
+				t.Fatal(err)
+			}
+			return rows
+		}},
 	} {
 		t.Run(tc.name, func(t *testing.T) {
 			path := filepath.Join(t.TempDir(), CredentialsFileName)
@@ -110,19 +131,31 @@
 			}
 			db := rawDB(t, path)
 			ddl := tableDDL(t, db)
-			requireExactTableNames(t, "upgraded credentials store", ddl, []string{"approvers", "episode_starts", "mint_requests", "session_credentials"})
+			requireExactTableNames(t, "upgraded credentials store", ddl, []string{"approvers", "episode_starts", "mint_requests", "session_credentials", "webauthn_credentials"})
 			if got, want := normalizeDDL(ddl["approvers"]), normalizeDDL(canonicalApproversDDL); got != want {
 				t.Fatalf("upgraded approvers DDL:\n got: %s\nwant: %s", got, want)
 			}
+			if got, want := normalizeDDL(ddl["webauthn_credentials"]), normalizeDDL(canonicalPasskeysDDL); got != want {
+				t.Fatalf("upgraded webauthn_credentials DDL:\n got: %s\nwant: %s", got, want)
+			}
 			ro, err := OpenCredentialsReadOnly(path)
 			if err != nil {
 				t.Fatal(err)
 			}
 			defer func() { _ = ro.Close() }()
-			if got := approverCount(t, ro); got != 0 {
+			ctx := testCtx(t)
+			if tc.v == 5 {
+				want := ApproverRow{Name: "kept", Format: "ssh", PublicKey: []byte("k"), EnrolledAt: 900, RevokedAt: 950}
+				if got, ok, err := ro.Approver(ctx, "kept"); err != nil || !ok || !reflect.DeepEqual(got, want) || approverCount(t, ro) != 1 {
+					t.Fatalf("approver after the upgrade: %+v ok %v err %v; want %+v alone", got, ok, err, want)
+				}
+			} else if got := approverCount(t, ro); got != 0 {
 				t.Fatalf("approvers after the upgrade = %d, want 0", got)
 			}
-			ctx := testCtx(t)
+			var passkeys int
+			if err := db.QueryRow(`SELECT count(*) FROM webauthn_credentials`).Scan(&passkeys); err != nil || passkeys != 0 {
+				t.Fatalf("passkeys after the upgrade = %d (%v), want 0", passkeys, err)
+			}
 			for _, r := range rows {
 				got, ok, err := ro.ResolveSession(ctx, r.CredentialID)
 				if err != nil || !ok || !reflect.DeepEqual(got, r) {
--- a/host/store/credentials.go
+++ b/host/store/credentials.go
@@ -38,8 +38,9 @@
 // credentialsSchemaVersion is the credentials store's own user_version,
 // independent of the world store's currentSchemaVersion. Version 2 added
 // profile bindings; version 3 adds immutable attended episode start history;
-// version 4 adds mint requests (row 201); version 5 adds approvers (row 201 A1b).
-const credentialsSchemaVersion = 5
+// version 4 adds mint requests (row 201); version 5 adds approvers (row 201 A1b); version 6 adds
+// passkeys (row 201 A2).
+const credentialsSchemaVersion = 6
 
 // CredentialsSchemaVersion exports credentialsSchemaVersion: serve compares a
 // store's version with it to decide whether to upgrade.
@@ -70,7 +71,9 @@
 // credentialsTables is the exact, sorted application table set of a credentials store at version.
 func credentialsTables(version int) []string {
 	switch {
-	case version >= 5:
+	case version >= 6:
+		return []string{approversTable, "episode_starts", mintRequestsTable, credentialsTable, passkeysTable}
+	case version == 5:
 		return []string{approversTable, "episode_starts", mintRequestsTable, credentialsTable}
 	case version == 4:
 		return []string{"episode_starts", mintRequestsTable, credentialsTable}
@@ -358,6 +361,9 @@
 	if version < 5 {
 		statements = append(statements, approversSQL)
 	}
+	if version < 6 {
+		statements = append(statements, passkeysSQL)
+	}
 	for _, stmt := range statements {
 		if _, err := tx.Exec(stmt); err != nil {
 			return fmt.Errorf("store: upgrade credentials %q: %w", display, err)
--- a/host/store/credentials_schema.sql
+++ b/host/store/credentials_schema.sql
@@ -58,8 +58,8 @@
 -- Approvers (version 5, row 201 A1b). A key enrolled at the daemon's terminal (a typed y, parent
 -- R3) may approve a mint request by signature. Revoking MARKS the row and never deletes it, so a
 -- name in mint_requests.approved_by ('ssh:<name>') names one key forever, and neither a name nor
--- a key is ever enrolled twice. format is checked in code (host/approval), not here, so A2 can add
--- 'webauthn' (and sign_count by ADD COLUMN) without rebuilding the table.
+-- a key is ever enrolled twice. format is checked in code (host/approval), not here, so A2 adds
+-- 'webauthn' without rebuilding the table (its passkey rows are webauthn_credentials, version 6).
 CREATE TABLE IF NOT EXISTS approvers (
     name         TEXT PRIMARY KEY CHECK (length(name) BETWEEN 1 AND 48),
     format       TEXT NOT NULL CHECK (format <> ''),
@@ -68,3 +68,18 @@
     last_used_at INTEGER CHECK (last_used_at IS NULL OR last_used_at >= enrolled_at),
     revoked_at   INTEGER CHECK (revoked_at IS NULL OR revoked_at >= enrolled_at)
 );
+
+-- Passkeys (version 6, row 201 A2). One row per webauthn approver (approvers.name, whose format is
+-- 'webauthn' and whose public_key is the canonical COSE key). The credential ID is never enrolled
+-- twice; the counter only moves forward, inside the consume that uses it. The credentials store
+-- does not enforce foreign keys, so EnrolPasskey writes both rows in one transaction.
+CREATE TABLE IF NOT EXISTS webauthn_credentials (
+    name            TEXT PRIMARY KEY CHECK (length(name) BETWEEN 1 AND 48),
+    credential_id   BLOB NOT NULL UNIQUE CHECK (length(credential_id) BETWEEN 1 AND 1023),
+    rp_id           TEXT NOT NULL CHECK (rp_id <> ''),
+    origin          TEXT NOT NULL CHECK (origin <> ''),
+    alg             INTEGER NOT NULL CHECK (alg IN (-7, -8)),
+    aaguid          BLOB NOT NULL CHECK (length(aaguid) = 16),
+    backup_eligible INTEGER NOT NULL CHECK (backup_eligible IN (0, 1)),
+    sign_count      INTEGER NOT NULL CHECK (sign_count BETWEEN 0 AND 4294967295)
+);
--- a/host/store/credentials_test.go
+++ b/host/store/credentials_test.go
@@ -61,10 +61,27 @@
     revoked_at   INTEGER CHECK (revoked_at IS NULL OR revoked_at >= enrolled_at)
 )`
 
+// canonicalPasskeysDDL is the current (version 6, row 201 A2) webauthn_credentials table.
+var canonicalPasskeysDDL = `CREATE TABLE webauthn_credentials (
+    name            TEXT PRIMARY KEY CHECK (length(name) BETWEEN 1 AND 48),
+    credential_id   BLOB NOT NULL UNIQUE CHECK (length(credential_id) BETWEEN 1 AND 1023),
+    rp_id           TEXT NOT NULL CHECK (rp_id <> ''),
+    origin          TEXT NOT NULL CHECK (origin <> ''),
+    alg             INTEGER NOT NULL CHECK (alg IN (-7, -8)),
+    aaguid          BLOB NOT NULL CHECK (length(aaguid) = 16),
+    backup_eligible INTEGER NOT NULL CHECK (backup_eligible IN (0, 1)),
+    sign_count      INTEGER NOT NULL CHECK (sign_count BETWEEN 0 AND 4294967295)
+)`
+
 func TestCredentialsSchemaMatchesCanonicalDDL(t *testing.T) {
 	c := openTestCredentials(t)
 	ddl := tableDDL(t, c.db)
-	requireExactTableNames(t, "fresh credentials store", ddl, []string{"approvers", "episode_starts", "mint_requests", "session_credentials"})
+	requireExactTableNames(t, "fresh credentials store", ddl, []string{"approvers", "episode_starts", "mint_requests", "session_credentials", "webauthn_credentials"})
+	for name, got := range map[string]string{"schema file": ddl["webauthn_credentials"], "passkeysSQL": passkeysSQL} {
+		if normalizeDDL(got) != normalizeDDL(canonicalPasskeysDDL) {
+			t.Fatalf("webauthn_credentials DDL (%s):\n got: %s\nwant: %s", name, got, canonicalPasskeysDDL)
+		}
+	}
 	if got, want := normalizeDDL(ddl["approvers"]), normalizeDDL(canonicalApproversDDL); got != want {
 		t.Fatalf("approvers DDL:\n got: %s\nwant: %s", got, want)
 	}
@@ -95,8 +112,8 @@
 		t.Fatalf("v4 world table differs from the version-1 credentials table: %s", got)
 	}
 	var version int
-	if err := c.db.QueryRow("PRAGMA user_version").Scan(&version); err != nil || version != credentialsSchemaVersion {
-		t.Fatalf("credentials user_version = %d %v, want %d", version, err, credentialsSchemaVersion)
+	if err := c.db.QueryRow("PRAGMA user_version").Scan(&version); err != nil || version != 6 {
+		t.Fatalf("credentials user_version = %d %v, want 6 (row 201 A2)", version, err)
 	}
 }
 
--- a/host/store/mint_requests_test.go
+++ b/host/store/mint_requests_test.go
@@ -129,7 +129,7 @@
 			}
 			db := rawDB(t, path)
 			ddl := tableDDL(t, db)
-			requireExactTableNames(t, "upgraded credentials store", ddl, []string{"approvers", "episode_starts", "mint_requests", "session_credentials"})
+			requireExactTableNames(t, "upgraded credentials store", ddl, []string{"approvers", "episode_starts", "mint_requests", "session_credentials", "webauthn_credentials"})
 			if got, want := normalizeDDL(ddl["mint_requests"]), normalizeDDL(canonicalMintRequestsDDL); got != want {
 				t.Fatalf("upgraded mint_requests DDL:\n got: %s\nwant: %s", got, want)
 			}
--- a/host/store/passkeys.go
+++ b/host/store/passkeys.go
@@ -0,0 +1,28 @@
+package store
+
+// Row 201 A2 (design_docs/planned/w-world-approvals-by-key-a2.md §2.10): passkeys (credentials
+// schema 6). A webauthn approver is an approvers row (format webauthn, public_key the canonical
+// COSE key) and one webauthn_credentials row, written together by EnrolPasskey after the typed
+// enrolment code and y. Like approvers, a passkey row is never deleted; its sign_count only moves
+// forward, inside the consume that uses it. The store does not import host/approval: these checks
+// are the store's own backstop for the pure law the CLI runs first.
+
+// ---- M4a ----
+
+const passkeysTable = "webauthn_credentials"
+
+// passkeysSchemaVersion is the first credentials schema with webauthn_credentials.
+const passkeysSchemaVersion = 6
+
+// passkeysSQL is the version 5 -> 6 upgrade: credentials_schema.sql's webauthn_credentials table,
+// column for column (TestCredentialsSchemaMatchesCanonicalDDL compares the two).
+const passkeysSQL = `CREATE TABLE webauthn_credentials (
+    name            TEXT PRIMARY KEY CHECK (length(name) BETWEEN 1 AND 48),
+    credential_id   BLOB NOT NULL UNIQUE CHECK (length(credential_id) BETWEEN 1 AND 1023),
+    rp_id           TEXT NOT NULL CHECK (rp_id <> ''),
+    origin          TEXT NOT NULL CHECK (origin <> ''),
+    alg             INTEGER NOT NULL CHECK (alg IN (-7, -8)),
+    aaguid          BLOB NOT NULL CHECK (length(aaguid) = 16),
+    backup_eligible INTEGER NOT NULL CHECK (backup_eligible IN (0, 1)),
+    sign_count      INTEGER NOT NULL CHECK (sign_count BETWEEN 0 AND 4294967295)
+)`
````

### Patch M4b — `EnrolPasskey`

````diff M4b
--- a/host/store/credentials_test.go
+++ b/host/store/credentials_test.go
@@ -293,6 +293,8 @@
 	_, _, errs["ApproverByKey"] = c.ApproverByKey(ctx, []byte("k"))
 	_, errs["ListApprovers"] = c.ListApprovers(ctx)
 	errs["RevokeApprover"] = c.RevokeApprover(ctx, "x", 1000)
+	a, p := testPasskey("x", 1)
+	errs["EnrolPasskey"] = c.EnrolPasskey(ctx, a, p)
 	for name, err := range errs {
 		if !errors.Is(err, ErrNoDeadline) {
 			t.Errorf("%s(Background) = %v, want ErrNoDeadline", name, err)
--- a/host/store/passkeys.go
+++ b/host/store/passkeys.go
@@ -7,6 +7,13 @@
 // forward, inside the consume that uses it. The store does not import host/approval: these checks
 // are the store's own backstop for the pure law the CLI runs first.
 
+import (
+	"context"
+	"database/sql"
+	"errors"
+	"fmt"
+)
+
 // ---- M4a ----
 
 const passkeysTable = "webauthn_credentials"
@@ -26,3 +33,107 @@
     backup_eligible INTEGER NOT NULL CHECK (backup_eligible IN (0, 1)),
     sign_count      INTEGER NOT NULL CHECK (sign_count BETWEEN 0 AND 4294967295)
 )`
+
+// ---- M4b ----
+
+// maxSignCount is the largest WebAuthn sign counter (32 bits); maxActivePasskeys caps ActivePasskeys.
+const (
+	maxSignCount      = 1<<32 - 1
+	maxActivePasskeys = 64
+)
+
+// PasskeyRow is one webauthn_credentials row.
+type PasskeyRow struct {
+	Name, RPID, Origin   string
+	CredentialID, AAGUID []byte
+	Alg                  int64 // -7 (ES256) or -8 (EdDSA)
+	BackupEligible       bool
+	SignCount            int64 // 0 .. 2^32-1
+}
+
+var (
+	ErrPasskeyCredentialEnrolled = errors.New("store: that passkey credential is already enrolled (a credential is never enrolled twice, even after revoke)")
+	ErrSignCountRegression       = errors.New("store: the passkey's sign counter did not advance (a cloned authenticator, or an older assertion)")
+	ErrPasskeyChanged            = errors.New("store: the approver's passkey is not the one that signed")
+)
+
+func (c *Credentials) olderForPasskeys() error {
+	return &OlderCredentialsStoreError{Path: c.display, Version: c.version, Need: passkeysSchemaVersion, Feature: "passkeys"}
+}
+
+// checkPasskeyRows is EnrolPasskey's validation (S1-S10): each refusal is ErrBadApprover.
+func checkPasskeyRows(a ApproverRow, p PasskeyRow) error {
+	switch {
+	case a.Format != "webauthn": // S1
+		return fmt.Errorf("%w: format %q is not webauthn", ErrBadApprover, a.Format)
+	case a.Name != p.Name: // S2
+		return fmt.Errorf("%w: the approver %q and its passkey %q differ", ErrBadApprover, a.Name, p.Name)
+	case a.Name == "" || len(a.PublicKey) == 0: // S3
+		return fmt.Errorf("%w: empty name or public key", ErrBadApprover)
+	case a.EnrolledAt <= 0: // S4
+		return fmt.Errorf("%w: enrolled_at %d is not positive", ErrBadApprover, a.EnrolledAt)
+	case a.LastUsedAt != 0 || a.RevokedAt != 0: // S5
+		return fmt.Errorf("%w: a new approver is neither used nor revoked", ErrBadApprover)
+	case len(p.CredentialID) < 1 || len(p.CredentialID) > 1023: // S6
+		return fmt.Errorf("%w: a credential ID is 1 to 1023 bytes, not %d", ErrBadApprover, len(p.CredentialID))
+	case p.Alg != -7 && p.Alg != -8: // S7
+		return fmt.Errorf("%w: alg %d is not -7 or -8", ErrBadApprover, p.Alg)
+	case len(p.AAGUID) != 16: // S8
+		return fmt.Errorf("%w: an AAGUID is 16 bytes, not %d", ErrBadApprover, len(p.AAGUID))
+	case p.RPID == "" || p.Origin == "": // S9
+		return fmt.Errorf("%w: empty relying party ID or origin", ErrBadApprover)
+	case p.SignCount < 0 || p.SignCount > maxSignCount: // S10
+		return fmt.Errorf("%w: sign count %d is not 0 to 2^32-1", ErrBadApprover, p.SignCount)
+	}
+	return nil
+}
+
+// EnrolPasskey stores one new, active webauthn approver and its passkey in one transaction. It
+// refuses a name, a key or a credential ID that any row holds, revoked rows included, naming the
+// holder; the UNIQUE constraints are the backstop and their error text is never parsed.
+func (c *Credentials) EnrolPasskey(ctx context.Context, a ApproverRow, p PasskeyRow) error {
+	if err := requireDeadline(ctx); err != nil {
+		return err
+	}
+	if c.version < passkeysSchemaVersion {
+		return c.olderForPasskeys()
+	}
+	if err := checkPasskeyRows(a, p); err != nil {
+		return err
+	}
+	tx, err := c.db.BeginTx(ctx, nil)
+	if err != nil {
+		return fmt.Errorf("store: enrol passkey: begin: %w", err)
+	}
+	defer func() { _ = tx.Rollback() }()
+	for _, pre := range []struct {
+		query string
+		arg   any
+		err   error
+	}{
+		{`SELECT name FROM approvers WHERE name = ?`, a.Name, ErrApproverNameTaken},
+		{`SELECT name FROM approvers WHERE public_key = ?`, a.PublicKey, ErrApproverKeyEnrolled},
+		{`SELECT name FROM webauthn_credentials WHERE credential_id = ?`, p.CredentialID, ErrPasskeyCredentialEnrolled},
+	} {
+		var holder string
+		err := tx.QueryRowContext(ctx, pre.query, pre.arg).Scan(&holder)
+		if err == nil {
+			return fmt.Errorf("%w: approver %q holds it", pre.err, holder)
+		}
+		if !errors.Is(err, sql.ErrNoRows) {
+			return fmt.Errorf("store: enrol passkey: %w", err)
+		}
+	}
+	if _, err := tx.ExecContext(ctx, `INSERT INTO approvers (name, format, public_key, enrolled_at) VALUES (?, ?, ?, ?)`,
+		a.Name, a.Format, a.PublicKey, a.EnrolledAt); err != nil {
+		return fmt.Errorf("store: enrol passkey: %w", err)
+	}
+	if _, err := tx.ExecContext(ctx, `INSERT INTO webauthn_credentials (name, credential_id, rp_id, origin, alg, aaguid, backup_eligible, sign_count) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
+		p.Name, p.CredentialID, p.RPID, p.Origin, p.Alg, p.AAGUID, p.BackupEligible, p.SignCount); err != nil {
+		return fmt.Errorf("store: enrol passkey: %w", err)
+	}
+	if err := tx.Commit(); err != nil {
+		return fmt.Errorf("store: enrol passkey: commit: %w", err)
+	}
+	return nil
+}
--- a/host/store/passkeys_test.go
+++ b/host/store/passkeys_test.go
@@ -0,0 +1,160 @@
+package store
+
+// Row 201 A2 sprint 2 M4b (design §6, PS2, PS5): schema 6's passkey rows, enrolled only with their
+// approver, and never twice.
+
+import (
+	"bytes"
+	"errors"
+	"fmt"
+	"path/filepath"
+	"strings"
+	"testing"
+)
+
+// testPasskey is a valid new webauthn approver and passkey: name, COSE key "cose-"+name, enrolled
+// at 1000; credential ID 32 bytes of id, rp localhost, ES256, counter 1.
+func testPasskey(name string, id byte) (ApproverRow, PasskeyRow) {
+	return ApproverRow{Name: name, Format: "webauthn", PublicKey: []byte("cose-" + name), EnrolledAt: 1000},
+		PasskeyRow{Name: name, CredentialID: bytes.Repeat([]byte{id}, 32), RPID: "localhost", Origin: "http://localhost:7644",
+			Alg: -7, AAGUID: make([]byte, 16), SignCount: 1}
+}
+
+// rowCounts is the number of approvers and webauthn_credentials rows in c's store.
+func rowCounts(t *testing.T, c *Credentials) string {
+	t.Helper()
+	var a, p int
+	if err := c.db.QueryRow(`SELECT (SELECT count(*) FROM approvers), (SELECT count(*) FROM webauthn_credentials)`).Scan(&a, &p); err != nil {
+		t.Fatal(err)
+	}
+	return fmt.Sprintf("%d approvers, %d passkeys", a, p)
+}
+
+// openV5 opens a schema-5 store (not upgraded).
+func openV5(t *testing.T) *Credentials {
+	t.Helper()
+	path := filepath.Join(t.TempDir(), CredentialsFileName)
+	makeV5Credentials(t, path, nil)
+	c, err := OpenCredentials(path)
+	if err != nil {
+		t.Fatal(err)
+	}
+	t.Cleanup(func() { _ = c.Close() })
+	return c
+}
+
+// isOlderPasskeys reports whether err is the passkeys refusal of a schema-5 store.
+func isOlderPasskeys(err error) bool {
+	var oe *OlderCredentialsStoreError
+	return errors.As(err, &oe) && oe.Version == 5 && oe.Need == 6 && oe.Feature == "passkeys" && strings.Contains(err.Error(), "passkeys need schema 6")
+}
+
+// PS2: every EnrolPasskey refusal returns its typed error and writes nothing to either table.
+func TestEnrolPasskeyRefusals(t *testing.T) {
+	ctx := testCtx(t)
+	a, p := testPasskey("mark", 1)
+	c5 := openV5(t)
+	if err := c5.EnrolPasskey(ctx, a, p); !isOlderPasskeys(err) {
+		t.Fatalf("v5 enrol: %v, want the passkeys refusal", err)
+	}
+	if n := approverCount(t, c5); n != 0 {
+		t.Fatalf("v5 enrol wrote %d approvers", n)
+	}
+
+	c := openTestCredentials(t)
+	bad := map[string]func(a *ApproverRow, p *PasskeyRow){
+		"S1-format-ssh":    func(a *ApproverRow, p *PasskeyRow) { a.Format = "ssh" },
+		"S1-format-empty":  func(a *ApproverRow, p *PasskeyRow) { a.Format = "" },
+		"S2-names-differ":  func(a *ApproverRow, p *PasskeyRow) { p.Name = "other" },
+		"S3-empty-name":    func(a *ApproverRow, p *PasskeyRow) { a.Name, p.Name = "", "" },
+		"S3-nil-key":       func(a *ApproverRow, p *PasskeyRow) { a.PublicKey = nil },
+		"S4-enrolled-0":    func(a *ApproverRow, p *PasskeyRow) { a.EnrolledAt = 0 },
+		"S4-enrolled--1":   func(a *ApproverRow, p *PasskeyRow) { a.EnrolledAt = -1 },
+		"S5-last-used":     func(a *ApproverRow, p *PasskeyRow) { a.LastUsedAt = 1500 },
+		"S5-revoked":       func(a *ApproverRow, p *PasskeyRow) { a.RevokedAt = 1500 },
+		"S6-cred-0":        func(a *ApproverRow, p *PasskeyRow) { p.CredentialID = nil },
+		"S6-cred-1024":     func(a *ApproverRow, p *PasskeyRow) { p.CredentialID = make([]byte, 1024) },
+		"S7-alg-rs256":     func(a *ApproverRow, p *PasskeyRow) { p.Alg = -257 },
+		"S7-alg-0":         func(a *ApproverRow, p *PasskeyRow) { p.Alg = 0 },
+		"S8-aaguid-15":     func(a *ApproverRow, p *PasskeyRow) { p.AAGUID = make([]byte, 15) },
+		"S8-aaguid-17":     func(a *ApproverRow, p *PasskeyRow) { p.AAGUID = make([]byte, 17) },
+		"S9-rp-empty":      func(a *ApproverRow, p *PasskeyRow) { p.RPID = "" },
+		"S9-origin-empty":  func(a *ApproverRow, p *PasskeyRow) { p.Origin = "" },
+		"S10-count--1":     func(a *ApproverRow, p *PasskeyRow) { p.SignCount = -1 },
+		"S10-count-2^32":   func(a *ApproverRow, p *PasskeyRow) { p.SignCount = 1 << 32 },
+		"S10-count-2^63-1": func(a *ApproverRow, p *PasskeyRow) { p.SignCount = 1<<63 - 1 },
+	}
+	for name, edit := range bad {
+		a, p := testPasskey("bad", 2)
+		edit(&a, &p)
+		if err := c.EnrolPasskey(ctx, a, p); !errors.Is(err, ErrBadApprover) {
+			t.Fatalf("%s: err %v, want ErrBadApprover", name, err)
+		}
+		if got := rowCounts(t, c); got != "0 approvers, 0 passkeys" {
+			t.Fatalf("%s wrote: %s", name, got)
+		}
+	}
+	// The bounds are inclusive: each of these enrols.
+	for i, edit := range []func(a *ApproverRow, p *PasskeyRow){
+		func(a *ApproverRow, p *PasskeyRow) { p.CredentialID = []byte{0xee} },
+		func(a *ApproverRow, p *PasskeyRow) { p.CredentialID = bytes.Repeat([]byte{0xef}, 1023) },
+		func(a *ApproverRow, p *PasskeyRow) { p.SignCount, p.Alg = 0, -8 },
+		func(a *ApproverRow, p *PasskeyRow) { p.SignCount = 1<<32 - 1 },
+	} {
+		a, p := testPasskey(fmt.Sprint("edge", i), byte(0xf0+i))
+		edit(&a, &p)
+		if err := c.EnrolPasskey(ctx, a, p); err != nil {
+			t.Fatalf("edge %d: %v", i, err)
+		}
+	}
+
+	if err := c.EnrolPasskey(ctx, a, p); err != nil {
+		t.Fatalf("good enrol: %v", err)
+	}
+	if err := c.EnrolApprover(ctx, testApprover("ssh-one", 1000)); err != nil {
+		t.Fatal(err)
+	}
+	const after = "6 approvers, 5 passkeys"
+	if got := rowCounts(t, c); got != after {
+		t.Fatalf("after the good enrols: %s, want %s", got, after)
+	}
+	check := func(when string) {
+		t.Helper()
+		for _, r := range []struct {
+			name   string
+			edit   func(a *ApproverRow, p *PasskeyRow)
+			want   error
+			holder string
+		}{
+			{"name-taken", func(a *ApproverRow, p *PasskeyRow) { a.Name, p.Name = "mark", "mark" }, ErrApproverNameTaken, "mark"},
+			{"name-taken-by-ssh", func(a *ApproverRow, p *PasskeyRow) { a.Name, p.Name = "ssh-one", "ssh-one" }, ErrApproverNameTaken, "ssh-one"},
+			{"key-enrolled", func(a *ApproverRow, p *PasskeyRow) { a.PublicKey = []byte("cose-mark") }, ErrApproverKeyEnrolled, "mark"},
+			{"key-enrolled-by-ssh", func(a *ApproverRow, p *PasskeyRow) { a.PublicKey = []byte("key-ssh-one") }, ErrApproverKeyEnrolled, "ssh-one"},
+			{"credential-enrolled", func(a *ApproverRow, p *PasskeyRow) { p.CredentialID = bytes.Repeat([]byte{1}, 32) }, ErrPasskeyCredentialEnrolled, "mark"},
+		} {
+			a, p := testPasskey("new", 3)
+			r.edit(&a, &p)
+			err := c.EnrolPasskey(ctx, a, p)
+			if !errors.Is(err, r.want) || !strings.Contains(err.Error(), `"`+r.holder+`"`) {
+				t.Fatalf("%s, %s: err %v, want %v naming %q", when, r.name, err, r.want, r.holder)
+			}
+			if got := rowCounts(t, c); got != after {
+				t.Fatalf("%s, %s wrote: %s", when, r.name, got)
+			}
+		}
+	}
+	check("active")
+	if _, err := c.db.Exec(`UPDATE approvers SET revoked_at = 1200 WHERE name IN ('mark', 'ssh-one')`); err != nil {
+		t.Fatal(err)
+	}
+	check("revoked")
+}
+
+// PS5: EnrolApprover still refuses every format but ssh, so a passkey enters only with its row.
+func TestEnrolApproverStillSSHOnly(t *testing.T) {
+	c := openTestCredentials(t)
+	a, _ := testPasskey("mark", 1)
+	if err := c.EnrolApprover(testCtx(t), a); !errors.Is(err, ErrBadApprover) || rowCounts(t, c) != "0 approvers, 0 passkeys" {
+		t.Fatalf("EnrolApprover(webauthn): %v, %s", err, rowCounts(t, c))
+	}
+}
````

### Patch M4c — the passkey lookups

````diff M4c
--- a/host/store/credentials_test.go
+++ b/host/store/credentials_test.go
@@ -295,6 +295,9 @@
 	errs["RevokeApprover"] = c.RevokeApprover(ctx, "x", 1000)
 	a, p := testPasskey("x", 1)
 	errs["EnrolPasskey"] = c.EnrolPasskey(ctx, a, p)
+	_, _, _, errs["PasskeyByCredentialID"] = c.PasskeyByCredentialID(ctx, p.CredentialID)
+	_, _, errs["Passkey"] = c.Passkey(ctx, "x")
+	_, errs["ActivePasskeys"] = c.ActivePasskeys(ctx)
 	for name, err := range errs {
 		if !errors.Is(err, ErrNoDeadline) {
 			t.Errorf("%s(Background) = %v, want ErrNoDeadline", name, err)
--- a/host/store/passkey_lookups_test.go
+++ b/host/store/passkey_lookups_test.go
@@ -0,0 +1,164 @@
+package store
+
+// Row 201 A2 sprint 2 M4c (design §6, PS3, RC6): the passkey lookups read back exactly what was
+// enrolled, and miss without an error.
+
+import (
+	"bytes"
+	"database/sql"
+	"fmt"
+	"reflect"
+	"sort"
+	"strings"
+	"testing"
+)
+
+// PS3: the three lookups return the enrolled rows exactly, miss without an error, and refuse a v5 store.
+func TestPasskeyLookups(t *testing.T) {
+	ctx := testCtx(t)
+	c := openTestCredentials(t)
+	for i, name := range []string{"c", "a", "b"} {
+		a, p := testPasskey(name, byte(10+i))
+		if err := c.EnrolPasskey(ctx, a, p); err != nil {
+			t.Fatal(err)
+		}
+	}
+	if err := c.EnrolApprover(ctx, testApprover("s", 1000)); err != nil {
+		t.Fatal(err)
+	}
+	if err := c.RevokeApprover(ctx, "b", 1100); err != nil {
+		t.Fatal(err)
+	}
+	if _, err := c.db.Exec(`UPDATE approvers SET last_used_at = 1050 WHERE name = 'a'`); err != nil {
+		t.Fatal(err)
+	}
+	wantA, wantP := testPasskey("a", 11)
+	wantA.LastUsedAt = 1050
+	if b, _, ok, err := c.PasskeyByCredentialID(ctx, bytes.Repeat([]byte{12}, 32)); err != nil || !ok || b.Name != "b" || b.RevokedAt != 1100 {
+		t.Fatalf("PasskeyByCredentialID(b's), revoked: %+v ok %v err %v", b, ok, err)
+	}
+	gotA, gotP, ok, err := c.PasskeyByCredentialID(ctx, wantP.CredentialID)
+	if err != nil || !ok || !reflect.DeepEqual(gotA, wantA) || !reflect.DeepEqual(gotP, wantP) {
+		t.Fatalf("PasskeyByCredentialID(a's): %+v %+v ok %v err %v", gotA, gotP, ok, err)
+	}
+	for name, id := range map[string][]byte{"prefix": wantP.CredentialID[:31], "longer": append(bytes.Clone(wantP.CredentialID), 0), "none": nil} {
+		if a, p, ok, err := c.PasskeyByCredentialID(ctx, id); ok || err != nil || a.Name != "" || p.Name != "" {
+			t.Fatalf("PasskeyByCredentialID(%s): %+v %+v ok %v err %v; want a miss", name, a, p, ok, err)
+		}
+	}
+	if p, ok, err := c.Passkey(ctx, "b"); err != nil || !ok || p.Name != "b" || p.CredentialID[0] != 12 {
+		t.Fatalf("Passkey(b), revoked: %+v ok %v err %v", p, ok, err)
+	}
+	for _, name := range []string{"s", "", "a%", "zz"} {
+		if p, ok, err := c.Passkey(ctx, name); ok || err != nil || p.Name != "" {
+			t.Fatalf("Passkey(%q): %+v ok %v err %v; want a miss", name, p, ok, err)
+		}
+	}
+	names := func() string {
+		t.Helper()
+		got, err := c.ActivePasskeys(ctx)
+		if err != nil {
+			t.Fatal(err)
+		}
+		var n []string
+		for _, e := range got {
+			if e.Approver.Name != e.Passkey.Name || e.Approver.Format != "webauthn" {
+				t.Fatalf("ActivePasskeys pair %+v", e)
+			}
+			n = append(n, e.Passkey.Name)
+		}
+		return strings.Join(n, ",")
+	}
+	if got := names(); got != "a,c" {
+		t.Fatalf("ActivePasskeys = %s, want a,c (b revoked, s ssh)", got)
+	}
+	// A passkey row beside an ssh approver (made by hand) is not an active passkey.
+	if _, err := c.db.Exec(`INSERT INTO webauthn_credentials (name, credential_id, rp_id, origin, alg, aaguid, backup_eligible, sign_count)
+		VALUES ('s', x'99', 'localhost', 'http://localhost', -7, zeroblob(16), 0, 0)`); err != nil {
+		t.Fatal(err)
+	}
+	if got := names(); got != "a,c" {
+		t.Fatalf("ActivePasskeys with an ssh row's passkey = %s, want a,c", got)
+	}
+	for i := 0; i < 70; i++ {
+		a, p := testPasskey(fmt.Sprintf("m%02d", i), 0)
+		p.CredentialID = []byte(a.Name)
+		if err := c.EnrolPasskey(ctx, a, p); err != nil {
+			t.Fatal(err)
+		}
+	}
+	if got := strings.Split(names(), ","); len(got) != 64 || got[0] != "a" || got[1] != "c" || got[63] != "m61" || !sort.StringsAreSorted(got) {
+		t.Fatalf("ActivePasskeys over 72 active: %d names, first %v, last %v; want 64 by name", len(got), got[:2], got[len(got)-1])
+	}
+
+	c5 := openV5(t)
+	_, _, _, e1 := c5.PasskeyByCredentialID(ctx, wantP.CredentialID)
+	_, _, e2 := c5.Passkey(ctx, "a")
+	_, e3 := c5.ActivePasskeys(ctx)
+	for i, err := range []error{e1, e2, e3} {
+		if !isOlderPasskeys(err) {
+			t.Fatalf("v5 lookup %d: %v, want the passkeys refusal", i, err)
+		}
+	}
+}
+
+// passkeyRowColumns maps every PasskeyRow field to its webauthn_credentials column.
+var passkeyRowColumns = map[string]string{
+	"Name": "name", "CredentialID": "credential_id", "RPID": "rp_id", "Origin": "origin", "Alg": "alg",
+	"AAGUID": "aaguid", "BackupEligible": "backup_eligible", "SignCount": "sign_count",
+}
+
+// RC6: PasskeyRow, the table and the rows the three lookups scan agree field for field, each a
+// distinct non-zero value.
+func TestPasskeyRowMatchesTheTable(t *testing.T) {
+	ctx := testCtx(t)
+	c := openTestCredentials(t)
+	typ := reflect.TypeOf(PasskeyRow{})
+	var fields, keys, cols, want []string
+	for i := 0; i < typ.NumField(); i++ {
+		fields = append(fields, typ.Field(i).Name)
+	}
+	for k, v := range passkeyRowColumns {
+		keys, want = append(keys, k), append(want, v)
+	}
+	rows, err := c.db.Query(`PRAGMA table_info(webauthn_credentials)`)
+	if err != nil {
+		t.Fatal(err)
+	}
+	for rows.Next() {
+		var cid, notnull, pk int
+		var name, ctype string
+		var dflt sql.NullString
+		if err := rows.Scan(&cid, &name, &ctype, &notnull, &dflt, &pk); err != nil {
+			t.Fatal(err)
+		}
+		cols = append(cols, name)
+	}
+	_ = rows.Close()
+	for _, s := range [][]string{fields, keys, cols, want} {
+		sort.Strings(s)
+	}
+	if !reflect.DeepEqual(fields, keys) || !reflect.DeepEqual(cols, want) {
+		t.Fatalf("PasskeyRow fields %v vs map %v; columns %v vs map %v", fields, keys, cols, want)
+	}
+	a := ApproverRow{Name: "rc6", Format: "webauthn", PublicKey: []byte("rc6-key"), EnrolledAt: 1000}
+	p := PasskeyRow{Name: "rc6", CredentialID: []byte{9, 8, 7}, RPID: "rp.ts.net", Origin: "https://rp.ts.net", Alg: -8,
+		AAGUID: bytes.Repeat([]byte{5}, 16), BackupEligible: true, SignCount: 4294967290}
+	if err := c.EnrolPasskey(ctx, a, p); err != nil {
+		t.Fatal(err)
+	}
+	_, byID, _, err1 := c.PasskeyByCredentialID(ctx, p.CredentialID)
+	byName, _, err2 := c.Passkey(ctx, "rc6")
+	active, err3 := c.ActivePasskeys(ctx)
+	if err1 != nil || err2 != nil || err3 != nil || len(active) != 1 {
+		t.Fatalf("lookups: %v %v %v %d", err1, err2, err3, len(active))
+	}
+	for _, got := range []PasskeyRow{byID, byName, active[0].Passkey} {
+		gv, wv := reflect.ValueOf(got), reflect.ValueOf(p)
+		for i := 0; i < typ.NumField(); i++ {
+			if gv.Field(i).IsZero() || !reflect.DeepEqual(gv.Field(i).Interface(), wv.Field(i).Interface()) {
+				t.Fatalf("field %s = %v, want %v", typ.Field(i).Name, gv.Field(i).Interface(), wv.Field(i).Interface())
+			}
+		}
+	}
+}
--- a/host/store/passkeys.go
+++ b/host/store/passkeys.go
@@ -137,3 +137,74 @@
 	}
 	return nil
 }
+
+// ---- M4c ----
+
+// EnrolledPasskey is an approvers row with its passkey row.
+type EnrolledPasskey struct {
+	Approver ApproverRow
+	Passkey  PasskeyRow
+}
+
+const passkeyColumns = `w.name, w.credential_id, w.rp_id, w.origin, w.alg, w.aaguid, w.backup_eligible, w.sign_count`
+
+func scanPasskey(scan func(...any) error, dst ...any) (PasskeyRow, error) {
+	var p PasskeyRow
+	err := scan(append(dst, &p.Name, &p.CredentialID, &p.RPID, &p.Origin, &p.Alg, &p.AAGUID, &p.BackupEligible, &p.SignCount)...)
+	return p, err
+}
+
+// enrolledPasskeys reads approvers joined with their passkeys, where says which.
+func (c *Credentials) enrolledPasskeys(ctx context.Context, where string, args ...any) ([]EnrolledPasskey, error) {
+	if err := requireDeadline(ctx); err != nil {
+		return nil, err
+	}
+	if c.version < passkeysSchemaVersion {
+		return nil, c.olderForPasskeys()
+	}
+	rows, err := c.db.QueryContext(ctx, `SELECT a.name, a.format, a.public_key, a.enrolled_at, a.last_used_at, a.revoked_at, `+passkeyColumns+
+		` FROM approvers a JOIN webauthn_credentials w ON w.name = a.name WHERE `+where, args...)
+	if err != nil {
+		return nil, fmt.Errorf("store: read passkeys: %w", err)
+	}
+	defer func() { _ = rows.Close() }()
+	var out []EnrolledPasskey
+	for rows.Next() {
+		var e EnrolledPasskey
+		var lastUsed, revoked sql.NullInt64
+		a := &e.Approver
+		if e.Passkey, err = scanPasskey(rows.Scan, &a.Name, &a.Format, &a.PublicKey, &a.EnrolledAt, &lastUsed, &revoked); err != nil {
+			return nil, fmt.Errorf("store: read passkeys: %w", err)
+		}
+		a.LastUsedAt, a.RevokedAt = lastUsed.Int64, revoked.Int64
+		out = append(out, e)
+	}
+	if err := rows.Err(); err != nil {
+		return nil, fmt.Errorf("store: read passkeys: %w", err)
+	}
+	return out, nil
+}
+
+// PasskeyByCredentialID reads the approver whose passkey has credential ID id, byte for byte, and
+// that passkey. ok is false when none.
+func (c *Credentials) PasskeyByCredentialID(ctx context.Context, id []byte) (ApproverRow, PasskeyRow, bool, error) {
+	got, err := c.enrolledPasskeys(ctx, `w.credential_id = ?`, id)
+	if err != nil || len(got) == 0 {
+		return ApproverRow{}, PasskeyRow{}, false, err
+	}
+	return got[0].Approver, got[0].Passkey, true, nil
+}
+
+// Passkey reads approver name's passkey, revoked approvers included. ok is false when none.
+func (c *Credentials) Passkey(ctx context.Context, name string) (PasskeyRow, bool, error) {
+	got, err := c.enrolledPasskeys(ctx, `w.name = ?`, name)
+	if err != nil || len(got) == 0 {
+		return PasskeyRow{}, false, err
+	}
+	return got[0].Passkey, true, nil
+}
+
+// ActivePasskeys lists the active webauthn approvers with their passkeys, ordered by name, at most 64.
+func (c *Credentials) ActivePasskeys(ctx context.Context) ([]EnrolledPasskey, error) {
+	return c.enrolledPasskeys(ctx, `a.revoked_at IS NULL AND a.format = 'webauthn' ORDER BY a.name LIMIT ?`, maxActivePasskeys)
+}
````

### Patch M5a — the consume binds the passkey (C13–C18)

````diff M5a
--- a/host/store/consume_binding_test.go
+++ b/host/store/consume_binding_test.go
@@ -190,13 +190,16 @@
 	"ApproverUse.Name":          "required (C2)",
 	"ApproverUse.Format":        "required (C2)",
 	"ApproverUse.PublicKey":     "required (C2)",
+	"ApproverUse.Passkey":       "optional with a meaning: set exactly for a webauthn approver (C13), bound by C14-C18",
+	"PasskeyUse.CredentialID":   "required (C14): 1 to 1023 bytes",
+	"PasskeyUse.SignCount":      "bounded (C14): 0 to 2^32-1; 0 is a counter (TestPasskeyConsumeCompleteness)",
 }
 
 // RC3 (M4d1): every RequestConsume and ApproverUse field has a disposition, and zeroing each
 // required one in an otherwise valid signed consume is ErrIncompleteConsume, before the clock.
 func TestRequestConsumeCompleteness(t *testing.T) {
 	n := 0
-	for _, typ := range []reflect.Type{reflect.TypeOf(RequestConsume{}), reflect.TypeOf(ApproverUse{})} {
+	for _, typ := range []reflect.Type{reflect.TypeOf(RequestConsume{}), reflect.TypeOf(ApproverUse{}), reflect.TypeOf(PasskeyUse{})} {
 		for i := 0; i < typ.NumField(); i++ {
 			if key := typ.Name() + "." + typ.Field(i).Name; requestConsumeFields[key] == "" {
 				t.Errorf("%s has no disposition in requestConsumeFields", key)
@@ -205,7 +208,7 @@
 		n += typ.NumField()
 	}
 	if n != len(requestConsumeFields) {
-		t.Fatalf("RequestConsume and ApproverUse have %d fields, requestConsumeFields %d", n, len(requestConsumeFields))
+		t.Fatalf("RequestConsume, ApproverUse and PasskeyUse have %d fields, requestConsumeFields %d", n, len(requestConsumeFields))
 	}
 	zero := map[string]func(*RequestConsume){
 		"RequestConsume.ID":        func(rc *RequestConsume) { rc.ID = "" },
@@ -216,6 +219,9 @@
 		"ApproverUse.Name":         func(rc *RequestConsume) { rc.Approver.Name, rc.ApprovedBy = "", "ssh:" },
 		"ApproverUse.Format":       func(rc *RequestConsume) { rc.Approver.Format, rc.ApprovedBy = "", ":x" },
 		"ApproverUse.PublicKey":    func(rc *RequestConsume) { rc.Approver.PublicKey = nil },
+		"PasskeyUse.CredentialID": func(rc *RequestConsume) {
+			rc.Approver.Format, rc.ApprovedBy, rc.Approver.Passkey = "webauthn", "webauthn:x", &PasskeyUse{SignCount: 2}
+		},
 	}
 	for key, d := range requestConsumeFields {
 		if _, ok := zero[key]; strings.HasPrefix(d, "required") != ok {
@@ -287,6 +293,7 @@
 		}, ErrApproverKeyChanged},
 		{"other-format", func(t *testing.T, c *Credentials, path string, rc *RequestConsume) {
 			rc.Approver.Format, rc.ApprovedBy = "webauthn", "webauthn:x"
+			rc.Approver.Passkey = &PasskeyUse{CredentialID: make([]byte, 32), SignCount: 1} // past C13-C15 (row 201 A2 §2.10)
 		}, ErrApproverKeyChanged},
 		{"enrolled-after-the-reading", func(t *testing.T, c *Credentials, path string, rc *RequestConsume) {
 			if _, err := rawDB(t, path).Exec(`UPDATE approvers SET enrolled_at = 1500 WHERE name = 'x'`); err != nil {
--- a/host/store/mint_requests.go
+++ b/host/store/mint_requests.go
@@ -294,6 +294,7 @@
 type ApproverUse struct {
 	Name, Format string
 	PublicKey    []byte
+	Passkey      *PasskeyUse // non-nil exactly when Format is "webauthn" (row 201 A2 C13)
 }
 
 // MintSessionForRequest consumes request rc and inserts session in one transaction: the request
@@ -323,7 +324,18 @@
 	}
 	if rc.Approver != nil && c.version < approversSchemaVersion { // C5
 		return c.olderForApprovers()
+	}
+	if a := rc.Approver; a != nil && (a.Format == "webauthn") != (a.Passkey != nil) { // C13
+		return ErrIncompleteConsume
 	}
+	if a := rc.Approver; a != nil && a.Passkey != nil {
+		if p := a.Passkey; len(p.CredentialID) < 1 || len(p.CredentialID) > 1023 || p.SignCount < 0 || p.SignCount > maxSignCount { // C14
+			return ErrIncompleteConsume
+		}
+		if c.version < passkeysSchemaVersion { // C15
+			return c.olderForPasskeys()
+		}
+	}
 	tx, err := c.db.BeginTx(ctx, nil)
 	if err != nil {
 		return fmt.Errorf("store: consume mint request: begin: %w", err)
@@ -349,6 +361,11 @@
 		if err := useApprover(ctx, tx, *rc.Approver, at); err != nil {
 			return err
 		}
+		if p := rc.Approver.Passkey; p != nil { // C16-C18: and its counter, forward only
+			if err := usePasskey(ctx, tx, rc.Approver.Name, *p); err != nil {
+				return err
+			}
+		}
 	}
 	if err := insertSessionRow(ctx, tx, session); err != nil {
 		return err
--- a/host/store/passkey_consume_test.go
+++ b/host/store/passkey_consume_test.go
@@ -0,0 +1,173 @@
+package store
+
+// Row 201 A2 sprint 2 M5 (design §6, CS7-CS9): a passkey consume binds the credential that signed
+// and moves its counter forward, in the transaction that consumes the request, or nothing is written.
+
+import (
+	"bytes"
+	"errors"
+	"path/filepath"
+	"testing"
+)
+
+// passkeyAt is a fresh schema-6 store, its path, and webauthn approver x enrolled at 900 with
+// credential cred and counter stored.
+func passkeyAt(t *testing.T, cred []byte, stored int64) (*Credentials, string) {
+	t.Helper()
+	path := filepath.Join(t.TempDir(), CredentialsFileName)
+	c, err := OpenCredentials(path)
+	if err != nil {
+		t.Fatal(err)
+	}
+	t.Cleanup(func() { _ = c.Close() })
+	a, p := testPasskey("x", 1)
+	a.EnrolledAt, p.CredentialID, p.SignCount = 900, cred, stored
+	if err := c.EnrolPasskey(testCtx(t), a, p); err != nil {
+		t.Fatal(err)
+	}
+	return c, path
+}
+
+var credX = bytes.Repeat([]byte{1}, 32)
+
+// passkeyConsume is okConsume(r, at, ...) for x's passkey assertion with counter count.
+func passkeyConsume(r MintRequestRow, at int64, clocks, checks *int, cred []byte, count int64) RequestConsume {
+	rc := okConsume(r, at, clocks, checks)
+	rc.Approver = &ApproverUse{Name: "x", Format: "webauthn", PublicKey: []byte("cose-x"), Passkey: &PasskeyUse{CredentialID: cred, SignCount: count}}
+	rc.ApprovedBy = "webauthn:x"
+	return rc
+}
+
+// passkeyState is x's counter and last_used_at.
+func passkeyState(t *testing.T, c *Credentials) (int64, int64) {
+	t.Helper()
+	p, ok, err := c.Passkey(testCtx(t), "x")
+	if err != nil || !ok {
+		t.Fatalf("x's passkey: ok %v err %v", ok, err)
+	}
+	return p.SignCount, approverX(t, c).LastUsedAt
+}
+
+// CS7: C13 (the passkey is set exactly for webauthn), C14 (its bounds) and C15 (schema 6), each
+// before the clock is read, leaving the request pending, no session, and x untouched.
+func TestPasskeyConsumeCompleteness(t *testing.T) {
+	for name, edit := range map[string]func(*RequestConsume){
+		"C13-webauthn-without-passkey": func(rc *RequestConsume) { rc.Approver.Passkey = nil },
+		"C13-ssh-with-passkey":         func(rc *RequestConsume) { rc.Approver.Format, rc.ApprovedBy = "ssh", "ssh:x" },
+		"C14-empty-credential":         func(rc *RequestConsume) { rc.Approver.Passkey.CredentialID = nil },
+		"C14-1024-byte-credential":     func(rc *RequestConsume) { rc.Approver.Passkey.CredentialID = make([]byte, 1024) },
+		"C14-counter-2^32":             func(rc *RequestConsume) { rc.Approver.Passkey.SignCount = 1 << 32 },
+		"C14-counter--1":               func(rc *RequestConsume) { rc.Approver.Passkey.SignCount = -1 },
+	} {
+		t.Run(name, func(t *testing.T) {
+			c, _ := passkeyAt(t, credX, 1)
+			r := putTestRequest(t, c)
+			var clocks, checks int
+			rc := passkeyConsume(r, 1010, &clocks, &checks, credX, 2)
+			edit(&rc)
+			if err := c.MintSessionForRequest(testCtx(t), consumeSession(1), rc); !errors.Is(err, ErrIncompleteConsume) || clocks != 0 {
+				t.Fatalf("err %v after %d clock reads, want ErrIncompleteConsume before the clock", err, clocks)
+			}
+			assertPendingNoSession(t, c, r)
+			if n, used := passkeyState(t, c); n != 1 || used != 0 {
+				t.Fatalf("x after a refusal: counter %d, last_used_at %d", n, used)
+			}
+		})
+	}
+	t.Run("C15-v5-store", func(t *testing.T) {
+		c := openV5(t)
+		r := putTestRequest(t, c)
+		var clocks, checks int
+		err := c.MintSessionForRequest(testCtx(t), consumeSession(1), passkeyConsume(r, 1010, &clocks, &checks, credX, 2))
+		if !isOlderPasskeys(err) || clocks != 0 {
+			t.Fatalf("err %v after %d clock reads, want the passkeys refusal before the clock", err, clocks)
+		}
+		assertPendingNoSession(t, c, r)
+	})
+}
+
+// CS8: the counter moves to the assertion's only forward (both zero: unchanged), with last_used_at,
+// approved_by and the session; otherwise ErrSignCountRegression and nothing changes.
+func TestPasskeyConsumeAdvancesTheCounter(t *testing.T) {
+	for _, tc := range []struct {
+		stored, count int64
+		cred          []byte
+		ok            bool
+	}{
+		{1, 2, credX, true}, {0, 0, credX, true}, {0, 7, credX, true}, {1, 1<<32 - 1, credX, true},
+		{1, 2, []byte{7}, true}, {1, 2, bytes.Repeat([]byte{7}, 1023), true},
+		{2, 2, credX, false}, {3, 2, credX, false}, {5, 0, credX, false}, {1<<32 - 1, 1<<32 - 1, credX, false},
+	} {
+		c, _ := passkeyAt(t, tc.cred, tc.stored)
+		r := putTestRequest(t, c)
+		var clocks, checks int
+		err := c.MintSessionForRequest(testCtx(t), consumeSession(1), passkeyConsume(r, 1010, &clocks, &checks, tc.cred, tc.count))
+		n, used := passkeyState(t, c)
+		if !tc.ok {
+			if !errors.Is(err, ErrSignCountRegression) || n != tc.stored || used != 0 {
+				t.Fatalf("stored %d, new %d: err %v, counter %d, last_used_at %d; want the regression and nothing changed", tc.stored, tc.count, err, n, used)
+			}
+			assertPendingNoSession(t, c, r)
+			continue
+		}
+		got, _, gerr := c.MintRequest(testCtx(t), r.RequestID)
+		sessions, _ := c.ListSessions(testCtx(t), "", MaxListSessions)
+		if err != nil || n != tc.count || used != 1010 || gerr != nil || got.ApprovedBy != "webauthn:x" || len(sessions) != 1 {
+			t.Fatalf("stored %d, new %d (cred %d bytes): err %v, counter %d, last_used_at %d, approved_by %q, %d sessions",
+				tc.stored, tc.count, len(tc.cred), err, n, used, got.ApprovedBy, len(sessions))
+		}
+	}
+	// The same counter again, on a fresh request, is a regression: the first consume moved it.
+	c, _ := passkeyAt(t, credX, 1)
+	var clocks, checks int
+	if err := c.MintSessionForRequest(testCtx(t), consumeSession(1), passkeyConsume(putTestRequest(t, c), 1010, &clocks, &checks, credX, 2)); err != nil {
+		t.Fatal(err)
+	}
+	r2 := testMintRequestRow(2, 1000)
+	if err := c.PutMintRequest(testCtx(t), r2, 1000); err != nil {
+		t.Fatal(err)
+	}
+	err := c.MintSessionForRequest(testCtx(t), consumeSession(2), passkeyConsume(r2, 1020, &clocks, &checks, credX, 2))
+	if n, used := passkeyState(t, c); !errors.Is(err, ErrSignCountRegression) || n != 2 || used != 1010 {
+		t.Fatalf("the replayed counter: err %v, counter %d, last_used_at %d; want the regression, 2, 1010", err, n, used)
+	}
+}
+
+// CS9: the passkey row is judged in the consume's transaction. A commit made on a second connection
+// after BeginTx (from the Clock hook, as A1b's CS3) that revokes x, rewrites its credential ID,
+// deletes its passkey row or raises its counter refuses the consume and rolls back everything.
+func TestPasskeyConsumeChecksTheRowInTheTransaction(t *testing.T) {
+	for _, a := range []struct {
+		name, sql string
+		want      error
+		counter   int64
+	}{
+		{"revoked", `UPDATE approvers SET revoked_at = 1005 WHERE name = 'x'`, ErrApproverRevoked, 1},
+		{"credential-rewritten", `UPDATE webauthn_credentials SET credential_id = x'02' WHERE name = 'x'`, ErrPasskeyChanged, 1},
+		{"passkey-deleted", `DELETE FROM webauthn_credentials WHERE name = 'x'`, ErrApproverUnknown, -1},
+		{"counter-raised", `UPDATE webauthn_credentials SET sign_count = 9 WHERE name = 'x'`, ErrSignCountRegression, 9},
+	} {
+		t.Run(a.name, func(t *testing.T) {
+			c, path := passkeyAt(t, credX, 1)
+			r := putTestRequest(t, c)
+			var clocks, checks int
+			rc := passkeyConsume(r, 1010, &clocks, &checks, credX, 2)
+			rc.Clock = func() int64 {
+				if _, err := rawDB(t, path).Exec(a.sql); err != nil {
+					t.Errorf("the second connection: %v", err)
+				}
+				return 1010
+			}
+			if err := c.MintSessionForRequest(testCtx(t), consumeSession(1), rc); !errors.Is(err, a.want) {
+				t.Fatalf("err %v, want %v", err, a.want)
+			}
+			assertPendingNoSession(t, c, r)
+			if x := approverX(t, c); x.LastUsedAt != 0 {
+				t.Fatalf("x's last_used_at %d after a refusal, want 0", x.LastUsedAt)
+			}
+			if p, ok, _ := c.Passkey(testCtx(t), "x"); a.counter >= 0 && (!ok || p.SignCount != a.counter) {
+				t.Fatalf("x's counter %d (ok %v), want %d", p.SignCount, ok, a.counter)
+			}
+		})
+	}
+}
--- a/host/store/passkeys.go
+++ b/host/store/passkeys.go
@@ -8,6 +8,7 @@
 // are the store's own backstop for the pure law the CLI runs first.
 
 import (
+	"bytes"
 	"context"
 	"database/sql"
 	"errors"
@@ -208,3 +209,40 @@
 func (c *Credentials) ActivePasskeys(ctx context.Context) ([]EnrolledPasskey, error) {
 	return c.enrolledPasskeys(ctx, `a.revoked_at IS NULL AND a.format = 'webauthn' ORDER BY a.name LIMIT ?`, maxActivePasskeys)
 }
+
+// ---- M5a ----
+
+// PasskeyUse is what a passkey consume binds: the credential that signed and its new counter.
+type PasskeyUse struct {
+	CredentialID []byte
+	SignCount    int64 // the assertion's counter, 0 .. 2^32-1
+}
+
+// usePasskey is the consume's counter binding (design §2.10): inside the consume's transaction,
+// after useApprover, it moves name's sign_count to p.SignCount only for this credential and only
+// forward (both zero: unchanged). Otherwise it names why (C16-C18), and the caller's rollback undoes
+// the request update and last_used_at.
+func usePasskey(ctx context.Context, tx *sql.Tx, name string, p PasskeyUse) error {
+	res, err := tx.ExecContext(ctx, `UPDATE webauthn_credentials SET sign_count = ? WHERE name = ? AND credential_id = ? AND (sign_count < ? OR (sign_count = 0 AND ? = 0))`,
+		p.SignCount, name, p.CredentialID, p.SignCount, p.SignCount)
+	if err != nil {
+		return fmt.Errorf("store: use passkey: %w", err)
+	}
+	if n, err := res.RowsAffected(); err != nil || n == 1 {
+		return err
+	}
+	var id []byte
+	var count int64
+	err = tx.QueryRowContext(ctx, `SELECT credential_id, sign_count FROM webauthn_credentials WHERE name = ?`, name).Scan(&id, &count)
+	switch {
+	case errors.Is(err, sql.ErrNoRows): // C16
+		return fmt.Errorf("%w: %q has no passkey", ErrApproverUnknown, name)
+	case err != nil:
+		return fmt.Errorf("store: use passkey: %w", err)
+	case !bytes.Equal(id, p.CredentialID): // C17
+		return fmt.Errorf("%w: %q", ErrPasskeyChanged, name)
+	default: // C18
+		return fmt.Errorf("%w: %q holds %d, this assertion %d; if you did not cause this, run: ailang-worldd approver revoke %s",
+			ErrSignCountRegression, name, count, p.SignCount, name)
+	}
+}
````

### Patch M5b — authority derives the passkey use (A6)

````diff M5b
--- a/host/authority/mint.go
+++ b/host/authority/mint.go
@@ -99,6 +99,9 @@
 	name, format, key, _ := rm.Signed.Use() // A3: Binds has shown it is verified
 	consume.ApprovedBy = rm.Signed.ApprovedBy()
 	consume.Approver = &store.ApproverUse{Name: name, Format: format, PublicKey: key}
+	if id, count, ok := rm.Signed.PasskeyUse(); ok { // A6: a passkey binds its credential and counter
+		consume.Approver.Passkey = &store.PasskeyUse{CredentialID: id, SignCount: int64(count)}
+	}
 	return consume, nil
 }
 
--- a/host/authority/mint_passkey_test.go
+++ b/host/authority/mint_passkey_test.go
@@ -0,0 +1,63 @@
+package authority
+
+// Row 201 A2 sprint 2 M5 (design §6, AU3): a verified passkey approval mints with approved_by
+// webauthn:<name> and moves the passkey's counter in the consume (A6). The passkey is approvaltest's
+// test authenticator (D-WORLD-109).
+
+import (
+	"encoding/hex"
+	"errors"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
+	"github.com/sunholo-data/ailang-world/host/store"
+)
+
+// AU3: A6 hands the store the credential and counter of a verified passkey approval; a zero
+// SignedApproval is still ErrNotVerified (A4) and writes nothing.
+func TestMintForRequestDerivesPasskeyUse(t *testing.T) {
+	const origin, rpID = "http://localhost:7644", "localhost"
+	pk := approvaltest.NewPasskey(approval.AlgES256, 7)
+	a, err := approval.NewPasskeyApprover(approval.RegistrationInput{Name: "mark-mac", Origin: origin, Challenge: [32]byte{9},
+		Response: pk.Register(origin, rpID, [32]byte{9}), Now: 900})
+	if err != nil {
+		t.Fatal(err)
+	}
+	p := a.Passkey
+	st, mr, row := au1Store(t)
+	if err := st.EnrolPasskey(boundedTestContext(t), store.ApproverRow{Name: a.Name, Format: a.Format, PublicKey: p.Key.Canonical(), EnrolledAt: 900},
+		store.PasskeyRow{Name: a.Name, CredentialID: p.CredentialID, RPID: p.RPID, Origin: p.Origin, Alg: p.Key.Alg, AAGUID: p.AAGUID[:], SignCount: int64(p.SignCount)}); err != nil {
+		t.Fatal(err)
+	}
+	d, _ := hex.DecodeString(row.Digest)
+	h, err := approval.MessageHash([32]byte(d), au1Challenge, row.ExpiresAt)
+	if err != nil {
+		t.Fatal(err)
+	}
+	s, err := approval.ApproveByAssertion(approval.AssertionInput{
+		Admit:     approval.AdmitInput{Now: 1010, CreatedAt: 1000, ExpiresAt: row.ExpiresAt, StoredDigest: [32]byte(d), BytesDigest: [32]byte(d)},
+		Challenge: au1Challenge, Response: pk.Assert(origin, rpID, h, 5), Approver: a,
+	})
+	if err != nil {
+		t.Fatal(err)
+	}
+	zero := au1RequestMint(mr, row)
+	zero.ExpiresAt, zero.Signed = row.ExpiresAt, &approval.SignedApproval{}
+	if err := au1Mint(t, st, zero); !errors.Is(err, approval.ErrNotVerified) {
+		t.Fatalf("a zero approval: %v, want ErrNotVerified", err)
+	}
+	assertNotMinted(t, st, row)
+	rm := au1RequestMint(mr, row)
+	rm.ExpiresAt, rm.Signed = row.ExpiresAt, s
+	if err := au1Mint(t, st, rm); err != nil {
+		t.Fatal(err)
+	}
+	ctx := boundedTestContext(t)
+	got, _, err1 := st.MintRequest(ctx, row.RequestID)
+	stored, _, err2 := st.Passkey(ctx, "mark-mac")
+	sessions, err3 := st.ListSessions(ctx, "", store.MaxListSessions)
+	if err1 != nil || err2 != nil || err3 != nil || got.ApprovedBy != "webauthn:mark-mac" || stored.SignCount != 5 || len(sessions) != 1 {
+		t.Fatalf("after the passkey mint: approved_by %q, counter %d, %d sessions (%v %v %v)", got.ApprovedBy, stored.SignCount, len(sessions), err1, err2, err3)
+	}
+}
````
