# Planner receipt — iteration 305, row 150 (w-marketing-capture-orchestration r2 → plan)

Role: planner (`opus fail-closed:env-pin`). Worktree `.wt-world-iter305-row150-orch`, branch `design/w-marketing-capture-orchestration`, HEAD `ff16bf65936f722021ef2cc0cad2e2faa53821dc` (`git rev-parse HEAD`). No git writes.

## What changed

**`design_docs/planned/w-marketing-capture-sprint-plan.md`.** 29 spots are marked `(amended iteration 305 per w-marketing-capture-orchestration §5)` or `§4.2`.
- **Header:** the orchestration doc governs M7c0–M8b, and its §5 table wins on disagreement. 22 runs / 563 min.
- **Environment:** the export at line 15 becomes `v0.54.15`. The pin prose at line 19 now names `4b686d8e` and cites `workspace.go:52`, replacing `:44`, which is a comment.
- **Landed M2b:** the two v0.52.1 lines (356, 394) are annotated as historical, not rewritten.
- **Ordered table:** M7c0 row added.
- **"## M7c — split":** iteration-305 governance paragraph and new chain. C13 carve-out, M7c0 after R0, M7c2a after M7c1 and M7c0, R1/P0/P3/A1 controller-owned.
- **New "## M7c0 — Shot geometry":** files, constants citing V51/R0, `TestMarketingShotGeometry` code, canonical list `png_test.go:462`, mutant M7c0-height, R0 dependency, acceptance commands.
- **"## M7c1": rewritten to be prescriptive.**
  - HOLD lifted; exact files and do-not-touch list; verified-facts table.
  - Interface; CLI argv/stdout/stderr/rc table.
  - Exact code blocks A, B and C, from the prototype below.
  - Subtest/assertion table; gate commands with expected results.
  - Mutants with the exact edit, the command and the observed red line.
  - C11 observation steps.
- **M7c2a:** routes carry the M7c0 dimensions. The `routes` subtest asserts literal 3700/2100/3000. The version strings are kept. The dependency now includes M7c0.
- **M7c2b:** `serve-argv` asserts that `--tool-ailang-bin` is followed by the operator's path. The shell must carry no hard-coded version (the grep must print 0).
- **M7c3:** "unchanged" note: the seeded log is contiguous (V40).
- **M7c4:**
  - local rig helper in `marketing_chrome_test.go`, with the probe reference; `demo_test.go` NOT edited;
  - corrected cites;
  - pinned env on the acceptance block;
  - completion semantics (a)–(e).
- **M8a:**
  - files: + `png_test.go`, + `png.go`, + `main.go`;
  - v0.54.15 environment;
  - captions (tall framing, narrowed provenance claim);
  - threshold ownership from R1/V57, with the literal cites and the gate.
- **M8b:** fake-Chrome gate `-timeout 300s`; `TestMarketingShotGeometry` gate added; pinned env.
- **Final disposition:** statuses (M1a–M7b completed with PRs); C8′/C12/C13/C14 belong only to controller/attended steps.

**`design_docs/planned/sprint_w-marketing-capture.json`.** The round-trip `json.dumps(indent=2, ensure_ascii=False)` was byte-identical before editing, so formatting is preserved. No new field names: only existing keys are used.
- **Top level:**
  - `estimate`;
  - `environment` (tool pin → v0.54.15);
  - `authority` (+ orchestration doc);
  - `controller_owned` (+ R0, R1, P0/P3, A1);
  - `planner_decisions.PD8`;
  - `coverage` (executor_runs 22, planned_tests +1, mutations +1, criteria + C8′/C12/C13/C14).
