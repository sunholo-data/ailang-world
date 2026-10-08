# w-worldd-agui-event-stream — row 148

- Status: **IMPLEMENTED — iteration 243, PR #231 → `38db793` (quorum r1/r2 BLOCKED → r3 narrow-refinement carve-out; sonnet judge r1 90 → r2 92, zero blocking).** Pick-time design, revision r1. Supersedes the row-148 parts of
  the attended DRAFT [`w-world-live-surface.md`](../planned/w-world-live-surface.md) (D2–D5, §5–§7). Rows 149–151 stay there.
- Base `d6334c2`; premises in §2; probes in [`../verification/world-row148-design/`](../verification/world-row148-design/).
  Binding: [coding-standards.md](../coding-standards.md) S2, S3, S6, S7; `D-WORLD-64`.

## 1. Problem and clause mapping

worldd has no live event stream. A human or agent that wants to see a commit has to poll
`GET /v1/log` or reload `/workbench` (V1). Because of that, P5 (attention choreography, "quiet is
health") cannot be rendered, and rows 149 (live workbench) and 151 (agent-generated views) have no
transport to consume.

**Row 148 adds one additive route, `POST /agui/`.** It speaks AG-UI 1.0 over SSE, and every event
is derived from committed log entries and nothing else.

**Clause mapping.** Clause 6 +5 (an open protocol, no new dialect); clause 5 indirectly via 149/151.
A **release row** under `D-WORLD-64`: it closes no UNMET clause, and the loop takes it only via that
ruling's scoped rule-(d) exception.

**Interface for 149/151 (named, not designed):** the D3 route contract; the `world/agui-state/v1`
state and its deltas; and the rule that any new `CUSTOM` name is `world.<noun>.<verb>` in
`host/agui`'s closed table, with a golden test (row 151's `world.view.a2ui` will be one).

## 2. Verification Log (first-party, at `d6334c2`)

Probes: `world-row148-design/probes_test.go.txt` + `probes-output.txt`, run as a scratch module
`replace`d onto this worktree (`probes.go.mod.txt`), Go 1.26.6 darwin/arm64.

