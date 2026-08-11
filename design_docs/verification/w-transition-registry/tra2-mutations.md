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
| MUT-DELETE-TR-A-TEST (AC2) | AC2 names `3->2` | see T8 section below | 0 | activated AC2 rc=1, observed count=2 | 0 after restore | match | KILLED |
| MUT-DELETE-TR-A-TEST (AC3) | AC3 names `4->3` | see T8 section below | 0 | activated AC3 rc=1, observed count=3 | 0 after restore | match | KILLED |

The semantic-tag arm exposes a plan-level coupling: changing the package semantic constant affects
many tests, so the required package inverse cannot be green. The CAS-swallow arm likewise changes
real CAS-conflict behavior covered by `TestPublishCASConflictPreservesWinner`; skipping only the
injected-error subtest cannot make the package green. Both mutants landed and built, and both named
tests fired, but neither is represented as an isolated mutation kill.
