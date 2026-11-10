# Sprint plan — w-world-approvals-by-key-a1b, sprint 3 (iteration 377, row 201)

> **THE CUT.** This sprint lands the design's **M4 whole** and, behind a measured gate, the **core of M5**.
> - **Committed — Runs 1–4, 1,314 changed lines, measured.** M4a/M4b, the pure signature law
>   (sprint-2 plan §2.5/§2.8, carried byte for byte and re-verified at `545fdaae`); **G1**, the sprint-2
>   judge's store-flag test gap (F-CREDS, AS-CREDFOR, F-DB); **M4c1** (atomic: the consume binds the row's
>   challenge, authority derives `approved_by`, and `approve.go` stops naming it — B-D1) and **M4c2** its
>   tests; **M4d1/M4d2**, the store's approver binding (C2–C12); **M4e**, authority's signature binding
>   (A3/A4).
> - **Gated — Runs 5–6, 596 changed lines, measured.** **M5a** `approve <id> --signature <file>` with
>   `provisionEpisode`'s no-terminal branch; **M5b** its refusals and the revoke race; **M5c**
>   `approve <id> --write-message <file>` and the flag rules. They run **only if** the committed
>   `cmd/ host/` numstat after M4e is **≤ 1,404** (= 2,000 − 596), and **all three run or none**
>   (`--signature` without `--write-message` gives an operator no way to make the 90-byte message).
> - **Appendix C (sprint 4+), not built:** `DecidedMintRequests` + `approve --list --decided` + `approver
>   revoke`'s session listing (K8/K9, B-D4/B-D5); `--sign-with` and the design's M6 (real tool, CI step,
>   QUICKSTART §15's last part).
>
> **Every line of code and test in this sprint is a unified diff in Appendix P**, generated from a scratch
> tree that was compiled, vetted, gofmt-clean and green at every sub-milestone, and re-applied to a fresh
> export of `545fdaae` with `patch -p1` (§1.2 C30). The executor applies the patch and runs the
> acceptance; it writes no code by hand. The size column is therefore measured, not estimated: total
> **1,910** `cmd/ host/` lines; the committed part with the requested 1.2 safety factor is 1,577.

- **Design**: [w-world-approvals-by-key-a1b.md](w-world-approvals-by-key-a1b.md), **r2** (commit
  `5e2bf772`). **Not reopened.** Where this plan and the design differ, this plan was measured against the
  tree at `545fdaae` and wins; each difference is a decision (§1.3) or a premise found false (§1.4).
- **Starting ground**: Appendix B of [the sprint-2 plan](w-world-approvals-by-key-a1b-sprint2-plan.md)
  (B-D1…B-D6, the measured premises, the carried mutants) and that plan's §2.5, §2.8 and §4 M4 rows.
  Sprint 2 (PR #539, merged `290084ad`, recorded `545fdaae`) built M3 and stopped at its size gate (1,896 >
  1,683), so M4a/M4b were not built.
- **Planner base**: `545fdaaecf6f23dd58072675964754ad716049f9` (= origin/dev). Worktree
  `/Users/voightkampff/dev/sunholo-data/.wt-world-iter377-row201-a1b-s3`, branch
  `sprint/world-iter377-row201-a1b-s3`.
- **Executor**: as routed by the controller (iterations 375/376: `claude-haiku-5-5`). **Six runs**, about
  two sub-milestones each: Run 1 (M4a, M4b), Run 2 (G1, M4c1), Run 3 (M4c2, M4d1), Run 4 (M4d2, M4e),
  **size gate**, Run 5 (M5a, M5b), Run 6 (M5c). Each sub-milestone is one `patch -p1` (§2) plus its
  acceptance commands (§3). **Apply, never retype.**
- **Touched** (committed part): `host/approval/signed.go` + `signed_test.go` (new); `host/store/mint_requests.go`,
  `host/store/approvers.go`, `host/store/mint_requests_test.go` (the consume callers' challenge),
  `host/store/consume_binding_test.go` (new); `host/authority/mint.go`, `host/authority/mint_request_test.go`
  (the consume callers), `host/authority/mint_derive_test.go` (new); `cmd/ailang-worldd/approve.go`,
  `cmd/ailang-worldd/mint_request.go` (field-map wording), `cmd/ailang-worldd/approver_stores_test.go` (new),
  `cmd/ailang-worldd/approve_sig_test.go` (new). **Gated part adds**: `cmd/ailang-worldd/approve_sig.go` (new),
  `cmd/ailang-worldd/{session_new.go, main.go, approver.go, help_test.go}`, `docs/QUICKSTART.md` (§15,
  attended-pending, docs only).
- **Not touched**: `cmd/ailang-worldd/{session.go, approver_test.go, approver_enroll_test.go,
  approver_revoke_test.go, approver_rig_test.go, approve_test.go}`, `host/store/{credentials.go,
  credentials_schema.sql, approvers_test.go, credentials_test.go, context_roots_test.go}`, `host/approval/`
  except `signed*.go`, `host/approval/testdata/**`, `host/approval/approvaltest/**`, `host/daemon/**`,
  `.github/workflows/ci.yml` (M6), `world/*.ail`, `packages/**`, `scripts/**`, `tools/launchd/*`, any skill,
  `~/.ssh`, `~/.ailang/world`, `~/.ailang/state/mission-v1*`.

---

## 0. Binding constraints (each is restated where it governs)

1. **Fakes, temp stores and test keys only (D-WORLD-109 = A).** Every store is under `t.TempDir()`. Keys
   are `approvaltest.NewSigner(7|11|12)` (fixed seeds, no randomness) or the checked-in fixture
   `host/approval/testdata/sshsig/ed25519.pub`. **No new `ssh-keygen` invocation** in this sprint's code or
   tests (the planner ran none; §6's existing A1a cross-check is CI's, optional locally), no fixture is generated or edited, and nothing opens `~/.ailang/world`, `~/.ssh`, a real
   credentials store, a real daemon or `/dev/tty`. G1's default-`--db` arms set `HOME` to a temp dir
   (`t.Setenv`), as `newEpisodeRun` already does (`episode_test.go:86`). No mint outside a temp store, no
   publish. QUICKSTART §15's new text is attended-pending and **is not run**.
2. **The FIFO rule (P1, carried; B-D2).** `--signature` is read by the existing `readRegularFileCapped(path,
   16<<10)` (`approver.go:84`: `O_RDONLY|O_NONBLOCK`, `f.Stat()` on the descriptor, `io.LimitReader(max+1)`).
3. **Atomic typed path (B-D1, planner P6).** C1 (challenge required), C6 (`AND challenge = ?`) and A1
   (`ErrApprovedBySetByCaller`) each refuse every typed `approve` unless `approve.go`'s consume changes in
   the **same commit**, together with every existing consume caller in the tests. That is **M4c1**, one
   patch, one commit: `host/store/mint_requests.go`, `host/authority/mint.go`, `cmd/ailang-worldd/approve.go`,
   `cmd/ailang-worldd/mint_request.go`, `host/store/mint_requests_test.go` (`okConsume` and
   `TestConcurrentApprovalsMintOnce` gain `Challenge`), `host/authority/mint_request_test.go` (`consumeAt`
   loses `ApprovedBy`; three `RequestMint` literals gain `Challenge`). Never split it.
4. **One mutant per new branch** (§4: 102 counted, 102 killed in scratch; 8 known-equivalent named, not counted). Each row names
   its killing test, and each was run in the planner's scratch tree (§4 results column).
5. **Exit codes by constant**: `exitOK = 0`, `exitUsage = 1`, `exitFatal = 2` (`main.go:72-74`). Every
   refusal before the write is `exitUsage`; a failure at or after the writer open is `exitFatal` (CR6).
6. **Reflective completeness** for the structs a check claims to cover completely: `SignedApproval` (RC1,
   M4b), `RequestConsume` **and its element type `ApproverUse`** (RC3, M4d1), `authority.RequestMint` (RC5,
   M4e). RC3 and RC5 also perturb each field alone and require a refusal.
7. **No `t.Parallel()`** in any new test (they swap package vars: `fetchDaemonHealth`, `HOME`).
8. **Bisectable.** Each sub-milestone ends gofmt-clean, `go build ./...` and `go vet` green (the `_test.go`
   compile fence), with its own tests and every earlier sub-milestone's tests green (measured, §1.2 C31).
9. **No new context root.** `approveSignature` takes `runApprove`'s existing 3 s root as a parameter, so
   `TestProductionContextRoots` is **not edited** and stays green (measured). K8's root (B-D4) is Appendix C.
10. **No git writes by the executor** (§7). The controller commits from snapshots.

---

## 1. Pre-flight

### 1.1 Baseline at `545fdaae` (measured by the planner, 2026-11-10, pristine worktree)

Every command ran in the worktree with `export AILANG_BIN=$HOME/.pinned-ailang/ailang` in the same shell
(`AILANG v0.41.0`, commit `24ee108`), `go1.26.6 darwin/arm64`.

| Command | rc | Time | Result |
|---|---|---|---|
| `git status --short` | 0 | — | empty (clean) |
| `go build ./... && go vet ./...` | **0** | 6 s | clean |
| `go test -run '^$' ./...` (compile fence for every `_test.go`) | **0** | 29 s | no FAIL line |
| `go test -count=1 ./cmd/ailang-worldd/ ./host/approval/... ./host/store/ ./host/authority/` | **0** | 152 s | `cmd/ailang-worldd` ok 151.2 s, `host/approval` ok 1.1 s, `approvaltest` ok 0.3 s, `host/store` ok 12.4 s, `host/authority` ok 1.2 s; **no pre-existing failure** |
| `go test ./host/store/ -run '^TestProductionContextRoots$' -count=1` | **0** | 2 s | ok |
| `go test ./host/daemon/ -run '^TestDaemonDependencyAllowlist$' -count=1` | **0** | 1 s | ok |

Not run by the planner: `./scripts/verify_ail.sh` (no `.ail` is touched), the whole-repo `go test ./...`,
`verify_go.sh` and the full race suite (the controller's gates, §6). Known at this base from iterations
375/376: the `host/broker` flake `TestExecProbeNet6ForeignDialIsRetriedNotRefused` and local `-race`
package timeouts under load in `host/broker`, `host/daemon`, `host/verifygate` — untouched packages.

### 1.2 Codebase claims checked (each at `545fdaae` unless "scratch"; scratch = `git archive HEAD` exported to `~/.ailang/state/world-iter377/plan-scratch`, outside the repo)

| # | Claim | Command → observed output |
|---|---|---|
| C1 | The base and its state | `git rev-parse HEAD` → `545fdaaecf6f23dd58072675964754ad716049f9`; `git status --short \| wc -l` → `0` |
| C2 | Exit codes | `grep -n "exitOK *= 0\|exitUsage *= 1\|exitFatal *= 2" cmd/ailang-worldd/main.go` → `72: exitOK = 0`, `73: exitUsage = 1`, `74: exitFatal = 2` |
| C3 | The consume, as sprint 2's Appendix B measured it (unchanged by sprint 2) | `grep -n` in `host/store/mint_requests.go` → `24: challenge TEXT NOT NULL CHECK (length(challenge) = 64)`, `56: ErrIncompleteConsume … "needs an id, a digest, approved_by, a clock and a check"`, `237: type RequestConsume struct`, `250: if rc.ID == "" \|\| rc.Digest == "" \|\| rc.ApprovedBy == "" …`, `256` version check, `259: tx, err := c.db.BeginTx(ctx, nil)`, `264: at := rc.Clock()`, `268: UPDATE mint_requests SET … WHERE request_id = ? AND digest = ? AND consumed_at IS NULL AND created_at <= ? AND expires_at > ?` (**no challenge**), `290: func classifyConsume` |
| C4 | The challenge column holds any 64 characters, so a non-hex value can be stored | C3's line 24: the only CHECK is `length(challenge) = 64`; scratch: CR11 rewrites it to `strings.Repeat("zz", 32)` with raw SQL and the row reads back (test passes, C29) |
| C5 | The mint core | `grep -n` in `host/authority/mint.go` → `60: type RequestMint struct` (`Request, RequestCreatedAt, Consume`), `114: case rm != nil:`, `115: consume := rm.Consume`, `124: err = st.MintSessionForRequest(ctx, row, consume)`; `ApprovedBy` is passed through unread |
| C6 | `approve` today | `grep -n` in `cmd/ailang-worldd/approve.go` → `88: func runApprove`, `153: _ = ro.Close()` (the read-only handle closes **before** admission), `163: mr, refusal := approveAdmission(row, now)`, `191: health, err := fetchDaemonHealth(mr.Daemon.Addr)`, `212: req.answerBy = row.ExpiresAt`, `214: … ApprovedBy: "tty", Clock: env.now}}` |
| C7 | `provisionEpisode` | `grep -n` in `cmd/ailang-worldd/session_new.go` → `28: // cannot drift apart and no other code creates a credential.`, `403: func provisionEpisode`, `406: term, err := env.openTerminal()`, `446: fmt.Fprint(term, r.presetPrompt)`, `507: … authority.MintForRequest(mintCtx, …)` |
| C8 | **Every** consume caller (the M4c1 caller list) | `grep -rn "RequestConsume{\|RequestMint{\|okConsume(" --include="*.go" cmd host` → `approve.go:213`; `host/authority/mint_request_test.go:56` (`ApprovedBy: "tty"`), `:67`, `:109`, `:132` (`RequestMint{…}` literals); `host/store/mint_requests_test.go:353` (`okConsume`'s body), `:582` (`TestConcurrentApprovalsMintOnce`'s literal); and `okConsume(` callers `approvers_test.go:186`, `credentials_test.go:273`, `mint_requests_test.go:402,426,441,491,534` — all of which go through `okConsume`, so editing `okConsume` once covers them |
| C9 | The approver store's errors exist; `ErrApproverKeyChanged` is declared and unused | `grep -rn ErrApproverKeyChanged --include="*.go" host cmd` → only `host/store/approvers.go:48` (the declaration); control: `grep -rn ErrApproverRevoked … host/store` → used at `approvers_test.go:571` and others |
| C10 | `ApproverByKey` exists | `approvers.go:139: func (c *Credentials) ApproverByKey(ctx, blob []byte) (ApproverRow, bool, error)` |
| C11 | The FIFO-safe reader exists | `approver.go:84: func readRegularFileCapped(path string, max int64) ([]byte, error)` |
| C12 | **The judge's gap: no approver test passes `--credentials-db` or relies on the default `--db`** | `grep -c -- "--credentials-db" cmd/ailang-worldd/approver{,_enroll,_revoke,_rig}_test.go` → `0 0 0 0`; control `… approve_test.go` → `5`. `grep -c 'Setenv("HOME"'` on the same four files → `0 0 0 0`; control `episode_test.go` → `1` |
| C13 | The gap's mutants, at `545fdaae` (scratch, pristine tests, `go test ./cmd/ailang-worldd/ -run '^(TestApprover\|TestApprove\|TestEpisodeStartRequest)' -count=1`) | MUT-F-CREDS (`approver.go:188` returns `(*credsFlag)[:0]`) **SURVIVED**; MUT-AS-CREDFOR (`:203` `credentialsFor(db, creds[:0])`) **SURVIVED**; MUT-F-DB (`:174` default `""`) **SURVIVED**; MUT-A-DB (`approve.go:98` default `def.workspaceRoot`) **SURVIVED**; MUT-A-CREDS (`approve.go:125` `(*credsDB)[:0]`) killed already by `TestApproveRefusesBadIDsAndMisuse` (P7). With G1 applied, all five are KILLED by G1's tests (§4) |
| C14 | The approver flag sites | `approver.go:174: dbFlag := fs.String("db", defaultWorldPaths().db, …)`, `:188: return *dbFlag, *credsFlag, positional, 0, false`, `:203: path, err := credentialsFor(db, creds)` |
| C15 | `defaultWorldPaths().db` is `$HOME/.ailang/world/world.db` | `setup.go:74-84` (`os.UserHomeDir()` → `filepath.Join(home, ".ailang", "world", "world.db")`); scratch: G1's `homeStores` control asserts it under `t.Setenv("HOME", tmp)` |
| C16 | The approver help still says an enrolled key approves nothing | `approver.go:124: arrives with 'approve --signature' (row 201 A1b, a later sprint); until then an` (M5a rewords it; the help is pinned only by the constant, `grep -c approverHelp cmd/ailang-worldd/approver_test.go` → 9 references, no golden text) |
| C17 | The typed effect check hard-codes `tty` | `approve_test.go:1057: if f.ApprovedBy != "tty" {` (CR2 filters that one line and asserts `ssh:mark-mac` itself) |
| C18 | The consume's clock is the **third** `env.now` call on the approve path | precedent `approve_test.go:698` `clockSeq(2, approveT0+1, approveT0+601)` makes the consume's reading expire (CL9); CR6 revokes on call 3 and asserts `calls == 3` |
| C19 | A second connection can commit inside the consume's `Clock` hook | `credentials.go:293: db.SetMaxOpenConns(1)`, C3's deferred `BeginTx(ctx, nil)` before `at := rc.Clock()`; scratch: CS3 `revoked-in-the-window` and CR6 revoke through `rawDB`/`rawCredentials` from inside the hook and get `ErrApproverRevoked` (pass, C29) |
| C20 | `approvaltest` is imported only by test files | `grep -rln 'approval/approvaltest"' --include="*.go" .` → `host/approval/approvaltest/signer_test.go`; scratch: `TestApprovalTestIsTestOnly` passes with the two new importers (`mint_derive_test.go`, `approve_sig_test.go`) |
| C21 | `shellQuote` is production code | `episode.go:627: func shellQuote(s string) string` |
| C22 | The Mkfifo precedent | `approver_test.go:261: if err := syscall.Mkfifo(fifo, 0o600); err != nil {` |
| C23 | **M4a/M4b at this HEAD**: the sprint-2 §2.5/§2.8 code, extracted verbatim, compiles, vets, formats and passes | scratch: `awk` extraction of §2.5 → `signed.go` (112 lines), §2.8 → `signed_test.go` (211 lines); `gofmt -l host/approval` empty; `go vet ./host/approval/` clean; `go test ./host/approval/... -count=1` ok; `-v` → PASS `TestApproveBySignatureOrder`, `TestSignedApprovalBinds`, `TestSignedApprovalAccessorsFailClosed`, `TestSignedApprovalFieldsAreAllBound`, `TestApprovalImportsOnlyTheStandardLibrary`. The M4a-only stage (file cut at `// ---- M4b ----`, test imports `errors`, `testing`) also passes; its `signed.go` needs `gofmt` (a trailing blank line) — the M4a patch is the formatted text |
| C24 | Sprint 2's 22 M4 mutants still die at this HEAD | scratch: all 22 KILLED (MUT-P1 … MUT-RC1-FIELD; P4-ORDER as a real move of P4 after P6); MUT-RET-DIGEST equivalent (SURVIVED, as sprint 2 recorded) |
| C25 | The design's file positions sprint 2 did not move | sprint 2's diff (`git diff --stat 115f5a2b 545fdaae`) touches only `cmd/ailang-worldd/{approver*.go, help_test.go, keyguard.go, main.go}`, `host/approval/approver{,_test}.go`, `host/store/{approvers_test.go, context_roots_test.go}`, `docs/QUICKSTART.md`, plan docs — so Appendix B's premises in `mint_requests.go`, `mint.go`, `approve.go`, `session_new.go`, `mint_request.go` hold (re-measured: C3, C5, C6, C7) |
| C26 | Sprint sizes (calibration) | sprint 1 1,706 vs 1,541 planned (1.11×); sprint 2 `git diff --numstat 115f5a2b 545fdaae -- cmd host` → **1,896** vs 1,637 planned for M3 (1.16×) |
| C27 | `approveEnv`, `promptTerm`, `rawCredentials`, `seedApprover`, `revokeSeeded`, `approverRows`, `dbPathOf`, `rawDB`, `makeV4Credentials`, `testApprover`, `boundedTestContext` exist with the signatures the new tests call | `grep -n` → `approve_test.go:399 func approveEnv`, `session_new_test.go:462 type promptTerm`, `approve_test.go:489 func rawCredentials`, `approver_rig_test.go:93 seedApprover`, `:108 revokeSeeded`, `:76 approverRows`, `session_test.go:51 dbPathOf`, `schema_version_test.go:169 rawDB`, `approvers_test.go:16 makeV4Credentials`, `:37 testApprover`, `m7b_deadline_test.go:10 boundedTestContext` (and the scratch compile, C29) |
| C28 | No name collisions | scratch: every new identifier compiles in its package (`go vet` clean at every stage, C31) |
| C29 | **The whole sprint compiles and passes in scratch** | scratch at M5c: `gofmt -l cmd host` empty; `go build ./... && go vet ./...` clean; `go test ./cmd/ailang-worldd/ ./host/approval/... ./host/store/ ./host/authority/ -count=1` → all `ok` (`cmd/ailang-worldd` 125.2 s, `host/approval` 0.7 s, `approvaltest` 0.4 s, `host/store` 12.9 s, `host/authority` 0.9 s) — run on the fresh export after applying the patches **extracted from this file** (C30); `TestProductionContextRoots` ok (unedited); `TestDaemonDependencyAllowlist` ok; `-race` on the new host and cmd tests (cmd ones `-count=3`) ok |
| C30 | **The patches reproduce the tree.** | fresh `git archive 545fdaae` export + `patch -p1` of Appendix P's eleven patches in order → `diff -r` of `cmd host docs` against the scratch tree: identical; **the round trip was run on the patches extracted from this plan file itself** (§2's `awk`), not only on the generator's copies |
| C31 | **Every sub-milestone is bisectable** | the same fresh export, after **each** patch: `gofmt -l` empty, `go build ./... && go vet` of the touched packages clean, `go test -run '^$' ./...` rc 0, and the §3 acceptance tests rc 0 — 53 steps, all rc 0, 0 failures; the per-stage numstat printed after each patch was 221, 323, 511, 602, 815, 993, 1125, 1314, 1577, 1752, 1910 (= §7's expected values); `diff -r` of the final `cmd host docs` against the scratch tree: identical; no `.orig`/`.rej`; final `go vet ./...` clean; the `-race` legs and a `-v` run with 0 `--- SKIP` |
| C32 | The measured numstat equals the plan's column | `git diff --no-index --numstat` of the scratch stage trees (no repository written): after M4e `cmd` 234 + `host` 1,080 = **1,314**; after M5c `cmd` 830 + `host` 1,080 = **1,910** |

### 1.3 Planner decisions where the design or Appendix B leaves a choice

- **S3-D1 — The cut** (§1.5). Committed: M4a, M4b, G1, M4c1, M4c2, M4d1, M4d2, M4e (1,314). Gated,
  all-or-nothing: M5a, M5b, M5c (596), only if the committed numstat after M4e is ≤ 1,404. Appendix C:
  everything that reads decided rows (`DecidedMintRequests`, `--list --decided`, K8/K9 — B-D5 says they land
  together) and `--sign-with` with M6.
- **S3-D2 — Patches are the executor's form.** Iteration 376's gap sat in prose the executor had to turn
  into code. Here every change is a unified diff generated from the compiled scratch stages and applied
  with `patch -p1` (§2). The executor cannot paraphrase; the judge reviews exactly what was compiled.
- **S3-D3 — G1 covers every store flag of every touched verb.** `approver enroll|list|revoke`: `--db` +
  `--credentials-db` naming a store in **another directory** (a verb that ignored `--credentials-db` looks
  beside `--db` and finds nothing; `assertNoStoreBeside` also proves nothing was created there), and **no
  store flags at all** under `HOME=<tmp>` (the default `--db`). `approve --list`: the same two arms.
  `approve <id>`'s `--db`/`--credentials-db` are the rig's (`a.args()`), exercised by every CR test; the
  new `--signature`/`--write-message` each get arms (CR2–CR6, CR10), and every pair rule (F1, F2) gets one
  (CR1, CR3). The flags `enroll --name/--format/--pubkey` are sprint 2's, already armed (CA5).
- **S3-D4 — Split for size, never for the typed path.** M4c1 is the atomic typed-path unit (91 lines);
  its new tests are M4c2. M4d is split at its seam: M4d1 = `ApproverUse` + C2–C5 (validation, before the
  transaction) + RC3 + CS2; M4d2 = C8–C12 (`useApprover`, inside the transaction) + CS3 + CS4. Between
  them no production caller sets `RequestConsume.Approver`, so M4d1 alone changes no behaviour.
- **S3-D5 — A1 also refuses a caller-set approver.** The design's A1 checks only `ApprovedBy != ""`. M4e
  makes it `consume.ApprovedBy != "" || consume.Approver != nil`: the caller names neither (AU2
  `caller-approver`, MUT-A1-APPROVER). Without it the store's C4 would still refuse, but as
  `ErrInconsistentApproval`, one layer later.
- **S3-D6 — `approve --signature` looks the approver up on its own read-only handle**, under `runApprove`'s
  existing 3 s root, **after** admission and F4 (P3). No new context root (§0.9).
- **S3-D7 — The key the signature is verified against is the enrolled row's blob**
  (`Key: approval.SSHPublicKey{Blob: enrolled.PublicKey}`), not the signature's own key, so
  `VerifySSHSIG`'s byte equality (`sshsig.go:269`) compares the signature's key with the row (design §2.3
  "Which approver"). `SSHPublicKey.Type` is not read by verification (`sshsig.go:262-318` parses the blob).
  The swap to the signature's key is a named equivalent mutant (the lookup is by exact blob).
- **S3-D8 — `answerBy` stays `row.ExpiresAt` on both paths** (design §2.7 says 0 on the signature path).
  `provisionEpisode` reads no answer on the signed path (V2), so the only effect is its existing
  "expired before the prompt" refusal firing, before any lock, for a request that expired after admission
  — the in-transaction expiry check stays the binding one. One fewer branch, no untestable mutant.
- **S3-D9 — `--write-message`** runs after admission and F4 and before the same-store check, writes with
  `O_WRONLY|O_CREATE|O_EXCL` mode `0644`, and prints the request's **own** stores (`mr.WorldDB`,
  `mr.CredentialsDB`) in its next-step line, so the printed command passes approve's same-store check.
  The mode is not asserted (umask-dependent); it is named, not counted.
- **S3-D10 — Help and docs move with the verb.** M5a rewrites `approverHelp`'s "until then an enrolled key
  approves nothing" (it would become false) and adds the `--signature` usage line to `approveHelp`,
  `main.go`'s usage and `help_test.go`; M5c adds `--write-message` to all three and QUICKSTART §15's
  second paragraph (attended-pending, not run). If the gate stops M5, the sprint-2 text stays (still true).
- **S3-D11 — No CR13** (direct `provisionEpisode` call with a zero `SignedApproval`). V1's mutants are killed
  by CR2 (signed path opens no terminal) and CR3 (typed path still does) — measured — and a zero
  `SignedApproval` is refused at authority (AU2 `zero-approval`, RC5 `Signed`).
- **S3-D12 — RC3 holds C2's arms.** The design's CS1 C2 arms are RC3's three `ApproverUse` zeroing arms,
  each with a consistent `ApprovedBy` so only C2 can refuse it (each C2 sub-condition mutant is killed
  alone, measured).
- **S3-D13 — Refusal texts are fixed by the patches** and pinned by the tests (exact `stderr` where §3 says
  "exact", `Contains` of the stated fragment otherwise).
- **S3-D14 — Commit order inside the gated tail.** M5a's commit already ships a usable flag only together
  with M5c; the controller commits M5a, M5b, M5c as three commits but merges none of them without the
  others (they ride one PR with the committed part).

### 1.4 Premises found false or incomplete at HEAD (each routed around minimally)

- **P1 — (carried, B-D2) a plain open blocks on a FIFO.** Routed by reusing `readRegularFileCapped` for
  `--signature` (F9). CR4 `fifo` arm; MUT-F9 (`os.ReadFile`) hangs and dies on the test timeout.
- **P2 — Appendix B's size premise.** It estimated sprint 3 = "M4c–M4e + M5 ≈1,635 if M4a/M4b landed
  (≈1,950 if not — then M5 splits)". Measured with exact code: M4a+M4b **325**, M4c–M4e **803**, and M5
  without `--sign-with`, `--decided` and K8/K9 already **596**; the whole of Appendix B's sprint-3 scope
  does not fit under 2,000. Routed by S3-D1.
- **P3 — Design §2.7 asks for two things HEAD cannot give together.** "look up the approver … on the same
  read-only handle, before it closes (`approve.go:153` moves after this lookup)" **and** CR4 "consumed and
  expired requests → admission's existing refusals, before the signature is read". At HEAD the handle
  closes at `:153` and admission runs at `:163`; keeping the lookup on that handle means reading the
  signature before admission. Routed by S3-D6: admission first, then a second read-only handle under the
  same root. CR4 `consumed-first` and `expired-first` pass a path that does not exist and still get the
  admission text (MUT-SIG-ORDER/… in §4).
- **P4 — Design CR4's "a revoked approver (F14: 'revoked at')".** `ApproveBySignature`'s P4 returns
  `approval.ErrApproverRevoked`, whose text is `approval: the approver is revoked` with no time. The plan
  pins that text; the revocation time is in `approver list`.
- **P5 — Design §2.6 V34: "the two authority tests that pass `ApprovedBy: "tty"` (`mint_request_test.go:56`,
  `:78`)".** At HEAD only `:56` passes it (`consumeAt`); `:78` is the assertion
  `got.ApprovedBy != "tty"`, which stays. The three `RequestMint{…}` literals (`:67`, `:109`, `:132`) must
  also gain `Challenge` (C8). M4c1 edits exactly these four sites.
- **P6 — (carried, sprint-2 P6) the design's M4 → M5 split is not bisectable.** Routed by M4c1 (§0.3).
- **P7 — Sprint 2's judge counted approve's `--credentials-db` with the approver gap; it is not one.**
  MUT-A-CREDS is already killed at HEAD by `TestApproveRefusesBadIDsAndMisuse` (`other-store`); the
  real approve gap is the default `--db` (MUT-A-DB survives, C13). G1 kills both directly.
- **P8 — Design §2.6's A1 does not cover a caller-set `Approver`** (it predates `RequestConsume.Approver`
  being reachable from authority). Routed by S3-D5.
- **P9 — Design §2.7 F4 is needed exactly as written**: the column's CHECK is only `length = 64` (C4), so a
  non-hex challenge is storable and `[32]byte(challenge)` would panic without F4 (MUT-F4 dies by panic).
- **P10 — Design §2.8's sketch keeps `answerBy` at 0 on the signature path.** Not false, but unneeded
  (S3-D8).
- **P11 — `gofmt` aligns a trailing `// F10` comment with the previous line's `// F9`.** Routed by putting
  F9's comment on its own line (cosmetic, in the M5a patch).

### 1.5 The size, and the cut

Lines are insertions + deletions under `cmd/ host/` (`git diff --numstat`, the design V24 measure; docs
not counted). **Every row is measured** from the scratch stages (C32), not estimated.

| Sub-milestone | Production | Test | Step total | Cumulative vs `545fdaae` | Run |
|---|---|---|---|---|---|
| M4a — `ApproveBySignature`, `Binds`, `NotAdmittedError`, `BindingError` (SA1, SA2) | 94 | 127 | **221** | 221 | 1 |
| M4b — `ApprovedBy`, `Use` (SA3, RC1) | 18 | 86 | **104** | 323 | 1 |
| G1 — the store-flag gap (approver ×3 and `approve --list`, each × {`--credentials-db` elsewhere, default `--db`}) | 0 | 188 | **188** | 511 | 2 |
| M4c1 — **atomic**: C1, C6, C7; A1, A2, A5 + `RequestMint.Challenge`; `approve.go` F4 and the consume; the six caller edits | 72 | 19 | **91** | 602 | 2 |
| M4c2 — CS1, CS5 (store), AU1 typed arms (authority), CR11 (cmd) | 0 | 213 | **213** | 815 | 3 |
| M4d1 — `ApproverUse`, `ErrInconsistentApproval`, C2–C5 (CS2, RC3) | 23 | 155 | **178** | 993 | 3 |
| M4d2 — C8–C12: `useApprover`, `classifyApproverUse` (CS3, CS4) | 47 | 85 | **132** | 1,125 | 4 |
| M4e — `RequestMint.{ExpiresAt, Signed}`, A1′, A3, A4 (AU1 signed, AU2, RC5) | 25 | 176 | **201** | **1,314** | 4 |
| **Committed** (numstat) | **267** | **1,047** | **1,314** | ×1.2 = 1,577 ≤ 2,000 | |
| *gate: run 5–6 only if the measured numstat after the M4e commit is ≤ 1,404* | | | | | |
| M5a — `approve --signature`, `approveSignature`, V1/V2, help (CR2) | 171 | 96 | **267** | 1,577 | 5 |
| M5b — CR3, CR4/CR5 (12 arms), CR6 | 0 | 175 | **175** | 1,752 | 5 |
| M5c — `--write-message`, F1/F2, QUICKSTART §15 (docs, not counted) (CR10, CR1) | 53 | 109 | **162** | **1,910** | 6 |
| **Sprint** (numstat) | **485** | **1,425** | **1,910 ≤ 2,000** | | |
| *Appendix C (sprint 4): `DecidedMintRequests` + `--list --decided` + K8/K9 + CS6, CR7, CR12* | ≈90 | ≈260 | ≈350 | estimate ×1.2 ≈ 420 | |
| *Appendix C (sprint 5, M6): `--sign-with` seam + `signWithSSHKeygen` + F7/F7a/F8/G1–G3 + XS1/XS2 + CI + QUICKSTART* | ≈110 | ≈400 | ≈510 (+20 YAML, +25 docs) | estimate ×1.2 ≈ 610 | |

(The step totals sum to 1,932; the cumulative column is 1,910 because a few lines are edited twice —
`mint.go`'s `RequestMint` and `requestConsume`, `approve.go`'s flag rules — which numstat against the base
counts once.)

---

## 2. How the executor applies a sub-milestone (the patch protocol)

Each sub-milestone `<M>` is the block fenced by a line `` ````diff <M> `` (four backticks: M5c's QUICKSTART hunk itself contains three-backtick fences) and a line of four backticks, in Appendix P. Paths are `a/<repo path>` and
`b/<repo path>`. Apply it **only** this way, in one `bash` shell (the controller's shell is zsh):

```bash
bash <<'EOF'
set -euo pipefail
M=M4a   # <- the sub-milestone, exactly as in the Appendix P heading
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter377-row201-a1b-s3
PLAN="$W/design_docs/planned/w-world-approvals-by-key-a1b-sprint3-plan.md"
T="$HOME/.ailang/state/world-iter377/exec-patches"; mkdir -p "$T"
awk -v m="$M" '$0 == "````diff " m {f=1; next} f && $0 == "````" {exit} f' "$PLAN" > "$T/$M.patch"
test -s "$T/$M.patch"
cd "$W"
patch -p1 --dry-run < "$T/$M.patch"
patch -p1 < "$T/$M.patch"
if find cmd host docs \( -name '*.orig' -o -name '*.rej' \) | grep .; then echo "STOP: patch left .orig/.rej"; exit 1; fi
git status --short
EOF
```

- If `--dry-run` fails, **stop and report** (do not edit by hand, do not use fuzz options). It can only
  fail if an earlier sub-milestone was skipped or changed.
- Then run that sub-milestone's **Acceptance** (§3) in one bash shell with
  `export AILANG_BIN=$HOME/.pinned-ailang/ailang`.
- Then snapshot (§7). Then compare the numstat with §7's expected value.

---

## 3. Milestones

Every sub-milestone's acceptance starts with the same three lines (abbreviated **[base]** below):

```
test -z "$(gofmt -l cmd/ailang-worldd host/approval host/store host/authority)"
go build ./... && go vet ./cmd/ailang-worldd/ ./host/approval/... ./host/store/ ./host/authority/
go test -run '^$' ./...
```

`CMDT` below is
`'^(TestApprove|TestEpisodeStartRequest|TestMintRequest|TestApprover|TestEveryVerbHasHelp|TestHelpIsNeverBlocked|TestExistingCallersReadUnbounded|TestReadTermLine|TestConfirmApproverAct|TestReadRegularFileCapped)'`.

### M4a — the signature law · 94 prod + 127 test (measured) · Run 1

**Patch** `M4a`: new `host/approval/signed.go` = sprint-2 §2.5 down to `Binds` (`SignedApproval`
{ok, digest, challenge, expiry, name, format, key}, `SignatureInput`, `NotAdmittedError` + `Error()`,
`ErrNotVerified`, `BindingError` + `Error()`, `ApproveBySignature` P1–P6, `Binds` B1–B4); new
`host/approval/signed_test.go` = sprint-2 §2.8's M4a block (`sigInput`, `verifiedApproval`, SA1
`TestApproveBySignatureOrder` 18 arms + `enrolled-at-now` + `clean`, SA2 `TestSignedApprovalBinds`), imports
`errors`, `testing`. Pure, standard library only (`bytes`, `errors`, `fmt`); no production caller yet.

**Acceptance**
```
[base]
go test ./host/approval/ -run '^(TestApproveBySignatureOrder|TestSignedApprovalBinds|TestApprovalImportsOnlyTheStandardLibrary)$' -count=1 -v
go test ./host/approval/... -count=1
```

### M4b — the accessors and the completeness check · 18 prod + 86 test · Run 1

**Patch** `M4b`: `ApprovedBy` (B5) and `Use` (B6, returns a copy of the key) after `// ---- M4b ----`;
SA3 `TestSignedApprovalAccessorsFailClosed`, `signedApprovalFields`, RC1 `TestSignedApprovalFieldsAreAllBound`;
test imports gain `bytes`, `fmt`, `reflect`.

**Acceptance**
```
[base]
go test ./host/approval/ -run '^(TestApproveBySignatureOrder|TestSignedApprovalBinds|TestSignedApprovalAccessorsFailClosed|TestSignedApprovalFieldsAreAllBound)$' -count=1 -v
go test ./host/approval/... -count=1 && go test -race ./host/approval/... -count=1
```

### G1 — the store flags every approver verb and `approve --list` use · 0 prod + 188 test · Run 2

**Patch** `G1`: new `cmd/ailang-worldd/approver_stores_test.go`:
- helpers `elsewhereStores(t) (db, creds)` (world store + a schema-5 credentials store in another temp dir),
  `homeStores(t) (db, creds)` (`t.Setenv("HOME", tmp)`, the default world store and the store beside it;
  asserts `defaultWorldPaths().db` equals it as a control), `assertNoStoreBeside(t, db)`,
  `putPendingRequest(t, creds)` (`firstRequestID`, episode `ep1`, created `approveT0`, expires `+600`).
- `TestApproverVerbsUseTheNamedStores`: six subtests `list/`, `enroll/`, `revoke/` × `credentials-db`
  (flags `--db db --credentials-db creds`) and `default-db` (no store flags). `list` → stdout exactly
  `mark-mac  ssh  ssh-ed25519 <edFingerprint>  enrolled 2027-01-15T08:00:00Z  last used never  active\n`;
  `enroll` (y) → the row is in `creds`; `revoke` (y) → `RevokedAt == approveT0+100` in `creds`; the
  `credentials-db` arms assert nothing was created beside `db`.
- `TestApproveListUsesTheNamedStores`: `credentials-db` and `default-db` → stdout exactly
  `mr-0102030405060708  episode ep1  digest abababababababab  expires 2027-01-15T08:10:00Z\n`.

**Acceptance**
```
[base]
go test ./cmd/ailang-worldd/ -run '^(TestApproverVerbsUseTheNamedStores|TestApproveListUsesTheNamedStores)$' -count=1 -v
go test ./cmd/ailang-worldd/ -run "$CMDT" -count=1
```

### M4c1 — ATOMIC: the challenge binding and `approved_by` derivation, with every caller · 72 prod + 19 test · Run 2

**One patch, one commit** (§0.3). Patch `M4c1`:
- `host/store/mint_requests.go`: `ErrIncompleteConsume`'s text gains "a challenge"; new
  `ErrRequestChallengeChanged = errors.New("store: the mint request's challenge is not the one signed")`;
  `RequestConsume.Challenge string`; **C1** `!isLowerHex64(rc.Challenge)` in the completeness `if`; **C6**
  `AND challenge = ?` (+ `rc.Challenge` argument) in the guarded `UPDATE`; **C7** `classifyConsume` reads
  `challenge` and adds `case challenge != rc.Challenge:` after the digest case.
- `host/authority/mint.go`: imports `errors`; `RequestMint.Challenge [32]byte`; `ErrApprovedBySetByCaller`
  (`"authority: approved_by is derived by the mint, never set by the caller"`), `ErrBadConsumeDigest`
  (`"authority: the consume's digest is not 64 hex"`); new `requestConsume(rm *RequestMint)
  (store.RequestConsume, error)`: **A1** `ApprovedBy != ""` → refuse, **A2** digest hex-decodes to 32 bytes
  or refuse, `Challenge = hex(rm.Challenge)`, **A5** `ApprovedBy = "tty"`; `mintRecorded`'s consume branch
  calls it.
- `cmd/ailang-worldd/approve.go`: **F4** after admission — `hex.DecodeString(row.Challenge)`, 32 bytes or
  `refuse("request %s's stored challenge is not 64 hex; refusing; fix: make a new request with episode start --request")`;
  the consume passes `Challenge: [32]byte(challenge)` and **no `ApprovedBy`**.
- `cmd/ailang-worldd/mint_request.go`: the `"consume"` field-map text.
- Tests (the caller list, C8): `host/store/mint_requests_test.go` `okConsume` gains `Challenge: row.Challenge`
  (covers its callers in `approvers_test.go:186`, `credentials_test.go:273` and five in the same file), and
  `TestConcurrentApprovalsMintOnce`'s literal gains `Challenge: r.Challenge`;
  `host/authority/mint_request_test.go` `consumeAt` drops `ApprovedBy: "tty"`, the three `RequestMint`
  literals gain `Challenge: au1Challenge`, and `au1Challenge` (32 × `0xab`, the row's `"ab"×32`) is declared.

**Acceptance** (the typed path must still be green end to end)
```
[base]
go test ./host/store/ ./host/authority/ -count=1
go test ./cmd/ailang-worldd/ -run "$CMDT" -count=1
go test ./cmd/ailang-worldd/ -run '^(TestApproveMintsExactlyTheRequest|TestApproveRefusesReplay|TestApproveRefusesExpired|TestApproveListShowsPending)$' -count=1 -v
```

### M4c2 — the tests of M4c1's branches · 0 prod + 213 test · Run 3

**Patch** `M4c2`:
- new `host/store/consume_binding_test.go`: `putTestRequest`, `assertPendingNoSession`; **CS1**
  `TestMintSessionForRequestNeedsTheChallenge` (`empty`, `63-hex`, `65-hex`, `upper`, `not-hex` →
  `ErrIncompleteConsume`, clock unread, request pending, 0 sessions); **CS5** `TestConsumeBindsTheChallenge`
  (`other-challenge` → `ErrRequestChallengeChanged`; `other-digest-and-challenge` → `ErrRequestDigestChanged`).
- new `host/authority/mint_derive_test.go`: `au1Store`, `au1RequestMint`, `au1Mint`, `assertNotMinted`;
  **AU1** `TestMintForRequestDerivesApprovedBy` (`caller-tty`, `caller-ssh` → `ErrApprovedBySetByCaller`;
  `digest-not-hex`, `digest-31-bytes` → `ErrBadConsumeDigest`; `other-challenge` →
  `store.ErrRequestChallengeChanged`; each not minted; `typed-records-tty` → consumed by `tty`).
- new `cmd/ailang-worldd/approve_sig_test.go`: **CR11** `TestApproveRefusesANonHexChallenge` (raw SQL sets
  the challenge to `"zz"×32` → `exitUsage`, exact stderr, no terminal, pending, 0 sessions).

**Acceptance**
```
[base]
go test ./host/store/ -run '^(TestMintSessionForRequestNeedsTheChallenge|TestConsumeBindsTheChallenge)$' -count=1 -v
go test ./host/authority/ -run '^TestMintForRequestDerivesApprovedBy$' -count=1 -v
go test ./cmd/ailang-worldd/ -run '^TestApproveRefusesANonHexChallenge$' -count=1 -v
go test ./host/approval/... ./host/store/ ./host/authority/ -count=1 && go test ./cmd/ailang-worldd/ -run "$CMDT" -count=1
```

### M4d1 — the approver a consume binds: validation (C2–C5) · 23 prod + 155 test · Run 3

**Patch** `M4d1`:
- `host/store/mint_requests.go`: `ErrInconsistentApproval = errors.New("store: approved_by does not match the approver the consume binds")`;
  `RequestConsume.Approver *ApproverUse` (nil = the typed y); `type ApproverUse struct{ Name, Format string;
  PublicKey []byte }`; **C2** an approver with an empty name, format or key → `ErrIncompleteConsume`; **C3**
  no approver and `ApprovedBy != "tty"` → `ErrInconsistentApproval`; **C4** an approver and
  `ApprovedBy != Format+":"+Name` → `ErrInconsistentApproval`; **C5** an approver on a store below
  version 5 → `c.olderForApprovers()` (`Need 5, Feature "approvers"`). All before `BeginTx`.
- `host/store/consume_binding_test.go` after `// ---- M4d1 ----`: `credentialsAt` (store + approver `x`,
  `testApprover("x", 900)`), `signedConsume`, `approverX`; **CS2** `TestApprovedByMatchesTheApprover`
  (`no-approver-ssh`, `approver-tty`, `approver-other-name`, `approver-other-format`,
  `approver-format-not-ssh` → `ErrInconsistentApproval`, clock unread; `v4-store` → `*OlderCredentialsStoreError{Need 5, "approvers"}`);
  `requestConsumeFields`; **RC3** `TestRequestConsumeCompleteness` (reflects `RequestConsume` + `ApproverUse`,
  10 dispositions; zeroes each of the 8 "required" fields in a valid signed consume → `ErrIncompleteConsume`).

**Acceptance**
```
[base]
go test ./host/store/ -run '^(TestApprovedByMatchesTheApprover|TestRequestConsumeCompleteness)$' -count=1 -v
go test ./host/approval/... ./host/store/ ./host/authority/ -count=1 && go test ./cmd/ailang-worldd/ -run "$CMDT" -count=1
```

### M4d2 — the approver a consume binds: inside the transaction (C8–C12) · 47 prod + 85 test · Run 4

**Patch** `M4d2`:
- `host/store/approvers.go`: imports `bytes`; `useApprover(ctx, tx, a ApproverUse, at int64) error` — **C8**
  `UPDATE approvers SET last_used_at = ? WHERE name = ? AND format = ? AND public_key = ? AND revoked_at IS NULL AND enrolled_at <= ?`
  with the consume's reading; not exactly one row → `classifyApproverUse`: **C9** absent →
  `ErrApproverUnknown`, **C10** revoked → `ErrApproverRevoked`, **C11** format or key differs →
  `ErrApproverKeyChanged`, **C12** `enrolled_at > at` → `ErrApproverClockBehind`.
- `host/store/mint_requests.go`: after the request `UPDATE` and before the session insert,
  `if rc.Approver != nil { useApprover(...) }` (an error rolls everything back).
- `host/store/consume_binding_test.go` after `// ---- M4d2 ----`: **CS3**
  `TestSignedConsumeChecksTheApproverInTheTransaction` (`active` → `ssh:x`, `last_used_at == 1010`, 1 session;
  `revoked-in-the-window` (a second connection revokes `x` from inside the `Clock` hook, C19),
  `unknown`, `other-key`, `other-format`, `enrolled-after-the-reading` → C10/C9/C11/C11/C12, each pending,
  0 sessions, `last_used_at` 0); **CS4** `TestTypedConsumeTouchesNoApprover`.

**Acceptance**
```
[base]
go test ./host/store/ -run '^(TestSignedConsumeChecksTheApproverInTheTransaction|TestTypedConsumeTouchesNoApprover)$' -count=1 -v
go test -race ./host/store/ -run '^(TestSignedConsumeChecksTheApproverInTheTransaction|TestConcurrentApprovalsMintOnce)$' -count=1
go test ./host/approval/... ./host/store/ ./host/authority/ -count=1 && go test ./cmd/ailang-worldd/ -run "$CMDT" -count=1
```

### M4e — authority binds a signature approval (A3, A4) · 25 prod + 176 test · Run 4

**Patch** `M4e`:
- `host/authority/mint.go`: `RequestMint.ExpiresAt int64`, `RequestMint.Signed *approval.SignedApproval`;
  `requestConsume`: **A1′** also refuses `consume.Approver != nil` (S3-D5); **A5** `tty` only when
  `rm.Signed == nil`; **A4** `rm.Signed.Binds(digest, rm.Challenge, rm.ExpiresAt)` or its error; **A3**
  `ApprovedBy = rm.Signed.ApprovedBy()`, `Approver = &store.ApproverUse{name, format, key}` from `Use()`.
- `host/authority/mint_derive_test.go` after `// ---- M4e ----`: `au1Signer` (`approvaltest.NewSigner(7)`),
  `au1Enrol`, `au1Signed`, `au1SignedMint`, `approverLastUsed`; **AU1** `TestMintForRequestRecordsTheSigner`
  (`ssh:mark-mac`, `last_used_at == 1020`); **AU2** `TestMintForRequestRefusesAnUnboundApproval`
  (`other-digest`, `other-challenge`, `other-expiry`, `expiry-not-passed` → `*approval.BindingError` naming
  the field; `zero-approval` → `approval.ErrNotVerified`; `caller-approver` → `ErrApprovedBySetByCaller`;
  each not minted, approver unused); `requestMintFields`; **RC5** `TestRequestMintFieldsAreSourced`
  (6 fields, each perturbed alone in a valid signed `RequestMint` → refused, not minted).

**Acceptance**
```
[base]
go test ./host/authority/ -run '^(TestMintForRequestDerivesApprovedBy|TestMintForRequestRecordsTheSigner|TestMintForRequestRefusesAnUnboundApproval|TestRequestMintFieldsAreSourced)$' -count=1 -v
go test ./host/approval/... -run '^TestApprovalTestIsTestOnly$' -count=1
go test ./host/approval/... ./host/store/ ./host/authority/ -count=1 && go test ./cmd/ailang-worldd/ -run "$CMDT" -count=1
go test -race ./host/approval/... ./host/store/ ./host/authority/ -count=1
```

**→ the size gate (§7) runs here.**

### M5a — `approve <id> --signature <file>` · 171 prod + 96 test · Run 5 (gated)

**Patch** `M5a`:
- new `cmd/ailang-worldd/approve_sig.go`: `approveMaxSignatureBytes = 16 << 10`;
  `approveSignature(ctx, credPath, path, row, now) (*approval.SignedApproval, string)`: **F9**
  `readRegularFileCapped`, **F10** read error → `"--signature: %v"`, **F11** `ParseSSHSIG` error →
  `"--signature %s is not an SSH signature: %v"`, the approver by the signature's key blob on a read-only
  handle (S3-D6), **F12** its error (below schema 5 names the restart), **F13** none →
  `"the signature's key %s is not an enrolled approver; see ailang-worldd approver list"`, then
  `approval.ApproveBySignature` over the stored row with the **row's** key (S3-D7), **F14** its error →
  `"the signature does not approve request %s: %v"`.
- `cmd/ailang-worldd/approve.go`: the help (usage line + a `--signature` paragraph), the `--signature` flag,
  **F2** `--list` with `--signature` → help + `exitUsage`; after F4, `approveSignature` (before the
  same-store check and any network call); the consume carries `ExpiresAt: row.ExpiresAt, Signed: signed`.
- `cmd/ailang-worldd/session_new.go` `provisionEpisode`: `signed := r.consume != nil && r.consume.Signed != nil`;
  **V1** the terminal is opened only when `!signed`; **V2** when signed, the request block and
  `"<action> and mint a session … expiry +<ttl>s: approved by <ssh:name>\n"` go to stdout and no answer is
  read; the typed branch is unchanged (re-indented). The step comment at `:28` names the signed path.
- `main.go` usage, `help_test.go` want, `approver.go`'s `approverHelp` (S3-D10), `mint_request.go` field map.
- `approve_sig_test.go` after `// ---- M5a ----`: `sigSigner` (`NewSigner(11)`), `otherSigner`
  (`NewSigner(12)`), `requestMessage`, `writeSignature`, `signedRig`; **CR2**
  `TestApproveBySignatureMintsWithNoTerminal` (no terminal opened, `ssh:mark-mac`, `last_used_at ==
  approveT0+1`, stdout has `mint request <id>\n` and `expiry +3600s: approved by ssh:mark-mac\n`, and every
  CL16 effect but `approved_by:` matches the request).

**Acceptance**
```
[base]
go test ./cmd/ailang-worldd/ -run '^(TestApproveBySignatureMintsWithNoTerminal|TestEveryVerbHasHelp|TestHelpIsNeverBlockedByTheRegistryKeyGuard|TestMintRequestFieldMapIsComplete)$' -count=1 -v
go test ./cmd/ailang-worldd/ -run "$CMDT" -count=1
go test ./cmd/ailang-worldd/ -run '^(TestSessionNew|TestSessionMint_|TestEpisodeStart)' -count=1
go test ./host/store/ -run '^TestProductionContextRoots$' -count=1
```

### M5b — the refusals and the revoke race · 0 prod + 175 test · Run 5 (gated)

**Patch** `M5b`, `approve_sig_test.go` after `// ---- M5b ----`: **CR3**
`TestApproveTypedPathUnchangedWithAnApprover` (the typed path with an approver enrolled opens the terminal,
records `tty`, leaves `last_used_at` 0; `--list --signature` → `exitUsage` + exactly the help); **CR4/CR5**
`TestApproveBySignatureRefusalsBeforeAnyEffect`, 12 arms, each `exitUsage`, no terminal, **0 daemon calls**
(a counting wrapper around `fetchDaemonHealth`), never `ssh:mark-mac`: `fifo`, `over-16k`, `missing`,
`garbage`, `unenrolled-key`, `revoked`, `other-request`, `enrolled-after-now`, `git-namespace`, `v4-store`,
`consumed-first`, `expired-first` (the last two pass a path that does not exist and get admission's
text); **CR6** `TestApproveRevokeRaceIsRefusedAtTheWrite` (the consume's clock reading — the third
`env.now` call — revokes the approver on a second connection → `exitFatal`, `store: the approver is
revoked`, `removed the worktree`, pending, 0 sessions, worktree gone, approver revoked and unused).

**Acceptance**
```
[base]
go test ./cmd/ailang-worldd/ -run '^(TestApproveTypedPathUnchangedWithAnApprover|TestApproveBySignatureRefusalsBeforeAnyEffect|TestApproveRevokeRaceIsRefusedAtTheWrite)$' -count=1 -v
go test -race ./cmd/ailang-worldd/ -run '^(TestApproveBySignature|TestApproveRevokeRace)' -count=1
go test ./cmd/ailang-worldd/ -run "$CMDT" -count=1
```

### M5c — `approve <id> --write-message <file>` and the flag rules · 53 prod + 109 test (+ docs) · Run 6 (gated)

**Patch** `M5c`:
- `approve_sig.go`: `writeApprovalMessage(stdout, path, row, mr) string`: `SignedMessage` of the stored
  digest, challenge and expiry; **F5** `os.OpenFile(path, O_WRONLY|O_CREATE|O_EXCL, 0o644)` error →
  `"--write-message: %v; nothing was written"`; writes the 90 bytes; prints exactly three lines (S3-D9).
- `approve.go`: the `--write-message` flag and help; **F1** `--signature` with `--write-message` → help +
  `exitUsage`; **F2** widened to `--list` with either; after F4, `--write-message` writes and returns
  `exitOK` (**F6**: no terminal, no store row, no network).
- `main.go` usage, `help_test.go` want; `docs/QUICKSTART.md` §15: the "arrives in a later sprint" sentence
  replaced, a second bash block (`episode start --request`, `approve --write-message`, `ssh-keygen -Y sign`
  of the scratch key, `approve --signature`, then `approver revoke`), and the refusal paragraph extended.
  **Docs only, attended-pending, not run.**
- `approve_sig_test.go` after `// ---- M5c ----`: **CR10** `TestApproveWriteMessage` (`writes-then-signs`:
  the file is exactly `requestMessage`, stdout exactly the three lines, no terminal, 0 daemon calls,
  pending; signing the file's bytes then approves as `ssh:mark-mac`; `existing-file` → exact stderr, file
  kept; `used-request` → admission text, no file); **CR1** `TestApproveFlagRules`
  (`signature-and-write-message`, `list-and-signature`, `list-and-write-message` → `exitUsage`, stderr
  exactly `approveHelp`, no credentials store touched).

**Acceptance**
```
[base]
go test ./cmd/ailang-worldd/ -run '^(TestApproveWriteMessage|TestApproveFlagRules|TestEveryVerbHasHelp)$' -count=1 -v
go test ./cmd/ailang-worldd/ -run "$CMDT" -count=1
go test ./cmd/ailang-worldd/ -run 'Quickstart|QuickStart|QUICKSTART' -count=1
```

---

## 4. Mutation table (102 counted mutants, 102 killed in the planner's scratch run; the controller and judge run them)

**Derivation.** One row per new production `if`, `case`, guard clause, SQL guard, copied field and
dispatch in Appendix P's production hunks, plus sprint 2's 22 M4 rows and the five store-flag mutants
G1 exists to kill. **Procedure** (as sprint 2 §4): copy the file, apply **one** edit, **check the file
changed** (`cmp` must differ), run the command, expect FAIL (red, panic or timeout), restore, `cmp` equal.
A mutant must compile (`false && <cond>`, `x[:0]`, `0*x`); the planner's harness ran `go vet` first and
recorded NOCOMPILE separately (three first drafts did not compile and two anchors had moved by M4e;
each was rewritten and re-run, and only the re-run counts). `export AILANG_BIN=$HOME/.pinned-ailang/ailang`. **Never commit or snapshot a mutant.**

Command keys: **A** = `go test ./host/approval/ -run '^<test>$' -count=1`; **S** = `go test ./host/store/ -run
'^<test>$' -count=1`; **U** = `go test ./host/authority/ -run '^<test>$' -count=1`; **C** = `go test
./cmd/ailang-worldd/ -run '^<test>$' -count=1` (MUT-F9: add `-timeout 90s`; it hangs on the FIFO by design).
Result = the planner's scratch run at the stated stage (all against the final M5c tree except the C13
controls).

| Id | Stage | Mutant (single edit) | Killing test (arm) | Cmd | Scratch |
|---|---|---|---|---|---|
| MUT-P1 | M4a | `signed.go` P1: `v != Admitted` → `false && v != Admitted` | SA1 `consumed`, `expired`, … | A `TestApproveBySignatureOrder` | KILLED |
| MUT-P1-ORDER | M4a | P1's `if` moved after P4 | SA1 `consumed-and-revoked` | A `TestApproveBySignatureOrder` | KILLED |
| MUT-P2 | M4a | P2's `if` → `if false && …` | SA1 `name-upper`, `name-empty` | A `TestApproveBySignatureOrder` | KILLED |
| MUT-P3 | M4a | P3's `if` → `if false && …` | SA1 `format-webauthn` | A `TestApproveBySignatureOrder` | KILLED |
| MUT-P4 | M4a | P4's `if` → `if false && …` | SA1 `revoked-good-signature` | A `TestApproveBySignatureOrder` | KILLED |
| MUT-P4-ORDER | M4a | P4's `if` moved after P6 | SA1 `revoked-bad-signature` | A `TestApproveBySignatureOrder` | KILLED |
| MUT-P5 | M4a | P5's `if` → `if false && …` | SA1 `enrolled-after-now` | A `TestApproveBySignatureOrder` | KILLED |
| MUT-P5-LE | M4a | P5 `<` → `<=` | SA1 `enrolled-at-now` | A `TestApproveBySignatureOrder` | KILLED |
| MUT-P6 | M4a | P6 `err != nil` → `false && err != nil` | SA1 `git-namespace`, `p256-signature`, `flipped-signature` | A `TestApproveBySignatureOrder` | KILLED |
| MUT-P6-CHALLENGE | M4a | `in.Challenge` → `[32]byte{}` in the `VerifyApproval` call | SA1 `clean` | A `TestApproveBySignatureOrder` | KILLED |
| MUT-P6-EXPIRY | M4a | `in.Admit.ExpiresAt` → `in.Admit.CreatedAt` in the call | SA1 `clean` | A `TestApproveBySignatureOrder` | KILLED |
| MUT-B1 | M4a | B1 `s == nil \|\| !s.ok` → `s == nil` | SA2 `zero` | A `TestSignedApprovalBinds` | KILLED |
| MUT-B1-NIL | M4a | B1 → `!s.ok` | SA2 `nil` (panic) | A `TestSignedApprovalBinds` | KILLED |
| MUT-B2 | M4a | B2's `if` → `if false && …` | SA2 digest | A `TestSignedApprovalBinds` | KILLED |
| MUT-B3 | M4a | B3's `if` → `if false && …` | SA2 challenge | A `TestSignedApprovalBinds` | KILLED |
| MUT-B4 | M4a | B4's `if` → `if false && …` | SA2 expiry ±1 | A `TestSignedApprovalBinds` | KILLED |
| MUT-A-CLONE | M4a | `key: bytes.Clone(in.Approver.Key.Blob)` → `key: in.Approver.Key.Blob` | SA3 (the input edit shows through) | A `TestSignedApprovalAccessorsFailClosed` | KILLED |
| MUT-RET-NAME | M4a | `name: in.Approver.Name` → `name: ""` | SA3 `ApprovedBy` | A `TestSignedApprovalAccessorsFailClosed` | KILLED |
| MUT-B5 | M4b | B5 → `s == nil` | SA3 `zero` | A `TestSignedApprovalAccessorsFailClosed` | KILLED |
| MUT-B6 | M4b | B6 → `s == nil` | SA3 `zero` | A `TestSignedApprovalAccessorsFailClosed` | KILLED |
| MUT-B6-CLONE | M4b | `bytes.Clone(s.key)` → `s.key` in `Use` | SA3 (copies) | A `TestSignedApprovalAccessorsFailClosed` | KILLED |
| MUT-RC1-FIELD | M4b | add `extra string` to `SignedApproval` | RC1 (no disposition) | A `TestSignedApprovalFieldsAreAllBound` | KILLED |
| MUT-F-CREDS | G1 | `approver.go:188` → `return *dbFlag, (*credsFlag)[:0], positional, 0, false` | G1 `*/credentials-db` | C `TestApproverVerbsUseTheNamedStores` | KILLED |
| MUT-AS-CREDFOR | G1 | `approver.go:203` → `credentialsFor(db, creds[:0])` | G1 `*/credentials-db` | C `TestApproverVerbsUseTheNamedStores` | KILLED |
| MUT-F-DB | G1 | `approver.go:174` default `defaultWorldPaths().db` → `""` | G1 `*/default-db` | C `TestApproverVerbsUseTheNamedStores` | KILLED |
| MUT-A-CREDS | G1 | `approve.go:125` → `credentialsFor(*db, (*credsDB)[:0])` | G2 `credentials-db` (also CL12 at HEAD, P7) | C `TestApproveListUsesTheNamedStores` | KILLED |
| MUT-A-DB | G1 | `approve.go:98` default `def.db` → `def.workspaceRoot` | G2 `default-db` | C `TestApproveListUsesTheNamedStores` | KILLED |
| MUT-C1 | M4c1 | C1: `!isLowerHex64(rc.Challenge) \|\|` → `false && !isLowerHex64(rc.Challenge) \|\|` | CS1 (all arms) | S `TestMintSessionForRequestNeedsTheChallenge` | KILLED |
| MUT-C1-EMPTY | M4c1 | C1: `!isLowerHex64(rc.Challenge)` → `rc.Challenge == ""` | CS1 `63-hex`, `upper`, `not-hex` | S `TestMintSessionForRequestNeedsTheChallenge` | KILLED |
| MUT-C6 | M4c1 | C6: `AND challenge = ?` → `AND (challenge = ? OR 1)` | CS5 `other-challenge` (mints) | S `TestConsumeBindsTheChallenge` | KILLED |
| MUT-C7 | M4c1 | C7: `case challenge != rc.Challenge:` → `case false && …` | CS5 `other-challenge` (generic error) | S `TestConsumeBindsTheChallenge` | KILLED |
| MUT-C7-ORDER | M4c1 | C7's case moved before the digest case | CS5 `other-digest-and-challenge` | S `TestConsumeBindsTheChallenge` | KILLED |
| MUT-A1 | M4c1 | A1: `consume.ApprovedBy != ""` → `false && …` (the M4e form keeps `\|\| consume.Approver != nil`) | AU1 `caller-tty`, `caller-ssh` | U `TestMintForRequestDerivesApprovedBy` | KILLED |
| MUT-A2 | M4c1 | A2: `err != nil \|\| len(digest) != 32` → `false && (…)` | AU1 `digest-not-hex`, `digest-31-bytes` | U `TestMintForRequestDerivesApprovedBy` | KILLED |
| MUT-A2-LEN | M4c1 | A2: `\|\| len(digest) != 32` → `\|\| false && len(digest) != 32` | AU1 `digest-31-bytes` | U `TestMintForRequestDerivesApprovedBy` | KILLED |
| MUT-A-CHALLENGE | M4c1 | `consume.Challenge = hex.EncodeToString(rm.Challenge[:])` → `consume.Challenge = consume.Digest` | AU1 `typed-records-tty` | U `TestMintForRequestDerivesApprovedBy` | KILLED |
| MUT-A5 | M4c1 | A5: `"tty"` → `"ttyx"` | AU1 `typed-records-tty` | U `TestMintForRequestDerivesApprovedBy` | KILLED |
| MUT-F4 | M4c1 | F4: `if err != nil \|\| len(challenge) != 32` → `if false && (…)` | CR11 (panic on the conversion) | C `TestApproveRefusesANonHexChallenge` | KILLED |
| MUT-TYPED-TTY | M4c1 | `approve.go`'s consume passes `ApprovedBy: "tty"` again | CL5 (A1 refuses the typed path) | C `TestApproveMintsExactlyTheRequest` | KILLED |
| MUT-APPROVE-CHALLENGE | M4c1 | `Challenge: [32]byte(challenge)` → `[32]byte{}` | CL5 (C7 refuses) | C `TestApproveMintsExactlyTheRequest` | KILLED |
| MUT-C2-NAME | M4d1 | C2: `a.Name == "" \|\|` → `false && a.Name == "" \|\|` | RC3 `ApproverUse.Name` | S `TestRequestConsumeCompleteness` | KILLED |
| MUT-C2-FORMAT | M4d1 | C2: `\|\| a.Format == "" \|\|` → `\|\| false && a.Format == "" \|\|` | RC3 `ApproverUse.Format` | S `TestRequestConsumeCompleteness` | KILLED |
| MUT-C2-KEY | M4d1 | C2: `\|\| len(a.PublicKey) == 0` → `\|\| false && …` | RC3 `ApproverUse.PublicKey` | S `TestRequestConsumeCompleteness` | KILLED |
| MUT-C3 | M4d1 | C3's `if` → `if false && …` | CS2 `no-approver-ssh` | S `TestApprovedByMatchesTheApprover` | KILLED |
| MUT-C4 | M4d1 | C4's `if` → `if false && …` | CS2 `approver-tty`, `approver-other-*` | S `TestApprovedByMatchesTheApprover` | KILLED |
| MUT-C4-NAMEONLY | M4d1 | C4 compares with `"ssh:"+rc.Approver.Name` | CS2 `approver-format-not-ssh` | S `TestApprovedByMatchesTheApprover` | KILLED |
| MUT-C5 | M4d1 | C5's `if` → `if false && …` | CS2 `v4-store` (no such table) | S `TestApprovedByMatchesTheApprover` | KILLED |
| MUT-RC3-FIELD | M4d1 | add `Label string` to `ApproverUse` | RC3 (no disposition) | S `TestRequestConsumeCompleteness` | KILLED |
| MUT-C8 | M4d2 | the C8 call `if rc.Approver != nil` → `if false && …` | CS3 `revoked-in-the-window` … (mints) | S `TestSignedConsumeChecksTheApproverInTheTransaction` | KILLED |
| MUT-C8-ALWAYS | M4d2 | the C8 call → `if rc.Approver != nil \|\| true` | CS4 (nil approver, panic) | S `TestTypedConsumeTouchesNoApprover` | KILLED |
| MUT-C8-ERR | M4d2 | `useApprover(…); err != nil` → `false && err != nil` | CS3 refusal arms | S `TestSignedConsumeChecksTheApproverInTheTransaction` | KILLED |
| MUT-C8-REVOKED | M4d2 | C8 SQL drops `AND revoked_at IS NULL` | CS3 `revoked-in-the-window` | S `TestSignedConsumeChecksTheApproverInTheTransaction` | KILLED |
| MUT-C8-FORMAT | M4d2 | C8 SQL `format = ?` → `(format = ? OR 1)` | CS3 `other-format` | S `TestSignedConsumeChecksTheApproverInTheTransaction` | KILLED |
| MUT-C8-KEY | M4d2 | C8 SQL `public_key = ?` → `(public_key = ? OR 1)` | CS3 `other-key` | S `TestSignedConsumeChecksTheApproverInTheTransaction` | KILLED |
| MUT-C8-ENROLLED | M4d2 | C8 SQL `enrolled_at <= ?` → `(enrolled_at <= ? OR 1)` | CS3 `enrolled-after-the-reading` | S `TestSignedConsumeChecksTheApproverInTheTransaction` | KILLED |
| MUT-C8-AT | M4d2 | C8 sets `last_used_at = at+1` | CS3 `active` (1010) | S `TestSignedConsumeChecksTheApproverInTheTransaction` | KILLED |
| MUT-C8-N | M4d2 | `if n != 1` → `if false && n != 1` | CS3 `unknown` (0 rows, mints) | S `TestSignedConsumeChecksTheApproverInTheTransaction` | KILLED |
| MUT-C9 | M4d2 | C9's case → `false && …` | CS3 `unknown` | S `TestSignedConsumeChecksTheApproverInTheTransaction` | KILLED |
| MUT-C10 | M4d2 | C10's case → `false && …` | CS3 `revoked-in-the-window` | S `TestSignedConsumeChecksTheApproverInTheTransaction` | KILLED |
| MUT-C11-FORMAT | M4d2 | C11 drops `row.Format != a.Format \|\|` | CS3 `other-format` | S `TestSignedConsumeChecksTheApproverInTheTransaction` | KILLED |
| MUT-C11-KEY | M4d2 | C11 drops `\|\| !bytes.Equal(…)` | CS3 `other-key` | S `TestSignedConsumeChecksTheApproverInTheTransaction` | KILLED |
| MUT-C12 | M4d2 | C12's case → `false && …` | CS3 `enrolled-after-the-reading` | S `TestSignedConsumeChecksTheApproverInTheTransaction` | KILLED |
| MUT-A1-APPROVER | M4e | A1′ drops `\|\| consume.Approver != nil` | AU2 `caller-approver` | U `TestMintForRequestRefusesAnUnboundApproval` | KILLED |
| MUT-A5-SIGNED | M4e | `if rm.Signed == nil` → `if true \|\| …` (a signed mint records tty) | AU1 signed | U `TestMintForRequestRecordsTheSigner` | KILLED |
| MUT-A4 | M4e | A4 `err != nil` → `false && err != nil` | AU2 `other-*`, `zero-approval` | U `TestMintForRequestRefusesAnUnboundApproval` | KILLED |
| MUT-A4-DIGEST | M4e | `Binds([32]byte(digest), …)` → `Binds(rm.Challenge, …)` | AU1 signed (BindingError digest) | U `TestMintForRequestRecordsTheSigner` | KILLED |
| MUT-A3-APPROVER | M4e | the bound approver's key loses its last byte | AU1 signed (C11 refuses) | U `TestMintForRequestRecordsTheSigner` | KILLED |
| MUT-A3-NOAPPROVER | M4e | `consume.Approver = …` deleted | AU1 signed (C4 refuses) | U `TestMintForRequestRecordsTheSigner` | KILLED |
| MUT-RC5-FIELD | M4e | add `Note string` to `RequestMint` | RC5 (no source) | U `TestRequestMintFieldsAreSourced` | KILLED |
| MUT-F2 | M5a | F2's `if` → `if false && …` | CR1 `list-and-*` (CR3 too) | C `TestApproveFlagRules` | KILLED |
| MUT-F2-SIG | M5a | F2 drops `*sigPath != ""` | CR1 `list-and-signature` | C `TestApproveFlagRules` | KILLED |
| MUT-F9 | M5a | F9 `readRegularFileCapped(…)` → `os.ReadFile(path)` (`-timeout 90s`) | CR4 `fifo` (hangs → timeout), `over-16k` | C `TestApproveBySignatureRefusalsBeforeAnyEffect` | KILLED (timeout) |
| MUT-F9-BOUND | M5a | `approveMaxSignatureBytes = 16 << 10` → `32 << 10` | CR4 `over-16k` | C `TestApproveBySignatureRefusalsBeforeAnyEffect` | KILLED |
| MUT-F10 | M5a | F10 `if err != nil` → `if false && err != nil` | CR4 `missing` | C `TestApproveBySignatureRefusalsBeforeAnyEffect` | KILLED |
| MUT-F11 | M5a | F11 → `if false && err != nil` | CR4 `garbage` (nil deref) | C `TestApproveBySignatureRefusalsBeforeAnyEffect` | KILLED |
| MUT-F12 | M5a | F12 → `if false && err != nil` | CR4 `v4-store` | C `TestApproveBySignatureRefusalsBeforeAnyEffect` | KILLED |
| MUT-F13 | M5a | F13 `if !ok` → `if false && !ok` | CR4 `unenrolled-key` | C `TestApproveBySignatureRefusalsBeforeAnyEffect` | KILLED |
| MUT-F14 | M5a | F14 → `if false && err != nil` | CR4 `revoked`, `other-request`, … (typed path opens a terminal) | C `TestApproveBySignatureRefusalsBeforeAnyEffect` | KILLED |
| MUT-SIG-REVOKEDAT | M5a | `RevokedAt: enrolled.RevokedAt` → `0*enrolled.RevokedAt` | CR4 `revoked` (reaches the network) | C `TestApproveBySignatureRefusalsBeforeAnyEffect` | KILLED |
| MUT-SIG-ENROLLEDAT | M5a | `EnrolledAt: enrolled.EnrolledAt` → `0 * …` | CR4 `enrolled-after-now` | C `TestApproveBySignatureRefusalsBeforeAnyEffect` | KILLED |
| MUT-SIG-EXPIRY | M5a | the law's `ExpiresAt: row.ExpiresAt` → `+ 1` | CR2 (bad signature) | C `TestApproveBySignatureMintsWithNoTerminal` | KILLED |
| MUT-SIG-CHALLENGE | M5a | the law's `Challenge` → the digest's bytes | CR2 | C `TestApproveBySignatureMintsWithNoTerminal` | KILLED |
| MUT-SIG-REFUSAL | M5a | `approveSignature(…); refusal != ""` → `false && …` | CR4 (proceeds) | C `TestApproveBySignatureRefusalsBeforeAnyEffect` | KILLED |
| MUT-SIG-BRANCH | M5a | `if *sigPath != ""` → `if false && …` (the flag ignored) | CR2 (typed path opens a terminal) | C `TestApproveBySignatureMintsWithNoTerminal` | KILLED |
| MUT-SIG-ORDER | M5a | the signature block moved after `checkDaemon` (after the network call) | CR4 (`daemon calls 1`) | C `TestApproveBySignatureRefusalsBeforeAnyEffect` | KILLED |
| MUT-CONSUME-SIGNED | M5a | the consume's `Signed: signed` → always nil | CR2 (terminal path) | C `TestApproveBySignatureMintsWithNoTerminal` | KILLED |
| MUT-CONSUME-EXPIRES | M5a | the consume's `ExpiresAt: row.ExpiresAt` → `row.CreatedAt` | CR2 (A4 BindingError expiry) | C `TestApproveBySignatureMintsWithNoTerminal` | KILLED |
| MUT-V1 | M5a | `signed := … && r.consume.Signed != nil && false` | CR2 (terminal opened) | C `TestApproveBySignatureMintsWithNoTerminal` | KILLED |
| MUT-V1-CONSUME | M5a | `signed := r.consume != nil` | CR3 (typed approve skips the terminal) | C `TestApproveTypedPathUnchangedWithAnApprover` | KILLED |
| MUT-V2-LINE | M5a | V2's line prints `""` for the approver | CR2 stdout | C `TestApproveBySignatureMintsWithNoTerminal` | KILLED |
| MUT-V2-PRESET | M5a | V2's `fmt.Fprint(stdout, r.presetPrompt)` deleted | CR2 stdout | C `TestApproveBySignatureMintsWithNoTerminal` | KILLED |
| MUT-C8-CLI | M5b | (store) the C8 call → `if false && …`, seen from the CLI | CR6 (mints) | C `TestApproveRevokeRaceIsRefusedAtTheWrite` | KILLED |
| MUT-F1 | M5c | F1's `if` → `if false && …` | CR1 `signature-and-write-message` | C `TestApproveFlagRules` | KILLED |
| MUT-F2-MSG | M5c | F2 drops `\|\| *msgPath != ""` | CR1 `list-and-write-message` | C `TestApproveFlagRules` | KILLED |
| MUT-MSG-BRANCH | M5c | `if *msgPath != ""` → `if false && …` | CR10 `writes-then-signs` | C `TestApproveWriteMessage` | KILLED |
| MUT-F5-EXCL | M5c | `O_EXCL` → `O_TRUNC` | CR10 `existing-file` (overwrites) | C `TestApproveWriteMessage` | KILLED |
| MUT-F5 | M5c | F5 → `if false && err != nil` | CR10 `existing-file` (nil file) | C `TestApproveWriteMessage` | KILLED |
| MUT-F6-RETURN | M5c | the `return exitOK` after writing deleted | CR10 `writes-then-signs` (opens a terminal) | C `TestApproveWriteMessage` | KILLED |
| MUT-F6-REFUSAL | M5c | `writeApprovalMessage(…); refusal != ""` → `false && …` | CR10 `existing-file` (exitOK) | C `TestApproveWriteMessage` | KILLED |
| MUT-F6-ORDER | M5c | the `--write-message` block moved before admission's refusal | CR10 `used-request` (a file appears) | C `TestApproveWriteMessage` | KILLED |
| MUT-MSG-EXPIRY | M5c | the message's expiry → `row.CreatedAt` | CR10 bytes | C `TestApproveWriteMessage` | KILLED |
| MUT-MSG-NEXT | M5c | the next-step line names `mr.WorldDB` twice | CR10 stdout | C `TestApproveWriteMessage` | KILLED |

**Known-equivalent (named, not counted; each run in scratch and SURVIVED as predicted):**
- **MUT-RET-DIGEST** (M4a): `digest: in.Admit.BytesDigest`. Equivalent because P1 has just proven it equal to `StoredDigest`. — scratch: SURVIVED
- **MUT-C8-AFTER-INSERT** (M4d2): the C8 call moved after `insertSessionRow` (still before `Commit`). Equivalent because the same transaction rolls both back. — scratch: SURVIVED
- **MUT-A3-APPROVEDBY** (M4e): `consume.ApprovedBy = "ssh:" + name`. Equivalent because `Binds` passed, so the approval came from P3, which fixes `format = ssh`. — scratch: SURVIVED
- **MUT-SIG-CONSUMED** (M5a): the law's `Consumed: row.Consumed` → `false`. Equivalent because `approveAdmission` already refused a consumed row. — scratch: SURVIVED
- **MUT-SIG-BYTES** (M5a): the law's `BytesDigest` → the stored digest. Equivalent because `approveAdmission` already proved the bytes match it. — scratch: SURVIVED
- **MUT-SIG-KEYSRC** (M5a): the law's key → the signature's own key. Equivalent because `ApproverByKey` matched that exact blob (S3-D7: defence in depth; C8's `public_key = ?` is the backstop). — scratch: SURVIVED
- **MUT-N10-ERR / MUT-N10-NL** (sprint 2): (carried, sprint-2 P3). Equivalent because `ReadString` returns an error iff the line lacks `\n`. — not re-run (sprint 2 measured)
- **MUT-MSG-MODE** (M5c): `0o644` → `0o600`. Equivalent because not asserted: the effective mode depends on the umask (S3-D9). — scratch: not run (umask)

**Propagation-only branches (named, not counted):** `approveSignature`'s `OpenCredentialsReadOnly` error
(the store vanished after `ro.MintRequest` succeeded); `writeApprovalMessage`'s `SignedMessage` error
(negative expiry, unreachable after admission) and its `f.Write`/`f.Close` error (removes the file);
`useApprover`'s `ExecContext`/`RowsAffected` errors; `classifyApproverUse`'s `err != nil` and `default`
arms; `requestConsume`'s `Use()` `ok` (shown true by `Binds`); the message-only refusal texts.

---

## 5. Files (final `git status --short` after M5c; after M4e if the gate stops the sprint)

| File | Sub-milestone | Change |
|---|---|---|
| `host/approval/signed.go`, `signed_test.go` | M4a, M4b | new |
| `cmd/ailang-worldd/approver_stores_test.go` | G1 | new |
| `host/store/mint_requests.go` | M4c1, M4d1, M4d2 | C1, C6, C7; `ApproverUse`, C2–C5; the C8 call |
| `host/authority/mint.go` | M4c1, M4e | `requestConsume` A1, A2, A5; `ExpiresAt`, `Signed`, A1′, A3, A4 |
| `cmd/ailang-worldd/approve.go` | M4c1, M5a, M5c | F4 + the consume; `--signature`, F2; `--write-message`, F1, F2′ |
| `cmd/ailang-worldd/mint_request.go` | M4c1, M5a | the `consume` field-map text |
| `host/store/mint_requests_test.go`, `host/authority/mint_request_test.go` | M4c1 | the consume callers |
| `host/store/consume_binding_test.go` | M4c2, M4d1, M4d2 | new: CS1, CS5; CS2, RC3; CS3, CS4 |
| `host/authority/mint_derive_test.go` | M4c2, M4e | new: AU1; AU1 signed, AU2, RC5 |
| `cmd/ailang-worldd/approve_sig_test.go` | M4c2, M5a, M5b, M5c | new: CR11; CR2; CR3, CR4/5, CR6; CR10, CR1 |
| `host/store/approvers.go` | M4d2 | `useApprover`, `classifyApproverUse` |
| `cmd/ailang-worldd/approve_sig.go` | M5a, M5c | new: `approveSignature`; `writeApprovalMessage` |
| `cmd/ailang-worldd/session_new.go` | M5a | V1, V2, the step comment |
| `cmd/ailang-worldd/main.go`, `help_test.go` | M5a, M5c | usage lines and wants |
| `cmd/ailang-worldd/approver.go` | M5a | `approverHelp` wording only |
| `docs/QUICKSTART.md` | M5c | §15 (attended-pending) |
| this plan and `sprint_w-world-approvals-by-key-a1b-s3.json` | planner | |

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
go test ./host/store/ -run '^TestProductionContextRoots$' -count=1         # UNEDITED, no new root
go test ./host/daemon/ -run '^TestDaemonDependencyAllowlist$' -count=1
go test -race ./host/approval/... ./host/store/ ./host/authority/ -count=1
go test -race ./cmd/ailang-worldd/ -run '^(TestApproveBySignature|TestApproveRevokeRace|TestApproveWriteMessage|TestApproveFlagRules|TestApproveTypedPathUnchanged|TestApproveRefusesANonHexChallenge|TestApproverVerbsUseTheNamedStores|TestApproveListUsesTheNamedStores)' -count=3
# optional locally (CI's ci.yml:141 step is the verdict; the existing test uses temp keys only, unchanged):
# WORLD_REQUIRE_SSH_KEYGEN=1 go test ./host/approval/ -run '^TestSSHSIGCrossCheckWithSSHKeygen$' -count=1 -v
go test ./cmd/ailang-worldd/ -count=1 -v -run '^(TestApprove|TestApprover|TestEpisodeStartRequest|TestSessionMint_|TestSessionNew|TestEveryVerbHasHelp|TestHelpIsNeverBlockedByTheRegistryKeyGuard)' 2>&1 | tee /tmp/w377-cl.log | tail -3
if grep -E '^\s*--- (SKIP|FAIL)' /tmp/w377-cl.log; then exit 1; fi
# AFTER the controller commits (the tracked-binary gate reads committed blobs):
VERIFY_GO_LEGS=plain ./scripts/verify_go.sh
git status --short                                                         # only the §5 files
```

If the gate stopped the sprint after M4e, drop the `TestApproveBySignature|TestApproveRevokeRace|TestApproveWriteMessage|TestApproveFlagRules|TestApproveTypedPathUnchanged`
names from the race line. **Known base reds outside this sprint** (§1.1) are the base; any other red, any
race report, or any red in a touched package is a regression. CI's `go-race-verify` is the race verdict.

**Must stay green, unedited**: every `TestApprove*` and `TestEpisodeStartRequest*` of `approve_test.go`
(the typed path, CL1–CL19, CT1 — including `TestApproveMintsExactlyTheRequest`, which kills
MUT-TYPED-TTY), `TestMintRequestFieldMapIsComplete`, `TestExistingCallersReadUnbounded`,
`TestReadTermLineWithin`, every `TestApprover*` of sprint 2, `TestSSHSIG*`, `TestCredentials*`,
`TestMintSessionForRequest*`, `TestMintForRequestChecksTheRowInsideTheTransaction`,
`TestOlderCredentialsStoreErrorText`, `TestProductionContextRoots`, `TestDaemonDependencyAllowlist`,
`TestApprovalTestIsTestOnly`, `TestApprovalImportsOnlyTheStandardLibrary`, the QUICKSTART tests. **The only
edited existing tests**: `okConsume` and `TestConcurrentApprovalsMintOnce` (+`Challenge`),
`TestMintForRequestChecksTheRowInsideTheTransaction` (`consumeAt` −`ApprovedBy`, +`Challenge` ×3),
`TestEveryVerbHasHelp` (+2 wants, gated part).

---

## 7. Executor directive notes

- **No git writes.** Never `git add`, `commit`, `stash`, `checkout`, `reset`, `restore`, `rebase`, `push`,
  `worktree` or `apply`. `patch -p1` (§2) is the only writer; read-only `git status`/`git diff` are fine.
- **Runs.** Run 1 = M4a, M4b; Run 2 = G1, M4c1; Run 3 = M4c2, M4d1; Run 4 = M4d2, M4e; **gate**; Run 5 =
  M5a, M5b; Run 6 = M5c. Each run reads §0, §2 and its own §3 sections, then applies its patches in order.
- **Expected numstat after each committed sub-milestone** (controller, after its commit):
  ```
  bash -c 'cd /Users/voightkampff/dev/sunholo-data/.wt-world-iter377-row201-a1b-s3 && git diff --numstat 545fdaaecf6f23dd58072675964754ad716049f9 HEAD -- cmd host | awk "{s+=\$1+\$2} END {print s}"'
  ```
  M4a **221**, M4b **323**, G1 **511**, M4c1 **602**, M4c2 **815**, M4d1 **993**, M4d2 **1125**, M4e
  **1314**, M5a **1577**, M5b **1752**, M5c **1910**. **Any other value means a patch was not applied as
  written: stop and report.**
- **The size gate (controller), after the M4e commit:** the numstat above is **≤ 1,404** → run Runs 5 and 6
  (all three of M5a, M5b, M5c). Over → do not run them; M5a–M5c move to Appendix C unchanged (their patches
  in Appendix P are final against the M4e tree) and the sprint closes after M4e. Expected: 1,314 (passes
  with 90 lines of margin).
- **Snapshot after every sub-milestone, cumulatively** (as sprint 2 §7, with `world-iter377`):
  ```
  bash -c 'set -euo pipefail; M=M4a; W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter377-row201-a1b-s3
    D=$HOME/.ailang/state/world-iter377/snap/$M; rm -rf "$D"; mkdir -p "$D"; cd "$W"
    git status --porcelain --untracked-files=all | while read -r st f; do
      case "$f" in design_docs/planned/w-world-approvals-by-key-a1b-sprint3-plan.md|design_docs/planned/sprint_w-world-approvals-by-key-a1b-s3.json) continue;; esac
      mkdir -p "$D/$(dirname "$f")"; cp -p "$f" "$D/$f"; done
    (cd "$D" && find . -type f ! -name MANIFEST.sha256 | sort | xargs shasum -a 256) > "$D/MANIFEST.sha256"
    cat "$D/MANIFEST.sha256"'
  ```
  (`M` = `M4a` … `M5c`.) No file is deleted, so the porcelain status has only `M` and `??`; stop on any
  other status. The controller rebuilds one commit per sub-milestone (iteration 376's `commit_ms.sh` with
  `WT` and the snapshot root pointed here). **M4c1 is one commit, never split** (§0.3).
- **Stop conditions.** A `patch --dry-run` failure; a numstat other than the expected value; a red
  earlier-milestone test; a `--- SKIP` in a `-v` run of `cmd/ailang-worldd`; any need to touch a file
  outside §5; any `ssh-keygen` invocation by the executor; any reference to `~/.ssh`, `~/.ailang/world` or a real
  `/dev/tty`. Stop and report; do not work around. **Do not run the QUICKSTART §15 commands** (attended).
- **Shell.** Wrap loops, globs and `$VAR` expansions in `bash -c` or a `bash <<'EOF'` heredoc (zsh does
  not word-split). Export `AILANG_BIN` in the same shell as the Go tests. Poll any background command's
  artifact in a bounded `date +%s` loop in the same turn; never end a turn announcing a wait.

---

## 8. Risks

- **R1 — Size.** Measured 1,910 against a 2,000 cap; the committed part 1,314 (×1.2 = 1,577). The gate and
  the per-step expected numstat make any drift visible before the tail runs.
- **R2 — `approve --signature` ships without `--sign-with` and `--list --decided`.** An operator signs out of
  band with `ssh-keygen -Y sign` (the `--write-message` output and QUICKSTART §15 show the line).
  `approve --list` shows pending requests only; the decided list is Appendix C. Until then
  `approved_by=ssh:<name>` is visible in the store and in admission's replay refusal (`approve.go:70`,
  "was already used at … by ssh:<name>"). Revoking an approver does not yet list what it approved
  (K8/K9, Appendix C); Q4's default (list, never auto-revoke) is unaffected.
- **R3 — R6 (SSHSIG proves possession, not presence)** is unchanged and stated in `approveHelp` and the
  enrolment y/N (AP5). Q5's attended exit run is Mark's.
- **R4 — Timing.** CR4's `fifo` arm returns at once on the production code (`O_NONBLOCK`); on MUT-F9 it hangs
  until the test timeout (use `-timeout 90s`). CR6 depends on the consume's clock being the third `env.now`
  call (C18) and asserts it, so a reordering fails loudly rather than passing vacuously.
- **R5 — `patch` portability.** The planner used `/usr/bin/patch` (`patch 2.0-12u11-Apple`). GNU patch
  applies the same unified diffs; `--dry-run` guards both.
- **R6 — Haiku and patches.** The executor transcribes nothing; the risks are a skipped patch (dry-run and
  numstat catch it) and a hand "fix" after a red test (forbidden by §7's stop conditions).

---

## 9. Acceptance criteria (the judge scores against these)

1. M4a/M4b: `host/approval/signed.go` and `signed_test.go` are sprint-2 §2.5/§2.8 byte for byte (the M4a
   cut is `gofmt`-formatted); SA1–SA3 and RC1 pass; `host/approval` imports only the standard library.
2. G1: every approver subcommand and `approve --list` honour `--credentials-db` (a store in another
   directory; nothing created beside `--db`) and the default `--db` (`HOME` temp); MUT-F-CREDS,
   MUT-AS-CREDFOR, MUT-F-DB and MUT-A-DB, which survive the pristine suite (C13), are killed.
3. M4c1 landed as **one** commit; after it, every unedited typed-path test is green, `approve` passes no
   `approved_by`, and `approved_by = tty` is derived by authority.
4. The consume refuses a missing or malformed challenge (C1), binds the row's challenge (C6/C7), refuses an
   inconsistent `approved_by` (C3/C4), an approver on schema 4 (C5), and, inside its own transaction,
   an approver that is unknown, revoked (including by a commit inside the window), re-keyed or ahead of
   the clock (C8–C12), leaving the request pending, no session and `last_used_at` unchanged.
5. Authority refuses a caller-set `approved_by` or approver (A1′), a non-hex digest (A2), an approval for
   another digest, challenge or expiry, and an unverified one (A4); records `ssh:<name>` and marks the
   approver used for a bound one (A3); `tty` otherwise (A5). RC3 and RC5 classify and perturb every field.
6. If the gate ran M5: `approve --signature` mints with **no terminal** only for an enrolled, unrevoked
   key's SSHSIG over the stored request's exact message in `world-approval@sunholo`; every refusal (CR4's
   12 arms) is `exitUsage`, before any terminal and any daemon call; a revoke committed inside the mint's
   transaction is `exitFatal` with the worktree removed (CR6); the typed path is unchanged with an approver
   enrolled (CR3); `--write-message` writes exactly the 90-byte message with `O_EXCL` and nothing else
   (CR10); F1/F2 hold (CR1); the help, usage and QUICKSTART §15 describe the paths.
7. Each §4 mutant, applied alone, makes its named test fail; the named equivalents survive for the stated
   reason.
8. Every §6 gate passes, including `VERIFY_GO_LEGS=plain ./scripts/verify_go.sh` after the commit;
   `TestProductionContextRoots` is unedited and green; no test skips; `git status --short` lists only §5.
9. The numstat after each sub-milestone equals §7's expected value.
10. No real key, real store, `~/.ssh` access, new `ssh-keygen` invocation, mint outside a temp store, or
    publish anywhere in the sprint (D-WORLD-109); the executor runs no `ssh-keygen` at all.

---

## Appendix C — carried to sprint 4 (and 5): the decided list, K8/K9, `--sign-with`, M6

Not executed in this sprint. If the §7 gate stopped Run 5, **M5a, M5b and M5c (Appendix P) head this
appendix unchanged**, with the expected numstat values 1,577 / 1,752 / 1,910 on top of the M4e tree.

**Measured premises at `545fdaae`** (re-verify; this sprint changes the files marked †):
- `host/store/mint_requests.go`† `PendingMintRequests` (`:165`, `ORDER BY created_at DESC, request_id DESC`,
  cap `MaxPendingMintRequests`) is the template for `DecidedMintRequests` (`consumed_at IS NOT NULL`,
  `ORDER BY consumed_at DESC, request_id DESC`, cap 64, works at v4; design §2.4 [r2]). No method lists
  consumed rows today (`grep -n "consumed_at IS NOT NULL" host/store/mint_requests.go` → `:83`, a field
  comment, and `:296`, `classifyConsume`'s scan; no query selects consumed rows).
- `cmd/ailang-worldd/approve.go`† `--list` (`:136-151` at HEAD) prints pending rows; `--decided` (F3: without
  `--list` → help) adds the decided list (F15 `no decided mint requests`, F16 one line per row: id, episode,
  `approved_by`, consumed RFC 3339, the first 16 hex of the credential id).
- `cmd/ailang-worldd/approver.go` `runApproverRevoke` (`:377`); K8 runs after the write root ends, under its
  own 3 s read-only root → `TestProductionContextRoots`'s pin `approver.go|runApproverRevoke|Background: 2 → 3`
  (B-D4); K9 `no session was approved by it`; Q4: it revokes no session.
- `cmd/ailang-worldd/session.go:25-38` `sessionEnv` (four seams; `sign` is the fifth, design §2.9);
  `realSessionEnv` `:46`. `--sign-with`'s F7 (`env.sign == nil` → `exitFatal`), F7a (`rem <= 0` → the
  expired text before calling the signer), F8 (signer error → `exitUsage`), and the `signForApproval`
  root (`approve.go|signForApproval|Background: 1`).
- `.github/workflows/ci.yml:141-155` the row 201 SSHSIG step (`sed -n 141,155p` → `- name: Row 201 SSHSIG cross-check with real ssh-keygen, verbose` … the PASS-grep loop over the parent and two subtests); M6 widens it.

**Decisions to carry:**
- **C-D1**: `DecidedMintRequests`, `--list --decided` and K8/K9 land together (B-D5), as sprint 4.
- **C-D2**: `--sign-with` lands whole with M6: the seam, `signWithSSHKeygen` (exec `ssh-keygen -Y sign -n
  world-approval@sunholo -f <key>` with the message on stdin, 16 KiB stdout cap), F7/F7a/F8, G1–G3, CR8/CR9
  with a fake signer, XS1/XS2 with the real tool under `WORLD_REQUIRE_SSH_KEYGEN=1`, the CI step, and
  QUICKSTART §15's `--sign-with`/`--list --decided` lines. Shipping the production exec before its
  real-tool tests is the thing to avoid; if sprint 4 has room after C-D1, the planner may pull M6 in whole.
- **C-D3**: the patch protocol (§2) and per-step expected numstat (§7) carry.
- **C-D4**: `approveEffectProblems` (`approve_test.go:1057`) stays typed-only; signed-path tests filter its
  `approved_by:` line as CR2 does.

**Estimate** (design × the measured ratio of this sprint's exact code to the design's numbers): sprint 4
(C-D1) ≈350 ×1.2 ≈420; M6 ≈510 (+20 YAML, +25 docs) ×1.2 ≈610. Together ≈1,030 — one sprint is possible.

**Mutants carried** (design §6 rows not in §4): F3, F15, F16, CS6 (filter, `request_id DESC` tiebreak), CR7,
K8 (filter by the wrong `approved_by`), K9, CR12; F7, F7a (`rem <= 0` vs `< 0`), F8, CR8 (the seam's output
trusted without verification), CR9 (the bound taken from `RequestLifetime`), G1–G3, XS1, XS2, the
`--sign-with` exclusivity arms of F1/F2.

---

## Appendix P — the patches (apply with §2; generated from the compiled scratch stages; do not edit)

### Patch M4a — the signature law

````diff M4a
--- a/host/approval/signed.go
+++ b/host/approval/signed.go
@@ -0,0 +1,94 @@
+package approval
+
+// Row 201 A1b (design_docs/planned/w-world-approvals-by-key-a1b.md §2.3): the pure signature
+// law. A mint request is approved by signature only when it is admitted, the approver is
+// enrolled, active and not ahead of the clock, and its key's SSHSIG over the request's exact
+// message (stored digest, row challenge, row expiry) verifies. No I/O here (coding standard S2).
+
+import (
+	"bytes"
+	"errors"
+	"fmt"
+)
+
+// SignedApproval is a verified signature approval. Its fields are unexported, so only
+// ApproveBySignature, in this package, can return one with ok set; the zero value binds nothing.
+type SignedApproval struct {
+	ok        bool
+	digest    [32]byte
+	challenge [32]byte
+	expiry    int64
+	name      string
+	format    string
+	key       []byte
+}
+
+// SignatureInput is what ApproveBySignature judges.
+type SignatureInput struct {
+	Admit     AdmitInput // the existing admission input (row times, consumed, both digests)
+	Challenge [32]byte   // the request row's challenge
+	Signature []byte     // the armored SSHSIG
+	Approver  Approver   // the enrolled row the signature's key names
+}
+
+// NotAdmittedError is admission's refusal on the signature path (P1): the request is consumed,
+// expired, ahead of the clock, or its bytes do not match their digest.
+type NotAdmittedError struct{ Verdict Verdict }
+
+func (e *NotAdmittedError) Error() string {
+	return fmt.Sprintf("approval: the request is not admitted (%s)", e.Verdict)
+}
+
+// ErrNotVerified is any use of a SignedApproval that ApproveBySignature did not return.
+var ErrNotVerified = errors.New("approval: not a verified signature approval")
+
+// BindingError is a verified approval used for a request it was not made for.
+type BindingError struct{ Field string } // "digest" | "challenge" | "expiry"
+
+func (e *BindingError) Error() string {
+	return fmt.Sprintf("approval: the signature approval was made for another request (%s differs)", e.Field)
+}
+
+// ApproveBySignature is the signature law (P1-P6, in this order): admission first, so a replay
+// says "already used"; the approver's state before the cryptography, so a revoked key says
+// "revoked" even when its signature is good.
+func ApproveBySignature(in SignatureInput) (*SignedApproval, error) {
+	if v := Admit(in.Admit); v != Admitted { // P1
+		return nil, &NotAdmittedError{Verdict: v}
+	}
+	if !ApproverNamePattern.MatchString(in.Approver.Name) { // P2
+		return nil, ErrApproverName
+	}
+	if in.Approver.Format != FormatSSH { // P3
+		return nil, ErrApproverFormat
+	}
+	if in.Approver.RevokedAt != 0 { // P4
+		return nil, ErrApproverRevoked
+	}
+	if in.Admit.Now < in.Approver.EnrolledAt { // P5
+		return nil, ErrApproverClockBehind
+	}
+	if err := VerifyApproval(in.Signature, in.Admit.StoredDigest, in.Challenge, in.Admit.ExpiresAt, in.Approver.Key); err != nil { // P6
+		return nil, err
+	}
+	return &SignedApproval{ok: true, digest: in.Admit.StoredDigest, challenge: in.Challenge,
+		expiry: in.Admit.ExpiresAt, name: in.Approver.Name, format: in.Approver.Format,
+		key: bytes.Clone(in.Approver.Key.Blob)}, nil
+}
+
+// Binds reports whether s is a verified approval of exactly this digest, challenge and expiry.
+func (s *SignedApproval) Binds(digest, challenge [32]byte, expiry int64) error {
+	if s == nil || !s.ok { // B1
+		return ErrNotVerified
+	}
+	if s.digest != digest { // B2
+		return &BindingError{Field: "digest"}
+	}
+	if s.challenge != challenge { // B3
+		return &BindingError{Field: "challenge"}
+	}
+	if s.expiry != expiry { // B4
+		return &BindingError{Field: "expiry"}
+	}
+	return nil
+}
--- a/host/approval/signed_test.go
+++ b/host/approval/signed_test.go
@@ -0,0 +1,127 @@
+package approval
+
+import (
+	"errors"
+
+	"testing"
+)
+
+// ---- M4a ----
+
+// sigInput is a valid SignatureInput over the checked-in fixtures: the ed25519 approver
+// "mark-mac", enrolled at fixtureExpiry-1000, admitted at fixtureExpiry-300 for a request made at
+// fixtureExpiry-600 that expires at fixtureExpiry, signed by testdata/sshsig/ed25519.sig.
+func sigInput(t *testing.T) SignatureInput {
+	t.Helper()
+	d := fixtureDigest(t)
+	return SignatureInput{
+		Admit:     AdmitInput{Now: fixtureExpiry - 300, CreatedAt: fixtureExpiry - 600, ExpiresAt: fixtureExpiry, StoredDigest: d, BytesDigest: d},
+		Challenge: fixtureChallenge(),
+		Signature: fixtureBytes(t, "ed25519.sig"),
+		Approver:  Approver{Name: "mark-mac", Format: FormatSSH, Key: fixtureKey(t, "ed25519.pub"), EnrolledAt: fixtureExpiry - 1000},
+	}
+}
+
+// verifiedApproval is ApproveBySignature(sigInput(t)), which must succeed.
+func verifiedApproval(t *testing.T) *SignedApproval {
+	t.Helper()
+	s, err := ApproveBySignature(sigInput(t))
+	if err != nil || s == nil {
+		t.Fatalf("the clean input: (%v, %v), want a verified approval", s, err)
+	}
+	return s
+}
+
+// SA1: one arm per branch of ApproveBySignature, in its order (P1-P6).
+func TestApproveBySignatureOrder(t *testing.T) {
+	other := fixtureDigest(t)
+	other[0] ^= 1
+	notAdmitted := func(v Verdict) func(error) bool {
+		return func(err error) bool {
+			var na *NotAdmittedError
+			return errors.As(err, &na) && na.Verdict == v
+		}
+	}
+	is := func(target error) func(error) bool { return func(err error) bool { return errors.Is(err, target) } }
+	arms := []struct {
+		name   string
+		mutate func(*SignatureInput)
+		want   func(error) bool
+	}{
+		{"consumed", func(in *SignatureInput) { in.Admit.Consumed = true }, notAdmitted(RefuseConsumed)},
+		{"consumed-and-revoked", func(in *SignatureInput) { in.Admit.Consumed = true; in.Approver.RevokedAt = fixtureExpiry - 900 }, notAdmitted(RefuseConsumed)},
+		{"clock-behind-request", func(in *SignatureInput) { in.Admit.Now = fixtureExpiry - 601 }, notAdmitted(RefuseClockBehind)},
+		{"expired", func(in *SignatureInput) { in.Admit.Now = fixtureExpiry }, notAdmitted(RefuseExpired)},
+		{"bytes-changed", func(in *SignatureInput) { in.Admit.BytesDigest = other }, notAdmitted(RefuseDigestMismatch)},
+		{"name-upper", func(in *SignatureInput) { in.Approver.Name = "Mark" }, is(ErrApproverName)},
+		{"name-empty", func(in *SignatureInput) { in.Approver.Name = "" }, is(ErrApproverName)},
+		{"format-webauthn", func(in *SignatureInput) { in.Approver.Format = FormatWebAuthn }, is(ErrApproverFormat)},
+		{"revoked-good-signature", func(in *SignatureInput) { in.Approver.RevokedAt = fixtureExpiry - 900 }, is(ErrApproverRevoked)},
+		{"revoked-bad-signature", func(in *SignatureInput) {
+			in.Approver.RevokedAt = fixtureExpiry - 900
+			in.Signature = withLastSigByteFlipped(t, in.Signature)
+		}, is(ErrApproverRevoked)},
+		{"enrolled-after-now", func(in *SignatureInput) { in.Approver.EnrolledAt = in.Admit.Now + 1 }, is(ErrApproverClockBehind)},
+		{"git-namespace", func(in *SignatureInput) { in.Signature = fixtureBytes(t, "git-ns-ed25519.sig") }, is(ErrNamespace)},
+		{"p256-signature", func(in *SignatureInput) { in.Signature = fixtureBytes(t, "p256.sig") }, is(ErrKeyMismatch)},
+		{"flipped-signature", func(in *SignatureInput) { in.Signature = withLastSigByteFlipped(t, in.Signature) }, is(ErrBadSignature)},
+		{"not-armored", func(in *SignatureInput) { in.Signature = []byte("x") }, is(ErrArmor)},
+		{"other-digest", func(in *SignatureInput) { in.Admit.StoredDigest, in.Admit.BytesDigest = other, other }, is(ErrBadSignature)},
+		{"other-challenge", func(in *SignatureInput) { in.Challenge[0] ^= 1 }, is(ErrBadSignature)},
+		{"other-expiry", func(in *SignatureInput) { in.Admit.ExpiresAt = fixtureExpiry + 1 }, is(ErrBadSignature)},
+	}
+	for _, a := range arms {
+		t.Run(a.name, func(t *testing.T) {
+			in := sigInput(t)
+			a.mutate(&in)
+			s, err := ApproveBySignature(in)
+			if s != nil || !a.want(err) {
+				t.Fatalf("(%v, %v): want a nil approval and the %s refusal", s, err, a.name)
+			}
+		})
+	}
+	t.Run("enrolled-at-now", func(t *testing.T) {
+		in := sigInput(t)
+		in.Approver.EnrolledAt = in.Admit.Now
+		if s, err := ApproveBySignature(in); err != nil || s == nil {
+			t.Fatalf("(%v, %v): an approver enrolled at now is active", s, err)
+		}
+	})
+	t.Run("clean", func(t *testing.T) {
+		if err := verifiedApproval(t).Binds(fixtureDigest(t), fixtureChallenge(), fixtureExpiry); err != nil {
+			t.Fatalf("the clean approval does not bind its own request: %v", err)
+		}
+	})
+}
+
+// SA2: Binds compares each of digest, challenge and expiry, and refuses an unverified value.
+func TestSignedApprovalBinds(t *testing.T) {
+	s := verifiedApproval(t)
+	d, c, e := fixtureDigest(t), fixtureChallenge(), int64(fixtureExpiry)
+	if err := s.Binds(d, c, e); err != nil {
+		t.Fatalf("Binds(own request) = %v", err)
+	}
+	var nilApproval *SignedApproval
+	for name, u := range map[string]*SignedApproval{"zero": {}, "nil": nilApproval} {
+		if err := u.Binds(d, c, e); !errors.Is(err, ErrNotVerified) {
+			t.Errorf("%s.Binds = %v, want ErrNotVerified", name, err)
+		}
+	}
+	d2, c2 := d, c
+	d2[31] ^= 1
+	c2[31] ^= 1
+	for _, a := range []struct {
+		field string
+		err   error
+	}{
+		{"digest", s.Binds(d2, c, e)},
+		{"challenge", s.Binds(d, c2, e)},
+		{"expiry", s.Binds(d, c, e+1)},
+		{"expiry", s.Binds(d, c, e-1)},
+	} {
+		var be *BindingError
+		if !errors.As(a.err, &be) || be.Field != a.field {
+			t.Errorf("Binds with %s changed = %v, want a *BindingError naming %s", a.field, a.err, a.field)
+		}
+	}
+}
````

### Patch M4b — the accessors and RC1

````diff M4b
--- a/host/approval/signed.go
+++ b/host/approval/signed.go
@@ -92,3 +92,21 @@
 	}
 	return nil
 }
+
+// ---- M4b ----
+
+// ApprovedBy is the approved_by value, format ":" name ("ssh:mark-mac"); "" unless verified.
+func (s *SignedApproval) ApprovedBy() string {
+	if s == nil || !s.ok { // B5
+		return ""
+	}
+	return s.format + ":" + s.name
+}
+
+// Use is what the consume binds in SQL: the approver's name, format and key (a copy).
+func (s *SignedApproval) Use() (name, format string, key []byte, ok bool) {
+	if s == nil || !s.ok { // B6
+		return "", "", nil, false
+	}
+	return s.name, s.format, bytes.Clone(s.key), true
+}
--- a/host/approval/signed_test.go
+++ b/host/approval/signed_test.go
@@ -1,8 +1,10 @@
 package approval
 
 import (
+	"bytes"
 	"errors"
-
+	"fmt"
+	"reflect"
 	"testing"
 )
 
@@ -122,6 +124,88 @@
 		var be *BindingError
 		if !errors.As(a.err, &be) || be.Field != a.field {
 			t.Errorf("Binds with %s changed = %v, want a *BindingError naming %s", a.field, a.err, a.field)
+		}
+	}
+}
+
+// ---- M4b ----
+
+// SA3: the accessors fail closed on an unverified value, and Use hands out copies of the key.
+func TestSignedApprovalAccessorsFailClosed(t *testing.T) {
+	var nilApproval *SignedApproval
+	for name, u := range map[string]*SignedApproval{"zero": {}, "nil": nilApproval} {
+		if got := u.ApprovedBy(); got != "" {
+			t.Errorf("%s.ApprovedBy() = %q, want \"\"", name, got)
+		}
+		if n, f, k, ok := u.Use(); ok || n != "" || f != "" || k != nil {
+			t.Errorf("%s.Use() = (%q, %q, %x, %v), want zero values and false", name, n, f, k, ok)
+		}
+	}
+	in := sigInput(t)
+	s, err := ApproveBySignature(in)
+	if err != nil {
+		t.Fatal(err)
+	}
+	if got := s.ApprovedBy(); got != "ssh:mark-mac" {
+		t.Errorf("ApprovedBy() = %q, want ssh:mark-mac", got)
+	}
+	want := bytes.Clone(in.Approver.Key.Blob)
+	n, f, k, ok := s.Use()
+	if !ok || n != "mark-mac" || f != FormatSSH || !bytes.Equal(k, want) {
+		t.Fatalf("Use() = (%q, %q, %x, %v), want (mark-mac, ssh, the ed25519 blob, true)", n, f, k, ok)
+	}
+	k[0] ^= 0xff                    // the caller's copy
+	in.Approver.Key.Blob[1] ^= 0xff // the input ApproveBySignature read
+	if _, _, k2, _ := s.Use(); !bytes.Equal(k2, want) {
+		t.Fatalf("Use() after editing a returned key and the input key = %x, want %x (copies)", k2, want)
+	}
+}
+
+// signedApprovalFields is the disposition of every SignedApproval field: what reads it.
+var signedApprovalFields = map[string]string{
+	"ok":        "B1/B5/B6: an unverified value binds nothing and has no approver",
+	"digest":    "B2: Binds",
+	"challenge": "B3: Binds",
+	"expiry":    "B4: Binds",
+	"name":      "ApprovedBy and Use",
+	"format":    "ApprovedBy and Use",
+	"key":       "Use",
+}
+
+// RC1: every SignedApproval field has a disposition, and changing any one of them in a verified
+// value changes what Binds, ApprovedBy or Use report.
+func TestSignedApprovalFieldsAreAllBound(t *testing.T) {
+	typ := reflect.TypeOf(SignedApproval{})
+	for i := 0; i < typ.NumField(); i++ {
+		if _, ok := signedApprovalFields[typ.Field(i).Name]; !ok {
+			t.Errorf("SignedApproval.%s has no disposition in signedApprovalFields", typ.Field(i).Name)
+		}
+	}
+	if typ.NumField() != len(signedApprovalFields) {
+		t.Fatalf("SignedApproval has %d fields, signedApprovalFields %d", typ.NumField(), len(signedApprovalFields))
+	}
+	perturb := map[string]func(*SignedApproval){
+		"ok":        func(s *SignedApproval) { s.ok = false },
+		"digest":    func(s *SignedApproval) { s.digest[0] ^= 1 },
+		"challenge": func(s *SignedApproval) { s.challenge[0] ^= 1 },
+		"expiry":    func(s *SignedApproval) { s.expiry++ },
+		"name":      func(s *SignedApproval) { s.name += "x" },
+		"format":    func(s *SignedApproval) { s.format += "x" },
+		"key":       func(s *SignedApproval) { s.key = bytes.Clone(s.key); s.key[0] ^= 0xff },
+	}
+	if len(perturb) != len(signedApprovalFields) {
+		t.Fatalf("%d perturbations for %d fields", len(perturb), len(signedApprovalFields))
+	}
+	observe := func(s *SignedApproval) string {
+		n, f, k, ok := s.Use()
+		return fmt.Sprintf("%v|%s|%s|%s|%x|%v", s.Binds(fixtureDigest(t), fixtureChallenge(), fixtureExpiry), s.ApprovedBy(), n, f, k, ok)
+	}
+	base := observe(verifiedApproval(t))
+	for field, p := range perturb {
+		s := verifiedApproval(t)
+		p(s)
+		if observe(s) == base {
+			t.Errorf("changing %s changed nothing Binds, ApprovedBy or Use report", field)
 		}
 	}
 }
````

### Patch G1 — the store-flag tests

````diff G1
--- a/cmd/ailang-worldd/approver_stores_test.go
+++ b/cmd/ailang-worldd/approver_stores_test.go
@@ -0,0 +1,188 @@
+package main
+
+// Row 201 A1b sprint 3 (plan design_docs/planned/w-world-approvals-by-key-a1b-sprint3-plan.md
+// G1): approver enroll|list|revoke and approve --list act on the stores their flags name, and on
+// the default --db when none is given. Sprint 2's judge found that no test passed
+// --credentials-db or relied on the default --db (mutants F-CREDS, AS-CREDFOR and F-DB survived).
+// Every store is under t.TempDir(), and HOME is a temp dir (D-WORLD-109).
+
+import (
+	"bytes"
+	"context"
+	"errors"
+	"os"
+	"path/filepath"
+	"strings"
+	"testing"
+	"time"
+
+	"github.com/sunholo-data/ailang-world/host/store"
+)
+
+// elsewhereStores is an existing world store and a fresh schema-5 credentials store in ANOTHER
+// directory: a verb that ignored --credentials-db would look beside db and find nothing there.
+func elsewhereStores(t *testing.T) (db, creds string) {
+	t.Helper()
+	db = dbPathOf(t)
+	creds = filepath.Join(t.TempDir(), "elsewhere-credentials.db")
+	st, err := store.OpenCredentials(creds)
+	if err != nil {
+		t.Fatal(err)
+	}
+	if err := st.Close(); err != nil {
+		t.Fatal(err)
+	}
+	return db, creds
+}
+
+// homeStores points HOME at a temp dir and makes the default world store
+// ($HOME/.ailang/world/world.db) and the credentials store beside it.
+func homeStores(t *testing.T) (db, creds string) {
+	t.Helper()
+	home := t.TempDir()
+	t.Setenv("HOME", home)
+	db = filepath.Join(home, ".ailang", "world", "world.db")
+	if got := defaultWorldPaths().db; got != db {
+		t.Fatalf("control: the default --db is %q, want %q", got, db)
+	}
+	if err := os.MkdirAll(filepath.Dir(db), 0o700); err != nil {
+		t.Fatal(err)
+	}
+	s, err := store.Open(db)
+	if err != nil {
+		t.Fatal(err)
+	}
+	if err := s.Close(); err != nil {
+		t.Fatal(err)
+	}
+	creds, err = credentialsFor(db, "")
+	if err != nil {
+		t.Fatal(err)
+	}
+	c, err := store.OpenCredentials(creds)
+	if err != nil {
+		t.Fatal(err)
+	}
+	if err := c.Close(); err != nil {
+		t.Fatal(err)
+	}
+	return db, creds
+}
+
+// assertNoStoreBeside fails if a credentials store exists beside db (the verb used the default).
+func assertNoStoreBeside(t *testing.T, db string) {
+	t.Helper()
+	beside, err := credentialsFor(db, "")
+	if err != nil {
+		t.Fatal(err)
+	}
+	if _, err := os.Stat(beside); !errors.Is(err, os.ErrNotExist) {
+		t.Fatalf("a credentials store exists beside the world store (%s, err %v): --credentials-db was ignored", beside, err)
+	}
+}
+
+// G1: each approver subcommand reads and writes the --credentials-db it is given, and the default
+// --db's store when no --db is given.
+func TestApproverVerbsUseTheNamedStores(t *testing.T) {
+	ed := fixtureKeyOf(t, "ed25519.pub")
+	wantLine := "mark-mac  ssh  ssh-ed25519 " + edFingerprint + "  enrolled 2027-01-15T08:00:00Z  last used never  active\n"
+	yes := func() sessionEnv {
+		return approveEnv(&lockedTerm{in: strings.NewReader("y\n")}, nil, atTime(approveT0+100))
+	}
+	type arm struct {
+		name   string
+		stores func(t *testing.T) (db, creds string)
+		flags  func(db, creds string) []string // the store flags this arm passes
+		beside bool                            // assert no store was created beside db
+	}
+	arms := []arm{
+		{"credentials-db", elsewhereStores, func(db, creds string) []string { return []string{"--db", db, "--credentials-db", creds} }, true},
+		{"default-db", homeStores, func(string, string) []string { return nil }, false},
+	}
+	for _, a := range arms {
+		t.Run("list/"+a.name, func(t *testing.T) {
+			db, creds := a.stores(t)
+			seedApprover(t, creds, "mark-mac", ed, approveT0)
+			var out, errw bytes.Buffer
+			code := runApprover(append([]string{"list"}, a.flags(db, creds)...), &out, &errw, approveEnv(nil, nil, atTime(approveT0)))
+			if code != exitOK || out.String() != wantLine {
+				t.Fatalf("code %d, stdout %q, want exitOK and %q; stderr %q", code, out.String(), wantLine, errw.String())
+			}
+			if a.beside {
+				assertNoStoreBeside(t, db)
+			}
+		})
+		t.Run("enroll/"+a.name, func(t *testing.T) {
+			db, creds := a.stores(t)
+			args := append([]string{"enroll", "--name", "mark-mac", "--format", "ssh", "--pubkey", fixturePub("ed25519.pub")}, a.flags(db, creds)...)
+			var out, errw bytes.Buffer
+			if code := runApprover(args, &out, &errw, yes()); code != exitOK {
+				t.Fatalf("code %d, want exitOK; stderr %q", code, errw.String())
+			}
+			rows := approverRows(t, creds)
+			if len(rows) != 1 || rows[0].Name != "mark-mac" || !bytes.Equal(rows[0].PublicKey, ed.Blob) {
+				t.Fatalf("approver rows in %s: %+v, want mark-mac's ed25519 key", creds, rows)
+			}
+			if a.beside {
+				assertNoStoreBeside(t, db)
+			}
+		})
+		t.Run("revoke/"+a.name, func(t *testing.T) {
+			db, creds := a.stores(t)
+			seedApprover(t, creds, "mark-mac", ed, approveT0)
+			args := append([]string{"revoke", "mark-mac"}, a.flags(db, creds)...)
+			var out, errw bytes.Buffer
+			if code := runApprover(args, &out, &errw, yes()); code != exitOK {
+				t.Fatalf("code %d, want exitOK; stderr %q", code, errw.String())
+			}
+			if rows := approverRows(t, creds); len(rows) != 1 || rows[0].RevokedAt != approveT0+100 {
+				t.Fatalf("approver rows in %s: %+v, want mark-mac revoked at %d", creds, rows, approveT0+100)
+			}
+			if a.beside {
+				assertNoStoreBeside(t, db)
+			}
+		})
+	}
+}
+
+// putPendingRequest stores one pending request for episode ep1 (id mr-0102030405060708) in creds.
+func putPendingRequest(t *testing.T, creds string) {
+	t.Helper()
+	st, err := store.OpenCredentials(creds)
+	if err != nil {
+		t.Fatal(err)
+	}
+	defer func() { _ = st.Close() }()
+	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
+	defer cancel()
+	row := store.MintRequestRow{RequestID: firstRequestID, Request: []byte("{}"), Digest: strings.Repeat("ab", 32),
+		Challenge: seqChallenge(), EpisodeID: "ep1", CreatedAt: approveT0, ExpiresAt: approveT0 + 600}
+	if err := st.PutMintRequest(ctx, row, approveT0); err != nil {
+		t.Fatal(err)
+	}
+}
+
+// G2: approve --list reads the --credentials-db it is given, and the default --db's store when no
+// --db is given.
+func TestApproveListUsesTheNamedStores(t *testing.T) {
+	want := firstRequestID + "  episode ep1  digest abababababababab  expires 2027-01-15T08:10:00Z\n"
+	t.Run("credentials-db", func(t *testing.T) {
+		db, creds := elsewhereStores(t)
+		putPendingRequest(t, creds)
+		var out, errw bytes.Buffer
+		code := runApprove([]string{"--list", "--db", db, "--credentials-db", creds}, &out, &errw, approveEnv(nil, nil, atTime(approveT0+1)))
+		if code != exitOK || out.String() != want {
+			t.Fatalf("code %d, stdout %q, want exitOK and %q; stderr %q", code, out.String(), want, errw.String())
+		}
+		assertNoStoreBeside(t, db)
+	})
+	t.Run("default-db", func(t *testing.T) {
+		_, creds := homeStores(t)
+		putPendingRequest(t, creds)
+		var out, errw bytes.Buffer
+		code := runApprove([]string{"--list"}, &out, &errw, approveEnv(nil, nil, atTime(approveT0+1)))
+		if code != exitOK || out.String() != want {
+			t.Fatalf("code %d, stdout %q, want exitOK and %q; stderr %q", code, out.String(), want, errw.String())
+		}
+	})
+}
````

### Patch M4c1 — ATOMIC: challenge binding, approved_by derivation, every caller

````diff M4c1
--- a/cmd/ailang-worldd/approve.go
+++ b/cmd/ailang-worldd/approve.go
@@ -164,6 +164,10 @@
 	if refusal != "" {
 		return refuse("%s", refusal)
 	}
+	challenge, err := hex.DecodeString(row.Challenge)
+	if err != nil || len(challenge) != 32 { // F4: the consume binds the row's challenge (row 201 A1b §2.7)
+		return refuse("request %s's stored challenge is not 64 hex; refusing; fix: make a new request with episode start --request", row.RequestID)
+	}
 
 	// 2. The same store: the request was made for this world store and credentials file.
 	w, werr := store.CanonicalPath(*db)
@@ -210,8 +214,8 @@
 	req.presetPrompt = renderRequest(row.RequestID, mr, row.Digest, row.ExpiresAt, now)
 	req.afterToken = mcpAfterToken(mr.MCPConfig, mr.Daemon.Addr)
 	req.answerBy = row.ExpiresAt
-	req.consume = &authority.RequestMint{Request: mr, RequestCreatedAt: row.CreatedAt, Consume: store.RequestConsume{
-		ID: row.RequestID, Digest: row.Digest, ApprovedBy: "tty", Clock: env.now}}
+	req.consume = &authority.RequestMint{Request: mr, RequestCreatedAt: row.CreatedAt, Challenge: [32]byte(challenge),
+		Consume: store.RequestConsume{ID: row.RequestID, Digest: row.Digest, Clock: env.now}}
 	if code := provisionEpisode(req, stdout, stderr, env); code != exitOK {
 		return code
 	}
--- a/cmd/ailang-worldd/mint_request.go
+++ b/cmd/ailang-worldd/mint_request.go
@@ -220,7 +220,7 @@
 	"startPoint":         "fixed: \"\" (set only by episode_preset.go)",
 	"existingBranch":     "fixed: false (set only by episode_preset.go)",
 	"presetPrompt":       "excluded: display only; approve sets renderRequest(request)",
-	"consume":            "excluded: set by approve from the stored row (id, digest, approved_by tty, env.now)",
+	"consume":            "excluded: set by approve from the stored row (id, digest, challenge, env.now); authority.MintForRequest derives approved_by",
 	"answerBy":           "excluded: set by approve to the stored row's expires_at",
 }
 
--- a/host/authority/mint.go
+++ b/host/authority/mint.go
@@ -5,6 +5,7 @@
 	"crypto/rand"
 	"encoding/hex"
 	"encoding/json"
+	"errors"
 	"fmt"
 	"io"
 
@@ -56,13 +57,38 @@
 }
 
 // RequestMint is what a mint from a stored request needs beyond mintRecorded's arguments
-// (row 201 §2.6): the decoded request, its row's created_at, and the consume.
+// (row 201 §2.6): the decoded request, its row's created_at and challenge, and the consume.
 type RequestMint struct {
 	Request          approval.MintRequest
 	RequestCreatedAt int64
-	Consume          store.RequestConsume // ID, Digest, ApprovedBy, Clock; MintForRequest sets Check
+	Challenge        [32]byte             // the stored row's challenge (row 201 A1b §2.6)
+	Consume          store.RequestConsume // ID, Digest, Clock; MintForRequest sets Challenge, ApprovedBy and Check
 }
 
+var (
+	// ErrApprovedBySetByCaller: approved_by is derived from how the request was approved, never
+	// named by the caller (row 201 A1b §2.6 A1).
+	ErrApprovedBySetByCaller = errors.New("authority: approved_by is derived by the mint, never set by the caller")
+	// ErrBadConsumeDigest: the consume's digest is not 64 hex (A2).
+	ErrBadConsumeDigest = errors.New("authority: the consume's digest is not 64 hex")
+)
+
+// requestConsume is what the store binds for rm (row 201 A1b §2.6): the caller's consume, with the
+// challenge from rm and approved_by derived here.
+func requestConsume(rm *RequestMint) (store.RequestConsume, error) {
+	consume := rm.Consume
+	if consume.ApprovedBy != "" { // A1
+		return store.RequestConsume{}, ErrApprovedBySetByCaller
+	}
+	digest, err := hex.DecodeString(consume.Digest)
+	if err != nil || len(digest) != 32 { // A2
+		return store.RequestConsume{}, ErrBadConsumeDigest
+	}
+	consume.Challenge = hex.EncodeToString(rm.Challenge[:])
+	consume.ApprovedBy = "tty" // A5: the typed y
+	return consume, nil
+}
+
 // MintForRequest mints exactly the session the request describes and consumes the request in the
 // same transaction. approval.MintMatches checks all seven columns of the row it built inside that
 // transaction, with the in-transaction clock reading as created_at's upper bound, before the row
@@ -112,7 +138,10 @@
 	}
 	switch {
 	case rm != nil:
-		consume := rm.Consume
+		var consume store.RequestConsume
+		if consume, err = requestConsume(rm); err != nil {
+			return "", "", store.SessionRow{}, err
+		}
 		consume.Check = func(consumedAt int64) error {
 			return approval.MintMatches(rm.Request, approval.MintFacts{
 				CredentialID: row.CredentialID, EpisodeID: row.EpisodeID, GrantsJSON: row.GrantsJSON,
--- a/host/authority/mint_request_test.go
+++ b/host/authority/mint_request_test.go
@@ -43,6 +43,10 @@
 	return mr, row
 }
 
+// au1Challenge is au1Request's row challenge, strings.Repeat("ab", 32), as bytes.
+var au1Challenge = [32]byte{0xab, 0xab, 0xab, 0xab, 0xab, 0xab, 0xab, 0xab, 0xab, 0xab, 0xab, 0xab, 0xab, 0xab, 0xab, 0xab,
+	0xab, 0xab, 0xab, 0xab, 0xab, 0xab, 0xab, 0xab, 0xab, 0xab, 0xab, 0xab, 0xab, 0xab, 0xab, 0xab}
+
 // au1Grants is the grant list provisionEpisode builds: ExpiresAt is now + ttl.
 func au1Grants(budget, expiresAt int64) []broker.Capability {
 	return []broker.Capability{{Effect: "Workspace.Read", Scope: "worktree", Budget: budget, ExpiresAt: expiresAt}}
@@ -53,7 +57,7 @@
 func TestMintForRequestChecksTheRowInsideTheTransaction(t *testing.T) {
 	good := ProfileBinding{Project: "proj", Digest: "sha256:" + strings.Repeat("c", 64)}
 	consumeAt := func(st *store.Credentials, mr approval.MintRequest, row store.MintRequestRow, at int64) store.RequestConsume {
-		return store.RequestConsume{ID: row.RequestID, Digest: row.Digest, ApprovedBy: "tty", Clock: func() int64 { return at }}
+		return store.RequestConsume{ID: row.RequestID, Digest: row.Digest, Clock: func() int64 { return at }}
 	}
 	t.Run("good", func(t *testing.T) {
 		st, err := store.OpenCredentials(filepath.Join(t.TempDir(), store.CredentialsFileName))
@@ -64,7 +68,7 @@
 		mr, row := au1Request(t, st)
 		ctx := boundedTestContext(t)
 		raw, credID, sess, err := MintForRequest(ctx, st, "ep", good, au1Grants(5, 1610), 600, 1010, nil,
-			RequestMint{Request: mr, RequestCreatedAt: 1000, Consume: consumeAt(st, mr, row, 1020)})
+			RequestMint{Request: mr, RequestCreatedAt: 1000, Challenge: au1Challenge, Consume: consumeAt(st, mr, row, 1020)})
 		if err != nil {
 			t.Fatal(err)
 		}
@@ -106,7 +110,7 @@
 			mr, row := au1Request(t, st)
 			ctx := boundedTestContext(t)
 			_, _, _, err = MintForRequest(ctx, st, m.episode, m.binding, m.grants, m.ttl, m.now, nil,
-				RequestMint{Request: mr, RequestCreatedAt: 1000, Consume: consumeAt(st, mr, row, m.clockAt)})
+				RequestMint{Request: mr, RequestCreatedAt: 1000, Challenge: au1Challenge, Consume: consumeAt(st, mr, row, m.clockAt)})
 			var mm *approval.MintMismatchError
 			if !errors.As(err, &mm) || mm.Column != m.column {
 				t.Fatalf("err %v, want a MintMismatchError on %s", err, m.column)
@@ -129,7 +133,7 @@
 		defer func() { _ = st.Close() }()
 		mr, row := au1Request(t, st)
 		_, _, _, err = MintForRequest(boundedTestContext(t), st, "ep", good, au1Grants(5, 1610), 600, 1010, nil,
-			RequestMint{Request: mr, RequestCreatedAt: 1000, Consume: consumeAt(st, mr, row, 1600)})
+			RequestMint{Request: mr, RequestCreatedAt: 1000, Challenge: au1Challenge, Consume: consumeAt(st, mr, row, 1600)})
 		if !errors.Is(err, store.ErrRequestExpired) {
 			t.Fatalf("err %v, want ErrRequestExpired", err)
 		}
--- a/host/store/mint_requests.go
+++ b/host/store/mint_requests.go
@@ -53,7 +53,10 @@
 	ErrRequestClockBehind     = errors.New("store: the clock is behind the mint request's creation time")
 	ErrRequestDigestChanged   = errors.New("store: the mint request's digest is not the one approved")
 	ErrBadMintRequest         = errors.New("store: invalid mint request row")
-	ErrIncompleteConsume      = errors.New("store: a mint request consume needs an id, a digest, approved_by, a clock and a check")
+	ErrIncompleteConsume      = errors.New("store: a mint request consume needs an id, a digest, a challenge, approved_by, a clock and a check")
+	// ErrRequestChallengeChanged: the row's challenge is not the one the approval was made for
+	// (row 201 A1b §2.5 C6, C7).
+	ErrRequestChallengeChanged = errors.New("store: the mint request's challenge is not the one signed")
 )
 
 // OlderCredentialsStoreError refuses an operation on a store below the version its Feature needs:
@@ -233,9 +236,11 @@
 // RequestConsume is what MintSessionForRequest consumes. Clock is read exactly once, inside the
 // transaction, after BeginTx (design §2.4 [r3]); that reading is consumed_at and the expiry
 // bound. Check runs with that reading before the guarded UPDATE; an error rolls everything back
-// (plan §1.3 D2: authority.MintForRequest puts approval.MintMatches here).
+// (plan §1.3 D2: authority.MintForRequest puts approval.MintMatches here). Challenge is the row's
+// challenge (64 lowercase hex), bound in the guarded UPDATE (row 201 A1b §2.5 C1, C6).
 type RequestConsume struct {
 	ID, Digest, ApprovedBy string
+	Challenge              string
 	Clock                  func() int64
 	Check                  func(consumedAt int64) error
 }
@@ -247,7 +252,7 @@
 	if err := requireDeadline(ctx); err != nil {
 		return err
 	}
-	if rc.ID == "" || rc.Digest == "" || rc.ApprovedBy == "" || rc.Clock == nil || rc.Check == nil {
+	if rc.ID == "" || rc.Digest == "" || !isLowerHex64(rc.Challenge) || rc.ApprovedBy == "" || rc.Clock == nil || rc.Check == nil { // C1
 		return ErrIncompleteConsume
 	}
 	if session.CredentialID == "" || session.EpisodeID == "" {
@@ -265,8 +270,8 @@
 	if err := rc.Check(at); err != nil {
 		return err
 	}
-	res, err := tx.ExecContext(ctx, `UPDATE mint_requests SET consumed_at = ?, approved_by = ?, credential_id = ? WHERE request_id = ? AND digest = ? AND consumed_at IS NULL AND created_at <= ? AND expires_at > ?`,
-		at, rc.ApprovedBy, session.CredentialID, rc.ID, rc.Digest, at, at)
+	res, err := tx.ExecContext(ctx, `UPDATE mint_requests SET consumed_at = ?, approved_by = ?, credential_id = ? WHERE request_id = ? AND digest = ? AND challenge = ? AND consumed_at IS NULL AND created_at <= ? AND expires_at > ?`,
+		at, rc.ApprovedBy, session.CredentialID, rc.ID, rc.Digest, rc.Challenge, at, at) // C6
 	if err != nil {
 		return fmt.Errorf("store: consume mint request: %w", err)
 	}
@@ -289,12 +294,12 @@
 // classifyConsume names why the guarded UPDATE changed no row, reading the request inside tx.
 func classifyConsume(ctx context.Context, tx *sql.Tx, rc RequestConsume, at int64, changed int64) error {
 	var (
-		digest           string
-		createdAt, expAt int64
-		consumed         bool
+		digest, challenge string
+		createdAt, expAt  int64
+		consumed          bool
 	)
-	err := tx.QueryRowContext(ctx, `SELECT digest, created_at, expires_at, consumed_at IS NOT NULL FROM mint_requests WHERE request_id = ?`, rc.ID).
-		Scan(&digest, &createdAt, &expAt, &consumed)
+	err := tx.QueryRowContext(ctx, `SELECT digest, challenge, created_at, expires_at, consumed_at IS NOT NULL FROM mint_requests WHERE request_id = ?`, rc.ID).
+		Scan(&digest, &challenge, &createdAt, &expAt, &consumed)
 	if errors.Is(err, sql.ErrNoRows) {
 		return fmt.Errorf("%w: %s", ErrRequestUnknown, rc.ID)
 	}
@@ -310,6 +315,8 @@
 		return fmt.Errorf("%w: %s", ErrRequestExpired, rc.ID)
 	case digest != rc.Digest:
 		return fmt.Errorf("%w: %s", ErrRequestDigestChanged, rc.ID)
+	case challenge != rc.Challenge: // C7
+		return fmt.Errorf("%w: %s", ErrRequestChallengeChanged, rc.ID)
 	}
 	return fmt.Errorf("store: mint request %s: the guarded update changed %d rows", rc.ID, changed)
 }
--- a/host/store/mint_requests_test.go
+++ b/host/store/mint_requests_test.go
@@ -347,12 +347,13 @@
 	}
 }
 
-// okConsume returns a RequestConsume for row with ApprovedBy "tty", a Clock returning at and
-// counting calls into *clocks, and a Check returning nil and counting into *checks.
+// okConsume returns a RequestConsume for row with its challenge, ApprovedBy "tty", a Clock
+// returning at and counting calls into *clocks, and a Check returning nil and counting into *checks.
 func okConsume(row MintRequestRow, at int64, clocks, checks *int) RequestConsume {
 	return RequestConsume{
 		ID:         row.RequestID,
 		Digest:     row.Digest,
+		Challenge:  row.Challenge,
 		ApprovedBy: "tty",
 		Clock:      func() int64 { *clocks++; return at },
 		Check:      func(int64) error { *checks++; return nil },
@@ -580,7 +581,7 @@
 			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
 			defer cancel()
 			rc := RequestConsume{
-				ID: r.RequestID, Digest: r.Digest, ApprovedBy: "tty",
+				ID: r.RequestID, Digest: r.Digest, Challenge: r.Challenge, ApprovedBy: "tty",
 				Clock: func() int64 { return 1010 },
 				Check: func(int64) error { return nil },
 			}
````

### Patch M4c2 — M4c1's tests

````diff M4c2
--- a/cmd/ailang-worldd/approve_sig_test.go
+++ b/cmd/ailang-worldd/approve_sig_test.go
@@ -0,0 +1,36 @@
+package main
+
+// Row 201 A1b sprint 3 (plan design_docs/planned/w-world-approvals-by-key-a1b-sprint3-plan.md §3
+// M4c2, M5a-M5c): approve binds the stored challenge, and approves by an enrolled key's signature.
+// Every store is under t.TempDir(); every key is approvaltest.NewSigner's or a checked-in fixture
+// (D-WORLD-109).
+
+import (
+	"bytes"
+	"strings"
+	"testing"
+)
+
+// CR11 (M4c2): a stored challenge that is not 64 hex is refused before any terminal opens (F4), and
+// the request stays pending.
+func TestApproveRefusesANonHexChallenge(t *testing.T) {
+	a := newApproveRig(t, "ep1", approveT0)
+	raw := rawCredentials(t, a.credPath)
+	if _, err := raw.Exec(`UPDATE mint_requests SET challenge = ? WHERE request_id = ?`, strings.Repeat("zz", 32), a.id); err != nil {
+		t.Fatal(err)
+	}
+	_ = raw.Close()
+	var opened bool
+	var stdout, stderr bytes.Buffer
+	code := runApprove(a.args(), &stdout, &stderr, approveEnv(&promptTerm{in: strings.NewReader("y\n")}, &opened, atTime(approveT0+1)))
+	if code != exitUsage || opened {
+		t.Fatalf("rc %d, terminal opened %v; want exitUsage and no terminal\nstderr %s", code, opened, stderr.String())
+	}
+	want := "ailang-worldd approve: request " + a.id + "'s stored challenge is not 64 hex; refusing; fix: make a new request with episode start --request\n"
+	if stderr.String() != want {
+		t.Fatalf("stderr %q, want %q", stderr.String(), want)
+	}
+	if a.row(t).Consumed || len(a.sessionRows(t)) != 0 {
+		t.Fatal("a request with a non-hex challenge was consumed or minted")
+	}
+}
--- a/host/authority/mint_derive_test.go
+++ b/host/authority/mint_derive_test.go
@@ -0,0 +1,91 @@
+package authority
+
+// Row 201 A1b sprint 3 (plan design_docs/planned/w-world-approvals-by-key-a1b-sprint3-plan.md §3
+// M4c2, M4e): MintForRequest derives approved_by and binds the row's challenge; its caller names
+// neither.
+
+import (
+	"errors"
+	"path/filepath"
+	"strings"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/store"
+)
+
+// au1Store is a fresh credentials store holding au1Request's row (created 1000, expires 1600).
+func au1Store(t *testing.T) (*store.Credentials, approval.MintRequest, store.MintRequestRow) {
+	t.Helper()
+	st, err := store.OpenCredentials(filepath.Join(t.TempDir(), store.CredentialsFileName))
+	if err != nil {
+		t.Fatal(err)
+	}
+	t.Cleanup(func() { _ = st.Close() })
+	mr, row := au1Request(t, st)
+	return st, mr, row
+}
+
+// au1RequestMint is the RequestMint approve builds for au1Request's row: no approved_by.
+func au1RequestMint(mr approval.MintRequest, row store.MintRequestRow) RequestMint {
+	return RequestMint{Request: mr, RequestCreatedAt: 1000, Challenge: au1Challenge,
+		Consume: store.RequestConsume{ID: row.RequestID, Digest: row.Digest}}
+}
+
+// au1Mint mints au1Request's session (episode ep, the profile binding, at 1010) through rm, the
+// consume's clock reading 1020.
+func au1Mint(t *testing.T, st *store.Credentials, rm RequestMint) error {
+	t.Helper()
+	good := ProfileBinding{Project: "proj", Digest: "sha256:" + strings.Repeat("c", 64)}
+	rm.Consume.Clock = func() int64 { return 1020 }
+	_, _, _, err := MintForRequest(boundedTestContext(t), st, "ep", good, au1Grants(5, 1610), 600, 1010, nil, rm)
+	return err
+}
+
+// assertNotMinted fails unless row is still pending in st and st holds no session row.
+func assertNotMinted(t *testing.T, st *store.Credentials, row store.MintRequestRow) {
+	t.Helper()
+	ctx := boundedTestContext(t)
+	if got, ok, err := st.MintRequest(ctx, row.RequestID); err != nil || !ok || got.Consumed {
+		t.Fatalf("request after a refused mint: %+v ok %v err %v, want it pending", got, ok, err)
+	}
+	if sessions, err := st.ListSessions(ctx, "", store.MaxListSessions); err != nil || len(sessions) != 0 {
+		t.Fatalf("%d session rows after a refused mint (err %v), want 0", len(sessions), err)
+	}
+}
+
+// AU1 (M4c2): a caller-set approved_by is refused (A1), a consume digest that is not 64 hex is
+// refused (A2), the row's challenge is bound, and the typed mint records tty (A5).
+func TestMintForRequestDerivesApprovedBy(t *testing.T) {
+	for _, a := range []struct {
+		name   string
+		mutate func(*RequestMint)
+		want   error
+	}{
+		{"caller-tty", func(rm *RequestMint) { rm.Consume.ApprovedBy = "tty" }, ErrApprovedBySetByCaller},
+		{"caller-ssh", func(rm *RequestMint) { rm.Consume.ApprovedBy = "ssh:mallory" }, ErrApprovedBySetByCaller},
+		{"digest-not-hex", func(rm *RequestMint) { rm.Consume.Digest = strings.Repeat("zz", 32) }, ErrBadConsumeDigest},
+		{"digest-31-bytes", func(rm *RequestMint) { rm.Consume.Digest = strings.Repeat("ab", 31) }, ErrBadConsumeDigest},
+		{"other-challenge", func(rm *RequestMint) { rm.Challenge[0] ^= 1 }, store.ErrRequestChallengeChanged},
+	} {
+		t.Run(a.name, func(t *testing.T) {
+			st, mr, row := au1Store(t)
+			rm := au1RequestMint(mr, row)
+			a.mutate(&rm)
+			if err := au1Mint(t, st, rm); !errors.Is(err, a.want) {
+				t.Fatalf("err %v, want %v", err, a.want)
+			}
+			assertNotMinted(t, st, row)
+		})
+	}
+	t.Run("typed-records-tty", func(t *testing.T) {
+		st, mr, row := au1Store(t)
+		if err := au1Mint(t, st, au1RequestMint(mr, row)); err != nil {
+			t.Fatal(err)
+		}
+		got, ok, err := st.MintRequest(boundedTestContext(t), row.RequestID)
+		if err != nil || !ok || !got.Consumed || got.ApprovedBy != "tty" {
+			t.Fatalf("request after the typed mint: %+v ok %v err %v, want consumed by tty", got, ok, err)
+		}
+	})
+}
--- a/host/store/consume_binding_test.go
+++ b/host/store/consume_binding_test.go
@@ -0,0 +1,86 @@
+package store
+
+// Row 201 A1b sprint 3 (plan design_docs/planned/w-world-approvals-by-key-a1b-sprint3-plan.md §3
+// M4c2, M4d1, M4d2): the consume binds the request's challenge and, on the signature path, its
+// approver's active row, in the transaction that consumes the request and inserts the session.
+
+import (
+	"errors"
+	"strings"
+	"testing"
+)
+
+// putTestRequest stores testMintRequestRow(1, 1000) in c and returns it.
+func putTestRequest(t *testing.T, c *Credentials) MintRequestRow {
+	t.Helper()
+	r := testMintRequestRow(1, 1000)
+	if err := c.PutMintRequest(testCtx(t), r, 1000); err != nil {
+		t.Fatal(err)
+	}
+	return r
+}
+
+// assertPendingNoSession fails unless r is still pending in c and c holds no session row.
+func assertPendingNoSession(t *testing.T, c *Credentials, r MintRequestRow) {
+	t.Helper()
+	ctx := testCtx(t)
+	if got, ok, err := c.MintRequest(ctx, r.RequestID); err != nil || !ok || got.Consumed {
+		t.Fatalf("request %s after a refused consume: %+v ok %v err %v, want it pending", r.RequestID, got, ok, err)
+	}
+	if rows, err := c.ListSessions(ctx, "", MaxListSessions); err != nil || len(rows) != 0 {
+		t.Fatalf("%d session rows after a refused consume (err %v), want 0", len(rows), err)
+	}
+}
+
+// CS1 (M4c2): a consume whose challenge is not 64 lowercase hex is incomplete (C1), refused
+// before the clock is read.
+func TestMintSessionForRequestNeedsTheChallenge(t *testing.T) {
+	for name, ch := range map[string]string{
+		"empty":   "",
+		"63-hex":  strings.Repeat("c", 63),
+		"65-hex":  strings.Repeat("c", 65),
+		"upper":   strings.Repeat("CD", 32),
+		"not-hex": strings.Repeat("zz", 32),
+	} {
+		t.Run(name, func(t *testing.T) {
+			c := openTestCredentials(t)
+			r := putTestRequest(t, c)
+			var clocks, checks int
+			rc := okConsume(r, 1010, &clocks, &checks)
+			rc.Challenge = ch
+			if err := c.MintSessionForRequest(testCtx(t), consumeSession(1), rc); !errors.Is(err, ErrIncompleteConsume) {
+				t.Fatalf("err %v, want ErrIncompleteConsume", err)
+			}
+			if clocks != 0 {
+				t.Fatalf("Clock read %d times before an incomplete-consume refusal, want 0", clocks)
+			}
+			assertPendingNoSession(t, c, r)
+		})
+	}
+}
+
+// CS5 (M4c2): the guarded UPDATE binds the row's challenge (C6), and the refusal names it (C7),
+// after the digest: a consume with both wrong names the digest.
+func TestConsumeBindsTheChallenge(t *testing.T) {
+	other := strings.Repeat("ef", 32)
+	for _, a := range []struct {
+		name   string
+		mutate func(*RequestConsume)
+		want   error
+	}{
+		{"other-challenge", func(rc *RequestConsume) { rc.Challenge = other }, ErrRequestChallengeChanged},
+		{"other-digest-and-challenge", func(rc *RequestConsume) { rc.Digest, rc.Challenge = other, other }, ErrRequestDigestChanged},
+	} {
+		t.Run(a.name, func(t *testing.T) {
+			c := openTestCredentials(t)
+			r := putTestRequest(t, c)
+			var clocks, checks int
+			rc := okConsume(r, 1010, &clocks, &checks)
+			a.mutate(&rc)
+			if err := c.MintSessionForRequest(testCtx(t), consumeSession(1), rc); !errors.Is(err, a.want) {
+				t.Fatalf("err %v, want %v", err, a.want)
+			}
+			assertPendingNoSession(t, c, r)
+		})
+	}
+}
````

### Patch M4d1 — approver validation C2–C5

````diff M4d1
--- a/host/store/consume_binding_test.go
+++ b/host/store/consume_binding_test.go
@@ -6,6 +6,8 @@
 
 import (
 	"errors"
+	"path/filepath"
+	"reflect"
 	"strings"
 	"testing"
 )
@@ -79,6 +81,159 @@
 			a.mutate(&rc)
 			if err := c.MintSessionForRequest(testCtx(t), consumeSession(1), rc); !errors.Is(err, a.want) {
 				t.Fatalf("err %v, want %v", err, a.want)
+			}
+			assertPendingNoSession(t, c, r)
+		})
+	}
+}
+
+// ---- M4d1 ----
+
+// credentialsAt is a fresh schema-5 credentials store, its path, and approver x
+// (testApprover("x", 900): format ssh, key "key-x") enrolled in it.
+func credentialsAt(t *testing.T) (*Credentials, string) {
+	t.Helper()
+	path := filepath.Join(t.TempDir(), CredentialsFileName)
+	c, err := OpenCredentials(path)
+	if err != nil {
+		t.Fatal(err)
+	}
+	t.Cleanup(func() { _ = c.Close() })
+	if err := c.EnrolApprover(testCtx(t), testApprover("x", 900)); err != nil {
+		t.Fatal(err)
+	}
+	return c, path
+}
+
+// signedConsume is okConsume(r, at, ...) for a signature approval by x: Approver x, approved_by ssh:x.
+func signedConsume(r MintRequestRow, at int64, clocks, checks *int) RequestConsume {
+	rc := okConsume(r, at, clocks, checks)
+	rc.Approver = &ApproverUse{Name: "x", Format: "ssh", PublicKey: []byte("key-x")}
+	rc.ApprovedBy = "ssh:x"
+	return rc
+}
+
+// approverX is approver x's row in c.
+func approverX(t *testing.T, c *Credentials) ApproverRow {
+	t.Helper()
+	row, ok, err := c.Approver(testCtx(t), "x")
+	if err != nil || !ok {
+		t.Fatalf("approver x: ok %v err %v", ok, err)
+	}
+	return row
+}
+
+// CS2 (M4d1): approved_by is tty exactly when no approver is bound (C3) and <format>:<name> of the
+// bound approver otherwise (C4); an approver needs schema 5 (C5). Each refusal comes before the
+// clock is read and leaves the request pending and no session.
+func TestApprovedByMatchesTheApprover(t *testing.T) {
+	for _, a := range []struct {
+		name   string
+		mutate func(*RequestConsume)
+	}{
+		{"no-approver-ssh", func(rc *RequestConsume) { rc.Approver = nil }},
+		{"approver-tty", func(rc *RequestConsume) { rc.ApprovedBy = "tty" }},
+		{"approver-other-name", func(rc *RequestConsume) { rc.ApprovedBy = "ssh:y" }},
+		{"approver-other-format", func(rc *RequestConsume) { rc.ApprovedBy = "webauthn:x" }},
+		{"approver-format-not-ssh", func(rc *RequestConsume) { rc.Approver.Format = "webauthn" }},
+	} {
+		t.Run(a.name, func(t *testing.T) {
+			c, _ := credentialsAt(t)
+			r := putTestRequest(t, c)
+			var clocks, checks int
+			rc := signedConsume(r, 1010, &clocks, &checks)
+			a.mutate(&rc)
+			if err := c.MintSessionForRequest(testCtx(t), consumeSession(1), rc); !errors.Is(err, ErrInconsistentApproval) {
+				t.Fatalf("err %v, want ErrInconsistentApproval", err)
+			}
+			if clocks != 0 {
+				t.Fatalf("Clock read %d times, want 0", clocks)
+			}
+			assertPendingNoSession(t, c, r)
+			if got := approverX(t, c); got.LastUsedAt != 0 {
+				t.Fatalf("x's last_used_at %d after a refusal, want 0", got.LastUsedAt)
+			}
+		})
+	}
+	t.Run("v4-store", func(t *testing.T) {
+		path := filepath.Join(t.TempDir(), CredentialsFileName)
+		makeV4Credentials(t, path, nil)
+		c, err := OpenCredentials(path)
+		if err != nil {
+			t.Fatal(err)
+		}
+		defer func() { _ = c.Close() }()
+		r := putTestRequest(t, c)
+		var clocks, checks int
+		err = c.MintSessionForRequest(testCtx(t), consumeSession(1), signedConsume(r, 1010, &clocks, &checks))
+		var older *OlderCredentialsStoreError
+		if !errors.As(err, &older) || older.Need != 5 || older.Feature != "approvers" {
+			t.Fatalf("err %v, want *OlderCredentialsStoreError{Need: 5, Feature: approvers}", err)
+		}
+		if clocks != 0 {
+			t.Fatalf("Clock read %d times on a schema-4 store, want 0", clocks)
+		}
+		assertPendingNoSession(t, c, r)
+	})
+}
+
+// requestConsumeFields is the disposition of every RequestConsume field and of its element type
+// ApproverUse's fields (RC3). "required" fields are zeroed one at a time below.
+var requestConsumeFields = map[string]string{
+	"RequestConsume.ID":         "required",
+	"RequestConsume.Digest":     "required",
+	"RequestConsume.Challenge":  "required (C1)",
+	"RequestConsume.ApprovedBy": "derived and checked: tty with no approver, <format>:<name> with one (C3, C4; TestApprovedByMatchesTheApprover)",
+	"RequestConsume.Approver":   "optional: nil is the typed y; set, it is bound by C2, C5 and C8",
+	"RequestConsume.Clock":      "required",
+	"RequestConsume.Check":      "required",
+	"ApproverUse.Name":          "required (C2)",
+	"ApproverUse.Format":        "required (C2)",
+	"ApproverUse.PublicKey":     "required (C2)",
+}
+
+// RC3 (M4d1): every RequestConsume and ApproverUse field has a disposition, and zeroing each
+// required one in an otherwise valid signed consume is ErrIncompleteConsume, before the clock.
+func TestRequestConsumeCompleteness(t *testing.T) {
+	n := 0
+	for _, typ := range []reflect.Type{reflect.TypeOf(RequestConsume{}), reflect.TypeOf(ApproverUse{})} {
+		for i := 0; i < typ.NumField(); i++ {
+			if key := typ.Name() + "." + typ.Field(i).Name; requestConsumeFields[key] == "" {
+				t.Errorf("%s has no disposition in requestConsumeFields", key)
+			}
+		}
+		n += typ.NumField()
+	}
+	if n != len(requestConsumeFields) {
+		t.Fatalf("RequestConsume and ApproverUse have %d fields, requestConsumeFields %d", n, len(requestConsumeFields))
+	}
+	zero := map[string]func(*RequestConsume){
+		"RequestConsume.ID":        func(rc *RequestConsume) { rc.ID = "" },
+		"RequestConsume.Digest":    func(rc *RequestConsume) { rc.Digest = "" },
+		"RequestConsume.Challenge": func(rc *RequestConsume) { rc.Challenge = "" },
+		"RequestConsume.Clock":     func(rc *RequestConsume) { rc.Clock = nil },
+		"RequestConsume.Check":     func(rc *RequestConsume) { rc.Check = nil },
+		"ApproverUse.Name":         func(rc *RequestConsume) { rc.Approver.Name, rc.ApprovedBy = "", "ssh:" },
+		"ApproverUse.Format":       func(rc *RequestConsume) { rc.Approver.Format, rc.ApprovedBy = "", ":x" },
+		"ApproverUse.PublicKey":    func(rc *RequestConsume) { rc.Approver.PublicKey = nil },
+	}
+	for key, d := range requestConsumeFields {
+		if _, ok := zero[key]; strings.HasPrefix(d, "required") != ok {
+			t.Errorf("%s: disposition %q, zeroing arm present %v", key, d, ok)
+		}
+	}
+	for key, z := range zero {
+		t.Run(key, func(t *testing.T) {
+			c, _ := credentialsAt(t)
+			r := putTestRequest(t, c)
+			var clocks, checks int
+			rc := signedConsume(r, 1010, &clocks, &checks)
+			z(&rc)
+			if err := c.MintSessionForRequest(testCtx(t), consumeSession(1), rc); !errors.Is(err, ErrIncompleteConsume) {
+				t.Fatalf("err %v, want ErrIncompleteConsume", err)
+			}
+			if clocks != 0 {
+				t.Fatalf("Clock read %d times, want 0", clocks)
 			}
 			assertPendingNoSession(t, c, r)
 		})
--- a/host/store/mint_requests.go
+++ b/host/store/mint_requests.go
@@ -57,6 +57,9 @@
 	// ErrRequestChallengeChanged: the row's challenge is not the one the approval was made for
 	// (row 201 A1b §2.5 C6, C7).
 	ErrRequestChallengeChanged = errors.New("store: the mint request's challenge is not the one signed")
+	// ErrInconsistentApproval: approved_by is "tty" exactly when no approver is bound, and
+	// "<format>:<name>" of the bound approver otherwise (row 201 A1b §2.5 C3, C4).
+	ErrInconsistentApproval = errors.New("store: approved_by does not match the approver the consume binds")
 )
 
 // OlderCredentialsStoreError refuses an operation on a store below the version its Feature needs:
@@ -241,10 +244,18 @@
 type RequestConsume struct {
 	ID, Digest, ApprovedBy string
 	Challenge              string
+	Approver               *ApproverUse // nil: the typed y; else the approver whose signature approved (§2.5 C8)
 	Clock                  func() int64
 	Check                  func(consumedAt int64) error
 }
 
+// ApproverUse is the approver a signature consume binds: the consume marks that approver's active
+// row used, with this exact format and key, in its own transaction, or refuses.
+type ApproverUse struct {
+	Name, Format string
+	PublicKey    []byte
+}
+
 // MintSessionForRequest consumes request rc and inserts session in one transaction: the request
 // is marked consumed (consumed_at, approved_by, credential_id) only if it is still pending, and
 // the session row lands with it or not at all.
@@ -255,11 +266,23 @@
 	if rc.ID == "" || rc.Digest == "" || !isLowerHex64(rc.Challenge) || rc.ApprovedBy == "" || rc.Clock == nil || rc.Check == nil { // C1
 		return ErrIncompleteConsume
 	}
+	if a := rc.Approver; a != nil && (a.Name == "" || a.Format == "" || len(a.PublicKey) == 0) { // C2
+		return ErrIncompleteConsume
+	}
+	if rc.Approver == nil && rc.ApprovedBy != "tty" { // C3
+		return ErrInconsistentApproval
+	}
+	if rc.Approver != nil && rc.ApprovedBy != rc.Approver.Format+":"+rc.Approver.Name { // C4
+		return ErrInconsistentApproval
+	}
 	if session.CredentialID == "" || session.EpisodeID == "" {
 		return fmt.Errorf("store: mint session: empty credential_id or episode_id")
 	}
 	if c.version < mintRequestsSchemaVersion {
 		return &OlderCredentialsStoreError{Path: c.display, Version: c.version, Need: mintRequestsSchemaVersion, Feature: "mint requests"}
+	}
+	if rc.Approver != nil && c.version < approversSchemaVersion { // C5
+		return c.olderForApprovers()
 	}
 	tx, err := c.db.BeginTx(ctx, nil)
 	if err != nil {
````

### Patch M4d2 — approver binding in the transaction C8–C12

````diff M4d2
--- a/host/store/approvers.go
+++ b/host/store/approvers.go
@@ -7,6 +7,7 @@
 // by the CLI first; the checks here are the store's own backstop.
 
 import (
+	"bytes"
 	"context"
 	"database/sql"
 	"errors"
@@ -232,5 +233,46 @@
 		return fmt.Errorf("%w: %q enrolled at %d, now %d", ErrApproverClockBehind, name, row.EnrolledAt, now)
 	default:
 		return fmt.Errorf("store: revoke approver %q: changed no row", name)
+	}
+}
+
+// ---- M4d (sprint 3) ----
+
+// useApprover is the consume's approver binding (row 201 A1b §2.5 C8): inside the consume's
+// transaction, at its clock reading, it marks the approver used only while its row is active,
+// holds this exact format and key, and was enrolled by at. Otherwise classifyApproverUse names
+// why (C9-C12), and the caller's rollback undoes the request update.
+func useApprover(ctx context.Context, tx *sql.Tx, a ApproverUse, at int64) error {
+	res, err := tx.ExecContext(ctx, `UPDATE approvers SET last_used_at = ? WHERE name = ? AND format = ? AND public_key = ? AND revoked_at IS NULL AND enrolled_at <= ?`,
+		at, a.Name, a.Format, a.PublicKey, at)
+	if err != nil {
+		return fmt.Errorf("store: use approver: %w", err)
+	}
+	n, err := res.RowsAffected()
+	if err != nil {
+		return fmt.Errorf("store: use approver: %w", err)
+	}
+	if n != 1 {
+		return classifyApproverUse(ctx, tx, a, at)
+	}
+	return nil
+}
+
+// classifyApproverUse names why useApprover changed no row, reading the row in the same transaction.
+func classifyApproverUse(ctx context.Context, tx *sql.Tx, a ApproverUse, at int64) error {
+	row, err := scanApprover(tx.QueryRowContext(ctx, `SELECT `+approverColumns+` FROM approvers WHERE name = ?`, a.Name).Scan)
+	switch {
+	case errors.Is(err, sql.ErrNoRows): // C9
+		return fmt.Errorf("%w: %q", ErrApproverUnknown, a.Name)
+	case err != nil:
+		return fmt.Errorf("store: use approver: %w", err)
+	case row.RevokedAt != 0: // C10
+		return fmt.Errorf("%w: %q at %d", ErrApproverRevoked, a.Name, row.RevokedAt)
+	case row.Format != a.Format || !bytes.Equal(row.PublicKey, a.PublicKey): // C11
+		return fmt.Errorf("%w: %q", ErrApproverKeyChanged, a.Name)
+	case row.EnrolledAt > at: // C12
+		return fmt.Errorf("%w: %q enrolled at %d, now %d", ErrApproverClockBehind, a.Name, row.EnrolledAt, at)
+	default:
+		return fmt.Errorf("store: use approver %q: changed no row", a.Name)
 	}
 }
--- a/host/store/consume_binding_test.go
+++ b/host/store/consume_binding_test.go
@@ -239,3 +239,88 @@
 		})
 	}
 }
+
+// ---- M4d2 ----
+
+// CS3 (M4d2): a signed consume marks its approver used at the in-transaction reading (C8). It is
+// refused, the request pending, no session and last_used_at unchanged, when the approver is
+// revoked by a commit made after BeginTx (from inside the Clock hook: the transaction is deferred,
+// so it holds no lock until its first statement, design V36/V37), unknown (C9), revoked (C10),
+// holding another key or format (C11), or enrolled after the reading (C12).
+func TestSignedConsumeChecksTheApproverInTheTransaction(t *testing.T) {
+	t.Run("active", func(t *testing.T) {
+		c, _ := credentialsAt(t)
+		r := putTestRequest(t, c)
+		var clocks, checks int
+		if err := c.MintSessionForRequest(testCtx(t), consumeSession(1), signedConsume(r, 1010, &clocks, &checks)); err != nil {
+			t.Fatal(err)
+		}
+		got, ok, err := c.MintRequest(testCtx(t), r.RequestID)
+		if err != nil || !ok || !got.Consumed || got.ApprovedBy != "ssh:x" {
+			t.Fatalf("request after the signed consume: %+v ok %v err %v, want consumed by ssh:x", got, ok, err)
+		}
+		if x := approverX(t, c); x.LastUsedAt != 1010 {
+			t.Fatalf("x's last_used_at %d, want 1010 (the in-transaction reading)", x.LastUsedAt)
+		}
+		if rows, err := c.ListSessions(testCtx(t), "", MaxListSessions); err != nil || len(rows) != 1 {
+			t.Fatalf("%d session rows (err %v), want 1", len(rows), err)
+		}
+	})
+	for _, a := range []struct {
+		name  string
+		setup func(t *testing.T, c *Credentials, path string, rc *RequestConsume)
+		want  error
+	}{
+		{"revoked-in-the-window", func(t *testing.T, c *Credentials, path string, rc *RequestConsume) {
+			rc.Clock = func() int64 {
+				if _, err := rawDB(t, path).Exec(`UPDATE approvers SET revoked_at = 1005 WHERE name = 'x'`); err != nil {
+					t.Errorf("revoke on a second connection: %v", err)
+				}
+				return 1010
+			}
+		}, ErrApproverRevoked},
+		{"unknown", func(t *testing.T, c *Credentials, path string, rc *RequestConsume) {
+			rc.Approver.Name, rc.ApprovedBy = "y", "ssh:y"
+		}, ErrApproverUnknown},
+		{"other-key", func(t *testing.T, c *Credentials, path string, rc *RequestConsume) {
+			rc.Approver.PublicKey = []byte("key-y")
+		}, ErrApproverKeyChanged},
+		{"other-format", func(t *testing.T, c *Credentials, path string, rc *RequestConsume) {
+			rc.Approver.Format, rc.ApprovedBy = "webauthn", "webauthn:x"
+		}, ErrApproverKeyChanged},
+		{"enrolled-after-the-reading", func(t *testing.T, c *Credentials, path string, rc *RequestConsume) {
+			if _, err := rawDB(t, path).Exec(`UPDATE approvers SET enrolled_at = 1500 WHERE name = 'x'`); err != nil {
+				t.Fatal(err)
+			}
+		}, ErrApproverClockBehind},
+	} {
+		t.Run(a.name, func(t *testing.T) {
+			c, path := credentialsAt(t)
+			r := putTestRequest(t, c)
+			var clocks, checks int
+			rc := signedConsume(r, 1010, &clocks, &checks)
+			a.setup(t, c, path, &rc)
+			if err := c.MintSessionForRequest(testCtx(t), consumeSession(1), rc); !errors.Is(err, a.want) {
+				t.Fatalf("err %v, want %v", err, a.want)
+			}
+			assertPendingNoSession(t, c, r)
+			if x := approverX(t, c); x.LastUsedAt != 0 {
+				t.Fatalf("x's last_used_at %d after a refusal, want 0", x.LastUsedAt)
+			}
+		})
+	}
+}
+
+// CS4 (M4d2): a typed consume (no approver, approved_by tty) mints on a schema-5 store that holds an
+// approver, and touches no approver row.
+func TestTypedConsumeTouchesNoApprover(t *testing.T) {
+	c, _ := credentialsAt(t)
+	r := putTestRequest(t, c)
+	var clocks, checks int
+	if err := c.MintSessionForRequest(testCtx(t), consumeSession(1), okConsume(r, 1010, &clocks, &checks)); err != nil {
+		t.Fatal(err)
+	}
+	if x := approverX(t, c); x.LastUsedAt != 0 {
+		t.Fatalf("x's last_used_at %d after a typed consume, want 0", x.LastUsedAt)
+	}
+}
--- a/host/store/mint_requests.go
+++ b/host/store/mint_requests.go
@@ -305,6 +305,11 @@
 	if n != 1 {
 		return classifyConsume(ctx, tx, rc, at, n)
 	}
+	if rc.Approver != nil { // C8: the approver's row, in this transaction, or nothing is consumed
+		if err := useApprover(ctx, tx, *rc.Approver, at); err != nil {
+			return err
+		}
+	}
 	if err := insertSessionRow(ctx, tx, session); err != nil {
 		return err
 	}
````

### Patch M4e — authority binds a signature approval

````diff M4e
--- a/host/authority/mint.go
+++ b/host/authority/mint.go
@@ -57,12 +57,15 @@
 }
 
 // RequestMint is what a mint from a stored request needs beyond mintRecorded's arguments
-// (row 201 §2.6): the decoded request, its row's created_at and challenge, and the consume.
+// (row 201 §2.6): the decoded request, its row's created_at, challenge and expires_at, the
+// verified signature approval if there is one, and the consume.
 type RequestMint struct {
 	Request          approval.MintRequest
 	RequestCreatedAt int64
-	Challenge        [32]byte             // the stored row's challenge (row 201 A1b §2.6)
-	Consume          store.RequestConsume // ID, Digest, Clock; MintForRequest sets Challenge, ApprovedBy and Check
+	Challenge        [32]byte                 // the stored row's challenge (row 201 A1b §2.6)
+	ExpiresAt        int64                    // the stored row's expires_at: what a signature approval signed
+	Signed           *approval.SignedApproval // nil: the typed y; else the verified signature approval
+	Consume          store.RequestConsume     // ID, Digest, Clock; MintForRequest sets Challenge, ApprovedBy, Approver and Check
 }
 
 var (
@@ -74,10 +77,11 @@
 )
 
 // requestConsume is what the store binds for rm (row 201 A1b §2.6): the caller's consume, with the
-// challenge from rm and approved_by derived here.
+// challenge from rm and approved_by (and, for a signature, the approver) derived here from rm.Signed,
+// the same field that decides whether the terminal was skipped.
 func requestConsume(rm *RequestMint) (store.RequestConsume, error) {
 	consume := rm.Consume
-	if consume.ApprovedBy != "" { // A1
+	if consume.ApprovedBy != "" || consume.Approver != nil { // A1
 		return store.RequestConsume{}, ErrApprovedBySetByCaller
 	}
 	digest, err := hex.DecodeString(consume.Digest)
@@ -85,7 +89,16 @@
 		return store.RequestConsume{}, ErrBadConsumeDigest
 	}
 	consume.Challenge = hex.EncodeToString(rm.Challenge[:])
-	consume.ApprovedBy = "tty" // A5: the typed y
+	if rm.Signed == nil { // A5: the typed y
+		consume.ApprovedBy = "tty"
+		return consume, nil
+	}
+	if err := rm.Signed.Binds([32]byte(digest), rm.Challenge, rm.ExpiresAt); err != nil { // A4
+		return store.RequestConsume{}, err
+	}
+	name, format, key, _ := rm.Signed.Use() // A3: Binds has shown it is verified
+	consume.ApprovedBy = rm.Signed.ApprovedBy()
+	consume.Approver = &store.ApproverUse{Name: name, Format: format, PublicKey: key}
 	return consume, nil
 }
 
--- a/host/authority/mint_derive_test.go
+++ b/host/authority/mint_derive_test.go
@@ -5,12 +5,15 @@
 // neither.
 
 import (
+	"encoding/hex"
 	"errors"
 	"path/filepath"
+	"reflect"
 	"strings"
 	"testing"
 
 	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
 	"github.com/sunholo-data/ailang-world/host/store"
 )
 
@@ -88,4 +91,177 @@
 			t.Fatalf("request after the typed mint: %+v ok %v err %v, want consumed by tty", got, ok, err)
 		}
 	})
+}
+
+// ---- M4e ----
+
+// au1Signer is the approver key the AU tests sign with: approvaltest's fixed seed (D-WORLD-109).
+var au1Signer = approvaltest.NewSigner(7)
+
+// au1Enrol enrols "mark-mac" with au1Signer's key in st, at 900.
+func au1Enrol(t *testing.T, st *store.Credentials) {
+	t.Helper()
+	row := store.ApproverRow{Name: "mark-mac", Format: approval.FormatSSH, PublicKey: au1Signer.PublicKey().Blob, EnrolledAt: 900}
+	if err := st.EnrolApprover(boundedTestContext(t), row); err != nil {
+		t.Fatal(err)
+	}
+}
+
+// au1Signed is au1Signer's verified approval of (digest, challenge, expiry), admitted at 1010 for a
+// request created at 1000; it fails the test unless ApproveBySignature returns one.
+func au1Signed(t *testing.T, digestHex string, challenge [32]byte, expiry int64) *approval.SignedApproval {
+	t.Helper()
+	d, err := hex.DecodeString(digestHex)
+	if err != nil || len(d) != 32 {
+		t.Fatalf("digest %q: %v", digestHex, err)
+	}
+	digest := [32]byte(d)
+	msg, err := approval.SignedMessage(digest, challenge, expiry)
+	if err != nil {
+		t.Fatal(err)
+	}
+	s, err := approval.ApproveBySignature(approval.SignatureInput{
+		Admit:     approval.AdmitInput{Now: 1010, CreatedAt: 1000, ExpiresAt: expiry, StoredDigest: digest, BytesDigest: digest},
+		Challenge: challenge,
+		Signature: au1Signer.Sign(approval.Namespace, msg),
+		Approver:  approval.Approver{Name: "mark-mac", Format: approval.FormatSSH, Key: au1Signer.PublicKey(), EnrolledAt: 900},
+	})
+	if err != nil || s == nil {
+		t.Fatalf("ApproveBySignature: (%v, %v), want a verified approval", s, err)
+	}
+	return s
+}
+
+// au1SignedMint is au1RequestMint carrying au1Signer's approval of the row as stored.
+func au1SignedMint(t *testing.T, mr approval.MintRequest, row store.MintRequestRow) RequestMint {
+	t.Helper()
+	rm := au1RequestMint(mr, row)
+	rm.ExpiresAt = row.ExpiresAt
+	rm.Signed = au1Signed(t, row.Digest, au1Challenge, row.ExpiresAt)
+	return rm
+}
+
+// approverLastUsed is mark-mac's last_used_at in st.
+func approverLastUsed(t *testing.T, st *store.Credentials) int64 {
+	t.Helper()
+	a, ok, err := st.Approver(boundedTestContext(t), "mark-mac")
+	if err != nil || !ok {
+		t.Fatalf("approver mark-mac: ok %v err %v", ok, err)
+	}
+	return a.LastUsedAt
+}
+
+// AU1 (M4e): a verified signature approval records approved_by=ssh:<name> and marks the approver
+// used at the consume's reading (A3).
+func TestMintForRequestRecordsTheSigner(t *testing.T) {
+	st, mr, row := au1Store(t)
+	au1Enrol(t, st)
+	if err := au1Mint(t, st, au1SignedMint(t, mr, row)); err != nil {
+		t.Fatal(err)
+	}
+	got, ok, err := st.MintRequest(boundedTestContext(t), row.RequestID)
+	if err != nil || !ok || !got.Consumed || got.ApprovedBy != "ssh:mark-mac" {
+		t.Fatalf("request after the signed mint: %+v ok %v err %v, want consumed by ssh:mark-mac", got, ok, err)
+	}
+	if at := approverLastUsed(t, st); at != 1020 {
+		t.Fatalf("mark-mac's last_used_at %d, want 1020 (the consume's reading)", at)
+	}
+}
+
+// AU2 (M4e): an approval made for another digest, challenge or expiry is refused by Binds (A4), a
+// zero SignedApproval is ErrNotVerified, and a caller-set approver is refused (A1). Nothing is
+// minted and the approver is not marked used.
+func TestMintForRequestRefusesAnUnboundApproval(t *testing.T) {
+	otherDigest := strings.Repeat("cd", 32)
+	otherChallenge := au1Challenge
+	otherChallenge[31] ^= 1
+	bindingError := func(field string) func(error) bool {
+		return func(err error) bool {
+			var be *approval.BindingError
+			return errors.As(err, &be) && be.Field == field
+		}
+	}
+	is := func(target error) func(error) bool { return func(err error) bool { return errors.Is(err, target) } }
+	for _, a := range []struct {
+		name   string
+		mutate func(t *testing.T, rm *RequestMint, row store.MintRequestRow)
+		want   func(error) bool
+	}{
+		{"other-digest", func(t *testing.T, rm *RequestMint, row store.MintRequestRow) {
+			rm.Signed = au1Signed(t, otherDigest, au1Challenge, row.ExpiresAt)
+		}, bindingError("digest")},
+		{"other-challenge", func(t *testing.T, rm *RequestMint, row store.MintRequestRow) {
+			rm.Signed = au1Signed(t, row.Digest, otherChallenge, row.ExpiresAt)
+		}, bindingError("challenge")},
+		{"other-expiry", func(t *testing.T, rm *RequestMint, row store.MintRequestRow) {
+			rm.Signed = au1Signed(t, row.Digest, au1Challenge, row.ExpiresAt+1)
+		}, bindingError("expiry")},
+		{"expiry-not-passed", func(t *testing.T, rm *RequestMint, row store.MintRequestRow) { rm.ExpiresAt = 0 }, bindingError("expiry")},
+		{"zero-approval", func(t *testing.T, rm *RequestMint, row store.MintRequestRow) { rm.Signed = &approval.SignedApproval{} }, is(approval.ErrNotVerified)},
+		{"caller-approver", func(t *testing.T, rm *RequestMint, row store.MintRequestRow) {
+			rm.Signed = nil
+			rm.Consume.Approver = &store.ApproverUse{Name: "mark-mac", Format: approval.FormatSSH, PublicKey: au1Signer.PublicKey().Blob}
+		}, is(ErrApprovedBySetByCaller)},
+	} {
+		t.Run(a.name, func(t *testing.T) {
+			st, mr, row := au1Store(t)
+			au1Enrol(t, st)
+			rm := au1SignedMint(t, mr, row)
+			a.mutate(t, &rm, row)
+			if err := au1Mint(t, st, rm); !a.want(err) {
+				t.Fatalf("err %v, want the %s refusal", err, a.name)
+			}
+			assertNotMinted(t, st, row)
+			if at := approverLastUsed(t, st); at != 0 {
+				t.Fatalf("mark-mac's last_used_at %d after a refused mint, want 0", at)
+			}
+		})
+	}
+}
+
+// requestMintFields is where every RequestMint field comes from, and what binds it (RC5).
+var requestMintFields = map[string]string{
+	"Request":          "the decoded stored request; MintMatches checks the session row against it",
+	"RequestCreatedAt": "the stored row's created_at; MintMatches bounds the session's created_at by it",
+	"Challenge":        "the stored row's challenge (approve F4); bound in the consume's UPDATE (C6) and by Binds (A4)",
+	"ExpiresAt":        "the stored row's expires_at; bound by Binds (A4)",
+	"Signed":           "approval.ApproveBySignature's verified approval, or nil for the typed y; derives approved_by (A3, A5)",
+	"Consume":          "the row's id and digest and env.now; approved_by and the approver are derived (A1)",
+}
+
+// RC5 (M4e): every RequestMint field has a source, and perturbing any one of them alone in a valid
+// signed RequestMint refuses the mint.
+func TestRequestMintFieldsAreSourced(t *testing.T) {
+	typ := reflect.TypeOf(RequestMint{})
+	for i := 0; i < typ.NumField(); i++ {
+		if requestMintFields[typ.Field(i).Name] == "" {
+			t.Errorf("RequestMint.%s has no source in requestMintFields", typ.Field(i).Name)
+		}
+	}
+	if typ.NumField() != len(requestMintFields) {
+		t.Fatalf("RequestMint has %d fields, requestMintFields %d", typ.NumField(), len(requestMintFields))
+	}
+	perturb := map[string]func(*RequestMint){
+		"Request":          func(rm *RequestMint) { rm.Request.Episode = "other" },
+		"RequestCreatedAt": func(rm *RequestMint) { rm.RequestCreatedAt = 1011 },
+		"Challenge":        func(rm *RequestMint) { rm.Challenge[0] ^= 1 },
+		"ExpiresAt":        func(rm *RequestMint) { rm.ExpiresAt++ },
+		"Signed":           func(rm *RequestMint) { rm.Signed = &approval.SignedApproval{} },
+		"Consume":          func(rm *RequestMint) { rm.Consume.ID = "mr-ffffffffffffffff" },
+	}
+	if len(perturb) != len(requestMintFields) {
+		t.Fatalf("%d perturbations for %d fields", len(perturb), len(requestMintFields))
+	}
+	for field, p := range perturb {
+		t.Run(field, func(t *testing.T) {
+			st, mr, row := au1Store(t)
+			au1Enrol(t, st)
+			rm := au1SignedMint(t, mr, row)
+			p(&rm)
+			if err := au1Mint(t, st, rm); err == nil {
+				t.Fatalf("a mint with %s perturbed succeeded", field)
+			}
+			assertNotMinted(t, st, row)
+		})
+	}
 }
````

### Patch M5a — approve --signature (gated)

````diff M5a
--- a/cmd/ailang-worldd/approve.go
+++ b/cmd/ailang-worldd/approve.go
@@ -22,11 +22,14 @@
 )
 
 const approveHelp = `usage: ailang-worldd approve <request-id> [--db <path>] [--credentials-db <path>]
+       ailang-worldd approve <request-id> --signature <file> [--db <path>] [--credentials-db <path>]
        ailang-worldd approve --list [--db <path>] [--credentials-db <path>]
 
-ATTENDED. Mints the session a stored mint request describes (row 201), after a
-y/N typed at /dev/tty. 'ailang-worldd episode start <episode> ... --request'
-makes the request without a terminal; nothing is minted until approve.
+ATTENDED, or SIGNED. Mints the session a stored mint request describes (row
+201), after a y/N typed at /dev/tty, or, with --signature, after an enrolled
+approver's SSH signature over the request verifies, with no terminal.
+'ailang-worldd episode start <episode> ... --request' makes the request without
+a terminal; nothing is minted until approve.
 
 approve <request-id> (mr- and 16 hex):
   1. reads the request and refuses, before any terminal opens, one that is
@@ -40,6 +43,12 @@
      request expires. The request is consumed in the mint's own transaction, so
      it mints at most once;
   4. prints the launch line, as episode start does. The token is never printed.
+--signature <file> replaces step 3's y/N. The file (at most 16 KiB) is an
+  SSHSIG over the request's exact message, in namespace world-approval@sunholo,
+  by a key enrolled with 'ailang-worldd approver enroll' and not revoked. It is
+  checked before any network call, the approver is re-checked in the mint's own
+  transaction, and the request records approved_by=ssh:<name>. A signature
+  proves that something holding the key signed, not that you were there.
 --list prints the pending requests, read-only. A request expires 10 minutes
 after it is made; at most 16 are pending per credentials store.
 
@@ -98,6 +107,7 @@
 	db := fs.String("db", def.db, "the world store the request was made for (default ~/.ailang/world/world.db)")
 	credsDB := fs.String("credentials-db", "", credentialsDBUsage)
 	list := fs.Bool("list", false, "print the pending requests, read-only")
+	sigPath := fs.String("signature", "", "approve with an enrolled approver's SSHSIG over the request (no terminal)")
 	positional, err := parseInterleaved(fs, args)
 	if errors.Is(err, flag.ErrHelp) {
 		fmt.Fprint(stdout, approveHelp)
@@ -111,6 +121,10 @@
 		fmt.Fprint(stderr, approveHelp)
 		return exitUsage
 	}
+	if *list && *sigPath != "" { // F2: --list is read-only
+		fmt.Fprint(stderr, approveHelp)
+		return exitUsage
+	}
 	id := ""
 	if !*list {
 		id = positional[0]
@@ -168,6 +182,13 @@
 	if err != nil || len(challenge) != 32 { // F4: the consume binds the row's challenge (row 201 A1b §2.7)
 		return refuse("request %s's stored challenge is not 64 hex; refusing; fix: make a new request with episode start --request", row.RequestID)
 	}
+	var signed *approval.SignedApproval
+	if *sigPath != "" { // the signature, judged before any network call (design §2.7 F9-F14)
+		var refusal string
+		if signed, refusal = approveSignature(ctx, credPath, *sigPath, row, now); refusal != "" {
+			return refuse("%s", refusal)
+		}
+	}
 
 	// 2. The same store: the request was made for this world store and credentials file.
 	w, werr := store.CanonicalPath(*db)
@@ -210,12 +231,13 @@
 		return refuse("the world changed since request %s was made: %s differs; it stays pending", id, approval.Diff(mr, fresh))
 	}
 
-	// 6. The mint: episode start's own path, bounded by the request's expiry.
+	// 6. The mint: episode start's own path, bounded by the request's expiry. With a verified
+	// signature, provisionEpisode skips the terminal (design §2.8).
 	req.presetPrompt = renderRequest(row.RequestID, mr, row.Digest, row.ExpiresAt, now)
 	req.afterToken = mcpAfterToken(mr.MCPConfig, mr.Daemon.Addr)
 	req.answerBy = row.ExpiresAt
 	req.consume = &authority.RequestMint{Request: mr, RequestCreatedAt: row.CreatedAt, Challenge: [32]byte(challenge),
-		Consume: store.RequestConsume{ID: row.RequestID, Digest: row.Digest, Clock: env.now}}
+		ExpiresAt: row.ExpiresAt, Signed: signed, Consume: store.RequestConsume{ID: row.RequestID, Digest: row.Digest, Clock: env.now}}
 	if code := provisionEpisode(req, stdout, stderr, env); code != exitOK {
 		return code
 	}
--- a/cmd/ailang-worldd/approve_sig.go
+++ b/cmd/ailang-worldd/approve_sig.go
@@ -0,0 +1,62 @@
+package main
+
+// Row 201 A1b (design_docs/planned/w-world-approvals-by-key-a1b.md §2.7; sprint-3 plan §3 M5a):
+// `approve <request-id> --signature <file>` approves a stored mint request with an enrolled
+// approver's SSHSIG over the request's exact message, with no terminal.
+
+import (
+	"context"
+	"crypto/sha256"
+	"encoding/hex"
+	"fmt"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/store"
+)
+
+// approveMaxSignatureBytes bounds the --signature file (design §2.7 F9): an Ed25519 or P-256
+// SSHSIG is under 1 KiB.
+const approveMaxSignatureBytes = 16 << 10
+
+// approveSignature is approve --signature's check, before any network call (design §2.7 F9-F14):
+// the file is read once through a descriptor check (F9, F10), parsed (F11), its key looked up
+// among the enrolled approvers on a read-only handle under approve's 3 s root (F12, F13), and
+// judged by the signature law against the stored row (F14). It returns the verified approval, or
+// the refusal text. The caller has already admitted row at now and decoded its challenge (F4).
+func approveSignature(ctx context.Context, credPath, path string, row store.MintRequestRow, now int64) (*approval.SignedApproval, string) {
+	// F9: one read, through a check of the open descriptor (a FIFO or a device is refused at once).
+	armored, err := readRegularFileCapped(path, approveMaxSignatureBytes)
+	if err != nil { // F10
+		return nil, fmt.Sprintf("--signature: %v", err)
+	}
+	sig, err := approval.ParseSSHSIG(armored)
+	if err != nil { // F11
+		return nil, fmt.Sprintf("--signature %s is not an SSH signature: %v", path, err)
+	}
+	ro, err := store.OpenCredentialsReadOnly(credPath)
+	if err != nil {
+		return nil, fmt.Sprintf("cannot read the credentials store %s: %v", credPath, err)
+	}
+	enrolled, ok, err := ro.ApproverByKey(ctx, sig.PublicKey.Blob)
+	_ = ro.Close()
+	if err != nil { // F12: below schema 5 names the serve restart
+		return nil, err.Error()
+	}
+	if !ok { // F13
+		return nil, fmt.Sprintf("the signature's key %s is not an enrolled approver; see ailang-worldd approver list", sig.PublicKey.Fingerprint())
+	}
+	stored, _ := hex.DecodeString(row.Digest)       // admission has shown it is 32 bytes
+	challenge, _ := hex.DecodeString(row.Challenge) // F4 has shown it is 32 bytes
+	signed, err := approval.ApproveBySignature(approval.SignatureInput{
+		Admit: approval.AdmitInput{Now: now, CreatedAt: row.CreatedAt, ExpiresAt: row.ExpiresAt, Consumed: row.Consumed,
+			StoredDigest: [32]byte(stored), BytesDigest: sha256.Sum256(row.Request)},
+		Challenge: [32]byte(challenge),
+		Signature: armored,
+		Approver: approval.Approver{Name: enrolled.Name, Format: enrolled.Format, Key: approval.SSHPublicKey{Blob: enrolled.PublicKey},
+			EnrolledAt: enrolled.EnrolledAt, LastUsedAt: enrolled.LastUsedAt, RevokedAt: enrolled.RevokedAt},
+	})
+	if err != nil { // F14
+		return nil, fmt.Sprintf("the signature does not approve request %s: %v", row.RequestID, err)
+	}
+	return signed, ""
+}
--- a/cmd/ailang-worldd/approve_sig_test.go
+++ b/cmd/ailang-worldd/approve_sig_test.go
@@ -7,8 +7,14 @@
 
 import (
 	"bytes"
+	"encoding/hex"
+	"os"
+	"path/filepath"
 	"strings"
 	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
 )
 
 // CR11 (M4c2): a stored challenge that is not 64 hex is refused before any terminal opens (F4), and
@@ -32,5 +38,93 @@
 	}
 	if a.row(t).Consumed || len(a.sessionRows(t)) != 0 {
 		t.Fatal("a request with a non-hex challenge was consumed or minted")
+	}
+}
+
+// ---- M5a ----
+
+// sigSigner is the enrolled approver's key in the approve --signature tests; otherSigner is a key
+// that is never enrolled. Both are approvaltest's fixed seeds (D-WORLD-109).
+var (
+	sigSigner   = approvaltest.NewSigner(11)
+	otherSigner = approvaltest.NewSigner(12)
+)
+
+// requestMessage is the 90-byte message the approval of the rig's stored row signs.
+func requestMessage(t *testing.T, a approveRig) []byte {
+	t.Helper()
+	row := a.row(t)
+	d, err := hex.DecodeString(row.Digest)
+	if err != nil {
+		t.Fatal(err)
+	}
+	c, err := hex.DecodeString(row.Challenge)
+	if err != nil {
+		t.Fatal(err)
+	}
+	msg, err := approval.SignedMessage([32]byte(d), [32]byte(c), row.ExpiresAt)
+	if err != nil {
+		t.Fatal(err)
+	}
+	return msg
+}
+
+// writeSignature writes sig to a new file under t.TempDir() and returns its path.
+func writeSignature(t *testing.T, sig []byte) string {
+	t.Helper()
+	p := filepath.Join(t.TempDir(), "request.sig")
+	if err := os.WriteFile(p, sig, 0o600); err != nil {
+		t.Fatal(err)
+	}
+	return p
+}
+
+// signedRig is a rig for ep1 at approveT0 with sigSigner enrolled as mark-mac at approveT0-100, and
+// the path of sigSigner's signature over the request's message.
+func signedRig(t *testing.T) (approveRig, string) {
+	t.Helper()
+	a := newApproveRig(t, "ep1", approveT0)
+	seedApprover(t, a.credPath, "mark-mac", sigSigner.PublicKey(), approveT0-100)
+	return a, writeSignature(t, sigSigner.Sign(approval.Namespace, requestMessage(t, a)))
+}
+
+// CR2 (M5a): a verified signature mints the stored request with no terminal: approved_by is
+// ssh:mark-mac, the approver's last_used_at is the consume's reading, stdout shows the request and
+// its approver, and every other effect is what the typed y produces (CL16).
+func TestApproveBySignatureMintsWithNoTerminal(t *testing.T) {
+	a, sig := signedRig(t)
+	var opened bool
+	var stdout, stderr bytes.Buffer
+	if code := runApprove(a.args("--signature", sig), &stdout, &stderr, approveEnv(nil, &opened, atTime(approveT0+1))); code != exitOK {
+		t.Fatalf("rc %d, want exitOK\nstderr %s", code, stderr.String())
+	}
+	if opened {
+		t.Fatal("a terminal was opened on the signature path")
+	}
+	row := a.row(t)
+	if !row.Consumed || row.ApprovedBy != "ssh:mark-mac" {
+		t.Fatalf("request after the signed approve: consumed %v by %q, want ssh:mark-mac", row.Consumed, row.ApprovedBy)
+	}
+	if rows := approverRows(t, a.credPath); len(rows) != 1 || rows[0].LastUsedAt != approveT0+1 {
+		t.Fatalf("approver rows %+v, want mark-mac last used at %d", rows, approveT0+1)
+	}
+	out := stdout.String()
+	for _, want := range []string{"mint request " + a.id + "\n", "expiry +3600s: approved by ssh:mark-mac\n"} {
+		if !strings.Contains(out, want) {
+			t.Errorf("stdout lacks %q:\n%s", want, out)
+		}
+	}
+	mr, err := approval.Decode(row.Request)
+	if err != nil {
+		t.Fatal(err)
+	}
+	var probs []string
+	for _, p := range approveEffectProblems(mr, collectApproveFacts(t, a, mr, out)) {
+		if !strings.HasPrefix(p, "approved_by:") { // the typed check wants tty; this path records ssh:mark-mac (asserted above)
+			probs = append(probs, p)
+		}
+	}
+	if len(probs) != 0 {
+		t.Errorf("the effects differ from the request:\n%s", strings.Join(probs, "\n"))
 	}
 }
