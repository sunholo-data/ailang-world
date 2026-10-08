# Sprint plan — `w-worldd-agui-event-stream`

**Iteration 243 · queue row 148 · clause-6 (+5; release row, D-WORLD-64)**

Design: [w-worldd-agui-event-stream.md](w-worldd-agui-event-stream.md), ratified quorum r3 narrow-refinement carve-out. Base: detached `0006819`; controller baseline `d6334c2` is code-identical except docs. Binding guidance: CLAUDE.md and coding-standards.md S2/S3/S6/S7. Estimate: **~1.8 executor days** (C1 0.3, C2 0.3, C3 1.0, C4 0.2).

Deliver one PR with **four commits C1–C4, one per M1–M4**. The EXECUTOR never runs git (including add/commit/stash/checkout); the controller commits from per-milestone snapshots. Save a reviewable green snapshot, tests-first evidence, command return codes, timings and mutation ledger at each boundary. This planner changes exactly this plan and the root JSON; implementation is executor work.

## Environment and measured baseline

```sh
export AILANG_BIN=$HOME/.pinned-ailang/ailang
export PATH=/opt/homebrew/bin:$PATH
```

Pinned gate binary: AILANG v0.41.0, commit 24ee108; Go 1.26.6 darwin/arm64. No .ail changes: the route/socket writes belong to the Go host boundary; the semantic source is the existing kernel log. No new dependencies. An AILANG package cannot own the HTTP route or socket writes.

| Observation | Evidence |
|---|---|
| tree | 0006819; controller baseline d6334c2 is code-identical except docs |
| source | Controller ran both full gates outside sandbox on pristine d6334c2; planner read both logs. |
| logs | ['~/.ailang/state/world-iter243/base_ail.log', '~/.ailang/state/world-iter243/base_go.log'] |
| ./scripts/verify_ail.sh | controller rc=0; 22 modules; 16/16 world identities; 138 execargs tests; 40 required named tests; 444 se-tools tests across 9 modules; world package 9/9 steps; PUB011 uncontracted_exports=3 |
| ./scripts/verify_go.sh | controller rc=0; Go 1.26.6 darwin/arm64; AILANG v0.41.0 (24ee108); 2133 tracked files, 0 binary blobs; decision ledger 55 rows; evidence 37 exact tests; plain daemon 179.509s, store 13.146s; race daemon 199.241s, store 26.767s, broker 145.978s, verifygate 169.697s; intentional race-control exit 66 with 2 races is expected |
| go vet ./... | planner rc=0; empty output at 0006819 |
| go test -count=1 ./host/store | planner rc=0; 11.297s; no socket refusal |
| sandbox_rule | Loopback bind: operation not permitted is UNINFORMATIVE UNDER SANDBOX, never pass or fail. No planner full-gate rerun or socket run; full-gate wall time not recorded in these logs. |

The intentional known-positive race-control output is expected, not a baseline test failure. The readable logs carry package timings but no reliable full-gate wall duration; do not invent one. **Any loopback `bind: operation not permitted` result is UNINFORMATIVE UNDER SANDBOX, never pass or fail.** Obtain socket evidence on the controller’s socket-capable host.

## Execution rules and protected surface

For each milestone, write every listed test before its implementation and run it on the **pre-milestone code**. Each new oracle must fail there; capture the specific red. Missing new package/API compile errors are build-red proof only; after implementation, mutation kills require the named assertion to fail on compiling code. For existing-route preconditions (especially AC3.14), first prove those preconditions green. No red may be attributed to missing sockets. M4 writes its documentation control before editing docs. Use byte-restored snapshots, never git operations, to reproduce a pre-milestone tree.

**Planner rule: any test that writes outside its own `t.TempDir()` is forbidden.** Put DBs, fixtures, HOME/state/cache roots, fake executables and subprocess outputs there; use `t.Setenv` when necessary before creating the daemon. No seeded writes to `$HOME` or `~/.ailang`. Do not copy the exemplar’s external log-write convention. Read-only pinned binary/schema sources are allowed. Keep executor evidence with controller snapshots; no extra tracked transcript files.

