# TR.B1 mutation transcript

All commands ran inside the workspace-write sandbox. Loopback denial makes whole-package inverse
arms uninformative; the controller must rerun them outside the sandbox. Every executed Go command
used `GOTOOLCHAIN=go1.25.6`; broker tests used the pinned `AILANG_BIN` and skipped the named base
flake.

## MUT-CAPS-ALIAS (session to snapshot)

- Anchor `grants: append`: before 2, after 1.
- SHA-256: before `8044a6503cb87b8683ae49fa9fb21806f7601bb81fd33f0d34ed13d4e31ad73d`;
  mutant `c2ce9afbb22d357c706608640c473e3d97152a9e556cc54f955f9d20fbb0a66c`.
- `go build ./...`: rc=0.
- KILLED: `TestCapabilitySnapshotEpochAndIsolation/snapshot_is_isolated_from_later_debit`
  failed with old snapshot budget 3, want 5; kill rc=1.
- Inverse `-skip 'TestCapabilitySnapshotEpochAndIsolation|TestHandlerTimeoutKillsTheWholeProcessGroup$'`:
  rc=1 solely because `TestAttendedPublishMintsThroughTheLandedTraversalAndSpendsExactlyOnce`
  panicked at `httptest.NewServer` with `bind: operation not permitted`. Sandbox result is
  uninformative, not a mutation survival or regression.
- Restore SHA-256 equaled backup SHA-256
  `8044a6503cb87b8683ae49fa9fb21806f7601bb81fd33f0d34ed13d4e31ad73d`.
- Post-restore `go vet ./host/...`: rc=0.

## Deferred arms

The following arms were not run in this executor turn and are explicitly deferred to the controller:

- `MUT-CAPS-ALIAS` (snapshot accessor seam)
- `MUT-CAPS-STATIC-EPOCH`
- `MUT-ALLOW-NAME` arms (a) and (b)
- `MUT-ALLOW-SCOPE` arms (a) and (b)
- `MUT-ALLOW-EXPIRED` arms (a) and (b)
- `MUT-ALLOW-BUDGET` arms (a) and (b)
- `MUT-MANIFEST-NEG-COST-OK`
- `MUT-MANIFEST-DUP-OK`
- `MUT-DECLARED-NAME-ONLY`
- `MUT-DECLARED-COST-ANY`
- `MUT-EPOCH-LIVE-ONLY`
- `MUT-ALLOW-NOW-ZERO`

Reason: only the first arm was completed under the full required protocol before the execution
window; no unrun arm is represented as killed. The plan's “~16 arms” is arithmetically inconsistent:
T3m has 11 doc-named arms (2 alias + 1 epoch + 8 allow) and six applicable rule-3j arms (J1-J6),
which is 17 before T7a and 18 including the delete-test arm. In addition, the implementation validates
`Manifest.Access.Cost` independently from declared costs, so its distinct refusal branch requires an
additional executor-enumerated mutation beyond J1 if retained.

## MUT-DELETE-TR-B-TEST (T7a)

- Neuter: renamed `TestAllowsUsesDecideAllFourDenials` so it was absent from `go test -list`.
- Anchor count before 1, after 0.
- SHA-256: before `cfbcb415d9e7483aac035875ad5045b31dfd85c6bb1f5d8a035959b830a47530`;
  mutant `59c589e7805345b04685aec2eb7f7f6f4cd8be4ffee35248ac5357cd542b6c47`.
- `go build ./...`: rc=0.
- KILLED: activated AC5 observed count **1**, required 2; kill rc=1.
- Inverse `-skip 'TestAllowsUsesDecideAllFourDenials|TestHandlerTimeoutKillsTheWholeProcessGroup$'`:
  rc=1 solely because the sandbox denied `httptest.NewServer` in
  `TestAttendedPublishMintsThroughTheLandedTraversalAndSpendsExactlyOnce`; uninformative.
- Restore and backup SHA-256 both
  `cfbcb415d9e7483aac035875ad5045b31dfd85c6bb1f5d8a035959b830a47530`.
- Post-restore `go vet ./host/...`: rc=0.

---

# Controller sweep — the 17 deferred arms, run outside the sandbox

The executor deferred 17 of its 19 required arms and enumerated them rather than
claiming them. The controller ran **18 arms** (the 17 deferred plus a re-run of
`MUT-CAPS-ALIAS` at the accessor seam), outside the `workspace-write` sandbox, so
every inverse arm is informative here where the executor's were not.

