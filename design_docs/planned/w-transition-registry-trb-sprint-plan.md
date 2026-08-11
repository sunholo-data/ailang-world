# Sprint plan — `TR.B` (queue item 11 `w-transition-registry`)

**Milestone**: `TR.B` — *capability snapshot and declared-effect confinement*
**Status**: PLANNED · **SPLIT RECOMMENDED** into `TR.B1` + `TR.B2` · **TR.B ONLY** (TR.A is landed; TR.C is out of scope)
**Design doc**: [`design_docs/planned/w-transition-registry.md`](w-transition-registry.md)
**Sibling plan**: [`w-transition-registry-tra-sprint-plan.md`](w-transition-registry-tra-sprint-plan.md) (TR.A1 + TR.A2, both landed)
**Base**: `dev` @ `66a1d63`, clean tree, CI green both jobs
**Planner**: mission-control iteration 73, opus, first-party measurement on this rig
**Executor**: sandboxed worktree, **NO git write permission**.
**THE CONTROLLER MAKES ALL COMMITS.** The executor never runs `git commit`, `git add`, `git push`,
`git checkout`, `git stash`, `git restore`, or `gh pr`. Restores are `cp` from a backup — see §7.

---

## 0. Planner's first-party verification of every premise

Every controller-supplied premise and every design-doc premise this plan leans on was re-run on this
rig at `66a1d63` before it entered the plan. Negative results are paired with a known-positive
control **in the same call**, per the standing rule that an empty result is a claim, not a fact.

| # | Premise | Verdict | Command → observed output |
|---|---|---|---|
| P1 | AC5/AC6/AC7 are all at their base-tolerant `count=0` arm | **CONFIRMED** | the three exact AC commands → `AC5_count=0`, `AC6_count=0`, `AC7_count=0`; **two** known-positive controls in the same call: the landed broker replay pair → `2`, the landed AC2 transitionreg trio → `3` |
| P2 | AC11 base count is 1 and must not move | **CONFIRMED** | exact AC11 command → `AC11_count=1` |
| P3 | `broker.Decide(c Capability, r EffectRequest) Decision` exists | **CONFIRMED** | `host/broker/decide.go:40`; the four denial labels are `decide.go:8-11` |
| P4 | `broker.Allows` / `CapabilitySnapshot` / `Bind` / `Requirement` do not exist | **CONFIRMED** | `grep -rn 'func Allows' host/ cmd/ \| wc -l` → `0`; `CapabilitySnapshot` → `0`; `func Bind` → `0`; `type Requirement` → `0`. Controls in the same call: `func Decide` → `2`, `func NewSession` → `1` |
| P5 | `broker.Session` `:46`, `NewSession` `:58`, `Session.Invoke` `:126` | **CONFIRMED** | `host/broker/broker.go` at exactly those lines |
| P6 | production `.Invoke(` call sites = exactly 3, all `publish_op.go` | **CONFIRMED** | `grep -rn '\.Invoke(' --include='*.go' host/ cmd/ \| grep -v _test.go` → `publish_op.go:135,162,279`; `_test.go` control arm → **83** |
| P7 | freshness of the doc's V-rows about `host/broker` | **CONFIRMED FRESH** | `git diff --name-only b0f323a HEAD \| grep -v '^design_docs/'` → **14** files (control incl. docs: **23**); `grep -c '^host/broker/'` of that list → **0** |
| P8 | nothing imports `host/transitionreg` yet | **CONFIRMED** | `grep -rn 'host/transitionreg' --include='*.go' host/ cmd/ \| grep -v '^host/transitionreg/'` → 0 lines, rc=1; control `host/broker` importers → **10** files |
| P9 | `host/transitionreg` has no package-scope `type Registry` (naming freeze) | **CONFIRMED** | `grep -c '^type Registry' host/transitionreg/*.go` → `0,0,0`; control `grep -c '^type Registry' host/broker/broker.go` → **1** |
| P10 | TR.A's store CAS primitive is landed and callable | **CONFIRMED** | `host/store/store.go:648 type RegistryCASConflict`, `:666 func IsRegistryCASConflict`, `:674 func (s *Store) CompareAndSetRegistryHead`; `store.TransitionRegistryV1` at `store.go:86` |
| P11 | grant debit happens at exactly 3 sites | **CONFIRMED** | `grep -rn 'grants\[.*\]\.Budget = ' host/broker/ \| grep -v _test.go` → `broker.go:208` (live allowed), `:367` (replay failed), `:380` (replay succeeded) |
| P12 | `go build ./...` and `go vet ./host/...` are **both clean at base** | **CONFIRMED** | `GOTOOLCHAIN=go1.25.6 go build ./...` → rc=0; `GOTOOLCHAIN=go1.25.6 go vet ./host/...` → **rc=0, zero findings** |
| P13 | AC8 (replay) base | **CONFIRMED** | `AC8_count=2`; `ok host/replay 2.054s`, rc=0 |
| P14 | AC9 repaired form, 4/11/14 | **CONFIRMED** | `rc=0 modules11=1 tests14=1 steps9=1` |
| P15 | `scripts/verify_go.sh` at base | **GREEN, rc=0** | `VERIFY_GO_BASE_RC=0`; plain leg `ok host/broker 48.814s`, `ok host/transitionreg 3.479s`; race leg `ok host/broker 92.337s`, `ok host/transitionreg 4.968s`; the **2** `DATA RACE` lines are the gate's own race-detector known-positive control |
| P16 | TR.C's hold surface is currently empty (nothing to break yet) | **CONFIRMED** | production `broker.Session` outside `host/broker` → **0**; `broker.NewSession(`/`NewReplaySession(` production callers outside → **0**; control: the same grep including `_test.go` → **1** |
| P17 | `rg` is not a binary | **CONFIRMED (inherited, re-checked)** | `whence -p rg` → `NO-BINARY`. **No command in this plan contains `rg`.** |

### Six measurements the controller did not make, five of which change the plan

#### (i) **TR.B as designed will make TR.C's AC11 impossible to satisfy, unless TR.B extracts an unexported `invoke`.**

This is the largest finding in the plan and it wants a controller decision (§8 R1).

TR.C (doc lines 398–410) pins **two** assertions: outside `host/broker`, *any* `Invoke` selector call
is RED; inside `host/broker`, the exemption count is **exactly 3**, at named identities in
`publish_op.go`. Measured at base, production `.Invoke(` sites are exactly those 3 (P6).

Decision 5 requires TR.B to supply a descriptor-bound invoker that "refuses an undeclared request
**before calling `Session.Invoke`**; declared requests still pass through `Session.Invoke`". Written
naively, that guard is a **fourth** production `.Invoke(` call site:

- put it in `host/transitionreg` → it is an `Invoke` selector call **outside** `host/broker` → TR.C's
  first assertion is RED, and the only repairs are a package wildcard (which the doc explicitly
  forbids: *"an exact legacy carve-out, not a package wildcard"*) or raising the exemption count
  (which the doc names as the mutation `MUT-BINDING-RAISE-COUNT`);
- put it in `host/broker` → the pinned exemption count of 3 becomes 4 → TR.C's second assertion is RED
  and is indistinguishable from `MUT-BINDING-FOURTH-INVOKE`.

Either way **an earlier milestone reds a later milestone's criteria**, which is the exact property the
TR.A split was designed to preserve, and which the doc claims for all three milestones
(*"All three milestones are independently mergeable"*, line 423).

**Freeze that removes the collision** (binding on T3):

1. Rename the body of `(*Session).Invoke` to unexported `(*Session).invoke` — same signature, same
   locking, no behaviour change. `Invoke` becomes a one-line wrapper `return s.invoke(ctx, req, payload)`.
2. The confined seam calls `b.s.invoke(...)`. An AST detector keyed on the selector name `Invoke`
   does not match `invoke`, and a `grep '\.Invoke('` does not match it either.
3. Production `.Invoke(` sites therefore stay at **exactly 3, all in `publish_op.go`**, and TR.C lands
   unchanged. This is added to §5 as a **hold criterion** of TR.B (`AC-INVOKE3`), because the doc's AC
   list has no criterion that protects it.

