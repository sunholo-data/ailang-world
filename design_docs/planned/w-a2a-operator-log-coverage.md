# w-a2a-operator-log-coverage — Missing projection operator lines (row 127)

**Status:** PARKED needs-human-review (iteration 215, D-WORLD-47); two quorum rounds BLOCKED. No implementation or implementation acceptance.  
**Target:** next World patch; **Priority:** clause-6 residual, inherited groom position 2.  
**Estimated:** 0.2d implementation + 0.2d verification. **Dependencies:** landed row125.  
**Planner-Lane:** codex-ok. **Measured base:** `4e2600e601711f40dae42479167251d43a8a3bf6`, 2026-10-01.

## Problem and scope

The same failure pattern crosses both projection routes: an operator line is present for A2A
registry and missing-coordinator failures, absent for A2A resolver failures and card read
failures, and two existing lines lack projection-level assertions (V1–V5). Complete this
small host-boundary observability surface without changing refusal, authority, registry,
receipt, quarantine, or deadline semantics. D-WORLD-46 default B permits this residual; it
hardens clause6 but does not satisfy the parked MCP end state (V0).

Related: [row125 disposition](w-reconcile-damage-disposition.md), specifically R-125-3/4;
[row123 wire plan](w-a2a-invocation-wire-coverage-sprint-plan.md). These cover dispatch
refusals; this note covers admission failures and the card route. Related-doc search returned
highest neural score 0.42 (row123 plan); no duplicate threshold fired. Row125's canonical
landed design remains in planned/, as identified by the charter (V0).

## Decisions and exact logging contract

Reuse `Config.ErrorLog` and the row125 A2A helper/format (V1). Card failures use
the reviewer-supplied accurate card prefix, with a single inline Fprintf in writeUnavailable:

`ailang-worldd: a2a refusal: <constant label> -: <Go %q of err.Error()>\n`

`fmt.Fprintf(h.errorLog, "ailang-worldd: card refusal: GET /.well-known/agent.json -: %q\n", err.Error())`

The final `\n` is one actual newline; embedded newlines, carriage returns, quotes and control
characters in the cause are escaped by `%q`. Assert full bytes using `fmt.Sprintf` with a
literal expected format and explicit injected error text, rather than invoking the product
formatter as the oracle. Labels and id are fixed on these pre-dispatch paths: retained req.Method arguments are
`tasks/send` by the preceding method guard; they do not use the decoded request id (V11):

| Path | Label | Placement / cardinality |
|---|---|---|
| A2A resolver returns error | `a2a/resolve` | immediately before existing A2AError; one line |
| A2A present-head snapshot or head-check error | `tasks/send` | retain existing allowedDescriptors error branch; one line |
| A2A admitted skill, nil coordinator | `tasks/send` | retain existing nil-coordinator branch; one line |
| Card resolver/head/snapshot error, 503 or 504 | `GET /.well-known/agent.json` | once at entry to writeUnavailable, before status classification |

Do not parse the A2A request to obtain a method on resolver failure. Do not add card log calls
at AgentCard call sites: the complete caller map contains only AgentCard resolver and descriptor failures, each
followed by return, so one inline card log in writeUnavailable cannot double-log (V10).
Keep the existing dispatch log format and invocation id unchanged. No new direct logging of
Authorization, session binding, request headers/body, task data, or registry payloads.

**Privacy boundary:** raw error cause remains quoted on the operator-only sink, as decided by
row125. “No secrets/no raw errors” applies to response bodies: injected diagnostic sentinels,
credential, payload, and store detail must never reach either route's wire. Quoting is line
sanitation, not redaction. Existing arbitrary cause text and caller-derived dispatch id can
contain sensitive text (V1/V6); the unquoted method/id positions also permit a caller task-id newline to create extra physical
lines on existing dispatch logs (V1/V6). Redesigning that private logging policy is a separate residual,
not silently folded into row127. New pre-dispatch lines use fixed labels and `-` only.

The existing `assertRefusalLog` helper (a2a_wire_test.go494) checks prefix/detail substrings
and newline count (V4). New tests keep it unchanged and use full-byte assertions because
constant labels, exact escaped cause, and the distinct card prefix need a stronger oracle.

