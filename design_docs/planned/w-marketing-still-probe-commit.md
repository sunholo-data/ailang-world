# w-marketing-still-probe-commit: make P0 item (1) a committed, asserting test on the capture store

Status: Planned (iteration 337), **r1**.
**Planner-Lane:** mission routing default
- **Item**: queue row 239 (`design_docs/world-mission.md:6311`), row 231's successor S-b
  (`w-marketing-position-checker.md` §7, and its carried finding (4) at line 278).
- **Kind**: release-hygiene. It moves no clause. No `.ail` changes.
- **Base**: `2f2935701959f4c727d4974db057ce7d9e6dd3d3` (= `origin/dev`), branch
  `sprint/world-iter337-still-probe-commit`. Toolchain `go1.26.6 darwin/arm64`, `AILANG v0.41.0`,
  tool `AILANG v0.54.15`, Chrome `156.0.8078.80` (V01).
- **Estimated**: 1 executor run. M1 is one new test file (~120 lines) plus a 5-line insertion. M2 is a
  one-line replacement in the orchestration doc.
- **Evidence**: `design_docs/verification/world-iter337/designer-*.txt`. The prototype ran in the
  worktree as `cmd/ailang-worldd/zz_iter337_still_proto_test.go` and was deleted. Its source is
  `designer-still-proto_test.go.txt`. Every `go test` ran under the process-group runner
  `world-iter336/run_bounded.sh.txt`.

## 1. Problem (measured at the base)

- No committed Go test asserts `repeat_equal` or `script_tags`. The match count over `*.go` is 0.
  Control: `requireProofPin` gives 13 matches (V02).
- §4.4 item (1) (`w-marketing-capture-orchestration.md:221`) tells the controller to `cp` the iter305
  probe `.txt` into `cmd/ailang-worldd/zz_iter305_probe_test.go`, run it and grep its log. The probe
  only `t.Logf`s `repeat_equal`/`script_tags` (probe lines 249–250). It also mints **its own** store,
  so the item-(1) record and the M7c4 drill (item (2), `TestMarketingChromeSixShots`) describe two
  different stores (iteration-315 judge NB4).
- `TestMarketingChromeSixShots` skips when `WORLD_CHROME` is unset (`marketing_chrome_test.go:52–54`).
  When it runs, it builds the store with `newMarketingChromeRig` (`:122`, defined `:159–246`). It
  builds `ailang-worldd` and `marketingcheck` (`:125–126`), then hands `rig.db` to
  `scripts/capture_marketing.sh` (`:136`). That script serves the db with `--workbench-still` and
  takes its routes from `marketingcheck routes` (`ShotRoutes`, `main.go:228`) (V03, V04). Item (3)'s
  rectangles enter as the `WORLD_MARKETING_RECTS` input (`:73`). The drill does not produce them.

## 2. Decisions (answers to the five questions)

### 2.1 Gating: AILANG pins, not Chrome, plus inside the Chrome drill (stronger than the row)

Item (1) is a check on the server's HTML and needs no Chrome. On the real `demo seed` S1–S6 store
the prototype passes in **8.81 s** with no Chrome (seed 5.73 s, `marketingcheck` build 0.32 s
warm or 7 s with a cold cache, still check 2.72 s) (V08, V11).

**Decision.** M1 commits two callers of one helper:
1. **`TestMarketingStillRoutes`**. It is gated only on what the real seed needs: `requireProofPin`
   (FATAL without `AILANG_BIN`) and `cliSeBins` (it SKIPs without `WORLD_TOOL_AILANG_BIN`, like
   every real-seed test, e.g. `demoM5cRealSeed`, `demo_test.go:1727–1732`) (V05). It runs in every
   executor, judge and controller `go test` that has both pins, Chrome or not.
2. **`TestMarketingChromeSixShots`** calls the same helper on its own rig before capture (§2.2).

Both callers run in every configuration where a Chrome-gated test would run, because the drill
FATALs without both pins (`marketing_chrome_test.go:69–71`). They also run where Chrome is
absent. That is strictly stronger than the row's "`WORLD_CHROME`-gated".

