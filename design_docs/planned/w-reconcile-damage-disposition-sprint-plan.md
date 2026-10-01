# Sprint plan — w-reconcile-damage-disposition (iteration 214)

**Authority:** approved revision-2 design [w-reconcile-damage-disposition.md](w-reconcile-damage-disposition.md), World queue row 125, following row 122. Its r1/r2 Quorum log answers were read, including V15–V22 and the newline, commit-output and construction-site refinements. **Base:** `4b87fdd`. **Planner lane:** `codex:gpt-6.1-sol`. **Scope:** two test-first milestones, ~0.25 day. **Landing:** M1 then M2, or one atomic commit; controller owns all git writes. The full prototype remains in this worktree, uncommitted.

Use the pinned binary for gates that exercise AILANG:

```sh
export AILANG_BIN=$HOME/.pinned-ailang/ailang
```

## 1. Problem

At HEAD, `committed()` verifies the row-122 hashes and bindings but wraps absent rows with a nil `%w`, returns untyped decode failures, and accepts hash-consistent `null` output as a successful reconcile. Integrity damage and other not-available dispatch failures are sanitized for the wire without an operator log. The approved design retains the existing client message, reports detail through the daemon-owned ErrorLog, and explicitly rejects quarantine.

This is host-boundary typing and daemon plumbing, not an AILANG package or kernel law. No `.ail`, store interface, store implementation, frozen driver, or public wire vocabulary changes.

## 2. The fix (M1 and M2, prototyped)

### M1 — Coordinator typing (test-first)

Files: `host/coordinator/coordinator.go`, `host/coordinator/errors.go`, `host/coordinator/coordinator_test.go`.

Add `IntegrityError.Kind` with `mismatch`, `absent`, `undecodable`; keep the zero-kind legacy message compatible. Error phrases are `does not match its reference`, `is absent`, `does not decode`. Split every store read's `err != nil` wrap from its `!ok` damage branch. Type verified record decode, output-ref parse, and non-object output failures. Add explicit mismatch kinds to all five existing checks. Do not wrap transient read errors in integrity types.

RED first: `TestReconcileDamageDisposition` has the exact design names:

| Subtest | Fixture and required assertion |
|---|---|
| `record_absent` | hide record; record/absent |
| `output_absent` | hide recorded output; output/absent |
| `world_absent` | hide recorded world; world/absent |
| `record_undecodable` | PutObject of hash-correct invalid record JSON; record/undecodable |
| `record_output_ref_unparsable` | PutObject of hash-correct record with bad output ref; record/undecodable |
| `output_not_object/null`, `output_not_object/[1]` | forged record → forged output → forged world; receipt points at the recalculated world ref; output/undecodable |
| `store_read_error_stays_untyped` | GetObject returns wrapped DeadlineExceeded; not IntegrityError, errors.Is survives |
| `commit_refuses_non_object_output/null`, `commit_refuses_non_object_output/[1]` | real forward Dispatch; OutputError, no receipt, head unchanged |

Every damage case checks invocation id, Object, Kind, no `%!`, no reconciled result or output bytes, fixture-hit count, must-not-run runner and unchanged head. Existing `TestReconcileRefusesDamagedRecord` rows additionally assert `Kind == "mismatch"`; its exact-result control stays green.

The first compile RED exposed the missing Kind field (zero RUN: **not a behavioral gate**). Adding only that field allowed the behavioral RED: six damage rows failed, including both output variants; `null` returned success at base. The deadline and forward-output controls were already green and their named mutations below make them red. GREEN after implementation: 27 RUN, zero failures/skips.

```sh
# Behavioral RED (Kind field scaffold only): rc=1, 19 RUN.
go test -count=1 -run '^TestReconcileDamageDisposition$' -v ./host/coordinator/
# GREEN: rc=0, 27 RUN.
go test -count=1 -run '^(TestReconcileDamageDisposition|TestReconcileRefusesDamagedRecord)$' -v ./host/coordinator/
```

M1 mutation list (all required rows assigned here):

