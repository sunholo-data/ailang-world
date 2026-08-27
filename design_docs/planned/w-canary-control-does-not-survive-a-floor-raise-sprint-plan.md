# Sprint plan — `P42` (queue row 42 `w-canary-control-does-not-survive-a-floor-raise`)

**Milestones**: `P42` — **ONE**, one commit, ~55 net LOC of Go test + three comment/prose edits.
**Status**: PLANNED · **NOT INFLATED** — a ≤0.2-day item stays a ≤0.2-day item.
**Design doc**: [`w-canary-control-does-not-survive-a-floor-raise.md`](w-canary-control-does-not-survive-a-floor-raise.md)
(573 lines, committed at `fc4776f`; two-round quorum, round 2 resolved under the narrow-refinement
carve-out by RUNNING the rejecting reviewer's own experiment → V24 → new queue row 48).
**The Decision is quorum-cleared. Nothing in this plan re-litigates it.**
**Base**: `sprint/w-canary-control-floor-raise` @ `fc4776f` (== `origin/dev`), tree clean, verified.
**Worktree**: `/Users/voightkampff/dev/sunholo-data/.wt-world-iter130` — a **sibling** of the repo,
never under `/tmp`.
**Planner**: mission-world iteration 130, opus sprint-planner, lane `opus fail-closed:env-pin`.
**Executor**: `codex` under `--sandbox workspace-write`, **no git write permission** (§6).
**THE CONTROLLER MAKES THE COMMIT.** The deliverable is an **UNCOMMITTED worktree diff**.
**Estimate**: **≤0.2 day** — honest for this shape; measured against the last five landed items (§7).

---

## 0. Planner's first-party verification at `fc4776f`

### 0.1 Base readings taken as GIVEN from the controller (rule 3e(a) — not re-derived)

The controller ran all of these on the pristine worktree at `fc4776f` this session. They are the
plan's base, not the doc's (the doc's Verification Log was measured one commit earlier, at
`2f727c7`).

| # | Command | Controller reading |
|---|---|---|
| B1 | `GOTOOLCHAIN=go1.26.6 go build ./...` | rc=0 |
| B2 | `GOTOOLCHAIN=go1.26.6 go vet ./host/verifygate/ ./host/store/` | rc=0 |
| B3 | `gofmt -l host/verifygate/ host/store/` | 0 lines |
| B4 | `AILANG_BIN=/tmp/ailang-v0300/ailang GOTOOLCHAIN=go1.26.6 go test ./host/verifygate/ ./host/store/ -count=1` | rc=0 |
| B5 | the two new test names, scoped | `no tests to run`, **rc=0** — the vacuity trap |
| B6 | `host/verifygate` **without** `AILANG_BIN`, package-wide | rc=1, 17 failures, **unrelated to any change** |

**No controller reading is challenged.** Each was spot-checked where the check was free, and every
spot-check reproduced (§0.2). **B6 is the reason every package-scoped `host/verifygate` gate line in
§4 carries `AILANG_BIN=/tmp/ailang-v0300/ailang` — with exactly one deliberate exception, AC4's
`env -u AILANG_BIN` line, whose whole purpose is to prove the two NEW tests do not need it. That
exception is scoped by `-run` to the two new names, so B6's 17 unrelated failures cannot enter it.**

### 0.2 What the planner re-ran anyway (cheap, and two of them are load-bearing)

| Claim | Command | Observed at `fc4776f` |
|---|---|---|
| B3 reproduces | `gofmt -l host/verifygate/ host/store/ \| wc -l` | `0` |
| **B5 reproduces exactly** | `GOTOOLCHAIN=go1.26.6 go test ./host/verifygate/ -run 'TestReproModuleFloorStaysBelowKnownBadToolchains\|TestCanaryDeclaresPositiveArmOnly' -count=1` | **rc=0**, `ok … 0.192s [no tests to run]` — the naive form is GREEN measuring nothing |
| AC4's lane is env-independent | `env -u AILANG_BIN GOTOOLCHAIN=go1.26.6 go test ./host/verifygate/ -run '^TestMiscompileInstrumentProbesPinnedToolchain$' -count=1 -v` | rc=0, 1 `=== RUN`, 1 `--- PASS` |
| the two base sha256s in the brief | `shasum -a 256 …/repro/go.mod …/run.sh` | `287cc106…` / `c9e2916c…` — **both confirmed** |
| all seven AC needle counts | `grep -c` per §4 | `GOTOOLCHAIN`=0, `stateRoot`=3, `POSITIVE ARM ONLY`=0, `LOAD-BEARING`=0, `^go 1.22$`=1, old-AC15-needle=1, new-AC15-phrase=0 — **identical to V8/V12/V13 one commit later** |
| helper `repoRoot` really exists | `grep -n 'repoRoot' host/verifygate/` | package-level `var repoRoot = findRepoRoot()` at `ail_binary_gate_test.go:27` — the doc's reuse claim holds |
| no identifier collides with `go/version` | `grep -n '\bversion\b' host/verifygate/toolchain_pin_gate_test.go` | only prose in comments and `go-version`/`goVersions` string work — **no shadowing** |
| the linked-worktree `.git` | `ls -la .git && cat .git` | `.git` is a **FILE**: `gitdir: /Users/voightkampff/dev/sunholo-data/ailang-world/.git/worktrees/-wt-world-iter130` — **outside the sandbox writable root** (§6, D4) |

---

## 1. Design-doc defects the planner found — five, none fatal, all cheap

Wanted, not discouraged. The doc's **Decision, ACs and mutation set are sound**; these are execution
hazards the doc's text would walk an executor straight into.

### D1 — **the drill's restore hashes are BASE hashes, and the sprint's own edits invalidate two of the three** · severity HIGH

The doc tells the sprint to restore every `repro/go.mod` arm to `287cc106…` — AC3 ("restore
`287cc106…` byte-identical (V10)"), M2b ("Restore `287cc106…` byte-identical, porcelain clean, per
arm") and the Non-Vacuity closing paragraph all say it. **That is the hash of the file BEFORE this
sprint adds the `LOAD-BEARING` comment block.** On the post-sprint tree the file's hash is
necessarily different, so an executor that restores "to `287cc106…`" either (a) fails the assertion
on a correct restore and reports a phantom red, or (b) *achieves* it by reverting §(b)'s comment —
**silently deleting the sprint's own AC7 deliverable inside the drill that was supposed to prove the
sprint correct.** The same defect applies to `host/store/toolchain_canary_test.go`
(base `5c08cc80a5cf8e33082caf91565fea9a0637d73c44c12ee660ab409682cc99b6`, changed by §(a)) for
M4/M5/M7.

**Binding reading, and it is a numbered step: `P42.7` captures the POST-SPRINT sha256 of both files
after the edits land and BEFORE the first mutation; every restore asserts against THAT value.**
`run.sh` is the one exception — **the sprint does not touch it**, so M3/M6 legitimately restore to
the doc's `c9e2916ca8ed5ce46c797e232e4fa6d19ee36e516e9ee76cc5a1f86bdea84779`, re-confirmed by the
planner at `fc4776f`. Recording that asymmetry explicitly is the point: one of the three hashes is
still valid, two are not.

### D2 — **`git status --porcelain` can never be `0` during this drill** · severity HIGH

The doc's house recipe says "porcelain 0 after every arm" (V1, V10, AC3, M2b). That was true for the
**designer**, who mutated a committed tree. It is false for the **sprint**, whose entire deliverable
is an *uncommitted* diff (§6): porcelain will legitimately show four modified paths for the whole
drill. An executor holding the doc's literal rule reports a red on a perfect restore.

**Binding reading:** `P42.7` snapshots `git status --porcelain` to
`/tmp/p42-porcelain-baseline.txt`, and every arm's post-restore check is
`diff <(git status --porcelain) /tmp/p42-porcelain-baseline.txt` → **empty**. Same guarantee,
correct instrument. (This also catches M7's rename leaking an untracked file.)

### D3 — **`host/verifygate` contains a scanner that reads every `.go` file in the package, and the doc's Conflict Surface never names it** · severity MEDIUM

`TestNoRigAbsolutePaths` (`host/verifygate/ail_binary_gate_test.go:553`) globs
`host/verifygate/*.go` and `t.Errorf`s on any file containing `/Users/`, `/tmp/ailang` or
`/home/runner/` — needles it deliberately **assembles from fragments** so it does not match its own
source. Two consequences:

1. **The new tests must resolve every path from `repoRoot`, never from a literal.** The sandboxed
   executor's most natural error — pasting the worktree path
   `/Users/voightkampff/dev/sunholo-data/.wt-world-iter130/design_docs/verification/...` into a
   `filepath.Join`, a fixture, or even a `t.Fatalf` message — reds a test **the design doc never
   mentions**, in the same package, for a reason unrelated to the item. This mission has red CI on
   exactly that defect once already (the test's own doc comment records it).
2. **AC1's `-run` filter cannot see it.** `TestNoRigAbsolutePaths` is not in the two-name pattern, so
   the scoped gates stay green while it reds. **Gate G3 (package-wide, with `AILANG_BIN`) is the
   only gate that catches it, and is therefore mandatory, not optional.**

Free bonus: it is an additional known-positive control on the new file — it must keep passing.

### D4 — **M7 says `git mv`; the executor cannot run `git mv`** · severity MEDIUM

M7 is specified as `git mv host/store/toolchain_canary_test.go host/store/toolchain_canary_moved_test.go`.
`git mv` writes the **index**, which in this linked worktree lives at
`/Users/voightkampff/dev/sunholo-data/ailang-world/.git/worktrees/-wt-world-iter130/` — measured
(§0.2), **outside `--sandbox workspace-write`'s writable root**. This is the same sandbox fact that
blocks `git commit`, and it is not an executor failure.

**Binding reading: M7 uses plain `mv`, both directions.** The kill is identical — Test B's
`os.ReadFile` on the named path fatals — and the restore is `mv` back, verified by D1's hash and D2's
porcelain diff.

### D5 — **the AC15 row's base text is longer than the doc's quotation of it** · severity LOW

The doc's Design Freeze and V13 quote the row as *"run the committed canary under deny-listed
`go1.26.5` | repro prints `BUG…`"*. The **actual** line 727 at `fc4776f` is:

```
| AC15 `MUT-CANARY-BLIND` | run the committed canary under deny-listed `go1.26.5` | repro prints `BUG: Field="" want "stateRoot"` (re-run 2026-08-25) — proving the detector still SEES the defect class, so `OK` under go1.26.6 is a measurement |
```

An executor doing a literal find-and-replace on the doc's quoted string finds nothing. **Binding
reading: replace the WHOLE of line 727 — the entire three-cell table row, leading `|` to trailing
`|` — with §(c)'s block.** AC6's two greps then read 0 and 1 exactly. Confirmed: the row is still at
**`:727`** at `fc4776f` (the design commit added no lines to `w-mcp-projection.md`), and the base
needle counts are still 1 / 0.

---

## 2. Milestone `P42` — ordered steps

**One milestone. One commit. The order is the deliverable**, because AC2 requires two separately
recorded runs and a plan that lets the executor apply everything at once destroys the evidence.

### `P42.1` — Test A and Test B, and **nothing else**

`host/verifygate/toolchain_pin_gate_test.go`, appended after
`TestMiscompileInstrumentProbesPinnedToolchain`. Add `"go/version"` to the existing import block
(**first in the stdlib group** — `go/version` sorts before `os`; gate G5 confirms). Reuse
`moduleGoFloor`, `shellAssignmentValues` and the package-level `repoRoot`; **redeclare none of
them** (measured: no name collision, V16, re-confirmed §0.2). Neither test calls `requirePinned` or
reads `AILANG_BIN` — that split is what AC4 buys.

**Test A — `TestReproModuleFloorStaysBelowKnownBadToolchains`**, the doc's Decision §A steps 1–4
verbatim: read the repro floor via `moduleGoFloor` → `version.IsValid` or instrument-failure fatal →
exactly one `KNOWN_BAD` assignment with non-empty fields → every token `version.IsValid` or
instrument-failure fatal → **`version.Compare(reproFloor, oldest) <= 0`**. Never string ordering
(V11). The failure message names the consequence and cites the V10 rehearsal, and — per V18, which
ran this test verbatim as a prototype — reads:

```
repro module floor "go1.26.1" is above the oldest KNOWN_BAD toolchain "go1.26.0": every deny-listed probe SKIPs, saw_bad stays 0, and run.sh reds for the wrong reason (the V10 rehearsal)
```

**Test B — `TestCanaryDeclaresPositiveArmOnly`**, the doc's Decision §B steps 1–3 verbatim:
known-positive control `strings.Count(src, "stateRoot") >= 2` FIRST (base 3) → zero-needle
`strings.Count(src, "GOTOOLCHAIN") == 0` (base 0) → marker `strings.Contains(src, "POSITIVE ARM
ONLY")` (base 0 — this is AC2's lever). Target resolved as
`filepath.Join(repoRoot, "host", "store", "toolchain_canary_test.go")` — **never a literal path**
(D3).

**Do not touch any comment, any `.md`, or `repro/go.mod` in this step.**

### `P42.2` — **RECORDED RUN (i): red-before-green.** Capture the transcript.

```bash
AILANG_BIN=/tmp/ailang-v0300/ailang GOTOOLCHAIN=go1.26.6 \
  go test ./host/verifygate/ \
  -run 'TestReproModuleFloorStaysBelowKnownBadToolchains|TestCanaryDeclaresPositiveArmOnly' \
  -count=1 -v
```

**Required observation — this is AC2's first leg and it is EXPECTED to be rc=1:**
- rc=**1**
- exactly **2** lines matching `^=== RUN`
- exactly **1** `--- PASS` — `TestReproModuleFloorStaysBelowKnownBadToolchains`
- exactly **1** `--- FAIL` — `TestCanaryDeclaresPositiveArmOnly`, **on the `POSITIVE ARM ONLY`
  marker clause**, not on the `stateRoot` control and not on the `GOTOOLCHAIN` zero-needle

**Do not "fix" this red.** A Test A red here means the test misreads a conformant tree (V18's green
control); a Test B red on the *control* clause means the canary's assertions have moved and the
whole fence is aimed at the wrong file — both are **stop-and-report**, not proceed.

### `P42.3` — production edit (a): the canary doc comment

Replace `host/store/toolchain_canary_test.go:5-7` with the doc's `POSITIVE ARM ONLY` block,
verbatim. **Zero assertion changes** — lines 8–47 are untouched. The block deliberately never spells
the token `GOTOOLCHAIN` (Test B's zero-needle counts it) and never adds a `stateRoot` occurrence
(the control stays at exactly 3).

### `P42.4` — production edit (b): the `repro/go.mod` comment

Insert the doc's `LOAD-BEARING` comment block **above** the `go 1.22` directive in
`design_docs/verification/w-race-gate-blindspot/repro/go.mod`. **The `go 1.22` line itself does not
move and is not re-indented** — AC7's second grep (`^go 1.22$` → 1) is what proves it. Every comment
line begins `//`, so `moduleGoFloor`'s `HasPrefix(line, "go ")` scan still finds **exactly one**
floor line; a comment line that began with `go ` would fatal both Test A and the sibling.

### `P42.5` — production edit (c): the `AC15 MUT-CANARY-BLIND` row

Replace **the whole of line 727** of `design_docs/planned/w-mcp-projection.md` (D5) with §(c)'s
single-line row, byte-for-byte.

### `P42.6` — **RECORDED RUN (ii): both green.** Capture the transcript.

Same command as `P42.2`. Required: rc=**0**, exactly **2** `=== RUN`, exactly **2** `--- PASS`,
**0** `--- FAIL`. Together with `P42.2`, this is AC2 discharged as **two separate recorded runs**.

### `P42.7` — freeze the drill's instruments (D1, D2)

```bash
shasum -a 256 design_docs/verification/w-race-gate-blindspot/repro/go.mod \
              host/store/toolchain_canary_test.go \
              design_docs/verification/w-race-gate-blindspot/run.sh \
  | tee /tmp/p42-postsprint.sha
git status --porcelain > /tmp/p42-porcelain-baseline.txt
cat /tmp/p42-porcelain-baseline.txt   # expect exactly 4 modified paths, no untracked
```

The `run.sh` line **must** print `c9e2916ca8ed5ce46c797e232e4fa6d19ee36e516e9ee76cc5a1f86bdea84779`
— if it does not, the sprint has touched a file it must not touch, and that is a stop-and-report.

### `P42.8` — the acceptance gate list (§4), in order

### `P42.9` — the mutation drill (§5), three batches, controls before and after each

### `P42.10` — hand off the **uncommitted** diff

`git diff --stat` must show **exactly four files** and no others (§3). Do not commit, do not push,
do not create a branch, do not touch the charter.

---

## 3. The exact file touch set — reconciled, and it is **NOT short**

The controller asked for this check because queue row 43 exists precisely because a floor-raise's
real touch set turned out to be **six** files against its doc's **two** (iteration 127, `P6.V`). The
planner therefore re-derived the closure by command rather than trusting the doc's table.

| # | File | Change | Steps |
|---|---|---|---|
| 1 | `host/verifygate/toolchain_pin_gate_test.go` | **MODIFY** · +~55 LOC (Test A, Test B) + `"go/version"` import | `P42.1` |
| 2 | `host/store/toolchain_canary_test.go` | **MODIFY** · doc comment `:5-7` replaced; **zero assertion changes** | `P42.3` |
| 3 | `design_docs/verification/w-race-gate-blindspot/repro/go.mod` | **MODIFY** · comment block only; directive byte-unchanged | `P42.4` |
| 4 | `design_docs/planned/w-mcp-projection.md` | **MODIFY** · line 727 replaced | `P42.5` |

**VERDICT: the doc's `Files to Create/Modify` table is COMPLETE at four. It is not short.** The
closure was searched, not assumed — here is what was checked and why each candidate is out:

- **`scripts/verify_go.sh`'s named-manifest gate does NOT cover this package.** Read end to end:
  `check_evidence_manifest` pins `REQUIRED_EVIDENCE_TESTS` with `EXACT_EVIDENCE_TESTS = 37`, and its
  discovery filter is `e.get("Package","").endswith("/host/evidence")` (`:79`). **Adding two tests to
  `host/verifygate` moves no pinned count anywhere.** This was the single most likely row-43-shaped
  surprise, and it does not fire. Repo-wide: `grep -rn 'EXACT_\|REQUIRED_' scripts/` returns hits in
  only `verify_go.sh` (evidence-only) and `verify_ail.sh` (`.ail` identities/tests — this sprint
  writes no `.ail`).
- **Nothing anywhere binds `w-mcp-projection.md`.** `git grep -n 'w-mcp-projection' -- '*.go' '*.sh'
  '*.yml'` → **zero hits**. (Contrast `P6.V`, where `docs/SELF_MOD_PUBLISH.md` was pinned by
  `host/runbook`'s AC28 — that is the mechanism that made row 43's touch set explode, and it is
  absent here.)
- **Nothing anywhere binds `toolchain_canary_test.go` by path except the new Test B.**
  `git grep -n 'toolchain_canary' -- '*.go' '*.sh' '*.yml'` → **zero hits** at base.
- **Nothing binds `repro/go.mod`** — V7's finding, re-confirmed: the only mention of the repro path
  in executable-adjacent files is the canary's own header comment, which `P42.3` rewrites.
- **`ci.yml` is untouched**, and `TestGoToolchainPinsAgreeAndMatchJobList` reads only `ci.yml` +
  root `go.mod`, neither of which moves.
- **`tools/launchd/*` is FROZEN CORE** and is not touched under any circumstances. `verify_go.sh`'s
  driver-drift gate must stay quiet; if it reds, that means *the fleet must commit*, never "absorb
  it into this change".
- **The charter (`design_docs/world-mission.md`) is controller-owned.** Marking row 42 LANDED and
  filing row 48 are **record** operations. Note for the record: at `fc4776f` the charter's queue ends
  at **row 47** — **row 48 does not exist yet**. The doc's Deferred Scope forward-references it. The
  executor must not create it.
- **Moving the doc + this plan from `planned/` to `implemented/`** is a controller record step
  (precedent: `2f727c7`), not sprint work.

---

## 4. Acceptance gates the executor runs

Run in this order. Every `host/verifygate` line carries `AILANG_BIN=/tmp/ailang-v0300/ailang`
because of base reading **B6** — *except* **G2**, whose entire purpose is the opposite and which is
`-run`-scoped to the two new names so B6's unrelated failures cannot enter it.

```bash
export PATH=/opt/homebrew/bin:$PATH
export AILANG_BIN=/tmp/ailang-v0300/ailang
export GOTOOLCHAIN=go1.26.6
NEW='TestReproModuleFloorStaysBelowKnownBadToolchains|TestCanaryDeclaresPositiveArmOnly'
```

| Gate | AC | Command | Required |
|---|---|---|---|
| **G0** | — | `GOTOOLCHAIN=go1.26.6 go build ./...` | rc=0 (base B1: rc=0) |
| **G1** | AC1 | `AILANG_BIN=… GOTOOLCHAIN=go1.26.6 go test ./host/verifygate/ -run "$NEW" -count=1 -v` | rc=0 **AND** `grep -c '^=== RUN'` = **2** **AND** `grep -c '^--- PASS'` = **2** **AND** `grep -c '^--- FAIL'` = **0** |
| **G1b** | AC1 | same, `-run TestNoSuchReproFloorTestZZZ` | prints `[no tests to run]` — the instrument SAYS SO rather than passing silently |
| **G2** | AC4 | `env -u AILANG_BIN GOTOOLCHAIN=go1.26.6 go test ./host/verifygate/ -run "$NEW" -count=1 -v` | rc=0 with G1's **same 2/2/0 counts** — no network, no solver, no pinned binary |
| **G3** | AC5 + D3 | `AILANG_BIN=… GOTOOLCHAIN=go1.26.6 go test ./host/verifygate/ ./host/store/ -count=1` | rc=0 (base B4: rc=0). **MANDATORY** — the only gate that sees `TestNoRigAbsolutePaths` and the sibling |
| **G4** | AC3 | `AILANG_BIN=… GOTOOLCHAIN=go1.26.6 go test ./host/verifygate/ -run '^TestMiscompileInstrumentProbesPinnedToolchain$' -count=1 -v` | rc=0, 1 `--- PASS` — the sibling's post-sprint re-confirmation |
| **G5** | AC5 | `GOTOOLCHAIN=go1.26.6 go vet ./host/verifygate/ ./host/store/` | rc=0 (base B2: rc=0) |
| **G6** | AC5 | `gofmt -l host/verifygate/ host/store/` | **0 lines** (base B3: 0) |
| **G7** | AC6 | `grep -c "run the committed canary under deny-listed" design_docs/planned/w-mcp-projection.md` | **0** (base 1 — red at base) |
| **G8** | AC6 | `grep -c "the repro PRINTS its verdict" design_docs/planned/w-mcp-projection.md` | **1** (base 0 — red at base) |
| **G9** | AC7 | `grep -c "LOAD-BEARING" …/repro/go.mod` | **1** (base 0 — red at base) |
| **G10** | AC7 | `grep -c '^go 1.22$' …/repro/go.mod` | **1** (base 1 — **green at base by design**: the pair proves the comment arrived *without* the directive moving) |
| **G11** | scope | `git diff --stat` | **exactly the four files of §3**, nothing else |
| **G12** | — | `AILANG_BIN=… ./scripts/verify_go.sh` | **`UNINFORMATIVE UNDER SANDBOX`** — `go test ./...` sweeps `host/daemon`, `host/broker`, `cmd/ailang-worldd`, which bind loopback. **Controller re-runs outside the sandbox; that re-run is the verdict.** |
| **G13** | — | `./scripts/verify_ail.sh` | Controller only. This sprint writes **no `.ail`**; the gate must be unchanged-green |

**Why G1's counting form and not a bare rc check.** Base reading **B5**, reproduced by the planner:
the verbatim `-run` on the two new names is **rc=0 with `[no tests to run]`** on a tree where neither
test exists. A gate that reads only the exit code is green while measuring nothing. **The
`=== RUN`/`--- PASS` counts are the gate; rc is a secondary clause.** The same applies to G2 — copy
the counting form, do not simplify it to `rc=0`.

Explicitly **rejected** as an acceptance gate: *"the full verify gate is green"* on its own. A
package-wide `ok` can print while the named tests never ran.

---

## 5. The mutation drill — 7 RED arms + 1 GREEN control

**Production side is mutated; the test helpers are never mutated.** Assertion coverage as the doc
maps it: A1-floor-read ← M1/M2a (+M7's shape), A2-single-assignment/non-empty ← M6, A3-validity ←
M3, A4-bound-compare(`<=`) ← M1/M2a **with M2b as A4's equality GREEN control**, B1-positive-control
← M7, B2-zero-needle ← M4, B3-marker ← M5.

> ### ⚠️ **M2b is a GREEN control, not a kill.**
> A run that reports M2b RED has **not** found a bug — it has found a broken plan or a
> mis-transcribed edit. At `repro` floor `go 1.26.0`, exactly equal to the oldest `KNOWN_BAD` token,
> the `<=` bound is **satisfied**, and V19 additionally measured that the instrument is *genuinely
> armed* there (`GOTOOLCHAIN=go1.26.0 go run .` → rc=0, `BUG: Field="" want "stateRoot"`, no SKIP).
> **A plan or a report that expects a red at M2b is wrong.** M2a and M2b are a boundary *pair*: M2a
> proves `<=` is enforced at all; M2b proves it is not accidentally `<`.

### 5.1 The house restore recipe (read before the first arm)

```bash
# --- BEFORE every batch: pristine known-positive control -------------------
AILANG_BIN=/tmp/ailang-v0300/ailang GOTOOLCHAIN=go1.26.6 \
  go test ./host/verifygate/ -run "$NEW" -count=1 -v      # rc=0, 2 RUN, 2 PASS

# --- per arm ---------------------------------------------------------------
cp <TARGET> /tmp/p42-$(basename <TARGET>).bak     # cp backup — NEVER `git checkout --`
<mutate>
<assert the mutation LANDED, by grep>              # a mutation that did not apply is a false green
<run the named command; record rc and the exact failure text>
cp /tmp/p42-$(basename <TARGET>).bak <TARGET>      # restore by cp
shasum -a 256 <TARGET>                             # == the P42.7 POST-SPRINT hash (D1)
diff <(git status --porcelain) /tmp/p42-porcelain-baseline.txt   # EMPTY (D2)

# --- AFTER every batch: the same pristine control again --------------------
```

**Never `git checkout -- <file>`.** In a sprint worktree that deletes the milestone's own work. This
is a standing house rule and it is doubly binding here, where the milestone's work *is* the
uncommitted diff.

**The three restore targets and their correct hashes (D1):**

| File | Restore-to hash | Source |
|---|---|---|
| `…/repro/go.mod` | **the `P42.7` post-sprint value** — **NOT** `287cc106…` | re-derived; `287cc106…` is the pre-comment base |
| `host/store/toolchain_canary_test.go` | **the `P42.7` post-sprint value** — **NOT** `5c08cc80…` | re-derived; `5c08cc80…` is the pre-comment base |
| `…/run.sh` | `c9e2916ca8ed5ce46c797e232e4fa6d19ee36e516e9ee76cc5a1f86bdea84779` | **unchanged by this sprint** — planner-confirmed at `fc4776f` |

### 5.2 Batch 1 — `repro/go.mod` (M1, M2a, M2b)

Target `T1 = design_docs/verification/w-race-gate-blindspot/repro/go.mod`.
Test command `C_A` = `AILANG_BIN=/tmp/ailang-v0300/ailang GOTOOLCHAIN=go1.26.6 go test
./host/verifygate/ -run '^TestReproModuleFloorStaysBelowKnownBadToolchains$' -count=1 -v`.

| Arm | Exact edit | Command | Expected exact result |
|---|---|---|---|
| **M1** | `sed -i '' 's/^go 1\.22$/go 1.26.6/' T1` · assert `grep -c '^go 1.26.6$' T1` = 1 | `C_A` | **rc=1, `--- FAIL`**, message names floor `"go1.26.6"` above oldest KNOWN_BAD `"go1.26.0"` and cites the V10 rehearsal. *The threat-shaped arm*: V10 measured that the **runtime** lane garbles this same edit into `INSTRUMENT FAILURE (or GOOD NEWS)`; the static lane must now **name** it |
| **M2a** | `sed -i '' 's/^go 1\.22$/go 1.26.1/' T1` · assert `grep -c '^go 1.26.1$' T1` = 1 | `C_A` | **rc=1, `--- FAIL`**, message names `go1.26.1` above `go1.26.0`. Boundary-shaped: proves `<=` is enforced at all rather than a tautology. **RUN at design time on the verbatim prototype (V18)** — reproduce, do not re-derive |
| **M2b** | `sed -i '' 's/^go 1\.22$/go 1.26.0/' T1` · assert `grep -c '^go 1.26.0$' T1` = 1 | (i) `C_A` · (ii) `cd design_docs/verification/w-race-gate-blindspot/repro && GOTOOLCHAIN=go1.26.0 go run .` | **(i) GREEN: rc=0, `--- PASS`.** (ii) rc=0, stdout `BUG: Field="" want "stateRoot"`, **no SKIP** — armed at equality (V19). **See the M2b box above: a red here is a plan/transcription error, not a finding** |

> **Sandbox note on M2b(ii) only.** `GOTOOLCHAIN=go1.26.0 go run .` needs `go1.26.0` in the local
> toolchain cache; a sandbox with no network that reports a *download* failure has produced
> **`UNINFORMATIVE UNDER SANDBOX`**, not a red. Report it as such and let the controller re-run.
> Every other arm in this drill is a pure static text scan and is fully informative in-sandbox.

### 5.3 Batch 2 — `run.sh` (M3, M6) · **the sprint does not edit this file**

Target `T2 = design_docs/verification/w-race-gate-blindspot/run.sh`, restore hash `c9e2916c…`.
`C_S` = the sibling command, `-run '^TestMiscompileInstrumentProbesPinnedToolchain$' … -v`.

| Arm | Exact edit | Command | Expected exact result |
|---|---|---|---|
| **M3** | `sed -i '' 's/go1\.26\.5"/go1.26.5x"/' T2` · assert `grep -c 'go1.26.5x' T2` = 1 | `C_A`, then `C_S` | `C_A`: **rc=1**, `t.Fatalf("instrument failure: KNOWN_BAD token \"go1.26.5x\" is not a valid Go version; version.Compare would misorder it")`. `C_S`: **rc=0, `--- PASS` — GREEN, and that is correct** (measured at design time, V20/C5; the sprint's job is to RE-CONFIRM it on the post-sprint tree, per AC3) |
| **M6** | `sed -i '' '/^KNOWN_BAD=/d' T2` · assert `grep -c '^KNOWN_BAD=' T2` = 0 | `C_A`, then `C_S` | `C_A`: **rc=1**, assignment count `0`, want `1` — fatal. `C_S`: **rc=1, `--- FAIL` with THREE clauses** — `:200` `does not contain known-positive control "KNOWN_BAD="`, `:220` `KNOWN_BAD assignment count=0, want 1`, `:236` `KNOWN_BAD must contain at least one toolchain` (measured, V21). **Shared fate is the expected shape here, not a defect** — both tests consume the same line |

### 5.4 Batch 3 — the canary file (M4, M5, M7)

Target `T3 = host/store/toolchain_canary_test.go`.
`C_B` = `AILANG_BIN=/tmp/ailang-v0300/ailang GOTOOLCHAIN=go1.26.6 go test ./host/verifygate/ -run
'^TestCanaryDeclaresPositiveArmOnly$' -count=1 -v`.

| Arm | Exact edit | Command | Expected exact result |
|---|---|---|---|
| **M4** | append **one comment line** carrying the token: `printf '// probe: GOTOOLCHAIN=go1.26.5\n' >> T3` · assert `grep -c 'GOTOOLCHAIN' T3` = 1. *A comment, deliberately: it keeps `host/store` compiling so the kill is attributable to the fence and to nothing else* | `C_B` | **rc=1**, zero-needle `1 ≠ 0`, message points the reader at the nested repro module. **ADDITION-shaped**: the fence fires on the re-add's first token |
| **M5** | delete the `POSITIVE ARM ONLY` block from `T3`'s doc comment (leave the rest of the comment and every assertion intact) · assert `grep -c 'POSITIVE ARM ONLY' T3` = 0 **and** `grep -c 'stateRoot' T3` = 3 | `C_B` | **rc=1**, marker needle absent. The `stateRoot=3` assertion above is what proves the kill is the **marker** clause and not a collapsed positive control |
| **M7** | `mv T3 host/store/toolchain_canary_moved_test.go` — **plain `mv`, NOT `git mv`** (D4) | `C_B` | **rc=1**, `os.ReadFile` fatal (`no such file or directory`) on the named path. Placement-shaped: a moved canary must move the fence in the same edit. Restore with `mv` back, then D1's hash **and** D2's porcelain diff (a stray untracked `…_moved_test.go` shows up in the porcelain diff, which is exactly what that check is for) |

### 5.5 Not mutated, and declared

The `repro/go.mod` comment (AC7's greps are the check) and the `AC15` row (AC6's greps are the
check) are prose that no test binds. They are **declared residual, not silent gaps** — the doc's
"What the gate CANNOT catch" section carries both, and this plan does not manufacture ceremony
needles for them.

---

## 6. Sandbox, worktree, and what "done" means

**Three binding sandbox facts. These travel into the executor's directive verbatim.**

1. **`--sandbox workspace-write` denies loopback binds.** Any gate whose result the executor cannot
   trust is labelled **`UNINFORMATIVE UNDER SANDBOX`** — **never** reported as a pass and **never**
   as a fail. **The controller re-runs those gates outside the sandbox before any verdict, and that
   re-run is the verdict.** For this item the labelled set is small and precisely known: **G12**
   (`verify_go.sh`, because it does `go test ./...` and sweeps `host/daemon`, `host/broker` and
   `cmd/ailang-worldd`), **G13** (`verify_ail.sh`), and **M2b(ii)** if it fails on a toolchain
   *download*. `host/verifygate` and `host/store` bind nothing, so G0–G11 and every other drill arm
   are **fully informative in-sandbox**.
2. **The executor CANNOT `git commit`, and cannot `git mv` either.** Measured: `.git` in this linked
   worktree is a **file** pointing at
   `/Users/voightkampff/dev/sunholo-data/ailang-world/.git/worktrees/-wt-world-iter130` — outside
   the writable root. Any operation that writes the index fails. **That is a sandbox fact, not an
   executor failure.** Consequences: M7 uses plain `mv` (D4), and **the deliverable is an
   UNCOMMITTED worktree diff that the controller reads and commits.** The executor leaves the four
   modified files in place, hands over the two AC2 transcripts and the drill log, and stops.
3. **Restores are `cp` from a `/tmp` backup taken before the mutation, verified by sha256** — never
   `git checkout -- <file>`, which in a sprint worktree deletes the milestone's work, and here would
   delete the *entire* deliverable.

**Worktree**: `/Users/voightkampff/dev/sunholo-data/.wt-world-iter130`, a sibling of the repo.
**Never under `/tmp`** — a `/tmp`-rooted checkout reds tests for its *location* rather than its code.

**Pinned binary**: `AILANG_BIN=/tmp/ailang-v0300/ailang`, `AILANG v0.30.0` — never a `-dirty` dev
build, never from memory. This sprint writes **no `.ail`**, so the binary is exercised only as
`host/verifygate`'s package-wide precondition (base reading B6).

---

## 7. Risks and velocity

### 7.1 Velocity — measured from the last five landed items

| Landed | Item | Doc estimate | Source diff (sprint-plan doc excluded) |
|---|---|---|---|
| `8e3c8cd` (iter-129) | `P41` row 41 | **≤0.2 d** | 2 files, **282** (+), 2 (−) — `toolchain_pin_gate_test.go` 270, `run.sh` 12/2 |
| `699f592` (iter-127) | `P6.V` row 5 | ~0.3 d | 6 files, **102** (+), 6 (−) |
| `8b196c3` (iter-126) | `P6.T` | ~0.1 d | 2 files, **5** (+), 5 (−) |
| `3dda87e` (iter-123) | `WB.K` item 14 | — | 1 file (drill-only milestone) |
| `b0d973c` (iter-122) | `WB.J` item 14 | — | 1 file, **124** (+) |

**Observed band for a ≤0.2 d item on this mission: 5–282 source insertions across 1–3 files, one
milestone PR per iteration.** This plan is **~55 LOC of Go test + three comment/prose edits across 4
files** — **at the small end of the band, and structurally identical to `P41`**, which is the item
that *built the very file this one extends* and shipped at ≤0.2 d with 270 insertions into it.

**Verdict: ≤0.2 day is HONEST for this shape, and the plan deliberately does not inflate it.** One
milestone, one commit. There is no second milestone to invent: the two tests share a file, the three
comment edits are each a handful of lines, and splitting AC2's ordering across two PRs would land a
knowingly-red `dev` — forbidden.

**The real cost is the drill, not the code.** Eight arms × (backup → mutate → assert-landed → run →
restore → hash → porcelain-diff) plus a pristine control before and after each of three batches. All
scoped `go test` runs are sub-second (measured: 0.118–0.192 s); the package-wide legs are the slow
ones (`host/verifygate` package-wide ≈ 47 s, V20). **That still fits ≤0.2 d comfortably**, which is
why the plan spends its budget on drill precision rather than on code.

### 7.2 Risks

| # | Risk | Mitigation |
|---|---|---|
| **R1** | **The executor restores to the doc's base hashes and silently reverts the sprint's own comments** (D1) | `P42.7` freezes post-sprint hashes *before* the first mutation; §5.1 names the correct target per file and flags `run.sh` as the one legitimate `c9e2916c…` |
| **R2** | **The executor reports a red on a perfect restore because porcelain is not 0** (D2) | `diff` against the `P42.7` baseline, not `wc -l == 0` |
| **R3** | **A rig-absolute path in the new tests reds `TestNoRigAbsolutePaths`** — a test the design doc never names (D3) | Every path from `repoRoot`; **G3 is mandatory** because the scoped gates cannot see it |
| **R4** | **`git mv` fails under the sandbox and M7 is reported as "blocked"** (D4) | M7 is plain `mv`, both directions |
| **R5** | **The AC15 find-and-replace finds nothing** because the doc quotes a truncation (D5) | Replace the whole of line 727; verify with G7/G8's two greps |
| **R6** | **The executor "fixes" `P42.2`'s expected red** and destroys AC2's evidence | `P42.2` states rc=1 is REQUIRED, and names the exact clause that must fail |
| **R7** | **M2b is reported as a failed kill** | The M2b box in §5, plus an explicit line in the arm table |
| **R8** | A comment line in `repro/go.mod` beginning `go ` would fatal `moduleGoFloor` and both tests | `P42.4` requires every inserted line to begin `//`; G10 (`^go 1.22$` → 1) and G1 together prove it |
| **R9** | Scope creep into `run.sh` or `ci.yml` — the tempting "just reword `INSTRUMENT FAILURE (or GOOD NEWS)`" | G11's four-file `git diff --stat`; §8 names the owner (row 44) |
| **R10** | A sandbox-load or download artefact read as a verdict | §6's `UNINFORMATIVE UNDER SANDBOX` label; controller re-run is the verdict |

---

## 8. What this sprint explicitly does **NOT** do

Absorbing any of the following is a **scope violation**, not initiative.

- **Row 44 `w-miscompile-instrument-inert-in-ci`** — `ci.yml:172`'s `continue-on-error: true`
  discarding `run.sh`'s exit code, and the script's darwin-only red in a linux CI lane. **Nothing in
  this sprint edits `run.sh` or `ci.yml`.** Rewording the `INSTRUMENT FAILURE (or GOOD NEWS)` block
  to distinguish a floor mismatch from toolchain unavailability is a `run.sh` edit and belongs
  there. Named, not absorbed.
- **Row 45 `w-toolchain-pin-normalizer-accepts-malformed-gotoolchain`** — `normalizeToolchainPin`
  is **reused unchanged**; not hardened, not touched.
- **Row 46 `w-worldd-cli-stderr-buffer-race`** — unrelated surface, untouched.
- **Row 48 `w-racecontrol-floor-bump-disarms-the-race-control` (NEW, and not yet in the charter)** —
  measured and **REFUTED** in the doc's V24: bumping `racecontrol/go.mod` to the root floor makes
  `GOTOOLCHAIN=auto` silently switch toolchains, and with the single variable `GOTOOLCHAIN=local`
  the same bump yields `go: go.mod requires go >= 1.26.6 (running go 1.26.4; GOTOOLCHAIN=local)` and
  **zero** `WARNING: DATA RACE` lines, so `verify_go.sh:232` would FATAL *"the race detector is not
  armed"*. It is the same class on a **second** module. **It is deferred deliberately**, because
  `racecontrol` has no `KNOWN_BAD` list to stay below and its real requirement — *"buildable by
  whatever toolchain `verify_go.sh` happens to run"* — is not statically knowable, so binding it
  here would be controller-invented design. **`racecontrol/go.mod` is not read, not written, and not
  mutated by any arm of this drill. This item does NOT claim systemic completeness.**
- **Row 43 `w-floor-raise-coupling-inventory`** — will *cite* the binding this sprint enforces; the
  inventory itself is not built here.
- **Row 39/40** — unrelated; row 40 is `BLOCKED on row 39`.
- **The in-module known-bad arm is NOT restored.** It is structurally unsatisfiable after `P6.T`
  (V4) — that is the finding, not a fix.
- **The root floor is not lowered or carved out.** The floor is the defect's fix, not its bug.
- **No side is derived from the other.** `repro/go.mod`'s directive and `run.sh`'s `KNOWN_BAD` stay
  independently authored artifacts; Test A only *compares* them (row 43's evaluator refutation,
  binding).
- **OD-1** (bind the oldest `KNOWN_GOOD` too) and **OD-2** (pin `KNOWN_BAD`'s exact token set) —
  both defaults are **NO**; do not implement either.
- **`tools/launchd/*` is frozen core.** Not touched under any circumstances.
- **No `.ail` is written or changed.** No charter edit. No commit. No push.

---

## 9. Open questions for the human

**None.** Deliberately empty — a question a measurement can answer is not a question. Every candidate
was measured instead: whether `verify_go.sh`'s named-manifest gate covers `host/verifygate`
(**read it — no, it filters on `/host/evidence` with `EXACT=37`**); whether anything binds
`w-mcp-projection.md` or `toolchain_canary_test.go` (**grepped with a firing control — zero hits**);
whether the AC15 row is still at `:727` after the design commit (**yes, and its base text is longer
than the doc quotes**); whether `repoRoot` exists as the doc claims (**yes,
`ail_binary_gate_test.go:27`**); whether `go/version` collides with anything in the package
(**no**); whether the linked worktree's `.git` is a file (**yes — which is what rules out `git mv`
in M7**); and whether the vacuity trap reproduces at this base (**ran it — rc=0,
`[no tests to run]`**).

The doc's Decision is quorum-cleared and its Open Decisions carry controller defaults. **There is
nothing here a human needs to rule on.**

---

**SPRINT_PLAN_PATH**: `design_docs/planned/w-canary-control-does-not-survive-a-floor-raise-sprint-plan.md`
**SPRINT_JSON_PATH**: `.ailang/state/sprints/w-canary-control-does-not-survive-a-floor-raise.plan.json`
