# Evaluator r1 - M7c3 (row 150), tree 47686112, base c448b15f

VERDICT: PASS. Score 90/100.

| Criterion | Max | Score |
|---|---|---|
| Plan conformance | 30 | 27 |
| Tests and gates | 25 | 24 |
| Mutation / non-vacuity | 20 | 16 |
| Premise and risk | 10 | 8 |
| Hygiene | 15 | 15 |

## Gates (run outside sandbox, AILANG_BIN=$HOME/.pinned-ailang/ailang)
All green. My first harness printed blank rc values (zsh PIPESTATUS), so rc comes from the tool output: no output from gofmt and vet, and every `go test` printed `ok`.
| Gate | Result |
|---|---|
| gofmt -l scripts/marketingcheck host/daemon | empty (clean) |
| go vet ./scripts/marketingcheck ./host/daemon | clean |
| go test -run '^$' ./... | all ok |
| TestMarketingCaptureInventoryMutation (host/daemon, -v) | PASS 0.62s, commits 0-2 200 |
| TestMarketingInventory* (scripts/marketingcheck, -v) | PASS: CompareBranches 14 subtests, ReadPaging 2 subtests |
| TestWorkbenchLiveCursorCrossesGap | ok |
| TestProductionGoSurface, TestProductionContextRoots (host/store) | ok |
| TestSubprocessSinkGate (host/verifygate) | ok |
| go test ./scripts/marketingcheck -count=1 -timeout 600s (full package) | ok 31.6s |

## Red-first
Base `main.go` (`git show c448b15f:scripts/marketingcheck/main.go` copied over the file):
- host/daemon test FAILS: `before inventory: rc 2: usage: marketingcheck check ...` (unknown subcommand).
- scripts/marketingcheck test does not compile: `undefined: Inventory / inventoryEntry / inventoryObject`.
- Restored; `sha256sum -c` OK and `git diff --quiet` rc 0.

## Mutants
Plan mutants (valid Go, run one at a time):
| Mutant | Result | Failing assertion |
|---|---|---|
| M7c-inventory (Compare returns nil after empty check) | KILLED | daemon: `marketing_inventory_test.go:76: extra commit invalidates capture: rc 0:` (also unit tests head/entry/count/object-*) |
| M7c3-paging (limit=1, stop after first page) | KILLED | daemon: `:88: after.json must list all 4 committed entries: got 1` (plus before.json 3 -> got 1). The unit test also reds on its `limit != 500` guard |

Own mutants, branch-complete over ReadInventory, CompareInventory and runInventory. NOCOMPILE and SKIP-APPLY variants were rewritten into valid Go. The a-series results below that were NOCOMPILE or SKIP-APPLY were rerun as b/c variants.
| Mutant | Result | Assertion |
|---|---|---|
| head json.Valid wrap off | KILLED | `inventory_test.go:104 text head must be kept as a JSON string` |
| head fallback replaced | KILLED | same assertion at :104 |
| no status check (non-200) | KILLED | `:125 a non-200 object read must fail the inventory` |
| object fetch err ignored | KILLED | same assertion at :125 |
| from += 1 | KILLED | `:107 every page must be read: got 9 entries, want 5` |
| limit 500 -> 5 | KILLED | server 400 on limit |
| objects unsorted | KILLED | `:120 objects must be the distinct refs, sorted` |
| digest of ref instead of body | KILLED | `:120` |
| empty-page break (`< 0`) | KILLED | test panic/timeouts |
| empty-check: drop null / blank / len==0 / no-entries / only-before | KILLED | `:60` per-case inventory-empty |
| head cmp / count cmp / entry hash cmp / entry ref cmp / obj len / obj hash cmp removed | KILLED | `:60` per-case inventory-mutated |
| code string changed | KILLED | `:60` |
| final CLI rc 1->2 | KILLED | daemon (rc 1 asserted) |
| compare skipped | KILLED | daemon |
| --out ignored | KILLED | daemon |
| TrimRight("/") removed | SURVIVED | cosmetic (no trailing-slash test) |
| head fetch err ignored | SURVIVED | NON-BLOCKING |
| log fetch err -> break | SURVIVED | NON-BLOCKING (see findings) |
| page unmarshal err ignored | SURVIVED | NON-BLOCKING (see findings) |
| TrimSpace in empty check | SURVIVED | NON-BLOCKING |
| subject label before/after in the inventory-empty message | SURVIVED | NON-BLOCKING |
| obj presence (`!ok`) removed | SURVIVED | equivalent mutant (missing ref gives hash "" which differs from any sha256) |
| inventory/compare CLI usage paths: NArg check, addr required, before/after required, ErrHelp rc, parse-err rc, usage rc 2->1 | SURVIVED | NON-BLOCKING (no CLI usage tests) |
| unreadable / inventory-json refusal wrap; write err ignored; unmarshal err ignored; before/after read swap | SURVIVED | NON-BLOCKING |
| schema string changed | SURVIVED | NON-BLOCKING |
| 10 s client timeout removed | SURVIVED (after rewrite to valid Go) | inherent, a time bound is untestable |