- **`gates`:** + M7c0; M7c1 and M8b command lists updated.
- **M1a–M7b:** already `"status": "completed"`. `completion_evidence` was added where it was absent, as `{pr, merge_commit, merged_utc, iteration, judge, source}`. Where it existed (M1a, M1b, M2a), only the missing keys were added (M1a gains `pr: 264`). PRs: #264, #266, #267, #273, #275, #277, #283, #289, #295, #353, #354, #359, #360, #361 (doc §2).
- **New M7c0 record:** the same 18 keys as M7c1, inserted before M7c1.
- **M7c1–M8b:** amendments mirrored. `files`, `depends_on`, `implementation`, assertions, commands, mutations and `measured_baseline` were updated as applicable.
- `jq -e . design_docs/planned/sprint_w-marketing-capture.json` → rc 0.

## M7c1 readiness: prototype outside the repo

Scratch tree `/tmp/m7c1-proto-iter305`. It holds copies of `scripts/marketingcheck/*`, `scripts/marketing_capture.py` and `scripts/capture_marketing.sh`, with symlinks to the worktree's `go.mod`, `go.sum`, `world`, `design_docs` and `tools`. Nothing was written into the worktree's source tree.

| Command (in the scratch package) | Result |
|---|---|
| `go vet .` | rc 0 |
| `go test -run '^TestMarketingProtectedObservationComparator$' . -count=1 -timeout 120s -v` | `--- PASS` top-level + 11 subtests, 0.80 s, load 1.97 |
| M7c-observation (`if b.SHA256 != a.SHA256 {` → `if false {`), `-run '…$/^same-line-count-change$'` | FAIL: `same-line-count changed copy must red: want Code=observation-digest, got <nil>` |
| M7c-line-observation (`if b.Lines != a.Lines {` → `if false {`), `-run '…$/^digest-equal-line-change$'` | FAIL: `line-count mismatch must refuse: want Code=observation-lines, got <nil>` |
| M7c1-set (the set guard → `if false {`), `-run '…$/^deleted$'` | FAIL: `deleted protected copy must red: want Code=observation-set, got observation-digest tools/launchd/testdata/deep/y.md: sha256 3c17e3eb… became ` (no panic: path-keyed) |
| restore `cp` + `cmp main.go /tmp/m7c1-main.orig` → "restored"; the three commands again | 3× PASS |
| `go test -run '^TestMarketingTwoPassInventory$' . -count=1 -timeout 120s -v` | PASS incl. `cli/unknown`, `cli/no-command`, `cli/missing-flags`, `cli/help` |
| `go test -race . -count=1 -timeout 240s` | `ok … 32.875s` at load 1.49. The first attempt failed only because the scratch tree lacked `../capture_marketing.sh` (`shell-*` "readiness file not published"); this is a prototype artifact, fixed by copying the script. |
| Reconstruction: apply the plan's block A to original `main.go` lines 1-16 + the `fs.Parse(args)` edit + block B appended; block C as the test | `cmp` vs prototype: **MAIN-IDENTICAL**, **TEST-IDENTICAL** |
| Python constant in block C vs plan heredoc lines 120-123 | `verbatim: True` |

M7c0 scratch probe (same tree): with `Shots` 3700/2100/3000 and `:462` updated, `go test . -count=1 -timeout 240s -v` gave all 8 tests PASS (17.7 s, load 1.67). With the packet height mutated back to 1400, `TestMarketingShotGeometry/packet.png` FAILed with `packet height must equal the rule value from measuredPacketBottom: got {… Height:1400 …}, want 1440x3700`, and `canonical-shots` FAILed with `canonical shot packet.png: …`. Restored green.

## Cited file:line verified at HEAD `ff16bf65`

