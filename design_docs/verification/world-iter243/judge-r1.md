# Judge r1 — row 148 `w-worldd-agui-event-stream` (iteration 243)

Worktree `.wt-world-iter243-eval` at `d2d1815` (detached, nothing committed or pushed). Pinned binary `~/.pinned-ailang/ailang`.

## Verdict: PASS — 90/100, zero BLOCKING findings

| Rubric | Max | Score |
|---|---|---|
| Correctness | 30 | 28 |
| AC coverage / test quality | 25 | 22 |
| Mutation resistance | 15 | 12 |
| Design / standards conformance | 15 | 14 |
| Docs / usability | 10 | 9 |
| Conflict-surface hygiene | 5 | 5 |
| Total | 100 | 90 |

## BLOCKING findings
None.

## NON-BLOCKING findings
1. **Surviving mutants (test gaps, not behaviour bugs)**, see table:
   - O19 `host/daemon/agui.go` `if len(page) == aguiPage { continue }` removed: the "read again immediately on a full page" rule (D5) is untested. A backlog of >100 entries would drain at 100 per 250 ms tick instead of at once. Fix: a test with ~250 entries and a short budget asserting all are delivered in well under 3 ticks.
   - O17 the `remaining < wait` clamp on the idle sleep removed: no test pins that the final sleep is clamped to the budget (overshoot at most 250 ms; still inside the 2 s margin). Fix: a test with budget not a multiple of tick, asserting `RUN_FINISHED` by budget + small slack.
   - O2 `<-watchDone` join dropped (`defer cancel()` only): the watcher still exits on cancel, so no real leak, but AC3.8's goroutine-baseline check does not detect an unjoined watcher. Low value.
   - O18 the `time.Since(t0) <= 25s` guard on the 408 write is not testable cheaply (needs a 25 s wait); the design records it as measured (V24). Acceptable.
2. `host/daemon/agui.go` pre-stream `GetLogEntry` failure caused by shutdown cancelling `runCtx` is reported through `writeInternalError` (HTTP 500) rather than a clean refusal. Benign (shutdown race, before any byte), but the error log gets a spurious "internal error" line.
3. `Last-Event-ID` uses `strconv.ParseInt`, which accepts a leading `+` (e.g. `+3`). Harmless; the state cursor path is stricter.
4. **Race-mode wall time.** `TestAGUISlowReaderCutIsResumable` takes 17.9 s under `-race` (about 5000-entry fixture seeding plus a 200 ms WriteTimeout cut). It is the only AGUI test over 2 s; none uses `time.Sleep` over 250 ms (sleeps at agui_test.go:731, 746, 944 are 30 to 250 ms). Whole-package `-race` for `host/daemon` is 238 s, `host/store` 25 s, total about 242 s, which fits CI's ~600 s go leg but leaves less headroom; consider a smaller fixture if the leg tightens.
5. QUICKSTART sample output shows fixture hashes (`sha256:000…064`, gapped 0,1,2,5,6) from the executor's seeded DB, not a real commit. The doc labels the CUSTOM value an "example" and the snippet says "your hashes and writer will differ", so it is honest, but the executed transcript's 5-entry store does not match the doc's "entries 0 and 1 committed (§2-§3)" precondition. A real-commit run needs a minted session (tty), which I could not do.