Protocol per arm, no exceptions: `cp` backup → mutate → assert **LANDED** (sha256
differs) → assert **BUILDS** (`go build ./...` rc=0) → scoped kill arm, recording
**which test and subtest** failed → inverse arm (`-skip` the killer plus the known
base flake, whole package) → `cp` restore → assert byte-identity.

**Result: 18 arms, 17 KILLED, 1 SURVIVED.**

| arm | file | verdict | killing subtest / note |
|---|---|---|---|
| `MUT-CAPS-STATIC-EPOCH` | broker.go | KILLED, isolated | `allowed_invoke_increments_epoch_exactly_once` + `replay_debit_increments_epoch` |
| `MUT-CAPS-ALIAS-ACCESSOR` | broker.go | KILLED, isolated | `grants_accessor_returns_a_fresh_copy` |
| `MUT-CAPS-SNAPSHOT-ALIAS` | broker.go | KILLED, isolated | `snapshot_is_isolated_from_later_debit` |
| `MUT-ALLOW-NOW-ZERO` | decide.go | KILLED, isolated | `uses_snapshot_now_not_a_wall_clock` (+ 2 more) |
| `MUT-ALLOW-EMPTY-LEDGER` | decide.go | KILLED, co-detected | `no_grants_is_denied_effect_name` |
| `MUT-ALLOW-NAME` | decide.go | KILLED, inverse unsatisfiable | `denied_effect_name`; co-detected by landed `decide_test.go`, NOT weakened |
| `MUT-ALLOW-SCOPE` | decide.go | KILLED, inverse unsatisfiable | `denied_scope`; co-detected by a landed test, NOT weakened |
| `MUT-ALLOW-EXPIRED` | decide.go | KILLED, inverse unsatisfiable | `denied_expired`; co-detected by a landed test, NOT weakened |
| `MUT-ALLOW-BUDGET` | decide.go | KILLED, inverse unsatisfiable | `denied_budget`; co-detected by a landed test, NOT weakened |
| `MUT-ALLOW-RANK` | decide.go | KILLED, **isolated** | `ranked_best_denial_across_three_grants` — see note below |
| `J1-MUT-MANIFEST-ACCESS-NEG-COST-OK` | confined.go | KILLED, isolated | `negative_access_cost` |
| `J1b-MUT-MANIFEST-DECLARED-NEG-COST-OK` | confined.go | KILLED, isolated | `negative_cost` |
| `J2-MUT-MANIFEST-DUP-OK` | confined.go | KILLED, isolated | `duplicate_declared` |
| `J3-MUT-DECLARED-NAME-ONLY` | confined.go | KILLED, isolated | `undeclared_scope` + `undeclared_cost` |
| `J4-MUT-DECLARED-COST-ANY` | confined.go | KILLED, isolated | `undeclared_cost` |
| `J5-MUT-DECLARED-SCOPE-ANY` | confined.go | KILLED, isolated | `undeclared_scope` |
| **`J6-MUT-BIND-DECLARED-ALIAS`** | confined.go | **SURVIVED** | see below |
| `J7-MUT-DECLARED-ACCESSOR-ALIAS` | confined.go | KILLED, isolated | `declared_accessor_is_a_fresh_copy` |

Every arm restored byte-identical from its `cp` backup.

## The survival: `J6-MUT-BIND-DECLARED-ALIAS`

Mutation: in `Bind`, `declared := append([]Requirement(nil), m.Declared...)`
becomes `declared := m.Declared`.

- **LANDED** — sha256 `c6fbdfc71786` → `62b0b1253242`.
- **BUILDS** — `go build ./...` rc=0.
- **SURVIVED** — the *entire* `host/broker` package is rc=0 with the defect present.

So this is a genuine survival, not a non-compiling mutant and not an instrument
failure.

**What it means.** `Bind`'s own doc comment says it "validates and **copies** a
descriptor authority envelope". Aliasing means the envelope is *not frozen at bind
time*: a caller that mutates its `Manifest.Declared` slice after `Bind` retroactively
changes what an already-bound invoker will accept. In a clause-3 substrate whose
whole claim is that authority is explicit and declared up front, a declaration set
that can be widened after the fact by an unrelated write is precisely the failure the
milestone exists to prevent.

**Why it was invisible.** The asymmetry is the finding. `J7` — the *output* side, the
`Declared()` accessor — was pinned by `declared_accessor_is_a_fresh_copy`, and it
kills cleanly. Nobody wrote the mirror assertion for the *input* side. That is this
repo's named recurring shape, **guard the helper, miss the call site**, arriving one
more time: the copy that is tested is the one going out, and the copy that is not
tested is the one coming in.

