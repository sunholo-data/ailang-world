# w-store-deadline-free-residue-owner — sprint plan

## Scope

Queue row 23, Go host plumbing only, based on detached `c01fb31` plus the designer's 28-file uncommitted prototype. Apply the approved [design](w-store-deadline-free-residue-owner.md) in order M1 → M2 → M3. No `.ail`, duration, runtime deadline guard, Commit contract, or policy decision is in this tranche. Zero is a syntax checkpoint; the caller must supply the intended lifetime. Mark still owns the strict guard and Commit policy under open row 23.

The prototype is the source of executable hunks. Its eight production files are `cmd/world-publish/main.go`, `host/broker/{approve,broker,publish_op}.go`, `host/daemon/daemon.go`, `host/registry/registry.go`, `host/replay/replay.go`, and `host/store/store.go`. The other 20 changed files are tests. M3 adds comments only to already touched production files.

## Baseline and staging invariant

The starting direct-read pin is `{approve.go:8, registry.go:2, replay.go:1}` (11). M1 must leave `{registry.go:2, replay.go:1}` (3). M2 must leave the **empty map** (0); M3 keeps it empty. The intermediate map edits to `host/store/context_read_test.go` are necessary for green cumulative snapshots. Keep its existing six-getter/Background regex in M1 and M2; move the final `TODO`/`GetVerifyResult` regex, empty-map explanatory comment, and correction of “six” to “seven” into M3. M2's zero is already the empty map, so §8's phrase “M3 land empty map” means finalize its guard and explanation, not postpone zero.

The source-root census changes 27 → 18 after M1 (eight getter roots and two hidden mint roots removed, one named CLI root added) → 15 after M2. Install the exact 15-root census only in M3, after the M2 roots exist. `TestProductionGoSurface` belongs with that census. Neither new M3 test may appear in M1 or M2.

## Milestones

### M1 — approvals

**Apply these prototype hunks together:**

| File | Exact M1 work |
|---|---|
| `host/broker/approve.go` | `HumanHandler.Execute` uses ctx on approve and poll; `DecideApproval`/`decideApproval`, `appendApprovalHead`, `findApprovalRequest`, `findApprovalDecision`, `walkApprovalHead`, `validatePublishApproval` take/forward ctx; replace all eight getter Background arguments. |
| `host/broker/broker.go` | `Session.invoke` forwards ctx to `validatePublishApproval`. |
| `host/broker/publish_op.go` | `MintAttendedApproval`/`mintAttendedApproval` take ctx; both session `Invoke` legs and private `decideApproval` receive it. |
| `cmd/world-publish/main.go` | `runApprove` supplies its explicit named `context.Background()` compatibility root and its two-line reason. |
| `host/broker/approve_test.go`, `episode_test.go`, `handlers_test.go`, `publish_op_test.go` | Apply all approval signature call-site edits. |
| `host/broker/registry_publish_test.go` | Apply only the five approval call-site edits (`decideAndPollLandedApproval`, `landApprovalOverRawRequest`, `landNonPublishApproval`, and two edits in `TestPublishApprovalRefusalSetWithALandedPositiveControl`). **Leave `driveReplayEntry` unchanged until M2.** |
| `host/broker/context_thread_test.go` | Add entire file: `TestApprovalContextIdentity` (both lifetime subtests, 5 object + 4 head mint reads and 2 publish object reads) and `TestApprovalCallerCancellation` (approve, poll, decide, mint, publish subtests). |
| `host/store/context_read_test.go` | Change only `deadlineFreeReadPins` to `{"host/registry/registry.go":2, "host/replay/replay.go":1}`; retain the original scanner and regex. |

**Boundary:** all `_test.go` files compile; `TestNoNewDeadlineFreeStoreReads`, `TestApprovalContextIdentity`, and `TestApprovalCallerCancellation` pass. Outside the sandbox, full `./host/broker ./cmd/world-publish ./host/store` suites pass with the pinned binary. AC1, AC5, AC8. Kills M01–M08, M14–M15, M19–M22.

### M2 — bootstrap and replay/cache

**Apply these prototype hunks together:**

