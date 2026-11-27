# w-daemon-race-package-timeout-headroom — the 8-minute cap encodes the fast runner; one in three race attempts draws a slower one

- **Status**: Planned (r1, iteration 443, unattended designer `claude:claude-opus-5-5`; r0 committed as `8256a83e`,
  quorum round 1 BLOCKED 3/3 on completeness/derivation, §11)
- **Item**: charter queue row 267, `w-daemon-race-package-timeout-headroom`, clause-2 hygiene. Picked under the
  standing red-dev rule (D-WORLD-121 OPEN, default B: Gate 1 found `dev`'s tip `6420426e` red at attempt 1)
- **Estimated**: one milestone, ~0.5 day, ~262 changed lines (§8)
- **Measurement base**: `6420426e` (dev tip, clean worktree `sprint/world-iter443-row267-race-headroom`), 2026-11-27
- **Instruments**: CI job logs of the `go-race-verify` job (ubuntu-24.04, Azure eastus2 hosted runners) for the
  last 40 workflow runs, 53 attempt files under `~/.ailang/state/world-iter443/racelogs/` (fetched by the
  controller); the controller's `race_table.py`/`race_table.txt`; this session's re-derivation scripts
  `~/.ailang/state/world-iter443/design-scratch/{legs,classes,atcap,slots}.py` (run with `python3 -I`; outputs
  `classes.out`, `atcap.out`, `slots.out`); controller measurements C1–C4 (iteration 443, at `8256a83e`); local `go1.26.6 darwin/arm64` source of `cmd/go` for the `go test` kill path.
  No `.ail` is touched, so the pinned `ailang` matters only to the gate runs (§6).
- **Files touched**: `scripts/verify_go.sh` (~+47/−14), `scripts/test_race_leg_bounds.sh` (new, ~80),
  `host/verifygate/race_job_gate_test.go` (~+95/−1), `host/verifygate/mission_config_gate_test.go` (~+12/−4),
  `.github/workflows/ci.yml` (~+5/−3). No Go production code, no `.ail`, nothing under `tools/launchd/`.

Quantities carry their instrument. `VERIFIED BY CONTROLLER` marks inherited controller facts; `D<n>` marks this
session's re-derivations, each with its command and output in §10.


## 0. r2 amendments (quorum round 2, applied by the controller under the narrow-refinement carve-out) — these SUPERSEDE every conflicting value in §1–§3

Round 2 (seats oc-glm-5-3 PASS, oc-kimi-k3 REJECT, codex:gpt-6.2-astra REJECT; no absent seat) left two blocking objections,
each with a concrete reviewer-authored `proposed_fix` and neither disputing the direction (option A, P = 900 s, strict layer
ordering). The controller applied those fixes verbatim (iteration 443, Gate 2 narrow-refinement carve-out); §11 logs them.

- **kimi (censored daemon slow tail).** `host/daemon` is capped in every slow-tail attempt, so its slow-tail per-`k` time is
  unmeasured. Two models: (i) class-invariant per-`k` → need 645.4 s (§3, r1); (ii) `cmd/ailang-worldd`'s measured inflation on the
  slow tail (+28%, 266 → 341 s per unit `k`) applied to daemon → need 407.2 × 1.28 × 1.585 ≈ **827.4 s**, uncapped extension
  827.4 − 480.1 ≈ **347 s**. The conservative stack is therefore A = 554 + 92.9 + 347 ≈ 993.9 → **A = 995 s**; the guard is
  L ≥ P + 60 + A = 900 + 60 + 995 = **1955 s**; **L = 2000 s** (default and CI); J = 40 still holds (2400 ≥ 2000 + 300 = 2300).
  P = 900 s is kept: it clears model (i) by 254.6 s (`k` ≤ 2.21) and model (ii) by ~73 s (`k` tolerance ≈ 1.73). The r1 margin
  claims in §3 ("254.6 s (1.39×)", "`k` ≤ 2.21") hold under model (i) only; under model (ii) they are ~73 s and ≈ 1.73.
  Ordering: 900 < 960 < 1955 ≤ 2000 < 2300 ≤ 2400. §4–§6 carry the r2 values (R1 text, R2 1955 accept, R3 1954 refuse,
  AC1/AC2, M4 `:-1900`, M8 = 813, the test's `daemonExtensionInflatedS = 347` and `P > daemonInflatedNeedS`).
  PENDING — planner to measure: daemon per-`k` versus `k` across the uncapped mid-class attempts (`k` 1.057–1.321) with
  `~/.ailang/state/world-iter443/design-scratch/classes.py`; if it is flat or falling, say so — the r2 values stand either way
  (they are the conservative reading).
- **astra (unbounded diagnostics on the timeout path).** `ps` gets a 2-second timeout and the process-group kill sits in a
  `finally`, so the listing can delay the SIGKILL by at most 2 s and never prevent it; new arm R7 (sleeping `ps` shim, 10 s
  watchdog) and mutations M20/M21 (§4, §5).
- **kimi catches.** AC6 gains an explicit `-timeout 45m` with a same-session base comparison; the drop-`-ww` mutant is declared
  unobserved everywhere rather than "observed under AC7".
- **glm (PASS) notes, applied.** AC8 classifies failure causes (§6, glm's text). §7's quoted row-243 clause is unverified
  in this doc's log — treat it as a paraphrase (PENDING — planner: `grep -n 'keep the package' design_docs/world-mission.md`).
  The `03:47:56.9` job-start timestamp in §3 C is unlogged and non-load-bearing (option C is deferred); disregard it.

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
   `-timeout 45m` (`world-iter442/readings.md:17`). `verify_go.sh`'s plain leg passes no `-timeout`: the script has exactly two
   `-timeout` lines, both in the race leg (`:350` banner, `:355` command) [C2, D12]. So the plain leg runs under
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
930–1034 s *with* the three heavy packages cut off at 480 s. **Which slot does `host/verifygate` take?** [D11]: in all 51 attempts it started *before* `host/daemon`
ended (by 91.4–172.7 s in the 15 attempts where both `cmd/ailang-worldd` and `host/daemon` read ≥ 440 s). In
those 15 it started 83.0 s before to 1.3 s after `cmd/ailang-worldd` ended. So on the measured schedules it takes
`cmd/ailang-worldd`'s slot, or another early package's, and never waits on `host/daemon`'s. r0's "delayed by
daemon's extension" premise is refuted by the data; the slot it does track is cmd's. Uncapping changes the schedule,
though, and the data cannot rule out a cascade, so the start allowance takes the stacked worst case (§3).
Uncapped, the leg grows by at most the sum of the extensions, which is a serial worst case:
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

1. **A runtime ordering guard.** A hung package's own panic (which names it), and `cmd/go`'s backup kill 60 s
   after it [D9], must both fire before the leg wrapper's SIGKILL. The latest heavy start is 554 s
   (`host/verifygate`, run 45787704457 a1) [D4]. Under the serial worst-case model already used for the leg (§2),
   that start can be delayed by every extension ahead of its slot: `cmd/ailang-worldd`'s 92.9 s plus
   `host/daemon`'s 165.3 s = 258.2 s. That gives a start allowance A = 554 + 258.2 = 812.2 → **813 s**. The
   measured schedules support only the cmd term (§2, [D11]: 554 + 93 = 647). The stacked value is kept because it
   costs 100 s of leg bound and nothing on a passing run. So L ≥ P + 60 + A = 900 + 60 + 813 = **1773 s** is
   required. `verify_go.sh` refuses (rc 2, before any leg runs) when that does not hold, so
   an override cannot silently invert the layers.
2. **The leg wrapper names what it kills.** A package that starts after L − P − 60 (> 940 s into the leg; only light
   packages, on a slow runner) and then hangs is still caught by L first. Today the wrapper prints only "leg timed
   out". The new wrapper lists the leg's process group (`ps -ww -eo pgid=,etime=,args=`) before the SIGKILL. `-ww` (unlimited width)
   is accepted by both BSD `ps` and procps. On darwin under a pipe, the output is identical with and without it: no
   truncation of a 146-character path [C3]. On Linux procps under a pipe it is unmeasured locally, so R6 running in
   the ubuntu-24.04 race job (AC7) is the Linux observation. A test
   binary's argv is `…/b<N>/<pkg>.test`, so the package is named. This also covers hangs outside a test binary
   (a stuck link step, a stuck `go` command).

**The bounds, with derivations (all seconds; measured inputs from §2 and [D3, D4]):**

| layer | value | derivation | margin | catches |
|---|---|---|---|---|
| per-package `go test -timeout` (P) | **900** (15 m) | ≥ 1.25 × slow-tail need 645.4 = 806.8; rounded up to 15 m | 254.6 s over need (1.39×); daemon tolerates `k` ≤ 2.21 | a test that hangs or runs away inside a test binary: `panic: test timed out after 15m0s` + test name + goroutines + `FAIL <pkg>` |
| `cmd/go` backup kill | P + 60 = 960 | Go's own (`test.go:840-841`) [D9] | — | a test binary whose alarm does not fire: `*** Test killed: ran too long (16m0s)` + `FAIL <pkg>` |
| start allowance (A) | **813** | latest heavy start 554 + cmd's 92.9 + daemon's 165.3 uncapped extensions (stacked) = 812.2 → 813 | 166 s over the measured-schedule value 647 [D11] | (ordering input only) |
| leg bound (L) | **1900** (31 m 40 s) | ≥ P + 60 + A = 1773 (ordering, incl. Go's backup kill); ≥ slow-leg serial worst case 1347.1 | 127 s over 1773; 552.9 s over 1347.1 (1.41×) | anything still running at 1900 s: a package started after 900 s that hangs, a stuck build/link or `go` process; names the running processes, SIGKILLs the group, rc 124 |
| step `timeout-minutes` (J) | **40** (2400 s) | J × 60 ≥ L + 300 = 2200 → J ≥ 36.7 → 40; 300 ≥ 5 × the 55 s max pre-leg time | 200 s over L + 300; 500 s over L (445 s after the measured 55 s pre-leg) | a hang before or around the leg (race-detector control, `go build ./...`, the wrapper itself); Actions fails the step |

Ordering: 900 < 960 < 1773 ≤ 1900 < 2200 ≤ 2400. Each inequality is asserted by a test (§5).

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
   # latest measured heavy-package start (554 s) + cmd's (93 s) + daemon's (166 s) uncapped extensions, stacked.
   # The guard adds cmd/go's backup kill (testKillTimeout = -timeout + 1m) so both fire before the leg SIGKILL.
   VERIFY_GO_RACE_PKG_TIMEOUT_S=900
   VERIFY_GO_RACE_START_ALLOWANCE_S=995
   VERIFY_GO_RACE_TIMEOUT_S="${VERIFY_GO_RACE_TIMEOUT_S:-2000}"
   check_race_bounds() {
     case "$VERIFY_GO_RACE_TIMEOUT_S" in
       ''|*[!0-9]*) echo "verify_go.sh: FATAL: VERIFY_GO_RACE_TIMEOUT_S must be a positive integer" >&2; return 2 ;;
     esac
     min=$((VERIFY_GO_RACE_PKG_TIMEOUT_S + 60 + VERIFY_GO_RACE_START_ALLOWANCE_S))
     if [ "$VERIFY_GO_RACE_TIMEOUT_S" -lt "$min" ]; then
       echo "verify_go.sh: FATAL: leg bound ${VERIFY_GO_RACE_TIMEOUT_S}s must be at least per-package ${VERIFY_GO_RACE_PKG_TIMEOUT_S}s + 60s backup kill + start allowance ${VERIFY_GO_RACE_START_ALLOWANCE_S}s (${min}s)" >&2
       return 2
     fi
   }
   ```
   (`min` must not shadow anything; the implementer checks with `grep -n '\bmin=' scripts/verify_go.sh`.)
   Everything this step adds above the subcommand dispatch is assignments and function definitions only: nothing
   runs, prints or exits there. That keeps `--evidence-manifest-check` byte-identical in behaviour and output for
   `TestEvidenceNamedManifestRejectsUnpinnedTest`, which copies and executes this file [C1]. The new heredoc must
   not contain the text `REQUIRED_EVIDENCE_TESTS = {`, which that test uses as a mutation anchor via `strings.Index`.
2. `run_race_leg BOUND CMD...`, a shell function wrapping the existing heredoc. Python reads the bound from
   `argv[1]` and the command from `argv[2:]`, via `python3 - "$@" <<'PY'`. It keeps
   `start_new_session=True`, `p.wait(timeout=bound)` with the exit code passed through, `os.killpg(…, SIGKILL)`,
   `print(f"verify_go.sh: FATAL: -race leg timed out after {bound}s", file=sys.stderr)` and `sys.exit(124)`. NEW,
   before the kill: list every process whose pgid is the leg's, each as
   `verify_go.sh: still running at the leg bound (elapsed <etime>): <args>`, read from
   `ps -ww -eo pgid=,etime=,args=` and split with `line.split(None, 2)` (procps left-pads `pgid`; `split()` absorbs it). **r2 (astra, verbatim fix):** give `ps` an explicit two-second timeout (`subprocess.run(..., timeout=2)`) and put the
   process-group kill in a `finally` block, so a stalled, missing or failing `ps` can delay the SIGKILL by at most 2 s and can
   never prevent it or the `exit 124`; document that maximum two-second diagnostic grace beside the bound. If `ps` fails or
   times out, print one `could not list` line and still kill and exit 124 (diagnostics never gate the refusal).
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
with an EXIT trap; prints `ok <arm>` per arm and `PASS 7/7 arms`; any failure prints the arm and its captured
output and exits 1):
- **R1 default**: `env -u VERIFY_GO_RACE_TIMEOUT_S bash scripts/verify_go.sh --race-bounds-check` → rc 0, stdout
  exactly `race bounds: per-package 900s, start allowance 995s, leg 2000s`.
- **R2 boundary accept**: `VERIFY_GO_RACE_TIMEOUT_S=1955 … --race-bounds-check` → rc 0.
- **R3 boundary refuse**: `VERIFY_GO_RACE_TIMEOUT_S=1954` → rc 2, stderr contains
  `must be at least per-package 900s + 60s backup kill + start allowance 995s (1955s)`.
- **R4 malformed**: `VERIFY_GO_RACE_TIMEOUT_S=15m` → rc 2, `must be a positive integer`.
- **R5 passthrough**: `--race-leg-wrapper-check 30 sh -c 'exit 3'` → rc 3, and stderr has no `FATAL`.
- **R6 hang**: a fake executable `$tmp/hang_probe.test` (`#!/bin/sh`: `sleep 300 & echo $! >"$tmp/gc.pid"; sleep 300`)
  run as `--race-leg-wrapper-check 2 "$tmp/hang_probe.test"` → rc 124. Stderr contains `still running at the
  leg bound`, `hang_probe.test`, and `FATAL: -race leg timed out after 2s`. Elapsed ≥ 2 s and < 30 s. Known-positive
  control: `gc.pid` is non-empty. After return, a **bounded poll** (up to 5 s in 0.1 s steps) waits for
  `kill -0 "$(cat gc.pid)"` to fail. The arm fails only if the grandchild is still alive at 5 s, so the latency of
  reparenting and reaping the orphan cannot redden it. Under M10 the grandchild sleeps 300 s, well past the poll.
  Wall ~2–3 s, ≤ 7 s worst case. `$tmp` sits under `$TMPDIR`; on darwin that is a long `/var/folders/…` path, which
  is the truncation case [C3].
- **R7 stalled `ps` (r2, astra, verbatim fix)**: an R6 variant with a sleeping `ps` shim first on `PATH` (`#!/bin/sh`:
  `sleep 300`) and an independent ten-second watchdog around the wrapper call. Require the wrapper to return 124 and
  terminate the probe group (bounded `kill -0` poll as in R6) before the watchdog fires; stderr contains
  `could not list` and `FATAL: -race leg timed out after 2s`. Kills M20 (`ps` timeout removed) and M21 (kill moved
  out of `finally` after the listing).

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
    `worlddExtensionS = 93`, `daemonExtensionS = 166`, `daemonExtensionInflatedS = 347` (r2, §0 kimi model ii), `daemonInflatedNeedS = 828`, `backupKillS = 60` (Go's `testKillTimeout`, [D9]),
    `preLegMaxS = 55`.
  - Asserts: `P*4 >= slowNeedS*5` (P ≥ 1.25 × need); `A >= latestHeavyStartS + worlddExtensionS + daemonExtensionInflatedS` (r2) and `P > daemonInflatedNeedS` (r2)
    (≥ 995); `L_ci == L_default`; `L_default >= P + backupKillS + A`; `J*60 >= L_ci + 300` and `300 >= 5*preLegMaxS`; the call line
    contains `-race` and contains none of `-short`, `-run`, `-skip`, `-failfast=false`. Each failure message names
    the layer and the inequality.

**`host/verifygate/mission_config_gate_test.go`** (~+12/−4): in `requireVerifyGoProductSequence`, the
`repaired` arm takes the new banner needle and the new call-line needle in place of
`cmd = ["go", "test", …, "8m"]`. It drops the two Python needles (`print(f"… timed out after {bound}s"` and
`    sys.exit(124)`) from the *ordered* sequence, because they now live in `run_race_leg`, which is defined above the
`AILANG_BIN` anchor, and asserts instead that each occurs exactly once anywhere. The embedded `565d0b2` fixture
(`:415-442`) is history and stays byte-identical.

**`.github/workflows/ci.yml`** (~+5/−3): step "go -race gate": `timeout-minutes: 40`,
`VERIFY_GO_RACE_TIMEOUT_S: "2000"`. The comment names the three layers and points at this doc. No new job, no new
step (the wrapper test runs inside `verify_go.sh`). The edit keeps the controls `ailang-verify:`, `go-verify:` and
`./scripts/verify_go.sh`, and adds no `Z3_VER`/`Z3_SHA` token, both of which `TestZ3PinDeclaredOnceAndInstalledInBothJobs`
counts [C1].

## 5. Mutation table

Each mutation is a single edit to the M1 tree. "Kills" names the test and the assertion that fires.

| id | mutation | killed by |
|---|---|---|
| M1 | revert P: `VERIFY_GO_RACE_PKG_TIMEOUT_S=480` | `TestRaceBoundsAreOrdered` (P ≥ 1.25 × 646); `test_race_leg_bounds.sh` R1 (stdout text) |
| M2 | call line hardcodes `-timeout 8m` (P defined but unused) | `TestRaceBoundsAreOrdered` (exact call line); `TestVerifyGoProductCheckSequence` repaired arm (call-line needle count 0) |
| M3 | CI keeps `VERIFY_GO_RACE_TIMEOUT_S: "1500"` | `TestRaceBoundsAreOrdered` (`L_ci == L_default`, and `L ≥ P + 60 + A`) |
| M4 | script default lowered to `:-1900` (< 1955; r1's 1900 now fails the guard, so it is the mutant) | `TestRaceBoundsAreOrdered` (`L_default ≥ P + 60 + A`); R1 (rc 2, not 0) |
| M5 | `timeout-minutes: 36` (J×60 = 2160 < L + 300 = 2300) | `TestRaceBoundsAreOrdered` (`J*60 ≥ L + 300`) |
| M6 | `check_race_bounds` ordering branch deleted (or `return 0` first) | R3 (rc 0, want 2) |
| M7 | ordering comparison `-lt` → `-le` (off by one) | R2 (rc 2 at 1955, want 0) |
| M8 | `VERIFY_GO_RACE_START_ALLOWANCE_S=813` (r1's value, without the censored-daemon inflation) | `TestRaceBoundsAreOrdered` (A ≥ 995); R1 text |
| M9 | wrapper's process listing deleted | R6 (`still running at the leg bound` / `hang_probe.test` absent) |
| M10 | `os.killpg(...)` → `p.kill()` (child only) | R6 (grandchild `kill -0` succeeds) |
| M11 | `sys.exit(p.wait(...))` → `p.wait(...); sys.exit(0)` | R5 (rc 0, want 3) |
| M12 | timeout path `sys.exit(124)` → `sys.exit(0)` | R6 (rc 0, want 124) |
| M13 | `-short` appended to the call line | `TestRaceBoundsAreOrdered` (forbidden flag, and the exact line) |
| M14 | `-race` removed from the call line | `TestRaceBoundsAreOrdered`; existing `TestRaceLegRunsInItsOwnJob` |
| M15 | `continue-on-error: true` on the race job/step | `TestRaceBoundsAreOrdered` (zero `continue-on-error`) |
| M16 | `bash scripts/test_race_leg_bounds.sh` line dropped from the race block | `TestVerifyGoProductCheckSequence` repaired arm, *only if* the implementer adds it as a needle; **required**: add it |
| M17 | `check_race_bounds || exit $?` dropped from the main flow | **required new needle** in the repaired sequence, the same as M16 |
| M18 | the `+ 60` backup-kill term dropped from `min` | R3 (rc 0 at 1954, want 2) |
| M19 | a top-level `check_race_bounds` call (or any echo) inserted above the subcommand dispatch | `TestEvidenceNamedManifestRejectsUnpinnedTest` if it changes `--evidence-manifest-check` output; R1 if it prints twice |
| M20 | the `timeout=2` on the `ps` call removed (r2) | R7 (watchdog fires before rc 124) |
| M21 | the process-group kill moved after the listing, outside `finally` (r2) | R7 (watchdog fires; probe group alive) |

M16/M17 are listed so the implementer adds the two needles; without them those edits would survive. There are no
declared equivalents. One edit is **not killable locally**, and is declared: dropping `-ww` from the `ps` call. On
darwin the output is identical either way [C3], so R6 stays green. **r2 (kimi):** AC7 runs R6 as written (with `-ww`) and never replays the mutation, so the drop-`-ww` mutant is declared UNOBSERVED everywhere (locally and in CI), not killed.

## 6. Acceptance criteria

All commands run from the worktree root on the rig with `AILANG_BIN=$HOME/.pinned-ailang/ailang` (v0.41.0,
released), `GOTOOLCHAIN=go1.26.6`.

- **AC1** `bash scripts/verify_go.sh --race-bounds-check; echo rc=$?` →
  `race bounds: per-package 900s, start allowance 995s, leg 2000s` / `rc=0`.
- **AC2** `VERIFY_GO_RACE_TIMEOUT_S=1954 bash scripts/verify_go.sh --race-bounds-check; echo rc=$?` → stderr
  `… leg bound 1954s must be at least per-package 900s + 60s backup kill + start allowance 995s (1955s)` / `rc=2`;
  and with `=1955` → `rc=0`.
- **AC3** `bash scripts/test_race_leg_bounds.sh; echo rc=$?` → six `ok R<n> …` lines, `PASS 7/7 arms`, `rc=0`.
  It also passes with `TMPDIR` under a 1777 parent and under CPU contention, per the iteration-442 Gate 5 process
  rule: `TMPDIR=$(mktemp -d /private/tmp/r267.XXXX) bash scripts/test_race_leg_bounds.sh` with 8 `yes >/dev/null`
  spinners running → `PASS 7/7 arms`.
- **AC4** `go vet ./host/verifygate/ && go test ./host/verifygate/ -count=1 -v -run 'TestRaceBoundsAreOrdered|TestRaceLegRunsInItsOwnJob|TestVerifyGoProductCheckSequence|TestVerifyGoRefusesUnknownLegs|TestEvidenceNamedManifestRejectsUnpinnedTest|TestZ3PinDeclaredOnceAndInstalledInBothJobs' 2>&1 | grep -E '^(--- |ok|FAIL)'`
  → six `--- PASS` lines and `ok`. (The last two are the C1 readers that copy-and-execute `verify_go.sh` and count
  `ci.yml` tokens. AC6's whole `go test ./...` runs them again.) (`go vet` is the compile fence for `_test.go`; `go build` is not.)
- **AC5** Mutation replay: for each of M1–M19, apply the single edit, run the killing test(s) named in §5, and
  observe red with the named assertion; restore and observe green. Record a table `id | killing test | rc mutated |
  rc restored`. Expected: 19/19 killed, 0 survivors (M19 as the specific insertion the replay applies:
  `check_race_bounds` called above the dispatch with an echo). The declared `-ww` case is excluded locally and
  observed under AC7.
- **AC6** Whole gates on the M1 tree: `./scripts/verify_ail.sh` rc 0; `go build ./... && go test ./... -count=1 -timeout 45m` (r2, kimi: explicit `-timeout`, because the rig's plain `cmd/ailang-worldd` reads ~696 s against Go's 600 s default; compared against a same-session base run on `6420426e`)
  with the fail set compared to base `6420426e`, so no new failure; `gofmt -l host/verifygate/` empty;
  `bash -n scripts/verify_go.sh scripts/test_race_leg_bounds.sh` rc 0.
- **AC7** CI on the PR: the race job's log contains `── legs: race (race bounds: per-package 900s, leg 2000s)`,
  the shell test's `PASS 7/7 arms`, and
  `── go test ./... -count=1 -race -timeout 900s (leg bound 2000s)`, and all three jobs are green. R6's `ok` line in
  that log is the Linux procps observation for `ps -ww … args=` [C3].
- **AC8 (observational; recorded, never claimed if unobserved; r2 text is glm's verbatim fix)**: over the first 10 race-job
  attempts after merge, (i) no `panic: test timed out after 15m0s` appears in any attempt, and (ii) every heavy-package
  instance reading > 480 s is classified as green or legitimately-red — red only with a non-timeout failure cause that
  reproduces on the same commit (e.g., in the plain leg or at base `6420426e`); record a per-attempt classification table
  (attempt | heavy>480s | result | cause class). If no attempt draws a non-normal runner (`host/store` ≤ 67.5 s in all 10),
  record AC8 as unobserved, not passed.

## 7. Conflict surface (who else reads these bounds or these lines)

Greps in [D8]:

- `host/verifygate/mission_config_gate_test.go:436-437, 500-501`: exact-text needles for the race banner and
  `cmd = [... "8m"]`. The fixture lines (`:436-437`) stay; the shared needles (`:500-501`) become `repaired`-
  parameterised (§4). **Must change in M1 or that test goes red.**
- `host/verifygate/race_job_gate_test.go:37, 55-57`: `VERIFY_GO_RACE_TIMEOUT_S:` presence and the Python `cmd`
  substring. The latter must change (§4).
- **The three D10 readers r0 left unclassified** [C1]:
  - `scripts/floor/canary.py:95`: a docstring only ("verify_go.sh's tracked-binary gate refuses .gz"). It pins
    no text, and nothing executes `canary.py` from CI or a script (grep: no match, with a positive control on
    `ci.yml`; [C1, D13]). Not affected; no AC needed.
  - `host/verifygate/evidence_manifest_gate_test.go` (`TestEvidenceNamedManifestRejectsUnpinnedTest`; helpers at
    `:98,106,169,187,212`): **copies and executes** `verify_go.sh --evidence-manifest-check`, and rewrites the text
    between `REQUIRED_EVIDENCE_TESTS = {` and `}\nEXACT_EVIDENCE_TESTS = int(sys.argv[2])`. Affected only if M1 adds
    executing top-level code before the dispatch, or reuses those anchors. §4 step 1 forbids both. Exercised by AC4
    (named) and AC6 (whole suite); M19 is its mutation.
  - `host/verifygate/ail_binary_gate_test.go:708` (`TestZ3PinDeclaredOnceAndInstalledInBothJobs`): requires the
    `ci.yml` controls `ailang-verify:`, `go-verify:` and `./scripts/verify_go.sh`, and counts `Z3_VER:`/`Z3_SHA:`.
    M1's `ci.yml` edit (two values and a comment in one step) touches none of them. Exercised by AC4 and AC6.
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

One milestone. ~+240/−22 ≈ 262 changed lines: `verify_go.sh` ~61, `test_race_leg_bounds.sh` ~80,
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
- Not shrinking A to the measured-schedule value 647 s [D11]. The stacked 813 s costs nothing on a passing run.
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
| D11 | slot cascade (C4, kimi's check): which slot `host/verifygate` takes | `python3 -I design-scratch/slots.py` (output `slots.out`; start = print offset − duration, as D4) | `both cmd&daemon >=440s: 15 \| vg_start - cmd_end range -83.0..+1.3 \| vg_start - daemon_end range -172.7..-91.4`; `all: 51 \| vg_start - cmd_end range -83.0..+16.7 \| vg_start before daemon_end in 51 of 51`; e.g. `45936239736_a1 cmd_end 534 dmn_end 620 \| vg_start 451` |
| D12 | the plain leg passes no `-timeout` (correction 6) | `grep -n -- '-timeout' scripts/verify_go.sh` (positive control: it matches the two race lines) | `350:echo "── go test ./... -count=1 -race -timeout 8m …"`, `355:cmd = [… "-timeout", "8m"]`; no other line; = C2 |
| D13 | C1 readers and their test names | `grep -n '^func Test' host/verifygate/evidence_manifest_gate_test.go`; `awk 'NR<=708 && /^func Test/{l=NR": "$0} END{print l}' host/verifygate/ail_binary_gate_test.go`; `grep -n verify_go scripts/floor/canary.py`; `grep -rn canary.py --include='*.sh' --include='*.yml' --include='*.go' --include=Makefile . \| grep -v ^./design_docs` with control `grep -rln verify_go.sh --include='*.yml' . \| grep -v design_docs` | `229:func TestEvidenceNamedManifestRejectsUnpinnedTest` (the file's only `Test`); `699: func TestZ3PinDeclaredOnceAndInstalledInBothJobs`; `95:` a docstring, `Evidence transcripts stay plain .jsonl: verify_go.sh's tracked-binary gate refuses .gz.`; **canary.py executors: none**, with control `.github/workflows/ci.yml` |
| C1 | readers (VERIFIED BY CONTROLLER) | as in the r1 brief: `grep -n 'verify_go' scripts/floor/canary.py host/verifygate/evidence_manifest_gate_test.go host/verifygate/ail_binary_gate_test.go` | canary docstring `:95`; evidence test copies+executes at `:98,106,169,187,212`; ail_binary controls `:708`; control `grep -c verify_go host/verifygate/mission_config_gate_test.go` → `25` |
| C2 | plain leg `-timeout` (VERIFIED BY CONTROLLER) | `grep -n -- '-timeout' scripts/verify_go.sh` | exactly `350` and `355` |
| C3 | `ps` truncation (VERIFIED BY CONTROLLER, darwin/arm64 only) | 146-char executable, `start_new_session=True`, listed from Python with `stdout=PIPE`, with and without `-ww` | full path present, longest line 3,282 chars, identical output; Linux procps unmeasured |
| C4 | start-allowance check | not measured by the controller; done here as D11 | — |
| D10 | existing pins of `verify_go.sh` text | `grep -rln 'verify_go' --include='*.go' --include='*.sh' --include='*.yml' --include='*.py' . \| grep -v '^./design_docs'` | `host/verifygate/{mission_config_gate,evidence_manifest_gate,ail_binary_gate,toolchain_pin_gate,race_job_gate}_test.go`, `host/pins/scan_test.go`, `host/{broker,transitionreg,runbook,replay}` (prose/comments), `scripts/{verify_go,tracked_binary_gate,test_tracked_binary_gate}.sh`, `scripts/floor/canary.py`, `.github/workflows/ci.yml` |

Inherited, VERIFIED BY CONTROLLER (iteration 443): `race_table.txt` and `racelogs/` (53 attempt logs, run ids,
shas, events, conclusions); the latest instance of §1. Each number used above was re-derived from the raw logs
(D2–D6), not from the table's text.

Not measured (each marked PENDING in the body): Linux procps `ps -ww` under a pipe (observed by AC7); the runner vCPU count; per-shard `-race` times (option C); a full
local `VERIFY_GO_LEGS=all` run on the M1 tree (AC6 covers it at implementation).

## 11. Quorum verification log

**Round 1** (on r0 `8256a83e`; author `claude:claude-opus-5-5` benched; seats gemini-4-carbon, oc-glm-5-3,
oc-kimi-k3, no absent seat): **BLOCKED 3/3**. No seat disputed the direction (option A, fixed per-package bound,
strict layer ordering).

| seat | objection (one line) | r1 change |
|---|---|---|
| gemini-4-carbon | `ps -eo … args=` may truncate argv (BSD without `-ww`, non-tty) and was asserted without a log row; A = 720 omits `cmd/ailang-worldd`'s 92.9 s extension | `ps -ww` adopted as a belt (§3, §4), with C3: darwin under a pipe does **not** truncate (premise refuted there), and Linux is observed by R6 in CI (AC7). The `-ww` drop is a declared locally-unkillable edit (§5). A recomputed with both extensions stacked: 813 (gemini proposed 820 with L ≥ P + A; r1 also adds the 60 s backup kill, per kimi) |
| oc-glm-5-3 | §7 omits three D10 readers (`canary.py`, `evidence_manifest_gate_test.go`, `ail_binary_gate_test.go`); correction 6's "plain leg passes no `-timeout`" was unverified | §7 classifies all three [C1, D13]: canary is a docstring with no executor; the evidence test copies and executes the script (§4 step 1 constraint, M19); the Z3 test counts ci.yml controls. AC4 names both Go tests and AC6 runs them. Correction 6 now cites C2/D12 |
| oc-kimi-k3 | A = 720 is under-derived against the doc's own stacked model; the guard omits Go's 60 s backup kill; R6's single `kill -0` can race the reaping | Slot cascade measured (D11): verifygate never waits on daemon's slot (51/51), so it tracks cmd's. The stacked A = 813 is kept anyway as the conservative reading. Guard is now L ≥ P + 60 + A = 1773, and L is raised to **1900** (127 s margin). J = 40 re-derived (2400 ≥ 2200). R1–R3, AC1/AC2, the test constants (`worlddExtensionS`, `backupKillS`), M7/M8 updated, M18 added. R6 is now a bounded 5 s poll |

r0 → r1 diff (derived from `diff design-scratch/r0.md` against this file): header/instruments; correction 6 citation;
§2 slot-cascade paragraph; §3 item 1 (A, guard formula) and item 2 (`ps -ww`, C3); bounds table rows A/L/J and the
ordering line; §4 step 1 constants, comment, guard, FATAL text and the C1 constraint; step 2 `ps -ww` and parsing;
R1/R2/R3/R6; test constants and asserts; ci.yml value 1900 and the C1 controls note; §5 M3/M4/M5/M7/M8, plus
M18/M19 and the declared `-ww` case; §6 AC1/AC2/AC4/AC5/AC7; §7 three-reader bullet; §8 sizing; §9 A note;
§10 D11–D13, C1–C4; §11.

### Round 2 (r1 `4c76f7c0`) — BLOCKED 2/3; resolved by the controller's bounded r2 under the narrow-refinement carve-out

Seats: oc-glm-5-3, oc-kimi-k3, codex:gpt-6.2-astra (reserve gemini-4-carbon). FLAG: gemini-4-carbon (round-1 seat) was not
re-seated — the controller passed a four-model `--reviewers` pool and the tool seated astra; the rule is the same three seats
every round. Artifact `w-daemon-race-package-timeout-headroom-2026-11-27T04-32-02Z.json`.

| seat | verdict | objection (one line) | r2 change |
|---|---|---|---|
| oc-glm-5-3 | PASS | AC8 mixes timeout-class and real-failure reds; §7's row-243 quote unlogged; one unlogged timestamp | AC8 replaced by glm's text; quote and timestamp marked in §0 |
| oc-kimi-k3 | REJECT | daemon slow-tail per-`k` is censored; under `cmd`'s measured inflation A ≈ 994, guard ≥ 1955, L = 1900 breaks | A 995, guard 1955, L 2000, test constants, M4/M8, R1–R3, AC1/AC2 (kimi's fix); AC6 `-timeout`; `-ww` mutant unobserved |
| codex:gpt-6.2-astra | REJECT | unbounded `ps` before `os.killpg` on the timeout path | `ps` timeout 2 s, kill in `finally`, R7 + M20/M21 (astra's fix) |

