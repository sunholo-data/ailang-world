# World iteration 303 — M7b independent evaluator recovery report

**EVALUATION_RESULT: pass**
**EVALUATION_SCORE: 93/100**
**EVALUATION_ROUND1**
**blocking: 0**

## Identity (verified)

| Field | Value |
|---|---|
| Reviewed SHA | `f615a9698370694664fbbc01431e4a2c549377cc` |
| Exact base | `c73c79048daf9c3a19ea05c96e793436eb45e9c5` |
| Merge-base | `c73c79048daf9c3a19ea05c96e793436eb45e9c5` (= exact base) |
| Worktree | `/Users/voightkampff/dev/sunholo-data/.eval-world-iter303-recovery` |
| Files diff | 7, +1432 / -22 |
| `git diff --exit-code` rc | 0 |
| 7/7 tracked SHAs original==final | yes (`evaluator-recovery-evidence/final-shas.json`) |
| Tracked mutations only in worktree | `evaluator-recovery-evidence/` (untracked); `.bin/` build cache in worktree |

Handshake: `session_protocol_ack` returned `acked=true`. No `ailang messages` run
(unattended mission stage per directive). This is independent evaluation under
existing approved design (`design_docs/planned/w-marketing-capture.md`) and M7b
sprint plan/JSON. No git writes; no mint/publish; no permanent source/test edits
or design/charter modifications. Mutation writes only via edit-script copies and
restored from `git show HEAD:...` exact bytes; no `.go` backup files in repo.

## Scope discipline

Diff inventory (seven files, all inside scope): 3 verifier additions
(`scripts/marketingcheck/{main.go,png.go,png_test.go}`), supervisor_test package
rename, two census allow-list entries (`host/store/context_roots_test.go`),
M7b completion bookkeeping (`sprint_w-marketing-capture.json`,
`w-marketing-capture-sprint-plan.md`). JSON Unicode escape normalization was
scoped within M7b JSON; decoded strings unchanged. No unexpected file/path/mode.

Out-of-scope manual checks: `tools/launchd/*` not touched; CLAUDE.md unchanged;
`scripts/marketingcheck/png.go` is the only png-side production change; the
inherited supervisor tests (Python/shell harness) were not re-executed.

## Refusal inventory coverage

Each production hunk in `scripts/marketingcheck/png.go` was audited for refusal
behaviour. Eight refusal sites map onto eight named mutations; the M7b criteria
enumerate exactly these eight assertions:

| Production hunk (png.go) | Refusal | Mutation drill | Assertion |
|---|---|---|---|
| `:25 FileCap = 1_048_576` | `file-cap-exceeded` | M7b-cap | `png_test.go:288` cap+1 |
| `:78-80 trailing bytes` | `trailing-bytes` | M7b-tail | `png_test.go:257` |
| `:97 ihdr-not-first`/`:99 ihdr-length` | typed | (cap/CRC/fmt fixtures) | covered by structure suite |
| `:100 crc32.ChecksumIEEE` | `bad-crc` | M7b-crc | `png_test.go:235` |
| `:103 forbiddenChunks[typ]` | `forbidden-chunk` | M7b-forbidden (with :106) | `png_test.go:239` tEXt |
| `:106 !allowedChunks[typ]` | `unknown-chunk` | (retained survivor) | refusal alone ≠ bypass proof |
| `:132 total > DirCap` | `dir-cap-exceeded` | M7b-total | `png_test.go:303` |
| `:246 too-few-colors` + `:286 blank-region` | `too-few-colors`/`blank-region` | M7b-uniform (both) | `png_test.go:332/368/379/396/454` |
| `:345 !present[s.Name]` (CheckPass) | `missing-shot` | M7b-missing | `png_test.go:583` |
| `:378 a[s.Name] != b[s.Name]` (CheckTwoPass) | `pass-mismatch` | M7b-altered | `png_test.go:622` ancillary+pixel, CLI mismatch |

## M7b acceptance criteria evaluation

Sprint-plan criteria (C8/C9 PNGStructure, C8 PNGContent, C8/C9 TwoPassInventory)
all PASS in evaluator tree. Each criterion was independently run with the named
command exactly as written in `w-marketing-capture-sprint-plan.md` line 938-942.

