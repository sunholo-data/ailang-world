# w-episode-outcome-capture (queue row 182, clause 5) — r1

**Status:** r1 — **RATIFIED by Mark Edmondson, attended 2026-10-16 ("ratified")**. Mark chose the shape before drafting: the
harness records outcomes, for both arms; ratification is Mark's, attended, followed by an
independent judge on the implementation (no model quorum: the OpenAI and Ollama seats are out of
quota).

## §1 Problem (measured)

Row 114 Q1 (PR #278) asked a World walk why codex's World arm failed `config_file_parser` 0/3 in
the 2026-10-14 FINAL. The walk answered "could not determine" (UNVERIFIED). Measured causes, from
its 55 commands (`Q1-config-file-parser/tool-calls.json`) and the K1 copy:

| id | measurement |
|---|---|
| M1 | K1 holds 15 semantic ids: effect intent/outcome/record/request/result/plan, invocation input/output/record, journal intent/outcome, transition source, genesis goal, epoch and transition registries. No object names a **task**, an **outcome** or a **grading rule**. |
| M2 | An invocation record carries `episodeId`, `skillId`, input/output/plan/effects. One episode spans all 23 tasks of a run, so nothing maps entries to tasks. |
| M3 | The walk ran `object find config_file_parser` (empty), then bisected entry indexes. Codex's attempts are at entries 99–101 (r1), **259–261 (r2)** and 428–430 (r3); the walk read 250 and then 262, so it never saw r2. |
| M4 | In r1 and r3 it found the decisive facts (an `ailang-write` of `app_config.json`, a solution that only reads it, `ailang-run` printing the expected line). It could not call them a cause: nothing says the grader runs only `benchmark/solution.ail` in an empty directory, or that the task failed. |
| M5 | The harness already holds everything needed at grade time: `run_cell` (`scripts/floor/run.py`) has the grade (class, stdout/stderr head, `solution_sha256`, argv via `grade.build_argv`), and for World rows `drv['world']['log_from'/'log_to']` = the task's entry range (`run_task.py:282,305`). |
| M6 | A one-object REST commit already exists: `demoProofWriter.PutObject` (`cmd/ailang-worldd/demo_proof.go`, row 150 M3a, judged 90) builds the next world and entry with `demoMakeCommit` and posts `/v1/commit` under a session token, mapping uncertain outcomes to `store.UncertainError`. |
| M7 | `POST /v1/commit` requires a resolved session binding (`handlers.go` handleCommit, 401 otherwise). Each FINAL episode's 0600 session file is on the rig (`private/<ep>.session`) and the harness already reads it in the paired phase (`cmd_final`). |
| M8 | `why` tolerates entries whose `transitionRef` is not an invocation record (`why.go` `record()` returns nil), and `why sha256:<ref>` dispatches unknown semantic ids to a transitionRef scan, so new entries do not break existing walks. |

## §2 Design

**D1 — the object.** `world/task-outcome/v1`, canonical JSON (sorted keys, no insignificant
whitespace), bounded at 16 KiB:

```json
{"agent":"codex","arm":"world","episodeId":"fl-cx-r2","run":2,"label":"final-paired-r2","task":"config_file_parser",
 "entries":{"from":259,"to":261},
 "outcome":{"class":"FAIL","errorCategory":"logic_error","compileOk":true,"runtimeOk":true,"stdoutOk":false,
            "stdoutHead":"configuration file not found\n","stderrHead":"","solutionSha256":"…","solutionState":"present"},
 "grading":{"rule":"only benchmark/solution.ail is copied into an empty temporary directory, plus the task's declared input_files; then run","argv":["run","--entry","main","--quiet","--relax-modules","--caps","IO,FS","benchmark/solution.ail"],
            "inputFiles":[],"expectedStdout":"Loaded MyApp v1.0.0 on port 8080 with 3 features\n","graderVersion":"floor-grade/1","toolSha256":"…"}}
```

`entries` is `null` for a shell row (no World calls). `stdoutHead`/`stderrHead` are the row's
existing 400-byte heads. `expectedStdout` is the task spec's, so a walk can compare. The
`semanticId` is the constant `world/task-outcome/v1` (discoverable with `object find`), and
`provenance` is `harness:floor`.

**D2 — the verb.** `ailang-worldd [--addr] outcome record --file <outcome.json|-> [--session <file|token>]`.
It strict-decodes the payload (unknown fields refused, required fields present, sizes bounded,
`arm` ∈ {shell, world}, `entries` null for shell and `from ≤ to` for world), canonicalizes it,
and commits it as ONE object through the existing one-object writer, refactored out of
`demo_proof.go` into a shared helper with `writtenBy` as a parameter (`harness:floor` here,
`demo` unchanged for row 150). It prints the object ref and the committed entry index. Exit codes:
0 committed, 1 refused before sending, 2 commit refused, 3 uncertain (reconcile; never retried).
No new route, no new authority: the commit is the episode's own session.

**D3 — the harness hook.** In `run_cell`, after grading and before the row is appended, when the
run has a World client and the agent's episode token (the paired phase and `iterate` with World
arms), the harness writes the D1 payload to a temp file and calls the rig's
`bin/ailang-worldd outcome record --session private/<ep>.session`. The record is outside every
timed window (after `wall_ms` is taken). The row gains `outcome_record: {ref, entry}` or
`{error}`. A record failure never changes the task's class; the phase record counts
`outcome_record_errors`, and a FINAL treats any non-zero count as a harness defect to report.
Gated by a new config knob `outcome_capture` (default true; it enters the FINAL config digest,
so a FINAL preregisters it). Phase 1 runs no daemon and records nothing (stated, not hidden).

