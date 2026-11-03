# w-marketing-still-probe-commit: make P0 item (1) a committed, asserting test on the capture store

Status: Planned (iteration 337), **r3**: a controller narrow-refinement carve-out after design quorum r2 (§7). r2 was the one designer revision after design quorum r1 was BLOCKED
(`mission-quorum-world/w-marketing-still-probe-commit-2026-11-03T00-14-32Z.json`: gpt6-1-sol and
oc-glm-5-3 REJECT, gemini-4-argon PASS, no absentees, no direction dispute). The fixes are in §6.
**Planner-Lane:** mission routing default
- **Item**: queue row 239 (`design_docs/world-mission.md:6311`), row 231's successor S-b
  (`w-marketing-position-checker.md` §7, and its carried finding (4) at line 278).
- **Kind**: release-hygiene. It moves no clause. No `.ail` changes.
- **Base**: `2f293570` (= `origin/dev`), branch `sprint/world-iter337-still-probe-commit`; go1.26.6,
  AILANG v0.41.0, tool v0.54.15, Chrome 156.0.8078.80 (V01).
- **Estimated**: 1 executor run. M1 is a new test file (~120 lines) plus a 5-line insertion; M2 is one line.
- **Evidence**: `design_docs/verification/world-iter337/designer-*.txt`. The prototypes ran as
  `cmd/ailang-worldd/zz_iter337_still_proto_test.go` and were deleted (sources banked as `.txt`).
  Every `go test` ran under `world-iter336/run_bounded.sh.txt`.

## 1. Problem (measured at the base)

- No committed Go test asserts `repeat_equal` or `script_tags`. The match count over `*.go` is 0.
  Control: `requireProofPin` gives 13 matches (V02).
- §4.4 item (1) (`w-marketing-capture-orchestration.md:221`) tells the controller to `cp` the iter305
  probe `.txt` into `cmd/ailang-worldd/zz_iter305_probe_test.go`, run it and grep its log.
  - The probe's only `t.Errorf` is `:241`, which fires on a still status other than 200.
  - `script_tags` and `repeat_equal` appear only in the `t.Logf`s at `:219` (DEFAULT) and `:249` (STILL).
  - It mints **its own** store: `t.TempDir()` at `:65`, then `demo seed` at `:139` (V20).

  So the item-(1) record and the M7c4 drill (item (2)) describe two different stores.
- `TestMarketingChromeSixShots` skips without `WORLD_CHROME` (`marketing_chrome_test.go:52–54`). Otherwise:
  - it builds the store with `newMarketingChromeRig` (`:122`, defined `:159–246`);
  - it builds `ailang-worldd` and `marketingcheck` (`:125–126`) and hands `rig.db` to `capture_marketing.sh` (`:136`);
  - that script serves the db `--workbench-still`, with routes from `marketingcheck routes` (`main.go:228`) (V03, V04);
  - item (3)'s rectangles are an input, `WORLD_MARKETING_RECTS` (`:73`).

## 2. Decisions (answers to the five questions)

### 2.1 Gating: AILANG pins, not Chrome, plus inside the Chrome drill (stronger than the row)

Item (1) checks server HTML and needs no Chrome: on the real `demo seed` S1–S6 store the prototype
passes in **8.81 s** (seed 5.73 s, `marketingcheck` build 0.32 s warm / 7 s cold, check 2.72 s) (V08, V11).

**Decision.** M1 commits two callers of one helper:
1. **`TestMarketingStillRoutes`**. It is gated only on what the real seed needs: `requireProofPin`
   (FATAL without `AILANG_BIN`) and `cliSeBins` (it SKIPs without `WORLD_TOOL_AILANG_BIN`, like
   every real-seed test, e.g. `demoM5cRealSeed`, `demo_test.go:1727–1732`) (V05). It runs in every
   executor, judge and controller `go test` that has both pins, Chrome or not.
2. **`TestMarketingChromeSixShots`** calls the same helper on its own rig before capture (§2.2).

