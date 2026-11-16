# Sprint plan — w-world-address-bar-verbs, phase V1 "one grammar for the bar and the CLI" (iteration 395, row 243)

> **THE SPRINT — row 243 phase V1, the one V1 sprint (design §0.3, §2.1–§2.4, §6, §8, §9), M1–M4:** a new pure
> package **`host/verbs`** (`Split`/`Quote`, the verb table `Lookup` with steps 1–5 and 4b, `ParseInterleaved` moved
> from `cmd/ailang-worldd/session.go`, `EpisodeReportFlags` and the bar's admission; VB1–VB5); the **CLI repoint** (the
> ten `parseInterleaved` call sites, `episode report` on the shared flag set, its printer moved to `host/daemon` as
> `WriteEpisodeReportText`); **`GET /surface/verb.json?q=<line>`** behind `approvalRoute`, read-only, with
> `d.episodeReport` extracted from the `/v1` handler (VR1–VR4, VS1); and **the page**: `go(raw)` asks the route first,
> **only a schema-valid `outcome:"address"` navigates**, every failure (non-2xx, network error, unreadable body, the
> **15 s budget**, a **superseded** ask) refuses with the view and the input unchanged (VP1–VP6), proven in **live
> Chrome** (XV1) and documented as **QUICKSTART §19**.
> - **Measured: 1,909 changed lines** — `git diff --numstat db2b7423 HEAD -- cmd host docs`, insertions + deletions,
>   after a commit of each of the eight scratch stages (**826 production** — 736 Go, a move counting both sides, + 90
>   page; **1,069 test** incl. the 20-line golden; **14 docs**; `-- cmd host` alone: **1,895**).
>   Margin to the 2,000 cap: **91**. **M1–M3 measure 1,491 ≤ 1,500: the design's fallback cut does not trigger;
>   nothing moves to V2** (V1-D3).
> - **After this sprint** an operator types `episode report <episode>` (or `--json`) into the World app's bar and reads
>   the CLI's own bytes; `approve <mr-id|sha256 ref>` opens row 201's card; every other operator verb shows the exact
>   CLI line and says the app cannot run it yet. Nothing on a page changes World.
> - **Eleven premises of the design or of HEAD are false or incomplete and are routed minimally**, and one of the
>   planner's own (P12) was caught by the node harness (§1.4).
> - **Every line of code, test and doc is a unified diff in Appendix P**, generated from compiled, vetted, gofmt-clean,
>   green scratch stages; extracted **from this file** with §2's `awk`, applied to a fresh `git archive db2b7423`
>   export and compared byte-for-byte with the stage tree and the mutant tree (§1.2 C17). The executor applies each
>   patch and runs its acceptance; it writes no code by hand.

