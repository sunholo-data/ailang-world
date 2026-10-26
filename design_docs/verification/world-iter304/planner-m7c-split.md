# Planner: w-marketing-capture M7c split (iteration 304)

Role: planner, lane `opus` (claude-opus-5-5), unattended. Worktree `.wt-world-iter304`, base `039bfa77`. No git write operations. No Chrome run. The design doc was not edited.

## What changed

**`design_docs/planned/w-marketing-capture-sprint-plan.md`**
- The single `## M7c` section is replaced by a shared-rules preamble and five children:
  - **M7c1**: observation comparator (C11).
  - **M7c2a**: Python primitives plus the Go `routes` subcommand.
  - **M7c2b**: the shell end-to-end, with a fake Chrome and a fake daemon.
  - **M7c3**: the frozen inventory guard, against a real store with a genuine commit.
  - **M7c4**: the real Chrome drill. The executor writes it; the controller runs it.
- Other edits:
  - The ordered-milestone table now lists the five children.
  - The header now reads 21 runs / 543 min (was 17 / 445).
  - M8a's dependency is now M7c4.
  - The M8b re-drill row names M7c1's mutant.
- M7a and M7b text is untouched.

**`design_docs/planned/sprint_w-marketing-capture.json`**
- The M7c milestone is replaced by five milestone objects in the same shape and fields as M7b.
- M8a now has `depends_on` `["M7c4"]`.
- The `gates` M7c entry is replaced by five entries.
- The `controller_owned` "Chrome155" item now states the design §2/V29 rule.
- `estimate` is now 21 runs / 543 min.
- `coverage` is now executor_runs 21, planned_tests 39, mutations 52.
- New `planner_decisions.PD7` records the split and the reasons for it.
- The file was rewritten with Python `indent=2, ensure_ascii=False`. I confirmed first that this round-trips the original byte-identically (raw non-ASCII, zero `\u` escapes). `jq -e .` passes. Every diff hunk is in an intended region.

**Departure from the designer's four-way split.** M7c2 is split again, into M7c2a and M7c2b, for size. That orchestration work is:
- three Python modes;
- a Go subcommand;
- shell flag parsing, flow and diagnostics;
- a fake `ailang-worldd` server;
- three mutants.

That does not fit claude-haiku-5-5 in ≤30 min and ≤~250 lines. M7c3 depends on M7c2b only to serialise edits to `main.go`.

**The test for the M7c-inventory mutant moves out of `./scripts/marketingcheck`.** `TestMarketingCaptureInventoryMutation` now lives in `host/daemon` (`package daemon`). That package already has genuine authenticated REST commit seams (`newHandlerDaemon`, `workbenchLiveREST`, `testCommit`, `authHeader`). The guard itself stays production code in `scripts/marketingcheck/main.go`, which the test drives as a built binary.

**The real drill moves to `cmd/ailang-worldd`.** That follows designer X1 option (a): the V11 seams are unexported there.

## Verification rows

| ID | Command | Observed |
|---|---|---|
| P1 | `sed -n 198-245p host/store/context_roots_test.go` | `TestProductionGoSurface` allow-list holds only `scripts/marketingcheck/main.go` and `png.go` (lines 204-206). Every child keeps its production Go in `main.go`, so no allow-list edit is prescribed. Executors are told to add an entry if they create a new non-test `.go` file. |
| P2 | `go test -run '^TestProductionGoSurface$' ./host/store -count=1 -timeout 60s` | `ok … 0.908s` at base. |
| P3 | `grep -n '^func \|^type ' scripts/marketingcheck/*.go` | No `Observe`, `CompareObservations`, `ShotRoutes` or `Inventory` exists. `run` gates on `args[0] != "check"` (`main.go:13`). Refusal stderr format at `main.go:36`. O_EXCL `--out` at `main.go:73`. |
| P4 | `grep -n '"docs"' scripts/marketingcheck/png_test.go` | `png_test.go:800`: `cli/unknown` uses `docs`, expecting rc 2. M7c1 must not add `docs`. **M8a's planned `docs` subcommand will red this subtest; M8a must update it.** |
| P5 | `sed -n 20,35p scripts/marketingcheck/png_test.go` | The `wantCode(t, err, code, assertion)` helper matches on `Refusal.Code`. That is the kill mechanism for the M7c1 mutants. |
| P6 | `sed -n 380,400p scripts/marketingcheck/supervisor_test.go`; `grep -rn TestSubprocessSinkGate` | File-sink precedent at `supervisor_test.go:388`. The gate is at `host/verifygate/subprocess_sink_gate_test.go:223`. |
| P7 | The plan's exact C11 Python one-liner run read-only in the worktree, output to `/tmp/iter304-planner-obs.json` | rc 0. A JSON array of 12 records with keys `lines`, `path`, `sha256`. Paths: `design_docs/world-mission.md`, `go.mod`, `go.sum`, 5 under `tools/launchd/`, and 4 `world/*.ail`. This is the format the M7c1 `compare` subcommand must accept and the `python-parity` subtest must regenerate. Python 3.14.5. |
| P8 | `sed -n 1033-1059p host/daemon/daemon.go`; `sed -n 619-661p host/daemon/handlers.go`; `grep 'type logEntryResponse' -A4` | Routes include `/v1/health`, `/v1/head`, `/v1/log` (`from`, `limit` clamped 1..500, default 100) and `/v1/objects/{ref}`; there is no list-all-objects route. `items` holds `{header, entryHash, transitionRef}`. The range **stops at the first missing index** (`handlers.go:656-657`). |
| P9 | `sed -n 74-90p host/daemon/workbench_live_test.go`; `sed -n 1053-1066p host/daemon/daemon_test.go` | `workbenchLiveREST` POSTs genuine authenticated commits through `requestRecorderAuth`. `authority.Mint` is exported (`host/authority/mint.go:28`). |
| P10 | `sed -n 245-262p cmd/ailang-worldd/main.go`; `grep -n sameRef docs/QUICKSTART.md` | `--proof-key-file` "Needs --ailang-bin and at least one --proof-required-identity". QUICKSTART 1557 and 1576-1577 name `sameRef` and `servesEntry`. |
| P11 | `sed -n 1,64p` and `380,512p scripts/marketingcheck/png.go` | `Shots` at lines 49-53. `BuildCaptureManifest` at 462 derives the routes. `CheckScheme` refuses `wrong-scheme` (272-280). `CheckPass` is at 323. |
| P12 | `cat -n scripts/marketing_capture.py` | Constants at 32-36. `shot()` argv at 91-97 always passes a scheme flag. Exit and verdict handling at 236-237. `daemon` mode is file-ready only (152) and takes at most one shot (164-165). |
| P13 | `grep -n '155.0.8059.41' design_docs/planned/w-marketing-capture.md` | Rows V25 (line 74) and V29 (line 78) give the real version strings that the M7c2a `version` subtest copies (155.0.8059.41 and 156.0.8078.47). |

