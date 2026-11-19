# Sprint plan — w-world-address-bar-verbs-v4, sprint S2 "the canonical action envelope and its key law" (iteration 408, row 243)

> **THE SPRINT — row 243 phase V4, sprint S2 (design §0, §1, §3.1, §3.2 Requested/Confirmed, §5.1 strict-JSON and
> request-limit rows, §5.2, §6, §7, §8 row S2, §9, §12 W15/W26/W36; S1 plan Appendix C); D-WORLD-119 and D-WORLD-120
> OPEN — S2 is answer-independent by design §8 ("Waits for: None. Pure law/adapters against synthetic fixtures") and by
> D-WORLD-119's default ("build S1 and S2 against fixtures only"):**
> **M1** — NEW `host/approval/operator_action.go`: the closed typed `world/operator-action/v1` envelope
> (`OperatorAction{Schema, ID, Kind, Episode, Store, Argv, Options, Resolution, Task, Review}`) canonicalizing S1's
> `OperatorIntent`, its strict decoder (over-limit, invalid UTF-8, NUL, non-object, duplicate key and trailing data refused
> **before** any decoding; then unknown key, type error, non-canonical bytes, invalid action), the review digest (semantic
> inputs, id blank) and the signature digest (final envelope), and the 97-byte message `world-operator-action/v1` NUL ‖
> digest ‖ challenge ‖ big-endian expiry under the NEW namespace `world-operator-action@sunholo`; a golden written by the real
> encoder; `TestOperatorStrictJSON`. **M2** — the 128 KiB / 8 KiB caps, every Validate rule, the two digests, and the join to
> S1's bar (`TestOperatorArgvIsTheBarsIntent`). **M3** — NEW `host/approval/operator_signed.go`: the pure key law
> `ConfirmOperatorAction` (the daemon's confirmation) and `CheckOperatorAction` (the runner's recheck from its own frozen
> action) returning the typed proof `*SignedOperatorAction` (unexported fields; only the two laws build one), with design
> §5.2's ConfirmBothKeys, ConfirmProof, ConfirmAuthenticity, ConfirmRevocationCounter, ConfirmBinding, ConfirmExpiry and the
> pure half of RunnerRechecks. **M4** — "pure entry points write nothing" (judge-407's retro rule), closing S1's survivor
> `J407-SETOUTPUT`. **No schema change, no daemon route, no `.ail`, `surfaceCSP` and `app.js` unchanged, no served or CLI
> behaviour change; the envelope and the law have no production caller; no pre-existing file changes.**
> - **Measured: 1,394 changed lines** — `git diff -M --numstat 8db22168 HEAD -- cmd host docs app` (insertions + deletions,
>   renames on; no rename occurs; **the same 1,394 with no path filter**; nothing lands outside `host/`), after a commit of
>   each of the four scratch stages: **499 production** (`operator_action.go` 335, `operator_signed.go` 164), **879 test**,
>   **16 golden**. **Cap 1,400: margin 6.**
> - **Cut (S2-D1).** Design 840 raw → measured **1,654** for everything (1.97×). **Moved out, compiled and green**: the
>   `host/authority` store adapters — `JudgeOperatorAction`, `RecheckOperatorAction` and the judge-only proof type
>   `VerifiedOperatorAction` with its AST construction test (**260 changed lines**,
>   `$HOME/.ailang/state/world-iter408/plan-scratch/patches/PARK-authority-judges.patch`, applies on this plan's M4; Appendix C
>   gives it to S3). What lands is the whole envelope and the whole law; what moves is the read-only store lookup around it.
> - **Four sub-milestones, each green at its commit** (§1.2 C14): **M1** envelope + strict decoder (530) → **M2** limits,
>   validation, digests, bar join (773) → **M3** the key law (1,304) → **M4** writes-nothing pins (1,394).
> - **Every line ships as a compiled patch** (Appendix P), extracted from this file and round-tripped onto a fresh
>   `git archive 8db22168` export, byte-identical to the stage trees (§1.2 C16). The executor applies, never retypes.
> - **Pre-existing files changed: 0; deleted: 0** (production or test). Every S2 path is added.

---

## 0. Binding constraints (each is restated where it governs)

- **D-WORLD-119 and D-WORLD-120 OPEN** (controller, iteration 408: `mission_decisions.sh --open` → both OPEN; `--check` valid,
  107 rows). The envelope and the law encode no answer: kind `start` and `finish` are both representable, as S1's
  `OperatorIntent` already is; which may be admitted stays S1's pure `OperatorPolicy` (zero value admits nothing). Nothing in
  S2 is reachable from a route, a store or the CLI.
- **No daemon route, no schema change, no `.ail`, `surfaceCSP` and `app.js` unchanged, no served-behaviour change**: the gates'
  `UNTOUCHED_VIOLATIONS` (every changed file other than a `_test.go`, the golden or the two NEW production files, plus any
  non-test file under `cmd`, `docs`, `host/daemon`, `host/store`, `host/authority`, `host/surfacepage`, `world`, `packages`,
  `scripts`, `tools`, `plugins` or any `.ail`) must say `none`; `DAEMON_ROUTES` base = final (47), `SURFACE_POSTS: 3`;
  `UNCHANGED:` app.js, `surface_page.go` (the CSP literal), `host/store/credentials.go` (the schema, `credentialsSchemaVersion =
  9`); `PREEXISTING_CHANGED: none`; `ENVELOPE_PRODUCTION_CALLERS: none` with a positive control.
- **D-WORLD-109**: no S2 test opens a store, a socket, a real key, `~/.ssh`, `~/.ailang/world`, `/dev/tty` or `gh`; every key
  is approvaltest's seeded test signer or passkey (`approvaltest.NewSigner(81)`, `NewPasskey(AlgEdDSA, 82)`, whose SSHSIG form
  is cross-checked against real `ssh-keygen` by the pre-existing `TestTestSignerAgreesWithSSHKeygen`, run in the gates' `inv`
  leg with `WORLD_REQUIRE_SSH_KEYGEN=1`); the only file a test writes is the golden, and only under
  `WORLD_UPDATE_OPERATOR_ACTION_GOLDEN=1` (never set by the executor or the gates).
- Coding standards: the envelope and the law are pure (data in, typed answer out; no I/O, clock or randomness — `now` is an
  input); host/approval imports the standard library only (`TestApprovalImportsOnlyTheStandardLibrary`, unedited, green).
- **Not touched**: `tools/launchd/*`, any skill, `world/*.ail`, `packages/`, `scripts/`, `cmd/` (all of it), `cmd/world-publish/*`,
  `plugins/`, `app.js`, `surface_page.go`, `daemon.go`, every `host/store`, `host/daemon`, `host/authority` file, every
  pre-existing file anywhere.

---

## 1. Pre-flight

### 1.1 Baseline at `8db22168` (measured by the planner, 2026-11-18)

`tools/gates.sh <base-tree> logs/base-gates base 8db22168` (the §6 script, byte-equal to the shipped one, in `base` mode on a
worktree of the planner's clone at `8db22168`; `GATES_SKIP_VERIFY_GO=1`), load ≈ 5–14 beside the planner's mutant shards (a
first run with the pre-final script at load ≈ 3: the same rcs, 96 s):

| leg | rc | time | note |
|---|---|---|---|
| `TOUCHED_PACKAGES` | — | — | *(empty at the base)* |
| `status`, `gofmt`, `build`, `vet` | 0 | 0 / 0 / 1 / 1 s | (cold: 4 / 3 s) |
| `verify_ail` | 0 | 37 s | `✓ verify gate PASSED: 16 required identities verified, 40 named tests pass, 503 se-tools named tests pass` — the final must match |
| `fence` | 0 | 0 s | `go test -run '^$' ./...` is the test-file compile fence (cold: 27 s) |
| `pkgs` | 0 | 3 s | `host/approval`, `host/verbs`, `host/authority` whole |
| `race` | 0 | 8 s | the same three whole under `-race` |
| `inv` | 0 | 5 s | **34 named `--- PASS`**, 0 SKIP, 0 FAIL (`WORLD_REQUIRE_SSH_KEYGEN=1`: the two `ssh-keygen` cross-checks ran) |
| `DAEMON_ROUTES` / `SURFACE_POSTS` | — | — | **47** (base 47) / **3** |
| `UNTOUCHED_VIOLATIONS` / `PREEXISTING_CHANGED` / `ENVELOPE_PRODUCTION_CALLERS` | — | — | none / none / none (control: `approval.CheckConfirmation` production callers 1) |
| `UNCHANGED` | — | — | app.js, surface_page.go, host/store/credentials.go |
| **`GATES_DONE 57s FAILS: none`** | | | `verify_go` is the controller's |

### 1.2 Codebase claims checked (at `8db22168` unless "work" — the planner's stage tree, tag `M4`)

| # | claim | command | observed |
|---|---|---|---|
| C1 | **Input 1 — the design's cites re-measured at `8db22168` (design measured at `41331b4d`): none moved** | `grep -n '^func \|^type ' host/authority/judge_verb.go`; `grep -n '^func \|humanMu\|commitBudget' host/daemon/surface_confirm.go; wc -l`; `git diff --stat 41331b4d 8db22168 -- host/authority host/approval host/daemon/surface_confirm.go host/daemon/surface_approve.go host/daemon/surface_approve_guard.go host/daemon/surface_page.go` | `judge_verb.go:30` `sshApprover`, `:44` `passkeyApprover`, `:63` `JudgeVerbSSH`, `:82` `JudgeVerbAssertion`, `:107` `RecheckVerb` — **identical** to W15/W36; `surface_confirm.go:97` `handleSurfaceConfirm`, `:129–131` `humanMu`/`commitBudget`, 215 lines (W26's 97–207 stands); the diff touches only `surface_page.go` (1 +/1 −, PR #604's attended change; `surfaceCSP` at `:19`, set at `:23`, W26's 19–23 stand). S1's symbols: `host/verbs/operator.go:48` `OperatorActionIDPattern = ^oa-[0-9a-f]{32}$`, `:85` `type OperatorIntent`, `:153` `fs.SetOutput(io.Discard)`; `host/store/credentials.go:45` `credentialsSchemaVersion = 9` |
| C2 | **Input 2 — the closest landed analogues, and what was mirrored** | `sed -n` of each | **Mirrored**: V3's verb envelope `host/approval/verb.go` (`VerbRequest` `:51`, `Validate`, `CanonicalVerb` `:120` — encoder, `SetEscapeHTML(false)`, no newline — `DecodeVerb` `:133` — `DisallowUnknownFields` + re-encode equality — `VerbMessage` `:152`, `VerbMessageHash` `:166`) → `OperatorAction`, `CanonicalOperatorAction`, `DecodeOperatorAction`, `OperatorActionMessage[Hash]`; V3's verb law `verb_signed.go` (`SignedVerb` `:18` with unexported fields, `ConfirmBySignature` `:74`, `ConfirmByAssertion` `:92`, `CheckConfirmation` `:121`, reusing `checkApprover` (`signed.go:75`, P2–P5) and `assertionSteps` (`passkey_assert.go:60`, W1–W14: credential, user handle, type, challenge, **origin W7**, cross-origin, **RP id W10**, **UP W11**, **UV W12**, BE, signature)) → `SignedOperatorAction`, `ConfirmOperatorAction`, `CheckOperatorAction`; `host/daemon/workbench_packet_evidence.go:61` `strictRosterJSON` (a token walk refusing duplicate keys and trailing values — the repo's only duplicate-key JSON helper, in a package host/approval may not import) → `strictObject`; the mint request's `stringLeaves` (`request.go:238`, UTF-8 and NUL per leaf) and errors `ErrSchema`, `ErrString`, `ErrEpisode`, `ErrPath`, `ErrDecode`, `ErrNotCanonical`, `ErrNegativeExpiry` (`request.go:78–90`), reused as is; the approvaltest keys and `paApprover` (`passkey_law_test.go:18`) for every test. **Differs**: (1) the decoder refuses the invalid classes **before** decoding, each with its own error (design §3.1), where `DecodeVerb` lets the re-encode comparison catch them all as one class; (2) the law takes the **action**, not a digest, and computes the signature digest itself, so a proof cannot be made for a digest that is not a valid action's; (3) the recheck checks expiry (V3's does not: its row is consumed) and that the confirmation's counter was recorded (R2, the pure half of design §3.1 "a counter already advanced at confirmation is the expected observation"); (4) one law per role with a `Format` input, as `CheckConfirmation` is, rather than one per key kind |
| C3 | **the package, by measurement** | `sed -n 1,80p host/approval/imports_test.go`; `grep -n pure host/daemon/surface_verb_test.go`; design §6 | design §6: "approval canonicalization and validators belong in `host/approval` and adapters in `host/authority`". `TestApprovalImportsOnlyTheStandardLibrary` allows host/approval any standard-library import but `os`, `os/exec`, `net`, `crypto/rand`, … — `operator_action.go` imports `bytes crypto/sha256 encoding/binary encoding/json errors fmt io path reflect regexp slices strings time unicode/utf8`, `operator_signed.go` `bytes errors fmt time`. host/approval cannot import host/verbs (not standard library), and host/verbs may not import host/approval (`TestSurfaceVerbIsReadOnly` pins its production imports to `errors flag fmt io regexp strings time unicode/utf8`), so the join to S1's `OperatorIntent` is a **test** in host/approval's external test package (S2-D4). The judges' `assertionSteps` and `checkApprover` are unexported in host/approval, so the law must live there; the store lookup (`ApproverByKey`, `PasskeyByCredentialID` on a read-only `*store.Credentials`) is host/authority's — the parked adapter |
| C4 | nothing named like S2's code exists | `grep -rnE '\b(OperatorAction\|ConfirmOperatorAction\|CheckOperatorAction\|SignedOperatorAction\|OperatorReviewDigest\|strictObject\|world-operator-action)\b' cmd host \| wc -l`; control `grep -rn 'func CheckConfirmation' host/approval` | `0`; control `verb_signed.go:121` |
| C5 | every type, field and function the patches name | `grep -n` at HEAD | **existing (reused)**: `stringLeaves`, `episodePattern`, `checkApprover`, `assertionSteps`, `ParseSSHSIG`, `VerifySSHSIG`, `Admit`, `AdmitInput{Now, CreatedAt, ExpiresAt, Consumed, StoredDigest, BytesDigest}`, `Admitted`, `RefuseExpired`, `NotAdmittedError{Verdict}`, `BindingError{Field}`, `ErrNotVerified`, `ErrVerbFormat`, `ErrSignCountRegression`, `Approver{Name, Format, Key, EnrolledAt, RevokedAt, Passkey}`, `Passkey{SignCount, CredentialID, Key}`, `FormatSSH`, `FormatWebAuthn`, `VerbPrefix`, `VerbNamespace`, `Namespace`, `DecisionNamespace`, `VerbMessage[Hash]`, `SignedMessage`, `VerbRequestSchema` (host/approval); `approvaltest.NewSigner/NewPasskey/Edit{Clear, CrossOrigin, Type, RawID, BadSignature}`, `paApprover`, `paOrigin`, `paRPID` (tests); `verbs.Split`, `AdmitOperator`, `OperatorPolicy`, `ValidOperatorActionID`, `admitLine`, `allYes` (host/verbs). **NEW** (host/approval): `OperatorActionSchema`, `OperatorActionPrefix`, `OperatorActionNamespace`, `OperatorActionMessageLen`, `OperatorActionLifetime`, `MaxOperatorActionBytes`, `MaxOperatorTaskBytes`, `OperatorOptions{Preset, Headless, ConfirmOnSurface, Ready}`, `OperatorResolved{Name, Value}`, `OperatorAction{…}`, `ErrOperatorAction`, `ErrOperatorTooLarge`, `ErrOperatorTask`, `ErrOperatorUTF8`, `ErrOperatorNUL`, `ErrOperatorShape`, `ErrOperatorDuplicateKey`, `ErrOperatorUnknownKey`, `ErrOperatorTrailing`, `OperatorArgv`, `(OperatorAction).Validate`, `CanonicalOperatorAction`, `DecodeOperatorAction`, `OperatorReviewDigest`, `OperatorSignatureDigest`, `OperatorActionMessage`, `OperatorActionMessageHash`, `SignedOperatorAction`, `OperatorConfirmInput`, `OperatorCheckInput`, `ErrOperatorLifetime`, `ErrOperatorCounter`, `ConfirmOperatorAction`, `CheckOperatorAction`, `(*SignedOperatorAction).Binds/ApprovedBy/PasskeyUse`; unexported `operatorActionIDPattern`, `operatorResolvedName`, `operatorResolvedMax`, `operatorPresetMax`, `validate`, `strictObject`, `operatorProof`; every test named in §3 |
| C6 | **Input 3 — fixtures from the real tool** | the test files | the golden `host/approval/testdata/operator-action.golden.json` (start and finish: canonical bytes, review digest, signature digest, message and message hash) is **written by the real encoder** — `WORLD_UPDATE_OPERATOR_ACTION_GOLDEN=1 go test ./host/approval/ -run '^TestOperatorActionGolden$'` — and read back; every hostile decoder input is the golden's bytes with one exact replacement (`strings.Count == 1` asserted); every SSHSIG is `approvaltest.Signer.Sign` (the repo's test signing path, byte-compared with `ssh-keygen -Y sign` by `TestTestSignerAgreesWithSSHKeygen`); every WebAuthn assertion is `approvaltest.Passkey.Assert` (Chrome's layout; the passkey is enrolled through the real `NewPasskeyApprover` by `paApprover`); the bar's tuples come from the real `verbs.Split` + `verbs.AdmitOperator`. No hand-typed bytes, digests or JSON |
| C7 | **Input 4 — tests that enumerate packages, imports, symbols, routes, writer sites, schema, CLI verbs, flags, help** | `grep -rn` by name; each run on the final tree (C17 `inv`) | **unchanged and green**: `TestApprovalImportsOnlyTheStandardLibrary` (host/approval production imports), `TestApprovalTestIsTestOnly`, `TestSurfaceVerbIsReadOnly` (host/verbs imports — S2 adds only a test there), `TestDaemonDependencyAllowlist`, `TestWorldBoundaryDependencyAllowlist`, `TestBareNetHTTPExemptionIsPerGroup`, `TestEveryCommandDirectoryIsAProtectedGroup`, `TestBoundaryASTWriteGuard`, `TestProductionContextRoots`, `TestProductionGoSurface`, `TestStoreSurfaceTakesContext`, `TestCredentialsSchemaMatchesCanonicalDDL`, `TestCredentialsSchema9VerbRequests`, `TestCredentialsNewerThanNineIsRefused`, `TestFrozenV1RouteTableMatchesMux`, `TestSurfaceAppRoutes`, `TestApprovalRoutesRegistered`, `TestDecisionRoutesRegistered`, `TestConfirmRoutesRegistered`, `TestDaemonCredentialWriterSites`, `TestSurfaceVerbRouteRegistered`, `TestSurfaceVerbOutcomes`, `TestEveryVerbHasHelp`, `TestVerbTableMatchesTheCLISwitch`, `TestCLICredentialWriterSites`, `TestOperatorFlagParity`. **Changed: none. New inventories: none** (the one AST construction test is in the parked patch) |
| C8 | **no route, no writer site, no schema change, no production caller — by test or caller search with a positive control** | the gates (C17) | `DAEMON_ROUTES: 47 (base 47) SURFACE_POSTS: 3`; `TestDaemonCredentialWriterSites`, `TestCLICredentialWriterSites` green unedited (no new writer); `UNCHANGED: host/store/credentials.go`; `ENVELOPE_PRODUCTION_CALLERS: none (control: approval.CheckConfirmation production callers 1)` — the search names all 21 NEW exported identifiers and excludes only the two files that define them |
| C9 | **S1's survivor reproduced on the base first** | `tools/mut.py <base-mut at 8db22168> base-whole J407-SETOUTPUT` | `J407-SETOUTPUT` (drop `fs.SetOutput(io.Discard)` in `AdmitOperator`): **SURVIVED** the whole `host/verbs` package at the base — judge-407's finding reproduced; **KILLED** on the final tree by `TestAdmitOperatorWritesNothing` (named and whole, §4) |
| C10 | **Input 5 — no asynchronous wait; time is an input** | the test files | no S2 test sleeps, ticks, polls or reads a clock: every law takes `Now`/`CreatedAt`/`ExpiresAt` as integers (the injected clock), and every time comparison's operator is its own mutant (`P1-ADMIT` via the pre-existing `Admit`, `C2-LIFETIME-AT/-OVER/-OFF`, `R1-EXPIRY-STRICT/-EARLY/-OFF`, `LIFETIME-CONST`, `C-EXPIRY-ROW`), each at the strict-equality boundary (`now == expiry` refused, `expiry − 1` admitted; lifetime 600 admitted, 601 refused). The writes-nothing tests use an `os.Pipe` read after the call returns |
| C11 | the new tests, work | `go test -count=1 -v` by name (`accept.sh`) | all `--- PASS`, none SKIP |
| C12 | **is the design's §3.1 decode order sound against `encoding/json`?** | work: `TestOperatorStrictJSON` | yes, and it needs the walk: `encoding/json` keeps the **last** of a repeated key (`{"kind":"start","kind":"start"}` decodes), reads only the **first** value (trailing `{}` is ignored by `Decode`), replaces invalid UTF-8 with U+FFFD, and matches keys **case-insensitively** (`"Kind"` fills `Kind`) — each would otherwise reach canonicalization. The walk refuses the first two by class; UTF-8 and NUL are byte checks before it; case folding, key order, spacing, HTML escaping and `null` for a string are caught by the canonical comparison (`ErrNotCanonical`) |
| C13 | the snapshot block (§7) on every sub-milestone | `tools/snap_rehearse.sh` (a worktree at `8db22168`, the patches in order, §7's block verbatim) | `M1 3/0/3`, `M2 4/0/4`, `M3 6/0/6`, `M4 8/0/8` (`STATUS` lines / `DELETED` lines / MANIFEST lines), each `MANIFEST_OK`, rc 0; every stage has only `host/` (the `find` guard skips the absent `cmd/` and `docs/`); the block extracted from this file is §7's verbatim but for `M`, `W` and the snapshot root (`BLOCK_VERBATIM_BUT_W_M_D`) — `logs/snap.txt`, re-run on this final plan file |
| C14 | **every sub-milestone is green when committed** | `tools/accept.sh` on each stage tree (`st-M1`…`st-M4` at the stage tags) | on the final stage trees (`logs/accept/*.txt`): `M1 fast rc=0 ACCEPT_OK` (15 named `--- PASS`), `M2` (19), `M3` (27), `M4` (31), each 5–7 s warm; cold, at load ≈ 10–14 beside five mutant shards, 34–94 s each — every part far inside one 600 s call. **Measured and changed before shipping**: (1) the first M1 draft was 769 lines alone and the whole sprint 1,654, so the adapters were cut (S2-D1) and the M1 test split into M1/M2; (2) the first full stage set measured **1,419**, then 1,398, then (after mutant run 1 added four rows) 1,401 — each over or at the cap, each brought back by compacting constants, comments and the two writes-nothing harnesses, never by removing a refusal row; (3) mutant run 1 left five non-equivalent whole-package survivors (`C-EXPIRY-ROW`, `SSH-MSG-EXPIRY`, `O11-RES-MAX`, `SSH-MSG-ERR`, `W15-ZERO-AD`), each answered by a new row (§4's preamble) |
| C15 | node, the pinned binaries | `$HOME/.pinned-ailang/ailang --version`; `go version`; `which ssh-keygen` | `AILANG v0.41.0`; `go1.26.6 darwin/arm64`; `/usr/bin/ssh-keygen` (the `inv` leg's `TestSSHSIGCrossCheckWithSSHKeygen` and `TestTestSignerAgreesWithSSHKeygen` run, not skip: `WORLD_REQUIRE_SSH_KEYGEN=1`) |
| C16 | **the round trip** | `tools/roundtrip.sh <this plan>` | `CMP_EQUAL` for all four (each block extracted from this file with §2's awk equals the generated patch byte for byte); each applied with `/usr/bin/patch -p1 -E` after a `--dry-run` to a fresh `git archive 8db22168` export, no `.orig`/`.rej` under `cmd host`, committed and measured: `NUMSTAT M1 530`, `M2 773`, `M3 1304`, `M4 1394` (each `unfiltered=` the same); after each, `diff -r` of `cmd` and `host` against `st-<M>` rc 0; then `diff -r -x .git` of the whole export against `st-M4` → **`ROUNDTRIP_IDENTICAL_TO_STAGE_TREE`**; `SCRIPT_EQUAL accept.sh`, `SCRIPT_EQUAL gates.sh` (the scripts extracted from this file are the ones that produced C14 and C17) — `logs/roundtrip.txt`, re-run on this final plan file. The controller's `ctl_roundtrip.sh` was read, not run (it refuses a reused export dir); its stage trees are `plan-scratch/st-M1` … `st-M4` |
| C17 | **the gates on the final tree** | `tools/gates.sh <final-tree> logs/final-gates final <sha>` | `tools/gates.sh <final-tree> logs/final-gates final fa8db2ad` (`final-tree` = a worktree of the planner's clone at the stage tag `M4`, the four stage commits on `8db22168`; `diff -r` equal to the round-trip export, C16; `GATES_SKIP_VERIFY_GO=1`): **`TOUCHED_PACKAGES: host/approval host/verbs`** (= the plan's list); `status` 0; **`verify_ail` rc 0** (the base's PASSED line, byte for byte); `gofmt`, `build`, `vet`, `fence` rc 0; **`pkgs` rc 0** (`host/approval` 0.5 s, `host/verbs` 0.3 s, `host/authority` 0.9 s); **`race` rc 0** (2.4 s, 1.4 s, 6.8 s); **`inv` rc 0: 51 named `--- PASS`** (the base's 34 + the sixteen host/approval `TestOperator*` + `TestAdmitOperatorWritesNothing`), 0 SKIP, 0 FAIL, no `NOT PASSED`; **`DAEMON_ROUTES: 47 (base 47)`, `SURFACE_POSTS: 3`, `UNTOUCHED_VIOLATIONS: none`, `PREEXISTING_CHANGED: none`, `ENVELOPE_PRODUCTION_CALLERS: none (control: approval.CheckConfirmation production callers 1)`, `UNCHANGED:` app.js, surface_page.go, host/store/credentials.go**; **`GATES_DONE 65s FAILS: none`**. `verify_go` is the controller's |

### 1.3 Planner decisions where the design or the directive leaves a choice

- **S2-D1 — the cut, by measurement.** Compiled: M1 530, M2 +243, M3 +531, M4 +90 = **1,394**; the authority adapters +260 =
  1,654 > 1,400. The envelope lands first (the law needs it), then the law (design §5.2's rows are its branches), then the
  writes-nothing pins (a retro rule). The adapters are the cut: they are a store lookup around the law (`ApproverByKey`,
  `PasskeyByCredentialID`, V3's `sshApprover`/`passkeyApprover` reused unchanged) and their only caller is S3/S7's confirm
  path, which needs S3's row type anyway. The parked patch is compiled, gofmt-clean, vetted and green on M4 (its three tests
  `--- PASS`; `go test ./host/authority/` ok). No refusal-branch test was cut: every test S2 lands kills its rows.
- **S2-D2 — four sub-milestones**, each independently green (C14): **M1** the envelope with its golden, message layout and
  strict-JSON tests (production lands with its decoder tests); **M2** test-only, the caps, Validate and digests; **M3** the law
  with its six §5.2 tests and the recheck law; **M4** test-only, writes-nothing. M1 (530) and M3 (531) are over the directive's
  ~250–400 guide because splitting either from its tests would land untested production.
- **S2-D3 — the package is `host/approval`, by measurement** (C3): the reused primitives are unexported there; the store
  lookup is `host/authority`'s (parked). The proof type `SignedOperatorAction` has unexported fields, so only
  `ConfirmOperatorAction` and `CheckOperatorAction` return one with `ok` set; the zero value and nil bind nothing
  (`TestOperatorConfirmBothKeys`). The parked `authority.VerifiedOperatorAction` adds what the audit asked for — a proof that
  can only come from a judge that read the approver from the enrolled store — with an AST test (one constructor literal,
  called only by the two judges, positive control) and a reflect test (every field unexported).
- **S2-D4 — the envelope canonicalizes S1's `OperatorIntent` field for field** (`Kind`, `Episode`, `Options{Preset, Headless,
  ConfirmOnSurface, Ready}`, `Argv`) and Validate requires `Argv == OperatorArgv(Kind, Episode, Options)` exactly — the
  command words are a function of the admitted options, so the two cannot disagree. `TestOperatorArgvIsTheBarsIntent` runs
  five real bar lines through `verbs.Split` + `verbs.AdmitOperator` and requires `OperatorArgv` to equal each `Intent.Argv` and
  the built action to validate; it also requires the envelope's id check to agree with `verbs.ValidOperatorActionID` on seven
  boundary ids. The bar line is never re-parsed.
- **S2-D5 — fields not yet resolved are typed, closed and bounded, not invented**: `Store` is the served credentials store's
  absolute clean path (design §2.2 "the canonical daemon store pair"; S4 records it); `Resolution` is the frozen resolution as
  1–64 `{name, value}` facts in strictly ascending name order (`^[a-z][a-z0-9_]{0,63}$`) so one resolution has one encoding —
  S5 names the facts (design §4.2's repository, base OID, preset digest, agent/model, …) without an envelope schema change;
  `Task` is 1–8,192 UTF-8 bytes on a start and empty on a finish; `Review` is the exact displayed review text (design §3.1
  "displayed review bytes"), non-empty.
- **S2-D6 — two digests, named and distinct**: `OperatorReviewDigest` = SHA-256 of the canonical bytes with `id` blank (validates
  everything but the id; S3's dedup key, design §3.1 "before assigning an action id or challenge"); `OperatorSignatureDigest` =
  SHA-256 of the final envelope (validates fully; what the message signs). `TestOperatorDigestsAreDistinct`: they differ for one
  action; another id moves only the signature digest; each semantic field moves both.
- **S2-D7 — the law's order** (ConfirmOperatorAction): the action's own signature digest (an invalid action has none), P1
  `Admit` (consumed, clock behind, `now ≥ expiry`, stored bytes ≠ digest column), C1 the row's digest column is **this**
  action's, C2 lifetime ≤ 600 s, then (operatorProof) the format, P2–P5 at now, the SSHSIG over the operator message under
  the operator namespace, or W1–W14 and W15 in full. `CheckOperatorAction` (the runner): its own frozen action's digest, R1
  `now ≥ expiry` refused, P2–P5 at now (revoked since → refused), the signature, and R2 — the stored counter must be at or past
  the artifact's (the daemon recorded it), never W15 again.
- **S2-D8 — S1's survivor is closed here** (cheap; same rule): `TestAdmitOperatorWritesNothing` in host/verbs, and the same
  rule for every new pure entry point (`TestOperatorPureEntryPointsWriteNothing`), each with a positive control proving the
  pipe sees a write. Judge-407 note (3), numeric flags typed as strings: **not acted on** — the envelope carries no numeric
  option (the admitted subset has none), so the typing does not reach it; carried in Appendix C.
- **S2-D9 — no Chrome, no node, no docs**: S2 adds no page, route, command, flag or output line. The gates have no `chrome` leg.

### 1.4 Premises found false or incomplete at HEAD (each routed around minimally)

| # | premise | measured | routing |
|---|---|---|---|
| F1 | directive: design cites `judge_verb.go:30,44,63,82,107` and `surface_confirm.go:97–207` may have moved | C1: none moved (#604 changed `surface_page.go` by one line, not the CSP's lines' meaning; `surfaceCSP` still `:19`) | none; cites stand |
| F2 | design §6 / audit: "typed verified-action proof … in host/authority with unexported constructors" | the primitives the law needs are unexported in host/approval (C3), so the **law and its proof type** must be host/approval's (as V3's `SignedVerb`); only the store adapter can be host/authority's | S2 lands the law's proof (`SignedOperatorAction`); the host/authority wrapper is the parked patch (S2-D1, S2-D3) |
| F3 | design §3.1 lists "duplicate/unknown keys, trailing JSON, invalid UTF-8, NUL" as if `encoding/json` refused them | it refuses none of duplicate keys, trailing values (after `Decode`), invalid UTF-8 (replaced), case-folded keys (C12) | the token walk and byte checks before decoding; the canonical comparison after |
| F4 | design §3.1 "challenge of 32 random bytes" in the envelope | the challenge is in the **message**, not the canonical action (V3's verb request likewise keeps its challenge in the row) | the envelope has no challenge; the law takes it as an input; binding is checked by `Binds(action, challenge, expiry)` |
| F5 | design §5.2 `TestOperatorConfirmRevocationCounter` "replayed/unchanged nonzero counter" | the reused W15 `(ad ≠ 0 ∨ p ≠ 0) ∧ ad ≤ p` is equivalent to `p ≠ 0 ∧ ad ≤ p` (if `p = 0`, `ad ≤ 0` forces `ad = 0`) — the same redundancy V3's `ConfirmByAssertion` carries | mirrored unchanged; mutant `W15-ZERO-P` is a named equivalent with this proof (§4) |
| F6 | directive: "fixtures … `ssh-keygen -Y sign` with a TEST key" | the repo's test signer is byte-compared with `ssh-keygen -Y sign` by `TestTestSignerAgreesWithSSHKeygen` (approvaltest), so its signatures are the real tool's bytes; a second `ssh-keygen` run per test would add subprocess time to every law test | approvaltest's signer; the cross-check runs in the gates' `inv` leg with `WORLD_REQUIRE_SSH_KEYGEN=1` |

**Files outside design §7**: none — §7 names `host/{approval,authority,store,verbs,daemon,surfacepage}`; S2 adds files to
`host/approval` and one test to `host/verbs`.

### 1.5 The size and the race-mode cost

| stage | production | test | changed lines | cumulative (numstat after commit) |
|---|---|---|---|---|
| M1 | 335 (`operator_action.go`) | 179 (`operator_action_test.go`) + 16 (golden) | 530 | **530** |
| M2 | 0 | 243 (`operator_action_valid_test.go`) | 243 | **773** |
| M3 | 164 (`operator_signed.go`) | 367 (`operator_signed_test.go`) | 531 | **1,304** |
| M4 | 0 | 55 (`operator_silent_test.go`) + 35 (`host/verbs/operator_silent_test.go`) | 90 | **1,394** |
| *parked* | 108 (`host/authority/judge_operator.go`) | 152 (`judge_operator_test.go`) | 260 | *(1,654)* |

Design 840 raw (envelope 360, proof and recheck 480) → compiled **1,394 landed (1.66×)**, 1,654 with the adapters (1.97×): the
directive's boundary-class matrix (a mutant per class, each killed) is most of the test weight.

**Race-mode wall time of the new tests** (directive 8): work tree, `go test -race -count=1 -v` (`logs/race-new.txt`, `logs/race-base.txt`): `host/approval` — the sixteen `TestOperator*` 0.20 s together (the slowest `TestOperatorRequestLimits` 0.05 s, `…ConfirmBinding` 0.04 s); package 1.9 s against 1.7 s at the base; `host/verbs` — `TestAdmitOperatorWritesNothing` < 0.01 s, package 1.3 s against 1.7 s at the base (noise). **≈ 0.2 s added in all**, in packages that take ≈ 2 s; host/daemon and host/broker gain nothing (no S2 test there), so CI's 8-minute bound per package is untouched. No key generation at run time: approvaltest's keys are derived from a seed byte, once per package-level `var`

---

## 2. How the executor applies a sub-milestone (the patch protocol)

Each sub-milestone `<M>` is the block fenced by a line `` ````diff <M> `` (four backticks) and a line of four backticks,
in Appendix P. Paths are `a/<repo path>` and `b/<repo path>`; every S2 file is new (`--- /dev/null`); **no sub-milestone
deletes or changes a pre-existing file**. Apply it **only** this way, in one `bash` shell (the controller's shell is zsh):

```bash
bash <<'EOF2'
set -euo pipefail
M=M1   # <- the sub-milestone, exactly as in the Appendix P heading, in order: M1 M2 M3 M4
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter408-row243-v4s2
PLAN="$W/design_docs/planned/w-world-address-bar-verbs-v4s2-sprint-plan.md"
T="$HOME/.ailang/state/world-iter408/exec-patches"; mkdir -p "$T"
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
  `bash "$T/accept.sh" "$W" <M> fast 2>&1 | tee "$HOME/.ailang/state/world-iter408/exec-logs/accept_<M>_fast.txt"` (never
  only through `tail`). **Never `nohup`, never background a test leg; pass a 600000 ms timeout to every long Bash call.**
  No part launches Chrome or regenerates the golden.
- Then snapshot (§7). The controller commits and compares the numstat with §3's expected value.

---

## 3. Milestones

**Acceptance.** Every part runs the same `[base]` checks first (gofmt over `cmd host`, build, vet of host/approval, host/verbs
and host/authority, the test-file compile fence `go test -run '^$' ./...`), then the named tests with `-v` (every name `---
PASS`, none SKIP or FAIL), then host/approval, approvaltest, host/verbs and host/authority whole, then three host/daemon
inventories. One part, `fast`, per sub-milestone (C14 measured each, well inside one 600 s call).

````bash accept.sh
#!/bin/bash
# accept.sh <tree> <M> [fast]: exactly plan §3's acceptance for sub-milestone M (one bash shell; rc 0 =
# ACCEPT_OK). One part per sub-milestone, each well inside one 600 s call. No part launches Chrome (S2 has
# no Chrome test) and no part regenerates the golden.
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
step go vet ./host/approval/ ./host/verbs/ ./host/authority/
if go test -run '^$' ./... 2>&1 | grep -E '^(FAIL|--- FAIL)'; then echo "ACCEPT_RED $M $P: FENCE_RED"; exit 1; fi
A1='TestOperatorActionGolden|TestOperatorActionMessage|TestOperatorStrictJSON'
A2='TestOperatorRequestLimits|TestOperatorActionValidate|TestOperatorDigestsAreDistinct|TestOperatorArgvIsTheBarsIntent'
A3='TestOperatorConfirmBothKeys|TestOperatorConfirmProof|TestOperatorConfirmAuthenticity|TestOperatorConfirmRevocationCounter|TestOperatorConfirmBinding|TestOperatorConfirmExpiry|TestOperatorConfirmNegativeExpiry|TestOperatorRunnerRecheckLaw'
A4='TestOperatorPureEntryPointsWriteNothing'
INV='TestApprovalImportsOnlyTheStandardLibrary|TestVerbDomainSeparation|TestSSHSIGCrossCheckWithSSHKeygen'
case "$M:$P" in
M1:fast) step v ./host/approval/ -run "^($A1|$INV)\$" ;;
M2:fast) step v ./host/approval/ -run "^($A1|$A2|$INV)\$" ;;
M3:fast) step v ./host/approval/ -run "^($A1|$A2|$A3|$INV)\$" ;;
M4:fast) step v ./host/approval/ -run "^($A1|$A2|$A3|$A4|$INV)\$"
         step v ./host/verbs/ -run '^(TestAdmitOperatorWritesNothing|TestOperatorGrammarOptions|TestOperatorGrammarArity)$' ;;
*) echo "unknown $M $P"; exit 2 ;;
esac
step p ./host/approval/ ./host/approval/approvaltest/ ./host/verbs/ ./host/authority/
step v ./host/daemon/ -run '^(TestSurfaceVerbIsReadOnly|TestDaemonDependencyAllowlist|TestConfirmRoutesRegistered)$'
echo "ACCEPT_OK $M $P"
````

