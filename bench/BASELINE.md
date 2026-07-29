# worldd day-1 benchmark baseline

The recorded day-1 kernel and broker performance budget (`w-worldd-m2`
Decision 6 / P3, `w-store-durability` Decision 7, and `w-effect-broker-m3`).
Later sprints diff against this file on the same development rig. CI asserts
only that the harness runs and reports (`scripts/bench_worldd.sh --smoke`);
thresholds are evaluated from the development-rig baseline. Noise-gating a
shared runner would be a dishonest gate (S6).

**All ten rows below come from one invocation.** The sprint executor's sandbox
denies the loopback binds five of the rows need, so the executor wrote
`<CONTROLLER-MEASURED>` into every measured cell and quoted the denied bind
verbatim; the controller performed the single complete measurement outside the
sandbox and replaced every row together.

- Machine: Mac Studio (Mac16,9), Apple M4 Max, 16 cores (12P/4E), 128 GB RAM
- Platform: `darwin/arm64`
- Go: `go version go1.26.4 darwin/arm64`
- Repository commit: the MJ.C branch (`sprint/w-effect-journal-mjc`, based on `b485ead`)
- AILANG pin: `/tmp/ailang-v0300/ailang`, v0.30.0, commit `e37b370`
- Invocation: `go test -bench . -benchtime 200x -run '^$' ./host/daemon/`
- **Rig load at measurement: `load averages: 5.22 4.99 5.91`** — the sibling V1
  mission's eval suite was running `ollama` + `llama-server` at 80–98% CPU
  throughout. See the MJ.C section: the absolute rows below are load-elevated,
  and this item's cost was therefore established by a same-rig A/B, not by a
  diff against the idle-rig M3.C numbers.

Percentiles are computed from per-iteration wall-clock samples and emitted as
`p50_ms` and `p95_ms`; `ns/op` remains in the raw output for reference. A mean
hides the tail the budget is actually about, which is why the percentiles are
the recorded numbers.

## MJ.C re-measure — and the delta that would have been a fabricated regression

MJ.C invoked all ten rows together with the command above. The sandbox denied
the five loopback benchmarks with:
`listen tcp 127.0.0.1:0: bind: operation not permitted`.
That run is **UNINFORMATIVE UNDER SANDBOX**, and the executor correctly wrote
`<CONTROLLER-MEASURED>` into every measured cell rather than mixing partial
numbers into a complete baseline. The controller performed the single complete
unsandboxed invocation and replaced every row together.

**Then the naive delta turned out to be an artefact of the rig, not of the code.**
Read against the idle-rig M3.C row, `BenchmarkBrokerFSRead` p95 moved
0.7472 ms → 4.529 ms, a **6.06× "regression"** that this item would have banked
as the effect journal's cost. It is not. The sibling **V1 mission's eval suite
was running on the same development rig** (`ollama` + `llama-server`, 80–98%
CPU, `load averages: 5.22 4.99 5.91`), and every row is inflated by it.

The cost was therefore established by a **same-rig A/B** — the identical
invocation on the pre-MJ.C parent commit `b485ead`, under the same load, minutes
apart — rather than by a diff against numbers taken under different conditions:

| `BenchmarkBrokerFSRead` p95 | Value |
|---|---:|
| Control — `b485ead`, pre-MJ.C, same loaded rig | **4.523 ms** |
| MJ.C run 1 (the recorded all-row invocation) | 4.529 ms |
| MJ.C run 2 | 4.610 ms |
| MJ.C run 3 | 4.604 ms |
| **MJ.C cost vs control** | **+0.13% — within run-to-run noise; no measurable cost** |

**MJ.C's measured cost is zero.** That is the expected result: MJ.C adds no
production code to the dispatch path — the effect-journal appends this row now
includes were already there from **MJ.B**, and MJ.B's own cost was never
isolated against a same-rig control.

The generalisable point, and the reason the A/B was run at all: **a delta
against a baseline captured under different conditions is not a measurement of
the change — it is a measurement of the conditions.** A benchmark is an
instrument, and it inherits the same burden of proof as any other instrument in
this repository; the control is what distinguishes a real regression from a busy
machine. Nothing in this file, in `scripts/bench_worldd.sh`, or in the harness
records the rig load at measurement time, so this confound is currently invisible
by default and was caught only because the ratio looked implausible. Raised as
queue item **4f `w-bench-load-confound`**.