**Not ungated in GitHub CI, and why.** CI's go gate exports `AILANG_BIN` to `GITHUB_ENV`
(`ci.yml:176`) and runs `verify_go.sh` (`:281`). It installs the v0.54.15 tool binary only later
(`:384`) and never exports `WORLD_TOOL_AILANG_BIN` (V05). In that environment the test SKIPs (V09),
exactly as every existing real-seed test does. Exporting the tool pin to the go gate would switch
on every `cliSeBins` test on Linux at once, which is unmeasured and out of this row's scope
(Residual R1).

CI already guards the **mechanism** on every push, on a synthetic store. `TestWorkbenchStillAllRegions`
(`host/daemon/workbench_still_test.go`) needs only `AILANG_BIN` and PASSES in 2.04 s (V10). It kills
production mutants MS1–MS4. `TestWorkbenchStillDefaultLiveScript` and `TestServeWorkbenchStillFlag`
kill MS5 (V14). What CI lacks is the real seeded store, the provenance route
(`?object=<invocation>&payload=1`) and the capture store. The new test supplies all three.

### 2.2 Same store as item (2): call the helper inside the drill

Insert these lines in `TestMarketingChromeSixShots` immediately after `:126`, which builds `check`.
That is after `newMarketingChromeRig` (`:122`) and before the listener and `capture_marketing.sh`
(`:127`, `:136`):

```go
	// P0 item (1) on this store, before capture (row 239).
	checkMarketingStill(t, rig, marketingStillRoutes(t, check, rig.manifest))
	if t.Failed() {
		t.FailNow()
	}
```

- **Routes.** `marketingStillRoutes` runs the freshly built `check routes --manifest rig.manifest`.
  That is the same `ShotRoutes` derivation the capture script feeds the supervisor (`:129`), so item
  (1) checks exactly the pages captured. It requires 6 rows, de-duplicates them in shot order, and
  `checkMarketingStill` FATALs unless 4 distinct routes remain.
- **Daemon config.** The helper's in-process daemons mirror the capture script's `serve` flags
  (`capture_marketing.sh:116–118`). They are opened and closed before capture starts.
- **Measured.** With this patch the drill PASSES in 20.00 s on real Chrome. It logs 4 `STILL … repeat_equal=true
  script_tags=0 equals_default_minus_tag=true default_live_tags=1` lines, and the capture's own
  `inventory: equal` assertion still holds (V12). Items (1) and (2) are then one run on one store.
- **`t.FailNow`.** This only saves the capture's up-to-660 s when item (1) has already failed. It is
  not load-bearing, because the `t.Errorf`s already fail the test.

Item (3)'s store relation is unchanged (Residual R2).

### 2.3 The 2 s sleep stays, and it is load-bearing

The render path reads no clock: there are 0 `time.Now`/`time.Since` calls in `host/workbench/*.go` or
`host/daemon/workbench*.go` outside tests. Control: `daemon.go` has 7 (V06). So the sleep guards
time dependence that is hidden or arrives later, which is what P0 exists to catch. It costs 2 s.

**Measured.** MS4 makes still HTML differ more than 1 s after its first still render. It is killed
only by the repeat check. With the sleep removed (MT4), MS4 **survives** (PASS) (V13). The sleep sits
after all first GETs, as in `workbench_still_test.go:95–98`.

### 2.4 The helper (reference: `designer-still-proto_test.go.txt` lines 20–112)

The executor writes the new file `cmd/ailang-worldd/marketing_still_test.go` from the prototype, with
these changes only:
- `TestIter337StillRoutes` becomes `TestMarketingStillRoutes`.
- The `PROTO_ROUTES_NIL` hook is dropped, together with its now-unused `os` import.
- The `TIMING` log is optional.

Per route, `checkMarketingStill` takes three bodies: `d`, a default daemon's GET; `a`, the first still
GET; and `b`, the still GET 2 s later. It asserts each of the following with its own `t.Errorf`:

| Id | Assertion | Prototype line |
|---|---|---|
| A0 (control) | `d` contains the live.js tag exactly once | 101 |
| A1 | `a == b` (`repeat_equal`) | 104 |
| A2 | `b` has 0 `<script` (`script_tags`) | 107 |
| A3 | `a == d` with the tag removed once (`equals_default_minus_tag`) | 110 |
| A4 (null case) | 4 distinct routes, else `t.Fatalf` | 52 |

It logs one `STILL <route> … repeat_equal=… script_tags=… equals_default_minus_tag=…
default_live_tags=… html_sha256=…` line per route. That line is the P0 record.

A3 compares the **first** still GET and A1 compares first with second. As a result, each of the
production mutants MS2–MS5 below fires exactly one assertion, and MS1 fires A2 and A3 (V13).

### 2.5 Mutation proof (the planner runs these; all were measured in V13)

Production mutants edit `host/workbench/render.go` (`:203`, the `{{if .Still}}` line, or `:281`,
`Render`). Each runs `go test -count=1 -timeout 240s -v -run '^TestMarketingStillRoutes$'
./cmd/ailang-worldd` with both pins, and is reverted afterwards. **A build failure is not a kill.**

| Id | Mutation | Measured result (prototype) |
|---|---|---|
| MS1 still keeps live.js | `:203` becomes the bare `<script src="/workbench/live.js" defer></script>` | FAIL: A2 and A3 on all 4 routes |
| MS2 extra script in both modes | append `<script type="application/json">{}</script>` after `:203` | FAIL: **A2 only**, 4 routes |
| MS3 still-only constant | `{{if .Still}}<meta name="still" content="1">{{else}}…` | FAIL: **A3 only**, 4 routes. An HTML-comment marker is a **no-op** under `html/template` (it PASSED in round 1), so do not use one |
| MS4 time-dependent still | `Render`: when `p.Still`, append `" "` to `p.Title` once 1 s has passed since the first still render | FAIL: **A1 only**, 4 routes |
| MS5 strip made vacuous | `:203` becomes empty, so live.js is gone in both modes | FAIL: **A0 only**, 4 routes |
| MT1 null routes | `routes := marketingStillRoutes(…)[:0]` | FAIL: A4 `want 4 distinct capture routes, got 0` |
| MT2 inverted repeat | `if !repeat` → `if repeat` | FAIL: A1 on the healthy code |
| MT3 + MS3 | A3 disabled | **PASS**, so A3 is MS3's sole killer |
| MT4 + MS4 | the sleep removed | **PASS**, so the sleep is load-bearing (§2.3) |
| MT5 + MS5 | A0 disabled | **PASS**, so A0 is MS5's sole killer |
| MI1 (Chrome) | delete the §2.2 call from the drill | the M2 pass-rule count `STILL … = 4` gives **0**, so item (1) fails. Derived, not run |

## 3. Milestones

Every command runs under `run_bounded.sh`, the copy of `world-iter336/run_bounded.sh.txt`
(process-group kill, rc 124 on expiry). A 124 stops the milestone. Limits: each `go test` 600 s
(the Chrome run 1300 s), `go vet` and `go build ./...` 600 s each, full `go test ./...` 2400 s, and
`./scripts/verify_ail.sh` 1200 s.

**M1: the test** (`cmd/ailang-worldd/marketing_still_test.go` new; `marketing_chrome_test.go` +5 lines).
1. With `AILANG_BIN=$HOME/.pinned-ailang/ailang WORLD_TOOL_AILANG_BIN=$HOME/.pinned-ailang-tools/v0.54.15/ailang`:
   `go test -count=1 -timeout 540s -v -run '^TestMarketingStillRoutes$' ./cmd/ailang-worldd > "$LOG" 2>&1; echo rc=$?`
   gives `rc=0`. Then `grep -c '^--- PASS: TestMarketingStillRoutes'` = 1,
   `grep -c -- '--- SKIP'` = 0, and
   `grep -cE 'STILL .* repeat_equal=true script_tags=0 equals_default_minus_tag=true default_live_tags=1 '` = 4.