| Command | Output (abridged) |
|---|---|
| `cat -n scripts/marketingcheck/main.go` | `13 if len(args) == 0 \|\| args[0] != "check" {`; `17 fs := flag.NewFlagSet("check", …)`; `23 if err := fs.Parse(args[1:]); …`; `36 fmt.Fprintf(os.Stderr, "marketingcheck: %s %s: %s\n", …)`; `73 os.OpenFile(out, os.O_WRONLY\|os.O_CREATE\|os.O_EXCL, 0600)`; 83 lines |
| `sed -n 15,60p scripts/marketingcheck/png.go` | `19 type Refusal struct{ Code, Name, Detail string }`; `22 func refuse(…)`; `27 // Content thresholds are provisional until M7c measures real renders.`; `28-35` MinColors 32 … MinRectStdDev 4.0; `49-53` Shots (packet 1400, graph 1100, provenance 1400) |
| `sed -n 1,40p png_test.go`; `awk 140-200`; `awk 455-470`; `awk 725-812` | `23 func wantCode`; `30 func pngOK`; `147 image.Rect(0, 0, s.Width, s.Height)`; `164 … W: min(s.Width-64, 900)`; `198 func demoFixture`; `461 t.Run("canonical-shots"`, `462 expected := []Shot{… {"packet.png", 1440, 1400, …} …}`; `734-741` go build -o bin with a 60 s ctx; `800 {"unknown", []string{"docs"}, 2}` |
| `awk 307-372 png_test.go`; `grep -n '1400\|1100\|Height'` | `321`, `325-326`, `336-337` 1440×1100 (workbench-light); `357` X: 1400; `366`, `383` 1440×1100. None tied to tall shots |
| `awk 396,438-441,444,455 png_test.go` | `396 "32 RGB colors rejected"`; `438-441` CheckScheme 160/96 accepted, 159/97 rejected; `444 for _, delta := range []int{7, 8}` |
| `sed -n 20-35p; 56p; 212p; awk 380-395` supervisor_test.go | `25 func pythonSupervisor`, `29 t.Fatal("python3 is mandatory: ", err)`; `56 func fakeScript`; `212 func assertDead`; `388 // A file sink: a bare buffer races os/exec's copier (TestSubprocessSinkGate).` |
| `awk 196-246 host/store/context_roots_test.go` | `198 func TestProductionGoSurface`; `205-206` main.go/png.go allow-listed |
| `awk 215-232,296-345 host/verifygate/subprocess_sink_gate_test.go` | `223 func TestSubprocessSinkGate`; flags bare buffers on `Start()`; `Run()` with a buffer is a negative control (`:302`) |
| `awk 30-37,88-98,230-241 scripts/marketing_capture.py` | `32-36` SHOT 30 / PASS 300 / STARTUP 30 / TERM 5 / KILL 5; `91-97` shot argv; `236-237` one JSON line, rc 0 only for ok |
| `cat -n capture_marketing.sh \| sed -n 18,40p`; `sed -n 80,85p` | ten flags at 25-34 incl. `--tool-ailang-bin`; `84 … capture refused`, `85 exit 1` |
| `grep -cE 'v0\.[0-9]+\.[0-9]+' scripts/capture_marketing.sh`; control: the same on the orchestration doc | `0`; control `22` |
| `grep -n '^func demoM5cRealSeed\|^func newDemoApprovalRig\|return r, m, rig.token\|^type demoSeRig' cmd/ailang-worldd/demo_test.go` | `548 type demoSeRig`; `561 func newDemoApprovalRig`; `1727 func demoM5cRealSeed`; `1830 return r, m, rig.token`. **The plan's `demo_test.go:1722-1760` was stale and is corrected.** |
| `grep -n '^func testMintEnv' cmd/ailang-worldd/session_test.go` | `29` |
| `grep -n 'ToolBinaryRelease' host/daemon/workspace.go` | `44` comment, `52 var ToolBinaryRelease = pins.VersionLine(pins.RoleTool)`. **The plan's `:44` was corrected to `:52`.** |
| probe grep (`iter305-seeded-still-chrome-probe_test.go.txt`) | `64 // --- rig: newDemoApprovalRig's body (demo_test.go:561), keeping the daemon handle.`; `139 demoRun(… "demo", "seed" …)`; `143-144 srv.Close(); d.Close()` |
| `awk 654-658,65-73 handlers.go`; `grep -n '^func newHandlerDaemon\|^func authHeader\|^func testCommit\|^func requestRecorderAuth' host/daemon/*_test.go`; `grep -n '^func workbenchLiveREST'` | `656 if !ok {` `657 break`; `65 type logEntryResponse`; handlers_test.go `21`, `47`, `59`, `78`; workbench_live_test.go `74` |
| plan C11 heredoc at the worktree root (Python 3.14.5) | 12 records (world-mission.md 6258 lines, go.mod 21, go.sum 59, 5 launchd files, 4 world/*.ail); `find tools/launchd world -type l \| wc -l` = 0 |
| `git grep -c 'func TestMarketingProtectedObservationComparator' -- '*.go' \| wc -l`; control `… TestMarketingTwoPassInventory` | `0`; `scripts/marketingcheck/png_test.go:1` |
| `go test -run '^TestMarketingTwoPassInventory$' ./scripts/marketingcheck -count=1 -timeout 120s -v` (worktree) | PASS, ok 1.325 s, load 8.25 |
| `go vet ./scripts/marketingcheck` (worktree) | rc 0 |
| `go test -run '^TestProductionGoSurface$' ./host/store …` (worktree) | PASS; `host/cmd=165 allow-listed reproducers=4` |
| `go test -run '^TestSubprocessSinkGate$' ./host/verifygate …` (worktree) | PASS; `violations=0 filesParsed=476 startSites=23` |
| `$HOME/.pinned-ailang-tools/v0.54.15/ailang --version` (40 s bound) | `AILANG v0.54.15`, commit `0384450` |
| `python3 --version`; `go version`; `grep '^go ' go.mod` | 3.14.5; go1.26.6 darwin/arm64; `go 1.26.6` |
| plan heredoc position after edits (`awk 118-125`) | lines 119-124, as cited by §M7c1 |

## Corrections beyond the doc's §5 table

1. **The M7c1 `cli` subtest depended on another subtest.** Case (vi) used "python one-liner output file from python-parity", which contradicted "every subtest builds its own fresh tree". It now calls `pythonObserve` itself, and the cases are reordered so that (vi) runs before the byte flip.
2. **The set guard in M7c1 could not see duplicates on both sides.** `slices.Equal` alone misses a path that both sides repeat, so a `slices.Compact` length check joins the same single guard. The refusal detail references both slices, so the `if false` mutant compiles. This was verified in the prototype.
3. **The M7c1 size estimate was low.** It said "~90 added lines"; the measured `main.go` goes from 83 to 297 lines.
4. **Two stale cites corrected:** `demo_test.go:1722-1760` → `:561` / `:1727` / `:1830`, and `workspace.go:44` → `:52`.
5. **M7c1 gates gained two commands:** `gofmt -l` and `TestSubprocessSinkGate`.

## Decisions flagged, not made

- **M8a vs R1.** Doc §5 says "M8a does not depend on R1", yet M8a must apply R1's threshold V-row (V57). The plan's interim rule is that the threshold step is reported DEFERRED if no R1 V-row exists at dispatch; the controller then either holds M8a or runs a follow-up. **Controller to decide.**
- **M8a `main.go`.** The doc adds `png_test.go` (and `png.go` via the threshold step) but not `main.go`. The `docs` subcommand cannot live anywhere else without an allow-list edit, so the planner added `main.go` to M8a's files. **Controller to confirm.**
- **M7c2a depends on M7c0** as well as M7c1, because its routes assert the M7c0 dims. This is a planner sequencing choice.
- **M7c0 `depends_on`.** It is `["M7b"]`, with the R0 gate stated in `implementation`, because R0 is not a milestone id and tooling may resolve `depends_on` ids. Plan §M7c0 states the R0 dependency in full.
- **The landed M2b text is annotated, not rewritten.** Its v0.52.1 lines record what was built then.