| # | Premise | Command | Output (abridged, verbatim values) |
|---|---|---|---|
| V1 | Route table | `grep -n 'HandleFunc' host/daemon/daemon.go` | 859–868, the ten `/v1` patterns: `GET /v1/health`, `GET /v1/head`, `GET /v1/worlds/{ref}`, `GET /v1/objects/{ref}`, `GET /v1/objects/by-semantic-id/{name...}`, **`GET /v1/log/{index}` (`:864`)**, `GET /v1/log`, `GET /v1/registry/{name...}`, `POST /v1/commit`, `GET /v1/receipts/{id}`. Then 869–876: `GET /workbench`, `GET /.well-known/agent.json`, `POST /a2a/`, `POST /mcp/`. No stream route. The draft's `:800-821` is **stale**: the table moved down 57 lines. |
| V2 | Read posture | `sed -n 877,887p host/daemon/daemon.go`; `sed -n 48,53p host/daemon/middleware.go` | `isProtected` (`:885`) is `POST && path == "/v1/commit"` only. `Wrap` does `if !protected(r) { mux.ServeHTTP(w, r); return }` (`middleware.go:50`). An unprotected GET **or POST** passes through untouched: no credential is read and none is rejected. |
| V3 | No per-session read filter exists | V2 plus `grep -n 'FromContext' host/daemon/handlers.go` | 1 hit, `:732`, inside `handleCommit` (positive control: the pattern does match). No `/v1` GET handler consults a binding. The draft D4's "filtered by session when a credential is presented" has **nothing to reuse**. That would be a new authority feature. |
| V4 | Transport deadlines | `sed -n 88,104p;143p host/daemon/daemon.go`; `:943-944` | `readHeaderTimeout 5s`, `readTimeout 30s`, **`writeTimeout 30s`**, `idleTimeout 120s`, `readDeadline 10s`, `shutdownTimeout 10s` (`:127`). Every value is wired into `newServer` and pinned by `TestBoundedWaitsAndBodyLimit`. |
| V5 | `WriteTimeout` caps a streamed response | probe `TestP1WriteTimeoutCutsStream` (server `WriteTimeout=2s`, SSE frame every 250 ms, `Flush` each) | `stream ended after 2.02s, last "id: 7", scanErr=unexpected EOF`. **Control** `TestP1ControlNoWriteTimeout` (same handler, no timeout): `ended after 3.52s, last "id: 13", scanErr=<nil>`. GOROOT `net/http/server.go:3002-3005`: "reset whenever a new request's header is read … does not let Handlers make decisions on a per-request basis". **Consequence: the draft's 30-min connection is impossible under frozen D7.** |
| V6 | Deadline relaxation is policed | `sed -n 1524,1545p host/projection/projection_test.go` | `TestProjection_NoDeadlineTampering` bans `ResponseController` / `SetWriteDeadline` / `SetReadDeadline` in projection sources (AC14). The same D7 posture applies to any new route. |
| V7 | Shutdown waits on active handlers | `sed -n 984,1012p host/daemon/daemon.go`; GOROOT `server.go:3128-3131` | `Shutdown()` = `drain(d.srv, d.drainTimeout)` → `srv.Shutdown(ctx)`, which waits "indefinitely for connections to return to idle". `grep -rn RegisterOnShutdown host` returns **0**. Control: `grep -c 'func (s \*Server) RegisterOnShutdown' $GOROOT/src/net/http/server.go` returns 1. A stream that ignores shutdown would hold the drain to its 10 s expiry. |
| V8 | Log read API | `sed -n 685,750p host/store/store.go`; `sed -n 120,140p` | `GetLogEntry(ctx, index)` (`:687`) is one `SELECT … WHERE entry_index = ?` (`:706`) and returns `ok=false` when the index is absent. `LogEntry = {Header{EntryIndex, SemanticsEpoch, TransitionFn, Interpreter, PrevEntryHash, WrittenBy}, EntryHash, TransitionRef}`. `entry_index INTEGER PRIMARY KEY` (`schema.sql:40`). Every read calls `requireDeadline` (`:239`). |
| V9 | `/v1/log` range semantics | `sed -n 604,665p host/daemon/handlers.go`; `:243` | `handleLogRange` loops `GetLogEntry(from+offset)` and **breaks at the first absent index**. That is a deliberate N+1 with clamp ≤500. Wire shape: `logJSON(entry)` → `{"header":{entryIndex,semanticsEpoch,transitionFn,interpreter,prevEntryHash,writtenBy},"entryHash","transitionRef"}`. The read seam is `readStore` (`daemon.go:475-485`). |
| V10 | The log can have index gaps | probe `TestP4IndexGap` | After entry 1, `Commit` at `entryIndex 5` gives `err=<nil>; entry2 present=false entry5 present=true`. A further commit with a wrong `prevEntryHash` also gives `err=<nil>`. The store enforces neither density nor linkage. `grep -rn 'EntryIndex !=\|LogHead !=' host` (non-test) finds no check. The positive-shape control `transitionreg.go:260` (`next.Revision != current.Revision+1`) shows what such a check looks like. End-to-end over REST: see V26 (a gapped `POST /v1/commit` lands). **Consequence: tail-by-`index+1` stalls silently at a gap. So does `/v1/log`.** |
| V11 | No change notification; poll cost | `grep -n 'SetMaxOpenConns' host/store/store.go`; probe `TestP2PollCost` (file store, 1000 entries, 20 000 iterations) | `db.SetMaxOpenConns(1)` (`:382`). There is no hook or notify API in `host/store` (rollback journal, no `journal_mode` pragma, 0 hits). Per op: `GetLogEntry(absent)=10.14µs`, `GetLogEntry(hit)=11.61µs`, `SelectedHead=8.28µs`. Cross-check: `bench/BASELINE.md:268` `BenchmarkLogRange/limit_100` = 2.47 ms over HTTP. |
| V12 | A "second process committing" cannot happen while worldd runs | probes `P3a`, `P3b`, `P3c`; `sed -n 1,25p host/store/writer_lock.go` | In-process second `store.Open`: `isWriterAlreadyActive=true`. Child **process** `store.Open`: `another process already holds the writer lock … isWriterAlreadyActive=true`. A read-only handle sees a writer's new entry on its next read (`ok=false` before, `ok=true` after, 34.7µs). Every log write goes through the daemon's one `store.Store`. **The draft D2 premise "a broadcaster misses writes by a second process" is false.** |
| V13 | No SSE writer and no `Flush` in production code | `grep -rln text/event-stream host cmd \| grep -v _test`; `grep -rn 'http\.Flusher\|\.Flush()' host cmd \| grep -v _test \| wc -l` | Only `cmd/ailang-worldd/mcpclient.go`, a **client** parser (`sseEvents`, `:131`). The flush count is **0**. Control: the same pattern on the probe file returns 2. The MCP server SSE is upstream and request-scoped (`TestProjection_*` source gates, `mcp_source_test.go:63`). The middleware passes `w` unchanged, so `http.Flusher` reaches the handler (`middleware.go:51`). |
| V14 | Decision and effect events have no log-derivable source | `grep -rln 'decision-packet\|DecisionPacket' host cmd \| grep -v _test`; `grep -rn commit_objects host/store/*.go` | Only `host/broker/testdata/…json`. Control: `epoch-registry/v1` hits 3 production files. `DecisionPacket` is typed (`world/types.ail:159`) but no host path commits one. `commit_objects` is read only object→entries (`object_commits.go:20`). There is no entry→objects read, and effect intents live in the journal, not the log. |
| V15 | Dependencies | `sed -n 1,10p go.mod` | Direct requirements are `ailang v0.47.2` and `modernc.org/sqlite v1.54.0`. The design adds no dependency. `net/http` is imported in production only by `host/{broker,daemon,pinfetch,projection}`. |
| V16 | AG-UI source pinned | `curl` of `ag-ui-protocol/ag-ui` at commit `903a9ab9a154a164e61c216d8206a760f44780dc` (2026-10-07) | `spec/1.0/schema.json`, 94 359 bytes, sha256 `4b5c9322…48e71a`, `$id https://ag-ui.com/spec/1.0/schema.json`, MIT licence. `@ag-ui/core` is 1.0.2. Banked under `~/.ailang/state/world-iter243-design/`. |
| V17 | Event shapes we rely on | `python3 -I` over `$defs` | `BaseEvent` requires `type` only (`timestamp` is optional). `RUN_STARTED` requires `type,threadId,runId` and allows `protocolVersion`, `parentRunId`, `input`. `RUN_FINISHED` requires `type,threadId,runId` and allows `result` (any JSON), `outcome`, `usage`. **`RUN_ERROR` requires `type,message` and allows only `code`, `usage`.** `STATE_SNAPSHOT` requires `snapshot` (any JSON). `STATE_DELTA` requires `delta` (RFC 6902). `CUSTOM` requires `name` and `value`. All events set `unevaluatedProperties:false`. `RunAgentInput` requires `threadId`, `runId` and `messages`, and has `forwardedProps` (any JSON). |
| V18 | SSE framing of the stock client | `client/src/transform/sse.ts:25,54`; `encoder/src/encoder.ts:28-29`; `client/src/agent/http.ts:55,59` | The encoder emits `data: <json>\n\n`. The client splits on `/\n\n/`, and "Non-data fields (event, id, retry) are ignored". `HttpAgent` sends `POST` with `Accept: text/event-stream`. **So a stock client never sends `Last-Event-ID`.** Resume needs an in-band path. `\r\n` framing would break this parser. |
| V19 | "CUSTOM names must be namespaced" | `grep -c -i namespac schema.json` | **0** (control `grep -c 'RFC 6902'` returns 11). The pinned spec does not require namespacing. `world.` is our own convention (D3), not an AG-UI rule. Whether the prose docs require it is **UNVERIFIED**: the docs site was not fetched. |
| V20 | Workbench CSP unchanged by this row | `sed -n 20p host/daemon/workbench.go` | `default-src 'none'; style-src 'unsafe-inline'; …`. This is row 149's change, not this row's. |
| V21 | Single-entry route and its bytes | `sed -n 580,602p host/daemon/handlers.go`; `sed -n 129,133p`; probe `encprobe_main.go.txt` | `handleLogEntry` ends `writeJSON(w, http.StatusOK, logJSON(entry))` (`:601`), a bare object. `writeJSON` = `json.NewEncoder(w).Encode(value)` (`:132`), and GOROOT `encoding/json/stream.go:200` says "followed by a newline character". The probe (default build: `stream.go` is `//go:build !goexperiment.jsonv2`, `GOEXPERIMENT` empty) gives `encode==marshal+\n: true`, with both HTML-escaping `<>&` identically. **So the route body is exactly `json.Marshal(logJSON(e))` followed by one `\n`.** |
| V22 | The write window opens before the handler, and the body read eats it | GOROOT `net/http/server.go:980-989` | `readRequest` sets `t0 := time.Now()` and `wholeReqDeadline = t0.Add(ReadTimeout)`, then `defer c.rwc.SetWriteDeadline(time.Now().Add(WriteTimeout))`. The 30 s write window therefore starts when the headers are parsed, **before** the handler runs. The body is read by the handler afterwards and is bounded only by the same 30 s whole-request read deadline. A slow body consumes the run budget. Found by quorum r1 (oc-kimi-k3), controller-measured, re-read here. |
| V23 | Stock clients send state back | `client/src/agent/agent.ts` at `903a9ab9` (sha256 `2e8c352c…5a162`, banked `~/.ailang/state/world-iter243/agui-src/`): `:238`, `:562-564`, `:578-607` | The constructor sets `this.state = structuredClone_(initialState ?? {})`. Applied `STATE_SNAPSHOT`/`STATE_DELTA` results land in `this.state` (`if (event.state !== undefined) { this.state = event.state; …`). `prepareRunAgentInput` sends `state: structuredClone_(this.state)` on **every** run. A first run sends `{}`. |
| V24 | A stalled request body gets zero bytes, on `/v1/commit` too | `r3_daemon_probes_test.go.txt` `TestZZProbeStalledBody` (run in-package, then removed): raw TCP `POST` with `Content-Length: 500` and no body, scratch server `ReadTimeout = WriteTimeout = 1 s` | Real `d.Handler()` `POST /v1/commit`: `0 bytes "" err=<nil> after 1s` (clean EOF, not a reset). The silent-arm handler gives the same result. **A handler that tries a typed 408 after the read error also gets `0 bytes`.** Control, with the body sent: `75 bytes "HTTP/1.1 200 OK"`. Parity with `/v1/commit` **holds**. |
| V25 | `aguiWriteMargin` against the measured write cost | `TestZZProbeMarginPageWrite`: 100 realistic entries (`logJSON` of `testCommit`), CUSTOM + `id:` + STATE_DELTA frames, one `Flush` per entry pair, loopback reading client, n = 300 pages | Page = 79 474 bytes. **median 0.82 ms, p99 3.59 ms, max 4.83 ms.** The 2 s margin is more than 400× the measured max. It is kept at 2 s because it also absorbs δ (V22) and scheduler or GC stalls, which this probe cannot bound. |
| V26 | Gapped commits land end-to-end over REST | `TestZZProbeGapEndToEnd`; path re-read: `handlers.go:856` (client `entryIndex` passed unchanged), `daemon.go:639` (`commits: s`, the `*store.Store` itself), `store.go:1142-1152` (the only ordering check is the `ObservedHead` CAS) | `POST /v1/commit` statuses `idx0=200 idx1=200 idx5(gap)=200`. `GET /v1/log?from=0` returns `2 items [0 1]`. `GET /v1/log/5` returns `200`. **V10's UNVERIFIED is discharged.** |

