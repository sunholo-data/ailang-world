# w-exec-failure-carries-a-reason — a failed `workspace-exec` call records, returns and logs one bounded failure line

- **Status**: Planned (designed iteration 360, **r1**). Queue row 255, in D-WORLD-108's scope after
  row 250 (attended 2026-11-06: rows 250, then 255 → 257). Rows 256 (`w-daemon-startup-under-load`)
  and 257 (`w-exec-cache-seed-at-bind`) are separate rows. §9 says how this row's classes serve
  them, and designs neither.
- **Item**: Clause 5 (MET; operator and agent UX). It adds an optional field to the
  `world/effect-record/v1` wire, so its codec choice (§3.1) is the decision a quorum should check.
- **Designer lane**: `claude:claude-opus-5-5` via Agent opus.
- **Estimated**: ~0.5 day. A complete patch is
  `design_docs/planned/w-exec-failure-carries-a-reason.patch` (1184 lines, 15 files, sha256
  prefix `b7dcb7821c3b389a`). It applies with `git apply --check` at `77c9e3d4`, and its tests pass
  on a fresh `git archive` of `77c9e3d4` with the patch applied (V14, V15). Production code is
  +108/−28 lines in 8 files plus one new 189-line file, `host/broker/failure.go`, of which about a
  third is comments. There are four new test files (478 lines), and `why_test.go` changes by +44/−1.
  `docs/QUICKSTART.md` gains one paragraph (S7).
- **Measurement base**: `77c9e3d4` (= origin/dev), worktree `sprint/world-iter360-exec-reason`.
  `handlers_exec.go` sha256 prefix `d85457256d78add8`, `record.go` `57a0ec0e8b09c0ea`, `broker.go`
  `788f6eaef667182a` (V1).
- **Instruments**: `go1.26.6 darwin/arm64`, matching `go.mod`'s `go 1.26.6`. Pinned
  `AILANG_BIN=$HOME/.pinned-ailang/ailang` (`AILANG v0.41.0`, commit `24ee108`). TMPDIR was left at
  its default (`/var/folders/…/T/`). The rig was shared, with load averages of 2.2–9.0 at the timed
  runs. No Linux host was used, so every Linux statement is **UNMEASURED** (§6 R6). Scratch copies
  were made under `/tmp/iter360-designer/` (`applied/`, `head/`). The prototype was built in the
  worktree, written out as the patch, and then reverted (V14). This doc and the patch are the only
  files this commit adds. Raw command outputs are quoted in §10 rather than banked as files,
  because the designer's brief allows no other file.

---

## 1. What row 255 says, and what this session measured

| # | Claim | Verdict at `77c9e3d4` | Evidence |
|---|---|---|---|
| 1 | A failed effect record holds `"failed": true, "denial": "", "resultRef": ""` and no reason | **CONFIRMED.** The broker writes the failed record from the request, budget and refs alone (`broker.go:292–296`). The pinned golden is `wantFailed` (`broker_test.go:622`), and the wire has no field that could hold a reason (`record.go:32–43`) | V2, V3 |
| 2 | The agent's MCP result says only `status: failed` | **CONFIRMED.** `requestEffect` copies only `failed.RecordRef` (`effectful.go:191–192`). The body of a failed one-effect call is `{}` (`:245`), and `worldEffect` has only `id`, `status` and `record` (`:213`) | V2 |
| 3 | `why` ends at `e1 failed … no result` | **CONFIRMED.** `why.go:906` sets `failed`, and `:911` writes `no result` for a zero `resultRef`. There is nothing else to print | V2 |
| 4 | The daemon log has no line | **CONFIRMED.** A failed effect is a *successful* `Dispatch`: the coordinator commits the failed invocation and returns its output. `mcp.go:131–133` logs only when `Dispatch` returns an error, so this path never logs | V2 |
| 5 | `Execute` returns bare Go errors from 7 named checks (`tools.resolve`, `sandbox.Verify`, `verifyFileDigest`, `verifyExecSettings`, `projectDir`, the seed, the budget) | **TRUE BUT INCOMPLETE.** There are **14** error returns: 11 in `Execute` (`handlers_exec.go:253, 256, 277, 291, 294, 297, 301, 304, 315, 356, 370`) and 3 in `prepareCache` (`:527, 538, 551`). The row's list misses the effect/scope guard (`:253`, `:256`), the cleanup-incomplete return (`:356`) and the post-spawn start failure (`:370`), and counts `prepareCache`'s directory creation and seed copy as one | V7 |
| 6 | The reason is lost | **Lost at the coordinator, not the broker.** On the live path the broker keeps the handler's error as `EffectFailedError.cause` (`broker.go:146`, unwrapped at `:154`). `requestEffect` drops it (`effectful.go:191–192`). On replay the cause cannot exist: `invokeReplay` rebuilds the error from the record alone (`broker.go:442–443`). So the reason must live **in the record**, or replay cannot reproduce the committed output (claim 7) | V2, V5 |
| 7 | (design constraint) Replay must reproduce the committed output byte for byte | **CONFIRMED.** `Coordinator.Replay` re-composes the output from the recorded effect records and refuses unless it equals the committed bytes (`effectful.go:485`). Failed and denied records are replayed by test today (`effectful_test.go:521`) | V5 |
| 8 | The 2026-11-03 cause is unattributed | **Not re-measured.** This session did not read the home store. The design gives each of the row's candidates its own class (§3.2), so the next instance names itself | — |

**Premises refined, not reversed.**

- The row's example class `budget-spent` collides with the broker's capability-budget denial
  label `denied:budget` (`decide.go:11`), which is a different thing: a *denial*, with nothing
  dispatched and the budget untouched. The exec class is named `time-budget-spent` (§3.4).
- The row says "keep paths that are already in the profile". The profile's `caches.*.seed` and
  `read_roots` are operator host paths that no agent-visible result contains. So this design uses
  a stricter set: it names a file only as a **successful** exec result already does (V9), or
  relative to the srt archive, or by its base name (§3.3).

## 2. The mechanism being changed

The failure path at `77c9e3d4`, hop by hop (V2):

1. **Handler.** `ExecHandler.Execute` (`handlers_exec.go:251`) returns `(nil, err)` from any of the
   14 sites (V7). A *refusal* is different: `execRefusal` (`:96–97`) returns
   `{"ok":false,"refused":…}` with a **nil** error.
2. **Broker.** `Session.invoke` calls the handler (`broker.go:261`). On a non-nil error it writes
   `EffectRecord{Allowed: true, Failed: true, …}` (`:292–296`) through `putRecord`, appends the
   `failed` outcome, and returns `&EffectFailedError{…, cause: err}` (`:307`).
3. **Coordinator.** `requestEffect` (`effectful.go:177`) maps `*EffectFailedError` to
   `EffectResult{Status: "failed", Record: …}` and returns a nil error (`:191–192`), so the cause
   ends here. `composeEffectOutput` builds `{}` plus the `world` block
   `{"effects":[{"id","status","record"}],"plan"}` (`:213`, `:245`). The call then commits
   normally.