**Wire stays exact (source-verified V13–V15; no new resolver HTTP probe executed).** Resolver A2A errors keep HTTP200, null id, -32603 and
`transition invocation is not available in this daemon`. Snapshot/nil-coordinator errors keep
the decoded request id and same code/message. Card failures keep injected Fail envelopes:
503 `ProjectionUnavailable` / `the projection read failed; retry the request`; wrapped
DeadlineExceeded remains 504 `ProjectionDeadlineExceeded` /
`the projection read exceeded its bounded deadline`. Card cancellation retains current 503
classification. No new Connection header or cancellation behavior (V2/V3).

**Deferred:** row127(4) receipt typing stays with row128's trust-boundary design. GetReceipt
errors remain wrapped with `%w`, including deadline identity (V7). This task neither
classifies I/O failures as integrity damage nor changes journal authentication. Parent-row
quarantine and wire decisions are reused. All decisions here are agent-resolvable; no human
freeze item. Unattended quorum is mandatory; controller owns independent author-excluding
review and must record its verdict before sprint routing.

## Implementation and verification

Why host code rather than a package: this is transport error rendering and an injected local
operator I/O sink, at the effect boundary (S2/S3); no World law or kernel export changes.

- `host/projection/projection.go` — one A2A helper call at resolver error and one inline card Fprintf at writeUnavailable entry (~2 lines).
- `host/projection/projection_test.go` — named `TestProjection_OperatorLogCoverage` table and deterministic counted error resolver/head/reader fixtures (~150–220 LOC).
- `docs/QUICKSTART.md` — short explanation of newly covered route labels, operator-only quoted diagnostics and unchanged retry/refusal behavior (~5–10 LOC).

Projection tests use bytes.Buffer through Config.ErrorLog; every injected resolver/head/reader
failure asserts that seam count exactly1, full log bytes, exactly1 physical newline, and exact wire status/code/message/
id or card envelope. Nil-coordinator is a configuration stimulus, not an injected error seam:
prove registry admission was reached and its descriptor was selected. Clean absent-head
controls count head and snapshot reads separately (one each), without an injected error-count
claim. Fixtures: resolver error with untrusted body method and credential;
head-check error; present-head snapshot error (force head present to avoid the legitimate
absent-head path); admitted descriptor + nil coordinator; card resolver/head/snapshot each
with ordinary error and wrapped DeadlineExceeded; card wrapped Canceled stays503. Error
text contains newline, carriage return, quote and a private diagnostic sentinel. Wire
assertions forbid that sentinel and credential. Use deterministic returned context errors,
not short wall-clock deadlines, for classification.

Negative controls on the same configured sink: successful card, clean absent-head zero-skills
card, all four session denial kinds on both routes, invalid JSON/params, unlisted skill and
successful dispatch write zero lines. Pair each zero-log fixture with a logging failure
positive control so a disconnected sink cannot pass. The absent-head arm's snapshot failure
is intentionally suppressed by current allowedDescriptors semantics (V3); pin its zero log
rather than treating it as an unavailable response.

### Load-bearing mutations (implementation-stage, NOT executed by this design-only task)

For each row: pristine focused suite green; compiling mutant landed and hash changed; named
assertions red; inverse run excluding ALL downstream killer groups green; byte-exact restore.
Record the complete failing-subtest set and exact inverse selection. Shared-helper mutations
may kill several existing and new groups; do not claim a sole killer or exclude just one.
For row-specific mutants below, expand exclusions if additional downstream killers fire.
Keep mutation work in an isolated checkout, never mutate shared harness or running sources.