2. With `WORLD_TOOL_AILANG_BIN` unset and `AILANG_BIN` set, the same command gives rc 0 and
   `--- SKIP: TestMarketingStillRoutes`. This is the CI-shaped control.
3. The §2.5 MS1–MS5 and MT1–MT5 rows reproduce, each as listed.
4. `gofmt -l cmd/ailang-worldd` (empty), `go vet ./cmd/ailang-worldd`, `go build ./... && go test ./...` (pins set)
   and `./scripts/verify_ail.sh` all give rc 0.
5. `git diff --name-only 2f293570 -- . ':!design_docs'` lists exactly the two test files.

**M2: §4.4 item (1) text** (`design_docs/planned/w-marketing-capture-orchestration.md` only; lands
with M1 or not at all).
- **Pre-paste check (hard stop).** `grep -cF '  1. **Route-level still determinism.**' <doc>` = 1
  (V07; misspelled-anchor control = 0). Any other count stops M2. Never paste by line number.
- **Paste.** Replace the **one whole line** that contains that anchor with the single line inside this
  fence, byte-exact. It begins with two ASCII spaces:

```text
  1. **Route-level still determinism.** Serve the store under WorkbenchStill and repeat-GET every capture route 2 s apart, asserting byte-identical HTML, `script_tags=0`, and still HTML equal to the default HTML minus exactly the live.js tag. The result is a required P0 record. Row 239 committed this check as `checkMarketingStill` (`cmd/ailang-worldd/marketing_still_test.go`; it `t.Errorf`s on `repeat_equal=false`, on `script_tags>0` and on a still page that is not the default minus exactly the live.js tag); the M7c4 drill runs it on its own seeded store before capture, so items (1) and (2) are one run on one store. Command, from the repository root at the capture commit, with both AILANG pins exported: `WORLD_CHROME=<Chrome> WORLD_MARKETING_RECTS=<rects> go test -count=1 -timeout 1200s -v -run '^TestMarketingChromeSixShots$' ./cmd/ailang-worldd > <log> 2>&1`. It passes only if all of these hold: rc 0; one `--- PASS: TestMarketingChromeSixShots`; and `grep -cE 'STILL .* repeat_equal=true script_tags=0 equals_default_minus_tag=true default_live_tags=1 '` = 4 (row 239).
```

- **Scope.** The change is confined to item (1)'s command and pass rule. Items (2)–(4), the rest of
  §4.4, §7 (C12 already reads "byte-identical repeat-GET HTML on every capture route,
  `script_tags=0`") and §9/§10 are **not touched**. Whether the orchestration doc's own quorum must
  review this one line is the controller's call (row 238 V27). This design does not presume it.

Acceptance (run all of these; pipe nothing):
- `git diff --numstat -- design_docs/planned/w-marketing-capture-orchestration.md` gives
  `1	1	design_docs/planned/w-marketing-capture-orchestration.md`.
- `grep -cF` gives `1` for each of `checkMarketingStill`, `zz_iter305_probe_test.go` and
  `TestIter305SeededStillCaptureProbe`. The last two are V40's history at `:558`; each is 2 at the base (V07).
- `go test -count=1 -timeout 300s -run '^TestMarketing(Position|Phone)Parity$' ./scripts/marketingcheck`
  gives rc 0. Its `pin-registry` subtests read this doc.
- The pasted command, with the drill's rig env (V12's), gives rc 0 and the count 4.
- Measured on a scratch-copy dry run (V16): numstat `1 1`, each `grep -cF` 1, parity rc 0 with both
  `pin-registry` subtests PASS. The pass rule applied to the V12 log gives 4 and 1 PASS line, and on
  the MS1 mutant log it gives 0 (V17).

## 4. Conflict surface, concurrent work, residuals

- **Files.** `cmd/ailang-worldd/marketing_still_test.go` (new), `cmd/ailang-worldd/marketing_chrome_test.go`
  (+5), and one line of the orchestration doc. Nothing changes in `.ail`, `go.mod`, `.github/`,
  `tools/launchd/*`, `host/` or `scripts/`. `render.go` is touched only by reverted mutants.
