# worldd day-1 benchmark baseline

The recorded day-1 kernel performance budget the charter demands (design doc
`w-worldd-m2` Decision 6 / P3). Later sprints diff against this file on the same
dev rig; CI asserts only that the harness RUNS and REPORTS
(`scripts/bench_worldd.sh --smoke`), never a threshold on shared-runner timings —
noise-gating would be a dishonest gate (S6).

Recorded 2026-07-27 on the development rig. **This table covers the FULL v1
surface**: M2.A recorded three rows and left three PENDING because the routes did
not exist; M2.B landed those routes and measured all six; **M2.C RE-MEASURED all
six on the closing branch**, so every row below comes from the same run rather
than being stitched across milestones.

- Machine: Mac Studio (Mac16,9), Apple M4 Max, 16 cores (12P/4E), 128 GB RAM
- Platform: `darwin/arm64`
- Go: `go version go1.26.4 darwin/arm64`
- Repository commit: the M2.C branch (`sprint/w-worldd-m2-c`, based on `1960188`)
- AILANG pin: `/tmp/ailang-v0300/ailang`, v0.30.0, commit `e37b370`
- Invocation: `go test -bench . -benchtime 200x -run '^$' ./host/daemon/`

Percentiles are computed by the harness itself from PER-ITERATION wall-clock
samples and emitted via `b.ReportMetric` as `p50_ms` / `p95_ms`; `ns/op` is Go's
own mean and is reported alongside for reference. A mean hides the tail the
budget is actually about, which is why the percentiles are the recorded numbers.

| Operation | Day-1 p95 target | Measured p50 | Measured p95 | Result |
|---|---:|---:|---:|---|
| Store commit (embedded `store.Commit`, the kernel floor) | ≤ 25 ms | 0.4026 ms | 0.4717 ms | **inside budget** (53× headroom) |
| REST commit (`POST /v1/commit`) | ≤ 35 ms | 0.4927 ms | 0.5682 ms | **inside budget** (61× headroom) |
| Head read (`GET /v1/head`) | ≤ 5 ms | 0.06504 ms | 0.08017 ms | **inside budget** (62× headroom) |
| Health (`GET /v1/health`) | ≤ 2 ms | 0.04517 ms | 0.06483 ms | **inside budget** (30× headroom) |
| Log range (`GET /v1/log`, limit=100 — the default page) | ≤ 30 ms | 0.9935 ms | 1.205 ms | **inside budget** (24× headroom) |
| Log range (`GET /v1/log`, limit=500 — the clamp max) | ≤ 120 ms | 4.784 ms | 4.996 ms | **inside budget** (24× headroom) |

## What the numbers say about the deliberate N+1

`GET /v1/log` is the surface's **only deliberate N+1**: a bounded loop over the
existing `GetLogEntry` rather than a new range query in the kernel (Decision 3).
The budget exists to measure that trade rather than hide it, so it is worth
reading the two log-range rows against each other:

- 100 entries → 1.205 ms p95; 500 entries → 4.996 ms p95. That is **4.15× the
  time for 5× the rows**, i.e. very close to linear with a small fixed transport
  cost amortised across the larger page (~10.0 µs/entry at 500 vs ~12.1 µs at 100).
  M2.B's independent run measured 3.94× with the same per-entry shape, so this is
  a **reproduced** result, not a single-run artifact.
- Linear-with-no-surprises is exactly the shape a bounded loop should have, and at
  the clamp max it sits **24× inside** its target. **No range-read store method is
  justified by this data.** The signal that would justify one is superlinearity or
  a p95 approaching the target at limit=500; neither is present, in either run. If
  a future sprint proposes a kernel range query, these rows are the evidence it
  must overturn.
- The REST commit costs **0.5682 ms p95 against the embedded floor's 0.4717 ms**.
  **Correcting an over-precise claim in the M2.B entry**: that run read the tax as
  "~0.03 ms", but across the two runs the same difference measured 0.03 ms and
  0.10 ms — the gap moved because the *floor* moved (store commit p95 0.5421 →
  0.4717 ms) while the REST row barely did (0.5763 → 0.5682 ms). At these
  sub-millisecond magnitudes a single run cannot resolve the tax to two decimal
  places. The claim the data actually supports is the weaker, durable one:
  **the loopback transport adds well under 0.1 ms, so commit cost is dominated by
  the kernel's fsync, not the daemon.**

## Raw benchmark evidence

```text
goos: darwin
goarch: arm64
pkg: github.com/sunholo-data/ailang-world/host/daemon
cpu: Apple M4 Max
BenchmarkStoreCommit-16    	     200	    408446 ns/op	         0.4026 p50_ms	         0.4717 p95_ms
BenchmarkHeadRead-16       	     200	     65850 ns/op	         0.06504 p50_ms	         0.08017 p95_ms
BenchmarkHealth-16         	     200	     47890 ns/op	         0.04517 p50_ms	         0.06483 p95_ms
BenchmarkRESTCommit-16     	     200	    495419 ns/op	         0.4927 p50_ms	         0.5682 p95_ms
BenchmarkLogRange/limit_100-16         	     200	   1018311 ns/op	         0.9935 p50_ms	         1.205 p95_ms
BenchmarkLogRange/limit_500-16         	     200	   4785594 ns/op	         4.784 p50_ms	         4.996 p95_ms
PASS
ok  	github.com/sunholo-data/ailang-world/host/daemon	2.605s
```

## Provenance and honesty notes

- `BenchmarkStoreCommit` opens a **fresh temp-file** SQLite store per run — never
  `:memory:` — so the number includes fsync reality and also exercises A1's
  fail-closed writer lock. Each iteration commits a DISTINCT world, chained on
  the previous head, so the compare-and-append path does real work every time.
- The four HTTP benchmarks measure a **real loopback round-trip** against a
  running daemon (ephemeral port, keep-alive connection warmed OUTSIDE the
  measured region), not an in-process handler call. `BenchmarkRESTCommit` chains
  each iteration on the previous head, so every request is a real
  compare-and-append rather than a repeated conflict.
- **Provenance is split, and the split is stated rather than smoothed over — for
  the third milestone running.** The sprint executor (codex `gpt-5.6-sol`) works
  in a sandbox that denies loopback `bind(2)`, so it can author the socket
  harnesses but cannot execute ANY of them. In M2.A, M2.B and again in M2.C it
  recorded those rows as UNAVAILABLE, quoted the sandbox error verbatim, and
  explicitly declined to invent values. That refusal is the only reason this table
  is trustworthy: a fabricated p95 here would poison every future sprint that
  diffs against this file, and it would be undetectable after the fact.
  **The controller measured all six rows on the dev rig outside that sandbox, in
  the single 200x invocation quoted verbatim above.**
- **Every row is re-measured each milestone, never carried forward.** Rows drift
  run to run on a shared machine (store commit p95 has read 0.6093 → 0.5421 →
  0.4717 ms across A3/M2.B/M2.C without any change to the commit path). Carrying a
  row forward would silently mix invocations and make exactly that drift look like
  a regression or an improvement. Re-measuring the whole table together is what
  makes the rows comparable to each other; comparing ACROSS milestones is only
  ever indicative, which is why the N+1 analysis above leans on the row *ratio*
  rather than on absolute movement.
- Every row is comfortably inside its day-1 target, in both the M2.B and M2.C
  runs. **No row constitutes a design signal for a range-read store method** — see
  the N+1 analysis above.
