# worldd day-1 benchmark baseline

The recorded day-1 kernel performance budget the charter demands (design doc
`w-worldd-m2` Decision 6 / P3). Later sprints diff against this file on the same
dev rig; CI asserts only that the harness RUNS and REPORTS
(`scripts/bench_worldd.sh --smoke`), never a threshold on shared-runner timings —
noise-gating would be a dishonest gate (S6).

Recorded 2026-07-27 on the development rig.

- Machine: Mac Studio (Mac16,9), Apple M4 Max, 16 cores (12P/4E), 128 GB RAM
- Platform: `darwin/arm64`
- Go: `go version go1.26.4 darwin/arm64`
- Repository commit: `bfbd94e4406a9d0371fe3f11097e0f91b4a077dc` (the A3 branch point)
- AILANG pin: `/tmp/ailang-v0300/ailang`, v0.30.0, commit `e37b370`
- Invocation: `go test -bench . -benchtime 200x -run '^$' ./host/daemon/`

Percentiles are computed by the harness itself from PER-ITERATION wall-clock
samples and emitted via `b.ReportMetric` as `p50_ms` / `p95_ms`; `ns/op` is Go's
own mean and is reported alongside for reference. A mean hides the tail the
budget is actually about, which is why the percentiles are the recorded numbers.

| Operation | Day-1 p95 target | Measured p50 | Measured p95 | Result |
|---|---:|---:|---:|---|
| Store commit (embedded `store.Commit`, the kernel floor) | ≤ 25 ms | 0.4981 ms | 0.6093 ms | **inside budget** (41× headroom) |
| REST commit (`POST /v1/commit`) | ≤ 35 ms | PENDING M2.B | PENDING M2.B | **PENDING M2.B** — route does not exist yet |
| Head read (`GET /v1/head`) | ≤ 5 ms | 0.06975 ms | 0.08596 ms | **inside budget** (58× headroom) |
| Health (`GET /v1/health`) | ≤ 2 ms | 0.04612 ms | 0.06288 ms | **inside budget** (32× headroom) |
| Log range (`GET /v1/log`, limit=100 — the default page) | ≤ 30 ms | PENDING M2.B | PENDING M2.B | **PENDING M2.B** — route does not exist yet |
| Log range (`GET /v1/log`, limit=500 — the clamp max) | ≤ 120 ms | PENDING M2.B | PENDING M2.B | **PENDING M2.B** — route does not exist yet |

The three PENDING rows are the surface's mutation and range-read paths. They land
with their routes in M2.B, and `BenchmarkLogRange` in particular exists because
log-range is the surface's **only deliberate N+1** (a bounded loop over
`GetLogEntry` rather than a new store method, D3): the budget must measure that
trade rather than hide it. This table is refreshed to the full surface in M2.C.

## Raw benchmark evidence

```text
goos: darwin
goarch: arm64
pkg: github.com/sunholo-data/ailang-world/host/daemon
cpu: Apple M4 Max
BenchmarkStoreCommit-16    	     200	    498815 ns/op	         0.4981 p50_ms	         0.6093 p95_ms
BenchmarkHeadRead-16       	     200	     71700 ns/op	         0.06975 p50_ms	         0.08596 p95_ms
BenchmarkHealth-16         	     200	     48088 ns/op	         0.04612 p50_ms	         0.06288 p95_ms
PASS
ok  	github.com/sunholo-data/ailang-world/host/daemon	0.378s
```

## Provenance and honesty notes

- `BenchmarkStoreCommit` opens a **fresh temp-file** SQLite store per run — never
  `:memory:` — so the number includes fsync reality and also exercises A1's
  fail-closed writer lock. Each iteration commits a DISTINCT world, chained on
  the previous head, so the compare-and-append path does real work every time.
- The two HTTP benchmarks measure a **real loopback round-trip** against a
  running daemon (ephemeral port, keep-alive connection warmed OUTSIDE the
  measured region), not an in-process handler call.
- The sprint executor (codex `gpt-5.6-sol`) authored the harness but ran under a
  sandbox that denies `bind(2)` on loopback, so it could measure only the
  store-commit row. It recorded the other rows as UNAVAILABLE and explicitly
  declined to invent values. **The controller measured the two HTTP rows on the
  same dev rig outside that sandbox and completed this table**; the numbers above
  are all real and all from the single 200x invocation quoted verbatim.
- Every row is comfortably inside its day-1 target, so no row constitutes a
  design signal for a range-read store method yet. The one that plausibly will is
  log-range at the clamp max, and it is PENDING by construction until M2.B.