4. **Projection.** `mcpAdapter.Invoke` returns `result.OutputBytes` as the tool result
   (`mcp.go:135`). The A2A path returns `result.Output` as the task artifact
   (`projection.go:363`, `:380`). Neither logs, because `Dispatch` succeeded.
5. **Readers.** `why` decodes the record (`why.go:896`) and prints `failed … no result`
   (`:906`, `:911`). `ailang-worldd call` prints `effect e1 failed <record>` (`mcpclient.go:642`).
   `invokeReplay` rebuilds `EffectFailedError` from the record (`broker.go:442`).

The record codec (`record.go`) is `world/effect-record/v1`: one fixed-order `recordWire` struct,
encoded with `json.Marshal`, and decoded with `DisallowUnknownFields` (`:60`). The record's
object hash is the sha256 of those bytes (`brokerObject`, `:115–122`).

## 3. Design

### 3.1 Where the line lives: a new optional `failure` field on the v1 record

| Option | Verdict | Why (measured) |
|---|---|---|
| **(a) A new field `failure`, `omitempty`, last in `recordWire`** | **CHOSEN** | `encoding/json` omits an empty string, so every record without a line encodes to the bytes it has today. That covers every success, every denial, and every failure of a handler that is not `ExecHandler`. Both existing goldens stay green (`TestRecordGoldenBytes`, V15), and `MUT-NO-OMITEMPTY` reddens them (§5). A stored record without the field decodes to `Failure: ""`. No stored object's hash changes, because content-addressed objects are never re-encoded. Only failed `Workspace.Exec` records written from now on carry it |
| (b) Reuse `denial` | **Rejected** | It breaks the record law. `RecordConsistent` requires `Denial == ""` on every allowed record, failed ones included (`record.go:100–101`), mirroring the sketch's ensures clause (`effectbroker.ail:172`). Replay refuses an allowed record with a non-empty denial (`broker.go:433`). It also blurs a capability denial with a handler failure (§3.4) |
| (c) A failure object referenced by `resultRef` | **Rejected** | It also breaks the law: a failed record must have a zero `resultRef` (`record.go:101`, `effectbroker.ail:172`). Readers treat a non-zero `resultRef` as a success result: `why` prints `result N B` and reads its worktree delta (`why.go:915–923`), and the episode report reads it as a handler answer (`episode_report.go:533`) |
| (d) A `world/effect-record/v2` | **Rejected** | It has the same downgrade property as (a) (R1) and costs much more. Ten non-test call sites decode records (V10), and several check `SemanticID == EffectRecordV1` (`why.go:525`, `broker.go:418`, `workbench_packet_evidence.go`, `episode_report.go`), so a v2 would need a dual-version reader in each. (a) keeps one codec, and old records stay valid v1 |

**What (a) does to existing state and tests** (V3, V15, V16):

- **Old stored records** (any kind) decode unchanged: the field is absent, so `Failure` is `""`.
- **Old committed outputs** replay byte-identically. Their records have no line, so the
  recomposed `world` block has no `reason` key (`omitempty`), which is exactly what was committed.
  TC's `unclassified` row pins this, and `MUT-WORLD-NO-OMITEMPTY` reddens it (§5).
- **`TestRecordGoldenBytes`** (`broker_test.go:598`): green and unchanged. It pins that a success
  and an unclassified failure keep their bytes.
- **`TestSketchRows`** (`decide_test.go`, rows 175–207) and **`TestRecordConsistentAllSketchArms`**:
  green and unchanged, because `RecordConsistent` is unchanged.
- **The sketch** `design_docs/sketches/effectbroker.ail` is **not** changed. Its `EffectRecord`
  type and `recordConsistent` law state the accounting (budget, allowed/failed/denied, refs), and
  the line does not touch the accounting. No test pins the Go struct's field list against the
  sketch type. The only sketch drift test is the row table in `decide_test.go` (V3). The new rule
  is a separate Go predicate, `FailureConsistent`, next to `RecordConsistent`, not inside it:
  `Failure == "" || (Allowed && Failed && ValidFailureLine(Failure))`. `putRecord` refuses a record
  that breaks it, and `DecodeRecord` refuses to decode one. **No `.ail` file is touched.** Whether
  the sketch should gain the field for documentation parity is open question Q3.
- **The one cost: downgrade.** A binary from before this change cannot decode a record that
  carries `failure`. `DecodeRecord` at `77c9e3d4` answers `json: unknown field "failure"`, and the
  same record without the field decodes (V16). This is R1.

### 3.2 The classes (a closed set of 10)

Every error `Execute` returns is now a `*HandlerFailure{Class, Detail, Err}` (V17: 0 unclassified
returns, against 10 classified). Its `Error()` and `Unwrap()` are the cause's, so every existing
caller and test that matches on the cause is unchanged (`TestExecHandlerReverifiesTheStackBeforeEachCall`
stays green, and `MUT-HF-ERROR-IS-LINE` reddens it). Only `Line()` is recorded.

