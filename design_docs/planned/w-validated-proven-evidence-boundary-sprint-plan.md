# Sprint plan — `PE` (queue item 17 `w-validated-proven-evidence-boundary`, **tranche 1 only**)

**Milestones**: `PE.A` … `PE.F` — six, each one PR, each independently CI-green.
**Status**: PLANNED · **SPLIT INTO SIX** · **TRANCHE 1 ONLY**
**Design doc**: [`w-validated-proven-evidence-boundary.md`](w-validated-proven-evidence-boundary.md)
(2259 lines, 13 quorum rounds, decision ledger **11 rows / 0 OPEN**)
**Base**: `dev` @ `2c9b5f3`, clean tree, **both gates measured green by the planner this session**
**Planner**: mission-control iteration 103, opus sprint-planner, first-party measurement on this rig
**Executor**: `codex:gpt-5.6-sol` under `--sandbox workspace-write`, sibling worktree, **NO git write permission**
**THE CONTROLLER MAKES ALL COMMITS.** Under that sandbox, in a *linked* worktree, `.git` is a **file**
pointing outside the sandbox, so `git commit` **cannot** succeed. That is a sandbox fact, not an
executor failure. The executor leaves the work **uncommitted**; the controller commits the diff.
Restores during mutation drills are `cp` from a `/tmp` backup taken **before** the mutation, verified
by sha256 — **never** `git checkout -- <file>`, which in a sprint worktree deletes the milestone's work.

**Cleared to plan by**: `D-WORLD-24` **arm A** (Mark Edmondson, attended, issue `#68`, bare comment
`A`, 2026-08-20T16:04:52Z) shed the bounded Z3 producer into new queue row **26**; round 13b applied
the narrow-refinement carve-out for both round-13 objections (→ AC21/AC22, M29/M30).
**Nothing in this plan re-litigates the design direction.**

**Headline price: 4.70 d across six PRs, none larger than 0.95 d.** That is the doc's own §9 number,
which the planner re-summed independently. It **exceeds the mission's 3–4 d sprint guardrail by
0.70 d** — a *stated, deliberate* overage, ratified scope, **not a new human ask**. §5 says how the
split absorbs it.

---

## 0. Planner's first-party verification

**Every base state below was re-derived BY COMMAND at `2c9b5f3` on this rig.** Nothing is transcribed
from the doc's §11 Verification Log — that log was measured at `bef0153`, `03c7892`, `52bc9ec`,
`7806cac`, `35fd875` and `516836f`: **five different, older epochs, by another author**. An empty
result is a *claim*, not a fact; **every negative below is paired with a known-positive control taken
in the same scope with the same instrument, and the control's observed value is recorded**.

`host/evidence` does not exist at HEAD. Confirmed: `ls -d host/evidence host/store` →
`ls: host/evidence: No such file or directory`, and the same call printed `host/store`.

### 0.1 The two things the planner *ran* rather than reasoned about

**(i) The AILANG change and its mutation were EXECUTED end to end against the pinned binary.**
This is the single riskiest AILANG assumption in the document — a seventh `gradeOf` arm under the
v0.30.0 Z3 encoding — and it was inherited, not measured, by every prior role.

On a scratch copy of `world/`, with `| ProofReceipt(HashRef)` appended, `ProofReceipt(_) => CLAIMED`
added to **both** the `ensures` postcondition and the body, and a seventh `gradeCode` tuple added:

| Run | Observed |
|---|---|
| control (unmodified) `ai-check` | `check.passed=true`, `errors=0`, `counterexample=0`, **6 verified** in `types.ail`, `gradeOf=verified` |
| **modified** `ai-check` | `check.passed=true`, `errors=0`, `counterexample=0`, **still 6 verified**, `gradeOf=verified` |
| **modified** `test --format json world/` | `len(tests[])` = **40**, `failed_tests=0`, seventh identity emitted as **`gradeCode_test_7`**, status `pass` |

→ **A seven-arm `gradeOf` still verifies. `EXACT_TOTAL_VERIFIED` stays 10** (`types.ail`'s own count
is unchanged at 6). `EXACT_TOTAL_TESTS` moves **39 → 40**. The identity to add to `REQUIRED_TESTS` is
**`gradeCode_test_7`** — observed, not assumed.

**M1, executed.** Mutating **only** the body arm to `PROVEN` (contract left at `CLAIMED`):

- **Leg 1** — `verify.counterexample = 1`, `gradeOf` flips `verified → counterexample`, `types.ail`'s
  verified count drops 6 → 5 (so `EXACT_TOTAL_VERIFIED=10` *also* breaks). Double kill.
- **Leg 2** — `failed_tests=1`, `gradeCode_test_7` **FAIL**, message **`test 0: expected 1, got 4`**.

> ⚠️ **The doc predicts the text `got 4, want 1`. The runner does not emit that.** An executor
> grepping for the doc's string would find nothing. See §1 defect **C4**.

**(ii) The go1.25.6 depth-bomb behaviour was re-run from scratch, four arms.**

| depth | bytes | result | elapsed |
|---:|---:|---|---:|
| 10 | 20 | **accept** | 0.23 ms |
| 9,999 | 19,998 | **accept** | 6.76 ms |
| 10,001 | 20,002 | refuse `invalid character '[' exceeded max depth` | 0.16 ms |
| **131,071** | **262,142** (inside the 262,144 cap) | refuse, same error | 0.15 ms |