| Mutation | File / change | Named test that must kill it | Verdict / RUN | Restored SHA-1 |
|---|---|---|---|---|
| MUT-ABSENT-REC | `host/coordinator/coordinator.go`: Record absent branch returns old untyped nil wrap | `'^TestReconcileDamageDisposition$/record_absent'` | **KILLED** / 3 | `e3d3f34d46ab5ae95482c3a93d14864326faf229` |
| MUT-ABSENT-OUT | `host/coordinator/coordinator.go`: Output absent branch returns old untyped nil wrap | `'^TestReconcileDamageDisposition$/output_absent'` | **KILLED** / 3 | `e3d3f34d46ab5ae95482c3a93d14864326faf229` |
| MUT-ABSENT-WORLD | `host/coordinator/coordinator.go`: World absent branch returns old untyped nil wrap | `'^TestReconcileDamageDisposition$/world_absent'` | **KILLED** / 3 | `e3d3f34d46ab5ae95482c3a93d14864326faf229` |
| MUT-REC-DECODE-UNTYPED | `host/coordinator/coordinator.go`: Verified record decode returns fmt.Errorf | `'^TestReconcileDamageDisposition$/record_undecodable'` | **KILLED** / 3 | `e3d3f34d46ab5ae95482c3a93d14864326faf229` |
| MUT-OUTREF-UNTYPED | `host/coordinator/coordinator.go`: Unparsable output ref returns fmt.Errorf | `'^TestReconcileDamageDisposition$/record_output_ref_unparsable'` | **KILLED** / 3 | `e3d3f34d46ab5ae95482c3a93d14864326faf229` |
| MUT-OUT-NIL-OBJ | `host/coordinator/coordinator.go`: Reconcile drops nil-map rejection | `'^TestReconcileDamageDisposition$/output_not_object/null'` | **KILLED** / 3 | `e3d3f34d46ab5ae95482c3a93d14864326faf229` |
| MUT-ERR-AS-DAMAGE | `host/coordinator/coordinator.go`: Store read error enters absent integrity branch | `'^TestReconcileDamageDisposition$/store_read_error_stays_untyped'` | **KILLED** / 3 | `e3d3f34d46ab5ae95482c3a93d14864326faf229` |
| MUT-KIND-SWAP | `host/coordinator/coordinator.go`: Record absent labelled mismatch | `'^TestReconcileDamageDisposition$/record_absent'` | **KILLED** / 3 | `e3d3f34d46ab5ae95482c3a93d14864326faf229` |
| MUT-COMMIT-NULL-OK | `host/coordinator/coordinator.go`: Forward parseOutput drops nil-map rejection | `'^TestReconcileDamageDisposition$/commit_refuses_non_object_output/null'` | **KILLED** / 3 | `e3d3f34d46ab5ae95482c3a93d14864326faf229` |

### M2 — Projection arm, ErrorLog and daemon wiring (test-first)

Files: `host/projection/projection.go`, `host/projection/projection_test.go`, `host/projection/a2a_wire_test.go`, `host/daemon/daemon.go`, `host/daemon/invoke_e2e_test.go`, `design_docs/planned/w-transition-invocation-coordinator.md`, `docs/QUICKSTART.md`.

Require `projection.Config.ErrorLog io.Writer` in New and retain it on Handler. Add the explicit IntegrityError arm returning `codeInternal, notAvailableMessage`; SourceError remains on default. Log precisely not-available dispatch failures with `fmt.Fprintf(h.errorLog, "ailang-worldd: a2a refusal: %s %s: %q\n", method, id, err.Error())`. Admission failure or unavailable coordinator uses `-` where no validated invocation was dispatched. Other refusal classes have no operator line. The production daemon passes its resolved `d.errLog`; both test-side config sites supply concrete writers. Update the parent's refusal row and the existing QUICKSTART stderr description.

RED first: `TestA2ADispatch/R13_reconcile_absent`; `TestA2ADispatchWire/R13_integrity` commits through the real coordinator and hides the record on resend; `R4_R6_default` checks its secret reaches the log; `store_error_newline` checks exactly one physical line plus escaped detail; `R15` and `R11` check empty logs; `TestProjection_ConfigValidation` strips ErrorLog; `TestA2AErrorLogWiring` constructs a real daemon with Config.ErrorLog and drives the production projection using a recorder. This last test deliberately needs no socket or interpreter, and kills io.Discard wiring.

