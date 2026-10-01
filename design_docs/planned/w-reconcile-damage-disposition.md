# w-reconcile-damage-disposition — What happens after reconcile refuses a damaged row (row 125)

Row 122 (`d7364a4`) made `Coordinator.committed()` refuse a resolved invocation whose rows fail to
re-hash or fail to name each other (`*IntegrityError`). It left three things open (row-122 plan
§6, OPEN-1 / OPEN-2 / "Not covered"). This note settles all three:

1. **Wire.** Keep the class. `*IntegrityError` stays on `-32603` with the constant
   `notAvailableMessage`. R6 does not move. **The public wire contract does not change.**
2. **Quarantine.** Do not quarantine. A damaged reconcile row is **reported to the operator**
   through a new one-line `ErrorLog` on the A2A projection, which the daemon wires to its
   existing `d.errLog`. The coordinator's `Store` seam does not change.
3. **Absent / undecodable rows.** Type them. An absent row, or verified bytes that do not decode,
   becomes `*IntegrityError` with a new `Kind` field (`mismatch` / `absent` / `undecodable`).
   A store **read error** (`err != nil`) stays wrapped and untyped, so `context.DeadlineExceeded`
   still maps to R11.

Parent design: [w-transition-invocation-coordinator.md](w-transition-invocation-coordinator.md)
(refusal table under "/a2a/ wiring (M5)"). Row-122 plan:
[w-invocation-reconcile-object-integrity-sprint-plan.md](w-invocation-reconcile-object-integrity-sprint-plan.md).

## Problem (measured at `8d86c4f`)

- **Absent and undecodable rows reach the wire with no type.** I ran a throwaway probe test
  (deleted before commit, V9). It resent a committed task id through a store that hides one row,
  or through a receipt that names the wrong object. Each case returned a `*fmt.wrapError`, not an
  `*IntegrityError`. When the row is absent (`ok=false`, `err=nil`), the message ends in
  **`%!w(<nil>)`**, because `coordinator.go:176/194/201` formats `%w` with a nil error. The
  "undecodable" cases are a receipt naming the transition source (`decode record: invalid
  character 'm'…`) and a hash-correct record whose `Output` is not a ref (`output ref: hashref:
  missing algorithm tag…`). Both are also untyped.
- **The operator sees nothing.** The `/v1/` routes write every sanitized 500 to `d.errLog`
  (`handlers.go:177-179`). The A2A projection has no error log (V6). The not-available class
  covers R4, R6, R13-integrity and every untyped store failure, and it hides the cause from the
  client by design, so today that cause reaches **no one**.
- **Quarantine is a lifecycle switch, not a damage response.** It is an in-memory `atomic.Bool`
  (`store.go:230`). Its doc comment says "refuses new operations while retained cleanup waits for
  workers" (`store.go:249`). Its one production caller is the startup budget abort
  (`daemon.go:679`). On a quarantined store, the probe's dispatch failed at the **registry head
  read**, before any reconcile (V9 `quarantined_store`). `GET /v1/health` reads no store
  (`daemon.go:743-750`), so a quarantined daemon would still report `"status":"ok"`.

## Verification Log

All rows: VERIFIED BY ME at HEAD `8d86c4f`, in worktree `.wt-world-iter214-design`.