| Class | Execute site at `77c9e3d4` → patched | Detail (built only from the words shown) | Example line (each one is a TX row) |
|---|---|---|---|
| `exec-misrouted` | `:253`, `:256` → one guard | `<effect>@<scope>` of the request (the plan's manifest triple) | `exec-misrouted: Workspace.Exec@/` |
| `exec-tools` | `:277` (`tools.resolve` error, not a drift) | fixed text | `exec-tools: the per-call re-check of the exec tools failed` |
| `exec-stack-drift` | `:291` (`sandbox.Verify`), `:294` (`verifyFileDigest`), and a drift found inside `:277` | `<named path>: <what>`. The path is named by §3.3. `<what>` is one of verifyArchive's own phrases (`changed since it was archived`, `missing from the archive`, `not in the archived manifest`), or `digest changed`, or `unreadable` | `exec-stack-drift: srt archive @anthropic-ai/sandbox-runtime/dist/cli.js: changed since it was archived`; `exec-stack-drift: <node>: digest changed` |
| `exec-settings-drift` | `:297` | `<settings file base name>: <what>`, where `<what>` is `not a regular 0600 file`, `bytes differ from World's rendering` or `unreadable` | `exec-settings-drift: ep1.proj.srt.json: bytes differ from World's rendering` |
| `project-root` | `:301` | `root "<profile root>" is not a directory inside the worktree` | `project-root: root "missing" is not a directory inside the worktree` |
| `exec-cache` | `:527`, `:538` (in `prepareCache`) | `<dir or cache name>: cannot be created` | `exec-cache: home: cannot be created` |
| `seed-cache` | `:551` (in `prepareCache`) | `<cache name>: <how>`, from the run error's **type** only: `timed out after <d>`, `the copy exited <n>`, `the copy was not reaped`, `the copy could not run` | `seed-cache: GOCACHE: the copy exited 1`; `seed-cache: GOCACHE: timed out after 1ms` |
| `time-budget-spent` | `:315` | `the call's time was spent before "<command id>" could start` | `time-budget-spent: the call's time was spent before "test" could start` |
| `cleanup-incomplete` | `:356` | fixed text | `cleanup-incomplete: the command's processes were not all reaped` (TK only: §6 R3) |
| `sandbox-start` | `:370` (post-spawn error that is not an exit, timeout or overflow) | fixed text | `sandbox-start: the sandbox could not be started` |

`prepareCache` classifies its own three sites. `execFailure`, the one classification point in
`Execute`, returns an existing `*HandlerFailure` unchanged. It then overrides the step's class
when the cause is a typed `*ExecStackDriftError` or `*ExecSettingsError`, and otherwise uses the
step's class and fixed detail.

### 3.3 Bounding and secrecy

**The rule (pure, `FailureLine(class, detail)`):**

- An unknown class gives `""`, so nothing unclassified is ever recorded.
- Otherwise the line is `class`, or `class + ": " + detail`.
- In the detail, invalid UTF-8 becomes `?`. Every control character (`unicode.IsControl`) and
  format character (`unicode.Cf`, for example the bidi override U+202E) becomes a space. Runs of
  whitespace, U+2028 included, collapse to one space.
- The line is at most **`FailureLineMax` = 200 bytes**, cut at a rune boundary with a trailing `…`.
  200 sits between `call`'s 160-byte elision (`mcpclient.go:595`) and the episode report's
  240-byte reason cap (`episode_report.go:155`). It fits the longest class (19 bytes), a typical
  absolute node path (about 60) and the longest fixed phrase (30).
- `ValidFailureLine(s)` holds exactly when `s` is non-empty and a **fixed point** of the rule
  (`FailureLine(Cut(s, ": ")) == s`). The decoder and `putRecord` use it, so a stored line is
  always one the rule could have produced.

**What a line may contain.** Only:

- World's own fixed words;
- the class constants;
- the request's effect and scope, which are the plan's manifest-validated triple;
- the profile's `root` and cache names;
- the command id, which has matched the id grammar and a profile entry before any failure site;
- file names produced by `execDisclosed.name`.

`execDisclosed.name(p)` returns `p` verbatim only when `p` is one of the four paths a successful
exec result already returns to the agent: `profile.argv0` and `argv0_target`, and `sandbox.node`
and `node_target` (`handlers_exec.go:231–241`, V9). A path under the srt archive becomes
`srt archive <rel>`. Any other path becomes its base name. The state dir, the operator `$HOME`,
seed paths and `read_roots` never appear. TX asserts that no line contains the fixture's state dir
or `$HOME`. **The Go error's text never enters a line.** Typed errors are read by field (`Path`,
`Why`), and `Why` passes only when it equals one of the three fixed phrases. The cause's text,
which can hold host paths, OS errors and child output (`exec_sandbox.go:107, 248`;
`exec_settings.go:147`; `handlers.go:74`, V8), stays in `Err`, which is never recorded.

### 3.4 Failure is not refusal, and neither is a denial

| | Who decides | Record | Budget | Agent sees | This row |
|---|---|---|---|---|---|
| **Denial** | the broker's capability decision (`decide`) | `allowed:false, denial:"denied:…"` (`broker.go:195–199`) | not debited | `status: "denied"` | unchanged |
| **Refusal** | the handler answers `{"ok":false,"refused":…}` with a nil error (`handlers_exec.go:96–97`) | `allowed:true, failed:false, resultRef:<the answer>` | debited | the handler's `refused` text, as output fields | unchanged; the profile, grammar and node-moved paths stay refusals |
| **Failure** | the handler returns an error | `allowed:true, failed:true, resultRef:""` (`broker.go:292–296`) | debited (CF-J-2 law) | `status: "failed"` (today with no reason) | **gains `failure` and `reason`** |

No call changes category. A call that was refused is still refused, and a call that failed still
fails. The classes describe *why World could not run* the command. A refusal is *World would not
run it*. That is why `time-budget-spent` is not named `budget-spent`: a reader who sees `budget`
should not have to ask which of the two budgets ran out.

### 3.5 The carriers (one line, four surfaces)

- **Broker** (`broker.go`): if the handler's error is (or wraps) a `*HandlerFailure`, then
  `reason := classified.Line()`. It goes into `EffectRecord.Failure` **and** into a new exported
  field, `EffectFailedError.Reason`. `invokeReplay` sets `Reason: rec.Failure`, so live and replay
  carry the same string (§4 claim D).
- **Coordinator** (`effectful.go`): `requestEffect` copies `failed.Reason` into a new
  `EffectResult.Reason` tagged `json:"-"`, so no finish phase's input bytes change.
  `composeEffectOutput` puts it in `worldEffect.Reason`, tagged `json:"reason,omitempty"`. The
  agent's MCP (and A2A) result is
  `{"world":{"effects":[{"id":"e1","reason":"<line>","record":"sha256:…","status":"failed"}],"plan":"sha256:…"}}`.
  Every `workspace-exec` plan has `"finish":false` (V6), so the body stays `{}`.
- **Log** (`projection.go`, new `logEffectFailures`): after a successful `Dispatch` on MCP
  (`mcp.go`) and A2A (`projection.go`), one line per effect that has a reason:
  `ailang-worldd: <surface> effect failed: <invocation id> <effect id> <record>: <line>`. It
  writes to the same `errorLog` as `logRefusal`. A `Reconciled` result (a resent A2A task, read back
  from the journal) logs nothing, so each call that ran logs exactly once.
- **`why`** (`why.go`): a record with a line prints `no result: <line>` where it printed
  `no result`. A record without one prints exactly as before.
- **`ailang-worldd call`** (`mcpclient.go`): `effect e1 failed sha256:… — <line>`.
- **Docs (S7)**: one `docs/QUICKSTART.md` paragraph under `workspace-exec` names the field, the
  classes, the three places the line appears and the refusal distinction.

### 3.6 Pure core, and why this is not a package

`FailureLine`, `ValidFailureLine`, `FailureConsistent`, `execFailure`, `execDisclosed.name`,
`driftWhat`, `settingsWhat` and `seedWhat` are pure functions of their arguments: no I/O, no
clock. TL, TG and TK test them with no process run. This is Go host code, and **no `.ail` is
touched**. S1's Z3 contracts bind `world/` and `.ail` code, so they do not apply here. The
change is in the effect handler and its broker, which are the host boundary (S2), not the
`world/` kernel. A package cannot hold it, because the line is made where the Go error is made
(S3).

## 4. The by-construction argument