## Restoration proof
`sha256sum -c /tmp/world-iter310/main.sha` printed `scripts/marketingcheck/main.go: OK` after every mutant batch. `git diff --quiet` in the worktree returned rc 0 (no tracked change). Only untracked mutant artifacts exist under /tmp/world-iter310/mut. Nothing was written under $HOME or ~/.ailang/state.

## Conformance check
- Files: main.go plus the new host/daemon test. The controller's `scripts/marketingcheck/inventory_test.go` is test-only. No new non-test .go file under scripts/. No context.Background/TODO/WithoutCancel in main.go (grep empty).
- Schema, paging (limit=500 from 0 until a zero page), distinct-ref sorted digests, and codes `inventory-mutated` / `inventory-empty`: match the plan. The CLI exits 0/1/2 like `check`, with `marketingcheck: <code> ...` on stderr.
- The daemon test follows the 8 steps in order. Subprocess sinks are `os.Create` files in t.TempDir. The 60 s build context and the assertion texts match the plan. The paging assertion "after.json must list all 4 committed entries" is present.
- Deviation: `writeObservation` was renamed to `writeJSON(any, out)`. This is a small refactor inside the permitted file, with `observe` still tested. Extra codes `inventory-unreadable` and `inventory-json` are not in the plan; this is reasonable and not forbidden.

## Premise and risk
- `handleHead` (host/daemon/daemon.go:1117-1133) writes `text/plain` `ref.String()`. Encoding the plain-text head as a JSON string is correct and necessary, and `ReadPaging/pages` pins it.
- `handleObject` (host/daemon/handlers.go:483-513) returns the envelope without payload unless `?payload=true`. The envelope carries the content hash, interfaceHash, semanticId and provenance. A no-payload digest therefore covers identity, and the content-addressed hash commits to the payload. This is sound and bounded (the handler rationale at :476-481 says the payload is unbounded).
- The guard detects any commit, since the head, the entry chain and the transition objects all change. It cannot detect writes that bypass the log, non-transition object puts, or an A-to-B-to-A revert (inherent). Reading head, log and objects is not atomic, so a commit landing between the three reads could skew one snapshot (low risk).
- `readCtx` timeout 404s are fatal (fail closed), which is good.

## Findings
BLOCKING: none.
NON-BLOCKING:
1. Head/log/page error handling survivors (main.go `ReadInventory`). If a mutant swallows the head or log error, the result is an empty or short inventory. `inventory-empty` and the compare refusals would still catch an empty one, so the severity is low. A short page-read failure could yield a truncated-but-equal pair only if both reads truncate identically. Suggest a test where /v1/log returns 500 or malformed JSON.
2. The CLI usage and flag-validation paths are untested (exit 2 for missing --addr/--before/--after, extra args, help). Plan says exit-code format matches `check`; rc 2 for these paths is unverified. Low severity.
3. Test coverage of the daemon test alone cannot attribute the refusal branch (controller's unit test compensates).
4. Executor's REPORT socket results are uninformative as the task states; independently re-verified here.
5. `inventory-empty` subject label (before/after) is unpinned: the unit test only checks detail text.