| Criterion | Test | Eval-treerc | Eval-treecommand | Plan note |
|---|---|---|---|---|
| C8/C9 | `TestMarketingPNGStructure` | 0 | `go test -run '^TestMarketingPNGStructure$' ./scripts/marketingcheck -count=1 -timeout 90s -v` | 0.354s, all 26 subtests PASS |
| C8 | `TestMarketingPNGContent` | 0 | `go test -run '^TestMarketingPNGContent$' ./scripts/marketingcheck -count=1 -timeout 90s -v` | 0.410s, all 44 subtests PASS |
| C8/C9 | `TestMarketingTwoPassInventory` | 0 | `go test -run '^TestMarketingTwoPassInventory$' ./scripts/marketingcheck -count=1 -timeout 90s -v` | 1.265s, all 42 subtests PASS |
| controller-fix | `TestProductionGoSurface` | 0 | `go test -run '^TestProductionGoSurface$' ./host/store -count=1 -v` | 0.336s |

Logs: `evaluator-recovery-evidence/baseline-TestMarketing*.log`.

## Mutation verification (8/8 KILLED)

Each recipe imported only as an edit; bytes SHA-proved before/after; build rc=0
in `evaluator-recovery-evidence/.bin/<ID>-bin`; named test rc=1 with the
sprint-plan-asserted error code; restored SHA identical to original; originalgreen
rc=0. No panic/timeout/build error/skip/empty selection in any red or green run.

| ID | Changed SHA | Build | Test rc | Assertion hit | Restored SHA | originalgreen |
|---|---|---|---|---|---|---|
| M7b-uniform | `428f9d5e…` | 0 | 1 | `png_test.go:332` too-few-colors | `ea9b43df…` | 0 |
| M7b-missing | `e99e0700…` | 0 | 1 | `png_test.go:583` missing-shot | `ea9b43df…` | 0 |
| M7b-altered | `81ca1a00…` | 0 | 1 | `png_test.go:622` pass-mismatch (ancillary, pixel, CLI) | `ea9b43df…` | 0 |
| M7b-forbidden | `6aa7d277…` | 0 | 1 | `png_test.go:239` forbidden-tEXt (zTXt/iTXt/eXIf still pass) | `ea9b43df…` | 0 |
| M7b-crc | `ba5ae08e…` | 0 | 1 | `png_test.go:235` bad-crc | `ea9b43df…` | 0 |
| M7b-tail | `a85426cb…` | 0 | 1 | `png_test.go:257` trailing-bytes | `ea9b43df…` | 0 |
| M7b-cap | `ece52300…` | 0 | 1 | `png_test.go:288` file-cap-exceeded | `ea9b43df…` | 0 |
| M7b-total | `cd2b9431…` | 0 | 1 | `png_test.go:303` dir-cap-exceeded | `ea9b43df…` | 0 |

All eight changed SHAs match the executor's reported SHAs to 16 hex chars
(`executor-recovery-evidence/results.json` lines 165-1108). All eight assertion
texts match the sprint plan and JSON. Logs: `evaluator-recovery-evidence/M7b-*.log`.

## Reach limits and honest findings (no overclaim)

- **M7b-uniform red stops at direct CheckColors assertion** (`png_test.go:332`).
  It does NOT freshly prove CheckContent downstream paths under mutation: the
  test bails at the first refusal. The other neutered branch (`:286`) was also
  independently exercised via the blank/low-variance/color-boundaries/variance-
  boundaries subtests, all of which killed correctly (5 subtest kills total).
- **M7b-missing red stops at CheckPass** (`png_test.go:583`). It does NOT
  freshly prove the CheckTwoPass path under this specific neuter; the missing
  shot is detected and refused inside CheckPass, before the cross-pass loop.
- **M7b-forbidden bypassed both** `forbiddenChunks[typ]` (`:103`) AND
  `!allowedChunks[typ]` (`:106`) for tEXt. Only tEXt subtest kills; zTXt/iTXt/
  eXIf still pass (unaffected because the bypass is tEXt-specific). This
  satisfies the directive's "actual tEXt guard bypass matters" — the typed
  `forbidden-chunk` refusal is genuinely exercised; it is not a generic
  `unknown-chunk` discrimination.
