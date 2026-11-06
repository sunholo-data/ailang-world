# w-exec-failure-carries-a-reason — a failed `workspace-exec` call records, returns and logs one bounded failure line

- **Status**: Planned (designed iteration 360, **r3**: r2 by the designer after quorum round 1 BLOCKED; r3 is the controller's narrow-refinement carve-out after round 2, §11). Queue row 255, in D-WORLD-108's scope after
  row 250 (attended 2026-11-06: rows 250, then 255 → 257). Rows 256 (`w-daemon-startup-under-load`)
  and 257 (`w-exec-cache-seed-at-bind`) are separate rows. §9 says how this row's classes serve
  them, and designs neither.
- **Item**: Clause 5 (MET; operator and agent UX). It adds an optional field to the
  `world/effect-record/v1` wire, so its codec choice (§3.1) is the decision a quorum should check.
  **[r2]** Because a binary from before this change cannot decode that field, r2 also raises the
  world store's schema version from 5 to 6 just before the first such record is written. Every
  already-shipped binary refuses a version-6 store at open, so an older binary stops cleanly
  instead of failing mid-read (§3.7).
- **Designer lane**: `claude:claude-opus-5-5` via Agent opus.
- **Estimated**: ~0.75 day. A complete patch is
  `design_docs/planned/w-exec-failure-carries-a-reason.patch`. **[r2]** It is 1538 lines over 22
  files, sha256 prefix `02259c3bdc1c942b`. It applies with `git apply --check` at `77c9e3d4`. Its
  18 named tests pass on a fresh `git archive` of `77c9e3d4` with the patch applied (V26). The full
  plain gate `verify_go.sh` passes on the byte-identical prototype (V24).
  - **Production code**: +162/−33 lines in 9 files. r2 adds +45/−5 in `host/store/store.go` and
    +9 in `broker.go`. There is also one new 189-line file, `host/broker/failure.go`, of which
    about a third is comments.
  - **Tests**: four new test files (561 lines). Edited tests: `why_test.go` +44/−1,
    `schema_version_test.go` +71/−3, `deadline_guard_test.go` +1, and four broker test stores
    +4 each.
  - **Docs**: one `docs/QUICKSTART.md` paragraph (S7).
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
| **(a) A new field `failure`, `omitempty`, last in `recordWire`, plus the store schema marker of §3.7 [r2]** | **CHOSEN** | `encoding/json` omits an empty string, so every record without a line encodes to the bytes it has today. That covers every success, every denial, and every failure of a handler that is not `ExecHandler`. Both existing goldens stay green (`TestRecordGoldenBytes`, V15), and `MUT-NO-OMITEMPTY` reddens them (§5). A stored record without the field decodes to `Failure: ""`. No stored object's hash changes, because content-addressed objects are never re-encoded. Only failed `Workspace.Exec` records written from now on carry it |
| (b) Reuse `denial` | **Rejected** | It breaks the record law. `RecordConsistent` requires `Denial == ""` on every allowed record, failed ones included (`record.go:100–101`), mirroring the sketch's ensures clause (`effectbroker.ail:172`). Replay refuses an allowed record with a non-empty denial (`broker.go:433`). It also blurs a capability denial with a handler failure (§3.4) |
| (c) A failure object referenced by `resultRef` | **Rejected** | It also breaks the law: a failed record must have a zero `resultRef` (`record.go:101`, `effectbroker.ail:172`). Readers treat a non-zero `resultRef` as a success result: `why` prints `result N B` and reads its worktree delta (`why.go:915–923`), and the episode report reads it as a handler answer (`episode_report.go:533`) |
| (d) A `world/effect-record/v2` | **Rejected** | An older binary cannot read a v2 record either, so (d) would need the same marker (§3.7), and it costs much more. Ten non-test call sites decode records (V10), and several check `SemanticID == EffectRecordV1` (`why.go:525`, `broker.go:418`, `workbench_packet_evidence.go`, `episode_report.go`), so a v2 would need a dual-version reader in each. (a) keeps one codec, and old records stay valid v1 |

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
  the sketch should gain the field for documentation parity is open question Q4.
- **The one cost: an older binary cannot decode the field.** `DecodeRecord` at `77c9e3d4`
  answers `json: unknown field "failure"`, and the same record without the field decodes (V16).
  **[r2]** r1 called this acceptable on the assertion that World never downgrades against a live
  store. No Verification Log row backed that assertion, and it is false: `service upgrade` rolls
  back to the previous build automatically (V23). §3.7 replaces the assertion with a measured
  answer and an enforced gate.

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
  from the journal) logs nothing. **[r2]** The helper writes one line per failed effect that
  carries a line. A call has at most one effect **by construction**, because plan law L1 refuses
  a second one before any effect runs (V22). So each call that ran logs at most once (Claim E).
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

### 3.7 [r2] Compatibility: who can read a store that holds a failure line

**What was measured** (V19–V23, V25). The facts that decide the design:

- **No written policy.** No store-compatibility, downgrade or rollback policy is written down. A
  search of the charter, `DESIGN.md`, the coding standards, `docs/` and `CLAUDE.md` finds one
  unrelated hit, the word "downgraded" used about an AC's claim (V19).
- **One enforced gate.** Every shipped binary checks the world store's SQLite `user_version` at
  open. Version 5 is current. Any higher version is refused with
  `store: schema version future: … has user_version N; binary supports 5; use a compatible
  binary` (`store.go:464–465`, message `:68`). HEAD's own test opens a store at version **6** and
  asserts exactly that refusal (`TestFutureVersionIsRejected`, `frozenFutureSchemaVersion = 6`,
  V21). `schema.sql` never sets the version (0 hits).
- **Rollback is real.** `service upgrade` keeps the old build as `.prev`
  (`service_upgrade.go:210`, `:219`). If the new build does not come up, or its upgrade cannot be
  recorded, it restores `.prev` over the binary and waits for that build to serve the **same
  store** (`:240`, `:256`, `:262`). If that also fails it reports
  `rollback failed too: the service is stopped` (`:287`). An operator can also put `.prev` back by
  hand.