A second-order consequence, also frozen: the confined seam's exported method must **not** be named
`Invoke` and must **not** expose `*broker.Session` to its caller, or TR.C's outside-`host/broker`
assertion reds `host/transitionreg`'s production code. Hence `Session.Bind(Manifest) (*BoundInvoker, error)`
and `(*BoundInvoker).Request(...)` in `host/broker`, consumed in `host/transitionreg` through a
locally-declared interface (`Binder`) so the identifier `broker.Session` never appears there. §3 T3/T4.

#### (ii) **`go test ./host/broker` is NOT deterministically green at base. Measured flake rate 2/11.**

`TestHandlerTimeoutKillsTheWholeProcessGroup` (`handlers_test.go:744`) fails intermittently on this
rig with *"Invoke took 5.24s for a 100ms bound: the kill missed the forked grandchild"*.

Measured, isolated, nothing else running:

| arm | runs | result |
|---|---:|---|
| `-run 'TestHandlerTimeoutKillsTheWholeProcessGroup$'` | 3 then 8 | **2 FAIL / 11** (≈18 %) |
| known-positive control `-run 'TestCanonicalLabelSet$'` | 3 | 3 PASS / 3 — the instrument is not just "the rig is red" |
| whole package, serial, `AILANG_BIN` set | 1 | `ok … 35.389s` |
| whole package, **under concurrent load** | 1 | FAIL on this same test |
| `scripts/verify_go.sh` (both legs, serial) | 1 | rc=0 |

**This is fatal to a naive mutation sweep.** §4.1 step 6 ("rc≠0 = kill") and step 7 (inverse arm must
be rc=0) are both corrupted by an 18 % random red in the package under mutation: it fakes kills and it
falsifies inverse arms. Every broker mutation arm in §4 therefore carries a mandatory
`-skip 'TestHandlerTimeoutKillsTheWholeProcessGroup$'` and **the FAIL line must be read, never the
exit code**. The flake is out of scope for TR.B (it is a pre-existing `host/broker` defect); it is
reported to the controller as a queue candidate in §8 R5.

#### (iii) **`go vet` catches the exact defect class TR.B's new code is most likely to contain, and `go test` does not. Both arms measured.**

TR.B returns a value snapshot from a `sync.Mutex`-bearing `*Session`. The single most likely slip is a
value receiver. I armed exactly that shape at base and ran both instruments on the identical tree:

```
host/broker/zz_probe.go:6:9: probeCapabilitySnapshot passes lock by value:
    …/host/broker.Session contains sync.Mutex
```

| instrument | same tree, probe armed | verdict |
|---|---|---|
| `GOTOOLCHAIN=go1.25.6 go vet ./host/broker` | **rc=1**, copylocks finding printed | sees it |
| `GOTOOLCHAIN=go1.25.6 go test ./host/broker -run 'TestNoSuchThing$'` | **rc=0**, `ok … [no tests to run]` | blind |

Probe removed; `git status --porcelain` → empty. With `go vet ./host/...` **rc=0 at base** (P12), the
gate is both non-vacuous and attributable: any TR.B vet finding is TR.B's. `go vet ./host/...` rc=0 is
an exit gate on **every** task in §3.

#### (iv) **`MUT-ALLOW-NAME/SCOPE/EXPIRED/BUDGET` have an inverse arm that is unsatisfiable by construction — and that is the *point*, not a defect.**

`Decide`'s four denial branches already have a landed co-detector: `decide_test.go:66-71`
(`TestSketchRows/line_227…233_decideLabel`) pins all four labels directly. So mutating a denial branch
**inside `Decide`** reds both the new `TestAllowsUsesDecideAllFourDenials` **and** the landed
`TestSketchRows`, and §4.1 step 7's green-package inverse arm cannot be satisfied without weakening a
landed test. This is the TR.A2 finding repeating exactly.

Do **not** weaken `TestSketchRows`. Instead each `MUT-ALLOW-*` is run as **two** arms:

- **arm (a) — in `Allows`** (e.g. `if d.Label == LabelDeniedScope { d.Allowed = true }`): kill =
  `TestAllowsUsesDecideAllFourDenials/denied_scope`; inverse `-skip` arm **must** be rc=0. This is the
  ordinary non-vacuity proof.
- **arm (b) — in `Decide`**: kill = the new test **and** `TestSketchRows`; the inverse arm is recorded
  **UNSATISFIABLE BY CONSTRUCTION** and not run to green. Arm (b) is *required*, because co-detection
  is the only positive evidence that `Allows` **delegates** to `Decide` rather than restating the four
  comparisons. An `Allows` that restated them correctly would pass arm (a) and survive arm (b). Decision 5
  forbids the second policy engine; arm (b) is how that is measured rather than asserted.

#### (v) **`TestAllowsUsesDecideAllFourDenials` is a near-duplicate of a landed test unless it is written as a delegation test.**

`decide_test.go` already enumerates the four denial labels as a table over `Decide`. A new test that
re-tables the same four rows over `Allows` is the "same concept duplicated" the Systemic-Issue Audit
forbids, and adds no information. The AC5 test must instead pin what is *new*: snapshot-sourced grants,
multi-grant ranked selection, the empty-grants arm, and `snap.Now` as the clock. See T2 for the
required subtest list.

#### (vi) **`grep -c 'return .*fmt.Errorf(.*%w'` would have missed 64 of TR.A's own 94 refusal returns. Calibrated on the landed diff.**

The controller's constraint is confirmed with a number. Over TR.A's own landed diff
(`git diff --unified=0 b0f323a..HEAD -- host/transitionreg host/store`):

| cut | count |
|---|---:|
| `%w`-only (the WRONG cut) | **30** |
| `(fmt\.Errorf\|errors\.New)\(` | **94** |
| typed struct returns `return .*&[A-Z]…{` | **2** |
| union of the last two | **96** |
| negative control, same cut over `design_docs/` | **0** |

TR.B returns typed error structs (`*DenialError` is landed; `*UndeclaredEffectError` and
`*AccessDeniedError` are new), so the union cut is the floor. §4.3 uses the union **and** says plainly
that the grep is a floor, not the enumeration: the enumeration is by reading TR.B's own diff.

### One non-finding, recorded because a reviewer will ask

**No Go `Proposal` type exists to reuse.** `grep -rn 'type Proposal'` over `host/ cmd/` → 0; the
`TransitionFn` hits are `host/store/journal.go:34`, `host/store/store.go:113/134`,
`host/replay/replay.go:69` and the two daemon DTOs — all *recorded-episode* structures, none carrying
`expectedEffects`/`requiredCaps`. The proposal shape lives in `world/types.ail:45,48,49` (V5). So
`transitionreg.Proposal` (T4) is a new host-side mirror of a pure-World shape, not a duplicated
concept. It must stay a plain agreement-checking value: it performs no authority decision.

---

## 1. What TR.B is, and what this plan explicitly does NOT do

**In scope** (the doc's Files table, TR.B rows, plus two files the doc's table omits):

