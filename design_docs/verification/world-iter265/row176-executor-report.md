Implemented the binding and QUICKSTART changes. **AC8 remains red**: serialization measures `max=1` and the unwrapped CONTROL measures `max=2`, but concurrent Dispatch calls conflict at the coordinator’s world-head commit. The required nil-error assertions remain intact.

`git status --porcelain` lists exactly:

```text
 M docs/QUICKSTART.md
 M host/daemon/daemon.go
 M host/daemon/decisionpacket_daemon_test.go
 M host/daemon/workbench_packet_test.go
```

Red-first evidence:

- **(a), rc=1:** `d.human undefined (type *Daemon has no field or method human)`, `d.humanMu undefined`, and `undefined: serializedHandler`.
- **(d), rc=1:** `workbench_packet_test.go:526: section lacks --grant Human.Approve=worktree:1`.

After implementation:

```text
--- PASS: TestDaemonInvocationMintsDecisionPacket (0.02s)
--- PASS: TestQuickstartDecisionPacketSection (0.03s)
```

The two-test gate returned **rc=1** after implementation and on final verification:

```text
serialized probe max = 1
CONTROL unwrapped probe max = 2
--- FAIL: TestHumanApproveIsSerializedAcrossEpisodes (0.43s)
```

Both AC8 arms reported `store: stale observed world head`; AC8 green was not obtained.

Mutation ledger:

| Mutation | Observed red | rc | Restoration |
|---|---|---:|---|
| MD2: delete registry insertion | `declares effects with no registered handler [Human.Approve]; effectful invocation is unavailable` | 1 | Reverted |
| MD3: delete Lock/Unlock, run 1 | `serialized probe max = 2, want 1` | 1 | Reverted after run 3 |
| MD3, run 2 | `serialized probe max = 2, want 1` | 1 | Reverted after run 3 |
| MD3, run 3 | `serialized probe max = 2, want 1` | 1 | Reverted |

After each mutation’s restoration, `git diff -- host/daemon/daemon.go` matched the saved implementation exactly. Final comparison also matched. No MD3 run was green.

Fast-check results:

| Check | rc | Result |
|---|---:|---|
| `go vet ./...` | 0 | Passed |
| `go test -run '^$' ./...` | 0 | Compile fence passed |
| Requested daemon subset | 1 | **UNINFORMATIVE UNDER SANDBOX**: `TestAGUISeesDirectStoreCommit` panicked on `bind: operation not permitted` |
| TR.C `TestRegistryDispatchBindingBoundary` | 0 | Passed unchanged: inside count 3, outside 0 |
| Requested grep | 0 | Completed |
| Discovered source checks | 0 | MCP mount, MCP budget, frozen route table, workbench live QUICKSTART all passed |
| Additional workspace/registry checks | 0 | Executed tests passed; six real-binary tests skipped |
| `git diff --check` | 0 | Passed |

TR.C’s green line:

```text
--- PASS: TestRegistryDispatchBindingBoundary (0.13s)
```

Scope verification confirmed `"worktree"` matches: [parseGrant](/Users/voightkampff/dev/sunholo-data/.wt-world-iter265-row176/cmd/ailang-worldd/session.go:394) preserves the supplied scope; [OpenBinder](/Users/voightkampff/dev/sunholo-data/.wt-world-iter265-row176/host/broker/binder.go:22) passes the grants into the session; [scopeMatches](/Users/voightkampff/dev/sunholo-data/.wt-world-iter265-row176/host/broker/decide.go:111) requires exact equality; [WorkspaceScope](/Users/voightkampff/dev/sunholo-data/.wt-world-iter265-row176/host/broker/handlers_ailang.go:40) is `"worktree"`.

**UNMEASURED:** outside-sandbox daemon subset and six skipped tests requiring unset `WORLD_TOOL_AILANG_BIN`: startup variants, SeTools A2A failure mapping, module imports, package imports, workspace sandbox, and production workspace reads. Race testing remains controller-owned.

Every Go test exported the required `AILANG_BIN`. No git write operations or changes to other files were made.