## Base focused gate

| Command | rc | Wall | Detail |
|---|---|---|---|
| `go test ./scripts/marketingcheck -count=1 -timeout 170s -v` | **1** | 18 s (package 17.568 s), load 8.24 | 6 of 7 tests PASS. **FAIL: `TestMarketingSupervisorFakeDaemonCleanup/shell-hup`** in 0.53 s, `supervisor_test.go:469: hangup trap must exit 129 after supervised cleanup: exit status 1` |
| `go test -run '^TestMarketingSupervisorFakeDaemonCleanup$' ./scripts/marketingcheck -count=1 -timeout 90s -v` ×3 | 0, 0, 0 | 5 s, 6 s, 6 s, load 13–15 | shell-hup PASS in 0.46–0.55 s each time |

This is a **second failure mode of the carried shell-hup flake**. The carried note said it exceeds its 10 s bound (29 s). This run instead got the wrong exit code fast.

Hypothesis, not established: the subtest passes unit-control waits `--term-wait 0.15 --kill-wait 0.3` (`supervisor_test.go:432`). Under load, the reap after KILL can overrun 0.3 s. Python then returns `CleanupFailed` (rc 1), and the trap's `wait` turns status into 1 (`capture_marketing.sh:13`). Those 0.3 s waits are below the ≥10× rule.

No child owns `supervisor_test.go`. I recommend the controller route a test-bound fix: raise the wait values, or assert `CleanupFailed`-tolerant. It must not be absorbed into M7c1.

## Open risks / controller decisions

1. **The M7c4 dependency gates M8a.** M7c4 completes when its executor deliverable lands: SKIP when unset, FAIL when set-but-invalid, plus the fake companion. The real run is a controller rehearsal that needs measured rectangles and is not a completion condition. If the controller wants M8a to wait for the rehearsal itself, it must say so.
2. **The base package gate is red** on the shell-hup flake (above). Executors running `go test -race ./scripts/marketingcheck` will see it intermittently. They are told to report it as carried and never edit `supervisor_test.go`.
3. **M8a will red `cli/unknown`** (P4) when it adds `docs`. M8a's file list does not include `png_test.go`, so its planner or controller must add it.
4. **The M8b focused list runs `TestMarketingCaptureFakeChromeSixShots` with `-timeout 120s`.** M7c2b prescribes 300 s, because the test spans two shell runs at ≤120 s each. I did not edit M8b's text; it should be aligned when M8b is next planned.
5. **Inventory paging stops at a gap** (P8). The seeded store's log must be contiguous; M7c4 inherits this assumption.
6. **M7c2b's free-port probe** (listen on :0, then close) is racy under parallel load. It is acceptable for a fake daemon, but a bind failure is UNINFORMATIVE, not PASS.
7. **M7c4's "rectangles required" guard lives in test code.** It is listed as a review check, not a mutant.
8. **A concurrent design quorum is running on §2.** This plan follows the doc as it stands (major ≥155, identical before and after, recorded). If the quorum changes the rule, M7c2a, M7c2b and M7c4 need the edit.
