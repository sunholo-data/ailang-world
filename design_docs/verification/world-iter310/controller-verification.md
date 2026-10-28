# World iteration 310 — controller verification of row 150 M7c3

Base `c448b15fd6ee9c9e5e221cb66bf099cb7bcb127c` (Gate 1, Gate-3 snap and Gate 4 all read the same SHA).
Sprint worktree `.wt-world-iter310-row150-m7c3`, branch `sprint/world-iter310-row150-m7c3`, PR #393.
All commands below ran OUTSIDE the codex sandbox with `AILANG_BIN=$HOME/.pinned-ailang/ailang` (v0.41.0), load 1–6.

## Baseline (pristine worktree at base, before the directive was sent)

| Gate | Result |
|---|---|
| `go vet ./scripts/marketingcheck ./host/daemon` | rc 0 |
| `TestWorkbenchLiveCursorCrossesGap` | ok |
| `TestProductionGoSurface` | ok |
| `TestProductionContextRoots`, `TestSubprocessSinkGate` | ok |
| `gofmt -l scripts/marketingcheck host/daemon` | empty |

Plan premises checked at HEAD: `newHandlerDaemon` (`handlers_test.go:21`), `authHeader` (:47), `requestRecorderAuth` (:59),
`testCommit` (:78), `workbenchLiveREST` (`workbench_live_test.go:74`), `logEntryResponse` (`handlers.go:65-73`),
`handleLogRange` stops at the first missing index (`handlers.go` loop `if !ok { break }`), `clampLimit` hard cap 500,
and `isProtected` (`daemon.go:1079`) guards only `POST /v1/commit`, so every read the inventory makes is unauthenticated.

## Executor

`codex:gpt-6.1-sol`, wrapper `/tmp/codex_run_iter310.sh` (probe rc 0; run rc 0, 349 s, 96,227 tok, not capped).
Its sandbox denied loopback binds, so the inventory test, the package run and all five of its mutants were
UNINFORMATIVE; it claimed no kill. Report: `executor-report.md`. Its snapshot (`snap/M7c3/`) was moved to
`/tmp/world-iter310/exec-snap/` before committing; it is not in the tree.

Diff read in full: `scripts/marketingcheck/main.go` +166 net (two subcommands, `Inventory`, `ReadInventory`,
`CompareInventory`, `runInventory`; `writeObservation` generalised to `writeJSON(any, out)`), and
`host/daemon/marketing_inventory_test.go` (91 lines). No write outside `--out`; no `context.Background/TODO/WithoutCancel`;
no new non-test `.go` under `scripts/`. Two refusal codes beyond the plan (`inventory-unreadable`, `inventory-json`)
cover I/O and malformed input, which the plan left unnamed. Deviation adjudicated (rule 3h): `/v1/head` serves
`text/plain` (`daemon.go:1117`), so the body is stored as a JSON string when it is not valid JSON — correct, and
`head-text-fallback-dropped` below shows the alternative fails the real daemon.

## Gates on the final tree `47686112`

| Gate | Result |
|---|---|
| `gofmt -l scripts/marketingcheck host/daemon` | empty |
| `go vet ./scripts/marketingcheck ./host/daemon` | rc 0 |
| `go test -run '^$' ./...` (compile fence) | rc 0, 34 packages |
| `TestMarketingCaptureInventoryMutation` | PASS (0.47 s) |
| `TestMarketingInventory*` (controller unit tests) | PASS |
| `TestWorkbenchLiveCursorCrossesGap` | ok |
| `TestProductionGoSurface`, `TestProductionContextRoots` | ok |
| `TestSubprocessSinkGate` | ok |
| `go test ./scripts/marketingcheck ./host/daemon` | ok (31.4 s / 250.7 s) |
| `go test -race ./scripts/marketingcheck` | ok (73.0 s) |
| `go test -race -run 'TestMarketingCaptureInventoryMutation|TestWorkbenchLive' ./host/daemon` | ok |

`.ail` untouched, so `verify_ail.sh` was not run (CI's `ailang-code verify gate` covers it).

## Red-first

Base `main.go` over the file (controller unit test held aside): the e2e test fails at step 5,
`before inventory: rc 2: usage: …` (unknown subcommand). Restored; `shasum -c` OK on both files.

## Mutation drill (`controller-mutants.py.txt`; one mutant at a time, anchors asserted unique, sha restored)

Round 1 ran the e2e test only. The two plan mutants were first written wrongly by the controller and re-run:
`return nil` after the empty loop fails `go vet` (unreachable code), so it became `if true { return nil }`; and the
plan's paging mutant is BOTH `limit=1` AND stop-after-first-page — `limit=1` alone still pages through everything.

| Mutant | e2e `TestMarketingCaptureInventoryMutation` | controller unit tests |
|---|---|---|
| M7c-inventory (plan) | KILLED — "extra commit invalidates capture: rc 0" | KILLED |
| M7c3-paging-full (plan) | KILLED — "after.json must list all 4 committed entries: got 1" | KILLED |
| head-not-compared | survived | KILLED — "head must refuse inventory-mutated (head differs)" |
| entries-count-only | survived | KILLED — "entry must refuse … (entry 1 differs)" |
| objects-len-ignored | survived | KILLED — "object-extra must refuse …" |
| objects-hash-ignored | survived | KILLED — "object-hash must refuse …" |
| non200-ignored | survived | KILLED — "a non-200 object read must fail the inventory" |
| empty-check-removed | survived | KILLED — "empty-head-before must refuse inventory-empty" |
| refusal-exit-0 | KILLED — "extra commit invalidates capture: rc 0" | survived (unit tests call functions) |
| no-objects-read | survived | KILLED — "objects must be the distinct refs, sorted, with body digests" |
| head-text-fallback-dropped | KILLED — "before inventory: rc 1 … invalid character 's'" | KILLED |
| paging-single-page | survived | KILLED — "every page must be read: got 2 entries, want 5" |
| M7c3-paging limit=1 only | survived (not the plan's mutant) | KILLED (fake server rejects limit≠500) |

Why the e2e test alone leaves branch survivors: one genuine extra commit changes the head AND the entry list, and the
test store reuses one transition object, so each comparison branch is redundant inside that single scenario (rule 3j:
for a refusal deliverable, the unit of mutation is the branch). Controller test-only commit `47686112` adds
`scripts/marketingcheck/inventory_test.go`: `TestMarketingInventoryCompareBranches` (14 subtests, one per refusal branch
and side) and `TestMarketingInventoryReadPaging` (a fake daemon whose page cap is 2, a plain-text head, two distinct
objects, and a 404 arm). After it, every mutant above is killed by at least one suite.

## Independent judge

sonnet r1 PASS 90/100, zero blocking (`evaluator-r1.md`). Controller reproduced its CLI finding first-party: the usage
paths are correct but untested (`inventory` with no `--addr` rc 2; `inventory-compare` missing `--after` rc 2; `-h` rc 0;
unreachable address rc 1 `inventory-unreadable`).