--- a/cmd/ailang-worldd/approver.go
+++ b/cmd/ailang-worldd/approver.go
@@ -120,10 +120,10 @@
        ailang-worldd approver revoke <name> [--db <path>] [--credentials-db <path>]
 
 ATTENDED (enroll, revoke). An approver is a name and one SSH public key, kept in
-the credentials store (row 201). Approving a mint request with its signature
-arrives with 'approve --signature' (row 201 A1b, a later sprint); until then an
-enrolled key approves nothing. enroll and revoke each ask y at /dev/tty and
-wait at most 2m0s for it; an agent cannot run them.
+the credentials store (row 201). An enrolled, unrevoked key approves a mint
+request with its signature: 'ailang-worldd approve <request-id> --signature
+<file>'. enroll and revoke each ask y at /dev/tty and wait at most 2m0s for it;
+an agent cannot run them.
 
 enroll   reads --pubkey once (one ssh-ed25519 or ecdsa-sha2-nistp256 key line,
          at most 8 KiB; a private key file is refused and never shown), shows
--- a/cmd/ailang-worldd/help_test.go
+++ b/cmd/ailang-worldd/help_test.go
@@ -41,7 +41,7 @@
 		"log tail", "registry get", "commit --file <commit.json> [--session <file|token>]",
 		"tools list", "] call <tool>", "] why ", "provenance", "setup", "doctor",
 		"session new <episode>", "session list", "session revoke", "session mint", "help [<verb>]",