| File | Exact M2 work |
|---|---|
| `host/registry/registry.go` | `Bootstrap(ctx, *store.Store, release)` delegates to private `bootstrap(ctx, bootstrapStore, release)`; both head/object reads use ctx. |
| `host/registry/registry_test.go`, `context_thread_test.go` | Migrate seven existing Bootstrap calls; add `TestBootstrapCallerCancellation` and `TestBootstrapSecondReadCancellation`. |
| `host/daemon/daemon.go` | `New(ctx, cfg)`, `Run`→`New(ctx,cfg)`, `New`→`registry.Bootstrap(ctx,...)`. |
| `host/daemon/bench_test.go`, `daemon_test.go`, `handlers_test.go`, `integrity_test.go`, `lock_wait_order_test.go`, `context_thread_test.go` | Migrate every `New` call, including benchmarks; add `TestStartupCallerCancellation` with `New` and `Run` subtests and writer reopen control. |
| `host/replay/replay.go` | Private `replayStore` seam; `ReplayEntry(ctx,...)` source/cache reads; `ReplayEpisode(ctx,...)` forwards to entry. Preserve `runPinnedTransition`'s separate subprocess root. |
| `host/replay/replay_test.go`, `context_thread_test.go` | Migrate all entry/episode/cache calls, including the Bootstrap call in `TestEpochRegistryCandidateCannotRedirect`; add `TestReplayCallerCancellation` and `TestReplayCacheCallerCancellation`. |
| `host/store/store.go`, `store_test.go`, `verify_context_test.go` | `GetVerifyResult(ctx,...)` with `QueryRowContext`, migrate existing cache tests, add `TestVerifyResultCallerCancellation`. The prototype's final-blank-line deletion in `store.go` is incidental. |
| `host/broker/registry_publish_test.go` | Apply the deferred `driveReplayEntry` call-site hunk. This file is deliberately split across M1/M2. |
| `host/store/context_read_test.go` | Set `deadlineFreeReadPins = map[string]int{}` now; keep M3's scanner/comment edits deferred. |

**Boundary:** all `_test.go` files compile; `TestNoNewDeadlineFreeStoreReads` reports zero; `TestBootstrapCallerCancellation`, `TestBootstrapSecondReadCancellation`, `TestStartupCallerCancellation`, `TestReplayCallerCancellation`, `TestReplayCacheCallerCancellation`, and `TestVerifyResultCallerCancellation` pass. Outside the sandbox, full `./host/daemon ./host/registry ./host/replay ./host/store ./host/broker` suites pass with the pinned binary. Broker is included because its deferred replay call now compiles against the changed API. AC2–AC4, AC6–AC7, AC13. Kills M09–M13, M16–M18, M28.

### M3 — honest zero and handoff

**Apply these prototype hunks:** `host/store/context_read_test.go` final zero-ratchet comment and regex; entire `host/store/context_roots_test.go`, containing `contextRootPins`, `TestProductionContextRoots` (15 named roots, Background/TODO/WithoutCancel counts, renamed import handling, dot-import refusal, package declaration scan, ≥58 floor and seven anchors) and `TestProductionGoSurface` (host/cmd plus the two exact required reproducers). The surface guard walks untracked files and rejects new non-test Go trees. Final scanner wording should call the getter set **seven**, since `GetVerifyResult` is added.

**API usage note, missing from the prototype:** add these exact one-line sentences to the existing exported doc comments, seven comment lines total. They are comments only and supply no default duration:

| Location | Add to its exported doc comment |
|---|---|
| `host/broker/approve.go`, `DecideApproval` | `// Pass the caller's intended read ctx; this API supplies no default timeout.` |
| `host/broker/publish_op.go`, `MintAttendedApproval` | `// Pass the caller's intended read ctx; this API supplies no default timeout.` |
| `host/daemon/daemon.go`, `New` | `// Pass the caller's intended bootstrap-read ctx; this API supplies no default timeout.` |
| `host/registry/registry.go`, `Bootstrap` | `// Pass the caller's intended read ctx; this API supplies no default timeout.` |
| `host/replay/replay.go`, `ReplayEntry` and `ReplayEpisode` (each) | `// Pass the caller's intended source/cache-read ctx; this API supplies no default timeout.` |
| `host/store/store.go`, `GetVerifyResult` | `// Pass the caller's intended cache-read ctx; this API supplies no default timeout.` |

For broker mint, replay entry, and daemon New, retain the prototype's existing narrower comments on writes, subprocess lifetime, and startup. These seven public API notes meet §11's migration instruction at the functions whose signatures changed; the test fixtures remain executable examples.

**Boundary:** all `_test.go` files compile; `TestNoNewDeadlineFreeStoreReads`, `TestProductionContextRoots`, and `TestProductionGoSurface` pass; full pinned out-of-sandbox `go test ./... -count=1` passes. AC9–AC12, AC14–AC22. Kills M23–M27, M29–M37. The four deadline-free guard impact paths and Mark's §5 A/B choice are documented handoff, with no runtime guard in this milestone.

