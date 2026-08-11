# Sprint plan — `TR.A` (queue item 11 `w-transition-registry`)

**Milestone**: `TR.A` — *immutable descriptor, canonical codec, eager snapshot, and CAS publication*
**Status**: PLANNED · one milestone · 2 days · **TR.A ONLY** (TR.B and TR.C are out of scope)
**Design doc**: [`design_docs/planned/w-transition-registry.md`](w-transition-registry.md) (landed iter-70, PR #58 → `11fb1fd`)
**Base**: `dev` @ `adfaa0b`, clean, CI green both jobs
**Planner**: mission-control iteration 71, opus, first-party measurement on this rig
**Executor**: sandboxed worktree, **NO git write permission**.
**THE CONTROLLER MAKES ALL COMMITS.** The executor never runs `git commit`, `git add`, `git push`,
`git checkout`, `git stash`, `git restore`, or `gh pr`. Restores are `cp` from a backup — see §7.

---

## 0. Planner's first-party verification of every premise

Every controller-supplied premise was re-run on this rig before it entered the plan. Negative
results are paired with a known-positive control **in the same call**, per the standing rule that an
empty result is a claim, not a fact.

| # | Premise | Verdict | Command → observed output |
|---|---|---|---|
| P1 | `host/transitionreg` is absent | **CONFIRMED** | `ls -d host/transitionreg` → `No such file or directory`, rc=1; control **in the same call** `ls -d host/registry` → `host/registry`, rc=0 |
| P2 | exactly 14 `host/` packages | **CONFIRMED** | `GOTOOLCHAIN=go1.25.6 go list ./host/... \| wc -l` → `14`: archive boundary broker canon capsule childenv daemon hashref pkgproj registry replay runbook store verifygate |
| P3 | store head + object primitives exist at the cited lines | **CONFIRMED** | `grep -n 'func (s \*Store) \(SetRegistryHead\|GetRegistryHead\|PutObject\|GetObject\|SelectedHead\|SelectHead\)' host/store/store.go` → `439 PutObject`, `463 GetObject`, `606 SetRegistryHead`, `624 GetRegistryHead`, `704 SelectedHead`, `731 SelectHead` |
| P4 | `schema.sql` declares exactly 8 tables; TR.A adds none | **CONFIRMED** | `grep -n 'CREATE TABLE' host/store/schema.sql` → `objects, worlds, log_entries, epoch_registry_heads, store_heads, verification_cache, journal, approval_claims` (8) |
| P5 | AC1 base | **CONFIRMED** | exact AC1 command → `count=0`, `rc=0` |
| P6 | AC2 base | **CONFIRMED** | exact AC2 command → `count=0`, `rc=0` |
| P7 | AC3 base | **CONFIRMED** | exact AC3 command → `count=0`, `rc=0` |
| P8 | AC4 base = 1, the existing control | **CONFIRMED** | exact AC4 command → `count=1`, `rc=0`; the listed name is `TestRegistryHeadRoundTrip` |
| P9 | AC10 base | **CONFIRMED** | exact AC10 command → build PASS, `ok host/store 0.192s`, `ok host/broker 0.234s`, `count=0`, `rc=0` |
| P10 | AC11 base = 1 (TR.C control, untouched by TR.A) | **CONFIRMED** | exact AC11 count → `1` |
| P11 | AC9 base 4/11/14 | **CONFIRMED** | `AILANG_BIN=/tmp/ailang-v0300/ailang ./scripts/verify_ail.sh` → rc=0, `✓ 4/4 required world/ identities verified across 11 module(s)`, `✓ all 14 required named tests pass`, `✓ world package gate PASSED: 9/9 steps`, `✓ verify gate PASSED` |
| P12 | **`scripts/verify_go.sh` at base** (controller flagged UNVERIFIED) | **GREEN, rc=0** | `GOTOOLCHAIN=go1.25.6 AILANG_BIN=/tmp/ailang-v0300/ailang ./scripts/verify_go.sh` → `VERIFY_GO_BASE_RC=0`, `✓ go gate PASSED: build clean, plain and race tests pass`. The two `WARNING: DATA RACE` lines are the gate's own **race-detector known-positive control** at log lines 13/31, emitted *before* `go build` — healthy, as documented |
| P13 | `rg` is not a binary and appears nowhere in CI or scripts | **CONFIRMED** | `whence -p rg` → `NO-BINARY`; `grep -c 'rg ' .github/workflows/ci.yml scripts/*.sh` → `0` in all six scripts and in `ci.yml` |
| P14 | V24 frozen v1 interface hash | **CONFIRMED, with a discriminating control** | `printf '%s' '<preimage>' \| shasum -a 256` → `743f39f470bf354ebab0ab196598b5ba72db80463d833325cb7672249d4734ac`; negative control (`v1`→`v2`, one byte) → `8e27fdca25327a2f…` — the instrument discriminates |
| P15 | V25 three production `Session.Invoke` sites | **CONFIRMED** | `grep -rn '\.Invoke(' --include='*.go' host/ cmd/ \| grep -v _test.go` → exactly 3, `host/broker/publish_op.go:135,162,279`; `_test.go` control arm → **83** |

### Four measurements the controller did not make, all of which change the plan

#### (i) **A second `store.Open` on one database is REFUSED. The doc's "two handles/racers" concurrency recipe is impossible.**

`host/store` enforces a single-writer model: `Open(path)` on a file-backed DB takes a non-waiting
exclusive OS lock *before* SQLite exists and returns `*WriterAlreadyActive` on contention — and the
landed test `writer_lock_test.go:507-523` proves this holds **in the same process**
(`if second, err := Open(dbPath); err == nil { t.Fatal("a second Open succeeded …") }`). The
in-memory carve-out is no escape either: an in-memory DSN is *per-connection*, so two
`Open(":memory:")` calls are two different databases.

The design doc's non-vacuity note says *"SQLite CAS conflicts are driven by two handles/racers"*
(w-transition-registry.md:617). The **two-handles half of that sentence is unimplementable.**

Measured replacement, on a probe I ran and deleted (tree verified clean afterwards):

```
PROBE-C: 8/8 concurrent SetRegistryHead on ONE handle succeeded
```

**`TestConcurrentPublishHasOneWinner` must use N goroutines on ONE `*Store` handle.** That is not a
weaker test: `db.SetMaxOpenConns(1)` (`store.go:293`) makes the single CAS transaction the sole
serialization point, which is exactly the mechanism under test, and `verify_go.sh`'s `-race` leg
exercises the Go-level half.

#### (ii) **A non-`tx` read inside an open transaction DEADLOCKS. This is the single most likely way to lose a day.**

`PutObject` (`store.go:450`) and `GetObject` (`store.go:463`) both issue their SQL on `s.db`, not on
a `*sql.Tx`. With `SetMaxOpenConns(1)`, an open `*sql.Tx` holds the **only** connection, so any
`s.db` query issued while that tx is open blocks — with no context, **forever**.

Measured, both arms, in one run:

```
PROBE-A: non-tx read INSIDE tx BLOCKED for 3s -> DEADLOCK CONFIRMED
PROBE-B CONTROL: read with no open tx COMPLETED (instrument works)
```

The control fires, so PROBE-A's block is the mechanism and not a broken probe.

**Binding constraint on task T1**: `CompareAndSetRegistryHead`'s *"requires `next` to name an
existing object"* check **must** be a `SELECT` on the `*sql.Tx`. Calling `s.GetObject(next)` from
inside the transaction hangs until `go test`'s 10-minute panic timeout, and the executor will read
that as "the CAS is slow", not as a deadlock. The repo already has the right idiom —
`selectedHeadTx` (`store.go:710`) takes a *querier interface* precisely so it can be driven by
either `s.db` or a `tx`. Follow it.

#### (iii) **`MUT-AIL-EMPTY-MODULE` is VACUOUS as specified, and so is AC9 as written.**

The doc's mutation table (line 607) claims the observable is *"exact module total/allowlist RED"*.
**There is no exact module total.** `scripts/verify_ail.sh:232-236` asserts only `checked -ne 0`
(non-vacuity); the `11` is *printed* at `:243` and never compared. The pinned quantities are
`EXACT_TOTAL_VERIFIED=4` (`:238`) and `EXACT_TOTAL_TESTS=14` (`:268`), neither of which an empty
module moves.

Measured with the mutation **proven to land** (`world/*.ail` 4 → 5, module count 11 → 12):

| arm | `world/*.ail` | modules reported | `AILANG_BIN=… ./scripts/verify_ail.sh` rc |
|---|---:|---:|---:|
| pristine base | 4 | 11 | **0** |
| `MUT-AIL-EMPTY-MODULE` armed | 5 | **12** | **0 — STILL PASSES** |

Both the step-2 *"exact `.ail` allowlist"* and step-8 *"tar contains exactly 6 allowlisted entries"*
gates also passed, because they enumerate `packages/world-core/`, not `world/`.

So **AC9 as written in the doc ("Command: `…verify_ail.sh`; Base observed: rc=0 … Delivered output
has the same 4/11/14 totals") cannot detect the P8 violation it exists to detect.** An executor
checking rc — which is all the command yields — passes a tree that added a `world/*.ail` module.
That is the mission's canonical "a green that can mean the check never ran".

**Repaired AC9, proven in BOTH arms this session:**

```bash
export PATH=/opt/homebrew/bin:$PATH
out=$(AILANG_BIN=/tmp/ailang-v0300/ailang ./scripts/verify_ail.sh 2>&1); rc=$?
m=$(printf '%s\n' "$out" | grep -c '4/4 required world/ identities verified across 11 module(s)')
t=$(printf '%s\n' "$out" | grep -c 'all 14 required named tests pass')
s=$(printf '%s\n' "$out" | grep -c 'world package gate PASSED: 9/9 steps')
[ "$rc" -eq 0 ] && [ "$m" -eq 1 ] && [ "$t" -eq 1 ] && [ "$s" -eq 1 ]
```

| arm | rc | modules11 | tests14 | steps9 | **repaired AC9** |
|---|---:|---:|---:|---:|---:|
| pristine base | 0 | 1 | 1 | 1 | **0 (PASS)** |
| `MUT-AIL-EMPTY-MODULE` armed | 0 | **0** | 1 | 1 | **1 (RED)** |

The discriminating observable is `modules11`. Tree restored and verified clean
(`git status --porcelain` → empty; `ls world/*.ail | wc -l` → 4).

#### (iv) **`MUT-CAS-EPOCH-HEAD`'s first offered observable does not fire.**

The doc offers *"epoch round-trip **or** CAS isolation test RED"*. The epoch round-trip test is
`TestRegistryHeadRoundTrip`, which calls only `SetRegistryHead`/`GetRegistryHead`. Hardcoding the
epoch registry name inside the **new** `CompareAndSetRegistryHead` does not touch `SetRegistryHead`,
so `TestRegistryHeadRoundTrip` stays **GREEN**. See §4 for the required observable.

### One non-finding worth recording, because a reviewer will ask

**TR-CJSON-1 is not a duplicate of anything landed.** Repo-wide search for an existing canonical
JSON facility found exactly two, and both are refuted as reuse candidates by inspection:

- `host/canon.Source` (`host/canon/source.go:51`) is the canonical **transition-source-byte**
  procedure — it has no JSON key sorting, no `UseNumber`, no duplicate detection
  (`grep -n 'sort\|UseNumber\|json.Marshal\|Decoder\|SetEscapeHTML' host/canon/*.go` → **0 hits**).
- `host/registry.Registry.Encode` (`registry.go:61-70`) is deterministic *because Go emits struct
  fields in declaration order* — it explicitly does **not** sort keys lexicographically, which is
  TR-CJSON-1's requirement, and it round-trips through `json.Unmarshal`, which Decision 4 forbids.
- `host/store/journal.go:186 canonicalJSON` is `json.Marshal` + `canon.Source` over **fixed Go
  structs**. V20 measured `encoding/json`'s value path accepting duplicate keys, lone surrogates and
  invalid UTF-8 — fine for a struct, unusable for caller-supplied schema bytes.

The Systemic-Issue Audit row *"Same concept duplicated: prevented"* is therefore **sustained**.

---

## 1. What TR.A is, and what this plan explicitly does NOT do

**In scope** (the doc's Files table, TR.A rows only):

| File | Purpose |
|---|---|
| `host/transitionreg/transitionreg.go` | `Descriptor`, `EffectRequirement`, `Reader`, `Snapshot`, `BuildNext`, `Publish` |
| `host/transitionreg/codec.go` | TR-CJSON-1 canonical encoding + validation |
| `host/transitionreg/transitionreg_test.go` | identity, codec, snapshot, publication tests |
| `host/store/store.go` | `CompareAndSetRegistryHead` — **one method, no table** |
| `host/store/store_test.go` | CAS conflict, absent-head, rollback, concurrency, isolation |
| `design_docs/planned/w-transition-registry.md` | **T8 only**: the zero-tolerance AC activation |

**Explicitly NOT in this sprint** — if the executor touches any of these, the milestone is wrong:

- no `host/broker` change of any kind (that is TR.B `broker.Allows` / `CapabilitySnapshot`, and TR.C's `invoke_boundary_test.go`);
- no `transitionreg.Bind`, no guarded invoker, no `DeclaredEffects` *enforcement* (TR.B) — the
  `DeclaredEffects` **field** and its **validation** are TR.A, its **enforcement** is not;
- no daemon, coordinator, `cmd/`, or protocol change;
- no `world/*.ail` addition (P8) — except transiently, under `MUT-AIL-EMPTY-MODULE`, restored in the same task;
- no `host/store/schema.sql` change — **the file is under a DDL gate**; TR.A adds zero tables;
- no CI workflow change;
- no `host/replay` change (AC8 must not move).

---

## 2. AC reconciliation — plan tasks diffed against the doc's AC list

The doc has 11 ACs. TR.A owns 5 and must hold 2 others still.

| AC | Owner | TR.A duty | Closing task |
|---|---|---|---|
| **AC1** identity, codec, schema validation | **TR.A** | **CLOSES** — activate to exactly **3** tests | T2, T3 → activated T8 |
| **AC2** eager one-head snapshot + copy isolation | **TR.A** | **CLOSES** — activate to exactly **3** | T5 → activated T8 |
| **AC3** CAS publication + deterministic ordering | **TR.A** | **CLOSES** — activate to exactly **4** | T6 → activated T8 |
| **AC4** generic store CAS preserves epoch registry | **TR.A** | **CLOSES** — activate from base **1** to exactly **2** | T1 → activated T8 |
| AC5 capability snapshot / denial arms | TR.B | untouched; base count 0 must remain 0 | — |
| AC6 declaration-honesty mechanism | TR.B | untouched; base count 0 must remain 0 | — |
| AC7 two-session fixture | TR.B | untouched; base count 0 must remain 0 | — |
| **AC8** replay stays hash-pinned | TR.A holds | **MUST NOT MOVE** — no `host/transitionreg` import in production replay code | verified every task |
| **AC9** AILANG gate totals | TR.A holds | **MUST NOT MOVE** — 4/11/14; **use the repaired form, §0(iii)** | verified every task; T7 mutation |
| **AC10** build + focused packages | **TR.A** | **CLOSES** — activate to exactly **1** listed + whole-package run | T3 → activated T8 |
| AC11 structural binding boundary | TR.C | untouched; base count **1** must remain 1 | — |

**TR.A closes AC1, AC2, AC3, AC4, AC10. It holds AC8 and AC9 unmoved. It leaves AC5, AC6, AC7 and
AC11 at their base counts (0, 0, 0, 1 respectively) — an executor that moves any of those four has
built TR.B or TR.C work and is out of scope.**

### Required test inventory after activation (the exact names the ACs list)

| AC | count | test names |
|---|---:|---|
| AC1 | 3 | `TestDescriptorIdentityAndContentUpdate`, `TestCodecGoldenRoundTrip`, `TestDescriptorValidationRefusals` |
| AC2 | 3 | `TestReadSnapshotReadsHeadOnce`, `TestSnapshotIsEagerAndCopyIsolated`, `TestReadSnapshotRefusals` |
| AC3 | 4 | `TestPublishCASConflictPreservesWinner`, `TestConcurrentPublishHasOneWinner`, `TestStableIDByteOrder`, `TestPublishRefusals` |
| AC4 | 2 | `TestRegistryHeadRoundTrip` (existing), `TestCompareAndSetRegistryHead` (new) |
| AC10 | 1 | `TestCodecGoldenRoundTrip` (listed) + the whole `./host/transitionreg` package runs |

### Doc contradictions found (reported, not smoothed over)

1. **Line 617 vs `writer_lock_test.go:507`** — "SQLite CAS conflicts are driven by two
   handles/racers" is unimplementable for two writer handles. §0(i). *Plan overrides the doc.*
2. **Line 607 vs `verify_ail.sh:232-243`** — `MUT-AIL-EMPTY-MODULE`'s stated observable
   ("exact module total/allowlist") does not exist, and AC9 is vacuous in consequence. §0(iii).
   *Plan repairs AC9 with a two-arm-measured command.*
3. **Line 590 vs `store_test.go:260`** — `MUT-CAS-EPOCH-HEAD`'s "epoch round-trip … RED" arm does
   not fire. §0(iv). *Plan pins the isolation observable instead.*
4. **Line 616 — "Every designed refusal is represented above. There is no claimed unreachable
   branch."** This is **FALSE** under rule 3j. Eleven refusal branches the design freezes have no
   named mutation. §4.2. *Plan adds them.*
5. **V13 (line 659) records `verify_go.sh` as rc=1.** That was the designer's `workspace-write`
   sandbox denying loopback binds, as the doc's own V18 half-notes for `host/daemon`. First-party,
   outside the sandbox, the whole gate is **rc=0** (P12). The doc leaves its log "as recorded" by
   policy; the *plan's* baseline is rc=0 and that is the number the executor must reproduce.

---

## 3. Day-by-day breakdown — 8 tasks, 8 commits

One task = one commit boundary. **The controller makes every commit**; the executor stops at each
boundary and reports. Every task ends by re-running the **hold set** (AC8, AC9-repaired, AC10 build
arm) — a task that closes its own AC while moving AC9 has not succeeded.

### Day 1 — store primitive, codec, descriptor

#### T1 — `host/store.CompareAndSetRegistryHead` + tests  (~140 impl / ~230 test LOC) → closes AC4

Add to `host/store/store.go`, immediately after `GetRegistryHead` (`:624`):

```go
// RegistryCASConflict is returned when the observed head is stale.
type RegistryCASConflict struct {
    Name     string
    Expected hashref.HashRef // zero means "expected absent"
    Actual   hashref.HashRef // zero means "was absent"
    HadHead  bool
}
func (e *RegistryCASConflict) Error() string { … }
func IsRegistryCASConflict(err error) bool { … }   // mirrors the landed IsConflict idiom (store.go:170)

// CompareAndSetRegistryHead sets name's head to next in ONE transaction, but only
// if the current head equals expected. A zero `expected` means "the name must have
// no head". `next` must name an existing object. Only the row for `name` changes.
func (s *Store) CompareAndSetRegistryHead(name string, expected, next hashref.HashRef) error
```

Implementation constraints, all load-bearing:

- `tx, err := s.db.Begin()`; `defer tx.Rollback()`; commit on the success path only.
- **Every read inside the tx goes through `tx.QueryRow`.** Do NOT call `s.GetObject`,
  `s.GetRegistryHead`, or any other `*Store` method between `Begin` and `Commit` — §0(ii), measured
  deadlock. Extract `registryHeadTx(q querier, name string)` in the `selectedHeadTx` (`:710`) style
  if a shared reader is wanted.
- Object-existence check: `SELECT 1 FROM objects WHERE hash_ref = ?` **on `tx`**.
- `validateRef("CompareAndSetRegistryHead", "next", next)` up front, matching `SetRegistryHead:607`.
  `expected` may be zero and is therefore validated only when non-zero.
- The `UPDATE`/`INSERT` names `name` explicitly. **Never a literal `EpochRegistryV1`.**
- Add `const TransitionRegistryV1 = "world/transition-registry/v1"` beside `EpochRegistryV1`
  (`store.go:82`) so the name is one shared constant, not two string literals.
- `SetRegistryHead` is untouched — the doc keeps it for the epoch registry.

New test `TestCompareAndSetRegistryHead` in `host/store/store_test.go`, using the existing
`openMem(t)` (`store_test.go:11`) and `obj(payload, semanticID)` (`:23`) helpers, with these named
subtests (each is a mutation target in §4):

| subtest | asserts |
|---|---|
| `absent_head_accepts_zero_expected` | zero `expected` + no head → success; head reads back as `next` |
| `absent_head_rejects_nonzero_expected` | zero head + non-zero `expected` → `*RegistryCASConflict`, `HadHead=false` |
| `stale_expected_conflicts` | head=A, expected=B → conflict carrying **both** heads; head still A |
| `dangling_next_refused` | `next` not in `objects` → error; head unchanged |
| `rollback_leaves_head_unchanged` | after every refusal above, `GetRegistryHead` is byte-identical to before |
| `epoch_registry_isolation` | CAS on `TransitionRegistryV1` leaves `GetRegistryHead(EpochRegistryV1)` **unchanged**, and vice versa — **this is the only observable that kills `MUT-CAS-EPOCH-HEAD`** (§0(iv)) |
| `concurrent_racers_one_winner` | 8 goroutines, ONE `*Store` handle, same `expected` → exactly 1 nil error and 7 `*RegistryCASConflict`; final head is the winner's `next` |

**Exit gate**: AC4 command → `count=2`, both tests PASS. Hold set green.

#### T2 — `host/transitionreg/codec.go`: TR-CJSON-1  (~470 LOC) → AC1 (codec half)

Dependency-free, `encoding/json.Decoder` token stream with `UseNumber`. `json.Unmarshal` and
`json.Marshal` are **forbidden** for canonicalization (Decision 4; V20 is the measurement that
justifies it).

Deliverables, each a separately mutatable branch:

- `validateUTF8`, duplicate-member rejection **at every depth**, escaped/literal surrogate
  rejection, root-must-be-object.
- Lexicographic key sort by **unescaped UTF-8 bytes**; arrays keep order; strings **not**
  normalized (NFC ≠ NFD is a deliberate distinction and gets a test).
- Escape set exactly `"`, `\`, U+0000–U+001F with lowercase `\u00xx` and the five short escapes;
  **no** HTML escaping of `<`, `>`, `&`.
- Arbitrary-precision decimal normalization: no `+`, no exponent, no leading integer zero, no
  trailing fractional zero, `-0` → `0`. Coefficient ≤ 1024 digits, |exponent| ≤ 10 000 **before**
  expansion, normalized token ≤ 16 384 bytes **after**.
- Bounds: schema raw ≤ 262 144 pre-parse, canonical ≤ 65 536 post; revision raw ≤ 16 777 216
  pre-decode, canonical ≤ 8 388 608 post; entries ≤ 1024. **All inclusive**; all enforced *before*
  hashing or store write.
- `Encode → Decode → Encode` is byte-identical.

`const SemanticIDV1 = "world/transition-registry/v1"` and the frozen interface hash live here.

> **Anti-vacuity requirement on the golden fixtures.** `TestCodecGoldenRoundTrip` must compare
> against a **committed literal** — the golden bytes as a literal string and the interface hash as
> the literal `"sha256:743f39f470bf354ebab0ab196598b5ba72db80463d833325cb7672249d4734ac"`. It must
> **never** derive the expectation from `SemanticIDV1` or from the code's own preimage builder. A
> fixture computed from the constant it is meant to freeze is set *alongside* the mechanism, cannot
> fail for the reason it claims (rule 3i), and would make `MUT-GO-CODEC-TAG` pass. The preimage is
> re-derivable independently: P14 above.

**Exit gate**: package builds; hold set green. (AC1 is not yet complete — T3 finishes it.)

#### T3 — `host/transitionreg/transitionreg.go`: descriptor, revision, ordering  (~300 LOC) → closes AC1, AC10

```go
type EffectRequirement struct{ Effect, Scope string; Cost int64 }
type Descriptor struct {
    ID string
    TransitionFn, Interpreter hashref.HashRef
    SemanticsEpoch int64
    InputSchema, OutputSchema []byte   // TR-CJSON-1 canonical
    Access EffectRequirement
    DeclaredEffects []EffectRequirement
    Title, Description string
}
type Revision struct{ SemanticID string; InterfaceHash hashref.HashRef; Revision int64; Parent hashref.HashRef; Entries []Descriptor }
```

ID grammar exactly `^[a-z0-9](?:[a-z0-9_-]{0,30}[a-z0-9])?(?:[./][a-z0-9](?:[a-z0-9_-]{0,30}[a-z0-9])?)*$`,
total 1–128 bytes, every `.`/`/` segment 1–32 bytes, ASCII so bytes == chars. Ordering is unsigned
bytewise; on a prefix the shorter sorts first.

Tests: `TestDescriptorIdentityAndContentUpdate` (same ID + new `TransitionFn` ⇒ new revision, old
revision unchanged and still addressable), `TestCodecGoldenRoundTrip` (incl. the **empty revision 1**
golden — the Open Decision's recommended default), `TestDescriptorValidationRefusals` (table-driven,
one named subtest per refusal branch in §4.2).

> **P1 freeze check, mechanical**: `grep -c '^type Registry' host/transitionreg/*.go` must be **0**.

**Exit gate**: AC1 → `count=3`, all PASS. AC10 → `count=1`, whole package PASS. Hold set green.

#### T4 — mutation sweep for T1–T3  (~26 arms, no production LOC)

Run the §4 protocol for every mutation whose mechanism landed in T1–T3. Deliverable is a
`design_docs/verification/w-transition-registry/tra-mutations.md` transcript recording, per arm:
anchor count before/after, differing sha256, `go build` rc, **the failing test name**, and the
inverse `-skip` arm rc=0.

### Day 2 — snapshot, publication, activation

#### T5 — `Reader` / `ReadSnapshot` / `Snapshot` + cache  (~300 LOC) → closes AC2

```go
type ObjectStore interface {                      // <- the seam MUT-READ-SWALLOW needs
    GetRegistryHead(name string) (hashref.HashRef, bool, error)
    GetObject(hashref.HashRef) (store.Object, bool, error)
    PutObject(store.Object) error
    CompareAndSetRegistryHead(name string, expected, next hashref.HashRef) error
}
type Reader interface{ ReadSnapshot(context.Context) (Snapshot, error) }
type Snapshot struct{ Head hashref.HashRef; Revision int64; entries []Descriptor }
func (Snapshot) Lookup(id string) (Descriptor, bool)   // binary search
func (Snapshot) List() []Descriptor                     // deep copy
```

`*store.Store` satisfies `ObjectStore` structurally. **This interface is mandatory**: without it,
`MUT-READ-SWALLOW` and `MUT-PUBLISH-SWALLOW` have no injection point and `MUT-SNAPSHOT-REREAD` has
no call counter. The doc says store errors are "driven by injected interfaces" but never names the
interface; this is that name.

`ReadSnapshot` does exactly the doc's five steps, in order, reading the head **once**. The
parsed-snapshot cache is keyed by the immutable head hash, is mutex-guarded, and **always** reads
the head first and **always** returns a deep copy.

Tests: `TestReadSnapshotReadsHeadOnce` (counting fake asserts head reads == 1 per call, and == 2
across two calls even on a cache hit), `TestSnapshotIsEagerAndCopyIsolated` (mutate every returned
slice and every schema byte slice; re-read; unchanged — covering construction, `List`, **and**
`Lookup`), `TestReadSnapshotRefusals` (table-driven, one subtest per §4.2 read branch).

**Exit gate**: AC2 → `count=3`, all PASS. Hold set green.

#### T6 — `BuildNext` / `Publish`  (~250 LOC) → closes AC3

`BuildNext(current Revision, changes []Change) (Revision, error)` — pure, side-effect free: validate
every candidate, apply replace/remove by stable ID, sort bytewise, set `Revision = current.Revision+1`
and `Parent = <captured head>`, encode.

`Publish(ctx, expectedHead hashref.HashRef, next Revision) (hashref.HashRef, error)` — `PutObject`
the immutable object, then `CompareAndSetRegistryHead(TransitionRegistryV1, expectedHead, objRef)`.
A failed CAS leaves the orphan object in place and returns the typed conflict **unwrapped enough
that `store.IsRegistryCASConflict` still reports true**.

Tests: `TestPublishCASConflictPreservesWinner`, `TestConcurrentPublishHasOneWinner` (**N goroutines,
ONE handle** — §0(i)), `TestStableIDByteOrder` (incl. the prefix case `a` < `ab` and a full unsigned
byte sweep), `TestPublishRefusals` (table-driven).

**Exit gate**: AC3 → `count=4`, all PASS. Hold set green.

#### T7 — mutation sweep for T5–T6, plus `MUT-GO-CODEC-TAG` and `MUT-AIL-EMPTY-MODULE`  (~14 arms)

`MUT-AIL-EMPTY-MODULE` is run against the **repaired** AC9 of §0(iii), never the doc's rc-only form.
Restore is `rm -f world/transitionregistry.ail` followed by `ls world/*.ail | wc -l` → `4`.

#### T8 — **ZERO-TOLERANCE ACTIVATION** (the merge criterion) → the milestone is not done without it

This is a *checkable task*, not a footnote. A count gate that still accepts 0 is satisfied by
deleting the tests.

**T8.a** — Edit `design_docs/planned/w-transition-registry.md` ACs 1, 2, 3, 4, 10 to delete the
base-tolerant arm and require the exact count:

| AC | remove | becomes |
|---|---|---|
| AC1 | `test "$count" -eq 0 \|\|` | `test "$count" -eq 3 && go test … -run … -count=1` |
| AC2 | `test "$count" -eq 0 \|\|` | `test "$count" -eq 3 && …` |
| AC3 | `test "$count" -eq 0 \|\|` | `test "$count" -eq 4 && …` |
| AC4 | `test "$count" -eq 1 \|\|` | `test "$count" -eq 2 && …` |
| AC10 | `test "$count" -eq 0 \|\|` | `test "$count" -eq 1 && go test ./host/transitionreg -count=1` |

**T8.b** — Machine check that no tolerant arm survives:

```bash
export PATH=/opt/homebrew/bin:$PATH
awk '/^1\. \*\*AC1/,/^5\. \*\*AC5/' design_docs/planned/w-transition-registry.md \
  | grep -c 'test "\$count" -eq 0'          # must be 0
```
with a **known-positive control in the same call** — the same `grep -c` over the AC5–AC7 range must
return **3** (TR.B's arms are untouched), proving the instrument can still see a tolerant arm.

**T8.c** — Record `MUT-DELETE-TR-A-TEST` RED. Delete **one** required test (do all five arms, one
per AC), re-run that AC, require rc≠0 **and** record the count observed (e.g. AC1 `count=2`, not
just "rc=1"). Restore from the `cp` backup and re-run to green.

**T8.d** — Also update AC9 in the doc to the repaired form of §0(iii), and correct the mutation-table
row for `MUT-AIL-EMPTY-MODULE` from "exact module total/allowlist RED" to the measured observable.

**T8.e** — Final full gate: `./scripts/verify_ail.sh` (repaired form) **and**
`./scripts/verify_go.sh`, both **outside** any sandbox, both rc=0.

---

## 4. Mutation discipline

### 4.1 Protocol — every arm, no exceptions

```
1. cp <file> /tmp/tra_backup/<file>.bak
2. record: anchor grep -c  AND  shasum -a 256 <file>      # BEFORE
3. apply the neuter  (if false && <cond>   or   _ = f(x)) — NEVER delete a block:
   a deleted block orphans an import and reds the BUILD, which is the colour you
   predicted for the wrong reason.
4. assert the mutation LANDED: anchor count changed AND sha256 DIFFERS from step 2.
   If either is unchanged, STOP — "did not red" and "never ran" are the same rc.
5. assert the mutant BUILDS: GOTOOLCHAIN=go1.25.6 go build ./...   rc MUST be 0.
6. kill arm:    go test <pkg> -run '<the named test>' -count=1     -> expect rc≠0
   RECORD THE FAILING TEST/SUBTEST NAME. rc=1 whose only FAIL is a pre-existing
   flake is not a kill.
7. inverse arm: go test <pkg> -skip '<the named test>' -count=1    -> MUST be rc=0.
   This is what proves YOUR test is the killer and not a bystander.
8. restore: cp /tmp/tra_backup/<file>.bak <file>   -- NEVER `git checkout -- <file>`;
   these files are uncommitted by construction in a sprint worktree and git checkout
   DELETES the executor's work.
9. assert the restore: sha256 equals step 2 exactly.
```

### 4.2 Mutation assignment — the 23 doc-named TR.A mutations

`arms` counts the "one at a time" splits the doc itself requires.

| Mutation | Task | arms | **Exact required observable** | rule 3i |
|---|---|---:|---|---|
| `MUT-ID-ACCEPT-EMPTY` | T4 | 1 | `TestDescriptorValidationRefusals/id_grammar` FAILs | downstream ✓ |
| `MUT-ID-ZERO-FN` | T4 | 1 | `…/zero_transition_fn` FAILs | ✓ |
| `MUT-ID-ZERO-INTERP` | T4 | 1 | `…/zero_interpreter` FAILs | ✓ |
| `MUT-SCHEMA-ANY-BYTES` | T4 | 1 | `…/schema_not_an_object` FAILs | ✓ |
| `MUT-SCHEMA-NO-LIMIT` | T4 | **2** | (a) `…/schema_raw_over_262144` (b) `…/schema_canonical_over_65536` | ✓ |
| `MUT-CODEC-INDENT` | T4 | 1 | `TestCodecGoldenRoundTrip` FAILs on a **committed literal** golden; `go build` rc=0 | ✓ **only if** the golden is a literal — see T2 |
| `MUT-GO-CODEC-TAG` | T7 | 1 | `TestCodecGoldenRoundTrip` **and** the interface-hash assertion FAIL; `go build ./...` rc=0 | ✓ **only if** the golden/hash are literals |
| `MUT-READ-EMPTY-OK` | T7 | 1 | `TestReadSnapshotRefusals/absent_head` FAILs | ✓ |
| `MUT-READ-SWALLOW` | T7 | 1 | `…/injected_read_error` FAILs (needs the `ObjectStore` seam, T5) | ✓ |
| `MUT-READ-ABSENT-OK` | T7 | 1 | `…/object_absent` FAILs | ✓ |
| `MUT-READ-NO-REHASH` | T7 | 1 | `…/corrupt_object_payload` FAILs | ✓ |
| `MUT-READ-ANY-TYPE` | T7 | **2** | (a) `…/wrong_semantic_id` (b) `…/wrong_interface_hash` | ✓ |
| `MUT-SNAPSHOT-ALIAS` | T7 | 1 | `TestSnapshotIsEagerAndCopyIsolated` FAILs — must cover construction, `List` **and** `Lookup` | ✓ |
| `MUT-SNAPSHOT-REREAD` | T7 | 1 | the counting fake's head-read count ≠ 1 → `TestReadSnapshotReadsHeadOnce` FAILs | ✓ counter lives in the fake, not the SUT |
| `MUT-SNAPSHOT-CACHE-BYPASS` | T7 | **2** | (a) head-read count stays 1 across two calls (b) copy-isolation FAILs | ✓ |
| `MUT-REVISION-SKIP` | T7 | **2** | (a) `TestPublishRefusals/revision_not_n_plus_1` (b) `…/parent_not_captured_head` | ✓ |
| `MUT-ORDER-INSERTION` | T7 | **2** | (a) `TestStableIDByteOrder` FAILs (b) `TestPublishRefusals/duplicate_id` FAILs | ✓ |
| `MUT-CAS-BLIND` | T4 | 1 | `TestCompareAndSetRegistryHead/stale_expected_conflicts` FAILs **and** `TestPublishCASConflictPreservesWinner` FAILs (T7 re-run) | ✓ head is read back after the CAS |
| `MUT-CAS-DANGLING` | T4 | 1 | `…/dangling_next_refused` FAILs | ✓ |
| `MUT-PUBLISH-SWALLOW` | T7 | **2** | (a) injected `PutObject` error (b) injected CAS error — each named subtest FAILs | ✓ |
| `MUT-CAS-EPOCH-HEAD` | T4 | 1 | **`TestCompareAndSetRegistryHead/epoch_registry_isolation` FAILs.** The doc's alternative "epoch round-trip test RED" is **REJECTED**: `TestRegistryHeadRoundTrip` never calls the new method and stays GREEN — §0(iv) | **rule 3i CATCH** |
| `MUT-AIL-EMPTY-MODULE` | T7 | 1 | **repaired AC9 → `modules11=0`.** The doc's "exact module total/allowlist RED" is **REJECTED**: measured rc=0 in both arms — §0(iii) | **rule 3i CATCH** |
| `MUT-DELETE-TR-A-TEST` | T8.c | **5** | one arm per activated AC; record the *count* observed, not just rc | ✓ |

**Doc-named subtotal: 23 mutations / 34 arms.**

### 4.3 Rule 3j — eleven refusal branches the doc freezes but never mutates

The doc asserts (line 616) *"Every designed refusal is represented above."* It is not. Each branch
below is specified in Decision 2 or Decision 4, will be written by T2/T3, and has **no** named
mutation. Under rule 3j the unit of mutation is the branch, so each needs its own arm. They are all
cheap `if false &&` neuters on code the executor is writing anyway.

| # | Frozen branch (doc reference) | Proposed mutation | Observable |
|---|---|---|---|
| 1 | ID length bounds 1–128 total / 1–32 per segment (Decision 1) | `MUT-ID-NO-LENGTH-BOUND` | `TestDescriptorValidationRefusals/id_too_long`, `…/segment_too_long` |
| 2 | `semanticsEpoch` non-negative (Decision 2) | `MUT-NEG-EPOCH-OK` | `…/negative_semantics_epoch` |
| 3 | requirement `cost` non-negative (Decision 2/5) | `MUT-NEG-COST-OK` | `…/negative_cost` |
| 4 | duplicate object member names at every depth (Decision 4) | `MUT-CJSON-DUP-KEY-OK` | `…/duplicate_schema_key_nested` |
| 5 | invalid UTF-8 / escaped or literal surrogate (Decision 4) | `MUT-CJSON-SURROGATE-OK` | `…/lone_surrogate`, `…/invalid_utf8` |
| 6 | number coefficient ≤1024 digits, \|exp\| ≤10 000, token ≤16 384 (Decision 4) | `MUT-CJSON-NO-NUMBER-BOUND` | `…/number_coefficient_overflow` |
| 7 | unknown or missing revision/descriptor key (Decision 2) | `MUT-CJSON-UNKNOWN-KEY-OK` | `…/unknown_revision_key`, `…/missing_descriptor_key` |
| 8 | entries ≤ 1024 (Decision 4) | `MUT-ENTRIES-NO-LIMIT` | `TestPublishRefusals/entries_over_1024` |
| 9 | revision raw ≤16 777 216 / canonical ≤8 388 608 (Decision 4) | `MUT-REVISION-NO-LIMIT` | `TestReadSnapshotRefusals/revision_raw_over_limit` |
| 10 | lexicographic key sort by unescaped UTF-8 bytes (Decision 4) | `MUT-CODEC-NO-KEYSORT` | `TestCodecGoldenRoundTrip` FAILs on the literal golden |
| 11 | decimal normalization `1`/`1.0`/`1e0` → `1`, `-0` → `0` (Decision 4) | `MUT-CODEC-NUMBER-RAW` | `TestCodecGoldenRoundTrip` FAILs; `…/number_spellings_collapse` |

**Rule-3j subtotal: 11 mutations / 13 arms.** Grand total **34 mutations / 47 arms**.

At the measured cost of a scoped `go test ./host/transitionreg` on this rig (sub-second; the whole
`host/store` package is 3.2 s), 47 arms scripted through §4.1 is roughly 1.5–2 hours of executor
wall time. That fits inside T4 + T7. **This is the honest cost of rule 3j and it is not optional:
if the controller wants to cut it, the cut must be a recorded decision, not a silent omission.**

### 4.4 One row I am flagging as UNCERTAIN

`MUT-CODEC-INDENT` and `MUT-CODEC-NO-KEYSORT` both kill via `TestCodecGoldenRoundTrip`. If the
executor implements the golden as *"encode, then decode, then encode, assert the two encodings
match"* — a pure round-trip with no committed literal — then **both mutations pass**, because an
indented or unsorted encoder is still idempotent with its own decoder. The round-trip assertion is
set alongside the mechanism, not downstream of it. The test therefore needs **two** independent
assertions: (a) byte-equality against a committed literal golden, and (b) round-trip idempotence.
Only (a) kills these two mutations. The executor must report which assertion fired.

---

## 5. Acceptance commands, as the executor must run them

All baselined this session on the pristine tree at `adfaa0b`; the base result is part of each
criterion. **No command contains `rg`** (P13).

| AC | base (measured) | TR.A delivered |
|---|---|---|
| AC1 | `count=0`, rc=0 | `count=3`, three tests PASS |
| AC2 | `count=0`, rc=0 | `count=3`, three tests PASS |
| AC3 | `count=0`, rc=0 | `count=4`, four tests PASS |
| AC4 | `count=1` (`TestRegistryHeadRoundTrip`), rc=0 | `count=2`, both PASS |
| AC8 | rc=0, `host/replay` PASS | rc=0, unchanged; no `transitionreg` import in production replay |
| AC9 **(repaired)** | rc=0, `modules11=1 tests14=1 steps9=1` | identical |
| AC10 | build PASS, store PASS, broker PASS, `count=0` | `count=1`, whole `./host/transitionreg` PASS |
| AC5/6/7 | `count=0` each | **still 0** — moving these is out of scope |
| AC11 | `count=1` | **still 1** |
| `verify_go.sh` | **rc=0**, `✓ go gate PASSED` | rc=0 |

**AC9 is used in the repaired form only.** The doc's rc-only form is vacuous (§0(iii)); the executor
who runs it as printed will report a pass that means nothing.

---

## 6. Estimates

| Task | impl LOC | test LOC | notes |
|---|---:|---:|---|
| T1 store CAS | 140 | 230 | 7 named subtests, one goroutine racer |
| T2 codec | 470 | — | TR-CJSON-1 is the bulk of the milestone |
| T3 descriptor/revision | 300 | 380 | incl. golden literals |
| T4 mutations (T1–T3) | 0 | 0 | ~26 arms + transcript |
| T5 reader/snapshot/cache | 300 | 300 | `ObjectStore` seam + counting fake |
| T6 buildnext/publish | 250 | 260 | |
| T7 mutations (T5–T6, +2) | 0 | 0 | ~16 arms |
| T8 activation | 0 | 0 | doc edit + 5 delete arms + full gates |
| **total** | **1460** | **1170** | **≈2630 LOC over 2 days** |

Repo velocity reference: the `VL.B` plan estimated 515 LOC at 0.5 day (~1000 LOC/day). 2630 LOC over
2 days is ~1300/day — **above** recent velocity. See risk R1.

---

## 7. Execution protocol

- **Worktree**: a sibling of the repo, e.g. `/Users/voightkampff/dev/sunholo-data/.wt-iter71`.
  **Never under `/tmp`** — cwd-relative path tests then fail for the location rather than the code,
  and a red CI can never reproduce it.
- **Every** Bash call starts `export PATH=/opt/homebrew/bin:$PATH`. Without it `go`, `gh` and `node`
  are rc=127 and look like a broken toolchain.
- **Every** `go` invocation carries `GOTOOLCHAIN=go1.25.6`. `verify_go.sh` rejects the ambient
  toolchain, so a gate already red at base measures the repo, not your change.
- `AILANG_BIN=/tmp/ailang-v0300/ailang` (v0.30.0), **outside** the repo. Never the PATH `ailang` —
  it is a `-dirty` dev build and CLAUDE.md forbids it.
- **zsh, not bash**, in this harness: `${PIPESTATUS[0]}` is empty (zsh spells it `${pipestatus[1]}`,
  1-indexed); an unquoted glob-shaped flag value (`--include=*.go`) aborts the command before it
  runs; zsh does **not** word-split an unquoted variable, so `cmd $FILES` passes ONE argument — use
  arrays and assert `${#arr[@]}`.
- **`rg` is not a binary.** Never put it in an acceptance command, a script, or a mutation
  observable. Use `grep`.
- **Restores are `cp` from `/tmp/tra_backup/`.** `git checkout -- <file>` deletes uncommitted work.
- **SANDBOX CAVEAT — read this before reporting any gate result.** A gate verdict obtained inside a
  `workspace-write` sandbox is **UNINFORMATIVE — neither a pass nor a fail.** Loopback socket binds
  are denied there, which both *invents* failures (`host/daemon`, `host/broker` red for
  `bind: operation not permitted`) and *hides* real ones. This is exactly what produced the design
  doc's V13 row claiming `verify_go.sh` rc=1, which is **rc=0** outside the sandbox (P12). Report
  sandbox results as "sandbox, uninformative"; **the controller re-runs every gate outside the
  sandbox** and that run is the verdict.

---

## 8. Risks — R1 and R4 want a controller decision BEFORE execution

| # | Risk | Assessment |
|---|---|---|
| **R1** | **2630 LOC over 2 days is ~1300 LOC/day, above measured velocity.** TR-CJSON-1 (T2) is a from-scratch canonical JSON profile with arbitrary-precision decimal normalization — the single largest and least compressible item. | **DECISION WANTED.** Options: (a) run TR.A as planned and accept a possible 2.5-day slip; (b) split at the T4/T5 boundary into TR.A1 (store CAS + codec + descriptor, closes AC1/AC4/AC10) and TR.A2 (snapshot + publication, closes AC2/AC3), each independently mergeable and each with its own zero-tolerance activation. **(b) is my recommendation** — the doc's own merge criterion is per-AC, so the split is clean, and it keeps each milestone ≤2 days as the doc requires. |
| **R2** | The `-race` leg of `verify_go.sh` has a hard 600 s kill (`verify_go.sh` python block). Base race-leg package time is **229 s of 600 s**. A concurrency test with goroutines and sleeps could eat the margin. | Mitigated: no sleeps, no wall-clock waits; racers synchronise on a `sync.WaitGroup` + closed-channel start barrier. Budget the new package at **< 5 s** race-leg. Executor must report `ok host/transitionreg <time>` from the race leg. |
| **R3** | The tx deadlock (§0(ii)) presents as a 10-minute hang, not an error, and reads as "slow CAS". | Mitigated by the explicit T1 constraint plus the measured probe. If any `go test ./host/store` exceeds 60 s, suspect a `s.db` call inside the tx **first**. |
| **R4** | **T8 edits the design doc's ACs.** `design_docs/planned/w-transition-registry.md` is a quorum-cleared reviewed artifact and this plan also proposes correcting three of its rows (AC9's form, and the `MUT-AIL-EMPTY-MODULE` / `MUT-CAS-EPOCH-HEAD` observables). | **DECISION WANTED.** The zero-tolerance activation is *mandated by the doc itself*, so T8.a/b/c are uncontroversial. T8.d is a **planner-initiated correction of a reviewed artifact** on the strength of a two-arm measurement. I recommend it land, but the controller should decide whether it goes in the TR.A PR or a separate doc-correction commit. |
| **R5** | `host/store` is one of the 4 protected `protectedGoGroups` in the boundary gate, and TR.A modifies `host/store/store.go`. | **Measured non-risk.** The group enumeration is explicit and pinned at 4 (`allowlist_world_test.go:882`); `host/transitionreg` is not swept, adds no group, and the CAS method adds no import beyond `database/sql`, already in the closure. The forbidden prefix is `…/host/registry`, which does not prefix-match `host/transitionreg`. |
| **R6** | Adding a `host/` package could trip an exact-count gate elsewhere. | **Measured non-risk.** `host/boundary`'s `wantFileCount = 1` (`:1163`) counts `host/boundary`'s own `.go` files. No repo-wide package-count assertion exists. |
| **R7** | AC5–AC7 and AC11 could drift if the executor "helpfully" starts TR.B/TR.C. | Their base counts (0/0/0/1) are in the §5 table as **hold** criteria and must be re-measured at T8.e. |

---

## 9. Handoff

- **Sprint plan**: `design_docs/planned/w-transition-registry-tra-sprint-plan.md` (this file)
- **Sprint JSON**: `.ailang/state/sprints/w-transition-registry-tra.plan.json`
- **Design doc**: `design_docs/planned/w-transition-registry.md`
- **Base**: `dev` @ `adfaa0b`
- Neither artifact is committed by the planner. **The controller commits.**

SPRINT_PLAN_PATH: design_docs/planned/w-transition-registry-tra-sprint-plan.md
SPRINT_JSON_PATH: .ailang/state/sprints/w-transition-registry-tra.plan.json