### M1 — the envelope and its strict decoder · 335 prod + 195 test/golden · part: `fast`

- `host/approval/operator_action.go` (NEW, S2-D4 … D6): the constants (schema, prefix, namespace, message length 97, lifetime
  600 s, caps 128 KiB / 8 KiB), `OperatorOptions`, `OperatorResolved`, `OperatorAction`, the nine errors, `OperatorArgv`,
  `Validate` (O1 schema, O2 every string leaf UTF-8 and NUL-free, O3 id, O4 kind, O5 episode, O6 store absolute and clean,
  O7 preset — not `se-tools`, ≤ 51, the episode grammar, O8 kind/option consistency, O9 argv exactly `OperatorArgv`, O10 task
  1–8,192 bytes on a start / empty on a finish, O11 resolution 1–64 strictly ascending valid names, O12 review non-empty, O13
  canonical ≤ 128 KiB), `CanonicalOperatorAction`, `DecodeOperatorAction` (D1 size, D2 UTF-8, D3 NUL, D4 object, D5 duplicate
  key at any depth, D6 trailing, D7 unknown key / type error, D8 canonical bytes, D9 Validate), `strictObject`, the two digests,
  `OperatorActionMessage[Hash]`.
- `host/approval/operator_action_test.go` (NEW): `opAction` (a start and a finish, the task holding `$(…)`, a backtick, `<b>`
  and `é`), `opGoldenOf`, `opGoldenBytes`, **`TestOperatorActionGolden`** (the two actions' canonical bytes, digests, message
  and hash equal the golden; the golden decodes back to each action), **`TestOperatorActionMessage`** (97 bytes; prefix, NUL at
  24, digest 25–56, challenge 57–88, big-endian expiry 89–96; negative expiry refused by both), **`TestOperatorStrictJSON`** (29
  hostile inputs, each the golden with one exact replacement or a shape: duplicate key top/nested/array element, unknown key
  top/nested, trailing object/word/newline, invalid and overlong UTF-8, raw and escaped NUL, null/array/string/number/empty,
  unclosed, five type confusions, case-folded key, reordered keys, spacing, HTML escape, canonical-but-invalid — each its own
  error class and the zero action).
- `host/approval/testdata/operator-action.golden.json` (NEW, written by the encoder, C6).
- Expected numstat after commit **530**.

### M2 — limits, validation, digests, the bar join · 0 prod + 243 test · part: `fast`

- `host/approval/operator_action_valid_test.go` (NEW): **`TestOperatorArgvIsTheBarsIntent`** (S2-D4), **`TestOperatorRequestLimits`**
  (the caps' values; canonical bytes exactly 128 KiB valid and decoded, one under valid, one over refused by Validate and by
  the decoder; one over by a single space refused **before any parse** as too large; at the cap with a space not canonical;
  task 8,192 bytes of two-byte runes valid, 8,193 refused, 8,191 + a two-byte rune refused, 1 valid, 0 refused on a start, any
  on a finish refused), **`TestOperatorActionValidate`** (37 rules refused one at a time, each with its error and the signature
  digest refused too; 1 and 64 facts and a 51-character preset valid), **`TestOperatorDigestsAreDistinct`** (S2-D6), and
  `opFieldEdits` (seven semantic edits that keep an action valid, shared with M3).
- Expected numstat after commit **773**.

### M3 — the key law · 164 prod + 367 test · part: `fast`

- `host/approval/operator_signed.go` (NEW, S2-D7): `SignedOperatorAction` (unexported fields), `OperatorConfirmInput`,
  `OperatorCheckInput`, `ErrOperatorLifetime`, `ErrOperatorCounter`, `ConfirmOperatorAction`, `CheckOperatorAction`,
  `operatorProof`, `Binds`, `ApprovedBy`, `PasskeyUse`.