## 1. AC coverage (checked by reading the test, not the name)
All 25 ACs have a real test.
- AC1.1 `TestGolden`/`TestDeterministic`: golden bytes over entries 0,1,2,5,6 with a gap; M-l/M-v-agui variants go red. Real.
- AC1.2 `TestEventsConformToPinnedSchema`: pinned schema sha, null controls (events > 0, 6 defs found). Real.
- AC1.3 `TestFramesSplitLikeStockClient`, AC1.5 `TestDeltaIdempotent`, AC1.4 `TestAGUIDependencyBoundary`: real (go list deps).
- AC2.1-2.4 `TestLogEntriesAfter*` and `TestAGUIReadStoreSeam`: gap crossing, exclusive cursor (M-d killed), bounds, deadline, quarantine (O20 killed), deep-equal to `GetLogEntry`.
- AC3.1 `TestAGUIGoldenRun`, AC3.2 `TestAGUIResumeExact` (all k, state vs header, malformed, conflict 400, unknown 404, `state:{}` genesis), AC3.2b `TestAGUIStockClientStateResume` (test-local RFC 6902 applier fed run-1 bytes, then more commits): non-vacuous, kills M-s.
- AC3.11 `TestAGUISlowBodyNeverTruncates`: both arms with timing assertions (kills M-q, M-r, O7).
- AC3.12 `TestAGUISeveranceResumable` (every byte offset), AC3.13 `TestAGUISlowReaderCutIsResumable`: real, kill M-v.
- AC3.14 `TestAGUIRESTGapCrossed`, AC3.3 byte equality to `GET /v1/log/{i}`, AC3.4 direct store commit, AC3.5 cap + release + Retry-After, AC3.6 derivation, AC3.7 shutdown (11 s failure under M-g shows it genuinely depends on the hook), AC3.8 disconnect, AC3.9 posture/405/inputs, AC3.10 store error: real.
- AC4.1 is a QUICKSTART grep (instrument-health only, as designed).
No tautological test found. The only untested behaviours are those in non-blocking finding 1.

## 2. Correctness review
- Semaphore: acquired after body/cursor validation, released by a single `defer` on every exit (O1 confirms test). Over-cap path writes 503 before any byte.
- Watcher goroutine: `defer func(){cancel(); <-watchDone}()` joins on every path past its creation. No lock is held across a write; there are no shared mutable fields besides the channel/Once.
- Shutdown: `RegisterOnShutdown` closes `aguiStop` once; run ends with `RUN_FINISHED` because the `r.Context().Err()==nil` gate is on the request context, not the run context. Correct.
- Budget/clock (D1): `t0` is the first statement, `end = t0 + aguiBudget` (18 s derived), loop re-checks the clock before each read, idle sleep clamped to remaining (see O17), each read capped at `readDeadline`, so the last frame is by t0+28 s of the 30 s window. Body reader checks the clock on every Read return including EOF; 408 only when it returns by t0+25 s, otherwise nothing is written. Matches V24.
- Cursor (D4): schema-matching state is explicit; non-object/other-schema state means genesis; missing/null/float/string/negative lastIndex gives 400; header and explicit state must agree; unknown non-negative index gives 404 from `GetLogEntry` before any header or byte. Verified live against a real `worldd` (404 on `Last-Event-ID: 1` of an empty store, 400 on `abc`).
- Error split: pre-stream failures are `APIError` JSON; post-`RUN_STARTED` store errors become `RUN_ERROR` with code `Internal`/`Timeout` and the constant `internalErrorMessage`; detail goes only to `errLog`. `RUN_ERROR` has no threadId (schema-checked).
- SSE framing: `\n` only; `id:` only on the entry's last frame; one flush per entry pair; write failure ends the run.
- Exposure: `CUSTOM.value` is `json.Marshal(logJSON(e))`, byte-equal to `GET /v1/log/{i}` (test + M-i); `isProtected` equality test (M-m killed); no CORS headers; read-only.
- `host/agui` is pure (no net/http, no store; dependency test). `LogEntriesAfter` requires a deadline, honours quarantine, bounds limit, uses the INTEGER PK.

## 3. Conflict surface
`git diff --stat d6334c2..HEAD -- host/daemon/handlers.go host/projection tools/launchd design_docs/world-mission.md` is EMPTY. `daemon.go` diff is exactly: `sync` import, six fields, one `readStore` method, constructor init plus `RegisterOnShutdown`, route comment, one `HandleFunc`. `isProtected` and D7 constants untouched. `host/store/deadline_guard_test.go` has a one-line census entry (expected).

