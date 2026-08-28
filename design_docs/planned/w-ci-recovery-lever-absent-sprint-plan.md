# Sprint plan — `w-ci-recovery-lever-absent` (queue row 47)

**Design doc**: `design_docs/planned/w-ci-recovery-lever-absent.md` (PLANNED, quorum-closed,
committed at `24e0ffc`)
**Planner**: mission-control iteration 136, `claude-opus-5` sprint-planner
**Date**: 2026-08-28
**Base commit**: `24e0ffc` (= `3a98c79` + the row-47 design-doc commit), tree clean
(`git status --porcelain` = **0** lines, planner-measured at planning entry AND exit)
**Milestones**: **2** (MS1, MS2). **Do not add a third.**
**Estimate**: ~0.10 day. **+1 line** in `.github/workflows/ci.yml`, **+113 lines** in one new
Go test file. Zero `.ail`, zero `tools/launchd/*`, zero `go.mod`, zero dependency change.
**Risk**: low. Every acceptance command and every mutation arm in this plan was **executed by
the planner** against a full working prototype before it was written down.

> **ID DISAMBIGUATION — READ FIRST.** The design doc uses `M1`/`M2` for its two **milestones**
> AND `M1`–`M9` for its nine **named RED mutations**. They are different things with the same
> names. In this plan the milestones are **`MS1`/`MS2`** and the mutation arms keep the doc's
> **`M1`–`M9`**. `MS1` = the doc's Milestone M1 (the lever line); `MS2` = the doc's Milestone M2
> (the gate). No renaming of mutation arms.

---

## 0. Planner method — what "measured" means in this document

The planner built a **full working prototype** of the shipped change in a throwaway copy of
this worktree at `/tmp/p47-probe` (`rsync -a --exclude='.git'`, so `findRepoRoot()`'s
`runtime.Caller`-derived `repoRoot` resolves to the probe and every mutation is isolated from
the real tree). Into it went (a) the lever line and (b) **the design doc's code sketch
extracted byte-verbatim** (`sed -n '174,286p' design_docs/planned/w-ci-recovery-lever-absent.md`).

Every mutation arm below was then **applied, proven landed by a moving count, and RUN** —
scoped and whole-package — and the resulting FAIL set enumerated from the actual output.
**No expected red set in this plan was predicted.** The real worktree was never mutated;
`git status --porcelain` was **0** before and after.

Environment (planner-measured, this session, this worktree):

| Fact | Measured |
|---|---|
| PATH | every command opens `export PATH=/opt/homebrew/bin:$PATH`; without it `go`/`gh` are rc=127 |
| shell | `zsh`. No `${PIPESTATUS[0]}` (bash-only, silently empty). Capture rc without a pipe: `cmd > /tmp/out 2>&1; rc=$?` |
| go | `go1.26.6 darwin/arm64` |
| gh | `2.92.0` |
| pinned AILANG | `/tmp/ailang-v0300/ailang` → `AILANG v0.30.0` (commit `e37b370`) |
| `ci.yml` sha256 at base | `9c16e64ca28c58f889a8be122a1b2f48c7d240f2a559fba6faf312ec5166239c` |
| `toolchain_pin_gate_test.go` sha256 at base | `fc8efffc641b19e93897d9ddd925cb5faea6b82d239f3fbf96e323c3384f3614` |
| `ail_binary_gate_test.go` sha256 at base | `76b88ee109a763ce8085266fe30e9df5fc1021a18163b31f172ba7d752353dab` |

---

## 1. THE FIRST-ORDER ENVIRONMENT FACT (D1) — carried forward, re-derived first-party

`go test ./host/verifygate/ -count=1` **with no `AILANG_BIN`** is, in this worktree at
`24e0ffc`, **rc=1 with 17 `--- FAIL` lines, every one reading `AILANG_BIN is unset`**.
The same command prefixed `AILANG_BIN=/tmp/ailang-v0300/ailang` is **rc=0, 0 FAIL, `ok … 29.0s`**.
The arms differ, so **the red is the environment, not the code**.

**Consequence, binding on every whole-package command in this sprint:** it carries the prefix,
or it is red at base and therefore vacuous. A red whose text says `AILANG_BIN is unset` is not
your diff. Do not "fix" it. Above all do not make those tests skip — they `t.Fatal` on purpose.

**AND THE COMPLEMENT, WHICH THE DESIGN DOC OVERSTATES (finding D13 below):** the doc says
"*Any command touching the gate package carries the `AILANG_BIN=…` prefix (P6)*". Measured
across all eleven probe arms: the **scoped** run
`go test ./host/verifygate/ -run TestEveryWorkflowDeclaresDispatchLever -count=1 -v` is
**deterministic UNPREFIXED** — rc=0/`--- PASS` on the clean prototype and rc=1 with the
expected message on every red arm. The new gate is a pure text scan with no subprocess and no
pinned-binary dependency. **Run the scoped arms UNPREFIXED so that lane-independence stays
visible**, exactly as the row-45 sprint did; prefix only whole-package runs.

---

## 2. Adjudication of the design doc against itself