## Deviations from design text

| Design text / prototype | Justification or correction |
|---|---|
| §8 says M2 reduces pins to zero while M3 “lands” the empty map. | The scanner compares maps exactly; zero must be an empty map in M2. M3 strengthens the regex and explanation. |
| §11 says production package comments carry API usage notes; prototype has only partial `ctx governs` comments and no “no default timeout” note for any changed API. | Add the seven exact exported-comment lines above in M3; no logic change. |
| Final `context_read_test.go` comment says “six context-first read getters” while its regex names seven. | Correct to seven in M3 so the guard description matches its expression. |
| §8 M3 names AC9–AC12/AC14, while the final §6 includes AC15–AC22. | M3 owns all scanner/surface witnesses, including the two controller mutations. |
| §8 M1's broker/CLI package gate and M2's daemon package gate cannot be run in this sandbox. | Run the named non-socket tests locally and require the full package gates outside the sandbox; socket failures here would be uninformative. |
| Prototype `store.go` deletes a trailing blank line. | Formatting-only, no behavior; it may ride M2's cache hunk. |
| §2's earlier F4 claim of three exported `DecideApproval` test calls is false. | The final design itself corrects this to one actual call; two textual hits are a comment and diagnostic. |

No milestone merge is required. The only shared-file split is `registry_publish_test.go` (approval hunks M1, replay hunk M2) and the three successive pin-map states in `context_read_test.go`. Applying entire final versions of either file in M1 would break compilation or the ratchet. `TestProductionContextRoots` cannot run until M2's 15-root state exists.

## Mutation tables

Rows M01–M35 are copied by ID, operation, and named killer from `mutations.json` / `mutations-run.txt`; every listed result is **KILLED**. M36/M37 are from the controller transcripts, each **KILLED**, restored and green on rerun. “Exact context” means `TestApprovalContextIdentity` and/or `TestApprovalCallerCancellation` as the JSON selector states.

### M1 — 14/14 killed

| ID | Mutation | Named killer test(s) |
|---|---|---|
| M01 | approval decision request GetObject → Background | `TestApprovalContextIdentity`, `TestApprovalCallerCancellation` |
| M02 | appendApprovalHead GetRegistryHead → Background | same two approval tests |
| M03 | walkApprovalHead GetRegistryHead → Background | same two approval tests |
| M04 | walkApprovalHead chain GetObject → Background | same two approval tests |
| M05 | walkApprovalHead decision GetObject → Background | same two approval tests |
| M06 | walkApprovalHead request GetObject → Background | same two approval tests |
| M07 | publish decision GetObject → Background | same two approval tests |
| M08 | publish request GetObject → Background | same two approval tests |
| M14 | Execute approve calls appendApprovalHead with Background | same two approval tests |
| M15 | Execute poll calls findApprovalDecision with Background | same two approval tests |
| M19 | DecideApproval wrapper calls private decide with Background | same two approval tests |
| M20 | MintAttendedApproval wrapper calls private mint with Background | same two approval tests |
| M21 | private mint calls decideApproval with Background | same two approval tests |
| M22 | Session.invoke calls validatePublishApproval with Background | same two approval tests |

### M2 — 9/9 killed

| ID | Mutation | Named killer test(s) |
|---|---|---|
| M09 | Bootstrap head read → Background | `TestBootstrapCallerCancellation`, `TestBootstrapSecondReadCancellation` |
| M10 | Bootstrap object read → Background | same two Bootstrap tests |
| M11 | ReplayEntry source read → Background | `TestReplayCallerCancellation`, `TestReplayCacheCallerCancellation` |
| M12 | ReplayEntry cache getter → Background | `TestReplayCacheCallerCancellation` |
| M13 | cache `QueryRowContext` → `QueryRow` | `TestVerifyResultCallerCancellation` |
| M16 | New calls Bootstrap with Background | `TestStartupCallerCancellation` |
| M17 | Run calls New with Background | `TestStartupCallerCancellation` |
| M18 | ReplayEpisode calls ReplayEntry with Background | `TestReplayCallerCancellation` |
| M28 | GetVerifyResult swallows the cache error | `TestVerifyResultCallerCancellation` |

### M3 — 14/14 killed

