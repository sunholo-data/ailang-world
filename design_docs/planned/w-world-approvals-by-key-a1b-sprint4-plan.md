# Sprint plan — w-world-approvals-by-key-a1b, sprint 4 (iteration 378, row 201)

> **THE CUT.** This sprint closes the design's code: the decided list, what a revoked approver
> approved, the sprint-3 judge's test gaps, and, behind a measured gate, `--sign-with` with the
> design's M6 whole.
> - **Committed — Runs 1–2, 359 changed lines, measured.** **T1**, the sprint-3 judge's three test
>   gaps (SIG-NOW, SIG-FORMAT, A4-EXP; tests only — item 4 is already killed at the base, §1.2 C13);
>   **D1** `DecidedMintRequests` (store, CS6); **D2** `approve --list --decided` (F3, F15, F16, CR7);
>   **D3** `approver revoke`'s listing (K8, K9, CR12; the `runApproverRevoke|Background` pin 2 → 3).
> - **Gated — Runs 3–4, 537 changed lines, measured.** **S1** the `--sign-with` seam (F1/F2 widened,
>   F7, F7a, F8, `signForApproval`'s root; no production signer yet); **S2** its fake-signer tests (CR8,
>   CR9); **S3** the production signer `signWithSSHKeygen` with its real-tool tests (G1–G3 with a fake
>   `ssh-keygen` on `PATH`, XS1 and XS2 with the real tool), the CI step widened, QUICKSTART §15's
>   `--sign-with` and `--list --decided` lines. They run **only if** the committed `cmd/ host/` numstat
>   after the D3 commit is **≤ 1,463** (= 2,000 − 537), and **all three run or none** (C-D2: the
>   production exec never lands without its real-tool tests; S1 deliberately ships no production signer,
>   so an S1-only tree refuses `--sign-with` with F7's internal error).
> - **After this sprint, A1b's code is complete** (M1–M6). What remains is attended: Mark's exit run
>   (design §8) and QUICKSTART §15's first verbatim run (Appendix C).
>
> **Every line of code and test is a unified diff in Appendix P**, generated from scratch stages that
> were compiled, vetted, gofmt-clean and green at every sub-milestone, and re-applied — **from this
> file**, with §2's `awk` — to a fresh export of `1fce6999` (§1.2 C30–C31). The executor applies the
> patch and runs the acceptance; it writes no code by hand. Total **896** `cmd/ host/` lines
> (+13 YAML, +5 docs).

- **Design**: [w-world-approvals-by-key-a1b.md](w-world-approvals-by-key-a1b.md), **r2** (commit
  `5e2bf772`). **Not reopened.** Where this plan and the design differ, this plan was measured against
  `1fce6999` and wins; each difference is a decision (§1.3) or a premise found false (§1.4).
- **Starting ground**: Appendix C of [the sprint-3 plan](w-world-approvals-by-key-a1b-sprint3-plan.md)
  (C-D1…C-D4, its measured premises and carried mutants). Sprint 3 landed whole (M4a–M5c) as PR #541,
  merge `f0ce26f7`, recorded `1fce6999`.
- **Planner base**: `1fce6999ebf875299254a1d77003e9e3351a3347` (= origin/dev). Worktree
  `/Users/voightkampff/dev/sunholo-data/.wt-world-iter378-row201-a1b-s4`, branch
  `sprint/world-iter378-row201-a1b-s4`.
- **Executor**: as routed by the controller. **Four runs**: Run 1 (T1, D1), Run 2 (D2, D3), **size
  gate**, Run 3 (S1, S2), Run 4 (S3). Each sub-milestone is one `patch -p1` (§2) plus its acceptance
  (§3). **Apply, never retype.**
- **Touched** (committed part): `host/store/mint_requests.go`, `host/store/decided_requests_test.go`
  (new), `host/store/context_roots_test.go` (one pin), `host/authority/mint_derive_test.go`,
  `cmd/ailang-worldd/{approve.go, approver.go, main.go, help_test.go, approve_sig_test.go,
  approver_revoke_test.go}`, `cmd/ailang-worldd/approve_decided_test.go` (new). **Gated part adds**:
  `cmd/ailang-worldd/{approve_sig.go, session.go}`, `cmd/ailang-worldd/approve_signwith_test.go` (new),
  `host/approval/approvaltest/sign_crosscheck_test.go` (new), `.github/workflows/ci.yml` (the row 201
  step only), `docs/QUICKSTART.md` (§15, attended-pending, docs only).
- **Not touched**: `cmd/ailang-worldd/{session_new.go, mint_request.go, approve_test.go,
  approver_test.go, approver_enroll_test.go, approver_rig_test.go, approver_stores_test.go,
  keyguard.go}`, `host/store/{approvers.go, credentials.go, credentials_schema.sql, approvers_test.go,
  mint_requests_test.go, consume_binding_test.go}`, `host/authority/mint.go`, `host/approval/` except
  the new `approvaltest/sign_crosscheck_test.go`, `host/approval/testdata/**`, `host/daemon/**`,
  `world/*.ail`, `packages/**`, `scripts/**`, `tools/launchd/*`, any skill, `~/.ssh`, `~/.ailang/world`,
  `~/.ailang/state/mission-v1*`.

---

## 0. Binding constraints (each is restated where it governs)

1. **Fakes, temp stores and test keys only (D-WORLD-109 = A).** Every store is under `t.TempDir()`.
   Keys are `approvaltest.NewSigner(7|11|12)` or, in XS1 only, keys the test makes with `ssh-keygen -q
   -t … -N '' -C xs1 -f <t.TempDir()>/k` (stdin `/dev/null`, `SSH_AUTH_SOCK=` and a temp `HOME` in the
   child). **No `ssh-keygen` runs except from inside the named Go tests** (XS1 `TestSignWithSSHKeygenEndToEnd`,
   XS2 `TestTestSignerAgreesWithSSHKeygen`, the existing XC1), all in temp dirs; the G1–G3 tests run a
   fake `ssh-keygen` shell script put first on `PATH`. The executor never types `ssh-keygen`. Nothing
   opens `~/.ailang/world`, `~/.ssh`, a real credentials store, a real daemon, a real agent or
   `/dev/tty`. No mint outside a temp store, no publish. QUICKSTART §15's new lines are
   attended-pending and **are not run**.
2. **The production exec lands only with its real-tool tests (C-D2).** `signWithSSHKeygen` and the
   `realSessionEnv` entry that wires it arrive in **S3**, together with G1–G3, XS1, XS2 and the CI step.
   S1/S2 build and test the seam with a fake signer; between S1 and S3 production `--sign-with` is F7's
   internal error (no exec path exists).
3. **One mutant per new production branch** (§4: 68 rows run, 66 killed, 2 named equivalents survived as
   predicted; plus 4 base controls). Each row names its single killing test.
4. **Exit codes by constant**: `exitOK = 0`, `exitUsage = 1`, `exitFatal = 2` (`main.go:72-74`). Every
   refusal before the write is `exitUsage`; F7 (no signer) is `exitFatal` (an internal error, like the
   existing `after == nil`); K8's read failing after a committed revoke is `exitFatal` (a failure after
   the write).