| V | Claim | Command / file:line → observed |
|---|---|---|
| V1 | `IntegrityError` has only `InvocationID`, `Object`; `SourceError` uses a `Kind` field | `grep -n 'type IntegrityError\|type SourceError' host/coordinator/errors.go` → `15`, `92`; `sed -n 88,100p` → fields `InvocationID`, `Object // "record" \| "output" \| "world"` |
| V2 | absent rows wrap a nil error | `grep -n 'ok=%v: %w' host/coordinator/coordinator.go` → `176`, `194`, `201` |
| V3 | the decode and output-ref failures are untyped | `grep -n 'decode record\|output ref:\|reconcile %s: output: %w' host/coordinator/coordinator.go` → `183`, `190`, `214` |
| V4 | reconcile has one entry point | `grep -n 'return c.committed(ctx, rc)' …/coordinator.go` → `243` |
| V5 | `dispatchError` has no arm for `IntegrityError`/`SourceError`; the ctx arm precedes `default` | `grep -n 'context.Canceled), errors.Is\|default:' host/projection/projection.go` → `427`, `433`; `default` returns `codeInternal, notAvailableMessage` (`434`); constant at `87` |
| V6 | the projection has no error log; the daemon's sanitized 500s do | `grep -n 'log\.\|slog\|Logger\|Printf' host/projection/projection.go host/coordinator/*.go \| grep -v _test` → **0 lines**; `grep -n errLog host/daemon/handlers.go` → `178 fmt.Fprintf(d.errLog, "ailang-worldd: internal error: %s %s: %v\n", …)`; `ErrorLog` doc `daemon.go:227-239` |
| V7 | the wire table pins R6 and R13-integrity to the not-available class | `host/projection/projection_test.go:822` `R6_corrupt_source`, `:823` `R13_reconcile_integrity`, both → `codeInternal, notAvailableMessage`; wire-level `a2a_wire_test.go:269` `R4_R6_default` |
| V8 | Quarantine is in memory, has one caller, and is lifecycle-scoped | `grep -rn 'Quarantine(' --include='*.go' host cmd \| grep -v _test` → `daemon/daemon.go:679 store.Quarantine(d.store)` only (plus `checkQuarantine` sites); `store.go:230 quarantined atomic.Bool`, `:249-251` |
| V9 | probe: absent / undecodable / quarantined behaviour at HEAD | throwaway `host/coordinator/zz_probe214_test.go` (`hideStore` + the row-122 `tamperStore`), `go test -count=1 -run '^TestProbe214$' -v ./host/coordinator/` → 7 `=== RUN` (parent + 6). `absent_record`/`absent_output`/`absent_world`: `type=*fmt.wrapError integrity=false msg="coordinator: reconcile a2a:ep1:t1: <row>: ok=false: %!w(<nil>)"`. `undecodable_record_names_source`: `…decode record: invalid character 'm' looking for beginning of value`. `record_bad_output_ref`: `…output ref: hashref: missing algorithm tag…`. `quarantined_store`: fails before dispatch, `NewRequest: … read transition registry head: store: quarantined`. File removed; `git status --short` empty afterwards |
| V10 | `/v1/health` reads no store | `sed -n 743,750p host/daemon/daemon.go` → body built from `Version`, `cfg.DBPath`, `interpreterRef/Version` only |
| V11 | startup damage is reported and the daemon keeps serving (precedent) | `daemon.go:932-941`: holes are written as `integrity_hole …` lines from a goroutine, then `Serve()` runs regardless |
| V12 | a resolved receipt and its three rows land in one transaction, and nothing deletes objects or worlds | `store.go:1118 beginDurable(ctx, "commit")` … `:1240 finishDurable`; the `objects`/`worlds` inserts and the `'outcome'` journal row all sit inside it; `plan.go:103 objects := []store.Object{in, out, rec}`; `grep -rn 'DELETE FROM' host \| grep -v _test` → only `session_credentials` (`store.go:1370`) |
| V13 | the store's own read failures are untyped `fmt.Errorf` (I/O and stored-text decode alike) | `store.go:600-606` (`get object`, `interface hash`), `:668-677` (`get world`, `state root`, `log head`) |
| V14 | the wire rig accepts any `coordinator.Store`, so a tampered store can be driven end to end | `a2a_wire_test.go:71 func (r *wireRig) handler(t, cs coordinator.Store, …)` |
| V15 | `parseOutput` refuses a nil map as well as a decode error (the check (3) mirrors) | `grep -n 'func (c \*Coordinator) parseOutput' -A10 host/coordinator/coordinator.go` → `148 func … parseOutput`; `154 if err := json.Unmarshal(out, &obj); err != nil \|\| obj == nil {`; `155 return nil, nil, &OutputError{Reason: "is not a JSON object"}` |
| V16 | `json.Unmarshal` into `map[string]any`: `null` → nil map, no error; `[1]` → error; `{}` → non-nil | throwaway `host/coordinator/zz_probe214b_test.go`, `go test -count=1 -run '^TestProbe214b$' -v ./host/coordinator/` → `unmarshal null: nil=true err=<nil>`; `unmarshal [1]: nil=true err=json: cannot unmarshal array into Go value of type map[string]interface {}`; `unmarshal {}: nil=false err=<nil>` |
| V17 | a hash-consistent forged chain (record → output → world, receipt pointed at it) carrying a `null` output **reconciles as success with a nil `Output`** at HEAD; `[1]` fails untyped; `{}` reconciles (the control, which proves the chain passes C1–C5 and reaches line 213) | same probe, 4 `=== RUN` (parent + 3), `--- PASS`, tamper `hits=1` in each. `output=null: err=<nil> integrity=false reconciled=true outputNil=true outputBytes="null"`; `output=[1]: err=coordinator: reconcile a2a:ep1:t1: output: json: cannot unmarshal array … integrity=false reconciled=false`; `output={}: err=<nil> … reconciled=true outputNil=false`. Probe file removed; `git status --short` empty afterwards |
| V18 | `projection.New` returns `(*Handler, error)` and already rejects nil required seams with a constant error; a table test pins each one | `grep -n 'func New(cfg Config)' -B3 -A22 host/projection/projection.go` → `167-168` Decision-6 comment; `169 func New(cfg Config) (*Handler, error)`; `171-180` `cfg.Resolver/Reader/Heads/Deny/Fail == nil` → `errors.New("projection: X is required")`; `181-184` MaxWait/InvokeWait. Test: `projection_test.go:1373 func TestProjection_ConfigValidation`, `strip` table at `:1397-1406` (`{"Fail", func(c *Config) { c.Fail = nil }}` …), each `New(c)` must error (`:1410`) |
| V19 | the `projection.Config` construction sites: one production, two test (revision 2: global search) | `grep -rn 'projection\.New(\|projection\.Config{' --include='*_test.go' host cmd` → `host/daemon/invoke_e2e_test.go:244` only (in-package tests call unqualified `New(`: `projection_test.go:226 New(cfg)` via `testConfig` `:211-212`, plus `TestProjection_ConfigValidation` `:1379-1414`); production:  `grep -rn 'projection.New(' --include='*.go' host \| grep -v _test` → `daemon/daemon.go:580` only; `grep -n 'Config{' host/projection/*_test.go` → `projection_test.go:212` (`testConfig`, which the wire rig reuses: `a2a_wire_test.go:86 cfg := testConfig(r.st)`); `host/daemon/invoke_e2e_test.go:244 projection.New(projection.Config{Resolver: d.resolver, …` |
| V20 | the ctx arm maps deadlines to R11, not `default` (so it never takes the not-available log path) | `sed -n 423,435p host/projection/projection.go` (controller, revision 2) → `427 case errors.Is(err, context.Canceled), errors.Is(err, context.DeadlineExceeded):` / `428 return codeInternal, "invocation exceeded its deadline"`; `433 default:` / `434 return codeInternal, notAvailableMessage` |
| V21 | every committed output row passed `parseOutput`, so a non-object output at reconcile is necessarily damage | `grep -rn 'parseOutput(\|planInvocation(\|OutputV1' --include='*.go' host cmd \| grep -v _test` (controller, revision 2) → `coordinator.go:308 outBytes, outObj, err := c.parseOutput(res.Stdout) // R12` (sole call site; `err != nil` returns before planning, `:309-311`), `:312 pl := planInvocation(…, outBytes, …)` (sole call site), `plan.go:84 out := object(OutputV1, output)`, `:103 objects := []store.Object{in, out, rec}`; no other `OutputV1` writer. Throwaway `TestProbe214c` (deleted): `parseOutput("null")`, `("[1]")`, `("null\n")` → `coordinator: transition output is not a JSON object`; `("{}")` → `err=<nil>` (control) |
| V22 | `Fail` is the card route's 503/504 renderer, a different seam from the operator log | `projection.go:136-137` → `// Fail renders the card route's non-denial 503/504 failures.` / `Fail ErrorWriter` |

## Design

### (1) Wire: keep the not-available class; R6 does not move

**Decision.** `*IntegrityError`, of every `Kind`, maps to `-32603`
`transition invocation is not available in this daemon`. M2 makes this an **explicit arm** in
`dispatchError` instead of relying on `default`, so it is a mutation unit of its own.
`*SourceError` keeps the `default` path. The message is unchanged and the public contract is
unchanged.

**Why the generic message is the honest one.** A damaged reconcile means the daemon cannot verify
its own journal for this id. In `receipt_names_other_record` (row 122) the damaged row is the
**journal intent itself**. So the daemon cannot truthfully say "committed" ("stored record is
damaged" would imply that). Nor can it say "not committed" (R15's message, which invites a new
task id and so a possible **second application** of a transition that did land). It cannot say
"resend" either (R16, which loops forever on deterministic damage). "Not available in this
daemon" asserts nothing about the outcome, and that is all the daemon can stand behind. The
client cannot repair a store; only the operator can, and (2) gives the operator the detail.

**Rejected: a dedicated message, for example "stored invocation record is damaged".** (a) It
changes the public contract pinned by `TestA2ADispatch` (`projection_test.go:822-823`) and by the
parent's refusal table. (b) It claims more than the daemon knows (above). (c) Moving R6 with it
would relabel "registry names a corrupt or absent source". That is publication-side damage with a
different remedy (republish), and it would merge two operator stories under one client string.
**Rejected: a new JSON-RPC code.** `-32603` already means server-side. A new code is the same
contract change with nothing gained.

### (2) Quarantine: no. Report to the operator instead

**Decision.** A detected reconcile integrity failure does **not** call `store.Quarantine`. The
projection gains `Config.ErrorLog io.Writer`, which `projection.New` requires to be non-nil. That
is one more `case cfg.ErrorLog == nil:` arm in the existing required-seam switch (V18), the
Decision-6 "construction error, never silent" pattern. The daemon passes its resolved `d.errLog`.
Whenever a dispatch refusal maps to `notAvailableMessage`, the handler writes **one line** via
`fmt.Fprintf(h.errorLog, "ailang-worldd: a2a refusal: %s %s: %q\n", method, id, errText)`. `%q`
escapes embedded newlines and control bytes, so the one-line guarantee holds for arbitrary store
and bind error text, not just coordinator strings. For bind failures (R4) where no invocation id
parsed, `id` renders as `-`. (Revision 2, kimi r2's proposed fix, applied verbatim.) `ErrorLog`
does not overload `Fail`: `Fail` renders the card route's 503/504 responses to the client (V22);
`ErrorLog` is the operator's sink and never reaches a client.

This uses the same audience split as `writeInternalError` (`handlers.go:159-179`): the wire gets
the constant, and the process owner gets the cause. With (3) in place, `<verbatim error>` reads,
for example, `coordinator: reconcile a2a:ep1:t1: recorded record is absent`. The line goes to
`ErrorLog` only, never to the announce writer.

**Reasons for not quarantining.**

- **Blast radius.** The probe shows a quarantined store refuses the registry head read (V9). That
  means every A2A call, the card route and every `/v1/` store route fail. Meanwhile `/v1/health`
  still says `ok` (V10). One damaged row would become a total, health-invisible outage.
- **No containment gain.** The damage cannot spread through reconcile. Every resend of the
  damaged id refuses again, deterministically: nothing is reported reconciled and nothing
  re-executes (`mustNotRun` in row 122's tests). The forward path never reads old
  record/output rows: its store reads are the receipt (`coordinator.go:236`), the transition source
  (`:274`, R6), the selected head (`:284`) and the head world (`:291`); `committed()` is reached only from the receipt branch (V4).
- **Not durable, so not recoverable either.** Quarantine is an in-memory bool (V8). A restart
  clears it and reopens the same damaged file, so it buys an outage, not a repair. No repair
  tooling exists for it to wait on.
- **Precedent.** The startup integrity sweep reports holes and keeps serving (V11). Runtime
  damage found on one id should not be treated more harshly than damage found at startup.
- **Semantics.** `Quarantine` is documented as a lifecycle action for retained cleanup (V8).
  Using it for damage would overload one bit with two meanings that need different clearing.

**What the operator sees and does.** The operator gets one `a2a refusal` line per refused resend,
naming the invocation id, the row (`record`/`output`/`world`) and the kind. The store keeps
serving every other id. Clearing is out of band (restore from backup). Until then the id keeps
refusing, which is the correct steady state. **Seam change:** `projection.Config.ErrorLog` plus
one daemon wiring line. The coordinator's `Store` interface (`coordinator.go:27-34`) does **not**
change.

**Rejected: a coordinator-side `Quarantine` seam, or a narrower "read-only" quarantine.** Both add
policy machinery and a clearing story for a failure that already contains itself.
**Rejected: logging in the coordinator.** That would add an effect sink to a component whose
effects are exactly its `Store` and `Runner` seams. The projection is the boundary that owns
"sanitize for the wire, detail for the owner".

### (3) Absent and undecodable rows: type them as `IntegrityError{Kind}`

**Decision.** `IntegrityError` gains `Kind string // "mismatch" | "absent" | "undecodable"`, after
`SourceError{Kind}`. `Error()` renders `reconcile <id>: recorded <object> <phrase>`, with phrases
`does not match its reference` / `is absent` / `does not decode`. In `committed()`:

| Site (HEAD line) | Today | After |
|---|---|---|
| record `GetObject` `ok=false, err=nil` (176) | wrap of nil | `IntegrityError{record, absent}` |
| output `GetObject` `ok=false, err=nil` (194) | wrap of nil | `IntegrityError{output, absent}` |
| world `GetWorld` `ok=false, err=nil` (201) | wrap of nil | `IntegrityError{world, absent}` |
| verified record bytes fail `json.Unmarshal` (183) | untyped | `IntegrityError{record, undecodable}` |
| `rec.Output` not a ref (190) | untyped | `IntegrityError{record, undecodable}` (the record names no output) |
| verified output bytes not a JSON object (214) | `[1]`: untyped; `null`: decodes to a nil map and **reconciles as success with a nil `Output`** (V16, V17) | `IntegrityError{output, undecodable}` when `err != nil \|\| obj == nil`, the same predicate `parseOutput` applies at commit time (`coordinator.go:154`, V15); every committed output passed it (V21), so this refuses only damage |
| any `err != nil` from the store (176/194/201) | wrap with `%!w` noise when nil | **unchanged class**: `fmt.Errorf("…: %w", err)`, split from the `!ok` branch |
| row-122 mismatch sites (179/186/197/207/210) | `IntegrityError{Object}` | add `Kind: "mismatch"` |

