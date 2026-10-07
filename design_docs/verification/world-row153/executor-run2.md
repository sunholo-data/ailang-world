# Row 153 executor run 2 log (C5, C6)

Executor: mission iteration 241, claude-sonnet-5-5, unattended. Branch `sprint/row153-plan-phase-deadline-under-load`, base
`b6f7dfb` (run 1's head). Raw logs: `~/.ailang/state/mission-world-iter241/executor-run2/` (`c5/` gates and mutation record,
`m3-bg8.log`, `m3-bg12.log` proxy runs, `scratch/` the probe-instrumented copy, never committed).

## Commits

| Plan | SHA | Subject |
|---|---|---|
| C5 | 06970ca | M3b: per-descriptor capsule template status, rebuild after listener-ready, template GC |
| C6 | (this commit) | M4: `design_docs/verification/world-row153/` evidence |

## What C5 added

- `host/daemon/templates.go` (no exec, coupling C1): `TemplateStatus`, `Daemon.CapsuleTemplates()`, `reportTemplates` (called
  from `New` after `bootstrapRegistry`: reads the registry head, one `CapsuleTemplateReady` stat per descriptor, one
  `ailang-worldd: capsule template <id>: ready|missing` line), `startTemplateMaintenance` (one goroutine, idempotent, started by
  `Run` after `Listen` and the announce write and before `Serve`), `stopTemplateMaintenance` (cancel plus a bounded join, called by
  `Close`), `maintainTemplates` (GC, then rebuild of `missing` entries at concurrency 2 via `archive.CheckSource`, each bounded by
  its 10 s check), and the unexported `Config.templateRebuildGate` test seam (nil in production).