5. **Context roots are pinned.** D3 moves `approver.go|runApproverRevoke|Background` 2 → 3 (K8's own 3 s
   read-only root, opened after the write's root and writer have ended — B-D4); S1 adds
   `approve_sig.go|signForApproval|Background: 1` (the signer's root, bounded by the request's remaining
   lifetime, made only after F7a has shown `rem > 0`). `TestProductionContextRoots` is edited for exactly
   these two pins.
6. **No `t.Parallel()`** in any new test (they swap package vars and env: `fetchDaemonHealth`, `HOME`,
   `PATH`, `SSH_AUTH_SOCK`).
7. **Bisectable.** Each sub-milestone ends gofmt-clean, `go build ./...` and `go vet` green, the `_test.go`
   compile fence green, with its own tests and every earlier sub-milestone's tests green (§1.2 C31).
8. **No git writes by the executor** (§7). The controller commits from snapshots.

---

## 1. Pre-flight

### 1.1 Baseline at `1fce6999` (measured by the planner, 2026-11-10)

Every command ran with `export AILANG_BIN=$HOME/.pinned-ailang/ailang` (`AILANG v0.41.0`), `go1.26.6
darwin/arm64`, in a `git archive 1fce6999` export at `~/.ailang/state/world-iter378/plan-scratch`
(byte-identical to the worktree; the worktree itself was only read).

| Command | rc | Time | Result |
|---|---|---|---|
| `git -C <worktree> status --short \| wc -l` | 0 | — | `0` (clean) |
| `go build ./... && go vet ./cmd/ailang-worldd/ ./host/approval/... ./host/store/ ./host/authority/` | **0** | — | clean |
| `go test -run '^$' ./...` (compile fence for every `_test.go`) | **0** | 27 s | no FAIL line |
| `go test -count=1 ./cmd/ailang-worldd/ ./host/approval/... ./host/store/ ./host/authority/` | **0** | 140 s | `cmd/ailang-worldd` ok 124.1 s, `host/approval` ok 0.4 s, `approvaltest` ok 0.6 s, `host/store` ok 13.3 s, `host/authority` ok 1.1 s; **no pre-existing failure** |

Not run by the planner at the base: `./scripts/verify_ail.sh` (no `.ail` is touched), the whole-repo
`go test ./...`, `verify_go.sh` and the full race suite (the controller's gates, §6). Known flakes from
iterations 375–377 (`host/broker` `TestExecProbeNet6ForeignDialIsRetriedNotRefused`; local `-race`
timeouts under load in `host/broker`, `host/daemon`, `host/verifygate`) are in untouched packages.

### 1.2 Codebase claims checked (at `1fce6999` unless "scratch"; scratch = the stage trees under `~/.ailang/state/world-iter378/`, outside the repo)

| # | Claim | Command → observed output |
|---|---|---|
| C1 | The base | `git rev-parse HEAD` → `1fce6999ebf875299254a1d77003e9e3351a3347`; `git status --short \| wc -l` → `0` |
| C2 | Exit codes | `grep -n "exitOK *= 0\|exitUsage *= 1\|exitFatal *= 2" cmd/ailang-worldd/main.go` → `72`, `73`, `74` |
| C3 | **No store method lists consumed rows** | `grep -n "consumed_at IS NOT NULL" host/store/mint_requests.go` → `89:` (a field comment) and `329:` (`classifyConsume`'s scan) only; control `grep -n "WHERE consumed_at IS NULL AND expires_at > ?"` → `134:` (the cap count), `182:` (`PendingMintRequests`) |
| C4 | The template and its constants | `grep -n` → `36: const mintRequestsSchemaVersion = 4`, `39: const MaxPendingMintRequests = 16`, `171: func (c *Credentials) PendingMintRequests(ctx, now, limit)` (`ORDER BY created_at DESC, request_id DESC`) |
| C5 | **`approve` has no `--decided` and no `--sign-with`** | `grep -c -- '--sign-with\|"sign-with"\|"decided"\|--decided'` on `approve.go`, `approve_sig.go`, `main.go` → `0 0 0`; control `grep -c -- '--write-message' approve.go` → `3`. Flags at `approve.go:113` (`list`), `:114` (`signature`), `:115` (`write-message`); `--list` is `:159-174` |
| C6 | `runApproverRevoke` today | `approver.go:377 func runApproverRevoke`; `:394` the read-only root `pre`; `:421` the write root `ctx` with `:422 defer cancel()` and `:428 defer func() { _ = st.Close() }()`; `:437` `fmt.Fprintf(stdout, "revoked approver %s\n", name) // K7` — so a third root after the write needs the write's defers made explicit (D3 does) |
| C7 | The roots pins | `grep -n` in `host/store/context_roots_test.go` → `20: …approve.go|runApprove|Background": 1`, `23: …approver.go|runApproverRevoke|Background": 2` |
| C8 | One existing test pins revoke's whole stdout | `approver_revoke_test.go:42: if out.String() != "revoked approver mark-mac\n" {` (D3 edits it to add K9's line); `approver_stores_test.go`'s `revoke/*` arms assert the rows only (read, `:129-141`) |
| C9 | `sessionEnv` has four seams, no `sign` | `sed -n 25,38p cmd/ailang-worldd/session.go` → `openTerminal`, `now`, `random`, `after`; `grep -c 'sign ' session.go` → `0`; `realSessionEnv` sets `after:  time.After,` at `:53` |
| C10 | The signature path reads then parses in one function | `approve_sig.go:30` `readRegularFileCapped(path, approveMaxSignatureBytes)`, `:34` `approval.ParseSSHSIG(armored)`, `:36` `"--signature %s is not an SSH signature: %v"` (S1 splits it: `approveArmored` takes the bytes and a source label) |
| C11 | The judge's gap lines | `approve_sig.go:53` `Admit: approval.AdmitInput{Now: now, CreatedAt: row.CreatedAt, ExpiresAt: row.ExpiresAt, …` (SIG-NOW); `:57` `… Format: enrolled.Format, …` (SIG-FORMAT); `host/authority/mint.go:96` `rm.Signed.Binds([32]byte(digest), rm.Challenge, rm.ExpiresAt)` (A4-EXP); `host/store/approvers.go:204` `UPDATE approvers SET revoked_at = ? WHERE name = ? AND revoked_at IS NULL AND enrolled_at <= ?` (item 4) |
| C12 | `session revoke` takes the 64-hex credential id and `--credentials-db` | `session.go:361 credsDB := fs.String("credentials-db", …)`, `:383 if !isHex64(credID) {` |
| C13 | **The gap mutants at the base** (`plan-scratch`, pristine tests, the §4 harness) | MUT-SIG-NOW with `-run '^TestApprove'` → **SURVIVED** (11.7 s); MUT-SIG-FORMAT with `^TestApproveBySignatureRefusalsBeforeAnyEffect$` → **SURVIVED**; MUT-A4-EXP with `^TestMintForRequestRefusesAnUnboundApproval$` → **SURVIVED**; MUT-REVOKE-ACTIVE (item 4) with `^TestRevokeApprover$` → **KILLED** at the base (its step 3, "already revoked: refused, the first revoke time kept"), so item 4 needs no new test. With T1 applied all four die by their single named test (§4) |
| C14 | The CI step | `ci.yml:141` `- name: Row 201 SSHSIG cross-check with real ssh-keygen, verbose`, `:142 timeout-minutes: 3`, `:150 -run '^TestSSHSIGCrossCheckWithSSHKeygen$'` over `./host/approval/` only, `:153` the PASS-grep loop over the parent and two subtests. It is the **first** compile step of the `go-verify` job (`:116`; the next Go build is `:295`) |
| C15 | QUICKSTART §15 | `grep -n '^## 15\. ' docs/QUICKSTART.md` → `1785: ## 15. Enrol an approver *(attended — pending first verbatim run)*`; `grep -c -- '--sign-with\|--decided'` → `0` |
| C16 | The tool on this machine | `command -v ssh-keygen` → `/usr/bin/ssh-keygen`; `ssh -V` → `OpenSSH_10.3p1, LibreSSL 3.3.6` |
| C17 | The helpers the new tests call exist with these signatures | `grep -n` → `approve_sig_test.go:52-53` `sigSigner`/`otherSigner`, `:57 requestMessage`, `:76 writeSignature`, `:87 signedRig`; `approve_test.go:57 clockSeq`, `:87 newApproveRigIn`, `:399 approveEnv`, `:489 rawCredentials`, `:499 atTime`; `approver_rig_test.go:43 approverStore`, `:93 seedApprover`; `approver_test.go:26 lockedTerm`; `host/store`: `mint_requests_test.go:30 makeV3Credentials`, `:162 testMintRequestRow`, `:175 testCtx`, `approvers_test.go:16 makeV4Credentials`, `credentials_test.go:104 openTestCredentials`; `host/authority/mint_derive_test.go:21 au1Store`, `:40 au1Mint`, `:102 au1Enrol`, `:136 au1SignedMint` |
| C18 | File sizes stay modest | `wc -l` at the base → `approve.go` 263, `approve_sig.go` 97, `approver.go` 439, `session.go` 540, `mint_requests.go` 350; after S3: 329, 182, 467, 544, 390 (World has no file-size gate: `ls scripts \| grep -i size` → none) |
| C19 | **Every stage compiles and passes in scratch** | at each stage: `gofmt -l cmd/ailang-worldd host/approval host/store host/authority` empty; `go build ./...`, `go vet` of the touched packages clean. At S3: `go vet ./...` rc 0; `go test -run '^$' ./...` rc 0; `go test -count=1 ./cmd/ailang-worldd/ ./host/approval/... ./host/store/ ./host/authority/` → all ok (`cmd/ailang-worldd` 127.1 s, `host/approval` 0.9 s, `approvaltest` 0.5 s, `host/store` 13.1 s, `host/authority` 1.1 s); `TestProductionContextRoots` ok; `TestDaemonDependencyAllowlist` ok; the QUICKSTART tests ok; `go test -race ./host/approval/... ./host/store/ ./host/authority/` ok; `-race -count=3` on the new and touched cmd tests ok (21.1 s); the §6 `-v` run: 0 `--- SKIP`, 0 `--- FAIL` |
| C20 | The real tool, through the tests only | `WORLD_REQUIRE_SSH_KEYGEN=1 go test ./host/approval/... ./cmd/ailang-worldd/ -run '^(TestSSHSIGCrossCheckWithSSHKeygen\|TestTestSignerAgreesWithSSHKeygen\|TestSignWithSSHKeygenEndToEnd)$' -v` → PASS for all seven names of the widened CI grep, 0 SKIP; the CI step's script, run verbatim locally without `AILANG_BIN` → `rc=0`, `ALL_PASS_GREP_OK` |
| C21 | The CI step's compile time (why S3 raises its timeout 3 → 8 min) | warm cache: the step's `go test` 1 s. Cold: `GOCACHE=$(mktemp -d) go test -count=1 -run '^$' ./host/approval/... ./cmd/ailang-worldd/` → `cold_rc=0 cold_t=24s` (§1.4 P6) |
| C30 | **The patches reproduce the tree** | fresh `git archive 1fce6999` export + `patch -p1` of Appendix P's seven patches **extracted from this file with §2's `awk`** → `diff -r` of `cmd host docs .github` against the scratch final tree: identical; no `.orig`/`.rej` |
| C31 | **Every sub-milestone is bisectable** | the same fresh export, after **each** patch: gofmt empty, build + vet clean, compile fence rc 0, the §3 acceptance commands rc 0, and the cumulative numstat (in a throwaway git repo, after a commit of the stage) = §7's value: **33, 181, 278, 359, 461, 615, 896** |

### 1.3 Planner decisions where the design or Appendix C leaves a choice

- **S4-D1 — The cut** (§1.5). Committed: T1, D1, D2, D3 (359). Gated, all-or-none: S1, S2, S3 (537),
  only if the committed numstat after D3 is ≤ 1,463. The directive's ×1.2 factor is not needed on
  measured numbers; the gate leaves 1,104 lines of margin.
- **S4-D2 — `DecidedMintRequests(ctx, approvedBy string, limit int)`.** The design's signature is
  `(ctx, limit)`. K8 needs "the decided rows whose `approved_by` is `ssh:<name>`"; filtering the newest 64
  in Go would silently miss an approver's older approvals once 64 newer rows exist. So the filter is in
  the SQL (`AND approved_by = ?` when non-empty); `""` means every decided row (`--list --decided`). The
  cap (64) and the order (`consumed_at DESC, request_id DESC`) are the design's.
- **S4-D3 — K8 filters by the row's own `format:name`** (`row.Format+":"+name`, from revoke's read-only
  lookup), not a hard-coded `ssh:`. The `"ssh:"+name` form is a named equivalent (P3 admits only `ssh`).
- **S4-D4 — K8's output** is `it approved N session(s), newest first; each stays valid until it expires
  or is revoked:` then per row `  <id>  episode <ep>  at <consumed RFC 3339>  credential <16 hex>` and
  `    ailang-worldd session revoke --db <db> --credentials-db <creds> <64-hex id>` (shell-quoted paths;
  CR12 runs the printed line through `runSessionRevoke` and checks it revokes exactly that session).
  K9: `no session was approved by it`. It revokes nothing itself (Q4). A K8 read failure after the
  committed revoke is `exitFatal` naming `approve --list --decided` as the fix (propagation-only; §4).
- **S4-D5 — F16's line**: `<id>  episode <ep>  approved by <approved_by>  at <consumed RFC 3339>
  credential <first 16 hex>`; F15: `no decided mint requests`. `--list --decided` reads on the existing 3 s
  `--list` root (no new root) and works on a schema-4 store (CS6 `v4`).
- **S4-D6 — F1/F2 as a count.** `modes` counts the non-empty `--signature`, `--sign-with`,
  `--write-message`; F1 is `modes > 1`, F2 `--list && modes != 0`. It replaces sprint 3's two-flag
  conditions (whose mutants MUT-F1/F2/F2-SIG/F2-MSG become MUT-F1-COUNT/-INC/-SIG/-SIGNWITH/-MSG and
  MUT-F2, §4).
- **S4-D7 — F7 is checked with `after`, before any store read**, and only when `--sign-with` is given:
  `env.sign == nil` → `exitFatal` "internal error: no signer (sessionEnv.sign)". F7a and F8 live in
  `signForApproval`.
- **S4-D8 — `--sign-with`'s bytes go through the same judge as `--signature`'s.** `approveSignature` is
  split: it reads the file (F9, F10) and calls the new `approveArmored(ctx, credPath, source, armored,
  row, now)` (F11–F14), which `--sign-with` also calls with source `ssh-keygen's output for --sign-with
  <key>`. Sprint 3's `--signature` texts are unchanged (`--signature <path> is not an SSH signature: …`).
  `approvalMessage(row)` is factored out of `writeApprovalMessage` and shared with `signForApproval`;
  `expiredRefusal(row)` is factored out of `approveAdmission` and shared with F7a, so F7a prints the
  standard text exactly (design §2.7 [r1]).
- **S4-D9 — `signWithSSHKeygen` lives in `approve_sig.go`** beside its only caller (design §2.9 put it in
  `session.go`); `session.go` gains the `sign` field (S1) and the `realSessionEnv` entry (S3). The root
  pin is therefore `approve_sig.go|signForApproval` (design: `approve.go|…`).
- **S4-D10 — `cmd.WaitDelay = time.Second`** in `signWithSSHKeygen` (not in the design). Measured: a
  signer that leaves a child holding stdout (the fake `cat; sleep 20`) keeps `cmd.Run` waiting ~20 s past
  the deadline without it, 2.1 s with it (`outlives-the-request` arm; MUT-WAITDELAY takes 23.7 s and dies
  on the arm's 10 s bound).
- **S4-D11 — G1–G3 are tested with a fake `ssh-keygen` on `PATH`, always** (no skip): the exec line and
  stdin are recorded by the script and asserted exactly; a non-zero exit (G1, with its stderr passed
  through), 16 KiB + 1 of output (G2), no output (G3) and a signer outliving the request each refuse
  with `exitUsage`. XS1 (real tool, both key types) also proves G1 with the real "No private key found".
  `cappedBuffer` has its own unit test (the memory bound is not visible at the CLI: MUT-CAP-WRITE).
- **S4-D12 — `TestRealSessionEnvSignsWithSSHKeygen`** pins the production wiring (a `reflect` pointer
  comparison), because every CLI test builds its own env (MUT-REAL-ENV).
- **S4-D13 — The CI step** keeps its name and grep form; `-run` becomes the three real-tool tests over
  `./host/approval/... ./cmd/ailang-worldd/`, the PASS list gains `TestTestSignerAgreesWithSSHKeygen`,
  `TestSignWithSSHKeygenEndToEnd` and its two subtests, and `timeout-minutes` 3 → 8 (it is now the job's
  first compile of the `cmd/ailang-worldd` test binary, C14/C21). Validated: `ruby -ryaml` loads the file,
  and the step's script ran verbatim locally (C20).
- **S4-D14 — No CR13** (carried S3-D11), and **no store change for the `approvers` format column**: SIG-FORMAT's
  CLI arm rewrites the row to `webauthn` with raw SQL (the column's CHECK is only `format <> ''`).
- **S4-D15 — Commit order inside the gated tail.** S1, S2, S3 are three commits that ride one PR with the
  committed part; none merges without the others.

### 1.4 Premises found false or incomplete at HEAD (each routed around minimally)

- **P1 — Appendix C's line numbers moved** (sprint 3 edited the files): `PendingMintRequests` is `:171`
  (not `:165`), approve's `--list` is `:159-174` (not `:136-151`). `runApproverRevoke` `:377` holds. No
  routing needed; the patches carry their own context.
- **P2 — The design's `DecidedMintRequests(ctx, limit)` cannot serve K8 correctly** (S4-D2): routed by
  the `approvedBy` parameter.
- **P3 — K8 cannot simply follow the write.** The write root and the writer are released by `defer` at
  the function's end (C6), so a third root "after the write's root has ended" needs them released
  explicitly first. Routed in D3: `cancel()` and `st.Close()` right after `approverRevoke`, then K8's own
  root. The roots count is per function, so K8's root stays in `runApproverRevoke` (pin 2 → 3, as B-D4).
- **P4 — The design pins `approve.go|signForApproval`; the signature code lives in `approve_sig.go`**
  since sprint 3. Routed by S4-D9.
- **P5 — Appendix C's estimate** (C-D1 ≈420, M6 ≈610 with the factor) vs measured: T1 + C-D1 **359**, the
  tail **537** (+13 YAML, +5 docs). Both fit; the whole design closes here.
- **P6 — The CI step's 3-minute timeout was set for `./host/approval/` alone** and the step runs before
  any other Go build in its job (C14). Widened to `cmd/ailang-worldd`, a cold compile of that test binary
  is the dominant cost. Measured locally (C21): **24 s** cold (`cold_rc=0`, Apple silicon, modules
already downloaded); a hosted 4-core runner is slower and may also fetch modules. Routed by 3 → 8 minutes
  (S4-D13); `actions/setup-go`'s module/build cache usually makes it warm.
- **P7 — The sprint-3 judge's item 4 premise** ("check whether any test kills the revoke SQL's `AND
  revoked_at IS NULL`"): it is already killed at the base by `TestRevokeApprover` step 3 (C13). No test
  is added; the mutant stays in §4 as MUT-REVOKE-ACTIVE.
- **P8 — The design's "existing test unchanged" list does not anticipate K9**: `TestApproverRevoke/yes`
  pins revoke's whole stdout (C8). Routed: D3 edits that one assertion to add K9's line.

### 1.5 The size, and the cut

Lines are insertions + deletions under `cmd/ host/` (`git diff --numstat`; docs and YAML not counted).
**Every row is measured** from the scratch stages, each committed in a throwaway repo (C31).

| Sub-milestone | Production | Test | Step total | Cumulative vs `1fce6999` (after its commit) | Run |
|---|---|---|---|---|---|
| T1 — SIG-NOW test, SIG-FORMAT arm, A4-EXP control arm | 0 | 33 | **33** | 33 | 1 |
| D1 — `DecidedMintRequests`, `MaxDecidedMintRequests` (CS6) | 40 | 108 | **148** | 181 | 1 |
| D2 — `approve --list --decided`, F3/F15/F16, help, usage (CR7, CR1 `decided-without-list`) | 37 | 60 | **97** | 278 | 2 |
| D3 — K8/K9 in `runApproverRevoke`, the roots pin 2 → 3 (CR12; `TestApproverRevoke/yes` edited) | 34 | 47 | **81** | **359** | 2 |
| **Committed** | **111** | **248** | **359** | | |
| *gate: Runs 3–4 only if the measured numstat after the D3 commit is ≤ 1,463* | | | | | |
| S1 — the seam, `--sign-with`, F1/F2 count, F7, F7a, F8, `signForApproval`, `approveArmored`, help, usage, the root pin, CR1's three `--sign-with` arms | 97 | 7 | **104** | 461 | 3 |
| S2 — CR8 `TestApproveSignWithFakeSigner`, CR9 `TestApproveSignWithDeadline` | 0 | 154 | **154** | 615 | 3 |
| S3 — `signWithSSHKeygen`, `cappedBuffer`, `realSessionEnv.sign`; G1–G3 + outlives (fake tool), `cappedBuffer` and wiring tests, XS1, XS2 (real tool); CI step (+13 YAML), QUICKSTART §15 (+5 docs) | 48 | 233 | **281** | **896** | 4 |
| **Sprint** | **256** | **640** | **896 ≤ 2,000** | | |

(The step totals sum to 898; the cumulative is 896 because S1 re-edits the one `help_test.go` line D2
edited — the usage wants, 1 + 1 each time — which numstat against the base counts once.)

---

## 2. How the executor applies a sub-milestone (the patch protocol)

Each sub-milestone `<M>` is the block fenced by a line `` ````diff <M> `` (four backticks) and a line of
four backticks, in Appendix P. Paths are `a/<repo path>` and `b/<repo path>`. Apply it **only** this
way, in one `bash` shell (the controller's shell is zsh):

```bash
bash <<'EOF'
set -euo pipefail
M=T1   # <- the sub-milestone, exactly as in the Appendix P heading
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter378-row201-a1b-s4
PLAN="$W/design_docs/planned/w-world-approvals-by-key-a1b-sprint4-plan.md"
T="$HOME/.ailang/state/world-iter378/exec-patches"; mkdir -p "$T"
awk -v m="$M" '$0 == "````diff " m {f=1; next} f && $0 == "````" {exit} f' "$PLAN" > "$T/$M.patch"
test -s "$T/$M.patch"
cd "$W"
patch -p1 --dry-run < "$T/$M.patch"
patch -p1 < "$T/$M.patch"
if find cmd host docs .github \( -name '*.orig' -o -name '*.rej' \) | grep .; then echo "STOP: patch left .orig/.rej"; exit 1; fi
git status --short
EOF
```

- If `--dry-run` fails, **stop and report** (no hand edits, no fuzz options). It can only fail if an
  earlier sub-milestone was skipped or changed.
- Then run that sub-milestone's **Acceptance** (§3) in one bash shell with
  `export AILANG_BIN=$HOME/.pinned-ailang/ailang`.
- Then snapshot (§7). The controller commits and compares the numstat with §7's expected value.

---

## 3. Milestones

Every sub-milestone's acceptance starts with the same three lines (abbreviated **[base]** below):

```
test -z "$(gofmt -l cmd/ailang-worldd host/approval host/store host/authority)"
go build ./... && go vet ./cmd/ailang-worldd/ ./host/approval/... ./host/store/ ./host/authority/
go test -run '^$' ./...
```

`CMDT` below is
`'^(TestApprove|TestEpisodeStartRequest|TestMintRequest|TestApprover|TestEveryVerbHasHelp|TestHelpIsNeverBlocked|TestExistingCallersReadUnbounded|TestReadTermLine|TestConfirmApproverAct|TestReadRegularFileCapped|TestSignWith|TestCappedBuffer|TestRealSessionEnv)'`
(the names that do not exist yet simply match nothing).

### T1 — the sprint-3 judge's test gaps · 0 prod + 33 test · Run 1

**Patch** `T1` (tests only; production is unchanged and was measured correct):
- `cmd/ailang-worldd/approve_sig_test.go`: CR4 gains the arm **`format-not-ssh`** (raw SQL sets
  mark-mac's `format` to `webauthn`; a good signature → `exitUsage`, `approval: approver format must be
  ssh`, no terminal, 0 daemon calls, pending) — kills SIG-FORMAT; new
  **`TestApproveBySignatureJudgesTheApproverAtNow`** (request at `approveT0`, approver enrolled at
  `approveT0+5`, `approve --signature` at `approveT0+10` → `exitOK`, no terminal, `ssh:mark-mac`) — kills
  SIG-NOW; its control (`ExpiresAt` → `row.CreatedAt` on the same line) dies too.
- `host/authority/mint_derive_test.go`: `TestMintForRequestRefusesAnUnboundApproval` gains the subtest
  **`bound-control`** (the unperturbed signed mint succeeds), so that test alone kills A4-EXP.

**Acceptance**
```
[base]
go test ./cmd/ailang-worldd/ -run '^(TestApproveBySignatureJudgesTheApproverAtNow|TestApproveBySignatureRefusalsBeforeAnyEffect)$' -count=1 -v
go test ./host/authority/ -run '^TestMintForRequestRefusesAnUnboundApproval$' -count=1 -v
go test ./host/store/ -run '^TestRevokeApprover$' -count=1
```

**Expected cumulative numstat after the controller's commit of T1: 33** (`git diff --numstat 1fce6999 HEAD -- cmd host`, §7).

### D1 — `DecidedMintRequests` · 40 prod + 108 test · Run 1

**Patch** `D1`:
- `host/store/mint_requests.go`: `const MaxDecidedMintRequests = 64`; `(*Credentials).DecidedMintRequests(ctx,
  approvedBy string, limit int) ([]MintRequestRow, error)`: `requireDeadline`; below schema 4 → the
  existing `*OlderCredentialsStoreError{Need: 4, Feature: "mint requests"}`; `limit <= 0 || limit > 64`
  → 64; `SELECT … WHERE consumed_at IS NOT NULL [AND approved_by = ?] ORDER BY consumed_at DESC,
  request_id DESC LIMIT ?`; scan with the existing `scanMintRequest`.
- new `host/store/decided_requests_test.go`: `putDecided` (a consumed row by raw insert); **CS6**
  `TestDecidedMintRequests`: the cap literal is 64; a pending row and no decided rows → empty; four
  decided rows come back `r3, r2, r1, r4` (equal `consumed_at` 3000 → `request_id DESC`), every column
  equal (`reflect.DeepEqual`); `ssh:x` → `r2, r4`; `ssh:nobody` → none; limit 2 → `r3, r2`; subtests
  `cap` (70 rows; limits 0, −1, 65 → 64), `v4` (a schema-4 store reads), `v3` (refused, need 4, names
  the restart), `no-deadline` (`ErrNoDeadline`).

**Acceptance**
```
[base]
go test ./host/store/ -run '^TestDecidedMintRequests$' -count=1 -v
go test ./host/store/ -count=1
```

**Expected cumulative numstat after the controller's commit of D1: 181** (`git diff --numstat 1fce6999 HEAD -- cmd host`, §7).

### D2 — `approve --list --decided` · 37 prod + 60 test · Run 2

**Patch** `D2`:
- `cmd/ailang-worldd/approve.go`: usage `approve --list [--decided] …`; help paragraph; the `--decided`
  flag; **F3** `--decided` without `--list` → help + `exitUsage`; `--list --decided` →
  `approveListDecided(ctx, ro, credPath, stdout, refuse)` on the existing read-only handle and 3 s root:
  **F15** `no decided mint requests`, **F16** one line per row (S4-D5); `credentialPrefix(id)`.
- `main.go` usage line; `help_test.go` want `approve --list [--decided]`.
- new `cmd/ailang-worldd/approve_decided_test.go`: `decidedRig` (one store: request a approved by the typed
  y at +1, b by mark-mac's signature at +2, c pending); **CR7** `TestApproveDecidedList` (stdout exactly
  b's line then a's, c absent; subtest `none` → exactly `no decided mint requests\n`).
- `approve_sig_test.go` CR1 gains `decided-without-list` (F3).

**Acceptance**
```
[base]
go test ./cmd/ailang-worldd/ -run '^(TestApproveDecidedList|TestApproveFlagRules|TestEveryVerbHasHelp|TestApproveListShowsPending|TestApproveListUsesTheNamedStores)$' -count=1 -v
go test ./cmd/ailang-worldd/ -run "$CMDT" -count=1
```

**Expected cumulative numstat after the controller's commit of D2: 278** (`git diff --numstat 1fce6999 HEAD -- cmd host`, §7).

### D3 — `approver revoke` lists what it approved (K8, K9) · 34 prod + 47 test · Run 2

**Patch** `D3`:
- `cmd/ailang-worldd/approver.go` `runApproverRevoke`: the write's root and writer are released right
  after `approverRevoke` (P3); after K7, **K8** under its own 3 s root opens the store read-only and calls
  `DecidedMintRequests(list, row.Format+":"+name, 0)`; **K9** none → `no session was approved by it`;
  else the S4-D4 lines. An open or read error → `exitFatal` naming `approve --list --decided`.
  `approverHelp`'s `revoke` paragraph says it lists and revokes nothing.
- `host/store/context_roots_test.go`: `runApproverRevoke|Background` **2 → 3**.
- `approver_revoke_test.go` `TestApproverRevoke/yes`: stdout `revoked approver mark-mac\nno session was
  approved by it\n` (P8).
- `approve_decided_test.go` after `// ---- D3 ----`: **CR12** `TestApproverRevokeListsWhatItApproved`
  (on `decidedRig`: revoke mark-mac → stdout exactly K7, the header, b's line and its `session revoke`
  line; revoke `idle` (enrolled, never used) → K7 + K9; both sessions still exist; the printed line, run
  through `runSessionRevoke`, revokes exactly b's session).

**Acceptance**
```
[base]
go test ./cmd/ailang-worldd/ -run '^(TestApproverRevokeListsWhatItApproved|TestApproverRevoke|TestApproverRevokeUnansweredPromptTimesOut|TestApproverVerbsUseTheNamedStores)$' -count=1 -v
go test ./host/store/ -run '^TestProductionContextRoots$' -count=1
go test ./cmd/ailang-worldd/ -run "$CMDT" -count=1
go test ./host/approval/... ./host/store/ ./host/authority/ -count=1
```

**Expected cumulative numstat after the controller's commit of D3: 359** (`git diff --numstat 1fce6999 HEAD -- cmd host`, §7).

**→ the size gate (§7) runs here.**

### S1 — the `--sign-with` seam · 97 prod + 7 test · Run 3 (gated)

**Patch** `S1`:
- `cmd/ailang-worldd/session.go`: `sessionEnv.sign func(ctx, key string, msg []byte, stderr io.Writer)
  ([]byte, error)` (the field only; `realSessionEnv` does not set it until S3).
- `cmd/ailang-worldd/approve.go`: usage + help for `--sign-with <key>` (R6's "a `.pub` path does not prove
  the agent signed"); the flag; **F1/F2** as a count (S4-D6); **F7** (S4-D7); after the `--signature`
  block, `--sign-with`: `signForApproval` then `approveArmored` (S4-D8); `expiredRefusal(row)`.
- `cmd/ailang-worldd/approve_sig.go`: `approveSignature` → reads, then `approveArmored(…, "--signature
  "+path, …)`; `approveArmored` (F11–F14, unchanged texts); `approvalMessage(row)`;
  **`signForApproval(env, key, row, stderr)`**: **F7a** `rem := row.ExpiresAt - env.now(); rem <= 0` →
  `expiredRefusal`; a new root `context.WithTimeout(context.Background(), rem s)`; `env.sign`; **F8**
  error → `--sign-with <key>: <err>; nothing was approved`.
- `main.go` usage line, `help_test.go` want `approve <request-id> --sign-with <key>`, `approverHelp`
  names `--sign-with`, `context_roots_test.go` `approve_sig.go|signForApproval|Background: 1`.
- `approve_sig_test.go` CR1 gains `signature-and-sign-with`, `sign-with-and-write-message`,
  `list-and-sign-with`.

**Acceptance**
```
[base]
go test ./cmd/ailang-worldd/ -run '^(TestApproveFlagRules|TestEveryVerbHasHelp|TestHelpIsNeverBlockedByTheRegistryKeyGuard|TestApproveBySignatureRefusalsBeforeAnyEffect|TestApproveWriteMessage|TestApproveBySignatureMintsWithNoTerminal)$' -count=1 -v
go test ./host/store/ -run '^TestProductionContextRoots$' -count=1
go test ./cmd/ailang-worldd/ -run "$CMDT" -count=1
```

**Expected cumulative numstat after the controller's commit of S1: 461** (`git diff --numstat 1fce6999 HEAD -- cmd host`, §7).

### S2 — the fake-signer tests (CR8, CR9) · 0 prod + 154 test · Run 3 (gated)

**Patch** `S2`: new `cmd/ailang-worldd/approve_signwith_test.go`: `fakeSigner` (records calls, key,
message and the deadline left; refuses an already-expired context), `signEnv`, `countDaemonCalls`;
**CR8** `TestApproveSignWithFakeSigner` (`no-signer` → `exitFatal`, exact stderr; `signer-error` →
`exitUsage`, exact stderr; `other-approver-other-message` → F14; `garbage` → F11 with the `--sign-with`
source; `good` → `exitOK`, `ssh:mark-mac`, the seam called once with `/keys/k` and exactly the request's
message; every refusal: no terminal, pending, no session, **0 daemon calls**); **CR9**
`TestApproveSignWithDeadline` (`37-seconds-left` → deadline within [36 s, 37 s]; `1-second-left` → (0, 1 s];
`expired-since-admission` (`clockSeq(1, +599, +600)`) → exact standard expired text, signer **never
called**; each pending).

**Acceptance**
```
[base]
go test ./cmd/ailang-worldd/ -run '^(TestApproveSignWithFakeSigner|TestApproveSignWithDeadline)$' -count=1 -v
go test -race ./cmd/ailang-worldd/ -run '^(TestApproveSignWithFakeSigner|TestApproveSignWithDeadline)$' -count=1
go test ./cmd/ailang-worldd/ -run "$CMDT" -count=1
```

**Expected cumulative numstat after the controller's commit of S2: 615** (`git diff --numstat 1fce6999 HEAD -- cmd host`, §7).

### S3 — the production signer with the real tool · 48 prod + 233 test (+13 YAML, +5 docs) · Run 4 (gated)

**Patch** `S3`:
- `approve_sig.go` after `// ---- S3 (sprint 4) ----`: `errSignatureTooLarge`, `errEmptySignature`;
  **`signWithSSHKeygen`**: `exec.CommandContext(ctx, "ssh-keygen", "-Y", "sign", "-n",
  approval.Namespace, "-f", key)`, stdin the message, stdout a `cappedBuffer` of 16 KiB, stderr the
  operator's, `WaitDelay` 1 s (S4-D10); **G1** run error → `ssh-keygen -Y sign: <err>`; **G2** over the
  cap; **G3** empty. `cappedBuffer`.
- `session.go` `realSessionEnv`: `sign: signWithSSHKeygen`.
- `approve_signwith_test.go` after `// ---- S3 ----`: `fakeSSHKeygen` (a `#!/bin/sh` script first on
  `PATH`); `TestSignWithSSHKeygenRefusesBadOutput` (`exec-line-and-no-output`: argv exactly `-Y sign -n
  world-approval@sunholo -f /keys/k`, stdin exactly the message, then G3; `failed`: G1 and the script's
  stderr passed through; `over-16k`: G2; `outlives-the-request`: `sleep 20` with 1 s left → `signal:
  killed` in under 10 s); `TestCappedBufferKeepsItsBound`; `TestRealSessionEnvSignsWithSSHKeygen`;
  **XS1** `TestSignWithSSHKeygenEndToEnd/{ed25519,ecdsa-p256}` (real keys in `t.TempDir()`, temp `HOME`,
  no agent; `--sign-with <k>` mints as `ssh:mark-mac`; `--sign-with <k>.pub` with `k` moved away →
  `exitUsage`, `No private key found`, `exit status 255`, pending; `SKIP COUNTED` unless
  `WORLD_REQUIRE_SSH_KEYGEN=1` makes a missing tool a failure).
- new `host/approval/approvaltest/sign_crosscheck_test.go`: **XS2** `TestTestSignerAgreesWithSSHKeygen`
  (`ssh-keygen -Y check-novalidate -n world-approval@sunholo` accepts `NewSigner(7)`'s signature over an
  approval message (`Good "world-approval@sunholo" signature`) and refuses it (rc 255) over the message
  with its last byte flipped; `ssh-keygen -l` of `AuthorizedKey()` prints `Fingerprint()`; same skip rule).
- `.github/workflows/ci.yml`: the row 201 step widened (S4-D13). `docs/QUICKSTART.md` §15: the
  `--sign-with` route, the `--list --decided` and K8 sentence, three new lines in the second bash block
  (a second request, `approve … --sign-with /tmp/se-world/demo-key`, `approve --list --decided`).
  **Docs only, attended-pending, not run.**

**Acceptance** (the Go tests below run `ssh-keygen` themselves, in temp dirs, §0.1; the executor runs no
`ssh-keygen` directly)
```
[base]
go test ./cmd/ailang-worldd/ -run '^(TestSignWithSSHKeygenRefusesBadOutput|TestCappedBufferKeepsItsBound|TestRealSessionEnvSignsWithSSHKeygen)$' -count=1 -v
WORLD_REQUIRE_SSH_KEYGEN=1 go test ./host/approval/... ./cmd/ailang-worldd/ -run '^(TestSSHSIGCrossCheckWithSSHKeygen|TestTestSignerAgreesWithSSHKeygen|TestSignWithSSHKeygenEndToEnd)$' -count=1 -v
go test ./cmd/ailang-worldd/ -run "$CMDT" -count=1
go test ./cmd/ailang-worldd/ -run 'Quickstart|QuickStart|QUICKSTART' -count=1
go test ./host/approval/... -run '^TestApprovalTestIsTestOnly$' -count=1
ruby -ryaml -e 'YAML.load_file(".github/workflows/ci.yml"); puts "yaml ok"'
```

**Expected cumulative numstat after the controller's commit of S3: 896** (`git diff --numstat 1fce6999 HEAD -- cmd host`, §7).

### S4 — revision: the signer's exec site moves to `host/sshsign`, scrubbed and driven (AC10) · Run 5 (after S3)

**Why** (the controller's §6 sweep at `784c9413`, `/tmp/world-iter378/gates-784c9413.log`): `go test
./host/broker/` → `--- FAIL: TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential` ("AC10(a): 13
files carry subprocess sites [... cmd/ailang-worldd/approve_sig.go ...] but 12 have drivers"). S3's
`signWithSSHKeygen` was a new `exec.CommandContext` site that (1) handed `ssh-keygen` the whole process
environment, registry credential included, against the AC10 policy every other World subprocess obeys,
and (2) lived in package `main`, which the census's driver map cannot import. The planner's §6 gate list
ran only the touched packages, so it missed `host/broker`, whose census enumerates every exec site of
`cmd/` and `host/`. **Planner miss, recorded.**

**Patch** `S4` (applies on top of the committed S3 tree, `784c9413`; it does not need §7.1's copied test
removed by hand — the patch removes it from `cmd/` and adds it to the new package):
- new `host/sshsign/sshsign.go` (the `host/toolprobe` pattern): `Sign(ctx, bin, key, msg, stderr)` runs
  `<bin> -Y sign -n world-approval@sunholo -f <key>` with **`cmd.Env = childenv.Scrubbed(os.Environ())`**
  (only the registry variables are dropped; `SSH_AUTH_SOCK`, `HOME` and `PATH` pass through), the message
  on stdin, stdout a 16 KiB `cappedBuffer`, stderr the operator's, `WaitDelay` 1 s; G1/G2/G3 with S3's
  exact texts (`ErrSignatureTooLarge`, `ErrEmptySignature`). It imports `host/approval` (stdlib-only, for
  `Namespace`) and `host/childenv`; it makes no context root (the caller's deadline).
- `cmd/ailang-worldd/approve_sig.go`: `signWithSSHKeygen` becomes `sshsign.Sign(ctx, "ssh-keygen", key,
  msg, stderr)`; the exec, `cappedBuffer`, the two errors and the imports `bytes`, `errors`, `os/exec` go.
  `realSessionEnv.sign` is unchanged, so `TestRealSessionEnvSignsWithSSHKeygen` and every CLI test
  (fake `ssh-keygen` on `PATH`, CR8, CR9, XS1) run unchanged.
- `cmd/ailang-worldd/approve_signwith_test.go`: `TestCappedBufferKeepsItsBound` (§7.1's copy) moves out.
- new `host/sshsign/sshsign_test.go`: `TestCappedBufferKeepsItsBound` (moved verbatim) and
  `TestSignHandsTheSignerAScrubbedEnvironment` (a temp shell-script signer dumps its env: no
  `AILANG_REGISTRY_API_KEY`, and `SSH_AUTH_SOCK`/`HOME` as set).
- `host/broker/registry_publish_test.go`: the driver entry `"host/sshsign/sshsign.go": {driveSSHSign}`
  and `driveSSHSign` (launches the production `Sign` with the census's env-dumping probe as the
  program; the probe's banner comes back as the "signature", and the launch and its env are the
  measurement), plus the `host/sshsign` import.

**Acceptance** (in scratch, all rc 0 on the tree this patch produces from a fresh `784c9413` export)
```
[base]
go test ./host/broker/ -run '^TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential$' -count=1 -v
go test ./host/sshsign/ -count=1 -v && go test -race ./host/sshsign/ -count=1
go test ./cmd/ailang-worldd/ -run '^(TestSignWithSSHKeygenRefusesBadOutput|TestRealSessionEnvSignsWithSSHKeygen|TestApproveSignWithFakeSigner|TestApproveSignWithDeadline)$' -count=1 -v
WORLD_REQUIRE_SSH_KEYGEN=1 go test ./host/approval/... ./cmd/ailang-worldd/ -run '^(TestSSHSIGCrossCheckWithSSHKeygen|TestTestSignerAgreesWithSSHKeygen|TestSignWithSSHKeygenEndToEnd)$' -count=1 -v
go test ./host/broker/ ./cmd/ailang-worldd/ ./host/approval/... ./host/store/ ./host/authority/ ./host/sshsign/ ./host/boundary/ -count=1
go test ./host/store/ -run '^TestProductionContextRoots$' -count=1
go test ./host/daemon/ -run '^TestDaemonDependencyAllowlist$' -count=1
go test ./host/verifygate/ -count=1        # repo-wide inventories (subprocess sink gate, fork-locked writes)
```
Measured in scratch on that tree: gofmt empty; `go vet ./...` rc 0; compile fence rc 0; the AC10 census
re-derives **13 sites in 13 files** (`host/sshsign/sshsign.go:47`) and passes with `host/sshsign/…/driver0`;
`host/broker` ok 264.8 s (whole package, incl. the load-sensitive `TestWorktreeSnapshotRefusesAConfiguredFilter`),
`cmd/ailang-worldd` ok 135.9 s, `host/approval` 1.0 s, `approvaltest` 0.3 s, `host/store` 14.5 s,
`host/authority` 1.4 s, `host/sshsign` 1.3 s (and `-race`), `host/boundary` 6.3 s, `host/verifygate` 235.6 s;
`TestProductionContextRoots` and `TestDaemonDependencyAllowlist` ok; the real-tool leg: XC1, XS2, XS1 and
both XS1 subtests PASS, 0 SKIP. Other repo-wide inventories checked: `host/boundary`
(`TestWorldBoundaryDependencyAllowlist`, `TestEveryCommandDirectoryIsAProtectedGroup`), `host/verifygate`
(`TestSubprocessSinkGate`, the fork-locked-write gates), `host/store` (`TestProductionContextRoots`,
`TestNoNewDeadlineFreeStoreReads`); no package list in `scripts/verify_go.sh` or `ci.yml` names peer
host packages (`grep -n 'host/toolprobe\|host/checkrun'` → none), so none needs `host/sshsign` added.

**Expected cumulative numstat after the controller's commit of S4: 1014** (`git diff --numstat 1fce6999
HEAD -- cmd host`; from 909 after S3 with §7.1's copied test: +154 host, +61 cmd; step 215).

**S4 mutants** (harness §4, `tools/mutants_s4.json`, run against the S4 tree; logs
`~/.ailang/state/world-iter378/logs/mutants_s4*.log`): 18 rows, 17 KILLED, 1 named equivalent SURVIVED.

| Id | Mutant (single edit, in `host/sshsign/sshsign.go` unless named) | Killing test (arm) | Result |
|---|---|---|---|
| MUT-ENV-SCRUB | the `cmd.Env = childenv.Scrubbed(os.Environ())` line dropped (written `_, _ = childenv.Scrubbed, os.Environ` so it compiles: the child inherits everything) | `host/broker` `TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential` (`host/sshsign/sshsign.go/driver0`: "the child observed AILANG_REGISTRY_API_KEY"); control: unmutated → ok | KILLED |
| MUT-ENV-SCRUB-UNIT | the same edit | `host/sshsign` `TestSignHandsTheSignerAScrubbedEnvironment` | KILLED |
| MUT-ENV-EMPTY | `Scrubbed(os.Environ())` → `Scrubbed(os.Environ()[:0])` (an empty env: no agent, no HOME) | `TestSignHandsTheSignerAScrubbedEnvironment` | KILLED |
| MUT-G1 / G2 / G3 | each check → `false && …` | cmd `TestSignWithSSHKeygenRefusesBadOutput` (failed / over-16k / exec-line) | KILLED ×3 |
| MUT-CAP-MAX | `max: MaxSignatureBytes + 1` | cmd `TestSignWithSSHKeygenRefusesBadOutput` (over-16k) | KILLED |
| MUT-CAP-WRITE | `c.buf.Write(p)` | `host/sshsign` `TestCappedBufferKeepsItsBound` | KILLED |
| MUT-ARGS-NS / -NS-REAL | `-n "git"+approval.Namespace[:0]` | cmd `TestSignWithSSHKeygenRefusesBadOutput` (exec-line) / `TestSignWithSSHKeygenEndToEnd` (real tool) | KILLED ×2 |
| MUT-ARGS-KEY | `-f key+".pub"` | cmd `TestSignWithSSHKeygenRefusesBadOutput` (exec-line) | KILLED |
| MUT-STDIN | stdin `msg[:0]` | same (exec-line) | KILLED |
| MUT-STDERR | stderr → `io.Discard` | same (failed) | KILLED |
| MUT-WAITDELAY | `cmd.WaitDelay = 0` | same (outlives-the-request, 25.9 s) | KILLED |
| MUT-CMD-CTX | `exec.CommandContext(context.WithoutCancel(ctx), bin, …)` | same (outlives-the-request, 24.0 s) | KILLED |
| MUT-SIGNER-BIN | `approve_sig.go`: `sshsign.Sign(ctx, "ssh-keygen-x", …)` | same (exec-line: no such program) | KILLED |
| MUT-SIGNER-ARGS | `approve_sig.go`: the message `msg[1:]` | same (exec-line: stdin) | KILLED |
| MUT-CAP-OVER-EQUIV | `c.over = n < len(p)` (equivalent, §4) | `TestCappedBufferKeepsItsBound` | SURVIVED (as predicted) |

The first harness pass recorded five NOCOMPILE drafts (a dropped line or literal left an import unused);
each was rewritten into the compiling form above and re-run, and only the re-run counts.

---

## 4. Mutation table (68 rows run in the planner's scratch: 66 KILLED, 2 named equivalents SURVIVED as predicted; plus 4 base controls, §1.2 C13)

**Derivation.** One row per new production `if`, guard, SQL clause, copied field, argument and dispatch in
Appendix P's production hunks; the sprint-3 judge's four gap mutants (with SIG-NOW's control); sprint 3's
F1/F2 rows re-derived for the count form (S4-D6) and MUT-F11 re-anchored in `approveArmored`; two
real-tool duplicates (the same edit, killed by XS1). **Procedure**: the harness
`~/.ailang/state/world-iter378/tools/mutate.py <tree> tools/mutants.json [ids…]` (outside the repo):
assert the anchor occurs exactly once, apply **one** edit, assert the file changed, `go vet` the package
(a mutant that does not compile is not a mutant — none did), run the command, restore, assert the
restore is byte-identical. Run under `WORLD_REQUIRE_SSH_KEYGEN=1` and `export
AILANG_BIN=$HOME/.pinned-ailang/ailang`. **Never commit or snapshot a mutant.**

Command keys: **C** = `go test ./cmd/ailang-worldd/ -run '^<test>$' -count=1`; **S** = `./host/store/`;
**U** = `./host/authority/`; **A** = `./host/approval/approvaltest/`. Result = the planner's run against
the final S3 tree (`~/.ailang/state/world-iter378/logs/mutants.log`).

| Id | Stage | Mutant (single edit) | Killing test (arm) | Cmd | Scratch |
|---|---|---|---|---|---|
| MUT-SIG-NOW | T1 | `approve_sig.go` the law's `Admit.Now: now` → `row.CreatedAt` (SIG-NOW) | `TestApproveBySignatureJudgesTheApproverAtNow` (enrolled at +5, approved at +10) | C | KILLED |
| MUT-SIG-NOW-CONTROL | T1 | *control*: the same line's `ExpiresAt: row.ExpiresAt` → `row.CreatedAt` | `TestApproveBySignatureJudgesTheApproverAtNow` (control: P1 expired) | C | KILLED |
| MUT-SIG-FORMAT | T1 | `Format: enrolled.Format` → `approval.FormatSSH` (SIG-FORMAT) | `TestApproveBySignatureRefusalsBeforeAnyEffect` (format-not-ssh) | C | KILLED |
| MUT-A4-EXP | T1 | `mint.go` A4 `Binds(…, rm.ExpiresAt)` → `rm.RequestCreatedAt` (A4-EXP) | `TestMintForRequestRefusesAnUnboundApproval` (bound-control) | U | KILLED |
| MUT-REVOKE-ACTIVE | T1 | `approvers.go` revoke SQL drops `AND revoked_at IS NULL` (item 4) | `TestRevokeApprover` (step 3, already revoked (exists at the base)) | S | KILLED |
| MUT-DM-DEADLINE | D1 | `requireDeadline` → `false && err != nil` | `TestDecidedMintRequests` (no-deadline) | S | KILLED |
| MUT-DM-VERSION | D1 | the version check → `false && …` | `TestDecidedMintRequests` (v3) | S | KILLED |
| MUT-DM-VERSION5 | D1 | `< mintRequestsSchemaVersion` → `< approversSchemaVersion` | `TestDecidedMintRequests` (v4) | S | KILLED |
| MUT-DM-LIMIT-LOW | D1 | `limit <= 0 \|\|` → `false && limit <= 0 \|\|` | `TestDecidedMintRequests` (limit 0 (LIMIT 0)) | S | KILLED |
| MUT-DM-LIMIT-HIGH | D1 | `\|\| limit > Max…` → `\|\| false && …` | `TestDecidedMintRequests` (cap, limit 65) | S | KILLED |
| MUT-DM-CAP | D1 | `MaxDecidedMintRequests = 64` → `65` | `TestDecidedMintRequests` (the cap literal) | S | KILLED |
| MUT-DM-FILTER | D1 | `WHERE consumed_at IS NOT NULL` → `WHERE 1` | `TestDecidedMintRequests` (the pending row) | S | KILLED |
| MUT-DM-ORDER | D1 | `consumed_at DESC` → `ASC` | `TestDecidedMintRequests` (order) | S | KILLED |
| MUT-DM-TIE | D1 | `request_id DESC` → `ASC` | `TestDecidedMintRequests` (r3 before r2) | S | KILLED |
| MUT-DM-TIE-DROP | D1 | `, request_id DESC` dropped | `TestDecidedMintRequests` (r3 before r2) | S | KILLED |
| MUT-DM-APPROVEDBY | D1 | `if approvedBy != ""` → `if false && …` | `TestDecidedMintRequests` (ssh:x) | S | KILLED |
| MUT-DM-APPROVEDBY-SQL | D1 | `approved_by = ?` → `approved_by <> ?` | `TestDecidedMintRequests` (ssh:x) | S | KILLED |
| MUT-F3 | D2 | F3 `if` → `if false && …` | `TestApproveFlagRules` (decided-without-list) | C | KILLED |
| MUT-F3-LIST | D2 | F3 → `if *decided` | `TestApproveDecidedList` (the list) | C | KILLED |
| MUT-DECIDED-BRANCH | D2 | `if *list && *decided` → `if *list && false && …` | `TestApproveDecidedList` (the list) | C | KILLED |
| MUT-DECIDED-FILTER | D2 | `DecidedMintRequests(ctx, "", 0)` → `(ctx, "tty", 0)` | `TestApproveDecidedList` (the list) | C | KILLED |
| MUT-F15 | D2 | F15 `if` → `if false && …` | `TestApproveDecidedList` (none) | C | KILLED |
| MUT-F16-BY | D2 | F16 prints the episode for `approved_by` | `TestApproveDecidedList` (the list) | C | KILLED |
| MUT-F16-AT | D2 | F16 prints `CreatedAt` for `ConsumedAt` | `TestApproveDecidedList` (the list) | C | KILLED |
| MUT-F16-CRED | D2 | `credentialPrefix` 16 → 17 | `TestApproveDecidedList` (the list) | C | KILLED |
| MUT-K8-FILTER | D3 | K8 `DecidedMintRequests(list, row.Format+":"+name, 0)` → `(list, "", 0)` | `TestApproverRevokeListsWhatItApproved` (mark-mac (the tty row appears)) | C | KILLED |
| MUT-K8-ROOT | D3 | K8's root `approveReadBudget` → `0` (an expired root) | `TestApproverRevokeListsWhatItApproved` (mark-mac (exitFatal)) | C | KILLED |
| MUT-K9 | D3 | K9 `if` → `if false && …` | `TestApproverRevokeListsWhatItApproved` (idle) | C | KILLED |
| MUT-K8-LINE | D3 | the `session revoke` line takes the 16-hex prefix | `TestApproverRevokeListsWhatItApproved` (the printed line) | C | KILLED |
| MUT-K8-AT | D3 | K8 prints `CreatedAt` for `ConsumedAt` | `TestApproverRevokeListsWhatItApproved` (mark-mac) | C | KILLED |
| MUT-K8-FORMAT-EQUIV | D3 | K8 filters `"ssh:"+name` (equivalent: P3 admits only `ssh`) | `TestApproverRevokeListsWhatItApproved` (equivalent) | C | SURVIVED |
| MUT-F1-COUNT | S1 | F1 `modes > 1` → `modes > 2` | `TestApproveFlagRules` (signature-and-*) | C | KILLED |
| MUT-F1-INC | S1 | `if v != ""` (the count) → `if false && …` | `TestApproveFlagRules` (every arm) | C | KILLED |
| MUT-F1-SIGNWITH | S1 | the count drops `*signWith` | `TestApproveFlagRules` (*-sign-with) | C | KILLED |
| MUT-F1-SIG | S1 | the count drops `*sigPath` | `TestApproveFlagRules` (signature-and-*, list-and-signature) | C | KILLED |
| MUT-F1-MSG | S1 | the count drops `*msgPath` | `TestApproveFlagRules` (*-write-message) | C | KILLED |
| MUT-F2 | S1 | F2 `if` → `if false && …` | `TestApproveFlagRules` (list-and-*) | C | KILLED |
| MUT-F7 | S1 | F7 `if` → `if false && …` | `TestApproveSignWithFakeSigner` (no-signer (nil call)) | C | KILLED |
| MUT-F7-ALWAYS | S1 | F7 → `if env.sign == nil` (any approve) | `TestApproveBySignatureMintsWithNoTerminal` (approveEnv has no signer) | C | KILLED |
| MUT-SW-BRANCH | S1 | `if *signWith != ""` → `if false && …` | `TestApproveSignWithFakeSigner` (good (typed path)) | C | KILLED |
| MUT-SW-REFUSAL | S1 | `signForApproval`'s refusal ignored | `TestApproveSignWithFakeSigner` (signer-error) | C | KILLED |
| MUT-SW-VERIFY | S1 | `approveArmored`'s refusal ignored (the signer's output trusted) | `TestApproveSignWithFakeSigner` (other-approver-other-message) | C | KILLED |
| MUT-SW-SOURCE | S1 | the source label → `"--sign-with "+key` | `TestApproveSignWithFakeSigner` (garbage) | C | KILLED |
| MUT-F11 | S1 | (sprint 3, re-anchored in `approveArmored`) F11 → `false && err != nil` | `TestApproveSignWithFakeSigner` (garbage (nil deref)) | C | KILLED |
| MUT-F7A | S1 | F7a `rem <= 0` → `rem < 0` | `TestApproveSignWithDeadline` (expired-since-admission) | C | KILLED |
| MUT-F7A-OFF | S1 | F7a → `if false && …` | `TestApproveSignWithDeadline` (expired-since-admission) | C | KILLED |
| MUT-REM-READING | S1 | `rem := ExpiresAt - env.now()` → `- row.CreatedAt` (the bound from the lifetime) | `TestApproveSignWithDeadline` (37-seconds-left) | C | KILLED |
| MUT-REM-BOUND | S1 | `time.Duration(rem)` → `time.Duration(rem+1)` | `TestApproveSignWithDeadline` (37-seconds-left, 1-second-left) | C | KILLED |
| MUT-SIGN-MSG | S1 | the seam gets `msg[1:]` | `TestApproveSignWithFakeSigner` (good) | C | KILLED |
| MUT-SIGN-KEY | S1 | the seam gets `"k"` for the key | `TestApproveSignWithFakeSigner` (good) | C | KILLED |
| MUT-SIGN-CTX | S1 | the seam gets `context.WithoutCancel(ctx)` (no deadline) | `TestApproveSignWithDeadline` (37-seconds-left) | C | KILLED |
| MUT-F8 | S1 | F8 → `if false && err != nil` | `TestApproveSignWithFakeSigner` (signer-error) | C | KILLED |
| MUT-G1 | S3 | G1 → `false && err != nil` | `TestSignWithSSHKeygenRefusesBadOutput` (failed) | C | KILLED |
| MUT-G1-REAL | S3 | the same G1 edit, against the real tool | `TestSignWithSSHKeygenEndToEnd` (k.pub, no private key) | C | KILLED |
| MUT-G2 | S3 | G2 → `if false && out.over` | `TestSignWithSSHKeygenRefusesBadOutput` (over-16k) | C | KILLED |
| MUT-G3 | S3 | G3 → `if false && …` | `TestSignWithSSHKeygenRefusesBadOutput` (exec-line-and-no-output) | C | KILLED |
| MUT-CAP-MAX | S3 | `cappedBuffer{max: …}` → `… + 1` | `TestSignWithSSHKeygenRefusesBadOutput` (over-16k) | C | KILLED |
| MUT-CAP-WRITE | S3 | `c.buf.Write(p[:n])` → `c.buf.Write(p)` (unbounded memory) | `TestCappedBufferKeepsItsBound` (—) | C | KILLED |
| MUT-CAP-OVER-EQUIV | S3 | `c.over = c.over \|\| n < len(p)` → `c.over = n < len(p)` (equivalent: a full buffer overflows on every later non-empty write) | `TestCappedBufferKeepsItsBound` (equivalent: after an overflow the buffer is full, so every later non-empty write overflows too) | C | SURVIVED |
| MUT-ARGS-NS | S3 | `-n approval.Namespace` → `-n git` | `TestSignWithSSHKeygenRefusesBadOutput` (exec-line-and-no-output) | C | KILLED |
| MUT-ARGS-NS-REAL | S3 | the same, against the real tool | `TestSignWithSSHKeygenEndToEnd` (mints) | C | KILLED |
| MUT-ARGS-KEY | S3 | `-f key` → `-f key+".pub"` | `TestSignWithSSHKeygenRefusesBadOutput` (exec-line-and-no-output) | C | KILLED |
| MUT-STDIN | S3 | stdin `msg` → `msg[:0]` | `TestSignWithSSHKeygenRefusesBadOutput` (exec-line-and-no-output) | C | KILLED |
| MUT-STDERR | S3 | the signer's stderr → `io.Discard` | `TestSignWithSSHKeygenRefusesBadOutput` (failed) | C | KILLED |
| MUT-WAITDELAY | S3 | `cmd.WaitDelay = time.Second` → `0` | `TestSignWithSSHKeygenRefusesBadOutput` (outlives-the-request (20 s)) | C | KILLED |
| MUT-CMD-CTX | S3 | `exec.CommandContext(ctx, …)` → `context.WithoutCancel(ctx)` | `TestSignWithSSHKeygenRefusesBadOutput` (outlives-the-request (20 s)) | C | KILLED |
| MUT-REAL-ENV | S3 | `realSessionEnv`'s `sign: signWithSSHKeygen` → `nil` | `TestRealSessionEnvSignsWithSSHKeygen` (—) | C | KILLED |
| MUT-TS-NAMESPACE | S3 | `approvaltest/signer.go` signs over `namespace+"x"` (self-consistent signer drift) | `TestTestSignerAgreesWithSSHKeygen` (check-novalidate) | A | KILLED |

**Base controls** (§1.2 C13, the pristine `1fce6999` tree): MUT-SIG-NOW (`-run '^TestApprove'`)
SURVIVED; MUT-SIG-FORMAT (`TestApproveBySignatureRefusalsBeforeAnyEffect`) SURVIVED; MUT-A4-EXP
(`TestMintForRequestRefusesAnUnboundApproval`) SURVIVED; MUT-REVOKE-ACTIVE (`TestRevokeApprover`) KILLED.

**Known-equivalent (named, not counted; each run and SURVIVED as predicted):**
- **MUT-K8-FORMAT-EQUIV** (D3): K8 filters `"ssh:"+name`. Equivalent because the signature law's P3
  admits only `format = ssh`, so no other `approved_by` prefix is ever written.
- **MUT-CAP-OVER-EQUIV** (S3): `c.over = n < len(p)`. Equivalent because after an overflow the buffer is
  full, so every later non-empty write overflows too (only a zero-length write could reset it; `os/exec`'s
  copy never makes one).

**Propagation-only branches (named, not counted):** `DecidedMintRequests`' `QueryContext`, scan and
`rows.Err` errors; `approveListDecided`'s read error; K8's read-only open and read errors (`exitFatal`);
`signForApproval`'s `approvalMessage` error (unreachable after admission); the message-only label of
`--signature`'s F11 text.

---

## 5. Files (final `git status --short` after S3; after D3 if the gate stops the sprint)

| File | Sub-milestone | Change |
|---|---|---|
| `cmd/ailang-worldd/approve_sig_test.go` | T1, D2, S1 | CR4 `format-not-ssh`, SIG-NOW test; CR1 `decided-without-list`; CR1 three `--sign-with` arms |
| `host/authority/mint_derive_test.go` | T1 | AU2 `bound-control` |
| `host/store/mint_requests.go` | D1 | `MaxDecidedMintRequests`, `DecidedMintRequests` |
| `host/store/decided_requests_test.go` | D1 | new: CS6 |
| `cmd/ailang-worldd/approve.go` | D2, S1 | `--decided`, F3, F15, F16; `--sign-with`, F1/F2 count, F7, `expiredRefusal` |
| `cmd/ailang-worldd/approve_decided_test.go` | D2, D3 | new: CR7; CR12 |
| `cmd/ailang-worldd/main.go`, `help_test.go` | D2, S1 | usage lines and wants |
| `cmd/ailang-worldd/approver.go` | D3, S1 | K8, K9, help; help names `--sign-with` |
| `cmd/ailang-worldd/approver_revoke_test.go` | D3 | `yes` stdout + K9 |
| `host/store/context_roots_test.go` | D3, S1 | `runApproverRevoke` 2 → 3; `signForApproval` 1 |
| `cmd/ailang-worldd/approve_sig.go` | S1, S3 | `approveArmored`, `approvalMessage`, `signForApproval`; `signWithSSHKeygen`, `cappedBuffer` |
| `cmd/ailang-worldd/session.go` | S1, S3 | `sessionEnv.sign`; `realSessionEnv.sign` |
| `cmd/ailang-worldd/approve_signwith_test.go` | S2, S3 | new: CR8, CR9; G1–G3, outlives, `cappedBuffer`, wiring, XS1 |
| `host/approval/approvaltest/sign_crosscheck_test.go` | S3 | new: XS2 |
| `.github/workflows/ci.yml` | S3 | the row 201 step only |
| `docs/QUICKSTART.md` | S3 | §15 (attended-pending) |
| this plan and `sprint_w-world-approvals-by-key-a1b-s4.json` | planner | |

---

## 6. Sprint gates (controller, after the last sub-milestone; each green at the base, §1.1)

In **one** `bash` shell, from the worktree:

```
set -euo pipefail
export AILANG_BIN=$HOME/.pinned-ailang/ailang
./scripts/verify_ail.sh                                                    # no .ail change: exactly green
test -z "$(gofmt -l cmd/ailang-worldd host/approval host/store host/authority)"
go build ./... && go vet ./...
go test -run '^$' ./...                                                    # compile fence for every _test.go
go test ./cmd/ailang-worldd/ ./host/approval/... ./host/store/ ./host/authority/ -count=1
go test ./host/store/ -run '^TestProductionContextRoots$' -count=1         # the two edited pins
go test ./host/daemon/ -run '^TestDaemonDependencyAllowlist$' -count=1
go test -race ./host/approval/... ./host/store/ ./host/authority/ -count=1
go test -race ./cmd/ailang-worldd/ -run '^(TestApproveDecidedList|TestApproverRevokeListsWhatItApproved|TestApproveSignWith|TestSignWithSSHKeygen|TestCappedBuffer|TestApproveBySignatureJudgesTheApproverAtNow|TestApproveFlagRules|TestApproverRevoke)' -count=3
WORLD_REQUIRE_SSH_KEYGEN=1 go test ./host/approval/... ./cmd/ailang-worldd/ -count=1 -v \
  -run '^(TestSSHSIGCrossCheckWithSSHKeygen|TestTestSignerAgreesWithSSHKeygen|TestSignWithSSHKeygenEndToEnd)$' 2>&1 | tee /tmp/w378-xs.log | tail -3
for t in TestSSHSIGCrossCheckWithSSHKeygen TestTestSignerAgreesWithSSHKeygen TestSignWithSSHKeygenEndToEnd/ed25519 TestSignWithSSHKeygenEndToEnd/ecdsa-p256; do
  grep -qE -- "^ *--- PASS: $t " /tmp/w378-xs.log || { echo "$t did not PASS"; exit 1; }; done
go test ./cmd/ailang-worldd/ -count=1 -v -run '^(TestApprove|TestApprover|TestEpisodeStartRequest|TestSessionMint_|TestSessionNew|TestEveryVerbHasHelp|TestHelpIsNeverBlockedByTheRegistryKeyGuard|TestSignWith|TestCappedBuffer|TestRealSessionEnv)' 2>&1 | tee /tmp/w378-cl.log | tail -3
if grep -E '^\s*--- (SKIP|FAIL)' /tmp/w378-cl.log; then exit 1; fi
# AFTER the controller commits (the tracked-binary gate reads committed blobs):
VERIFY_GO_LEGS=plain ./scripts/verify_go.sh
git status --short                                                         # only the §5 files
```

If the gate stopped the sprint after D3: drop the `TestApproveSignWith|TestSignWithSSHKeygen|TestCappedBuffer`
names, the `WORLD_REQUIRE_SSH_KEYGEN` leg and `TestSignWith|TestCappedBuffer|TestRealSessionEnv` from the
`-v` run. **Known base reds outside this sprint** (§1.1) are the base; any other red, any race report, or
any red in a touched package is a regression. CI's `go-race-verify` and the widened row 201 step are the
verdicts.

**Must stay green, unedited**: every sprint-3 `TestApprove*` (CR2–CR6, CR10, CR11), every a0a1
`TestApprove*`/`TestEpisodeStartRequest*`, `TestApproveListShowsPending`, `TestApproveListUsesTheNamedStores`,
`TestApproverVerbsUseTheNamedStores`, `TestApproverRevokeUnansweredPromptTimesOut`, `TestMintRequestFieldMapIsComplete`,
`TestSSHSIG*`, `TestCredentials*`, `TestMintSessionForRequest*`, `TestRevokeApprover`,
`TestDaemonDependencyAllowlist`, `TestApprovalTestIsTestOnly`, the QUICKSTART tests. **The only edited
existing tests**: `TestApproveBySignatureRefusalsBeforeAnyEffect` (+1 arm), `TestApproveFlagRules` (+4
arms), `TestMintForRequestRefusesAnUnboundApproval` (+1 subtest), `TestApproverRevoke/yes` (K9 line),
`TestEveryVerbHasHelp` (+2 wants, one per part), `TestProductionContextRoots` (two pins).

---

## 7. Executor directive notes

- **No git writes.** Never `git add`, `commit`, `stash`, `checkout`, `reset`, `restore`, `rebase`, `push`,
  `worktree` or `apply`. `patch -p1` (§2) is the only writer; read-only `git status`/`git diff` are fine.
- **Runs.** Run 1 = T1, D1; Run 2 = D2, D3; **gate**; Run 3 = S1, S2; Run 4 = S3. Each run reads §0, §2
  and its own §3 sections, then applies its patches in order.
- **Expected numstat after each sub-milestone, measured AFTER the controller's commit of it** (an
  uncommitted tree hides the new untracked files from numstat):
  ```
  bash -c 'cd /Users/voightkampff/dev/sunholo-data/.wt-world-iter378-row201-a1b-s4 && git diff --numstat 1fce6999ebf875299254a1d77003e9e3351a3347 HEAD -- cmd host | awk "{s+=\$1+\$2} END {print s}"'
  ```
  After the controller's commit of T1 **33**, D1 **181**, D2 **278**, D3 **359**, S1 **461**, S2 **615**,
  S3 **896**. **Any other value means a patch was not applied as written: stop and report.**
- **The size gate (controller), after the D3 commit:** the numstat above is **≤ 1,463** → run Runs 3 and 4
  (all of S1, S2, S3). Over → do not run them; S1–S3 move to the next sprint unchanged (their patches are
  final against the D3 tree) and the sprint closes after D3. Expected: 359 (1,104 lines of margin).
- **Snapshot after every sub-milestone, cumulatively**:
  ```
  bash -c 'set -euo pipefail; M=T1; W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter378-row201-a1b-s4
    D=$HOME/.ailang/state/world-iter378/snap/$M; rm -rf "$D"; mkdir -p "$D"; cd "$W"
    git status --porcelain --untracked-files=all | while read -r st f; do
      case "$f" in design_docs/planned/w-world-approvals-by-key-a1b-sprint4-plan.md|design_docs/planned/sprint_w-world-approvals-by-key-a1b-s4.json) continue;; esac
      mkdir -p "$D/$(dirname "$f")"; cp -p "$f" "$D/$f"; done
    (cd "$D" && find . -type f ! -name MANIFEST.sha256 | sort | xargs shasum -a 256) > "$D/MANIFEST.sha256"
    cat "$D/MANIFEST.sha256"'
  ```
  (`M` = `T1` … `S3`.) No file is deleted, so the porcelain status has only `M` and `??`; stop on any other
  status. The controller rebuilds one commit per sub-milestone (iteration 377's `commit_ms.sh` with `WT`
  and the snapshot root pointed here).
- **Stop conditions.** A `patch --dry-run` failure; a numstat other than the expected value; a red
  earlier-milestone test; a `--- SKIP` in a `-v` run of `cmd/ailang-worldd` or of the
  `WORLD_REQUIRE_SSH_KEYGEN=1` leg; any need to touch a file outside §5; the executor running
  `ssh-keygen` itself (only the named Go tests do, §0.1); any reference to `~/.ssh`, `~/.ailang/world` or
  a real `/dev/tty`. Stop and report; do not work around. **Do not run the QUICKSTART §15 commands**
  (attended).
- **Timing notes.** `outlives-the-request` takes ~2 s and leaves a `sleep 20` child of the fake script
  running for up to 20 s after the test (harmless; it holds nothing). The real-tool tests take < 1 s.
- **Shell.** Wrap loops, globs and `$VAR` expansions in `bash -c` or a `bash <<'EOF'` heredoc (zsh does not
  word-split). Export `AILANG_BIN` in the same shell as the Go tests. Poll any background command's
  artifact in a bounded `date +%s` loop in the same turn; never end a turn announcing a wait.

---

### 7.1 Controller correction after Run 4 (iteration 378)

The S3 patch in Appendix P omits one test the planner compiled and mutation-tested in its scratch tree:
`TestCappedBufferKeepsItsBound` (13 lines in `cmd/ailang-worldd/approve_signwith_test.go`), the §4 killer of
MUT-CAP-WRITE. The executor found it (the §3 `-run` regex matched nothing for that name and still exited 0).
The controller compared all 17 sprint files with the planner's tree `~/.ailang/state/world-iter378/work/`:
16 were byte-identical and this file differed only by that function. It was copied from the planner's tree,
not rewritten. MUT-CAP-WRITE (`c.buf.Write(p[:n])` → `c.buf.Write(p)`) was re-run by the controller against
that single test: `--- FAIL` (killed), then the source was restored byte-identical. Appendix P is left as
generated; **the expected numstat after S3 is therefore 909, not 896.**

**Revision S4 (planner, after the controller's §6 sweep at `784c9413`):** Run 5 = S4 (§3 S4, Appendix P
`diff S4`), applied with §2 on top of `784c9413`. Expected numstat after the controller's commit of S4:
**1014**. §6's gate list gains `go test ./host/broker/ ./host/sshsign/ ./host/verifygate/ -count=1` (the
repo-wide AC10 census lives in `host/broker`; the planner's original list missed it).

## 8. Risks

- **R1 — Size.** Measured 896 against a 2,000 cap; the committed part 359. The per-step expected numstat
  makes drift visible before the tail runs.
- **R2 — R6 (SSHSIG proves possession, not presence)** is unchanged, and `--sign-with` makes it
  sharper: anything that can run `ssh-keygen` with the operator's agent can sign. The help says a `.pub`
  path does not prove the agent signed. Q5's attended exit run is Mark's.
- **R3 — `PATH` is not a trust boundary** (design §2.9): a substituted `ssh-keygen` can only refuse, or
  produce bytes that `approveArmored` verifies against the enrolled key over World's own message
  (MUT-SW-VERIFY is killed). Its output is capped at 16 KiB in memory (MUT-CAP-WRITE) and its run at the
  request's remaining lifetime plus 1 s (`WaitDelay`).
- **R4 — A passphrase prompt.** `ssh-keygen` prompts on `/dev/tty` itself; World's deadline (≤ the
  request's remaining lifetime) bounds it, and the kill is `signal: killed` → `exitUsage`. Not run by the
  loop (attended).
- **R5 — CI's ssh-keygen and compile time.** The step already required `ssh-keygen` and passed at
  `1fce6999`; the new tests use the same flags as XC1. The timeout is raised for the cold compile (P6).
- **R6 — K8 shows at most the newest 64** sessions an approver approved. A session lives an hour, so older
  rows are expired; the listing does not say it was capped (residual, Appendix C).
- **R7 — Haiku and patches.** The executor transcribes nothing; a skipped patch fails `--dry-run` or the
  numstat; a hand "fix" after a red test is forbidden by §7.

---

## 9. Acceptance criteria (the judge scores against these)

1. T1: SIG-NOW, SIG-FORMAT and A4-EXP, which survive their tests at the base (C13), are each killed by
   one named test; the revoke SQL mutant is killed (at the base already).
2. `DecidedMintRequests` lists consumed rows only, newest `consumed_at` first, `request_id DESC` on ties,
   every column, at most 64, filtered by `approved_by` when asked, on schema 4 and above, refusing below
   4 and without a deadline (CS6).
3. `approve --list --decided` prints exactly F16's lines (typed `tty` and signed `ssh:<name>`, newest
   first) or F15; `--decided` without `--list` is the help and `exitUsage` (F3).
4. `approver revoke` lists exactly the sessions that approver approved, each with a `session revoke` line
   that works, revokes none of them (Q4), says K9 when there are none, and reads under its own root after
   the write's has ended (`TestProductionContextRoots` pin 3).
5. If the gate ran S1–S3: `--sign-with` signs exactly the request's message under a deadline equal to
   the request's remaining lifetime, refuses before signing a request that expired after admission (F7a),
   refuses a signer error (F8) and a missing signer (F7, `exitFatal`), judges the signer's output exactly as
   `--signature` (another approver's signature, garbage), mints with no terminal as `ssh:<name>`; F1/F2 hold
   for all three flags; the production signer runs `ssh-keygen -Y sign -n world-approval@sunholo -f <key>`
   with the message on stdin and refuses a failed run, over 16 KiB, no output and a run outliving the
   request; with the real tool it mints for ed25519 and ecdsa-p256 keys and refuses a `.pub` with no private
   key (XS1); the test signer agrees with `ssh-keygen` (XS2); the CI step runs and greps all three.
6. Each §4 mutant, applied alone, makes its named test fail; the named equivalents survive for the stated
   reason.
7. Every §6 gate passes, including `VERIFY_GO_LEGS=plain ./scripts/verify_go.sh` after the commit; no
   test skips; `git status --short` lists only §5.
8. The numstat after each sub-milestone's commit equals §7's value.
9. No real key, real store, `~/.ssh` access, real agent, mint outside a temp store, or publish anywhere in
   the sprint (D-WORLD-109); the executor runs no `ssh-keygen` itself.

---

## Appendix C — after this sprint (A1b is code-complete)

Not built here; none of it is a code milestone of the design:
- **The attended exit run** (design §8, D-WORLD-109): Mark enrols his own key, approves with `--sign-with`,
  reads `--list --decided`, revokes, and sees the revoked key refused; QUICKSTART §15 is run verbatim and
  its *attended — pending* tag removed.
- **Residuals to consider, not scheduled:** K8 / `--list --decided` do not say when the 64-row cap was hit
  (R6); design CR4's [r2] "size-grew" arm (a test reader seam for a file whose stream outgrows its
  `fstat` size) was not built in sprint 3 either — `readRegularFileCapped`'s `LimitReader` +1 check is
  sprint 2's and is covered by its own tests; design CR13 stays dropped (S3-D11).

---

## Appendix P — the patches (apply with §2; generated from the compiled scratch stages; do not edit)

### Patch T1 — the sprint-3 judge's test gaps

````diff T1
--- a/cmd/ailang-worldd/approve_sig_test.go
+++ b/cmd/ailang-worldd/approve_sig_test.go
@@ -203,6 +203,13 @@
 				t.Fatal(err)
 			}
 		}, sig: good, want: "approval: the clock is behind the approver's enrolment time"},
+		{name: "format-not-ssh", setup: func(t *testing.T, a approveRig) { // SIG-FORMAT (sprint 4 T1): the law reads the row's format
+			raw := rawCredentials(t, a.credPath)
+			defer func() { _ = raw.Close() }()
+			if _, err := raw.Exec(`UPDATE approvers SET format = 'webauthn' WHERE name = 'mark-mac'`); err != nil {
+				t.Fatal(err)
+			}
+		}, sig: good, want: "approval: approver format must be ssh"},
 		{name: "git-namespace", sig: func(t *testing.T, a approveRig) string {
 			return writeSignature(t, sigSigner.Sign("git", requestMessage(t, a)))
 		}, want: "approval: sshsig: wrong namespace"},
@@ -408,5 +415,24 @@
 				t.Fatalf("a credentials store was touched (err %v)", err)
 			}
 		})
+	}
+}
+
+// ---- T1 (sprint 4) ----
+
+// SIG-NOW (sprint 4 T1, the sprint-3 judge's gap): the signature law judges the approver at
+// approve's own clock reading, not at the request's creation: an approver enrolled after the
+// request was made and before the approve approves it.
+func TestApproveBySignatureJudgesTheApproverAtNow(t *testing.T) {
+	a := newApproveRig(t, "ep1", approveT0)
+	seedApprover(t, a.credPath, "mark-mac", sigSigner.PublicKey(), approveT0+5)
+	sig := writeSignature(t, sigSigner.Sign(approval.Namespace, requestMessage(t, a)))
+	var opened bool
+	var stdout, stderr bytes.Buffer
+	if code := runApprove(a.args("--signature", sig), &stdout, &stderr, approveEnv(nil, &opened, atTime(approveT0+10))); code != exitOK || opened {
+		t.Fatalf("rc %d, terminal opened %v; want exitOK with no terminal for an approver enrolled at +5 and an approve at +10\nstderr %s", code, opened, stderr.String())
 	}
+	if row := a.row(t); !row.Consumed || row.ApprovedBy != "ssh:mark-mac" {
+		t.Fatalf("request consumed %v by %q, want ssh:mark-mac", row.Consumed, row.ApprovedBy)
+	}
 }
--- a/host/authority/mint_derive_test.go
+++ b/host/authority/mint_derive_test.go
@@ -182,6 +182,13 @@
 		}
 	}
 	is := func(target error) func(error) bool { return func(err error) bool { return errors.Is(err, target) } }
+	t.Run("bound-control", func(t *testing.T) { // A4-EXP (sprint 4 T1): the unperturbed approval binds, so each arm below is refused for its own field
+		st, mr, row := au1Store(t)
+		au1Enrol(t, st)
+		if err := au1Mint(t, st, au1SignedMint(t, mr, row)); err != nil {
+			t.Fatalf("the bound approval: %v, want it minted", err)
+		}
+	})
 	for _, a := range []struct {
 		name   string
 		mutate func(t *testing.T, rm *RequestMint, row store.MintRequestRow)
````

### Patch D1 — `DecidedMintRequests`

````diff D1
--- a/host/store/decided_requests_test.go
+++ b/host/store/decided_requests_test.go
@@ -0,0 +1,108 @@
+package store
+
+// Row 201 A1b sprint 4 (plan design_docs/planned/w-world-approvals-by-key-a1b-sprint4-plan.md §3
+// D1): DecidedMintRequests, the consumed rows approve --list --decided and approver revoke's K8
+// read. Every store is under t.TempDir().
+
+import (
+	"context"
+	"database/sql"
+	"errors"
+	"fmt"
+	"path/filepath"
+	"reflect"
+	"strings"
+	"testing"
+)
+
+// putDecided inserts request n as consumed at `at` by approvedBy, with credential id "%064x" of n.
+func putDecided(t *testing.T, db *sql.DB, n int, at int64, approvedBy string) MintRequestRow {
+	t.Helper()
+	r := testMintRequestRow(n, at-10)
+	r.Consumed, r.ConsumedAt, r.ApprovedBy, r.CredentialID = true, at, approvedBy, fmt.Sprintf("%064x", n)
+	if _, err := db.Exec(`INSERT INTO mint_requests (`+mintRequestColumns+`) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
+		r.RequestID, r.Request, r.Digest, r.Challenge, r.EpisodeID, r.CreatedAt, r.ExpiresAt, at, approvedBy, r.CredentialID); err != nil {
+		t.Fatal(err)
+	}
+	return r
+}
+
+// CS6: consumed rows only, newest consumed_at first and request_id descending among equal times,
+// every column returned; approvedBy keeps only its own rows; the cap holds; schema 4 is read, and
+// a store below it is refused.
+func TestDecidedMintRequests(t *testing.T) {
+	if MaxDecidedMintRequests != 64 {
+		t.Fatalf("MaxDecidedMintRequests = %d, want 64 (design §2.4 [r2])", MaxDecidedMintRequests)
+	}
+	ctx := testCtx(t)
+	c := openTestCredentials(t)
+	if err := c.PutMintRequest(ctx, testMintRequestRow(100, 1000), 1000); err != nil { // pending: never listed
+		t.Fatal(err)
+	}
+	if got, err := c.DecidedMintRequests(ctx, "", 0); err != nil || len(got) != 0 {
+		t.Fatalf("no decided rows: (%+v, %v), want none", got, err)
+	}
+	r1 := putDecided(t, c.db, 1, 2000, "tty")
+	r2 := putDecided(t, c.db, 2, 3000, "ssh:x") // r2 and r3 share consumed_at: request_id DESC puts r3 first
+	r3 := putDecided(t, c.db, 3, 3000, "tty")
+	r4 := putDecided(t, c.db, 4, 1500, "ssh:x")
+	got, err := c.DecidedMintRequests(ctx, "", 0)
+	if err != nil {
+		t.Fatal(err)
+	}
+	if want := []MintRequestRow{r3, r2, r1, r4}; !reflect.DeepEqual(got, want) {
+		t.Fatalf("decided rows\n got %+v\nwant %+v", got, want)
+	}
+	if got, err := c.DecidedMintRequests(ctx, "ssh:x", 0); err != nil || !reflect.DeepEqual(got, []MintRequestRow{r2, r4}) {
+		t.Fatalf("approved by ssh:x: (%+v, %v), want r2 and r4", got, err)
+	}
+	if got, err := c.DecidedMintRequests(ctx, "ssh:nobody", 0); err != nil || len(got) != 0 {
+		t.Fatalf("approved by ssh:nobody: (%+v, %v), want none", got, err)
+	}
+	if got, err := c.DecidedMintRequests(ctx, "", 2); err != nil || !reflect.DeepEqual(got, []MintRequestRow{r3, r2}) {
+		t.Fatalf("limit 2: (%+v, %v), want r3 and r2", got, err)
+	}
+
+	t.Run("cap", func(t *testing.T) {
+		c := openTestCredentials(t)
+		for n := 1; n <= MaxDecidedMintRequests+6; n++ {
+			putDecided(t, c.db, n, int64(2000+n), "tty")
+		}
+		for _, limit := range []int{0, -1, MaxDecidedMintRequests + 1} {
+			if got, err := c.DecidedMintRequests(ctx, "", limit); err != nil || len(got) != MaxDecidedMintRequests {
+				t.Fatalf("limit %d: %d rows (%v), want %d", limit, len(got), err, MaxDecidedMintRequests)
+			}
+		}
+	})
+	t.Run("v4", func(t *testing.T) {
+		path := filepath.Join(t.TempDir(), CredentialsFileName)
+		makeV4Credentials(t, path, nil)
+		c4, err := OpenCredentials(path)
+		if err != nil {
+			t.Fatal(err)
+		}
+		defer func() { _ = c4.Close() }()
+		r := putDecided(t, c4.db, 7, 2000, "tty")
+		if got, err := c4.DecidedMintRequests(ctx, "", 0); err != nil || !reflect.DeepEqual(got, []MintRequestRow{r}) {
+			t.Fatalf("v4: (%+v, %v), want the one decided row", got, err)
+		}
+	})
+	t.Run("v3", func(t *testing.T) {
+		path := filepath.Join(t.TempDir(), CredentialsFileName)
+		makeV3Credentials(t, path, nil)
+		c3, err := OpenCredentials(path)
+		if err != nil {
+			t.Fatal(err)
+		}
+		defer func() { _ = c3.Close() }()
+		var older *OlderCredentialsStoreError
+		if _, err := c3.DecidedMintRequests(ctx, "", 0); !errors.As(err, &older) || older.Need != 4 || !strings.Contains(err.Error(), "restart `ailang-worldd serve`") {
+			t.Fatalf("v3: err %v, want the mint-requests refusal (need 4)", err)
+		}
+	})
+	t.Run("no-deadline", func(t *testing.T) {
+		if _, err := c.DecidedMintRequests(context.Background(), "", 0); !errors.Is(err, ErrNoDeadline) {
+			t.Fatalf("err %v, want ErrNoDeadline", err)
+		}
+	})
+}
--- a/host/store/mint_requests.go
+++ b/host/store/mint_requests.go
@@ -195,7 +195,47 @@
 	}
 	if err := rows.Err(); err != nil {
 		return nil, fmt.Errorf("store: list mint requests: %w", err)
+	}
+	return out, nil
+}
+
+// MaxDecidedMintRequests caps one DecidedMintRequests listing (row 201 A1b design §2.4 [r2]).
+const MaxDecidedMintRequests = 64
+
+// DecidedMintRequests lists the consumed requests, newest consumed_at first and request_id
+// descending among equal times (design §2.4 [r2]), at most MaxDecidedMintRequests (limit <= 0
+// or above it means the cap). A non-empty approvedBy keeps only the rows it approved
+// (approver revoke's K8). It reads every store with mint requests, schema 4 included.
+func (c *Credentials) DecidedMintRequests(ctx context.Context, approvedBy string, limit int) ([]MintRequestRow, error) {
+	if err := requireDeadline(ctx); err != nil {
+		return nil, err
+	}
+	if c.version < mintRequestsSchemaVersion {
+		return nil, &OlderCredentialsStoreError{Path: c.display, Version: c.version, Need: mintRequestsSchemaVersion, Feature: "mint requests"}
+	}
+	if limit <= 0 || limit > MaxDecidedMintRequests {
+		limit = MaxDecidedMintRequests
+	}
+	q, args := `SELECT `+mintRequestColumns+` FROM mint_requests WHERE consumed_at IS NOT NULL`, []any{}
+	if approvedBy != "" {
+		q, args = q+` AND approved_by = ?`, append(args, approvedBy)
 	}
+	rows, err := c.db.QueryContext(ctx, q+` ORDER BY consumed_at DESC, request_id DESC LIMIT ?`, append(args, limit)...)
+	if err != nil {
+		return nil, fmt.Errorf("store: list decided mint requests: %w", err)
+	}
+	defer func() { _ = rows.Close() }()
+	var out []MintRequestRow
+	for rows.Next() {
+		row, err := scanMintRequest(rows.Scan)
+		if err != nil {
+			return nil, fmt.Errorf("store: list decided mint requests: %w", err)
+		}
+		out = append(out, row)
+	}
+	if err := rows.Err(); err != nil {
+		return nil, fmt.Errorf("store: list decided mint requests: %w", err)
+	}
 	return out, nil
 }
 