## 4. Mutation drills (mine, re-run at HEAD with `-race`, backup/restore verified by shasum, tree clean afterwards)
| Mutant | Edit | Test(s) | Result |
|---|---|---|---|
| M-c | tail via `GetLogEntry(after+1)` | GoldenRun/ResumeExact | KILLED (golden bytes differ) |
| M-g | remove `RegisterOnShutdown` | ShutdownEndsRuns | KILLED |
| M-q | `end` anchored at `time.Now()` | SlowBodyNeverTruncates | KILLED |
| M-s | state schema never recognised | StockClientStateResume | KILLED (got all, want [7 8]) |
| M-v | `id:` before CUSTOM frame | SeveranceResumable, TestGolden (agui) | KILLED (both) |
| M-w | SQL `= ?+1` (dense loop) | AGUIRESTGapCrossed | KILLED |
| M-d | SQL `>=` | LogEntriesAfterCrossesGap | KILLED |
| M-f | semaphore check bypassed | GlobalCap | KILLED |
| M-k | `\r\n\r\n` framing | TestGolden | KILLED |
| O1 | drop slot release | GlobalCap | KILLED |
| O3 | Last-Event-ID lower bound -5 | ResumeExact | KILLED |
| O4 | state lastIndex lower bound -5 | ResumeExact | KILLED |
| O5 | unknown cursor 404 removed | ResumeExact | KILLED |
| O6 | Retry-After removed | GlobalCap | KILLED |
| O7 | budget x3 | SlowBodyNeverTruncates | KILLED |
| O8 | Timeout code collapsed to Internal | StoreErrorIsRunError | KILLED |
| O9 | body cap 64 KiB to 1 MiB | InputErrors | KILLED |
| O10/O11 | nosniff / no-store removed | InputErrors | KILLED (both) |
| O12 | `Finished` lastIndex off by one | GoldenRun | KILLED |
| O13 | no RUN_FINISHED at end | GoldenRun | KILLED |
| O15 | messages-array check removed | InputErrors | KILLED |
| O16 | threadId nil check removed | InputErrors | KILLED (nil-deref panic) |
| O20 | quarantine check removed | LogEntriesAfterBounds | KILLED |
| O14 | snapshot entry dropped | n/a | INCONCLUSIVE (compile error, unused var); not counted |
| O2 | watcher join dropped | all TestAGUI | SURVIVED |
| O17 | idle-sleep clamp removed | all TestAGUI | SURVIVED |
| O18 | 25 s guard on 408 removed | all TestAGUI | SURVIVED (untestable cheaply) |
| O19 | no immediate re-read on full page | all TestAGUI | SURVIVED |
(An "M-l" row I attempted was a no-op edit and is discarded.) Executor's own ledger (`executor-mutants.jsonl`) lists M-a, M-c, M-e..M-i, M-m..M-t, M-v, M-w as killed; my re-runs of M-c, M-g, M-q, M-s, M-v, M-w, M-f agree. 9 design mutants and 22 of 26 of mine killed: 4 genuine survivors (finding 1).

## 5. Gates
| Command | rc | Wall |
|---|---|---|
| `go vet ./...` | 0 | 1 s |
| `go test -count=1 ./host/agui ./host/store ./host/daemon` | 0 | 147 s (agui 2.1, store 10.7, daemon 146.6) |
| `go test -race -count=1 ./host/agui ./host/store ./host/daemon` | 0 | 242 s (agui 3.2, store 24.9, daemon 238.6) |
New AGUI tests under `-race` (`-v` run): total 26 s; slowest TestAGUISlowReaderCutIsResumable 17.9 s, SeveranceResumable 2.8 s, DependencyBoundary 1.7 s, GoldenRun 1.1 s, everything else under 1 s.

## 6. Docs ("Watch the world live", QUICKSTART section 10)
Built `go build -o /tmp/worldd-eval ./cmd/ailang-worldd`, served a fresh temp store on a free loopback port. The documented curl produced exactly the documented first frames (`RUN_STARTED`, `STATE_SNAPSHOT lastIndex -1 logHead null`) then `RUN_FINISHED {"lastIndex":-1}` after 18 s wall, as documented. `Last-Event-ID: 1` gave 404 `unknown cursor`, `abc` gave 400, GET gave 405, `{}` gave 400 `BadRequest`, all before any stream byte, matching the described behaviour. The CUSTOM/STATE_DELTA frames could not be reproduced on my store (committing needs a minted session, which needs a tty); they match the executor's transcript, which used a seeded fixture (finding 5). Docs match.