The three MJ.C runs above also satisfy this file's own "a ratio within 2× of 1.0
is run three times with all three reported" rule for the commit-with-receipt row
(+50.5%, +72.1%, +46.1%; control +35.1%) — a spread that is itself the load, and
a standing exceedance now in its sixth milestone, not a new one.

| Operation | Day-1 p95 target | Measured p50 | Measured p95 | Result |
|---|---:|---:|---:|---|
| Store commit (embedded `store.Commit`, kernel floor) | ≤ 25 ms | 0.4808 ms | 1.015 ms | **inside budget** (24.6× headroom) |
| Journal intent append (`store.AppendIntent`) | ≤ 10 ms | 0.7145 ms | 0.9647 ms | **inside budget** (10.4× headroom) |
| Commit with in-transaction receipt | ≤ 120% of store-commit p95 (≤ 1.218 ms) | 1.019 ms | 1.528 ms | **TARGET EXCEEDED — +50.5%**, see below |
| REST commit (`POST /v1/commit`) | ≤ 35 ms | 0.6105 ms | 1.297 ms | **inside budget** (27× headroom) |
| Head read (`GET /v1/head`) | ≤ 5 ms | 0.09896 ms | 0.1449 ms | **inside budget** (34.5× headroom) |
| Health (`GET /v1/health`) | ≤ 2 ms | 0.06546 ms | 0.08962 ms | **inside budget** (22.3× headroom) |
| Log range (`GET /v1/log`, limit=100 — the default page) | ≤ 30 ms | 3.064 ms | 3.850 ms | **inside budget** (7.8× headroom) |
| Log range (`GET /v1/log`, limit=500 — the clamp max) | ≤ 120 ms | 15.65 ms | 17.55 ms | **inside budget** (6.8× headroom) |
| Pure broker decision (`broker.Decide`) | ≤ 0.1 ms | 0.0000830 ms | 0.0000840 ms | **inside budget, but RESOLUTION-LIMITED — read the note below, not the number** |
| Brokered `FS.Read` full pipeline | ≤ 10 ms | 3.618 ms | 4.529 ms | **inside budget** (2.2× headroom) — load-elevated; see the A/B above |

**Every absolute row above is load-elevated** and must not be read as a
regression against the idle-rig M3.C table. The headroom figures are honest for
*this* run's conditions and are the conservative reading; the ratios and the A/B
are the load-independent signal.

## The pure decision is below this harness's clock resolution — a bound, not a measurement

`BenchmarkBrokerDecide` reports p50 **and** p95 as **exactly** `0.0000420 ms` in
all three runs. Two numbers that agree to three significant figures across three
independent runs are not a tail measurement; they are a **quantization
artifact**. 0.0000420 ms = **42 ns**, and darwin/arm64's `mach_absolute_time`
timebase ticks at **41.67 ns** — so every per-sample `time.Since` reading is
exactly ONE tick, and the percentile over a constant is that constant.

The resolvable number is Go's own aggregate: **78.55 ns/op** (77.50 / 81.45 in
the other two runs), which averages over 200 iterations and therefore does not
hit the tick floor — though it also includes the harness's own `time.Now()` and
`append` per iteration, so the decision itself is somewhat cheaper still.

**The honest claim: the pure capability/budget decision costs on the order of
80 ns, and its p95 cannot be resolved by this harness — it is at or below one
clock tick.** Against the 0.1 ms target that is roughly **1,200× inside
budget**, so the budget question is not close; but recording `0.0000420 ms p95`
as if it were a measured tail would be the same class of error as quoting a
ratio from a 50-sample run. The percentile harness is fine for the
sub-millisecond store rows it was built for and is simply the wrong instrument
three orders of magnitude further down. If a future sprint needs a real tail for
this row, it must batch N decisions per sample rather than time one.

## The receipt tax — a target blown by ~50%, recorded rather than relaxed

Decision 7 set the receipt target at **within +20%** of a bare store commit. The
measured tax is **+47.5%** (0.4610 → 0.6800 ms p95). The number is recorded as a
design signal, exactly as Decision 7 requires; **no threshold was relaxed and no
row was re-run until it agreed.**