````

### Patch D2 — `approve --list --decided`

````diff D2
--- a/cmd/ailang-worldd/approve.go
+++ b/cmd/ailang-worldd/approve.go
@@ -24,7 +24,7 @@
 const approveHelp = `usage: ailang-worldd approve <request-id> [--db <path>] [--credentials-db <path>]
        ailang-worldd approve <request-id> --signature <file> [--db <path>] [--credentials-db <path>]
        ailang-worldd approve <request-id> --write-message <file> [--db <path>] [--credentials-db <path>]
-       ailang-worldd approve --list [--db <path>] [--credentials-db <path>]
+       ailang-worldd approve --list [--decided] [--db <path>] [--credentials-db <path>]
 
 ATTENDED, or SIGNED. Mints the session a stored mint request describes (row
 201), after a y/N typed at /dev/tty, or, with --signature, after an enrolled
@@ -55,6 +55,9 @@
   and writes no store row: the request stays pending.
 --list prints the pending requests, read-only. A request expires 10 minutes
 after it is made; at most 16 are pending per credentials store.
+--list --decided prints the used requests instead, newest first (at most 64):
+  who approved each (tty, or ssh:<name>), when, and its credential's first 16
+  hex. 'ailang-worldd session revoke' takes the credential id.
 
   --db <path>               the world store the request was made for
                             (default ~/.ailang/world/world.db)
@@ -111,6 +114,7 @@
 	db := fs.String("db", def.db, "the world store the request was made for (default ~/.ailang/world/world.db)")
 	credsDB := fs.String("credentials-db", "", credentialsDBUsage)
 	list := fs.Bool("list", false, "print the pending requests, read-only")
+	decided := fs.Bool("decided", false, "with --list: print the used requests and who approved them")
 	sigPath := fs.String("signature", "", "approve with an enrolled approver's SSHSIG over the request (no terminal)")
 	msgPath := fs.String("write-message", "", "write the request's approval message to this new file, to sign elsewhere")
 	positional, err := parseInterleaved(fs, args)
@@ -134,6 +138,10 @@
 		fmt.Fprint(stderr, approveHelp)
 		return exitUsage
 	}
+	if *decided && !*list { // F3: --decided is a --list
+		fmt.Fprint(stderr, approveHelp)
+		return exitUsage
+	}
 	id := ""
 	if !*list {
 		id = positional[0]
@@ -156,6 +164,9 @@
 	if err != nil {
 		return refuse("cannot read the credentials store %s: %v", credPath, err)
 	}
+	if *list && *decided {
+		return approveListDecided(ctx, ro, credPath, stdout, refuse)
+	}
 	if *list {
 		rows, err := ro.PendingMintRequests(ctx, env.now(), 0)
 		_ = ro.Close()
@@ -261,3 +272,25 @@
 		task: mr.Task, grants: toCapabilities(mr.Grants)})
 	return exitOK
 }
+
+// approveListDecided is approve --list --decided (design §2.7 F15, F16): the used requests,
+// newest first, each with who approved it, when, and its credential id's first 16 hex. It closes ro.
+func approveListDecided(ctx context.Context, ro *store.Credentials, credPath string, stdout io.Writer, refuse func(string, ...any) int) int {
+	rows, err := ro.DecidedMintRequests(ctx, "", 0)
+	_ = ro.Close()
+	if err != nil {
+		return refuse("cannot list the used requests in %s: %v", credPath, err)
+	}
+	if len(rows) == 0 { // F15
+		fmt.Fprintln(stdout, "no decided mint requests")
+		return exitOK
+	}
+	for _, r := range rows { // F16
+		fmt.Fprintf(stdout, "%s  episode %s  approved by %s  at %s  credential %s\n", r.RequestID, r.EpisodeID, r.ApprovedBy,
+			approverTime(r.ConsumedAt), credentialPrefix(r.CredentialID))
+	}
+	return exitOK
+}
+
+// credentialPrefix is a credential id's first 16 characters (all of a shorter one).
+func credentialPrefix(id string) string { return id[:min(16, len(id))] }
--- a/cmd/ailang-worldd/approve_decided_test.go
+++ b/cmd/ailang-worldd/approve_decided_test.go
@@ -0,0 +1,56 @@
+package main
+
+// Row 201 A1b sprint 4 (plan design_docs/planned/w-world-approvals-by-key-a1b-sprint4-plan.md §3
+// D2, D3): approve --list --decided, and what approver revoke lists. Every store is under
+// t.TempDir(); every key is approvaltest's fixed seed (D-WORLD-109).
+
+import (
+	"bytes"
+	"strings"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+)
+
+// decidedRig is one store with three requests: a approved by the typed y at +1, b approved by
+// mark-mac's signature at +2, and c still pending.
+func decidedRig(t *testing.T) (a, b, c approveRig) {
+	t.Helper()
+	r := newEpisodeRun(t)
+	a = newApproveRigIn(t, r, "ep1", approveT0)
+	b = newApproveRigIn(t, r, "ep2", approveT0)
+	c = newApproveRigIn(t, r, "ep3", approveT0)
+	seedApprover(t, b.credPath, "mark-mac", sigSigner.PublicKey(), approveT0-100)
+	var stdout, stderr bytes.Buffer
+	if code := runApprove(a.args(), &stdout, &stderr, approveEnv(&promptTerm{in: strings.NewReader("y\n")}, nil, atTime(approveT0+1))); code != exitOK {
+		t.Fatalf("typed approve of a: rc %d\nstderr %s", code, stderr.String())
+	}
+	sig := writeSignature(t, sigSigner.Sign(approval.Namespace, requestMessage(t, b)))
+	if code := runApprove(b.args("--signature", sig), &stdout, &stderr, approveEnv(nil, nil, atTime(approveT0+2))); code != exitOK {
+		t.Fatalf("signed approve of b: rc %d\nstderr %s", code, stderr.String())
+	}
+	return a, b, c
+}
+
+// CR7: --list --decided prints the used requests only, newest first, each with who approved it,
+// when, and its credential's first 16 hex (F16); with none it says so (F15).
+func TestApproveDecidedList(t *testing.T) {
+	a, b, _ := decidedRig(t)
+	ar, br := a.row(t), b.row(t)
+	want := b.id + "  episode ep2  approved by ssh:mark-mac  at 2027-01-15T08:00:02Z  credential " + br.CredentialID[:16] + "\n" +
+		a.id + "  episode ep1  approved by tty  at 2027-01-15T08:00:01Z  credential " + ar.CredentialID[:16] + "\n"
+	var stdout, stderr bytes.Buffer
+	args := []string{"--list", "--decided", "--db", a.r.f.db, "--credentials-db", a.credPath}
+	if code := runApprove(args, &stdout, &stderr, approveEnv(nil, nil, atTime(approveT0+3))); code != exitOK || stdout.String() != want {
+		t.Fatalf("rc %d, stdout\n%s\nwant\n%s\nstderr %s", code, stdout.String(), want, stderr.String())
+	}
+
+	t.Run("none", func(t *testing.T) {
+		db, creds := approverStore(t)
+		var stdout, stderr bytes.Buffer
+		code := runApprove([]string{"--list", "--decided", "--db", db, "--credentials-db", creds}, &stdout, &stderr, approveEnv(nil, nil, atTime(approveT0)))
+		if code != exitOK || stdout.String() != "no decided mint requests\n" {
+			t.Fatalf("rc %d, stdout %q; want exitOK and %q\nstderr %s", code, stdout.String(), "no decided mint requests\n", stderr.String())
+		}
+	})
+}
--- a/cmd/ailang-worldd/approve_sig_test.go
+++ b/cmd/ailang-worldd/approve_sig_test.go
@@ -405,6 +405,8 @@
 		"signature-and-write-message": {firstRequestID, "--signature", "s", "--write-message", "m", "--db", db},
 		"list-and-signature":          {"--list", "--signature", "s", "--db", db},
 		"list-and-write-message":      {"--list", "--write-message", "m", "--db", db},
+		// F3 (sprint 4 D2):
+		"decided-without-list": {firstRequestID, "--decided", "--db", db},
 	} {
 		t.Run(name, func(t *testing.T) {
 			var stdout, stderr bytes.Buffer
--- a/cmd/ailang-worldd/help_test.go
+++ b/cmd/ailang-worldd/help_test.go
@@ -41,7 +41,7 @@
 		"log tail", "registry get", "commit --file <commit.json> [--session <file|token>]",
 		"tools list", "] call <tool>", "] why ", "provenance", "setup", "doctor",
 		"session new <episode>", "session list", "session revoke", "session mint", "help [<verb>]",
-		"approve <request-id>", "approve <request-id> --signature <file>", "approve <request-id> --write-message <file>", "approve --list", "approver enroll --name", "approver list", "approver revoke <name>",
+		"approve <request-id>", "approve <request-id> --signature <file>", "approve <request-id> --write-message <file>", "approve --list [--decided]", "approver enroll --name", "approver list", "approver revoke <name>",
 	} {
 		if !strings.Contains(out.String(), want) {
 			t.Errorf("top-level usage lacks %q", want)
--- a/cmd/ailang-worldd/main.go
+++ b/cmd/ailang-worldd/main.go
@@ -124,7 +124,7 @@
   ailang-worldd approve <request-id> [--db <path>] [--credentials-db <path>]
   ailang-worldd approve <request-id> --signature <file> [--db <path>] [--credentials-db <path>]
   ailang-worldd approve <request-id> --write-message <file> [--db <path>] [--credentials-db <path>]
-  ailang-worldd approve --list [--db <path>] [--credentials-db <path>]
+  ailang-worldd approve --list [--decided] [--db <path>] [--credentials-db <path>]
   ailang-worldd approver enroll --name <name> --format ssh --pubkey <file> [--db <path>] [--credentials-db <path>]
   ailang-worldd approver list [--db <path>] [--credentials-db <path>]
   ailang-worldd approver revoke <name> [--db <path>] [--credentials-db <path>]
````

### Patch D3 — `approver revoke` lists what it approved (K8, K9)

````diff D3
--- a/cmd/ailang-worldd/approve_decided_test.go
+++ b/cmd/ailang-worldd/approve_decided_test.go
@@ -54,3 +54,44 @@
 		}
 	})
 }
+
+// ---- D3 ----
+
+// CR12: approver revoke lists the sessions the approver approved, each with a session revoke line
+// that works, and revokes none of them (K8, Q4); an approver that approved nothing says so (K9).
+func TestApproverRevokeListsWhatItApproved(t *testing.T) {
+	a, b, _ := decidedRig(t)
+	seedApprover(t, a.credPath, "idle", otherSigner.PublicKey(), approveT0-100)
+	br := b.row(t)
+	revoke := func(name string) string {
+		t.Helper()
+		var stdout, stderr bytes.Buffer
+		args := []string{"revoke", name, "--db", a.r.f.db, "--credentials-db", a.credPath}
+		if code := runApprover(args, &stdout, &stderr, approveEnv(&lockedTerm{in: strings.NewReader("y\n")}, nil, atTime(approveT0+10))); code != exitOK {
+			t.Fatalf("revoke %s: rc %d\nstderr %s", name, code, stderr.String())
+		}
+		return stdout.String()
+	}
+	line := "ailang-worldd session revoke --db " + shellQuote(a.r.f.db) + " --credentials-db " + shellQuote(a.credPath) + " " + br.CredentialID
+	want := "revoked approver mark-mac\n" +
+		"it approved 1 session(s), newest first; each stays valid until it expires or is revoked:\n" +
+		"  " + b.id + "  episode ep2  at 2027-01-15T08:00:02Z  credential " + br.CredentialID[:16] + "\n" +
+		"    " + line + "\n"
+	if got := revoke("mark-mac"); got != want {
+		t.Fatalf("stdout\n%s\nwant\n%s", got, want)
+	}
+	if got, want := revoke("idle"), "revoked approver idle\nno session was approved by it\n"; got != want {
+		t.Fatalf("stdout %q, want %q", got, want)
+	}
+	if n := len(a.sessionRows(t)); n != 2 {
+		t.Fatalf("%d session rows after the revokes, want 2: approver revoke revokes no session", n)
+	}
+	var stdout, stderr bytes.Buffer
+	if code := runSessionRevoke(strings.Fields(line)[3:], &stdout, &stderr); code != exitOK {
+		t.Fatalf("the printed line %q: rc %d\nstderr %s", line, code, stderr.String())
+	}
+	rows := a.sessionRows(t)
+	if len(rows) != 1 || rows[0].CredentialID == br.CredentialID {
+		t.Fatalf("session rows after the printed line: %+v, want only a's", rows)
+	}
+}
--- a/cmd/ailang-worldd/approver.go
+++ b/cmd/ailang-worldd/approver.go
@@ -133,6 +133,8 @@
          fingerprint, enrolled, last used, and active or revoked.
 revoke   marks the approver revoked after y. The row is kept: it never
          approves again, and its name and key can never be enrolled again.
+         It then lists the sessions the approver approved (at most 64, newest
+         first), each with its 'session revoke' line. It revokes none of them.
 
   --name <name>             1-48 of a-z 0-9 . _ - (the first a letter or digit);
                             its approvals record approved_by=ssh:<name>
@@ -417,16 +419,17 @@
 		return code
 	}
 
-	// The write (K5b, K6, K7).
+	// The write (K5b, K6, K7). Its root and its writer end before K8 reads.
 	ctx, cancel := context.WithTimeout(context.Background(), approveReadBudget)
-	defer cancel()
 	st, err := openCredentialsForSession(ctx, db, credPath)
 	if err != nil { // K5b
+		cancel()
 		fmt.Fprintf(stderr, "%scannot open the credentials store for writing: %v\n", prefix, err)
 		return exitFatal
 	}
-	defer func() { _ = st.Close() }()
 	err = approverRevoke(st, ctx, name, env.now())
+	_ = st.Close()
+	cancel()
 	switch {
 	case errors.Is(err, store.ErrApproverRevoked), errors.Is(err, store.ErrApproverUnknown): // K6
 		return approverRefuse(stderr, prefix, "%v (another revoke ran first); nothing was revoked", err)
@@ -435,5 +438,30 @@
 		return exitFatal
 	}
 	fmt.Fprintf(stdout, "revoked approver %s\n", name) // K7
+
+	// K8, K9 (design §2.10 [r2]): what it approved, under its own read-only root. Q4: it revokes
+	// no session; the operator decides.
+	list, cancelList := context.WithTimeout(context.Background(), approveReadBudget)
+	defer cancelList()
+	ro, err = store.OpenCredentialsReadOnly(credPath)
+	if err != nil {
+		fmt.Fprintf(stderr, "%srevoked, but cannot read the credentials store %s to list what it approved: %v; fix: ailang-worldd approve --list --decided\n", prefix, credPath, err)
+		return exitFatal
+	}
+	rows, err := ro.DecidedMintRequests(list, row.Format+":"+name, 0)
+	_ = ro.Close()
+	if err != nil {
+		fmt.Fprintf(stderr, "%srevoked, but cannot list what it approved: %v; fix: ailang-worldd approve --list --decided\n", prefix, err)
+		return exitFatal
+	}
+	if len(rows) == 0 { // K9
+		fmt.Fprintln(stdout, "no session was approved by it")
+		return exitOK
+	}
+	fmt.Fprintf(stdout, "it approved %d session(s), newest first; each stays valid until it expires or is revoked:\n", len(rows)) // K8
+	for _, r := range rows {
+		fmt.Fprintf(stdout, "  %s  episode %s  at %s  credential %s\n    ailang-worldd session revoke --db %s --credentials-db %s %s\n",
+			r.RequestID, r.EpisodeID, approverTime(r.ConsumedAt), credentialPrefix(r.CredentialID), shellQuote(db), shellQuote(credPath), r.CredentialID)
+	}
 	return exitOK
 }
--- a/cmd/ailang-worldd/approver_revoke_test.go
+++ b/cmd/ailang-worldd/approver_revoke_test.go
@@ -39,8 +39,8 @@
 		if code != exitOK {
 			t.Fatalf("code %d, want exitOK; stderr %q", code, errw.String())
 		}
-		if out.String() != "revoked approver mark-mac\n" {
-			t.Fatalf("stdout %q, want %q", out.String(), "revoked approver mark-mac\n")
+		if want := "revoked approver mark-mac\nno session was approved by it\n"; out.String() != want { // K7, K9 (sprint 4 D3)
+			t.Fatalf("stdout %q, want %q", out.String(), want)
 		}
 		wantPrompt := "Revoke approver mark-mac (ssh-ed25519 " + edFingerprint + ")? It will never approve again, and its name and key can never be enrolled again. [y/N] "
 		if term.Output() != wantPrompt {
--- a/host/store/context_roots_test.go
+++ b/host/store/context_roots_test.go
@@ -20,7 +20,7 @@
 	"cmd/ailang-worldd/approve.go|runApprove|Background":              1, // row 201 M6: approve's read-only request lookup or --list (3 s)
 	"cmd/ailang-worldd/approver.go|runApproverList|Background":        1, // row 201 A1b M3: read-only ListApprovers (3 s)
 	"cmd/ailang-worldd/approver.go|runApproverEnroll|Background":      2, // row 201 A1b M3: the read-only pre-check (3 s), then the write after y (3 s)
-	"cmd/ailang-worldd/approver.go|runApproverRevoke|Background":      2, // row 201 A1b M3: the read-only lookup (3 s), then the write after y (3 s)
+	"cmd/ailang-worldd/approver.go|runApproverRevoke|Background":      3, // row 201 A1b M3: the read-only lookup (3 s), the write after y (3 s); sprint 4 K8: the read-only listing (3 s)
 	"cmd/ailang-worldd/episode_launch.go|launchEpisode|Background":    3, // row 217 M2: HTTP preflight (8 s), agent (wall cap + margin), append (3 s)
 
 	"cmd/ailang-worldd/cli.go|get|Background":                                1,
````

### Patch S1 — the `--sign-with` seam (gated)

````diff S1
--- a/cmd/ailang-worldd/approve.go
+++ b/cmd/ailang-worldd/approve.go
@@ -23,12 +23,14 @@
 
 const approveHelp = `usage: ailang-worldd approve <request-id> [--db <path>] [--credentials-db <path>]
        ailang-worldd approve <request-id> --signature <file> [--db <path>] [--credentials-db <path>]