`GOEXPERIMENT` empty; the refusing limit is the **unexported** `const maxNestingDepth = 10000` at
`GOROOT/src/encoding/json/scanner.go:148`, used at `:182`. The accepting arms are the discrimination
control. V50 reproduces exactly — **and the pin is not redundant**, because the limit really is an
internal nothing asserts.

### 0.2 Every acceptance criterion's base state, by command

| AC | Check (observed at `2c9b5f3`) | Firing control | Verdict |
|---|---|---|---|
| AC1/4/13/14/15/21 | `host/evidence` **absent** | same `ls` printed `host/store` | base = absent; every token new |
| AC3 | `ReadObject` **0** occ / **0** files; `ObjectMeta` **0** / **0** | **`GetObject` 80 occ / 23 files** — a *different* token, same scope, same instrument | CONFIRMED, **identical to V44's 23/80** eleven commits later |
| AC3 | `ObjectTooLargeError` **0** repo-wide | the same grep form returned 80 for `GetObject` | base = absent |
| AC16 | `ObjectReadTimeout` **0** | `context.WithTimeout` **25** all / **10** non-test | CONFIRMED; **doc's control figure is stale** → C3 |
| AC18 | `ErrUnorderedTimeouts` **0**; `func (s *Store) BusyTimeout(` **0** | `busyTimeoutMillis` **5 lines** (`writer_lock.go:173,179,196`; `context_read_test.go:209,210`) | CONFIRMED *exactly*, incl. V49's 5 and the cited line numbers |
| AC18 | `handlers.go:299–301` — *"an ORDERING nothing in this code asserts, not a guarantee"* | 13 populated lines printed | CONFIRMED at the cited offset |
| AC16 | `store.go:298` `db.SetMaxOpenConns(1)` | 3 hits, one the live statement | CONFIRMED |
| AC16 | DR-2 residue pinned at **8 / 2 / 1 = 11** (`approve.go`, `registry.go`, `replay.go`) | independent enumeration: **13** `.GetObject(` call sites and **4** interface declarations, *counted separately* | CONFIRMED — same per-file breakdown as V41/V43 |
| AC17 | exactly **9** statements naming `objects`: 3 `INSERT OR IGNORE` + 6 reads; **0** UPDATE/DELETE/DROP/ALTER; 0 triggers; 0 cascades | **positive enumeration**, because the negative control (`DELETE FROM` anything, non-test) reads **0** and *cannot fire* | CONFIRMED = V47's 9/3/6/0/0/0 |
| AC19 | `maxNestingDepth = 10000` at `scanner.go:148` + the 4-arm probe above | depth-10 and depth-9999 **accept** in the same call | CONFIRMED both directions |
| AC6/7/10 | `world/types.ail` and its projection **byte-identical**; gate: 4/4 projections, 4 exports, empty effects, **6 tar entries**, golden byte-for-byte | the gate's own nine steps all report non-zero work | CONFIRMED |
| AC10 | **no** Go test / script / workflow pins the live world/core hashes (`06acbb83`/`d16cc882`/`5823edcf`/`7856`) | the same `sha256:` grep form **does** fire at `host/pkgproj/pkgproj_test.go:25,29,43` — synthetic `t.TempDir()` fixtures | CONFIRMED; **removes a latent risk the doc never discusses** |
| AC8 | `scripts/verify_go.sh` has **zero** `REQUIRED_*`/`EXACT_*` pins | the same grep on `verify_ail.sh` returns lines 274, 323, 347, 363 | CONFIRMED absent → C6 |
| AC11 | AIL gate **PASSED** (10 identities / 39 tests / 9-of-9); Go gate **EXIT=0** | the race-detector known-positive control **fired** (`Found 2 data race(s)`) | **BASE IS GREEN**, measured outside any sandbox |
| — | queue row 22 `w-daemon-lock-wait-not-deadline-bound` carries **no** LANDED/COMPLETE marker | **2** other rows in the same file *do* match the marker pattern | **OPEN** — the AC16 residual has a live owner (`D-WORLD-23` obligation (i)) |

### 0.3 Measurements the controller did not make

- `host/pkgproj/pkgproj_test.go`'s three hardcoded `sha256:` literals were **read**, not assumed: they
  hash synthetic `t.TempDir()` fixtures (`"beta\n"`, `"alpha\n"`, `"manifest\n"`). PE.A's golden
  regeneration therefore reds **no** Go test.
- `TestNoNewDeadlineFreeStoreReads`' detector regex was read **in full** — see defect **C2**, the
  finding this planning pass exists to catch.
- `scripts/verify_go.sh` was read end to end: no manifest pattern to extend, and its final leg wraps
  the race run in a `python3` 600 s process-group killer that PE.F must not disturb.
- The socket-binding surface was enumerated, which makes the sandbox story *precise* rather than
  blanket: **only** `cmd/ailang-worldd`, `host/broker` and `host/daemon` bind loopback.
- `.github/workflows/ci.yml` was read: jobs `ailang-code verify gate` and `go host build + test gate`,
  both `GOTOOLCHAIN: go1.25.6`, both installing a pinned v0.30.0 binary (`:102`, `:144`). **§8.2's
  "confirm the Go job exports the pinned binary and Go toolchain during implementation" is
  discharged here — it does.**
