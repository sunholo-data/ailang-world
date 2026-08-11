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
