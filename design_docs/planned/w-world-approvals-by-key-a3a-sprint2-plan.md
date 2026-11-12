# Sprint plan — w-world-approvals-by-key-a3a, sprint 2 (iteration 385, row 201)

> **THE SPRINT — A3a sprint 2, "the routes": M4 the read routes, the surface origin set `O`, the guards
> and `serve --surface-origin` (AM1–AM3, RG1, RT1–RT3, RT7, CF1), and M5 the POST (RT4, RT5, RT8–RT12,
> E2E1), as the A3 child design r2 §6 "Sprint 2" and §9 list them, under the design's own pre-named
> fallback cut: the SSHSIG arm of the POST (R11's `sshsig`, RT6, RT13, AW6's route half) moves to sprint 3
> (Appendix C).** Plus iteration 384's carried test obligation: `TestApproveAwaitJudgesTheApproverAtNow`
> (AW-NOW), which kills the minter judging at the request's creation instead of the poll's `now`.
> - **Measured: 1,995 changed lines under `cmd/ host/`** — `git diff --numstat <base> HEAD -- cmd host`,
>   insertions + deletions, after a commit of each of the eight scratch stages (689 production: the pure
>   model 267, the guards and the origin set 91, `daemon.go` 18, the routes 291, `main.go` 22; 1,306 test).
>   Margin to the 2,000 cap: **5**. Without the cut the design's own estimate (≈ 1,900 + RT6/RT13's ≈ 110)
>   and this tree's measured density put the sprint over the cap (§1.5), so **the fallback cut fires**,
>   exactly as design §9 names it; the guards R2–R5, the origin set `O` and `--surface-origin`/CF1 are not
>   cut. A POST carrying `sshsig` is refused 422 `NotApproved` naming `approve <id> --signature <file>`
>   (RT4 arm `R11-sshsig-cut`), so the cut is a refusal, never a silent accept.
> - **Safety order is the design's**: the pure model (M4a–M4b) and the guards with `O` and the flag
>   (M4c) land before any route; the read routes (M4d–M4e) before the route that writes (M5a); the POST
>   judges with the A2 law **before** it records (RT4's `MUT-RECORD-BEFORE-LAW` is killed by the 0-rows
>   check of every law arm), opens a **non-waiting** writer per POST and closes it at once (RT10, RT12:
>   the daemon's production writer sites are exactly `ensureCredentialsStore`'s upgrade and the POST's
>   record).
> - **Every line of code and test is a unified diff in Appendix P**, generated LAST from compiled,
>   vetted, gofmt-clean, green scratch stages; the patches were then extracted **from this file** with §2's
>   `awk`, applied to a fresh `git archive a8c66655` export, each sub-milestone's acceptance run, and the
>   result compared byte-for-byte with the mutant tree (§1.2 C24). The executor applies each patch and runs
>   its acceptance; it writes no code by hand.

- **Design**: [w-world-approvals-by-key-a3.md](w-world-approvals-by-key-a3.md), **r2** (quorum r1+r2 +
  controller carve-out + independent audit VALID 84; **not reopened**). D-WORLD-114 (no Tauri half in A3a;
  the CLI minter re-judges and mints) is OPEN with default A = keep both, applied since iteration 384:
  built as designed. Where this plan and the design differ, this plan was measured against `a8c66655` and
  wins; each difference is a decision (§1.3) or a premise found false (§1.4).