**Why that is stronger than the row's "`WORLD_CHROME`-gated".** When `WORLD_CHROME` is set, the drill
FATALs without both pins (`:71`). It does so only after its skip (`:52–54`) and Chrome checks
(`:56–67`), and was exercised to rc 1 at `:71` (V18). So any passing `WORLD_CHROME` run has both pins,
and `TestMarketingStillRoutes` runs there too. It also runs where Chrome is absent.

**Not ungated in GitHub CI, and why.**
- CI exports `AILANG_BIN` (`ci.yml:176`) before the go gate's `verify_go.sh` (`:281`). It installs the
  tool binary later (`:384`) and never exports `WORLD_TOOL_AILANG_BIN` (V05).
- So the test SKIPs there (V09), like every real-seed test. Exporting the pin would switch on every
  `cliSeBins` test on Linux at once, which is unmeasured (R1).
- CI already guards the **mechanism** on a synthetic store. `TestWorkbenchStillAllRegions` (2.04 s,
  V10) kills MS1–MS4, and two other CI tests kill MS5 (V14).
- What CI lacks is the real seed, the provenance route (`?object=<invocation>&payload=1`) and the
  capture store. The new test adds all three.

### 2.2 Same store as item (2): call the helper inside the drill

Insert after `:126` (which builds `check`): after the rig (`:122`), before the listener and capture (`:127`, `:136`):

```go
	// P0 item (1) on this store, before capture (row 239).
	checkMarketingStill(t, rig, marketingStillRoutes(t, check, rig.manifest))
	if t.Failed() {
		t.FailNow()
	}
```

- **Routes.** `marketingStillRoutes` runs the built `check routes --manifest rig.manifest`, the same
  `ShotRoutes` the capture uses (`:129`). It requires 6 rows and de-duplicates them; the helper needs 4.
- **Daemon config.** The helper's in-process daemons take the capture daemon's flags
  (`capture_marketing.sh:116–118`, the whole `serve` command, V19) field by field:

  | `serve` flag (script) | `daemon.Config` field (helper, prototype `:54–55`) |
  |---|---|
  | `--db "$db"` | `DBPath: rig.db` |
  | `--bind "$bind"` (must match `127.0.0.1:<port>`, `:101`) | none. `httptest.NewServer` listens on `127.0.0.1:<random port>`, and `BindHost: daemon.DefaultBindHost` is set |
  | `--workbench-still` | `WorkbenchStill: true` (still daemon only) |
  | `--ailang-bin`, `--tool-ailang-bin` | `AilangBin: rig.interp`, `ToolAilangBin: rig.tool` |
  | `--workspace-root` | `WorkspaceRoot: rig.ws` |
  | `--proof-key-file` | `ProofKeyFile: rig.key` |
  | `--proof-required-identity sameRef`, `… servesEntry` | `ProofRequiredIdentities: {"sameRef", "servesEntry"}` |

  **The bind cannot change the HTML, measured (V21).**
  - Every A3 comparison already crosses two ports (`:59766` default, `:59784` still).
  - The real `ailang-worldd serve` was started with the script's exact flags on a third port, on the
    same db, after the helper returned. Its still HTML is **byte-equal** to the helper's on 4/4 routes,
    and no body contains the bind or the port.
  - That also covers the `serve` defaults the helper leaves zero (`ExamplesDir`, `CredentialsPath`).
  - Control: `serve` without `--workbench-still` gives 0/4 equal, each 48 B longer, which is the tag's length.
- **Measured.** The patched drill PASSES in 20.00 s on real Chrome, with 4 passing `STILL` lines and
  its `inventory: equal` check intact (V12). Items (1) and (2) are then one run on one store.
- **`t.FailNow`** only skips a capture that is already doomed. It is not load-bearing.

Item (3)'s store relation is unchanged (Residual R2).

### 2.3 The 2 s sleep stays: this test catches delayed variation only with it

