# M7c3 — Frozen public inventory guard

Implementation is ready for controller verification outside the sandbox. Socket-dependent gates and all mutant runs are **UNINFORMATIVE UNDER SANDBOX**; no acceptance pass or mutant kill is claimed for them.

## Files and behavior

- `scripts/marketingcheck/main.go`: 521 total lines, +166 net lines from the 355-line base. Adds `Inventory`, `ReadInventory`, `CompareInventory`, dispatch and help for `inventory` and `inventory-compare`. Reuses the existing JSON writer by generalizing and renaming it to `writeJSON`; its exclusive file-create behavior is preserved.
- `host/daemon/marketing_inventory_test.go`: new, 91 lines. Implements all eight steps in order, genuine authenticated REST commits, file sinks, 60-second subprocess bounds, equality/refusal assertions, and before/after counts 3/4.
- Combined net addition: 257 lines (approximately 250). No required assertion was shortened or omitted.
- Full final source snapshots are `scripts/marketingcheck/main.go.txt` and `host/daemon/marketing_inventory_test.go.txt`, relative to this directory. Evidence logs, commands/return codes, mutation diffs and restoration records are under `evidence/`.
- No git write operation was performed. No changes to `scripts/capture_marketing.sh` or other production/test paths. All task artifacts are inside this worktree.

## Red-first and sandbox record

The requested initial probe `python3 -c "import socket; s=socket.socket(); s.bind(('127.0.0.1',0)); print('bind ok')"` returned rc 1: `PermissionError: [Errno 1] Operation not permitted`.

The test was written before production edits and run against base production code. It logged three genuine REST commits at 200, then returned rc 1 at the server bind. Exact obstruction:
```text
panic: httptest: failed to listen on a port: listen tcp6 [::1]:0: bind: operation not permitted [recovered, repanicked]
```
This red-first run is **UNINFORMATIVE UNDER SANDBOX**: it did not reach the inventory invocation. Raw output: `evidence/red-first.log`.

A separate base binary build returned rc 0. Invoking that unmodified binary as `inventory --addr http://127.0.0.1:1 --out snap/M7c3/evidence/base-before.json` returned rc 2 and printed the old usage, containing only `check`, `routes`, `observe`, and `compare`. The dispatch evidence confirms the intended unknown-subcommand red without claiming the eight-step test reached it. Temporary binaries were removed.

## Final gates

Each command below ran on the final source, before the final isolated mutant sequence. Return codes are those of the command itself, also recorded in `evidence/gate-*.json`.

| Command | rc | Result |
|---|---:|---|
| `gofmt -l scripts/marketingcheck host/daemon` | 0 | PASS; no output |
| `go vet ./scripts/marketingcheck ./host/daemon` | 0 | PASS |
| `go test -run '^$' ./...` | 0 | PASS |
| `go test -run '^TestMarketingCaptureInventoryMutation$' ./host/daemon -count=1 -timeout 180s -v` | 1 | UNINFORMATIVE UNDER SANDBOX — denied loopback bind |
| `go test -run '^TestWorkbenchLiveCursorCrossesGap$' ./host/daemon -count=1 -timeout 90s -v` | 0 | PASS |
| `go test -run '^TestProductionGoSurface$|^TestProductionContextRoots$' ./host/store -count=1 -timeout 60s -v` | 0 | PASS |
| `go test -run '^TestSubprocessSinkGate$' ./host/verifygate -count=1 -timeout 120s` | 0 | PASS |
| `go test ./scripts/marketingcheck -count=1 -timeout 300s` | 1 | UNINFORMATIVE UNDER SANDBOX — denied loopback bind |

The marketingcheck package stops in the pre-existing `TestMarketingCapturePassFakeChrome/ready` at the same denied bind. The new daemon test stops before building the CLI. These are environmental obstructions, not test verdicts.

Additional CLI checks: both `-h` variants rc 0; missing required arguments and unknown flag rc 2; malformed before and after files rc 1 with `inventory-json`; empty comparison rc 1 with `marketingcheck: inventory-empty before: head or entries empty`. Exact outputs: `evidence/cli-checks.json`. `git diff --check` returned rc 0.

## Mutants

Applied one at a time in production `main.go`; only the focused `TestMarketingCaptureInventoryMutation` command was run for each. Each source was restored from the final `main.go.txt` snapshot and immediately checked using `cmp`. Mutation patches, raw logs and command records are preserved.

| Mutant | Edit | Test rc | Exact failing assertion | Verdict | restoration cmp rc |
|---|---|---:|---|---|---:|
| M7c-inventory | Return nil after both empty-input checks | 1 | None reached: denied bind | UNINFORMATIVE UNDER SANDBOX | 0 |
| M7c3-paging | Request limit=1 and stop after first nonempty page | 1 | None reached: denied bind | UNINFORMATIVE UNDER SANDBOX | 0 |
| own-head-unchecked | Disable head comparison condition | 1 | None reached: denied bind | UNINFORMATIVE UNDER SANDBOX | 0 |
| own-entries-count-only | Disable per-entry equality condition; retain count check | 1 | None reached: denied bind | UNINFORMATIVE UNDER SANDBOX | 0 |
| own-objects-unchecked | Return nil before object-set comparison | 1 | None reached: denied bind | UNINFORMATIVE UNDER SANDBOX | 0 |

Exact obstruction for every mutant is the panic quoted above. None reached a named failing assertion; none is classified killed or survived. The controller must rerun these outside the sandbox.

Named expected assertions, not observed here:
- M7c-inventory: `extra commit invalidates capture`.
- M7c3-paging: `after.json must list all 4 committed entries` (and `before.json must list all 3 committed entries`).

Coverage finding from source inspection, not measured mutant survival: the extra genuine commit changes head, count and objects together, so isolated removal of head, per-entry or object comparison is not guaranteed to be detected by this integration test. The three additional branch mutants record that remaining coverage question honestly.

## Plan interpretation and remaining verification

- Real `/v1/head` is canonical **plain text**, per `host/daemon/daemon.go:1108-1140`, not JSON. The inventory encodes this exact body as a JSON string; already-valid JSON bodies remain JSON. This necessary representation makes the planned JSON inventory work with the genuine store. Empty/null/empty-string head and zero-entry inventories refuse in `CompareInventory` for either input.
- The sandbox prevents observing the full red-first failure at unknown `inventory`, the acceptance test, and every mutant assertion. The prescribed test order and assertions remain intact; no skips or socket workaround were added.
- The optional empty-inventory test arm was not added; a direct CLI empty-input check was recorded instead. No scope was added for shell wiring.
- Controller must rerun socket-dependent gates and all mutants outside the sandbox. This is not real-Chrome capture evidence and does not mark C8/C9 met.
