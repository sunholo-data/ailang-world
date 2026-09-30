# Sprint plan — w-invocation-reconcile-object-integrity (iteration 212)

**Authority:** World mission queue row 122 (clause-6 residual of row 106). **Parent design:** [w-transition-invocation-coordinator.md](w-transition-invocation-coordinator.md), step 1b (reconcile) and R13. The finding was filed as **R-106-12** in the mission log, not in the parent doc (V12). **Scope:** one milestone, about 0.25 day. **Landing:** one commit. The planner prototyped the full change in the worktree. Every gate and mutation below was run on that prototype and nothing was committed.

Run this prefix before every command. It pins ailang v0.41.0 and keeps TMPDIR outside the repo:

```sh
source /Users/voightkampff/.ailang/state/mission-world-iter212-evidence/env.sh
```

## 1. Problem

`Coordinator.committed` (base `c8da93f`, `host/coordinator/coordinator.go:163-190`) answers a resent task id from the journal. It reads the record object at `rc.Intent.TransitionRef` (`:164`), the output object at `rec.Output` (`:176`) and the world row at `rc.Intent.WorldRef` (`:180`), then returns `Result{Reconciled: true}`. **It re-hashes no payload.** The forward path does re-hash the transition source and refuses a mismatch (`:253-254`, `SourceError{Kind:"corrupt"}`) (V1).

The store does not make up the difference:

- `store.GetObject` returns the stored payload with `Hash: ref`, echoing the ref it was asked for (`host/store/store.go:607-608`, V2).
- `store.GetWorld` echoes `Ref: ref` in the same way (`store.go:682`, V2).
- `store.Commit` "verifies neither EntryHash nor NextWorld.Ref" (`host/store/durable.go:135`, V3).

A damaged record, output or world row is therefore reported to the client as the committed result.

## 2. The fix (M1, prototyped)

Add five checks in `committed`, in read order. Each check runs before the next row is read, and all of them run before any result is built:

| # | Check | Refusal |
|---|---|---|
| C1 | record payload hashes to `rc.Intent.TransitionRef`, under the ref's own algorithm | `IntegrityError{Object:"record"}` |
| C2 | the decoded record names this invocation (`rec.InvocationID == rc.InvocationID`) | `IntegrityError{Object:"record"}` |
| C3 | output payload hashes to `rec.Output` | `IntegrityError{Object:"output"}` |
| C4 | the world row re-derives the ref the journal names (`worldRef(w) == rc.Intent.WorldRef`) | `IntegrityError{Object:"world"}` |
| C5 | the world's state root is the recorded output (`w.StateRoot == outRef`), which is the `applyRevision` law | `IntegrityError{Object:"world"}` |

Why each check beyond the row text (C1+C3) is needed. These answer the directive's UNVERIFIED items, measured:

- **C2:** a journal row can name the wrong record even when that record's hash is intact. The test `receipt_names_other_record` builds this case, and without C2 it is caught only by accident, as a "world" refusal. MUT-REC-ID-BIND-DROP proves C2 is load-bearing.
- **C4:** `GetWorld` echoes the ref (V2), so `w.Ref` proves nothing. Re-deriving the ref from the row is the only way to verify the reported `EntryIndex` (`w.Revision`). This is sound only because "a2a:" invocation ids are coordinator-exclusive. `/v1/commit` accepts only "rest:" ids (V4), so every world reconciled here was planned by `planInvocation`. The world-ref derivation is moved out of `planInvocation` into a shared `worldRef(store.World)` (`plan.go:49`, used at `:102`), so the forward and reconcile paths cannot drift apart.
- **C5:** this binds the world to the record. Without it, a journal that names another invocation's (self-consistent) world passes C4. That case is `receipt_names_other_world`, and MUT-WORLD-STATEROOT-DROP proves C5 is load-bearing.

Absent rows and undecodable payloads keep their existing untyped errors. This row does not change them.

### Error type and wire mapping (decision, with evidence)