CI runs `-race` with a ~600s Go-leg wall budget. Shrink per-instance run budgets (typically 20–50ms, golden/body tests 300ms), ticks (5–10ms), and body bound (100ms); leave production defaults and frozen D7 unchanged. Hold cap-test streams concurrently; use start/read/cancel channels, bounded contexts, and joins instead of long sleeps. No new daemon test may sleep more than ~2s total. Resume loops close/cancel after the expected entries instead of burning a whole budget per cursor/cut. Severance enumeration is pure and does not open a socket for every byte offset. Report aggregate focused test wall and daemon race wall; never raise the CI cap.

The following must have **no delivered hunk**:

- `GET /v1/health`
- `GET /v1/head`
- `GET /v1/worlds/{ref}`
- `GET /v1/objects/{ref}`
- `GET /v1/objects/by-semantic-id/{name...}`
- `GET /v1/log/{index}`
- `GET /v1/log`
- `GET /v1/registry/{name...}`
- `POST /v1/commit`
- `GET /v1/receipts/{id}`
- `host/daemon/handlers.go (all /v1 handlers unchanged)`
- `host/daemon/daemon.go isProtected body`
- `all D7 constants and TestBoundedWaitsAndBodyLimit in host/daemon/daemon_test.go`
- `host/daemon/workbench.go workbenchCSP`
- `host/projection/*`
- `tools/launchd/*`
- `design_docs/world-mission*.md`

All existing `/v1` route patterns and handlers stay byte-unchanged; `host/daemon/handlers.go` is excluded. `daemon.go` changes only the section-8 seam, additive route, new AGUI fields/initialization, shutdown-hook registration and route comment. Existing `newServer` transport settings remain frozen. No `ResponseController`, `SetWriteDeadline` or `SetReadDeadline` in new sources. Scope audit compares controller snapshot with base; read-only `git diff` is controller work. M-m is an isolated temporary mutant, never a deliverable hunk.

## Numbered planner decisions

**PD1.** Golden boundary is underspecified (AC3.1 says golden plus RUN_FINISHED). Store stream_fixture.golden as Start + entries 0,1,2,5,6, excluding terminal; TestGolden pins that prefix and separately pins Finished and Error bytes. TestAGUIGoldenRun appends the independently expected RUN_FINISHED(lastIndex=6). TestDeterministic also asserts timestamp absent in every frame: a fixed added timestamp need not break two-encode equality (M-b).

**PD2.** AC3.12/3.13 using only delta state cannot kill early id placement (M-v). Keep the state-cursor recovery oracle, and additionally track complete SSE id-bearing frames at every cut; require each id on its matching STATE_DELTA and last complete id == last fully delivered delta index. Resume through Last-Event-ID too. This makes M-v red even if state-only recovery would succeed.

**PD3.** D6 fixes aguiTick at 250ms but the brief requires shrunk ticks. Add an unexported per-daemon aguiTick field initialized from the unchanged new constant (250ms); use 5–10ms in most tests, 10ms for slow-body tests. Add aguiBodyBound and aguiBudget fields as the design requests. Pin defaults in TestAGUIBudgetDerivation; no D7 constant changes or package-global mutable test seams.

**PD4.** D5 says reuse GetLogEntry row parse, but section 8 does not permit store.go refactoring. In log_after.go faithfully reuse the same column order, hashref.Parse operations, header construction and error handling algorithm locally. Do not change store.go or call GetLogEntry N times. TestLogEntriesAfterMatchesGetLogEntry deep-equality and malformed-row arms bind the local parser to the existing representation. A shared-helper refactor is deferred.

**PD5.** M4 has only a grep instrument-health AC, so red-first needs an explicit oracle. Before docs edits add TestAGUIQuickstartControl in host/daemon/agui_test.go: require the live-watch heading, runnable exact POST payload, first-three-frame example, Last-Event-ID example and 18s re-POST text. It must fail on C3 docs. The grep and this control are documentation instrument-health checks, not load-bearing transport proof; execute examples verbatim separately.

**PD6.** D4 header/state conflict wording must not make header-only resume conflict with implicit genesis. Compare header only against an explicitly supplied recognized-schema state cursor; absent or schema-less state is genesis only when no header is supplied. Test header-only, empty state plus header, matching/conflicting recognized state, malformed/overflow header and untrusted logHead.

**PD7.** Shutdown may occur during a 10s in-flight store read, so merely selecting aguiStop during idle wait does not prove within-one-tick shutdown. Use a cancellable child of r.Context for run reads, cancel it on aguiStop, join the watcher on all exits, and treat shutdown cancellation as RUN_FINISHED rather than RUN_ERROR. AC3.7 includes a blocking read seam and idle stream, AC3.8 tests actual read-context cancellation (kills M-h). Keep server drain and readDeadline unchanged.

