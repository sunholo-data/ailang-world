# w-finish-approve-record-guard-pin: no test pins the approve guard's `|| input.Record != nil` arm

- **Status**: Planned (r1, iteration 335). Designed iteration 335, unattended.
- **Item**: queue row 237 (charter `design_docs/world-mission.md:6309`), clause-3 (MET; hygiene),
  controller-filed in iteration 327 from the row-236 judge's NB1. The guard arm is correct. Nothing
  tests it. The fix is test-only: one new `_test.go` file, **0 production lines**.
- **Estimated**: ~0.1 day. A **prototype** is in the worktree at
  `host/daemon/episode_finish_approve_guard_test.go` (61 lines, sha256
  `f5df08b346dd833250ea444ed15427fb4cec1fe818f05273e0fde5187134d2e5`, V7). It is a prototype for
  measuring the drills. The planner pins the final diff.
- **Measurement base**: `77e0dd7c` (worktree `sprint/world-iter335-finish-approve-record-guard-pin`).
- **Instruments**: `go1.26.6 darwin/arm64`. All runs are `go test -count=1` from the worktree root.
  `AILANG_BIN=$HOME/.pinned-ailang/ailang` (v0.41.0) was set for the whole-package run only. Raw
  output is in `verification/world-iter335/designer/`: `drills.sh` (applies each mutant, runs it,
  restores it and prints the sha256 before, edited and restored), `drills.txt`, `facts.txt`, `pkg.txt`.

---

## 1. The gap

`handleEpisodeFinish` (`host/daemon/episode_finish.go`) runs in this order for the approve phase (V2):

1. `:160` `consumeOperatorProof` checks the single-use proof file and deletes it (`:129`
   `os.Remove(path)`).
2. `:166` the approve branch is entered.
3. `:169` the guard runs:
   `if err != nil || scope.ExpiresAt <= now || scope.ExpiresAt > now+3600 || input.Record != nil {`.
4. `:170` it refuses with `operator route: invalid or expired finish scope`.
5. `:173` otherwise `broker.RecordOperatorApproval` writes the approval.

The last arm stops an approve request from also carrying a `record`. Without it, such a request is
approved, and the record is silently ignored. The arm was added with the route in `8c5a14e6`
(row 218 M2, V3).

**The row's measurement reproduces (V4).** With the mutant
`s/ || input.Record != nil {/ {/` applied, every existing Finish test in `host/daemon` passes:
`go test -count=1 -run Finish ./host/daemon/` → `ok`. The controller's premise file
(`verification/world-iter335/premise-measure.txt`) records the same for `./cmd/ailang-worldd/` and for
the whole `./host/daemon/` package.

**The refusal text is not unique to this arm.** All four clauses of `:169` share one refusal (`:170`),
so a 403 with that text cannot show *which* clause fired. Only a control that changes `Record` alone
can show that.

## 2. The decision: one new test with a refused arm, a no-write check and a control

### 2.1 Where it lives

The test is a **new function in a new file**: `host/daemon/episode_finish_approve_guard_test.go`,
`TestEpisodeFinishApproveRefusesRecord`, with one helper, `approvalWrites`.

- **Why not an arm in `TestEpisodeFinishScopeExpiry`'s table.** That table shares one daemon across
  its arms, and its 200 arms write approvals (`{1, 200}` runs before `{3601, 403}`). So it cannot
  assert that *no* approval exists. It also checks the status code only.
- **Why a new file rather than appending to `episode_finish_test.go`.** The new file conflicts with
  nothing (§5) and gives the planner a sha for the whole change. If the planner prefers it, the same
  code can go at the end of `episode_finish_test.go` (it uses only that file's helpers and
  `newHandlerDaemon`/`boundedTestContext` from `handlers_test.go`). Behaviour is unaffected either way.

### 2.2 The test

It uses one fresh daemon (`newHandlerDaemon`), the pinned clock `pinFinishClock(d, finishClockAt)`
and one scope, `finishScopeAt(finishClockAt + 900)`. That scope is valid on the pinned clock:
`TestEpisodeFinishScopeExpiry` accepts offsets from +1 to +3600.

**Arm 1, refused (runs first, on an empty store).**
- Proof: `finishProofFile(t, d, "approve", finishClockAt+100, 0o600)`.
- Request: `FinishRequest{Phase: "approve", Scope: scope, Record: &EpisodeFinish{Episode: "ep1"}}`.
  Any non-nil record will do. `Record` is a pointer with `omitempty`, so this non-nil value is
  sent on the wire as an object and decodes back to non-nil (V2).