The pair is fixture-matched on purpose: both use temp-file SQLite stores and
chained, distinct commits, and both time ONLY the `store.Commit` call.
`BenchmarkCommitWithReceipt` stages its durable intents *before* the measured
region, so the number is the marginal **in-transaction** cost — one indexed
journal lookup, the eight-field intent compare, the outcome encode, and two
extra row inserts.

Ratios close to unity are unstable at low sample counts, so the comparison is
run three times at `-benchtime 200x` and all three are reported:

| Run | store commit p95 | commit-with-receipt p95 | ratio |
|---|---:|---:|---:|
| 1 (the recorded invocation) | 0.4610 ms | 0.6800 ms | **1.475×** |
| 2 | 0.4520 ms | 0.6823 ms | 1.510× |
| 3 | 0.4634 ms | 0.7043 ms | 1.520× |

This is the **second independent reproduction** of the tax: `w-store-durability`
SD.C measured 1.51× / 1.49× / 1.46× on the same rig, and M3.C measures 1.475× /
1.510× / 1.520×. Six runs across two milestones put it in a **1.46×–1.52×** band
with no overlap with the +20% target. The target IS blown, and it is blown
reproducibly rather than noisily — which is precisely why it is recorded instead
of adjusted.

## Per-commit receipt tax amortized across brokered effects

This section answers the question SD.C's baseline explicitly handed to
`w-effect-broker-m3`. The receipt tax is paid once per **commit**, not once per
effect. Under every M3.D option, each brokered effect pays the two
content-addressed `PutObject` writes priced by `BenchmarkBrokerFSRead`; the
in-transaction receipt is paid once at the episode's commit boundary.

**Both figures below are the M3.C IDLE-rig numbers, and deliberately so** — this
analysis is a ratio between two rows, so it is only valid when both rows come
from the same conditions. MJ.C's rows are load-elevated (see the MJ.C section),
so substituting them here would change the arithmetic without changing the
conclusion. Re-derive this section from a clean-rig invocation when queue item
**4f `w-bench-load-confound`** lands a load gate.

- Measured receipt delta: 0.6800 − 0.4610 = **0.2190 ms** per commit (p95, run 1).
- Measured per-effect cost: **0.7472 ms** p95 (`BenchmarkBrokerFSRead`).
- N=1 brokered effect: 0.2190 / 1 = **0.2190 ms** per effect → **+29.3%** on top
  of that effect's own cost.
- N=3 brokered effects: 0.2190 / 3 = **0.0730 ms** per effect → **+9.8%**.
- N=6 brokered effects — **the acceptance episode's actual effect count**, which
  is 6, not the 3 the sprint plan predicted: 0.2190 / 6 = **0.0365 ms** per
  effect → **+4.9%**.

**What the data supports:** the receipt tax is a *per-episode* constant, and an
episode doing real work amortizes it into the noise — at the acceptance
episode's own N=6 it is under 5% of one effect's cost, against a per-effect
pipeline that already sits 13× inside its own budget. Decision 7's +20% bound
was written against a *per-commit* comparison and is genuinely exceeded there;
whether +20% was ever the right shape of bound for two extra indexed inserts
inside an existing transaction is a **design signal for a future document change
with its own evidence, NOT a threshold this sprint edits.** This table is the
evidence that discussion starts from.

## What the numbers say about the deliberate N+1

`GET /v1/log` is the surface's **only deliberate N+1**: a bounded loop over the
existing `GetLogEntry` rather than a new range query in the kernel
(`w-worldd-m2` Decision 3).

- 100 entries → 1.409 ms p95; 500 entries → 6.149 ms p95. That is **4.36× the
  time for 5× the rows**, i.e. close to linear with a small fixed transport cost
  amortised across the larger page (~12.3 µs/entry at 500 vs ~14.1 µs at 100).
  M2.B measured 3.94×, M2.C 4.15×, SD.C 4.34× and this run 4.36× — a
  **four-times-reproduced** result with the same per-entry shape.
- At the clamp max it sits **20× inside** its target. **No range-read store
  method is justified by this data.** The signal that would justify one is
  superlinearity or a p95 approaching the target at limit=500; neither is present
  in any of the four runs.
