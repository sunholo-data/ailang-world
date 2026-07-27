# w-worldd-m2 — `ailang-worldd` Local Daemon (CLI + REST over the M1 Host)

**Status**: **ARM A RATIFIED by Mark 2026-07-27 (attended) — ENFORCE single-writer**: reopen the M1 `host/store` freeze (kernel change hereby ratified) for a fail-closed `OpenWriter` (non-waiting OS-backed exclusive lock, `WriterAlreadyActive` on contention, distinct read-only path, subprocess tests). Next iteration applies arm A as the r3 designer revision → sprint-planner. (design DIRECTION accepted by both quorum
reviewers across two rounds; one remaining objection is a ratification-class kernel decision — see
the Park box below). NOT yet routed to sprint-planner.
**Date**: 2026-07-27

> ## ✅ PARK BOX — RESOLVED: Mark picked **(A) Enforce it** (attended, 2026-07-27; recorded in the charter STATUS — the ratified-mission-state channel). Rationale: authority as enforcement, not convention — an embedded writer bypassing the daemon's capability/budget checks is exactly the ambient-authority pattern clause 3 exists to end; the M1 kernel change is ratified by this decision.
> ## (original park box, kept for the record — mission iteration 17, 2026-07-27)
>
> This doc passed the pick-time quorum's DIRECTION check (both `gpt5-6-sol` + `gemini-3-1-pro`
> accept the daemon-shell-over-M1-host direction). Round-1's two objections (bounded waits /
> body-limit; the log-range N+1 hidden from the perf budget) were RESOLVED in revision r1 and
> re-verified on the pinned binary. Re-quorum r2 (the single re-quorum the gate allows) surfaced
> **one new blocking objection from `gpt5-6-sol`** that is **ratification-class** and cannot be
> resolved headless (the mission guardrail: *kernel changes require explicit human ratification,
> quorum evidence attached*). The narrow-refinement carve-out does NOT apply — the objection
> offers a FORK requiring controller judgment, and one arm opens the LANDED M1 `host/store`.
>
> **The objection (verbatim, `gpt5-6-sol`):** *"The single-writer guarantee is asserted but not
> enforced. `ailang-worldd` takes no database lease or process lock, `store.Open` remains
> available to embedded writers, and a second daemon can open the same SQLite file. An operational
> instruction to 'never' do so cannot support the claimed sole-handle model or A6 safe-concurrency
> conclusion."*
>
> **The decision needed (the reviewer's own fork):**
> - **(A) Enforce it** — open the M1 `host/store` freeze to add a fail-closed `OpenWriter`
>   (non-waiting OS-backed exclusive lock, `WriterAlreadyActive` on contention, distinct read-only
>   path, subprocess tests). This is a **kernel change → ratification-class** by the mission
>   guardrail.
> - **(B) Withdraw the claim** — drop the sole-handle model, downgrade the A6 conclusion to rely on
>   SQLite's own transaction serialization (compare-and-append already serializes writers; a second
>   writer sees the updated head → structured 409), and document + test bounded multi-process
>   behavior instead. No M1 change.
>
> On Mark's answer, the next iteration UNPARKS: applies the chosen arm as an r3 designer revision,
> re-verifies, and routes to sprint-planner. Non-blocking nit to fold in then (`gemini-3-1-pro`,
> who PASSED): make CLI `--addr` a global client flag, not only on `health|head`.
**Charter clause**: clause-2 (local-first daemon)
**Verified against**: **`AILANG v0.30.0`** — the pinned released binary at `/tmp/ailang-v0300/ailang`
(`AILANG v0.30.0`, commit `e37b370`, clean — no `-dirty` suffix; the same release the charter's
Verification Log and CI pin, `sha256:ac3174e0f27692bb091d341a518b9473bb78010a4234cbff792aab63c67bb4d3`).
Every `.ail` claim in this doc was checked live on that binary; the checked artifact is
[`sketches/worlddapi.ail`](../sketches/worlddapi.ail).
**Traces to**: [DESIGN.md](../DESIGN.md) §15 (definitive physical architecture), §13.5 (kernel
performance budget), §14 (self-modification boundaries), §17 M2; charter clause 2 + Conflict
Surface ([world-mission.md](../world-mission.md))
**Depends on**: [w-world-library-m1.md](../implemented/w-world-library-m1.md) (**LANDED** —
all 6 milestones, dev `a07ac96`, CI green)
**Estimated**: ~2 days

> **Scope note.** M2 delivers exactly the DESIGN §15 "single local daemon" layer: a long-running
> `ailang-worldd` process that EXPOSES the already-landed M1 host packages over **CLI + REST**,
> binds **localhost by default**, and carries **zero cloud dependencies in the core**. The effect
> broker, capability/budget checks, effect-result recording, and worker isolation are **clause-3
> (`w-effect-broker-m3`) — out of scope**; the MCP/A2A protocol projection is **clause-6
> (`w-mcp-projection`) — out of scope**. Scope creep into either is a defect of this design.

---

## Motivation

Clause 2 of the ratified bar: *"`ailang-worldd` runs on one machine over SQLite with CLI + REST,
zero cloud dependencies in the core."* M1 landed the entire semantic substrate — pure AILANG
kernel (`world/*`), SQLite store with atomic compare-and-append (`host/store`), content
addressing (`host/hashref`, `host/canon`), interpreter archive (`host/archive`), epoch registry
(`host/registry`), and the authoritative replay engine (`host/replay`) — but only as **embedded
Go libraries opened by tests**. Nothing serves them; there is no process, no API, no operator
surface (verified below: zero `net/http` imports, zero `worldd` code, no `cmd/` directory).

M2 turns the library into a service **without moving any logic**: the daemon is a transport
shell. Store, replay, hashing, and registry semantics stay in the host packages where M1 proved
them (S3). And per DESIGN §13.5 and the charter guardrail, the kernel **performance budget is
measured and recorded from the daemon's first commit** — commit latency and read latency are
design constraints now, because fixed kernel overhead looms proportionally larger as models
accelerate.

## Premises (hard constraints — each verified in the Premise Verification Log)

- **P1 — worldd is an in-repo binary in the existing Go module, NOT an upstream `ailang world`
  subcommand.** The charter's Conflict Surface marks this "OPEN for ratification … Revisit only
  on concrete binary-distribution pain" and records the coordinator-recommended DEFAULT as the
  in-repo module (keeps the language repo frozen, lets World iterate at its own cadence). There
  is **no binary-distribution pain** — this repo distributes no binary at all today (no `cmd/`
  directory exists; verified). M1 already created the module (`github.com/sunholo-data/ailang-world`,
  go 1.26.4) and landed six host packages in it. M2 proceeds with the ratified-default layout;
  a human may still reverse this before implementation.