- **Concurrent work (V15).**
  - There are 0 open PRs. Control: the merged list shows #462 and #461.
  - Of the 311 origin branches, 7 differ on these paths. All 7 differ only in `host/workbench/render.go`,
    and all 7 are **merged** PRs from 09-24 to 10-13 (#143–#249). This row does not change that file.
  - Control: the same diff form on `95179277` lists the orchestration doc.
  - Neighbouring rows name other files: 240 names `scripts/marketingcheck` (`phone.go:157`, `main.go`),
    246 names `scripts/marketingcheck/position.go:142` and 247 names `scripts/marketingcheck/capture_test.go`.
  - On any new change to these paths, re-base and re-run V07 before M2.
- **R1, unmeasured.** GitHub CI does not run the seeded check (V05, V09). A follow-up row could add one
  CI step after `ci.yml:384`:
  `WORLD_TOOL_AILANG_BIN=$HOME/.pinned-ailang-tools/v0.54.15/ailang go test -run '^TestMarketingStillRoutes$' ./cmd/ailang-worldd`.
  Whether `demo seed` runs on the Linux runner is UNMEASURED. It cannot be measured without pushing.
- **R2.** Item (3), the V60 rectangle probe, stays a separate instrument. The drill consumes
  `WORLD_MARKETING_RECTS` (`:73`) and builds its store under `t.TempDir()` (`:162`), so the drill
  cannot hand its store to item (3). This row does not change that.
- **R3, flake.** The prototype ran healthy 2× standalone (8.81 s, 7.96 s) and 1× in the drill
  (V08, V12, V13). That is not a flake-rate measurement.

## 5. Verification Log (base `2f293570`; files under `design_docs/verification/world-iter337/`)

| V | Command | Observed |
|---|---|---|
| V01 | `git rev-parse HEAD origin/dev; go version; ailang --version` ×2; `Chrome --version` (`designer-env.txt`) | `2f293570…` twice; `go1.26.6 darwin/arm64`; `AILANG v0.41.0`; `AILANG v0.54.15`; `Google Chrome 156.0.8078.80` |
| V02 | `git grep -n -E 'repeat_equal\|script_tags' 2f293570 -- '*.go' \| wc -l`; control `git grep -n requireProofPin … \| wc -l` | `0`; control `13` |
| V03 | `grep -n -E 'WORLD_CHROME\|newMarketingChromeRig\|"go", "build"\|capture_marketing.sh\|…' marketing_chrome_test.go` | `52/54` WORLD_CHROME skip; `73` rects env; `122` rig; `125–126` builds; `136` capture; `155` rig type; `159` def; `245` return |
| V04 | `grep -n -E 'workbench-still\|step routes\|step check' scripts/capture_marketing.sh`; `grep -n 'func ShotRoutes\|func runRoutes' main.go` | `116 … serve --db "$db" … --workbench-still`; `129 step routes … routes --manifest`; `141 step check … --rectangles`; `228 func ShotRoutes`; `240 func runRoutes` |
| V05 | `grep -n -E 'GITHUB_ENV\|verify_go.sh\|Install PINNED v0.54.15\|WORLD_TOOL_AILANG_BIN' ci.yml`; `grep -A8 'func cliSeBins'`; `grep -A5 'func demoM5cRealSeed'` | `176 AILANG_BIN=… >> GITHUB_ENV`; `281 run: ./scripts/verify_go.sh`; `384 Install PINNED v0.54.15`; **no** `WORLD_TOOL_AILANG_BIN` line. `cliSeBins` `t.Skipf` if either env is empty; `demoM5cRealSeed` calls `requireProofPin` + `cliSeBins` |
| V06 | `git grep -n -E 'time\.(Now\|Since)' -- 'host/workbench/*.go' 'host/daemon/workbench*.go' \| grep -v _test.go \| wc -l`; control `git grep -c … host/daemon/daemon.go`; `sed -n 203p;281p render.go` | `0`; control `7`; `{{if .Still}}{{else}}<script src="/workbench/live.js" defer></script>{{end}}`; `func Render(…) { return pageTemplate.Execute(w, p) }` |
| V07 | `grep -cF` anchors in the orchestration doc | `  1. **Route-level still determinism.**` 1 (line 221); `  2. **Six-shot two-pass.**` 1; `zz_iter305_probe_test.go` 2 (221, 558); `TestIter305SeededStillCaptureProbe` 2; `checkMarketingStill` 0; misspelled control 0 |
| V08 | prototype `TestIter337StillRoutes`, both pins (`designer-proto-runs.txt`) | rc 0; 4× `STILL … repeat_equal=true script_tags=0 equals_default_minus_tag=true default_live_tags=1`; `TIMING seed=5.73s build=0.32s still=2.72s total=8.77s`; `--- PASS (8.81s)` |
| V09 | same plus `TestMarketingChromeSixShots`, `AILANG_BIN` only (CI-shaped) | rc 0; `--- SKIP: TestMarketingChromeSixShots`; `--- SKIP: TestIter337StillRoutes` with `WORLD_TOOL_AILANG_BIN set=false` |
| V10 | `go test -v -run '^TestWorkbenchStillAllRegions$' ./host/daemon`, `AILANG_BIN` only | rc 0; `repeat GETs byte-identical on all three routes after 2.006s`; `--- PASS (2.04s)` |
| V11 | `GOCACHE=$(mktemp -d) go build ./scripts/marketingcheck`; `go vet ./cmd/ailang-worldd` with the prototype present | cold build rc 0, 7 s; vet rc 0 |
| V12 | drill with the §2.2 patch, real Chrome, `WORLD_MARKETING_RECTS=world-iter312/r1-rects.json` (`designer-sixshots-with-still.txt`) | rc 0, wall 22 s; 4 `STILL` pass lines (`marketing_chrome_test.go:128`); `--- PASS: TestMarketingChromeSixShots (20.00s)`; patch reverted (`git status` clean apart from the prototype) |
| V13 | 14 `run.sh` runs over `mutate.py` (`designer-mutants.txt`) | healthy PASS; MS1 A2+A3; MS2 A2; MS3(meta) A3; MS4 A1; MS5 A0; MT1 A4; MT2 A1; MT3+MS3, MT4+MS4, MT5+MS5 PASS. Round-1 MS3 (`<!--still-->`) PASSED, and round-1 MT1 was a build failure, so both were redone |
| V14 | MS1–MS5 against `TestWorkbenchStillAllRegions`; MS5 against `TestWorkbenchStillDefaultLiveScript` and `TestServeWorkbenchStillFlag` | MS1/MS2 `:54 zero script tags`; MS3 `:83 changed HTML beyond the script`; MS4 `:106 repeat GET … varies`; MS5 **PASS** there, but FAIL in `render_test.go:675` and `main_test.go:217` |
| V15 | `gh pr list --state open`; a 311-branch `git diff --name-only 2f293570...$b -- <paths>` under `bash -c`; control `95179277^...95179277`; the `gh pr list --head` state of each hit (`designer-sweep.txt`) | 0 open, control #462/#461; row 240 names `phone.go:157`, `main.go`, `scripts/marketingcheck`; 7 hits, all `render.go` only, all MERGED (#143, #145, #157, #158, #160, #234, #249); control lists the orchestration doc. A first zsh run did not split `$P`, and its empty control exposed that |
| V16 | M2 dry run in a `git archive` scratch copy: paste the fenced line over the anchor line; `git diff --no-index --numstat`; `grep -cF` ×4; parity tests (`designer-m2-dryrun.txt`) | pasted at line 221; `1 1`; `checkMarketingStill`, `zz_iter305_probe_test.go`, `TestIter305SeededStillCaptureProbe` and the item-(2) anchor each 1; parity rc 0, both `pin-registry` PASS |
| V17 | the M2 pass-rule `grep -cE 'STILL .* repeat_equal=true … default_live_tags=1 '` and the `--- PASS` count on the V12 log; the same grep on the MS1 log as a control | `4`, `1`; control `0` |