- **Starting ground**: A3a sprint 1 merged (`b44b5f7d`, PR #554; recorded `621ebf9c`): credentials schema 7
  `surface_approvals` (`host/store/surface_approvals.go`), the shared judge (`host/authority/judge.go`) and
  the minter `approve <id> --await` (`cmd/ailang-worldd/approve_await.go`). The sprint-1 plan's §1.4 P1–P6 are
  facts of HEAD this plan builds on (P3: the minter's root is `approve_await.go|pollSurfaceApproval`; P4:
  the judge's steps include `request` and `store`).
- **Planner base**: `a8c666554f284651de33d2f6a2c85a16f5b72c02` (= origin/dev). Worktree
  `/Users/voightkampff/dev/sunholo-data/.wt-world-iter385-a3a-s2`, branch
  `sprint/world-iter385-row201-a3a-s2`. Planner scratch (outside the repo, never under `/tmp`):
  `$HOME/.ailang/state/world-iter385/plan-scratch/` (`base/` the pristine export, `work/` the stage repo
  (tags `M4a` … `M5c`) and the tree the mutants were cut from, `patches/`, `tools/`, `logs/`, `roundtrip/`,
  `mut/` the mutant workers' copies).
- **Executor**: as routed by the controller. **Three runs**: Run 1 (M4a, M4b, M4c), Run 2 (M4d, M4e), Run 3
  (M5a, M5b, M5c). Each sub-milestone is one `patch -p1` (§2) plus its acceptance (§3). **Apply, never
  retype.**
- **Touched** (§5): new `host/surfacepage/{approve_model.go, approve_model_test.go}`; new
  `cmd/ailang-worldd/approve_await_now_test.go` (M4a); `host/daemon/daemon.go` (`Config.SurfaceOrigins`, its
  startup check, the `approveNow` clock, three `HandleFunc`), new `host/daemon/{surface_approve_guard.go,
  surface_approve_guard_test.go}`, `cmd/ailang-worldd/main.go` (`--surface-origin`), new
  `cmd/ailang-worldd/serve_surface_origin_test.go` (M4c); new `host/daemon/{surface_approve.go,
  surface_approve_rig_test.go, surface_approve_test.go}`, `host/daemon/surface_app_test.go` (M4d); new
  `host/daemon/surface_approve_read_test.go` (M4e); new `host/daemon/surface_approve_post_test.go` (M5a); new
  `host/daemon/surface_approve_post_guard_test.go` (M5b); new `cmd/ailang-worldd/approve_surface_e2e_test.go`
  (M5c). **Outside the design's §8 file list**: `approve_await_now_test.go` (the carried obligation),
  `surface_approve_rig_test.go`/`surface_approve_read_test.go`/`surface_approve_post_test.go`/
  `surface_approve_post_guard_test.go`/`serve_surface_origin_test.go`/`approve_surface_e2e_test.go` (the
  design's "+ tests", split by sub-milestone), and `host/daemon/surface_app_test.go` (named by design §6
  RT1/§7 as an edit, though §8 does not list it).
- **Not touched**: every `host/approval` file (the law stays pure and standard-library-only), `host/store`
  and `host/authority` (sprint 1's table, record and judge are used as they are), `surfaceCSP`,
  `/surface/passkey`, `passkey.js`, `app.js`, `isProtected`, every `/v1` route, `world/*.ail`, `packages/**`,
  `scripts/**`, `tools/launchd/*`, `cmd/world-publish/*`, `docs/**`, any skill, `~/.ssh`, `~/.ailang/world`,
  `~/.ailang/state/mission-v1*`.

---

## 0. Binding constraints (each is restated where it governs)

1. **Test keys and temp stores only (D-WORLD-109 = A).** Every store is under `t.TempDir()`; every key is
   `approvaltest.NewPasskey`/`NewSigner`'s fixed seed; every request reaches the daemon through
   `Handler().ServeHTTP` with an `httptest` recorder and a context-carried local address (no listener, no
   socket); the minter's clock and timer are sprint 1's `awaitEnv` fakes. Nothing opens `~/.ailang/world`,
   `~/.ssh`, a real credentials store or `/dev/tty`; nothing mints into a real store, publishes, runs a
   browser, node or `ssh-keygen`.
2. **Pure core, effects at the boundary (coding standards S1–S3).** The JSON shapes, the POST body grammar,
   the admission texts (`host/surfacepage/approve_model.go`) and the guards with the origin set
   (`host/daemon/surface_approve_guard.go`) are pure functions, tested without a daemon or a store (AM1–AM3,
   RG1). Effects (HTTP, the read-only store handles, the one short writer) live in
   `host/daemon/surface_approve.go`. `host/approval` is not touched. No `.ail` (design §3), no new Go module,
   package, production exec site, `/v1` route or context root (the handlers use the request context).
3. **The daemon never keeps the credentials writer.** The POST opens `store.OpenCredentials` (non-waiting)
   after every check has passed, records, and closes it before it answers; a busy writer is `503 Busy` with
   `Retry-After: 1`, never a wait (RT4 `R13-busy`, RT10, RT12). No CORS header anywhere; `isProtected` is
   unchanged; no cookie or session is read (RT11).
4. **One mutant per new production branch, anchored to the diff** (§4), each run against its single named
   test **and** its package suite, named equivalents with their reasons; every run appends a ledger line.
5. **No `t.Parallel()`** in any new test (`TestNoParallelWriteForkPackagesOutsideVerifygate`).
6. **Bisectable and safe at every commit.** Each sub-milestone ends gofmt-clean, `go build ./...` and `go
   vet` green, the `_test.go` compile fence (`go test -run '^$' ./...`) green, its own tests and every earlier
   one green. Guards and the Host/Origin checks land before any route that writes; the POST judges before it
   records.
7. **No git writes by the executor** (§7). The controller commits from snapshots.

---

## 1. Pre-flight

### 1.1 Baseline at `a8c66655` (measured by the planner, 2026-11-12)

Every command ran with `export AILANG_BIN=$HOME/.pinned-ailang/ailang`, `go1.26.6 darwin/arm64`, in a `git
archive a8c66655` export at `plan-scratch/base` (a throwaway git repo; `diff -r cmd host` against the worktree:
identical, `BASE_EQUALS_WORKTREE`; the worktree itself was only read). The pre-commit legs are
`tools/gates.sh <tree> <logs> base` = §6 without the sprint's 20 new tests (they do not exist at the base).
Logs: `plan-scratch/logs/baseline/`.

| Command | rc | Time | Result |
|---|---|---|---|
| `git -C <worktree> status --porcelain \| wc -l` | 0 | — | `0` (clean at `a8c66655`) |
| `./scripts/verify_ail.sh` | **0** | — | `✓ verify gate PASSED: 16 required identities verified, 40 named tests pass, 503 se-tools named tests pass` |
| `gofmt -l cmd host`; `go build ./...`; `go vet ./...`; `go test -run '^$' ./...` | **0** | — | empty / clean / empty / no FAIL line |
| `go test -count=1 -timeout 20m ./host/surfacepage/ ./host/store/ ./host/authority/ ./host/approval/... ./cmd/ailang-worldd/ ./host/daemon/` | **0** | — | all `ok` (`host/surfacepage` 1.6 s, `host/store` 13.5 s, `cmd/ailang-worldd` 127.0 s, `host/daemon` 286.2 s) |
| `-race` legs (§6's three, at the base's tests of the same patterns) | **0** | — | `ok host/surfacepage` 3.0 s, `ok host/daemon` 1.7 s, `ok cmd/ailang-worldd` 29.5 s; no race report |
| §6's 35 named inventory, schema, reflective, help and quickstart tests, `-v` | **0** | — | 35 `--- PASS`, 0 SKIP, 0 FAIL; the AW leg 8 `--- PASS: TestApproveAwait` |
| the whole pre-commit script (`tools/gates.sh base`) | **0** | **457 s** | `PRE_COMMIT_GATES_GREEN 457s` |
| `VERIFY_GO_LEGS=plain ./scripts/verify_go.sh` (the base export, committed) | **1** under load; the two red packages **0** alone | — | under run1's mutant load: `host/evidence` and `host/verifygate` red (load-sensitive); alone: `ok host/evidence 10.888s`, `ok host/verifygate 167.830s` (§1.2 C23) |

**Rule (carried):** a red at the base measures the repo, not the sprint. The controller re-runs a red
untouched package alone once and does not chase it. **Never start a gate under `nohup`.**
`VERIFY_GO_LEGS=plain ./scripts/verify_go.sh` reads committed blobs and runs `go test ./...`; it is the
after-commit leg (§6).

### 1.2 Codebase claims checked (at `a8c66655` unless "work" — the planner's final scratch tree)

| # | Claim | Command → observed output |
|---|---|---|
| C1 | The base | `git rev-parse HEAD origin/dev` (worktree) → `a8c666554f284651de33d2f6a2c85a16f5b72c02` twice; `git status --porcelain \| wc -l` → `0` |
| C2 | **Nothing of sprint 2 has landed** | `grep -rn 'requests.json\|"/surface/approve' --include='*.go' cmd host \| grep -v _test \| wc -l` → `0`; `grep -rn 'SurfaceOrigins\|surface-origin' --include='*.go' cmd host \| grep -v _test \| wc -l` → `0`; control in the same call: `grep -rn '/surface/passkey' --include='*.go' cmd host \| grep -v _test \| wc -l` → `11` |
| C3 | **Design drift since `0140bd7f`** | `shasum -a 256` of design V1's twelve files: `daemon.go 9a9357770a4b0360`, `surface_page.go f769ee06b047f113`, `surfacepage/model.go 4ace112b84f959b4`, `app.js 8ba97c83d7e78a69`, `passkey.js 0064c6c90099a356`, `session_new.go 587f39f99ace013d`, `passkey_cdp.mjs 2a336e5e810f8ed6` **equal the design's**; `approve.go 6a38843161c7d75f`, `approve_passkey.go e17e59d114ef5648`, `approve_sig.go c50db2326f6a52ef`, `mint_requests.go a1501959398c48eb`, `credentials.go 4273fad7537d9905` differ (sprint 1, `git diff --stat 0a25c0d6 a8c66655 -- cmd host` → `26 files changed, 1467 insertions(+), 76 deletions(-)`). So every `host/daemon` line the design cites is current; `main.go`'s are +1 (C9) |
| C4 | The route block and `isProtected` (design §1.6) | `grep -n 'mux.HandleFunc("GET /surface/passkey\|stack.json\|^func (d \*Daemon) Handler\|return r.Method == http.MethodPost' host/daemon/daemon.go` → `1090:func (d *Daemon) Handler()`, `1115`/`1116` the passkey page and script, `1119` `GET /surface/stack.json` (the last surface route), `1140: return r.Method == http.MethodPost && r.URL.Path == "/v1/commit"` |
| C5 | **The daemon's credentials-writer sites** (design §2.2, RT12) | `grep -n "store\.\(OpenCredentials\|UpgradeCredentials\|OpenCredentialsReadOnly\)(" host/daemon/*.go \| grep -v _test` → `credentials.go:66: upgraded, err := store.UpgradeCredentials(path)` and five `OpenCredentialsReadOnly` (`binding.go:84`, `credentials.go:60`, `episode_start.go:17`, `surface_follow_episodes.go:68`, `surface_page.go:199`); **no `store.OpenCredentials(` call** — the premise of §1.4 P1 |
| C6 | `surfaceCSP` and the headers | `host/daemon/surface_page.go:19 const surfaceCSP = "default-src 'none'; … frame-ancestors 'none'"`; unchanged on the final tree (`git diff base M5c -- host/daemon/surface_page.go host/daemon/workbench.go \| wc -l` → `0`, work) |
| C7 | **The connection's local address is the bound port** | `grep -rn LocalAddrContextKey --include='*.go' host cmd \| wc -l` → `0` at the base (no production reader yet); net/http's `Server.Serve` stores `http.LocalAddrContextKey` (the accepted connection's local `net.Addr`, i.e. the listener's address, ephemeral port included) in every request context; an `httptest.NewRequest` carries none (work: RT7's "no local address" arm → `421`) — the premise of S2-D2 |
| C8 | **The `mint_requests` table holds the challenge to 64 characters** | `host/store/credentials_schema.sql:48` and `mint_requests.go:24`: `challenge TEXT NOT NULL CHECK (length(challenge) = 64)`; work: a raw `UPDATE … SET challenge = 'abab…ab'` (62 chars) → `CHECK constraint failed: length(challenge) = 64 (275)` — so R10 is "decodes to 32 bytes", the only distinguishable form (S2-D9) |
| C9 | `serve`'s flags and `Config` (design §1.16) | `cmd/ailang-worldd/main.go:509 func runServe`, `:613 cfg := daemon.Config{DBPath: …` (design: `:508`/`:612`, +1 since sprint 1); `host/approval/origin.go:23 func ParseCeremonyOrigin` |
| C10 | **The CLI's admission texts live in `package main`** | `cmd/ailang-worldd/approve.go:98 func approveAdmission(row store.MintRequestRow, now int64) (approval.MintRequest, string)` — the daemon cannot call it (design §2.4 R8 "the CLI's own refusal text"); §1.4 P3 |
| C11 | **The minter's two judge calls read `now`** (the carried obligation) | `cmd/ailang-worldd/approve_await.go:91 … JudgeSSHSignature(ctx, ro, p.row, p.rec.Artifact, now)`, `:94 … JudgeAssertion(ctx, ro, p.row, p.rec.Artifact, now)` |
| C12 | **AW-NOW kills the 384 mutation** (work) | `sed -i '' 's/Artifact, now)/Artifact, p.row.CreatedAt)/' cmd/ailang-worldd/approve_await.go` → `git diff --stat` `2 insertions(+), 2 deletions(-)`; `go test ./cmd/ailang-worldd/ -run '^TestApproveAwaitJudgesTheApproverAtNow$'` → **`--- FAIL`** both subtests: `refused the surface's approval: … approval: the clock is behind the approver's enrolment time; waiting for another` then `no approval for request … by 2027-01-15T08:00:15Z`; restored → `--- PASS` (0.47 s) (`logs/awnow_proof.txt`); the two single-line halves are §4's `MUT-AW-JUDGE-AT-CREATION`/`-PK` |
| C13 | No test enumerates `Config` fields or `serve` flags beyond the named ones | `grep -rln 'proof-required-identity' --include='*_test.go' cmd host` → `cmd/ailang-worldd/main_test.go` (`TestServeProofFlags`: per-flag regexes, not a full list), `demo_test.go` (a comment); `grep -rln 'VisitAll(' --include='*_test.go' cmd host` → `cmd/world-publish/wiring_test.go` only (another binary) |
| C14 | **Route-enumerating tests** | `TestFrozenV1RouteTableMatchesMux` (`route_table_test.go`: `/v1/` patterns only, 13); `TestSurfaceAppRoutes` (`surface_app_test.go:150-158`: a 404 list under `/surface/`); `TestMCPMountMethodSource` (`mcp_test.go:18`: AST of `Handler`, `d.projection.*` selector args only — a `d.approvalRoute(...)` call argument is skipped); `TestReadCtxCancelledAfterHandler` (`read_deadline_test.go:536`: world-store reads only). No other test lists the mux |
| C15 | `host/daemon`'s import of `host/approval` is new but adds no module | `go list -f '{{join .Imports " "}}' ./host/daemon` (base) has no `host/approval`; `host/authority` (imported) imports it; `TestDaemonDependencyAllowlist` classifies modules, not packages (`daemon_test.go:899`); work: `TestDaemonDependencyAllowlist`, `TestWorldBoundaryDependencyAllowlist` `--- PASS` (§6) |
| C16 | `host/surfacepage` stays transport- and store-free | base imports: `bytes embed fmt …/host/surface html/template math net/url regexp sort strconv strings time unicode unicode/utf8`; work adds `crypto/sha256 encoding/base64 encoding/hex encoding/json errors io …/host/approval` (standard library + the pure law); `go list -deps ./host/workbench/...` holds no `host/surfacepage` (the workbench-closure gate is unaffected) |
| C17 | The store APIs the routes read | `host/store/mint_requests.go`: `:46 MintRequestIDPattern`, `:157 MintRequest`, `:177 PendingMintRequests` (cap 16), `:209 MaxDecidedMintRequests = 64`, `:215 DecidedMintRequests(ctx, approvedBy, limit)` (`""` = no filter, design V30); `passkeys.go:209 ActivePasskeys`; `surface_approvals.go:80 RecordSurfaceApproval`, `:130 SurfaceApproval`; `writer_lock.go:61 IsWriterAlreadyActive` |
| C18 | The judge's steps (sprint 1 P4) | `host/authority/judge.go:20-28` `request, parse, store, unknown-key, stored-key, law`; `JudgeError.Unwrap` → `Err` (`:44`) |
| C19 | **Every repo-wide inventory, on the final tree** | the §6 script (`tools/gates.sh <tree> <logs> final` = §6's pre-commit block) on the round-trip tree (byte-identical to the mutant tree, C24) → `PRE_COMMIT_GATES_GREEN 347s`; pkgs: `host/surfacepage` ok 0.556s, `host/store` ok 12.047s, `host/authority` ok 0.708s, `host/approval` ok 0.822s, `host/approval/approvaltest` ok 1.163s, `cmd/ailang-worldd` ok 101.464s, `host/daemon` ok 270.743s (`logs/final_gates/`): design §7's inventories by name — `TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential` (no exec site added), `TestDaemonDependencyAllowlist`, `TestWorldBoundaryDependencyAllowlist`, `TestEveryCommandDirectoryIsAProtectedGroup`, `TestProductionContextRoots` (no root added), `TestProductionGoSurface`, `TestNoNewDeadlineFreeStoreReads`, `TestCredentialsMethodsRequireADeadline`, `TestFrozenV1RouteTableMatchesMux` (13), `TestSurfaceAppRoutes` (edited), `TestReadCtxCancelledAfterHandler`, `TestApprovalImportsOnlyTheStandardLibrary`, `TestNoParallelWriteForkPackagesOutsideVerifygate`, `TestSubprocessSinkGate`, `TestEveryVerbHasHelp`, `TestServeProofFlags`, the quickstart tests, sprint 1's schema/judge tests and the 20 new tests incl. **`TestDaemonCredentialWriterSites`** — 55 named `--- PASS`, 0 SKIP/FAIL; the AW leg 9 `--- PASS` |
| C20 | **Every mutant, on the final tree** | 115 rows, the last ledger line of each: **114 KILLED (named and suite), 1 named equivalent SURVIVED (named and suite)**, 0 not as expected; 230 ledger lines in all (§4; `logs/mutant_ledger.tsv`). run1 (`logs/mut_run1.out`): 105 as expected, 8 `NOCOMPILE` (re-cut to compiling edits), `MUT-HEX32-LEN` alive at its named test and **`MUT-R11-NOW` alive at both** (no rig passkey was enrolled after a request) — both fixed in the tests (AM2 `challenge-short`; the rig enrols at `apT0-10`), then run2-final: `DONE 115 mutants, 115 MATCH` |
| C21 | `go vet ./...`, gofmt, `-race`, `verify_ail.sh` on the final tree | same §6 run: `go vet ./...` empty; `gofmt -l cmd host` empty; `-race` → `ok host/surfacepage`, `ok host/daemon`, `ok cmd/ailang-worldd`; `./scripts/verify_ail.sh` → `✓ verify gate PASSED: 16 required identities verified, 40 named tests pass, 503 se-tools named tests pass` (no `.ail` change) |
| C22 | **The whole repository on the final tree** (committed) | `VERIFY_GO_LEGS=plain ./scripts/verify_go.sh` on the committed round-trip tree at M5c, nothing else running → **rc 0**: `✓ 0 other binary blobs among 6776 tracked files`; `✓ toolchain floor gate: go1.26.6 >= root module floor go1.26.6`; `✓ all 64 required top-level evidence tests passed exactly once`; `✓ go gate PASSED (plain): build clean; legs run pass with pinned AILANG_BIN (AILANG v0.41.0` (`logs/final_verify_go.log`) |
| C23 | The base's `verify_go.sh` | `VERIFY_GO_LEGS=plain ./scripts/verify_go.sh` on the committed base export, run while run1's 4-worker mutant loop loaded the machine → **rc 1**: `✓ 0 other binary blobs among 6763 tracked files`, `✓ toolchain floor gate`, `✓ all 64 required top-level evidence tests passed exactly once`, and two untouched packages red under that load — `host/evidence` (`--- FAIL: TestProducerCallerCancelIsAnError`) and `host/verifygate` (`TestMissionConfigDefaultFlowIntegration/valid_ledger` and `/invalid_ledger`, then the 600 s package timeout in `TestExecGateArmsAreLoadBearing`) (`logs/base_verify_go.log`). Re-run alone at the base, nothing else running: `go test -count=1 -timeout 20m ./host/evidence/ ./host/verifygate/` → `ok host/evidence 10.888s`, `ok host/verifygate 167.830s` (`logs/base_rerun_alone.log`): load, not the base (§6's re-run rule; `TestMissionConfigDefaultFlowIntegration/invalid_ledger` and the local `host/verifygate` are on the known list). The final tree's run (C22) was made with nothing else running |
| C24 | **The patches reproduce the mutant tree** | a fresh `git archive a8c66655` export + `patch -p1` of Appendix P's eight patches **extracted from this file with §2's `awk`** (each `cmp`-equal to the generated patch; `--dry-run` first; no `.orig`/`.rej`) → `diff -r cmd host` against `plan-scratch/work` (the tree §4's mutants were cut from): **identical** (`ROUNDTRIP_IDENTICAL_TO_MUTANT_TREE`) (`tools/roundtrip.sh`, `logs/roundtrip.log`); nothing uncommitted after M5c |
| C25 | **Every sub-milestone is bisectable** | the same fresh export, after **each** patch: `tools/accept.sh <tree> <M>` = §3 verbatim → `ACCEPT M4a rc=0` 38s, `ACCEPT M4b rc=0` 16s, `ACCEPT M4c rc=0` 18s, `ACCEPT M4d rc=0` 270s, `ACCEPT M4e rc=0` 6s, `ACCEPT M5a rc=0` 257s, `ACCEPT M5b rc=0` 10s, `ACCEPT M5c rc=0` 106s (`logs/rt/`); the cumulative numstat after a commit of each = §7: M4a **350**, M4b **517**, M4c **793**, M4d **1257**, M4e **1396**, M5a **1741**, M5b **1912**, M5c **1995** (`NUMSTAT` lines of `logs/roundtrip.log`) |
| C26 | The gate script | §6's pre-commit block is `tools/gates.sh <tree> <logs> final`; `base` mode (§1.1) omits the 20 new names and expects 8 `TestApproveAwait` passes |

### 1.3 Planner decisions where the design leaves a choice

- **S2-D1 — The fallback cut fires (design §9).** The SSHSIG arm moves to sprint 3: R11's `sshsig` branch,
  RT6, RT13 and AW6's route half (Appendix C). Measured: the sprint as built is 1,995 lines (§1.5); RT6 and
  RT13 alone are ≈ 110 lines at this tree's density (the RT5 twin is 40, the configured-origin arms of RT7/RT8
  are 30), which the design itself adds on top of its ≈ 1,900. The guards R2–R5, `O`, `--surface-origin` and
  CF1 stay; so does `O`'s configured member in the route context (RT7's configured Host, RT8's configured
  Origin through the loopback Host — the proxy case of RT13, with a passkey). The POST refuses `sshsig` with
  `422 NotApproved` naming `approve <id> --signature <file>` (RT4).
- **S2-D2 — The loopback member's port is the request's connection's local port**
  (`http.LocalAddrContextKey`, C7): it is the listener's port whether `--bind` named it or the kernel chose
  it (design V21), it is what `tailscale serve` connects to, and it needs no field on `Daemon`. A request
  with no local address (a recorder, a direct handler call) has **no** loopback member: the check fails
  closed (RT7).
- **S2-D3 — R2 runs before R1.** All three routes are registered through `d.approvalRoute(handler)`, which
  computes `O` once per request (one short read-only open for the enrolled origins) and answers `421
  MisdirectedRequest` before the handler sees the path; a foreign Host learns nothing, not even the id
  grammar. The handler signature `func(w, r, o []string)` makes an unwrapped approval route a compile error,
  and RT1 asserts every `/surface/requests*`/`/surface/approve*` registration is wrapped.
- **S2-D4 — Host matching** is exact `host[:port]` equality with an origin of `O` after lower-casing the
  Host (DNS names are case-insensitive; browsers send lower case); no prefix, suffix or default-port
  folding (RG1: `localhost:76440`, `127.0.0.1:7644.evil`, `p.tail0.ts.net` without its port and
  `s.tail0.ts.net:443` are refused). Origin matching is exact string membership (no case folding: a browser
  serializes origins in lower case).
- **S2-D5 — The 421 and 403 texts name the fix only for a `*.ts.net` name** (`start serve with
  --surface-origin https://<host>`), never for a foreign or loopback name (RG1).
- **S2-D6 — R8 reproduces the CLI's admission texts in the pure model** (`surfacepage.AdmissionRefusal`,
  byte-identical wording to `approveAdmission`, C10, pinned by AM3 `TestApproveModelAdmission`); the CLI is
  not edited.
- **S2-D7 — R11's refusal** is `422 NotApproved` with `the passkey's assertion does not approve request <id>:
  <the judge's error>` for every judge step but `store`, which (a store below the passkeys' schema 6) is
  `503 Restart`, as R13/R14's below-7 is. R12 (`Binds`) is kept as written in the design; after a verified
  judgement of the same row it cannot fail (named equivalent, §4).
- **S2-D8 — R14's race** is any of the five request sentinels `RecordSurfaceApproval` returns (consumed,
  expired, clock behind, digest changed, unknown) → `409 Conflict` (`isRequestRace`, pure, RG1 pins all five);
  the record reads the clock **again** (`received_at` is when it is recorded), so a request that expires
  between the checks and the write is a lost race, which RT4's `R14-race` arm reproduces with a two-step
  clock. Other record errors go through `approvalStoreError` (below 7 → `503 Restart`, busy → `503 Busy`,
  else the sanitized 500).
- **S2-D9 — R10 is `len(hex.DecodeString(challenge)) == 32`**: the table's CHECK holds the column to 64
  characters (C8), so "not hex" and "not 32 bytes" are one condition; RT4 reaches it with a hand-edited
  non-hex challenge.
- **S2-D10 — R6's body bound is checked by `http.MaxBytesReader`**; any read error is `413 TooLarge`.
- **S2-D11 — The routes' clock is a `Daemon` field `approveNow`** (New: the wall clock; tests pin it), like
  `finishNow`. E2E1 uses the wall clock through `daemon.New` (its request is created at `time.Now()`).
- **S2-D12 — `requests.json`'s `recorded`** is the `surface_approvals` row's `approver` column (display
  only, design §2.3) or `null`; below schema 7 (no table) every row is `null`, not `unavailable` (RT2's v6
  arm). A store that cannot be opened or read answers `200` with `"unavailable":"the mint requests could not
  be read"` and `[]` lists, and logs the error for the operator. The card on a store below 6 has
  `"passkeys":[]` (RT3's v5 arm); a store that cannot be opened is the sanitized `500 Internal` (RT3).
- **S2-D13 — R7's store error is the sanitized 500**: `serve` upgrades the store to 7 at start, so a store
  below 4 at request time means the store was replaced under the daemon.
- **S2-D14 — `Config.SurfaceOrigins` is checked twice**: by the flag (`fs.Func`, so `serve` exits
  `exitUsage` naming `-surface-origin` before `runDaemon`, CF1) and by `New` (`StageConfig`, naming
  `--surface-origin`, RG1), so a `Config` built in Go cannot bypass the grammar.
- **S2-D15 — RT1/RT7's "five routes" are this sprint's three** (§1.4 P2): `GET /surface/requests.json`,
  `GET /surface/requests/{id}`, `POST /surface/approve/{id}`. RT1 asserts exactly these three are registered
  behind `approvalRoute`, that `/surface/approve.js` and `/surface/approve` are still `404` (sprint 3 serves
  them), and that no other method is answered on the read routes; `TestSurfaceAppRoutes` gains the two
  sprint-3 names in its 404 list plus `/surface/requests` and `/surface/requests.json/x`, and asserts the
  two read routes answer its non-surface Host with `421`. (A `GET /surface/approve/<id>` is the mux's `405`
  until sprint 3 registers the page.)
- **S2-D16 — The carried obligation lands first (M4a)**: it needs no route; both formats (the two
  single-line halves of iteration 384's two-line `sed` are separate §4 mutants, each killed by its subtest).
- **S2-D17 — Sub-milestones** (§1.5): M4a the read-side model + AW-NOW (350), M4b the POST-side model (167),
  M4c the guards, `O`, the flag (276), M4d the read routes and the Host check + RT1/RT7 (464), M4e RT2/RT3
  (139), M5a the POST + RT4/RT5 (345), M5b RT8–RT12 (171), M5c E2E1 (83). M4d is above ~350: it is one
  production unit (the routes, their wrapper and registration, 175 lines) with the shared test rig (163) and
  the two tests that pin its guard (RT1, RT7); splitting the rig away would land the routes untested.

### 1.4 Premises found false or incomplete at HEAD (each routed around minimally)

- **P1 — Design §2.2/§6 RT12: "exactly two production `store.OpenCredentials` call sites in `host/daemon`:
  `ensureCredentialsStore` and the POST".** False at HEAD: `ensureCredentialsStore` has no
  `store.OpenCredentials` call; its short writable open is `store.UpgradeCredentials(path)`
  (`credentials.go:66`, C5). Routed: `TestDaemonCredentialWriterSites` pins the writer-opening calls
  (`store.OpenCredentials`, `store.UpgradeCredentials`) by enclosing function to exactly
  `{ensureCredentialsStore: UpgradeCredentials, handleSurfaceApprove: OpenCredentials}`, with a labelled
  control that the same walk finds the read-only opens (≥ 5). The design's intent ("the second short writable
  open in the daemon, beside `ensureCredentialsStore`'s") holds.
- **P2 — Design §0/§6 RT1/RT7/T7: "the five new routes"** includes `GET /surface/approve/{id}` and `GET
  /surface/approve.js`, which design §2.7/§9 assign to sprint 3 (M6) and the directive excludes. Routed by
  S2-D15: sprint 2 registers three; RT1/RT7 test those three; the sprint-3 plan extends RT1's list and
  RT7's targets with the page and its script (Appendix C).
- **P3 — Design §2.4 R8: "409 Conflict with the CLI's own refusal text (`approveAdmission`)".** The text
  cannot be shared without editing `package main` (C10). Routed by S2-D6 (the same words, pinned in AM3).
- **P4 — Design §6 RT10/T12: "the page retries 3 times then says so"** is sprint 3's page; this sprint pins
  the daemon half (503 Busy, Retry-After: 1, no wait, no kept writer).
- **P5 — Design §2.4 R13: "below v7 → 503 Restart"** is detected by `RecordSurfaceApproval`'s own S1 check
  (`*OlderCredentialsStoreError`) after the non-waiting open, not by a version read before it; same status
  and class (RT4 `R14-v6`). A store below 6 is refused one step earlier, by the judge's `store` step (RT4
  `R11-store-v5`).
- **P6 — Design §6 RT11 "a valid session Bearer and a cookie with no signature → refused (R6/R11)"**: the
  refusal classes are the handler's own (`422 NotApproved` for a garbage assertion, `400 BadRequest` for no
  body), proving the middleware did not stop or authorize the request; `isProtected` is asserted unchanged
  directly, and a source check pins that `surface_approve.go` reads no `Cookie`/`Authorization` and calls no
  `authority.FromContext`.
- **P7 — Design §6 RT2 "pending/decided caps"**: the pending cap (16) is the store's `PutMintRequest` refusal
  (sprint A0) and cannot be exceeded to observe a route cap; RT2 pins the decided cap (21 consumed → 20
  listed, newest first) and `PendingMintRequests(now, 16)` is the call.

### 1.5 The size

Lines are insertions + deletions under `cmd/ host/`: **`git diff --numstat <base> HEAD -- cmd host`**
(`docs/` and `design_docs/` are outside the filter). **Every row is measured** from the scratch stages, each
committed in a throwaway repo (`plan-scratch/work`, tags `M4a` … `M5c`; `tools/build_patches.py`).

| Sub-milestone | Production | Test | Cumulative vs `a8c66655` (after its commit) | Run |
|---|---|---|---|---|
| M4a — the read-side model (summaries, card); AM1, AM2; **AW-NOW** (carried) | 171 | 179 | **350** | 1 |
| M4b — the POST-side model (body grammar, admission texts, refusal/success bodies); AM3 | 96 | 71 | **517** | 1 |
| M4c — `O`, the guards, `isRequestRace`, `Config.SurfaceOrigins` + its startup check, `serve --surface-origin`; RG1, CF1 | 122 | 154 | **793** | 1 |
| M4d — the read routes behind the Host check (`approvalRoute`, `approveNow`, two `HandleFunc`); the rig; RT1, RT7; `TestSurfaceAppRoutes` edit | 175 | 289 | **1,257** | 2 |
| M4e — RT2, RT3 | 0 | 139 | **1,396** | 2 |
| M5a — `POST /surface/approve/{id}` (R1, R3–R14), its `HandleFunc`; RT4, RT5; RT1's list +1 | 127 | 222 | **1,741** | 3 |
| M5b — RT8, RT9, RT10, RT11, RT12 | 0 | 171 | **1,912** | 3 |
| M5c — E2E1 | 0 | 83 | **1,995** | 3 |
| **Sprint 2** | **689** | **1,306** | **1,995 ≤ 2,000** (margin 5) | |

(Per-row production/test are the step's own numstat; the cumulative column is measured after each commit,
so a line changed in two steps counts once: the step rows sum to 691 + 1,308.) The design estimated sprint 2 at ≈ 1,900 **plus** ≈ 200 for
`--surface-origin` and RT13 (§0 header) — ≈ 2,100 with the SSHSIG arm; with the cut it measured 1,995.

---

## 2. How the executor applies a sub-milestone (the patch protocol)

Each sub-milestone `<M>` is the block fenced by a line `` ````diff <M> `` (four backticks) and a line of
four backticks, in Appendix P. Paths are `a/<repo path>` and `b/<repo path>`; a new file's patch has an
`@@ -0,0 …` hunk. Apply it **only** this way, in one `bash` shell (the controller's shell is zsh):

```bash
bash <<'EOF'
set -euo pipefail
M=M4a   # <- the sub-milestone, exactly as in the Appendix P heading: M4a M4b M4c M4d M4e M5a M5b M5c
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter385-a3a-s2
PLAN="$W/design_docs/planned/w-world-approvals-by-key-a3a-sprint2-plan.md"
T="$HOME/.ailang/state/world-iter385/exec-patches"; mkdir -p "$T"
awk -v m="$M" '$0 == "````diff " m {f=1; next} f && $0 == "````" {exit} f' "$PLAN" > "$T/$M.patch"
test -s "$T/$M.patch"
cd "$W"
patch -p1 --dry-run < "$T/$M.patch"
patch -p1 < "$T/$M.patch"
if find cmd host \( -name '*.orig' -o -name '*.rej' \) | grep .; then echo "STOP: patch left .orig/.rej"; exit 1; fi
git status --short
EOF
```

- If `--dry-run` fails, **stop and report** (no hand edits, no fuzz options). It can only fail if an earlier
  sub-milestone was skipped or changed.
- Then run that sub-milestone's **Acceptance** (§3) in one bash shell with `export
  AILANG_BIN=$HOME/.pinned-ailang/ailang` and **`unset WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER
  WORLD_PASSKEY_RECORD_DIR`**, each command's output `tee`'d to a log under
  `$HOME/.ailang/state/world-iter385/exec-logs/` (never only through `tail`).
- Then snapshot (§7). The controller commits and compares the numstat with §7's expected value.

---

## 3. Milestones

Every sub-milestone's acceptance starts with the same three lines (abbreviated **[base]** below):

```
test -z "$(gofmt -l cmd host)"
go build ./... && go vet ./host/surfacepage/ ./host/daemon/ ./cmd/ailang-worldd/
go test -run '^$' ./...
```

Each `-v` run below must show `--- PASS` for every test it names and **no `--- FAIL` and no `--- SKIP`**;
each package run must end `ok` with no `FAIL` line. The planner's `tools/accept.sh <tree> <M> <logs>` is
exactly these lines (§1.2 C25).

### M4a — the read-side model and the carried obligation · 171 prod + 179 test · Run 1

**Patch** `M4a`:
- new `host/surfacepage/approve_model.go` (read side): `MaxDecidedListed = 20`, `StoredRequest`,
  `PendingRequest`, `DecidedRequest`, `RequestsJSON`, `hex32`, `decodeStored` (bytes hash to the digest,
  decode canonically, the challenge is hex, `approval.MessageHash`), `Pending` (summary or `refused` with no
  summary; `recorded` or `null`; `hash16` = the CLI's message-hash prefix), `Decided` (`credential16`),
  `CardPasskey`, `CardPasskeys`, `CardSSH`, `Card`, `RequestCard` (base64url bytes, the row's digest/nonce/
  times, passkeys grouped by origin, the SSH namespace, **no message-hash field**).
- new `host/surfacepage/approve_model_test.go`: **AM1** `TestApproveModelSummaries`, **AM2**
  `TestApproveModelRefusesWhatDoesNotDecode` (digest, canonical form, non-hex and 31-byte challenge, negative expiry) and
  `TestApproveModelCard`.
- new `cmd/ailang-worldd/approve_await_now_test.go`: **AW-NOW** `TestApproveAwaitJudgesTheApproverAtNow`
  (an approver enrolled at +5, polls from +10, both formats, mints; iteration 384's carried obligation, C12).

**Acceptance**
```
[base]
go test ./host/surfacepage/ -run '^(TestApproveModelSummaries|TestApproveModelRefusesWhatDoesNotDecode|TestApproveModelCard)$' -count=1 -v
go test ./cmd/ailang-worldd/ -run '^(TestApproveAwaitJudgesTheApproverAtNow|TestApproveAwaitMintsARecordedPasskeyApproval|TestApproveAwaitSSHSIG)$' -count=1 -v
go test ./host/surfacepage/ -count=1
```

**Expected cumulative numstat after the controller's commit of M4a: 350**.

### M4b — the POST-side model · 96 prod + 71 test · Run 1

**Patch** `M4b`: `host/surfacepage/approve_model.go` gains `MaxApproveBodyBytes = 20 << 10`, `ApproveBody`,
`ErrApproveBody`, `ParseApproveBody` (strict JSON: no unknown field, exactly one value, `v == 1`, exactly one
of `webauthn`/`sshsig`, the digest 64 lowercase hex), `lowerHex64`, `AdmissionRefusal` (S2-D6), `rfc3339`,
`ApproveRefusal`, `ApproveRecorded`, `Recorded`. `approve_model_test.go`: **AM3** `TestApproveModelBodyGrammar`,
`TestApproveModelAdmission`.

**Acceptance**
```
[base]
go test ./host/surfacepage/ -run '^(TestApproveModelSummaries|TestApproveModelRefusesWhatDoesNotDecode|TestApproveModelCard|TestApproveModelBodyGrammar|TestApproveModelAdmission)$' -count=1 -v
go test ./host/surfacepage/ -count=1
```

**Expected cumulative numstat after the controller's commit of M4b: 517**.

### M4c — the origin set, the guards and `--surface-origin` · 122 prod + 154 test · Run 1

**Patch** `M4c`:
- new `host/daemon/surface_approve_guard.go` (pure): `surfaceOrigins(port, configured, enrolled)` (= `O`),
  `checkSurfaceOrigin` (`approval.ParseCeremonyOrigin`, naming `--surface-origin`), `allowedApprovalHost`
  (R2, S2-D4), `allowedApprovalOrigin` (R3), `sameOriginFetch` (R4), `jsonMediaType` (R5), `isRequestRace`
  (R14, S2-D8), `misdirectedReason`, `crossOriginReason` (S2-D5). Inert: nothing calls them until M4d.
- `host/daemon/daemon.go`: `Config.SurfaceOrigins`; `newDaemon` refuses an invalid one at `StageConfig`
  (S2-D14).
- `cmd/ailang-worldd/main.go`: `--surface-origin` (repeatable, `fs.Func` with `ParseCeremonyOrigin`), wired
  into `Config.SurfaceOrigins`; the usage line and the flag's help paragraph.
- new `host/daemon/surface_approve_guard_test.go`: **RG1** `TestApprovalGuards`. New
  `cmd/ailang-worldd/serve_surface_origin_test.go`: **CF1** `TestServeSurfaceOriginFlag`.

**Acceptance**
```
[base]
go test ./host/daemon/ -run '^(TestApprovalGuards|TestDaemonDependencyAllowlist|TestFrozenV1RouteTableMatchesMux)$' -count=1 -v
go test ./cmd/ailang-worldd/ -run '^(TestServeSurfaceOriginFlag|TestServeProofFlags|TestEveryVerbHasHelp|TestQuickstartServiceSectionMatchesTheCLIHelp)$' -count=1 -v
go test ./cmd/ailang-worldd/ -run '^(TestServe|TestEveryVerbHasHelp|TestQuickstart|TestRun)' -count=1
```

**Expected cumulative numstat after the controller's commit of M4c: 793**.

### M4d — the read routes behind the Host check · 175 prod + 289 test · Run 2

**Patch** `M4d`:
- new `host/daemon/surface_approve.go` (read routes): `storedRequest`, `writeApprovalJSON` (no-store,
  nosniff, no CORS), `writeApprovalRefusal`, `approvalInternal`, `isOlderStore`, `approvalOrigins` (S2-D2:
  the local port, `--surface-origin`, active passkeys' origins on a read-only handle), `approvalRoute` (R2,
  S2-D3), `handleSurfaceRequests` + `readRequests` (`GET /surface/requests.json`, S2-D12),
  `handleSurfaceRequestCard` (`GET /surface/requests/{id}`).
- `host/daemon/daemon.go`: the `approveNow` field and its wall-clock default (S2-D11); `mux.HandleFunc("GET
  /surface/requests.json", d.approvalRoute(…))` and `("GET /surface/requests/{id}", …)` after `stack.json`.
- new `host/daemon/surface_approve_rig_test.go` (the rig: a daemon on temp stores with a pinned clock;
  `putRequest`, `consume`, `enrolPasskey` (at `apT0-10`, after the tests' requests, so a judge at a request's creation refuses), `serve` with a context-carried local address
  `127.0.0.1:7644`, `approvalHeaders`, `requireRefusal`, `setRequestColumn`, `setVersion`).
- new `host/daemon/surface_approve_test.go`: **RT1** `TestApprovalRoutesRegistered`, **RT7**
  `TestApprovalRoutesRefuseForeignHost` (control: `/surface/model.json`, `/surface/passkey`, `/surface/app.js`
  still answer a foreign Host).
- `host/daemon/surface_app_test.go` (`TestSurfaceAppRoutes`): the doc comment, four more 404 names, the two
  read routes `421` for its Host (S2-D15).

**Acceptance** (~6–8 min: the whole `host/daemon` runs here, because `Handler()` changed)
```
[base]
go test ./host/daemon/ -run '^(TestApprovalRoutesRegistered|TestApprovalRoutesRefuseForeignHost|TestApprovalGuards|TestSurfaceAppRoutes|TestFrozenV1RouteTableMatchesMux|TestPasskeyPageRoutes|TestDaemonDependencyAllowlist)$' -count=1 -v
go test -count=1 -timeout 20m ./host/daemon/
```

**Expected cumulative numstat after the controller's commit of M4d: 1257**.

### M4e — what the read routes say · 139 test · Run 2

**Patch** `M4e`: new `host/daemon/surface_approve_read_test.go`: **RT2** `TestRequestsJSON` (pending newest
first with `recorded`/`null`, a non-decoding row `refused`, the expired one absent, 21 decided → 20 newest;
no request bytes, challenge, artifact, token, digest or path; a v6 store lists with `recorded: null`; a
missing store → `unavailable` with `[]` lists), **RT3** `TestRequestCardJSON` (the stored bytes, digest, nonce
and times; passkeys grouped by origin with the revoked one left out; no message hash in hex or base64url; 400
bad id, 404 unknown, consumed shown; a v5 store's card has `passkeys: []`; a missing store → 500 `Internal`).

**Acceptance**
```
[base]
go test ./host/daemon/ -run '^(TestRequestsJSON|TestRequestCardJSON|TestApprovalRoutesRegistered|TestApprovalRoutesRefuseForeignHost)$' -count=1 -v
go test ./host/daemon/ -run '^(TestApprov|TestRequest|TestSurface|TestPasskeyPage)' -count=1
```

**Expected cumulative numstat after the controller's commit of M4e: 1396**.

### M5a — the POST · 127 prod + 222 test · Run 3

**Patch** `M5a`:
- `host/daemon/surface_approve.go`: `approvalStoreError` (below the schema → `503 Restart`; busy → `503 Busy`
  + `Retry-After: 1`; else the sanitized 500) and `handleSurfaceApprove`: R1 id; R3 Origin ∈ `O`; R4
  `Sec-Fetch-Site`; R5 media type; R6 `MaxBytesReader` + `ParseApproveBody`; then on one read-only handle R7
  the row, R8 `AdmissionRefusal`, R9 the body's digest, R10 the challenge, R11 **the sprint-2 cut** (an
  `sshsig` body → 422) and the law `authority.JudgeAssertion` (S2-D7), R12 `Binds`; only then R13 the
  non-waiting writer and R14 `RecordSurfaceApproval` at a fresh clock reading, the writer closed before the
  answer (S2-D8).
- `host/daemon/daemon.go`: `mux.HandleFunc("POST /surface/approve/{id}", d.approvalRoute(d.handleSurfaceApprove))`.
- `host/daemon/surface_approve_test.go`: RT1's route list gains the POST.
- new `host/daemon/surface_approve_post_test.go`: **RT4** `TestApprovePOSTRefusalArms` (21 arms: R1, R3, R4,
  R5, R6 size and grammar, R7, R8 expired and bytes, R9, R10, R11 another request / unenrolled / counter /
  garbage / `sshsig` cut / v5 store, R13 busy, R14 v6 store and race — each its status and class, **0
  `surface_approvals` rows and 0 sessions**), **RT5** `TestApprovePOSTRecordsAVerifiedPasskey`.

**Acceptance** (~6–8 min: the whole `host/daemon`)
```
[base]
go test ./host/daemon/ -run '^(TestApprovePOSTRefusalArms|TestApprovePOSTRecordsAVerifiedPasskey|TestApprovalRoutesRegistered|TestApprovalRoutesRefuseForeignHost|TestRequestsJSON|TestRequestCardJSON|TestApprovalGuards)$' -count=1 -v
go test -count=1 -timeout 20m ./host/daemon/
```

**Expected cumulative numstat after the controller's commit of M5a: 1741**.

### M5b — the POST's guards, replay, writer and inventory · 171 test · Run 3

**Patch** `M5b`: new `host/daemon/surface_approve_post_guard_test.go`: **RT8** `TestApprovePOSTCrossOrigin`
(no Origin, `null`, another port, an unconfigured `*.ts.net` naming the flag, `same-site`, `cross-site`,
`text/plain`, form-encoded → 403/415, 0 rows; `OPTIONS` → 405, no `Access-Control-*`; a configured origin
through the loopback Host → 200), **RT9** `TestApprovePOSTReplay`, **RT10** `TestApprovePOSTWriterBusy`,
**RT11** `TestApprovePOSTNoSessionNoCookie`, **RT12** `TestDaemonCredentialWriterSites` (P1).

**Acceptance**
```
[base]
go test ./host/daemon/ -run '^(TestApprovePOSTCrossOrigin|TestApprovePOSTReplay|TestApprovePOSTWriterBusy|TestApprovePOSTNoSessionNoCookie|TestDaemonCredentialWriterSites)$' -count=1 -v
go test ./host/daemon/ -run '^(TestApprov|TestRequest|TestDaemon|TestSurface|TestPasskeyPage|TestFrozen|TestSessionMiddleware)' -count=1
```

**Expected cumulative numstat after the controller's commit of M5b: 1912**.

### M5c — the end-to-end seam · 83 test · Run 3

**Patch** `M5c`: new `cmd/ailang-worldd/approve_surface_e2e_test.go`: **E2E1**
`TestSurfaceApproveThenAwaitMints` (`episode start --request` on temp stores; a real `daemon.New` +
`Handler()` on the request's credentials store; the card GET; the message hash computed from the card's
digest, nonce and expiry as the page will; a test passkey's assertion POSTed → 200; `approve --await` mints
`webauthn:mark-mac` with no terminal; `requests.json` lists it under `decided`).

**Acceptance** (~4–5 min: the whole `cmd/ailang-worldd`)
```
[base]
go test ./cmd/ailang-worldd/ -run '^(TestSurfaceApproveThenAwaitMints|TestApproveAwaitJudgesTheApproverAtNow|TestServeSurfaceOriginFlag)$' -count=1 -v
go test -count=1 -timeout 20m ./cmd/ailang-worldd/
```

**Expected cumulative numstat after the controller's commit of M5c: 1995**.

---

## 4. Mutation table

**115 rows run on the planner's final tree: 114 KILLED (named and suite), 1 named equivalent SURVIVED as predicted, 0 not as expected.** 230 ledger lines in all (runs: run1 115, run2-final 115).

**Derivation.** Anchored to Appendix P's production hunks, one row per new branch: `approve_model.go` (every `decodeStored` refusal, the recorded/refused/summary arms of `Pending`, `hash16`, the decided truncation and time, the card's encoding, grouping, rpID, empty list and namespace; the body grammar's unknown-field, trailing-value, `v`, both/neither, digest and lower-case arms and the body bound; every admission verdict and its clock), `surface_approve_guard.go` (each member of `O`, the `--surface-origin` grammar, the Host's case folding and exact match against both prefix directions, the Origin's exactness and case, `Sec-Fetch-Site`'s absent/`same-origin`/`same-site` arms, the media type's exactness and parse error, each of the five race sentinels and the race predicate, the two `*.ts.net` fix texts), `daemon.go` (New's origin check, each of the three registrations, the wall clock), `main.go` (the flag's grammar, its append and its wiring), `surface_approve.go` (the Host check, the local port, the enrolled and configured members, no-store and nosniff, every read-route branch: the id, unknown, unavailable, the empty lists, the older-store arms, recorded, the decided cap and filter, the row mapping, the internal status; every POST check R1, R3–R14 in order, R8's and R11's clock, R11's store step and the sprint-2 `sshsig` cut, **`MUT-RECORD-BEFORE-LAW`** (a record inserted before the judge), the received-at clock, the race arm, each `approvalStoreError` arm and `Retry-After`, the recorded artifact, approver and format, the success body, **the kept writer** and **a third writer site**), and **iteration 384's carried mutation** (`approve_await.go`'s two judge calls at `p.row.CreatedAt`, one row per line; the two-line `sed` together: §1.2 C12). **Procedure**: `plan-scratch/tools/mutate.py <run> --workers 4` (outside the repo; the table is `tools/mutants.py`, and every anchor is asserted to occur exactly once in the final tree): for each row, in its own worker copy of the tree, apply **one** edit, build the test binary of the named test's package with `-vet=off` (a mutant that does not compile is recorded `NOCOMPILE`, never counted; run1's eight were re-cut to compiling edits for run2), run **the single named test** and then **its package suite** — `host/surfacepage` whole; `host/daemon` `-test.run '^(TestApprov|TestRequest|TestDaemonCredentialWriterSites|TestSurfaceAppRoutes|TestFrozenV1RouteTableMatchesMux|TestPasskeyPageRoutes)'` and `cmd/ailang-worldd` `-test.run '^(TestApprove|TestServe|TestSurfaceApprove|TestEveryVerbHasHelp|TestQuickstartServiceSectionMatchesTheCLIHelp)'` (the whole packages are 340 s and 230 s per mutant and hold §6's load-sensitive tests) — restore, assert the restore byte-identical, and **append one line to the ledger** `plan-scratch/logs/mutant_ledger.tsv` (UTC time, run, id, stage, named, suite, expected, MATCH/MISMATCH). A row is KILLED only when both fail. **Never commit or snapshot a mutant.** In the table tabs are dropped, a newline prints as ⏎, a backquote as ' and a pipe as \|.

**run1 found two gaps, fixed in the tests (§1.2 C20)**: `MUT-HEX32-LEN` survived its named test (no AM2 row had a 31-byte challenge; AM2 gained `challenge-short`) and `MUT-R11-NOW` (the POST judging at the request's creation) survived everywhere — every rig passkey was enrolled before the requests; the rig now enrols at `apT0-10`, after the tests' requests (`apT0-30`), the POST twin of PK-NOW and AW-NOW.

**Named equivalent** (survives by construction): `MUT-R12` — R11's judge verified the assertion against this very row's digest, challenge and expiry (`ApproveByAssertion` checks the signed message hash, W6), and `Binds` re-checks those same three values (B2–B4), so it cannot fail after a verified judgement; it stays as design §2.4's defense in depth (a skipped judge is caught by it: `MUT-R11-SKIP-LAW` is killed through `Binds`' `ErrNotVerified`).

| Id | Stage | Mutant (single edit) | Named test | Named | Suite |
|---|---|---|---|---|---|
| MUT-AW-JUDGE-AT-CREATION | M4a | approve_await.go: p.rec.Artifact, now); jerr != nil {⏎p.refusal = sshRefusal → p.rec.Artifact, p.row.CreatedAt); jerr != nil {⏎p.refusal = sshRefusal | `TestApproveAwaitJudgesTheApproverAtNow` | KILLED | KILLED |
| MUT-AW-JUDGE-AT-CREATION-PK | M4a | approve_await.go: p.rec.Artifact, now); jerr != nil {⏎p.refusal = passkeyRefusal → p.rec.Artifact, p.row.CreatedAt); jerr != nil {⏎p.refusal = passkeyRefusal | `TestApproveAwaitJudgesTheApproverAtNow` | KILLED | KILLED |
| MUT-HEX32-LEN | M4a | approve_model.go: if err != nil \|\| len(b) != 32 { → if err != nil { | `TestApproveModelRefusesWhatDoesNotDecode` | KILLED | KILLED |
| MUT-DEC-DIGEST | M4a | approve_model.go: if !ok \|\| sha256.Sum256(r.Request) != digest { → if !ok { | `TestApproveModelRefusesWhatDoesNotDecode` | KILLED | KILLED |
| MUT-DEC-CANON | M4a | approve_model.go: mr, err := approval.Decode(r.Request)⏎if err != nil { → mr, err := approval.Decode(r.Request)⏎if false { | `TestApproveModelRefusesWhatDoesNotDecode` | KILLED | KILLED |
| MUT-DEC-CHALLENGE | M4a | approve_model.go: challenge, ok := hex32(r.Challenge)⏎if !ok { → challenge, ok := hex32(r.Challenge)⏎if false && !ok { | `TestApproveModelRefusesWhatDoesNotDecode` | KILLED | KILLED |
| MUT-DEC-EXPIRY | M4a | approve_model.go: h, err := approval.MessageHash(digest, challenge, r.Expires)⏎if err != nil { → h, err := approval.MessageHash(digest, challenge, r.Expires)⏎if false { | `TestApproveModelRefusesWhatDoesNotDecode` | KILLED | KILLED |
| MUT-PEND-RECORDED | M4a | approve_model.go: if recorded != "" { → if true { | `TestApproveModelSummaries` | KILLED | KILLED |
| MUT-PEND-REFUSED | M4a | approve_model.go: mr, h, err := decodeStored(r)⏎if err != nil { → mr, h, err := decodeStored(r)⏎if false { | `TestApproveModelRefusesWhatDoesNotDecode` | KILLED | KILLED |
| MUT-PEND-GRANTS | M4a | approve_model.go: len(mr.Grants), &ttl → len(mr.Grants)+1, &ttl | `TestApproveModelSummaries` | KILLED | KILLED |
| MUT-PEND-HASH16 | M4a | approve_model.go: hex.EncodeToString(h[:8]) → hex.EncodeToString(h[:7]) | `TestApproveModelSummaries` | KILLED | KILLED |
| MUT-DECIDED-GUARD | M4a | approve_model.go: if len(c) > 16 { → if true { | `TestApproveModelSummaries` | KILLED | KILLED |
| MUT-DECIDED-AT | M4a | approve_model.go: At: r.ConsumedAt, → At: r.Created, | `TestApproveModelSummaries` | KILLED | KILLED |
| MUT-CARD-B64 | M4a | approve_model.go: Request: base64.RawURLEncoding.EncodeToString(r.Request) → Request: base64.StdEncoding.EncodeToString(r.Request) | `TestApproveModelCard` | KILLED | KILLED |
| MUT-CARD-EMPTY | M4a | approve_model.go: Passkeys: []CardPasskeys{}, SSH → Passkeys: nil, SSH | `TestApproveModelCard` | KILLED | KILLED |
| MUT-CARD-NS | M4a | approve_model.go: SSH: CardSSH{Namespace: approval.Namespace} → SSH: CardSSH{} | `TestApproveModelCard` | KILLED | KILLED |
| MUT-CARD-GROUP | M4a | approve_model.go: i, ok := at[p.Origin]⏎if !ok { → i, ok := at[p.Origin]⏎if true \|\| !ok { | `TestApproveModelCard` | KILLED | KILLED |
| MUT-CARD-RP | M4a | approve_model.go: RP: p.RPID, Allow → RP: p.Origin, Allow | `TestApproveModelCard` | KILLED | KILLED |
| MUT-BODY-UNKNOWN | M4b | approve_model.go: dec.DisallowUnknownFields()⏎ → ⏎ | `TestApproveModelBodyGrammar` | KILLED | KILLED |
| MUT-BODY-TRAILING | M4b | approve_model.go: if _, err := dec.Token(); err != io.EOF { → if _, err := dec.Token(); false && err != io.EOF { | `TestApproveModelBodyGrammar` | KILLED | KILLED |
| MUT-BODY-V | M4b | approve_model.go: case body.V != 1: → case body.V > 1: | `TestApproveModelBodyGrammar` | KILLED | KILLED |
| MUT-BODY-BOTH | M4b | approve_model.go: case (body.WebAuthn != nil) == (body.SSHSIG != nil): → case body.WebAuthn == nil && body.SSHSIG == nil: | `TestApproveModelBodyGrammar` | KILLED | KILLED |
| MUT-BODY-NEITHER | M4b | approve_model.go: case (body.WebAuthn != nil) == (body.SSHSIG != nil): → case body.WebAuthn != nil && body.SSHSIG != nil: | `TestApproveModelBodyGrammar` | KILLED | KILLED |
| MUT-BODY-DIGEST | M4b | approve_model.go: case !lowerHex64(body.Digest): → case false: | `TestApproveModelBodyGrammar` | KILLED | KILLED |
| MUT-BODY-LOWER | M4b | approve_model.go: return ok && hex.EncodeToString(b[:]) == s → return ok && len(b) == 32 | `TestApproveModelBodyGrammar` | KILLED | KILLED |
| MUT-BODY-LIMIT | M4b | approve_model.go: const MaxApproveBodyBytes = 20 << 10 → const MaxApproveBodyBytes = 21 << 10 | `TestApprovePOSTRefusalArms` (host/daemon) | KILLED | KILLED |
| MUT-ADM-DECODE | M4b | approve_model.go: if _, err := approval.Decode(r.Request); err != nil { → if _, err := approval.Decode(r.Request); false && err != nil { | `TestApproveModelAdmission` | KILLED | KILLED |
| MUT-ADM-ADMITTED | M4b | approve_model.go: case approval.Admitted:⏎ → case approval.Verdict(98):⏎ | `TestApproveModelAdmission` | KILLED | KILLED |
| MUT-ADM-CONSUMED | M4b | approve_model.go: case approval.RefuseConsumed:⏎ → case approval.Verdict(99):⏎ | `TestApproveModelAdmission` | KILLED | KILLED |
| MUT-ADM-CLOCK | M4b | approve_model.go: case approval.RefuseClockBehind:⏎ → case approval.Verdict(97):⏎ | `TestApproveModelAdmission` | KILLED | KILLED |
| MUT-ADM-EXPIRED | M4b | approve_model.go: case approval.RefuseExpired:⏎ → case approval.Verdict(96):⏎ | `TestApproveModelAdmission` | KILLED | KILLED |
| MUT-ADM-NOW | M4b | approve_model.go: approval.AdmitInput{Now: now, → approval.AdmitInput{Now: now - 1, | `TestApproveModelAdmission` | KILLED | KILLED |
| MUT-RECORDED-NEXT | M4b | approve_model.go: Recorded: true, Request: id, → Recorded: true, Request: "", | `TestApproveModelAdmission` | KILLED | KILLED |
| MUT-O-PORT | M4c | surface_approve_guard.go: if port != "" { → if true { | `TestApprovalGuards` | KILLED | KILLED |
| MUT-O-127 | M4c | surface_approve_guard.go: "http://localhost:"+port, "http://127.0.0.1:"+port) → "http://localhost:"+port) | `TestApprovalGuards` | KILLED | KILLED |
| MUT-O-CONFIGURED | M4c | surface_approve_guard.go: o = append(o, configured...)⏎ → ⏎ | `TestApprovalGuards` | KILLED | KILLED |
| MUT-O-ENROLLED | M4c | surface_approve_guard.go: return append(o, enrolled...) → return o | `TestApprovalGuards` | KILLED | KILLED |
| MUT-ORIGIN-GRAMMAR | M4c | surface_approve_guard.go: if _, err := approval.ParseCeremonyOrigin(origin); err != nil {⏎return fmt.Errorf → if _, err := approval.ParseCeremonyOrigin(origin); false && err != nil {⏎return fmt.Errorf | `TestApprovalGuards` | KILLED | KILLED |
| MUT-HOST-CASE | M4c | surface_approve_guard.go: host = strings.ToLower(host)⏎ → ⏎ | `TestApprovalGuards` | KILLED | KILLED |
| MUT-HOST-PREFIX-OF | M4c | surface_approve_guard.go: host == hostPort { → strings.HasPrefix(hostPort, host) { | `TestApprovalGuards` | KILLED | KILLED |
| MUT-HOST-HAS-PREFIX | M4c | surface_approve_guard.go: host == hostPort { → strings.HasPrefix(host, hostPort) { | `TestApprovalGuards` | KILLED | KILLED |
| MUT-ORIGIN-NONEMPTY | M4c | surface_approve_guard.go: return slices.Contains(o, origin) } → return slices.Contains(o, origin) \|\| origin != "null" } | `TestApprovalGuards` | KILLED | KILLED |
| MUT-ORIGIN-FOLD | M4c | surface_approve_guard.go: return slices.Contains(o, origin) } → return slices.ContainsFunc(o, func(x string) bool { return strings.EqualFold(x, origin) }) } | `TestApprovalGuards` | KILLED | KILLED |
| MUT-SFS-SAMESITE | M4c | surface_approve_guard.go: return site == "" \|\| site == "same-origin" } → return site == "" \|\| site == "same-origin" \|\| site == "same-site" } | `TestApprovalGuards` | KILLED | KILLED |
| MUT-SFS-ABSENT | M4c | surface_approve_guard.go: return site == "" \|\| site == "same-origin" } → return site == "same-origin" } | `TestApprovalGuards` | KILLED | KILLED |
| MUT-CT-PREFIX | M4c | surface_approve_guard.go: return err == nil && mt == "application/json" → return err == nil && strings.HasPrefix(mt, "application/json") | `TestApprovalGuards` | KILLED | KILLED |
| MUT-CT-PARSEERR | M4c | surface_approve_guard.go: return err == nil && mt == "application/json" → return err == err && mt == "application/json" | `TestApprovalGuards` | KILLED | KILLED |
| MUT-RACE-CONSUMED | M4c | surface_approve_guard.go: []error{store.ErrRequestConsumed, store.ErrRequestExpired → []error{store.ErrRequestExpired | `TestApprovalGuards` | KILLED | KILLED |
| MUT-RACE-EXPIRED | M4c | surface_approve_guard.go: store.ErrRequestConsumed, store.ErrRequestExpired,  → store.ErrRequestConsumed,  | `TestApprovalGuards` | KILLED | KILLED |
| MUT-RACE-CLOCK | M4c | surface_approve_guard.go: store.ErrRequestClockBehind, store.ErrRequestDigestChanged → store.ErrRequestDigestChanged | `TestApprovalGuards` | KILLED | KILLED |
| MUT-RACE-DIGEST | M4c | surface_approve_guard.go: store.ErrRequestDigestChanged, store.ErrRequestUnknown} → store.ErrRequestUnknown} | `TestApprovalGuards` | KILLED | KILLED |
| MUT-RACE-UNKNOWN | M4c | surface_approve_guard.go: store.ErrRequestDigestChanged, store.ErrRequestUnknown} → store.ErrRequestDigestChanged} | `TestApprovalGuards` | KILLED | KILLED |
| MUT-RACE-ANY | M4c | surface_approve_guard.go: if errors.Is(err, race) {⏎ → if err != nil \|\| errors.Is(err, race) {⏎ | `TestApprovalGuards` | KILLED | KILLED |
| MUT-MISDIRECTED-TSNET | M4c | surface_approve_guard.go: strings.HasSuffix(name, ".ts.net") { → strings.HasSuffix(name, ".ts.net.never") { | `TestApprovalGuards` | KILLED | KILLED |
| MUT-CROSSORIGIN-TSNET | M4c | surface_approve_guard.go: err == nil && strings.HasSuffix(rp, ".ts.net") { → err == nil && strings.HasSuffix(rp, "") { | `TestApprovalGuards` | KILLED | KILLED |
| MUT-NEW-ORIGINS | M4c | daemon.go: if err := checkSurfaceOrigin(origin); err != nil { → if err := checkSurfaceOrigin(origin); false && err != nil { | `TestApprovalGuards` | KILLED | KILLED |
| MUT-CLI-GRAMMAR | M4c | main.go: if _, err := approval.ParseCeremonyOrigin(v); err != nil { → if _, err := approval.ParseCeremonyOrigin(v); false && err != nil { | `TestServeSurfaceOriginFlag` | KILLED | KILLED |
| MUT-CLI-APPEND | M4c | main.go: surfaceOrigins = append(surfaceOrigins, v) → _ = v | `TestServeSurfaceOriginFlag` | KILLED | KILLED |
| MUT-CLI-WIRE | M4c | main.go: SurfaceOrigins: surfaceOrigins,⏎ → ⏎ | `TestServeSurfaceOriginFlag` | KILLED | KILLED |
| MUT-ROUTE-REQUESTS | M4d | daemon.go: mux.HandleFunc("GET /surface/requests.json", d.approvalRoute(d.handleSurfaceRequests))⏎ → ⏎ | `TestApprovalRoutesRegistered` | KILLED | KILLED |
| MUT-ROUTE-CARD | M4d | daemon.go: mux.HandleFunc("GET /surface/requests/{id}", d.approvalRoute(d.handleSurfaceRequestCard))⏎ → ⏎ | `TestApprovalRoutesRegistered` | KILLED | KILLED |
| MUT-HOST-CHECK | M4d | surface_approve.go: if !allowedApprovalHost(r.Host, o) { → if false { | `TestApprovalRoutesRefuseForeignHost` | KILLED | KILLED |
| MUT-LOCALADDR | M4d | surface_approve.go: if addr, ok := r.Context().Value(http.LocalAddrContextKey).(net.Addr); ok { → if addr, ok := r.Context().Value(http.LocalAddrContextKey).(net.Addr); false && ok { | `TestApprovalRoutesRefuseForeignHost` | KILLED | KILLED |
| MUT-ENROLLED-ORIGINS | M4d | surface_approve.go: enrolled = append(enrolled, e.Passkey.Origin)⏎ → _ = e⏎ | `TestApprovalRoutesRefuseForeignHost` | KILLED | KILLED |
| MUT-CONFIGURED-ORIGINS | M4d | surface_approve.go: surfaceOrigins(port, d.cfg.SurfaceOrigins, enrolled) → surfaceOrigins(port, nil, enrolled) | `TestApprovalRoutesRefuseForeignHost` | KILLED | KILLED |
| MUT-NOSTORE | M4d | surface_approve.go: w.Header().Set("Cache-Control", "no-store")⏎ → ⏎ | `TestApprovalRoutesRegistered` | KILLED | KILLED |
| MUT-NOSNIFF | M4d | surface_approve.go: w.Header().Set("X-Content-Type-Options", "nosniff")⏎ → ⏎ | `TestApprovalRoutesRegistered` | KILLED | KILLED |
| MUT-APPROVENOW-WALL | M4d | daemon.go: approveNow:   func() int64 { return time.Now().Unix() }, → approveNow:   func() int64 { return 0 }, | `TestSurfaceApproveThenAwaitMints` (cmd/ailang-worldd) | KILLED | KILLED |
| MUT-CARD-ID | M4d | surface_approve.go: id := r.PathValue("id")⏎if !store.MintRequestIDPattern.MatchString(id) {⏎writeApprovalRefusal(w, http.StatusBadRequest, "BadRequest", "a request id is mr- and 16 lowercase hex")⏎return⏎}⏎ctx, cancel := d.readCtx(r) → id := r.PathValue("id")⏎if false {⏎writeApprovalRefusal(w, http.StatusBadRequest, "BadRequest", "a request id is mr- and 16 lowercase hex")⏎return⏎}⏎ctx, cancel := d.readCtx(r) | `TestRequestCardJSON` | KILLED | KILLED |
| MUT-REQ-UNAVAILABLE | M4e | surface_approve.go: if err := d.readRequests(ctx, now, &out); err != nil { → if err := d.readRequests(ctx, now, &out); false && err != nil { | `TestRequestsJSON` | KILLED | KILLED |
| MUT-REQ-PENDING-EMPTY | M4e | surface_approve.go: if out.Pending == nil { → if false { | `TestRequestsJSON` | KILLED | KILLED |
| MUT-REQ-DECIDED-EMPTY | M4e | surface_approve.go: if out.Decided == nil { → if false { | `TestRequestsJSON` | KILLED | KILLED |
| MUT-REQ-OLDER-RECORD | M4e | surface_approve.go: if err != nil && !isOlderStore(err) {⏎return err → if err != nil {⏎return err | `TestRequestsJSON` | KILLED | KILLED |
| MUT-REQ-RECORDED | M4e | surface_approve.go: surfacepage.Pending(storedRequest(row), rec.Approver) → surfacepage.Pending(storedRequest(row), rec.Approver[:0]) | `TestRequestsJSON` | KILLED | KILLED |
| MUT-REQ-DECIDED-CAP | M4e | surface_approve.go: DecidedMintRequests(ctx, "", surfacepage.MaxDecidedListed) → DecidedMintRequests(ctx, "", store.MaxDecidedMintRequests) | `TestRequestsJSON` | KILLED | KILLED |
| MUT-REQ-DECIDED-ALL | M4e | surface_approve.go: DecidedMintRequests(ctx, "", surfacepage.MaxDecidedListed) → DecidedMintRequests(ctx, "tty", surfacepage.MaxDecidedListed) | `TestSurfaceApproveThenAwaitMints` (cmd/ailang-worldd) | KILLED | KILLED |
| MUT-ROW-CONSUMEDAT | M4e | surface_approve.go: ConsumedAt: row.ConsumedAt, → ConsumedAt: row.CreatedAt, | `TestRequestsJSON` | KILLED | KILLED |
| MUT-ROW-CREDENTIAL | M4e | surface_approve.go: CredentialID: row.CredentialID} → CredentialID: ""} | `TestRequestsJSON` | KILLED | KILLED |
| MUT-ROW-CONSUMED | M4e | surface_approve.go: Consumed: row.Consumed, → Consumed: false, | `TestRequestCardJSON` | KILLED | KILLED |
| MUT-CARD-UNKNOWN | M4e | surface_approve.go: if !ok {⏎writeApprovalRefusal(w, http.StatusNotFound, "NotFound", "no mint request "+id)⏎return⏎}⏎active → if !ok && false {⏎writeApprovalRefusal(w, http.StatusNotFound, "NotFound", "no mint request "+id)⏎return⏎}⏎active | `TestRequestCardJSON` | KILLED | KILLED |
| MUT-CARD-OLDER | M4e | surface_approve.go: active, err := ro.ActivePasskeys(ctx)⏎if err != nil && !isOlderStore(err) { → active, err := ro.ActivePasskeys(ctx)⏎if err != nil { | `TestRequestCardJSON` | KILLED | KILLED |
| MUT-CARD-PASSKEY-RP | M4e | surface_approve.go: RPID: e.Passkey.RPID, → RPID: e.Passkey.Origin, | `TestRequestCardJSON` | KILLED | KILLED |
| MUT-INTERNAL-STATUS | M4e | surface_approve.go: writeApprovalRefusal(w, http.StatusInternalServerError, "Internal", internalErrorMessage) → writeApprovalRefusal(w, http.StatusOK, "Internal", internalErrorMessage) | `TestRequestCardJSON` | KILLED | KILLED |
| MUT-ROUTE-POST | M5a | daemon.go: mux.HandleFunc("POST /surface/approve/{id}", d.approvalRoute(d.handleSurfaceApprove))⏎ → ⏎ | `TestApprovalRoutesRegistered` | KILLED | KILLED |
| MUT-R1 | M5a | surface_approve.go: if !store.MintRequestIDPattern.MatchString(id) { // R1 → if false { // R1 | `TestApprovePOSTRefusalArms` | KILLED | KILLED |
| MUT-R3 | M5a | surface_approve.go: !allowedApprovalOrigin(origin, o) { // R3 → false { // R3 | `TestApprovePOSTRefusalArms` | KILLED | KILLED |
| MUT-R4 | M5a | surface_approve.go: if !sameOriginFetch(r.Header.Get("Sec-Fetch-Site")) { // R4 → if false { // R4 | `TestApprovePOSTRefusalArms` | KILLED | KILLED |
| MUT-R5 | M5a | surface_approve.go: if !jsonMediaType(r.Header.Get("Content-Type")) { // R5 → if false { // R5 | `TestApprovePOSTRefusalArms` | KILLED | KILLED |
| MUT-R6-SIZE | M5a | surface_approve.go: http.MaxBytesReader(w, r.Body, surfacepage.MaxApproveBodyBytes) → http.MaxBytesReader(w, r.Body, 2*surfacepage.MaxApproveBodyBytes) | `TestApprovePOSTRefusalArms` | KILLED | KILLED |
| MUT-R6-GRAMMAR | M5a | surface_approve.go: body, err := surfacepage.ParseApproveBody(raw)⏎if err != nil { → body, err := surfacepage.ParseApproveBody(raw)⏎if false && err != nil { | `TestApprovePOSTRefusalArms` | KILLED | KILLED |
| MUT-R7 | M5a | surface_approve.go: if !ok {⏎writeApprovalRefusal(w, http.StatusNotFound, "NotFound", "no mint request "+id)⏎return⏎}⏎if reason → if !ok && false {⏎writeApprovalRefusal(w, http.StatusNotFound, "NotFound", "no mint request "+id)⏎return⏎}⏎if reason | `TestApprovePOSTRefusalArms` | KILLED | KILLED |
| MUT-R8 | M5a | surface_approve.go: reason != "" { // R8 → false { // R8 | `TestApprovePOSTRefusalArms` | KILLED | KILLED |
| MUT-R8-NOW | M5a | surface_approve.go: surfacepage.AdmissionRefusal(storedRequest(row), now) → surfacepage.AdmissionRefusal(storedRequest(row), row.CreatedAt) | `TestApprovePOSTRefusalArms` | KILLED | KILLED |
| MUT-R9 | M5a | surface_approve.go: if body.Digest != row.Digest { // R9 → if false { // R9 | `TestApprovePOSTRefusalArms` | KILLED | KILLED |
| MUT-R10 | M5a | surface_approve.go: if len(challenge) != 32 { → if false { | `TestApprovePOSTRefusalArms` | KILLED | KILLED |
| MUT-R11-SSHSIG-CUT | M5a | surface_approve.go: if body.SSHSIG != nil { // R11 → if false { // R11 | `TestApprovePOSTRefusalArms` | KILLED | KILLED |
| MUT-R11-SKIP-LAW | M5a | surface_approve.go: signed, err := authority.JudgeAssertion(ctx, ro, row, body.WebAuthn, now) // R11: the law⏎if err != nil { → signed, err := authority.JudgeAssertion(ctx, ro, row, body.WebAuthn, now) // R11: the law⏎if false && err != nil { | `TestApprovePOSTRefusalArms` | KILLED | KILLED |
| MUT-R11-NOW | M5a | surface_approve.go: authority.JudgeAssertion(ctx, ro, row, body.WebAuthn, now) → authority.JudgeAssertion(ctx, ro, row, body.WebAuthn, row.CreatedAt) | `TestApprovePOSTRecordsAVerifiedPasskey` | KILLED | KILLED |
| MUT-R11-STORE-STEP | M5a | surface_approve.go: errors.As(err, &je) && je.Step == authority.JudgeStepStore { → errors.As(err, &je) && je.Step == "never" { | `TestApprovePOSTRefusalArms` | KILLED | KILLED |
| MUT-RECORD-BEFORE-LAW | M5a | surface_approve.go: signed, err := authority.JudgeAssertion( → if w2, e2 := store.OpenCredentials(d.credentialsPath); e2 == nil {⏎_ = w2.RecordSurfaceApproval(ctx, store.SurfaceApprovalRow{RequestID: id, Format: approval.FormatWebAuthn, Approver: "webauthn:unjudged", Artifact: body.WebAuthn, Digest: row.Digest}, now)⏎_ = w2.Close()⏎}⏎signed, err := authority.JudgeAssertion( | `TestApprovePOSTRefusalArms` | KILLED | KILLED |
| MUT-R12 | M5a | surface_approve.go: err := signed.Binds([32]byte(digest), [32]byte(challenge), row.ExpiresAt); err != nil { // R12 → err := signed.Binds([32]byte(digest), [32]byte(challenge), row.ExpiresAt); false && err != nil { // R12 | `TestApprovePOSTRefusalArms` | SURVIVED | SURVIVED |
| MUT-R13-OPEN | M5a | surface_approve.go: wr, err := store.OpenCredentials(d.credentialsPath) // R13: non-waiting, closed at once⏎if err != nil { → wr, err := store.OpenCredentials(d.credentialsPath) // R13: non-waiting, closed at once⏎if false && err != nil { | `TestApprovePOSTRefusalArms` | KILLED | KILLED |
| MUT-R14-RACE | M5a | surface_approve.go: case isRequestRace(err):⏎ → case false:⏎ | `TestApprovePOSTRefusalArms` | KILLED | KILLED |
| MUT-R14-RECEIVED-NOW | M5a | surface_approve.go: Digest: row.Digest}, d.approveNow()) // R14 → Digest: row.Digest}, now) // R14 | `TestApprovePOSTRefusalArms` | KILLED | KILLED |
| MUT-R14-ERR | M5a | surface_approve.go: case err != nil:⏎d.approvalStoreError(w, r, err)⏎ → case false:⏎d.approvalStoreError(w, r, err)⏎ | `TestApprovePOSTRefusalArms` | KILLED | KILLED |
| MUT-STOREERR-OLDER | M5a | surface_approve.go: case isOlderStore(err):⏎ → case false:⏎ | `TestApprovePOSTRefusalArms` | KILLED | KILLED |
| MUT-STOREERR-BUSY | M5a | surface_approve.go: case store.IsWriterAlreadyActive(err):⏎ → case false:⏎ | `TestApprovePOSTRefusalArms` | KILLED | KILLED |
| MUT-RETRY-AFTER | M5a | surface_approve.go: w.Header().Set("Retry-After", "1")⏎ → ⏎ | `TestApprovePOSTRefusalArms` | KILLED | KILLED |
| MUT-OLDER-AS | M5a | surface_approve.go: var older *store.OlderCredentialsStoreError⏎return errors.As(err, &older) → var older *store.OlderCredentialsStoreError⏎return errors.As(err, &older) && false | `TestApprovePOSTRefusalArms` | KILLED | KILLED |
| MUT-REC-ARTIFACT | M5a | surface_approve.go: Artifact: body.WebAuthn, Digest: row.Digest}, d.approveNow()) → Artifact: []byte("{}"), Digest: row.Digest}, d.approveNow()) | `TestApprovePOSTRecordsAVerifiedPasskey` | KILLED | KILLED |
| MUT-REC-APPROVER | M5a | surface_approve.go: Approver: signed.ApprovedBy(),⏎ → Approver: "webauthn:someone",⏎ | `TestApprovePOSTRecordsAVerifiedPasskey` | KILLED | KILLED |
| MUT-REC-FORMAT | M5a | surface_approve.go: Format: approval.FormatWebAuthn, Approver: signed.ApprovedBy() → Format: approval.FormatSSH, Approver: signed.ApprovedBy() | `TestApprovePOSTRecordsAVerifiedPasskey` | KILLED | KILLED |
| MUT-SUCCESS-BODY | M5a | surface_approve.go: surfacepage.Recorded(id, signed.ApprovedBy()) → surfacepage.Recorded(id, "") | `TestApprovePOSTRecordsAVerifiedPasskey` | KILLED | KILLED |
| MUT-WRITER-KEPT | M5b | surface_approve.go: _ = wr.Close()⏎ → _ = wr⏎ | `TestApprovePOSTWriterBusy` | KILLED | KILLED |
| MUT-THIRD-WRITER | M5b | surface_approve.go: now := d.approveNow()⏎out := surfacepage.RequestsJSON{V: 1, Now: now}⏎ → now := d.approveNow()⏎if w2, e2 := store.OpenCredentials(d.credentialsPath); e2 == nil {⏎_ = w2.Close()⏎}⏎out := surfacepage.RequestsJSON{V: 1, Now: now}⏎ | `TestDaemonCredentialWriterSites` | KILLED | KILLED |

---

## 5. Files (final `git status --short` after M5c)

| File | Sub-milestone | Change |
|---|---|---|
| `host/surfacepage/approve_model.go`, `approve_model_test.go` | M4a, M4b | new: the read side, AM1, AM2 (M4a); the POST side, AM3 (M4b) |
| `cmd/ailang-worldd/approve_await_now_test.go` | M4a | new: AW-NOW (carried) |
| `host/daemon/surface_approve_guard.go`, `surface_approve_guard_test.go` | M4c | new: `O`, the guards; RG1 |
| `host/daemon/daemon.go` | M4c, M4d, M5a | `Config.SurfaceOrigins` + check (M4c); `approveNow`, two GET routes (M4d); the POST route (M5a) |
| `cmd/ailang-worldd/main.go`, `serve_surface_origin_test.go` | M4c | `--surface-origin`; new CF1 |
| `host/daemon/surface_approve.go` | M4d, M5a | new: the read routes (M4d); the POST (M5a) |
| `host/daemon/surface_approve_rig_test.go`, `surface_approve_test.go`, `surface_app_test.go` | M4d (+M5a) | new rig; RT1, RT7; the 404 list (RT1's list +1 in M5a) |
| `host/daemon/surface_approve_read_test.go` | M4e | new: RT2, RT3 |
| `host/daemon/surface_approve_post_test.go` | M5a | new: RT4, RT5 |
| `host/daemon/surface_approve_post_guard_test.go` | M5b | new: RT8–RT12 |
| `cmd/ailang-worldd/approve_surface_e2e_test.go` | M5c | new: E2E1 |
| this plan and `sprint_w-world-approvals-by-key-a3a-s2.json` | planner | |

---

## 6. Sprint gates (controller, after the last sub-milestone; each measured at the base, §1.1)

In **one** `bash` shell, from the worktree, **not under `nohup`**. Every leg `tee`s its full output to a log
under `$HOME/.ailang/state/world-iter385/gates/`; a verdict is read from its log, never from `tail` alone. No
browser leg: nothing in this sprint serves a page that a browser drives (sprint 3).

```
bash <<'EOF'
set -euo pipefail
export AILANG_BIN=$HOME/.pinned-ailang/ailang
unset WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER WORLD_PASSKEY_RECORD_DIR
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter385-a3a-s2; cd "$W"
L=$HOME/.ailang/state/world-iter385/gates; mkdir -p "$L"
s0=$(date +%s)
./scripts/verify_ail.sh 2>&1 | tee "$L/verify_ail.log"                      # no .ail change: exactly green
out=$(gofmt -l cmd host); echo "$out" | tee "$L/gofmt.log"; test -z "$out"
go build ./... 2>&1 | tee "$L/build.log"
go vet ./... 2>&1 | tee "$L/vet.log"
go test -run '^$' ./... 2>&1 | tee "$L/fence.log"                           # compile fence for every _test.go
go test -count=1 -timeout 20m ./host/surfacepage/ ./host/store/ ./host/authority/ ./host/approval/... ./cmd/ailang-worldd/ ./host/daemon/ 2>&1 | tee "$L/pkgs.log"
if grep -E '^(FAIL|--- FAIL)' "$L/pkgs.log"; then exit 1; fi
go test -race -count=1 ./host/surfacepage/ 2>&1 | tee "$L/race.log"
go test -race -count=1 ./host/daemon/ -run '^(TestApprov|TestRequest|TestDaemonCredentialWriterSites|TestSurfaceAppRoutes)' 2>&1 | tee -a "$L/race.log"
go test -race -count=1 ./cmd/ailang-worldd/ -run '^(TestApprove|TestServeSurfaceOrigin|TestSurfaceApprove|TestEveryVerbHasHelp)' 2>&1 | tee -a "$L/race.log"
if grep -E '^(FAIL|--- FAIL)|WARNING: DATA RACE' "$L/race.log"; then exit 1; fi
NEW_SP="TestApproveModelSummaries TestApproveModelRefusesWhatDoesNotDecode TestApproveModelCard TestApproveModelBodyGrammar TestApproveModelAdmission"
NEW_D="TestApprovalGuards TestApprovalRoutesRegistered TestApprovalRoutesRefuseForeignHost TestRequestsJSON TestRequestCardJSON TestApprovePOSTRefusalArms TestApprovePOSTRecordsAVerifiedPasskey TestApprovePOSTCrossOrigin TestApprovePOSTReplay TestApprovePOSTWriterBusy TestApprovePOSTNoSessionNoCookie TestDaemonCredentialWriterSites"
NEW_C="TestApproveAwaitJudgesTheApproverAtNow TestServeSurfaceOriginFlag TestSurfaceApproveThenAwaitMints"
alt() { local x="$*"; echo "${x// /|}"; }
# the repo-wide inventories (design §7), the schema pins, the reflective, help and quickstart tests and the 20 new tests, by name, -v:
go test -v -count=1 ./host/broker/ -run '^(TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential|TestQuickstartSection10ProfilesLoad|TestQuickstartProtectedContract)$' 2>&1 | tee "$L/inv.log"
go test -v -count=1 ./host/daemon/ -run "^(TestDaemonDependencyAllowlist|TestFrozenV1RouteTableMatchesMux|TestSurfaceAppRoutes|TestDaemonConstructorCallerInventory|TestQuickstartDecisionPacketSection|TestPasskeyPageRoutes|TestReadCtxCancelledAfterHandler|$(alt $NEW_D))\$" 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./host/boundary/ -run '^(TestWorldBoundaryDependencyAllowlist|TestBareNetHTTPExemptionIsPerGroup|TestEveryCommandDirectoryIsAProtectedGroup)$' 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./host/store/ -run '^(TestProductionContextRoots|TestProductionGoSurface|TestNoNewDeadlineFreeStoreReads|TestReadGettersHonorContext|TestCredentialsMethodsRequireADeadline|TestCredentialsSchemaMatchesCanonicalDDL|TestCredentialsSchema7SurfaceApprovals|TestRecordSurfaceApprovalGuards|TestSurfaceApprovalPrunedWithItsRequest)$' 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./host/verifygate/ -run '^(TestNoParallelWriteForkPackagesOutsideVerifygate|TestSubprocessSinkGate)$' 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./host/approval/... -run '^(TestApprovalImportsOnlyTheStandardLibrary|TestApprovalTestIsTestOnly)$' 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./host/authority/ -run '^(TestJudgeAssertionIsTheCLIJudge|TestJudgeSSHSignatureIsTheCLIJudge)$' 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./cmd/ailang-worldd/ -run "^(TestEveryVerbHasHelp|TestServeProofFlags|TestQuickstartAttendedStepsBoxMatchesSection9|TestQuickstartSection9FlagsMatchTheCLI|TestQuickstartSection10FlagsMatchTheCLI|TestQuickstartSection11FlagsMatchTheCLI|TestQuickstartServiceSectionMatchesTheCLIHelp|$(alt $NEW_C))\$" 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./host/surfacepage/ -run "^($(alt $NEW_SP))\$" 2>&1 | tee -a "$L/inv.log"
for t in TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential TestQuickstartSection10ProfilesLoad TestQuickstartProtectedContract \
  TestDaemonDependencyAllowlist TestFrozenV1RouteTableMatchesMux TestSurfaceAppRoutes TestDaemonConstructorCallerInventory TestQuickstartDecisionPacketSection TestPasskeyPageRoutes TestReadCtxCancelledAfterHandler \
  TestWorldBoundaryDependencyAllowlist TestBareNetHTTPExemptionIsPerGroup TestEveryCommandDirectoryIsAProtectedGroup \
  TestProductionContextRoots TestProductionGoSurface TestNoNewDeadlineFreeStoreReads TestReadGettersHonorContext TestCredentialsMethodsRequireADeadline TestCredentialsSchemaMatchesCanonicalDDL \
  TestCredentialsSchema7SurfaceApprovals TestRecordSurfaceApprovalGuards TestSurfaceApprovalPrunedWithItsRequest \
  TestNoParallelWriteForkPackagesOutsideVerifygate TestSubprocessSinkGate TestApprovalImportsOnlyTheStandardLibrary TestApprovalTestIsTestOnly \
  TestJudgeAssertionIsTheCLIJudge TestJudgeSSHSignatureIsTheCLIJudge \
  TestEveryVerbHasHelp TestServeProofFlags TestQuickstartAttendedStepsBoxMatchesSection9 TestQuickstartSection9FlagsMatchTheCLI TestQuickstartSection10FlagsMatchTheCLI \
  TestQuickstartSection11FlagsMatchTheCLI TestQuickstartServiceSectionMatchesTheCLIHelp $NEW_SP $NEW_D $NEW_C; do
  grep -qE -- "^--- PASS: $t " "$L/inv.log" || { echo "$t did not PASS"; exit 1; }; done
if grep -E '^\s*--- (SKIP|FAIL)' "$L/inv.log"; then exit 1; fi
go test -v -count=1 ./cmd/ailang-worldd/ -run '^TestApproveAwait' 2>&1 | tee "$L/aw.log"
test "$(grep -cE '^--- PASS: TestApproveAwait' "$L/aw.log")" = 9                       # sprint 1's 8 + AW-NOW
if grep -E '^\s*--- (SKIP|FAIL)' "$L/aw.log"; then exit 1; fi
echo "PRE_COMMIT_GATES_GREEN $(( $(date +%s) - s0 ))s" | tee "$L/pre_commit.ok"
EOF
# AFTER the controller commits (the tracked-binary gate reads committed blobs; this leg runs go test ./...):
bash -c 'set -o pipefail; export AILANG_BIN=$HOME/.pinned-ailang/ailang; cd /Users/voightkampff/dev/sunholo-data/.wt-world-iter385-a3a-s2 && VERIFY_GO_LEGS=plain ./scripts/verify_go.sh 2>&1 | tee $HOME/.ailang/state/world-iter385/gates/verify_go.log; echo "verify_go rc=$?" | tee -a $HOME/.ailang/state/world-iter385/gates/verify_go.log'
git -C /Users/voightkampff/dev/sunholo-data/.wt-world-iter385-a3a-s2 status --short   # only the two plan files, until the controller commits them
```

The pre-commit block is the planner's `tools/gates.sh <tree> <logs> final` (§1.2 C26); `base` mode omits
the 20 new names and expects 8 `TestApproveAwait` passes.

**At the base**: every pre-commit leg is green at `a8c66655` in 457 s (§1.1; 35 named tests `--- PASS`, the AW leg 8); `verify_go.sh` at the base was red only under the concurrent mutant load and green alone for its two red packages (§1.2 C23). **On the final tree**: `PRE_COMMIT_GATES_GREEN 347s`; pkgs: `host/surfacepage` ok 0.556s, `host/store` ok 12.047s, `host/authority` ok 0.708s, `host/approval` ok 0.822s, `host/approval/approvaltest` ok 1.163s, `cmd/ailang-worldd` ok 101.464s, `host/daemon` ok 270.743s (55 named tests `--- PASS`, the AW leg 9 `--- PASS`, `-race` ok host/surfacepage, ok host/daemon, ok cmd/ailang-worldd; §1.2 C19, C21), and the after-commit `verify_go.sh` **rc 0**: `✓ 0 other binary blobs among 6776 tracked files`; `✓ toolchain floor gate: go1.26.6 >= root module floor go1.26.6`; `✓ all 64 required top-level evidence tests passed exactly once`; `✓ go gate PASSED (plain): build clean; legs run pass with pinned AILANG_BIN (AILANG v0.41.0` (§1.2 C22).

**Known load- or environment-sensitive (re-run once, alone; do not chase)**:
`TestSeedRestartIgnoresASurvivingOldCopy` (host/broker), `TestHostsLaunchersExecHost`, `host/broker`'s
exec-probe timeouts (incl. `TestExecProtectedCheckTimeout/reserve`) and its local `-race`, the local `-race`
of `host/daemon` and `host/verifygate` (this sprint's `-race` leg is the approval subset only),
`TestDemoSeedPreflightRefusals`, `TestWorktreeSnapshotRefusesAConfiguredFilter`,
`TestMissionConfigDefaultFlowIntegration/invalid_ledger`, `scripts/marketingcheck`'s supervisor tests and
`TestMarketingCaptureFakeChromeSixShots/inventory`, `host/daemon`'s `TestAGUISlowReaderCutIsResumable`,
`TestSessionSelfGit`, `TestServeStartsWhenTheProbeSpendsTheStartupClock`, `TestRegistryPhaseHasItsOwnClock`,
`cmd/ailang-worldd`'s `TestWorldContextAndStatusline` and `TestCLIRealSubprocessEpisode`, and anything run
under `nohup`. Any other red, any race report, any red in `host/surfacepage`, `host/store`, `host/authority`,
`host/approval`, in `host/daemon`'s `TestApprov*`/`TestRequest*`/`TestSurfaceApp*`/`TestDaemonCredential*`
or in `cmd/ailang-worldd`'s `TestApprove*`/`TestServe*`/`TestSurfaceApprove*`, is a regression. CI's
`go-race-verify` is the race verdict.

**Must stay green, with only the named edits**: every a0a1, A1b, A2 and A3a-sprint-1 test; the only edit to
an existing test is `TestSurfaceAppRoutes` (S2-D15). In particular `TestFrozenV1RouteTableMatchesMux` (13,
untouched), `TestPasskeyPageRoutes` (the A2 page unchanged), `TestServeProofFlags`, `TestEveryVerbHasHelp`
(unedited: `serve --help` still prints its usage) and every §7 inventory.

---

## 7. Executor directive notes

- **No git writes.** Never `git add`, `commit`, `stash`, `checkout`, `reset`, `restore`, `rebase`, `push`,
  `worktree` or `apply`. `patch -p1` (§2) is the only writer; read-only `git status`/`git diff` are fine.
- **No browser, no node, no `ssh-keygen`, no real store.** Nothing in this sprint needs Chrome or node; `go
  test` with `WORLD_CHROME` unset (§2) leaves the A2 browser tests as their counted skips, which only the
  whole-package runs include (not a §3 `-v` run).
- **Runs.** Run 1 = M4a, M4b, M4c; Run 2 = M4d, M4e; Run 3 = M5a, M5b, M5c. Each run reads §0, §2 and its own §3
  sections, then applies its patches in order.
- **Expected numstat after each sub-milestone, measured AFTER the controller's commit of it** (an uncommitted
  tree hides the new untracked files from numstat):
  ```
  bash -c 'cd /Users/voightkampff/dev/sunholo-data/.wt-world-iter385-a3a-s2 && git diff --numstat a8c666554f284651de33d2f6a2c85a16f5b72c02 HEAD -- cmd host | awk "{s+=\$1+\$2} END {print s}"'
  ```
  After the controller's commit of M4a **350**, M4b **517**, M4c **793**, M4d **1257**, M4e **1396**, M5a
  **1741**, M5b **1912**, M5c **1995**. **Any other value means a patch was not applied as written: stop and
  report.**
- **Snapshot after every sub-milestone, cumulatively**:
  ```
  bash -c 'set -euo pipefail; M=M4a; W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter385-a3a-s2
    D=$HOME/.ailang/state/world-iter385/snap/$M; rm -rf "$D"; mkdir -p "$D"; cd "$W"
    git status --porcelain --untracked-files=all | while read -r st f; do
      case "$f" in design_docs/planned/w-world-approvals-by-key-a3a-sprint2-plan.md|design_docs/planned/sprint_w-world-approvals-by-key-a3a-s2.json) continue;; esac
      mkdir -p "$D/$(dirname "$f")"; cp -p "$f" "$D/$f"; done
    (cd "$D" && find . -type f ! -name MANIFEST.sha256 | sort | xargs shasum -a 256) > "$D/MANIFEST.sha256"
    cat "$D/MANIFEST.sha256"'
  ```
  (`M` = `M4a` … `M5c`.) No file is deleted, so the porcelain status has only `M` and `??`; stop on any other
  status. The controller rebuilds one commit per sub-milestone (`$HOME/.ailang/state/world-iter385/commit_ms.sh`
  with `WT` and the snapshot root pointed here; commit the plan BEFORE running it — it deletes untracked files).
- **Stop conditions.** A `patch --dry-run` failure; a numstat other than the expected value; a red
  earlier-milestone test; a `--- FAIL`, or a `--- SKIP` in a §3 `-v` run; any need to touch a file outside §5;
  running Chrome, node, `ssh-keygen` or any browser; any reference to `~/.ssh`, `~/.ailang/world` or a real
  `/dev/tty`. Stop and report; do not work around.
- **Timing.** M4d's and M5a's acceptance ~6–8 min each (the whole `host/daemon`, 340 s at the planner's
  load); M5c's whole `cmd/ailang-worldd` ~4–5 min; M4c's CLI subset ~1 min; every other acceptance under 1
  min plus the compile fence (~25–60 s).
- **Shell.** Wrap loops, globs and `$VAR` expansions in `bash -c` or a `bash <<'EOF'` heredoc (zsh does not
  word-split). Export `AILANG_BIN` in the same shell as the Go tests. `tee` every gate's output to a file;
  never pipe a gate's only copy through `tail`. **Never start a gate under `nohup`.** Poll any background
  command's artifact in a bounded `date +%s` loop in the same turn; never end a turn announcing a wait.

---

## 8. Risks

- **R1 — Margin 5 lines.** The sprint is at the cap after the fallback cut; any executor-side change is a
  stop (§7), so the margin only matters if the controller asks for a fix: a fix that adds lines must cut
  elsewhere or move a test to sprint 3.
- **R2 — The minter still waits for a page that does not exist yet.** After this sprint a POST can record an
  approval and `approve --await` mints it (E2E1), but the card page that makes the POST in a browser is
  sprint 3; until then only a test (or a hand-made POST) reaches the route.
- **R3 — The read routes are as open as the existing surface on the loopback** (design T6/R1): any loopback
  process may read the pending requests' summaries and cards; no token, secret, artifact or message hash is
  served (RT2, RT3). Row 199 closes it.
- **R4 — An SSH-only operator cannot approve from the surface until sprint 3** (the cut); `approve <id>
  --signature` at a terminal is unchanged, and the POST's refusal names it.
- **R5 — Behind `tailscale serve`, which `Host` arrives is UNVERIFIED** (design §2.3); both the `*.ts.net`
  Host (when its origin is in `O`) and `127.0.0.1:<p>` (the local port) pass; RT8 pins the proxy case
  (configured Origin through the loopback Host). The attended run records it (design §10.2).
- **R6 — Load-sensitive whole-package runs** (§6's list): M4d and M5a run the whole `host/daemon`; re-run a
  red untouched test alone once.

---

## 9. Acceptance criteria (the judge scores against these)

1. **The read model is pure and exact** (AM1, AM2): a pending line is decoded from the stored bytes (agent,
   model, grant count, TTL, expiry, `hash16` = the message hash's first 8 bytes) with `recorded` the approver
   or `null`; a row whose bytes do not hash to the digest, do not decode canonically, whose challenge is not
   hex or whose expiry is negative is listed with `refused` and no summary; the card is the stored bytes in
   base64url, the row's digest, nonce and times, passkeys grouped by origin, the SSH namespace, and **no
   message-hash field**; the decided line truncates the credential to 16.
2. **The POST body grammar** (AM3): strict JSON, no unknown field, exactly one value, `v == 1`, exactly one of
   `webauthn`/`sshsig`, the digest 64 lowercase hex; R8's admission texts equal the CLI's.
3. **The origin set and the guards** (RG1): `O` = loopback at the arriving port + `--surface-origin` +
   enrolled passkeys' origins; Host exact (case-insensitive), Origin exact, `Sec-Fetch-Site` absent or
   `same-origin`, media type `application/json` (parameters allowed); an invalid `--surface-origin` is
   refused by `New` and by `serve` (exitUsage naming the flag, before the daemon starts; CF1); the 421/403
   texts name `--surface-origin` only for a `*.ts.net` name; the five request races are 409s.
4. **The read routes** (RT1–RT3, RT7): exactly `GET /surface/requests.json`, `GET /surface/requests/{id}` and
   `POST /surface/approve/{id}` are registered, each behind the Host check; `/surface/approve.js` stays 404;
   a foreign Host is 421 on each, loopback/enrolled/configured hosts pass, and an existing route
   (`/surface/model.json`) still answers a foreign Host; the JSON carries no request bytes (list), challenge
   (list), artifact, token or message hash; caps and `unavailable` as §1.3 S2-D12; `/v1` stays 13.
5. **The POST judges before it records** (RT4, RT5): R1, R3–R11, R13, R14 each refuse with their status and
   class and leave **0 `surface_approvals` rows and 0 sessions**; a verified passkey assertion records one
   row (`webauthn`, `webauthn:<name>`, the posted envelope's bytes, the request's digest, `received_at` =
   the daemon's clock at the write) and mints nothing.
6. **Cross-origin, replay, writer, session** (RT8–RT11): every cross-origin, cross-site or non-JSON POST is
   refused, a preflight is 405 with no CORS header; a replay while pending replaces one row and after the
   consume is 409 "already used"; a held writer is 503 `Busy` with `Retry-After: 1` and no wait, and after
   any response the daemon holds no writer; a session Bearer or cookie approves nothing and `isProtected` is
   unchanged.
7. **The writer inventory** (RT12): the daemon's production writer opens are exactly
   `ensureCredentialsStore`'s `UpgradeCredentials` and the POST's `OpenCredentials`.
8. **End to end** (E2E1): the surface's POST records, `approve --await` mints it with no terminal, and
   `requests.json` lists it as decided by `webauthn:<name>`.
9. **The carried obligation** (AW-NOW): the minter judges at the poll's `now`; the 384 mutation (both halves)
   is killed.
10. **The fallback cut is explicit**: an `sshsig` POST is a 422 naming the terminal path; RT6, RT13 and AW6's
    route half are Appendix C's sprint-3 scope.
11. Each §4 mutant, applied alone, makes its named test **and** its package suite fail; the named equivalents
    survive for the stated reasons; the ledger has a line per run.
12. Every §6 gate passes, including every inventory by name, the `-race` legs and `VERIFY_GO_LEGS=plain
    ./scripts/verify_go.sh` after the commit (modulo §6's re-run rule); `git status --short` lists only §5.
13. The numstat after each sub-milestone's commit equals §7's value; the sprint total is ≤ 2,000.
14. `host/approval`, `host/store`, `host/authority`, `surfaceCSP`, `isProtected`, `/surface/passkey`,
    `passkey.js` and `app.js` unchanged; no new production exec site, module, package, context root or `/v1`
    route; no CORS header; no `t.Parallel`.
15. No real key, store, `/dev/tty`, `~/.ssh` access, browser, node, mint into a real store, publish or
    `ssh-keygen` anywhere in the sprint (D-WORLD-109).

---

## Appendix C — after this sprint (sprint-3 scope moved by the fallback cut; no re-design)

- **Moved here by design §9's fallback cut, unchanged from design §2.4/§6**: R11's `sshsig` arm
  (`authority.JudgeSSHSignature` on `*body.SSHSIG`, recorded as `ssh`/`ssh:<name>`; it replaces this sprint's
  `if body.SSHSIG != nil` refusal and its RT4 arm `R11-sshsig-cut`), **RT6**
  `TestApprovePOSTRecordsAVerifiedSSHSIG`, **RT13** `TestApprovePOSTSSHSIGViaTailnetOriginWithNoPasskey` (one
  SSH approver, no passkey, `SurfaceOrigins: ["https://studio.tail0.ts.net"]`, both Hosts, the card GET, the
  unconfigured 403/421, `other.tail0.ts.net` 403) and **AW6's route half** (an SSHSIG POSTed, then `approve
  --await` mints `ssh:<name>`, as E2E1 does for the passkey).
- **Sprint 3's own scope (design §9)**: the card page `GET /surface/approve/{id}` and `GET
  /surface/approve.js` — register both through `d.approvalRoute` (the Host check), add them to RT1's
  `approvalRoutes` and RT7's targets, and drop `/surface/approve.js` and `/surface/approve` from
  `TestSurfaceAppRoutes`' 404 list (S2-D15); `app.js`'s Needs you, Decided and the notification; XA1–XA5,
  PG2, PG3, PG3b; QUICKSTART §17.
- **Carry**: the page's POST body is surfacepage's `ApproveBody`; a `503 Busy` carries `Retry-After: 1`
  (the page's ≤ 3 retries, design §2.7 [r2]); refusals are `{v, error, reason}` with the classes of §2.4's
  table plus `TooLarge` (413) and `Internal` (500).
- **D-WORLD-114** stays open (default A applied); if Mark vetoes (b), the minter's mint location returns to
  design.

---

## Appendix P — the patches (apply with §2; generated from the compiled scratch stages; do not edit)

### Patch M4a — the read-side model and the carried obligation

````diff M4a
--- a/cmd/ailang-worldd/approve_await_now_test.go
+++ b/cmd/ailang-worldd/approve_await_now_test.go
@@ -0,0 +1,41 @@
+package main
+
+// Row 201 A3 sprint 2 M4a (iteration 384's carried obligation): the minter judges a recorded
+// artifact at the poll's own clock reading, never at the request's creation. A judge at
+// p.row.CreatedAt left every sprint-1 await test green (controller-reproduced, iteration 384).
+
+import (
+	"bytes"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+)
+
+// AW-NOW: an approver enrolled after the request was made (+5) and before the poll (+10) approves
+// it through the surface, in both formats: the await twin of PK-NOW and A1b's SIG-NOW.
+func TestApproveAwaitJudgesTheApproverAtNow(t *testing.T) {
+	for _, format := range []string{"webauthn", "ssh"} {
+		t.Run(format, func(t *testing.T) {
+			a := newApproveRig(t, "ep1", approveT0)
+			var artifact []byte
+			if format == "ssh" {
+				seedApprover(t, a.credPath, "mark-mac", sigSigner.PublicKey(), approveT0+5)
+				artifact = sigSigner.Sign(approval.Namespace, requestMessage(t, a))
+			} else {
+				seedPasskey(t, a.credPath, "mark-mac", passkeySigner, approveT0+5)
+				artifact = passkeySigner.Assert(pkOrigin, pkRPID, requestHashOf(t, a), 2)
+			}
+			var opened bool
+			env, _ := awaitEnv(approveT0+10, &opened, func(n, now int64) {
+				if n == 1 {
+					recordSurface(t, a, format, format+":mark-mac", artifact, now)
+				}
+			})
+			var stdout, stderr bytes.Buffer
+			if code := runApprove(a.args("--await", "--await-for", "5s"), &stdout, &stderr, env); code != exitOK || opened {
+				t.Fatalf("rc %d, terminal opened %v; want exitOK with no terminal for an approver enrolled at +5 and polls from +10\nstderr %s", code, opened, stderr.String())
+			}
+			requireMintedBy(t, a, format+":mark-mac", stdout.String())
+		})
+	}
+}
--- a/host/surfacepage/approve_model.go
+++ b/host/surfacepage/approve_model.go
@@ -0,0 +1,171 @@
+package surfacepage
+
+// Row 201 A3 sprint 2 (design_docs/planned/w-world-approvals-by-key-a3.md §2.3, §2.4): the JSON
+// the approval routes speak, built here without a daemon or a store. GET /surface/requests.json
+// lists the pending requests (Needs you) and the decided ones; GET /surface/requests/{id} is the
+// card's data (the stored bytes and the raw nonce, never a message hash); POST
+// /surface/approve/{id} takes ApproveBody and answers ApproveRecorded or ApproveRefusal.
+
+import (
+	"crypto/sha256"
+	"encoding/base64"
+	"encoding/hex"
+	"errors"
+	"fmt"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+)
+
+// MaxDecidedListed is the decision list's length (design §2.3, Q3).
+const MaxDecidedListed = 20
+
+// StoredRequest is one mint_requests row as the routes read it (store.MintRequestRow's fields).
+type StoredRequest struct {
+	ID, Episode, Digest, Challenge string
+	Request                        []byte
+	Created, Expires, ConsumedAt   int64
+	Consumed                       bool
+	ApprovedBy, CredentialID       string
+}
+
+// PendingRequest is one Needs-you line: a summary decoded from the stored bytes, or, when they do
+// not decode or do not hash to the digest, Refused and no summary (never hidden).
+type PendingRequest struct {
+	ID       string  `json:"id"`
+	Episode  string  `json:"episode"`
+	Agent    string  `json:"agent,omitempty"`
+	Model    string  `json:"model,omitempty"`
+	Grants   int     `json:"grants,omitempty"`
+	TTL      *int64  `json:"ttl_s,omitempty"`
+	Expires  int64   `json:"expires"`
+	Hash16   string  `json:"hash16,omitempty"`
+	Recorded *string `json:"recorded"` // the verified approver of a recorded approval, or null
+	Refused  string  `json:"refused,omitempty"`
+}
+
+// DecidedRequest is one line of the decision list.
+type DecidedRequest struct {
+	ID           string `json:"id"`
+	Episode      string `json:"episode"`
+	ApprovedBy   string `json:"approvedBy"`
+	At           int64  `json:"at"`
+	Credential16 string `json:"credential16"`
+}
+
+// RequestsJSON is GET /surface/requests.json. It carries no request bytes, challenge, artifact,
+// token or path.
+type RequestsJSON struct {
+	V           int              `json:"v"`
+	Now         int64            `json:"now"`
+	Pending     []PendingRequest `json:"pending"`
+	Decided     []DecidedRequest `json:"decided"`
+	Unavailable string           `json:"unavailable"`
+}
+
+// hex32 decodes 64 hex characters.
+func hex32(s string) ([32]byte, bool) {
+	b, err := hex.DecodeString(s)
+	if err != nil || len(b) != 32 {
+		return [32]byte{}, false
+	}
+	return [32]byte(b), true
+}
+
+// decodeStored is a stored request's bytes as a canonical request whose SHA-256 is its digest,
+// and its message hash (approval.MessageHash, what a signer signs).
+func decodeStored(r StoredRequest) (approval.MintRequest, [32]byte, error) {
+	digest, ok := hex32(r.Digest)
+	if !ok || sha256.Sum256(r.Request) != digest {
+		return approval.MintRequest{}, [32]byte{}, errors.New("the stored bytes do not match the request's digest")
+	}
+	mr, err := approval.Decode(r.Request)
+	if err != nil {
+		return approval.MintRequest{}, [32]byte{}, fmt.Errorf("the stored bytes are not a canonical mint request (%v)", err)
+	}
+	challenge, ok := hex32(r.Challenge)
+	if !ok {
+		return approval.MintRequest{}, [32]byte{}, errors.New("the stored challenge is not 64 hex")
+	}
+	h, err := approval.MessageHash(digest, challenge, r.Expires)
+	if err != nil {
+		return approval.MintRequest{}, [32]byte{}, err
+	}
+	return mr, h, nil
+}
+
+// Pending is r's Needs-you line; recorded is the verified approver of its recorded approval, or "".
+func Pending(r StoredRequest, recorded string) PendingRequest {
+	p := PendingRequest{ID: r.ID, Episode: r.Episode, Expires: r.Expires}
+	if recorded != "" {
+		p.Recorded = &recorded
+	}
+	mr, h, err := decodeStored(r)
+	if err != nil {
+		p.Refused = err.Error()
+		return p
+	}
+	ttl := mr.TTLSeconds
+	p.Agent, p.Model, p.Grants, p.TTL, p.Hash16 = mr.Agent, mr.Model, len(mr.Grants), &ttl, hex.EncodeToString(h[:8])
+	return p
+}
+
+// Decided is r's line in the decision list.
+func Decided(r StoredRequest) DecidedRequest {
+	c := r.CredentialID
+	if len(c) > 16 {
+		c = c[:16]
+	}
+	return DecidedRequest{ID: r.ID, Episode: r.Episode, ApprovedBy: r.ApprovedBy, At: r.ConsumedAt, Credential16: c}
+}
+
+// CardPasskey is one active passkey as the card offers it.
+type CardPasskey struct {
+	Origin, RPID string
+	CredentialID []byte
+}
+
+// CardPasskeys is the active passkeys of one origin: its rpID and their credential IDs (base64url).
+type CardPasskeys struct {
+	Origin string   `json:"origin"`
+	RP     string   `json:"rp"`
+	Allow  []string `json:"allow"`
+}
+
+// CardSSH is what an SSH signer needs besides the message: the namespace.
+type CardSSH struct {
+	Namespace string `json:"namespace"`
+}
+
+// Card is GET /surface/requests/{id}: the stored bytes and the raw nonce, from which the page
+// computes the digest and the message hash itself. There is no message-hash field.
+type Card struct {
+	V          int            `json:"v"`
+	ID         string         `json:"id"`
+	Request    string         `json:"request"` // base64url of the stored canonical bytes
+	Digest     string         `json:"digest"`
+	Challenge  string         `json:"challenge"`
+	Created    int64          `json:"created"`
+	Expires    int64          `json:"expires"`
+	Consumed   bool           `json:"consumed"`
+	ApprovedBy string         `json:"approvedBy"`
+	Passkeys   []CardPasskeys `json:"passkeys"`
+	SSH        CardSSH        `json:"ssh"`
+}
+
+// RequestCard is r's card, with the active passkeys grouped by origin in the order given.
+func RequestCard(r StoredRequest, passkeys []CardPasskey) Card {
+	c := Card{V: 1, ID: r.ID, Request: base64.RawURLEncoding.EncodeToString(r.Request), Digest: r.Digest, Challenge: r.Challenge,
+		Created: r.Created, Expires: r.Expires, Consumed: r.Consumed, ApprovedBy: r.ApprovedBy,
+		Passkeys: []CardPasskeys{}, SSH: CardSSH{Namespace: approval.Namespace}}
+	at := map[string]int{}
+	for _, p := range passkeys {
+		i, ok := at[p.Origin]
+		if !ok {
+			i = len(c.Passkeys)
+			at[p.Origin] = i
+			c.Passkeys = append(c.Passkeys, CardPasskeys{Origin: p.Origin, RP: p.RPID, Allow: []string{}})
+		}
+		c.Passkeys[i].Allow = append(c.Passkeys[i].Allow, base64.RawURLEncoding.EncodeToString(p.CredentialID))
+	}
+	return c
+}
--- a/host/surfacepage/approve_model_test.go
+++ b/host/surfacepage/approve_model_test.go
@@ -0,0 +1,138 @@
+package surfacepage
+
+// Row 201 A3 sprint 2 M4a (design §2.3, §2.4, §6 AM1-AM3): the approval routes' JSON, pure.
+
+import (
+	"bytes"
+	"crypto/sha256"
+	"encoding/base64"
+	"encoding/hex"
+	"encoding/json"
+	"reflect"
+	"sort"
+	"strings"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+)
+
+const amCreated, amExpires = int64(1800000000), int64(1800000600)
+
+var amChallenge = strings.Repeat("5a", 32)
+
+// amRequest is a stored request whose bytes are canonical and hash to its digest.
+func amRequest(t *testing.T) StoredRequest {
+	t.Helper()
+	b, err := approval.Canonical(approval.MintRequest{Schema: approval.RequestSchema, Episode: "w201x",
+		WorldDB: "/ws/world.db", WorkspaceRoot: "/ws/root", Worktree: "/ws/root/w201x", Git: "/usr/bin/git", Branch: "world/w201x",
+		Grants:     []approval.Grant{{Effect: "Workspace.Read", Scope: "w201x"}, {Effect: "Workspace.Write", Scope: "w201x", Budget: 5}, {Effect: "Exec", Scope: "w201x"}},
+		TTLSeconds: 43200, ProfilesDir: "/ws/profiles", CredentialsDB: "/ws/credentials.db", SessionFile: "/ws/s/session.token",
+		MCPConfig: "/ws/s/mcp.json", Agent: "claude", Model: "claude-opus<5>&x", Task: "/ws/task.md", Daemon: approval.DaemonRef{Addr: "127.0.0.1:7644"}})
+	if err != nil {
+		t.Fatal(err)
+	}
+	d := sha256.Sum256(b)
+	return StoredRequest{ID: "mr-0123456789abcdef", Episode: "w201x", Request: b, Digest: hex.EncodeToString(d[:]), Challenge: amChallenge,
+		Created: amCreated, Expires: amExpires}
+}
+
+// jsonKeys is the sorted top-level keys of v's JSON.
+func jsonKeys(t *testing.T, v any) []string {
+	t.Helper()
+	b, err := json.Marshal(v)
+	if err != nil {
+		t.Fatal(err)
+	}
+	var m map[string]any
+	if err := json.Unmarshal(b, &m); err != nil {
+		t.Fatal(err)
+	}
+	var keys []string
+	for k := range m {
+		keys = append(keys, k)
+	}
+	sort.Strings(keys)
+	return keys
+}
+
+// AM1: a pending request's line is decoded from its bytes; hash16 is the CLI's message hash
+// prefix; recorded is null until an approval is recorded; a decided line truncates the credential.
+func TestApproveModelSummaries(t *testing.T) {
+	r := amRequest(t)
+	p := Pending(r, "")
+	d, _ := hex32(r.Digest)
+	h, err := approval.MessageHash(d, [32]byte(bytes.Repeat([]byte{0x5a}, 32)), amExpires)
+	if err != nil {
+		t.Fatal(err)
+	}
+	if p.ID != r.ID || p.Episode != "w201x" || p.Agent != "claude" || p.Model != "claude-opus<5>&x" || p.Grants != 3 ||
+		p.TTL == nil || *p.TTL != 43200 || p.Expires != amExpires || p.Hash16 != hex.EncodeToString(h[:8]) || p.Recorded != nil || p.Refused != "" {
+		t.Fatalf("pending = %+v", p)
+	}
+	b, _ := json.Marshal(p)
+	if !strings.Contains(string(b), `"recorded":null`) {
+		t.Fatalf("an unrecorded line is %s, want recorded null", b)
+	}
+	if got := Pending(r, "webauthn:mark"); got.Recorded == nil || *got.Recorded != "webauthn:mark" {
+		t.Fatalf("recorded = %v", got.Recorded)
+	}
+	if want := []string{"agent", "episode", "expires", "grants", "hash16", "id", "model", "recorded", "ttl_s"}; !reflect.DeepEqual(jsonKeys(t, p), want) {
+		t.Fatalf("pending keys %v, want %v (no bytes, challenge, artifact or path)", jsonKeys(t, p), want)
+	}
+	r.Consumed, r.ConsumedAt, r.ApprovedBy, r.CredentialID = true, amCreated+9, "webauthn:mark", "cred-0123456789abcdefXYZ"
+	if got := Decided(r); got != (DecidedRequest{ID: r.ID, Episode: "w201x", ApprovedBy: "webauthn:mark", At: amCreated + 9, Credential16: "cred-0123456789a"}) {
+		t.Fatalf("decided = %+v", got)
+	}
+	r.CredentialID = "short"
+	if got := Decided(r); got.Credential16 != "short" {
+		t.Fatalf("a short credential id became %q", got.Credential16)
+	}
+}
+
+// AM2: a row whose bytes do not hash to its digest, do not decode canonically, or whose challenge
+// is not hex is listed with refused and no summary, never hidden.
+func TestApproveModelRefusesWhatDoesNotDecode(t *testing.T) {
+	for name, edit := range map[string]func(*StoredRequest){
+		"digest": func(r *StoredRequest) { r.Digest = strings.Repeat("0", 64) },
+		"canonical": func(r *StoredRequest) {
+			r.Request = append(r.Request, ' ')
+			d := sha256.Sum256(r.Request)
+			r.Digest = hex.EncodeToString(d[:])
+		},
+		"challenge":       func(r *StoredRequest) { r.Challenge = "zz" },
+		"challenge-short": func(r *StoredRequest) { r.Challenge = amChallenge[2:] }, // 31 bytes
+		"expiry":          func(r *StoredRequest) { r.Expires = -1 },
+	} {
+		r := amRequest(t)
+		edit(&r)
+		p := Pending(r, "ssh:mark")
+		if p.Refused == "" || p.ID != r.ID || p.Episode != "w201x" || p.Agent != "" || p.Grants != 0 || p.TTL != nil || p.Hash16 != "" || p.Recorded == nil {
+			t.Errorf("%s: %+v, want refused with no summary", name, p)
+		}
+		if want := []string{"episode", "expires", "id", "recorded", "refused"}; !reflect.DeepEqual(jsonKeys(t, p), want) {
+			t.Errorf("%s: keys %v, want %v", name, jsonKeys(t, p), want)
+		}
+	}
+}
+
+// AM2 (the card): the stored bytes as base64url, the row's digest, nonce and times, the active
+// passkeys grouped by origin, the SSH namespace, and no message-hash field.
+func TestApproveModelCard(t *testing.T) {
+	r := amRequest(t)
+	c := RequestCard(r, []CardPasskey{{Origin: "https://s.tail0.ts.net", RPID: "s.tail0.ts.net", CredentialID: []byte{1}},
+		{Origin: "http://localhost:7644", RPID: "localhost", CredentialID: []byte{2}}, {Origin: "https://s.tail0.ts.net", RPID: "s.tail0.ts.net", CredentialID: []byte{3}}})
+	if b, err := base64.RawURLEncoding.DecodeString(c.Request); err != nil || !bytes.Equal(b, r.Request) {
+		t.Fatalf("request does not decode to the stored bytes: %v", err)
+	}
+	want := []CardPasskeys{{Origin: "https://s.tail0.ts.net", RP: "s.tail0.ts.net", Allow: []string{"AQ", "Aw"}}, {Origin: "http://localhost:7644", RP: "localhost", Allow: []string{"Ag"}}}
+	if c.V != 1 || c.ID != r.ID || c.Digest != r.Digest || c.Challenge != amChallenge || c.Created != amCreated || c.Expires != amExpires ||
+		c.Consumed || c.ApprovedBy != "" || !reflect.DeepEqual(c.Passkeys, want) || c.SSH.Namespace != "world-approval@sunholo" {
+		t.Fatalf("card = %+v", c)
+	}
+	if want := []string{"approvedBy", "challenge", "consumed", "created", "digest", "expires", "id", "passkeys", "request", "ssh", "v"}; !reflect.DeepEqual(jsonKeys(t, c), want) {
+		t.Fatalf("card keys %v, want %v (no message hash)", jsonKeys(t, c), want)
+	}
+	if c := RequestCard(r, nil); c.Passkeys == nil || len(c.Passkeys) != 0 {
+		t.Fatalf("no passkeys = %#v, want []", c.Passkeys)
+	}
+}
````

### Patch M4b — the POST-side model

````diff M4b
--- a/host/surfacepage/approve_model.go
+++ b/host/surfacepage/approve_model.go
@@ -7,11 +7,15 @@ package surfacepage
 // /surface/approve/{id} takes ApproveBody and answers ApproveRecorded or ApproveRefusal.
 
 import (
+	"bytes"
 	"crypto/sha256"
 	"encoding/base64"
 	"encoding/hex"
+	"encoding/json"
 	"errors"
 	"fmt"
+	"io"
+	"time"
 
 	"github.com/sunholo-data/ailang-world/host/approval"
 )
@@ -19,6 +23,9 @@ import (
 // MaxDecidedListed is the decision list's length (design §2.3, Q3).
 const MaxDecidedListed = 20
 
+// MaxApproveBodyBytes bounds a POST /surface/approve/{id} body (design §2.4 R6).
+const MaxApproveBodyBytes = 20 << 10
+
 // StoredRequest is one mint_requests row as the routes read it (store.MintRequestRow's fields).
 type StoredRequest struct {
 	ID, Episode, Digest, Challenge string
@@ -169,3 +176,92 @@ func RequestCard(r StoredRequest, passkeys []CardPasskey) Card {
 	}
 	return c
 }
+
+// ApproveBody is a POST /surface/approve/{id} body: v 1, the digest the page displayed, and
+// exactly one of a WebAuthn assertion envelope (A2 §2.7) and an armored SSHSIG.
+type ApproveBody struct {
+	V        int             `json:"v"`
+	Digest   string          `json:"digest"`
+	WebAuthn json.RawMessage `json:"webauthn"`
+	SSHSIG   *string         `json:"sshsig"`
+}
+
+// ErrApproveBody refuses a body that is not ApproveBody's grammar (R6).
+var ErrApproveBody = errors.New("the body is not {\"v\":1,\"digest\":<64 hex>} with exactly one of webauthn and sshsig")
+
+// ParseApproveBody is R6's grammar: strict JSON (no unknown field, exactly one value), v 1,
+// exactly one format, the digest 64 lowercase hex.
+func ParseApproveBody(b []byte) (ApproveBody, error) {
+	dec := json.NewDecoder(bytes.NewReader(b))
+	dec.DisallowUnknownFields()
+	var body ApproveBody
+	if err := dec.Decode(&body); err != nil {
+		return ApproveBody{}, fmt.Errorf("%w: %v", ErrApproveBody, err)
+	}
+	if _, err := dec.Token(); err != io.EOF {
+		return ApproveBody{}, fmt.Errorf("%w: more than one JSON value", ErrApproveBody)
+	}
+	switch {
+	case body.V != 1:
+		return ApproveBody{}, fmt.Errorf("%w: v is %d", ErrApproveBody, body.V)
+	case (body.WebAuthn != nil) == (body.SSHSIG != nil):
+		return ApproveBody{}, fmt.Errorf("%w: it has %s", ErrApproveBody, map[bool]string{true: "both", false: "neither"}[body.SSHSIG != nil])
+	case !lowerHex64(body.Digest):
+		return ApproveBody{}, fmt.Errorf("%w: digest is not 64 lowercase hex", ErrApproveBody)
+	}
+	return body, nil
+}
+
+// lowerHex64 is exactly 64 lowercase hex characters.
+func lowerHex64(s string) bool {
+	b, ok := hex32(s)
+	return ok && hex.EncodeToString(b[:]) == s
+}
+
+// AdmissionRefusal is approve's admission (approval.Admit) of r at now, in the CLI's words
+// (cmd/ailang-worldd approveAdmission); "" when r is admitted (R8).
+func AdmissionRefusal(r StoredRequest, now int64) string {
+	if _, err := approval.Decode(r.Request); err != nil {
+		return fmt.Sprintf("request %s is not a canonical mint request (%v); fix: make a new request with episode start --request", r.ID, err)
+	}
+	digestRefusal := fmt.Sprintf("request %s's stored bytes do not match its digest; refusing; fix: make a new request with episode start --request", r.ID)
+	stored, _ := hex32(r.Digest) // not 64 hex: the zero array, which no bytes hash to (the digest refusal)
+	switch approval.Admit(approval.AdmitInput{Now: now, CreatedAt: r.Created, ExpiresAt: r.Expires, Consumed: r.Consumed,
+		StoredDigest: stored, BytesDigest: sha256.Sum256(r.Request)}) {
+	case approval.Admitted:
+		return ""
+	case approval.RefuseConsumed:
+		return fmt.Sprintf("request %s was already used at %s by %s (credential %s); fix: make a new request with episode start --request",
+			r.ID, rfc3339(r.ConsumedAt), r.ApprovedBy, r.CredentialID)
+	case approval.RefuseClockBehind:
+		return fmt.Sprintf("the clock is behind request %s's creation time (%s); refusing", r.ID, rfc3339(r.Created))
+	case approval.RefuseExpired:
+		return fmt.Sprintf("request %s expired at %s; fix: make a new request with episode start --request", r.ID, rfc3339(r.Expires))
+	default:
+		return digestRefusal
+	}
+}
+
+func rfc3339(t int64) string { return time.Unix(t, 0).UTC().Format(time.RFC3339) }
+
+// ApproveRefusal is every refusal of the approval routes: inert JSON, nothing written.
+type ApproveRefusal struct {
+	V      int    `json:"v"`
+	Error  string `json:"error"`
+	Reason string `json:"reason"`
+}
+
+// ApproveRecorded is a POST's success: the verified approval is recorded for the minter.
+type ApproveRecorded struct {
+	V          int    `json:"v"`
+	Recorded   bool   `json:"recorded"`
+	Request    string `json:"request"`
+	ApprovedBy string `json:"approvedBy"`
+	Next       string `json:"next"`
+}
+
+// Recorded is the success body for request id approved by approvedBy.
+func Recorded(id, approvedBy string) ApproveRecorded {
+	return ApproveRecorded{V: 1, Recorded: true, Request: id, ApprovedBy: approvedBy,
+		Next: "the minter (approve " + id + " --await) mints it; it then appears under Decided"}
+}
--- a/host/surfacepage/approve_model_test.go
+++ b/host/surfacepage/approve_model_test.go
@@ -8,6 +8,7 @@ import (
 	"encoding/base64"
 	"encoding/hex"
 	"encoding/json"
+	"errors"
 	"reflect"
 	"sort"
 	"strings"
@@ -136,3 +137,73 @@ func TestApproveModelCard(t *testing.T) {
 		t.Fatalf("no passkeys = %#v, want []", c.Passkeys)
 	}
 }
+
+// AM3: the POST body's grammar (R6).
+func TestApproveModelBodyGrammar(t *testing.T) {
+	dg := strings.Repeat("ab", 32)
+	ok := map[string]string{
+		"webauthn": `{"v":1,"digest":"` + dg + `","webauthn":{"id":"x"}}`,
+		"sshsig":   `{"v":1,"digest":"` + dg + `","sshsig":"-----BEGIN SSH SIGNATURE-----\n"}`,
+	}
+	for name, body := range ok {
+		b, err := ParseApproveBody([]byte(body))
+		if err != nil || b.Digest != dg || (name == "webauthn") != (b.WebAuthn != nil) || (name == "sshsig") != (b.SSHSIG != nil) {
+			t.Errorf("%s: %+v, %v", name, b, err)
+		}
+	}
+	if b, _ := ParseApproveBody([]byte(ok["webauthn"])); string(b.WebAuthn) != `{"id":"x"}` {
+		t.Errorf("the envelope is %q, want the posted bytes", b.WebAuthn)
+	}
+	for name, body := range map[string]string{
+		"both":          `{"v":1,"digest":"` + dg + `","webauthn":{},"sshsig":"x"}`,
+		"neither":       `{"v":1,"digest":"` + dg + `"}`,
+		"null-and-ssh":  `{"v":1,"digest":"` + dg + `","webauthn":null,"sshsig":"x"}`,
+		"unknown-field": `{"v":1,"digest":"` + dg + `","sshsig":"x","approver":"webauthn:mark"}`,
+		"trailing":      `{"v":1,"digest":"` + dg + `","sshsig":"x"} {}`,
+		"v2":            `{"v":2,"digest":"` + dg + `","sshsig":"x"}`,
+		"v-missing":     `{"digest":"` + dg + `","sshsig":"x"}`,
+		"digest-upper":  `{"v":1,"digest":"` + strings.ToUpper(dg) + `","sshsig":"x"}`,
+		"digest-short":  `{"v":1,"digest":"` + dg[2:] + `","sshsig":"x"}`,
+		"not-json":      `v=1`,
+		"empty":         ``,
+	} {
+		if _, err := ParseApproveBody([]byte(body)); !errors.Is(err, ErrApproveBody) {
+			t.Errorf("%s: err %v, want ErrApproveBody", name, err)
+		}
+	}
+}
+
+// AM3 (admission, R8): approve's admission in the CLI's words; "" only when admitted.
+func TestApproveModelAdmission(t *testing.T) {
+	r := amRequest(t)
+	if got := AdmissionRefusal(r, amCreated); got != "" {
+		t.Fatalf("admitted at creation: %q", got)
+	}
+	if got := AdmissionRefusal(r, amExpires-1); got != "" {
+		t.Fatalf("admitted a second before expiry: %q", got)
+	}
+	for name, tc := range map[string]struct {
+		edit func(*StoredRequest)
+		now  int64
+		want string
+	}{
+		"expired": {func(*StoredRequest) {}, amExpires, "request mr-0123456789abcdef expired at 2027-01-15T08:10:00Z; fix: make a new request with episode start --request"},
+		"clock":   {func(*StoredRequest) {}, amCreated - 1, "the clock is behind request mr-0123456789abcdef's creation time (2027-01-15T08:00:00Z); refusing"},
+		"consumed": {func(r *StoredRequest) {
+			r.Consumed, r.ConsumedAt, r.ApprovedBy, r.CredentialID = true, amCreated+5, "ssh:mark", "c1"
+		}, amCreated + 9, "request mr-0123456789abcdef was already used at 2027-01-15T08:00:05Z by ssh:mark (credential c1); fix: make a new request with episode start --request"},
+		"digest":     {func(r *StoredRequest) { r.Digest = strings.Repeat("0", 64) }, amCreated, "request mr-0123456789abcdef's stored bytes do not match its digest; refusing; fix: make a new request with episode start --request"},
+		"digest-hex": {func(r *StoredRequest) { r.Digest = "zz" }, amCreated, "request mr-0123456789abcdef's stored bytes do not match its digest; refusing; fix: make a new request with episode start --request"},
+		"canonical":  {func(r *StoredRequest) { r.Request = append(r.Request, ' ') }, amCreated, "request mr-0123456789abcdef is not a canonical mint request (approval: request bytes are not canonical); fix: make a new request with episode start --request"},
+	} {
+		r := amRequest(t)
+		tc.edit(&r)
+		if got := AdmissionRefusal(r, tc.now); got != tc.want {
+			t.Errorf("%s: %q, want %q", name, got, tc.want)
+		}
+	}
+	if got := Recorded("mr-0123456789abcdef", "webauthn:mark"); got != (ApproveRecorded{V: 1, Recorded: true, Request: "mr-0123456789abcdef", ApprovedBy: "webauthn:mark",
+		Next: "the minter (approve mr-0123456789abcdef --await) mints it; it then appears under Decided"}) {
+		t.Fatalf("recorded = %+v", got)
+	}
+}
````

### Patch M4c — the origin set, the guards and --surface-origin

````diff M4c
--- a/cmd/ailang-worldd/main.go
+++ b/cmd/ailang-worldd/main.go
@@ -9,6 +9,7 @@
 //	ailang-worldd serve --db <path> [--bind 127.0.0.1:7644] [--ailang-bin <path>]
 //	                    [--workspace-root <dir> --tool-ailang-bin <path>]
 //	                    [--examples-dir <dir>] [--surface-index <file>] [--surface-follow <file>]
+//	                    [--surface-origin https://<host>.ts.net ...]
 //	                    [--workspace-module-root REL] [--workspace-episode-module-root EP=REL ...]
 //	                    [--workspace-package-cache <dir>]
 //	                    [--run-allow-caps Env,Net,Declassify]
@@ -63,6 +64,7 @@ import (
 	"strconv"
 	"strings"
 
+	"github.com/sunholo-data/ailang-world/host/approval"
 	"github.com/sunholo-data/ailang-world/host/broker"
 	"github.com/sunholo-data/ailang-world/host/daemon"
 )
@@ -85,6 +87,7 @@ Usage:
                       [--workbench-still]
                       [--workspace-root <dir> --tool-ailang-bin <path>]
                       [--examples-dir <dir>] [--surface-index <file>] [--surface-follow <file>]
+                      [--surface-origin https://<host>.ts.net ...]
                       [--workspace-module-root REL] [--workspace-episode-module-root EP=REL ...]
                       [--workspace-package-cache <dir>]
                       [--run-allow-caps Env,Net,Declassify]
@@ -282,6 +285,12 @@ serve flags:
                        world:episode:<ep>); the page re-reads
                        /surface/model.json every 15 s while visible and each
                        repo's line says when its local branch was last checked
+  --surface-origin <origin>
+                       an origin the approval routes also answer at, beside
+                       loopback and the enrolled passkeys' origins: exactly
+                       https://<host>.ts.net[:port] or http(s)://localhost[:port]
+                       (repeatable; serve refuses any other). An operator who
+                       approves with SSH through tailscale serve needs it
   --proof-key-file <file>
                        a file holding exactly 32 bytes: the MAC key the
                        workbench grade validator authenticates stored proof
@@ -520,6 +529,14 @@ func runServe(args []string, stdout, stderr io.Writer) int {
 	examplesDir := fs.String("examples-dir", "", "AILANG examples corpus for examples-search (default ~/.ailang/examples when it exists)")
 	surfaceIndex := fs.String("surface-index", "", "generated surface index outside --workspace-root (optional; may be built later)")
 	surfaceFollow := fs.String("surface-follow", "", "follow config outside --workspace-root; mutually exclusive with --surface-index")
+	var surfaceOrigins []string
+	fs.Func("surface-origin", "an origin the approval routes also answer at, https://<host>.ts.net[:port] (repeatable)", func(v string) error {
+		if _, err := approval.ParseCeremonyOrigin(v); err != nil {
+			return err
+		}
+		surfaceOrigins = append(surfaceOrigins, v)
+		return nil
+	})
 	var moduleRoot string
 	var episodeModuleRoots []string
 	fs.Func("workspace-module-root", "REL under each episode worktree that is the AILANG sandbox and module root (default .)", func(v string) error {
@@ -613,8 +630,9 @@ func runServe(args []string, stdout, stderr io.Writer) int {
 	cfg := daemon.Config{DBPath: *dbPath, CredentialsPath: *credsDB, BindHost: host, BindPort: port, AilangBin: *ailangBin,
 		WorkbenchStill: *workbenchStill,
 		WorkspaceRoot:  *workspaceRoot, ToolAilangBin: *toolAilangBin, SurfaceIndex: *surfaceIndex, SurfaceFollow: *surfaceFollow,
-		ExamplesDir:  resolveExamplesDefault(*examplesDir, os.UserHomeDir),
-		RunAllowCaps: runAllowCaps, RunNetAllow: runNetAllow, RunNetAllowHTTP: *runNetAllowHTTP,
+		SurfaceOrigins: surfaceOrigins,
+		ExamplesDir:    resolveExamplesDefault(*examplesDir, os.UserHomeDir),
+		RunAllowCaps:   runAllowCaps, RunNetAllow: runNetAllow, RunNetAllowHTTP: *runNetAllowHTTP,
 		ExecProfiles: execProfiles, ExecProfilesDir: resolveProfilesDirDefault(*execProfilesDir, execProfiles, *execSandbox),
 		ExecEpisodeProjects: execEpisodeProjects, ExecSandbox: *execSandbox,
 		ExecNode: *execNode, ExecMaxOutputBytes: execMaxOutput,
--- a/cmd/ailang-worldd/serve_surface_origin_test.go
+++ b/cmd/ailang-worldd/serve_surface_origin_test.go
@@ -0,0 +1,46 @@
+package main
+
+// Row 201 A3 sprint 2 M4c (design §2.3 [r1], §6 CF1): serve --surface-origin.
+
+import (
+	"bytes"
+	"reflect"
+	"regexp"
+	"strings"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/daemon"
+)
+
+// CF1: --surface-origin is in serve's help, repeatable, wired into Config.SurfaceOrigins (nil when
+// unset), and a value outside A2's ceremony-origin grammar ends serve with exitUsage naming the
+// flag, before the daemon is started.
+func TestServeSurfaceOriginFlag(t *testing.T) {
+	var stdout, stderr bytes.Buffer
+	if code := run([]string{"serve", "--help"}, &stdout, &stderr); code != 0 ||
+		!regexp.MustCompile(`(?m)^\s+--surface-origin <origin>$`).MatchString(stdout.String()+stderr.String()) {
+		t.Fatalf("serve --help: exit %d, lacks the --surface-origin flag", code)
+	}
+	got := captureServeConfig(t)
+	stdout.Reset()
+	stderr.Reset()
+	if code := run([]string{"serve", "--db", "/x/world.db", "--surface-origin", "https://s.tail0.ts.net", "--surface-origin", "http://localhost:7644"}, &stdout, &stderr); code != 0 {
+		t.Fatalf("serve with two origins: exit %d, stderr %q", code, stderr.String())
+	}
+	if want := []string{"https://s.tail0.ts.net", "http://localhost:7644"}; !reflect.DeepEqual(got.SurfaceOrigins, want) {
+		t.Fatalf("SurfaceOrigins = %#v, want %#v", got.SurfaceOrigins, want)
+	}
+	*got = daemon.Config{}
+	if code := run([]string{"serve", "--db", "/x/world.db"}, &stdout, &stderr); code != 0 || got.DBPath != "/x/world.db" || got.SurfaceOrigins != nil {
+		t.Fatalf("unset: exit %d, DBPath %q, SurfaceOrigins %#v; want nil", code, got.DBPath, got.SurfaceOrigins)
+	}
+	for _, bad := range []string{"https://100.64.0.1", "https://x.ts.net:443", "https://example.com", "http://s.tail0.ts.net", ""} {
+		*got = daemon.Config{}
+		stdout.Reset()
+		stderr.Reset()
+		if code := run([]string{"serve", "--db", "/x/world.db", "--surface-origin", "https://s.tail0.ts.net", "--surface-origin", bad}, &stdout, &stderr); code != exitUsage ||
+			!strings.Contains(stderr.String(), `invalid value "`+bad+`" for flag -surface-origin: approval: `) || got.DBPath != "" {
+			t.Errorf("--surface-origin %q: exit %d, daemon started %v, stderr %q; want exitUsage naming the flag before any start", bad, code, got.DBPath != "", stderr.String())
+		}
+	}
+}
--- a/host/daemon/daemon.go
+++ b/host/daemon/daemon.go
@@ -258,6 +258,10 @@ type Config struct {
 	SurfaceIndex string
 	// SurfaceFollow enrolls read-only sources; mutually exclusive with SurfaceIndex.
 	SurfaceFollow string
+	// SurfaceOrigins (`--surface-origin`, repeatable; row 201 A3) are the origins, beside
+	// loopback and the enrolled passkeys' origins, at which the approval routes answer and accept
+	// a POST (design §2.3 [r1]). Each is A2's ceremony origin; startup refuses any other.
+	SurfaceOrigins []string
 	// RunAllowCaps, RunNetAllow and RunNetAllowHTTP (`--run-allow-caps`,
 	// `--run-net-allow`, `--run-net-allow-http`; row 135 §4.3/§4.6) are the
 	// operator's ailang-run allowlist: the extra capabilities (Declassify,
@@ -710,6 +714,11 @@ func newDaemon(ctx, phaseParent context.Context, registryBudget time.Duration, c
 	if cfg.SurfaceIndex != "" && cfg.SurfaceFollow != "" {
 		return nil, &StartupError{Stage: StageConfig, Detail: "surface follow refused: --surface-follow and --surface-index are mutually exclusive; next: ailang-worldd serve --help"}
 	}
+	for _, origin := range cfg.SurfaceOrigins {
+		if err := checkSurfaceOrigin(origin); err != nil {
+			return nil, &StartupError{Stage: StageConfig, Detail: err.Error() + "; next: ailang-worldd serve --help"}
+		}
+	}
 	if !isLoopbackHost(cfg.BindHost) {
 		return nil, &StartupError{
 			Stage: StageBindPolicy,
--- a/host/daemon/surface_approve_guard.go
+++ b/host/daemon/surface_approve_guard.go
@@ -0,0 +1,91 @@
+package daemon
+
+// Row 201 A3 sprint 2 (design_docs/planned/w-world-approvals-by-key-a3.md §2.3, §2.4 R2-R5): the
+// approval routes' guards, as pure functions of the header values and the surface origin set O.
+// Only the new approval routes apply them; every existing route stays as row 199 leaves it.
+
+import (
+	"errors"
+	"fmt"
+	"mime"
+	"slices"
+	"strings"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/store"
+)
+
+// surfaceOrigins is the surface origin set O (design §2.3 [r1]): the loopback origins at port
+// (the port the request arrived on; none when unknown), every --surface-origin, and every active
+// passkey's origin. It does not depend on which format is POSTed.
+func surfaceOrigins(port string, configured, enrolled []string) []string {
+	var o []string
+	if port != "" {
+		o = append(o, "http://localhost:"+port, "http://127.0.0.1:"+port)
+	}
+	o = append(o, configured...)
+	return append(o, enrolled...)
+}
+
+// checkSurfaceOrigin is --surface-origin's grammar: A2's ceremony origin
+// (https://<host>.ts.net[:port] or http(s)://localhost[:port], no default port, no IP literal).
+func checkSurfaceOrigin(origin string) error {
+	if _, err := approval.ParseCeremonyOrigin(origin); err != nil {
+		return fmt.Errorf("--surface-origin %q: %w", origin, err)
+	}
+	return nil
+}
+
+// allowedApprovalHost reports whether host (the Host header, any case) is exactly the host[:port]
+// of an origin in o (R2: DNS rebinding sends its own name).
+func allowedApprovalHost(host string, o []string) bool {
+	host = strings.ToLower(host)
+	for _, origin := range o {
+		if _, hostPort, _ := strings.Cut(origin, "://"); host == hostPort {
+			return true
+		}
+	}
+	return false
+}
+
+// allowedApprovalOrigin reports whether origin (the Origin header) is exactly a member of o (R3);
+// a missing Origin and "null" never are.
+func allowedApprovalOrigin(origin string, o []string) bool { return slices.Contains(o, origin) }
+
+// sameOriginFetch is R4: Sec-Fetch-Site, when present, is exactly same-origin (two localhost
+// ports are same-site, design V17).
+func sameOriginFetch(site string) bool { return site == "" || site == "same-origin" }
+
+// jsonMediaType is R5: the media type is application/json, so a browser preflights every
+// cross-origin POST, and the daemon answers no preflight with CORS.
+func jsonMediaType(contentType string) bool {
+	mt, _, err := mime.ParseMediaType(contentType)
+	return err == nil && mt == "application/json"
+}
+
+// isRequestRace reports a record refused because the request changed after the checks (R14): it
+// was consumed, expired, deleted, or its clock or digest moved, between the judge and the write.
+func isRequestRace(err error) bool {
+	for _, race := range []error{store.ErrRequestConsumed, store.ErrRequestExpired, store.ErrRequestClockBehind, store.ErrRequestDigestChanged, store.ErrRequestUnknown} {
+		if errors.Is(err, race) {
+			return true
+		}
+	}
+	return false
+}
+
+// misdirectedReason is R2's refusal text; an unlisted *.ts.net Host names the setting that adds it.
+func misdirectedReason(host string) string {
+	if name, _, _ := strings.Cut(strings.ToLower(host), ":"); strings.HasSuffix(name, ".ts.net") {
+		return fmt.Sprintf("this daemon does not serve the approval surface at %s; to approve there, start serve with --surface-origin https://%s", host, strings.ToLower(host))
+	}
+	return fmt.Sprintf("this daemon does not serve the approval surface at host %q", host)
+}
+
+// crossOriginReason is R3's refusal text; a well-formed *.ts.net origin names the setting that adds it.
+func crossOriginReason(origin string) string {
+	if rp, err := approval.ParseCeremonyOrigin(origin); err == nil && strings.HasSuffix(rp, ".ts.net") {
+		return fmt.Sprintf("the origin %s is not one this daemon serves the approval surface at; to approve from it, start serve with --surface-origin %s", origin, origin)
+	}
+	return fmt.Sprintf("the origin %q is not one this daemon serves the approval surface at (loopback, a configured surface origin or an enrolled passkey's origin)", origin)
+}
--- a/host/daemon/surface_approve_guard_test.go
+++ b/host/daemon/surface_approve_guard_test.go
@@ -0,0 +1,108 @@
+package daemon
+
+// Row 201 A3 sprint 2 M4c (design §2.3, §2.4 R2-R5, §6 RG1): the approval routes' guards and the
+// surface origin set O, pure.
+
+import (
+	"errors"
+	"fmt"
+	"path/filepath"
+	"reflect"
+	"strings"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/store"
+)
+
+// RG1: O's three members; the Host, Origin, Sec-Fetch-Site and media-type predicates; the
+// --surface-origin grammar, at the predicate and at New.
+func TestApprovalGuards(t *testing.T) {
+	o := surfaceOrigins("7644", []string{"https://s.tail0.ts.net"}, []string{"https://p.tail0.ts.net:8443"})
+	if want := []string{"http://localhost:7644", "http://127.0.0.1:7644", "https://s.tail0.ts.net", "https://p.tail0.ts.net:8443"}; !reflect.DeepEqual(o, want) {
+		t.Fatalf("O = %v, want %v", o, want)
+	}
+	if got := surfaceOrigins("", nil, nil); len(got) != 0 {
+		t.Fatalf("O with no port, flag or passkey = %v, want empty", got)
+	}
+	loopbackOnly := surfaceOrigins("7644", nil, nil)
+	for _, tc := range []struct {
+		host string
+		in   []string
+		want bool
+	}{
+		{"localhost:7644", o, true}, {"127.0.0.1:7644", o, true}, {"LOCALHOST:7644", o, true},
+		{"s.tail0.ts.net", o, true}, {"S.Tail0.TS.net", o, true}, {"p.tail0.ts.net:8443", o, true},
+		{"rebind.example:7644", o, false}, {"localhost:7645", o, false}, {"localhost", o, false}, {"localhost:76440", o, false},
+		{"127.0.0.1:7644.evil", o, false}, {"s.tail0.ts.net.evil", o, false}, {"evil.s.tail0.ts.net", o, false},
+		{"p.tail0.ts.net", o, false}, {"s.tail0.ts.net:443", o, false}, {"", o, false}, {"[::1]:7644", o, false},
+		{"s.tail0.ts.net", loopbackOnly, false}, {"127.0.0.1:7644", loopbackOnly, true},
+	} {
+		if got := allowedApprovalHost(tc.host, tc.in); got != tc.want {
+			t.Errorf("Host %q in %v = %v, want %v", tc.host, tc.in, got, tc.want)
+		}
+	}
+	for _, tc := range []struct {
+		origin string
+		in     []string
+		want   bool
+	}{
+		{"http://localhost:7644", o, true}, {"http://127.0.0.1:7644", o, true}, {"https://s.tail0.ts.net", o, true}, {"https://p.tail0.ts.net:8443", o, true},
+		{"", o, false}, {"null", o, false}, {"http://localhost:7645", o, false}, {"https://localhost:7644", o, false},
+		{"http://s.tail0.ts.net", o, false}, {"https://other.tail0.ts.net", o, false}, {"HTTP://LOCALHOST:7644", o, false},
+		{"http://localhost:7644/", o, false}, {"https://p.tail0.ts.net", o, false},
+		{"https://s.tail0.ts.net", loopbackOnly, false}, {"http://127.0.0.1:7644", loopbackOnly, true},
+	} {
+		if got := allowedApprovalOrigin(tc.origin, tc.in); got != tc.want {
+			t.Errorf("Origin %q in %v = %v, want %v", tc.origin, tc.in, got, tc.want)
+		}
+	}
+	for site, want := range map[string]bool{"": true, "same-origin": true, "same-site": false, "cross-site": false, "none": false, "Same-Origin": false} {
+		if got := sameOriginFetch(site); got != want {
+			t.Errorf("Sec-Fetch-Site %q = %v, want %v", site, got, want)
+		}
+	}
+	for ct, want := range map[string]bool{"application/json": true, "application/json; charset=utf-8": true, "Application/JSON": true,
+		"text/plain": false, "": false, "application/json-patch+json": false, "multipart/form-data; boundary=x": false,
+		"application/x-www-form-urlencoded": false, "application/json;;": false} {
+		if got := jsonMediaType(ct); got != want {
+			t.Errorf("Content-Type %q = %v, want %v", ct, got, want)
+		}
+	}
+	for _, good := range []string{"https://s.tail0.ts.net", "https://s.tail0.ts.net:8443", "http://localhost:7644", "https://localhost"} {
+		if err := checkSurfaceOrigin(good); err != nil {
+			t.Errorf("--surface-origin %q refused: %v", good, err)
+		}
+	}
+	for _, bad := range []string{"https://100.64.0.1", "https://x.ts.net:443", "https://example.com", "http://s.tail0.ts.net", "https://s.tail0.ts.net/", "s.tail0.ts.net", ""} {
+		if err := checkSurfaceOrigin(bad); err == nil || !strings.Contains(err.Error(), "--surface-origin") {
+			t.Errorf("--surface-origin %q: %v, want a refusal naming the flag", bad, err)
+		}
+		_, err := New(boundedTestContext(t), Config{DBPath: filepath.Join(t.TempDir(), "world.db"), BindHost: DefaultBindHost, SurfaceOrigins: []string{"http://localhost:7644", bad}})
+		var se *StartupError
+		if !errors.As(err, &se) || se.Stage != StageConfig || !strings.Contains(se.Detail, "--surface-origin") {
+			t.Errorf("New with --surface-origin %q: %v, want a config refusal naming the flag", bad, err)
+		}
+	}
+	for _, race := range []error{store.ErrRequestConsumed, store.ErrRequestExpired, store.ErrRequestClockBehind, store.ErrRequestDigestChanged, store.ErrRequestUnknown} {
+		if !isRequestRace(fmt.Errorf("%w: mr-0123456789abcdef", race)) {
+			t.Errorf("%v is not a race (R14 409)", race)
+		}
+	}
+	if isRequestRace(store.ErrBadSurfaceApproval) || isRequestRace(errors.New("disk I/O error")) {
+		t.Error("a non-race record failure is a race")
+	}
+	if r := misdirectedReason("S.tail0.ts.net"); !strings.Contains(r, "--surface-origin https://s.tail0.ts.net") {
+		t.Errorf("an unlisted ts.net Host is refused with %q, want the fix named", r)
+	}
+	if r := misdirectedReason("rebind.example:7644"); strings.Contains(r, "--surface-origin") {
+		t.Errorf("a foreign Host is refused with %q, which names --surface-origin", r)
+	}
+	if r := crossOriginReason("https://s.tail0.ts.net"); !strings.Contains(r, "--surface-origin https://s.tail0.ts.net") {
+		t.Errorf("an unlisted ts.net Origin is refused with %q, want the fix named", r)
+	}
+	for _, other := range []string{"null", "http://localhost:7645", "https://100.64.0.1"} {
+		if r := crossOriginReason(other); strings.Contains(r, "--surface-origin") {
+			t.Errorf("Origin %q is refused with %q, which names --surface-origin", other, r)
+		}
+	}
+}
````

### Patch M4d — the read routes behind the Host check

````diff M4d
--- a/host/daemon/daemon.go
+++ b/host/daemon/daemon.go
@@ -459,6 +459,10 @@ type Daemon struct {
 	// finish-route time is ever compared with a coordinator time.
 	finishNow func() int64
 
+	// approveNow is the approval routes' clock (row 201 A3 sprint 2), Unix seconds; New sets the
+	// wall clock, and a test may pin it between synchronous handler calls.
+	approveNow func() int64
+
 	// readDeadline bounds every store read a GET handler performs. New always
 	// sets it from the readDeadline constant; it is a field rather than a
 	// direct constant reference for the same two reasons drainTimeout is —
@@ -823,6 +827,7 @@ func newDaemon(ctx, phaseParent context.Context, registryBudget time.Duration, c
 		aguiBudget: aguiRunBudget, aguiBodyBound: aguiBodyBound, aguiTick: aguiTick,
 		aguiStop: make(chan struct{}), aguiSlots: make(chan struct{}, aguiCap),
 		readDeadline: readDeadline, errLog: resolveErrorLog(cfg.ErrorLog), finishNow: func() int64 { return time.Now().Unix() },
+		approveNow:   func() int64 { return time.Now().Unix() },
 		scanPageSize: integrityScanPageSize, scanRowBudget: integrityScanRowBudget,
 		scanTimeBudget: integrityScanTimeBudget, workspace: workspace,
 		resolver: authority.New(authority.ReadOnlyAt(credPath)), credentialsPath: credPath,
@@ -1126,6 +1131,9 @@ func (d *Daemon) Handler() http.Handler {
 	mux.HandleFunc("GET /surface/fonts/{name...}", d.handleSurfaceFont)
 	mux.HandleFunc("GET /surface/model.json", d.handleSurfaceModel)
 	mux.HandleFunc("GET /surface/stack.json", d.handleSurfaceStack)
+	// The approval routes (row 201 A3): Host-checked against the surface origin set; none under /v1.
+	mux.HandleFunc("GET /surface/requests.json", d.approvalRoute(d.handleSurfaceRequests))
+	mux.HandleFunc("GET /surface/requests/{id}", d.approvalRoute(d.handleSurfaceRequestCard))
 	// The two A2A projection routes (w-a2a-session-projection P6.B-A2A-CARD)
 	// are ADDITIVE: the frozen /v1/ table above is untouched, and the routes
 	// are NOT in isProtected — the projection handler resolves the session
--- a/host/daemon/surface_app_test.go
+++ b/host/daemon/surface_app_test.go
@@ -16,7 +16,9 @@ import (
 // TestSurfaceAppRoutes: row 188 N2b. GET /surface is the app shell (the
 // same bytes whatever the index holds), its assets are the exact embedded
 // bytes, and every other name under /surface/ is a 404 (row 201 A2 adds
-// /surface/passkey and /surface/passkey.js, PG1's, and no name under them).
+// /surface/passkey and /surface/passkey.js, PG1's, and no name under them; row 201 A3
+// sprint 2 adds the approval routes /surface/requests.json and /surface/requests/{id},
+// which answer only an allowed Host, and /surface/approve.js stays a 404 until sprint 3).
 func TestSurfaceAppRoutes(t *testing.T) {
 	path := filepath.Join(t.TempDir(), "operator-index.json")
 	d, err := New(boundedTestContext(t), Config{DBPath: filepath.Join(t.TempDir(), "world.db"), BindHost: DefaultBindHost, SurfaceIndex: path})
@@ -150,10 +152,16 @@ func TestSurfaceAppRoutes(t *testing.T) {
 	t.Run("404", func(t *testing.T) {
 		for _, target := range []string{"/surface/", "/surface/nope", "/surface/app.jsx", "/surface/index.html", "/surface/data.js",
 			"/surface/spike-data.json", "/surface/app/app.js", "/surface/static/x", "/surface/fonts.css",
-			"/surface/passkey/", "/surface/passkey/index.html", "/surface/passkey/passkey.js", "/surface/passkey.html"} {
+			"/surface/passkey/", "/surface/passkey/index.html", "/surface/passkey/passkey.js", "/surface/passkey.html",
+			"/surface/approve.js", "/surface/approve", "/surface/requests", "/surface/requests.json/x"} {
 			if w := request("GET", target); w.Code != 404 {
 				t.Errorf("%s = %d, want 404", target, w.Code)
 			}
 		}
+		for _, target := range []string{"/surface/requests.json", "/surface/requests/mr-0123456789abcdef"} {
+			if w := request("GET", target); w.Code != 421 { // registered; example.com is no Host of the surface's
+				t.Errorf("%s = %d, want 421", target, w.Code)
+			}
+		}
 	})
 }
--- a/host/daemon/surface_approve.go
+++ b/host/daemon/surface_approve.go
@@ -0,0 +1,167 @@
+package daemon
+
+// Row 201 A3 sprint 2 (design_docs/planned/w-world-approvals-by-key-a3.md §2.3, §2.4): the
+// approval routes under /surface, none under /v1. GET /surface/requests.json (Needs you and the
+// decision list) and GET /surface/requests/{id} (the card's data) read the credentials store on a
+// fresh read-only handle. Each route answers only a Host in the surface origin set O (R2); every
+// existing route stays exactly as row 199 leaves it.
+
+import (
+	"context"
+	"errors"
+	"net"
+	"net/http"
+
+	"github.com/sunholo-data/ailang-world/host/store"
+	"github.com/sunholo-data/ailang-world/host/surfacepage"
+)
+
+// storedRequest is row as the approval model reads it.
+func storedRequest(row store.MintRequestRow) surfacepage.StoredRequest {
+	return surfacepage.StoredRequest{ID: row.RequestID, Episode: row.EpisodeID, Digest: row.Digest, Challenge: row.Challenge, Request: row.Request,
+		Created: row.CreatedAt, Expires: row.ExpiresAt, ConsumedAt: row.ConsumedAt, Consumed: row.Consumed, ApprovedBy: row.ApprovedBy, CredentialID: row.CredentialID}
+}
+
+// writeApprovalJSON writes an approval route's JSON: no-store, nosniff, and never a CORS header.
+func writeApprovalJSON(w http.ResponseWriter, status int, v any) {
+	w.Header().Set("Cache-Control", "no-store")
+	w.Header().Set("X-Content-Type-Options", "nosniff")
+	writeJSON(w, status, v)
+}
+
+// writeApprovalRefusal is every refusal of the approval routes: inert JSON, nothing written.
+func writeApprovalRefusal(w http.ResponseWriter, status int, class, reason string) {
+	writeApprovalJSON(w, status, surfacepage.ApproveRefusal{V: 1, Error: class, Reason: reason})
+}
+
+// approvalInternal logs err for the operator and answers the sanitized 500.
+func (d *Daemon) approvalInternal(w http.ResponseWriter, r *http.Request, err error) {
+	d.writeInternalErrorLog(r, err)
+	writeApprovalRefusal(w, http.StatusInternalServerError, "Internal", internalErrorMessage)
+}
+
+// isOlderStore reports a store below the schema a read needs (which reads as nothing there).
+func isOlderStore(err error) bool {
+	var older *store.OlderCredentialsStoreError
+	return errors.As(err, &older)
+}
+
+// approvalOrigins is O for r (design §2.3): the loopback origins at the port r arrived on (its
+// connection's local address, so an ephemeral port too), every --surface-origin, and the origin of
+// every active passkey on a fresh read-only handle. A store it cannot read adds no passkey origin.
+func (d *Daemon) approvalOrigins(ctx context.Context, r *http.Request) []string {
+	var port string
+	if addr, ok := r.Context().Value(http.LocalAddrContextKey).(net.Addr); ok {
+		_, port, _ = net.SplitHostPort(addr.String())
+	}
+	var enrolled []string
+	if ro, err := store.OpenCredentialsReadOnly(d.credentialsPath); err == nil {
+		active, _ := ro.ActivePasskeys(ctx)
+		_ = ro.Close()
+		for _, e := range active {
+			enrolled = append(enrolled, e.Passkey.Origin)
+		}
+	}
+	return surfaceOrigins(port, d.cfg.SurfaceOrigins, enrolled)
+}
+
+// approvalRoute is an approval route behind R2: a Host outside O is 421 before anything else, and
+// the handler gets O for its own checks.
+func (d *Daemon) approvalRoute(h func(http.ResponseWriter, *http.Request, []string)) http.HandlerFunc {
+	return func(w http.ResponseWriter, r *http.Request) {
+		ctx, cancel := d.readCtx(r)
+		o := d.approvalOrigins(ctx, r)
+		cancel()
+		if !allowedApprovalHost(r.Host, o) {
+			writeApprovalRefusal(w, http.StatusMisdirectedRequest, "MisdirectedRequest", misdirectedReason(r.Host))
+			return
+		}
+		h(w, r, o)
+	}
+}
+
+// handleSurfaceRequests is GET /surface/requests.json: the pending requests, each with its
+// recorded approver or null, and the last MaxDecidedListed decided ones. A store that cannot be
+// read answers 200 with unavailable and empty lists, as /surface/model.json reports its sessions.
+func (d *Daemon) handleSurfaceRequests(w http.ResponseWriter, r *http.Request, _ []string) {
+	now := d.approveNow()
+	out := surfacepage.RequestsJSON{V: 1, Now: now}
+	ctx, cancel := d.readCtx(r)
+	defer cancel()
+	if err := d.readRequests(ctx, now, &out); err != nil {
+		d.writeInternalErrorLog(r, err)
+		out.Pending, out.Decided, out.Unavailable = nil, nil, "the mint requests could not be read"
+	}
+	if out.Pending == nil {
+		out.Pending = []surfacepage.PendingRequest{}
+	}
+	if out.Decided == nil {
+		out.Decided = []surfacepage.DecidedRequest{}
+	}
+	writeApprovalJSON(w, http.StatusOK, out)
+}
+
+// readRequests fills out's lists on one read-only handle. Below schema 7 no approval is recorded.
+func (d *Daemon) readRequests(ctx context.Context, now int64, out *surfacepage.RequestsJSON) error {
+	ro, err := store.OpenCredentialsReadOnly(d.credentialsPath)
+	if err != nil {
+		return err
+	}
+	defer func() { _ = ro.Close() }()
+	pending, err := ro.PendingMintRequests(ctx, now, store.MaxPendingMintRequests)
+	if err != nil {
+		return err
+	}
+	for _, row := range pending {
+		rec, _, err := ro.SurfaceApproval(ctx, row.RequestID) // none: rec.Approver is ""
+		if err != nil && !isOlderStore(err) {
+			return err
+		}
+		out.Pending = append(out.Pending, surfacepage.Pending(storedRequest(row), rec.Approver))
+	}
+	decided, err := ro.DecidedMintRequests(ctx, "", surfacepage.MaxDecidedListed)
+	if err != nil {
+		return err
+	}
+	for _, row := range decided {
+		out.Decided = append(out.Decided, surfacepage.Decided(storedRequest(row)))
+	}
+	return nil
+}
+
+// handleSurfaceRequestCard is GET /surface/requests/{id}: the stored bytes, the digest, the raw
+// nonce and the times, and the active passkeys by origin (none below schema 6). Never a message hash.
+func (d *Daemon) handleSurfaceRequestCard(w http.ResponseWriter, r *http.Request, _ []string) {
+	id := r.PathValue("id")
+	if !store.MintRequestIDPattern.MatchString(id) {
+		writeApprovalRefusal(w, http.StatusBadRequest, "BadRequest", "a request id is mr- and 16 lowercase hex")
+		return
+	}
+	ctx, cancel := d.readCtx(r)
+	defer cancel()
+	ro, err := store.OpenCredentialsReadOnly(d.credentialsPath)
+	if err != nil {
+		d.approvalInternal(w, r, err)
+		return
+	}
+	defer func() { _ = ro.Close() }()
+	row, ok, err := ro.MintRequest(ctx, id)
+	if err != nil {
+		d.approvalInternal(w, r, err)
+		return
+	}
+	if !ok {
+		writeApprovalRefusal(w, http.StatusNotFound, "NotFound", "no mint request "+id)
+		return
+	}
+	active, err := ro.ActivePasskeys(ctx)
+	if err != nil && !isOlderStore(err) {
+		d.approvalInternal(w, r, err)
+		return
+	}
+	var passkeys []surfacepage.CardPasskey
+	for _, e := range active {
+		passkeys = append(passkeys, surfacepage.CardPasskey{Origin: e.Passkey.Origin, RPID: e.Passkey.RPID, CredentialID: e.Passkey.CredentialID})
+	}
+	writeApprovalJSON(w, http.StatusOK, surfacepage.RequestCard(storedRequest(row), passkeys))
+}
--- a/host/daemon/surface_approve_rig_test.go
+++ b/host/daemon/surface_approve_rig_test.go
@@ -0,0 +1,178 @@
+package daemon
+
+// Row 201 A3 sprint 2: the approval routes' test rig. Every store is under t.TempDir(); every key
+// is approvaltest's (D-WORLD-109); no listener: each request carries the local address
+// 127.0.0.1:7644 in its context, as net/http sets it for a connection.
+
+import (
+	"bytes"
+	"context"
+	"crypto/sha256"
+	"database/sql"
+	"encoding/hex"
+	"encoding/json"
+	"fmt"
+	"io"
+	"net"
+	"net/http"
+	"net/http/httptest"
+	"path/filepath"
+	"strings"
+	"testing"
+	"time"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
+	"github.com/sunholo-data/ailang-world/host/store"
+	"github.com/sunholo-data/ailang-world/host/surfacepage"
+)
+
+// apT0 is the approval routes' pinned clock; apLoopback the Host of a request on the loopback port.
+const apT0, apLoopback = int64(1800000000), "127.0.0.1:7644"
+
+type approvalRig struct {
+	t *testing.T
+	d *Daemon
+}
+
+func newApprovalRig(t *testing.T) approvalRig {
+	t.Helper()
+	d, err := New(boundedTestContext(t), Config{DBPath: filepath.Join(t.TempDir(), "world.db"), BindHost: DefaultBindHost, ErrorLog: io.Discard})
+	if err != nil {
+		t.Fatal(err)
+	}
+	t.Cleanup(func() { _ = d.Close() })
+	d.approveNow = func() int64 { return apT0 }
+	return approvalRig{t: t, d: d}
+}
+
+// write runs f on the credentials writer, opened and closed around it.
+func (a approvalRig) write(f func(ctx context.Context, c *store.Credentials) error) {
+	a.t.Helper()
+	c, err := store.OpenCredentials(a.d.credentialsPath)
+	if err != nil {
+		a.t.Fatal(err)
+	}
+	defer func() { _ = c.Close() }()
+	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
+	defer cancel()
+	if err := f(ctx, c); err != nil {
+		a.t.Fatal(err)
+	}
+}
+
+// putRequest stores a pending, canonical request id for ep, created at created, expiring 600 s
+// later; its challenge is the SHA-256 of its id.
+func (a approvalRig) putRequest(id, ep string, created int64) store.MintRequestRow {
+	a.t.Helper()
+	b, err := approval.Canonical(approval.MintRequest{Schema: approval.RequestSchema, Episode: ep, WorldDB: "/ws/world.db", WorkspaceRoot: "/ws/root",
+		Worktree: "/ws/root/" + ep, Git: "/usr/bin/git", Branch: "world/" + ep, Grants: []approval.Grant{{Effect: "Workspace.Read", Scope: ep}, {Effect: "Workspace.Write", Scope: ep, Budget: 5}},
+		TTLSeconds: 3600, ProfilesDir: "/ws/profiles", CredentialsDB: "/ws/credentials.db", SessionFile: "/ws/sessions/" + ep + "/session.token",
+		MCPConfig: "/ws/sessions/" + ep + "/mcp.json", Agent: "claude", Model: "claude-opus-5-5", Task: "/ws/task.md", Daemon: approval.DaemonRef{Addr: apLoopback}})
+	if err != nil {
+		a.t.Fatal(err)
+	}
+	d, c := sha256.Sum256(b), sha256.Sum256([]byte(id))
+	row := store.MintRequestRow{RequestID: id, Request: b, Digest: hex.EncodeToString(d[:]), Challenge: hex.EncodeToString(c[:]), EpisodeID: ep, CreatedAt: created, ExpiresAt: created + 600}
+	a.write(func(ctx context.Context, c *store.Credentials) error { return c.PutMintRequest(ctx, row, created) })
+	return row
+}
+
+// consume consumes row by the typed y at `at`, minting session credential cred.
+func (a approvalRig) consume(row store.MintRequestRow, at int64, cred string) {
+	a.t.Helper()
+	a.write(func(ctx context.Context, c *store.Credentials) error {
+		return c.MintSessionForRequest(ctx, store.SessionRow{CredentialID: cred, EpisodeID: row.EpisodeID, GrantsJSON: "[]", ExpiresAt: at + 3600, CreatedAt: at},
+			store.RequestConsume{ID: row.RequestID, Digest: row.Digest, Challenge: row.Challenge, ApprovedBy: "tty",
+				Clock: func() int64 { return at }, Check: func(int64) error { return nil }})
+	})
+}
+
+// enrolPasskey enrols pk as name at origin (rpID: the origin's host) at apT0-10: after the tests'
+// requests (apT0-30), so a judge at the request's creation instead of now refuses.
+func (a approvalRig) enrolPasskey(name string, pk *approvaltest.Passkey, origin string) {
+	a.t.Helper()
+	rp, err := approval.ParseCeremonyOrigin(origin)
+	if err != nil {
+		a.t.Fatal(err)
+	}
+	k := pk.Key()
+	a.write(func(ctx context.Context, c *store.Credentials) error {
+		return c.EnrolPasskey(ctx, store.ApproverRow{Name: name, Format: approval.FormatWebAuthn, PublicKey: k.Canonical(), EnrolledAt: apT0 - 10},
+			store.PasskeyRow{Name: name, RPID: rp, Origin: origin, CredentialID: pk.CredentialID(), AAGUID: make([]byte, 16), Alg: k.Alg, SignCount: 1})
+	})
+}
+
+// serve sends one request through the daemon's handler as if it arrived on 127.0.0.1:7644, with
+// Host host and the header pairs hdr.
+func (a approvalRig) serve(method, target, host string, body []byte, hdr ...string) *httptest.ResponseRecorder {
+	a.t.Helper()
+	req := httptest.NewRequest(method, target, bytes.NewReader(body))
+	req = req.WithContext(context.WithValue(req.Context(), http.LocalAddrContextKey, net.Addr(&net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 7644})))
+	req.Host = host
+	for i := 0; i+1 < len(hdr); i += 2 {
+		req.Header.Set(hdr[i], hdr[i+1])
+	}
+	rec := httptest.NewRecorder()
+	a.d.Handler().ServeHTTP(rec, req)
+	return rec
+}
+
+// approvalHeaders are an approval route's headers: JSON, no-store, nosniff, and no CORS header.
+func approvalHeaders(rec *httptest.ResponseRecorder) error {
+	for k, v := range map[string]string{"Content-Type": "application/json", "Cache-Control": "no-store", "X-Content-Type-Options": "nosniff"} {
+		if got := rec.Header().Get(k); got != v {
+			return fmt.Errorf("%s = %q, want %q", k, got, v)
+		}
+	}
+	for k := range rec.Header() {
+		if strings.HasPrefix(strings.ToLower(k), "access-control-") {
+			return fmt.Errorf("a CORS header %s", k)
+		}
+	}
+	return nil
+}
+
+// requireRefusal fails unless rec is an inert refusal with status and class.
+func requireRefusal(t *testing.T, what string, rec *httptest.ResponseRecorder, status int, class string) surfacepage.ApproveRefusal {
+	t.Helper()
+	var r surfacepage.ApproveRefusal
+	if err := json.Unmarshal(rec.Body.Bytes(), &r); err != nil || rec.Code != status || r.V != 1 || r.Error != class || r.Reason == "" {
+		t.Fatalf("%s: %d %s, want %d %s", what, rec.Code, rec.Body.String(), status, class)
+	}
+	if err := approvalHeaders(rec); err != nil {
+		t.Fatalf("%s: %v", what, err)
+	}
+	return r
+}
+
+// setRequestColumn sets one column of request id straight in the credentials file (a hand edit,
+// or a row no store method writes).
+func setRequestColumn(t *testing.T, path, id, column string, v any) {
+	t.Helper()
+	db, err := sql.Open("sqlite", path)
+	if err != nil {
+		t.Fatal(err)
+	}
+	defer func() { _ = db.Close() }()
+	if res, err := db.Exec(`UPDATE mint_requests SET `+column+` = ? WHERE request_id = ?`, v, id); err != nil {
+		t.Fatal(err)
+	} else if n, _ := res.RowsAffected(); n != 1 {
+		t.Fatalf("%d rows updated", n)
+	}
+}
+
+// setVersion relabels the credentials store at path to an older schema by hand (DROP TABLE, PRAGMA).
+func setVersion(t *testing.T, path string, statements ...string) {
+	t.Helper()
+	db, err := sql.Open("sqlite", path)
+	if err != nil {
+		t.Fatal(err)
+	}
+	defer func() { _ = db.Close() }()
+	for _, s := range statements {
+		if _, err := db.Exec(s); err != nil {
+			t.Fatal(err)
+		}
+	}
+}
--- a/host/daemon/surface_approve_test.go
+++ b/host/daemon/surface_approve_test.go
@@ -0,0 +1,99 @@
+package daemon
+
+// Row 201 A3 sprint 2 M4d (design §2.3, §6 RT1, RT7): the approval routes and their Host check.
+
+import (
+	"net/http/httptest"
+	"os"
+	"reflect"
+	"regexp"
+	"sort"
+	"strings"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
+)
+
+// approvalRoutes is every approval route Handler registers this sprint (the card page and its
+// script are sprint 3's).
+var approvalRoutes = []string{"GET /surface/requests.json", "GET /surface/requests/{id}"}
+
+// RT1: exactly the approval routes are registered, under /surface, each behind approvalRoute; no
+// other method is answered on them; the card page and its script are not served yet.
+func TestApprovalRoutesRegistered(t *testing.T) {
+	src, err := os.ReadFile("daemon.go")
+	if err != nil {
+		t.Fatal(err)
+	}
+	body := strings.SplitN(strings.SplitN(string(src), "func (d *Daemon) Handler() http.Handler {", 2)[1], "\n}\n", 2)[0]
+	var got []string
+	for _, m := range regexp.MustCompile(`mux\.HandleFunc\("([A-Z]+ /surface/(?:requests|approve)[^"]*)", (d\.approvalRoute\()?`).FindAllStringSubmatch(body, -1) {
+		if m[2] == "" {
+			t.Errorf("%s is not behind approvalRoute", m[1])
+		}
+		got = append(got, m[1])
+	}
+	sort.Strings(got)
+	if want := append([]string(nil), approvalRoutes...); !reflect.DeepEqual(got, want) {
+		t.Fatalf("approval routes %v, want %v", got, want)
+	}
+	a := newApprovalRig(t)
+	if rec := a.serve("GET", "/surface/requests.json", apLoopback, nil); rec.Code != 200 {
+		t.Fatalf("requests.json = %d", rec.Code)
+	}
+	requireRefusal(t, "unknown card", a.serve("GET", "/surface/requests/mr-0123456789abcdef", apLoopback, nil), 404, "NotFound")
+	for _, target := range []string{"/surface/requests.json", "/surface/requests/mr-0123456789abcdef"} {
+		for _, method := range []string{"POST", "PUT", "DELETE", "OPTIONS"} {
+			if rec := a.serve(method, target, apLoopback, nil, "Origin", "http://localhost:7644"); rec.Code != 405 || rec.Header().Get("Access-Control-Allow-Origin") != "" {
+				t.Errorf("%s %s = %d, want 405 with no CORS header", method, target, rec.Code)
+			}
+		}
+	}
+	for _, target := range []string{"/surface/approve.js", "/surface/approve", "/surface/requests", "/surface/requests/mr-0123456789abcdef/x"} {
+		if rec := a.serve("GET", target, apLoopback, nil); rec.Code != 404 {
+			t.Errorf("GET %s = %d, want 404 (not served this sprint)", target, rec.Code)
+		}
+	}
+}
+
+// RT7: a Host outside O is 421 on every approval route; loopback at the arriving port, an enrolled
+// passkey's host and a configured origin's host pass; with no local address there is no loopback
+// member. Control: the existing routes still answer a foreign Host (row 199's posture).
+func TestApprovalRoutesRefuseForeignHost(t *testing.T) {
+	a := newApprovalRig(t)
+	row := a.putRequest("mr-00000000000000b1", "w201b", apT0-30)
+	targets := []string{"/surface/requests.json", "/surface/requests/" + row.RequestID}
+	for _, target := range targets {
+		for _, host := range []string{"rebind.example:7644", "localhost:7645", "localhost", "127.0.0.1:7644.evil", "studio.tail0.ts.net", "cfg.tail0.ts.net"} {
+			requireRefusal(t, host+target, a.serve("GET", target, host, nil), 421, "MisdirectedRequest")
+		}
+		for _, host := range []string{"127.0.0.1:7644", "localhost:7644", "LOCALHOST:7644"} {
+			if rec := a.serve("GET", target, host, nil); rec.Code != 200 {
+				t.Errorf("GET %s with Host %s = %d, want 200", target, host, rec.Code)
+			}
+		}
+	}
+	if r := requireRefusal(t, "ts.net", a.serve("GET", targets[0], "studio.tail0.ts.net", nil), 421, "MisdirectedRequest"); !strings.Contains(r.Reason, "--surface-origin https://studio.tail0.ts.net") {
+		t.Errorf("the unlisted ts.net Host's refusal %q does not name the fix", r.Reason)
+	}
+	a.enrolPasskey("b-phone", approvaltest.NewPasskey(approval.AlgES256, 62), "https://studio.tail0.ts.net")
+	a.d.cfg.SurfaceOrigins = []string{"https://cfg.tail0.ts.net"}
+	for _, target := range targets {
+		for _, host := range []string{"studio.tail0.ts.net", "cfg.tail0.ts.net"} {
+			if rec := a.serve("GET", target, host, nil); rec.Code != 200 {
+				t.Errorf("GET %s with Host %s = %d, want 200", target, host, rec.Code)
+			}
+		}
+	}
+	rec := httptest.NewRecorder() // no local address: no loopback member
+	req := httptest.NewRequest("GET", targets[0], nil)
+	req.Host = apLoopback
+	a.d.Handler().ServeHTTP(rec, req)
+	requireRefusal(t, "no local address", rec, 421, "MisdirectedRequest")
+	for _, target := range []string{"/surface/model.json", "/surface/passkey", "/surface/app.js"} {
+		if rec := a.serve("GET", target, "rebind.example:7644", nil); rec.Code != 200 {
+			t.Errorf("control: GET %s with a foreign Host = %d, want 200 (existing routes unchanged)", target, rec.Code)
+		}
+	}
+}
````

### Patch M4e — what the read routes say

````diff M4e
--- a/host/daemon/surface_approve_read_test.go
+++ b/host/daemon/surface_approve_read_test.go
@@ -0,0 +1,139 @@
+package daemon
+
+// Row 201 A3 sprint 2 M4e (design §2.3, §6 RT2, RT3): what the read routes say.
+
+import (
+	"bytes"
+	"context"
+	"encoding/base64"
+	"encoding/hex"
+	"encoding/json"
+	"os"
+	"reflect"
+	"sort"
+	"strings"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
+	"github.com/sunholo-data/ailang-world/host/store"
+	"github.com/sunholo-data/ailang-world/host/surfacepage"
+)
+
+// RT2: requests.json lists the pending requests newest first with their recorded approver or
+// null, a row that does not decode as refused, and the decided ones (at most 20, newest first);
+// it carries no request bytes, challenge, artifact or token; a store it cannot read is 200 with
+// unavailable and empty lists.
+func TestRequestsJSON(t *testing.T) {
+	a := newApprovalRig(t)
+	a.putRequest("mr-00000000000000a1", "w201a", apT0-700) // expired at apT0-100: not listed
+	r1 := a.putRequest("mr-00000000000000b1", "w201b", apT0-30)
+	r2 := a.putRequest("mr-00000000000000c1", "w201c", apT0-20)
+	bad := a.putRequest("mr-00000000000000d1", "w201d", apT0-10)
+	artifact := []byte("ARTIFACT-NEVER-LISTED")
+	a.write(func(ctx context.Context, c *store.Credentials) error {
+		return c.RecordSurfaceApproval(ctx, store.SurfaceApprovalRow{RequestID: r2.RequestID, Format: "webauthn", Approver: "webauthn:mark", Artifact: artifact, Digest: r2.Digest}, apT0-5)
+	})
+	setRequestColumn(t, a.d.credentialsPath, bad.RequestID, "request", append(bad.Request, ' ')) // no longer hashes to its digest
+	for i := 0; i < 21; i++ {
+		row := a.putRequest("mr-0000000000000e"+hex.EncodeToString([]byte{byte(i)}), "w201e", apT0-1000+int64(i))
+		a.consume(row, apT0-900+int64(i), strings.Repeat(hex.EncodeToString([]byte{byte(i)}), 32))
+	}
+	rec := a.serve("GET", "/surface/requests.json", apLoopback, nil)
+	if err := approvalHeaders(rec); rec.Code != 200 || err != nil {
+		t.Fatalf("requests.json = %d, %v", rec.Code, err)
+	}
+	var got surfacepage.RequestsJSON
+	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil {
+		t.Fatal(err)
+	}
+	if got.V != 1 || got.Now != apT0 || got.Unavailable != "" || len(got.Pending) != 3 || len(got.Decided) != surfacepage.MaxDecidedListed {
+		t.Fatalf("requests.json = %+v", got)
+	}
+	p := got.Pending
+	if p[0].ID != bad.RequestID || p[0].Refused == "" || p[0].Agent != "" || p[1].ID != r2.RequestID || p[1].Recorded == nil || *p[1].Recorded != "webauthn:mark" ||
+		p[2].ID != r1.RequestID || p[2].Recorded != nil || p[2].Grants != 2 || p[2].Agent != "claude" || p[2].Hash16 == "" || p[2].Expires != apT0+570 {
+		t.Fatalf("pending = %+v", p)
+	}
+	if d := got.Decided; d[0].ID != "mr-0000000000000e14" || d[0].ApprovedBy != "tty" || d[0].At != apT0-880 || d[0].Credential16 != "1414141414141414" || d[19].ID != "mr-0000000000000e01" {
+		t.Fatalf("decided = %+v … %+v", d[0], d[19])
+	}
+	for _, leak := range []string{`"request"`, `"challenge"`, `"artifact"`, `"token"`, string(artifact), r1.Challenge, r1.Digest, "/ws/", "session"} {
+		if strings.Contains(rec.Body.String(), leak) {
+			t.Errorf("requests.json carries %q", leak)
+		}
+	}
+	setVersion(t, a.d.credentialsPath, `DROP TABLE surface_approvals`, `PRAGMA user_version = 6`) // below 7: listed, nothing recorded
+	if body := a.serve("GET", "/surface/requests.json", apLoopback, nil).Body.String(); strings.Count(body, `"recorded":null`) != 3 || !strings.Contains(body, `"unavailable":""`) {
+		t.Fatalf("a schema-6 store: %s", body)
+	}
+	if err := os.Remove(a.d.credentialsPath); err != nil {
+		t.Fatal(err)
+	}
+	rec = a.serve("GET", "/surface/requests.json", apLoopback, nil)
+	if rec.Code != 200 || !strings.Contains(rec.Body.String(), `"pending":[],"decided":[],"unavailable":"the mint requests could not be read"`) {
+		t.Fatalf("a missing store: %d %s", rec.Code, rec.Body.String())
+	}
+}
+
+// RT3: the card is the stored bytes (base64url), the row's digest, nonce and times, and the active
+// passkeys by origin (a revoked one left out); it has no message-hash field. A bad id is 400, an
+// unknown one 404, a consumed one says so.
+func TestRequestCardJSON(t *testing.T) {
+	a := newApprovalRig(t)
+	row := a.putRequest("mr-00000000000000b1", "w201b", apT0-30)
+	a.enrolPasskey("a-mac", approvaltest.NewPasskey(approval.AlgEdDSA, 61), "http://localhost:7644")
+	a.enrolPasskey("b-phone", approvaltest.NewPasskey(approval.AlgES256, 62), "https://studio.tail0.ts.net")
+	a.enrolPasskey("c-key", approvaltest.NewPasskey(approval.AlgEdDSA, 63), "http://localhost:7644")
+	a.enrolPasskey("d-gone", approvaltest.NewPasskey(approval.AlgEdDSA, 64), "https://gone.tail0.ts.net")
+	a.write(func(ctx context.Context, c *store.Credentials) error { return c.RevokeApprover(ctx, "d-gone", apT0-5) })
+	rec := a.serve("GET", "/surface/requests/"+row.RequestID, apLoopback, nil)
+	if err := approvalHeaders(rec); rec.Code != 200 || err != nil {
+		t.Fatalf("card = %d, %v", rec.Code, err)
+	}
+	var keys map[string]json.RawMessage
+	if err := json.Unmarshal(rec.Body.Bytes(), &keys); err != nil {
+		t.Fatal(err)
+	}
+	var names []string
+	for k := range keys {
+		names = append(names, k)
+	}
+	sort.Strings(names)
+	if want := []string{"approvedBy", "challenge", "consumed", "created", "digest", "expires", "id", "passkeys", "request", "ssh", "v"}; !reflect.DeepEqual(names, want) {
+		t.Fatalf("card keys %v, want %v", names, want)
+	}
+	var c surfacepage.Card
+	if err := json.Unmarshal(rec.Body.Bytes(), &c); err != nil {
+		t.Fatal(err)
+	}
+	b, err := base64.RawURLEncoding.DecodeString(c.Request)
+	b64 := func(pk *approvaltest.Passkey) string { return base64.RawURLEncoding.EncodeToString(pk.CredentialID()) }
+	want := []surfacepage.CardPasskeys{{Origin: "http://localhost:7644", RP: "localhost", Allow: []string{b64(approvaltest.NewPasskey(approval.AlgEdDSA, 61)), b64(approvaltest.NewPasskey(approval.AlgEdDSA, 63))}},
+		{Origin: "https://studio.tail0.ts.net", RP: "studio.tail0.ts.net", Allow: []string{b64(approvaltest.NewPasskey(approval.AlgES256, 62))}}}
+	if err != nil || !bytes.Equal(b, row.Request) || c.Digest != row.Digest || c.Challenge != row.Challenge || c.Created != row.CreatedAt ||
+		c.Expires != row.ExpiresAt || c.Consumed || !reflect.DeepEqual(c.Passkeys, want) || c.SSH.Namespace != approval.Namespace {
+		t.Fatalf("card = %+v", c)
+	}
+	d, _ := hex.DecodeString(row.Digest)
+	n, _ := hex.DecodeString(row.Challenge)
+	h, err := approval.MessageHash([32]byte(d), [32]byte(n), row.ExpiresAt)
+	if err != nil || strings.Contains(rec.Body.String(), hex.EncodeToString(h[:])) || strings.Contains(rec.Body.String(), base64.RawURLEncoding.EncodeToString(h[:])) {
+		t.Fatalf("the card carries the message hash (%v)", err)
+	}
+	requireRefusal(t, "bad id", a.serve("GET", "/surface/requests/MR-00000000000000B1", apLoopback, nil), 400, "BadRequest")
+	requireRefusal(t, "short id", a.serve("GET", "/surface/requests/mr-00", apLoopback, nil), 400, "BadRequest")
+	requireRefusal(t, "unknown", a.serve("GET", "/surface/requests/mr-00000000000000f1", apLoopback, nil), 404, "NotFound")
+	a.consume(row, apT0-10, strings.Repeat("ab", 32))
+	if err := json.Unmarshal(a.serve("GET", "/surface/requests/"+row.RequestID, apLoopback, nil).Body.Bytes(), &c); err != nil || !c.Consumed || c.ApprovedBy != "tty" {
+		t.Fatalf("consumed card = %+v, %v", c, err)
+	}
+	setVersion(t, a.d.credentialsPath, `DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `PRAGMA user_version = 5`)
+	if rec := a.serve("GET", "/surface/requests/"+row.RequestID, apLoopback, nil); rec.Code != 200 || !strings.Contains(rec.Body.String(), `"passkeys":[]`) {
+		t.Fatalf("a schema-5 store's card = %d %s, want 200 with no passkeys", rec.Code, rec.Body.String())
+	}
+	if err := os.Remove(a.d.credentialsPath); err != nil {
+		t.Fatal(err)
+	}
+	requireRefusal(t, "no store", a.serve("GET", "/surface/requests/"+row.RequestID, apLoopback, nil), 500, "Internal")
+}
````

### Patch M5a — the POST

````diff M5a
--- a/host/daemon/daemon.go
+++ b/host/daemon/daemon.go
@@ -1134,6 +1134,7 @@ func (d *Daemon) Handler() http.Handler {
 	// The approval routes (row 201 A3): Host-checked against the surface origin set; none under /v1.
 	mux.HandleFunc("GET /surface/requests.json", d.approvalRoute(d.handleSurfaceRequests))
 	mux.HandleFunc("GET /surface/requests/{id}", d.approvalRoute(d.handleSurfaceRequestCard))
+	mux.HandleFunc("POST /surface/approve/{id}", d.approvalRoute(d.handleSurfaceApprove))
 	// The two A2A projection routes (w-a2a-session-projection P6.B-A2A-CARD)
 	// are ADDITIVE: the frozen /v1/ table above is untouched, and the routes
 	// are NOT in isProtected — the projection handler resolves the session
--- a/host/daemon/surface_approve.go
+++ b/host/daemon/surface_approve.go
@@ -3,15 +3,21 @@ package daemon
 // Row 201 A3 sprint 2 (design_docs/planned/w-world-approvals-by-key-a3.md §2.3, §2.4): the
 // approval routes under /surface, none under /v1. GET /surface/requests.json (Needs you and the
 // decision list) and GET /surface/requests/{id} (the card's data) read the credentials store on a
-// fresh read-only handle. Each route answers only a Host in the surface origin set O (R2); every
+// fresh read-only handle; POST /surface/approve/{id} records a passkey approval only after the
+// A2 law verified it. Each route answers only a Host in the surface origin set O (R2); every
 // existing route stays exactly as row 199 leaves it.
 
 import (
 	"context"
+	"encoding/hex"
 	"errors"
+	"fmt"
+	"io"
 	"net"
 	"net/http"
 
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/authority"
 	"github.com/sunholo-data/ailang-world/host/store"
 	"github.com/sunholo-data/ailang-world/host/surfacepage"
 )
@@ -165,3 +171,121 @@ func (d *Daemon) handleSurfaceRequestCard(w http.ResponseWriter, r *http.Request
 	}
 	writeApprovalJSON(w, http.StatusOK, surfacepage.RequestCard(storedRequest(row), passkeys))
 }
+
+// approvalStoreError answers a store failure of the POST: a store below the schema a step needs is
+// 503 naming the restart, a busy writer 503 Busy with Retry-After (never a wait), anything else
+// the sanitized 500.
+func (d *Daemon) approvalStoreError(w http.ResponseWriter, r *http.Request, err error) {
+	switch {
+	case isOlderStore(err):
+		writeApprovalRefusal(w, http.StatusServiceUnavailable, "Restart",
+			"the credentials store predates surface approvals: restart `ailang-worldd serve` with this build once (it upgrades the store), then approve again")
+	case store.IsWriterAlreadyActive(err):
+		w.Header().Set("Retry-After", "1")
+		writeApprovalRefusal(w, http.StatusServiceUnavailable, "Busy", "the credentials store is busy (a mint or a session verb is writing it); try again")
+	default:
+		d.approvalInternal(w, r, err)
+	}
+}
+
+// handleSurfaceApprove is POST /surface/approve/{id} (design §2.4): after R2 (approvalRoute), the
+// checks R1, R3-R12 in order on a read-only handle, each refusal writing nothing; then R13 a
+// non-waiting writer, opened for R14's RecordSurfaceApproval and closed at once. It reads no
+// Cookie and no Authorization, answers no CORS, and never mints: the minter (approve --await)
+// judges the recorded artifact again.
+func (d *Daemon) handleSurfaceApprove(w http.ResponseWriter, r *http.Request, o []string) {
+	id := r.PathValue("id")
+	if !store.MintRequestIDPattern.MatchString(id) { // R1
+		writeApprovalRefusal(w, http.StatusBadRequest, "BadRequest", "a request id is mr- and 16 lowercase hex")
+		return
+	}
+	if origin := r.Header.Get("Origin"); !allowedApprovalOrigin(origin, o) { // R3
+		writeApprovalRefusal(w, http.StatusForbidden, "CrossOrigin", crossOriginReason(origin))
+		return
+	}
+	if !sameOriginFetch(r.Header.Get("Sec-Fetch-Site")) { // R4
+		writeApprovalRefusal(w, http.StatusForbidden, "CrossOrigin", "Sec-Fetch-Site is not same-origin")
+		return
+	}
+	if !jsonMediaType(r.Header.Get("Content-Type")) { // R5
+		writeApprovalRefusal(w, http.StatusUnsupportedMediaType, "UnsupportedMediaType", "the body must be application/json")
+		return
+	}
+	raw, err := io.ReadAll(http.MaxBytesReader(w, r.Body, surfacepage.MaxApproveBodyBytes)) // R6
+	if err != nil {
+		writeApprovalRefusal(w, http.StatusRequestEntityTooLarge, "TooLarge", fmt.Sprintf("the body is over %d bytes", surfacepage.MaxApproveBodyBytes))
+		return
+	}
+	body, err := surfacepage.ParseApproveBody(raw)
+	if err != nil {
+		writeApprovalRefusal(w, http.StatusBadRequest, "BadRequest", err.Error())
+		return
+	}
+	now := d.approveNow()
+	ctx, cancel := d.readCtx(r)
+	defer cancel()
+	ro, err := store.OpenCredentialsReadOnly(d.credentialsPath)
+	if err != nil {
+		d.approvalInternal(w, r, err)
+		return
+	}
+	defer func() { _ = ro.Close() }()
+	row, ok, err := ro.MintRequest(ctx, id) // R7 (serve upgraded the store to 7 at start)
+	if err != nil {
+		d.approvalInternal(w, r, err)
+		return
+	}
+	if !ok {
+		writeApprovalRefusal(w, http.StatusNotFound, "NotFound", "no mint request "+id)
+		return
+	}
+	if reason := surfacepage.AdmissionRefusal(storedRequest(row), now); reason != "" { // R8
+		writeApprovalRefusal(w, http.StatusConflict, "Conflict", reason)
+		return
+	}
+	if body.Digest != row.Digest { // R9
+		writeApprovalRefusal(w, http.StatusConflict, "Conflict", "the request changed since the card was shown; reload it")
+		return
+	}
+	challenge, _ := hex.DecodeString(row.Challenge) // R10 (the table's CHECK holds it to 64 characters)
+	if len(challenge) != 32 {
+		writeApprovalRefusal(w, http.StatusConflict, "Conflict", "request "+id+"'s stored challenge is not 64 hex; make a new request")
+		return
+	}
+	if body.SSHSIG != nil { // R11, sprint 2's fallback cut: the SSHSIG arm is sprint 3's
+		writeApprovalRefusal(w, http.StatusUnprocessableEntity, "NotApproved",
+			"this daemon records only passkey approvals from the surface; approve an SSH signature at a terminal: ailang-worldd approve "+id+" --signature <file>")
+		return
+	}
+	signed, err := authority.JudgeAssertion(ctx, ro, row, body.WebAuthn, now) // R11: the law
+	if err != nil {
+		var je *authority.JudgeError
+		if errors.As(err, &je) && je.Step == authority.JudgeStepStore {
+			d.approvalStoreError(w, r, err)
+			return
+		}
+		writeApprovalRefusal(w, http.StatusUnprocessableEntity, "NotApproved", "the passkey's assertion does not approve request "+id+": "+err.Error())
+		return
+	}
+	digest, _ := hex.DecodeString(row.Digest)                                                  // R8 has shown it is 64 hex
+	if err := signed.Binds([32]byte(digest), [32]byte(challenge), row.ExpiresAt); err != nil { // R12
+		writeApprovalRefusal(w, http.StatusUnprocessableEntity, "NotApproved", "the verified approval is not bound to request "+id+": "+err.Error())
+		return
+	}
+	wr, err := store.OpenCredentials(d.credentialsPath) // R13: non-waiting, closed at once
+	if err != nil {
+		d.approvalStoreError(w, r, err)
+		return
+	}
+	err = wr.RecordSurfaceApproval(ctx, store.SurfaceApprovalRow{RequestID: id, Format: approval.FormatWebAuthn, Approver: signed.ApprovedBy(),
+		Artifact: body.WebAuthn, Digest: row.Digest}, d.approveNow()) // R14, at the time it is received
+	_ = wr.Close()
+	switch {
+	case isRequestRace(err):
+		writeApprovalRefusal(w, http.StatusConflict, "Conflict", "request "+id+" changed while it was approved: "+err.Error())
+	case err != nil:
+		d.approvalStoreError(w, r, err)
+	default:
+		writeApprovalJSON(w, http.StatusOK, surfacepage.Recorded(id, signed.ApprovedBy()))
+	}
+}
--- a/host/daemon/surface_approve_post_test.go
+++ b/host/daemon/surface_approve_post_test.go
@@ -0,0 +1,220 @@
+package daemon
+
+// Row 201 A3 sprint 2 M5a (design §2.4, §6 RT4, RT5): POST /surface/approve/{id} judges with the
+// A2 law before it records, and each refusal writes nothing.
+
+import (
+	"bytes"
+	"context"
+	"database/sql"
+	"encoding/hex"
+	"encoding/json"
+	"net/http/httptest"
+	"strings"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
+	"github.com/sunholo-data/ailang-world/host/store"
+	"github.com/sunholo-data/ailang-world/host/surfacepage"
+)
+
+// apPasskey is the enrolled passkey (mark) of the POST tests, at the loopback origin apOrigin.
+var apPasskey = approvaltest.NewPasskey(approval.AlgEdDSA, 71)
+
+const apOrigin = "http://localhost:7644"
+
+// newPostRig is a rig with apPasskey enrolled as mark and one pending request.
+func newPostRig(t *testing.T) (approvalRig, store.MintRequestRow) {
+	t.Helper()
+	a := newApprovalRig(t)
+	a.enrolPasskey("mark", apPasskey, apOrigin)
+	return a, a.putRequest("mr-00000000000000b1", "w201b", apT0-30)
+}
+
+// assertionFor is pk's assertion (counter n) at apOrigin over row's message.
+func assertionFor(t *testing.T, pk *approvaltest.Passkey, row store.MintRequestRow, n uint32) []byte {
+	t.Helper()
+	d, _ := hex.DecodeString(row.Digest)
+	c, _ := hex.DecodeString(row.Challenge)
+	h, err := approval.MessageHash([32]byte(d), [32]byte(c), row.ExpiresAt)
+	if err != nil {
+		t.Fatal(err)
+	}
+	return pk.Assert(apOrigin, "localhost", h, n)
+}
+
+// approveBody is the page's POST body for digest with a webauthn assertion.
+func approveBody(digest string, assertion []byte) []byte {
+	b, _ := json.Marshal(map[string]any{"v": 1, "digest": digest, "webauthn": json.RawMessage(assertion)})
+	return b
+}
+
+// post POSTs body to /surface/approve/<id> as the card page would (loopback Host and Origin,
+// same-origin, JSON); hdr pairs override those headers ("" removes one).
+func (a approvalRig) post(id string, body []byte, hdr ...string) *httptest.ResponseRecorder {
+	a.t.Helper()
+	return a.serve("POST", "/surface/approve/"+id, apLoopback, body,
+		append([]string{"Origin", apOrigin, "Sec-Fetch-Site", "same-origin", "Content-Type", "application/json"}, hdr...)...)
+}
+
+// recordedAndSessions counts surface_approvals and session_credentials rows (straight from the file).
+func (a approvalRig) recordedAndSessions() (recorded, sessions int) {
+	a.t.Helper()
+	db, err := sql.Open("sqlite", a.d.credentialsPath)
+	if err != nil {
+		a.t.Fatal(err)
+	}
+	defer func() { _ = db.Close() }()
+	_ = db.QueryRow(`SELECT count(*) FROM surface_approvals`).Scan(&recorded) // below schema 7: 0
+	if err := db.QueryRow(`SELECT count(*) FROM session_credentials`).Scan(&sessions); err != nil {
+		a.t.Fatal(err)
+	}
+	return recorded, sessions
+}
+
+// RT4: R1-R14, each refusal with its status and class, and 0 surface_approvals rows and 0
+// sessions after it. R2's 421 and R12 are RT7's and the law's (named equivalent, plan §4).
+func TestApprovePOSTRefusalArms(t *testing.T) {
+	other := approvaltest.NewPasskey(approval.AlgES256, 72)
+	for _, arm := range []struct {
+		name, id, class, reason string
+		status                  int
+		setup                   func(a approvalRig, row store.MintRequestRow) (body []byte, hdr []string)
+	}{
+		{"R1-id", "MR-00000000000000B1", "BadRequest", "mr- and 16 lowercase hex", 400, nil},
+		{"R3-origin", "", "CrossOrigin", "http://localhost:7645", 403, func(a approvalRig, row store.MintRequestRow) ([]byte, []string) {
+			return approveBody(row.Digest, assertionFor(t, apPasskey, row, 2)), []string{"Origin", "http://localhost:7645"}
+		}},
+		{"R4-site", "", "CrossOrigin", "Sec-Fetch-Site", 403, func(a approvalRig, row store.MintRequestRow) ([]byte, []string) {
+			return approveBody(row.Digest, assertionFor(t, apPasskey, row, 2)), []string{"Sec-Fetch-Site", "same-site"}
+		}},
+		{"R5-type", "", "UnsupportedMediaType", "application/json", 415, func(a approvalRig, row store.MintRequestRow) ([]byte, []string) {
+			return approveBody(row.Digest, assertionFor(t, apPasskey, row, 2)), []string{"Content-Type", "text/plain"}
+		}},
+		{"R6-size", "", "TooLarge", "20480", 413, func(a approvalRig, row store.MintRequestRow) ([]byte, []string) {
+			return approveBody(row.Digest, []byte(`"`+strings.Repeat("x", surfacepage.MaxApproveBodyBytes)+`"`)), nil
+		}},
+		{"R6-grammar", "", "BadRequest", "both", 400, func(a approvalRig, row store.MintRequestRow) ([]byte, []string) {
+			return []byte(`{"v":1,"digest":"` + row.Digest + `","webauthn":{},"sshsig":"x"}`), nil
+		}},
+		{"R7-unknown", "mr-00000000000000f1", "NotFound", "no mint request mr-00000000000000f1", 404, func(a approvalRig, row store.MintRequestRow) ([]byte, []string) {
+			return approveBody(row.Digest, assertionFor(t, apPasskey, row, 2)), nil
+		}},
+		{"R8-expired", "", "Conflict", "expired at", 409, func(a approvalRig, row store.MintRequestRow) ([]byte, []string) {
+			a.d.approveNow = func() int64 { return row.ExpiresAt }
+			return approveBody(row.Digest, assertionFor(t, apPasskey, row, 2)), nil
+		}},
+		{"R8-bytes", "", "Conflict", "do not match its digest", 409, func(a approvalRig, row store.MintRequestRow) ([]byte, []string) {
+			setRequestColumn(t, a.d.credentialsPath, row.RequestID, "digest", strings.Repeat("0", 64))
+			return approveBody(strings.Repeat("0", 64), assertionFor(t, apPasskey, row, 2)), nil
+		}},
+		{"R9-digest", "", "Conflict", "changed since the card was shown", 409, func(a approvalRig, row store.MintRequestRow) ([]byte, []string) {
+			return approveBody(strings.Repeat("ab", 32), assertionFor(t, apPasskey, row, 2)), nil
+		}},
+		{"R10-challenge", "", "Conflict", "challenge is not 64 hex", 409, func(a approvalRig, row store.MintRequestRow) ([]byte, []string) {
+			setRequestColumn(t, a.d.credentialsPath, row.RequestID, "challenge", strings.Repeat("zz", 32))
+			return approveBody(row.Digest, assertionFor(t, apPasskey, row, 2)), nil
+		}},
+		{"R11-other-request", "", "NotApproved", "does not approve request", 422, func(a approvalRig, row store.MintRequestRow) ([]byte, []string) {
+			return approveBody(row.Digest, assertionFor(t, apPasskey, a.putRequest("mr-00000000000000c1", "w201c", apT0-20), 2)), nil
+		}},
+		{"R11-unenrolled", "", "NotApproved", "no enrolled approver", 422, func(a approvalRig, row store.MintRequestRow) ([]byte, []string) {
+			return approveBody(row.Digest, assertionFor(t, other, row, 2)), nil
+		}},
+		{"R11-counter", "", "NotApproved", "does not approve request", 422, func(a approvalRig, row store.MintRequestRow) ([]byte, []string) {
+			return approveBody(row.Digest, assertionFor(t, apPasskey, row, 1)), nil
+		}},
+		{"R11-garbage", "", "NotApproved", "does not approve request", 422, func(a approvalRig, row store.MintRequestRow) ([]byte, []string) {
+			return approveBody(row.Digest, []byte(`{"id":"x"}`)), nil
+		}},
+		{"R11-sshsig-cut", "", "NotApproved", "--signature", 422, func(a approvalRig, row store.MintRequestRow) ([]byte, []string) {
+			return []byte(`{"v":1,"digest":"` + row.Digest + `","sshsig":"-----BEGIN SSH SIGNATURE-----\n"}`), nil
+		}},
+		{"R11-store-v5", "", "Restart", "restart `ailang-worldd serve`", 503, func(a approvalRig, row store.MintRequestRow) ([]byte, []string) {
+			body := approveBody(row.Digest, assertionFor(t, apPasskey, row, 2))
+			setVersion(t, a.d.credentialsPath, `DROP TABLE surface_approvals`, `DROP TABLE webauthn_credentials`, `PRAGMA user_version = 5`)
+			return body, nil
+		}},
+		{"R13-busy", "", "Busy", "busy", 503, func(a approvalRig, row store.MintRequestRow) ([]byte, []string) {
+			held, err := store.OpenCredentials(a.d.credentialsPath)
+			if err != nil {
+				t.Fatal(err)
+			}
+			t.Cleanup(func() { _ = held.Close() })
+			return approveBody(row.Digest, assertionFor(t, apPasskey, row, 2)), nil
+		}},
+		{"R14-v6", "", "Restart", "restart `ailang-worldd serve`", 503, func(a approvalRig, row store.MintRequestRow) ([]byte, []string) {
+			setVersion(t, a.d.credentialsPath, `DROP TABLE surface_approvals`, `PRAGMA user_version = 6`)
+			return approveBody(row.Digest, assertionFor(t, apPasskey, row, 2)), nil
+		}},
+		{"R14-race", "", "Conflict", "changed while it was approved", 409, func(a approvalRig, row store.MintRequestRow) ([]byte, []string) {
+			reads := 0 // the record reads the clock again: the request expired meanwhile
+			a.d.approveNow = func() int64 { reads++; return map[bool]int64{true: apT0, false: row.ExpiresAt}[reads == 1] }
+			return approveBody(row.Digest, assertionFor(t, apPasskey, row, 2)), nil
+		}},
+	} {
+		t.Run(arm.name, func(t *testing.T) {
+			a, row := newPostRig(t)
+			var body []byte
+			var hdr []string
+			if arm.setup != nil {
+				body, hdr = arm.setup(a, row)
+			}
+			id := map[bool]string{true: row.RequestID, false: arm.id}[arm.id == ""]
+			rec := a.post(id, body, hdr...)
+			if r := requireRefusal(t, arm.name, rec, arm.status, arm.class); !strings.Contains(r.Reason, arm.reason) {
+				t.Fatalf("reason %q lacks %q", r.Reason, arm.reason)
+			}
+			if arm.name == "R13-busy" && rec.Header().Get("Retry-After") != "1" {
+				t.Fatalf("Retry-After = %q, want 1", rec.Header().Get("Retry-After"))
+			}
+			if n, s := a.recordedAndSessions(); n != 0 || s != 0 {
+				t.Fatalf("%d recorded, %d sessions after a refusal; want 0 and 0", n, s)
+			}
+		})
+	}
+}
+
+// RT5: a passkey assertion the law verifies is recorded (one row: webauthn, webauthn:mark, the
+// posted envelope's bytes, the request's digest, received at the daemon's now) and nothing mints;
+// requests.json then shows the approver. A refused one before it recorded nothing.
+func TestApprovePOSTRecordsAVerifiedPasskey(t *testing.T) {
+	a, row := newPostRig(t)
+	if rec := a.post(row.RequestID, approveBody(row.Digest, assertionFor(t, apPasskey, row, 1))); rec.Code != 422 {
+		t.Fatalf("a stale counter = %d, want 422", rec.Code)
+	}
+	if n, _ := a.recordedAndSessions(); n != 0 {
+		t.Fatalf("a refused assertion was recorded (%d rows)", n)
+	}
+	assertion := assertionFor(t, apPasskey, row, 2)
+	a.d.approveNow = func() int64 { return apT0 + 3 }
+	rec := a.post(row.RequestID, approveBody(row.Digest, assertion))
+	var got surfacepage.ApproveRecorded
+	if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil || rec.Code != 200 || got != surfacepage.Recorded(row.RequestID, "webauthn:mark") {
+		t.Fatalf("POST = %d %s", rec.Code, rec.Body.String())
+	}
+	if err := approvalHeaders(rec); err != nil {
+		t.Fatal(err)
+	}
+	ro, err := store.OpenCredentialsReadOnly(a.d.credentialsPath)
+	if err != nil {
+		t.Fatal(err)
+	}
+	defer func() { _ = ro.Close() }()
+	ctx, cancel := context.WithTimeout(context.Background(), 5e9)
+	defer cancel()
+	sa, ok, err := ro.SurfaceApproval(ctx, row.RequestID)
+	if err != nil || !ok || sa.Format != "webauthn" || sa.Approver != "webauthn:mark" || !bytes.Equal(sa.Artifact, assertion) || sa.Digest != row.Digest || sa.ReceivedAt != apT0+3 {
+		t.Fatalf("recorded = %+v %v %v", sa, ok, err)
+	}
+	if again, _, _ := ro.MintRequest(ctx, row.RequestID); again.Consumed {
+		t.Fatal("the POST consumed the request")
+	}
+	if n, s := a.recordedAndSessions(); n != 1 || s != 0 {
+		t.Fatalf("%d recorded, %d sessions; want 1 and 0", n, s)
+	}
+	if body := a.serve("GET", "/surface/requests.json", apLoopback, nil).Body.String(); !strings.Contains(body, `"recorded":"webauthn:mark"`) {
+		t.Fatalf("requests.json does not show the approver: %s", body)
+	}
+}
--- a/host/daemon/surface_approve_test.go
+++ b/host/daemon/surface_approve_test.go
@@ -17,7 +17,7 @@ import (
 
 // approvalRoutes is every approval route Handler registers this sprint (the card page and its
 // script are sprint 3's).
-var approvalRoutes = []string{"GET /surface/requests.json", "GET /surface/requests/{id}"}
+var approvalRoutes = []string{"GET /surface/requests.json", "GET /surface/requests/{id}", "POST /surface/approve/{id}"}
 
 // RT1: exactly the approval routes are registered, under /surface, each behind approvalRoute; no
 // other method is answered on them; the card page and its script are not served yet.
````

### Patch M5b — the POST's guards, replay, writer and inventory

````diff M5b
--- a/host/daemon/surface_approve_post_guard_test.go
+++ b/host/daemon/surface_approve_post_guard_test.go
@@ -0,0 +1,171 @@
+package daemon
+
+// Row 201 A3 sprint 2 M5b (design §2.4, §4 T2, T3, T10, T12, §6 RT8-RT12): the POST's cross-origin
+// guards, replay, the writer it never keeps, no session or cookie, and the writer-site inventory.
+
+import (
+	"context"
+	"go/ast"
+	"go/parser"
+	"go/token"
+	"net/http/httptest"
+	"os"
+	"path/filepath"
+	"reflect"
+	"strings"
+	"testing"
+	"time"
+
+	"github.com/sunholo-data/ailang-world/host/authority"
+	"github.com/sunholo-data/ailang-world/host/store"
+)
+
+// RT8: a POST from outside O, cross-site or not JSON is refused and records nothing; a preflight
+// is 405 with no CORS header. A configured origin's page may POST through the loopback Host
+// (behind a proxy); unconfigured, the same origin is refused naming --surface-origin.
+func TestApprovePOSTCrossOrigin(t *testing.T) {
+	a, row := newPostRig(t)
+	body := approveBody(row.Digest, assertionFor(t, apPasskey, row, 2))
+	for _, tc := range []struct {
+		hdr    []string
+		status int
+	}{
+		{[]string{"Origin", ""}, 403}, {[]string{"Origin", "null"}, 403}, {[]string{"Origin", "http://localhost:7645"}, 403},
+		{[]string{"Sec-Fetch-Site", "same-site"}, 403}, {[]string{"Sec-Fetch-Site", "cross-site"}, 403},
+		{[]string{"Content-Type", "text/plain"}, 415}, {[]string{"Content-Type", "application/x-www-form-urlencoded"}, 415},
+	} {
+		if rec := a.post(row.RequestID, body, tc.hdr...); rec.Code != tc.status {
+			t.Errorf("%v = %d, want %d", tc.hdr, rec.Code, tc.status)
+		}
+	}
+	if r := requireRefusal(t, "unconfigured ts.net", a.post(row.RequestID, body, "Origin", "https://cfg.tail0.ts.net"), 403, "CrossOrigin"); !strings.Contains(r.Reason, "--surface-origin https://cfg.tail0.ts.net") {
+		t.Errorf("reason %q does not name the fix", r.Reason)
+	}
+	if rec := a.serve("OPTIONS", "/surface/approve/"+row.RequestID, apLoopback, nil, "Origin", "http://localhost:7645", "Access-Control-Request-Method", "POST"); rec.Code != 405 ||
+		rec.Header().Get("Access-Control-Allow-Origin")+rec.Header().Get("Access-Control-Allow-Methods") != "" {
+		t.Fatalf("preflight = %d %v, want 405 and no CORS header", rec.Code, rec.Header())
+	}
+	if n, _ := a.recordedAndSessions(); n != 0 {
+		t.Fatalf("%d recorded after refusals", n)
+	}
+	a.d.cfg.SurfaceOrigins = []string{"https://cfg.tail0.ts.net"}
+	if rec := a.post(row.RequestID, body, "Origin", "https://cfg.tail0.ts.net"); rec.Code != 200 {
+		t.Fatalf("a configured origin through the loopback Host = %d %s, want 200", rec.Code, rec.Body.String())
+	}
+}
+
+// RT9: the same body twice while pending is recorded once (replaced); after the request is
+// consumed the same body is refused as already used.
+func TestApprovePOSTReplay(t *testing.T) {
+	a, row := newPostRig(t)
+	body := approveBody(row.Digest, assertionFor(t, apPasskey, row, 2))
+	for i := 0; i < 2; i++ {
+		if rec := a.post(row.RequestID, body); rec.Code != 200 {
+			t.Fatalf("POST %d = %d %s", i+1, rec.Code, rec.Body.String())
+		}
+	}
+	if n, _ := a.recordedAndSessions(); n != 1 {
+		t.Fatalf("%d recorded, want 1", n)
+	}
+	a.consume(row, apT0-1, strings.Repeat("cd", 32))
+	if r := requireRefusal(t, "after the consume", a.post(row.RequestID, body), 409, "Conflict"); !strings.Contains(r.Reason, "already used") {
+		t.Fatalf("reason %q", r.Reason)
+	}
+}
+
+// RT10: a held writer is 503 Busy with Retry-After and nothing recorded, never a wait; released,
+// the POST records; after every response the daemon holds no writer.
+func TestApprovePOSTWriterBusy(t *testing.T) {
+	a, row := newPostRig(t)
+	body := approveBody(row.Digest, assertionFor(t, apPasskey, row, 2))
+	held, err := store.OpenCredentials(a.d.credentialsPath)
+	if err != nil {
+		t.Fatal(err)
+	}
+	requireRefusal(t, "busy", a.post(row.RequestID, body), 503, "Busy") // RT4 R13-busy pins Retry-After
+	_ = held.Close()
+	if n, _ := a.recordedAndSessions(); n != 0 {
+		t.Fatalf("%d recorded while busy", n)
+	}
+	if rec := a.post(row.RequestID, body); rec.Code != 200 {
+		t.Fatalf("released: %d %s", rec.Code, rec.Body.String())
+	}
+	w, err := store.OpenCredentials(a.d.credentialsPath)
+	if err != nil {
+		t.Fatalf("the daemon kept the writer after the POST: %v", err)
+	}
+	_ = w.Close()
+}
+
+// RT11: neither a valid session Bearer nor a cookie approves: the handler answers its own
+// refusals; isProtected is unchanged (POST /v1/commit only); the handler reads neither header.
+func TestApprovePOSTNoSessionNoCookie(t *testing.T) {
+	a, row := newPostRig(t)
+	var tok string
+	a.write(func(ctx context.Context, c *store.Credentials) (err error) {
+		tok, _, _, err = authority.Mint(ctx, c, "w201b", grantForAuth(), 3600, time.Now().Unix(), nil)
+		return err
+	})
+	auth := []string{"Authorization", "Bearer " + tok, "Cookie", "session=" + tok}
+	requireRefusal(t, "session, garbage", a.post(row.RequestID, approveBody(row.Digest, []byte(`{"id":"x"}`)), auth...), 422, "NotApproved")
+	requireRefusal(t, "session, no body", a.post(row.RequestID, nil, auth...), 400, "BadRequest")
+	if n, s := a.recordedAndSessions(); n != 0 || s != 1 {
+		t.Fatalf("%d recorded, %d sessions; want 0 and the one minted here", n, s)
+	}
+	if a.d.isProtected(httptest.NewRequest("POST", "/surface/approve/"+row.RequestID, nil)) || !a.d.isProtected(httptest.NewRequest("POST", "/v1/commit", nil)) {
+		t.Fatal("isProtected changed")
+	}
+	src, err := os.ReadFile("surface_approve.go")
+	if err != nil {
+		t.Fatal(err)
+	}
+	for _, read := range []string{`"Cookie"`, `"Authorization"`, `.Cookie(`, `.Cookies()`, "authority.FromContext"} {
+		if strings.Contains(string(src), read) {
+			t.Errorf("surface_approve.go reads %s", read)
+		}
+	}
+}
+
+// RT12: the daemon's production code opens the credentials writer at exactly two sites:
+// ensureCredentialsStore's upgrade at serve start and the POST's record. Control: the read-only
+// opens are found by the same walk.
+func TestDaemonCredentialWriterSites(t *testing.T) {
+	files, err := filepath.Glob("*.go")
+	if err != nil {
+		t.Fatal(err)
+	}
+	sites, readOnly := map[string]string{}, 0
+	for _, f := range files {
+		if strings.HasSuffix(f, "_test.go") {
+			continue
+		}
+		file, err := parser.ParseFile(token.NewFileSet(), f, nil, 0)
+		if err != nil {
+			t.Fatal(err)
+		}
+		for _, decl := range file.Decls {
+			fn, ok := decl.(*ast.FuncDecl)
+			if !ok {
+				continue
+			}
+			ast.Inspect(fn, func(n ast.Node) bool {
+				sel, ok := n.(*ast.SelectorExpr)
+				if !ok {
+					return true
+				}
+				if pkg, isIdent := sel.X.(*ast.Ident); isIdent && pkg.Name == "store" {
+					switch sel.Sel.Name {
+					case "OpenCredentials", "UpgradeCredentials":
+						sites[fn.Name.Name] += sel.Sel.Name
+					case "OpenCredentialsReadOnly":
+						readOnly++
+					}
+				}
+				return true
+			})
+		}
+	}
+	if want := map[string]string{"ensureCredentialsStore": "UpgradeCredentials", "handleSurfaceApprove": "OpenCredentials"}; !reflect.DeepEqual(sites, want) || readOnly < 5 {
+		t.Fatalf("writer sites %v (read-only opens %d), want exactly %v and at least 5 read-only opens", sites, readOnly, want)
+	}
+}
````

### Patch M5c — the end-to-end seam

````diff M5c
--- a/cmd/ailang-worldd/approve_surface_e2e_test.go
+++ b/cmd/ailang-worldd/approve_surface_e2e_test.go
@@ -0,0 +1,83 @@
+package main
+
+// Row 201 A3 sprint 2 M5c (design §2.1, §6 E2E1): the surface approves, the minter mints, in one
+// process: a real daemon Handler() on the request's temp credentials store, no listener.
+
+import (
+	"bytes"
+	"context"
+	"encoding/hex"
+	"encoding/json"
+	"io"
+	"net"
+	"net/http"
+	"net/http/httptest"
+	"path/filepath"
+	"testing"
+	"time"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/daemon"
+)
+
+// E2E1: episode start --request; the card GET; the page's own message hash from the card's
+// digest, nonce and expiry; a test passkey's assertion POSTed (200, recorded); approve --await
+// mints it with no terminal; requests.json then lists it under decided, by webauthn:mark-mac.
+func TestSurfaceApproveThenAwaitMints(t *testing.T) {
+	now := time.Now().Unix() // the daemon's routes read the wall clock
+	a := newApproveRig(t, "ep1", now)
+	seedPasskey(t, a.credPath, "mark-mac", passkeySigner, now-100)
+	d, err := daemon.New(boundedTestContext(t), daemon.Config{DBPath: filepath.Join(t.TempDir(), "world.db"), CredentialsPath: a.credPath,
+		BindHost: daemon.DefaultBindHost, ErrorLog: io.Discard})
+	if err != nil {
+		t.Fatal(err)
+	}
+	defer func() { _ = d.Close() }()
+	serve := func(method, target string, body []byte, out any) int {
+		req := httptest.NewRequest(method, target, bytes.NewReader(body))
+		req = req.WithContext(context.WithValue(req.Context(), http.LocalAddrContextKey, net.Addr(&net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 7644})))
+		req.Host = "localhost:7644"
+		for k, v := range map[string]string{"Origin": pkOrigin, "Sec-Fetch-Site": "same-origin", "Content-Type": "application/json"} {
+			req.Header.Set(k, v)
+		}
+		rec := httptest.NewRecorder()
+		d.Handler().ServeHTTP(rec, req)
+		if err := json.Unmarshal(rec.Body.Bytes(), out); err != nil {
+			t.Fatalf("%s %s: %d %s", method, target, rec.Code, rec.Body.String())
+		}
+		return rec.Code
+	}
+	var card struct {
+		Digest, Challenge string
+		Expires           int64
+	}
+	if code := serve("GET", "/surface/requests/"+a.id, nil, &card); code != 200 || card.Digest != a.row(t).Digest {
+		t.Fatalf("card = %d %+v", code, card)
+	}
+	dg, _ := hex.DecodeString(card.Digest)
+	nonce, _ := hex.DecodeString(card.Challenge)
+	h, err := approval.MessageHash([32]byte(dg), [32]byte(nonce), card.Expires)
+	if err != nil {
+		t.Fatal(err)
+	}
+	body, _ := json.Marshal(map[string]any{"v": 1, "digest": card.Digest, "webauthn": json.RawMessage(passkeySigner.Assert(pkOrigin, pkRPID, h, 2))})
+	var recorded struct{ Recorded bool }
+	if code := serve("POST", "/surface/approve/"+a.id, body, &recorded); code != 200 || !recorded.Recorded {
+		t.Fatalf("POST = %d %+v", code, recorded)
+	}
+	var opened bool
+	env, _ := awaitEnv(time.Now().Unix(), &opened, nil)
+	var stdout, stderr bytes.Buffer
+	if code := runApprove(a.args("--await", "--await-for", "5s"), &stdout, &stderr, env); code != exitOK || opened {
+		t.Fatalf("approve --await: rc %d, terminal opened %v\nstderr %s", code, opened, stderr.String())
+	}
+	requireMintedBy(t, a, "webauthn:mark-mac", stdout.String())
+	var list struct {
+		Pending []json.RawMessage
+		Decided []struct{ ID, ApprovedBy string }
+	}
+	if code := serve("GET", "/surface/requests.json", nil, &list); code != 200 || len(list.Pending) != 0 || len(list.Decided) != 1 ||
+		list.Decided[0].ID != a.id || list.Decided[0].ApprovedBy != "webauthn:mark-mac" {
+		t.Fatalf("requests.json = %d %+v", code, list)
+	}
+}
````