Schema source inspected read-only: `~/.ailang/state/world-iter243-design/agui_schema.json`, upstream `903a9ab9a154a164e61c216d8206a760f44780dc`, 94,359 bytes, SHA-256 `4b5c93226838a0e72d88e6c5df20633c686c49fcb75d9be2815c6cbf9e48e71a`. Executor copies these exact bytes into the named testdata path (preserve MIT provenance in test comments); no live network lookup needed. Existing readStore test wrappers were inspected: they embed `*store.Store` or `readStore`; the new method is promoted. Override blocking/error/recording behavior explicitly in new AGUI seams rather than modifying unrelated tests.

## C1 / M1 — pure host/agui encoder

**Exact file paths:**

- `host/agui/agui.go`
- `host/agui/agui_test.go`
- `host/agui/testdata/agui-1.0-schema.json`
- `host/agui/testdata/stream_fixture.golden`

**Tests FIRST; every new test must be red on the pre-milestone snapshot:**

- **AC1.1 — `TestGolden, TestDeterministic`:** Fixture 0,1,2,5,6 including gap; compare exact LF bytes with golden, independently pinned terminal/error bytes, and two encodes. Assert no timestamp (PD1); closed CUSTOM table allows only world.entry.committed, protocolVersion 1.0, state schema world/agui-state/v1.
- **AC1.2 — `TestEventsConformToPinnedSchema`:** Read pinned schema; assert exact SHA-256, six nonzero checked event types/defs, EventType enum and def const, required keys, allowed def/BaseEvent/Attributable keys, and delta op/path/value. Include Error so extra threadId reds. No claim of a general JSON Schema validator.
- **AC1.3 — `TestFramesSplitLikeStockClient`:** Use exact LF LF splitter and data-line parser; reject CR; compare all parsed events with independent expected list and drop unterminated tail.
- **AC1.4 — `TestAGUIDependencyBoundary`:** Test invokes go list -deps ./host/agui and refuses net/http and host/store; also go vet ./host/agui. Independently inspect dependency list.
- **AC1.5 — `TestDeltaIdempotent`:** Every delta op is replace; independently apply each delta twice to post-entry state and assert equality, including both lastIndex and logHead.

**Implementation after red proof:** Implement pure Entry/Run, Start/EntryFrames/Finished/Error byte encoders, fixed field order, six event types, closed CUSTOM table and no timestamps. id appears only on STATE_DELTA, last frame of each entry; one LF LF boundary per frame. No HTTP or store import.

**Acceptance commands, in order** (use environment above; every named test must show PASS, no SKIP/zero-match):

```sh
go test -count=1 -v ./host/agui
go vet ./host/agui
go list -deps ./host/agui | awk '/^(net\/http|github.com\/sunholo-data\/ailang-world\/host\/store)$/ {bad=1; print} END {exit bad}'
go vet ./...
```

**Mutation drills on this milestone snapshot:**

| Mutant | Isolated edit | Test that must red / assertion |
|---|---|---|
| M-b | add timestamp, including a constant timestamp | `TestDeterministic`: timestamp-absence and golden bytes (also TestAGUIGoldenRun at C3) |
| M-j | RUN_ERROR carries threadId | `TestEventsConformToPinnedSchema`: allowed properties excludes threadId on Error |
| M-k | CRLF CRLF framing | `TestFramesSplitLikeStockClient`: noCR/exact stock-split event list |
| M-l | move id to CUSTOM frame | `TestGolden`: exact frame bytes |
| M-u | non-replace delta e.g. add /entries/- | `TestDeltaIdempotent`: every op replace and double application equal |

**Controller commit message:** `feat(agui): add deterministic AG-UI 1.0 SSE encoder`

## C2 / M2 — keyset Store.LogEntriesAfter and readStore seam

**Exact file paths:**

- `host/store/log_after.go`
- `host/store/log_after_test.go`
- `host/daemon/daemon.go`
- `host/daemon/agui_test.go`

**Tests FIRST; every new test must be red on the pre-milestone snapshot:**