- §9's tranche column re-summed independently: `0.65+0.75+0.50+0.75+0.35+0.45+0.50+0.30+0.15+0.10+0.20
  = **4.70** exactly`; decomposition `4.70+1.0+1.0+3.0+2.0 = **11.70** exactly`. **The arithmetic is
  sound as delivered and is not touched by this plan.**

---

## 1. Design-doc contradictions found

Wanted, not discouraged — and every one re-derived before it is asserted.

### C2 — **the ratchet the document leans on four times cannot see `ReadObject`** · severity HIGH

The doc rests AC16, §3.3 step 2, §8.2 and V41 on `TestNoNewDeadlineFreeStoreReads` as *the* thing
that makes the deadline-free residue mechanically observable (11 → 0). That test's detector is:

```go
var deadlineFreeReadCall = regexp.MustCompile(
	`\.(GetObject|GetWorld|GetLogEntry|GetRegistryHead|SelectedHead)\(\s*context\.Background\(\)`)
```

— **exactly five getters.** This tranche adds a **sixth** production store read method, `ReadObject`,
and the document never adds it to that set. A future `.ReadObject(context.Background())` anywhere
under `host/` or `cmd/` would land **silently**, past a ratchet whose entire stated purpose is to make
exactly that impossible — and it would do so in *the one method this tranche introduced specifically
to be deadline-bounded*.

**Cost to close: one token.** Add `|ReadObject` to the alternation, in the same commit that introduces
the method. **Measured count impact: zero** — `host/evidence` passes `readCtx`, never
`context.Background()`, and no other caller exists, so the pins stay 8/2/1 = 11.
**Folded into PE.B as task B5.** This is not scope creep: the doc already *claims* the ratchet covers
the residue, and without B5 that claim becomes false the moment `ReadObject` lands.

### C1 — AC18 still specifies a **live** `PRAGMA` three paragraphs before withdrawing it · MEDIUM

AC18's opening paragraph: *"read from the LIVE `PRAGMA busy_timeout` on the store's own connection"*.
AC18's own round-11b amendment, three paragraphs later: *"a NONBLOCKING CACHED property resolved once
during `Open` from the DSN … **not** a live `PRAGMA` query"*. §3.2 and §8.2 carry only the cached
spelling. The superseded sentence survived the amendment **in AC18 alone**, and an executor reading
top-down implements the wrong mechanism.

**Binding reading: CACHED.** Resolved once at `Open` from the DSN, thereafter a field read — no
connection acquisition, no query, no error channel. Recorded, not escalated: the amendment is later,
more specific, and reviewer-authored verbatim.

### C3 — AC16's `context.WithTimeout` control figure is stale · LOW

AC16 asserts *"the base `context.WithTimeout` control count of 8"*. That is the **non-test** count at
the doc's own base `35fd875` (verified: `git grep … 35fd875 … | grep -v _test.go | wc -l` → 8;
all-files → 22). At HEAD it is **10 non-test / 25 all**. The control still fires; the executor
**re-derives at its own base and never asserts 8**.

### C4 — M1's predicted failure text is not what the instrument emits · LOW

Predicted `got 4, want 1`; measured **`test 0: expected 1, got 4`**. The plan asserts M1's kill on
`failed_tests >= 1` **and** `gradeCode_test_7 != pass` **and** (independently) `verify.counterexample
>= 1` on `gradeOf` — never on a grep for the doc's prose.

### C5 — the `host/evidence → host/store` import edge is unstated but forcing · LOW

§3.2 says the package depends on `host/hashref` "and the minimal object reader seam", listing what it
does *not* depend on. It never states the `host/store` edge — yet `ObjectReader.ReadObject` returns
`store.ObjectMeta` and M24 requires `errors.As(err, &tooLarge)` against `*store.ObjectTooLargeError`.
No cycle exists (`host/store` imports only `host/hashref`), but this **forces the milestone order**:
PE.B must land before PE.C/PE.D can compile.

### C6 — there is no `verify_go.sh` manifest pattern to extend · INFORMATIONAL

§3.6/AC8 read as though the named-test leg extends an existing pattern. `scripts/verify_go.sh`
(160 lines) contains **zero** `REQUIRED_*`/`EXACT_*` pins. The templates are `scripts/verify_ail.sh`
Leg 2 (bounded run → `python3` JSON parse → required-name set → exact count) and
`host/verifygate/module_manifest_gate_test.go` (three anti-vacuity arms). PE.F writes the leg from
scratch against those. Priced inside the doc's existing 0.45 d row; **no re-pricing**.

---

## 2. Milestones

| ID | Name | Days | ACs closed | Mutations RED here |
|---|---|---:|---|---|
| **PE.A** | Kernel receipt arm, AILANG pins, projection & golden | **0.35** | AC6, AC7, AC10, AC2 (AILANG half) | **M1, M18, M19** |
| **PE.B** | Bounded one-snapshot store read + cached wait-bound accessor | **0.83** | AC3 (store half), AC17, AC18 (accessor half) | **M25** |
| **PE.C** | `host/evidence` codecs, byte caps, nesting-depth pin | **0.80** | AC2, AC19, AC3 (canonical-bytes half) | **M27** |
| **PE.D** | Validator, sealed mint authority, resolved grade, constructor refusals | **0.92** | AC1, AC4, AC13, AC14, AC15, AC21, AC12; constructor halves of AC16/AC18/AC22 | **M2, M3, M5, M7–M14, M16, M20, M21, M29** (15) |
| **PE.E** | Real-store integration proofs — the four kills no fake may make | **0.85** | AC16, AC18, AC22, AC3 (head reading) | **M4, M22, M23, M24, M26, M30** |
| **PE.F** | Persistent named-manifest gate, its self-mutation gate, full drill | **0.95** | AC8, AC9, AC11, AC12 | **M17** + re-drill of all 27 |
| | **Total** | **4.70** | 20 live ACs | **27 live rows** |

**Execution order is `PE.A → PE.B → PE.C → PE.D → PE.E → PE.F`**, and three of those edges are hard:

- PE.B **before** PE.C/PE.D — compile dependency (**C5**).
- PE.C **before** PE.D — the validator consumes the codecs.
- **PE.F last, without exception** — `EXACT_EVIDENCE_TESTS` pins the *observed* test count in
  `host/evidence`, so any test landed after it reds the gate.

PE.A is independent of every Go milestone and sits first only because it is fully measured, cheapest,
and gives the item an early green PR.

### PE.A — kernel receipt arm, AILANG pins, projection and golden · 0.35 d

Files: `world/types.ail`, `packages/world-core/world/types.ail`, `scripts/verify_ail.sh`,
`scripts/world_package_ready_packet.golden.json`. ~12 net source lines + a regenerated golden.

**The five coupled AILANG moves land in ONE commit** (§8.1): a projection that lags its source reds
step 3/9; a stale golden reds step 9/9.

Requirements, all pre-verified by the planner: `ai-check` gives `check.passed=true`, `errors=0`,
`counterexample=0`, `gradeOf=verified`, and `types.ail`'s verified count **stays 6** so repo-wide
`EXACT_TOTAL_VERIFIED` **stays 10**; the test leg gives `len(tests[])=40` with the seventh identity
emitted as `gradeCode_test_7` (**use the observed name**); the two `types.ail` files stay byte-identical;
`verify_ail.sh` passes all nine world-package steps; the Go gate stays `EXIT=0` (measured
pre-condition: nothing pins the live world/core hashes). Do **not** touch `REQUIRED_VERIFIED` or
`EXACT_TOTAL_VERIFIED`; do **not** invent `EXACT_TOTAL_MODULES`.

Mutations — **M1** measured red on both legs (see §0.1); **M18** edit only `world/types.ail` → step
3/9 projection hash mismatch (*confirm the exact wording during implementation*); **M19** rebuild the
projection but retain the old golden → step 9/9 byte-for-byte failure.

**Sandbox: fully informative** (no Go tests, no sockets). Only the Go-gate re-run needs the controller.

### PE.B — bounded one-snapshot store read + cached wait-bound accessor · 0.83 d

`host/store` gains exactly **two** additive surfaces, neither changing an existing signature.

**(1) `ReadObject(ctx, ref, maxBytes) (ObjectMeta, []byte, error)`** — the `D-WORLD-21` arm-A ruled
spelling. Two ordered statements on **one reserved connection inside one read transaction**: a probe
`SELECT interface_hash_ref, semantic_id, provenance, length(payload) FROM objects WHERE hash_ref = ?`
**whose select list omits the payload column**, so no payload byte crosses the driver before the
guard; then a typed `*ObjectTooLargeError` carrying the probed size and **no payload**; and only under
the guard, `SELECT payload …`. **Both** statements through the transaction's `QueryRowContext` with
the **supplied** context, and the connection reservation **and** transaction begin taking that same
context. The probe is the transaction's first statement and fixes the snapshot. Plus a package-private
scheduling hook, **nil in production**, that *wraps* the interval between the two statements and
replaces nothing.

**(2) `BusyTimeout() time.Duration`** — a **nonblocking cached property** resolved once during `Open`
from the DSN, then a field read (**C1**). It must **not** re-export `busyTimeoutMillis`: a
caller-supplied DSN that set its own window (which `withBusyTimeout` deliberately never overrides,
`writer_lock.go:181–198`) must be reported at *its* value. Test both arms. **"Nonblocking" is a
measurement, not an adjective**: with the sole pooled connection occupied by a decoy, `BusyTimeout()`
must return within a bound far below the hold.

**B5 — the C2 fix.** Add `|ReadObject` to `deadlineFreeReadCall`. Zero count impact, measured.

**AC17 / M25.** `TestConcurrentMutationCannotDesyncProbeAndPayload` against a file-backed **real**
store: the hook fires between probe and payload and attempts `UPDATE objects SET payload = …` through
a **second** `sql.DB` handle on the same file. Assertion: the returned payload's length equals the
probed `ObjectMeta` length **and** its recomputed hash equals the requested ref. **Either** writer
outcome is accepted — committed-after-snapshot under WAL, busy-refused under a rollback journal —
because the journal mode is deliberately unmeasured. **The hook must be asserted to have FIRED and the
writer's outcome recorded**: a green can never come from a writer that never ran. Two-armed control,
both required. The mutating SQL is test-authored **by necessity** — the planner re-verified that
production has 9 statements on `objects`, none of them a mutation.

> **R1, the specific risk here.** `db.SetMaxOpenConns(1)` is measured at `store.go:298`. Reserve
> **once** via `db.Conn(ctx)` and begin the transaction **on that connection**. A second `db.Conn` or
> a `db.BeginTx` alongside a held conn will hang.

**Sandbox: informative** for `go test ./host/store` (binds no sockets). B4's occupied-connection arm
and B3's interleaving are load-sensitive → controller re-runs.

### PE.C — `host/evidence` codecs, byte caps, nesting-depth pin · 0.80 d

Strict canonical codecs for `ProofReportV1` (nine fields, **in order**) and the envelope (exactly two
fields in order: `report`, base64url-no-padding of the canonical report bytes; `mac`,
base64url-no-padding 32-byte tag). Unknown / duplicate / missing / non-canonical / trailing /
invalid-UTF-8 / over-limit input is **rejected**, with **one** classification exception: **the
envelope's `mac` member reaches the authentication step**, so every tag failure carries the single
stable reason `unauthenticated_report`. Decode-then-re-encode is **byte-equal** for each.

Caps: raw envelope 256 KiB **owned by `ReadObject`'s `maxBytes` and not re-checked here** (a duplicate
check would make **M4 observable-identical and therefore unkillable**); decoded report 256 KiB;
verified list 256 unique identities; every string 1 KiB. `DecodeProposal` caps its raw input at
256 KiB **before any parse** — assert on an input that would *also* fail parse, so the refusal cannot
be attributed to the parser.

**AC19 / M27.** `TestNestingDepthBombWithinByteCapIsRefused` — depth-131,071 / 262,142 bytes must
produce the typed decode refusal and **no** `ClaimedEvidence`; a depth-10 control **in the same test**
must decode. No depth guard of our own is added; the pin exists because the refusing limit is an
unexported stdlib internal (§0.1(ii)).

**Sandbox: fully informative.**

### PE.D — validator, sealed mint authority, resolved grade, constructor refusals · 0.92 d

The largest milestone: ~1150 LOC, **15 mutations**, and the whole authority argument.

`ObjectReader` is the **mandatory two-method seam** — `ReadObject` **and** `BusyTimeout()`. There is
**no** optional `busyWindowReporter`, no type assertion, no skip path: round 13b deleted the
type-assert-and-skip whose *absence of capability* was silently read as *absence of lock-wait*.

`NewValidator(key [32]byte, reader, compilerConfig, requiredIdentities []string)` with **three
distinct** constructor refusals — non-positive `ObjectReadTimeout` **and** empty/unset identities
**and** negative/unknown `BusyTimeout()` all reuse `ErrInvalidValidatorConfig`; a *positive*
`ObjectReadTimeout` ≤ a *positive* reported window is the dedicated `ErrUnorderedTimeouts`, naming
both values. **The ordering is pinned against the values actually used at run time, never two
literals** — comparing two constants is a tautology that stays green while the configured timeout
changes. AC16's non-positive guard is **weaker** and must stay green on a positive-but-unordered
stimulus; that is what makes the arms distinguishable rather than one check wearing two names.

**The mint-identity mechanism, exactly as specified.** An unexported pointer to a per-instance
**non-zero-size** heap allocation made only inside `NewValidator` — non-zero-size because Go permits
two distinct zero-size allocations to share an address, which would let two identities **collide by
accident**. Pointing it at the validator's own heap-held key copy satisfies this. **Assert the
distinctness in a test; do not merely comment it.**

**`Resolve` performs two ORDERED checks, separately observable.** (1) **Mint validity** — nil receiver
identity *or* nil `mintedBy` → `ErrUnmintedAuthority`, **before any comparison**. (2) **Binding** —
mismatch → `ErrForeignSeal`. Each sentinel is produced by exactly one check and nowhere else, and
neither is a member of §2.5's `UnsupportedReason` set. **If check 1 were merged into check 2, the
zero-zero pair would compare EQUAL and the forge would resolve** — which is precisely M21.

`ResolutionResult` has unexported fields, mutually exclusive `Proven() (ResolvedGrade, bool)` and
`Err() error`, and a **zero value that is an explicit refusal**. There is **no** `NewValidatedEvidence`,
struct literal, `SetGrade`, raw-hash resolver, receipt resolver, exported unseal, or package-level
grade resolver of any name: the free `GradeOfValidated` is **dropped, not renamed**, with no
deprecated alias. External-package tests live in `host/evidence/authority_test.go` so they compile
with **only** the public API.

**Copy semantics are asserted, not avoided.** `v2 := *v1` carries the same identity and key and **is**
the same mint authority. The guarantee is **per-identity, not per-Go-variable**. Self-minting into a
caller-constructed validator remains possible and is **accepted** — no library can stop a caller lying
to itself.

> **The fake-audit rule (§6.1) binds every test here.** A prescribed fake is *admissible* where it
> supplies **input** the mutated mechanism consumes; *inadmissible* where it supplies the **property**
> the mutation is supposed to expose. Round 8's reject was exactly that defect. The retired round-7b
> `TestBlockingObjectReadReturnsWithinObjectReadTimeout` with its context-observing fake is
> **deleted, not retained beside the honest test — do not write it.**

Shared control for M2–M14 (and M24): the test **first** validates one good authenticated report and
observes the minting validator's `Resolve` return `ResolvedGradeProven`, so a mutant cannot pass
merely because the test never reached minting.

**Sandbox: fully informative** (fakes only; no database, no sockets).

### PE.E — real-store integration proofs · 0.85 d

Test-only, ~700 lines. **No fake participates in any kill in this milestone, by construction.**

**AC16 — `TestRealStoreBlockedObjectReadReturnsWithinObjectReadTimeout`.** A decoy goroutine holds the
store's **single** pooled connection with one in-flight `GetObject` of a decoy object sized so a single
read outlasts 20× `ObjectReadTimeout` — **the hold duration is asserted in-test, not assumed.**
`ValidateProof` then runs under `context.Background()` against a file-backed store; `ReadObject` waits
for the pooled connection, and `database/sql`'s connection wait honours the supplied context **by
stdlib contract**, so the attempt must return the operational error within the bound, carrying no seal.

The classification is **race-honest and non-negotiable**: a seal below the bound means the attempt won
the freed connection → retry, **bounded at 5, with exhaustion a loud instrument failure and never a
pass**; a seal at ≥ 2× the bound is the **mutant signature** and reds immediately; a 20× watchdog makes
the hang mode a red rather than a hang. The **no-decoy control runs first, in the same test**.

**The blocking mechanism is chosen deliberately.** Lock contention is **rejected** as the stimulus —
iteration 94 measured a lock-blocked read returning at `busy_timeout` (2.043 s under a 300 ms
deadline), bounded by *the wrong constant*, which is the queue-row-22 composition this tranche must not
absorb. Interrupting the target read's own blob transfer is **also rejected** — iteration 94's
in-flight interruption measurement was a many-opcode query, and asserting the same of one
blob-column materialization would generalize past the evidence. The pool wait depends only on the
stdlib contract plus the measured single-connection pool.

**AC3 head reading / M4 / M24 — `TestOversizeProofReportIsRefused`** through the **real** store.
Base64's 4/3 inflation makes the fixture constructible: a ~250 KiB correctly-MAC'd report yields a
~333 KiB envelope, over the read bound while the decoded report stays inside its own 256 KiB cap.
**M4 and M24 are paired by design** — M4 proves the *store guard* is load-bearing, M24 proves the
*validator's reason mapping* is; the byte-slice seam separates what the round-7b stream unified, so
§5's every-refusal-branch rule demands both.

**AC18 — `TestConstructorPinsBusyTimeoutBelowObjectReadTimeout`**, real store. Ordered arm first (3 s
> 2000 ms → constructs, validates, resolves `ResolvedGradeProven`); pin arm (a **positive** 1 s ≤
2000 ms → `ErrUnorderedTimeouts`). Plus the **accessor cross-check moved from the construction path to
the test** (round-11b): the test issues a direct `PRAGMA busy_timeout` itself, on a store whose
connection is free, and requires equality with `BusyTimeout()` — so the cached value cannot drift into
a stale literal while **no production path performs the query**. Plus the **required nonblocking arm**
(§E8), *"what makes nonblocking a measurement instead of an adjective."*

**AC22 — `TestReaderWaitBoundsCannotBeLostThroughWrapper`**, wrapping the real store. Unknown/negative
arm cannot construct; forwarding arm reaches `ErrUnorderedTimeouts`. The wrapper is the **stimulus,
not a substitute** for the store.

Falsifiability readings (AC3, AC16, AC18) are **re-derived at the executor's own base**, never
transcribed — including AC16's control, where the doc's `8` is stale (**C3**).

> ⚠️ **This is the milestone that most needs an out-of-sandbox re-run — and it is not a socket
> problem.** `host/evidence` and `host/store` bind nothing. But E1/E2/E8 are **wall-clock classified**,
> and a sandboxed, load-contended executor can produce a false red **indistinguishable from the
> M22/M23 mutant signature**. Every timing verdict here is labelled `UNINFORMATIVE UNDER SANDBOX` and
> re-run by the controller. This mission already burned an iteration on exactly this (iter-97: 64
> unkilled spinners produced a false *"dev is red"* in two independent roles).

### PE.F — persistent named-manifest gate, self-mutation gate, full drill · 0.95 d

A focused leg in `scripts/verify_go.sh`, **before** the broad plain/race runs: `go test -json
./host/evidence -count=1`, parsing **only terminal `Action=pass` events** for `Test…` identities,
requiring the exact `REQUIRED_EVIDENCE_TESTS` set, requiring **that set and the observed set are both
non-empty**, and pinning `EXACT_EVIDENCE_TESTS` to the observed count. It fails on missing, skipped,
failed, **duplicate**, or **extra** tests. Do not disturb the existing `python3` 600 s process-group
killer wrapping the race leg.

**No precedent to extend (C6)** — write it against `verify_ail.sh` Leg 2's shape and
`module_manifest_gate_test.go`'s three anti-vacuity arms. **The anti-vacuity floor is explicit
(§5): one discovered package, a non-empty required set, and at least one terminal named-test pass.
A shell grep over source is not acceptance evidence.**

`TestEvidenceNamedManifestRejectsUnpinnedTest` in `host/verifygate` drives the leg **in isolation**
over a synthetic observed set, and must also fail loudly on an empty required set and an empty
enumeration.

**AC9 — the full 27-row drill.** Order: AILANG rows first (cheapest to restore), then the pure-Go
validator rows, then the real-store rows last (slowest, most load-sensitive). Record for **every** row:
the exact edit, the exact command, the observed **RED** text, and the observed **GREEN** of its named
control. Restore by `cp` from a pre-mutation `/tmp` backup verified by sha256.

**AC12** — `git diff --stat` over the whole item shows **zero** changes under `host/daemon/`, `cmd/`,
`host/replay/` and any renderer surface, with a firing control (the diff is non-empty elsewhere).

**Sandbox: mixed.** The leg and the isolated gate are informative; the full-gate runs are
`UNINFORMATIVE UNDER SANDBOX` because `verify_go.sh` does `go test ./...`, sweeping `host/daemon`,
`host/broker` and `cmd/ailang-worldd`.

---

## 3. Mutation → milestone map

**Live set: M1–M5, M7–M14, M16–M27, M29, M30 = 27 rows.**

**Mapping rule: a row belongs to the milestone containing its NAMED KILLER TEST, not the milestone
containing the mutated line.** M4 mutates `host/store` (PE.B's code) but is killed by
`TestOversizeProofReportIsRefused` in `host/evidence`, so it discharges in PE.E. M26/M30 neuter
constructor code written in PE.D but are killed by real-store tests written in PE.E. This is the only
mapping under which every row is actually RED-provable in the PR that claims it.

| Milestone | Rows | Count |
|---|---|---:|
| PE.A | M1, M18, M19 | 3 |
| PE.B | M25 | 1 |
| PE.C | M27 | 1 |
| PE.D | M2, M3, M5, M7, M8, M9, M10, M11, M12, M13, M14, M16, M20, M21, M29 | 15 |
| PE.E | M4, M22, M23, M24, M26, M30 | 6 |
| PE.F | M17 (+ re-drill of all 27) | 1 |
| | **Total** | **27** |

**Cannot be killed without the REAL store** (§6.1): **M4** (the mutated line *is* store code), **M22**
and **M23** (no fake participates by construction — round 8's reject was precisely a fake that observed
the context, making the mutant die for a property the real reader had never been shown to have),
**M25** (the hook supplies a scheduling *input*; snapshot isolation is supplied by the real SQLite
transaction), **M26** (only the real store reports a live busy window), **M30** (the wrapper
*surrounds* the real store).

**M22's "EVERY context-taking call" is load-bearing.** Round 9b's one-snapshot spelling moved the
connection **wait** — the very wait AC16's stimulus blocks on — out of `QueryRowContext` and into the
**reservation**. A mutant detaching only the two statement sites would leave that wait context-bounded
and the row would survive **vacuously**.

### Deliberately excluded — do not plan, do not renumber

- **M6** — tranche 3 (replay). `host/replay` untouched. Until tranche 3 lands, replay has **no** route
  to `PROVEN`, which is the stronger default.
- **M15, M28** — **shed with the producer into queue row 26** by `D-WORLD-24` arm A. **Their numbers
  are deliberately absent and MUST NOT be renumbered or reused.**
- **AC5, AC20** — **declared gaps, not oversights**; moved to row 26 with their original identifiers.
  AC20 carries the connection-pool-versus-lock-wait vacuity obligation with it.

---

## 4. What this sprint explicitly does **not** do

- **Tranche 2** `w-proven-evidence-production-key-wiring` (1.0 d) — **no** production MAC key
  provisioned or loaded, **no** production composition root wired, **no** startup/restart behaviour
  claimed. `NewValidator` takes an exactly-32-byte key **directly from its caller** and accepts no key
  path. **Tranche 1 is explicitly library-only, NON-PRODUCTION, and not deployable.**
- **Tranche 3** `w-validated-replay-evidence-boundary` (3.0 d) — no replay evidence; **M6 is specified
  in the doc but deliberately not planned, implemented, or discharged here**; `RecordedEffect` stays
  `ATTESTED`.
- **Tranche 4** `w-proven-evidence-renderer-consumption` (2.0 d) — no renderer, daemon route, CLI verb,
  or agent tool. AC12: **no surface may display `PROVEN`.**
- **Queue row 26** `w-bounded-z3-report-producer` — the shed producer: `proof_producer.go`,
  `NewProducer`, `Producer.GenerateProof`, `ObjectWriter`, the additive store `WriteObject`, bounded
  process-group execution, pinned-executable-byte checks. **Also travelling to row 26 and deliberately
  not answered here: both unresolved round-12 producer objections** (the missing-reader constructor
  gap, and AC20's vacuity obligation).
- **Item 18's 11 deadline-free production store reads** (approve 8, registry 2, replay 1) — item 18's
  follow-on obligation, not absorbed. This tranche imposes its own deadline at its **own** call site
  and inherits none.
- **Queue row 22** `w-daemon-lock-wait-not-deadline-bound` (verified **OPEN** by command this session)
  owns the residual: under real lock contention the store's wait is bounded by `busy_timeout`, not by
  `ObjectReadTimeout`. This tranche **pins the ordering** that keeps the composition safe (AC18) and
  does **not** fix the composition.
- **`GetObject` is not modified**; the 13 call sites and 4 interface declarations are untouched. No
  schema change, migration, or registry head. **`tools/launchd/*` is frozen core** and is not touched
  under any circumstances.
- **No Z3 contract is attached to `Proposal`** (the measured verifier limitation).
- **§9's arithmetic is not re-opened.** Both columns re-summed exactly. The 0.70 d overage is a stated
  pricing fact, not a new ask.

---

## 5. Pricing, velocity, and the guardrail

**Velocity is measured, not guessed** — the last six landed milestone squashes:

| Commit | Date | Size |
|---|---|---|
| `7ad24ea` item 18 M1 | 2026-08-18 | 30 files, 1747 (+), 125 (−) |
| `b3c5de0` item 18 M2 | 2026-08-18 | 6 files, 988 (+), 62 (−) |
| `d21754f` item 18 M3 | 2026-08-18 | 5 files, 338 (+), 10 (−) |
| `6c2a537` item 19 | 2026-08-19 | 5 files, 2553 (+), 2 (−) |
| `9fa2647` row 21 | 2026-08-20 | 7 files, 1177 (+), 21 (−) |
| `912009d` row 20 | 2026-08-20 | 5 files, 769 (+), 40 (−) |

**Observed band 338–2553 insertions per landed milestone PR, median ≈1080, roughly one milestone PR
per mission iteration.** This plan's six PRs total **3642** net new source lines — average **607**,
max **1150** (PE.D). **Nothing here is out of family.**

**Days reconcile to the doc exactly.** `0.35 + 0.83 + 0.80 + 0.92 + 0.85 + 0.95 = 4.70`. Every §9
tranche row is assigned in full to exactly one milestone, except the 0.30 d timeout-ordering row,
split **0.08** (store accessor → PE.B) / **0.07** (constructor refusal → PE.D) / **0.15** (real-store
AC18 test + M26 → PE.E). No row dropped, duplicated, or re-priced.

**The guardrail.** The mission's sprint guardrail is 3–4 d; **4.70 d exceeds it by 0.70 d**. §9 states
that deliberately (*"round 9's own standard refuses to disguise the number by leaving 4.0 because it
was 4.0"*). **This plan does not re-price it.**

**How the split absorbs it — by sequencing, the only lever a planner legitimately has.** Six
independently landable PRs mean the guardrail never binds a single unmergeable stretch: the largest
continuous commitment is **PE.D at 0.92 d**. The mission has the precedent — queue item 11
`w-transition-registry` shipped as TR.A / TR.B / TR.C across three plans and three iterations.

**If the iteration budget forces a cut**, the honest seam is **after PE.D**: `PE.A+PE.B+PE.C+PE.D =
2.90 d` delivers the kernel arm, the bounded store read, the codecs and the full authority boundary
with **20 of 27** mutations discharged; `PE.E+PE.F = 1.80 d` completes the real-store proofs and the
persistent gate. **That cut is NOT recommended** — it leaves the named manifest unlanded and AC8/AC9
open — but it is the seam if one is needed.

---

## 6. Gates, sandbox, worktree

```bash
# AILANG gate — the shipped, pinned binary IS the gate. Profile is ailang-code, not go-compiler.
export AILANG_BIN=/tmp/ailang-v0300/ailang
./scripts/verify_ail.sh

# Go host gate — THE EXPORT IS MANDATORY. Without it the gate FAILS CLOSED BY DESIGN:
# host/replay's pinnedBinary(t) t.Skip()s silently and a bare `go test` reports `ok`
# with the load-bearing assertions never running.
export AILANG_BIN=/tmp/ailang-v0300/ailang
export GOTOOLCHAIN=go1.25.6
./scripts/verify_go.sh

# The subset that IS informative inside --sandbox workspace-write (none of these bind sockets):
GOTOOLCHAIN=go1.25.6 go test ./host/evidence ./host/store ./host/verifygate -count=1
```

Base results, measured this session at `2c9b5f3`: AIL gate **PASSED** (10 identities / 39 tests / 9-of-9
world-package steps); Go gate **EXIT=0** with the race-detector known-positive control firing. Binary:
`AILANG v0.30.0`, commit `e37b370` — never a `-dirty` dev build, never from memory. CI: workflow `CI`,
jobs `ailang-code verify gate` and `go host build + test gate`.

**Sandbox — three binding facts.** (1) **No commits**: `.git` is a file pointing outside the sandbox in
a linked worktree. (2) **No loopback binds** — measured surface: `net.Listen`/`httptest` appear **only**
in `cmd/ailang-worldd/cli_test.go`, `host/broker/registry_publish_test.go`,
`host/broker/registry_reconcile_test.go` and `host/daemon/*`. So `host/evidence`, `host/store`,
`host/verifygate` and `world/` are **informative**, while **any** `./scripts/verify_go.sh` run is
**`UNINFORMATIVE UNDER SANDBOX`** because it does `go test ./...`. (3) **PE.E is the special case and
it is not a socket problem** — it is wall-clock classification under load. Every such verdict is
labelled and re-run by the controller outside the sandbox; **that re-run is the verdict.**

**Worktree**: a **sibling of the repo**, e.g. `/Users/voightkampff/dev/sunholo-data/.wt-iter103`.
**Never under `/tmp`** — a `/tmp`-rooted checkout reds tests for its *location* rather than its code,
and a red CI can never reproduce it. This item is especially exposed: PE.E opens file-backed SQLite
stores whose journal-mode behaviour under `/tmp` is not the behaviour CI sees.

---

## 7. Open questions for the human

**None.** Deliberately empty. No human is present for this run, and *a question a measurement can
answer is not a question*. Every candidate was measured instead: whether a seven-arm `gradeOf` still
verifies (**ran it — yes**); whether `EXACT_TOTAL_VERIFIED` moves (**ran it — no, stays 10**); the
emitted seventh identity (**ran it — `gradeCode_test_7`**); whether regenerating the golden reds any Go
test (**grepped it with a firing control — no pin exists**); whether the stdlib depth limit still holds
at the pinned toolchain (**ran it — yes, `scanner.go:148`**); whether queue row 22 is still OPEN
(**grepped it with a firing control — yes**). The 0.70 d overage is a **stated pricing fact** from §9
under ratified scope, **not a new ask** — manufacturing a decision here would cost a whole iteration
for nothing.

---

**SPRINT_PLAN_PATH**: `design_docs/planned/w-validated-proven-evidence-boundary-sprint-plan.md`
**SPRINT_JSON_PATH**: `.ailang/state/sprints/w-validated-proven-evidence-boundary.plan.json`
