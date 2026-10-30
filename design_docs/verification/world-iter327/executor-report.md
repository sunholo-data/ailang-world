# Executor report — w-finish-scope-expiry-boundary (row 236), iteration 327

**Worktree:** `/Users/voightkampff/dev/sunholo-data/ailang-world/.wt-world-iter327`
**Branch:** `sprint/world-iter327-finish-scope-expiry`
**Base:** `f98ef9ef` (Go tree = `4d43d8e3`)
**Result:** M1 committed (`d7a54db5`). M2 **stopped on a failure**: AC-4 `go test ./...` has a `host/daemon` failure (default 10-minute package timeout) and a `host/broker` failure. Per the plan, the `host/daemon` failure is a stop.
**Not touched:** `design_docs/world-mission*.md`, `tools/launchd/*`, `~/.ailang/state/mission-v1*`, the main checkout. Nothing pushed, no PR, no merge.

## Pins

| Pin | Observed |
|---|---|
| go | `go1.26.6 darwin/arm64` |
| AILANG_BIN | `AILANG v0.41.0` (commit `24ee108`) |

## Base sha256 (pristine, before the first edit)

Matches the plan §3 base column for all three files.

| File | sha256 |
|---|---|
| `host/daemon/daemon.go` | `5b09012f8d19…` (matches) |
| `host/daemon/episode_finish.go` | `ea44759896b5…` (matches) |
| `host/daemon/episode_finish_test.go` | `094884e036f7…` (matches) |

Tool scripts extracted with the plan's `x()` function: all four hashes equal §M1.1 (`allmuts.sh` `96e27ec7…`, `harness.sh` `c2531dca…`, `mut.py` `715ab302…`, `runmut.sh` `bd8aed15…`). Diff block `d3661012…` matches.

## M1 — red at base, the diff, green on the fix, mutation drill

| AC | Command (see evidence file) | Expected | Observed | Verdict |
|---|---|---|---|---|
| AC-1 base plain | `harness.sh` `:164`, 20, `^TestEpisodeFinishScopeExpiry$` | rc 1, FAIL 20 | rc=1 PASS=0 FAIL=20; first `episode_finish_test.go:171: expiry …: 200`; sha `094884e036f7` equal | PASS |
| AC-1 base race | same, `-race` | rc 1, FAIL 20 | rc=1 PASS=0 FAIL=20; sha equal | PASS |
| AC-1 base control | `SPIN` `> 100_000_000` | rc 0, PASS 20 | rc=0 PASS=20 FAIL=0 | PASS |
| AC-1b base plain | `harness.sh` `:79`, `^TestEpisodeFinishOperatorProofAC11$` | rc 1, FAIL 20 | rc=1 PASS=0 FAIL=20; first `episode_finish_test.go:94: 200 …` | PASS |
| AC-1b base race | same, `-race` | rc 1, FAIL 20 | rc=1 PASS=0 FAIL=20 | PASS |
| AC-1b base control | `SPIN` control | rc 0, PASS 20 | rc=0 PASS=20 FAIL=0 | PASS |
| apply / hashes | `patch -p1 --dry-run` then apply | dry rc 0, 3 files, no offset/fuzz | dry rc 0; apply rc 0; hashes `e6b70d7e3f3f` / `e429375dc3e0` / `b872a38cc091` | PASS |
| G1 | `gofmt -l host/daemon/` | empty | empty | PASS |
| G2 | `go vet ./host/daemon/` | rc 0 | rc 0 | PASS |
| G3 plain | `go test -timeout 280s ./host/daemon/ -run ^TestEpisodeFinish -count=1 -v` | rc 0, 5 PASS | rc 0; 5 PASS (AC11, ScopeExpiry, OperatorProofExpiryBoundary, RecordRechecksScopeExpiry, AuthorizationPrecedesBody) | PASS |
| G3 race | same, `-race` | rc 0, 5 PASS, 0 FAIL/race | rc 0; PASS=5 FAIL/RACE=0 | PASS |
| AC-1 fix plain | `harness.sh` `:179`, 20 | rc 0, PASS 20 | rc 0 PASS=20 FAIL=0, 19.8 s; sha `b872a38cc091` equal | PASS |
| AC-1 fix race | same, `-race` | rc 0, PASS 20 | rc 0 PASS=20 FAIL=0, 24.0 s | PASS |
| AC-1b fix plain | `harness.sh` `:79`, 20 | rc 0, PASS 20 | rc 0 PASS=20 FAIL=0, 21.3 s | PASS |
| AC-1b fix race | same, `-race` | rc 0, PASS 20 | rc 0 PASS=20 FAIL=0, 23.0 s | PASS |
| AC-2 | `allmuts.sh` (21 mutants + S0 control) | 21 × `rc=1 build_err=0` named killer; S0 `rc=0`; every `sha` equal; `scratch removed: yes` | see tally below; S0 rc 0; all 22 `sha` pairs equal; `scratch removed: yes`; `shasum -c` OK ×3 | PASS |

**AC-2 tally (21 mutants, all killed as named, build_err 0):** P1–P6 (T1 `ScopeExpiry`), Q1–Q4, Q6 (T3 `RecordRechecksScopeExpiry`), R1, R3, R4 (T2 `OperatorProofExpiryBoundary`), R2 and R6 (AC11 and cascade, plan D4), S1 and S2 (T1+T2+T3), S3 (T3), W1 (AC11 `:106: 403`), TM1 (T1 `:187`). S0 control `rc=0 fails=[]`. Restore check: every `sha=` pair equal (`e429375dc3e0` for P/Q/R/S, `e6b70d7e3f3f` for W1, `b872a38cc091` for TM1).