| Mutant | Named killer / assertion | Negative control |
|---|---|---|
| delete A2A resolver log | resolver_error: exact expected line | wire status/id remains correct; exclude all resolver-error arms; remaining suite green |
| delete existing registry log | a2a_snapshot_error: exact expected line | counted snapshot reached1; exclude all A2A registry-failure arms |
| delete existing nil-coordinator log | a2a_nil_coordinator: exact expected line | admitted descriptor proves nil branch; exclude all projection and daemon nil-coordinator logging arms |
| delete card centralized log | card ordinary + deadline arms: exact expected line | unchanged envelopes; exclude ALL card failure arms (resolver/head/snapshot, ordinary/deadline/canceled) |
| add duplicate call at one AgentCard caller | card_resolver_error: newline count1/full bytes | exclude all card-resolver-error variants that take this caller |
| change %q to %s | error_with_controls: full escaped bytes and newline count1 | exclude every control-character/quote log assertion, including existing store_error_newline and any affected integrity/corrupt-source arms; benign exact-%q arms may also kill %s and must be included |
| use attacker method at resolver logging | resolver_error: constant a2a/resolve label | body parsing observer remains0, malicious body method absent from log |
| log on ordinary denial/success | matching zero-log arm: log.Len()==0 | paired error control still yields1 |
| expose err.Error() on either wire | resolver/card sentinel exclusion and exact envelope | operator cause still logged; retain operator assertions when excluding wire killer |
| map wrapped deadline to503 | card_wrapped_deadline: status504/exact envelope | ordinary error stays503 |

The inverse is evidence for attribution, not a substitute for the counted reached seam. A
compile failure is not a kill. Add distinct mutations for head versus snapshot branches when
implementation placement can separate them.

Success: every above failure writes exactly1 line and all controls write0; all stated wire
values unchanged; all mutations killed with inverse controls; documentation updated;
`gofmt`, `go vet ./...`, `go build ./...`, full `go test ./...`, projection/daemon `-race`, and
`./scripts/verify_ail.sh` pass under pinned released binary, with complete out-of-sandbox
local runs and remote CI. Planner must establish the pristine baseline for its EXACT command
before attributing a red to this change; an existing red cannot be waved through. Full local
verify_go green is NOT established: the controller's pristine run was red in capsule, pkgproj,
replay and verifygate (V16). Those failures are recorded without attributing them to row126
or to load absent a comparator. No new .ail snippets or language claims. No product changes or mutation results are claimed in this note.

## Verification log

Commands executed in iteration215 worktree; source output read before making claims.

- **V0:** `git rev-parse HEAD` → `4e2600e601711f40dae42479167251d43a8a3bf6`; `sed -n '1910,1922p' design_docs/world-mission.md` → rows125/127/128 with landed parent path, missing paths, deferred trust design; `rg -n 'D-WORLD-46' design_docs/world-mission.md` → default B clause-residual order, clauses4/5/6 parked in latest STATUS

- **V1:** `rg -n 'logRefusal|ErrorLog|err.Error' host/projection/projection.go` → Config sink139/140, constructor required186, calls319/326/359, helper466/467: `fmt.Fprintf(h.errorLog, "ailang-worldd: a2a refusal: %s %s: %q\n", method, id, err.Error())`

- **V2:** `sed -n '275,282p;457,468p' host/projection/projection.go` → resolver error calls only A2AError with nil id; writeUnavailable only calls Fail, deadline504 else503. Positive control same read includes logRefusal and Fprintf, proving logging search/read live

- **V3:** `cat host/projection/projection.go` → AgentCard230–243 both error branches delegate writeUnavailable; A2A317–329 snapshot/nil branches log then error; allowedDescriptors388–403 suppresses snapshot error only if !hasHead; dispatchError context arm precedes default

- **V4:** `rg -n 'assertRefusalLog|bytes.Count|log.String|ErrorLog' host/projection/projection_test.go` → only ErrorLog io.Discard218 and nil-config strip1408; no log-byte assertion in this file. Same-command positive query matched both sink sites; separate `rg -n 'assertRefusalLog|ErrorLog' host/projection/*test.go` matches wire helper494 and callers212/226/349, proving assertion-query instrument health. `sed -n '470,515p' host/projection/a2a_wire_test.go` → helper494–499 checks newline count, prefix, detail

- **V5:** `sed -n '380,402p' host/daemon/invoke_e2e_test.go` → TestA2AErrorLogWiring seeds descriptor, sends with nil daemon coordinator, checks not-available response and exactly1 newline with tasks/send - prefix; `rg -n 'ErrorLog' host/daemon/daemon.go` → required injection at586 `ErrorLog: d.errLog`, resolved daemon sink527; Config docs227/239

