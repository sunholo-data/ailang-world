# w-daemon-race-package-timeout-headroom — the 8-minute cap encodes the fast runner; one in three race attempts draws a slower one

- **Status**: Planned (r0, iteration 443, unattended designer `claude:claude-opus-5-5`)
- **Item**: charter queue row 267, `w-daemon-race-package-timeout-headroom`, clause-2 hygiene. Picked under the
  standing red-dev rule (D-WORLD-121 OPEN, default B: Gate 1 found `dev`'s tip `6420426e` red at attempt 1)
- **Estimated**: one milestone, ~0.5 day, ~250 changed lines (§8)
- **Measurement base**: `6420426e` (dev tip, clean worktree `sprint/world-iter443-row267-race-headroom`), 2026-11-27
- **Instruments**: CI job logs of the `go-race-verify` job (ubuntu-24.04, Azure eastus2 hosted runners) for the
  last 40 workflow runs, 53 attempt files under `~/.ailang/state/world-iter443/racelogs/` (fetched by the
  controller); the controller's `race_table.py`/`race_table.txt`; this session's re-derivation scripts
  `~/.ailang/state/world-iter443/design-scratch/{legs,classes,atcap}.py` (run with `python3 -I`; outputs
  `classes.out`, `atcap.out`); local `go1.26.6 darwin/arm64` source of `cmd/go` for the `go test` kill path.
  No `.ail` is touched, so the pinned `ailang` matters only to the gate runs (§6).
- **Files touched**: `scripts/verify_go.sh` (~+45/−14), `scripts/test_race_leg_bounds.sh` (new, ~75),
  `host/verifygate/race_job_gate_test.go` (~+95/−1), `host/verifygate/mission_config_gate_test.go` (~+12/−4),
  `.github/workflows/ci.yml` (~+5/−3). No Go production code, no `.ail`, nothing under `tools/launchd/`.

Quantities carry their instrument. `VERIFIED BY CONTROLLER` marks inherited controller facts; `D<n>` marks this
session's re-derivations, each with its command and output in §10.

## 1. Problem as measured (and where row 267 is wrong or stale)

The race job runs `go test ./... -count=1 -race -timeout 8m` (`scripts/verify_go.sh:355`) inside a Python
wrapper that SIGKILLs the whole leg at `VERIFY_GO_RACE_TIMEOUT_S` = 1500 s (`verify_go.sh:222`, `ci.yml:525`),
under a step ceiling of `timeout-minutes: 35` (`ci.yml:522`) [D1]. The 480 s per-package bound is the one that
fires.

**The measured population** (51 complete, distinct race-job attempts, 2026-11-25 → 2026-11-27; the controller's
53 rows include one in-progress attempt and one byte-identical duplicate, both dropped) [D2, D3]. I index runner
speed per attempt as `k` = the summed `-race` seconds of every package that never reaches the cap (all packages
except `cmd/ailang-worldd`, `host/daemon`, `host/verifygate`), divided by that sum's median on the normal
cluster (827.9 s):

| class | attempts | `k` | `host/store` s | pre-leg s | leg wall s | attempts with a timeout panic |
|---|---|---|---|---|---|---|
| normal | 24 | 0.986–1.010 | 58.2–67.5 | 34–38 | 643–673 | 0 |
| mid | 13 | 1.057–1.321 | 74.7–97.8 | 34–43 | 676–875 | 4 |
| slow | 14 | 1.403–1.585 | 82.4–102.2 | 46–55 | 930–1034 | 14 |

[D3, `classes.out`]. **18 of 51 attempts (35 %) contain at least one `panic: test timed out after 8m0s`; there are
29 panics in all: `host/daemon` 18, `cmd/ailang-worldd` 6, `host/verifygate` 5** [D3, D5]. On the normal cluster
the three heavy packages read `host/daemon` 353.7–360.3 s (median 356.8), `cmd/ailang-worldd` 259.5–284.0 s
(median 266.3), `host/verifygate` 321.5–334.9 s (median 328.6) [D3].

**Corrections to the row and the brief:**

1. *"Typical `host/daemon` race time ~342 s"* (row 267, written iteration 430) is **stale**: it is 356.8 s median
   on the normal cluster now (24 attempts) [D3].
2. *"A runner about 1.4× slower"* holds only for the slow tail. The cap also fires at `k` 1.25–1.32 (4 of 13 mid
   attempts, e.g. run 45782990414 a1, `k` 1.250, `host/daemon` F480.077) [D3]. The cause is not one second
   machine class.
3. *"The runner population looks BIMODAL"* (brief, controller reading) is **not what the data shows**. There is
   one tight normal cluster (24 attempts, `k` spread 2.4 %), then a broad continuous tail from 1.057 to 1.585
   (27 attempts), with no gap between mid and slow [D3]. What *is* bimodal is the normal-versus-not split: `host/store`
   jumps from ≤ 67.5 s to ≥ 74.7 s.
4. *"Is `host/store` a usable in-run measure of runner speed?"* **No**, for two measured reasons. (a) It separates
   normal from not-normal but does not order the tail: `host/store` read 82.4 s at `k` 1.525 (run 45743770171 a1)
   and 97.8 s at `k` 1.321 (run 45936239736 a1) [D3]. (b) It arrives too late. `host/store` completes 459–646 s
   into the leg, but the heavy packages start 43–166 s in (`cmd/ailang-worldd`, `host/daemon`) and 322–554 s in
   (`host/verifygate`) [D4]. `go test` fixes `-timeout` when it launches each test binary, so a store reading cannot
   set the bound of a binary that is already running (§3 B).
5. *"Moving the slowest `host/daemon` tests out of the shared bound"* (row 267's second UNMEASURED candidate)
   **has no single culprit to move.** When the daemon cap fires, the test it interrupts has itself run only 1–25 s
   (`TestWorkbenchTimelinePaging` 15 of 18 daemon panics, `TestEpisodePolicyProtectedRebased` 2,
   `TestWorkspacePackageCacheCoversLock` 1). The six `cmd/ailang-worldd` panics fall on six different tests, each
   at 0 s elapsed [D5]. The package's total is what overruns, not one test.
6. *"Iteration 442 measured a full local `cmd/ailang-worldd` plain run at ~696 s under load, so local plain runs
   already pass their own `-timeout`"* (brief) is **false as stated**. The 696 s reading was taken with an explicit
   `-timeout 45m` (`world-iter442/readings.md:17`). `verify_go.sh`'s plain leg passes no `-timeout`, so it runs under
   Go's default 10 m = 600 s, which a 696 s package would exceed [D9]. This item does not change the plain leg (§9).
   I record the plain leg's 600 s exposure under rig load as an open observation, not a fix.
7. The brief counts 52 attempts. The table has 53 rows: 51 complete distinct attempts, the in-progress
   `6420426e` a2, and a duplicate of `45792033970` [D2].

The latest instance, VERIFIED BY CONTROLLER and re-read [D6]: dev `6420426e` (docs-only), run 45936239736 a1:
`FAIL cmd/ailang-worldd 480.085s`, `FAIL host/daemon 480.072s`, `host/verifygate 424.143s`,
`host/broker 333.040s`, `host/store 97.766s` (`k` 1.321, a mid-class runner).

## 2. Mechanism

`go test ./...` builds the test binaries and runs them `-p` at a time. It launches each binary with the one
`-timeout`, and the binary's alarm panics at that bound, naming the running test and dumping goroutines.
`cmd/go` adds a backup kill at `-timeout` + 1 minute (`testKillTimeout = testTimeout + 1*time.Minute` when the
WaitDelay is under a minute, `go1.26.6 src/cmd/go/internal/test/test.go:840-841`) [D9].

A package's `-race` time scales with the runner. On this population the scaling is **not uniform across packages**.
Dividing each uncapped package time by `k` gives `host/daemon` 355.3–407.2 s, `cmd/ailang-worldd` 262.6–361.6 s
and `host/verifygate` 288.7–337.5 s [D3]. So the heavy packages grow faster than the light ones: `cmd/ailang-worldd`
on the slow tail reads 305–341 per unit `k`, against 266 on the normal cluster. A bound sized to the normal
cluster's 357 s, at 1.34× headroom (480/357), therefore fails whenever `k` exceeds ~1.18 for `host/daemon`
(480/407.2), which is why the cap fires on mid runners too, not only on the slow tail.

**What the slow tail actually needs** (the censored packages, extrapolated with each package's *worst* measured
per-`k` time to the worst measured `k` = 1.585, run 45787704457 a1) [D3]:

| package | worst dur/`k` (uncapped obs.) | need at `k` 1.585 |
|---|---|---|
| `host/daemon` | 407.2 (run 45898245166 a2: 468.6 s at `k` 1.151) | **645.4 s** |
| `cmd/ailang-worldd` | 361.6 | 573.1 s |
| `host/verifygate` | 337.5 | 534.9 s |

Every heavy package's need on the slow tail exceeds 480 s; the largest (645.4 s) sits 254.6 s below 900 s.

**The leg's schedule** [D4]: `cmd/ailang-worldd` starts 43–73 s into the leg, `host/daemon` 105–166 s, and
`host/verifygate` (alphabetically late) 322–554 s, once an earlier package frees a slot. The slow-class leg ends at
930–1034 s *with* the three heavy packages cut off at 480 s. Uncapped, the leg grows by at most the sum of their
extensions, which is a serial worst case:
1034 + (645.4 − 480.1) + (573.1 − 480.2) + (534.9 − 480.0) = **1347.1 s**.

## 3. Options considered

**(A) A larger fixed per-package bound, with the leg bound and step ceiling recomputed. CHOSEN, with a
process-naming leg wrapper and a runtime ordering guard.** These measurements decide it:

- The slow tail's worst need is 645.4 s (table above). P = **900 s** (15 m) clears it by 254.6 s (1.39×). It
  tolerates `host/daemon` up to `k` = 900/407.2 = 2.21 at its worst measured per-`k` time, against a worst
  observed `k` of 1.585. On the normal cluster it is 2.50× the slowest daemon reading (360.3 s).
- Nothing in the data argues for a smaller P with more machinery. The normal cluster's leg (643–673 s) is
  unchanged by any P above 480, because no normal-cluster package comes near it.
- Cost: a genuinely hung test is detected later, at its start + 900 s instead of start + 480 s. For `host/daemon`
  (start ≤ 166 s) that is ≤ 1066 s into the leg instead of ≤ 646 s, i.e. +7 minutes on a run that is already red.
  It is still loud, and still names the test and package.

**(B) A per-package bound derived in-run from a measured stimulus. Rejected on measurement.** I take the
verification rule seriously ("where a bound and its stimulus both scale with the machine, the bound must be
derived from the measured stimulus"), and here the bound *is* derived from a measured stimulus: the per-package,
per-`k` table of 51 attempts, extrapolated to the worst observed runner. It is not derived in-run, because no
in-run stimulus is both early enough and faithful enough:
- `host/store`, or any uncapped package, finishes 459–646 s into the leg, after the heavy binaries were launched
  with their `-timeout` fixed [D4]. Using it would mean a separate serial calibration run first, adding 58–102 s to
  every attempt (`host/store`'s own `-race` time) and measuring a different contention state from the one the bound
  must cover.
- The only stimulus available before launch is the pre-leg time (race control + `go build ./...`). It separates
  the slow class (46–55 s against 34–38 s) but **misses the mid class**: attempts at `k` 1.057–1.152 read 34–39 s,
  the same as normal [D3]. A bound derived from it would under-predict exactly where 4 of the 18 failing attempts
  sit. The pre-leg time also depends on the module/build cache state, which this design does not control.
- A derived bound therefore needs a ≥ 1.4× safety factor *and* a floor, which is fixed bound (A) with extra moving
  parts and a new failure mode (a calibration that mis-measures silently changes the hang bound).
- Rule 3m's target is a *test assertion* whose tight ceiling encodes one machine. This bound is a hang detector
  whose only requirements are "finite, loud, above every honest run". The residual exposure is stated, not hidden:
  a runner beyond `k` ≈ 2.2 would hit 900 s, loudly and named (§9).

**(C) Splitting the race job across runners. Rejected for M1, not refuted.** What is measured: each extra shard
pays ~15 s of job setup (job start 03:47:56.9 → verify step 03:48:11.5, run 45936239736 a1) plus the 34–55 s
pre-leg [D3, D6]. The repo is PUBLIC [D7], so standard-runner minutes are not billed, and the cost is runner
concurrency, not money. What is *not* measured, and decides whether C helps: whether `host/daemon` alone on a slow
runner falls below 480 s. The per-`k` table says a slow runner's own slowdown carries the heavy packages to
535–645 s, and sharding does not change the runner. So C cannot be shown to remove the cap without paid
experiments; at best it shortens normal-class wall (verifygate would start at ~0 s, not ~327 s). C also brings two
gate risks that would need their own guards: (i) a package missing from every shard silently drops out of `-race`,
and (ii) `host/verifygate/toolchain_pin_gate_test.go:159` pins the exact job set `{ailang-verify, go-race-verify,
go-verify}` [D8]. It is left as a wall-time follow-up (§9), to be measured as
`PENDING — controller to measure: one manual workflow_dispatch of a 2-shard race job on a draft branch, reading
per-package -race seconds per shard`.

**(D) Moving specific slow tests out of the shared bound. Rejected on measurement.** There is no slow test to move
(§1 correction 5). The interrupted test has run ≤ 25 s, and the cmd panics spread over six tests at 0 s. Moving
tests would also change package boundaries in production-test files for no measured gain.

**Why A needs two small additions to stay honest (constraint 1):**

1. **A runtime ordering guard.** A hung package's own panic (which names it) must fire before the leg wrapper's
   SIGKILL. The latest heavy start is 554 s (`host/verifygate`, run 45787704457 a1) [D4]. Uncapping `host/daemon`
   can delay that slot by at most its 165.3 s extension, giving a start allowance A = 554 + 166 = **720 s**. So
   L ≥ P + A = 1620 s is required. `verify_go.sh` refuses (rc 2, before any leg runs) when that does not hold, so
   an override cannot silently invert the layers.
2. **The leg wrapper names what it kills.** A package that starts after L − P (> 900 s into the leg; only light
   packages, on a slow runner) and then hangs is still caught by L first. Today the wrapper prints only "leg timed
   out". The new wrapper lists the leg's process group (`ps -eo pgid=,etime=,args=`) before the SIGKILL. A test
   binary's argv is `…/b<N>/<pkg>.test`, so the package is named. This also covers hangs outside a test binary
   (a stuck link step, a stuck `go` command).

**The bounds, with derivations (all seconds; measured inputs from §2 and [D3, D4]):**

| layer | value | derivation | margin | catches |
|---|---|---|---|---|
| per-package `go test -timeout` (P) | **900** (15 m) | ≥ 1.25 × slow-tail need 645.4 = 806.8; rounded up to 15 m | 254.6 s over need (1.39×); daemon tolerates `k` ≤ 2.21 | a test that hangs or runs away inside a test binary: `panic: test timed out after 15m0s` + test name + goroutines + `FAIL <pkg>` |
| `cmd/go` backup kill | P + 60 = 960 | Go's own (`test.go:840-841`) [D9] | — | a test binary whose alarm does not fire: `*** Test killed: ran too long (16m0s)` + `FAIL <pkg>` |
| start allowance (A) | **720** | latest heavy start 554 + daemon's max uncapped extension 165.3 → 720 | — | (ordering input only) |
| leg bound (L) | **1800** (30 m) | ≥ P + A = 1620 (ordering); ≥ slow-leg serial worst case 1347.1 | 180 s over P + A; 452.9 s over 1347.1 (1.34×) | anything still running at 30 m: a package started after 900 s that hangs, a stuck build/link or `go` process; names the running processes, SIGKILLs the group, rc 124 |
| step `timeout-minutes` (J) | **40** (2400 s) | ≥ L + 300; 300 ≥ 5 × the 55 s max pre-leg time | 600 s over L (545 s after the measured 55 s pre-leg) | a hang before or around the leg (race-detector control, `go build ./...`, the wrapper itself); Actions fails the step |

Ordering: 900 < 960 < 1620 ≤ 1800 < 2100 ≤ 2400. Each inequality is asserted by a test (§5).

Expected effect, from the table: all 18 attempts that panicked at 480 s fall below 900 s on the extrapolation (max
need 645.4 s), so none would have needed a re-run for this class. The 5 attempts with real test failures at 0
timeouts (`c1f848de`, `18743820`, `81789324`, and the cmd failures beside a daemon timeout in `3cd68ecc`, `255413ef`)
are a different class and stay red [D3]. The worst slow-class job wall becomes ≈ 15 + 55 + 1347 ≈ 1417 s (~24 min,
the serial upper bound). Today the same attempt fails at ~18 min and then pays a ~12 min re-run.

## 4. M1 — the change (one milestone)

**`scripts/verify_go.sh`** (~+45/−14):

1. Top-level constants and a validator, beside `check_mission_config` (before the subcommand dispatch, so
   subcommands can use them):
   ```bash
   # Row 267 (design_docs/planned/w-daemon-race-package-timeout-headroom.md §3): per-package -race bound P,
   # derived from the measured slow-runner need (645 s at k=1.585) with 1.39x margin; start allowance A =
   # latest measured heavy-package start (554 s) + the daemon's uncapped extension (166 s).
   VERIFY_GO_RACE_PKG_TIMEOUT_S=900
   VERIFY_GO_RACE_START_ALLOWANCE_S=720
   VERIFY_GO_RACE_TIMEOUT_S="${VERIFY_GO_RACE_TIMEOUT_S:-1800}"
   check_race_bounds() {
     case "$VERIFY_GO_RACE_TIMEOUT_S" in
       ''|*[!0-9]*) echo "verify_go.sh: FATAL: VERIFY_GO_RACE_TIMEOUT_S must be a positive integer" >&2; return 2 ;;
     esac
     min=$((VERIFY_GO_RACE_PKG_TIMEOUT_S + VERIFY_GO_RACE_START_ALLOWANCE_S))
     if [ "$VERIFY_GO_RACE_TIMEOUT_S" -lt "$min" ]; then
       echo "verify_go.sh: FATAL: leg bound ${VERIFY_GO_RACE_TIMEOUT_S}s must be at least per-package ${VERIFY_GO_RACE_PKG_TIMEOUT_S}s + start allowance ${VERIFY_GO_RACE_START_ALLOWANCE_S}s (${min}s)" >&2
       return 2
     fi
   }
   ```
   (`min` must not shadow anything; the implementer checks with `grep -n '\bmin=' scripts/verify_go.sh`.)
2. `run_race_leg BOUND CMD...`, a shell function wrapping the existing heredoc. Python reads the bound from
   `argv[1]` and the command from `argv[2:]`, via `python3 - "$@" <<'PY'`. It keeps
   `start_new_session=True`, `p.wait(timeout=bound)` with the exit code passed through, `os.killpg(…, SIGKILL)`,
   `print(f"verify_go.sh: FATAL: -race leg timed out after {bound}s", file=sys.stderr)` and `sys.exit(124)`. NEW,
   before the kill: list every process whose pgid is the leg's, each as
   `verify_go.sh: still running at the leg bound (elapsed <etime>): <args>`. If `ps` itself fails, print one
   `could not list` line and still kill and exit 124 (diagnostics never gate the refusal).
3. Two subcommands beside `--evidence-manifest-check` and `--mission-config-check`, both exiting before the
   `AILANG_BIN` check:
   - `--race-bounds-check` (no arguments): runs `check_race_bounds`; on success prints exactly
     `race bounds: per-package ${P}s, start allowance ${A}s, leg ${L}s` and exits 0.
   - `--race-leg-wrapper-check BOUND CMD...`: `run_race_leg "$@"`, exits with its rc.
4. Main flow: replace lines 222–226 with `check_race_bounds || exit $?` and the banner
   `echo "── legs: $VERIFY_GO_LEGS (race bounds: per-package ${VERIFY_GO_RACE_PKG_TIMEOUT_S}s, leg ${VERIFY_GO_RACE_TIMEOUT_S}s)"`.
   (The `VERIFY_GO_LEGS` case stays where it is.)
5. Race leg block (`:348–364`): first line `bash scripts/test_race_leg_bounds.sh` (the wrapper's own adversarial
   arms run before the wrapper is trusted, the same precedent as `test_tracked_binary_gate.sh` at `:239`). Then the
   banner `echo "── go test ./... -count=1 -race -timeout ${VERIFY_GO_RACE_PKG_TIMEOUT_S}s (leg bound ${VERIFY_GO_RACE_TIMEOUT_S}s)"`
   and the single call line
   `run_race_leg "$VERIFY_GO_RACE_TIMEOUT_S" go test ./... -count=1 -race -timeout "${VERIFY_GO_RACE_PKG_TIMEOUT_S}s"`.
   Update the stale comment block above it: "2-core runner" is unverified (nproc not in the logs, §9), and the
   new comment points at this doc.
6. The edit must add no `v0.41.0` token (`host/pins/scan_test.go:158` counts exactly 4 in this file) [D8], and
   must not touch the P1/P2 toolchain lines read by `toolchain_pin_gate_test.go` [D8].

**`scripts/test_race_leg_bounds.sh`** (new, ~75; `set -euo pipefail`, scratch under `mktemp -d "${TMPDIR:-/tmp}/…"`
with an EXIT trap; prints `ok <arm>` per arm and `PASS 6/6 arms`; any failure prints the arm and its captured
output and exits 1):
- **R1 default**: `env -u VERIFY_GO_RACE_TIMEOUT_S bash scripts/verify_go.sh --race-bounds-check` → rc 0, stdout
  exactly `race bounds: per-package 900s, start allowance 720s, leg 1800s`.
- **R2 boundary accept**: `VERIFY_GO_RACE_TIMEOUT_S=1620 … --race-bounds-check` → rc 0.
- **R3 boundary refuse**: `VERIFY_GO_RACE_TIMEOUT_S=1619` → rc 2, stderr contains `must be at least per-package 900s`.
- **R4 malformed**: `VERIFY_GO_RACE_TIMEOUT_S=15m` → rc 2, `must be a positive integer`.
- **R5 passthrough**: `--race-leg-wrapper-check 30 sh -c 'exit 3'` → rc 3, and stderr has no `FATAL`.
- **R6 hang**: a fake executable `$tmp/hang_probe.test` (`#!/bin/sh`: `sleep 300 & echo $! >"$tmp/gc.pid"; sleep 300`)
  run as `--race-leg-wrapper-check 2 "$tmp/hang_probe.test"` → rc 124. Stderr contains `still running at the
  leg bound`, `hang_probe.test`, and `FATAL: -race leg timed out after 2s`. Elapsed ≥ 2 s and < 30 s. Known-positive
  control: `gc.pid` is non-empty. After return, `kill -0 "$(cat gc.pid)"` fails (the grandchild died with the
  group). Wall ~2–3 s.

**`host/verifygate/race_job_gate_test.go`** (~+95/−1):
- Existing `TestRaceLegRunsInItsOwnJob`: replace the needle `"go", "test", "./...", "-count=1", "-race"` with
  the new call line. Keep `FATAL: -race leg timed out after {bound}s` and the race-control needle.
- NEW `TestRaceBoundsAreOrdered` (pure file parse, milliseconds; it never needs `AILANG_BIN`). It reads
  `verify_go.sh` and `ci.yml`, and requires each parse to match **exactly once**, which is its non-vacuity control:
  - from `verify_go.sh`: `^VERIFY_GO_RACE_PKG_TIMEOUT_S=(\d+)$` → P;
    `^VERIFY_GO_RACE_START_ALLOWANCE_S=(\d+)$` → A; `^VERIFY_GO_RACE_TIMEOUT_S="\$\{VERIFY_GO_RACE_TIMEOUT_S:-(\d+)\}"$`
    → L_default; the call line, exactly, as a whole line;
  - from the `go-race-verify` job body (the same `job()` splitter the existing test uses): the step containing
    `run: ./scripts/verify_go.sh` → `timeout-minutes: (\d+)` → J and `VERIFY_GO_RACE_TIMEOUT_S: "(\d+)"` →
    L_ci. Control: the job contains `runs-on: ubuntu-24.04`, and has zero `continue-on-error`.
  - Measured constants, named with their sources (this doc §2/§3): `slowNeedS = 646`, `latestHeavyStartS = 554`,
    `daemonExtensionS = 166`, `preLegMaxS = 55`.
  - Asserts: `P*4 >= slowNeedS*5` (P ≥ 1.25 × need); `A >= latestHeavyStartS + daemonExtensionS`;
    `L_ci == L_default`; `L_default >= P + A`; `J*60 >= L_ci + 300` and `300 >= 5*preLegMaxS`; the call line
    contains `-race` and contains none of `-short`, `-run`, `-skip`, `-failfast=false`. Each failure message names
    the layer and the inequality.

**`host/verifygate/mission_config_gate_test.go`** (~+12/−4): in `requireVerifyGoProductSequence`, the
`repaired` arm takes the new banner needle and the new call-line needle in place of
`cmd = ["go", "test", …, "8m"]`. It drops the two Python needles (`print(f"… timed out after {bound}s"` and
`    sys.exit(124)`) from the *ordered* sequence, because they now live in `run_race_leg`, which is defined above the
`AILANG_BIN` anchor, and asserts instead that each occurs exactly once anywhere. The embedded `565d0b2` fixture
(`:415-442`) is history and stays byte-identical.

**`.github/workflows/ci.yml`** (~+5/−3): step "go -race gate": `timeout-minutes: 40`,
`VERIFY_GO_RACE_TIMEOUT_S: "1800"`. The comment names the three layers and points at this doc. No new job, no new
step (the wrapper test runs inside `verify_go.sh`).

## 5. Mutation table

Each mutation is a single edit to the M1 tree. "Kills" names the test and the assertion that fires.

| id | mutation | killed by |
|---|---|---|
| M1 | revert P: `VERIFY_GO_RACE_PKG_TIMEOUT_S=480` | `TestRaceBoundsAreOrdered` (P ≥ 1.25 × 646); `test_race_leg_bounds.sh` R1 (stdout text) |
| M2 | call line hardcodes `-timeout 8m` (P defined but unused) | `TestRaceBoundsAreOrdered` (exact call line); `TestVerifyGoProductCheckSequence` repaired arm (call-line needle count 0) |
| M3 | CI keeps `VERIFY_GO_RACE_TIMEOUT_S: "1500"` | `TestRaceBoundsAreOrdered` (`L_ci == L_default`, and `L ≥ P + A`) |
| M4 | script default back to `:-1500` | `TestRaceBoundsAreOrdered` (`L_default ≥ P + A`); R1 (rc 2, not 0) |
| M5 | `timeout-minutes: 30` (J×60 = 1800 ≤ L) | `TestRaceBoundsAreOrdered` (`J*60 ≥ L + 300`) |
| M6 | `check_race_bounds` ordering branch deleted (or `return 0` first) | R3 (rc 0, want 2) |
| M7 | ordering comparison `-lt` → `-le` (off by one) | R2 (rc 2 at 1620, want 0) |
| M8 | `VERIFY_GO_RACE_START_ALLOWANCE_S=0` | `TestRaceBoundsAreOrdered` (A ≥ 720); R1 text |
| M9 | wrapper's process listing deleted | R6 (`still running at the leg bound` / `hang_probe.test` absent) |
| M10 | `os.killpg(...)` → `p.kill()` (child only) | R6 (grandchild `kill -0` succeeds) |
| M11 | `sys.exit(p.wait(...))` → `p.wait(...); sys.exit(0)` | R5 (rc 0, want 3) |
| M12 | timeout path `sys.exit(124)` → `sys.exit(0)` | R6 (rc 0, want 124) |
| M13 | `-short` appended to the call line | `TestRaceBoundsAreOrdered` (forbidden flag, and the exact line) |
| M14 | `-race` removed from the call line | `TestRaceBoundsAreOrdered`; existing `TestRaceLegRunsInItsOwnJob` |
| M15 | `continue-on-error: true` on the race job/step | `TestRaceBoundsAreOrdered` (zero `continue-on-error`) |
| M16 | `bash scripts/test_race_leg_bounds.sh` line dropped from the race block | `TestVerifyGoProductCheckSequence` repaired arm, *only if* the implementer adds it as a needle; **required**: add it |
| M17 | `check_race_bounds || exit $?` dropped from the main flow | **required new needle** in the repaired sequence, the same as M16 |

M16/M17 are listed so the implementer adds the two needles; without them those edits would survive. There are no
declared equivalents.

## 6. Acceptance criteria

All commands run from the worktree root on the rig with `AILANG_BIN=$HOME/.pinned-ailang/ailang` (v0.41.0,
released), `GOTOOLCHAIN=go1.26.6`.

- **AC1** `bash scripts/verify_go.sh --race-bounds-check; echo rc=$?` →
  `race bounds: per-package 900s, start allowance 720s, leg 1800s` / `rc=0`.
- **AC2** `VERIFY_GO_RACE_TIMEOUT_S=1619 bash scripts/verify_go.sh --race-bounds-check; echo rc=$?` → stderr
  `… leg bound 1619s must be at least per-package 900s + start allowance 720s (1620s)` / `rc=2`.
- **AC3** `bash scripts/test_race_leg_bounds.sh; echo rc=$?` → six `ok R<n> …` lines, `PASS 6/6 arms`, `rc=0`.
  It also passes with `TMPDIR` under a 1777 parent and under CPU contention, per the iteration-442 Gate 5 process
  rule: `TMPDIR=$(mktemp -d /private/tmp/r267.XXXX) bash scripts/test_race_leg_bounds.sh` with 8 `yes >/dev/null`
  spinners running → `PASS 6/6 arms`.
- **AC4** `go vet ./host/verifygate/ && go test ./host/verifygate/ -count=1 -v -run 'TestRaceBoundsAreOrdered|TestRaceLegRunsInItsOwnJob|TestVerifyGoProductCheckSequence|TestVerifyGoRefusesUnknownLegs' 2>&1 | grep -E '^(--- |ok|FAIL)'`
  → four `--- PASS` lines and `ok`. (`go vet` is the compile fence for `_test.go`; `go build` is not.)
- **AC5** Mutation replay: for each of M1–M17, apply the single edit, run the killing test(s) named in §5, and
  observe red with the named assertion; restore and observe green. Record a table `id | killing test | rc mutated |
  rc restored`. Expected: 17/17 killed, 0 survivors.
- **AC6** Whole gates on the M1 tree: `./scripts/verify_ail.sh` rc 0; `go build ./... && go test ./... -count=1`
  with the fail set compared to base `6420426e`, so no new failure; `gofmt -l host/verifygate/` empty;
  `bash -n scripts/verify_go.sh scripts/test_race_leg_bounds.sh` rc 0.
- **AC7** CI on the PR: the race job's log contains `── legs: race (race bounds: per-package 900s, leg 1800s)`,
  the shell test's `PASS 6/6 arms`, and
  `── go test ./... -count=1 -race -timeout 900s (leg bound 1800s)`, and all three jobs are green.
- **AC8 (observational; recorded, never claimed if unobserved)**: over the first 10 race-job attempts after merge,
  any attempt in which a heavy package reads > 480 s is green for that package, and no
  `panic: test timed out after 15m0s` appears. If none of the 10 draws a non-normal runner (`host/store` ≤ 67.5 s
  in all), AC8 is recorded as **unobserved**, not passed.

## 7. Conflict surface (who else reads these bounds or these lines)

Greps in [D8]:

- `host/verifygate/mission_config_gate_test.go:436-437, 500-501`: exact-text needles for the race banner and
  `cmd = [... "8m"]`. The fixture lines (`:436-437`) stay; the shared needles (`:500-501`) become `repaired`-
  parameterised (§4). **Must change in M1 or that test goes red.**
- `host/verifygate/race_job_gate_test.go:37, 55-57`: `VERIFY_GO_RACE_TIMEOUT_S:` presence and the Python `cmd`
  substring. The latter must change (§4).
- `host/pins/scan_test.go:158`: counts `v0.41.0` tokens in `verify_go.sh` (n = 4). M1 adds none.
- `host/verifygate/toolchain_pin_gate_test.go:123, 159-160, 1189-1340`: reads `verify_go.sh`'s P1/P2 toolchain
  bindings and pins the exact CI job set. M1 touches neither, and is a reason option C would need more edits.
- `host/archive/archive.go:66`: a comment citing "(-timeout 8m race under a 600 s outer wait)". It is **already
  stale** (the outer wait has been 1500 s since row 181) and stays true in substance (10 s ≪ 900 s). It is Go
  production code and is not edited by M1. A follow-up comment fix is listed in §9.
- No test reads any `timeout-minutes` [D8]. `tools/launchd/*` does not reference `verify_go` or
  `go-race-verify` [D8]. No `design_docs/verification/*/gates.sh` sets `VERIFY_GO_RACE_TIMEOUT_S` (only
  historical CI-log copies mention it) [D8].
- Row 243 V4 S6 (the next product sprint) adds tests to `cmd/ailang-worldd`. With P = 900 its directive's
  "keep the package's race time flat or move the bound" clause relaxes to the §3 margin: `cmd/ailang-worldd`'s
  slow-tail need is 573.1 s, so 326.9 s of headroom.

## 8. Sizing

One milestone. ~+232/−22 ≈ 254 changed lines: `verify_go.sh` ~59, `test_race_leg_bounds.sh` ~75,
`race_job_gate_test.go` ~96, `mission_config_gate_test.go` ~16, `ci.yml` ~8. Within the < 300 brief. Shell and
Go test only, no `.ail`, no production Go. CI cost: +~3 s for the shell arms in each race-leg run, and +ms for the
static test in both Go jobs. `host/verifygate`'s race time (a heavy package) grows by the static test only, which
is file reads.

## 9. What this item is NOT doing

- Not dropping any package from `-race`, not adding `-short`/`-skip`/`-run`, no `continue-on-error`, no retry
  loop. The race-detector known-positive control and the leg wrapper's loud rc 124 stay.
- Not sharding the race job (option C). That is a wall-time follow-up gated on the PENDING measurement in §3 C.
- Not moving or editing any test in `host/daemon`, `cmd/ailang-worldd` or `host/verifygate` (option D). There is no
  measured culprit.
- Not changing the plain leg. Its implicit Go default 600 s per-package bound under rig load (the 696 s reading in
  §1 correction 6) is recorded as an open observation for a separate row, if the controller confirms it bites.
- Not fixing `host/archive/archive.go:66`'s stale comment (production file; follow-up).
- Not claiming the runner's vCPU count: `nproc` is not in the logs. The existing "2-core runner" comment is
  unverified, and M1's comment says so. `PENDING — controller to measure: add `nproc; go env GOMAXPROCS` to the race
  job's "go version" step on any PR and read it.`
- Residual exposure, stated: a runner slower than `k` ≈ 2.2 (39 % beyond the worst observed) would reach 900 s in
  `host/daemon`. It fails loudly and named, the same as today's failure, at a 1.39× higher bar.

## 10. Verification Log

Scripts: `~/.ailang/state/world-iter443/design-scratch/`. Logs: `~/.ailang/state/world-iter443/racelogs/`.

| id | claim | command | observed (trimmed) |
|---|---|---|---|
| D1 | the current bounds and their lines | `grep -n 'timeout", "8m"\|VERIFY_GO_RACE_TIMEOUT_S="\${\|-race -timeout 8m\|os.killpg\|sys.exit(124)' scripts/verify_go.sh; grep -n 'timeout-minutes: 35\|VERIFY_GO_RACE_TIMEOUT_S: "1500"\|continue-on-error' .github/workflows/ci.yml; grep -c 'timeout-minutes' .github/workflows/ci.yml` | `222:VERIFY_GO_RACE_TIMEOUT_S="${VERIFY_GO_RACE_TIMEOUT_S:-1500}"`, `350:echo "── go test ./... -count=1 -race -timeout 8m …"`, `355:cmd = ["go", "test", "./...", "-count=1", "-race", "-timeout", "8m"]`, `360: os.killpg(…)`, `362: sys.exit(124)`; ci `522: timeout-minutes: 35`, `525: VERIFY_GO_RACE_TIMEOUT_S: "1500"`; **zero `continue-on-error`**, with positive control `19` `timeout-minutes` lines in the same file |
| D2 | 53 rows = 51 distinct complete + 1 in-progress + 1 duplicate | `tail -n +2 race_table.txt \| grep -c .; wc -l < racelogs/jobs.txt; cmp racelogs/job_45792033970_a1.log racelogs/job_45792033970_a2.log && echo IDENTICAL; wc -c racelogs/job_45936239736_a2.log` | `53`, `53`, `IDENTICAL`, a2 log 215 bytes (in progress at fetch) |
| D3 | class table, per-package medians, dur/`k` ratios, `k`max 1.585 | `python3 -I design-scratch/classes.py` (output `classes.out`) | `attempts (deduped, complete): 51`; `normal 24 k 0.986-1.010 pre 34-38 leg 643-673 store 58.2-67.5 timeouts>0: 0`; `mid 13 k 1.057-1.321 pre 34-43 leg 676-875 store 74.7-97.8 timeouts>0: 4`; `slow 14 k 1.403-1.585 pre 46-55 leg 930-1034 store 82.4-102.2 timeouts>0: 14`; `host/daemon normal median 356.8 … dur/k … max 407.2`; `cmd/ailang-worldd normal median 266.3 … max 361.6`; `host/verifygate normal median 328.6 … max 337.5`; `kmax 1.585` (45787704457 a1) |
| D4 | leg schedule: start offsets, store completion | `python3 -I design-scratch/legs.py racelogs` and the inline start-offset script (print offset − duration; valid because `go test` prints in package order and each heavy package is printed at its own completion) | `cmd/ailang-worldd start offset min 43 max 73`; `host/daemon … min 105 max 166`; `host/verifygate … min 322 max 554`; `host/store` printed at 459–646 s into the leg (`store(dur@end)` column, e.g. `98@620`, `63@459`) |
| D5 | which tests the cap interrupts | `python3 -I design-scratch/atcap.py` (output `atcap.out`), plus `grep -o 'TestWorkbenchTimelinePaging([0-9]*s)' atcap.out \| sort \| uniq -c` | `29 timeout panics`; `15 ('host/daemon','TestWorkbenchTimelinePaging')`, `2 … TestEpisodePolicyProtectedRebased`, `1 … TestWorkspacePackageCacheCoversLock`; six `cmd/ailang-worldd` tests, one each, all `(0s)`; `host/verifygate` 5 (subtests only); TimelinePaging elapsed 1–25 s; `grep -c 'panic: test timed out after 8m0s' racelogs/job_*.log \| awk -F: '$2>0' \| wc -l` → `18` |
| D6 | latest instance (VERIFIED BY CONTROLLER, re-read) | `sed 's/\x1b\[[0-9;]*m//g' racelogs/job_45936239736_a1.log \| grep -E 'FAIL\|ok .*(verifygate\|broker\|store)\|── go test\|Run ./scripts/verify_go.sh\|Run actions/checkout'` | `03:47:57.75 Run actions/checkout@v4`, `03:48:11.46 Run ./scripts/verify_go.sh`, `03:48:53.43 ── go test … -timeout 8m (leg bound 1500s)`, `FAIL …/cmd/ailang-worldd 480.085s`, `ok …/host/broker 333.040s`, `FAIL …/host/daemon 480.072s`, `ok …/host/store 97.766s`, `04:03:28.50 ok …/host/verifygate 424.143s` |
| D7 | repo visibility (minute cost of C) | `gh repo view sunholo-data/ailang-world --json visibility -q .visibility` | `PUBLIC` |
| D8 | readers of these lines/bounds | `grep -rn 'VERIFY_GO_RACE_TIMEOUT_S' --include='*.go' --include='*.sh' --include='*.yml' --include='*.py' . \| grep -v '^./design_docs/'`; `grep -rln 'timeout-minutes' --include='*.go' --include='*.sh' --include='*.py' . \| grep -v design_docs` with control `grep -rln 'ci.yml' --include='*_test.go' host`; `grep -rn 'verify_go\|go-race-verify' tools/launchd` with control `ls tools/launchd \| wc -l`; `grep -n 'verify_go' host/pins/scan_test.go`; `grep -rln VERIFY_GO_RACE_TIMEOUT_S design_docs/verification` with control `grep -rln VERIFY_GO_LEGS design_docs/verification` | `race_job_gate_test.go:37`, `mission_config_gate_test.go:436,500`, `verify_go.sh:222-226,350-354`, `ci.yml:521,525`; **timeout-minutes readers: none**, with control listing 9 `_test.go` files that read `ci.yml`; **tools/launchd: no match**, with control `5` files present; `scan_test.go:158 "scripts/verify_go.sh": {… n: 4 …}`; verification hits only historical CI copies (`world-iter386/ci_merge_race_attempt1.log`, `world-iter329/plan-M*-ci.yml.txt`), with control `world-iter397/gates.sh` matching `VERIFY_GO_LEGS` |
| D9 | Go's timeout semantics; the brief's 696 s premise | `go help testflag \| grep -A3 -- '-timeout d'`; `grep -n 'testKillTimeout\|WaitDelay' $(go env GOROOT)/src/cmd/go/internal/test/test.go`; `grep -n 696 design_docs/verification/world-iter442/readings.md` | `The default is 10 minutes (10m).`; `840: if testWaitDelay < 1*time.Minute {` / `841: testKillTimeout = testTimeout + 1*time.Minute`; `1740: "*** Test killed with %v: ran too long (%v).\n"`; readings `17: base2/cmd-plain (-timeout 45m): complete, 696 s …` |
| D10 | existing pins of `verify_go.sh` text | `grep -rln 'verify_go' --include='*.go' --include='*.sh' --include='*.yml' --include='*.py' . \| grep -v '^./design_docs'` | `host/verifygate/{mission_config_gate,evidence_manifest_gate,ail_binary_gate,toolchain_pin_gate,race_job_gate}_test.go`, `host/pins/scan_test.go`, `host/{broker,transitionreg,runbook,replay}` (prose/comments), `scripts/{verify_go,tracked_binary_gate,test_tracked_binary_gate}.sh`, `scripts/floor/canary.py`, `.github/workflows/ci.yml` |

Inherited, VERIFIED BY CONTROLLER (iteration 443): `race_table.txt` and `racelogs/` (53 attempt logs, run ids,
shas, events, conclusions); the latest instance of §1. Each number used above was re-derived from the raw logs
(D2–D6), not from the table's text.

Not measured (each marked PENDING in the body): the runner vCPU count; per-shard `-race` times (option C); a full
local `VERIFY_GO_LEGS=all` run on the M1 tree (AC6 covers it at implementation).

## 11. Quorum verification log

Pending.