The projection behavioral RED had 32 RUN: missing-log rows and nil-ErrorLog validation failed. The new wire-mapping row was already green via default; its explicit-arm mutation supplies its RED. Initial daemon compile attempts were uninformative while the new seam was incomplete; the unwiring mutation then gave its behavioral RED, and the restored wiring test gave GREEN. `assertWireRefusal` checks hidden strings before the constant, so MUT-LOG-TO-WIRE fires the specified sanitizer assertion. All pre-existing rows remain unchanged on the wire.

```sh
# Projection RED rc=1 and GREEN rc=0; each 32 RUN, zero skips.
go test -count=1 -run '^(TestA2ADispatchWire|TestProjection_ConfigValidation|TestA2ADispatch)$' -v ./host/projection/
# Wiring GREEN rc=0, 1 RUN; unwired RED rc=1, 1 RUN (mutation below).
go test -count=1 -run '^TestA2AErrorLogWiring$' -v ./host/daemon/
```

M2 mutation list (all remaining design rows assigned here):

| Mutation | File / change | Named test that must kill it | Verdict / RUN | Restored SHA-1 |
|---|---|---|---|---|
| MUT-WIRE-INTEGRITY | `host/projection/projection.go`: Explicit integrity arm returns R15 message | `'^TestA2ADispatch$/(R13_reconcile_absent\|R13_reconcile_integrity)'` | **KILLED** / 3 | `6acbbfb21910fcf00c7fd0187a5904f6894500be` |
| MUT-NOLOG | `host/projection/projection.go`: Delete operator write | `'^TestA2ADispatchWire$/(R13_integrity\|R4_R6_default)'` | **KILLED** / 3 | `6acbbfb21910fcf00c7fd0187a5904f6894500be` |
| MUT-LOG-ALL | `host/projection/projection.go`: Log every dispatch class | `'^TestA2ADispatchWire$/R15'` | **KILLED** / 2 | `6acbbfb21910fcf00c7fd0187a5904f6894500be` |
| MUT-LOG-TO-WIRE | `host/projection/projection.go`: Interpolate cause into wire message | `'^TestA2ADispatchWire$/R13_integrity'` | **KILLED** / 2 | `6acbbfb21910fcf00c7fd0187a5904f6894500be` |
| MUT-LOG-NEWLINE | `host/projection/projection.go`: Log cause with %v instead of %q | `'^TestA2ADispatchWire$/store_error_newline'` | **KILLED** / 2 | `6acbbfb21910fcf00c7fd0187a5904f6894500be` |
| MUT-LOG-DEADLINE | `host/projection/projection.go`: Also log R11 deadline class | `'^TestA2ADispatchWire$/R11'` | **KILLED** / 2 | `6acbbfb21910fcf00c7fd0187a5904f6894500be` |
| MUT-ERRLOG-NIL-OK | `host/projection/projection.go`: Accept nil ErrorLog | `'^TestProjection_ConfigValidation$'` | **KILLED** / 1 | `6acbbfb21910fcf00c7fd0187a5904f6894500be` |
| MUT-ERRLOG-UNWIRED | `host/daemon/daemon.go`: Daemon passes io.Discard | `'^TestA2AErrorLogWiring$'` | **KILLED** / 1 | `2b66316b92c94c141dbd38c0c21870d66c7cd7a1` |

## 3. Acceptance criteria

| AC | Discharge |
|---|---|
| AC1 | Six named damage fixtures refuse with Object/Kind/id and no execution/result/head movement; existing mismatch controls remain green. |
| AC2 | Wrapped deadline remains errors.Is-compatible and untyped; all M1 refusal text rejects `%!`. |
| AC3 | Existing dispatch and wire rows pass; R13 absent and mismatch keep -32603 / notAvailableMessage; explicit arm independently mutated. |
| AC4 | One escaped physical line for not-available causes, R15/R11 empty logs, sanitized body and real daemon sink; announce is a separate existing writer and never injected here. |
| AC5 | Production Quarantine caller census unchanged: only startup abort. Coordinator Store seam diff is empty. |
| AC6 | vet, formatting, pinned verify_ail and coordinator/projection race pass; broader host race has sandbox and unrelated failures, explicitly OPEN below. |

## 4. Mutation table (run on the prototype)