**Claim A (no existing bytes move).** A record's bytes differ from what `77c9e3d4` would write
only if the record is failed **and** its handler returned a `*HandlerFailure`. *Proof.* `Failure`
is set in exactly one place, from `classified.Line()` under `errors.As(err, &classified)`
(`broker.go:298`, patched). The field is `omitempty`. `HandlerFailure` is constructed in no
non-test file except `failure.go` and `handlers_exec.go` (V17: 0 other files). So only
`Workspace.Exec` failures change. The same holds for outputs: `reason` is `omitempty`, and
`EffectResult.Reason` is `json:"-"`. ∎ Pinned by `TestRecordGoldenBytes`, unchanged and green,
and by TC's `unclassified` row (no `reason` key).

**Claim B (no Go error text reaches a line).** *Proof.* `Line()` reads only `Class` and `Detail`.
Every `Detail` in `failure.go` and `handlers_exec.go` is a concatenation of string literals, the
profile's `root` and cache names, the request's effect and scope, the validated command id, and
`execDisclosed.name(…)`, `driftWhat(…)`, `settingsWhat(…)` or `seedWhat(…)`. The last three return
only literals, a whitelisted literal, or `Duration.String()`/`strconv.Itoa` of a typed field. No
`Detail` is built from `err.Error()`. ∎ TK and TX kill `MUT-ERR-TEXT`, `MUT-RAW-PATH` and
`MUT-DRIFT-WHY-RAW`.

**Claim C (every stored line is bounded and one line).** *Proof.* `putRecord` refuses, and
`DecodeRecord` rejects, any record whose `failure` is not a fixed point of `FailureLine`, and every
such fixed point is at most 200 bytes with no control or format character. ∎ TG covers the
decode side. The put side is unreachable by construction (`MUT-PUT-NO-CHECK` survives, §5),
because a `Line()` is always a fixed point (TL asserts that for every row it produces).

**Claim D (live and replay agree).** *Proof.* The live `EffectFailedError.Reason` is the same
`reason` variable written to the record. The replayed one is `rec.Failure`, read from that
record. The coordinator composes the output from `Reason` in both modes, so the replayed output
equals the committed one. ∎ TC (`Replay` equals the committed bytes) kills
`MUT-REPLAY-DROPS-REASON` and `MUT-RECORD-DROPS-LINE`.

**Claim E (logged once per call).** *Proof.* The log line is written only after a `Dispatch`
that returned a result that is not `Reconciled`, once per effect with a non-empty `reason`, and a
`workspace-exec` plan has exactly one effect. ∎ TP sends the same A2A task twice, which runs once,
and expects exactly one line.

## 5. Test plan and mutation table

| Id | Test (file) | What it asserts |
|---|---|---|
| **TL** | `TestFailureLineBoundsAndSanitises` (`host/broker/failure_test.go`) | the closed set has 10 classes; the rule's output on clean, class-only, control, invalid-UTF-8, bidi/U+2028 and unknown-class inputs; a 300-rune detail is cut to ≤ 200 bytes, ends in `…`, is valid UTF-8 and is a fixed point; 5 malformed lines are not valid |
| **TG** | `TestRecordFailureGoldenBytes` (same file) | the exact bytes of a failed record with a line (`…,"resultRef":"","failure":"seed-cache: GOCACHE: timed out after 9s"}`) and their round trip; decoding refuses a line on a success, on a denial, a line that is not a fixed point, and an unknown class |
| **TK** | `TestExecFailureClassifies` (same file) | `execFailure` with no process run: a typed drift or settings error overrides the step's class; archive paths are relative, the node path is verbatim, other paths become base names; a drift found through `resolve` is a drift; an existing `*HandlerFailure` passes unchanged; `cleanup-incomplete`; an untyped cause gets the step's fixed detail; a planted `/Users/op/.secret` never appears; every failure unwraps to its cause; `seedWhat` on a timeout, a joined timeout and cleanup error, cleanup-incomplete, and a start error |
| **TX** | `TestExecFailureClassesReachTheRecord` (same file) | **one Execute-driven row per reachable class, 11 rows**: misrouted, tools (closed tool state), stack drift (archive, node), settings drift, project root, exec cache, seed (exited, timed out), time budget, sandbox start. Each runs through a **live broker session over a real store** (the CF-J-2 rig) and asserts the stored record's `failure` and `EffectFailedError.Reason` both equal the expected line; the record is `RecordConsistent`; the line names neither the state dir nor `$HOME`; the cause's own text is kept |
| **TC** | `TestEffectFailureReasonReachesTheWorldBlockAndReplays` (`host/coordinator/effect_failure_reason_test.go`) | a `*HandlerFailure` from the handler reaches `world.effects[0].reason` and the record, the cause's text does not reach the output, and `Replay` reproduces the committed bytes. Row `unclassified`: a plain error leaves **no** `reason` key and an empty `failure`, and also replays |
| **TP** | `TestMCPEffectFailureLineReturnedAndLoggedOnce` (`host/projection/effect_failure_log_test.go`) | a real MCP `tools/call` over the projection: the tool result carries the line and no cause text, and the error log holds exactly one line matching `^ailang-worldd: mcp effect failed: mcp:ep-a:<64 hex> e1 sha256:<64 hex>: <line>\n$`. Then the same A2A task sent twice: both answers carry the line, and the log holds exactly one `a2a` line |
| **TW** | `TestWhyShowsTheFailureLine` (`cmd/ailang-worldd/why_test.go`) | with a committed failed `Workspace.Exec` invocation whose record carries a line, `why` verifies all 7 links and prints `no result: <line>` on the effect link; with a failed record without a line, it prints as before. The `walkRig` helper gains a `failure` argument; existing callers pass `nil` |
| **TR** | `TestCallPrintsTheFailureLine` (`cmd/ailang-worldd/call_failure_line_test.go`) | `printCallResult` prints `effect e1 failed sha256:1111 — <line>`, and without a reason prints as before |

The row asks for "one test per failure class in `Execute`, showing that the record, the MCP result
and `why` each carry it". The design meets it in two layers. **TX runs once per class** at the
place the class is decided, and asserts the record and the `EffectFailedError.Reason` that the
coordinator copies into the MCP result. **TC, TP and TW** then prove each downstream surface
carries the line it is given. Past the broker the line is an opaque string, so repeating those
three per class would test the same carrier ten times. `MUT-COORD-DROPS-REASON`,
`MUT-WHY-NO-LINE` and `MUT-LOG-MCP-NONE` show each carrier is pinned.

**Mutation run** (V12; the harness is quoted in V12). Each mutant is one exact-string edit of the
prototype. The 13 tests run once (`-count=1`) across the four packages, and then the file is
restored, with the tracked diff compared to its pre-run copy (`cmp`: identical). The columns list
the tests that went red. `Gold` = `TestRecordGoldenBytes` and `Reverify` =
`TestExecHandlerReverifiesTheStackBeforeEachCall`, both existing. `CFJ2`, `Sketch` and `Replay`
(`TestEffectReplay`) were in the set and stayed green under every mutant.