V06 found no direct time.Now/time.Since calls in non-test files matching host/workbench/*.go or
host/daemon/workbench*.go. This is not a transitive audit of clock dependence. We retain the 2 s
delay because the measured MS4/MT4 comparison shows that this test detects the injected delayed
variation only with the delay present; it does not prove independence from all possible
time-dependent behavior. The delay costs 2 s.

**Measured.** MS4 (still HTML changes 1 s after the first still render) is killed only by A1; with
the sleep removed (MT4) it **survives** (V13). The sleep follows all first GETs (`workbench_still_test.go:95–98` idiom).

### 2.4 The helper (reference: `designer-still-proto_test.go.txt` lines 20–112)

The executor writes the new file `cmd/ailang-worldd/marketing_still_test.go` from the prototype, with
these changes only: `TestIter337StillRoutes` becomes `TestMarketingStillRoutes`; the
`PROTO_ROUTES_NIL` hook and its `os` import are dropped; the `TIMING` log is optional.

Key signatures, as they will appear in the new file:

```go
const marketingStillLiveTag = `<script src="/workbench/live.js" defer></script>` // render.go:203's tag
func marketingStillRoutes(t *testing.T, check, manifest string) []string        // `check routes`, 6 rows, de-duplicated
func checkMarketingStill(t *testing.T, rig marketingChromeRig, routes []string) // P0 item (1); t.Errorf per assertion
func TestMarketingStillRoutes(t *testing.T)                                     // requireProofPin + cliSeBins + newMarketingChromeRig
```

**Required: no open handles on `rig.db` once the helper returns** (gemini-4-argon).
- `checkMarketingStill` opens each in-process daemon and its `httptest` server in a per-mode scope,
  and closes both with a `defer` in that scope. In the prototype that is the `serve` closure (`:56`),
  whose `:64` is `defer func() { srv.Close(); _ = d.Close() }()`.
- Both scopes (`:79`, `:84`) end before the assertions run and before the helper returns. Every
  response body is closed (`:70`).
- **Measured.** A separate `serve` process reopened `rig.db` after the helper returned (V21), and so
  did the drill's capture (V12).

Per route, `checkMarketingStill` takes three bodies: `d`, a default daemon's GET; `a`, the first still
GET; and `b`, the still GET 2 s later. It asserts each of the following with its own `t.Errorf`:

| Id | Assertion | Prototype line |
|---|---|---|
| A0 (control) | `d` contains the live.js tag exactly once | 101 |
| A1 | `a == b` (`repeat_equal`) | 104 |
| A2 | `b` has 0 `<script` (`script_tags`) | 107 |
| A3 | `a == d` with the tag removed once (`equals_default_minus_tag`) | 110 |
| A4 (null case) | 4 distinct routes, else `t.Fatalf` | 52 |

Per route it logs one `STILL <route> … repeat_equal=… script_tags=… equals_default_minus_tag=…
default_live_tags=… html_sha256=…` line, the P0 record. A3 uses the **first** still GET, so MS2–MS5
each fire exactly one assertion and MS1 fires A2 and A3 (V13).

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

### 2.6 Overlap with the existing still tests: a second caller, not a shared checker (r3)

Added by the controller under the narrow-refinement carve-out, answering gpt6-1-sol's quorum-r2
objection with its own proposed fix (§7). Measured in `controller-r3-overlap.txt` (V22).

| Existing test | Package | What it shares with A0–A4 | Why `cmd/ailang-worldd` cannot call it |
|---|---|---|---|
| `TestWorkbenchStillAllRegions` (`host/daemon/workbench_still_test.go`) | `daemon` (internal test) | A1 (repeat equality after a 2 s sleep, `:95–107`), A2 (zero `<script`, `:52`), A3 (default minus the tag, `:82`) | The file declares **1** function, the test itself (control: `demo_proof_test.go` declares 33), so its assertions are inline, not a helper. It reaches the daemon through unexported package-`daemon` helpers (`newHandlerDaemon`, `requestRecorder` in `handlers_test.go`; `workbenchLiveBody`, `workbenchNamedSection`, `packetSection` in three other `_test.go` files), flips `d.cfg.WorkbenchStill` on one daemon and reads `d.store`. Go does not export `_test.go` code across packages. |
| `TestWorkbenchStillDefaultLiveScript` (`host/workbench/render_test.go:655`) | `workbench` | A0, A2 on `Render` alone; no store, no HTTP | A test body in another package. |
| `TestServeWorkbenchStillFlag` (`cmd/ailang-worldd/main_test.go:193`) | `main` | A0, A2 on `Render` with the parsed flag | Same package, but it renders an empty `Page`, not a served store; there is no helper to share. |

**Decision: keep `checkMarketingStill` as a separate implementation.** Sharing would mean a new
non-test helper package (the repo has the pattern: `host/proctest`, `host/ptytest`) holding the three
byte comparisons (A1–A3, about 10 lines), plus rewriting `TestWorkbenchStillAllRegions` to call it.
That widens this row from two `cmd/ailang-worldd` test files to `host/`, edits a CI-guarding test,
and changes that test's failure semantics: it stops at the first failure with a byte-offset diff
(`t.Fatalf`, `:103`), while this row requires one `t.Errorf` per route and assertion so the P0 record
names every failing route. Ten shared lines do not pay for that.

**What must stay synchronized**, and what catches drift:
- **The live.js tag.** The literal appears at `render.go:203` and at 5 test sites (V22); the new file
  adds a sixth as `marketingStillLiveTag`. If `render.go`'s tag changes, the new test's A0 control
  fails (it requires the tag exactly once in the default page), and so does `workbench_still_test.go:82`.
  MS5 (§2.5) shows A0 is the sole killer of a vacuous strip in the new test.
- **The zero-`<script` predicate** (A2 here, `:52` there) and **the 2 s spacing** (§2.3 here,
  `:95–98` there).
- M1 adds a header comment to `marketing_still_test.go` naming these three invariants and the sibling
  sites `host/workbench/render.go:203` and `host/daemon/workbench_still_test.go:52,82,95`.

## 3. Milestones

Every command runs under a copy of `world-iter336/run_bounded.sh.txt` (process-group kill, rc 124
stops the milestone). Limits: each `go test` 600 s (Chrome run 1300 s); vet and build 600 s each;
`go test ./...` 2400 s; `verify_ail.sh` 1200 s.

**M1: the test** (`cmd/ailang-worldd/marketing_still_test.go` new; `marketing_chrome_test.go` +5 lines).
1. With `AILANG_BIN=$HOME/.pinned-ailang/ailang WORLD_TOOL_AILANG_BIN=$HOME/.pinned-ailang-tools/v0.54.15/ailang`:
   `go test -count=1 -timeout 540s -v -run '^TestMarketingStillRoutes$' ./cmd/ailang-worldd > "$LOG" 2>&1; echo rc=$?`
   gives `rc=0`; `grep -c '^--- PASS: TestMarketingStillRoutes'` = 1; `grep -c -- '--- SKIP'` = 0; `grep -cE 'STILL .* repeat_equal=true script_tags=0 equals_default_minus_tag=true default_live_tags=1 '` = 4.
2. With `WORLD_TOOL_AILANG_BIN` unset and `AILANG_BIN` set, the same command gives rc 0 and
   `--- SKIP: TestMarketingStillRoutes`. This is the CI-shaped control.
3. The §2.5 MS1–MS5 and MT1–MT5 rows reproduce, each as listed.
4. `gofmt -l cmd/ailang-worldd` (empty), `go vet ./cmd/ailang-worldd`, `go build ./... && go test ./...` (pins set)
   and `./scripts/verify_ail.sh` all give rc 0.
5. `git diff --name-only 2f293570 -- . ':!design_docs'` lists exactly the two test files.

**M2: §4.4 item (1) text** (the orchestration doc only; lands with M1 or not at all).
- **Pre-paste check (hard stop).** `grep -cF '  1. **Route-level still determinism.**' <doc>` = 1
  (V07; misspelled-anchor control = 0). Any other count stops M2. Never paste by line number.
- **Paste.** Replace the **one whole line** that contains that anchor with the single line inside this
  fence, byte-exact. It begins with two ASCII spaces:

```text
  1. **Route-level still determinism.** Serve the store under WorkbenchStill and repeat-GET every capture route 2 s apart, asserting byte-identical HTML, `script_tags=0`, and still HTML equal to the default HTML minus exactly the live.js tag. The result is a required P0 record. Row 239 committed this check as `checkMarketingStill` (`cmd/ailang-worldd/marketing_still_test.go`; it `t.Errorf`s on `repeat_equal=false`, on `script_tags>0` and on a still page that is not the default minus exactly the live.js tag); the M7c4 drill runs it on its own seeded store before capture, so items (1) and (2) are one run on one store. Command, from the repository root at the capture commit, with both AILANG pins exported: `WORLD_CHROME=<Chrome> WORLD_MARKETING_RECTS=<rects> go test -count=1 -timeout 1200s -v -run '^TestMarketingChromeSixShots$' ./cmd/ailang-worldd > <log> 2>&1`. It passes only if all of these hold: rc 0; one `--- PASS: TestMarketingChromeSixShots`; and `grep -cE 'STILL .* repeat_equal=true script_tags=0 equals_default_minus_tag=true default_live_tags=1 '` = 4 (row 239).
```

- **Scope.** The change is confined to item (1)'s command and pass rule. Items (2)–(4), the rest of §4.4,
  §7 (C12 already matches) and §9/§10 are **not touched**. Whether the orchestration quorum reviews
  this line is the controller's call (row 238 V27).

Acceptance (run all of these; pipe nothing):
- `git diff --numstat -- design_docs/planned/w-marketing-capture-orchestration.md` gives
  `1	1	design_docs/planned/w-marketing-capture-orchestration.md`.
- `grep -cF` gives `1` for each of `checkMarketingStill`, `zz_iter305_probe_test.go` and
  `TestIter305SeededStillCaptureProbe`. The last two are V40's history at `:558`; each is 2 at the base (V07).
- `go test -count=1 -timeout 300s -run '^TestMarketing(Position|Phone)Parity$' ./scripts/marketingcheck`
  gives rc 0. Its `pin-registry` subtests read this doc.
- The pasted command, with the drill's rig env (V12's), gives rc 0 and the count 4.
- Dry run (V16, V17): numstat `1 1`, each `grep -cF` gives 1, parity rc 0. The pass rule gives 4 on
  the V12 log and 0 on the MS1 log.

## 4. Conflict surface, concurrent work, residuals

- **Files.** `marketing_still_test.go` (new), `marketing_chrome_test.go` (+5), one orchestration line.
  No `.ail`, `go.mod`, `.github/`, `tools/launchd/*`, `host/` or `scripts/` change (mutants are reverted).
- **Concurrent work (V15).**
  - 0 open PRs (control: the merged list shows #462 and #461).
  - 7 of 311 origin branches differ on these paths, all only in `render.go`, and all are **merged**
    PRs (#143–#249). Control: the same diff on `95179277` lists the orchestration doc.
  - Rows 240, 246 and 247 name only `scripts/marketingcheck` files.
  - On any new change to these paths, re-base and re-run V07 before M2.
- **R1, UNMEASURED.** GitHub CI does not run the seeded check (V05, V09). A follow-up row could add a
  CI step after `ci.yml:384` running `TestMarketingStillRoutes` with `WORLD_TOOL_AILANG_BIN` set.
  Whether `demo seed` runs on the Linux runner cannot be measured without pushing.
- **R2.** Item (3), the V60 rectangle probe, stays a separate instrument. The drill consumes
  `WORLD_MARKETING_RECTS` (`:73`) and builds its store under `t.TempDir()` (`:162`), so the drill
  cannot hand its store to item (3). This row does not change that.
- **R3, flake.** The helper ran healthy 4× standalone and 1× in the drill (V08, V12, V13, V21). That is
  not a flake-rate measurement.

## 5. Verification Log (base `2f293570`; files under `design_docs/verification/world-iter337/`)

| V | Command | Observed |
|---|---|---|
| V01 | `git rev-parse HEAD origin/dev; go version; ailang --version` ×2; `Chrome --version` (`designer-env.txt`) | `2f293570…` twice; `go1.26.6 darwin/arm64`; `AILANG v0.41.0`; `AILANG v0.54.15`; `Google Chrome 156.0.8078.80` |
| V02 | `git grep -n -E 'repeat_equal\|script_tags' 2f293570 -- '*.go' \| wc -l`; control `git grep -n requireProofPin … \| wc -l` | `0`; control `13` |
| V03 | `grep -n -E 'WORLD_CHROME\|newMarketingChromeRig\|"go", "build"\|capture_marketing.sh\|…' marketing_chrome_test.go` | `52/54` WORLD_CHROME skip; `73` rects env; `122` rig; `125–126` builds; `136` capture; `155` rig type; `159` def; `245` return |
| V04 | `grep -n -E 'workbench-still\|step routes\|step check' scripts/capture_marketing.sh`; `grep -n 'func ShotRoutes\|func runRoutes' main.go` | `116 … serve --db "$db" … --workbench-still`; `129 step routes … routes --manifest`; `141 step check … --rectangles`; `228 func ShotRoutes`; `240 func runRoutes` |
| V05 | `grep -n -E 'GITHUB_ENV\|verify_go.sh\|Install PINNED v0.54.15\|WORLD_TOOL_AILANG_BIN' ci.yml`; `grep -A8 'func cliSeBins'`; `grep -A5 'func demoM5cRealSeed'` | `176 AILANG_BIN=… >> GITHUB_ENV`; `281 run: ./scripts/verify_go.sh`; `384 Install PINNED v0.54.15`; **no** `WORLD_TOOL_AILANG_BIN` line. `cliSeBins` `t.Skipf` if either env is empty; `demoM5cRealSeed` calls `requireProofPin` + `cliSeBins` |
| V06 | `git grep -n -E 'time\.(Now\|Since)' -- 'host/workbench/*.go' 'host/daemon/workbench*.go' \| grep -v _test.go \| wc -l`; control `git grep -c … host/daemon/daemon.go`; `sed -n 203p;281p render.go` | `0` direct calls in those non-test files. This is a textual result, not a transitive audit. Control `7` (it shows the search finds direct calls, not that its scope covers the render path); `{{if .Still}}{{else}}<script src="/workbench/live.js" defer></script>{{end}}`; `func Render(…) { return pageTemplate.Execute(w, p) }` |
| V07 | `grep -cF` anchors in the orchestration doc | `  1. **Route-level still determinism.**` 1 (line 221); `  2. **Six-shot two-pass.**` 1; `zz_iter305_probe_test.go` 2 (221, 558); `TestIter305SeededStillCaptureProbe` 2; `checkMarketingStill` 0; misspelled control 0 |
| V08 | prototype `TestIter337StillRoutes`, both pins (`designer-proto-runs.txt`) | rc 0; 4× `STILL … repeat_equal=true script_tags=0 equals_default_minus_tag=true default_live_tags=1`; `TIMING seed=5.73s build=0.32s still=2.72s total=8.77s`; `--- PASS (8.81s)` |
| V09 | same plus `TestMarketingChromeSixShots`, `AILANG_BIN` only (CI-shaped) | rc 0; `--- SKIP: TestMarketingChromeSixShots`; `--- SKIP: TestIter337StillRoutes` with `WORLD_TOOL_AILANG_BIN set=false` |
| V10 | `go test -v -run '^TestWorkbenchStillAllRegions$' ./host/daemon`, `AILANG_BIN` only | rc 0; `repeat GETs byte-identical on all three routes after 2.006s`; `--- PASS (2.04s)` |
| V11 | `GOCACHE=$(mktemp -d) go build ./scripts/marketingcheck`; `go vet ./cmd/ailang-worldd` with the prototype present | cold build rc 0, 7 s; vet rc 0 |
| V12 | drill with the §2.2 patch, real Chrome, `WORLD_MARKETING_RECTS=world-iter312/r1-rects.json` (`designer-sixshots-with-still.txt`) | rc 0, wall 22 s; 4 `STILL` pass lines (`marketing_chrome_test.go:128`); `--- PASS: TestMarketingChromeSixShots (20.00s)`; patch reverted (`git status` clean apart from the prototype) |
| V13 | 14 `run.sh` runs over `mutate.py` (`designer-mutants.txt`). **Round 1, 11 runs:** (1) healthy, (2) MS1, (3) MS2, (4) MS3 with a comment marker, (5) MS4, (6) MS5, (7) MT1 with `if true {`, (8) MT2, (9) MT3+MS3 with a comment marker, (10) MT4+MS4, (11) MT5+MS5. **Round 2, 3 runs:** (12) MS3 with a `<meta>` marker, (13) MT1 with `[:0]`, (14) MT3+MS3 with a `<meta>` marker | (1) PASS; (2) A2+A3; (3) A2; (4) **PASS**, because html/template strips comments, so it is void; (5) A1; (6) A0; (7) **build failure** (`"os" imported and not used`), not a kill, so void; (8) A1; (9) PASS, void with (4); (10) PASS; (11) PASS; (12) A3; (13) A4 `got 0`; (14) PASS. The §2.5 results are (1)–(3), (5), (6), (8), (10)–(14) |
| V14 | MS1–MS5 against `TestWorkbenchStillAllRegions`; MS5 against `TestWorkbenchStillDefaultLiveScript` and `TestServeWorkbenchStillFlag` | MS1/MS2 `:54 zero script tags`; MS3 `:83 changed HTML beyond the script`; MS4 `:106 repeat GET … varies`; MS5 **PASS** there, but FAIL in `render_test.go:675` and `main_test.go:217` |
| V15 | `gh pr list --state open`; a 311-branch `git diff --name-only 2f293570...$b -- <paths>` under `bash -c`; control `95179277^...95179277`; the `gh pr list --head` state of each hit (`designer-sweep.txt`) | 0 open, control #462/#461; row 240 names `phone.go:157`, `main.go`, `scripts/marketingcheck`; 7 hits, all `render.go` only, all MERGED (#143, #145, #157, #158, #160, #234, #249); control lists the orchestration doc. A first zsh run did not split `$P`, and its empty control exposed that |
| V16 | M2 dry run in a `git archive` scratch copy: paste the fenced line over the anchor line; `git diff --no-index --numstat`; `grep -cF` ×4; parity tests (`designer-m2-dryrun.txt`) | pasted at line 221; `1 1`; `checkMarketingStill`, `zz_iter305_probe_test.go`, `TestIter305SeededStillCaptureProbe` and the item-(2) anchor each 1; parity rc 0, both `pin-registry` PASS |
| V17 | the M2 pass-rule `grep -cE 'STILL .* repeat_equal=true … default_live_tags=1 '` and the `--- PASS` count on the V12 log; the same grep on the MS1 log as a control | `4`, `1`; control `0` |
| V18 | `sed -n '69,71p' marketing_chrome_test.go`; `awk 'NR>=52 && NR<=71 && /Fatal\|Skip/'`; the drill with `WORLD_CHROME` set and `WORLD_TOOL_AILANG_BIN` unset (`designer-r2-premises.txt`) | `:71 t.Fatal("AILANG_BIN and WORLD_TOOL_AILANG_BIN are required")`; Skip/Fatal lines `54` Skip, then `58`, `62`, `67`, `71` Fatal; exercised rc 1 at `marketing_chrome_test.go:71`, `--- FAIL (0.23s)` |
| V19 | `sed -n '114,120p;101p' scripts/capture_marketing.sh` (same file) | the `serve` command spans `116–118`: `--db --bind --workbench-still --ailang-bin --workspace-root --tool-ailang-bin --proof-key-file --proof-required-identity sameRef --proof-required-identity servesEntry`, then `> …daemon-verdict.json &` at `119`; `101` requires `^127\.0\.0\.1:([0-9]+)$` |
| V20 | `grep -n -E 't\.(Logf\|Errorf\|Fatalf)\("(DEFAULT\|STILL)\|t\.Errorf\|TempDir\|"demo", "seed"\|…'` on the iter305 probe `.txt`; `grep -c 't.Errorf'` (same file) | `219 t.Logf("DEFAULT … script_tags=%d")`; `241 t.Errorf("STILL %s status %d")`; `249 t.Logf("STILL … repeat_equal=%v script_tags=%d …")`; `65 t.TempDir()`; `139 "demo", "seed", "--db", db`; `t.Errorf` count `1` |
| V21 | `TestIter337ServeParity` + `TestIter337StillRoutes` (`designer-r2-serve-parity-proto_test.go.txt`, deleted after the run); control `PROTO_SERVE_DEFAULT=1` (`designer-r2-serve-parity.txt`) | rc 0. `LISTENER` ports `59766`/`59784` (default/still), and `60302`/`60316` in the parity test; `SERVE bind=127.0.0.1:60394`; `SERVE-PARITY … equal=true contains_bind=false contains_port=false` on 4/4 routes. Control rc 1: `0/4 routes byte-equal`, bytes `19245/18702/10323/16702` vs `19197/18654/10275/16654` (+48 each, and the tag is 48 bytes). A first control attempt dropped the wrong argv and is recorded as void |

## 6. Quorum r1 dispositions (r2)

| # | Reviewer: objection | Change | Where |
|---|---|---|---|
| 1 | gpt6-1-sol (REJECT): §2.3's "the render path reads no clock" goes beyond V06's two-glob textual search | §2.3's opening replaced with the reviewer's text verbatim, and its heading narrowed. V06's observed column now reads "direct calls … not a transitive audit", and its control is scoped the same way | §2.3, V06 |
| 2a | oc-glm-5-3 (REJECT): "the drill FATALs without both pins" has no V row, and the FATAL path was never exercised | V18 records `:69–71` and the Skip/Fatal order over `:52–71`, and exercises it (rc 1 at `:71`). §2.1 now says it FATALs **when `WORLD_CHROME` is set**, after the skip and the Chrome checks. The superset argument is unchanged | §2.1, V18 |
| 2b | oc-glm-5-3: "mirror the capture script's serve flags" rests on line 116 only | V19 records the whole command (`116–118`). §2.2 gains a flag→field table and says why the bind is not mirrored. V21 measures byte-parity against the real `serve` with those exact flags (4/4), with a failing control (0/4, +48 B) | §2.2, V19, V21 |
| 2c | oc-glm-5-3: the probe's log-only behaviour and own store rest on judge NB4 | Replaced by direct observation: `:241` is the only `t.Errorf`; `:219`/`:249` are `t.Logf`; the store comes from `:65` TempDir and the `:139` seed | §1, V20 |
| 2d | oc-glm-5-3: V13 says 14 runs but enumerates 13 | All 14 are numbered, with which ones were void and why | V13 |
| 3 | gemini-4-argon (PASS): the helper lives only in the sidecar; daemon teardown before capture is not required | §2.4 inlines the four signatures and **requires** per-mode `defer` closure of both daemons and both `httptest` servers before return (prototype `:56`/`:64`/`:70`). It is measured by V21 (a separate `serve` reopened `rig.db`) and V12 | §2.4 |

## 7. Quorum r2 dispositions (r3, controller narrow-refinement carve-out)

Quorum r2 (`quorum-r2.json`; artifact `w-marketing-still-probe-commit-2026-11-03T00-23-09Z.json`):
gemini-4-argon PASS, oc-glm-5-3 PASS (its r1 objections closed), gpt6-1-sol REJECT, no absentees. The
one remaining objection carries a concrete reviewer-authored fix and does not dispute the design
direction (it accepts a separate implementation if justified), so the controller applied that fix and
routes to the planner without a third quorum round.

| Reviewer: objection | Change | Where |
|---|---|---|
| gpt6-1-sol (REJECT): the conflict surface does not justify re-implementing assertions that `TestWorkbenchStillAllRegions` and two other tests already make. Proposed fix: a verified machinery-overlap subsection mapping A0–A4; reuse a shared checker or justify separate code; name what must stay synchronized | §2.6: per-test sharing table; the measured reason none is callable from `package main`; the extraction cost; the decision to keep a separate implementation; the three synchronized invariants with their drift guards; an M1 header-comment requirement | §2.6, V22 |

| V | Command | Observed |
|---|---|---|
| V22 | `grep -c "^func "` on `workbench_still_test.go` (control `demo_proof_test.go`); `git grep -n -F` the live.js tag over `*.go`; `git grep -l "^func <helper>("` for the five helpers; non-test helper packages (`controller-r3-overlap.txt`) | `1` (control `33`); tag at `render.go:203` + 5 test sites (`main_test.go:216`, `workbench_live_test.go:295,304`, `workbench_still_test.go:82`, `render_test.go:674`); helpers in `handlers_test.go` ×2, `workbench_live_test.go`, `workbench_graph_test.go`, `workbench_packet_test.go`; packages `host/proctest`, `host/ptytest` |