- **Historical M7b-truncated-chunk panic survivor** (per executor ledger and
  design audit E11) is retained as a residual per the skill rule
  ("truncated-chunk panic is not an assertion kill"). Not freshly drilled.
- **Unknown-chunk refusal alone** is retained as a limit (does not prove
  forbidden-entry bypass). The M7b-forbidden drill above is the actual
  forbidden-entry bypass evidence.
- **Redundant-condition masking and filesystem-race survivors** retained
  per skill rule; not freshly drilled.
- **No socket denials encountered** in selected runs (sandbox `Operation not
  permitted` for `/tmp/eval-mut-bin-*` write redirected to in-worktree
  `evaluator-recovery-evidence/.bin/`, so no UNINFORMATIVE UNDER SANDBOX
  verdicts arose).

## Full-gate disposition (honest)

The directive says: "controller outside-sandbox full gates complete before this
judge. Do NOT rerun whole suites. ... A socket denial is UNINFORMATIVE UNDER
SANDBOX."

Controller gates (`/tmp/world-mc-20261023/controller/gates.json`, run rc
recorded; outside-sandbox): pins rc=0, focused rc=0, ail rc=0 (27.4s pinned
v0.41), go rc=1 (448.0s). PR361 CI (`pr361-before-record.json`, snapshot
before this evaluator's run) reports the headRefOid matching
`f615a9698370694664fbbc01431e4a2c549377cc` with `mergeStateStatus: CLEAN`,
`mergeable: MERGEABLE`, and three CI checks all SUCCESS: `ailang-code verify
gate`, `go host build + test gate`, `go race gate (go test -race, own job)`.

The controller's local go rc=1 is decomposed:

1. **Two unchanged-package reds** (`red-controls/candidate.log` and
   `base.log`): `TestCheckSourceUnderPinnedInterpreter` (host/archive) and
   `TestTransitionsVerbRefusesImportBearingSourceUnderPinnedInterpreter`
   (cmd/world-publish). The pristine exact-base control (`base.log`) shows
   BOTH identical failures at base with the same "want LDR001" diagnostic.
   `red-controls/results.json` records `candidate rc=1 2.4s` and `base rc=1
   2.8s` — they are pre-existing instrument/environment failures that
   predate M7b. Neither host/archive nor cmd/world-publish appears in the
   seven-file M7b diff. Per planner-recovery-plan and skill: park as
   findings, do not waive. Not M7b regressions.
2. **Marketing fake-timeout flake** in full `go test ./...` (rc=1 due to
   `TestMarketingSupervisorFakeChromeTimeoutKillsGroup` / supervisor
   group-recording variance). `red-controls/marketing-results.json` shows
   `marketing-candidate rc=0 17.9s`, `marketing-base rc=0 8.6s`,
   `marketing-race rc=0 35.5s` — all three isolated arms PASS plain AND
   race, including the suspect test. Timing-sensitive flake from
   whole-suite resource contention, not a regression introduced by M7b.

Honest disposition: M7b-specific code passes in evaluator tree, in focused
controller arms, and in PR361 CI. The two unchanged-package reds and one
timing flake are not M7b regressions and do not block M7b's own acceptance
bar. The directive's instruction "Do not call fullsuite green" is honoured:
this evaluator does not certify `go test ./...` for unrelated packages.

## Scoring (out of 100, rubric at `ailang/.claude/skills/sprint-evaluator/resources/scoring_rubric.md`)

| Category | Max | Score | Notes |
|---|---|---|---|
| Tests Pass | 20 | 20 | All 4 named M7b tests PASS in eval tree + PR361 CI |
| Lint Clean | 10 | 10 | `go vet ./scripts/marketingcheck/` rc=0; controller `ail` rc=0 |
| Acceptance Criteria | 30 | 30 | 4/4 criteria met (PNGStructure, PNGContent, TwoPassInventory, ProductionGoSurface) |
| Code Quality | 15 | 8 | png_test.go 804 lines (-5, file exceeds 800); CheckStructure 52 lines (-2, fn exceeds 50); no TODO/HACK/FIXME; sprint JSON complete with status/acceptance_criteria; M7b status "completed" in plan |
| Documentation | 15 | 15 | Code comments at non-obvious points (decoder-independence, provisional thresholds, RGB-pixel early stop); sprint plan M7b section is detailed with assertions + commands; design doc reflects M7b implementation |
| Design Fidelity | 10 | 10 | CRC fixture direct parser test present; CLI surface matches design (`check --pass-a --pass-b --manifest --rectangles [--out]`); structural parser separated from image decoder and six-set gate for non-masked mutation controls |
| Regression Surface | n/a | n/a | Not triggered (no parser/lexer/AST touched) |
| Performance | n/a | n/a | Not a perf sprint |

**Total: 93/100** (PASS; threshold 70).

### Concrete deductions

- `-5`: `scripts/marketingcheck/png_test.go` is 804 lines (>800). This is the
  M7b-added test file. The 800-line rule was applied as written; the file's
  bulk is intentional subtest coverage (~70 subtests across three top-level
  tests), but the rubric does not distinguish source vs test.
- `-2`: `scripts/marketingcheck/png.go:65 CheckStructure` is 52 lines
  (>50). One function slightly over the threshold; chunk-stream parsing is
  unavoidably linear in chunk size.

### Reach-limit notes (already priced into "honest findings"; do not double-deduct)

- M7b-uniform first assertion kill at `png_test.go:332`; later CheckContent
  paths under the same mutation NOT freshly exercised. Five independent
  subtest kills (uniform, blank, low-variance, color-boundaries,
  variance-boundaries) still confirm both neutered branches are guarded.
- M7b-missing kills at CheckPass assertion (`:583`); CheckTwoPass path not
  exercised under this neuter. The clause "Missing shot and altered valid
  PNG ... refuse specifically equality" is split across two mutations:
  missing = CheckPass branch, altered = CheckTwoPass branch. Altered
  independently kills both ancillary/pixel pairs and CLI mismatch.
- Historical M7b-truncated-chunk panic residual not re-drilled (per skill).

## Provider / cost (this evaluator run)

Provider: openrouter via pi session; explicit model label per session start.
This evaluator run performed: 4 baseline `go test` invocations, 8 mutation
cycles × (1 build + 1 test red + 1 SHA check + 1 git show restore + 1
test green) ≈ 32 sub-commands, plus SHA diffs, plus go vet. Wall time:
~5 minutes. No `ailang messages` send; no `op` secrets; no Mint.
(Supervisor appends actual NDJSON receipts after completion per directive.)

## Residual / non-blocking findings

- Pre-existing controller full `go test ./...` rc=1 with three failures: two
  unchanged-package reds (host/archive, cmd/world-publish) reproducible at
  pristine base; one marketing fake-timeout flake (isolated plain+race
  green). NOT M7b regressions; not blocking M7b. Park per planner.
- M7b-truncated-chunk panic survivor (historical, retained per skill rule).
- Unknown-chunk refusal alone doesn't prove forbidden-entry bypass
  (skill limitation; actual bypass proven by M7b-forbidden tEXt drill).
- Filesystem race and redundant-condition survivor categories (skill
  limitations; not freshly drilled).

## Report path

`/Users/voightkampff/dev/sunholo-data/.eval-world-iter303-recovery/evaluator-recovery-report.md`

## Supervisor transport and billing receipt

Canonical wrapper terminal rc0, typed verdict `ok`, fenced=true, pi_rc0; actual elapsed446s. Actual NDJSON provider reports run115assistant calls, input189639/output28591/cacheRead6135424, cost$0.45932634. Probe input1228/output17/cacheRead128, cost$0.00039648. Combined evaluator cost$0.45972282. Receipt `/tmp/world-mc-20261023/evaluator-303/supervisor-final-receipt.json`. Supervisor independently rechecked HEAD exact, tracked diff rc0 and7/7 changed tracked SHA256 equal HEAD after judge exit. Report's approximate5minute duration is superseded by446second canonical receipt.
