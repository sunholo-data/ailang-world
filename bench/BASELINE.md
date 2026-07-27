# worldd day-1 benchmark baseline

The recorded day-1 kernel performance budget the charter demands (design doc
`w-worldd-m2` Decision 6 / P3). Later sprints diff against this file on the same
dev rig; CI asserts only that the harness RUNS and REPORTS
(`scripts/bench_worldd.sh --smoke`), never a threshold on shared-runner timings —
noise-gating would be a dishonest gate (S6).

Recorded 2026-07-27 on the development rig. **This table now covers the FULL v1
surface**: M2.A recorded three rows and left three PENDING because the routes did
not exist; M2.B landed those routes and re-measured all six in ONE invocation, so
every row below comes from the same run rather than being stitched across two.

- Machine: Mac Studio (Mac16,9), Apple M4 Max, 16 cores (12P/4E), 128 GB RAM
- Platform: `darwin/arm64`
- Go: `go version go1.26.4 darwin/arm64`
- Repository commit: the M2.B branch point (`f61aafb`) plus the M2.B REST surface
- AILANG pin: `/tmp/ailang-v0300/ailang`, v0.30.0, commit `e37b370`
- Invocation: `go test -bench . -benchtime 200x -run '^$' ./host/daemon/`

Percentiles are computed by the harness itself from PER-ITERATION wall-clock
samples and emitted via `b.ReportMetric` as `p50_ms` / `p95_ms`; `ns/op` is Go's
own mean and is reported alongside for reference. A mean hides the tail the
budget is actually about, which is why the percentiles are the recorded numbers.

| Operation | Day-1 p95 target | Measured p50 | Measured p95 | Result |
|---|---:|---:|---:|---|
| Store commit (embedded `store.Commit`, the kernel floor) | ≤ 25 ms | 0.4715 ms | 0.5421 ms | **inside budget** (46× headroom) |
| REST commit (`POST /v1/commit`) | ≤ 35 ms | 0.5000 ms | 0.5763 ms | **inside budget** (61× headroom) |
| Head read (`GET /v1/head`) | ≤ 5 ms | 0.06617 ms | 0.08033 ms | **inside budget** (62× headroom) |
| Health (`GET /v1/health`) | ≤ 2 ms | 0.04579 ms | 0.06796 ms | **inside budget** (29× headroom) |
| Log range (`GET /v1/log`, limit=100 — the default page) | ≤ 30 ms | 0.9824 ms | 1.248 ms | **inside budget** (24× headroom) |
| Log range (`GET /v1/log`, limit=500 — the clamp max) | ≤ 120 ms | 4.738 ms | 4.915 ms | **inside budget** (24× headroom) |

## What the numbers say about the deliberate N+1

`GET /v1/log` is the surface's **only deliberate N+1**: a bounded loop over the
existing `GetLogEntry` rather than a new range query in the kernel (Decision 3).
The budget exists to measure that trade rather than hide it, so it is worth
reading the two log-range rows against each other:

- 100 entries → 1.248 ms p95; 500 entries → 4.915 ms p95. That is **3.94× the
  time for 5× the rows**, i.e. very close to linear with a small fixed transport
  cost amortised across the larger page (~9.8 µs/entry at 500 vs ~12.5 µs at 100).
- Linear-with-no-surprises is exactly the shape a bounded loop should have, and at
  the clamp max it sits **24× inside** its target. **No range-read store method is
  justified by this data.** The signal that would justify one is superlinearity or
  a p95 approaching the target at limit=500; neither is present. If a future
  sprint proposes a kernel range query, this row is the evidence it must overturn.
- The REST commit costs **0.5763 ms p95 against the embedded floor's 0.5421 ms** —
  a ~0.03 ms transport tax over a real loopback socket, so essentially all of the
  commit cost is the kernel's fsync, not the daemon.

## Raw benchmark evidence

```text
goos: darwin
goarch: arm64
pkg: github.com/sunholo-data/ailang-world/host/daemon
cpu: Apple M4 Max
BenchmarkStoreCommit-16    	     200	    473582 ns/op	         0.4715 p50_ms	         0.5421 p95_ms
BenchmarkHeadRead-16       	     200	     66947 ns/op	         0.06617 p50_ms	         0.08033 p95_ms
BenchmarkHealth-16         	     200	     48305 ns/op	         0.04579 p50_ms	         0.06796 p95_ms
BenchmarkRESTCommit-16     	     200	    504435 ns/op	         0.5000 p50_ms	         0.5763 p95_ms
BenchmarkLogRange/limit_100-16         	     200	   1015356 ns/op	         0.9824 p50_ms	         1.248 p95_ms
BenchmarkLogRange/limit_500-16         	     200	   4729270 ns/op	         4.738 p50_ms	         4.915 p95_ms
PASS
ok  	github.com/sunholo-data/ailang-world/host/daemon	2.606s
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
- **Provenance is split, and the split is stated rather than smoothed over.** The
  M2.B sprint executor (codex `gpt-5.6-sol`) authored the REST-commit and
  log-range harnesses but ran under a sandbox that denies `bind(2)` on loopback,
  so it could not execute ANY of the socket benchmarks. It recorded them as
  UNAVAILABLE, quoted the sandbox error, and explicitly declined to invent values
  — which is the only reason this table is trustworthy. **The controller measured
  all six rows on the dev rig outside that sandbox, in the single 200x invocation
  quoted verbatim above.** The three M2.A rows shifted slightly from their A3
  values (e.g. store commit p95 0.6093 → 0.5421 ms) purely because this is a
  different run on a differently-loaded machine; they were re-measured rather than
  carried forward so that no row in this table comes from a different invocation
  than its neighbours.
- Every row is comfortably inside its day-1 target. **No row constitutes a design
  signal for a range-read store method** — see the N+1 analysis above.