- **One binary, read over HTTP.** The daemon (`serve`) and the reader verbs (`call`, `why`,
  `episode`) are one binary, `cmd/ailang-worldd` (`main.go:354, 395, 398, 434`). `why` and `call`
  do not open the store. They read the daemon's `/v1/objects/…` (`why.go:136`) and decode records
  in the CLI process. Only `episode start` refuses a CLI build that differs from the installed one
  (`episode.go:330`).
- **Replay has no production caller.** `Coordinator.Replay`, `OpenReplayBinder` and `host/replay`
  have no non-test caller (V21: 0 `.Replay(` calls, 0 importers of `host/replay`, against 58
  importers of `host/store`).

**Why a carrier that old readers tolerate does not exist.** Two facts close it:

- Every record reader decodes with `DisallowUnknownFields` (`record.go:60`), and the record law
  forbids reusing `denial` or `resultRef` (§3.1). So no record carrier exists that an old binary
  can read.
- Leaving the record unchanged and putting the line only in the committed output fails the other
  way. The output must be recomposable from the records, byte for byte (`effectful.go:485`), so a
  reason present only in the output could not be replayed.

The row asks for the line in the record, so the remaining question is how to stop an old binary
before it misreads the record.

**The gate: raise the existing schema version, lazily.**

- A new constant, `failureLineSchemaVersion = 6`, means "schema 5's tables, and the store may hold
  an effect record with a `failure` line". New binaries open 5 and 6, both writable and read-only,
  and refuse 7 and above with `binary supports 6`.
- Fresh stores still start at 5.
- Before the broker writes the **first** record that carries a line, it calls
  `Store.MarkFailureLines`. That method raises `user_version` from 5 to 6 under the writer and
  re-reads it. It is idempotent at 6, and it refuses any other version, including a read-only
  handle.
- **[r3]** If `MarkFailureLines` fails, the broker must not downgrade the record to the legacy
  shape or report a normally committed classified failure. It surfaces an explicit
  persistence/compatibility error through the broker's verified storage-error path: it writes no
  record and returns `&OutcomeWriteError{Effect, Scope, Err: <the mark error>}`, exactly as a failed
  `putRecord` does at `broker.go:297–300`. Because the handler has already executed, this path
  preserves the debit/accounting invariants (the grant was debited at `broker.go:260`, before
  `Execute`, and is not refunded) and does not automatically rerun the effect. A record with a line
  therefore never sits in a store that an older binary would still open (Claim F), and no failure
  is ever committed without its line. (r2 wrote the record without the line here; all three r2
  reviewers rejected that as a silent fallback, §11 round 2.)

The gate needs no new enforcement code in old binaries. It uses the refusal they already ship.
The marker is a version number only: no table, index or row changes. It extends what
`user_version` means, from "the DDL" to "the DDL, and which record shapes may be present". Q2
asks Mark to ratify that extension.

**Supported combinations** (V20 is a live measurement; the rest are from code and tests):

| Readers × store | Store at 5 (no line written yet) | Store at 6 (a line written) | How it is held |
|---|---|---|---|
| this build: daemon, `doctor`, `session list`, `call`, `why` | reads and writes as today | reads and writes; `why` and `call` print the line | TS (reopen as writer and reader at 6) |
| a build from **before** this change, opening the store (`serve`, `doctor`, read-only verbs) | works: no record can carry a line | **refused at open**: `store: schema version future: "…" has user_version 6; binary supports 5; use a compatible binary`, as writer and as reader | the shipped binary's own gate (`store.go:464`); measured with HEAD code on a store the patched broker marked (V20) |
| `service upgrade` **to** this build, then the automatic **rollback** | works: a build that never came up wrote no line | the old build refuses the store, so the rollback ends `rollback failed too: the service is stopped`, with both builds kept and named. A line exists only if the new build served a failed exec call before the rollback | the same gate. This is the trade-off Q2 names |
| a CLI from **before** this change (`why`, `call`) against a daemon of this build, over HTTP | unchanged | **not prevented**: `call` prints the effect without the reason (its decode is not strict). `why` marks that one effect link ✗ with `broker: decode effect record: json: unknown field "failure"`; the other links still verify | no gate: a CLI does not open the store. Q1 |
| a CLI of this build against an older daemon | unchanged | not reachable: an older daemon refuses a version-6 store | — |
| `Coordinator.Replay` of this build | old and new invocations replay byte-identically (TC, `TestEffectReplay`) | the same | no production caller (V21) |
| a future build, at 7 or above | refused by this build with `binary supports 6` | refused | `TestFutureVersionIsRejected`, its future version moved from 6 to 7 |

**Rejected alternatives.**

- **Raise 5 to 6 at every writer open.** It is simpler, but it would break every `service upgrade`
  rollback to an older build, including rollbacks where no line was ever written. The lazy mark
  breaks only the ones where a line exists.
- **A policy with no enforcement** ("a downgrade is unsupported"). That would be a policy change
  that only Mark can ratify, and it would leave old readers failing mid-read.
- **A new table or flag that an old binary does not check.** Old binaries check only
  `user_version`, so it would not be enforced.

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

**Claim E (logged at most once per call).** *Proof.* `logEffectFailures` writes only after a
`Dispatch` whose result is not `Reconciled`, and writes one line per world-block effect with a
non-empty `reason`. **[r2]** That world block has at most one effect, by construction:

- The only non-test code that builds a plan's effect list is `parsePlan`'s loop
  (`effectplan.go:135`).
- That loop runs only after `if !effectCountOk(n)` has refused n > 1 under law L1 (`:125–126`;
  `maxPlanEffectsV1 = 1`, `:32`).
- `parsePlan` is the gate for both the live path (`effectful.go:289`) and replay (`:430`).
- L1 is a Z3-verified contract in `design_docs/sketches/effectplan.ail`
  (`ensures { result == (n >= 0 && n <= 1) }`, `:40–41`), bound to the Go by the drift test
  (`effectplan_drift_test.go:35`).