- (a) It asserts `w.Code == 403` and that the body contains
  `operator route: invalid or expired finish scope`.
- (b) It asserts that no approval was recorded, using `approvalWrites`, which reads the store
  directly. Three things must all hold:
  - the `broker.ApprovalsV1` registry head (`world/approvals/v1`) is **absent**. `appendApprovalHead`
    moves this head on every `RecordOperatorApproval` (V5);
  - `d.store.ObjectsBySemanticID(ctx, broker.ApprovalRequestV1, "", 10)` is empty;
  - the same call for `broker.ApprovalDecisionV1` is empty.

  These are the exact writes `RecordOperatorApproval` makes: `PutObject(request)`, then
  `appendApprovalHead`, then `decideApproval` (V5). The check reads the store, not the response,
  so a handler that writes and then refuses cannot pass it (M3a).

**Arm 2, the control (runs second).**
- Proof: a **fresh** proof, same expiry.
- Request: the identical request with `Record: nil`.
- (c) It asserts `200`, decodes the `FinishApproval`, and then checks that `approvalWrites`
  sees:
  - the head **present**;
  - **exactly one** request object, whose hash equals `approval.Request`;
  - **exactly one** decision object, whose hash equals `approval.Approval`.

  This makes the 403 in arm 1 attributable to `Record` alone, not to the scope, clock, proof or
  episode. It also shows the instrument works: the same helper that reads "nothing" in arm 1 does
  see the writes when they happen.

**Why a fresh proof per arm.** `consumeOperatorProof` deletes the proof at `:129`, and it runs
before the guard at `:160`. So the **refused** arm also uses up its proof, and reusing it would
get a 403 from the proof check, not from the guard. AC11's `replay` check pins that
(`episode_finish_test.go:112-113`, V10). Each arm writes its own `ep1.json`, as the existing tests do.

**Why the refused arm must run first.** With the clock pinned, arm 1 (under mutants M1/M3a) and the
clean control write the **same** request and decision hashes. `drills.txt` shows `sha256:9d1d990a…`
and `sha256:ff79441f…` in both. If the control ran first, the object lists could not tell a write
from no write. Only the head's wire, which carries `PreviousHead`, would differ. Running arm 1 on an
empty store makes "absent / empty" unambiguous.

## 3. Mutation drills (all first-party, `drills.sh`, output `drills.txt`)

Production mutants were applied to `host/daemon/episode_finish.go` (sha `4e320c537e68`) and restored
with `git checkout --`. Test-side mutants edited the prototype (sha `f5df08b346dd`) and were restored
from a copy. Every edited tree passed `go vet ./host/daemon/` before the run. After the drills,
`git status --porcelain` showed only the prototype and `verification/world-iter335/`.
Killer lines are prototype lines.

| ID | Mutant | Run | Verdict | Killer |
|---|---|---|---|---|
| M1 | **the row's mutant**: `s/ \|\| input.Record != nil {/ {/` (1+/1−) | new test | **KILLED** rc=1 | `:45 approve with record: 200 {"request":"sha256:9d1d…","approval":"sha256:ff79…"}` |
| M1′ | M1, existing Finish tests only (`-run Finish -skip <new>`) | 6 existing | SURVIVED rc=0 (reproduces the row) | none |
| M3a | **write-then-refuse**: M1 plus `if input.Record != nil { refuse(<same text>); return }` inserted after `RecordOperatorApproval` (5+/1−). A 403 with the right text | new test | **KILLED** rc=1 | `:48 refused approve recorded an approval: head=true requests=[sha256:9d1d…] decisions=[sha256:ff79…]` |
| M3c | **a different early refusal**: `if true \|\| err != nil \|\| …` at `:169`, so approve always refuses with the same text | new test | **KILLED** rc=1 | `:55 control approve without record: 403 {… invalid or expired finish scope …}` |
| M2 | test: control arm (prototype `:51-60`) deleted, `encoding/json` import dropped | clean prod | PASS rc=0 (expected) | none |
| M2+M3c | M2's test against M3c | M3c prod | **SURVIVED** rc=0. **The control is load-bearing.** | none |
| M2b | test: the control's `!= 200` inverted to `!= 403` | clean prod | FAIL rc=1. The clean route really does approve without a record | `:55 control approve without record: 200 {…}` |
| M4+M3a | test: the no-approval assertion (`:47-49`) deleted, run against M3a | M3a prod | **SURVIVED** rc=0. **The no-write assertion is load-bearing.** | none |

Each of the prototype's three assertions has a mutant that only it kills. The 403 kills M1, the
no-write check kills M3a, and the control kills M3c.