- **V6:** `rg -n 'func InvocationID' host/coordinator/*.go` → coordinator.go101 concatenates `"a2a:" + episodeID + ":" + taskID`; `sed -n '205,240p' host/projection/a2a_wire_test.go` → store_error_newline injects SECRET newline and requires quoted cause operator-side, while assertWireRefusal forbids SECRET on wire. Existing tests require this policy

- **V7:** `sed -n '235,260p' host/coordinator/coordinator.go` → GetReceipt245, err wrapping247 `coordinator: receipt: %w`; no typed replacement in that branch. `sed -n '298,304p' design_docs/planned/w-reconcile-damage-disposition.md` → R-125-3 receipt stays untyped and R-125-4 card follow-up

- **V8:** `sed -n '610,645p' host/projection/projection_test.go` → TestAgentCard_HeadReadFailuresAre5xx checks head failure503/envelope and present-head snapshot503; fixtures165–194 expose fixedHeads/errReader. No logging behavior claimed from these existing tests

- **V9:** `go test -count=1 -run '^(TestAgentCard_HeadReadFailuresAre5xx|TestProjection_ConfigValidation|TestProjection_DeniedNeverTouchesRegistry|TestA2ADispatchWire)$' -v ./host/projection/` → rc0, 4 named parent tests, all listed children PASS; package2.715s. This is baseline evidence, not proof of future coverage

- **V10:** `rg -n 'writeUnavailable|logRefusal' host/projection --glob '*.go'` → writeUnavailable call sites ONLY232/241, definition457/comment453; logRefusal calls319/326/359, definition466/comment465. Same query positive controls include all existing log calls. `sed -n '20,90p;225,245p;275,330p;355,365p;453,468p' host/projection/projection.go` →232 resolver failure and241 allowedDescriptors failure both return immediately; denied branch235–237 calls Deny instead, successful branch builds card. This is the complete in-package source caller map, not inference from response output.

- **V11:** same source read → method guard302–305 rejects anything except tasks/send before319/326; those calls pass req.Method and literal "-". Dispatch359 passes req.Method and InvocationID(episode, params.ID), while A2AError320/327 separately passes req.ID. Fixed pre-dispatch log labels follow control flow; caller-derived newline risk is in dispatch id, not decoded JSON-RPC request id.

- **V12:** `sed -n '30,90p' host/authority/resolver.go` → enum35–43: Absent, Malformed, Unknown, Expired; String maps SessionAbsent, InvalidSession, SessionUnknown, SessionExpired. `sed -n '78,118p' host/daemon/middleware.go` → writeSessionDenial83: Absent/Unknown/Expired401, Malformed400, constant messages107–114. `sed -n '580,590p' host/daemon/daemon.go` → Deny injected as writeSessionDenial584, Fail writeAPIError585. Source read in V10 shows all four A2A switch arms281–291 and card delegation235–237; both denial routes bypass unavailable logging.

- **V13:** `go list -f '{{.Dir}}' github.com/sunholo-data/ailang/serveapi/protocol` → `/Users/voightkampff/go/pkg/mod/github.com/sunholo-data/ailang@v0.33.2/serveapi/protocol`; `sed -n '1,55p' /Users/voightkampff/go/pkg/mod/github.com/sunholo-data/ailang@v0.33.2/serveapi/protocol/a2a_wire.go` → A2AError29 signature id json.RawMessage, response30 contains jsonrpc2.0/id/error{code,message}, Content-Type application/json31, WriteHeader StatusOK32, json.Encoder33. No omitempty on map id; A2ARequest11 ID json.RawMessage. Source read in V10 shows codeInternal68=-32603 and notAvailableMessage88 literal, with resolver277 nil id and registry320/nil-coordinator327 req.ID, all same code/message.