- The control shows the search is not empty: two existing tests submit two-effect plans and get L1
  refusals (`TestParsePlanLaws/L1/two-effects`, `TestProducerMintsNoPacketForMultiEffectPlan`;
  both PASS at HEAD, V22).

So a call logs at most one line. If L1 is ever widened (R-SE-10 reserves that), the helper's rule
is still one line per failed effect that carries a line. ∎ TP sends the same A2A task twice, which
runs once, and expects exactly one line.

**Claim F [r2, r3] (no lined record in a store an older binary would open, and no silent
fallback).** *Proof.* The broker computes `reason` and then, before `putRecord`, calls
`MarkFailureLines`. **[r3]** If the mark fails, the broker writes no record and returns
`*OutcomeWriteError` wrapping the mark error, through the same path a failed `putRecord` takes
(`broker.go:297–300`); the coordinator handles that type at `effectful.go:184`, and a `Dispatch`
error is logged by `logRefusal` (`mcp.go:131–133`), so the failure is loud. `MarkFailureLines`
returns nil only when the store's `user_version` reads 6 after the call. Every binary from before
this change refuses a store at 6 at open (`store.go:464`, V20). ∎ TV pins the order
(`mark,record`), the mark-failure error (no record written, `*OutcomeWriteError` returned, the
handler executed exactly once, the grant still debited) and the unchanged version for a failure
with no line. TS pins the marker itself.

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
| **TV** [r2] | `TestFailureLineMarksTheStoreFirst` (`host/broker/failure_test.go`) | over a file-backed store: a lined failure raises `user_version` to 6, with the call order exactly `mark,record`; a plain failure leaves it at 5; **[r3]** a store whose mark fails (fault-injected) gets **no** record and no committed legacy failure, the call returns `*OutcomeWriteError` wrapping the mark error, the handler executed exactly once, and the grant stays debited |
| **TS** [r2] | `TestMarkFailureLinesRaisesFiveToSixOnce` (`host/store/schema_version_test.go`) | a fresh store is 5; a read-only handle cannot mark (it stays 5); marking twice gives 6 both times; `WorldSchemaVersion` reads 6; a writer and a reader of this build reopen it without rewriting the version; a schema-4 store refuses to be marked |
| edited [r2] | `TestFutureVersionIsRejected`, `TestOpenReadOnlyEnforcesVersionFive` (via `frozenFutureSchemaVersion` 6 → 7, `Current` 5 → 6); `TestStoreDeadlineGuardDynamic` (one fixture: the new method is quarantine- and deadline-guarded like every store method); 4 broker test stores delegate `MarkFailureLines` | the store still refuses a version above what it supports, now 7 and up |

The row asks for "one test per failure class in `Execute`, showing that the record, the MCP result
and `why` each carry it". The design meets it in two layers. **TX runs once per class** at the
place the class is decided, and asserts the record and the `EffectFailedError.Reason` that the
coordinator copies into the MCP result. **TC, TP and TW** then prove each downstream surface
carries the line it is given. Past the broker the line is an opaque string, so repeating those
three per class would test the same carrier ten times. `MUT-COORD-DROPS-REASON`,
`MUT-WHY-NO-LINE` and `MUT-LOG-MCP-NONE` show each carrier is pinned.

**Mutation run** (V12 for r1; **[r2]** V25 re-ran every mutant on r2 with 18 tests across 5
packages, adding TV, TS, `Future`, `ROFive` and `Guard`. Every r1 kill pattern reproduced; the
only changes are that TV joins the reds of MUT-RECORD-DROPS-LINE and MUT-HEAD-SHAPE). Each mutant
is one exact-string edit of the prototype (MUT-MARK-AFTER-RECORD is two). The 13 r1 tests run once (`-count=1`) across the four packages, and then the file is
restored, with the tracked diff compared to its pre-run copy (`cmp`: identical). The columns list
the tests that went red. `Gold` = `TestRecordGoldenBytes` and `Reverify` =
`TestExecHandlerReverifiesTheStackBeforeEachCall`, both existing. `CFJ2`, `Sketch` and `Replay`
(`TestEffectReplay`) were in the set and stayed green under every mutant.