The doc was revised twice (a designer revision + the controller's carve-out applying reviewers'
verbatim fixes, which added mutations M7–M9). Residue is present. **Where the plan and the doc
disagree, the DOC wins** — so each delta below quotes the doc verbatim and states what the
executor does, rather than silently rewriting the doc.

### D0 — THE HEADLINE: the doc's code sketch is shippable BYTE-VERBATIM. (severity: INFO, but load-bearing)

The planner extracted doc lines 174–286 into
`/tmp/p47-probe/host/verifygate/dispatch_lever_gate_test.go` with **zero edits** and measured:
`gofmt -l host/verifygate/` → **empty** (and `gofmt -d` → empty); `go vet ./host/verifygate/`
→ **rc=0**; scoped run → **rc=0, one `=== RUN`, one `--- PASS`**. It compiles, formats and
passes as written, with the sibling import set (P15) and no new imports.

**Instruction to the executor: ship the sketch verbatim.** Do not "improve" the parser, do not
rename the helper, do not re-word a message. Every AC and every mutation arm below is keyed on
the sketch's exact message strings; a re-wording silently voids the drill.

### D1 — `AILANG_BIN`. (HIGH) — see §1. Resolution: prefix on whole-package runs; scoped runs unprefixed and stated as such.

### D2 — Three stale `M1–M6` ranges left over from the M7–M9 carve-out. (MEDIUM)

The doc's mutation table has **nine** arms. Three sites still say six. Verbatim:

- Thesis, line 31: "*pins it with a static gate in the existing `host/verifygate` family so it
  cannot be deleted silently (the removal-shaped and addition-shaped mutations, **M1–M6**)*"
- Milestones, line 305: "*rehearse mutations **M1–M6** in probe worktrees (restore
  byte-identical, porcelain 0)*"
- AC2, line 326: "***AC2's teeth are M1–M6.***"

