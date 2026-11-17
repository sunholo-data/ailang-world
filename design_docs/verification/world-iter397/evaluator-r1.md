# Evaluator report r1 — iter 397, row 243 V2a "the seam" (base 18d30a17, head 5d71dc69)

## Verdict: PASS — 90/100. No BLOCKING findings.

| Category | Score |
|---|---|
| Correctness vs design/plan §9 (behaviour neutrality first) | 28/30 |
| Tests: non-vacuous, kill mutants, cover every operand | 25/30 |
| Security / fail-closed / purity | 20/20 |
| Coding standards | 9/10 |
| Scope and size | 8/10 |

## Behaviour neutrality (my own read of the rename diff + CLI file)
- `git diff -M 18d30a17 HEAD -- cmd/ailang-worldd/why.go host/why/why.go` is a mechanical move: identifiers exported, `net/http`/`os`/`flag`/`time` removed, `r.c.do(GET)` replaced by `r.Get`, `http.Status*` by literal 200/404, flag/help/target-read/--expect-store code moved verbatim into `why_cli.go`. I found no semantic change in `GetJSON`, `Head`, `ScanEntries`, `Locate`, `BuildChain`, `Elide`.
- `Run` buffers into two builders and the deferred closure assigns `a.Stdout/a.Stderr` after each `return Answer{Exit:…}`. The CLI prints stdout, then stderr, and returns Exit. The split and exits (0/1/3) hold on every path. Verified by mutants E-RUN-SWAP, E-RUN-INTEG-EXIT, E-RUN-DEADLINE-EXIT, E-CLI-NOSTDERR and E-CLI-NOSTDOUT, all killed.
- `ParseTarget("-")` falls to the integer error. The CLI intercepts `-` before calling it, as the base ordering did. Order in the CLI (`--result` first, then `-`, then ParseTarget) matches the base switch.
- VC1 golden is a real base capture: I copied `why_characterise_test.go` and the golden into a throwaway worktree at 18d30a17 and `TestWhyCharacterisedPaths` PASSES there (the worktree is removed). Golden also passes at head.
- Shims: `entryWire`/`worldWire` are defined types over `why.EntryWire/WorldWire` (vet clean); `sum`, `sumJSON`, the 13 delegations and `whyTarget`→`why.Target` copy every field. No semantic change found. `elideAt` is removed; `elide` is a one-line delegation (VW5).

## Real-stack leg (RAN, not SKIP)
`AILANG_BIN=$HOME/.pinned-ailang/ailang WORLD_TOOL_AILANG_BIN=$HOME/.pinned-ailang-tools/v0.54.15/ailang go test -count=1 -v ./cmd/ailang-worldd/ -run '^(TestDemoSeedWhyNamedLinks|TestCLIAgainstSeToolsDaemon)$'`: rc=0; `--- PASS: TestCLIAgainstSeToolsDaemon (6.42s)`, `--- PASS: TestDemoSeedWhyNamedLinks/real-stack (5.61s)` plus its 14 control subtests. Log: `judge/real.log`.

## Test runs (AILANG_BIN pinned)
- `go test -count=1 -v ./host/why/`: rc=0 (28 PASS lines).
- `go test -race -count=1 ./host/why/`: ok.
- `go test -count=1 -timeout 30m ./cmd/ailang-worldd/`: rc=0, 117.9 s.
- `go vet ./host/... ./cmd/...`: rc=0. `go test -run '^$' ./...`: no failures.
- Not run, per directive: verify_go.sh, whole `go test ./...`.

## Pre-existing test edit audit
- `why_input_test.go` renamed to `host/why/` (R095): 3+/3-, only the package line and `elide`→`Elide` (two sites). Nothing weakened.
- `why_protected_test.go`: 66 lines removed (18 pure subtests of `TestWhyProtectedLine` and `TestWhyProtectedJSON`), moved to `host/why/why_protected_test.go` (104 added, `whyLink`→`Link`). The two e2e subtests stay in place.
- `go test -v -run TestWhy` pass-line set, base vs head: base 33, head 40. `comm -23 base head` is EMPTY (every base test and subtest name is present at head). The 7 new names are TestWhyCharacterisedPaths, ConstantsMatchTheCLI, HeadIsOneReader, ImportsNoEffect, RunGetterClasses, RunReplaysTheCLI, plus the moved TestWhyProtectedLine in the second package. Files: `judge/base_why.txt`, `judge/head_why.txt`.

