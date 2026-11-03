# w-exec-probe-hello-wait-pin-survivors: pin pollLoopback's real dial cadence to execProbePollEvery, and settle the row-230 judge's survivors

Status: Planned (iteration 339), **r1**: designer draft; design quorum r1 PROCEED 3/3 (§7).
**Planner-Lane:** mission routing default
- **Item**: queue row 241 (`design_docs/world-mission.md:6313`). Source: the row-230 judge
  (`verification/world-iter329/evaluator-sprint-r1.md`, non-blocking survivors 1–5 at lines 10–14 and its
  mutation table at lines 20–30). Parent design: `planned/w-exec-probe-reverse-leg-pins.md` (row 230, which
  added `execProbeHelloWait`'s two bounds, `execProbePollEvery`, T10 and T11). Row 242 (the cross-profile
  `net6` hoist) is **out of scope**: it changes probe structure and needs its own design.
- **Kind**: clause-3 hygiene (MET). It moves no clause. No `.ail` changes. **Test-only**: zero production
  bytes (§2.4 says why there is no production seam).
- **Base**: `dcb2366e` (= `origin/dev`), branch `sprint/world-iter339-hello-wait-cadence`; go1.26.6
  darwin/arm64, 16 cores, pinned AILANG v0.41.0 (V01).
- **Estimated**: ~0.15d, 1 executor run.
  - `host/broker/exec_probe_test.go`: +59/−4. One new test (T12), T10's ceiling re-referenced, one field on
    `helloConn`.
  - `.github/workflows/ci.yml`: one test name added to one list.
- **Evidence**: `design_docs/verification/world-iter339/designer-*.txt`.
  - **Where.** Every measurement ran on a `git archive` scratch copy of `dcb2366e` under
    `/tmp/world-iter339/designer/`, never in the worktree. `base` is HEAD. `proto` is HEAD plus the §2 test
    diff (`designer-proto_test.go.patch.txt`, sha256 `bbfccede…`). `mutA`/`mutB`/`mutC` are copies of `proto`
    for mutants.
  - **Mutants.** Each mutant is an exact single-occurrence replacement. `designer-mut.py.txt` exits 3 if the
    site count is not 1. Each mutant was restored from a pristine copy and checked with `cmp`.
  - **Environment.** Every `go test` exported `AILANG_BIN=$HOME/.pinned-ailang/ailang`.
  - **Load.** N spinning `yes` processes, started by the runner and killed by a trap. The runner prints
    `yes-alive-after=0` after each run.
  - **The rig was not quiet.** From 09:25 to about 11:08, another session's AILANG release builds held the
    16-core rig at load 80–214. Every mutant line carries its `uptime`. Those runs are harsher than the
    designed conditions, not easier (§2.3).

## 1. Problem (measured at the base)