Restore check: the sha pairs are equal on all 22 `sha=` lines (20 on `episode_finish.go` `e429375dc3e0`, 1 on `daemon.go` for W1 `e6b70d7e3f3f`, 1 on the test file for TM1 `b872a38cc091`). No protocol error.

## M1 commit

`d7a54db5 fix(row 236 M1): finishNow clock seam; exact-edge finish expiry tests`: 3 files, +84 / −13, only the `host/daemon` Go files. Co-Authored-By line: the Haiku 5.5 line from the harness attribution rule (the brief said Opus 5.5; see Deviations).

## M2 — AC-3 (package), AC-4 (full gate)

| AC | Command | Expected | Observed | Verdict |
|---|---|---|---|---|
| AC-3 vet | `go vet ./host/daemon/` | rc 0 | rc 0 | PASS |
| AC-3 plain | `go test -timeout 880s ./host/daemon/ -count=1` | ok, 0 FAIL | `ok host/daemon 441.258s`, 0 `--- FAIL` | PASS |
| AC-3 race | `go test -race -timeout 1450s ./host/daemon/ -count=1` | ok, 0 DATA RACE | `ok host/daemon 524.659s`, 0 DATA RACE, 0 FAIL | PASS |
| AC-3 race count=200 | `go test -race -timeout 880s ./host/daemon/ -run ^TestEpisodeFinish -count=200` | ok | `ok host/daemon 52.032s` | PASS |
| AC-3 worldd Finish | `go test -timeout 560s ./cmd/ailang-worldd/ -run Finish -count=1 -v` | rc 0, ok, 0 FAIL | rc 0, `ok 9.307s`; 55 `--- PASS` lines (includes subtests under `-v`; the planner's 19 was without `-v`), 0 FAIL | PASS |
| AC-4 full gate | `go build ./... && go test ./...` (default `-timeout` 10m per package) | rc 0 | **rc 1, elapsed 606 s.** `go build` rc 0. See failures below. | **FAIL** |
| AC-4 verify_ail | `./scripts/verify_ail.sh` with `AILANG_BIN` | rc 0 | rc 0; "verify gate PASSED: 16 required identities verified, 40 named tests pass, 503 se-tools named tests pass" | PASS |

### AC-4 failures (stopped; not classified)

1. **`host/daemon` — STOP per plan.** `panic: test timed out after 10m0s`, `FAIL github.com/sunholo-data/ailang-world/host/daemon 603.261s`. At the timeout the running test was `TestWorkbenchPacketRosterTestedAndProven` (3m54s, subtest `test-record-size` running), in `workbench_grade_test.go:65` (`gradeDaemonConfig`). The package alone ran 441 s (plain) and 525 s (race) in AC-3 under concurrent load from other sessions (`go test -race ./cmd/ailang-worldd ./host/daemon ./host/surfacepage` and `daemon.test` processes from other trials were running on this machine). The gate command omitted `-timeout`, so the 10-minute default applied. This is most likely a **timeout of the gate invocation, not a test assertion**: no assertion failure appears in the `host/daemon` output. It is still a `host/daemon` failure under the plan's rule, and I did not re-run or patch it. The controller should decide whether to re-run the full gate with `-timeout` raised (for example `go test -timeout 1500s ./...`). I did not run that.
2. **`host/broker` — not classified.** `--- FAIL: TestExecProbeForeignListenerOnTheReversePortDoesNotRefuse` (11.60 s; subtest `chatty`), and `--- FAIL: TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential` (65.66 s; subtests `host/agentrun/agentrun.go`, `host/checkrun/checkrun.go` driver0). Neither name contains `Finish`, and neither package imports `host/daemon`. Per §M2.3 these need the alone re-run ×3 and base baseline; I did not run them because the `host/daemon` stop takes precedence. Not yet known whether this is load or pre-existing.

Full output: `executor-ac4-full.txt`, `executor-ac4-full-summary.txt`. verify_ail output: `executor-ac4-ail.txt`.

## Deviations from the plan

- **AC-4 gate invocation:** the plan's §M2.2 command is `perl alarm 2700 bash -c 'go build ./... && go test ./...'` with no `-timeout`. The default 10-minute per-package timeout applied. The plan's own §M2.3 step 2 uses explicit `-timeout`; the gate command did not. This is what caused the `host/daemon` timeout. The plan's `host/daemon` timing (`335.8 s` plain, V-P8) is close to the 600 s default under load.
- **Co-Authored-By:** the brief named Opus 5.5. The harness attribution rule names Haiku 5.5 and says it replaces earlier guidance. I used the Haiku 5.5 line.
- **AC-3 worldd count:** the plan's `19 --- PASS` figure was without `-v`. With `-v` the run reports 55 PASS lines including subtests. Still 0 FAIL.
- **Mutant sha tally:** all pairs equal; the count differs from the plan's per-file prose only in wording.
- The background tool notification in M2 made some runs asynchronous; I polled them with bounded waits. No second package-scale test ran at the same time as mine.

## Worktree state at the stop

M1 committed. M2 evidence files written under `design_docs/verification/world-iter327/executor-*`. `git status --porcelain` at the stop is listed in the message to the controller. No mutant residue: `host/daemon/zz_mut_wall.go` absent; `shasum -c` against `fixed.sha` OK for all three files.
