# Controller gates — iteration 378

## Sweep 1 at `784c9413` (S3 head): rc 1

```
69:ok  	github.com/sunholo-data/ailang-world/host/broker	(cached) [no tests to run]
213:--- FAIL: TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential (0.53s)
229:--- FAIL: TestWorktreeSnapshotRefusesAConfiguredFilter (0.19s)
231:FAIL
232:FAIL	github.com/sunholo-data/ailang-world/host/broker	340.227s
262:FAIL
--- FAIL: TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential (0.53s)
    registry_publish_test.go:1165: AC10(a) re-derived N = 13 production subprocess launch sites in 13 files:
    registry_publish_test.go:1168:     cmd/ailang-worldd/approve_sig.go:153 exec.CommandContext
    registry_publish_test.go:1168:     host/agentrun/agentrun.go:46 exec.CommandContext
    registry_publish_test.go:1168:     host/archive/archive.go:479 exec.CommandContext
    registry_publish_test.go:1168:     host/archive/check.go:95 exec.CommandContext
    registry_publish_test.go:1168:     host/broker/handlers.go:159 exec.CommandContext
    registry_publish_test.go:1168:     host/capsule/capsule.go:236 exec.CommandContext
    registry_publish_test.go:1168:     host/checkrun/checkrun.go:98 exec.CommandContext
    registry_publish_test.go:1168:     host/evidence/proof_producer.go:178 exec.CommandContext
    registry_publish_test.go:1168:     host/pkgproj/iface.go:171 exec.CommandContext
    registry_publish_test.go:1168:     host/pkgproj/pkgproj.go:259 exec.Command
    registry_publish_test.go:1168:     host/replay/replay.go:336 exec.CommandContext
    registry_publish_test.go:1168:     host/toolprobe/toolprobe.go:43 exec.CommandContext
    registry_publish_test.go:1168:     host/worktree/worktree.go:138 exec.CommandContext
    registry_publish_test.go:1228: AC10(a): 13 files carry subprocess sites [cmd/ailang-worldd/approve_sig.go host/agentrun/agentrun.go host/archive/archive.go host/archive/check.go host/broker/handlers.go host/capsule/capsule.go host/checkrun/checkrun.go host/evidence/proof_producer.go host/pkgproj/iface.go host/pkgproj/pkgproj.go host/replay/replay.go host/toolprobe/toolprobe.go host/worktree/worktree.go] but 12 have drivers; every enumerated site must be DRIVEN
--- FAIL: TestWorktreeSnapshotRefusesAConfiguredFilter (0.19s)
    worktree_delta_test.go:432: control: plain git status did not run the clean filter (stat /private/var/folders/kr/h0wr2sj94vd6ljtmsxv8jkt00000gn/T/TestWorktreeSnapshotRefusesAConfiguredFilter4090633584/002/filter-ran: no such file or directory); the arm below would prove nothing
FAIL
```

- AC10 census red = this sprint's S3 (new exec site in package main, unscrubbed env, undriven) → fixed by S4.
- TestWorktreeSnapshotRefusesAConfiguredFilter: known load-sensitive control arm; 5/5 ok alone at the same tree (`go test ./host/broker/ -run '^TestWorktreeSnapshotRefusesAConfiguredFilter$' -count=5`); sprint touches neither host/broker production code nor host/worktree.

## Sweep 2 at `5e85f8bb` (S4 head): rc 0

```
rc=0 for 5e85f8bb
ok lines: 114
FAIL lines: 0
ok  	github.com/sunholo-data/ailang-world/cmd/ailang-worldd	0.490s [no tests to run]
ok  	github.com/sunholo-data/ailang-world/host/broker	0.260s [no tests to run]
ok  	github.com/sunholo-data/ailang-world/host/sshsign	0.572s [no tests to run]
ok  	github.com/sunholo-data/ailang-world/cmd/ailang-worldd	116.789s
ok  	github.com/sunholo-data/ailang-world/cmd/ailang-worldd	17.535s
ok  	github.com/sunholo-data/ailang-world/cmd/ailang-worldd	0.489s
ok  	github.com/sunholo-data/ailang-world/cmd/ailang-worldd	7.695s
ok  	github.com/sunholo-data/ailang-world/cmd/ailang-worldd	224.273s
ok  	github.com/sunholo-data/ailang-world/host/broker	344.330s
ok  	github.com/sunholo-data/ailang-world/host/sshsign	3.021s
   ✓ canonical JSON equals committed golden byte-for-byte
✓ world package gate PASSED: 9/9 steps performed non-zero work
✓ verify gate PASSED: 16 required identities verified, 40 named tests pass, 503 se-tools named tests pass
   ✓ 3 pinned font(s) validated, 1017220 bytes (host/surfacepage/fonts/)
   ✓ 0 other binary blobs among 6648 tracked files
   ✓ toolchain floor gate: go1.26.6 >= root module floor go1.26.6
   ✓ all 64 required top-level evidence tests passed exactly once
✓ go gate PASSED (plain): build clean; legs run pass with pinned AILANG_BIN (AILANG v0.41.0
```

The race-detector known-positive control printed its expected `WARNING: DATA RACE` (racecontrol/main.go:16); not a finding.

## Unattributed single red

First `commit_ms.sh S1` run: `go test ./cmd/ailang-worldd/` FAIL at 175.6 s; log lost (piped through `tail`). Re-run on the byte-identical S1 tree (8/8 files cmp-equal to the snapshot): ok 112.6 s; later full runs of the package at S2, S3, S4 and both sweeps all ok. Recorded, not attributed.
