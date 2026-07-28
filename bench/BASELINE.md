# worldd day-1 benchmark baseline

The recorded day-1 kernel performance budget the charter demands
(`w-worldd-m2` Decision 6 / P3, extended by `w-store-durability` Decision 7).
Later sprints diff against this file on the same development rig. CI asserts
only that the harness runs and reports (`scripts/bench_worldd.sh --smoke`);
thresholds are evaluated from the development-rig baseline. Noise-gating a
shared runner would be a dishonest gate (S6).

**All eight rows below come from one invocation.** The sprint executor's sandbox
denies the loopback binds five of the rows need, so the controller performed the
single complete measurement outside the sandbox and replaced every row together.

- Machine: Mac Studio (Mac16,9), Apple M4 Max, 16 cores (12P/4E), 128 GB RAM
- Platform: `darwin/arm64`
- Go: `go version go1.26.4 darwin/arm64`
- Repository commit: the SD.C branch (`sprint/w-store-durability-sdc`, based on `14635f8`)
- AILANG pin: `/tmp/ailang-v0300/ailang`, v0.30.0, commit `e37b370`
- Invocation: `go test -bench . -benchtime 200x -run '^$' ./host/daemon/`

Percentiles are computed by the harness from per-iteration wall-clock samples
and emitted through `b.ReportMetric` as `p50_ms` and `p95_ms`. `ns/op` remains
in the raw output for reference. A mean hides the tail the budget is actually
about, which is why the percentiles are the recorded numbers.

| Operation | Day-1 p95 target | Measured p50 | Measured p95 | Result |
|---|---:|---:|---:|---|
| Store commit (embedded `store.Commit`, kernel floor) | ≤ 25 ms | 0.3962 ms | 0.4537 ms | **inside budget** (55× headroom) |
| Journal intent append (`store.AppendIntent`) | ≤ 10 ms | 0.3880 ms | 0.4599 ms | **inside budget** (22× headroom) |
| Commit with in-transaction receipt | ≤ 120% of store-commit p95 (≤ 0.5444 ms) | 0.5657 ms | 0.6846 ms | **TARGET EXCEEDED — +50.9%**, see below |
| REST commit (`POST /v1/commit`) | ≤ 35 ms | 0.5051 ms | 0.5900 ms | **inside budget** (59× headroom) |
| Head read (`GET /v1/head`) | ≤ 5 ms | 0.06617 ms | 0.09167 ms | **inside budget** (55× headroom) |
| Health (`GET /v1/health`) | ≤ 2 ms | 0.04667 ms | 0.06492 ms | **inside budget** (31× headroom) |
| Log range (`GET /v1/log`, limit=100 — the default page) | ≤ 30 ms | 1.077 ms | 1.276 ms | **inside budget** (24× headroom) |
| Log range (`GET /v1/log`, limit=500 — the clamp max) | ≤ 120 ms | 5.347 ms | 5.538 ms | **inside budget** (22× headroom) |

## The receipt tax — a target blown by 51%, recorded rather than relaxed

Decision 7 set the receipt target at **within +20%** of a bare store commit. The
measured tax is **+50.9%** (0.4537 → 0.6846 ms p95). The number is recorded as a
design signal, exactly as Decision 7 requires; **no threshold was relaxed and no
row was re-run until it agreed.**

The pair is fixture-matched on purpose so the comparison means something: both
use temp-file SQLite stores and chained, distinct commits, and both time ONLY the
`store.Commit` call. `BenchmarkCommitWithReceipt` stages its durable intents
*before* the measured region, so the number below is the marginal
**in-transaction** cost — one indexed journal lookup, the eight-field intent
compare, the outcome encode, and two extra row inserts — not the cost of the
intent append, which is priced separately in its own row.

**The tax is reproduced, and an earlier reading of it was wrong.** The sprint
executor measured the same pair at `-benchtime 50x` inside its sandbox and
reported **2.8×**. Re-measured on the dev rig at 200x, three times:

| Run | store commit p95 | commit-with-receipt p95 | ratio |
|---|---:|---:|---:|
| 1 (the recorded invocation) | 0.4537 ms | 0.6846 ms | **1.51×** |
| 2 | 0.4446 ms | 0.6629 ms | 1.49× |
| 3 | 0.4545 ms | 0.6638 ms | 1.46× |

So the target IS blown, but by **~50%, not ~180%** — the 50x reading overstated
it by nearly 2×. At these sub-millisecond magnitudes a 50-sample run cannot
resolve the ratio, which is the same lesson the REST-commit row already carried.
The durable claim the data supports: **the in-transaction receipt costs roughly
half a bare commit again, and both remain ~50× inside the absolute kernel
budget.** Whether +20% was ever the right target for two extra indexed inserts
inside an existing transaction is a question for `w-effect-broker-m3`, which is
the first component to pay this cost on a real dispatch path; this table is the
evidence that discussion starts from.

## What the numbers say about the deliberate N+1