The `## Named RED mutations` table (lines 366–377) and the quorum record (line 489: "*landed as
**M7**, **M8**, **M9***") are the later, reviewed text. **Resolution: the rehearsal set is
M1–M9, nine arms.** The three `M1–M6` strings are stale range text, not a scope decision — the
doc itself added M7–M9 under a ratified carve-out and never swept the ranges.

### D3 — The mutation preamble misnames the test-file arm, twice in one sentence. (MEDIUM)

Verbatim, doc lines 363–364: "*The **two test-file arms (M5)** use `go vet ./host/verifygate/`
as their typecheck, because `go build ./...` does not compile `_test.go` (P8).*"

Measured: **M5 is a YAML-addition arm** (it creates `.github/workflows/deploy.yml`), not a
test-file arm. The **only** test-file arm in the table is **M6** (delete the anti-vacuity
floor), and there is **one** of it, not two. The rule the sentence states is right and
important; its subject is wrong. **Resolution: the `go vet ./host/verifygate/` typecheck
attaches to M6 (and to M6b, its combined form). Planner-measured on M6: `go vet` rc=0 AND
`go build ./...` rc=0 — the mutant compiles; and P8's warning stands, because a *broken* test
file would leave `go build ./...` at rc=0 while `go vet` reds.**

### D4 — MEASURED: the doc's multi-owner list is INCOMPLETE. M8 is multi-owner as literally described. (HIGH)

The doc's Conflict Surface and green-control paragraph name the multi-owner arms as **M4, M5**
(plus M6b). Verbatim, line 383: "*where the observable is genuinely multi-owned (M4/M5 also
move the sibling's `[ci.yml]` list, M6b also breaks sibling `os.ReadFile`)*".

**M8 is not on that list, and measurement says it belongs there — conditionally.** M8 is
"*Duplicate the top-level `on:` block (**a second column-0 `on:` later in the file**)*". The
planner ran both readings of "later in the file":

| M8 placement | `-skip TestEveryWorkflowDeclaresDispatchLever` inverse | Enumerated FAIL set |
|---|---|---|
| appended at **end of file** (after `jobs:`) | **rc=1, 1 FAIL** → MULTI-OWNER | `TestEveryWorkflowDeclaresDispatchLever` + `TestGoToolchainPinsAgreeAndMatchJobList` (`toolchain_pin_gate_test.go:138`: `ci.yml: enumerated jobs=[ailang-verify go-verify push], want [ailang-verify go-verify]`) |
| inserted **immediately before `jobs:`** | **rc=0, 0 FAIL** → SINGLE-OWNER | `TestEveryWorkflowDeclaresDispatchLever` only |

The mechanism: the duplicated block's `  push:` line matches the sibling's job regexp
`^  ([a-z0-9-]+):$`, and the sibling counts job lines only **after** it has seen `jobs:`. The
doc's P13 finding ("*the trigger line does NOT match `[a-z0-9-]+` (underscore excluded)*") is
**true and re-derived** — but it covers `workflow_dispatch:`, not the `push:` that M8's
duplicated block drags along. **Resolution: M8's canonical placement is BEFORE `jobs:`, which
makes it single-owner and gives it the `-skip` rc=0 inverse. The after-`jobs:` variant is
recorded above with its measured two-member red set so that an executor who plants it there
does not read the sibling's red as its own gate's.**

### D5 — MEASURED: the shipped sketch emits TWO messages on M9, contradicting the Decision prose. (MEDIUM)

Verbatim, doc line 291: "*and EXACTLY ONE attributed message per defect (the `t.Errorf` names
the file and the absent lever; **it never cascades** — P5's precedent)*".

Measured on M9 (`  workflow_dispatch: garbage`), scoped run, **two** messages:

```
dispatch_lever_gate_test.go:101: <abs>/.github/workflows/ci.yml: `workflow_dispatch:` has scalar value "garbage"; want an empty key or a mapping
dispatch_lever_gate_test.go:103: ci.yml: on-block triggers=[push pull_request] lack workflow_dispatch — a dropped push to dev is permanently unverifiable; every workflow file must declare the lever
```

This is structural: the scalar-value branch `continue`s **without** appending the key, so
`slices.Contains(triggers, "workflow_dispatch")` is then false and the outer `t.Errorf` fires
too. **The doc's own M9 row does not require a single message** — it requires only "*RED naming
`has scalar value "garbage"`*" — so the code and the mutation table agree; it is the Decision
paragraph's "never cascades" that overreaches. **Resolution: ship the code unchanged (the doc
wins, and the doc's M9 row is satisfied). AC9 asserts the scalar-value literal is present and
message count is 2, so the second message is a MEASURED expectation and not a surprise.**

### D6 — Line-number citation `:206` never appears in any output; the message is emitted at `:207`. (LOW)

P5 and the Conflict Surface cite `toolchain_pin_gate_test.go:206` for the `[ci.yml]`-exact
assertion. Read first-party: `:206` is `if !slices.Equal(workflowFiles, []string{"ci.yml"}) {`
and `:207` is the `t.Errorf`. **Go reports the `t.Errorf` line.** Measured under M4:
`toolchain_pin_gate_test.go:207: workflow files=[ci.yml deploy.yml], want exactly [ci.yml]; …`.
**Resolution: assert on the message TEXT (`want exactly [ci.yml]`), never on `:206`.**

### D7 — The helper's messages carry an ABSOLUTE path; the loop's message carries the base name. (MEDIUM — an assertion trap)

`onBlockTriggerKeys(t, m, src)` is called with `m`, the **full Glob match**, so its three
messages print e.g. `/tmp/p47-probe/.github/workflows/ci.yml`, while the loop's final
`t.Errorf` uses `filepath.Base(m)` → `ci.yml`. Both were observed.
**Resolution: every AC below asserts SUBSTRINGS (`ci.yml`, `deploy.yml`, `has scalar value`,
`no top-level`, `instrument failure: no workflow files enumerated`). Never assert a whole
message; never assert an absolute path — it is rig-specific and will not reproduce in CI.**

### D8 — AC6's `git status --porcelain` → 0 clause is FALSE BY CONSTRUCTION. (HIGH)

Verbatim, AC6, line 351: "*`git status --porcelain` → **0 lines** after all rehearsals*".

This cannot hold during the sprint: the deliverables are **uncommitted by construction** (the
executor performs no git writes), and the executor's cumulative milestone snapshots land in
`.snap/M<k>/`, which the planner measured is **NOT gitignored**
(`git check-ignore -v .snap/M1/x` → rc=1; same-call control `git check-ignore -v .ailang/state/x`
→ rc=0 matching `.gitignore:3:**/.ailang/` — the instrument fires).

**Resolution: porcelain is replaced by VENUE INTEGRITY.** After every rehearsal, the two
mutation venues must equal their base sha256 (`ci.yml` → `9c16e64c…6239c`;
`toolchain_pin_gate_test.go` → `fc8efffc…3614`), `.github/workflows/` must hold exactly
`ci.yml`, and `deploy.yml` must not exist. That is the check that actually detects an unrestored
probe. Porcelain would report a false red on a correct landing.

**Same-session good news, planner-measured:** `.snap/` is inert to every code gate — the Go tool
skips dot-prefixed directories, and so does `gofmt`. With a full `.snap/M2/host/verifygate/*.go`
copy planted in the probe: `go vet ./...` rc=0, `go build ./...` rc=0, `gofmt -l .` **empty**,
`go list ./... | grep -c snap` → **0**. Snapshots cannot red a gate.

### D9 — AC5's landed-proof produces its `0` from an ERROR, not from an empty directory. (LOW)

Verbatim, AC5: "*`mv .github/workflows /tmp/wf_backup_$$` so the Glob returns empty (LANDED:
`ls .github/workflows/ | wc -l` 1→0)*". After the `mv` the **directory does not exist**, so
`ls` writes to stderr and `wc -l` counts zero lines — the right number for the wrong reason,
and it would read `0` for a typo'd path too.
**Resolution: landed-proof becomes `test -d .github/workflows; echo $?` **0→1** (the directory
genuinely went away) paired with `ls -1 .github/workflows 2>/dev/null | wc -l` **1→0**.**

### D10 — M6 has no observable on its own; the doc's "M6b" is an unlisted arm. (MEDIUM)

The M6 row's Predicted-result cell already says the proof needs "*M6 + M6b*", and its
landed-proof cell introduces `M6b` — but `M6b` is not a row in the table. Measured, and the doc
is RIGHT that M6 alone proves nothing:

| arm | scoped rc | whole-package (prefixed) |
|---|---|---|
| **M6 alone** (floor deleted, `.github/workflows/` intact) | **rc=0, `--- PASS`** | **rc=0, 0 FAIL** — no gate anywhere in the repo moves |
| **M6b alone** (floor intact, workflows dir moved away) = AC5 | rc=1, floor message at `:89` | rc=1, **4** FAIL (enumerated in AC5 below) |
| **M6+M6b** (floor deleted AND dir moved away) | **rc=0, `--- PASS`** ← the proof | rc=1, **3** FAIL, **none of them mine** (enumerated below) |

**Resolution: M6b is promoted to an explicit arm. M6's verdict is read on the SCOPED run only
(`--- PASS` with the floor gone vs. the floor message with it present); the whole-package rc is
uninformative for M6+M6b because three siblings red on the missing `ci.yml` for reasons that
have nothing to do with the floor.**

### D11 — "names `ci.yml` and `workflow_dispatch` exactly once" is a per-arm claim, not a template. (LOW)

AC3's phrasing is **measured correct for M1** (`ci.yml` ×1, `workflow_dispatch` ×1, one
message). It does **not** generalise, and reusing it would produce dead or false assertions:

| arm | occurrences of `ci.yml` in scoped output | occurrences of `workflow_dispatch` | messages |
|---|---|---|---|
| M1 | 1 | **1** | 1 |
| M2 | 1 | **2** — `workflow_dispatcht` *contains* `workflow_dispatch` | 1 |
| M3 | 1 | 1 | 1 |
| M7 | 1 | **0** — the message never names the trigger | 1 |
| M9 | 2 | 2 | **2** |

**Resolution: each AC below carries its own measured counts. Count OCCURRENCES with
`grep -o … | wc -l`, never `grep -c` (which counts lines and would silently collapse M9's two
messages if they ever shared a line).**

### D12 — Two premise numbers have drifted since the doc was written; neither is load-bearing. (INFO)

- **P12** measured "*whole repo → 15, all in `design_docs/world-mission*.md` prose*". Re-measured
  at `24e0ffc`: **64**. The delta is the design doc's own commit (it says `workflow_dispatch`
  many times). **The load-bearing half is unchanged and re-derived: code surfaces
  (`--include='*.go' --include='*.sh' --include='*.yml' --include='*.yaml'`) → 0.** No competing
  mechanism writes the observable.
- Every other premise the plan leans on was re-derived first-party and **HOLDS**: P1 (1 file,
  `ci.yml`), P2 (`workflow_dispatch` → 0 rc=1; controls `pull_request` → 1, `^  push:` → 1),
  P3 (`on:` at `:3`, `push:` `:4`, `pull_request:` `:6`), P4 (`ail_binary_gate_test.go:668`,
  `toolchain_pin_gate_test.go:106`), P5 (Glob at `:197`), P7 (gofmt/vet/build clean), P9
  (2 `actionlint` hits, both in `host/verifygate/*_test.go` comments), P14 (`^on:$` → 1,
  control `^  pull_request:$` → 1), P15 (import set as quoted).

### D13 — See §1: the doc's blanket "every command carries the prefix" is over-broad for the scoped arms. (LOW)

### D14 — MEASURED CONFIRMATION, not a defect: the lever is inert to every existing gate. (INFO)

With the lever line landed and the new test file present, the planner ran the **whole**
`host/verifygate` package prefixed in the probe: **rc=0, 0 `--- FAIL`, `ok … 45.8s`.** The
doc's Decision claim that the addition is purely additive, and P13's claim that the trigger line
cannot corrupt the sibling's job set, are both **confirmed by execution**, not just by regex
probe.

**Nothing in the design doc was found to be false about the repo** beyond the drifted count in
D12 and the internal residue in D2/D3/D4/D5/D6/D8/D9/D10. Its Premises table, its Conflict
Surface, its Decision sketch and its nine mutation arms all describe a change that the planner
built and ran end to end.

### Does the code sketch, the ACs and M1–M9 describe the SAME parser? — the specific check asked for

| Property | Sketch | ACs | Mutations | Verdict |
|---|---|---|---|---|
| exact `l == "on:"` column-0 anchoring | line 198 `if l == "on:"` | not named in any AC | **M7** exercises it | AGREE. Gap: no AC names M7. **Repaired by AC7 below.** |
| duplicate top-level `on:` detection | lines 206–215, `t.Fatalf … want exactly 1` | not named in any AC | **M8** exercises it | AGREE. Gap: no AC names M8. **Repaired by AC8 below.** |
| scalar-value rejection for `workflow_dispatch` | lines 232–239, `has scalar value %q` | not named in any AC | **M9** exercises it | AGREE. Gap: no AC names M9. **Repaired by AC9 below.** |
| anti-vacuity floor | lines 267–270 | **AC5** | **M6/M6b** | AGREE |
| runtime enumeration (every workflow file) | line 263 `filepath.Glob` | **AC4** | **M4/M5** | AGREE |
| trigger-depth / comment handling | lines 218–230 | not named | **M2/M3** | AGREE. **Repaired by AC3b below.** |