- **Design**: [w-world-address-bar-verbs.md](w-world-address-bar-verbs.md), **r2** (quorum r1+r2, controller
  narrow-refinement carve-out r2, independent audit CARVE-OUT VALID 84; iteration 394; **not reopened**). Where this
  plan and the design differ, this plan was measured against `db2b7423` and wins; each difference is a decision
  (§1.3) or a premise found false (§1.4). The two calls the r2 audit left to the planner are **V1-D1** (a superseded
  ask) and **V1-D2** (VP2's clock; the budget is not a tunable).
- **Authority**: D-WORLD-107 = A, D-WORLD-109 = A (test stores, temp dirs and test keys only), D-WORLD-89 (the agent
  never approves; `approve` only *opens* the card). No ledger ask (design Q5).
- **Planner base**: `db2b7423327d5223cd1a37cf40e0794c2019026a` (= origin/dev, PR #580's merge). Worktree
  `/Users/voightkampff/dev/sunholo-data/.wt-world-iter395-row243-v1`, branch `sprint/world-iter395-row243-v1`.
  The design was measured at `48c92a52`; the drift since is PRs #577–#580 (§1.2 C2) and every design line this plan
  relies on was re-measured at `db2b7423`.
- **Planner scratch** (outside the repo, never under `/tmp`): `$HOME/.ailang/state/world-iter395/plan-scratch/`
  (`base/` the pristine export (a throwaway git repo, tag `BASE`), `work/` the final tree the mutant runs' sources
  came from, `work2/` the stage repo (tags `M1a` … `M4b`, `WORK2_EQUALS_FINAL`), `final-tree/` (work2 at `M4b`, the
  final gate run), `patches/`, `tools/` (`STAGES`, `stages.sh`, `mid_cmd_ailang-worldd_episode_report.go.py`,
  `mkpatches.sh`, `accept.sh`, `accept_all.sh`, `gates.sh`, `mutants.py`, `mutate.py`, `roundtrip.sh`,
  `assemble.py`, `mkjson.py`), `logs/` (`base_gates/`, `final_gates/`, `accept_<M>.log`, `mutant_ledger.tsv`,
  `mut/`, `roundtrip.log`)).
- **Executor**: as routed by the controller. **Three runs**: Run 1 (M1a, M1b, M1c — `host/verbs` and the CLI
  repoint), Run 2 (M2, M3a, M3b — the printer, the route, its proofs), Run 3 (M4a, M4b — the page, live Chrome,
  QUICKSTART). Each sub-milestone is one `patch -p1` (§2) plus its acceptance (§3). **Apply, never retype.**
- **Touched** (§5): new `host/verbs/{split,lookup,flags}.go` + tests, `host/daemon/{surface_verb.go,
  episode_report_text.go}` + tests, `host/surfacepage/testdata/verb-route.golden.json`, the new cmd tests; edited
  `cmd/ailang-worldd/{approve,approver,episode,episode_finish,episode_launch,preset,session,session_new}.go` (the
  call sites), `cmd/ailang-worldd/episode_report.go` (flag set; printer moved out), `cmd/ailang-worldd/episode_markdown.go`
  (§1.4 P1), `host/daemon/{episode_report.go (the extraction), daemon.go (+1 route)}`,
  `host/surfacepage/app/{app.js, index.html}`, `host/surfacepage/app_test.go` (+2 −1, §1.4 P2), `docs/QUICKSTART.md`
  (+14, §19). **Outside the design's §8 file list**: `host/daemon/episode_report_text.go` (the printer's home, design
  §2.2 names `daemon.WriteEpisodeReportText`), `cmd/ailang-worldd/episode_markdown.go` (P1), the eight call-site files
  other than `session.go` (design §2.1(c) names them, §8 lists only `session.go`), `host/surfacepage/app_test.go` (P2),
  the golden (V1-D12), `docs/QUICKSTART.md` (design §7) and every new `_test.go`.
- **Not touched**: `host/surfacepage/surface.js` (the static page's bar, `parseAddress`), `host/surfacepage/{approve,
  passkey,decide}/**`, every row-201 handler/law (`surface_approve*.go`, `surface_decide*.go`), the `approve` CLI's
  behaviour (only its `parseInterleaved` call is repointed), `surfaceCSP`, `world/*.ail`, `packages/**`, `scripts/**`,
  `tools/launchd/*`, `cmd/world-publish/*`, any skill, `~/.ssh`, `~/.ailang/world`, `~/.ailang/state/mission-v1*`.
  **No `/v1` route** (`TestFrozenV1RouteTableMatchesMux` stays 13), **no CORS header**, **no POST, store write,
  credentials-store open, process spawn or `.ail`**.

---

## 0. Binding constraints (each is restated where it governs)

1. **Test stores, temp dirs and test keys only (D-WORLD-109 = A).** Every store is under `t.TempDir()`; every daemon
   is `daemon.New` behind `httptest` (loopback) or an `httptest.NewRecorder`; Chrome is the A2 driver
   (`testdata/passkey_cdp.mjs`) with its own temp `--user-data-dir`. No real store, `~/.ssh`, `/dev/tty`, browser
   profile, mint, publish or listener beyond loopback.
2. **Pure core, effects at the boundary (coding standards S1–S3).** `host/verbs` imports only `errors flag fmt io
   regexp strings time unicode/utf8` — no I/O, no clock call, no randomness (VS1 pins the import set). The one effect
   V1 adds is the daemon's existing report read under `d.readCtx`.
3. **Read-only route.** `GET /surface/verb.json` reaches exactly `d.readCtx`, `d.episodeReport`,
   `d.writeInternalErrorLog`, `d.readDeadline` and its own helper (VS1 pins the selector set and the import set: no
   `host/store`, no `os/exec`); no POST arm; behind `approvalRoute` (R2's 421).
4. **No silent fallback.** Only a 200, schema-valid `outcome:"address"` reaches the page's address path; every
   failure class says why and leaves the view and the input (§2.3 of the design; VP2; mutants `J-*`).
5. **The CLI is unchanged.** Every existing test of `cmd/ailang-worldd` passes unedited (the whole package runs in
   M1c's and M2's acceptance); the only edited existing test file in the sprint is `host/surfacepage/app_test.go`
   (§1.4 P2).
6. **One mutant per new production branch, anchored to the diff** (§4): the named test and the behavioural suite with
   every inventory/static test `-skip`ped; a suite survivor in `host/daemon`/`cmd` is re-run on the whole package;
   every run appends a ledger line; the table's ids are the ledger's ids.
7. **No `t.Parallel()`** in any new test (`TestNoParallelWriteForkPackagesOutsideVerifygate`).
8. **Bisectable and safe at every commit.** Each sub-milestone ends gofmt-clean, `go build ./...` and `go vet` green,
   the `_test.go` compile fence green, its own tests and every earlier one green (§1.2 C15).
9. **No git writes by the executor** (§7). The controller commits from snapshots.

---

## 1. Pre-flight

### 1.1 Baseline at `db2b7423` (measured by the planner, 2026-11-16)

| Leg (`tools/gates.sh <tree> <logs> base`) | rc | Time |
|---|---|---|
| `status` | **0** | 0 s |
| `verify_ail` | **0** | 34 s |
| `gofmt` | **0** | 0 s |
| `build` | **0** | 3 s |
| `vet` | **0** | 3 s |
| `fence` | **0** | 27 s |
| `pkgs` | **1** | 540 s |
| `race` | **0** | 107 s |
| `inv` | **0** | 42 s |
| `chrome` | **0** | 53 s |
| named PASS: 51 | | |
| `verify_go` | **1** | 642 s |
| GATES_DONE 1452s FAILS: pkgs verify_go | | |

Every command ran with `export AILANG_BIN=$HOME/.pinned-ailang/ailang` (AILANG v0.41.0), `go1.26.6 darwin/arm64`,
node v26.13.0, in a `git archive db2b7423` export at `plan-scratch/base` (a throwaway git repo; `diff -r cmd`, `host`,
`docs` against the worktree: identical; the worktree itself was only read), with §6's script in `base` mode (the
sprint's new tests and `host/verbs` left out). Load average 12–20 throughout (this planner's stage, acceptance and
mutant runs alongside). Logs: `plan-scratch/logs/base_gates/` (`legs.txt`, one log per leg).

- `verify_ail`: `✓ verify gate PASSED: 16 required identities verified, 40 named tests pass, 503 se-tools named tests pass`.
- `pkgs` rc 1: `ok host/surfacepage`, `ok host/daemon 539.6s`; `cmd/ailang-worldd` red **only** on
  `TestHostTranscriptRefusesSymlink/claude` ("`--version failed or timed out (F10)`") and `TestWorldContextAndStatusline`
  ("`world-context: signal: killed`") — 5 s subprocess deadlines under load in untouched host-launcher tests (the
  directive's `TestHostsLaunchersExecHost` family).
- `inv` rc 0: **51** named `--- PASS`, 0 SKIP, 0 FAIL, none missing. `chrome` rc 0: the twelve page tests and the two
  `host/daemon` app-shell drills `--- PASS`.
- `verify_go` rc 1: **56** packages `ok`; red only in untouched packages, every one a deadline under load:
  `cmd/ailang-worldd` (`TestCLIRealSubprocessEpisode`, `TestDoctorQuickIsFast`, `TestWorldContextAndStatusline`,
  `TestHostsLaunchersExecHost`), `cmd/world-publish` (`TestEpochNominateRefusals`), `host/broker`
  (`TestExecProbeNet6ForeignDialIsRetriedNotRefused`, `TestExecProtectedBaselineFailure`), `host/daemon`
  (`TestServeStartsWhenTheProbeSpendsTheStartupClock`, `TestRegistryPhaseHasItsOwnClock`), `host/verifygate`
  (`TestMissionConfigDefaultFlowIntegration`), `scripts/marketingcheck` (`TestMarketingSubprocessDeadline`,
  `TestMarketingSupervisorFakeChromeTimeoutKillsGroup`) — the directive's known-flaky families.
- Whole script: `GATES_DONE 1452s FAILS: pkgs verify_go`.

**Rule (carried):** a red at the base measures the repo, not the sprint. The controller re-runs a red untouched
package alone once and does not chase it. **Never start a gate under `nohup`.**

### 1.2 Codebase claims checked (at `db2b7423` unless "work" — the planner's final scratch tree)

| # | Claim | Command → observed output |
|---|---|---|
| C1 | The base | `git rev-parse HEAD` (worktree) → `db2b7423327d5223cd1a37cf40e0794c2019026a`; `git status --porcelain` → only the two untracked plan files once written; `$HOME/.pinned-ailang/ailang --version` → `AILANG v0.41.0`; `go version` → `go1.26.6 darwin/arm64`; `node --version` → `v26.13.0` |
| C2 | **Drift since the design's base** (#577–#580) | `git diff --shortstat 48c92a52 db2b7423 -- cmd host docs` → `12 files changed, 272 insertions(+), 3 deletions(-)` (the controller's figure); none of the files this sprint edits but `daemon.go` (+1 route, `GET /surface/home-summary.json`, line 1126), `app_test.go` (+16, not the export list), `app.css` (untouched here) and `QUICKSTART.md` (+3 −1: §18 moved to `1862`, file `1,897` lines) |
| C3 | **`parseInterleaved`: ten call sites + the definition** (design §1.2 "11 lines in 9 files", unchanged) | `grep -c 'parseInterleaved(' cmd/ailang-worldd/*.go \| grep -v ':0'` → `approve.go:1 approver.go:1 episode_launch.go:1 episode.go:1 episode_report.go:1 episode_finish.go:1 preset.go:1 session_new.go:2 session.go:2`; the definition `session.go:181`; no `_test.go` names it (`grep -rn parseInterleaved --include='*_test.go' cmd host` → empty, with the known positive `grep -c` above in the same session) |
| C4 | **The printer's other caller** (§1.4 P1) | `grep -rn 'printReport\|execCacheLine\|readsLine\|printRefusals\|sortedCounts' cmd host --include='*.go' \| grep -v '^cmd/ailang-worldd/episode_report.go'` → `episode_report_test.go:102,131,256,273,279` (`printReport`, kept as a one-line wrapper) and `episode_markdown.go:86` `for _, kv := range sortedCounts(r.CallsByTool)` with `kv.name, kv.n`, `:94-95` `printRefusals(&b, …)` |
| C5 | **The switch has 28 arms** (design V33, re-derived) | `sed -n 375,520p cmd/ailang-worldd/main.go \| grep -cE '^\s*case "'` → `28`; arms `approve approver call commit demo doctor episode head health help hosts log object outcome preset profile provenance registry serve service session setup surface-index tools version whoami why world`; work: VB4 `TestVerbTableMatchesTheCLISwitch` derives the same 28 with `go/ast` |
| C6 | **`handleEpisodeReport` maps only the builder's error to a read timeout** (§1.4 P6) | `episode_report.go:707` `func (d *Daemon) handleEpisodeReport`; `:715-720` `rep, ok, err := buildEpisodeReport(…); if err != nil { if timedOut(ctx, err) { writeReadTimeout(…)`; `:724-733` the start and finish read errors → `d.writeInternalError` (no timeout check) |
| C7 | **`approvalRoute` and the route pins** | `surface_approve.go:76` `func (d *Daemon) approvalRoute(…)`; `daemon.go:1139-1147` the approval and decision routes, the last `GET /surface/decide.js` at `:1147`; `TestApprovalRoutesRegistered`'s regex `/surface/(?:requests\|approve)`, `TestDecisionRoutesRegistered`'s `decide\|decisions` — neither matches `/surface/verb.json`; `TestSurfaceAppRoutes`'s 404/421 lists name no `verb` path (work: all three PASS unedited) |
| C8 | **No CORS header; no process in the daemon** (design V31, V24) | `grep -rniE 'Access-Control\|cors' host/daemon --include='*.go' \| grep -v _test` → the three comment lines; `grep -rln 'exec.Command' host/daemon --include='*.go' \| grep -v _test \| wc -l` → `0`, control `grep -rln 'exec.Command' host --include='*.go' \| grep -v _test \| head -2` → `host/agentrun/agentrun.go`, `host/sshsign/sshsign.go`; work: unchanged (`surface_verb.go` sets no header; VR1 asserts no `Access-Control-*` on 405) |
| C9 | **The page's bar at HEAD** (§1.4 P3, P4) | `app.js:1201` `async function go(raw)`; `:1430` `if (e.key === 'Enter') { go(e.target.value); e.target.value = ''; …` (the input is cleared before `go` runs); `:928` `$('addr').placeholder = st.d !== 0 ? '' : … : 'go to a place — row-189@6 · shared@0ffbb9c'` (renderCrumbs overwrites `index.html:67`'s placeholder); `:931` `function say(t)` fades after `2600` ms |
| C10 | **`TestSurfaceAppPure` pins app.js's exact export list** (§1.4 P2) | `app_test.go:174` `assert.deepEqual(Object.keys(app).sort(), ['LEVELS', 'adaptModel', 'addressTarget', …, 'verdict']);` |
| C11 | **No asset-hash pin over `app.js`/`index.html`** | `grep -rn 'AppScript\|app\.js' --include='*_test.go' host cmd \| grep -i 'sha256\|hash'` → empty; known positive in the same tree: `surface_approve_page_test.go:54` pins `passkey.js`/`approve.js` (untouched here) |
| C12 | **A real daemon in `cmd` tests reports nothing for the walk rig** (§1.4 P5) | `why_test.go:67` `r.post(genesis)` / `:78` `POST /v1/commit` — no journal intent; `grep -n AppendIntent cmd/ailang-worldd/*_test.go` → empty (known positive: `host/daemon/episode_report_test.go:129` `s.d.store.AppendIntent`); work: a first VR3 on `newWalkRig` failed `exit 1 … episode w1 has no committed calls`; `seedReportCalls` (VR3) then PASS |
| C13 | **`store.MintRequestIDPattern` is `^mr-[0-9a-f]{16}$`** (`host/store/mint_requests.go:46`); `host/verbs` may not import `host/store` (S1) | work: `lookup.go` carries the same pattern; `TestSurfaceVerbIsReadOnly` checks `verbs.ApproveTarget` against `store.MintRequestIDPattern` on six ids (mutant `A-MR-PAT` killed) |
| C14 | **The page and the route ship in one binary; `surface.js` untouched** | `app.go:15` `//go:embed app/app.js`; work: `git diff --name-only BASE M4b -- host/surfacepage/surface.js host/surfacepage/approve host/surfacepage/passkey host/surfacepage/decide` → empty |
| C15 | **Every sub-milestone is green when committed** | `tools/accept_all.sh M1a … M4b` (an APFS clone of the stage repo checked out at each tag; `logs/accept_<M>.log`): `M1a` ACCEPT_OK M1a; `M1b` ACCEPT_OK M1b; `M1c` ACCEPT_OK M1c; `M3b` ACCEPT_OK M3b; `M2` ACCEPT_OK M2; `M4a` ACCEPT_OK M4a; `M4b` ACCEPT_OK M4b; `M3a` ACCEPT_OK M3a — **`ACCEPT_OK` for all eight** on the final stage repo (M1a–M3b are byte-identical across the stage-repo rebuilds; only M4a changed). One earlier run was red and is recorded: M2's first run
(`ACCEPT_RED M2: p … ./cmd/ailang-worldd/`) failed only `TestWorldContextAndStatusline`, `TestHostsLaunchersExecHost`,
`TestPiLauncherRefusesMissingAdapter/missing_at_init` — each "`signal: killed`" at a 5 s subprocess deadline at load
≈ 19, untouched host-launcher tests (red at the base too, §1.1); re-run: `ACCEPT_OK M2`. M4a/M4b were re-run after
§1.4 P12's fix and the `constructor`-outcome case (both `ACCEPT_OK`). The live-Chrome lines (M4a: the two *Needs you*
page tests; M4b: XV1 + one *Needs you*) ran with `WORLD_REQUIRE_PASSKEY_BROWSER=1`: `--- PASS`, none SKIP |
| C16 | **Every sub-milestone is committed and measured** | `tools/stages.sh` commits each stage in `work2` and measures `git diff --numstat BASE <tag> -- cmd host docs`: `NUMSTAT M1a 203`, `M1b 510`, `M1c 697`, `M2 905`, `M3a 1248`, `M3b 1491`, `M4a 1762`, `M4b 1909 (cmd host: 1895)`, ending `WORK2_EQUALS_FINAL`; the round-trip repo agrees: M1a 203 M1b 510 M1c 697 M2 905 M3a 1248 M3b 1491 M4a 1762 M4b 1909 |
| C17 | **The patches reproduce the stage tree and the mutant tree** | `tools/roundtrip.sh <this plan>`: a fresh `git archive db2b7423` export + `patch -p1` (after `--dry-run`) of Appendix P's eight patches **extracted from this file with §2's `awk`** (each `cmp`-equal to the generated patch; no `.orig`/`.rej`), each committed and measured (C16) → `diff -r` of `cmd`, `host`, `docs` against `plan-scratch/work2` (`M4b`, the stage tree the gates and the final mutant runs used) and against `plan-scratch/work` (the tree the stages were cut from): **identical** (`ROUNDTRIP_IDENTICAL_TO_STAGE_TREE_AND_MUTANT_TREE`, `logs/roundtrip.log`) |
| C18 | **Every mutant; the table's ids are the ledger's** | 120 rows (106 K, 7 N, 2 L, 0 E, 5 controls); 224 ledger lines; the last line of each row behaves as expected for 120 of 120 run rows — **106 `K` rows killed by the named test and by the behavioural suite** (inventory/static tests `-skip`ped); **7 `N`
rows** (`S-UTF8`, `F-HELP`, `F-JSON`, `E-AS`, `E-NOTFOUND`, `R-NOTFOUND`, `R-JSON`: an unused variable or import —
not mutants) each **rewritten to compile** (`…b`) and then killed, or, for `E-ASb`, the named survivor; **2 `L` named
survivors** (`E-ASb`, `JW-GO-EMPTY`, §3.1, Appendix C); **0 `E`** (the three branches that would have been
equivalent — a `back` arm, the route's own length check, a JSON trim — were removed from production instead, V1-D8,
V1-D10, V1-D4); **5 unmutated controls** green (`CTL-D`, `CTL-C` also green on the whole package); **0 not as
expected**. The ledger's id set equals the table's (every ledger id is a row; every row has a line; checked by
`tools/assemble.py`, which fills this table from the ledger's last line per id). The four live-Chrome `JW-*` rows: three killed by XV1 (the
route skipped, the bar cleared on Enter, the panel's bytes changed), one survivor (`JW-GO-EMPTY`) |
| C19 | **The gates, on the final tree** | `tools/gates.sh <final-tree> logs/final_gates final` (`final-tree` = `work2` at `M4b`, a git repo with every stage committed, load average 15–25): `status` 0; `verify_ail` rc 0 → `✓ verify gate PASSED: 16 required identities verified, 40 named tests pass, 503 se-tools named tests pass` — the base's log line for line but one count (`cache control fired on 46 non-.ail file(s)`, base 41: the control counts non-`.ail` files, and the sprint adds Go files; no `.ail` changed); `gofmt`, `build`, `vet`, `fence` rc 0; **`pkgs` rc 0** (`ok` host/verbs, host/surfacepage, host/daemon 548 s, cmd/ailang-worldd — the base's two cmd load reds did not recur); **`race` rc 0** (`ok` all four; no `DATA RACE`); **`inv` rc 0: 65 named `--- PASS`** (the base's 51 + 14 new), 0 SKIP, 0 FAIL, none missing — incl. `TestFrozenV1RouteTableMatchesMux` (13), `TestSurfaceAppRoutes`, `TestApprovalRoutesRegistered`, `TestDecisionRoutesRegistered`, `TestDaemonDependencyAllowlist`, `TestWorldBoundaryDependencyAllowlist`, `TestProductionContextRoots`, `TestProductionGoSurface`, `TestDaemonCredentialWriterSites`, `TestEveryVerbHasHelp`, `TestSurfaceAppPure`, `TestEpisodeReportRESTRoute` and the quickstart tests; **`chrome` rc 0: 15 `--- PASS`** (the twelve page tests + XV1 in `cmd`, the two app-shell drills in `host/daemon`); `verify_go` rc 1: **61** packages `ok`, red only `host/capsule` (`TestCapsuleTemplateIsCopiedNotShared`, an archived-interpreter template check under load) and `host/verifygate` (`panic: test timed out after 10m0s` in `TestSeToolsGateArmsAreLoadBearing/plan-bytes`) — both untouched, both **green re-run alone** (`ok host/capsule 32.7s`, `ok host/verifygate 371.2s`; `logs/final_capsule_alone.log`, `logs/final_verifygate_alone.log`); the base's twelve `verify_go` reds did not recur. **`GATES_DONE 1969s FAILS: verify_go`**, every red an untouched load test green alone; `ALL_DONE.<sha>` written (`logs/final_gates/`) |

### 1.3 Planner decisions where the design leaves a choice

- **V1-D1 — A superseded ask (the r2 audit's first call).** A second submit while one is pending **aborts the first**
  (`askVerb` marks the pending job `superseded`, then `ctl.abort()`), and the first ask **ends with `{superseded:
  true}`, which `verbAct` acts on not at all**: no message, no panel, no navigation, the input untouched by it. The
  second ask is now the bar's question and answers for itself — explicitly — within its own 15 s budget. This differs
  from the budget timeout by design: **timeout** = the operator's *only* question went unanswered → `say('the verb
  route did not answer in time')`, view and input kept; **superseded** = the operator asked something newer → the old
  answer is moot and saying anything about it would race the newer answer's message. It is the app's existing rule for
  states (`openState`, `app.js:1165`: "a newer open supersedes an older one"). It cannot weaken "no silent fallback":
  a superseded ask can never navigate (the abort rejects both of its awaits through one `stop` promise raced against
  the fetch and the body, and `verbAct` returns on `superseded` before every arm), and no failure is ever converted
  into an address. An answer that had already resolved when the newer ask started (possible when both settle in one
  turn) is dropped too, by `askVerb`'s `job.superseded ? {superseded:true} : …` after its awaits. Proven in VP2's two
  superseded cases — (A) the first ask pending and released after the second starts: its `AbortSignal` is aborted, it
  ends `{superseded:true}`, `verbAct` calls nothing, the second's refusal stands; (B) the first ask's stub resolved
  at once and the second started in the same turn: still `{superseded:true}` — mutants `J-SUPERSEDE`,
  `J-CATCH-SUPERSEDED`, `J-SUPERSEDE-CHECK`, `J-ACT-SUPERSEDED`. *Rejected*: (a)
  ignoring the second submit while one is pending (a wedged route would lock the bar for 15 s); (b) letting both run
  and acting on the last to *arrive* (a stale answer could navigate after a newer refusal).
- **V1-D2 — VP2's clock; the budget is not a tunable (the audit's second call).** `askVerb(raw, io)` takes its
  effects as `io = {fetch, AbortController, setTimeout, clearTimeout}`; production `go(raw)` passes the window's
  own, the node harness passes **a controlled clock** whose `setTimeout` records `(fn, ms)` and schedules nothing.
  VP2 asserts the recorded `ms === app.VERB_ROUTE_BUDGET_MS === 15000`, "advances the clock" by calling the recorded
  `fn`, and sees the timeout refusal for a stalled fetch and for a stalled body, with `clearTimeout` called with the
  timer id on every path (success, network error, timeout). No 15 s wall is spent. **`VERB_ROUTE_BUDGET_MS` stays the
  design's named constant, exported read-only for the test — not a tunable**: no query parameter, storage key or
  config feeds it, because a knob is a way to make the wait unbounded (the axiom the budget answers) and nothing needs
  one (the server's own bound is `readDeadline = 10 s`, `daemon.go:144`, so a slow but healthy report answers first).
  Mutants `J-BUDGET` (the constant) and `J-BUDGET-USE` (the timer not using it) are killed.
- **V1-D3 — Scope; the fallback cut is measured, not taken.** M1–M3 measure **1,491** (`-- cmd host docs`, after
  M3b's commit) ≤ the design's 1,500 threshold, so XV1 and the `approve` navigate arm **stay**; the sprint measures
  **1,909** ≤ 2,000. Had a later re-measure crossed 1,500 at M3b, the cut is M4b's XV1 file (78 lines) and the
  `navigate` arm (`ApproveTarget` → `terminal`); it did not, so nothing moves to V2 (Appendix C).
- **V1-D4 — The answer's wire** (design §2.2 leaves member shapes open). `{v:1, schema:"world/surface-verb/v1", q,
  argv, verb, class, outcome}` on every answer (`argv` `[]`, never null); `reason` leads with the refusal's name for a
  `Split` refusal (`TooLong: a line is at most 1024 bytes`); `result` carries exactly one of `text` (the CLI's table),
  `json` (the CLI's exact `--json` stdout: the `/v1` encoder's compact bytes and one newline) or `stderr` + `exit:1`
  (the CLI's not-found line, newline included); `cli` = `Quote(["ailang-worldd", argv…])`. A flagged or malformed
  `approve` answers `outcome:"refused"` with `class:"navigate"`. Every answer is written by A3a's
  `writeApprovalJSON` (no-store, nosniff); a missing `q` is A3a's `400 BadRequest` refusal; an empty `q` answers
  `address` (the page never sends it: empty is home without a fetch).
- **V1-D5 — What the page shows.** `refused` → `say(reason)` in the header (design §2.3; VP6's "refusal panel" is read
  as that line — the reason of a never/other verb already ends `at a terminal: <cli>`); `result` and `terminal` → one
  panel (`#verb`, a `.pop` built by `verbPanel` with `textContent` only: header `$ <cli>`, a `<pre>` body; `terminal`
  adds **Copy the line**); Escape closes it. The panel is created by script, so `index.html` gains only the
  placeholder words. `navigate` → `location.assign(href)` only for `/surface/approve/mr-…` or
  `/surface/decide/sha256:…` (VP5).
- **V1-D6 — The input is kept on a refusal** (design "the input kept"; §1.4 P3). Enter no longer clears the bar;
  `verbAct` clears it only when it acts (address, navigation, panel) and only if it still holds the submitted line
  (`clear: () => { if ($('addr').value === raw) $('addr').value = ''; }`).
- **V1-D7 — `Split`'s grammar beyond the design's sentence.** Adjacent segments join into one word (`a'b'"c"` →
  `abc`), so `Quote` can print a `'` as `'"'"'` and `Split(Quote(a)) == a` holds for every admissible `a` (VB2:
  2,000 seeded cases over an alphabet with every shell character, plus every VB3 argv). An unquoted backslash is
  literal; `Quote` single-quotes every word outside `[A-Za-z0-9@%+=:,./_-]`, so a printed `cli` line is safe in a
  shell whatever the word holds.
- **V1-D8 — Step 3 is `world` with one or two words only; `back` reaches `address` through step 5** (it is not a CLI
  word), which is the design's outcome for every `back` input; keeping a separate `back` arm would be an equivalent
  mutant. Step 4b seizes `episode` (no subcommand row matched) and `world` (three or more words, not `world get`); every
  other seized word has a one-word row that already matches (design 4b's "the row's own refusal text applies").
- **V1-D9 — The extraction keeps `/v1` byte-exact.** `d.episodeReport(ctx, ep) (EpisodeReport, bool, error)` wraps
  only the builder's failure in `reportBuildError`; `reportTimedOut(ctx, err)` is true only for it, so the `/v1`
  route still answers 503 `Timeout` for a builder timeout and 500 `Internal` for a start/finish read failure, as at
  HEAD (§1.4 P6). The verb route uses the same classification for its refusal text (`read deadline (10s) exceeded` or
  `internal store failure`) and logs the error for the operator; it never shows the error.
- **V1-D10 — The 1,024-byte bound is `Split`'s first statement**, before any tokenizing, so the route needs no second
  check of its own (design §2.2 step 1's property holds; a duplicate check was an equivalent mutant and was removed).
- **V1-D11 — `approve`'s id pattern** is `host/verbs`' own copy of `store.MintRequestIDPattern` (`host/verbs` stays
  standard-library-only); VS1 pins them equal on six ids.
- **V1-D12 — Fixtures are generated, never hand-typed.** VR3 runs the real CLI (`run([... "episode","report",…])`)
  and the real route against one `daemon.New` on a store seeded with real journal intents and commits
  (`seedReportCalls`, the layout of `host/daemon`'s `reportSeed`), for **two episodes** — `w1` with a plan refusal and
  no handler refusal, `w2` with a handler refusal and no plan refusal (each table's optional `#N tool: reason` lines
  present in one and absent in the other) — in table and `--json` form, plus the not-found case. The page's node test
  replays **`testdata/verb-route.golden.json`**, which `TestSurfaceVerbOutcomes` writes from the real route under
  `WORLD_UPDATE_VERB_GOLDEN=1` and otherwise compares byte for byte (18 answers, incl. a report whose tool name and
  refusal are hostile markup, VP3).
- **V1-D13 — Sub-milestones** (§1.5): eight, each ≤ 350 lines, each green at commit. `host/verbs` is split by part
  (`split.go`, `lookup.go`, `flags.go`) so M1a–M1c each carry their own code and tests.

### 1.4 Premises found false or incomplete at HEAD (each routed around minimally)

- **P1 — Design §7 "any other caller (e.g. `episode_markdown.go`) keeps compiling through the moved, exported
  names".** It does not without an edit: `episode_markdown.go:86-95` reads the unexported `namedCount` fields
  (`kv.name`, `kv.n`) and calls `printRefusals`. Routed: the moved helpers are exported as `daemon.SortedCounts`
  (`NamedCount{Name, N}`) and `daemon.WriteReportRefusals`; `episode_markdown.go` +4 −4 (M2). Its tests
  (`TestFinishMarkdownOrderCapAndAC5`, `TestFinishMarkdownCLIRefusalsLoopback`, `TestFinishSharedScrubVectors`) pass
  unedited.
- **P2 — Design §2.3 "`verbOutcome` is exported to the node harness like `addressTarget`".** `TestSurfaceAppPure`
  pins the exact export list (`app_test.go:174`), so any export edits it. Routed: +2 −1 (the five new names appended,
  the list sorted); the only edited existing test in the sprint (M4a).
- **P3 — Design §2.3 "the input kept".** At HEAD the Enter handler clears the bar before `go` runs (`app.js:1430`).
  Routed: V1-D6.
- **P4 — Design §2.3 "The placeholder gains `· episode report <episode>`".** `renderCrumbs` (`app.js:928`) replaces
  `index.html`'s placeholder on every timeline render. Routed: both gain the words (M4a).
- **P5 — Design §6 VR3 "run the real CLI … against the loopback daemon" with "the existing report fixtures'
  episode".** The CLI's report tests use a fake daemon (`fakeReportDaemon`); the cmd walk rig commits over `POST
  /v1/commit` with no journal intent, so a real daemon answers "no committed calls" (§1.2 C12); `host/daemon`'s
  `reportSeed` is package-private. Routed: `seedReportCalls` (55 lines, VR3's file) seeds the store with the exported
  store API before `daemon.New`; XV1 reuses it.
- **P6 — Design §2.2 "`handleEpisodeReport`'s body is extracted unchanged into `d.episodeReport`".** The body's
  timeout mapping depends on *which* read failed (§1.2 C6). Routed: V1-D9 (`reportBuildError`); named survivor
  `E-AS` (§4).
- **P7 — Design §7 "QUICKSTART §19 after §18 (`docs/QUICKSTART.md:1861`)".** §18 is at `1862` (C2), and both §18
  tests read from §18 to the end of the file (`TestQuickstartSection18CardMatchesThePage` checks every bold label
  against the card's page; `TestQuickstartSection18FlagsMatchTheCLI` counts `ailang-worldd approve sha256:`/`--list`
  lines). Routed: §19 has no bold label and no such line; both §18 tests PASS unedited.
- **P8 — Design §2.1(a) "`Quote(argv)` … `Split(Quote(a)) == a`" with `'…'` literal.** A word holding `'` cannot be
  single-quoted alone; routed by V1-D7 (segment joining).
- **P9 — Design §2.4's `back` row says step 3.** Outcome identical through step 5 (V1-D8).
- **P10 — Design §2.2 step 1 "checked on the raw query value before `Split`".** `Split`'s first statement is that
  check (V1-D10).
- **P12 — (the planner's own premise) "an answer cannot beat the abort".** The planner's first `app.js` dropped the
  post-await `superseded` check as unreachable in a browser; the node harness's case (B) — a stub that resolves in the
  same turn — showed the answer is acted on without it (`TestSurfaceAppVerbs` red on that stage tree; the first
  mutant run's `app.js` lines are void, §4). Restored (V1-D1) and pinned (`J-SUPERSEDE-CHECK`).
- **P11 — Design §2.3 "`say(reason)`" as the refusal display.** `say` fades the header line after 2.6 s (§1.2 C9);
  the input is kept (V1-D6), so pressing Enter again re-shows it. Residual RS-SAY (Appendix C).

### 1.5 The size

Lines are insertions + deletions: **`git diff --numstat <base> HEAD -- cmd host docs`** (the directive's filter;
`-- cmd host` alone in its own column). **Every row is measured** from the scratch stages, each committed in the
throwaway repo `work2` (tags `M1a` … `M4b`), and re-measured on the round-trip tree (§1.2 C16).

| Sub-milestone | Production | Test | Docs | Own | Cumulative vs `db2b7423` (after its commit) | `-- cmd host` | Run |
|---|---|---|---|---|---|---|---|
| M1a — `Split`, `Quote`; VB1, VB2 | 110 | 93 | 0 | 203 | **203** | 203 | 1 |
| M1b — `Lookup`, the table, `ApproveTarget`; VB3, VB4 | 133 | 174 | 0 | 307 | **510** | 510 | 1 |
| M1c — `ParseInterleaved` moved, `EpisodeReportFlags` + admission; the ten call sites; VB5 | 129 | 58 | 0 | 187 | **697** | 697 | 1 |
| M2 — `episode report`'s printer moved to `host/daemon` | 208 | 0 | 0 | 208 | **905** | 905 | 2 |
| M3a — `d.episodeReport`; `GET /surface/verb.json`; VR1, VR4, VS1 | 156 | 187 | 0 | 343 | **1,248** | 1,248 | 2 |
| M3b — VR2 + the page's golden; VR3 (the real CLI) | 0 | 243 | 0 | 243 | **1,491** | 1,491 | 2 |
| M4a — the page: `go` asks the route; VP1–VP6 | 90 | 181 | 0 | 271 | **1,762** | 1,762 | 3 |
| M4b — XV1 in live Chrome; QUICKSTART §19 and its test | 0 | 133 | 14 | 147 | **1,909** | 1,895 | 3 |
| **The sprint** | **826** | **1,069** | **14** | | **1,909 ≤ 2,000** (margin 91) | 1,895 | |

(Production counts a move on both sides: the printer is −99 in `cmd` and +105 in `host/daemon`; `ParseInterleaved`
−15 in `session.go`, +15 in `flags.go`. M3b's "test" includes the 20-line golden.) Per file (+/−):
`split.go` +110, `split_test.go` +93; `lookup.go` +133, `lookup_test.go` +94, `verbs_parity_test.go` +80; `flags.go`
+77, `flags_test.go` +58, `approve.go`/`approver.go`/`episode.go`/`episode_finish.go`/`episode_launch.go`/`preset.go`
+2 −1 each, `session.go` +2 −17, `session_new.go` +3 −2, `episode_report.go` (cmd) +4 −1 at M1c and +2 −98 at M2;
`episode_markdown.go` +4 −4, `episode_report_text.go` +105; `episode_report.go` (daemon) +37 −15, `surface_verb.go`
+103, `daemon.go` +1, `surface_verb_test.go` +187; `surface_verb_outcomes_test.go` +96, `verb-route.golden.json`
+20, `verb_report_e2e_test.go` +127; `app.js` +83 −5, `index.html` +1 −1, `app_test.go` +2 −1, `app_verb_test.go`
+178; `verb_bar_chrome_test.go` +78, `quickstart_verbs_test.go` +55, `QUICKSTART.md` +14. Against the design's
≈ 1,750: M1 came in at 697 (≈ 700), M2 at 208 (≈ 250), M3 at 586 (≈ 500; VR3's store seeding, §1.4 P5), M4 at 418
(≈ 300; the controlled-clock and superseded cases, V1-D1/D2). **If a later re-measure exceeded 2,000**, the cut is the
design's: M4b's XV1 (78) and the navigate arm, then M1b's VB4 file (80, an AST inventory) — never VP2.

---

## 2. How the executor applies a sub-milestone (the patch protocol)

Each sub-milestone `<M>` is the block fenced by a line `` ````diff <M> `` (four backticks) and a line of four
backticks, in Appendix P. Paths are `a/<repo path>` and `b/<repo path>`; a new file's patch has `--- /dev/null`
and an `@@ -0,0 …` hunk. Apply it **only** this way, in one `bash` shell (the controller's shell is zsh):

```bash
bash <<'EOF2'
set -euo pipefail
M=M1a   # <- the sub-milestone, exactly as in the Appendix P heading: M1a M1b M1c M2 M3a M3b M4a M4b
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter395-row243-v1
PLAN="$W/design_docs/planned/w-world-address-bar-verbs-v1-sprint-plan.md"
T="$HOME/.ailang/state/world-iter395/exec-patches"; mkdir -p "$T"
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
- Then run that sub-milestone's **Acceptance** (§3): extract `accept.sh` from this plan and run it, in one bash shell,
  `tee`'d to a log (never only through `tail`). **Never `nohup`, never background a test leg; pass a 600000 ms
  timeout to every long Bash call.**
- Then snapshot (§7). The controller commits and compares the numstat with §3's expected value.

---

## 3. Milestones

**Acceptance.** Every sub-milestone's acceptance is the arm of this script named for it, preceded by the same
**[base]** lines (gofmt, `go build ./...`, `go vet` of the four touched packages, the `_test.go` compile fence). Each
`v` line must show `--- PASS` for every test it names and **no `--- FAIL`, no `--- SKIP`**; each `p` line must end
`ok`. The script ends `ACCEPT_OK <M>` or prints `ACCEPT_RED <M>: <step>` and exits 1. The planner's run of it on every
stage is §1.2 C15. Extract and run it (`M` and `W` as in §2):

```bash
bash <<'EOF2'
set -uo pipefail
M=M1a
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter395-row243-v1
PLAN="$W/design_docs/planned/w-world-address-bar-verbs-v1-sprint-plan.md"
L="$HOME/.ailang/state/world-iter395/exec-logs"; mkdir -p "$L"
awk '$0 == "````bash accept.sh" {f=1; next} f && $0 == "````" {exit} f' "$PLAN" > "$L/accept.sh"
/bin/bash "$L/accept.sh" "$W" "$M" 2>&1 | tee "$L/accept_$M.log"
EOF2
```

````bash accept.sh
#!/bin/bash
# accept.sh <tree> <M>: exactly plan §3's acceptance for sub-milestone M (one bash shell; rc 0 = ACCEPT_OK).
T=$1; M=$2
cd "$T" || exit 2
export AILANG_BIN=$HOME/.pinned-ailang/ailang WORLD_EXEC_NODE=/opt/homebrew/bin/node
unset WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER WORLD_UPDATE_VERB_GOLDEN
CH="/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
set -o pipefail
chk() { tee /dev/stderr | { out=$(cat); echo "$out" | grep -qE -- '--- (FAIL|SKIP)|^FAIL' && exit 1; echo "$out" | grep -q -- '--- PASS' || exit 1; }; }
v() { go test -count=1 -v "$@" 2>&1 | chk; }                                                            # -v: every named test PASS, none SKIP/FAIL
vc() { WORLD_CHROME="$CH" WORLD_REQUIRE_PASSKEY_BROWSER=1 go test -count=1 -v "$@" 2>&1 | chk; }       # the same in live Chrome
p() { go test -count=1 "$@"; }
step() { echo "== $*"; "$@" || { echo "ACCEPT_RED $M: $*"; exit 1; }; }
# [base]
step test -z "$(gofmt -l cmd host)"
step go build ./...
step go vet ./host/verbs/ ./host/surfacepage/ ./host/daemon/ ./cmd/ailang-worldd/
if go test -run '^$' ./... 2>&1 | grep -E '^(FAIL|--- FAIL)'; then echo "ACCEPT_RED $M: FENCE_RED"; exit 1; fi
case $M in
M1a) step v ./host/verbs/ -run '^(TestSplit|TestQuoteRoundTrip)$'
     step p ./host/verbs/ ;;
M1b) step v ./host/verbs/ -run '^(TestLookup|TestApproveTarget|TestSplit|TestQuoteRoundTrip)$'
     step v ./cmd/ailang-worldd/ -run '^(TestVerbTableMatchesTheCLISwitch|TestEveryVerbHasHelp)$'
     step p ./host/verbs/ ;;
M1c) step v ./host/verbs/ -run '^(TestEpisodeReportFlags|TestParseInterleaved|TestLookup)$'
     step p ./host/verbs/
     step p -timeout 30m ./cmd/ailang-worldd/ ;;
M2)  step v ./cmd/ailang-worldd/ -run '^(TestEpisodeReportUsage|TestEpisodeReportTablePrintsEveryClass|TestEpisodeReportTablePrintsTheReadsLine|TestEpisodeReportJSONAndTableLoopback|TestEpisodeReportFollowReprintsOnChangeLoopback|TestEpisodeReportPrintsExecCache|TestEpisodeGlobalAddrRefusedForStartOnly|TestGlobalBareAddrIsNormalisedForEveryClientVerb|TestFinishMarkdownOrderCapAndAC5|TestFinishMarkdownCLIRefusalsLoopback)$'
     step p -timeout 30m ./cmd/ailang-worldd/
     step p ./host/daemon/ -run '^TestEpisodeReport' ;;
M3a) step v ./host/daemon/ -run '^(TestSurfaceVerbRouteRegistered|TestSurfaceVerbBoundsAndErrors|TestSurfaceVerbIsReadOnly|TestEpisodeReportRESTRoute|TestFrozenV1RouteTableMatchesMux|TestSurfaceAppRoutes|TestApprovalRoutesRegistered|TestApprovalRoutesRefuseForeignHost|TestDecisionRoutesRegistered|TestDaemonCredentialWriterSites|TestDaemonDependencyAllowlist)$'
     step p ./host/daemon/ -run '^(TestSurface|TestEpisode|TestApprov|TestDecision|TestDaemon|TestFrozen|TestGETRoutes|TestReadCtx)' ;;   # the whole package is §6's pkgs leg (> 600 s under load)
M3b) step v ./host/daemon/ -run '^TestSurfaceVerbOutcomes$'
     step v ./cmd/ailang-worldd/ -run '^(TestVerbRouteReportIsTheCLIs|TestEpisodeReportJSONAndTableLoopback)$'
     step p ./host/daemon/ -run '^(TestSurfaceVerb|TestEpisodeReport)' ;;
M4a) step v ./host/surfacepage/ -run '^(TestSurfaceAppVerbs|TestSurfaceAppPure|TestNeedsYouListsDecisions|TestSurfaceAppShellAssembly|TestSurfaceAppMatchesWorkingCopy)$'
     step p ./host/surfacepage/
     step v ./host/daemon/ -run '^(TestSurfaceAppRoutes|TestSurfaceVerbOutcomes)$'
     step vc ./cmd/ailang-worldd/ -run '^(TestNeedsYouListsPendingRequests|TestNeedsYouNotifiesOnlyWhenAskedOnce)$' ;;
M4b) step vc ./cmd/ailang-worldd/ -run '^(TestVerbBarEpisodeReportInChrome|TestNeedsYouListsPendingRequests)$'
     step v ./cmd/ailang-worldd/ -run '^(TestQuickstartSection19MatchesTheBar|TestQuickstartSection18FlagsMatchTheCLI|TestQuickstartSection18CardMatchesThePage|TestQuickstartAttendedStepsBoxMatchesSection9)$'
     step p ./cmd/ailang-worldd/ -run '^TestQuickstart'
     step p ./host/daemon/ -run '^(TestQuickstartDecisionPacketSection|TestAGUIQuickstartControl|TestMCPQuickstartPayloadsVerbatim)$'
     step p ./host/broker/ -run '^(TestQuickstartSection10ProfilesLoad|TestQuickstartProtectedContract)$' ;;
*) echo "unknown $M"; exit 2 ;;
esac
echo "ACCEPT_OK $M"
````

### M1a — `Split`, `Quote` · 110 prod + 93 test · Run 1

**Patch** `M1a`: new `host/verbs/split.go` — the package doc; `MaxLine` (1,024), `MaxWords` (32); the five named
refusals `ErrTooLong`, `ErrTooMany`, `ErrNotText`, `ErrUnclosed`, `ErrShell` (each message leads with its name); the
nine unquoted shell characters; `Split` (spaces separate; `'…'` literal; `"…"` literal but `\"`, `\\`; adjacent
segments join; the bound checked first); `Quote` (safe words bare, others single-quoted, `'` as `'"'"'`). New
`host/verbs/split_test.go` — **VB1** `TestSplit` (each refusal at and past each bound: 1,024/1,025 bytes, 32/33
words, `\n` `\t` 0x7f 0x00 0x1f and invalid UTF-8, three unclosed forms, each of the nine characters unquoted refused
and accepted inside `'…'` and `"…"`), **VB2** `TestQuoteRoundTrip` (2,000 seeded argv over an alphabet with every
shell character, quotes, backslashes and spaces; the bare and quoted forms).
**Expected cumulative numstat after the controller's commit of M1a: 203**.

### M1b — the verb table · 133 prod + 174 test · Run 1

**Patch** `M1b`: new `host/verbs/lookup.go` — `Class` (address, read, navigate, mutate, host, never, other, refused),
`Verb{Words, Class, Reason, CLI}`, `Binary`, the three word sets `OperatorWords` (3), `NeverWords` (5),
`OtherWords` (19) with each never verb's reason, `EpisodeNeedsSubcommand`, `AddressHelp` (surface.js's text),
`Lookup` (steps 1, 2, 3, 3b, 4 longest-match, 4b, 5; V1-D8), `ApproveTarget` (`mr-<16 hex>` → `/surface/approve/…`,
`sha256:<64 hex>` → `/surface/decide/…`, anything else refused). New `host/verbs/lookup_test.go` — **VB3**
`TestLookup` (30 rows: every table row and every §2.4 trace row incl. bare `log`, `world session`, `ailang-worldd
world get world-1`, `world get`, `episode`/`episode frobnicate`/`episode w1`, `session frob`, `world a b c`, global
flags; each argv also round-trips through `Quote` (VB2 over VB3); the 27-word shadow set bare, with an argument and
after `ailang-worldd` never `address`, and each reachable as `world NAME`), `TestApproveTarget`. New
`cmd/ailang-worldd/verbs_parity_test.go` — **VB4** `TestVerbTableMatchesTheCLISwitch` (run's switch arms by `go/ast`,
28; arms minus `world` = operator ∪ never ∪ other, each once; three controls fail the check; every word has a help
path via `ownHelpVerbs`, `genericVerbHelp` or `help`; no `run()`).
**Expected cumulative numstat after M1b: 510**.

### M1c — one flag parse: `ParseInterleaved`, `EpisodeReportFlags` · 129 prod + 58 test · Run 1

**Patch** `M1c`: new `host/verbs/flags.go` — `ParseInterleaved` (moved verbatim from `session.go:179-193`),
`EpisodeReportArgs`, `EpisodeReportFlags` (exactly today's five flags, names, defaults, usage strings), the bar
profile `barRefusal` and `AdmitEpisodeReport` (allow `--json`; refuse `--follow`, `--interval`, `--markdown`,
`--addr`, `-h`/`--help`, a parse error; exactly one non-empty episode). The ten call sites call
`verbs.ParseInterleaved` (`approve.go`, `approver.go`, `episode.go`, `episode_finish.go`, `episode_launch.go`,
`episode_report.go`, `preset.go`, `session.go` ×1 (the definition removed), `session_new.go` ×2); `runEpisodeReport`
declares its flags through `verbs.EpisodeReportFlags`. New `host/verbs/flags_test.go` — **VB5**
`TestEpisodeReportFlags` (the five flags and defaults by `VisitAll`; 14 admission rows), `TestParseInterleaved`.
The CLI's tests run whole and unedited.
**Expected cumulative numstat after M1c: 697**.

### M2 — `episode report`'s printer, moved · 208 prod · Run 2

**Patch** `M2`: new `host/daemon/episode_report_text.go` — `WriteEpisodeReportText` (the CLI's `printReport`,
moved), `execCacheLine`, `readsLine`, `WriteReportRefusals`, `NamedCount`, `SortedCounts` (moved; two exported for
P1). `cmd/ailang-worldd/episode_report.go`: `printReport` becomes `daemon.WriteEpisodeReportText(w, r)`; the moved
helpers and the `sort` import leave. `cmd/ailang-worldd/episode_markdown.go` +4 −4 (P1). No test changes: the
report's nine CLI tests and the three markdown tests pass unedited.
**Expected cumulative numstat after M2: 905**.

### M3a — `GET /surface/verb.json` · 156 prod + 187 test · Run 2

**Patch** `M3a`: `host/daemon/episode_report.go` — `handleEpisodeReport` calls the extracted `d.episodeReport`;
`reportBuildError`, `reportTimedOut` (V1-D9). New `host/daemon/surface_verb.go` — `VerbSchema`, the answer,
`handleSurfaceVerb` (400 without `q`; `Split` refusal → 200 `refused`; `Lookup` → address / result / navigate /
terminal / refused) and `verbEpisodeReport` (admission, the 256-byte episode bound, `d.episodeReport` under
`d.readCtx`, the CLI's text or JSON or not-found line; a store error or an expired deadline → `refused` with the `/v1`
text, logged). `host/daemon/daemon.go` +1: `mux.HandleFunc("GET /surface/verb.json",
d.approvalRoute(d.handleSurfaceVerb))`. New `host/daemon/surface_verb_test.go` — **VR1**
`TestSurfaceVerbRouteRegistered` (the one route behind `approvalRoute` in `Handler()`'s source; 421 for three foreign
Hosts; 405 with no CORS header for POST/PUT/DELETE/OPTIONS; 404 around it), **VR4** `TestSurfaceVerbBoundsAndErrors`
(400 without `q`; 200 `TooLong` at 1,025 bytes and a normal answer at 1,024; `Unclosed`; an expired read deadline →
`read deadline (-1ns) exceeded`; a closed store → `internal store failure`, no raw error, logged; the `/v1` route's
500 from the same builder), **VS1** `TestSurfaceVerbIsReadOnly` (host/verbs' import set; `surface_verb.go`'s imports
and its exact `d.` selector set; the id pattern equals `store.MintRequestIDPattern`).
**Expected cumulative numstat after M3a: 1,248**.

### M3b — the proofs: every outcome; the CLI's bytes · 243 test · Run 2

**Patch** `M3b`: new `host/daemon/surface_verb_outcomes_test.go` — **VR2** `TestSurfaceVerbOutcomes` (31 inputs end
to end through `Handler()`: every class and outcome of design §2.4, each with its verb, class, `cli`, `href` and
reason; result members only on `result`; the not-found line and `exit 1`; the hostile report text; the JSON's one
line) and the page's **golden**: 18 real answers written to `host/surfacepage/testdata/verb-route.golden.json` under
`WORLD_UPDATE_VERB_GOLDEN=1`, compared byte for byte otherwise. New `cmd/ailang-worldd/verb_report_e2e_test.go` —
`seedReportCalls` (§1.4 P5) and **VR3** `TestVerbRouteReportIsTheCLIs` (the real CLI and the real route against one
loopback daemon: `w1` and `w2` tables and `--json` byte-identical; not-found: the CLI's stderr and exit code are the
route's `stderr` and `exit`).
**Expected cumulative numstat after M3b: 1,491** (M1–M3 ≤ 1,500: V1-D3).

### M4a — the page · 90 prod + 181 test · Run 3

**Patch** `M4a`: `host/surfacepage/app/app.js` — the pure verb part before the exports (`VERB_ROUTE_BUDGET_MS =
15000`; `verbOutcome(status, text)` with the schema check; `askVerb(raw, io)` with one `AbortController`, one timer,
the fetch and the body raced against the abort, the superseding rule and the post-await check (V1-D1),
`clearTimeout` in `finally`;
`verbAct(act, raw, fx)`; `verbPanel(doc, p)` with `textContent` only), five new exports; `go(raw)` — empty → home
without a fetch, else `askVerb` + `verbAct` with the real effects (`goAddress` is the old `go` body, unchanged);
`showVerbPanel`; Enter keeps the input (V1-D6); Escape closes the panel; the placeholder and one `?` keys row.
`index.html` placeholder +`· episode report <episode>`. `app_test.go` +2 −1 (P2). New
`host/surfacepage/app_verb_test.go` — **`TestSurfaceAppVerbs`** (node; VP1–VP6 over the golden's real answers;
VP2's every failure class × `why`, `log`, `row-189@6`, `world shared@0ffbb9c`, the controlled-clock hang cases and the
superseded cases; spies prove `address`, `assign`, `panel` and `clear` are never called on a failure).
**Expected cumulative numstat after M4a: 1,762**.

### M4b — live Chrome and QUICKSTART §19 · 133 test + 14 docs · Run 3

**Patch** `M4b`: new `cmd/ailang-worldd/verb_bar_chrome_test.go` — **XV1** `TestVerbBarEpisodeReportInChrome` (the A2
driver `passkey_cdp.mjs` and its skip rule, unchanged; a real daemon on a `seedReportCalls` store behind loopback
`httptest`; the real shell at `/surface`: `episode report w1` + Enter → the panel's `<pre>` text equals the CLI's stdout
byte for byte and the bar is cleared; `why` → the panel names `ailang-worldd why`; `log` → no panel, the header names
`ailang-worldd log`, the bar keeps `log`, the location is unchanged). `docs/QUICKSTART.md` +14: **§19 Operator verbs
in the address bar** *(attended — pending first verbatim run)*. New `cmd/ailang-worldd/quickstart_verbs_test.go` —
`TestQuickstartSection19MatchesTheBar` (§19 follows §18; every quoted line classifies as §19 says; the shell's
placeholder names `episode report`).
**Expected cumulative numstat after M4b: 1,909**.

### 3.1 Reachability: every new step, path and refusal, from its entry point (directive retro 4)

Entry points: **BAR** = the app's `#addr` + Enter → `go(raw)` → `askVerb` → `verbAct` (live only in Chrome; node drives
`askVerb`/`verbAct`/`verbPanel` with the real golden bodies); **ROUTE** = `GET /surface/verb.json` through
`d.Handler()` (`approvalRoute` → `handleSurfaceVerb`); **CLI** = `run(argv)`; **UNIT** = a direct call into
`host/verbs`. "E2E" = driven through that entry point.

| Branch (file) | Input class → outcome | Driven |
|---|---|---|
| R2 (approvalRoute) | foreign `Host` → 421 before `q` is read | E2E ROUTE `TestSurfaceVerbRouteRegistered`; the page's 421 → refusal: node VP2 |
| no `q` → 400 | `/surface/verb.json` | E2E ROUTE VR4; page 400 → refusal: node VP2 |
| `Split`: TooLong / TooMany / NotText / Unclosed / Shell | each | UNIT VB1 (all); E2E ROUTE: TooLong, Unclosed (VR4), Shell (VR2); BAR TooLong golden → `say` (node). **TooMany, NotText not E2E**: the page cannot type a control byte into `#addr` (`<input type=text>` strips newlines), and 33 words fit in 1,024 bytes only by typing them — UNIT only |
| step 1 strip | `ailang-worldd …` | UNIT VB3; E2E ROUTE VR2 (`ailang-worldd approve sha256:…`, `ailang-worldd world get world-1`), VR3 n/a |
| step 2 global flag | `--addr x episode report w1` | UNIT VB3; E2E ROUTE VR2 |
| step 3 navigation | `world`, `world NAME`, `world get` | UNIT VB3; E2E ROUTE VR2; BAR golden → `fx.address` (node), Chrome n/a |
| step 3b `world get <ref>` | 3+ words → other, refused with `cli` | UNIT VB3; E2E ROUTE VR2; BAR golden → `say` (node) |
| step 4 rows (read/run, navigate, mutate, host, never, other) | each row | UNIT VB3; E2E ROUTE VR2; BAR: run (Chrome XV1 + node), terminal (Chrome `why` + node), refused (Chrome `log` + node), navigate (node; Chrome n/a — it would load row 201's card, proven by A3a/A4) |
| step 4b `episode` (bare, unknown, `episode w1`), `world a b c` | refused | UNIT VB3; E2E ROUTE VR2; BAR golden `episode`, `world a b c` → `say` (node) |
| step 5 anything else | `row-189@6`, `back`, `World a b c` | UNIT VB3; E2E ROUTE VR2; BAR golden → address path → `addressTarget` gives today's answers (node) |
| admission: `--json` allowed; `--follow`/`--interval`/`--markdown`/`--addr`/`-h`/parse error/positional count refused | each | UNIT VB5 (all); E2E ROUTE VR2 (`--follow`, `--markdown`, `--json`); CLI parses the same set: every CLI report test |
| episode > 256 bytes | refused with the `/v1` text | E2E ROUTE VR2 |
| report found → `text` / `json`; not found → `stderr` + `exit 1` | `w1`, `w1 --json`, `nobody` | E2E ROUTE VR2, VR3 (vs the real CLI); BAR Chrome XV1 (text) + node (all three) |
| report error: deadline / store failure | refused, `/v1` text, logged | E2E ROUTE VR4 |
| `approve` shape: id / ref / other | navigate / navigate / refused | UNIT `TestApproveTarget`; E2E ROUTE VR2 |
| `/v1` extraction: builder timeout → 503; start/finish error → 500; not found → 404 | each | E2E `TestEpisodeReportRESTRoute` (503, 404), VR4 (500 via a closed store) — the start/finish-error-*after-deadline* arm is **not reachable** in a test (the builder fails first under an expired deadline; survivor `E-AS`) |
| page: non-2xx / network error / not JSON / no schema / `v≠1` / unknown or prototype outcome / missing fields / two result members / `stderr` without `exit 1` | `say(...)`, nothing else | node `TestSurfaceAppVerbs` (each) |
| page: budget timeout (headers stalled; body stalled) | `say('…did not answer in time')` | node (controlled clock) |
| page: superseded | nothing at all | node (two cases) |
| page: navigate prefix | `https://evil/x`, `/surface/model.json`, `//evil/…` refused | node VP5 |
| page: panel text | hostile markup as text | node VP3 (fake DOM that throws on `innerHTML`), Chrome XV1 (`textContent` equals the CLI's bytes) |
| page: `go` empty → home, no fetch | `''` | **Chrome n/a** (XV1 types non-empty lines; survivor `JW-GO-EMPTY`) — the branch precedes any fetch and is the HEAD behaviour |

---

## 4. Mutation table

**Procedure** (`tools/mutate.py <run> <tree> <workers> [ids…]`, rows in `tools/mutants.py`): for each row, an APFS
clone of the stage repo at `M4b` (`cp -cR work2`), one exact-string replacement whose anchor occurs exactly once in the
file (else `ANCHOR<n>`, not a run), then `go test -count=1` of the **named test** and of the **behavioural suite**
with **every inventory/static test `-skip`ped** (`TestSurfaceVerbIsReadOnly`, `TestVerbTableMatchesTheCLISwitch`,
`TestSurfaceAppRoutes`, `TestFrozenV1RouteTableMatchesMux`, `TestDaemonDependencyAllowlist`, `TestEveryVerbHasHelp`,
`TestSurfaceAppShellAssembly`, `TestSurfaceAppMatchesWorkingCopy`, `TestQuickstart.*`,
`TestApprovalRoutesRegistered`, `TestDecisionRoutesRegistered`, `TestMarketingChrome.*`) — so **no kill in the
suite column rests on an inventory** (directive retro 3). Suites: `host/verbs`, `host/surfacepage` whole (node with
`WORLD_EXEC_NODE`); `host/daemon` `^(TestSurfaceVerb|TestEpisodeReport|TestSurfaceApp|TestApprovalRoutes)`; `cmd`
`^(TestVerb|TestEpisodeReport|TestEpisodeGlobal|TestGlobalBareAddr|TestApprove|TestSession|TestPreset|TestEpisode)`; a
suite survivor in `host/daemon`/`cmd` is re-run on the **whole package** (`(whole)`). The page-wiring rows (`JW-*`)
run in live Chrome (`WORLD_CHROME`, `WORLD_REQUIRE_PASSKEY_BROWSER=1`) with suite
`^(TestVerbBarEpisodeReportInChrome|TestNeedsYouListsPendingRequests|TestNeedsYouNotifiesOnlyWhenAskedOnce)$`. Go rows
ran with `GOFLAGS=-vet=off`. A build failure is `NOCOMPILE` and **is not a mutant**: the row is rewritten to compile
(`…b`) and re-run, and both lines stay in the ledger. Every run appends one line per row to
`plan-scratch/logs/mutant_ledger.tsv` (time, run, tree, id, file, from → to, named regex, named result, suite result,
failing tests, seconds, expected); per-row output in `logs/mut/<run>-<id>.log`. **The ids below are the ledger's
ids.** Runs: **V1** (`host/verbs` and `app.js` rows on the first stage repo; its `app.js` lines are **void** — that
tree's `TestSurfaceAppVerbs` was red unmutated (§1.4 P12), which is why every later run starts with the unmutated
`CTL-*` controls), **V2** (the three rewritten `…b` rows), **V3** (`app.js` rows after the fix), **V4** (every row,
controls first, on the final stage repo `work2` at `M4b`, byte-identical to the round-trip tree); the table shows each
row's last line. Expected: `K` killed by both; `N` a non-compiling first form (rewritten); `E` named equivalent; `L` named
survivor (reason given). `VR1`'s kill of `R-ROUTE-HOST` is its **421 assertion** (behavioural), recorded from the
row's log, not its source regex.

| id | file | mutation (from → to) | named test | named | suite (whole on a survivor) — failing tests | exp | run |
|---|---|---|---|---|---|---|---|
| `CTL-V` | `split.go` | func Split(line string) → func Split(line string) | `TestSplit\|TestLookup\|TestEpisodeReportFlags` | SURVIVED | SURVIVED — – | S | V4 |
| `CTL-D` | `surface_verb.go` | func (d *Daemon) handleSurfaceVerb( → func (d *Daemon) handleSurfaceVerb( | `TestSurfaceVerbOutcomes\|TestSurfaceVerbBoundsAndErrors\|TestEpisodeReportRESTRoute` | SURVIVED | SURVIVED(whole) — – | S | V4 |
| `CTL-C` | `surface_verb.go` | func (d *Daemon) handleSurfaceVerb( → func (d *Daemon) handleSurfaceVerb( | `TestVerbRouteReportIsTheCLIs\|TestEpisodeReportTablePrintsEveryClass` | SURVIVED | SURVIVED(whole) — – | S | V4 |
| `CTL-SP` | `app.js` | const VERB_ROUTE_BUDGET_MS = 15000; → const VERB_ROUTE_BUDGET_MS = 15000; | `TestSurfaceAppVerbs` | SURVIVED | SURVIVED — – | S | V4 |
| `CTL-CH` | `app.js` | const VERB_ROUTE_BUDGET_MS = 15000; → const VERB_ROUTE_BUDGET_MS = 15000; | `TestVerbBarEpisodeReportInChrome` | SURVIVED | SURVIVED — – | S | V4 |
| `S-UTF8b` | `split.go` | if !utf8.ValidString(line) { → if !utf8.ValidString(line) && false { | `TestSplit` | KILLED | KILLED — TestSplit | K | V4 |
| `S-LONG` | `split.go` | if len(line) > MaxLine { → if len(line) > MaxLine+1 { | `TestSplit` | KILLED | KILLED — TestSplit | K | V4 |
| `S-UTF8` | `split.go` | if !utf8.ValidString(line) { → if false { | `TestSplit` | NOCOMPILE | NOCOMPILE — – | N | V4 |
| `S-CTRL` | `split.go` | if line[i] < 0x20 \|\| line[i] == 0x7f { → if line[i] < 0x09 \|\| line[i] == 0x7f { | `TestSplit` | KILLED | KILLED — TestSplit | K | V4 |
| `S-DEL` | `split.go` | if line[i] < 0x20 \|\| line[i] == 0x7f { → if line[i] < 0x20 { | `TestSplit` | KILLED | KILLED — TestSplit | K | V4 |
| `S-SPACE` | `split.go` | case c == ' ': → case c == '\t': | `TestSplit` | KILLED | KILLED — TestLookup,TestQuoteRoundTrip,TestSplit | K | V4 |
| `S-SQ-OPEN` | `split.go` | if j < 0 {⏎return nil, ErrUnclosed → if j < -1 {⏎return nil, ErrUnclosed | `TestSplit` | KILLED | KILLED — TestSplit | K | V4 |
| `S-DQ-ESC-BS` | `split.go` | (line[j+1] == '"' \|\| line[j+1] == '\\') → (line[j+1] == '"') | `TestSplit` | KILLED | KILLED — TestSplit | K | V4 |
| `S-DQ-ESC-Q` | `split.go` | (line[j+1] == '"' \|\| line[j+1] == '\\') → (line[j+1] == '\\') | `TestSplit` | KILLED | KILLED — TestSplit | K | V4 |
| `S-DQ-OPEN` | `split.go` | if j == len(line) { → if j > len(line) { | `TestSplit` | KILLED | KILLED — TestSplit | K | V4 |
| `S-SHELL` | `split.go` | case strings.IndexByte(shellChars, c) >= 0: → case false: | `TestSplit` | KILLED | KILLED — TestSplit | K | V4 |
| `S-SHELL-SET` | `split.go` | const shellChars = "$';\|&<>()" → const shellChars = "$';\|&<>)" | `TestSplit` | KILLED | KILLED — TestSplit | K | V4 |
| `S-FLUSH` | `split.go` | if in {⏎words = append(words, w.String())⏎}⏎if len(words) → if len(words) | `TestSplit` | KILLED | KILLED — TestLookup,TestQuoteRoundTrip,TestSplit | K | V4 |
| `S-EMPTYWORD` | `split.go` | in = true⏎} → in = w.Len() > 0⏎} | `TestSplit` | KILLED | KILLED — TestQuoteRoundTrip,TestSplit | K | V4 |
| `S-TOOMANY` | `split.go` | if len(words) > MaxWords { → if len(words) > MaxWords+1 { | `TestSplit` | KILLED | KILLED — TestSplit | K | V4 |
| `Q-SAFE` | `split.go` | '^[A-Za-z0-9@%+=:,./_-]+$' → '^[A-Za-z0-9@%+=:,./_ -]+$' | `TestQuoteRoundTrip` | KILLED | KILLED — TestQuoteRoundTrip | K | V4 |
| `Q-EMPTY` | `split.go` | '^[A-Za-z0-9@%+=:,./_-]+$' → '^[A-Za-z0-9@%+=:,./_-]*$' | `TestQuoteRoundTrip` | KILLED | KILLED — TestQuoteRoundTrip | K | V4 |
| `Q-ESC` | `split.go` | strings.ReplaceAll(a, "'", ''"'"'') → a | `TestQuoteRoundTrip` | KILLED | KILLED — TestQuoteRoundTrip | K | V4 |
| `L-STRIP` | `lookup.go` | if len(argv) > 0 && argv[0] == Binary { // 1 → if false { // 1 | `TestLookup` | KILLED | KILLED — TestLookup | K | V4 |
| `L-EMPTY` | `lookup.go` | if len(argv) == 0 {⏎return Verb{Class: Address}, nil → if len(argv) == 0 {⏎return Verb{Class: Refused}, nil | `TestLookup` | KILLED | KILLED — TestLookup | K | V4 |
| `L-FLAG` | `lookup.go` | if strings.HasPrefix(first, "-") { // 2 → if strings.HasPrefix(first, "--") { // 2 | `TestLookup` | KILLED | KILLED — TestLookup | K | V4 |
| `L-FLAG-ORDER` | `lookup.go` | if strings.HasPrefix(first, "-") { // 2 → if strings.HasPrefix(first, "-") && len(argv) == 1 { // 2 | `TestLookup` | KILLED | KILLED — TestLookup | K | V4 |
| `L-NAV-WIDE` | `lookup.go` | if first == "world" && len(argv) <= 2 { // 3 → if first == "world" && len(argv) <= 3 { // 3 | `TestLookup` | KILLED | KILLED — TestLookup | K | V4 |
| `L-NAV-NARROW` | `lookup.go` | if first == "world" && len(argv) <= 2 { // 3 → if first == "world" && len(argv) <= 1 { // 3 | `TestLookup` | KILLED | KILLED — TestLookup | K | V4 |
| `L-3B` | `lookup.go` | if first == "world" && argv[1] == "get" { // 3b → if false { // 3b | `TestLookup` | KILLED | KILLED — TestLookup | K | V4 |
| `L-EP-REPORT` | `lookup.go` | case "report": → case "reportx": | `TestLookup` | KILLED | KILLED — TestLookup | K | V4 |
| `L-EP-FINISH` | `lookup.go` | case "start", "finish": → case "start": | `TestLookup` | KILLED | KILLED — TestLookup | K | V4 |
| `L-EP-LAUNCH` | `lookup.go` | case "launch": → case "launchx": | `TestLookup` | KILLED | KILLED — TestLookup | K | V4 |
| `L-WHY` | `lookup.go` | case first == "why": → case first == "whyx": | `TestLookup` | KILLED | KILLED — TestLookup | K | V4 |
| `L-APPROVE` | `lookup.go` | case first == "approve": → case first == "approvex": | `TestLookup` | KILLED | KILLED — TestLookup | K | V4 |
| `L-NEVER` | `lookup.go` | case member(NeverWords, first): → case false: | `TestLookup` | KILLED | KILLED — TestLookup | K | V4 |
| `L-OTHER` | `lookup.go` | case member(OtherWords, first): → case false: | `TestLookup` | KILLED | KILLED — TestLookup | K | V4 |
| `L-MEMBER` | `lookup.go` | if s == w { → if s != w { | `TestLookup` | KILLED | KILLED — TestLookup | K | V4 |
| `L-SET-LOG` | `lookup.go` | "hosts", "log", "object" → "hosts", "object" | `TestLookup` | KILLED | KILLED — TestLookup | K | V4 |
| `L-4B-EPISODE` | `lookup.go` | case "episode":⏎return Verb{Words: "episode", Class: Refused → case "episodex":⏎return Verb{Words: "episode", Class: Refuse | `TestLookup` | KILLED | KILLED — TestLookup | K | V4 |
| `L-4B-WORLD` | `lookup.go` | case "world":⏎return Verb{Words: "world", Class: Refused → case "worldx":⏎return Verb{Words: "world", Class: Refused | `TestLookup` | KILLED | KILLED — TestLookup | K | V4 |
| `L-CLI` | `lookup.go` | cli := Quote(append([]string{Binary}, argv...)) → cli := Quote(argv) | `TestLookup` | KILLED | KILLED — TestLookup | K | V4 |
| `A-ONE` | `lookup.go` | if len(args) == 1 && mintRequestID.MatchString(args[0]) { → if len(args) >= 1 && mintRequestID.MatchString(args[0]) { | `TestApproveTarget` | KILLED | KILLED — TestApproveTarget | K | V4 |
| `A-MR-PAT` | `lookup.go` | '^mr-[0-9a-f]{16}$' → '^mr-[0-9a-fA-F]{16}$' | `TestApproveTarget` | KILLED | KILLED — TestApproveTarget | K | V4 |
| `A-REF` | `lookup.go` | if len(args) == 1 && decisionRef.MatchString(args[0]) { → if len(args) == 1 && strings.HasPrefix(args[0], "sha256:") { | `TestApproveTarget` | KILLED | KILLED — TestApproveTarget | K | V4 |
| `F-PI` | `flags.go` | if fs.NArg() == 0 { → if fs.NArg() <= 1 { | `TestParseInterleaved` | KILLED | KILLED — TestEpisodeReportFlags,TestParseInterleaved | K | V4 |
| `F-DEFAULT` | `flags.go` | 2*time.Second → 3*time.Second | `TestEpisodeReportFlags` | KILLED | KILLED — TestEpisodeReportFlags | K | V4 |
| `F-HELP` | `flags.go` | if errors.Is(err, flag.ErrHelp) { → if false { | `TestEpisodeReportFlags` | NOCOMPILE | NOCOMPILE — – | N | V4 |
| `F-HELPb` | `flags.go` | if errors.Is(err, flag.ErrHelp) { → if errors.Is(err, flag.ErrHelp) && false { | `TestEpisodeReportFlags` | KILLED | KILLED — TestEpisodeReportFlags | K | V4 |
| `F-ERR` | `flags.go` | if err != nil {⏎return "", false, Binary → if false {⏎return "", false, Binary | `TestEpisodeReportFlags` | KILLED | KILLED — TestEpisodeReportFlags | K | V4 |
| `F-FOLLOW` | `flags.go` | {"follow", "the app shows one report; it re-reads on Enter"} → (deleted) | `TestEpisodeReportFlags` | KILLED | KILLED — TestEpisodeReportFlags | K | V4 |
| `F-INTERVAL` | `flags.go` | {"interval", "the app shows one report; it re-reads on Enter → (deleted) | `TestEpisodeReportFlags` | KILLED | KILLED — TestEpisodeReportFlags | K | V4 |
| `F-MARKDOWN` | `flags.go` | {"markdown", "the PR body is rendered by episode finish"},⏎ → (deleted) | `TestEpisodeReportFlags` | KILLED | KILLED — TestEpisodeReportFlags | K | V4 |
| `F-ADDR` | `flags.go` | {"addr", "the app reads its own daemon"},⏎ → (deleted) | `TestEpisodeReportFlags` | KILLED | KILLED — TestEpisodeReportFlags | K | V4 |
| `F-VISIT` | `flags.go` | set[f.Name] = true → set[f.Name] = false | `TestEpisodeReportFlags` | KILLED | KILLED — TestEpisodeReportFlags | K | V4 |
| `F-POS-COUNT` | `flags.go` | if len(pos) != 1 \|\| pos[0] == "" { → if len(pos) < 1 \|\| pos[0] == "" { | `TestEpisodeReportFlags` | KILLED | KILLED — TestEpisodeReportFlags | K | V4 |
| `F-POS-EMPTY` | `flags.go` | if len(pos) != 1 \|\| pos[0] == "" { → if len(pos) != 1 { | `TestEpisodeReportFlags` | KILLED | KILLED — TestEpisodeReportFlags | K | V4 |
| `F-JSON` | `flags.go` | return pos[0], *a.JSON, "" → return pos[0], false, "" | `TestEpisodeReportFlags` | NOCOMPILE | NOCOMPILE — – | N | V4 |
| `F-JSONb` | `flags.go` | return pos[0], *a.JSON, "" → return pos[0], *a.JSON && false, "" | `TestEpisodeReportFlags` | KILLED | KILLED — TestEpisodeReportFlags | K | V4 |
| `F-PI-CLI` | `flags.go` | if fs.NArg() == 0 { → if fs.NArg() <= 1 { | `TestEpisodeReportJSONAndTableLoopback` | KILLED | KILLED — TestApproverRevokeListsWhatItApproved | K | V4 |
| `P-TOOL-WIDTH` | `episode_report_text.go` | "  %-18s %6d\n", kv.Name, kv.N → "  %-17s %6d\n", kv.Name, kv.N | `TestEpisodeReportTablePrintsEveryClass` | KILLED | KILLED — TestEpisodeReportTablePrintsEveryClass,TestEpisodeReportTablePrintsTheReadsLine | K | V4 |
| `P-SORT` | `episode_report_text.go` | return out[i].N > out[j].N → return out[i].N < out[j].N | `TestEpisodeReportTablePrintsTheReadsLine` | KILLED | KILLED — TestEpisodeReportTablePrintsTheReadsLine | K | V5 |
| `P-READS` | `episode_report_text.go` | line += fmt.Sprintf(", %d outside the worktree", r.Unlisted) → line += fmt.Sprintf(", %d outside", r.Unlisted) | `TestEpisodeReportTablePrintsTheReadsLine` | KILLED | KILLED — TestEpisodeReportTablePrintsTheReadsLine | K | V4 |
| `E-FOUND` | `episode_report.go` | return rep, ok \|\| start != nil \|\| rep.Finish != nil \|\| len(r → return rep, ok, nil | `TestEpisodeReportShowsBindSeedFailure` | KILLED | KILLED — TestEpisodeReportShowsBindSeedFailure | K | V5 |
| `E-WRAP` | `episode_report.go` | return rep, false, reportBuildError{err} → return rep, false, err | `TestEpisodeReportRESTRoute` | KILLED | KILLED — TestEpisodeReportRESTRoute,TestSurfaceVerbBoundsAndErrors | K | V4 |
| `E-TIMEOUT` | `episode_report.go` | if reportTimedOut(ctx, err) {⏎writeReadTimeout(w, d.readDead → if false {⏎writeReadTimeout(w, d.readDeadline) | `TestEpisodeReportRESTRoute` | KILLED | KILLED — TestEpisodeReportRESTRoute | K | V4 |
| `E-AS` | `episode_report.go` | return errors.As(err, &build) && timedOut(ctx, build.err) → return err != nil && timedOut(ctx, err) | `TestEpisodeReportRESTRoute` | NOCOMPILE | NOCOMPILE — – | N | V4 |
| `E-ASb` | `episode_report.go` | return errors.As(err, &build) && timedOut(ctx, build.err) → return (errors.As(err, &build) \|\| err != nil) && timedOut(ct | `TestEpisodeReportRESTRoute` | SURVIVED | SURVIVED(whole) — – | L | V5 |
| `E-NOTFOUND` | `episode_report.go` | if !found {⏎writeAPIError(w, "NotFound" → if false {⏎writeAPIError(w, "NotFound" | `TestEpisodeReportRESTRoute` | NOCOMPILE | NOCOMPILE — – | N | V4 |
| `E-NOTFOUNDb` | `episode_report.go` | if !found {⏎writeAPIError(w, "NotFound" → if !found && false {⏎writeAPIError(w, "NotFound" | `TestEpisodeReportRESTRoute` | KILLED | KILLED — TestEpisodeReportRESTRoute | K | V5 |
| `R-Q` | `surface_verb.go` | if !query.Has("q") { → if false { | `TestSurfaceVerbBoundsAndErrors` | KILLED | KILLED — TestSurfaceVerbBoundsAndErrors | K | V4 |
| `R-SPLIT` | `surface_verb.go` | if err != nil {⏎a.Reason = err.Error() → if false {⏎a.Reason = err.Error() | `TestSurfaceVerbBoundsAndErrors` | KILLED | KILLED — TestSurfaceVerbBoundsAndErrors,TestSurfaceVerbOutcomes | K | V4 |
| `R-ADDRESS` | `surface_verb.go` | a.Outcome = "address" → a.Outcome = "refused" | `TestSurfaceVerbOutcomes` | KILLED | KILLED — TestSurfaceVerbBoundsAndErrors,TestSurfaceVerbOutcomes | K | V4 |
| `R-REPORT` | `surface_verb.go` | case v.Words == "episode report": → case v.Words == "episode reportx": | `TestSurfaceVerbOutcomes` | KILLED | KILLED — TestSurfaceVerbBoundsAndErrors,TestSurfaceVerbOutcomes | K | V4 |
| `R-NAV` | `surface_verb.go` | ; a.Href != "" { → ; true { | `TestSurfaceVerbOutcomes` | KILLED | KILLED — TestSurfaceVerbOutcomes | K | V4 |
| `R-TERM-HOST` | `surface_verb.go` | v.Class == verbs.Read \|\| v.Class == verbs.Mutate \|\| v.Class  → v.Class == verbs.Read \|\| v.Class == verbs.Mutate: | `TestSurfaceVerbOutcomes` | KILLED | KILLED — TestSurfaceVerbOutcomes | K | V4 |
| `R-TERM-READ` | `surface_verb.go` | v.Class == verbs.Read \|\| v.Class == verbs.Mutate \|\| v.Class  → v.Class == verbs.Mutate \|\| v.Class == verbs.Host: | `TestSurfaceVerbOutcomes` | KILLED | KILLED — TestSurfaceVerbBoundsAndErrors,TestSurfaceVerbOutcomes | K | V4 |
| `R-DEFAULT-OUTCOME` | `surface_verb.go` | Class: string(verbs.Refused), Outcome: "refused"} → Class: string(verbs.Refused), Outcome: "address"} | `TestSurfaceVerbOutcomes` | KILLED | KILLED — TestSurfaceVerbBoundsAndErrors,TestSurfaceVerbOutcomes | K | V4 |
| `R-ADMIT` | `surface_verb.go` | if reason != "" {⏎a.Reason = reason → if false {⏎a.Reason = reason | `TestSurfaceVerbOutcomes` | KILLED | KILLED — TestSurfaceVerbOutcomes | K | V4 |
| `R-EPBOUND` | `surface_verb.go` | if len(episode) > maxReportEpisodeBytes { → if len(episode) > maxReportEpisodeBytes+1 { | `TestSurfaceVerbOutcomes` | KILLED | KILLED — TestSurfaceVerbOutcomes | K | V4 |
| `R-ERR-TIMEOUT` | `surface_verb.go` | if reportTimedOut(ctx, err) {⏎a.Reason → if false {⏎a.Reason | `TestSurfaceVerbBoundsAndErrors` | KILLED | KILLED — TestSurfaceVerbBoundsAndErrors | K | V4 |
| `R-ERR-LEAK` | `surface_verb.go` | a.Reason = internalErrorMessage⏎ → a.Reason = err.Error()⏎ | `TestSurfaceVerbBoundsAndErrors` | KILLED | KILLED — TestSurfaceVerbBoundsAndErrors | K | V4 |
| `R-ERR-LOG` | `surface_verb.go` | d.writeInternalErrorLog(r, err)⏎ → (deleted) | `TestSurfaceVerbBoundsAndErrors` | KILLED | KILLED — TestSurfaceVerbBoundsAndErrors | K | V4 |
| `R-NOTFOUND` | `surface_verb.go` | case !found:⏎ → case false:⏎ | `TestSurfaceVerbOutcomes` | NOCOMPILE | NOCOMPILE — – | N | V4 |
| `R-NOTFOUNDb` | `surface_verb.go` | case !found:⏎ → case !found && false:⏎ | `TestSurfaceVerbOutcomes` | KILLED | KILLED — TestSurfaceVerbOutcomes | K | V5 |
| `R-JSON` | `surface_verb.go` | case asJSON:⏎ → case false:⏎ | `TestSurfaceVerbOutcomes` | NOCOMPILE | NOCOMPILE — – | N | V4 |
| `R-JSONb` | `surface_verb.go` | case asJSON:⏎ → case asJSON && false:⏎ | `TestSurfaceVerbOutcomes` | KILLED | KILLED — TestSurfaceVerbOutcomes | K | V5 |
| `R-JSON-INDENT` | `surface_verb.go` | _ = json.NewEncoder(&b).Encode(rep) → e := json.NewEncoder(&b); e.SetIndent("", " "); _ = e.Encode | `TestVerbRouteReportIsTheCLIs` | KILLED | KILLED — TestVerbRouteReportIsTheCLIs | K | V4 |
| `R-HEADERS` | `surface_verb.go` | defer func() { writeApprovalJSON(w, http.StatusOK, a) }() → defer func() { writeJSON(w, http.StatusOK, a) }() | `TestSurfaceVerbOutcomes` | KILLED | KILLED — TestSurfaceVerbBoundsAndErrors,TestSurfaceVerbOutcomes | K | V4 |
| `R-ARGV` | `surface_verb.go` | Argv: []string{}, Class: → Argv: nil, Class: | `TestSurfaceVerbBoundsAndErrors` | KILLED | KILLED — TestSurfaceVerbBoundsAndErrors,TestSurfaceVerbOutcomes | K | V4 |
| `R-ROUTE-HOST` | `daemon.go` | mux.HandleFunc("GET /surface/verb.json", d.approvalRoute(d.h → mux.HandleFunc("GET /surface/verb.json", func(w http.Respons | `TestSurfaceVerbRouteRegistered` | KILLED | KILLED — TestSurfaceVerbRouteRegistered | K | V4 |
| `J-404-FALLBACK` | `app.js` | if (status !== 200) return { say: → if (status === 404) return { address: true }; if (status !== | `TestSurfaceAppVerbs` | KILLED | KILLED — TestSurfaceAppVerbs | K | V4 |
| `J-ANY-FALLBACK` | `app.js` |     fx.say(act.say);⏎  } →     fx.clear(); return fx.address(raw);⏎  } | `TestSurfaceAppVerbs` | KILLED | KILLED — TestSurfaceAppVerbs | K | V4 |
| `J-SCHEMA` | `app.js` | b.schema === 'world/surface-verb/v1' &&  → (deleted) | `TestSurfaceAppVerbs` | KILLED | KILLED — TestSurfaceAppVerbs | K | V4 |
| `J-V` | `app.js` | b.v === 1 &&  → (deleted) | `TestSurfaceAppVerbs` | KILLED | KILLED — TestSurfaceAppVerbs | K | V4 |
| `J-OWN-OUTCOME` | `app.js` | Object.prototype.hasOwnProperty.call(VERB_NEEDS, b.outcome) → b.outcome in VERB_NEEDS | `TestSurfaceAppVerbs` | KILLED | KILLED — TestSurfaceAppVerbs | K | V4 |
| `J-FIELDS` | `app.js` | if (!need \|\| !need.every(k => typeof b[k] === 'string' && b[ → if (!need) return UNREADABLE; | `TestSurfaceAppVerbs` | KILLED | KILLED — TestSurfaceAppVerbs | K | V4 |
| `J-NOTJSON` | `app.js` | try { b = JSON.parse(text); } catch (_) { return UNREADABLE; → try { b = JSON.parse(text); } catch (_) { return { address:  | `TestSurfaceAppVerbs` | KILLED | KILLED — TestSurfaceAppVerbs | K | V4 |
| `J-NAV-PREFIX` | `app.js` | b.href.startsWith('/surface/approve/mr-') \|\| b.href.startsWi → true | `TestSurfaceAppVerbs` | KILLED | KILLED — TestSurfaceAppVerbs | K | V4 |
| `J-TERMINAL` | `app.js` | if (b.outcome === 'terminal') return { panel: → if (b.outcome === 'terminal') return { say: b.reason, panel: | `TestSurfaceAppVerbs` | KILLED | KILLED — TestSurfaceAppVerbs | K | V4 |
| `J-RESULT-ONE` | `app.js` | if (shown.length !== 1 \|\| (b.stderr !== undefined && b.exit  → if (shown.length < 1) return UNREADABLE; | `TestSurfaceAppVerbs` | KILLED | KILLED — TestSurfaceAppVerbs | K | V4 |
| `J-SIGNAL` | `app.js` | { signal: job.ctl.signal, headers: → { headers: | `TestSurfaceAppVerbs` | KILLED | KILLED — TestSurfaceAppVerbs | K | V4 |
| `J-LATE` | `app.js` | io.setTimeout(() => { job.late = true; job.ctl.abort(); } → io.setTimeout(() => { job.ctl.abort(); } | `TestSurfaceAppVerbs` | KILLED | KILLED — TestSurfaceAppVerbs | K | V4 |
| `J-BUDGET` | `app.js` | const VERB_ROUTE_BUDGET_MS = 15000; → const VERB_ROUTE_BUDGET_MS = 1500000; | `TestSurfaceAppVerbs` | KILLED | KILLED — TestSurfaceAppVerbs | K | V4 |
| `J-BUDGET-USE` | `app.js` | }, VERB_ROUTE_BUDGET_MS); → }, 60000); | `TestSurfaceAppVerbs` | KILLED | KILLED — TestSurfaceAppVerbs | K | V4 |
| `J-BODY-RACE` | `app.js` | const text = await Promise.race([res.text(), stop]); → const text = await res.text(); | `TestSurfaceAppVerbs` | KILLED | KILLED — TestSurfaceAppVerbs | K | V4 |
| `J-SUPERSEDE-CHECK` | `app.js` | return job.superseded ? { superseded: true } : verbOutcome(r → return verbOutcome(res.status, text); | `TestSurfaceAppVerbs` | KILLED | KILLED — TestSurfaceAppVerbs | K | V4 |
| `J-CLEAR` | `app.js` |       io.clearTimeout(timer);⏎ → (deleted) | `TestSurfaceAppVerbs` | KILLED | KILLED — TestSurfaceAppVerbs | K | V4 |
| `J-SUPERSEDE` | `app.js` | if (verbJob) { verbJob.superseded = true; verbJob.ctl.abort( → (deleted) | `TestSurfaceAppVerbs` | KILLED | KILLED — TestSurfaceAppVerbs | K | V4 |
| `J-CATCH-SUPERSEDED` | `app.js` |       if (job.superseded) return { superseded: true };⏎ → (deleted) | `TestSurfaceAppVerbs` | KILLED | KILLED — TestSurfaceAppVerbs | K | V4 |
| `J-ACT-SUPERSEDED` | `app.js` |     if (act.superseded) return;⏎ → (deleted) | `TestSurfaceAppVerbs` | KILLED | KILLED — TestSurfaceAppVerbs | K | V4 |
| `J-ACT-CLEAR` | `app.js` | if (act.address) { fx.clear(); return fx.address(raw); } → if (act.address) { return fx.address(raw); } | `TestSurfaceAppVerbs` | KILLED | KILLED — TestSurfaceAppVerbs | K | V4 |
| `J-REFUSE-CLEARS` | `app.js` |     fx.say(act.say);⏎  } →     fx.clear(); fx.say(act.say);⏎  } | `TestSurfaceAppVerbs` | KILLED | KILLED — TestSurfaceAppVerbs | K | V4 |
| `J-PANEL-TEXT` | `app.js` | body.textContent = p.body; → body.innerHTML = p.body; | `TestSurfaceAppVerbs` | KILLED | KILLED — TestSurfaceAppVerbs | K | V4 |
| `J-PANEL-HEAD` | `app.js` | head.textContent = p.head; → head.innerHTML = p.head; | `TestSurfaceAppVerbs` | KILLED | KILLED — TestSurfaceAppVerbs | K | V4 |
| `J-URL` | `app.js` | '/surface/verb.json?q=' + encodeURIComponent(raw) → '/surface/verb.json?q=' + raw | `TestSurfaceAppVerbs` | KILLED | KILLED — TestSurfaceAppVerbs | K | V4 |
| `JW-GO-ROUTE` | `app.js` |     const act = await askVerb(raw, →     if (true) return goAddress(raw);⏎    const act = await a | `TestVerbBarEpisodeReportInChrome` | KILLED | KILLED — TestVerbBarEpisodeReportInChrome | K | V4 |
| `JW-ENTER-CLEARS` | `app.js` | if (e.key === 'Enter') { go(e.target.value); $('root') → if (e.key === 'Enter') { go(e.target.value); e.target.value  | `TestVerbBarEpisodeReportInChrome` | KILLED | KILLED — TestVerbBarEpisodeReportInChrome | K | V4 |
| `JW-PANEL-TEXT` | `app.js` | const old = $('verb'), el = verbPanel(document, p); → const old = $('verb'), el = verbPanel(document, p); el.query | `TestVerbBarEpisodeReportInChrome` | KILLED | KILLED — TestVerbBarEpisodeReportInChrome | K | V4 |
| `JW-GO-EMPTY` | `app.js` |     if (String(raw == null ? '' : raw).trim() === '') return → (deleted) | `TestVerbBarEpisodeReportInChrome` | SURVIVED | SURVIVED — – | L | V4 |

---

## 5. Files (final `git status --short --untracked-files=all` after M4b)

```
 M cmd/ailang-worldd/approve.go
 M cmd/ailang-worldd/approver.go
 M cmd/ailang-worldd/episode.go
 M cmd/ailang-worldd/episode_finish.go
 M cmd/ailang-worldd/episode_launch.go
 M cmd/ailang-worldd/episode_markdown.go
 M cmd/ailang-worldd/episode_report.go
 M cmd/ailang-worldd/preset.go
 M cmd/ailang-worldd/session.go
 M cmd/ailang-worldd/session_new.go
 M docs/QUICKSTART.md
 M host/daemon/daemon.go
 M host/daemon/episode_report.go
 M host/surfacepage/app/app.js
 M host/surfacepage/app/index.html
 M host/surfacepage/app_test.go
?? cmd/ailang-worldd/quickstart_verbs_test.go
?? cmd/ailang-worldd/verb_bar_chrome_test.go
?? cmd/ailang-worldd/verb_report_e2e_test.go
?? cmd/ailang-worldd/verbs_parity_test.go
?? host/daemon/episode_report_text.go
?? host/daemon/surface_verb.go
?? host/daemon/surface_verb_outcomes_test.go
?? host/daemon/surface_verb_test.go
?? host/surfacepage/app_verb_test.go
?? host/surfacepage/testdata/verb-route.golden.json
?? host/verbs/flags.go
?? host/verbs/flags_test.go
?? host/verbs/lookup.go
?? host/verbs/lookup_test.go
?? host/verbs/split.go
?? host/verbs/split_test.go
```
Plus the two plan files (`design_docs/planned/w-world-address-bar-verbs-v1-sprint-plan.md`,
`design_docs/planned/sprint_w-world-address-bar-verbs-v1.json`), which the controller commits. Nothing under
`world/`, `packages/`, `scripts/`, `tools/`, `host/store/`, `host/broker/`, `host/surfacepage/{approve,passkey,decide}/`,
`host/surfacepage/surface.js`.

---

## 6. Sprint gates (controller, after the last sub-milestone's commit; each measured at the base, §1.1)

Run it in **one bash shell, in the foreground, never under `nohup`** (macOS `/bin/bash` 3.2; no `declare -A`; every
leg records its own rc and every leg runs — no `set -e`; every leg is `tee`'d), after `mkdir -p
$HOME/.ailang/state/world-iter395/gates`, extracting the script from this plan:

```bash
bash <<'EOF2'
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter395-row243-v1
G="$HOME/.ailang/state/world-iter395/gates"; mkdir -p "$G"
awk '$0 == "````bash gates.sh" {f=1; next} f && $0 == "````" {exit} f' "$W/design_docs/planned/w-world-address-bar-verbs-v1-sprint-plan.md" > "$G/gates.sh"
/bin/bash "$G/gates.sh" "$W" "$G" final "$(git -C "$W" rev-parse HEAD)"
EOF2
```

It writes `gates/legs.txt`, one log per leg, and **`gates/ALL_DONE.<sha>`** when every leg has finished (the judge
waits on it). The `verify_go` leg (`VERIFY_GO_LEGS=plain AILANG_BIN=$HOME/.pinned-ailang/ailang
./scripts/verify_go.sh`) reads committed blobs, so the script runs **after** the last commit. The legs: `status`,
`verify_ail` (no `.ail` changed: the base's result byte for byte, §1.1), `gofmt`, `build`, `vet`, `fence` (`go test
-run '^$' ./...`, the `_test.go` compile fence), `pkgs` (the four touched packages whole), `race` (`-race` on
`host/verbs`, `host/surfacepage` whole and the touched tests of `host/daemon` and `cmd`), `inv` (every named inventory
`-v` + every new test: none missing, none SKIP), `chrome` (the twelve A2/A3/A4/Needs-you page tests + XV1 in `cmd`,
and the two `host/daemon` app-shell Chrome drills), `verify_go`. The script, verbatim:

````bash gates.sh
#!/bin/bash
# Row 243 V1 gates (plan §6). usage: gates.sh <tree> <logdir> <base|final> [<sha>]
# The controller: bash gates.sh <worktree> $HOME/.ailang/state/world-iter395/gates final <HEAD sha>
# macOS /bin/bash 3.2: no declare -A. Every leg records its own rc; every leg runs (no set -e).
# Never start this under nohup. Writes <logdir>/ALL_DONE.<sha> when every leg has finished.
T=$1; L=$2; MODE=$3
mkdir -p "$L"; cd "$T" || exit 2
SHA=${4:-$(git rev-parse HEAD 2>/dev/null || echo unknown)}
export AILANG_BIN=$HOME/.pinned-ailang/ailang
export WORLD_EXEC_NODE=/opt/homebrew/bin/node
unset WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER WORLD_PASSKEY_RECORD_DIR WORLD_UPDATE_VERB_GOLDEN
: > "$L/legs.txt"
FAILS=""
s0=$(date +%s)
leg() { # leg <name> <command string>: tee'd to $L/<name>.log, rc recorded in legs.txt
  local name=$1 cmd=$2 s rc
  s=$(date +%s)
  bash -c "$cmd" 2>&1 | tee "$L/$name.log"
  rc=${PIPESTATUS[0]}
  echo "LEG $name rc=$rc $(( $(date +%s) - s ))s" | tee -a "$L/legs.txt"
  [ "$rc" -eq 0 ] || FAILS="$FAILS $name"
}
leg status 'git status --porcelain | wc -l'
leg verify_ail './scripts/verify_ail.sh'
leg gofmt 'out=$(gofmt -l cmd host); echo "$out"; test -z "$out"'
leg build 'go build ./...'
leg vet 'go vet ./...'
leg fence "go test -run '^\$' ./... > '$L/fence.raw' 2>&1; rc=\$?; cat '$L/fence.raw'; if grep -E '^(FAIL|--- FAIL)' '$L/fence.raw' >/dev/null; then exit 1; fi; exit \$rc"
PK="./host/surfacepage/ ./host/daemon/ ./cmd/ailang-worldd/"
[ "$MODE" = base ] || PK="./host/verbs/ $PK"
leg pkgs "go test -count=1 -timeout 30m $PK"
RACEV=""; [ "$MODE" = base ] || RACEV="go test -race -count=1 ./host/verbs/ && "
leg race "$RACEV go test -race -count=1 ./host/surfacepage/ && go test -race -count=1 ./host/daemon/ -run '^(TestSurface|TestEpisodeReport|TestApprov|TestDaemonCredentialWriterSites)' && go test -race -count=1 ./cmd/ailang-worldd/ -run '^(TestEpisodeReport|TestVerb|TestSession|TestApprove|TestPreset|TestEpisode)'"
NEW_V="TestSplit TestQuoteRoundTrip TestLookup TestEpisodeReportFlags TestParseInterleaved TestApproveTarget"
NEW_D="TestSurfaceVerbRouteRegistered TestSurfaceVerbOutcomes TestSurfaceVerbBoundsAndErrors TestSurfaceVerbIsReadOnly"
NEW_SP="TestSurfaceAppVerbs"
NEW_C="TestVerbTableMatchesTheCLISwitch TestVerbRouteReportIsTheCLIs TestQuickstartSection19MatchesTheBar"
NEW_X="TestVerbBarEpisodeReportInChrome"
if [ "$MODE" = base ]; then NEW_V=""; NEW_D=""; NEW_SP=""; NEW_C=""; NEW_X=""; fi
alt() { local x="$*"; x="${x// /|}"; echo "${x:+|$x}"; }
INV_DAEMON="TestDaemonDependencyAllowlist TestFrozenV1RouteTableMatchesMux TestSurfaceAppRoutes TestApprovalRoutesRegistered TestDecisionRoutesRegistered TestReadCtxCancelledAfterHandler TestMCPMountMethodSource TestDaemonCredentialWriterSites TestQuickstartDecisionPacketSection TestDaemonConstructorCallerInventory TestApprovePageRoutes TestPasskeyPageRoutes TestApprovalRoutesRefuseForeignHost TestGETRoutesRejectOtherMethods TestDaemonCallerInventoryAliasControl TestEpisodeReportRESTRoute"
INV_BOUNDARY="TestWorldBoundaryDependencyAllowlist TestBareNetHTTPExemptionIsPerGroup TestEveryCommandDirectoryIsAProtectedGroup TestBoundaryASTWriteGuard TestWorkbenchPackageRemainsTransportFree"
INV_STORE="TestProductionContextRoots TestProductionGoSurface TestNoNewDeadlineFreeStoreReads TestReadGettersHonorContext TestCredentialsMethodsRequireADeadline"
INV_VG="TestNoParallelWriteForkPackagesOutsideVerifygate TestSubprocessSinkGate TestVerifygateTestWritesAreForkLocked"
INV_SP="TestSurfaceAppPure TestSurfaceScriptPure TestSurfaceScriptAssetSeams TestSurfaceAppMatchesWorkingCopy TestSurfaceAppShellAssembly TestNeedsYouListsDecisions"
INV_CMD="TestEveryVerbHasHelp TestEpisodeReportUsage TestEpisodeReportTablePrintsEveryClass TestEpisodeReportTablePrintsTheReadsLine TestEpisodeReportJSONAndTableLoopback TestEpisodeReportFollowReprintsOnChangeLoopback TestEpisodeReportPrintsExecCache TestEpisodeGlobalAddrRefusedForStartOnly TestGlobalBareAddrIsNormalisedForEveryClientVerb TestQuickstartAttendedStepsBoxMatchesSection9 TestQuickstartSection9FlagsMatchTheCLI TestQuickstartSection10FlagsMatchTheCLI TestQuickstartSection11FlagsMatchTheCLI TestQuickstartSection18FlagsMatchTheCLI TestQuickstartSection18CardMatchesThePage TestQuickstartServiceSectionMatchesTheCLIHelp"
leg inv "set -o pipefail
go test -v -count=1 ./host/daemon/ -run '^($(echo $INV_DAEMON | tr ' ' '|')$(alt $NEW_D))\$' &&
go test -v -count=1 ./host/boundary/ -run '^($(echo $INV_BOUNDARY | tr ' ' '|'))\$' &&
go test -v -count=1 ./host/store/ -run '^($(echo $INV_STORE | tr ' ' '|'))\$' &&
go test -v -count=1 ./host/verifygate/ -run '^($(echo $INV_VG | tr ' ' '|'))\$' &&
go test -v -count=1 ./host/surfacepage/ -run '^($(echo $INV_SP | tr ' ' '|')$(alt $NEW_SP))\$' &&
go test -v -count=1 ./cmd/ailang-worldd/ -run '^($(echo $INV_CMD | tr ' ' '|')$(alt $NEW_C))\$' &&
{ [ -z '$NEW_V' ] || go test -v -count=1 ./host/verbs/ -run '^($(echo $NEW_V | tr ' ' '|'))\$'; }"
XB="TestApprovePageInChrome TestApprovePageRefusesASwappedRequest TestApprovePageSSHSIGInChrome TestApprovePageShowsAHostileValueAsText TestPasskeyPageAbortsAPendingCeremony TestPasskeyPageInChrome TestApproveDecideWebAuthnLinkInChrome TestDecidePageApproveWithPasskey TestDecidePageFailsClosed TestDecidePageSSHSIGPaste TestNeedsYouListsPendingRequests TestNeedsYouNotifiesOnlyWhenAskedOnce"
XD="TestSurfaceLiveChrome TestSurfaceFollowChrome"
CH="'/Applications/Google Chrome.app/Contents/MacOS/Google Chrome'"
leg chrome "set -o pipefail; WORLD_CHROME=$CH WORLD_REQUIRE_PASSKEY_BROWSER=1 WORLD_REQUIRE_SSH_KEYGEN=1 go test ./cmd/ailang-worldd/ -run '^($(echo $XB | tr ' ' '|')$(alt $NEW_X))\$' -count=1 -v && WORLD_CHROME=$CH go test ./host/daemon/ -run '^($(echo $XD | tr ' ' '|'))\$' -count=1 -v"
for t in $XB $NEW_X $XD; do grep -qE -- "^--- PASS: $t " "$L/chrome.log" || { echo "CHROME NOT PASSED: $t" | tee -a "$L/legs.txt"; FAILS="$FAILS chrome_missing"; }; done
MISSING=""
for t in $INV_DAEMON $NEW_D $INV_BOUNDARY $INV_STORE $INV_VG $INV_SP $NEW_SP $INV_CMD $NEW_C $NEW_V; do
  grep -qE -- "^--- PASS: $t " "$L/inv.log" || MISSING="$MISSING $t"
done
if grep -E '^\s*--- (SKIP|FAIL)' "$L/inv.log"; then FAILS="$FAILS inv_skip_or_fail"; fi
[ -z "$MISSING" ] || { echo "NOT PASSED:$MISSING" | tee -a "$L/legs.txt"; FAILS="$FAILS inv_missing"; }
echo "named PASS: $(grep -cE '^--- PASS' "$L/inv.log")" | tee -a "$L/legs.txt"
leg verify_go 'VERIFY_GO_LEGS=plain AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_go.sh'
echo "GATES_DONE $(( $(date +%s) - s0 ))s FAILS:${FAILS:- none}" | tee -a "$L/legs.txt"
echo "GATES_DONE $(( $(date +%s) - s0 ))s FAILS:${FAILS:- none}" > "$L/ALL_DONE.$SHA"
````

**Expected:** every leg rc 0 except a known-load red, re-run alone once (§1.1's rule); `inv` prints **named PASS**
for every inventory and every new test; `chrome` prints `--- PASS` for all fifteen. The planner's run of this script
on the final tree: §1.2 C19.

---

## 7. Executor directive notes

- **No git writes of any kind** (no add, commit, stash, checkout, reset, branch). **`patch -p1` (§2) is the only
  writer** of the tree; no hand edit, no `gofmt -w`, no fuzz, no `WORLD_UPDATE_VERB_GOLDEN=1` (the golden ships in
  M3b's patch). A failed `--dry-run` or a red acceptance → **stop and report** with the log; do not repair.
- **Never `nohup`, never background a test leg; pass a 600000 ms timeout to every long Bash call** (the iteration-389
  executor lost a backgrounded acceptance shell to the 120 s default). M1c's and M2's acceptance run the whole
  `cmd/ailang-worldd` package (≈ 2–4 min at load 10–20); M3a runs `host/daemon`'s route, report, approval and
  inventory tests by regex (≈ 20 s) — the whole `host/daemon` package (≈ 9–10 min under load, over one call's 600 s)
  is §6's `pkgs` leg. If a call still nears 600 s, run `accept.sh` again after a timeout rather than backgrounding it.
- **Runs**: Run 1 = M1a, M1b, M1c; Run 2 = M2, M3a, M3b; Run 3 = M4a, M4b. After each sub-milestone's acceptance,
  snapshot to `$HOME/.ailang/state/world-iter395/snap/<M>/`, **copying only the files `git status --porcelain`
  reports** (cumulative against the base: modified and untracked files, plus a `DELETED` list — none are expected),
  **never an `rsync` of whole directories**:
  ```bash
  bash <<'EOF2'
  set -euo pipefail
  M=M1a   # the sub-milestone just accepted
  W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter395-row243-v1
  D="$HOME/.ailang/state/world-iter395/snap/$M"; rm -rf "$D"; mkdir -p "$D"
  cd "$W"
  git status --porcelain --untracked-files=all | grep -v -e 'w-world-address-bar-verbs-v1-sprint-plan.md$' -e 'sprint_w-world-address-bar-verbs-v1.json$' > "$D/STATUS"
  awk '$1 != "D" {print $NF}' "$D/STATUS" | while read -r p; do mkdir -p "$D/$(dirname "$p")"; cp -p "$p" "$D/$p"; done
  awk '$1 == "D" {print $NF}' "$D/STATUS" > "$D/DELETED"
  (cd "$D" && find cmd host docs -type f -print0 2>/dev/null | sort -z | xargs -0 shasum -a 256) > "$D/MANIFEST.sha256"
  cat "$D/STATUS"; wc -l < "$D/MANIFEST.sha256"
  EOF2
  ```
  The two plan files are excluded. The controller commits each snapshot as one commit and checks the numstat against
  §3's expected value (`git diff --numstat db2b7423 HEAD -- cmd host docs | awk '{s+=$1+$2} END {print s}'`): M1a
  203, M1b 510, M1c 697, M2 905, M3a 1,248, M3b 1,491, M4a 1,762, M4b 1,909. A different number → stop.
- `export AILANG_BIN=$HOME/.pinned-ailang/ailang` and `WORLD_EXEC_NODE=/opt/homebrew/bin/node` in every shell that
  runs Go tests (`accept.sh` sets both). The controller's shell is zsh: run each block under `bash <<'EOF2' … EOF2`.
  The live-Chrome lines carry `WORLD_CHROME="/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
  WORLD_REQUIRE_PASSKEY_BROWSER=1` (the `vc` lines of `accept.sh`; the A2 driver makes its own temp
  `--user-data-dir`; never a real profile). With `WORLD_CHROME` set, run only named tests (never the whole `cmd`
  package: `TestMarketingChromeSixShots` needs other env).
- **Stop conditions**: a dry-run failure; any `--- FAIL`, `--- SKIP` or `FAIL` line in an acceptance (`ACCEPT_RED`);
  gofmt output; a numstat mismatch; any file outside §5 changed (`git status --short` after each patch must equal §5's
  list up to that stage).
- **Known load-sensitive tests** (the directive's list) are not this sprint's: a red there is re-run once alone.
- D-WORLD-109: nothing touches `~/.ailang/world`, `~/.ssh`, a real `/dev/tty`, a real browser profile, or a real
  daemon (§0.1).

---

## 8. Risks

- **R1 — RK1, the shadow set (design §11).** 27 bare words (`log`, `episode`, `session`, …) that today navigate to a
  brane so named now answer the CLI's refusal; such a brane is reached as `world NAME`. Deliberate and enumerated; VB3
  and VB4 pin the set to the CLI's switch, so a new CLI verb without a class fails the build.
- **R2 — The bar depends on the route for navigation too.** While the route fails, nothing navigates from the bar
  (design §2.3's honest state); the timeline, the crumbs and the `?`/`N` keys still work. XV1 proves the healthy path.
- **R3 — RS-SAY**: a refusal is the header's transient `say` line (2.6 s); the input is kept so Enter re-shows it.
- **R4 — Live Chrome only on the controller's machine** (CI has no Chrome step): XV1 and the `JW-*` mutants are gated
  locally as the A3/A4 pages were; a CI run skips them, counted.
- **R5 — Size margin 91.** Any executor-side growth is a stop (§7's numstat check), never a hand edit.
- **R6 — Load.** The planner's runs saw load 12–20; every base red was a load timeout in an untouched test (§1.1).

---

## 9. Acceptance criteria (the judge scores against these)

1. **One grammar** (`host/verbs`, pure, standard-library-only): `Split`'s five named refusals at their exact bounds and
   the nine shell characters; `Split(Quote(a)) == a` (VB1, VB2).
2. **The table** classifies exactly the CLI switch's 28 arms (VB4, AST-derived, no `run()`): `world`/`world NAME` →
   address; `world get <ref>` → refused with its CLI line; `episode report` → runs; `approve` → navigate;
   `why`/`episode start`/`episode finish` → terminal; `episode launch` → terminal (host); the five never words and
   nineteen other words → refused with the CLI line; bare/unknown `episode` and `world a b c` → refused (step 4b);
   global flags refused; anything else → address (VB3).
3. **One flag parse**: the CLI and the bar run `verbs.EpisodeReportFlags` + `verbs.ParseInterleaved`; the bar's
   profile allows `--json` only (VB5); the ten CLI call sites use the moved `ParseInterleaved`; every existing CLI
   test passes unedited.
4. **`GET /surface/verb.json`** is one route behind `approvalRoute` (421 for a foreign Host), GET only, no CORS
   header, read-only (VS1), `/v1` still 13 routes (`TestFrozenV1RouteTableMatchesMux`); 400 without `q`; an operator
   typo or a 1,025-byte line is a 200 refusal (VR1, VR4).
5. **The same bytes**: the route's `text` and `json` equal the real CLI's stdout for two different episodes, and its
   not-found answer is the CLI's stderr line and exit code (VR3); every outcome class answers as design §2.4 (VR2).
6. **The `/v1` route is unchanged** by the extraction (503 for a builder timeout, 500 for a store failure, 404, 200
   bodies; `TestEpisodeReportRESTRoute` unedited).
7. **No silent fallback** (VP2): for `why`, `log`, `row-189@6`, `world shared@0ffbb9c`, each of 404/400/421/500, a
   network error, a non-JSON 200, a schema-less 200, an unknown outcome — an explicit message, and no navigation,
   panel or input change; only the golden `address` answers reach the address path, with today's `addressTarget`
   answers.
8. **Bounded wait** (VP2): `VERB_ROUTE_BUDGET_MS = 15000`, enforced across fetch and body with one `AbortController`;
   the timeout says so and changes nothing; the timer is cleared on every path; tested with a controlled clock
   (V1-D2).
9. **Superseded ask** (V1-D1): the older ask is aborted and its answer never acted on; the newer answers.
10. **The page's other outcomes** (VP1, VP3–VP6): result/terminal panels set by `textContent` (hostile text stays
    text); navigate only to the two card prefixes; `world why` → address and `why` → the route; the shadow set at page
    level from the route's real answers.
11. **Live Chrome (XV1)**: typing `episode report w1` into `#addr` shows the CLI's exact bytes; `why` shows its CLI
    line; `log` is refused with the bar and the place kept.
12. **QUICKSTART §19** documents the bar (`TestQuickstartSection19MatchesTheBar`); the §17/§18 tests pass unedited.
13. **Inventories** all PASS (`TestFrozenV1RouteTableMatchesMux` 13, `TestSurfaceAppRoutes`,
    `TestApprovalRoutesRegistered`, `TestDecisionRoutesRegistered`, `TestDaemonDependencyAllowlist`,
    `TestWorldBoundaryDependencyAllowlist`, `TestProductionContextRoots`, `TestProductionGoSurface`,
    `TestDaemonCredentialWriterSites`, `TestEveryVerbHasHelp`, `TestSurfaceAppPure` (edited +2 −1, P2), the
    quickstart tests); the only edited existing test is `app_test.go`'s export list.
14. **Size** ≤ 2,000 (`-- cmd host docs`): 1,909 measured; each sub-milestone's numstat as §3; M1–M3 1,491 (no cut).
15. **Every mutant in §4 behaves as expected**, its id in the ledger; no suite kill rests on an inventory test.
16. **The round trip is identical** (§1.2 C17) and the gates (§6) pass, `verify_go` after commit included;
    `./scripts/verify_ail.sh` gives the base's result.

**What an operator can do after this sprint:** open the World app (`/surface`), press `/`, type `episode report
<episode>` (or paste `ailang-worldd episode report <episode> --json`) and read, in a panel, exactly what the terminal
prints; type `approve mr-…` or `approve sha256:…` and land on row 201's card for that request; type any other operator
verb (`episode finish w1 --preset se-tools`, `why …`) and see the exact CLI line with a copy button and that the app
cannot run it yet; a CLI word typed bare (`log`) is refused with its CLI line, and `world log` still opens the brane.

**What the attended exit run (design §10) still waits for:** V2 (`why` in the bar), V3 (the verb-approval law: a
key-confirmed `episode start`/`finish`, its ledger ask on A4 L1's reason), V4 (the host runner; herdr hand-off), V5
(the attended script), and Mark's own run with his enrolled key (D-WORLD-109). V1's optional attended check is
QUICKSTART §19 verbatim (design §10, last line).

---

## Appendix C — after this sprint (V1's close and the V2 handoff)

- **Moved to V2: nothing.** M1–M3 measured 1,491 ≤ 1,500 (V1-D3); the fallback cut (XV1, 78 lines, and the `approve`
  navigate arm, ≈ 12 lines) is not taken.
- **Named survivors carried**: `E-ASb` (the `/v1` extraction's start/finish-read error *after* an expired deadline would answer 503 instead of 500 —
unreachable in a test: under an expired deadline the builder fails first; HEAD's behaviour is kept by
`reportBuildError`, V1-D9); `JW-GO-EMPTY` (`go('')` would ask the route instead of going home without a fetch — the
route answers `address` for an empty `q` and the address path goes home, so the operator sees the same view; only a
request count would see it; ≈ 15 lines in XV1 if wanted).
- **Residuals**: RS-SAY (a refusal fades from the header after 2.6 s; a persistent refusal line is ≈ 10 lines in V2
  if Mark wants it); R2 (the bar navigates only through the route).
- **V2 handoff** (`why <target>` in the bar, design §0.3): `host/verbs` already classifies `why` as `read`/terminal;
  V2 moves `why.go`'s walk behind a reader interface and adds a `why` arm in `handleSurfaceVerb` beside
  `verbEpisodeReport`, reusing the same answer wire (`result` with `text`/`json`) and the same page path (no page
  change needed beyond a golden refresh). V2's child design must decide `why --result <file>` (T3: a path flag) —
  refused in the bar's profile as `--addr` is here. V3 migrates `episode start`/`finish`'s flag sets into
  `host/verbs` the way M1c migrated `episode report`'s.

---

## Appendix P — the patches (apply with §2; generated from the compiled scratch stages; do not edit)

### Patch M1a — `Split`, `Quote` (VB1, VB2)

````diff M1a
diff --git a/host/verbs/split.go b/host/verbs/split.go
new file mode 100644
index 0000000..d7f5d00
--- /dev/null
+++ b/host/verbs/split.go
@@ -0,0 +1,110 @@
+// Package verbs is the operator verbs' one grammar (row 243 V1, design §2.1): the address bar's
+// tokenizer (Split, Quote), the verb table (Lookup), the CLI's interleaved flag parse
+// (ParseInterleaved) and episode report's flag set with the bar's admission. It is pure: no I/O, no
+// clock, no randomness. The CLI and the daemon's GET /surface/verb.json both call it, so the bar
+// and a terminal read one grammar.
+package verbs
+
+import (
+	"errors"
+	"regexp"
+	"strings"
+	"unicode/utf8"
+)
+
+// MaxLine and MaxWords bound one line of the bar (design §2.1(a)).
+const (
+	MaxLine  = 1024
+	MaxWords = 32
+)
+
+// The Split refusals, each named; the name leads the message.
+var (
+	ErrTooLong  = errors.New("TooLong: a line is at most 1024 bytes")
+	ErrTooMany  = errors.New("TooMany: a line is at most 32 words")
+	ErrNotText  = errors.New("NotText: a line is UTF-8 text with no control character")
+	ErrUnclosed = errors.New("Unclosed: a quote is not closed")
+	ErrShell    = errors.New("Shell: the bar runs no shell; quote $ ` ; | & < > ( ) to use them as text")
+)
+
+// shellChars are refused unquoted: the bar runs no shell, so what it shows is what a shell would run.
+const shellChars = "$`;|&<>()"
+
+// Split cuts one line into words: spaces separate them, '…' is literal, "…" is literal but for \"
+// and \\, and adjacent segments join into one word.
+func Split(line string) ([]string, error) {
+	if len(line) > MaxLine {
+		return nil, ErrTooLong
+	}
+	if !utf8.ValidString(line) {
+		return nil, ErrNotText
+	}
+	for i := 0; i < len(line); i++ {
+		if line[i] < 0x20 || line[i] == 0x7f {
+			return nil, ErrNotText
+		}
+	}
+	words := []string{}
+	var w strings.Builder
+	in := false
+	for i := 0; i < len(line); i++ {
+		c := line[i]
+		switch {
+		case c == ' ':
+			if in {
+				words = append(words, w.String())
+				w.Reset()
+				in = false
+			}
+			continue
+		case c == '\'':
+			j := strings.IndexByte(line[i+1:], '\'')
+			if j < 0 {
+				return nil, ErrUnclosed
+			}
+			w.WriteString(line[i+1 : i+1+j])
+			i += j + 1
+		case c == '"':
+			j := i + 1
+			for ; j < len(line) && line[j] != '"'; j++ {
+				if line[j] == '\\' && j+1 < len(line) && (line[j+1] == '"' || line[j+1] == '\\') {
+					j++
+				}
+				w.WriteByte(line[j])
+			}
+			if j == len(line) {
+				return nil, ErrUnclosed
+			}
+			i = j
+		case strings.IndexByte(shellChars, c) >= 0:
+			return nil, ErrShell
+		default:
+			w.WriteByte(c)
+		}
+		in = true
+	}
+	if in {
+		words = append(words, w.String())
+	}
+	if len(words) > MaxWords {
+		return nil, ErrTooMany
+	}
+	return words, nil
+}
+
+// safeWord is a word Quote prints bare: nothing a shell or Split reads specially.
+var safeWord = regexp.MustCompile(`^[A-Za-z0-9@%+=:,./_-]+$`)
+
+// Quote prints argv as one line a shell and Split both read back as argv: a safe word bare, any
+// other in single quotes, a ' as '"'"'.
+func Quote(argv []string) string {
+	out := make([]string, len(argv))
+	for i, a := range argv {
+		if safeWord.MatchString(a) {
+			out[i] = a
+			continue
+		}
+		out[i] = "'" + strings.ReplaceAll(a, "'", `'"'"'`) + "'"
+	}
+	return strings.Join(out, " ")
+}
diff --git a/host/verbs/split_test.go b/host/verbs/split_test.go
new file mode 100644
index 0000000..3fc864a
--- /dev/null
+++ b/host/verbs/split_test.go
@@ -0,0 +1,93 @@
+package verbs
+
+// Row 243 V1 (design §6 VB1, VB2).
+
+import (
+	"math/rand"
+	"reflect"
+	"strings"
+	"testing"
+)
+
+// VB1: Split's words and each named refusal, at and past each bound.
+func TestSplit(t *testing.T) {
+	for _, c := range []struct {
+		line string
+		want []string
+	}{
+		{"", []string{}},
+		{"  episode   report  w1 ", []string{"episode", "report", "w1"}},
+		{`why 'a b' "c d"`, []string{"why", "a b", "c d"}},
+		{`'' x`, []string{"", "x"}},
+		{`a'b'"c"d`, []string{"abcd"}},
+		{`"q\"uo\\te\n"`, []string{`q"uo\te\n`}},
+		{`'it'"'"'s'`, []string{"it's"}},
+		{`a\b`, []string{`a\b`}},
+		{`'$x' "a;b|c&d<e>f(g)h` + "`" + `"`, []string{"$x", "a;b|c&d<e>f(g)h`"}},
+		{"é ü", []string{"é", "ü"}},
+		{strings.Repeat("x", MaxLine), []string{strings.Repeat("x", MaxLine)}},
+		{strings.TrimSpace(strings.Repeat("w ", MaxWords)), strings.Fields(strings.Repeat("w ", MaxWords))},
+	} {
+		got, err := Split(c.line)
+		if err != nil || !reflect.DeepEqual(got, c.want) {
+			t.Errorf("Split(%q) = %q, %v; want %q", c.line, got, err, c.want)
+		}
+	}
+	for _, c := range []struct {
+		line string
+		want error
+	}{
+		{strings.Repeat("x", MaxLine+1), ErrTooLong},
+		{strings.TrimSpace(strings.Repeat("w ", MaxWords+1)), ErrTooMany},
+		{"a\nb", ErrNotText}, {"a\tb", ErrNotText}, {"a\x7fb", ErrNotText}, {"a\x00", ErrNotText}, {"a\x1f", ErrNotText},
+		{"a \xff", ErrNotText}, {`'a`, ErrUnclosed}, {`"a`, ErrUnclosed}, {`"a\"`, ErrUnclosed}, {`a 'b' "c`, ErrUnclosed},
+	} {
+		if got, err := Split(c.line); err != c.want || got != nil {
+			t.Errorf("Split(%.40q) = %q, %v; want %v", c.line, got, err, c.want)
+		}
+	}
+	for _, ch := range "$`;|&<>()" { // the nine
+		if _, err := Split("why a" + string(ch) + "b"); err != ErrShell {
+			t.Errorf("unquoted %q: %v, want Shell", ch, err)
+		}
+		for _, q := range []string{"'a" + string(ch) + "b'", `"a` + string(ch) + `b"`} {
+			if got, err := Split(q); err != nil || len(got) != 1 || got[0] != "a"+string(ch)+"b" {
+				t.Errorf("quoted %s = %q, %v", q, got, err)
+			}
+		}
+	}
+	if !strings.HasPrefix(ErrTooLong.Error(), "TooLong") || !strings.HasPrefix(ErrShell.Error(), "Shell") {
+		t.Error("a refusal's message does not lead with its name")
+	}
+}
+
+// VB2: Quote is Split's inverse over generated argv, and a safe word stays bare.
+func TestQuoteRoundTrip(t *testing.T) {
+	alphabet := []rune("ab-_./:@=+%,'\"\\ $`;|&<>()*?~#!{}[]éx")
+	rng := rand.New(rand.NewSource(243))
+	check := func(argv []string) {
+		t.Helper()
+		line := Quote(argv)
+		got, err := Split(line)
+		if err != nil || !reflect.DeepEqual(got, argv) {
+			t.Fatalf("Split(Quote(%q)) = Split(%s) = %q, %v", argv, line, got, err)
+		}
+	}
+	for i := 0; i < 2000; i++ {
+		argv := make([]string, 1+rng.Intn(6))
+		for j := range argv {
+			w := make([]rune, rng.Intn(8))
+			for k := range w {
+				w[k] = alphabet[rng.Intn(len(alphabet))]
+			}
+			argv[j] = string(w)
+		}
+		check(argv)
+	}
+	if got := Quote([]string{"ailang-worldd", "episode", "report", "w1", "--json"}); got != "ailang-worldd episode report w1 --json" {
+		t.Errorf("safe words quoted: %s", got)
+	}
+	if got := Quote([]string{"a b", "it's", "$x", ""}); got != `'a b' 'it'"'"'s' '$x' ''` {
+		t.Errorf("Quote = %s", got)
+	}
+}
````

### Patch M1b — the verb table (VB3, VB4)

````diff M1b
diff --git a/cmd/ailang-worldd/verbs_parity_test.go b/cmd/ailang-worldd/verbs_parity_test.go
new file mode 100644
index 0000000..de4d084
--- /dev/null
+++ b/cmd/ailang-worldd/verbs_parity_test.go
@@ -0,0 +1,80 @@
+package main
+
+// Row 243 V1 (design §6 VB4): host/verbs' three word sets are exactly run's switch arms but the
+// navigation token world, each in one set, each with a help path — read with go/ast, never run().
+
+import (
+	"go/ast"
+	"go/parser"
+	"go/token"
+	"sort"
+	"strconv"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/verbs"
+)
+
+// switchArms is the case string literals of run's `switch verb := rest[0]`.
+func switchArms(t *testing.T) (arms []string) {
+	f, err := parser.ParseFile(token.NewFileSet(), "main.go", nil, 0)
+	if err != nil {
+		t.Fatal(err)
+	}
+	ast.Inspect(f, func(n ast.Node) bool {
+		if fn, ok := n.(*ast.FuncDecl); ok {
+			return fn.Name.Name == "run"
+		}
+		sw, ok := n.(*ast.SwitchStmt)
+		if !ok || sw.Init == nil {
+			return true
+		}
+		if as, ok := sw.Init.(*ast.AssignStmt); !ok || as.Lhs[0].(*ast.Ident).Name != "verb" {
+			return true
+		}
+		for _, s := range sw.Body.List {
+			for _, e := range s.(*ast.CaseClause).List {
+				if lit, ok := e.(*ast.BasicLit); ok && lit.Kind == token.STRING {
+					w, _ := strconv.Unquote(lit.Value)
+					arms = append(arms, w)
+				}
+			}
+		}
+		return false
+	})
+	sort.Strings(arms)
+	return arms
+}
+
+// classified: arms minus world equal the union of sets, each word in exactly one set.
+func classified(arms []string, sets ...[]string) bool {
+	seen, n := map[string]int{}, 0
+	for _, set := range sets {
+		for _, w := range set {
+			seen[w]++
+		}
+	}
+	for _, a := range arms {
+		if a != "world" && seen[a] != 1 {
+			return false
+		} else if a != "world" {
+			n++
+		}
+	}
+	return n == len(seen)
+}
+
+func TestVerbTableMatchesTheCLISwitch(t *testing.T) {
+	arms, op, nv, ot := switchArms(t), verbs.OperatorWords, verbs.NeverWords, verbs.OtherWords
+	if len(arms) != 28 || !classified(arms, op, nv, ot) {
+		t.Fatalf("run's switch arms %q (want 28) are not exactly operator ∪ never ∪ other, each once, world excluded", arms)
+	}
+	// Controls: an extra arm, a word the table lacks, a word in two sets each fail the check.
+	if classified(append(arms, "nonsuch"), op, nv, ot) || classified(arms, op, nv, ot[1:]) || classified(arms, op, nv, append(ot, "why")) {
+		t.Error("a control passed the parity check")
+	}
+	for _, w := range append(append(append([]string{}, op...), nv...), ot...) {
+		if _, ok := genericVerbHelp([]string{w, "--help"}); !ownHelpVerbs[w] && !ok && w != "help" {
+			t.Errorf("%s has no help path", w)
+		}
+	}
+}
diff --git a/host/verbs/lookup.go b/host/verbs/lookup.go
new file mode 100644
index 0000000..399e275
--- /dev/null
+++ b/host/verbs/lookup.go
@@ -0,0 +1,133 @@
+package verbs
+
+import (
+	"regexp"
+	"strings"
+)
+
+// Class is what the bar may do with a verb (design §2.1(b)).
+type Class string
+
+const (
+	Address  Class = "address"  // navigation: the page's own parseAddress path
+	Read     Class = "read"     // reads World: episode report runs; why answers terminal (V2)
+	Navigate Class = "navigate" // approve opens row 201's card
+	Mutate   Class = "mutate"   // changes World: terminal (V3/V4)
+	Host     Class = "host"     // needs a host process: terminal (V4)
+	Never    Class = "never"    // typed-only, never from a page
+	Other    Class = "other"    // a CLI command that is not an operator verb
+	Refused  Class = "refused"  // a global flag, or a seized first word with no row
+)
+
+// Verb is one answer of Lookup: the table row's words (empty for an address), its class, the
+// reason the bar gives, and the CLI line (empty when the input is not a CLI line).
+type Verb struct {
+	Words  string
+	Class  Class
+	Reason string
+	CLI    string
+}
+
+// Binary is the CLI's name; a pasted line may start with it.
+const Binary = "ailang-worldd"
+
+// The table's word sets. They are exactly the CLI switch's arms but world (VB4 asserts it).
+var (
+	OperatorWords = []string{"why", "episode", "approve"}
+	NeverWords    = []string{"approver", "session", "service", "serve", "setup"}
+	OtherWords    = []string{"commit", "outcome", "call", "demo", "doctor", "head", "health", "help", "hosts", "log", "object",
+		"preset", "profile", "provenance", "registry", "surface-index", "tools", "version", "whoami"}
+)
+
+// neverWhy is each never verb's reason (D-WORLD-89 R3; host administration stays at a terminal).
+var neverWhy = map[string]string{
+	"approver": "enrolling and revoking an approver key is typed at a terminal",
+	"session":  "minting and revoking a session is typed at a terminal",
+	"service":  "host administration stays at a terminal",
+	"serve":    "host administration stays at a terminal",
+	"setup":    "host administration stays at a terminal",
+}
+
+// EpisodeNeedsSubcommand is step 4b's refusal for a bare or unknown episode form.
+const EpisodeNeedsSubcommand = Binary + " episode needs a subcommand: start | launch | finish | report"
+
+// AddressHelp is the page's own address help (surface.js ADDRESS_HELP), step 4b's text for world.
+const AddressHelp = `unknown address: try "world", "world NAME", "world NAME@N", "world NAME@HASH" or "back"`
+
+func member(set []string, w string) bool {
+	for _, s := range set {
+		if s == w {
+			return true
+		}
+	}
+	return false
+}
+
+// Lookup classifies argv (design §2.1(b), steps 1–5) and returns the verb and the words after it.
+func Lookup(argv []string) (Verb, []string) {
+	if len(argv) > 0 && argv[0] == Binary { // 1
+		argv = argv[1:]
+	}
+	if len(argv) == 0 {
+		return Verb{Class: Address}, nil
+	}
+	cli := Quote(append([]string{Binary}, argv...))
+	first := argv[0]
+	if strings.HasPrefix(first, "-") { // 2
+		return Verb{Class: Refused, Reason: "global flags are not accepted: the app talks to its own daemon", CLI: cli}, nil
+	}
+	if first == "world" && len(argv) <= 2 { // 3 (back is not a CLI word: step 5 keeps it an address)
+		return Verb{Class: Address}, nil
+	}
+	if first == "world" && argv[1] == "get" { // 3b
+		return Verb{Words: "world get", Class: Other, CLI: cli,
+			Reason: "world get is a CLI command the app does not run; at a terminal: " + cli}, argv[2:]
+	}
+	if first == "episode" && len(argv) > 1 { // 4: the two-word rows
+		switch argv[1] {
+		case "report":
+			return Verb{Words: "episode report", Class: Read, CLI: cli}, argv[2:]
+		case "start", "finish":
+			return Verb{Words: "episode " + argv[1], Class: Mutate, CLI: cli,
+				Reason: "episode " + argv[1] + " changes World. The app cannot run it yet."}, argv[2:]
+		case "launch":
+			return Verb{Words: "episode launch", Class: Host, CLI: cli,
+				Reason: "episode launch needs a host process. The app cannot run it yet."}, argv[2:]
+		}
+	}
+	switch { // 4: the one-word rows
+	case first == "why":
+		return Verb{Words: "why", Class: Read, CLI: cli, Reason: "why reads World, but its walk runs in the CLI. The app cannot run it yet."}, argv[1:]
+	case first == "approve":
+		return Verb{Words: "approve", Class: Navigate, CLI: cli}, argv[1:]
+	case member(NeverWords, first):
+		return Verb{Words: first, Class: Never, CLI: cli,
+			Reason: first + " is never run from the app: " + neverWhy[first] + "; at a terminal: " + cli}, argv[1:]
+	case member(OtherWords, first):
+		return Verb{Words: first, Class: Other, CLI: cli,
+			Reason: first + " is a CLI command the app does not run; at a terminal: " + cli}, argv[1:]
+	}
+	switch first { // 4b: a seized first word with no row
+	case "episode":
+		return Verb{Words: "episode", Class: Refused, Reason: EpisodeNeedsSubcommand, CLI: cli}, argv[1:]
+	case "world":
+		return Verb{Words: "world", Class: Refused, Reason: AddressHelp}, argv[1:]
+	}
+	return Verb{Class: Address}, nil // 5
+}
+
+var (
+	mintRequestID = regexp.MustCompile(`^mr-[0-9a-f]{16}$`)
+	decisionRef   = regexp.MustCompile(`^sha256:[0-9a-f]{64}$`)
+)
+
+// ApproveTarget is approve's navigate arm: one request id or decision ref, no flag, opens its card.
+func ApproveTarget(args []string) (href, reason string) {
+	if len(args) == 1 && mintRequestID.MatchString(args[0]) {
+		return "/surface/approve/" + args[0], ""
+	}
+	if len(args) == 1 && decisionRef.MatchString(args[0]) {
+		return "/surface/decide/" + args[0], ""
+	}
+	return "", "the app opens approve <request-id|request-ref> as its card; any flag or other form runs at a terminal"
+}
diff --git a/host/verbs/lookup_test.go b/host/verbs/lookup_test.go
new file mode 100644
index 0000000..04a48a0
--- /dev/null
+++ b/host/verbs/lookup_test.go
@@ -0,0 +1,94 @@
+package verbs
+
+// Row 243 V1 (design §6 VB3): the verb table and approve's navigate arm.
+
+import (
+	"reflect"
+	"strings"
+	"testing"
+)
+
+var lookupCases = []struct {
+	line, words string
+	class       Class
+	reason      string
+	cli         bool
+	rest        []string
+}{
+	{"row-189@6", "", Address, "", false, nil},
+	{"world", "", Address, "", false, nil},
+	{"back", "", Address, "", false, nil},
+	{"back x", "", Address, "", false, nil},
+	{"world why", "", Address, "", false, nil},
+	{"world session", "", Address, "", false, nil},
+	{"world get", "", Address, "", false, nil},
+	{"World a b c", "", Address, "", false, nil},
+	{"ailang-worldd", "", Address, "", false, nil},
+	{"why", "why", Read, "why reads World", true, []string{}},
+	{"why e1 --json", "why", Read, "why reads World", true, []string{"e1", "--json"}},
+	{"episode report w1", "episode report", Read, "", true, []string{"w1"}},
+	{"ailang-worldd episode report w1 --json", "episode report", Read, "", true, []string{"w1", "--json"}},
+	{"approve mr-0123456789abcdef", "approve", Navigate, "", true, []string{"mr-0123456789abcdef"}},
+	{"episode start w1", "episode start", Mutate, "changes World", true, []string{"w1"}},
+	{"episode finish w1 --preset p", "episode finish", Mutate, "changes World", true, []string{"w1", "--preset", "p"}},
+	{"episode launch w1", "episode launch", Host, "needs a host process", true, []string{"w1"}},
+	{"approver enroll --name x", "approver", Never, "approver is never run from the app: enrolling", true, []string{"enroll", "--name", "x"}},
+	{"session frob", "session", Never, "minting and revoking a session", true, []string{"frob"}},
+	{"serve", "serve", Never, "host administration", true, []string{}},
+	{"log", "log", Other, "log is a CLI command the app does not run; at a terminal: ailang-worldd log", true, []string{}},
+	{"log tail", "log", Other, "at a terminal: ailang-worldd log tail", true, []string{"tail"}},
+	{"ailang-worldd world get world-1", "world get", Other, "world get is a CLI command the app does not run; at a terminal: ailang-worldd world get world-1", true, []string{"world-1"}},
+	{"world get a b", "world get", Other, "world get is a CLI command", true, []string{"a", "b"}},
+	{"--addr x episode report w1", "", Refused, "global flags are not accepted", true, nil},
+	{"ailang-worldd -h", "", Refused, "global flags are not accepted", true, nil},
+	{"episode", "episode", Refused, EpisodeNeedsSubcommand, true, []string{}},
+	{"episode frobnicate", "episode", Refused, EpisodeNeedsSubcommand, true, []string{"frobnicate"}},
+	{"episode w1", "episode", Refused, EpisodeNeedsSubcommand, true, []string{"w1"}},
+	{"world a b c", "world", Refused, AddressHelp, false, []string{"a", "b", "c"}},
+}
+
+// VB3: every row of the table and of the trace (design §2.1(b), §2.4), the shadow set, step 4b.
+func TestLookup(t *testing.T) {
+	for _, c := range lookupCases {
+		argv, err := Split(c.line)
+		if err != nil {
+			t.Fatal(err)
+		}
+		if back, err := Split(Quote(argv)); err != nil || !reflect.DeepEqual(back, argv) { // VB2 over every VB3 argv
+			t.Errorf("Split(Quote(%q)) = %q, %v", argv, back, err)
+		}
+		v, rest := Lookup(argv)
+		cli := strings.TrimPrefix(c.line, "ailang-worldd ")
+		if v.Words != c.words || v.Class != c.class || !strings.Contains(v.Reason, c.reason) || (c.reason == "" && v.Reason != "") ||
+			(v.CLI != "") != c.cli || (c.cli && v.CLI != "ailang-worldd "+cli) || !reflect.DeepEqual(rest, c.rest) {
+			t.Errorf("Lookup(%q) = %+v, %q", c.line, v, rest)
+		}
+	}
+	// The shadow set: every one of the 27 words is a verb, bare and with an argument; never an address.
+	all := append(append(append([]string{}, OperatorWords...), NeverWords...), OtherWords...)
+	if len(all) != 27 {
+		t.Fatalf("shadow set has %d words, want 27", len(all))
+	}
+	for _, w := range all {
+		for _, argv := range [][]string{{w}, {w, "zzz"}, {Binary, w, "zzz"}} {
+			if v, _ := Lookup(argv); v.Class == Address || v.CLI == "" {
+				t.Errorf("Lookup(%q) = %+v: a shadow-set word navigates or has no CLI line", argv, v)
+			}
+		}
+		if v, _ := Lookup([]string{"world", w}); v.Class != Address {
+			t.Errorf("world %s = %+v: a shadowed name is not reachable as world NAME", w, v)
+		}
+	}
+}
+
+// ApproveTarget: an id or a ref opens its card; anything else is refused.
+func TestApproveTarget(t *testing.T) {
+	ref := "sha256:" + strings.Repeat("ab", 32)
+	for args, want := range map[string]string{"mr-0123456789abcdef": "/surface/approve/mr-0123456789abcdef", ref: "/surface/decide/" + ref,
+		"mr-0123456789abcdef --deny": "", "--list": "", "": "", ref + "0": "", "sha256:" + strings.Repeat("AB", 32): "", "mr-0123456789abcdef x": "", "mr-0123456789ABCDEF": ""} {
+		href, reason := ApproveTarget(strings.Fields(args))
+		if href != want || (href == "") == (reason == "") {
+			t.Errorf("ApproveTarget(%q) = %q, %q", args, href, reason)
+		}
+	}
+}
````

### Patch M1c — one flag parse: the ten call sites, `EpisodeReportFlags` (VB5)

````diff M1c
diff --git a/cmd/ailang-worldd/approve.go b/cmd/ailang-worldd/approve.go
index cc6e8c5..b793bf0 100644
--- a/cmd/ailang-worldd/approve.go
+++ b/cmd/ailang-worldd/approve.go
@@ -19,6 +19,7 @@ import (
 	"github.com/sunholo-data/ailang-world/host/approval"
 	"github.com/sunholo-data/ailang-world/host/authority"
 	"github.com/sunholo-data/ailang-world/host/store"
+	"github.com/sunholo-data/ailang-world/host/verbs"
 )
 
 const approveHelp = `usage: ailang-worldd approve <request-id> [--db <path>] [--credentials-db <path>]
@@ -181,7 +182,7 @@ func runApprove(args []string, stdout, stderr io.Writer, env sessionEnv) int {
 	deny := fs.Bool("deny", false, "with a request ref: deny it instead of approving it")
 	daemonAddr := fs.String("daemon", "", "with a request ref: the daemon serving --db (default: the installed service's address)")
 	expiry := fs.Int64("expiry", 0, "with a request ref and --signature or --webauthn-assertion: the expiry the signed decision message carries")
-	positional, err := parseInterleaved(fs, args)
+	positional, err := verbs.ParseInterleaved(fs, args)
 	if errors.Is(err, flag.ErrHelp) {
 		fmt.Fprint(stdout, approveHelp)
 		return exitOK
diff --git a/cmd/ailang-worldd/approver.go b/cmd/ailang-worldd/approver.go
index 841d700..1a2e30a 100644
--- a/cmd/ailang-worldd/approver.go
+++ b/cmd/ailang-worldd/approver.go
@@ -20,6 +20,7 @@ import (
 
 	"github.com/sunholo-data/ailang-world/host/approval"
 	"github.com/sunholo-data/ailang-world/host/store"
+	"github.com/sunholo-data/ailang-world/host/verbs"
 )
 
 // ---- M3a ----
@@ -216,7 +217,7 @@ func approverFlags(name string, args []string, stdout, stderr io.Writer, extra f
 	if extra != nil {
 		extra(fs)
 	}
-	positional, err := parseInterleaved(fs, args)
+	positional, err := verbs.ParseInterleaved(fs, args)
 	if errors.Is(err, flag.ErrHelp) {
 		fmt.Fprint(stdout, approverHelp)
 		return "", "", nil, exitOK, true
diff --git a/cmd/ailang-worldd/episode.go b/cmd/ailang-worldd/episode.go
index 3200082..2f803a6 100644
--- a/cmd/ailang-worldd/episode.go
+++ b/cmd/ailang-worldd/episode.go
@@ -45,6 +45,7 @@ import (
 	"github.com/sunholo-data/ailang-world/host/daemon"
 	"github.com/sunholo-data/ailang-world/host/pins"
 	"github.com/sunholo-data/ailang-world/host/store"
+	"github.com/sunholo-data/ailang-world/host/verbs"
 )
 
 const episodeHelp = `episode finish <episode> --preset <name>: attended local checks, commit and body; see episode finish --help.
@@ -214,7 +215,7 @@ func runEpisodeStart(args []string, stdout, stderr io.Writer, env sessionEnv) in
 	budget := fs.Int64("budget", 50, "budget of each preset grant")
 	sessionsDir := fs.String("sessions-dir", defaultSessionsDir(), "where the session file and the MCP config are written")
 	addr := fs.String("addr", defaultServiceBind(), "the serving daemon, host:port")
-	positional, err := parseInterleaved(fs, args)
+	positional, err := verbs.ParseInterleaved(fs, args)
 	if err != nil {
 		if errors.Is(err, flag.ErrHelp) {
 			fmt.Fprint(stdout, episodeHelp)
diff --git a/cmd/ailang-worldd/episode_finish.go b/cmd/ailang-worldd/episode_finish.go
index 1d20909..01b6c9f 100644
--- a/cmd/ailang-worldd/episode_finish.go
+++ b/cmd/ailang-worldd/episode_finish.go
@@ -21,6 +21,7 @@ import (
 	"github.com/sunholo-data/ailang-world/host/checkrun"
 	"github.com/sunholo-data/ailang-world/host/daemon"
 	"github.com/sunholo-data/ailang-world/host/store"
+	"github.com/sunholo-data/ailang-world/host/verbs"
 	"github.com/sunholo-data/ailang-world/host/worktree"
 )
 
@@ -68,7 +69,7 @@ func runEpisodeFinishWith(args []string, stdout, stderr io.Writer, env sessionEn
 	dry := fs.Bool("dry-run", false, "print then reset local commit")
 	var allowed multiFlag
 	fs.Var(&allowed, "allow-failing", "allow this failing check")
-	pos, err := parseInterleaved(fs, args)
+	pos, err := verbs.ParseInterleaved(fs, args)
 	if errors.Is(err, flag.ErrHelp) {
 		fmt.Fprint(stdout, episodeFinishHelp)
 		return exitOK
diff --git a/cmd/ailang-worldd/episode_launch.go b/cmd/ailang-worldd/episode_launch.go
index 79c3081..ddfeffd 100644
--- a/cmd/ailang-worldd/episode_launch.go
+++ b/cmd/ailang-worldd/episode_launch.go
@@ -25,6 +25,7 @@ import (
 	"github.com/sunholo-data/ailang-world/host/childenv"
 	"github.com/sunholo-data/ailang-world/host/daemon"
 	"github.com/sunholo-data/ailang-world/host/store"
+	"github.com/sunholo-data/ailang-world/host/verbs"
 )
 
 const episodeLaunchHelp = `usage: ailang-worldd episode launch <episode> --preset <name> [--presets-dir D] [--hosts-dir D] [--headless]
@@ -234,7 +235,7 @@ func runEpisodeLaunch(args []string, stdout, stderr io.Writer) int {
 	dir := fs.String("presets-dir", defaultPresetsDir(), "private presets")
 	headless := fs.Bool("headless", false, "headless launch")
 	hostsDir := fs.String("hosts-dir", defaultHostsDir(), "hosts init output directory")
-	pos, err := parseInterleaved(fs, args)
+	pos, err := verbs.ParseInterleaved(fs, args)
 	if errors.Is(err, flag.ErrHelp) {
 		fmt.Fprint(stdout, episodeLaunchHelp)
 		return exitOK
diff --git a/cmd/ailang-worldd/episode_report.go b/cmd/ailang-worldd/episode_report.go
index 25030ae..515c574 100644
--- a/cmd/ailang-worldd/episode_report.go
+++ b/cmd/ailang-worldd/episode_report.go
@@ -21,6 +21,7 @@ import (
 	"time"
 
 	"github.com/sunholo-data/ailang-world/host/daemon"
+	"github.com/sunholo-data/ailang-world/host/verbs"
 )
 
 const episodeReportHelp = `usage: ailang-worldd [--addr <url>] episode report <episode> [--json | --markdown] [--follow] [--interval 2s] [--addr <url>]
@@ -50,12 +51,9 @@ func runEpisodeReport(globalAddr string, args []string, stdout, stderr io.Writer
 	fs := flag.NewFlagSet("ailang-worldd episode report", flag.ContinueOnError)
 	fs.SetOutput(stderr)
 	fs.Usage = func() {}
-	markdown := fs.Bool("markdown", false, "render a scrubbed PR body")
-	asJSON := fs.Bool("json", false, "print the route's JSON")
-	follow := fs.Bool("follow", false, "re-print when the episode changes, until Ctrl-C")
-	interval := fs.Duration("interval", 2*time.Second, "poll interval with --follow")
-	addr := fs.String("addr", "", "the daemon (default: the global --addr)")
-	positional, err := parseInterleaved(fs, args)
+	flags := verbs.EpisodeReportFlags(fs) // row 243: the bar's admission parses the same set
+	markdown, asJSON, follow, interval, addr := flags.Markdown, flags.JSON, flags.Follow, flags.Interval, flags.Addr
+	positional, err := verbs.ParseInterleaved(fs, args)
 	if err != nil {
 		if errors.Is(err, flag.ErrHelp) {
 			fmt.Fprint(stdout, episodeReportHelp)
diff --git a/cmd/ailang-worldd/preset.go b/cmd/ailang-worldd/preset.go
index 9c54de4..85273dd 100644
--- a/cmd/ailang-worldd/preset.go
+++ b/cmd/ailang-worldd/preset.go
@@ -18,6 +18,7 @@ import (
 
 	"github.com/sunholo-data/ailang-world/host/broker"
 	"github.com/sunholo-data/ailang-world/host/store"
+	"github.com/sunholo-data/ailang-world/host/verbs"
 	"github.com/sunholo-data/ailang-world/host/worktree"
 )
 
@@ -528,7 +529,7 @@ func runPreset(args []string, out, errOut io.Writer) int {
 	remote := fs.String("remote", "origin", "remote name")
 	agent := fs.String("agent", "codex", "agent")
 	profile := fs.String("profile", "", "exec profile")
-	pos, err := parseInterleaved(fs, args[1:])
+	pos, err := verbs.ParseInterleaved(fs, args[1:])
 	fail := func(err error) int { fmt.Fprintln(errOut, err); return exitUsage }
 	if errors.Is(err, flag.ErrHelp) {
 		fmt.Fprint(out, presetHelp)
diff --git a/cmd/ailang-worldd/session.go b/cmd/ailang-worldd/session.go
index 549ed3c..a274755 100644
--- a/cmd/ailang-worldd/session.go
+++ b/cmd/ailang-worldd/session.go
@@ -16,6 +16,7 @@ import (
 	"github.com/sunholo-data/ailang-world/host/authority"
 	"github.com/sunholo-data/ailang-world/host/broker"
 	"github.com/sunholo-data/ailang-world/host/store"
+	"github.com/sunholo-data/ailang-world/host/verbs"
 )
 
 // sessionEnv is the injectable human-act seam a session verb performs. It lets
@@ -176,22 +177,6 @@ func openCredentialsForSession(ctx context.Context, db, creds string) (*store.Cr
 	return st, nil
 }
 
-// parseInterleaved parses fs over args, allowing flags and positional
-// arguments in any order (Go's flag package stops at the first positional).
-func parseInterleaved(fs *flag.FlagSet, args []string) ([]string, error) {
-	var positional []string
-	for {
-		if err := fs.Parse(args); err != nil {
-			return nil, err
-		}
-		if fs.NArg() == 0 {
-			return positional, nil
-		}
-		positional = append(positional, fs.Arg(0))
-		args = fs.Args()[1:]
-	}
-}
-
 // multiFlag is a repeatable string flag value (used by the repeatable --grant).
 type multiFlag []string
 
@@ -369,7 +354,7 @@ func runSessionRevoke(args []string, stdout, stderr io.Writer) int {
 	dbPath := fs.String("db", "", "world store database (required)")
 	credsDB := fs.String("credentials-db", "", credentialsDBUsage)
 	episode := fs.String("episode", "", "revoke every credential of this episode")
-	rest, err := parseInterleaved(fs, args)
+	rest, err := verbs.ParseInterleaved(fs, args)
 	if err != nil {
 		if errors.Is(err, flag.ErrHelp) {
 			fmt.Fprint(stdout, sessionRevokeHelp)
diff --git a/cmd/ailang-worldd/session_new.go b/cmd/ailang-worldd/session_new.go
index fb420c2..784b5f8 100644
--- a/cmd/ailang-worldd/session_new.go
+++ b/cmd/ailang-worldd/session_new.go
@@ -65,6 +65,7 @@ import (
 	"github.com/sunholo-data/ailang-world/host/authority"
 	"github.com/sunholo-data/ailang-world/host/broker"
 	"github.com/sunholo-data/ailang-world/host/store"
+	"github.com/sunholo-data/ailang-world/host/verbs"
 	"github.com/sunholo-data/ailang-world/host/worktree"
 )
 
@@ -607,7 +608,7 @@ func runSessionNew(args []string, stdout, stderr io.Writer, env sessionEnv) int
 	outPath := fs.String("out", "", "write the token once to this file (0600)")
 	profile := fs.String("profile", "", "bind the session to this exec profile's project")
 	profilesDir := fs.String("exec-profiles-dir", defaultExecProfilesDir(), "where --profile looks")
-	positional, err := parseInterleaved(fs, args)
+	positional, err := verbs.ParseInterleaved(fs, args)
 	if err != nil {
 		if errors.Is(err, flag.ErrHelp) {
 			fmt.Fprint(stdout, sessionNewHelp)
@@ -781,7 +782,7 @@ func runSessionList(args []string, stdout, stderr io.Writer, env sessionEnv) int
 	wsRoot := fs.String("workspace-root", def.workspaceRoot, "workspace root (marks degraded episodes)")
 	episode := fs.String("episode", "", "only this episode")
 	asJSON := fs.Bool("json", false, "print JSON")
-	positional, err := parseInterleaved(fs, args)
+	positional, err := verbs.ParseInterleaved(fs, args)
 	if err != nil {
 		if errors.Is(err, flag.ErrHelp) {
 			fmt.Fprint(stdout, sessionListHelp)
diff --git a/host/verbs/flags.go b/host/verbs/flags.go
new file mode 100644
index 0000000..8792023
--- /dev/null
+++ b/host/verbs/flags.go
@@ -0,0 +1,77 @@
+package verbs
+
+import (
+	"errors"
+	"flag"
+	"fmt"
+	"io"
+	"time"
+)
+
+// ParseInterleaved parses fs over args, allowing flags and positional
+// arguments in any order (Go's flag package stops at the first positional).
+func ParseInterleaved(fs *flag.FlagSet, args []string) ([]string, error) {
+	var positional []string
+	for {
+		if err := fs.Parse(args); err != nil {
+			return nil, err
+		}
+		if fs.NArg() == 0 {
+			return positional, nil
+		}
+		positional = append(positional, fs.Arg(0))
+		args = fs.Args()[1:]
+	}
+}
+
+// EpisodeReportArgs are episode report's flags.
+type EpisodeReportArgs struct {
+	Markdown, JSON, Follow *bool
+	Interval               *time.Duration
+	Addr                   *string
+}
+
+// EpisodeReportFlags registers episode report's five flags on fs; the CLI and the bar share it.
+func EpisodeReportFlags(fs *flag.FlagSet) *EpisodeReportArgs {
+	return &EpisodeReportArgs{
+		Markdown: fs.Bool("markdown", false, "render a scrubbed PR body"),
+		JSON:     fs.Bool("json", false, "print the route's JSON"),
+		Follow:   fs.Bool("follow", false, "re-print when the episode changes, until Ctrl-C"),
+		Interval: fs.Duration("interval", 2*time.Second, "poll interval with --follow"),
+		Addr:     fs.String("addr", "", "the daemon (default: the global --addr)"),
+	}
+}
+
+// barRefusal is the bar profile's reason for each flag the bar does not run.
+var barRefusal = []struct{ flag, why string }{
+	{"follow", "the app shows one report; it re-reads on Enter"},
+	{"interval", "the app shows one report; it re-reads on Enter"},
+	{"markdown", "the PR body is rendered by episode finish"},
+	{"addr", "the app reads its own daemon"},
+}
+
+// AdmitEpisodeReport runs episode report's flag set and parse over args, then the bar profile:
+// --json is allowed, every other flag and -h/--help refused, exactly one episode required.
+func AdmitEpisodeReport(args []string) (episode string, asJSON bool, reason string) {
+	fs := flag.NewFlagSet(Binary+" episode report", flag.ContinueOnError)
+	fs.SetOutput(io.Discard)
+	a := EpisodeReportFlags(fs)
+	pos, err := ParseInterleaved(fs, args)
+	if errors.Is(err, flag.ErrHelp) {
+		return "", false, "see " + Binary + " episode report --help at a terminal"
+	}
+	if err != nil {
+		return "", false, Binary + " episode report: " + err.Error()
+	}
+	set := map[string]bool{}
+	fs.Visit(func(f *flag.Flag) { set[f.Name] = true })
+	for _, r := range barRefusal {
+		if set[r.flag] {
+			return "", false, fmt.Sprintf("--%s is not run from the app: %s", r.flag, r.why)
+		}
+	}
+	if len(pos) != 1 || pos[0] == "" {
+		return "", false, Binary + " episode report needs exactly one <episode>"
+	}
+	return pos[0], *a.JSON, ""
+}
diff --git a/host/verbs/flags_test.go b/host/verbs/flags_test.go
new file mode 100644
index 0000000..4ae6111
--- /dev/null
+++ b/host/verbs/flags_test.go
@@ -0,0 +1,58 @@
+package verbs
+
+// Row 243 V1 (design §6 VB5): episode report's one flag set and the bar's admission.
+
+import (
+	"flag"
+	"reflect"
+	"strings"
+	"testing"
+	"time"
+)
+
+// VB5: the CLI's flag set exactly, and the bar profile's one allow and its refusals.
+func TestEpisodeReportFlags(t *testing.T) {
+	fs := flag.NewFlagSet("x", flag.ContinueOnError)
+	EpisodeReportFlags(fs)
+	got := map[string]string{}
+	fs.VisitAll(func(f *flag.Flag) { got[f.Name] = f.DefValue })
+	if want := map[string]string{"markdown": "false", "json": "false", "follow": "false", "interval": (2 * time.Second).String(), "addr": ""}; !reflect.DeepEqual(got, want) {
+		t.Errorf("flags = %v, want %v", got, want)
+	}
+	for _, c := range []struct {
+		args          []string
+		episode       string
+		json          bool
+		reasonContain string
+	}{
+		{[]string{"w1"}, "w1", false, ""},
+		{[]string{"--json", "w1"}, "w1", true, ""},
+		{[]string{"w1", "--json=false"}, "w1", false, ""},
+		{[]string{"w1", "--follow"}, "", false, "--follow is not run from the app: the app shows one report"},
+		{[]string{"w1", "--interval", "5s"}, "", false, "--interval is not run from the app"},
+		{[]string{"w1", "--markdown"}, "", false, "--markdown is not run from the app: the PR body"},
+		{[]string{"w1", "--addr", "http://x"}, "", false, "--addr is not run from the app: the app reads its own daemon"},
+		{[]string{"w1", "-h"}, "", false, "see ailang-worldd episode report --help"},
+		{[]string{"--help"}, "", false, "see ailang-worldd episode report --help"},
+		{[]string{"w1", "--frob"}, "", false, "flag provided but not defined: -frob"},
+		{[]string{"w1", "--interval", "x"}, "", false, "invalid value"},
+		{[]string{}, "", false, "needs exactly one <episode>"},
+		{[]string{"w1", "w2"}, "", false, "needs exactly one <episode>"},
+		{[]string{""}, "", false, "needs exactly one <episode>"},
+	} {
+		ep, asJSON, reason := AdmitEpisodeReport(c.args)
+		if ep != c.episode || asJSON != c.json || !strings.Contains(reason, c.reasonContain) || (c.reasonContain == "") != (reason == "") {
+			t.Errorf("AdmitEpisodeReport(%q) = %q, %v, %q", c.args, ep, asJSON, reason)
+		}
+	}
+}
+
+// ParseInterleaved keeps the CLI's order-free parse (its CLI callers' tests pin the rest).
+func TestParseInterleaved(t *testing.T) {
+	fs := flag.NewFlagSet("x", flag.ContinueOnError)
+	j := fs.Bool("json", false, "")
+	pos, err := ParseInterleaved(fs, []string{"a", "--json", "b"})
+	if err != nil || !*j || !reflect.DeepEqual(pos, []string{"a", "b"}) {
+		t.Errorf("= %q %v %v", pos, *j, err)
+	}
+}
````

### Patch M2 — `episode report`'s printer, moved

````diff M2
diff --git a/cmd/ailang-worldd/episode_markdown.go b/cmd/ailang-worldd/episode_markdown.go
index b44029a..2c8b09f 100644
--- a/cmd/ailang-worldd/episode_markdown.go
+++ b/cmd/ailang-worldd/episode_markdown.go
@@ -83,16 +83,16 @@ func renderReportMarkdown(r daemon.EpisodeReport, db string, add reportAdditions
 	introEnd := b.Len()
 	b.WriteString("## What the agent did\n\n")
 	fmt.Fprintf(&b, "%d calls.\n\nCalls by tool:\n", r.TotalCalls)
-	for _, kv := range sortedCounts(r.CallsByTool) {
-		fmt.Fprintf(&b, "- %s: %d\n", kv.name, kv.n)
+	for _, kv := range daemon.SortedCounts(r.CallsByTool) {
+		fmt.Fprintf(&b, "- %s: %d\n", kv.Name, kv.N)
 	}
 	b.WriteString("\nEffects by decision:\n")
 	for _, e := range r.Effects {
 		fmt.Fprintf(&b, "- %s %s: %d\n", e.Effect, e.Decision, e.Count)
 	}
 	b.WriteString("\n")
-	printRefusals(&b, "Plan refusals", r.PlanRefusals)
-	printRefusals(&b, "Handler refusals", r.HandlerRefusals)
+	daemon.WriteReportRefusals(&b, "Plan refusals", r.PlanRefusals)
+	daemon.WriteReportRefusals(&b, "Handler refusals", r.HandlerRefusals)
 	b.WriteString("\nBudgets (used / left):\n")
 	for _, v := range r.Budgets {
 		fmt.Fprintf(&b, "- %s@%s: %d / %d\n", v.Effect, v.Scope, v.Used, v.Left)
diff --git a/cmd/ailang-worldd/episode_report.go b/cmd/ailang-worldd/episode_report.go
index 515c574..c1d7d89 100644
--- a/cmd/ailang-worldd/episode_report.go
+++ b/cmd/ailang-worldd/episode_report.go
@@ -16,7 +16,6 @@ import (
 	"io"
 	"net/http"
 	"net/url"
-	"sort"
 	"strings"
 	"time"
 
@@ -219,95 +218,5 @@ func printFollowed(w io.Writer, o reportOpts, body []byte, rep daemon.EpisodeRep
 	printReport(w, rep)
 }
 
-// printReport renders the compact human table.
-func printReport(w io.Writer, r daemon.EpisodeReport) {
-	if r.Start != nil {
-		data, _ := json.MarshalIndent(r.Start, "  ", "  ")
-		fmt.Fprintf(w, "start:\n  %s\n", data)
-	}
-	fmt.Fprintf(w, "episode %s: %d calls, entries %d-%d, %s .. %s\n",
-		r.Episode, r.TotalCalls, r.Entries[0], r.Entries[1], r.FirstAt, r.LastAt)
-	fmt.Fprintln(w, "calls by tool:")
-	for _, kv := range sortedCounts(r.CallsByTool) {
-		fmt.Fprintf(w, "  %-18s %6d\n", kv.name, kv.n)
-	}
-	fmt.Fprintln(w, readsLine(r.Reads))
-	fmt.Fprintln(w, "effects:")
-	for _, e := range r.Effects {
-		fmt.Fprintf(w, "  %-22s %-16s %6d\n", e.Effect, e.Decision, e.Count)
-	}
-	fmt.Fprintln(w, "budgets (minted / used / left):")
-	for _, b := range r.Budgets {
-		fmt.Fprintf(w, "  %-22s %6d %6d %6d\n", b.Effect+"@"+b.Scope, b.Minted, b.Used, b.Left)
-	}
-	printRefusals(w, "plan refusals (refused before any effect)", r.PlanRefusals)
-	printRefusals(w, "handler refusals (effect allowed, handler ok:false)", r.HandlerRefusals)
-	if r.Uncommitted > 0 {
-		fmt.Fprintf(w, "uncommitted intents: %d\n", r.Uncommitted)
-	}
-	for _, e := range r.ExecCache {
-		fmt.Fprintln(w, execCacheLine(e))
-	}
-}
-
-// execCacheLine is one exec_cache entry of the report table (row 257).
-func execCacheLine(e daemon.ReportExecCache) string {
-	line := fmt.Sprintf("exec cache   %s %s %s", e.Project, e.Cache, e.State)
-	if e.Attempts > 0 {
-		line += fmt.Sprintf(" at %s (attempt %d of 3)", e.By, e.Attempts)
-	}
-	if e.State == "ready" && e.Attempts > 0 {
-		line += fmt.Sprintf(" in %dms", e.TookMS)
-	}
-	if e.Failure != "" {
-		line += ": " + e.Failure
-	}
-	return line
-}
-
-// readsLine is the one-line reads summary; --json carries the file list.
-func readsLine(r daemon.ReportReads) string {
-	line := fmt.Sprintf("reads: %d calls, %d files", r.Count, len(r.Files))
-	if r.Truncated > 0 {
-		line += fmt.Sprintf(", %d more not listed", r.Truncated)
-	}
-	if r.Unlisted > 0 {
-		line += fmt.Sprintf(", %d outside the worktree", r.Unlisted)
-	}
-	return line
-}
-
-func printRefusals(w io.Writer, title string, r daemon.ReportRefusals) {
-	parts := make([]string, 0, len(r.ByTool))
-	for _, kv := range sortedCounts(r.ByTool) {
-		parts = append(parts, fmt.Sprintf("%s %d", kv.name, kv.n))
-	}
-	line := fmt.Sprintf("%s: %d", title, r.Count)
-	if len(parts) > 0 {
-		line += " (" + strings.Join(parts, ", ") + ")"
-	}
-	fmt.Fprintln(w, line)
-	for _, f := range r.First {
-		fmt.Fprintf(w, "  #%d %s: %s\n", f.Entry, f.Tool, f.Reason)
-	}
-}
-
-type namedCount struct {
-	name string
-	n    int
-}
-
-// sortedCounts orders a count map by count, then name.
-func sortedCounts(m map[string]int) []namedCount {
-	out := make([]namedCount, 0, len(m))
-	for k, n := range m {
-		out = append(out, namedCount{k, n})
-	}
-	sort.Slice(out, func(i, j int) bool {
-		if out[i].n != out[j].n {
-			return out[i].n > out[j].n
-		}
-		return out[i].name < out[j].name
-	})
-	return out
-}
+// printReport renders the compact human table: the daemon's own text, which the app's bar shows too.
+func printReport(w io.Writer, r daemon.EpisodeReport) { daemon.WriteEpisodeReportText(w, r) }
diff --git a/host/daemon/episode_report_text.go b/host/daemon/episode_report_text.go
new file mode 100644
index 0000000..dd859ae
--- /dev/null
+++ b/host/daemon/episode_report_text.go
@@ -0,0 +1,105 @@
+package daemon
+
+import (
+	"encoding/json"
+	"fmt"
+	"io"
+	"sort"
+	"strings"
+)
+
+// WriteEpisodeReportText renders the compact human table: episode report's text at a terminal and
+// in the app's bar (row 243 V1, moved from cmd/ailang-worldd so both print one function's bytes).
+func WriteEpisodeReportText(w io.Writer, r EpisodeReport) {
+	if r.Start != nil {
+		data, _ := json.MarshalIndent(r.Start, "  ", "  ")
+		fmt.Fprintf(w, "start:\n  %s\n", data)
+	}
+	fmt.Fprintf(w, "episode %s: %d calls, entries %d-%d, %s .. %s\n",
+		r.Episode, r.TotalCalls, r.Entries[0], r.Entries[1], r.FirstAt, r.LastAt)
+	fmt.Fprintln(w, "calls by tool:")
+	for _, kv := range SortedCounts(r.CallsByTool) {
+		fmt.Fprintf(w, "  %-18s %6d\n", kv.Name, kv.N)
+	}
+	fmt.Fprintln(w, readsLine(r.Reads))
+	fmt.Fprintln(w, "effects:")
+	for _, e := range r.Effects {
+		fmt.Fprintf(w, "  %-22s %-16s %6d\n", e.Effect, e.Decision, e.Count)
+	}
+	fmt.Fprintln(w, "budgets (minted / used / left):")
+	for _, b := range r.Budgets {
+		fmt.Fprintf(w, "  %-22s %6d %6d %6d\n", b.Effect+"@"+b.Scope, b.Minted, b.Used, b.Left)
+	}
+	WriteReportRefusals(w, "plan refusals (refused before any effect)", r.PlanRefusals)
+	WriteReportRefusals(w, "handler refusals (effect allowed, handler ok:false)", r.HandlerRefusals)
+	if r.Uncommitted > 0 {
+		fmt.Fprintf(w, "uncommitted intents: %d\n", r.Uncommitted)
+	}
+	for _, e := range r.ExecCache {
+		fmt.Fprintln(w, execCacheLine(e))
+	}
+}
+
+// execCacheLine is one exec_cache entry of the report table (row 257).
+func execCacheLine(e ReportExecCache) string {
+	line := fmt.Sprintf("exec cache   %s %s %s", e.Project, e.Cache, e.State)
+	if e.Attempts > 0 {
+		line += fmt.Sprintf(" at %s (attempt %d of 3)", e.By, e.Attempts)
+	}
+	if e.State == "ready" && e.Attempts > 0 {
+		line += fmt.Sprintf(" in %dms", e.TookMS)
+	}
+	if e.Failure != "" {
+		line += ": " + e.Failure
+	}
+	return line
+}
+
+// readsLine is the one-line reads summary; --json carries the file list.
+func readsLine(r ReportReads) string {
+	line := fmt.Sprintf("reads: %d calls, %d files", r.Count, len(r.Files))
+	if r.Truncated > 0 {
+		line += fmt.Sprintf(", %d more not listed", r.Truncated)
+	}
+	if r.Unlisted > 0 {
+		line += fmt.Sprintf(", %d outside the worktree", r.Unlisted)
+	}
+	return line
+}
+
+// WriteReportRefusals is one refusals line and its first entries.
+func WriteReportRefusals(w io.Writer, title string, r ReportRefusals) {
+	parts := make([]string, 0, len(r.ByTool))
+	for _, kv := range SortedCounts(r.ByTool) {
+		parts = append(parts, fmt.Sprintf("%s %d", kv.Name, kv.N))
+	}
+	line := fmt.Sprintf("%s: %d", title, r.Count)
+	if len(parts) > 0 {
+		line += " (" + strings.Join(parts, ", ") + ")"
+	}
+	fmt.Fprintln(w, line)
+	for _, f := range r.First {
+		fmt.Fprintf(w, "  #%d %s: %s\n", f.Entry, f.Tool, f.Reason)
+	}
+}
+
+// NamedCount is one name and its count.
+type NamedCount struct {
+	Name string
+	N    int
+}
+
+// SortedCounts orders a count map by count, then name.
+func SortedCounts(m map[string]int) []NamedCount {
+	out := make([]NamedCount, 0, len(m))
+	for k, n := range m {
+		out = append(out, NamedCount{k, n})
+	}
+	sort.Slice(out, func(i, j int) bool {
+		if out[i].N != out[j].N {
+			return out[i].N > out[j].N
+		}
+		return out[i].Name < out[j].Name
+	})
+	return out
+}
````

### Patch M3a — `GET /surface/verb.json`, `d.episodeReport` (VR1, VR4, VS1)

````diff M3a
diff --git a/host/daemon/daemon.go b/host/daemon/daemon.go
index d443564..8a2c0d4 100644
--- a/host/daemon/daemon.go
+++ b/host/daemon/daemon.go
@@ -1145,6 +1145,7 @@ func (d *Daemon) Handler() http.Handler {
 	mux.HandleFunc("POST /surface/decide/{ref}", d.approvalRoute(d.handleSurfaceDecide))
 	mux.HandleFunc("GET /surface/decide/{ref}", d.approvalRoute(d.handleSurfaceDecidePage))
 	mux.HandleFunc("GET /surface/decide.js", d.approvalRoute(d.handleSurfaceDecideScript))
+	mux.HandleFunc("GET /surface/verb.json", d.approvalRoute(d.handleSurfaceVerb)) // row 243 V1: the bar's verbs, read-only
 	// The two A2A projection routes (w-a2a-session-projection P6.B-A2A-CARD)
 	// are ADDITIVE: the frozen /v1/ table above is untouched, and the routes
 	// are NOT in isProtected — the projection handler resolves the session
diff --git a/host/daemon/episode_report.go b/host/daemon/episode_report.go
index 841ba23..27d393c 100644
--- a/host/daemon/episode_report.go
+++ b/host/daemon/episode_report.go
@@ -31,6 +31,7 @@ package daemon
 import (
 	"context"
 	"encoding/json"
+	"errors"
 	"fmt"
 	"net/http"
 	"path/filepath"
@@ -712,30 +713,51 @@ func (d *Daemon) handleEpisodeReport(w http.ResponseWriter, r *http.Request) {
 	}
 	ctx, cancel := d.readCtx(r)
 	defer cancel()
-	rep, ok, err := buildEpisodeReport(ctx, d.intents, d.reads, episode)
+	rep, found, err := d.episodeReport(ctx, episode)
+	if reportTimedOut(ctx, err) {
+		writeReadTimeout(w, d.readDeadline)
+		return
+	}
 	if err != nil {
-		if timedOut(ctx, err) {
-			writeReadTimeout(w, d.readDeadline)
-			return
-		}
 		d.writeInternalError(w, r, err)
 		return
 	}
-	start, startErr := d.readEpisodeStart(ctx, episode)
-	if startErr != nil {
-		d.writeInternalError(w, r, startErr)
+	if !found {
+		writeAPIError(w, "NotFound", "episode "+episode+" has no committed calls", http.StatusNotFound)
 		return
 	}
+	writeJSON(w, http.StatusOK, rep)
+}
+
+// reportBuildError is the builder's own failure: the one the report maps to a read timeout.
+type reportBuildError struct{ err error }
+
+func (e reportBuildError) Error() string { return e.err.Error() }
+func (e reportBuildError) Unwrap() error { return e.err }
+
+// reportTimedOut reports a builder failure past the read deadline (the start and finish reads'
+// failures stay internal errors, as the route has always answered them).
+func reportTimedOut(ctx context.Context, err error) bool {
+	var build reportBuildError
+	return errors.As(err, &build) && timedOut(ctx, build.err)
+}
+
+// episodeReport is the report GET /v1/episodes/{episode}/report serves and the app's bar runs (row
+// 243 V1: extracted unchanged from the handler); found is false when the episode has nothing.
+func (d *Daemon) episodeReport(ctx context.Context, episode string) (EpisodeReport, bool, error) {
+	rep, ok, err := buildEpisodeReport(ctx, d.intents, d.reads, episode)
+	if err != nil {
+		return rep, false, reportBuildError{err}
+	}
+	start, err := d.readEpisodeStart(ctx, episode)
+	if err != nil {
+		return rep, false, err
+	}
 	rep.Start = start
 	rep.Finish, err = d.readEpisodeFinish(ctx, episode)
 	if err != nil {
-		d.writeInternalError(w, r, err)
-		return
+		return rep, false, err
 	}
 	rep.ExecCache = d.workspace.execCacheStatus(episode)
-	if !ok && start == nil && rep.Finish == nil && len(rep.ExecCache) == 0 {
-		writeAPIError(w, "NotFound", "episode "+episode+" has no committed calls", http.StatusNotFound)
-		return
-	}
-	writeJSON(w, http.StatusOK, rep)
+	return rep, ok || start != nil || rep.Finish != nil || len(rep.ExecCache) > 0, nil
 }
diff --git a/host/daemon/surface_verb.go b/host/daemon/surface_verb.go
new file mode 100644
index 0000000..8fb95e4
--- /dev/null
+++ b/host/daemon/surface_verb.go
@@ -0,0 +1,103 @@
+package daemon
+
+import (
+	"bytes"
+	"encoding/json"
+	"fmt"
+	"net/http"
+	"strings"
+
+	"github.com/sunholo-data/ailang-world/host/verbs"
+)
+
+// VerbSchema names GET /surface/verb.json's answer (row 243 V1, design §2.2).
+const VerbSchema = "world/surface-verb/v1"
+
+// verbAnswer is one answer of the verb route. Every answer carries v, schema, q, argv, verb, class
+// and outcome; the rest belong to their outcome.
+type verbAnswer struct {
+	V       int      `json:"v"`
+	Schema  string   `json:"schema"`
+	Q       string   `json:"q"`
+	Argv    []string `json:"argv"`
+	Verb    string   `json:"verb"`
+	Class   string   `json:"class"`
+	Outcome string   `json:"outcome"`
+	Reason  string   `json:"reason,omitempty"`
+	CLI     string   `json:"cli,omitempty"`
+	Href    string   `json:"href,omitempty"`
+	Text    string   `json:"text,omitempty"`
+	JSON    string   `json:"json,omitempty"`
+	Stderr  string   `json:"stderr,omitempty"`
+	Exit    int      `json:"exit,omitempty"`
+}
+
+// handleSurfaceVerb is GET /surface/verb.json?q=<line>, behind approvalRoute: the bar's line through
+// host/verbs' one grammar. It reads (episode report, the /v1 route's own builder) and writes nothing:
+// no store write, no credentials store, no process. An operator's typo is a 200 refusal.
+func (d *Daemon) handleSurfaceVerb(w http.ResponseWriter, r *http.Request, _ []string) {
+	query := r.URL.Query()
+	if !query.Has("q") {
+		writeApprovalRefusal(w, http.StatusBadRequest, "BadRequest", "q is required")
+		return
+	}
+	a := verbAnswer{V: 1, Schema: VerbSchema, Q: query.Get("q"), Argv: []string{}, Class: string(verbs.Refused), Outcome: "refused"}
+	defer func() { writeApprovalJSON(w, http.StatusOK, a) }()
+	argv, err := verbs.Split(a.Q) // its first check is the 1,024-byte bound: no oversized line is tokenized
+	if err != nil {
+		a.Reason = err.Error()
+		return
+	}
+	v, args := verbs.Lookup(argv)
+	a.Argv, a.Verb, a.Class, a.CLI, a.Reason = argv, v.Words, string(v.Class), v.CLI, v.Reason
+	switch {
+	case v.Class == verbs.Address:
+		a.Outcome = "address"
+	case v.Words == "episode report":
+		d.verbEpisodeReport(r, args, &a)
+	case v.Class == verbs.Navigate:
+		if a.Href, a.Reason = verbs.ApproveTarget(args); a.Href != "" {
+			a.Outcome = "navigate"
+		}
+	case v.Class == verbs.Read || v.Class == verbs.Mutate || v.Class == verbs.Host:
+		a.Outcome = "terminal"
+	}
+}
+
+// verbEpisodeReport runs episode report for the bar: the CLI's flag set and parse, the bar's
+// admission, then the /v1 route's builder; the text is the CLI's own printer's.
+func (d *Daemon) verbEpisodeReport(r *http.Request, args []string, a *verbAnswer) {
+	episode, asJSON, reason := verbs.AdmitEpisodeReport(args)
+	if reason != "" {
+		a.Reason = reason
+		return
+	}
+	if len(episode) > maxReportEpisodeBytes {
+		a.Reason = fmt.Sprintf("episode must be 1..%d bytes", maxReportEpisodeBytes)
+		return
+	}
+	ctx, cancel := d.readCtx(r)
+	defer cancel()
+	rep, found, err := d.episodeReport(ctx, episode)
+	if err != nil {
+		d.writeInternalErrorLog(r, err)
+		a.Reason = internalErrorMessage
+		if reportTimedOut(ctx, err) {
+			a.Reason = fmt.Sprintf("read deadline (%s) exceeded", d.readDeadline)
+		}
+		return
+	}
+	a.Outcome = "result"
+	switch {
+	case !found:
+		a.Exit, a.Stderr = 1, "ailang-worldd episode report: episode "+episode+" has no committed calls\n"
+	case asJSON:
+		var b bytes.Buffer
+		_ = json.NewEncoder(&b).Encode(rep) // the /v1 route's bytes: compact JSON and one newline, as the CLI prints them
+		a.JSON = b.String()
+	default:
+		var b strings.Builder
+		WriteEpisodeReportText(&b, rep)
+		a.Text = b.String()
+	}
+}
diff --git a/host/daemon/surface_verb_test.go b/host/daemon/surface_verb_test.go
new file mode 100644
index 0000000..40c2ad7
--- /dev/null
+++ b/host/daemon/surface_verb_test.go
@@ -0,0 +1,187 @@
+package daemon
+
+// Row 243 V1 (design §2.2, §6 VR1, VR4, VS1): GET /surface/verb.json.
+
+import (
+	"bytes"
+	"encoding/json"
+	"go/ast"
+	"go/parser"
+	"go/token"
+	"net/url"
+	"os"
+	"path/filepath"
+	"reflect"
+	"regexp"
+	"sort"
+	"strconv"
+	"strings"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/store"
+	"github.com/sunholo-data/ailang-world/host/verbs"
+)
+
+// verbRoutes is every verb route Handler registers.
+var verbRoutes = []string{"GET /surface/verb.json"}
+
+func verbRig(t *testing.T) (approvalRig, *reportSeed) {
+	t.Helper()
+	s := newReportSeed(t)
+	s.call("mcp", "w1", "1", "ailang-read", read(50, map[string]any{"ok": true, "content": "x"}), nil)
+	s.call("mcp", "w1", "2", `<img src=x onerror=alert(1)>`, nil, map[string]any{"ok": false, "refused": `<b>hostile</b> "reason"`})
+	return approvalRig{t: t, d: s.d}, s
+}
+
+func verbGet(a approvalRig, q string) (*verbAnswerT, string) {
+	a.t.Helper()
+	rec := a.serve("GET", "/surface/verb.json?q="+url.QueryEscape(q), apLoopback, nil)
+	if err := approvalHeaders(rec); rec.Code != 200 || err != nil {
+		a.t.Fatalf("%q: %d %s %v", q, rec.Code, rec.Body, err)
+	}
+	var v verbAnswerT
+	if err := json.Unmarshal(rec.Body.Bytes(), &v); err != nil || v.V != 1 || v.Schema != VerbSchema || v.Q != q || v.Argv == nil {
+		a.t.Fatalf("%q: not a v1 answer: %s", q, rec.Body)
+	}
+	return &v, rec.Body.String()
+}
+
+// verbAnswerT decodes an answer with every member optional, so a missing one is seen.
+type verbAnswerT struct {
+	V, Exit                                            int
+	Schema, Q, Verb, Class, Outcome, Reason, CLI, Href string
+	Argv                                               []string
+	Text, JSON, Stderr                                 *string
+}
+
+// VR1: one route, behind approvalRoute: a foreign Host is 421; GET only; nothing else under it.
+func TestSurfaceVerbRouteRegistered(t *testing.T) {
+	src, err := os.ReadFile("daemon.go")
+	if err != nil {
+		t.Fatal(err)
+	}
+	body := strings.SplitN(strings.SplitN(string(src), "func (d *Daemon) Handler() http.Handler {", 2)[1], "\n}\n", 2)[0]
+	var got []string
+	for _, m := range regexp.MustCompile(`mux\.HandleFunc\("([A-Z]+ /surface/verb[^"]*)", (d\.approvalRoute\()?`).FindAllStringSubmatch(body, -1) {
+		if m[2] == "" {
+			t.Errorf("%s is not behind approvalRoute", m[1])
+		}
+		got = append(got, m[1])
+	}
+	sort.Strings(got)
+	if !reflect.DeepEqual(got, verbRoutes) {
+		t.Fatalf("verb routes %v, want %v", got, verbRoutes)
+	}
+	a := newApprovalRig(t)
+	for _, host := range []string{"rebind.example:7644", "localhost:7645", "studio.tail0.ts.net"} {
+		requireRefusal(t, host, a.serve("GET", "/surface/verb.json?q=why", host, nil), 421, "MisdirectedRequest")
+	}
+	for _, method := range []string{"POST", "PUT", "DELETE", "OPTIONS"} {
+		if rec := a.serve(method, "/surface/verb.json?q=why", apLoopback, nil, "Origin", "http://localhost:7644"); rec.Code != 405 || rec.Header().Get("Access-Control-Allow-Origin") != "" {
+			t.Errorf("%s = %d, want 405 with no CORS header", method, rec.Code)
+		}
+	}
+	for _, target := range []string{"/surface/verb", "/surface/verb.json/x", "/surface/verb.jsonx"} {
+		if rec := a.serve("GET", target, apLoopback, nil); rec.Code != 404 {
+			t.Errorf("GET %s = %d, want 404", target, rec.Code)
+		}
+	}
+}
+
+// VR4: q is required; the line bound is an operator refusal, checked on the raw value; a store
+// failure or an expired read deadline answers a refusal with the /v1 route's text, never the error.
+func TestSurfaceVerbBoundsAndErrors(t *testing.T) {
+	a, s := verbRig(t)
+	requireRefusal(t, "no q", a.serve("GET", "/surface/verb.json", apLoopback, nil), 400, "BadRequest")
+	requireRefusal(t, "other key", a.serve("GET", "/surface/verb.json?line=why", apLoopback, nil), 400, "BadRequest")
+	if v, _ := verbGet(a, ""); v.Outcome != "address" {
+		t.Errorf("empty q = %+v", *v)
+	}
+	over := "why " + strings.Repeat("x", verbs.MaxLine-3)
+	if v, _ := verbGet(a, over); v.Outcome != "refused" || v.Reason != verbs.ErrTooLong.Error() || len(v.Argv) != 0 || v.Class != "refused" {
+		t.Errorf("1,025 bytes = %+v", *v)
+	}
+	if v, _ := verbGet(a, over[:verbs.MaxLine]); v.Outcome != "terminal" || len(v.Argv) != 2 {
+		t.Errorf("1,024 bytes = %+v", *v)
+	}
+	if v, _ := verbGet(a, `why "a`); v.Outcome != "refused" || v.Reason != verbs.ErrUnclosed.Error() {
+		t.Errorf("unclosed = %+v", *v)
+	}
+	var log bytes.Buffer
+	s.d.errLog = &log
+	s.d.readDeadline = expiredReadDeadline
+	if v, raw := verbGet(a, "episode report w1"); v.Outcome != "refused" || v.Reason != "read deadline (-1ns) exceeded" || v.Text != nil {
+		t.Errorf("expired deadline = %s", raw)
+	}
+	s.d.readDeadline = readDeadline
+	if err := s.d.store.Close(); err != nil {
+		t.Fatal(err)
+	}
+	v, raw := verbGet(a, "episode report w1")
+	if v.Outcome != "refused" || v.Reason != internalErrorMessage || strings.Contains(raw, "closed") {
+		t.Errorf("store failure = %s", raw)
+	}
+	if !strings.Contains(log.String(), "internal error: GET /surface/verb.json") {
+		t.Errorf("the failure is not logged for the operator: %q", log.String())
+	}
+	// The /v1 route still answers its own classes from the extracted builder.
+	assertErrorClass(t, requestRecorder(t, s.d, "GET", "/v1/episodes/w1/report", nil), 500, "Internal")
+}
+
+// VS1: host/verbs imports no I/O package; the route's file reaches no store, credentials store or
+// process; approve's id pattern is host/store's.
+func TestSurfaceVerbIsReadOnly(t *testing.T) {
+	importsOf := func(path string) (map[string]bool, *ast.File) {
+		t.Helper()
+		f, err := parser.ParseFile(token.NewFileSet(), path, nil, 0)
+		if err != nil {
+			t.Fatal(err)
+		}
+		out := map[string]bool{}
+		for _, im := range f.Imports {
+			p, _ := strconv.Unquote(im.Path.Value)
+			out[p] = true
+		}
+		return out, f
+	}
+	files, err := filepath.Glob(filepath.Join("..", "verbs", "*.go"))
+	if err != nil || len(files) == 0 {
+		t.Fatalf("host/verbs: %v %v", files, err)
+	}
+	pure := map[string]bool{"errors": true, "flag": true, "fmt": true, "io": true, "regexp": true, "strings": true, "time": true, "unicode/utf8": true}
+	for _, f := range files {
+		if strings.HasSuffix(f, "_test.go") {
+			continue
+		}
+		ims, _ := importsOf(f)
+		for p := range ims {
+			if !pure[p] {
+				t.Errorf("%s imports %s", f, p)
+			}
+		}
+	}
+	ims, f := importsOf("surface_verb.go")
+	for p := range ims {
+		if !map[string]bool{"bytes": true, "encoding/json": true, "fmt": true, "net/http": true, "strings": true,
+			"github.com/sunholo-data/ailang-world/host/verbs": true}[p] {
+			t.Errorf("surface_verb.go imports %s", p)
+		}
+	}
+	calls := map[string]bool{}
+	ast.Inspect(f, func(n ast.Node) bool {
+		if sel, ok := n.(*ast.SelectorExpr); ok {
+			if id, ok := sel.X.(*ast.Ident); ok && id.Name == "d" {
+				calls[sel.Sel.Name] = true
+			}
+		}
+		return true
+	})
+	if want := map[string]bool{"readCtx": true, "episodeReport": true, "writeInternalErrorLog": true, "verbEpisodeReport": true, "readDeadline": true}; !reflect.DeepEqual(calls, want) {
+		t.Errorf("surface_verb.go reaches d.%v, want exactly %v", calls, want)
+	}
+	for _, id := range []string{"mr-0123456789abcdef", "mr-0123456789ABCDEF", "mr-0123456789abcde", "mr-0123456789abcdef0", "xr-0123456789abcdef", "mr-012345678 abcdef"} {
+		if href, _ := verbs.ApproveTarget([]string{id}); (href != "") != store.MintRequestIDPattern.MatchString(id) {
+			t.Errorf("approve %q: href %q, but store's pattern says %v", id, href, store.MintRequestIDPattern.MatchString(id))
+		}
+	}
+}
````

### Patch M3b — every outcome, the page's golden (VR2); the CLI's bytes (VR3)

````diff M3b
diff --git a/cmd/ailang-worldd/verb_report_e2e_test.go b/cmd/ailang-worldd/verb_report_e2e_test.go
new file mode 100644
index 0000000..4366eef
--- /dev/null
+++ b/cmd/ailang-worldd/verb_report_e2e_test.go
@@ -0,0 +1,127 @@
+package main
+
+// Row 243 V1 (design §6 VR3): the bar's episode report is the CLI's, byte for byte — the real CLI
+// and the real route against one loopback daemon on one test store, two episodes and the not-found
+// case, every fixture generated by the run itself.
+
+import (
+	"encoding/json"
+	"fmt"
+	"io"
+	"net/http"
+	"net/http/httptest"
+	"net/url"
+	"path/filepath"
+	"strings"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/broker"
+	"github.com/sunholo-data/ailang-world/host/coordinator"
+	"github.com/sunholo-data/ailang-world/host/daemon"
+	"github.com/sunholo-data/ailang-world/host/hashref"
+	"github.com/sunholo-data/ailang-world/host/store"
+)
+
+// seedReportCalls journals and commits coordinator-shaped calls on a fresh store, as a served
+// episode leaves them (host/daemon's reportSeed layout): each call is (episode, task, tool, effect
+// or "" for a plan refusal, handler ok).
+func seedReportCalls(t *testing.T, db string, calls [][5]string) {
+	t.Helper()
+	st, err := store.Open(db)
+	if err != nil {
+		t.Fatal(err)
+	}
+	defer func() { _ = st.Close() }()
+	ctx := boundedTestContext(t)
+	obj := func(semantic string, v any) store.Object {
+		b := mustJSONT(t, v)
+		return store.Object{Hash: hashref.SumSHA256(b), InterfaceHash: hashref.SumSHA256([]byte(semantic)), SemanticID: semantic, Provenance: "verb-report-test", Payload: b}
+	}
+	head := store.World{Ref: hashref.SumSHA256([]byte("g")), StateRoot: hashref.SumSHA256([]byte("gs")), LogHead: hashref.SumSHA256([]byte("gl"))}
+	if err := st.PutWorld(ctx, head); err != nil {
+		t.Fatal(err)
+	}
+	if err := st.SelectHead(ctx, head.Ref); err != nil {
+		t.Fatal(err)
+	}
+	for i, c := range calls {
+		id, index := "mcp:"+c[0]+":"+c[1], int64(i+1)
+		var objects []store.Object
+		refs, planned, result := []string{}, []map[string]any{}, any(nil)
+		if c[3] == "" {
+			result = map[string]any{"ok": false, "refused": "missing <path>"}
+		} else {
+			req := obj(broker.EffectRequestV1, map[string]any{"req": id})
+			res := obj(broker.EffectResultV1, map[string]any{"ok": c[4] == "ok", "refused": c[4]})
+			rec := obj(broker.EffectRecordV1, nil)
+			rec.Payload = broker.EncodeRecord(broker.EffectRecord{Effect: c[3], Scope: "worktree", Cost: 1, BudgetBefore: 9, BudgetAfter: 8, Allowed: true, RequestRef: req.Hash, ResultRef: res.Hash})
+			rec.Hash = hashref.SumSHA256(rec.Payload)
+			objects, refs = append(objects, req, res, rec), []string{rec.Hash.String()}
+			planned = []map[string]any{{"id": "e1", "effect": c[3], "scope": "worktree", "cost": 1}}
+		}
+		in, out := obj(coordinator.InputV1, map[string]any{"task": id}), obj(coordinator.OutputV1, map[string]any{"out": id})
+		plan := obj(coordinator.EffectPlanV1, map[string]any{"plan": "world/effect-plan/v1", "effects": planned, "finish": false, "result": result})
+		fn, interp := hashref.SumSHA256([]byte("fn-"+c[2])), hashref.SumSHA256([]byte("interpreter"))
+		record := obj(coordinator.RecordV2, map[string]any{"invocationId": id, "episodeId": c[0], "skillId": c[2], "transitionFn": fn.String(),
+			"interpreter": interp.String(), "semanticsEpoch": 1, "input": in.Hash.String(), "output": out.Hash.String(), "plan": plan.Hash.String(), "effects": refs})
+		entry := hashref.SumSHA256([]byte(fmt.Sprintf("entry-%d", index)))
+		next := store.World{Ref: hashref.SumSHA256([]byte(fmt.Sprintf("world-%d", index))), Revision: index, StateRoot: out.Hash, LogHead: entry}
+		if _, _, err := st.AppendIntent(ctx, id, store.JournalIntent{InvocationID: id, WorldRef: next.Ref, EntryHash: entry, ObservedHead: head.Ref,
+			PrevEntryHash: head.LogHead, TransitionFn: fn, TransitionRef: record.Hash, Interpreter: interp, LogicalTime: 1_800_000_000 + index}); err != nil {
+			t.Fatal(err)
+		}
+		if err := st.Commit(ctx, store.Commit{InvocationID: id, ObservedHead: head.Ref, Objects: append(objects, in, out, plan, record), NextWorld: next,
+			Entry: store.LogEntry{Header: store.LogHeader{EntryIndex: index, SemanticsEpoch: 1, TransitionFn: fn, Interpreter: interp,
+				PrevEntryHash: head.LogHead, WrittenBy: "coordinator:mcp"}, EntryHash: entry, TransitionRef: record.Hash}}); err != nil {
+			t.Fatal(err)
+		}
+		head = next
+	}
+}
+
+func TestVerbRouteReportIsTheCLIs(t *testing.T) {
+	db := filepath.Join(t.TempDir(), "world.db")
+	seedReportCalls(t, db, [][5]string{{"w1", "1", "ailang-read", broker.EffectWorkspaceRead, "ok"}, {"w2", "1", "workspace-exec", broker.EffectWorkspaceExec, "the handler <b>refused</b>"},
+		{"w1", "2", "ailang-read", "", ""}, {"w2", "2", "ailang-read", broker.EffectWorkspaceRead, "ok"}})
+	d, err := daemon.New(boundedTestContext(t), daemon.Config{DBPath: db, BindHost: daemon.DefaultBindHost, ErrorLog: io.Discard})
+	if err != nil {
+		t.Fatal(err)
+	}
+	srv := httptest.NewServer(d.Handler())
+	t.Cleanup(func() { srv.Close(); _ = d.Close() })
+	route := func(q string) map[string]any {
+		t.Helper()
+		resp, err := http.Get(srv.URL + "/surface/verb.json?q=" + url.QueryEscape(q))
+		if err != nil {
+			t.Fatal(err)
+		}
+		defer resp.Body.Close()
+		body, _ := io.ReadAll(resp.Body)
+		var a map[string]any
+		if resp.StatusCode != 200 || json.Unmarshal(body, &a) != nil || a["outcome"] != "result" || a["cli"] != "ailang-worldd "+q {
+			t.Fatalf("%q: %d %s", q, resp.StatusCode, body)
+		}
+		return a
+	}
+	for _, c := range []struct{ episode, line string }{{"w1", "plan refusals (refused before any effect): 1 (ailang-read 1)"}, {"w2", "handler refusals (effect allowed, handler ok:false): 1 (workspace-exec 1)"}} {
+		code, stdout, stderr := runCLI(t, srv.URL, "episode", "report", c.episode)
+		if code != exitOK || !strings.Contains(stdout, c.line) {
+			t.Fatalf("CLI %s: exit %d %q %s", c.episode, code, stdout, stderr)
+		}
+		if a := route("episode report " + c.episode); a["text"] != stdout || a["json"] != nil || a["stderr"] != nil {
+			t.Errorf("%s: the route's text differs from the CLI's:\nroute %q\n  cli %q", c.episode, a["text"], stdout)
+		}
+		code, stdout, _ = runCLI(t, srv.URL, "episode", "report", c.episode, "--json")
+		if code != exitOK || !strings.HasPrefix(stdout, `{"episode":"`+c.episode+`"`) {
+			t.Fatalf("CLI %s --json: exit %d %q", c.episode, code, stdout)
+		}
+		if a := route("episode report --json " + c.episode); a["json"] != stdout || a["text"] != nil {
+			t.Errorf("%s --json: the route's json differs from the CLI's:\nroute %q\n  cli %q", c.episode, a["json"], stdout)
+		}
+	}
+	code, stdout, stderr := runCLI(t, srv.URL, "episode", "report", "nobody")
+	if a := route("episode report nobody"); code != exitUsage || stdout != "" || !strings.Contains(stderr, "episode nobody has no committed calls") ||
+		a["stderr"] != stderr || a["exit"] != float64(code) || a["text"] != nil {
+		t.Errorf("not found: CLI exit %d %q; route %v", code, stderr, a)
+	}
+}
diff --git a/host/daemon/surface_verb_outcomes_test.go b/host/daemon/surface_verb_outcomes_test.go
new file mode 100644
index 0000000..3fd0fb8
--- /dev/null
+++ b/host/daemon/surface_verb_outcomes_test.go
@@ -0,0 +1,96 @@
+package daemon
+
+// Row 243 V1 (design §2.4, §6 VR2): every outcome class end to end, and the page's replayed answers.
+
+import (
+	"bytes"
+	"encoding/json"
+	"os"
+	"path/filepath"
+	"strings"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/verbs"
+)
+
+// verbGolden is the route's real answers the page's node test replays (VP6): regenerate with
+// WORLD_UPDATE_VERB_GOLDEN=1; otherwise every answer must equal it byte for byte.
+var verbGolden = filepath.Join("..", "surfacepage", "testdata", "verb-route.golden.json")
+
+// VR2: each outcome end to end through the handler, and the page's replayed answers are these.
+func TestSurfaceVerbOutcomes(t *testing.T) {
+	a, _ := verbRig(t)
+	const mr, ref = "mr-0123456789abcdef", "sha256:abababababababababababababababababababababababababababababababab"
+	for _, c := range []struct{ q, outcome, class, verb, cli, href, reason string }{
+		{"row-189@6", "address", "address", "", "", "", ""},
+		{"world shared@0ffbb9c", "address", "address", "", "", "", ""},
+		{"back", "address", "address", "", "", "", ""},
+		{"world why", "address", "address", "", "", "", ""},
+		{"world session", "address", "address", "", "", "", ""},
+		{"world get", "address", "address", "", "", "", ""},
+		{"world episode", "address", "address", "", "", "", ""},
+		{"why", "terminal", "read", "why", "ailang-worldd why", "", "why reads World"},
+		{"episode finish w1 --preset p", "terminal", "mutate", "episode finish", "ailang-worldd episode finish w1 --preset p", "", "changes World"},
+		{"episode start w1", "terminal", "mutate", "episode start", "ailang-worldd episode start w1", "", "changes World"},
+		{"episode launch w1 --preset p", "terminal", "host", "episode launch", "ailang-worldd episode launch w1 --preset p", "", "needs a host process"},
+		{"approve " + mr, "navigate", "navigate", "approve", "ailang-worldd approve " + mr, "/surface/approve/" + mr, ""},
+		{"ailang-worldd approve " + ref, "navigate", "navigate", "approve", "ailang-worldd approve " + ref, "/surface/decide/" + ref, ""},
+		{"approve " + mr + " --deny", "refused", "navigate", "approve", "ailang-worldd approve " + mr + " --deny", "", "the app opens approve"},
+		{"approver enroll --name x", "refused", "never", "approver", "ailang-worldd approver enroll --name x", "", "never run from the app"},
+		{"session frob", "refused", "never", "session", "ailang-worldd session frob", "", "never run from the app"},
+		{"log", "refused", "other", "log", "ailang-worldd log", "", "log is a CLI command the app does not run; at a terminal: ailang-worldd log"},
+		{"log tail", "refused", "other", "log", "ailang-worldd log tail", "", "CLI command"},
+		{"ailang-worldd world get world-1", "refused", "other", "world get", "ailang-worldd world get world-1", "", "world get is a CLI command the app does not run; at a terminal: ailang-worldd world get world-1"},
+		{"episode", "refused", "refused", "episode", "ailang-worldd episode", "", verbs.EpisodeNeedsSubcommand},
+		{"episode frobnicate", "refused", "refused", "episode", "ailang-worldd episode frobnicate", "", verbs.EpisodeNeedsSubcommand},
+		{"episode w1", "refused", "refused", "episode", "ailang-worldd episode w1", "", verbs.EpisodeNeedsSubcommand},
+		{"world a b c", "refused", "refused", "world", "", "", verbs.AddressHelp},
+		{"--addr x episode report w1", "refused", "refused", "", "ailang-worldd --addr x episode report w1", "", "global flags are not accepted"},
+		{"episode report w1 --follow", "refused", "read", "episode report", "ailang-worldd episode report w1 --follow", "", "--follow is not run from the app"},
+		{"episode report w1 --markdown", "refused", "read", "episode report", "ailang-worldd episode report w1 --markdown", "", "--markdown is not run from the app"},
+		{"episode report " + strings.Repeat("e", 257), "refused", "read", "episode report", "ailang-worldd episode report " + strings.Repeat("e", 257), "", "episode must be 1..256 bytes"},
+		{"episode report w1; rm -rf ~", "refused", "refused", "", "", "", "Shell:"},
+		{"episode report nobody", "result", "read", "episode report", "ailang-worldd episode report nobody", "", ""},
+		{"episode report w1", "result", "read", "episode report", "ailang-worldd episode report w1", "", ""},
+		{"ailang-worldd episode report w1 --json", "result", "read", "episode report", "ailang-worldd episode report w1 --json", "", ""},
+	} {
+		v, _ := verbGet(a, c.q)
+		if v.Outcome != c.outcome || v.Class != c.class || v.Verb != c.verb || v.CLI != c.cli || v.Href != c.href || !strings.Contains(v.Reason, c.reason) || (c.reason == "" && v.Reason != "") {
+			t.Errorf("%q = %+v, want outcome %s class %s verb %q cli %q href %q reason %q", c.q, *v, c.outcome, c.class, c.verb, c.cli, c.href, c.reason)
+		}
+		if (v.Text != nil || v.JSON != nil || v.Stderr != nil) != (c.outcome == "result") {
+			t.Errorf("%q: result members on a %s answer", c.q, v.Outcome)
+		}
+	}
+	if v, _ := verbGet(a, "episode report nobody"); v.Stderr == nil || *v.Stderr != "ailang-worldd episode report: episode nobody has no committed calls\n" || v.Exit != 1 || v.Text != nil {
+		t.Errorf("not found = %+v", *v)
+	}
+	if v, _ := verbGet(a, "episode report w1"); v.Text == nil || !strings.Contains(*v.Text, "episode w1: 2 calls") || !strings.Contains(*v.Text, `#2 <img src=x onerror=alert(1)>: <b>hostile</b> "reason"`) || v.JSON != nil || v.Exit != 0 {
+		t.Errorf("text = %+v", *v)
+	}
+	if v, _ := verbGet(a, "episode report w1 --json"); v.JSON == nil || !strings.HasPrefix(*v.JSON, `{"episode":"w1",`) || !strings.HasSuffix(*v.JSON, "}\n") || strings.Count(*v.JSON, "\n") != 1 {
+		t.Errorf("json = %+v", *v)
+	}
+	// The page's node test (VP6, VP1–VP3) replays these real answers, never hand-typed ones.
+	golden := map[string]string{}
+	for _, q := range []string{"row-189@6", "world shared@0ffbb9c", "back", "world why", "why", "log", "world session",
+		"ailang-worldd world get world-1", "episode", "world episode", "world a b c", "episode finish w1 --preset p",
+		"approve mr-0123456789abcdef", "approve " + ref, "episode report nobody", "episode report w1", "episode report w1 --json",
+		strings.Repeat("x", verbs.MaxLine+1)} {
+		_, raw := verbGet(a, q)
+		golden[q] = raw
+	}
+	want, err := json.MarshalIndent(golden, "", "  ")
+	if err != nil {
+		t.Fatal(err)
+	}
+	want = append(want, '\n')
+	if os.Getenv("WORLD_UPDATE_VERB_GOLDEN") == "1" {
+		if err := os.WriteFile(verbGolden, want, 0o644); err != nil {
+			t.Fatal(err)
+		}
+	}
+	if got, err := os.ReadFile(verbGolden); err != nil || !bytes.Equal(got, want) {
+		t.Errorf("%s is not the route's answers (WORLD_UPDATE_VERB_GOLDEN=1 regenerates it): %v", verbGolden, err)
+	}
+}
diff --git a/host/surfacepage/testdata/verb-route.golden.json b/host/surfacepage/testdata/verb-route.golden.json
new file mode 100644
index 0000000..d4c2c92
--- /dev/null
+++ b/host/surfacepage/testdata/verb-route.golden.json
@@ -0,0 +1,20 @@
+{
+  "ailang-worldd world get world-1": "{\"v\":1,\"schema\":\"world/surface-verb/v1\",\"q\":\"ailang-worldd world get world-1\",\"argv\":[\"ailang-worldd\",\"world\",\"get\",\"world-1\"],\"verb\":\"world get\",\"class\":\"other\",\"outcome\":\"refused\",\"reason\":\"world get is a CLI command the app does not run; at a terminal: ailang-worldd world get world-1\",\"cli\":\"ailang-worldd world get world-1\"}\n",
+  "approve mr-0123456789abcdef": "{\"v\":1,\"schema\":\"world/surface-verb/v1\",\"q\":\"approve mr-0123456789abcdef\",\"argv\":[\"approve\",\"mr-0123456789abcdef\"],\"verb\":\"approve\",\"class\":\"navigate\",\"outcome\":\"navigate\",\"cli\":\"ailang-worldd approve mr-0123456789abcdef\",\"href\":\"/surface/approve/mr-0123456789abcdef\"}\n",
+  "approve sha256:abababababababababababababababababababababababababababababababab": "{\"v\":1,\"schema\":\"world/surface-verb/v1\",\"q\":\"approve sha256:abababababababababababababababababababababababababababababababab\",\"argv\":[\"approve\",\"sha256:abababababababababababababababababababababababababababababababab\"],\"verb\":\"approve\",\"class\":\"navigate\",\"outcome\":\"navigate\",\"cli\":\"ailang-worldd approve sha256:abababababababababababababababababababababababababababababababab\",\"href\":\"/surface/decide/sha256:abababababababababababababababababababababababababababababababab\"}\n",
+  "back": "{\"v\":1,\"schema\":\"world/surface-verb/v1\",\"q\":\"back\",\"argv\":[\"back\"],\"verb\":\"\",\"class\":\"address\",\"outcome\":\"address\"}\n",
+  "episode": "{\"v\":1,\"schema\":\"world/surface-verb/v1\",\"q\":\"episode\",\"argv\":[\"episode\"],\"verb\":\"episode\",\"class\":\"refused\",\"outcome\":\"refused\",\"reason\":\"ailang-worldd episode needs a subcommand: start | launch | finish | report\",\"cli\":\"ailang-worldd episode\"}\n",
+  "episode finish w1 --preset p": "{\"v\":1,\"schema\":\"world/surface-verb/v1\",\"q\":\"episode finish w1 --preset p\",\"argv\":[\"episode\",\"finish\",\"w1\",\"--preset\",\"p\"],\"verb\":\"episode finish\",\"class\":\"mutate\",\"outcome\":\"terminal\",\"reason\":\"episode finish changes World. The app cannot run it yet.\",\"cli\":\"ailang-worldd episode finish w1 --preset p\"}\n",
+  "episode report nobody": "{\"v\":1,\"schema\":\"world/surface-verb/v1\",\"q\":\"episode report nobody\",\"argv\":[\"episode\",\"report\",\"nobody\"],\"verb\":\"episode report\",\"class\":\"read\",\"outcome\":\"result\",\"cli\":\"ailang-worldd episode report nobody\",\"stderr\":\"ailang-worldd episode report: episode nobody has no committed calls\\n\",\"exit\":1}\n",
+  "episode report w1": "{\"v\":1,\"schema\":\"world/surface-verb/v1\",\"q\":\"episode report w1\",\"argv\":[\"episode\",\"report\",\"w1\"],\"verb\":\"episode report\",\"class\":\"read\",\"outcome\":\"result\",\"cli\":\"ailang-worldd episode report w1\",\"text\":\"episode w1: 2 calls, entries 1-2, 2027-01-15T08:00:01Z .. 2027-01-15T08:00:02Z\\ncalls by tool:\\n  \\u003cimg src=x onerror=alert(1)\\u003e      1\\n  ailang-read             1\\nreads: 0 calls, 0 files, 1 outside the worktree\\neffects:\\n  Workspace.Read         allowed               1\\nbudgets (minted / used / left):\\n  Workspace.Read@worktree     50      1     49\\nplan refusals (refused before any effect): 1 (\\u003cimg src=x onerror=alert(1)\\u003e 1)\\n  #2 \\u003cimg src=x onerror=alert(1)\\u003e: \\u003cb\\u003ehostile\\u003c/b\\u003e \\\"reason\\\"\\nhandler refusals (effect allowed, handler ok:false): 0\\n\"}\n",
+  "episode report w1 --json": "{\"v\":1,\"schema\":\"world/surface-verb/v1\",\"q\":\"episode report w1 --json\",\"argv\":[\"episode\",\"report\",\"w1\",\"--json\"],\"verb\":\"episode report\",\"class\":\"read\",\"outcome\":\"result\",\"cli\":\"ailang-worldd episode report w1 --json\",\"json\":\"{\\\"episode\\\":\\\"w1\\\",\\\"entries\\\":[1,2],\\\"total_calls\\\":2,\\\"calls_by_tool\\\":{\\\"\\\\u003cimg src=x onerror=alert(1)\\\\u003e\\\":1,\\\"ailang-read\\\":1},\\\"effects\\\":[{\\\"effect\\\":\\\"Workspace.Read\\\",\\\"decision\\\":\\\"allowed\\\",\\\"count\\\":1}],\\\"budgets\\\":[{\\\"effect\\\":\\\"Workspace.Read\\\",\\\"scope\\\":\\\"worktree\\\",\\\"minted\\\":50,\\\"used\\\":1,\\\"left\\\":49}],\\\"plan_refusals\\\":{\\\"count\\\":1,\\\"by_tool\\\":{\\\"\\\\u003cimg src=x onerror=alert(1)\\\\u003e\\\":1},\\\"first\\\":[{\\\"entry\\\":2,\\\"tool\\\":\\\"\\\\u003cimg src=x onerror=alert(1)\\\\u003e\\\",\\\"reason\\\":\\\"\\\\u003cb\\\\u003ehostile\\\\u003c/b\\\\u003e \\\\\\\"reason\\\\\\\"\\\"}]},\\\"handler_refusals\\\":{\\\"count\\\":0,\\\"by_tool\\\":{},\\\"first\\\":[]},\\\"writes\\\":{\\\"written\\\":[],\\\"produced\\\":[],\\\"deltas_truncated\\\":[],\\\"delta_errors\\\":[],\\\"written_truncated\\\":0,\\\"produced_truncated\\\":0,\\\"unlisted\\\":0,\\\"deltas_truncated_omitted\\\":0,\\\"delta_errors_omitted\\\":0},\\\"reads\\\":{\\\"count\\\":0,\\\"files\\\":[],\\\"truncated\\\":0,\\\"unlisted\\\":1},\\\"first_at\\\":\\\"2027-01-15T08:00:01Z\\\",\\\"last_at\\\":\\\"2027-01-15T08:00:02Z\\\",\\\"calls_by_surface\\\":{\\\"mcp\\\":2},\\\"uncommitted_intents\\\":0,\\\"notes\\\":{\\\"budgets\\\":\\\"per (effect, scope) with at least one effect record: left = latest budgetAfter, used = cost of allowed records, minted = left + used\\\",\\\"entries\\\":\\\"first and last log entries of the episode's committed calls\\\",\\\"first_at\\\":\\\"coordinator logical time (unix seconds) of the first and last committed call, UTC\\\",\\\"handler_refusals\\\":\\\"calls whose effect was allowed and ran, but whose handler result is ok:false (the effect record alone shows these as allowed)\\\",\\\"plan_refusals\\\":\\\"calls whose plan named no effect and answered ok:false (refused before any effect ran)\\\",\\\"reads\\\":\\\"committed ailang-read calls with a v2 record whose every effect was allowed and did not fail and whose handler result is not ok:false; the path is the record input's \\\\\\\"path\\\\\\\" (a ranged read counts); a path that is not local and relative counts in unlisted and is never shown; files are ordered by first entry, at most 1000\\\",\\\"writes\\\":\\\"written: successful committed ailang-write/edit/delete calls, from input.path; produced: allowed non-failed effect result worktree_delta.changes; deltas_truncated and delta_errors name evidence gaps; local relative paths only, each path list ordered by first entry and capped at 1000; excess counts and unlisted are explicit\\\"}}\\n\"}\n",
+  "log": "{\"v\":1,\"schema\":\"world/surface-verb/v1\",\"q\":\"log\",\"argv\":[\"log\"],\"verb\":\"log\",\"class\":\"other\",\"outcome\":\"refused\",\"reason\":\"log is a CLI command the app does not run; at a terminal: ailang-worldd log\",\"cli\":\"ailang-worldd log\"}\n",
+  "row-189@6": "{\"v\":1,\"schema\":\"world/surface-verb/v1\",\"q\":\"row-189@6\",\"argv\":[\"row-189@6\"],\"verb\":\"\",\"class\":\"address\",\"outcome\":\"address\"}\n",
+  "why": "{\"v\":1,\"schema\":\"world/surface-verb/v1\",\"q\":\"why\",\"argv\":[\"why\"],\"verb\":\"why\",\"class\":\"read\",\"outcome\":\"terminal\",\"reason\":\"why reads World, but its walk runs in the CLI. The app cannot run it yet.\",\"cli\":\"ailang-worldd why\"}\n",
+  "world a b c": "{\"v\":1,\"schema\":\"world/surface-verb/v1\",\"q\":\"world a b c\",\"argv\":[\"world\",\"a\",\"b\",\"c\"],\"verb\":\"world\",\"class\":\"refused\",\"outcome\":\"refused\",\"reason\":\"unknown address: try \\\"world\\\", \\\"world NAME\\\", \\\"world NAME@N\\\", \\\"world NAME@HASH\\\" or \\\"back\\\"\"}\n",
+  "world episode": "{\"v\":1,\"schema\":\"world/surface-verb/v1\",\"q\":\"world episode\",\"argv\":[\"world\",\"episode\"],\"verb\":\"\",\"class\":\"address\",\"outcome\":\"address\"}\n",
+  "world session": "{\"v\":1,\"schema\":\"world/surface-verb/v1\",\"q\":\"world session\",\"argv\":[\"world\",\"session\"],\"verb\":\"\",\"class\":\"address\",\"outcome\":\"address\"}\n",
+  "world shared@0ffbb9c": "{\"v\":1,\"schema\":\"world/surface-verb/v1\",\"q\":\"world shared@0ffbb9c\",\"argv\":[\"world\",\"shared@0ffbb9c\"],\"verb\":\"\",\"class\":\"address\",\"outcome\":\"address\"}\n",
+  "world why": "{\"v\":1,\"schema\":\"world/surface-verb/v1\",\"q\":\"world why\",\"argv\":[\"world\",\"why\"],\"verb\":\"\",\"class\":\"address\",\"outcome\":\"address\"}\n",
+  "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx": "{\"v\":1,\"schema\":\"world/surface-verb/v1\",\"q\":\"xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx\",\"argv\":[],\"verb\":\"\",\"class\":\"refused\",\"outcome\":\"refused\",\"reason\":\"TooLong: a line is at most 1024 bytes\"}\n"
+}
````

### Patch M4a — the page (VP1–VP6)

````diff M4a
diff --git a/host/surfacepage/app/app.js b/host/surfacepage/app/app.js
index 39b4502..a32fad7 100644
--- a/host/surfacepage/app/app.js
+++ b/host/surfacepage/app/app.js
@@ -480,8 +480,72 @@
       + (d.length > DECIDED_FEW ? `<button type="button" class="nl run" id="decided-all" aria-expanded="${!!all}">${all ? 'Show fewer' : 'Show all ' + d.length}</button>` : '') + '</div>';
   }
 
+  // ---------- The bar's verbs (row 243 V1): GET /surface/verb.json decides; only a schema-valid address navigates ----------
+  const VERB_ROUTE_BUDGET_MS = 15000;   // above the daemon's 10 s read deadline; a named constant, not a setting
+  const VERB_NEEDS = { address: [], result: ['cli'], navigate: ['href'], terminal: ['cli', 'reason'], refused: ['reason'] };
+  const UNREADABLE = { say: 'the verb route answered something this page cannot read' };
+  // verbOutcome maps the route's answer (HTTP status, body text) to one action; anything else is a refusal.
+  function verbOutcome(status, text) {
+    if (status !== 200) return { say: 'the verb route did not answer: HTTP ' + status };
+    let b = null;
+    try { b = JSON.parse(text); } catch (_) { return UNREADABLE; }
+    const need = b && typeof b === 'object' && b.schema === 'world/surface-verb/v1' && b.v === 1 && Object.prototype.hasOwnProperty.call(VERB_NEEDS, b.outcome) ? VERB_NEEDS[b.outcome] : null;
+    if (!need || !need.every(k => typeof b[k] === 'string' && b[k] !== '')) return UNREADABLE;
+    if (b.outcome === 'address') return { address: true };
+    if (b.outcome === 'refused') return { say: b.reason };
+    if (b.outcome === 'navigate') return b.href.startsWith('/surface/approve/mr-') || b.href.startsWith('/surface/decide/sha256:') ? { assign: b.href } : { say: 'refused: unexpected target' };
+    if (b.outcome === 'terminal') return { panel: { head: '$ ' + b.cli, body: b.reason + '\nAt a terminal: ' + b.cli, copy: b.cli } };
+    const shown = [b.text, b.json, b.stderr].filter(x => typeof x === 'string' && x !== '');
+    if (shown.length !== 1 || (b.stderr !== undefined && b.exit !== 1)) return UNREADABLE;
+    return { panel: { head: '$ ' + b.cli, body: b.stderr !== undefined ? b.stderr + 'exit 1' : shown[0], copy: '' } };
+  }
+  // askVerb asks the route under one budget, across the fetch and the body (io carries fetch, the
+  // clock and AbortController). A newer ask supersedes an older one: the older is aborted and its
+  // answer dropped, never acted on; the newer answers for the bar.
+  let verbJob = null;
+  async function askVerb(raw, io) {
+    if (verbJob) { verbJob.superseded = true; verbJob.ctl.abort(); }
+    const job = { ctl: new io.AbortController(), superseded: false, late: false };
+    verbJob = job;
+    const stop = new Promise((_, reject) => job.ctl.signal.addEventListener('abort', () => reject(new Error('aborted'))));
+    stop.catch(() => {});
+    const timer = io.setTimeout(() => { job.late = true; job.ctl.abort(); }, VERB_ROUTE_BUDGET_MS);
+    try {
+      const res = await Promise.race([io.fetch('/surface/verb.json?q=' + encodeURIComponent(raw), { signal: job.ctl.signal, headers: { Accept: 'application/json' } }), stop]);
+      const text = await Promise.race([res.text(), stop]);
+      return job.superseded ? { superseded: true } : verbOutcome(res.status, text);   // an answer that beat the abort is still dropped
+    } catch (e) {
+      if (job.superseded) return { superseded: true };
+      if (job.late) return { say: 'the verb route did not answer in time' };
+      return { say: 'the verb route is unreachable: ' + (e && e.message ? e.message : String(e)) };
+    } finally {
+      io.clearTimeout(timer);
+      if (verbJob === job) verbJob = null;
+    }
+  }
+  // verbAct runs one action through fx: only an address reaches the address path; only an address,
+  // a navigation or a panel clears the bar; a refusal says why and leaves the view and the input.
+  function verbAct(act, raw, fx) {
+    if (act.superseded) return;
+    if (act.address) { fx.clear(); return fx.address(raw); }
+    if (act.assign) { fx.clear(); return fx.assign(act.assign); }
+    if (act.panel) { fx.clear(); return fx.panel(act.panel); }
+    fx.say(act.say);
+  }
+  // verbPanel builds the panel's nodes with textContent only: a report or a reason is never markup.
+  function verbPanel(doc, p) {
+    const el = doc.createElement('div'), head = doc.createElement('div'), body = doc.createElement('pre');
+    el.id = 'verb'; el.className = 'pop'; el.setAttribute('role', 'dialog'); el.setAttribute('aria-label', 'Verb');
+    head.className = 'mono'; head.textContent = p.head;
+    body.style.cssText = 'margin:10px 0 0;white-space:pre-wrap;max-height:60vh;overflow:auto;font-size:12px'; body.textContent = p.body;
+    el.append(head, body);
+    if (p.copy) { const b = doc.createElement('button'); b.type = 'button'; b.className = 'hbtn'; b.textContent = 'Copy the line'; b.dataset.copy = p.copy; el.append(b); }
+    return el;
+  }
+
   if (typeof module !== 'undefined') module.exports = { needsYouHTML, decidedHTML, notifications, isSharedKey, sourcesHTML, openRunningHTML, fetchCoherent, serviceLineHTML, splitFrames, applyDelta, nextBackoff, liveRequest, liveStatus, esc, cleanTitle, shortTitle, prettyModel, agentLine, summaryOf, rowLabel, verdict, LEVELS, levelXs, blendXS,
-    adaptModel, layoutGraph, stackBars, barHTML, frontHTML, railHTML, appPath, stackPath, parseAppQuery, addressTarget, lookOf, crumbsHTML, geometry: { PL, PR, PW } };
+    adaptModel, layoutGraph, stackBars, barHTML, frontHTML, railHTML, appPath, stackPath, parseAppQuery, addressTarget, lookOf, crumbsHTML, geometry: { PL, PR, PW },
+    VERB_ROUTE_BUDGET_MS, verbOutcome, askVerb, verbAct, verbPanel };
   if (typeof document === 'undefined') return;
 
   // ================= The page =================
@@ -925,7 +989,7 @@
   const addrOf = () => F.kind === 'open' ? String(F.id).replace(/:episode:/, ':') : `${String(F.brane).replace(/^([a-z][a-z0-9-]*):brane:/, '$1:')}@${railIdx() + 1}`;
   function renderCrumbs() {
     $('crumbs').innerHTML = crumbsHTML(st.d === 1 && F ? addrOf() : null);
-    $('addr').placeholder = st.d !== 0 ? '' : T && T.GROUPS[0].id ? 'go to a place — ' + T.GROUPS.map(g => g.id).join(' · ') + ' · ' + T.GROUPS[0].id + ':row-N@1' : 'go to a place — row-189@6 · shared@0ffbb9c';
+    $('addr').placeholder = st.d !== 0 ? '' : (T && T.GROUPS[0].id ? 'go to a place — ' + T.GROUPS.map(g => g.id).join(' · ') + ' · ' + T.GROUPS[0].id + ':row-N@1' : 'go to a place — row-189@6 · shared@0ffbb9c') + ' · episode report <episode>';
   }
   let msgT;
   function say(t) { const m = $('hmsg'); m.textContent = t; m.style.opacity = '1'; clearTimeout(msgT); msgT = setTimeout(() => { m.style.opacity = '0'; }, 2600); }
@@ -1198,7 +1262,21 @@
     if (!F || !RAIL[k] || RAIL[k].id === F.id) return;
     openState(F.brane, RAIL[k].id);
   }
+  // go is the bar's Enter: empty is home (no fetch); anything else asks the verb route first.
   async function go(raw) {
+    if (String(raw == null ? '' : raw).trim() === '') return out();
+    const act = await askVerb(raw, { fetch: (u, o) => fetch(u, o), AbortController, setTimeout: (f, ms) => setTimeout(f, ms), clearTimeout: t => clearTimeout(t) });
+    return verbAct(act, raw, { address: goAddress, assign: h => location.assign(h), panel: showVerbPanel, say,
+      clear: () => { if ($('addr').value === raw) $('addr').value = ''; } });
+  }
+  function showVerbPanel(p) {
+    const old = $('verb'), el = verbPanel(document, p);
+    if (old) old.replaceWith(el); else $('stage').appendChild(el);
+    const b = el.querySelector('button');
+    if (b) b.addEventListener('click', () => { try { navigator.clipboard.writeText(b.dataset.copy).then(() => say('copied'), () => say('copy it from the panel')); } catch (_) { say('copy it from the panel'); } });
+  }
+  // goAddress is the address bar's own path (M4a's grammar), reached only through the route's address answer.
+  async function goAddress(raw) {
     const parse = SEAMS.parseAddress;
     let r = addressTarget(raw, null, parse);
     if (r.error && typeof parse === 'function') {
@@ -1210,7 +1288,7 @@
   }
 
   // ---------- Popovers ----------
-  const KEYROWS = [['← →', 'back · forward a page', '—'], ['⇧← ⇧→', 'previous · next state', '—'], ['↑ ↓', 'lane above · below', 'back in time · toward now'], ['+ −', 'zoom in · out', '—'], ['End', 'back to now', 'most recent'], ['Home', 'the oldest state', 'earliest'], ['Enter', 'open brane', '—'], ['Esc', '—', 'back to timeline'], ['N', 'what needs you', 'what needs you'], ['/', 'address bar', 'address bar'], [',', 'appearance', 'appearance'], ['?', 'these keys', 'these keys']];
+  const KEYROWS = [['← →', 'back · forward a page', '—'], ['⇧← ⇧→', 'previous · next state', '—'], ['↑ ↓', 'lane above · below', 'back in time · toward now'], ['+ −', 'zoom in · out', '—'], ['End', 'back to now', 'most recent'], ['Home', 'the oldest state', 'earliest'], ['Enter', 'open brane', '—'], ['Esc', '—', 'back to timeline'], ['N', 'what needs you', 'what needs you'], ['/', 'address bar', 'address bar'], ['Enter', 'in the bar: a place, episode report E, approve ID', 'the same'], [',', 'appearance', 'appearance'], ['?', 'these keys', 'these keys']];
   function renderKeys() {
     $('keys').innerHTML = `<div class="kgrid"><span></span><span style="font-size:12px;font-weight:700;color:${st.d === 0 ? 'var(--text)' : 'var(--faint)'}">Timeline</span><span style="font-size:12px;font-weight:700;color:${st.d === 1 ? 'var(--text)' : 'var(--faint)'}">Brane</span>` +
       KEYROWS.map(([k, a, b]) => `<span><span class="kbd">${k}</span></span><span>${a}</span><span>${b}</span>`).join('') +
@@ -1427,7 +1505,7 @@
   $('look').addEventListener('click', e => { const b = e.target.closest('[data-look]'); if (!b) return; const [k, v] = b.dataset.look.split(':'); look[k] = v; saveLook(); applyLook(); renderLook(); });
   $('addr').addEventListener('keydown', e => {
     e.stopPropagation();
-    if (e.key === 'Enter') { go(e.target.value); e.target.value = ''; $('root').focus({ preventScroll: true }); }
+    if (e.key === 'Enter') { go(e.target.value); $('root').focus({ preventScroll: true }); }
     if (e.key === 'Escape') { e.target.value = ''; $('root').focus({ preventScroll: true }); }
   });
   const onRow = row => Object.keys(DISP).filter(id => DISP[id].row === row).sort((a, b) => DISP[a].x - DISP[b].x);
@@ -1439,7 +1517,7 @@
     if (k === '?') toggle('keys');
     else if (k === ',') toggle('look');
     else if (k === '/') $('addr').focus();
-    else if (k === 'Escape') { if (!$('keys').hidden || !$('look').hidden) { $('keys').hidden = true; $('look').hidden = true; } else out(); }
+    else if (k === 'Escape') { if ($('verb')) $('verb').remove(); else if (!$('keys').hidden || !$('look').hidden) { $('keys').hidden = true; $('look').hidden = true; } else out(); }
     else if (k === 'n' || k === 'N') { const a = document.querySelector('[data-needs] a'); if (a) { if (st.d === 1) out(); a.focus(); } }   // design §2.8: N jumps to what needs you
     else if (st.d === 1 && k === 'ArrowUp') openRail(railIdx() - 1);
     else if (st.d === 1 && k === 'ArrowDown') openRail(railIdx() + 1);
diff --git a/host/surfacepage/app/index.html b/host/surfacepage/app/index.html
index 7823159..f7a4fc9 100644
--- a/host/surfacepage/app/index.html
+++ b/host/surfacepage/app/index.html
@@ -64,7 +64,7 @@
     <nav class="crumbs" id="crumbs" aria-label="Where you are"></nav>
     <label class="addr">
       <span class="prompt">world&gt;</span>
-      <input id="addr" type="text" aria-label="Go to an address" spellcheck="false" autocomplete="off" placeholder="row-189@6 · shared@0ffbb9c · back">
+      <input id="addr" type="text" aria-label="Go to an address" spellcheck="false" autocomplete="off" placeholder="row-189@6 · shared@0ffbb9c · back · episode report &lt;episode&gt;">
     </label>
     <span class="hmsg" id="hmsg"></span>
     <button type="button" class="hbtn mono" id="keys-btn" aria-label="Keyboard shortcuts" aria-keyshortcuts="?">?</button>
diff --git a/host/surfacepage/app_test.go b/host/surfacepage/app_test.go
index bc1d07d..e191553 100644
--- a/host/surfacepage/app_test.go
+++ b/host/surfacepage/app_test.go
@@ -172,7 +172,8 @@ const model = JSON.parse(fs.readFileSync('model.json', 'utf8'));
 const stacks = JSON.parse(fs.readFileSync('stacks.json', 'utf8'));
 const index = JSON.parse(fs.readFileSync('index.json', 'utf8'));
 assert.deepEqual(Object.keys(app).sort(), ['LEVELS', 'adaptModel', 'addressTarget', 'agentLine', 'appPath', 'applyDelta', 'barHTML', 'blendXS', 'cleanTitle', 'crumbsHTML', 'decidedHTML', 'esc', 'fetchCoherent', 'frontHTML',
-  'geometry', 'isSharedKey', 'layoutGraph', 'levelXs', 'liveRequest', 'liveStatus', 'lookOf', 'needsYouHTML', 'nextBackoff', 'notifications', 'openRunningHTML', 'parseAppQuery', 'prettyModel', 'railHTML', 'rowLabel', 'serviceLineHTML', 'shortTitle', 'sourcesHTML', 'splitFrames', 'stackBars', 'stackPath', 'summaryOf', 'verdict']);
+  'geometry', 'isSharedKey', 'layoutGraph', 'levelXs', 'liveRequest', 'liveStatus', 'lookOf', 'needsYouHTML', 'nextBackoff', 'notifications', 'openRunningHTML', 'parseAppQuery', 'prettyModel', 'railHTML', 'rowLabel', 'serviceLineHTML', 'shortTitle', 'sourcesHTML', 'splitFrames', 'stackBars', 'stackPath', 'summaryOf', 'verdict',
+  'VERB_ROUTE_BUDGET_MS', 'askVerb', 'verbAct', 'verbOutcome', 'verbPanel'].sort());
 
 assert.equal(app.serviceLineHTML(undefined), '');
 // Row 201 A3: Needs you is the pending requests, each a link to its card; a refused row shows why, with no link.
diff --git a/host/surfacepage/app_verb_test.go b/host/surfacepage/app_verb_test.go
new file mode 100644
index 0000000..7255c24
--- /dev/null
+++ b/host/surfacepage/app_verb_test.go
@@ -0,0 +1,178 @@
+package surfacepage
+
+import (
+	"context"
+	"os"
+	"os/exec"
+	"path/filepath"
+	"strings"
+	"testing"
+	"time"
+)
+
+// TestSurfaceAppVerbs (row 243 V1, design §2.3, §6 VP1–VP6) runs the bar's verb part of app.js under
+// node against GET /surface/verb.json's real answers (testdata/verb-route.golden.json, which
+// host/daemon's TestSurfaceVerbOutcomes checks byte for byte): only a schema-valid address answer
+// reaches the address path; every failure, timeout and superseded ask leaves the view and the input.
+// It skips only when no node exists; with WORLD_EXEC_NODE set, a harness that cannot run fails.
+func TestSurfaceAppVerbs(t *testing.T) {
+	node := os.Getenv("WORLD_EXEC_NODE")
+	if node == "" {
+		var err error
+		node, err = exec.LookPath("node")
+		if err != nil {
+			t.Skip("node unavailable locally")
+		}
+	}
+	dir := t.TempDir()
+	files := map[string][]byte{"app.js": AppScript, "surface.js": Script, "harness.js": []byte(verbHarness)}
+	for name, src := range map[string]string{"golden.json": filepath.Join("testdata", "verb-route.golden.json"),
+		"index.json": filepath.Join("..", "surface", "testdata", "index-rig.json")} {
+		b, err := os.ReadFile(src)
+		if err != nil {
+			t.Fatal(err)
+		}
+		files[name] = b
+	}
+	for name, data := range files {
+		if err := os.WriteFile(filepath.Join(dir, name), data, 0o600); err != nil {
+			t.Fatal(err)
+		}
+	}
+	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
+	defer cancel()
+	cmd := exec.CommandContext(ctx, node, filepath.Join(dir, "harness.js"))
+	cmd.Dir = dir
+	out, err := cmd.CombinedOutput()
+	if err != nil {
+		t.Fatalf("node harness: %v\n%s", err, out)
+	}
+	if !strings.Contains(string(out), "APP VERBS PASS") {
+		t.Fatalf("harness manifest missing: %s", out)
+	}
+}
+
+const verbHarness = `
+const assert = require('node:assert/strict');
+const fs = require('node:fs');
+const app = require('./app.js');
+const seams = require('./surface.js');
+const index = JSON.parse(fs.readFileSync('index.json', 'utf8'));
+const G = JSON.parse(fs.readFileSync('golden.json', 'utf8'));
+const B = q => JSON.parse(G[q]);
+const LONG = 'x'.repeat(1025), MR = 'approve mr-0123456789abcdef', REF = 'approve sha256:' + 'ab'.repeat(32);
+const UNREAD = 'the verb route answered something this page cannot read', LATE = 'the verb route did not answer in time';
+// io: a fetch stub and a controlled clock (setTimeout records, nothing is scheduled).
+function rig(fetchImpl) {
+  const io = { calls: [], timers: [], cleared: [], AbortController };
+  io.fetch = (u, o) => { io.calls.push({ u, o }); return fetchImpl(u, o); };
+  io.setTimeout = (fn, ms) => { io.timers.push({ fn, ms }); return io.timers.length; };
+  io.clearTimeout = id => { io.cleared.push(id); };
+  return io;
+}
+const reply = (status, text) => () => Promise.resolve({ status, text: () => Promise.resolve(text) });
+const golden = q => reply(200, G[q]);
+function spies() { const fx = { log: [] }; for (const k of ['address', 'assign', 'panel', 'say', 'clear']) fx[k] = (...a) => { fx.log.push([k, ...a]); }; return fx; }
+async function run(raw, fetchImpl) {
+  const io = rig(fetchImpl), act = await app.askVerb(raw, io), fx = spies();
+  app.verbAct(act, raw, fx);
+  assert.equal(io.calls.length, 1, raw); // the page never short-circuits on its own word list (VP4)
+  assert.equal(io.calls[0].u, '/surface/verb.json?q=' + encodeURIComponent(raw));
+  assert(io.calls[0].o.signal instanceof AbortSignal, raw);
+  assert.deepEqual(io.timers.map(x => x.ms), [app.VERB_ROUTE_BUDGET_MS], raw);
+  assert.deepEqual(io.cleared, [1], 'the timer is cleared: ' + raw);
+  return { act, log: fx.log };
+}
+(async () => {
+  assert.equal(app.VERB_ROUTE_BUDGET_MS, 15000);
+  // VP1, VP4, VP6: the address answers reach the address path with the raw line, after clearing the bar.
+  for (const q of ['row-189@6', 'world shared@0ffbb9c', 'back', 'world why', 'world session', 'world episode']) {
+    const { act, log } = await run(q, golden(q));
+    assert.deepEqual(act, { address: true }, q); assert.deepEqual(log, [['clear'], ['address', q]], q);
+  }
+  // Must-live: that path gives today's answers.
+  const at = raw => app.addressTarget(raw, index, seams.parseAddress);
+  assert.deepEqual(at('row-189@6'), { open: { brane: 'row-189', at: '652c341d5072c5f49fefe29f9409c93a5119a5be' } });
+  assert.deepEqual(at('world shared@0ffbb9c'), { open: { brane: 'shared', at: '0ffbb9c9d386f83aaa45b98e4d4b7763a5c5b845' } });
+  assert.deepEqual(at('back'), { home: true });
+  // Results: the CLI's text, its JSON, its not-found line with exit 1.
+  for (const [q, body] of [['episode report w1', B('episode report w1').text], ['episode report w1 --json', B('episode report w1 --json').json],
+    ['episode report nobody', B('episode report nobody').stderr + 'exit 1']]) {
+    const { act, log } = await run(q, golden(q));
+    assert.deepEqual(act, { panel: { head: '$ ' + B(q).cli, body, copy: '' } }, q); assert.deepEqual(log, [['clear'], ['panel', act.panel]], q);
+  }
+  assert(B('episode report nobody').stderr.endsWith('has no committed calls\n'));
+  // Terminal: the CLI line, with a copy.
+  for (const q of ['why', 'episode finish w1 --preset p']) {
+    const { act, log } = await run(q, golden(q)), b = B(q);
+    assert.deepEqual(act, { panel: { head: '$ ' + b.cli, body: b.reason + '\nAt a terminal: ' + b.cli, copy: b.cli } }, q);
+    assert.deepEqual(log, [['clear'], ['panel', act.panel]], q);
+  }
+  // Navigate: the two card prefixes (VP5's must-live).
+  for (const [q, href] of [[MR, '/surface/approve/mr-0123456789abcdef'], [REF, '/surface/decide/sha256:' + 'ab'.repeat(32)]]) {
+    const { act, log } = await run(q, golden(q));
+    assert.equal(B(q).href, href); assert.deepEqual(act, { assign: href }); assert.deepEqual(log, [['clear'], ['assign', href]]);
+  }
+  // VP5: any other target is refused.
+  for (const href of ['https://evil/x', '/surface/model.json', '//evil/surface/approve/mr-']) {
+    const { log } = await run(MR, reply(200, JSON.stringify({ ...B(MR), href })));
+    assert.deepEqual(log, [['say', 'refused: unexpected target']], href);
+  }
+  // Refused (VP6, the shadow set; TooLong): the reason is said; nothing else happens.
+  for (const [q, has] of [['log', 'at a terminal: ailang-worldd log'], ['ailang-worldd world get world-1', 'ailang-worldd world get world-1'],
+    ['episode', 'needs a subcommand'], ['world a b c', 'unknown address: try "world"'], [LONG, 'TooLong']]) {
+    const { act, log } = await run(q, golden(q));
+    assert.deepEqual(act, { say: B(q).reason }, q); assert.deepEqual(log, [['say', B(q).reason]], q); assert(B(q).reason.includes(has), q);
+  }
+  assert(B(LONG).reason.startsWith('TooLong'));
+  // VP2: no silent fallback. Every failure says its message and calls nothing else, for verbs and addresses alike.
+  for (const q of ['why', 'log', 'row-189@6', 'world shared@0ffbb9c']) {
+    const b = B(q), drop = k => { const c = { ...b }; delete c[k]; return JSON.stringify(c); };
+    const failures = [404, 400, 421, 500].map(s => [reply(s, G[q]), 'the verb route did not answer: HTTP ' + s]).concat([
+      [() => Promise.reject(new TypeError('Failed to fetch')), 'the verb route is unreachable: Failed to fetch'],
+      [reply(200, '<html>not json</html>'), UNREAD], [reply(200, drop('schema')), UNREAD], [reply(200, JSON.stringify({ ...b, outcome: 'run' })), UNREAD], [reply(200, JSON.stringify({ ...b, outcome: 'constructor' })), UNREAD],
+      [reply(200, JSON.stringify({ ...b, v: 2 })), UNREAD], [reply(200, 'null'), UNREAD], [reply(200, JSON.stringify([b])), UNREAD]]);
+    for (const [f, msg] of failures) { const { log } = await run(q, f); assert.deepEqual(log, [['say', msg]], q + ': ' + msg); }
+  }
+  // An answer missing its outcome's own members is unreadable too.
+  for (const body of [{ ...B('why'), cli: '' }, { ...B('why'), reason: undefined }, { ...B(MR), href: 7 }, { ...B('log'), reason: '' },
+    { ...B('episode report w1'), json: '{}' }, { ...B('episode report nobody'), exit: 0 }, { ...B('episode report w1'), text: '' }]) {
+    const { log } = await run('why', reply(200, JSON.stringify(body)));
+    assert.deepEqual(log, [['say', UNREAD]], JSON.stringify(body).slice(0, 80));
+  }
+  // VP2 hangs, on a controlled clock: stalled headers, then a stalled body. Firing the budget's timer aborts the ask.
+  for (const fetchImpl of [() => new Promise(() => {}), () => Promise.resolve({ status: 200, text: () => new Promise(() => {}) })]) {
+    const io = rig(fetchImpl), p = app.askVerb('why', io);
+    await new Promise(r => setImmediate(r));
+    assert.deepEqual(io.timers.map(x => x.ms), [15000]); assert.deepEqual(io.cleared, []); assert(!io.calls[0].o.signal.aborted);
+    io.timers[0].fn();
+    const act = await p, fx = spies(); app.verbAct(act, 'why', fx);
+    assert.deepEqual(fx.log, [['say', LATE]]); assert.deepEqual(io.cleared, [1]); assert(io.calls[0].o.signal.aborted);
+  }
+  // Superseded (planner S-D1): a second ask aborts the first, whose answer is dropped (no say, no navigation),
+  // whether it was still waiting (A) or had already arrived (B); the second answers for the bar.
+  let release; const io1 = rig(() => new Promise(r => { release = r; }));
+  const p1 = app.askVerb('row-189@6', io1), p2 = app.askVerb('log', rig(golden('log')));
+  assert(io1.calls[0].o.signal.aborted);
+  release({ status: 200, text: () => Promise.resolve(G['row-189@6']) });
+  const a1 = await p1, a2 = await p2;
+  assert.deepEqual(a1, { superseded: true }); assert.deepEqual(io1.cleared, [1]); assert.deepEqual(a2, { say: B('log').reason });
+  { const fx = spies(); app.verbAct(a1, 'row-189@6', fx); assert.deepEqual(fx.log, []); }
+  const io3 = rig(golden('row-189@6')), p3 = app.askVerb('row-189@6', io3), p4 = app.askVerb('why', rig(golden('why')));
+  assert.deepEqual(await p3, { superseded: true }); assert.deepEqual(io3.cleared, [1]); assert(Boolean((await p4).panel));
+  // VP3: the panel is built with textContent only; a hostile report or reason stays text.
+  const fakeDoc = () => { const doc = { made: [], createElement(tag) {
+    const el = { tag, children: [], dataset: {}, style: {}, attrs: {}, textContent: '', setAttribute(k, v) { this.attrs[k] = v; }, append(...c) { this.children.push(...c); } };
+    Object.defineProperty(el, 'innerHTML', { get() { throw new Error('innerHTML read'); }, set() { throw new Error('innerHTML set'); } });
+    doc.made.push(el); return el; } }; return doc; };
+  { const doc = fakeDoc(), el = app.verbPanel(doc, app.verbOutcome(200, G['episode report w1']).panel);
+    assert.deepEqual(doc.made.map(e => e.tag), ['div', 'div', 'pre']); assert.equal(el.id, 'verb'); assert.equal(el.children[0].textContent, '$ ailang-worldd episode report w1');
+    assert(el.children[1].textContent.includes('#2 <img src=x onerror=alert(1)>: <b>hostile</b> "reason"')); assert.equal(el.children[1].textContent, B('episode report w1').text); }
+  { const doc = fakeDoc(), el = app.verbPanel(doc, app.verbOutcome(200, G.why).panel);
+    assert.deepEqual(doc.made.map(e => e.tag), ['div', 'div', 'pre', 'button']); assert.equal(el.children[2].dataset.copy, 'ailang-worldd why');
+    assert.equal(el.children[2].textContent, 'Copy the line'); }
+  { const doc = fakeDoc(), el = app.verbPanel(doc, { head: '$ <b>x</b>', body: '<img src=x onerror=1>', copy: '' });
+    assert.equal(el.children[1].textContent, '<img src=x onerror=1>'); assert.equal(doc.made.length, 3); }
+  console.log('APP VERBS PASS');
+})().catch(e => { console.error(e); process.exit(1); });
+`
````

### Patch M4b — live Chrome (XV1); QUICKSTART §19

````diff M4b
diff --git a/cmd/ailang-worldd/quickstart_verbs_test.go b/cmd/ailang-worldd/quickstart_verbs_test.go
new file mode 100644
index 0000000..37dbf59
--- /dev/null
+++ b/cmd/ailang-worldd/quickstart_verbs_test.go
@@ -0,0 +1,55 @@
+package main
+
+// Row 243 V1 (design §7, §10): QUICKSTART §19 says what the bar does with each line it shows, and
+// host/verbs, the page and the daemon agree.
+
+import (
+	"os"
+	"strings"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/surfacepage"
+	"github.com/sunholo-data/ailang-world/host/verbs"
+)
+
+func TestQuickstartSection19MatchesTheBar(t *testing.T) {
+	raw, err := os.ReadFile("../../docs/QUICKSTART.md")
+	if err != nil {
+		t.Fatal(err)
+	}
+	guide := string(raw)
+	s18, s19 := strings.Index(guide, "\n## 18. "), strings.Index(guide, "\n## 19. Operator verbs in the address bar *(attended — pending first verbatim run)*\n")
+	if s18 < 0 || s19 < s18 {
+		t.Fatalf("§19 is missing or not after §18 (%d, %d)", s18, s19)
+	}
+	section := guide[s19:]
+	_, block, _ := strings.Cut(section, "```text\n")
+	block, _, _ = strings.Cut(block, "```")
+	want := map[string]verbs.Class{"episode report": verbs.Read, "approve": verbs.Navigate, "episode finish": verbs.Mutate, "": verbs.Address}
+	seen := map[verbs.Class]bool{}
+	for _, line := range strings.Split(strings.TrimSpace(block), "\n") {
+		argv, err := verbs.Split(line)
+		if err != nil {
+			t.Fatalf("§19 line %q: %v", line, err)
+		}
+		v, _ := verbs.Lookup(argv)
+		if c, ok := want[v.Words]; !ok || c != v.Class {
+			t.Errorf("§19 line %q is %q %s, which §19 does not describe", line, v.Words, v.Class)
+		}
+		seen[v.Class] = true
+	}
+	if len(seen) != len(want) {
+		t.Errorf("§19's lines cover %v, want each of %v", seen, want)
+	}
+	for _, w := range verbs.NeverWords {
+		if !strings.Contains(section, "`"+w+"`") {
+			t.Errorf("§19 does not name %s among the commands that stay at a terminal", w)
+		}
+	}
+	if !strings.Contains(section, "15 seconds") || !strings.Contains(string(surfacepage.AppScript), "const VERB_ROUTE_BUDGET_MS = 15000;") {
+		t.Error("§19's 15 seconds is not the page's budget")
+	}
+	if !strings.Contains(string(surfacepage.AppShell), "episode report") || !strings.Contains(section, "GET /surface/verb.json") {
+		t.Error("the bar's placeholder or §19 does not name the verb route's first verb")
+	}
+}
diff --git a/cmd/ailang-worldd/verb_bar_chrome_test.go b/cmd/ailang-worldd/verb_bar_chrome_test.go
new file mode 100644
index 0000000..c4e8d21
--- /dev/null
+++ b/cmd/ailang-worldd/verb_bar_chrome_test.go
@@ -0,0 +1,78 @@
+package main
+
+// Row 243 V1 (design §6 XV1): the app's address bar runs episode report in a live Chrome, on the
+// passkey pages' driver (testdata/passkey_cdp.mjs) and skip rule; the daemon is a real one on a test
+// store seeded by VR3's seedReportCalls, served on the loopback.
+
+import (
+	"fmt"
+	"io"
+	"net/http/httptest"
+	"path/filepath"
+	"strings"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/broker"
+	"github.com/sunholo-data/ailang-world/host/daemon"
+)
+
+// XV1: `episode report w1` typed into #addr shows the CLI's own bytes in the verb panel and clears the
+// bar; `why` shows its CLI line; `log` is refused in the header with its line, the bar and the place kept.
+func TestVerbBarEpisodeReportInChrome(t *testing.T) {
+	b := startPasskeyBrowser(t)
+	b.do(t, nil, "tab", map[string]any{})
+	db := filepath.Join(t.TempDir(), "world.db")
+	seedReportCalls(t, db, [][5]string{{"w1", "1", "ailang-read", broker.EffectWorkspaceRead, "ok"}, {"w1", "2", "ailang-read", "", ""}})
+	d, err := daemon.New(boundedTestContext(t), daemon.Config{DBPath: db, BindHost: daemon.DefaultBindHost, ErrorLog: io.Discard})
+	if err != nil {
+		t.Fatal(err)
+	}
+	srv := httptest.NewServer(d.Handler())
+	t.Cleanup(func() { srv.Close(); _ = d.Close() })
+	code, want, stderr := runCLI(t, srv.URL, "episode", "report", "w1")
+	if code != exitOK || !strings.Contains(want, "episode w1: 2 calls") {
+		t.Fatalf("the CLI: exit %d stdout %q stderr %q", code, want, stderr)
+	}
+	b.do(t, nil, "openCard", map[string]any{"url": srv.URL + "/surface", "raw": true})
+	b.do(t, nil, "waitFor", map[string]any{"expression": `Boolean(window.World && document.getElementById('addr'))`})
+	enter := func(line string) {
+		t.Helper()
+		b.do(t, nil, "waitFor", map[string]any{"expression": fmt.Sprintf(`(() => { const e = document.getElementById('addr'); e.focus(); e.value = %q; return document.activeElement === e; })()`, line)})
+		b.do(t, nil, "key", map[string]any{"key": "Enter"})
+	}
+	panel := func(has string) string {
+		t.Helper()
+		var s string
+		b.do(t, &s, "waitFor", map[string]any{"expression": fmt.Sprintf(`(() => { const p = document.querySelector('#verb pre'); return p && p.textContent.includes(%q) && p.textContent; })()`, has)})
+		return s
+	}
+	enter("episode report w1")
+	if got := panel("episode w1: 2 calls"); got != want {
+		t.Fatalf("the bar shows\n%q\nthe CLI printed\n%q", got, want)
+	}
+	var state struct {
+		Addr, Head, Location string
+		Verb                 bool
+	}
+	read := func() {
+		t.Helper()
+		b.do(t, &state, "waitFor", map[string]any{"expression": `({ addr: document.getElementById('addr').value, head: (document.querySelector('#verb .mono') || {}).textContent || '',
+			location: location.href, verb: Boolean(document.getElementById('verb')) })`})
+	}
+	if read(); state.Addr != "" || state.Head != "$ ailang-worldd episode report w1" {
+		t.Fatalf("after episode report: %+v", state)
+	}
+	enter("why")
+	if got := panel("At a terminal: ailang-worldd why"); !strings.HasPrefix(got, "why reads World") {
+		t.Fatalf("why shows %q", got)
+	}
+	read()
+	before := state.Location
+	b.do(t, nil, "waitFor", map[string]any{"expression": `(document.getElementById('verb').remove(), 1)`})
+	enter("log")
+	b.do(t, nil, "waitFor", map[string]any{"expression": `document.getElementById('hmsg').textContent.includes('at a terminal: ailang-worldd log')`})
+	if read(); state.Verb || state.Addr != "log" || state.Location != before {
+		t.Fatalf("after log: %+v (was at %s)", state, before)
+	}
+	b.do(t, nil, "openCard", map[string]any{"url": "about:blank", "raw": true}) // ends the page's live stream before the server closes
+}
diff --git a/docs/QUICKSTART.md b/docs/QUICKSTART.md
index 59b8a46..706b940 100644
--- a/docs/QUICKSTART.md
+++ b/docs/QUICKSTART.md
@@ -1895,3 +1895,17 @@ ailang-worldd approve --list --db /tmp/se-world/world.db
 ```
 
 The line opens the decision card at `/surface/decide/<request-ref>`. The card reads the request and its packet from the daemon, checks that both hash to their refs and that the packet carries the request, shows every field of both, and refuses a field it cannot show. **Approve with passkey** or **Deny with passkey** signs that verb's decision message with a §16 passkey, and the card says which verb it is signing. **Save the approve message** (or the deny message) downloads the 91 bytes for the `ssh-keygen -Y sign` line the card prints; paste the signature and press **Approve with this signature** (or **Deny with this signature**). The daemon checks every decision as it checks the CLI's. A publish- or finish-shaped request is listed with its typed owner and has no link. A request the card cannot decide shows the daemon's reason and no button.
+
+## 19. Operator verbs in the address bar *(attended — pending first verbatim run)*
+
+The surface app's address bar (`/` focuses it) takes the CLI's own operator verbs, in the CLI's words; a line pasted with a leading `ailang-worldd` works too. The daemon reads the line with the same grammar and flag set the CLI uses (`GET /surface/verb.json`), so the bar and a terminal cannot disagree. No verb changes World from the page.
+
+```text
+episode report w1
+episode report w1 --json
+approve mr-0123456789abcdef
+episode finish w1 --preset se-tools
+world shared@0ffbb9c
+```
+
+The first prints the report exactly as `ailang-worldd episode report w1` prints it at a terminal, and `--json` prints its JSON; an episode with no committed call shows the CLI's own line and `exit 1`. `approve <request-id|request-ref>` opens that request's card (§17, §18), where the decision is your key's. Every other verb shows the exact CLI line and says the app cannot run it yet, with a copy button: `episode start`, `episode finish` and `why` today. Global flags, `--follow`, `--markdown` and `--addr` are refused, as are the commands that stay at a terminal (`approver`, `session`, `service`, `serve`, `setup`). A brane named like a CLI verb (`log`, `episode`, …) is reached as `world NAME`, which always navigates. If the daemon does not answer within 15 seconds, or answers something the page cannot read, the bar says so and keeps your line; it never falls back to navigation.
````