| ID | Mutation | Named killer test(s) |
|---|---|---|
| M23 | direct-read scanner searches no host/cmd roots | `TestNoNewDeadlineFreeStoreReads` |
| M24 | root census searches no host/cmd roots | `TestProductionContextRoots` |
| M25 | hoist Background into Execute's ctx | `TestProductionContextRoots` |
| M26 | direct getter takes `context.TODO()` | `TestNoNewDeadlineFreeStoreReads` |
| M27 | direct getter takes `context.Background()` | `TestNoNewDeadlineFreeStoreReads` |
| M29 | root census skips `approve.go` | `TestProductionContextRoots` |
| M30 | renamed context import `c` and `c.Background()` | `TestProductionContextRoots` |
| M31 | add `context.WithoutCancel(ctx)` | `TestProductionContextRoots` |
| M32 | function-valued `context.Background` | `TestProductionContextRoots` |
| M33 | package-level `context.Background()` | `TestProductionContextRoots` |
| M34 | standalone `context.TODO()` | `TestProductionContextRoots` |
| M35 | new top-level non-test Go tree | `TestProductionGoSurface` |
| M36 | dot-import context, bare `Background()`; counts stay 15/0/0 | `TestProductionContextRoots` (`dot context import`) |
| M37 | remove allow-listed `repro/main.go` | `TestProductionGoSurface` (`allow-listed reproducer missing`) |

## Test discipline

The identity tests compare the exact context object and read counts. Cancellation tests use `errors.Is(context.Canceled)` and live controls, not elapsed time. The direct getter regex is intentionally syntactic; the independent AST root census catches the drilled hoists, aliases, detached contexts, package roots, and dot imports. The surface guard covers new production Go trees and requires the two existing reproducers. None proves an actual deadline, all possible dataflow, SQLite lock-wait dominance, writes, full startup, or replay subprocess termination.

All milestone Go commands must be bounded. Use `-timeout 300s` on `go test`; put a finite process timeout around `go vet` and the entire gate sequence. Use `AILANG_BIN=$HOME/.pinned-ailang/ailang` for behavioral suites outside the sandbox. Do not run socket-dependent broker, daemon, or cmd suites here as evidence; any accidental result is **UNINFORMATIVE UNDER SANDBOX**. The controller's out-of-sandbox final-prototype gate already reports 23 `ok`, 0 `FAIL`.

## Prototype gate results

Executed once in this worktree with a finite Python subprocess timeout around each prescribed command (and `-timeout 300s` on both Go tests):

| Gate | Observed |
|---|---|
| `go vet ./...` | `vet rc=0`, no diagnostics |
| `go test ./... -run '^$' -count=1 -timeout 300s` | `compile rc=0`; 23 `ok`, 0 `FAIL` (all `_test.go` compiled) |
| `go test ./host/store ./host/registry ./host/replay -count=1 -timeout 300s` | `test rc=0`; three `ok`, 0 `FAIL` |

Verbatim targeted behavioral `ok` lines:

```text
ok  	github.com/sunholo-data/ailang-world/host/store	6.052s
ok  	github.com/sunholo-data/ailang-world/host/registry	0.359s
ok  	github.com/sunholo-data/ailang-world/host/replay	0.685s
```

Verbatim compile `ok` lines (all 23):

```text
ok  	github.com/sunholo-data/ailang-world/cmd/ailang-worldd	0.344s [no tests to run]
ok  	github.com/sunholo-data/ailang-world/cmd/world-publish	0.616s [no tests to run]
ok  	github.com/sunholo-data/ailang-world/host/archive	0.810s [no tests to run]
ok  	github.com/sunholo-data/ailang-world/host/authority	1.268s [no tests to run]
ok  	github.com/sunholo-data/ailang-world/host/boundary	1.017s [no tests to run]
ok  	github.com/sunholo-data/ailang-world/host/broker	1.218s [no tests to run]
ok  	github.com/sunholo-data/ailang-world/host/canon	1.620s [no tests to run]
ok  	github.com/sunholo-data/ailang-world/host/capsule	1.856s [no tests to run]
ok  	github.com/sunholo-data/ailang-world/host/childenv	2.687s [no tests to run]
ok  	github.com/sunholo-data/ailang-world/host/daemon	3.226s [no tests to run]
ok  	github.com/sunholo-data/ailang-world/host/evidence	2.085s [no tests to run]
ok  	github.com/sunholo-data/ailang-world/host/hashref	2.279s [no tests to run]
ok  	github.com/sunholo-data/ailang-world/host/pkgproj	1.419s [no tests to run]
ok  	github.com/sunholo-data/ailang-world/host/procbound	3.412s [no tests to run]
ok  	github.com/sunholo-data/ailang-world/host/proctest	2.485s [no tests to run]
ok  	github.com/sunholo-data/ailang-world/host/projection	2.944s [no tests to run]
ok  	github.com/sunholo-data/ailang-world/host/registry	3.632s [no tests to run]
ok  	github.com/sunholo-data/ailang-world/host/replay	3.577s [no tests to run]
ok  	github.com/sunholo-data/ailang-world/host/runbook	3.577s [no tests to run]
ok  	github.com/sunholo-data/ailang-world/host/store	3.621s [no tests to run]
ok  	github.com/sunholo-data/ailang-world/host/transitionreg	3.615s [no tests to run]
ok  	github.com/sunholo-data/ailang-world/host/verifygate	3.548s [no tests to run]
ok  	github.com/sunholo-data/ailang-world/host/workbench	3.570s [no tests to run]
```