## 4. Acceptance criteria

- **AC-1 (kills the row's mutant).** Apply M1 and run
  `go test -count=1 -run '^TestEpisodeFinishApproveRefusesRecord$' ./host/daemon/`. Want: FAIL at the
  403 assertion with `approve with record: 200`. Restore.
- **AC-2 (no-write).** Apply M3a, same command. Want: FAIL at the no-write assertion with
  `refused approve recorded an approval: head=true`. Restore.
- **AC-3 (control).** Apply M3c, same command. Want: FAIL at the control with
  `control approve without record: 403`. Restore.
- **AC-4 (green and stable).** On the clean tree:
  - the new test passes `-count=20` and `-race -count=5`;
  - `go test -count=1 -run Finish ./host/daemon/` shows 7 top-level PASS;
  - `go vet ./host/daemon/` and `gofmt -l host/daemon/` are clean.
- **AC-5 (package and full gate).**
  - `AILANG_BIN=$HOME/.pinned-ailang/ailang go test -count=1 ./host/daemon/` → `ok`.
  - Then `go build ./... && go test ./...` and `./scripts/verify_ail.sh` per CLAUDE.md. CI is green
    on the fix head.
- **AC-6 (production untouched).** `git diff --stat <base>.. -- ':!*_test.go' ':!design_docs' ':!verification'`
  is empty.

## 5. Files touched, conflict surface, scope

| File | Change |
|---|---|
| `host/daemon/episode_finish_approve_guard_test.go` (new) | `approvalWrites` helper + `TestEpisodeFinishApproveRefusesRecord`, about 61 lines |

- **Name collisions: none.** `approvalWrites` and `TestEpisodeFinishApproveRefusesRecord` occur only
  in the prototype (V6), and `go vet` compiles the package with them (V7).
- **Helpers reused, not changed:** `finishProofFile`, `finishRoute`, `pinFinishClock`,
  `finishClockAt`, `finishScopeAt` (`episode_finish_test.go`), `newHandlerDaemon`, and
  `boundedTestContext`.
- **S1/S3/S4/S5/S7:** n/a (test-only, host, no `.ail`, no usage surface). **S6:** M1/M3a/M3c (§3).

**Out of scope:**
- **`cmd/ailang-worldd`.** The CLI's Finish tests (5 files, V6) drive the operator flow end to
  end. The guard is a daemon-route property, so it is pinned at the route.
- **The refused arm uses up the proof.** This is existing behaviour, and it is arguably right:
  a proof is spent on any well-formed attempt. The row does not ask to change it.
- **The record phase's `input.Scope != ""` arm (`:182`)** is the mirror check (a record that
  carries a scope). It was not measured here. If it is also unpinned, that is a separate row,
  filed only after the same mutant measurement.

## 6. Ruled out

- **Asserting the 403 alone.** M3c and M3a both return that exact 403 and text. The status code
  cannot tell this arm from the other three clauses of `:169`, or from write-then-refuse.
- **Reading approvals through an HTTP route instead of the store.** The finish response is the
  only route output for this write. A store read is what M3a cannot fake, and `d.store` is already
  used this way by AC11 (`episode_finish_test.go:96`, which already asserts the `ApprovalsV1` head absent after its proof refusals, V10) and `TestEpisodeFinishRecordReuseAndSingleUse`.
- **Head-only check.** The head alone would kill M3a too. The two object lists cost two lines, and
  they let the control bind the exact hashes the response names, which proves the helper observes
  real writes.
- **Running the control first.** See §2.2: the objects are content-addressed and identical under
  the pinned clock.

## 7. Premise check

- **"All 5 `host/daemon` Finish tests" is off by one.** `-run Finish` matches **6** top-level tests at
  base: AC11, ScopeExpiry, OperatorProofExpiryBoundary, RecordRechecksScopeExpiry,
  AuthorizationPrecedesBody and RecordReuseAndSingleUse (V4). The controller's premise file also counts
  6. The conclusion holds: all 6 pass under M1.
- The rest holds: the arm exists at `:169`, the mutant survives, and it pre-dates base
  (`8c5a14e6`, V3).

## 8. Verification Log

