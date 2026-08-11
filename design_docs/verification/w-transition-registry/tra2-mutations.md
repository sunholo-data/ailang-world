# TR.A2 mutation transcript (iteration 72)

All Go arms were applied one at a time after copying the target to
`/tmp/tra2_backup/<file>.bak`. Each mutant had a changed branch-specific anchor and SHA-256,
compiled with `GOTOOLCHAIN=go1.25.6 go build ./...`, ran the named kill, ran the package inverse,
and was restored with `cp`. The restored production SHAs were:

- `host/transitionreg/transitionreg.go`: `5aecd46b3bb9080dc4fb5fe2dadfdd14458ad17444a90f779b4db484b8504611`
- `host/transitionreg/codec.go`: `11ad144ec17e6291cecdf2540f8c74031cd9b7c48353705116867d168aa5ba49`

Anchor counts below are `before -> mutant`. “Inverse” is the package run with the named test
skipped. Two arms have non-zero inverses and are explicitly not claimed as isolated kills.

| mutation / arm | anchor | SHA before -> mutant | build | kill and observed failing test | inverse | restore | result |
|---|---:|---|---:|---|---:|---|---|
| MUT-GO-CODEC-TAG | v1 semantic literal `1->0` | `11ad144e... -> f2c3cb3...` | 0 | rc=1 `TestCodecGoldenRoundTrip` (literal golden byte assertion; interface-hash literal assertion did not fire) | **1** | `11ad144e...` match | KILLED, **not isolated** |
| MUT-READ-EMPTY-OK | absent-head error `1->0` | `5aecd46b... -> 3ce3cb1d...` | 0 | rc=1 `TestReadSnapshotRefusals/absent_head` | 0 | `5aecd46b...` match | KILLED |
| MUT-READ-SWALLOW | injected head wrapper `1->0` | `5aecd46b... -> 301ca282...` | 0 | rc=1 `TestReadSnapshotRefusals/injected_read_error` | 0 | match | KILLED |
| MUT-READ-ABSENT-OK | object-absent error `1->0` | `5aecd46b... -> 5b2a9fab...` | 0 | rc=1 `TestReadSnapshotRefusals/object_absent` | 0 | match | KILLED |
| MUT-READ-NO-REHASH | live hash condition `1->0` | `5aecd46b... -> 58a09a2f...` | 0 | rc=1 `TestReadSnapshotRefusals/corrupt_object_payload` | 0 | match | KILLED |
| MUT-READ-ANY-TYPE (semantic) | live semantic guard `1->0` | `5aecd46b... -> ffda8806...` | 0 | rc=1 `TestReadSnapshotRefusals/wrong_semantic_id` | 0 | match | KILLED |
| MUT-READ-ANY-TYPE (interface) | live interface guard `1->0` | `5aecd46b... -> 364464a7...` | 0 | rc=1 `TestReadSnapshotRefusals/wrong_interface_hash` | 0 | match | KILLED |
| MUT-SNAPSHOT-ALIAS | three clone assignments `3->0` | `5aecd46b... -> 7836721c...` | 0 | rc=1 `TestSnapshotIsEagerAndCopyIsolated` (construction/List/Lookup coverage) | 0 | match | KILLED |
| MUT-SNAPSHOT-REREAD | head calls `1->2` | `5aecd46b... -> 41c05b24...` | 0 | rc=1 `TestReadSnapshotReadsHeadOnce` | 0 | match | KILLED |
| MUT-SNAPSHOT-CACHE-BYPASS (head) | pre-head cache branch `0->1` | `5aecd46b... -> 8ef4f78c...` | 0 | rc=1 `TestReadSnapshotReadsHeadOnce` | 0 | match | KILLED |
| MUT-SNAPSHOT-CACHE-BYPASS (copy) | cached clone return `1->0` | `5aecd46b... -> dd1446f5...` | 0 | rc=1 `TestSnapshotIsEagerAndCopyIsolated` | 0 | match | KILLED |
| MUT-REVISION-SKIP (number) | live revision guard `1->0` | `5aecd46b... -> 5fc6b37f...` | 0 | rc=1 `TestPublishRefusals/revision_not_n_plus_1` | 0 | match | KILLED |
| MUT-REVISION-SKIP (parent) | live parent guard `1->0` | `5aecd46b... -> 9b238bbe...` | 0 | rc=1 `TestPublishRefusals/parent_not_captured_head` | 0 | match | KILLED |
| MUT-ORDER-INSERTION (sort) | disabled-sort wrapper `0->1` | `5aecd46b... -> 18448969...` | 0 | rc=1 `TestStableIDByteOrder` | 0 | match | KILLED |
| MUT-ORDER-INSERTION (duplicate) | `>= 0` guard `1->0` | `5aecd46b... -> 1a48903e...` | 0 | rc=1 `TestPublishRefusals/duplicate_id` | 0 | match | KILLED |
| MUT-PUBLISH-SWALLOW (PutObject) | put error wrapper `1->0` | `5aecd46b... -> 03fb1228...` | 0 | rc=1 `TestPublishRefusals/injected_put_error` | 0 | match | KILLED |
| MUT-PUBLISH-SWALLOW (CAS) | CAS error wrapper `1->0` | `5aecd46b... -> 71fde8f1...` | 0 | rc=1 `TestPublishRefusals/injected_cas_error` | **1** (`TestPublishCASConflictPreservesWinner` also fails) | match | KILLED, **not isolated** |
| MUT-ENTRIES-NO-LIMIT | live entries guard `1->0` | `5aecd46b... -> 1244ba67...` | 0 | rc=1 `TestPublishRefusals/entries_over_1024` | 0 | match | KILLED |
| MUT-REVISION-NO-LIMIT | max-raw call `1->0` | `11ad144e... -> fea7b1df...` | 0 | rc=1 `TestReadSnapshotRefusals/revision_raw_over_limit` | 0 | `11ad144e...` match | KILLED |
| MUT-AIL-EMPTY-MODULE | `world/*.ail` `4->5` | absent -> `cdbc91b7...` | 0 | repaired AC9 rc=1; `modules11=0`, `tests14=1`, `steps9=1` | raw script rc=0 | absent; file count 4 | KILLED |
| MUT-DELETE-TR-A-TEST (AC2) | AC2 names `3->2` | `33bd76d3... -> c611ce74...` | 0 | activated AC2 rc=1, observed count=2 | 0 | `33bd76d3...` match; restored gate rc=0/count=3 | KILLED |
| MUT-DELETE-TR-A-TEST (AC3) | AC3 names `4->3` | `33bd76d3... -> 4946d4b4...` | 0 | activated AC3 rc=1, observed count=3 | 0 | `33bd76d3...` match; restored gate rc=0/count=4 | KILLED |