+       ailang-worldd approve <request-id> --sign-with <key> [--db <path>] [--credentials-db <path>]
        ailang-worldd approve <request-id> --write-message <file> [--db <path>] [--credentials-db <path>]
        ailang-worldd approve --list [--decided] [--db <path>] [--credentials-db <path>]
 
 ATTENDED, or SIGNED. Mints the session a stored mint request describes (row
-201), after a y/N typed at /dev/tty, or, with --signature, after an enrolled
-approver's SSH signature over the request verifies, with no terminal.
+201), after a y/N typed at /dev/tty, or, with --signature or --sign-with,
+after an enrolled approver's SSH signature over the request verifies, with no
+terminal.
 'ailang-worldd episode start <episode> ... --request' makes the request without
 a terminal; nothing is minted until approve.
 
@@ -50,6 +52,12 @@
   checked before any network call, the approver is re-checked in the mint's own
   transaction, and the request records approved_by=ssh:<name>. A signature
   proves that something holding the key signed, not that you were there.
+--sign-with <key> signs that message by running ssh-keygen -Y sign -f <key>
+  (a private key file, or a .pub whose key ssh-agent holds; ssh-keygen may ask
+  for a passphrase itself) and checks the signature as --signature does. It
+  waits at most until the request expires. World never reads the private key.
+  A .pub path does not prove the agent signed: ssh-keygen uses the private
+  file beside it when there is one.
 --write-message <file> writes that message (90 bytes) to a NEW file and prints
   the ssh-keygen -Y sign line and the --signature step. It opens no terminal
   and writes no store row: the request stays pending.