- **The transport claim is falsified again, by its own next sample — for the
  second milestone running.** M2.C concluded "the loopback transport adds well
  under 0.1 ms"; SD.C falsified that at 0.136 ms and replaced it with "on the
  order of 0.1 ms and never more than ~0.15 ms". This run measures
  **0.1709 ms** (0.4610 → 0.6319 ms p95), which falsifies the replacement too.
  Across four samples the delta reads 0.03 / 0.10 / 0.136 / **0.171** ms — it
  has risen monotonically every time it has been measured, and the "never more
  than X" form has now failed twice. The durable claim is the one that does not
  assert a ceiling it has not earned: **loopback transport adds on the order of
  0.1–0.2 ms, and commit cost remains dominated by the kernel's fsync rather
  than by the daemon.** A future milestone that wants a ceiling here must
  measure enough samples to justify one, rather than reading the maximum of its
  own three or four points as a bound.

## Raw benchmark evidence

**MJ.C — the recorded all-row invocation (LOADED rig, `load averages: 5.22 4.99 5.91`).**
These are the numbers the summary table above reports.

```text
goos: darwin
goarch: arm64
pkg: github.com/sunholo-data/ailang-world/host/daemon
cpu: Apple M4 Max
BenchmarkStoreCommit-16          	     200	    619414 ns/op	         0.4808 p50_ms	         1.015 p95_ms
BenchmarkJournalAppend-16        	     200	    703076 ns/op	         0.7145 p50_ms	         0.9647 p95_ms
BenchmarkCommitWithReceipt-16    	     200	    944072 ns/op	         1.019 p50_ms	         1.528 p95_ms
BenchmarkHeadRead-16             	     200	    107971 ns/op	         0.09896 p50_ms	         0.1449 p95_ms
BenchmarkHealth-16               	     200	     66019 ns/op	         0.06546 p50_ms	         0.08962 p95_ms
BenchmarkRESTCommit-16           	     200	    803222 ns/op	         0.6105 p50_ms	         1.297 p95_ms
BenchmarkLogRange/limit_100-16   	     200	   2465289 ns/op	         3.064 p50_ms	         3.850 p95_ms
BenchmarkLogRange/limit_500-16   	     200	  12167745 ns/op	        15.65 p50_ms	        17.55 p95_ms
BenchmarkBrokerDecide-16         	     200	       141.2 ns/op	         0.0000830 p50_ms	         0.0000840 p95_ms
BenchmarkBrokerFSRead-16         	     200	   3206511 ns/op	         3.618 p50_ms	         4.529 p95_ms
PASS
ok  	github.com/sunholo-data/ailang-world/host/daemon
```

**The MJ.C A/B control — the identical invocation on the pre-MJ.C parent
`b485ead`, same loaded rig, minutes apart.** This is what establishes that MJ.C
costs nothing; it is not a milestone baseline and must not be diffed as one.

```text
BenchmarkStoreCommit-16          	     200	    610850 ns/op	         0.4409 p50_ms	         1.022 p95_ms
BenchmarkJournalAppend-16        	     200	    613433 ns/op	         0.6074 p50_ms	         0.9161 p95_ms
BenchmarkCommitWithReceipt-16    	     200	    946877 ns/op	         1.034 p50_ms	         1.381 p95_ms
BenchmarkHeadRead-16             	     200	     70529 ns/op	         0.06671 p50_ms	         0.08904 p95_ms
BenchmarkHealth-16               	     200	     67325 ns/op	         0.06692 p50_ms	         0.09046 p95_ms
BenchmarkRESTCommit-16           	     200	    924774 ns/op	         0.9967 p50_ms	         1.316 p95_ms
BenchmarkLogRange/limit_100-16   	     200	   2497106 ns/op	         3.128 p50_ms	         3.749 p95_ms
BenchmarkLogRange/limit_500-16   	     200	  12177548 ns/op	        15.71 p50_ms	        17.31 p95_ms
BenchmarkBrokerDecide-16         	     200	       172.3 ns/op	         0.0000830 p50_ms	         0.0000840 p95_ms
BenchmarkBrokerFSRead-16         	     200	   3296938 ns/op	         3.661 p50_ms	         4.523 p95_ms
```

**M3.C — the last IDLE-rig all-row invocation, retained deliberately.** The
amortisation analysis above is computed from these numbers, and they are the
only clean-rig reference this file still holds. **They are NOT this milestone's
measurement** — diffing the loaded MJ.C rows against them is exactly the
fabricated 6.06× "regression" the MJ.C section documents.