The semantic-tag arm exposes a plan-level coupling: changing the package semantic constant affects
many tests, so the required package inverse cannot be green. The CAS-swallow arm likewise changes
real CAS-conflict behavior covered by `TestPublishCASConflictPreservesWinner`; skipping only the
injected-error subtest cannot make the package green. Both mutants landed and built, and both named
tests fired, but neither is represented as an isolated mutation kill.

## Controller adjudication of the two non-isolated arms (iteration 72, rule 3h)

The executor flagged both departures rather than quietly weakening a test, so both were adjudicated
by measurement in **both** directions rather than accepted on the strength of its report.

**The checkable proposition:** *the package inverse is red because ADDITIONAL legitimate tests
detect the mutation — over-coverage — and not because the named test's failure is a bystander
effect of unrelated breakage.* The discriminating command is the inverse arm widened to skip the
COMPLETE detector set. If the mechanism is over-coverage that run is rc=0; if the mutant is
breaking something unrelated it stays red. Both arms were then re-run with the same `-skip` on the
**unmutated** tree, so an rc=0 cannot be an artifact of the skip set concealing a pre-existing red
(rule 3d's negative control, aimed at a mutation inverse).

| arm | anchor | build | named kill | executor inverse | full detector set | **skip-all inverse** | same skip, UNMUTATED | verdict |
|---|---|---:|---|---:|---|---:|---:|---|
| `MUT-PUBLISH-SWALLOW` (CAS) | `1 -> 0`, sha `5aecd46b…` -> `5de5c683…` | 0 | rc=1 `TestPublishRefusals/injected_cas_error` | rc=1 | `TestPublishRefusals`, `TestPublishCASConflictPreservesWinner`, **`TestConcurrentPublishHasOneWinner`** | **rc=0** | rc=0 | **KILL, over-coverage** |
| `MUT-GO-CODEC-TAG` | `1 -> 0`, sha `11ad144e…` -> `f2c3cb37…` | 0 | rc=1 `TestCodecGoldenRoundTrip` | rc=1 | `TestCodecGoldenRoundTrip`, `TestReadSnapshotRefusals` | **rc=0** | rc=0 | **KILL, over-coverage** |

Both production files were restored from a `cp` backup and re-hashed: `codec.go`
`11ad144ec17e6291cecdf2540f8c74031cd9b7c48353705116867d168aa5ba49` and `transitionreg.go`
`5aecd46b3bb9080dc4fb5fe2dadfdd14458ad17444a90f779b4db484b8504611` — byte-identical to the
executor's delivered tree, so this adjudication changed nothing.

**Two corrections to the executor's characterisation, both measured, neither changing the verdict:**

1. The CAS arm names **one** co-detector. There are **two** — `TestConcurrentPublishHasOneWinner`
   also fails, because a swallowed CAS error makes every racer believe it won. That is a
   *strengthening*: the concurrency invariant is pinned by a second, independent test.
2. The semantic-tag arm says the constant "affects many tests". Enumerated with `-v`, it is exactly
   **two**, and the second one is worth naming because it is a consequence of this milestone's own
   design decision: `TestReadSnapshotRefusals/wrong_semantic_id` pins the *measured message*
   `semantic ID "wrong/semantic" is not "world/transition-registry/v1"`, which **embeds the
   constant**. Per-branch message pinning — mandated here after TR.A1's two survivals — therefore
   couples the refusal table to the identity literal. The coupling is benign and arguably desirable
   (a silent semantic-ID change now reds two tests, not one), but it is the reason this arm cannot
   have a green single-test inverse, and it would be invisible to anyone reading only the
   mutation's name.

## Evaluator finding — THREE refusal branches shipped with zero coverage (rule 3j), now closed

The `sonnet` evaluator (PASS 86/100) filed two BLOCKING findings that neither the executor nor the
controller caught. Both were reproduced first-party before being acted on, per the standing rule
that a judge's finding is a claim like any other — and both are **genuine survivals**: mutant
LANDED (differing sha256), mutant BUILDS (`go build ./...` rc=0), and the **entire
`host/transitionreg` package rc=0 with the defect present**.

| branch | file | survived as delivered |
|---|---|---|
| `revision == 0 && !parent.IsZero()` | `transitionreg.go:109` | package rc=0 |
| `revision > 1 && parent.IsZero()` | `transitionreg.go:112` | package rc=0 |
| `GetObject` error wrapper | `transitionreg.go:89` | package rc=0 |

The first two are Decision 3's parent/revision chain rules; with both neutered together,
`ReadSnapshot` **silently accepts** a `Revision{Revision: 5, Parent: zero}` as a valid `Snapshot`,
so a corrupted or tampered object with a broken revision chain reads as valid. The third is the
generic store-error wrapper: `fakeObjectStore.objectErr` was declared and wired into the fake's
`GetObject` but **never set to non-nil by any test**, so an injected read error was swallowed
undetected. That is the same defect class `MUT-READ-SWALLOW` exists to catch, one call deeper.

**Why the sweep missed them, which is the part worth keeping.** The plan's §4.3 rule-3j audit was
genuinely thorough — eleven previously-unmutated frozen branches — but it enumerated branches
frozen by Decisions **1, 2 and 4**. Decision **3**'s own chain rules were never audited, and
neither were the store-error wrappers **T5 itself introduces**. So the gap is not an oversight
*within* the audit's scope; it is the audit's scope. Rule 3j says the unit of mutation is the
BRANCH, and an audit anchored to a list of *decisions* silently inherits that list's boundaries —
including branches written during the sprint, which by construction no pre-sprint enumeration can
contain.

Three subtests added to `TestReadSnapshotRefusals` (`injected_object_read_error`,
`revision_zero_with_parent`, `revision_after_one_without_parent`), each pinning its own measured
message. AC2's count is unchanged at **3** — these are subtests, not new top-level names.

| arm | LANDED | BUILDS | KILL | failing subtest | INVERSE | RESTORE |
|---|---|---:|---:|---|---:|---|
| `MUT-REV-CHAIN-ZERO-PARENT` | `5aecd46b…` -> `fa9739aa…` | 0 | 1 | `…/revision_zero_with_parent` | **0** | match |
| `MUT-REV-CHAIN-MISSING-PARENT` | `5aecd46b…` -> `ec80cf7e…` | 0 | 1 | `…/revision_after_one_without_parent` | **0** | match |
| `MUT-READ-OBJECT-SWALLOW` | `5aecd46b…` -> `609ace1f…` | 0 | 1 | `…/injected_object_read_error` | **0** | match |

All three are **isolated** kills — each redded by exactly its own subtest, each with a green
package inverse — so unlike the two arms adjudicated above they satisfy §4.1 step 7 unmodified.

**Also fixed, non-blocking, and instructive about the gate rather than the code:**
`go vet ./host/transitionreg/...` was **rc=1** with five `copylocks` findings — `f := *base`
copies `fakeObjectStore`, which embeds a `sync.Mutex`. It escaped review because
`go test`'s default vet subset **excludes** the `copylocks` analyzer and `verify_go.sh` never
invokes `go vet` directly, so both the local gate and CI were green on it. Replaced with an
explicit `clone()` that constructs a fresh fake; `go vet` is now rc=0. Worth noting as a standing
gap: **a green `go test` is not a green `go vet`**, and nothing in this repo's gate currently
closes the difference.

**The general point this pair settles, and the reason it is recorded rather than waved through:**
the plan's §4.1 step 7 requires a green package inverse to prove *your* test is the killer rather
than a bystander. That is the right check for a mutation whose blast radius is one branch, and it
is **unsatisfiable by construction** for a mutation that alters a shared constant or a real
behavioural invariant — the mutant is *supposed* to be visible to more than one test. Read
strictly, step 7 would have scored two correct, well-pinned kills as failures, and the cheapest way
to make them "pass" is to weaken the co-detectors. The executor declined to do that and said so,
which is what made the correct adjudication available. The generalised form of the check, which
holds for both narrow and broad mutations: **skip the complete set of tests that legitimately
detect the mutation and require rc=0, with the same skip on the unmutated tree as the control.**