@@ -89,8 +97,7 @@
 		return approval.MintRequest{}, fmt.Sprintf("the clock is behind request %s's creation time (%s); refusing",
 			row.RequestID, time.Unix(row.CreatedAt, 0).UTC().Format(time.RFC3339))
 	case approval.RefuseExpired:
-		return approval.MintRequest{}, fmt.Sprintf("request %s expired at %s; fix: make a new request with episode start --request",
-			row.RequestID, time.Unix(row.ExpiresAt, 0).UTC().Format(time.RFC3339))
+		return approval.MintRequest{}, expiredRefusal(row)
 	case approval.RefuseDigestMismatch:
 		return approval.MintRequest{}, digestRefusal
 	default:
@@ -98,6 +105,12 @@
 	}
 }
 
+// expiredRefusal is the refusal of a request that expired before it was approved.
+func expiredRefusal(row store.MintRequestRow) string {
+	return fmt.Sprintf("request %s expired at %s; fix: make a new request with episode start --request",
+		row.RequestID, time.Unix(row.ExpiresAt, 0).UTC().Format(time.RFC3339))
+}
+
 // runApprove is `ailang-worldd approve <request-id>` and `approve --list` (row 201 §2.6). Every
 // refusal before the mint names its fix and exits exitUsage; a failure at or after the write is
 // provisionEpisode's exitFatal.
