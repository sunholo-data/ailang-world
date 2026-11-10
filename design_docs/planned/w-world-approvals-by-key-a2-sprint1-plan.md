# Sprint plan — w-world-approvals-by-key-a2, sprint 1 (iteration 379, row 201)

> **THE CUT — the design's §9 fallback fired (both clauses).** Sprint 1 is "the verification law",
> registration half: **M0** the A1b judge's carry-forwards (K8a, K8b, SG1); **M1** the bounded CBOR
> subset and ES256/EdDSA COSE keys, with the Chrome-recorded fixtures pinned (XW0); **M2**
> authenticator data, `clientDataJSON`, the ceremony origin (with r2's default-port refusals) and the
> response envelope; **M3a** the registration law `NewPasskeyApprover`, the 128-bit enrolment code,
> `RenderPasskeyEnrolment`, `Approver.Passkey`, E2's rename and RC4.
> - **Measured: 1,958 changed lines under `cmd/ host/`** (4 sub-milestones; +1 line `.gitattributes`,
>   not counted) — **≤ 2,000**, 42 lines of margin.
> - **Clause 1 (the generator is not committed):** the fixtures are the designer's iteration-379
>   Chrome recording, normalised, SHA-pinned by XW0, with `SOURCE.txt`; `generate.sh` is a stub that
>   says *regeneration lands in sprint 3 (M7)*.
> - **Clause 2 (still over 2,000):** the whole sprint as designed measures ≈ 2,400 (M0–M3a measured
>   1,958 + the assertion half estimated ≈ 450 from these files' measured density). So **M3's PA\*
>   half moves to the head of sprint 2**: `ApproveByAssertion` (P1–P5, P3w, W1–W15, `checkApprover`),
>   `SignedApproval.PasskeyUse` (B7, B8), the test authenticator `approvaltest.Passkey`, PA1–PA4,
>   RC1, TP1, XW1 and XW2's assertion part (Appendix C lists them exactly).
>
> **Every line of code and test is a unified diff in Appendix P**, generated LAST from scratch stages
> that were compiled, vetted, gofmt-clean and green at every sub-milestone, after every mutant ran on
> the final tree; the patches were then extracted **from this file** with §2's `awk`, applied to a
> fresh export of `03d78b6a`, and the result compared byte-for-byte (`diff -r`) with the tree the
> mutants ran on (§1.2 C30–C31). The executor applies each patch and runs its acceptance; it writes no
> code by hand.

- **Design**: [w-world-approvals-by-key-a2.md](w-world-approvals-by-key-a2.md), **r2** (commit
  `03d78b6a`; r0 `72e20eb4`, r1 `967ed666`). **Not reopened.** Where this plan and the design differ,
  this plan was measured against `03d78b6a` and wins; each difference is a decision (§1.3) or a premise
  found false (§1.4).
- **Starting ground**: A1b is code-complete (PRs #537, #539, #541, #543; record #544). The iteration
  378 judge's two carry-forwards are M0. Iteration 378's retro directives are applied: every
  repo-wide inventory ran against the final tree (§1.2 C25–C27), the patches were regenerated last
  and round-tripped against the mutant tree (C30), and no gate's only output went through `tail`.
- **Planner base**: `03d78b6a16070a0fc91e8f6eb87d08e820df0cad` (= origin/dev `98374909` + three
  design-doc commits; `cmd/` and `host/` are byte-identical to `98374909`). Worktree
  `/Users/voightkampff/dev/sunholo-data/.wt-world-iter379-row201-a2`, branch
  `sprint/world-iter379-row201-a2`.
- **Executor**: as routed by the controller (a small model). **Two runs**: Run 1 (M0, M1), Run 2
  (M2, M3a). Each sub-milestone is one `patch -p1` (§2) plus its acceptance (§3). **Apply, never
  retype.**
- **Touched**: `cmd/ailang-worldd/{approver.go, approver_revoke_test.go, approver_enroll_test.go}`,
  `host/sshsign/sshsign_test.go`, `host/approval/{approver.go, approver_test.go}`, new
  `host/approval/{cbor,cose,authdata,clientdata,origin,passkey_json,passkey}.go`, new
  `host/approval/{cbor,cose,passkey,webauthn_fixture,webauthn_parts}_test.go`, new
  `host/approval/testdata/webauthn/{chrome-es256.json, chrome-eddsa.json, SOURCE.txt, generate.sh}`,
  `.gitattributes` (+1 line).
- **Not touched**: `host/approval/{signed.go, signed_test.go, request.go, admit.go, sshsig*.go}`,
  `host/approval/approvaltest/**` (moves to sprint 2 with the PA half), `host/sshsign/sshsign.go`,
  `host/store/**`, `host/authority/**`, `host/daemon/**`, `host/surfacepage/**`, every other
  `cmd/ailang-worldd` file, `docs/QUICKSTART.md`, `world/*.ail`, `packages/**`, `scripts/**`,
  `tools/launchd/*`, any skill, `~/.ssh`, `~/.ailang/world`, `~/.ailang/state/mission-v1*`.

---

## 0. Binding constraints (each is restated where it governs)

1. **Fakes, temp stores and test keys only (D-WORLD-109 = A).** The only keys are the fixtures'
   (Chrome 157's **virtual** authenticator, a test key, recorded by the designer in iteration 379) and
   values the tests derive from them. No browser, Chrome, Node, `ssh-keygen` or real authenticator
   runs in this sprint (the browser driver is sprint 3's M7). Every store is under `t.TempDir()`. The
   only subprocess a new test starts is SG1's temp `#!/bin/sh` script in `t.TempDir()`. Nothing opens
   `~/.ailang/world`, `~/.ssh`, a real credentials store, a real daemon or `/dev/tty`.
2. **Pure core, effects at the boundary (coding standards S1–S3).** Every new production file is in
   `host/approval`, standard library only, no I/O: `TestApprovalImportsOnlyTheStandardLibrary` stays
   green (new imports: `encoding/base32`, `encoding/hex`, `io`, `math`, `strconv`; none is forbidden).
   No new Go module (`TestDaemonDependencyAllowlist`, `TestWorldBoundaryDependencyAllowlist`), no new
   package, no production exec site, no route, no context root, no `.ail`.
3. **One mutant per new production branch** (§4: 197 rows run on the final tree, 184 KILLED, 13 named
   equivalents SURVIVED as predicted; plus 3 base controls, §1.2 C13). Each row names its single
   killing test.
4. **Typed errors, asserted with `errors.Is`, and the zero value on every refusal.** Every refusal arm
   checks the typed error **and** that nothing was returned (the zero `Approver`, `AuthData`,
   `ClientData`, `COSEKey`, response). M0's arms assert exit codes by constant (`exitOK = 0`,
   `exitUsage = 1`, `exitFatal = 2`).
5. **No `t.Parallel()`** in any new test (design §6; M0's tests swap package vars).
6. **Bisectable.** Each sub-milestone ends gofmt-clean, `go build ./...` and `go vet` green, the
   `_test.go` compile fence green, with its own tests and every earlier sub-milestone's tests green
   (§1.2 C31).
7. **No git writes by the executor** (§7). The controller commits from snapshots.

---

## 1. Pre-flight

### 1.1 Baseline at `03d78b6a` (measured by the planner, 2026-11-10)

Every command ran with `export AILANG_BIN=$HOME/.pinned-ailang/ailang` (`AILANG v0.41.0`), `go1.26.6
darwin/arm64`, in a `git archive 03d78b6a` export at
`/Users/voightkampff/dev/sunholo-data/.planner-scratch-iter379/base` (`diff -r` of `cmd/` and `host/`
against the worktree: identical; the worktree itself was only read). Each command's full output is in
`.planner-scratch-iter379/logs/baseline/<name>.log`.

| Command | rc | Time | Result |
|---|---|---|---|
| `git -C <worktree> status --porcelain \| wc -l` | 0 | — | `0` (clean at `03d78b6a`) |
| `test -z "$(gofmt -l cmd/ailang-worldd host/approval host/sshsign)"` | **0** | — | clean |
| `go build ./... && go vet ./cmd/ailang-worldd/ ./host/approval/... ./host/sshsign/` | **0** | 3 s | clean |
| `go test -run '^$' ./...` (compile fence for every `_test.go`) | **0** | 20 s | no FAIL line |
| `go test -count=1 ./host/approval/...` | **0** | 1 s | `host/approval` ok, `approvaltest` ok |
| `go test -count=1 ./host/sshsign/` | **0** | 1 s | ok |
| `go test -count=1 ./cmd/ailang-worldd/ -run '^TestApproverRevoke'` | **0** | 1 s | ok |
| `go test -count=1 ./cmd/ailang-worldd/ -run '^(TestApprover\|TestApprove\|TestConfirmApproverAct\|TestEveryVerbHasHelp)'` | **0** | 6 s | ok |
| `go test -count=1 ./cmd/ailang-worldd/` (whole package) | **0** | 114 s | ok |
| `go test -race -count=1 ./host/approval/... ./host/sshsign/` | **0** | 3 s | ok |
| the §7 inventories by name with `-v` (design §9's five commands, plus `TestBareNetHTTPExemptionIsPerGroup`, `TestReadGettersHonorContext`, and `host/approval/...`'s two guards) | **0** | 12 s | all 16 named tests `--- PASS` (`logs/baseline/inv_v.log`) |

**No red at the base.** Known load-sensitive legs (not chased, rule 3e): the local `-race` leg of
`host/broker`, `host/daemon`, `host/verifygate`, and `host/broker`'s
`TestWorktreeSnapshotRefusesAConfiguredFilter`. None is in a package this sprint touches.

### 1.2 Codebase claims checked (at `03d78b6a` unless "work" — the planner's final scratch tree `.planner-scratch-iter379/work`, outside the repo)

| # | Claim | Command → observed output |
|---|---|---|
| C1 | The base | `git rev-parse HEAD` → `03d78b6a16070a0fc91e8f6eb87d08e820df0cad`; `git rev-parse origin/dev` → `98374909268b9eab35382c5bf1779a3784edbf5b`; `git status --porcelain \| wc -l` → `0` |
| C2 | K8's two exits | `grep -n "return exitFatal" cmd/ailang-worldd/approver.go` → … `449`, `455`; `:451 rows, err := ro.DecidedMintRequests(list, row.Format+":"+name, 0)`; the seam pattern `:376 var approverRevoke = (*store.Credentials).RevokeApprover` |
| C3 | sshsign's G1 | `host/sshsign/sshsign.go:53 if err := cmd.Run(); err != nil { // G1`; its tests `:18 TestCappedBufferKeepsItsBound`, `:32 TestSignHandsTheSignerAScrubbedEnvironment` |
| C4 | E2 today | `host/approval/approver.go:39 ErrFormatNotYet = … "arrive with row 201 A2 …"`, `:53-54` E2; pinned by `approver_test.go:150-151` and `cmd/ailang-worldd/approver_enroll_test.go:230` (`want: []string{"row 201 A2"}`) |
| C5 | The fixture message hash | `request_test.go:35 const fixtureMessageHashHex = "29181d224d…fef5"` (the Chrome probe's challenge, design V21) |
| C6 | RC4 cannot take a pointer field | `approver_test.go:317 func approverFieldPaths` recurses into `reflect.Struct` only and `:372 default: t.Fatalf("… unhandled kind …")`: `Approver.Passkey *Passkey` must be walked (M3a edits it) |
| C7 | Nothing of A2 exists | `grep -rln "cbor\|webauthn" --include='*.go' -i host/approval` → `approver.go`, `approver_test.go`, `signed_test.go` only (the `FormatWebAuthn` constant and its refusal arms) |
| C8 | `.gitattributes` | last lines `host/approval/testdata/sshsig/* -text`, `host/approval/testdata/approvers/* -text`; no `webauthn` line |
| C9 | The designer's recording | `/tmp/w379-designer/browser/`: `es256.json` sha256 `c0370982…caf6`, `eddsa.json` `1544721f…f719`, `probe.mjs` `01a17c0d…37cf`; `jq '.result \| keys'` → `asr asr2 origin reg`; **`asr2` holds only `authenticatorData`** (`probe.mjs`: `const asr2 = { authenticatorData: … }`, §1.4 P2); origins `http://localhost:56547` (ES256) and `:56629` (EdDSA); `userHandle` null |
| C10 | Go 1.26 APIs | `go doc crypto/ecdsa ParseUncompressedPublicKey` → `(curve elliptic.Curve, data []byte) (*PublicKey, error)`, refuses a point not on the curve |
| C11 | The import guard | `imports_test.go:14-17` forbids `os`, `os/exec`, `os/signal`, `net`, `net/http`, `syscall`, `io/ioutil`, `crypto/rand`, `math/rand`, `math/rand/v2`, `unsafe` and any dotted first segment; the base's `host/approval` production imports include `crypto/ecdsa`, `crypto/ed25519`, `crypto/elliptic`, `encoding/json`, `unicode/utf8` |
| C12 | No other reflective walk of `approval.Approver` | `grep -rn "reflect.TypeOf(approval\.\|TypeOf(Approver" --include='*_test.go' cmd host` → `approver_test.go:343` (RC4) and `mint_request_test.go:57` (`MintRequest`, unaffected) |
| C13 | **The carry-forward mutants survive at the base** (the §4 harness, `tools/base_controls.json`, on the pristine base copy) | MUT-K8A with `-run '^TestApproverRevoke'` → **SURVIVED** (2.5 s); MUT-K8B → **SURVIVED** (2.4 s); MUT-SG1 with `./host/sshsign/` (all tests) → **SURVIVED** (0.7 s). Design V3/V4 re-measured TRUE. With M0 applied each dies by its single named test (§4) |
| C14 | `go build` is no fence for tests | each stage below also runs `go vet` of the touched packages and `go test -run '^$' ./...` |
| C20 | The enrolment code on the committed fixture | `EnrolmentCode` of `chrome-es256.json`'s registration = `YDOR4X36PQNFRTXBHWAV4W2XG4`, grouped `YDOR4-X36PQ-NFRTX-BHWAV-4W2XG4` — **equal to design V35's probe value**, because the fixture is that recording, normalised (PE4 pins it) |
| C21 | The fixtures are NUL-free JSON text | `LC_ALL=C tr -cd '\000' < <file> \| wc -c` → `0` for all four files; 10 lines each (2,151 and 1,984 bytes); sha256 `924e9f79…7639` (ES256) and `4c3fe46c…7b40` (EdDSA), pinned by XW0 |
| C22 | The fuzz leg | `go test ./host/approval/ -run '^$' -fuzz '^FuzzDecodeCBOR$' -fuzztime 30s` → `elapsed: 31s, execs: 1689399 … PASS`; no corpus written. A throwaway 20 s fuzz of `ParseClientData`, `ParseAuthData` and `NewPasskeyApprover`/`ParseAssertionResponse` (deleted afterwards) → PASS: no panic, no hang (§1.4 P7 records the hang it caught in a draft) |
| C25 | **Every repo-wide inventory, on the final tree** (work, `-v`, tee'd; each named test printed `--- PASS`) | `TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential` (AC10 census, `host/broker`); `TestDaemonDependencyAllowlist`, `TestFrozenV1RouteTableMatchesMux`, `TestSurfaceAppRoutes`, `TestDaemonConstructorCallerInventory` (`host/daemon`); `TestWorldBoundaryDependencyAllowlist`, `TestBareNetHTTPExemptionIsPerGroup`, `TestEveryCommandDirectoryIsAProtectedGroup` (`host/boundary`); `TestProductionContextRoots`, `TestProductionGoSurface`, `TestNoNewDeadlineFreeStoreReads`, `TestReadGettersHonorContext`, `TestCredentialsMethodsRequireADeadline` (`host/store`); `TestNoParallelWriteForkPackagesOutsideVerifygate` (`host/verifygate`); `TestApprovalImportsOnlyTheStandardLibrary`, `TestApprovalTestIsTestOnly` (`host/approval/...`) → all 16 printed `--- PASS`, 0 `--- SKIP`/`--- FAIL` (`logs/final_gates/inv_all.log`); no inventory changed: no exec site, module, package, route, root or `t.Parallel` was added |
| C26 | **Other tests that walk `cmd/` or `host/`** (directive (a)): `grep -rln --include='*_test.go' -E 'filepath\.WalkDir\|filepath\.Walk\(\|fs\.WalkDir\|os\.ReadDir' cmd host` → 46 files; the source-tree walkers among them are `host/broker/{registry_publish,invoke_boundary,handlers_parallel_guard}_test.go`, `host/daemon/startup_callers_test.go`, `host/store/{context_read,context_roots,deadline_guard,writer_lock_held}_test.go`, `host/boundary/allowlist_world_test.go`, `host/verifygate/{subprocess_sink_gate,forklocked_write,evidence_manifest_gate,module_manifest_gate,mission_config_gate}_test.go`, `host/approval/{imports_test.go, approvaltest/signer_test.go}` | all ran whole-package, concurrently: `go test -count=1 ./host/broker/ ./host/daemon/ ./host/boundary/ ./host/store/ ./host/verifygate/ ./host/approval/... ./host/sshsign/ ./cmd/ailang-worldd/` → every package ok except two **load-sensitive controls**: `host/broker` `TestWorktreeSnapshotRefusesAConfiguredFilter` (the known one, rule 3e) and `host/verifygate` `TestMissionConfigDefaultFlowIntegration/invalid_ledger` ("default-flow fixture timed out": a 10 s fixture deadline with five heavy packages running at once). Both PASS alone on the same tree (0.25 s; 7.30 s) and inside C27's sequential `go test ./...`. Neither package is touched |
| C27 | **The whole repository on the final tree** | `VERIFY_GO_LEGS=plain ./scripts/verify_go.sh` (tracked-binary gate on the committed fixtures, toolchain floor, race-detector control, `go build ./...`, the `host/evidence` manifest, `go test ./... -count=1`) → **rc 0** in 390 s: `✓ 0 other binary blobs among 6670 tracked files` (the fixtures are text), `✓ all 64 required top-level evidence tests passed exactly once`, 62 packages `ok` (`host/broker` 335 s, `host/daemon` 360 s, `host/verifygate` 333 s, `cmd/ailang-worldd` 219 s), `✓ go gate PASSED (plain)`; and `./scripts/verify_ail.sh` → rc 0 (54 s, `✓ verify gate PASSED`) |
| C28 | `go vet ./...` and gofmt on the final tree | `go vet ./...` rc 0; `gofmt -l cmd host` empty; `go test -race ./host/approval/... ./host/sshsign/` rc 0; `go test -race ./cmd/ailang-worldd/ -run '^(TestApproverRevoke\|TestApproverEnroll)'` rc 0 (16 s) |
| C30 | **The patches reproduce the mutant tree** | a fresh `git archive 03d78b6a` export + `patch -p1` of Appendix P's four patches **extracted from this file with §2's `awk`** (each `cmp`-equal to the generated patch) → `diff -r cmd host` and `cmp .gitattributes` against the tree §4's mutants ran on (`.planner-scratch-iter379/work`): **identical** (`ROUNDTRIP_IDENTICAL_TO_MUTANT_TREE`, `tools/roundtrip.sh`, `logs/roundtrip.log`); no `.orig`/`.rej`; nothing left uncommitted |
| C31 | **Every sub-milestone is bisectable** | the same fresh export, after **each** patch: gofmt empty, build + vet clean, compile fence rc 0, the §3 acceptance commands rc 0 with no `--- SKIP`/`--- FAIL`, the 30 s fuzz leg PASS, `host/approval/... host/sshsign` ok; the cumulative numstat after a commit of each stage = §7's value: **75, 905, 1510, 1958** (`NUMSTAT` lines of `logs/roundtrip.log`); every `-v` log free of `--- SKIP`/`--- FAIL`; the fuzz leg `PASS` with no corpus |

### 1.3 Planner decisions where the design leaves a choice

- **S1-D1 — The cut** (§1.5): the §9 fallback's two clauses, in order. M3 is split into **M3a**
  (registration) in this sprint and **M3b** (the assertion half) at the head of sprint 2.
- **S1-D2 — Fixtures land in M1, not M3**, with XW0 and `SOURCE.txt`, because M1's own tests (CB3,
  CK2, CK3) read Chrome's bytes; a sub-milestone must be green on its own (§0.6).
- **S1-D3 — The fixture form.** `chrome-<alg>.json` holds `source`, `chrome`, `origin`, `rp_id`,
  `challenge_hex`, and **`registration` and `assertion` exactly as Chrome's `toJSON()` returned them**
  (so EN1 and PE2 also prove the envelope ignores `toJSON()`'s extra members), and `assertion2` =
  `{authenticatorData}` (all the probe kept, §1.4 P2). One JSON value per line (10 lines each).
  Produced from the designer's probe outputs by `.planner-scratch-iter379/tools/mkfixtures.py`
  (deterministic). `generate.sh` is a non-executable stub (`sh generate.sh` exits 1 naming sprint 3);
  `patch(1)` cannot set a mode bit, and the M7 driver replaces it.
- **S1-D4 — AD1, PE1, PE3's backup arms and RC4 build their inputs by editing Chrome's recorded
  bytes**, inside package `approval`, with a 40-line test CBOR encoder (`cb`, M1). A registration
  carries no signature (`fmt: none`), so every R-arm (type, challenge, origin, `crossOrigin`,
  `topOrigin`, `fmt`, `attStmt`, flags, rpID hash, `rawId`) is a faithful edit of real browser
  output. The test authenticator is first needed by a signing law (PA1) and moves with it (§1.4 P4).
- **S1-D5 — Typed errors carry detail by wrapping** (`fmt.Errorf("%w: …", ErrCOSEKey)`); R9 wraps
  the CBOR error too (`ErrAttestationObject` *and* `ErrCBORTrailing`), so PE1's R9 arm can tell R9's
  decode branch from its shape branch.
- **S1-D6 — `Approver.Key` stays the zero `SSHPublicKey` for a passkey approver**; the passkey's key
  is `Passkey.Key`, and the stored `approvers.public_key` (sprint 2) is `Passkey.Key.Canonical()`.
- **S1-D7 — `ParseRegistrationResponse` / `ParseAssertionResponse` are exported** (design §2.7 names
  the checks, not the functions): sprint 3's CLI needs the assertion's `rawId` before the store lookup
  (WF3, WF5). `ParseAssertionResponse` lands now with EN1 because E2–E4 are one envelope.
- **S1-D8 — `RenderPasskeyEnrolment`'s layout** (the design lists its contents):
  `approver / format / passkey <ALG> for <rpID> at <origin> / cred <16 hex> / aaguid <32 hex> / backup
  <synced passkey|backup-eligible|device-bound> / counter <n> / records approved_by=webauthn:<name>`,
  a blank line, then the §2.8 paragraph verbatim (`PasskeyEnrolmentWarning`). PE3 pins it whole.
- **S1-D9 — RC4's dispositions for the new fields**: `Passkey.RPID`, `.Origin`, `.SignCount` shown;
  `.CredentialID` "shown as the first 16 hex of its SHA-256"; `.Key.Alg` "shown as ES256 or EdDSA";
  `.AAGUID` "shown in hex"; `.BackupEligible`, `.BackupState` "shown as synced passkey,
  backup-eligible or device-bound"; `.Key.X`, `.Key.Y` "excluded: bound by the typed enrolment code, a
  fingerprint of the registration that carries it". `Passkey.*` paths are perturbed on a
  backup-eligible passkey through `RenderPasskeyEnrolment`; the rest on the ssh base as before.
- **S1-D10 — The ceremony origin's extra arms**: a 253-character host accepted and a 254-character
  host refused (the `len(host) > 253` branch), `:65535` accepted and `:65536`, `:+8443`, `:` refused,
  `ftp://localhost` and a bare `localhost` refused. `http://127.0.0.1:7644`, `https://evil.example`
  and `https://ts.net` are `ErrOriginHost`; every other refusal is `ErrOrigin`.
- **S1-D11 — CB4's property** is "never panics, and the returned remainder is never longer than the
  input": production has no general CBOR encoder to re-encode with (only `COSEKey.Canonical`).
- **S1-D12 — The registration law's extra arms** beyond one per R: R9 as an array of six, with an
  extra key, and with `authData` as text; R10 for `fmt` and for `attStmt`; R8 for `crossOrigin` and
  for `topOrigin`; and PE2's second registration (the same bytes re-made for
  `https://x.ts.net:8443`, flags BE|BS, counter 7) so that every copied `Passkey` field is observable.

### 1.4 Premises found false or incomplete at HEAD (each routed around minimally)

- **P1 — Design §2.15 and §9: "regenerated by the committed generator".** The generator is M7's (sprint
  3) and is not committed. Routed by the §9 fallback's first clause (S1-D3): the iteration-379
  recording is committed, pinned by XW0, `generate.sh` says regeneration lands in sprint 3.
- **P2 — Design §2.15: each fixture records "an `assertion2` (counter 3)".** The designer's probe kept
  only `assertion2`'s `authenticatorData` (`probe.mjs`: `const asr2 = { authenticatorData:
  b64u(a2.response.authenticatorData) }`; C9). There is no second `clientDataJSON` or signature to
  verify. Routed: AD2 reads its counter (3) from the recorded bytes; sprint 2's PA3 last arm becomes
  "the stored counter is assertion2's recorded counter (3); `assertion` (2) is refused by W15" —
  the same real-byte refusal, without verifying `assertion2` itself. Sprint 3's M7 re-records a full
  `assertion2` with the committed driver.
- **P3 — Design §8 puts `cmd/ailang-worldd/approver_enroll_test.go:230` in sprint 3**, but PE5's E2
  rename (design §9: M3, sprint 1) changes the CLI's stderr, so that arm (`want: "row 201 A2"`) goes red
  in sprint 1. Routed: M3a edits that one `want` to `--registration` (design §6's named edit, moved
  forward); WN1's usage refusal for the same arm stays sprint 3. Between sprints 1 and 3, `approver
  enroll --format webauthn` prints E2's new text, which names `--origin`/`--registration` flags that
  sprint 3 adds (no release ships between the sprints of row 201).
- **P4 — Design §6 builds AD1 "from `approvaltest.Passkey` knobs", but §9 lands the test
  authenticator in M3, after M2**; and package `approval`'s internal tests cannot import
  `approvaltest`, which imports `approval` (Go refuses the cycle). Routed by S1-D4: AD1 edits Chrome's
  bytes in package `approval`; `approvaltest.Passkey` lands with PA1 (sprint 2), in an external
  `approval_test` file.
- **P5 — Design §9: sprint 1 ≈ 1,800 lines.** Measured M0 75 + M1 830 + M2 605 + M3a 448 = **1,958**
  without the assertion half (design: M1 ≈ 700, M2 ≈ 610, M3 ≈ 970 + fixtures). Routed by the second
  fallback clause (S1-D1).
- **P6 — Design §3: "The new imports are `crypto/subtle` … and packages it already imports".** Sprint 1
  adds `encoding/base32`, `encoding/hex`, `io`, `math` and `strconv` (all standard library, none
  forbidden, C11); `crypto/subtle` arrives with WN10 in sprint 3. No routing needed.
- **P7 — A draft's assumption, caught by measurement:** "a `json.Decoder` that failed keeps failing"
  is false for `More()` (it peeks bytes, not the error): a draft `skipJSON` that ignored token errors
  looped forever on `{"x":{1:2}}`. Routed: every `Token()` error in `skipJSON` is refused (J3), and the
  arms `J3-unknown-member-no-value`, `-number-key` and `-unterminated` pin it (MUT-SKIP-TOK,
  MUT-SKIP-CLOSE).
- **P8 — Design §6's "refused before verification" (CK2: a 73-byte ES256 and a 63-byte Ed25519
  signature) is not observable**: the library refuses the same inputs. The explicit bounds are kept
  (design §2.4) and their mutants are named equivalents (MUT-ES-LEN-EQUIV, MUT-ED-LEN-EQUIV); CK2
  still asserts both inputs are refused.

### 1.5 The size, and the cut

Lines are insertions + deletions under `cmd/ host/` (`git diff --numstat`; `.gitattributes` not
counted). **Every row is measured** from the scratch stages, each committed in a throwaway repo (C31).

| Sub-milestone | Production | Test | Testdata | Step total | Cumulative vs `03d78b6a` (after its commit) | Run |
|---|---|---|---|---|---|---|
| M0 — K8a, K8b (the `approverDecided` seam), SG1 | 6 | 69 | 0 | **75** | 75 | 1 |
| M1 — CBOR subset, COSE keys, the fixtures, XW0, CB1–CB4, CK1–CK4 | 323 | 460 | 47 | **830** | 905 | 1 |
| M2 — authData, clientData, origin, envelope; AD1–AD3, CD1–CD2, OR1, EN1 | 343 | 262 | 0 | **605** | 1,510 | 2 |
| M3a — the registration law, the enrolment code, the prompt, `Approver.Passkey`, E2; PE1–PE5, RC4 | 179 | 269 | 0 | **448** | **1,958** | 2 |
| **Sprint** | **851** | **1,060** | **47** | **1,958 ≤ 2,000** | | |
| *moved to sprint 2's head (§1.3 S1-D1, Appendix C)* — M3b, the assertion half | ≈ 200 | ≈ 250 | — | *≈ 450 (estimate)* | *≈ 2,410 if kept* | — |

---

## 2. How the executor applies a sub-milestone (the patch protocol)

Each sub-milestone `<M>` is the block fenced by a line `` ````diff <M> `` (four backticks) and a line of
four backticks, in Appendix P. Paths are `a/<repo path>` and `b/<repo path>`. Apply it **only** this
way, in one `bash` shell (the controller's shell is zsh):

```bash
bash <<'EOF'
set -euo pipefail
M=M0   # <- the sub-milestone, exactly as in the Appendix P heading: M0, M1, M2 or M3a
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter379-row201-a2
PLAN="$W/design_docs/planned/w-world-approvals-by-key-a2-sprint1-plan.md"
T="$HOME/.ailang/state/world-iter379/exec-patches"; mkdir -p "$T"
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
  `$HOME/.ailang/state/world-iter379/exec-logs/` (never only through `tail`).
- Then snapshot (§7). The controller commits and compares the numstat with §7's expected value.

---

## 3. Milestones

Every sub-milestone's acceptance starts with the same three lines (abbreviated **[base]** below):

```
test -z "$(gofmt -l cmd/ailang-worldd host/approval host/sshsign)"
go build ./... && go vet ./cmd/ailang-worldd/ ./host/approval/... ./host/sshsign/
go test -run '^$' ./...
```

Each `-v` run below must show `--- PASS` for every test it names and **no `--- SKIP` or `--- FAIL`**.

### M0 — the carry-forwards K8a, K8b, SG1 · 6 prod + 69 test · Run 1

**Patch** `M0`:
- `cmd/ailang-worldd/approver.go`: `var approverDecided = (*store.Credentials).DecidedMintRequests`
  (a test seam beside `approverRevoke`; production never reassigns it), and K8's read calls
  `approverDecided(ro, list, row.Format+":"+name, 0)`. No behaviour change.
- `cmd/ailang-worldd/approver_revoke_test.go`: **K8a** `TestApproverRevokeListingReopenFailureIsFatal`
  (the `approverRevoke` seam runs the real revoke, then renames the credentials file away → `exitFatal`;
  stdout exactly `revoked approver mark-mac\n`; stderr has `revoked, but cannot read the credentials
  store` and the `approve --list --decided` fix); **K8b**
  `TestApproverRevokeListingQueryFailureIsFatal` (the `approverDecided` seam fails → `exitFatal`,
  `revoked, but cannot list what it approved: disk on fire`; the row is revoked at `approveT0+100`).
- `host/sshsign/sshsign_test.go`: **SG1** `TestSignRefusesAFailedSignerEvenWithOutput` (a temp
  `#!/bin/sh` signer prints `partial-signature` and exits 3 → nil bytes and an error wrapping
  `*exec.ExitError` with exit code 3).

**Acceptance**
```
[base]
go test ./cmd/ailang-worldd/ -run '^(TestApproverRevokeListingReopenFailureIsFatal|TestApproverRevokeListingQueryFailureIsFatal|TestApproverRevoke|TestApproverRevokeUnansweredPromptTimesOut|TestApproverRevokeListsWhatItApproved)$' -count=1 -v
go test ./host/sshsign/ -count=1 -v
```

**Expected cumulative numstat after the controller's commit of M0: 75**.

### M1 — CBOR and COSE, and the recorded fixtures · 323 prod + 460 test + 47 testdata · Run 1

**Patch** `M1`:
- new `host/approval/cbor.go`: `cborValue` (int, bytes, text, array, map in read order, bool, null),
  `cborLimits`, `attestationLimits = {maxDepth 4, maxItemsPerContainer 16, maxTotalItems 64,
  maxStringBytes 2048}`, `decodeCBOR`, `decodeCBORExact`; K1–K16 each its own branch with its typed
  error (`ErrCBORTruncated`, `…Malformed`, `…Indefinite`, `…NonCanonical`, `…TooLarge`, `…TooDeep`,
  `…Unsupported`, `…DuplicateKey`, `…Trailing`). Map key order is not checked; duplicates are.
- new `host/approval/cose.go`: `COSEKey{Alg, X, Y}`, `AlgES256 = -7`, `AlgEdDSA = -8`, `ParseCOSEKey`
  (Q1–Q8, `ErrCOSEKey`/`ErrCOSEAlg`), `Canonical()` (CTAP2 order: what sprint 2 stores),
  `AlgName()`, `verify` (ES256: ≤ 72-byte DER over SHA-256 with `ecdsa.VerifyASN1`; EdDSA: 64 bytes,
  `ed25519.Verify`; else `ErrCOSEAlg`).
- new `host/approval/testdata/webauthn/{chrome-es256.json, chrome-eddsa.json, SOURCE.txt,
  generate.sh}` (S1-D3) and `.gitattributes` `host/approval/testdata/webauthn/* -text`.
- new `host/approval/webauthn_fixture_test.go`: `chromeFixture` (the decoded recording),
  `regAuthData`, `coseKeyBytes`; **XW0** `TestPasskeyFixturesAreTheGeneratedBytes` (both sha256 pins;
  `SOURCE.txt` names `Chrome/157.0.8103.50`, `virtual authenticator`, `D-WORLD-109`; rp `localhost`,
  challenge = `fixtureMessageHashHex`); the test encoder `cb`.
- new `cbor_test.go`: **CB1** `TestCBORRefusals` (35 arms over K1–K16, incl. `0x18 0x17` for K4's
  one-byte boundary), **CB2** `TestCBORBoundsAreInclusive` (2,048 bytes/text, 16 entries, depth 4,
  64 items, `int64` limits, keys `0` and `""` distinct; decoded values), **CB3**
  `TestCBORAcceptsChromesBytes`, **CB4** `FuzzDecodeCBOR`.
- new `cose_test.go`: **CK1** `TestCOSEKeyRefusals` (18 arms; the zero key each time; both controls
  accepted), **CK2** `TestCOSEKeyVerifies` (Chrome's signatures over `authData ‖ SHA-256(cdj)`;
  tampered, `authData`-only, 73-byte and 63-byte arms refused), **CK3** `TestCOSECanonicalIsChromesBytes`
  (and a permuted map), **CK4** `TestCOSEKeyHasOnlyDeclaredFields`.

**Acceptance**
```
[base]
go test ./host/approval/ -run 'CBOR|COSE|PasskeyFixtures' -count=1 -v
go test ./host/approval/ -run '^$' -fuzz '^FuzzDecodeCBOR$' -fuzztime 30s     # CB4: must end PASS; must write no testdata/fuzz corpus
git status --porcelain -- host/approval/testdata/fuzz                         # must print nothing
go test ./host/approval/... -count=1
```

**Expected cumulative numstat after the controller's commit of M1: 905**.

### M2 — authenticator data, clientDataJSON, the origin and the envelope · 343 prod + 262 test · Run 2

**Patch** `M2`:
- new `authdata.go`: `AuthData`, `ParseAuthData(b, registration)` (D1–D8; flags UP 0x01, UV 0x04,
  BE 0x08, BS 0x10, AT 0x40, ED 0x80; `MaxCredentialIDBytes = 1023`, floor 1, **[r1]**).
- new `clientdata.go`: `ClientData`, `ParseClientData` (J1–J7 by `json.Decoder` tokens;
  `MaxClientDataBytes = 2048`; unknown members skipped to depth 4; every token error refused).
- new `origin.go`: `ParseCeremonyOrigin` (O1 grammar, O2 host, http for localhost only, O3 a bad port
  **and r2's default port** with guidance `leave out the default port :443 …`).
- new `passkey_json.go`: `RegistrationResponse`, `AssertionResponse`, `ParseRegistrationResponse`,
  `ParseAssertionResponse` (E2–E4; strict unpadded base64url; `userHandle` null/absent → none),
  `MaxPasskeyResponseBytes = 16 KiB` (E1 is the CLI's, sprint 3).
- new `webauthn_parts_test.go`: **AD1** `TestAuthDataRefusals` (14 arms; credential IDs of 0 and 1,024
  bytes refused, 1 and 1,023 accepted), **AD2** `TestAuthDataReadsChrome` (flags, counters 1/2/3,
  AAGUID, the 32-byte credential ID = `rawId`, the key), **AD3**
  `TestAuthDataIgnoresExtensionsContents`, **CD1** `TestClientDataRefusals` (20 arms + 3 acceptance
  controls + the 2,048-byte boundary), **CD2** `TestClientDataReadsChrome`, **OR1**
  `TestCeremonyOrigin` (8 accepted, 20 refused; the three r2 default-port arms each name the rule),
  **EN1** `TestPasskeyEnvelope` (registration and assertion arms; `userHandle` absent/null/present;
  both Chrome `toJSON()` envelopes).

**Acceptance**
```
[base]
go test ./host/approval/ -run 'AuthData|ClientData|CeremonyOrigin|PasskeyEnvelope' -count=1 -v
go test ./host/approval/... -count=1
```

**Expected cumulative numstat after the controller's commit of M2: 1510**.

### M3a — the registration law, the enrolment code and the prompt · 179 prod + 269 test · Run 2

**Patch** `M3a`:
- new `passkey.go`: `Passkey`, `RegistrationInput`, `NewPasskeyApprover` (R1–R15 in the design's
  order; R16 is Q3 via R11), `attestedAuthData` (R9 decode and shape, R10 `fmt: none` + empty
  `attStmt`), `EnrolmentCode` (26 characters, 128 bits, **[r1]**), `GroupEnrolmentCode` (5-5-5-5-6),
  `PasskeyEnrolmentWarning` (the §2.8 paragraph), `RenderPasskeyEnrolment` (S1-D8); errors
  `ErrClientDataType`, `ErrChallenge`, `ErrOriginMismatch`, `ErrCrossOrigin`, `ErrAttestationObject`,
  `ErrAttestationFormat`, `ErrRPIDHash`, `ErrUserPresence`, `ErrUserVerification`,
  `ErrCredentialIDMismatch`.
- `approver.go`: `Approver.Passkey *Passkey` (nil for ssh); `ErrFormatNotYet` → **`ErrUseRegistration`**
  with the design's text; E2 returns it.
- `approver_test.go`: **PE5** (the two webauthn arms expect `ErrUseRegistration`); **RC4** edited
  (pointer-to-struct walk; the S1-D9 dispositions; `Bool`, `Uint32`, `Array` perturbations).
- `cmd/ailang-worldd/approver_enroll_test.go`: the `webauthn` arm's want `row 201 A2` → `--registration`
  (§1.4 P3).
- new `passkey_test.go`: **PE1** `TestNewPasskeyApproverRefusals` (21 arms over R1–R15, the zero
  `Approver` each time, plus two controls), **PE2** `TestNewPasskeyApproverAcceptsChrome` (both
  recordings, every field; the tailnet/synced/counter-7 re-make), **PE3**
  `TestRenderPasskeyEnrolmentGolden` (golden; `synced passkey`, `backup-eligible`, `EdDSA`), **PE4**
  `TestEnrolmentCodeGolden` (`YDOR4X36PQNFRTXBHWAV4W2XG4`, `YDOR4-X36PQ-NFRTX-BHWAV-4W2XG4`; one byte
  of either input, or swapping them, changes it).

**Acceptance**
```
[base]
go test ./host/approval/... -run 'Passkey|Enrolment|NewApprover|Imports|ApprovalTest' -count=1 -v
go test ./cmd/ailang-worldd/ -run '^TestApproverEnroll' -count=1 -v
go test ./host/approval/... ./host/sshsign/ -count=1
```

**Expected cumulative numstat after the controller's commit of M3a: 1958**.

---

## 4. Mutation table (197 rows run on the planner's final tree: 184 KILLED, 13 named equivalents SURVIVED as predicted; plus 3 base controls, §1.2 C13)

**Derivation.** Anchored to Appendix P's production hunks: one row per new production `if`, `case`,
condition operand that a test can tell apart, bound, flag bit, copied field and constant, plus the
design §6 named mutants (R14 with a UV-clear arm; the 16-byte credential-ID floor; Q3 accepting
−257; the duplicate-key detection; `<`/`<=` at every bound). **Procedure**: the harness
`.planner-scratch-iter379/tools/mutate.py <tree> tools/mutants.json [ids…]` (outside the repo; the
iteration-378 harness, unchanged): assert the anchor occurs exactly once, apply **one** edit, assert the
file changed, `go vet` the package (a mutant that does not compile is not a mutant — the first pass
had 18 NOCOMPILE drafts, each rewritten into a compiling form and re-run; only the re-run counts; it also had
two unexpected survivors, MUT-D4A (the truncated arms were sub-slices whose capacity let the deleted bound
read past the cut) and MUT-E2-RAWID (no arm lacked `rawId`): both arms were strengthened (`bytes.Clone`,
`E2-no-rawId`) and the whole table re-run on the final tree),
run **the single named test** (`go test <pkg> -run '^<test>$' -count=1`), restore, assert the restore is
byte-identical. Run with `export AILANG_BIN=$HOME/.pinned-ailang/ailang`. Logs:
`.planner-scratch-iter379/logs/mutants_final.log`. **Never commit or snapshot a mutant.** A KILLED
row marked "(panics)" fails its test by a runtime panic: the deleted check was the only guard.

| Id | Stage | Mutant (single edit) | Killing test (arm) | Scratch |
|---|---|---|---|---|
| MUT-K8A | M0 | cmd/approver.go: `to list what it approved: %v; fix: ailang-worldd approve --list --decided\n", prefix, credPath, err)⏎return exitFatal` → `to list what it approved: %v; fix: ailang-worldd approve --list --decided\n", prefix, credPath, err)⏎return exitOK` | `TestApproverRevokeListingReopenFailureIsFatal` (the file moved away after the revoke) | KILLED |
| MUT-K8B | M0 | cmd/approver.go: `cannot list what it approved: %v; fix: ailang-worldd approve --list --decided\n", prefix, err)⏎return exitFatal` → `cannot list what it approved: %v; fix: ailang-worldd approve --list --decided\n", prefix, err)⏎return exitOK` | `TestApproverRevokeListingQueryFailureIsFatal` (the seam fails) | KILLED |
| MUT-K8-SEAM | M0 | cmd/approver.go: `rows, err := approverDecided(ro, list,` → `rows, err := ro.DecidedMintRequests(list,` | `TestApproverRevokeListingQueryFailureIsFatal` (the seam bypassed) | KILLED |
| MUT-SG1 | M0 | host/sshsign/sshsign.go: `if err := cmd.Run(); err != nil { // G1` → `if err := cmd.Run(); false && err != nil { // G1` | `TestSignRefusesAFailedSignerEvenWithOutput` (prints partial-signature, exits 3) | KILLED |
| MUT-K1-EMPTY | M1 | cbor.go: `if len(d.b) == 0 { // K1` → `if false { // K1` | `TestCBORRefusals` (K1-empty (panics)) | KILLED |
| MUT-K9 | M1 | cbor.go: `case major == 6: // K9` → `case false: // K9` | `TestCBORRefusals` (K9-tag) | KILLED |
| MUT-K10 | M1 | cbor.go: `if ai < 20 \|\| ai > 22 { // K10` → `if false { // K10` | `TestCBORRefusals` (K10-*) | KILLED |
| MUT-K10-LOW | M1 | cbor.go: `if ai < 20 \|\| ai > 22 { // K10` → `if ai < 19 \|\| ai > 22 { // K10` | `TestCBORRefusals` (K10-simple-19) | KILLED |
| MUT-K10-HIGH | M1 | cbor.go: `if ai < 20 \|\| ai > 22 { // K10` → `if ai < 20 \|\| ai > 23 { // K10` | `TestCBORRefusals` (K10-undefined) | KILLED |
| MUT-K3 | M1 | cbor.go: `case ai == 31: // K3` → `case false: // K3` | `TestCBORRefusals` (K3-indefinite-*) | KILLED |
| MUT-K2 | M1 | cbor.go: `case ai > 27: // K2` → `case false: // K2` | `TestCBORRefusals` (K2-reserved) | KILLED |
| MUT-K1-ARG | M1 | cbor.go: `if len(d.b) < 1+n { // K1` → `if false { // K1` | `TestCBORRefusals` (K1-head-argument (panics)) | KILLED |
| MUT-K4 | M1 | cbor.go: `if arg < [4]uint64{24, 1 << 8, 1 << 16, 1 << 32}[ai-24] { // K4` → `if false && arg < [4]uint64{24, 1 << 8, 1 << 16, 1 << 32}[ai-24] { // K4` | `TestCBORRefusals` (K4-*) | KILLED |
| MUT-K4-24 | M1 | cbor.go: `[4]uint64{24, 1 << 8,` → `[4]uint64{23, 1 << 8,` | `TestCBORRefusals` (K4-one-byte-23) | KILLED |
| MUT-K4-8 | M1 | cbor.go: `{24, 1 << 8, 1 << 16,` → `{24, 1<<8 - 1, 1 << 16,` | `TestCBORRefusals` (K4-two-byte (255)) | KILLED |
| MUT-K4-16 | M1 | cbor.go: `1 << 16, 1 << 32}` → `1<<16 - 1, 1 << 32}` | `TestCBORRefusals` (K4-four-byte (65535)) | KILLED |
| MUT-K4-32 | M1 | cbor.go: `1 << 32}[ai-24]` → `1<<32 - 1}[ai-24]` | `TestCBORRefusals` (K4-eight-byte (2^32-1)) | KILLED |
| MUT-K8-OFF | M1 | cbor.go: `if depth > d.lim.maxDepth { // K8` → `if false { // K8` | `TestCBORRefusals` (K8-five-nested) | KILLED |
| MUT-K8-BOUND | M1 | cbor.go: `if depth > d.lim.maxDepth { // K8` → `if depth >= d.lim.maxDepth { // K8` | `TestCBORBoundsAreInclusive` (depth-4) | KILLED |
| MUT-K15-OFF | M1 | cbor.go: `if d.total++; d.total > d.lim.maxTotalItems { // K15` → `if d.total++; false { // K15` | `TestCBORRefusals` (K15-65-items) | KILLED |
| MUT-K15-BOUND | M1 | cbor.go: `if d.total++; d.total > d.lim.maxTotalItems { // K15` → `if d.total++; d.total >= d.lim.maxTotalItems { // K15` | `TestCBORBoundsAreInclusive` (64-items) | KILLED |
| MUT-K14-OFF | M1 | cbor.go: `if arg > math.MaxInt64 { // K14` → `if false && arg > math.MaxInt64 { // K14` | `TestCBORRefusals` (K14-*) | KILLED |
| MUT-K14-BOUND | M1 | cbor.go: `if arg > math.MaxInt64 { // K14` → `if arg >= math.MaxInt64 { // K14` | `TestCBORBoundsAreInclusive` (int64-limits) | KILLED |
| MUT-NEG | M1 | cbor.go: `i: -1 - int64(arg)}` → `i: -int64(arg)}` | `TestCBORBoundsAreInclusive` (the decoded -1) | KILLED |
| MUT-K6-OFF | M1 | cbor.go: `if arg > uint64(d.lim.maxStringBytes) { // K6` → `if false { // K6` | `TestCBORRefusals` (K6-2049-*) | KILLED |
| MUT-K6-BOUND | M1 | cbor.go: `if arg > uint64(d.lim.maxStringBytes) { // K6` → `if arg >= uint64(d.lim.maxStringBytes) { // K6` | `TestCBORBoundsAreInclusive` (2048-bytes) | KILLED |
| MUT-K5-OFF | M1 | cbor.go: `if arg > uint64(len(d.b)) { // K5` → `if false { // K5` | `TestCBORRefusals` (K5-* (panics)) | KILLED |
| MUT-K11 | M1 | cbor.go: `if major == 3 && !utf8.Valid(s) { // K11` → `if false && major == 3 && !utf8.Valid(s) { // K11` | `TestCBORRefusals` (K11-utf8) | KILLED |
| MUT-TEXT-KIND | M1 | cbor.go: `kind: cborBytes + cborKind(major-2)` → `kind: cborBytes` | `TestCBORBoundsAreInclusive` (the decoded text) | KILLED |
| MUT-K7-OFF | M1 | cbor.go: `if arg > uint64(d.lim.maxItemsPerContainer) { // K7` → `if false { // K7` | `TestCBORRefusals` (K7-17-*) | KILLED |
| MUT-K7-BOUND | M1 | cbor.go: `if arg > uint64(d.lim.maxItemsPerContainer) { // K7` → `if arg >= uint64(d.lim.maxItemsPerContainer) { // K7` | `TestCBORBoundsAreInclusive` (16-entries) | KILLED |
| MUT-MAP-KIND | M1 | cbor.go: `kind: cborArray + cborKind(major-4)` → `kind: cborArray` | `TestCBORAcceptsChromesBytes` (the attestationObject is a map) | KILLED |
| MUT-K12 | M1 | cbor.go: `if e.kind != cborInt && e.kind != cborText { // K12` → `if false { // K12` | `TestCBORRefusals` (K12-*) | KILLED |
| MUT-K13 | M1 | cbor.go: `if _, dup := v.get(e); dup { // K13` → `if _, dup := v.get(e); false && dup { // K13` | `TestCBORRefusals` (K13-*) | KILLED |
| MUT-NULL | M1 | cbor.go: `if arg == 22 {` → `if arg == 21 {` | `TestCBORBoundsAreInclusive` (true and null decoded) | KILLED |
| MUT-TRUE | M1 | cbor.go: `t: arg == 21}` → `t: arg == 20}` | `TestCBORBoundsAreInclusive` (true decoded) | KILLED |
| MUT-K16 | M1 | cbor.go: `if len(rest) != 0 { // K16` → `if false && len(rest) != 0 { // K16` | `TestCBORRefusals` (K16-trailing) | KILLED |
| MUT-GET-KIND | M1 | cbor.go: `for j := 0; v.kind == cborMap && j < len(v.items); j += 2 {` → `for j := 0; j < len(v.items); j += 2 {` | `TestNewPasskeyApproverRefusals` (R9-array-of-six) | KILLED |
| MUT-SAMEKEY-KIND | M1 | cbor.go: `return a.kind == b.kind && a.i == b.i` → `return a.i == b.i` | `TestCBORBoundsAreInclusive` (int-0-and-empty-text keys) | KILLED |
| MUT-SAMEKEY-INT | M1 | cbor.go: `a.kind == b.kind && a.i == b.i && string` → `a.kind == b.kind && string` | `TestCBORAcceptsChromesBytes` (COSE labels 1 and 3 collide (K13)) | KILLED |
| MUT-SAMEKEY-TEXT | M1 | cbor.go: `&& string(a.b) == string(b.b)` → (nothing) | `TestCBORAcceptsChromesBytes` (fmt and attStmt collide (K13)) | KILLED |
| MUT-Q1-EQUIV | M1 | cose.go: `if v.kind != cborMap { // Q1` → `if false { // Q1` | `TestCOSEKeyRefusals` (equivalent) | SURVIVED |
| MUT-Q2 | M1 | cose.go: `if !ok1 \|\| !ok2 \|\| !ok3 \|\| kty.kind != cborInt \|\| alg.kind != cborInt \|\| crv.kind != cborInt { // Q2` → `if false && (!ok1 \|\| !ok2 \|\| !ok3 \|\| kty.kind != cborInt \|\| alg.kind != cborInt \|\| crv.kind != cborInt) { // Q2` | `TestCOSEKeyRefusals` (Q2-no-alg, Q2-alg-text (become ErrCOSEAlg)) | KILLED |
| MUT-Q2-ALGKIND | M1 | cose.go: `\|\| alg.kind != cborInt \|\|` → `\|\|` | `TestCOSEKeyRefusals` (Q2-alg-text) | KILLED |
| MUT-Q3-OFF | M1 | cose.go: `default: // Q3⏎return COSEKey{}, fmt.Errorf("%w (it is %d)", ErrCOSEAlg, k.Alg)⏎` → `default: // Q3⏎` | `TestCOSEKeyRefusals` (Q3-RS256, Q3-ES384) | KILLED |
| MUT-Q3-ACCEPT-RS256 | M1 | cose.go: `case AlgEdDSA:⏎if kty.i != 1` → `case AlgEdDSA, -257:⏎if kty.i != 1` | `TestCOSEKeyRefusals` (Q3-RS256) | KILLED |
| MUT-Q4 | M1 | cose.go: `if kty.i != 2 \|\| crv.i != 1 { // Q4` → `if false { // Q4` | `TestCOSEKeyRefusals` (Q4-*) | KILLED |
| MUT-Q4-CRV | M1 | cose.go: `if kty.i != 2 \|\| crv.i != 1 { // Q4` → `if kty.i != 2 { // Q4` | `TestCOSEKeyRefusals` (Q4-ES256-crv-2) | KILLED |
| MUT-Q5 | M1 | cose.go: `if kty.i != 1 \|\| crv.i != 6 { // Q5` → `if false { // Q5` | `TestCOSEKeyRefusals` (Q5-*) | KILLED |
| MUT-Q5-KTY | M1 | cose.go: `if kty.i != 1 \|\| crv.i != 6 { // Q5` → `if crv.i != 6 { // Q5` | `TestCOSEKeyRefusals` (Q5-EdDSA-kty-2) | KILLED |
| MUT-ALLOW-Y | M1 | cose.go: `allowed[-3] = true` → `allowed[-3] = false` | `TestCOSEKeyRefusals` (the es256 control) | KILLED |
| MUT-Q8 | M1 | cose.go: `if key := v.items[j]; key.kind != cborInt \|\| !allowed[key.i] { // Q8` → `if key := v.items[j]; false && (key.kind != cborInt \|\| !allowed[key.i]) { // Q8` | `TestCOSEKeyRefusals` (Q8-kid, Q8-EdDSA-y) | KILLED |
| MUT-Q8-TEXT-EQUIV | M1 | cose.go: `key.kind != cborInt \|\| !allowed[key.i]` → `!allowed[key.i]` | `TestCOSEKeyRefusals` (equivalent) | SURVIVED |
| MUT-SKIP-Y | M1 | cose.go: `if !allowed[c.n] {⏎continue` → `if false {⏎continue` | `TestCOSEKeyRefusals` (the eddsa control (no -3)) | KILLED |
| MUT-Q6 | M1 | cose.go: `if !ok \|\| b.kind != cborBytes \|\| len(b.b) != 32 { // Q6` → `if false && (!ok \|\| b.kind != cborBytes \|\| len(b.b) != 32) { // Q6` | `TestCOSEKeyRefusals` (Q6-EdDSA-31-byte-x) | KILLED |
| MUT-Q6-KIND | M1 | cose.go: `if !ok \|\| b.kind != cborBytes \|\| len(b.b) != 32 { // Q6` → `if !ok \|\| len(b.b) != 32 { // Q6` | `TestCOSEKeyRefusals` (Q6-EdDSA-x-text) | KILLED |
| MUT-Q6-LEN | M1 | cose.go: `if !ok \|\| b.kind != cborBytes \|\| len(b.b) != 32 { // Q6` → `if !ok \|\| b.kind != cborBytes { // Q6` | `TestCOSEKeyRefusals` (Q6-EdDSA-31-byte-x) | KILLED |
| MUT-Q6-OK-EQUIV | M1 | cose.go: `if !ok \|\| b.kind != cborBytes \|\| len(b.b) != 32 { // Q6` → `if b.kind != cborBytes \|\| len(b.b) != 32 \|\| false && !ok { // Q6` | `TestCOSEKeyRefusals` (equivalent) | SURVIVED |
| MUT-Q7 | M1 | cose.go: `k.point()); err != nil { // Q7` → `k.point()); false && err != nil { // Q7` | `TestCOSEKeyRefusals` (Q7-off-curve) | KILLED |
| MUT-CANON-ED | M1 | cose.go: `if k.Alg == AlgEdDSA {⏎n, kty, crv = 0xa4, 1, 6` → `if false {⏎n, kty, crv = 0xa4, 1, 6` | `TestCOSECanonicalIsChromesBytes` (eddsa) | KILLED |
| MUT-CANON-Y | M1 | cose.go: `if k.Alg == AlgES256 {⏎b = append(append(b, 0x22` → `if false {⏎b = append(append(b, 0x22` | `TestCOSECanonicalIsChromesBytes` (es256) | KILLED |
| MUT-CANON-ALG | M1 | cose.go: `byte(0x20 \| (-1 - k.Alg))` → `byte(0x26)` | `TestCOSECanonicalIsChromesBytes` (eddsa) | KILLED |
| MUT-ALGNAME | M1 | cose.go: `if k.Alg == AlgES256 {⏎return "ES256"` → `if false {⏎return "ES256"` | `TestCOSECanonicalIsChromesBytes` (AlgName) | KILLED |
| MUT-VERIFY-ES | M1 | cose.go: `\|\| !ecdsa.VerifyASN1(pub, h[:], sig) {` → `\|\| false && !ecdsa.VerifyASN1(pub, h[:], sig) {` | `TestCOSEKeyVerifies` (es256 tampered-message) | KILLED |
| MUT-VERIFY-HASH | M1 | cose.go: `h := sha256.Sum256(signed)` → `h := sha256.Sum224(signed)` | `TestCOSEKeyVerifies` (es256 Chrome's signature) | KILLED |
| MUT-ES-LEN-EQUIV | M1 | cose.go: `if len(sig) > 72 \|\| err != nil` → `if err != nil` | `TestCOSEKeyVerifies` (equivalent) | SURVIVED |
| MUT-ES-ERR-EQUIV | M1 | cose.go: `\|\| err != nil \|\| !ecdsa` → `\|\| false && err != nil \|\| !ecdsa` | `TestCOSEKeyVerifies` (equivalent) | SURVIVED |
| MUT-VERIFY-ED | M1 | cose.go: `\|\| !ed25519.Verify(k.X, signed, sig) {` → `\|\| false && !ed25519.Verify(k.X, signed, sig) {` | `TestCOSEKeyVerifies` (eddsa tampered-message) | KILLED |
| MUT-ED-LEN-EQUIV | M1 | cose.go: `if len(sig) != ed25519.SignatureSize \|\| ` → `if ` | `TestCOSEKeyVerifies` (equivalent) | SURVIVED |
| MUT-ED-XLEN-EQUIV | M1 | cose.go: `\|\| len(k.X) != ed25519.PublicKeySize \|\|` → `\|\|` | `TestCOSEKeyVerifies` (equivalent) | SURVIVED |
| MUT-VERIFY-DEFAULT | M1 | cose.go: `default:⏎return ErrCOSEAlg⏎}` → `default:⏎}` | `TestCOSEKeyVerifies` (an unknown algorithm) | KILLED |
| MUT-D1 | M2 | authdata.go: `if len(b) < 37 { // D1` → `if false { // D1` | `TestAuthDataRefusals` (D1-36-bytes (panics)) | KILLED |
| MUT-D2 | M2 | authdata.go: `if a.BS && !a.BE { // D2` → `if false { // D2` | `TestAuthDataRefusals` (D2-BS-without-BE) | KILLED |
| MUT-D3 | M2 | authdata.go: `if (f&0x40 != 0) != registration { // D3` → `if false { // D3` | `TestAuthDataRefusals` (D3-*) | KILLED |
| MUT-UP-BIT | M2 | authdata.go: `a.UP, a.UV, a.BE, a.BS = f&0x01 != 0,` → `a.UP, a.UV, a.BE, a.BS = f&0x02 != 0,` | `TestAuthDataReadsChrome` (UP) | KILLED |
| MUT-UV-BIT | M2 | authdata.go: `f&0x04 != 0, f&0x08 != 0` → `f&0x02 != 0, f&0x08 != 0` | `TestAuthDataReadsChrome` (UV read from bit 1) | KILLED |
| MUT-BE-BIT | M2 | authdata.go: `f&0x08 != 0, f&0x10 != 0` → `f&0x20 != 0, f&0x10 != 0` | `TestRenderPasskeyEnrolmentGolden` (flags 0x08 backup-eligible) | KILLED |
| MUT-BS-BIT | M2 | authdata.go: `f&0x10 != 0⏎` → `f&0x20 != 0⏎` | `TestAuthDataRefusals` (D2-BS-without-BE) | KILLED |
| MUT-D4A | M2 | authdata.go: `if len(rest) < 18 { // D4` → `if false { // D4` | `TestAuthDataRefusals` (D4-inside-aaguid (panics)) | KILLED |
| MUT-D5-LOW | M2 | authdata.go: `if n < 1 \|\| n > MaxCredentialIDBytes { // D5` → `if n > MaxCredentialIDBytes { // D5` | `TestAuthDataRefusals` (D5-0-byte-id) | KILLED |
| MUT-D5-HIGH | M2 | authdata.go: `if n < 1 \|\| n > MaxCredentialIDBytes { // D5` → `if n < 1 { // D5` | `TestAuthDataRefusals` (D5-1024-byte-id) | KILLED |
| MUT-D5-BOUND | M2 | authdata.go: `if n < 1 \|\| n > MaxCredentialIDBytes { // D5` → `if n < 1 \|\| n >= MaxCredentialIDBytes { // D5` | `TestAuthDataRefusals` (1023 accepted) | KILLED |
| MUT-D5-FLOOR16 | M2 | authdata.go: `if n < 1 \|\| n > MaxCredentialIDBytes { // D5` → `if n < 16 \|\| n > MaxCredentialIDBytes { // D5` | `TestAuthDataRefusals` (1 accepted ([r1]: the 16-byte floor)) | KILLED |
| MUT-MAXCRED | M2 | authdata.go: `const MaxCredentialIDBytes = 1023` → `const MaxCredentialIDBytes = 1024` | `TestAuthDataRefusals` (D5-1024-byte-id) | KILLED |
| MUT-D4B | M2 | authdata.go: `if len(rest) < 18+n { // D4` → `if false { // D4` | `TestAuthDataRefusals` (D4-inside-credential-id (panics)) | KILLED |
| MUT-D6 | M2 | authdata.go: `if err != nil { // D6` → `if false { // D6` | `TestAuthDataRefusals` (D6-*) | KILLED |
| MUT-D6-PARSE | M2 | authdata.go: `if err == nil {⏎a.Key, err = ParseCOSEKey(v)` → `if false {⏎a.Key, err = ParseCOSEKey(v)` | `TestAuthDataRefusals` (D6-RS256-key) | KILLED |
| MUT-D7 | M2 | authdata.go: `if err != nil \|\| v.kind != cborMap { // D7` → `if false && (err != nil \|\| v.kind != cborMap) { // D7` | `TestAuthDataRefusals` (D7-*) | KILLED |
| MUT-D7-KIND | M2 | authdata.go: `if err != nil \|\| v.kind != cborMap { // D7` → `if err != nil \|\| false && v.kind != cborMap { // D7` | `TestAuthDataRefusals` (D7-integer-extensions) | KILLED |
| MUT-D7-ERR-EQUIV | M2 | authdata.go: `if err != nil \|\| v.kind != cborMap { // D7` → `if v.kind != cborMap \|\| false && err != nil { // D7` | `TestAuthDataRefusals` (equivalent) | SURVIVED |
| MUT-ED-BIT | M2 | authdata.go: `if f&0x80 != 0 {` → `if f&0x40 != 0 {` | `TestAuthDataReadsChrome` (the registration (AT read as ED)) | KILLED |
| MUT-D8 | M2 | authdata.go: `if len(rest) != 0 { // D8` → `if false { // D8` | `TestAuthDataRefusals` (D8-*) | KILLED |
| MUT-SIGNCOUNT | M2 | authdata.go: `a.SignCount = binary.BigEndian.Uint32(b[33:37])` → `a.SignCount = binary.LittleEndian.Uint32(b[33:37])` | `TestAuthDataReadsChrome` (counters 1, 2, 3) | KILLED |
| MUT-J1-EMPTY-EQUIV | M2 | clientdata.go: `if len(b) == 0 \|\| len(b) > MaxClientDataBytes { // J1` → `if len(b) > MaxClientDataBytes { // J1` | `TestClientDataRefusals` (equivalent) | SURVIVED |
| MUT-J1-MAX | M2 | clientdata.go: `if len(b) == 0 \|\| len(b) > MaxClientDataBytes { // J1` → `if len(b) == 0 { // J1` | `TestClientDataRefusals` (J1-2049-bytes) | KILLED |
| MUT-J1-BOUND | M2 | clientdata.go: `if len(b) == 0 \|\| len(b) > MaxClientDataBytes { // J1` → `if len(b) == 0 \|\| len(b) >= MaxClientDataBytes { // J1` | `TestClientDataRefusals` (the 2048-byte control) | KILLED |
| MUT-J2 | M2 | clientdata.go: `if !utf8.Valid(b) { // J2` → `if false && !utf8.Valid(b) { // J2` | `TestClientDataRefusals` (J2-not-utf8) | KILLED |
| MUT-J3-FIRST | M2 | clientdata.go: `if tok, err := dec.Token(); err != nil \|\| tok != json.Delim('{') { // J3` → `if _, err := dec.Token(); err != nil { // J3` | `TestClientDataRefusals` (J3-array-of-pairs) | KILLED |
| MUT-J3-KEY-EQUIV | M2 | clientdata.go: `if !ok { // J3: not a member name` → `if false && !ok { // J3: not a member name` | `TestClientDataRefusals` (equivalent) | SURVIVED |
| MUT-J4 | M2 | clientdata.go: `if seen[key] { // J4` → `if false { // J4` | `TestClientDataRefusals` (J4-*) | KILLED |
| MUT-J4-MARK | M2 | clientdata.go: `seen[key] = true` → `seen[key] = key == "type"` | `TestClientDataRefusals` (J4-challenge-twice) | KILLED |
| MUT-J5-TYPE | M2 | clientdata.go: `if !isStr { // J5` → `if false { // J5` | `TestClientDataRefusals` (J5-type-number) | KILLED |
| MUT-J5-ASSIGN | M2 | clientdata.go: `*strs[key] = s` → `_ = s` | `TestClientDataReadsChrome` (fields) | KILLED |
| MUT-J6-CROSS | M2 | clientdata.go: `if !isBool { // J6` → `if false && !isBool { // J6` | `TestClientDataRefusals` (J6-crossOrigin-string) | KILLED |
| MUT-J6-CROSS-ASSIGN | M2 | clientdata.go: `c.CrossOrigin = bv` → `c.CrossOrigin = bv && false` | `TestClientDataRefusals` (the crossOrigin true control) | KILLED |
| MUT-J6-TOP | M2 | clientdata.go: `if !isStr { // J6` → `if false { // J6` | `TestClientDataRefusals` (J6-topOrigin-bool) | KILLED |
| MUT-HASTOP | M2 | clientdata.go: `c.HasTopOrigin = true` → `c.HasTopOrigin = false` | `TestClientDataRefusals` (the topOrigin control) | KILLED |
| MUT-J3-CLOSE | M2 | clientdata.go: `if _, err := dec.Token(); err != nil { // J3: the closing brace` → `if false { // J3: the closing brace` | `TestClientDataRefusals` (J3-unterminated) | KILLED |
| MUT-J3-EOF | M2 | clientdata.go: `if _, err := dec.Token(); err != io.EOF { // J3: one object` → `if _, err := dec.Token(); false && err != io.EOF { // J3: one object` | `TestClientDataRefusals` (J3-two-objects, J3-trailing) | KILLED |
| MUT-J5-MISSING | M2 | clientdata.go: `if !seen["type"] \|\| !seen["challenge"] \|\| !seen["origin"] { // J5` → `if false { // J5` | `TestClientDataRefusals` (J5-no-origin) | KILLED |
| MUT-J5-ORIGIN | M2 | clientdata.go: ` \|\| !seen["origin"] { // J5` → ` { // J5` | `TestClientDataRefusals` (J5-no-origin) | KILLED |
| MUT-SKIP-TOK | M2 | clientdata.go: `tok, err := dec.Token()⏎if err != nil { // J3⏎` → `tok, err := dec.Token()⏎if false && err != nil { // J3⏎` | `TestClientDataRefusals` (J3-unknown-member-no-value) | KILLED |
| MUT-J7 | M2 | clientdata.go: `if depth > 4 { // J7` → `if false { // J7` | `TestClientDataRefusals` (J7-depth-5) | KILLED |
| MUT-J7-BOUND | M2 | clientdata.go: `if depth > 4 { // J7` → `if depth > 3 { // J7` | `TestClientDataRefusals` (the depth-4 acceptance control) | KILLED |
| MUT-SKIP-NEST | M2 | clientdata.go: `if err := skipJSON(dec, depth+1); err != nil {` → `if err := skipJSON(dec, depth); err != nil {` | `TestClientDataRefusals` (J7-depth-5) | KILLED |
| MUT-SKIP-CLOSE | M2 | clientdata.go: `if _, err := dec.Token(); err != nil { // J3: the closing delimiter` → `if false { // J3: the closing delimiter` | `TestClientDataRefusals` (J3-unknown-member-unterminated) | KILLED |
| MUT-O1 | M2 | origin.go: `if !ok \|\| (scheme != "https" && scheme != "http") \|\| len(host) > 253 \|\| !originHost.MatchString(host) { // O1` → `if false && (!ok \|\| (scheme != "https" && scheme != "http") \|\| len(host) > 253 \|\| !originHost.MatchString(host)) { // O1` | `TestCeremonyOrigin` (O1 refusals) | KILLED |
| MUT-O1-CUT-EQUIV | M2 | origin.go: `if !ok \|\| (scheme != "https" && scheme != "http") \|\| len(host) > 253 \|\| !originHost.MatchString(host) { // O1` → `if (scheme != "https" && scheme != "http") \|\| len(host) > 253 \|\| !originHost.MatchString(host) \|\| false && !ok { // O1` | `TestCeremonyOrigin` (equivalent) | SURVIVED |
| MUT-O1-SCHEME | M2 | origin.go: `if !ok \|\| (scheme != "https" && scheme != "http") \|\| len(host) > 253 \|\| !originHost.MatchString(host) { // O1` → `if !ok \|\| len(host) > 253 \|\| !originHost.MatchString(host) { // O1` | `TestCeremonyOrigin` (ftp://localhost) | KILLED |
| MUT-O1-LEN | M2 | origin.go: `if !ok \|\| (scheme != "https" && scheme != "http") \|\| len(host) > 253 \|\| !originHost.MatchString(host) { // O1` → `if !ok \|\| (scheme != "https" && scheme != "http") \|\| !originHost.MatchString(host) { // O1` | `TestCeremonyOrigin` (the 254-character host) | KILLED |
| MUT-O1-LEN-BOUND | M2 | origin.go: `if !ok \|\| (scheme != "https" && scheme != "http") \|\| len(host) > 253 \|\| !originHost.MatchString(host) { // O1` → `if !ok \|\| (scheme != "https" && scheme != "http") \|\| len(host) > 252 \|\| !originHost.MatchString(host) { // O1` | `TestCeremonyOrigin` (the 253-character control) | KILLED |
| MUT-O1-GRAMMAR | M2 | origin.go: `if !ok \|\| (scheme != "https" && scheme != "http") \|\| len(host) > 253 \|\| !originHost.MatchString(host) { // O1` → `if !ok \|\| (scheme != "https" && scheme != "http") \|\| len(host) > 253 { // O1` | `TestCeremonyOrigin` (LOCALHOST, x.ts.net/, [::1], u@x.ts.net, a..ts.net) | KILLED |
| MUT-O2 | M2 | origin.go: `if host != "localhost" && !strings.HasSuffix(host, ".ts.net") { // O2` → `if false { // O2` | `TestCeremonyOrigin` (evil.example, 127.0.0.1) | KILLED |
| MUT-O2-SUFFIX | M2 | origin.go: `strings.HasSuffix(host, ".ts.net")` → `strings.HasSuffix(host, "ts.net")` | `TestCeremonyOrigin` (https://ts.net) | KILLED |
| MUT-HTTP-TSNET | M2 | origin.go: `if scheme == "http" && host != "localhost" { // O1` → `if false { // O1` | `TestCeremonyOrigin` (http://x.ts.net) | KILLED |
| MUT-O3 | M2 | origin.go: `if perr != nil \|\| n < 1 \|\| n > 65535 \|\| port != strconv.Itoa(n) { // O3` → `if false && (perr != nil \|\| n < 1 \|\| n > 65535 \|\| port != strconv.Itoa(n)) { // O3` | `TestCeremonyOrigin` (bad ports) | KILLED |
| MUT-O3-PERR-EQUIV | M2 | origin.go: `if perr != nil \|\| n < 1 \|\| n > 65535 \|\| port != strconv.Itoa(n) { // O3` → `if n < 1 \|\| n > 65535 \|\| port != strconv.Itoa(n) \|\| false && perr != nil { // O3` | `TestCeremonyOrigin` (equivalent) | SURVIVED |
| MUT-O3-ZERO | M2 | origin.go: `if perr != nil \|\| n < 1 \|\| n > 65535 \|\| port != strconv.Itoa(n) { // O3` → `if perr != nil \|\| n > 65535 \|\| port != strconv.Itoa(n) { // O3` | `TestCeremonyOrigin` (:0) | KILLED |
| MUT-O3-MAX | M2 | origin.go: `if perr != nil \|\| n < 1 \|\| n > 65535 \|\| port != strconv.Itoa(n) { // O3` → `if perr != nil \|\| n < 1 \|\| port != strconv.Itoa(n) { // O3` | `TestCeremonyOrigin` (:65536) | KILLED |
| MUT-O3-MAX-BOUND | M2 | origin.go: `if perr != nil \|\| n < 1 \|\| n > 65535 \|\| port != strconv.Itoa(n) { // O3` → `if perr != nil \|\| n < 1 \|\| n > 65534 \|\| port != strconv.Itoa(n) { // O3` | `TestCeremonyOrigin` (the :65535 control) | KILLED |
| MUT-O3-CANON | M2 | origin.go: `if perr != nil \|\| n < 1 \|\| n > 65535 \|\| port != strconv.Itoa(n) { // O3` → `if perr != nil \|\| n < 1 \|\| n > 65535 { // O3` | `TestCeremonyOrigin` (:065535, :+8443) | KILLED |
| MUT-O3-DEFAULT | M2 | origin.go: `if (scheme == "https" && n == 443) \|\| (scheme == "http" && n == 80) { // O3` → `if false { // O3` | `TestCeremonyOrigin` (:443, :80) | KILLED |
| MUT-O3-443 | M2 | origin.go: `if (scheme == "https" && n == 443) \|\| (scheme == "http" && n == 80) { // O3` → `if scheme == "http" && n == 80 { // O3` | `TestCeremonyOrigin` (https://x.ts.net:443) | KILLED |
| MUT-O3-80 | M2 | origin.go: `if (scheme == "https" && n == 443) \|\| (scheme == "http" && n == 80) { // O3` → `if scheme == "https" && n == 443 { // O3` | `TestCeremonyOrigin` (http://localhost:80) | KILLED |
| MUT-E2 | M2 | passkey_json.go: `if err := json.Unmarshal(b, &e); err != nil \|\| e.ID == nil \|\| e.RawID == nil \|\| e.Type == nil \|\| e.Response == nil { // E2` → `if err := json.Unmarshal(b, &e); false && err != nil { // E2` | `TestPasskeyEnvelope` (E2-* (panics)) | KILLED |
| MUT-E2-ID | M2 | passkey_json.go: `if err := json.Unmarshal(b, &e); err != nil \|\| e.ID == nil \|\| e.RawID == nil \|\| e.Type == nil \|\| e.Response == nil { // E2` → `if err := json.Unmarshal(b, &e); err != nil \|\| e.RawID == nil \|\| e.Type == nil \|\| e.Response == nil { // E2` | `TestPasskeyEnvelope` (E2-no-id (panics)) | KILLED |
| MUT-E2-TYPE | M2 | passkey_json.go: `if err := json.Unmarshal(b, &e); err != nil \|\| e.ID == nil \|\| e.RawID == nil \|\| e.Type == nil \|\| e.Response == nil { // E2` → `if err := json.Unmarshal(b, &e); err != nil \|\| e.ID == nil \|\| e.RawID == nil \|\| e.Response == nil { // E2` | `TestPasskeyEnvelope` (E2-no-type (panics)) | KILLED |
| MUT-E2-RESPONSE | M2 | passkey_json.go: `if err := json.Unmarshal(b, &e); err != nil \|\| e.ID == nil \|\| e.RawID == nil \|\| e.Type == nil \|\| e.Response == nil { // E2` → `if err := json.Unmarshal(b, &e); err != nil \|\| e.ID == nil \|\| e.RawID == nil \|\| e.Type == nil { // E2` | `TestPasskeyEnvelope` (E2-no-response (panics)) | KILLED |
| MUT-E2-RAWID | M2 | passkey_json.go: `if err := json.Unmarshal(b, &e); err != nil \|\| e.ID == nil \|\| e.RawID == nil \|\| e.Type == nil \|\| e.Response == nil { // E2` → `if err := json.Unmarshal(b, &e); err != nil \|\| e.ID == nil \|\| e.Type == nil \|\| e.Response == nil { // E2` | `TestPasskeyEnvelope` (E2-no-rawId (panics)) | KILLED |
| MUT-E3-TYPE | M2 | passkey_json.go: `if *e.Type != "public-key" \|\| *e.ID != *e.RawID { // E3` → `if *e.ID != *e.RawID { // E3` | `TestPasskeyEnvelope` (E3-type-x) | KILLED |
| MUT-E3-ID | M2 | passkey_json.go: `if *e.Type != "public-key" \|\| *e.ID != *e.RawID { // E3` → `if *e.Type != "public-key" { // E3` | `TestPasskeyEnvelope` (E3-id-not-rawId) | KILLED |
| MUT-B64-NIL | M2 | passkey_json.go: `if s == nil { // E2` → `if false { // E2` | `TestPasskeyEnvelope` (E2-no-clientDataJSON (panics)) | KILLED |
| MUT-E4 | M2 | passkey_json.go: `if err != nil { // E4` → `if false && err != nil { // E4` | `TestPasskeyEnvelope` (E4-*) | KILLED |
| MUT-E4-STRICT | M2 | passkey_json.go: `base64.RawURLEncoding.Strict().DecodeString(*s)` → `base64.RawURLEncoding.DecodeString(*s)` | `TestPasskeyEnvelope` (E4-non-canonical) | KILLED |
| MUT-REG-CDJ | M2 | passkey_json.go: `if r.ClientDataJSON, err = b64Field(e.Response.ClientDataJSON); err != nil {` → `if r.ClientDataJSON, err = b64Field(e.Response.ClientDataJSON); false && err != nil {` | `TestPasskeyEnvelope` (E2-no-clientDataJSON) | KILLED |
| MUT-REG-ATT | M2 | passkey_json.go: `if r.AttestationObject, err = b64Field(e.Response.AttestationObject); err != nil {` → `if r.AttestationObject, err = b64Field(e.Response.AttestationObject); false && err != nil {` | `TestPasskeyEnvelope` (E2-null-attestation) | KILLED |
| MUT-ASR-FIELDS | M2 | passkey_json.go: `if *f.dst, err = b64Field(f.src); err != nil {` → `if *f.dst, err = b64Field(f.src); false && err != nil {` | `TestPasskeyEnvelope` (E2-no-signature) | KILLED |
| MUT-UH-BRANCH | M2 | passkey_json.go: `if e.Response.UserHandle != nil {` → `if false {` | `TestPasskeyEnvelope` (userHandle "bWFyaw") | KILLED |
| MUT-UH-ERR | M2 | passkey_json.go: `if r.UserHandle, err = b64Field(e.Response.UserHandle); err != nil {` → `if r.UserHandle, err = b64Field(e.Response.UserHandle); false && err != nil {` | `TestPasskeyEnvelope` (E4-userHandle) | KILLED |
| MUT-UH-FLAG | M2 | passkey_json.go: `r.HasUserHandle = true` → `r.HasUserHandle = false` | `TestPasskeyEnvelope` (userHandle "bWFyaw") | KILLED |
| MUT-R1 | M3a | passkey.go: `if !ApproverNamePattern.MatchString(in.Name) { // R1` → `if false && !ApproverNamePattern.MatchString(in.Name) { // R1` | `TestNewPasskeyApproverRefusals` (R1-name) | KILLED |
| MUT-R2 | M3a | passkey.go: `rpID, err := ParseCeremonyOrigin(in.Origin) // R2⏎if err != nil {` → `rpID, err := ParseCeremonyOrigin(in.Origin) // R2⏎if false && err != nil {` | `TestNewPasskeyApproverRefusals` (R2-ip-origin) | KILLED |
| MUT-R3 | M3a | passkey.go: `resp, err := ParseRegistrationResponse(in.Response) // R3⏎if err != nil {` → `resp, err := ParseRegistrationResponse(in.Response) // R3⏎if false && err != nil {` | `TestNewPasskeyApproverRefusals` (R3-envelope) | KILLED |
| MUT-R4 | M3a | passkey.go: `cd, err := ParseClientData(resp.ClientDataJSON) // R4⏎if err != nil {` → `cd, err := ParseClientData(resp.ClientDataJSON) // R4⏎if false && err != nil {` | `TestNewPasskeyApproverRefusals` (R4-client-data) | KILLED |
| MUT-R5 | M3a | passkey.go: `if cd.Type != "webauthn.create" { // R5` → `if false && cd.Type != "webauthn.create" { // R5` | `TestNewPasskeyApproverRefusals` (R5-type-get) | KILLED |
| MUT-R6 | M3a | passkey.go: `if cd.Challenge != b64Challenge(in.Challenge[:]) { // R6` → `if false && cd.Challenge != b64Challenge(in.Challenge[:]) { // R6` | `TestNewPasskeyApproverRefusals` (R6-challenge) | KILLED |
| MUT-R7 | M3a | passkey.go: `if cd.Origin != in.Origin { // R7` → `if false && cd.Origin != in.Origin { // R7` | `TestNewPasskeyApproverRefusals` (R7-origin) | KILLED |
| MUT-R8-CROSS | M3a | passkey.go: `if cd.CrossOrigin \|\| cd.HasTopOrigin { // R8` → `if cd.HasTopOrigin { // R8` | `TestNewPasskeyApproverRefusals` (R8-cross-origin) | KILLED |
| MUT-R8-TOP | M3a | passkey.go: `if cd.CrossOrigin \|\| cd.HasTopOrigin { // R8` → `if cd.CrossOrigin { // R8` | `TestNewPasskeyApproverRefusals` (R8-top-origin) | KILLED |
| MUT-R9R10 | M3a | passkey.go: `raw, err := attestedAuthData(resp.AttestationObject) // R9, R10⏎if err != nil {` → `raw, err := attestedAuthData(resp.AttestationObject) // R9, R10⏎if false && err != nil {` | `TestNewPasskeyApproverRefusals` (R10-packed) | KILLED |
| MUT-R11 | M3a | passkey.go: `ad, err := ParseAuthData(raw, true) // R11⏎if err != nil {` → `ad, err := ParseAuthData(raw, true) // R11⏎if false && err != nil {` | `TestNewPasskeyApproverRefusals` (R11-no-AT) | KILLED |
| MUT-R12 | M3a | passkey.go: `if ad.RPIDHash != sha256.Sum256([]byte(rpID)) { // R12` → `if false && ad.RPIDHash != sha256.Sum256([]byte(rpID)) { // R12` | `TestNewPasskeyApproverRefusals` (R12-rp-id) | KILLED |
| MUT-R12-HASH | M3a | passkey.go: `sha256.Sum256([]byte(rpID)) { // R12` → `sha256.Sum256([]byte(in.Origin)) { // R12` | `TestNewPasskeyApproverRefusals` (the control) | KILLED |
| MUT-R13 | M3a | passkey.go: `if !ad.UP { // R13` → `if false && !ad.UP { // R13` | `TestNewPasskeyApproverRefusals` (R13-no-UP) | KILLED |
| MUT-R14 | M3a | passkey.go: `if !ad.UV { // R14` → `if false && !ad.UV { // R14` | `TestNewPasskeyApproverRefusals` (R14-no-UV) | KILLED |
| MUT-R15 | M3a | passkey.go: `if !bytes.Equal(resp.RawID, ad.CredentialID) { // R15` → `if false && !bytes.Equal(resp.RawID, ad.CredentialID) { // R15` | `TestNewPasskeyApproverRefusals` (R15-raw-id) | KILLED |
| MUT-R9-DECODE | M3a | passkey.go: `if err != nil { // R9⏎` → `if false && err != nil { // R9⏎` | `TestNewPasskeyApproverRefusals` (R9-trailing-byte (wraps ErrCBORTrailing)) | KILLED |
| MUT-R9-SHAPE | M3a | passkey.go: `if len(v.items) != 6 \|\| !okF \|\| !okS \|\| !okA \|\| f.kind != cborText \|\| st.kind != cborMap \|\| ad.kind != cborBytes { // R9` → `if false && (len(v.items) != 6 \|\| !okF \|\| !okS \|\| !okA \|\| f.kind != cborText \|\| st.kind != cborMap \|\| ad.kind != cborBytes) { // R9` | `TestNewPasskeyApproverRefusals` (R9-array) | KILLED |
| MUT-R9-LEN | M3a | passkey.go: `if len(v.items) != 6 \|\| !okF \|\| !okS \|\| !okA \|\| f.kind != cborText \|\| st.kind != cborMap \|\| ad.kind != cborBytes { // R9` → `if !okF \|\| !okS \|\| !okA \|\| f.kind != cborText \|\| st.kind != cborMap \|\| ad.kind != cborBytes { // R9` | `TestNewPasskeyApproverRefusals` (R9-extra-key) | KILLED |
| MUT-R9-ADKIND | M3a | passkey.go: `if len(v.items) != 6 \|\| !okF \|\| !okS \|\| !okA \|\| f.kind != cborText \|\| st.kind != cborMap \|\| ad.kind != cborBytes { // R9` → `if len(v.items) != 6 \|\| !okF \|\| !okS \|\| !okA \|\| f.kind != cborText \|\| st.kind != cborMap { // R9` | `TestNewPasskeyApproverRefusals` (R9-authData-text) | KILLED |
| MUT-R10 | M3a | passkey.go: `if string(f.b) != "none" \|\| len(st.items) != 0 { // R10` → `if false { // R10` | `TestNewPasskeyApproverRefusals` (R10-*) | KILLED |
| MUT-R10-FMT | M3a | passkey.go: `if string(f.b) != "none" \|\| len(st.items) != 0 { // R10` → `if len(st.items) != 0 { // R10` | `TestNewPasskeyApproverRefusals` (R10-packed) | KILLED |
| MUT-R10-STMT | M3a | passkey.go: `if string(f.b) != "none" \|\| len(st.items) != 0 { // R10` → `if string(f.b) != "none" { // R10` | `TestNewPasskeyApproverRefusals` (R10-attStmt) | KILLED |
| MUT-COPY-FORMAT | M3a | passkey.go: `Format: FormatWebAuthn, EnrolledAt: in.Now` → `Format: FormatSSH, EnrolledAt: in.Now` | `TestNewPasskeyApproverAcceptsChrome` (es256) | KILLED |
| MUT-COPY-NOW | M3a | passkey.go: `Format: FormatWebAuthn, EnrolledAt: in.Now` → `Format: FormatWebAuthn, EnrolledAt: 0` | `TestNewPasskeyApproverAcceptsChrome` (es256) | KILLED |
| MUT-COPY-RPID | M3a | passkey.go: `RPID: rpID,` → `RPID: "localhost",` | `TestNewPasskeyApproverAcceptsChrome` (tailnet) | KILLED |
| MUT-COPY-AAGUID | M3a | passkey.go: `AAGUID: ad.AAGUID,` → `AAGUID: [16]byte{},` | `TestNewPasskeyApproverAcceptsChrome` (es256) | KILLED |
| MUT-COPY-KEY | M3a | passkey.go: `Key: ad.Key,` → `Key: COSEKey{Alg: ad.Key.Alg},` | `TestNewPasskeyApproverAcceptsChrome` (es256) | KILLED |
| MUT-COPY-BE | M3a | passkey.go: `BackupEligible: ad.BE,` → `BackupEligible: false,` | `TestNewPasskeyApproverAcceptsChrome` (tailnet (synced)) | KILLED |
| MUT-COPY-BS | M3a | passkey.go: `BackupState: ad.BS,` → `BackupState: false,` | `TestNewPasskeyApproverAcceptsChrome` (tailnet (synced)) | KILLED |
| MUT-COPY-COUNT | M3a | passkey.go: `SignCount: ad.SignCount}}` → `SignCount: 1}}` | `TestNewPasskeyApproverAcceptsChrome` (tailnet (counter 7)) | KILLED |
| MUT-CODE-PREFIX | M3a | passkey.go: `"world-enrolment-code/v1\x00"` → `"world-enrolment-code/v2\x00"` | `TestEnrolmentCodeGolden` (the golden) | KILLED |
| MUT-CODE-ONE-INPUT | M3a | passkey.go: `c, a := sha256.Sum256(clientDataJSON), sha256.Sum256(attestationObject)` → `c, a := sha256.Sum256(clientDataJSON), sha256.Sum256(clientDataJSON)` | `TestEnrolmentCodeGolden` (attestationObject changed) | KILLED |
| MUT-CODE-SWAP | M3a | passkey.go: `c, a := sha256.Sum256(clientDataJSON), sha256.Sum256(attestationObject)` → `a, c := sha256.Sum256(clientDataJSON), sha256.Sum256(attestationObject)` | `TestEnrolmentCodeGolden` (the golden) | KILLED |
| MUT-CODE-TRUNC | M3a | passkey.go: `EncodeToString(h[:16])` → `EncodeToString(h[:15])` | `TestEnrolmentCodeGolden` (length 26) | KILLED |
| MUT-GROUP | M3a | passkey.go: `code[15:20] + "-" + code[20:]` → `code[15:21] + "-" + code[21:]` | `TestEnrolmentCodeGolden` (grouping) | KILLED |
| MUT-BACKUP-SYNCED | M3a | passkey.go: `case p.BackupEligible && p.BackupState:` → `case false:` | `TestRenderPasskeyEnrolmentGolden` (flags 0x18) | KILLED |
| MUT-BACKUP-ELIGIBLE | M3a | passkey.go: `case p.BackupEligible:` → `case false:` | `TestRenderPasskeyEnrolmentGolden` (flags 0x08) | KILLED |
| MUT-RENDER-AAGUID | M3a | passkey.go: `hex.EncodeToString(p.AAGUID[:])` → `hex.EncodeToString(p.AAGUID[:0])` | `TestRenderPasskeyEnrolmentGolden` (the golden) | KILLED |
| MUT-RENDER-CRED | M3a | passkey.go: `hex.EncodeToString(cred[:8])` → `hex.EncodeToString(cred[:7])` | `TestRenderPasskeyEnrolmentGolden` (the golden) | KILLED |
| MUT-RENDER-PARAGRAPH | M3a | passkey.go: `backup, p.SignCount, a.Format, a.Name, PasskeyEnrolmentWarning)` → `backup, p.SignCount, a.Format, a.Name, "")` | `TestRenderPasskeyEnrolmentGolden` (the golden) | KILLED |
| MUT-RC4-AAGUID | M3a | passkey.go: `hex.EncodeToString(p.AAGUID[:])` → `hex.EncodeToString(p.AAGUID[1:])` | `TestEnrolmentPromptShowsEveryPersistedField` (Passkey.AAGUID (byte 0)) | KILLED |
| MUT-RC4-ORIGIN | M3a | passkey.go: `p.Key.AlgName(), p.RPID, p.Origin,` → `p.Key.AlgName(), p.RPID, p.RPID,` | `TestEnrolmentPromptShowsEveryPersistedField` (Passkey.Origin) | KILLED |
| MUT-RC4-COUNT | M3a | passkey.go: `backup, p.SignCount, a.Format` → `backup, 1, a.Format` | `TestEnrolmentPromptShowsEveryPersistedField` (Passkey.SignCount) | KILLED |
| MUT-E2W | M3a | approver.go: `if format == FormatWebAuthn { // E2` → `if false { // E2` | `TestNewApproverRefusals` (webauthn, webauthn-empty-key) | KILLED |
| MUT-E2W-CLI | M3a | approver.go: `if format == FormatWebAuthn { // E2` → `if false { // E2` | `TestApproverEnrollRefusals` (webauthn (stderr names --registration)) | KILLED |
| MUT-E2-ERR-EQUIV | M2 | passkey_json.go: `if err := json.Unmarshal(b, &e); err != nil \|\| e.ID == nil \|\| e.RawID == nil \|\| e.Type == nil \|\| e.Response == nil { // E2` → `if _ = json.Unmarshal(b, &e); e.ID == nil \|\| e.RawID == nil \|\| e.Type == nil \|\| e.Response == nil { // E2` | `TestPasskeyEnvelope` (equivalent) | SURVIVED |

**Known-equivalent (named, not counted as kills; each run and SURVIVED as predicted):**
- **MUT-Q1-EQUIV** (M1): equivalent because a non-map has no labels (get reads maps only), so Q2 refuses it with ErrCOSEKey.
- **MUT-Q8-TEXT-EQUIV** (M1): equivalent because a text label has i = 0, and 0 is never an allowed label.
- **MUT-Q6-OK-EQUIV** (M1): equivalent because a missing label is the zero value, whose kind is not bytes.
- **MUT-ES-LEN-EQUIV** (M1): equivalent because no DER P-256 signature exceeds 72 bytes, and VerifyASN1 refuses trailing bytes (long-signature still refused).
- **MUT-ES-ERR-EQUIV** (M1): equivalent because every COSEKey that reaches verify passed Q7, so the parse cannot fail.
- **MUT-ED-LEN-EQUIV** (M1): equivalent because ed25519.Verify returns false for a signature that is not 64 bytes.
- **MUT-ED-XLEN-EQUIV** (M1): equivalent because Q6 makes every EdDSA X 32 bytes.
- **MUT-D7-ERR-EQUIV** (M2): equivalent because a failed decode returns the zero value, which is not a map.
- **MUT-J1-EMPTY-EQUIV** (M2): equivalent because an empty input's first token is io.EOF, which J3 refuses with ErrClientData.
- **MUT-J3-KEY-EQUIV** (M2): equivalent because a failed read leaves the decoder at the same byte, so the unknown-member skip of key "" fails the same way (J3).
- **MUT-O1-CUT-EQUIV** (M2): equivalent because without "://" the scheme is the whole input, which is "http" or "https" only with an empty host, which the grammar refuses.
- **MUT-O3-PERR-EQUIV** (M2): equivalent because a failed Atoi returns 0 (n < 1) or a value out of range (n > 65535).
- **MUT-E2-ERR-EQUIV** (M2): equivalent because json.Unmarshal checks the whole input before decoding, and a type error leaves its pointer field nil, so every error also leaves a required member nil.

**Propagation-only (named, not mutated):** `decodeCBOR`'s and `ParseCOSEKey`'s callers returning the
callee's error unchanged (`ParseAuthData` D6 is mutated; R2–R4, R11 are mutated); `chromeFixture`'s
test-only decoding.

**Design §6's mutant list (copied, rule 3i), and where each lands:** K1–K16, Q1–Q8, D1–D8, J1–J7,
O1–O3, E2–E4, R1–R15 and the M0 carry-forwards — **this sprint** (above). P3w, W1–W15, B7–B8 —
**sprint 2's head** (M3b, Appendix C). S1–S9, the three `EnrolPasskey` pre-checks, C13–C18, A6 —
sprint 2. WN1–WN11, WF1–WF6, the widened F1, L4–L5 — sprint 3.

---

## 5. Files (final `git status --short` after M3a)

| File | Sub-milestone | Change |
|---|---|---|
| `cmd/ailang-worldd/approver.go` | M0 | the `approverDecided` seam |
| `cmd/ailang-worldd/approver_revoke_test.go` | M0 | K8a, K8b |
| `host/sshsign/sshsign_test.go` | M0 | SG1 |
| `.gitattributes` | M1 | `host/approval/testdata/webauthn/* -text` |
| `host/approval/cbor.go`, `cbor_test.go` | M1 | new |
| `host/approval/cose.go`, `cose_test.go` | M1 | new |
| `host/approval/webauthn_fixture_test.go` | M1 | new: the fixture loader, XW0, the test encoder |
| `host/approval/testdata/webauthn/{chrome-es256.json, chrome-eddsa.json, SOURCE.txt, generate.sh}` | M1 | new |
| `host/approval/{authdata,clientdata,origin,passkey_json}.go`, `webauthn_parts_test.go` | M2 | new |
| `host/approval/passkey.go`, `passkey_test.go` | M3a | new |
| `host/approval/approver.go`, `approver_test.go` | M3a | `Approver.Passkey`, `ErrUseRegistration`; PE5, RC4 |
| `cmd/ailang-worldd/approver_enroll_test.go` | M3a | one `want` (§1.4 P3) |
| this plan and `sprint_w-world-approvals-by-key-a2-s1.json` | planner | |

---

## 6. Sprint gates (controller, after the last sub-milestone; each green at the base, §1.1)

In **one** `bash` shell, from the worktree. Every leg `tee`s its full output to a log; a leg's verdict
is read from its log, never from `tail` alone:

```
set -euo pipefail
export AILANG_BIN=$HOME/.pinned-ailang/ailang
L=$HOME/.ailang/state/world-iter379/gates; mkdir -p "$L"
./scripts/verify_ail.sh 2>&1 | tee "$L/verify_ail.log"                     # no .ail change: exactly green
test -z "$(gofmt -l cmd/ailang-worldd host/approval host/sshsign)"
go build ./... && go vet ./... 2>&1 | tee "$L/vet.log"
go test -run '^$' ./... 2>&1 | tee "$L/fence.log"                          # compile fence for every _test.go
go test ./cmd/ailang-worldd/ ./host/approval/... ./host/sshsign/ -count=1 2>&1 | tee "$L/pkgs.log"
go test -race ./host/approval/... ./host/sshsign/ -count=1 2>&1 | tee "$L/race.log"
go test -race ./cmd/ailang-worldd/ -run '^(TestApproverRevoke|TestApproverEnroll)' -count=1 2>&1 | tee "$L/race-cmd.log"
# the repo-wide inventories (design §7, iteration 378 retro (a)), by name, -v:
go test -v -count=1 ./host/broker/ -run '^TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential$' 2>&1 | tee "$L/inv.log"
go test -v -count=1 ./host/daemon/ -run '^(TestDaemonDependencyAllowlist|TestFrozenV1RouteTableMatchesMux|TestSurfaceAppRoutes|TestDaemonConstructorCallerInventory)$' 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./host/boundary/ -run '^(TestWorldBoundaryDependencyAllowlist|TestBareNetHTTPExemptionIsPerGroup|TestEveryCommandDirectoryIsAProtectedGroup)$' 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./host/store/ -run '^(TestProductionContextRoots|TestProductionGoSurface|TestNoNewDeadlineFreeStoreReads|TestReadGettersHonorContext|TestCredentialsMethodsRequireADeadline)$' 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./host/verifygate/ -run '^TestNoParallelWriteForkPackagesOutsideVerifygate$' 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./host/approval/... -run '^(TestApprovalImportsOnlyTheStandardLibrary|TestApprovalTestIsTestOnly)$' 2>&1 | tee -a "$L/inv.log"
for t in TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential TestDaemonDependencyAllowlist TestFrozenV1RouteTableMatchesMux \
  TestSurfaceAppRoutes TestDaemonConstructorCallerInventory TestWorldBoundaryDependencyAllowlist TestBareNetHTTPExemptionIsPerGroup \
  TestEveryCommandDirectoryIsAProtectedGroup TestProductionContextRoots TestProductionGoSurface TestNoNewDeadlineFreeStoreReads \
  TestReadGettersHonorContext TestCredentialsMethodsRequireADeadline TestNoParallelWriteForkPackagesOutsideVerifygate \
  TestApprovalImportsOnlyTheStandardLibrary TestApprovalTestIsTestOnly; do
  grep -qE -- "^--- PASS: $t " "$L/inv.log" || { echo "$t did not PASS"; exit 1; }; done
if grep -E '^\s*--- (SKIP|FAIL)' "$L/inv.log"; then exit 1; fi
# AFTER the controller commits (the tracked-binary gate reads committed blobs; this leg runs go test ./...):
VERIFY_GO_LEGS=plain ./scripts/verify_go.sh 2>&1 | tee "$L/verify_go.log"
git status --short                                                         # only the §5 files
```

**Known base reds outside this sprint**: none at `03d78b6a` (§1.1). The local `-race` leg of
`host/broker`, `host/daemon`, `host/verifygate`, `TestWorktreeSnapshotRefusesAConfiguredFilter` and
`TestMissionConfigDefaultFlowIntegration/invalid_ledger` (§1.2 C26) are load-sensitive (rule 3e): a
timeout there is re-run once, alone, not chased. Any other red, any race report,
or any red in a touched package is a regression. CI's `go-race-verify` is the race verdict.

**Must stay green, unedited**: every a0a1 and A1b test, except the named edits:
`approver_test.go`'s two webauthn arms (PE5) and RC4; `approver_enroll_test.go`'s webauthn `want`
(§1.4 P3). In particular `TestSignedApprovalFieldsAreAllBound`, `TestApproveBySignatureOrder`, every
`TestSSHSIG*`, `TestApprover*`, `TestApproverRevoke*`, `TestRenderEnrolmentGolden`,
`TestNewApproverAcceptsTheRealFixtures`, `TestApprovalTestIsTestOnly`, `TestCappedBufferKeepsItsBound`,
`TestSignHandsTheSignerAScrubbedEnvironment`, and every §7 inventory.

---

## 7. Executor directive notes

- **No git writes.** Never `git add`, `commit`, `stash`, `checkout`, `reset`, `restore`, `rebase`, `push`,
  `worktree` or `apply`. `patch -p1` (§2) is the only writer; read-only `git status`/`git diff` are fine.
- **Runs.** Run 1 = M0, M1; Run 2 = M2, M3a. Each run reads §0, §2 and its own §3 sections, then
  applies its patches in order.
- **Expected numstat after each sub-milestone, measured AFTER the controller's commit of it** (an
  uncommitted tree hides the new untracked files from numstat):
  ```
  bash -c 'cd /Users/voightkampff/dev/sunholo-data/.wt-world-iter379-row201-a2 && git diff --numstat 03d78b6a16070a0fc91e8f6eb87d08e820df0cad HEAD -- cmd host | awk "{s+=\$1+\$2} END {print s}"'
  ```
  After the controller's commit of M0 **75**, M1 **905**, M2 **1510**, M3a **1958**. **Any other value
  means a patch was not applied as written: stop and report.**
- **Snapshot after every sub-milestone, cumulatively**:
  ```
  bash -c 'set -euo pipefail; M=M0; W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter379-row201-a2
    D=$HOME/.ailang/state/world-iter379/snap/$M; rm -rf "$D"; mkdir -p "$D"; cd "$W"
    git status --porcelain --untracked-files=all | while read -r st f; do
      case "$f" in design_docs/planned/w-world-approvals-by-key-a2-sprint1-plan.md|design_docs/planned/sprint_w-world-approvals-by-key-a2-s1.json) continue;; esac
      mkdir -p "$D/$(dirname "$f")"; cp -p "$f" "$D/$f"; done
    (cd "$D" && find . -type f ! -name MANIFEST.sha256 | sort | xargs shasum -a 256) > "$D/MANIFEST.sha256"
    cat "$D/MANIFEST.sha256"'
  ```
  (`M` = `M0` … `M3a`.) No file is deleted, so the porcelain status has only `M` and `??`; stop on any
  other status. The controller rebuilds one commit per sub-milestone (iteration 377's `commit_ms.sh`
  with `WT` and the snapshot root pointed here).
- **Stop conditions.** A `patch --dry-run` failure; a numstat other than the expected value; a red
  earlier-milestone test; a `--- SKIP` in any `-v` run; the CB4 fuzz leg failing or writing a corpus
  under `host/approval/testdata/fuzz`; any need to touch a file outside §5; running Chrome, Node or any
  browser (sprint 3); any reference to `~/.ssh`, `~/.ailang/world` or a real `/dev/tty`. Stop and
  report; do not work around.
- **Timing.** Every new `host/approval` test runs in under 1 s; the fuzz leg takes 30 s; M0's cmd tests
  ~1 s; the compile fence ~20 s.
- **Shell.** Wrap loops, globs and `$VAR` expansions in `bash -c` or a `bash <<'EOF'` heredoc (zsh does
  not word-split). Export `AILANG_BIN` in the same shell as the Go tests. `tee` every gate's output to a
  file; never pipe a gate's only copy through `tail`. Poll any background command's artifact in a
  bounded `date +%s` loop in the same turn; never end a turn announcing a wait.

---

## 8. Risks

- **R1 — Size.** 1,958 of 2,000, measured; the per-step expected numstat makes drift visible at once.
  The assertion half is the next sprint's head, so sprint 2's own estimate (≈ 1,600) grows by ≈ 450:
  its planner must apply §9's sprint-2 fallback early.
- **R2 — Strictness against other browsers (design R3).** K4 (shortest heads), J4 (duplicate members),
  the strict base64url envelope and `rawId == authData`'s credential ID are measured only on Chrome
  157's virtual authenticator. The attended exit run (design §10) is the check; a fix is a narrow
  change to one refusal.
- **R3 — The fixtures are one recording.** XW0 pins them; P2 (no full `assertion2`) is routed. Sprint
  3's committed driver re-records, and its XB1 is the end-to-end check.
- **R4 — A transitional message.** Between sprints 1 and 3, `approver enroll --format webauthn` names
  flags that do not exist yet (§1.4 P3). Nothing ships between the sprints of one row.
- **R5 — Haiku and patches.** The executor transcribes nothing; a skipped patch fails `--dry-run` or the
  numstat; a hand "fix" after a red test is forbidden by §7.

---

## 9. Acceptance criteria (the judge scores against these)

1. M0: K8's two post-write failures are `exitFatal` with K7's line already printed, each pinned by its
   own test (K8a, K8b); a signer that prints and fails is refused with nil bytes (SG1). The three
   mutants survive at the base (C13) and die with M0.
2. The CBOR subset refuses exactly K1–K16 with their typed errors, accepts each bound inclusively, and
   decodes Chrome's attestation objects and COSE keys whole (CB1–CB3); the fuzz leg ran 30 s with no
   failure (CB4).
3. COSE keys: only ES256 (kty 2, crv 1, on-curve) and EdDSA (kty 1, crv 6) with exactly their labels;
   `Canonical()` is Chrome's bytes and order-independent; Chrome's signatures verify and tampered ones
   do not (CK1–CK4).
4. Authenticator data, client data, origin and envelope refuse exactly D1–D8, J1–J7, O1–O3 (with r2's
   default ports refused, each naming the rule) and E2–E4, and read Chrome's bytes as recorded
   (AD1–AD3, CD1–CD2, OR1, EN1); credential IDs are 1–1023 bytes **[r1]**.
5. `NewPasskeyApprover` refuses each of R1–R15 with its typed error and the zero `Approver`, accepts both
   Chrome registrations with every field copied, and binds origin, rpID, challenge, UP and **UV**
   (PE1, PE2); the enrolment code is the 26-character, 128-bit fingerprint pinned on the recording,
   grouped 5-5-5-5-6 (PE4); the y/N text is golden and shows every persisted field or excludes it
   with a reason (PE3, RC4); `NewApprover` points webauthn at the registration path (PE5).
6. Each §4 mutant, applied alone, makes its named test fail; the 13 named equivalents survive for the
   stated reasons.
7. Every §6 gate passes, including every §7 inventory by name and `VERIFY_GO_LEGS=plain
   ./scripts/verify_go.sh` after the commit; no test skips; `git status --short` lists only §5.
8. The numstat after each sub-milestone's commit equals §7's value; the sprint total is ≤ 2,000.
9. No real key, browser, authenticator, store, `~/.ssh` access, mint or publish anywhere in the sprint
   (D-WORLD-109); `host/approval` stays pure and standard-library-only.

---

## Appendix C — carried to sprint 2's head (M3b, the §9 fallback's second clause)

Not built here. Sprint 2's planner takes these **before** M4, unchanged from the design except §1.4 P2
and P4:
- **`host/approval/passkey_assert.go`**: `AssertionInput`, `ApproveByAssertion` (P1, P2–P5 via the
  shared `checkApprover(admit, a, format)` with **P3w** `Passkey != nil`, W1–W15; W6 against
  `MessageHash(StoredDigest, Challenge, ExpiresAt)`; W14 over `authenticatorData ‖
  SHA-256(clientDataJSON)` with `COSEKey.verify`; W15 "both zero passes, else strictly greater",
  `ErrSignCountRegression`, `ErrUserHandle`, `ErrBackupEligibility`). `ApproveBySignature` moves to
  `checkApprover` (its P2–P5 mutants are then killed by both laws' tests).
- **`signed.go`**: `credentialID`, `signCount`, `PasskeyUse()` (B7, B8); `Use()`'s key is
  `Passkey.Key.Canonical()`.
- **`host/approval/approvaltest/passkey.go`** (`NewPasskey(alg, seed)`, `Register`, `Assert`, knobs) in an
  external `approval_test` package file (§1.4 P4); `TestApprovalTestIsTestOnly` already guards it.
- **Tests**: PA1–PA4, RC1 (edited), TP1, XW1, and XW2's assertion part. PA3's last arm uses
  `assertion2`'s recorded counter (§1.4 P2). The recorded `assertion` (counter 2) approves the
  existing fixture request (`fixtureDigest`, `fixtureChallenge`, `fixtureExpiry`) against PE2's
  passkey (stored counter 1).
- **What sprint 1 leaves for it**: `ParseAssertionResponse` (EN1-tested), `ParseClientData`,
  `ParseAuthData(…, false)`, `COSEKey.verify` (CK2-tested on Chrome's assertions), the fixtures with
  both assertions' bytes, `Approver.Passkey`, `b64Challenge`, and the error values R5–R15 share with
  W5–W14 (`ErrClientDataType`, `ErrChallenge`, `ErrOriginMismatch`, `ErrCrossOrigin`, `ErrRPIDHash`,
  `ErrUserPresence`, `ErrUserVerification`, `ErrCredentialIDMismatch`).
- **Size**: ≈ 450 lines (estimate from this sprint's measured density); with design sprint 2's ≈ 1,600
  the sprint-2 planner must apply its own fallback (AU3 and A6 to the head of sprint 3) early.

---

## Appendix P — the patches (apply with §2; generated from the compiled scratch stages; do not edit)

### Patch M0 — the carry-forwards K8a, K8b, SG1

````diff M0
--- a/cmd/ailang-worldd/approver.go
+++ b/cmd/ailang-worldd/approver.go
@@ -374,6 +374,10 @@
 // approverRevoke is the revoke write, a test seam (plan A-D4): CA8 replaces it to inject a race or
 // a failure; production never reassigns it.
 var approverRevoke = (*store.Credentials).RevokeApprover
+
+// approverDecided is K8's listing read, a test seam (row 201 A2 sprint 1 M0, K8b): a test
+// replaces it to inject a failure; production never reassigns it.
+var approverDecided = (*store.Credentials).DecidedMintRequests
 
 // runApproverRevoke is `approver revoke <name>` (design §2.10 K1-K7).
 func runApproverRevoke(args []string, stdout, stderr io.Writer, env sessionEnv) int {
@@ -448,7 +452,7 @@
 		fmt.Fprintf(stderr, "%srevoked, but cannot read the credentials store %s to list what it approved: %v; fix: ailang-worldd approve --list --decided\n", prefix, credPath, err)
 		return exitFatal
 	}
-	rows, err := ro.DecidedMintRequests(list, row.Format+":"+name, 0)
+	rows, err := approverDecided(ro, list, row.Format+":"+name, 0)
 	_ = ro.Close()
 	if err != nil {
 		fmt.Fprintf(stderr, "%srevoked, but cannot list what it approved: %v; fix: ailang-worldd approve --list --decided\n", prefix, err)
--- a/cmd/ailang-worldd/approver_revoke_test.go
+++ b/cmd/ailang-worldd/approver_revoke_test.go
@@ -10,6 +10,7 @@
 	"errors"
 	"fmt"
 	"io"
+	"os"
 	"reflect"
 	"strings"
 	"sync"
@@ -332,5 +333,57 @@
 	}
 	if rows := approverRows(t, creds); !reflect.DeepEqual(rows, before) {
 		t.Fatalf("no timer: rows changed: %+v", rows)
+	}
+}
+
+// ---- row 201 A2 sprint 1 M0 (K8-ERRCODE): K8's read failing after a committed revoke is
+// exitFatal (a failure after the write), with the revoke already printed.
+
+// assertRevokedThenFatal runs a confirmed revoke and requires K7's line, exitFatal and want.
+func assertRevokedThenFatal(t *testing.T, db, want string) {
+	t.Helper()
+	term := &lockedTerm{in: strings.NewReader("y\n")}
+	var out, errw bytes.Buffer
+	code := runApprover([]string{"revoke", "mark-mac", "--db", db}, &out, &errw, approveEnv(term, nil, atTime(approveT0+100)))
+	if code != exitFatal {
+		t.Fatalf("code %d, want exitFatal; stderr %q", code, errw.String())
 	}
+	if out.String() != "revoked approver mark-mac\n" {
+		t.Fatalf("stdout %q, want only K7's line", out.String())
+	}
+	if !strings.Contains(errw.String(), want) || !strings.Contains(errw.String(), "fix: ailang-worldd approve --list --decided") {
+		t.Fatalf("stderr %q lacks %q and the fix", errw.String(), want)
+	}
 }
+
+// K8a: the real revoke commits, then the credentials file is moved away, so K8's read-only
+// reopen fails.
+func TestApproverRevokeListingReopenFailureIsFatal(t *testing.T) {
+	db, creds, _ := revokeRig(t)
+	orig := approverRevoke
+	t.Cleanup(func() { approverRevoke = orig })
+	approverRevoke = func(c *store.Credentials, ctx context.Context, name string, at int64) error {
+		if err := orig(c, ctx, name, at); err != nil {
+			return err
+		}
+		return os.Rename(creds, creds+".away")
+	}
+	assertRevokedThenFatal(t, db, "revoked, but cannot read the credentials store")
+	if _, err := os.Stat(creds + ".away"); err != nil {
+		t.Fatalf("the seam did not run: %v", err)
+	}
+}
+
+// K8b: K8's listing query fails after the revoke committed.
+func TestApproverRevokeListingQueryFailureIsFatal(t *testing.T) {
+	db, creds, _ := revokeRig(t)
+	orig := approverDecided
+	t.Cleanup(func() { approverDecided = orig })
+	approverDecided = func(*store.Credentials, context.Context, string, int) ([]store.MintRequestRow, error) {
+		return nil, errors.New("disk on fire")
+	}
+	assertRevokedThenFatal(t, db, "revoked, but cannot list what it approved: disk on fire")
+	if rows := approverRows(t, creds); len(rows) != 1 || rows[0].RevokedAt != approveT0+100 {
+		t.Fatalf("rows %+v, want mark-mac revoked at %d", rows, approveT0+100)
+	}
+}
--- a/host/sshsign/sshsign_test.go
+++ b/host/sshsign/sshsign_test.go
@@ -5,8 +5,10 @@
 // ssh-keygen and no key: the program is a temp shell script.
 
 import (
+	"errors"
 	"io"
 	"os"
+	"os/exec"
 	"path/filepath"
 	"strings"
 	"testing"
@@ -61,3 +63,17 @@
 		}
 	}
 }
+
+// SG1 (row 201 A2 sprint 1 M0): a signer that prints something and then fails is refused (G1),
+// and its output is never returned as a signature.
+func TestSignRefusesAFailedSignerEvenWithOutput(t *testing.T) {
+	bin := filepath.Join(t.TempDir(), "signer")
+	if err := os.WriteFile(bin, []byte("#!/bin/sh\ncat > /dev/null\necho partial-signature\nexit 3\n"), 0o700); err != nil {
+		t.Fatal(err)
+	}
+	out, err := Sign(t.Context(), bin, "k", []byte("msg"), io.Discard)
+	var exit *exec.ExitError
+	if out != nil || !errors.As(err, &exit) || exit.ExitCode() != 3 {
+		t.Fatalf("Sign = (%q, %v), want nil bytes and an error wrapping exit status 3", out, err)
+	}
+}
````

### Patch M1 — CBOR and COSE, and the recorded fixtures

````diff M1
--- a/.gitattributes
+++ b/.gitattributes
@@ -1,2 +1,3 @@
 host/approval/testdata/sshsig/* -text
 host/approval/testdata/approvers/* -text
+host/approval/testdata/webauthn/* -text
--- a/host/approval/cbor.go
+++ b/host/approval/cbor.go
@@ -0,0 +1,195 @@
+package approval
+
+// Row 201 A2 (design_docs/planned/w-world-approvals-by-key-a2.md §2.3): the bounded CBOR subset a
+// WebAuthn registration needs (RFC 8949 major types 0-5 and the simple values false/true/null),
+// read from a byte slice with fixed bounds. Every refusal is its own branch (K1-K16). No I/O.
+
+import (
+	"encoding/binary"
+	"errors"
+	"math"
+	"unicode/utf8"
+)
+
+var (
+	ErrCBORTruncated    = errors.New("approval: cbor: truncated")
+	ErrCBORMalformed    = errors.New("approval: cbor: malformed")
+	ErrCBORIndefinite   = errors.New("approval: cbor: indefinite length is not accepted")
+	ErrCBORNonCanonical = errors.New("approval: cbor: an argument is not in its shortest form")
+	ErrCBORTooLarge     = errors.New("approval: cbor: over a bound")
+	ErrCBORTooDeep      = errors.New("approval: cbor: nested too deep")
+	ErrCBORUnsupported  = errors.New("approval: cbor: a tag, float, simple value or map key type outside the subset")
+	ErrCBORDuplicateKey = errors.New("approval: cbor: a duplicate map key")
+	ErrCBORTrailing     = errors.New("approval: cbor: trailing bytes")
+)
+
+type cborKind uint8
+
+const (
+	cborInt cborKind = iota + 1
+	cborBytes
+	cborText
+	cborArray
+	cborMap
+	cborBool
+	cborNull
+)
+
+// cborValue is one decoded item. A map keeps its entries in the order they were read.
+type cborValue struct {
+	kind  cborKind
+	i     int64
+	b     []byte // bytes, or text's UTF-8
+	t     bool
+	items []cborValue // array elements, or map keys and values alternating
+}
+
+// get returns the value a map holds under an integer or text key.
+func (v cborValue) get(key cborValue) (cborValue, bool) {
+	for j := 0; v.kind == cborMap && j < len(v.items); j += 2 {
+		if sameKey(v.items[j], key) {
+			return v.items[j+1], true
+		}
+	}
+	return cborValue{}, false
+}
+
+func sameKey(a, b cborValue) bool {
+	return a.kind == b.kind && a.i == b.i && string(a.b) == string(b.b)
+}
+
+type cborLimits struct{ maxDepth, maxItemsPerContainer, maxTotalItems, maxStringBytes int }
+
+var attestationLimits = cborLimits{maxDepth: 4, maxItemsPerContainer: 16, maxTotalItems: 64, maxStringBytes: 2048}
+
+// decodeCBOR reads one item and returns it and the bytes after it.
+func decodeCBOR(b []byte, lim cborLimits) (cborValue, []byte, error) {
+	d := &cborDecoder{b: b, lim: lim}
+	v, err := d.item(1)
+	if err != nil {
+		return cborValue{}, nil, err
+	}
+	return v, d.b, nil
+}
+
+// decodeCBORExact reads one item and refuses trailing bytes.
+func decodeCBORExact(b []byte, lim cborLimits) (cborValue, error) {
+	v, rest, err := decodeCBOR(b, lim)
+	if err != nil {
+		return cborValue{}, err
+	}
+	if len(rest) != 0 { // K16
+		return cborValue{}, ErrCBORTrailing
+	}
+	return v, nil
+}
+
+type cborDecoder struct {
+	b     []byte
+	lim   cborLimits
+	total int
+}
+
+// head reads an initial byte and its argument.
+func (d *cborDecoder) head() (major byte, arg uint64, err error) {
+	if len(d.b) == 0 { // K1
+		return 0, 0, ErrCBORTruncated
+	}
+	major, ai := d.b[0]>>5, d.b[0]&0x1f
+	switch {
+	case major == 6: // K9
+		return 0, 0, ErrCBORUnsupported
+	case major == 7:
+		if ai < 20 || ai > 22 { // K10
+			return 0, 0, ErrCBORUnsupported
+		}
+		d.b = d.b[1:]
+		return major, uint64(ai), nil
+	case ai < 24:
+		d.b = d.b[1:]
+		return major, uint64(ai), nil
+	case ai == 31: // K3
+		return 0, 0, ErrCBORIndefinite
+	case ai > 27: // K2
+		return 0, 0, ErrCBORMalformed
+	}
+	n := 1 << (ai - 24)
+	if len(d.b) < 1+n { // K1
+		return 0, 0, ErrCBORTruncated
+	}
+	var buf [8]byte
+	copy(buf[8-n:], d.b[1:1+n])
+	arg = binary.BigEndian.Uint64(buf[:])
+	if arg < [4]uint64{24, 1 << 8, 1 << 16, 1 << 32}[ai-24] { // K4
+		return 0, 0, ErrCBORNonCanonical
+	}
+	d.b = d.b[1+n:]
+	return major, arg, nil
+}
+
+func (d *cborDecoder) item(depth int) (cborValue, error) {
+	if depth > d.lim.maxDepth { // K8
+		return cborValue{}, ErrCBORTooDeep
+	}
+	if d.total++; d.total > d.lim.maxTotalItems { // K15
+		return cborValue{}, ErrCBORTooLarge
+	}
+	major, arg, err := d.head()
+	if err != nil {
+		return cborValue{}, err
+	}
+	switch major {
+	case 0, 1:
+		if arg > math.MaxInt64 { // K14
+			return cborValue{}, ErrCBORTooLarge
+		}
+		if major == 1 {
+			return cborValue{kind: cborInt, i: -1 - int64(arg)}, nil
+		}
+		return cborValue{kind: cborInt, i: int64(arg)}, nil
+	case 2, 3:
+		if arg > uint64(d.lim.maxStringBytes) { // K6
+			return cborValue{}, ErrCBORTooLarge
+		}
+		if arg > uint64(len(d.b)) { // K5
+			return cborValue{}, ErrCBORTruncated
+		}
+		s := d.b[:arg:arg]
+		d.b = d.b[arg:]
+		if major == 3 && !utf8.Valid(s) { // K11
+			return cborValue{}, ErrCBORMalformed
+		}
+		return cborValue{kind: cborBytes + cborKind(major-2), b: s}, nil
+	case 4, 5:
+		if arg > uint64(d.lim.maxItemsPerContainer) { // K7
+			return cborValue{}, ErrCBORTooLarge
+		}
+		v := cborValue{kind: cborArray + cborKind(major-4)}
+		for range int(arg) {
+			e, err := d.item(depth + 1)
+			if err != nil {
+				return cborValue{}, err
+			}
+			if major == 5 {
+				if e.kind != cborInt && e.kind != cborText { // K12
+					return cborValue{}, ErrCBORUnsupported
+				}
+				if _, dup := v.get(e); dup { // K13
+					return cborValue{}, ErrCBORDuplicateKey
+				}
+				val, err := d.item(depth + 1)
+				if err != nil {
+					return cborValue{}, err
+				}
+				v.items = append(v.items, e, val)
+				continue
+			}
+			v.items = append(v.items, e)
+		}
+		return v, nil
+	}
+	if arg == 22 {
+		return cborValue{kind: cborNull}, nil
+	}
+	return cborValue{kind: cborBool, t: arg == 21}, nil
+}
--- a/host/approval/cbor_test.go
+++ b/host/approval/cbor_test.go
@@ -0,0 +1,142 @@
+package approval
+
+// Row 201 A2 sprint 1 M1 (design §6, CB1-CB4): the bounded CBOR subset.
+
+import (
+	"bytes"
+	"errors"
+	"testing"
+)
+
+func nested(depth int) []byte { return append(bytes.Repeat([]byte{0x81}, depth-1), 0x80) }
+
+// intMap is a map of n entries i -> i.
+func intMap(n int) cbMap {
+	m := cbMap{}
+	for i := range n {
+		m = append(m, i, i)
+	}
+	return m
+}
+
+// CB1: one arm per refusal K1-K16, each errors.Is its own error.
+func TestCBORRefusals(t *testing.T) {
+	arms := []struct {
+		name  string
+		in    []byte
+		exact bool
+		want  error
+	}{
+		{"K1-head-argument", []byte{0x18}, false, ErrCBORTruncated},
+		{"K1-empty", nil, false, ErrCBORTruncated},
+		{"K1-in-a-map", []byte{0xa1, 0x01}, false, ErrCBORTruncated},
+		{"K2-reserved", []byte{0x1c}, false, ErrCBORMalformed},
+		{"K2-reserved-30", []byte{0x5e}, false, ErrCBORMalformed},
+		{"K3-indefinite-bytes", []byte{0x5f, 0x41, 0x00, 0xff}, false, ErrCBORIndefinite},
+		{"K3-indefinite-map", []byte{0xbf, 0xff}, false, ErrCBORIndefinite},
+		{"K4-one-byte", []byte{0x18, 0x05}, false, ErrCBORNonCanonical},
+		{"K4-one-byte-23", []byte{0x18, 0x17}, false, ErrCBORNonCanonical},
+		{"K4-two-byte", []byte{0x19, 0x00, 0xff}, false, ErrCBORNonCanonical},
+		{"K4-four-byte", []byte{0x1a, 0x00, 0x00, 0xff, 0xff}, false, ErrCBORNonCanonical},
+		{"K4-eight-byte", []byte{0x1b, 0, 0, 0, 0, 0xff, 0xff, 0xff, 0xff}, false, ErrCBORNonCanonical},
+		{"K4-length", []byte{0x58, 0x01, 0xaa}, false, ErrCBORNonCanonical},
+		{"K5-bytes", []byte{0x58, 0x20, 1, 2, 3}, false, ErrCBORTruncated},
+		{"K5-text", []byte{0x62, 'a'}, false, ErrCBORTruncated},
+		{"K6-2049-bytes", cb(make([]byte, 2049)), false, ErrCBORTooLarge},
+		{"K6-2049-text", cb(string(make([]byte, 2049))), false, ErrCBORTooLarge},
+		{"K7-17-entry-map", cb(intMap(17)), false, ErrCBORTooLarge},
+		{"K7-17-element-array", cb(cbArr{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}), false, ErrCBORTooLarge},
+		{"K8-five-nested", nested(5), false, ErrCBORTooDeep},
+		{"K9-tag", []byte{0xc0, 0x00}, false, ErrCBORUnsupported},
+		{"K10-half-float", []byte{0xf9, 0x00, 0x00}, false, ErrCBORUnsupported},
+		{"K10-undefined", []byte{0xf7}, false, ErrCBORUnsupported},
+		{"K10-simple-19", []byte{0xf3}, false, ErrCBORUnsupported},
+		{"K10-break", []byte{0xff}, false, ErrCBORUnsupported},
+		{"K11-utf8", []byte{0x61, 0xff}, false, ErrCBORMalformed},
+		{"K12-bytes-key", cb(cbMap{[]byte{1}, 1}), false, ErrCBORUnsupported},
+		{"K12-bool-key", cb(cbMap{true, 1}), false, ErrCBORUnsupported},
+		{"K13-duplicate-int", []byte{0xa2, 0x01, 0x01, 0x01, 0x02}, false, ErrCBORDuplicateKey},
+		{"K13-duplicate-text", cb(cbMap{"a", 1, "a", 2}), false, ErrCBORDuplicateKey},
+		{"K14-uint-over-int64", []byte{0x1b, 0x80, 0, 0, 0, 0, 0, 0, 0}, false, ErrCBORTooLarge},
+		{"K14-uint-max", []byte{0x1b, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff}, false, ErrCBORTooLarge},
+		{"K14-negative", []byte{0x3b, 0x80, 0, 0, 0, 0, 0, 0, 0}, false, ErrCBORTooLarge},
+		{"K15-65-items", cb(cbArr{intMap(16), intMap(15)}), false, ErrCBORTooLarge},
+		{"K16-trailing", []byte{0x00, 0x00}, true, ErrCBORTrailing},
+	}
+	for _, a := range arms {
+		t.Run(a.name, func(t *testing.T) {
+			var err error
+			if a.exact {
+				_, err = decodeCBORExact(a.in, attestationLimits)
+			} else {
+				_, _, err = decodeCBOR(a.in, attestationLimits)
+			}
+			if !errors.Is(err, a.want) {
+				t.Fatalf("decode(% x) = %v, want %v", a.in, err, a.want)
+			}
+		})
+	}
+	if _, rest, err := decodeCBOR([]byte{0x00, 0x00}, attestationLimits); err != nil || len(rest) != 1 {
+		t.Fatalf("decodeCBOR leaves trailing bytes to its caller: rest %d, %v", len(rest), err)
+	}
+}
+
+// CB2: every bound is inclusive, and the subset's values decode to what they encode.
+func TestCBORBoundsAreInclusive(t *testing.T) {
+	accepted := map[string][]byte{
+		"2048-bytes":   cb(make([]byte, 2048)),
+		"2048-text":    cb(string(make([]byte, 2048))),
+		"16-entries":   cb(intMap(16)),
+		"depth-4":      nested(4),
+		"64-items":     cb(cbArr{intMap(16), intMap(14), 0}),
+		"0-and-empty":  cb(cbMap{0, 1, "", 2}),
+		"int64-limits": {0x83, 0x1b, 0x7f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x3b, 0x7f, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x18, 0x18},
+	}
+	for name, in := range accepted {
+		if _, err := decodeCBORExact(in, attestationLimits); err != nil {
+			t.Errorf("%s: %v, want accepted", name, err)
+		}
+	}
+	v, err := decodeCBORExact([]byte{0x86, 0x20, 0x17, 0x43, 1, 2, 3, 0x62, 'h', 'i', 0xf5, 0xf6}, attestationLimits)
+	if err != nil || len(v.items) != 6 {
+		t.Fatalf("decode: %v, %+v", err, v)
+	}
+	it := v.items
+	if it[0].kind != cborInt || it[0].i != -1 || it[1].i != 23 || it[2].kind != cborBytes || !bytes.Equal(it[2].b, []byte{1, 2, 3}) ||
+		it[3].kind != cborText || string(it[3].b) != "hi" || it[4].kind != cborBool || !it[4].t || it[5].kind != cborNull {
+		t.Fatalf("decoded %+v", it)
+	}
+	if f, err := decodeCBORExact([]byte{0xf4}, attestationLimits); err != nil || f.kind != cborBool || f.t {
+		t.Fatalf("false decoded as %+v, %v", f, err)
+	}
+}
+
+// CB3: Chrome's attestationObjects and COSE keys decode with nothing left over.
+func TestCBORAcceptsChromesBytes(t *testing.T) {
+	for alg, keyLen := range map[string]int{"es256": 77, "eddsa": 42} {
+		c := chromeFixture(t, alg)
+		v, err := decodeCBORExact(c.AttObj, attestationLimits)
+		if err != nil || v.kind != cborMap || len(v.items) != 6 {
+			t.Fatalf("%s attestationObject: %v (%d items)", alg, err, len(v.items))
+		}
+		if k := coseKeyBytes(t, c); len(k) != keyLen {
+			t.Fatalf("%s COSE key is %d bytes, want %d", alg, len(k), keyLen)
+		} else if _, err := decodeCBORExact(k, attestationLimits); err != nil {
+			t.Fatalf("%s COSE key: %v", alg, err)
+		}
+	}
+}
+
+// CB4: never panics; an accepted item is never longer than its input.
+func FuzzDecodeCBOR(f *testing.F) {
+	f.Add([]byte{0x18})
+	f.Add(cb(intMap(16)))
+	f.Add(nested(5))
+	f.Add([]byte{0xa2, 0x01, 0x01, 0x01, 0x02})
+	f.Fuzz(func(t *testing.T, b []byte) {
+		_, rest, err := decodeCBOR(b, attestationLimits)
+		if err == nil && len(rest) > len(b) {
+			t.Fatalf("rest %d > input %d", len(rest), len(b))
+		}
+	})
+}
--- a/host/approval/cose.go
+++ b/host/approval/cose.go
@@ -0,0 +1,128 @@
+package approval
+
+// Row 201 A2 (design §2.4): COSE public keys, ES256 (P-256) and EdDSA (Ed25519) only, as an
+// authenticator emits them inside authenticator data. What is stored is Canonical(), so one key
+// has one encoding. No I/O.
+
+import (
+	"crypto/ecdsa"
+	"crypto/ed25519"
+	"crypto/elliptic"
+	"crypto/sha256"
+	"errors"
+	"fmt"
+)
+
+const (
+	AlgES256 int64 = -7
+	AlgEdDSA int64 = -8
+)
+
+var (
+	ErrCOSEKey = errors.New("approval: not an ES256 (P-256) or EdDSA (Ed25519) COSE key")
+	ErrCOSEAlg = errors.New("approval: the passkey's algorithm is not ES256 (-7) or EdDSA (-8)")
+)
+
+// COSEKey is a parsed passkey public key.
+type COSEKey struct {
+	Alg int64  // AlgES256 or AlgEdDSA
+	X   []byte // 32 bytes
+	Y   []byte // 32 bytes for ES256; nil for EdDSA
+}
+
+func coseInt(n int64) cborValue { return cborValue{kind: cborInt, i: n} }
+
+// ParseCOSEKey judges one decoded COSE key (Q1-Q8).
+func ParseCOSEKey(v cborValue) (COSEKey, error) {
+	if v.kind != cborMap { // Q1
+		return COSEKey{}, fmt.Errorf("%w: not a map", ErrCOSEKey)
+	}
+	label := func(n int64) (cborValue, bool) { return v.get(coseInt(n)) }
+	kty, ok1 := label(1)
+	alg, ok2 := label(3)
+	crv, ok3 := label(-1)
+	if !ok1 || !ok2 || !ok3 || kty.kind != cborInt || alg.kind != cborInt || crv.kind != cborInt { // Q2
+		return COSEKey{}, fmt.Errorf("%w: kty, alg and crv must be integers", ErrCOSEKey)
+	}
+	k := COSEKey{Alg: alg.i}
+	allowed := map[int64]bool{1: true, 3: true, -1: true, -2: true}
+	switch k.Alg {
+	case AlgES256:
+		if kty.i != 2 || crv.i != 1 { // Q4
+			return COSEKey{}, fmt.Errorf("%w: ES256 needs kty 2 (EC2) and crv 1 (P-256)", ErrCOSEKey)
+		}
+		allowed[-3] = true
+	case AlgEdDSA:
+		if kty.i != 1 || crv.i != 6 { // Q5
+			return COSEKey{}, fmt.Errorf("%w: EdDSA needs kty 1 (OKP) and crv 6 (Ed25519)", ErrCOSEKey)
+		}
+	default: // Q3
+		return COSEKey{}, fmt.Errorf("%w (it is %d)", ErrCOSEAlg, k.Alg)
+	}
+	for j := 0; j < len(v.items); j += 2 {
+		if key := v.items[j]; key.kind != cborInt || !allowed[key.i] { // Q8
+			return COSEKey{}, fmt.Errorf("%w: an unexpected label", ErrCOSEKey)
+		}
+	}
+	for _, c := range []struct {
+		n   int64
+		dst *[]byte
+	}{{-2, &k.X}, {-3, &k.Y}} {
+		if !allowed[c.n] {
+			continue
+		}
+		b, ok := label(c.n)
+		if !ok || b.kind != cborBytes || len(b.b) != 32 { // Q6
+			return COSEKey{}, fmt.Errorf("%w: a coordinate is not 32 bytes", ErrCOSEKey)
+		}
+		*c.dst = append([]byte(nil), b.b...)
+	}
+	if k.Alg == AlgES256 {
+		if _, err := ecdsa.ParseUncompressedPublicKey(elliptic.P256(), k.point()); err != nil { // Q7
+			return COSEKey{}, fmt.Errorf("%w: %v", ErrCOSEKey, err)
+		}
+	}
+	return k, nil
+}
+
+func (k COSEKey) point() []byte { return append(append([]byte{4}, k.X...), k.Y...) }
+
+// Canonical is the key's CTAP2 canonical encoding (labels 1, 3, -1, -2[, -3]): what is stored.
+func (k COSEKey) Canonical() []byte {
+	n, kty, crv := byte(0xa5), byte(2), byte(1)
+	if k.Alg == AlgEdDSA {
+		n, kty, crv = 0xa4, 1, 6
+	}
+	b := append([]byte{n, 0x01, kty, 0x03, byte(0x20 | (-1 - k.Alg)), 0x20, crv, 0x21, 0x58, byte(len(k.X))}, k.X...)
+	if k.Alg == AlgES256 {
+		b = append(append(b, 0x22, 0x58, byte(len(k.Y))), k.Y...)
+	}
+	return b
+}
+
+// AlgName is "ES256" or "EdDSA".
+func (k COSEKey) AlgName() string {
+	if k.Alg == AlgES256 {
+		return "ES256"
+	}
+	return "EdDSA"
+}
+
+// verify checks sig over signed: ES256 is ASN.1 DER over SHA-256(signed), EdDSA is Ed25519.
+func (k COSEKey) verify(signed, sig []byte) error {
+	switch k.Alg {
+	case AlgES256:
+		pub, err := ecdsa.ParseUncompressedPublicKey(elliptic.P256(), k.point())
+		h := sha256.Sum256(signed)
+		if len(sig) > 72 || err != nil || !ecdsa.VerifyASN1(pub, h[:], sig) {
+			return ErrBadSignature
+		}
+	case AlgEdDSA:
+		if len(sig) != ed25519.SignatureSize || len(k.X) != ed25519.PublicKeySize || !ed25519.Verify(k.X, signed, sig) {
+			return ErrBadSignature
+		}
+	default:
+		return ErrCOSEAlg
+	}
+	return nil
+}
--- a/host/approval/cose_test.go
+++ b/host/approval/cose_test.go
@@ -0,0 +1,153 @@
+package approval
+
+// Row 201 A2 sprint 1 M1 (design §6, CK1-CK4): ES256 and EdDSA COSE keys.
+
+import (
+	"bytes"
+	"crypto/sha256"
+	"errors"
+	"reflect"
+	"testing"
+)
+
+// chromeKey is the parsed COSE key of a recorded registration.
+func chromeKey(t *testing.T, alg string) (COSEKey, []byte) {
+	t.Helper()
+	raw := coseKeyBytes(t, chromeFixture(t, alg))
+	v, err := decodeCBORExact(raw, attestationLimits)
+	if err != nil {
+		t.Fatal(err)
+	}
+	k, err := ParseCOSEKey(v)
+	if err != nil {
+		t.Fatalf("%s: %v", alg, err)
+	}
+	return k, raw
+}
+
+func parseCOSEBytes(b []byte) (COSEKey, error) {
+	v, err := decodeCBORExact(b, attestationLimits)
+	if err != nil {
+		return COSEKey{}, err
+	}
+	return ParseCOSEKey(v)
+}
+
+// CK1: Q1-Q8, each refused with its error.
+func TestCOSEKeyRefusals(t *testing.T) {
+	ec, _ := chromeKey(t, "es256")
+	ed, _ := chromeKey(t, "eddsa")
+	x, y, ex := ec.X, ec.Y, ed.X
+	offCurve := bytes.Clone(y)
+	offCurve[31] ^= 1
+	arms := []struct {
+		name string
+		key  any
+		want error
+	}{
+		{"Q1-array", cbArr{1, 2}, ErrCOSEKey},
+		{"Q2-no-kty", cbMap{3, -7, -1, 1, -2, x, -3, y}, ErrCOSEKey},
+		{"Q2-no-alg", cbMap{1, 2, -1, 1, -2, x, -3, y}, ErrCOSEKey},
+		{"Q2-alg-text", cbMap{1, 2, 3, "ES256", -1, 1, -2, x, -3, y}, ErrCOSEKey},
+		{"Q2-no-crv", cbMap{1, 1, 3, -8, -2, ex}, ErrCOSEKey},
+		{"Q3-RS256", cbMap{1, 3, 3, -257, -1, 1, -2, x, -3, y}, ErrCOSEAlg},
+		{"Q3-ES384", cbMap{1, 2, 3, -35, -1, 2, -2, x, -3, y}, ErrCOSEAlg},
+		{"Q4-ES256-kty-1", cbMap{1, 1, 3, -7, -1, 1, -2, x, -3, y}, ErrCOSEKey},
+		{"Q4-ES256-crv-2", cbMap{1, 2, 3, -7, -1, 2, -2, x, -3, y}, ErrCOSEKey},
+		{"Q5-EdDSA-crv-1", cbMap{1, 1, 3, -8, -1, 1, -2, ex}, ErrCOSEKey},
+		{"Q5-EdDSA-kty-2", cbMap{1, 2, 3, -8, -1, 6, -2, ex}, ErrCOSEKey},
+		{"Q6-EdDSA-31-byte-x", cbMap{1, 1, 3, -8, -1, 6, -2, ex[:31]}, ErrCOSEKey},
+		{"Q6-EdDSA-x-text", cbMap{1, 1, 3, -8, -1, 6, -2, string(bytes.Repeat([]byte("a"), 32))}, ErrCOSEKey},
+		{"Q6-ES256-no-y", cbMap{1, 2, 3, -7, -1, 1, -2, x}, ErrCOSEKey},
+		{"Q7-off-curve", cbMap{1, 2, 3, -7, -1, 1, -2, x, -3, offCurve}, ErrCOSEKey},
+		{"Q8-kid", cbMap{1, 2, 2, []byte{1}, 3, -7, -1, 1, -2, x, -3, y}, ErrCOSEKey},
+		{"Q8-EdDSA-y", cbMap{1, 1, 3, -8, -1, 6, -2, ex, -3, ex}, ErrCOSEKey},
+		{"Q8-text-label", cbMap{1, 1, 3, -8, -1, 6, -2, ex, "x", 1}, ErrCOSEKey},
+	}
+	for _, a := range arms {
+		t.Run(a.name, func(t *testing.T) {
+			k, err := parseCOSEBytes(cb(a.key))
+			if !errors.Is(err, a.want) || (a.want == ErrCOSEKey && errors.Is(err, ErrCOSEAlg)) || !reflect.DeepEqual(k, COSEKey{}) {
+				t.Fatalf("(%+v, %v), want the zero key and %v", k, err, a.want)
+			}
+		})
+	}
+	for name, m := range map[string]cbMap{"es256": {1, 2, 3, -7, -1, 1, -2, x, -3, y}, "eddsa": {1, 1, 3, -8, -1, 6, -2, ex}} {
+		if _, err := parseCOSEBytes(cb(m)); err != nil {
+			t.Errorf("control %s: %v", name, err)
+		}
+	}
+}
+
+// CK2: Chrome's keys verify Chrome's assertion signatures over authData || SHA-256(clientDataJSON).
+func TestCOSEKeyVerifies(t *testing.T) {
+	for _, alg := range []string{"es256", "eddsa"} {
+		c := chromeFixture(t, alg)
+		k, _ := chromeKey(t, alg)
+		h := sha256.Sum256(c.AsrCDJ)
+		signed := append(bytes.Clone(c.AsrAuthData), h[:]...)
+		if err := k.verify(signed, c.AsrSig); err != nil {
+			t.Fatalf("%s: Chrome's signature: %v", alg, err)
+		}
+		tampered := bytes.Clone(signed)
+		tampered[len(tampered)-1] ^= 1
+		for name, in := range map[string][2][]byte{
+			"tampered-message": {tampered, c.AsrSig},
+			"authData-only":    {c.AsrAuthData, c.AsrSig},
+			"long-signature":   {signed, append(bytes.Clone(c.AsrSig), 0)},
+			"short-signature":  {signed, c.AsrSig[:len(c.AsrSig)-1]},
+		} {
+			if err := k.verify(in[0], in[1]); !errors.Is(err, ErrBadSignature) {
+				t.Errorf("%s %s: %v, want ErrBadSignature", alg, name, err)
+			}
+		}
+	}
+	if err := (COSEKey{Alg: -257}).verify(nil, nil); !errors.Is(err, ErrCOSEAlg) {
+		t.Errorf("an unknown algorithm: %v, want ErrCOSEAlg", err)
+	}
+}
+
+// CK3: Canonical() is Chrome's key bytes, and a reordered map re-encodes to them.
+func TestCOSECanonicalIsChromesBytes(t *testing.T) {
+	for _, alg := range []string{"es256", "eddsa"} {
+		k, raw := chromeKey(t, alg)
+		if !bytes.Equal(k.Canonical(), raw) {
+			t.Fatalf("%s: Canonical %x, Chrome %x", alg, k.Canonical(), raw)
+		}
+		m := cbMap{-2, k.X, -1, 6, 3, -8, 1, 1}
+		if alg == "es256" {
+			m = cbMap{-3, k.Y, -2, k.X, -1, 1, 3, -7, 1, 2}
+		}
+		p, err := parseCOSEBytes(cb(m))
+		if err != nil || !bytes.Equal(p.Canonical(), raw) {
+			t.Fatalf("%s: a permuted map: %v, %x", alg, err, p.Canonical())
+		}
+	}
+	if es, ed := (COSEKey{Alg: AlgES256}).AlgName(), (COSEKey{Alg: AlgEdDSA}).AlgName(); es != "ES256" || ed != "EdDSA" {
+		t.Fatalf("AlgName %q %q", es, ed)
+	}
+}
+
+// CK4: COSEKey holds exactly Alg, X and Y, and Canonical reads each.
+func TestCOSEKeyHasOnlyDeclaredFields(t *testing.T) {
+	typ := reflect.TypeOf(COSEKey{})
+	var names []string
+	for i := 0; i < typ.NumField(); i++ {
+		names = append(names, typ.Field(i).Name)
+	}
+	if !reflect.DeepEqual(names, []string{"Alg", "X", "Y"}) {
+		t.Fatalf("COSEKey fields %v, want [Alg X Y] (a new field must be encoded by Canonical)", names)
+	}
+	k, _ := chromeKey(t, "es256")
+	for name, p := range map[string]func(*COSEKey){
+		"Alg": func(c *COSEKey) { c.Alg = AlgEdDSA },
+		"X":   func(c *COSEKey) { c.X = bytes.Clone(c.X); c.X[0] ^= 1 },
+		"Y":   func(c *COSEKey) { c.Y = bytes.Clone(c.Y); c.Y[0] ^= 1 },
+	} {
+		c := k
+		p(&c)
+		if bytes.Equal(c.Canonical(), k.Canonical()) {
+			t.Errorf("changing %s does not change Canonical()", name)
+		}
+	}
+}
--- a/host/approval/testdata/webauthn/SOURCE.txt
+++ b/host/approval/testdata/webauthn/SOURCE.txt
@@ -0,0 +1,20 @@
+Row 201 A2 recorded WebAuthn fixtures (design_docs/planned/w-world-approvals-by-key-a2.md §2.15).
+
+Recorded in iteration 379 by the designer's probe (probe.mjs, sha256
+01a17c0d7bab9333dee53452159fa67a18b7acf80965aac163cfb6132f8f37cf): Google Chrome
+Chrome/157.0.8103.50 headless with a temporary profile, driven over the DevTools protocol, with a
+CDP virtual authenticator (WebAuthn.addVirtualAuthenticator: ctap2 2.1, internal transport,
+resident key, user verification, isUserVerified). A virtual authenticator is a TEST key: no real
+authenticator, passkey or credentials store was used (D-WORLD-109). Origin http://localhost:<port>,
+rpId localhost, attestation "none", user verification required, and the challenge of both the
+create and the get is the existing fixture message hash 29181d22...fef5 (request_test.go).
+
+Each chrome-<alg>.json holds the recording's public values: the Chrome version, origin, rp_id,
+challenge_hex, the registration and the first assertion exactly as PublicKeyCredential.toJSON()
+returned them, and the second assertion's authenticatorData (the probe kept only that part).
+Normalised from the probe outputs by the sprint-1 planner (TestPasskeyFixturesAreTheGeneratedBytes
+pins the bytes):
+  es256.json  sha256 c0370982dbca416a7677b4fa1ff2cd64e9e533cefc608d3dcab42a1000d6caf6
+  eddsa.json  sha256 1544721fc45dd65a553a5baab8a43fcd4df0375b0b8c2c897a199d2526f6f719
+
+Regeneration lands in sprint 3 (M7, the committed DevTools driver); see generate.sh.
--- a/host/approval/testdata/webauthn/chrome-eddsa.json
+++ b/host/approval/testdata/webauthn/chrome-eddsa.json
@@ -0,0 +1,10 @@
+{
+  "source": "row 201 A2, iteration 379: Chrome headless, CDP virtual authenticator (a test key, D-WORLD-109); see SOURCE.txt",
+  "chrome": "Chrome/157.0.8103.50",
+  "origin": "http://localhost:56629",
+  "rp_id": "localhost",
+  "challenge_hex": "29181d224d982f659d1bd4703238d3f8489e26968ffe2a22bd786fa3bf15fef5",
+  "registration": {"authenticatorAttachment":"platform","clientExtensionResults":{},"id":"15HgLiYRkyGZ2RT27yosO_nqoohHe4OogLJlNUK-bbY","rawId":"15HgLiYRkyGZ2RT27yosO_nqoohHe4OogLJlNUK-bbY","response":{"attestationObject":"o2NmbXRkbm9uZWdhdHRTdG10oGhhdXRoRGF0YViBSZYN5YgOjGh0NBcPZHZgW4_krrmihjLHmVzzuoMdl2NFAAAAAQECAwQFBgcIAQIDBAUGBwgAINeR4C4mEZMhmdkU9u8qLDv56qKIR3uDqICyZTVCvm22pAEBAycgBiFYICbpoZHRqpAlkXWdFG6utclYPxkX5JYG3RT0ECxxD2ro","authenticatorData":"SZYN5YgOjGh0NBcPZHZgW4_krrmihjLHmVzzuoMdl2NFAAAAAQECAwQFBgcIAQIDBAUGBwgAINeR4C4mEZMhmdkU9u8qLDv56qKIR3uDqICyZTVCvm22pAEBAycgBiFYICbpoZHRqpAlkXWdFG6utclYPxkX5JYG3RT0ECxxD2ro","clientDataJSON":"eyJ0eXBlIjoid2ViYXV0aG4uY3JlYXRlIiwiY2hhbGxlbmdlIjoiS1JnZElrMllMMldkRzlSd01qalQtRWllSnBhUF9pb2l2WGh2bzc4Vl92VSIsIm9yaWdpbiI6Imh0dHA6Ly9sb2NhbGhvc3Q6NTY2MjkiLCJjcm9zc09yaWdpbiI6ZmFsc2V9","publicKey":"MCowBQYDK2VwAyEAJumhkdGqkCWRdZ0Ubq61yVg_GRfklgbdFPQQLHEPaug","publicKeyAlgorithm":-8,"transports":["internal"]},"type":"public-key"},
+  "assertion": {"authenticatorAttachment":"platform","clientExtensionResults":{},"id":"15HgLiYRkyGZ2RT27yosO_nqoohHe4OogLJlNUK-bbY","rawId":"15HgLiYRkyGZ2RT27yosO_nqoohHe4OogLJlNUK-bbY","response":{"authenticatorData":"SZYN5YgOjGh0NBcPZHZgW4_krrmihjLHmVzzuoMdl2MFAAAAAg","clientDataJSON":"eyJ0eXBlIjoid2ViYXV0aG4uZ2V0IiwiY2hhbGxlbmdlIjoiS1JnZElrMllMMldkRzlSd01qalQtRWllSnBhUF9pb2l2WGh2bzc4Vl92VSIsIm9yaWdpbiI6Imh0dHA6Ly9sb2NhbGhvc3Q6NTY2MjkiLCJjcm9zc09yaWdpbiI6ZmFsc2V9","signature":"ltfy2YVJXOPy9P3aKQ9qroe4jxGQg-6Q0lbl98_F5ARiAiNckkJpgZ6rLSzNtM7fbTasGLYOcO1DdMjhdbZ4Aw"},"type":"public-key"},
+  "assertion2": {"authenticatorData":"SZYN5YgOjGh0NBcPZHZgW4_krrmihjLHmVzzuoMdl2MFAAAAAw"}
+}
--- a/host/approval/testdata/webauthn/chrome-es256.json
+++ b/host/approval/testdata/webauthn/chrome-es256.json
@@ -0,0 +1,10 @@
+{
+  "source": "row 201 A2, iteration 379: Chrome headless, CDP virtual authenticator (a test key, D-WORLD-109); see SOURCE.txt",
+  "chrome": "Chrome/157.0.8103.50",
+  "origin": "http://localhost:56547",
+  "rp_id": "localhost",
+  "challenge_hex": "29181d224d982f659d1bd4703238d3f8489e26968ffe2a22bd786fa3bf15fef5",
+  "registration": {"authenticatorAttachment":"platform","clientExtensionResults":{},"id":"18ZYwHpS0jBgdaH6F89nTjl0LuAaxzs-k2kc0eJYNDI","rawId":"18ZYwHpS0jBgdaH6F89nTjl0LuAaxzs-k2kc0eJYNDI","response":{"attestationObject":"o2NmbXRkbm9uZWdhdHRTdG10oGhhdXRoRGF0YVikSZYN5YgOjGh0NBcPZHZgW4_krrmihjLHmVzzuoMdl2NFAAAAAQECAwQFBgcIAQIDBAUGBwgAINfGWMB6UtIwYHWh-hfPZ045dC7gGsc7PpNpHNHiWDQypQECAyYgASFYIM1s_xRHouyhiYBh_0krF24rmScb5nMISw64W9OWNHjtIlggrJhbdhYWyyY2wdvW-VJduL4pGlGHLeOMCGZvq3aZ_hA","authenticatorData":"SZYN5YgOjGh0NBcPZHZgW4_krrmihjLHmVzzuoMdl2NFAAAAAQECAwQFBgcIAQIDBAUGBwgAINfGWMB6UtIwYHWh-hfPZ045dC7gGsc7PpNpHNHiWDQypQECAyYgASFYIM1s_xRHouyhiYBh_0krF24rmScb5nMISw64W9OWNHjtIlggrJhbdhYWyyY2wdvW-VJduL4pGlGHLeOMCGZvq3aZ_hA","clientDataJSON":"eyJ0eXBlIjoid2ViYXV0aG4uY3JlYXRlIiwiY2hhbGxlbmdlIjoiS1JnZElrMllMMldkRzlSd01qalQtRWllSnBhUF9pb2l2WGh2bzc4Vl92VSIsIm9yaWdpbiI6Imh0dHA6Ly9sb2NhbGhvc3Q6NTY1NDciLCJjcm9zc09yaWdpbiI6ZmFsc2V9","publicKey":"MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEzWz_FEei7KGJgGH_SSsXbiuZJxvmcwhLDrhb05Y0eO2smFt2FhbLJjbB29b5Ul24vikaUYct44wIZm-rdpn-EA","publicKeyAlgorithm":-7,"transports":["internal"]},"type":"public-key"},
+  "assertion": {"authenticatorAttachment":"platform","clientExtensionResults":{},"id":"18ZYwHpS0jBgdaH6F89nTjl0LuAaxzs-k2kc0eJYNDI","rawId":"18ZYwHpS0jBgdaH6F89nTjl0LuAaxzs-k2kc0eJYNDI","response":{"authenticatorData":"SZYN5YgOjGh0NBcPZHZgW4_krrmihjLHmVzzuoMdl2MFAAAAAg","clientDataJSON":"eyJ0eXBlIjoid2ViYXV0aG4uZ2V0IiwiY2hhbGxlbmdlIjoiS1JnZElrMllMMldkRzlSd01qalQtRWllSnBhUF9pb2l2WGh2bzc4Vl92VSIsIm9yaWdpbiI6Imh0dHA6Ly9sb2NhbGhvc3Q6NTY1NDciLCJjcm9zc09yaWdpbiI6ZmFsc2V9","signature":"MEYCIQCd2FJZsQPiF1HAHoE5TQu6hukSvK0XtddumKSgfXh9NwIhALyLnjuwVR1iBPOGSm0I6h2l_e4nYIBr4vShU-2imih6"},"type":"public-key"},
+  "assertion2": {"authenticatorData":"SZYN5YgOjGh0NBcPZHZgW4_krrmihjLHmVzzuoMdl2MFAAAAAw"}
+}
--- a/host/approval/testdata/webauthn/generate.sh
+++ b/host/approval/testdata/webauthn/generate.sh
@@ -0,0 +1,7 @@
+#!/bin/sh
+# Run as: sh generate.sh. Re-records chrome-es256.json and chrome-eddsa.json with Chrome's
+# virtual authenticator (a test key, D-WORLD-109). Row 201 A2 sprint 1 commits the iteration-379
+# recording (SOURCE.txt); the driver this script runs, host/daemon/testdata/passkey_cdp.mjs,
+# lands in sprint 3 (M7).
+echo "generate.sh: regeneration lands in row 201 A2 sprint 3 (M7); see SOURCE.txt" >&2
+exit 1
--- a/host/approval/webauthn_fixture_test.go
+++ b/host/approval/webauthn_fixture_test.go
@@ -0,0 +1,165 @@
+package approval
+
+// Row 201 A2: the recorded Chrome fixtures (testdata/webauthn, design §2.15) and the test CBOR
+// encoder the refusal arms are built with.
+
+import (
+	"crypto/sha256"
+	"encoding/base64"
+	"encoding/binary"
+	"encoding/hex"
+	"encoding/json"
+	"os"
+	"path/filepath"
+	"strings"
+	"testing"
+)
+
+// chromeCeremony is one recording, with its base64url fields decoded.
+type chromeCeremony struct {
+	Origin, RPID, ChallengeHex  string
+	Registration, Assertion     []byte // the toJSON() envelopes, as recorded
+	RegRawID, RegCDJ, AttObj    []byte
+	AsrCDJ, AsrAuthData, AsrSig []byte
+	Asr2AuthData                []byte
+}
+
+func b64u(t *testing.T, s string) []byte {
+	t.Helper()
+	b, err := base64.RawURLEncoding.Strict().DecodeString(s)
+	if err != nil {
+		t.Fatal(err)
+	}
+	return b
+}
+
+// chromeFixture loads testdata/webauthn/chrome-<alg>.json ("es256" or "eddsa").
+func chromeFixture(t *testing.T, alg string) chromeCeremony {
+	t.Helper()
+	raw, err := os.ReadFile(filepath.Join("testdata", "webauthn", "chrome-"+alg+".json"))
+	if err != nil {
+		t.Fatal(err)
+	}
+	type resp struct{ ClientDataJSON, AttestationObject, AuthenticatorData, Signature string }
+	var f struct {
+		Origin                  string
+		RPID                    string `json:"rp_id"`
+		ChallengeHex            string `json:"challenge_hex"`
+		Registration, Assertion json.RawMessage
+		Assertion2              struct{ AuthenticatorData string }
+	}
+	var reg, asr struct {
+		RawID    string
+		Response resp
+	}
+	if json.Unmarshal(raw, &f) != nil || json.Unmarshal(f.Registration, &reg) != nil || json.Unmarshal(f.Assertion, &asr) != nil {
+		t.Fatalf("chrome-%s.json does not decode", alg)
+	}
+	return chromeCeremony{Origin: f.Origin, RPID: f.RPID, ChallengeHex: f.ChallengeHex,
+		Registration: f.Registration, Assertion: f.Assertion, RegRawID: b64u(t, reg.RawID),
+		RegCDJ: b64u(t, reg.Response.ClientDataJSON), AttObj: b64u(t, reg.Response.AttestationObject),
+		AsrCDJ: b64u(t, asr.Response.ClientDataJSON), AsrAuthData: b64u(t, asr.Response.AuthenticatorData),
+		AsrSig: b64u(t, asr.Response.Signature), Asr2AuthData: b64u(t, f.Assertion2.AuthenticatorData)}
+}
+
+// regAuthData is the authData inside a recorded attestationObject; coseKeyBytes is the COSE key
+// inside it (after rpIdHash, flags, counter, aaguid, the length and the credential ID).
+func regAuthData(t *testing.T, c chromeCeremony) []byte {
+	t.Helper()
+	v, err := decodeCBORExact(c.AttObj, attestationLimits)
+	if err != nil {
+		t.Fatal(err)
+	}
+	ad, ok := v.get(cborValue{kind: cborText, b: []byte("authData")})
+	if !ok {
+		t.Fatal("no authData")
+	}
+	return ad.b
+}
+
+func coseKeyBytes(t *testing.T, c chromeCeremony) []byte {
+	ad := regAuthData(t, c)
+	return ad[55+int(binary.BigEndian.Uint16(ad[53:55])):]
+}
+
+// XW0: the fixtures are the recorded bytes, and SOURCE.txt names their origin.
+func TestPasskeyFixturesAreTheGeneratedBytes(t *testing.T) {
+	for name, want := range map[string]string{
+		"chrome-es256.json": "924e9f79e873d487d83e93f0df22da342e1130db85c3e5dddc6bec9129e27639",
+		"chrome-eddsa.json": "4c3fe46cb62898bd33c83244728ea84ace5ca25cb7df550e3759842928467b40",
+	} {
+		b, err := os.ReadFile(filepath.Join("testdata", "webauthn", name))
+		if err != nil {
+			t.Fatal(err)
+		}
+		if got := sha256.Sum256(b); hex.EncodeToString(got[:]) != want {
+			t.Errorf("%s sha256 %x, want %s (fixtures are recorded, never hand-edited)", name, got, want)
+		}
+	}
+	src, err := os.ReadFile(filepath.Join("testdata", "webauthn", "SOURCE.txt"))
+	if err != nil {
+		t.Fatal(err)
+	}
+	for _, want := range []string{"Chrome/157.0.8103.50", "virtual authenticator", "D-WORLD-109"} {
+		if !strings.Contains(string(src), want) {
+			t.Errorf("SOURCE.txt does not name %q", want)
+		}
+	}
+	for _, alg := range []string{"es256", "eddsa"} {
+		c := chromeFixture(t, alg)
+		if c.RPID != "localhost" || c.ChallengeHex != fixtureMessageHashHex || !strings.HasPrefix(c.Origin, "http://localhost:") {
+			t.Errorf("%s: rp %q challenge %q origin %q", alg, c.RPID, c.ChallengeHex, c.Origin)
+		}
+	}
+}
+
+// ---- the test encoder: shortest heads only, so a non-canonical arm is written by hand.
+
+type cbArr []any
+type cbMap []any // key, value, key, value, ...
+
+func cbHead(major byte, n uint64) []byte {
+	switch {
+	case n < 24:
+		return []byte{major<<5 | byte(n)}
+	case n < 1<<8:
+		return []byte{major<<5 | 24, byte(n)}
+	case n < 1<<16:
+		return binary.BigEndian.AppendUint16([]byte{major<<5 | 25}, uint16(n))
+	case n < 1<<32:
+		return binary.BigEndian.AppendUint32([]byte{major<<5 | 26}, uint32(n))
+	}
+	return binary.BigEndian.AppendUint64([]byte{major<<5 | 27}, n)
+}
+
+func cb(v any) []byte {
+	switch x := v.(type) {
+	case int:
+		if x < 0 {
+			return cbHead(1, uint64(-1-x))
+		}
+		return cbHead(0, uint64(x))
+	case []byte:
+		return append(cbHead(2, uint64(len(x))), x...)
+	case string:
+		return append(cbHead(3, uint64(len(x))), x...)
+	case bool:
+		if x {
+			return []byte{0xf5}
+		}
+		return []byte{0xf4}
+	case cbArr:
+		b := cbHead(4, uint64(len(x)))
+		for _, e := range x {
+			b = append(b, cb(e)...)
+		}
+		return b
+	case cbMap:
+		b := cbHead(5, uint64(len(x)/2))
+		for _, e := range x {
+			b = append(b, cb(e)...)
+		}
+		return b
+	}
+	panic("cb: unsupported type")
+}
````

### Patch M2 — authenticator data, clientDataJSON, the origin and the envelope

````diff M2
--- a/host/approval/authdata.go
+++ b/host/approval/authdata.go
@@ -0,0 +1,85 @@
+package approval
+
+// Row 201 A2 (design §2.5): authenticator data, laid out as WebAuthn L3 §6.1 says: rpIdHash (32),
+// flags (1), signCount (4, big-endian), then for a registration the attested credential data,
+// then any extensions map. No I/O.
+
+import (
+	"encoding/binary"
+	"errors"
+	"fmt"
+)
+
+// MaxCredentialIDBytes is L3 §7.1 step 25's bound; there is no lower bound but 1 (design §2.5 [r1]).
+const MaxCredentialIDBytes = 1023
+
+var (
+	ErrAuthDataShort      = errors.New("approval: authenticator data is truncated")
+	ErrAuthDataFlags      = errors.New("approval: authenticator data has inconsistent flags")
+	ErrCredentialID       = errors.New("approval: a credential ID must be 1 to 1023 bytes")
+	ErrAuthDataExtensions = errors.New("approval: authenticator data extensions are not one CBOR map")
+	ErrAuthDataTrailing   = errors.New("approval: authenticator data has trailing bytes")
+)
+
+// AuthData is parsed authenticator data. The last three fields are a registration's only.
+type AuthData struct {
+	RPIDHash       [32]byte
+	UP, UV, BE, BS bool // flags bits 0, 2, 3, 4
+	SignCount      uint32
+	CredentialID   []byte
+	AAGUID         [16]byte
+	Key            COSEKey
+}
+
+// ParseAuthData judges authenticator data (D1-D8); registration says whether attested
+// credential data must be present (a registration) or absent (an assertion).
+func ParseAuthData(b []byte, registration bool) (AuthData, error) {
+	if len(b) < 37 { // D1
+		return AuthData{}, ErrAuthDataShort
+	}
+	var a AuthData
+	copy(a.RPIDHash[:], b[:32])
+	f := b[32]
+	a.UP, a.UV, a.BE, a.BS = f&0x01 != 0, f&0x04 != 0, f&0x08 != 0, f&0x10 != 0
+	if a.BS && !a.BE { // D2
+		return AuthData{}, fmt.Errorf("%w: backup state without backup eligibility", ErrAuthDataFlags)
+	}
+	if (f&0x40 != 0) != registration { // D3
+		return AuthData{}, fmt.Errorf("%w: attested credential data is for a registration only", ErrAuthDataFlags)
+	}
+	a.SignCount = binary.BigEndian.Uint32(b[33:37])
+	rest := b[37:]
+	if registration {
+		if len(rest) < 18 { // D4
+			return AuthData{}, ErrAuthDataShort
+		}
+		copy(a.AAGUID[:], rest[:16])
+		n := int(binary.BigEndian.Uint16(rest[16:18]))
+		if n < 1 || n > MaxCredentialIDBytes { // D5
+			return AuthData{}, ErrCredentialID
+		}
+		if len(rest) < 18+n { // D4
+			return AuthData{}, ErrAuthDataShort
+		}
+		a.CredentialID = append([]byte(nil), rest[18:18+n]...)
+		v, r, err := decodeCBOR(rest[18+n:], attestationLimits)
+		if err == nil {
+			a.Key, err = ParseCOSEKey(v)
+		}
+		if err != nil { // D6
+			return AuthData{}, fmt.Errorf("approval: the passkey's public key: %w", err)
+		}
+		rest = r
+	}
+	if f&0x80 != 0 {
+		v, r, err := decodeCBOR(rest, attestationLimits)
+		if err != nil || v.kind != cborMap { // D7
+			return AuthData{}, ErrAuthDataExtensions
+		}
+		rest = r
+	}
+	if len(rest) != 0 { // D8
+		return AuthData{}, ErrAuthDataTrailing
+	}
+	return a, nil
+}
--- a/host/approval/clientdata.go
+++ b/host/approval/clientdata.go
@@ -0,0 +1,115 @@
+package approval
+
+// Row 201 A2 (design §2.6): clientDataJSON, read token by token so a duplicated member is
+// refused (Unmarshal would keep the last one silently). Unknown members are allowed, bounded.
+
+import (
+	"bytes"
+	"encoding/json"
+	"errors"
+	"io"
+	"unicode/utf8"
+)
+
+// MaxClientDataBytes bounds clientDataJSON (J1).
+const MaxClientDataBytes = 2048
+
+var (
+	ErrClientData          = errors.New("approval: clientDataJSON is not a valid client data object")
+	ErrClientDataDuplicate = errors.New("approval: clientDataJSON repeats a member")
+)
+
+// ClientData is the part of CollectedClientData the laws read (L3 §5.8.1).
+type ClientData struct {
+	Type, Challenge, Origin string
+	CrossOrigin             bool
+	HasTopOrigin            bool
+}
+
+// ParseClientData judges clientDataJSON (J1-J7).
+func ParseClientData(b []byte) (ClientData, error) {
+	if len(b) == 0 || len(b) > MaxClientDataBytes { // J1
+		return ClientData{}, ErrClientData
+	}
+	if !utf8.Valid(b) { // J2
+		return ClientData{}, ErrClientData
+	}
+	dec := json.NewDecoder(bytes.NewReader(b))
+	if tok, err := dec.Token(); err != nil || tok != json.Delim('{') { // J3
+		return ClientData{}, ErrClientData
+	}
+	var c ClientData
+	seen := map[string]bool{}
+	for dec.More() {
+		tok, _ := dec.Token()
+		key, ok := tok.(string)
+		if !ok { // J3: not a member name (a failed read is nil)
+			return ClientData{}, ErrClientData
+		}
+		strs := map[string]*string{"type": &c.Type, "challenge": &c.Challenge, "origin": &c.Origin}
+		if _, known := strs[key]; !known && key != "crossOrigin" && key != "topOrigin" {
+			if err := skipJSON(dec, 2); err != nil {
+				return ClientData{}, err
+			}
+			continue
+		}
+		if seen[key] { // J4
+			return ClientData{}, ErrClientDataDuplicate
+		}
+		seen[key] = true
+		val, _ := dec.Token() // a failed read is nil: neither a string nor a bool
+		s, isStr := val.(string)
+		bv, isBool := val.(bool)
+		switch {
+		case strs[key] != nil:
+			if !isStr { // J5
+				return ClientData{}, ErrClientData
+			}
+			*strs[key] = s
+		case key == "crossOrigin":
+			if !isBool { // J6
+				return ClientData{}, ErrClientData
+			}
+			c.CrossOrigin = bv
+		default: // topOrigin
+			if !isStr { // J6
+				return ClientData{}, ErrClientData
+			}
+			c.HasTopOrigin = true
+		}
+	}
+	if _, err := dec.Token(); err != nil { // J3: the closing brace
+		return ClientData{}, ErrClientData
+	}
+	if _, err := dec.Token(); err != io.EOF { // J3: one object, nothing after it
+		return ClientData{}, ErrClientData
+	}
+	if !seen["type"] || !seen["challenge"] || !seen["origin"] { // J5
+		return ClientData{}, ErrClientData
+	}
+	return c, nil
+}
+
+// skipJSON reads one value of an unknown member (a member name is read as a value: the token
+// stream is the same), refusing nesting over depth 4 (J7) and any token that does not read (J3).
+func skipJSON(dec *json.Decoder, depth int) error {
+	tok, err := dec.Token()
+	if err != nil { // J3
+		return ErrClientData
+	}
+	if tok != json.Delim('{') && tok != json.Delim('[') {
+		return nil
+	}
+	if depth > 4 { // J7
+		return ErrClientData
+	}
+	for dec.More() {
+		if err := skipJSON(dec, depth+1); err != nil {
+			return err
+		}
+	}
+	if _, err := dec.Token(); err != nil { // J3: the closing delimiter
+		return ErrClientData
+	}
+	return nil
+}
--- a/host/approval/origin.go
+++ b/host/approval/origin.go
@@ -0,0 +1,45 @@
+package approval
+
+// Row 201 A2 (design §2.2): the one origin a passkey ceremony may run at, and its rpID. Only
+// localhost and a tailnet *.ts.net host, never an IP literal; rpID is the host exactly.
+
+import (
+	"errors"
+	"fmt"
+	"regexp"
+	"strconv"
+	"strings"
+)
+
+var (
+	ErrOrigin     = errors.New("approval: a ceremony origin is scheme://host[:port], with https for *.ts.net and http or https for localhost")
+	ErrOriginHost = errors.New("approval: a ceremony origin's host must be localhost or a tailnet *.ts.net name (an IP address cannot be a WebAuthn origin: use http://localhost:<port>)")
+)
+
+var originHost = regexp.MustCompile(`^[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?(\.[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?)*$`)
+
+// ParseCeremonyOrigin accepts exactly https://<host>.ts.net[:port] and http(s)://localhost[:port]
+// (O1-O3) and returns the rpID, which is the host.
+func ParseCeremonyOrigin(s string) (rpID string, err error) {
+	scheme, rest, ok := strings.Cut(s, "://")
+	host, port, hasPort := strings.Cut(rest, ":")
+	if !ok || (scheme != "https" && scheme != "http") || len(host) > 253 || !originHost.MatchString(host) { // O1
+		return "", ErrOrigin
+	}
+	if host != "localhost" && !strings.HasSuffix(host, ".ts.net") { // O2
+		return "", ErrOriginHost
+	}
+	if scheme == "http" && host != "localhost" { // O1: http is for localhost only
+		return "", fmt.Errorf("%w (http is accepted for localhost only)", ErrOrigin)
+	}
+	if hasPort {
+		n, perr := strconv.Atoi(port)
+		if perr != nil || n < 1 || n > 65535 || port != strconv.Itoa(n) { // O3: a bad port
+			return "", fmt.Errorf("%w (the port must be 1 to 65535 with no leading zero)", ErrOrigin)
+		}
+		if (scheme == "https" && n == 443) || (scheme == "http" && n == 80) { // O3: a default port
+			return "", fmt.Errorf("%w (leave out the default port :%d: the browser's origin omits it, so it could never match)", ErrOrigin, n)
+		}
+	}
+	return host, nil
+}
--- a/host/approval/passkey_json.go
+++ b/host/approval/passkey_json.go
@@ -0,0 +1,98 @@
+package approval
+
+// Row 201 A2 (design §2.7): the response files the ceremony page saves: the subset of L3's
+// PublicKeyCredential.toJSON() the laws read. Unknown members are ignored; no envelope value is
+// trusted unchecked (rawId is checked against authData or selects the key; the byte fields are
+// covered by the enrolment code or the signature).
+
+import (
+	"encoding/base64"
+	"encoding/json"
+	"errors"
+)
+
+// MaxPasskeyResponseBytes bounds a response file (E1, checked by the CLI's reader).
+const MaxPasskeyResponseBytes = 16 << 10
+
+var ErrPasskeyResponse = errors.New("approval: not a passkey response file (id, rawId, type public-key, response)")
+
+// RegistrationResponse is a parsed registration file.
+type RegistrationResponse struct{ RawID, ClientDataJSON, AttestationObject []byte }
+
+// AssertionResponse is a parsed assertion file. UserHandle is nil when absent or null.
+type AssertionResponse struct {
+	RawID, ClientDataJSON, AuthenticatorData, Signature, UserHandle []byte
+	HasUserHandle                                                   bool
+}
+
+type passkeyEnvelope struct {
+	ID, RawID, Type *string
+	Response        *struct {
+		ClientDataJSON, AttestationObject, AuthenticatorData, Signature, UserHandle *string
+	}
+}
+
+// decodeEnvelope parses the envelope (E2, E3) and returns rawId's bytes.
+func decodeEnvelope(b []byte) (passkeyEnvelope, []byte, error) {
+	var e passkeyEnvelope
+	if err := json.Unmarshal(b, &e); err != nil || e.ID == nil || e.RawID == nil || e.Type == nil || e.Response == nil { // E2
+		return passkeyEnvelope{}, nil, ErrPasskeyResponse
+	}
+	if *e.Type != "public-key" || *e.ID != *e.RawID { // E3
+		return passkeyEnvelope{}, nil, ErrPasskeyResponse
+	}
+	raw, err := b64Field(e.RawID)
+	return e, raw, err
+}
+
+// b64Field decodes a required member as strict, unpadded base64url (E2 when missing, E4).
+func b64Field(s *string) ([]byte, error) {
+	if s == nil { // E2
+		return nil, ErrPasskeyResponse
+	}
+	b, err := base64.RawURLEncoding.Strict().DecodeString(*s)
+	if err != nil { // E4
+		return nil, ErrPasskeyResponse
+	}
+	return b, nil
+}
+
+// ParseRegistrationResponse reads a registration file.
+func ParseRegistrationResponse(b []byte) (RegistrationResponse, error) {
+	e, raw, err := decodeEnvelope(b)
+	if err != nil {
+		return RegistrationResponse{}, err
+	}
+	r := RegistrationResponse{RawID: raw}
+	if r.ClientDataJSON, err = b64Field(e.Response.ClientDataJSON); err != nil {
+		return RegistrationResponse{}, err
+	}
+	if r.AttestationObject, err = b64Field(e.Response.AttestationObject); err != nil {
+		return RegistrationResponse{}, err
+	}
+	return r, nil
+}
+
+// ParseAssertionResponse reads an assertion file.
+func ParseAssertionResponse(b []byte) (AssertionResponse, error) {
+	e, raw, err := decodeEnvelope(b)
+	if err != nil {
+		return AssertionResponse{}, err
+	}
+	r := AssertionResponse{RawID: raw}
+	for _, f := range []struct {
+		src *string
+		dst *[]byte
+	}{{e.Response.ClientDataJSON, &r.ClientDataJSON}, {e.Response.AuthenticatorData, &r.AuthenticatorData}, {e.Response.Signature, &r.Signature}} {
+		if *f.dst, err = b64Field(f.src); err != nil {
+			return AssertionResponse{}, err
+		}
+	}
+	if e.Response.UserHandle != nil {
+		if r.UserHandle, err = b64Field(e.Response.UserHandle); err != nil {
+			return AssertionResponse{}, err
+		}
+		r.HasUserHandle = true
+	}
+	return r, nil
+}
--- a/host/approval/webauthn_parts_test.go
+++ b/host/approval/webauthn_parts_test.go
@@ -0,0 +1,262 @@
+package approval
+
+// Row 201 A2 sprint 1 M2 (design §6, AD1-AD3, CD1-CD2, OR1, EN1): authenticator data,
+// clientDataJSON, the ceremony origin and the response envelope.
+
+import (
+	"bytes"
+	"crypto/sha256"
+	"encoding/binary"
+	"errors"
+	"reflect"
+	"strings"
+	"testing"
+)
+
+// withCredentialID rebuilds a registration's authData around an n-byte credential ID.
+func withCredentialID(reg []byte, n int) []byte {
+	l := int(binary.BigEndian.Uint16(reg[53:55]))
+	b := binary.BigEndian.AppendUint16(bytes.Clone(reg[:53]), uint16(n))
+	return append(append(b, bytes.Repeat([]byte{7}, n)...), reg[55+l:]...)
+}
+
+func withFlags(b []byte, set, clear byte) []byte {
+	b = bytes.Clone(b)
+	b[32] = b[32]&^clear | set
+	return b
+}
+
+// AD1: D1-D8, with the credential-ID bounds at 0, 1, 1023 and 1024 bytes.
+func TestAuthDataRefusals(t *testing.T) {
+	c := chromeFixture(t, "es256")
+	reg, asr := regAuthData(t, c), c.AsrAuthData
+	key := coseKeyBytes(t, c)
+	badKey := append(bytes.Clone(reg[:len(reg)-len(key)]), cb(cbMap{1, 3, 3, -257, -1, 1})...)
+	arms := []struct {
+		name         string
+		in           []byte
+		registration bool
+		want         error
+	}{
+		{"D1-36-bytes", bytes.Clone(asr[:36]), false, ErrAuthDataShort},
+		{"D2-BS-without-BE", withFlags(asr, 0x10, 0), false, ErrAuthDataFlags},
+		{"D3-registration-without-AT", withFlags(reg, 0, 0x40), true, ErrAuthDataFlags},
+		{"D3-assertion-with-AT", withFlags(asr, 0x40, 0), false, ErrAuthDataFlags},
+		{"D4-inside-aaguid", bytes.Clone(reg[:37+17]), true, ErrAuthDataShort},
+		{"D4-inside-credential-id", bytes.Clone(reg[:55+10]), true, ErrAuthDataShort},
+		{"D5-0-byte-id", withCredentialID(reg, 0), true, ErrCredentialID},
+		{"D5-1024-byte-id", withCredentialID(reg, 1024), true, ErrCredentialID},
+		{"D6-RS256-key", badKey, true, ErrCOSEAlg},
+		{"D6-truncated-key", bytes.Clone(reg[:len(reg)-1]), true, ErrCBORTruncated},
+		{"D7-integer-extensions", append(withFlags(reg, 0x80, 0), 0x01), true, ErrAuthDataExtensions},
+		{"D7-no-extensions", withFlags(asr, 0x80, 0), false, ErrAuthDataExtensions},
+		{"D8-assertion-trailing", append(bytes.Clone(asr), 0), false, ErrAuthDataTrailing},
+		{"D8-registration-trailing", append(bytes.Clone(reg), 0), true, ErrAuthDataTrailing},
+	}
+	for _, a := range arms {
+		t.Run(a.name, func(t *testing.T) {
+			if ad, err := ParseAuthData(a.in, a.registration); !errors.Is(err, a.want) || !reflect.DeepEqual(ad, AuthData{}) {
+				t.Fatalf("(%+v, %v), want the zero value and %v", ad, err, a.want)
+			}
+		})
+	}
+	for _, n := range []int{1, 1023} {
+		ad, err := ParseAuthData(withCredentialID(reg, n), true)
+		if err != nil || len(ad.CredentialID) != n {
+			t.Errorf("a %d-byte credential ID: %v (%d bytes)", n, err, len(ad.CredentialID))
+		}
+	}
+}
+
+// AD2: Chrome's registration and assertions read as recorded (V21).
+func TestAuthDataReadsChrome(t *testing.T) {
+	for _, alg := range []string{"es256", "eddsa"} {
+		c := chromeFixture(t, alg)
+		ad, err := ParseAuthData(regAuthData(t, c), true)
+		if err != nil {
+			t.Fatal(err)
+		}
+		k, _ := chromeKey(t, alg)
+		aaguid := [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8}
+		if ad.RPIDHash != sha256.Sum256([]byte("localhost")) || !ad.UP || !ad.UV || ad.BE || ad.BS || ad.SignCount != 1 ||
+			ad.AAGUID != aaguid || len(ad.CredentialID) != 32 || !bytes.Equal(ad.CredentialID, c.RegRawID) || !reflect.DeepEqual(ad.Key, k) {
+			t.Fatalf("%s registration: %+v", alg, ad)
+		}
+		for want, b := range map[uint32][]byte{2: c.AsrAuthData, 3: c.Asr2AuthData} {
+			as, err := ParseAuthData(b, false)
+			if err != nil || len(b) != 37 || !as.UP || !as.UV || as.BE || as.BS || as.SignCount != want || as.CredentialID != nil {
+				t.Fatalf("%s assertion %d: %+v, %v", alg, want, as, err)
+			}
+		}
+	}
+}
+
+// AD3: an extensions map is parsed past and its contents ignored.
+func TestAuthDataIgnoresExtensionsContents(t *testing.T) {
+	c := chromeFixture(t, "eddsa")
+	ext := cb(cbMap{"credProtect", 2})
+	ad, err := ParseAuthData(append(withFlags(regAuthData(t, c), 0x80, 0), ext...), true)
+	if k, _ := chromeKey(t, "eddsa"); err != nil || !reflect.DeepEqual(ad.Key, k) {
+		t.Fatalf("registration with extensions: %v, %+v", err, ad.Key)
+	}
+	if _, err := ParseAuthData(append(withFlags(c.AsrAuthData, 0x80, 0), ext...), false); err != nil {
+		t.Fatalf("assertion with extensions: %v", err)
+	}
+}
+
+const cdjBase = `"type":"webauthn.get","challenge":"abc","origin":"http://localhost:7644"`
+
+// CD1: J1-J7, and what is allowed: unknown members, any member order, topOrigin and crossOrigin.
+func TestClientDataRefusals(t *testing.T) {
+	arms := []struct {
+		name, in string
+		want     error
+	}{
+		{"J1-empty", "", ErrClientData},
+		{"J1-2049-bytes", "{" + cdjBase + `,"x":"` + strings.Repeat("a", 2049-len(cdjBase)-9) + `"}`, ErrClientData},
+		{"J2-not-utf8", `{"type":"webauthn.get","challenge":"abc","origin":"http://localhost:7644","x":"` + "\xff" + `"}`, ErrClientData},
+		{"J3-array", `[]`, ErrClientData},
+		{"J3-array-of-pairs", `["type","webauthn.get","challenge","abc","origin","http://localhost:7644"]`, ErrClientData},
+		{"J3-unterminated", "{" + cdjBase, ErrClientData},
+		{"J3-two-objects", "{" + cdjBase + "}{}", ErrClientData},
+		{"J3-trailing", "{" + cdjBase + "}x", ErrClientData},
+		{"J3-unknown-member-unterminated", "{" + cdjBase + `,"x":[1`, ErrClientData},
+		{"J3-number-key", `{1:2}`, ErrClientData},
+		{"J3-unknown-member-no-value", "{" + cdjBase + `,"x":}`, ErrClientData},
+		{"J3-unknown-member-number-key", "{" + cdjBase + `,"x":{1:2}}`, ErrClientData},
+		{"J3-known-member-no-value", `{"type":}`, ErrClientData},
+		{"J4-challenge-twice", "{" + cdjBase + `,"challenge":"xyz"}`, ErrClientDataDuplicate},
+		{"J4-crossOrigin-twice", "{" + cdjBase + `,"crossOrigin":false,"crossOrigin":true}`, ErrClientDataDuplicate},
+		{"J5-no-origin", `{"type":"webauthn.get","challenge":"abc"}`, ErrClientData},
+		{"J5-type-number", `{"type":1,"challenge":"abc","origin":"http://localhost:7644"}`, ErrClientData},
+		{"J6-crossOrigin-string", "{" + cdjBase + `,"crossOrigin":"false"}`, ErrClientData},
+		{"J6-topOrigin-bool", "{" + cdjBase + `,"topOrigin":true}`, ErrClientData},
+		{"J7-depth-5", "{" + cdjBase + `,"x":{"a":{"b":{"c":{"d":1}}}}}`, ErrClientData},
+	}
+	for _, a := range arms {
+		t.Run(a.name, func(t *testing.T) {
+			if c, err := ParseClientData([]byte(a.in)); !errors.Is(err, a.want) || c != (ClientData{}) {
+				t.Fatalf("(%+v, %v), want the zero value and %v", c, err, a.want)
+			}
+		})
+	}
+	for in, want := range map[string]ClientData{
+		"{" + cdjBase + `,"other_keys_can_be_added_here":"x","x":{"a":{"b":[1]}},"x":2}`:                {"webauthn.get", "abc", "http://localhost:7644", false, false},
+		`{"origin":"http://localhost:7644","crossOrigin":true,"challenge":"abc","type":"webauthn.get"}`: {"webauthn.get", "abc", "http://localhost:7644", true, false},
+		"{" + cdjBase + `,"topOrigin":"https://x.example"} `:                                            {"webauthn.get", "abc", "http://localhost:7644", false, true},
+	} {
+		if c, err := ParseClientData([]byte(in)); err != nil || c != want {
+			t.Errorf("%s: (%+v, %v), want %+v", in, c, err, want)
+		}
+	}
+	if in := "{" + cdjBase + `,"x":"` + strings.Repeat("a", 2048-len(cdjBase)-9) + `"}`; len(in) != 2048 {
+		t.Fatalf("the 2048-byte control is %d bytes", len(in))
+	} else if _, err := ParseClientData([]byte(in)); err != nil {
+		t.Errorf("2048 bytes: %v, want accepted", err)
+	}
+}
+
+// CD2: Chrome's clientDataJSON, as recorded.
+func TestClientDataReadsChrome(t *testing.T) {
+	for _, alg := range []string{"es256", "eddsa"} {
+		c := chromeFixture(t, alg)
+		for typ, b := range map[string][]byte{"webauthn.create": c.RegCDJ, "webauthn.get": c.AsrCDJ} {
+			want := ClientData{Type: typ, Challenge: "KRgdIk2YL2WdG9RwMjjT-EieJpaP_ioivXhvo78V_vU", Origin: c.Origin}
+			if got, err := ParseClientData(b); err != nil || got != want {
+				t.Errorf("%s %s: (%+v, %v)", alg, typ, got, err)
+			}
+		}
+	}
+}
+
+// OR1: the accepted origins and their rpID; every refusal with its error.
+func TestCeremonyOrigin(t *testing.T) {
+	long := strings.Repeat("a", 63) + "." + strings.Repeat("b", 63) + "." + strings.Repeat("c", 63) + "." + strings.Repeat("d", 54) + ".ts.net"
+	for in, rp := range map[string]string{
+		"http://localhost": "localhost", "http://localhost:7644": "localhost", "https://localhost:8443": "localhost",
+		"https://localhost": "localhost", "https://voights-mac-studio.tail97eda0.ts.net": "voights-mac-studio.tail97eda0.ts.net",
+		"https://x.ts.net:8443": "x.ts.net", "https://" + long: long, "https://x.ts.net:65535": "x.ts.net",
+	} {
+		if got, err := ParseCeremonyOrigin(in); err != nil || got != rp {
+			t.Errorf("%s: (%q, %v), want %q", in, got, err, rp)
+		}
+	}
+	for in, want := range map[string]error{
+		"https://x.ts.net:443": ErrOrigin, "https://localhost:443": ErrOrigin, "http://localhost:80": ErrOrigin,
+		"http://127.0.0.1:7644": ErrOriginHost, "https://evil.example": ErrOriginHost, "https://ts.net": ErrOriginHost,
+		"https://[::1]": ErrOrigin, "https://LOCALHOST": ErrOrigin, "https://x.ts.net/": ErrOrigin,
+		"https://x.ts.net:0": ErrOrigin, "https://x.ts.net:065535": ErrOrigin, "https://x.ts.net:65536": ErrOrigin,
+		"https://x.ts.net:+8443": ErrOrigin, "https://x.ts.net:": ErrOrigin, "https://a..ts.net": ErrOrigin,
+		"http://x.ts.net": ErrOrigin, "https://u@x.ts.net": ErrOrigin, "ftp://localhost": ErrOrigin, "localhost": ErrOrigin,
+		"https://" + strings.Replace(long, ".ts.net", "d.ts.net", 1): ErrOrigin,
+	} {
+		got, err := ParseCeremonyOrigin(in)
+		if !errors.Is(err, want) || got != "" || (want == ErrOrigin && errors.Is(err, ErrOriginHost)) {
+			t.Errorf("%s: (%q, %v), want %v", in, got, err, want)
+		}
+		if strings.HasSuffix(in, ":443") || strings.HasSuffix(in, ":80") {
+			if !strings.Contains(err.Error(), "default port") {
+				t.Errorf("%s: %v does not name the default-port rule", in, err)
+			}
+		}
+	}
+}
+
+// EN1: E2-E4, and Chrome's toJSON() output with its extra members.
+func TestPasskeyEnvelope(t *testing.T) {
+	const reg = `{"id":"AQI","rawId":"AQI","type":"public-key","response":{"clientDataJSON":"e30","attestationObject":"oA"}}`
+	for name, in := range map[string]string{
+		"E2-not-json":          `{`,
+		"E2-no-response":       `{"id":"AQI","rawId":"AQI","type":"public-key"}`,
+		"E2-no-id":             strings.Replace(reg, `"id":"AQI",`, "", 1),
+		"E2-no-rawId":          strings.Replace(reg, `"rawId":"AQI",`, "", 1),
+		"E2-no-type":           strings.Replace(reg, `"type":"public-key",`, "", 1),
+		"E2-no-clientDataJSON": strings.Replace(reg, `"clientDataJSON":"e30",`, "", 1),
+		"E2-null-attestation":  strings.Replace(reg, `"oA"`, "null", 1),
+		"E3-type-x":            strings.Replace(reg, `"public-key"`, `"x"`, 1),
+		"E3-id-not-rawId":      strings.Replace(reg, `"id":"AQI"`, `"id":"AQM"`, 1),
+		"E4-padded":            strings.Replace(reg, `"oA"`, `"oA=="`, 1),
+		"E4-non-canonical":     strings.Replace(reg, `"e30"`, `"e31"`, 1),
+		"E4-rawId-standard":    strings.ReplaceAll(reg, `"AQI"`, `"A+I"`),
+	} {
+		if r, err := ParseRegistrationResponse([]byte(in)); !errors.Is(err, ErrPasskeyResponse) || !reflect.DeepEqual(r, RegistrationResponse{}) {
+			t.Errorf("%s: (%+v, %v), want ErrPasskeyResponse", name, r, err)
+		}
+	}
+	if r, err := ParseRegistrationResponse([]byte(reg)); err != nil || !bytes.Equal(r.RawID, []byte{1, 2}) || string(r.ClientDataJSON) != "{}" || !bytes.Equal(r.AttestationObject, []byte{0xa0}) {
+		t.Fatalf("control: (%+v, %v)", r, err)
+	}
+	asr := `{"id":"AQI","rawId":"AQI","type":"public-key","response":{"clientDataJSON":"e30","authenticatorData":"AQ","signature":"Ag"%s}}`
+	for name, in := range map[string]string{
+		"E2-no-signature":         strings.Replace(asr, `,"signature":"Ag"`, "", 1),
+		"E2-no-authenticatorData": strings.Replace(asr, `"authenticatorData":"AQ",`, "", 1),
+		"E4-signature":            strings.Replace(asr, `"Ag"`, `"Ag=="`, 1),
+		"E4-userHandle":           strings.Replace(asr, "%s", `,"userHandle":"A="`, 1),
+		"E3-assertion-type":       strings.Replace(asr, `"public-key"`, `"x"`, 1),
+	} {
+		if r, err := ParseAssertionResponse([]byte(strings.Replace(in, "%s", "", 1))); !errors.Is(err, ErrPasskeyResponse) || r.RawID != nil {
+			t.Errorf("%s: (%+v, %v), want ErrPasskeyResponse", name, r, err)
+		}
+	}
+	for uh, want := range map[string]AssertionResponse{
+		"":                   {RawID: []byte{1, 2}, ClientDataJSON: []byte("{}"), AuthenticatorData: []byte{1}, Signature: []byte{2}},
+		`,"userHandle":null`: {RawID: []byte{1, 2}, ClientDataJSON: []byte("{}"), AuthenticatorData: []byte{1}, Signature: []byte{2}},
+		`,"userHandle":"bWFyaw"`: {RawID: []byte{1, 2}, ClientDataJSON: []byte("{}"), AuthenticatorData: []byte{1}, Signature: []byte{2},
+			UserHandle: []byte("mark"), HasUserHandle: true},
+	} {
+		if r, err := ParseAssertionResponse([]byte(strings.Replace(asr, "%s", uh, 1))); err != nil || !reflect.DeepEqual(r, want) {
+			t.Errorf("userHandle %q: (%+v, %v), want %+v", uh, r, err, want)
+		}
+	}
+	for _, alg := range []string{"es256", "eddsa"} {
+		c := chromeFixture(t, alg)
+		r, err := ParseRegistrationResponse(c.Registration)
+		if err != nil || !bytes.Equal(r.RawID, c.RegRawID) || !bytes.Equal(r.ClientDataJSON, c.RegCDJ) || !bytes.Equal(r.AttestationObject, c.AttObj) {
+			t.Fatalf("%s registration toJSON: (%+v, %v)", alg, r, err)
+		}
+		a, err := ParseAssertionResponse(c.Assertion)
+		if err != nil || !bytes.Equal(a.RawID, c.RegRawID) || !bytes.Equal(a.Signature, c.AsrSig) || a.HasUserHandle {
+			t.Fatalf("%s assertion toJSON: (%+v, %v)", alg, a, err)
+		}
+	}
+}
````

### Patch M3a — the registration law, the enrolment code and the prompt

````diff M3a
--- a/cmd/ailang-worldd/approver_enroll_test.go
+++ b/cmd/ailang-worldd/approver_enroll_test.go
@@ -227,7 +227,7 @@
 				}
 				return enrollArgs(db, "mark-mac", p)
 			}},
-		{name: "webauthn", setup: std, want: []string{"row 201 A2"}, checkRows: true,
+		{name: "webauthn", setup: std, want: []string{"--registration"}, checkRows: true,
 			args: func(t *testing.T, db string) []string {
 				return []string{"enroll", "--name", "mark-mac", "--format", "webauthn", "--pubkey", ed, "--db", db}
 			}},
--- a/host/approval/approver.go
+++ b/host/approval/approver.go
@@ -12,9 +12,9 @@
 )
 
 const (
-	// FormatSSH is the only approver format A1b enrols.
+	// FormatSSH is an SSH key approver (A1b).
 	FormatSSH = "ssh"
-	// FormatWebAuthn is named so it can be refused with a message naming row 201 A2.
+	// FormatWebAuthn is a passkey approver (row 201 A2): it enrols only through NewPasskeyApprover.
 	FormatWebAuthn = "webauthn"
 	// MaxKeyFileBytes caps a --pubkey file (one authorized_keys line).
 	MaxKeyFileBytes = 8 << 10
@@ -27,16 +27,17 @@
 // Approver is one enrolled approver as the pure law sees it.
 type Approver struct {
 	Name       string
-	Format     string       // FormatSSH
-	Key        SSHPublicKey // Type + wire Blob
+	Format     string       // FormatSSH or FormatWebAuthn
+	Key        SSHPublicKey // Type + wire Blob (ssh approvers)
 	EnrolledAt int64
-	LastUsedAt int64 // 0: never used
-	RevokedAt  int64 // 0: active
+	LastUsedAt int64    // 0: never used
+	RevokedAt  int64    // 0: active
+	Passkey    *Passkey // webauthn approvers only (row 201 A2); nil for ssh
 }
 
 var (
 	ErrApproverName        = errors.New("approval: approver name must match ^[a-z0-9][a-z0-9._-]{0,47}$")
-	ErrFormatNotYet        = errors.New("approval: webauthn approvers arrive with row 201 A2; A1b enrols ssh keys only")
+	ErrUseRegistration     = errors.New("approval: webauthn approvers enrol from a browser registration (approver enroll --format webauthn --origin … --registration …)")
 	ErrApproverFormat      = errors.New("approval: approver format must be ssh")
 	ErrApproverKeyFile     = errors.New("approval: the public key file must hold exactly one key line of at most 8 KiB")
 	ErrPrivateKeyFile      = errors.New("approval: that is a PRIVATE key file; pass the .pub file (its contents were not shown or kept)")
@@ -50,8 +51,8 @@
 	if !ApproverNamePattern.MatchString(name) { // E1
 		return Approver{}, ErrApproverName
 	}
-	if format == FormatWebAuthn { // E2 (before E3, so webauthn names A2)
-		return Approver{}, ErrFormatNotYet
+	if format == FormatWebAuthn { // E2 (before E3, so webauthn names its own path)
+		return Approver{}, ErrUseRegistration
 	}
 	if format != FormatSSH { // E3
 		return Approver{}, ErrApproverFormat
--- a/host/approval/approver_test.go
+++ b/host/approval/approver_test.go
@@ -147,8 +147,8 @@
 		{name: "name-49", call: ssh(strings.Repeat("a", 49), good), want: ErrApproverName},
 		{name: "name-escape", call: ssh("x\x1b[2J", good), want: ErrApproverName},
 		{name: "name-newline", call: ssh("mark\n", good), want: ErrApproverName},
-		{name: "webauthn", call: func() (Approver, error) { return NewApprover("mark", "webauthn", good, now) }, want: ErrFormatNotYet},
-		{name: "webauthn-empty-key", call: func() (Approver, error) { return NewApprover("mark", "webauthn", nil, now) }, want: ErrFormatNotYet},
+		{name: "webauthn", call: func() (Approver, error) { return NewApprover("mark", "webauthn", good, now) }, want: ErrUseRegistration},
+		{name: "webauthn-empty-key", call: func() (Approver, error) { return NewApprover("mark", "webauthn", nil, now) }, want: ErrUseRegistration},
 		{name: "format-gpg", call: func() (Approver, error) { return NewApprover("mark", "gpg", good, now) }, want: ErrApproverFormat},
 		{name: "format-empty", call: func() (Approver, error) { return NewApprover("mark", "", good, now) }, want: ErrApproverFormat},
 		{name: "format-SSH", call: func() (Approver, error) { return NewApprover("mark", "SSH", good, now) }, want: ErrApproverFormat},
@@ -302,24 +302,39 @@
 }
 
 // enrolmentPromptFields is the disposition of every Approver field (and every field of its
-// element type SSHPublicKey): shown on the enrolment y/N, or excluded with the reason.
+// element types SSHPublicKey and, row 201 A2, Passkey and COSEKey): shown on the enrolment y/N
+// (RenderEnrolment for ssh, RenderPasskeyEnrolment for Passkey.*), or excluded with the reason.
 var enrolmentPromptFields = map[string]string{
-	"Name":       "shown",
-	"Format":     "shown",
-	"Key.Type":   "shown",
-	"Key.Blob":   "shown as its SHA256 fingerprint",
-	"EnrolledAt": "excluded: a timestamp the store sets",
-	"LastUsedAt": "excluded: a timestamp the store sets",
-	"RevokedAt":  "excluded: a timestamp the store sets",
+	"Name":                   "shown",
+	"Format":                 "shown",
+	"Key.Type":               "shown",
+	"Key.Blob":               "shown as its SHA256 fingerprint",
+	"EnrolledAt":             "excluded: a timestamp the store sets",
+	"LastUsedAt":             "excluded: a timestamp the store sets",
+	"RevokedAt":              "excluded: a timestamp the store sets",
+	"Passkey.RPID":           "shown",
+	"Passkey.Origin":         "shown",
+	"Passkey.CredentialID":   "shown as the first 16 hex of its SHA-256",
+	"Passkey.Key.Alg":        "shown as ES256 or EdDSA",
+	"Passkey.Key.X":          "excluded: bound by the typed enrolment code, a fingerprint of the registration that carries it",
+	"Passkey.Key.Y":          "excluded: bound by the typed enrolment code, a fingerprint of the registration that carries it",
+	"Passkey.AAGUID":         "shown in hex",
+	"Passkey.BackupEligible": "shown as synced passkey, backup-eligible or device-bound",
+	"Passkey.BackupState":    "shown as synced passkey, backup-eligible or device-bound",
+	"Passkey.SignCount":      "shown",
 }
 
-// approverFieldPaths lists the dotted path of every leaf field of typ, walking struct fields.
+// approverFieldPaths lists the dotted path of every leaf field of typ, walking struct fields and
+// pointers to structs.
 func approverFieldPaths(typ reflect.Type, prefix string) []string {
 	var out []string
 	for i := 0; i < typ.NumField(); i++ {
 		f := typ.Field(i)
-		if f.Type.Kind() == reflect.Struct {
-			out = append(out, approverFieldPaths(f.Type, prefix+f.Name+".")...)
+		if ft := f.Type; ft.Kind() == reflect.Struct || (ft.Kind() == reflect.Pointer && ft.Elem().Kind() == reflect.Struct) {
+			if ft.Kind() == reflect.Pointer {
+				ft = ft.Elem()
+			}
+			out = append(out, approverFieldPaths(ft, prefix+f.Name+".")...)
 			continue
 		}
 		out = append(out, prefix+f.Name)
@@ -331,6 +346,9 @@
 func fieldByPath(a *Approver, path string) reflect.Value {
 	v := reflect.ValueOf(a).Elem()
 	for _, seg := range strings.Split(path, ".") {
+		if v.Kind() == reflect.Pointer {
+			v = v.Elem()
+		}
 		v = v.FieldByName(seg)
 	}
 	return v
@@ -356,9 +374,17 @@
 	base.LastUsedAt = 1500
 	base.RevokedAt = 1700
 	baseText := RenderEnrolment(base)
+	pkBase := passkeyWithFlags(t, 0x08) // backup-eligible, so BackupState's flip shows
+	pkText := RenderPasskeyEnrolment(pkBase)
 	for _, p := range paths {
-		b := base
+		b, render, text := base, RenderEnrolment, baseText
 		b.Key.Blob = append([]byte(nil), base.Key.Blob...)
+		if strings.HasPrefix(p, "Passkey.") {
+			pk := *pkBase.Passkey
+			pk.CredentialID, pk.Key.X, pk.Key.Y = bytes.Clone(pk.CredentialID), bytes.Clone(pk.Key.X), bytes.Clone(pk.Key.Y)
+			b, render, text = pkBase, RenderPasskeyEnrolment, pkText
+			b.Passkey = &pk
+		}
 		fv := fieldByPath(&b, p)
 		switch fv.Kind() {
 		case reflect.String:
@@ -369,10 +395,16 @@
 			fv.SetBytes(nb)
 		case reflect.Int64:
 			fv.SetInt(fv.Int() + 7)
+		case reflect.Uint32:
+			fv.SetUint(fv.Uint() + 7)
+		case reflect.Bool:
+			fv.SetBool(!fv.Bool())
+		case reflect.Array:
+			fv.Index(0).SetUint(fv.Index(0).Uint() ^ 0xff)
 		default:
 			t.Fatalf("field %s has unhandled kind %s", p, fv.Kind())
 		}
-		changed := RenderEnrolment(b) != baseText
+		changed := render(b) != text
 		shown := strings.HasPrefix(enrolmentPromptFields[p], "shown")
 		if shown && !changed {
 			t.Errorf("%s is %q but perturbing it does not change the prompt", p, enrolmentPromptFields[p])
--- a/host/approval/passkey.go
+++ b/host/approval/passkey.go
@@ -0,0 +1,160 @@
+package approval
+
+// Row 201 A2 (design §2.8): the registration law. A passkey approver is enrolled from one browser
+// registration (fmt "none": nothing is attested), checked against the enrol process's own
+// challenge and origin, and tied to the operator by the enrolment code he types (a 128-bit
+// fingerprint of the exact registration bytes) and his typed y. No I/O.
+
+import (
+	"bytes"
+	"crypto/sha256"
+	"encoding/base32"
+	"encoding/base64"
+	"encoding/hex"
+	"errors"
+	"fmt"
+)
+
+// Passkey is an enrolled passkey: what an assertion must match.
+type Passkey struct {
+	RPID, Origin   string
+	CredentialID   []byte
+	Key            COSEKey
+	AAGUID         [16]byte
+	BackupEligible bool   // stored; an assertion must agree (W13)
+	BackupState    bool   // shown at enrolment only (synced or device-bound); not stored
+	SignCount      uint32 // the registration's counter
+}
+
+// RegistrationInput is what NewPasskeyApprover judges.
+type RegistrationInput struct {
+	Name, Origin string
+	Challenge    [32]byte // the enrol process's fresh challenge
+	Response     []byte   // the registration file
+	Now          int64
+}
+
+var (
+	ErrClientDataType       = errors.New("approval: clientDataJSON has the wrong type for this ceremony")
+	ErrChallenge            = errors.New("approval: the response was made for another challenge")
+	ErrOriginMismatch       = errors.New("approval: the response was made at another origin")
+	ErrCrossOrigin          = errors.New("approval: the response was made in a cross-origin frame")
+	ErrAttestationObject    = errors.New("approval: the attestationObject is not a map of exactly fmt, attStmt and authData")
+	ErrAttestationFormat    = errors.New(`approval: the registration must use attestation "none" with an empty attStmt`)
+	ErrRPIDHash             = errors.New("approval: the authenticator signed for another relying party ID")
+	ErrUserPresence         = errors.New("approval: the authenticator did not report user presence (UP)")
+	ErrUserVerification     = errors.New("approval: the authenticator did not verify the user (UV: Face ID, Touch ID, a PIN)")
+	ErrCredentialIDMismatch = errors.New("approval: the response's rawId is not the passkey's credential ID")
+)
+
+func b64Challenge(c []byte) string { return base64.RawURLEncoding.EncodeToString(c) }
+
+// NewPasskeyApprover is the registration law (R1-R15, in this order; the first that fails is
+// returned). R16 (the algorithm) is ParseCOSEKey's Q3, reached through R11.
+func NewPasskeyApprover(in RegistrationInput) (Approver, error) {
+	if !ApproverNamePattern.MatchString(in.Name) { // R1
+		return Approver{}, ErrApproverName
+	}
+	rpID, err := ParseCeremonyOrigin(in.Origin) // R2
+	if err != nil {
+		return Approver{}, err
+	}
+	resp, err := ParseRegistrationResponse(in.Response) // R3
+	if err != nil {
+		return Approver{}, err
+	}
+	cd, err := ParseClientData(resp.ClientDataJSON) // R4
+	if err != nil {
+		return Approver{}, err
+	}
+	if cd.Type != "webauthn.create" { // R5
+		return Approver{}, ErrClientDataType
+	}
+	if cd.Challenge != b64Challenge(in.Challenge[:]) { // R6
+		return Approver{}, ErrChallenge
+	}
+	if cd.Origin != in.Origin { // R7
+		return Approver{}, ErrOriginMismatch
+	}
+	if cd.CrossOrigin || cd.HasTopOrigin { // R8
+		return Approver{}, ErrCrossOrigin
+	}
+	raw, err := attestedAuthData(resp.AttestationObject) // R9, R10
+	if err != nil {
+		return Approver{}, err
+	}
+	ad, err := ParseAuthData(raw, true) // R11
+	if err != nil {
+		return Approver{}, err
+	}
+	if ad.RPIDHash != sha256.Sum256([]byte(rpID)) { // R12
+		return Approver{}, ErrRPIDHash
+	}
+	if !ad.UP { // R13
+		return Approver{}, ErrUserPresence
+	}
+	if !ad.UV { // R14
+		return Approver{}, ErrUserVerification
+	}
+	if !bytes.Equal(resp.RawID, ad.CredentialID) { // R15
+		return Approver{}, ErrCredentialIDMismatch
+	}
+	return Approver{Name: in.Name, Format: FormatWebAuthn, EnrolledAt: in.Now, Passkey: &Passkey{
+		RPID: rpID, Origin: in.Origin, CredentialID: ad.CredentialID, Key: ad.Key, AAGUID: ad.AAGUID,
+		BackupEligible: ad.BE, BackupState: ad.BS, SignCount: ad.SignCount}}, nil
+}
+
+// attestedAuthData checks the attestationObject's shape (R9) and format (R10) and returns its
+// authData.
+func attestedAuthData(b []byte) ([]byte, error) {
+	v, err := decodeCBORExact(b, attestationLimits)
+	if err != nil { // R9
+		return nil, fmt.Errorf("%w: %w", ErrAttestationObject, err)
+	}
+	f, okF := v.get(cborValue{kind: cborText, b: []byte("fmt")})
+	st, okS := v.get(cborValue{kind: cborText, b: []byte("attStmt")})
+	ad, okA := v.get(cborValue{kind: cborText, b: []byte("authData")})
+	if len(v.items) != 6 || !okF || !okS || !okA || f.kind != cborText || st.kind != cborMap || ad.kind != cborBytes { // R9
+		return nil, ErrAttestationObject
+	}
+	if string(f.b) != "none" || len(st.items) != 0 { // R10
+		return nil, ErrAttestationFormat
+	}
+	return ad.b, nil
+}
+
+// EnrolmentCode is the 26-character, 128-bit fingerprint of a registration (design §2.8 [r1]):
+// base32 (A-Z, 2-7, unpadded) of the first 16 bytes of SHA-256("world-enrolment-code/v1" || 0x00
+// || SHA-256(clientDataJSON) || SHA-256(attestationObject)). It is public, not a secret.
+func EnrolmentCode(clientDataJSON, attestationObject []byte) string {
+	c, a := sha256.Sum256(clientDataJSON), sha256.Sum256(attestationObject)
+	h := sha256.Sum256(append(append([]byte("world-enrolment-code/v1\x00"), c[:]...), a[:]...))
+	return base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(h[:16])
+}
+
+// GroupEnrolmentCode shows an EnrolmentCode as XXXXX-XXXXX-XXXXX-XXXXX-XXXXXX.
+func GroupEnrolmentCode(code string) string {
+	return code[:5] + "-" + code[5:10] + "-" + code[10:15] + "-" + code[15:20] + "-" + code[20:]
+}
+
+// PasskeyEnrolmentWarning is what every passkey enrolment y/N shows, verbatim (design §2.8).
+const PasskeyEnrolmentWarning = "This passkey will approve mint requests with no terminal, each time after your device verifies\n" +
+	"you (Face ID, Touch ID, a PIN or a security-key touch). World cannot check which device holds\n" +
+	"it; the 26-character enrolment code you typed, which matched this exact registration, is what ties\n" +
+	"it to the browser you just used. Enrol only a passkey you created yourself just now.\n"
+
+// RenderPasskeyEnrolment is the body of a passkey enrolment y/N; a is NewPasskeyApprover's.
+func RenderPasskeyEnrolment(a Approver) string {
+	p := a.Passkey
+	backup := "device-bound"
+	switch {
+	case p.BackupEligible && p.BackupState:
+		backup = "synced passkey"
+	case p.BackupEligible:
+		backup = "backup-eligible"
+	}
+	cred := sha256.Sum256(p.CredentialID)
+	return fmt.Sprintf("approver  %s\nformat    %s\npasskey   %s for %s at %s\ncred      %s\naaguid    %s\nbackup    %s\ncounter   %d\nrecords   approved_by=%s:%s\n\n%s",
+		a.Name, a.Format, p.Key.AlgName(), p.RPID, p.Origin, hex.EncodeToString(cred[:8]), hex.EncodeToString(p.AAGUID[:]),
+		backup, p.SignCount, a.Format, a.Name, PasskeyEnrolmentWarning)
+}
--- a/host/approval/passkey_test.go
+++ b/host/approval/passkey_test.go
@@ -0,0 +1,205 @@
+package approval
+
+// Row 201 A2 sprint 1 M3 (design §6, PE1-PE4): the registration law over Chrome's recorded
+// registrations. A registration carries no signature (fmt "none"), so each refusal arm edits the
+// recorded bytes directly; the enrolment code is what binds them (PE4, sprint 3's WN10).
+
+import (
+	"bytes"
+	"crypto/sha256"
+	"encoding/base64"
+	"encoding/binary"
+	"encoding/hex"
+	"encoding/json"
+	"errors"
+	"reflect"
+	"regexp"
+	"strings"
+	"testing"
+)
+
+// regFile is a registration file around rawId, clientDataJSON and attestationObject.
+func regFile(rawID, cdj, attObj []byte) []byte {
+	e := base64.RawURLEncoding.EncodeToString
+	b, _ := json.Marshal(map[string]any{"id": e(rawID), "rawId": e(rawID), "type": "public-key",
+		"response": map[string]string{"clientDataJSON": e(cdj), "attestationObject": e(attObj)}})
+	return b
+}
+
+func attObj(fmtName string, attStmt cbMap, authData []byte) []byte {
+	return cb(cbMap{"fmt", fmtName, "attStmt", attStmt, "authData", authData})
+}
+
+// chromeRegistration is the recorded registration as NewPasskeyApprover's input.
+func chromeRegistration(t *testing.T, alg string) (chromeCeremony, RegistrationInput) {
+	t.Helper()
+	c := chromeFixture(t, alg)
+	var ch [32]byte
+	copy(ch[:], unhex(t, c.ChallengeHex))
+	return c, RegistrationInput{Name: "mark-mac", Origin: c.Origin, Challenge: ch, Response: c.Registration, Now: 1000}
+}
+
+func unhex(t *testing.T, s string) []byte {
+	b, err := hex.DecodeString(s)
+	if err != nil {
+		t.Fatal(err)
+	}
+	return b
+}
+
+// PE1: one arm per R1-R15 on an otherwise valid registration, each with the zero Approver.
+func TestNewPasskeyApproverRefusals(t *testing.T) {
+	c, base := chromeRegistration(t, "es256")
+	ad := regAuthData(t, c)
+	cdj := func(old, new string) []byte { return []byte(strings.Replace(string(c.RegCDJ), old, new, 1)) }
+	file := func(cd []byte, ao []byte) []byte { return regFile(c.RegRawID, cd, ao) }
+	arms := []struct {
+		name string
+		edit func(*RegistrationInput)
+		want error
+	}{
+		{"R1-name", func(in *RegistrationInput) { in.Name = "Mark" }, ErrApproverName},
+		{"R2-ip-origin", func(in *RegistrationInput) { in.Origin = strings.Replace(c.Origin, "localhost", "127.0.0.1", 1) }, ErrOriginHost},
+		{"R3-envelope", func(in *RegistrationInput) {
+			in.Response = bytes.Replace(c.Registration, []byte(`"public-key"`), []byte(`"x"`), 1)
+		}, ErrPasskeyResponse},
+		{"R4-client-data", func(in *RegistrationInput) { in.Response = file([]byte("{"), c.AttObj) }, ErrClientData},
+		{"R5-type-get", func(in *RegistrationInput) { in.Response = file(cdj("webauthn.create", "webauthn.get"), c.AttObj) }, ErrClientDataType},
+		{"R6-challenge", func(in *RegistrationInput) { in.Challenge[0] ^= 1 }, ErrChallenge},
+		{"R7-origin", func(in *RegistrationInput) { in.Origin = "http://localhost:7644" }, ErrOriginMismatch},
+		{"R8-cross-origin", func(in *RegistrationInput) {
+			in.Response = file(cdj(`"crossOrigin":false`, `"crossOrigin":true`), c.AttObj)
+		}, ErrCrossOrigin},
+		{"R8-top-origin", func(in *RegistrationInput) {
+			in.Response = file(cdj(`"crossOrigin":false`, `"crossOrigin":false,"topOrigin":"https://x.example"`), c.AttObj)
+		}, ErrCrossOrigin},
+		{"R9-trailing-byte", func(in *RegistrationInput) { in.Response = file(c.RegCDJ, append(bytes.Clone(c.AttObj), 0)) }, ErrCBORTrailing},
+		{"R9-array", func(in *RegistrationInput) { in.Response = file(c.RegCDJ, cb(cbArr{"none", cbMap{}, ad})) }, ErrAttestationObject},
+		{"R9-array-of-six", func(in *RegistrationInput) {
+			in.Response = file(c.RegCDJ, cb(cbArr{"fmt", "none", "attStmt", cbMap{}, "authData", ad}))
+		}, ErrAttestationObject},
+		{"R9-extra-key", func(in *RegistrationInput) {
+			in.Response = file(c.RegCDJ, cb(cbMap{"fmt", "none", "attStmt", cbMap{}, "authData", ad, "x", 1}))
+		}, ErrAttestationObject},
+		{"R9-authData-text", func(in *RegistrationInput) {
+			in.Response = file(c.RegCDJ, cb(cbMap{"fmt", "none", "attStmt", cbMap{}, "authData", "x"}))
+		}, ErrAttestationObject},
+		{"R10-packed", func(in *RegistrationInput) { in.Response = file(c.RegCDJ, attObj("packed", cbMap{}, ad)) }, ErrAttestationFormat},
+		{"R10-attStmt", func(in *RegistrationInput) { in.Response = file(c.RegCDJ, attObj("none", cbMap{"alg", -7}, ad)) }, ErrAttestationFormat},
+		{"R11-no-AT", func(in *RegistrationInput) {
+			in.Response = file(c.RegCDJ, attObj("none", cbMap{}, withFlags(ad, 0, 0x40)))
+		}, ErrAuthDataFlags},
+		{"R12-rp-id", func(in *RegistrationInput) {
+			in.Origin = "https://x.ts.net"
+			in.Response = file(cdj(c.Origin, in.Origin), c.AttObj)
+		}, ErrRPIDHash},
+		{"R13-no-UP", func(in *RegistrationInput) {
+			in.Response = file(c.RegCDJ, attObj("none", cbMap{}, withFlags(ad, 0, 0x01)))
+		}, ErrUserPresence},
+		{"R14-no-UV", func(in *RegistrationInput) {
+			in.Response = file(c.RegCDJ, attObj("none", cbMap{}, withFlags(ad, 0, 0x04)))
+		}, ErrUserVerification},
+		{"R15-raw-id", func(in *RegistrationInput) { in.Response = regFile(c.RegRawID[1:], c.RegCDJ, c.AttObj) }, ErrCredentialIDMismatch},
+	}
+	for _, a := range arms {
+		t.Run(a.name, func(t *testing.T) {
+			in := base
+			a.edit(&in)
+			if got, err := NewPasskeyApprover(in); !errors.Is(err, a.want) || !reflect.DeepEqual(got, Approver{}) {
+				t.Fatalf("(%+v, %v), want the zero Approver and %v", got, err, a.want)
+			}
+		})
+	}
+	if _, err := NewPasskeyApprover(base); err != nil {
+		t.Fatalf("control: %v", err)
+	}
+	if _, err := NewPasskeyApprover(RegistrationInput{Name: "mark-mac", Origin: base.Origin, Challenge: base.Challenge,
+		Response: file(c.RegCDJ, attObj("none", cbMap{}, ad)), Now: 1}); err != nil {
+		t.Fatalf("control (re-encoded attestationObject): %v", err)
+	}
+}
+
+// PE2: both Chrome registrations enrol, every Passkey field as recorded.
+func TestNewPasskeyApproverAcceptsChrome(t *testing.T) {
+	for _, alg := range []string{"es256", "eddsa"} {
+		c, in := chromeRegistration(t, alg)
+		k, _ := chromeKey(t, alg)
+		want := Approver{Name: "mark-mac", Format: FormatWebAuthn, EnrolledAt: 1000, Passkey: &Passkey{RPID: "localhost",
+			Origin: c.Origin, CredentialID: c.RegRawID, Key: k, AAGUID: [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8}, SignCount: 1}}
+		got, err := NewPasskeyApprover(in)
+		if err != nil || !reflect.DeepEqual(got, want) {
+			t.Fatalf("%s: (%+v %+v, %v), want %+v", alg, got, got.Passkey, err, want.Passkey)
+		}
+		// The same bytes re-made for a tailnet origin, synced (BE, BS), counter 7.
+		ad := withFlags(regAuthData(t, c), 0x18, 0)
+		h := sha256.Sum256([]byte("x.ts.net"))
+		ad = binary.BigEndian.AppendUint32(append(h[:], ad[32]), 7)[:37:37]
+		ad = append(ad, regAuthData(t, c)[37:]...)
+		in.Origin, in.Response = "https://x.ts.net:8443", regFile(c.RegRawID, []byte(strings.Replace(string(c.RegCDJ), c.Origin, "https://x.ts.net:8443", 1)), attObj("none", cbMap{}, ad))
+		want.Passkey.RPID, want.Passkey.Origin, want.Passkey.BackupEligible, want.Passkey.BackupState, want.Passkey.SignCount = "x.ts.net", in.Origin, true, true, 7
+		if got, err := NewPasskeyApprover(in); err != nil || !reflect.DeepEqual(got, want) {
+			t.Fatalf("%s tailnet: (%+v, %v), want %+v", alg, got.Passkey, err, want.Passkey)
+		}
+	}
+}
+
+// passkeyWithFlags is the ES256 registration with authData's backup flags set, enrolled.
+func passkeyWithFlags(t *testing.T, set byte) Approver {
+	c, in := chromeRegistration(t, "es256")
+	in.Response = regFile(c.RegRawID, c.RegCDJ, attObj("none", cbMap{}, withFlags(regAuthData(t, c), set, 0)))
+	a, err := NewPasskeyApprover(in)
+	if err != nil {
+		t.Fatal(err)
+	}
+	return a
+}
+
+// PE3: the enrolment y/N text, golden, and its backup labels.
+func TestRenderPasskeyEnrolmentGolden(t *testing.T) {
+	_, in := chromeRegistration(t, "es256")
+	a, err := NewPasskeyApprover(in)
+	if err != nil {
+		t.Fatal(err)
+	}
+	const want = "approver  mark-mac\nformat    webauthn\npasskey   ES256 for localhost at http://localhost:56547\n" +
+		"cred      914d608a1de4d7bf\naaguid    01020304050607080102030405060708\nbackup    device-bound\ncounter   1\n" +
+		"records   approved_by=webauthn:mark-mac\n\n" +
+		"This passkey will approve mint requests with no terminal, each time after your device verifies\n" +
+		"you (Face ID, Touch ID, a PIN or a security-key touch). World cannot check which device holds\n" +
+		"it; the 26-character enrolment code you typed, which matched this exact registration, is what ties\n" +
+		"it to the browser you just used. Enrol only a passkey you created yourself just now.\n"
+	if got := RenderPasskeyEnrolment(a); got != want {
+		t.Fatalf("RenderPasskeyEnrolment:\n%s\nwant:\n%s", got, want)
+	}
+	for set, label := range map[byte]string{0x18: "backup    synced passkey\n", 0x08: "backup    backup-eligible\n"} {
+		if got := RenderPasskeyEnrolment(passkeyWithFlags(t, set)); !strings.Contains(got, label) {
+			t.Errorf("flags %#x: %q lacks %q", set, got, label)
+		}
+	}
+	_, ed := chromeRegistration(t, "eddsa")
+	if e, err := NewPasskeyApprover(ed); err != nil || !strings.Contains(RenderPasskeyEnrolment(e), "passkey   EdDSA for localhost") {
+		t.Fatalf("EdDSA: %v", err)
+	}
+}
+
+// PE4: the enrolment code is 26 characters of A-Z2-7, pinned on the recorded registration, and
+// changes with one byte of either input.
+func TestEnrolmentCodeGolden(t *testing.T) {
+	c := chromeFixture(t, "es256")
+	code := EnrolmentCode(c.RegCDJ, c.AttObj)
+	if code != "YDOR4X36PQNFRTXBHWAV4W2XG4" || !regexp.MustCompile(`^[A-Z2-7]{26}$`).MatchString(code) {
+		t.Fatalf("EnrolmentCode = %q", code)
+	}
+	if g := GroupEnrolmentCode(code); g != "YDOR4-X36PQ-NFRTX-BHWAV-4W2XG4" {
+		t.Fatalf("GroupEnrolmentCode = %q", g)
+	}
+	for name, in := range map[string][2][]byte{
+		"clientDataJSON":    {append(bytes.Clone(c.RegCDJ[:len(c.RegCDJ)-1]), '!'), c.AttObj},
+		"attestationObject": {c.RegCDJ, append(bytes.Clone(c.AttObj[:len(c.AttObj)-1]), c.AttObj[len(c.AttObj)-1]^1)},
+		"swapped":           {c.AttObj, c.RegCDJ},
+	} {
+		if got := EnrolmentCode(in[0], in[1]); got == code || len(got) != 26 {
+			t.Errorf("a changed %s gives %q", name, got)
+		}
+	}
+}
````