**17/17 KILLED; 0 SURVIVED.** Every row above was applied alone, its exact named subtest run with `-count=1 -v`, then the production file restored from its saved bytes. Each run had nonzero discovery and an assertion failure (never merely compile failure). Before and after each mutation, `shasum <file>` returned the SHA-1 in the table. End-of-harness observation: `TREE-IDENTICAL-AFTER-HARNESS`. No mutation survives in the worktree. The sanitizer assertion-order refinement was followed by a second isolated MUT-LOG-TO-WIRE kill and identical shasum.

Exact mutation commands, in execution order (ordinary RE2 pipes, never escaped alternation):

```sh
# MUT-ABSENT-REC: rc=1, 3 RUN, KILLED; restored byte-identical
go test -count=1 -run '^TestReconcileDamageDisposition$/record_absent' -v ./host/coordinator/
# MUT-ABSENT-OUT: rc=1, 3 RUN, KILLED; restored byte-identical
go test -count=1 -run '^TestReconcileDamageDisposition$/output_absent' -v ./host/coordinator/
# MUT-ABSENT-WORLD: rc=1, 3 RUN, KILLED; restored byte-identical
go test -count=1 -run '^TestReconcileDamageDisposition$/world_absent' -v ./host/coordinator/
# MUT-REC-DECODE-UNTYPED: rc=1, 3 RUN, KILLED; restored byte-identical
go test -count=1 -run '^TestReconcileDamageDisposition$/record_undecodable' -v ./host/coordinator/
# MUT-OUTREF-UNTYPED: rc=1, 3 RUN, KILLED; restored byte-identical
go test -count=1 -run '^TestReconcileDamageDisposition$/record_output_ref_unparsable' -v ./host/coordinator/
# MUT-OUT-NIL-OBJ: rc=1, 3 RUN, KILLED; restored byte-identical
go test -count=1 -run '^TestReconcileDamageDisposition$/output_not_object/null' -v ./host/coordinator/
# MUT-ERR-AS-DAMAGE: rc=1, 3 RUN, KILLED; restored byte-identical
go test -count=1 -run '^TestReconcileDamageDisposition$/store_read_error_stays_untyped' -v ./host/coordinator/
# MUT-KIND-SWAP: rc=1, 3 RUN, KILLED; restored byte-identical
go test -count=1 -run '^TestReconcileDamageDisposition$/record_absent' -v ./host/coordinator/
# MUT-COMMIT-NULL-OK: rc=1, 3 RUN, KILLED; restored byte-identical
go test -count=1 -run '^TestReconcileDamageDisposition$/commit_refuses_non_object_output/null' -v ./host/coordinator/
# MUT-WIRE-INTEGRITY: rc=1, 3 RUN, KILLED; restored byte-identical
go test -count=1 -run '^TestA2ADispatch$/(R13_reconcile_absent|R13_reconcile_integrity)' -v ./host/projection/
# MUT-NOLOG: rc=1, 3 RUN, KILLED; restored byte-identical
go test -count=1 -run '^TestA2ADispatchWire$/(R13_integrity|R4_R6_default)' -v ./host/projection/
# MUT-LOG-ALL: rc=1, 2 RUN, KILLED; restored byte-identical
go test -count=1 -run '^TestA2ADispatchWire$/R15' -v ./host/projection/
# MUT-LOG-TO-WIRE: rc=1, 2 RUN, KILLED; restored byte-identical
go test -count=1 -run '^TestA2ADispatchWire$/R13_integrity' -v ./host/projection/
# MUT-LOG-NEWLINE: rc=1, 2 RUN, KILLED; restored byte-identical
go test -count=1 -run '^TestA2ADispatchWire$/store_error_newline' -v ./host/projection/
# MUT-LOG-DEADLINE: rc=1, 2 RUN, KILLED; restored byte-identical
go test -count=1 -run '^TestA2ADispatchWire$/R11' -v ./host/projection/
# MUT-ERRLOG-NIL-OK: rc=1, 1 RUN, KILLED; restored byte-identical
go test -count=1 -run '^TestProjection_ConfigValidation$' -v ./host/projection/
# MUT-ERRLOG-UNWIRED: rc=1, 1 RUN, KILLED; restored byte-identical
go test -count=1 -run '^TestA2AErrorLogWiring$' -v ./host/daemon/
```

Assertion log excerpts (durable evidence; transient harness files are removed):