@@ -117,6 +130,7 @@
 	decided := fs.Bool("decided", false, "with --list: print the used requests and who approved them")
 	sigPath := fs.String("signature", "", "approve with an enrolled approver's SSHSIG over the request (no terminal)")
 	msgPath := fs.String("write-message", "", "write the request's approval message to this new file, to sign elsewhere")
+	signWith := fs.String("sign-with", "", "approve by signing the request with ssh-keygen -Y sign -f <key> (no terminal)")
 	positional, err := parseInterleaved(fs, args)
 	if errors.Is(err, flag.ErrHelp) {
 		fmt.Fprint(stdout, approveHelp)
@@ -130,11 +144,17 @@
 		fmt.Fprint(stderr, approveHelp)
 		return exitUsage
 	}
-	if *sigPath != "" && *msgPath != "" { // F1: one of --signature and --write-message
+	modes := 0
+	for _, v := range []string{*sigPath, *signWith, *msgPath} {
+		if v != "" {
+			modes++
+		}
+	}
+	if modes > 1 { // F1: at most one of --signature, --sign-with and --write-message
 		fmt.Fprint(stderr, approveHelp)
 		return exitUsage
 	}
-	if *list && (*sigPath != "" || *msgPath != "") { // F2: --list is read-only
+	if *list && modes != 0 { // F2: --list is read-only
 		fmt.Fprint(stderr, approveHelp)
 		return exitUsage
 	}
@@ -152,6 +172,10 @@
 			fmt.Fprintf(stderr, "%sinternal error: no prompt timer (sessionEnv.after)\n", prefix)
 			return exitFatal
 		}
+		if *signWith != "" && env.sign == nil { // F7
+			fmt.Fprintf(stderr, "%sinternal error: no signer (sessionEnv.sign)\n", prefix)
+			return exitFatal
+		}
 	}
 	credPath, err := credentialsFor(*db, *credsDB)
 	if err != nil {
@@ -215,6 +239,15 @@
 			return refuse("%s", refusal)
 		}
 	}
+	if *signWith != "" { // the signer's output, judged as --signature's (design §2.7 step 3, F7a, F8)
+		armored, refusal := signForApproval(env, *signWith, row, stderr)
+		if refusal != "" {
+			return refuse("%s", refusal)
+		}
+		if signed, refusal = approveArmored(ctx, credPath, "ssh-keygen's output for --sign-with "+*signWith, armored, row, now); refusal != "" {
+			return refuse("%s", refusal)
+		}
+	}
 
 	// 2. The same store: the request was made for this world store and credentials file.
 	w, werr := store.CanonicalPath(*db)
--- a/cmd/ailang-worldd/approve_sig.go
+++ b/cmd/ailang-worldd/approve_sig.go
@@ -11,6 +11,7 @@
 	"fmt"
 	"io"
 	"os"
+	"time"
 
 	"github.com/sunholo-data/ailang-world/host/approval"
 	"github.com/sunholo-data/ailang-world/host/store"
@@ -31,9 +32,16 @@
 	if err != nil { // F10
 		return nil, fmt.Sprintf("--signature: %v", err)
 	}
+	return approveArmored(ctx, credPath, "--signature "+path, armored, row, now)
+}
+
+// approveArmored judges an armored signature from source (--signature's file or --sign-with's
+// signer): parsed (F11), its key looked up among the enrolled approvers (F12, F13), and judged by
+// the signature law against the stored row (F14).
+func approveArmored(ctx context.Context, credPath, source string, armored []byte, row store.MintRequestRow, now int64) (*approval.SignedApproval, string) {
 	sig, err := approval.ParseSSHSIG(armored)
 	if err != nil { // F11
-		return nil, fmt.Sprintf("--signature %s is not an SSH signature: %v", path, err)
+		return nil, fmt.Sprintf("%s is not an SSH signature: %v", source, err)
 	}
 	ro, err := store.OpenCredentialsReadOnly(credPath)
 	if err != nil {
@@ -70,9 +78,7 @@
 // and prints the out-of-band signing line and the next step. It opens no terminal and writes no
 // store row, so the request stays pending. It returns the refusal text, or "".
 func writeApprovalMessage(stdout io.Writer, path string, row store.MintRequestRow, mr approval.MintRequest) string {
-	stored, _ := hex.DecodeString(row.Digest)       // admission has shown it is 32 bytes
-	challenge, _ := hex.DecodeString(row.Challenge) // F4 has shown it is 32 bytes
-	msg, err := approval.SignedMessage([32]byte(stored), [32]byte(challenge), row.ExpiresAt)
+	msg, err := approvalMessage(row)
 	if err != nil {
 		return fmt.Sprintf("request %s: %v", row.RequestID, err)
 	}
@@ -95,3 +101,35 @@
 		row.RequestID, shellQuote(path), shellQuote(mr.WorldDB), shellQuote(mr.CredentialsDB))
 	return ""
 }
+
+// approvalMessage is the admitted row's 90-byte approval message (design §2.3).
+func approvalMessage(row store.MintRequestRow) ([]byte, error) {
+	stored, _ := hex.DecodeString(row.Digest)       // admission has shown it is 32 bytes
+	challenge, _ := hex.DecodeString(row.Challenge) // F4 has shown it is 32 bytes
+	return approval.SignedMessage([32]byte(stored), [32]byte(challenge), row.ExpiresAt)
+}
+
+// ---- S1 (sprint 4) ----
+
+// signForApproval is approve --sign-with's signing (design §2.7 step 3, §2.9): the admitted row's
+// message, signed by env.sign under a new root bounded by the request's remaining lifetime. A
+// request that expired since admission is refused before the signer runs (F7a), so the root is
+// never made with a non-positive duration; a signer error is refused (F8). It returns the armored
+// signature or the refusal text. The caller has checked env.sign (F7).
+func signForApproval(env sessionEnv, key string, row store.MintRequestRow, stderr io.Writer) ([]byte, string) {
+	rem := row.ExpiresAt - env.now()
+	if rem <= 0 { // F7a
+		return nil, expiredRefusal(row)
+	}
+	msg, err := approvalMessage(row)
+	if err != nil {
+		return nil, fmt.Sprintf("request %s: %v", row.RequestID, err)
+	}
+	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(rem)*time.Second)
+	defer cancel()
+	armored, err := env.sign(ctx, key, msg, stderr)
+	if err != nil { // F8
+		return nil, fmt.Sprintf("--sign-with %s: %v; nothing was approved", key, err)
+	}
+	return armored, ""
+}
--- a/cmd/ailang-worldd/approve_sig_test.go
+++ b/cmd/ailang-worldd/approve_sig_test.go
@@ -407,6 +407,10 @@
 		"list-and-write-message":      {"--list", "--write-message", "m", "--db", db},
 		// F3 (sprint 4 D2):
 		"decided-without-list": {firstRequestID, "--decided", "--db", db},