- **AC2.1 — `TestLogEntriesAfterCrossesGap`:** Store commits 0,1,5: after -1 => [0,1,5], after 1 => [5], after 5 => []; assert ascending exclusive ordering; add >100 entries and limit 1/100/500 paging so full-page loss is observable.
- **AC2.2 — `TestLogEntriesAfterBounds`:** Limits 0,501 => InvalidLimitError; after < -1 refused; missing deadline => ErrNoDeadline; quarantined store refused; expired/canceled context propagates error. Seed and quarantine only a t.TempDir store.
- **AC2.3 — `TestLogEntriesAfterMatchesGetLogEntry`:** Deep-equal each result with GetLogEntry same index; include realistic hashes/header fields and malformed hash parse errors. Verify one ordered SQL query, no per-row getter calls.
- **AC2.4 — `TestAGUIReadStoreSeam`:** Add compile-time *store.Store satisfies readStore assertion plus call LogEntriesAfter via readStore on temporary fixture. Existing fake stores embed *store.Store or readStore and inherit the added method; no unrelated test edits required. New blocking/recording/failing overrides live in agui_test.go.

**Implementation after red proof:** Implement LogEntriesAfter with quarantine/deadline validation, after>=-1 and limits1..500, exclusive ordered SQL keyset read and faithful existing row representation (PD4). Add readStore method in daemon.go and seam test in agui_test.go; no route yet.

**Acceptance commands, in order** (use environment above; every named test must show PASS, no SKIP/zero-match):

```sh
go test -count=1 -v ./host/store -run '^TestLogEntriesAfter'
go test -count=1 -v ./host/daemon -run '^TestAGUIReadStoreSeam$'
go test -count=1 ./host/daemon ./host/store
go vet ./...
```

**Mutation drills on this milestone snapshot:**

| Mutant | Isolated edit | Test that must red / assertion |
|---|---|---|
| M-d | LogEntriesAfter WHERE >= | `TestLogEntriesAfterCrossesGap`: after1 excludes1 |

**Controller commit message:** `feat(store): add bounded keyset log reads for AG-UI`

## C3 / M3 — bounded POST /agui/ event stream

**Exact file paths:**

- `host/daemon/agui.go`
- `host/daemon/agui_test.go`
- `host/daemon/daemon.go`

**Tests FIRST; every new test must be red on the pre-milestone snapshot:**