| Mutant | Edit | Red | Expected → observed |
|---|---|---|---|
| CONTROL | none | — (13 of 13 PASS) | green → green |
| **MUT-HEAD-SHAPE** (positive control) | the record drops the line **and** the coordinator drops the reason: today's behaviour inside the new structure | **TX, TC, TP** | red → red. TX's text: `record = broker.EffectRecord{… Failed:true, Denial:"" …}` with no failure, 11 of 11 rows; TC's output `{"world":{"effects":[{"id":"e1","record":"sha256:…","status":"failed"}],…}` |
| MUT-NO-OMITEMPTY | `json:"failure,omitempty"` → `json:"failure"` | **Gold** | red → red: the two existing goldens move, so every existing record's bytes are pinned |
| MUT-RECORD-DROPS-LINE | the failed record omits `Failure: reason` | TX, TC | red → red (TC through replay divergence) |
| MUT-REPLAY-DROPS-REASON | `invokeReplay` omits `Reason: rec.Failure` | TC | red → red |
| MUT-COORD-DROPS-REASON | `requestEffect` does not copy `failed.Reason` | TC, TP | red → red |
| MUT-WORLD-NO-OMITEMPTY | `json:"reason,omitempty"` → `json:"reason"` | TC | red → red (`unclassified` grows a key, so old outputs would no longer replay byte-for-byte) |
| MUT-ERR-TEXT | the untyped fallback's detail is `err.Error()` | TK, TX | red → red |
| MUT-RAW-PATH | `execDisclosed.name` returns every path verbatim | TK, TX | red → red |
| MUT-DRIFT-WHY-RAW | `driftWhat` returns `Why` verbatim | TK, TX | red → red |
| MUT-NO-SANITISE | no character mapping | TL | red → red (first anchor did not compile; re-anchored, V12) |
| MUT-CONTROL-ONLY | maps control characters but not format characters | TL | red → red (the U+202E row) |
| MUT-NO-BOUND | no length cut | TL | red → red |
| MUT-CLASS-OPEN | an unknown class is accepted | TL, TG | red → red |
| MUT-DECODE-NO-CHECK | `DecodeRecord` skips `FailureConsistent` | TG | red → red |
| MUT-FAILURE-ANY-ARM | a line is allowed on any arm | TG | red → red |
| MUT-HF-ERROR-IS-LINE | `HandlerFailure.Error()` returns the line, not the cause | **Reverify**, TX | red → red: existing callers that read the cause are pinned |
| MUT-SEED-UNCLASSIFIED | the seed copy returns a plain error | TX | red → red |
| MUT-SEED-TIMEOUT-AS-EXIT | `seedWhat` skips the timeout case | TK, TX | red → red |
| MUT-LOG-MCP-NONE | no log call on MCP | TP | red → red |
| MUT-LOG-A2A-NONE | no log call on A2A | TP | red → red |
| MUT-LOG-RECONCILED | a reconciled result also logs | TP | red → red (two `a2a` lines) |
| MUT-WHY-NO-LINE | `why` ignores `Failure` | TW | red → red |
| MUT-CALL-NO-LINE | `call` ignores `reason` | TR | red → red |
| MUT-PUT-NO-CHECK | `putRecord` skips `FailureConsistent` | — | **SURVIVES.** Unreachable by construction: the only writer passes `Line()`, which is always `""` or a fixed point (TL). The check stays as a guard against a future writer; decoding is pinned by TG |
| MUT-CLEANUP-UNWIRED | `:356`'s return goes back to the bare `runErr` | — | **SURVIVES.** No seam leaves a child unreaped inside `Execute` without a process-kill fault injector, which this item does not add. The class's text is pinned by TK, and the wiring is one line in the same pattern as the 9 sites TX covers (R3) |

## 6. Residuals and risks

- **R1, downgrade (stated, not fixed).** After this lands, a daemon or CLI from before it cannot
  decode a failed `Workspace.Exec` record that carries a line (V16). In that store, an older
  `why` marks that link ✗, an older episode report stops at that entry, and an older replay of that
  invocation is a replay gap. Records without a line are unaffected. World ships forward-only
  binaries against one store, and no item downgrades the daemon over a live store. A v2 record
  (option (d)) has the same property, because an old binary cannot read v2. Q1 asks the quorum to
  accept this.
- **R2, the operator log carries the bounded line, not the full Go error.** The row says "the same
  line … logged". The full cause stays in `HandlerFailure.Err` and is dropped at `requestEffect`,
  as at HEAD. When a line is not enough (for example `exec-tools` or `sandbox-start`), the operator
  still has to reproduce the failure. Plumbing the cause to the projection's operator-only log is
  Q2.
- **R3, `cleanup-incomplete` is not driven through `Execute`** (`MUT-CLEANUP-UNWIRED` survives).
  TK pins its class and text. A seam would need a child that survives the group kill past the
  cleanup deadline.