```text
MUT-ABSENT-REC: coordinator_test.go:1072: damage: *fmt.wrapError coordinator: reconcile a2a:ep1:t1: record: ok=false: %!w(<nil>) want record/absent
MUT-ABSENT-OUT: coordinator_test.go:1072: damage: *fmt.wrapError coordinator: reconcile a2a:ep1:t1: output: ok=false: %!w(<nil>) want output/absent
MUT-ABSENT-WORLD: coordinator_test.go:1072: damage: *fmt.wrapError coordinator: reconcile a2a:ep1:t1: world: ok=false: %!w(<nil>) want world/absent
MUT-REC-DECODE-UNTYPED: coordinator_test.go:1072: damage: *fmt.wrapError decode record: invalid character 'i' looking for beginning of value want record/undecodable
MUT-OUTREF-UNTYPED: coordinator_test.go:1072: damage: *fmt.wrapError output ref: hashref: missing algorithm tag (expected algo:digest): "bad" want record/undecodable
MUT-OUT-NIL-OBJ: coordinator_test.go:1072: damage: <nil> <nil> want output/undecodable
MUT-ERR-AS-DAMAGE: coordinator_test.go:1069: read error: *coordinator.IntegrityError coordinator: reconcile a2a:ep1:t1: recorded record is absent
MUT-KIND-SWAP: coordinator_test.go:1072: damage: *coordinator.IntegrityError coordinator: reconcile a2a:ep1:t1: recorded record does not match its reference want record/absent
MUT-COMMIT-NULL-OK: coordinator_test.go:986: output accepted: {Output:map[] OutputBytes:[110 117 108 108] InvocationID:a2a:ep1:t1 WorldRef:sha256:08aff506beecaec2ac11b686f1def1513547d79ad62b1f4f6c829363151945d8 EntryIndex:1 RecordRef:sha256:2a1a58f77deebe031f88444967c55ce653d8f5d2a61becd494b28fd51a6fd29d Reconciled:false} <nil>
MUT-WIRE-INTEGRITY: projection_test.go:838: got -32603 "invocation was not committed; send a new task id", want -32603 "transition invocation is not available in this daemon"
MUT-NOLOG: a2a_wire_test.go:212: operator log="" want one line carrying "is absent"
MUT-LOG-ALL: a2a_wire_test.go:236: R15 logged: "ailang-worldd: a2a refusal: tasks/send a2a:ep-a:t1: \"coordinator: invocation a2a:ep-a:t1 was not committed; send a new task id\"\n"
MUT-LOG-TO-WIRE: a2a_wire_test.go:208: body interpolates inner error detail "coordinator:": {"error":{"code":-32603,"message":"coordinator: reconcile a2a:ep-a:t1: recorded record is absent"},"id":"retry","jsonrpc":"2.0"}
MUT-LOG-NEWLINE: a2a_wire_test.go:226: operator log="ailang-worldd: a2a refusal: tasks/send a2a:ep-a:t1: coordinator: load source: SECRET\nsecond line\n" want one line carrying "SECRET\\nsecond line"
MUT-LOG-DEADLINE: a2a_wire_test.go:279: deadline logged: "ailang-worldd: a2a refusal: tasks/send a2a:ep-a:t-r11: \"coordinator: execute: context deadline exceeded\"\n"
MUT-ERRLOG-NIL-OK: projection_test.go:1414: New with nil ErrorLog: want an error (the projection never invents its own seams)
MUT-ERRLOG-UNWIRED: invoke_e2e_test.go:398: daemon ErrorLog=""
```

## 5. Gates and prototype inventory

| Gate | Observed |
|---|---|
| `gofmt -l host cmd` | empty; zero lines |
| `git diff --check` | rc=0, empty |
| `go vet ./...` | rc=0, empty |
| `go build ./...` | rc=0 |
| `AILANG_BIN=$HOME/.pinned-ailang/ailang go test -race -count=1 -v ./host/coordinator/ ./host/projection/` | rc=0; 84 coordinator + 85 projection = 169 RUN; zero skips/failures |
| `go test -count=1 -run '^TestA2AErrorLogWiring$' -v ./host/daemon/` | rc=0; 1 RUN, zero skips |
| `AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_ail.sh` | rc=0; 16 identities, 40 named tests, 9/9 package steps, PUB011 uncontracted_exports=3 |
| `AILANG_BIN=$HOME/.pinned-ailang/ailang go test -race -count=1 -v ./host/...` | rc=1; 1361 RUN before package panics; socket arms UNINFORMATIVE UNDER SANDBOX; separate timing failures detailed in OPEN. |