- **AC3.14 — `TestAGUIRESTGapCrossed`:** FIRST M3 test written and run before any M3 production edit: authenticated POST /v1/commit commits 0,1,5 with valid observed-head CAS, each 200; GET /v1/log?from=0 gives exactly [0,1], GET /v1/log/5 200; POST /agui/ genesis emits exactly 0,1,5. Baseline red is missing /agui/, never failure of the commit precondition.
- **AC3.1 — `TestAGUIGoldenRun`:** Real daemon reads 0,1,2,5,6 fixture, aguiBudget=300ms; byte-exact golden prefix plus independently expected RUN_FINISHED, same input twice identical bytes. Add >100-entry subtest to catch skipped row of full page (M-a).
- **AC3.2 — `TestAGUIResumeExact`:** For k=-1,0,1,2,5,6 compare header-only and state schema/lastIndex exact bytes, only entries >k; state:{} genesis; empty state + header allowed (PD6); malformed lastIndex string 2, -2,1.5 =>400; conflicting explicit cursors =>400; unknown 3 =>404, before any data byte. Include matching cursors, forged logHead ignored, malformed header, beyond-head cursor, forwardedProps ignored.
- **AC3.2b — `TestAGUIStockClientStateResume`:** Parse run 1 bytes exactly like stock splitter; local RFC6902 applier supports replace/add only, rejects all other ops; apply snapshot/deltas, commit two more entries, POST resulting state verbatim without header; receive exactly newer entries.
- **AC3.3 — `TestAGUIEntryValueEqualsLogRoute`:** Retain raw CUSTOM value token bytes using json.RawMessage, append one LF, compare to GET /v1/log/{i} body exactly; no TrimSpace/remarshal/semantic equality. Include <>& escaping control and all fixture entries.
- **AC3.4 — `TestAGUISeesDirectStoreCommit`:** Wait for stream start, commit through d.store.Commit bypassing handlers; see it within two shrunk ticks; repeat POST /v1/commit. Record polls with channel handshake, never rely on arbitrary sleep; no commit signal may be sole source.
- **AC3.5 — `TestAGUIGlobalCap`:** Hold 16 streams with start handshakes (one shared short budget); 17th =>503 APIError StreamLimit and Retry-After:1 before any SSE byte; cancel/join one then next admitted; verify slot cleanup on refusals, store errors and disconnect.
- **AC3.6 — `TestAGUIBudgetDerivation`:** Pin aguiRunBudget == writeTimeout-readDeadline-aguiWriteMargin, >0 and 18s; new d.aguiBudget matches, margin=2s, body bound=2s, tick=250ms, page=100 and cap=16. Do not shrink defaults before checking them.
- **AC3.7 — `TestAGUIShutdownEndsRuns`:** Real server run open: d.Shutdown returns nil in <2s; stream ends RUN_FINISHED within one shrunk tick plus scheduling slack. Include idle and blocked-read seam, observe context cancellation; no wait for frozen 10s drain expiry (PD7).
- **AC3.8 — `TestAGUIClientDisconnectReturns`:** Block read in request-rooted seam, cancel client; require read ctx canceled and handler returned within one tick plus bounded scheduling slack; join all goroutines/stop watchers, goroutine count returns to measured baseline under race; release fallback blocker only in cleanup so M-h reds promptly.
- **AC3.9 — `TestAGUIReadPostureMatchesLog, TestAGUIRejectsGET, TestAGUIInputErrors`:** Assert isProtected(POST /agui/)==isProtected(GET /v1/log)==false. Compare no auth/invalid/live/expired auth responses against log read posture; no filtering. GET =>405. Bad/trailing JSON, missing runId/threadId, non-string IDs, messages non-array/missing/null =>400; >64KiB =>413 PayloadTooLarge; optional messages contents/tools/context/forwardedProps accepted/ignored. Pin stream headers/no CORS and pre-stream APIError shape.
- **AC3.10 — `TestAGUIStoreErrorIsRunError`:** After start fail read seam with secret-bearing error: terminal RUN_ERROR code Internal constant sanitized message, no threadId/no store text; context deadline variant =>Timeout; pre-stream failure =>500 APIError; emit no second terminal frame. Capture detail only via existing internal log path.
- **AC3.11 — `TestAGUISlowBodyNeverTruncates`:** BodyBound=100ms, budget=300ms, tick=10ms; 1 byte/30ms body returning past bound =>408 SlowBody and no data; valid body completed at 80ms =>RUN_FINISHED by entry t0+300ms+tick+50ms. Reading client sees complete APIError or terminal stream. Add scaled scratch-server stalled-body arm =>zero bytes/clean EOF; no deadline setter. Total deliberate waits <2s.
- **AC3.12 — `TestAGUISeveranceResumable`:** Every byte offset of golden run: LF LF parser drops partial tail; recover snapshot/last complete delta cursor; append replay entry frames and prove no entry loss/gap (duplicate CUSTOM before its delta may be replayed). Also SSE-id oracle and header resume per PD2; never claim exactly-once CUSTOM delivery.
- **AC3.13 — `TestAGUISlowReaderCutIsResumable`:** Scratch server wraps d.Handler with WriteTimeout=200ms; about 5000 temporary-store entries whose frames exceed socket buffers. Read 1KiB then pause ~250ms, then drain delivered prefix with bounded client deadline; assert write-timeout cut and incomplete terminal run rather than normal finish. Ensure run budget exceeds scratch write timeout. Resume from last complete state and id per PD2/AC3.12, verify expected remaining entries. Real sockets mandatory; recorder alone cannot prove a cut.

**FIRST M3 action:** write and run `TestAGUIRESTGapCrossed` alone on C2 before any M3 production edit. Require the three commit statuses and frozen GET-range assertions to pass, then require the `/agui/` assertion red. After code, run this exact test first again; do not build the tail on the dense range loop.

**Implementation after red proof:** Anchor t0 at the first handler statement. Parse at most64KiB with clock checked on every Read return, body bound2s; typed408 if overdue Read returns by entry+25s, later =>write zero bytes. No new read after entry+18s; each read bounded by existing readDeadline and rooted in request/stop cancellation. Validate cursor and cap before stream; snapshot is cursor entry, not current head. Ignore auth exactly like log GET, ignore optional non-state fields. Keyset page100, immediately repoll full pages, wait cancellably only on short page; flush once per entry pair. Open RUN_STARTED+STATE_SNAPSHOT, terminal Finished(lastIndex) or sanitized Error; release semaphore on every exit. RegisterOnShutdown closes stop idempotently. No commit signal dependence or deadlines relaxation.

**Acceptance commands, in order** (use environment above; every named test must show PASS, no SKIP/zero-match):