| File | Purpose | Milestone half |
|---|---|---|
| `host/broker/decide.go` | `Requirement`, `Allows`, shared `decideOver` | TR.B1 |
| `host/broker/broker.go` | `epoch` field, `CapabilitySnapshot`, single `debitGrant` site, `Invoke`→`invoke` extraction | TR.B1 |
| `host/broker/confined.go` **(new; doc's table omits it)** | `Manifest`, `Session.Bind`, `BoundInvoker.Request`, `UndeclaredEffectError` | TR.B1 |
| `host/broker/broker_test.go` | AC5's two tests + the confined-seam tests | TR.B1 |
| `host/transitionreg/bind.go` **(new; doc's table folds this into `transitionreg.go`)** | `Bind`, `Bound`, `Proposal`, `Request`, `NewRequest` | TR.B2 |
| `host/transitionreg/transitionreg_test.go` | AC6's three + AC7's three tests | TR.B2 |
| `design_docs/planned/w-transition-registry.md` | **T7 only**: the zero-tolerance AC5/AC6/AC7 activation | TR.B2 |

**Explicitly NOT in this sprint** — if the executor touches any of these, the milestone is wrong:

- no `host/broker/invoke_boundary_test.go` and no AST binding assertion — that is **TR.C**, and AC11's
  base count of 1 must still read 1 at merge;
- no change to `Decide`'s four comparisons, to `denialRank`, or to `TestSketchRows`. `Allows` **delegates**;
- no new production `.Invoke(` call site anywhere (§0(i), `AC-INVOKE3`);
- no `host/store` change of any kind — TR.A's CAS primitive is complete and TR.B only calls it;
- no change to `host/transitionreg/codec.go`, to the frozen `InterfaceHashV1`, or to any golden byte:
  TR.B "adds authority integration **without changing its wire format**" (doc line 424). The
  `Descriptor`/`Revision` structs are read-only inputs to TR.B;
- no `world/*.ail` addition, no store table, no REST route, no projection package, no MCP/A2A code;
- no CI workflow change; no `host/replay` change (AC8 must not move);
- no fix for the `TestHandlerTimeoutKillsTheWholeProcessGroup` flake (§0(ii)) — it is reported, not
  repaired, and TR.B must not silence it.

---

## 2. AC reconciliation — plan tasks diffed against the doc's AC list

The doc has 11 ACs. TR.B owns 3, holds 5, and must leave 3 at their TR.A-activated values.

| AC | Owner | TR.B duty | Closing task |
|---|---|---|---|
| AC1 identity/codec/schema | TR.A (landed) | **HOLD at exactly 3** — TR.B must not touch the codec | verified every task |
| AC2 eager snapshot | TR.A (landed) | **HOLD at exactly 3** | verified every task |
| AC3 CAS publication | TR.A (landed) | **HOLD at exactly 4** | verified every task |
| AC4 store CAS | TR.A (landed) | **HOLD at exactly 2** | verified every task |
| **AC5** capability snapshot + four denial arms | **TR.B1** | **CLOSES** — activate to exactly **2** | T1, T2 → activated T7 |
| **AC6** declaration-honesty mechanism | **TR.B2** | **CLOSES** — activate to exactly **3** | T3, T4 → activated T7 |
| **AC7** two-session fixture | **TR.B2** | **CLOSES** — activate to exactly **3** | T5 → activated T7 |
| AC8 replay hash-pinned | TR.B holds | **MUST NOT MOVE** — no `transitionreg`/`broker` import added to production replay | verified every task |
| AC9 AILANG gate totals | TR.B holds | **MUST NOT MOVE** — repaired form, 4/11/14 | verified every task |
| AC10 build + focused packages | TR.B holds | **MUST NOT MOVE** — `count=1`, whole `./host/transitionreg` PASS | verified every task |
| AC11 structural binding boundary | TR.C | untouched; base count **1** must remain 1 | — |

### Delta: three criteria the doc's AC list does not have, which TR.B needs

| new criterion | why the doc's list cannot cover it | form |
|---|---|---|
| **`AC-INVOKE3`** — production `.Invoke(` sites stay exactly 3, all in `publish_op.go` | The doc's AC11 pins this, but AC11 is TR.C's criterion and is *absent* at TR.B merge time (count=1). Nothing in AC1–AC10 notices a 4th call site. §0(i) shows TR.B is the milestone most likely to add one. | §5, run every task |
| **`AC-VET`** — `go vet ./host/...` rc=0 | `go test`'s default vet subset excludes `copylocks`, and `verify_go.sh` never calls `go vet`; §0(iii) measures both arms. TR.B's `CapabilitySnapshot` is precisely the shape that trips it. | §5, run every task |
| **`AC-NOFLAKE`** — every broker gate result is attributed by FAIL line, with the known flake skipped | The doc assumes a deterministic package. §0(ii) measures 2/11. | §4.1 step 6a |

### Doc contradictions and gaps found (reported, not smoothed over)

1. **Lines 398–410 (TR.C) vs Decision 5 (line 278–283)** — the descriptor-bound invoker's call into
   `Session.Invoke` is a fourth production `Invoke` selector call under *either* placement, so TR.B as
   designed reds TR.C. §0(i). *Plan resolves it with the unexported `invoke` extraction and adds
   `AC-INVOKE3`. This is the one item I want ratified before execution (§8 R1).*
2. **Files table (lines 436–438) is incomplete for TR.B.** It names `broker.go`, `decide.go`,
   `broker_test.go` only. The confined seam needs its own file, and the entire AC6/AC7 surface
   (`Bind`, `Bound`, `Proposal`, `Request`) lands in `host/transitionreg`, whose TR.B rows the table
   marks only as `transitionreg_test.go` (line 433) — i.e. **the table has TR.B writing AC6/AC7 tests
   with no TR.B production file to test.** *Plan adds `host/broker/confined.go` and
   `host/transitionreg/bind.go` (§1).*
3. **Line 273 spells `broker.Allows(snapshot, requirement)` but no `requirement` type exists in
   `host/broker`,** and `transitionreg.EffectRequirement` cannot be used because `host/broker` must not
   import `host/transitionreg` (that is the cycle: `transitionreg` imports `broker` for
   `EffectRequest`/`Manifest`). *Plan freezes `broker.Requirement{Effect, Scope string; Cost int64}`
   with a one-line adapter in `transitionreg`. §3 T2.*
4. **Line 272 leaves the epoch's increment rule ambiguous** — "returns an immutable, monotonically
   incremented epoch" is satisfied by incrementing on *every snapshot call*, which would make
   `MUT-CAPS-STATIC-EPOCH` killable by a version counter that carries no information about the ledger.
   Line 387 ("epoch increments **on debit**") is the binding reading. *Plan freezes: the epoch
   increments if and only if a grant budget is mutated, at the single extracted `debitGrant` site, and
   adds the `denied_invoke_does_not_increment_epoch` subtest that a per-call counter would fail.*
5. **Line 596 `MUT-CAPS-ALIAS` names one branch where the mechanism has two aliasing seams** — the
   session→snapshot copy and the snapshot→caller copy. Neutering one leaves the other. *Plan splits it
   into two arms.*
6. **Line 621 ("Eleven frozen refusal branches … the TR.A sprint plan adds mutations for all eleven")
   is now stale in a way that matters to TR.B**: the doc's rule-3j sentence is scoped to Decisions 1/2/4,
   i.e. to TR.A. TR.B's own refusal branches (manifest membership, proposal agreement, access denial
   propagation, request capture) are **not** enumerated anywhere. This is precisely the TR.A2 failure —
   an audit anchored to the doc's decision list cannot see the branches the sprint itself writes.
   *Plan enumerates them from TR.B's own diff instead. §4.3.*
7. **Line 623 "SQLite CAS conflicts are driven by two handles/racers"** — already refuted in the TR.A
   plan §0(i) and still uncorrected in the doc. TR.B does not depend on it; noted for the record.

---

## 3. Day-by-day breakdown — 7 tasks, 7 commits, across a recommended 2-milestone split

One task = one commit boundary. **The controller makes every commit**; the executor stops at each
boundary and reports. Every task ends by re-running the **hold set**: AC1–AC4 counts, AC8, AC9-repaired,
AC10, AC11=1, `AC-INVOKE3`, `AC-VET`. *A task that closes its own AC while moving a held one has not
succeeded.*

### `TR.B1` — broker capability snapshot and delegating predicate (closes AC5)

Touches **only** `host/broker`. Leaves AC6/AC7 at `count=0`, where their base-tolerant arm keeps them
green — so TR.B1 is independently mergeable.

#### T1 — `CapabilitySnapshot` + ledger epoch  (~90 impl / ~150 test LOC)

In `host/broker/broker.go`:

```go
// Session gains ONE field:
epoch int64

// CapabilitySnapshot is an immutable copy of the session ledger at one instant.
// Now is captured so every downstream decision uses one clock reading and the
// broker still never reads a wall clock.
type CapabilitySnapshot struct {
    Epoch  int64
    Now    int64
    grants []Capability   // unexported: the snapshot cannot be mutated in place
}
func (s *Session) CapabilitySnapshot(now int64) CapabilitySnapshot  // POINTER receiver — §0(iii)
func (c CapabilitySnapshot) Grants() []Capability                   // fresh copy every call
func (c CapabilitySnapshot) Len() int
```

Implementation constraints, all load-bearing:

- **Pointer receiver on `CapabilitySnapshot()`.** A value receiver copies `sync.Mutex` and `go vet`
  reds — measured, §0(iii). `go test` does not.
- **Never call `CapabilitySnapshot` from inside a `s.mu`-held region.** `sync.Mutex` is not reentrant;
  the deadlock presents as a 10-minute `go test` panic timeout, not an error. This is the TR.B analogue
  of TR.A's §0(ii). In particular `invoke` (T3) holds the lock for its whole body.
- **Extract the three debit sites (P11: `broker.go:208,367,380`) into ONE method**
  `func (s *Session) debitGrant(i int, remaining int64) { s.grants[i].Budget = remaining; s.epoch++ }`.
  One mechanism, one mutation site. Leaving three sites means `MUT-CAPS-STATIC-EPOCH` has three
  homes and a neuter of one is invisible.
- The epoch increments **iff** a budget is mutated. A denial does not debit and must not increment.

Test `TestCapabilitySnapshotEpochAndIsolation` in `broker_test.go`, reusing the landed
`openTestStore(t)` (`broker_test.go:15`) and `echoRegistry(&counter)` (`:29`) helpers, with these named
subtests (each is a mutation target in §4):

| subtest | asserts |
|---|---|
| `fresh_session_epoch_is_zero` | a new session snapshots at epoch 0 |
| `two_snapshots_without_a_debit_share_an_epoch` | pins "on debit", kills a per-call counter |
| `allowed_invoke_increments_epoch_exactly_once` | epoch before+1 == after |
| `denied_invoke_does_not_increment_epoch` | denial path leaves the ledger version alone |
| `replay_debit_increments_epoch` | the replay debit sites are on the same mechanism |
| `snapshot_is_isolated_from_later_debit` | snapshot taken, then a debiting `Invoke`, then `Grants()` still shows the old budget — **`MUT-CAPS-ALIAS` arm (a)** |
| `grants_accessor_returns_a_fresh_copy` | mutate the returned slice, call `Grants()` again, unchanged — **`MUT-CAPS-ALIAS` arm (b)** |
| `snapshot_now_is_the_caller_supplied_reading` | `Now` round-trips; no wall clock |

**Exit gate**: `AC-VET` rc=0. Hold set green. (AC5 is not yet complete — T2 finishes it.)

#### T2 — `Requirement`, `Allows`, and the shared ranked selector  (~50 impl / ~180 test LOC) → closes AC5

In `host/broker/decide.go`:

```go
// Requirement is what a registry descriptor asks for. It carries no expiry and
// no budget: those belong to session grants, never to registry metadata.
type Requirement struct { Effect, Scope string; Cost int64 }

// decideOver is the ONE ranked-selection mechanism. Session.decide and Allows
// both call it; neither restates a comparison.
func decideOver(grants []Capability, r EffectRequest) (int, Decision)

// Allows feeds the snapshot's grants and the requirement to the landed Decide.
// An absent matching grant is a denial, not an error.
func Allows(snap CapabilitySnapshot, req Requirement) Decision
```

- `Session.decide` becomes `return decideOver(s.grants, req)` — a pure extraction, behaviour-identical.
  The landed broker suite is the regression net; run it whole (with the §0(ii) skip) at this task's exit.
- `Allows` builds `EffectRequest{Effect: req.Effect, Scope: req.Scope, Cost: req.Cost, Now: snap.Now}`
  and returns `decideOver(snap.grants, …)`'s `Decision`. **Zero comparisons of effect, scope, expiry or
  budget appear in `Allows`.** A mechanical check, run at this task's exit:
  `awk '/^func Allows/,/^}/' host/broker/decide.go | grep -cE '==|<|>|ExpiresAt|Budget'` → **0**, with the
  known-positive control `awk '/^func Decide/,/^}/' host/broker/decide.go | grep -cE '==|<|>|!'` → **≥4**
  in the same call.

Test `TestAllowsUsesDecideAllFourDenials` — **not** a re-table of `Decide` (§0(v)); it pins what is new:

| subtest | asserts |
|---|---|
| `denied_effect_name` / `denied_scope` / `denied_expired` / `denied_budget` | one grant, one requirement, exact `Label` — the four `MUT-ALLOW-*` targets |
| `allowed_returns_remaining_label` | `allowed:<remaining>` matches `Decide`'s |
| `no_grants_is_denied_effect_name` | the empty-ledger arm is a denial, not an error |
| `ranked_best_denial_across_three_grants` | multi-grant selection matches `Session.decide`'s ranking |
| `first_allowing_grant_wins` | short-circuit parity with the landed selector |
| `uses_snapshot_now_not_a_wall_clock` | same grants, two snapshots at `now=5` and `now=50`, expiry 10 → allowed then `denied:expired` |
| `label_agrees_with_decide_row_for_row` | for a table of (capability, requirement), `Allows(...).Label == Decide(...).Label` — the delegation assertion whose *teeth* are `MUT-ALLOW-*` arm (b), §0(iv) |

**Exit gate**: AC5 → `count=2`, both PASS. `AC-VET` rc=0. Hold set green.

#### T3 — the confined seam + `invoke` extraction  (~110 impl / ~160 test LOC) → makes TR.C possible

New file `host/broker/confined.go`:

```go
// Manifest is one descriptor's authority envelope, copied out of the registry.
type Manifest struct { Access Requirement; Declared []Requirement }

// BoundInvoker is the ONLY dispatch surface exported outside host/broker. Its
// method is deliberately NOT named Invoke and it does not expose *Session, so
// TR.C's AST boundary can stay an exact three-site carve-out (§0(i)).
type BoundInvoker struct { s *Session; declared []Requirement }

func (s *Session) Bind(m Manifest) (*BoundInvoker, error)
func (b *BoundInvoker) Request(ctx context.Context, req EffectRequest, payload []byte) ([]byte, hashref.HashRef, error)
func (b *BoundInvoker) Declared() []Requirement   // fresh copy

type UndeclaredEffectError struct{ Effect, Scope string; Cost int64 }
func (e *UndeclaredEffectError) Error() string   // MEASURED message, pinned by test
```

- **`Session.Invoke` body → unexported `(*Session).invoke`**, `Invoke` a one-line wrapper. §0(i).
  The lock stays inside `invoke`; `Invoke` must not lock. `Request` calls `b.s.invoke`.
- Membership is on the **whole triple** `{Effect, Scope, Cost}`, not the effect name. A declared
  `{FS.Read, /a, 1}` does not authorize `{FS.Read, /b, 1}` or `{FS.Read, /a, 99}`. Three separate
  refusal branches, three separate arms (§4.3 J3/J4).
- The guard runs **before** `invoke`, so a refused request performs zero handler calls, zero store
  writes and zero debits. The tests assert the handler counter is 0 **and** the epoch is unchanged.
- `Bind` validates the manifest: non-negative costs, duplicate-free `Declared` (§4.3 J1/J2).
- **`Request` adds no capability decision.** Authority stays `invoke`'s landed pipeline; a declared
  request with no live grant must still come back as the landed `*DenialError`.

Tests (broker-side, not AC-listed, therefore they move no count):
`TestBoundInvokerRefusesUndeclaredRequest`, `TestBoundInvokerRequestStillRunsTheLandedPipeline`,
`TestBindRefusesMalformedManifest`.

**Exit gate**: `AC-INVOKE3` → exactly 3, all `publish_op.go`. Whole `./host/broker` package green
(serial, `AILANG_BIN` set, §0(ii) skip applied). `AC-VET` rc=0. AC5 still `count=2`. Hold set green.

#### T3m — mutation sweep for TR.B1  (~16 arms, no production LOC)

Run §4.1 for every mutation whose mechanism landed in T1–T3. Deliverable is
`design_docs/verification/w-transition-registry/trb1-mutations.md` recording per arm: anchor count
before/after, differing sha256, `go build` rc, **the failing test AND subtest name**, the inverse
`-skip` arm rc, and — for the four `MUT-ALLOW-*` arm (b)s — the words *"inverse arm unsatisfiable by
construction: co-detected by TestSketchRows (landed), not weakened"*.

#### T7a — **TR.B1 zero-tolerance activation** (AC5 only)

Edit AC5 in the design doc: delete `test "$count" -eq 0 ||` and require `test "$count" -eq 2 && …`.
**Leave AC6 and AC7 tolerant arms intact** — deleting them here would red TR.B1's own merge.
Machine check, with a known-positive control in the same call:

```bash
export PATH=/opt/homebrew/bin:$PATH
awk '/^5\. \*\*AC5/,/^6\. \*\*AC6/' design_docs/planned/w-transition-registry.md \
  | grep -c 'test "\$count" -eq 0'      # must be 0
awk '/^6\. \*\*AC6/,/^8\. \*\*AC8/' design_docs/planned/w-transition-registry.md \
  | grep -c 'test "\$count" -eq 0'      # control: must still be 2
```

Record `MUT-DELETE-TR-B-TEST` RED for AC5: delete one of the two required tests, re-run AC5, require
rc≠0 **and record the observed count (1, not merely "rc=1")**; restore from the `cp` backup, re-run green.

### `TR.B2` — descriptor-bound confinement and the two-session fixture (closes AC6, AC7)

Touches **only** `host/transitionreg` (production) plus the doc's AC6/AC7 rows. Depends on TR.B1.
Adds nothing to `host/broker`, so AC5 stays exactly 2 and TR.B1's activation cannot be regressed.

#### T4 — `Bind`, the manifest guard, and proposal agreement  (~200 impl / ~380 test LOC) → closes AC6

New file `host/transitionreg/bind.go`:

```go
// Binder is satisfied by *broker.Session. Declaring it HERE is what keeps the
// identifier broker.Session out of this package's production code (§0(i)).
type Binder interface { Bind(broker.Manifest) (*broker.BoundInvoker, error) }

// CapabilitySource is satisfied by *broker.Session; a counting fake satisfies it
// in tests, which is what gives MUT-CAPS-REREAD an observable (T5).
type CapabilitySource interface { CapabilitySnapshot(now int64) broker.CapabilitySnapshot }

type Proposal struct {
    TransitionFn, Interpreter hashref.HashRef
    SemanticsEpoch            int64
    RequiredCaps              EffectRequirement
    ExpectedEffects           []EffectRequirement
}

type Bound struct { /* descriptor + *broker.BoundInvoker */ }

func Bind(snap Snapshot, id string, caps broker.CapabilitySnapshot, target Binder) (*Bound, error)
func (b *Bound) Check(p Proposal) error
func (b *Bound) Request(ctx context.Context, req broker.EffectRequest, payload []byte) ([]byte, hashref.HashRef, error)
func (b *Bound) Descriptor() Descriptor   // deep copy
```

`Bind`'s ordered refusals, each its own named error with a **measured** message:

1. `snap.Lookup(id)` misses → `*TransitionAbsentError` — **`MUT-BIND-MISSING`**;
2. `broker.Allows(caps, requirementOf(d.Access))` not allowed → `*AccessDeniedError{Label: dec.Label}`,
   carrying **the broker's own label verbatim** — no relabelling, no collapsing four labels into one
   (§4.3 J7);
3. `target.Bind(broker.Manifest{…})` error → wrapped and returned, never swallowed.

`Check(p)` compares the proposal against the bound descriptor: `TransitionFn` (`MUT-PROPOSAL-FN`),
`RequiredCaps` (`MUT-PROPOSAL-CAPS`), `ExpectedEffects` as an **ordered, exact** set
(`MUT-PROPOSAL-EFFECTS`), plus `Interpreter` and `SemanticsEpoch` (§4.3 J9/J10 — Decision 7 pins all
three execution selectors, and a proposal that agrees on the source but not the interpreter is exactly
the case Decision 7 exists to forbid).

Tests, each refusal pinning **its own measured message** (the TR.A1 lesson: a message-agnostic
assertion stays green when the named guard is neutered, because a second refuser stands behind it):

- `TestGuardedSessionRefusesUndeclaredEffect` — handler counter **0**, epoch unchanged, error is
  `*broker.UndeclaredEffectError` with the exact message, one subtest each for undeclared *name*,
  undeclared *scope*, undeclared *cost*;
- `TestGuardedSessionStillRequiresBrokerGrant` — declared effect, no live grant → landed
  `*broker.DenialError` with the exact label; handler counter **0**; a second subtest where the grant is
  live and the same request succeeds, proving the refusal was the grant and not the guard;
- `TestProposalDescriptorAgreementRefusals` — table-driven, one named subtest per branch above plus an
  `all_fields_agree` positive control.

**Exit gate**: AC6 → `count=3`, all PASS. AC5 still `count=2`. `AC-INVOKE3` → 3. `AC-VET` rc=0. Hold set green.

#### T5 — the two-session request fixture  (~120 impl / ~300 test LOC) → closes AC7

```go
// Request is one consumer request: exactly one registry head and exactly one
// capability reading, both captured at construction and never re-read.
type Request struct { Registry Snapshot; Caps broker.CapabilitySnapshot }

func NewRequest(ctx context.Context, r Reader, caps CapabilitySource, now int64) (Request, error)
func (q Request) Allowed() []Descriptor   // snapshot byte order preserved; deep copies
```

- `NewRequest` calls `r.ReadSnapshot(ctx)` once and `caps.CapabilitySnapshot(now)` **once**. A reader
  error propagates (§4.3 J11) — it never degrades to an empty `Request`.
- `Allowed()` filters by `broker.Allows(q.Caps, requirementOf(d.Access)).Allowed` **without** re-reading
  either source. It is a pure function of the captured `Request`.

Tests:

- `TestTwoSessionExactOrderedSets` — one registry snapshot, two sessions with different grants; each
  `Allowed()` is the **exact** ordered set, asserted element-by-element in both directions (present and
  absent), plus a descriptor allowed to neither. **`MUT-SESSION-UNION`.**
- `TestNextReadObservesNewHeadWithoutRestart` — real store; publish revision 2 through the same
  `*StoreReader`; a **new** `NewRequest` returns `Head`=rev2 and the new entry, with no reader
  reconstruction. **`MUT-STARTUP-CACHE`.**
- `TestSingleRequestKeepsCapturedEpochs` — a counting `CapabilitySource` fake asserts exactly **1**
  call per `NewRequest` and **0** further calls across repeated `Allowed()`; between construction and
  the assertion the test both debits the session (epoch moves) and publishes a new head (registry moves),
  and `q.Caps.Epoch` / `q.Registry.Head` are unchanged. **`MUT-CAPS-REREAD`.** The counter lives in the
  fake, not in the code under test (rule 3i).

Reuse the landed `fakeObjectStore` (`transitionreg_test.go:25`), `openTransitionStore` (`:203`),
`seedRevision` (`:213`) and `validDescriptor` (`:21`) helpers rather than writing new ones.

**Exit gate**: AC7 → `count=3`, all PASS. AC6 still 3, AC5 still 2. `AC-VET` rc=0. Hold set green.

#### T6 — mutation sweep for TR.B2  (~17 arms, no production LOC)

As T3m, into `trb2-mutations.md`.

#### T7b — **TR.B2 ZERO-TOLERANCE ACTIVATION** (the merge criterion) → the milestone is not done without it

A count gate that still accepts 0 is satisfied by deleting the tests.

**T7b.a** — Edit AC6 and AC7: delete `test "$count" -eq 0 ||`, require `-eq 3` and run the tests.

**T7b.b** — Machine check that no tolerant arm survives in the AC5–AC7 range, with a known-positive
control in the same call:

```bash
export PATH=/opt/homebrew/bin:$PATH
awk '/^5\. \*\*AC5/,/^8\. \*\*AC8/' design_docs/planned/w-transition-registry.md \
  | grep -c 'test "\$count" -eq 0'    # must be 0
awk '/^11\. \*\*AC11/,/^## Non-Vacuity/' design_docs/planned/w-transition-registry.md \
  | grep -c 'test "\$count" -eq 1'    # control: must still be 1 (TR.C's arm, untouched)
```

**T7b.c** — Record `MUT-DELETE-TR-B-TEST` RED for AC6 and AC7: delete one required test per AC, re-run
that AC, require rc≠0 **and record the observed count** (AC6 `count=2`, AC7 `count=2`), restore from the
`cp` backup, re-run green.

**T7b.d** — Correct the doc's `MUT-CAPS-ALIAS` row to the two measured arms (§2 contradiction 5) and
add TR.B's rule-3j branch list (§4.3) to the Non-Vacuity table, so the next milestone's auditor is not
handed the same blind list TR.A2's auditor was.

**T7b.e** — Final full gate: `./scripts/verify_ail.sh` (repaired form), `./scripts/verify_go.sh`, and
`go vet ./host/...`, all **outside** any sandbox, all rc=0.

---

## 4. Mutation discipline

### 4.1 Protocol — every arm, no exceptions

```
1. cp <file> /tmp/trb_backup/<file>.bak
2. record: anchor grep -c  AND  shasum -a 256 <file>        # BEFORE
3. apply the neuter  (if false && <cond>   or   _ = f(x)) — NEVER delete a block:
   a deleted block orphans an import and reds the BUILD, which is the colour you
   predicted for the wrong reason.
4. assert the mutation LANDED: anchor count changed AND sha256 DIFFERS from step 2.
   If either is unchanged, STOP — "did not red" and "never ran" are the same rc.
5. assert the mutant BUILDS: GOTOOLCHAIN=go1.25.6 go build ./...   rc MUST be 0.
   "The mutant does not compile" is a THIRD fact wearing the same exit code as a kill.
6. kill arm:  go test <pkg> -run '<the named test>' -count=1        -> expect rc≠0
   RECORD THE FAILING TEST **AND SUBTEST** NAME. rc≠0 whose only FAIL is a
   pre-existing flake is not a kill.
6a. BROKER PACKAGE ONLY: every arm carries
      -skip 'TestHandlerTimeoutKillsTheWholeProcessGroup$'
    and AILANG_BIN=/tmp/ailang-v0300/ailang. Measured 2 failures / 11 isolated runs
    at BASE (§0(ii)); without the skip an 18% coin-flip is indistinguishable from
    your kill and from your inverse arm.
7. inverse arm: go test <pkg> -skip '<your test>|TestHandlerTimeout…$' -count=1
                -> MUST be rc=0. This is what proves YOUR test is the killer.
   EXCEPTION, declared in advance: for a mutation that alters a SHARED mechanism
   with a landed co-detector, this arm is UNSATISFIABLE BY CONSTRUCTION. Record it
   with that phrase and the co-detector's name. DO NOT weaken the co-detector.
   Applies to: MUT-ALLOW-* arm (b) (co-detector TestSketchRows, §0(iv)),
   and any arm in the extracted decideOver / debitGrant shared sites.
8. restore: cp /tmp/trb_backup/<file>.bak <file>  -- NEVER `git checkout -- <file>`;
   these files are uncommitted by construction in a sprint worktree and git checkout
   DELETES the executor's work.
9. assert the restore: sha256 equals step 2 exactly.
10. after every arm: GOTOOLCHAIN=go1.25.6 go vet ./host/... rc=0.
```

### 4.2 Mutation assignment — the 16 doc-named TR.B mutations

| Mutation | Task | arms | **Exact required observable** | rule 3i |
|---|---|---:|---|---|
| `MUT-CAPS-ALIAS` | T3m | **2** | (a) `…/snapshot_is_isolated_from_later_debit` (b) `…/grants_accessor_returns_a_fresh_copy` — the doc names one branch; the mechanism has two seams (§2.5) | ✓ |
| `MUT-CAPS-STATIC-EPOCH` | T3m | 1 | `…/allowed_invoke_increments_epoch_exactly_once` FAILs; single `debitGrant` site | ✓ |
| `MUT-ALLOW-NAME` | T3m | **2** | (a) in `Allows` → `…/denied_effect_name`, inverse rc=0 (b) in `Decide` → new test **and** `TestSketchRows` FAIL, inverse **unsatisfiable** (§0(iv)) | ✓ (b) is the delegation proof |
| `MUT-ALLOW-SCOPE` | T3m | **2** | as above, `…/denied_scope` | ✓ |
| `MUT-ALLOW-EXPIRED` | T3m | **2** | as above, `…/denied_expired` | ✓ |
| `MUT-ALLOW-BUDGET` | T3m | **2** | as above, `…/denied_budget` | ✓ |
| `MUT-BIND-MISSING` | T6 | 1 | `TestGuardedSessionRefusesUndeclaredEffect/absent_transition` FAILs on the **measured message**, not merely on "an error" | ✓ |
| `MUT-EFFECT-UNDECLARED` | T6 | 1 | handler counter becomes 1 → `TestGuardedSessionRefusesUndeclaredEffect` FAILs | ✓ counter in the fake |
| `MUT-EFFECT-BYPASS-BROKER` | T6 | 1 | `Request` calls the handler directly → `TestGuardedSessionStillRequiresBrokerGrant` FAILs (no `*DenialError`, counter 1) | ✓ |
| `MUT-PROPOSAL-FN` | T6 | 1 | `TestProposalDescriptorAgreementRefusals/transition_fn_mismatch` | ✓ |
| `MUT-PROPOSAL-CAPS` | T6 | 1 | `…/required_caps_mismatch` | ✓ |
| `MUT-PROPOSAL-EFFECTS` | T6 | 1 | `…/expected_effects_mismatch` | ✓ |
| `MUT-SESSION-UNION` | T6 | 1 | `TestTwoSessionExactOrderedSets` FAILs — the assertion is element-wise in **both** directions, so a union is caught by the *absent* half | ✓ |
| `MUT-STARTUP-CACHE` | T6 | 1 | `TestNextReadObservesNewHeadWithoutRestart` FAILs | ✓ |
| `MUT-CAPS-REREAD` | T6 | 1 | counting fake's `CapabilitySnapshot` call count ≠ 1 → `TestSingleRequestKeepsCapturedEpochs` FAILs | ✓ counter in the fake |
| `MUT-DELETE-TR-B-TEST` | T7a/T7b.c | **3** | one arm per activated AC; record the *count* observed, not just rc | ✓ |

**Doc-named subtotal: 16 mutations / 23 arms.**

### 4.3 Rule 3j — refusal branches enumerated from TR.B's OWN DIFF, not from the doc's decision list

**This is the section TR.A2 got wrong and the reason it shipped three uncovered branches.** The doc's
rule-3j sentence (line 621) is scoped to Decisions 1/2/4 — i.e. to TR.A — and *by construction* cannot
contain a branch TR.B writes. The enumeration below is anchored to what TR.B's diff adds.

**Cut instrument** (calibrated in §0(vi); `grep`, never `rg`):

```bash
export PATH=/opt/homebrew/bin:$PATH
git diff --unified=0 -- host/broker host/transitionreg \
  | grep -cE '^\+.*return .*(fmt\.Errorf|errors\.New)\(|^\+.*return .*&[A-Z][A-Za-z]*\{'
```

Known-positive control, same call: the identical cut over TR.A's landed diff
(`git diff --unified=0 b0f323a..HEAD -- host/transitionreg host/store`) → **96**. Negative control over
`design_docs/` → **0**. The `%w`-only cut returns **30** on the same TR.A diff and would have missed
**64 of 94** — do not use it. **The grep is a floor, not the enumeration.** The enumeration is by
reading the diff: every `return` on an error path, including bare `return err`, `return false`, and
early `continue`/`break` guards, gets an arm.

| # | Branch TR.B's diff adds | Mutation | Observable |
|---|---|---|---|
| J1 | `Session.Bind` rejects a negative-cost requirement | `MUT-MANIFEST-NEG-COST-OK` | `TestBindRefusesMalformedManifest/negative_cost` |
| J2 | `Session.Bind` rejects a duplicated declared requirement | `MUT-MANIFEST-DUP-OK` | `…/duplicate_declared` |
| J3 | `Request` membership is on the triple, not the effect **name** | `MUT-DECLARED-NAME-ONLY` | `TestGuardedSessionRefusesUndeclaredEffect/undeclared_scope` |
| J4 | …and not on name+scope with a free **cost** | `MUT-DECLARED-COST-ANY` | `…/undeclared_cost` |
| J5 | the epoch increments on the **replay** debit path too | `MUT-EPOCH-LIVE-ONLY` | `TestCapabilitySnapshotEpochAndIsolation/replay_debit_increments_epoch` |
| J6 | `Allows` uses `snap.Now`, not `0` and not a fresh clock | `MUT-ALLOW-NOW-ZERO` | `TestAllowsUsesDecideAllFourDenials/uses_snapshot_now_not_a_wall_clock` |
| J7 | `Bind` propagates the broker's **exact** denial label, not a collapsed one | `MUT-BIND-COLLAPSE-LABEL` | `TestGuardedSessionStillRequiresBrokerGrant` label subtests (4) |
| J8 | `Bind` refuses a zero `Snapshot` (no head captured) | `MUT-BIND-EMPTY-SNAPSHOT-OK` | `TestGuardedSessionRefusesUndeclaredEffect/zero_snapshot` |
| J9 | `Check` compares `Interpreter` (Decision 7's second selector) | `MUT-PROPOSAL-INTERP` | `TestProposalDescriptorAgreementRefusals/interpreter_mismatch` |
| J10 | `Check` compares `SemanticsEpoch` (third selector) | `MUT-PROPOSAL-EPOCH` | `…/semantics_epoch_mismatch` |
| J11 | `NewRequest` propagates a reader error rather than degrading to empty | `MUT-REQUEST-SWALLOW` | `TestSingleRequestKeepsCapturedEpochs/injected_read_error` |
| J12 | `Allowed()` returns deep copies (schema and declared-effect slices) | `MUT-REQUEST-ALIAS` | `TestTwoSessionExactOrderedSets/returned_descriptors_are_copies` |
| J13 | `Allowed()` preserves the snapshot's bytewise ID order | `MUT-REQUEST-REORDER` | `TestTwoSessionExactOrderedSets/order_is_the_snapshot_order` |

**Rule-3j subtotal: 13 mutations / 13 arms. Grand total 29 mutations / 36 arms.**

Cost: the scoped `./host/transitionreg` package is 0.31 s and a filtered `./host/broker -run …` is
~0.3 s on this rig, so 36 arms scripted through §4.1 is roughly 1–1.5 hours of executor wall time,
split T3m/T6. **This is the honest cost of rule 3j and it is not optional: if the controller wants to
cut it, the cut must be a recorded decision, not a silent omission.**

### 4.4 Two rows flagged UNCERTAIN

1. **`MUT-EFFECT-BYPASS-BROKER`.** The doc's neuter is "invoke handler directly". If the executor
   implements `BoundInvoker.Request` as *"check membership, then `s.invoke`"*, the only way to bypass
   the broker is to reach into `s.registry` from `Request`, which is a rewrite rather than a neuter and
   risks a non-compiling mutant (§4.1 step 5). The compiling form is
   `handler, _ := b.s.registry[req.Effect]; return handler.Execute(ctx, req, payload)` inside `Request`,
   with `if false &&` on the guard. The executor must report which form was used and that step 5 passed.
2. **`MUT-SESSION-UNION` vs a set assertion that only checks membership.** If
   `TestTwoSessionExactOrderedSets` asserts only *"every returned descriptor is allowed"*, a union
   passes. The assertion must be **exact and two-sided**: the full expected slice compared element-wise
   (`reflect.DeepEqual` on IDs), plus an explicit assertion that each session's set does **not** contain
   the other's exclusive descriptor. The executor must report which assertion fired.

---

## 5. Acceptance commands, as the executor must run them

All baselined this session on the pristine tree at `66a1d63`; the base result is part of each criterion.
**No command contains `rg`** (P17). Every command starts `export PATH=/opt/homebrew/bin:$PATH` and
carries `GOTOOLCHAIN=go1.25.6`.

| AC | base (measured at `66a1d63`) | TR.B delivered |
|---|---|---|
| **AC5** | `count=0`, rc=0 | **`count=2`**, both tests PASS (TR.B1) |
| **AC6** | `count=0`, rc=0 | **`count=3`**, three tests PASS (TR.B2) |
| **AC7** | `count=0`, rc=0 | **`count=3`**, three tests PASS (TR.B2) |
| AC1 / AC2 / AC3 / AC4 | `3 / 3 / 4 / 2`, all PASS | **unchanged** |
| AC8 | `count=2`, `ok host/replay 2.054s` | unchanged; no `broker`/`transitionreg` import in production replay |
| AC9 (repaired) | rc=0, `modules11=1 tests14=1 steps9=1` | identical |
| AC10 | build PASS, `count=1`, whole `./host/transitionreg` PASS | unchanged |
| AC11 | `count=1` | **still 1** — moving it means TR.C was built |
| `verify_go.sh` | **rc=0**; plain `host/broker 48.8s`, `host/transitionreg 3.5s`; race `92.3s` / `5.0s` | rc=0; new broker+transitionreg race time must stay **< 120 s** combined |

### The three new criteria (§2 delta)

```bash
# AC-INVOKE3 — protects TR.C's exact three-site carve-out. Base: 3, all publish_op.go.
export PATH=/opt/homebrew/bin:$PATH
n=$(grep -rn '\.Invoke(' --include='*.go' host/ cmd/ | grep -v _test.go | wc -l | tr -d ' ')
p=$(grep -rn '\.Invoke(' --include='*.go' host/ cmd/ | grep -v _test.go | grep -c 'host/broker/publish_op.go')
t=$(grep -rn '\.Invoke(' --include='*.go' host/ cmd/ | grep -c _test.go)   # known-positive control
[ "$n" -eq 3 ] && [ "$p" -eq 3 ] && [ "$t" -gt 0 ]
```

```bash
# AC-VET — copylocks and friends live OUTSIDE go test's default vet subset (§0(iii)). Base: rc=0.
export PATH=/opt/homebrew/bin:$PATH
GOTOOLCHAIN=go1.25.6 go vet ./host/...
```

```bash
# AC-NOFLAKE — the whole broker package, serial, with the known base flake skipped and
# AILANG_BIN set. A bare `go test ./host/broker` is red at base ~18% of the time (§0(ii))
# and is red 100% of the time without AILANG_BIN (TestEpisodeLiveReplayThreeArmsAndEvidence).
export PATH=/opt/homebrew/bin:$PATH
GOTOOLCHAIN=go1.25.6 AILANG_BIN=/tmp/ailang-v0300/ailang \
  go test ./host/broker -skip 'TestHandlerTimeoutKillsTheWholeProcessGroup$' -count=1
```

---

## 6. Estimates and the split verdict

| Task | milestone | impl LOC | test LOC | notes |
|---|---|---:|---:|---|
| T1 capability snapshot + epoch | TR.B1 | 90 | 150 | 8 named subtests; single `debitGrant` extraction |
| T2 `Requirement` / `Allows` / `decideOver` | TR.B1 | 50 | 180 | 10 named subtests; delegation, not re-tabling |
| T3 confined seam + `invoke` extraction | TR.B1 | 110 | 160 | the TR.C-compatibility freeze |
| T3m mutations (TR.B1) | TR.B1 | 0 | 0 | 16 arms + transcript |
| T7a activation (AC5) | TR.B1 | 0 | 0 | doc edit + 1 delete arm |
| **TR.B1 subtotal** | | **250** | **490** | **≈740 LOC** |
| T4 `Bind` / guard / proposal agreement | TR.B2 | 200 | 380 | 3 AC tests, ~14 subtests |
| T5 two-session request fixture | TR.B2 | 120 | 300 | 3 AC tests, counting fakes |
| T6 mutations (TR.B2) | TR.B2 | 0 | 0 | 20 arms + transcript |
| T7b activation (AC6, AC7) | TR.B2 | 0 | 0 | doc edit + 2 delete arms + full gates |
| **TR.B2 subtotal** | | **320** | **680** | **≈1000 LOC** |
| **TR.B total** | | **570** | **1170** | **≈1740 LOC** |

### Verdict: **the doc's "1 day" is wrong by roughly a factor of two. TR.B needs a split.**

- Reference velocity: `VL.B` priced 515 LOC at 0.5 day → **~1000 LOC/day**.
- TR.A was priced at 2630 LOC over 2 days (**~1315/day**), its own planner called that above velocity,
  recommended a split, and **the split held** across two landed milestones.
- TR.B at 1740 LOC in **1 day** is **~1740/day** — higher than the number TR.A was split for, on a
  milestone the doc prices at *half* TR.A's duration. Add 36 mutation arms and a design collision with
  TR.C that has to be resolved inside the same day.

**Recommended split: `TR.B1` (T1–T3, T3m, T7a) ≈ 740 LOC ≈ 0.75 day; `TR.B2` (T4–T6, T7b) ≈ 1000 LOC
≈ 1 day.** Total ≈ 1.75–2 days.

**Why this boundary is safe — the same per-AC, directory-independent reasoning that made TR.A's split safe:**

| property | check |
|---|---|
| **No earlier milestone can red a later one's criteria** | TR.B1 touches only `host/broker`. AC6/AC7 grep `./host/transitionreg` and stay at `count=0`, where their **base-tolerant arm keeps them green** — so TR.B1 merges without them. TR.B2 touches only `host/transitionreg` production, so AC5's `count=2` cannot move after TR.B1's activation. |
| **Each half is independently mergeable** | TR.B1 delivers a working capability snapshot + delegating predicate + confined seam; nothing in it is dead code for TR.C, which needs exactly the `invoke` extraction TR.B1 performs. TR.B2 is a pure consumer of TR.B1's exports. |
| **The dependency is one-directional and total** | TR.B2 cannot compile without `broker.Allows`, `broker.CapabilitySnapshot`, `broker.Manifest`, `broker.BoundInvoker`. TR.B1 references nothing in `host/transitionreg`. |
| **The activation split is clean** | T7a deletes AC5's tolerant arm only; T7b deletes AC6's and AC7's. The doc's merge criterion is per-AC, so neither half leaves a half-activated criterion. |
| **The doc's own constraint is preserved** | "All three milestones are independently mergeable and at most 2 days" — both halves are well under 2 days, and TR.C is unaffected by the split (it needs TR.B1 only). |

**If the controller declines the split**, the honest fallback is: run TR.B whole and accept a likely
2-day slip, with T3 (the TR.C-compatibility freeze) executed **first** so the collision is resolved
before any test is written. Do **not** cut the mutation sweep to make one day.

---

## 7. Execution protocol

- **Worktree**: a sibling of the repo, e.g. `/Users/voightkampff/dev/sunholo-data/.wt-iter73`.
  **Never under `/tmp`** — cwd-relative path tests then fail for the location rather than the code, and
  a red CI can never reproduce it. (`git worktree list` at base shows only the main checkout.)
- **Every** Bash call starts `export PATH=/opt/homebrew/bin:$PATH`. Without it `go`, `gh` and `node` are
  rc=127 and look like a broken toolchain.
- **Every** `go` invocation carries `GOTOOLCHAIN=go1.25.6`.
- **Every whole-package `./host/broker` run carries `AILANG_BIN=/tmp/ailang-v0300/ailang`.** Measured at
  base: without it, `TestEpisodeLiveReplayThreeArmsAndEvidence` fails with *"AILANG_BIN must name the
  pinned released interpreter"* — a red that measures the environment, not the change. `AILANG_BIN`
  points **outside** the repo at v0.30.0 (`ailang --version` → `AILANG v0.30.0`, commit `e37b370`).
  Never the PATH `ailang`: it is a `-dirty` dev build and CLAUDE.md forbids it.
- **The broker package has a measured 18 % base flake** (§0(ii)). Skip it, and read FAIL lines, never
  exit codes.
- **`go vet ./host/...` after every task and every mutation arm.** `go test` is blind to `copylocks`
  (§0(iii)) and `verify_go.sh` never calls `go vet`.
- **zsh, not bash**, in this harness: `${PIPESTATUS[0]}` is silently empty (zsh spells it
  `${pipestatus[1]}`, 1-indexed) — or capture to a file; an unquoted glob-shaped flag value
  (`--include=*.go`) aborts the command before it runs; zsh does **not** word-split an unquoted
  variable, so `cmd $FILES` passes ONE argument — use arrays and assert `${#arr[@]}`; brace a variable
  followed by a colon (`"${rev}:path"`, never `"$rev:path"`); `echo` interprets backslash escapes — use
  `printf '%s'` / `od -c`.
- **`rg` is not a binary.** Never in an acceptance command, a script, or a mutation observable. Use `grep`.
- **Restores are `cp` from `/tmp/trb_backup/`.** `git checkout -- <file>` deletes uncommitted work.
- **Never touch** `~/.ailang/state/mission-v1*` or the V1 checkout.
- **SANDBOX CAVEAT — read this before reporting any gate result.** A gate verdict obtained inside a
  `workspace-write` sandbox is **UNINFORMATIVE — neither a pass nor a fail.** Loopback socket binds are
  denied there, which both *invents* failures (`host/daemon`, `host/broker` red for
  `bind: operation not permitted`) and *hides* real ones. That is what produced the design doc's V13 row
  claiming `verify_go.sh` rc=1, which is **rc=0** outside the sandbox (P15). Report sandbox results as
  "sandbox, uninformative"; **the controller re-runs every gate outside the sandbox** and that run is the
  verdict.

---

## 8. Risks — R1 and R2 want a controller decision BEFORE execution

| # | Risk | Assessment |
|---|---|---|
| **R1** | **TR.B as designed reds TR.C's AC11.** §0(i): the descriptor-bound invoker's call into the broker pipeline is a fourth production `Invoke` selector call under either placement, and TR.C pins "exactly 3, outside-broker zero". | **DECISION WANTED.** Options: **(a) freeze the unexported `invoke` extraction + `BoundInvoker.Request` seam + the new `AC-INVOKE3` hold criterion** — my recommendation, it costs ~15 LOC and leaves TR.C's design untouched; (b) let TR.B add a fourth in-broker site and amend TR.C's pinned exemption to 4 with named identities — cheaper to write, but it makes an earlier milestone move a later one's criterion and blurs `MUT-BINDING-FOURTH-INVOKE`; (c) defer the bound invoker to TR.C — but then AC6 has no mechanism and TR.B closes nothing. |
| **R2** | **The doc prices TR.B at 1 day; I price it at 1740 LOC ≈ 2 days.** §6. | **DECISION WANTED.** Options: **(a) split at the T3/T4 boundary into TR.B1 (AC5) and TR.B2 (AC6+AC7)** — my recommendation, justified per-AC and directory-independently in §6; (b) run whole and accept the slip, T3 first. Not an option: cutting the mutation sweep. |
| **R3** | The `sync.Mutex` on `Session` is not reentrant; `CapabilitySnapshot` called from inside `invoke`'s locked region deadlocks and presents as a 10-minute `go test` panic, not an error. This is TR.B's analogue of TR.A's tx deadlock. | Mitigated by the T1 constraint. If any `go test ./host/broker` exceeds 120 s, suspect a lock-inside-lock **first**, before suspecting the flake. |
| **R4** | The `Invoke` → `invoke` extraction touches the single most load-bearing method in the repo (the decide/debit/dispatch/record pipeline, including the `EffectRegistryPublish` attended-approval special case at `broker.go:179-201`). | Mitigated: it is a pure rename plus a one-line wrapper, with no reordering and no lock change. The regression net is the whole landed broker suite (25 top-level tests, `ok … 35.4 s` serial) plus `verify_go.sh`'s race leg, both run at T3's exit. A descriptor may legitimately declare `Registry.Publish`; the landed `validatePublishApproval` still runs unchanged inside `invoke`, so TR.B adds no new branch there and needs no new test for it. |
| **R5** | **`TestHandlerTimeoutKillsTheWholeProcessGroup` is flaky at base — 2 failures in 11 isolated runs (§0(ii)).** It is a pre-existing `host/broker` defect, not TR.B's. | **Reported, not repaired.** Out of scope; TR.B must not silence it. **Recommend the controller file it as a queue item**: a test that reds ~18 % of the time in the package every future broker milestone will mutate is a standing tax on every mutation sweep, and it is exactly the "load confound" class the `w-bench-load-confound` sprint already studied. |
| **R6** | `Allows` could quietly become the second policy engine Decision 5 forbids, by restating the four comparisons instead of delegating. A behaviourally-correct restatement passes every ordinary test. | Mitigated by two independent instruments: the mechanical `awk`+`grep` check in T2 (zero comparison operators inside `Allows`, with a firing control over `Decide`), and `MUT-ALLOW-*` **arm (b)**, whose co-detection is the only positive evidence of delegation (§0(iv)). |
| **R7** | AC1–AC4, AC10 and AC11 could drift if the executor "helpfully" touches the codec or starts TR.C. | Their values (3/3/4/2, `count=1`, `count=1`) are in the §5 table as **hold** criteria and are re-measured at every task exit and at T7b.e. TR.B changes no wire format and no golden byte (§1). |
| **R8** | The race leg has a hard 600 s kill and `host/broker` alone is already **92.3 s** of it (P15). | Budget the new broker and transitionreg tests at **< 30 s** combined on the race leg; no sleeps, no wall-clock waits. The executor must report `ok host/broker <time>` and `ok host/transitionreg <time>` from the race leg at T7b.e. |

---

## 9. Handoff

- **Sprint plan**: `design_docs/planned/w-transition-registry-trb-sprint-plan.md` (this file)
- **Sprint JSON**: `.ailang/state/sprints/w-transition-registry-trb.plan.json`
- **Design doc**: `design_docs/planned/w-transition-registry.md`
- **Base**: `dev` @ `66a1d63`
- Neither artifact is committed by the planner. **The controller commits.**

SPRINT_PLAN_PATH: design_docs/planned/w-transition-registry-trb-sprint-plan.md
SPRINT_JSON_PATH: .ailang/state/sprints/w-transition-registry-trb.plan.json
