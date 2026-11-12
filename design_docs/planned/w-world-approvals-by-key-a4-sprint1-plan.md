# Sprint plan — w-world-approvals-by-key-a4, sprint 1 "the law and the typed root" (iteration 389, row 201)

> **THE SPRINT — A4 sprint 1 (design §0 "Sprint 1", §6 "Sprint 1", §9 "Sprint 1 = M1–M3"): the 91-byte decision
> message; the pure decision law D1–D10 (with D4a/D4b) and the reserved-scope predicate L1; the complete,
> request-wide packet discovery; the broker's first-decision-final rule and the publish consumer's key-decision
> refusal L2; the caller census XG1; `GET /surface/decisions/{ref}` and `POST /surface/decide/{ref}` with the
> typed arm only. The typed CLI (CD1–CD4) is moved to sprint 2's head by the design's priced fallback cut,
> forced by the measured size (§1.5, S1-D1).**
> - **Measured: 1,985 changed lines under `cmd/ host/`** — `git diff --numstat <base> HEAD -- cmd host`,
>   insertions + deletions, after a commit of each of the eight scratch stages (**587 production**: `decision.go`
>   69, `decision_law.go` 132, `decisionpacket_discover.go` 63, `approve.go` 36, `surface_decide.go` 190,
>   `decide_model.go` 92, `daemon.go` 5; **1,398 test**). Margin to the 2,000 cap: **15**. The typed CLI
>   (design's ≈ 400 for CD1–CD4) does not fit: **fallback cut fired** (CD1–CD4 → sprint 2's head).
> - **No `/v1` route, no `.ail`, no CSP, no CORS, no new dependency, no credentials writer.** Two `/surface`
>   routes behind A3a's `approvalRoute`; `TestFrozenV1RouteTableMatchesMux` stays 13; the daemon still has one
>   `store.OpenCredentials(` site (`TestDaemonCredentialWriterSites` unedited).
> - **Two premises of the design are false at HEAD and are routed minimally** (§1.4): the first-decision-final
>   rule *does* reach a production decider — `RecordOperatorApproval` re-approves a byte-identical finish request
>   in the same second (an existing, unedited finish test does exactly that), so an identical decision is
>   answered idempotently (P2); and an existing read-census test pins the mint path's reads (P3, one edited line).
> - **Every line of code and test is a unified diff in Appendix P**, generated from compiled, vetted,
>   gofmt-clean, green scratch stages; extracted **from this file** with §2's `awk`, applied to a fresh
>   `git archive 9a6b959d` export and compared byte-for-byte with the mutant tree (§1.2 C30). The executor
>   applies each patch and runs its acceptance; it writes no code by hand.

- **Design**: [w-world-approvals-by-key-a4.md](w-world-approvals-by-key-a4.md), **r2** (quorum r1+r2 BLOCKED on
  narrow defects, controller carve-out, independent sonnet audit VALID 84, iteration 388; **not reopened**).
  **D-WORLD-115** (the decision message with a derived challenge) is OPEN; its default **A = keep it** applies from
  this iteration: built as designed. Sprint 1 does not depend on the answer (§1.3 S1-D2: the message bytes are
  built and tested here, but nothing in sprint 1 verifies a signature over them; under B only the challenge slot
  and sprint 2 change). Where this plan and the design differ, this plan was measured against `9a6b959d` and
  wins; each difference is a decision (§1.3) or a premise found false (§1.4).
