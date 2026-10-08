# Row 153 verification evidence (executor side)

Design: `design_docs/planned/w-plan-phase-deadline-under-load.md`. Plan: `...-sprint-plan.md`. Logs: `executor-run1.md` (C0-C4b,
baseline, run-1 mutation ledger), `executor-run2.md` (C5, C6, run-2 ledger), `mut-run2.log` (verbatim), `proxy/` (scripts and
trimmed logs). Not here: AC4.2 (5/5 consecutive CI attempts) and AC4.3 (records), which belong to the controller.

## AC to evidence

| AC | Evidence |
|---|---|
| AC1.1-AC1.7 (typed plan/finish timeout, wire, diagnosable rig) | run 1: C1/C2, mutation rows MUT-RELABEL ... MUT-SILENT-ERROR-b |
| AC2.1-AC2.4 (derived 4 s cap) | run 1: C3, `TestPlanPhaseBudgetIsDerived`, `TestSlowPlanLeavesHandlerItsCap` |
| AC3.1 build/error/racer/retry | run 1: C4a tests; run 2: `TestPublishWithoutATemplateIsRefused` (MUT-SOFT-PUBLISH), GC (`TestTemplateGCPrunesStaleInterpDigests`, `TestPruneCapsuleCacheKeepsListedDigests`; MUT-NO-GC, MUT-GC-*) |
| AC3.2-AC3.6 (copy, byte equality, warmth witness, tripwire) | run 1: C4b |
| AC3.7 (i)-(v) startup status, `CapsuleTemplates()`, rebuild after listener-ready, cold label, stub refusal | run 2: `TestDaemonStartReportsMissingTemplates` + `TestPublishWithoutATemplateIsRefused`; MUT-SILENT-COLD, MUT-BLOCKING-PREWARM(-b), MUT-SOFT-PUBLISH, MUT-STATUS-LIE, MUT-NO-REBUILD |
| AC3.8, AC4.1 (throttled-load proxy) | table below |

## Gates

Every commit boundary C0-C5: `verify_ail` rc 0, `go vet` rc 0, `go test -race` rc 0, srt list 20/20 `--- PASS`, `gofmt` empty, no
binary blobs. `go build && go test ./...` is rc 1 at **every** boundary including the untouched base C0, solely for `host/broker`
`TestRunBoundedHeadTailTimeoutKeepsPartialOutput` (a 0.30 s test that fails under parallel package load on this host, passes
alone; not row 153). Per-boundary tables: `executor-run1.md`, `executor-run2.md`.

## Proxy: `proxy/bg_load.sh 8 4` on the C5 head (race build, probe-instrumented scratch copy, `PlanPhaseBudget` = shipped 4 s)

Control: the design's V13 on the base tree, `proxy/v13-base-bg8-control.trimmed.log`.

| Reading | Base (V13, 2 s cap, cold) | C5 head (this run) |
|---|---|---|
| Tests that reached a plan | 4 of 6 | 6 of 8 (the other 2 failed at publication, below) |
| Plan phase `took`, min-max | 2.03-2.96 s | **1.33-1.66 s** (12 plans, all `err=<nil>`) |
| Plan timeouts | **4 of 4** (`context deadline exceeded`) | **0 of 12** |
| `cold compile` lines | n/a (no label) | **0** (every capsule run `cold=false`) |
| First-plan child (V15: 1.96-2.98 s cold) | 1.96-2.98 s (V15, uncapped) | **0.50-0.82 s** (<= 1 s, AC3.8c) |
| Interpreter verify step | 376-660 ms (V15) | 526-720 ms |
| Startup probe (AC3.8a) | n/a | PASS in all 6 tests that reached serving |
| Per-descriptor template status after publish (AC3.8b) | n/a | 9 of 9 `ready` in every one of the 6 (`ailang-check ailang-cli ailang-edit ailang-read ailang-run ailang-write builtins-search examples-search workspace-exec`) |
| Publish duration (9 sources, serial, under 8 burners) | per-source check 3.1-5.0 s (V34) | **28.4-28.9 s** total = 3.2 s/source, inside V34's per-check range |
| Time to serving (rig setup, AC3.8d) | not probed on base (tests ended 47-57 s, failing at the plan) | 48.1-52.4 s |
| Test outcome | 6 FAIL | 6 PASS, 2 FAIL |

The 2 FAILs of the C5 run are `PublishSet ... is not loadable under its pinned interpreter` at `exec_e2e_srt_test.go:143/210`:
publication (28.4 s typical) exceeded the rig's 30 s `boundedTestContext` once per test. This is the V16 class (the base run had 2 of
6 failing at publication too), a rig-only load failure R3 recorded, with no plan and no template involved.

**12 burners** (`proxy/m3-head-bg12-count3.trimmed.log`, V16 control): 6 of 6 FAIL before any plan: 3 startup-probe failures
(`arm7-toolchain` subprocess timeout x2, `arm6-network` x1) and 3 publication failures. 0 plan timeouts, 0 plans reached, 0
`cold compile`, no template failure. This is the AC3.8 requirement (failures only at the startup probe or publication).

Readings that are NOT evidence: AC4.1's "0 plan timeouts and 0 cold lines" holds on this proxy, but the proxy is one Apple-silicon
host under QoS throttling, not the CI runner. The base-tree publish duration and time-to-serving were not measured (the base
probe did not print them), so (b) and (d) are compared with V34 and the test wall only. The CI race leg is the controller's
measurement (AC4.2, row 154).

## Files

- `proxy/bg_load.sh`: the parameterised script that ran (burners, count, out, test binary, package dir).
- `proxy/probe-phase-timing.m3.patch.txt`: the hand-adapted probe (uncommitted in source; archived as inert text).
- `proxy/m3-head-bg8-count4.trimmed.log`, `proxy/m3-head-bg12-count3.trimmed.log`, `proxy/v13-base-bg8-control.trimmed.log`.