-		"approve <request-id>", "approve --list", "approver enroll --name", "approver list", "approver revoke <name>",
+		"approve <request-id>", "approve <request-id> --signature <file>", "approve --list", "approver enroll --name", "approver list", "approver revoke <name>",
 	} {
 		if !strings.Contains(out.String(), want) {
 			t.Errorf("top-level usage lacks %q", want)
--- a/cmd/ailang-worldd/main.go
+++ b/cmd/ailang-worldd/main.go
@@ -122,6 +122,7 @@
                     [--db <path>] [--credentials-db <path>] [--workspace-root <dir>]
                     [--sessions-dir <dir>] [--exec-profiles-dir <dir>] [--addr 127.0.0.1:7644 (or service.json's address)] [--request]
   ailang-worldd approve <request-id> [--db <path>] [--credentials-db <path>]
+  ailang-worldd approve <request-id> --signature <file> [--db <path>] [--credentials-db <path>]
   ailang-worldd approve --list [--db <path>] [--credentials-db <path>]
   ailang-worldd approver enroll --name <name> --format ssh --pubkey <file> [--db <path>] [--credentials-db <path>]
   ailang-worldd approver list [--db <path>] [--credentials-db <path>]
--- a/cmd/ailang-worldd/mint_request.go
+++ b/cmd/ailang-worldd/mint_request.go
@@ -220,7 +220,7 @@
 	"startPoint":         "fixed: \"\" (set only by episode_preset.go)",
 	"existingBranch":     "fixed: false (set only by episode_preset.go)",
 	"presetPrompt":       "excluded: display only; approve sets renderRequest(request)",
-	"consume":            "excluded: set by approve from the stored row (id, digest, challenge, env.now); authority.MintForRequest derives approved_by",
+	"consume":            "excluded: set by approve from the stored row (id, digest, challenge, expires_at, env.now) and, with --signature, the verified approval; authority.MintForRequest derives approved_by",
 	"answerBy":           "excluded: set by approve to the stored row's expires_at",
 }
 
--- a/cmd/ailang-worldd/session_new.go
+++ b/cmd/ailang-worldd/session_new.go
@@ -25,7 +25,8 @@
 // start` (episode.go, row 189 M3), its launch-preset and resume paths, and
 // `approve` (approve.go, row 201, which consumes a stored mint request in the
 // mint's own transaction) run it too, after their own validation, so the verbs
-// cannot drift apart and no other code creates a credential.
+// cannot drift apart and no other code creates a credential. `approve --signature` runs steps
+// 4-7 without steps 2-3's terminal, and only with a verified SSHSIG approval (row 201 A1b).
 //
 // Row 189 M2 (D-189-2 A): --profile <project> binds the credential to the
 // exec profile <exec-profiles-dir>/<project>.json. The file is parsed (never
@@ -402,14 +403,19 @@
 // new and episode start. It returns the exit code.
 func provisionEpisode(r *mintRequest, stdout, stderr io.Writer, env sessionEnv) int {
 	prefix := "ailang-worldd " + r.verb + ": "
-	// 2. The mint fence, unchanged, before any side effect.
-	term, err := env.openTerminal()
-	if err != nil {
-		fmt.Fprintf(stderr, "%srefusing: no controlling terminal (%v); "+
-			"provisioning a session credential requires one human act at a terminal. %s\n", prefix, err, sessionTTYFix)
-		return exitUsage
+	// 2. The mint fence, unchanged, before any side effect, unless a verified signature approved
+	// this request (row 201 A1b §2.8: Signed is set only by approve, from ApproveBySignature).
+	signed := r.consume != nil && r.consume.Signed != nil
+	var term io.ReadWriteCloser
+	if !signed { // V1
+		var err error
+		if term, err = env.openTerminal(); err != nil {
+			fmt.Fprintf(stderr, "%srefusing: no controlling terminal (%v); "+
+				"provisioning a session credential requires one human act at a terminal. %s\n", prefix, err, sessionTTYFix)
+			return exitUsage
+		}
+		defer func() { _ = term.Close() }()
 	}
-	defer func() { _ = term.Close() }()
 
 	// 3. y/N naming the worktree, every grant and the exec profile binding.
 	now := env.now()
@@ -442,30 +448,36 @@
 		fmt.Fprintf(stderr, "%srequest %s expired before the prompt; make a new request with episode start --request\n", prefix, requestID)
 		return exitUsage
 	}
-	if r.presetPrompt != "" {
-		fmt.Fprint(term, r.presetPrompt)
-	}
 	if r.existingBranch {
 		action += " on existing branch " + r.branch
 	}
-	fmt.Fprintf(term, "%s and mint a session for episode %s with %d grant(s) [%s]%s, expiry +%ds? [y/N] ",
-		action, r.episode, len(grants), grantList(grants), bound, r.ttl)
-	var answer string
-	if r.answerBy != 0 {
-		remaining := time.Duration(r.answerBy-now) * time.Second
-		line, ok := readTermLineWithin(term, remaining, env.after)
-		if !ok {
-			fmt.Fprintf(stderr, "%sno answer within the request's remaining lifetime (%s); nothing was minted and request %s is still pending until %s\n",
-				prefix, remaining, requestID, time.Unix(r.answerBy, 0).UTC().Format(time.RFC3339))
+	if signed { // V2: no prompt; the request and its approver go to stdout
+		fmt.Fprint(stdout, r.presetPrompt)
+		fmt.Fprintf(stdout, "%s and mint a session for episode %s with %d grant(s) [%s]%s, expiry +%ds: approved by %s\n",
+			action, r.episode, len(grants), grantList(grants), bound, r.ttl, r.consume.Signed.ApprovedBy())
+	} else {
+		if r.presetPrompt != "" {
+			fmt.Fprint(term, r.presetPrompt)
+		}
+		fmt.Fprintf(term, "%s and mint a session for episode %s with %d grant(s) [%s]%s, expiry +%ds? [y/N] ",
+			action, r.episode, len(grants), grantList(grants), bound, r.ttl)
+		var answer string
+		if r.answerBy != 0 {
+			remaining := time.Duration(r.answerBy-now) * time.Second
+			line, ok := readTermLineWithin(term, remaining, env.after)
+			if !ok {
+				fmt.Fprintf(stderr, "%sno answer within the request's remaining lifetime (%s); nothing was minted and request %s is still pending until %s\n",
+					prefix, remaining, requestID, time.Unix(r.answerBy, 0).UTC().Format(time.RFC3339))
+				return exitUsage
+			}
+			answer = line
+		} else {
+			answer = readTermLine(term)
+		}
+		if !isYes(answer) {
+			fmt.Fprintf(stderr, "%saborted (not confirmed). fix: run it again and type y at the prompt\n", prefix)
 			return exitUsage
 		}
-		answer = line
-	} else {
-		answer = readTermLine(term)
-	}
-	if !isYes(answer) {
-		fmt.Fprintf(stderr, "%saborted (not confirmed). fix: run it again and type y at the prompt\n", prefix)
-		return exitUsage
 	}
 
 	// 4. The credentials writer lock, before the worktree exists.
````

### Patch M5b — refusals and the revoke race (gated)

````diff M5b
--- a/cmd/ailang-worldd/approve_sig_test.go
+++ b/cmd/ailang-worldd/approve_sig_test.go
@@ -8,13 +8,16 @@
 import (
 	"bytes"
 	"encoding/hex"
+	"errors"
 	"os"
 	"path/filepath"
 	"strings"
+	"syscall"
 	"testing"
 
 	"github.com/sunholo-data/ailang-world/host/approval"
 	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
+	"github.com/sunholo-data/ailang-world/host/daemon"
 )
 
 // CR11 (M4c2): a stored challenge that is not 64 hex is refused before any terminal opens (F4), and
@@ -126,5 +129,177 @@
 	}
 	if len(probs) != 0 {
 		t.Errorf("the effects differ from the request:\n%s", strings.Join(probs, "\n"))
+	}
+}
+
+// ---- M5b ----
+
+// CR3 (M5b): with an approver enrolled, the typed path still opens the terminal, records tty and
+// leaves the approver unused; --list refuses --signature (F2).
+func TestApproveTypedPathUnchangedWithAnApprover(t *testing.T) {
+	a, sig := signedRig(t)
+	var opened bool
+	var stdout, stderr bytes.Buffer
+	term := &promptTerm{in: strings.NewReader("y\n")}
+	if code := runApprove(a.args(), &stdout, &stderr, approveEnv(term, &opened, atTime(approveT0+1))); code != exitOK || !opened {
+		t.Fatalf("typed approve: rc %d, terminal opened %v; want exitOK through the terminal\nstderr %s", code, opened, stderr.String())
+	}
+	if row := a.row(t); row.ApprovedBy != "tty" {
+		t.Fatalf("approved_by %q, want tty", row.ApprovedBy)
+	}
+	if rows := approverRows(t, a.credPath); len(rows) != 1 || rows[0].LastUsedAt != 0 {
+		t.Fatalf("approver rows %+v, want mark-mac never used", rows)
+	}
+	stderr.Reset()
+	listArgs := []string{"--list", "--signature", sig, "--db", a.r.f.db, "--credentials-db", a.credPath}
+	if code := runApprove(listArgs, &stdout, &stderr, approveEnv(nil, nil, atTime(approveT0+2))); code != exitUsage || stderr.String() != approveHelp {
+		t.Fatalf("--list --signature: rc %d, stderr %q; want exitUsage and the help", code, stderr.String())
+	}
+}
+
+// CR4, CR5 (M5b): each signature refusal exits exitUsage with its text, before any terminal and
+// before any network call, and leaves the request as it was. A consumed or expired request is
+// refused by admission before the signature file is read (the path does not exist).
+func TestApproveBySignatureRefusalsBeforeAnyEffect(t *testing.T) {
+	good := func(t *testing.T, a approveRig) string {
+		return writeSignature(t, sigSigner.Sign(approval.Namespace, requestMessage(t, a)))
+	}
+	absent := func(t *testing.T, a approveRig) string { return filepath.Join(t.TempDir(), "absent.sig") }
+	for _, arm := range []struct {
+		name     string
+		setup    func(t *testing.T, a approveRig)
+		sig      func(t *testing.T, a approveRig) string
+		now      int64 // 0: approveT0+1
+		want     string
+		consumed bool // the setup consumed the request (one session row)
+	}{
+		{name: "fifo", sig: func(t *testing.T, a approveRig) string {
+			p := filepath.Join(t.TempDir(), "sig.fifo")
+			if err := syscall.Mkfifo(p, 0o600); err != nil {
+				t.Fatal(err)
+			}
+			return p
+		}, want: "sig.fifo is not a regular file"},
+		{name: "over-16k", sig: func(t *testing.T, a approveRig) string {
+			return writeSignature(t, bytes.Repeat([]byte("a"), 16<<10+1))
+		}, want: "is 16385 bytes, over the 16384-byte bound"},
+		{name: "missing", sig: absent, want: "--signature: open "},
+		{name: "garbage", sig: func(t *testing.T, a approveRig) string { return writeSignature(t, []byte("not a signature\n")) },
+			want: "is not an SSH signature: approval: sshsig: bad armor"},
+		{name: "unenrolled-key", sig: func(t *testing.T, a approveRig) string {
+			return writeSignature(t, otherSigner.Sign(approval.Namespace, requestMessage(t, a)))
+		}, want: "the signature's key " + otherSigner.PublicKey().Fingerprint() + " is not an enrolled approver; see ailang-worldd approver list"},
+		{name: "revoked", setup: func(t *testing.T, a approveRig) { revokeSeeded(t, a.credPath, "mark-mac", approveT0-50) },
+			sig: good, want: "approval: the approver is revoked"},
+		{name: "other-request", sig: func(t *testing.T, a approveRig) string {
+			msg := requestMessage(t, a)
+			msg[len(msg)-9] ^= 1 // the challenge's last byte
+			return writeSignature(t, sigSigner.Sign(approval.Namespace, msg))
+		}, want: "approval: sshsig: signature does not verify"},
+		{name: "enrolled-after-now", setup: func(t *testing.T, a approveRig) {
+			raw := rawCredentials(t, a.credPath)
+			defer func() { _ = raw.Close() }()
+			if _, err := raw.Exec(`UPDATE approvers SET enrolled_at = ? WHERE name = 'mark-mac'`, approveT0+50); err != nil {
+				t.Fatal(err)
+			}
+		}, sig: good, want: "approval: the clock is behind the approver's enrolment time"},
+		{name: "git-namespace", sig: func(t *testing.T, a approveRig) string {
+			return writeSignature(t, sigSigner.Sign("git", requestMessage(t, a)))
+		}, want: "approval: sshsig: wrong namespace"},
+		{name: "v4-store", setup: func(t *testing.T, a approveRig) {
+			raw := rawCredentials(t, a.credPath)
+			defer func() { _ = raw.Close() }()
+			for _, stmt := range []string{`DROP TABLE approvers`, `PRAGMA user_version = 4`} {
+				if _, err := raw.Exec(stmt); err != nil {
+					t.Fatal(err)
+				}
+			}
+		}, sig: good, want: "approvers need schema 5"},
+		{name: "consumed-first", setup: func(t *testing.T, a approveRig) {
+			var stdout, stderr bytes.Buffer
+			if code := runApprove(a.args(), &stdout, &stderr, approveEnv(&promptTerm{in: strings.NewReader("y\n")}, nil, atTime(approveT0+1))); code != exitOK {
+				t.Fatalf("typed approve: rc %d\nstderr %s", code, stderr.String())
+			}
+		}, sig: absent, want: "was already used", consumed: true},
+		{name: "expired-first", sig: absent, now: approveT0 + 600, want: "expired at"},
+	} {
+		t.Run(arm.name, func(t *testing.T) {
+			a := newApproveRig(t, "ep1", approveT0)
+			seedApprover(t, a.credPath, "mark-mac", sigSigner.PublicKey(), approveT0-100)
+			if arm.setup != nil {
+				arm.setup(t, a)
+			}
+			calls := 0
+			inner := fetchDaemonHealth
+			t.Cleanup(func() { fetchDaemonHealth = inner })
+			fetchDaemonHealth = func(addr string) (daemon.HealthResponse, error) { calls++; return inner(addr) }
+			now := arm.now
+			if now == 0 {
+				now = approveT0 + 1
+			}
+			var opened bool
+			var stdout, stderr bytes.Buffer
+			code := runApprove(a.args("--signature", arm.sig(t, a)), &stdout, &stderr, approveEnv(nil, &opened, atTime(now)))
+			if code != exitUsage || opened || calls != 0 {
+				t.Fatalf("rc %d, terminal opened %v, daemon calls %d; want exitUsage, none, 0\nstderr %s", code, opened, calls, stderr.String())
+			}
+			if !strings.Contains(stderr.String(), arm.want) {
+				t.Fatalf("stderr %q lacks %q", stderr.String(), arm.want)
+			}
+			if row := a.row(t); row.Consumed != arm.consumed || row.ApprovedBy == "ssh:mark-mac" {
+				t.Fatalf("request consumed %v by %q, want consumed %v and never by ssh:mark-mac", row.Consumed, row.ApprovedBy, arm.consumed)
+			}
+			want := 0
+			if arm.consumed {
+				want = 1
+			}
+			if n := len(a.sessionRows(t)); n != want {
+				t.Fatalf("%d session rows, want %d", n, want)
+			}
+		})
+	}
+}
+
+// CR6 (M5b): a revoke that commits after approve verified the signature, inside the mint's own
+// transaction (the consume's in-transaction clock reading: design V36/V37), refuses the mint at
+// the write: exitFatal, the request pending, no session, the worktree this run made removed.
+func TestApproveRevokeRaceIsRefusedAtTheWrite(t *testing.T) {
+	a, sig := signedRig(t)
+	calls := 0
+	clock := func() int64 {
+		calls++
+		if calls == 3 { // 1 admission, 2 provisionEpisode, 3 the consume's reading
+			raw := rawCredentials(t, a.credPath)
+			defer func() { _ = raw.Close() }()
+			if _, err := raw.Exec(`UPDATE approvers SET revoked_at = ? WHERE name = 'mark-mac'`, approveT0+1); err != nil {
+				t.Errorf("the racing revoke: %v", err)
+			}
+		}
+		return approveT0 + 1
+	}
+	var opened bool
+	var stdout, stderr bytes.Buffer
+	if code := runApprove(a.args("--signature", sig), &stdout, &stderr, approveEnv(nil, &opened, clock)); code != exitFatal {
+		t.Fatalf("rc %d, want exitFatal\nstderr %s", code, stderr.String())
+	}
+	for _, want := range []string{"store: the approver is revoked", "removed the worktree"} {
+		if !strings.Contains(stderr.String(), want) {
+			t.Errorf("stderr lacks %q:\n%s", want, stderr.String())
+		}
+	}
+	if calls != 3 || opened {
+		t.Errorf("clock read %d times (want 3), terminal opened %v", calls, opened)
+	}
+	if row := a.row(t); row.Consumed {
+		t.Error("the request was consumed by a revoked approver")
+	}
+	if n := len(a.sessionRows(t)); n != 0 {
+		t.Errorf("%d session rows, want 0", n)
+	}
+	if _, err := os.Lstat(a.worktree()); !errors.Is(err, os.ErrNotExist) {
+		t.Errorf("the worktree exists (err %v)", err)
+	}
+	if rows := approverRows(t, a.credPath); len(rows) != 1 || rows[0].LastUsedAt != 0 || rows[0].RevokedAt != approveT0+1 {
+		t.Errorf("approver rows %+v, want mark-mac revoked and never used", rows)
 	}
 }
````

### Patch M5c — approve --write-message and the flag rules (gated)

````diff M5c
--- a/cmd/ailang-worldd/approve.go
+++ b/cmd/ailang-worldd/approve.go
@@ -23,6 +23,7 @@
 
 const approveHelp = `usage: ailang-worldd approve <request-id> [--db <path>] [--credentials-db <path>]
        ailang-worldd approve <request-id> --signature <file> [--db <path>] [--credentials-db <path>]
+       ailang-worldd approve <request-id> --write-message <file> [--db <path>] [--credentials-db <path>]
        ailang-worldd approve --list [--db <path>] [--credentials-db <path>]
 
 ATTENDED, or SIGNED. Mints the session a stored mint request describes (row
@@ -49,6 +50,9 @@
   checked before any network call, the approver is re-checked in the mint's own
   transaction, and the request records approved_by=ssh:<name>. A signature
   proves that something holding the key signed, not that you were there.
+--write-message <file> writes that message (90 bytes) to a NEW file and prints
+  the ssh-keygen -Y sign line and the --signature step. It opens no terminal
+  and writes no store row: the request stays pending.
 --list prints the pending requests, read-only. A request expires 10 minutes
 after it is made; at most 16 are pending per credentials store.
 
@@ -108,6 +112,7 @@
 	credsDB := fs.String("credentials-db", "", credentialsDBUsage)
 	list := fs.Bool("list", false, "print the pending requests, read-only")
 	sigPath := fs.String("signature", "", "approve with an enrolled approver's SSHSIG over the request (no terminal)")
+	msgPath := fs.String("write-message", "", "write the request's approval message to this new file, to sign elsewhere")
 	positional, err := parseInterleaved(fs, args)
 	if errors.Is(err, flag.ErrHelp) {
 		fmt.Fprint(stdout, approveHelp)
@@ -121,10 +126,14 @@
 		fmt.Fprint(stderr, approveHelp)
 		return exitUsage
 	}
-	if *list && *sigPath != "" { // F2: --list is read-only
+	if *sigPath != "" && *msgPath != "" { // F1: one of --signature and --write-message
 		fmt.Fprint(stderr, approveHelp)
 		return exitUsage
 	}
+	if *list && (*sigPath != "" || *msgPath != "") { // F2: --list is read-only
+		fmt.Fprint(stderr, approveHelp)
+		return exitUsage
+	}
 	id := ""
 	if !*list {
 		id = positional[0]
@@ -181,6 +190,12 @@
 	challenge, err := hex.DecodeString(row.Challenge)
 	if err != nil || len(challenge) != 32 { // F4: the consume binds the row's challenge (row 201 A1b §2.7)
 		return refuse("request %s's stored challenge is not 64 hex; refusing; fix: make a new request with episode start --request", row.RequestID)
+	}
+	if *msgPath != "" { // F5, F6: the message to sign elsewhere, and nothing else
+		if refusal := writeApprovalMessage(stdout, *msgPath, row, mr); refusal != "" {
+			return refuse("%s", refusal)
+		}
+		return exitOK
 	}
 	var signed *approval.SignedApproval
 	if *sigPath != "" { // the signature, judged before any network call (design §2.7 F9-F14)
--- a/cmd/ailang-worldd/approve_sig.go
+++ b/cmd/ailang-worldd/approve_sig.go
@@ -9,6 +9,8 @@
 	"crypto/sha256"
 	"encoding/hex"
 	"fmt"
+	"io"
+	"os"
 
 	"github.com/sunholo-data/ailang-world/host/approval"
 	"github.com/sunholo-data/ailang-world/host/store"
@@ -59,4 +61,37 @@
 		return nil, fmt.Sprintf("the signature does not approve request %s: %v", row.RequestID, err)
 	}
 	return signed, ""
+}
+
+// ---- M5c ----
+
+// writeApprovalMessage is approve --write-message (design §2.7 F5, F6): it writes the admitted
+// request's 90-byte approval message to a NEW file (O_EXCL, mode 0644: the message is not secret)
+// and prints the out-of-band signing line and the next step. It opens no terminal and writes no
+// store row, so the request stays pending. It returns the refusal text, or "".
+func writeApprovalMessage(stdout io.Writer, path string, row store.MintRequestRow, mr approval.MintRequest) string {
+	stored, _ := hex.DecodeString(row.Digest)       // admission has shown it is 32 bytes
+	challenge, _ := hex.DecodeString(row.Challenge) // F4 has shown it is 32 bytes
+	msg, err := approval.SignedMessage([32]byte(stored), [32]byte(challenge), row.ExpiresAt)
+	if err != nil {
+		return fmt.Sprintf("request %s: %v", row.RequestID, err)
+	}
+	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
+	if err != nil { // F5: an existing file is never overwritten
+		return fmt.Sprintf("--write-message: %v; nothing was written", err)
+	}
+	_, err = f.Write(msg)
+	if cerr := f.Close(); err == nil {
+		err = cerr
+	}
+	if err != nil {
+		_ = os.Remove(path)
+		return fmt.Sprintf("--write-message: %v; nothing was written", err)
+	}
+	fmt.Fprintf(stdout, "wrote request %s's %d-byte approval message to %s\n", row.RequestID, len(msg), path)
+	fmt.Fprintf(stdout, "sign it with an enrolled key (the request expires at %s):\n  ssh-keygen -Y sign -n %s -f <your key> < %s > %s.sig\n",
+		approverTime(row.ExpiresAt), approval.Namespace, shellQuote(path), shellQuote(path))
+	fmt.Fprintf(stdout, "then approve it:\n  ailang-worldd approve %s --signature %s.sig --db %s --credentials-db %s\n",
+		row.RequestID, shellQuote(path), shellQuote(mr.WorldDB), shellQuote(mr.CredentialsDB))
+	return ""
 }
--- a/cmd/ailang-worldd/approve_sig_test.go
+++ b/cmd/ailang-worldd/approve_sig_test.go
@@ -303,3 +303,110 @@
 		t.Errorf("approver rows %+v, want mark-mac revoked and never used", rows)
 	}
 }
+
+// ---- M5c ----
+
+// CR10 (M5c): --write-message writes exactly the request's 90-byte message to a new file and the
+// two next steps, opens no terminal, calls no daemon and leaves the request pending (F6); signing
+// that file approves the request. An existing file is refused and kept (F5); a used request is
+// refused by admission and nothing is written.
+func TestApproveWriteMessage(t *testing.T) {
+	t.Run("writes-then-signs", func(t *testing.T) {
+		a, _ := signedRig(t)
+		calls := 0
+		inner := fetchDaemonHealth
+		t.Cleanup(func() { fetchDaemonHealth = inner })
+		fetchDaemonHealth = func(addr string) (daemon.HealthResponse, error) { calls++; return inner(addr) }
+		path := filepath.Join(t.TempDir(), "request.msg")
+		var opened bool
+		var stdout, stderr bytes.Buffer
+		if code := runApprove(a.args("--write-message", path), &stdout, &stderr, approveEnv(nil, &opened, atTime(approveT0+1))); code != exitOK {
+			t.Fatalf("rc %d, want exitOK\nstderr %s", code, stderr.String())
+		}
+		got, err := os.ReadFile(path)
+		if err != nil {
+			t.Fatal(err)
+		}
+		if want := requestMessage(t, a); !bytes.Equal(got, want) || len(got) != 90 {
+			t.Fatalf("the file holds %x, want the 90-byte message %x", got, want)
+		}
+		row := a.row(t)
+		mr, err := approval.Decode(row.Request)
+		if err != nil {
+			t.Fatal(err)
+		}
+		want := "wrote request " + a.id + "'s 90-byte approval message to " + path + "\n" +
+			"sign it with an enrolled key (the request expires at 2027-01-15T08:10:00Z):\n" +
+			"  ssh-keygen -Y sign -n world-approval@sunholo -f <your key> < " + path + " > " + path + ".sig\n" +
+			"then approve it:\n" +
+			"  ailang-worldd approve " + a.id + " --signature " + path + ".sig --db " + shellQuote(mr.WorldDB) + " --credentials-db " + shellQuote(mr.CredentialsDB) + "\n"
+		if stdout.String() != want {
+			t.Fatalf("stdout:\n%s\nwant:\n%s", stdout.String(), want)
+		}
+		if opened || calls != 0 || row.Consumed || len(a.sessionRows(t)) != 0 {
+			t.Fatalf("terminal opened %v, daemon calls %d, consumed %v; want none, 0, pending", opened, calls, row.Consumed)
+		}
+		sig := writeSignature(t, sigSigner.Sign(approval.Namespace, got))
+		stdout.Reset()
+		stderr.Reset()
+		if code := runApprove(a.args("--signature", sig), &stdout, &stderr, approveEnv(nil, &opened, atTime(approveT0+2))); code != exitOK {
+			t.Fatalf("approve --signature over the written message: rc %d\nstderr %s", code, stderr.String())
+		}
+		if row := a.row(t); row.ApprovedBy != "ssh:mark-mac" {
+			t.Fatalf("approved_by %q, want ssh:mark-mac", row.ApprovedBy)
+		}
+	})
+	t.Run("existing-file", func(t *testing.T) {
+		a, _ := signedRig(t)
+		path := filepath.Join(t.TempDir(), "request.msg")
+		if err := os.WriteFile(path, []byte("keep\n"), 0o600); err != nil {
+			t.Fatal(err)
+		}
+		var stdout, stderr bytes.Buffer
+		if code := runApprove(a.args("--write-message", path), &stdout, &stderr, approveEnv(nil, nil, atTime(approveT0+1))); code != exitUsage {
+			t.Fatalf("rc %d, want exitUsage\nstderr %s", code, stderr.String())
+		}
+		if want := "ailang-worldd approve: --write-message: open " + path + ": file exists; nothing was written\n"; stderr.String() != want {
+			t.Fatalf("stderr %q, want %q", stderr.String(), want)
+		}
+		if got, _ := os.ReadFile(path); string(got) != "keep\n" {
+			t.Fatalf("the existing file now holds %q", got)
+		}
+	})
+	t.Run("used-request", func(t *testing.T) {
+		a, _ := signedRig(t)
+		var stdout, stderr bytes.Buffer
+		if code := runApprove(a.args(), &stdout, &stderr, approveEnv(&promptTerm{in: strings.NewReader("y\n")}, nil, atTime(approveT0+1))); code != exitOK {
+			t.Fatalf("typed approve: rc %d\nstderr %s", code, stderr.String())
+		}
+		path := filepath.Join(t.TempDir(), "request.msg")
+		stderr.Reset()
+		if code := runApprove(a.args("--write-message", path), &stdout, &stderr, approveEnv(nil, nil, atTime(approveT0+2))); code != exitUsage || !strings.Contains(stderr.String(), "was already used") {
+			t.Fatalf("rc %d, stderr %q; want exitUsage and the admission refusal", code, stderr.String())
+		}
+		if _, err := os.Lstat(path); !errors.Is(err, os.ErrNotExist) {
+			t.Fatalf("a message was written for a used request (err %v)", err)
+		}
+	})
+}
+
+// CR1 (M5c): --signature with --write-message (F1), and --list with either (F2), print the help and
+// exit exitUsage before any store is read.
+func TestApproveFlagRules(t *testing.T) {
+	db := dbPathOf(t)
+	for name, args := range map[string][]string{
+		"signature-and-write-message": {firstRequestID, "--signature", "s", "--write-message", "m", "--db", db},
+		"list-and-signature":          {"--list", "--signature", "s", "--db", db},
+		"list-and-write-message":      {"--list", "--write-message", "m", "--db", db},
+	} {
+		t.Run(name, func(t *testing.T) {
+			var stdout, stderr bytes.Buffer
+			if code := runApprove(args, &stdout, &stderr, approveEnv(nil, nil, atTime(approveT0))); code != exitUsage || stderr.String() != approveHelp {
+				t.Fatalf("rc %d, stderr %q; want exitUsage and exactly the help", code, stderr.String())
+			}
+			if _, err := os.Lstat(filepath.Join(filepath.Dir(db), "credentials.db")); !errors.Is(err, os.ErrNotExist) {
+				t.Fatalf("a credentials store was touched (err %v)", err)
+			}
+		})
+	}
+}
--- a/cmd/ailang-worldd/help_test.go
+++ b/cmd/ailang-worldd/help_test.go
@@ -41,7 +41,7 @@
 		"log tail", "registry get", "commit --file <commit.json> [--session <file|token>]",
 		"tools list", "] call <tool>", "] why ", "provenance", "setup", "doctor",
 		"session new <episode>", "session list", "session revoke", "session mint", "help [<verb>]",
-		"approve <request-id>", "approve <request-id> --signature <file>", "approve --list", "approver enroll --name", "approver list", "approver revoke <name>",
+		"approve <request-id>", "approve <request-id> --signature <file>", "approve <request-id> --write-message <file>", "approve --list", "approver enroll --name", "approver list", "approver revoke <name>",
 	} {
 		if !strings.Contains(out.String(), want) {
 			t.Errorf("top-level usage lacks %q", want)
--- a/cmd/ailang-worldd/main.go
+++ b/cmd/ailang-worldd/main.go
@@ -123,6 +123,7 @@
                     [--sessions-dir <dir>] [--exec-profiles-dir <dir>] [--addr 127.0.0.1:7644 (or service.json's address)] [--request]
   ailang-worldd approve <request-id> [--db <path>] [--credentials-db <path>]
   ailang-worldd approve <request-id> --signature <file> [--db <path>] [--credentials-db <path>]
+  ailang-worldd approve <request-id> --write-message <file> [--db <path>] [--credentials-db <path>]
   ailang-worldd approve --list [--db <path>] [--credentials-db <path>]
   ailang-worldd approver enroll --name <name> --format ssh --pubkey <file> [--db <path>] [--credentials-db <path>]
   ailang-worldd approver list [--db <path>] [--credentials-db <path>]
--- a/docs/QUICKSTART.md
+++ b/docs/QUICKSTART.md
@@ -1784,14 +1784,23 @@
 
 ## 15. Enrol an approver *(attended — pending first verbatim run)*
 
-An approver is a name and one SSH public key, kept in the credentials store. `ailang-worldd approver enroll` reads the `.pub` file once, shows its SHA256 fingerprint, and stores it only after a `y` typed at the terminal within 2 minutes; `approver revoke` marks it revoked the same way, and keeps the row. Approving a mint request with the key's signature (`approve --signature`) arrives in a later sprint; until then an enrolled key approves nothing. This guide uses the §0 scratch store `/tmp/se-world/world.db`, with its daemon serving (`serve` upgrades the credentials store to schema 5), and a scratch key made for it: never a key from `~/.ssh` and never a real store.
+An approver is a name and one SSH public key, kept in the credentials store. `ailang-worldd approver enroll` reads the `.pub` file once, shows its SHA256 fingerprint, and stores it only after a `y` typed at the terminal within 2 minutes; `approver revoke` marks it revoked the same way, and keeps the row. An enrolled key approves a mint request with its signature, with no terminal (below). This guide uses the §0 scratch store `/tmp/se-world/world.db`, with its daemon serving (`serve` upgrades the credentials store to schema 5), and a scratch key made for it: never a key from `~/.ssh` and never a real store.
 
 ```bash
 ssh-keygen -t ed25519 -N '' -f /tmp/se-world/demo-key < /dev/null
 ssh-keygen -l -f /tmp/se-world/demo-key.pub
 ailang-worldd approver enroll --name demo --format ssh --pubkey /tmp/se-world/demo-key.pub --db /tmp/se-world/world.db
 ailang-worldd approver list --db /tmp/se-world/world.db
+```
+
+To approve a §14 request with that key instead of a `y`: write the request's 90-byte message, sign it out of band with `ssh-keygen` in World's namespace, and pass the signature. `approve --signature` opens no terminal; it checks the signature against the enrolled, unrevoked key before any network call, re-checks the approver inside the mint's own transaction, and records `approved_by=ssh:demo`. `--write-message` prints these same two commands with the request's own paths.
+
+```bash
+ailang-worldd episode start sig-demo --repo "$PWD" --db /tmp/se-world/world.db --request
+ailang-worldd approve mr-<the 16 hex it printed> --write-message /tmp/se-world/req.msg --db /tmp/se-world/world.db
+ssh-keygen -Y sign -n world-approval@sunholo -f /tmp/se-world/demo-key < /tmp/se-world/req.msg > /tmp/se-world/req.msg.sig
+ailang-worldd approve mr-<the same id> --signature /tmp/se-world/req.msg.sig --db /tmp/se-world/world.db
 ailang-worldd approver revoke demo --db /tmp/se-world/world.db
 ```
 
-Compare the fingerprint `ssh-keygen -l` prints with the one the y/N shows before you type `y`. `enroll` refuses, before any terminal opens, a private key file (without showing it), a key type other than ssh-ed25519 or ecdsa-sha2-nistp256, a name or a key already enrolled (revoked ones included), and a credentials store below schema 5. A revoked name and key can never be enrolled again.
+Compare the fingerprint `ssh-keygen -l` prints with the one the y/N shows before you type `y`. `enroll` refuses, before any terminal opens, a private key file (without showing it), a key type other than ssh-ed25519 or ecdsa-sha2-nistp256, a name or a key already enrolled (revoked ones included), and a credentials store below schema 5. A revoked name and key can never be enrolled again, and a revoked key approves nothing: `approve --signature` refuses it, and a revoke that lands while an approve is minting refuses that mint. `approve --signature` also refuses a file over 16 KiB or not a regular file, a signature in another namespace or over another request's message, and a key that is not enrolled.
````