```sh
go test -count=1 -v ./host/daemon -run '^TestAGUIRESTGapCrossed$'
go test -count=1 -v ./host/daemon -run '^TestAGUI'
go test -count=1 ./host/agui ./host/store
go build ./...
go test -count=1 ./...
go test -race -count=1 ./host/agui ./host/daemon ./host/store
go vet ./...
```

**Mutation drills on this milestone snapshot:**

| Mutant | Isolated edit | Test that must red / assertion |
|---|---|---|
| M-a | cursor +2 / skip final row of full page | `TestAGUIGoldenRun`, `TestAGUIResumeExact`: expected indexes and >100-entry page-boundary arm |
| M-c | tail via GetLogEntry(cursor+1) | `TestAGUIResumeExact`: k=2 must deliver entry5 |
| M-e | budget=writeTimeout | `TestAGUIBudgetDerivation`: derived equality/default budget |
| M-f | remove semaphore check | `TestAGUIGlobalCap`: 17th is503 StreamLimit |
| M-g | remove RegisterOnShutdown hook | `TestAGUIShutdownEndsRuns`: shutdown cancels blocked read and ends RUN_FINISHED promptly |
| M-h | read ctx from context.Background | `TestAGUIClientDisconnectReturns`: actual read ctx cancels and handler returns |
| M-i | CUSTOM adds transition payload field | `TestAGUIEntryValueEqualsLogRoute`: raw value+LF equals single-entry route bytes |
| M-m | add /agui/ to isProtected | `TestAGUIReadPostureMatchesLog`: read-posture equality false |
| M-n | header wins conflicting cursor | `TestAGUIResumeExact`: conflict =>400 before SSE |
| M-o | tail only on in-handler commit signal | `TestAGUISeesDirectStoreCommit`: direct store commit appears within2ticks |
| M-p | copy store error text into message | `TestAGUIStoreErrorIsRunError`: constant message excludes secret sentinel |
| M-q | anchor run clock after body | `TestAGUISlowBodyNeverTruncates`: 80ms body still finishes within entry-anchored bound |
| M-r | ignore body bound | `TestAGUISlowBodyNeverTruncates`: slow body =>408 no data |
| M-s | ignore state cursor | `TestAGUIStockClientStateResume`: state-only replay gives only newer entries |
| M-t | reject schema-less state | `TestAGUIResumeExact`: state:{} =>genesis200 |
| M-v | emit id before final entry frame | `TestAGUISeveranceResumable`, `TestAGUISlowReaderCutIsResumable`: id-bearing frame matches completed delta at every cut (PD2) |
| M-w | rebuild tail using dense /v1/log loop | `TestAGUIRESTGapCrossed`: AGUI crosses1to5 while frozen log range stops |

**Controller commit message:** `feat(daemon): stream committed log events over POST /agui/`

## C4 / M4 — S7 live-watch usage surface

**Exact file paths:**

- `docs/QUICKSTART.md`
- `host/daemon/daemon.go`
- `host/daemon/agui_test.go`

**Tests FIRST; every new test must be red on the pre-milestone snapshot:**

- **AC4.1 — `TestAGUIQuickstartControl`:** Documentation instrument-health only (PD5): QUICKSTART Watch the world live, complete payload curl -N POST, first three frames RUN_STARTED/STATE_SNAPSHOT/CUSTOM, Last-Event-ID resume and 18s re-POST. Execute both printed curl examples verbatim against temporary daemon fixture; record transcript for controller/PR, not a new repo artifact. Route doc comment updated.

**Implementation after red proof:** Add Watch the world live and route comment, exact curl POST payload, first three frames, header resume and18s re-POST contract. Execute printed examples verbatim on a temporary daemon/store, with payload construction included; controller retains actual bytes/status/transcript in PR. Do not touch mission bookkeeping.

**Acceptance commands, in order** (use environment above; every named test must show PASS, no SKIP/zero-match):

```sh
go test -count=1 -v ./host/daemon -run '^TestAGUIQuickstartControl$'
rg -n 'Watch the world live|Last-Event-ID|18 s' docs/QUICKSTART.md
curl -N -X POST -H 'Content-Type: application/json' -d '{"threadId":"t","runId":"r","messages":[]}' http://127.0.0.1:7644/agui/
curl -N -X POST -H 'Content-Type: application/json' -H 'Last-Event-ID: 1' -d '{"threadId":"t","runId":"r-resume","messages":[]}' http://127.0.0.1:7644/agui/
go test -race -count=1 ./host/agui ./host/daemon ./host/store
go vet ./...
```