**Why `absent` is damage, not "not found".** The resolved receipt and all three rows are written
in **one** transaction, and nothing in `host/` deletes objects or worlds (V12). A resolved receipt
whose row is missing therefore names something the commit wrote, so it is damage.
**Why `undecodable` is damage.** It is only reached **after** the bytes re-hash to the ref the
journal names, so the journal names the wrong object. That is a binding failure, like row 122's C2.

**Why store read errors stay untyped.** The store wraps I/O failures and stored-text parse
failures (`state root`, `log head`, `interface hash`) with the same untyped `fmt.Errorf` (V13).
Typing them at the coordinator would mislabel a transient `SQLITE_BUSY` or a deadline as
permanent damage. It would also break R11, because `context.DeadlineExceeded` must keep reaching
the ctx arm (`projection.go:427`). Distinguishing stored-row damage inside the store is a store
change; see Residuals.

**Rejected: separate `ReconcileAbsentError` / `ReconcileDecodeError` types.** These would mean
three wire arms for one client class and one operator story. **Rejected: mapping absent to R15
"not committed".** That is false (V12), and it invites a second application.

## Milestones (test-first; each ≤ ~150 production lines)

- **M1: coordinator typing (~25 production lines, `host/coordinator/{errors,coordinator}.go`).**
  RED first: extend `TestReconcileRefusesDamagedRecord` with subtests `record_absent`,
  `output_absent`, `world_absent` (a `hideStore` that returns `ok=false, err=nil` for one semantic
  id or world, with a hit counter), `record_undecodable` (the receipt names the transition source),
  `record_output_ref_unparsable` (a hash-correct forged record put via `PutObject`),
  `output_not_object` (two payloads, `[1]` and `null`: a forged record → forged output → forged world, with the receipt's
  `WorldRef = worldRef(forged)`), and `store_read_error_stays_untyped` (`GetObject` returns
  `fmt.Errorf("x: %w", context.DeadlineExceeded)`: assert `!errors.As(*IntegrityError)` and
  `errors.Is(DeadlineExceeded)`), and `commit_refuses_non_object_output` (real coordinator + real
  store, no forgery: a runner whose raw output is `null`, then `[1]`; assert `*OutputError`, no
  receipt, head unchanged), which pins V21 so the mirror argument is tested rather than assumed.
  Every existing case also asserts `Kind == "mismatch"`. Then implement the table above.