- **P2 — local-first is inviolable.** The REST server binds `127.0.0.1` by default and refuses a
  non-loopback bind (M2 ships no override flag); the daemon core imports only the Go standard
  library, this module's packages, and the already-pinned pure-Go SQLite chain. No cloud SDK, no
  egress, enforced by a dependency-allowlist test (D4).
- **P3 — kernel performance budget from the FIRST daemon commit** (DESIGN §13.5 + charter
  guardrail). The first daemon PR (milestone M2.A) lands the benchmark harness, a committed
  baseline artifact (`bench/BASELINE.md`), and a non-vacuous CI bench-smoke leg. The recorded
  budget must cover **every access pattern the surface ships — including the deliberate N+1
  log-range read** (`BenchmarkLogRange`, lands with the route in M2.B and enters the baseline
  before close-out): a day-1 budget that structurally omits the slowest read path is not a
  budget. This is an M2 acceptance criterion, not a later optimization.
- **P4 — slim kernel / package-first (S3).** Every daemon capability answers "why is this not a
  package?" — answered in full below. REST/CLI are transports; store/replay logic stays in
  `host/*`.
- **P5 — §14 boundary.** The LIVE `ailang-worldd` binary is **excluded from World's
  self-modification scope**; self-replacement follows the attended atomic-mv discipline until a
  shutdown-handoff protocol is designed and proven. This doc restates that boundary; M2 builds
  no self-mod machinery.
- **P6 — verify gate extends automatically.** `scripts/verify_go.sh` (anti-false-green guard:
  fails loudly if `AILANG_BIN` is unset or ≠ v0.30.0) and the CI `go-verify` job already run
  `go build ./... && go test ./...` — new `cmd/` + `host/daemon` packages are swept with no gate
  change. The ONLY CI addition is the bench-smoke step (P3), landed in the same PR that
  introduces the daemon (the M1/M6 same-PR pattern).

## High-Impact Decisions

| Decision | Why High Impact | Chosen By | Deadline | Change Cost |
|----------|-----------------|-----------|----------|-------------|
| worldd = in-repo binary `cmd/ailang-worldd` in the existing module | Resolves the charter's OPEN placement item along the ratified default | M2 (P1) | compile | medium |
| Daemon holds the SOLE store handle while serving; CLI is a REST client | One writer process over SQLite; no second write path to corrupt compare-and-append | M2 | runtime | high |
| REST surface is worldd-NATIVE `/v1/*` (not MCP/A2A) | Draws the clause-6 boundary; prevents accidental serve-api dependency | M2 | API | high |
| Loopback-only bind, no override in M2 | Makes local-first structurally true, not configuration-true | M2 | runtime | low |
| No capability/budget checks in M2 — a loopback trusted-operator surface | Draws the clause-3 boundary explicitly instead of half-building the broker | M2 | API | medium |
| Perf harness + committed baseline + CI smoke from the first daemon commit | The §13.5 budget is unfakeable only if measured before optimization pressure exists | M2 (P3) | first PR | medium |
| Every daemon wait & allocation explicitly bounded (server timeouts, body cap, client deadline, shutdown deadline) | Standing Rule 6 / DESIGN cost discipline: loopback trust does not make unbounded waits or allocations safe | M2 (D7) | first PR | low |

### Design Freeze

- [ ] worldd lives at `cmd/ailang-worldd` + `host/daemon` inside `github.com/sunholo-data/ailang-world` (P1).
- [ ] The daemon process is the single writer; every CLI mutation goes through REST to it.
- [ ] REST v1 is the seven-route table frozen in `sketches/worlddapi.ail` `routes()`.
- [ ] Error mapping is the sketch's `httpStatus`: conflict→409, not-found→404, malformed→400, oversized body→413, internal→500.
- [ ] Default bind `127.0.0.1:7644`; startup refuses any bindHost failing `isLoopbackHost` (Z3-proven predicate).
- [ ] Every wait and allocation is bounded (D7): `http.Server` carries `ReadHeaderTimeout`/`ReadTimeout`/`WriteTimeout`/`IdleTimeout`; commit bodies capped at `maxCommitBytes` (8 MiB, Z3-proven `withinCommitBytes`); CLI client calls carry `context.WithTimeout(defaultClientTimeout)`; shutdown drains under `shutdownTimeout` then hard-closes.
- [ ] Benchmarks report p50/p95 for store-commit, REST-commit, head-read, health, **and log-range (the only deliberate N+1) at limit=100 and the clamp max 500**; baseline committed at `bench/BASELINE.md`.
- [ ] M2 adds **zero** columns, tables, or methods to `host/store`'s schema and **zero** semantics to `world/*`.

---

## Decision 1 — Process & Package Layout (and "why is this not a package?")

**Decision.** Two new Go packages, both transports:

| Path | Role |
|------|------|
| `cmd/ailang-worldd` | `main`: flag parsing, subcommand dispatch (`serve` + client verbs), exit codes |
| `host/daemon` | HTTP server lifecycle, route wiring, JSON codecs, error mapping — importing `host/store`, `host/archive`, `host/registry`, `host/hashref` unchanged |

**S3 answer — why is this not a package?** The daemon is precisely the thing that CANNOT be an
AILANG package: it is the OS-process boundary that future packages are served *through* (DESIGN
§15 places it above the store and below nothing). It contains no semantics: every request
terminates in an existing M1 host call (`store.Commit`, `store.GetObject`, `store.SelectedHead`,
`registry` reads) or a pure predicate already frozen in AILANG. Anything that would grow
semantics here — policies, new transition kinds, projections — is ruled out of M2 and routed to
packages/broker/clause-6 by the Deferred Scope table. The kernel stays thin: ~2 transport
packages, 0 new store methods, 0 new `world/*` functions.

**What the daemon reuses (never reinvents):**

| Need | Reused M1 surface |
|------|-------------------|
| Open/serve the store | `store.Open(path)` → `Store` methods |
| Commit transactions | `store.Commit(Commit)` + `store.IsConflict` / `*ConflictError` |
| Reads | `GetObject/GetWorld/GetLogEntry/GetRegistryHead/SelectedHead/GetVerifyResult` |
| Hash text at the boundary | `hashref.Parse` / `HashRef.String` (reject → 400) |
| Interpreter archival at startup | `archive.New(dbPath)` + `Archive(execPath)` + `ReadManifest` |
| Epoch registry bootstrap | `registry.Bootstrap` (idempotent, divergent-head → structured error) |