**Fixed in-PR** by the subtest `bind_copies_the_caller_declaration_slice`, which pins
both halves — that `Declared()` still shows the envelope frozen at `Bind` time, and,
the authority-bearing half, that `Request` still REFUSES the triple the caller tried
to inject afterwards, with its own measured message. Re-run against the identical
mutant (same sha256 `62b0b1253242`, `go build` rc=0): **KILLS**, failing exactly
`TestBindRefusesMalformedManifest/bind_copies_the_caller_declaration_slice`, and the
inverse arm (`-skip` that test, whole package) is **rc=0 with 0 FAIL** — so the new
subtest is the killer, not a bystander. Restored byte-identical.

It is a SUBTEST, not a new top-level test, so AC5 stays at exactly the 2 this
milestone just activated.

## One planner prediction refuted

`MUT-ALLOW-RANK` was expected to be co-detected by a landed test (the four
`MUT-ALLOW-*` arms are, which is what makes their inverse arms unsatisfiable by
construction). It was **isolated** instead: the landed `decide_test.go` covers
`Decide` on a single capability and does not exercise the ranked selector across a
multi-grant ledger. So the delegation evidence for *ranking* rests solely on the new
`ranked_best_denial_across_three_grants`, not on co-detection with landed coverage.
Recorded rather than smoothed over — it is weaker evidence than the plan assumed,
though the pin itself is real and non-vacuous.

---

# Evaluator finding — a SECOND uncovered call site, in the mechanism the sweep audited

The independent evaluator (`sonnet`, generator≠judge — the executor ran on
`codex:gpt-5.6-sol`) was handed the `debitGrant` extraction as a named target and found
a blocking gap the controller's own 18-arm sweep missed. **Reproduced first-party before
acting on it**, per the rule that a judge's finding is a claim like any other:

- Mutation: at the **failed**-replay branch (`broker.go:390`) only, replace
  `s.debitGrant(grantIndex, decision.Remaining)` with a direct
  `s.grants[grantIndex].Budget = decision.Remaining` — i.e. keep the budget write, drop
  the epoch bump. The live site (`:245`) and the succeeded-replay site (`:403`) are left
  untouched (control: `grep -c 'debitGrant(grantIndex'` = **2** after the mutation).
- **LANDED** — sha256 `8044a6503cb8` → `dec185739972`.
- **BUILDS** — `go build ./...` rc=0.
- **SURVIVED** — whole `./host/broker` rc=0, **0 FAIL**, `ok 35.554s` with the defect
  present.

`debitGrant` has **three** call sites. `replay_debit_increments_epoch` covers the
succeeded-replay one; the live one is covered by
`allowed_invoke_increments_epoch_exactly_once`. The **failed**-replay one was covered by
nothing — `TestReplayOfFailedRecordReproducesTheFailure` exercises that branch but never
asserts `Epoch` or `Budget`. T1's exit-gate language claims the subtest pins that "the
replay debit **sites** are on the same mechanism", plural; only one of the two was
actually pinned.

**This is the same shape as `J6`, one layer down, and that is the durable finding.** The
production code is correct today — all three sites really do call the one mechanism — so
every direct observation agrees with the claim. What is missing is the assertion that
would notice a future edit breaking only the third site. A sweep that unifies three call
sites into one mechanism naturally tests *the mechanism*, and the thing it stops testing
is *the sites*. Two instances in one milestone (`J6`'s input-side copy, and this), plus
three in the previous one, is now this repo's most reliably recurring defect class:
**guard the helper, miss the call site.**

**Fixed in-PR** by the subtest `replay_debit_increments_epoch_on_failure`, which drives a
live invoke through a failing handler, replays the resulting failed record, and asserts
the epoch advances by exactly 1 across the failed-replay debit. Re-run against the
identical mutant (same sha256 `dec185739972`, `go build` rc=0): **KILLS**, failing exactly
`TestCapabilitySnapshotEpochAndIsolation/replay_debit_increments_epoch_on_failure`, and
the inverse arm (`-skip` that test, whole package) is **rc=0 with 0 FAIL** — the new
subtest is the killer, not a bystander. Restored byte-identical. A SUBTEST, so AC5 stays
at exactly 2.

Post-fix gates, all outside the sandbox: `verify_go.sh` **rc=0** (0 FAIL; broker 48.180s
plain / 90.886s race, transitionreg 2.860s / 2.999s), `go vet ./host/...` **rc=0**,
AC5 **count=2**.