+		// F1, F2 with --sign-with (sprint 4 S2):
+		"signature-and-sign-with":     {firstRequestID, "--signature", "s", "--sign-with", "k", "--db", db},
+		"sign-with-and-write-message": {firstRequestID, "--sign-with", "k", "--write-message", "m", "--db", db},
+		"list-and-sign-with":          {"--list", "--sign-with", "k", "--db", db},
 	} {
 		t.Run(name, func(t *testing.T) {
 			var stdout, stderr bytes.Buffer
--- a/cmd/ailang-worldd/approver.go
+++ b/cmd/ailang-worldd/approver.go
@@ -122,7 +122,7 @@
 ATTENDED (enroll, revoke). An approver is a name and one SSH public key, kept in
 the credentials store (row 201). An enrolled, unrevoked key approves a mint
 request with its signature: 'ailang-worldd approve <request-id> --signature
-<file>'. enroll and revoke each ask y at /dev/tty and wait at most 2m0s for it;
+<file>' or '--sign-with <key>'. enroll and revoke each ask y at /dev/tty and wait at most 2m0s for it;
 an agent cannot run them.
 
 enroll   reads --pubkey once (one ssh-ed25519 or ecdsa-sha2-nistp256 key line,
--- a/cmd/ailang-worldd/help_test.go
+++ b/cmd/ailang-worldd/help_test.go
@@ -41,7 +41,7 @@
 		"log tail", "registry get", "commit --file <commit.json> [--session <file|token>]",
 		"tools list", "] call <tool>", "] why ", "provenance", "setup", "doctor",
 		"session new <episode>", "session list", "session revoke", "session mint", "help [<verb>]",
-		"approve <request-id>", "approve <request-id> --signature <file>", "approve <request-id> --write-message <file>", "approve --list [--decided]", "approver enroll --name", "approver list", "approver revoke <name>",
+		"approve <request-id>", "approve <request-id> --signature <file>", "approve <request-id> --sign-with <key>", "approve <request-id> --write-message <file>", "approve --list [--decided]", "approver enroll --name", "approver list", "approver revoke <name>",
 	} {
 		if !strings.Contains(out.String(), want) {
 			t.Errorf("top-level usage lacks %q", want)
--- a/cmd/ailang-worldd/main.go
+++ b/cmd/ailang-worldd/main.go
@@ -123,6 +123,7 @@
                     [--sessions-dir <dir>] [--exec-profiles-dir <dir>] [--addr 127.0.0.1:7644 (or service.json's address)] [--request]
   ailang-worldd approve <request-id> [--db <path>] [--credentials-db <path>]
   ailang-worldd approve <request-id> --signature <file> [--db <path>] [--credentials-db <path>]
+  ailang-worldd approve <request-id> --sign-with <key> [--db <path>] [--credentials-db <path>]
   ailang-worldd approve <request-id> --write-message <file> [--db <path>] [--credentials-db <path>]
   ailang-worldd approve --list [--decided] [--db <path>] [--credentials-db <path>]
   ailang-worldd approver enroll --name <name> --format ssh --pubkey <file> [--db <path>] [--credentials-db <path>]
--- a/cmd/ailang-worldd/session.go
+++ b/cmd/ailang-worldd/session.go
@@ -35,6 +35,9 @@
 	random func(b []byte) error
 	// after is time.After in production: the timer that bounds approve's y/N (row 201 §2.6a).
 	after func(time.Duration) <-chan time.Time
+	// sign is approve --sign-with's signer, signWithSSHKeygen in production (row 201 A1b §2.9):
+	// the armored SSHSIG of msg by key. World never reads the private key; ssh-keygen does.
+	sign func(ctx context.Context, key string, msg []byte, stderr io.Writer) ([]byte, error)
 }
 
 // sessionTTYFix ends every terminal refusal of the attended session verbs. It
--- a/host/store/context_roots_test.go
+++ b/host/store/context_roots_test.go
@@ -18,6 +18,7 @@
 	"cmd/ailang-worldd/episode_resume.go|runEpisodeResume|Background": 1, // row 217 M3: resume preflight (8 s)
 	"cmd/ailang-worldd/mint_request.go|createMintRequest|Background":  1, // row 201 M5: --request's store open and insert (3 s)
 	"cmd/ailang-worldd/approve.go|runApprove|Background":              1, // row 201 M6: approve's read-only request lookup or --list (3 s)
+	"cmd/ailang-worldd/approve_sig.go|signForApproval|Background":     1, // row 201 A1b sprint 4: --sign-with's signer, bounded by the request's remaining lifetime
 	"cmd/ailang-worldd/approver.go|runApproverList|Background":        1, // row 201 A1b M3: read-only ListApprovers (3 s)
 	"cmd/ailang-worldd/approver.go|runApproverEnroll|Background":      2, // row 201 A1b M3: the read-only pre-check (3 s), then the write after y (3 s)
 	"cmd/ailang-worldd/approver.go|runApproverRevoke|Background":      3, // row 201 A1b M3: the read-only lookup (3 s), the write after y (3 s); sprint 4 K8: the read-only listing (3 s)
````

### Patch S2 — the fake-signer tests CR8, CR9 (gated)

````diff S2
--- a/cmd/ailang-worldd/approve_signwith_test.go
+++ b/cmd/ailang-worldd/approve_signwith_test.go
@@ -0,0 +1,154 @@
+package main
+
+// Row 201 A1b sprint 4 (plan design_docs/planned/w-world-approvals-by-key-a1b-sprint4-plan.md §3
+// S2, S3): approve --sign-with. Every store is under t.TempDir(); the fake signer signs with
+// approvaltest's fixed seeds, and the real-tool tests make their keys with ssh-keygen in
+// t.TempDir() under a temp HOME with no agent (D-WORLD-109).
+
+import (
+	"bytes"
+	"context"
+	"errors"
+	"io"
+	"strings"
+	"testing"
+	"time"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/daemon"
+)
+
+// fakeSigner is a sign seam: it records what approve passed and returns ret's answer.
+type fakeSigner struct {
+	calls int
+	key   string
+	msg   []byte
+	left  time.Duration // the context's deadline, measured from the call
+	ret   func(msg []byte) ([]byte, error)
+}
+
+func (f *fakeSigner) sign(ctx context.Context, key string, msg []byte, _ io.Writer) ([]byte, error) {
+	f.calls++
+	f.key, f.msg = key, bytes.Clone(msg)
+	if d, ok := ctx.Deadline(); ok {
+		f.left = time.Until(d)
+	}
+	if err := ctx.Err(); err != nil {
+		return nil, err
+	}
+	return f.ret(msg)
+}
+
+// signEnv is approveEnv with f as the sign seam (nil f: no seam at all).
+func signEnv(f *fakeSigner, opened *bool, clock func() int64) sessionEnv {
+	env := approveEnv(nil, opened, clock)
+	if f != nil {
+		env.sign = f.sign
+	}
+	return env
+}
+
+// countDaemonCalls counts fetchDaemonHealth calls for the rest of the test.
+func countDaemonCalls(t *testing.T) *int {
+	t.Helper()
+	calls := 0
+	inner := fetchDaemonHealth
+	t.Cleanup(func() { fetchDaemonHealth = inner })
+	fetchDaemonHealth = func(addr string) (daemon.HealthResponse, error) { calls++; return inner(addr) }
+	return &calls
+}
+
+// CR8: --sign-with signs the request's exact message through the seam and judges the output as
+// --signature does: no seam is an internal error (F7), a signer error is refused (F8), another
+// approver's valid signature over another message and garbage are refused (F14, F11), and the
+// right signature mints as ssh:mark-mac with no terminal.
+func TestApproveSignWithFakeSigner(t *testing.T) {
+	for _, arm := range []struct {
+		name string
+		ret  func(msg []byte) ([]byte, error) // nil: no sign seam
+		code int
+		want string // exact stderr when it ends in "\n", else a fragment
+	}{
+		{name: "no-signer", code: exitFatal, want: "ailang-worldd approve: internal error: no signer (sessionEnv.sign)\n"},
+		{name: "signer-error", ret: func([]byte) ([]byte, error) { return nil, errors.New("agent refused") }, code: exitUsage,
+			want: "ailang-worldd approve: --sign-with /keys/k: agent refused; nothing was approved\n"},
+		{name: "other-approver-other-message", ret: func([]byte) ([]byte, error) {
+			return otherSigner.Sign(approval.Namespace, []byte("another message")), nil
+		}, code: exitUsage, want: "the signature does not approve request "},
+		{name: "garbage", ret: func([]byte) ([]byte, error) { return []byte("not a signature\n"), nil }, code: exitUsage,
+			want: "ssh-keygen's output for --sign-with /keys/k is not an SSH signature: approval: sshsig: bad armor"},
+		{name: "good", ret: func(msg []byte) ([]byte, error) { return sigSigner.Sign(approval.Namespace, msg), nil }, code: exitOK},
+	} {
+		t.Run(arm.name, func(t *testing.T) {
+			a, _ := signedRig(t)
+			seedApprover(t, a.credPath, "other", otherSigner.PublicKey(), approveT0-100)
+			calls := countDaemonCalls(t)
+			var f *fakeSigner
+			if arm.ret != nil {
+				f = &fakeSigner{ret: arm.ret}
+			}
+			var opened bool
+			var stdout, stderr bytes.Buffer
+			code := runApprove(a.args("--sign-with", "/keys/k"), &stdout, &stderr, signEnv(f, &opened, atTime(approveT0+1)))
+			if code != arm.code || opened {
+				t.Fatalf("rc %d, terminal opened %v; want %d and no terminal\nstderr %s", code, opened, arm.code, stderr.String())
+			}
+			row := a.row(t)
+			if arm.code != exitOK {
+				if strings.HasSuffix(arm.want, "\n") && stderr.String() != arm.want || !strings.Contains(stderr.String(), arm.want) {
+					t.Fatalf("stderr %q, want %q", stderr.String(), arm.want)
+				}
+				if row.Consumed || len(a.sessionRows(t)) != 0 || *calls != 0 {
+					t.Fatalf("consumed %v, daemon calls %d; want pending, no session, no network call", row.Consumed, *calls)
+				}
+				return
+			}
+			if !row.Consumed || row.ApprovedBy != "ssh:mark-mac" {
+				t.Fatalf("consumed %v by %q, want ssh:mark-mac", row.Consumed, row.ApprovedBy)
+			}
+			if f.calls != 1 || f.key != "/keys/k" || !bytes.Equal(f.msg, requestMessage(t, a)) {
+				t.Fatalf("the seam: %d call(s), key %q, msg %x; want 1, /keys/k and the request's message", f.calls, f.key, f.msg)
+			}
+		})
+	}
+}
+
+// CR9: the signer's context ends when the request expires: 37 s left gives a 37 s deadline, 1 s
+// gives 1 s, and a request that expired after admission is refused with the standard text before
+// the signer is called (F7a).
+func TestApproveSignWithDeadline(t *testing.T) {
+	stop := func([]byte) ([]byte, error) { return nil, errors.New("stop") }
+	for _, arm := range []struct {
+		name     string
+		clock    func() int64
+		calls    int
+		min, max time.Duration
+	}{
+		{"37-seconds-left", atTime(approveT0 + 600 - 37), 1, 36 * time.Second, 37 * time.Second},
+		{"1-second-left", clockSeq(1, approveT0+598, approveT0+599), 1, 1, time.Second},
+		{"expired-since-admission", clockSeq(1, approveT0+599, approveT0+600), 0, 0, 0},
+	} {
+		t.Run(arm.name, func(t *testing.T) {
+			a, _ := signedRig(t)
+			f := &fakeSigner{ret: stop}
+			var stdout, stderr bytes.Buffer
+			if code := runApprove(a.args("--sign-with", "/keys/k"), &stdout, &stderr, signEnv(f, nil, arm.clock)); code != exitUsage {
+				t.Fatalf("rc %d, want exitUsage\nstderr %s", code, stderr.String())
+			}
+			if f.calls != arm.calls {
+				t.Fatalf("the signer was called %d time(s), want %d\nstderr %s", f.calls, arm.calls, stderr.String())
+			}
+			if arm.calls == 0 {
+				want := "ailang-worldd approve: request " + a.id + " expired at 2027-01-15T08:10:00Z; fix: make a new request with episode start --request\n"
+				if stderr.String() != want {
+					t.Fatalf("stderr %q, want %q", stderr.String(), want)
+				}
+			} else if f.left < arm.min || f.left > arm.max {
+				t.Fatalf("the signer's deadline was %s ahead, want within [%s, %s]", f.left, arm.min, arm.max)
+			}
+			if a.row(t).Consumed {
+				t.Fatal("the request was consumed")
+			}
+		})
+	}
+}
````

### Patch S3 — the production signer with the real tool, the CI step, QUICKSTART §15 (gated)

````diff S3
--- a/.github/workflows/ci.yml
+++ b/.github/workflows/ci.yml
@@ -135,22 +135,25 @@
           go env GOVERSION
 
       # Row 201 (w-world-approvals-by-key-a0a1 §2.8, test XC1): the SSHSIG verifier is cross-checked
-      # against REAL ssh-keygen. WORLD_REQUIRE_SSH_KEYGEN=1 turns a missing ssh-keygen into a red,
+      # against REAL ssh-keygen; row 201 A1b sprint 4 adds the test signer (XS2) and approve
+      # --sign-with end to end with the production signer (XS1), each with temp keys and no agent.
+      # WORLD_REQUIRE_SSH_KEYGEN=1 turns a missing ssh-keygen into a red,
       # and the PASS grep turns a SKIP into a red. Whether ubuntu-24.04 ships ssh-keygen is
       # UNVERIFIED; if it does not, this step fails naming the apt-get line.
       - name: Row 201 SSHSIG cross-check with real ssh-keygen, verbose
-        timeout-minutes: 3
+        timeout-minutes: 8
         run: |
           set -euo pipefail
           command -v ssh-keygen || { echo "::error::ssh-keygen missing (apt-get install openssh-client)"; exit 1; }
           ssh -V
           log="$RUNNER_TEMP/sshsig-crosscheck.log"
           rc=0
-          WORLD_REQUIRE_SSH_KEYGEN=1 go test ./host/approval/ -count=1 -v \
-            -run '^TestSSHSIGCrossCheckWithSSHKeygen$' >"$log" 2>&1 || rc=$?
+          WORLD_REQUIRE_SSH_KEYGEN=1 go test ./host/approval/... ./cmd/ailang-worldd/ -count=1 -v \
+            -run '^(TestSSHSIGCrossCheckWithSSHKeygen|TestTestSignerAgreesWithSSHKeygen|TestSignWithSSHKeygenEndToEnd)$' >"$log" 2>&1 || rc=$?
           cat "$log"
           [ "$rc" -eq 0 ] || exit "$rc"
-          for t in TestSSHSIGCrossCheckWithSSHKeygen TestSSHSIGCrossCheckWithSSHKeygen/ed25519 TestSSHSIGCrossCheckWithSSHKeygen/ecdsa-p256; do
+          for t in TestSSHSIGCrossCheckWithSSHKeygen TestSSHSIGCrossCheckWithSSHKeygen/ed25519 TestSSHSIGCrossCheckWithSSHKeygen/ecdsa-p256 \
+            TestTestSignerAgreesWithSSHKeygen TestSignWithSSHKeygenEndToEnd TestSignWithSSHKeygenEndToEnd/ed25519 TestSignWithSSHKeygenEndToEnd/ecdsa-p256; do
             grep -qE -- "^ *--- PASS: $t " "$log" || { echo "::error::$t did not PASS (SKIP?)"; exit 1; }
           done
 
--- a/cmd/ailang-worldd/approve_sig.go
+++ b/cmd/ailang-worldd/approve_sig.go
@@ -5,12 +5,15 @@
 // approver's SSHSIG over the request's exact message, with no terminal.
 
 import (
+	"bytes"
 	"context"
 	"crypto/sha256"
 	"encoding/hex"
+	"errors"
 	"fmt"
 	"io"
 	"os"
+	"os/exec"
 	"time"
 
 	"github.com/sunholo-data/ailang-world/host/approval"
@@ -133,3 +136,47 @@
 	}
 	return armored, ""
 }
+
+// ---- S3 (sprint 4) ----
+
+var (
+	errSignatureTooLarge = fmt.Errorf("ssh-keygen -Y sign printed over %d bytes", approveMaxSignatureBytes)
+	errEmptySignature    = errors.New("ssh-keygen -Y sign printed nothing")
+)
+
+// signWithSSHKeygen is the production signer (design §2.9): `ssh-keygen -Y sign -n
+// world-approval@sunholo -f <key>` with the message on stdin, so the signature comes back on
+// stdout and nothing is written beside the key. Its stderr (a passphrase prompt goes to /dev/tty)
+// is the operator's. PATH is not a trust boundary: whatever runs, approveArmored verifies its
+// output against the enrolled key over World's own message.
+func signWithSSHKeygen(ctx context.Context, key string, msg []byte, stderr io.Writer) ([]byte, error) {
+	cmd := exec.CommandContext(ctx, "ssh-keygen", "-Y", "sign", "-n", approval.Namespace, "-f", key)
+	cmd.Stdin = bytes.NewReader(msg)
+	out := &cappedBuffer{max: approveMaxSignatureBytes}
+	cmd.Stdout, cmd.Stderr = out, stderr
+	cmd.WaitDelay = time.Second       // the deadline holds even if a child of the signer keeps a pipe open
+	if err := cmd.Run(); err != nil { // G1
+		return nil, fmt.Errorf("ssh-keygen -Y sign: %w", err)
+	}
+	if out.over { // G2
+		return nil, errSignatureTooLarge
+	}
+	if out.buf.Len() == 0 { // G3
+		return nil, errEmptySignature
+	}
+	return out.buf.Bytes(), nil
+}
+
+// cappedBuffer keeps the first max bytes written to it and records whether more came.
+type cappedBuffer struct {
+	buf  bytes.Buffer
+	max  int
+	over bool
+}
+
+func (c *cappedBuffer) Write(p []byte) (int, error) {
+	n := min(len(p), c.max-c.buf.Len())
+	c.buf.Write(p[:n])
+	c.over = c.over || n < len(p)
+	return len(p), nil
+}
--- a/cmd/ailang-worldd/approve_signwith_test.go
+++ b/cmd/ailang-worldd/approve_signwith_test.go
@@ -10,6 +10,10 @@
 	"context"
 	"errors"
 	"io"
+	"os"
+	"os/exec"
+	"path/filepath"
+	"reflect"
 	"strings"
 	"testing"
 	"time"
@@ -152,3 +156,150 @@
 		})
 	}
 }
+
+// ---- S3 ----
+
+// fakeSSHKeygen puts an executable `ssh-keygen` running body (sh) first on PATH for the rest of the
+// test, and returns its directory. The real tool is never run by it.
+func fakeSSHKeygen(t *testing.T, body string) string {
+	t.Helper()
+	dir := t.TempDir()
+	if err := os.WriteFile(filepath.Join(dir, "ssh-keygen"), []byte("#!/bin/sh\n"+body), 0o755); err != nil {
+		t.Fatal(err)
+	}
+	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
+	return dir
+}
+
+// G1-G3 and the exec line, with a fake ssh-keygen first on PATH: the production signer runs
+// `ssh-keygen -Y sign -n world-approval@sunholo -f <key>` with the message on stdin and its stderr
+// passed through, and refuses a failed run (G1), over 16 KiB of output (G2), no output (G3) and a
+// signer still running when the request expires. Each is exitUsage, pending.
+func TestSignWithSSHKeygenRefusesBadOutput(t *testing.T) {
+	for _, arm := range []struct {
+		name, body string
+		now        int64
+		want       []string
+	}{
+		{"exec-line-and-no-output", "printf '%s\\n' \"$@\" > \"$(dirname \"$0\")/args\"\ncat > \"$(dirname \"$0\")/stdin\"\n", approveT0 + 1,
+			[]string{"--sign-with /keys/k: ssh-keygen -Y sign printed nothing; nothing was approved"}},
+		{"failed", "cat > /dev/null\necho 'No private key found' >&2\nexit 255\n", approveT0 + 1,
+			[]string{"No private key found\n", "--sign-with /keys/k: ssh-keygen -Y sign: exit status 255; nothing was approved"}},
+		{"over-16k", "cat > /dev/null\nhead -c 16385 /dev/zero | tr '\\000' a\n", approveT0 + 1,
+			[]string{"--sign-with /keys/k: ssh-keygen -Y sign printed over 16384 bytes; nothing was approved"}},
+		{"outlives-the-request", "cat > /dev/null\nsleep 20\n", approveT0 + 599, // 1 s left
+			[]string{"--sign-with /keys/k: ssh-keygen -Y sign: signal: killed; nothing was approved"}},
+	} {
+		t.Run(arm.name, func(t *testing.T) {
+			a, _ := signedRig(t)
+			dir := fakeSSHKeygen(t, arm.body)
+			env := approveEnv(nil, nil, atTime(arm.now))
+			env.sign = signWithSSHKeygen
+			var stdout, stderr bytes.Buffer
+			start := time.Now()
+			if code := runApprove(a.args("--sign-with", "/keys/k"), &stdout, &stderr, env); code != exitUsage {
+				t.Fatalf("rc %d, want exitUsage\nstderr %s", code, stderr.String())
+			}
+			if took := time.Since(start); took > 10*time.Second {
+				t.Fatalf("approve took %s: the signer outlived the request's deadline", took)
+			}
+			for _, want := range arm.want {
+				if !strings.Contains(stderr.String(), want) {
+					t.Fatalf("stderr %q lacks %q", stderr.String(), want)
+				}
+			}
+			if a.row(t).Consumed {
+				t.Fatal("the request was consumed")
+			}
+			if arm.name != "exec-line-and-no-output" {
+				return
+			}
+			args, err := os.ReadFile(filepath.Join(dir, "args"))
+			if want := "-Y\nsign\n-n\nworld-approval@sunholo\n-f\n/keys/k\n"; err != nil || string(args) != want {
+				t.Fatalf("ssh-keygen's arguments %q (%v), want %q", args, err, want)
+			}
+			if in, err := os.ReadFile(filepath.Join(dir, "stdin")); err != nil || !bytes.Equal(in, requestMessage(t, a)) {
+				t.Fatalf("ssh-keygen's stdin %x (%v), want the request's message", in, err)
+			}
+		})
+	}
+}
+
+// The production env signs with signWithSSHKeygen: without it, --sign-with is the F7 internal error.
+func TestRealSessionEnvSignsWithSSHKeygen(t *testing.T) {
+	if got, want := reflect.ValueOf(realSessionEnv().sign).Pointer(), reflect.ValueOf(signWithSSHKeygen).Pointer(); got != want {
+		t.Fatalf("realSessionEnv().sign is %#x, want signWithSSHKeygen (%#x)", got, want)
+	}
+}
+
+// sshKeygen runs the real ssh-keygen with no agent (SSH_AUTH_SOCK emptied) and stdin from
+// /dev/null, under the test's temp HOME; a run over 30 s or a failure to start is t.Fatal.
+func sshKeygen(t *testing.T, args ...string) {
+	t.Helper()
+	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
+	defer cancel()
+	cmd := exec.CommandContext(ctx, "ssh-keygen", args...)
+	cmd.Env = append(os.Environ(), "SSH_AUTH_SOCK=")
+	if out, err := cmd.CombinedOutput(); err != nil {
+		t.Fatalf("ssh-keygen %v: %v\n%s", args, err, out)
+	}
+}
+
+// XS1: approve --sign-with with the production signer and the REAL ssh-keygen, for each key type:
+// a key made in t.TempDir() under a temp HOME, with no agent, signs and mints as ssh:mark-mac; its
+// .pub with the private file moved away and no agent is refused, naming ssh-keygen's failure (G1).
+// With WORLD_REQUIRE_SSH_KEYGEN=1 a missing ssh-keygen is a failure, not a skip (CI).
+func TestSignWithSSHKeygenEndToEnd(t *testing.T) {
+	if _, err := exec.LookPath("ssh-keygen"); err != nil {
+		if os.Getenv("WORLD_REQUIRE_SSH_KEYGEN") == "1" {
+			t.Fatal("WORLD_REQUIRE_SSH_KEYGEN=1 but ssh-keygen was not found")
+		}
+		t.Skip("SKIP COUNTED: ssh-keygen not found; the --sign-with end-to-end test did NOT run")
+	}
+	for _, kt := range []struct {
+		name string
+		args []string
+	}{{"ed25519", []string{"-t", "ed25519"}}, {"ecdsa-p256", []string{"-t", "ecdsa", "-b", "256"}}} {
+		t.Run(kt.name, func(t *testing.T) {
+			a := newApproveRig(t, "ep1", approveT0) // its episode run sets HOME to a temp dir
+			t.Setenv("SSH_AUTH_SOCK", "")
+			key := filepath.Join(t.TempDir(), "k")
+			sshKeygen(t, append(append([]string{"-q"}, kt.args...), "-N", "", "-C", "xs1", "-f", key)...)
+			pub, err := os.ReadFile(key + ".pub")
+			if err != nil {
+				t.Fatal(err)
+			}
+			k, err := approval.ParseAuthorizedKey(pub)
+			if err != nil {
+				t.Fatal(err)
+			}
+			seedApprover(t, a.credPath, "mark-mac", k, approveT0-100)
+			env := approveEnv(nil, nil, atTime(approveT0+1))
+			env.sign = signWithSSHKeygen
+			var stdout, stderr bytes.Buffer
+			if code := runApprove(a.args("--sign-with", key), &stdout, &stderr, env); code != exitOK {
+				t.Fatalf("approve --sign-with %s: rc %d, want exitOK\nstderr %s", key, code, stderr.String())
+			}
+			if row := a.row(t); row.ApprovedBy != "ssh:mark-mac" {
+				t.Fatalf("approved_by %q, want ssh:mark-mac", row.ApprovedBy)
+			}
+
+			b := newApproveRigIn(t, a.r, "ep2", approveT0)
+			if err := os.Rename(key, key+".moved"); err != nil {
+				t.Fatal(err)
+			}
+			stderr.Reset()
+			if code := runApprove(b.args("--sign-with", key+".pub"), &stdout, &stderr, env); code != exitUsage {
+				t.Fatalf("approve --sign-with k.pub, no private key and no agent: rc %d, want exitUsage\nstderr %s", code, stderr.String())
+			}
+			for _, want := range []string{"No private key found", "ssh-keygen -Y sign: exit status 255; nothing was approved"} {
+				if !strings.Contains(stderr.String(), want) {
+					t.Errorf("stderr lacks %q:\n%s", want, stderr.String())
+				}
+			}
+			if b.row(t).Consumed {
+				t.Fatal("the request was consumed with no key to sign it")
+			}
+		})
+	}
+}
--- a/cmd/ailang-worldd/session.go
+++ b/cmd/ailang-worldd/session.go
@@ -54,6 +54,7 @@
 		now:    func() int64 { return time.Now().Unix() },
 		random: func(b []byte) error { _, err := rand.Read(b); return err },
 		after:  time.After,
+		sign:   signWithSSHKeygen,
 	}
 }
 