- **R4, base-name disclosure.** A drifted file outside the four disclosed paths and the archive
  (for example `libnode.127.dylib`, or a tool's identity target) is shown by base name. The name of
  an operator file reaches the agent. The full path never does.
- **R5, the other two exec handlers.** `ExecUnconfiguredHandler` and `ExecRefusalHandler` still
  return bare errors for a misrouted effect or scope (`handlers_exec.go:42–45`, `:88–91`). Both
  are unreachable through the registry: the effect is the registry key, and the scope is the
  manifest's. Both answer every real call with a refusal.
- **R6, Linux (UNMEASURED).** TX's `seed-cache exited` row expects GNU `cp -a --reflink=auto` on a
  missing seed to exit 1, and its `sandbox-start` row expects `EACCES` from a 0644 node. Both run
  in CI's plain `go test ./...` leg on Linux, which is the acceptance gate (§7).
- **R7, `why --json`.** The line is in the effect link's detail text, not a structured field.
  Adding one is not asked for.

## 7. Milestones

**M1 (single, ~0.5 day).** Apply the patch, then prove each step. Leave TMPDIR at its default: a
TMPDIR under `~/dev` makes an unrelated test red.

```bash
cd <worktree>
git apply --check design_docs/planned/w-exec-failure-carries-a-reason.patch && git apply design_docs/planned/w-exec-failure-carries-a-reason.patch
# expect: rc 0, no output
go build ./... && go vet ./host/broker/ ./host/coordinator/ ./host/projection/ ./cmd/ailang-worldd/ && test -z "$(gofmt -l host cmd)" && echo OK
# expect: OK
go test -count=1 -v -run '^(TestFailureLineBoundsAndSanitises|TestRecordFailureGoldenBytes|TestExecFailureClassifies|TestExecFailureClassesReachTheRecord|TestRecordGoldenBytes|TestExecHandlerReverifiesTheStackBeforeEachCall|TestCFJ2HandlerErrorKeepsDebitAndWritesOneFailureRecord|TestSketchRows|TestEffectFailureReasonReachesTheWorldBlockAndReplays|TestEffectReplay|TestMCPEffectFailureLineReturnedAndLoggedOnce|TestWhyShowsTheFailureLine|TestCallPrintsTheFailureLine)$' \
  ./host/broker/ ./host/coordinator/ ./host/projection/ ./cmd/ailang-worldd/ | grep -E '^(--- |ok|FAIL)'
# expect: 13 lines "--- PASS: <name>", 4 "ok", no FAIL (V15)
go test -count=3 -run '^(TestFailureLineBoundsAndSanitises|TestRecordFailureGoldenBytes|TestExecFailureClassifies|TestExecFailureClassesReachTheRecord|TestRecordGoldenBytes|TestEffectFailureReasonReachesTheWorldBlockAndReplays|TestEffectReplay|TestMCPEffectFailureLineReturnedAndLoggedOnce|TestWhyShowsTheFailureLine|TestCallPrintsTheFailureLine)$' \
  ./host/broker/ ./host/coordinator/ ./host/projection/ ./cmd/ailang-worldd/
# expect: 4 "ok" (V13)
AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_ail.sh | tail -1
# expect: ✓ verify gate PASSED: 16 required identities verified, 40 named tests pass, 503 se-tools named tests pass
VERIFY_GO_LEGS=plain AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_go.sh | grep 'go gate PASSED'
# expect: ✓ go gate PASSED (plain): …
```

Then reproduce at least these kill patterns from §5, each red on its named test:
`MUT-HEAD-SHAPE` (TX, TC, TP), `MUT-NO-OMITEMPTY` (Gold), `MUT-REPLAY-DROPS-REASON` (TC),
`MUT-ERR-TEXT` (TK, TX), `MUT-RAW-PATH` (TK, TX), `MUT-LOG-RECONCILED` (TP) and
`MUT-HF-ERROR-IS-LINE` (Reverify). Record the two survivors with their reasons.

**Acceptance gate: CI on the PR.** The plain `go build + test gate` step (`verify_go.sh`) runs TX
on Linux (R6). The item is not accepted until that step is green on the PR head. A local green is
not a substitute.

## 8. Files

| File | Change |
|---|---|
| `host/broker/failure.go` | **new**, 189 lines: `FailureLineMax`, the 10 class constants and their closed set, `HandlerFailure` (`Error`/`Unwrap` = the cause's, `Line`), `FailureLine`, `ValidFailureLine`, `execDisclosed` and `name`, `execFailure`, `driftWhat`, `settingsWhat`, `seedWhat` |
| `host/broker/record.go` | +21/−4: `EffectRecord.Failure`, `recordWire.Failure` (`omitempty`, last), the codec's two positional literals, `DecodeRecord`'s check, new `FailureConsistent` |
| `host/broker/broker.go` | +13/−5: `EffectFailedError.Reason`; the failed record and error carry `classified.Line()`; replay sets `Reason: rec.Failure`; `putRecord` checks `FailureConsistent` |
| `host/broker/handlers_exec.go` | +24/−17: the 11 `Execute` returns go through `execFailure` (the effect and scope guards merge into one); `prepareCache`'s 3 returns build `*HandlerFailure` |
| `host/coordinator/effectful.go` | +9/−2: `EffectResult.Reason` (`json:"-"`), `worldEffect.Reason` (`omitempty`), the copy in `requestEffect` and `composeEffectOutput` |
| `host/projection/projection.go`, `mcp.go` | +32, +1: `logEffectFailures`, called after a successful `Dispatch` on A2A and MCP |
| `cmd/ailang-worldd/why.go`, `mcpclient.go` | +4, +4: `no result: <line>`; `call`'s `— <line>` |
| `docs/QUICKSTART.md` | +11: the "When a call fails" paragraph (S7) |
| `host/broker/failure_test.go` | new, 278 lines: TL, TG, TK, TX and the `expiredCtx` wrapper |
| `host/coordinator/effect_failure_reason_test.go` | new, 62 lines: TC |
| `host/projection/effect_failure_log_test.go` | new, 106 lines: TP |
| `cmd/ailang-worldd/call_failure_line_test.go` | new, 32 lines: TR |
| `cmd/ailang-worldd/why_test.go` | +44/−1: `invokeEffect` (a `failure` argument; `invokeEffectResult` wraps it) and TW |

## 9. What this item is NOT doing

- **Not row 256.** No change to the startup probe, its budget, the registry bootstrap, or how a
  left-out profile is reported. Row 256's "refuse with that name" for a left-out profile would be
  a **refusal** (§3.4), not one of these failures. What this row gives it: at call time, a call
  starved by load reads `time-budget-spent` or `seed-cache: …: timed out after …`. Neither is
  a stack or settings drift, so an operator can tell load apart from a confinement change.
- **Not row 257.** The seed is still cloned on the first call, inside its budget. Row 257 can
  report a bind-time seed failure in these exact words by reusing `FailSeedCache` and `seedWhat`.
  The per-call check it keeps would then report `exec-cache` or `seed-cache` lines.
- No change to which calls fail and which are refused, to budgets, or to the CF-J-2 debit law.
- No change to the sketch `effectbroker.ail`, to any `.ail`, to `RecordConsistent`, or to
  the record's semantic id (Q3).
- No episode-report change. It still counts `failed` effects without their class. Row 257's test
  for a visible bind-time failure is that row's.
- No reason for handlers other than `ExecHandler`. Their failed records keep today's bytes
  (Claim A). Other handlers can opt in later by returning a `*HandlerFailure` with their own
  classes added to the closed set.
- No retry, no backoff, no new seam in production code. The tests use the existing fixture, the
  existing `ExecToolState.Close`, and a test-only `expiredCtx` wrapper.
- No change to `tools/launchd/*`, skills, `~/.ailang/state/*` or `HOST_CHANGELOG.md` (that belongs
  to the release that ships this).

**Open questions for the quorum.**

- **Q1.** Accept R1, the one-way decode of new failed exec records by older binaries? (Default:
  yes, as for any additive record field; (d) has the same property.)
- **Q2.** Should the operator-only log also carry the full cause, `%q`-quoted like `logRefusal`?
  It would need `EffectFailedError`'s cause carried past `requestEffect`. (Default: no, the row
  says "the same line".)
- **Q3.** Should `effectbroker.ail`'s `EffectRecord` gain `failure: string` for documentation
  parity (S4)? That is an `.ail` change with sketch-test churn and no change to the law. (Default:
  no.)

## 10. Verification Log

Commands were run from `/Users/voightkampff/dev/sunholo-data/.wt-world-iter360` unless they say
otherwise. A bare line number refers to `77c9e3d4`. The base-commit greps ran under `bash`,
because zsh treats `$B:h…` as a history modifier, and this session lost one run to that before
switching.

| # | Claim | Command | Observed |
|---|---|---|---|
| V1 | Base and instruments | `git rev-parse --short HEAD`; `shasum -a 256 host/broker/{handlers_exec,record,broker}.go \| cut -c1-16`; `go version`; `grep '^go ' go.mod`; `echo $TMPDIR`; `uptime`; `$HOME/.pinned-ailang/ailang --version` | `77c9e3d4`; `d85457256d78add8`, `57a0ec0e8b09c0ea`, `788f6eaef667182a`; `go1.26.6 darwin/arm64` = `go 1.26.6`; `/var/folders/kr/…/T/`; loads 2.23–2.41 at start; `AILANG v0.41.0`, `Commit: 24ee108` |
| V2 | The failure path, hop by hop | `bash /tmp/iter360-designer/greps.sh`: `git show 77c9e3d4:<file> \| grep -n …` over `broker.go`, `effectful.go`, `mcp.go`, `projection.go`, `why.go`, `mcpclient.go` | `broker.go:146 cause error`, `:154 return e.cause`, `:261 result, err = handler.Execute(…)`, `:295 Allowed: true, Failed: true, RequestRef: requestRef,`, `:307 return nil, ref, &EffectFailedError{`, `:442–443 … Effect: rec.Effect, Scope: rec.Scope, RecordRef: recordRef,`; `effectful.go:191 case errors.As(err, &failed):`, `:192 res.Status, ref = StatusFailed, failed.RecordRef`, `:213 worldEffect{ID: r.ID, Status: r.Status, Record: r.Record}`, `` :245 body = []byte(`{}`) // denied or failed … ``, `:308 requestEffect(hctx, …)`; `mcp.go:130 coord.Dispatch`, `:132 a.h.logRefusal("mcp", "tools/call", …, err)` (inside `if err != nil`), `:135 return protocol.InvocationResult{Value: result.OutputBytes}, nil`; `projection.go:167 errorLog io.Writer`, `:363` A2A `Dispatch`, `:510 func (h *Handler) logRefusal`; `why.go:528`, `:896` `DecodeRecord`, `:906 verdict = "failed"`, `:911 resText := "no result"`; `mcpclient.go:642 "  effect     %s %s %s\n"` |
| V3 | The codec, the law and their pins | same script, record section | `record.go:42 ResultRef string json:"resultRef"` (the last field), `:46 func EncodeRecord`, `:60 dec.DisallowUnknownFields()`, `:78` positional literal, `:100–101` failed arm requires `Denial == ""` and a zero `ResultRef`; `effectbroker.ail:154 export type EffectRecord`, `:164 resultRef: HashRef`, `:170 recordConsistent`, `:172` the failed arm; `broker_test.go:598 TestRecordGoldenBytes`, `:605 const want`, `:622 const wantFailed`; `decide_test.go:88 {187, EffectRecord{… Failed: true}, true}`; `git grep effectbroker 77c9e3d4 -- '*_test.go'` → only a comment (`effectplan_drift_test.go:23`) and the gate manifest (`module_manifest_gate_test.go:566`): no field-list drift test |
| V4 | Refusal, denial and failure are distinct today | same script | `handlers_exec.go:96 func execRefusal`, `:97 {"ok": false, "refused": why}`; `broker.go:198 Denial: decision.Label,` (denied arm), `:311 resultObj := resultObject(result)` (a refusal is a success with a result); `decide.go:11 LabelDeniedBudget = "denied:budget"`; `episode_report.go:155 reportReasonMax = 240`, `:533 func handlerRefusal` |
| V5 | Replay demands byte equality, and failed records already replay | same script | `effectful.go:485 if !bytes.Equal(outBytes, committed.OutputBytes) {`; `effectful_test.go:521 t.Run("denied_and_failed_records_replay", …)` |
| V6 | `workspace-exec` plans have no finish phase | `git show 77c9e3d4:host/coordinator/setools_exec_test.go \| grep -n '"finish":false'` | `` 30: payload + `}],"finish":false,"result":null}` `` |
| V7 | Execute's error returns at HEAD | `bash -c 'git show 77c9e3d4:host/broker/handlers_exec.go \| grep -nE "return (nil, )?(err\|runErr\|fmt\.Errorf)" \| awk -F: "(\$1>=251 && \$1<=391) \|\| (\$1>=522 && \$1<=555)"'` | 14 lines: `253`, `256` (`fmt.Errorf`), `277`, `291`, `294`, `297`, `301`, `304` (`return nil, err`), `315` (`fmt.Errorf … budget is spent`), `356`, `370` (`return nil, runErr`), `527`, `538` (`return err`), `551` (`fmt.Errorf … seed cache %s from %s: %w`) |
| V8 | Cause texts carry host paths and OS or child text | same script | `exec_sandbox.go:107 "exec confinement stack drifted: %s: %s"`, `:208` and `:248 Why: err.Error()`; `exec_settings.go:147 "exec settings %s is not World's rendering: %s"`; `handlers.go:74 "handler subprocess failed: %v (output %q)"` |
| V9 | The paths a success result already discloses | same script | `handlers_exec.go:231 Argv0 json:"argv0"`, `:232 Argv0Target json:"argv0_target"`, `:240 Node json:"node"`, `:241 NodeTarget json:"node_target"` |
| V10 | Nothing at HEAD can hold a reason (the search has a positive control); the record's readers | same script, last section; `git grep -n "DecodeRecord(" 77c9e3d4 -- "*.go" \| grep -v _test.go \| grep -v "func DecodeRecord" \| wc -l`; `git show 77c9e3d4:cmd/ailang-worldd/why.go \| sed -n 912,925p` | `failure field in record.go: 0 (control json:"denial": 1)`; `reason in worldEffect: 0 (control json:"record": 2)`; `10` decode call sites; `why.go:915 if !d.ResultRef.IsZero() {`, `:918 resText = fmt.Sprintf("result %d B", …)`, `:920 delta = resultWorktreeDelta(res.Payload)`; `why.go:525 … obj.SemanticID != broker.EffectRecordV1`; `broker.go:433 rec.Allowed && (rec.Denial != "" \|\| …)` |
| V11 | The prototype passes the full gate (r1, before the U+202E mapping and the A2A half of TP) | `VERIFY_GO_LEGS=plain AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_go.sh` → `/tmp/iter360-designer/verify_go.txt` (a first run without `AILANG_BIN` was red only on the tests that require it, e.g. `demo_test.go:292: required AILANG_BIN absent`; it is discarded) | `✓ go gate PASSED (plain): build clean; legs run pass with pinned AILANG_BIN (AILANG v0.41.0`; every package `ok`, including `host/broker 306.763s`, `host/daemon 341.595s`, `host/verifygate 321.620s`, `cmd/ailang-worldd 189.658s`; `RC=0` |
| V12 | The mutation table | `/tmp/iter360-designer/mutate.py`: for each mutant, exact-string replace (anchor count must be 1), then `go test -count=1 -timeout 300s -v -run '^(<13 names>)$' ./host/broker/ ./host/coordinator/ ./host/projection/ ./cmd/ailang-worldd/`, then parse `--- FAIL:`/`--- PASS:` and restore → `/tmp/iter360-designer/mutants.txt`; loads 6.0–9.0 | `CONTROL: rc=0 pass=13 red=-`; `MUT-HEAD-SHAPE: rc=1 pass=10 red=TC,TP,TX`; `MUT-NO-OMITEMPTY: red=Gold`; `MUT-RECORD-DROPS-LINE: red=TC,TX`; `MUT-REPLAY-DROPS-REASON: red=TC`; `MUT-COORD-DROPS-REASON: red=TC,TP`; `MUT-WORLD-NO-OMITEMPTY: red=TC`; `MUT-ERR-TEXT`, `MUT-RAW-PATH`, `MUT-DRIFT-WHY-RAW: red=TK,TX`; `MUT-CONTROL-ONLY`, `MUT-NO-BOUND: red=TL`; `MUT-CLASS-OPEN: red=TG,TL`; `MUT-DECODE-NO-CHECK`, `MUT-FAILURE-ANY-ARM: red=TG`; `MUT-HF-ERROR-IS-LINE: red=Reverify,TX`; `MUT-SEED-UNCLASSIFIED: red=TX`; `MUT-SEED-TIMEOUT-AS-EXIT: red=TK,TX`; `MUT-LOG-MCP-NONE`, `MUT-LOG-A2A-NONE`, `MUT-LOG-RECONCILED: red=TP`; `MUT-WHY-NO-LINE: red=TW`; `MUT-CALL-NO-LINE: red=TR`; `MUT-PUT-NO-CHECK: rc=0 pass=13 red=-`; `MUT-CLEANUP-UNWIRED: rc=0 pass=13 red=-`. `MUT-NO-SANITISE`'s first anchor (`if false {`) left `unicode` unused (`BUILD-FAIL`); re-anchored as `if unicode.IsControl(r) && false {` → `red=TL`. Restore: `cmp proto-before-mutants.diff proto-after-mutants.diff` → `tracked diff identical after mutants`; `failure.go identical` |
| V13 | Stability and package runs on the final prototype | `go test -count=3 -run '^(<10 names>)$' ./host/broker/ ./host/coordinator/ ./host/projection/ ./cmd/ailang-worldd/`; `go test -count=1 ./host/projection/`; `gofmt -l host cmd`; `go vet` of the four packages | `ok` ×4 (broker 0.728s, coordinator 0.761s, projection 0.398s, cmd 0.802s); projection package `ok 13.596s`; gofmt empty; `VET-OK` |
| V14 | The patch applies at HEAD, and the tree went back to the doc state | `git add -N <5 new files> && git diff > design_docs/planned/w-exec-failure-carries-a-reason.patch && git reset -q <5 new files>`; `git checkout -- <10 modified files>`; `rm <5 new files>`; `git status --porcelain`; `git apply --check -v design_docs/planned/w-exec-failure-carries-a-reason.patch` | patch 1184 lines, sha256 prefix `b7dcb7821c3b389a`, 15 `diff --git` sections; numstat `mcpclient.go 4/0`, `why.go 4/0`, `why_test.go 44/1`, `QUICKSTART.md 11/0`, `broker.go 13/5`, `handlers_exec.go 24/17`, `record.go 21/4`, `effectful.go 9/2`, `mcp.go 1/0`, `projection.go 32/0`; new files 189 + 278 + 62 + 106 + 32 lines; `git status --porcelain` → `?? design_docs/planned/w-exec-failure-carries-a-reason.patch` only; `Checking patch …` ×15, `apply-check rc=0` at `77c9e3d4` |
| V15 | The patch, applied to a fresh `git archive 77c9e3d4`, builds and its tests pass | in `/tmp/iter360-designer/applied`: `git apply <patch>`; `go build ./... && go vet <4 pkgs>`; `test -z "$(gofmt -l host cmd)"`; `go test -count=1 -v -run '^(<13 names>)$' <4 pkgs>`; then `AILANG_BIN=$HOME/.pinned-ailang/ailang go test -count=1 ./host/broker/ ./host/coordinator/ ./host/projection/ ./cmd/ailang-worldd/ ./host/daemon/` | `APPLY-OK`, `BUILD-VET-OK`, `GOFMT-OK`; 13 `--- PASS` (TL, TG, TK, TX 0.18s, Gold, Sketch, CFJ2, Reverify 2.04s, TC, Replay, TP, TW, TR) and `ok` ×4; full packages `ok host/broker 244.644s`, `ok host/coordinator 65.213s`, `ok host/projection 13.905s`, `ok cmd/ailang-worldd 125.260s`, `ok host/daemon 291.581s` (it holds the QUICKSTART-bound tests), `RC=0` |
| V16 | Downgrade: HEAD cannot decode a record carrying `failure` (R1), but decodes one without it | in `/tmp/iter360-designer/head` (`git archive 77c9e3d4`) a throwaway `zz_downgrade_probe_test.go` calls HEAD's `DecodeRecord` on TG's golden, with and without the field; `go test ./host/broker/ -run '^TestZZHeadDecodesRow255Record$' -v` | `err = broker: decode effect record: json: unknown field "failure"`; without the field `err = <nil>` |
| V17 | Every `Execute` error is classified; `HandlerFailure` is built only in two files | in `applied`: `awk 'NR>=254 && NR<=400' host/broker/handlers_exec.go \| grep "return nil, " \| grep -vc "execFailure\|execRefusal"`, with the classified count as control; `grep -rn "HandlerFailure{" --include="*.go" host cmd \| grep -v _test.go \| grep -v "failure.go\|handlers_exec.go" \| wc -l`, with test files as control | `unclassified returns in Execute: 0 (control, classified: 10)`; `prepareCache` bare returns 0; other non-test constructors `0`; control: `failure_test.go`, `effect_failure_reason_test.go`, `effect_failure_log_test.go` |
| V18 | No `.ail` change; the `.ail` gate is unaffected | `AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_ail.sh` (worktree at HEAD plus the untracked patch); `grep -c "^diff --git.*\.ail" <patch>` with control `"^diff --git.*\.go"` | `✓ verify gate PASSED: 16 required identities verified, 40 named tests pass, 503 se-tools named tests pass`; rc 0. `.ail` paths in the patch `0` (control `.go` paths `14`) |
| — | **UNMEASURED** | — | Linux (R6: GNU `cp` exit code and `EACCES` on the CI runner; the acceptance gate is the PR's CI); the home store's entries 54–64 for episode `worldd-2611031015` (row claim 8); `verify_go.sh` on the **final** patch (V11 ran on r1; V15 ran the five touched packages in full on the final patch); a `-race` run |

## 11. Quorum log

*(empty: r1, awaiting quorum)*