**Draft claims false or stale at HEAD:** line numbers (V1, V2); "no store change" (V10, V26);
"a broadcaster misses second-process writes" (V12); a 30-min connection (V5); session-filtered
reads (V3); "inherits R1 by construction" (V2); "namespaced as AG-UI requires" (V19);
`Last-Event-ID` resume for stock clients (V18); decision/effect mappings (V14). Quorum r1 fixed this
doc's r0: the unanchored clock (V22) and the ignored `state` (V23).

## 3. Decisions, with alternatives

### D1 — A run is bounded inside the frozen write window; resume is the normal way to keep watching (REVISES draft D5)

One `POST /agui/` is one AG-UI run of at most **`aguiRunBudget = writeTimeout − readDeadline −
aguiWriteMargin` = 30 s − 10 s − 2 s = 18 s**. After the budget, no new store read starts.
The run ends with `RUN_FINISHED` whose `result` is `{"lastIndex": k}`. The client re-POSTs with
cursor `k`. A test pins the derivation, not just the value (M3).

**The clock is anchored at handler entry (r1 fix, V22).** `t0` is the handler's first statement.
The write deadline `W` was set at header parse, before dispatch, so `W = t0 + 30 s − δ`. Here δ is
the in-process dispatch gap: the mux plus a middleware that does no I/O on unprotected routes
(V2). δ is microseconds, and the 2 s margin absorbs it. **The body read counts against the
budget.**