## Execution notes

Executor applies prototype hunks, not rewrites. The executor writes **no commits**: the controller builds one commit per milestone from cumulative `.snap/M<k>/` snapshots. Stage partial hunks exactly as above, including both pin-map transitions. Do not use `git add -A`, `git checkout --`, or the designer's `mutate.py` restore in this uncommitted worktree; its `finally` rewrites prototype files. The mutation table is already measured. M3's seven API-note lines are the sole non-prototype production additions authorized by this plan; the scanner's six→seven comment correction is test-only.

At **each** boundary run bounded `go vet ./...` and bounded `go test ./... -run '^$' -count=1 -timeout 300s`. Then run the narrowest relevant behavioral packages outside the sandbox with `AILANG_BIN` and `-count=1 -timeout 300s`: M1 `./host/broker ./cmd/world-publish ./host/store`; M2 `./host/daemon ./host/registry ./host/replay ./host/store ./host/broker`; M3 `./host/store` for the three scanner tests, then full `./...` as AC14. In the sandbox, use only named non-socket tests for broker/daemon and the full store/registry/replay packages. Capture each snapshot's pin log and `ok`/`FAIL` lines. For M3, record the 15-root table and two-reproducer surface, plus the four-path guard impact from design §4; do not turn the guard on.

## Verification Log

| Claim | Command run | Observed output |
|---|---|---|
| Prototype manifest | `git status --porcelain`, `git diff --numstat` | 22 modified + 6 new Go files; eight modified production Go files. |
| Design and precedent | `cat design_docs/planned/w-store-deadline-free-residue-owner.md`; `cat design_docs/implemented/w-daemon-lock-wait-not-deadline-bound-sprint-plan.md` | §2/§3/§6/§7/§8/§10/§11 and precedent sections read; design's final manifest is 28 files. |
| Prototype sources/tests | `git diff --no-ext-diff -- ...`; `cat` of the six new test files; targeted source reads; Python `read_bytes()` over every Go path from `git status --porcelain` | Read 28 Go files, 564294 bytes; confirmed shared `registry_publish_test.go` has approval and replay hunks, final ratchet map is empty, and root test pins 15. |
| Designer mutations | `cat .../mutations.json .../mutations-run.txt` | M01–M35 each `KILLED`; `TALLY 35 35`. |
| Controller mutations | `cat .../M36-controller.txt .../M37-controller.txt` | M36 fails dot-import rule with counts 15/0/0; M37 fails missing repro producer; both have `--- FAIL:`. |
| API-note absence with positive control | `rg -n 'read ctx only|no default timeout' host --glob '*.go'; rg -n 'ctx governs|ctx is threaded' host/broker/publish_op.go host/replay/replay.go host/daemon/daemon.go` | First search: no match; positive control finds partial comments at publish_op.go:123, replay.go:157, daemon.go:454. |
| Local gates | finite Python `subprocess.run` wrapper, three commands in Prototype gate results | vet 0; compile 0 with 23 `ok`/0 `FAIL`; target test 0 with three `ok`/0 `FAIL`. |
| Controller full suite | `cat /Users/voightkampff/.ailang/state/world-iter194/ctl-proto-fulltest.txt` | 23 `ok`/0 `FAIL` outside sandbox on final prototype; controller separately reports vet rc=0. |

The false design/controller claims to carry forward are the already corrected F4 “three exported test calls,” the §11 claim that the prototype already has package-level API guidance, and the final scanner's “six getters” comment. No local gate contradicted the controller's 23-package full-suite result.