- **Type:** a new `*coordinator.IntegrityError{InvocationID, Object}` (`errors.go:92`). Its message is `coordinator: reconcile <id>: recorded <object> does not match its reference`.
- **Why not reuse `SourceError`:** its message names "transition source", which would be false here.
- **Wire:** `projection.dispatchError` gets **no new case**. The error falls to `default:` → `-32603 "transition invocation is not available in this daemon"` (`projection.go:434`, `:87`). This is the class R6's `SourceError{corrupt}` already takes. `dispatchError` names neither type (count 0, with the positive control `NotCommittedError` at count 1, V5).
- **Why the existing wire classes would be wrong:**
  - "outcome not confirmed; resend the same task id" is false: a resend meets the same damaged rows every time.
  - "not committed; send a new task id" is false: the commit landed.
- **Test pin:** two new `TestA2ADispatch` table rows lock this mapping: `R6_corrupt_source`, which pins the precedent that was previously covered only by an untyped stand-in, and `R13_reconcile_integrity`.

The public wire contract is **unchanged**. Adding the exported Go type is additive.

## 3. Acceptance criteria

| AC | Criterion | Command |
|---|---|---|
| AC1 | damaged record, output or world is refused with the named `Object`; no reconciled result is returned; nothing re-executes; the head does not move | `go test -count=1 -run '^TestReconcileRefusesDamagedRecord$' -v ./host/coordinator/` |
| AC2 | an untampered resend still reconciles to the first call's exact result (control, plus the existing R13/R16) | the AC1 command (control subtest), plus the command below the table |
| AC3 | wire mapping is `-32603` not-available for `IntegrityError` and `SourceError` | `go test -count=1 -run '^TestA2ADispatch$' ./host/projection/` |
| AC4 | every mutation in §4 is KILLED by its named test, and the tree is byte-identical afterwards | harness `…/mut/run_mutations.py` |
| AC5 | the gates are green: vet, gofmt, race on coordinator/projection/daemon, the full suite, verify_ail | §5 |

AC2 command (outside the table, so the alternation is a plain RE2 pipe). A run with zero `=== RUN` lines is not green:

```sh
go test -count=1 -v -run 'TestDispatchRefuses/(R13|R16)' ./host/coordinator/ | grep -c '=== RUN'
```

### Tests (all in the diff)

`TestReconcileRefusesDamagedRecord` in `host/coordinator/coordinator_test.go` sets up t1, t2 and, in one subtest, t3, each committed through the real store. A `tamperStore` wrapper then damages the reconcile reads. Its hook counts every tamper hit, and a subtest whose tamper never fired fails ("fixture broken"). Each subtest asserts:

- the error type, its `Object` and its `InvocationID`;
- no `Reconciled` and no `OutputBytes`;
- a `mustNotRun` runner (nothing re-executes);
- the head is unchanged.

| Subtest | Damage | Kills |
|---|---|---|
| `control_untampered_reconciles` | none: must return t1's exact result | MUT-OUT-WRONGREF-SWAP, MUT-WORLD-WRONGREF-RECORD (any always-refuse) |
| `record_payload_corrupt` | record payload byte flipped | MUT-REC-REHASH-DROP, MUT-REC-WRONGREF-SELF, MUT-HASHES-TAUTOLOGY, MUT-REC-UNTYPED |
| `output_payload_corrupt` | output payload byte flipped | MUT-OUT-REHASH-DROP, MUT-HASHES-TAUTOLOGY, MUT-OUT-LABEL |
| `receipt_names_other_record` | receipt's TransitionRef → t2's intact record | MUT-REC-ID-BIND-DROP |
| `receipt_names_other_world` | receipt's WorldRef → t2's intact world (different state root) | MUT-WORLD-STATEROOT-DROP |
| `world_row_corrupt` | world row Revision+1 | MUT-WORLD-REDERIVE-DROP, MUT-WORLDREF-REVISION-UNBOUND |
| `world_row_of_other_invocation` | GetWorld returns t3's self-consistent row (same output, so same state root) | MUT-WORLD-WRONGREF |
| `TestA2ADispatch/R13_reconcile_integrity` (projection) | — | MUT-WIRE-RESEND |

## 4. Mutation table (run on the prototype)

Harness: `/Users/voightkampff/.ailang/state/mission-world-iter212-evidence/mut/run_mutations.py`. It applies one mutation, runs the named test only (`-count=1 -v`), restores from a byte copy, and compares sha1 before and after. Log: `…/mut/mutations.log`. After the whole run, the tree matched the pre-harness sha1 set (`TREE-IDENTICAL-AFTER-HARNESS`).