- The body must be fully read by `t0 + aguiBodyBound` (**2 s**). A counting reader checks the
  clock on every `Read` return. Once past the bound it stops reading, and the handler answers
  **408 `SlowBody`** before any stream byte. Any `Read` that returns by `t0 + 25 s` yields that
  typed refusal, which is a small JSON write at ≤ `t0 + 25 s` and leaves ≥ 5 s − δ of `W`.
- A `Read` that returns later (a client that sends nothing until the server's whole-request
  read deadline fails it) gets **no bytes at all**, and the connection gets a clean EOF. This is
  measured (V24), with parity to `POST /v1/commit`. The handler writes nothing.
- If the body is done by `t0 + 2 s`, the stream starts. No read starts after `t0 + 18 s`. Each
  read is bounded by 10 s, so the final frame is written by `t0 + 28 s`, which leaves 2 s − δ.
  The 2 s margin is more than 400× the measured page write (V25). No `ResponseController` or
  `SetReadDeadline` is used (V6).
- **Invariant, conditional on a peer that keeps reading (r3):** the response is a complete
  `APIError`, a stream ending in `RUN_FINISHED`/`RUN_ERROR`, or zero bytes. The server **cannot**
  stop a stalled *reader* being cut by `WriteTimeout`, possibly mid-frame. Severance is safe for
  recovery under two tested assumptions. **(a)** The stock splitter (V18) drops an unterminated
  trailing partial frame, so `state.lastIndex` never advances past a fully delivered entry
  (AC3.12). **(b)** Every `STATE_DELTA` op is `replace`, so replaying an entry after severance is
  idempotent (AC1.5).

*Alternatives (rejected):* (a) `ResponseController` deadline extension relaxes policed D7 (V6);
(b) the draft's 30-min budget cannot work (V5); (c) a second `http.Server` without `WriteTimeout`
is a D7 exemption. Cost: a reconnect every 18 s, lossless because resume is exact (D4); AG-UI runs
are finite by design (V17).

### D2 — Events come from the log only; v1 emits exactly six event types (REVISES draft D3 mapping)

The six types are `RUN_STARTED`, `STATE_SNAPSHOT`, `CUSTOM`, `STATE_DELTA`, `RUN_FINISHED` and `RUN_ERROR`. For each committed entry `i`, in index order, the run emits two frames:

```
data: {"type":"CUSTOM","name":"world.entry.committed","value":<logJSON(entry i)>}

id: <i>
data: {"type":"STATE_DELTA","delta":[{"op":"replace","path":"/lastIndex","value":<i>},{"op":"replace","path":"/logHead","value":"<entryHash>"}]}

```

The run opens with `RUN_STARTED` `{threadId, runId, protocolVersion:"1.0"}` (echoing the input)
and `STATE_SNAPSHOT` `{"snapshot":{"schema":"world/agui-state/v1","lastIndex":n,"logHead":<hash of
entry n, or null when n = -1>}}`. Here `n` is the resume cursor, not the current head. That
makes the whole run a function of `(input, entry n, entries > n)`.

- The **`value` is byte-identical to `GET /v1/log/{i}`'s body** (`logJSON`, V9). The stream
  therefore exposes exactly what an unauthenticated GET already exposes. This is the
  no-new-exposure proof, and it is a test (M3), not a claim.
- **`id:` sits on the entry's last frame only.** `Last-Event-ID: i` therefore means entry `i`
  was fully delivered.
- **No `timestamp`** (optional, V17): it would break byte-identity. There are no text, tool-call,
  decision or effect events. Decision and effect events have no log-derivable source (V14). Row
  151 adds `world.view.a2ui` through the name table. Head world ref and revision are left out of
  the state, because no entry→world read exists. Clients read `GET /v1/head`.

*Alternative:* a homegrown `GET /v1/events` dialect. It violates §3.7 / clause 6 and is
unusable by AG-UI clients.

### D3 — Route `POST /agui/`, AG-UI request in, SSE out (RE-VALIDATES draft D3)

- **Route.** `POST /agui/` is additive and outside the frozen `/v1` table, the same pattern as
  `/a2a/` and `/mcp/`. `GET` would admit `EventSource`, but the stock AG-UI client only POSTs
  (V18), and row 149 uses `fetch`. A `GET` on `/agui/` is a mux 405.