- `host/archive/capsulecache.go`: `PruneCapsuleCache(keep, tmpMaxAge, now)` and `StaleTemplateTmpAge = 2 x checkTimeout` (file
  I/O only). GC keeps the digests of (daemon pin plus every head descriptor's interpreter) (P11), prunes only 64-hex digest
  directories and `.tmp-*` older than 20 s, and is skipped when the head or a source could not be read.
- QUICKSTART: the status, rebuild and GC lines.
- Tests: `TestDaemonStartReportsMissingTemplates` (AC3.7 i-iv, drives `Run`, rebuild held by the gate),
  `TestPublishWithoutATemplateIsRefused` (AC3.7 v / AC3.1), `TestTemplateGCPrunesStaleInterpDigests`,
  `TestPruneCapsuleCacheKeepsListedDigests` (archive).

## Gate results per boundary (rc; wall in s)

| Boundary | g1 verify_ail | g2 vet | g3 build+test | g4 race | g5 srt (-p 1 -v) | gofmt | blobs |
|---|---|---|---|---|---|---|---|
| C5 06970ca | 0 (23) | 0 | 1* (258) | 0 (518) | 0, 20/20 | empty | empty |
| C6 (docs only; code tree == C5) | 0 | 0 | carried from C5 (no code change) | carried from C5 | carried from C5 | empty | empty |

`*` = only `host/broker` `TestRunBoundedHeadTailTimeoutKeepsPartialOutput` (0.30 s budget), the run-1 baseline red (F1), still
red under the parallel `go test ./...` here and not attributable to row 153; every other package `ok`.
New-test PASS evidence: `c5/new-plain.log` (3 PASS), `c5/new-race.log` (3 PASS under `-race`), `c5/new-archive.log` (1 PASS).
Local supplement (`WORLD_TOOL_AILANG_BIN` = v0.52.1, `go test ./host/daemon/ ./cmd/ailang-worldd/ -v`): rc=0, 238 PASS, 0 FAIL.
Its 6 SKIPs are the three baseline ones plus the three srt e2e tests (that invocation did not export
`WORLD_EXEC_SRT_NODE_MODULES`; they are PASS in g5). One `cold compile` string appears in it: the deliberate cold call inside
`TestDaemonStartReportsMissingTemplates` (the test logs daemon #2's operator log); no rig tripwire fired.
`bench_worldd.sh --smoke` rc=0.

Per-package wall (g3 plain / g4 race), the row-154 headroom figure:

| Boundary | coordinator | projection | daemon | archive |
|---|---|---|---|---|
| C4b (run 1) | 52.0 / 54.7 | 16.8 / 19.6 | 229.3 / 255.5 | n/a |
| C5 | 53.1 / 53.8 | 16.1 / 19.7 | 232.8 / 257.7 | 41.2 / 35.3 |

Race `daemon` +2.2 s over C4b (255.5 to 257.7), whole race leg 512 to 518 s. C5 adds three daemon tests (~3 s each under load).
The CI runner is the measurement (row 154).

## Mutation ledger (C5; verbatim record in `mut-run2.log`)

Method as in run 1: a one-line edit to C5's own production file, the named test `-count=1 -v`, restore from a byte-identical
backup (`cmp` clean after every batch).

| Mutation | File | Test | Red? | Failing assertion |
|---|---|---|---|---|
| MUT-SILENT-COLD (startup): status `Fprintf` -> `_ = state` | daemon/templates.go | TestDaemonStartReportsMissingTemplates | yes | `startup lines = map[]` (:133) |
| MUT-BLOCKING-PREWARM as planned: `d.maintainTemplates(ctx)` in `New` | daemon/daemon.go | same | yes, but at the (ii) arm, not the announce bound (dev. 2) | `CapsuleTemplates() ... ailang-read:rebuilt, want missing` (:130) |
| MUT-BLOCKING-PREWARM-b: `d.maintainTemplates(ctx)` before `d.Listen()` in `Run` | daemon/daemon.go | same | yes, at the announce bound | `Run announced nothing within 7s with the template rebuild held` (:175) |
| MUT-SOFT-PUBLISH: `false && (manifest check)` | archive/check.go | TestPublishWithoutATemplateIsRefused | yes | `PublishSet ... = <nil>, want *archive.TemplateBuildError` |
| MUT-NO-GC: `if true { continue }` before the prune | archive/capsulecache.go | TestTemplateGCPrunesStaleInterpDigests | yes | stale digest and `.tmp-old` not pruned |
| MUT-GC-KEEP-INVERTED: `keep[name]` | archive/capsulecache.go | TestPruneCapsuleCacheKeepsListedDigests | yes | removed the two kept digests |
| MUT-GC-FRESH-TMP: age test negated | archive/capsulecache.go | TestTemplateGCPrunesStaleInterpDigests | yes | `a fresh .tmp-new ... was pruned` |
| MUT-NO-REBUILD: skip `rebuildTemplate` | daemon/templates.go | TestDaemonStartReportsMissingTemplates | yes | timed out (60 s) waiting for the rebuilt line |
| MUT-STATUS-LIE: every descriptor `ready` | daemon/templates.go | same | yes | `CapsuleTemplates() ... ailang-read:ready, want missing` |

Every mutation was restored before the gates ran; the gates ran on the committed tree.

## Proxy (AC3.8 / AC4.1)

See `README.md` (table) and `proxy/`.

## Deviations from the plan (each measured)

1. **Status and rebuild cover every descriptor, not only effectful ones.** The plan step and AC3.7(i) say "effectful
   descriptor". The capsule runs a pure transition from the same template mechanism (design U4), so a missing template for one
   would run cold and go unreported. The se-tools head is 9 effectful descriptors, so (i)'s count is unchanged; the superset is the
   safer reading.
2. **MUT-BLOCKING-PREWARM as planned reds at the (ii) arm first** (the synchronous rebuild in `New` flips the status to `rebuilt`
   before `CapsuleTemplates()` is read), so the announce-bound tooth was not the one that fired. MUT-BLOCKING-PREWARM-b (the same
   synchronous call placed in `Run` before `Listen`) is the mutant only the announce bound can catch; both are red.
3. **The rebuild does not skip a template that became ready after startup.** In the AC3.7 flow the cold call (iv) promotes the
   template lazily, so by the time the gate is released the template exists; the planned `rebuilt in` line is still emitted
   because the rebuild re-checks unconditionally (`PromoteTemplate` treats the racer as success). Cost: one extra `check` per
   missing template.
4. **GC is skipped when the head or any descriptor source cannot be read** (it could not enumerate what to keep). Not in the plan;
   fail-safe, it only widens what survives.
5. **Probe patch.** `probe-phase-timing.patch` was adapted by hand (the code it touched changed in M1-M3) and keeps
   `PlanPhaseBudget` at the shipped 4 s (the original raised it to 20 s for V15), so the proxy measures the real cap. It also
   times publication, prints per-descriptor template state, and adds the cold flag. It lives in a scratch copy of the C5 tree
   (`git archive HEAD`), is archived as `proxy/probe-phase-timing.m3.patch.txt` (inert text), and was never applied to this branch.
6. **Proxy script parameterised** (binary and package dir as arguments), as the plan asked; the 12-burner arm used count 3
   (as V16), the 8-burner arm count 4.