## Size / scope
- `git diff -M --numstat 18d30a17 HEAD -- cmd host docs` sums to **1,299** (cap 1,300). Per milestone: M0 77, M12 1,039, M3t 1,299 — matches the plan. The `why.go` rename is `233 285` and `why_input_test.go` is `3 3`, as §9.2 requires.
- `git diff --name-status` shows only §5 files plus the two plan files. Nothing under tools/, world/, packages/, scripts/, world-publish, host/verbs, daemon, surfacepage, or any `.ail`. No `t.Parallel`, no `exec.Command` added. `context.Background` appears only in `_test.go` files; no production context root added.
- Purity: `TestWhyImportsNoEffect` passes (direct imports exclude os, os/exec, net, net/http, io/fs, database/sql and host/store, with a broker control).
- Deduction: 1,299 against a 1,300 cap leaves no slack. The WhyFlags cut to V2b (V2a-D2) is honest and coherent: V2a ships a working seam and V2b picks up the flags.

## Mutation testing (my own; ledger `$HOME/.ailang/state/world-iter397/judge/mutants.tsv`, 28 lines, all my ids; `comm` vs the plan's ledger ids: 0 overlap, so every one of my ids is absent from the plan's table)
Each run was whole `./host/why/` first, then a narrow cmd `-run`, and a WHOLE-package run before any SURVIVED verdict.