Initial broader host run omitted AILANG_BIN and failed mandatory archive/broker/pkgproj interpreter checks. It is superseded by the pinned run, not counted as a product failure. Initial targeted race passed 169 RUN with three interpreter skips; the final pinned run above has zero skips.

Measured production diff (`git diff --stat` / `git diff --numstat`, excluding tests/docs):

| Milestone | Files | Add / delete | Changed / net lines |
|---|---|---|---|
| M1 | coordinator.go +24/−15; errors.go +9/−1 | +33/−16 | 49 changed / +17 net |
| M2 | projection.go +25/−0; daemon.go +1/−0 | +26/−0 | 26 changed / +26 net |

Both milestones are below 150 production changed lines. Tests and usage/refusal docs are included in the prototype. No git writes, message operations, deployment or publication performed.

## OPEN

- **OPEN-1 — broad gate outside sandbox.** AC6's all-host race gate cannot be certified here: loopback binds in broker/daemon return `operation not permitted`. Those socket results are **UNINFORMATIVE UNDER SANDBOX**, neither PASS nor product FAIL. Controller must rerun outside sandbox. No all-repository `go test ./...` gate is claimed; it also remains for controller certification.
- **OPEN-2 — unrelated timing failures observed.** The pinned broad run reported `host/broker/TestModelHandlerTimeoutWritesFailureRecord`: elapsed 5.514203291s, expected ≤2s for a 40ms bound; and `host/daemon/TestCommitInvocationReceiptIdentity/A_landed_uncertain_then_B_on_top`: Timeout instead of CommitUncertain under its 100ms deadline. These are real recorded assertion failures, not attributed to socket permissions or asserted fixed. The daemon race suite was already measured flaky at base by the controller; no attempt to chase it or alter those tests was made.
- Design policy decisions are settled: no dedicated wire message, no quarantine. Design residuals R-125-0 through R-125-5 (self-consistent journal forgery, store-level typing, forward-head trust, receipt typing, card-route logs, repair tooling) remain outside this sprint, with their owners in the approved design.

## Verification Log

All source premises below were re-checked against HEAD `4b87fdd` (read-only git show/git grep), not assumed from the design's older `8d86c4f`. Behavioral observations are explicitly labelled scaffold/prototype, not falsely attributed to pristine HEAD.