T10 (`TestExecProbeReverseHelloWaitBounds`, `exec_probe_test.go:892`) bounds the hello wait from above by
`ceiling := listenWindow - execProbePollEvery` (`:896`). The premise is that the first connected dial lands at
most one poll period after the child binds.
- **What T10 subtracts.** The constant `execProbePollEvery` (`exec_probe.go:131`).
- **What the poll waits.** Whatever its only timer is given: `case <-time.After(execProbePollEvery):`
  (`exec_probe.go:763`, the file's only `time.After`).
- **The gap.** The constant is read in exactly those two places (V02). No test observes the interval the
  poll waits, so the two can drift apart.

The judge's five items, re-measured:

| Item | Mutant | At the base | V |
|---|---|---|---|
| 1 | `:763` → `time.After(500 * time.Millisecond)` (M1) | **rc 0** on the reverse-leg subset; control rc 0; positive control rc 1 (constant 2 s). M1's broad-regex FAIL list on `proto` adds only T12 (plus one load flake in an unchanged real-socket test), so the base suite cannot see M1 | V03 (controller), V05 |
| 2 | `:613` SILENT text `within %d ms` → `within 500 ms`, argument dropped | rc 0 (equivalent at 500 ms) | V06 |
| 3 | `:128` wait 951 ms, 1 ms over T10's ceiling | Killed 1,200/1,200 subtests at 0, 32 and 160 `yes`. But the hand-off T10 measures through (classifier clock read → poll exit → test wake-up) reached **1.713 ms under 160 `yes`**, more than the 1 ms margin, on a 955 ms run. So a stalled run can let a just-over mutant survive: the judge's survivor is real, and it just did not land on a 951 run | V07 |
| 4 | T10's arithmetic (ceiling without the poll, floor 100/500, arrival at start, `sets`/`wrote` dropped) | each rc 0 alone | V06 |
| 5 | constant 100 ms / 25 ms, listen 1100 | rc 0 (in bounds) | V05 |

## 2. Design

### 2.1 T12 `TestExecProbeReversePollCadence` (new, in `exec_probe_test.go` after T11)

**What it observes.** The real loop, through the existing `dial` seam. The fake dial counts the dial,
records `time.Now()` for the first 80, and returns ECONNREFUSED. So `classifyReverseDial` never runs, and each
loop iteration is dial, then the timer, then the next dial.

**The clock is a reference loop, not wall time.** The test goroutine waits `time.After(execProbePollEvery)`
20 times, the same timer the poll should run, and then stops the poll. Load that delays the poll's wake-ups
delays the reference's too. §2.3 shows, by measurement, why wall time cannot be the clock.

- **Count gate.** A poll that waits `execProbePollEvery` dials once at the start and then once per period:
  about 21 times while the reference runs. The test requires at least `revCadencePeriods/4+1` = 6.
  - M1 makes 3 dials (`3 dials while a reference waited execProbePollEvery (50ms) 20 times`) and reds (V05).
  - A correct tree reds only if the poll averages about four reference periods per dial.
- **Floor gate.** No gap between consecutive dials may be under `execProbePollEvery`. It reports once, with
  the count and the fastest gap.
  - It kills a faster cadence: `/10`, 1 ms, 45 ms, a busy loop.
  - It kills a 50 ms literal left behind when the constant moves to 100 ms (V05).
- **Record cap.** Timestamps are recorded only up to `4*revCadencePeriods`. A busy-loop mutant made
  6,753,851 dials in 2.1 s. The first prototype's per-gap `t.Errorf` wrote a 1.86 GB log for it (V14). The
  cap and the single floor message fix that.

The reference code is the exact T12 in `designer-proto_test.go.patch.txt`, the file every final measurement
ran:

```go
func TestExecProbeReversePollCadence(t *testing.T) {
	const revCadencePeriods = 20
	refused := revDialErr(syscall.ECONNREFUSED)
	var mu sync.Mutex
	dials, at := 0, make([]time.Time, 0, 4*revCadencePeriods)
	stop := pollLoopback(40000, "WORLD-PROBE-NONCE "+strings.Repeat("a", 32), func(string) (net.Conn, error) {
		mu.Lock()
		if dials++; len(at) < cap(at) { // a runaway mutant must not grow the slice without bound
			at = append(at, time.Now())
		}
		mu.Unlock()
		return nil, refused
	})
	start := time.Now()
	for range revCadencePeriods {
		<-time.After(execProbePollEvery)
	}
	span := time.Since(start)
	stop()
	mu.Lock()
	defer mu.Unlock()
	if dials < revCadencePeriods/4+1 {
		t.Errorf("%d dials while a reference waited execProbePollEvery (%v) %d times (%v); want at least %d: the poll does not wait execProbePollEvery", dials, execProbePollEvery, revCadencePeriods, span, revCadencePeriods/4+1)
	}
	short, fastest, slowest := 0, time.Duration(0), time.Duration(0)
	for i := 1; i < len(at); i++ {
		gap := at[i].Sub(at[i-1])
		if i == 1 || gap < fastest {
			fastest = gap
		}
		slowest = max(slowest, gap)
		if gap < execProbePollEvery {
			short++
		}
	}
	if short > 0 {
		t.Errorf("%d of %d dial gaps under execProbePollEvery (%v), the fastest %v; want none: the poll waits less than execProbePollEvery", short, len(at)-1, execProbePollEvery, fastest)
	}
	t.Logf("CADENCE dials=%d span=%v fastest=%v slowest=%v", dials, span, fastest, slowest)
}
```

The doc comment above it (in the patch) states the reference-clock argument and both gates.
`exec_probe_test.go` already imports `sync`.

**No bound duplicates a constant.** The reference and the floor both read `execProbePollEvery`. The test
parameters (20 periods, the factor ¼, the cap of 4 × periods) are test shape, not production values. Change
the constant and every bound moves with it: K-const100, K-const25 and the constant at 1 ms stay green (V05,
V08).

**What T12 does not pin: the accepted band.** Literal cadences of 2×, 3× and 4× `execProbePollEvery` make
11, 8–9 and 6 dials, and survive 30/30 each. 5× makes 5 dials and reds 30/30 (V08).
- Inside the band, T10's ceiling (950 ms) can be off by up to three poll periods (150 ms).
- Today's wait sits 450 ms under the ceiling.
- The row's mutant (10×) and anything at or above 5× are killed.

### 2.2 T10: reference the ceiling to the instant the deadline was set (survivor 3)

`helloConn.SetDeadline` also records `setAt = time.Now()`. T10's ceiling check becomes
`conn.deadline.After(conn.setAt.Add(ceiling))`. The `after := time.Now()` stop-time reference is removed,
and the message says `set+` instead of `stop+`.

- **A correct tree cannot red, at any load.** The classifier computes `deadline = Now_c + wait` and then
  calls `SetDeadline`, which reads `setAt` on the same goroutine. So `setAt >= Now_c`, and
  `deadline − setAt <= wait` for every tree. With wait 500 ms and ceiling 950 ms the check cannot fire.
  Wait 950 ms, the boundary, stayed green 400/400 at each load (V07).
- **It is strictly stronger than the base check.** `stop()` returns only after the poll goroutine closes
  `finished`, which follows the classify call. So `setAt <= after`, and
  `deadline − after <= deadline − setAt`. Every tree red at the base is red here.
- **The kill margin, measured (V07).** For the 951 ms mutant, the margin is how far the reported excess sits
  above the 950 ms ceiling.
  The margin is lost when the reference lags the classifier's clock read by more than 1 ms. Across 1,600 red
  subtests per tree and load (waits 951/952/955/960), the worst lag was
  (`designer-s3-aggregate.txt`):

  | `yes` | Base (stop-referenced) | Proposed (setAt-referenced) |
  |---|---|---|
  | 0 | 80.0 µs | 0.3 µs |
  | 32 | 110.1 µs | 0.3 µs |
  | 160 | **1,713.3 µs** | 0.5 µs |

  - At the base, a lag of 1.713 ms at 160 `yes` exceeds the 1 ms margin. A 951 ms mutant landing on such a
    run survives. The judge's survivor 3 is real.
  - The new reference removes the goroutine hand-off from the margin. What remains is two adjacent clock
    reads on one goroutine: ≤ 0.5 µs measured, three orders of magnitude inside 1 ms.

### 2.3 Load robustness (each bound, why it cannot false-red, the worst value observed)

| Bound | Derived from | Why a correct tree cannot red | Worst observed (V10) |
|---|---|---|---|
| T12 floor: each dial gap ≥ `execProbePollEvery` | `execProbePollEvery` | A Go timer never fires before its duration, and the next dial follows the timer. Both use the runtime's monotonic clock. Load only lengthens a gap | Smallest gap seen, all versions and loads: 50.013 ms (v2, v3); v1 50.032 ms. 0 floor reds in 3,678 correct-tree runs (v1 1,700, v2 1,178, v3 800) |
| T12 count: ≥ 6 dials in 20 reference periods | `execProbePollEvery` (the reference's timer) | Load lands on both loops' timers. A false red needs the poll to average about 4 reference periods per dial | v3: min 18 dials (gate 6) in 800 runs, worst lag 1.18 periods per dial (gate admits 4.0). v2 under ambient 150–214: 2.00 |
| T10 ceiling: `deadline − setAt <= listenWindow − execProbePollEvery` | `execProbeListenMS`, `execProbePollEvery` | `deadline − setAt <= wait` by construction (§2.2) | 0 reds in every load run (v1 1,700, v2 1,178, v3 800 runs; 2 subtests each), including v1's 76-false-red ambient run |
| T10 floor (unchanged) | `execProbeHelloFloor` | `helloConn` has no clock and no sleep | 0 reds |

**How T12 got here: two rejected prototypes, both measured.**

- **v1, an absolute ceiling (V09, V11).** T12 v1 (`designer-proto-v1_test.go.patch.txt`) required the
  fastest of up to 8 dial gaps to be under 2 × `execProbePollEvery`. It collected them inside one listen
  window, 1 s of wall time.
  - It held for 1,500 runs at load up to 77, including `-race` with `GOMAXPROCS=2`. The worst fastest gap
    was 52.08 ms.
  - At ambient load 150–170, with 32 more `yes`, `-race` and `GOMAXPROCS=2`, v1 **false-red 76 of 200
    runs**: 35 saw one dial in the window, and 41 had a fastest gap ≥ 100 ms. T10 passed 200/200 in the same
    runs.
  - Single gaps reached 2.8 s. Any absolute wall-clock ceiling fails under starvation.
- **v2, a reference clock with 10 periods and a gate of 6 (V12,
  `designer-proto-v2_test.go.patch.txt`).**
  - It held at idle (min 10 dials), under 32 `yes` (min 8; one 500 ms reference took 9.2 s), under `-race`
    (min 9), and under 64 `yes` with `-race`, `GOMAXPROCS=2` and ambient load ~206 (200 runs, min 9).
  - Under 32 `yes` with `-race` at ambient load 150–214, 78 runs completed, all green. Then the harness's
    10-minute binary timeout fired while one 500 ms reference loop took 32.6 s; this was an instrument
    timeout, not a verdict. In one completed run, the reference took 62.6 s with a 29.1 s single gap and
    still saw 11 dials. **But the minimum was 6 dials, exactly the gate.**
  - Zero margin is not robust. v3 doubles the averaging window (20 periods) and halves the gate ratio (¼).
    v2's worst case (6 dials in 10 periods, two periods per dial) would give about 11 dials against v3's
    gate of 6.
- **v3, the shipped T12 (V10).** 800 runs, 0 FAIL (T12 and T10).
  - Conditions: ambient load 10–24 (200 runs); `-race` + 32 `yes` (200; end load 53); `-race` + 64 `yes`
    + `GOMAXPROCS=2` (200; end 87); `-race` + 160 `yes` + `GOMAXPROCS=2` (100; end 187); 160 `yes` (100; end
    194).
  - Minimum 18 dials against the gate of 6. Worst lag 1.18 periods per dial. Floor minimum 50.013 ms.
    Slowest single gap 1.29 s, and one reference took 6.1 s; the count stayed ≥ 18.
  - **Not measured:** v3 under the other session's 150–214 ambient load, which was what broke v1 and took
    v2 to its gate. That load ended before v3 existed. The only evidence for v3 there is v2's 2.00 lag, which
    is half of v3's 4.0 admission.

**Load erodes the M1 kill, toward survival, never toward a false red.** Each wake-up's delay accumulates
20 times in the reference but only 2–4 times in a 10× poll. So under load, M1 makes more dials:
- 3 dials at load ≤ 140 (V05);
- 5 dials in 2.35 s at load 92, with the gate removed (TM1-nocount+M1, V13).

By that arithmetic, M1 survives once the average wake-up delay nears 2 × `execProbePollEvery` (100 ms). The
same load cannot red a correct tree, because the correct poll accumulates the same 20 delays as the
reference.

**What the reference clock does not cover.** If the OS starved only the poll's thread while the test's
thread ran, the poll could fall behind the reference. The gate admits about a 4× lag. The worst lag measured
across every v2 and v3 run is 2.00 (v2, ambient 150–214; v3's worst is 1.18). This is a measured margin, not a proof. §5 records it as the
residual.

### 2.4 Why test-only (production alternatives considered)

- **A timer seam.** This would be a `wait func(time.Duration) <-chan time.Time` parameter beside `dial`, or a
  package variable.
  - For: a test could record the requested duration exactly, with no timing at all.
  - Against: a test that injects its own timer cannot see the seam's **default**. M1 just moves into the
    default (`time.After` replaced by a 10× wrapper) and survives. The `dial` seam's own default, the
    200 ms `DialTimeout` at `:735`, is unpinned for the same reason (argued; not mutated here).
  - Against: it adds a seam to a fail-closed startup gate.
  - This is argued, not prototyped. T12 already kills M1 by observing the real default path, so no
    production change is needed.
- **A max-gap ceiling.** Rejected by measurement. The slowest single gap was 107–231 ms even with no added
  load, and up to 29.1 s under load (V09, V12). A correct tree would red.
- **A dial count, or a fastest gap, against wall time.** This is v1, which false-red 76/200 (V11).

### 2.5 The judge's survivors

1. **Killed.** M1 → T12's count gate, measured with T12 alone in the FAIL list (V05).
   - Also killed by the count gate: 10×, 5×, and 60 s (1 dial).
   - Also killed by the floor gate: `/10`, 1 ms, 45 ms, a busy loop, and a 50 ms literal left behind when the
     constant moves to 100 ms.
   - 2×–4× survive as the accepted band (§2.1).
2. **Equivalent at 500 ms; accepted, no pin.**
   - The mutant alone is rc 0 (S2-msg).
   - It becomes observable only when the wait changes. The wait alone at 600 ms is rc 0, but the pair reds T7
     on all three rows: `… within 500 ms …` against the wanted `within 600 ms` (S2-msg+wait600).
   - Dropping the argument while keeping `%d` fails the `go test` build through vet's printf check:
     `fmt.Sprintf format %d reads arg #2, but call has 1 arg` (S2-argdrop).
   - So every observable form is killed. The only survivor is a text that is true at today's value. A pin
     would have to be a source-text pin of a log string; it would buy nothing behavioural (V06).
3. **Killed by re-referencing T10.** T10 now references the instant the deadline was set (§2.2). The
   judge's "~0.9 ms" is measured (V07):
   - The base's worst hand-off lag was 1.713 ms under 160 `yes`, which loses the 1 ms margin.
   - The proposed reference lags by ≤ 0.5 µs. The 951 ms mutant is red 1,200/1,200 on both trees in the
     sampled runs.
   - TM11 moves the proposed T10's reference back to a clock read after the poll. It still kills 951 in a
     single run (V13); the margin argument, not this run, is what the change buys.
4. **Accepted.** Each T10-arithmetic mutant is green alone, because at today's wait it differs from T10 only
   when some production mutant is also present. Each clause was paired with the production mutant it
   guards (V06):
   - **Ceiling** (`ceiling := listenWindow`): with it, wait 960 survives; without it, wait 960 reds
     (`set+959.99996ms … no later than set+950ms`). The ceiling clause is wait 960's sole killer. The stricter
     `− 2*execProbePollEvery` is green (equivalent at 500 ms).
   - **Floor value** (100 ms): with it, wait 249 survives; without it, wait 249 reds. Floor 500 is green alone
     (equivalent at a 500 ms wait). Floor 600 reds a correct tree, so the floor can never be raised past the
     wait.
   - **Arrival** (`arrival: before`): with it, wait 5 survives; without it, wait 5 reds.
   - **`wrote`**: an empty hello (`io.WriteString(c, "")`) reds T10 **and** T7's `silent-after-chatty`
     (`picker handed out 1 ports, want at least 2`). With `wrote` dropped, T7 still kills it. So the judge's
     "the only guard that the hello is written" is **refuted**: there are two guards.
   - **`sets == 0`**: redundant with `deadline.IsZero()`. Removing `SetDeadline` reds T10 with `sets` dropped,
     and is green only with both clauses dropped. That redundancy is harmless and is kept.
5. **Equivalent in bounds; accepted.** The constant at 100 ms, 25 ms and 1 ms (30/30), listen 1100 and wait
   950 all stay green, with T12 following the constant (V05, V08). Listen 500 reds T10, which is row 230's
   intended K9 coupling (judge table; not re-run here).

### 2.6 Not in scope

- Row 242's `net6` hoist.
- The dial timeout literal `200*time.Millisecond` (`exec_probe.go:735`). It is not the cadence, and no
  survivor names it.
- Any production byte, comments included. `exec_probe.go:130` already names the constant as the cadence, and
  T12 now enforces it.
- The real-socket reverse-leg tests' behaviour at load ≥ 90 (§5). Their unchanged code went red in the
  control run at load 92.

## 3. Mutation plan (final file; `designer-mutants-{A,B,C,borderline}.txt`)

Regexes:
- **F** (focused): `^(TestExecProbeReversePollCadence|TestExecProbeReverseHelloWaitBounds|TestExecProbeReverseLegSilentHolderRefusesSticky|TestExecProbeReversePollRecordsConnectReset|TestExecProbeReverseDialErrIsSilent)$`.
- **W** (broad, the controller's): `^(TestExecProbeReverse|TestExecProbeForeignListener|TestExecProbeChildListener)`.
- **T10**: `^TestExecProbeReverseHelloWaitBounds$`.
- **B**: T12 alone at `-count=30`.

Sites are in `host/broker/exec_probe.go` (P) or `exec_probe_test.go` (T) on `proto`. Apply each mutant with
`designer-mut.py.txt FILE OLD NEW` (exact, single occurrence) and restore from `git show HEAD:<file>`.

**Production mutants**

| # | Mutant (old → new) | Must red | Measured |
|---|---|---|---|
| M1 | P `case <-time.After(execProbePollEvery):` → `…(500 * time.Millisecond):` | T12 count | F: T12 only, `3 dials … want at least 6`. W: T12 plus `ExhaustionRefusesDistinctly/chatty-every-time` at load 137 (an unchanged real-socket test; see C0) |
| M1-10x | → `(10 * execProbePollEvery)` | T12 count | 3 dials, T12 only |
| M1-5x | → `(5 * execProbePollEvery)` | T12 count | 30/30 red, 5 dials |
| M1-60s | → `(60 * time.Second)` | T12 count | 1 dial, T12 only |
| M1-div10 | → `(execProbePollEvery / 10)` | T12 floor | `76 of 79 dial gaps under …`, T12 only |
| M1-1ms | → `(time.Millisecond)` | T12 floor | `40 of 55 …`; at load 142 also T7 `silent-after-chatty` (not in the earlier round at load 75: load-sensitive, `designer-mutants-A-v3pre.txt`) |
| M1-45ms | → `(45 * time.Millisecond)` | T12 floor | 30/30 red, `22 of 22 …` |
| M1-busy | `case <-time.After(execProbePollEvery):` → `default:` | T12 floor | `79 of 79 …`, 6,753,851 dials, T12 only |
| B-2x / 3x / 4x | → `(k * execProbePollEvery)`, k = 2, 3, 4 | none (accepted band) | 30/30 green each; 11, 8–9, 6 dials |
| M1-ctl50 | → `(50 * time.Millisecond)` | none (equivalent) | green |
| M1-ctl50 + const 100 | the literal, and `const execProbePollEvery = 100 * time.Millisecond` | T12 floor | `23 of 32 … (100ms)`, T12 only |
| PC (positive control) | `const execProbePollEvery = 2000 * time.Millisecond` | T10 ceiling | T10 both rows `… no later than set+-1s`. T12 green: it follows the constant (21 dials in 40.7 s) |
| K-const100 / 25 / 1ms | the constant | none | green; 1 ms 30/30 |
| K-listen1100 | `const execProbeListenMS = 1100` | none | green |
| K-wait950 | `const execProbeHelloWait = 950 * time.Millisecond` | none (boundary) | green; 400/400 in V07 |
| K-wait951 | `… = 951 …` | T10 ceiling | 400/400 red at each of 0, 32 and 160 `yes`, on both trees (V07) |
| K-wait960 | `… = 960 …` | T10 ceiling | `set+959.999958ms … set+950ms` |
| K-wait249 / K-wait5 | `… = 249 / 5 …` | T10 floor | `seen {silent:true} … not silent` |
| S2-msg | `within %d ms; … refuses", port, execProbeHelloWait.Milliseconds())` → `within 500 ms; … refuses", port)` | none (equivalent) | green |
| S2-wait600 | wait 600 | none | green |
| S2-msg + wait600 | both | T7 ×3 | `… within 500 ms …` vs `within 600 ms` |
| S2-argdrop | `refuses", port, execProbeHelloWait.Milliseconds())` → `refuses", port)` | build | vet printf; `[build failed]` |
| P-nohello | `io.WriteString(c, execProbeHello)` → `io.WriteString(c, "")` | T10 `wrote`, T7 | T10 both rows (`wrote ""`) and T7 `silent-after-chatty` |
| P-nodeadline | `_ = c.SetDeadline(time.Now().Add(execProbeHelloWait))` → `_ = 0` | T10 | `set+0s (sets 0)` |

**Mutants of the test code (T12 new, T10 edited; survivors are expected and say which clause is load-bearing)**

| # | Mutant | Expect | Measured |
|---|---|---|---|
| TM1 | `if dials < revCadencePeriods/4+1 {` → `if false && …` | green; + M1 survives | green; + M1 green (5 dials at load 92) |
| TM2 | `<` → `>=` in the count gate | red on a correct tree | red, `21 dials …` |
| TM3a | `if short > 0 {` → `if false && …`, + `/10` | `/10` survives | T12 green; T7 `reset-holder` timed out 8 s at load 136 (load) |
| TM3b | `if gap < execProbePollEvery {` → `>=` | red on a correct tree | red, `20 of 20 …` |
| TM4 | `const revCadencePeriods = 20` → `= 1`, + M1 | M1 survives (window shortened) | green, 1 dial ≥ gate 1 |
| TM5 | gate `/4+1` → `/10`, + M1 | M1 survives | green, 4 dials ≥ 2 |
| TM6 | the reference's `<-time.After(execProbePollEvery)` → `(500 * time.Millisecond)` | green; + M1 survives | green (122 dials); + M1 green (21 dials) |
| TM7 | `for range revCadencePeriods {` → `for range 0 {` | red on a correct tree | red, `1 dials … (42ns)` |
| TM8 | cap `4*revCadencePeriods` → `1`, + `/10` | `/10` survives (no gaps recorded) | T12 green; T7 `silent-after-chatty` red at load 109 |
| TM9 | `setAt` not recorded | red on a correct tree | red, `set+2562047h…` |
| TM10 | `setAt = t` (the deadline itself), + wait 960 | 960 survives | green; green |
| TM11 | the T10 check against `time.Now()` instead of `setAt`, + wait 951 | red | red, `set+950.999959ms` |
| S4-ceil | `ceiling := listenWindow - execProbePollEvery` → `listenWindow`, alone / + wait 960 | green / green | green / green (and wait 960 alone red) |
| S4-ceil2 | `… - 2*execProbePollEvery` | green (stricter) | green |
| S4-floor | `execProbeHelloFloor` 100 (alone, + wait 249), 500, 600 | green, green, green, red | as expected |
| S4-arrival | `arrival: before`, alone / + wait 5 | green / green | as expected |
| S4-wrote | drop `conn.wrote != execProbeHello`, alone / + P-nohello | green / **red (T7)** | as expected |
| S4-sets | drop `conn.sets == 0`, alone / + P-nodeadline; also drop `IsZero` + P-nodeadline | green / red / green | as expected |

**Controls.** C0 (W, unmutated, start of round, load 92) was rc 1. The FAIL list was
`ForeignListenerOnTheReversePortDoesNotRefuse/chatty` and `RealReachSurvivesAForeignAttempt`
(`broker: handler subprocess timed out after 8s`). T12 and T10 passed. C1 (W, unmutated, end of round, load
76) was rc 0 with 12 PASS. Both failing tests are unchanged by this diff (§5). The earlier v3 round's C0 (W,
load 105) was rc 0 (`designer-mutants-A-v3pre.txt`).

**Headline.**
- **29 production mutants.** 18 killed. 8 equivalent or in-bounds and green as required. 3 survive: the
  2×–4× accepted band.
- **29 test-code mutants and pairs.** 8 red as expected. 21 green as expected: test weakenings, each paired
  to show the clause it removes is load-bearing.

## 4. Milestone M1 (test-only) and acceptance

M1 applies `designer-proto_test.go.patch.txt` to `host/broker/exec_probe_test.go` and adds the CI name.
Export `AILANG_BIN=$HOME/.pinned-ailang/ailang` for every `host/broker` test.

1. `go test -count=1 -v -run '^(TestExecProbeReversePollCadence|TestExecProbeReverseHelloWaitBounds)$' ./host/broker/`
   → rc 0, with `--- PASS: TestExecProbeReversePollCadence` and `--- PASS: TestExecProbeReverseHelloWaitBounds`
   (2 subtests), a `CADENCE dials=` log line, and 0 FAIL/SKIP.
2. `go test -count=1 ./host/broker/` → rc 0: 336 top-level PASS (the judge's 335 + T12), 0 FAIL, 30 top-level SKIP (V15).
3. These §3 rows reproduce, on a scratch copy: M1, M1-div10, M1-ctl50+const100, PC, K-wait951, K-wait950,
   TM2 and TM9. M1 and M1-div10 red only T12 under regex F. PC reds T10 and not T12. K-wait950 stays green.
4. **Load gate.** Under 32 spinning `yes` (killed by a trap, 0 alive after):
   `go test -count=100 -race -run '^(TestExecProbeReversePollCadence|TestExecProbeReverseHelloWaitBounds)$' ./host/broker/`
   → rc 0. Record the CADENCE minimum dial count; it must be ≥ 6 (by construction it is, if rc 0).
5. `gofmt -l host/broker` prints nothing; `go vet ./host/broker/` → rc 0.
6. `go build ./... && go test ./...` → rc 0; `./scripts/verify_ail.sh` → rc 0.
7. Diff shape:
   - `git diff --name-only dcb2366e -- . ':!design_docs'` lists exactly `.github/workflows/ci.yml` and
     `host/broker/exec_probe_test.go`.
   - `git diff --quiet dcb2366e -- host/broker/exec_probe.go` → rc 0.
   - The `ci.yml` diff is one changed line (`:247`) that appends ` TestExecProbeReversePollCadence` after
     `TestExecProbeReverseHelloWaitBounds`.

**The CI list edit.** `ci.yml:247` belongs to the step "Row 140 Workspace.Exec real-srt tests, verbose (linux
bubblewrap)" in the `go host build + test gate` job. That step greps a `--- PASS: <name>` line for each
listed name.
- The name is added there, as row 230 added T10 and T11.
- The race job (`go-race-verify`) needs no edit. Its `VERIFY_GO_LEGS: race` runs `verify_go.sh`'s
  `go test ./... -count=1 -race` (V02), which already includes T12.

## 5. Conflict surface and residuals

- **Files.** One test file and one CI line. No production `.go`, no `.ail`, `go.mod`, `tools/launchd/*` or
  docs change.
- **Concurrent work.** 0 open PRs (V02; control: merged `#466`, `#465`). Row 242 will restructure
  `ProbeExecProfilesEach`/`probe()`, not `pollLoopback` or these tests, so no textual conflict is expected.
- **The count gate's residual.** It is a measured margin, not a proof: thread-level starvation of the poll
  alone could exceed a 4× lag. Worst measured lag: 2.00 (v2); v3 1.18. If T12 ever reds with `N dials while a
  reference …` on an unmutated tree:
  - treat it as this residual;
  - record the CADENCE line and the load;
  - raise `revCadencePeriods`, not the gate ratio.
- **Load erodes the M1 kill** (§2.3). This degrades toward survival only.
- **Time.** T12 costs about 1.05 s idle (20 × 50 ms) and up to the reference's span under load.
- **Unchanged real-socket tests are load-sensitive.** At load 92–142, `ForeignListener…/chatty`,
  `RealReachSurvivesAForeignAttempt`, `ExhaustionRefusesDistinctly/chatty-every-time` and T7's rows timed out
  (`handler subprocess timed out after 8s`) in mutant runs whose mutants cannot reach them, and in C0. This
  diff does not touch them, and C1 passed at load 76. The base was not run at load 92, so "pre-existing" is
  an argument, not a measurement: the failing tests and the production code they run are byte-identical to
  the base. This is outside the row. It is recorded for the controller to triage.

## 6. Verification Log (base `dcb2366e`; files under `design_docs/verification/world-iter339/`)

| V | Command | Observed |
|---|---|---|
| V01 | `git rev-parse HEAD origin/dev`; `go version`; `$HOME/.pinned-ailang/ailang --version`; `sysctl -n hw.ncpu` (`designer-sites.txt`) | `dcb2366e…` twice; `go1.26.6 darwin/arm64`; `AILANG v0.41.0`; `16` |
| V02 | `grep -n -F` of every site; `grep -c -F 'time.After('`; `grep -rn execProbePollEvery --include=*.go .`; the seam's users; T10 sites; `t.Parallel` count; CI names and the race leg; `gh pr list --state open` with merged control (`designer-sites.txt`) | `exec_probe.go` `108` listen, `128` wait, `131` poll, `613` SILENT text, `693` SetDeadline, `694` hello write, `735` 200 ms dial, `763` the timer. `time.After(` count `1`. The constant is read at `:763` and `exec_probe_test.go:896` only (plus comments). `pollLoopback(` callers `:598`, tests `:829`, `:908`. T10 `:851`, `:883`, `:892`, `:896`, `:909`, `:910`, `:913`. 0 files with `t.Parallel`, against 21 `func Test` (control). `ci.yml:247` lists `TestExecProbeReverseHelloWaitBounds`; `:501` `VERIFY_GO_LEGS: race`; `verify_go.sh:355` `go test ./... -count=1 -race`. Open PRs: none; control `#466`, `#465` |
| V03 | controller premise on a `git archive` of `dcb2366e` (`controller-premise.sh.txt` / `.txt`), **controller-measured, not re-run by the designer** | control `ok 10.965s` rc 0; M1 `ok 10.856s` rc 0; constant 2 s rc 1 `hello deadline at stop+499.999166ms … no later than stop+-1s` |
| V04 | final file: `git show HEAD:host/broker/exec_probe_test.go`, then `patch` with `designer-proto_test.go.patch.txt`; `shasum -a 256`; `cmp` of `exec_probe.go` with HEAD; then on the final tree `gofmt -l host/broker`, `go vet ./host/broker/`, and `go test -count=2 -v -run '^(TestExecProbeReversePollCadence\|TestExecProbeReverseHelloWaitBounds)$'` (`designer-final-check.txt`) | patch rc 0; sha256 `bbfccede19236d17…` equals `proto`'s; patch `+59 −4`; `exec_probe.go identical to HEAD`; gofmt list empty; vet rc 0; `--- PASS` T10 (both subtests) ×2 and T12 ×2 (`CADENCE dials=21 span=1.03s fastest=50.16ms`), `ok 2.430s` |
| V05 | mutation round A (production cadence mutants, controls, in-bounds constants; `designer-mutants-A.txt`, `designer-mutants.sh.txt`) | as §3. 16 entries (A) and 6 (borderline, `-count=30`), every one `restored(cmp ok)`; no `SITE ERROR` |
| V06 | mutation round B (survivors 2 and 4; `designer-mutants-B.txt`) | as §2.5 items 2 and 4 and §3. 23 entries, every one `restored(cmp ok)` |
| V07 | survivor 3: wait 500/950/951/952/955/960 on `base` (T10 stop-referenced) and `proto` (setAt), `-count=200` each, at 0 (ambient 41–52), 32 and 160 `yes` (`designer-s3.sh.txt`, `designer-s3-load{0,32,160}.txt`); the worst lag per tree and load (`designer-s3-aggregate.txt`) | Every load, both trees: 500 and 950 green 400/400; 951, 952, 955 and 960 red 400/400. Worst lag of the reference behind the classifier's clock read: base 80.0 / 110.1 / **1,713.3 µs** (0 / 32 / 160 `yes`; the 1.71 ms is on a 955 ms run, reported `stop+953.2867ms`); proposed 0.3 / 0.3 / 0.5 µs. `yes-alive-after=0` each time |
| V08 | borderline, T12 alone `-count=30` (`designer-mutants-borderline.txt`) | 2×: 30/30 green (11 dials); 3×: 30/30 green (8–9); 4×: 30/30 green (6); 5×: 30/30 red (5 dials); 45 ms: 30/30 red (`22 of 22 …`); constant 1 ms: 30/30 green |
| V09 | T12 v1 load matrix, `-count` 300/200 (`designer-v1-load.txt`, `designer-v1-load_all.sh.txt`; v1's summary parser was the earlier inline version of `designer-run_load.sh.txt`) | 1,500 runs, 0 FAIL: idle, 32 `yes` (end load 35.8), `-race`, `-race` + 32 `yes` (end 41.6), 64 `yes` (end 76.1), `-race` + 64 `yes` + `GOMAXPROCS=2` (end 77.2). Fastest gap 50.032–52.080 ms; slowest single gap 85–390 ms |
| V10 | T12 v3 (final) load matrix, T12 and T10 together (`designer-v3-load.txt`, `designer-phase6.sh.txt`, `designer-run_load.sh.txt`, `designer-cadence_summary.py.txt`); worst lag across v2 and v3 (`designer-worst_ratio.py.txt`, appended to the same file) | 800 runs, rc 0 each, 0 FAIL: ambient (200; load 24→10), `-race` + 32 `yes` (200; →53), `-race` + 64 `yes` + `GOMAXPROCS=2` (200; →87), `-race` + 160 `yes` + `GOMAXPROCS=2` (100; →187), 160 `yes` (100; →194). Min dials 20/18/20/20/20; floor min 50.013 ms; slowest gap 1,291 ms; reference span max 6.14 s; worst lag 1.18 (v3) and 2.00 (v2, ambient 150–214) |
| V11 | T12 v1 at ambient load 150–170 + 32 `yes`, `-race`, `GOMAXPROCS=2`, `-count=200` (`designer-v1-ambient-false-red.txt`) | T12 v1 124 PASS / **76 FAIL** (35 `1 dial(s) in one listen window`, 41 `fastest of k dial gaps ≥ 100ms`); T10 (setAt) 200 PASS / 0 FAIL |
| V12 | T12 v2 load matrix (`designer-v2-load.txt`, `designer-v2-phase3.sh.txt`) | idle 300 (min 10 dials of 11); 32 `yes` 300 (min 8; span max 9.2 s); `-race` 300 (min 9); `-race` + 32 `yes` at ambient 150–214: 78 complete runs, 0 FAIL, **min 6 = gate**, then the 10-minute binary timeout (span 32.6 s; one run span 62.6 s with a 29.1 s gap still saw 11); `-race` + 64 `yes` + `GOMAXPROCS=2` at ambient ~206: 200 runs, 0 FAIL, min 9, `ok 486.787s` |
| V13 | mutation round C (test-code mutants; `designer-mutants-C.txt`) | as §3. 15 entries, every one `restored(cmp ok)` |
| V14 | the earlier v3 round with per-gap floor errors (`designer-mutants-A-v3pre.txt`, `designer-mutants-C-v3pre.txt`) | same verdicts as V05/V13 where both ran (C0 at load 105 rc 0). M1-busy wrote a 1.86 GB log (15,555,776 floor lines) before the round was stopped. This motivated the single floor message and the record cap. The file ends mid-entry |
| V15 | `go test ./host/broker/ -count=1 -v` on the final `proto` tree (`designer-fullpkg.txt`) | rc 0, `ok … 200.153s`, 336 top-level `--- PASS` (1,094 at all levels), 0 FAIL, 30 top-level SKIP; T12 `(1.04s)`. Instrument control: the same run on `mutB`, a copy made without `design_docs/`, was rc 1 with only `TestExecArgsSketchDrift` (`open ../../design_docs/sketches/execargs.ail: no such file`). That is a scratch-shape artifact, gone on the full copy |

## 7. Quorum log

- **r1, 2026-11-03T10:34:29Z: PROCEED, 3 of 3 seats present** ($0.3840; author `claude:claude-opus-5-5`, so
  the Claude seat sat out; seats `gemini-4-argon`, `oc-glm-5-3`, `oc-kimi-k3`; `absent_reviewers` `[]`;
  artifact `~/.ailang/state/mission-quorum-world/w-exec-probe-hello-wait-pin-survivors-2026-11-03T10-34-29Z.json`).
  All three passed. `oc-glm-5-3` and `oc-kimi-k3` both noted, without blocking, that the shipped T12 (v3)
  was never run under the 150–214 ambient load that false-redded v1, and that §4 step 4's 32-`yes` load
  gate is far below it. Controller disposition: the plan adds a heavy-load acceptance run of the shipped T12
  (§4 step 4 at 160 `yes`, with `-race`), measured outside the sandbox before the judge, and records its
  minimum dial count against the gate of 6.