`GET /v1/log` is the surface's **only deliberate N+1**: a bounded loop over the
existing `GetLogEntry` rather than a new range query in the kernel
(`w-worldd-m2` Decision 3). The budget exists to measure that trade rather than
hide it, so it is worth reading the two log-range rows against each other:

- 100 entries → 1.276 ms p95; 500 entries → 5.538 ms p95. That is **4.34× the
  time for 5× the rows**, i.e. close to linear with a small fixed transport cost
  amortised across the larger page (~11.1 µs/entry at 500 vs ~12.8 µs at 100).
  M2.B measured 3.94×, M2.C 4.15×, and this run 4.34× with the same per-entry
  shape — a **three-times-reproduced** result, not a single-run artifact.
- Linear-with-no-surprises is exactly the shape a bounded loop should have, and
  at the clamp max it sits **22× inside** its target. **No range-read store
  method is justified by this data.** The signal that would justify one is
  superlinearity or a p95 approaching the target at limit=500; neither is present
  in any of the three runs. If a future sprint proposes a kernel range query,
  these rows are the evidence it must overturn.
- **Correcting the M2.C entry's transport claim.** That entry concluded "the
  loopback transport adds well under 0.1 ms". Across three runs the same
  difference has now measured 0.03 ms, 0.10 ms and — here — **0.136 ms**
  (0.4537 → 0.5900 ms p95), so the "well under 0.1 ms" form is already falsified
  by its own third sample. The gap moves mostly because the *floor* moves. The
  claim the data actually supports is the weaker, durable one: **loopback
  transport adds on the order of 0.1 ms and never more than ~0.15 ms, so commit
  cost is dominated by the kernel's fsync, not by the daemon.**

## Raw benchmark evidence

```text
goos: darwin
goarch: arm64
pkg: github.com/sunholo-data/ailang-world/host/daemon
cpu: Apple M4 Max
BenchmarkStoreCommit-16          	     200	    392958 ns/op	         0.3962 p50_ms	         0.4537 p95_ms
BenchmarkJournalAppend-16        	     200	    394353 ns/op	         0.3880 p50_ms	         0.4599 p95_ms
BenchmarkCommitWithReceipt-16    	     200	    590346 ns/op	         0.5657 p50_ms	         0.6846 p95_ms
BenchmarkHeadRead-16             	     200	     69467 ns/op	         0.06617 p50_ms	         0.09167 p95_ms
BenchmarkHealth-16               	     200	     48321 ns/op	         0.04667 p50_ms	         0.06492 p95_ms
BenchmarkRESTCommit-16           	     200	    511095 ns/op	         0.5051 p50_ms	         0.5900 p95_ms
BenchmarkLogRange/limit_100-16   	     200	   1104855 ns/op	         1.077 p50_ms	         1.276 p95_ms
BenchmarkLogRange/limit_500-16   	     200	   5335346 ns/op	         5.347 p50_ms	         5.538 p95_ms
PASS
ok  	github.com/sunholo-data/ailang-world/host/daemon	3.044s
```

## Provenance and honesty notes

- Every benchmark opens a **fresh temp-file** SQLite store, never `:memory:`, so
  every store-level number includes fsync reality and also exercises A1's
  fail-closed writer lock. Each iteration commits a DISTINCT world, chained on
  the previous head, so the compare-and-append path does real work every time.
- The five HTTP rows measure a **real loopback round-trip** against a running
  daemon (ephemeral port, keep-alive connection warmed OUTSIDE the measured
  region), not an in-process handler call. `BenchmarkRESTCommit` chains each
  iteration on the previous head, so every request is a real compare-and-append
  rather than a repeated conflict.
- **Provenance is split, and the split is stated rather than smoothed over — for
  the fourth milestone running.** The sprint executor (codex `gpt-5.6-sol`) works
  in a sandbox that denies loopback `bind(2)`, so it can author the socket
  harnesses but cannot execute them. In M2.A, M2.B, M2.C and again in SD.C it
  wrote `<CONTROLLER-MEASURED>` into every cell, quoted the sandbox error
  verbatim, and explicitly declined to invent values. That refusal is the only
  reason this table is trustworthy: a fabricated p95 here would poison every
  future sprint that diffs against this file, and it would be undetectable after
  the fact. **The controller measured all eight rows on the dev rig outside that
  sandbox, in the single 200x invocation quoted verbatim above.**
- **Every row is re-measured each milestone, never carried forward.** Rows drift
  run to run on a shared machine (store commit p95 has read 0.6093 → 0.5421 →
  0.4717 → 0.4537 ms across A3/M2.B/M2.C/SD.C without any regression in the
  commit path). Carrying a row forward would silently mix invocations and make
  exactly that drift look like a regression or an improvement. Re-measuring the
  whole table together is what makes the rows comparable **to each other**;
  comparing ACROSS milestones is only ever indicative, which is why both analyses
  above lean on row *ratios* rather than on absolute movement.
- Seven of the eight rows are comfortably inside their day-1 targets. The eighth
  — commit-with-receipt — **is not**, and is recorded that way rather than
  re-targeted.