- **M2: projection arm, error log and daemon wiring (~25 production lines, `host/projection/projection.go`,
  `host/daemon/daemon.go`).** RED first: `TestA2ADispatch` gains `R13_reconcile_absent`
  (`IntegrityError{Kind:"absent"}` → `-32603 notAvailableMessage`). `TestA2ADispatchWire` gains
  `R13_integrity`: the real coordinator, first dispatch commits, the resend goes through a store
  that hides the record. `assertWireRefusal` checks the constant and the body hides `coordinator:`
  and the invocation id, while the test's `ErrorLog` buffer holds exactly one `a2a refusal` line
  containing `is absent`. `R4_R6_default` additionally asserts its secret appears in the log and
  not the body. A host/daemon test runs a real daemon with a `Config.ErrorLog` buffer and asserts
  the line arrives there (this pins the wiring). `projection.New` with a nil `ErrorLog` is a
  construction error: one `{"ErrorLog", …}` row in `TestProjection_ConfigValidation`'s `strip`
  table (V18).
  Production-line count, from V18/V19: a `Config` field with its doc comment (~4), a `Handler`
  field and its init (2), a `New` arm (2), the explicit `IntegrityError` arm (3), the log write (~4)
  and one daemon field at `daemon.go:580` (1), about 16 in total. Test-side config edits: `testConfig`
  (`projection_test.go:212`, which also feeds the wire rig) and `invoke_e2e_test.go:244`.