| Mutant | Edit | Red | Expected → observed |
|---|---|---|---|
| CONTROL | none | — (13 of 13 PASS) | green → green |
| **MUT-HEAD-SHAPE** (positive control) | the record drops the line **and** the coordinator drops the reason: today's behaviour inside the new structure | **TX, TC, TP** (r2: **and TV**) | red → red. TX's text: `record = broker.EffectRecord{… Failed:true, Denial:"" …}` with no failure, 11 of 11 rows; TC's output `{"world":{"effects":[{"id":"e1","record":"sha256:…","status":"failed"}],…}` |
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
| **MUT-NO-MARK** [r2] | the broker never calls `MarkFailureLines` | TV | red → red |
| **MUT-MARK-ERROR-IGNORED** [r2] | a failed mark still writes the line | TV | red → red |
| **MUT-MARK-FAIL-SILENT** [r3] | a failed mark writes the record without the line and returns the ordinary `*EffectFailedError` (r2's fallback) | TV | red → red (the planner re-runs it on the r3 patch) |
| **MUT-MARK-AFTER-RECORD** [r2] | the mark moves after `putRecord` (a crash between them would leave a lined record in a version-5 store) | TV | red → red (`call order = [record mark]`) |
| **MUT-STORE-REJECTS-6** [r2] | `enforceSchemaVersion` does not accept 6 | TS | red → red |
| **MUT-MARK-NOOP** [r2] | `MarkFailureLines` returns nil without writing | TS, TV | red → red |
| **MUT-MARK-ANY-VERSION** [r2] | `MarkFailureLines` tries to mark any version | TS | red → red (the schema-4 row) |
| MUT-FUTURE-BOUND-5 [r2] | the future check reads `> currentSchemaVersion` | — | **SURVIVES, equivalent.** Version 6 returns from the branch above it, so `> 5` and `> 6` refuse the same set (7 and up). No test can tell them apart |

## 6. Residuals and risks

- **R1 [r2, replaces r1's R1], older binaries.** r1 asserted that World never downgrades. That
  was unsupported, and false (V23). r2 makes an older build **refuse** a store that holds a line,
  at open, through the gate it already ships (§3.7, V20). Two costs remain, both stated in §3.7's
  table:
  - Once a line exists, a `service upgrade` rollback to a build from before this change ends with
    the service stopped and both builds kept. Before this change it would have served. Q2.
  - An older **CLI** reading a newer daemon over HTTP is not prevented. Its `why` marks one link ✗
    naming the unknown field. Q1.
- **R2, the operator log carries the bounded line, not the full Go error.** The row says "the same
  line … logged". The full cause stays in `HandlerFailure.Err` and is dropped at `requestEffect`,
  as at HEAD. When a line is not enough (for example `exec-tools` or `sandbox-start`), the operator
  still has to reproduce the failure. Plumbing the cause to the projection's operator-only log is
  Q3.
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
go build ./... && go vet ./host/broker/ ./host/coordinator/ ./host/projection/ ./cmd/ailang-worldd/ ./host/store/ && test -z "$(gofmt -l host cmd)" && echo OK
# expect: OK
go test -count=1 -v -run '^(TestFailureLineBoundsAndSanitises|TestRecordFailureGoldenBytes|TestExecFailureClassifies|TestExecFailureClassesReachTheRecord|TestFailureLineMarksTheStoreFirst|TestRecordGoldenBytes|TestExecHandlerReverifiesTheStackBeforeEachCall|TestCFJ2HandlerErrorKeepsDebitAndWritesOneFailureRecord|TestSketchRows|TestEffectFailureReasonReachesTheWorldBlockAndReplays|TestEffectReplay|TestMCPEffectFailureLineReturnedAndLoggedOnce|TestWhyShowsTheFailureLine|TestCallPrintsTheFailureLine|TestMarkFailureLinesRaisesFiveToSixOnce|TestFutureVersionIsRejected|TestOpenReadOnlyEnforcesVersionFive|TestStoreDeadlineGuardDynamic)$' \
  ./host/broker/ ./host/coordinator/ ./host/projection/ ./cmd/ailang-worldd/ ./host/store/ | grep -E '^(--- |ok|FAIL)'
# expect: 18 lines "--- PASS: <name>", 5 "ok", no FAIL (V26)
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
`MUT-ERR-TEXT` (TK, TX), `MUT-RAW-PATH` (TK, TX), `MUT-LOG-RECONCILED` (TP),
`MUT-HF-ERROR-IS-LINE` (Reverify) and, **[r2]**, `MUT-NO-MARK` (TV), `MUT-MARK-AFTER-RECORD` (TV)
and `MUT-STORE-REJECTS-6` (TS). Record the three survivors with their reasons.

**[r2] Old-binary refusal, on the applied patch:** reproduce V20. Mark a store with the patched
broker, then open it with a `git archive 77c9e3d4` build's `store.Open` and `store.OpenReadOnly`.
Both must return `store: schema version future: … has user_version 6; binary supports 5; use a
compatible binary`. Record the output in the sprint's Verification Log.

**Acceptance gate: CI on the PR.** The plain `go build + test gate` step (`verify_go.sh`) runs TX
on Linux (R6). The item is not accepted until that step is green on the PR head. A local green is
not a substitute.

## 8. Files

| File | Change |
|---|---|
| `host/broker/failure.go` | **new**, 189 lines: `FailureLineMax`, the 10 class constants and their closed set, `HandlerFailure` (`Error`/`Unwrap` = the cause's, `Line`), `FailureLine`, `ValidFailureLine`, `execDisclosed` and `name`, `execFailure`, `driftWhat`, `settingsWhat`, `seedWhat` |
| `host/broker/record.go` | +21/−4: `EffectRecord.Failure`, `recordWire.Failure` (`omitempty`, last), the codec's two positional literals, `DecodeRecord`'s check, new `FailureConsistent` |
| `host/broker/broker.go` | +22/−5: `EffectFailedError.Reason`; the failed record and error carry `classified.Line()`; **[r2]** `objectStore` gains `MarkFailureLines`, which is called before a lined record ([r3] a failed mark returns `*OutcomeWriteError` and writes no record); replay sets `Reason: rec.Failure`; `putRecord` checks `FailureConsistent` |
| `host/store/store.go` [r2] | +45/−5: `failureLineSchemaVersion = 6`; `enforceSchemaVersion` accepts 6 like 5, re-reads the version it found after applying the schema, and refuses only above 6; new guarded `(*Store).MarkFailureLines` |
| `host/store/schema_version_test.go`, `deadline_guard_test.go` [r2] | +71/−3, +1: TS; `frozenFutureSchemaVersion` 6 → 7 and the future `Current` 5 → 6; the guard fixture for the new method |
| `host/broker/{broker,handlers,handler_error_repro,registry_publish}_test.go` [r2] | +4 each: the test stores delegate `MarkFailureLines` |
| `host/broker/handlers_exec.go` | +24/−17: the 11 `Execute` returns go through `execFailure` (the effect and scope guards merge into one); `prepareCache`'s 3 returns build `*HandlerFailure` |
| `host/coordinator/effectful.go` | +9/−2: `EffectResult.Reason` (`json:"-"`), `worldEffect.Reason` (`omitempty`), the copy in `requestEffect` and `composeEffectOutput` |
| `host/projection/projection.go`, `mcp.go` | +32, +1: `logEffectFailures`, called after a successful `Dispatch` on A2A and MCP |
| `cmd/ailang-worldd/why.go`, `mcpclient.go` | +4, +4: `no result: <line>`; `call`'s `— <line>` |
| `docs/QUICKSTART.md` | +15: the "When a call fails" paragraph (S7); **[r2]** it names the 5 → 6 marker and the rollback consequence |
| `host/broker/failure_test.go` | new, 361 lines: TL, TG, TK, TX, **[r2]** TV, and the `expiredCtx` wrapper |
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
  the record's semantic id (Q4).
- No episode-report change. It still counts `failed` effects without their class. Row 257's test
  for a visible bind-time failure is that row's.
- No reason for handlers other than `ExecHandler`. Their failed records keep today's bytes
  (Claim A). Other handlers can opt in later by returning a `*HandlerFailure` with their own
  classes added to the closed set.
- No retry, no backoff, no new seam in production code. The tests use the existing fixture, the
  existing `ExecToolState.Close`, and a test-only `expiredCtx` wrapper.
- No change to `tools/launchd/*`, skills, `~/.ailang/state/*` or `HOST_CHANGELOG.md` (that belongs
  to the release that ships this).

- **[r2]** No build-skew check for the read verbs (`why`, `call`). Today only `episode start`
  refuses a CLI build other than the installed one. Q1.

**Open questions.** Q1 and Q2 are **for Mark**: they touch compatibility policy, which no written
document states (V19).

- **Q1 (Mark).** Should the read verbs (`why`, `call`) refuse to read a daemon of a different
  build, as `episode start` does (`episode.go:330`)? Without that, an older CLI reading a newer
  daemon marks the failed effect's link ✗ naming the unknown field. (Default for this row: no
  enforcement. The effect is one visible ✗ link, and a build check for every read verb is a
  separate change.)
- **Q2 (Mark).** Ratify §3.7's use of the world store's `user_version`. It would mean "which
  record shapes may be present" as well as the DDL. The consequence: after the first failure line,
  the store refuses builds from before this change, so a `service upgrade` rollback across this
  change stops the service instead of serving. (Default: yes. It is the only refusal the shipped
  binaries already enforce, and the alternative is old readers failing mid-read.)
- **Q3.** Should the operator-only log also carry the full cause, `%q`-quoted like `logRefusal`?
  It would need `EffectFailedError`'s cause carried past `requestEffect`. (Default: no, the row
  says "the same line".)
- **Q4.** Should `effectbroker.ail`'s `EffectRecord` gain `failure: string` for documentation
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
| V11 | (r1) The prototype passes the full gate (r1, before the U+202E mapping and the A2A half of TP) | `VERIFY_GO_LEGS=plain AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_go.sh` → `/tmp/iter360-designer/verify_go.txt` (a first run without `AILANG_BIN` was red only on the tests that require it, e.g. `demo_test.go:292: required AILANG_BIN absent`; it is discarded) | `✓ go gate PASSED (plain): build clean; legs run pass with pinned AILANG_BIN (AILANG v0.41.0`; every package `ok`, including `host/broker 306.763s`, `host/daemon 341.595s`, `host/verifygate 321.620s`, `cmd/ailang-worldd 189.658s`; `RC=0` |
| V12 | (r1) The mutation table | `/tmp/iter360-designer/mutate.py`: for each mutant, exact-string replace (anchor count must be 1), then `go test -count=1 -timeout 300s -v -run '^(<13 names>)$' ./host/broker/ ./host/coordinator/ ./host/projection/ ./cmd/ailang-worldd/`, then parse `--- FAIL:`/`--- PASS:` and restore → `/tmp/iter360-designer/mutants.txt`; loads 6.0–9.0 | `CONTROL: rc=0 pass=13 red=-`; `MUT-HEAD-SHAPE: rc=1 pass=10 red=TC,TP,TX`; `MUT-NO-OMITEMPTY: red=Gold`; `MUT-RECORD-DROPS-LINE: red=TC,TX`; `MUT-REPLAY-DROPS-REASON: red=TC`; `MUT-COORD-DROPS-REASON: red=TC,TP`; `MUT-WORLD-NO-OMITEMPTY: red=TC`; `MUT-ERR-TEXT`, `MUT-RAW-PATH`, `MUT-DRIFT-WHY-RAW: red=TK,TX`; `MUT-CONTROL-ONLY`, `MUT-NO-BOUND: red=TL`; `MUT-CLASS-OPEN: red=TG,TL`; `MUT-DECODE-NO-CHECK`, `MUT-FAILURE-ANY-ARM: red=TG`; `MUT-HF-ERROR-IS-LINE: red=Reverify,TX`; `MUT-SEED-UNCLASSIFIED: red=TX`; `MUT-SEED-TIMEOUT-AS-EXIT: red=TK,TX`; `MUT-LOG-MCP-NONE`, `MUT-LOG-A2A-NONE`, `MUT-LOG-RECONCILED: red=TP`; `MUT-WHY-NO-LINE: red=TW`; `MUT-CALL-NO-LINE: red=TR`; `MUT-PUT-NO-CHECK: rc=0 pass=13 red=-`; `MUT-CLEANUP-UNWIRED: rc=0 pass=13 red=-`. `MUT-NO-SANITISE`'s first anchor (`if false {`) left `unicode` unused (`BUILD-FAIL`); re-anchored as `if unicode.IsControl(r) && false {` → `red=TL`. Restore: `cmp proto-before-mutants.diff proto-after-mutants.diff` → `tracked diff identical after mutants`; `failure.go identical` |
| V13 | (r1) Stability and package runs on the r1 prototype | `go test -count=3 -run '^(<10 names>)$' ./host/broker/ ./host/coordinator/ ./host/projection/ ./cmd/ailang-worldd/`; `go test -count=1 ./host/projection/`; `gofmt -l host cmd`; `go vet` of the four packages | `ok` ×4 (broker 0.728s, coordinator 0.761s, projection 0.398s, cmd 0.802s); projection package `ok 13.596s`; gofmt empty; `VET-OK` |
| V14 | (r1) The r1 patch applies at HEAD, and the tree went back to the doc state | `git add -N <5 new files> && git diff > design_docs/planned/w-exec-failure-carries-a-reason.patch && git reset -q <5 new files>`; `git checkout -- <10 modified files>`; `rm <5 new files>`; `git status --porcelain`; `git apply --check -v design_docs/planned/w-exec-failure-carries-a-reason.patch` | patch 1184 lines, sha256 prefix `b7dcb7821c3b389a`, 15 `diff --git` sections; numstat `mcpclient.go 4/0`, `why.go 4/0`, `why_test.go 44/1`, `QUICKSTART.md 11/0`, `broker.go 13/5`, `handlers_exec.go 24/17`, `record.go 21/4`, `effectful.go 9/2`, `mcp.go 1/0`, `projection.go 32/0`; new files 189 + 278 + 62 + 106 + 32 lines; `git status --porcelain` → `?? design_docs/planned/w-exec-failure-carries-a-reason.patch` only; `Checking patch …` ×15, `apply-check rc=0` at `77c9e3d4` |
| V15 | (r1) The r1 patch, applied to a fresh `git archive 77c9e3d4`, builds and its tests pass | in `/tmp/iter360-designer/applied`: `git apply <patch>`; `go build ./... && go vet <4 pkgs>`; `test -z "$(gofmt -l host cmd)"`; `go test -count=1 -v -run '^(<13 names>)$' <4 pkgs>`; then `AILANG_BIN=$HOME/.pinned-ailang/ailang go test -count=1 ./host/broker/ ./host/coordinator/ ./host/projection/ ./cmd/ailang-worldd/ ./host/daemon/` | `APPLY-OK`, `BUILD-VET-OK`, `GOFMT-OK`; 13 `--- PASS` (TL, TG, TK, TX 0.18s, Gold, Sketch, CFJ2, Reverify 2.04s, TC, Replay, TP, TW, TR) and `ok` ×4; full packages `ok host/broker 244.644s`, `ok host/coordinator 65.213s`, `ok host/projection 13.905s`, `ok cmd/ailang-worldd 125.260s`, `ok host/daemon 291.581s` (it holds the QUICKSTART-bound tests), `RC=0` |
| V16 | Downgrade: HEAD cannot decode a record carrying `failure` (R1), but decodes one without it | in `/tmp/iter360-designer/head` (`git archive 77c9e3d4`) a throwaway `zz_downgrade_probe_test.go` calls HEAD's `DecodeRecord` on TG's golden, with and without the field; `go test ./host/broker/ -run '^TestZZHeadDecodesRow255Record$' -v` | `err = broker: decode effect record: json: unknown field "failure"`; without the field `err = <nil>` |
| V17 | Every `Execute` error is classified; `HandlerFailure` is built only in two files | in `applied`: `awk 'NR>=254 && NR<=400' host/broker/handlers_exec.go \| grep "return nil, " \| grep -vc "execFailure\|execRefusal"`, with the classified count as control; `grep -rn "HandlerFailure{" --include="*.go" host cmd \| grep -v _test.go \| grep -v "failure.go\|handlers_exec.go" \| wc -l`, with test files as control | `unclassified returns in Execute: 0 (control, classified: 10)`; `prepareCache` bare returns 0; other non-test constructors `0`; control: `failure_test.go`, `effect_failure_reason_test.go`, `effect_failure_log_test.go` |
| V18 | No `.ail` change; the `.ail` gate is unaffected | `AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_ail.sh` (worktree at HEAD plus the untracked patch); `grep -c "^diff --git.*\.ail" <patch>` with control `"^diff --git.*\.go"` | `✓ verify gate PASSED: 16 required identities verified, 40 named tests pass, 503 se-tools named tests pass`; rc 0. `.ail` paths in the patch `0` (control `.go` paths `14`) |
| V19 | [r2] No written store-compatibility, downgrade or rollback policy exists | `grep -rniEc "downgrade\|older binary\|newer binary\|store (format\|schema) version\|forward.compat\|backward.compat\|roll ?back to\|compatible binary" design_docs/world-mission.md design_docs/DESIGN.md design_docs/coding-standards.md docs/*.md CLAUDE.md`, control `grep -c " the " design_docs/world-mission.md` | `world-mission.md:1` (line 2379: "AC11's claim is downgraded in-doc", unrelated); every other file `0`; control `3214`. The controller's narrower search found the same (0 hits; control `the` in DESIGN.md = 236) |
| V20 | [r2] An older binary refuses a store that this change has marked, as writer and as reader | a throwaway test in the r2 prototype opens `/tmp/iter360-designer/marked.db` with `store.Open` and drives one lined failure through a live broker session; then, in a `git archive 77c9e3d4` copy, a throwaway test calls HEAD's `store.Open` and `store.OpenReadOnly` on that file → `/tmp/iter360-designer/old-binary-probe.txt` (both throwaway files deleted) | prototype: `record sha256:53633d27…`, `marked store user_version = 6`; HEAD: `Open(marked) = false, err = store: schema version future: "/tmp/iter360-designer/marked.db" has user_version 6; binary supports 5; use a compatible binary`; `OpenReadOnly(marked) = false`, the same error. (A first attempt under zsh aborted on a glob before the store was made, so HEAD opened a fresh path. That run was discarded and the probe was redone under `bash` with the file deleted first) |
| V21 | [r2] The shipped gate, its own test at version 6, and no production replay | run under `bash`, output in `/tmp/iter360-designer/compat-greps.txt`: `git show 77c9e3d4:host/store/store.go \| grep -n …`; `git show 77c9e3d4:host/store/schema_version_test.go \| grep -n …`; `git show 77c9e3d4:host/store/schema.sql \| grep -c user_version`; `git grep -n "\.Replay(\|OpenReplayBinder(" 77c9e3d4 -- "*.go" \| grep -v _test.go`; importers of `host/replay` with control `host/store` | `store.go:46 const currentSchemaVersion = 5`, `:464 if version > currentSchemaVersion {`, `:68 "… binary supports %d; use a compatible binary"`; `schema_version_test.go:18 const frozenFutureSchemaVersion = 6`, `:418 … future.Current != 5`; `schema.sql` user_version `0`; the only `.Replay(`/`OpenReplayBinder(` hit is the definition `binder.go:37`; `host/replay` non-test importers `0` (control `host/store` `58`) |
| V22 | [r2] A plan has at most one effect, by construction (kimi's fix), with a control | `/tmp/iter360-designer/l1-greps.txt`: `git show 77c9e3d4:host/coordinator/effectplan.go \| grep -n …`; the same over `effectful.go`, `design_docs/sketches/effectplan.ail` and `effectplan_drift_test.go`; `git grep -n "\.Effects = append\|Effects: \[\]PlannedEffect" 77c9e3d4 -- "host/coordinator/*.go" \| grep -v _test.go`; control: `git grep -n "L1/two-effects\|func TestProducerMintsNoPacketForMultiEffectPlan"`, and in the `77c9e3d4` copy `go test ./host/coordinator/ -run '^(TestProducerMintsNoPacketForMultiEffectPlan\|TestParsePlanLaws.*)$' -v` | `effectplan.go:32 maxPlanEffectsV1 = 1 // L1`, `:125 if !effectCountOk(n) {`, `:126 lawErr(LawL1, "%d effects; a v1 plan carries at most %d", …)`, `:323 func effectCountOk(n int) bool { return n >= 0 && n <= maxPlanEffectsV1 }`; `effectful.go:289` (live) and `:430` (replay) `parsePlan`, `:307` and `:469` iterate `pl.Effects`, `:246 default: // unreachable under L1`; sketch `effectplan.ail:38` "L1: a v1 plan carries at most ONE effect", `:40 export func effectCountOk`, `:41 ensures { result == (n >= 0 && n <= 1) }`; drift test `:35` binds `effectCountOk`; the only non-test construction of a plan's effects is `effectplan.go:135 plan.Effects = append(plan.Effects, e)` (the other hit, `effectful.go:213`, is the world block); control: `effectplan_test.go:63 "L1/two-effects"` and `decisionpacket_producer_test.go:316`, both `--- PASS` (each submits a two-effect plan and gets the L1 refusal) |
| V23 | [r2] Rollback is real; one binary; the CLI reads over HTTP | `git show 77c9e3d4:cmd/ailang-worldd/{service_upgrade,main,why,episode}.go \| grep -n …` | `service_upgrade.go:210 fresh, previous := …".new", …".prev"`, `:219 os.Link(e.paths.Binary, previous)`, `:240 return e.rollbackUpgrade(…)`, `:256 err = os.Rename(previous, e.paths.Binary)`, `:262 e.waitForBuild(r.Address, from)`, `:287 "rollback failed too: the service is stopped"`; `main.go:354 case "serve":`, `:395 case "call":`, `:398 case "why":`, `:434 case "episode":` (one `package main`); `why.go:136 r.getJSON("/v1/objects/"+…)`; `episode.go:330 if self.ExecutableSHA256 != installed {` |
| V24 | [r2] The r2 prototype passes the full plain gate | `VERIFY_GO_LEGS=plain AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_go.sh` → `/tmp/iter360-designer/r2-verify_go.txt`. The prototype is byte-identical to the r2 patch: no source edit came between this run and `git diff`, except the mutation run, whose restore V25 checks | `✓ go gate PASSED (plain): …`; `RC=0`; 35 packages `ok`, including `host/broker 315.957s`, `host/store 15.408s`, `host/daemon 347.692s`, `host/verifygate 311.481s`, `cmd/ailang-worldd 197.173s` |
| V25 | [r2] The mutation table on r2 | `/tmp/iter360-designer/mutate.py` (r2: 18 tests over `host/broker`, `host/coordinator`, `host/projection`, `cmd/ailang-worldd`, `host/store`; multi-edit mutants) → `/tmp/iter360-designer/r2-mutants.txt`; then `git diff` compared with its pre-run copy | `CONTROL: rc=0 pass=18 red=-`; every r1 mutant's reds as in V12, plus TV on `MUT-RECORD-DROPS-LINE` (`red=TC,TV,TX`) and `MUT-HEAD-SHAPE` (`red=TC,TP,TV,TX`); `MUT-NO-MARK`, `MUT-MARK-ERROR-IGNORED`, `MUT-MARK-AFTER-RECORD: red=TV`; `MUT-STORE-REJECTS-6`, `MUT-MARK-ANY-VERSION: red=TS`; `MUT-MARK-NOOP: red=TS,TV`; survivors `MUT-PUT-NO-CHECK`, `MUT-CLEANUP-UNWIRED`, `MUT-FUTURE-BOUND-5` (`rc=0 pass=18`); `tracked diff identical after mutants`; `-count=3` of TV, TX, TC, TP, TW, TS and `Future` → `ok` ×5 |
| V26 | [r2] The r2 patch applies at HEAD and on a fresh archive, and its 18 named tests pass there | `git add -N <5 new files>; git diff -- . ':(exclude)<patch>' > <patch>; git reset -q`; `git checkout --` the 17 modified files; `rm` the 5 new; `git status --porcelain`; `git apply --check <patch>` in the worktree; then in `/tmp/iter360-designer/applied2` (`git archive 77c9e3d4`): `git apply --check`, `git apply`, `go build ./...`, `go vet` (5 packages), `gofmt -l host cmd`, `go test -count=1 -v -run '^(<18 names>)$'` (5 packages) | patch 1538 lines, sha256 prefix `02259c3bdc1c942b`, 22 `diff --git` sections; numstat `store.go 45/5`, `broker.go 22/5`, `schema_version_test.go 71/3`, `deadline_guard_test.go 1/0`, the four test stores `4/0` each, `QUICKSTART.md 15/0`, and the r1 rows unchanged; `failure_test.go` 361 lines; the worktree is back to `?? <patch>` only; `apply-check-rc=0`; `APPLY-CHECK-ARCHIVE rc=0`, `APPLY-OK`, `BUILD-VET-OK`, `GOFMT-OK`; 18 `--- PASS` and 5 `ok` |
| — | **UNMEASURED** | — | **[r2]** a live `service upgrade` with its rollback onto a marked store (§3.7 reads the code path, V23; launchd was not driven); an older **CLI** against a daemon of this build (derived from V16 and `why.go:896–899`, not run); Linux (R6: GNU `cp` exit code and `EACCES` on the CI runner; the acceptance gate is the PR's CI); the home store's entries 54–64 for episode `worldd-2611031015` (row claim 8); [r3] the r3 delta to the patch (the mark-failure arm, TV's fault-injection rows and MUT-MARK-FAIL-SILENT) is specified here and pinned by the sprint plan, not by the r2 patch file; the r2 patch was gated by V24 (`verify_go.sh` on r2) and V15; a `-race` run |

## 11. Quorum log

**Round 1 (r1, `89c1967e`): BLOCKED, $0.2805.** Seats: gemini-4-argon **PASS** ($0.0706),
oc-kimi-k3 **REJECT** ($0.1767), gpt6-1-sol **REJECT** ($0.0332). claude was benched as the
author's vendor; no seat was absent. Artifact:
`~/.ailang/state/mission-quorum-world/w-exec-failure-carries-a-reason-2026-11-06T21-30-57Z.json`
(read through the controller's message).

- **kimi: "logged once per call" rested on one fixture** (V6, a `"finish":false` hit). Its fix:
  show by construction that the effects array has one element, with a control; otherwise rewrite
  Claim E and add a two-effect test. **r2's answer:** construction is shown. L1 (`effectplan.go:32`,
  `:125–126`) gates the only non-test builder of a plan's effects (`:135`). It runs on the live and
  replay paths and is a Z3-verified sketch contract bound by a drift test. The control is the two
  existing two-effect refusals, both PASS (V22). Claim E and §3.5 now cite it and say "at most once
  per call". They also state the helper's per-effect rule, so a future widened L1 logs once per
  failed effect. No two-effect TP fixture is added, because L1 makes that plan unconstructible.
- **gpt: R1's "World never downgrades" had no V-log row, and Q1 defaulted to accepting it.** Its
  fix: a compatibility gate with cited policy and enumerated combinations, or a redesign. **r2's
  answer:**
  - Measured: no written policy exists (V19), and rollback is a shipped feature (V23), so r1's
    assertion was false and is withdrawn.
  - New §3.7 shows that no carrier is tolerated by old readers, then adds the smallest enforced
    gate: the store's existing `user_version` refusal, raised lazily from 5 to 6 just before the
    first lined record. It is measured end to end against HEAD code (V20) and pinned by TS, TV and
    6 new mutants (V25).
  - §3.7's table enumerates daemon, CLI, replay and rollback combinations.
  - The two remaining policy points become Q1 and Q2, addressed to Mark with defaults, not
    asserted.
- **gemini (pass):** no change requested.

r2 changes the patch (r1 `b7dcb7821c3b389a`, 1184 lines → r2 `02259c3bdc1c942b`, 1538 lines). The
r1 carrier, classes, bounding rule and surfaces are unchanged. r2 adds the store marker, its
broker call and tests, and one QUICKSTART sentence block.

**Round 2 (r2, `dad6beae`): BLOCKED, $0.2615.** Seats unchanged: oc-kimi-k3 **REJECT** ($0.1118),
gpt6-1-sol **REJECT** ($0.0472), gemini-4-argon **REJECT** ($0.1026); claude benched; no seat absent.
Artifact: `~/.ailang/state/mission-quorum-world/w-exec-failure-carries-a-reason-2026-11-06T21-55-19Z.json`.
All three objections land on **one surface**, §3.7's mark-failure arm (r2: write the record without
the line), as a silent fallback. None disputes the design direction (the `failure` field, the
classes, the bounding rule, the lazy schema marker), and r1's two objections were not re-raised.

- **gpt6-1-sol, verbatim fix:** "If MarkFailureLines fails, the broker must not downgrade the record
  to the legacy shape or report a normally committed classified failure. It surfaces an explicit
  persistence/compatibility error through the broker's verified storage-error path. Because the
  handler may already have executed, this path must preserve debit/accounting invariants and must
  not automatically rerun the effect." Plus fault-injection assertions in TV.
- **gemini-4-argon, verbatim fix:** "if `MarkFailureLines` fails when a non-empty `reason` is present,
  fail the broker record-write path and return the store error (just as a failed `putRecord` does)
  rather than clearing `reason = ""` and continuing. Update §3.7, §4 Claim F, test TV, and
  `MUT-MARK-ERROR-IGNORED` ... Also clean up the stale r1 note in §10's UNMEASURED row."
- **oc-kimi-k3:** the same silent fallback; its fix keeps the line-free record but logs the
  suppression. Its objection is the silence, which the gpt/gemini fix removes, because no record is
  written and the error reaches the logged `Dispatch` error path.

**r3 — controller narrow-refinement carve-out (iteration 360).** Every remaining objection carries a
concrete reviewer-authored fix and none disputes the direction, so the controller applied the
gpt/gemini fix (two of three seats, the same text) verbatim to §3.7, Claim F, TV and the mutation
table, added MUT-MARK-FAIL-SILENT, and fixed the stale §10 note. Premise measured first-party by
the controller at `77c9e3d4`: a failed `putRecord` after `Execute` returns `*OutcomeWriteError`
(`broker.go:297–300`) with the grant already debited (`debitGrant`, `broker.go:260`, before
`Execute` at `:261`) and nothing re-run; the coordinator handles `*OutcomeWriteError`
(`effectful.go:184`), and a `Dispatch` error is logged (`mcp.go:131–133`). kimi's alternative
(keep a line-free record and log a suppression) was not applied: it would still commit a failure
without its line, which is what the other two seats rejected. The patch file stays at r2
(`02259c3bdc1c942b`); the r3 delta (the mark-failure arm, TV's fault-injection rows and the new
mutant) is pinned by the sprint plan. No re-quorum: the carve-out routes straight to the planner.