- **Input.** Body ≤ **64 KiB** (`http.MaxBytesReader`, 413 `PayloadTooLarge`). The body is JSON
  with `threadId` (string), `runId` (string) and `messages` (array), as V17 requires. `messages`,
  `tools`, `context` and `forwardedProps` are accepted and ignored. `state` is the resume
  cursor (D4). `threadId` is opaque and echoed. The
  store has one world, so it names nothing (REVISES the draft's "threadId names the world").
- **Response.** `200`, `Content-Type: text/event-stream`, `Cache-Control: no-store`,
  `X-Content-Type-Options: nosniff`. Frames use `\n` only, and there is one `Flush` per entry pair.
- **Error split.** Anything wrong **before the first byte** (bad JSON, missing field, slow body,
  bad or conflicting cursor, unknown cursor index, cap reached) is a normal `APIError` HTTP response.
  The stock client surfaces `!response.ok` (V18). Anything after `RUN_STARTED` becomes `RUN_ERROR`
  `{message, code}`, with `code` taken from the daemon's classes (`Timeout`, `Internal`) and a
  **constant** message. Internal detail goes only to `writeInternalError`'s log path.
  `RUN_ERROR` never carries `threadId` (V17: `unevaluatedProperties:false`).

### D4 — Resume by entry index, read from the standard `state` field (REVISES draft D2; r1 fix)

A stock client sends back the state it built from our snapshot and deltas, on every run (V23).
The cursor is therefore read from that standard field, with no custom client code:

- `state.schema == "world/agui-state/v1"`: then `after = state.lastIndex`, which must be an
  integer ≥ −1, or the request gets 400 `BadRequest`. `state.logHead` is not trusted, because the
  snapshot re-derives it from the store.
- `state` absent, or without our schema (a first run sends `{}`): from genesis, `after = −1`.
- `Last-Event-ID` (standard SSE) is kept for `curl` and hand-rolled clients. If both are present
  and differ, the request gets 400. `forwardedProps` is **not** a cursor (dropped from r0: it was
  a private dialect).
- A non-negative `after` must name an **existing** entry (one `GetLogEntry`, which also feeds
  the snapshot). Otherwise the answer is 404 `NotFound` "unknown cursor", never a silent empty
  stream. The run then delivers exactly the entries with index `> after`, in order.

### D5 — Tail by keyset read, one statement per tick; one small store read is added (REVISES draft D2 "no store change")

Add `Store.LogEntriesAfter(ctx, after int64, limit int) ([]LogEntry, error)`. It is
`SELECT … FROM log_entries WHERE entry_index > ? ORDER BY entry_index LIMIT ?` on the INTEGER
PRIMARY KEY, with `1 ≤ limit ≤ 500` (`InvalidLimitError`) and `after ≥ -1`, and it reuses
`GetLogEntry`'s row parse. It joins the `readStore` seam. The run loop works as follows:
`page := LogEntriesAfter(cursor, 100)`. Emit the page. If the page is full, read again
immediately. If it is short, sleep `aguiTick = 250 ms`, unless the run is stopping.

- **Why a store read:** V10. Composing `GetLogEntry(cursor+1)` stalls silently forever at a gap,
  which is the draft's own "never a silent stall" failure mode. A keyset read crosses gaps by
  construction, which V26 confirms is reachable over REST. Even with no gap, it is justified
  because one statement replaces an N+1 per-index poll. The frozen `GET /v1/log` is **not**
  changed. Its first-gap stop is recorded as a finding (§9), not fixed here.
- **Why poll, not a broadcaster:** all writes are in-process (V12), so the only reason left is
  coverage. A hook must reach every commit path (REST, MCP, A2A, coordinator, future ones). A poll
  sees every row by construction, at ≈ 10 µs per idle tick (V11): 16 × 4 × 10 µs ≈ 0.6 ms per
  second at the cap, with ≤ 250 ms latency. Each read uses `context.WithTimeout(r.Context(),
  d.readDeadline)`, so a disconnect cancels it.

### D6 — Bounds (REVISES draft D5 numbers)

| Bound | Value | Enforced by |
|---|---|---|
| Run lifetime | 18 s from **handler entry**, derived (D1) | `d.aguiBudget` field set from the constant (test-shrinkable, like `drainTimeout`) |
| Global concurrent runs | 16 | buffered-channel semaphore. Over cap → 503 class `StreamLimit` with `Retry-After: 1`, before any byte |
| Entries per read | 100 (≤ the 500 store max) | `LogEntriesAfter` limit |
| Idle tick | 250 ms | constant |
| Request body | 64 KiB, fully read by `t0 + 2 s` | `MaxBytesReader` plus the clock-checking reader (D1). Over the bound → 408 `SlowBody` |
| Each store read | `readDeadline` (10 s) | `readCtx`-equivalent |
| Slow or stalled client | ≤ `writeTimeout` (30 s) total | the unchanged D7 server deadline (V5) |
| Shutdown | the run ends with `RUN_FINISHED` ≤ 1 tick after `Shutdown()` starts | `srv.RegisterOnShutdown` closes `d.aguiStop` (V7) |

### D7 — Read posture equals `GET /v1/log`, pinned by a test, not by "construction" (REVISES draft D4)

`isProtected(POST /agui/)` is false today, the same as every `/v1` GET (V2). The stream ignores
any `Authorization` header, as `GET /v1/log` does. There is **no session filtering** in this
row (V3: there is no filter to reuse, and adding one is an authority feature, not a transport).
A test asserts `isProtected(POST /agui/) == isProtected(GET /v1/log)`. When the residual-R1
row flips GET protection, that test goes red and forces the `/agui/` decision to be made
deliberately. Cross-origin: the route sets no CORS headers, and a JSON `POST` is preflighted, so
other origins cannot read the stream. Nothing here is actionable. The stream is read-only.

## 4. Pure core vs effect boundary

- **`host/agui` (new, pure; imports neither `net/http` nor `host/store`):** `Entry{Index, EntryHash,
  Value json.RawMessage}`, `Run{ThreadID, RunID}`, the pure functions `Start(run, after, snap
  *Entry)`, `EntryFrames(e)`, `Finished(run, lastIndex)` and `Error(code, msg)`, each returning
  `[]byte`, the closed `CUSTOM` name table, `ProtocolVersion = "1.0"` and `StateSchema`. It is
  deterministic: fixed field order, `json.Marshal`, no clock, no map iteration, `\n` framing.
- **`host/daemon/agui.go` (effects):** parsing, the cursor, the semaphore, the poll loop, `Flush`,
  shutdown and the route. It converts `store.LogEntry` → `agui.Entry` with `Value =
  json.Marshal(logJSON(e))`.
- **Pinning:** the golden `host/agui/testdata/stream_fixture.golden` (entries 0, 1, 2, 5, 6,
  including the gap), determinism (encoder and daemon run twice), schema conformance, and the
  `\n\n` split.
- **S1/S3:** no `.ail`. This is host-boundary wire formatting like `/mcp/` and `/a2a/`. *Why not
  a package:* an AILANG package cannot own a route or a socket write, and the semantics (which
  entries exist) are already the kernel log. There is no Z3-encodable invariant beyond the budget
  arithmetic, which is pinned by a Go derivation test (M-e).

## 5. Milestones (tests first; each ≤ 1 executor day)

**M1 — `host/agui` pure encoder.**
- AC1.1 `go test ./host/agui -run 'TestGolden|TestDeterministic'`: the fixture bytes equal the
  golden, and two encodes are equal.
- AC1.2 `TestEventsConformToPinnedSchema`: it loads `host/agui/testdata/agui-1.0-schema.json`
  (V16 copy, sha256 asserted in the test). For every emitted event, `type` is in `EventType.enum`
  and equals the def's `const`, `required ⊆ keys`, and `keys ⊆ def.properties ∪ BaseEvent ∪
  Attributable`. Every delta op has `op`, `path` and `value`. Null control: the test fails if it
  checked 0 events or found 0 of the 6 expected defs.
- AC1.3 `TestFramesSplitLikeStockClient`: splitting the bytes on `"\n\n"` and taking `data:`
  lines (V18's algorithm) yields exactly the event list. Each frame contains no `\r`.
- AC1.5 `TestDeltaIdempotent`: every emitted delta op is `replace`. For each fixture entry,
  applying its delta twice to the post-entry state yields that state unchanged.
- AC1.4 `go vet ./host/agui` passes, and `go list -deps ./host/agui` contains neither `net/http`
  nor `host/store`.

**M2 — `Store.LogEntriesAfter` + seam.**
- AC2.1 `TestLogEntriesAfterCrossesGap` (entries 0,1,5) checks that after −1 returns
  [0,1,5], after 1 returns [5], and after 5 returns [].
- AC2.2 Limit bounds: 0 and 501 give `InvalidLimitError`, and `after < -1` is refused. A context
  without a deadline gives `ErrNoDeadline`. A quarantined store is refused.
- AC2.3 `TestLogEntriesAfterMatchesGetLogEntry` checks that each row deep-equals `GetLogEntry`
  for the same index.
- AC2.4 `readStore` gains the method. `go test ./host/daemon ./host/store` passes.

**M3 — `POST /agui/` handler.**
- AC3.1 `TestAGUIGoldenRun`: a daemon over the fixture store, with `aguiBudget` shrunk to
  300 ms. The body equals the golden plus `RUN_FINISHED`, and it is identical on a second run.
- AC3.2 `TestAGUIResumeExact`: for each k in {−1,0,1,2,5,6}, `Last-Event-ID: k` yields exactly
  the entries > k. `state:{schema,lastIndex:k}` gives the same bytes. `state:{}` means genesis.
  A malformed `lastIndex` (`"2"`, `-2`, `1.5`) gives 400, a state/header conflict gives 400, and
  an unknown k (3) gives 404.
- AC3.2b `TestAGUIStockClientStateResume`: take run 1's exact frames, parse them with the
  `\n\n` splitter, and apply `STATE_SNAPSHOT` plus each `STATE_DELTA` with a test-local RFC 6902
  applier (`replace`/`add` only; any other op fails the test). Commit 2 more entries, POST the
  resulting state verbatim as `state`, and get exactly the entries > that state's `lastIndex`.
- AC3.11 `TestAGUISlowBodyNeverTruncates`: bounds are shrunk (`aguiBodyBound` 100 ms, budget
  300 ms; both are daemon fields like `drainTimeout`). A body dripping 1 byte every 30 ms that
  ends after the bound gives 408 `SlowBody` with no `data:` byte. A body that ends at 80 ms gives
  a stream whose `RUN_FINISHED` arrives by `t0 + 300 ms + 1 tick + 50 ms` slack. With a client
  that **keeps reading**, every response body is either one parseable `APIError` or ends with a
  `RUN_FINISHED`/`RUN_ERROR` frame. The slow-reader case is explicitly excluded (AC3.13 owns it).
- AC3.12 `TestAGUISeveranceResumable`: truncate the golden run at **every byte offset** and split
  it with the V18 algorithm. The last complete `STATE_DELTA`'s `lastIndex` (or the snapshot's) is
  a valid cursor. Resuming from it and concatenating the entry frames reproduces the golden
  entries with no loss and no gap.
- AC3.13 `TestAGUISlowReaderCutIsResumable`: a scratch server wraps `d.Handler()` with
  `WriteTimeout` 200 ms over a fixture of about 5 000 entries (more than the socket buffers). The
  client reads 1 KiB, then stops. The oracle: the connection is cut, and the last COMPLETE
  frame's `lastIndex` resumes per AC3.12.
- AC3.14 `TestAGUIRESTGapCrossed` (**executor runs it first**, V26): gapped commits 0, 1, 5 via
  `POST /v1/commit` land. `GET /v1/log?from=0` returns [0 1]. `/agui/` from genesis delivers
  0, 1, 5.
- AC3.3 `TestAGUIEntryValueEqualsLogRoute`: for each entry, the `CUSTOM` value bytes followed by
  `"\n"` equal the body of `GET /v1/log/{i}` byte-for-byte (V21). This is not a trimmed or
  semantic comparison.
- AC3.4 `TestAGUISeesDirectStoreCommit`: mid-run, `d.store.Commit` (which bypasses every
  handler) appears within 2 ticks. So does a `POST /v1/commit`.
- AC3.5 `TestAGUIGlobalCap`: 16 runs held open, the 17th gets 503 `StreamLimit`, and after one
  closes the next is admitted.
- AC3.6 `TestAGUIBudgetDerivation`: `aguiRunBudget == writeTimeout - readDeadline -
  aguiWriteMargin`, `aguiRunBudget > 0`, and a new daemon's `d.aguiBudget == aguiRunBudget`.
- AC3.7 `TestAGUIShutdownEndsRuns`: with a run open, `d.Shutdown()` returns nil in < 2 s, and the
  body ends with `RUN_FINISHED`.
- AC3.8 `TestAGUIClientDisconnectReturns`: cancelling the request context makes the handler
  return within 1 tick. The goroutine count returns to baseline, run under `-race`.
- AC3.9 `TestAGUIReadPostureMatchesLog` (D7 equality) and `TestAGUIRejectsGET` (405).
  `TestAGUIInputErrors` covers bad JSON, a missing `runId`, `messages` that is not an array, and
  a body over 64 KiB (413).
- AC3.10 `TestAGUIStoreErrorIsRunError`: a seam wrapper fails mid-run, and the stream ends with
  `RUN_ERROR` `code:"Internal"` and a constant message that contains no store text.
- Gate: `go build ./... && go test ./... && go test -race ./host/daemon ./host/agui ./host/store`.

**M4 — S7 usage surface.** QUICKSTART gains "Watch the world live": `curl -N -X POST -H
'Content-Type: application/json' -d '{"threadId":"t","runId":"r","messages":[]}'
http://127.0.0.1:7644/agui/`, the first three frames, a `Last-Event-ID` resume, and the 18 s
re-POST contract. It is executed verbatim, with the transcript in the PR. The route doc comment is
updated. AC4.1 is a QUICKSTART grep, an instrument-health control only.

## 6. Failure modes

| Failure | Behaviour |
|---|---|
| Client disconnects | `r.Context()` cancels the in-flight read. The handler returns and the semaphore is released (AC3.8) |
| Slow or stalled request body | 408 `SlowBody` before any byte, or zero bytes if no `Read` returns by `t0 + 25 s`. Never a truncated stream (AC3.11) |
| Client stops reading | Cut by the unchanged 30 s `WriteTimeout`, possibly mid-frame. The stock splitter drops the partial frame, and the client resumes from the last complete `lastIndex` (AC3.12, AC3.13). There is no buffering beyond one entry pair |
| Store read error or deadline | `RUN_ERROR` (`Internal`/`Timeout`), constant message, stream closed (AC3.10) |
| Log index gap | Crossed by the keyset read (AC2.1, AC3.2 with k=2) |
| Cursor beyond or absent from the log | 404 before any byte (AC3.2) |
| Daemon shutdown with runs open | `RUN_FINISHED` within 1 tick. The drain is not held to its 10 s expiry (AC3.7) |
| Cap reached | 503 `StreamLimit`, `Retry-After: 1` (AC3.5) |
| AG-UI spec bump | The pinned schema copy and sha make it a deliberate one-file update (AC1.2) |
| Quarantined store | The store refuses reads, which surfaces as `RUN_ERROR Internal` or a pre-stream 500 |

## 7. Load-bearing mutations (S6)

| Mutant | Test that must go red |
|---|---|
| M-a cursor advances by `+2` / skips the last row of a full page | `TestAGUIGoldenRun`, `TestAGUIResumeExact` |
| M-b encoder adds `timestamp` | `TestDeterministic` (encoder), `TestAGUIGoldenRun` |
| M-c tail by `GetLogEntry(cursor+1)` instead of the keyset read | `TestAGUIResumeExact` (k=2 times out with no entry 5) |
| M-d `LogEntriesAfter` uses `>=` | `TestLogEntriesAfterCrossesGap` |
| M-e budget = `writeTimeout` (no derivation) | `TestAGUIBudgetDerivation` |
| M-f semaphore check removed | `TestAGUIGlobalCap` |
| M-g `RegisterOnShutdown` hook removed | `TestAGUIShutdownEndsRuns` |
| M-h read ctx from `context.Background()` | `TestAGUIClientDisconnectReturns` |
| M-i `CUSTOM` value adds a field (for example the transition payload) | `TestAGUIEntryValueEqualsLogRoute` |
| M-j `RUN_ERROR` carries `threadId` | `TestEventsConformToPinnedSchema` |
| M-k frames end `\r\n\r\n` | `TestFramesSplitLikeStockClient` |
| M-l `id:` moved to the `CUSTOM` frame | `TestGolden` |
| M-m `/agui/` added to `isProtected` | `TestAGUIReadPostureMatchesLog` |
| M-n conflicting cursors: header wins silently | `TestAGUIResumeExact` (400 arm) |
| M-q run clock anchored **after** the body read | `TestAGUISlowBodyNeverTruncates` (`RUN_FINISHED` late) |
| M-r body bound not enforced (stream starts after a slow body) | `TestAGUISlowBodyNeverTruncates` (408 arm) |
| M-s `state` ignored, cursor only from `Last-Event-ID` | `TestAGUIStockClientStateResume` |
| M-t a schema-less `state` treated as an error, not genesis | `TestAGUIResumeExact` (`state:{}` arm) |
| M-u a delta op that is not `replace` (e.g. `add` to `/entries/-`) | `TestDeltaIdempotent` |
| M-v `id:` emitted before the entry's last frame | `TestAGUISeveranceResumable`, `TestAGUISlowReaderCutIsResumable` |
| M-w the `/agui/` tail rebuilt on the dense `/v1/log` loop | `TestAGUIRESTGapCrossed` |
| M-o tail only on an in-handler commit signal (no poll) | `TestAGUISeesDirectStoreCommit` |
| M-p store error text copied into `RUN_ERROR.message` | `TestAGUIStoreErrorIsRunError` |

## 8. Conflict surface

- **New:** `host/agui/{agui.go,agui_test.go,testdata/*}`, `host/daemon/agui.go`,
  `host/daemon/agui_test.go`, `host/store/log_after.go` and its test.
- **Edited:** `host/daemon/daemon.go` (one `HandleFunc` line, the `readStore` method, the
  `aguiBudget`/`aguiStop`/semaphore fields, `RegisterOnShutdown` in `newServer`'s caller, and the
  route doc comment), and `docs/QUICKSTART.md`.
- **Must stay unchanged** (no hunk; evaluator checks `git diff --stat`): the ten `/v1` patterns and
  handlers, `isProtected`'s body, every D7 constant and `TestBoundedWaitsAndBodyLimit`,
  `workbenchCSP`, `host/projection/*`, `tools/launchd/*`. No open critical-path row names these
  files at `d6334c2`; re-check at merge.

## 9. Findings to file (not fixed here)

1. **The store and REST commit path accept a gapped `entryIndex` and an unlinked `prevEntryHash`**
   (V10); `GET /v1/log` then stops at the gap. A log-integrity ledger question / queue candidate.
2. `w-world-live-surface.md` should point here for row 148; 149/151 consume §1's interface.

## 10. Non-goals

- Rows 149–151: the CSP change, `workbench.js`, demo seed, captures, A2UI views.
- `world.decision.*`/`world.effect.*` (no log source, V14); session filtering; closing residual R1;
  text/tool-call/reasoning/`MESSAGES_SNAPSHOT` events; protobuf; CORS; any change to `GET /v1/log`.
- A "start at tail" cursor. Row 149 server-renders the current last index and passes it as
  `after`.

## 11. Quorum verification log

**Round 1: BLOCKED.** 3 of 4 seats present. `gpt6-1-sol` was absent (unreachable). The controller
measured every objection premise first-party before this revision (V21–V23).

| Seat | Verdict | Objection (one line) | Premise | Resolution |
|---|---|---|---|---|
| oc-kimi-k3 | BLOCK | The run clock is unanchored: body-read time eats the 18 s budget, and a slow body can push the final frame past `WriteTimeout` (a silent cut) | REAL (V22) | D1 anchors at handler entry with a 2 s body bound, typed 408 and the write-window arithmetic. AC3.11, M-q, M-r |
| gemini-3-1-pro | BLOCK | Ignoring `state` and inventing `forwardedProps.world.after` is a dialect; stock clients send `state` back every run | REAL (V23) | D4 reads the cursor from `state`. `forwardedProps` cursor dropped. AC3.2/3.2b, M-s, M-t |
| oc-glm-5-3 | BLOCK | `GET /v1/log/{index}` does not exist, so the "value equals route body" proof is unfounded | FALSE (V1 `:864`, V21) | V1 now names all ten patterns. AC3.3 states the exact bytes compared (`Marshal` + `\n`) |

**Round 2: BLOCKED.** 3 of 4 present (`gpt6-1-sol` absent again). No objection disputed direction.
**r3 is the narrow-refinement carve-out: the reviewers' own fixes are applied verbatim, with no
redesign.**

| Seat | Objection (one line) | r3 fix |
|---|---|---|
| oc-kimi-k3, oc-glm-5-3 | The "never truncated" invariant is unconditional, yet a stalled reader is cut | D1 invariant made conditional. Assumptions (a) and (b) tested by AC1.5, AC3.12, AC3.13. AC3.11 excludes slow readers. M-u, M-v |
| oc-kimi-k3 | The 2 s `aguiWriteMargin` is unmeasured | Probed: max 4.83 ms per 100-entry page (V25). 2 s kept, justified |
| oc-glm-5-3 | The stalled-body zero-byte arm is analysis only | Probed: 0 bytes and a clean EOF, parity with `/v1/commit` (V24) |
| gemini-3-1-pro | V10 is UNVERIFIED end-to-end | Path re-read and the REST gap probe run (V26). AC3.14 runs first. M-w |