| id | edit | verdict (killing test) |
|---|---|---|
| CTL-DIE (must-die) | ElideAt 160→159 | KILLED TestElideBounds |
| CTL-LIVE (must-live) | package doc comment | SURVIVED (correct) |
| E-NF-OP1 | drop errScanExhausted operand | KILLED TestWhyRunReplaysTheCLI (+ScanBound, golden) |
| E-NF-OP2 | drop DeadlineExceeded operand | KILLED TestWhyRunGetterClasses |
| E-HEAD-404 | 404 → present=true | KILLED TestHeadAndGetJSONClasses |
| E-GETJSON-404 | 404 → true | KILLED TestWhyRunReplaysTheCLI, RunGetterClasses (+golden) |
| E-JSON-NEWLINE | Fprintln→Fprint | KILLED by golden only (see below) |
| E-JSON-NEWLINE-noGolden | same, `-skip TestWhyCharacterisedPaths` | SURVIVED |
| E-PT-MCP | drop mcp: operand | KILLED TestParseTargetForms |
| E-PT-ZERO | n<0→n<=0 | KILLED TestParseTargetForms |
| E-RUN-SWAP | swap Stdout/Stderr | KILLED TestDemoSeedWhyNamedLinks, TestWhyCharacterisedPaths |
| E-RUN-INTEG-EXIT | ✗ exit 3→1 | KILLED TestWhyProtectedLine/e2e-tampered, TestWhyRunReplaysTheCLI |
| E-RUN-DEADLINE-EXIT | not-found exit 1→3 | KILLED TestWhyRunReplaysTheCLI, RunGetterClasses |
| E-CLI-NOSTDERR | drop stderr print | KILLED TestWhyRunReplaysTheCLI, TestWhyScanBound (+golden) |
| E-CLI-NOSTDOUT | drop stdout print | KILLED TestDemoSeedWhyNamedLinks, TestWhyCharacterisedPaths |
| E-CLI-DASH | `-` → `--` | KILLED TestWhyTargetFormsResolveToOneEntry, TestWhyBrokenLinkIsIntegrityRefusal |
| E-CLI-RESULT-UNWRAP | drop errors.Unwrap | KILLED by golden only |
| E-CLI-RESULT-UNWRAP-noGolden | same, golden skipped | SURVIVED |
| E-GETTER-METHOD | GET→HEAD in clientGetter | KILLED TestDemoSeedWhyNamedLinks, TestEpisodeReportFollowReprintsOnChangeLoopback |
| E-MC-ELIDE | `elide` delegation returns s | KILLED TestWhyHeadIsOneReader (a static AST test, but the same mutant also fails behavioural elide users) |
| E-LOC-HEADMSG | "no world head…" text | SURVIVED whole pkg |
| E-LOC-REST-OP2 | drop `ResultRef == ""` operand | SURVIVED whole pkg |
| E-LOC-FALLBACK-OP2 | drop `Plan == ""` operand | SURVIVED whole pkg |
| E-ELIDE-UTF8 | drop ToValidUTF8 | SURVIVED whole pkg |
| E-SCANFOR-NOHEAD | errScanExhausted→other error when no head | SURVIVED whole pkg |
| E-CLI-SCAN-OP1 | drop `*scan < 1` operand | SURVIVED whole pkg |
| E-SHIM-HEALTHCACHE | reader shim drops HealthCache | SURVIVED (equivalent: plan's C-HEALTHCACHE) |
| E-VH-REGROW | headState re-inlined as closure still calling why.Head | SURVIVED (equivalent: still one reader; a true re-grown switch copy was not tried) |

Plan §4 sample (the planner's own ledger, not re-run by me except where overlapping): the plan's rows W-ELIDE-AT/COUNT, W-EXIT-*, W-RUN-STREAMS, W-HEAD-*, W-PT-* are covered by the same tests my analogous mutants hit. I judged them on the table and my own kills, not by re-running 10 of them (time budget).

Survivor assessment. The five unlisted survivors (HEADMSG, REST-OP2, FALLBACK-OP2, ELIDE-UTF8, SCANFOR-NOHEAD, and CLI-SCAN-OP1) are all UNCHANGED base lines: branches the base's tests never pinned, and the move does not alter them. They are not neutrality regressions, but VC1 did not close them, so a later edit there is unguarded. Non-blocking, worth adding to VC1 or V2b (`--scan 0`, `why head` on an empty store, `rest:` on an unresolved receipt, a mid-rune Elide cut).

## Non-vacuity
- VC1 (M0): the golden test passes at the base and fails on E-JSON-NEWLINE, E-RUN-SWAP and E-CLI-RESULT-UNWRAP, so it is sensitive to the moved walk's output.
- VW2/VW3 (M3t): stdout/stderr swap, exit 3→1 and the deadline-exit mutants are killed by TestWhyRunReplaysTheCLI / RunGetterClasses, plus the older e2e tests.
- VW5: it kills a delegation that returns its input. It does NOT distinguish a closure that still calls `why.Head` (E-VH-REGROW), which is behaviourally equal.
- Golden-only kills (E-JSON-NEWLINE, E-CLI-RESULT-UNWRAP): with the golden skipped they survive. The golden is a byte comparison against a real base capture, so it is a legitimate behavioural pin, but no second test backs these two details (the trailing newline of `--json`, and the `--result` error text with `errors.Unwrap`). `TestWhyRunReplaysTheCLI` is partly circular here, since the CLI prints `Run`'s own output.

## Named survivors/equivalents (W-LOGPAGE, W-MAXRESP, W-SCANMAX, C-HEALTHCACHE, C-READER-CTX)
Accepted as non-blocking. They are constant pins and equivalent mutants, honestly described in the plan; my E-SHIM-HEALTHCACHE agrees with the C-HEALTHCACHE equivalence.

## Findings
BLOCKING: none.
Non-blocking:
1. Pre-existing unpinned branches listed above (`--scan 0`, `rest:` receipt state, fallback-plan branch, head-absent paths, mid-rune cut). VC1 could add them cheaply.
2. Two details are protected only by the golden (see above).
3. Size sits at 1,299/1,300, so V2b inherits no headroom for V2a follow-ups.
4. Process note from me: I started one mutant batch under `nohup` in the background by mistake (violating the directive), then killed it, ran `git checkout -- .` in my worktree and re-ran the same mutants in the foreground. No result from the aborted batch is in the ledger or this report. My worktree `git status --short` is clean; the throwaway base worktree is removed.

Artifacts: `$HOME/.ailang/state/world-iter397/judge/` (mutants.tsv, mut.sh, real.log, cmd.log, hostwhy.log, base_why.txt, head_why.txt).