**Conclusion: the three parsers agree.** The only defect is coverage: the doc ships six ACs
against nine mutations, so **M2, M3, M7, M8 and M9 — every construct the carve-out added — have
no acceptance criterion.** This plan adds AC3b, AC7, AC8 and AC9 to close that, each one
base-measured. This is a coverage repair inside the doc's declared scope; it adds **no code**.

---

## 3. Acceptance criteria — every one BASELINED on the pristine tree

`export PATH=/opt/homebrew/bin:$PATH` opens every command. `PREFIX` =
`AILANG_BIN=/tmp/ailang-v0300/ailang`. `SCOPED` =
`go test ./host/verifygate/ -run TestEveryWorkflowDeclaresDispatchLever -count=1 -v`
(**unprefixed** — D13). Baselines below were run by the planner on the pristine tree at
`24e0ffc` (ACs 1, 2n, 6) or on the byte-verbatim prototype (all red arms).

| AC | Command | **BASE (measured)** | REQUIRED after the sprint | Repaired? |
|---|---|---|---|---|
| **AC1** | `grep -cE '^  workflow_dispatch:$' .github/workflows/ci.yml` with same-call controls `grep -c 'pull_request'` and `grep -cE '^  push:'` | **0** (rc=1), controls **1**, **1** | **1**, controls still **1**, **1** | no — doc is correct; red at base is the finding being reversed |
| **AC2** | `SCOPED` | **N/A** — the test does not exist at HEAD. Nonsense control base measured: `go test ./host/verifygate/ -run '^TestNoSuchGateZZZ$' -count=1 -v` → **rc=0, `[no tests to run]`, 0 `=== RUN`** | rc=0, exactly **1** `=== RUN`, **1** `--- PASS`; nonsense control unchanged; `gofmt -l host/verifygate/` empty | amended: base of the control measured; prefix dropped (D13) |
| **AC3** | M1 arm, then `SCOPED` | N/A | rc=1; occurrences `ci.yml`=**1**, `workflow_dispatch`=**1**, messages=**1**; **`-skip` inverse `PREFIX go test ./host/verifygate/ -count=1 -skip TestEveryWorkflowDeclaresDispatchLever` → rc=0, 0 FAIL** (single-owner, measured) | **repaired**: occurrence counting method fixed (D11); `-skip` inverse added as the single-owner proof |
| **AC3b** *(new)* | M2 and M3 arms, then `SCOPED` | N/A | both rc=1 with the `lack workflow_dispatch` message; M2 triggers list reads `[push pull_request workflow_dispatcht]`, M3's reads `[push pull_request]`; **`-skip` inverse rc=0/0 FAIL for both** | **added** — the doc has no AC for M2/M3 |
| **AC4** | M4 arm, then `SCOPED`, then `PREFIX go test ./host/verifygate/ -count=1` | N/A | `SCOPED` rc=1 naming **`deploy.yml`**; **ENUMERATED whole-package red set, exactly 2**: `TestEveryWorkflowDeclaresDispatchLever` (mine, `deploy.yml … lack workflow_dispatch`) + `TestGoToolchainPinsAgreeAndMatchJobList` (`toolchain_pin_gate_test.go:207: workflow files=[ci.yml deploy.yml], want exactly [ci.yml]`). **No `-skip` inverse exists — it is rc=1** | **repaired**: red set enumerated by execution, both members explained; `:206`→`:207` (D6) |
| **AC5** | M6b arm (dir moved away, floor intact), then `SCOPED`, then `PREFIX` whole package | N/A | `SCOPED` rc=1 containing `instrument failure: no workflow files enumerated`; **ENUMERATED whole-package red set, exactly 4**: mine (`dispatch_lever_gate_test.go:89`) + `TestZ3PinDeclaredOnceAndInstalledInBothJobs` (`ail_binary_gate_test.go:672: open …/ci.yml: no such file or directory`) + `TestGoToolchainPinsAgreeAndMatchJobList` (`toolchain_pin_gate_test.go:110`, same cause) + `TestMiscompileInstrumentStepIsGatedInCI` (`toolchain_pin_gate_test.go:404`, same cause) | **repaired**: landed-proof fixed (D9); 4-member red set enumerated |
| **AC6** | `gofmt -l host/verifygate/`; `go vet ./host/verifygate/`; `go build ./...`; `PREFIX go test ./host/verifygate/ -count=1` | `gofmt` **empty rc=0**; `vet` **rc=0**; `build` **rc=0**; prefixed package **rc=0, 0 FAIL, ok 29.0s** (unprefixed: **rc=1, 17 FAIL, 17 `AILANG_BIN is unset`**) | all four identical; plus `grep -c 'AILANG_BIN is unset'` on the test output = **0**; plus **venue integrity** (below) | **repaired**: `git status --porcelain → 0` replaced by venue integrity (D8) |
| **AC7** *(new)* | M7 arm, then `SCOPED` | N/A | rc=1, exactly **1** message containing ``has no top-level `on:` trigger block``; occurrences of `workflow_dispatch` in output = **0**; **`-skip` inverse rc=0/0 FAIL** | **added** — no AC covered the column-0 anchoring the carve-out shipped |
| **AC8** *(new)* | M8 arm **inserted immediately before `jobs:`**, then `SCOPED` | N/A | rc=1, exactly **1** message containing ``has 2 top-level `on:` blocks, want exactly 1``; **`-skip` inverse rc=0/0 FAIL** (single-owner **at this placement only** — D4) | **added** |
| **AC9** *(new)* | M9 arm, then `SCOPED` | N/A | rc=1, **2** messages (D5): one containing `` has scalar value "garbage" ``, one containing `lack workflow_dispatch`; **`-skip` inverse rc=0/0 FAIL** | **added** |