| # | Claim | Command | Observed |
|---|---|---|---|
| V1 | Base and tools | `git log --oneline -1`; `git status --porcelain`; `go version` | `77e0dd7c Merge pull request #458 …`; only `?? verification/world-iter335/` before work; `go1.26.6 darwin/arm64` |
| V2 | Route order, guard, refusal text, request type | `sed -n 1,260p host/daemon/episode_finish.go`; `grep -n "consumeOperatorProof(episode\|input.Record != nil {\|os.Remove(path)" host/daemon/episode_finish.go` (`facts.txt`) | `:50-54 FinishRequest{Phase; Scope …omitempty; Record *EpisodeFinish json:"record,omitempty"}`; `:129 os.Remove(path)`; `:160 consumeOperatorProof(…)`; `:169 if err != nil \|\| scope.ExpiresAt <= now \|\| scope.ExpiresAt > now+3600 \|\| input.Record != nil {`; `:170 refuse("operator route: invalid or expired finish scope")`; `:173 broker.RecordOperatorApproval(…)` |
| V3 | The arm's origin | `git log --oneline -S'\|\| input.Record != nil' -- host/daemon/episode_finish.go` | `8c5a14e6 feat(row 218 M2): finish's approval, operator route, …` (the only commit) |
| V4 | Row premise: M1 survives the existing tests; how many there are | `grep -n "^func Test[A-Za-z0-9_]*Finish" host/daemon/*_test.go`; `go test -count=1 -v -run Finish -skip '^TestEpisodeFinishApproveRefusesRecord$' ./host/daemon/` (clean, `facts.txt`); drill M1′ (`drills.txt`) | 6 funcs (m3 `:11`; `episode_finish_test.go:61,177,194,213,246`); clean: 6 × `--- PASS`, `ok`; M1′: `ok … 0.403s rc=0` |
| V5 | What `RecordOperatorApproval` writes | `sed -n 686,700p host/broker/approve.go`; `grep -n "func appendApprovalHead\|func decideApproval" -A25 host/broker/approve.go`; `grep -n "ApprovalsV1 *=\|ApprovalRequestV1 *=\|ApprovalDecisionV1 *=" host/broker/*.go`; `sed -n 1,60p host/store/objects_by_semantic_id.go` | `PutObject(obj)` (ApprovalRequestV1), `appendApprovalHead` → `SetRegistryHead(ctx, ApprovalsV1, …)` with `PreviousHead` in the wire, `decideApproval(…"approve"…)` → ApprovalDecisionV1 object; ids `world/approvals/v1`, `world/approval-request/v1`, `world/approval-decision/v1`; `ObjectsBySemanticID`: "An id no object carries yields an empty slice", limit 1–500 |
| V6 | No name collisions; CLI Finish files | `grep -rn "episode_finish_approve_guard\|TestEpisodeFinishApproveRefusesRecord\|approvalWrites" --include='*.go' --include='*.md' .` (control: `grep -rln finishProofFile host/daemon` → 3 files incl. the prototype); `grep -ln "func Test[A-Za-z0-9_]*Finish" cmd/ailang-worldd/*_test.go` | matches only inside the prototype; control 3 files; CLI: 5 files (`episode_finish_{,m2_,m3_,r2_}test.go`, `episode_markdown_test.go`) |
| V7 | Prototype compiles, is formatted, passes, is stable | `go vet ./host/daemon/`; `gofmt -l host/daemon/`; `shasum -a 256 …approve_guard_test.go`; `-count=20`; `-race -count=5`; `-v` single run (`facts.txt`, `drills.txt`) | `VET_OK`; gofmt empty; `f5df08b346dd833250ea444ed15427fb4cec1fe818f05273e0fde5187134d2e5`, 61 lines; `ok 0.621s`; `ok 2.017s`; `--- PASS (0.02s)` |
| V8 | Mutation drills | `bash verification/world-iter335/designer/drills.sh \| tee …/drills.txt` | §3 table, verbatim lines; prod restored `4e320c537e68`, test restored `f5df08b346dd`; final `git status --porcelain`: `?? host/daemon/episode_finish_approve_guard_test.go`, `?? verification/world-iter335/` |
| V9 | Whole package with the prototype | `AILANG_BIN=$HOME/.pinned-ailang/ailang go test -count=1 -timeout 900s ./host/daemon/` (`pkg.txt`) | `ok … 282.817s`, `rc=0 secs=284` |
| V10 | Existing uses of the same observables; Finish count with the prototype | `grep -n 'replay\|d.store.GetRegistryHead' host/daemon/episode_finish_test.go`; `go test -count=1 -v -run Finish ./host/daemon/ \| grep -c '^--- PASS'` (`facts.txt`) | `:96 … GetRegistryHead(…, head); err != nil \|\| exists` (loop over `broker.ApprovalsV1`, `finishHead("ep1")`); `:113 t.Fatalf("replay: %d", w.Code)`; `7` |