Estimate: ~0.25 d, matching the row. Revision 1 re-derived it from V18/V19 and it stands.

## Acceptance Criteria

- **AC1.** Each of the six M1 damage subtests returns `*IntegrityError` with the stated
  `Object`/`Kind`. None reports `Reconciled`, none re-executes (`mustNotRun`), and the head is
  unchanged (`assertHead`).
- **AC2.** A store read error during reconcile is not an `IntegrityError` and still satisfies
  `errors.Is(err, context.DeadlineExceeded)`. No `%!w(` appears in any reconcile error text: the
  M1 subtests assert `!strings.Contains(err.Error(), "%!")`.
- **AC3.** The wire is byte-unchanged for every pre-existing `TestA2ADispatch` row and every
  `TestA2ADispatchWire` subtest. The new rows map to `-32603 notAvailableMessage`.
- **AC4.** Every not-available dispatch refusal writes exactly one `ErrorLog` line that carries the
  verbatim error text with newlines escaped; a store error containing `"\n"` still yields exactly
  one physical line (asserted via `bytes.Count(buf, "\n") == 1`). No other refusal class writes
  one: the R15 row and an R11 deadline row (V20) each assert an empty log. Nothing goes to the
  announce writer.
- **AC5.** No call to `store.Quarantine` is added:
  `grep -rn 'Quarantine(' --include='*.go' host cmd | grep -v _test` still prints only `daemon.go`'s
  existing line. The coordinator `Store` interface is unchanged.
- **AC6.** Gates: `go vet ./...`, `go test -race -count=1 ./host/...`, `./scripts/verify_ail.sh`
  with `AILANG_BIN` set to the pinned binary.

## Non-Vacuity (one named RED mutation per new branch / AC)