## Decision 2 — Single-Writer Process Model

**Decision.** While `ailang-worldd serve` runs, its process owns the only open handle to the
SQLite database. The CLI client verbs (D5) never open the DB; they speak REST to the daemon.
Rationale: `host/store`'s concurrency boundary is the single compare-and-append transaction
(M1 Decision 4) — one serving process serializes commits through it and surfaces stale heads as
structured 409s, instead of two processes contending on SQLite file locks. Embedded use
(tests, offline tools) remains legitimate **when the daemon is not serving that DB** — the store
stays a library (S3); the daemon is not a mandatory gateway.

Commit flow per request: decode JSON → `hashref.Parse` every ref (fail → 400) → assemble
`store.Commit` → call it → `ConflictError` → 409 with both heads (sketch `conflictBody` shape:
observed + selected, so the caller re-plans against `SelectedHead`) → success → 200 with the new
selected head. The daemon adds **no locking of its own**; `net/http` concurrency is safe because
the store transaction is the serialization point (proved by M1's `TestCommitConflictOnStaleHead`).

## Decision 3 — REST Surface v1 (worldd-native)

**Decision.** The frozen route table (checked artifact: `routes()` in
[`sketches/worlddapi.ail`](../sketches/worlddapi.ail) — ai-check green on the pinned binary):

| Method | Path | Backing call | Notes |
|--------|------|--------------|-------|
| GET | `/v1/health` | — | daemon version, db path, archived-interpreter `HashRef` + `ailang --version` from the manifest |
| GET | `/v1/head` | `SelectedHead()` | canonical `algo:digest` text |
| GET | `/v1/worlds/{ref}` | `GetWorld` | 404 on absent, 400 on malformed ref |
| GET | `/v1/objects/{ref}` | `GetObject` | envelope always; payload base64 only with `?payload=true` (bounded response by default) |
| GET | `/v1/log/{index}` | `GetLogEntry` | also `GET /v1/log?from=N&limit=M`: a loop over `GetLogEntry`, `limit` clamped by the Z3-proven `clampLimit` (1..500, default 100) — **no new store method** |
| GET | `/v1/registry/{name}` | `GetRegistryHead` | e.g. `world/epoch-registry/v1` |
| POST | `/v1/commit` | `store.Commit` | request mirrors `store.Commit`: `observedHead`, `objects[]` (payload base64), `nextWorld`, `entry{header…}`; 409 = `ConflictError`; body capped at `maxCommitBytes` via `http.MaxBytesReader` → oversized = 413 (D7) |

Encodings: every digest field is canonical `HashRef` text (`algo:digest`) — the same single
representation M1 fixed for SQLite columns; the daemon parses at the boundary and never invents
a second hash encoding. Error classes are the sketch's `ApiError` ADT with its tests-covered
`httpStatus` mapping (409/404/400/413/500). All list reads are bounded (`clampLimit`,
`ensures { result >= 1 && result <= 500 }`, Z3-verified) and all request bodies are bounded
(`withinCommitBytes`, `ensures { result == (n >= 0 && n <= 8388608) }`, Z3-verified — D7).

This surface is **worldd-native**: its consumers in M2 are the CLI and local operators. It is
NOT the protocol-native boundary — MCP/A2A projection over `ailang serve-api` is clause-6
(`w-mcp-projection`) and depends on nothing here beyond the store.

## Decision 4 — Local-First Bind Policy & Zero-Cloud Enforcement

**Decision.** Default bind `127.0.0.1:7644` (operational default, `--bind` configurable).
Startup validates the host part with the Z3-proven `isLoopbackHost` predicate
(`ensures { result == (host == "127.0.0.1" || host == "::1" || host == "localhost") }`) and
**refuses to start** otherwise — M2 ships no override flag, so local-first is structural, not
advisory. The Go implementation mirrors the frozen predicate and cites the sketch.

Zero-cloud is enforced, not asserted: `host/daemon` gets a dependency-allowlist test
(`TestDaemonDependencyAllowlist`) that walks `go list -deps` and fails if any import falls
outside {Go stdlib, `github.com/sunholo-data/ailang-world/...`, the existing
`modernc.org/sqlite` chain already pinned in `go.mod`}. A cloud SDK physically cannot enter the
daemon core without turning CI red (S6: the gate fails loudly on its null case — an empty
dependency list is also a failure).

The daemon core performs **no network egress**: it only listens. (The CLI client dials the
loopback daemon — that is ingress to worldd, not egress to any external system.)

## Decision 5 — CLI Surface

**Decision.** One binary, subcommand style, mapping 1:1 onto the frozen route table:

```
ailang-worldd serve   --db <path> [--bind 127.0.0.1:7644] [--ailang-bin <path>]
ailang-worldd health|head                     [--addr http://127.0.0.1:7644]
ailang-worldd world get <ref> | object get <ref> [--payload]
ailang-worldd log get <index> | log range --from N [--limit M]
ailang-worldd registry get <name>
ailang-worldd commit --file <commit.json>
```

`serve` lifecycle: `store.Open` → (if `--ailang-bin` given) `archive.Archive` + manifest probe →
`registry.Bootstrap` (idempotent; divergent head = fatal structured error, never silent) → bind
(loopback-checked) → serve (with the D7 server timeouts set) → on SIGINT/SIGTERM, **bounded**
graceful drain via `http.Server.Shutdown(ctx)` under `shutdownTimeout` (10 s), then hard
`Close()` — never an unbounded drain → `store.Close`. Client verbs are thin REST callers sharing
one JSON codec with the server, and **every client call dials under
`context.WithTimeout(defaultClientTimeout)` (30 s, D7)** — a single logic path, so CLI use
continuously exercises the REST surface (the M2 end-to-end test drives the daemon exclusively
through the CLI client functions).

## Decision 6 — Kernel Performance Budget (measured from the FIRST commit)

**Decision.** Milestone M2.A — the PR that creates the daemon — lands all three of:

1. **Harness**: `host/daemon/bench_test.go` — Go benchmarks that, per operation, collect
   per-iteration wall-clock samples and report **p50/p95** via `b.ReportMetric` (ns/op comes
   free). Operations: `BenchmarkStoreCommit` (embedded `store.Commit`, the kernel floor),
   `BenchmarkRESTCommit` (end-to-end POST `/v1/commit` over loopback — added M2.B when the
   route exists), `BenchmarkHeadRead` (GET `/v1/head`), `BenchmarkHealth`, and
   `BenchmarkLogRange` (GET `/v1/log?from&limit` — added M2.B with the route, run at
   **limit=100 (the default page) and limit=500 (the clamp max)** over a pre-populated log).
   `BenchmarkLogRange` exists precisely because log-range is the surface's **only deliberate
   N+1** (a bounded loop over `GetLogEntry` instead of a new store method, D3): the D6 budget
   must measure that trade, not hide it — omitting the slowest read path would fake A9/P3.
   Fresh temp-file SQLite store per benchmark; no `:memory:` (the budget must include fsync
   reality).
2. **Committed baseline artifact**: `bench/BASELINE.md` — machine (darwin/arm64 dev rig), Go
   version, commit hash, and the measured p50/p95 table — **including log-range rows at
   limit=100 and limit=500** — refreshed in M2.C with the full surface. This is the recorded
   day-1 budget the charter demands; later sprints diff against it, and the N+1 pattern is in
   the diff base from day 1.
3. **CI bench-smoke leg** (same PR, `go-verify` job): `scripts/bench_worldd.sh --smoke` runs
   `go test -bench . -benchtime 1x -run '^$' ./host/daemon/` and **fails loudly if zero
   benchmark lines appear in the output** (S6 — `go test -bench` exits 0 on no-match, the
   vacuous-pass class this repo has been burned by twice: V27, B1).

**Initial budgets (targets to validate and record, honest about enforcement):** store commit
p95 ≤ 25 ms; REST commit p95 ≤ 35 ms; head read p95 ≤ 5 ms; health p95 ≤ 2 ms; log-range
(limit=100) p95 ≤ 30 ms; log-range (limit=500, clamp max) p95 ≤ 120 ms — on the dev rig. If the
measured N+1 numbers blow these targets, that is a **recorded design signal** for a future
range-read store method (a doc'd decision, not a drive-by) — exactly what a day-1 baseline is for.
CI **asserts the harness runs and reports** (mechanism, non-vacuous); it does **not** threshold
on shared-runner timings, which would be noise-gating (a dishonest gate per S6). Threshold
regression tracking happens against the committed baseline on the dev rig, where the numbers
are reproducible. Fork-cost (§13.5 / M5) inherits this harness: `BenchmarkStoreCommit`'s
temp-store setup path is the seed measurement for future fork benchmarks.

## Decision 7 — Bounded Waits & Allocations (Standing Rule 6)

**Decision.** Loopback exposure is a trust statement, not a resource-safety statement: a trusted
operator can still wedge an unbounded server with one hung connection or one giant payload, and
the daemon is designed to run unattended for days. Therefore **every wait and every
request-driven allocation in the daemon is explicitly bounded by a named constant** — no
unbounded `http.Server`, no unbounded client call, no unbounded drain, no uncapped body read.
The named defaults (Go `const`/`var` block in `host/daemon/daemon.go`, each citing this
section):

| Constant | Default | Bounds |
|----------|---------|--------|
| `readHeaderTimeout` | 5 s | `http.Server.ReadHeaderTimeout` — slow-header connections |
| `readTimeout` | 30 s | `http.Server.ReadTimeout` — full request read incl. commit bodies |
| `writeTimeout` | 30 s | `http.Server.WriteTimeout` — response write incl. `?payload=true` |
| `idleTimeout` | 120 s | `http.Server.IdleTimeout` — keep-alive connection lifetime |
| `maxCommitBytes` | 8 MiB (8_388_608) | `http.MaxBytesReader` on every body-reading route (v1: only `POST /v1/commit`) — oversized → **413** |
| `defaultClientTimeout` | 30 s | `context.WithTimeout` on **every** CLI REST client call |
| `shutdownTimeout` | 10 s | `http.Server.Shutdown(ctx)` deadline; on expiry → hard `Close()` |

Mechanics:

- The `http.Server` literal sets all four timeouts at construction; a zero value in any of them
  is a test failure, not a default.
- `POST /v1/commit` wraps `r.Body` in `http.MaxBytesReader(w, r.Body, maxCommitBytes)` before
  decoding; the resulting `*http.MaxBytesError` maps to the sketch's `PayloadTooLarge` →
  **413** (`httpStatus` in `sketches/worlddapi.ail`, tests-covered like the other arms). The
  bound itself is frozen semantically: `withinCommitBytes` is **Z3-proven** on the pinned
  binary (`ensures { result == (n >= 0 && n <= 8388608) }`) — the Go constant cites the sketch
  and a unit test asserts `maxCommitBytes == 8388608` so the two cannot drift silently.
  8 MiB comfortably holds M1-scale commits (fixture episodes are KBs) while capping the
  allocation a single request can force; raising it is a doc + sketch change, not a tweak.
- The CLI client owns one `http.Client` and derives `context.WithTimeout(ctx,
  defaultClientTimeout)` per call — no client call can hang past the deadline (bounded even if
  a future non-CLI consumer forgets its own context).
- Shutdown: signal → `Shutdown(ctx)` with `shutdownTimeout` → if the deadline expires, hard
  `Close()` and exit non-zero (drain incomplete is *reported*, never waited out forever).
- **Test — `TestBoundedWaitsAndBodyLimit`** (non-vacuous per S6): (a) constructs the server and
  asserts all four timeout fields equal the named constants and are non-zero (M2.A); (b) POSTs
  a `maxCommitBytes+1`-byte body to `/v1/commit` and asserts **HTTP 413** with the
  `PayloadTooLarge` error class in the JSON body, then POSTs a valid small commit to prove the
  cap doesn't reject legitimate traffic (M2.B, when the route exists); (c) asserts
  `maxCommitBytes == 8388608` matching the Z3-proven sketch bound. The null case fails loudly:
  if the handler ignores the cap and reads the oversized body, (b) sees 200/400 — red.

The unbounded-wait class this rules out is exactly Standing Rule 6's target; D7 makes it a
frozen design property with named constants, a sketch-frozen 413 class, and a test — not an
implementation courtesy.

---

## Milestones (each independently CI-green; ~2d total)

### M2.A — Daemon shell: serve + health/head + perf harness (~0.75d)

- **files**: `cmd/ailang-worldd/main.go` (~120 LOC), `host/daemon/daemon.go` (~240 — config,
  loopback guard, **D7 named bound constants + all four `http.Server` timeouts**, lifecycle
  with bounded `Shutdown(shutdownTimeout)`-then-`Close()`, `/v1/health`, `/v1/head`),
  `host/daemon/daemon_test.go` (~220 — httptest: health/head round-trip, non-loopback bind
  refused, bounded graceful shutdown, dependency allowlist, **`TestBoundedWaitsAndBodyLimit`
  part (a): the four server timeouts equal the D7 constants and are non-zero, and
  `maxCommitBytes == 8388608` matches the Z3-proven sketch bound**),
  `host/daemon/bench_test.go` (~120), `scripts/bench_worldd.sh` (~40),
  `bench/BASELINE.md` (~30), `.github/workflows/ci.yml` (+~8: bench-smoke step in `go-verify`)
- **acceptance_checks**: daemon starts on a temp store, serves health+head, refuses
  `--bind 0.0.0.0:…`; `TestDaemonDependencyAllowlist` green; server timeouts asserted set
  (D7); shutdown drain bounded by `shutdownTimeout` then hard-close; bench smoke emits ≥1
  benchmark line and the committed `bench/BASELINE.md` records store-commit/head-read/health
  p50+p95 (log-range rows land with the route in M2.B — noted as pending in the baseline)
- **verify_commands**: `AILANG_BIN=/tmp/ailang-v0300/ailang ./scripts/verify_go.sh` ·
  `./scripts/verify_ail.sh` · `./scripts/bench_worldd.sh --smoke`
- **ci_green_boundary**: both existing CI jobs + the new bench-smoke step green on the PR; no
  REST mutations exist yet, so nothing depends on M2.B

### M2.B — Full REST v1: reads + commit + conflict semantics (~0.75d)

- **files**: `host/daemon/handlers.go` (~290 — worlds/objects/log/log-range/registry GETs,
  POST commit **body-capped via `http.MaxBytesReader(maxCommitBytes)`**, `ApiError`→status
  mapping per sketch incl. `PayloadTooLarge`→413), `host/daemon/handlers_test.go` (~350 —
  read round-trips; malformed ref→400; absent→404; **conflict→409 carrying observed+selected
  heads**; log-range clamped at 500; payload gated behind `?payload=true`;
  **`TestBoundedWaitsAndBodyLimit` part (b): `maxCommitBytes+1` body → 413 with
  `PayloadTooLarge` class, then a valid small commit succeeds**),
  `host/daemon/bench_test.go` (+~70: `BenchmarkRESTCommit` + **`BenchmarkLogRange` at
  limit=100 and limit=500 over a pre-populated log** — the deliberate N+1 enters the harness
  the moment the route exists)
- **acceptance_checks**: a genesis+commit episode driven entirely over REST reproduces the
  store-level result byte-for-byte; stale-head commit returns 409 whose body lets the client
  re-plan (asserted by re-planning and committing successfully); oversized commit body → 413;
  no new store methods appeared (`git diff` on `host/store/` is empty); bench smoke output now
  contains `BenchmarkLogRange` lines
- **verify_commands**: same two gates + `./scripts/bench_worldd.sh --smoke`
- **ci_green_boundary**: full test suite green including 409/400/404/413 paths; bench smoke now
  includes REST commit + log-range

### M2.C — CLI client + end-to-end + baseline refresh + close-out (~0.5d)

- **files**: `cmd/ailang-worldd/cli.go` (~210 — client verbs, shared JSON codec, **every call
  under `context.WithTimeout(defaultClientTimeout)` (D7)**), `cmd/ailang-worldd/cli_test.go`
  (~180 — end-to-end: spawn `serve` on an ephemeral loopback port, drive every verb through
  the CLI paths, commit + conflict round-trip), `bench/BASELINE.md` (refresh with
  full-surface numbers **including the log-range N+1 rows at limit=100/500**), `README.md`
  (+~15 operator quickstart)
- **acceptance_checks**: every route in the frozen table is reachable via a CLI verb;
  end-to-end test passes against a real subprocess daemon; client calls carry the bounded
  deadline (no unbounded dial in `cli.go` — asserted by construction in review + a deadline
  test); refreshed baseline committed with log-range p50/p95; all acceptance boxes below
  checkable
- **verify_commands**: same gates; final sweep `AILANG_BIN=/tmp/ailang-v0300/ailang
  ./scripts/verify_go.sh && ./scripts/verify_ail.sh && ./scripts/bench_worldd.sh --smoke`
- **ci_green_boundary**: this milestone LANDS the item — doc moves to `implemented/` with all
  boxes checked

## Files to Create/Modify (aggregate)

| File | Est. LOC | Change |
|------|---------:|--------|
| `cmd/ailang-worldd/main.go` | ~120 | new: entrypoint + subcommand dispatch |
| `cmd/ailang-worldd/cli.go` + `cli_test.go` | ~390 | new: REST-client verbs (bounded deadlines, D7) + end-to-end test |
| `host/daemon/daemon.go` | ~240 | new: config, lifecycle, loopback guard, D7 bound constants + server timeouts, health/head |
| `host/daemon/handlers.go` | ~290 | new: REST v1 reads + body-capped commit + error mapping (incl. 413) |
| `host/daemon/daemon_test.go` + `handlers_test.go` | ~570 | new: unit + conflict + allowlist + `TestBoundedWaitsAndBodyLimit` |
| `host/daemon/bench_test.go` | ~190 | new: p50/p95 perf harness incl. `BenchmarkLogRange` (the N+1) |
| `scripts/bench_worldd.sh` | ~40 | new: non-vacuous bench runner (fails on zero benchmarks) |
| `bench/BASELINE.md` | ~30 | new: committed day-1 budget baseline (incl. log-range rows) |
| `.github/workflows/ci.yml` | ~8 | bench-smoke step in the existing `go-verify` job |
| `design_docs/sketches/worlddapi.ail` | 123 | **already created + checked with this doc** (r1: + `withinCommitBytes`, `PayloadTooLarge`→413) |

Estimated implementation total: ~1,880 LOC. `host/store`, `host/hashref`, `host/canon`,
`host/archive`, `host/registry`, `host/replay`, `world/*`, `scripts/verify_ail.sh`,
`scripts/verify_go.sh`: **unchanged**.

## Conflict Surface

The iteration-0 quorum's standing objection demands this boundary be explicit.

- **vs `ailang serve-api --mcp/--a2a` (the load-bearing boundary).** Already-verified charter
  premises (code inspection + live test 2026-07-23, charter Premise Verification Log — not
  re-derived here): `internal/apiserver` is **stateless** (no persistence/scheduler API
  package-wide; request-scoped serving of module exports) and is a Go **`internal/` package,
  not importable cross-repo**. Reuse therefore has exactly three paths — (a) primary: expose
  World's registry AS `.ail` modules served by serve-api, (b) sidecar process, (c) upstream
  export issue — and ALL THREE belong to **clause-6 / `w-mcp-projection`, not M2**. The M2
  boundary: **worldd's REST API is worldd-NATIVE** — its own `/v1/*` HTTP surface over the
  store, serving state that apiserver has no machinery for (store, log, heads, atomic commit).
  M2 does **not** import, spawn, wrap, or depend on serve-api; it does not speak MCP or A2A; it
  reserves no serve-api ports. Conversely M2 must not preempt clause-6: no "temporary" MCP
  endpoint, no tool-shaped projections of transitions. The daemon justification runs exactly as
  the charter states: *a new daemon is justified by state, not by protocol* — and M2 adds only
  the state-serving layer.
- **vs the effect broker (clause-3 / `w-effect-broker-m3`).** `POST /v1/commit` enforces ONLY
  what M1 already enforces: store-level compare-and-append + content verification + the pure
  `world/*` predicates the caller ran. **No capability checks, no budget checks, no effect
  execution, no effect-result recording, no worker isolation.** Consequence, stated honestly:
  the M2 REST surface is a **loopback-bound trusted-operator surface** — the same trust level as
  M1's embedded library, now reachable over localhost HTTP. That is pre-authority by design;
  bolting a partial capability check onto M2 would half-build the broker in the wrong layer and
  is ruled out. Clause-3 adds the broker BETWEEN daemon and external systems (DESIGN §15 stack)
  without changing this REST surface's store semantics.
- **vs Coordinator / Collaboration Hub.** Pattern overlap only (the future approval inbox is
  clause-5, not an M2 endpoint). Patterns port; schemas do not; **no shared database** — worldd
  serves exactly one SQLite world store. No Hub code or schema is read, imported, or extended.
- **vs the §14 self-modification boundary.** The live daemon binary is outside World's self-mod
  scope (P5). M2 ships no self-update, no hot-reload, no in-band binary replacement. Replacing
  a running worldd follows the attended atomic-mv discipline.
- **vs M1's embedded consumers.** The store remains an embeddable library; M2 does not make the
  daemon a mandatory gateway. The one operational rule: never run the daemon and an embedded
  writer against the same DB file concurrently (D2 single-writer). Tests and `host/replay`
  fixtures keep opening temp stores directly — the M1 test suite is untouched and must stay
  green.
- **vs the Go namespace of this repo.** Verified negative-existence (log below): zero `net/http`
  imports, zero `worldd` identifiers in Go code, no `cmd/` directory — the daemon collides with
  nothing. The new sketch enters `verify_ail.sh`'s sweep with an EMPTY required set (sketches
  cannot mask required world/ identities nor perturb the exact totals — verified against the
  script's manifest logic and by a live full-sweep run).

## Systemic-Issue Audit

Is a daemon-shaped gap being patched one endpoint at a time? No — the audit runs the other way:
M2's risk is the **OS gravity well** (DESIGN §12.4), kernel growth disguised as convenience
endpoints. The design counters it structurally: one JSON codec, one error-mapping ADT (frozen in
the checked sketch), one backing host call per route, zero new store methods, and a Deferred
Scope table that names where each tempting addition actually belongs. The `log?from&limit` range
read is the deliberate test case: implemented as a bounded loop over the EXISTING `GetLogEntry`
rather than a new store query — the pattern every future read endpoint follows — **and its N+1
cost is measured, not hidden: `BenchmarkLogRange` puts it in the D6 harness and the committed
baseline from the moment the route exists**. Endpoint count
is capped by the frozen `routes()` table; growing it is a doc change, not a drive-by commit.

## Deferred Scope

| Item | Belongs to | Boundary |
|------|-----------|----------|
| Capability/budget checks, effect execution + recording, worker isolation | `w-effect-broker-m3` (clause-3) | M2 commit = store semantics only; trusted-operator surface |
| MCP tools / A2A card / serve-api reuse paths (a)/(b)/(c) | `w-mcp-projection` (clause-6) | M2 REST is worldd-native; no protocol projection |
| Approval inbox, provenance walk endpoints | `w-approval-inbox` (clause-5) | not an M2 route |
| Non-loopback bind, auth, TLS | post-broker ratification | M2 refuses non-loopback, ships no override |
| Replay-over-REST | later, if demanded | replay stays a bounded embedded/CLI operation via `host/replay` |
| Fork-cost benchmark (M5 speculation) | M5, seeded by this harness | M2 records commit/read budget only |

## Acceptance Criteria

- [ ] `ailang-worldd serve` runs over a SQLite store using ONLY landed M1 packages — `git diff`
  on `host/{store,hashref,canon,archive,registry,replay}/` and `world/` is empty.
- [ ] All seven frozen routes implemented; CLI verbs cover the table 1:1; end-to-end test drives
  a genesis→commit→read episode through the CLI against a real subprocess daemon.
- [ ] Stale-head commit over REST returns 409 with observed+selected heads; the test re-plans
  from the body and commits successfully.
- [ ] Malformed `HashRef` → 400; absent object/world/log/registry → 404; log-range limit
  clamped to ≤500; oversized commit body (> `maxCommitBytes`) → 413 with the `PayloadTooLarge`
  class.
- [ ] **Bounded waits (D7)**: the `http.Server` sets `ReadHeaderTimeout`/`ReadTimeout`/
  `WriteTimeout`/`IdleTimeout` to the named constants; `POST /v1/commit` reads through
  `http.MaxBytesReader(maxCommitBytes)`; every CLI client call carries
  `context.WithTimeout(defaultClientTimeout)`; shutdown is `Shutdown(ctx≤shutdownTimeout)` then
  hard `Close()`. `TestBoundedWaitsAndBodyLimit` green, non-vacuously: timeouts asserted
  non-zero AND oversized-body → 413 AND `maxCommitBytes == 8388608` (the Z3-proven sketch bound).
- [ ] Non-loopback `--bind` refused at startup; default bind is `127.0.0.1:7644`.
- [ ] `TestDaemonDependencyAllowlist` proves zero cloud/network-egress dependencies in
  `host/daemon` + `cmd/ailang-worldd`.
- [ ] Startup archives `--ailang-bin` (when given) via `host/archive` and bootstraps
  `world/epoch-registry/v1` via `registry.Bootstrap`, idempotently.
- [ ] Perf harness lands in the FIRST daemon PR; `bench/BASELINE.md` committed with p50/p95 for
  store-commit, REST-commit, head-read, health, **and log-range at limit=100 and limit=500
  (`BenchmarkLogRange` — the deliberate N+1 is in the day-1 budget, not hidden from it)**; CI
  bench-smoke step fails on zero benchmarks.
- [ ] `sketches/worlddapi.ail` stays green in `verify_ail.sh` alongside the prior 8 modules;
  the 4 required world/ identities and 14 required named tests are unperturbed.
- [ ] `scripts/verify_go.sh` (with pinned `AILANG_BIN`) and both CI jobs green on every
  milestone PR.
- [ ] M2 stops before broker, capabilities, budgets, isolation, MCP, A2A, and any serve-api
  dependency.

## Axiom Compliance

| Axiom | Score | Justification |
|-------|-------|---------------|
| A1: Determinism | +1 | Transport shell adds no nondeterminism; commit semantics remain M1's compare-and-append; replay machinery untouched |
| A2: Replayability | +1 | Same immutable objects/log via the same store calls; REST commit is byte-equivalent to embedded commit (tested) |
| A3: Effect Legibility | 0 | Daemon I/O is host-boundary Go (S2-conformant); broker-recorded effects are explicitly clause-3 |
| A4: Explicit Authority | 0 | Honest pre-authority: loopback-bound trusted-operator surface, no ambient path beyond what M1's library already granted locally; enforcement lands with the broker |
| A5: Bounded Verification | +1 | Verify cache untouched; all list reads clamped by a Z3-proven bound (`clampLimit`) and all request bodies by a second Z3-proven bound (`withinCommitBytes`, D7) |
| A6: Safe Concurrency | +2 | Single-writer daemon + store transaction as the sole serialization point; stale heads surface as structured 409 re-plan signals |
| A7: Machines First | +2 | Route table, error ADT, and bind/clamp predicates frozen in a compiler-checked artifact; JSON everywhere; canonical hash text only |
| A8: Minimal Syntax | 0 | No language surface |
| A9: Cost Visibility | +2 | p50/p95 budget measured and committed from the first daemon commit, **covering every shipped access pattern including the deliberate N+1 log-range read (`BenchmarkLogRange` at limit=100/500)**; every daemon wait/allocation bounded by named constants (D7); non-vacuous CI smoke |
| A10: Composability | +1 | Transports over unchanged packages; store stays embeddable; clause-3/6 stack on top without rework |
| A11: Structured Failure | +1 | `ApiError` ADT mirrors structured host errors; startup failures (bad bind, divergent registry head) are fatal and named |
| A12: System Boundary | +2 | Daemon = process boundary exactly where DESIGN §15 draws it; semantics stay in `world/*` + `host/*` |

**Net Score: +13** ✅ — hard axioms A1/A3/A4/A7 all ≥ 0.

## Premise Verification Log (live evidence, this session, 2026-07-27)

| Claim | Check | Result |
|-------|-------|--------|
| Pinned binary is clean released v0.30.0 | `/tmp/ailang-v0300/ailang --version` | ✓ `AILANG v0.30.0`, commit `e37b370`, no `-dirty` suffix |
| P1: placement is charter-OPEN with in-repo default; no distribution pain | Read charter Conflict Surface ("OPEN for ratification … Revisit only on concrete binary-distribution pain"; coordinator default = in-repo module) | ✓ quoted; repo ships no binary today (no `cmd/`), so no distribution pain exists |
| Go module already exists in-repo | Read `go.mod` | ✓ `module github.com/sunholo-data/ailang-world`, `go 1.26.4`, `modernc.org/sqlite v1.54.0` (pure-Go) |
| M1 store API surface as relied on by D1–D3 | `grep '^func\|^type' host/store/store.go` | ✓ `Open`, `PutObject/GetObject`, `PutWorld/GetWorld`, `GetLogEntry`, `SetRegistryHead/GetRegistryHead`, `PutVerifyResult/GetVerifyResult`, `SelectedHead/SelectHead`, `Commit`, `ConflictError`+`IsConflict` all present |
| M1 archive/registry APIs as relied on by D5 | `grep '^func\|^type' host/archive/archive.go`; registry.go | ✓ `archive.New(storeDBPath)`, `(*Archive).Archive(execPath)`, `Resolve`, `ReadManifest`, `Manifest`; `registry.Bootstrap(s, releaseString)` idempotent |
| **Negative-existence: no daemon/HTTP/cmd code exists** | `grep -rl "net/http" --include="*.go" .` → 0 files; `grep -rl "worldd" --include="*.go" .` → 0; `ls cmd` → "No such file or directory" | ✓ all three empty — the daemon collides with nothing; M2 creates these names fresh |
| `scripts/verify_go.sh` exists with anti-false-green guard | Read script | ✓ fails rc=1 if `AILANG_BIN` unset or binary ≠ v0.30.0; runs `go build ./... && go test ./... -count=1` repo-root-relative |
| CI `go-verify` job pattern to extend | Read `.github/workflows/ci.yml` | ✓ pins v0.30.0 by TAG, sha256-verifies, asserts version, exports `AILANG_BIN`, calls `verify_go.sh`; bench-smoke step slots after it |
| serve-api stateless + `internal/` (Conflict-Surface premises) | Charter Premise Verification Log rows (code inspection + live MCP test 2026-07-23) | ✓ cited, not re-derived: no persistence/scheduler API; not importable cross-repo; path (a) live-tested for static exports |
| Sketch compiles + contracts prove on the pinned binary (re-run after r1) | `cd design_docs && /tmp/ailang-v0300/ailang ai-check -timeout 5s sketches/worlddapi.ail` | ✓ `check.passed: true`, `verify: {verified: 3, counterexample: 0, skipped: 0, errors: 0}` — `isLoopbackHost` + `clampLimit` + `withinCommitBytes` Z3-verified (purity via explicit `! {}`, the world/logepoch pattern) |
| Sketch inline tests pass (re-run after r1) | `/tmp/ailang-v0300/ailang test --format json sketches/worlddapi.ail` | ✓ `failed_tests: 0`; 15/15 named inline tests pass (+3 contract-derived properties = 18 passed) — incl. `PayloadTooLarge`→413 and the three `withinCommitBytes` boundary cases |
| Contract-on-ADT limitation honored (httpStatus tests-only) | `w-m1-ailang-hardening` documented limitation: ADT-bearing sorts Z3-error as `unknown sort`, exit 0 silently | ✓ `httpStatus` (ADT param) is tests-covered by policy; no verify.errors in the sketch's JSON confirms nothing silently failed |
| New sketch cannot perturb the required-check manifest | Read `scripts/verify_ail.sh` (`REQUIRED_VERIFIED` keys world/-only; sketches "carry EMPTY required sets and are excluded from the total"; Leg 2 runs `world/` only) + full live sweep | ✓ `AILANG_BIN=/tmp/ailang-v0300/ailang ./scripts/verify_ail.sh` → "4/4 required world/ identities verified across **9** module(s)" + "all 14 required named tests pass" + "verify gate PASSED" |
| §15 places worldd above the store, below broker/workers | Read DESIGN.md §15 stack diagram | ✓ `worldd → semantic store → effect broker → workers → external systems`; storage = SQLite; interfaces "CLI first, REST second" |
| §13.5 requires day-1 perf budget | Read DESIGN.md §13.5 + charter guardrail | ✓ "the kernel carries a performance budget from day 1, not as a later optimization"; fork cost named a kernel design constraint |
| §14 excludes the live daemon from self-mod | Read DESIGN.md §14 hard boundary 3 | ✓ "The live daemon binary… attended atomic-mv discipline… until a shutdown-handoff protocol is designed and proven" |
| No duplicate/overlapping design doc | `ls design_docs/planned/ design_docs/implemented/` | ✓ planned/ holds only `w-log-epoch-decision.md` (settled); M1 doc's Deferred Scope explicitly routes daemon/REST/CLI to THIS item |

### Quorum revision r1 (2026-07-27) — both blocking objections resolved

- **Objection 1 (gpt5-6-sol, bounded waits — Standing Rule 6 / DESIGN cost discipline):**
  unbounded `http.Server`, unbounded client calls, unbounded graceful shutdown, uncapped
  base64 commit bodies. **Resolved by new Decision 7** — named bound constants for all four
  server timeouts (`readHeaderTimeout` 5 s / `readTimeout` 30 s / `writeTimeout` 30 s /
  `idleTimeout` 120 s), `http.MaxBytesReader` at `maxCommitBytes` = 8 MiB on `POST /v1/commit`
  → **413** (new `PayloadTooLarge` arm in the sketch's `ApiError`/`httpStatus`),
  `context.WithTimeout(defaultClientTimeout` = 30 s`)` on every CLI client call, and
  `Shutdown(ctx ≤ shutdownTimeout` = 10 s`)` then hard `Close()`. The body bound is
  **Z3-proven** in the sketch (`withinCommitBytes`, verified on the pinned binary); enforced
  by the non-vacuous `TestBoundedWaitsAndBodyLimit` (M2.A part a/c, M2.B part b) and new
  acceptance criteria. Design Freeze, D3, D5, milestones, files table, and A5/A9 rows updated.
- **Objection 2 (gemini-3-1-pro, N+1 missing from the perf harness):** the deliberate
  log-range N+1 (`GET /v1/log?from&limit`, a bounded loop over `GetLogEntry`) was absent from
  the D6 harness/baseline while claiming A9/P3. **Resolved:** `BenchmarkLogRange` added to
  `host/daemon/bench_test.go` at **limit=100 (default page) and limit=500 (clamp max)**,
  landing in **M2.B the moment the route exists**; `bench/BASELINE.md` MUST carry log-range
  p50/p95 rows (M2.C refresh); initial targets recorded (p95 ≤ 30 ms @100, ≤ 120 ms @500) with
  the honest escape hatch that a blown target is a recorded design signal for a future
  range-read store method. D6, P3, M2.B/M2.C bullets, the A9 axiom row, the Systemic-Issue
  Audit, and the Acceptance Criteria all updated.
- **Direction unchanged:** daemon shell + CLI/REST over the landed M1 host; broker = clause-3,
  MCP/A2A = clause-6 remain out of scope. Both gates re-run on the pinned binary after the
  sketch change — tails below.

### Gate Output Tails (re-run after r1, pinned binary, 2026-07-27)

`cd design_docs && /tmp/ailang-v0300/ailang ai-check -timeout 5s sketches/worlddapi.ail`:

```text
  "check": { "passed": true, "error_count": 0, "errors": [] },
  "verify": {
    "available": true,
    "verified": 3,
    "counterexample": 0,
    "skipped": 0,
    "errors": 0,
    "results": [
      { "function": "isLoopbackHost",    "status": "verified", ... },
      { "function": "clampLimit",        "status": "verified", ... },
      { "function": "withinCommitBytes", "status": "verified", ... }
```

(`ailang test --format json sketches/worlddapi.ail` → `"failed_tests": 0`,
`"passed_tests": 18`, `"success": true` — 15 named tests + 3 contract properties.)

`AILANG_BIN=/tmp/ailang-v0300/ailang ./scripts/verify_ail.sh` (full sweep, r1 sketch included):

```text
   ai-check design_docs/sketches/transitions.ail
   ai-check design_docs/sketches/worlddapi.ail
   ai-check design_docs/sketches/worldkernel.ail
   ai-check design_docs/sketches/worldtypes.ail
   ai-check world/contracts.ail
   ai-check world/logepoch.ail
   ai-check world/transitions.ail
   ai-check world/types.ail
   ✓ 4/4 required world/ identities verified across 9 module(s)
── Leg 2: named inline tests
   ✓ all 14 required named tests pass (failed_tests=0)
✓ verify gate PASSED: 4 required identities verified, 14 named tests pass
```

## Related Documents

- [w-world-library-m1.md](../implemented/w-world-library-m1.md) — the LANDED host this daemon
  wraps; its Deferred Scope table names this item as the daemon/REST/CLI boundary
- [w-m1-ailang-hardening.md](../implemented/w-m1-ailang-hardening.md) — verifier facts this doc
  obeys: explicit `! {}` for provable contracts, ADT-sort Z3 limitation, non-vacuous manifest gates
- [w-log-epoch-decision.md](w-log-epoch-decision.md) — settled identities the REST surface
  renders (canonical `HashRef` text, frozen `LogHeader`)
- [world-mission.md](../world-mission.md) — clause 2, Conflict Surface (serve-api premises,
  placement OPEN item), guardrails
- [DESIGN.md](../DESIGN.md) — §13.5, §14, §15 (definitive architecture), §16, §17 M2
- [sketches/worlddapi.ail](../sketches/worlddapi.ail) — the compiler-checked typed API artifact
  for this doc (ai-check green, 3 Z3-verified contracts incl. the D7 `withinCommitBytes` body
  bound, 15 named tests incl. `PayloadTooLarge`→413)
