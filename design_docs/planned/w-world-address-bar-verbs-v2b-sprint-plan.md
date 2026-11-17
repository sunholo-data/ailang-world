# Sprint plan — w-world-address-bar-verbs-v2, sprint V2b "the bar" (iteration 398, row 243)

> **THE SPRINT — row 243 phase V2, sprint V2b (design §0.3 A1–A5, §0.4, §2.2–§2.7, §3–§5, §6 "V2b — the bar", §7, §8
> M4–M6, §9, §10, §13.2; V2a plan Appendix C):** typing **`why <target>`** in the World app's address bar answers with
> **the bytes `ailang-worldd why <target>` prints** — its stdout (text, or JSON with `--json`), or its stderr line, and
> its exit (0, 1 or 3) — produced by **`why.Run`, the CLI's own walk**, which the daemon calls **in-process** through a
> private mux over the seven `/v1` GET handlers the walk reads (`GET /surface/verb.json`, V1's route; no new route, no
> POST, no store write, no process, no `.ail`). The CLI's `why` parse moves to `host/verbs` (design §2.2, moved here from
> V2a) so the bar and the terminal share one parse; `why -`, `--result` and `--expect-store` are refused with the CLI line.
> A walk that passes the daemon's 10 s read deadline, or a read that fails, is a **refusal with the CLI line** — never
> "not found", never a ✗ chain. Every `result` answer is capped at **1 MiB** (A2). The page reads `exit 3` (A2), opens
> the panel for a refusal that carries a CLI line (A3) and lets **every Enter own the bar** (A4).
> - **Measured: 1,255 changed lines** — `git diff -M --numstat c09e0cdb HEAD -- cmd host docs` (insertions +
>   deletions, git's rename detection on; no rename occurs), after a commit of each of the five scratch stages:
>   **343 production** (`host/daemon/surface_verb_why.go` 163, `host/verbs/flags.go` 85, `app.js` 58, `why_cli.go` 25,
>   `surface_verb.go` 6, `lookup.go` 4, `index.html` 2), **908 test** (incl. the 6-line golden change), **4 docs**.
>   **Cap 1,400: margin 145. No cut** (budget ≈ 1,135; the over-run is test, §1.5).
> - **Five sub-milestones, each green at its commit** (§1.2 C11, incl. live Chrome where it applies): **M4** admission
>   (§2.2 + `AdmitWhy`) → **M6p** the page (A2–A4; the route's `why` is still `terminal`) → **M5g** the getter and the
>   cap → **M5r** the route's `why` → **M6e** the CLI's bytes (VA3, VE1, XV2) and QUICKSTART §19.
> - **Audit inputs 4–6 and the V2a judge's notes are resolved by measurement** (§1.3 V2b-D3 … D9, §1.2 C5–C9); V2a's
>   three named survivors are re-run here: `W-MAXRESP` and `W-SCANMAX` are now **killed behaviourally** (VR5, VA3);
>   `W-LOGPAGE` survives as before (Appendix C).
> - **Every line of code, test, golden and doc is a unified diff in Appendix P**, generated from compiled, vetted,
>   gofmt-clean, green scratch stages; extracted **from this file** with §2's `awk`, applied with `patch -p1 -E` to a
>   fresh `git archive c09e0cdb` export and compared byte for byte with the stage tree (§1.2 C12). The executor applies;
>   it writes nothing.

- **Design**: [w-world-address-bar-verbs-v2.md](w-world-address-bar-verbs-v2.md), **r1** (quorum r2 PROCEED 3/3, no
  absentee; independent audit PASS 84; iteration 396; **not reopened**). Where this plan and the design differ, this plan
  was measured at `c09e0cdb` and wins; each difference is a decision (§1.3) or a premise found false (§1.4).
- **Previous sprint**: V2a (the seam) landed at iteration 397 (PR #585, merge `2d530cbe`, judge PASS 90):
  [w-world-address-bar-verbs-v2a-sprint-plan.md](w-world-address-bar-verbs-v2a-sprint-plan.md); its Appendix C is this
  sprint's handoff.
- **Authority**: D-WORLD-107 = A, D-WORLD-109 = A (test stores, temp dirs and test keys only), D-WORLD-89 R1–R5
  untouched (no POST, write or key path). No ledger ask (design Q4). A3 within D-WORLD-107 (design Q3).
- **Planner base**: `c09e0cdb3c8a3eb25c208dbacedddfb690fe8ca0` (= origin/dev, PR #586's merge). Worktree
  `/Users/voightkampff/dev/sunholo-data/.wt-world-iter398-row243-v2b`, branch `sprint/world-iter398-row243-v2b`.
  `git diff --stat 2d530cbe c09e0cdb -- cmd host docs` is empty (§1.2 C2): V2a's merged tree is the code; every design
  line this plan relies on was re-measured at `c09e0cdb` (§1.2 C3).
- **Planner scratch** (outside the repo, never under `/tmp`): `$HOME/.ailang/state/world-iter398/plan-scratch/`
  (`base/` the pristine `git archive c09e0cdb` export, a throwaway repo tagged `BASE`; `work/` the final tree;
  `work2/` the stage repo, tags `M4 M6p M5g M5r M6e`, the tree every mutant cloned; `final-tree/` its clone for the
  final gates; `rt/` the round trip; `patches/`; `snaptest/` the snapshot rehearsal; `tools/` (`stages.py`,
  `mkpatches.sh`, `accept.sh`, `accept_all.sh`, `gates.sh`, `mutants.py`, `mutate.py`, `roundtrip.sh`,
  `snap_snippet.sh`, `assemble.py`, `plan_template.md`); `logs/` (`base_gates/`, `final_gates/`, `accept_<M>.log`,
  `accept_all.txt`, `mutant_ledger.tsv`, `mut/`, `roundtrip.log`)).
- **Executor**: as routed by the controller. **Two runs**: Run 1 = M4, M6p, M5g; Run 2 = M5r, M6e. Each sub-milestone is
  one `patch -p1 -E` (§2) plus its acceptance (§3). **Apply, never retype.**
- **Touched** (§5): new `host/verbs/{why_admit_test.go, bounds_test.go}`, `host/daemon/{surface_verb_why.go,
  surface_verb_why_test.go, surface_verb_why_route_test.go}`, `cmd/ailang-worldd/verb_why_e2e_test.go`; edited
  `host/verbs/{flags.go, lookup.go, lookup_test.go}`, `cmd/ailang-worldd/{why_cli.go, why_seam_test.go,
  verb_bar_chrome_test.go, quickstart_verbs_test.go}`, `host/daemon/{surface_verb.go, surface_verb_test.go,
  surface_verb_outcomes_test.go}`, `host/surfacepage/{app/app.js, app/index.html, app_test.go, app_verb_test.go,
  testdata/verb-route.golden.json}`, `docs/QUICKSTART.md`. **Outside the design's §7 lists** (each in §1.4 P1):
  `surface_verb_test.go`'s 1,024-byte row, `app_test.go`'s export list, `verb_bar_chrome_test.go` (XV1),
  `quickstart_verbs_test.go`, `why_seam_test.go` (§2.2's VW4 constants), and the split of the design's single
  `surface_verb_why_test.go` into two files (V2b-D1). **Not created**: `host/surfacepage/testdata/verb-why.golden.json`
  (V2b-D7), a separate XV2 file (XV2 is in `verb_why_e2e_test.go`).
- **Not touched**: `host/why/*` (V2a's seam is used as is), `host/daemon/daemon.go` (no `Daemon` field: V2b-D3), every
  `/v1` handler, `readStore`, `approvalRoute`, the CSP, `cmd/ailang-worldd/testdata/passkey_cdp.mjs` (row 201's driver:
  V2b-D13), row 201's approval-card files, `world/*.ail`, `packages/**`, `scripts/**`, `tools/launchd/*`,
  `cmd/world-publish/*`, any skill, `~/.ssh`, `~/.ailang/world`, `~/.ailang/state/mission-v1*`.

---

## 0. Binding constraints (each is restated where it governs)

1. **Test stores, temp dirs and test keys only (D-WORLD-109 = A).** Every daemon is `daemon.New`/`newHandlerDaemon` on a
   `t.TempDir()` store, behind `httptest` loopback or called in-process; the live-Chrome tests use the passkey driver's
   temp profile (`startPasskeyBrowser`). No real store, key, `~/.ssh`, `/dev/tty`, browser profile, mint or publish.
2. **Read-only (design §0.2, §3, §4 T1).** No new `/v1` route (`TestFrozenV1RouteTableMatchesMux` unedited and green), no
   POST, no store write, no process spawn from the daemon, no `.ail`. VS2 pins `surface_verb_why.go`'s imports and every
   `d.*` it reaches to exact sets; VR4/VR6 prove the private mux serves exactly the seven GET handlers `Handler()`
   serves; VR8 proves none of them reads middleware state. The bar never runs `why -`, `--result` or `--expect-store`.
3. **The CLI is behaviour-neutral (design RV3).** `cmd/ailang-worldd/testdata/why-characterise.golden.jsonl` (V2a's VC1,
   generated at V2a's base) passes **unedited**, as do `why_test.go` and every caller test; the real-stack why tests
   (`TestDemoSeedWhyNamedLinks/real-stack`, `TestCLIAgainstSeToolsDaemon`) pass with `WORLD_TOOL_AILANG_BIN` set (M4,
   M6e acceptance; §6 `realstack`).
4. **Pure core, effects at the boundary (S1–S3).** `host/verbs` stays in VS1's pure import set; `host/why` is untouched
   (VS2a); the getter's effects (the in-process dispatch) live in `host/daemon`.
5. **One mutant per new or changed production branch, anchored to the diff** (§4); the named test and the behavioural
   suites with every static test `-skip`ped; every run appends a ledger line; the table's ids are the ledger's.
6. **No `t.Parallel()`** in any new test (`TestNoParallelWriteForkPackagesOutsideVerifygate`).
7. **No git writes by the executor** (§7). The controller commits from snapshots.

---

## 1. Pre-flight

### 1.1 Baseline at `c09e0cdb` (measured by the planner, 2026-11-17)

| Leg (`tools/gates.sh <base-export> <logs> base`, §6's script in base mode) | rc | Time |
|---|---|---|
| `status` | **0** | 0s |
| `verify_ail` | **0** | 37s |
| `gofmt` | **0** | 0s |
| `build` | **0** | 4s |
| `vet` | **0** | 3s |
| `fence` | **0** | 25s |
| `pkgs` | **0** | 285s |
| `race` | **0** | 32s |
| `inv` | **0** | 16s |
| named PASS: 55 | | |
| `realstack` | **0** | 12s |
| `chrome` | **0** | 69s |
| `verify_go` | **1** | 481s |
| `GATES_DONE 964s FAILS: verify_go` | | |

Every command ran with `export AILANG_BIN=$HOME/.pinned-ailang/ailang` (AILANG v0.41.0), `go1.26.6 darwin/arm64`, node
v26.13.0, in a `git archive c09e0cdb` export at `plan-scratch/base` (a throwaway git repo; `diff -r cmd host docs` against
the worktree: identical), load average 3–21 (the planner's stage, acceptance and mutant runs alongside).
Logs: `plan-scratch/logs/base_gates/`.

- `verify_ail`: `✓ verify gate PASSED: 16 required identities verified, 40 named tests pass, 503 se-tools named tests pass`.
- `pkgs` rc 0 (`host/verbs`, `host/why`, `host/surfacepage`, `host/daemon`, `cmd/ailang-worldd` whole); `race` rc 0;
  `inv` rc 0 (55 named `--- PASS`, none missing); **`realstack` rc 0**: `--- PASS: TestDemoSeedWhyNamedLinks/real-stack
  (5.93s)`, `--- PASS: TestCLIAgainstSeToolsDaemon (5.22s)` (with `WORLD_TOOL_AILANG_BIN=$HOME/.pinned-ailang-tools/
  v0.54.15/ailang`; without it CI skips them); **`chrome` rc 0**: 22 `--- PASS` (the twelve passkey/approve/decide/Needs-you
  page tests, XV1 `TestVerbBarEpisodeReportInChrome`, the two `host/daemon` app-shell drills, and their subtests).
- `verify_go` rc 1: **60** packages `ok`; red only in untouched packages, each a load timing test: `cmd/ailang-worldd`
  `TestHostsLaunchersExecHost` (the directive's list), `host/broker` `TestPublishDenialsPersistAndNeverDispatchWithLivePositiveControl`,
  `host/daemon` `TestServeStartsWhenTheProbeSpendsTheStartupClock` (red at V1's and V2a's bases too) and
  `TestRegistryPhaseHasItsOwnClock`, `host/verifygate` `TestMissionConfigDefaultFlowIntegration` (the directive's list).

**Rule (carried):** a red at the base measures the repo, not the sprint. The controller re-runs a red untouched
package alone once and does not chase it. **Never start a gate under `nohup`.**

### 1.2 Codebase claims checked (at `c09e0cdb` unless "work" — the planner's final stage tree `work2` at `M6e`)

| # | Claim | Command → observed output |
|---|---|---|
| C1 | The base and the tools | `git rev-parse HEAD` (worktree) → `c09e0cdb3c8a3eb25c208dbacedddfb690fe8ca0`; `git status --porcelain \| wc -l` → `0` (before the two plan files); `$HOME/.pinned-ailang/ailang --version` → `AILANG v0.41.0`; `go version` → `go1.26.6 darwin/arm64`; `/usr/bin/patch --version` → `patch 2.0-12u11-Apple`; `node --version` → `v26.13.0` |
| C2 | **No code drift since V2a's merge** (negative, with control) | `git diff --stat 2d530cbe c09e0cdb -- cmd host docs \| tail -1` → *(empty)*; control `… -- design_docs \| tail -1` → `18 files changed, 833 insertions(+), 131 deletions(-)` |
| C3 | **The design's M4–M6 line cites, re-measured** (all unchanged: V2a touched none of these files) | `grep -n` → `surface_verb.go:45` `defer func() { writeApprovalJSON(w, http.StatusOK, a) }()`, `:52` `a.Argv, a.Verb, a.Class, a.CLI, a.Reason = argv, v.Words, string(v.Class), v.CLI, v.Reason`, `:56` `case v.Words == "episode report":`, `:63` `a.Outcome = "terminal"`; `lookup.go:99-100` (the `why` row); `app.js:485` `const VERB_NEEDS = { address: [], result: ['cli'], …`, `:495` `if (b.outcome === 'refused') return { say: b.reason };`, `:499-500` (the `exit 1` rule and body), `:523` `if (verbJob === job) verbJob = null;`, `:995` `say` (2,600 ms), `:1267-1268` (`go`: empty → `out()` before `askVerb`), `:1283` `await loadIndex()`, `:1508` Enter, `:1520` Escape; `handlers.go:287` `readCtx`, `:349` `writeReadTimeout`; `daemon.go:144` `readDeadline = 10 * time.Second`, `:1110-1119` the seven GET registrations (and `:1109` `GET /v1/session/self`, `:1114` `by-semantic-id`, `:1117` `registry`, `:1120` episode report, which the walk does not read); `host/why/why.go:35` `ScanMax = 5000`, `:36` `LogPageMax = 500`, `:38` `MaxResponseBytes = 8388608`; `split.go:17-18` `MaxLine = 1024`, `MaxWords = 32` |
| C4 | **Audit input 5: the stale reason** | `grep -n 'case first == "why"' -A1 host/verbs/lookup.go` → `99: case first == "why":` `100: return Verb{Words: "why", Class: Read, CLI: cli, Reason: "why reads World, but its walk runs in the CLI. The app cannot run it yet."}, argv[1:]`; `surface_verb.go:52` copies `v.Reason` to `a.Reason`. Work: the row is `Verb{Words: "why", Class: Read, CLI: cli}` (V2b-D5); VR1 asserts `Reason == ""` on every `result` row; mutant `L-REASON` (the old text restored) is killed by VR1 (§4) |
| C5 | **Audit input 4: VS2's `d.*` set** | work: `TestSurfaceVerbWhyIsReadOnly` (`ast.Inspect` + `DeepEqual`, VS1's shape) → PASS with exactly `{handleHealth, handleHead, handleWorld, handleObject, handleLogEntry, handleLogRange, handleReceipt, whyRoutes, whyMux, whyGetter, readCtx, readDeadline, writeInternalErrorLog}`; no `sync.Once` field exists (V2b-D3). VS1 (`TestSurfaceVerbIsReadOnly`) gains only `verbWhy` for `surface_verb.go`. Its firing: the first V2b draft's VS1 run → `surface_verb.go reaches d.map[… verbWhy:true …], want exactly map[…]` (FAIL, before the set was widened) |
| C6 | **Audit input 6: the cap on raw bytes vs the encoded body** | work: `TestVerbAnswerCap` writes a 1,048,576-byte `\x01` text through `writeApprovalJSON` and asserts `1048576 < len(body) ≤ 6×1048576 + 1024` → PASS (each control byte encodes as `\u0001`; `writeJSON` is `json.NewEncoder`, HTML-escaping `<>&` to six bytes too) |
| C7 | **The report seed answers every §2.7 route row** (VR1's fixture; a probe test, removed) | work: on `verbRig` + `s.call("rest", "w1", "3", …)` + `s.pending("rest", "w1", "9", …)`: `why 1` → `result`, text `why: entry 1 (matched by entry index)\n  ✗ world …` (the seed's embedded genesis makes the world links ✗: **exit 3** end to end); `why head --json` → `json` `{\n  "entry": 3,\n  "matchedBy": "head",…`; `why rest:w1:3` → `why: entry 3 (matched by rest receipt → world)` (**RS-REST's resolved arm**); `why rest:w1:9` → stderr `ailang-worldd why: receipt rest:w1:9 is indeterminate: nothing committed\n`, exit 1; `why rest:nobody` → `… is not-started: …`; `why a2a:nobody:x --scan 1` → `… not found in the last 1 entries (invocationId)\n`; `why 99` → `… log entry 99 not found\n` |
| C8 | **The route golden is deterministic** (VP9's fixture, V2b-D7) | work: `WORLD_UPDATE_VERB_GOLDEN=1 go test ./host/daemon/ -run '^TestSurfaceVerbOutcomes$'` then two fresh runs without it → `ok` ×2; `why 2`'s text contains `<img src=x onerror=alert(1)>` (the seed's hostile skill name) |
| C9 | **VC1's golden carries the CLI's pre-move bytes VE1 needs** (V2a judge note 2) | `cut -c1-60 cmd/ailang-worldd/testdata/why-characterise.golden.jsonl` → rows `no target`, `malformed sha256`, `absent sha256`, `flag before target` (`--json 1`), `two targets`, `rest: unknown receipt`, `--expect-store this store` (the `1` chain), `index past head`; VE1 rebuilds VC1's rig exactly (`newWalkRig` + `invoke("ep1", "ailang-read", "task-1", "hello\n")`) and compares the route with those lines before any other commit |
| C10 | **Every sub-milestone is committed and measured** | `tools/stages.py` (commits each stage in `work2`, measures `git diff -M --numstat BASE <tag> -- cmd host docs`): `NUMSTAT M4 255`, `M6p 382`, `M5g 781`, `M5r 1106`, `M6e 1255` (`-- cmd host`: 255, 382, 781, 1106, 1251); `WORK2_EQUALS_WORK` |
| C11 | **Every sub-milestone is green when committed** | `tools/accept_all.sh` (an APFS clone of `work2` checked out at each tag, §3's `accept.sh`; `logs/accept_<M>.log`): first run, on the stages as first shipped: `M4 rc=0 ACCEPT_OK M4`; `M6p rc=0 ACCEPT_OK M6p`; `M5g rc=0 ACCEPT_OK M5g`; `M5r rc=0 ACCEPT_OK M5r`; `M6e rc=0 ACCEPT_OK M6e`. After VR2's combined row changed (§4 R2; M5r's test file only), on the final stages: `M5r rc=1 ACCEPT_RED M5r: x ^(TestVerbBarEpisodeReportInChrome)$` — XV1's first `waitFor` (the `episode report` panel) timed out at **load average 28.5** while the planner's final-gate `pkgs` leg ran beside it; the same test alone at `M5r`: `--- PASS: TestVerbBarEpisodeReportInChrome (1.60s)`; then `M6e rc=0 ACCEPT_OK M6e`; `M5r rc=0 ACCEPT_OK M5r` (re-run whole, load 8). **`ACCEPT_OK` for all five on the final stages**; each arm ran in ≈ 1–6 minutes (the five together ≈ 11 minutes at load 6–25), inside one 600 s call (§8 R6 for the Chrome legs under load) |
| C12 | **The patches reproduce the stage tree** | `tools/roundtrip.sh <this plan>`: a fresh `git archive c09e0cdb` export (from the worktree) + `patch -p1 -E` (after `--dry-run`) of Appendix P's five patches **extracted from this file with §2's `awk`**, in order (each `cmp`-equal to the generated patch; no `.orig`/`.rej`), each committed and measured → `CMP_EQUAL M4`; `NUMSTAT M4 255`; `CMP_EQUAL M6p`; `NUMSTAT M6p 382`; `CMP_EQUAL M5g`; `NUMSTAT M5g 781`; `CMP_EQUAL M5r`; `NUMSTAT M5r 1106`; `CMP_EQUAL M6e`; `NUMSTAT M6e 1255`; `ROUNDTRIP_IDENTICAL_TO_STAGE_TREE` |
| C13 | **Every mutant; the table's ids are the ledger's** | 99 rows (3 `CTL`, 92 `K`, 1 `N`, 2 `S`, 1 `X`); 110 ledger lines; 99 of 99 rows run; the last line of each row behaves as expected for **99 of 99**. Table ids = ledger ids (`tools/assemble.py` fills §4 from the ledger's last line per id) |
| C14 | **The gates, on the final tree** | `tools/gates.sh <final-tree> logs/final_gates final` (`final-tree` = `work2` at `M6e`, a git repo with every stage committed, byte-identical to the round-trip tree; load 8–29, the planner's mutant and acceptance runs beside it): `status` 0; `verify_ail` rc 0 (no `.ail` changed; `✓ verify gate PASSED: 16 required identities verified, 40 named tests pass, 503 se-tools named tests pass`, as the base); `gofmt`, `build`, **`vet`**, `fence` rc 0; **`pkgs` rc 1** — `ok` `host/verbs`, `host/why`, `host/surfacepage`, `cmd/ailang-worldd` (172.8 s); red only `host/daemon`'s `TestServeStartsWhenTheProbeSpendsTheStartupClock` and `TestRegistryPhaseHasItsOwnClock` (startup-clock timing tests, untouched, **both red at the base's `verify_go` too**, §1.1) — **green re-run alone** (`--- PASS … (1.94s)`, `--- PASS … (9.42s)`, `logs/final_daemon_reds_alone.txt`) and green in this run's own `verify_go` (`ok host/daemon 378.1s`); **`race` rc 0**; **`inv` rc 0: 71 named `--- PASS`**, 0 SKIP, 0 FAIL, none missing (the base's 55 + the 16 new: `TestWhyFlags`, `TestAdmitWhy`, `TestNeverAndOtherReasonsEndWithTheCLILine`, `TestSplitBoundsAreLiterals`, the ten new `host/daemon` tests, `TestWhyBarAdmissionIsTheCLIs`, `TestVerbRouteWhyIsTheCLIs`) — incl. `TestDaemonDependencyAllowlist`, `TestFrozenV1RouteTableMatchesMux`, `TestWorldBoundaryDependencyAllowlist`, `TestProductionContextRoots`, `TestProductionGoSurface`, `TestEveryVerbHasHelp`, `TestSurfaceAppPure`, VS1, VS2, `TestSurfaceVerbOutcomes` (the golden); **`realstack` rc 0** (`--- PASS: TestDemoSeedWhyNamedLinks/real-stack (5.07s)`, `--- PASS: TestCLIAgainstSeToolsDaemon (4.90s)`); **`chrome` rc 0: 23 `--- PASS`** (incl. `TestVerbBarEpisodeReportInChrome (1.19s)`, `TestVerbBarWhyInChrome (0.98s)`); **`verify_go` rc 0: 64 packages `ok`** (incl. `host/daemon`, `cmd/ailang-worldd 236.0s`, `host/verbs`, `host/why`, `host/surfacepage`). **`GATES_DONE 1054s FAILS: pkgs`** — the one red two untouched load-timing tests, green alone; `ALL_DONE.<sha>` written (`logs/final_gates/`) |
| C15 | **Exactly these pre-existing test files change** | work: `git diff --name-status BASE M6e -- '*_test.go'` → `M cmd/ailang-worldd/quickstart_verbs_test.go`, `M cmd/ailang-worldd/verb_bar_chrome_test.go`, `M cmd/ailang-worldd/why_seam_test.go`, `M host/daemon/surface_verb_outcomes_test.go`, `M host/daemon/surface_verb_test.go`, `M host/surfacepage/app_test.go`, `M host/surfacepage/app_verb_test.go`, `M host/verbs/lookup_test.go` (+ 5 `A`). Each reason: §9 criterion 4 |
| C16 | **No new `/v1` route; `daemon.go` untouched** | work: `git diff --name-only BASE M6e -- host/daemon/daemon.go host/daemon/handlers.go` → *(empty)*; `TestFrozenV1RouteTableMatchesMux`, `TestSurfaceVerbRouteRegistered` → PASS unedited (C14's `inv`) |
| C17 | **The snapshot line works on a sub-milestone with no `docs/`** | `tools/snap_snippet.sh` (§7's block) on a clone at `M4` (`cmd/`, `host/` only): `STATUS` 5 lines, `DELETED: 0`, `MANIFEST: 5`, no error |
| C18 | **Every name the patches use exists at HEAD or is new** (396 retro) | at HEAD: `verbs.Quote` `split.go:100`, `verbs.Binary` `lookup.go:32`, `verbs.Split` `split.go:35`, `verbs.Lookup` `lookup.go:67`, `verbs.ErrTooLong` `split.go:23`, `why.Run` `why.go:330`, `why.Reader` `:100`, `why.ParseTarget` `:295`, `why.Answer` `:323`, `why.Options` `:317`, `why.Getter` `:97`, `why.MaxResponseBytes` `:38`, `writeApprovalJSON` `surface_approve.go:32`, `verbAnswer` `surface_verb.go:18`, `d.writeInternalErrorLog` `workbench.go:136`, `d.readCtx` `handlers.go:287`, `d.readDeadline` `daemon.go:472`, `readDeadline` `daemon.go:144`, `d.handleHealth` `daemon.go:1173` (and the other six, `daemon.go:1110-1119`), `d.handleSessionSelf` `session_self.go:46`, `internalErrorMessage` `handlers.go:150`, `readStore` `daemon.go:557`; tests: `newHandlerDaemon` `handlers_test.go:21`, `verbRig`/`verbGet`/`verbAnswerT` `surface_verb_test.go:28,36`, `approvalRig.serve` `surface_approve_rig_test.go:108`, `apLoopback` `:31`, `expiredReadDeadline` `read_deadline_test.go:35`, `aguiFailReads` `agui_test.go:626`, `newReportSeed`/`call`/`pending`/`read` `episode_report_test.go:56,125,149,157`, `boundedTestContext` (package helper); cmd: `newWalkRig` `why_test.go:52`, `whyRun` `:208`, `invoke`/`invokeAs` `:109,116`, `synthInv` `:43`, `contains` `why_cli.go`, `whyGolden` `why_characterise_test.go:15`, `startPasskeyBrowser` `passkey_page_driver_test.go:63`, `exitUsage`/`exitOK` `main.go:74-75`; page: `askVerb`, `verbAct`, `verbOutcome`, `verbPanel`, `addressTarget`, `SEAMS.parseAddress`, `loadIndex`, `openState`, `out`, `say`, `showVerbPanel` (`app.js`). **New**: `verbs.WhyScanDefault`, `WhyScanMax`, `WhyArgs`, `WhyFlags`, `WhyLeadingTarget`, `WhyTrailingTarget`, `WhyScanOK`, `AdmitWhy` (+ unexported `whyBarRefusal`); daemon `MaxVerbAnswerBytes`, `verbAnswerCap`, `capVerbAnswer`, `whyRoutes`, `whyMux`, `whyTally`, `boundedRecorder`, `whyGetter`, `verbWhy`, `whyResult`; page `dropVerb`, `barTurn`, `goWith`, `addressGo`; test helpers `whyRig`, `whyCLI`, `whyFailReads`, `recordingGetter`, `healthPin`, `methodName`, `mapKeys`, `whySorted`, `verbRoute` |

### 1.3 Planner decisions where the design leaves a choice

- **V2b-D1 — Five sub-milestones, page before route.** **M4** (design M4 + §2.2), **M6p** (design M6's page half: A2–A4,
  VP7, VP8, VP10, VP11; the route still answers `why` with `terminal`, so the page change is proven on V1's real
  answers first), **M5g** (design M5's getter and A2 cap, unrouted: VR5, VR7, VR8, VS2, `whyResult`), **M5r** (the route's
  `why` branch: Lookup's row, VR1–VR4, VR6, VR7 end to end, the golden's `why` rows, VP9), **M6e** (VA3, VE1/VE3, XV2,
  QUICKSTART §19). Order chosen so that each commit is green **including live Chrome**: XV1's `log` arm changes with
  A3 (M6p) and its `why` arm with A1 (M5r) — two small edits, never a red commit. The design's single
  `surface_verb_why_test.go` is two files (`…_why_test.go` M5g, `…_why_route_test.go` M5r) so M5g lands its tests.
- **V2b-D2 — Design §2.2 taken from V2a's probe, re-measured.** `v2b-probe-2.2.patch` applied to `c09e0cdb` cleanly (an
  offset in `why_seam_test.go`); its `flags.go` half is **verbatim** (46 lines), its `why_cli.go` repoint verbatim (+9
  −16), VW4's two constants verbatim (+3 −1); VA1 (`TestWhyFlags`) moves from `flags_test.go` to the new
  `why_admit_test.go` (no pre-existing `host/verbs` test edited for it) and adds the two constants as literals. §2.2
  measures ≈ 120 here (probe 111). VC1, `why_test.go` and the real-stack leg pass unedited after the repoint (M4).
- **V2b-D3 — The private mux is built per walk, not once (audit input 4).** `d.whyMux()` builds a fresh
  `http.ServeMux` from `d.whyRoutes()` per `verbWhy` call (seven registrations, microseconds; the walk then makes up to
  thousands of reads through it). Reason: no `Daemon` field, so `daemon.go` is untouched and there is no `sync.Once`
  state to pin or race; VS2's set is therefore the seven handler method values + `whyRoutes`, `whyMux`, `whyGetter`,
  `readCtx`, `readDeadline`, `writeInternalErrorLog` (C5). **T1's read-only argument rests on VS2 (static: nothing
  else is reachable) and on VR4/VR6/VR8 (behavioural: the mux serves exactly the seven GET handlers `Handler()` serves,
  each byte-identical over HTTP, none reading middleware state).**
- **V2b-D4 — `whyRoutes()` is the one table.** A method returning `map[pattern]http.HandlerFunc`; `whyMux` registers from
  it; VR6 compares its keys with the patterns the recorded walks reach (generated); VS2 binds each pattern to the method
  name `Handler()` registers for it in `daemon.go`'s source (static, labelled).
- **V2b-D5 — Lookup's `why` row loses its reason (audit input 5).** `Verb{Words: "why", Class: Read, CLI: cli}`, as
  `episode report`'s row; the route copies `""`; every refusal sets its own reason. The stale text can therefore ride on
  no answer; VR1 asserts `Reason == ""` on all nine `result` rows and `TestLookup`'s two `why` rows assert `""`.
- **V2b-D6 — The A2 cap (audit input 6).** `capVerbAnswer(&a)` runs first in `handleSurfaceVerb`'s deferred writer for
  every outcome (it acts only on `result`), measuring `len(Text)+len(JSON)+len(Stderr)` raw bytes against
  `verbAnswerCap` (= `MaxVerbAnswerBytes` = 1,048,576; a package variable so VR7's end-to-end test can lower it).
  **RV2, amended:** the cap is on raw bytes; the JSON-encoded body can be larger — a control byte or `<`, `>`, `&`
  escapes to six bytes (`\u00XX`) — and is still bounded by 6 × 1 MiB plus the answer's other members (C6). The
  comment on `MaxVerbAnswerBytes` says the same.
- **V2b-D7 — The page's `why` rows live in V1's route golden** (`verb-route.golden.json`, generated by
  `TestSurfaceVerbOutcomes` under `WORLD_UPDATE_VERB_GOLDEN=1`), not in a new `verb-why.golden.json` (design VE3).
  Reason: one generator (V1's), real answers from the route on `verbRig` — a real **exit 3** chain (`why 2`, whose skill
  is the seed's hostile `<img …>` name, VP9's hostile row), a real `--json` answer (`why head --json`), a real exit-1
  stderr (`why 9`) and a real `--result` refusal — deterministic (C8). VE1 proves the route's bytes are the CLI's on the
  walk rig (exit 0 chains), so the page replays route bytes that are CLI bytes.
- **V2b-D8 — VE2 (exit-3 parity) is carried by VR1 + VR4.** The walk rig makes ✗ only through `tamperProxy`, which sits
  in front of HTTP and cannot reach the in-process getter (§1.4 P7). The report seed's ✗ chains give: VR1 — `exit 3`
  with the CLI printer's literal `… link(s) BROKEN — integrity refusal (exit 3)\n` end to end; VR4 — `why.Run` over an
  HTTP getter against the same daemon's `Handler()` equals the in-process answer for the ✗ targets (`1`, `head --json`,
  `rest:w1:3`, the world ref). Mutant `D-RES-EXITSET` (exit dropped) is killed by both.
- **V2b-D9 — The V2a judge's notes (directive input 4).** (1) the six unpinned base lines: `--scan 0` → VA2 + VA3 +
  VR1 (refusal) and the CLI's own `TestWhyScanBound`; `why head` with no head, the head-absent error, and `ScanFor`'s "no
  head" stay residuals (a store whose head names no world cannot be built through `POST /v1/commit`, V2a §3.1); the
  `rest:` receipt state operand → VR1 (`indeterminate`, `not-started`) and VE1/VC1 (`not-started`); the resolved
  `rest:` arm (RS-REST) → VR1 + VR4 (`rest:w1:3` → `matched by rest receipt → world`) — **RS-REST closed**; the
  fallback-plan operand (`- not a result`) is the CLI's `-`, which the bar refuses — VC1 keeps it; a mid-rune `Elide`
  cut stays as V2a's `TestElideBounds` covers it (host/why untouched). (2) the bar's `--json` answer is pinned by VE1
  against VC1's `flag before target` line (generated at V2a's base, not by `why.Run`), trailing newline included, plus
  a literal `\n}\n` suffix check; `--result`'s error text is not reachable from the bar (refused). (4) VW5 is unchanged
  by V2b (`host/why` and `mcpclient.go` untouched): no action.
- **V2b-D10 — Precedence (a) is tested.** The design says nobody reads (a); but without it a cancelled walk falls through
  to `whyResult` and answers `ailang-worldd why: context canceled` as a `result`. VR2 drives a cancelled request
  through `Handler()` and asserts `the request ended before the walk did` (mutant `D-PREC-A`).
- **V2b-D11 — `whyResult` is a function** (design §2.5 d and f) so (d), unreachable through the real `why.Run`, is
  unit-tested over crafted `why.Answer`s (`TestWhyResult`, ten rows).
- **V2b-D12 — `addressGo(raw, current, fx)`** holds `goAddress`'s body above the page's `document` guard, so A4's
  post-`loadIndex` check is node-testable; `goAddress` is a one-line wrapper with the real effects (§1.4 P6).
- **V2b-D13 — XV2's Escape is a page-dispatched `KeyboardEvent`** on `document` (the app's own keydown listener runs);
  the passkey driver's `KEYS` table (`passkey_cdp.mjs:285`, row 201's file) has no Escape and is not edited.
- **V2b-D14 — New tests in new files** wherever the design allows (`why_admit_test.go`, `bounds_test.go`, the two daemon
  files, `verb_why_e2e_test.go`); a pre-existing test file changes only where V2b's behaviour change demands it (C15).
- **V2b-D15 — VP10 on episode-report rows** (route-real since V1) lands with the page in M6p; VP9's `why` rows land with
  the route in M5r.
- **V2b-D16 — QUICKSTART §19** gains `why head` and `why 3 --json` in its verbatim block and one paragraph (the 10 s note,
  the refusals, A3); `TestQuickstartSection19MatchesTheBar` learns `why` (Read) and counts words instead of classes
  (two `Read` words).
- **V2b-D17 — No cut.** 1,255 ≤ 1,400. Had it been needed, the next cut is M6e's VE1 second loop (the generated
  CLI-vs-route sweep, ≈ 25 lines), never VR1–VR5, XV2 or VP7–VP11.

### 1.4 Premises found false or incomplete at HEAD (each routed around minimally)

- **P1 — Design §7 "Pre-existing tests edited: `surface_verb_test.go` VS1 call set; `surface_verb_outcomes_test.go`;
  `lookup_test.go`; `app_verb_test.go`. Nothing else."** Also needed: `surface_verb_test.go:104` (the 1,024-byte
  `why xxx…` line expected `terminal`; now a `ParseTarget` refusal — the row keeps its point, "1,024 bytes is not
  `TooLong`"), `host/surfacepage/app_test.go:176` (`TestSurfaceAppPure` pins `module.exports`; `goWith`, `addressGo`
  join), `cmd/ailang-worldd/verb_bar_chrome_test.go` (XV1 asserted `why` → terminal and `log` → header say; A1 and A3
  change both), `cmd/ailang-worldd/quickstart_verbs_test.go` (§19's word set), `cmd/ailang-worldd/why_seam_test.go`
  (§2.2's two constants in VW4, from the V2a probe). Routed: each is a minimal edit (§9 criterion 4).
- **P2 — Design VR2 "a `d.reads` wrapper that blocks past the deadline".** With an expired `readDeadline` the walk's
  `GetJSON` returns the context error before any read (V2a's `W-CTX` branch), so no read starts — measured: 0
  `GetObject` calls (VR2 asserts it). Routed: VR2 = `expiredReadDeadline` + a wrapper answering
  `context.DeadlineExceeded` (the handlers' 503) with the walk's context alive; no blocking wrapper.
- **P3 — Design §2.7 "`why rest:x` with a pending receipt → `receipt rest:x is pending`".** Measured: an intent that
  never commits is `indeterminate`; an unknown id `not-started` (C7). Routed: VR1 asserts the measured words (the
  design marks §2.7's texts read-not-run).
- **P4 — Design VR8's control "`approvalRoute`'s `LocalAddrContextKey` read".** The read is in `approvalOrigins`
  (`surface_approve.go:60`); `approvalRoute` (`:76`) only calls it. Routed: VR8's control names `approvalOrigins`.
- **P5 — Design §6 "V1 page rows re-pointed … use `episode launch w1 --preset p`".** Not in the golden; re-pointed to
  `episode finish w1 --preset p` (a golden row, `terminal`).
- **P6 — Design §2.6 "`goAddress(raw, current)` re-checks `current()`" as a node-tested change (VP7).** `goAddress` sits
  below `if (typeof document === 'undefined') return;` (`app.js:549`), unreachable by the node harness. Routed:
  V2b-D12.
- **P7 — Design VE2 "a `makeCLICommit`-built entry the commit route accepts" with a ✗.** The walk rig's ✗ tests use
  `tamperProxy` (HTTP only). Routed: V2b-D8.
- **P8 — Design XV2 "then Escape".** The CDP driver cannot press Escape (P/D13). Routed: V2b-D13.
- **P9 — Design §8 "V2b ≈ 930" (+ §2.2's 111 ≈ 1,041).** Measured 1,255: production 343 is within the design's
  estimate; tests are 908 (the daemon's two files 509) because every §2.5 precedence row, every getter class and every
  §2.3 form has its own behavioural row (directive retro 3–4). Within the cap; no cut (V2b-D17).
- **P10 — Audit input 4 "plus any `sync.Once` field".** None exists (V2b-D3).
- **P11 — Design §2.5 (a) "nobody reads it".** True on the wire, but untested it is a silent fall-through (V2b-D10).

### 1.5 The size

Lines are insertions + deletions: **`git diff -M --numstat c09e0cdb HEAD -- cmd host docs`** (git's rename detection on;
V2b renames nothing). **Every row is measured** from the scratch stages, each committed in `work2`, and re-measured on
the round-trip tree (C10, C12).

| Sub-milestone (in order) | Production | Test | Docs | Own | Cumulative vs `c09e0cdb` (after its commit) | `-- cmd host` | Run |
|---|---|---|---|---|---|---|---|
| M4 — §2.2 + `AdmitWhy`; VA1, VA2, VB6, VB7 | 110 | 145 | 0 | 255 | **255** | 255 | 1 |
| M6p — the page: A2–A4; VP7, VP8, VP10, VP11; XV1's `log` arm | 60 | 67 | 0 | 127 | **382** | 382 | 1 |
| M5g — the getter, `verbWhy`, `whyResult`, the A2 cap; VR5, VR7, VR8, VS2 | 165 | 234 | 0 | 399 | **781** | 781 | 1 |
| M5r — the route's `why`; Lookup's row; VR1–VR4, VR6, VR7 e2e; golden; VP9; XV1's `why` arm | 8 | 317 | 0 | 325 | **1,106** | 1,106 | 2 |
| M6e — VA3, VE1/VE3, XV2; QUICKSTART §19 | 0 | 145 | 4 | 149 | **1,255** | 1,251 | 2 |
| **The sprint** | **343** | **908** | **4** | | **1,255 ≤ 1,400** (margin 145) | 1,251 | |

(Own lines of the files two stages touch — `app_verb_test.go`, `verb_bar_chrome_test.go`, `surface_verb.go` — are
counted where they land; the cumulative column is authoritative.) Per file, cumulative (+/−): `surface_verb_why.go`
+163; `surface_verb_why_route_test.go` +285; `surface_verb_why_test.go` +234; `verb_why_e2e_test.go` +139;
`why_admit_test.go` +103; `flags.go` +85; `app_verb_test.go` +53 −6; `bounds_test.go` +38; `app.js` +35 −23;
`verb_bar_chrome_test.go` +11 −6; `why_cli.go` +9 −16; `golden` +5 −1; `surface_verb.go` +4 −2; `why_seam_test.go`
+3 −1; `QUICKSTART.md` +3 −1; `quickstart_verbs_test.go` +3 −3; `surface_verb_test.go`, `surface_verb_outcomes_test.go`,
`lookup.go`, `lookup_test.go` +2 −2 each; `index.html`, `app_test.go` +1 −1 each.

---

## 2. How the executor applies a sub-milestone (the patch protocol)

Each sub-milestone `<M>` is the block fenced by a line `` ````diff <M> `` (four backticks) and a line of four backticks,
in Appendix P. Paths are `a/<repo path>` and `b/<repo path>`; a new file's patch has `--- /dev/null` and an
`@@ -0,0 …` hunk. V2b deletes no file; `-E` is kept (the protocol is V2a's, and it is harmless here). Apply it **only**
this way, in one `bash` shell (the controller's shell is zsh):

```bash
bash <<'EOF2'
set -euo pipefail
M=M4   # <- the sub-milestone, exactly as in the Appendix P heading, in order: M4 M6p M5g M5r M6e
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter398-row243-v2b
PLAN="$W/design_docs/planned/w-world-address-bar-verbs-v2b-sprint-plan.md"
T="$HOME/.ailang/state/world-iter398/exec-patches"; mkdir -p "$T"
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
**[base]** lines (gofmt, `go build ./...`, `go vet` of the five touched packages, the `_test.go` compile fence). Each
`v`/`x`/`rs` line must show `--- PASS` for every test it names and **no `--- FAIL`, no `--- SKIP`** (`x` is live Chrome
with `WORLD_REQUIRE_PASSKEY_BROWSER=1`, so a missing browser is a FAIL, not a skip; `rs` is the real-stack why leg with
`WORLD_TOOL_AILANG_BIN`); each `p` line must end `ok`. The script ends `ACCEPT_OK <M>` or prints `ACCEPT_RED <M>:
<step>` and exits 1. Each arm fits a 600 s call (measured: the planner's runs, C11). The planner's run of it on every
stage is §1.2 C11. Extract and run it (`M` and `W` as in §2):

```bash
bash <<'EOF2'
set -uo pipefail
M=M4
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter398-row243-v2b
PLAN="$W/design_docs/planned/w-world-address-bar-verbs-v2b-sprint-plan.md"
L="$HOME/.ailang/state/world-iter398/exec-logs"; mkdir -p "$L"
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
export WORLD_EXEC_NODE=/opt/homebrew/bin/node
unset WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER WORLD_UPDATE_VERB_GOLDEN WORLD_UPDATE_WHY_GOLDEN WORLD_TOOL_AILANG_BIN
CH='/Applications/Google Chrome.app/Contents/MacOS/Google Chrome'
TOOL=$HOME/.pinned-ailang-tools/v0.54.15/ailang
set -o pipefail
chk() { tee /dev/stderr | { out=$(cat); echo "$out" | grep -qE -- '--- (FAIL|SKIP)|^FAIL' && exit 1; echo "$out" | grep -q -- '--- PASS' || exit 1; }; }
v() { go test -count=1 -v "$@" 2>&1 | chk; }   # -v: every named test PASS, none SKIP/FAIL
p() { go test -count=1 "$@"; }
x() { WORLD_CHROME="$CH" WORLD_REQUIRE_PASSKEY_BROWSER=1 go test -count=1 -v ./cmd/ailang-worldd/ -run "$1" 2>&1 | chk; }   # live Chrome, temp profile
rs() { WORLD_TOOL_AILANG_BIN=$TOOL go test -count=1 -v -timeout 20m ./cmd/ailang-worldd/ -run '^(TestDemoSeedWhyNamedLinks|TestCLIAgainstSeToolsDaemon)$' 2>&1 | chk; }
step() { echo "== $*"; "$@" || { echo "ACCEPT_RED $M: $*"; exit 1; }; }
# [base]
step test -z "$(gofmt -l cmd host)"
step go build ./...
step go vet ./host/verbs/ ./host/why/ ./host/daemon/ ./host/surfacepage/ ./cmd/ailang-worldd/
if go test -run '^$' ./... 2>&1 | grep -E '^(FAIL|--- FAIL)'; then echo "ACCEPT_RED $M: FENCE_RED"; exit 1; fi
WHYT='TestWhyCharacterisedPaths|TestWhyTargetFormsResolveToOneEntry|TestWhyResolvesMCPInvocationID|TestWhyBrokenLinkIsIntegrityRefusal|TestWhyScanBound|TestWhyShowsTheWorktreeDelta|TestWhyShowsTheFailureLine|TestWhyRunReplaysTheCLI|TestWhyRunGetterClasses|TestWhyConstantsMatchTheCLI|TestWhyHeadIsOneReader|TestWhyProtectedLine|TestEveryVerbHasHelp|TestVerbTableMatchesTheCLISwitch'
VERBS='TestWhyFlags|TestAdmitWhy|TestNeverAndOtherReasonsEndWithTheCLILine|TestSplitBoundsAreLiterals|TestLookup|TestSplit|TestQuoteRoundTrip|TestEpisodeReportFlags|TestParseInterleaved|TestApproveTarget'
PAGE='TestSurfaceAppVerbs|TestSurfaceAppPure|TestSurfaceAppMatchesWorkingCopy|TestSurfaceAppShellAssembly'
V1D='TestSurfaceVerbRouteRegistered|TestSurfaceVerbOutcomes|TestSurfaceVerbBoundsAndErrors|TestSurfaceVerbIsReadOnly|TestDaemonDependencyAllowlist|TestFrozenV1RouteTableMatchesMux|TestEpisodeReportRESTRoute'
GETD='TestWhyGetterUnroutedAndOversize|TestVerbAnswerCap|TestWhyHandlersReadNoMiddlewareState|TestSurfaceVerbWhyIsReadOnly'
ROUTED='TestSurfaceVerbWhyOutcomes|TestSurfaceVerbWhyDeadline|TestSurfaceVerbWhyReadFailure|TestWhyGetterTransportParity|TestVerbAnswerCapEndToEnd'
case $M in
M4)  step v ./host/verbs/ -run "^($VERBS)\$"
     step v ./host/daemon/ -run '^(TestSurfaceVerbIsReadOnly)$'
     step v ./cmd/ailang-worldd/ -run "^($WHYT)\$"
     step rs
     step p ./host/verbs/
     step p ./cmd/ailang-worldd/ -run '^(TestWhy|TestDemo|TestLogTail|TestProvenance|TestServicePointers|TestCLIWalk|TestVerb)' ;;
M6p) step v ./host/surfacepage/ -run "^($PAGE)\$"
     step v ./host/daemon/ -run "^($V1D)\$"
     step x '^(TestVerbBarEpisodeReportInChrome)$'
     step p ./host/surfacepage/ ;;
M5g) step v ./host/daemon/ -run "^($GETD|$V1D)\$"
     step v ./host/boundary/ -run '^(TestWorldBoundaryDependencyAllowlist|TestBareNetHTTPExemptionIsPerGroup|TestBoundaryASTWriteGuard)$'
     step v ./host/store/ -run '^(TestProductionContextRoots|TestProductionGoSurface)$'
     step p ./host/daemon/ -run '^(TestSurface|TestWhy|TestVerb|TestEpisodeReport|TestApprov|TestGET|TestFrozen|TestHandler)' ;;
M5r) step v ./host/daemon/ -run "^($ROUTED|$GETD|$V1D)\$"
     step v ./host/verbs/ -run "^($VERBS)\$"
     step v ./host/surfacepage/ -run "^($PAGE)\$"
     step x '^(TestVerbBarEpisodeReportInChrome)$'
     step v ./cmd/ailang-worldd/ -run '^(TestVerbRouteReportIsTheCLIs|TestQuickstartSection19MatchesTheBar)$'
     step p ./host/daemon/ -run '^(TestSurface|TestWhy|TestVerb|TestEpisodeReport|TestApprov|TestGET|TestFrozen|TestHandler)'
     step p ./host/surfacepage/ ./host/verbs/ ;;
M6e) step v ./cmd/ailang-worldd/ -run "^(TestWhyBarAdmissionIsTheCLIs|TestVerbRouteWhyIsTheCLIs|TestQuickstartSection19MatchesTheBar|TestVerbRouteReportIsTheCLIs|$WHYT)\$"
     step x '^(TestVerbBarEpisodeReportInChrome|TestVerbBarWhyInChrome)$'
     step rs
     step p ./cmd/ailang-worldd/ -run '^(TestWhy|TestVerb|TestQuickstart|TestDemo|TestLogTail|TestProvenance)' ;;
*) echo "unknown $M"; exit 2 ;;
esac
echo "ACCEPT_OK $M"
````

### M4 — admission: design §2.2 + `AdmitWhy` (design M4 without the route) · 110 prod + 145 test · Run 1

**Patch** `M4`: **`host/verbs/flags.go`** — `WhyScanDefault`/`WhyScanMax` (500/5000), `WhyArgs`, `WhyFlags` (the four
flags, the CLI's usage strings verbatim), `WhyLeadingTarget` (`runWhy`'s leading-target rule), `WhyTrailingTarget`
(`unexpected argument %q`), `WhyScanOK` (V2a's probe, verbatim); `whyBarRefusal` and **`AdmitWhy(args) (target, scan,
asJSON, reason)`**: the CLI's order (leading target → `fs.Parse` on a fresh set, output discarded → trailing target →
`--scan` bound) then the bar profile (`--result`, `--expect-store`, `-`, in that order, each `… is not run from the app:
<why>; at a terminal: <cli>` with `cli` = `Quote(["ailang-worldd","why",args…])`, Lookup's line); help → `see
ailang-worldd why --help at a terminal`; a flag error → the flag package's line (no prefix). **`why_cli.go`** — the CLI's
`runWhy` uses the five (V2a probe, verbatim; output unchanged, VC1). **`why_seam_test.go`** — VW4 gains
`verbs.WhyScan* == why.Scan*`. New **`host/verbs/why_admit_test.go`** — **VA1** `TestWhyFlags`; **VA2** `TestAdmitWhy`
(31 rows: every form of design §2.3 incl. `--scan 5000`/`5001`/`0`/`-1`, `-3`, `--bogus`, `--scan x`, `-h`, `--help`
before `--result`, the three profile refusals and their order, an empty `--result` as at a terminal; each refusal's CLI
line is `Lookup`'s). New **`host/verbs/bounds_test.go`** — **VB6** `TestNeverAndOtherReasonsEndWithTheCLILine` (24
words × bare/with an argument, `HasSuffix "; at a terminal: " + cli`), **VB7** `TestSplitBoundsAreLiterals` (32, 1,024
as literals; 33 words and 1,025 bytes refused).
**Expected cumulative numstat after the controller's commit of M4: 255.**

### M6p — the page: A2, A3, A4 (design M6, page half) · 60 prod + 67 test · Run 1

**Patch** `M6p`: **`app.js`** — `verbOutcome`: a `refused` with a non-empty `cli` → `{panel: {head: '$ ' + cli, body:
reason, copy: cli}}`, without → `say` (A3); a `result`: exactly one member, `exit` ∈ {absent, 1, 3}, `stderr ⇔ exit
1`, body = member + (`exit` ? `'exit ' + exit` : '') (A2); `dropVerb()` (askVerb's first line, extracted); `barTurn` +
**`goWith(raw, fx)`** (A4: an empty line drops a pending ask and goes home; an answer after a newer Enter is dropped;
the address path gets the turn's `current`); **`addressGo(raw, current, fx)`** (`goAddress`'s body, re-checking
`current()` after `await loadIndex()`, V2b-D12); `go` and `goAddress` become wrappers over them with the real effects;
`module.exports` += `goWith, addressGo`; the placeholder and the `?` keys row name `why`. **`index.html`** — the static
placeholder names `why <target>`. **`app_test.go`** — `TestSurfaceAppPure`'s export list += the two. **`app_verb_test.go`**
— **VP11** (refusals with `cli` → panel for `log`, `world get`, `episode`; without → `say` for `world a b c`,
`TooLong`; a `cli` of `''` is no `cli`), **VP10** (exit 3 under text; no `cli`, exit 3 with stderr, exit 2, exit 1 with
text, two members → unreadable, on episode-report rows), **VP8** (A pending, B supersedes, A's `finally`, C aborts
B), **VP7** (an empty Enter aborts the pending ask and goes home; a late `address` after a newer Enter does nothing;
the address path's `current` ends with a newer Enter; `addressGo` does not act after a newer Enter, and acts
without one), the superseded `log` row expects the panel. **`verb_bar_chrome_test.go`** (XV1) — `log` now opens the
panel with its line; `world a b c` is the header-`say` row (bar and place kept).
**Expected cumulative numstat after the controller's commit of M6p: 382.**

### M5g — the getter and the A2 cap (design M5, unrouted) · 165 prod + 234 test · Run 1

**Patch** `M5g`: new **`host/daemon/surface_verb_why.go`** — `MaxVerbAnswerBytes` (1 MiB; the raw-vs-encoded comment,
V2b-D6), `verbAnswerCap`, `capVerbAnswer`; `whyRoutes` (the seven patterns → `d.handle*`), `whyMux` (per walk,
V2b-D3); `whyTally`, `boundedRecorder` (status once; body capped at `max`, over → `over`); **`whyGetter(r, mux,
tally)`** (request build error, `pattern == ""` → `unrouted in-process read: <path>`, over the cap → `response from
<path> exceeds 8388608 bytes`: each an error, `tally.failed`, `d.writeInternalErrorLog`; 503 → `tally.timeout`; other
≥ 500 → `tally.failed`; else `(status, body, nil)`, 404 included); **`verbWhy`** (AdmitWhy → reason; `ParseTarget` →
`ailang-worldd why: <err>`; `d.readCtx(r)`; `why.Run`; precedence (a) cancelled → `the request ended before the walk
did`, (b) `DeadlineExceeded` or a 503 → `read deadline (10s) exceeded: the walk did not finish in the app; at a
terminal it has 60 s: <cli>`, (c) a failed read → `the daemon could not read World during the walk, so no chain is
shown; at a terminal: <cli>`, then `whyResult`); **`whyResult`** ((d) `the walk answered something the app cannot show;
at a terminal: <cli>` unless exactly stdout+0/3 or stderr+1; (f) text/json/stderr + exit). **`surface_verb.go`** — the
deferred writer runs `capVerbAnswer(&a)` first (episode report is capped from here on). New
**`host/daemon/surface_verb_why_test.go`** — **VR5** `TestWhyGetterUnroutedAndOversize` (five unrouted paths incl.
`Handler()`'s own `session/self`, `registry`, `by-semantic-id`, episode report → error, never 404, logged; an
unparseable path; a 404 is an answer; 8,388,608 bytes pass, 8,388,609 fail (literals); 503 → timeout, 500/502 →
failed, 404/400 → clean); **VR7** `TestVerbAnswerCap` (1,048,576 passes; +1 on text, json, stderr, or summed → refused
with the size and CLI line, members cleared; a refusal is never capped; the encoded body bound, C6); `TestWhyResult`
(ten shapes, (d)/(f)); **VR8** `TestWhyHandlersReadNoMiddlewareState` (AST of each `whyRoutes` method, by its
`runtime` name; control `approvalOrigins` found); **VS2** `TestSurfaceVerbWhyIsReadOnly` (imports and `d.*` exact
sets; each pattern bound to `Handler()`'s method in `daemon.go`).
**Expected cumulative numstat after the controller's commit of M5g: 781.**

### M5r — the route's `why` (design M5, routed) · 8 prod + 317 test · Run 2

**Patch** `M5r`: **`surface_verb.go`** — `case v.Words == "why": d.verbWhy(r, args, &a)` before the `terminal` case;
the doc comment names why's walk. **`lookup.go`** — the `why` row without its reason (V2b-D5); `Read`'s comment.
**`lookup_test.go`** — the two `why` rows expect no reason. **`surface_verb_outcomes_test.go`** — the bare `why` row is
`refused` with `ailang-worldd why: a target is required`; the golden list gains `why 2`, `why head --json`, `why 9`,
`why --result f`. **`verb-route.golden.json`** — regenerated by that test under `WORLD_UPDATE_VERB_GOLDEN=1` (the
`why` row changed, four rows added; V2b-D7). **`surface_verb_test.go`** — the 1,024-byte row (P1) and VS1's set +
`verbWhy`. New **`host/daemon/surface_verb_why_route_test.go`** — **VR1** `TestSurfaceVerbWhyOutcomes` (nine `result`
rows: exit 3 text and json, the resolved `rest:` arm, the world ref, five exit-1 stderr lines exactly; ten `refused`
rows with exact reasons and the CLI line, no result members; `world why` still navigates); **VR2**
`TestSurfaceVerbWhyDeadline` (expired deadline on three forms incl. the scan → (b), never "not found", and no read
starts; a 503 mid-walk → (b); a 503 and a 500 → (b) before (c); a cancelled request → (a)); **VR3**
`TestSurfaceVerbWhyReadFailure` (every `GetObject` of the walk failing in turn, then a log read → (c), no store text,
no ✗, logged); **VR4 + VR6** `TestWhyGetterTransportParity` (seven targets: `why.Run` over HTTP equals in-process,
tally clean; every recorded path byte-equal, `/v1/health` on the embedded members; `whyRoutes`' keys = the patterns
the recorded paths reach — generated); **VR7 e2e** `TestVerbAnswerCapEndToEnd` (the cap lowered to `why 9`'s size:
`why 1`, `why head --json`, `episode report w1` → refused with their sizes, no members; `why 9` at the cap exactly →
result). **`app_verb_test.go`** — the terminal rows lose `why`; VP11 gains `why` and `why --result f`; **VP9** (`why 2`
text + `exit 3`, `why head --json` json + `exit 3`, `why 9` stderr + `exit 1`; the hostile name stays text).
**`verb_bar_chrome_test.go`** (XV1) — `why` now shows `ailang-worldd why: a target is required` in the panel.
**Expected cumulative numstat after the controller's commit of M5r: 1,106.**

### M6e — the CLI's bytes, XV2, QUICKSTART (design M6, rest) · 145 test + 4 docs · Run 2

**Patch** `M6e`: new **`cmd/ailang-worldd/verb_why_e2e_test.go`** — **VA3** `TestWhyBarAdmissionIsTheCLIs` (ten usage
forms: the bar's reason equals the real CLI's first stderr line, exit 1, generated); **VE1/VE3**
`TestVerbRouteWhyIsTheCLIs` (first, on VC1's rig, eight route answers equal VC1's **pre-move** golden lines —
`--json 1`'s JSON with its one trailing newline, `1`'s chain, `9`, `rest:nobody`, the absent `sha256:`, and the first
line of `why`, `why sha256:zz`, `why 1 2`; then seventeen forms — index, `head`, `head --json`, `--json 1`, an `a2a:` and
an `mcp:` id, six `sha256:` kinds, the world ref, `9`, an exhausted `--scan 1`, `rest:nobody`, `--scan 5000` — each the
CLI's stdout/stderr/exit byte for byte on the same daemon); **XV2** `TestVerbBarWhyInChrome` (`why head` in `#addr` →
the panel's text equals the CLI's stdout; Escape removes `#verb`). **`QUICKSTART.md` §19** — `why head`, `why 3 --json`
in the block; one paragraph (V2b-D16). **`quickstart_verbs_test.go`** — `why` in the word set.
**Expected cumulative numstat after the controller's commit of M6e: 1,255.**

### 3.1 Every input and failure class, traced (design §2.7; directive retro 4)

Route = `GET /surface/verb.json` through `Handler()` (approvalRoute, V1 steps, Lookup, `verbWhy`); page = node
harness (`TestSurfaceAppVerbs`) over the route's real answers; Chrome = live (`x`).

| design §2.7 input / condition | route (test → asserted) | page / Chrome |
|---|---|---|
| `why 3` (✓ chain) | VE1 (`1`, `head`, six `sha256:` kinds, ids: text = CLI stdout, exit 0) | XV2 (`why head`: panel = CLI stdout) |
| `why head --json` | VR1 (json, exit 3), VE1 (`head --json`, `--json 1` = CLI and = VC1's pre-move JSON) | VP9 (json + `exit 3`) |
| `why <index>` with a ✗ link | VR1 (`why 1`: text, exit 3, `BROKEN — integrity refusal (exit 3)`), VR4 (= over HTTP) | VP9 (`why 2` + `exit 3`) |
| `why 999` | VR1 (`why 99` stderr exactly, exit 1), VE1 (`9` = CLI = VC1) | VP9 (`why 9` + `exit 1`) |
| `why a2a:nobody:x` (scan exhausted) | VR1 (`--scan 1`), VE1 (`--scan 1` = CLI) | via VP9's stderr shape |
| `why head` with no selected head | **residual** (a store with no head cannot be walked to a head; the walk rig and the seed always have one) | – |
| `why sha256:<unknown>` | VR1 (exact stderr), VE1 (= CLI = VC1) | – |
| `why rest:x` pending / unknown / resolved | VR1 (`indeterminate`, `not-started`, resolved → `matched by rest receipt → world`), VR4, VE1 (`not-started` = CLI = VC1) | – |
| `why`, `why x`, `why sha256:zz`, `why --scan 0 head`, `why 1 2`, `why --bogus` | VA2 (unit), VR1 (exact), VA3 (= the CLI's first stderr line), VE1 (= VC1) | VP11 (`why` → panel), XV1 (`why` → panel) |
| `why -`, `why --result f`, `why --expect-store p head`, `why --help` | VA2, VR1 | VP11 (`why --result f` → panel) |
| a read 500s mid-walk | VR3 (every `GetObject` in turn; a log read) → (c) | refusal panel (VP11's shape) |
| a read times out (503) / the walk passes 10 s | VR2 (503 → (b); expired deadline → (b), no read) | refusal panel |
| a 503 and a 500 in one walk | VR2 → (b) | – |
| the page aborts (superseded / its 15 s) | VR2 (cancelled request → (a)) | VP7, VP8, V1's supersession rows |
| an unrouted path / a request that cannot be built | VR5 (error, logged, never 404) | – |
| a sub-response over 8 MiB | VR5 (8,388,609 → error; 8,388,608 passes) | – |
| the answer over 1 MiB | VR7 (unit at 1,048,576/+1), VR7 e2e (why text, why json, episode report) | refusal panel |
| (d) the walk's answer has another shape | `TestWhyResult` (unit; unreachable through `why.Run`) | VP10 (the page refuses the same shapes) |
| daemon stalls / HTTP 4xx-5xx / network / unreadable | – (V1) | V1's VP2 rows, unchanged |
| a newer Enter while `why` is pending | – | VP8, V1's two-ask rows |
| an empty Enter while an ask is pending | – | VP7 (abort + home; no late action) |
| an Enter during `goAddress`'s `await loadIndex()` | – | VP7 (`addressGo` with a newer Enter: nothing; without: `openState`) |
| `world why` | VR1 (`address`) | V1's VP1 row |

**Residuals** (Appendix C): `why head` with no head and `ScanFor`'s "no head" (V2a's); `ScanEntries`' forward page cap;
`W-LOGPAGE`; RS-FOCUS (Enter's refocus, design §9); RV1–RV5.

---

## 4. Mutation table

**Procedure** (`tools/mutate.py <run> <tree> <workers> [ids…]`, rows in `tools/mutants.py`, enumerated from the
diff `BASE..M6e`: every new or changed production branch, constant, flag registration, route mapping, precedence row,
cap comparison and page branch, plus V2a's three named survivors): for each row, an APFS clone of `work2` at `M6e`
(`cp -cR`; byte-identical to the round-trip tree, C12), one exact-string replacement whose anchor occurs exactly once
(checked before the run: 0 bad anchors), then `go test -count=1` of the **named test**, then the **behavioural suites**
— `host/verbs` whole, `host/daemon -run '^(TestSurfaceVerb|TestWhyGetter|TestVerbAnswerCap|TestWhyResult)'`,
`host/surfacepage -run '^(TestSurfaceAppVerbs)$'`, `cmd/ailang-worldd -run '^(TestWhy|TestVerbRoute|TestQuickstartSection19)'`
— with **every static test `-skip`ped**: `TestSurfaceVerbWhyIsReadOnly`, `TestSurfaceVerbIsReadOnly`,
`TestWhyHandlersReadNoMiddlewareState` (AST/import), `TestWhyConstantsMatchTheCLI`, `TestWhyHeadIsOneReader`,
`TestWhyImportsNoEffect`, `TestSplitBoundsAreLiterals`, `TestWhyFlags` (literal/registration pins), `TestSurfaceAppPure`
(export inventory) — and `TestDemoSeedPreflightRefusals` (load list). So **no kill in the suite column rests on a
static test** (directive retro 3). The one Chrome row (`P-GO-PANEL`) runs its named test with `WORLD_CHROME` and
`WORLD_REQUIRE_PASSKEY_BROWSER=1`. Go rows ran with `GOFLAGS=-vet=off`, three workers, at load 6–25. A build failure
is `NOCOMPILE` and not a mutant. Every run appends one line per row to `plan-scratch/logs/mutant_ledger.tsv` (time,
run, tree, id, file, from, to, named regex, named result, named failures, suite result, suite failures, seconds,
expected); per-row output in `logs/mut/<run>-<id>.log`. **The ids below are the ledger's ids**; the table shows each
row's last line. Expected: `K` killed by the named test **and** by the suites; `S` killed only by a static pin (the
named test), the suites survive; `E` a named equivalent; `CTL` unmutated control, both green.

**Runs.** **R1**: every row (98) on the stage tree as first shipped. **R2**: the controls and every row whose R1 result
was not as expected or whose test changed, on the **final** stage tree (byte-identical to the round-trip tree):
`D-PREC-BC` ((c) moved before (b)) **survived R1** — VR2's "a 503 and a 500" row failed the walk's *first* object read
(the record), which ends the walk, so only one tally flag was ever set; VR2 now fails the 2nd object read (500) and then
the 3rd (503), both inside the chain's ✗-and-continue reads, and R2 re-runs `D-PREC-BC` with the rows on the same path
(`D-PREC-A`, `D-PREC-B-CTX`, `D-PREC-B-TALLY`, `D-TALLY-503`) and `CTL-D`; `D-RUN-SCAN` did not compile (`scan`
unused — `N`, not a mutant) and is rewritten as `D-RUN-SCANb`; `C-SCAN` (the CLI's `--scan` bound check removed) was
named against VC1/`TestWhyScanBound`, which **do not drive the CLI's `--scan` usage bound** (both survived R1), while the
suite's **VA3** killed it — R2 names VA3, which is therefore the CLI's own first behavioural pin of `--scan 0`/`5001`;
`P-GO-PANEL` (`go`'s DOM wiring: the panel effect replaced by `say`) is killed by live-Chrome **XV2** and by nothing in
the node or Go suites, which cannot reach the page's DOM half — class **`X`** (behavioural, live Chrome only; CI has no
Chrome, §8 R3). Every R1 line stays in the ledger; the table shows each id's last line, its `Exp` the current class
(the ledger line keeps the class it was run with).

**Classes.** `K` killed by the named test and by the suites (statics skipped); `S` killed only by a static pin —
`V-FLAG-USAGE` (a usage string: only `TestWhyFlags`' registration pin and the CLI's `--help` text, which is `whyHelp`, not
the flag usage, see it) and `W-LOGPAGE` (V2a's survivor: `LogPageMax` 500 → 499 changes no walk on any store a test
builds; VW4's literal); `X` killed only by a live-Chrome test (`P-GO-PANEL`); `N` a first form that did not compile;
`CTL` unmutated controls. **No named equivalents (`E`)**: every V2b branch has a behavioural kill.

**V2a's named survivors, re-run here** (directive input 5; R1 on the first stage tree, R2 again on the final tree, the same results): `W-MAXRESP` (`why.MaxResponseBytes` 8,388,608 → 8,388,607)
is **killed by VR5** (the 8,388,608-byte literal row now errs) and by the suites; `W-SCANMAX` (`why.ScanMax` 5,000 →
4,999) is **killed by VA3** (the CLI's `--scan 5001` line now says `1..4999` while the bar's — `verbs.WhyScanMax` — says
`1..5000`) and by the suites; `W-LOGPAGE` survives the suites as in V2a (`S`).

**Static-only kills, re-run without the static test** (directive retro 3): the suites column *is* that re-run — every
static test is `-skip`ped there — so a `K` row's suite kill never rests on VS1, VS2, VR8, VW4, VW5, VS2a, VB7, VA1 or
`TestSurfaceAppPure`. The route-binding rows (`D-R-*`, `D-MUX-REG`) are killed in the suites by the behavioural tests that
walk through the private mux — VR4 (in-process vs `Handler()` over HTTP), VR1, VE1 — not by VS2's source regex.

| Id | File | Mutation (from → to, ⏎ = newline) | Named test | Named | Suites (static skipped) | Failing tests (suites) | Exp |
|---|---|---|---|---|---|---|---|
| `CTL-V` | `flags.go` | `-` → `-` | `.` | SURVIVED | SURVIVED |  | CTL |
| `CTL-D` | `surface_verb_why.go` | `-` → `-` | `^(TestSurfaceVerb\|TestWhy\|TestVerbAnswer)` | SURVIVED | SURVIVED |  | CTL |
| `CTL-P` | `app.js` | `-` → `-` | `^(TestSurfaceAppVerbs)$` | SURVIVED | SURVIVED |  | CTL |
| `V-SCANDEF` | `flags.go` | ` WhyScanDefault = 500` → ` WhyScanDefault = 501` | `^(TestAdmitWhy)$` | KILLED | KILLED | TestAdmitWhy, TestWhyScanBound | K |
| `V-SCANMAX` | `flags.go` | ` WhyScanMax     = 5000` → ` WhyScanMax     = 4999` | `^(TestAdmitWhy)$` | KILLED | KILLED | TestAdmitWhy, TestSurfaceVerbWhyOutcomes, TestVerbRouteWhyIsTheCLIs, TestWhyBarAdmissionIsTheCLIs | K |
| `V-FLAG-ES` | `flags.go` | `fs.String("expect-store", ""` → `fs.String("expect-stor", ""` | `^(TestAdmitWhy)$` | KILLED | KILLED | TestAdmitWhy, TestSurfaceVerbWhyOutcomes, TestWhyCharacterisedPaths, TestWhyExpectStore | K |
| `V-FLAG-RESULT` | `flags.go` | `fs.String("result", ""` → `fs.String("results", ""` | `^(TestAdmitWhy)$` | KILLED | KILLED | TestAdmitWhy, TestSurfaceVerbOutcomes, TestSurfaceVerbWhyOutcomes, TestWhyCharacterisedPaths (+1) | K |
| `V-FLAG-SCAN` | `flags.go` | `fs.Int("scan", WhyScanDefault,` → `fs.Int("scan", 50,` | `^(TestAdmitWhy)$` | KILLED | KILLED | TestAdmitWhy, TestWhyScanBound | K |
| `V-FLAG-JSON` | `flags.go` | `fs.Bool("json", false, "print the chain as JSON")` → `fs.Bool("json", true, "print the chain as JSON")` | `^(TestAdmitWhy)$` | KILLED | KILLED | TestAdmitWhy, TestSurfaceVerbOutcomes, TestSurfaceVerbWhyOutcomes, TestVerbRouteWhyIsTheCLIs (+11) | K |
| `V-FLAG-USAGE` | `flags.go` | `"entries to scan back from head (max 5000)"` → `"entries to scan back (max 5000)"` | `^(TestWhyFlags)$` | KILLED | SURVIVED |  | S |
| `V-LEAD-DASH` | `flags.go` | `(!strings.HasPrefix(args[0], "-") \|\| args[0] == "-")` → `(!strings.HasPrefix(args[0], "-"))` | `^(TestAdmitWhy)$` | KILLED | KILLED | TestAdmitWhy, TestWhyCharacterisedPaths | K |
| `V-LEAD-FLAG` | `flags.go` | `(!strings.HasPrefix(args[0], "-") \|\| args[0] == "-")` → `(strings.HasPrefix(args[0], "-") \|\| args[0] == "-")` | `^(TestAdmitWhy)$` | KILLED | KILLED | TestAdmitWhy, TestSurfaceVerbOutcomes, TestSurfaceVerbWhyOutcomes, TestVerbRouteWhyIsTheCLIs (+9) | K |
| `V-TRAIL-NONE` | `flags.go` | ` if fs.NArg() == 0 {⏎  return target, nil` → ` if fs.NArg() == 0 {⏎  return "", nil` | `^(TestAdmitWhy)$` | KILLED | KILLED | TestAdmitWhy, TestSurfaceVerbOutcomes, TestSurfaceVerbWhyOutcomes, TestVerbRouteWhyIsTheCLIs (+12) | K |
| `V-TRAIL-BOTH` | `flags.go` | `if target != "" \|\| fs.NArg() > 1 {` → `if fs.NArg() > 1 {` | `^(TestAdmitWhy)$` | KILLED | KILLED | TestAdmitWhy, TestSurfaceVerbWhyOutcomes, TestVerbRouteWhyIsTheCLIs, TestWhyBarAdmissionIsTheCLIs (+1) | K |
| `V-TRAIL-TWO` | `flags.go` | `if target != "" \|\| fs.NArg() > 1 {` → `if target != "" {` | `^(TestAdmitWhy)$` | KILLED | KILLED | TestAdmitWhy, TestWhyBarAdmissionIsTheCLIs | K |
| `V-TRAIL-LAST` | `flags.go` | `fs.Arg(fs.NArg()-1))` → `fs.Arg(0))` | `^(TestAdmitWhy)$` | KILLED | KILLED | TestAdmitWhy | K |
| `V-TRAIL-TAKE` | `flags.go` | ` return fs.Arg(0), nil` → ` return target, nil` | `^(TestAdmitWhy)$` | KILLED | KILLED | TestAdmitWhy, TestVerbRouteWhyIsTheCLIs, TestWhyCharacterisedPaths | K |
| `V-SCANOK-LO` | `flags.go` | `return n >= 1 && n <= WhyScanMax` → `return n >= 0 && n <= WhyScanMax` | `^(TestAdmitWhy)$` | KILLED | KILLED | TestAdmitWhy, TestSurfaceVerbWhyOutcomes, TestWhyBarAdmissionIsTheCLIs | K |
| `V-SCANOK-HI` | `flags.go` | `return n >= 1 && n <= WhyScanMax` → `return n >= 1 && n < WhyScanMax` | `^(TestAdmitWhy)$` | KILLED | KILLED | TestAdmitWhy, TestVerbRouteWhyIsTheCLIs | K |
| `A-HELP` | `flags.go` | `  if errors.Is(err, flag.ErrHelp) {⏎   return "", 0, false, "see " + B` → `  if false && errors.Is(err, flag.ErrHelp) {⏎   return "", 0, false, "` | `^(TestAdmitWhy)$` | KILLED | KILLED | TestAdmitWhy, TestSurfaceVerbWhyOutcomes | K |
| `A-FLAGERR` | `flags.go` | `  return "", 0, false, err.Error()⏎` → `  return "", 0, false, Binary + " why: " + err.Error()⏎` | `^(TestAdmitWhy)$` | KILLED | KILLED | TestAdmitWhy, TestSurfaceVerbWhyOutcomes, TestWhyBarAdmissionIsTheCLIs | K |
| `A-TRAIL` | `flags.go` | `  return "", 0, false, Binary + " why: " + err.Error()⏎ }⏎ if !WhyScan` → `  return target, 500, false, ""⏎ }⏎ if !WhyScanOK` | `^(TestAdmitWhy)$` | KILLED | KILLED | TestAdmitWhy, TestSurfaceVerbWhyOutcomes, TestVerbRouteWhyIsTheCLIs, TestWhyBarAdmissionIsTheCLIs | K |
| `A-SCAN` | `flags.go` | ` if !WhyScanOK(*a.Scan) {⏎  return "", 0, false,` → ` if false {⏎  return "", 0, false,` | `^(TestAdmitWhy)$` | KILLED | KILLED | TestAdmitWhy, TestSurfaceVerbWhyOutcomes, TestWhyBarAdmissionIsTheCLIs | K |
| `A-RESULT` | `flags.go` | `[]bool{*a.Result != "", *a.ExpectStore != "", target == "-"}` → `[]bool{false, *a.ExpectStore != "", target == "-"}` | `^(TestAdmitWhy)$` | KILLED | KILLED | TestAdmitWhy, TestSurfaceVerbOutcomes, TestSurfaceVerbWhyOutcomes | K |
| `A-EXPECT` | `flags.go` | `[]bool{*a.Result != "", *a.ExpectStore != "", target == "-"}` → `[]bool{*a.Result != "", false, target == "-"}` | `^(TestAdmitWhy)$` | KILLED | KILLED | TestAdmitWhy, TestSurfaceVerbWhyOutcomes | K |
| `A-DASH` | `flags.go` | `[]bool{*a.Result != "", *a.ExpectStore != "", target == "-"}` → `[]bool{*a.Result != "", *a.ExpectStore != "", false}` | `^(TestAdmitWhy)$` | KILLED | KILLED | TestAdmitWhy, TestSurfaceVerbWhyOutcomes | K |
| `A-ORDER` | `flags.go` | ` {"--result", "the app reads no file on the daemon's host"},⏎ {"--expe` → ` {"--expect-store", "the app reads its own daemon's store"},⏎ {"--resu` | `^(TestAdmitWhy)$` | KILLED | KILLED | TestAdmitWhy, TestSurfaceVerbOutcomes, TestSurfaceVerbWhyOutcomes | K |
| `A-CLI` | `flags.go` | `cli := Quote(append([]string{Binary, "why"}, args...))` → `cli := Quote(append([]string{Binary}, args...))` | `^(TestAdmitWhy)$` | KILLED | KILLED | TestAdmitWhy, TestSurfaceVerbOutcomes, TestSurfaceVerbWhyOutcomes | K |
| `A-RET-JSON` | `flags.go` | `return target, *a.Scan, *a.JSON, ""` → `return target, *a.Scan, false, ""` | `^(TestAdmitWhy)$` | KILLED | KILLED | TestAdmitWhy, TestSurfaceVerbOutcomes, TestSurfaceVerbWhyOutcomes, TestVerbRouteWhyIsTheCLIs | K |
| `A-RET-SCAN` | `flags.go` | `return target, *a.Scan, *a.JSON, ""` → `return target, WhyScanDefault, *a.JSON, ""` | `^(TestAdmitWhy)$` | KILLED | KILLED | TestAdmitWhy, TestSurfaceVerbWhyOutcomes, TestVerbRouteWhyIsTheCLIs | K |
| `L-REASON` | `lookup.go` | `return Verb{Words: "why", Class: Read, CLI: cli}, argv[1:]` → `return Verb{Words: "why", Class: Read, CLI: cli, Reason: "why reads Wo` | `^(TestSurfaceVerbWhyOutcomes)$` | KILLED | KILLED | TestLookup, TestSurfaceVerbOutcomes, TestSurfaceVerbWhyOutcomes, TestVerbRouteWhyIsTheCLIs | K |
| `C-LEAD` | `why_cli.go` | `target, args := verbs.WhyLeadingTarget(args)` → `target, args := "", args` | `^(TestWhyCharacterisedPaths\|TestWhyScanBound\|TestWhyTargetFormsResolveToOneEntry)$` | KILLED | KILLED | TestVerbRouteWhyIsTheCLIs, TestWhyBarAdmissionIsTheCLIs, TestWhyCharacterisedPaths, TestWhyExpectStore (+7) | K |
| `C-TRAIL` | `why_cli.go` | ` target, err := verbs.WhyTrailingTarget(fs, target)⏎ if err != nil {` → ` target, err := verbs.WhyTrailingTarget(fs, target)⏎ if false && err !` | `^(TestWhyCharacterisedPaths\|TestWhyScanBound\|TestWhyTargetFormsResolveToOneEntry)$` | KILLED | KILLED | TestWhyBarAdmissionIsTheCLIs, TestWhyCharacterisedPaths | K |
| `C-SCAN` | `why_cli.go` | `if !verbs.WhyScanOK(*a.Scan) {` → `if false {` | `^(TestWhyBarAdmissionIsTheCLIs)$` | KILLED | KILLED | TestWhyBarAdmissionIsTheCLIs | K |
| `D-R-HEALTH` | `surface_verb_why.go` | `"GET /v1/health":        d.handleHealth,` → `"GET /v1/health":        d.handleHead,` | `^(TestWhyGetterTransportParity)$` | KILLED | KILLED | TestSurfaceVerbOutcomes, TestSurfaceVerbWhyDeadline, TestSurfaceVerbWhyOutcomes, TestSurfaceVerbWhyReadFailure (+2) | K |
| `D-R-HEAD` | `surface_verb_why.go` | `"GET /v1/head":          d.handleHead,` → `"GET /v1/head":          d.handleHealth,` | `^(TestWhyGetterTransportParity)$` | KILLED | KILLED | TestSurfaceVerbOutcomes, TestSurfaceVerbWhyOutcomes, TestWhyGetterTransportParity, TestVerbRouteWhyIsTheCLIs | K |
| `D-R-WORLD` | `surface_verb_why.go` | `"GET /v1/worlds/{ref}":  d.handleWorld,` → `"GET /v1/worlds/{ref}":  d.handleObject,` | `^(TestWhyGetterTransportParity)$` | KILLED | KILLED | TestSurfaceVerbOutcomes, TestSurfaceVerbWhyOutcomes, TestWhyGetterTransportParity, TestVerbRouteWhyIsTheCLIs | K |
| `D-R-OBJECT` | `surface_verb_why.go` | `"GET /v1/objects/{ref}": d.handleObject,` → `"GET /v1/objects/{ref}": d.handleWorld,` | `^(TestWhyGetterTransportParity)$` | KILLED | KILLED | TestSurfaceVerbOutcomes, TestSurfaceVerbWhyDeadline, TestSurfaceVerbWhyOutcomes, TestSurfaceVerbWhyReadFailure (+2) | K |
| `D-R-LOGENTRY` | `surface_verb_why.go` | `"GET /v1/log/{index}":   d.handleLogEntry,` → `"GET /v1/log/{index}":   d.handleLogRange,` | `^(TestWhyGetterTransportParity)$` | KILLED | KILLED | TestSurfaceVerbOutcomes, TestSurfaceVerbWhyOutcomes, TestVerbRouteWhyIsTheCLIs | K |
| `D-R-LOGRANGE` | `surface_verb_why.go` | `"GET /v1/log":           d.handleLogRange,` → `"GET /v1/log":           d.handleLogEntry,` | `^(TestWhyGetterTransportParity)$` | KILLED | KILLED | TestSurfaceVerbWhyOutcomes, TestWhyGetterTransportParity, TestVerbRouteWhyIsTheCLIs | K |
| `D-R-RECEIPT` | `surface_verb_why.go` | `"GET /v1/receipts/{id}": d.handleReceipt,` → `"GET /v1/receipts/{id}": d.handleHealth,` | `^(TestWhyGetterTransportParity)$` | KILLED | KILLED | TestSurfaceVerbWhyOutcomes, TestWhyGetterTransportParity, TestVerbRouteWhyIsTheCLIs | K |
| `D-R-DROP` | `surface_verb_why.go` | `  "GET /v1/receipts/{id}": d.handleReceipt,⏎` → `-` | `^(TestWhyGetterTransportParity)$` | KILLED | KILLED | TestSurfaceVerbWhyOutcomes, TestVerbRouteWhyIsTheCLIs | K |
| `D-R-EXTRA` | `surface_verb_why.go` | `  "GET /v1/receipts/{id}": d.handleReceipt,⏎` → `  "GET /v1/receipts/{id}": d.handleReceipt,⏎  "GET /v1/session/self": ` | `^(TestWhyGetterTransportParity)$` | KILLED | KILLED | TestWhyGetterTransportParity, TestWhyGetterUnroutedAndOversize | K |
| `D-MUX-REG` | `surface_verb_why.go` | `  m.HandleFunc(pattern, h)` → `  if pattern != "GET /v1/head" {⏎   m.HandleFunc(pattern, h)⏎  }` | `^(TestWhyGetterTransportParity)$` | KILLED | KILLED | TestSurfaceVerbOutcomes, TestSurfaceVerbWhyOutcomes, TestVerbRouteWhyIsTheCLIs | K |
| `D-GET-UNROUTED` | `surface_verb_why.go` | `if _, pattern := mux.Handler(req); pattern == "" {` → `if _, pattern := mux.Handler(req); false && pattern == "" {` | `^(TestWhyGetterUnroutedAndOversize)$` | KILLED | KILLED | TestWhyGetterUnroutedAndOversize | K |
| `D-GET-REQERR` | `surface_verb_why.go` | `  if err == nil {⏎   if _, pattern` → `  err = nil⏎  if err == nil {⏎   if _, pattern` | `^(TestWhyGetterUnroutedAndOversize)$` | KILLED | KILLED | TestWhyGetterUnroutedAndOversize | K |
| `D-GET-CAP` | `surface_verb_why.go` | `max: why.MaxResponseBytes}` → `max: why.MaxResponseBytes + 1}` | `^(TestWhyGetterUnroutedAndOversize)$` | KILLED | KILLED | TestWhyGetterUnroutedAndOversize | K |
| `D-REC-CMP` | `surface_verb_why.go` | `if b.body.Len()+len(p) > b.max {` → `if b.body.Len()+len(p) >= b.max {` | `^(TestWhyGetterUnroutedAndOversize)$` | KILLED | KILLED | TestWhyGetterUnroutedAndOversize | K |
| `D-REC-STATUS` | `surface_verb_why.go` | ` if b.status == 0 {⏎  b.status = status` → ` if true {⏎  b.status = status` | `^(TestWhyGetterUnroutedAndOversize)$` | KILLED | KILLED | TestSurfaceVerbOutcomes, TestSurfaceVerbWhyOutcomes, TestVerbRouteWhyIsTheCLIs | K |
| `D-REC-WRITEHDR` | `surface_verb_why.go` | ` b.WriteHeader(http.StatusOK)⏎ if b.body.Len()` → ` if b.body.Len()` | `^(TestWhyGetterUnroutedAndOversize)$` | KILLED | KILLED | TestSurfaceVerbOutcomes, TestSurfaceVerbWhyDeadline, TestSurfaceVerbWhyOutcomes, TestSurfaceVerbWhyReadFailure (+3) | K |
| `D-GET-OVER` | `surface_verb_why.go` | `   if rec.over {` → `   if false && rec.over {` | `^(TestWhyGetterUnroutedAndOversize)$` | KILLED | KILLED | TestWhyGetterUnroutedAndOversize | K |
| `D-TALLY-ERR` | `surface_verb_why.go` | `   tally.failed = true⏎   d.writeInternalErrorLog(r, err)` → `   d.writeInternalErrorLog(r, err)` | `^(TestWhyGetterUnroutedAndOversize)$` | KILLED | KILLED | TestWhyGetterUnroutedAndOversize | K |
| `D-LOG-ERR` | `surface_verb_why.go` | `   tally.failed = true⏎   d.writeInternalErrorLog(r, err)` → `   tally.failed = true` | `^(TestWhyGetterUnroutedAndOversize)$` | KILLED | KILLED | TestWhyGetterUnroutedAndOversize | K |
| `D-TALLY-503` | `surface_verb_why.go` | `case rec.status == http.StatusServiceUnavailable:` → `case false:` | `^(TestSurfaceVerbWhyDeadline)$` | KILLED | KILLED | TestSurfaceVerbWhyDeadline, TestWhyGetterUnroutedAndOversize | K |
| `D-TALLY-5XX` | `surface_verb_why.go` | `case rec.status >= 500:` → `case rec.status > 500:` | `^(TestSurfaceVerbWhyReadFailure)$` | KILLED | KILLED | TestSurfaceVerbWhyReadFailure, TestWhyGetterUnroutedAndOversize | K |
| `D-ADMIT` | `surface_verb_why.go` | ` if reason != "" {⏎  a.Reason = reason⏎  return` → ` if false && reason != "" {⏎  a.Reason = reason⏎  return` | `^(TestSurfaceVerbWhyOutcomes)$` | KILLED | KILLED | TestSurfaceVerbOutcomes, TestSurfaceVerbWhyOutcomes, TestVerbRouteWhyIsTheCLIs, TestWhyBarAdmissionIsTheCLIs | K |
| `D-PARSE` | `surface_verb_why.go` | `a.Reason = "ailang-worldd why: " + err.Error()` → `a.Reason = err.Error()` | `^(TestSurfaceVerbWhyOutcomes)$` | KILLED | KILLED | TestSurfaceVerbOutcomes, TestSurfaceVerbWhyOutcomes, TestVerbRouteWhyIsTheCLIs, TestWhyBarAdmissionIsTheCLIs | K |
| `D-PARSE-RET` | `surface_verb_why.go` | `  a.Reason = "ailang-worldd why: " + err.Error()⏎  return⏎` → `  a.Reason = "ailang-worldd why: " + err.Error()⏎` | `^(TestSurfaceVerbWhyOutcomes)$` | KILLED | KILLED | TestSurfaceVerbOutcomes, TestSurfaceVerbWhyOutcomes, TestVerbRouteWhyIsTheCLIs, TestWhyBarAdmissionIsTheCLIs | K |
| `D-READCTX` | `surface_verb_why.go` | `ctx, cancel := d.readCtx(r)` → `ctx, cancel := context.WithCancel(r.Context())` | `^(TestSurfaceVerbWhyDeadline)$` | KILLED | KILLED | TestSurfaceVerbWhyDeadline | K |
| `D-RUN-JSON` | `surface_verb_why.go` | `why.Options{Scan: scan, JSON: asJSON}` → `why.Options{Scan: scan}` | `^(TestSurfaceVerbWhyOutcomes)$` | KILLED | KILLED | TestSurfaceVerbOutcomes, TestSurfaceVerbWhyOutcomes, TestVerbRouteWhyIsTheCLIs | K |
| `D-RUN-SCAN` | `surface_verb_why.go` | `why.Options{Scan: scan, JSON: asJSON}` → `why.Options{Scan: 500, JSON: asJSON}` | `^(TestSurfaceVerbWhyOutcomes)$` | NOCOMPILE | NOCOMPILE |  | N (ledger: K) |
| `D-RUN-SCANb` | `surface_verb_why.go` | `why.Options{Scan: scan, JSON: asJSON}` → `why.Options{Scan: scan / scan * 500, JSON: asJSON}` | `^(TestSurfaceVerbWhyOutcomes)$` | KILLED | KILLED | TestSurfaceVerbWhyOutcomes, TestVerbRouteWhyIsTheCLIs | K |
| `D-PREC-A` | `surface_verb_why.go` | `case r.Context().Err() != nil:` → `case false:` | `^(TestSurfaceVerbWhyDeadline)$` | KILLED | KILLED | TestSurfaceVerbWhyDeadline | K |
| `D-PREC-B-CTX` | `surface_verb_why.go` | `case errors.Is(ctx.Err(), context.DeadlineExceeded) \|\| tally.timeout:` → `case tally.timeout:` | `^(TestSurfaceVerbWhyDeadline)$` | KILLED | KILLED | TestSurfaceVerbWhyDeadline | K |
| `D-PREC-B-TALLY` | `surface_verb_why.go` | `case errors.Is(ctx.Err(), context.DeadlineExceeded) \|\| tally.timeout:` → `case errors.Is(ctx.Err(), context.DeadlineExceeded):` | `^(TestSurfaceVerbWhyDeadline)$` | KILLED | KILLED | TestSurfaceVerbWhyDeadline | K |
| `D-PREC-C` | `surface_verb_why.go` | `case tally.failed: // (c)` → `case false: // (c)` | `^(TestSurfaceVerbWhyReadFailure)$` | KILLED | KILLED | TestSurfaceVerbWhyReadFailure | K |
| `D-PREC-BC` | `surface_verb_why.go` | ` case errors.Is(ctx.Err(), context.DeadlineExceeded) \|\| tally.timeout:` → ` case tally.failed:⏎  a.Reason = "the daemon could not read World duri` | `^(TestSurfaceVerbWhyDeadline)$` | KILLED | KILLED | TestSurfaceVerbWhyDeadline | K |
| `D-RES-BOTH` | `surface_verb_why.go` | `if out == errOut \|\|` → `if out && errOut \|\|` | `^(TestWhyResult)$` | KILLED | KILLED | TestWhyResult | K |
| `D-RES-EXIT` | `surface_verb_why.go` | `(out && ans.Exit != 0 && ans.Exit != 3)` → `(out && ans.Exit != 0)` | `^(TestWhyResult)$` | KILLED | KILLED | TestSurfaceVerbOutcomes, TestSurfaceVerbWhyOutcomes, TestVerbAnswerCapEndToEnd, TestWhyResult | K |
| `D-RES-ERREXIT` | `surface_verb_why.go` | `\|\| (errOut && ans.Exit != 1) {` → `{` | `^(TestWhyResult)$` | KILLED | KILLED | TestWhyResult | K |
| `D-RES-STDERR` | `surface_verb_why.go` | ` case errOut:⏎  a.Stderr = ans.Stderr` → ` case false:⏎  a.Stderr = ans.Stderr` | `^(TestWhyResult)$` | KILLED | KILLED | TestSurfaceVerbOutcomes, TestSurfaceVerbWhyOutcomes, TestVerbRouteWhyIsTheCLIs | K |
| `D-RES-JSON` | `surface_verb_why.go` | ` case asJSON:⏎  a.JSON = ans.Stdout` → ` case false:⏎  a.JSON = ans.Stdout` | `^(TestWhyResult)$` | KILLED | KILLED | TestSurfaceVerbOutcomes, TestSurfaceVerbWhyOutcomes, TestWhyResult, TestVerbRouteWhyIsTheCLIs | K |
| `D-RES-EXITSET` | `surface_verb_why.go` | `a.Outcome, a.Exit = "result", ans.Exit` → `a.Outcome = "result"` | `^(TestWhyResult)$` | KILLED | KILLED | TestSurfaceVerbOutcomes, TestSurfaceVerbWhyOutcomes, TestWhyResult, TestVerbRouteWhyIsTheCLIs | K |
| `D-CAP-CONST` | `surface_verb_why.go` | `const MaxVerbAnswerBytes = 1 << 20` → `const MaxVerbAnswerBytes = 1<<20 + 1` | `^(TestVerbAnswerCap)$` | KILLED | KILLED | TestVerbAnswerCap | K |
| `D-CAP-CMP` | `surface_verb_why.go` | `n > verbAnswerCap {` → `n >= verbAnswerCap {` | `^(TestVerbAnswerCap)$` | KILLED | KILLED | TestVerbAnswerCap, TestVerbAnswerCapEndToEnd | K |
| `D-CAP-OUTCOME` | `surface_verb_why.go` | `a.Outcome == "result" && n > verbAnswerCap` → `n > verbAnswerCap` | `^(TestVerbAnswerCap)$` | KILLED | KILLED | TestVerbAnswerCap | K |
| `D-CAP-SUM` | `surface_verb_why.go` | `len(a.Text) + len(a.JSON) + len(a.Stderr);` → `len(a.Text);` | `^(TestVerbAnswerCap)$` | KILLED | KILLED | TestVerbAnswerCap, TestVerbAnswerCapEndToEnd | K |
| `D-CAP-CLEAR` | `surface_verb_why.go` | `a.Outcome, a.Text, a.JSON, a.Stderr, a.Exit = "refused", "", "", "", 0` → `a.Outcome = "refused"` | `^(TestVerbAnswerCap)$` | KILLED | KILLED | TestVerbAnswerCap, TestVerbAnswerCapEndToEnd | K |
| `D-CAP-WIRE` | `surface_verb.go` | `defer func() { capVerbAnswer(&a); writeApprovalJSON(w, http.StatusOK, ` → `defer func() { writeApprovalJSON(w, http.StatusOK, a) }()` | `^(TestVerbAnswerCapEndToEnd)$` | KILLED | KILLED | TestVerbAnswerCapEndToEnd | K |
| `D-ROUTE-CASE` | `surface_verb.go` | ` case v.Words == "why":⏎  d.verbWhy(r, args, &a)⏎` → `-` | `^(TestSurfaceVerbWhyOutcomes)$` | KILLED | KILLED | TestSurfaceVerbBoundsAndErrors, TestSurfaceVerbOutcomes, TestSurfaceVerbWhyOutcomes, TestVerbRouteWhyIsTheCLIs (+1) | K |
| `W-MAXRESP` | `why.go` | ` MaxResponseBytes = 8388608` → ` MaxResponseBytes = 8388607` | `^(TestWhyGetterUnroutedAndOversize)$` | KILLED | KILLED | TestWhyGetterUnroutedAndOversize | K |
| `W-SCANMAX` | `why.go` | ` ScanMax     = 5000` → ` ScanMax     = 4999` | `^(TestWhyBarAdmissionIsTheCLIs)$` | KILLED | KILLED | TestWhyBarAdmissionIsTheCLIs | K |
| `W-LOGPAGE` | `why.go` | ` LogPageMax  = 500` → ` LogPageMax  = 499` | `^(TestWhyConstantsMatchTheCLI)$` | KILLED | SURVIVED |  | S |
| `P-REF-PANEL` | `app.js` | `return typeof b.cli === 'string' && b.cli !== '' ? { panel:` → `return false && typeof b.cli === 'string' && b.cli !== '' ? { panel:` | `^(TestSurfaceAppVerbs)$` | KILLED | KILLED | TestSurfaceAppVerbs | K |
| `P-REF-CLI-EMPTY` | `app.js` | `return typeof b.cli === 'string' && b.cli !== '' ? { panel:` → `return typeof b.cli === 'string' ? { panel:` | `^(TestSurfaceAppVerbs)$` | KILLED | KILLED | TestSurfaceAppVerbs | K |
| `P-REF-BODY` | `app.js` | `{ panel: { head: '$ ' + b.cli, body: b.reason, copy: b.cli } }` → `{ panel: { head: '$ ' + b.cli, body: b.reason, copy: '' } }` | `^(TestSurfaceAppVerbs)$` | KILLED | KILLED | TestSurfaceAppVerbs | K |
| `P-EXIT-3` | `app.js` | `(b.exit !== undefined && b.exit !== 1 && b.exit !== 3)` → `(b.exit !== undefined && b.exit !== 1)` | `^(TestSurfaceAppVerbs)$` | KILLED | KILLED | TestSurfaceAppVerbs | K |
| `P-EXIT-2` | `app.js` | `(b.exit !== undefined && b.exit !== 1 && b.exit !== 3)` → `(b.exit !== undefined && b.exit !== 1 && b.exit !== 3 && b.exit !== 2)` | `^(TestSurfaceAppVerbs)$` | KILLED | KILLED | TestSurfaceAppVerbs | K |
| `P-STDERR-IFF` | `app.js` | `(b.stderr !== undefined) !== (b.exit === 1)` → `(b.stderr !== undefined && b.exit !== 1)` | `^(TestSurfaceAppVerbs)$` | KILLED | KILLED | TestSurfaceAppVerbs | K |
| `P-BODY-EXIT` | `app.js` | `body: shown[0] + (b.exit ? 'exit ' + b.exit : '')` → `body: shown[0]` | `^(TestSurfaceAppVerbs)$` | KILLED | KILLED | TestSurfaceAppVerbs | K |
| `P-RESULT-CLI` | `app.js` | `result: ['cli']` → `result: []` | `^(TestSurfaceAppVerbs)$` | KILLED | KILLED | TestSurfaceAppVerbs | K |
| `P-DROP` | `app.js` | `function dropVerb() { if (verbJob) { verbJob.superseded = true; verbJo` → `function dropVerb() { }` | `^(TestSurfaceAppVerbs)$` | KILLED | KILLED | TestSurfaceAppVerbs | K |
| `P-GO-EMPTY-DROP` | `app.js` | `{ dropVerb(); return fx.out(); }` → `{ return fx.out(); }` | `^(TestSurfaceAppVerbs)$` | KILLED | KILLED | TestSurfaceAppVerbs | K |
| `P-GO-CURRENT` | `app.js` | `    if (!current()) return;⏎    return verbAct` → `    return verbAct` | `^(TestSurfaceAppVerbs)$` | KILLED | KILLED | TestSurfaceAppVerbs | K |
| `P-GO-ADDRESS` | `app.js` | `address: r => fx.address(r, current)` → `address: r => fx.address(r, () => true)` | `^(TestSurfaceAppVerbs)$` | KILLED | KILLED | TestSurfaceAppVerbs | K |
| `P-GO-TURN` | `app.js` | `current = () => turn === barTurn;` → `current = () => true;` | `^(TestSurfaceAppVerbs)$` | KILLED | KILLED | TestSurfaceAppVerbs | K |
| `P-ADDR-CURRENT` | `app.js` | `      if (!current()) return;⏎` → `-` | `^(TestSurfaceAppVerbs)$` | KILLED | KILLED | TestSurfaceAppVerbs | K |
| `P-FINALLY` | `app.js` | `if (verbJob === job) verbJob = null;` → `verbJob = null;` | `^(TestSurfaceAppVerbs)$` | KILLED | KILLED | TestSurfaceAppVerbs | K |
| `P-GO-PANEL` | `app.js` | `act: { assign: h => location.assign(h), panel: showVerbPanel, say,` → `act: { assign: h => location.assign(h), panel: say, say,` | `^(TestVerbBarWhyInChrome)$` | KILLED | SURVIVED |  | X |

---

## 5. Files (final `git status --short --untracked-files=all` after M6e, against the base)

```
 M cmd/ailang-worldd/quickstart_verbs_test.go
 M cmd/ailang-worldd/verb_bar_chrome_test.go
 M cmd/ailang-worldd/why_cli.go
 M cmd/ailang-worldd/why_seam_test.go
 M docs/QUICKSTART.md
 M host/daemon/surface_verb.go
 M host/daemon/surface_verb_outcomes_test.go
 M host/daemon/surface_verb_test.go
 M host/surfacepage/app/app.js
 M host/surfacepage/app/index.html
 M host/surfacepage/app_test.go
 M host/surfacepage/app_verb_test.go
 M host/surfacepage/testdata/verb-route.golden.json
 M host/verbs/flags.go
 M host/verbs/lookup.go
 M host/verbs/lookup_test.go
?? cmd/ailang-worldd/verb_why_e2e_test.go
?? host/daemon/surface_verb_why.go
?? host/daemon/surface_verb_why_route_test.go
?? host/daemon/surface_verb_why_test.go
?? host/verbs/bounds_test.go
?? host/verbs/why_admit_test.go
```
Plus the two plan files, which the controller commits. After each sub-milestone, exactly the files its §3 paragraph
names (cumulatively). Nothing under `world/`, `packages/`, `scripts/`, `tools/`, `host/why/`, `host/store/`,
`host/broker/`, nor `host/daemon/daemon.go`, `handlers.go`, `cmd/ailang-worldd/testdata/`.

---

## 6. Sprint gates (controller, after the last sub-milestone's commit; each measured at the base, §1.1)

Run it in **one bash shell, in the foreground, never under `nohup`** (macOS `/bin/bash` 3.2; every leg records its own rc
and every leg runs — no `set -e`; every leg is `tee`'d to `gates/<leg>.txt`), extracting the script from this plan. The
executor sets `GATES_SKIP_VERIFY_GO=1`; the controller runs it without, after the last commit (`verify_go` reads
committed blobs):

```bash
bash <<'EOF2'
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter398-row243-v2b
G="$HOME/.ailang/state/world-iter398/gates"; mkdir -p "$G"
awk '$0 == "````bash gates.sh" {f=1; next} f && $0 == "````" {exit} f' "$W/design_docs/planned/w-world-address-bar-verbs-v2b-sprint-plan.md" > "$G/gates.sh"
/bin/bash "$G/gates.sh" "$W" "$G" final "$(git -C "$W" rev-parse HEAD)"
EOF2
```

It writes `gates/legs.txt`, one `<leg>.txt` per leg, and **`gates/ALL_DONE.<sha>`** (`GATES_DONE <s> FAILS: …`) when every
leg has finished (the judge waits on it). The legs: `status`, `verify_ail` (no `.ail` changed), `gofmt`, `build`, `vet`
(`go vet ./...`), `fence` (`go test -run '^$' ./...`, the test-file compile fence), `pkgs` (`host/verbs`, `host/why`,
`host/surfacepage`, `host/daemon`, `cmd/ailang-worldd` whole), `race` (`-race` on the three small packages whole and on
the touched tests of `host/daemon` and `cmd`), `inv` (every named inventory `-v` + every new test: none missing, none
SKIP — incl. `TestDaemonDependencyAllowlist`, `TestWorldBoundaryDependencyAllowlist`, `TestProductionContextRoots`,
`TestProductionGoSurface`, `TestEveryVerbHasHelp`, `TestFrozenV1RouteTableMatchesMux`, `TestSurfaceAppPure`, VS1, VS2,
the golden test), **`realstack`** (`TestDemoSeedWhyNamedLinks/real-stack` and `TestCLIAgainstSeToolsDaemon` with
`WORLD_TOOL_AILANG_BIN=$HOME/.pinned-ailang-tools/v0.54.15/ailang`: `--- PASS`, never SKIP), **`chrome`** (V1's twelve
passkey/approve/decide/Needs-you page tests, XV1, **XV2**, and the two `host/daemon` app-shell drills, with
`WORLD_CHROME` and `WORLD_REQUIRE_PASSKEY_BROWSER=1`: temp profiles only), `verify_go` (`VERIFY_GO_LEGS=plain
AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_go.sh`). No production context root is added (V2b adds no
`context.Background`/`TODO` to production code), so `TestProductionContextRoots` is an inventory, not a changed gate. The
script, verbatim:

````bash gates.sh
#!/bin/bash
# Row 243 V2b gates (plan §6). usage: gates.sh <tree> <logdir> <base|final> [<sha>]
# The controller: bash gates.sh <worktree> $HOME/.ailang/state/world-iter398/gates final <HEAD sha>
# macOS /bin/bash 3.2: no declare -A. Every leg records its own rc; every leg runs (no set -e); every leg is tee'd
# to <logdir>/<leg>.txt. GATES_SKIP_VERIFY_GO=1 skips verify_go (the executor sets it; the controller runs it).
# Never start this under nohup. Writes <logdir>/ALL_DONE.<sha> with FAILS: when every leg has finished.
T=$1; L=$2; MODE=$3
mkdir -p "$L"; cd "$T" || exit 2
SHA=${4:-$(git rev-parse HEAD 2>/dev/null || echo unknown)}
export AILANG_BIN=$HOME/.pinned-ailang/ailang
export WORLD_EXEC_NODE=/opt/homebrew/bin/node
unset WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER WORLD_PASSKEY_RECORD_DIR WORLD_UPDATE_VERB_GOLDEN WORLD_UPDATE_WHY_GOLDEN WORLD_TOOL_AILANG_BIN
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
leg status 'git status --porcelain | wc -l'
leg verify_ail './scripts/verify_ail.sh'
leg gofmt 'out=$(gofmt -l cmd host); echo "$out"; test -z "$out"'
leg build 'go build ./...'
leg vet 'go vet ./...'
leg fence "go test -run '^\$' ./... > '$L/fence.raw' 2>&1; rc=\$?; cat '$L/fence.raw'; if grep -E '^(FAIL|--- FAIL)' '$L/fence.raw' >/dev/null; then exit 1; fi; exit \$rc"
leg pkgs "go test -count=1 -timeout 30m ./host/verbs/ ./host/why/ ./host/surfacepage/ ./host/daemon/ ./cmd/ailang-worldd/"
leg race "go test -race -count=1 ./host/verbs/ ./host/why/ ./host/surfacepage/ && go test -race -count=1 ./host/daemon/ -run '^(TestSurface|TestWhy|TestVerb|TestEpisodeReport)' && go test -race -count=1 ./cmd/ailang-worldd/ -run '^(TestWhy|TestVerb|TestQuickstartSection19)'"
NEW_V="TestWhyFlags TestAdmitWhy TestNeverAndOtherReasonsEndWithTheCLILine TestSplitBoundsAreLiterals"
NEW_D="TestSurfaceVerbWhyOutcomes TestSurfaceVerbWhyDeadline TestSurfaceVerbWhyReadFailure TestWhyGetterTransportParity TestWhyGetterUnroutedAndOversize TestVerbAnswerCap TestVerbAnswerCapEndToEnd TestWhyHandlersReadNoMiddlewareState TestSurfaceVerbWhyIsReadOnly TestWhyResult"
NEW_C="TestWhyBarAdmissionIsTheCLIs TestVerbRouteWhyIsTheCLIs"
if [ "$MODE" = base ]; then NEW_V=""; NEW_D=""; NEW_C=""; fi
INV_DAEMON="TestDaemonDependencyAllowlist TestFrozenV1RouteTableMatchesMux TestSurfaceAppRoutes TestApprovalRoutesRegistered TestDecisionRoutesRegistered TestReadCtxCancelledAfterHandler TestDaemonCredentialWriterSites TestDaemonConstructorCallerInventory TestGETRoutesRejectOtherMethods TestDaemonCallerInventoryAliasControl TestEpisodeReportRESTRoute TestSurfaceVerbRouteRegistered TestSurfaceVerbOutcomes TestSurfaceVerbBoundsAndErrors TestSurfaceVerbIsReadOnly"
INV_BOUNDARY="TestWorldBoundaryDependencyAllowlist TestBareNetHTTPExemptionIsPerGroup TestEveryCommandDirectoryIsAProtectedGroup TestBoundaryASTWriteGuard TestWorkbenchPackageRemainsTransportFree TestPinfetchStaysOutOfStoreReplayAndPublish"
INV_STORE="TestProductionContextRoots TestProductionGoSurface TestNoNewDeadlineFreeStoreReads TestReadGettersHonorContext"
INV_VG="TestNoParallelWriteForkPackagesOutsideVerifygate TestSubprocessSinkGate TestVerifygateTestWritesAreForkLocked"
INV_SP="TestSurfaceAppPure TestSurfaceScriptPure TestSurfaceAppMatchesWorkingCopy TestSurfaceAppShellAssembly TestSurfaceAppVerbs"
INV_VERBS="TestSplit TestQuoteRoundTrip TestLookup TestEpisodeReportFlags TestParseInterleaved TestApproveTarget"
INV_WHY="TestWhyImportsNoEffect TestParseTargetForms TestHeadAndGetJSONClasses TestElideBounds"
INV_CMD="TestEveryVerbHasHelp TestVerbTableMatchesTheCLISwitch TestVerbRouteReportIsTheCLIs TestQuickstartSection19MatchesTheBar TestWhyCharacterisedPaths TestWhyRunReplaysTheCLI TestWhyRunGetterClasses TestWhyConstantsMatchTheCLI TestWhyHeadIsOneReader TestWhyTargetFormsResolveToOneEntry TestWhyScanBound TestWhyBrokenLinkIsIntegrityRefusal"
leg inv "set -o pipefail
go test -v -count=1 ./host/daemon/ -run '^($(echo $INV_DAEMON | tr ' ' '|')$(alt $NEW_D))\$' &&
go test -v -count=1 ./host/boundary/ -run '^($(echo $INV_BOUNDARY | tr ' ' '|'))\$' &&
go test -v -count=1 ./host/store/ -run '^($(echo $INV_STORE | tr ' ' '|'))\$' &&
go test -v -count=1 ./host/verifygate/ -run '^($(echo $INV_VG | tr ' ' '|'))\$' &&
go test -v -count=1 ./host/surfacepage/ -run '^($(echo $INV_SP | tr ' ' '|'))\$' &&
go test -v -count=1 ./host/verbs/ -run '^($(echo $INV_VERBS | tr ' ' '|')$(alt $NEW_V))\$' &&
go test -v -count=1 ./host/why/ -run '^($(echo $INV_WHY | tr ' ' '|'))\$' &&
go test -v -count=1 ./cmd/ailang-worldd/ -run '^($(echo $INV_CMD | tr ' ' '|')$(alt $NEW_C))\$'"
need inv $INV_DAEMON $NEW_D $INV_BOUNDARY $INV_STORE $INV_VG $INV_SP $INV_VERBS $NEW_V $INV_WHY $INV_CMD $NEW_C
if grep -E '^\s*--- (SKIP|FAIL)' "$L/inv.txt"; then FAILS="$FAILS inv_skip_or_fail"; fi
echo "named PASS: $(grep -cE '^--- PASS' "$L/inv.txt")" | tee -a "$L/legs.txt"
# the real-stack why leg: CI skips these (WORLD_TOOL_AILANG_BIN unset); here they must PASS, not SKIP
leg realstack "set -o pipefail; WORLD_TOOL_AILANG_BIN=\$HOME/.pinned-ailang-tools/v0.54.15/ailang go test -count=1 -v -timeout 20m ./cmd/ailang-worldd/ -run '^(TestDemoSeedWhyNamedLinks|TestCLIAgainstSeToolsDaemon)\$'"
need realstack 'TestDemoSeedWhyNamedLinks/real-stack' TestDemoSeedWhyNamedLinks TestCLIAgainstSeToolsDaemon
if grep -E -- '--- (SKIP|FAIL)' "$L/realstack.txt"; then FAILS="$FAILS realstack_skip_or_fail"; fi
XB="TestApprovePageInChrome TestApprovePageRefusesASwappedRequest TestApprovePageSSHSIGInChrome TestApprovePageShowsAHostileValueAsText TestPasskeyPageAbortsAPendingCeremony TestPasskeyPageInChrome TestApproveDecideWebAuthnLinkInChrome TestDecidePageApproveWithPasskey TestDecidePageFailsClosed TestDecidePageSSHSIGPaste TestNeedsYouListsPendingRequests TestNeedsYouNotifiesOnlyWhenAskedOnce TestVerbBarEpisodeReportInChrome"
NEW_X="TestVerbBarWhyInChrome"; [ "$MODE" = base ] && NEW_X=""
XD="TestSurfaceLiveChrome TestSurfaceFollowChrome"
CH="'/Applications/Google Chrome.app/Contents/MacOS/Google Chrome'"
leg chrome "set -o pipefail; WORLD_CHROME=$CH WORLD_REQUIRE_PASSKEY_BROWSER=1 WORLD_REQUIRE_SSH_KEYGEN=1 go test ./cmd/ailang-worldd/ -run '^($(echo $XB | tr ' ' '|')$(alt $NEW_X))\$' -count=1 -v && WORLD_CHROME=$CH go test ./host/daemon/ -run '^($(echo $XD | tr ' ' '|'))\$' -count=1 -v"
need chrome $XB $NEW_X $XD
if [ "${GATES_SKIP_VERIFY_GO:-0}" = 1 ]; then echo "LEG verify_go SKIPPED (GATES_SKIP_VERIFY_GO=1)" | tee -a "$L/legs.txt"
else leg verify_go 'VERIFY_GO_LEGS=plain AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_go.sh'; fi
echo "GATES_DONE $(( $(date +%s) - s0 ))s FAILS:${FAILS:- none}" | tee -a "$L/legs.txt"
echo "GATES_DONE $(( $(date +%s) - s0 ))s FAILS:${FAILS:- none}" > "$L/ALL_DONE.$SHA"
````

**Expected:** every leg rc 0 except a known-load red in `verify_go`, re-run alone once (§1.1's rule); `inv` prints
**named PASS** for every inventory and every new test; `realstack` and `chrome` show every named test `--- PASS`.
The planner's run of this script on the final tree: §1.2 C14.

---

## 7. Executor directive notes

- **No git writes of any kind** (no add, rm, mv, commit, stash, checkout, reset, branch). **`patch -p1 -E` (§2) is the
  only writer** of the tree; no hand edit, no `gofmt -w`, no fuzz, no `WORLD_UPDATE_VERB_GOLDEN=1` /
  `WORLD_UPDATE_WHY_GOLDEN=1` (the golden ships in M5r's patch). A failed `--dry-run` or a red acceptance → **stop and
  report** with the log; do not repair.
- **Never `nohup`, never background a test leg; pass a 600000 ms timeout to every long Bash call.** Each acceptance arm
  fits one call (C11); if a call nears 600 s under load, run `accept.sh` for that sub-milestone again rather than
  backgrounding it.
- **Runs**: **Run 1** — M4 → M6p → M5g; **Run 2** — M5r → M6e. **The order is fixed** (M6p's patch assumes M4's tree;
  M5r's assumes M6p's and M5g's). After each sub-milestone's acceptance, snapshot **cumulatively against the base** to
  `$HOME/.ailang/state/world-iter398/snap/<M>/`: every file added or changed since `c09e0cdb` (repo-relative paths),
  a **`DELETED`** file (one repo path per line; V2b deletes nothing, so it is empty), a **`STATUS`** file, and
  `MANIFEST.sha256` (the two plan files, `STATUS` and `DELETED` are not in it). **Never an `rsync` of whole
  directories.** The `find` takes only the top-level dirs that exist (M4–M5r have no `docs/`), rehearsed in scratch on
  M4 (C17):
  ```bash
  bash <<'EOF2'
  set -euo pipefail
  M=M4   # the sub-milestone just accepted
  B=c09e0cdb3c8a3eb25c208dbacedddfb690fe8ca0
  W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter398-row243-v2b
  D="$HOME/.ailang/state/world-iter398/snap/$M"; rm -rf "$D"; mkdir -p "$D"
  cd "$W"
  { git diff --name-status --no-renames "$B" -- cmd host docs; git ls-files --others --exclude-standard -- cmd host docs | sed 's/^/A\t/'; } | sort -u -t "$(printf '\t')" -k2 > "$D/STATUS"
  awk -F'\t' '$1 != "D" {print $2}' "$D/STATUS" | while read -r p; do mkdir -p "$D/$(dirname "$p")"; cp -p "$p" "$D/$p"; done
  awk -F'\t' '$1 == "D" {print $2}' "$D/STATUS" > "$D/DELETED"
  (cd "$D" && find $(for d in cmd host docs; do [ -d "$d" ] && echo "$d"; done) -type f -print0 | sort -z | xargs -0 shasum -a 256) > "$D/MANIFEST.sha256"
  cat "$D/STATUS"; echo "DELETED: $(wc -l < "$D/DELETED")"; echo "MANIFEST: $(wc -l < "$D/MANIFEST.sha256")"
  EOF2
  ```
  Expected (`STATUS` lines / `DELETED` / MANIFEST): **M4** 5 / 0 / 5; **M6p** 10 / 0 / 10; **M5g** 13 / 0 / 13; **M5r**
  19 / 0 / 19; **M6e** 22 / 0 / 22. The controller copies the snapshot over the worktree, `git rm`s each `DELETED` path
  (none), commits each snapshot as one commit, and checks the numstat (`git diff -M --numstat c09e0cdb HEAD -- cmd host
  docs | awk '{s+=$1+$2} END {print s}'`): **M4 255, M6p 382, M5g 781, M5r 1,106, M6e 1,255**. A different number → stop.
- `export AILANG_BIN=$HOME/.pinned-ailang/ailang` in every shell that runs Go tests (`accept.sh` sets it, and
  `WORLD_EXEC_NODE=/opt/homebrew/bin/node`). The controller's shell is zsh: run each block under `bash <<'EOF2' … EOF2`.
- **Live Chrome** (M6p, M5r, M6e acceptance; the `chrome` gate): `/Applications/Google Chrome.app/Contents/MacOS/Google
  Chrome` through the passkey driver, which starts it on a `t.TempDir()` profile — never a real profile.
- **Stop conditions**: a dry-run failure; any `--- FAIL`, `--- SKIP` or `FAIL` line in an acceptance (`ACCEPT_RED`; a
  live-Chrome `x` step that times out under load is re-run once, alone, first — §8 R6);
  gofmt output; a numstat mismatch; any file outside §5 changed (`git status --short` after each patch must equal the
  cumulative list of §3's paragraphs).
- **Known load-sensitive tests** (the directive's list) are not this sprint's: a red there is re-run once alone.
- D-WORLD-109: nothing touches `~/.ailang/world`, `~/.ssh`, a real `/dev/tty`, a real browser profile, or a real daemon.

---

## 8. Risks

- **R1 — Size margin 145.** The patches are exact; any executor-side growth is a stop (§7).
- **R2 — Load.** The `host/daemon` and `cmd` packages carry the directive's load-sensitive tests; a red outside the
  touched tests is re-run alone once (§1.1).
- **R3 — Live Chrome only on the controller's machine** (CI has no Chrome step): XV1/XV2 and `P-GO-PANEL` are local
  proofs; the node harness (VP7–VP11) carries the page's branches in CI.
- **R4 — The route golden grows** by four `why` answers (≈ 6 KB). It is regenerated only by its own test.
- **R6 — Live Chrome under heavy load.** The drivers' `waitFor` gives a page 15 s; at load ≈ 28 the planner saw XV1's
  first panel miss it once (C11), green alone in 1.6 s. A Chrome red in acceptance is re-run once alone (the
  directive's load rule); a second red is a stop.
- **R5 — The in-process getter mirrors `Handler()`'s seven registrations.** A future change that re-binds a `/v1` GET
  handler in `daemon.go` fails VS2 (static) and VR4 (behavioural); a future walk read outside the seven is an error,
  never a 404 (VR5).

---

## 9. Acceptance criteria (the judge scores against these)

1. **The bar answers `why` with the CLI's bytes**: for every target form the walk rig resolves, the route's
   `text`/`json`/`stderr`/`exit` equal `ailang-worldd why`'s stdout/stderr/exit on the same daemon (VE1), and the bar's
   `--json 1`, `1`, `9`, `rest:nobody`, absent-`sha256:` answers equal V2a's pre-move VC1 golden (VE1, not derived from
   `why.Run`); exit 3 reaches the page with its BROKEN line (VR1, VP9).
2. **One parse**: the CLI's `runWhy` uses `verbs.WhyFlags`/`WhyLeadingTarget`/`WhyTrailingTarget`/`WhyScanOK` (M4); every
   usage refusal of the bar equals the CLI's first stderr line (VA3); every §2.3 form has an exact row (VA2); `-`,
   `--result`, `--expect-store` are refused with the CLI line before `ParseTarget` (VA2, VR1).
3. **CLI behaviour-neutral**: VC1 (`why-characterise.golden.jsonl`) and `why_test.go` pass **unedited**; the real-stack
   leg passes (`--- PASS: TestDemoSeedWhyNamedLinks/real-stack`, `TestCLIAgainstSeToolsDaemon`).
4. **Exactly these pre-existing test files change**, each for a named reason: `lookup_test.go` (the `why` row has no
   reason, V2b-D5), `surface_verb_outcomes_test.go` (bare `why` is refused; four `why` golden rows), `surface_verb_test.go`
   (VS1 + `verbWhy`; the 1,024-byte `why` line is a refusal, not `terminal`), `app_verb_test.go` (A2–A4 rows; `why`
   no longer the terminal example), `app_test.go` (two exports), `verb_bar_chrome_test.go` (XV1: A1/A3 change its `why`
   and `log` arms), `quickstart_verbs_test.go` (§19 names `why`), `why_seam_test.go` (VW4: §2.2's constants). The golden
   is regenerated by its test only.
5. **Read-only**: no new `/v1` route (`TestFrozenV1RouteTableMatchesMux`, `TestSurfaceVerbRouteRegistered` unedited, green),
   no POST, no store write, no process; VS1 (`surface_verb.go` + `verbWhy` only) and VS2 (exact imports and `d.*`
   set: the seven handlers, `whyRoutes`, `whyMux`, `whyGetter`, `readCtx`, `readDeadline`, `writeInternalErrorLog`;
   no `sync.Once`, V2b-D3) pass; VR4/VR6 (the private mux = the walk's seven routes, byte-identical to `Handler()`) and
   VR8 (no middleware state read) pass — T1 rests on all four.
6. **Precedence (design §2.5)**: (a) cancelled, (b) deadline or 503 — before (c), (c) a failed read, (d) an unshowable
   answer, (e) the 1 MiB cap in the deferred writer for `why` and `episode report`, (f) the result — each driven by its
   own test (VR2, VR3, `TestWhyResult`, VR7, VR1); never "not found" for a deadline, never ✗ for a read failure.
7. **Audit inputs 4–6 resolved**: VS2's full set (C5); the stale reason gone with a failing-on-stale test (C4, `L-REASON`);
   the cap on raw bytes with the encoded bound tested (C6) and stated in RV2 (V2b-D6) and the code comment.
8. **The page (A2–A4)**: exit 3 shown under its text/json, `stderr ⇔ exit 1`, `result` needs `cli` (VP9, VP10);
   refusals with a CLI line open the panel, without stay a `say` (VP11, XV1); every Enter owns the bar — an empty Enter
   aborts a pending ask, a late answer and a post-`loadIndex` address do nothing (VP7); the third-ask supersession
   (VP8); Escape closes the panel in live Chrome (XV2).
9. **The V1 judge notes and V2a's notes disposed** as design §9 and V2b-D9 say; VB6, VB7 pass.
10. **Live Chrome**: XV2 shows `why head`'s panel equal to the CLI's stdout; XV1 green with its re-pointed arms.
11. **Docs**: QUICKSTART §19 shows `why head` / `why 3 --json` and states the 10 s deadline and the refusals;
    `TestQuickstartSection19MatchesTheBar` passes.
12. **Size** ≤ 1,400 (`git diff -M --numstat c09e0cdb HEAD -- cmd host docs`): 1,255 measured; each sub-milestone's
    numstat as §3.
13. **Every mutant in §4 behaves as expected**, its id in the ledger; no suite kill rests on a static test; V2a's
    `W-MAXRESP` and `W-SCANMAX` are killed behaviourally.
14. **The round trip is identical** (C12) and the gates (§6) pass — `vet`, `realstack`, `chrome`, and `verify_go` after
    commit included; `./scripts/verify_ail.sh` gives the base's result.

**What an operator can do after this sprint:** type `why head`, `why 3`, `why sha256:<ref>`, `why a2a:<id>` (with
`--scan N`, `--json`) in the World app's address bar and read the same chain the terminal prints, with `exit 3` when a
link is broken; a walk the app cannot finish in 10 s, or a read that fails, is refused with the exact CLI line to copy.

---

## Appendix C — after this sprint (V2b's close; what V3 inherits)

- **Moved out: nothing.** 1,255 ≤ 1,400 (V2b-D17).
- **Residuals carried**: **RV1** (the bar answers only walks that finish in 10 s; past it, (b) with the CLI line),
  **RV2** (amended, V2b-D6: the 1 MiB cap is on raw bytes; the encoded body is bounded by 6× plus members; a large
  `--json` worktree delta is refused with its size), **RV3** (the CLI keeps P4/P5: at a terminal a deadline prints "not
  found" and a read failure a ✗; unchanged and pinned by VW3/VC1), **RV4** (V2a's compatibility layer), **RV5** (rename
  pairing, V2a only); **RS-FOCUS** (Enter's refocus after `go`, a DOM-only effect; untested, as design §9);
  **`why head` with no selected head** / `ScanFor`'s "no head" / `HeadIndex`'s headless world (no store reachable
  through `POST /v1/commit` has them); `ScanEntries`' forward page cap (> 500 entries forward); **`W-LOGPAGE`**
  (`LogPageMax` 500 → 499 changes no walk on any store a test builds — every store < 499 entries; pinned by VW4's
  literal). **Closed here**: **RS-REST** (VR1, VR4: the resolved `rest:` arm walks through the route), V2a's
  `W-MAXRESP` (VR5) and `W-SCANMAX` (VA3).
- **The V2a judge's unpinned lines**: `--scan 0` (pinned: VA2, VR1, VA3), the `rest:` state operand (pinned: VR1, VE1),
  the head-absent error and `why head` with no head (residual, above), the fallback-plan operand (the CLI's `-` only;
  the bar refuses `-`; VC1 keeps it), a mid-rune `Elide` cut (V2a's `TestElideBounds`; host/why untouched).
- **What row 243 V3 (the verb-approval law) inherits**: the route's outcome vocabulary is now `address | navigate |
  result | terminal | refused` with `result` carrying `exit` 0/1/3 and capped at 1 MiB in one place
  (`capVerbAnswer`); a `refused` with `cli` opens the panel (A3), so V3's "approve at a terminal" refusals carry their
  line for free; `goWith` is the bar's one Enter path (A4) — a V3 verb that navigates or posts must go through
  `verbAct`'s `act` effects and respect `current()`; `AdmitWhy`'s shape (the CLI's flag set + a bar profile that refuses
  reads/effects the app does not make) is the template for admitting a mutating verb's flags; the in-process getter is
  GET-only by construction (VS2, VR6) — V3's writes need their own seam, not this mux.

---

## Appendix P — the patches (apply with §2, in order; generated from the compiled scratch stages; do not edit)

### Patch M4 — admission: design §2.2 + `AdmitWhy`; VA1, VA2, VB6, VB7

````diff M4
diff --git a/cmd/ailang-worldd/why_cli.go b/cmd/ailang-worldd/why_cli.go
--- a/cmd/ailang-worldd/why_cli.go
+++ b/cmd/ailang-worldd/why_cli.go
@@ -12,10 +12,10 @@ import (
 	"io"
 	"net/http"
 	"os"
-	"strings"
 	"time"
 
 	"github.com/sunholo-data/ailang-world/host/broker"
+	"github.com/sunholo-data/ailang-world/host/verbs"
 	"github.com/sunholo-data/ailang-world/host/why"
 )
 
@@ -98,31 +98,24 @@ func trimOneNewline(b []byte) []byte {
 }
 
 func runWhy(addr string, args []string, stdin io.Reader, stdout, stderr io.Writer) int {
-	target := ""
-	if len(args) > 0 && (!strings.HasPrefix(args[0], "-") || args[0] == "-") {
-		target, args = args[0], args[1:]
-	}
+	target, args := verbs.WhyLeadingTarget(args)
 	fs := flag.NewFlagSet("ailang-worldd why", flag.ContinueOnError)
 	fs.SetOutput(stderr)
 	fs.Usage = func() {}
-	expectStore := fs.String("expect-store", "", "refuse unless the daemon serves this store (World-Provenance store= path)")
-	resultFile := fs.String("result", "", "a 'call --json-out' output file")
-	scan := fs.Int("scan", whyScanDefault, "entries to scan back from head (max 5000)")
-	asJSON := fs.Bool("json", false, "print the chain as JSON")
+	a := verbs.WhyFlags(fs)
 	if code, done := parseVerbFlags(fs, args, whyHelp, stdout, stderr); done {
 		return code
 	}
-	if fs.NArg() > 0 {
-		if target != "" || fs.NArg() > 1 {
-			fmt.Fprintf(stderr, "ailang-worldd why: unexpected argument %q\n", fs.Arg(fs.NArg()-1))
-			return exitUsage
-		}
-		target = fs.Arg(0)
+	target, err := verbs.WhyTrailingTarget(fs, target)
+	if err != nil {
+		fmt.Fprintf(stderr, "ailang-worldd why: %v\n", err)
+		return exitUsage
 	}
-	if *scan < 1 || *scan > whyScanMax {
+	if !verbs.WhyScanOK(*a.Scan) {
 		fmt.Fprintf(stderr, "ailang-worldd why: --scan must be 1..%d\n", whyScanMax)
 		return exitUsage
 	}
+	expectStore, resultFile, scan, asJSON := a.ExpectStore, a.Result, a.Scan, a.JSON
 	t, err := parseWhyTarget(target, *resultFile, stdin)
 	if err != nil {
 		fmt.Fprintf(stderr, "ailang-worldd why: %v\n%s", err, whyHelp)
diff --git a/cmd/ailang-worldd/why_seam_test.go b/cmd/ailang-worldd/why_seam_test.go
--- a/cmd/ailang-worldd/why_seam_test.go
+++ b/cmd/ailang-worldd/why_seam_test.go
@@ -12,6 +12,7 @@ import (
 	"strings"
 	"testing"
 
+	"github.com/sunholo-data/ailang-world/host/verbs"
 	"github.com/sunholo-data/ailang-world/host/why"
 )
 
@@ -124,7 +125,8 @@ func TestWhyRunGetterClasses(t *testing.T) {
 // VW4: the moved constants, as literals, so a change on both sides still fails.
 func TestWhyConstantsMatchTheCLI(t *testing.T) {
 	if why.MaxResponseBytes != maxClientResponseBytes || maxClientResponseBytes != 8388608 ||
-		whyScanDefault != 500 || whyScanMax != 5000 || logPageMax != 500 || why.ElideAt != 160 {
+		whyScanDefault != 500 || whyScanMax != 5000 || logPageMax != 500 || why.ElideAt != 160 ||
+		verbs.WhyScanDefault != why.ScanDefault || verbs.WhyScanMax != why.ScanMax {
 		t.Fatalf("constants drifted: %d %d %d %d %d", why.MaxResponseBytes, whyScanDefault, whyScanMax, logPageMax, why.ElideAt)
 	}
 }
diff --git a/host/verbs/bounds_test.go b/host/verbs/bounds_test.go
new file mode 100644
--- /dev/null
+++ b/host/verbs/bounds_test.go
@@ -0,0 +1,38 @@
+package verbs
+
+// Row 243 V2b (design §9, V1 judge notes 2 and 4b): VB6, VB7.
+
+import (
+	"strings"
+	"testing"
+)
+
+// VB6: every never and other word's refusal, bare and with an argument, ends with the CLI line.
+func TestNeverAndOtherReasonsEndWithTheCLILine(t *testing.T) {
+	n := 0
+	for _, w := range append(append([]string{}, NeverWords...), OtherWords...) {
+		for _, argv := range [][]string{{w}, {w, "zzz"}} {
+			v, _ := Lookup(argv)
+			if want := "; at a terminal: ailang-worldd " + strings.Join(argv, " "); v.CLI == "" || !strings.HasSuffix(v.Reason, want) {
+				t.Errorf("Lookup(%q).Reason = %q, want it to end %q", argv, v.Reason, want)
+			}
+			n++
+		}
+	}
+	if n != 48 {
+		t.Fatalf("%d rows, want 48 (24 words, bare and with an argument)", n)
+	}
+}
+
+// VB7: the bar's two input bounds, as literals, so a bound changed with its test still fails.
+func TestSplitBoundsAreLiterals(t *testing.T) {
+	if MaxWords != 32 || MaxLine != 1024 {
+		t.Fatalf("MaxWords %d MaxLine %d, want 32 and 1024", MaxWords, MaxLine)
+	}
+	if _, err := Split(strings.TrimSpace(strings.Repeat("w ", 33))); err == nil {
+		t.Error("33 words split")
+	}
+	if _, err := Split(strings.Repeat("x", 1025)); err != ErrTooLong {
+		t.Errorf("1,025 bytes: %v", err)
+	}
+}
diff --git a/host/verbs/flags.go b/host/verbs/flags.go
--- a/host/verbs/flags.go
+++ b/host/verbs/flags.go
@@ -5,6 +5,7 @@ import (
 	"flag"
 	"fmt"
 	"io"
+	"strings"
 	"time"
 )
 
@@ -75,3 +76,87 @@ func AdmitEpisodeReport(args []string) (episode string, asJSON bool, reason stri
 	}
 	return pos[0], *a.JSON, ""
 }
+
+// WhyScanDefault and WhyScanMax are why's --scan default and bound (host/why's ScanDefault, ScanMax).
+const (
+	WhyScanDefault = 500
+	WhyScanMax     = 5000
+)
+
+// WhyArgs are why's flags.
+type WhyArgs struct {
+	ExpectStore, Result *string
+	Scan                *int
+	JSON                *bool
+}
+
+// WhyFlags registers why's four flags on fs; the CLI and the bar share it.
+func WhyFlags(fs *flag.FlagSet) *WhyArgs {
+	return &WhyArgs{
+		ExpectStore: fs.String("expect-store", "", "refuse unless the daemon serves this store (World-Provenance store= path)"),
+		Result:      fs.String("result", "", "a 'call --json-out' output file"),
+		Scan:        fs.Int("scan", WhyScanDefault, "entries to scan back from head (max 5000)"),
+		JSON:        fs.Bool("json", false, "print the chain as JSON"),
+	}
+}
+
+// WhyLeadingTarget takes a leading target: a first argument that is not a flag, or exactly -.
+func WhyLeadingTarget(args []string) (target string, rest []string) {
+	if len(args) > 0 && (!strings.HasPrefix(args[0], "-") || args[0] == "-") {
+		return args[0], args[1:]
+	}
+	return "", args
+}
+
+// WhyTrailingTarget takes the one positional left after the flags, unless a target was already given.
+func WhyTrailingTarget(fs *flag.FlagSet, target string) (string, error) {
+	if fs.NArg() == 0 {
+		return target, nil
+	}
+	if target != "" || fs.NArg() > 1 {
+		return "", fmt.Errorf("unexpected argument %q", fs.Arg(fs.NArg()-1))
+	}
+	return fs.Arg(0), nil
+}
+
+// WhyScanOK reports whether n is a --scan the CLI runs.
+func WhyScanOK(n int) bool { return n >= 1 && n <= WhyScanMax }
+
+// whyBarRefusal is the bar profile's reason for each of why's reads the app does not make (design §2.3),
+// in the order the bar checks them.
+var whyBarRefusal = []struct{ what, why string }{
+	{"--result", "the app reads no file on the daemon's host"},
+	{"--expect-store", "the app reads its own daemon's store"},
+	{"why -", "the app has no stdin"},
+}
+
+// AdmitWhy runs why's flag set and parse over args in the CLI's order (leading target, flags, trailing
+// target, the --scan bound), then the bar profile: --result, --expect-store and - are refused with the
+// CLI line. The target is then why.ParseTarget's (the route's), as at a terminal.
+func AdmitWhy(args []string) (target string, scan int, asJSON bool, reason string) {
+	cli := Quote(append([]string{Binary, "why"}, args...))
+	target, rest := WhyLeadingTarget(args)
+	fs := flag.NewFlagSet(Binary+" why", flag.ContinueOnError)
+	fs.SetOutput(io.Discard)
+	a := WhyFlags(fs)
+	if err := fs.Parse(rest); err != nil {
+		if errors.Is(err, flag.ErrHelp) {
+			return "", 0, false, "see " + Binary + " why --help at a terminal"
+		}
+		return "", 0, false, err.Error()
+	}
+	target, err := WhyTrailingTarget(fs, target)
+	if err != nil {
+		return "", 0, false, Binary + " why: " + err.Error()
+	}
+	if !WhyScanOK(*a.Scan) {
+		return "", 0, false, fmt.Sprintf("%s why: --scan must be 1..%d", Binary, WhyScanMax)
+	}
+	for i, used := range []bool{*a.Result != "", *a.ExpectStore != "", target == "-"} {
+		if used {
+			r := whyBarRefusal[i]
+			return "", 0, false, r.what + " is not run from the app: " + r.why + "; at a terminal: " + cli
+		}
+	}
+	return target, *a.Scan, *a.JSON, ""
+}
diff --git a/host/verbs/why_admit_test.go b/host/verbs/why_admit_test.go
new file mode 100644
--- /dev/null
+++ b/host/verbs/why_admit_test.go
@@ -0,0 +1,103 @@
+package verbs
+
+// Row 243 V2b (design §2.2 VA1, §2.3 VA2): why's one flag set and parse, and the bar's admission.
+
+import (
+	"flag"
+	"reflect"
+	"strings"
+	"testing"
+)
+
+// VA1: why's flag set exactly, its two positional rules and its --scan bound.
+func TestWhyFlags(t *testing.T) {
+	fs := flag.NewFlagSet("x", flag.ContinueOnError)
+	WhyFlags(fs)
+	got := map[string]string{}
+	fs.VisitAll(func(f *flag.Flag) { got[f.Name] = f.DefValue + " " + f.Usage })
+	if want := map[string]string{"expect-store": " refuse unless the daemon serves this store (World-Provenance store= path)",
+		"result": " a 'call --json-out' output file", "scan": "500 entries to scan back from head (max 5000)",
+		"json": "false print the chain as JSON"}; !reflect.DeepEqual(got, want) {
+		t.Errorf("flags = %v", got)
+	}
+	for _, c := range []struct {
+		args   []string
+		target string
+		rest   int
+	}{{[]string{"3", "--json"}, "3", 1}, {[]string{"-", "--scan", "2"}, "-", 2}, {[]string{"--json", "3"}, "", 2}, {[]string{"-3"}, "", 1}, {nil, "", 0}} {
+		if tg, rest := WhyLeadingTarget(c.args); tg != c.target || len(rest) != c.rest {
+			t.Errorf("WhyLeadingTarget(%q) = %q %q", c.args, tg, rest)
+		}
+	}
+	for _, c := range []struct {
+		args         []string
+		target, want string
+	}{{[]string{"3"}, "", "3"}, {nil, "head", "head"}, {[]string{"2"}, "1", `unexpected argument "2"`}, {[]string{"1", "2"}, "", `unexpected argument "2"`}} {
+		fs := flag.NewFlagSet("x", flag.ContinueOnError)
+		_ = fs.Parse(c.args)
+		tg, err := WhyTrailingTarget(fs, c.target)
+		if (err == nil && tg != c.want) || (err != nil && err.Error() != c.want) {
+			t.Errorf("WhyTrailingTarget(%q, %q) = %q, %v", c.args, c.target, tg, err)
+		}
+	}
+	if WhyScanOK(0) || !WhyScanOK(1) || !WhyScanOK(5000) || WhyScanOK(5001) || WhyScanDefault != 500 || WhyScanMax != 5000 {
+		t.Error("WhyScanOK bound")
+	}
+}
+
+// VA2: one row per argument form of design §2.3, the reason exactly; the target, --scan and --json
+// of every admitted form.
+func TestAdmitWhy(t *testing.T) {
+	const at = "; at a terminal: ailang-worldd why "
+	ref := "sha256:" + strings.Repeat("ab", 32)
+	for _, c := range []struct {
+		line, target string
+		scan         int
+		json         bool
+		reason       string
+	}{
+		{"", "", 500, false, ""}, // no target: the route's why.ParseTarget answers "a target is required"
+		{"3", "3", 500, false, ""},
+		{"head", "head", 500, false, ""},
+		{ref, ref, 500, false, ""},
+		{"sha256:zz", "sha256:zz", 500, false, ""}, // malformed: why.ParseTarget's, at the route
+		{"a2a:e:t", "a2a:e:t", 500, false, ""},
+		{"mcp:e:t", "mcp:e:t", 500, false, ""},
+		{"rest:x", "rest:x", 500, false, ""},
+		{"x", "x", 500, false, ""},
+		{"head --json", "head", 500, true, ""},
+		{"--json head", "head", 500, true, ""},
+		{"a2a:e:t --scan 1", "a2a:e:t", 1, false, ""},
+		{"--scan 5000 a2a:e:t", "a2a:e:t", 5000, false, ""},
+		{"--scan 0 head", "", 0, false, "ailang-worldd why: --scan must be 1..5000"},
+		{"head --scan 5001", "", 0, false, "ailang-worldd why: --scan must be 1..5000"},
+		{"head --scan -1", "", 0, false, "ailang-worldd why: --scan must be 1..5000"},
+		{"1 2", "", 0, false, `ailang-worldd why: unexpected argument "2"`},
+		{"--json 1 2", "", 0, false, `ailang-worldd why: unexpected argument "2"`},
+		{"-3", "", 0, false, "flag provided but not defined: -3"},
+		{"--bogus", "", 0, false, "flag provided but not defined: -bogus"},
+		{"--scan x head", "", 0, false, `invalid value "x" for flag -scan: parse error`},
+		{"-h", "", 0, false, "see ailang-worldd why --help at a terminal"},
+		{"--help", "", 0, false, "see ailang-worldd why --help at a terminal"},
+		{"head --help --result f", "", 0, false, "see ailang-worldd why --help at a terminal"},
+		{"-", "", 0, false, "why - is not run from the app: the app has no stdin" + at + "-"},
+		{"- --scan 3", "", 0, false, "why - is not run from the app: the app has no stdin" + at + "- --scan 3"},
+		{"--result f", "", 0, false, "--result is not run from the app: the app reads no file on the daemon's host" + at + "--result f"},
+		{"--expect-store p head", "", 0, false, "--expect-store is not run from the app: the app reads its own daemon's store" + at + "--expect-store p head"},
+		{"- --result f --expect-store p", "", 0, false, "--result is not run from the app: the app reads no file on the daemon's host" + at + "- --result f --expect-store p"},
+		{"- --expect-store p", "", 0, false, "--expect-store is not run from the app: the app reads its own daemon's store" + at + "- --expect-store p"},
+		{"--result '' head", "head", 500, false, ""}, // an empty --result is no --result, as at a terminal
+	} {
+		args, err := Split(c.line)
+		if err != nil {
+			t.Fatal(err)
+		}
+		target, scan, asJSON, reason := AdmitWhy(args)
+		if reason != c.reason || target != c.target || scan != c.scan || asJSON != c.json {
+			t.Errorf("AdmitWhy(%q) = %q %d %v %q; want %q %d %v %q", c.line, target, scan, asJSON, reason, c.target, c.scan, c.json, c.reason)
+		}
+		if v, _ := Lookup(append([]string{"why"}, args...)); strings.Contains(reason, "; at a terminal: ") && !strings.HasSuffix(reason, "; at a terminal: "+v.CLI) {
+			t.Errorf("AdmitWhy(%q): the CLI line %q is not Lookup's %q", c.line, reason, v.CLI)
+		}
+	}
+}
````

### Patch M6p — the page: A2–A4; VP7, VP8, VP10, VP11; XV1's `log` arm

````diff M6p
diff --git a/cmd/ailang-worldd/verb_bar_chrome_test.go b/cmd/ailang-worldd/verb_bar_chrome_test.go
--- a/cmd/ailang-worldd/verb_bar_chrome_test.go
+++ b/cmd/ailang-worldd/verb_bar_chrome_test.go
@@ -17,7 +17,8 @@ import (
 )
 
 // XV1: `episode report w1` typed into #addr shows the CLI's own bytes in the verb panel and clears the
-// bar; `why` shows its CLI line; `log` is refused in the header with its line, the bar and the place kept.
+// bar; `why` shows its CLI line; `log` is refused in the panel with its CLI line (V2b A3); `world a b c`, a refusal
+// with no CLI line, is said in the header, the bar and the place kept.
 func TestVerbBarEpisodeReportInChrome(t *testing.T) {
 	b := startPasskeyBrowser(t)
 	b.do(t, nil, "tab", map[string]any{})
@@ -68,11 +69,15 @@ func TestVerbBarEpisodeReportInChrome(t *testing.T) {
 	}
 	read()
 	before := state.Location
-	b.do(t, nil, "waitFor", map[string]any{"expression": `(document.getElementById('verb').remove(), 1)`})
 	enter("log")
-	b.do(t, nil, "waitFor", map[string]any{"expression": `document.getElementById('hmsg').textContent.includes('at a terminal: ailang-worldd log')`})
-	if read(); state.Verb || state.Addr != "log" || state.Location != before {
-		t.Fatalf("after log: %+v (was at %s)", state, before)
+	if got := panel("ailang-worldd log"); got != "log is a CLI command the app does not run; at a terminal: ailang-worldd log" {
+		t.Fatalf("log shows %q", got)
+	}
+	b.do(t, nil, "waitFor", map[string]any{"expression": `(document.getElementById('verb').remove(), 1)`})
+	enter("world a b c")
+	b.do(t, nil, "waitFor", map[string]any{"expression": `document.getElementById('hmsg').textContent.includes('unknown address: try')`})
+	if read(); state.Verb || state.Addr != "world a b c" || state.Location != before {
+		t.Fatalf("after world a b c: %+v (was at %s)", state, before)
 	}
 	b.do(t, nil, "openCard", map[string]any{"url": "about:blank", "raw": true}) // ends the page's live stream before the server closes
 }
diff --git a/host/surfacepage/app/app.js b/host/surfacepage/app/app.js
--- a/host/surfacepage/app/app.js
+++ b/host/surfacepage/app/app.js
@@ -492,19 +492,20 @@
     const need = b && typeof b === 'object' && b.schema === 'world/surface-verb/v1' && b.v === 1 && Object.prototype.hasOwnProperty.call(VERB_NEEDS, b.outcome) ? VERB_NEEDS[b.outcome] : null;
     if (!need || !need.every(k => typeof b[k] === 'string' && b[k] !== '')) return UNREADABLE;
     if (b.outcome === 'address') return { address: true };
-    if (b.outcome === 'refused') return { say: b.reason };
+    if (b.outcome === 'refused') return typeof b.cli === 'string' && b.cli !== '' ? { panel: { head: '$ ' + b.cli, body: b.reason, copy: b.cli } } : { say: b.reason };
     if (b.outcome === 'navigate') return b.href.startsWith('/surface/approve/mr-') || b.href.startsWith('/surface/decide/sha256:') ? { assign: b.href } : { say: 'refused: unexpected target' };
     if (b.outcome === 'terminal') return { panel: { head: '$ ' + b.cli, body: b.reason + '\nAt a terminal: ' + b.cli, copy: b.cli } };
     const shown = [b.text, b.json, b.stderr].filter(x => typeof x === 'string' && x !== '');
-    if (shown.length !== 1 || (b.stderr !== undefined && b.exit !== 1)) return UNREADABLE;
-    return { panel: { head: '$ ' + b.cli, body: b.stderr !== undefined ? b.stderr + 'exit 1' : shown[0], copy: '' } };
+    if (shown.length !== 1 || (b.exit !== undefined && b.exit !== 1 && b.exit !== 3) || (b.stderr !== undefined) !== (b.exit === 1)) return UNREADABLE;
+    return { panel: { head: '$ ' + b.cli, body: shown[0] + (b.exit ? 'exit ' + b.exit : ''), copy: '' } };
   }
   // askVerb asks the route under one budget, across the fetch and the body (io carries fetch, the
   // clock and AbortController). A newer ask supersedes an older one: the older is aborted and its
   // answer dropped, never acted on; the newer answers for the bar.
   let verbJob = null;
+  function dropVerb() { if (verbJob) { verbJob.superseded = true; verbJob.ctl.abort(); } }
   async function askVerb(raw, io) {
-    if (verbJob) { verbJob.superseded = true; verbJob.ctl.abort(); }
+    dropVerb();
     const job = { ctl: new io.AbortController(), superseded: false, late: false };
     verbJob = job;
     const stop = new Promise((_, reject) => job.ctl.signal.addEventListener('abort', () => reject(new Error('aborted'))));
@@ -532,6 +533,28 @@
     if (act.panel) { fx.clear(); return fx.panel(act.panel); }
     fx.say(act.say);
   }
+  // goWith is the bar's Enter over fx (design A4): every Enter owns the bar. An empty line supersedes a pending ask and goes
+  // home; an answer that arrives after a newer Enter is dropped, and the address path re-checks after its own await.
+  let barTurn = 0;
+  async function goWith(raw, fx) {
+    const turn = ++barTurn, current = () => turn === barTurn;
+    if (String(raw == null ? '' : raw).trim() === '') { dropVerb(); return fx.out(); }
+    const act = await fx.ask(raw);
+    if (!current()) return;
+    return verbAct(act, raw, { ...fx.act, address: r => fx.address(r, current) });
+  }
+  // addressGo is the address bar's own path (M4a's grammar), reached only through the route's address answer; after the
+  // index loads it acts only for the current Enter.
+  async function addressGo(raw, current, fx) {
+    let r = addressTarget(raw, null, fx.parse);
+    if (r.error && typeof fx.parse === 'function') {
+      try { r = addressTarget(raw, await fx.loadIndex(), fx.parse); } catch (e) { r = { error: 'the surface index could not be read: ' + e.message }; }
+      if (!current()) return;
+    }
+    if (r.home) return fx.out();
+    if (r.open) return fx.openState(r.open.brane, r.open.at);
+    fx.say(r.error);
+  }
   // verbPanel builds the panel's nodes with textContent only: a report or a reason is never markup.
   function verbPanel(doc, p) {
     const el = doc.createElement('div'), head = doc.createElement('div'), body = doc.createElement('pre');
@@ -545,7 +568,7 @@
 
   if (typeof module !== 'undefined') module.exports = { needsYouHTML, decidedHTML, notifications, isSharedKey, sourcesHTML, openRunningHTML, fetchCoherent, serviceLineHTML, splitFrames, applyDelta, nextBackoff, liveRequest, liveStatus, esc, cleanTitle, shortTitle, prettyModel, agentLine, summaryOf, rowLabel, verdict, LEVELS, levelXs, blendXS,
     adaptModel, layoutGraph, stackBars, barHTML, frontHTML, railHTML, appPath, stackPath, parseAppQuery, addressTarget, lookOf, crumbsHTML, geometry: { PL, PR, PW },
-    VERB_ROUTE_BUDGET_MS, verbOutcome, askVerb, verbAct, verbPanel };
+    VERB_ROUTE_BUDGET_MS, verbOutcome, askVerb, verbAct, verbPanel, goWith, addressGo };
   if (typeof document === 'undefined') return;
 
   // ================= The page =================
@@ -989,7 +1012,7 @@
   const addrOf = () => F.kind === 'open' ? String(F.id).replace(/:episode:/, ':') : `${String(F.brane).replace(/^([a-z][a-z0-9-]*):brane:/, '$1:')}@${railIdx() + 1}`;
   function renderCrumbs() {
     $('crumbs').innerHTML = crumbsHTML(st.d === 1 && F ? addrOf() : null);
-    $('addr').placeholder = st.d !== 0 ? '' : (T && T.GROUPS[0].id ? 'go to a place — ' + T.GROUPS.map(g => g.id).join(' · ') + ' · ' + T.GROUPS[0].id + ':row-N@1' : 'go to a place — row-189@6 · shared@0ffbb9c') + ' · episode report <episode>';
+    $('addr').placeholder = st.d !== 0 ? '' : (T && T.GROUPS[0].id ? 'go to a place — ' + T.GROUPS.map(g => g.id).join(' · ') + ' · ' + T.GROUPS[0].id + ':row-N@1' : 'go to a place — row-189@6 · shared@0ffbb9c') + ' · episode report <episode> · why <target>';
   }
   let msgT;
   function say(t) { const m = $('hmsg'); m.textContent = t; m.style.opacity = '1'; clearTimeout(msgT); msgT = setTimeout(() => { m.style.opacity = '0'; }, 2600); }
@@ -1263,11 +1286,10 @@
     openState(F.brane, RAIL[k].id);
   }
   // go is the bar's Enter: empty is home (no fetch); anything else asks the verb route first.
-  async function go(raw) {
-    if (String(raw == null ? '' : raw).trim() === '') return out();
-    const act = await askVerb(raw, { fetch: (u, o) => fetch(u, o), AbortController, setTimeout: (f, ms) => setTimeout(f, ms), clearTimeout: t => clearTimeout(t) });
-    return verbAct(act, raw, { address: goAddress, assign: h => location.assign(h), panel: showVerbPanel, say,
-      clear: () => { if ($('addr').value === raw) $('addr').value = ''; } });
+  function go(raw) {
+    return goWith(raw, { out, address: goAddress,
+      ask: q => askVerb(q, { fetch: (u, o) => fetch(u, o), AbortController, setTimeout: (f, ms) => setTimeout(f, ms), clearTimeout: t => clearTimeout(t) }),
+      act: { assign: h => location.assign(h), panel: showVerbPanel, say, clear: () => { if ($('addr').value === raw) $('addr').value = ''; } } });
   }
   function showVerbPanel(p) {
     const old = $('verb'), el = verbPanel(document, p);
@@ -1275,20 +1297,10 @@
     const b = el.querySelector('button');
     if (b) b.addEventListener('click', () => { try { navigator.clipboard.writeText(b.dataset.copy).then(() => say('copied'), () => say('copy it from the panel')); } catch (_) { say('copy it from the panel'); } });
   }
-  // goAddress is the address bar's own path (M4a's grammar), reached only through the route's address answer.
-  async function goAddress(raw) {
-    const parse = SEAMS.parseAddress;
-    let r = addressTarget(raw, null, parse);
-    if (r.error && typeof parse === 'function') {
-      try { r = addressTarget(raw, await loadIndex(), parse); } catch (e) { r = { error: 'the surface index could not be read: ' + e.message }; }
-    }
-    if (r.home) return out();
-    if (r.open) return openState(r.open.brane, r.open.at);
-    say(r.error);
-  }
+  function goAddress(raw, current) { return addressGo(raw, current, { parse: SEAMS.parseAddress, loadIndex, out, openState, say }); }
 
   // ---------- Popovers ----------
-  const KEYROWS = [['← →', 'back · forward a page', '—'], ['⇧← ⇧→', 'previous · next state', '—'], ['↑ ↓', 'lane above · below', 'back in time · toward now'], ['+ −', 'zoom in · out', '—'], ['End', 'back to now', 'most recent'], ['Home', 'the oldest state', 'earliest'], ['Enter', 'open brane', '—'], ['Esc', '—', 'back to timeline'], ['N', 'what needs you', 'what needs you'], ['/', 'address bar', 'address bar'], ['Enter', 'in the bar: a place, episode report E, approve ID', 'the same'], [',', 'appearance', 'appearance'], ['?', 'these keys', 'these keys']];
+  const KEYROWS = [['← →', 'back · forward a page', '—'], ['⇧← ⇧→', 'previous · next state', '—'], ['↑ ↓', 'lane above · below', 'back in time · toward now'], ['+ −', 'zoom in · out', '—'], ['End', 'back to now', 'most recent'], ['Home', 'the oldest state', 'earliest'], ['Enter', 'open brane', '—'], ['Esc', '—', 'back to timeline'], ['N', 'what needs you', 'what needs you'], ['/', 'address bar', 'address bar'], ['Enter', 'in the bar: a place, episode report E, why T, approve ID', 'the same'], [',', 'appearance', 'appearance'], ['?', 'these keys', 'these keys']];
   function renderKeys() {
     $('keys').innerHTML = `<div class="kgrid"><span></span><span style="font-size:12px;font-weight:700;color:${st.d === 0 ? 'var(--text)' : 'var(--faint)'}">Timeline</span><span style="font-size:12px;font-weight:700;color:${st.d === 1 ? 'var(--text)' : 'var(--faint)'}">Brane</span>` +
       KEYROWS.map(([k, a, b]) => `<span><span class="kbd">${k}</span></span><span>${a}</span><span>${b}</span>`).join('') +
diff --git a/host/surfacepage/app/index.html b/host/surfacepage/app/index.html
--- a/host/surfacepage/app/index.html
+++ b/host/surfacepage/app/index.html
@@ -64,7 +64,7 @@
     <nav class="crumbs" id="crumbs" aria-label="Where you are"></nav>
     <label class="addr">
       <span class="prompt">world&gt;</span>
-      <input id="addr" type="text" aria-label="Go to an address" spellcheck="false" autocomplete="off" placeholder="row-189@6 · shared@0ffbb9c · back · episode report &lt;episode&gt;">
+      <input id="addr" type="text" aria-label="Go to an address" spellcheck="false" autocomplete="off" placeholder="row-189@6 · shared@0ffbb9c · back · episode report &lt;episode&gt; · why &lt;target&gt;">
     </label>
     <span class="hmsg" id="hmsg"></span>
     <button type="button" class="hbtn mono" id="keys-btn" aria-label="Keyboard shortcuts" aria-keyshortcuts="?">?</button>
diff --git a/host/surfacepage/app_test.go b/host/surfacepage/app_test.go
--- a/host/surfacepage/app_test.go
+++ b/host/surfacepage/app_test.go
@@ -173,7 +173,7 @@ const stacks = JSON.parse(fs.readFileSync('stacks.json', 'utf8'));
 const index = JSON.parse(fs.readFileSync('index.json', 'utf8'));
 assert.deepEqual(Object.keys(app).sort(), ['LEVELS', 'adaptModel', 'addressTarget', 'agentLine', 'appPath', 'applyDelta', 'barHTML', 'blendXS', 'cleanTitle', 'crumbsHTML', 'decidedHTML', 'esc', 'fetchCoherent', 'frontHTML',
   'geometry', 'isSharedKey', 'layoutGraph', 'levelXs', 'liveRequest', 'liveStatus', 'lookOf', 'needsYouHTML', 'nextBackoff', 'notifications', 'openRunningHTML', 'parseAppQuery', 'prettyModel', 'railHTML', 'rowLabel', 'serviceLineHTML', 'shortTitle', 'sourcesHTML', 'splitFrames', 'stackBars', 'stackPath', 'summaryOf', 'verdict',
-  'VERB_ROUTE_BUDGET_MS', 'askVerb', 'verbAct', 'verbOutcome', 'verbPanel'].sort());
+  'VERB_ROUTE_BUDGET_MS', 'askVerb', 'verbAct', 'verbOutcome', 'verbPanel', 'goWith', 'addressGo'].sort());
 
 assert.equal(app.serviceLineHTML(undefined), '');
 // Row 201 A3: Needs you is the pending requests, each a link to its card; a refused row shows why, with no link.
diff --git a/host/surfacepage/app_verb_test.go b/host/surfacepage/app_verb_test.go
--- a/host/surfacepage/app_verb_test.go
+++ b/host/surfacepage/app_verb_test.go
@@ -102,6 +102,15 @@ async function run(raw, fetchImpl) {
     assert.deepEqual(act, { panel: { head: '$ ' + B(q).cli, body, copy: '' } }, q); assert.deepEqual(log, [['clear'], ['panel', act.panel]], q);
   }
   assert(B('episode report nobody').stderr.endsWith('has no committed calls\n'));
+  // VP10 (A2): exit 3 rides on text or json, under it; a result's shape is otherwise the route's or it is unreadable: no cli;
+  // exit 3 with stderr; exit 2; exit 1 with text; two members.
+  { const { act } = await run('episode report w1', reply(200, JSON.stringify({ ...B('episode report w1'), exit: 3 })));
+    assert.equal(act.panel.body, B('episode report w1').text + 'exit 3'); }
+  for (const body of [{ ...B('episode report w1'), cli: undefined }, { ...B('episode report nobody'), exit: 3 }, { ...B('episode report w1'), exit: 2 },
+    { ...B('episode report w1'), exit: 1 }, { ...B('episode report nobody'), text: 'x' }]) {
+    const { log } = await run('episode report w1', reply(200, JSON.stringify(body)));
+    assert.deepEqual(log, [['say', UNREAD]], JSON.stringify(body).slice(0, 80));
+  }
   // Terminal: the CLI line, with a copy.
   for (const q of ['why', 'episode finish w1 --preset p']) {
     const { act, log } = await run(q, golden(q)), b = B(q);
@@ -118,13 +127,17 @@ async function run(raw, fetchImpl) {
     const { log } = await run(MR, reply(200, JSON.stringify({ ...B(MR), href })));
     assert.deepEqual(log, [['say', 'refused: unexpected target']], href);
   }
-  // Refused (VP6, the shadow set; TooLong): the reason is said; nothing else happens.
+  // VP11 (A3): a refusal with a CLI line opens the panel with the line to copy; one without (the world 4b shape, TooLong) is said.
   for (const [q, has] of [['log', 'at a terminal: ailang-worldd log'], ['ailang-worldd world get world-1', 'ailang-worldd world get world-1'],
-    ['episode', 'needs a subcommand'], ['world a b c', 'unknown address: try "world"'], [LONG, 'TooLong']]) {
+    ['episode', 'needs a subcommand']]) {
+    const { act, log } = await run(q, golden(q)), b = B(q);
+    assert.deepEqual(act, { panel: { head: '$ ' + b.cli, body: b.reason, copy: b.cli } }, q); assert.deepEqual(log, [['clear'], ['panel', act.panel]], q); assert(b.reason.includes(has), q);
+  }
+  for (const [q, has] of [['world a b c', 'unknown address: try "world"'], [LONG, 'TooLong']]) {
+    assert.deepEqual(app.verbOutcome(200, JSON.stringify({ ...B(q), cli: '' })), { say: B(q).reason }, q);
     const { act, log } = await run(q, golden(q));
-    assert.deepEqual(act, { say: B(q).reason }, q); assert.deepEqual(log, [['say', B(q).reason]], q); assert(B(q).reason.includes(has), q);
+    assert.deepEqual(act, { say: B(q).reason }, q); assert.deepEqual(log, [['say', B(q).reason]], q); assert(B(q).reason.includes(has), q); assert.equal(B(q).cli, undefined, q);
   }
-  assert(B(LONG).reason.startsWith('TooLong'));
   // VP2: no silent fallback. Every failure says its message and calls nothing else, for verbs and addresses alike.
   for (const q of ['why', 'log', 'row-189@6', 'world shared@0ffbb9c']) {
     const b = B(q), drop = k => { const c = { ...b }; delete c[k]; return JSON.stringify(c); };
@@ -156,10 +169,37 @@ async function run(raw, fetchImpl) {
   assert(io1.calls[0].o.signal.aborted);
   release({ status: 200, text: () => Promise.resolve(G['row-189@6']) });
   const a1 = await p1, a2 = await p2;
-  assert.deepEqual(a1, { superseded: true }); assert.deepEqual(io1.cleared, [1]); assert.deepEqual(a2, { say: B('log').reason });
+  assert.deepEqual(a1, { superseded: true }); assert.deepEqual(io1.cleared, [1]); assert.deepEqual(a2, { panel: { head: '$ ailang-worldd log', body: B('log').reason, copy: 'ailang-worldd log' } });
   { const fx = spies(); app.verbAct(a1, 'row-189@6', fx); assert.deepEqual(fx.log, []); }
   const io3 = rig(golden('row-189@6')), p3 = app.askVerb('row-189@6', io3), p4 = app.askVerb('why', rig(golden('why')));
   assert.deepEqual(await p3, { superseded: true }); assert.deepEqual(io3.cleared, [1]); assert(Boolean((await p4).panel));
+  // VP8 (X-JS-FINALLY-NULL): A pending, B supersedes A, A's finally runs; C must still abort B.
+  { const ioA = rig(() => new Promise(() => {})), ioB = rig(() => new Promise(() => {}));
+    const pA = app.askVerb('a', ioA), pB = app.askVerb('b', ioB);
+    assert.deepEqual(await pA, { superseded: true });
+    const pC = app.askVerb('log', rig(golden('log')));
+    assert(ioB.calls[0].o.signal.aborted, 'C aborts B'); assert.deepEqual(await pB, { superseded: true }); assert(Boolean((await pC).panel)); }
+  // VP7 (A4): every Enter owns the bar. An empty Enter aborts a pending ask and goes home, and that ask's answer does nothing;
+  // an answer arriving after a newer Enter is dropped; the address path does not act after a newer Enter.
+  { const log = [], fx = { out: () => log.push('out'), address: r => log.push('address ' + r), act: spies() };
+    let io; fx.ask = q => app.askVerb(q, io = rig(() => new Promise(() => {})));
+    const p = app.goWith('row-189@6', fx);
+    await app.goWith('  ', fx);
+    assert(io.calls[0].o.signal.aborted, 'the empty Enter aborts the pending ask'); await p;
+    let release; fx.ask = () => new Promise(r => { release = r; });
+    const p2 = app.goWith('row-189@6', fx);
+    await app.goWith('', fx); release({ address: true }); await p2;
+    assert.deepEqual(log, ['out', 'out']); assert.deepEqual(fx.act.log, []);
+    let cur; fx.address = (r, c) => { log.push('address ' + r); cur = c; };
+    fx.ask = () => Promise.resolve({ address: true }); await app.goWith('row-189@6', fx);
+    assert.deepEqual(log, ['out', 'out', 'address row-189@6']); assert.deepEqual(fx.act.log, [['clear']]);
+    assert.equal(cur(), true); await app.goWith('', fx); assert.equal(cur(), false, 'a newer Enter ends the address path\'s turn'); }
+  for (const [newer, want] of [[true, []], [false, [['openState', 'row-189', '652c341d5072c5f49fefe29f9409c93a5119a5be']]]]) {
+    const calls = [], ax = { parse: seams.parseAddress, loadIndex: () => Promise.resolve(index) };
+    for (const k of ['out', 'openState', 'say']) ax[k] = (...a) => calls.push([k, ...a]);
+    await app.addressGo('row-189@6', () => !newer, ax);
+    assert.deepEqual(calls, want, 'a newer Enter: ' + newer);
+  }
   // VP3: the panel is built with textContent only; a hostile report or reason stays text.
   const fakeDoc = () => { const doc = { made: [], createElement(tag) {
     const el = { tag, children: [], dataset: {}, style: {}, attrs: {}, textContent: '', setAttribute(k, v) { this.attrs[k] = v; }, append(...c) { this.children.push(...c); } };
````

### Patch M5g — the getter, `verbWhy`, `whyResult`, the A2 cap; VR5, VR7, VR8, VS2

````diff M5g
diff --git a/host/daemon/surface_verb.go b/host/daemon/surface_verb.go
--- a/host/daemon/surface_verb.go
+++ b/host/daemon/surface_verb.go
@@ -42,7 +42,7 @@ func (d *Daemon) handleSurfaceVerb(w http.ResponseWriter, r *http.Request, _ []s
 		return
 	}
 	a := verbAnswer{V: 1, Schema: VerbSchema, Q: query.Get("q"), Argv: []string{}, Class: string(verbs.Refused), Outcome: "refused"}
-	defer func() { writeApprovalJSON(w, http.StatusOK, a) }()
+	defer func() { capVerbAnswer(&a); writeApprovalJSON(w, http.StatusOK, a) }()
 	argv, err := verbs.Split(a.Q) // its first check is the 1,024-byte bound: no oversized line is tokenized
 	if err != nil {
 		a.Reason = err.Error()
diff --git a/host/daemon/surface_verb_why.go b/host/daemon/surface_verb_why.go
new file mode 100644
--- /dev/null
+++ b/host/daemon/surface_verb_why.go
@@ -0,0 +1,163 @@
+package daemon
+
+// Row 243 V2b (design_docs/planned/w-world-address-bar-verbs-v2.md §2.4, §2.5): why in the bar. The
+// walk is host/why's, the CLI's own; it reads World through a private mux over the seven /v1 GET
+// handlers it needs, in-process. No route, POST, store write or process is added.
+
+import (
+	"bytes"
+	"context"
+	"errors"
+	"fmt"
+	"net/http"
+
+	"github.com/sunholo-data/ailang-world/host/verbs"
+	"github.com/sunholo-data/ailang-world/host/why"
+)
+
+// MaxVerbAnswerBytes caps a result answer's text, json and stderr together (design A2): 1 MiB of raw
+// bytes. The JSON-encoded body can be larger (a control byte or <, >, & escapes to six bytes, \u00XX)
+// and is still bounded, by six times the cap plus the answer's other members.
+const MaxVerbAnswerBytes = 1 << 20
+
+// verbAnswerCap is the cap the route applies; a test lowers it to drive the refusal end to end.
+var verbAnswerCap = MaxVerbAnswerBytes
+
+// capVerbAnswer turns a result over the cap into a refusal with its size and the CLI line (design §2.5 e).
+func capVerbAnswer(a *verbAnswer) {
+	if n := len(a.Text) + len(a.JSON) + len(a.Stderr); a.Outcome == "result" && n > verbAnswerCap {
+		a.Outcome, a.Text, a.JSON, a.Stderr, a.Exit = "refused", "", "", "", 0
+		a.Reason = fmt.Sprintf("the answer is %d bytes, over the app's 1 MiB cap; at a terminal: %s", n, a.CLI)
+	}
+}
+
+// whyRoutes is the seven /v1 GET patterns why's walk reads, each to the method value Handler registers.
+func (d *Daemon) whyRoutes() map[string]http.HandlerFunc {
+	return map[string]http.HandlerFunc{
+		"GET /v1/health":        d.handleHealth,
+		"GET /v1/head":          d.handleHead,
+		"GET /v1/worlds/{ref}":  d.handleWorld,
+		"GET /v1/objects/{ref}": d.handleObject,
+		"GET /v1/log/{index}":   d.handleLogEntry,
+		"GET /v1/log":           d.handleLogRange,
+		"GET /v1/receipts/{id}": d.handleReceipt,
+	}
+}
+
+// whyMux is a private mux over exactly whyRoutes, built per walk.
+func (d *Daemon) whyMux() *http.ServeMux {
+	m := http.NewServeMux()
+	for pattern, h := range d.whyRoutes() {
+		m.HandleFunc(pattern, h)
+	}
+	return m
+}
+
+// whyTally is what the walk's reads saw: a read deadline (503) or another failure.
+type whyTally struct{ timeout, failed bool }
+
+// boundedRecorder is one in-process response, its body capped at max bytes.
+type boundedRecorder struct {
+	header http.Header
+	status int
+	body   bytes.Buffer
+	max    int
+	over   bool
+}
+
+func (b *boundedRecorder) Header() http.Header { return b.header }
+func (b *boundedRecorder) WriteHeader(status int) {
+	if b.status == 0 {
+		b.status = status
+	}
+}
+func (b *boundedRecorder) Write(p []byte) (int, error) {
+	b.WriteHeader(http.StatusOK)
+	if b.body.Len()+len(p) > b.max {
+		b.over = true
+		return 0, errors.New("response over the in-process cap")
+	}
+	return b.body.Write(p)
+}
+
+// whyGetter is why's Getter over mux, in-process: a 404 is an answer; an unrouted path, a request that
+// cannot be built or a body over why.MaxResponseBytes is an error, logged for the operator. tally
+// records each 503 (the handlers' read deadline) and each other failure.
+func (d *Daemon) whyGetter(r *http.Request, mux *http.ServeMux, tally *whyTally) why.Getter {
+	return func(ctx context.Context, path string) (int, []byte, error) {
+		req, err := http.NewRequestWithContext(ctx, http.MethodGet, path, nil)
+		if err == nil {
+			if _, pattern := mux.Handler(req); pattern == "" {
+				err = fmt.Errorf("unrouted in-process read: %s", path)
+			}
+		}
+		var rec *boundedRecorder
+		if err == nil {
+			rec = &boundedRecorder{header: http.Header{}, max: why.MaxResponseBytes}
+			mux.ServeHTTP(rec, req)
+			if rec.over {
+				err = fmt.Errorf("response from %s exceeds %d bytes", path, why.MaxResponseBytes)
+			}
+		}
+		if err != nil {
+			tally.failed = true
+			d.writeInternalErrorLog(r, err)
+			return 0, nil, err
+		}
+		switch {
+		case rec.status == http.StatusServiceUnavailable:
+			tally.timeout = true
+		case rec.status >= 500:
+			tally.failed = true
+		}
+		return rec.status, rec.body.Bytes(), nil
+	}
+}
+
+// verbWhy runs why for the bar (design §2.5): the CLI's parse and the bar's admission, then the walk
+// under the daemon's read deadline, then the outcome by precedence.
+func (d *Daemon) verbWhy(r *http.Request, args []string, a *verbAnswer) {
+	target, scan, asJSON, reason := verbs.AdmitWhy(args)
+	if reason != "" {
+		a.Reason = reason
+		return
+	}
+	t, err := why.ParseTarget(target)
+	if err != nil {
+		a.Reason = "ailang-worldd why: " + err.Error()
+		return
+	}
+	ctx, cancel := d.readCtx(r)
+	defer cancel()
+	var tally whyTally
+	ans := why.Run(why.Reader{Ctx: ctx, Get: d.whyGetter(r, d.whyMux(), &tally)}, t, why.Options{Scan: scan, JSON: asJSON})
+	switch {
+	case r.Context().Err() != nil: // (a) the page went away: nobody reads this
+		a.Reason = "the request ended before the walk did"
+	case errors.Is(ctx.Err(), context.DeadlineExceeded) || tally.timeout: // (b)
+		a.Reason = fmt.Sprintf("read deadline (%s) exceeded: the walk did not finish in the app; at a terminal it has 60 s: %s", d.readDeadline, a.CLI)
+	case tally.failed: // (c)
+		a.Reason = "the daemon could not read World during the walk, so no chain is shown; at a terminal: " + a.CLI
+	default: // (d), (f)
+		whyResult(a, ans, asJSON)
+	}
+}
+
+// whyResult is the walk's answer as the bar's result (design §2.5 f), or (d) a refusal when it is not
+// exactly one of the CLI's three shapes: stdout with exit 0 or 3, stderr with exit 1.
+func whyResult(a *verbAnswer, ans why.Answer, asJSON bool) {
+	out, errOut := ans.Stdout != "", ans.Stderr != ""
+	if out == errOut || (out && ans.Exit != 0 && ans.Exit != 3) || (errOut && ans.Exit != 1) {
+		a.Reason = "the walk answered something the app cannot show; at a terminal: " + a.CLI
+		return
+	}
+	a.Outcome, a.Exit = "result", ans.Exit
+	switch {
+	case errOut:
+		a.Stderr = ans.Stderr
+	case asJSON:
+		a.JSON = ans.Stdout
+	default:
+		a.Text = ans.Stdout
+	}
+}
diff --git a/host/daemon/surface_verb_why_test.go b/host/daemon/surface_verb_why_test.go
new file mode 100644
--- /dev/null
+++ b/host/daemon/surface_verb_why_test.go
@@ -0,0 +1,234 @@
+package daemon
+
+// Row 243 V2b (design §2.4, §2.5 e, §6 VR5, VR7, VR8, VS2): why's in-process getter and the A2 cap.
+
+import (
+	"bytes"
+	"fmt"
+	"go/ast"
+	"go/parser"
+	"go/token"
+	"net/http"
+	"net/http/httptest"
+	"os"
+	"path/filepath"
+	"reflect"
+	"regexp"
+	"runtime"
+	"strconv"
+	"strings"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/why"
+)
+
+func methodName(h http.HandlerFunc) string {
+	n := runtime.FuncForPC(reflect.ValueOf(h).Pointer()).Name()
+	return strings.TrimSuffix(n[strings.LastIndex(n, ".")+1:], "-fm")
+}
+
+// VR5: a path outside the seven is an error (never a 404, never another handler); a 404 is an answer;
+// the body cap is 8 MiB exactly; a 503 is a timeout and another 5xx a failure.
+func TestWhyGetterUnroutedAndOversize(t *testing.T) {
+	d := newHandlerDaemon(t)
+	var log bytes.Buffer
+	d.errLog = &log
+	ctx := boundedTestContext(t)
+	r := httptest.NewRequest("GET", "/surface/verb.json?q=why", nil)
+	for _, p := range []string{"/v1/session/self", "/v1/nowhere", "/v1/registry/x", "/v1/objects/by-semantic-id/x", "/v1/episodes/w1/report"} {
+		var tally whyTally
+		if st, body, err := d.whyGetter(r, d.whyMux(), &tally)(ctx, p); err == nil || err.Error() != "unrouted in-process read: "+p || st != 0 || body != nil || !tally.failed || tally.timeout {
+			t.Errorf("%s = %d %q %v %+v", p, st, body, err, tally)
+		}
+	}
+	if !strings.Contains(log.String(), "internal error: GET /surface/verb.json: unrouted in-process read: /v1/session/self") {
+		t.Errorf("not logged: %q", log.String())
+	}
+	var tally whyTally
+	if _, _, err := d.whyGetter(r, d.whyMux(), &tally)(ctx, "/v1/log/%zz"); err == nil || !tally.failed {
+		t.Errorf("an unparseable path = %v %+v", err, tally)
+	}
+	tally = whyTally{}
+	if st, _, err := d.whyGetter(r, d.whyMux(), &tally)(ctx, "/v1/objects/sha256:"+strings.Repeat("0", 64)+"?payload=true"); st != 404 || err != nil || tally != (whyTally{}) {
+		t.Errorf("absent object = %d %v %+v", st, err, tally)
+	}
+	stub := http.NewServeMux()
+	for _, c := range []struct {
+		n      int
+		status int
+	}{{8388608, 200}, {8388609, 0}} {
+		n := c.n
+		stub.HandleFunc(fmt.Sprintf("GET /big/%d", n), func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write(bytes.Repeat([]byte("x"), n)) })
+		var tally whyTally
+		st, body, err := d.whyGetter(r, stub, &tally)(ctx, fmt.Sprintf("/big/%d", n))
+		if c.status == 200 && (st != 200 || len(body) != n || err != nil || tally.failed) {
+			t.Errorf("%d bytes = %d %d %v", n, st, len(body), err)
+		}
+		if c.status == 0 && (err == nil || err.Error() != "response from /big/8388609 exceeds 8388608 bytes" || !tally.failed) {
+			t.Errorf("%d bytes = %d %v %+v", n, st, err, tally)
+		}
+	}
+	for status, want := range map[int]whyTally{503: {timeout: true}, 500: {failed: true}, 502: {failed: true}, 404: {}, 400: {}} {
+		code := status
+		stub.HandleFunc(fmt.Sprintf("GET /s/%d", code), func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(code) })
+		var tally whyTally
+		if st, _, err := d.whyGetter(r, stub, &tally)(ctx, fmt.Sprintf("/s/%d", code)); st != code || err != nil || tally != want {
+			t.Errorf("status %d = %d %v %+v", code, st, err, tally)
+		}
+	}
+}
+
+// VR7 (unit): the A2 cap on text+json+stderr together, at 1 MiB exactly; members cleared; only a result
+// is capped. The JSON-encoded body of an answer at the cap is larger than the cap and bounded by 6×.
+func TestVerbAnswerCap(t *testing.T) {
+	third := strings.Repeat("y", 1048576/2)
+	for _, c := range []struct {
+		a    verbAnswer
+		over bool
+	}{
+		{verbAnswer{Outcome: "result", Text: strings.Repeat("x", 1048576)}, false},
+		{verbAnswer{Outcome: "result", Text: strings.Repeat("x", 1048577)}, true},
+		{verbAnswer{Outcome: "result", JSON: strings.Repeat("x", 1048577)}, true},
+		{verbAnswer{Outcome: "result", Stderr: strings.Repeat("x", 1048576) + "\n", Exit: 1}, true},
+		{verbAnswer{Outcome: "result", Text: third, JSON: third, Stderr: "z", Exit: 3}, true},
+		{verbAnswer{Outcome: "result", Text: third, JSON: third}, false},
+		{verbAnswer{Outcome: "refused", Reason: "r", Text: strings.Repeat("x", 1048577)}, false},
+	} {
+		a := c.a
+		a.CLI = "ailang-worldd why 1"
+		n := len(a.Text) + len(a.JSON) + len(a.Stderr)
+		capVerbAnswer(&a)
+		if over := a.Outcome == "refused" && a.Reason == fmt.Sprintf("the answer is %d bytes, over the app's 1 MiB cap; at a terminal: ailang-worldd why 1", n); over != c.over {
+			t.Errorf("%d bytes: %s %q", n, a.Outcome, a.Reason)
+		}
+		if c.over && (a.Text != "" || a.JSON != "" || a.Stderr != "" || a.Exit != 0) {
+			t.Errorf("%d bytes: members left on the refusal", n)
+		}
+		if !c.over && c.a.Outcome == "result" && (a.Outcome != "result" || a.Text != c.a.Text) {
+			t.Errorf("%d bytes refused", n)
+		}
+	}
+	w := httptest.NewRecorder()
+	writeApprovalJSON(w, 200, verbAnswer{Outcome: "result", Text: strings.Repeat("\x01", 1048576)})
+	if n := w.Body.Len(); n <= 1048576 || n > 6*1048576+1024 {
+		t.Errorf("the encoded body is %d bytes", n)
+	}
+}
+
+// VR8: no handler the private mux dispatches reads middleware state (a context value), so the
+// in-process path cannot silently read nil where HTTP would read something; control: approvalOrigins's
+// LocalAddrContextKey read is found by the same instrument.
+func TestWhyHandlersReadNoMiddlewareState(t *testing.T) {
+	fset := token.NewFileSet()
+	files, _ := filepath.Glob("*.go")
+	bodies := map[string]*ast.FuncDecl{}
+	for _, f := range files {
+		if strings.HasSuffix(f, "_test.go") {
+			continue
+		}
+		af, err := parser.ParseFile(fset, f, nil, 0)
+		if err != nil {
+			t.Fatal(err)
+		}
+		for _, d := range af.Decls {
+			if fd, ok := d.(*ast.FuncDecl); ok && fd.Recv != nil {
+				bodies[fd.Name.Name] = fd
+			}
+		}
+	}
+	readsState := func(fd *ast.FuncDecl) bool {
+		found := false
+		ast.Inspect(fd, func(n ast.Node) bool {
+			switch x := n.(type) {
+			case *ast.SelectorExpr:
+				found = found || x.Sel.Name == "Value"
+			case *ast.Ident:
+				found = found || x.Name == "activityRequestKey"
+			}
+			return true
+		})
+		return found
+	}
+	if fd := bodies["approvalOrigins"]; fd == nil || !readsState(fd) {
+		t.Fatal("control: approvalOrigins's context read is not found")
+	}
+	d := newHandlerDaemon(t)
+	for p, h := range d.whyRoutes() {
+		if fd := bodies[methodName(h)]; fd == nil || readsState(fd) {
+			t.Errorf("%s: %s reads middleware state (or is not found)", p, methodName(h))
+		}
+	}
+}
+
+// VS2: surface_verb_why.go's imports and d.* reach are exact sets: the seven read handlers (VR4/VR6
+// prove which), the read deadline, the error log and its own three helpers; each pattern is bound to the
+// method Handler registers for it. No store, credentials store, process or write.
+func TestSurfaceVerbWhyIsReadOnly(t *testing.T) {
+	f, err := parser.ParseFile(token.NewFileSet(), "surface_verb_why.go", nil, 0)
+	if err != nil {
+		t.Fatal(err)
+	}
+	ims := map[string]bool{}
+	for _, im := range f.Imports {
+		p, _ := strconv.Unquote(im.Path.Value)
+		ims[p] = true
+	}
+	if want := map[string]bool{"bytes": true, "context": true, "errors": true, "fmt": true, "net/http": true,
+		"github.com/sunholo-data/ailang-world/host/verbs": true, "github.com/sunholo-data/ailang-world/host/why": true}; !reflect.DeepEqual(ims, want) {
+		t.Errorf("surface_verb_why.go imports %v", ims)
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
+	want := map[string]bool{"handleHealth": true, "handleHead": true, "handleWorld": true, "handleObject": true, "handleLogEntry": true,
+		"handleLogRange": true, "handleReceipt": true, "whyRoutes": true, "whyMux": true, "whyGetter": true, "readCtx": true,
+		"readDeadline": true, "writeInternalErrorLog": true}
+	if !reflect.DeepEqual(calls, want) {
+		t.Errorf("surface_verb_why.go reaches d.%v, want exactly %v", calls, want)
+	}
+	d := newHandlerDaemon(t)
+	src, err := os.ReadFile("daemon.go")
+	if err != nil {
+		t.Fatal(err)
+	}
+	for p, h := range d.whyRoutes() {
+		if !regexp.MustCompile(`mux\.HandleFunc\("` + regexp.QuoteMeta(p) + `", d\.` + methodName(h) + `\)`).Match(src) {
+			t.Errorf("%s is bound to %s, which Handler does not register for it", p, methodName(h))
+		}
+	}
+}
+
+// VR1 (unit, design §2.5 d and f): the walk's answer is shown only in one of the CLI's three shapes —
+// stdout with exit 0 or 3 (text, or json with --json), stderr with exit 1; anything else is (d).
+func TestWhyResult(t *testing.T) {
+	const d = "the walk answered something the app cannot show; at a terminal: ailang-worldd why 1"
+	for _, c := range []struct {
+		ans    why.Answer
+		asJSON bool
+		want   verbAnswer
+	}{
+		{why.Answer{Stdout: "o\n"}, false, verbAnswer{Outcome: "result", Text: "o\n"}},
+		{why.Answer{Stdout: "o\n", Exit: 3}, false, verbAnswer{Outcome: "result", Text: "o\n", Exit: 3}},
+		{why.Answer{Stdout: "{}\n", Exit: 3}, true, verbAnswer{Outcome: "result", JSON: "{}\n", Exit: 3}},
+		{why.Answer{Stderr: "e\n", Exit: 1}, true, verbAnswer{Outcome: "result", Stderr: "e\n", Exit: 1}},
+		{why.Answer{}, false, verbAnswer{Outcome: "refused", Reason: d}},
+		{why.Answer{Stdout: "o", Stderr: "e", Exit: 1}, false, verbAnswer{Outcome: "refused", Reason: d}},
+		{why.Answer{Stdout: "o", Exit: 1}, false, verbAnswer{Outcome: "refused", Reason: d}},
+		{why.Answer{Stdout: "o", Exit: 2}, false, verbAnswer{Outcome: "refused", Reason: d}},
+		{why.Answer{Stderr: "e", Exit: 3}, false, verbAnswer{Outcome: "refused", Reason: d}},
+		{why.Answer{Stderr: "e"}, false, verbAnswer{Outcome: "refused", Reason: d}},
+	} {
+		a := verbAnswer{Outcome: "refused", CLI: "ailang-worldd why 1"}
+		whyResult(&a, c.ans, c.asJSON)
+		c.want.CLI = a.CLI
+		if !reflect.DeepEqual(a, c.want) {
+			t.Errorf("%+v json=%v: %+v, want %+v", c.ans, c.asJSON, a, c.want)
+		}
+	}
+}
````

### Patch M5r — the route's `why`; Lookup's row; VR1–VR4, VR6, VR7 e2e; the golden; VP9; XV1's `why` arm

````diff M5r
diff --git a/cmd/ailang-worldd/verb_bar_chrome_test.go b/cmd/ailang-worldd/verb_bar_chrome_test.go
--- a/cmd/ailang-worldd/verb_bar_chrome_test.go
+++ b/cmd/ailang-worldd/verb_bar_chrome_test.go
@@ -17,7 +17,7 @@ import (
 )
 
 // XV1: `episode report w1` typed into #addr shows the CLI's own bytes in the verb panel and clears the
-// bar; `why` shows its CLI line; `log` is refused in the panel with its CLI line (V2b A3); `world a b c`, a refusal
+// bar; `why` and `log` are refused in the panel with their CLI lines (V2b A3); `world a b c`, a refusal
 // with no CLI line, is said in the header, the bar and the place kept.
 func TestVerbBarEpisodeReportInChrome(t *testing.T) {
 	b := startPasskeyBrowser(t)
@@ -64,7 +64,7 @@ func TestVerbBarEpisodeReportInChrome(t *testing.T) {
 		t.Fatalf("after episode report: %+v", state)
 	}
 	enter("why")
-	if got := panel("At a terminal: ailang-worldd why"); !strings.HasPrefix(got, "why reads World") {
+	if got := panel("a target is required"); got != "ailang-worldd why: a target is required" {
 		t.Fatalf("why shows %q", got)
 	}
 	read()
diff --git a/host/daemon/surface_verb.go b/host/daemon/surface_verb.go
--- a/host/daemon/surface_verb.go
+++ b/host/daemon/surface_verb.go
@@ -33,7 +33,7 @@ type verbAnswer struct {
 }
 
 // handleSurfaceVerb is GET /surface/verb.json?q=<line>, behind approvalRoute: the bar's line through
-// host/verbs' one grammar. It reads (episode report, the /v1 route's own builder) and writes nothing:
+// host/verbs' one grammar. It reads (episode report, the /v1 route's own builder; why, host/why's walk) and writes nothing:
 // no store write, no credentials store, no process. An operator's typo is a 200 refusal.
 func (d *Daemon) handleSurfaceVerb(w http.ResponseWriter, r *http.Request, _ []string) {
 	query := r.URL.Query()
@@ -55,6 +55,8 @@ func (d *Daemon) handleSurfaceVerb(w http.ResponseWriter, r *http.Request, _ []s
 		a.Outcome = "address"
 	case v.Words == "episode report":
 		d.verbEpisodeReport(r, args, &a)
+	case v.Words == "why":
+		d.verbWhy(r, args, &a)
 	case v.Class == verbs.Navigate:
 		if a.Href, a.Reason = verbs.ApproveTarget(args); a.Href != "" {
 			a.Outcome = "navigate"
diff --git a/host/daemon/surface_verb_outcomes_test.go b/host/daemon/surface_verb_outcomes_test.go
--- a/host/daemon/surface_verb_outcomes_test.go
+++ b/host/daemon/surface_verb_outcomes_test.go
@@ -29,7 +29,7 @@ func TestSurfaceVerbOutcomes(t *testing.T) {
 		{"world session", "address", "address", "", "", "", ""},
 		{"world get", "address", "address", "", "", "", ""},
 		{"world episode", "address", "address", "", "", "", ""},
-		{"why", "terminal", "read", "why", "ailang-worldd why", "", "why reads World"},
+		{"why", "refused", "read", "why", "ailang-worldd why", "", "ailang-worldd why: a target is required"},
 		{"episode finish w1 --preset p", "terminal", "mutate", "episode finish", "ailang-worldd episode finish w1 --preset p", "", "changes World"},
 		{"episode start w1", "terminal", "mutate", "episode start", "ailang-worldd episode start w1", "", "changes World"},
 		{"episode launch w1 --preset p", "terminal", "host", "episode launch", "ailang-worldd episode launch w1 --preset p", "", "needs a host process"},
@@ -76,7 +76,7 @@ func TestSurfaceVerbOutcomes(t *testing.T) {
 	for _, q := range []string{"row-189@6", "world shared@0ffbb9c", "back", "world why", "why", "log", "world session",
 		"ailang-worldd world get world-1", "episode", "world episode", "world a b c", "episode finish w1 --preset p",
 		"approve mr-0123456789abcdef", "approve " + ref, "episode report nobody", "episode report w1", "episode report w1 --json",
-		strings.Repeat("x", verbs.MaxLine+1)} {
+		strings.Repeat("x", verbs.MaxLine+1), "why 2", "why head --json", "why 9", "why --result f"} {
 		_, raw := verbGet(a, q)
 		golden[q] = raw
 	}
diff --git a/host/daemon/surface_verb_test.go b/host/daemon/surface_verb_test.go
--- a/host/daemon/surface_verb_test.go
+++ b/host/daemon/surface_verb_test.go
@@ -101,7 +101,7 @@ func TestSurfaceVerbBoundsAndErrors(t *testing.T) {
 	if v, _ := verbGet(a, over); v.Outcome != "refused" || v.Reason != verbs.ErrTooLong.Error() || len(v.Argv) != 0 || v.Class != "refused" {
 		t.Errorf("1,025 bytes = %+v", *v)
 	}
-	if v, _ := verbGet(a, over[:verbs.MaxLine]); v.Outcome != "terminal" || len(v.Argv) != 2 {
+	if v, _ := verbGet(a, over[:verbs.MaxLine]); v.Outcome != "refused" || v.Reason == verbs.ErrTooLong.Error() || len(v.Argv) != 2 {
 		t.Errorf("1,024 bytes = %+v", *v)
 	}
 	if v, _ := verbGet(a, `why "a`); v.Outcome != "refused" || v.Reason != verbs.ErrUnclosed.Error() {
@@ -176,7 +176,7 @@ func TestSurfaceVerbIsReadOnly(t *testing.T) {
 		}
 		return true
 	})
-	if want := map[string]bool{"readCtx": true, "episodeReport": true, "writeInternalErrorLog": true, "verbEpisodeReport": true, "readDeadline": true}; !reflect.DeepEqual(calls, want) {
+	if want := map[string]bool{"readCtx": true, "episodeReport": true, "writeInternalErrorLog": true, "verbEpisodeReport": true, "readDeadline": true, "verbWhy": true}; !reflect.DeepEqual(calls, want) {
 		t.Errorf("surface_verb.go reaches d.%v, want exactly %v", calls, want)
 	}
 	for _, id := range []string{"mr-0123456789abcdef", "mr-0123456789ABCDEF", "mr-0123456789abcde", "mr-0123456789abcdef0", "xr-0123456789abcdef", "mr-012345678 abcdef"} {
diff --git a/host/daemon/surface_verb_why_route_test.go b/host/daemon/surface_verb_why_route_test.go
new file mode 100644
--- /dev/null
+++ b/host/daemon/surface_verb_why_route_test.go
@@ -0,0 +1,285 @@
+package daemon
+
+// Row 243 V2b (design §2.5, §2.7, §6 VR1–VR4, VR6, VR7): why in the bar, end to end through the route.
+
+import (
+	"bytes"
+	"context"
+	"encoding/json"
+	"errors"
+	"fmt"
+	"io"
+	"net"
+	"net/http"
+	"net/http/httptest"
+	"reflect"
+	"sort"
+	"strconv"
+	"strings"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/hashref"
+	"github.com/sunholo-data/ailang-world/host/store"
+	"github.com/sunholo-data/ailang-world/host/why"
+)
+
+// whyRig is verbRig plus a committed REST call (a resolved rest: receipt) and a REST intent that never commits.
+func whyRig(t *testing.T) (approvalRig, *reportSeed) {
+	t.Helper()
+	a, s := verbRig(t)
+	s.call("rest", "w1", "3", "ailang-read", read(50, map[string]any{"ok": true, "content": "y"}), nil)
+	s.pending("rest", "w1", "9", "ailang-read")
+	return a, s
+}
+
+// whyCLI is the CLI line Lookup gives a why line.
+func whyCLI(q string) string { return "ailang-worldd " + q }
+
+// VR1: every outcome of design §2.7 end to end through approvalRoute; a result carries no reason
+// (Lookup's V1 "cannot run it yet" text must not ride on an answer); a refusal carries the CLI line.
+func TestSurfaceVerbWhyOutcomes(t *testing.T) {
+	a, s := whyRig(t)
+	zero := "sha256:" + strings.Repeat("0", 64)
+	for _, c := range []struct {
+		q, member, prefix, suffix string
+		exit                      int
+	}{
+		{"why 1", "text", "why: entry 1 (matched by entry index)\n  ✗ world ", " link(s) BROKEN — integrity refusal (exit 3)\n", 3},
+		{"why head --json", "json", "{\n  \"entry\": 3,\n  \"matchedBy\": \"head\",\n", "\n}\n", 3},
+		{"why rest:w1:3", "text", "why: entry 3 (matched by rest receipt → world)\n", "(exit 3)\n", 3},
+		{"why " + s.head.Ref.String(), "text", "why: entry 3 (matched by world ref)\n", "(exit 3)\n", 3},
+		{"why 99", "stderr", "ailang-worldd why: log entry 99 not found\n", "", 1},
+		{"why a2a:nobody:x --scan 1", "stderr", "ailang-worldd why: not found in the last 1 entries (invocationId)\n", "", 1},
+		{"why rest:w1:9", "stderr", "ailang-worldd why: receipt rest:w1:9 is indeterminate: nothing committed\n", "", 1},
+		{"why rest:nobody", "stderr", "ailang-worldd why: receipt rest:nobody is not-started: nothing committed\n", "", 1},
+		{"why " + zero, "stderr", "ailang-worldd why: " + zero + " is neither an object nor a world in this store\n", "", 1},
+	} {
+		v, raw := verbGet(a, c.q)
+		got := map[string]*string{"text": v.Text, "json": v.JSON, "stderr": v.Stderr}
+		if v.Outcome != "result" || v.Class != "read" || v.Verb != "why" || v.CLI != whyCLI(c.q) || v.Reason != "" || v.Exit != c.exit ||
+			got[c.member] == nil || !strings.HasPrefix(*got[c.member], c.prefix) || !strings.HasSuffix(*got[c.member], c.suffix) {
+			t.Errorf("%q = %s", c.q, raw)
+		}
+		if c.member == "stderr" && *got[c.member] != c.prefix {
+			t.Errorf("%q stderr = %q", c.q, *got[c.member])
+		}
+		for k, p := range got {
+			if k != c.member && p != nil {
+				t.Errorf("%q carries %s too: %s", c.q, k, raw)
+			}
+		}
+	}
+	for _, c := range []struct{ q, reason string }{
+		{"why", "ailang-worldd why: a target is required"},
+		{"why x", `ailang-worldd why: "x" is not an entry index, head, sha256:<ref>, a2a:/mcp:/rest:<id>, or -`},
+		{"why sha256:zz", `ailang-worldd why: "sha256:zz" is not a content address: hashref: digest is not lowercase hexadecimal: "zz"`},
+		{"why --scan 0 head", "ailang-worldd why: --scan must be 1..5000"},
+		{"why 1 2", `ailang-worldd why: unexpected argument "2"`},
+		{"why --bogus", "flag provided but not defined: -bogus"},
+		{"why --help", "see ailang-worldd why --help at a terminal"},
+		{"why -", "why - is not run from the app: the app has no stdin; at a terminal: ailang-worldd why -"},
+		{"why --result f", "--result is not run from the app: the app reads no file on the daemon's host; at a terminal: ailang-worldd why --result f"},
+		{"why --expect-store p head", "--expect-store is not run from the app: the app reads its own daemon's store; at a terminal: ailang-worldd why --expect-store p head"},
+	} {
+		if v, raw := verbGet(a, c.q); v.Outcome != "refused" || v.Reason != c.reason || v.CLI != whyCLI(c.q) || v.Text != nil || v.JSON != nil || v.Stderr != nil || v.Exit != 0 {
+			t.Errorf("%q = %s, want refused %q", c.q, raw, c.reason)
+		}
+	}
+	if v, _ := verbGet(a, "world why"); v.Outcome != "address" { // RK1: brane why still navigates
+		t.Errorf("world why = %+v", *v)
+	}
+}
+
+// VR2: an exhausted read deadline is (b), never a result, never "not found" — whether the walk's own
+// context expires before any read (the scan) or a handler answers 503 while the walk's context lives.
+func TestSurfaceVerbWhyDeadline(t *testing.T) {
+	a, s := whyRig(t)
+	s.d.errLog = io.Discard
+	s.d.readDeadline = expiredReadDeadline
+	for _, q := range []string{"why 1", "why a2a:nobody:x", "why head --json"} {
+		if v, raw := verbGet(a, q); v.Outcome != "refused" || v.Reason != "read deadline (-1ns) exceeded: the walk did not finish in the app; at a terminal it has 60 s: "+whyCLI(q) || strings.Contains(raw, "not found") {
+			t.Errorf("expired: %q = %s", q, raw)
+		}
+	}
+	base := s.d.reads
+	count := &whyFailReads{readStore: base}
+	s.d.reads = count
+	if v, _ := verbGet(a, "why sha256:"+strings.Repeat("0", 64)); v.Outcome != "refused" || count.n != 0 { // the walk's own deadline: no read starts
+		t.Errorf("expired: %d object reads, %+v", count.n, *v)
+	}
+	s.d.readDeadline = readDeadline
+	for _, reads := range []readStore{&whyFailReads{readStore: base, nth: 1, err: context.DeadlineExceeded}, // a 503 mid-walk
+		&whyFailReads{readStore: &whyFailReads{readStore: base, nth: 2, err: errors.New("x")}, nth: 3, err: context.DeadlineExceeded}} { // a 500, then a 503: (b) first
+		s.d.reads = reads
+		if v, raw := verbGet(a, "why 1"); v.Outcome != "refused" || !strings.HasPrefix(v.Reason, "read deadline (10s) exceeded: the walk did not finish in the app") {
+			t.Errorf("a 503 mid-walk = %s", raw)
+		}
+	}
+	// (a): the page went away; the answer says so, never a deadline or a result (nobody reads it).
+	s.d.reads = base
+	ctx, cancel := context.WithCancel(context.WithValue(context.Background(), http.LocalAddrContextKey, net.Addr(&net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 7644})))
+	cancel()
+	req := httptest.NewRequest("GET", "/surface/verb.json?q=why+1", nil).WithContext(ctx)
+	req.Host = apLoopback
+	rec := httptest.NewRecorder()
+	s.d.Handler().ServeHTTP(rec, req)
+	if !strings.Contains(rec.Body.String(), `"reason":"the request ended before the walk did"`) {
+		t.Errorf("a cancelled request = %d %s", rec.Code, rec.Body)
+	}
+}
+
+// whyFailReads fails the nth GetObject (counting from 1) with err; n counts every call.
+type whyFailReads struct {
+	readStore
+	nth, n int
+	err    error
+}
+
+func (f *whyFailReads) GetObject(ctx context.Context, ref hashref.HashRef) (store.Object, bool, error) {
+	if f.n++; f.n == f.nth {
+		return store.Object{}, false, f.err
+	}
+	return f.readStore.GetObject(ctx, ref)
+}
+
+// VR3: a read failure anywhere in the walk — every GetObject in turn, then a log read — is (c): no
+// chain, no ✗, no "not found", no store error text; the operator's log has it.
+func TestSurfaceVerbWhyReadFailure(t *testing.T) {
+	a, s := whyRig(t)
+	base := s.d.reads
+	count := &whyFailReads{readStore: base}
+	s.d.reads = count
+	verbGet(a, "why 1")
+	if count.n < 5 {
+		t.Fatalf("why 1 read %d objects", count.n)
+	}
+	const want = "the daemon could not read World during the walk, so no chain is shown; at a terminal: ailang-worldd why 1"
+	for nth := 1; nth <= count.n; nth++ {
+		var log bytes.Buffer
+		s.d.errLog = &log
+		s.d.reads = &whyFailReads{readStore: base, nth: nth, err: errors.New("SECRET-store-path")}
+		if v, raw := verbGet(a, "why 1"); v.Outcome != "refused" || v.Reason != want || strings.Contains(raw, "SECRET") || strings.Contains(raw, "✗") {
+			t.Errorf("GetObject #%d fails: %s", nth, raw)
+		}
+		if !strings.Contains(log.String(), "SECRET-store-path") {
+			t.Errorf("GetObject #%d: not logged: %q", nth, log.String())
+		}
+	}
+	s.d.errLog = io.Discard
+	s.d.reads = aguiFailReads{base, errors.New("SECRET-store-path")}
+	if v, raw := verbGet(a, "why 1"); v.Outcome != "refused" || v.Reason != want || strings.Contains(raw, "not found") {
+		t.Errorf("log read fails: %s", raw)
+	}
+}
+
+// recordingGetter is why's Getter over HTTP to url, keeping every path's (status, body).
+func recordingGetter(t *testing.T, url string, rec map[string][2]string) why.Getter {
+	return func(ctx context.Context, path string) (int, []byte, error) {
+		req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url+path, nil)
+		resp, err := http.DefaultClient.Do(req)
+		if err != nil {
+			return 0, nil, err
+		}
+		defer resp.Body.Close()
+		b, err := io.ReadAll(resp.Body)
+		rec[path] = [2]string{strconv.Itoa(resp.StatusCode), string(b)}
+		return resp.StatusCode, b, err
+	}
+}
+
+// VR4, VR6: every path the walk reads over HTTP (every target kind) answers byte for byte the same
+// in-process, /v1/health on the members an answer embeds; the walks' answers are equal; the private
+// mux's patterns are exactly the ones those paths reach, each bound to the method Handler registers.
+func TestWhyGetterTransportParity(t *testing.T) {
+	_, s := whyRig(t)
+	srv := httptest.NewServer(s.d.Handler())
+	defer srv.Close()
+	rec := map[string][2]string{}
+	ctx := boundedTestContext(t)
+	for _, c := range []struct {
+		target string
+		o      why.Options
+	}{{"1", why.Options{Scan: 500}}, {"head", why.Options{Scan: 500, JSON: true}}, {"rest:w1:3", why.Options{Scan: 500}},
+		{"rest:nobody", why.Options{Scan: 500}}, {s.head.Ref.String(), why.Options{Scan: 500}}, {"a2a:nobody:x", why.Options{Scan: 500}}, {"99", why.Options{Scan: 500}}} {
+		tg, err := why.ParseTarget(c.target)
+		if err != nil {
+			t.Fatal(err)
+		}
+		var tally whyTally
+		overHTTP := why.Run(why.Reader{Ctx: ctx, Get: recordingGetter(t, srv.URL, rec)}, tg, c.o)
+		inProc := why.Run(why.Reader{Ctx: ctx, Get: s.d.whyGetter(httptest.NewRequest("GET", "/surface/verb.json", nil), s.d.whyMux(), &tally)}, tg, c.o)
+		if overHTTP != inProc || tally != (whyTally{}) {
+			t.Errorf("%s: in-process %+v %+v, over HTTP %+v", c.target, inProc, tally, overHTTP)
+		}
+	}
+	mux, reached := s.d.whyMux(), map[string]bool{}
+	for path, want := range rec {
+		var tally whyTally
+		status, body, err := s.d.whyGetter(httptest.NewRequest("GET", "/surface/verb.json", nil), mux, &tally)(ctx, path)
+		got := [2]string{strconv.Itoa(status), string(body)}
+		if strings.HasPrefix(path, "/v1/health") {
+			got[1], want[1] = healthPin(t, got[1]), healthPin(t, want[1])
+		}
+		if err != nil || got != want {
+			t.Errorf("%s: in-process %v %v, over HTTP %v", path, got, err, want)
+		}
+		req, _ := http.NewRequest("GET", path, nil)
+		_, pattern := mux.Handler(req)
+		reached[pattern] = true
+	}
+	var patterns []string
+	for p := range s.d.whyRoutes() {
+		patterns = append(patterns, p)
+	}
+	if want := mapKeys(reached); !reflect.DeepEqual(whySorted(patterns), want) {
+		t.Errorf("whyRoutes %v, but the recorded walks reach %v", whySorted(patterns), want)
+	}
+}
+
+func healthPin(t *testing.T, body string) string {
+	var h struct {
+		Ref     string `json:"interpreter_ref"`
+		Version string `json:"interpreter_version"`
+	}
+	if err := json.Unmarshal([]byte(body), &h); err != nil {
+		t.Fatalf("health %q: %v", body, err)
+	}
+	return h.Ref + "|" + strings.SplitN(h.Version, "\n", 2)[0]
+}
+
+func mapKeys(m map[string]bool) []string {
+	var out []string
+	for k := range m {
+		out = append(out, k)
+	}
+	return whySorted(out)
+}
+
+func whySorted(s []string) []string { sort.Strings(s); return s }
+
+// VR7 (end to end): the cap runs in the route's deferred writer, for why (text and json) and for
+// episode report alike; the refusal carries the size and the CLI line and none of the members.
+func TestVerbAnswerCapEndToEnd(t *testing.T) {
+	a, _ := whyRig(t)
+	sizes := map[string]int{}
+	for _, q := range []string{"why 1", "why head --json", "episode report w1", "why 99"} {
+		v, _ := verbGet(a, q)
+		for _, m := range []*string{v.Text, v.JSON, v.Stderr} {
+			if m != nil {
+				sizes[q] += len(*m)
+			}
+		}
+	}
+	verbAnswerCap = sizes["why 99"]
+	t.Cleanup(func() { verbAnswerCap = MaxVerbAnswerBytes })
+	for _, q := range []string{"why 1", "why head --json", "episode report w1"} {
+		v, raw := verbGet(a, q)
+		if v.Outcome != "refused" || v.Reason != fmt.Sprintf("the answer is %d bytes, over the app's 1 MiB cap; at a terminal: ailang-worldd %s", sizes[q], q) ||
+			strings.Contains(raw, `"text"`) || strings.Contains(raw, `"json"`) || strings.Contains(raw, `"stderr"`) || strings.Contains(raw, `"exit"`) {
+			t.Errorf("%q over the cap = %s", q, raw)
+		}
+	}
+	if v, raw := verbGet(a, "why 99"); v.Outcome != "result" || v.Stderr == nil { // at the cap exactly
+		t.Errorf("at the cap = %s", raw)
+	}
+}
diff --git a/host/surfacepage/app_verb_test.go b/host/surfacepage/app_verb_test.go
--- a/host/surfacepage/app_verb_test.go
+++ b/host/surfacepage/app_verb_test.go
@@ -112,7 +112,7 @@ async function run(raw, fetchImpl) {
     assert.deepEqual(log, [['say', UNREAD]], JSON.stringify(body).slice(0, 80));
   }
   // Terminal: the CLI line, with a copy.
-  for (const q of ['why', 'episode finish w1 --preset p']) {
+  for (const q of ['episode finish w1 --preset p']) {
     const { act, log } = await run(q, golden(q)), b = B(q);
     assert.deepEqual(act, { panel: { head: '$ ' + b.cli, body: b.reason + '\nAt a terminal: ' + b.cli, copy: b.cli } }, q);
     assert.deepEqual(log, [['clear'], ['panel', act.panel]], q);
@@ -129,7 +129,7 @@ async function run(raw, fetchImpl) {
   }
   // VP11 (A3): a refusal with a CLI line opens the panel with the line to copy; one without (the world 4b shape, TooLong) is said.
   for (const [q, has] of [['log', 'at a terminal: ailang-worldd log'], ['ailang-worldd world get world-1', 'ailang-worldd world get world-1'],
-    ['episode', 'needs a subcommand']]) {
+    ['episode', 'needs a subcommand'], ['why', 'a target is required'], ['why --result f', 'reads no file']]) {
     const { act, log } = await run(q, golden(q)), b = B(q);
     assert.deepEqual(act, { panel: { head: '$ ' + b.cli, body: b.reason, copy: b.cli } }, q); assert.deepEqual(log, [['clear'], ['panel', act.panel]], q); assert(b.reason.includes(has), q);
   }
@@ -138,6 +138,16 @@ async function run(raw, fetchImpl) {
     const { act, log } = await run(q, golden(q));
     assert.deepEqual(act, { say: B(q).reason }, q); assert.deepEqual(log, [['say', B(q).reason]], q); assert(B(q).reason.includes(has), q); assert.equal(B(q).cli, undefined, q);
   }
+  // VP9 (A2): why's answers, generated by the route: text and json with exit 3, stderr with exit 1, each status on its own line;
+  // a hostile skill name in a chain stays text.
+  for (const [q, body] of [['why 2', B('why 2').text + 'exit 3'], ['why head --json', B('why head --json').json + 'exit 3'], ['why 9', B('why 9').stderr + 'exit 1']]) {
+    const { act } = await run(q, golden(q));
+    assert.deepEqual(act, { panel: { head: '$ ' + B(q).cli, body, copy: '' } }, q); assert.equal(B(q).exit, q === 'why 9' ? 1 : 3, q);
+  }
+  { const { act } = await run('why 2', reply(200, JSON.stringify({ ...B('why 2'), exit: undefined })));
+    assert.equal(act.panel.body, B('why 2').text); assert(B('why 2').text.includes('<img src=x onerror=alert(1)>')); }
+
+  assert(B(LONG).reason.startsWith('TooLong'));
   // VP2: no silent fallback. Every failure says its message and calls nothing else, for verbs and addresses alike.
   for (const q of ['why', 'log', 'row-189@6', 'world shared@0ffbb9c']) {
     const b = B(q), drop = k => { const c = { ...b }; delete c[k]; return JSON.stringify(c); };
@@ -148,7 +158,7 @@ async function run(raw, fetchImpl) {
     for (const [f, msg] of failures) { const { log } = await run(q, f); assert.deepEqual(log, [['say', msg]], q + ': ' + msg); }
   }
   // An answer missing its outcome's own members is unreadable too.
-  for (const body of [{ ...B('why'), cli: '' }, { ...B('why'), reason: undefined }, { ...B(MR), href: 7 }, { ...B('log'), reason: '' },
+  for (const body of [{ ...B('episode finish w1 --preset p'), cli: '' }, { ...B('why'), reason: undefined }, { ...B(MR), href: 7 }, { ...B('log'), reason: '' },
     { ...B('episode report w1'), json: '{}' }, { ...B('episode report nobody'), exit: 0 }, { ...B('episode report w1'), text: '' }]) {
     const { log } = await run('why', reply(200, JSON.stringify(body)));
     assert.deepEqual(log, [['say', UNREAD]], JSON.stringify(body).slice(0, 80));
diff --git a/host/surfacepage/testdata/verb-route.golden.json b/host/surfacepage/testdata/verb-route.golden.json
--- a/host/surfacepage/testdata/verb-route.golden.json
+++ b/host/surfacepage/testdata/verb-route.golden.json
@@ -10,7 +10,11 @@
   "episode report w1 --json": "{\"v\":1,\"schema\":\"world/surface-verb/v1\",\"q\":\"episode report w1 --json\",\"argv\":[\"episode\",\"report\",\"w1\",\"--json\"],\"verb\":\"episode report\",\"class\":\"read\",\"outcome\":\"result\",\"cli\":\"ailang-worldd episode report w1 --json\",\"json\":\"{\\\"episode\\\":\\\"w1\\\",\\\"entries\\\":[1,2],\\\"total_calls\\\":2,\\\"calls_by_tool\\\":{\\\"\\\\u003cimg src=x onerror=alert(1)\\\\u003e\\\":1,\\\"ailang-read\\\":1},\\\"effects\\\":[{\\\"effect\\\":\\\"Workspace.Read\\\",\\\"decision\\\":\\\"allowed\\\",\\\"count\\\":1}],\\\"budgets\\\":[{\\\"effect\\\":\\\"Workspace.Read\\\",\\\"scope\\\":\\\"worktree\\\",\\\"minted\\\":50,\\\"used\\\":1,\\\"left\\\":49}],\\\"plan_refusals\\\":{\\\"count\\\":1,\\\"by_tool\\\":{\\\"\\\\u003cimg src=x onerror=alert(1)\\\\u003e\\\":1},\\\"first\\\":[{\\\"entry\\\":2,\\\"tool\\\":\\\"\\\\u003cimg src=x onerror=alert(1)\\\\u003e\\\",\\\"reason\\\":\\\"\\\\u003cb\\\\u003ehostile\\\\u003c/b\\\\u003e \\\\\\\"reason\\\\\\\"\\\"}]},\\\"handler_refusals\\\":{\\\"count\\\":0,\\\"by_tool\\\":{},\\\"first\\\":[]},\\\"writes\\\":{\\\"written\\\":[],\\\"produced\\\":[],\\\"deltas_truncated\\\":[],\\\"delta_errors\\\":[],\\\"written_truncated\\\":0,\\\"produced_truncated\\\":0,\\\"unlisted\\\":0,\\\"deltas_truncated_omitted\\\":0,\\\"delta_errors_omitted\\\":0},\\\"reads\\\":{\\\"count\\\":0,\\\"files\\\":[],\\\"truncated\\\":0,\\\"unlisted\\\":1},\\\"first_at\\\":\\\"2027-01-15T08:00:01Z\\\",\\\"last_at\\\":\\\"2027-01-15T08:00:02Z\\\",\\\"calls_by_surface\\\":{\\\"mcp\\\":2},\\\"uncommitted_intents\\\":0,\\\"notes\\\":{\\\"budgets\\\":\\\"per (effect, scope) with at least one effect record: left = latest budgetAfter, used = cost of allowed records, minted = left + used\\\",\\\"entries\\\":\\\"first and last log entries of the episode's committed calls\\\",\\\"first_at\\\":\\\"coordinator logical time (unix seconds) of the first and last committed call, UTC\\\",\\\"handler_refusals\\\":\\\"calls whose effect was allowed and ran, but whose handler result is ok:false (the effect record alone shows these as allowed)\\\",\\\"plan_refusals\\\":\\\"calls whose plan named no effect and answered ok:false (refused before any effect ran)\\\",\\\"reads\\\":\\\"committed ailang-read calls with a v2 record whose every effect was allowed and did not fail and whose handler result is not ok:false; the path is the record input's \\\\\\\"path\\\\\\\" (a ranged read counts); a path that is not local and relative counts in unlisted and is never shown; files are ordered by first entry, at most 1000\\\",\\\"writes\\\":\\\"written: successful committed ailang-write/edit/delete calls, from input.path; produced: allowed non-failed effect result worktree_delta.changes; deltas_truncated and delta_errors name evidence gaps; local relative paths only, each path list ordered by first entry and capped at 1000; excess counts and unlisted are explicit\\\"}}\\n\"}\n",
   "log": "{\"v\":1,\"schema\":\"world/surface-verb/v1\",\"q\":\"log\",\"argv\":[\"log\"],\"verb\":\"log\",\"class\":\"other\",\"outcome\":\"refused\",\"reason\":\"log is a CLI command the app does not run; at a terminal: ailang-worldd log\",\"cli\":\"ailang-worldd log\"}\n",
   "row-189@6": "{\"v\":1,\"schema\":\"world/surface-verb/v1\",\"q\":\"row-189@6\",\"argv\":[\"row-189@6\"],\"verb\":\"\",\"class\":\"address\",\"outcome\":\"address\"}\n",
-  "why": "{\"v\":1,\"schema\":\"world/surface-verb/v1\",\"q\":\"why\",\"argv\":[\"why\"],\"verb\":\"why\",\"class\":\"read\",\"outcome\":\"terminal\",\"reason\":\"why reads World, but its walk runs in the CLI. The app cannot run it yet.\",\"cli\":\"ailang-worldd why\"}\n",
+  "why": "{\"v\":1,\"schema\":\"world/surface-verb/v1\",\"q\":\"why\",\"argv\":[\"why\"],\"verb\":\"why\",\"class\":\"read\",\"outcome\":\"refused\",\"reason\":\"ailang-worldd why: a target is required\",\"cli\":\"ailang-worldd why\"}\n",
+  "why --result f": "{\"v\":1,\"schema\":\"world/surface-verb/v1\",\"q\":\"why --result f\",\"argv\":[\"why\",\"--result\",\"f\"],\"verb\":\"why\",\"class\":\"read\",\"outcome\":\"refused\",\"reason\":\"--result is not run from the app: the app reads no file on the daemon's host; at a terminal: ailang-worldd why --result f\",\"cli\":\"ailang-worldd why --result f\"}\n",
+  "why 2": "{\"v\":1,\"schema\":\"world/surface-verb/v1\",\"q\":\"why 2\",\"argv\":[\"why\",\"2\"],\"verb\":\"why\",\"class\":\"read\",\"outcome\":\"result\",\"cli\":\"ailang-worldd why 2\",\"text\":\"why: entry 2 (matched by entry index)\\n  ✗ world   sha256:a475ce7680e560a4dbd0261379cf8a726f08ae28f71dbbb95dc7ec2a96a5957c  revision 2, stateRoot sha256:288610f1a5eaedc5e22ad00d8f7f782b9cda1e5e3c4d970c911c4bf0eadaccec, logHead sha256:e9c64feca5b5574877037083bce69af595d409ba7adcd2effd74d16b84bff21a (recomputed) — ✗ GET /v1/worlds/\\u003crecomputed ref\\u003e does not return this world\\n  ✗ entry   sha256:e9c64feca5b5574877037083bce69af595d409ba7adcd2effd74d16b84bff21a  #2 writtenBy coordinator:mcp prev sha256:20c8af45f11bffa9cfc9e7032bd569fdaf965de4927e48eb2a071b07e91e2f80 transitionFn sha256:294988fdc44414843c4d02e87709c1cfd9e59cfd2e380e2de0053448aa558b95 interpreter sha256:9666d9e8899447735cb9897b77dbb121754fd4d503609758755bd3fdae3a4b22 (not the daemon's pinned interpreter) — ✗ entry hash recomputes to sha256:38bfdd341a6ba2d15875096838074e75eb57d6eda9b12c88f4206ddba1db9967\\n  ✓ record  sha256:40c5088e8c9d3b5ab6139ea904d09dee16c93e1348a25651b54e136083a37f64  world/invocation-record/v2 invocation mcp:w1:2 episode w1 skill \\u003cimg src=x onerror=alert(1)\\u003e\\n  ✓ input   sha256:3450b301b391e55f7dc4e3d8e42403687b96fc8ccc07e48a3bf92d178cba4f68  {\\\"task\\\":\\\"mcp:w1:2\\\"}\\n  ✓ plan    sha256:058f9d7804781e6e970d48dec8ca7507fe50eb3e04de5b94523c116a4ea5a8cc  no effects; result {\\\"ok\\\":false,\\\"refused\\\":\\\"\\\\u003cb\\\\u003ehostile\\\\u003c/b\\\\u003e \\\\\\\"reason\\\\\\\"\\\"}\\n  ✗ output  sha256:288610f1a5eaedc5e22ad00d8f7f782b9cda1e5e3c4d970c911c4bf0eadaccec  18 B = world.stateRoot — ✗ the output is not the committed world's stateRoot\\n3 of 6 link(s) BROKEN — integrity refusal (exit 3)\\n\",\"exit\":3}\n",
+  "why 9": "{\"v\":1,\"schema\":\"world/surface-verb/v1\",\"q\":\"why 9\",\"argv\":[\"why\",\"9\"],\"verb\":\"why\",\"class\":\"read\",\"outcome\":\"result\",\"cli\":\"ailang-worldd why 9\",\"stderr\":\"ailang-worldd why: log entry 9 not found\\n\",\"exit\":1}\n",
+  "why head --json": "{\"v\":1,\"schema\":\"world/surface-verb/v1\",\"q\":\"why head --json\",\"argv\":[\"why\",\"head\",\"--json\"],\"verb\":\"why\",\"class\":\"read\",\"outcome\":\"result\",\"cli\":\"ailang-worldd why head --json\",\"json\":\"{\\n  \\\"entry\\\": 2,\\n  \\\"matchedBy\\\": \\\"head\\\",\\n  \\\"coordinator\\\": true,\\n  \\\"links\\\": [\\n    {\\n      \\\"name\\\": \\\"world\\\",\\n      \\\"ok\\\": false,\\n      \\\"ref\\\": \\\"sha256:a475ce7680e560a4dbd0261379cf8a726f08ae28f71dbbb95dc7ec2a96a5957c\\\",\\n      \\\"detail\\\": \\\"revision 2, stateRoot sha256:288610f1a5eaedc5e22ad00d8f7f782b9cda1e5e3c4d970c911c4bf0eadaccec, logHead sha256:e9c64feca5b5574877037083bce69af595d409ba7adcd2effd74d16b84bff21a (recomputed) — ✗ GET /v1/worlds/\\\\u003crecomputed ref\\\\u003e does not return this world\\\"\\n    },\\n    {\\n      \\\"name\\\": \\\"entry\\\",\\n      \\\"ok\\\": false,\\n      \\\"ref\\\": \\\"sha256:e9c64feca5b5574877037083bce69af595d409ba7adcd2effd74d16b84bff21a\\\",\\n      \\\"detail\\\": \\\"#2 writtenBy coordinator:mcp prev sha256:20c8af45f11bffa9cfc9e7032bd569fdaf965de4927e48eb2a071b07e91e2f80 transitionFn sha256:294988fdc44414843c4d02e87709c1cfd9e59cfd2e380e2de0053448aa558b95 interpreter sha256:9666d9e8899447735cb9897b77dbb121754fd4d503609758755bd3fdae3a4b22 (not the daemon's pinned interpreter) — ✗ entry hash recomputes to sha256:38bfdd341a6ba2d15875096838074e75eb57d6eda9b12c88f4206ddba1db9967\\\"\\n    },\\n    {\\n      \\\"name\\\": \\\"record\\\",\\n      \\\"ok\\\": true,\\n      \\\"ref\\\": \\\"sha256:40c5088e8c9d3b5ab6139ea904d09dee16c93e1348a25651b54e136083a37f64\\\",\\n      \\\"detail\\\": \\\"world/invocation-record/v2 invocation mcp:w1:2 episode w1 skill \\\\u003cimg src=x onerror=alert(1)\\\\u003e\\\"\\n    },\\n    {\\n      \\\"name\\\": \\\"input\\\",\\n      \\\"ok\\\": true,\\n      \\\"ref\\\": \\\"sha256:3450b301b391e55f7dc4e3d8e42403687b96fc8ccc07e48a3bf92d178cba4f68\\\",\\n      \\\"detail\\\": \\\"{\\\\\\\"task\\\\\\\":\\\\\\\"mcp:w1:2\\\\\\\"}\\\"\\n    },\\n    {\\n      \\\"name\\\": \\\"plan\\\",\\n      \\\"ok\\\": true,\\n      \\\"ref\\\": \\\"sha256:058f9d7804781e6e970d48dec8ca7507fe50eb3e04de5b94523c116a4ea5a8cc\\\",\\n      \\\"detail\\\": \\\"no effects; result {\\\\\\\"ok\\\\\\\":false,\\\\\\\"refused\\\\\\\":\\\\\\\"\\\\\\\\u003cb\\\\\\\\u003ehostile\\\\\\\\u003c/b\\\\\\\\u003e \\\\\\\\\\\\\\\"reason\\\\\\\\\\\\\\\"\\\\\\\"}\\\"\\n    },\\n    {\\n      \\\"name\\\": \\\"output\\\",\\n      \\\"ok\\\": false,\\n      \\\"ref\\\": \\\"sha256:288610f1a5eaedc5e22ad00d8f7f782b9cda1e5e3c4d970c911c4bf0eadaccec\\\",\\n      \\\"detail\\\": \\\"18 B = world.stateRoot — ✗ the output is not the committed world's stateRoot\\\"\\n    }\\n  ],\\n  \\\"ok\\\": false\\n}\\n\",\"exit\":3}\n",
   "world a b c": "{\"v\":1,\"schema\":\"world/surface-verb/v1\",\"q\":\"world a b c\",\"argv\":[\"world\",\"a\",\"b\",\"c\"],\"verb\":\"world\",\"class\":\"refused\",\"outcome\":\"refused\",\"reason\":\"unknown address: try \\\"world\\\", \\\"world NAME\\\", \\\"world NAME@N\\\", \\\"world NAME@HASH\\\" or \\\"back\\\"\"}\n",
   "world episode": "{\"v\":1,\"schema\":\"world/surface-verb/v1\",\"q\":\"world episode\",\"argv\":[\"world\",\"episode\"],\"verb\":\"\",\"class\":\"address\",\"outcome\":\"address\"}\n",
   "world session": "{\"v\":1,\"schema\":\"world/surface-verb/v1\",\"q\":\"world session\",\"argv\":[\"world\",\"session\"],\"verb\":\"\",\"class\":\"address\",\"outcome\":\"address\"}\n",
diff --git a/host/verbs/lookup.go b/host/verbs/lookup.go
--- a/host/verbs/lookup.go
+++ b/host/verbs/lookup.go
@@ -10,7 +10,7 @@ type Class string
 
 const (
 	Address  Class = "address"  // navigation: the page's own parseAddress path
-	Read     Class = "read"     // reads World: episode report runs; why answers terminal (V2)
+	Read     Class = "read"     // reads World: episode report and why run
 	Navigate Class = "navigate" // approve opens row 201's card
 	Mutate   Class = "mutate"   // changes World: terminal (V3/V4)
 	Host     Class = "host"     // needs a host process: terminal (V4)
@@ -97,7 +97,7 @@ func Lookup(argv []string) (Verb, []string) {
 	}
 	switch { // 4: the one-word rows
 	case first == "why":
-		return Verb{Words: "why", Class: Read, CLI: cli, Reason: "why reads World, but its walk runs in the CLI. The app cannot run it yet."}, argv[1:]
+		return Verb{Words: "why", Class: Read, CLI: cli}, argv[1:]
 	case first == "approve":
 		return Verb{Words: "approve", Class: Navigate, CLI: cli}, argv[1:]
 	case member(NeverWords, first):
diff --git a/host/verbs/lookup_test.go b/host/verbs/lookup_test.go
--- a/host/verbs/lookup_test.go
+++ b/host/verbs/lookup_test.go
@@ -24,8 +24,8 @@ var lookupCases = []struct {
 	{"world get", "", Address, "", false, nil},
 	{"World a b c", "", Address, "", false, nil},
 	{"ailang-worldd", "", Address, "", false, nil},
-	{"why", "why", Read, "why reads World", true, []string{}},
-	{"why e1 --json", "why", Read, "why reads World", true, []string{"e1", "--json"}},
+	{"why", "why", Read, "", true, []string{}},
+	{"why e1 --json", "why", Read, "", true, []string{"e1", "--json"}},
 	{"episode report w1", "episode report", Read, "", true, []string{"w1"}},
 	{"ailang-worldd episode report w1 --json", "episode report", Read, "", true, []string{"w1", "--json"}},
 	{"approve mr-0123456789abcdef", "approve", Navigate, "", true, []string{"mr-0123456789abcdef"}},
````

### Patch M6e — the CLI's bytes: VA3, VE1/VE3, XV2; QUICKSTART §19

````diff M6e
diff --git a/cmd/ailang-worldd/quickstart_verbs_test.go b/cmd/ailang-worldd/quickstart_verbs_test.go
--- a/cmd/ailang-worldd/quickstart_verbs_test.go
+++ b/cmd/ailang-worldd/quickstart_verbs_test.go
@@ -25,8 +25,8 @@ func TestQuickstartSection19MatchesTheBar(t *testing.T) {
 	section := guide[s19:]
 	_, block, _ := strings.Cut(section, "```text\n")
 	block, _, _ = strings.Cut(block, "```")
-	want := map[string]verbs.Class{"episode report": verbs.Read, "approve": verbs.Navigate, "episode finish": verbs.Mutate, "": verbs.Address}
-	seen := map[verbs.Class]bool{}
+	want := map[string]verbs.Class{"episode report": verbs.Read, "why": verbs.Read, "approve": verbs.Navigate, "episode finish": verbs.Mutate, "": verbs.Address}
+	seen := map[string]bool{}
 	for _, line := range strings.Split(strings.TrimSpace(block), "\n") {
 		argv, err := verbs.Split(line)
 		if err != nil {
@@ -36,7 +36,7 @@ func TestQuickstartSection19MatchesTheBar(t *testing.T) {
 		if c, ok := want[v.Words]; !ok || c != v.Class {
 			t.Errorf("§19 line %q is %q %s, which §19 does not describe", line, v.Words, v.Class)
 		}
-		seen[v.Class] = true
+		seen[v.Words] = true
 	}
 	if len(seen) != len(want) {
 		t.Errorf("§19's lines cover %v, want each of %v", seen, want)
diff --git a/cmd/ailang-worldd/verb_why_e2e_test.go b/cmd/ailang-worldd/verb_why_e2e_test.go
new file mode 100644
--- /dev/null
+++ b/cmd/ailang-worldd/verb_why_e2e_test.go
@@ -0,0 +1,139 @@
+package main
+
+// Row 243 V2b (design §6 VA3, VE1–VE3, XV2): the bar's why is the CLI's — its usage lines, its
+// stdout, stderr and exit — on the same daemon, and in a live Chrome.
+
+import (
+	"bufio"
+	"encoding/json"
+	"fmt"
+	"io"
+	"net/http"
+	"net/url"
+	"os"
+	"strings"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/coordinator"
+	"github.com/sunholo-data/ailang-world/host/verbs"
+)
+
+// verbRoute answers q from the rig's daemon's GET /surface/verb.json.
+func verbRoute(t *testing.T, base, q string) map[string]any {
+	t.Helper()
+	resp, err := http.Get(base + "/surface/verb.json?q=" + url.QueryEscape(q))
+	if err != nil {
+		t.Fatal(err)
+	}
+	defer resp.Body.Close()
+	body, _ := io.ReadAll(resp.Body)
+	var a map[string]any
+	if resp.StatusCode != 200 || json.Unmarshal(body, &a) != nil || a["cli"] != "ailang-worldd "+q {
+		t.Fatalf("%q: %d %s", q, resp.StatusCode, body)
+	}
+	return a
+}
+
+// VA3: every usage refusal of design §2.3 is the CLI's own first stderr line, generated by the CLI.
+func TestWhyBarAdmissionIsTheCLIs(t *testing.T) {
+	r := newWalkRig(t)
+	r.invoke("ep1", "ailang-read", "task-1", "hello\n")
+	for _, args := range [][]string{nil, {"x"}, {"sha256:zz"}, {"--scan", "0", "head"}, {"head", "--scan", "5001"}, {"1", "2"},
+		{"--json", "1", "2"}, {"-3"}, {"--bogus"}, {"--scan", "x", "head"}} {
+		code, stdout, stderr := whyRun(t, r.url, "", args...)
+		q := verbs.Quote(append([]string{"why"}, args...))
+		a := verbRoute(t, r.url, q)
+		if first, _, _ := strings.Cut(stderr, "\n"); code != exitUsage || stdout != "" || a["outcome"] != "refused" || a["reason"] != first {
+			t.Errorf("%q: the bar says %q; the CLI exits %d with %q", q, a["reason"], code, first)
+		}
+	}
+}
+
+// VE1, VE3: for every target form the walk rig resolves, the bar's answer is the CLI's — text or json
+// = stdout, stderr = stderr, exit = the exit — byte for byte on the same daemon; and the answers VC1
+// recorded from the CLI before the walk moved (testdata/why-characterise.golden.jsonl, never derived
+// from why.Run) are the bar's too, its JSON's one trailing newline included.
+func TestVerbRouteWhyIsTheCLIs(t *testing.T) {
+	r := newWalkRig(t)
+	first := r.invoke("ep1", "ailang-read", "task-1", "hello\n")
+	golden := map[string]struct{ Stdout, Stderr string }{}
+	f, err := os.Open(whyGolden)
+	if err != nil {
+		t.Fatal(err)
+	}
+	defer f.Close()
+	for sc := bufio.NewScanner(f); sc.Scan(); {
+		var g struct{ Name, Stdout, Stderr string }
+		if err := json.Unmarshal(sc.Bytes(), &g); err != nil {
+			t.Fatal(err)
+		}
+		golden[g.Name] = struct{ Stdout, Stderr string }{g.Stdout, g.Stderr}
+	}
+	firstLine := func(s string) string { l, _, _ := strings.Cut(s, "\n"); return l }
+	for q, want := range map[string]string{
+		"why --json 1": golden["flag before target"].Stdout, "why 1": golden["--expect-store this store"].Stdout,
+		"why 9": golden["index past head"].Stderr, "why rest:nobody": golden["rest: unknown receipt"].Stderr,
+		"why sha256:" + strings.Repeat("0", 64): golden["absent sha256"].Stderr,
+		"why":                                   firstLine(golden["no target"].Stderr), "why sha256:zz": firstLine(golden["malformed sha256"].Stderr),
+		"why 1 2": firstLine(golden["two targets"].Stderr),
+	} {
+		a := verbRoute(t, r.url, q)
+		got := a["reason"]
+		for _, m := range []string{"text", "json", "stderr"} {
+			if a["outcome"] == "result" && a[m] != nil {
+				got = a[m]
+			}
+		}
+		if want == "" || got != want {
+			t.Errorf("%q: the bar answers %q; the CLI printed before the move %q", q, got, want)
+		}
+	}
+	if a := verbRoute(t, r.url, "why --json 1"); !strings.HasSuffix(fmt.Sprint(a["json"]), "\n}\n") || strings.HasSuffix(fmt.Sprint(a["json"]), "}\n\n") {
+		t.Errorf("--json's trailing newline: %q", a["json"])
+	}
+	r.invoke("ep1", "ailang-read", "task-2", "second\n")
+	mcp := r.invokeAs(coordinator.SurfaceMCP, "ep2", "ailang-read", "x", "over mcp\n")
+	for _, args := range [][]string{{"1"}, {"head"}, {"head", "--json"}, {"--json", "1"}, {first.id}, {mcp.id}, {mcp.output}, {mcp.record},
+		{mcp.input}, {mcp.plan}, {mcp.effect}, {mcp.result}, {r.tip.Ref}, {"9"}, {"a2a:nobody:x", "--scan", "1"}, {"rest:nobody"}, {"--scan", "5000", "head"}} {
+		code, stdout, stderr := whyRun(t, r.url, "", args...)
+		q := verbs.Quote(append([]string{"why"}, args...))
+		a := verbRoute(t, r.url, q)
+		out := a["text"]
+		if contains(args, "--json") {
+			out = a["json"]
+		}
+		exit, _ := a["exit"].(float64)
+		if a["outcome"] != "result" || a["reason"] != nil || fmt.Sprint(out) != fmt.Sprint(map[bool]any{true: stdout, false: nil}[stdout != ""]) ||
+			fmt.Sprint(a["stderr"]) != fmt.Sprint(map[bool]any{true: stderr, false: nil}[stderr != ""]) || int(exit) != code {
+			t.Errorf("%q: the bar %v; the CLI exit %d stdout %q stderr %q", q, a, code, stdout, stderr)
+		}
+	}
+}
+
+// XV2: `why head` typed into #addr shows the CLI's exact stdout in the verb panel on the same daemon;
+// Escape closes the panel (V1 judge note 5).
+func TestVerbBarWhyInChrome(t *testing.T) {
+	b := startPasskeyBrowser(t)
+	b.do(t, nil, "tab", map[string]any{})
+	r := newWalkRig(t)
+	r.invoke("ep1", "ailang-read", "task-1", "hello\n")
+	code, want, stderr := whyRun(t, r.url, "", "head")
+	if code != exitOK || !strings.Contains(want, "all 7 link(s) verified") {
+		t.Fatalf("the CLI: exit %d stdout %q stderr %q", code, want, stderr)
+	}
+	b.do(t, nil, "openCard", map[string]any{"url": r.url + "/surface", "raw": true})
+	b.do(t, nil, "waitFor", map[string]any{"expression": `Boolean(window.World && document.getElementById('addr'))`})
+	b.do(t, nil, "waitFor", map[string]any{"expression": `(() => { const e = document.getElementById('addr'); e.focus(); e.value = 'why head'; return document.activeElement === e; })()`})
+	b.do(t, nil, "key", map[string]any{"key": "Enter"})
+	var got string
+	b.do(t, &got, "waitFor", map[string]any{"expression": `(() => { const p = document.querySelector('#verb pre'); return p && p.textContent.includes('link(s) verified') && p.textContent; })()`})
+	if got != want {
+		t.Fatalf("the bar shows\n%q\nthe CLI printed\n%q", got, want)
+	}
+	var gone bool
+	b.do(t, &gone, "waitFor", map[string]any{"expression": `(() => { document.dispatchEvent(new KeyboardEvent('keydown', { key: 'Escape', bubbles: true })); return !document.getElementById('verb'); })()`})
+	if !gone {
+		t.Fatal("Escape left the verb panel open")
+	}
+	b.do(t, nil, "openCard", map[string]any{"url": "about:blank", "raw": true}) // ends the page's live stream before the server closes
+}
diff --git a/docs/QUICKSTART.md b/docs/QUICKSTART.md
--- a/docs/QUICKSTART.md
+++ b/docs/QUICKSTART.md
@@ -1903,9 +1903,11 @@ The surface app's address bar (`/` focuses it) takes the CLI's own operator verb
 ```text
 episode report w1
 episode report w1 --json
+why head
+why 3 --json
 approve mr-0123456789abcdef
 episode finish w1 --preset se-tools
 world shared@0ffbb9c
 ```
 
-The first prints the report exactly as `ailang-worldd episode report w1` prints it at a terminal, and `--json` prints its JSON; an episode with no committed call shows the CLI's own line and `exit 1`. `approve <request-id|request-ref>` opens that request's card (§17, §18), where the decision is your key's. Every other verb shows the exact CLI line and says the app cannot run it yet, with a copy button: `episode start`, `episode finish` and `why` today. Global flags, `--follow`, `--markdown` and `--addr` are refused, as are the commands that stay at a terminal (`approver`, `session`, `service`, `serve`, `setup`). A brane named like a CLI verb (`log`, `episode`, …) is reached as `world NAME`, which always navigates. If the daemon does not answer within 15 seconds, or answers something the page cannot read, the bar says so and keeps your line; it never falls back to navigation.
+The first prints the report exactly as `ailang-worldd episode report w1` prints it at a terminal, and `--json` prints its JSON; an episode with no committed call shows the CLI's own line and `exit 1`. `approve <request-id|request-ref>` opens that request's card (§17, §18), where the decision is your key's. `why head` (or any target the CLI takes: an entry index, `sha256:<ref>`, `a2a:`/`mcp:`/`rest:<id>`, with `--scan N` and `--json`) prints the chain exactly as `ailang-worldd why head` prints it, with `exit 3` under it when a link is broken and the CLI's own line with `exit 1` when the target is not found. The app walks under the daemon's 10 s read deadline (a terminal has 60 s): a walk that needs longer, or a read that fails, is refused with the CLI line to copy, never shown as not found or as a broken link; `why -`, `--result` and `--expect-store` read stdin, a file or another store and run only at a terminal. A refusal that names a CLI line opens the panel with the line to copy. Every other verb shows the exact CLI line and says the app cannot run it yet, with a copy button: `episode start` and `episode finish` today. Global flags, `--follow`, `--markdown` and `--addr` are refused, as are the commands that stay at a terminal (`approver`, `session`, `service`, `serve`, `setup`). A brane named like a CLI verb (`log`, `episode`, …) is reached as `world NAME`, which always navigates. If the daemon does not answer within 15 seconds, or answers something the page cannot read, the bar says so and keeps your line; it never falls back to navigation.
````