| Mutation | Change | Must red |
|---|---|---|
| MUT-ABSENT-REC | record `!ok` branch returns the old untyped wrap | `record_absent` |
| MUT-ABSENT-OUT | same, output | `output_absent` |
| MUT-ABSENT-WORLD | same, world | `world_absent` |
| MUT-REC-DECODE-UNTYPED | record decode branch returns `fmt.Errorf` | `record_undecodable` |
| MUT-OUTREF-UNTYPED | output-ref branch returns `fmt.Errorf` | `record_output_ref_unparsable` |
| MUT-OUT-NIL-OBJ | drop `\|\| obj == nil` | `output_not_object` (forged `null` variant) |
| MUT-ERR-AS-DAMAGE | merge `err != nil` into the `!ok` integrity branch | `store_read_error_stays_untyped` |
| MUT-KIND-SWAP | `absent` labelled `mismatch` | `record_absent` (Kind assert) |
| MUT-WIRE-INTEGRITY | explicit arm returns R15's message | `TestA2ADispatch/R13_reconcile_absent` and `/R13_reconcile_integrity` |
| MUT-NOLOG | delete the `ErrorLog` write | `TestA2ADispatchWire/R13_integrity`, `/R4_R6_default` |
| MUT-LOG-ALL | log every refusal class | R15 row's empty-log assert |
| MUT-LOG-TO-WIRE | interpolate the error into the message | `assertWireRefusal` hidden-string check |
| MUT-LOG-NEWLINE | log write interpolates `%v` without escaping | the newline-carrying store-error subtest (`bytes.Count == 1`) |
| MUT-LOG-DEADLINE | log the R11 ctx arm too | the R11 row's empty-log assert |
| MUT-COMMIT-NULL-OK | `parseOutput` drops `\|\| obj == nil` | `commit_refuses_non_object_output` (`null` arm) |
| MUT-ERRLOG-UNWIRED | daemon passes `io.Discard` | the host/daemon wiring test |
| MUT-ERRLOG-NIL-OK | `New` accepts a nil `ErrorLog` | the `New`-rejects row |

The executor records each as KILLED with a log, as row 122 did. A survivor blocks landing.

## Conflict Surface

- `host/coordinator/coordinator_test.go`: row 122's table gains a `Kind` assertion and new rows.
  No existing expectation flips.
- `host/projection/projection_test.go`: `testConfig` (`:212`) supplies `ErrorLog`. The wire rig
  reuses it (`a2a_wire_test.go:86`), so it needs no edit of its own. `TestProjection_ConfigValidation`
  gains one strip row (V18).
- `host/daemon/invoke_e2e_test.go:244`: a direct `projection.Config` literal that must add
  `ErrorLog` (V19).
- `host/daemon/daemon.go`: one field in the `projection.Config` literal (around line 580).
- Parent design's refusal table, the "R4 bind error, R6 source error, other store errors" row:
  append "R13 reconcile integrity (row 122 / row 125: mismatch, absent, undecodable)". Doc-only.
- **Public wire contract: unchanged.** No `/v1/` vocabulary change. No store change.

## Axiom Compliance (coding-standards.md)

- **S1**: no `world/` function changes. The new branches are host-boundary integrity checks, not
  kernel laws.
- **S2**: the only new effect (the log write) lives at the host boundary (projection), on a
  daemon-owned writer. The coordinator gains no effect.
- **S3**: no package question arises: this is daemon plumbing.
- **S4/S5**: this note contains no `.ail` code, and no AILANG syntax claim is made.
- **S6**: every new branch has a named mutation (above). AC5 is a measured absence, not an assumed
  one.
- **S7**: no operator command changes. If QUICKSTART describes stderr content, add the
  `a2a refusal` line shape there.

## Residuals

- **R-125-0 (noted in revision 1, not widened):** V17's `{}` control shows that a forged chain
  which is hash-consistent at every link reconciles. Row 122's checks bind rows to the refs the
  journal names; they cannot detect a journal intent that was rewritten to a self-consistent
  forgery. That is out of this row's text; owner: a future journal-authentication row, if wanted.

- **R-125-1: store-level damage typing.** The store does not tell a stored-text parse failure
  apart from I/O (V13). Owner: a future store row, if a consumer needs it.