| Mutation | File: change | Named test | Verdict | Byte-identical after revert |
|---|---|---|---|---|
| MUT-REC-REHASH-DROP | coordinator.go C1 `if !hashes(rc.Intent.TransitionRef…` → `if false && …` | record_payload_corrupt | **KILLED** | yes |
| MUT-OUT-REHASH-DROP | C3 → `if false && …` | output_payload_corrupt | **KILLED** | yes |
| MUT-REC-WRONGREF-SELF | C1 compares against `hashref.SumSHA256(recObj.Payload)` (a self-referential tautology) | record_payload_corrupt | **KILLED** | yes |
| MUT-OUT-WRONGREF-SWAP | C3 compares against `rc.Intent.TransitionRef` | control_untampered_reconciles | **KILLED** | yes |
| MUT-WORLD-WRONGREF | C4 compares against `w.Ref` (the store's echo) | world_row_of_other_invocation | **KILLED** | yes |
| MUT-WORLD-WRONGREF-RECORD | C4 compares against `rc.Intent.TransitionRef` | control_untampered_reconciles | **KILLED** | yes |
| MUT-HASHES-TAUTOLOGY | `hashes`: `err == nil && sum == ref` → `\|\|` | record_/output_payload_corrupt | **KILLED** | yes |
| MUT-REC-ID-BIND-DROP | C2 → `if false && …` | receipt_names_other_record | **KILLED** | yes |
| MUT-WORLD-REDERIVE-DROP | C4 → `if false && …` | world_row_corrupt | **KILLED** | yes |
| MUT-WORLD-STATEROOT-DROP | C5 → `if false && …` | receipt_names_other_world | **KILLED** | yes |
| MUT-OUT-LABEL | output refusal labelled `"record"` | output_payload_corrupt | **KILLED** | yes |
| MUT-REC-UNTYPED | C1 returns an untyped `fmt.Errorf` | record_payload_corrupt | **KILLED** | yes |
| MUT-WORLDREF-REVISION-UNBOUND | plan.go `worldRef`: `Revision: w.Revision` → `0` | world_row_corrupt | **KILLED** | yes |
| MUT-WIRE-RESEND | projection.go: map `IntegrityError` into the "resend the same task id" case | TestA2ADispatch/R13_reconcile_integrity | **KILLED** | yes |

**14/14 KILLED.** History: the first draft of the test table had no `world_row_of_other_invocation` subtest, and MUT-WORLD-WRONGREF **SURVIVED** against `world_row_corrupt`. The survivor was equivalent under the real store, because `GetWorld` echoes `Ref` (V2) and a Revision tamper keeps the echoed ref. The wrong-key subtest was added and the mutant is now KILLED (`…/mut/history.txt`).

## 5. Gates (prototype, logs in the evidence dir)

| Gate | Result | Log |
|---|---|---|
| `go vet ./...` | rc=0 | `gate_vet.log` |
| `gofmt -l host cmd` | empty (0 lines) | `gate_gofmt.log` |
| `go test -race -count=1 -v ./host/coordinator/ ./host/projection/ ./host/daemon/` | rc=0; `=== RUN` 65 / 69 / 298 (432 total); 0 `--- FAIL` | `gate_race_targeted.log` |
| `go test -count=1 ./...` | rc=0; 24 `ok`, 0 `FAIL` (no load-flake re-runs needed) | `gate_full.log` |
| `bash scripts/verify_ail.sh` | rc=0; 16 identities, 40 named tests, 9/9 package steps, PUB011 ratchet | `gate_verify_ail.log` |

Prototype diff (`git diff --numstat`): `coordinator.go` +29/−1, `coordinator_test.go` +164/−0, `errors.go` +13/−0, `plan.go` +8/−3, `projection_test.go` +4/−0. In total 5 files, +218/−4. `projection.go` is **unchanged**.

## 6. OPEN (not settled here)

- **OPEN-1 (wire, human judgement):** should store damage found during reconcile have its **own** JSON-RPC message, for example "stored invocation record is damaged", instead of the generic not-available one? That would change the public wire contract, and R6's corrupt source would presumably move with it. The prototype keeps the existing class. The table pin makes any change deliberate.
- **OPEN-2 (policy):** should a detected integrity failure quarantine the store (`store.Quarantine`, used today only by the daemon's startup budget abort, `daemon.go:679`)? The coordinator's `Store` seam does not expose quarantine, and R6 does not quarantine either. Not done.
- **Not covered:** absent record/output/world rows during reconcile stay untyped errors, which map to the same wire class. Typing them is out of this row's text.

## 7. Verification Log

The raw transcript is `/Users/voightkampff/.ailang/state/mission-world-iter212-evidence/verification_log_raw.txt`.

| V | Claim | Command → observed |
|---|---|---|
| V1 | base `committed` reads three rows with no re-hash; the source path re-hashes | `git show c8da93f:host/coordinator/coordinator.go \| grep -n 'func (c \*Coordinator) committed\|GetObject(ctx, rc.Intent.TransitionRef)\|GetObject(ctx, outRef)\|GetWorld(ctx, rc.Intent.WorldRef)\|SourceError{Kind: "corrupt"}\|hashref.Sum('` → `163`, `164`, `176`, `180`, `253 … hashref.Sum(d.TransitionFn.Algo()…`, `254 … SourceError{Kind: "corrupt"}`. `hashref.Sum(` occurs only at 253, which is the positive control that the grep can see a re-hash. |
| V2 | GetObject/GetWorld echo the requested ref | `sed -n 607,608p host/store/store.go; sed -n 682p …` → `return Object{` / `Hash: ref,` / `return World{Ref: ref, …}` |
| V3 | the store verifies no world ref | `grep -n "verifies neither" host/store/durable.go` → `135: … the store verifies neither EntryHash nor NextWorld.Ref;` |
| V4 | "a2a:" is coordinator-exclusive | `grep -n '"a2a:"\|restInvocationPrefix =' host/coordinator/coordinator.go host/daemon/handlers.go` → `coordinator.go:101 … return "a2a:" + …`; `handlers.go:353 // … cannot collide with the coordinator's "a2a:" ids`; `:354 const restInvocationPrefix = "rest:"` |
| V5 | dispatchError has no case for SourceError/IntegrityError, so both fall to default | `grep -c 'coordinator.SourceError\|coordinator.IntegrityError' host/projection/projection.go` → `0`; control `grep -c 'coordinator.NotCommittedError' …` → `1`; `grep -n 'return codeInternal, notAvailableMessage'` → `434`; `grep -n 'notAvailableMessage *='` → `87: … "transition invocation is not available in this daemon"` |
| V6 | the prototype's checks are at the stated lines | `grep -n 'hashes(\|worldRef(w)\|w.StateRoot != outRef\|rec.InvocationID != rc.InvocationID\|IntegrityError{' host/coordinator/coordinator.go` → `163 func hashes`, `178/179` C1, `185/186` C2, `196/197` C3, `206/207` C4, `209/210` C5 |
| V7 | shared world-ref derivation | `grep -n 'func worldRef\|next.Ref = worldRef' host/coordinator/plan.go` → `49`, `102` |
| V8 | new error type | `grep -n 'type IntegrityError' host/coordinator/errors.go` → `92` |
| V9 | reconcile is entered only from Dispatch step 1b | `grep -n 'committed(ctx, rc)' host/coordinator/coordinator.go` → `243` (single site) |
| V10 | quarantine has one production caller | `grep -rn 'Quarantine(' --include='*.go' host \| grep -v _test` → `daemon/daemon.go:679 store.Quarantine(d.store)`, plus the definition `store.go:251` (all other hits are `checkQuarantine`) |
| V11 | pinned binary | `$HOME/.pinned-ailang/ailang --version` → `AILANG v0.41.0` |
| V12 | R-106-12 lives in the mission log, not the parent doc | `grep -c R-106-12 design_docs/planned/w-transition-invocation-coordinator.md` → `0`; control `grep -c R-106-11 …` → `8`; `grep -n R-106-12 design_docs/world-mission-log.md` → `866: - \`committed()\` does not re-hash the objects it reads back … → R-106-12`, and `1309` |
| V13 | charter row | `grep -n '^122\. \*\*w-invocation-reconcile' design_docs/world-mission.md` → `1909` |