| V | Premise | Command → observed |
|---|---|---|
| V1 | base, queue and binding instructions | `git rev-parse --short HEAD`; `cat CLAUDE.md design_docs/coding-standards.md design_docs/planned/w-reconcile-damage-disposition.md design_docs/planned/w-invocation-reconcile-object-integrity-sprint-plan.md sprint_w-world-core-0-1-1-release.json`; `rg -n '^125\.' design_docs/world-mission.md` → 4b87fdd, approved design AC1–AC6/M1/M2/17 MUT rows and r1/r2 answers read; row 125 at 1913. Charter start and DESIGN §§1/14 also read. |
| V2 | base IntegrityError fields and SourceError Kind precedent | `git show HEAD:host/coordinator/errors.go`; `git grep -n -E 'SourceError\|IntegrityError' HEAD -- host/coordinator/errors.go` → IntegrityError at 92 has InvocationID/Object only; SourceError at 15 has Kind. |
| V3 | absent/decode paths and single reconcile entry | `git show HEAD:host/coordinator/coordinator.go` → absent nil wraps 176/194/201; untyped decode 183, output ref 190, output decode 214; sole committed call 243. |
| V4 | row-122 checks precede result | `git show HEAD:host/coordinator/coordinator.go` → rehash record/output, bind invocation, rederive world, bind state root before Result; final output decode at 213 lacks nil-map check. |
| V5 | forward non-object rejection is sole output writer path | `git grep -n -E 'json.Unmarshal\(out\|parseOutput\(res\|pl := planInvocation' HEAD -- host/coordinator/coordinator.go`; `rg -n 'parseOutput\(\|planInvocation\(\|OutputV1' host cmd -g '*.go' -g '!**/*_test.go'` → parseOutput nil predicate at HEAD 154, call 308, error return before planning 312; planInvocation sole output construction, plan.go 84 and objects list 103. |
| V6 | wire mapping and ctx precedence | `git show HEAD:host/projection/projection.go` (391–437) → no IntegrityError or SourceError arm; ctx arm 427 returns R11 deadline, default 434 returns constant notAvailableMessage. |
| V7 | no operator log; Fail has separate role | `git grep -n -E 'log\.\|slog\|Logger\|Printf' HEAD -- host/projection/projection.go host/coordinator/coordinator.go host/coordinator/errors.go` → zero hits; `git grep -n -E 'errLog\|Fail renders' HEAD -- host/daemon/handlers.go host/projection/projection.go` → daemon sanitized 500 detail at handlers.go 178; Fail is card 503/504 renderer at projection.go 136. |
| V8 | required-seam validation and complete config-site census | `git show HEAD:host/projection/projection.go` (167–186); `git show HEAD:host/projection/projection_test.go` (1397–1414); `git grep -n -E 'projection.New\(\|projection.Config\{' HEAD -- host cmd` → New returns Handler/error and validates five nil seams; strip table pins them; one production site daemon.go 580 and one qualified test site invoke_e2e_test.go 244; in-package testConfig at projection_test.go 212 feeds wire rig. |
| V9 | real wire rig accepts Store and pins precedent | `git grep -n -E 'R6_corrupt_source\|R13_reconcile_integrity\|R4_R6_default\|func.*handler\(t.*coordinator.Store' HEAD -- host/projection/projection_test.go host/projection/a2a_wire_test.go` → Store seam at wire-test 71; R6/R13 constant rows 822/823; default wire row 269. |
| V10 | quarantine lifecycle, unchanged caller census and no deletion | `git grep -n -E 'Quarantine\(\|DELETE FROM\|quarantined atomic.Bool\|refuses new operations' HEAD -- host/daemon/daemon.go host/store/store.go`; `rg -n 'DELETE FROM\|Quarantine\(' host cmd -g '*.go' -g '!**/*_test.go'` → bool 230, lifecycle doc 249, sole actual production caller startup abort (HEAD 679; prototype 680); only DELETE is session_credentials at 1370. checkQuarantine hits are checks, not new callers. |
| V11 | health ignores store and startup damage does not stop serving | `sed -n '743,750p' host/daemon/daemon.go`; `sed -n '932,946p' host/daemon/daemon.go` → static ok/version/path/interpreter health; integrity report lines asynchronously written to announce, then Serve runs. This is reporting/serving precedent, not an ErrorLog claim. |
| V12 | resolved rows transaction and echoing refs | `rg -n 'beginDurable\(ctx, "commit"\)\|finishDurable\|objects :=' host/store/store.go host/coordinator/plan.go`; `git show HEAD:host/store/store.go` (600–608, 668–682) → commit transaction 1118–1240; object list 103; GetObject echoes Hash/ref and GetWorld echoes Ref/ref. |
| V13 | store parse/I/O errors remain indistinguishable | same `git show HEAD:host/store/store.go` ranges → get-object/interface-hash and get-world/state-root/log-head errors are untyped fmt.Errorf wraps. |
| V14 | null/array/object behavior, with real chain evidence | behavioral M1 RED command in §2 (Kind scaffold only) → null forged chain returned err=nil; array decode untyped; six damage rows refused with wrong types. M1 GREEN control `{}` plus row-122 exact-result control reconcile; null/[1] now typed, forward null/[1] OutputError. No standalone old probe is claimed re-run. |
| V15 | no Store-seam change and operator writer distinct from announce | `git diff HEAD -- host/coordinator/coordinator.go` → only committed body changes; no interface diff; `sed -n '225,240p' host/daemon/daemon.go` and `rg -n 'announce\|Announce' host/daemon/daemon.go` → ErrorLog resolved to process stderr separately; Run takes announce independently. Wiring test receives exactly one refusal line in Config.ErrorLog. |
| V16 | pinned interpreter | `$HOME/.pinned-ailang/ailang --version` → AILANG v0.41.0, commit 24ee108 (released, not dirty). |
| V17 | non-vacuity/restoration/gates | isolated commands in §4 → 17 named mutants KILLED with nonzero RUN, same before/after shasum, TREE-IDENTICAL-AFTER-HARNESS; final pinned race → 169 RUN, zero skips; wiring → 1 RUN. Gate results in §5 are measured, with broader failures left OPEN. |