For these exact curl rehearsal commands, the controller provides a disposable fixture daemon on `127.0.0.1:7644` (entries 0,1,2,5,6; store wholly inside the rehearsal temporary directory), never the live/user store. These manual S7 commands exercise the default 18s run contract; they are not sleeps inside daemon tests. AC4.1 grep is **instrument-health only**; it does not discharge stream behavior. `TestAGUIQuickstartControl` is likewise a documentation control. The executed curl rehearsal is separate S7 evidence.

**Mutation drills on this milestone snapshot:**

| Mutant | Isolated edit | Test that must red / assertion |
|---|---|---|
| None assigned by design | Documentation health/control and executed S7 rehearsal | Transport mutants remain bound to C1–C3 |

**Controller commit message:** `docs(agui): document and rehearse live world watching`

## Gates, mutation protocol and handoff

Milestone gates are the exact commands above: narrowest affected packages plus `go vet ./...` at every C1–C4 boundary. M2 retains the design’s daemon/store gate; M3 retains its build, full plain tests and focused race gate. Final C4 handoff runs:

```sh
export AILANG_BIN=$HOME/.pinned-ailang/ailang
export PATH=/opt/homebrew/bin:$PATH
./scripts/verify_ail.sh
./scripts/verify_go.sh
go test -race -count=1 ./host/agui ./host/daemon ./host/store
```

Full verify_go includes build, exact evidence manifest, plain tests and race tests (8m package timeout, 600s race-leg process budget). Controller runs CI and performs commits; executor does not rerun CI. Do not run full verification concurrently with mutation timing drills. Capture each rc immediately, without piping away the actual command status; record named PASS/FAIL/SKIP and durations. Require every planned test name to be discovered exactly once (subtests additional); named manifests are instrument-health controls, not replacements for assertions.

For every M-a through M-w, apply the exact isolated edit on its own Cn snapshot, run the root JSON’s named mutation command, require the specified assertion to turn red, restore the original bytes from snapshot without git, re-run green and compare bytes before the next edit. No compile error, socket denial, timeout of the test harness or unrelated failure counts as a kill. Record mutant, test, fired assertion and restore evidence. A SURVIVED mutant is an acceptance finding; report it and block handoff, never silently alter the mutant or weaken the test. M-a has two independent forms; drill both cursor+2 and full-page final-row loss. M-v drills early id framing with both pure truncation and real slow-reader tests. M-b runs against constant as well as clock timestamp additions.

## Risks and controller work

- CI Go race leg has ~600s wall budget (8m per-package timeout). Keep new tests fast: shrunk per-daemon budgets/ticks, concurrent held streams, channel handshakes, no test sleeping >~2s total. Report focused AGUI and full daemon race wall; do not widen CI budgets.
- Slow-reader test requires actual sockets and fixture bytes larger than buffers; bounded post-pause drain distinguishes write-timeout severance from normal finish. Never skip this oracle or substitute recorder evidence.
- Timing oracle for M-q needs handler-entry timestamp and 80ms body stimulus; 50ms slack must remain below the 80ms shift. Report race/load flakes; do not silently increase slack to let mutant survive.
- Cancel and join every stream/read/stop watcher before test return; shutdown channel closure must be idempotent; semaphore cleanup on every exit.
- Schema asset must be exact V16 bytes; custom names are World convention, not a schema namespace requirement. No new dependency or package manifest change.
- Golden fixture contains gaps; replay can repeat CUSTOM from an interrupted entry pair. State replace operations are idempotent; promise lossless resumability, not exactly-once event delivery.
- Frozen /v1/log retains first-gap stop, and commit density/linkage remain unfixed. Controller records design section9 findings separately; executor changes no mission doc.

Controller owns row148 landing state and design §9 follow-ups: density/linkage/gap finding and row148 reference from w-world-live-surface. These require no executor edits to `design_docs/world-mission*.md` or the design itself. Merge checks re-evaluate conflict surface against concurrent rows.

**Coverage census:** 4 milestones, 25 AC labels (AC1.1–1.5; AC2.1–2.4; AC3.1–3.14 plus AC3.2b; AC4.1), 23 mutants M-a…M-w, 7 numbered planner decisions. No AC or mutant deferred.