- `host/approval/operator_signed_test.go` (NEW): `opIn` (a valid confirmation by mark-ssh or mark-mac: created 1000, expires
  1600, now 1010, stored counter 1, assertion 2), `opRefused` (asserts the error class **and no proof**), and
  **`TestOperatorConfirmBothKeys`**, **`…Proof`** (typed y, no proof, malformed armor/assertion, the wrong artifact for the
  format, no or a typed format, bad challenge both kinds, invalid action), **`…Authenticity`** (another key, forged assertion,
  unknown credential; five wrong namespaces incl. V3's, A3's, A4's and `git`; no UV, no UP, another RP id, another origin,
  cross-origin, a registration), **`…RevocationCounter`** (revoked both kinds, enrolled after now, enrolled after the request
  but before now **admitted**, a non-passkey approver, counters {1,1} {5,2} {5,0} refused, {1,2} {0,0} {0,1} admitted),
  **`…Binding`** (each of nine fields of the action and the expiry, both kinds; the row's digest column of another action; a
  V3 verb and an A3 mint SSHSIG over the same digest; the verb message under the operator namespace; a V3 verb assertion),
  **`…Expiry`** (both kinds: `now == expiry`, past, `expiry − 1`, before creation, at creation, lifetime 600, 300 and 601,
  consumed, stored bytes differ), **`TestOperatorConfirmNegativeExpiry`**, **`TestOperatorRunnerRecheckLaw`** (both kinds: the
  control with the counter recorded; `now == expiry`; revoked since; another frozen action; another challenge; an invalid
  action; no format; stored counter 1 behind the assertion's 2 refused, 2 and 3 admitted).
- Expected numstat after commit **1,304**.

### M4 — pure entry points write nothing · 0 prod + 90 test · part: `fast`

- `host/approval/operator_silent_test.go` (NEW): `silent` (os.Stdout, os.Stderr and the log on one pipe) and
  **`TestOperatorPureEntryPointsWriteNothing`** (control: a write to each is seen; then the decoder on valid, duplicate-key
  and invalid-UTF-8 input, both digests, the message hash on a negative expiry, both laws on all three formats accepting and
  refusing → nothing written).
- `host/verbs/operator_silent_test.go` (NEW): **`TestAdmitOperatorWritesNothing`** (control; then six bar lines including an
  undefined flag, a bad bool, a missing value, `--help` and `-h` → nothing written; kills `J407-SETOUTPUT`, C9).
- Expected numstat after commit **1,394**.

### 3.1 Every design §5.1/§5.2 row S2 drives, traced (directive retro 4)

| design row | branch | test (and the no-proof / zero-action assertion) |
|---|---|---|
| §5.1 Unknown/duplicate keys, trailing JSON, invalid byte or NUL | D1–D9 | `TestOperatorStrictJSON` (each class its own error; `opDecodeRefused` asserts the zero action) |
| §5.1 Empty or over-cap request/body/task (the envelope half) | O10, O13, D1 | `TestOperatorRequestLimits` (one under, at, one over, both caps) |
| §5.1 Non-object, non-array argv, non-string element, type-confused (the envelope half) | D4, D7, O9 | `TestOperatorStrictJSON` (null, array, string, number, empty; argv a string, an int element, `null`; options an array; headless a string) |
| §5.1 Shell metacharacters inside literal task (the signing half) | O2, canonical encoding | `opAction`'s task carries `$(rm -rf ~)`, a backtick and `<b>`; the golden pins them literal (no HTML escape: mutant `CANON-HTML`) |
| §5.2 Good passkey and good enrolled SSHSIG | the whole law | `TestOperatorConfirmBothKeys` (approver, binding, passkey use) |
| §5.2 Plain y, no proof, malformed proof, bad challenge | parse, format, P6 / W6 | `TestOperatorConfirmProof` (no proof each) |
| §5.2 Forged, unknown, wrong SSH namespace, WebAuthn UV/RP/origin | P6, W2, W7, W8, W10–W12 | `TestOperatorConfirmAuthenticity` (each independently) |
| §5.2 Revoked key or replayed/unchanged nonzero counter (pure half) | P4, P5, P3w, W15 | `TestOperatorConfirmRevocationCounter` |
| §5.2 Proof from another action, episode, kind, store, argv, preset or task; a V3/A3 artifact | digest, C1, namespace, prefix | `TestOperatorConfirmBinding` |
| §5.2 Expiry equal to now, past, outside issue/lifetime bounds | P1, C2 | `TestOperatorConfirmExpiry`, `TestOperatorConfirmNegativeExpiry` |
| §5.2 Forged DB row / mismatch after claim (pure half: own bytes, current key, recorded counter) | R1, R2, P2–P5, P6 | `TestOperatorRunnerRecheckLaw` |
| judge-407 / directive 12: pure entry points write nothing | — | `TestOperatorPureEntryPointsWriteNothing`, `TestAdmitOperatorWritesNothing` |

**Not S2** (Appendix C names each owner): the store adapters and the approver-identity comparison (parked → S3); every
CAS/claim/crash/replay/race/recovery/status/prune/cancel row (S3, S4, S6); every route row (S7); every page row (S8).

---

## 4. Mutation table

Enumerated **from the diff** (`tools/mut.py`; the ids below are the ledger's): every branch of `operator_action.go` and
`operator_signed.go` — each Validate rule removed, each boundary operator moved (`>` ↔ `>=`, ±1), each cap constant ±1, each
regex at every boundary class (each anchor, `(?m)`, `(?i)`, `(?s)`, length under/over, the prefix, the hex class), each decoder
class removed (size, UTF-8, NUL, shape, duplicate at the top only, nested walk, trailing, unknown-key on and its class, the
canonical comparison, Validate), the encoder's settings (HTML escape, trailing newline), the review digest's blank id and the
signature digest's id check, **every field of the envelope dropped from the encoding** (`json:"-"`, eleven) and one reorder,
**every field of the signed message** (prefix, NUL, digest/challenge order, the challenge dropped, endianness, width, the
negative-expiry guard, the hash's error), the namespace, each law step (validate, P1, C1, C2 ×3, the lifetime constant, the
row clock at P5, the W15/R2 selector both ways, R0, R1 ×3, the format ×2, P2–P5 off and at a fixed clock, SSHSIG parse,
namespace, verify, the message's challenge, expiry and error, the hash's challenge and kind, W1–W14 off, W15 ×5, R2 ×2), the
proof's fields (count, digest, challenge, expiry, format), every `Binds`/`ApprovedBy`/`PasskeyUse` branch; the **reused** checks
the new law must reach (W2, W7, W8, W10, W11, W12, P3w, P4, P5 — each killed through the new law's named test); and S1's
`J407-SETOUTPUT`. Each is one exact replacement (text not found exactly once → `NM-text`; a mutant that does not compile →
`NM`, then rewritten to compile and re-run). Each runs against its **named** test, then the **WHOLE** package of the mutated
file (`-run '.'`: every test). A mutant whose only whole-package killer is `TestOperatorActionGolden` (a fixture comparison)
is re-run with that test skipped and the leg recorded (`w/o golden`). No S2 test is an inventory, AST or string-grep test.
Ledger: `$HOME/.ailang/state/world-iter408/plan-scratch/mutants/ledger.txt` (one line per run; run 1 on the first M4 draft —
its survivors drove the rows added to `…RevocationCounter` ({5,0}, enrolled between creation and now), `…Expiry` (lifetime 300),
`TestOperatorConfirmNegativeExpiry` and the 65-fact row of `…Validate`, and its nine `NM` were rewritten to compile — stays in the
ledger above the `=== run 2` line). The table is run 2, on the final stage tree.

**Base reproduction (C9; tag `base-whole`, a worktree at `8db22168`)**: `J407-SETOUTPUT` — named SURVIVED, whole `host/verbs` **SURVIVED**.

**Final (run 2, tag `final2`, the M4 stage tree)**:

| id | stage · file | mutation | named test | named | whole package | verdict |
|---|---|---|---|---|---|---|
| `O1-SCHEMA` | M1 · `operator_action.go` | `` if a.Schema != OperatorActionSchema { // O1 `` → `` if false && a.Schema != OperatorActionSchema { // O1 `` | `TestOperatorActionValidate` | killed (TestOperatorActionValidate) | killed (TestOperatorActionValidate,TestOperatorConfirmBothKeys) | K |
| `O2-UTF8` | M1 · `operator_action.go` | `` if !utf8.ValidString(s) \|\| strings.ContainsRune(s, 0) { `` → `` if strings.ContainsRune(s, 0) { `` | `TestOperatorActionValidate` | killed (TestOperatorActionValidate) | killed (TestOperatorActionValidate) | K |
| `O2-NUL` | M1 · `operator_action.go` | `` if !utf8.ValidString(s) \|\| strings.ContainsRune(s, 0) { `` → `` if !utf8.ValidString(s) { `` | `TestOperatorStrictJSON` | killed (TestOperatorStrictJSON) | killed (TestOperatorActionValidate,TestOperatorStrictJSON) | K |
| `O3-ID-OFF` | M1 · `operator_action.go` | `` if withID && !operatorActionIDPattern.MatchString(a.ID) { `` → `` if false && withID && !operatorActionIDPattern.MatchString(a.ID) { `` | `TestOperatorActionValidate` | killed (TestOperatorActionValidate) | killed (TestOperatorActionValidate,TestOperatorArgvIsTheBarsIntent,TestOperatorConfirmProof,TestOperatorDigestsAreDistinct) | K |
| `O3-WITHID` | M1 · `operator_action.go` | `` if withID && !operatorActionIDPattern.MatchString(a.ID) { `` → `` if !operatorActionIDPattern.MatchString(a.ID) { `` | `TestOperatorDigestsAreDistinct` | killed (TestOperatorDigestsAreDistinct) | killed (TestOperatorActionGolden,TestOperatorDigestsAreDistinct) | K |
| `O3-REVIEW-VALIDATES-ID` | M1 · `operator_action.go` | `` if err := a.validate(false); err != nil { `` → `` if err := a.validate(true); err != nil { `` | `TestOperatorDigestsAreDistinct` | killed (TestOperatorDigestsAreDistinct) | killed (TestOperatorActionGolden,TestOperatorDigestsAreDistinct) | K |
| `ID-CARET` | M1 · `operator_action.go` | `` `^oa-[0-9a-f]{32}$` `` → `` `oa-[0-9a-f]{32}$` `` | `TestOperatorActionValidate` | killed (TestOperatorActionValidate) | killed (TestOperatorActionValidate) | K |
| `ID-DOLLAR` | M1 · `operator_action.go` | `` `^oa-[0-9a-f]{32}$` `` → `` `^oa-[0-9a-f]{32}` `` | `TestOperatorActionValidate` | killed (TestOperatorActionValidate) | killed (TestOperatorActionValidate,TestOperatorArgvIsTheBarsIntent) | K |
| `ID-MULTILINE` | M1 · `operator_action.go` | `` `^oa-[0-9a-f]{32}$` `` → `` `(?m)^oa-[0-9a-f]{32}$` `` | `TestOperatorActionValidate` | killed (TestOperatorActionValidate) | killed (TestOperatorActionValidate,TestOperatorArgvIsTheBarsIntent) | K |
| `ID-CASE` | M1 · `operator_action.go` | `` `^oa-[0-9a-f]{32}$` `` → `` `(?i)^oa-[0-9a-f]{32}$` `` | `TestOperatorActionValidate` | killed (TestOperatorActionValidate) | killed (TestOperatorActionValidate,TestOperatorArgvIsTheBarsIntent) | K |
| `ID-DOTALL` | M1 · `operator_action.go` | `` `^oa-[0-9a-f]{32}$` `` → `` `(?s)^oa-[0-9a-f]{32}$` `` | `TestOperatorActionValidate` | survived | survived — **equivalent**: `(?s)` changes only `.`, and the pattern has none | E |
| `ID-LEN-UNDER` | M1 · `operator_action.go` | `` `^oa-[0-9a-f]{32}$` `` → `` `^oa-[0-9a-f]{31,32}$` `` | `TestOperatorActionValidate` | killed (TestOperatorActionValidate) | killed (TestOperatorActionValidate,TestOperatorArgvIsTheBarsIntent) | K |
| `ID-LEN-OVER` | M1 · `operator_action.go` | `` `^oa-[0-9a-f]{32}$` `` → `` `^oa-[0-9a-f]{32,33}$` `` | `TestOperatorActionValidate` | killed (TestOperatorActionValidate) | killed (TestOperatorActionValidate,TestOperatorArgvIsTheBarsIntent) | K |
| `ID-PREFIX` | M1 · `operator_action.go` | `` `^oa-[0-9a-f]{32}$` `` → `` `^[a-z][a-z]-[0-9a-f]{32}$` `` | `TestOperatorActionValidate` | killed (TestOperatorActionValidate) | killed (TestOperatorActionValidate) | K |
| `ID-UPPERHEX` | M1 · `operator_action.go` | `` `^oa-[0-9a-f]{32}$` `` → `` `^oa-[0-9a-fA-F]{32}$` `` | `TestOperatorArgvIsTheBarsIntent` | killed (TestOperatorArgvIsTheBarsIntent) | killed (TestOperatorArgvIsTheBarsIntent) | K |
| `O4-KIND-ANY` | M1 · `operator_action.go` | `` if a.Kind != "start" && a.Kind != "finish" { // O4 `` → `` if a.Kind == "" { // O4 `` | `TestOperatorActionValidate` | killed (TestOperatorActionValidate) | killed (TestOperatorActionValidate) | K |
| `O4-KIND-FOLD` | M1 · `operator_action.go` | `` if a.Kind != "start" && a.Kind != "finish" { // O4 `` → `` if !strings.EqualFold(a.Kind, "start") && !strings.EqualFold(a.Kind, … `` | `TestOperatorActionValidate` | killed (TestOperatorActionValidate) | killed (TestOperatorActionValidate) | K |
| `O5-EPISODE` | M1 · `operator_action.go` | `` if !episodePattern.MatchString(a.Episode) { // O5 `` → `` if false && !episodePattern.MatchString(a.Episode) { // O5 `` | `TestOperatorActionValidate` | killed (TestOperatorActionValidate) | killed (TestOperatorActionValidate) | K |
| `O6-ABS` | M1 · `operator_action.go` | `` if !path.IsAbs(a.Store) \|\| path.Clean(a.Store) != a.Store { // O6 `` → `` if path.Clean(a.Store) != a.Store { // O6 `` | `TestOperatorActionValidate` | killed (TestOperatorActionValidate) | killed (TestOperatorActionValidate,TestOperatorRunnerRecheckLaw) | K |
| `O6-CLEAN` | M1 · `operator_action.go` | `` if !path.IsAbs(a.Store) \|\| path.Clean(a.Store) != a.Store { // O6 `` → `` if !path.IsAbs(a.Store) { // O6 `` | `TestOperatorActionValidate` | killed (TestOperatorActionValidate) | killed (TestOperatorActionValidate) | K |
| `O7-SETOOLS` | M1 · `operator_action.go` | `` p == "se-tools" \|\| len(p) > operatorPresetMax `` → `` len(p) > operatorPresetMax `` | `TestOperatorActionValidate` | killed (TestOperatorActionValidate) | killed (TestOperatorActionValidate) | K |
| `O7-LEN-OVER` | M1 · `operator_action.go` | `` len(p) > operatorPresetMax \|\| `` → `` len(p) > operatorPresetMax+1 \|\| `` | `TestOperatorActionValidate` | killed (TestOperatorActionValidate) | killed (TestOperatorActionValidate) | K |
| `O7-LEN-AT` | M1 · `operator_action.go` | `` len(p) > operatorPresetMax \|\| `` → `` len(p) >= operatorPresetMax \|\| `` | `TestOperatorActionValidate` | killed (TestOperatorActionValidate) | killed (TestOperatorActionValidate) | K |
| `O7-GRAMMAR` | M1 · `operator_action.go` | `` \|\| !episodePattern.MatchString(p) { // O7 `` → `` { // O7 `` | `TestOperatorActionValidate` | killed (TestOperatorActionValidate) | killed (TestOperatorActionValidate) | K |
| `O8-START-READY` | M1 · `operator_action.go` | `` (a.Options.ConfirmOnSurface \|\| a.Options.Ready) { // O8 `` → `` (a.Options.ConfirmOnSurface) { // O8 `` | `TestOperatorActionValidate` | killed (TestOperatorActionValidate) | killed (TestOperatorActionValidate) | K |
| `O8-START-COS` | M1 · `operator_action.go` | `` (a.Options.ConfirmOnSurface \|\| a.Options.Ready) { // O8 `` → `` (a.Options.Ready) { // O8 `` | `TestOperatorActionValidate` | killed (TestOperatorActionValidate) | killed (TestOperatorActionValidate) | K |
| `O8-FINISH-HEADLESS` | M1 · `operator_action.go` | `` if a.Kind == "finish" && a.Options.Headless { `` → `` if false && a.Kind == "finish" && a.Options.Headless { `` | `TestOperatorActionValidate` | killed (TestOperatorActionValidate) | killed (TestOperatorActionValidate) | K |
| `O9-ARGV-LEN` | M1 · `operator_action.go` | `` !slices.Equal(a.Argv, want) { // O9 `` → `` len(a.Argv) != len(want) && slices.Equal(want, want) { // O9 `` | `TestOperatorActionValidate` | killed (TestOperatorActionValidate) | killed (TestOperatorActionValidate) | K |
| `O9-ARGV-OFF` | M1 · `operator_action.go` | `` !slices.Equal(a.Argv, want) { // O9 `` → `` false && !slices.Equal(a.Argv, want) { // O9 `` | `TestOperatorActionValidate` | killed (TestOperatorActionValidate) | killed (TestOperatorActionValidate,TestOperatorStrictJSON) | K |
| `ARGV-LAUNCH` | M1 · `operator_action.go` | `` argv = append(argv, "--launch") `` → `` argv = argv[:len(argv):len(argv)] `` | `TestOperatorArgvIsTheBarsIntent` | killed (TestOperatorArgvIsTheBarsIntent) | killed (TestOperatorActionGolden,TestOperatorActionValidate,TestOperatorArgvIsTheBarsIntent,TestOperatorStrictJSON) | K |
| `ARGV-HEADLESS` | M1 · `operator_action.go` | `` argv = append(argv, "--headless") `` → `` argv = argv[:len(argv):len(argv)] `` | `TestOperatorArgvIsTheBarsIntent` | killed (TestOperatorArgvIsTheBarsIntent) | killed (TestOperatorActionGolden,TestOperatorActionValidate,TestOperatorArgvIsTheBarsIntent,TestOperatorStrictJSON) | K |
| `ARGV-COS` | M1 · `operator_action.go` | `` argv = append(argv, "--confirm-on-surface") `` → `` argv = argv[:len(argv):len(argv)] `` | `TestOperatorArgvIsTheBarsIntent` | killed (TestOperatorArgvIsTheBarsIntent) | killed (TestOperatorActionGolden,TestOperatorActionValidate,TestOperatorArgvIsTheBarsIntent) | K |
| `ARGV-READY` | M1 · `operator_action.go` | `` argv = append(argv, "--ready") `` → `` argv = argv[:len(argv):len(argv)] `` | `TestOperatorArgvIsTheBarsIntent` | killed (TestOperatorArgvIsTheBarsIntent) | killed (TestOperatorArgvIsTheBarsIntent) | K |
| `ARGV-ORDER` | M1 · `operator_action.go` | `` argv := []string{"ailang-worldd", "episode", kind, episode, "--preset… `` → `` argv := []string{"ailang-worldd", "episode", kind, "--preset", o.Pres… `` | `TestOperatorArgvIsTheBarsIntent` | killed (TestOperatorArgvIsTheBarsIntent) | killed (TestOperatorActionGolden,TestOperatorArgvIsTheBarsIntent,TestOperatorStrictJSON) | K |
| `O10-TASK-EMPTY` | M1 · `operator_action.go` | `` (len(a.Task) < 1 \|\| len(a.Task) > MaxOperatorTaskBytes) `` → `` (len(a.Task) > MaxOperatorTaskBytes) `` | `TestOperatorRequestLimits` | killed (TestOperatorRequestLimits) | killed (TestOperatorRequestLimits) | K |
| `O10-TASK-OVER` | M1 · `operator_action.go` | `` len(a.Task) > MaxOperatorTaskBytes) { // O10 `` → `` len(a.Task) > MaxOperatorTaskBytes+1) { // O10 `` | `TestOperatorRequestLimits` | killed (TestOperatorRequestLimits) | killed (TestOperatorRequestLimits) | K |
| `O10-TASK-AT` | M1 · `operator_action.go` | `` len(a.Task) > MaxOperatorTaskBytes) { // O10 `` → `` len(a.Task) >= MaxOperatorTaskBytes) { // O10 `` | `TestOperatorRequestLimits` | killed (TestOperatorRequestLimits) | killed (TestOperatorRequestLimits) | K |
| `O10-TASK-RUNES` | M1 · `operator_action.go` | `` len(a.Task) > MaxOperatorTaskBytes) { // O10 `` → `` utf8.RuneCountInString(a.Task) > MaxOperatorTaskBytes) { // O10 `` | `TestOperatorRequestLimits` | killed (TestOperatorRequestLimits) | killed (TestOperatorRequestLimits) | K |
| `O10-FINISH-TASK` | M1 · `operator_action.go` | `` if a.Kind == "finish" && a.Task != "" { `` → `` if false && a.Kind == "finish" && a.Task != "" { `` | `TestOperatorRequestLimits` | killed (TestOperatorRequestLimits) | killed (TestOperatorRequestLimits) | K |
| `TASK-CAP` | M1 · `operator_action.go` | `` MaxOperatorTaskBytes   = 8 << 10 `` → `` MaxOperatorTaskBytes   = 8<<10 + 1 `` | `TestOperatorRequestLimits` | killed (TestOperatorRequestLimits) | killed (TestOperatorRequestLimits) | K |
| `O11-RES-MIN` | M1 · `operator_action.go` | `` if len(a.Resolution) < 1 \|\| len(a.Resolution) > operatorResolvedMax {… `` → `` if len(a.Resolution) > operatorResolvedMax { // O11 `` | `TestOperatorActionValidate` | killed (TestOperatorActionValidate) | killed (TestOperatorActionValidate) | K |
| `O11-RES-MAX` | M1 · `operator_action.go` | `` if len(a.Resolution) < 1 \|\| len(a.Resolution) > operatorResolvedMax {… `` → `` if len(a.Resolution) < 1 { // O11 `` | `TestOperatorActionValidate` | killed (TestOperatorActionValidate) | killed (TestOperatorActionValidate) | K |
| `O11-RES-MAX-AT` | M1 · `operator_action.go` | `` len(a.Resolution) > operatorResolvedMax { // O11 `` → `` len(a.Resolution) >= operatorResolvedMax { // O11 `` | `TestOperatorActionValidate` | killed (TestOperatorActionValidate) | killed (TestOperatorActionValidate) | K |
| `O11-NAME` | M1 · `operator_action.go` | `` if !operatorResolvedName.MatchString(r.Name) { `` → `` if false && !operatorResolvedName.MatchString(r.Name) { `` | `TestOperatorActionValidate` | killed (TestOperatorActionValidate) | killed (TestOperatorActionValidate) | K |
| `O11-ORDER-EQ` | M1 · `operator_action.go` | `` a.Resolution[i-1].Name >= r.Name `` → `` a.Resolution[i-1].Name > r.Name `` | `TestOperatorActionValidate` | killed (TestOperatorActionValidate) | killed (TestOperatorActionValidate) | K |
| `O11-ORDER-OFF` | M1 · `operator_action.go` | `` if i > 0 && a.Resolution[i-1].Name >= r.Name { `` → `` if false && i > 0 && a.Resolution[i-1].Name >= r.Name { `` | `TestOperatorActionValidate` | killed (TestOperatorActionValidate) | killed (TestOperatorActionValidate) | K |
| `RES-NAME-CASE` | M1 · `operator_action.go` | `` `^[a-z][a-z0-9_]{0,63}$` `` → `` `(?i)^[a-z][a-z0-9_]{0,63}$` `` | `TestOperatorActionValidate` | killed (TestOperatorActionValidate) | killed (TestOperatorActionValidate) | K |
| `O12-REVIEW` | M1 · `operator_action.go` | `` if a.Review == "" { // O12 `` → `` if false { // O12 `` | `TestOperatorActionValidate` | killed (TestOperatorActionValidate) | killed (TestOperatorActionValidate) | K |
| `O13-SIZE-OFF` | M1 · `operator_action.go` | `` if len(c) > MaxOperatorActionBytes { // O13 `` → `` if false { // O13 `` | `TestOperatorRequestLimits` | killed (TestOperatorRequestLimits) | killed (TestOperatorRequestLimits) | K |
| `O13-SIZE-AT` | M1 · `operator_action.go` | `` if len(c) > MaxOperatorActionBytes { // O13 `` → `` if len(c) >= MaxOperatorActionBytes { // O13 `` | `TestOperatorRequestLimits` | killed (TestOperatorRequestLimits) | killed (TestOperatorRequestLimits) | K |
| `O13-SIZE-OVER` | M1 · `operator_action.go` | `` if len(c) > MaxOperatorActionBytes { // O13 `` → `` if len(c) > MaxOperatorActionBytes+1 { // O13 `` | `TestOperatorRequestLimits` | killed (TestOperatorRequestLimits) | killed (TestOperatorRequestLimits) | K |
| `ACTION-CAP` | M1 · `operator_action.go` | `` MaxOperatorActionBytes = 128 << 10 `` → `` MaxOperatorActionBytes = 128<<10 + 1 `` | `TestOperatorRequestLimits` | killed (TestOperatorRequestLimits) | killed (TestOperatorRequestLimits) | K |
| `CANON-HTML` | M1 · `operator_action.go` | `` enc.SetEscapeHTML(false)⏎if err := enc.Encode(a) `` → `` enc.SetEscapeHTML(true)⏎if err := enc.Encode(a) `` | `TestOperatorStrictJSON` | killed (TestOperatorStrictJSON) | killed (TestOperatorActionGolden,TestOperatorStrictJSON) | K |
| `CANON-NEWLINE` | M1 · `operator_action.go` | `` return bytes.TrimSuffix(buf.Bytes(), []byte("\⏎")), nil⏎}⏎⏎// DecodeO… `` → `` return buf.Bytes(), nil⏎}⏎⏎// DecodeOperatorAction `` | `TestOperatorStrictJSON` | killed (TestOperatorStrictJSON) | killed (TestOperatorActionGolden,TestOperatorStrictJSON) | K |
| `D1-SIZE-OFF` | M1 · `operator_action.go` | `` if len(b) > MaxOperatorActionBytes { // D1 `` → `` if false { // D1 `` | `TestOperatorRequestLimits` | killed (TestOperatorRequestLimits) | killed (TestOperatorRequestLimits) | K |
| `D1-SIZE-AT` | M1 · `operator_action.go` | `` if len(b) > MaxOperatorActionBytes { // D1 `` → `` if len(b) >= MaxOperatorActionBytes { // D1 `` | `TestOperatorRequestLimits` | killed (TestOperatorRequestLimits) | killed (TestOperatorRequestLimits) | K |
| `D1-SIZE-OVER` | M1 · `operator_action.go` | `` if len(b) > MaxOperatorActionBytes { // D1 `` → `` if len(b) > MaxOperatorActionBytes+1 { // D1 `` | `TestOperatorRequestLimits` | killed (TestOperatorRequestLimits) | killed (TestOperatorRequestLimits) | K |
| `D2-UTF8` | M1 · `operator_action.go` | `` if !utf8.Valid(b) { // D2 `` → `` if false { // D2 `` | `TestOperatorStrictJSON` | killed (TestOperatorStrictJSON) | killed (TestOperatorStrictJSON) | K |
| `D3-NUL` | M1 · `operator_action.go` | `` if bytes.IndexByte(b, 0) >= 0 { // D3 `` → `` if false { // D3 `` | `TestOperatorStrictJSON` | killed (TestOperatorStrictJSON) | killed (TestOperatorStrictJSON) | K |
| `D4-SHAPE` | M1 · `operator_action.go` | `` tok != json.Delim('{') { // D4 `` → `` false { // D4 `` | `TestOperatorStrictJSON` | killed (TestOperatorStrictJSON) | killed (TestOperatorStrictJSON) | K |
| `D4-SHAPE-ERR` | M1 · `operator_action.go` | `` if tok, err := dec.Token(); err != nil \|\| tok != json.Delim('{') { //… `` → `` if tok, _ := dec.Token(); tok != json.Delim('{') { // D4 `` | `TestOperatorStrictJSON` | survived | survived — **equivalent**: `dec.Token()` returns a nil token whenever it returns an error, so `tok != Delim('{')` already refuses | E |
| `D5-DUP` | M1 · `operator_action.go` | `` if seen[k] { // D5 `` → `` if false { // D5 `` | `TestOperatorStrictJSON` | killed (TestOperatorStrictJSON) | killed (TestOperatorStrictJSON) | K |
| `D5-DUP-TOPONLY` | M1 · `operator_action.go` | `` seen := map[string]bool{}⏎for dec.More() { `` → `` seen := topSeen⏎for dec.More() { `` | `TestOperatorStrictJSON` | killed (TestOperatorStrictJSON) | killed (TestOperatorActionGolden,TestOperatorRequestLimits,TestOperatorStrictJSON) | K |
| `D5-WALK-NESTED` | M1 · `operator_action.go` | `` if err := walk(next); err != nil { `` → `` _ = next⏎if err := walk(nil); err != nil { `` | `TestOperatorStrictJSON` | killed (TestOperatorStrictJSON) | killed (TestOperatorActionGolden,TestOperatorRequestLimits,TestOperatorStrictJSON) | K |
| `D6-TRAILING` | M1 · `operator_action.go` | `` if _, err := dec.Token(); err != io.EOF { // D6 `` → `` if _, err := dec.Token(); false && err != io.EOF { // D6 `` | `TestOperatorStrictJSON` | killed (TestOperatorStrictJSON) | killed (TestOperatorStrictJSON) | K |
| `D7-UNKNOWN-ON` | M1 · `operator_action.go` | `` dec.DisallowUnknownFields()⏎var a OperatorAction `` → `` var a OperatorAction `` | `TestOperatorStrictJSON` | killed (TestOperatorStrictJSON) | killed (TestOperatorStrictJSON) | K |
| `D7-UNKNOWN-CLASS` | M1 · `operator_action.go` | `` if strings.HasPrefix(err.Error(), "json: unknown field ") { `` → `` if false { `` | `TestOperatorStrictJSON` | killed (TestOperatorStrictJSON) | killed (TestOperatorStrictJSON) | K |
| `D8-CANON` | M1 · `operator_action.go` | `` if !bytes.Equal(c, b) { // D8 `` → `` if false && !bytes.Equal(c, b) { // D8 `` | `TestOperatorStrictJSON` | killed (TestOperatorStrictJSON) | killed (TestOperatorRequestLimits,TestOperatorStrictJSON) | K |
| `D9-VALIDATE` | M1 · `operator_action.go` | `` if err := a.Validate(); err != nil { // D9 `` → `` if false { // D9 `` | `TestOperatorStrictJSON` | killed (TestOperatorStrictJSON) | killed (TestOperatorStrictJSON) | K |
| `RD-BLANK-ID` | M1 · `operator_action.go` | `` a.ID = ""⏎if err := a.validate(false) `` → `` if err := a.validate(false) `` | `TestOperatorDigestsAreDistinct` | killed (TestOperatorDigestsAreDistinct) | killed (TestOperatorActionGolden,TestOperatorDigestsAreDistinct) | K |
| `SD-VALIDATE` | M1 · `operator_action.go` | `` func OperatorSignatureDigest(a OperatorAction) ([32]byte, error) {⏎if… `` → `` func OperatorSignatureDigest(a OperatorAction) ([32]byte, error) {⏎if… `` | `TestOperatorDigestsAreDistinct` | killed (TestOperatorDigestsAreDistinct) | killed (TestOperatorActionValidate,TestOperatorConfirmProof,TestOperatorDigestsAreDistinct) | K |
| `MSG-NEG` | M1 · `operator_action.go` | `` if expiry < 0 {⏎return nil, ErrNegativeExpiry⏎}⏎m := make([]byte, 0, … `` → `` m := make([]byte, 0, OperatorActionMessageLen) `` | `TestOperatorActionMessage` | killed (TestOperatorActionMessage) | killed (TestOperatorActionMessage,TestOperatorConfirmNegativeExpiry) | K |
| `MSG-PREFIX` | M1 · `operator_action.go` | `` m = append(m, OperatorActionPrefix...) `` → `` m = append(m, VerbPrefix...) `` | `TestOperatorActionMessage` | killed (TestOperatorActionMessage) | killed (TestOperatorActionGolden,TestOperatorActionMessage,TestOperatorConfirmBinding) | K |
| `MSG-NUL` | M1 · `operator_action.go` | `` m = append(m, OperatorActionPrefix...)⏎m = append(m, 0) `` → `` m = append(m, OperatorActionPrefix...) `` | `TestOperatorActionMessage` | killed (TestOperatorActionMessage) | killed (TestOperatorActionGolden,TestOperatorActionMessage) | K |
| `MSG-ORDER` | M1 · `operator_action.go` | `` m = append(m, digest[:]...)⏎m = append(m, challenge[:]...) `` → `` m = append(m, challenge[:]...)⏎m = append(m, digest[:]...) `` | `TestOperatorActionMessage` | killed (TestOperatorActionMessage) | killed (TestOperatorActionGolden,TestOperatorActionMessage) | K |
| `MSG-DROP-CHALLENGE` | M1 · `operator_action.go` | `` m = append(m, challenge[:]...)⏎m = binary.BigEndian.AppendUint64(m, u… `` → `` m = binary.BigEndian.AppendUint64(m, uint64(expiry))⏎return m, nil⏎}⏎… `` | `TestOperatorActionMessage` | killed (TestOperatorActionMessage) | killed (TestOperatorActionGolden,TestOperatorActionMessage,TestOperatorConfirmProof,TestOperatorRunnerRecheckLaw) | K |
| `MSG-ENDIAN` | M1 · `operator_action.go` | `` m = binary.BigEndian.AppendUint64(m, uint64(expiry))⏎return m, nil⏎}⏎… `` → `` m = binary.LittleEndian.AppendUint64(m, uint64(expiry))⏎return m, nil… `` | `TestOperatorActionMessage` | killed (TestOperatorActionMessage) | killed (TestOperatorActionGolden,TestOperatorActionMessage) | K |
| `MSG-WIDTH` | M1 · `operator_action.go` | `` m = binary.BigEndian.AppendUint64(m, uint64(expiry))⏎return m, nil⏎}⏎… `` → `` m = binary.BigEndian.AppendUint32(m, uint32(expiry))⏎return m, nil⏎}⏎… `` | `TestOperatorActionMessage` | killed (TestOperatorActionMessage) | killed (TestOperatorActionGolden,TestOperatorActionMessage) | K |
| `MSGHASH-ERR` | M1 · `operator_action.go` | `` m, err := OperatorActionMessage(digest, challenge, expiry)⏎if err != … `` → `` m, _ := OperatorActionMessage(digest, challenge, expiry) `` | `TestOperatorActionMessage` | killed (TestOperatorActionMessage) | killed (TestOperatorActionMessage) | K |
| `NAMESPACE` | M1 · `operator_action.go` | `` OperatorActionNamespace  = "world-operator-action@sunholo" `` → `` OperatorActionNamespace  = "world-verb@sunholo" `` | `TestOperatorConfirmAuthenticity` | killed (TestOperatorConfirmAuthenticity) | killed (TestOperatorConfirmAuthenticity,TestOperatorConfirmBinding) | K |
| `FIELD-DROP-store` | M1 · `operator_action.go` | `` Store      string             `json:"store"` `` → `` Store      string             `json:"-"` `` | `TestOperatorConfirmBinding` | killed (TestOperatorConfirmBinding) | killed (TestOperatorActionGolden,TestOperatorConfirmBinding,TestOperatorDigestsAreDistinct,TestOperatorRequestLimits,TestOperatorStrictJSON) | K |
| `FIELD-DROP-task` | M1 · `operator_action.go` | `` Task       string             `json:"task"` `` → `` Task       string             `json:"-"` `` | `TestOperatorConfirmBinding` | killed (TestOperatorConfirmBinding) | killed (TestOperatorActionGolden,TestOperatorConfirmBinding,TestOperatorConfirmBothKeys,TestOperatorDigestsAreDistinct,TestOperatorRequestLimits,TestOperatorRunnerRecheckLaw,TestOperatorStrictJSON) | K |
| `FIELD-DROP-review` | M1 · `operator_action.go` | `` Review     string             `json:"review"` `` → `` Review     string             `json:"-"` `` | `TestOperatorConfirmBinding` | killed (TestOperatorConfirmBinding) | killed (TestOperatorActionGolden,TestOperatorConfirmBinding,TestOperatorDigestsAreDistinct,TestOperatorRequestLimits,TestOperatorStrictJSON) | K |
| `FIELD-DROP-resolution` | M1 · `operator_action.go` | `` Resolution []OperatorResolved `json:"resolution"` `` → `` Resolution []OperatorResolved `json:"-"` `` | `TestOperatorConfirmBinding` | killed (TestOperatorConfirmBinding) | killed (TestOperatorActionGolden,TestOperatorConfirmBinding,TestOperatorDigestsAreDistinct,TestOperatorRequestLimits,TestOperatorStrictJSON) | K |
| `FIELD-DROP-episode` | M1 · `operator_action.go` | `` Episode    string             `json:"episode"` `` → `` Episode    string             `json:"-"` `` | `TestOperatorActionGolden` | killed (TestOperatorActionGolden) | killed (TestOperatorActionGolden,TestOperatorRequestLimits,TestOperatorStrictJSON) | K |
| `FIELD-DROP-id` | M1 · `operator_action.go` | `` ID         string             `json:"id"` `` → `` ID         string             `json:"-"` `` | `TestOperatorConfirmBinding` | killed (TestOperatorConfirmBinding) | killed (TestOperatorActionGolden,TestOperatorConfirmBinding,TestOperatorDigestsAreDistinct,TestOperatorRequestLimits,TestOperatorStrictJSON) | K |
| `FIELD-DROP-kind` | M1 · `operator_action.go` | `` Kind       string             `json:"kind"` `` → `` Kind       string             `json:"-"` `` | `TestOperatorActionGolden` | killed (TestOperatorActionGolden) | killed (TestOperatorActionGolden,TestOperatorRequestLimits,TestOperatorStrictJSON) | K |
| `FIELD-DROP-argv` | M1 · `operator_action.go` | `` Argv       []string           `json:"argv"` `` → `` Argv       []string           `json:"-"` `` | `TestOperatorActionGolden` | killed (TestOperatorActionGolden) | killed (TestOperatorActionGolden,TestOperatorRequestLimits,TestOperatorStrictJSON) | K |
| `FIELD-DROP-options` | M1 · `operator_action.go` | `` Options    OperatorOptions    `json:"options"` `` → `` Options    OperatorOptions    `json:"-"` `` | `TestOperatorActionGolden` | killed (TestOperatorActionGolden) | killed (TestOperatorActionGolden,TestOperatorRequestLimits,TestOperatorStrictJSON) | K |
| `FIELD-DROP-preset` | M1 · `operator_action.go` | `` Preset           string `json:"preset"` `` → `` Preset           string `json:"-"` `` | `TestOperatorActionGolden` | killed (TestOperatorActionGolden) | killed (TestOperatorActionGolden,TestOperatorRequestLimits,TestOperatorStrictJSON) | K |
| `FIELD-DROP-headless` | M1 · `operator_action.go` | `` Headless         bool   `json:"headless"` `` → `` Headless         bool   `json:"-"` `` | `TestOperatorActionGolden` | killed (TestOperatorActionGolden) | killed (TestOperatorActionGolden,TestOperatorRequestLimits,TestOperatorStrictJSON) | K |
| `FIELD-DROP-value` | M1 · `operator_action.go` | `` Value string `json:"value"` `` → `` Value string `json:"-"` `` | `TestOperatorActionGolden` | killed (TestOperatorActionGolden) | killed (TestOperatorActionGolden,TestOperatorStrictJSON) | K |
| `FIELD-REORDER` | M1 · `operator_action.go` | `` Kind       string             `json:"kind"`       // "start" or "fini… `` → `` Episode    string             `json:"episode"`    // ^[a-z0-9][a-z0-9… `` | `TestOperatorActionGolden` | killed (TestOperatorActionGolden) | killed (TestOperatorActionGolden,TestOperatorStrictJSON) | K |
| `C0-VALIDATE` | M3 · `operator_signed.go` | `` func ConfirmOperatorAction(in OperatorConfirmInput) (*SignedOperatorA… `` → `` func ConfirmOperatorAction(in OperatorConfirmInput) (*SignedOperatorA… `` | `TestOperatorConfirmProof` | killed (TestOperatorConfirmProof) | killed (TestOperatorConfirmProof) | K |
| `P1-ADMIT` | M3 · `operator_signed.go` | `` if v := Admit(in.Admit); v != Admitted { // P1 `` → `` if v := Admit(in.Admit); false && v != Admitted { // P1 `` | `TestOperatorConfirmExpiry` | killed (TestOperatorConfirmExpiry) | killed (TestOperatorConfirmExpiry) | K |
| `C1-ROW-DIGEST` | M3 · `operator_signed.go` | `` if in.Admit.StoredDigest != d { // C1 `` → `` if false { // C1 `` | `TestOperatorConfirmBinding` | killed (TestOperatorConfirmBinding) | killed (TestOperatorConfirmBinding) | K |
| `C2-LIFETIME-OFF` | M3 · `operator_signed.go` | `` if in.Admit.ExpiresAt-in.Admit.CreatedAt > int64(OperatorActionLifeti… `` → `` if false && in.Admit.ExpiresAt-in.Admit.CreatedAt > int64(OperatorAct… `` | `TestOperatorConfirmExpiry` | killed (TestOperatorConfirmExpiry) | killed (TestOperatorConfirmExpiry) | K |
| `C2-LIFETIME-AT` | M3 · `operator_signed.go` | `` in.Admit.ExpiresAt-in.Admit.CreatedAt > int64( `` → `` in.Admit.ExpiresAt-in.Admit.CreatedAt >= int64( `` | `TestOperatorConfirmExpiry` | killed (TestOperatorConfirmExpiry) | killed (TestOperatorConfirmAuthenticity,TestOperatorConfirmBinding,TestOperatorConfirmBothKeys,TestOperatorConfirmExpiry,TestOperatorConfirmProof,TestOperatorConfirmRevocationCounter) | K |
| `C2-LIFETIME-OVER` | M3 · `operator_signed.go` | `` int64(OperatorActionLifetime/time.Second) { // C2 `` → `` int64(OperatorActionLifetime/time.Second)+1 { // C2 `` | `TestOperatorConfirmExpiry` | killed (TestOperatorConfirmExpiry) | killed (TestOperatorConfirmExpiry) | K |
| `LIFETIME-CONST` | M1 · `operator_action.go` | `` OperatorActionLifetime   = 10 * time.Minute `` → `` OperatorActionLifetime   = 11 * time.Minute `` | `TestOperatorConfirmExpiry` | killed (TestOperatorConfirmExpiry) | killed (TestOperatorConfirmExpiry) | K |
| `C-EXPIRY-ROW` | M3 · `operator_signed.go` | `` return operatorProof(d, in.Challenge, in.Admit.ExpiresAt, in.Admit.No… `` → `` return operatorProof(d, in.Challenge, in.Admit.ExpiresAt+0*in.Admit.N… `` | `TestOperatorConfirmRevocationCounter` | killed (TestOperatorConfirmRevocationCounter) | killed (TestOperatorConfirmRevocationCounter) | K |
| `C-W15-SELECT` | M3 · `operator_signed.go` | `` in.Artifact, in.Format, in.Approver, true) `` → `` in.Artifact, in.Format, in.Approver, false) `` | `TestOperatorConfirmRevocationCounter` | killed (TestOperatorConfirmRevocationCounter) | killed (TestOperatorConfirmBothKeys,TestOperatorConfirmExpiry,TestOperatorConfirmRevocationCounter) | K |
| `R0-VALIDATE` | M3 · `operator_signed.go` | `` func CheckOperatorAction(in OperatorCheckInput) (*SignedOperatorActio… `` → `` func CheckOperatorAction(in OperatorCheckInput) (*SignedOperatorActio… `` | `TestOperatorRunnerRecheckLaw` | killed (TestOperatorRunnerRecheckLaw) | killed (TestOperatorRunnerRecheckLaw) | K |
| `R1-EXPIRY-OFF` | M3 · `operator_signed.go` | `` if in.Now >= in.Expiry { // R1 `` → `` if false { // R1 `` | `TestOperatorRunnerRecheckLaw` | killed (TestOperatorRunnerRecheckLaw) | killed (TestOperatorRunnerRecheckLaw) | K |
| `R1-EXPIRY-STRICT` | M3 · `operator_signed.go` | `` if in.Now >= in.Expiry { // R1 `` → `` if in.Now > in.Expiry { // R1 `` | `TestOperatorRunnerRecheckLaw` | killed (TestOperatorRunnerRecheckLaw) | killed (TestOperatorRunnerRecheckLaw) | K |
| `R1-EXPIRY-EARLY` | M3 · `operator_signed.go` | `` if in.Now >= in.Expiry { // R1 `` → `` if in.Now >= in.Expiry-1 { // R1 `` | `TestOperatorRunnerRecheckLaw` | killed (TestOperatorRunnerRecheckLaw) | killed (TestOperatorRunnerRecheckLaw) | K |
| `R-W15-SELECT` | M3 · `operator_signed.go` | `` in.Artifact, in.Format, in.Approver, false) `` → `` in.Artifact, in.Format, in.Approver, true) `` | `TestOperatorRunnerRecheckLaw` | killed (TestOperatorRunnerRecheckLaw) | killed (TestOperatorRunnerRecheckLaw) | K |
| `F-FORMAT-OFF` | M3 · `operator_signed.go` | `` if format != FormatSSH && format != FormatWebAuthn { `` → `` if format == "" { `` | `TestOperatorConfirmProof` | killed (TestOperatorConfirmProof) | killed (TestOperatorConfirmProof) | K |
| `F-FORMAT-NONE` | M3 · `operator_signed.go` | `` if format != FormatSSH && format != FormatWebAuthn { `` → `` if false { `` | `TestOperatorConfirmProof` | killed (TestOperatorConfirmProof) | killed (TestOperatorConfirmProof,TestOperatorRunnerRecheckLaw) | K |
| `P2-P5-OFF` | M3 · `operator_signed.go` | `` if err := checkApprover(AdmitInput{Now: now}, ap, format); err != nil… `` → `` if err := checkApprover(AdmitInput{Now: now}, ap, format); false && e… `` | `TestOperatorConfirmRevocationCounter` | killed (TestOperatorConfirmRevocationCounter) | killed (TestOperatorConfirmRevocationCounter,TestOperatorRunnerRecheckLaw) | K |
| `P2-P5-NOW` | M3 · `operator_signed.go` | `` checkApprover(AdmitInput{Now: now}, ap, format) `` → `` checkApprover(AdmitInput{Now: 1 << 40}, ap, format) `` | `TestOperatorConfirmRevocationCounter` | killed (TestOperatorConfirmRevocationCounter) | killed (TestOperatorConfirmRevocationCounter) | K |
| `SSH-PARSE` | M3 · `operator_signed.go` | `` sig, err := ParseSSHSIG(artifact)⏎if err != nil {⏎return nil, err⏎} `` → `` sig, _ := ParseSSHSIG(artifact)⏎if sig == nil {⏎return nil, ErrBadSig… `` | `TestOperatorConfirmProof` | killed (TestOperatorConfirmProof) | killed (TestOperatorConfirmProof) | K |
| `SSH-NAMESPACE` | M3 · `operator_signed.go` | `` VerifySSHSIG(sig, OperatorActionNamespace, msg, ap.Key) `` → `` VerifySSHSIG(sig, sig.Namespace, msg, ap.Key) `` | `TestOperatorConfirmAuthenticity` | killed (TestOperatorConfirmAuthenticity) | killed (TestOperatorConfirmAuthenticity,TestOperatorConfirmBinding) | K |
| `SSH-VERIFY-OFF` | M3 · `operator_signed.go` | `` if err := VerifySSHSIG(sig, OperatorActionNamespace, msg, ap.Key); er… `` → `` if err := VerifySSHSIG(sig, OperatorActionNamespace, msg, ap.Key); fa… `` | `TestOperatorConfirmAuthenticity` | killed (TestOperatorConfirmAuthenticity) | killed (TestOperatorConfirmAuthenticity,TestOperatorConfirmBinding,TestOperatorConfirmProof,TestOperatorRunnerRecheckLaw) | K |
| `SSH-MSG-CHALLENGE` | M3 · `operator_signed.go` | `` msg, err := OperatorActionMessage(d, challenge, expiry) `` → `` msg, err := OperatorActionMessage(d, [32]byte{}, expiry) `` | `TestOperatorConfirmBothKeys` | killed (TestOperatorConfirmBothKeys) | killed (TestOperatorConfirmBothKeys,TestOperatorConfirmExpiry,TestOperatorConfirmRevocationCounter,TestOperatorRunnerRecheckLaw) | K |
| `SSH-MSG-EXPIRY` | M3 · `operator_signed.go` | `` msg, err := OperatorActionMessage(d, challenge, expiry) `` → `` msg, err := OperatorActionMessage(d, challenge, 1600) `` | `TestOperatorConfirmExpiry` | killed (TestOperatorConfirmExpiry) | killed (TestOperatorConfirmExpiry,TestOperatorConfirmNegativeExpiry) | K |
| `SSH-MSG-ERR` | M3 · `operator_signed.go` | `` msg, err := OperatorActionMessage(d, challenge, expiry)⏎if err != nil… `` → `` msg, _ := OperatorActionMessage(d, challenge, expiry) `` | `TestOperatorConfirmNegativeExpiry` | killed (TestOperatorConfirmNegativeExpiry) | killed (TestOperatorConfirmNegativeExpiry) | K |
| `PK-HASH-CHALLENGE` | M3 · `operator_signed.go` | `` h, hErr := OperatorActionMessageHash(d, challenge, expiry) `` → `` h, hErr := OperatorActionMessageHash(d, [32]byte{}, expiry) `` | `TestOperatorConfirmBothKeys` | killed (TestOperatorConfirmBothKeys) | killed (TestOperatorConfirmAuthenticity,TestOperatorConfirmBothKeys,TestOperatorConfirmExpiry,TestOperatorConfirmRevocationCounter,TestOperatorRunnerRecheckLaw) | K |
| `PK-HASH-VERB` | M3 · `operator_signed.go` | `` h, hErr := OperatorActionMessageHash(d, challenge, expiry) `` → `` h, hErr := VerbMessageHash(d, challenge, expiry) `` | `TestOperatorConfirmBinding` | killed (TestOperatorConfirmBinding) | killed (TestOperatorConfirmAuthenticity,TestOperatorConfirmBinding,TestOperatorConfirmBothKeys,TestOperatorConfirmExpiry,TestOperatorConfirmRevocationCounter,TestOperatorRunnerRecheckLaw) | K |
| `PK-STEPS-OFF` | M3 · `operator_signed.go` | `` ad, err := assertionSteps(artifact, ap, h, hErr) // W1-W14⏎if err != … `` → `` ad, err := assertionSteps(artifact, ap, h, hErr) // W1-W14⏎if false &… `` | `TestOperatorConfirmAuthenticity` | killed (TestOperatorConfirmAuthenticity) | killed (TestOperatorConfirmAuthenticity,TestOperatorConfirmBinding,TestOperatorConfirmNegativeExpiry,TestOperatorConfirmProof,TestOperatorRunnerRecheckLaw) | K |
| `W15-OFF` | M3 · `operator_signed.go` | `` if confirm && (ad.SignCount != 0 \|\| p.SignCount != 0) && ad.SignCount… `` → `` if false { // W15 `` | `TestOperatorConfirmRevocationCounter` | killed (TestOperatorConfirmRevocationCounter) | killed (TestOperatorConfirmRevocationCounter) | K |
| `W15-STRICT` | M3 · `operator_signed.go` | `` && ad.SignCount <= p.SignCount { // W15 `` → `` && ad.SignCount < p.SignCount { // W15 `` | `TestOperatorConfirmRevocationCounter` | killed (TestOperatorConfirmRevocationCounter) | killed (TestOperatorConfirmRevocationCounter) | K |
| `W15-ZERO` | M3 · `operator_signed.go` | `` (ad.SignCount != 0 \|\| p.SignCount != 0) && ad.SignCount <= p.SignCount `` → `` ad.SignCount <= p.SignCount `` | `TestOperatorConfirmRevocationCounter` | killed (TestOperatorConfirmRevocationCounter) | killed (TestOperatorConfirmRevocationCounter) | K |
| `W15-ZERO-AD` | M3 · `operator_signed.go` | `` (ad.SignCount != 0 \|\| p.SignCount != 0) && ad.SignCount <= p.SignCount `` → `` (ad.SignCount != 0) && ad.SignCount <= p.SignCount `` | `TestOperatorConfirmRevocationCounter` | killed (TestOperatorConfirmRevocationCounter) | killed (TestOperatorConfirmRevocationCounter) | K |
| `W15-ZERO-P` | M3 · `operator_signed.go` | `` (ad.SignCount != 0 \|\| p.SignCount != 0) && ad.SignCount <= p.SignCount `` → `` (p.SignCount != 0) && ad.SignCount <= p.SignCount `` | `TestOperatorConfirmRevocationCounter` | survived | survived — **equivalent**: the reused W15 is `p ≠ 0 ∧ ad ≤ p` already: with `p = 0`, `ad ≤ 0` forces `ad = 0` and the first disjunct is false (F5); V3 carries the same form | E |
| `R2-OFF` | M3 · `operator_signed.go` | `` if !confirm && ad.SignCount > p.SignCount { // R2 `` → `` if false { // R2 `` | `TestOperatorRunnerRecheckLaw` | killed (TestOperatorRunnerRecheckLaw) | killed (TestOperatorRunnerRecheckLaw) | K |
| `R2-AT` | M3 · `operator_signed.go` | `` if !confirm && ad.SignCount > p.SignCount { // R2 `` → `` if !confirm && ad.SignCount >= p.SignCount { // R2 `` | `TestOperatorRunnerRecheckLaw` | killed (TestOperatorRunnerRecheckLaw) | killed (TestOperatorRunnerRecheckLaw) | K |
| `PROOF-COUNT` | M3 · `operator_signed.go` | `` s.credentialID, s.signCount = bytes.Clone(p.CredentialID), ad.SignCou… `` → `` s.credentialID, s.signCount = bytes.Clone(p.CredentialID), p.SignCount `` | `TestOperatorConfirmBothKeys` | killed (TestOperatorConfirmBothKeys) | killed (TestOperatorConfirmBothKeys) | K |
| `PROOF-DIGEST` | M3 · `operator_signed.go` | `` s := &SignedOperatorAction{ok: true, digest: d, `` → `` s := &SignedOperatorAction{ok: true, digest: [32]byte{}, `` | `TestOperatorConfirmBothKeys` | killed (TestOperatorConfirmBothKeys) | killed (TestOperatorConfirmBothKeys,TestOperatorConfirmExpiry,TestOperatorRunnerRecheckLaw) | K |
| `PROOF-CHALLENGE` | M3 · `operator_signed.go` | `` challenge: challenge, expiry: expiry, name: ap.Name `` → `` challenge: [32]byte{}, expiry: expiry, name: ap.Name `` | `TestOperatorConfirmBothKeys` | killed (TestOperatorConfirmBothKeys) | killed (TestOperatorConfirmBothKeys,TestOperatorConfirmExpiry,TestOperatorRunnerRecheckLaw) | K |
| `PROOF-EXPIRY` | M3 · `operator_signed.go` | `` expiry: expiry, name: ap.Name `` → `` expiry: 0, name: ap.Name `` | `TestOperatorConfirmBothKeys` | killed (TestOperatorConfirmBothKeys) | killed (TestOperatorConfirmBothKeys,TestOperatorConfirmExpiry,TestOperatorRunnerRecheckLaw) | K |
| `PROOF-FORMAT` | M3 · `operator_signed.go` | `` name: ap.Name, format: format} `` → `` name: ap.Name, format: ap.Format + "x"} `` | `TestOperatorConfirmBothKeys` | killed (TestOperatorConfirmBothKeys) | killed (TestOperatorConfirmBothKeys) | K |
| `BINDS-NIL` | M3 · `operator_signed.go` | `` func (s *SignedOperatorAction) Binds(a OperatorAction, challenge [32]… `` → `` func (s *SignedOperatorAction) Binds(a OperatorAction, challenge [32]… `` | `TestOperatorConfirmBothKeys` | killed (TestOperatorConfirmBothKeys) | killed (TestOperatorConfirmBothKeys) | K |
| `BINDS-ERR` | M3 · `operator_signed.go` | `` case err != nil:⏎ `` → `` case err != nil && false:⏎ `` | `TestOperatorConfirmBothKeys` | killed (TestOperatorConfirmBothKeys) | killed (TestOperatorConfirmBothKeys) | K |
| `BINDS-DIGEST` | M3 · `operator_signed.go` | `` case s.digest != d:⏎ `` → `` case s.digest != d && false:⏎ `` | `TestOperatorConfirmBothKeys` | killed (TestOperatorConfirmBothKeys) | killed (TestOperatorConfirmBothKeys) | K |
| `BINDS-CHALLENGE` | M3 · `operator_signed.go` | `` case s.challenge != challenge:⏎return &BindingError{Field: "challenge… `` → ``  `` | `TestOperatorConfirmBothKeys` | killed (TestOperatorConfirmBothKeys) | killed (TestOperatorConfirmBothKeys) | K |
| `BINDS-EXPIRY` | M3 · `operator_signed.go` | `` case s.expiry != expiry:⏎return &BindingError{Field: "expiry"}⏎ `` → ``  `` | `TestOperatorConfirmBothKeys` | killed (TestOperatorConfirmBothKeys) | killed (TestOperatorConfirmBothKeys) | K |
| `APPROVEDBY-OK` | M3 · `operator_signed.go` | `` func (s *SignedOperatorAction) ApprovedBy() string {⏎if s == nil \|\| !… `` → `` func (s *SignedOperatorAction) ApprovedBy() string {⏎if s == nil { `` | `TestOperatorConfirmBothKeys` | killed (TestOperatorConfirmBothKeys) | killed (TestOperatorConfirmBothKeys) | K |
| `PASSKEYUSE-FORMAT` | M3 · `operator_signed.go` | `` if s == nil \|\| !s.ok \|\| s.format != FormatWebAuthn { `` → `` if s == nil \|\| !s.ok { `` | `TestOperatorConfirmBothKeys` | killed (TestOperatorConfirmBothKeys) | killed (TestOperatorConfirmBothKeys) | K |
| `PASSKEYUSE-OK` | M3 · `operator_signed.go` | `` if s == nil \|\| !s.ok \|\| s.format != FormatWebAuthn { `` → `` if s == nil \|\| s.format != FormatWebAuthn { `` | `TestOperatorConfirmBothKeys` | survived | survived — **equivalent**: no reachable value has `format == webauthn` and `ok == false`: the fields are unexported and both laws set `ok` on every value they build; nil and the zero value are refused by the remaining conditions | E |
| `REUSED-W7-ORIGIN` | — · `passkey_assert.go (reused)` | `` if cd.Origin != p.Origin { // W7 `` → `` if false { // W7 `` | `TestOperatorConfirmAuthenticity` | killed (TestOperatorConfirmAuthenticity) | killed (TestApproveByAssertionOrder,TestConfirmByAssertionLaw,TestDecideByAssertionSteps,TestOperatorConfirmAuthenticity) | K |
| `REUSED-W10-RPID` | — · `passkey_assert.go (reused)` | `` if ad.RPIDHash != sha256.Sum256([]byte(p.RPID)) { // W10 `` → `` if false { // W10 `` | `TestOperatorConfirmAuthenticity` | killed (TestOperatorConfirmAuthenticity) | killed (TestApproveByAssertionOrder,TestConfirmByAssertionLaw,TestDecideByAssertionSteps,TestOperatorConfirmAuthenticity) | K |
| `REUSED-W11-UP` | — · `passkey_assert.go (reused)` | `` if !ad.UP { // W11 `` → `` if false { // W11 `` | `TestOperatorConfirmAuthenticity` | killed (TestOperatorConfirmAuthenticity) | killed (TestApproveByAssertionOrder,TestConfirmByAssertionLaw,TestDecideByAssertionSteps,TestOperatorConfirmAuthenticity) | K |
| `REUSED-W12-UV` | — · `passkey_assert.go (reused)` | `` if !ad.UV { // W12 `` → `` if false { // W12 `` | `TestOperatorConfirmAuthenticity` | killed (TestOperatorConfirmAuthenticity) | killed (TestApproveByAssertionOrder,TestConfirmByAssertionLaw,TestDecideByAssertionSteps,TestOperatorConfirmAuthenticity) | K |
| `REUSED-W8-CROSS` | — · `passkey_assert.go (reused)` | `` if cd.CrossOrigin \|\| cd.HasTopOrigin { // W8 `` → `` if false { // W8 `` | `TestOperatorConfirmAuthenticity` | killed (TestOperatorConfirmAuthenticity) | killed (TestApproveByAssertionOrder,TestConfirmByAssertionLaw,TestDecideByAssertionSteps,TestOperatorConfirmAuthenticity) | K |
| `REUSED-W2-CRED` | — · `passkey_assert.go (reused)` | `` if !bytes.Equal(resp.RawID, p.CredentialID) { // W2 `` → `` if false { // W2 `` | `TestOperatorConfirmAuthenticity` | killed (TestOperatorConfirmAuthenticity) | killed (TestApproveByAssertionOrder,TestConfirmByAssertionLaw,TestDecideByAssertionRefusesAMissingMessageAtW6,TestDecideByAssertionSteps,TestOperatorConfirmAuthenticity) | K |
| `REUSED-P4-REVOKED` | — · `signed.go (reused)` | `` if a.RevokedAt != 0 { // P4 `` → `` if false { // P4 `` | `TestOperatorConfirmRevocationCounter` | killed (TestOperatorConfirmRevocationCounter) | killed (TestApproveByAssertionOrder,TestApproveBySignatureOrder,TestCheckConfirmation,TestConfirmByAssertionLaw,TestConfirmBySignatureLaw,TestDecideByAssertionSteps,TestDecideBySignatureApprover,TestOperatorConfirmRevocationCounter,TestOperatorRunnerRecheckLaw) | K |
| `REUSED-P5-ENROLLED` | — · `signed.go (reused)` | `` if admit.Now < a.EnrolledAt { // P5 `` → `` if false { // P5 `` | `TestOperatorConfirmRevocationCounter` | killed (TestOperatorConfirmRevocationCounter) | killed (TestApproveByAssertionOrder,TestApproveBySignatureOrder,TestCheckConfirmation,TestConfirmByAssertionLaw,TestConfirmBySignatureLaw,TestDecideByAssertionSteps,TestDecideBySignatureApprover,TestOperatorConfirmRevocationCounter) | K |
| `REUSED-P3W` | — · `signed.go (reused)` | `` if format == FormatWebAuthn && (a.Format != format \|\| a.Passkey == ni… `` → `` if format == FormatWebAuthn && a.Passkey == nil { // P3w `` | `TestOperatorConfirmRevocationCounter` | killed (TestOperatorConfirmRevocationCounter) | killed (TestApproveByAssertionOrder,TestCheckConfirmation,TestConfirmByAssertionLaw,TestDecideByAssertionSteps,TestOperatorConfirmRevocationCounter) | K |
| `J407-SETOUTPUT` | S1 · `host/verbs operator.go` | `` fs.SetOutput(io.Discard)⏎ `` → `` _ = io.Discard⏎ `` | `TestAdmitOperatorWritesNothing` | killed (TestAdmitOperatorWritesNothing) | killed (TestAdmitOperatorWritesNothing) | K |

**Totals: 150 mutants (NM excluded; 0 NM, 0 NM-text in run 2), 146 killed by their named test and by the whole package, 4 named equivalents, 0 survivors, 0 pending.** No kill rests on the golden alone (the `w/o golden` leg never triggered: every golden kill is shared with a behavioural test).

---

## 5. Files (final `git diff --name-status 8db22168 HEAD`)

```
A	host/approval/operator_action.go
A	host/approval/operator_action_test.go
A	host/approval/operator_action_valid_test.go
A	host/approval/operator_signed.go
A	host/approval/operator_signed_test.go
A	host/approval/operator_silent_test.go
A	host/approval/testdata/operator-action.golden.json
A	host/verbs/operator_silent_test.go
```

Two new production files; no pre-existing file of any kind; nothing deleted; no `cmd`, `docs`, `app.js`, `daemon.go`,
`surface_page.go`, `host/store`, `host/authority`, `.ail`, `world`, `packages`, `scripts`, `tools` or `plugins` file.

---

## 6. Sprint gates (controller, after the last sub-milestone's commit; each measured at the base, §1.1)

`bash gates.sh <worktree> $HOME/.ailang/state/world-iter408/gates final <HEAD sha>` (one bash shell; `GATES_SKIP_VERIFY_GO=1`
skips `verify_go` when the controller runs it separately — `VERIFY_GO_LEGS=plain AILANG_BIN=$HOME/.pinned-ailang/ailang
./scripts/verify_go.sh` after commit). Each leg is its own subshell with its own rc, `tee`'d to `<logdir>/<leg>.txt`; no leg is
chained to another; `ALL_DONE.<sha>` holds `FAILS: …`. `TOUCHED_PACKAGES` is derived from `git diff --name-only` and must equal
`host/approval host/verbs`. `-race` runs on both touched packages whole and on their dependant host/authority whole. S2 has no
Chrome or node test: no `chrome` or `node` leg. Never start it under `nohup`.

````bash gates.sh
#!/bin/bash
# Row 243 V4 S2 gates (plan §6). usage: gates.sh <tree> <logdir> <base|final> [<sha>]
# The controller: bash gates.sh <worktree> $HOME/.ailang/state/world-iter408/gates final <HEAD sha>
# macOS /bin/bash 3.2: no declare -A. Every leg records its own rc; every leg runs (no set -e, no leg
# chained to another with &&); every leg is tee'd to <logdir>/<leg>.txt. GATES_SKIP_VERIFY_GO=1 skips
# verify_go (the executor sets it; the controller runs it). Never start this under nohup. Writes
# <logdir>/ALL_DONE.<sha> with FAILS: when every leg has finished. S2 has no Chrome test: no chrome leg.
T=$1; L=$2; MODE=$3
mkdir -p "$L"; cd "$T" || exit 2
SHA=${4:-$(git rev-parse HEAD 2>/dev/null || echo unknown)}
B=8db22168689f31e19349b075e0199c8d9c940e61
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
# host/approval host/verbs; a difference FAILS.
TOUCHED=$(git diff --name-only "$B" HEAD -- cmd host 2>/dev/null | grep '\.go$' | xargs -n1 dirname 2>/dev/null | sort -u | tr '\n' ' ')
echo "TOUCHED_PACKAGES: $TOUCHED" | tee -a "$L/legs.txt"
if [ "$MODE" = final ] && [ "$TOUCHED" != "host/approval host/verbs " ]; then FAILS="$FAILS touched_packages_differ"; fi
leg status 'git status --porcelain | wc -l'
leg verify_ail './scripts/verify_ail.sh'
leg gofmt 'out=$(gofmt -l cmd host); echo "$out"; test -z "$out"'
leg build 'go build ./...'
leg vet 'go vet ./...'
leg fence "go test -run '^\$' ./... > '$L/fence.raw' 2>&1; rc=\$?; cat '$L/fence.raw'; if grep -E '^(FAIL|--- FAIL)' '$L/fence.raw' >/dev/null; then exit 1; fi; exit \$rc"
# the two touched packages whole and their dependant host/authority (it imports host/approval), each its own rc
leg pkgs "r=0; for p in ./host/approval/ ./host/verbs/ ./host/authority/; do go test -count=1 -timeout 30m \$p; r=\$(( r | \$? )); done; exit \$r"
# -race on every touched package whole, and on the dependant host/authority whole
leg race "r=0; for p in ./host/approval/ ./host/verbs/ ./host/authority/; do go test -race -count=1 -timeout 30m \$p; r=\$(( r | \$? )); done; exit \$r"
NEW_A="TestOperatorActionGolden TestOperatorActionMessage TestOperatorStrictJSON TestOperatorRequestLimits TestOperatorActionValidate TestOperatorDigestsAreDistinct TestOperatorArgvIsTheBarsIntent TestOperatorConfirmBothKeys TestOperatorConfirmProof TestOperatorConfirmAuthenticity TestOperatorConfirmRevocationCounter TestOperatorConfirmBinding TestOperatorConfirmExpiry TestOperatorConfirmNegativeExpiry TestOperatorRunnerRecheckLaw TestOperatorPureEntryPointsWriteNothing"
NEW_V="TestAdmitOperatorWritesNothing"
if [ "$MODE" = base ]; then NEW_A=""; NEW_V=""; fi
INV_APPROVAL="TestApprovalImportsOnlyTheStandardLibrary TestSSHSIGCrossCheckWithSSHKeygen TestVerbDomainSeparation"
INV_APPROVALTEST="TestApprovalTestIsTestOnly TestTestSignerAgreesWithSSHKeygen"
INV_DAEMON="TestDaemonDependencyAllowlist TestFrozenV1RouteTableMatchesMux TestSurfaceAppRoutes TestApprovalRoutesRegistered TestDecisionRoutesRegistered TestConfirmRoutesRegistered TestDaemonCredentialWriterSites TestSurfaceVerbRouteRegistered TestSurfaceVerbIsReadOnly TestSurfaceVerbOutcomes"
INV_BOUNDARY="TestWorldBoundaryDependencyAllowlist TestBareNetHTTPExemptionIsPerGroup TestEveryCommandDirectoryIsAProtectedGroup TestBoundaryASTWriteGuard"
INV_STORE="TestProductionContextRoots TestProductionGoSurface TestStoreSurfaceTakesContext TestCredentialsSchemaMatchesCanonicalDDL TestCredentialsSchema9VerbRequests TestCredentialsNewerThanNineIsRefused"
INV_VERBS="TestLookup TestApproveTarget TestOperatorGrammarCase TestOperatorIdentifierAnchors TestOperatorNeverDispatch"
INV_CMD="TestEveryVerbHasHelp TestVerbTableMatchesTheCLISwitch TestCLICredentialWriterSites TestOperatorFlagParity"
leg inv "set -o pipefail; r=0
go test -v -count=1 ./host/approval/ -run '^($(bar $INV_APPROVAL)$(alt $NEW_A))\$'; r=\$(( r | \$? ))
go test -v -count=1 ./host/approval/approvaltest/ -run '^($(bar $INV_APPROVALTEST))\$'; r=\$(( r | \$? ))
go test -v -count=1 ./host/daemon/ -run '^($(bar $INV_DAEMON))\$'; r=\$(( r | \$? ))
go test -v -count=1 ./host/boundary/ -run '^($(bar $INV_BOUNDARY))\$'; r=\$(( r | \$? ))
go test -v -count=1 ./host/store/ -run '^($(bar $INV_STORE))\$'; r=\$(( r | \$? ))
go test -v -count=1 ./host/verbs/ -run '^($(bar $INV_VERBS)$(alt $NEW_V))\$'; r=\$(( r | \$? ))
go test -v -count=1 ./cmd/ailang-worldd/ -run '^($(bar $INV_CMD))\$'; r=\$(( r | \$? ))
exit \$r"
need inv $INV_APPROVAL $NEW_A $INV_APPROVALTEST $INV_DAEMON $INV_BOUNDARY $INV_STORE $INV_VERBS $NEW_V $INV_CMD
if grep -E '^\s*--- (SKIP|FAIL)' "$L/inv.txt"; then FAILS="$FAILS inv_skip_or_fail"; fi
echo "named PASS: $(grep -cE '^--- PASS' "$L/inv.txt")" | tee -a "$L/legs.txt"
# No new daemon route: every mux registration in daemon.go, counted (base = final), and the /surface POSTs (3)
ROUTES=$(grep -cE 'mux\.(HandleFunc|Handle)\(' host/daemon/daemon.go); POSTS=$(grep -cE 'mux\.HandleFunc\("POST /surface/' host/daemon/daemon.go)
BROUTES=$(git show "$B:host/daemon/daemon.go" | grep -cE 'mux\.(HandleFunc|Handle)\(')
echo "DAEMON_ROUTES: $ROUTES (base $BROUTES) SURFACE_POSTS: $POSTS" | tee -a "$L/legs.txt"
if [ "$ROUTES" != "$BROUTES" ] || [ "$POSTS" != 3 ]; then FAILS="$FAILS routes_changed"; fi
# Production bytes: S2's only production files are host/approval/operator_action.go and operator_signed.go (NEW).
# Every other changed file is a _test.go or the golden; nothing under host/daemon, host/store, host/authority,
# host/surfacepage, cmd, docs, world, packages, scripts, tools, plugins; no .ail.
NOPROD=$(git diff --name-only "$B" HEAD | grep -vE '(_test\.go$|^host/approval/operator_(action|signed)\.go$|^host/approval/testdata/operator-action\.golden\.json$|^design_docs/planned/(w-world-address-bar-verbs-v4s2-sprint-plan\.md|sprint_w-world-address-bar-verbs-v4s2\.json)$)' | tr '\n' ' ')
NOPROD="$NOPROD$(git diff --name-only "$B" HEAD -- cmd docs host/daemon host/store host/authority host/surfacepage world packages scripts tools plugins '*.ail' | grep -v '_test\.go$' | tr '\n' ' ')"
echo "UNTOUCHED_VIOLATIONS: ${NOPROD:-none}" | tee -a "$L/legs.txt"
if [ "$MODE" = final ] && [ -n "$NOPROD" ]; then FAILS="$FAILS untouched_violated"; fi
# Pre-existing files changed or deleted: none (every S2 path is added).
PRE=$(git diff --name-status "$B" HEAD | awk '$1 != "A"' | tr '\n' ' ')
echo "PREEXISTING_CHANGED: ${PRE:-none}" | tee -a "$L/legs.txt"
if [ "$MODE" = final ] && [ -n "$PRE" ]; then FAILS="$FAILS preexisting_changed"; fi
# The envelope and the law have no production caller (a caller search with a positive control in the same call).
CALLERS=$(grep -rnE '\b(OperatorAction|OperatorActionSchema|OperatorActionPrefix|OperatorActionNamespace|OperatorActionMessage|OperatorActionMessageHash|OperatorReviewDigest|OperatorSignatureDigest|DecodeOperatorAction|CanonicalOperatorAction|OperatorArgv|OperatorOptions|OperatorResolved|ConfirmOperatorAction|CheckOperatorAction|SignedOperatorAction|OperatorConfirmInput|OperatorCheckInput|MaxOperatorActionBytes|MaxOperatorTaskBytes|OperatorActionLifetime)\b' --include='*.go' cmd host | grep -v '_test\.go:' | grep -vE '^host/approval/operator_(action|signed)\.go:' | tr '\n' ' ')
CONTROL=$(grep -rnE '\bapproval\.CheckConfirmation\(' --include='*.go' cmd host | grep -v '_test\.go:' | wc -l | tr -d ' ')
echo "ENVELOPE_PRODUCTION_CALLERS: ${CALLERS:-none} (control: approval.CheckConfirmation production callers $CONTROL)" | tee -a "$L/legs.txt"
if [ "$MODE" = final ] && { [ -n "$CALLERS" ] || [ "$CONTROL" -lt 1 ]; }; then FAILS="$FAILS envelope_has_a_caller"; fi
# Served bytes unchanged: app.js, the CSP literal's file and the store's schema
for f in host/surfacepage/app/app.js host/daemon/surface_page.go host/store/credentials.go; do
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
  only writer** of the tree; no hand edit, no `gofmt -w`, **no golden regeneration** (`WORLD_UPDATE_OPERATOR_ACTION_GOLDEN`
  is unset by `accept.sh`), no fuzz. A failed `--dry-run` or a red acceptance → **stop and report** with the log; do not repair.
- **Never `nohup`, never background a test leg; pass a 600000 ms timeout to every long Bash call.** Each acceptance part is
  one call (C14 measured each under 2 minutes). **The executor never launches Chrome** and never sets `WORLD_CHROME`.
- **Runs**: **one run, M1 → M2 → M3 → M4** (each patch assumes the previous tree; the controller may commit between
  sub-milestones). One part each: `fast`. `accept.sh` exports `AILANG_BIN=$HOME/.pinned-ailang/ailang`,
  `WORLD_EXEC_NODE=/opt/homebrew/bin/node` and `WORLD_REQUIRE_SSH_KEYGEN=1`, and unsets `WORLD_CHROME`,
  `WORLD_REQUIRE_PASSKEY_BROWSER` and every golden-update variable.
- **Commit-time test commands** (`commit_ms.sh <M> <msgfile> <test-cmd…>`; each carries its env through `accept.sh`; with
  `T=$HOME/.ailang/state/world-iter408/exec-patches W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter408-row243-v4s2`):
  **M1** `bash $T/accept.sh $W M1 fast`; **M2** `bash $T/accept.sh $W M2 fast`; **M3** `bash $T/accept.sh $W M3 fast`; **M4**
  `bash $T/accept.sh $W M4 fast`. `commit_ms.sh` vets the directories of the changed `.go` files only (`host/approval`, and
  `host/verbs` at M4); the golden is not `.go` and is copied with the snapshot.
- After each sub-milestone's acceptance, snapshot **cumulatively against the base** to
  `$HOME/.ailang/state/world-iter408/snap/<M>/`: every file added or changed since `8db22168` (repo-relative paths), a
  **`DELETED`** file (one repo path per line; **empty for every S2 stage**), a **`STATUS`** file, and `MANIFEST.sha256`
  (the two plan files, `STATUS` and `DELETED` are not in it). **Never an `rsync` of whole directories.** The `find` takes
  only the top-level dirs that exist (every S2 stage has only `host/`; no stage has `cmd/` or `docs/`), rehearsed on all four
  (C13):
  ```bash
  bash <<'EOF2'
  set -euo pipefail
  M=M1   # the sub-milestone just accepted
  B=8db22168689f31e19349b075e0199c8d9c940e61
  W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter408-row243-v4s2
  D="$HOME/.ailang/state/world-iter408/snap/$M"; rm -rf "${D:?}"; mkdir -p "$D"
  cd "$W"
  { git diff --name-status --no-renames "$B" -- cmd host docs; git ls-files --others --exclude-standard -- cmd host docs | sed 's/^/A\t/'; } | sort -u -t "$(printf '\t')" -k2 > "$D/STATUS"
  awk -F'\t' '$1 != "D" {print $2}' "$D/STATUS" | while read -r p; do mkdir -p "$D/$(dirname "$p")"; cp -p "$p" "$D/$p"; done
  awk -F'\t' '$1 == "D" {print $2}' "$D/STATUS" > "$D/DELETED"
  (cd "$D" && find $(for d in cmd host docs; do [ -d "$d" ] && echo "$d"; done) -type f -print0 | sort -z | xargs -0 shasum -a 256) > "$D/MANIFEST.sha256"
  cat "$D/STATUS"; echo "DELETED: $(wc -l < "$D/DELETED")"; echo "MANIFEST: $(wc -l < "$D/MANIFEST.sha256")"
  EOF2
  ```
  Expected (`STATUS` lines / `DELETED` / MANIFEST): **M1** 3 / 0 / 3; **M2** 4 / 0 / 4; **M3** 6 / 0 / 6; **M4** 8 / 0 / 8. The
  controller copies each snapshot over the worktree, commits it as one commit, and checks the numstat (`git diff -M --numstat
  8db22168 HEAD -- cmd host docs app | awk '{s+=$1+$2} END {print s}'`): **M1 530, M2 773, M3 1304, M4 1394**. A different
  number → stop.
- **Stop conditions**: a dry-run failure; any `--- FAIL`, `--- SKIP` or `FAIL` line in an acceptance (`ACCEPT_RED`); gofmt
  output; a numstat mismatch; any file outside §5 changed (`git status --short` after each patch must equal the cumulative list
  of §3's paragraphs); a numstat above 1,400.
- **Known load-sensitive tests** (the directive's list) are not this sprint's; none is in an S2 acceptance part (the
  packages run whole are host/approval, approvaltest, host/verbs and host/authority, each a few seconds).
- D-WORLD-109: nothing touches `~/.ailang/world`, `~/.ssh`, a real `/dev/tty` or a real key; no S2 test opens a store.

---

## 8. Risks

- **R1 — margin 6 lines under the cap.** The patches are exact, so the numstat is exact (C16 measured it after a real
  `patch -p1 -E` and commit of each stage). Any executor-side edit is a stop condition anyway.
- **R2 — the envelope's `Store` and `Resolution` are typed now, filled later** (S2-D5). S4/S5 choose the values; a new fact is a
  name, not a schema change. If S5 needs a structured value, the envelope schema moves to `v2` — a deliberate, visible change.
- **R3 — the law is reachable only through a caller that supplies the approver.** That is why the parked adapter exists: until
  S3 lands it, nothing can call the law in production (C8), so the risk of a caller trusting a hand-built `Approver` does not
  exist yet; S3 must land the adapter before any route uses the law (Appendix C).
- **R4 — `strings.HasPrefix(err.Error(), "json: unknown field ")`** classifies the unknown-key error by text (encoding/json has
  no typed error for it). If Go changes the text, the class falls back to `ErrDecode` — still a refusal, never an admission;
  `TestOperatorStrictJSON` would fail by name, and `D7-UNKNOWN-CLASS` pins it.
- **R5 — the writes-nothing tests swap process-global `os.Stdout`/`os.Stderr`.** Neither test is `t.Parallel`, so Go runs
  them with no other test of the package concurrently; under `-race` both packages passed (§1.5).

---

## 9. Acceptance criteria (the judge scores against these)

1. **No new daemon route, no schema change, no `.ail`, `surfaceCSP` and `app.js` unchanged, no served-behaviour change**:
   `DAEMON_ROUTES: 47 (base 47)`, `SURFACE_POSTS: 3`, `UNTOUCHED_VIOLATIONS: none`, `UNCHANGED:` app.js / surface_page.go /
   host/store/credentials.go, `verify_ail` identical to the base's; `TestFrozenV1RouteTableMatchesMux`,
   `TestConfirmRoutesRegistered`, `TestDaemonCredentialWriterSites`, `TestCLICredentialWriterSites`,
   `TestSurfaceVerbIsReadOnly`, `TestApprovalImportsOnlyTheStandardLibrary` green unedited.
2. **The envelope and the law have no production caller**: `ENVELOPE_PRODUCTION_CALLERS: none` over all 21 NEW exported
   identifiers, with the positive control ≥ 1.
3. **The proof cannot be built outside its law**: `SignedOperatorAction`'s fields are unexported (a composite literal with a
   field outside host/approval does not compile; the parked adapter's reflect test checks it), and the zero value and nil bind
   nothing (`TestOperatorConfirmBothKeys`: `ErrNotVerified`, `ApprovedBy() == ""`, no passkey use; mutants `BINDS-NIL`,
   `APPROVEDBY-OK`, `PASSKEYUSE-FORMAT` killed).
4. **The canonical bytes and both digests are golden-pinned** by the real encoder (C6), the message layout field by field
   (`TestOperatorActionMessage`), and every envelope field is bound (each `FIELD-DROP-*` killed; those also bound through
   `argv` are killed by the golden and by `TestOperatorStrictJSON`/`…RequestLimits`, behavioural fixture tests).
5. **Every invalid input class is refused before canonicalization, each by its own error, with no action** (`TestOperatorStrictJSON`;
   `D1`–`D9` killed), and both caps at the exact boundary (`TestOperatorRequestLimits`; `O10-*`, `O13-*`, `D1-*`, `TASK-CAP`,
   `ACTION-CAP` killed).
6. **Each §5.2 check is killed independently by a named test, whole package**: UV (`REUSED-W12-UV`), RP id (`…W10-RPID`), origin
   (`…W7-ORIGIN`), namespace (`SSH-NAMESPACE`, `NAMESPACE`), enrollment (`REUSED-P5-ENROLLED`, `P2-P5-NOW`), revocation
   (`REUSED-P4-REVOKED`, `P2-P5-OFF`), counter (`W15-*`, `R2-*`), challenge (`SSH-MSG-CHALLENGE`, `PK-HASH-CHALLENGE`), expiry
   (`P1-ADMIT`, `C2-*`, `R1-*`, `SSH-MSG-EXPIRY`), each binding field (`FIELD-DROP-*`, `C1-ROW-DIGEST`), the domain separation
   (`PK-HASH-VERB`, `MSG-PREFIX`).
7. **Every time comparison's operator is its own mutant at the strict-equality boundary** (C10), each killed.
8. **Every refusal asserts its class and that no proof (or no action) is produced** (`opRefused`, `opDecodeRefused`).
9. **Every pure entry point writes nothing** (`TestOperatorPureEntryPointsWriteNothing`), and S1's survivor `J407-SETOUTPUT`
   is killed by `TestAdmitOperatorWritesNothing` after reproducing SURVIVED on the base (C9).
10. **Size**: `git diff -M --numstat 8db22168 HEAD -- cmd host docs app` sums to **1,394** (≤ 1,400; the same unfiltered); each
    sub-milestone's numstat as §7.
11. **Exactly these files change: the eight added in §5. No pre-existing file — production or test — changes or is deleted**
    (`PREEXISTING_CHANGED: none`).
12. **Mutants**: every `K` row killed by its named test and by the whole package; the `E` rows named with reasons; the plan's
    ids are the ledger's.
13. **Gates**: §6's script `FAILS: none` (or only a directive-listed load-sensitive test, green alone), `verify_go` rc 0.

---

## Appendix C — after this sprint (S2's close; what S3–S12 own)

- **Moved out (S2-D1)**: the `host/authority` adapters — `JudgeOperatorAction(ctx, ro, OperatorJudgeInput, artifact)`,
  `RecheckOperatorAction(ctx, ro, action, challenge, expiry, now, artifact, approvedBy)` (the runner's approver-identity
  comparison, `ErrOperatorApprover`), `VerifiedOperatorAction` (unexported fields; `Binds`, `ApprovedBy`, `PasskeyUse`,
  `Rechecked`) and `TestOperatorJudgeBothKeys`, `TestOperatorRunnerRechecks`, `TestVerifiedOperatorActionOnlyFromTheJudges`
  (AST: one constructor literal, called by exactly the two judges, positive control). **260 changed lines, compiled, vetted,
  gofmt-clean, green on this plan's M4**: `$HOME/.ailang/state/world-iter408/plan-scratch/patches/PARK-authority-judges.patch`.
  **S3 lands it first** (it needs no schema; its input `OperatorJudgeInput` is S3's row), with its own mutant pass.
- **Closed here**: design §3.1's envelope, limits, digests, message and namespace; §5.1 `TestOperatorStrictJSON` and the envelope
  half of `TestOperatorRequestLimits` and `TestOperatorWireTypes`; §5.2 `TestOperatorConfirmBothKeys`, `…Proof`,
  `…Authenticity`, `…RevocationCounter` (pure), `…Binding`, `…Expiry` and the pure half of `TestOperatorRunnerRechecks`;
  judge-407's `J407-SETOUTPUT`.
- **§5.1/§5.2 rows S3–S12 own** (by design §8): **S3** schema 10 `operator_actions`, `operator_runner_state`,
  `operator_action_steps`, the receipts and CAS — `TestOperatorAdmissionBoundaries`, `TestOperatorDuplicateTap` (keyed on
  `OperatorReviewDigest`), `TestOperatorReconnect` (durable dedup), `TestOperatorConfirmReplay`, `TestOperatorConfirmCrash`,
  `TestOperatorClaimBeforeExpiry`, `TestOperatorClaimRace`, `TestOperatorTerminalExclusion`, `TestOperatorPrune`, the
  store half of `TestOperatorRunnerRechecks` (with the parked adapter); **S4** `TestOperatorRunnerUnavailable`,
  `TestOperatorCatalogueAdmission`, `TestOperatorRecoverUntouched`/`…Reserved` (runner side), `TestOperatorStatusDisconnected`,
  `TestOperatorCancelInflight`, the service rows (§5.5); **S5–S6** §5.3; **S7** the routes — `TestOperatorMisdirected`,
  `…CrossOrigin`, `…FetchSite`, `…MediaType`, `…BodyBound`, `…UnsupportedMethod`, `…RequestLimits` (body and task at the
  route), `…LiteralTask` (rendering), `…BarStartReview`, `…BarFinishReview`, `…OtherVerbParity` at the route; **S8**
  `…PanelNetwork`, `…PanelLatest`, `…LinkAllowlist`, `…WireTypes` (page), S1's F8 (`null` elements in `app.js` arrays);
  **S9–S11** §5.4; **S12** the herdr hand-off.
- **What S3 inherits from S2**: `OperatorAction` and its canonical bytes are what the `operator_actions` row stores (design §3.1
  "exact canonical resolved-intent bytes and digest"): store the canonical bytes, `OperatorSignatureDigest` as the digest
  column and `OperatorReviewDigest` as the dedup key; decode only with `DecodeOperatorAction`; `OperatorActionLifetime` (600 s)
  is the request expiry the row's `expires_at − created_at` may not exceed (C2); the confirm CAS binds
  `SignedOperatorAction.PasskeyUse()` and advances the counter once (R2 then requires the runner to see it recorded); the
  challenge (32 random bytes) and expiry are row columns, never envelope fields (F4). S3 also needs D-WORLD-119 A/B (design §8)
  and stays start-only unless D-WORLD-120 is answered (the audit residual).
- **Carried (not acted on)**: judge-407 note (3) — numeric CLI flags mirrored as strings refuse `--ttl x` as terminal not usage;
  the envelope carries no numeric option, so it does not matter here (S7's route renders the class). Judge-407 note (2) / S1 F8 —
  `null` array elements throw in `app.js`: S8's. Iteration 407 retro — an attended merge may land mid-sprint: S2 touches no
  shared file (`app.js` untouched).
- **Audit residuals (iteration 406, unchanged)**: an uncertain action has no attended resolver; a headless child outliving the
  runner is unowned; S3's schema 10 start-only unless D-WORLD-120 is answered; the second credentials writer vs `LOCK_NB` and
  the writer-site inventories move together; the herdr plugin hard-codes `codex=gpt-6.1-sol` (`world-config-init:96`).
- **Controller decisions requested**: (1) accept the S2-D1 cut — the adapters to S3 as its first milestone (or a short S2b
  before S3, 260 lines, answer-independent). (2) None on the ledger: D-WORLD-119/120 remain the ledger's; S2 is
  answer-independent.
- The controller records; this plan moves no design doc and edits no charter.

---

## Appendix P — the patches (apply with §2, in order; generated from the compiled scratch stages; do not edit)

### Patch M1 — the envelope and its strict decoder

````diff M1
diff --git a/host/approval/operator_action.go b/host/approval/operator_action.go
new file mode 100644
index 00000000..148bc6d5
--- /dev/null
+++ b/host/approval/operator_action.go
@@ -0,0 +1,335 @@
+package approval
+
+// Row 243 V4 sprint 2 M1 (design_docs/planned/w-world-address-bar-verbs-v4.md §3.1): the operator action
+// an enrolled key approves from the World app (host/verbs' OperatorIntent, canonicalized, never a re-parsed
+// bar line) and the 97-byte message it signs, with its own schema, prefix and SSHSIG namespace. Pure: no
+// I/O, no clock, no randomness. No production caller yet (S3 stores it, S7 serves it).
+
+import (
+	"bytes"
+	"crypto/sha256"
+	"encoding/binary"
+	"encoding/json"
+	"errors"
+	"fmt"
+	"io"
+	"path"
+	"reflect"
+	"regexp"
+	"slices"
+	"strings"
+	"time"
+	"unicode/utf8"
+)
+
+// The operator action's schema, signed-message prefix and SSHSIG namespace; the message is
+// 24 + 1 + 32 + 32 + 8 bytes; a signed request lives at most 600 s (design §3.1).
+const (
+	OperatorActionSchema     = "world/operator-action/v1"
+	OperatorActionPrefix     = "world-operator-action/v1"
+	OperatorActionNamespace  = "world-operator-action@sunholo"
+	OperatorActionMessageLen = len(OperatorActionPrefix) + 1 + 32 + 32 + 8
+	OperatorActionLifetime   = 10 * time.Minute
+)
+
+// MaxOperatorActionBytes bounds an action's canonical bytes; MaxOperatorTaskBytes bounds its task text,
+// in UTF-8 bytes, separately (design §3.1).
+const (
+	MaxOperatorActionBytes = 128 << 10
+	MaxOperatorTaskBytes   = 8 << 10
+)
+
+// operatorActionIDPattern is host/verbs' OperatorActionIDPattern (host/approval imports only the standard
+// library; TestOperatorArgvIsTheBarsIntent joins the two); operatorResolvedName is a resolution fact's
+// name; operatorPresetMax is host/verbs' MaxPresetName (the CLI's maxPresetNameLen).
+var (
+	operatorActionIDPattern = regexp.MustCompile(`^oa-[0-9a-f]{32}$`)
+	operatorResolvedName    = regexp.MustCompile(`^[a-z][a-z0-9_]{0,63}$`)
+)
+
+const operatorResolvedMax, operatorPresetMax = 64, 51
+
+// OperatorOptions are the admitted options of OperatorIntent.
+type OperatorOptions struct {
+	Preset           string `json:"preset"`
+	Headless         bool   `json:"headless"`           // start: --launch --headless
+	ConfirmOnSurface bool   `json:"confirm_on_surface"` // finish
+	Ready            bool   `json:"ready"`              // finish
+}
+
+// OperatorResolved is one frozen resolution fact (S5 names them); entries are in strictly ascending
+// name order, so one resolution has one encoding.
+type OperatorResolved struct {
+	Name  string `json:"name"`
+	Value string `json:"value"`
+}
+
+// OperatorAction is one operator action. Its canonical JSON (CanonicalOperatorAction), in this field
+// order, is what is digested, stored and signed.
+type OperatorAction struct {
+	Schema     string             `json:"schema"`     // OperatorActionSchema
+	ID         string             `json:"id"`         // ^oa-[0-9a-f]{32}$; "" only inside OperatorReviewDigest
+	Kind       string             `json:"kind"`       // "start" or "finish"
+	Episode    string             `json:"episode"`    // ^[a-z0-9][a-z0-9-]{0,63}$
+	Store      string             `json:"store"`      // the served credentials store, absolute and clean
+	Argv       []string           `json:"argv"`       // exactly OperatorArgv(Kind, Episode, Options)
+	Options    OperatorOptions    `json:"options"`    // the admitted options
+	Resolution []OperatorResolved `json:"resolution"` // the frozen resolution, 1 to 64 facts
+	Task       string             `json:"task"`       // start: 1 to 8 KiB of UTF-8; finish: ""
+	Review     string             `json:"review"`     // the exact review text the card displays
+}
+
+var (
+	ErrOperatorAction       = errors.New("approval: bad operator action")
+	ErrOperatorTooLarge     = errors.New("approval: operator action is over 128 KiB")
+	ErrOperatorTask         = errors.New("approval: operator task is empty, over 8 KiB or on a finish")
+	ErrOperatorUTF8         = errors.New("approval: operator action is not valid UTF-8")
+	ErrOperatorNUL          = errors.New("approval: operator action holds a NUL byte")
+	ErrOperatorShape        = errors.New("approval: operator action is not a JSON object")
+	ErrOperatorDuplicateKey = errors.New("approval: operator action repeats a key")
+	ErrOperatorUnknownKey   = errors.New("approval: operator action has an unknown key")
+	ErrOperatorTrailing     = errors.New("approval: operator action has data after its object")
+)
+
+// OperatorArgv is the canonical CLI line of an action: host/verbs AdmitOperator's Intent.Argv.
+func OperatorArgv(kind, episode string, o OperatorOptions) []string {
+	argv := []string{"ailang-worldd", "episode", kind, episode, "--preset", o.Preset}
+	switch kind {
+	case "start":
+		argv = append(argv, "--launch")
+		if o.Headless {
+			argv = append(argv, "--headless")
+		}
+	case "finish":
+		if o.ConfirmOnSurface {
+			argv = append(argv, "--confirm-on-surface")
+		}
+		if o.Ready {
+			argv = append(argv, "--ready")
+		}
+	}
+	return argv
+}
+
+// Validate refuses an action the daemon or the runner must never act on. Each refusal wraps an Err
+// above, ErrSchema, ErrString, ErrEpisode or ErrPath, and names its JSON field.
+func (a OperatorAction) Validate() error { return a.validate(true) }
+
+func (a OperatorAction) validate(withID bool) error {
+	bad := func(format string, args ...any) error {
+		return fmt.Errorf("%w: "+format, append([]any{ErrOperatorAction}, args...)...)
+	}
+	if a.Schema != OperatorActionSchema { // O1
+		return fmt.Errorf("%w: schema is %q, want %q", ErrSchema, a.Schema, OperatorActionSchema)
+	}
+	if err := stringLeaves(reflect.ValueOf(a), "", func(p, s string) error { // O2
+		if !utf8.ValidString(s) || strings.ContainsRune(s, 0) {
+			return fmt.Errorf("%w: %s", ErrString, p)
+		}
+		return nil
+	}); err != nil {
+		return err
+	}
+	if withID && !operatorActionIDPattern.MatchString(a.ID) { // O3
+		return bad("id %q does not match %s", a.ID, operatorActionIDPattern)
+	}
+	if a.Kind != "start" && a.Kind != "finish" { // O4
+		return bad("kind is %q, want start or finish", a.Kind)
+	}
+	if !episodePattern.MatchString(a.Episode) { // O5
+		return fmt.Errorf("%w: episode %q", ErrEpisode, a.Episode)
+	}
+	if !path.IsAbs(a.Store) || path.Clean(a.Store) != a.Store { // O6
+		return fmt.Errorf("%w: store %q", ErrPath, a.Store)
+	}
+	if p := a.Options.Preset; p == "se-tools" || len(p) > operatorPresetMax || !episodePattern.MatchString(p) { // O7
+		return bad("options.preset %q is not a launch preset name", p)
+	}
+	if a.Kind == "start" && (a.Options.ConfirmOnSurface || a.Options.Ready) { // O8
+		return bad("options: a start sets only preset and headless")
+	}
+	if a.Kind == "finish" && a.Options.Headless {
+		return bad("options: a finish has no headless")
+	}
+	if want := OperatorArgv(a.Kind, a.Episode, a.Options); !slices.Equal(a.Argv, want) { // O9
+		return bad("argv is %q, want %q", a.Argv, want)
+	}
+	if a.Kind == "start" && (len(a.Task) < 1 || len(a.Task) > MaxOperatorTaskBytes) { // O10
+		return fmt.Errorf("%w: task is %d bytes, want 1 to %d", ErrOperatorTask, len(a.Task), MaxOperatorTaskBytes)
+	}
+	if a.Kind == "finish" && a.Task != "" {
+		return fmt.Errorf("%w: a finish carries no task", ErrOperatorTask)
+	}
+	if len(a.Resolution) < 1 || len(a.Resolution) > operatorResolvedMax { // O11
+		return bad("resolution has %d facts, want 1 to %d", len(a.Resolution), operatorResolvedMax)
+	}
+	for i, r := range a.Resolution {
+		if !operatorResolvedName.MatchString(r.Name) {
+			return bad("resolution[%d].name %q", i, r.Name)
+		}
+		if i > 0 && a.Resolution[i-1].Name >= r.Name {
+			return bad("resolution[%d].name %q is not after %q", i, r.Name, a.Resolution[i-1].Name)
+		}
+	}
+	if a.Review == "" { // O12
+		return bad("review is empty")
+	}
+	c, err := CanonicalOperatorAction(a)
+	if err != nil {
+		return err
+	}
+	if len(c) > MaxOperatorActionBytes { // O13
+		return fmt.Errorf("%w: %d bytes, at most %d", ErrOperatorTooLarge, len(c), MaxOperatorActionBytes)
+	}
+	return nil
+}
+
+// CanonicalOperatorAction is the compact JSON of a in struct field order, HTML escaping off, no
+// trailing newline (CanonicalVerb's settings). It does not call Validate.
+func CanonicalOperatorAction(a OperatorAction) ([]byte, error) {
+	var buf bytes.Buffer
+	enc := json.NewEncoder(&buf)
+	enc.SetEscapeHTML(false)
+	if err := enc.Encode(a); err != nil {
+		return nil, err
+	}
+	return bytes.TrimSuffix(buf.Bytes(), []byte("\n")), nil
+}
+
+// DecodeOperatorAction is the strict inverse of CanonicalOperatorAction. Before any decoding it
+// refuses over-limit input, invalid UTF-8, a NUL byte, a non-object, a repeated key and trailing data;
+// then an unknown key, a type error, bytes that do not re-encode identically and an invalid action.
+func DecodeOperatorAction(b []byte) (OperatorAction, error) {
+	if len(b) > MaxOperatorActionBytes { // D1
+		return OperatorAction{}, fmt.Errorf("%w: %d bytes, at most %d", ErrOperatorTooLarge, len(b), MaxOperatorActionBytes)
+	}
+	if !utf8.Valid(b) { // D2
+		return OperatorAction{}, ErrOperatorUTF8
+	}
+	if bytes.IndexByte(b, 0) >= 0 { // D3
+		return OperatorAction{}, ErrOperatorNUL
+	}
+	if err := strictObject(b); err != nil { // D4-D6
+		return OperatorAction{}, err
+	}
+	dec := json.NewDecoder(bytes.NewReader(b))
+	dec.DisallowUnknownFields()
+	var a OperatorAction
+	if err := dec.Decode(&a); err != nil { // D7
+		if strings.HasPrefix(err.Error(), "json: unknown field ") {
+			return OperatorAction{}, fmt.Errorf("%w: %v", ErrOperatorUnknownKey, err)
+		}
+		return OperatorAction{}, fmt.Errorf("%w: %v", ErrDecode, err)
+	}
+	c, err := CanonicalOperatorAction(a)
+	if err != nil {
+		return OperatorAction{}, err
+	}
+	if !bytes.Equal(c, b) { // D8
+		return OperatorAction{}, ErrNotCanonical
+	}
+	if err := a.Validate(); err != nil { // D9
+		return OperatorAction{}, err
+	}
+	return a, nil
+}
+
+// strictObject walks b's tokens (host/daemon strictRosterJSON's form): one JSON object (D4), no key
+// twice in any object (D5) and nothing after it (D6). encoding/json alone keeps the last of a repeated
+// key and stops reading after the first value.
+func strictObject(b []byte) error {
+	dec := json.NewDecoder(bytes.NewReader(b))
+	var walk func(tok json.Token) error
+	walk = func(tok json.Token) error {
+		d, ok := tok.(json.Delim)
+		if !ok {
+			return nil
+		}
+		seen := map[string]bool{}
+		for dec.More() {
+			if d == '{' {
+				key, err := dec.Token()
+				if err != nil {
+					return fmt.Errorf("%w: %v", ErrDecode, err)
+				}
+				k, _ := key.(string)
+				if seen[k] { // D5
+					return fmt.Errorf("%w: %q", ErrOperatorDuplicateKey, k)
+				}
+				seen[k] = true
+			}
+			next, err := dec.Token()
+			if err != nil {
+				return fmt.Errorf("%w: %v", ErrDecode, err)
+			}
+			if err := walk(next); err != nil {
+				return err
+			}
+		}
+		if _, err := dec.Token(); err != nil { // the closing delimiter
+			return fmt.Errorf("%w: %v", ErrDecode, err)
+		}
+		return nil
+	}
+	if tok, err := dec.Token(); err != nil || tok != json.Delim('{') { // D4
+		return ErrOperatorShape
+	} else if err := walk(tok); err != nil {
+		return err
+	}
+	if _, err := dec.Token(); err != io.EOF { // D6
+		return ErrOperatorTrailing
+	}
+	return nil
+}
+
+// OperatorReviewDigest is the SHA-256 of a's canonical bytes with its id blank: the semantic inputs
+// before an action id or challenge is assigned (design §3.1), the key S3 deduplicates intents on. It
+// validates every field but the id.
+func OperatorReviewDigest(a OperatorAction) ([32]byte, error) {
+	a.ID = ""
+	if err := a.validate(false); err != nil {
+		return [32]byte{}, err
+	}
+	c, err := CanonicalOperatorAction(a)
+	if err != nil {
+		return [32]byte{}, err
+	}
+	return sha256.Sum256(c), nil
+}
+
+// OperatorSignatureDigest is the SHA-256 of the final envelope's canonical bytes, id included: the
+// digest OperatorActionMessage signs. It validates a.
+func OperatorSignatureDigest(a OperatorAction) ([32]byte, error) {
+	if err := a.Validate(); err != nil {
+		return [32]byte{}, err
+	}
+	c, err := CanonicalOperatorAction(a)
+	if err != nil {
+		return [32]byte{}, err
+	}
+	return sha256.Sum256(c), nil
+}
+
+// OperatorActionMessage is the 97-byte message a confirmation signs (design §3.1):
+// OperatorActionPrefix ‖ 0x00 ‖ signature digest ‖ challenge ‖ expiry (uint64, big-endian).
+func OperatorActionMessage(digest, challenge [32]byte, expiry int64) ([]byte, error) {
+	if expiry < 0 {
+		return nil, ErrNegativeExpiry
+	}
+	m := make([]byte, 0, OperatorActionMessageLen)
+	m = append(m, OperatorActionPrefix...)
+	m = append(m, 0)
+	m = append(m, digest[:]...)
+	m = append(m, challenge[:]...)
+	m = binary.BigEndian.AppendUint64(m, uint64(expiry))
+	return m, nil
+}
+
+// OperatorActionMessageHash is the SHA-256 of OperatorActionMessage: a passkey's WebAuthn challenge.
+func OperatorActionMessageHash(digest, challenge [32]byte, expiry int64) ([32]byte, error) {
+	m, err := OperatorActionMessage(digest, challenge, expiry)
+	if err != nil {
+		return [32]byte{}, err
+	}
+	return sha256.Sum256(m), nil
+}
diff --git a/host/approval/operator_action_test.go b/host/approval/operator_action_test.go
new file mode 100644
index 00000000..d110a939
--- /dev/null
+++ b/host/approval/operator_action_test.go
@@ -0,0 +1,179 @@
+package approval_test
+
+// Row 243 V4 sprint 2 M1 (design §3.1, §5.1 TestOperatorStrictJSON): the canonical bytes and both digests,
+// pinned by a golden the real encoder writes (WORLD_UPDATE_OPERATOR_ACTION_GOLDEN=1 go test ./host/approval/
+// -run '^TestOperatorActionGolden$'); every invalid input class, each the golden's bytes with one change,
+// refused before canonicalization by its own error.
+
+import (
+	"bytes"
+	"encoding/binary"
+	"encoding/hex"
+	"encoding/json"
+	"errors"
+	"os"
+	"path/filepath"
+	"reflect"
+	"strings"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+)
+
+// opAction is a valid action of kind start or finish.
+func opAction(kind string) approval.OperatorAction {
+	a := approval.OperatorAction{Schema: approval.OperatorActionSchema, ID: "oa-0123456789abcdef0123456789abcdef", Kind: kind,
+		Episode: "fixture-ep", Store: "/ws/state/credentials.db", Options: approval.OperatorOptions{Preset: "sol-headless"},
+		Resolution: []approval.OperatorResolved{{Name: "agent", Value: "codex"}, {Name: "base_oid", Value: "8db22168689f31e19349b075e0199c8d9c940e61"},
+			{Name: "model", Value: "gpt-6.2-sol<&>"}, {Name: "repo", Value: "/src/ailang-world"}},
+		Task: "Fix the flaky test; $(rm -rf ~) `x` <b>é</b>\n", Review: "Start fixture-ep with sol-headless, headless.\n"}
+	if kind == "start" {
+		a.Options.Headless = true
+	} else {
+		a.Options.ConfirmOnSurface, a.Task, a.Review = true, "", "Finish fixture-ep: check, commit, then confirm the PR.\n"
+	}
+	a.Argv = approval.OperatorArgv(a.Kind, a.Episode, a.Options)
+	return a
+}
+
+var opChallenge = [32]byte{0x6f, 0x61, 1, 2, 3}
+
+const opExpiry = int64(1767225600)
+
+type opGolden struct {
+	Canonical, ReviewDigest, SignatureDigest, Message, MessageHash string
+}
+
+func opGoldenOf(t *testing.T, a approval.OperatorAction) opGolden {
+	t.Helper()
+	c, err := approval.CanonicalOperatorAction(a)
+	rd, err2 := approval.OperatorReviewDigest(a)
+	sd, err3 := approval.OperatorSignatureDigest(a)
+	m, err4 := approval.OperatorActionMessage(sd, opChallenge, opExpiry)
+	h, err5 := approval.OperatorActionMessageHash(sd, opChallenge, opExpiry)
+	if err := errors.Join(err, err2, err3, err4, err5); err != nil {
+		t.Fatal(err)
+	}
+	return opGolden{string(c), hex.EncodeToString(rd[:]), hex.EncodeToString(sd[:]), hex.EncodeToString(m), hex.EncodeToString(h[:])}
+}
+
+// opGoldenBytes is the start action's canonical bytes as the golden file holds them.
+func opGoldenBytes(t *testing.T) []byte {
+	t.Helper()
+	b, err := os.ReadFile(filepath.Join("testdata", "operator-action.golden.json"))
+	if err != nil {
+		t.Fatal(err)
+	}
+	var g map[string]opGolden
+	if err := json.Unmarshal(b, &g); err != nil {
+		t.Fatal(err)
+	}
+	return []byte(g["start"].Canonical)
+}
+
+func TestOperatorActionGolden(t *testing.T) {
+	got := map[string]opGolden{"start": opGoldenOf(t, opAction("start")), "finish": opGoldenOf(t, opAction("finish"))}
+	file := filepath.Join("testdata", "operator-action.golden.json")
+	if os.Getenv("WORLD_UPDATE_OPERATOR_ACTION_GOLDEN") == "1" {
+		b, _ := json.MarshalIndent(got, "", " ")
+		if err := os.WriteFile(file, append(b, '\n'), 0o644); err != nil {
+			t.Fatal(err)
+		}
+	}
+	b, err := os.ReadFile(file)
+	if err != nil {
+		t.Fatal(err)
+	}
+	var want map[string]opGolden
+	if err := json.Unmarshal(b, &want); err != nil {
+		t.Fatal(err)
+	}
+	if !reflect.DeepEqual(got, want) {
+		t.Fatalf("the encoder's bytes or digests moved:\n got %+v\nwant %+v", got, want)
+	}
+	for kind, g := range want {
+		a, err := approval.DecodeOperatorAction([]byte(g.Canonical))
+		if err != nil || !reflect.DeepEqual(a, opAction(kind)) {
+			t.Errorf("%s: the golden does not decode to its action: %v", kind, err)
+		}
+	}
+}
+
+// The message's layout, field by field: prefix, NUL, digest, challenge, big-endian expiry.
+func TestOperatorActionMessage(t *testing.T) {
+	d := [32]byte{0xd1, 0xd2}
+	m, err := approval.OperatorActionMessage(d, opChallenge, opExpiry)
+	if err != nil || approval.OperatorActionMessageLen != 97 || len(m) != 97 {
+		t.Fatalf("len %d, %v; want 97", len(m), err)
+	}
+	if string(m[:24]) != "world-operator-action/v1" || m[24] != 0 || !bytes.Equal(m[25:57], d[:]) || !bytes.Equal(m[57:89], opChallenge[:]) ||
+		binary.BigEndian.Uint64(m[89:]) != uint64(opExpiry) {
+		t.Fatalf("layout: %x", m)
+	}
+	if _, err := approval.OperatorActionMessage(d, opChallenge, -1); !errors.Is(err, approval.ErrNegativeExpiry) {
+		t.Fatalf("negative expiry: %v", err)
+	}
+	if _, err := approval.OperatorActionMessageHash(d, opChallenge, -1); !errors.Is(err, approval.ErrNegativeExpiry) {
+		t.Fatalf("negative expiry hash: %v", err)
+	}
+}
+
+// opDecodeRefused decodes b and wants err class want and the zero action.
+func opDecodeRefused(t *testing.T, name string, b []byte, want error) {
+	t.Helper()
+	a, err := approval.DecodeOperatorAction(b)
+	if !errors.Is(err, want) || !reflect.DeepEqual(a, approval.OperatorAction{}) {
+		t.Errorf("%s: %v (%+v); want %v and no action", name, err, a.ID, want)
+	}
+}
+
+func TestOperatorStrictJSON(t *testing.T) {
+	g := opGoldenBytes(t)
+	if _, err := approval.DecodeOperatorAction(g); err != nil {
+		t.Fatalf("control: %v", err)
+	}
+	s := string(g)
+	rep := func(old, new string) []byte {
+		if strings.Count(s, old) != 1 {
+			t.Fatalf("fixture: %q not once in %s", old, s)
+		}
+		return []byte(strings.Replace(s, old, new, 1))
+	}
+	for _, tc := range []struct {
+		name string
+		b    []byte
+		want error
+	}{
+		{"duplicate top-level key", rep(`"kind":"start",`, `"kind":"start","kind":"start",`), approval.ErrOperatorDuplicateKey},
+		{"duplicate nested key", rep(`"headless":true,`, `"headless":true,"headless":true,`), approval.ErrOperatorDuplicateKey},
+		{"duplicate key in an array element", rep(`{"name":"agent",`, `{"name":"agent","name":"agent",`), approval.ErrOperatorDuplicateKey},
+		{"unknown top-level key", rep(`"kind":"start",`, `"kind":"start","shell":"x",`), approval.ErrOperatorUnknownKey},
+		{"unknown nested key", rep(`"headless":true,`, `"headless":true,"env":"x",`), approval.ErrOperatorUnknownKey},
+		{"trailing object", append(bytes.Clone(g), `{}`...), approval.ErrOperatorTrailing},
+		{"trailing word", append(bytes.Clone(g), ` x`...), approval.ErrOperatorTrailing},
+		{"trailing newline", append(bytes.Clone(g), '\n'), approval.ErrNotCanonical},
+		{"invalid UTF-8", rep(`é`, "\xc3"), approval.ErrOperatorUTF8},
+		{"overlong UTF-8", rep(`é`, "\xc0\xa9"), approval.ErrOperatorUTF8},
+		{"raw NUL", rep(`<b>`, "<\x00b>"), approval.ErrOperatorNUL},
+		{"escaped NUL", rep(`<b>`, `<\u0000b>`), approval.ErrString},
+		{"null", []byte(`null`), approval.ErrOperatorShape},
+		{"array", []byte(`[` + s + `]`), approval.ErrOperatorShape},
+		{"string", []byte(`"x"`), approval.ErrOperatorShape},
+		{"number", []byte(`1`), approval.ErrOperatorShape},
+		{"empty", nil, approval.ErrOperatorShape},
+		{"unclosed", g[:len(g)-1], approval.ErrDecode},
+		{"argv a string", rep(`"argv":["ailang-worldd","episode","start","fixture-ep","--preset","sol-headless","--launch","--headless"]`, `"argv":"ailang-worldd"`), approval.ErrDecode},
+		{"argv an int element", rep(`"--headless"]`, `7]`), approval.ErrDecode},
+		{"options an array", rep(`"options":{"preset":"sol-headless","headless":true,"confirm_on_surface":false,"ready":false}`, `"options":[]`), approval.ErrDecode},
+		{"headless a string", rep(`"headless":true`, `"headless":"true"`), approval.ErrDecode},
+		{"argv null", rep(`"argv":["ailang-worldd","episode","start","fixture-ep","--preset","sol-headless","--launch","--headless"]`, `"argv":null`), approval.ErrOperatorAction},
+		{"id null", rep(`"id":"oa-0123456789abcdef0123456789abcdef"`, `"id":null`), approval.ErrNotCanonical},
+		{"case-folded key", rep(`"kind":`, `"Kind":`), approval.ErrNotCanonical},
+		{"reordered keys", rep(`"schema":"world/operator-action/v1","id":"oa-0123456789abcdef0123456789abcdef"`, `"id":"oa-0123456789abcdef0123456789abcdef","schema":"world/operator-action/v1"`), approval.ErrNotCanonical},
+		{"spaced", rep(`"kind":"start"`, `"kind": "start"`), approval.ErrNotCanonical},
+		{"HTML-escaped", rep(`<b>`, `\u003cb>`), approval.ErrNotCanonical},
+		{"canonical but invalid", rep(`"kind":"start"`, `"kind":"launch"`), approval.ErrOperatorAction},
+	} {
+		opDecodeRefused(t, tc.name, tc.b, tc.want)
+	}
+}
diff --git a/host/approval/testdata/operator-action.golden.json b/host/approval/testdata/operator-action.golden.json
new file mode 100644
index 00000000..f76b8930
--- /dev/null
+++ b/host/approval/testdata/operator-action.golden.json
@@ -0,0 +1,16 @@
+{
+ "finish": {
+  "Canonical": "{\"schema\":\"world/operator-action/v1\",\"id\":\"oa-0123456789abcdef0123456789abcdef\",\"kind\":\"finish\",\"episode\":\"fixture-ep\",\"store\":\"/ws/state/credentials.db\",\"argv\":[\"ailang-worldd\",\"episode\",\"finish\",\"fixture-ep\",\"--preset\",\"sol-headless\",\"--confirm-on-surface\"],\"options\":{\"preset\":\"sol-headless\",\"headless\":false,\"confirm_on_surface\":true,\"ready\":false},\"resolution\":[{\"name\":\"agent\",\"value\":\"codex\"},{\"name\":\"base_oid\",\"value\":\"8db22168689f31e19349b075e0199c8d9c940e61\"},{\"name\":\"model\",\"value\":\"gpt-6.2-sol\u003c\u0026\u003e\"},{\"name\":\"repo\",\"value\":\"/src/ailang-world\"}],\"task\":\"\",\"review\":\"Finish fixture-ep: check, commit, then confirm the PR.\\n\"}",
+  "ReviewDigest": "b9c16e206af876c6f553aed3c7116479993f42034a90dbeb0b906d9a739a4660",
+  "SignatureDigest": "bd143eecafe4c902ac54c4a1f593557aef84dc61cfe8cafdc45283e1af150f82",
+  "Message": "776f726c642d6f70657261746f722d616374696f6e2f763100bd143eecafe4c902ac54c4a1f593557aef84dc61cfe8cafdc45283e1af150f826f61010203000000000000000000000000000000000000000000000000000000000000006955b900",
+  "MessageHash": "1bc8f82c4b6941035c9cd49f2335947815fede831ee2337e57ebe10db58b6d58"
+ },
+ "start": {
+  "Canonical": "{\"schema\":\"world/operator-action/v1\",\"id\":\"oa-0123456789abcdef0123456789abcdef\",\"kind\":\"start\",\"episode\":\"fixture-ep\",\"store\":\"/ws/state/credentials.db\",\"argv\":[\"ailang-worldd\",\"episode\",\"start\",\"fixture-ep\",\"--preset\",\"sol-headless\",\"--launch\",\"--headless\"],\"options\":{\"preset\":\"sol-headless\",\"headless\":true,\"confirm_on_surface\":false,\"ready\":false},\"resolution\":[{\"name\":\"agent\",\"value\":\"codex\"},{\"name\":\"base_oid\",\"value\":\"8db22168689f31e19349b075e0199c8d9c940e61\"},{\"name\":\"model\",\"value\":\"gpt-6.2-sol\u003c\u0026\u003e\"},{\"name\":\"repo\",\"value\":\"/src/ailang-world\"}],\"task\":\"Fix the flaky test; $(rm -rf ~) `x` \u003cb\u003eé\u003c/b\u003e\\n\",\"review\":\"Start fixture-ep with sol-headless, headless.\\n\"}",
+  "ReviewDigest": "bb4a1f15e74d078556c2e366145885110652daf63088f0b7e6dd2da0d2c4eeb0",
+  "SignatureDigest": "f2b9bab25f0af92eee921f92a0cb3b7cf8425f0d855daeea8f1ffc3bcb219b40",
+  "Message": "776f726c642d6f70657261746f722d616374696f6e2f763100f2b9bab25f0af92eee921f92a0cb3b7cf8425f0d855daeea8f1ffc3bcb219b406f61010203000000000000000000000000000000000000000000000000000000000000006955b900",
+  "MessageHash": "9cab2a653aba7a7dd10981c335718ab449199227c57e6dcc597eb7cf7921028d"
+ }
+}
````

### Patch M2 — limits, validation, digests, the bar join

````diff M2
diff --git a/host/approval/operator_action_valid_test.go b/host/approval/operator_action_valid_test.go
new file mode 100644
index 00000000..5170f46b
--- /dev/null
+++ b/host/approval/operator_action_valid_test.go
@@ -0,0 +1,243 @@
+package approval_test
+
+// Row 243 V4 sprint 2 M2 (design §3.1, §5.1 the envelope half of TestOperatorRequestLimits): the
+// 128 KiB and 8 KiB caps at the exact boundary (one under, at, one over), every Validate rule refused
+// on its own from a valid action, and the two digests: the review digest covers the semantic inputs
+// and not the id, the signature digest covers both.
+
+import (
+	"bytes"
+	"errors"
+	"fmt"
+	"slices"
+	"strings"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/verbs"
+)
+
+// The envelope canonicalizes host/verbs' admitted tuple: OperatorArgv is the bar's Intent.Argv for every
+// admitted line, an action built from the Intent validates, and the two action-id grammars agree.
+func TestOperatorArgvIsTheBarsIntent(t *testing.T) {
+	yes := verbs.OperatorPolicy{Start: "A", Finish: "A", HandoffProven: true}
+	for _, line := range []string{"episode start e1 --preset p --launch --headless", "episode start e1 --launch=true --preset=p",
+		"episode finish e1 --preset p", "episode finish e1 --ready --preset p --confirm-on-surface", "episode finish e1 --preset p --ready"} {
+		w, _ := verbs.Split(line)
+		in := verbs.AdmitOperator(w, yes).Intent
+		a := opAction(in.Kind)
+		a.Episode, a.Options = in.Episode, approval.OperatorOptions{Preset: in.Preset, Headless: in.Headless, ConfirmOnSurface: in.ConfirmOnSurface, Ready: in.Ready}
+		if a.Argv = in.Argv; !slices.Equal(approval.OperatorArgv(a.Kind, a.Episode, a.Options), in.Argv) || a.Validate() != nil {
+			t.Errorf("%s: %q, %v", line, approval.OperatorArgv(a.Kind, a.Episode, a.Options), a.Validate())
+		}
+	}
+	for _, id := range []string{"oa-0123456789abcdef0123456789abcdef", "oa-0123456789abcdef0123456789abcde", "oa-0123456789abcdef0123456789abcdef0",
+		"OA-0123456789abcdef0123456789abcdef", "oa-0123456789ABCDEF0123456789abcdef", "xoa-0123456789abcdef0123456789abcde", "oa-0123456789abcdef0123456789abcdef\n"} {
+		a := opAction("start")
+		a.ID = id
+		if (a.Validate() == nil) != verbs.ValidOperatorActionID(id) {
+			t.Errorf("id %q: the envelope and the bar disagree", id)
+		}
+	}
+}
+
+// opSized is the start action with its review padded so its canonical bytes are exactly n long.
+func opSized(t *testing.T, n int) approval.OperatorAction {
+	t.Helper()
+	a := opAction("start")
+	a.Review = "r"
+	c, _ := approval.CanonicalOperatorAction(a)
+	a.Review = strings.Repeat("r", n-len(c)+1)
+	if c, _ = approval.CanonicalOperatorAction(a); len(c) != n {
+		t.Fatalf("sized %d, want %d", len(c), n)
+	}
+	return a
+}
+
+func TestOperatorRequestLimits(t *testing.T) {
+	max := approval.MaxOperatorActionBytes
+	if max != 131072 || approval.MaxOperatorTaskBytes != 8192 {
+		t.Fatalf("caps %d %d", max, approval.MaxOperatorTaskBytes)
+	}
+	at := opSized(t, max)
+	c, _ := approval.CanonicalOperatorAction(at)
+	if err := at.Validate(); err != nil {
+		t.Fatalf("exactly 128 KiB: %v", err)
+	}
+	if _, err := approval.DecodeOperatorAction(c); err != nil {
+		t.Fatalf("decode exactly 128 KiB: %v", err)
+	}
+	for _, n := range []int{max - 1, max} {
+		if err := opSized(t, n).Validate(); err != nil {
+			t.Errorf("%d bytes: %v", n, err)
+		}
+	}
+	if err := opSized(t, max+1).Validate(); !errors.Is(err, approval.ErrOperatorTooLarge) {
+		t.Errorf("one over: %v", err)
+	}
+	over, _ := approval.CanonicalOperatorAction(opSized(t, max+1))
+	opDecodeRefused(t, "decode one over", over, approval.ErrOperatorTooLarge)
+	opDecodeRefused(t, "decode one over by a space, before any parse", append(bytes.Clone(c), ' '), approval.ErrOperatorTooLarge)
+	under, _ := approval.CanonicalOperatorAction(opSized(t, max-1))
+	opDecodeRefused(t, "at the cap with a space", append(under, ' '), approval.ErrNotCanonical)
+	for _, tc := range []struct {
+		name, kind, task string
+		ok               bool
+	}{
+		{"8 KiB of two-byte runes", "start", strings.Repeat("é", 4096), true},
+		{"one byte", "start", "x", true},
+		{"8 KiB + 1", "start", strings.Repeat("é", 4096) + "x", false},
+		{"8 KiB - 1 + a two-byte rune", "start", strings.Repeat("x", 8191) + "é", false},
+		{"empty start", "start", "", false},
+		{"a finish's task", "finish", "x", false},
+		{"no finish task", "finish", "", true},
+	} {
+		a := opAction(tc.kind)
+		a.Task = tc.task
+		err := a.Validate()
+		if tc.ok != (err == nil) || (!tc.ok && !errors.Is(err, approval.ErrOperatorTask)) {
+			t.Errorf("task %s: %v", tc.name, err)
+		}
+	}
+}
+
+func TestOperatorActionValidate(t *testing.T) {
+	for _, kind := range []string{"start", "finish"} {
+		if err := opAction(kind).Validate(); err != nil {
+			t.Fatalf("control %s: %v", kind, err)
+		}
+	}
+	res := func(names ...string) []approval.OperatorResolved {
+		r := []approval.OperatorResolved{}
+		for _, n := range names {
+			r = append(r, approval.OperatorResolved{Name: n, Value: "v"})
+		}
+		return r
+	}
+	many := []string{}
+	for i := 0; i < 65; i++ {
+		many = append(many, fmt.Sprintf("n%02d", i))
+	}
+	argv := func(a *approval.OperatorAction) { a.Argv = approval.OperatorArgv(a.Kind, a.Episode, a.Options) }
+	for _, tc := range []struct {
+		name, kind string
+		edit       func(*approval.OperatorAction)
+		want       error
+	}{
+		{"schema", "start", func(a *approval.OperatorAction) { a.Schema = approval.VerbRequestSchema }, approval.ErrSchema},
+		{"NUL in a leaf", "start", func(a *approval.OperatorAction) { a.Resolution[0].Value = "a\x00" }, approval.ErrString},
+		{"invalid UTF-8 in a leaf", "start", func(a *approval.OperatorAction) { a.Review = "\xff" }, approval.ErrString},
+		{"no id", "start", func(a *approval.OperatorAction) { a.ID = "" }, approval.ErrOperatorAction},
+		{"id with a prefix", "start", func(a *approval.OperatorAction) { a.ID = "x" + a.ID }, approval.ErrOperatorAction},
+		{"id with a suffix", "start", func(a *approval.OperatorAction) { a.ID += "0" }, approval.ErrOperatorAction},
+		{"id with a newline", "start", func(a *approval.OperatorAction) { a.ID += "\n" }, approval.ErrOperatorAction},
+		{"id upper-case", "start", func(a *approval.OperatorAction) { a.ID = strings.ToUpper(a.ID) }, approval.ErrOperatorAction},
+		{"id 31 hex", "start", func(a *approval.OperatorAction) { a.ID = a.ID[:34] }, approval.ErrOperatorAction},
+		{"id mr-", "start", func(a *approval.OperatorAction) { a.ID = "mr" + a.ID[2:] }, approval.ErrOperatorAction},
+		{"kind launch", "start", func(a *approval.OperatorAction) { a.Kind = "launch"; argv(a) }, approval.ErrOperatorAction},
+		{"kind Start", "start", func(a *approval.OperatorAction) { a.Kind = "Start"; argv(a) }, approval.ErrOperatorAction},
+		{"episode", "start", func(a *approval.OperatorAction) { a.Episode = "-ep"; argv(a) }, approval.ErrEpisode},
+		{"store relative", "start", func(a *approval.OperatorAction) { a.Store = "state/credentials.db" }, approval.ErrPath},
+		{"store unclean", "start", func(a *approval.OperatorAction) { a.Store = "/ws/x/../credentials.db" }, approval.ErrPath},
+		{"store empty", "start", func(a *approval.OperatorAction) { a.Store = "" }, approval.ErrPath},
+		{"preset empty", "start", func(a *approval.OperatorAction) { a.Options.Preset = ""; argv(a) }, approval.ErrOperatorAction},
+		{"preset se-tools", "start", func(a *approval.OperatorAction) { a.Options.Preset = "se-tools"; argv(a) }, approval.ErrOperatorAction},
+		{"preset 52", "start", func(a *approval.OperatorAction) { a.Options.Preset = strings.Repeat("p", 52); argv(a) }, approval.ErrOperatorAction},
+		{"preset upper-case", "start", func(a *approval.OperatorAction) { a.Options.Preset = "Sol"; argv(a) }, approval.ErrOperatorAction},
+		{"start ready", "start", func(a *approval.OperatorAction) { a.Options.Ready = true; argv(a) }, approval.ErrOperatorAction},
+		{"start confirm-on-surface", "start", func(a *approval.OperatorAction) { a.Options.ConfirmOnSurface = true; argv(a) }, approval.ErrOperatorAction},
+		{"finish headless", "finish", func(a *approval.OperatorAction) { a.Options.Headless = true; argv(a) }, approval.ErrOperatorAction},
+		{"argv a word dropped", "start", func(a *approval.OperatorAction) { a.Argv = a.Argv[:7] }, approval.ErrOperatorAction},
+		{"argv a word added", "start", func(a *approval.OperatorAction) { a.Argv = append(a.Argv, "--grant") }, approval.ErrOperatorAction},
+		{"argv reordered", "start", func(a *approval.OperatorAction) { a.Argv[6], a.Argv[7] = a.Argv[7], a.Argv[6] }, approval.ErrOperatorAction},
+		{"argv another episode", "start", func(a *approval.OperatorAction) { a.Argv[3] = "other-ep" }, approval.ErrOperatorAction},
+		{"argv without --headless", "start", func(a *approval.OperatorAction) { a.Options.Headless = false }, approval.ErrOperatorAction},
+		{"argv without --confirm-on-surface", "finish", func(a *approval.OperatorAction) { a.Options.ConfirmOnSurface = false }, approval.ErrOperatorAction},
+		{"argv with --ready", "finish", func(a *approval.OperatorAction) { a.Argv = append(a.Argv, "--ready") }, approval.ErrOperatorAction},
+		{"resolution nil", "start", func(a *approval.OperatorAction) { a.Resolution = nil }, approval.ErrOperatorAction},
+		{"resolution empty", "start", func(a *approval.OperatorAction) { a.Resolution = res() }, approval.ErrOperatorAction},
+		{"resolution 65", "start", func(a *approval.OperatorAction) { a.Resolution = res(many...) }, approval.ErrOperatorAction},
+		{"resolution unsorted", "start", func(a *approval.OperatorAction) { a.Resolution = res("b", "a") }, approval.ErrOperatorAction},
+		{"resolution repeated", "start", func(a *approval.OperatorAction) { a.Resolution = res("a", "a") }, approval.ErrOperatorAction},
+		{"resolution name", "start", func(a *approval.OperatorAction) { a.Resolution = res("Agent") }, approval.ErrOperatorAction},
+		{"review empty", "start", func(a *approval.OperatorAction) { a.Review = "" }, approval.ErrOperatorAction},
+	} {
+		a := opAction(tc.kind)
+		a.Resolution = append([]approval.OperatorResolved(nil), a.Resolution...)
+		tc.edit(&a)
+		if err := a.Validate(); !errors.Is(err, tc.want) {
+			t.Errorf("%s: %v, want %v", tc.name, err, tc.want)
+		}
+		if _, err := approval.OperatorSignatureDigest(a); !errors.Is(err, tc.want) {
+			t.Errorf("%s: the signature digest of an invalid action: %v", tc.name, err)
+		}
+	}
+	for _, ok := range [][]string{{"a"}, many[:64]} {
+		a := opAction("start")
+		a.Resolution = res(ok...)
+		if err := a.Validate(); err != nil {
+			t.Errorf("resolution of %d: %v", len(a.Resolution), err)
+		}
+	}
+	a := opAction("start")
+	a.Options.Preset = strings.Repeat("p", 51)
+	a.Argv = approval.OperatorArgv(a.Kind, a.Episode, a.Options)
+	if err := a.Validate(); err != nil {
+		t.Errorf("preset 51: %v", err)
+	}
+}
+
+// The review digest covers the semantic inputs and not the id; the signature digest covers both.
+func TestOperatorDigestsAreDistinct(t *testing.T) {
+	a := opAction("start")
+	rd, _ := approval.OperatorReviewDigest(a)
+	sd, _ := approval.OperatorSignatureDigest(a)
+	if rd == sd || rd == [32]byte{} || sd == [32]byte{} {
+		t.Fatalf("review %x signature %x", rd, sd)
+	}
+	other := a
+	other.ID = "oa-ffffffffffffffffffffffffffffffff"
+	rd2, _ := approval.OperatorReviewDigest(other)
+	sd2, _ := approval.OperatorSignatureDigest(other)
+	if rd2 != rd || sd2 == sd {
+		t.Fatalf("another id: review moved %v, signature kept %v", rd2 != rd, sd2 == sd)
+	}
+	noID := a
+	noID.ID = ""
+	if d, err := approval.OperatorReviewDigest(noID); err != nil || d != rd {
+		t.Fatalf("the review digest needs no id: %v", err)
+	}
+	if _, err := approval.OperatorSignatureDigest(noID); !errors.Is(err, approval.ErrOperatorAction) {
+		t.Fatalf("the signature digest needs an id: %v", err)
+	}
+	bad := a
+	bad.Kind = "launch"
+	if _, err := approval.OperatorReviewDigest(bad); !errors.Is(err, approval.ErrOperatorAction) {
+		t.Fatalf("the review digest of an invalid action: %v", err)
+	}
+	for name, edit := range opFieldEdits() {
+		b := opAction("start")
+		edit(&b)
+		r, err := approval.OperatorReviewDigest(b)
+		s, err2 := approval.OperatorSignatureDigest(b)
+		if err != nil || err2 != nil || r == rd || s == sd {
+			t.Errorf("%s: review moved %v, signature moved %v (%v, %v)", name, r != rd, s != sd, err, err2)
+		}
+	}
+}
+
+// opFieldEdits changes one semantic field of a valid start action, keeping it valid.
+func opFieldEdits() map[string]func(*approval.OperatorAction) {
+	re := func(a *approval.OperatorAction) { a.Argv = approval.OperatorArgv(a.Kind, a.Episode, a.Options) }
+	return map[string]func(*approval.OperatorAction){
+		"episode":  func(a *approval.OperatorAction) { a.Episode = "other-ep"; re(a) },
+		"store":    func(a *approval.OperatorAction) { a.Store = "/other/credentials.db" },
+		"preset":   func(a *approval.OperatorAction) { a.Options.Preset = "sol-two"; re(a) },
+		"headless": func(a *approval.OperatorAction) { a.Options.Headless = false; re(a) },
+		"resolution": func(a *approval.OperatorAction) {
+			a.Resolution = []approval.OperatorResolved{{Name: "agent", Value: "claude"}}
+		},
+		"task":   func(a *approval.OperatorAction) { a.Task += "!" },
+		"review": func(a *approval.OperatorAction) { a.Review += "!" },
+	}
+}
````

### Patch M3 — the key law

````diff M3
diff --git a/host/approval/operator_signed.go b/host/approval/operator_signed.go
new file mode 100644
index 00000000..42be6748
--- /dev/null
+++ b/host/approval/operator_signed.go
@@ -0,0 +1,164 @@
+package approval
+
+// Row 243 V4 sprint 2 M3 (design §3.1, §5.2): the pure law of an operator-action confirmation, verb_signed.go's
+// shape over the operator message. ConfirmOperatorAction is the daemon's confirmation; CheckOperatorAction
+// is the runner's recheck from its own frozen action (no P1, the counter already moved, the key judged at
+// now). No I/O; now is an input.
+
+import (
+	"bytes"
+	"errors"
+	"fmt"
+	"time"
+)
+
+// SignedOperatorAction is a verified operator-action confirmation. Its fields are unexported, so only
+// ConfirmOperatorAction and CheckOperatorAction return one with ok set; the zero value binds nothing.
+type SignedOperatorAction struct {
+	ok                bool
+	digest, challenge [32]byte
+	expiry            int64
+	name, format      string
+	credentialID      []byte
+	signCount         uint32
+}
+
+// OperatorConfirmInput is what ConfirmOperatorAction judges: the stored action (strictly decoded), the
+// row's admission input (StoredDigest is the row's digest column, BytesDigest the SHA-256 of its stored
+// bytes), its challenge, the artifact, its format and the enrolled approver it names.
+type OperatorConfirmInput struct {
+	Action    OperatorAction
+	Admit     AdmitInput
+	Challenge [32]byte
+	Artifact  []byte
+	Format    string // FormatSSH or FormatWebAuthn
+	Approver  Approver
+}
+
+// OperatorCheckInput is what CheckOperatorAction rechecks: the runner's own frozen action, the
+// challenge and expiry it recorded, the clock, the stored artifact, its format and the approver.
+type OperatorCheckInput struct {
+	Action      OperatorAction
+	Challenge   [32]byte
+	Expiry, Now int64
+	Artifact    []byte
+	Format      string
+	Approver    Approver
+}
+
+var (
+	ErrOperatorLifetime = errors.New("approval: the action's request outlives 600 s")
+	ErrOperatorCounter  = errors.New("approval: the confirmation's passkey counter was never recorded")
+)
+
+// ConfirmOperatorAction is the operator confirmation law, in this order: the action is valid (its
+// signature digest), P1 admission (consumed, clock behind, expiry at or before now, stored bytes), the
+// row's digest is this action's, the lifetime is at most OperatorActionLifetime, P2-P5, then the
+// signature over the operator message, and for a passkey W1-W14 and W15.
+func ConfirmOperatorAction(in OperatorConfirmInput) (*SignedOperatorAction, error) {
+	d, err := OperatorSignatureDigest(in.Action)
+	if err != nil {
+		return nil, err
+	}
+	if v := Admit(in.Admit); v != Admitted { // P1
+		return nil, &NotAdmittedError{Verdict: v}
+	}
+	if in.Admit.StoredDigest != d { // C1
+		return nil, &BindingError{Field: "digest"}
+	}
+	if in.Admit.ExpiresAt-in.Admit.CreatedAt > int64(OperatorActionLifetime/time.Second) { // C2
+		return nil, ErrOperatorLifetime
+	}
+	return operatorProof(d, in.Challenge, in.Admit.ExpiresAt, in.Admit.Now, in.Artifact, in.Format, in.Approver, true)
+}
+
+// CheckOperatorAction is the runner's recheck before an effect (design §3.1): its own action's digest,
+// the recorded challenge and expiry, expiry at or before now refused (R1), P2-P5 at now (a key revoked
+// since is refused), the signature, and for a passkey W1-W14 and a counter the confirmation recorded
+// (R2: the stored counter is at or past the artifact's; never advanced again).
+func CheckOperatorAction(in OperatorCheckInput) (*SignedOperatorAction, error) {
+	d, err := OperatorSignatureDigest(in.Action)
+	if err != nil {
+		return nil, err
+	}
+	if in.Now >= in.Expiry { // R1
+		return nil, &NotAdmittedError{Verdict: RefuseExpired}
+	}
+	return operatorProof(d, in.Challenge, in.Expiry, in.Now, in.Artifact, in.Format, in.Approver, false)
+}
+
+// operatorProof is P2-P6 / W1-W15 for both laws; confirm selects W15 (the counter must move) over R2.
+func operatorProof(d, challenge [32]byte, expiry, now int64, artifact []byte, format string, ap Approver, confirm bool) (*SignedOperatorAction, error) {
+	if format != FormatSSH && format != FormatWebAuthn {
+		return nil, fmt.Errorf("%w: %q", ErrVerbFormat, format)
+	}
+	if err := checkApprover(AdmitInput{Now: now}, ap, format); err != nil { // P2-P5, P3w
+		return nil, err
+	}
+	s := &SignedOperatorAction{ok: true, digest: d, challenge: challenge, expiry: expiry, name: ap.Name, format: format}
+	if format == FormatSSH {
+		sig, err := ParseSSHSIG(artifact)
+		if err != nil {
+			return nil, err
+		}
+		msg, err := OperatorActionMessage(d, challenge, expiry)
+		if err != nil {
+			return nil, err
+		}
+		if err := VerifySSHSIG(sig, OperatorActionNamespace, msg, ap.Key); err != nil { // P6
+			return nil, err
+		}
+		return s, nil
+	}
+	h, hErr := OperatorActionMessageHash(d, challenge, expiry)
+	ad, err := assertionSteps(artifact, ap, h, hErr) // W1-W14
+	if err != nil {
+		return nil, err
+	}
+	p := ap.Passkey
+	if confirm && (ad.SignCount != 0 || p.SignCount != 0) && ad.SignCount <= p.SignCount { // W15
+		return nil, fmt.Errorf("%w: enrolled %d, this assertion %d; if you did not cause this, run: ailang-worldd approver revoke %s",
+			ErrSignCountRegression, p.SignCount, ad.SignCount, ap.Name)
+	}
+	if !confirm && ad.SignCount > p.SignCount { // R2
+		return nil, fmt.Errorf("%w: enrolled %d, this assertion %d", ErrOperatorCounter, p.SignCount, ad.SignCount)
+	}
+	s.credentialID, s.signCount = bytes.Clone(p.CredentialID), ad.SignCount
+	return s, nil
+}
+
+// Binds reports whether s is a verified confirmation of exactly this action, challenge and expiry.
+func (s *SignedOperatorAction) Binds(a OperatorAction, challenge [32]byte, expiry int64) error {
+	if s == nil || !s.ok {
+		return ErrNotVerified
+	}
+	d, err := OperatorSignatureDigest(a)
+	switch {
+	case err != nil:
+		return err
+	case s.digest != d:
+		return &BindingError{Field: "digest"}
+	case s.challenge != challenge:
+		return &BindingError{Field: "challenge"}
+	case s.expiry != expiry:
+		return &BindingError{Field: "expiry"}
+	}
+	return nil
+}
+
+// ApprovedBy is format ":" name ("webauthn:mark-mac"); "" unless verified.
+func (s *SignedOperatorAction) ApprovedBy() string {
+	if s == nil || !s.ok {
+		return ""
+	}
+	return s.format + ":" + s.name
+}
+
+// PasskeyUse is what S3's confirm CAS binds for a passkey: the credential ID (a copy) and the
+// assertion's counter; ok only for a verified webauthn confirmation.
+func (s *SignedOperatorAction) PasskeyUse() (credentialID []byte, signCount uint32, ok bool) {
+	if s == nil || !s.ok || s.format != FormatWebAuthn {
+		return nil, 0, false
+	}
+	return bytes.Clone(s.credentialID), s.signCount, true
+}
diff --git a/host/approval/operator_signed_test.go b/host/approval/operator_signed_test.go
new file mode 100644
index 00000000..675e85dc
--- /dev/null
+++ b/host/approval/operator_signed_test.go
@@ -0,0 +1,367 @@
+package approval_test
+
+// Row 243 V4 sprint 2 M3 (design §3.1, §5.2 ConfirmBothKeys, ConfirmProof, ConfirmAuthenticity,
+// ConfirmRevocationCounter, ConfirmBinding, ConfirmExpiry, and the pure half of RunnerRechecks): the
+// operator confirmation law with approvaltest keys and passkeys (D-WORLD-109), every check refused on
+// its own with its own error and no proof, against a valid control in the same test.
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
+var (
+	opSigner  = approvaltest.NewSigner(81)
+	opPasskey = approvaltest.NewPasskey(approval.AlgEdDSA, 82)
+	opC       = [32]byte{0x6f, 0x63}
+)
+
+const opCreated, opExp, opNow = int64(1000), int64(1600), int64(1010)
+
+func opDigest(t *testing.T, a approval.OperatorAction) [32]byte {
+	t.Helper()
+	d, err := approval.OperatorSignatureDigest(a)
+	if err != nil {
+		t.Fatal(err)
+	}
+	return d
+}
+
+// opSSH and opAssert sign action a's message (challenge opC, expiry exp).
+func opSSH(t *testing.T, s *approvaltest.Signer, ns string, a approval.OperatorAction, exp int64) []byte {
+	t.Helper()
+	m, err := approval.OperatorActionMessage(opDigest(t, a), opC, exp)
+	if err != nil {
+		t.Fatal(err)
+	}
+	return s.Sign(ns, m)
+}
+
+func opAssert(t *testing.T, a approval.OperatorAction, exp int64, count uint32, e ...approvaltest.Edit) []byte {
+	t.Helper()
+	h, err := approval.OperatorActionMessageHash(opDigest(t, a), opC, exp)
+	if err != nil {
+		t.Fatal(err)
+	}
+	return opPasskey.Assert(paOrigin, paRPID, h, count, e...)
+}
+
+// opIn is a valid confirmation of the start action: mark-ssh (opSigner) or mark-mac (opPasskey, stored
+// counter 1, assertion 2), created 1000, expires 1600, now 1010.
+func opIn(t *testing.T, format string) approval.OperatorConfirmInput {
+	t.Helper()
+	a := opAction("start")
+	d := opDigest(t, a)
+	in := approval.OperatorConfirmInput{Action: a, Challenge: opC, Format: format,
+		Admit: approval.AdmitInput{Now: opNow, CreatedAt: opCreated, ExpiresAt: opExp, StoredDigest: d, BytesDigest: d}}
+	if format == approval.FormatSSH {
+		in.Artifact, in.Approver = opSSH(t, opSigner, approval.OperatorActionNamespace, a, opExp), approval.Approver{Name: "mark-ssh",
+			Format: approval.FormatSSH, Key: opSigner.PublicKey(), EnrolledAt: 900}
+	} else {
+		in.Artifact, in.Approver = opAssert(t, a, opExp, 2), paApprover(t, opPasskey)
+	}
+	return in
+}
+
+// opRefused runs the law on in and wants err class want and no proof.
+func opRefused(t *testing.T, name string, in approval.OperatorConfirmInput, want error) {
+	t.Helper()
+	s, err := approval.ConfirmOperatorAction(in)
+	if s != nil || !errors.Is(err, want) {
+		t.Errorf("%s: %v, %v; want %v and no proof", name, s, err, want)
+	}
+}
+
+func TestOperatorConfirmBothKeys(t *testing.T) {
+	a := opAction("start")
+	for format, by := range map[string]string{approval.FormatSSH: "ssh:mark-ssh", approval.FormatWebAuthn: "webauthn:mark-mac"} {
+		s, err := approval.ConfirmOperatorAction(opIn(t, format))
+		if err != nil || s.ApprovedBy() != by || s.Binds(a, opC, opExp) != nil {
+			t.Fatalf("%s: %v %q", format, err, s.ApprovedBy())
+		}
+		cid, n, ok := s.PasskeyUse()
+		if (format == approval.FormatWebAuthn) != ok || (ok && (n != 2 || !bytes.Equal(cid, opPasskey.CredentialID()))) {
+			t.Fatalf("%s: passkey use %x %d %v", format, cid, n, ok)
+		}
+		var be *approval.BindingError
+		other := a
+		other.Task += "!"
+		for field, err := range map[string]error{"digest": s.Binds(other, opC, opExp), "challenge": s.Binds(a, [32]byte{1}, opExp),
+			"expiry": s.Binds(a, opC, opExp+1)} {
+			if !errors.As(err, &be) || be.Field != field {
+				t.Errorf("%s: Binds another %s: %v", format, field, err)
+			}
+		}
+		if err := s.Binds(approval.OperatorAction{}, opC, opExp); !errors.Is(err, approval.ErrSchema) {
+			t.Errorf("%s: Binds an invalid action: %v", format, err)
+		}
+	}
+	var zero *approval.SignedOperatorAction
+	for _, s := range []*approval.SignedOperatorAction{zero, {}} {
+		_, _, ok := s.PasskeyUse()
+		if !errors.Is(s.Binds(a, opC, opExp), approval.ErrNotVerified) || s.ApprovedBy() != "" || ok {
+			t.Fatal("an unverified proof binds something")
+		}
+	}
+}
+
+func TestOperatorConfirmProof(t *testing.T) {
+	ssh, pk := opIn(t, approval.FormatSSH), opIn(t, approval.FormatWebAuthn)
+	edit := func(in approval.OperatorConfirmInput, f func(*approval.OperatorConfirmInput)) approval.OperatorConfirmInput {
+		f(&in)
+		return in
+	}
+	opRefused(t, "typed y", edit(ssh, func(in *approval.OperatorConfirmInput) { in.Artifact = []byte("y\n") }), approval.ErrArmor)
+	opRefused(t, "no proof (ssh)", edit(ssh, func(in *approval.OperatorConfirmInput) { in.Artifact = nil }), approval.ErrArmor)
+	opRefused(t, "no proof (passkey)", edit(pk, func(in *approval.OperatorConfirmInput) { in.Artifact = nil }), approval.ErrPasskeyResponse)
+	opRefused(t, "malformed armor", edit(ssh, func(in *approval.OperatorConfirmInput) { in.Artifact = in.Artifact[:len(in.Artifact)/2] }), approval.ErrArmor)
+	opRefused(t, "malformed assertion", edit(pk, func(in *approval.OperatorConfirmInput) { in.Artifact = []byte("{}") }), approval.ErrPasskeyResponse)
+	opRefused(t, "an assertion as ssh", edit(ssh, func(in *approval.OperatorConfirmInput) { in.Artifact = pk.Artifact }), approval.ErrArmor)
+	opRefused(t, "an SSHSIG as passkey", edit(pk, func(in *approval.OperatorConfirmInput) { in.Artifact = ssh.Artifact }), approval.ErrPasskeyResponse)
+	opRefused(t, "no format", edit(ssh, func(in *approval.OperatorConfirmInput) { in.Format = "" }), approval.ErrVerbFormat)
+	opRefused(t, "typed format", edit(ssh, func(in *approval.OperatorConfirmInput) { in.Format = "typed" }), approval.ErrVerbFormat)
+	opRefused(t, "bad challenge (ssh)", edit(ssh, func(in *approval.OperatorConfirmInput) { in.Challenge[0] ^= 1 }), approval.ErrBadSignature)
+	opRefused(t, "bad challenge (passkey)", edit(pk, func(in *approval.OperatorConfirmInput) { in.Challenge[0] ^= 1 }), approval.ErrChallenge)
+	opRefused(t, "an invalid action", edit(ssh, func(in *approval.OperatorConfirmInput) { in.Action.ID = "" }), approval.ErrOperatorAction)
+}
+
+func TestOperatorConfirmAuthenticity(t *testing.T) {
+	a := opAction("start")
+	ssh, pk := opIn(t, approval.FormatSSH), opIn(t, approval.FormatWebAuthn)
+	for _, tc := range []struct {
+		name string
+		in   approval.OperatorConfirmInput
+		want error
+	}{
+		{"another key's SSHSIG", func() approval.OperatorConfirmInput {
+			in := ssh
+			in.Artifact = opSSH(t, approvaltest.NewSigner(83), approval.OperatorActionNamespace, a, opExp)
+			return in
+		}(), approval.ErrKeyMismatch},
+		{"forged assertion", func() approval.OperatorConfirmInput {
+			in := pk
+			in.Artifact = opAssert(t, a, opExp, 2, approvaltest.Edit{BadSignature: true})
+			return in
+		}(), approval.ErrBadSignature},
+		{"unknown credential", func() approval.OperatorConfirmInput {
+			in := pk
+			in.Artifact = opAssert(t, a, opExp, 2, approvaltest.Edit{RawID: []byte("another-credential")})
+			return in
+		}(), approval.ErrCredentialIDMismatch},
+	} {
+		opRefused(t, tc.name, tc.in, tc.want)
+	}
+	for _, ns := range []string{approval.VerbNamespace, approval.Namespace, approval.DecisionNamespace, "git", "world-operator-action@sunholo.x"} {
+		in := ssh
+		in.Artifact = opSSH(t, opSigner, ns, a, opExp)
+		opRefused(t, "namespace "+ns, in, approval.ErrNamespace)
+	}
+	h, _ := approval.OperatorActionMessageHash(opDigest(t, a), opC, opExp)
+	for _, tc := range []struct {
+		name, origin, rpID string
+		e                  approvaltest.Edit
+		want               error
+	}{
+		{"no UV", paOrigin, paRPID, approvaltest.Edit{Clear: 0x04}, approval.ErrUserVerification},
+		{"no UP", paOrigin, paRPID, approvaltest.Edit{Clear: 0x01}, approval.ErrUserPresence},
+		{"another RP id", paOrigin, "evil.ts.net", approvaltest.Edit{}, approval.ErrRPIDHash},
+		{"another origin", "https://x.ts.net:8444", paRPID, approvaltest.Edit{}, approval.ErrOriginMismatch},
+		{"cross-origin", paOrigin, paRPID, approvaltest.Edit{CrossOrigin: true}, approval.ErrCrossOrigin},
+		{"a registration", paOrigin, paRPID, approvaltest.Edit{Type: "webauthn.create"}, approval.ErrClientDataType},
+	} {
+		in := pk
+		in.Artifact = opPasskey.Assert(tc.origin, tc.rpID, h, 2, tc.e)
+		opRefused(t, tc.name, in, tc.want)
+	}
+}
+
+func TestOperatorConfirmRevocationCounter(t *testing.T) {
+	ssh, pk := opIn(t, approval.FormatSSH), opIn(t, approval.FormatWebAuthn)
+	a := opAction("start")
+	revoked := ssh
+	revoked.Approver.RevokedAt = 950
+	opRefused(t, "revoked ssh", revoked, approval.ErrApproverRevoked)
+	pr := opIn(t, approval.FormatWebAuthn)
+	pr.Approver.RevokedAt = 950
+	opRefused(t, "revoked passkey", pr, approval.ErrApproverRevoked)
+	late := ssh
+	late.Approver.EnrolledAt = opNow + 1
+	opRefused(t, "enrolled after now", late, approval.ErrApproverClockBehind)
+	if late.Approver.EnrolledAt = opCreated + 5; func() error { _, err := approval.ConfirmOperatorAction(late); return err }() != nil {
+		t.Error("enrolled after the request, before now: refused") // P5 at now, not at the row's creation
+	}
+	asSSH := pk
+	asSSH.Approver.Format = approval.FormatSSH
+	opRefused(t, "a passkey approver that is not a passkey", asSSH, approval.ErrApproverNotPasskey)
+	for _, c := range []struct {
+		stored, asserted uint32
+		ok               bool
+	}{{1, 1, false}, {5, 2, false}, {5, 0, false}, {1, 2, true}, {0, 0, true}, {0, 1, true}} {
+		in := opIn(t, approval.FormatWebAuthn)
+		in.Approver.Passkey.SignCount = c.stored
+		in.Artifact = opAssert(t, a, opExp, c.asserted)
+		s, err := approval.ConfirmOperatorAction(in)
+		if c.ok != (err == nil) || (!c.ok && (s != nil || !errors.Is(err, approval.ErrSignCountRegression))) {
+			t.Errorf("counter stored %d asserted %d: %v", c.stored, c.asserted, err)
+		}
+	}
+}
+
+func TestOperatorConfirmBinding(t *testing.T) {
+	edits := opFieldEdits()
+	edits["id"] = func(a *approval.OperatorAction) { a.ID = "oa-ffffffffffffffffffffffffffffffff" }
+	edits["kind"] = func(a *approval.OperatorAction) { *a = opAction("finish") }
+	for _, format := range []string{approval.FormatSSH, approval.FormatWebAuthn} {
+		want := map[string]error{approval.FormatSSH: approval.ErrBadSignature, approval.FormatWebAuthn: approval.ErrChallenge}[format]
+		for name, edit := range edits {
+			in := opIn(t, format) // signs the start action
+			edit(&in.Action)
+			d := opDigest(t, in.Action)
+			in.Admit.StoredDigest, in.Admit.BytesDigest = d, d // a consistent row for another action
+			opRefused(t, format+" another "+name, in, want)
+		}
+		in := opIn(t, format)
+		in.Admit.StoredDigest = opDigest(t, opAction("finish"))
+		in.Admit.BytesDigest = in.Admit.StoredDigest
+		var be *approval.BindingError
+		if s, err := approval.ConfirmOperatorAction(in); s != nil || !errors.As(err, &be) || be.Field != "digest" {
+			t.Errorf("%s: the row's digest is another action's: %v", format, err)
+		}
+		in = opIn(t, format)
+		if format == approval.FormatSSH {
+			in.Artifact = opSSH(t, opSigner, approval.OperatorActionNamespace, opAction("start"), opExp+1)
+		} else {
+			in.Artifact = opAssert(t, opAction("start"), opExp+1, 2)
+		}
+		opRefused(t, format+" another expiry", in, want)
+	}
+	a := opAction("start")
+	d := opDigest(t, a)
+	vm, _ := approval.VerbMessage(d, opC, opExp)
+	mm, _ := approval.SignedMessage(d, opC, opExp)
+	for name, sig := range map[string][]byte{"a V3 verb SSHSIG": opSigner.Sign(approval.VerbNamespace, vm), "an A3 mint SSHSIG": opSigner.Sign(approval.Namespace, mm)} {
+		in := opIn(t, approval.FormatSSH)
+		in.Artifact = sig
+		opRefused(t, name, in, approval.ErrNamespace)
+	}
+	in := opIn(t, approval.FormatSSH)
+	in.Artifact = opSigner.Sign(approval.OperatorActionNamespace, vm)
+	opRefused(t, "the verb message under the operator namespace", in, approval.ErrBadSignature)
+	vh, _ := approval.VerbMessageHash(d, opC, opExp)
+	in = opIn(t, approval.FormatWebAuthn)
+	in.Artifact = opPasskey.Assert(paOrigin, paRPID, vh, 2)
+	opRefused(t, "a V3 verb assertion", in, approval.ErrChallenge)
+}
+
+func TestOperatorConfirmExpiry(t *testing.T) {
+	a := opAction("start")
+	for _, format := range []string{approval.FormatSSH, approval.FormatWebAuthn} {
+		for _, tc := range []struct {
+			name                string
+			now, created, exp   int64
+			consumed, otherByte bool
+			want                error
+		}{
+			{"now at expiry", opExp, opCreated, opExp, false, false, &approval.NotAdmittedError{Verdict: approval.RefuseExpired}},
+			{"now past expiry", opExp + 1, opCreated, opExp, false, false, &approval.NotAdmittedError{Verdict: approval.RefuseExpired}},
+			{"now one before expiry", opExp - 1, opCreated, opExp, false, false, nil},
+			{"now before creation", opCreated - 1, opCreated, opExp, false, false, &approval.NotAdmittedError{Verdict: approval.RefuseClockBehind}},
+			{"now at creation", opCreated, opCreated, opExp, false, false, nil},
+			{"lifetime 600", opNow, opCreated, opCreated + 600, false, false, nil},
+			{"lifetime 300", opNow, opCreated, opCreated + 300, false, false, nil},
+			{"lifetime 601", opNow, opCreated, opCreated + 601, false, false, approval.ErrOperatorLifetime},
+			{"consumed", opNow, opCreated, opExp, true, false, &approval.NotAdmittedError{Verdict: approval.RefuseConsumed}},
+			{"stored bytes differ", opNow, opCreated, opExp, false, true, &approval.NotAdmittedError{Verdict: approval.RefuseDigestMismatch}},
+		} {
+			in := opIn(t, format)
+			in.Admit.Now, in.Admit.CreatedAt, in.Admit.ExpiresAt, in.Admit.Consumed = tc.now, tc.created, tc.exp, tc.consumed
+			if tc.otherByte {
+				in.Admit.BytesDigest[0] ^= 1
+			}
+			if format == approval.FormatSSH {
+				in.Artifact = opSSH(t, opSigner, approval.OperatorActionNamespace, a, tc.exp)
+			} else {
+				in.Artifact = opAssert(t, a, tc.exp, 2)
+			}
+			s, err := approval.ConfirmOperatorAction(in)
+			var na, wna *approval.NotAdmittedError
+			switch {
+			case tc.want == nil:
+				if err != nil || s.Binds(a, opC, tc.exp) != nil {
+					t.Errorf("%s %s: %v", format, tc.name, err)
+				}
+			case errors.As(tc.want, &wna):
+				if s != nil || !errors.As(err, &na) || na.Verdict != wna.Verdict {
+					t.Errorf("%s %s: %v, want %v", format, tc.name, err, wna.Verdict)
+				}
+			default:
+				if s != nil || !errors.Is(err, tc.want) {
+					t.Errorf("%s %s: %v, want %v", format, tc.name, err, tc.want)
+				}
+			}
+		}
+	}
+}
+
+// A negative expiry has no message: refused as such, never verified over an empty message.
+func TestOperatorConfirmNegativeExpiry(t *testing.T) {
+	for format, want := range map[string]error{approval.FormatSSH: approval.ErrNegativeExpiry, approval.FormatWebAuthn: approval.ErrChallenge} {
+		in := opIn(t, format)
+		in.Admit.Now, in.Admit.CreatedAt, in.Admit.ExpiresAt, in.Approver.EnrolledAt = -50, -100, -5, -100
+		opRefused(t, format+" negative expiry", in, want)
+	}
+}
+
+// The runner's recheck from its own frozen action: no admission row, the counter already recorded,
+// the key and the expiry judged at now.
+func TestOperatorRunnerRecheckLaw(t *testing.T) {
+	a := opAction("start")
+	check := func(format string) approval.OperatorCheckInput {
+		in := opIn(t, format)
+		c := approval.OperatorCheckInput{Action: a, Challenge: opC, Expiry: opExp, Now: opExp - 1, Artifact: in.Artifact, Format: format, Approver: in.Approver}
+		if format == approval.FormatWebAuthn {
+			c.Approver.Passkey.SignCount = 2 // the confirmation's counter, recorded
+		}
+		return c
+	}
+	for _, format := range []string{approval.FormatSSH, approval.FormatWebAuthn} {
+		if s, err := approval.CheckOperatorAction(check(format)); err != nil || s.Binds(a, opC, opExp) != nil {
+			t.Fatalf("%s control: %v", format, err)
+		}
+		for _, tc := range []struct {
+			name string
+			edit func(*approval.OperatorCheckInput)
+			want error
+		}{
+			{"now at expiry", func(c *approval.OperatorCheckInput) { c.Now = opExp }, &approval.NotAdmittedError{}},
+			{"revoked since", func(c *approval.OperatorCheckInput) { c.Approver.RevokedAt = 1500 }, approval.ErrApproverRevoked},
+			{"another frozen action", func(c *approval.OperatorCheckInput) { c.Action.Task += "!" }, map[string]error{approval.FormatSSH: approval.ErrBadSignature, approval.FormatWebAuthn: approval.ErrChallenge}[format]},
+			{"another challenge", func(c *approval.OperatorCheckInput) { c.Challenge[1] ^= 1 }, map[string]error{approval.FormatSSH: approval.ErrBadSignature, approval.FormatWebAuthn: approval.ErrChallenge}[format]},
+			{"an invalid action", func(c *approval.OperatorCheckInput) { c.Action.Store = "x" }, approval.ErrPath},
+			{"no format", func(c *approval.OperatorCheckInput) { c.Format = "" }, approval.ErrVerbFormat},
+		} {
+			c := check(format)
+			tc.edit(&c)
+			var na *approval.NotAdmittedError
+			s, err := approval.CheckOperatorAction(c)
+			if _, isNA := tc.want.(*approval.NotAdmittedError); s != nil || (isNA && (!errors.As(err, &na) || na.Verdict != approval.RefuseExpired)) || (!isNA && !errors.Is(err, tc.want)) {
+				t.Errorf("%s %s: %v, %v", format, tc.name, s, err)
+			}
+		}
+	}
+	for _, c := range []struct {
+		stored uint32
+		ok     bool
+	}{{1, false}, {2, true}, {3, true}} {
+		in := check(approval.FormatWebAuthn)
+		in.Approver.Passkey.SignCount = c.stored
+		s, err := approval.CheckOperatorAction(in)
+		if c.ok != (err == nil) || (!c.ok && (s != nil || !errors.Is(err, approval.ErrOperatorCounter))) {
+			t.Errorf("recheck, stored counter %d against the assertion's 2: %v", c.stored, err)
+		}
+	}
+}
````

### Patch M4 — pure entry points write nothing

````diff M4
diff --git a/host/approval/operator_silent_test.go b/host/approval/operator_silent_test.go
new file mode 100644
index 00000000..a5814802
--- /dev/null
+++ b/host/approval/operator_silent_test.go
@@ -0,0 +1,55 @@
+package approval_test
+
+// Row 243 V4 sprint 2 M4 (judge-407's retro rule): every new pure entry point a later sprint calls from
+// a server writes nothing to stdout, stderr or the log, on its accepting and its refusing paths.
+
+import (
+	"io"
+	"log"
+	"os"
+	"strings"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+)
+
+// silent runs f with os.Stdout, os.Stderr and the log on one pipe and returns what was written.
+func silent(t *testing.T, f func()) string {
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
+func TestOperatorPureEntryPointsWriteNothing(t *testing.T) {
+	if got := silent(t, func() { _, _ = os.Stderr.WriteString("x"); _, _ = os.Stdout.WriteString("y"); log.Print("z") }); !strings.HasPrefix(got, "xy") || !strings.HasSuffix(got, " z\n") {
+		t.Fatalf("control: the harness saw %q, want x, y and the log line", got)
+	}
+	got := silent(t, func() {
+		for _, b := range [][]byte{opGoldenBytes(t), []byte(`{"a":1,"a":2}`), []byte("\xff")} {
+			a, _ := approval.DecodeOperatorAction(b)
+			_, _ = approval.OperatorReviewDigest(a)
+			_, _ = approval.OperatorSignatureDigest(a)
+		}
+		_, _ = approval.OperatorActionMessageHash([32]byte{}, [32]byte{}, -1)
+		for _, format := range []string{approval.FormatSSH, approval.FormatWebAuthn, "typed"} {
+			in := opIn(t, format)
+			_, _ = approval.ConfirmOperatorAction(in)
+			_, _ = approval.CheckOperatorAction(approval.OperatorCheckInput{Action: in.Action, Challenge: opC, Expiry: opExp, Now: opNow,
+				Artifact: []byte("y\n"), Format: format, Approver: in.Approver})
+		}
+	})
+	if got != "" {
+		t.Fatalf("a pure entry point wrote %q", got)
+	}
+}
diff --git a/host/verbs/operator_silent_test.go b/host/verbs/operator_silent_test.go
new file mode 100644
index 00000000..21720a50
--- /dev/null
+++ b/host/verbs/operator_silent_test.go
@@ -0,0 +1,35 @@
+package verbs
+
+// Row 243 V4 sprint 2 M4 (judge-407 survivor: AdmitOperator without fs.SetOutput(io.Discard) survived
+// host/verbs): the pure bar parser writes nothing to stdout or stderr, on the lines whose parse the flag
+// package would otherwise report (an undefined flag, a bad bool, a missing value, help).
+
+import (
+	"io"
+	"os"
+	"testing"
+)
+
+func TestAdmitOperatorWritesNothing(t *testing.T) {
+	capture := func(f func()) string {
+		r, w, _ := os.Pipe()
+		out, errOut := os.Stdout, os.Stderr
+		os.Stdout, os.Stderr = w, w
+		f()
+		os.Stdout, os.Stderr = out, errOut
+		_ = w.Close()
+		b, _ := io.ReadAll(r)
+		return string(b)
+	}
+	if got := capture(func() { _, _ = os.Stderr.WriteString("x") }); got != "x" {
+		t.Fatalf("control: the harness saw %q", got)
+	}
+	if got := capture(func() {
+		for _, line := range []string{"episode start e1 --bogus", "episode start e1 --preset p --launch=maybe", "episode finish e1 --preset",
+			"episode start e1 --help", "episode finish e1 -h", "episode start e1 --preset p --launch --headless"} {
+			admitLine(t, line, allYes)
+		}
+	}); got != "" {
+		t.Fatalf("AdmitOperator wrote %q", got)
+	}
+}
````
