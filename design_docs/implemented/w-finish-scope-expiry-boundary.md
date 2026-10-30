# w-finish-scope-expiry-boundary: `TestEpisodeFinishScopeExpiry` races two wall-clock reads across a second edge

- **Status**: Implemented iteration 327 (M1 `d7a54db5`, M2 `75d07431`, PR #436; judge sonnet PASS 93). Designed iteration 327, unattended. **r3**: r2 after quorum r1 BLOCKED, then the quorum-r2 fixes applied by the controller (§10). This is
  the one mandated revision (§10). It is a red on `dev` (`4d43d8e3`, race gate, attempt 1), so it
  outranks the queue. The controller numbers the row (the charter's last row at base is 235, V12).
- **Item**: a test defect from row 218 M2 (`8c5a14e6`, 2026-10-29, w218b, Codex, attended track,
  V2). The production bounds are correct. The fix adds one unexported, route-scoped clock field
  (`finishNow`) to the daemon so the tests can pin the exact edges.
- **Estimated**: ~0.1 day. About +14/−5 LOC production (`daemon.go`, `episode_finish.go`), about
  +70/−8 test (`episode_finish_test.go`). Prototype diff: V9.
- **Measurement base**: `4d43d8e3` (worktree `sprint/world-iter327-finish-scope-expiry`, V1).
- **Instruments**: `go1.26.6 darwin/arm64`, 16 cores; `AILANG_BIN=$HOME/.pinned-ailang/ailang`
  (v0.41.0, the CI pin). Mutants, harness runs and the prototype ran in throwaway copies
  (`/tmp/w327d/base`, `/tmp/w327d/fix`, both `git archive HEAD`), never in the worktree. The
  scripts are in `verification/world-iter327/designer-tools.txt`.

---

## 1. The defect

`TestEpisodeFinishScopeExpiry` (`host/daemon/episode_finish_test.go:162`) builds both arms
before its loop:

```go
for _, expiry := range []int64{time.Now().Unix() - 1, time.Now().Unix() + 3601} {
```

The handler reads its own clock later, at `episode_finish.go:168`, and refuses when
`scope.ExpiresAt <= now || scope.ExpiresAt > now+3600` (`:169`). Call the test's read `c` and the
handler's read `s`, so `s >= c`. The `+3601` arm is refused only while `c+3601 > s+3600`, that is
only while `s == c`. **If any whole-second edge falls between the two reads, the arm is accepted
and the route returns 200.** The gap is not one request long. It covers the whole first iteration
(proof-file write, request, handler for the `-1` arm), plus the second proof write and the second
request's path up to `:168`.

- **CI (V3).** Run 41093918891 attempt 1, "go race gate": `--- FAIL: TestEpisodeFinishScopeExpiry
  (0.15s)  episode_finish_test.go:170: expiry 1793383038: 200`. Attempt 2 is green. The parents
  were green, and the merge was docs-only.
- **Forced (V5).** If a spin to `Nanosecond() >= 999_950_000` is placed before `:164`, the test
  fails 20/20 plain and 20/20 under `-race`. The same spin to just after an edge (`<= 100 ms`)
  passes 20/20. Idle and unmodified, it passed 400/400 under `-race`. So the red is a function of
  where the second edge falls, not of load as such. Load only widens the gap.
- **The `-1` arm is safe by direction.** A later `s` only makes `c-1 <= s` more true.

The same shape exists once more, in the same file (§7). It is `TestEpisodeFinishOperatorProofAC11`
`future` (`:79`, a proof at `c+301` against `:117`'s `p.ExpiresAt > now+300`). It fails 20/20
under the same harness (V5).

## 2. The decision: a daemon clock seam, with exact-edge arms on a pinned clock

**Production (minimal).** Add an unexported `finishNow func() int64` field to `Daemon`, in Unix
seconds. `New` sets it to `func() int64 { return time.Now().Unix() }`. The three reads in
`episode_finish.go` (`:116`, `:168`, `:270`) become `now := d.finishNow()`, and the file drops its
`time` import. Nothing else in the daemon changes: the other 17 `time.Now()` reads in
`host/daemon` stay as they are (V11). The fields `drainTimeout` and `readDeadline` already exist so
that tests can reach a shipped branch (`daemon.go:441-454`), and `cmd/ailang-worldd`'s
`sessionEnv.now` is a component-local clock seam of the same kind (`session_test.go:29-33`).

**[r2] Reuse or a new seam: the coordinator's clock (oc-kimi-k3).** The daemon already wires one
injected clock, at `daemon.go:835-836`:
`coordinator.New(coordinator.Config{…, Now: func() int64 { return time.Now().Unix() }, …})`. It is
also `int64` Unix seconds. It is read at `coordinator.go:354` (`planInvocation`) and
`effectful.go:284` ("one logical time: the effect requests and the intent"), both on the invocation
path. I measured three options (V16):
- **(i) Reuse the coordinator's clock: impossible as the code stands.**
  - The clock lives in the coordinator's unexported `cfg Config` (`coordinator.go:59`), and
    `*Coordinator` has no clock accessor (0 matches; the control finds 2 exported methods).
  - `d.coord` exists only when `cfg.AilangBin != ""` (`daemon.go:818`). The finish route has to
    work without it: `newHandlerDaemon` (`handlers_test.go:23`) sets no `AilangBin`, and its finish
    tests pass.
  - So reuse would need a new coordinator API and would still leave the route with no clock when no
    coordinator is configured.
- **(ii) One daemon clock, which the daemon also hands to the coordinator**
  (`Now: func() int64 { return d.finishNow() }`). It is possible, but rejected for this row.
  - Production would be byte-for-byte the same, because both clocks are the wall clock.
  - The tests would change. Pinning the shared clock moves grant liveness and intent stamps for
    every invocation on that daemon. The finish tests pin to 2033 precisely so that wall-clock
    reads show up, and a shared pin would make every wall-clock-minted grant read as expired.
  - The coordinator reads its clock inside `Dispatch`. The MCP handler (`projection/mcp.go:130`)
    and the A2A handler (`projection/projection.go:363`) call `Dispatch`, and they run on the
    server's per-request goroutines, which tests drive concurrently (V16). A shared clock that a
    test can write would then need an atomic, or a rule that it is pinned before serving starts.
  - The wiring would need its own killer: a coordinator-equipped daemon (`AilangBin`) pinned and
    then checked against an intent stamp. That is a heavier fixture than this flake row warrants.
- **(iii) Two seams, chosen.** The stated reason is that **no comparison crosses them**.
  - Every read of an approval's request or decision stamp, repo-wide, is in
    `validatePublishApproval` (`broker/approve.go:517`, lines 629-644). That function is reached only
    for publish approvals.
  - The finish approval's stamps (written by `RecordOperatorApproval` with the route's `now`) are
    read back only by `finishDecisionScope`. It reads `scope.ExpiresAt`, never the stamps.
  - Decisions are looked up by request ref (`findApprovalDecision`, `approve.go:252`). Only the
    finish route holds a finish approval's request ref.
  - So two clocks pinned differently cannot produce a wrong verdict. In production they are the
    same function.
  - The field is named `finishNow`, not `now`, and its doc comment says it is not the coordinator's
    clock, so that a future reader does not take it for a daemon-wide clock.
  - The daemon's own coordinator clock is not pinnable after `New` at all. Tests that need a pinned
    coordinator clock build their own `coordinator.New` (7 sites in `host/daemon` tests).
  - So the daemon carries exactly one clock that a test can pin.

  If a later row needs a finish-route time to be compared with a coordinator time, that row must
  first move to (ii).

**[r2] Concurrency (gemini-4-argon).** All three reads of `d.finishNow` run on the request path
while `d.humanMu` is held: the lock is taken at `:158`, before `consumeOperatorProof` (`:160`) and
before both phases. No background goroutine reads the field. The tests write it only between
synchronous `ServeHTTP` calls on an `httptest.ResponseRecorder`, so no handler is in flight. The
field's doc comment states this rule: write it only before or between synchronous handler calls.
The prototype is green under `-race` (V10, V17).

**S2:** the clock is a host effect and stays in host Go. The seam only moves where it is read.

**Tests.** The clock is pinned at `finishClockAt = 2_000_000_000`, which is about 6.6 years from
the wall clock. So a handler that ignores the seam refuses every accept arm (mutants S1–S3).
`pinFinishClock(d, at)` sets `d.finishNow = func() int64 { return at }`.
- **T1** rewrites `TestEpisodeFinishScopeExpiry` (`:169`). Scope expiry at `now` plus −1, 0, +1,
  +3600 and +3601 must give 403, 403, 200, 200 and 403.
- **T2 `TestEpisodeFinishOperatorProofExpiryBoundary`** is new (`:117`). Proof expiry at `now`
  plus −1, 0, +1, +300 and +301 must give 403, 403, 200, 200 and 403.
- **T3 `TestEpisodeFinishRecordRechecksScopeExpiry`** is new (`:271`). **[r2, measured by coverage,
  V15]:** today `:271`'s condition is evaluated, by AC11's record path, but its refusal branch
  (`return bad()`, `:271-273`) is executed by no test. The expiry-bound mutants of `:271` survive
  every test (V6). It approves a scope at `T+3600` with the clock at `T`. Then it records a
  `push-failed` outcome (not single-use) with the clock at T−1, T, T+3599 and T+3600, which must
  give 403, 200, 200 and 403.
- **AC11 `future`** changes `+301` to `+3600`. The exact edge now belongs to T2, so AC11 keeps a
  wall-clock arm with a margin in the safe direction.

Alternatives, with their mutation cost measured:
- **(a) A wider margin (`+3601` to `+3660`), test only.** It ends the flake: 20/20 green under
  the harness (V7). But the boundary stops killing anything inside the margin. P4 (`3601`) and
  P4b (`3659`) **survive** (V7). P1 and P3 stay unkilled as they are today, and `:271` stays
  untested. Rejected, because it trades a flake for a blind band.
- **(a′) Re-read the clock right before each request.** This narrows the gap but does not close it.
  The harness spin still lands before the read. Rejected.
- **(c) Derive the edge from a clock read after the request, or retry or skip when the second
  changed.** The first cannot give an exact value. A retry hides flakes, and a skip is an S6
  vacuity. Rejected.

## 3. Risk

`New` always sets `d.finishNow`. **[r3, controller] Every construction of a `host/daemon.Daemon`, by a type-aware audit (V18):** `go/packages` + `go/types` over `./...` with tests (108 packages, 0 load errors) enumerates composite literals whose resolved type is `daemon.Daemon` under any import alias, `new(daemon.Daemon)`, and zero-value declarations. It finds **11**: `daemon.go:784` inside `New`, and the 10 test literals below; **0** `new()` and **0** zero-value declarations. A scratch control file (`package daemon_test`, alias `dd`) adds exactly 3 (one per branch), so each branch can fire. None of the 6 files holding the 10 test literals mentions the finish route (0 `finish` matches each; control `episode_finish_test.go` 10). V13/V14 below are textual-search results and are kept as such; V18 is the construction audit. **[r2] Textual search (V14):**
- **Non-test:** exactly one, `daemon.go:784`, inside `New` itself. Outside the package there are 0
  `daemon.Daemon{` or `new(daemon.Daemon)` sites; the control, any `daemon.<Type>{` literal, finds
  85. `Daemon`'s fields are unexported, so no other package can build a usable one.
- **Test:** ten sites, all in `host/daemon` (V13).

A literal-built daemon has a nil `finishNow` and would panic, loudly, if it hit the finish route.
Coverage shows the handler runs in-process only under `host/daemon`'s own tests (V15).
`cmd/ailang-worldd`'s finish tests talk to a stub route (`finishFixture.serve`,
`episode_finish_test.go:156`), and `cmd/world-publish` never executes `episode_finish.go`. In
`host/daemon` only `episode_finish_test.go` posts to the route, and it goes through `New`. I add
no nil guard, because a silent fallback to the wall clock would bring the flake back unseen. W1 proves that `New`'s wiring is load-bearing. The production
behaviour is byte-for-byte the same function of the wall clock. The whole `host/daemon` package
passes plain and under `-race` on the prototype (V10). **[r3] Package green verified on the pre-rename prototype; AC-3 at execution is the first package run of the r2 diff.**

## 4. Mutation table

The killer is the assertion that fires. Base results are V4/V6; prototype results are V8 (the r2
rename to `finishNow` was re-run for S1–S3 and W1: V17). **[r2]** Base lines are `:117/:169/:271`. **[r3] P/Q/R/TM1 kill attributions verified on the pre-rename prototype; the r2 diff's first verification is AC-2 at execution time.**
The S rows give **prototype** lines `:115/:167/:269`, one lower because the `time` import is
removed. The seam decision (§2, option iii) adds no mutant: the coordinator's wiring is unchanged.

| ID | Mutant | Base `4d43d8e3` (current tests) | Prototype killer |
|---|---|---|---|
| P1 | `:169` `<= now` → `< now` | **survives** | T1 arm `now+0` (`:187`, `200, want 403`) |
| P2 | `:169` lower arm dropped | killed (`:170`, the `-1` arm) | T1 arm `now-1` |
| P3 | `:169` `> now+3600` → `>=` | **survives** | T1 arm `now+3600` (`403, want 200`); also T3 |
| P4 | `:169` `3600` → `3601` | killed (the `+3601` arm) | T1 arm `now+3601` |
| P5 | `:169` `3600` → `3660` | killed | T1 arm `now+3601` |
| P6 | `:169` upper arm dropped | killed | T1 arm `now+3601` |
| Q1, Q2 | `:271` `<= now` → `<`; lower arm dropped | **survive** (both finish test files, V6; the refusal branch is never executed, V15) | T3 arm `approve+3600` (`:231`, `200, want 403`) |
| Q3 | `:271` `>` → `>=` | **survives** | T3 arm `approve+0` (`403, want 200`) |
| Q4, Q6 | `:271` `3600` → `3601`; upper arm dropped | **survive** | T3 arm `approve-1` (`200, want 403`) |
| R1 | `:117` `<= now` → `<` | **survives** | T2 arm `now+0` (`:204`) |
| R2 | `:117` lower arm dropped | killed (AC11 `expired`) | AC11 `expired`; T2 |
| R3 | `:117` `> now+300` → `>=` | **survives** | T2 arm `now+300` (`403, want 200`) |
| R4 | `:117` `300` → `301` | killed (AC11 `future`, flaky arm) | T2 arm `now+301` (AC11's `+3600` no longer kills it) |
| R6 | `:117` upper arm dropped | killed (AC11 `future`) | AC11 `future`; T2 |
| S1/S2/S3 | `d.finishNow()` → wall clock at prototype `:115`/`:167`/`:269` | n/a | S1, S2: T1, T2, T3. S3: T3 arm `approve+0` |
| W1 | `New` wires `finishNow` to `func() int64 { return 0 }` | n/a | AC11 final approve (`:106`, `403`) |
| TM1 | T1 without `pinFinishClock` | n/a | T1 arm `now+1` (`403, want 200`): the pin is load-bearing |
| TM2 | option (a): `+3601` → `+3660` (base test) | P4 and P4b (`3659`) **survive**; P5 and P6 are killed | (rejected design) |

Base results: at `:169` the current test kills 4 of 6 (P2, P4, P5, P6). At `:271` it kills 0 of 5.
At `:117`, AC11 kills 3 of 5 (R2, R4, R6). The prototype kills all 16 bound mutants, plus S1–S3,
W1 and TM1 (V8).

## 5. Acceptance criteria

- **AC-1 (red at base, green after; deterministic).** The harness inserts
  `for time.Now().Nanosecond() < 999_950_000 {}` on its own line before the arm's first clock
  read. At base that is before `:164`. On the fix it is before `pinFinishClock(d, finishClockAt)`
  in T1. Then run
  `go test -timeout 300s ./host/daemon/ -run '^TestEpisodeFinishScopeExpiry$' -count=20 -v`,
  plain and with `-race`, and restore the file. (`designer-tools.txt` has `harness.sh`, which does
  this and checks the file's sha256 before and after.)
  - Want at base: 20 `--- FAIL`, first line `expiry <n>: 200`. Measured: 20/20 plain and 20/20
    race (V5).
  - Want on the fix: 20 `--- PASS`, 0 FAIL. Measured: 20/20 plain and 20/20 race (V8).
  - Control: at base, the spin `for time.Now().Nanosecond() > 100_000_000 {}` gives 20/20 PASS
    (V5). This proves the edge, not the edit, causes the red.
- **AC-1b.** The same harness before AC11's `future` line (`:79`), with
  `-run '^TestEpisodeFinishOperatorProofAC11$'`. Base: 20/20 FAIL (V5). Fix: 20/20 PASS (V8).
  Control at base: 20/20 PASS (V5).
- **AC-2.** Apply every mutant in §4 **except TM2**, one at a time, to the fixed tree. TM2 is the
  rejected design's measurement, on the base test. Each one reds exactly the
  killer named there, and the file is restored byte-identical. Run
  `go test ./host/daemon/ -run '^TestEpisodeFinish' -count=1 -v`. For S1–S3, a scratch file
  defines `wallNow()`, and the scratch file is removed afterwards. The prototype gave all killed
  (V8).
- **AC-3.** `go vet ./host/daemon/`. Then `go test ./host/daemon/ -count=1` and the same with
  `-race`. Then `go test ./host/daemon/ -run '^TestEpisodeFinish' -race -count=200`. Then
  `go test ./cmd/ailang-worldd/ -run Finish`. All are green (prototype: V10). **[r3] Package green verified on the pre-rename prototype; AC-3 at execution is the first package run of the r2 diff.**
- **AC-4.** The full gate: `go build ./... && go test ./...` and `./scripts/verify_ail.sh` (no
  `.ail` is touched). CI on the fix head must have all 3 jobs green, with
  `ok  github.com/sunholo-data/ailang-world/host/daemon` in the race gate. **[r2] I did not run the
  full `./...` or `verify_ail.sh` (V10).** The executor's AC-4 run is therefore the first
  full-suite check of this change, and CI on the fix head is the second. Before that, the evidence
  is per package: `host/daemon` whole, plain and race; `cmd/ailang-worldd -run Finish`; and
  `go build ./...`.
- **No stress AC.** Whether the test fails depends only on whether a second edge falls inside the
  read gap. AC-1 forces that, and after the fix no arm reads the wall clock. Any load bound would
  be a guess.

## 6. Files touched, conflict surface

| File | Change |
|---|---|
| `host/daemon/daemon.go` | the `finishNow func() int64` field and its comment (next to `drainTimeout`; it states the request-path, between-calls write rule and that this is not the coordinator's clock); `finishNow:` in `New`'s literal (`:790`). The coordinator wiring at `:835-836` is unchanged |
| `host/daemon/episode_finish.go` | `:116`, `:168` and `:270` use `d.finishNow()`; the `time` import is removed |
| `host/daemon/episode_finish_test.go` | `finishClockAt`, `pinFinishClock`, `finishScopeAt`; T1 replaces `:162-173`; T2 and T3; AC11 `future` `+301` becomes `+3600` |

Rows 218 and 219 are attended and Codex-owned. An open branch that edits `episode_finish.go` or
the three tests would conflict there. **S1/S4/S5:** n/a. **S3:** host only, not kernel. **S6:**
TM1 and S1–S3 show the pin is consumed, so the arms cannot pass by never reaching the edge.
**S7:** n/a (an unexported field, no usage surface).

## 7. Sibling sweep: client clock against the handler's later clock

**The hazard class:** a test arm built from its own wall-clock read at the **exact edge** of a bound
that the code under test checks with a **later** read of its own. Two things make a site SAFE:
- **Direction.** An "expired" arm at `c-k` stays expired, and an "accepted" arm at `c+B` (when the
  bound is `> now+B`) stays accepted.
- **Margin.** An arm at least 100 s from the edge would need a read gap that long.

An **injected** clock is safe by construction.

| Sweep (V11) | Hits | Control | Verdict |
|---|---|---|---|
| A: `time.Now().Unix…() ± <digits>` in `*_test.go` | 27 | `episode_finish_test.go:164` ✓ | **2 lines SAME-HAZARD:** `:164` `+3601` (this row; its `-1` arm is SAFE by direction) and `:79` `+301` (V5, fixed in this row). **6 lines SAFE by direction:** `:77` `-1`; `agui_test.go:458` (ttl 1 at `-10`); `session_middleware_test.go:114` (`-10`); `daemon_test.go:1123,1151` (ttl 60 at `-3600`); `projection_test.go:97` (`-3600`). **19 lines SAFE by margin:** `episode_finish_test.go:57` (+900); `:74,81,83,85,89,103,133,144,150,167` (+100 against a 300 s bound, so a 100 s gap would be needed); `cmd/…/episode_finish_test.go:93`, `session_git_test.go:56`, `credentials_split_test.go:51` and `projection_test.go:235` (+3600); `daemon_test.go:1059` and `projection_test.go:81,96` (+7200); `binding_test.go:185` (+600) |
| B: `<now-var> ± <digits>` in `*_test.go` | 39 | `home_service_test.go:66` ✓ (a hand-read site; the failing line is A-shaped, not B-shaped) | **0 SAME-HAZARD.** **14 injected:** authority `resolveNow`/`Mint` take `now` (`mint_test.go` ×4, `resolver_test.go` ×2, `resolvecontext_test.go` ×1); `quickstart_se_test.go` ×4, whose edges `now+3599/3600/3601` go to `NewCapabilitySnapshot(…, at)`; `session_new_test.go` ×3, through `testMintEnv(…, now)`. **7 direction:** `session_self_test.go:78`, `surface_page_test.go:121`, `surface_model_test.go:96`, `demo_test.go:60,67`, `home_service_test.go:64,66` (and `activity.begin` takes `now`). **18 margin:** +3600, +7200 and +1000 alive arms, including `home_service_test.go:68`; the `working-episode` `last=now` against the health snapshot's `last >= now-300` is also margin |
| C: production bounds (`expires`/`deadline` comparisons, `> now±k`) | 29 | `episode_finish.go:169` ✓ | Wall-clock reads that tests reach at an edge: only `episode_finish.go:117/169/271`, and this row fixes all three. The rest take `now` as a parameter (`broker/approve.go`, `decide.go`, `authority/resolver.go`, `binding.go`, `surface_page.go`, `session_new.go`), or no test probes their edge (`doctor_hosts.go:88`, `episode_finish_operator.go:121`, `activity` snapshot via `daemon.go:1118`, margin 300) |
| D: `time.Now().Add(` in `*_test.go` | 29 | n/a (shape-only) | SAFE: poll and conn deadlines, already-past contexts (`-time.Second`, `-time.Hour`), and ±1 h mtimes (`templates_test.go:283`, `worktree_delta_test.go:575`). None feeds an expiry edge |

**Result:** 2 same-hazard sites, both in `episode_finish_test.go` and both fixed here. No follow-up
row is needed. Row 235 (`TestHumanApproveIsSerializedAcrossEpisodes`) is a different mechanism (a
200 ms wait budget) and is out of scope.

## 8. Residuals

- **R1. The production client sits exactly on the proof's upper edge.**
  `cmd/ailang-worldd/episode_finish_operator.go:48` writes `time.Now().Unix() + 300`. That is
  accepted only because the CLI's read comes before the daemon's on the same host. A wall-clock
  step backwards between the two (an NTP slew) would refuse the proof. The failure is closed and
  retryable, so it is noted, not filed.
- **R2. The seam covers the finish route only.** After the fix, 17 non-test `time.Now()` reads
  remain in `host/daemon`: 20 at base, less the 3 this row moves (V11). The sweep found no other exact-edge test. A future exact-edge test elsewhere
  should get its own route-scoped seam, or move to §2 option (ii) if it must share time with the
  coordinator. It should not race the wall clock.
- **R3. Margin is not impossibility.** The SAFE-by-margin sites (§7) would need a read gap of 100 s
  or more. CI's whole `host/daemon` race run took 352 s, but one request takes milliseconds.
- **R4. The sweep is by shape.** A–D find literal offsets. An expiry computed inside a helper with
  no literal at the call site (for example `authority.Mint`'s ttl) was judged from the ttl arguments
  that A and B show, not from an exhaustive helper walk.

## 9. Verification Log

| # | Claim | Command | Observed |
|---|---|---|---|
| V1 | Base and tools | `git log --oneline -1`; `git status --porcelain`; `go version`; `~/.pinned-ailang/ailang --version` | `4d43d8e3 Merge pull request #432 …`; only `?? design_docs/verification/world-iter327/`; `go1.26.6 darwin/arm64`; `AILANG v0.41.0` |
| V2 | Origin | `git blame -L162,174 host/daemon/episode_finish_test.go`; `-L166,172 episode_finish.go`; `git log -1 8c5a14e6` | every line `8c5a14e6c … 2026-10-29 14:34:07 +0100`; `feat(row 218 M2): finish's approval, operator route, push and PR dispatch, outcome records, revoke (w218b, Codex)` |
| V3 | The CI red | `gh api …/actions/runs/41093918891/attempts/1 --jq .conclusion`; `…/attempts/1/jobs`; `gh api …/runs/41093918891` | attempt 1 `failure`, `go race gate` failure, other 2 success; run now `attempt 2`, `success`. FAIL text from the controller's `gate1-2-readings.txt` |
| V4 | Base mutants (P, R) | `runmut.sh /tmp/w327d/base '^TestEpisodeFinish' 3 <id> host/daemon/episode_finish.go <line> <old> <new>` (`designer-base-mutants.txt`) | P1 rc=0, P3 rc=0, R1 rc=0, R3 rc=0 (survive); P2, P4, P5, P6 3/3 `TestEpisodeFinishScopeExpiry` (`:170 … 200`); R2, R4, R6 3/3 `TestEpisodeFinishOperatorProofAC11`. sha256 `ea44759896b5` before and after each |
| V5 | Harness at base, and controls | `harness.sh /tmp/w327d/base host/daemon/episode_finish_test.go <164\|79> 20 <run> [-race]`; control `SPIN='for time.Now().Nanosecond() > 100_000_000 {}'`; idle `go test -race … '^TestEpisodeFinishScopeExpiry$' -count=400` (`designer-base-harness.txt`) | `:164`: rc=1, FAIL 20/20 plain and race, `episode_finish_test.go:171: expiry 1793394729: 200`. `:79`: rc=1, FAIL 20/20 plain and race, `:94: 200 {"request":…}`. Controls: PASS 20/20 at both lines. Idle: `ok … 13.495s` |
| V6 | **[r2, renamed] Mutation sensitivity of `:271`**: no finish-selected test in either package detects a change to its bounds. This is detection, not execution; execution is V15 | `grep -rlnE '/finish"\|finishRoute\(\|finishDecisionScope\|handleEpisodeFinish' --include='*_test.go'`; Q1–Q6 as V4; both `:271` arms dropped (`false \|\| now < 0`), `go test ./host/daemon/ ./cmd/ailang-worldd/ -run Finish` | route callers: `cmd/ailang-worldd/episode_finish_test.go`, `host/daemon/episode_finish_test.go` (plus `capsule_template_test.go`, a `Phase: "record"` text hit). Q1–Q6 rc=0. Both dropped: `ok host/daemon 0.765s`, `ok cmd/ailang-worldd 8.227s` |
| V7 | Option (a) costs | base test `:164` `+3601`→`+3660`; harness ×20; P4/P4b/P5/P6 ×3 (`designer-option-a-margin.txt`) | harness rc=0, PASS 20. P4 rc=0, P4b (`3659`) rc=0 (**survive**); P5, P6 rc=1 |
| V8 | Prototype: mutants, harness | `/tmp/w327d/fix` = base + `designer-prototype.diff`; `runmut.sh /tmp/w327d/fix '^TestEpisodeFinish' 1 …`; S1–S3 with scratch `zz_designer_mut_wall.go`; `harness.sh` before `:179` (`pinFinishClock`) and `:79` (`designer-prototype-mutants.txt`, `-harness.txt`) | all 16 bound mutants rc=1, with killers as §4; S0 control (no-op replace) rc=0; S1, S2 kill T1/T2/T3; S3 kills T3 `approve+0`; W1 kills AC11 `:106: 403`; TM1 kills T1 `now+1: 403, want 200`. Harness PASS 20/20 at both lines, plain and race. sha256 restored each time (`632cef2e39e2`, `9a8d3f8c2b6f`, `d641a586ebd3`) |
| V9 | Prototype size | `diff -u base fix` over the 3 files; `diffstat` | r1: `daemon.go 8`, `episode_finish.go 7`, `episode_finish_test.go 77`; 79 insertions, 13 deletions. **r2** (`finishNow` and its longer comment): `daemon.go 12`, `episode_finish.go 7`, `episode_finish_test.go 77`; 83 insertions, 13 deletions |
| V10 | Prototype packages | `go vet ./host/daemon/`; `go test ./host/daemon/ -count=1`; `-race`; `go build ./...`; `go test ./cmd/ailang-worldd/ -run Finish`; `-run '^TestEpisodeFinish' -race -count=200`; `AILANG_BIN` pinned (`designer-prototype-packages.txt`) | vet rc=0; `ok … 279.811s`; race `ok … 521.323s`; build rc=0; `ok cmd/ailang-worldd 9.946s`; `ok … 62.627s`. **Not run:** the full `go test ./...` and `verify_ail.sh` |
| V11 | Sweeps A–D; seam precedents | the greps saved as `designer-sweep{A,B,C,D}-*.txt` (the patterns are in each file's command, §7); `grep -n 'Now *func' host/coordinator/coordinator.go`; `grep -n 'now:' cmd/ailang-worldd/session_test.go`; `grep -c 'time.Now()' host/daemon/*.go` (non-test) | A 27 (control `:164` 1); B 39 (control `home_service_test.go:66` 1); C 29 (control `episode_finish.go:169` 1); D 29. `coordinator.go:52 Now func() int64`; `session_test.go:32 now: func() int64 { return now }`; 20 non-test reads in `host/daemon`, 3 of them in `episode_finish.go` |
| V12 | Row numbering | `grep -nE '^23[5-9]\. ' design_docs/world-mission.md` | only `6304:235.`, so 236 is free |
| V13 | Struct-literal daemons in tests | `grep -n '&Daemon{' host/daemon/*_test.go` | 10 sites (`binding_test.go:174,202`, `guarded_chains_test.go:24,60`, `health_exec_tools_test.go:48`, `health_exec_tool_status_test.go:80`, `home_service_test.go:41,63`, `startup_budget_test.go:36,125`); the whole package is green on the prototype (V10) |
| V14 | **[r2]** Every `Daemon` literal, repo-wide (`designer-r2-daemon-literals.txt`) | `grep -rnE '&Daemon\{\|Daemon\{\}' --include='*.go' . \| grep -v _test.go`; `grep -rnE 'daemon\.Daemon\{\|new\(daemon\.Daemon\)\|new\(Daemon\)' --include='*.go' . \| wc -l`; control `grep -rnE 'daemon\.[A-Z][A-Za-z]*\{' --include='*.go' . \| wc -l` | non-test: `host/daemon/daemon.go:784: d := &Daemon{` (inside `New`), plus 6 `presetDaemon{}` false positives in `cmd/ailang-worldd/episode_finish.go`; external constructions `0`; control `85`; test literals `10` |
| V15 | **[r2]** Execution of `:271`, by coverage, with `:169` as the positive control (`designer-r2-coverage-271.txt`) | packages whose tests import `host/daemon`: `go list -test -f '{{.ImportPath}} {{join .Imports " "}}' ./...` gives `cmd/ailang-worldd`, `cmd/world-publish` and `host/daemon.test`. For each, `go test -c -coverpkg=./host/daemon -o <bin>`, then run from the package dir in 6 round-robin shards (`-test.run '^(<names>)$' -test.timeout=1500s -test.coverprofile=…`, mode `set`), summed per block with `covsum.py`. The one-shot `go test ./host/daemon/ -coverpkg=… -timeout 570s` **timed out** (`FAIL … 571.209s`), so sharding was needed | `host/daemon` (302 tests, 6/6 shards rc=0): block `episode_finish.go:269.2,271.72` (the condition) **count 1** (shard 0, which holds AC11); the refusal block `271.72,273.3` **count 0**. Control: `166.30,169.96` count 2, and the refusal `169.96,172.4` count 1 (shard 1, which holds `TestEpisodeFinishScopeExpiry`). `cmd/ailang-worldd` (246 tests): every `episode_finish.go` block is 0, because its finish tests use a stub route (`episode_finish_test.go:156`). It had 12 FAILs, none of them finish tests and not triaged (terminal-less background binary). `cmd/world-publish`: 0 for every block; 1 FAIL, `/dev/tty` absent. **Scope:** in-process only. A daemon run as a subprocess by a CLI test is not measured (`GOCOVERDIR not set`), but no CLI test posts to a real `/finish` (V6) |
| V16 | **[r2]** The coordinator clock and where approval stamps are read (`designer-r2-clock-seams.txt`) | `sed -n 818p;835,836p host/daemon/daemon.go`; `grep -n 'cfg Config\|cfg\.Now' host/coordinator/{coordinator,effectful}.go`; `grep -rnE 'func \(c \*Coordinator\) (Now\|Clock\|Config)' host/coordinator \| wc -l` (control: any exported method); `grep -rn 'coordinator.New(' host/daemon/*_test.go \| wc -l`; `grep -rnE '\b(request\|decision)\.Now\b' --include='*.go' host cmd \| grep -v _test`; `grep -n '^func ' host/broker/approve.go`; `grep -rn 'Dispatch(' host/projection/*.go` | `:818 if cfg.AilangBin != "" {`; `:836 … Now: func() int64 { return time.Now().Unix() } …`; `coordinator.go:59 cfg Config`, reads at `:354` and `effectful.go:284`; clock accessors `0` (control `2`); `7` test-built coordinators; stamp reads only at `approve.go:629-642`, inside `validatePublishApproval` (`:517`); `findApprovalDecision(…, requestRef)` at `:252`; `Dispatch` callers `mcp.go:130`, `projection.go:363` |
| V17 | **[r2]** Prototype renamed to `finishNow` | `/tmp/w327d/fix` after the rename: `go vet ./host/daemon/`; `go test ./host/daemon/ -run '^TestEpisodeFinish' -race -count=50`; S1–S3 and W1 via `runmut.sh` (`designer-r2-prototype-seam-mutants.txt`) | vet rc=0; `ok … 15.549s`; S1, S2 kill T1/T2/T3 (`:187 scope expiry now+1: 403, want 200`); S3 kills T3 (`:231 record at approve+0: 403, want 200`); W1 kills AC11 (`:106: 403`). sha256 restored each time (`e429375dc3e0`, `e6b70d7e3f3f`) |
| V18 | **[r3, controller]** Type-aware `Daemon` construction audit (`controller-r3-daemon-construction-audit.txt`, source included) | `/tmp/w327/audit/audit <worktree>` (`go/packages` `Tests: true`, `./...`, resolves `types.Unalias` of each composite literal, `new(T)` argument and argument-less `ValueSpec` name against `host/daemon.Daemon`); control: scratch `zz_w327_auditctl_test.go` with alias `dd`, removed after | base: `packages=108 load-errors=0 constructions=11` (`daemon.go:784` + 10 `host/daemon` test literals, all composite); control: `constructions=14`, one hit per branch (`new()`, `zero-value var`, `composite-literal`); `git status` clean after |

## 10. Quorum verification log

**Round 2** (artifact `verification/world-iter327/quorum/r2.json`, 2026-10-30T21:05:42Z, $0.3578): **BLOCKED**. Same seats: `oc-kimi-k3` REJECT, `gpt6-1-sol` REJECT, `gemini-4-argon` PASS; no seat absent; controller pass. Neither reject disputes the direction (both are evidence completeness), and each carries a concrete `proposed_fix`, so the controller applied them under the **narrow-refinement carve-out** (r3, controller-authored, verbatim):
- **oc-kimi-k3:** the reviewer's own minimum acceptable revision, verbatim: §4 header "P/Q/R/TM1 kill attributions verified on the pre-rename prototype; the r2 diff's first verification is AC-2 at execution time"; §3 and AC-3 "package green verified on the pre-rename prototype; AC-3 at execution is the first package run of the r2 diff".
- **gpt6-1-sol:** "a type-aware audit at the measurement commit: enumerate composite literals, new expressions, and zero-value declarations whose resolved type is host/daemon.Daemon, including package aliases. Add a verification-log row listing each construction and whether it can reach the finish handler" → V18 (11 constructions, 0 reach the finish route except via `New`); V13/V14 are now labelled textual-search results; the no-silent-fallback behaviour is retained.
- No third quorum round (one re-quorum per doc); the planner routes from r3.

**Round 1** (artifact `verification/world-iter327/quorum/r1.json`, 2026-10-30T20:36:12Z, $0.2474):
**BLOCKED**. The seats were `oc-kimi-k3` REJECT, `gpt6-1-sol` REJECT and `gemini-4-argon` PASS.
No seat was absent. The controller's in-session reading was a pass.

- **oc-kimi-k3 (REJECT): reuse versus a new seam.** Quoted: "the doc would add a second,
  independently pinnable `Now func() int64` inside the same daemon, so finish-route tests and
  coordinator behavior can be pinned to different notions of 'now' in the same process with no
  error." **Measured before deciding (V16).**
  - (i), reuse, is impossible as the code stands: the clock is in an unexported `cfg` with no
    accessor, and `d.coord` is nil without `AilangBin`.
  - (ii), one shared clock, is possible but would put a far-future pin on every coordinator
    invocation, and it needs a heavier killer fixture.
  - **r2 chooses (iii).** Only one clock on the daemon is pinnable. It is renamed `finishNow` and
    documented as not the coordinator's clock. The measured reason is that no comparison crosses
    the two clocks: approval stamps are compared only in `validatePublishApproval`. That is §2.
  - Their catches are applied:
    - every `Daemon` literal repo-wide is measured (V14, §3);
    - AC-2 excludes TM2;
    - AC-4 now says the executor's full-suite run is the first one.
- **gpt6-1-sol (REJECT): reachability was claimed from mutation survival.** Quoted: "Surviving
  mutations demonstrate missing detection, not missing execution." **Measured: correct.** Coverage
  (V15) shows `:271`'s condition **is** executed, by AC11's record path, but its refusal branch
  never is. The `:169` control is executed in both its blocks. r2 makes three corrections:
  - "no test reaches `:271`" becomes "the condition is evaluated, the refusal branch is never
    executed, and no test detects a bound change" (§2 T3, §4);
  - V6 is renamed as a mutation-sensitivity row;
  - V15 is added.
- **gemini-4-argon (PASS), non-blocking notes, both applied.**
  - **Concurrency.** The `finishNow` reads are on the request path under `humanMu`, and tests write
    it only between synchronous calls. This is stated in the field's comment and in §2.
  - **Line numbers.** The S rows' `:115/:167/:269` are labelled as prototype lines, shifted one by
    the removed import (§4).