```text
BenchmarkStoreCommit-16          	     200	    399910 ns/op	         0.3980 p50_ms	         0.4610 p95_ms
BenchmarkJournalAppend-16        	     200	    395345 ns/op	         0.3906 p50_ms	         0.4542 p95_ms
BenchmarkCommitWithReceipt-16    	     200	    572273 ns/op	         0.5543 p50_ms	         0.6800 p95_ms
BenchmarkHeadRead-16             	     200	     71929 ns/op	         0.07033 p50_ms	         0.08904 p95_ms
BenchmarkHealth-16               	     200	     50341 ns/op	         0.04646 p50_ms	         0.06813 p95_ms
BenchmarkRESTCommit-16           	     200	    536061 ns/op	         0.5256 p50_ms	         0.6319 p95_ms
BenchmarkLogRange/limit_100-16   	     200	   1183422 ns/op	         1.151 p50_ms	         1.409 p95_ms
BenchmarkLogRange/limit_500-16   	     200	   5770270 ns/op	         5.759 p50_ms	         6.149 p95_ms
BenchmarkBrokerDecide-16         	     200	        78.55 ns/op	         0.0000420 p50_ms	         0.0000420 p95_ms
BenchmarkBrokerFSRead-16         	     200	    639719 ns/op	         0.6368 p50_ms	         0.7472 p95_ms
PASS
ok  	github.com/sunholo-data/ailang-world/host/daemon	3.433s
```

## Provenance and honesty notes

- Every store-using benchmark opens a **fresh temp-file** SQLite store, never
  `:memory:`, so every store-level number includes fsync reality and also
  exercises A1's fail-closed writer lock. Each commit iteration commits a
  DISTINCT world, chained on the previous head, so the compare-and-append path
  does real work every time.
- The five HTTP rows measure a **real loopback round-trip** against a running
  daemon (ephemeral port, keep-alive connection warmed OUTSIDE the measured
  region), not an in-process handler call.
- The broker rows: `BenchmarkBrokerDecide` is pure (no store, no handler);
  `BenchmarkBrokerFSRead` is the FULL pipeline — decision, ledger debit, handler
  dispatch, `PutObject` of the result bytes and `PutObject` of the effect record
  — with distinct file bytes each iteration so no read is served from a warm no-op.
- **Provenance is split, and the split is stated rather than smoothed over — for
  the fifth milestone running** (M2.A, M2.B, M2.C, SD.C, M3.C). The sprint
  executor (codex `gpt-5.6-sol`) works in a sandbox that denies loopback
  `bind(2)`, so it can author the socket harnesses but cannot execute them. It
  wrote `<CONTROLLER-MEASURED>` into all 44 measured fields, quoted
  `listen tcp 127.0.0.1:0: bind: operation not permitted` verbatim, and
  explicitly declined to invent values. That refusal is the only reason this
  table is trustworthy: a fabricated p95 here would poison every future sprint
  that diffs against this file and would be undetectable after the fact.
  **MJ.C repeated the pattern for a sixth milestone** — and added the missing
  half: declining to invent a number protects against fabricated *values*, but
  not against a fabricated *delta*. See the MJ.C A/B above, where a real
  measurement compared against a baseline taken under different rig conditions
  produced a 6.06× regression that does not exist.
- **Every row is re-measured each milestone, never carried forward.** Rows drift
  run to run on a shared machine (store commit p95 has read 0.6093 → 0.5421 →
  0.4717 → 0.4537 → 0.4610 ms across A3/M2.B/M2.C/SD.C/M3.C without any
  regression in the commit path). Carrying a row forward would silently mix
  invocations and make that drift look like a regression. Comparing ACROSS
  milestones is only ever indicative, which is why the analyses above lean on
  row *ratios* rather than on absolute movement.
- **A number at insufficient resolution is a claim, exactly as a ratio at
  insufficient sample count is.** Sub-millisecond rows use at least
  `-benchtime 200x`, a ratio within 2× of 1.0 is run three times with all three
  reported, and a percentile whose p50 equals its p95 to three significant
  figures is reported as a resolution bound rather than a tail (see the
  `BenchmarkBrokerDecide` note above).
- Eight of the ten rows are comfortably inside their day-1 targets. The
  commit-with-receipt row **is not**, and is recorded that way rather than
  re-targeted. The broker decision row is inside by a wide margin but is
  reported as a bound rather than a measurement.