- **R-125-2: the forward path trusts the head world row.** `Dispatch` reads `GetWorld(head)`
  without re-deriving its ref (`coordinator.go` R7 block). Damage there is planned over, not
  refused. Owner: a new clause-6 row, if wanted; the store's `Commit` compare-and-append bounds it.
- **R-125-3: damage in the receipt read.** `GetReceipt` failures stay untyped
  (`coordinator: receipt: %w`). With M2 they are logged, but not typed.
- **R-125-4: card-route failures do not log either.** `writeUnavailable` maps to 503/504 without
  an operator line. M2's `ErrorLog` makes this a one-line follow-up.
- **R-125-5: repair tooling.** There is no operator command that marks or repairs a damaged
  invocation. This is out of scope until a real incident needs one.

## Open decisions for the human

**None required.** The note keeps the public wire contract, adds no charter policy and calls no
new lifecycle action. If Mark later wants a dedicated client message (OPEN-1 reopened), it is a
one-arm change in `dispatchError` plus the two pinned rows. The recommended default stays "keep
the class", for the honesty reason given in (1).

## Quorum log

### Round 1 (iteration 214): BLOCKED 3/3 (glm-5.3, kimi-k3, gemini-3-1-pro; gpt6-1-sol absent)

All three objections were of the form "a load-bearing premise is asserted without a V-row". None
disputed the design direction.

| Objection | Reviewer | Answer (revision 1) |
|---|---|---|
| The (3) table's line-214 row claims a hash-verified `null` output "decodes to a nil map and succeeds", with no V-row | glm, gemini | **Measured and confirmed.** V16 (standalone unmarshal) and V17 (`null` driven through a real reconcile on a forged hash-consistent chain → `err=<nil> reconciled=true outputNil=true`). The table row now cites both. |
| "`\|\| obj == nil` mirrors `parseOutput`", but `parseOutput` is never cited | glm, gemini | **Measured and confirmed.** V15 cites `coordinator.go:148-158`, predicate at `:154`. The table row now cites it. |
| M2's "`projection.New` enforces required config as a construction error" and its "~25 lines / one field in daemon.go" estimate have no V-row | kimi | **Measured and confirmed.** V18 pins `func New(cfg Config) (*Handler, error)` and its existing nil-seam switch plus `TestProjection_ConfigValidation`. V19 counts the construction sites: one production, plus the test-side `invoke_e2e_test.go:244` literal that the first draft missed (now in M2 and the Conflict Surface). The estimate is re-derived at about 16 production lines and stands. |

### Round 2 (iteration 214): BLOCKED 3/3 present (glm-5.3, kimi-k3, gemini-3-1-pro; gpt6-1-sol absent, `openai error (429): You have no credits remaining`)

Objections landed on three NEW surfaces, none disputing the direction, each with a concrete
reviewer-authored `proposed_fix`. Under Gate 2's narrow-refinement carve-out the controller made a
bounded revision 2 applying those fixes, with the controller's own measurements (V20–V22):

| Objection | Reviewer | Answer |
|---|---|---|
| line-214 branch rests on an unverified invariant that every committed output passed `parseOutput` | glm | **Measured and confirmed** (V21: sole call site, dataflow to `plan.go:103`, probe refuses `null`/`[1]`). Added M1 RED `commit_refuses_non_object_output` and MUT-COMMIT-NULL-OK verbatim from the fix. glm's catches: `Fail`'s role stated (V22); the forward-path claim cited; the R11 arm does not map to not-available (V20), now pinned by an empty-log row and MUT-LOG-DEADLINE. |
| AC4's one-line guarantee breaks on a `\n` in store/bind error text | kimi | Format bullet replaced with the proposed `%q` line and `-` for an unparsed id; AC4 extended with `bytes.Count == 1`; MUT-LOG-NEWLINE added; V20 quotes the full ctx arm. |
| V19's test-site bound used a directory-scoped grep | gemini | V19 re-run with the proposed global grep: `invoke_e2e_test.go:244` is the only qualified test site; in-package tests use `testConfig`/`New(`. Conflict Surface unchanged. |