--- a/docs/QUICKSTART.md
+++ b/docs/QUICKSTART.md
@@ -1793,13 +1793,16 @@
 ailang-worldd approver list --db /tmp/se-world/world.db
 ```
 
-To approve a §14 request with that key instead of a `y`: write the request's 90-byte message, sign it out of band with `ssh-keygen` in World's namespace, and pass the signature. `approve --signature` opens no terminal; it checks the signature against the enrolled, unrevoked key before any network call, re-checks the approver inside the mint's own transaction, and records `approved_by=ssh:demo`. `--write-message` prints these same two commands with the request's own paths.
+To approve a §14 request with that key instead of a `y`, either let `approve --sign-with` run `ssh-keygen -Y sign` on the key for you, or write the request's 90-byte message, sign it out of band with `ssh-keygen` in World's namespace, and pass the signature. `approve --signature` opens no terminal; it checks the signature against the enrolled, unrevoked key before any network call, re-checks the approver inside the mint's own transaction, and records `approved_by=ssh:demo`. `--write-message` prints these same two commands with the request's own paths. `approve --sign-with` checks `ssh-keygen`'s output the same way; World never reads the private key, and the signer waits at most until the request expires. `approve --list --decided` shows who approved each used request, and `approver revoke` lists the sessions the approver approved, each with its `session revoke` line; it revokes none of them.
 
 ```bash
 ailang-worldd episode start sig-demo --repo "$PWD" --db /tmp/se-world/world.db --request
 ailang-worldd approve mr-<the 16 hex it printed> --write-message /tmp/se-world/req.msg --db /tmp/se-world/world.db
 ssh-keygen -Y sign -n world-approval@sunholo -f /tmp/se-world/demo-key < /tmp/se-world/req.msg > /tmp/se-world/req.msg.sig
 ailang-worldd approve mr-<the same id> --signature /tmp/se-world/req.msg.sig --db /tmp/se-world/world.db
+ailang-worldd episode start sig-demo-2 --repo "$PWD" --db /tmp/se-world/world.db --request
+ailang-worldd approve mr-<the 16 hex it printed> --sign-with /tmp/se-world/demo-key --db /tmp/se-world/world.db
+ailang-worldd approve --list --decided --db /tmp/se-world/world.db
 ailang-worldd approver revoke demo --db /tmp/se-world/world.db
 ```
 
--- a/host/approval/approvaltest/sign_crosscheck_test.go
+++ b/host/approval/approvaltest/sign_crosscheck_test.go
@@ -0,0 +1,82 @@
+package approvaltest_test
+
+// Row 201 A1b sprint 4 (plan design_docs/planned/w-world-approvals-by-key-a1b-sprint4-plan.md §3
+// S3, XS2): the test signer is SSHSIG as the real ssh-keygen reads it, so every test that signs
+// with it tests the format the operator's tool produces. Test keys only; no agent, temp HOME.
+
+import (
+	"bytes"
+	"context"
+	"errors"
+	"os"
+	"os/exec"
+	"path/filepath"
+	"strings"
+	"testing"
+	"time"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
+)
+
+// keygen runs the real ssh-keygen with stdin, no agent and the test's HOME; it returns the
+// combined output and the exit code. A run over 30 s or a failure to start is t.Fatal.
+func keygen(t *testing.T, stdin []byte, args ...string) (string, int) {
+	t.Helper()
+	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
+	defer cancel()
+	cmd := exec.CommandContext(ctx, "ssh-keygen", args...)
+	cmd.Env = append(os.Environ(), "SSH_AUTH_SOCK=")
+	cmd.Stdin = bytes.NewReader(stdin)
+	out, err := cmd.CombinedOutput()
+	var ee *exec.ExitError
+	switch {
+	case err == nil:
+		return string(out), 0
+	case errors.As(err, &ee):
+		return string(out), ee.ExitCode()
+	default:
+		t.Fatalf("ssh-keygen %v: %v", args, err)
+		return "", 0
+	}
+}
+
+// XS2: ssh-keygen -Y check-novalidate accepts the test signer's signature over an approval
+// message in World's namespace and refuses it over a message one byte different; ssh-keygen -l
+// prints the signer's own fingerprint for its AuthorizedKey. With WORLD_REQUIRE_SSH_KEYGEN=1 a
+// missing ssh-keygen is a failure, not a skip (CI).
+func TestTestSignerAgreesWithSSHKeygen(t *testing.T) {
+	if _, err := exec.LookPath("ssh-keygen"); err != nil {
+		if os.Getenv("WORLD_REQUIRE_SSH_KEYGEN") == "1" {
+			t.Fatal("WORLD_REQUIRE_SSH_KEYGEN=1 but ssh-keygen was not found")
+		}
+		t.Skip("SKIP COUNTED: ssh-keygen not found; the test-signer cross-check did NOT run")
+	}
+	t.Setenv("HOME", t.TempDir())
+	dir := t.TempDir()
+	s := approvaltest.NewSigner(7)
+	msg, err := approval.SignedMessage([32]byte(bytes.Repeat([]byte{0x11}, 32)), [32]byte(bytes.Repeat([]byte{0x22}, 32)), 1800000600)
+	if err != nil {
+		t.Fatal(err)
+	}
+	sig := filepath.Join(dir, "msg.sig")
+	pub := filepath.Join(dir, "k.pub")
+	if err := os.WriteFile(sig, s.Sign(approval.Namespace, msg), 0o600); err != nil {
+		t.Fatal(err)
+	}
+	if err := os.WriteFile(pub, s.AuthorizedKey(), 0o600); err != nil {
+		t.Fatal(err)
+	}
+	out, rc := keygen(t, msg, "-Y", "check-novalidate", "-n", approval.Namespace, "-s", sig)
+	if want := `Good "` + approval.Namespace + `" signature`; rc != 0 || !strings.Contains(out, want) {
+		t.Fatalf("check-novalidate of the message: rc %d, output %q; want 0 and %q", rc, out, want)
+	}
+	flipped := bytes.Clone(msg)
+	flipped[len(flipped)-1] ^= 1
+	if out, rc := keygen(t, flipped, "-Y", "check-novalidate", "-n", approval.Namespace, "-s", sig); rc != 255 {
+		t.Fatalf("check-novalidate of a flipped message: rc %d (%q), want 255", rc, out)
+	}
+	if out, rc := keygen(t, nil, "-l", "-f", pub); rc != 0 || !strings.Contains(out, s.PublicKey().Fingerprint()) {
+		t.Fatalf("ssh-keygen -l: rc %d, output %q; want 0 and %s", rc, out, s.PublicKey().Fingerprint())
+	}
+}
````

### Patch S4 — revision: the signer moves to `host/sshsign`, scrubbed and driven (AC10)

````diff S4
--- a/cmd/ailang-worldd/approve_sig.go
+++ b/cmd/ailang-worldd/approve_sig.go
@@ -5,18 +5,16 @@
 // approver's SSHSIG over the request's exact message, with no terminal.
 
 import (
-	"bytes"
 	"context"
 	"crypto/sha256"
 	"encoding/hex"
-	"errors"
 	"fmt"
 	"io"
 	"os"
-	"os/exec"
 	"time"
 
 	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/sshsign"
 	"github.com/sunholo-data/ailang-world/host/store"
 )
 
@@ -139,44 +137,10 @@
 
 // ---- S3 (sprint 4) ----
 
-var (
-	errSignatureTooLarge = fmt.Errorf("ssh-keygen -Y sign printed over %d bytes", approveMaxSignatureBytes)
-	errEmptySignature    = errors.New("ssh-keygen -Y sign printed nothing")
-)
-
-// signWithSSHKeygen is the production signer (design §2.9): `ssh-keygen -Y sign -n
-// world-approval@sunholo -f <key>` with the message on stdin, so the signature comes back on
-// stdout and nothing is written beside the key. Its stderr (a passphrase prompt goes to /dev/tty)
-// is the operator's. PATH is not a trust boundary: whatever runs, approveArmored verifies its
-// output against the enrolled key over World's own message.
+// signWithSSHKeygen is the production signer (design §2.9): ssh-keygen from PATH through
+// host/sshsign's driven, environment-scrubbed site (AC10). PATH is not a trust boundary:
+// whatever runs, approveArmored verifies its output against the enrolled key over World's own
+// message.
 func signWithSSHKeygen(ctx context.Context, key string, msg []byte, stderr io.Writer) ([]byte, error) {
-	cmd := exec.CommandContext(ctx, "ssh-keygen", "-Y", "sign", "-n", approval.Namespace, "-f", key)
-	cmd.Stdin = bytes.NewReader(msg)
-	out := &cappedBuffer{max: approveMaxSignatureBytes}
-	cmd.Stdout, cmd.Stderr = out, stderr
-	cmd.WaitDelay = time.Second       // the deadline holds even if a child of the signer keeps a pipe open
-	if err := cmd.Run(); err != nil { // G1
-		return nil, fmt.Errorf("ssh-keygen -Y sign: %w", err)
-	}
-	if out.over { // G2
-		return nil, errSignatureTooLarge
-	}
-	if out.buf.Len() == 0 { // G3
-		return nil, errEmptySignature
-	}
-	return out.buf.Bytes(), nil
+	return sshsign.Sign(ctx, "ssh-keygen", key, msg, stderr)
 }
-
-// cappedBuffer keeps the first max bytes written to it and records whether more came.
-type cappedBuffer struct {
-	buf  bytes.Buffer
-	max  int
-	over bool
-}
-
-func (c *cappedBuffer) Write(p []byte) (int, error) {
-	n := min(len(p), c.max-c.buf.Len())
-	c.buf.Write(p[:n])
-	c.over = c.over || n < len(p)
-	return len(p), nil
-}
--- a/cmd/ailang-worldd/approve_signwith_test.go
+++ b/cmd/ailang-worldd/approve_signwith_test.go
@@ -225,19 +225,6 @@
 	}
 }
 
-// cappedBuffer holds at most its bound, however much is written, and records the overflow (G2).
-func TestCappedBufferKeepsItsBound(t *testing.T) {
-	c := &cappedBuffer{max: 4}
-	for _, w := range []string{"ab", "cdef", "g"} {
-		if n, err := c.Write([]byte(w)); n != len(w) || err != nil {
-			t.Fatalf("Write(%q) = (%d, %v), want (%d, nil)", w, n, err, len(w))
-		}
-	}
-	if c.buf.String() != "abcd" || !c.over {
-		t.Fatalf("buffer %q, over %v; want \"abcd\", true", c.buf.String(), c.over)
-	}
-}
-
 // The production env signs with signWithSSHKeygen: without it, --sign-with is the F7 internal error.
 func TestRealSessionEnvSignsWithSSHKeygen(t *testing.T) {
 	if got, want := reflect.ValueOf(realSessionEnv().sign).Pointer(), reflect.ValueOf(signWithSSHKeygen).Pointer(); got != want {
--- a/host/broker/registry_publish_test.go
+++ b/host/broker/registry_publish_test.go
@@ -32,6 +32,7 @@
 	"github.com/sunholo-data/ailang-world/host/hashref"
 	"github.com/sunholo-data/ailang-world/host/pkgproj"
 	"github.com/sunholo-data/ailang-world/host/replay"
+	"github.com/sunholo-data/ailang-world/host/sshsign"
 	"github.com/sunholo-data/ailang-world/host/store"
 	"github.com/sunholo-data/ailang-world/host/toolprobe"
 	"github.com/sunholo-data/ailang-world/host/worktree"
@@ -1223,6 +1224,9 @@
 		"host/toolprobe/toolprobe.go": {driveToolprobeOutput},
 		"host/checkrun/checkrun.go":   {driveCheckrun},
 		"host/agentrun/agentrun.go":   {driveAgentrun},
+		// Row 201 A1b sprint 4: `approve --sign-with` runs ssh-keygen -Y sign through
+		// host/sshsign's one site.
+		"host/sshsign/sshsign.go": {driveSSHSign},
 	}
 	if len(drivers) != len(files) {
 		t.Fatalf("AC10(a): %d files carry subprocess sites %v but %d have drivers; "+
@@ -1375,6 +1379,16 @@
 	if _, err := toolprobe.Output(boundedTestContext(t), probe, t.TempDir(),
 		"env", "GOROOT", "GOMODCACHE", "GOVERSION"); err != nil {
 		t.Fatalf("toolprobe.Output via the probe: %v", err)
+	}
+}
+
+// driveSSHSign launches approve --sign-with's signer with the probe standing in for
+// ssh-keygen. The probe prints a version banner, which Sign returns as the "signature" (the CLI
+// would then refuse it as not an SSHSIG); the launch and its environment are the measurement.
+func driveSSHSign(t *testing.T, probe string) {
+	t.Helper()
+	if _, err := sshsign.Sign(boundedTestContext(t), probe, filepath.Join(t.TempDir(), "k"), []byte("message"), io.Discard); err != nil {
+		t.Fatalf("sshsign.Sign via the probe: %v", err)
 	}
 }
 
--- a/host/sshsign/sshsign.go
+++ b/host/sshsign/sshsign.go
@@ -0,0 +1,77 @@
+// Package sshsign is the one bounded, environment-scrubbed subprocess site of `ailang-worldd
+// approve --sign-with` (row 201 A1b, design_docs/planned/w-world-approvals-by-key-a1b.md §2.9):
+// `ssh-keygen -Y sign -n world-approval@sunholo -f <key>` with the message on stdin, so the
+// signature comes back on stdout and nothing is written beside the key.
+//
+// Why a package: the AC10 census (host/broker registry_publish_test.go) requires every
+// production subprocess site to be DRIVEN and to scrub the registry credential, and a site in
+// package main cannot be driven from there, so the site lives here (as host/toolprobe's does)
+// with an injectable binary path the census points at its env-dumping probe.
+//
+// The child environment is the process environment minus every registry variable
+// (childenv.Scrubbed): SSH_AUTH_SOCK, HOME and PATH, which ssh-keygen needs, pass through.
+// PATH is not a trust boundary: whatever runs, the caller verifies its output against the
+// enrolled key over World's own message.
+package sshsign
+
+import (
+	"bytes"
+	"context"
+	"errors"
+	"fmt"
+	"io"
+	"os"
+	"os/exec"
+	"time"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/childenv"
+)
+
+// MaxSignatureBytes bounds the signer's output (G2): an Ed25519 or P-256 SSHSIG is under 1 KiB.
+const MaxSignatureBytes = 16 << 10
+
+// waitDelay bounds Wait after the deadline's kill, even if a child of the signer keeps a pipe open.
+const waitDelay = time.Second
+
+var (
+	ErrSignatureTooLarge = fmt.Errorf("ssh-keygen -Y sign printed over %d bytes", MaxSignatureBytes)
+	ErrEmptySignature    = errors.New("ssh-keygen -Y sign printed nothing")
+)
+
+// Sign runs `<bin> -Y sign -n world-approval@sunholo -f <key>` under ctx with msg on stdin and
+// returns its standard output, the armored signature. stderr (where ssh-keygen reports, and
+// from where a passphrase prompt goes to /dev/tty) is the operator's. A failed run (G1), over
+// MaxSignatureBytes of output (G2) and no output (G3) are errors.
+func Sign(ctx context.Context, bin, key string, msg []byte, stderr io.Writer) ([]byte, error) {
+	cmd := exec.CommandContext(ctx, bin, "-Y", "sign", "-n", approval.Namespace, "-f", key)
+	cmd.Env = childenv.Scrubbed(os.Environ())
+	cmd.Stdin = bytes.NewReader(msg)
+	out := &cappedBuffer{max: MaxSignatureBytes}
+	cmd.Stdout, cmd.Stderr = out, stderr
+	cmd.WaitDelay = waitDelay
+	if err := cmd.Run(); err != nil { // G1
+		return nil, fmt.Errorf("ssh-keygen -Y sign: %w", err)
+	}
+	if out.over { // G2
+		return nil, ErrSignatureTooLarge
+	}
+	if out.buf.Len() == 0 { // G3
+		return nil, ErrEmptySignature
+	}
+	return out.buf.Bytes(), nil
+}
+
+// cappedBuffer keeps the first max bytes written to it and records whether more came.
+type cappedBuffer struct {
+	buf  bytes.Buffer
+	max  int
+	over bool
+}
+
+func (c *cappedBuffer) Write(p []byte) (int, error) {
+	n := min(len(p), c.max-c.buf.Len())
+	c.buf.Write(p[:n])
+	c.over = c.over || n < len(p)
+	return len(p), nil
+}
--- a/host/sshsign/sshsign_test.go
+++ b/host/sshsign/sshsign_test.go
@@ -0,0 +1,63 @@
+package sshsign
+
+// Row 201 A1b sprint 4 S4 (plan design_docs/planned/w-world-approvals-by-key-a1b-sprint4-plan.md
+// §3 S4): the signer's buffer bound, and the child environment it hands ssh-keygen. No real
+// ssh-keygen and no key: the program is a temp shell script.
+
+import (
+	"io"
+	"os"
+	"path/filepath"
+	"strings"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/childenv"
+)
+
+// cappedBuffer holds at most its bound, however much is written, and records the overflow (G2).
+func TestCappedBufferKeepsItsBound(t *testing.T) {
+	c := &cappedBuffer{max: 4}
+	for _, w := range []string{"ab", "cdef", "g"} {
+		if n, err := c.Write([]byte(w)); n != len(w) || err != nil {
+			t.Fatalf("Write(%q) = (%d, %v), want (%d, nil)", w, n, err, len(w))
+		}
+	}
+	if c.buf.String() != "abcd" || !c.over {
+		t.Fatalf("buffer %q, over %v; want \"abcd\", true", c.buf.String(), c.over)
+	}
+}
+
+// The signer's child sees what ssh-keygen needs (SSH_AUTH_SOCK, HOME, PATH) and never the
+// registry credential (AC10, childenv.Scrubbed).
+func TestSignHandsTheSignerAScrubbedEnvironment(t *testing.T) {
+	dir := t.TempDir()
+	dump := filepath.Join(dir, "child.env")
+	bin := filepath.Join(dir, "signer")
+	if err := os.WriteFile(bin, []byte("#!/bin/sh\n/usr/bin/env > '"+dump+"'\ncat > /dev/null\necho signature\n"), 0o700); err != nil {
+		t.Fatal(err)
+	}
+	t.Setenv(childenv.CredentialVariable, "ac10-sentinel")
+	t.Setenv("SSH_AUTH_SOCK", "/tmp/agent.sock")
+	t.Setenv("HOME", dir)
+	out, err := Sign(t.Context(), bin, filepath.Join(dir, "k"), []byte("msg"), io.Discard)
+	if err != nil || string(out) != "signature\n" {
+		t.Fatalf("Sign = (%q, %v), want the signer's output", out, err)
+	}
+	b, err := os.ReadFile(dump)
+	if err != nil {
+		t.Fatal(err)
+	}
+	env := strings.Split(strings.TrimSpace(string(b)), "\n")
+	if childenv.Has(env, childenv.CredentialVariable) {
+		t.Fatalf("the signer observed %s", childenv.CredentialVariable)
+	}
+	for _, want := range []string{"SSH_AUTH_SOCK=/tmp/agent.sock", "HOME=" + dir} {
+		found := false
+		for _, kv := range env {
+			found = found || kv == want
+		}
+		if !found {
+			t.Errorf("the signer's environment lacks %s", want)
+		}
+	}
+}
````