- **Starting ground**: A3a code-complete (sprints 1–4 merged; the last `4134ce89`), the A4 design merged
  (`9a6b959d`, PR #562). The design was measured at `b0fbcd3e`; `git diff --stat b0fbcd3e 9a6b959d -- cmd host` is
  **empty** and the 12 Go files of the design's V1 hash list re-measure identical (§1.2 C3), so the design's line numbers hold.
- **Planner base**: `9a6b959deed573d7966863a1f166f3a8c5700f2e` (= origin/dev). Worktree
  `/Users/voightkampff/dev/sunholo-data/.wt-world-iter389-a4-s1`, branch `sprint/world-iter389-row201-a4-s1`.
  Planner scratch (outside the repo, never under `/tmp`): `$HOME/.ailang/state/world-iter389/plan-scratch/`
  (`base/` the pristine export, `work/` the stage repo (tags `M1a` … `M3d`) and the tree the mutants were cut
  from, `patches/`, `tools/` (`gates.sh`, `mutants.py`, `mutate.py`, `roundtrip.sh`), `logs/`, `roundtrip/`,
  `mut/`).
- **Executor**: as routed by the controller. **Two runs**: Run 1 (M1a, M1b, M1c, M2 — the pure core and the
  broker), Run 2 (M3a, M3b, M3c, M3d — the routes and their tests). Each sub-milestone is one `patch -p1` (§2)
  plus its acceptance (§3). **Apply, never retype.**
- **Touched** (§5): new `host/approval/decision.go` (+test), new `host/broker/decision_law.go` (+test),
  `host/broker/approve.go`, new `host/broker/decide_once_test.go`, `host/broker/context_thread_test.go` (one
  census line, §1.4 P3), new `host/coordinator/decisionpacket_discover.go` (+test), `host/daemon/daemon.go` (two
  `HandleFunc` lines + the doc comment), new `host/daemon/surface_decide.go`, new
  `host/surfacepage/decide_model.go`, new `host/daemon/surface_decide_{routes,rig,subject,post,discovery}_test.go`.
  **Outside the design's §8 file list**: `host/approval/decision_test.go`, `host/broker/decision_law_test.go`,
  `host/broker/decide_once_test.go`, `host/broker/context_thread_test.go` (P3), `host/coordinator/
  decisionpacket_discover_test.go` and the five `surface_decide_*_test.go` (the design names tests by id, not
  file; §3 table lists `decisionpacket_discover.go`, `decision_law.go`, `decision.go`, `surface_decide.go`,
  `decide_model.go`, all built).
- **Not touched**: `cmd/**` (the typed CLI moved, S1-D1), `host/store/**`, `host/authority/**`,
  `host/approval/{request,signed,sshsig,passkey_assert}.go`, `host/daemon/{episode_finish,surface_approve,
  surface_approve_guard,surface_page}.go` (reused, unchanged), `host/surfacepage/{approve_model.go,app/**,
  approve/**,passkey/**}`, `docs/QUICKSTART.md` (§18 is sprint 3's), `world/*.ail`, `packages/**`, `scripts/**`,
  `tools/launchd/*`, `cmd/world-publish/*`, any skill, `~/.ssh`, `~/.ailang/world`, `~/.ailang/state/mission-v1*`.

---

## 0. Binding constraints (each is restated where it governs)

1. **Test keys and temp stores only (D-WORLD-109 = A).** Every store is under `t.TempDir()`; every packet is a
   forged, validated commit built with the existing `workbench_packet_test.go` helpers (`approvalRequestInChain`,
   `newForgedChain`, `commitForgedChain`); every operator proof is a file the test writes as
   `writeOperatorProof` does. No key is used at all in sprint 1. Nothing opens `~/.ailang/world`, `~/.ssh`, a real
   store or a real `/dev/tty`; no `ssh-keygen`; no mint into a real store; no publish; no listener beyond
   `httptest` recorders.
2. **Pure core, effects at the boundary (coding standards S1–S3).** The message (`host/approval/decision.go`,
   standard library only, `TestApprovalImportsOnlyTheStandardLibrary`) and the law
   (`host/broker/decision_law.go`: no I/O, no clock; `now` is a fact) are pure. The discovery reads the store
   through a narrow interface (`DiscoveryReader`). The daemon gathers facts (store, clock, body) and writes once
   (`broker.DecideApproval`) under `humanMu`. No `.ail` (design §3 S1: the law mirrors no sketch law; no `/v1`
   route, so `worlddapi.ail` is untouched).
3. **`decidedBy` is derived, never accepted.** The body grammar has no `decidedBy` (strict decode refuses it);
   the typed arm passes the constant `operator:terminal`.
4. **The daemon never keeps a credentials writer and opens none here**; no CORS header anywhere; `surfaceCSP`
   unchanged; `/v1` stays 13.
5. **One mutant per new production branch, anchored to the diff** (§4), each run with its named test and its
   package suite (inventory and static tests excluded); named equivalents with reasons; every run appends a
   ledger line.
6. **No `t.Parallel()`** in any new test (`TestNoParallelWriteForkPackagesOutsideVerifygate`).
7. **Bisectable and safe at every commit.** Each sub-milestone ends gofmt-clean, `go build ./...` and `go vet`
   green, the `_test.go` compile fence (`go test -run '^$' ./...`) green, its own tests and every earlier one
   green.
8. **No git writes by the executor** (§7). The controller commits from snapshots.

---

## 1. Pre-flight

### 1.1 Baseline at `9a6b959d` (measured by the planner, 2026-11-12)

Every command ran with `export AILANG_BIN=$HOME/.pinned-ailang/ailang` (AILANG v0.41.0), `go1.26.6 darwin/arm64`,
in a `git archive 9a6b959d` export at `plan-scratch/base` (a throwaway git repo; `diff -r cmd` and `diff -r host`
against the worktree: identical, `BASE_EQUALS_WORKTREE`; the worktree itself was only read). The pre-commit legs
are `tools/gates.sh <tree> <logs> base` = §6's pre-commit block without the sprint's new tests. **The rig carried a
load average of 90–190 on 16 cores from other sessions throughout** (`uptime`, measured at 17:35, 18:06, 18:33);
every timing below is under that load. Logs: `plan-scratch/logs/baseline/`.

| Command | rc | Time | Result |
|---|---|---|---|
| `git -C <worktree> status --porcelain \| wc -l` | 0 | — | `0` (clean at `9a6b959d`) |
| `./scripts/verify_ail.sh` | **0** | — | `✓ verify gate PASSED: 16 required identities verified, 40 named tests pass, 503 se-tools named tests pass` (`baseline/verify_ail.log`) |
| `gofmt -l cmd host`; `go build ./...`; `go vet ./...`; `go test -run '^$' ./...` | **0** | — | empty / clean / empty / no FAIL line |
| `go test -count=1 -timeout 30m ./host/approval/... ./host/coordinator/ ./host/surfacepage/ ./host/daemon/` | **1** (every red green alone but one environmental) | — | `ok host/approval 21.3s`, `ok …/approvaltest 23.7s`, `ok host/surfacepage 15.2s`; `FAIL host/coordinator 616s` on `TestCapsuleWarmEqualsColdForSeTools` (`ailang-write: CheckSource = (… Passed:false}, <nil>)`: the interpreter under load); `FAIL host/daemon 1040s` on `TestMarketingCaptureInventoryMutation`, `TestServeStartsWhenTheProbeSpendsTheStartupClock`, `TestRegistryPhaseHasItsOwnClock` (`baseline/pkgs.log`). **Re-run alone** (`logs/base_rerun_alone.log`): `--- PASS: TestCapsuleWarmEqualsColdForSeTools (56.14s)`, `--- PASS: TestMarketingCaptureInventoryMutation (4.58s)`, `--- PASS: TestRegistryPhaseHasItsOwnClock (5.66s)`; `--- FAIL: TestServeStartsWhenTheProbeSpendsTheStartupClock (1.82s)` again: `daemon startup failed at config: … --exec-node …/node --version: broker: handler subprocess timed out after 5s` — the rig's `node --version` exceeding 5 s, environmental, untouched by the sprint |
| `go test -count=1 -timeout 30m ./host/broker/` (reported, not gating) | **1** | 673.5 s | only `--- FAIL: TestExecProtectedCheckTimeout` (known list) (`baseline_rest/broker.log`) |
| `-race` legs (§6's three) + every named inventory, `-v` (base mode: no sprint test) | **0** | 265 s | `ok host/approval 3.0s`, `ok host/coordinator 142.5s`, `ok host/broker 7.0s`, `ok host/daemon 27.0s`, no race report; **33** named `--- PASS`, 0 SKIP, 0 FAIL; `PRE_COMMIT_GATES_GREEN 265s` (`baseline_tail/`) |
| `VERIFY_GO_LEGS=plain ./scripts/verify_go.sh` (after commit) | — | — | **not measured by this planner** (the iteration's time budget went to the mutant set under a load average of 90–190); it is the controller's after-commit leg (§6), with the rule below. The A3a sprint-4 plan measured it at its base under load: rc 1 on three known-list reds, each green alone |

**Rule (carried):** a red at the base measures the repo, not the sprint. The controller re-runs a red untouched
package alone once and does not chase it. **Never start a gate under `nohup`.** `VERIFY_GO_LEGS=plain
./scripts/verify_go.sh` reads committed blobs and runs `go test ./...`; it is the after-commit leg (§6).

### 1.2 Codebase claims checked (at `9a6b959d` unless "work" — the planner's final scratch tree)

| # | Claim | Command → observed output |
|---|---|---|
| C1 | The base | `git rev-parse HEAD origin/dev` (worktree) → `9a6b959deed573d7966863a1f166f3a8c5700f2e` twice; `git status --porcelain \| wc -l` → `0` |
| C2 | **Nothing of A4 has landed** | `grep -rln "DecisionMessage\|DecisionRefusal\|surface/decide\|DiscoverRequestPackets\|ErrApprovalAlreadyDecided" host cmd` → only `host/surface/validate.go`, `host/surface/episodes.go`, `host/surface/episodes_test.go`, all on the identifier `DecisionRefusal` (`validate.go:183`, `episodes.go:50`, `:604`) of the **unrelated** package `host/surface` (no collision: the new one is `broker.DecisionRefusal`, §1.4 P9); `grep -rn "/surface/decisions" host cmd` → no match; control in the same call: `grep -c "func SignedMessage" host/approval/request.go` → `1` |
| C3 | **No code drift since the design's base** | `git diff --stat b0fbcd3e 9a6b959d -- cmd host` → empty; `shasum -a 256` of the 12 Go files of the design's V1 list (`app.js` not re-hashed: sprint 1 does not touch it) → `aeef311293e1c42d approve.go`, `198bb050061c92b9 approvals_view.go`, `ff50de9b6af8bf92 decisionpacket.go`, `ac56423782a187aa decisionpacket_validate.go`, `0495878c6aedce4b daemon.go`, `4e320c537e689620 episode_finish.go`, `6f3cafda604f4c23 surface_approve.go`, `d3a6620f0aa00022 judge.go`, `cf9897ec3f10881e request.go`, `a40850200782b3e4 sshsig.go`, `6a38843161c7d75f cmd/…/approve.go`, `0adb363562679cd5 episode_finish_operator.go` — **every prefix equals the design's header**; so every design line number below holds |
| C4 | **The decide entry point's callers** (controller-verified, re-measured) | `grep -rn "DecideApproval(\|decideApproval(" --include="*.go" host cmd \| grep -v _test.go` → `approve.go:181` (def), `:188` (the wrapper's call), `:191` (def), `:696` (`RecordOperatorApproval`), `publish_op.go:153` (`mintAttendedApproval`); the controller's `grep -rl "DecideApproval(" cmd host --include="*.go"` → `approve.go` + 5 `_test.go` files |
| C5 | **`decideApproval` checks existence and the chain only** (design §1.6) | `approve.go:191 func decideApproval`, `:208 if _, found, err := findApprovalRequest(ctx, s, requestRef)`, `:213 payload := mustApprovalJSON(approvalDecisionWire{` — no decision read between them; `findApprovalDecision` (`:276`) walks the whole chain newest first |
| C6 | **ValidatePacket's duplicate scan is a 100-packet prefix** (design V33) | `decisionpacket_validate.go:34 const packetScanLimit = 100`, `:190 ObjectsBySemanticID(ctx, DecisionPacketV1, "", packetScanLimit+1)`, `:194 res.UniqueTruncated = …`, `:195 for _, o := range others[:min(…)]`; `samePair` (`decisionpacket.go:76`) requires the same request **and** proposal, so two packets with one request and two proposals both validate with code 0 (work: RD9b asserts `ValidatePacket` of the first is `0` while the second exists) |
| C7 | **The keyset page contract** (design V34) | `objects_by_semantic_id.go:13 const MaxSemanticIDPage = 500`, `:20 … WHERE semantic_id = ? AND hash_ref > ? ORDER BY hash_ref LIMIT ?`; payload nil (callers `GetObject`) |
| C8 | **The operator-proof fence reused verbatim** (design §1.11) | `episode_finish.go:92 func (d *Daemon) consumeOperatorProof(episode, phase, proof string) error`, `:102 path := filepath.Join(dir, episode+".json")`, `:117 … p.Phase != phase \|\| p.ExpiresAt <= now \|\| p.ExpiresAt > now+300` with `now := d.finishNow()`; `:158 d.humanMu.Lock()`; `daemon.go:440 humanMu sync.Mutex`, `:460 finishNow`, `:464 approveNow`, `:150 commitBudget = 3 * time.Second` |
| C9 | **A3a's guards are reusable predicates** (design V27) | `surface_approve_guard.go:41 allowedApprovalHost`, `:53 allowedApprovalOrigin`, `:57 sameOriginFetch`, `:61 jsonMediaType`, `:86 crossOriginReason`; `surface_approve.go:76 func (d *Daemon) approvalRoute`, `:215 io.ReadAll(http.MaxBytesReader(w, r.Body, surfacepage.MaxApproveBodyBytes)) // R6` (20 KiB, `approve_model.go:27`) |
| C10 | **Route registration; `/v1` is 13** | `daemon.go:1137-1141` the five A3a routes, each `d.approvalRoute(`; `grep -c 'mux.HandleFunc("[A-Z]* /v1/' daemon.go` → `13`; work: the two new lines follow `:1141`, `/v1` still `13` (`TestFrozenV1RouteTableMatchesMux` PASS) |
| C11 | **The publish consumer never reads `decidedBy`** (design V11) | `approve.go:517 func validatePublishApproval`, `:544 if decision.Decision != "approve" {`; `DecidedBy` occurs only at `:77` (the wire) and `:214` (the decide write) |
| C12 | **Route inventories do not list the decision names** (so no inventory edit) | `surface_approve_test.go:20 var approvalRoutes = […]` matched by ``mux\.HandleFunc\("([A-Z]+ /surface/(?:requests\|approve)[^"]*)"`` — `decide`/`decisions` are outside it; `surface_app_test.go:154-161` the 404/421 lists name no `/surface/deci…` path; work: `TestApprovalRoutesRegistered`, `TestSurfaceAppRoutes` PASS unedited |
| C13 | **No other test enumerates routes, exported symbols, files or verbs this sprint touches** | `grep -rn 'ReadDir("\.")' host/coordinator host/broker host/daemon host/surfacepage --include='*_test.go'` → only `surface_approve_post_guard_test.go` (`TestDaemonCredentialWriterSites`, which counts `OpenCredentials(`: none added); `TestProductionGoSurface` and `TestProductionContextRoots` scan files outside `host/`/`cmd/` and context roots (none added); `TestEveryVerbHasHelp` (no CLI change, S1-D1); `TestReadCtxCancelledAfterHandler` covers its `seedReadRoutes` list (the GET uses `readCtx` + `defer cancel()`, the POST `context.WithTimeout(r.Context(), d.commitBudget)` + `defer cancel()`); every inventory by name PASS on the final tree (C26) |
| C14 | **The test rig exists** (no new store helper) | `workbench_packet_test.go:157 forgeObject`, `:177 approvalRequestInChain`, `:207 newForgedChain`, `:247 commitForgedChain` (a forged commit whose packet validates, `wantValidated` `:293`); `surface_approve_rig_test.go:38 newApprovalRig`, `:107 serve`, `:137 requireRefusal`; `episode_finish_test.go:18 finishProofFile` (the proof file shape) |
| C15 | **Identical requests collide, and a production decider re-decides one** (design §1.5; §1.4 P2) | `episode_finish_test.go:194-206 TestEpisodeFinishOperatorProofExpiryBoundary`: one `pinFinishClock(d, finishClockAt)`, one `input` scope, arms `{1, 200}, {300, 200}` → `RecordOperatorApproval` (`approve.go:686-697`) puts the same request object twice and decides it twice at the same `now`; work, with the design's literal rule (refuse any second decision): `episode_finish_test.go:204: proof expiry now+300: 500, want 200: {"error":{"class":"Internal",…}}` (`logs/m2_daemon.log`); with S1-D10: PASS |
| C16 | **A read-census test pins the mint path's reads** (§1.4 P3) | `context_thread_test.go:61 if w.objects != 5 \|\| w.heads != 4 {` ("Exact read counts prevent dead arms"); work, with BR1: `mint read census objects=6 heads=5, want 5/4` (`logs/m2_pkgs.log`); BR1's walk adds exactly one head read and one object read |
| C17 | **`host/approval` stays standard-library only** | `imports_test.go` forbids non-stdlib and `os`/`net`/`crypto/rand`/…; `decision.go` imports `crypto/sha256`, `encoding/binary`, `errors`; work: `TestApprovalImportsOnlyTheStandardLibrary` PASS |
| C18 | **The broker cannot import the coordinator** (S1-D3) | `decisionpacket_validate.go:8 "github.com/sunholo-data/ailang-world/host/broker"` — the coordinator imports the broker, so `Discovery` and `DecisionDiscoveryLimit` live in `host/broker` and the coordinator returns `broker.Discovery` |
| C19 | **The typed CLI does not fit** (S1-D1) | work, `git diff --numstat` per stage: M1–M3 without CD1–CD4 = **1,985**; design §9 prices CD1–CD4 at ≈ 400 (1,950 − 1,550); 1,985 + 400 > 2,000 |
| C20 | **Every repo-wide inventory and every new test, on the final tree, by name** | `tools/gates_tail.sh <round-trip tree> logs/final_tail final` (= §6's `-race` and inventory legs) → `PRE_COMMIT_GATES_GREEN 312s`; **68 named `--- PASS`**, 0 SKIP, 0 FAIL (33 carried inventories, as at the base, + 35 new): `TestApprovalImportsOnlyTheStandardLibrary`, `TestApprovalTestIsTestOnly`, `TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential`, `TestQuickstartSection10ProfilesLoad`, `TestQuickstartProtectedContract`, `TestApprovalContextIdentity` (6/5, P3), `TestValidatePacketBytesCheckNeedsThePacketLabel`, `TestDaemonDependencyAllowlist`, `TestFrozenV1RouteTableMatchesMux` (13), `TestSurfaceAppRoutes`, `TestApprovalRoutesRegistered`, `TestReadCtxCancelledAfterHandler`, `TestMCPMountMethodSource`, `TestDaemonCredentialWriterSites` (unchanged: one site), `TestQuickstartDecisionPacketSection`, `TestDaemonConstructorCallerInventory`, `TestEpisodeFinishOperatorProofExpiryBoundary` (P2), `TestWorldBoundaryDependencyAllowlist`, `TestBareNetHTTPExemptionIsPerGroup`, `TestEveryCommandDirectoryIsAProtectedGroup`, `TestProductionContextRoots`, `TestProductionGoSurface`, `TestNoNewDeadlineFreeStoreReads`, `TestReadGettersHonorContext`, `TestCredentialsMethodsRequireADeadline`, `TestNoParallelWriteForkPackagesOutsideVerifygate`, `TestSubprocessSinkGate`, `TestEveryVerbHasHelp`, the five quickstart CLI tests; and every sprint test (DM1, DM2, DL1–DL12, DL4a, DL4b, R4a, DS1, BR1, BR2, R4b, R4c, XG1, RD1–RD9, RD9b–RD9d, R4d) — e.g. `--- PASS: TestDecideDiscoveryRefusesAnOverflow (3.84s)`, `--- PASS: TestDecideApprovalCallerCensus (0.19s)` (`logs/final_tail/inv.log`) |
| C21 | `-race` on the touched packages, final tree | same run: `ok host/approval 2.1s`, `ok host/coordinator 154.3s`, `ok host/broker 9.5s` (the sprint's and the finish/publish approval tests), `ok host/daemon 40.3s` (`TestDecision*`, `TestDecide*`, `TestEpisodeFinish*`, `TestApprov*`, `TestRequest*`, …); no race report (`logs/final_tail/race.log`) |
| C22 | The packages, final tree | `tools/gates.sh <round-trip tree> logs/final_gates2 final`: `verify_ail.sh` → `✓ verify gate PASSED: 16 required identities verified, 40 named tests pass, 503 se-tools named tests pass` (no `.ail` changed); gofmt empty, build, `go vet ./...` clean, the compile fence clean; `ok host/approval 0.9s`, `ok …/approvaltest 0.4s`, `ok host/coordinator 73.4s`, `ok host/surfacepage 2.7s`; `host/daemon` **one** red, on the known list: `--- FAIL: TestMCPPostBudgetProductionConstants/batch (15.83s)` (405.9 s run); every sprint test and every other daemon test PASS in that run (`logs/final_gates2/pkgs.log`). Re-run alone: `--- PASS: TestMCPPostBudgetProductionConstants (42.45s)` with `/single` and `/batch` PASS, `ok host/daemon 42.7s` (`logs/final_rerun_alone.log`) |
| C23 | `host/broker` whole, final tree (reported, not gating) | on a clone of the mutant tree (`diff -r` of `host` against the round-trip tree: identical; the round-trip tree had been re-created under the first attempt): `ok host/broker 356.8s`, rc 0 — every `registry_publish_test.go` arm, `TestApprovalContextIdentity` (6/5) and every sprint broker test PASS (`logs/final_broker.log`). At the base the same leg was rc 1 on `TestExecProtectedCheckTimeout` alone (known list, §1.1) |
| C24 | **Every mutant, on the final tree** | 96 rows, the last ledger line of each: **88 KILLED (named and suite), 8 `E` rows** (3 named equivalents, 5 named survivors, each with its reason), 0 not as expected, 0 `NOCOMPILE`; 104 ledger lines (runs `run1` 4 stopped, `run2` 96, `run3` 3, `run4` 1) (§4, `logs/mutant_ledger.tsv`) |
| C25 | **Nothing outside §5 moved** | round-trip tree: `git diff --stat HEAD~8 HEAD -- world packages scripts tools docs design_docs cmd \| wc -l` → `0`; production files: `git diff --stat HEAD~8 HEAD -- '*.go' ':!*_test.go'` → `7 files changed, 586 insertions(+), 1 deletion(-)` (`decision.go`, `decision_law.go`, `decisionpacket_discover.go`, `approve.go`, `daemon.go`, `surface_decide.go`, `decide_model.go`); added lines matching `t.Parallel\|Access-Control\|OpenCredentials(\|context.Background()` → 2, both in tests (`decisionpacket_discover_test.go` `ctx := context.Background()` for the fake reader; `surface_decide_routes_test.go` sends `Access-Control-Request-Method` to prove no CORS answer) |
| C26 | **Every sub-milestone is committed and measured** | the round-trip tree commits each patch in turn: `NUMSTAT M1a 162`, `M1b 508`, `M1c 733`, `M2 1004`, `M3a 1356`, `M3b 1628`, `M3c 1901`, `M3d 1985` (`logs/roundtrip.log`); each stage's own tests were green when the stage was built (`work`), and the final tree is green as above. The per-stage acceptance commands of §3 were **not** re-run one by one on the round-trip tree (time; the stages are cumulative and their tests all PASS on the final tree, C20) |
| C27 | **The patches reproduce the mutant tree** | `tools/roundtrip.sh`: a fresh `git archive 9a6b959d` export + `patch -p1` of Appendix P's eight patches **extracted from this file with §2's `awk`** (each `cmp`-equal to the generated patch; `--dry-run` first; no `.orig`/`.rej`) → `diff -r` of `cmd`, `host`, `docs` against `plan-scratch/work` (the tree §4's mutants were cut from): **identical** (`ROUNDTRIP_IDENTICAL_TO_MUTANT_TREE`, `logs/roundtrip.log`); re-run after this file's final assembly: `ROUNDTRIP_IDENTICAL_TO_MUTANT_TREE` with the same eight `NUMSTAT` values (`logs/roundtrip_final.log`) |

### 1.3 Planner decisions where the design leaves a choice

- **S1-D1 — The fallback cut fires: CD1–CD4 (the typed CLI) move to sprint 2's head** (design §9's priced cut).
  Measured M1–M3 without the CLI: **1,985** (§1.5). The CLI is the `approve <sha256:…> [--deny]` dispatch, GET +
  hash checks, every field printed, `openTerminal`, `writeOperatorProof` with phase `decide-<verb>` and the name
  `decision-<hex Q>`, the POST, F1 rules, help lines and `TestEveryVerbHasHelp` (design ≈ 400). Never cut and
  built: M1, M2, the discovery, D4a, RD9b–RD9d. **Consequence:** after this sprint the decide path exists and is
  proven over HTTP, but no shipped command writes the `decision-<hex>` proof, so an operator cannot yet decide at
  a terminal without hand-writing the proof (§9 AC12, Appendix C). Design Q2 (`approve --list` decisions) moves
  with the CLI.
- **S1-D2 — D-WORLD-115 independence, checked.** Sprint 1 builds `DecisionMessage`/`DecisionMessageHash`/
  `DecisionNamespace` and their layout tests, and nothing verifies a signature over them; the typed arm signs
  nothing. Under B (a stored daemon challenge) the 32-byte P slot would carry the challenge; DM1's goldens and
  `DecisionMessage`'s signature change in sprint 2, nothing else of sprint 1 (no route, law step or test reads
  the message). Recorded, not a blocker.
- **S1-D3 — `broker.Discovery` and `broker.DecisionDiscoveryLimit = 2000`** live with the law (C18);
  `coordinator.DiscoverRequestPackets(ctx, DiscoveryReader, request) (broker.Discovery, carrying store.Object,
  error)` pages `ObjectsBySemanticID(DecisionPacketV1, after, store.MaxSemanticIDPage)`; a packet past the
  2,000th is the sentinel (Complete stays false); a short page is the end (Complete true); a payload that is
  missing or does not decode counts in `Undecodable`.
- **S1-D4 — The law's form.** `LawRefusal{Step, Status, Reason}` (Step `""` accepts) instead of the design's
  `(code int, reason string)`, so every test asserts the step id; `SubjectRefusal(f)` is D1–D7 (what GET shows)
  and `DecisionRefusal(f)` = `SubjectRefusal` then D8–D10. The daemon **gathers every fact it can independently**
  (request, chain state, discovery, and the packet verdict only for exactly one carrying packet) and calls the
  law once; the order lives only in the law (DL11). D4b's three sentences: zero carrying, N carrying, the
  validator's sentence (`packetRefusalSentence`, unchanged). D7 prints the deadline as RFC 3339 UTC (pure
  formatting, no clock read).
- **S1-D5 — GET `/surface/decisions/{ref}`** answers a D1–D4 refusal with its status and the inert refusal JSON
  (there is no subject without exactly one validated packet); otherwise 200 with `{v, request, requestBytes,
  packet, packetBytes, entry, proposal, now, window{min,max}, refusal, reserved, decidedBy, decision,
  ssh{namespace}}` (bytes base64url without padding; no message, no hash). `refusal` is D5–D7 at `now`
  (display). **`passkeys` arrives with sprint 2's key arms** (it needs the read-only credentials handle; §1.4 P7).
  GET takes no lock (display; the POST re-reads under `humanMu`).
- **S1-D6 — Sprint 1's body grammar is typed only:** `{"v":1,"verb","packet":<64 hex>,"request":<64 hex>,
  "expiry","typed":true}`, strict (`DisallowUnknownFields`, one value), so `sshsig`, `webauthn` and `decidedBy`
  are unknown fields (400) until sprint 2 widens the grammar to "exactly one of three". The verb and the expiry
  are **not** grammar: D9 and D10 are reachable end to end (RD4).
- **S1-D7 — R0** = any `Authorization` or `Cookie` header, even empty (`len(r.Header.Values(…)) != 0`), checked
  after R2 (approvalRoute wraps the handler) and before R3.
- **S1-D8 — The POST reads and writes on `d.store` under `humanMu`** with `context.WithTimeout(r.Context(),
  d.commitBudget)` (the finish route's posture, `episode_finish.go:158-163`); the law's `now` and the decision's
  `now` are `d.approveNow()` (the approval routes' clock); the proof's expiry is checked by the reused
  `consumeOperatorProof` against `d.finishNow()` (C8).
- **S1-D9 — N4, the handler order sprint 2's K0 slots into (fixed here so sprint 2 reorders nothing).** In
  `handleSurfaceDecide`: **R2** (approvalRoute) → **R0** → **R3** → **R4** → **R5** → **R6** (read ≤ 20 KiB +
  grammar) → `humanMu.Lock()` → `readDecision` (facts) → **law D1–D10** → **[the arm slot, D11]** →
  **`broker.DecideApproval`** → 200. Sprint 1's arm: `consumeOperatorProof("decision-"+hex(Q),
  "decide-"+verb, header)` (unlinks before the write). **Sprint 2 fills the same slot for a key format, in
  design §2.7's order, and moves no sprint-1 line:** judge on a read-only credentials handle (W1–W14, P2–P5,
  `Binds`; **W15 deferred**) → non-waiting writer open (busy → 503 `Busy`, `Retry-After: 1`) →
  `RecordKeyDecision`'s one transaction {K1 deadline/schema 8, K2 grammar, **K0**: reload the approver
  (`useApprover`, both paths) and the passkey row; an unlanded row whose artifact and signed subject fields
  are byte-equal → skip W15/`usePasskey`; otherwise W15 + `usePasskey`; K5 insert or replace the unlanded row}
  → the shared `broker.DecideApproval` → `LandKeyDecision` → writer close → 200. So K0 sits **after** the law
  (D6 still refuses a landed retry first, as design §2.7 says) and **after** the read-only judge, and **before**
  `DecideApproval`. Sprint 2's only edits to sprint-1 code: the grammar (S1-D6), a `switch` at the arm slot, and
  a post-`DecideApproval` land call for the key arms. The law's inputs need nothing new for K0.
- **S1-D10 — An identical decision is answered idempotently (routes §1.4 P2).** `decideApproval` builds the
  decision payload first; if the chain already holds a decision of the request: byte-identical payload (same
  verb, decider and `now`) → return its existing ref, write nothing; otherwise `ErrApprovalAlreadyDecided`,
  write nothing. This keeps `RecordOperatorApproval`'s same-second re-approval (C15) answering the same refs it
  answered before, writes no second decision, and refuses every different second decision (BR1). A4's own POST
  never reaches it with a decided request (D6 refuses first, §1.4 P8).
- **S1-D11 — `ApprovalChainState(ctx, *store.Store, ref) (onChain, decided, verb, decidedBy, err)`** is the
  broker's exported read of D3 and D6 (it wraps the unexported walks; no signature of theirs changes).
- **S1-D12 — RD1 is a new test** (`TestDecisionRoutesRegistered`); `TestSurfaceAppRoutes` and
  `TestApprovalRoutesRegistered` stay unedited (C12; design §7 "TestSurfaceAppRoutes +2" is not needed).
- **S1-D13 — Fillers for RD9b/RD9c/RD9d are stored packet objects** (`PutObject`, no commit): the discovery
  counts every stored packet whatever its validity (design §2.4: "whatever its validity"), so only the carrying
  packets need to be validated commits. RD9b/RD9d pick a proposal label whose packet hash sorts after position
  101 and assert the position through the store's own keyset order before posting.
- **S1-D14 — `approval.DecisionVerb(word)`** (the word → byte map) ships with the message for sprint 2's CLI and
  judge; it is pinned by DM1/DM2 and has no sprint-1 production caller.
- **S1-D15 — R4d's must-die half** ("with L1 compiled out, L2 still refuses the resulting key decision") needs a
  key decision, which sprint 1 cannot make: in sprint 1 the L1-deleted mutant (`LAW-D5`) is killed by R4d, and L2
  is proven at the consumer by R4b/R4c; sprint 2's R4e chains the two (design §6 already says "sprint 2 extends
  it with a key").
- **S1-D16 — Sub-milestones** (§1.5): M1a the message (162), M1b the law and L1 (346), M1c the discovery (225),
  M2 the broker rules (271), M3a the routes and handler (352), M3b the rig and RD3 (272), M3c the POST tests
  (273), M3d the discovery end to end (84). M1b and M3a are at the top of the 250–350 band; each is one coherent
  unit (a law and its one-test-per-branch table; a handler and its route inventory) and splitting either would
  ship production code without its tests.

### 1.4 Premises found false or incomplete at HEAD (each routed around minimally)

- **P1 — Design §7 "`TestSurfaceAppRoutes`: +2 (sprint 1)" / §6 RD1 "`TestSurfaceAppRoutes` edited".** Neither
  it nor `TestApprovalRoutesRegistered` lists the decision names (C12); both stay unedited and RD1 is new (S1-D12).
- **P2 — Design §0/§2.4 "the two production deciders decide a request they have just created, so they are
  unaffected".** False when the just-created request is byte-identical to an earlier decided one (identical
  requests collide, design §1.5): `RecordOperatorApproval` re-approves the same finish scope at the same second,
  and the unedited `TestEpisodeFinishOperatorProofExpiryBoundary` does exactly that (C15: 500 with the literal
  rule). Routed by S1-D10 (identical → idempotent, different → refused). The design's "must stay green,
  unedited" list (`episode_finish*_test.go`) holds.
- **P3 — An existing read-census test pins the mint path's reads** (C16). `context_thread_test.go:61` is edited
  from `5/4` to `6/5` with a one-line comment; it is not in the design's "must stay unedited" list. This is the
  only edit to an existing test.
- **P4 — Design §2.5 "`DecisionRefusal(f DecisionFacts) (code int, reason string)`" and `Discovery` in the
  coordinator.** Built as `LawRefusal` (S1-D4) and `broker.Discovery` (S1-D3, C18). Same facts, same order.
- **P5 — Design §0/§9 "after sprint 1 an operator at a terminal can approve or deny a packet".** Not after this
  sprint: the CLI moved (S1-D1). The route and its tests stay; design §10's typed step waits for sprint 2's head.
- **P6 — Design §6 R4d's must-die control needs a key decision** (S1-D15); sprint 2's R4e carries it.
- **P7 — Design §2.8's subject JSON lists `passkeys`**: built in sprint 2 with the key arms (S1-D5).
- **P8 — BR1's refusal is not reachable end to end through the new POST**: under `humanMu`, D6 refuses every
  decided request before `DecideApproval` (the iteration-388 audit's N1 pattern). BR1 is proven at the broker
  (unit) and is the guard for the other two deciders and any future caller (XG1 pins who they are).
- **P9 — An unrelated `DecisionRefusal` identifier exists in `host/surface`** (C2); no collision (package
  `broker`), recorded so a grep reader is not misled.

### 1.5 The size

Lines are insertions + deletions under `cmd/ host/`: **`git diff --numstat <base> HEAD -- cmd host`**
(`docs/` and `design_docs/` are outside the filter). **Every row is measured** from the scratch stages, each
committed in a throwaway repo (`plan-scratch/work`, tags `M1a` … `M3d`), and re-measured on the round-trip tree
(`logs/roundtrip.log` `NUMSTAT` lines).

| Sub-milestone | Production | Test | Own | Cumulative vs `9a6b959d` (after its commit) | Run |
|---|---|---|---|---|---|
| M1a — the decision message (`decision.go`); DM1, DM2 | 69 | 93 | 162 | **162** | 1 |
| M1b — the law D1–D10 and L1 (`decision_law.go`); DL1–DL12, DL4a, DL4b, R4a | 132 | 214 | 346 | **508** | 1 |
| M1c — the discovery (`decisionpacket_discover.go`); DS1 | 63 | 162 | 225 | **733** | 1 |
| M2 — BR1 (+ idempotent identical), `ApprovalChainState`, L2; BR1, BR2, R4b, R4c, XG1; the census line (P3) | 36 | 235 | 271 | **1,004** | 1 |
| M3a — `surface_decide.go`, `decide_model.go`, two routes; RD1; XG1 gains the handler | 287 | 65 | 352 | **1,356** | 2 |
| M3b — the rig; RD3 | 0 | 272 | 272 | **1,628** | 2 |
| M3c — RD2, RD4, RD5, RD6, RD7, RD8, R4d | 0 | 273 | 273 | **1,901** | 2 |
| M3d — RD9, RD9b, RD9c, RD9d | 0 | 84 | 84 | **1,985** | 2 |
| **Sprint 1** (cumulative, `git diff --numstat 9a6b959d M3d -- cmd host`) | **587** | **1,398** | | **1,985 ≤ 2,000** (margin 15) | |

Per file, cumulative: `decision.go` +69, `decision_test.go` +93, `decision_law.go` +132, `decision_law_test.go`
+214, `decisionpacket_discover.go` +63, `decisionpacket_discover_test.go` +162, `approve.go` +36,
`context_thread_test.go` +3 −2, `decide_once_test.go` +231, `daemon.go` +4 −1, `surface_decide.go` +190,
`decide_model.go` +92, `surface_decide_routes_test.go` +64, `surface_decide_rig_test.go` +164,
`surface_decide_subject_test.go` +108, `surface_decide_post_test.go` +273, `surface_decide_discovery_test.go`
+84. Against the design's ≈ 1,950 with the CLI (≈ 1,550 without): the built part came in ≈ 435 over the design's
price, almost all test (the guards and law arms each assert 0 decisions, an unmoved head and an unconsumed proof).

---

## 2. How the executor applies a sub-milestone (the patch protocol)

Each sub-milestone `<M>` is the block fenced by a line `` ````diff <M> `` (four backticks) and a line of four
backticks, in Appendix P. Paths are `a/<repo path>` and `b/<repo path>`; a new file's patch has `--- /dev/null`
and an `@@ -0,0 …` hunk. Apply it **only** this way, in one `bash` shell (the controller's shell is zsh):

```bash
bash <<'EOF2'
set -euo pipefail
M=M1a   # <- the sub-milestone, exactly as in the Appendix P heading: M1a M1b M1c M2 M3a M3b M3c M3d
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter389-a4-s1
PLAN="$W/design_docs/planned/w-world-approvals-by-key-a4-sprint1-plan.md"
T="$HOME/.ailang/state/world-iter389/exec-patches"; mkdir -p "$T"
awk -v m="$M" '$0 == "````diff " m {f=1; next} f && $0 == "````" {exit} f' "$PLAN" > "$T/$M.patch"
test -s "$T/$M.patch"
cd "$W"
patch -p1 --dry-run < "$T/$M.patch"
patch -p1 < "$T/$M.patch"
if find cmd host docs \( -name '*.orig' -o -name '*.rej' \) | grep .; then echo "STOP: patch left .orig/.rej"; exit 1; fi
git status --short
EOF2
```

- If `--dry-run` fails, **stop and report** (no hand edits, no fuzz options). It can only fail if an earlier
  sub-milestone was skipped or changed.
- Then run that sub-milestone's **Acceptance** (§3) in one bash shell with `export
  AILANG_BIN=$HOME/.pinned-ailang/ailang`, each command's output `tee`'d to a log under
  `$HOME/.ailang/state/world-iter389/exec-logs/` (never only through `tail`). **Never `nohup`, never background a
  test leg.**
- Then snapshot (§7). The controller commits and compares the numstat with §7's expected value.

---

## 3. Milestones

Every sub-milestone's acceptance starts with the same three lines (abbreviated **[base]** below):

```
test -z "$(gofmt -l cmd host)"
go build ./... && go vet ./host/approval/ ./host/broker/ ./host/coordinator/ ./host/surfacepage/ ./host/daemon/
go test -run '^$' ./...
```

Each `-v` run below must show `--- PASS` for every test it names and **no `--- FAIL` and no `--- SKIP`**; each
package run must end `ok` with no `FAIL` line. The planner's `tools/accept.sh <tree> <M>` is exactly these lines.

### M1a — the decision message · 69 prod + 93 test · Run 1

**Patch** `M1a`: new `host/approval/decision.go` — `DecisionPrefix = "world-decision/v1"`, `DecisionNamespace =
"world-decision@sunholo"`, `DecisionMessageLen = 91`, `VerbApprove = 0x01`, `VerbDeny = 0x02`, `ErrDecisionVerb`,
`DecisionVerb(word)`, `DecisionMessage(verb, P, Q, expiry)` (prefix ‖ 0x00 ‖ verb ‖ P ‖ Q ‖ u64be(expiry); verb
and negative expiry refused), `DecisionMessageHash`. New `decision_test.go`: DM1 `TestDecisionMessageLayout`
(91 bytes, every offset, golden hex and SHA-256 for two subjects computed by an independent encoder, domain
separation from `SignedMessage`/`Namespace`), DM2 `TestDecisionMessageRefusesBadInput`.

**Acceptance**
```
[base]
go test ./host/approval/... -run '^(TestDecisionMessage|TestApprovalImportsOnlyTheStandardLibrary)' -count=1 -v
go test ./host/approval/... -count=1
```
**Expected cumulative numstat after the controller's commit of M1a: 162**.

### M1b — the decision law and the reserved scope · 132 prod + 214 test · Run 1

**Patch** `M1b`: new `host/broker/decision_law.go` — `MaxDecisionWindow = 600`, `Discovery`, `DecisionFacts`,
`LawRefusal`, `DecisionWindowMax`, `SubjectRefusal` (D1–D7), `DecisionRefusal` (D1–D10), `DecisionDiscoveryLimit
= 2000`, `ReservedApprovalScope` (L1: publish parse / `registry:` / `effect=Registry.Publish`; finish parse /
`github-pr:`), `reservedReason`. New `decision_law_test.go`: `TestDecisionRefusalD1` … `D10`, `D4a`, `D4b`,
`TestDecisionRefusalOrder` (every pair of steps), `TestDecisionRefusalMustLive`, `TestReservedApprovalScope`.

**Acceptance**
```
[base]
go test ./host/broker/ -run '^(TestDecisionRefusal|TestReservedApprovalScope)' -count=1 -v
```
**Expected cumulative numstat after the controller's commit of M1b: 508**.

### M1c — the complete packet discovery · 63 prod + 162 test · Run 1

**Patch** `M1c`: new `host/coordinator/decisionpacket_discover.go` — `DiscoveryReader`,
`DiscoverRequestPackets` (S1-D3). New `decisionpacket_discover_test.go`: DS1
`TestDiscoverRequestPacketsCompleteness` over a fake keyset reader (0, 499, 500, 1,999, 2,000 complete; 2,001 not;
page counts and every cursor; two carrying across pages; one undecodable; one missing payload; a page error and a
payload-read error returned).

**Acceptance**
```
[base]
go test ./host/coordinator/ -run '^TestDiscoverRequestPackets' -count=1 -v
go test ./host/coordinator/ -count=1
```
**Expected cumulative numstat after the controller's commit of M1c: 733**.

### M2 — the broker rules · 36 prod + 235 test · Run 1

**Patch** `M2`: `host/broker/approve.go` — `ErrApprovalAlreadyDecided`; `decideApproval` refuses a second,
different decision and answers an identical one idempotently (S1-D10); `ApprovalChainState` (S1-D11); L2 in
`validatePublishApproval` (a `decidedBy` starting `ssh:` or `webauthn:` → `ErrPublishApprovalDenied`, "publishing
stays typed-only (D-WORLD-89 R4): decision <ref> was made by <decidedBy>"). New `decide_once_test.go`: BR1
`TestDecideApprovalRefusesASecondDecision`, BR2 `TestDecideApprovalFirstDecisionMustLive`, R4b
`TestPublishApprovalRefusesAKeyDecision`, R4c `TestPublishApprovalTypedStampMustLive`, XG1
`TestDecideApprovalCallerCensus` (AST census of the module's non-test calls). `context_thread_test.go`: the census
line `5/4` → `6/5` (P3).

**Acceptance**
```
[base]
go test ./host/broker/ -run '^(TestDecideApproval|TestPublishApproval|TestApprovalContextIdentity|TestDecisionRefusal|TestReservedApprovalScope)' -count=1 -v
go test ./host/daemon/ -run '^TestEpisodeFinish' -count=1 -v
```
(`TestPublishApproval` matches the existing `registry_publish_test.go` arms too: they must PASS unedited.)
**Expected cumulative numstat after the controller's commit of M2: 1004**.

### M3a — the decision routes, the guards, the typed arm · 287 prod + 65 test · Run 2

**Patch** `M3a`: new `host/surfacepage/decide_model.go` (`DecisionSubject`, `DecisionWindow`, `DecisionSSH`,
`DecisionMade`, `DecideBody`, `ErrDecideBody`, `ParseDecideBody`); new `host/daemon/surface_decide.go`
(`readDecision`, `lawClass`, `handleSurfaceDecision`, `handleSurfaceDecide`, S1-D4–D9); `host/daemon/daemon.go`
two `HandleFunc` lines behind `approvalRoute` + the doc comment; new `surface_decide_routes_test.go` (RD1); the
census gains `host/daemon/surface_decide.go:handleSurfaceDecide`.

**Acceptance**
```
[base]
go test ./host/daemon/ -run '^(TestDecisionRoutesRegistered|TestFrozenV1RouteTableMatchesMux|TestSurfaceAppRoutes|TestApprovalRoutesRegistered|TestDaemonCredentialWriterSites|TestReadCtxCancelledAfterHandler)$' -count=1 -v
go test ./host/broker/ -run '^TestDecideApprovalCallerCensus$' -count=1 -v
```
**Expected cumulative numstat after the controller's commit of M3a: 1356**.

### M3b — the rig and the subject · 272 test · Run 2

**Patch** `M3b`: new `surface_decide_rig_test.go` (`newDecideRig`, `packet`, `proof`, `decideBody`, `post`,
`decisions`, `fillPackets`, `packetAfter`, `packetPosition`); new `surface_decide_subject_test.go` (RD3).

**Acceptance**
```
[base]
go test ./host/daemon/ -run '^TestDecisionSubjectJSON$' -count=1 -v
```
**Expected cumulative numstat after the controller's commit of M3b: 1628**.

### M3c — the POST: guards, law arms, the proof, replay, R4 · 273 test · Run 2

**Patch** `M3c`: new `surface_decide_post_test.go`: RD2 `TestDecidePOSTGuards`, RD4 `TestDecidePOSTLawArms`, RD5
`TestDecideTypedLands`, RD6 `TestDecideTypedProofBindsTheVerb`, RD7 `TestDecidePOSTStaleDigests`, RD8
`TestDecideReplayAfterRestart`, R4d `TestDecideRefusesReservedScopesTyped`.

**Acceptance**
```
[base]
go test ./host/daemon/ -run '^(TestDecidePOST|TestDecideTyped|TestDecideReplay|TestDecideRefusesReserved)' -count=1 -v
```
**Expected cumulative numstat after the controller's commit of M3c: 1901**.

### M3d — the discovery end to end · 84 test · Run 2

**Patch** `M3d`: new `surface_decide_discovery_test.go`: RD9 `TestDecideRefusesASharedRequest`, RD9b
`TestDecideDiscoveryFindsAHiddenSecondPacket` (150 packets), RD9c `TestDecideDiscoveryRefusesAnOverflow` (2,001),
RD9d `TestDecideDiscoveryFindsTheOnlyPacketBeyondThePrefix` (150, must live).

**Acceptance**
```
[base]
go test ./host/daemon/ -run '^(TestDecideRefusesASharedRequest|TestDecideDiscovery)' -count=1 -v
go test ./host/daemon/ -count=1 -timeout 30m
```
**Expected cumulative numstat after the controller's commit of M3d: 1985**.

### 3.1 Reachability: every new refusal and branch, from the entry point (directive retro 4)

Entry point: `POST /surface/decide/{ref}` (the mux matches any one-segment `{ref}`) → `approvalRoute` (R2) →
`handleSurfaceDecide`; or `GET /surface/decisions/{ref}` → `approvalRoute` → `handleSurfaceDecision`. "E2E" = a
`host/daemon` test drives it through `d.Handler()`; "unit" = only a package test reaches it.

| Branch | Input class → trace | Driven by |
|---|---|---|
| R2 | foreign `Host` → approvalRoute 421 before the handler | E2E RD1, RD2 |
| R0 | `Authorization` (also empty) / `Cookie` → 403 before the body is read | E2E RD2 |
| R3 / R4 / R5 | no, `null`, other `Origin`; `Sec-Fetch-Site: same-site`; `text/plain` | E2E RD2 |
| R6 size / grammar | > 20 KiB → 413; unknown field, `decidedBy`, `sshsig`, `typed:false`, no format, `v:2`, upper-case or short hex, two values → 400 | E2E RD2 |
| D1 | `sha256:XYZ` → guards pass → `readDecision` reads nothing (`RefOK` false) → law D1 400 | E2E RD4, RD3 (GET) |
| D2 | absent ref; a packet's ref; a `Human.PollApproval`-effect request object → 404 | E2E RD4, RD3 |
| D3 | a request object stored but never on the chain → 404 | E2E RD4 |
| D4a sentinel | 2,001 stored packets, one carrying → 409 | E2E RD9c; unit DS1, DL4a |
| D4a undecodable | a stored packet that does not decode | **unit only** (DS1, DL4a): no E2E test stores an undecodable packet; the trace is the same `readDecision` → law path as the sentinel |
| D4b zero / two / refused packet | a chain request with no packet; two forged commits of one request (also one hidden past position 101); a failed-effect packet (validator code 4) → 409 | E2E RD4, RD3, RD9, RD9b |
| D5 | publish-shaped, `registry:` near miss, finish-shaped, `github-pr:` near miss → 403 | E2E R4d, RD3 (display) |
| D6 | a decided request (also after a restart) → 409 | E2E RD8, RD3 (display) |
| D7 | `now = deadline` → 409 | E2E RD4 |
| D8 | another packet's P; another request's Q → 409 | E2E RD4, RD7 |
| D9 | `"verb":"maybe"` → 400 | E2E RD4 |
| D10 | `expiry = now`; `expiry > deadline` → 409 | E2E RD4 |
| D11 typed | deny proof for approve; finish phase; wrong secret; no header; another request's proof → 403; a valid proof → 200 and the file unlinked | E2E RD5, RD6 |
| `DecideApproval` error after the law | **unreachable without a store fault**: D6 refuses every decided request under `humanMu` first (§1.4 P8); named survivor `H-DECIDE-ERR` | — |
| BR1 refuse | a second, different decision of a decided request | **unit only** (BR1): unreachable through the POST (D6 first, P8); reachable by `RecordOperatorApproval`/`mintAttendedApproval` only with a different payload for a colliding request |
| BR1 identical → idempotent | the finish route re-approving one scope at one second | E2E `TestEpisodeFinishOperatorProofExpiryBoundary` (existing, unedited); unit BR1 |
| L2 | a publish payload naming a key-decided approval → `Session.Invoke` → `validatePublishApproval` | **unit** (R4b, R4c at `validatePublishApproval`); E2E needs `world-publish publish` with `--decided-by ssh:…`, outside this sprint's files (`cmd/world-publish/*` untouched) |
| GET D1–D4 vs 200 | as D1–D4 above → refusal; D5–D7 → 200 with the display refusal | E2E RD3 |

---

## 4. Mutation table

**96 rows on the planner's final tree (`work` = M3d; the last ledger line of each row, runs `run2`–`run4`): 88 KILLED (named and suite) as expected, 8 `E` rows (3 named equivalents: `RES-FINISH-PARSE`, `DIS-NOTFOUND`, `R-DISC-GATE`; 5 named survivors with reasons: `BR1-READ-ERR`, `H-LOCK`, `H-DECIDE-ERR`, `R-D2-SEMANTIC`, `R-DISCOVERY-ERR`), 0 not as expected, 0 `NOCOMPILE` among the last lines.** 104 ledger lines in all (`plan-scratch/logs/mutant_ledger.tsv`; runs: run1 4, run2 96, run3 3, run4 1). Run `run1` (4 rows, the first 16 minutes) was stopped and re-run: its daemon suite (`TestApprov*`, `TestRequest*`, `TestWorkbenchPacket*`) took ≈ 980 s per mutant under load, and its runner treated a named test that the suite's `-skip` excluded as run (the two `ROUTE-*` rows read `SURVIVED` for the named RD1, which kills both when run, `logs/`); `run2` narrows the daemon and broker suites to the sprint's own tests plus the existing finish-proof tests and runs a skipped named test on its own. Its four lines stay in the ledger. Four rows did not compile in `run2` (`H-POSTED-P`, `H-POSTED-Q`, `G-RESERVED`: an unused variable; `DIS-NOTFOUND`: `ok` unused) and are **not counted**: each was rewritten to compile (`; _ = packet`, `; _ = request`, a constant scope argument, `if _ = ok; …`) and re-run (`run3`, `run4`); the NOCOMPILE lines stay in the ledger. **No kill rests on an inventory or static test**: the daemon suite skips the static RD1 (`TestDecisionRoutesRegistered`) and the broker suite the census XG1 (`TestDecideApprovalCallerCensus`); the two `ROUTE-*` rows name RD1, run it on its own (KILLED) and are KILLED by the behavioural suite too (RD3, RD2, RD4… answer 404/405 for a missing route); every other kill's failing-test column names the row's behavioural test.

**Derivation.** Anchored to Appendix P's production hunks, one row per new branch, enumerated from the diff
(`git diff 9a6b959d M3d -- host ':!*_test.go'`): the message's verb and expiry refusals, its separator byte,
field order, endianness and prefix, the word map (M1a); each law step deleted (the unit is the branch), D4a
reading `Carrying` instead of `Complete` and dropping `Undecodable`, D4b `Carrying >= 1` and the validator code,
each edge (D7 `>=`, D10 both ends, the window's deadline), the subject-first call, the owner text, each L1 arm
(M1b); the discovery's sentinel (`>=` → `>` and deleted), cursor, undecodable count, not-found arm, carrying
test, short-page edge and both error returns (M1c); BR1's refusal, its idempotent arm both ways, its read error,
the chain state's three branches, each L2 prefix, L2 matching everything and a loosened prefix (M2); both
routes, R0 (both headers), R3–R6, the grammar's strictness, one-value, `v`, `typed` and hex checks, the law
call, both posted digests, the proof consumption, the phase's verb, the proof's name, the derived `decidedBy`,
the decision's time, the lock, the decide error, D1's guard, D2's effect and semantic checks, the discovery gate,
the packet code, the deadline and digest facts, the discovery error, GET's no-subject set, window, reserved,
refusal and bytes, and the refusal classes (M3a).
**Procedure**: `plan-scratch/tools/mutate.py <run> <tree> <workers>` (the table is `tools/mutants.py`; every
anchor is asserted to occur exactly once in the tree it edits): for each row, in its own worker copy of the tree
(`cp -cR`), apply **one** edit, run the package's suite with `-v` (which includes the named test, or the named
test alone first when the suite excludes it), restore, assert the restore byte-identical, and append one ledger
line (UTC time, run, tree, id, stage, file, from → to, named test, named, suite, the suite's failing tests,
seconds, expected). `GOFLAGS=-vet=off`; a mutant that does not compile is `NOCOMPILE`, never counted. Suites
(inventories and static tests excluded): `host/approval` whole minus `TestApprovalImportsOnlyTheStandardLibrary`;
`host/coordinator` whole; `host/broker` `^(TestDecision|TestDecide|TestReserved|TestPublishApproval|
TestApprovalContextIdentity)` minus `TestDecideApprovalCallerCensus`; `host/daemon` `^(TestDecision|TestDecide|
TestEpisodeFinish|TestApprov|TestRequest|TestWorkbenchPacket)` minus `TestDecisionRoutesRegistered`,
`TestApprovalRoutesRegistered`, `TestDaemonCredentialWriterSites`. A row is KILLED only when both fail. The two
route rows name the static RD1, which the daemon suite excludes, so their suite verdict is behavioural (RD3, RD2…
fail on the 404/405 a missing route answers). **Never commit or snapshot a mutant.** In the table a newline
prints as ⏎ and a pipe as \|.

| Id | Stage | file: from → to | Named test | Named | Suite | Suite's failing tests | Expected |
|---|---|---|---|---|---|---|---|
| DM-VERB-CHECK | M1a | decision.go: if verb != VerbApprove && verb != VerbDeny { → if false { | `TestDecisionMessageRefusesBadInput` | KILLED | KILLED | TestDecisionMessageRefusesBadInput | K |
| DM-NEG-EXPIRY | M1a | decision.go: if expiry < 0 { → if false { | `TestDecisionMessageRefusesBadInput` | KILLED | KILLED | TestDecisionMessageRefusesBadInput | K |
| DM-ZERO-BYTE | M1a | decision.go: m = append(m, 0)⏎	m = append(m, verb) → m = append(m, 1)⏎	m = append(m, verb) | `TestDecisionMessageLayout` | KILLED | KILLED | TestDecisionMessageLayout | K |
| DM-FIELD-ORDER | M1a | decision.go: m = append(m, packet[:]...)⏎	m = append(m, request[:]...) → m = append(m, request[:]...)⏎	m = append(m, packet[:]...) | `TestDecisionMessageLayout` | KILLED | KILLED | TestDecisionMessageLayout | K |
| DM-LE-EXPIRY | M1a | decision.go: binary.BigEndian.AppendUint64(m, uint64(expiry)) → binary.LittleEndian.AppendUint64(m, uint64(expiry)) | `TestDecisionMessageLayout` | KILLED | KILLED | TestDecisionMessageLayout | K |
| DM-PREFIX | M1a | decision.go: const DecisionPrefix = "world-decision/v1" → const DecisionPrefix = "world-approval/v1" | `TestDecisionMessageLayout` | KILLED | KILLED | TestDecisionMessageLayout | K |
| DM-HASH-ERR | M1a | decision.go: m, err := DecisionMessage(verb, packet, request, expiry)⏎	if err != nil { → m, err := DecisionMessage(verb, packet, request, expiry)⏎	if false { | `TestDecisionMessageRefusesBadInput` | KILLED | KILLED | TestDecisionMessageRefusesBadInput | K |
| DM-WORD-APPROVE | M1a | decision.go: case "approve":⏎		return VerbApprove, nil → case "approve":⏎		return VerbDeny, nil | `TestDecisionMessageLayout` | KILLED | KILLED | TestDecisionMessageLayout | K |
| DM-WORD-DENY | M1a | decision.go: case "deny": → case "denied": | `TestDecisionMessageLayout` | KILLED | KILLED | TestDecisionMessageLayout | K |
| LAW-D1 | M1b | decision_law.go: case !f.RefOK: → case false: | `TestDecisionRefusalD1` | KILLED | KILLED | TestDecisionRefusalD1,TestDecisionRefusalOrder | K |
| LAW-D2 | M1b | decision_law.go: case !f.RequestOK: → case false: | `TestDecisionRefusalD2` | KILLED | KILLED | TestDecisionRefusalD2,TestDecisionRefusalOrder | K |
| LAW-D3 | M1b | decision_law.go: case !f.OnChain: → case false: | `TestDecisionRefusalD3` | KILLED | KILLED | TestDecisionRefusalD3,TestDecisionRefusalOrder | K |
| LAW-D4A | M1b | decision_law.go: case !f.Discovery.Complete \|\| f.Discovery.Undecodable > 0: → case false: | `TestDecisionRefusalD4a` | KILLED | KILLED | TestDecisionRefusalD4a,TestDecisionRefusalOrder | K |
| LAW-D4A-READS-CARRYING | M1b | decision_law.go: case !f.Discovery.Complete \|\| f.Discovery.Undecodable > 0: → case f.Discovery.Carrying != 1 \|\| f.Discovery.Undecodable > 0: | `TestDecisionRefusalD4a` | KILLED | KILLED | TestDecisionRefusalD4a,TestDecisionRefusalD4b,TestDecisionRefusalOrder | K |
| LAW-D4A-UNDECODABLE | M1b | decision_law.go: case !f.Discovery.Complete \|\| f.Discovery.Undecodable > 0: → case !f.Discovery.Complete: | `TestDecisionRefusalD4a` | KILLED | KILLED | TestDecisionRefusalD4a | K |
| LAW-D4B-ZERO | M1b | decision_law.go: case f.Discovery.Carrying == 0: → case false: | `TestDecisionRefusalD4b` | KILLED | KILLED | TestDecisionRefusalD4b | K |
| LAW-D4B-ATLEAST1 | M1b | decision_law.go: case f.Discovery.Carrying != 1: → case f.Discovery.Carrying < 1: | `TestDecisionRefusalD4b` | KILLED | KILLED | TestDecisionRefusalD4b,TestDecisionRefusalOrder | K |
| LAW-D4B-CODE | M1b | decision_law.go: case f.PacketCode != 0: → case false: | `TestDecisionRefusalD4b` | KILLED | KILLED | TestDecisionRefusalD4b | K |
| LAW-D5 | M1b | decision_law.go: ReservedApprovalScope(f.Scope); reserved { → ReservedApprovalScope(f.Scope); reserved && false { | `TestDecisionRefusalD5` | KILLED | KILLED | TestDecisionRefusalD5,TestDecisionRefusalOrder | K |
| LAW-D6 | M1b | decision_law.go: case f.Decided: → case false: | `TestDecisionRefusalD6` | KILLED | KILLED | TestDecisionRefusalD6,TestDecisionRefusalOrder | K |
| LAW-D7 | M1b | decision_law.go: case f.Now >= f.Deadline: → case false: | `TestDecisionRefusalD7` | KILLED | KILLED | TestDecisionRefusalD7,TestDecisionRefusalOrder | K |
| LAW-D7-EDGE | M1b | decision_law.go: case f.Now >= f.Deadline: → case f.Now > f.Deadline: | `TestDecisionRefusalD7` | KILLED | KILLED | TestDecisionRefusalD7 | K |
| LAW-D8 | M1b | decision_law.go: case f.PostedPacket != f.Packet \|\| f.PostedRequest != f.Request: → case false: | `TestDecisionRefusalD8` | KILLED | KILLED | TestDecisionRefusalD8,TestDecisionRefusalOrder | K |
| LAW-D8-P | M1b | decision_law.go: case f.PostedPacket != f.Packet \|\| f.PostedRequest != f.Request: → case f.PostedRequest != f.Request: | `TestDecisionRefusalD8` | KILLED | KILLED | TestDecisionRefusalD8,TestDecisionRefusalOrder | K |
| LAW-D8-Q | M1b | decision_law.go: case f.PostedPacket != f.Packet \|\| f.PostedRequest != f.Request: → case f.PostedPacket != f.Packet: | `TestDecisionRefusalD8` | KILLED | KILLED | TestDecisionRefusalD8 | K |
| LAW-D9 | M1b | decision_law.go: case f.Verb != "approve" && f.Verb != "deny": → case false: | `TestDecisionRefusalD9` | KILLED | KILLED | TestDecisionRefusalD9,TestDecisionRefusalOrder | K |
| LAW-D10 | M1b | decision_law.go: case f.Expiry <= f.Now \|\| f.Expiry > DecisionWindowMax(f.Now, f.Deadline): → case false: | `TestDecisionRefusalD10` | KILLED | KILLED | TestDecisionRefusalD10 | K |
| LAW-D10-LOW | M1b | decision_law.go: case f.Expiry <= f.Now \|\| f.Expiry > DecisionWindowMax(f.Now, f.Deadline): → case f.Expiry < f.Now \|\| f.Expiry > DecisionWindowMax(f.Now, f.Deadline): | `TestDecisionRefusalD10` | KILLED | KILLED | TestDecisionRefusalD10 | K |
| LAW-D10-HIGH | M1b | decision_law.go: case f.Expiry <= f.Now \|\| f.Expiry > DecisionWindowMax(f.Now, f.Deadline): → case f.Expiry <= f.Now \|\| f.Expiry >= DecisionWindowMax(f.Now, f.Deadline): | `TestDecisionRefusalD10` | KILLED | KILLED | TestDecisionRefusalD10,TestDecisionRefusalD7,TestDecisionRefusalMustLive | K |
| LAW-WINDOW-DEADLINE | M1b | decision_law.go: return min(now+MaxDecisionWindow, deadline) → return now + MaxDecisionWindow | `TestDecisionRefusalD10` | KILLED | KILLED | TestDecisionRefusalD10 | K |
| LAW-SUBJECT-FIRST | M1b | decision_law.go: if r := SubjectRefusal(f); r.Step != "" { → if r := SubjectRefusal(f); false && r.Step != "" { | `TestDecisionRefusalD1` | KILLED | KILLED | TestDecisionRefusalD1,TestDecisionRefusalD2,TestDecisionRefusalD3,TestDecisionRefusalD4a,TestDecisionRefusalD4b,TestDecisionRefusalD5,TestDecisionRefusalD6,Test | K |
| LAW-OWNER-TEXT | M1b | decision_law.go: func reservedReason(owner string) string {⏎	if owner == "publish" { → func reservedReason(owner string) string {⏎	if owner != "publish" { | `TestDecisionRefusalD5` | KILLED | KILLED | TestDecisionRefusalD5 | K |
| RES-PUBLISH-PARSE | M1b | decision_law.go: if _, err := ParsePublishApprovalScope(scope); err == nil \|\| → if _, err := ParsePublishApprovalScope(scope); err == nil && false \|\| | `TestReservedApprovalScope` | KILLED | KILLED | TestReservedApprovalScope | K |
| RES-REGISTRY-PREFIX | M1b | decision_law.go: strings.HasPrefix(scope, "registry:") → false | `TestReservedApprovalScope` | KILLED | KILLED | TestDecisionRefusalD5,TestReservedApprovalScope | K |
| RES-EFFECT-ARM | M1b | decision_law.go: strings.Contains(scope, "effect="+EffectRegistryPublish) → false | `TestReservedApprovalScope` | KILLED | KILLED | TestReservedApprovalScope | K |
| RES-FINISH-PARSE | M1b | decision_law.go: if _, err := ParseFinishApprovalScope(scope); err == nil \|\| → if _, err := ParseFinishApprovalScope(scope); err == nil && false \|\| | `TestReservedApprovalScope` | SURVIVED | SURVIVED | — | E: equivalent: `finishScopePattern` (approve.go) starts `^github-pr:`, so every scope the finish parse arm accepts the `github-pr:` prefix arm accepts too; kept as the parser's own word if the grammar ever changes |
| RES-FINISH-PREFIX | M1b | decision_law.go: strings.HasPrefix(scope, "github-pr:") → false | `TestReservedApprovalScope` | KILLED | KILLED | TestDecisionRefusalD5,TestDecisionRefusalOrder,TestReservedApprovalScope | K |
| DIS-SENTINEL-EDGE | M1c | decisionpacket_discover.go: if d.Scanned >= broker.DecisionDiscoveryLimit { → if d.Scanned > broker.DecisionDiscoveryLimit { | `TestDiscoverRequestPacketsCompleteness` | KILLED | KILLED | TestDiscoverRequestPacketsCompleteness | K |
| DIS-SENTINEL-DEL | M1c | decisionpacket_discover.go: if d.Scanned >= broker.DecisionDiscoveryLimit { → if false { | `TestDiscoverRequestPacketsCompleteness` | KILLED | KILLED | TestDiscoverRequestPacketsCompleteness | K |
| DIS-CURSOR | M1c | decisionpacket_discover.go: after = page[len(page)-1].Hash.String() → after = "" | `TestDiscoverRequestPacketsCompleteness` | KILLED | KILLED | TestDiscoverRequestPacketsCompleteness | K |
| DIS-UNDECODABLE | M1c | decisionpacket_discover.go: d.Undecodable++⏎ → ⏎ | `TestDiscoverRequestPacketsCompleteness` | KILLED | KILLED | TestDiscoverRequestPacketsCompleteness | K |
| DIS-NOTFOUND | M1c | decisionpacket_discover.go: if !ok \|\| decodeErr != nil { → if _ = ok; decodeErr != nil { | `TestDiscoverRequestPacketsCompleteness` | SURVIVED | SURVIVED | — | E: equivalent: a payload `GetObject` did not find is the zero Object, whose nil payload never decodes, so `decodeErr != nil` already counts it; kept to name the case |
| DIS-CARRYING | M1c | decisionpacket_discover.go: if p.RequestRef == request { → if p.RequestRef != request { | `TestDiscoverRequestPacketsCompleteness` | KILLED | KILLED | TestDiscoverRequestPacketsCompleteness | K |
| DIS-COMPLETE-EDGE | M1c | decisionpacket_discover.go: if len(page) < store.MaxSemanticIDPage { → if len(page) <= store.MaxSemanticIDPage { | `TestDiscoverRequestPacketsCompleteness` | KILLED | KILLED | TestDiscoverRequestPacketsCompleteness | K |
| DIS-PAGE-ERR | M1c | decisionpacket_discover.go: page, err := r.ObjectsBySemanticID(ctx, DecisionPacketV1, after, store.MaxSemanticIDPage)⏎		if err != nil { → page, err := r.ObjectsBySemanticID(ctx, DecisionPacketV1, after, store.MaxSemanticIDPage)⏎		if false && err != | `TestDiscoverRequestPacketsCompleteness` | KILLED | KILLED | TestDiscoverRequestPacketsCompleteness | K |
| DIS-GET-ERR | M1c | decisionpacket_discover.go: full, ok, err := r.GetObject(ctx, o.Hash)⏎			if err != nil { → full, ok, err := r.GetObject(ctx, o.Hash)⏎			if false && err != nil { | `TestDiscoverRequestPacketsCompleteness` | KILLED | KILLED | TestDiscoverRequestPacketsCompleteness | K |
| BR1-REFUSE | M2 | approve.go: } else if decided {⏎		return hashref.HashRef{}, ErrApprovalAlreadyDecided → } else if false {⏎		return hashref.HashRef{}, ErrApprovalAlreadyDecided | `TestDecideApprovalRefusesASecondDecision` | KILLED | KILLED | TestDecideApprovalRefusesASecondDecision | K |
| BR1-ANY-IS-IDENTICAL | M2 | approve.go: } else if decided && bytes.Equal(existing.Payload, payload) { → } else if decided { | `TestDecideApprovalRefusesASecondDecision` | KILLED | KILLED | TestDecideApprovalRefusesASecondDecision | K |
| BR1-NO-IDEMPOTENCE | M2 | approve.go: } else if decided && bytes.Equal(existing.Payload, payload) { → } else if false { | `TestDecideApprovalRefusesASecondDecision` | KILLED | KILLED | TestDecideApprovalRefusesASecondDecision | K |
| BR1-READ-ERR | M2 | approve.go: if existing, decided, err := findApprovalDecision(ctx, s, requestRef); err != nil { → if existing, decided, err := findApprovalDecision(ctx, s, requestRef); false && err != nil { | `TestDecideApprovalRefusesASecondDecision` | SURVIVED | SURVIVED | — | E: survivor (untested defensive arm): `findApprovalDecision` fails only on a store error after `findApprovalRequest` read the same chain successfully; no host/broker fake injects a fault between the two walks. Mutated, a read error would read as undecided and the write would proceed |
| CHAIN-ONCHAIN | M2 | approve.go: err != nil \|\| !onChain { → err != nil { | `TestDecideApprovalFirstDecisionMustLive` | KILLED | KILLED | TestDecideApprovalFirstDecisionMustLive | K |
| CHAIN-UNDECIDED | M2 | approve.go: if err != nil \|\| !decided {⏎		return true, false → if err != nil {⏎		return true, false | `TestDecideApprovalFirstDecisionMustLive` | KILLED | KILLED | TestDecideApprovalFirstDecisionMustLive | K |
| CHAIN-VERB-BY | M2 | approve.go: return true, true, w.Decision, w.DecidedBy, nil → return true, true, w.DecidedBy, w.Decision, nil | `TestDecideApprovalRefusesASecondDecision` | KILLED | KILLED | TestDecideApprovalFirstDecisionMustLive,TestDecideApprovalRefusesASecondDecision | K |
| L2-SSH | M2 | approve.go: if strings.HasPrefix(decision.DecidedBy, "ssh:") \|\| → if false \|\| | `TestPublishApprovalRefusesAKeyDecision` | KILLED | KILLED | TestPublishApprovalRefusesAKeyDecision | K |
| L2-WEBAUTHN | M2 | approve.go: \|\| strings.HasPrefix(decision.DecidedBy, "webauthn:") { → \|\| false { | `TestPublishApprovalRefusesAKeyDecision` | KILLED | KILLED | TestPublishApprovalRefusesAKeyDecision | K |
| L2-EVERY | M2 | approve.go: if strings.HasPrefix(decision.DecidedBy, "ssh:") \|\| strings.HasPrefix(decision.DecidedBy, "webauthn:") { → if true { | `TestPublishApprovalTypedStampMustLive` | KILLED | KILLED | TestApprovalContextIdentity,TestApprovalContextIdentity/bounded,TestPublishApprovalScopeRefusesASmuggledFragment,TestPublishApprovalTypedStampMustLive | K |
| L2-LOOSE-PREFIX | M2 | approve.go: strings.HasPrefix(decision.DecidedBy, "ssh:") → strings.HasPrefix(decision.DecidedBy, "ssh") | `TestPublishApprovalTypedStampMustLive` | KILLED | KILLED | TestPublishApprovalTypedStampMustLive | K |
| ROUTE-GET | M3a | daemon.go: 	mux.HandleFunc("GET /surface/decisions/{ref}", d.approvalRoute(d.handleSurfaceDecision))⏎ → (deleted) | `TestDecisionRoutesRegistered` | KILLED | KILLED | TestDecisionSubjectJSON | K |
| ROUTE-POST | M3a | daemon.go: 	mux.HandleFunc("POST /surface/decide/{ref}", d.approvalRoute(d.handleSurfaceDecide))⏎ → (deleted) | `TestDecisionRoutesRegistered` | KILLED | KILLED | TestDecideDiscoveryFindsAHiddenSecondPacket,TestDecideDiscoveryFindsTheOnlyPacketBeyondThePrefix,TestDecideDiscoveryRefusesAnOverflow,TestDecidePOSTGuards,TestD | K |
| H-R0-AUTH | M3a | surface_decide.go: if len(r.Header.Values("Authorization")) != 0 \|\| → if false \|\| | `TestDecidePOSTGuards` | KILLED | KILLED | TestDecidePOSTGuards | K |
| H-R0-COOKIE | M3a | surface_decide.go: \|\| len(r.Header.Values("Cookie")) != 0 { // R0 → \|\| false { // R0 | `TestDecidePOSTGuards` | KILLED | KILLED | TestDecidePOSTGuards | K |
| H-R3 | M3a | surface_decide.go: if origin := r.Header.Get("Origin"); !allowedApprovalOrigin(origin, o) { // R3 → if origin := r.Header.Get("Origin"); false && !allowedApprovalOrigin(origin, o) { // R3 | `TestDecidePOSTGuards` | KILLED | KILLED | TestDecidePOSTGuards | K |
| H-R4 | M3a | surface_decide.go: if !sameOriginFetch(r.Header.Get("Sec-Fetch-Site")) { // R4 → if false { // R4 | `TestDecidePOSTGuards` | KILLED | KILLED | TestDecidePOSTGuards | K |
| H-R5 | M3a | surface_decide.go: if !jsonMediaType(r.Header.Get("Content-Type")) { // R5 → if false { // R5 | `TestDecidePOSTGuards` | KILLED | KILLED | TestDecidePOSTGuards | K |
| H-R6-SIZE | M3a | surface_decide.go: io.ReadAll(http.MaxBytesReader(w, r.Body, surfacepage.MaxApproveBodyBytes)) // R6 → io.ReadAll(r.Body) // R6 | `TestDecidePOSTGuards` | KILLED | KILLED | TestDecidePOSTGuards | K |
| H-R6-PARSE | M3a | surface_decide.go: body, packet, request, err := surfacepage.ParseDecideBody(raw)⏎	if err != nil { → body, packet, request, err := surfacepage.ParseDecideBody(raw)⏎	if false && err != nil { | `TestDecidePOSTGuards` | KILLED | KILLED | TestDecidePOSTGuards | K |
| M-STRICT | M3a | decide_model.go: 	dec.DisallowUnknownFields()⏎ → ⏎ | `TestDecidePOSTGuards` | KILLED | KILLED | TestDecidePOSTGuards | K |
| M-TWO-VALUES | M3a | decide_model.go: if _, err := dec.Token(); err != io.EOF { → if _, err := dec.Token(); false && err != io.EOF { | `TestDecidePOSTGuards` | KILLED | KILLED | TestDecidePOSTGuards | K |
| M-V | M3a | decide_model.go: case body.V != 1: → case false: | `TestDecidePOSTGuards` | KILLED | KILLED | TestDecidePOSTGuards | K |
| M-TYPED | M3a | decide_model.go: case !body.Typed: → case false: | `TestDecidePOSTGuards` | KILLED | KILLED | TestDecidePOSTGuards | K |
| M-HEX | M3a | decide_model.go: case !lowerHex64(body.Packet) \|\| !lowerHex64(body.Request): → case false: | `TestDecidePOSTGuards` | KILLED | KILLED | TestDecidePOSTGuards | K |
| H-LAW | M3a | surface_decide.go: if law := broker.DecisionRefusal(f); law.Step != "" { // 1. the law → if law := broker.DecisionRefusal(f); false && law.Step != "" { // 1. the law | `TestDecidePOSTLawArms` | KILLED | KILLED | TestDecideDiscoveryFindsAHiddenSecondPacket,TestDecideDiscoveryRefusesAnOverflow,TestDecidePOSTLawArms,TestDecideRefusesASharedRequest | K |
| H-POSTED-P | M3a | surface_decide.go: f.PostedPacket, f.PostedRequest, f.Verb, f.Expiry = packet, request, body.Verb, body.Expiry → f.PostedPacket, f.PostedRequest, f.Verb, f.Expiry = f.Packet, request, body.Verb, body.Expiry; _ = packet | `TestDecidePOSTStaleDigests` | KILLED | KILLED | TestDecidePOSTLawArms,TestDecidePOSTStaleDigests | K |
| H-POSTED-Q | M3a | surface_decide.go: f.PostedPacket, f.PostedRequest, f.Verb, f.Expiry = packet, request, body.Verb, body.Expiry → f.PostedPacket, f.PostedRequest, f.Verb, f.Expiry = packet, f.Request, body.Verb, body.Expiry; _ = request | `TestDecidePOSTStaleDigests` | KILLED | KILLED | TestDecidePOSTLawArms,TestDecidePOSTStaleDigests | K |
| H-PROOF-SKIPPED | M3a | surface_decide.go: r.Header.Get("X-World-Operator-Proof")); err != nil { → r.Header.Get("X-World-Operator-Proof")); false && err != nil { | `TestDecideTypedProofBindsTheVerb` | KILLED | KILLED | TestDecideTypedProofBindsTheVerb | K |
| H-PHASE-NO-VERB | M3a | surface_decide.go: "decide-"+body.Verb, r.Header → "decide-approve", r.Header | `TestDecideTypedLands` | KILLED | KILLED | TestDecideDiscoveryFindsTheOnlyPacketBeyondThePrefix,TestDecideReplayAfterRestart,TestDecideTypedLands | K |
| H-PROOF-NAME | M3a | surface_decide.go: d.consumeOperatorProof("decision-"+q.Digest(), → d.consumeOperatorProof("decision", | `TestDecideTypedLands` | KILLED | KILLED | TestDecideDiscoveryFindsTheOnlyPacketBeyondThePrefix,TestDecidePOSTGuards,TestDecideReplayAfterRestart,TestDecideTypedLands | K |
| H-DECIDEDBY | M3a | surface_decide.go: const typedDecider = "operator:terminal" → const typedDecider = "operator:surface" | `TestDecideTypedLands` | KILLED | KILLED | TestDecideReplayAfterRestart,TestDecideTypedLands | K |
| H-DECISION-NOW | M3a | surface_decide.go: broker.DecideApproval(ctx, d.store, q, body.Verb, typedDecider, now) → broker.DecideApproval(ctx, d.store, q, body.Verb, typedDecider, 0) | `TestDecideTypedLands` | KILLED | KILLED | TestDecideTypedLands | K |
| H-LOCK | M3a | surface_decide.go: 	d.humanMu.Lock()⏎	defer d.humanMu.Unlock()⏎ → ⏎ | `TestDecideTypedLands` | SURVIVED | SURVIVED | — | E: survivor (not observable sequentially): every sprint-1 test posts one request at a time; the lock is `humanMu`, shared with `serializedHandler` and the finish route (row 172/236); a concurrency test is carried to sprint 2 (§8 R4, Appendix C) |
| H-DECIDE-ERR | M3a | surface_decide.go: typedDecider, now) // 3.⏎	if err != nil { → typedDecider, now) // 3.⏎	if false && err != nil { | `TestDecideReplayAfterRestart` | SURVIVED | SURVIVED | — | E: survivor (unreachable without a store fault): D6 refuses every decided request under `humanMu` before `DecideApproval` (§1.4 P8), and no daemon test injects a world-store write failure (sprint 2's RD13 adds that seam) |
| R-D1 | M3a | surface_decide.go: if !decisionRefPattern.MatchString(ref) { // D1 → if false { // D1 | `TestDecidePOSTLawArms` | KILLED | KILLED | TestDecidePOSTLawArms | K |
| R-D2-EFFECT | M3a | surface_decide.go: err == nil && req.Effect == broker.EffectHumanApprove { // D2 → err == nil { // D2 | `TestDecidePOSTLawArms` | KILLED | KILLED | TestDecidePOSTLawArms | K |
| R-D2-SEMANTIC | M3a | surface_decide.go: ok && obj.SemanticID == broker.ApprovalRequestV1 && err == nil → ok && err == nil | `TestDecidePOSTLawArms` | SURVIVED | SURVIVED | — | E: survivor (content addressing): an object whose payload decodes as a Human.Approve request has that request's hash, so it is the request itself or it is not on the approvals chain (D3 then refuses with 404 and its own reason); no test can store a second object under the same hash |
| R-DISC-GATE | M3a | surface_decide.go: if f.Discovery.Complete && f.Discovery.Undecodable == 0 && f.Discovery.Carrying == 1 { // D4b, D7, D8 → if f.Discovery.Carrying == 1 { // D4b, D7, D8 | `TestDecideDiscoveryRefusesAnOverflow` | SURVIVED | SURVIVED | — | E: equivalent: with an incomplete or undecodable discovery the law refuses at D4a before it reads any packet fact; reading the one carrying packet anyway changes no response |
| R-PACKET-CODE | M3a | surface_decide.go: f.PacketCode, f.PacketSentence, f.Deadline, f.Packet = res.Code, → f.PacketCode, f.PacketSentence, f.Deadline, f.Packet = 0, | `TestDecidePOSTLawArms` | KILLED | KILLED | TestDecidePOSTLawArms | K |
| R-DEADLINE | M3a | surface_decide.go: packetRefusalSentence(res), res.Packet.DeadlineAt, [32]byte(pd) → packetRefusalSentence(res), res.Packet.DeadlineAt+1, [32]byte(pd) | `TestDecidePOSTLawArms` | KILLED | KILLED | TestDecidePOSTLawArms | K |
| R-PACKET-DIGEST | M3a | surface_decide.go: pd, _ := hex.DecodeString(carrying.Hash.Digest()) → pd, _ := hex.DecodeString(q.Digest()) | `TestDecideTypedLands` | KILLED | KILLED | TestDecideDiscoveryFindsTheOnlyPacketBeyondThePrefix,TestDecidePOSTGuards,TestDecidePOSTLawArms,TestDecideReplayAfterRestart,TestDecideTypedLands,TestDecideType | K |
| R-DISCOVERY-ERR | M3a | surface_decide.go: f.Discovery, carrying, err = coordinator.DiscoverRequestPackets(ctx, d.store, q); err != nil { → f.Discovery, carrying, err = coordinator.DiscoverRequestPackets(ctx, d.store, q); false && err != nil { | `TestDecideDiscoveryRefusesAnOverflow` | SURVIVED | SURVIVED | — | E: survivor (unreachable without a store fault): the discovery's store errors are proven at the coordinator (DS1's page and payload-read errors, rows DIS-PAGE-ERR and DIS-GET-ERR); the daemon has no fault seam for `ObjectsBySemanticID` under the POST |
| G-NO-SUBJECT-D4B | M3a | surface_decide.go: case "D1", "D2", "D3", "D4a", "D4b": → case "D1", "D2", "D3", "D4a": | `TestDecisionSubjectJSON` | KILLED | KILLED | TestDecisionSubjectJSON | K |
| G-WINDOW | M3a | surface_decide.go: Max: broker.DecisionWindowMax(now, p.DeadlineAt) → Max: now + 600 | `TestDecisionSubjectJSON` | KILLED | KILLED | TestDecisionSubjectJSON | K |
| G-RESERVED | M3a | surface_decide.go: reserved, _ := broker.ReservedApprovalScope(rd.facts.Scope) → reserved, _ := broker.ReservedApprovalScope("worktree") | `TestDecisionSubjectJSON` | KILLED | KILLED | TestDecisionSubjectJSON | K |
| G-REFUSAL | M3a | surface_decide.go: Refusal: law.Reason, → Refusal: "", | `TestDecisionSubjectJSON` | KILLED | KILLED | TestDecisionSubjectJSON | K |
| G-REQUEST-BYTES | M3a | surface_decide.go: RequestBytes: base64.RawURLEncoding.EncodeToString(rd.request.Payload) → RequestBytes: base64.RawURLEncoding.EncodeToString(rd.packet.Payload) | `TestDecisionSubjectJSON` | KILLED | KILLED | TestDecisionSubjectJSON | K |
| LAWCLASS-FORBIDDEN | M3a | surface_decide.go: case http.StatusForbidden:⏎		return "Forbidden" → case http.StatusForbidden:⏎		return "Conflict" | `TestDecideRefusesReservedScopesTyped` | KILLED | KILLED | TestDecideRefusesReservedScopesTyped | K |
| LAWCLASS-NOTFOUND | M3a | surface_decide.go: case http.StatusNotFound:⏎		return "NotFound" → case http.StatusNotFound:⏎		return "Conflict" | `TestDecidePOSTLawArms` | KILLED | KILLED | TestDecidePOSTLawArms,TestDecisionSubjectJSON | K |

---

## 5. Files (final `git status --short` after M3d)

| File | Sub-milestone | Change |
|---|---|---|
| `host/approval/decision.go`, `decision_test.go` | M1a | new: the message; DM1, DM2 |
| `host/broker/decision_law.go`, `decision_law_test.go` | M1b | new: the law, L1; DL*, R4a |
| `host/coordinator/decisionpacket_discover.go`, `decisionpacket_discover_test.go` | M1c | new: the discovery; DS1 |
| `host/broker/approve.go` | M2 | BR1 (+ identical idempotent), `ApprovalChainState`, L2 |
| `host/broker/decide_once_test.go` | M2, M3a | new: BR1, BR2, R4b, R4c, XG1; M3a adds the handler to the census |
| `host/broker/context_thread_test.go` | M2 | the read census 5/4 → 6/5 (P3) |
| `host/daemon/daemon.go` | M3a | two `HandleFunc` lines, the Handler doc comment |
| `host/daemon/surface_decide.go`, `host/surfacepage/decide_model.go` | M3a | new: the routes, the grammar |
| `host/daemon/surface_decide_routes_test.go` | M3a | new: RD1 |
| `host/daemon/surface_decide_rig_test.go`, `surface_decide_subject_test.go` | M3b | new: the rig; RD3 |
| `host/daemon/surface_decide_post_test.go` | M3c | new: RD2, RD4–RD8, R4d |
| `host/daemon/surface_decide_discovery_test.go` | M3d | new: RD9, RD9b–RD9d |
| this plan and `sprint_w-world-approvals-by-key-a4-s1.json` | planner | |

---

## 6. Sprint gates (controller, after the last sub-milestone; each measured at the base, §1.1)

In **one** `bash` shell, from the worktree, **not under `nohup`**, nothing else running. Every leg `tee`s its full
output to a log under `$HOME/.ailang/state/world-iter389/gates/`; a verdict is read from its log, never from `tail`
alone. The pre-commit block is the planner's `tools/gates.sh <tree> <logs> final` with `W`/`L` set as below; `base`
mode empties the four `NEW_*` lists.

```
bash <<'EOF2'
set -euo pipefail
export AILANG_BIN=$HOME/.pinned-ailang/ailang
unset WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER WORLD_PASSKEY_RECORD_DIR
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter389-a4-s1; cd "$W"
L=$HOME/.ailang/state/world-iter389/gates; mkdir -p "$L"; MODE=final
s0=$(date +%s)
./scripts/verify_ail.sh 2>&1 | tee "$L/verify_ail.log"
out=$(gofmt -l cmd host); echo "$out" | tee "$L/gofmt.log"; test -z "$out"
go build ./... 2>&1 | tee "$L/build.log"
go vet ./... 2>&1 | tee "$L/vet.log"
go test -run '^$' ./... 2>&1 | tee "$L/fence.log"
if grep -E '^(FAIL|--- FAIL)' "$L/fence.log"; then exit 1; fi
go test -count=1 -timeout 30m ./host/approval/... ./host/coordinator/ ./host/surfacepage/ ./host/daemon/ 2>&1 | tee "$L/pkgs.log"
if grep -E '^(FAIL|--- FAIL)' "$L/pkgs.log"; then exit 1; fi
go test -count=1 -timeout 30m ./host/broker/ 2>&1 | tee "$L/broker.log" || true   # reported, not gating (§6)
if grep -E '^(FAIL|--- FAIL)' "$L/broker.log"; then echo "BROKER_RED (re-run each red alone; known-list reds measure the repo)"; fi
go test -race -count=1 ./host/approval/ ./host/coordinator/ 2>&1 | tee "$L/race.log"
go test -race -count=1 ./host/broker/ -run '^(TestDecideApproval|TestPublishApproval|TestDecisionRefusal|TestReservedApprovalScope|TestApprovalContextIdentity)' 2>&1 | tee -a "$L/race.log"
go test -race -count=1 ./host/daemon/ -run '^(TestDecision|TestDecide|TestEpisodeFinish|TestApprov|TestRequest|TestDaemonCredentialWriterSites|TestSurfaceAppRoutes)' 2>&1 | tee -a "$L/race.log"
if grep -E '^(FAIL|--- FAIL)|WARNING: DATA RACE' "$L/race.log"; then exit 1; fi
NEW_A="TestDecisionMessageLayout TestDecisionMessageRefusesBadInput"
NEW_B="TestDecisionRefusalD1 TestDecisionRefusalD2 TestDecisionRefusalD3 TestDecisionRefusalD4a TestDecisionRefusalD4b TestDecisionRefusalD5 TestDecisionRefusalD6 TestDecisionRefusalD7 TestDecisionRefusalD8 TestDecisionRefusalD9 TestDecisionRefusalD10 TestDecisionRefusalOrder TestDecisionRefusalMustLive TestReservedApprovalScope TestDecideApprovalRefusesASecondDecision TestDecideApprovalFirstDecisionMustLive TestPublishApprovalRefusesAKeyDecision TestPublishApprovalTypedStampMustLive TestDecideApprovalCallerCensus"
NEW_C="TestDiscoverRequestPacketsCompleteness"
NEW_D="TestDecisionRoutesRegistered TestDecisionSubjectJSON TestDecidePOSTGuards TestDecidePOSTLawArms TestDecideTypedLands TestDecideTypedProofBindsTheVerb TestDecidePOSTStaleDigests TestDecideReplayAfterRestart TestDecideRefusesReservedScopesTyped TestDecideRefusesASharedRequest TestDecideDiscoveryFindsAHiddenSecondPacket TestDecideDiscoveryRefusesAnOverflow TestDecideDiscoveryFindsTheOnlyPacketBeyondThePrefix"
if [ "$MODE" = base ]; then NEW_A=""; NEW_B=""; NEW_C=""; NEW_D=""; fi
alt() { local x="$*"; x="${x// /|}"; echo "${x:+|$x}"; }
go test -v -count=1 ./host/approval/... -run "^(TestApprovalImportsOnlyTheStandardLibrary|TestApprovalTestIsTestOnly$(alt $NEW_A))\$" 2>&1 | tee "$L/inv.log"
go test -v -count=1 ./host/broker/ -run "^(TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential|TestQuickstartSection10ProfilesLoad|TestQuickstartProtectedContract|TestApprovalContextIdentity$(alt $NEW_B))\$" 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./host/coordinator/ -run "^(TestValidatePacketBytesCheckNeedsThePacketLabel$(alt $NEW_C))\$" 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./host/daemon/ -run "^(TestDaemonDependencyAllowlist|TestFrozenV1RouteTableMatchesMux|TestSurfaceAppRoutes|TestApprovalRoutesRegistered|TestReadCtxCancelledAfterHandler|TestMCPMountMethodSource|TestDaemonCredentialWriterSites|TestQuickstartDecisionPacketSection|TestDaemonConstructorCallerInventory|TestEpisodeFinishOperatorProofExpiryBoundary$(alt $NEW_D))\$" 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./host/boundary/ -run '^(TestWorldBoundaryDependencyAllowlist|TestBareNetHTTPExemptionIsPerGroup|TestEveryCommandDirectoryIsAProtectedGroup)$' 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./host/store/ -run '^(TestProductionContextRoots|TestProductionGoSurface|TestNoNewDeadlineFreeStoreReads|TestReadGettersHonorContext|TestCredentialsMethodsRequireADeadline)$' 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./host/verifygate/ -run '^(TestNoParallelWriteForkPackagesOutsideVerifygate|TestSubprocessSinkGate)$' 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./cmd/ailang-worldd/ -run '^(TestEveryVerbHasHelp|TestQuickstartAttendedStepsBoxMatchesSection9|TestQuickstartSection9FlagsMatchTheCLI|TestQuickstartSection10FlagsMatchTheCLI|TestQuickstartSection11FlagsMatchTheCLI|TestQuickstartServiceSectionMatchesTheCLIHelp)$' 2>&1 | tee -a "$L/inv.log"
for t in TestApprovalImportsOnlyTheStandardLibrary TestApprovalTestIsTestOnly TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential TestQuickstartSection10ProfilesLoad TestQuickstartProtectedContract TestApprovalContextIdentity \
  TestValidatePacketBytesCheckNeedsThePacketLabel TestDaemonDependencyAllowlist TestFrozenV1RouteTableMatchesMux TestSurfaceAppRoutes TestApprovalRoutesRegistered TestReadCtxCancelledAfterHandler TestMCPMountMethodSource \
  TestDaemonCredentialWriterSites TestQuickstartDecisionPacketSection TestDaemonConstructorCallerInventory TestEpisodeFinishOperatorProofExpiryBoundary \
  TestWorldBoundaryDependencyAllowlist TestBareNetHTTPExemptionIsPerGroup TestEveryCommandDirectoryIsAProtectedGroup \
  TestProductionContextRoots TestProductionGoSurface TestNoNewDeadlineFreeStoreReads TestReadGettersHonorContext TestCredentialsMethodsRequireADeadline \
  TestNoParallelWriteForkPackagesOutsideVerifygate TestSubprocessSinkGate TestEveryVerbHasHelp TestQuickstartAttendedStepsBoxMatchesSection9 \
  TestQuickstartSection9FlagsMatchTheCLI TestQuickstartSection10FlagsMatchTheCLI TestQuickstartSection11FlagsMatchTheCLI TestQuickstartServiceSectionMatchesTheCLIHelp \
  $NEW_A $NEW_B $NEW_C $NEW_D; do
  grep -qE -- "^--- PASS: $t " "$L/inv.log" || { echo "$t did not PASS"; exit 1; }; done
if grep -E '^\s*--- (SKIP|FAIL)' "$L/inv.log"; then exit 1; fi
echo "named PASS: $(grep -cE '^--- PASS' "$L/inv.log")" | tee -a "$L/inv.log"
echo "PRE_COMMIT_GATES_GREEN $(( $(date +%s) - s0 ))s" | tee "$L/pre_commit.ok"
EOF2
# AFTER the controller commits (the tracked-binary gate reads committed blobs; this leg runs go test ./...):
bash -c 'set -o pipefail; export AILANG_BIN=$HOME/.pinned-ailang/ailang; unset WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER; cd /Users/voightkampff/dev/sunholo-data/.wt-world-iter389-a4-s1 && VERIFY_GO_LEGS=plain AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_go.sh 2>&1 | tee $HOME/.ailang/state/world-iter389/gates/verify_go.log; echo "verify_go rc=$?" | tee -a $HOME/.ailang/state/world-iter389/gates/verify_go.log'
git -C /Users/voightkampff/dev/sunholo-data/.wt-world-iter389-a4-s1 status --short   # only the two plan files, until the controller commits them
```

**The `host/broker` whole-package leg** (`broker.log`) is reported, not gating, inside the script: under load it
carries the directive's known exec-probe/timing reds and its default 10-minute timeout was hit once in the
planner's M2 run (`panic: test timed out after 10m0s` in `TestExecProbePassThroughShimFailsEveryConfinementArm`,
`logs/m2_pkgs.log`); the script gives it 30 minutes. **The controller reads `broker.log`**: any red outside the
known list below, or any red among `TestDecide*`, `TestDecision*`, `TestReserved*`, `TestPublishApproval*`,
`TestApprovalContextIdentity`, `TestApprove*`, is a regression; a known-list red is re-run alone once.

**At the base** (§1.1): `verify_ail.sh`, gofmt, build, vet and the fence green; the packages red only on load/environment-class tests, each green alone except `TestServeStartsWhenTheProbeSpendsTheStartupClock` (the rig's `node --version` > 5 s, at the base too); `-race` and **33** inventories green in 265 s. **On the final tree** (round-trip tree, byte-identical to the mutant tree): `verify_ail.sh`, gofmt, build, vet and the fence green; packages `ok` except one known-list red (`TestMCPPostBudgetProductionConstants/batch`); `-race` ok ×4 legs; **68** named `--- PASS` (33 carried + 35 new) in 312 s (§1.2 C20–C23). Because `set -o pipefail` stops the script at a red package leg, the planner ran the `-race` and inventory legs as `tools/gates_tail.sh` (the same lines) on both trees; the controller runs the whole block once and, on a known-list red, re-runs that test alone and then the remaining legs. The after-commit `verify_go.sh` was not measured by the planner (§1.1).

**Known load- or environment-sensitive (re-run once, alone; do not chase)**:
`TestSeedRestartIgnoresASurvivingOldCopy` (host/broker), `TestHostsLaunchersExecHost`, `host/broker`'s exec-probe
timeouts (incl. `TestExecProtectedCheckTimeout/reserve`, and, measured by this planner under load,
`TestExecProbeNet6*`, `TestPipeCloseBoundsEscapedDescendant`, `TestFailedKillBoundsCallAndChildIsReapedLater`,
`TestReservationIsExactUnderConcurrency`), `TestPublishApprovalRefusalSetWithALandedPositiveControl` (host/broker),
local `-race` of host/broker, host/daemon, host/verifygate (incl. `TestSeToolsGateArmsAreLoadBearing/plan-bytes`),
`TestMCPPostBudgetProductionConstants/batch` (host/daemon), `TestBackoff` (host/surfacefollow),
`TestDemoSeedPreflightRefusals`, `TestWorktreeSnapshotRefusesAConfiguredFilter`,
`TestMissionConfigDefaultFlowIntegration/invalid_ledger` (and `/valid_ledger`), `scripts/marketingcheck`
supervisor tests, `TestMarketingCaptureFakeChromeSixShots/inventory`, `TestWorldContextAndStatusline`, and
`host/coordinator`'s `TestCapsuleWarmEqualsColdForSeTools` and `host/daemon`'s `TestMarketingCaptureInventoryMutation`, `TestRegistryPhaseHasItsOwnClock` and `TestServeStartsWhenTheProbeSpendsTheStartupClock` (measured red at the base under load, §1.1; the last red alone too: `node --version` > 5 s), and `cmd/world-publish`'s interpreter-timeout reds `TestSeToolsManifestPublishes` and
`TestWorldCoreManifestMatchesTheCommittedGolden` (measured by this planner under load: "interpreter check did not
finish within 10s", "context deadline exceeded"; `cmd/world-publish` is not touched).

**Must stay green, unedited**: every `TestApprove*`, `TestPasskey*`, `TestApprovePage*`, `TestSurfaceApprove*`,
`registry_publish_test.go`, `episode_finish*_test.go`, `decisionpacket*_test.go`, `TestFrozenV1RouteTableMatchesMux`
(13), `TestDaemonDependencyAllowlist`, `TestSurfaceAppRoutes`, `TestApprovalRoutesRegistered`,
`TestDaemonCredentialWriterSites`. The only edit to an existing test is `context_thread_test.go`'s census line (P3).

---

## 7. Executor directive notes

- **No git writes.** Never `git add`, `commit`, `stash`, `checkout`, `reset`, `restore`, `rebase`, `push`,
  `worktree` or `apply`. `patch -p1` (§2) is the only writer; read-only `git status`/`git diff` are fine.
- **Never `nohup`, never background a test leg.** Every acceptance command runs in the foreground of one bash
  shell, `tee`'d to `$HOME/.ailang/state/world-iter389/exec-logs/`. Poll nothing in the background: there is
  nothing to background.
- **Runs.** Run 1 = M1a, M1b, M1c, M2; Run 2 = M3a, M3b, M3c, M3d. Each run reads §0, §2 and its own §3 sections,
  then applies its patches in order.
- **Expected numstat after each sub-milestone, measured AFTER the controller's commit of it** (an uncommitted
  tree hides the new untracked files from numstat):
  ```
  bash -c 'cd /Users/voightkampff/dev/sunholo-data/.wt-world-iter389-a4-s1 && git diff --numstat 9a6b959deed573d7966863a1f166f3a8c5700f2e HEAD -- cmd host | awk "{s+=\$1+\$2} END {print s}"'
  ```
  After the controller's commit of M1a **162**, M1b **508**, M1c **733**, M2 **1004**, M3a **1356**, M3b
  **1628**, M3c **1901**, M3d **1985**. **Any other value means a patch was not applied as written: stop and
  report.**
- **Snapshot after every sub-milestone, cumulatively**:
  ```
  bash -c 'set -euo pipefail; M=M1a; W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter389-a4-s1
    D=$HOME/.ailang/state/world-iter389/snap/$M; rm -rf "$D"; mkdir -p "$D"; cd "$W"
    git status --porcelain --untracked-files=all | while read -r st f; do
      case "$f" in design_docs/planned/w-world-approvals-by-key-a4-sprint1-plan.md|design_docs/planned/sprint_w-world-approvals-by-key-a4-s1.json) continue;; esac
      mkdir -p "$D/$(dirname "$f")"; cp -p "$f" "$D/$f"; done
    (cd "$D" && find . -type f ! -name MANIFEST.sha256 | sort | xargs shasum -a 256) > "$D/MANIFEST.sha256"
    cat "$D/MANIFEST.sha256"'
  ```
  (`M` = `M1a` … `M3d`.) No file is deleted, so the porcelain status has only `M` and `??`; stop on any other
  status. The controller rebuilds one commit per sub-milestone (`commit_ms.sh` with `WT` and the snapshot root
  pointed here; **commit the plan BEFORE running it — it deletes untracked files**).
- **Stop conditions.** A `patch --dry-run` failure; a numstat other than the expected value; a red
  earlier-milestone test; a `--- FAIL`, or a `--- SKIP` in a §3 `-v` run; any need to touch a file outside §5; any
  reference to `~/.ssh`, `~/.ailang/world`, a real `/dev/tty`, `ssh-keygen` or a real store. Stop and report; do
  not work around. A red in an untouched package that is on §6's known list is re-run alone once.
- **Timing (under the planner's load).** Each `go test` of one package compiles in ~20–60 s; M1a–M3c acceptance
  ~1–2 min each; M3d's whole `host/daemon` run ~5–10 min (RD9c stores 2,000 packets in ~2 s).
- **Shell.** Wrap loops, globs and `$VAR` expansions in `bash -c` or a `bash <<'EOF'` heredoc (zsh does not
  word-split). Export `AILANG_BIN` in the same shell as the Go tests. `tee` every gate's output to a file; never
  pipe a gate's only copy through `tail`. macOS has no `timeout` command.

---

## 8. Risks

- **R1 — The typed decide path has no shipped client until sprint 2's head** (S1-D1): the route is proven over
  HTTP with test proofs; an operator would have to hand-write a proof file to use it. Sprint 2 must land CD1–CD4
  first.
- **R2 — RK3b (design): more than 2,000 stored packets refuse every A4 decision** (D4a). Measured cost: RD9c's
  2,001-packet discovery inside one POST ran within the default deadline under load (§1.2 results).
- **R3 — The identical-decision idempotence (S1-D10)** is a behaviour the design did not name: it answers the
  same ref for a byte-identical re-decision and writes nothing (previously it appended a duplicate head). No
  reader of the chain distinguishes the two (the fold keeps the newest decision per request, which is the same
  object).
- **R4 — `humanMu` is not exercised concurrently by a sprint-1 test** (named survivor `H-LOCK`): the lock is the
  finish route's and `serializedHandler`'s (row 172/236), and the design relies on it for D6-then-write
  atomicity; a concurrency test belongs with sprint 2's key arms (the writer adds a second resource).
- **R5 — `host/broker`'s whole package is slow and load-sensitive** (§6); the sprint's broker tests are fast and
  named.

---

## 9. Acceptance criteria (the judge scores against these)

1. **The message** (M1a): `DecisionMessage` is exactly 91 bytes `"world-decision/v1" ‖ 0x00 ‖ verb ‖ P ‖ Q ‖
   u64be(expiry)`, golden for two subjects; verb ∉ {1, 2} and a negative expiry are refused; namespace
   `world-decision@sunholo`; the mint message and namespace are unchanged; `host/approval` stays stdlib-only.
2. **The law** (M1b): `DecisionRefusal` refuses in order D1 400, D2 404, D3 404, D4a 409, D4b 409, D5 403, D6
   409, D7 409, D8 409, D9 400, D10 409 with the design's reasons; D4a is decided from `Discovery` alone (never
   `Carrying`); the law is pure (no I/O, no clock).
3. **L1** (M1b): publish-shaped (parse, `registry:`, `effect=Registry.Publish`) and finish-shaped (parse,
   `github-pr:`) scopes are reserved for every format, typed included.
4. **The discovery** (M1c): complete only when a page comes back short before 2,000 packets; the 2,001st is the
   sentinel; the cursor follows `after`; undecodable and missing payloads are counted; store errors are returned.
5. **The broker** (M2): a second, different decision of a request is refused with nothing written; a
   byte-identical one answers its ref with nothing written; a key `decidedBy` never stamps a publish (L2); every
   `registry_publish_test.go` and `episode_finish*_test.go` arm passes unedited; the non-test callers of the
   decide entry point are exactly `DecideApproval`'s wrapper, `RecordOperatorApproval`, `mintAttendedApproval` and
   `handleSurfaceDecide` (XG1).
6. **The routes** (M3a–M3d): `GET /surface/decisions/{ref}` and `POST /surface/decide/{ref}` exist under `/surface`
   behind `approvalRoute`; `/v1` stays 13; GET sends bytes and never a message or hash; the POST refuses R0
   (`Authorization`/`Cookie`), R2–R6, and every law step end to end with 0 decisions, an unmoved head and an
   unconsumed proof; a valid typed proof decides once, with `decidedBy = operator:terminal`, and is unlinked; the
   proof's phase binds the verb and its name binds the request; a replay after a daemon restart is 409 D6; a
   hidden second packet past position 101 is found (409), an overflow is refused (409 D4a), a lone packet past
   the prefix decides.
7. **No new `/v1` route, no `.ail`, no CSP change, no CORS header, no new dependency, no new
   `store.OpenCredentials(` site, no context root, no `t.Parallel`**; `decidedBy` is never read from a caller.
8. Each §4 mutant, applied alone, makes its named test **and** its package suite fail (inventories excluded); the
   named equivalents and survivors carry their reasons; the ledger has a line per run.
9. Every §6 gate passes, including every inventory by name, the `-race` legs and `VERIFY_GO_LEGS=plain
   ./scripts/verify_go.sh` after the commit (modulo §6's re-run rule); `git status --short` lists only §5.
10. The numstat after each sub-milestone's commit equals §7's value; the sprint total is ≤ 2,000 (1,985).
11. **The fallback cut is recorded**: CD1–CD4 (≈ 400 by the design's price) are sprint 2's head (Appendix C).
12. **What an operator can do after this sprint**: nothing new at a terminal yet — the typed decide path exists
    and is proven over HTTP, but no shipped command writes its `decision-<hex>` proof (S1-D1). What the agent can
    do: nothing new (no grant reaches `POST /surface/decide`: R0 refuses its session, it cannot write the proof
    directory, XG1 pins the callers). **The attended exit run (design §10) still waits for**: sprint 2 (the typed
    CLI CD1–CD4, then the SSH and passkey arms), sprint 3 (the card, *Needs you*, QUICKSTART §18).
13. No real key, store, `/dev/tty`, `~/.ssh` access, mint into a real store, publish or `ssh-keygen` anywhere in
    the sprint (D-WORLD-109).

---

## Appendix C — after this sprint (sprint 2's head and handoff)

- **Moved to sprint 2's head (S1-D1), ≈ 400 lines by the design's price**: CD1 `TestApproveDecideTyped`, CD2
  `TestApproveDecideNoTerminal`, CD3 `TestApproveDecideFlagRules`, CD4 `TestApproveDecideRefusesAMismatchedSubject`
  and `cmd/ailang-worldd/approve_decide.go`: the `sha256:` + 64 hex id shape in `runApprove` (beside the `mr-`
  check, `approve.go:197-201`), `--deny` (refused with an `mr-` id), `--await` refused with a request ref, GET
  `/surface/decisions/{ref}` from the daemon address (how the CLI finds the daemon for a request ref is sprint
  2's to decide: the mint requests carry `Daemon.Addr`, a decision request does not), the client-side hash
  checks of both byte fields and `packet.requestRef`, every field printed with the verb in capitals and the
  deadline, `y` at `/dev/tty` (`openTerminal`), `writeOperatorProof(credentials, "decision-"+hex(Q),
  "decide-"+verb)`, the POST with `Origin: http://127.0.0.1:<port>` and `typed: true`; help lines
  (`TestEveryVerbHasHelp`), `TestProductionContextRoots` (+1 root). Design Q2 (`approve --list` decisions) goes
  with it. With the CLI, sprint 2 is the design's ≈ 1,800 + ≈ 400 = ≈ 2,200 — **over the cap: sprint 2's planner
  must price its own fallback (the design's: the webauthn arm to sprint 3)**.
- **N4 — K0's slot (S1-D9), for sprint 2's planner verbatim**: law D1–D10 (unchanged; D6 refuses a landed retry)
  → arm slot: judge (read-only; W1–W14, P2–P5, `Binds`; W15 deferred) → writer open (non-waiting; 503 `Busy`) →
  `RecordKeyDecision` {K1, K2, **K0** (reload approver via `useApprover` on both paths + passkey row; exact match
  of an unlanded row's artifact and signed fields skips only W15/`usePasskey`), else W15 + `usePasskey`, K5} →
  `broker.DecideApproval` → `LandKeyDecision` → close. Sprint 2 edits sprint 1 only at: `ParseDecideBody` (exactly
  one of `sshsig`, `webauthn`, `typed`), the arm slot (a `switch` on the format), a post-`DecideApproval` land
  call, `DecisionSubject.Passkeys` (P7), and RD9b/RD9c's SSHSIG arms + "0 `key_decisions` rows" (design §6).
- **Carried for sprint 2**: R4e (L1 compiled out → L2 refuses the key decision, S1-D15); a concurrent-decide test
  for `humanMu` (R4, named survivor `H-LOCK`); `TestDaemonCredentialWriterSites` 1 → 2.
- **For the attended exit run** (design §10, Mark; NOT run by the loop, D-WORLD-109): unchanged; its typed step 2
  needs sprint 2's CD1–CD4.
- **D-WORLD-115** stays open (default A applied); sprint 1 does not depend on it (S1-D2).

---

## Appendix P — the patches (apply with §2; generated from the compiled scratch stages; do not edit)

### Patch M1a — the decision message

````diff M1a
diff --git a/host/approval/decision.go b/host/approval/decision.go
new file mode 100644
index 0000000..c6635a9
--- /dev/null
+++ b/host/approval/decision.go
@@ -0,0 +1,69 @@
+package approval
+
+// Row 201 A4 (design_docs/planned/w-world-approvals-by-key-a4.md §2.3): the message a decision of
+// a Human.Approve request signs. It is not the mint message (SignedMessage): a decision is approve
+// or deny, so the verb is signed, and a distinct prefix and SSHSIG namespace keep a mint signature
+// from deciding and a decision signature from minting. Pure: no I/O, no clock.
+
+import (
+	"crypto/sha256"
+	"encoding/binary"
+	"errors"
+)
+
+// DecisionPrefix starts every signed decision message (design §2.3).
+const DecisionPrefix = "world-decision/v1"
+
+// DecisionNamespace is the SSHSIG namespace every decision is signed under.
+const DecisionNamespace = "world-decision@sunholo"
+
+// DecisionMessageLen is the length of DecisionMessage: 17 + 1 + 1 + 32 + 32 + 8.
+const DecisionMessageLen = len(DecisionPrefix) + 1 + 1 + 32 + 32 + 8
+
+// The signed verb bytes; anything else is refused.
+const (
+	VerbApprove byte = 0x01
+	VerbDeny    byte = 0x02
+)
+
+// ErrDecisionVerb is a verb byte other than VerbApprove and VerbDeny.
+var ErrDecisionVerb = errors.New("approval: a decision verb is approve (0x01) or deny (0x02)")
+
+// DecisionVerb is the signed byte of a decision word ("approve" or "deny").
+func DecisionVerb(word string) (byte, error) {
+	switch word {
+	case "approve":
+		return VerbApprove, nil
+	case "deny":
+		return VerbDeny, nil
+	}
+	return 0, ErrDecisionVerb
+}
+
+// DecisionMessage is the 91-byte message a decision signs (design §2.3):
+// DecisionPrefix ‖ 0x00 ‖ verb ‖ packet digest P ‖ request digest Q ‖ expiry (uint64, big-endian).
+func DecisionMessage(verb byte, packet, request [32]byte, expiry int64) ([]byte, error) {
+	if verb != VerbApprove && verb != VerbDeny {
+		return nil, ErrDecisionVerb
+	}
+	if expiry < 0 {
+		return nil, ErrNegativeExpiry
+	}
+	m := make([]byte, 0, DecisionMessageLen)
+	m = append(m, DecisionPrefix...)
+	m = append(m, 0)
+	m = append(m, verb)
+	m = append(m, packet[:]...)
+	m = append(m, request[:]...)
+	m = binary.BigEndian.AppendUint64(m, uint64(expiry))
+	return m, nil
+}
+
+// DecisionMessageHash is the SHA-256 of DecisionMessage: a passkey's WebAuthn challenge.
+func DecisionMessageHash(verb byte, packet, request [32]byte, expiry int64) ([32]byte, error) {
+	m, err := DecisionMessage(verb, packet, request, expiry)
+	if err != nil {
+		return [32]byte{}, err
+	}
+	return sha256.Sum256(m), nil
+}
diff --git a/host/approval/decision_test.go b/host/approval/decision_test.go
new file mode 100644
index 0000000..bd6cd34
--- /dev/null
+++ b/host/approval/decision_test.go
@@ -0,0 +1,93 @@
+package approval
+
+// Row 201 A4 sprint 1 (design §2.3, §6 DM1, DM2): the decision message's layout and refusals.
+
+import (
+	"bytes"
+	"crypto/sha256"
+	"encoding/binary"
+	"encoding/hex"
+	"errors"
+	"testing"
+)
+
+// DM1: 91 bytes; the prefix, 0x00, the verb at 18, P at 19, Q at 51, the big-endian expiry at 83;
+// golden hex for two subjects (computed by an independent encoder, not by DecisionMessage).
+func TestDecisionMessageLayout(t *testing.T) {
+	p1, q1 := sha256.Sum256([]byte("packet-1")), sha256.Sum256([]byte("request-1"))
+	var p2, q2 [32]byte
+	for i := range p2 {
+		p2[i], q2[i] = byte(i), byte(0xff-i)
+	}
+	for _, tc := range []struct {
+		verb         byte
+		p, q         [32]byte
+		expiry       int64
+		golden, hash string
+	}{
+		{VerbApprove, p1, q1, 1794414540,
+			"776f726c642d6465636973696f6e2f76310001c48e8d2fd69d694f90962ad5c6da898c6e83f070c32804601d4c8b4b0285412719f1064b619d49d35392eac7261cd7266c720671fc594f4b226f32bf0bee74ba000000006af497cc",
+			"efa6cc1b6c482505379b440b52639ae3db80351938d4e3addf1119553b28c4b7"},
+		{VerbDeny, p2, q2, 0x0102030405060708,
+			"776f726c642d6465636973696f6e2f76310002000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1ffffefdfcfbfaf9f8f7f6f5f4f3f2f1f0efeeedecebeae9e8e7e6e5e4e3e2e1e00102030405060708",
+			"080bf79e5fa2f6a74c8726f47460ad398f2794235f5e92acbe1e51b35cc35ec2"},
+	} {
+		m, err := DecisionMessage(tc.verb, tc.p, tc.q, tc.expiry)
+		if err != nil {
+			t.Fatal(err)
+		}
+		if len(m) != 91 || DecisionMessageLen != 91 {
+			t.Fatalf("len = %d (DecisionMessageLen %d), want 91", len(m), DecisionMessageLen)
+		}
+		if got := hex.EncodeToString(m); got != tc.golden {
+			t.Errorf("message = %s, want %s", got, tc.golden)
+		}
+		if string(m[:17]) != "world-decision/v1" || m[17] != 0 || m[18] != tc.verb ||
+			!bytes.Equal(m[19:51], tc.p[:]) || !bytes.Equal(m[51:83], tc.q[:]) || binary.BigEndian.Uint64(m[83:]) != uint64(tc.expiry) {
+			t.Errorf("layout of %x", m)
+		}
+		h, err := DecisionMessageHash(tc.verb, tc.p, tc.q, tc.expiry)
+		if err != nil || hex.EncodeToString(h[:]) != tc.hash {
+			t.Errorf("hash = %x, %v; want %s", h, err, tc.hash)
+		}
+	}
+	// Domain separation from the mint message (§2.3): a different prefix and length.
+	mint, err := SignedMessage(p1, q1, 1794414540)
+	if err != nil {
+		t.Fatal(err)
+	}
+	m, _ := DecisionMessage(VerbApprove, p1, q1, 1794414540)
+	if len(mint) == len(m) || bytes.HasPrefix(m, mint[:18]) || DecisionNamespace == Namespace || DecisionNamespace != "world-decision@sunholo" {
+		t.Error("the decision message is not separated from the mint message")
+	}
+	if v, err := DecisionVerb("approve"); v != VerbApprove || err != nil {
+		t.Errorf("approve = %d, %v", v, err)
+	}
+	if v, err := DecisionVerb("deny"); v != VerbDeny || err != nil {
+		t.Errorf("deny = %d, %v", v, err)
+	}
+}
+
+// DM2: verb 0 and 3 and a negative expiry are refused, by the message and by its hash.
+func TestDecisionMessageRefusesBadInput(t *testing.T) {
+	var p, q [32]byte
+	for _, verb := range []byte{0, 3, 0xff} {
+		if m, err := DecisionMessage(verb, p, q, 1); m != nil || !errors.Is(err, ErrDecisionVerb) {
+			t.Errorf("verb %d: %x, %v", verb, m, err)
+		}
+		if _, err := DecisionMessageHash(verb, p, q, 1); !errors.Is(err, ErrDecisionVerb) {
+			t.Errorf("hash verb %d: %v", verb, err)
+		}
+	}
+	if m, err := DecisionMessage(VerbDeny, p, q, -1); m != nil || !errors.Is(err, ErrNegativeExpiry) {
+		t.Errorf("negative expiry: %x, %v", m, err)
+	}
+	if m, err := DecisionMessage(VerbApprove, p, q, 0); err != nil || len(m) != 91 {
+		t.Errorf("expiry 0 (control): %v", err)
+	}
+	for _, word := range []string{"", "Approve", "approve ", "yes"} {
+		if _, err := DecisionVerb(word); !errors.Is(err, ErrDecisionVerb) {
+			t.Errorf("word %q: %v", word, err)
+		}
+	}
+}
````

### Patch M1b — the decision law and the reserved scope

````diff M1b
diff --git a/host/broker/decision_law.go b/host/broker/decision_law.go
new file mode 100644
index 0000000..c9d1f04
--- /dev/null
+++ b/host/broker/decision_law.go
@@ -0,0 +1,132 @@
+package broker
+
+// Row 201 A4 (design_docs/planned/w-world-approvals-by-key-a4.md §2.4–§2.6): the law an operator's
+// decision of a Human.Approve request passes before anything is written. Pure: the daemon reads the
+// facts (the store, the clock, the request body) and this file only orders them; the first failing
+// fact wins (D1–D10), and the format arm (D11) runs after it in the daemon.
+
+import (
+	"fmt"
+	"strings"
+	"time"
+)
+
+// MaxDecisionWindow is how far ahead of now a decision's expiry may be (D10), in seconds.
+const MaxDecisionWindow = 600
+
+// Discovery is the result of the request-wide decision-packet discovery (design §2.4 [r1]).
+// Complete is true only when the keyset scan ran out before its overflow sentinel; Carrying counts
+// every decoded packet naming the request; Undecodable counts packets that did not decode.
+type Discovery struct {
+	Complete                       bool
+	Scanned, Carrying, Undecodable int
+}
+
+// DecisionFacts are what the daemon read, under humanMu, for one decision of request Ref.
+type DecisionFacts struct {
+	Ref            string    // the request ref as given (the reasons name it)
+	RefOK          bool      // D1: sha256: and 64 lowercase hex
+	RequestOK      bool      // D2: a strictly decoding world/approval-request/v1 of Human.Approve
+	OnChain        bool      // D3: on the approvals chain
+	Discovery      Discovery // D4a
+	PacketCode     int       // D4b: the one carrying packet's ValidatePacket code
+	PacketSentence string    // D4b: that code's sentence when it is not 0
+	Scope          string    // D5: the request's scope
+	Decided        bool      // D6: the chain holds a decision of the request
+	DecidedVerb    string    // D6
+	DecidedBy      string    // D6
+	Now            int64     // D7, D10: the daemon's clock, Unix seconds
+	Deadline       int64     // D7, D10: the packet's deadlineAt
+	Packet         [32]byte  // D8: the stored packet object digest P
+	Request        [32]byte  // D8: the stored request object digest Q
+	PostedPacket   [32]byte  // D8: the P the caller saw
+	PostedRequest  [32]byte  // D8: the Q the caller saw
+	Verb           string    // D9
+	Expiry         int64     // D10: the decision's expiry
+}
+
+// LawRefusal is the law's verdict: Step "" (and Status 0) accepts; otherwise the first failing
+// step, its HTTP status and a reason that names the fix.
+type LawRefusal struct {
+	Step   string
+	Status int
+	Reason string
+}
+
+// DecisionWindowMax is the latest expiry a decision made at now may carry: min(now + 600, deadline).
+func DecisionWindowMax(now, deadline int64) int64 {
+	return min(now+MaxDecisionWindow, deadline)
+}
+
+// SubjectRefusal is the law's D1–D7: what refuses a request before any decision of it is posted.
+// GET /surface/decisions/{ref} shows it; DecisionRefusal starts with it.
+func SubjectRefusal(f DecisionFacts) LawRefusal {
+	switch {
+	case !f.RefOK:
+		return LawRefusal{"D1", 400, fmt.Sprintf("%q is not a request ref (sha256: and 64 lowercase hex)", f.Ref)}
+	case !f.RequestOK:
+		return LawRefusal{"D2", 404, "no such approval request " + f.Ref}
+	case !f.OnChain:
+		return LawRefusal{"D3", 404, "request " + f.Ref + " is not an approval request of this world (it is not on the approvals chain)"}
+	case !f.Discovery.Complete || f.Discovery.Undecodable > 0:
+		return LawRefusal{"D4a", 409, fmt.Sprintf("the decision packets carrying request %s cannot be counted completely (more than %d stored packets, or one that does not decode); it cannot be decided here", f.Ref, DecisionDiscoveryLimit)}
+	case f.Discovery.Carrying == 0:
+		return LawRefusal{"D4b", 409, "no committed decision packet carries request " + f.Ref + "; it cannot be decided here"}
+	case f.Discovery.Carrying != 1:
+		return LawRefusal{"D4b", 409, fmt.Sprintf("request %s is carried by %d decision packets; deciding it would decide all of them; it cannot be decided here", f.Ref, f.Discovery.Carrying)}
+	case f.PacketCode != 0:
+		return LawRefusal{"D4b", 409, "the decision packet carrying request " + f.Ref + " is refused: " + f.PacketSentence}
+	}
+	if owner, reserved := ReservedApprovalScope(f.Scope); reserved {
+		return LawRefusal{"D5", 403, reservedReason(owner)}
+	}
+	switch {
+	case f.Decided:
+		return LawRefusal{"D6", 409, fmt.Sprintf("request %s was already decided %s by %s", f.Ref, f.DecidedVerb, f.DecidedBy)}
+	case f.Now >= f.Deadline:
+		return LawRefusal{"D7", 409, "the packet's deadline passed at " + time.Unix(f.Deadline, 0).UTC().Format(time.RFC3339) + "; its policy is Cancel"}
+	}
+	return LawRefusal{}
+}
+
+// DecisionRefusal is the whole law, D1–D10, over a posted decision.
+func DecisionRefusal(f DecisionFacts) LawRefusal {
+	if r := SubjectRefusal(f); r.Step != "" {
+		return r
+	}
+	switch {
+	case f.PostedPacket != f.Packet || f.PostedRequest != f.Request:
+		return LawRefusal{"D8", 409, "the request changed since it was shown; reload it"}
+	case f.Verb != "approve" && f.Verb != "deny":
+		return LawRefusal{"D9", 400, fmt.Sprintf("bad verb %q: a decision is approve or deny", f.Verb)}
+	case f.Expiry <= f.Now || f.Expiry > DecisionWindowMax(f.Now, f.Deadline):
+		return LawRefusal{"D10", 409, "the decision's expiry is stale or too far ahead (at most 10 minutes, and not after the packet's deadline)"}
+	}
+	return LawRefusal{}
+}
+
+// DecisionDiscoveryLimit is how many stored decision packets the discovery reads before its
+// overflow sentinel (design §2.4 [r1]); one more means the set is not established.
+const DecisionDiscoveryLimit = 2000
+
+// ReservedApprovalScope reports whether scope belongs to a typed owner that is not A4's card
+// (design §2.6 L1): publish-shaped scopes (world-publish approve, R4) and finish-shaped scopes
+// (episode finish). Near-miss shapes fail closed: a future parser change cannot make a refused
+// shape decidable here.
+func ReservedApprovalScope(scope string) (owner string, reserved bool) {
+	if _, err := ParsePublishApprovalScope(scope); err == nil ||
+		strings.HasPrefix(scope, "registry:") || strings.Contains(scope, "effect="+EffectRegistryPublish) {
+		return "publish", true
+	}
+	if _, err := ParseFinishApprovalScope(scope); err == nil || strings.HasPrefix(scope, "github-pr:") {
+		return "finish", true
+	}
+	return "", false
+}
+
+func reservedReason(owner string) string {
+	if owner == "publish" {
+		return "publish approvals are minted by `world-publish approve` at a terminal (R4); this request cannot be decided here"
+	}
+	return "finish approvals are made by `ailang-worldd episode finish`; this request cannot be decided here"
+}
diff --git a/host/broker/decision_law_test.go b/host/broker/decision_law_test.go
new file mode 100644
index 0000000..fb53b00
--- /dev/null
+++ b/host/broker/decision_law_test.go
@@ -0,0 +1,214 @@
+package broker
+
+// Row 201 A4 sprint 1 (design §2.5, §2.6, §6 DL1–DL12, DL4a, DL4b, R4a): the decision law, one
+// test per refusal (the unit is the branch), its order, its must-live control, and L1.
+
+import (
+	"crypto/sha256"
+	"strings"
+	"testing"
+)
+
+const lawRef = "sha256:00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff"
+
+// lawFacts is a decidable approve: every fact holds.
+func lawFacts() DecisionFacts {
+	p, q := sha256.Sum256([]byte("p")), sha256.Sum256([]byte("q"))
+	return DecisionFacts{Ref: lawRef, RefOK: true, RequestOK: true, OnChain: true,
+		Discovery: Discovery{Complete: true, Scanned: 3, Carrying: 1}, Scope: "worktree",
+		Now: 1000, Deadline: 5000, Packet: p, Request: q, PostedPacket: p, PostedRequest: q, Verb: "approve", Expiry: 1600}
+}
+
+// wantLaw fails unless the law refuses f at step with status and a reason containing each of want.
+func wantLaw(t *testing.T, what string, f DecisionFacts, step string, status int, want ...string) {
+	t.Helper()
+	r := DecisionRefusal(f)
+	if r.Step != step || r.Status != status {
+		t.Fatalf("%s: %+v, want %s %d", what, r, step, status)
+	}
+	for _, w := range want {
+		if !strings.Contains(r.Reason, w) {
+			t.Errorf("%s: reason %q lacks %q", what, r.Reason, w)
+		}
+	}
+}
+
+func TestDecisionRefusalD1(t *testing.T) {
+	f := lawFacts()
+	f.RefOK, f.Ref = false, "sha256:XYZ"
+	wantLaw(t, "D1", f, "D1", 400, `"sha256:XYZ" is not a request ref`)
+}
+
+func TestDecisionRefusalD2(t *testing.T) {
+	f := lawFacts()
+	f.RequestOK = false
+	wantLaw(t, "D2", f, "D2", 404, "no such approval request "+lawRef)
+}
+
+func TestDecisionRefusalD3(t *testing.T) {
+	f := lawFacts()
+	f.OnChain = false
+	wantLaw(t, "D3", f, "D3", 404, "not on the approvals chain")
+}
+
+// DL4a: a scan that hit its sentinel, or saw an undecodable packet, is refused even when it saw
+// exactly one carrying packet (one visible match is never exactly one match).
+func TestDecisionRefusalD4a(t *testing.T) {
+	f := lawFacts()
+	f.Discovery = Discovery{Complete: false, Scanned: 2000, Carrying: 1}
+	wantLaw(t, "incomplete", f, "D4a", 409, "cannot be counted completely", "more than 2000")
+	f.Discovery = Discovery{Complete: true, Scanned: 4, Carrying: 1, Undecodable: 1}
+	wantLaw(t, "undecodable", f, "D4a", 409, "cannot be counted completely")
+	f.Discovery = Discovery{Complete: false, Scanned: 2000, Carrying: 2}
+	wantLaw(t, "incomplete before the count", f, "D4a", 409)
+}
+
+// DL4b: a complete count of 0 or 2, or one carrying packet that does not validate.
+func TestDecisionRefusalD4b(t *testing.T) {
+	f := lawFacts()
+	f.Discovery.Carrying = 0
+	wantLaw(t, "zero", f, "D4b", 409, "no committed decision packet carries request "+lawRef)
+	f.Discovery.Carrying = 2
+	wantLaw(t, "two", f, "D4b", 409, "is carried by 2 decision packets")
+	f = lawFacts()
+	f.PacketCode, f.PacketSentence = 5, "its approval request is not among the newest 40 approval heads, so its approvals-chain membership is not established"
+	wantLaw(t, "code 5", f, "D4b", 409, "is refused: ", "not established")
+}
+
+func TestDecisionRefusalD5(t *testing.T) {
+	f := lawFacts()
+	f.Scope = "registry:https://r.example/package:v/n/version:1.0.0"
+	wantLaw(t, "publish", f, "D5", 403, "world-publish approve", "(R4)")
+	f.Scope = "github-pr:o/r"
+	wantLaw(t, "finish", f, "D5", 403, "ailang-worldd episode finish")
+}
+
+func TestDecisionRefusalD6(t *testing.T) {
+	f := lawFacts()
+	f.Decided, f.DecidedVerb, f.DecidedBy = true, "deny", "operator:terminal"
+	wantLaw(t, "D6", f, "D6", 409, "request "+lawRef+" was already decided deny by operator:terminal")
+}
+
+func TestDecisionRefusalD7(t *testing.T) {
+	f := lawFacts()
+	f.Now = f.Deadline
+	wantLaw(t, "at the deadline", f, "D7", 409, "deadline passed at 1970-01-01T01:23:20Z", "Cancel")
+	f.Now = f.Deadline - 1
+	f.Expiry = f.Deadline
+	if r := DecisionRefusal(f); r.Step != "" {
+		t.Errorf("one second before the deadline (control): %+v", r)
+	}
+}
+
+func TestDecisionRefusalD8(t *testing.T) {
+	f := lawFacts()
+	f.PostedPacket[0] ^= 1
+	wantLaw(t, "P", f, "D8", 409, "reload it")
+	f = lawFacts()
+	f.PostedRequest[31] ^= 1
+	wantLaw(t, "Q", f, "D8", 409, "reload it")
+}
+
+func TestDecisionRefusalD9(t *testing.T) {
+	for _, verb := range []string{"", "Approve", "yes", "approve "} {
+		f := lawFacts()
+		f.Verb = verb
+		wantLaw(t, "verb "+verb, f, "D9", 400, "a decision is approve or deny")
+	}
+}
+
+// D10: now < expiry <= min(now + 600, deadline), each edge on both sides.
+func TestDecisionRefusalD10(t *testing.T) {
+	for _, tc := range []struct {
+		now, deadline, expiry int64
+		ok                    bool
+	}{
+		{1000, 5000, 1000, false}, {1000, 5000, 1001, true}, {1000, 5000, 1600, true}, {1000, 5000, 1601, false},
+		{1000, 1300, 1300, true}, {1000, 1300, 1301, false}, {1000, 5000, 999, false},
+	} {
+		f := lawFacts()
+		f.Now, f.Deadline, f.Expiry = tc.now, tc.deadline, tc.expiry
+		r := DecisionRefusal(f)
+		if tc.ok != (r.Step == "") || !tc.ok && (r.Step != "D10" || r.Status != 409 || !strings.Contains(r.Reason, "at most 10 minutes")) {
+			t.Errorf("now %d deadline %d expiry %d: %+v, want ok=%v", tc.now, tc.deadline, tc.expiry, r, tc.ok)
+		}
+	}
+	if got := DecisionWindowMax(1000, 5000); got != 1600 {
+		t.Errorf("window max = %d, want 1600", got)
+	}
+	if got := DecisionWindowMax(1000, 1300); got != 1300 {
+		t.Errorf("window max = %d, want 1300", got)
+	}
+}
+
+// DL11: with two facts false, the lower step wins (each adjacent pair, D1 ... D10).
+func TestDecisionRefusalOrder(t *testing.T) {
+	breaks := []func(*DecisionFacts){
+		func(f *DecisionFacts) { f.RefOK = false },
+		func(f *DecisionFacts) { f.RequestOK = false },
+		func(f *DecisionFacts) { f.OnChain = false },
+		func(f *DecisionFacts) { f.Discovery.Complete = false },
+		func(f *DecisionFacts) { f.Discovery.Carrying = 2 },
+		func(f *DecisionFacts) { f.Scope = "github-pr:o/r" },
+		func(f *DecisionFacts) { f.Decided = true },
+		func(f *DecisionFacts) { f.Now = f.Deadline + 1; f.Expiry = f.Deadline + 2 },
+		func(f *DecisionFacts) { f.PostedPacket[0] ^= 1 },
+		func(f *DecisionFacts) { f.Verb = "maybe" },
+		func(f *DecisionFacts) { f.Expiry = 1 },
+	}
+	steps := []string{"D1", "D2", "D3", "D4a", "D4b", "D5", "D6", "D7", "D8", "D9", "D10"}
+	for i := 0; i+1 < len(breaks); i++ {
+		for j := i + 1; j < len(breaks); j++ {
+			f := lawFacts()
+			breaks[j](&f)
+			breaks[i](&f)
+			if r := DecisionRefusal(f); r.Step != steps[i] {
+				t.Errorf("%s and %s false: %s won", steps[i], steps[j], r.Step)
+			}
+		}
+	}
+	// The subject law is the decision law's first seven steps, and stops there.
+	f := lawFacts()
+	f.Verb, f.PostedPacket = "maybe", [32]byte{}
+	if r := SubjectRefusal(f); r.Step != "" {
+		t.Errorf("SubjectRefusal reads a posted fact: %+v", r)
+	}
+	f.Decided = true
+	if r := SubjectRefusal(f); r.Step != "D6" {
+		t.Errorf("SubjectRefusal D6: %+v", r)
+	}
+}
+
+// DL12: a well-formed approve and a well-formed deny are accepted (the control).
+func TestDecisionRefusalMustLive(t *testing.T) {
+	for _, verb := range []string{"approve", "deny"} {
+		f := lawFacts()
+		f.Verb = verb
+		if r := DecisionRefusal(f); r != (LawRefusal{}) {
+			t.Errorf("%s: %+v, want accepted", verb, r)
+		}
+		if r := SubjectRefusal(f); r != (LawRefusal{}) {
+			t.Errorf("subject %s: %+v, want accepted", verb, r)
+		}
+	}
+}
+
+// R4a: L1's publish and finish shapes, their near misses, and unreserved controls.
+func TestReservedApprovalScope(t *testing.T) {
+	id := PublishIdentity{RegistryOrigin: "https://registry.example", Vendor: "v", Name: "n", Version: "1.0.0"}
+	for _, tc := range []struct{ scope, owner string }{
+		{PublishApprovalScopeFor(id, PublishHashes{TarballSHA256: "t", ContentHash: "c", InterfaceHash: "i"}, 9), "publish"},
+		{"pkg#effect=Other&manifest=m&tarball=t&content=c&interface=i&expires=9", "publish"}, // the parse arm alone
+		{"registry:x", "publish"},                       // the prefix arm alone
+		{"worktree#effect=Registry.Publish", "publish"}, // the effect arm alone
+		{FormatFinishApprovalScope(FinishApprovalScope{Repository: "o/r", Remote: "https://github.com/o/r.git", Head: "world/e", SHA: strings.Repeat("a", 40),
+			Base: "dev", TitleSHA256: strings.Repeat("b", 64), BodySHA256: strings.Repeat("c", 64), ExpiresAt: 9}), "finish"},
+		{"github-pr:x", "finish"}, // the finish prefix arm alone
+		{"worktree", ""}, {"world", ""}, {"", ""}, {"worktree#x", ""}, {"Registry.Publish", ""},
+	} {
+		owner, reserved := ReservedApprovalScope(tc.scope)
+		if owner != tc.owner || reserved != (tc.owner != "") {
+			t.Errorf("%q: %q %v, want %q", tc.scope, owner, reserved, tc.owner)
+		}
+	}
+}
````

### Patch M1c — the complete packet discovery

````diff M1c
diff --git a/host/coordinator/decisionpacket_discover.go b/host/coordinator/decisionpacket_discover.go
new file mode 100644
index 0000000..54ebffc
--- /dev/null
+++ b/host/coordinator/decisionpacket_discover.go
@@ -0,0 +1,63 @@
+package coordinator
+
+// Row 201 A4 (design_docs/planned/w-world-approvals-by-key-a4.md §2.4 [r1], D4a): the complete,
+// request-wide decision-packet discovery a decision runs under humanMu. Packets have no request
+// index and ValidatePacket's duplicate scan is a 100-packet prefix, so a decision never trusts a
+// bounded view: it pages through every stored packet, up to DecisionDiscoveryLimit plus one
+// sentinel read, and reports whether the set it counted is complete.
+
+import (
+	"context"
+
+	"github.com/sunholo-data/ailang-world/host/broker"
+	"github.com/sunholo-data/ailang-world/host/hashref"
+	"github.com/sunholo-data/ailang-world/host/store"
+)
+
+// DiscoveryReader is the read surface the discovery needs; PacketReader and the daemon's
+// readStore satisfy it.
+type DiscoveryReader interface {
+	GetObject(ctx context.Context, ref hashref.HashRef) (store.Object, bool, error)
+	ObjectsBySemanticID(ctx context.Context, id, after string, limit int) ([]store.Object, error)
+}
+
+// DiscoverRequestPackets counts the stored decision packets that name request, in keyset pages of
+// store.MaxSemanticIDPage following the cursor. Complete is true only when a page came back short
+// before the scan reached broker.DecisionDiscoveryLimit packets; a packet past the limit is the
+// overflow sentinel, and the count is then not established. carrying is the last carrying packet
+// read (the one, when Carrying is 1). Store errors are returned, never a count.
+func DiscoverRequestPackets(ctx context.Context, r DiscoveryReader, request hashref.HashRef) (broker.Discovery, store.Object, error) {
+	var d broker.Discovery
+	var carrying store.Object
+	after := ""
+	for {
+		page, err := r.ObjectsBySemanticID(ctx, DecisionPacketV1, after, store.MaxSemanticIDPage)
+		if err != nil {
+			return broker.Discovery{}, store.Object{}, err
+		}
+		for _, o := range page {
+			if d.Scanned >= broker.DecisionDiscoveryLimit {
+				return d, carrying, nil // the sentinel: more packets than the limit
+			}
+			d.Scanned++
+			full, ok, err := r.GetObject(ctx, o.Hash)
+			if err != nil {
+				return broker.Discovery{}, store.Object{}, err
+			}
+			p, decodeErr := DecodeDecisionPacket(full.Payload)
+			if !ok || decodeErr != nil {
+				d.Undecodable++
+				continue
+			}
+			if p.RequestRef == request {
+				d.Carrying++
+				carrying = full
+			}
+		}
+		if len(page) < store.MaxSemanticIDPage {
+			d.Complete = true
+			return d, carrying, nil
+		}
+		after = page[len(page)-1].Hash.String()
+	}
+}
diff --git a/host/coordinator/decisionpacket_discover_test.go b/host/coordinator/decisionpacket_discover_test.go
new file mode 100644
index 0000000..7d993d6
--- /dev/null
+++ b/host/coordinator/decisionpacket_discover_test.go
@@ -0,0 +1,162 @@
+package coordinator
+
+// Row 201 A4 sprint 1 (design §2.4 [r1], §6 DS1): the discovery's completeness, its sentinel, its
+// cursor and its undecodable count, over a fake keyset reader with the store's page contract.
+
+import (
+	"context"
+	"fmt"
+	"sort"
+	"strings"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/broker"
+	"github.com/sunholo-data/ailang-world/host/hashref"
+	"github.com/sunholo-data/ailang-world/host/store"
+)
+
+// keysetPackets is a fake DiscoveryReader: packets sorted by hash_ref, served in keyset pages of at
+// most store.MaxSemanticIDPage; it records every after it was asked for.
+type keysetPackets struct {
+	objs   []store.Object
+	afters []string
+	gets   int
+}
+
+func (k *keysetPackets) ObjectsBySemanticID(_ context.Context, id, after string, limit int) ([]store.Object, error) {
+	if id != DecisionPacketV1 || limit < 1 || limit > store.MaxSemanticIDPage {
+		return nil, fmt.Errorf("unexpected page %q %d", id, limit)
+	}
+	k.afters = append(k.afters, after)
+	var page []store.Object
+	for _, o := range k.objs {
+		if o.Hash.String() > after && len(page) < limit {
+			page = append(page, store.Object{Hash: o.Hash, SemanticID: o.SemanticID})
+		}
+	}
+	return page, nil
+}
+
+func (k *keysetPackets) GetObject(_ context.Context, ref hashref.HashRef) (store.Object, bool, error) {
+	k.gets++
+	for _, o := range k.objs {
+		if o.Hash == ref {
+			return o, true, nil
+		}
+	}
+	return store.Object{}, false, nil
+}
+
+// discoveryRig stores n packets, the ones whose index is in carrying naming request, the rest
+// naming their own request; undecodable payloads replace the packets whose index is in bad.
+func discoveryRig(n int, request hashref.HashRef, carrying, bad map[int]bool) *keysetPackets {
+	k := &keysetPackets{}
+	for i := 0; i < n; i++ {
+		req := hashref.SumSHA256([]byte(fmt.Sprintf("other-request-%d", i)))
+		if carrying[i] {
+			req = request
+		}
+		payload := EncodeDecisionPacket(DecisionPacket{ProposalHash: hashref.SumSHA256([]byte(fmt.Sprintf("plan-%d", i))), RequestRef: req,
+			CreatedAt: 1, DeadlineAt: 2, Policy: "Cancel"})
+		if bad[i] {
+			payload = []byte(fmt.Sprintf(`{"not":"a packet %d"}`, i))
+		}
+		k.objs = append(k.objs, store.Object{Hash: hashref.SumSHA256(payload), SemanticID: DecisionPacketV1, Payload: payload})
+	}
+	sort.Slice(k.objs, func(i, j int) bool { return k.objs[i].Hash.String() < k.objs[j].Hash.String() })
+	return k
+}
+
+// DS1: 0, 499, 500, 1,999 and exactly 2,000 packets are complete (the sentinel read returns none);
+// 2,001 are not; the cursor follows after across pages; an undecodable payload is counted.
+func TestDiscoverRequestPacketsCompleteness(t *testing.T) {
+	ctx := context.Background()
+	request := hashref.SumSHA256([]byte("the-request"))
+	for _, tc := range []struct {
+		n, carryAt, wantScanned, wantPages int
+		complete                           bool
+	}{
+		{0, -1, 0, 1, true}, {499, 0, 499, 1, true}, {500, 499, 500, 2, true}, {1999, 1998, 1999, 4, true},
+		{2000, 1999, 2000, 5, true}, {2001, 1999, 2000, 5, false},
+	} {
+		k := discoveryRig(tc.n, request, map[int]bool{tc.carryAt: true}, nil)
+		d, carrying, err := DiscoverRequestPackets(ctx, k, request)
+		if err != nil {
+			t.Fatal(err)
+		}
+		wantCarrying := 0
+		if tc.carryAt >= 0 {
+			wantCarrying = 1
+		}
+		if d.Complete != tc.complete || d.Scanned != tc.wantScanned || d.Carrying != wantCarrying || d.Undecodable != 0 || len(k.afters) != tc.wantPages {
+			t.Errorf("%d packets: %+v after %d pages, want complete=%v scanned=%d carrying=%d pages=%d", tc.n, d, len(k.afters), tc.complete, tc.wantScanned, wantCarrying, tc.wantPages)
+		}
+		if wantCarrying == 1 && d.Complete {
+			p, err := DecodeDecisionPacket(carrying.Payload)
+			if err != nil || p.RequestRef != request {
+				t.Errorf("%d packets: the carrying packet is %v, %v", tc.n, p.RequestRef, err)
+			}
+		}
+		// The cursor: each page starts strictly after the last hash of the one before.
+		for i := 1; i < len(k.afters); i++ {
+			if want := k.objs[i*store.MaxSemanticIDPage-1].Hash.String(); k.afters[i] != want {
+				t.Errorf("%d packets: page %d after %q, want %q", tc.n, i, k.afters[i], want)
+			}
+		}
+		if k.gets != d.Scanned {
+			t.Errorf("%d packets: %d payload reads for %d scanned", tc.n, k.gets, d.Scanned)
+		}
+	}
+	// Two carrying packets, one past the first page, are both counted.
+	k := discoveryRig(1200, request, map[int]bool{3: true, 1100: true}, nil)
+	if d, _, err := DiscoverRequestPackets(ctx, k, request); err != nil || d != (broker.Discovery{Complete: true, Scanned: 1200, Carrying: 2}) {
+		t.Errorf("two carrying: %+v, %v", d, err)
+	}
+	// An undecodable payload is counted, not skipped; so is a listed packet whose payload is gone.
+	k = discoveryRig(600, request, map[int]bool{7: true}, map[int]bool{550: true})
+	if d, _, err := DiscoverRequestPackets(ctx, k, request); err != nil || !d.Complete || d.Undecodable != 1 || d.Carrying != 1 || d.Scanned != 600 {
+		t.Errorf("undecodable: %+v, %v", d, err)
+	}
+	k = discoveryRig(3, request, map[int]bool{0: true}, nil)
+	gone := k.objs[2].Hash
+	if d, _, err := DiscoverRequestPackets(ctx, &missingPayload{k, gone}, request); err != nil || d.Undecodable != 1 || d.Scanned != 3 {
+		t.Errorf("payload gone: %+v, %v", d, err)
+	}
+	// A store error is returned, never a count.
+	if _, _, err := DiscoverRequestPackets(ctx, failingPages{}, request); err == nil || !strings.Contains(err.Error(), "page failed") {
+		t.Errorf("store error: %v", err)
+	}
+	if _, _, err := DiscoverRequestPackets(ctx, failingGets{discoveryRig(3, request, nil, nil)}, request); err == nil || !strings.Contains(err.Error(), "read failed") {
+		t.Errorf("payload read error: %v", err)
+	}
+}
+
+// failingGets lists packets but fails every payload read.
+type failingGets struct{ *keysetPackets }
+
+func (failingGets) GetObject(context.Context, hashref.HashRef) (store.Object, bool, error) {
+	return store.Object{}, false, fmt.Errorf("read failed")
+}
+
+// missingPayload lists every packet but answers not-found for one of them.
+type missingPayload struct {
+	*keysetPackets
+	gone hashref.HashRef
+}
+
+func (m *missingPayload) GetObject(ctx context.Context, ref hashref.HashRef) (store.Object, bool, error) {
+	if ref == m.gone {
+		return store.Object{}, false, nil
+	}
+	return m.keysetPackets.GetObject(ctx, ref)
+}
+
+type failingPages struct{}
+
+func (failingPages) ObjectsBySemanticID(context.Context, string, string, int) ([]store.Object, error) {
+	return nil, fmt.Errorf("page failed")
+}
+
+func (failingPages) GetObject(context.Context, hashref.HashRef) (store.Object, bool, error) {
+	return store.Object{}, false, nil
+}
````

### Patch M2 — the broker rules (first decision final, L2, the caller census)

````diff M2
diff --git a/host/broker/approve.go b/host/broker/approve.go
index cb4674a..33a789b 100644
--- a/host/broker/approve.go
+++ b/host/broker/approve.go
@@ -26,6 +26,8 @@ const (
 var (
 	ErrApprovalRequestNotFound = errors.New("broker: approval request not found")
 	ErrInvalidApprovalDecision = errors.New("broker: invalid approval decision")
+	// ErrApprovalAlreadyDecided: the first decision of a request is final (row 201 A4 §2.4).
+	ErrApprovalAlreadyDecided = errors.New("broker: approval request already decided")
 
 	// The five publish-approval refusal classes (SM.B2b / AC9). They are
 	// separate sentinels rather than one because a caller that cannot tell
@@ -213,6 +215,16 @@ func decideApproval(ctx context.Context,
 	payload := mustApprovalJSON(approvalDecisionWire{
 		RequestRef: requestRef.String(), Decision: decision, DecidedBy: decidedBy, Now: now,
 	})
+	// Row 201 A4 §2.4: the first decision is final. Nothing is written for a second one; the
+	// same decision again (the same verb, decider and time: byte-identical, so the same object)
+	// answers its existing ref, as it did before the rule (two identical requests share one ref).
+	if existing, decided, err := findApprovalDecision(ctx, s, requestRef); err != nil {
+		return hashref.HashRef{}, err
+	} else if decided && bytes.Equal(existing.Payload, payload) {
+		return existing.Hash, nil
+	} else if decided {
+		return hashref.HashRef{}, ErrApprovalAlreadyDecided
+	}
 	obj := brokerObject(ApprovalDecisionV1, payload)
 	if err := s.PutObject(ctx, obj); err != nil {
 		return hashref.HashRef{}, fmt.Errorf("broker: put approval decision: %w", err)
@@ -223,6 +235,24 @@ func decideApproval(ctx context.Context,
 	return obj.Hash, nil
 }
 
+// ApprovalChainState reads requestRef's place on the approvals chain (row 201 A4 §2.5 D3, D6):
+// whether it is an approval request of the chain, and the decision that landed for it, if any.
+// Pass the caller's intended read ctx; this API supplies no default timeout.
+func ApprovalChainState(ctx context.Context, s *store.Store, requestRef hashref.HashRef) (onChain, decided bool, verb, decidedBy string, err error) {
+	if _, onChain, err = findApprovalRequest(ctx, s, requestRef); err != nil || !onChain {
+		return false, false, "", "", err
+	}
+	obj, decided, err := findApprovalDecision(ctx, s, requestRef)
+	if err != nil || !decided {
+		return true, false, "", "", err
+	}
+	var w approvalDecisionWire
+	if err := decodeApprovalJSON(obj.Payload, &w); err != nil {
+		return false, false, "", "", err
+	}
+	return true, true, w.Decision, w.DecidedBy, nil
+}
+
 func appendApprovalHead(ctx context.Context, s approvalStore, requestRef, decisionRef hashref.HashRef) error {
 	previous, ok, err := s.GetRegistryHead(ctx, ApprovalsV1)
 	if err != nil {
@@ -545,6 +575,12 @@ func validatePublishApproval(ctx context.Context,
 		return hashref.HashRef{}, fmt.Errorf(
 			"%w: decision %s says %q", ErrPublishApprovalDenied, id.ApprovalRef.String(), decision.Decision)
 	}
+	// Row 201 A4 §2.6 L2: a decision a key made never stamps a publish.
+	if strings.HasPrefix(decision.DecidedBy, "ssh:") || strings.HasPrefix(decision.DecidedBy, "webauthn:") {
+		return hashref.HashRef{}, fmt.Errorf(
+			"%w: publishing stays typed-only (D-WORLD-89 R4): decision %s was made by %s",
+			ErrPublishApprovalDenied, id.ApprovalRef.String(), decision.DecidedBy)
+	}
 
 	requestRef, err := hashref.Parse(decision.RequestRef)
 	if err != nil {
diff --git a/host/broker/context_thread_test.go b/host/broker/context_thread_test.go
index 1dbf134..6973960 100644
--- a/host/broker/context_thread_test.go
+++ b/host/broker/context_thread_test.go
@@ -57,8 +57,9 @@ func TestApprovalContextIdentity(t *testing.T) {
 			}
 			// Mint traverses Execute, append, decide, findRequest, walk's request and
 			// decision branches, and both sessions. Exact read counts prevent dead arms.
-			if w.objects != 5 || w.heads != 4 {
-				t.Fatalf("mint read census objects=%d heads=%d, want 5/4", w.objects, w.heads)
+			// Row 201 A4 BR1: decide's first-decision-final walk reads one head and one object more.
+			if w.objects != 6 || w.heads != 5 {
+				t.Fatalf("mint read census objects=%d heads=%d, want 6/5", w.objects, w.heads)
 			}
 			w.objects = 0
 			w.heads = 0
diff --git a/host/broker/decide_once_test.go b/host/broker/decide_once_test.go
new file mode 100644
index 0000000..a764394
--- /dev/null
+++ b/host/broker/decide_once_test.go
@@ -0,0 +1,230 @@
+package broker
+
+// Row 201 A4 sprint 1 (design §2.4, §2.6 L2, §5.1, §6 BR1, BR2, R4b, R4c, XG1): the first
+// decision of a request is final; a key decision never stamps a publish; who may call the decide
+// entry point.
+
+import (
+	"errors"
+	"go/ast"
+	"go/parser"
+	"go/token"
+	"os"
+	"path/filepath"
+	"reflect"
+	"runtime"
+	"sort"
+	"strings"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/hashref"
+	"github.com/sunholo-data/ailang-world/host/store"
+)
+
+// openDecideStore is a world store under t.TempDir().
+func openDecideStore(t *testing.T) *store.Store {
+	t.Helper()
+	s, err := store.Open(filepath.Join(t.TempDir(), "world.db"))
+	if err != nil {
+		t.Fatal(err)
+	}
+	t.Cleanup(func() { _ = s.Close() })
+	return s
+}
+
+// pendingRequest files one Human.Approve request through the landed handler.
+func pendingRequest(t *testing.T, s approvalStore, scope, requester string, now int64) hashref.HashRef {
+	t.Helper()
+	out, err := newHumanHandler(s).Execute(boundedTestContext(t), EffectRequest{Effect: EffectHumanApprove, Scope: scope, Cost: 1, Now: now},
+		mustApprovalJSON(approvalInputWire{Requester: requester}))
+	if err != nil {
+		t.Fatal(err)
+	}
+	ref, err := ParsePendingApproval(out)
+	if err != nil {
+		t.Fatal(err)
+	}
+	return ref
+}
+
+// BR1: deny, then approve → ErrApprovalAlreadyDecided, the head unchanged, the poll still deny (the
+// design's V9 probe as a test); a deny by another decider or at another time is refused too; the
+// byte-identical decision again answers its existing ref and writes nothing (plan §1.4 P2).
+func TestDecideApprovalRefusesASecondDecision(t *testing.T) {
+	s := openDecideStore(t)
+	ctx := boundedTestContext(t)
+	ref := pendingRequest(t, s, WorkspaceScope, "agent", 100)
+	first, err := DecideApproval(ctx, s, ref, "deny", "operator:terminal", 101)
+	if err != nil {
+		t.Fatal(err)
+	}
+	head, _, err := s.GetRegistryHead(ctx, ApprovalsV1)
+	if err != nil {
+		t.Fatal(err)
+	}
+	for _, tc := range []struct {
+		verb, by string
+		now      int64
+	}{{"approve", "ssh:mallory", 102}, {"approve", "operator:terminal", 101}, {"deny", "operator:terminal", 102}, {"deny", "ssh:mallory", 101}} {
+		if got, err := DecideApproval(ctx, s, ref, tc.verb, tc.by, tc.now); !errors.Is(err, ErrApprovalAlreadyDecided) || !got.IsZero() {
+			t.Fatalf("second decision %+v: %v, %v", tc, got, err)
+		}
+	}
+	if again, err := DecideApproval(ctx, s, ref, "deny", "operator:terminal", 101); err != nil || again != first {
+		t.Fatalf("the identical decision again: %v, %v; want %v", again, err, first)
+	}
+	if after, _, err := s.GetRegistryHead(ctx, ApprovalsV1); err != nil || after != head {
+		t.Fatalf("head moved: %v -> %v (%v)", head, after, err)
+	}
+	out, err := newHumanHandler(s).Execute(ctx, EffectRequest{Effect: EffectHumanPollApproval, Scope: WorkspaceScope, Now: 103},
+		mustApprovalJSON(approvalInputWire{RequestRef: ref.String()}))
+	if err != nil || !strings.Contains(string(out), `"decision":"deny","decidedBy":"operator:terminal"`) {
+		t.Fatalf("poll = %s, %v", out, err)
+	}
+	if on, decided, verb, by, err := ApprovalChainState(ctx, s, ref); err != nil || !on || !decided || verb != "deny" || by != "operator:terminal" {
+		t.Fatalf("chain state = %v %v %q %q %v", on, decided, verb, by, err)
+	}
+}
+
+// BR2: an undecided request decides (the control), and the chain state reads it before and after;
+// a ref that is not on the chain reads as absent.
+func TestDecideApprovalFirstDecisionMustLive(t *testing.T) {
+	s := openDecideStore(t)
+	ctx := boundedTestContext(t)
+	ref := pendingRequest(t, s, WorkspaceScope, "agent", 100)
+	if on, decided, _, _, err := ApprovalChainState(ctx, s, ref); err != nil || !on || decided {
+		t.Fatalf("pending chain state = %v %v %v", on, decided, err)
+	}
+	other := pendingRequest(t, s, WorkspaceScope, "other", 100)
+	if _, err := DecideApproval(ctx, s, other, "deny", "operator:terminal", 101); err != nil {
+		t.Fatal(err)
+	}
+	dec, err := DecideApproval(ctx, s, ref, "approve", "operator:terminal", 102)
+	if err != nil || dec.IsZero() {
+		t.Fatalf("first decision: %v, %v", dec, err)
+	}
+	if on, decided, verb, by, err := ApprovalChainState(ctx, s, ref); err != nil || !on || !decided || verb != "approve" || by != "operator:terminal" {
+		t.Fatalf("decided chain state = %v %v %q %q %v", on, decided, verb, by, err)
+	}
+	if on, decided, _, _, err := ApprovalChainState(ctx, s, hashref.SumSHA256([]byte("nowhere"))); err != nil || on || decided {
+		t.Fatalf("absent chain state = %v %v %v", on, decided, err)
+	}
+	if _, _, _, _, err := ApprovalChainState(ctx, openDecideStore(t), ref); err != nil {
+		t.Fatalf("an empty chain: %v", err)
+	}
+}
+
+// landPublishDecidedBy lands a publish stamp for f decided by decidedBy, through the landed request
+// handler and decideApproval, and returns the payload naming it.
+func landPublishDecidedBy(t *testing.T, base approvalStore, f publishFixture, decidedBy string) []byte {
+	t.Helper()
+	request := pendingRequest(t, base, PublishApprovalScopeFor(f.identity, f.hashes, 100), "r4-"+decidedBy, 10)
+	decision, err := decideApproval(boundedTestContext(t), base, request, "approve", decidedBy, 11)
+	if err != nil {
+		t.Fatal(err)
+	}
+	f.identity.ApprovalRef = decision
+	return EncodePublishPayload(f.identity, f.hashes)
+}
+
+// R4b: a publish-shaped request decided by a key (ssh:, webauthn:) never stamps a publish.
+func TestPublishApprovalRefusesAKeyDecision(t *testing.T) {
+	base := openPublishStore(t)
+	f := newPublishFixture(t, "http://127.0.0.1:1", "r4b")
+	req := EffectRequest{Effect: EffectRegistryPublish, Scope: f.scope, Cost: PublishCost, Now: 50}
+	for _, by := range []string{"ssh:mark", "webauthn:mark-mac"} {
+		_, err := validatePublishApproval(boundedTestContext(t), base, landPublishDecidedBy(t, base, f, by), req)
+		if !errors.Is(err, ErrPublishApprovalDenied) || !strings.Contains(err.Error(), "publishing stays typed-only (D-WORLD-89 R4)") || !strings.Contains(err.Error(), "was made by "+by) {
+			t.Errorf("%s: %v", by, err)
+		}
+	}
+}
+
+// R4c: the same stamp decided at a terminal, or by a free-text name, is accepted (the control).
+func TestPublishApprovalTypedStampMustLive(t *testing.T) {
+	base := openPublishStore(t)
+	f := newPublishFixture(t, "http://127.0.0.1:1", "r4c")
+	req := EffectRequest{Effect: EffectRegistryPublish, Scope: f.scope, Cost: PublishCost, Now: 50}
+	for _, by := range []string{"operator:terminal", "mark", "attended-operator", "ssh", "sshmark", "webauthn"} {
+		if _, err := validatePublishApproval(boundedTestContext(t), base, landPublishDecidedBy(t, base, f, by), req); err != nil {
+			t.Errorf("%s: %v", by, err)
+		}
+	}
+}
+
+// decideCallers is the census of non-test calls of DecideApproval and decideApproval in the module,
+// as "<file>:<enclosing func>" (design §5.1).
+var decideCallers = []string{
+	"host/broker/approve.go:DecideApproval",
+	"host/broker/approve.go:RecordOperatorApproval",
+	"host/broker/publish_op.go:mintAttendedApproval",
+}
+
+// XG1: the only non-test callers of the decide entry point are the ones named above: no effect
+// handler, MCP, AG-UI or A2A route can decide a request.
+func TestDecideApprovalCallerCensus(t *testing.T) {
+	_, file, _, ok := runtime.Caller(0)
+	if !ok {
+		t.Fatal("runtime.Caller")
+	}
+	root := filepath.Dir(filepath.Dir(filepath.Dir(file)))
+	fset := token.NewFileSet()
+	var got []string
+	scanned := 0
+	err := filepath.WalkDir(root, func(path string, e os.DirEntry, err error) error {
+		if err != nil {
+			return err
+		}
+		if e.IsDir() {
+			if path != root && (strings.HasPrefix(e.Name(), ".") || e.Name() == "node_modules" || e.Name() == "vendor" || e.Name() == "testdata") {
+				return filepath.SkipDir
+			}
+			return nil
+		}
+		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
+			return nil
+		}
+		f, err := parser.ParseFile(fset, path, nil, parser.SkipObjectResolution)
+		if err != nil {
+			return err
+		}
+		scanned++
+		rel, _ := filepath.Rel(root, path)
+		for _, decl := range f.Decls {
+			fn, ok := decl.(*ast.FuncDecl)
+			if !ok || fn.Body == nil {
+				continue
+			}
+			ast.Inspect(fn.Body, func(n ast.Node) bool {
+				call, ok := n.(*ast.CallExpr)
+				if !ok {
+					return true
+				}
+				name := ""
+				switch fun := call.Fun.(type) {
+				case *ast.Ident:
+					name = fun.Name
+				case *ast.SelectorExpr:
+					name = fun.Sel.Name
+				}
+				if name == "DecideApproval" || name == "decideApproval" {
+					got = append(got, filepath.ToSlash(rel)+":"+fn.Name.Name)
+				}
+				return true
+			})
+		}
+		return nil
+	})
+	if err != nil {
+		t.Fatal(err)
+	}
+	if scanned < 100 {
+		t.Fatalf("scanned %d Go files under %s; the census did not see the module", scanned, root)
+	}
+	sort.Strings(got)
+	want := append([]string(nil), decideCallers...)
+	sort.Strings(want)
+	if !reflect.DeepEqual(got, want) {
+		t.Fatalf("decide callers = %q, want %q", got, want)
+	}
+}
````

### Patch M3a — the decision routes, the guards, the typed arm

````diff M3a
diff --git a/host/broker/decide_once_test.go b/host/broker/decide_once_test.go
index a764394..dce1495 100644
--- a/host/broker/decide_once_test.go
+++ b/host/broker/decide_once_test.go
@@ -158,6 +158,7 @@ var decideCallers = []string{
 	"host/broker/approve.go:DecideApproval",
 	"host/broker/approve.go:RecordOperatorApproval",
 	"host/broker/publish_op.go:mintAttendedApproval",
+	"host/daemon/surface_decide.go:handleSurfaceDecide",
 }
 
 // XG1: the only non-test callers of the decide entry point are the ones named above: no effect
diff --git a/host/daemon/daemon.go b/host/daemon/daemon.go
index 5a34ff8..ec62186 100644
--- a/host/daemon/daemon.go
+++ b/host/daemon/daemon.go
@@ -1100,7 +1100,8 @@ func releaseFromVersion(version string) string {
 // /surface/app.css and /surface/ailang-logo.svg (the app, N2b), and GET /surface/passkey with
 // /surface/passkey.js (the passkey ceremony page, row 201 A2: static, it reads nothing), and the
 // approval routes of row 201 A3 (requests.json, requests/{id}, the approve POST, the card page
-// GET /surface/approve/{id} and its approve.js), each behind approvalRoute's Host check.
+// GET /surface/approve/{id} and its approve.js), each behind approvalRoute's Host check, and
+// row 201 A4's decision routes GET /surface/decisions/{ref} and POST /surface/decide/{ref}.
 // The registry pattern deliberately uses a multi-segment wildcard: registry semantic IDs
 // such as "world/epoch-registry/v1" contain slashes.
 func (d *Daemon) Handler() http.Handler {
@@ -1139,6 +1140,8 @@ func (d *Daemon) Handler() http.Handler {
 	mux.HandleFunc("POST /surface/approve/{id}", d.approvalRoute(d.handleSurfaceApprove))
 	mux.HandleFunc("GET /surface/approve/{id}", d.approvalRoute(d.handleSurfaceApprovePage))
 	mux.HandleFunc("GET /surface/approve.js", d.approvalRoute(d.handleSurfaceApproveScript))
+	mux.HandleFunc("GET /surface/decisions/{ref}", d.approvalRoute(d.handleSurfaceDecision))
+	mux.HandleFunc("POST /surface/decide/{ref}", d.approvalRoute(d.handleSurfaceDecide))
 	// The two A2A projection routes (w-a2a-session-projection P6.B-A2A-CARD)
 	// are ADDITIVE: the frozen /v1/ table above is untouched, and the routes
 	// are NOT in isProtected — the projection handler resolves the session
diff --git a/host/daemon/surface_decide.go b/host/daemon/surface_decide.go
new file mode 100644
index 0000000..7f2ccf4
--- /dev/null
+++ b/host/daemon/surface_decide.go
@@ -0,0 +1,190 @@
+package daemon
+
+// Row 201 A4 sprint 1 (design_docs/planned/w-world-approvals-by-key-a4.md §2.4, §2.5, §2.8): an
+// operator decides a Human.Approve request whose decision packet the coordinator committed.
+// GET /surface/decisions/{ref} serves the subject (the request's and the packet's stored bytes, the
+// expiry window, the law's display verdict); POST /surface/decide/{ref} runs A3a's guards, then,
+// under humanMu, the decision law (D1-D10) over facts read from the world store, then the format
+// arm (D11), then broker.DecideApproval with a decidedBy the daemon derives. Sprint 1's only arm is
+// typed: the row-236 operator proof the CLI writes after a typed y. Neither route is under /v1.
+
+import (
+	"context"
+	"encoding/base64"
+	"encoding/hex"
+	"fmt"
+	"io"
+	"net/http"
+	"regexp"
+
+	"github.com/sunholo-data/ailang-world/host/approval"
+	"github.com/sunholo-data/ailang-world/host/broker"
+	"github.com/sunholo-data/ailang-world/host/coordinator"
+	"github.com/sunholo-data/ailang-world/host/hashref"
+	"github.com/sunholo-data/ailang-world/host/store"
+	"github.com/sunholo-data/ailang-world/host/surfacepage"
+)
+
+// decisionRefPattern is D1's grammar: a request ref is sha256: and 64 lowercase hex.
+var decisionRefPattern = regexp.MustCompile(`^sha256:[0-9a-f]{64}$`)
+
+// typedDecider is the decidedBy of a decision made by the typed arm.
+const typedDecider = "operator:terminal"
+
+// decisionRead is what readDecision found: the law's facts and, when exactly one packet carries
+// the request, the subject's objects.
+type decisionRead struct {
+	facts   broker.DecisionFacts
+	request store.Object
+	packet  store.Object
+	entry   int64
+}
+
+// readDecision reads every fact the law needs for request ref at now from the world store (design
+// §2.5): the request object, its place on the approvals chain, the complete packet discovery
+// (§2.4 [r1]) and, for exactly one carrying packet, its ValidatePacket verdict. A ref that is not
+// D1's grammar reads nothing. Store errors are returned, never a fact.
+func (d *Daemon) readDecision(ctx context.Context, ref string, now int64) (decisionRead, error) {
+	rd := decisionRead{facts: broker.DecisionFacts{Ref: ref, Now: now}}
+	f := &rd.facts
+	if !decisionRefPattern.MatchString(ref) { // D1
+		return rd, nil
+	}
+	f.RefOK = true
+	q := hashref.MustParse(ref)
+	digest, _ := hex.DecodeString(q.Digest())
+	f.Request = [32]byte(digest)
+	obj, ok, err := d.store.GetObject(ctx, q)
+	if err != nil {
+		return rd, err
+	}
+	if req, err := broker.DecodeApprovalRequest(obj.Payload); ok && obj.SemanticID == broker.ApprovalRequestV1 && err == nil && req.Effect == broker.EffectHumanApprove { // D2
+		f.RequestOK, f.Scope, rd.request = true, req.Scope, obj
+	}
+	if f.OnChain, f.Decided, f.DecidedVerb, f.DecidedBy, err = broker.ApprovalChainState(ctx, d.store, q); err != nil { // D3, D6
+		return rd, err
+	}
+	var carrying store.Object
+	if f.Discovery, carrying, err = coordinator.DiscoverRequestPackets(ctx, d.store, q); err != nil { // D4a
+		return rd, err
+	}
+	if f.Discovery.Complete && f.Discovery.Undecodable == 0 && f.Discovery.Carrying == 1 { // D4b, D7, D8
+		res, err := coordinator.ValidatePacket(ctx, d.store, carrying)
+		if err != nil {
+			return rd, err
+		}
+		pd, _ := hex.DecodeString(carrying.Hash.Digest())
+		f.PacketCode, f.PacketSentence, f.Deadline, f.Packet = res.Code, packetRefusalSentence(res), res.Packet.DeadlineAt, [32]byte(pd)
+		rd.packet, rd.entry = carrying, res.Entry
+	}
+	return rd, nil
+}
+
+// lawClass is the refusal class of a law status.
+func lawClass(status int) string {
+	switch status {
+	case http.StatusBadRequest:
+		return "BadRequest"
+	case http.StatusForbidden:
+		return "Forbidden"
+	case http.StatusNotFound:
+		return "NotFound"
+	}
+	return "Conflict"
+}
+
+// handleSurfaceDecision is GET /surface/decisions/{ref}: the subject of a decision, as bytes. A
+// request that cannot have a subject (D1-D4: no such request, not on the chain, not exactly one
+// validated packet) is the law's refusal; otherwise 200 with the law's D5-D7 verdict at now as
+// display (refusal, reserved, the decision when one landed). The POST runs the law again.
+func (d *Daemon) handleSurfaceDecision(w http.ResponseWriter, r *http.Request, _ []string) {
+	ref := r.PathValue("ref")
+	now := d.approveNow()
+	ctx, cancel := d.readCtx(r)
+	defer cancel()
+	rd, err := d.readDecision(ctx, ref, now)
+	if err != nil {
+		d.approvalInternal(w, r, err)
+		return
+	}
+	law := broker.SubjectRefusal(rd.facts)
+	switch law.Step {
+	case "D1", "D2", "D3", "D4a", "D4b":
+		writeApprovalRefusal(w, law.Status, lawClass(law.Status), law.Reason)
+		return
+	}
+	p, _ := coordinator.DecodeDecisionPacket(rd.packet.Payload) // validated: code 0
+	reserved, _ := broker.ReservedApprovalScope(rd.facts.Scope)
+	writeApprovalJSON(w, http.StatusOK, surfacepage.DecisionSubject{V: 1,
+		Request: ref, RequestBytes: base64.RawURLEncoding.EncodeToString(rd.request.Payload),
+		Packet: rd.packet.Hash.String(), PacketBytes: base64.RawURLEncoding.EncodeToString(rd.packet.Payload),
+		Entry: rd.entry, Proposal: p.ProposalHash.String(), Now: now,
+		Window:  surfacepage.DecisionWindow{Min: now + 1, Max: broker.DecisionWindowMax(now, p.DeadlineAt)},
+		Refusal: law.Reason, Reserved: reserved, DecidedBy: rd.facts.DecidedBy, Decision: rd.facts.DecidedVerb,
+		SSH: surfacepage.DecisionSSH{Namespace: approval.DecisionNamespace}})
+}
+
+// handleSurfaceDecide is POST /surface/decide/{ref} (design §2.8). After R2 (approvalRoute): R0 any
+// Authorization or Cookie is refused (the agent holds a session, so a session-bearing request is
+// one it could make), then A3a's R3-R6 verbatim. Then, under humanMu, in this order (sprint 2's
+// key arms take the arm's slot and nothing else moves; plan §1.3 S1-D9):
+//  1. the law, D1-D10, over facts read now (a refusal writes nothing and consumes no proof);
+//  2. the arm (D11): typed consumes the operator proof for this request and verb, unlinking it;
+//  3. broker.DecideApproval, which lands the first decision only, with decidedBy derived here.
+func (d *Daemon) handleSurfaceDecide(w http.ResponseWriter, r *http.Request, o []string) {
+	if len(r.Header.Values("Authorization")) != 0 || len(r.Header.Values("Cookie")) != 0 { // R0
+		writeApprovalRefusal(w, http.StatusForbidden, "Forbidden", "a decision is never made with a session or a cookie (the agent holds a session); decide at a terminal with ailang-worldd approve <request-ref>")
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
+	body, packet, request, err := surfacepage.ParseDecideBody(raw)
+	if err != nil {
+		writeApprovalRefusal(w, http.StatusBadRequest, "BadRequest", err.Error())
+		return
+	}
+	ref := r.PathValue("ref")
+	d.humanMu.Lock()
+	defer d.humanMu.Unlock()
+	ctx, cancel := context.WithTimeout(r.Context(), d.commitBudget)
+	defer cancel()
+	now := d.approveNow()
+	rd, err := d.readDecision(ctx, ref, now)
+	if err != nil {
+		d.approvalInternal(w, r, err)
+		return
+	}
+	f := rd.facts
+	f.PostedPacket, f.PostedRequest, f.Verb, f.Expiry = packet, request, body.Verb, body.Expiry
+	if law := broker.DecisionRefusal(f); law.Step != "" { // 1. the law
+		writeApprovalRefusal(w, law.Status, lawClass(law.Status), law.Reason)
+		return
+	}
+	q := hashref.MustParse(ref) // D1 held
+	// 2. the arm (D11). The grammar admits typed only in sprint 1; the phase carries the verb.
+	if err := d.consumeOperatorProof("decision-"+q.Digest(), "decide-"+body.Verb, r.Header.Get("X-World-Operator-Proof")); err != nil {
+		writeApprovalRefusal(w, http.StatusForbidden, "Forbidden", "the typed decision is not proven: "+err.Error()+"; decide at a terminal with ailang-worldd approve "+ref)
+		return
+	}
+	decision, err := broker.DecideApproval(ctx, d.store, q, body.Verb, typedDecider, now) // 3.
+	if err != nil {
+		d.approvalInternal(w, r, err)
+		return
+	}
+	writeApprovalJSON(w, http.StatusOK, surfacepage.DecisionMade{V: 1, Decided: body.Verb, Request: ref, Decision: decision.String(), DecidedBy: typedDecider})
+}
diff --git a/host/daemon/surface_decide_routes_test.go b/host/daemon/surface_decide_routes_test.go
new file mode 100644
index 0000000..5c508c2
--- /dev/null
+++ b/host/daemon/surface_decide_routes_test.go
@@ -0,0 +1,64 @@
+package daemon
+
+// Row 201 A4 sprint 1 (design §2.8, §6 RD1): the decision routes, under /surface, behind
+// approvalRoute; no other method; no CORS; /v1 is untouched (TestFrozenV1RouteTableMatchesMux).
+
+import (
+	"os"
+	"reflect"
+	"regexp"
+	"sort"
+	"strings"
+	"testing"
+)
+
+// RD1: exactly the two decision routes are registered, each behind approvalRoute; a foreign Host is
+// 421 on both; other methods are 405 with no CORS header; no other name under them is served.
+func TestDecisionRoutesRegistered(t *testing.T) {
+	src, err := os.ReadFile("daemon.go")
+	if err != nil {
+		t.Fatal(err)
+	}
+	body := strings.SplitN(strings.SplitN(string(src), "func (d *Daemon) Handler() http.Handler {", 2)[1], "\n}\n", 2)[0]
+	var got []string
+	for _, m := range regexp.MustCompile(`mux\.HandleFunc\("([A-Z]+ /surface/deci(?:de|sions)[^"]*)", (d\.approvalRoute\()?`).FindAllStringSubmatch(body, -1) {
+		if m[2] == "" {
+			t.Errorf("%s is not behind approvalRoute", m[1])
+		}
+		got = append(got, m[1])
+	}
+	sort.Strings(got)
+	if want := []string{"GET /surface/decisions/{ref}", "POST /surface/decide/{ref}"}; !reflect.DeepEqual(got, want) {
+		t.Fatalf("decision routes %v, want %v", got, want)
+	}
+	r := newApprovalRig(t)
+	ref := "sha256:" + strings.Repeat("ab", 32)
+	for _, tc := range []struct{ method, target string }{{"GET", "/surface/decisions/" + ref}, {"POST", "/surface/decide/" + ref}} {
+		requireRefusal(t, "foreign "+tc.target, r.serve(tc.method, tc.target, "rebind.example:7644", nil, "Origin", "http://127.0.0.1:7644"), 421, "MisdirectedRequest")
+	}
+	requireRefusal(t, "unknown subject", r.serve("GET", "/surface/decisions/"+ref, apLoopback, nil), 404, "NotFound")
+	for _, tc := range []struct {
+		target  string
+		methods []string
+	}{
+		{"/surface/decisions/" + ref, []string{"POST", "PUT", "DELETE", "OPTIONS"}},
+		{"/surface/decide/" + ref, []string{"GET", "PUT", "DELETE", "OPTIONS"}},
+	} {
+		for _, method := range tc.methods {
+			rec := r.serve(method, tc.target, apLoopback, nil, "Origin", "http://localhost:7644", "Access-Control-Request-Method", "POST")
+			if rec.Code != 405 {
+				t.Errorf("%s %s = %d, want 405", method, tc.target, rec.Code)
+			}
+			for k := range rec.Header() {
+				if strings.HasPrefix(strings.ToLower(k), "access-control-") {
+					t.Errorf("%s %s wrote %s", method, tc.target, k)
+				}
+			}
+		}
+	}
+	for _, target := range []string{"/surface/decisions", "/surface/decisions/", "/surface/decisions/" + ref + "/x", "/surface/decide", "/surface/decide/" + ref + "/x"} {
+		if rec := r.serve("GET", target, apLoopback, nil); rec.Code != 404 && rec.Code != 405 {
+			t.Errorf("GET %s = %d, want 404 or 405", target, rec.Code)
+		}
+	}
+}
diff --git a/host/surfacepage/decide_model.go b/host/surfacepage/decide_model.go
new file mode 100644
index 0000000..f6d2135
--- /dev/null
+++ b/host/surfacepage/decide_model.go
@@ -0,0 +1,92 @@
+package surfacepage
+
+// Row 201 A4 sprint 1 (design_docs/planned/w-world-approvals-by-key-a4.md §2.8): the decision
+// routes' JSON shapes and the POST body's grammar, pure. Sprint 1 decides with the typed arm only;
+// the key formats (sshsig, webauthn) arrive with sprint 2, so until then they are unknown fields.
+
+import (
+	"bytes"
+	"encoding/json"
+	"errors"
+	"fmt"
+	"io"
+)
+
+// DecisionSubject is GET /surface/decisions/{ref}: the stored bytes of the request and of the one
+// packet carrying it (never a precomputed message or hash), the window an expiry must fall in, and
+// the law's display verdict at now (the POST runs the law again).
+type DecisionSubject struct {
+	V            int            `json:"v"`
+	Request      string         `json:"request"`
+	RequestBytes string         `json:"requestBytes"` // base64url, no padding
+	Packet       string         `json:"packet"`
+	PacketBytes  string         `json:"packetBytes"` // base64url, no padding
+	Entry        int64          `json:"entry"`
+	Proposal     string         `json:"proposal"`
+	Now          int64          `json:"now"`
+	Window       DecisionWindow `json:"window"`
+	Refusal      string         `json:"refusal"`
+	Reserved     string         `json:"reserved"`
+	DecidedBy    string         `json:"decidedBy"`
+	Decision     string         `json:"decision"`
+	SSH          DecisionSSH    `json:"ssh"`
+}
+
+// DecisionWindow is the expiry window of a decision made now: min = now + 1, max = min(now + 600, deadline).
+type DecisionWindow struct {
+	Min int64 `json:"min"`
+	Max int64 `json:"max"`
+}
+
+// DecisionSSH names the SSHSIG namespace a decision is signed under.
+type DecisionSSH struct {
+	Namespace string `json:"namespace"`
+}
+
+// DecisionMade is the POST's success: the verb, the request, the decision object and who decided.
+type DecisionMade struct {
+	V         int    `json:"v"`
+	Decided   string `json:"decided"`
+	Request   string `json:"request"`
+	Decision  string `json:"decision"`
+	DecidedBy string `json:"decidedBy"`
+}
+
+// DecideBody is the POST body (design §2.8), sprint 1's typed arm: the verb, the packet digest P
+// and request digest Q the caller saw (64 lowercase hex), the expiry, and typed: true.
+type DecideBody struct {
+	V       int    `json:"v"`
+	Verb    string `json:"verb"`
+	Packet  string `json:"packet"`
+	Request string `json:"request"`
+	Expiry  int64  `json:"expiry"`
+	Typed   bool   `json:"typed"`
+}
+
+// ErrDecideBody refuses a body that is not DecideBody's grammar.
+var ErrDecideBody = errors.New(`the body is not {"v":1,"verb":…,"packet":<64 hex>,"request":<64 hex>,"expiry":…,"typed":true}`)
+
+// ParseDecideBody is the POST's grammar: strict JSON (no unknown field, so no decidedBy and no key
+// format yet; exactly one value), v 1, typed true, both digests 64 lowercase hex. The verb and the
+// expiry are the law's (D9, D10), not the grammar's.
+func ParseDecideBody(b []byte) (body DecideBody, packet, request [32]byte, err error) {
+	dec := json.NewDecoder(bytes.NewReader(b))
+	dec.DisallowUnknownFields()
+	if err := dec.Decode(&body); err != nil {
+		return DecideBody{}, packet, request, fmt.Errorf("%w: %v", ErrDecideBody, err)
+	}
+	if _, err := dec.Token(); err != io.EOF {
+		return DecideBody{}, packet, request, fmt.Errorf("%w: more than one JSON value", ErrDecideBody)
+	}
+	switch {
+	case body.V != 1:
+		return DecideBody{}, packet, request, fmt.Errorf("%w: v is %d", ErrDecideBody, body.V)
+	case !body.Typed:
+		return DecideBody{}, packet, request, fmt.Errorf("%w: it has no format (this build decides with typed only)", ErrDecideBody)
+	case !lowerHex64(body.Packet) || !lowerHex64(body.Request):
+		return DecideBody{}, packet, request, fmt.Errorf("%w: packet and request are each 64 lowercase hex", ErrDecideBody)
+	}
+	packet, _ = hex32(body.Packet)
+	request, _ = hex32(body.Request)
+	return body, packet, request, nil
+}
````

### Patch M3b — the rig and the subject (RD3)

````diff M3b
diff --git a/host/daemon/surface_decide_rig_test.go b/host/daemon/surface_decide_rig_test.go
new file mode 100644
index 0000000..f097c34
--- /dev/null
+++ b/host/daemon/surface_decide_rig_test.go
@@ -0,0 +1,164 @@
+package daemon
+
+// Row 201 A4 sprint 1: the decision routes' test rig. Every store is under t.TempDir(); every
+// packet is a forged, validated commit (workbench_packet_test.go's forgedChain); every operator
+// proof is written as the CLI's writeOperatorProof writes it (D-WORLD-109: no real store or tty).
+
+import (
+	"crypto/sha256"
+	"fmt"
+	"net/http/httptest"
+	"os"
+	"path/filepath"
+	"sort"
+	"testing"
+	"time"
+
+	"github.com/sunholo-data/ailang-world/host/broker"
+	"github.com/sunholo-data/ailang-world/host/coordinator"
+	"github.com/sunholo-data/ailang-world/host/hashref"
+	"github.com/sunholo-data/ailang-world/host/store"
+)
+
+// decNow is the decision routes' pinned clock: after the packets' createdAt (packetNow), before
+// their deadline (packetNow + 60).
+const decNow = packetNow + 10
+
+type decideRig struct{ approvalRig }
+
+func newDecideRig(t *testing.T) decideRig {
+	t.Helper()
+	a := newApprovalRig(t)
+	a.d.approveNow = func() int64 { return decNow }
+	return decideRig{a}
+}
+
+// packet files a pending Human.Approve request (scope, requester) in the approvals chain and
+// commits a validated packet carrying it, whose proposal is label's.
+func (r decideRig) packet(label, scope, requester string) (hashref.HashRef, store.Object) {
+	r.t.Helper()
+	q := approvalRequestInChain(r.t, r.d, scope, requester, packetNow)
+	p, _ := commitForgedChain(r.t, r.d, newForgedChain(label, q, scope))
+	return q, p
+}
+
+// proof writes the operator proof for a decision of q with phase, beside the credentials store,
+// as writeOperatorProof does (0700 directory, 0600 file, expiry within 300 s); it returns the secret.
+func (r decideRig) proof(q hashref.HashRef, phase string) string {
+	r.t.Helper()
+	secret := hashref.SumSHA256([]byte(phase + q.String())).Digest()
+	dir := filepath.Join(filepath.Dir(r.d.credentialsPath), "operator-proofs")
+	if err := os.MkdirAll(dir, 0o700); err != nil {
+		r.t.Fatal(err)
+	}
+	data := mustMarshal(r.t, operatorProof{SHA256: hashref.SumSHA256([]byte(secret)).Digest(), Phase: phase, ExpiresAt: time.Now().Unix() + 200})
+	if err := os.WriteFile(r.proofPath(q), data, 0o600); err != nil {
+		r.t.Fatal(err)
+	}
+	return secret
+}
+
+func (r decideRig) proofPath(q hashref.HashRef) string {
+	return filepath.Join(filepath.Dir(r.d.credentialsPath), "operator-proofs", "decision-"+q.Digest()+".json")
+}
+
+// proofExists reports whether q's proof file is still there (nothing consumed it).
+func (r decideRig) proofExists(q hashref.HashRef) bool {
+	_, err := os.Lstat(r.proofPath(q))
+	return err == nil
+}
+
+// decideBody is a typed POST body for (verb, p, q, expiry).
+func decideBody(verb string, p store.Object, q hashref.HashRef, expiry int64) []byte {
+	return []byte(fmt.Sprintf(`{"v":1,"verb":%q,"packet":%q,"request":%q,"expiry":%d,"typed":true}`, verb, p.Hash.Digest(), q.Digest(), expiry))
+}
+
+// post is POST /surface/decide/<ref> from the loopback with A3a's headers and the proof; hdr adds
+// or replaces headers.
+func (r decideRig) post(ref string, body []byte, proof string, hdr ...string) *httptest.ResponseRecorder {
+	r.t.Helper()
+	h := []string{"Origin", "http://127.0.0.1:7644", "Content-Type", "application/json"}
+	if proof != "" {
+		h = append(h, "X-World-Operator-Proof", proof)
+	}
+	return r.serve("POST", "/surface/decide/"+ref, apLoopback, body, append(h, hdr...)...)
+}
+
+// decisions is the number of stored approval decisions and the approvals head.
+func (r decideRig) decisions() (int, hashref.HashRef) {
+	r.t.Helper()
+	ctx := boundedTestContext(r.t)
+	objs, err := r.d.store.ObjectsBySemanticID(ctx, broker.ApprovalDecisionV1, "", store.MaxSemanticIDPage)
+	if err != nil {
+		r.t.Fatal(err)
+	}
+	head, _, err := r.d.store.GetRegistryHead(ctx, broker.ApprovalsV1)
+	if err != nil {
+		r.t.Fatal(err)
+	}
+	return len(objs), head
+}
+
+// fillPackets stores n decodable decision packets that carry other requests (no commit: the
+// discovery reads every stored packet) and returns their hashes, sorted.
+func (r decideRig) fillPackets(label string, n int) []string {
+	r.t.Helper()
+	ctx := boundedTestContext(r.t)
+	var hashes []string
+	for i := 0; i < n; i++ {
+		payload := coordinator.EncodeDecisionPacket(coordinator.DecisionPacket{ProposalHash: hashref.SumSHA256([]byte(fmt.Sprintf("%s-plan-%d", label, i))),
+			RequestRef: hashref.SumSHA256([]byte(fmt.Sprintf("%s-request-%d", label, i))), CreatedAt: packetNow, DeadlineAt: packetNow + 60, Policy: "Cancel"})
+		o := forgeObject(coordinator.DecisionPacketV1, payload)
+		if err := r.d.store.PutObject(ctx, o); err != nil {
+			r.t.Fatal(err)
+		}
+		hashes = append(hashes, o.Hash.String())
+	}
+	sort.Strings(hashes)
+	return hashes
+}
+
+// packetAfter commits a validated packet carrying q whose hash sorts after position pos among
+// fillers (sorted), trying proposal labels until one does; it returns the packet.
+func (r decideRig) packetAfter(label string, q hashref.HashRef, scope string, fillers []string, pos int) store.Object {
+	r.t.Helper()
+	for i := 0; i < 200; i++ {
+		f := newForgedChain(fmt.Sprintf("%s-%d", label, i), q, scope)
+		h := hashref.SumSHA256(coordinator.EncodeDecisionPacket(f.packet)).String()
+		if sort.SearchStrings(fillers, h) >= pos {
+			p, _ := commitForgedChain(r.t, r.d, f)
+			return p
+		}
+	}
+	r.t.Fatal("no proposal label sorts the packet late enough")
+	return store.Object{}
+}
+
+// packetPosition is p's 0-based position among every stored packet in hash_ref order.
+func (r decideRig) packetPosition(p store.Object) int {
+	r.t.Helper()
+	ctx := boundedTestContext(r.t)
+	pos, after := 0, ""
+	for {
+		page, err := r.d.store.ObjectsBySemanticID(ctx, coordinator.DecisionPacketV1, after, store.MaxSemanticIDPage)
+		if err != nil {
+			r.t.Fatal(err)
+		}
+		for _, o := range page {
+			if o.Hash == p.Hash {
+				return pos
+			}
+			pos++
+		}
+		if len(page) < store.MaxSemanticIDPage {
+			r.t.Fatalf("packet %s is not stored", p.Hash)
+		}
+		after = page[len(page)-1].Hash.String()
+	}
+}
+
+// sha256Of is the ref of payload.
+func sha256Of(payload []byte) string {
+	sum := sha256.Sum256(payload)
+	return "sha256:" + fmt.Sprintf("%x", sum[:])
+}
diff --git a/host/daemon/surface_decide_subject_test.go b/host/daemon/surface_decide_subject_test.go
new file mode 100644
index 0000000..43e731e
--- /dev/null
+++ b/host/daemon/surface_decide_subject_test.go
@@ -0,0 +1,108 @@
+package daemon
+
+// Row 201 A4 sprint 1 (design §2.8, §6 RD3): GET /surface/decisions/{ref} is the subject as bytes.
+
+import (
+	"encoding/base64"
+	"encoding/json"
+	"reflect"
+	"sort"
+	"strings"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/broker"
+	"github.com/sunholo-data/ailang-world/host/coordinator"
+	"github.com/sunholo-data/ailang-world/host/surfacepage"
+)
+
+// RD3: the bytes decode to exactly the stored payloads (the client re-hashes them); no message or
+// hash field; the window; refusal/reserved/decision set by the law; D1-D4 answer the refusal.
+func TestDecisionSubjectJSON(t *testing.T) {
+	r := newDecideRig(t)
+	q, p := r.packet("subject", broker.WorkspaceScope, "agent")
+	rec := r.serve("GET", "/surface/decisions/"+q.String(), apLoopback, nil)
+	if rec.Code != 200 {
+		t.Fatalf("subject = %d %s", rec.Code, rec.Body.String())
+	}
+	if err := approvalHeaders(rec); err != nil {
+		t.Fatal(err)
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
+	if want := []string{"decidedBy", "decision", "entry", "now", "packet", "packetBytes", "proposal", "refusal", "request", "requestBytes", "reserved", "ssh", "v", "window"}; !reflect.DeepEqual(names, want) {
+		t.Fatalf("subject keys %v, want %v (no message, no hash)", names, want)
+	}
+	var s surfacepage.DecisionSubject
+	if err := json.Unmarshal(rec.Body.Bytes(), &s); err != nil {
+		t.Fatal(err)
+	}
+	reqBytes, err1 := base64.RawURLEncoding.DecodeString(s.RequestBytes)
+	pktBytes, err2 := base64.RawURLEncoding.DecodeString(s.PacketBytes)
+	if err1 != nil || err2 != nil || sha256Of(reqBytes) != q.String() || sha256Of(pktBytes) != p.Hash.String() || s.Request != q.String() || s.Packet != p.Hash.String() {
+		t.Fatalf("subject bytes do not hash to their refs: %+v", s)
+	}
+	stored, ok, err := r.d.store.GetObject(boundedTestContext(t), q)
+	if err != nil || !ok || string(stored.Payload) != string(reqBytes) || string(p.Payload) != string(pktBytes) {
+		t.Fatal("subject bytes are not the stored payloads")
+	}
+	packet, err := coordinator.DecodeDecisionPacket(pktBytes)
+	if err != nil || packet.RequestRef != q || s.Proposal != packet.ProposalHash.String() || s.Entry < 1 {
+		t.Fatalf("packet %+v, %v; subject %+v", packet, err, s)
+	}
+	if s.V != 1 || s.Now != decNow || s.Window != (surfacepage.DecisionWindow{Min: decNow + 1, Max: packetNow + 60}) ||
+		s.Refusal != "" || s.Reserved != "" || s.Decision != "" || s.DecidedBy != "" || s.SSH.Namespace != "world-decision@sunholo" {
+		t.Fatalf("subject %+v", s)
+	}
+	// A wide deadline: window.max is now + 600.
+	r.d.approveNow = func() int64 { return packetNow - 600 }
+	if err := json.Unmarshal(r.serve("GET", "/surface/decisions/"+q.String(), apLoopback, nil).Body.Bytes(), &s); err != nil || s.Window.Max != packetNow {
+		t.Fatalf("window %+v, %v", s.Window, err)
+	}
+	r.d.approveNow = func() int64 { return decNow }
+	// After a decision: the decision, and the law's D6 refusal as display.
+	if _, err := broker.DecideApproval(boundedTestContext(t), r.d.store, q, "deny", "operator:terminal", decNow); err != nil {
+		t.Fatal(err)
+	}
+	s = surfacepage.DecisionSubject{}
+	if err := json.Unmarshal(r.serve("GET", "/surface/decisions/"+q.String(), apLoopback, nil).Body.Bytes(), &s); err != nil ||
+		s.Decision != "deny" || s.DecidedBy != "operator:terminal" || !strings.Contains(s.Refusal, "was already decided deny by operator:terminal") {
+		t.Fatalf("decided subject %+v, %v", s, err)
+	}
+	// Reserved: a finish-shaped request shows its owner and the D5 refusal.
+	fq, _ := r.packet("subject-finish", "github-pr:o/r", "agent")
+	s = surfacepage.DecisionSubject{}
+	if err := json.Unmarshal(r.serve("GET", "/surface/decisions/"+fq.String(), apLoopback, nil).Body.Bytes(), &s); err != nil ||
+		s.Reserved != "finish" || !strings.Contains(s.Refusal, "ailang-worldd episode finish") {
+		t.Fatalf("reserved subject %+v, %v", s, err)
+	}
+	// D7 at the deadline is display too.
+	r.d.approveNow = func() int64 { return packetNow + 60 }
+	s = surfacepage.DecisionSubject{}
+	if err := json.Unmarshal(r.serve("GET", "/surface/decisions/"+fq.String(), apLoopback, nil).Body.Bytes(), &s); err != nil || !strings.Contains(s.Refusal, "ailang-worldd episode finish") {
+		t.Fatalf("reserved before the deadline: %+v", s)
+	}
+	r.d.approveNow = func() int64 { return decNow }
+	// D1-D4: no subject, the law's refusal.
+	pending := approvalRequestInChain(t, r.d, broker.WorkspaceScope, "no-packet", packetNow)
+	for _, tc := range []struct {
+		ref, class, reason string
+		status             int
+	}{
+		{"sha256:XYZ", "BadRequest", "is not a request ref", 400},
+		{"sha256:" + strings.Repeat("0", 64), "NotFound", "no such approval request", 404},
+		{p.Hash.String(), "NotFound", "no such approval request", 404},
+		{pending.String(), "Conflict", "no committed decision packet carries request", 409},
+	} {
+		got := requireRefusal(t, tc.ref, r.serve("GET", "/surface/decisions/"+tc.ref, apLoopback, nil), tc.status, tc.class)
+		if !strings.Contains(got.Reason, tc.reason) {
+			t.Errorf("%s: %q lacks %q", tc.ref, got.Reason, tc.reason)
+		}
+	}
+}
````

### Patch M3c — the POST's guards, law arms, proof, replay and R4 (RD2, RD4–RD8, R4d)

````diff M3c
diff --git a/host/daemon/surface_decide_post_test.go b/host/daemon/surface_decide_post_test.go
new file mode 100644
index 0000000..928d4ab
--- /dev/null
+++ b/host/daemon/surface_decide_post_test.go
@@ -0,0 +1,273 @@
+package daemon
+
+// Row 201 A4 sprint 1 (design §2.4, §2.5, §2.8, §6 RD2, RD4-RD8, R4d): POST /surface/decide/{ref}
+// with the typed arm: the guards, the law's arms, the proof, replay after a restart, R4 by L1.
+
+import (
+	"encoding/json"
+	"io"
+	"os"
+	"path/filepath"
+	"strings"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/broker"
+	"github.com/sunholo-data/ailang-world/host/hashref"
+	"github.com/sunholo-data/ailang-world/host/store"
+	"github.com/sunholo-data/ailang-world/host/surfacepage"
+)
+
+// requireNothingDecided fails unless no decision was stored, the head did not move and the proof
+// of q is still unconsumed.
+func (r decideRig) requireNothingDecided(what string, q hashref.HashRef, head hashref.HashRef) {
+	r.t.Helper()
+	if n, h := r.decisions(); n != 0 || h != head {
+		r.t.Fatalf("%s: %d decisions, head %v (was %v)", what, n, h, head)
+	}
+	if !r.proofExists(q) {
+		r.t.Fatalf("%s: the proof was consumed", what)
+	}
+}
+
+// RD2: R0 Authorization, R0 Cookie, R2 a foreign Host, R3 no/null/other Origin, R4 same-site, R5
+// text/plain, R6 an unknown field, decidedBy, a key format, typed false, v 2, a bad digest, two
+// values, too large: each refused before the law, 0 decisions, the proof unconsumed; then the
+// same proof decides (the control).
+func TestDecidePOSTGuards(t *testing.T) {
+	r := newDecideRig(t)
+	q, p := r.packet("guards", broker.WorkspaceScope, "agent")
+	proof := r.proof(q, "decide-approve")
+	good := decideBody("approve", p, q, decNow+50)
+	_, head := r.decisions()
+	field := func(extra string) []byte { return []byte(strings.Replace(string(good), `"typed":true`, extra, 1)) }
+	for _, tc := range []struct {
+		name   string
+		host   string
+		body   []byte
+		hdr    []string
+		status int
+		class  string
+	}{
+		{"R0 Authorization", apLoopback, good, []string{"Authorization", "Bearer agent-session"}, 403, "Forbidden"},
+		{"R0 empty Authorization", apLoopback, good, []string{"Authorization", ""}, 403, "Forbidden"},
+		{"R0 Cookie", apLoopback, good, []string{"Cookie", "s=1"}, 403, "Forbidden"},
+		{"R2 foreign Host", "rebind.example:7644", good, nil, 421, "MisdirectedRequest"},
+		{"R3 no Origin", apLoopback, good, []string{"Origin", ""}, 403, "CrossOrigin"},
+		{"R3 null Origin", apLoopback, good, []string{"Origin", "null"}, 403, "CrossOrigin"},
+		{"R3 other Origin", apLoopback, good, []string{"Origin", "http://localhost:7645"}, 403, "CrossOrigin"},
+		{"R4 same-site", apLoopback, good, []string{"Sec-Fetch-Site", "same-site"}, 403, "CrossOrigin"},
+		{"R5 text/plain", apLoopback, good, []string{"Content-Type", "text/plain"}, 415, "UnsupportedMediaType"},
+		{"R6 unknown field", apLoopback, field(`"typed":true,"x":1`), nil, 400, "BadRequest"},
+		{"R6 decidedBy", apLoopback, field(`"typed":true,"decidedBy":"ssh:mark"`), nil, 400, "BadRequest"},
+		{"R6 a key format", apLoopback, field(`"typed":true,"sshsig":"-----BEGIN SSH SIGNATURE-----"`), nil, 400, "BadRequest"},
+		{"R6 typed false", apLoopback, field(`"typed":false`), nil, 400, "BadRequest"},
+		{"R6 no format", apLoopback, []byte(strings.Replace(string(good), `,"typed":true`, "", 1)), nil, 400, "BadRequest"},
+		{"R6 v 2", apLoopback, []byte(strings.Replace(string(good), `"v":1`, `"v":2`, 1)), nil, 400, "BadRequest"},
+		{"R6 upper-case packet", apLoopback, []byte(strings.Replace(string(good), p.Hash.Digest(), strings.ToUpper(p.Hash.Digest()), 1)), nil, 400, "BadRequest"},
+		{"R6 short request", apLoopback, []byte(strings.Replace(string(good), q.Digest(), q.Digest()[:62], 1)), nil, 400, "BadRequest"},
+		{"R6 two values", apLoopback, append(append([]byte{}, good...), good...), nil, 400, "BadRequest"},
+		{"R6 too large", apLoopback, []byte(`{"v":1,"verb":"` + strings.Repeat("a", surfacepage.MaxApproveBodyBytes) + `"}`), nil, 413, "TooLarge"},
+	} {
+		hdr := append([]string{"Origin", "http://127.0.0.1:7644", "Content-Type", "application/json", "X-World-Operator-Proof", proof}, tc.hdr...)
+		rec := r.serve("POST", "/surface/decide/"+q.String(), tc.host, tc.body, hdr...)
+		requireRefusal(t, tc.name, rec, tc.status, tc.class)
+		r.requireNothingDecided(tc.name, q, head)
+	}
+	if rec := r.post(q.String(), good, proof); rec.Code != 200 {
+		t.Fatalf("control: %d %s", rec.Code, rec.Body.String())
+	}
+}
+
+// RD4: a store built per D1-D10 → the law's status and class, 0 decisions, the proof unconsumed
+// (the law runs before the arm).
+func TestDecidePOSTLawArms(t *testing.T) {
+	r := newDecideRig(t)
+	q, p := r.packet("law", broker.WorkspaceScope, "agent")
+	proof := r.proof(q, "decide-approve")
+	_, head := r.decisions()
+	ctx := boundedTestContext(t)
+	offChain := forgeObject(broker.ApprovalRequestV1, []byte(`{"effect":"Human.Approve","scope":"worktree","cost":1,"requester":"off","now":1700}`))
+	notApprove := forgeObject(broker.ApprovalRequestV1, []byte(`{"effect":"Human.PollApproval","scope":"worktree","cost":1,"requester":"x","now":1700}`))
+	for _, o := range []store.Object{offChain, notApprove} {
+		if err := r.d.store.PutObject(ctx, o); err != nil {
+			t.Fatal(err)
+		}
+	}
+	pending := approvalRequestInChain(t, r.d, broker.WorkspaceScope, "no-packet", packetNow)
+	failed := approvalRequestInChain(t, r.d, broker.WorkspaceScope, "failed", packetNow)
+	fc := newForgedChain("law-failed", failed, broker.WorkspaceScope)
+	fc.failed = true
+	fp, _ := commitForgedChain(t, r.d, fc)
+	_, head = r.decisions()
+	for _, tc := range []struct {
+		name, ref string
+		body      []byte
+		now       int64
+		status    int
+		class     string
+		reason    string
+	}{
+		{"D1", "sha256:XYZ", decideBody("approve", p, q, decNow+50), decNow, 400, "BadRequest", "is not a request ref"},
+		{"D2 absent", "sha256:" + strings.Repeat("0", 64), decideBody("approve", p, q, decNow+50), decNow, 404, "NotFound", "no such approval request"},
+		{"D2 a packet", p.Hash.String(), decideBody("approve", p, q, decNow+50), decNow, 404, "NotFound", "no such approval request"},
+		{"D2 another effect", notApprove.Hash.String(), decideBody("approve", p, q, decNow+50), decNow, 404, "NotFound", "no such approval request"},
+		{"D3", offChain.Hash.String(), decideBody("approve", p, q, decNow+50), decNow, 404, "NotFound", "not on the approvals chain"},
+		{"D4b zero", pending.String(), decideBody("approve", p, pending, decNow+50), decNow, 409, "Conflict", "no committed decision packet carries request"},
+		{"D4b refused packet", failed.String(), decideBody("approve", fp, failed, decNow+50), decNow, 409, "Conflict", "is refused: not a validated decision packet: requestRef"},
+		{"D7", q.String(), decideBody("approve", p, q, packetNow+61), packetNow + 60, 409, "Conflict", "deadline passed"},
+		{"D8 packet", q.String(), decideBody("approve", fp, q, decNow+50), decNow, 409, "Conflict", "reload it"},
+		{"D8 request", q.String(), []byte(strings.Replace(string(decideBody("approve", p, q, decNow+50)), q.Digest(), pending.Digest(), 1)), decNow, 409, "Conflict", "reload it"},
+		{"D9", q.String(), decideBody("maybe", p, q, decNow+50), decNow, 400, "BadRequest", "a decision is approve or deny"},
+		{"D10 stale", q.String(), decideBody("approve", p, q, decNow), decNow, 409, "Conflict", "at most 10 minutes"},
+		{"D10 past the deadline", q.String(), decideBody("approve", p, q, packetNow+61), decNow, 409, "Conflict", "at most 10 minutes"},
+	} {
+		r.d.approveNow = func() int64 { return tc.now }
+		got := requireRefusal(t, tc.name, r.post(tc.ref, tc.body, proof), tc.status, tc.class)
+		if !strings.Contains(got.Reason, tc.reason) {
+			t.Errorf("%s: %q lacks %q", tc.name, got.Reason, tc.reason)
+		}
+		r.requireNothingDecided(tc.name, q, head)
+	}
+	// D6 is RD8's; D5 is R4d's; D4a and the counts are RD9*'s.
+}
+
+// RD5: a typed proof decides: one decision by operator:terminal, the proof unlinked; deny too.
+func TestDecideTypedLands(t *testing.T) {
+	r := newDecideRig(t)
+	for _, verb := range []string{"approve", "deny"} {
+		q, p := r.packet("typed-"+verb, broker.WorkspaceScope, "agent-"+verb)
+		before, _ := r.decisions()
+		rec := r.post(q.String(), decideBody(verb, p, q, decNow+50), r.proof(q, "decide-"+verb))
+		var got surfacepage.DecisionMade
+		if err := json.Unmarshal(rec.Body.Bytes(), &got); err != nil || rec.Code != 200 || got.V != 1 || got.Decided != verb || got.Request != q.String() || got.DecidedBy != "operator:terminal" {
+			t.Fatalf("%s: %d %s", verb, rec.Code, rec.Body.String())
+		}
+		if err := approvalHeaders(rec); err != nil {
+			t.Fatal(err)
+		}
+		on, decided, v, by, err := broker.ApprovalChainState(boundedTestContext(t), r.d.store, q)
+		if err != nil || !on || !decided || v != verb || by != "operator:terminal" {
+			t.Fatalf("%s: chain %v %v %q %q %v", verb, on, decided, v, by, err)
+		}
+		obj, ok, err := r.d.store.GetObject(boundedTestContext(t), hashref.MustParse(got.Decision))
+		if err != nil || !ok || obj.SemanticID != broker.ApprovalDecisionV1 || !strings.Contains(string(obj.Payload), `"now":1710`) {
+			t.Fatalf("%s: decision object %s %v %v", verb, obj.Payload, ok, err)
+		}
+		if after, _ := r.decisions(); after != before+1 || r.proofExists(q) {
+			t.Fatalf("%s: decisions %d -> %d, proof left %v", verb, before, after, r.proofExists(q))
+		}
+	}
+}
+
+// RD6: the proof binds the verb and the request: a decide-deny proof cannot approve; another
+// request's proof, a wrong secret and no proof are refused; 0 decisions.
+func TestDecideTypedProofBindsTheVerb(t *testing.T) {
+	r := newDecideRig(t)
+	q, p := r.packet("verb", broker.WorkspaceScope, "agent")
+	other, _ := r.packet("verb-other", broker.WorkspaceScope, "other")
+	_, head := r.decisions()
+	body := decideBody("approve", p, q, decNow+50)
+	for _, tc := range []struct{ name, phase, proof string }{
+		{"deny proof", "decide-deny", ""}, {"finish phase", "approve", ""}, {"wrong secret", "decide-approve", strings.Repeat("e", 64)}, {"no header", "decide-approve", "-"},
+	} {
+		secret := r.proof(q, tc.phase)
+		switch tc.proof {
+		case "":
+		case "-":
+			secret = ""
+		default:
+			secret = tc.proof
+		}
+		got := requireRefusal(t, tc.name, r.post(q.String(), body, secret), 403, "Forbidden")
+		if !strings.Contains(got.Reason, "the typed decision is not proven: operator proof") {
+			t.Errorf("%s: %q", tc.name, got.Reason)
+		}
+		if n, h := r.decisions(); n != 0 || h != head {
+			t.Fatalf("%s: %d decisions", tc.name, n)
+		}
+	}
+	// Another request's proof is not this request's.
+	_ = os.Remove(r.proofPath(q))
+	got := requireRefusal(t, "other request", r.post(q.String(), body, r.proof(other, "decide-approve")), 403, "Forbidden")
+	if !strings.Contains(got.Reason, "operator proof is missing") || !r.proofExists(other) {
+		t.Errorf("other request: %q", got.Reason)
+	}
+}
+
+// RD7: P or Q of another packet → 409 D8.
+func TestDecidePOSTStaleDigests(t *testing.T) {
+	r := newDecideRig(t)
+	q, p := r.packet("stale", broker.WorkspaceScope, "agent")
+	q2, p2 := r.packet("stale-2", broker.WorkspaceScope, "agent-2")
+	proof := r.proof(q, "decide-approve")
+	_, head := r.decisions()
+	for name, body := range map[string][]byte{"P": decideBody("approve", p2, q, decNow+50), "Q": decideBody("approve", p, q2, decNow+50)} {
+		got := requireRefusal(t, name, r.post(q.String(), body, proof), 409, "Conflict")
+		if !strings.Contains(got.Reason, "the request changed since it was shown; reload it") {
+			t.Errorf("%s: %q", name, got.Reason)
+		}
+		r.requireNothingDecided(name, q, head)
+	}
+}
+
+// RD8: decide, close the daemon, open a new one on the same stores, POST again with a fresh proof
+// → 409 D6, still one decision.
+func TestDecideReplayAfterRestart(t *testing.T) {
+	r := newDecideRig(t)
+	q, p := r.packet("replay", broker.WorkspaceScope, "agent")
+	body := decideBody("deny", p, q, decNow+50)
+	if rec := r.post(q.String(), body, r.proof(q, "decide-deny")); rec.Code != 200 {
+		t.Fatalf("first: %d %s", rec.Code, rec.Body.String())
+	}
+	n, head := r.decisions()
+	db := r.d.cfg.DBPath
+	if err := r.d.Close(); err != nil {
+		t.Fatal(err)
+	}
+	d2, err := New(boundedTestContext(t), Config{DBPath: db, BindHost: DefaultBindHost, ErrorLog: io.Discard})
+	if err != nil {
+		t.Fatal(err)
+	}
+	t.Cleanup(func() { _ = d2.Close() })
+	d2.approveNow = func() int64 { return decNow + 1 }
+	r2 := decideRig{approvalRig{t: t, d: d2}}
+	if filepath.Dir(d2.credentialsPath) != filepath.Dir(r.d.credentialsPath) {
+		t.Fatal("the restarted daemon has another credentials directory")
+	}
+	for _, verb := range []string{"deny", "approve"} {
+		got := requireRefusal(t, "replay "+verb, r2.post(q.String(), decideBody(verb, p, q, decNow+50), r2.proof(q, "decide-"+verb)), 409, "Conflict")
+		if !strings.Contains(got.Reason, "was already decided deny by operator:terminal") {
+			t.Errorf("replay %s: %q", verb, got.Reason)
+		}
+		r2.requireNothingDecidedSince(n, head)
+	}
+}
+
+func (r decideRig) requireNothingDecidedSince(n int, head hashref.HashRef) {
+	r.t.Helper()
+	if m, h := r.decisions(); m != n || h != head {
+		r.t.Fatalf("%d decisions (was %d), head %v (was %v)", m, n, h, head)
+	}
+}
+
+// R4d: publish- and finish-shaped requests are refused for the typed arm too (L1, D5), with the
+// owning command; 0 decisions, the proof unconsumed.
+func TestDecideRefusesReservedScopesTyped(t *testing.T) {
+	r := newDecideRig(t)
+	publish := broker.PublishApprovalScopeFor(broker.PublishIdentity{RegistryOrigin: "https://registry.example", Vendor: "v", Name: "n", Version: "1.0.0",
+		ManifestRef: hashref.SumSHA256([]byte("m"))}, broker.PublishHashes{TarballSHA256: "t", ContentHash: "c", InterfaceHash: "i"}, packetNow+60)
+	finish := broker.FormatFinishApprovalScope(broker.FinishApprovalScope{Repository: "o/r", Remote: "https://github.com/o/r.git", Head: "world/e", SHA: strings.Repeat("a", 40),
+		Base: "dev", TitleSHA256: strings.Repeat("b", 64), BodySHA256: strings.Repeat("c", 64), ExpiresAt: packetNow + 60})
+	for _, tc := range []struct{ name, scope, owner string }{
+		{"publish", publish, "`world-publish approve` at a terminal (R4)"}, {"publish near miss", "registry:x", "`world-publish approve`"},
+		{"finish", finish, "`ailang-worldd episode finish`"}, {"finish near miss", "github-pr:x", "`ailang-worldd episode finish`"},
+	} {
+		q, p := r.packet("reserved-"+tc.name, tc.scope, "agent")
+		_, head := r.decisions()
+		got := requireRefusal(t, tc.name, r.post(q.String(), decideBody("approve", p, q, decNow+50), r.proof(q, "decide-approve")), 403, "Forbidden")
+		if !strings.Contains(got.Reason, tc.owner) {
+			t.Errorf("%s: %q lacks %q", tc.name, got.Reason, tc.owner)
+		}
+		r.requireNothingDecided(tc.name, q, head)
+	}
+}
````

### Patch M3d — the discovery end to end (RD9, RD9b–RD9d)

````diff M3d
diff --git a/host/daemon/surface_decide_discovery_test.go b/host/daemon/surface_decide_discovery_test.go
new file mode 100644
index 0000000..cdccc0c
--- /dev/null
+++ b/host/daemon/surface_decide_discovery_test.go
@@ -0,0 +1,84 @@
+package daemon
+
+// Row 201 A4 sprint 1 (design §2.4 [r1], §6 RD9, RD9b-RD9d): the POST counts the packets carrying
+// a request over every stored packet, never over ValidatePacket's 100-packet prefix.
+
+import (
+	"strings"
+	"testing"
+	"time"
+
+	"github.com/sunholo-data/ailang-world/host/broker"
+	"github.com/sunholo-data/ailang-world/host/coordinator"
+)
+
+// RD9: two invocations whose requests collide (one ref) and two packets → 409 D4b, 0 decisions.
+func TestDecideRefusesASharedRequest(t *testing.T) {
+	r := newDecideRig(t)
+	q, p := r.packet("shared-1", broker.WorkspaceScope, "agent")
+	q2, _ := r.packet("shared-2", broker.WorkspaceScope, "agent")
+	if q2 != q {
+		t.Fatalf("identical requests did not collide: %v %v", q, q2)
+	}
+	_, head := r.decisions()
+	got := requireRefusal(t, "shared", r.post(q.String(), decideBody("approve", p, q, decNow+50), r.proof(q, "decide-approve")), 409, "Conflict")
+	if !strings.Contains(got.Reason, "is carried by 2 decision packets; deciding it would decide all of them") {
+		t.Errorf("shared: %q", got.Reason)
+	}
+	r.requireNothingDecided("shared", q, head)
+}
+
+// RD9b: 150 packets, two carrying the request, the second sorting after the first 101 (where
+// ValidatePacket's duplicate scan stops; it validates the first with code 0) → 409 D4b, 0
+// decisions, the head unchanged.
+func TestDecideDiscoveryFindsAHiddenSecondPacket(t *testing.T) {
+	r := newDecideRig(t)
+	q, p := r.packet("hidden-1", broker.WorkspaceScope, "agent")
+	fillers := r.fillPackets("hidden", 147)
+	late := r.packetAfter("hidden-2", q, broker.WorkspaceScope, fillers, 102)
+	if pos := r.packetPosition(late); pos < 101 {
+		t.Fatalf("the second packet is at position %d, inside the first 101", pos)
+	}
+	if res, err := coordinator.ValidatePacket(boundedTestContext(t), r.d.store, p); err != nil || res.Code != 0 {
+		t.Fatalf("ValidatePacket of the first = %d, %v; want 0 (the scan cannot see the second)", res.Code, err)
+	}
+	_, head := r.decisions()
+	got := requireRefusal(t, "hidden", r.post(q.String(), decideBody("approve", p, q, decNow+50), r.proof(q, "decide-approve")), 409, "Conflict")
+	if !strings.Contains(got.Reason, "is carried by 2 decision packets") {
+		t.Errorf("hidden: %q", got.Reason)
+	}
+	r.requireNothingDecided("hidden", q, head)
+}
+
+// RD9c: 2,001 packets, exactly one carrying the request → 409 D4a (the overflow sentinel), 0
+// decisions.
+func TestDecideDiscoveryRefusesAnOverflow(t *testing.T) {
+	r := newDecideRig(t)
+	r.d.commitBudget = time.Minute // 2,001 payload reads in one POST: the budget is not what this test measures (a loaded rig exceeded 3 s)
+	q, p := r.packet("overflow", broker.WorkspaceScope, "agent")
+	r.fillPackets("overflow", broker.DecisionDiscoveryLimit)
+	_, head := r.decisions()
+	got := requireRefusal(t, "overflow", r.post(q.String(), decideBody("approve", p, q, decNow+50), r.proof(q, "decide-approve")), 409, "Conflict")
+	if !strings.Contains(got.Reason, "cannot be counted completely (more than 2000 stored packets") {
+		t.Errorf("overflow: %q", got.Reason)
+	}
+	r.requireNothingDecided("overflow", q, head)
+}
+
+// RD9d (must live): 150 packets, exactly one carrying the request, sorting after the first 101 →
+// the typed POST lands one decision.
+func TestDecideDiscoveryFindsTheOnlyPacketBeyondThePrefix(t *testing.T) {
+	r := newDecideRig(t)
+	fillers := r.fillPackets("beyond", 149)
+	q := approvalRequestInChain(t, r.d, broker.WorkspaceScope, "agent", packetNow)
+	p := r.packetAfter("beyond", q, broker.WorkspaceScope, fillers, 102)
+	if pos := r.packetPosition(p); pos < 101 {
+		t.Fatalf("the packet is at position %d, inside the first 101", pos)
+	}
+	if rec := r.post(q.String(), decideBody("deny", p, q, decNow+50), r.proof(q, "decide-deny")); rec.Code != 200 {
+		t.Fatalf("beyond the prefix: %d %s", rec.Code, rec.Body.String())
+	}
+	if n, _ := r.decisions(); n != 1 {
+		t.Fatalf("%d decisions, want 1", n)
+	}
+}
````