**D4 — no back-fill.** K1–K3 and every store captured before this lands stay as they are; the verb
refuses nothing about age, but row 114's protocol forbids seeding a declared cohort (§6 INVALID).

**Why not a package (S3)?** Nothing here touches `world/` or the kernel: it is an operator CLI verb
at the host boundary plus a harness change. The object is data committed through the existing
commit surface.

**Why not task context in MCP?** Labelling every invocation with its task would need the
coordinator and session to carry a task id from the client — a broader product change. Mark chose
the harness record; per-invocation labels stay a possible later row.

## §3 Acceptance criteria

| AC | Criterion |
|---|---|
| AC1 | `outcome record` commits one `world/task-outcome/v1` object; `object find world/task-outcome/v1` and `object get <ref> --payload` return it byte-equal to the canonical payload; the entry's `writtenBy` is `harness:floor`. |
| AC2 | Refusals before any send (exit 1): unknown field, missing required field, payload > 16 KiB, bad `arm`, shell row with entries, world row with `from > to` or missing entries, non-canonical numbers. Each has a test. |
| AC3 | No session → the daemon's 401 is surfaced (exit 2), nothing committed. Uncertain commit → exit 3. |
| AC4 | `demo` seed/proof behaviour unchanged (existing row-150 tests pass unmodified). |
| AC5 | Harness: a paired `iterate` (fake world or a real daemon in the test rig) appends one outcome per graded cell, both arms; World rows carry the task's `log_from..log_to`; a record failure leaves `class` unchanged and is counted. |
| AC6 | `why` and `provenance` keep working over a log holding outcome entries (`why head` when head is an outcome entry; `provenance --episode`). |
| AC7 | End-to-end on a K1 COPY is NOT allowed (D4). Instead, a scratch store: run a real paired `iterate` smoke cell with a known-failing task, then a walk on the served copy finds the outcome by `object find` and reads the rule and the entry range. |

## §4 Mutation controls (the judge runs them)

MU-1 drop `entries` for a world row → AC2 refusal test fails. MU-2 skip canonicalization → AC1
byte-equality fails. MU-3 swallow a record failure into the row's class → AC5 fails. MU-4 hardcode
`writtenBy` to `demo` → AC1 fails. MU-5 record before grading (outcome absent) → AC5 fails.

## §5 Risks

| risk | mitigation |
|---|---|
| A commit between tasks shifts the next task's `log_from` | `log_from` is read before each spawn (`run_task.py:282`), after the previous record |
| The record adds latency to the run | outside timed windows; one commit per cell (~ms locally) |
| Row 169 (REST commit accepts an unlinked `prevEntryHash`) | the shared writer always links to the current head entry (as `demoProofWriter` does); not widened here |
| A future FINAL's digest changes | intended; the knob is preregistered |

## §6 Still unknown

- Whether a fresh walk answers a Q1-type question from a captured store within 300 s. Only a
  real, harvested question from a store captured after this lands can show it (row 114).