**Venue integrity** (replaces AC6's porcelain clause, D8), asserted at every milestone boundary
and after every rehearsal:

```
shasum -a 256 .github/workflows/ci.yml host/verifygate/toolchain_pin_gate_test.go
# post-MS1 ci.yml differs from base by exactly the one added line; toolchain_pin_gate_test.go
# must still be fc8efffc641b19e93897d9ddd925cb5faea6b82d239f3fbf96e323c3384f3614
ls -1 .github/workflows/          # exactly: ci.yml
test -e .github/workflows/deploy.yml; echo $?   # 1
test -d .github/workflows;         echo $?      # 0
```

**Baseline tally: 6 doc ACs baselined, 4 repaired (AC3, AC4, AC5, AC6), 1 amended (AC2), 1
confirmed sound as written (AC1). 4 ACs added (AC3b, AC7, AC8, AC9) to cover the five
carve-out mutations the doc left uncovered. No AC in this plan is red at base.**

---

## 4. Mutation drill — blast radius CLASSIFIED BY EXECUTION

Every row was run. `-skip inverse` = `PREFIX go test ./host/verifygate/ -count=1 -skip
TestEveryWorkflowDeclaresDispatchLever`; **rc=0 there proves the arm is SINGLE-OWNED** (nothing
but my gate moved). Where it is rc≠0 the arm is MULTI-OWNER and carries an enumerated set
instead. All landed-proofs are counts that MOVE in both directions.

| # | Mutation (exact) | Landed-proof (MEASURED, both directions) | `-skip` inverse | Enumerated whole-package FAIL set (MEASURED) | Type |
|---|---|---|---|---|---|
| **M1** | delete the `  workflow_dispatch:` line from `ci.yml` | `grep -cE '^  workflow_dispatch:$'` **1→0** AND `grep -c 'workflow_dispatch'` **1→0**; control `grep -c 'pull_request'` stays **1** | **rc=0, 0 FAIL** | `{TestEveryWorkflowDeclaresDispatchLever}` — `ci.yml: on-block triggers=[push pull_request] lack workflow_dispatch …` | SINGLE-OWNER / NEW TEETH |
| **M2** | `  workflow_dispatch:` → `  workflow_dispatcht:` | `grep -cE '^  workflow_dispatch:$'` **1→0** AND `grep -cE '^  workflow_dispatcht:$'` **0→1** (one call) | **rc=0, 0 FAIL** | `{mine}` — `triggers=[push pull_request workflow_dispatcht] lack workflow_dispatch` | SINGLE-OWNER / NEW TEETH |
| **M3** | `  workflow_dispatch:` → `  # workflow_dispatch: recovery lever` | keyed `grep -cE '^  workflow_dispatch:$'` **1→0** WHILE raw `grep -c 'workflow_dispatch'` **stays 1** (a moving count *and* a held count — proves comment, not deletion) | **rc=0, 0 FAIL** | `{mine}` — `triggers=[push pull_request] lack workflow_dispatch` | SINGLE-OWNER / NEW TEETH |
| **M4** | create `.github/workflows/deploy.yml` (bytes in §5) with `on:`/`  push:` and **no** lever | `ls -1 .github/workflows \| wc -l` **1→2** AND `grep -cE '^  workflow_dispatch:$' deploy.yml` = **0** | **rc=1 — NO inverse** | **2 members**: `{mine (deploy.yml … lack workflow_dispatch), TestGoToolchainPinsAgreeAndMatchJobList (:207 workflow files=[ci.yml deploy.yml], want exactly [ci.yml])}`. The sibling member is multi-owner **by construction** (Conflict Surface): it owns the `[ci.yml]`-exact policy. Assert on **my** `deploy.yml` message text, never on rc | MULTI-OWNER / NEW TEETH |
| **M5** | same file **with** `  workflow_dispatch:` | `ls -1 … \| wc -l` **1→2** AND `grep -cE '^  workflow_dispatch:$' deploy.yml` **0→1** | n/a | `SCOPED` rc=**0**, **`--- PASS`** ← the assertion. Whole package rc=1 with exactly **1** member, `{TestGoToolchainPinsAgreeAndMatchJobList :207}` — the sibling alone, explained above | MULTI-OWNER / **MUST-STAY-GREEN** |
| **M6** | delete the `if len(matches) == 0 { t.Fatal(…) }` guard from the new test file | `grep -c 'instrument failure: no workflow files enumerated' host/verifygate/dispatch_lever_gate_test.go` **1→0** | n/a | typecheck `go vet ./host/verifygate/` **rc=0** (`go build ./...` also rc=0 — vacuous for test files, P8/D3). Whole package **rc=0, 0 FAIL** — **M6 alone has no observable**; it must be combined with M6b | TEST-FILE / needs M6b |
| **M6b** | `mv .github/workflows /tmp/p47-wfbak` (floor **intact**) — this is AC5's arm | `test -d .github/workflows; echo $?` **0→1** AND `ls -1 .github/workflows 2>/dev/null \| wc -l` **1→0** | **rc=1 — NO inverse** | **4 members** (enumerated in AC5). Assert only on **my** `instrument failure: no workflow files enumerated` at `dispatch_lever_gate_test.go:89` | MULTI-OWNER / FLOOR PROOF |
| **M6+M6b** | both together — **the arm that proves the floor is load-bearing** | both landed-proofs above | n/a | `SCOPED` rc=**0**, **`--- PASS`** ← my gate prints a checkmark on an empty enumeration once the floor is gone. Whole package rc=1 with **3** members, **none of them mine**: `{TestZ3PinDeclaredOnceAndInstalledInBothJobs :672, TestGoToolchainPinsAgreeAndMatchJobList :110, TestMiscompileInstrumentStepIsGatedInCI :404}`, all `open …/ci.yml: no such file or directory`. **Read the SCOPED result only** | COMBINED / FLOOR PROOF |
| **M7** | delete the whole top-level `on:` block; append a nested decoy step carrying `        on:` / `          workflow_dispatch:` | `grep -c '^on:$'` **1→0** AND `grep -c 'workflow_dispatch'` stays **≥1** (the decoy is still there — the pair proves *decoy planted*, not *lever deleted*) | **rc=0, 0 FAIL** | `{mine}` — ``dispatch_lever_gate_test.go:101: instrument failure: <abs path>/ci.yml has no top-level `on:` trigger block``. Occurrences of `workflow_dispatch` in output = **0** | SINGLE-OWNER / NEW TEETH |
| **M8** | insert a second column-0 `on:` block **immediately before `jobs:`** | `grep -c '^on:$'` **1→2** | **rc=0, 0 FAIL** | `{mine}` — ``:101: instrument failure: <abs>/ci.yml has 2 top-level `on:` blocks, want exactly 1``. **If instead appended after `jobs:`: rc=1, 2 members** (D4) | SINGLE-OWNER at the stated placement / NEW TEETH |
| **M9** | `  workflow_dispatch:` → `  workflow_dispatch: garbage` | `grep -cE '^  workflow_dispatch:$'` **1→0** AND `grep -c '^  workflow_dispatch: garbage$'` **0→1** (one call) | **rc=0, 0 FAIL** | `{mine}`, **2 messages** (D5): `:101 … has scalar value "garbage"; want an empty key or a mapping` and `:103 … lack workflow_dispatch` | SINGLE-OWNER / NEW TEETH |

**Reading rule.** Six arms (M1, M2, M3, M7, M8, M9) are single-owned and may use the `-skip`
rc=0 inverse as their blast-radius proof. Three (M4, M6b, M6+M6b) and the green control (M5)
are multi-owner; for those, **rc cannot distinguish which gate reded** — assert on my test's
message text and check the enumerated set member-by-member. Every set above was produced by
running the arm, not by prediction.

---

## 5. Milestones

### MS1 — the lever (0.05d, **+1 line**, `.github/workflows/ci.yml`)

Insert `  workflow_dispatch:` at trigger depth (lead 2), after `  pull_request:` at `:6`, so
the block reads exactly:

```yaml
on:
  push:
    branches: [dev]
  pull_request:
  workflow_dispatch:
```

No `inputs:` (Option B). Nothing else in `ci.yml` changes.

**Gate**: AC1 (base **0** → **1**, controls held at 1/1) + `PREFIX go test
./host/verifygate/ -count=1` → **rc=0, 0 FAIL** (planner-measured with the lever landed and no
new test present: the lever alone reds nothing) + venue integrity.

### MS2 — the gate (0.05d, **+113 lines**, new file `host/verifygate/dispatch_lever_gate_test.go`)

Ship the design doc's code sketch (doc lines 174–286) **byte-verbatim**: `package verifygate`,
the five-import block, `onBlockTriggerKeys`, and `TestEveryWorkflowDeclaresDispatchLever` with
its full doc comment (including the branch-protection / static-scan / enumerator residual
paragraph — that text IS the declared residual, do not trim it). Planner-measured: gofmt-clean,
vet-clean, passes, no import changes (P15).

**Gate**: AC2, AC3, AC3b, AC4, AC5, AC6, AC7, AC8, AC9 + venue integrity.

**Rehearsal venues** (mutation-only; each restored from a `cp` backup, never `git checkout --`,
since nothing is committed):

- `.github/workflows/ci.yml` — M1, M2, M3, M7, M8, M9
- `.github/workflows/` (the directory) — M6b
- `host/verifygate/dispatch_lever_gate_test.go` — M6 (the only test-file arm, D3)
- `.github/workflows/deploy.yml` (created, then **deleted**) — M4, M5

`deploy.yml` bytes the planner used (M4 — the M5 variant adds `  workflow_dispatch:` after
`  push:`):

```yaml
name: Deploy

on:
  push:

jobs:
  noop:
    runs-on: ubuntu-latest
    steps:
      - run: 'true'
```

---

## 6. Executor constraints

- **NO git write operations.** No `add`, `commit`, `stash`, `checkout`, `restore`, `reset`,
  `branch`, `worktree`, `git mv`. Git **reads** (`show`, `status`, `diff`, `log`,
  `check-ignore`, `ls-files`) are required and allowed. **The controller builds the commits.**
- **Restores use `cp` from a `/tmp` backup taken before the first mutation.** `git checkout --
  host/verifygate/dispatch_lever_gate_test.go` would **delete the deliverable** — the file is
  untracked, so git has nothing to restore it from.
- **Snapshots**: after MS1, copy the cumulative changed set to `.snap/MS1/`; after MS2, to
  `.snap/MS2/` (cumulative, i.e. MS2's snapshot contains both files). Measured inert to every
  gate (D8): the Go tool and `gofmt` both skip dot-prefixed directories.
- **Files this sprint MAY touch**: `.github/workflows/ci.yml` (MS1, +1 line);
  `host/verifygate/dispatch_lever_gate_test.go` (MS2, new).
- **Files this sprint MUST NOT touch**: `host/verifygate/toolchain_pin_gate_test.go` (sibling —
  its `[ci.yml]`-exact list at `:206`/`:207` is a **declared residual 4**, a future row's
  policy decision, **not** something to "fix" when M4/M5 red it);
  `host/verifygate/ail_binary_gate_test.go`; `scripts/verify_go.sh`; `scripts/verify_ail.sh`;
  `go.mod`/`go.sum`; `tools/launchd/*` (frozen core, FLEET-owned); any `.ail` file; any other
  workflow file (`deploy.yml` exists only inside a rehearsal and is deleted after).
- **Message text is FROZEN.** The four strings the drill keys on must appear byte-exactly as in
  the doc's sketch: ``has no top-level `on:` trigger block``, ``top-level `on:` blocks, want exactly 1``,
  `` `workflow_dispatch:` has scalar value %q; want an empty key or a mapping ``,
  `lack workflow_dispatch — a dropped push to dev is permanently unverifiable`, and the floor
  `instrument failure: no workflow files enumerated under .github/workflows/`.
- **Authority**: THE DESIGN DOC WINS over this plan, **except** on D1, D2, D3, D4, D5, D8, D9,
  D10 and D13, where the planner's first-party measurement showed the doc's literal text would
  either red a correct landing (D1, D8) or pass vacuously / mis-scope a result (D2, D3, D4, D5,
  D9, D10, D13). On those nine, this plan wins.
- **Do not widen any gate to `go test ./...`.** The narrowest gate that can fail for this diff
  is `go test ./host/verifygate/ -count=1`.
- **Cost centres, planner-measured**: whole `host/verifygate` package **29.0s** (base) /
  **45.8s** (probe) *with* `AILANG_BIN`; scoped run **~0.15–0.4s**; `go build ./...` ~1s;
  `go vet ./host/verifygate/` <1s. Budget the full drill at ~12 min of wall clock.

---

## 7. Sandbox posture (executor: `codex:gpt-5.6-sol`, `--sandbox workspace-write`)

`workspace-write` **denies loopback socket binds**. Nothing in this sprint's own gate needs a
socket — `TestEveryWorkflowDeclaresDispatchLever` is a pure text scan over YAML with no
subprocess and no network — but the **whole-package** run
(`AILANG_BIN=… go test ./host/verifygate/ -count=1`) drives `scripts/verify_ail.sh` through the
shim arms, and the repo's wider suites include socket-touching daemon tests.

**Rule:** any gate result obtained inside the sandbox for a socket-touching suite must be
labelled **`UNINFORMATIVE UNDER SANDBOX`** in the executor's log — reported, never silently
dropped, and never read as a pass or a fail. **The controller re-runs those gates outside the
sandbox** and the controller's numbers are the ones of record. The scoped arms (AC2, AC3, AC3b,
AC7, AC8, AC9 and the `SCOPED` half of AC4/AC5) are sandbox-safe and their results stand as
measured.

`gh` is not required by any acceptance criterion in this sprint (P16/P17 are design-time
premises, already measured by the controller; the live dispatch rehearsal is **Deferred Scope**).

---

## 8. Residuals this plan inherits unchanged

All seven of the doc's Declared residuals stand, unmodified — in particular **residual 1** (the
lever buys **a verdict on the tip of a named ref**, never a mergeable PR: a dispatch run's
checks do not satisfy branch protection) and **residual 7** (the recoverable window closes the
moment `dev` advances past the unverified commit). The gate's green means **"every enumerated
workflow file's TEXT declares the lever"** and nothing stronger. Nobody may read it as "CI can
always be re-triggered".

**SPRINT_PLAN_PATH**: `design_docs/planned/w-ci-recovery-lever-absent-sprint-plan.md`
**SPRINT_JSON_PATH**: `design_docs/planned/sprint_w-ci-recovery-lever-absent.json`