- **V14:** `go env GOROOT` → `/Users/voightkampff/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.6.darwin-arm64`; `sed -n '259,272p' /Users/voightkampff/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.6.darwin-arm64/src/encoding/json/stream.go` → RawMessage.MarshalJSON264–269 returns []byte("null") when nil. Thus nil RawMessage in A2AError map encodes id:null; fixture oracle is source-derived, not a measured HTTP resolver probe.

- **V15:** `rg -n 'classUnavailable|msgUnavailable|classDeadlineExceeded|msgDeadlineExceeded' host/projection/projection.go` →96–99 exact card strings stated above;459 deadline504,462 otherwise503. Positive control: same query matches both definitions and live Fail call sites. `sed -n '25,35p;151,163p' host/daemon/handlers.go` → APIError.Error JSON error; detail Class/Message JSON class/message; writeAPIError151–153 delegates status and those exact strings to writeJSON. V10 source confirms errors.Is(DeadlineExceeded) only, so cancellation follows503.

- **V16 (controller-run, read by designer):** `rg -n 'FAIL|TestF5WallClockTimeoutHasElapsedBound|TestRunContextCancel|TestInterfaceV2MovesWhenAnExportedADTGainsAConstructor|TestUnknownEntryInterpreterFailsAbsent|valid_ledger' /tmp/world-iter215-base-verifygo.log` → failures capsule68/70, pkgproj79, replay87 (archive --version10s deadline), verifygate102–103 (valid_ledger). Controller reports pristine outside-sandbox pinnedv0.41.0 verify_go terminal rc1; this designer did not execute that full command. These are baseline observations, not causation findings.

- **V17 (controller-reported, not executed by designer):** `AILANG_BIN=$HOME/.pinned-ailang/ailang go test ./host/capsule ./host/pkgproj ./host/replay ./host/verifygate -p 1 -count=1 -timeout 300s -run '^(TestF5WallClockTimeoutHasElapsedBound|TestRunContextCancel|TestInterfaceV2MovesWhenAnExportedADTGainsAConstructor|TestUnknownEntryInterpreterFailsAbsent|TestMissionConfigDefaultFlowIntegration)$'` → rc0; capsule3.678s, pkgproj9.230s, replay1.211s, verifygate0.681s. Same pristine base/five selected tests diverged under narrowed serialized execution; this proves neither full gate green nor a unique cause for V16.

## Quorum revision record

Round1 BLOCKED3/3 at `/tmp/world-iter215-quorum/w-a2a-operator-log-coverage-2026-10-01T10-29-00Z.json`, read in full; Sonnet absent/quota. GLM: added complete caller/label/denial proofs and helper reuse rationale (V10–V12). Kimi: added actual imported protocol, constant/status/ID-serialization and card envelope evidence (V13–V15). Gemini: accepted supplied scoped fix verbatim—separate inline `card refusal` format in writeUnavailable; A2A helper unchanged. Corrected counted-fixture wording and recorded broader pristine baseline reds (V16). One author revision spent; re-quorum pending.

## Axiom compliance / conflict surface

| Axiom | Score | Reason |
|---|---|---|
| A1 Determinism | 0 | Semantic result unchanged |
| A2 Replayability | 0 | Journal unchanged |
| A3 Effect legibility | +1 | Missing operator effect paths become explicit |
| A4 Explicit authority | 0 | Admission unchanged |
| A5 Bounded verification | +1 | Deterministic scoped fixtures and mutations |
| A6 Safe concurrency | 0 | Same sink and synchronous call sites; no new workers |
| A7 Machines first | 0 | Constant machine wire preserved |
| A8 Minimal syntax | 0 | No syntax changes |
| A9 Cost visibility | 0 | No cost semantic change |
| A10 Composability | +1 | Existing injected sink reused |
| A11 Structured failure | +1 | Operator-path observability pinned |
| A12 System boundary | 0 | Host transport boundary retained |

Net +4; no hard violation. No parser/typechecker/codegen conflict surface. Transport conflict
surface is ordinary session denial vs resolver error, clean absent head vs unreadable present
head, deadline vs cancellation, nil coordinator vs dispatch failure: each distinction is
preserved and gets its named negative control. Independent quorum is still pending; this
artifact does not authorize a sprint by itself.
