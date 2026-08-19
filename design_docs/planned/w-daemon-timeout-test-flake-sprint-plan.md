# Sprint plan — `w-daemon-timeout-test-flake` (queue row 19)

**Item**: queue row 19, `w-daemon-timeout-test-flake` — the design said "already-expired"; the test spelled it "1 ns", and in Go "1 ns" is a timer
**Status**: PLANNED · **ONE milestone** · **0.25 days** · one commit
**Design doc**: [`w-daemon-timeout-test-flake.md`](w-daemon-timeout-test-flake.md) (706 lines; two full quorum rounds + a controller carve-out — **the doc is authoritative; where anything here and the doc disagree, the doc wins, and every place this plan departs from it is named in §0.5 with the measurement that forced it**)
**Base**: `sprint/w-daemon-timeout-test-flake` @ **`5380e6f`**, clean tree, **both gates measured GREEN first-party this session** (§0.2)
**Planner**: mission-control iteration 94, opus sprint-planner, first-party measurement on this rig
**Executor**: sandboxed worktree, **NO git write permission**. **THE CONTROLLER MAKES ALL COMMITS.** The executor never runs `git commit`, `git add`, `git push`, `git checkout`, `git stash`, `git restore`, or `gh pr`. Restores are `cp` from a backup — §4.1.

**Base conditions on EVERY command in this plan** (controller-verified, re-verified by the planner):

```bash
export AILANG_BIN=/tmp/ailang-v0300/ailang
export GOTOOLCHAIN=go1.25.6
```

`verify_go.sh` FATALs without `AILANG_BIN` and **denylists the local `go1.26.4`** (`scripts/verify_go.sh:113-121` lists go1.26.0–go1.26.5). Omitting `AILANG_BIN` makes 17 `host/verifygate` arms fail **LOUDLY BY DESIGN** — that is a base condition, **not a regression**.

---

## 0. Planner's first-party baseline — every acceptance command, run on the pristine tree

**P2 discipline**: a gate already red at base measures the repo, not the change. Every command in
§5 was run by the planner at `5380e6f` before it entered this plan, and its observed output is
recorded below. **Zero repo bytes were written to produce any of it**: the head-tree and mutant
readings of §0.4 were taken with `go build/vet/test -overlay` against scratch files in `/tmp`, and
`git status --porcelain` is empty and both target files' sha256 unchanged after the whole session
(§0.6).

### 0.1 Instrument pins

| Pin | Command | Observed |
|---|---|---|
| HEAD | `git rev-parse --short HEAD` | **`5380e6f`**, `git status --porcelain` empty |
| Toolchain | `GOTOOLCHAIN=go1.25.6 go version` | `go1.25.6 darwin/arm64`; `GOROOT` = `…/toolchain@v0.0.1-go1.25.6.darwin-arm64` (present, no download) |
| Pinned binary | `/tmp/ailang-v0300/ailang --version` | `AILANG v0.30.0, Commit e37b370` |
| Observatory | (stderr of every `ailang` call) | `Observatory: 302MB (warn threshold: 200MB)` — the iteration-89 stderr-merge repair holds; `verify_ail.sh` is rc=0 with the warning present. **Not a regression, do not "fix" it.** |
| sha256 PRE | `shasum -a 256 host/daemon/read_deadline_test.go host/daemon/handlers.go` | `59c8f03df154edd1c659f308e13402ae2a895ae8dd3c398b6c5eb4785d6cd5f4` / `cc9f206c0472287551c5885ab654efb221502ba252d5235c814621ddde94c708` |

The test file's PRE hash `59c8f03df154edd1…` **is byte-identical to the hash M4 recorded when it
restored the file** after the controller's iteration-93 probe — independent confirmation that the
controller's one-shot left no trace and that this plan and that measurement share a base.

### 0.2 Base gate state (AC1's base arm), measured

```
./scripts/verify_ail.sh   -> rc=0   (4 s)    ✓ 10 required identities verified, 39 named tests pass; world package gate 9/9
./scripts/verify_go.sh    -> rc=0   (162 s)  ✓ 0 binary blobs among 215 tracked files
                                             ✓ driver drift gate: 6 tracked driver files, working tree matches HEAD
                                             ✓ go gate PASSED: build clean, plain and race tests pass with pinned AILANG_BIN
```

**The driver drift gate is GREEN in this worktree** — `git status --porcelain -- tools/launchd/
scripts/mission_decisions.sh` is empty here, unlike the iteration-91 checkout. So unlike the
previous sprint, **`verify_go.sh` IS this sprint's gate and AC1 runs it.** Do not substitute a bare
`go build && go test`.

**AC1's honest base condition, carried forward from the doc (§8 AC1) and re-measured (§0.3 rows
AC3/AC4):** a full `verify_go.sh` runs both affected tests, so **the base gate reds on exactly these
two tests with a small probability from the very flake this item fixes.** The doc prices that at
≈2.9% (two legs); the planner's own base rates are ~1.6× the doc's, and the planner also measured
that the **race leg is effectively immune** (0 failures in 1,000 runs under `-race`), so the honest
number is **≈2.1%, concentrated in the plain leg**. A base red matching `status = 200, want 503` /
`a timed-out read answered 200` is **attributed AT BASE** (rule 3d), the gate re-run once, and the
event recorded — it is **not** a regression and **not** "discovered" as one. At head that shape is
extinct and any red is real.

### 0.3 AC-by-AC baseline table — every §5 command, run at `5380e6f`

| AC | Command (abridged; §5 has the exact form) | Base reading, OBSERVED | Notes |
|---|---|---|---|
| **AC1** | `./scripts/verify_ail.sh`; `./scripts/verify_go.sh` | **rc=0 / rc=0**, first attempt, no re-run needed | 4 s and 162 s |
| **AC2** | `go test ./host/daemon -run '^TestExpiredReadDeadlineExpiresAtConstruction$' -v -count=1` | **rc=0**, `--- PASS:` lines = **0**, `--- FAIL:` = 0, output `testing: warning: no tests to run` | **The vacuous rc=0.** This is exactly why AC2's verdict is the COUNTED `--- PASS:` line and never rc |
| **AC3** | `… -run '^TestDaemonReadDeadline$/^real-store-expired-deadline$' -count=1000 -v` | **rc=1**; sample A: **985 PASS / 15 failing runs**; sample B: **994 PASS / 6 failing runs** → **21/2000 = 1.05%**. Failure text: `read_deadline_test.go:258: status = 200, want 503; body={…real data…}` on `world`, `object`, `log range`, `registry`, `head` | The AC **can fail** at base, with ~99.998% probability at n=1000 and the measured rate. Wall clock **3 s per sample** |
| **AC4** | `… -run '^TestTimeoutStatusMirrorsSketch$' -count=2000 -v` | **rc=1**; sample A: **1974 PASS / 26 FAIL** in 2000; sample B (n=1000): **994 / 6** → **32/3000 = 1.07%** | Wall clock **6 s** (n=2000) and **3 s** (n=1000) |
| **AC5(a)** | `sed 's://.*::' host/daemon/*.go \| grep -c 'Nanosecond\|Microsecond'` | **2** | the two `= 1 * time.Nanosecond` stimulus lines (`read_deadline_test.go:254`, `:526`), which are code and survive the strip |
| **AC5(b)** | `find host cmd -name '*.go' -exec sed 's://.*::' {} + \| grep -c …` | **2**, measured at that scope (not assumed from (a)); rest-of-scope (`-not -path 'host/daemon/*'`) = **0** | |
| AC5 control | `sed 's://.*::' host/daemon/handlers.go \| grep -c readDeadline` vs raw `grep -c` | **6** vs **7** | the strip removes the one comment mention and keeps all six code hits — the instrument reads |
| AC5 scope | `find host/daemon -type d` | `host/daemon` only — **no subdirectory**, so the `*.go` glob equals the old recursive scope | |
| **AC6(a)** | `grep -c 'LIMITATION(w-daemon-late-read-503)' host/daemon/handlers.go`; control `grep -c 'writeReadTimeout' …` | **0** (rc=1) with control **7** (rc=0) | the zero carries its known-positive control in the same call |
| **AC6(b)** | `grep -c 'err.Error()' host/daemon/handlers.go` | **5** | item 18's AC5 live count; base-green by design |
| **AC6(c)** | `go build ./host/daemon/` | **rc=0** | |
| **AC7 pre** | `go vet ./host/daemon` | **rc=0** | the mutant-vets baseline |
| **AC8** | `git status --porcelain`; `shasum -a 256 <both files>` | empty; hashes in §0.1 | |
| name | `grep -rn 'TestExpiredReadDeadline' host/` | **rc=1**, no hits | the test name is unallocated |
| sites | `grep -n 'time.Nanosecond' host/daemon/read_deadline_test.go` | `254`, `526` | M1/D2 re-derived |
| seam | `grep -n 'func (d \*Daemon) readCtx' host/daemon/handlers.go`; `grep -n 'readDeadline = 10 \* time.Second' host/daemon/daemon.go` | `handlers.go:269-270`; `daemon.go:128` | D16 re-derived — **both identifiers in scope, so MU-DEADLINE-DETACH compiles** (proven in §0.4) |

### 0.4 HEAD PRE-REGISTRATION — the changed tree and all three mutants, measured before the executor starts

The planner applied §4.1/§4.2/§4.3 **verbatim** to scratch copies in `/tmp/head94/` and ran the
whole acceptance set and the whole mutation table through `go {build,vet,test} -overlay`. **The repo
was never written.** Every number below is therefore a **pre-registered value, not an expectation**:
if the executor's run reports anything else, the edit is not the one probed here — **STOP and
report**, do not adjust a pin.

**Head acceptance readings (measured):**

| AC | Head reading, PRE-REGISTERED | Wall clock |
|---|---|---|
| AC2 | rc=0, **exactly 1** `--- PASS: TestExpiredReadDeadlineExpiresAtConstruction` | 0.4 s |
| AC3 | rc=0, **1000** `--- PASS: TestDaemonReadDeadline/real-store-expired-deadline`, **0** `--- FAIL` | **3 s** |
| AC4 | rc=0, **2000** `--- PASS: TestTimeoutStatusMirrorsSketch`, **0** `--- FAIL` | **5 s** |
| AC5(a) | **1** — the matched line is `const expiredReadDeadline = -1 * time.Nanosecond` (leading `-` present) | instant |
| AC5(b) | **1** (changed package 1 + rest-of-`host/ cmd/` 0 — composed, both halves measured) | instant |
| AC5 control | stripped **6** / raw **7** on `handlers.go` — **unmoved** by the §4.3 comment | instant |
| AC5 raw secondary (**declared NOT the gate**) | **2** at base and **2** at head — cannot move, cannot fail for the right reason. Recorded, never gating | instant |
| AC6(a) | **1**, control **7** | instant |
| AC6(b) | **5** — unmoved | instant |
| AC6(c) | rc=0 | 1 s |
| full package | `go test ./host/daemon -count=1` rc=0 | 2 s |

**Mutation arms (measured on the same scratch head):**

| Mutation | vet / build | KILL arm | INVERSE arm |
|---|---|---|---|
| MU-STIM-POSITIVE | vet rc=0 | rc=1, **exactly 1** `--- FAIL: TestExpiredReadDeadlineExpiresAtConstruction` | `-skip` the killer, package `-count=1` → **rc=0** (observed) |
| MU-DEADLINE-DETACH | **build rc=0**, vet rc=0 (the mutant compiles — D16's two-identifiers-in-scope claim, executed) | rc=1, **exactly 1** `--- FAIL:`, message `readCtx under an already-expired deadline returned a LIVE context …` | rc=1 with **4** `--- FAIL:` lines — the **corrected** red set, §0.5(i) |
| MU-SITE-REVERT | vet rc=0 | **none — and the pin test still PASSES under it (measured: 1 `--- PASS:`)**, which is the row's whole point | AC5(a) reads **2** (against expected 1) → RED; AC5(b) reads **2** → RED; package `-count=1` → rc=0 green |

### 0.5 Findings that CHANGE the plan

#### (i) **PLAN-AFFECTING DOC DEFECT — MU-DEADLINE-DETACH's inverse red SET is under-counted by one. `blocking-store` reds too, and the doc says any red outside its set "fails the arm".**

Design doc §6, MU-DEADLINE-DETACH inverse arm: *"`TestDaemonReadDeadline/real-store-expired-deadline`
and `TestTimeoutStatusMirrorsSketch` red deterministically …, `normal-deadline-answers-200` and
everything else green. **Any red outside that set is unexplained and fails the arm.**"*

**Measured, twice — at the base stimulus and again at the head stimulus** (`go test -overlay … -skip
'^TestExpiredReadDeadlineExpiresAtConstruction$' -count=1 -v`):

```
--- FAIL: TestDaemonReadDeadline (2.01s)
    --- FAIL: TestDaemonReadDeadline/real-store-expired-deadline (0.00s)
    --- FAIL: TestDaemonReadDeadline/blocking-store (2.01s)          <-- NOT in the doc's set
--- FAIL: TestTimeoutStatusMirrorsSketch (0.01s)
    --- PASS: TestDaemonReadDeadline/normal-deadline-answers-200 (0.00s)
rc=1, exactly 4 `--- FAIL:` lines
```

**Why `blocking-store` reds, and why it belongs in the set**: that subtest sets
`d.readDeadline = 50 * time.Millisecond` (`read_deadline_test.go:271`) and parks a getter until
`ctx.Done()`. The mutant makes `readCtx` ignore the **field** and use the 10 s package constant, so
the park is never released inside test time and the subtest's own 2 s watchdog fires. That is the
**same phenotype** as the other two reds — "the field is ignored" — not an unexplained red. An
executor following the doc's sentence verbatim would score a **correct** mutant as a **failed arm**.

**Plan decision**: §4.2's MU-DEADLINE-DETACH row states the **corrected, measured** red set
(three subtests/tests, four `--- FAIL:` lines, `normal-deadline-answers-200` green, rc=1) and
requires the executor to **enumerate** the reds and compare to that set. **This is the only place
this plan overrides a literal sentence of the doc, it is overridden by measurement rather than by
argument, and the doc's underlying intent — "the red set is the assertion, an unexplained red fails
the arm" — is preserved exactly.** Report it to the controller for the record commit.

#### (ii) **`gofmt` reflows the doc's §4.3 comment block and one line of §4.2. Pasting them verbatim leaves the tree non-gofmt-clean.**

Both target files are gofmt-clean at base (`gofmt -l` → empty). Pasting §4.2 and §4.3 verbatim makes
**both** files dirty under `gofmt -l`. Measured (`gofmt -d`), two changes, both cosmetic:

1. **§4.3's numbered residual list**: Go 1.19+ doc-comment formatting treats the indented
   `  (i) …` / ` (ii) …` lines as a code block — gofmt inserts a blank `//` before and after the
   list and re-indents its lines with a tab. The final landed form is quoted **verbatim and
   gofmt-clean** in §3, task T1.3. **AC-neutral**: it is all comment, so AC5(a)/(b) (comment-
   stripped) and AC6(b) cannot see it, and AC6(a)'s token is untouched — all four **re-measured on
   the gofmt'd file**: 1 / 1 / 5 / 1-with-control-7.
2. **§4.2's second `t.Fatalf`**: gofmt wants `"… returned a LIVE context — the " +` (space before
   `+`); the doc writes `the "+`. The *first* `Fatalf` in the same test is left alone by gofmt.
   The gofmt-clean form is quoted in §3, task T1.2.

**No gate enforces gofmt** — `grep -rln 'gofmt' .github scripts` returns nothing, and
`verify_go.sh` has no fmt arm — so this reds nothing. It is hygiene, and the repo is currently
clean, so the plan prescribes `gofmt -w` on both files (T1.4) and **pre-registers the reflow as
EXPECTED**, so that neither the executor "fixes" it by reverting nor the evaluator reads the
whitespace delta from the doc's rendering as a mis-paste.

#### (iii) **The high-count ACs are ~50× faster on this rig than the controller's brief states. Budget from the measurement, not from the brief.**

The controller's brief states `-count=1000` of `TestTimeoutStatusMirrorsSketch` took **138 s**.
**That did not reproduce.** Four first-party samples under `GOTOOLCHAIN=go1.25.6` on this rig:

| Command | Elapsed |
|---|---|
| `-run '^TestTimeoutStatusMirrorsSketch$' -count=1000` (base) | **3 s** (`2.707s` reported by the test binary) |
| `-run '^TestTimeoutStatusMirrorsSketch$' -count=2000` (base) | **6 s** |
| same at head (overlay) `-count=2000` | **5 s** |
| `-run '^TestDaemonReadDeadline$/^real-store-expired-deadline$' -count=1000` (base ×2, head ×1) | **3 s** each |

The plan budgets from these (§5), with a generous ceiling. **The executor must not treat a 3-second
AC3 as "it didn't run"** — the verdict is the counted `--- PASS:` line (1000 of them), not the
clock. Conversely, if either AC takes minutes on the executor's rig, that is *load*, not a hang;
the honest ceiling is: if AC3 or AC4 exceeds **5 minutes**, stop and report the discrepancy rather
than killing the run.

#### (iv) **The base flake is a PLAIN-leg phenomenon: 0 failures in 1,000 runs under `-race`.**

`go test -race ./host/daemon -run '^TestTimeoutStatusMirrorsSketch$' -count=1000` → **rc=0, 1000
PASS, 0 FAIL, 18 s**. The race detector slows the store read enough that the armed timer wins every
time. Consequences, both stated because both matter:

- **AC1's base-red probability is lower than the doc's two-leg model** and is concentrated in the
  plain leg: ≈2.1%, not ≈2.9%. The direction is favourable; the doc's number is not "wrong", its
  per-leg independence assumption over-counts the race leg.
- **The race leg is NOT evidence about this flake in either direction.** A green race leg at base
  says nothing; the plain leg is where the shape lives.

#### (v) **The planner's own near-miss, recorded because it is the class this repo keeps paying for.**

Building the MU-SITE-REVERT mutant, the planner's first edit script targeted "the occurrence of
`d.readDeadline = expiredReadDeadline` after `func TestTimeoutStatusMirrorsSketch`" — and there are
**two** after that point in the head file: the mirrors-sketch stimulus **and the new §4.2 pin test's
own field write**. The script asserted `count == 1`, the assert fired, and **nothing was written** —
so the sweep that followed read the *unmutated* value **1** and would have scored as a clean arm if
the assert had not been there. **A mutation that never landed prints the same green as one that
worked.** This is why §4.1 step 3's sha256 landed-proof AND step 3b's token grep are both
mandatory, and why T-MU3 identifies its site **by enclosing function**, never by a bare pattern:
after this change the head file contains **three** `d.readDeadline = expiredReadDeadline` writes
(the two stimulus sites plus the pin test's).

### 0.6 The planner wrote zero repo bytes

After the entire baseline + pre-registration session:

```
git status --porcelain                                            -> empty
shasum -a 256 host/daemon/read_deadline_test.go host/daemon/handlers.go
   59c8f03df154edd1c659f308e13402ae2a895ae8dd3c398b6c5eb4785d6cd5f4  (unchanged)
   cc9f206c0472287551c5885ab654efb221502ba252d5235c814621ddde94c708  (unchanged)
```

Every head/mutant measurement used `-overlay` with scratch files under `/tmp/head94/`,
`/tmp/mu_stim/`, `/tmp/mu_detach/`, `/tmp/mu_revert/`. **The executor does NOT need `-overlay`** —
it edits the real files and restores from `cp` backups (§4.1). The overlay was a planner-only
technique for probing a tree it is not allowed to write.

---

## 1. Scope

### 1.1 What this sprint delivers — ONE milestone, one commit

**M1_DETERMINISTIC_STIMULUS** (~0.25 d, ~70 changed lines across 2 files, 1 commit).

1. `host/daemon/read_deadline_test.go`: one named constant `expiredReadDeadline = -1 *
   time.Nanosecond` with the doc's comment block, replacing **BOTH** `1 * time.Nanosecond` stimulus
   lines (`:254` inside `TestDaemonReadDeadline/real-store-expired-deadline` and `:526` inside
   `TestTimeoutStatusMirrorsSketch`) with `d.readDeadline = expiredReadDeadline` (doc §4.1).
2. Same file: the new mechanism pin `TestExpiredReadDeadlineExpiresAtConstruction` (doc §4.2).
3. `host/daemon/handlers.go`: the `LIMITATION(w-daemon-late-read-503)` block appended to
   `timedOut`'s doc comment — **comment-only, zero behavioural bytes** (doc §4.3).

Measured size of exactly that change (planner's scratch application, gofmt-clean):
`read_deadline_test.go` **+41 / −2**, `handlers.go` **+27 / −0**.

### 1.2 What this sprint does NOT do — from doc §11, do not "fix" any of it

- **Not** changing production behaviour. `handlers.go` gains a comment and nothing else;
  `daemon.go` is untouched. ARM B (a post-read expiry check) and ARM C are **rejected on the
  record** (doc §3) — do not implement, do not "improve toward".
- **Not** weakening, moving, or re-wording either test's assertions. The 503/`Timeout` pins are
  byte-unchanged; **only the stimulus moves.**
- **Not** adding a retry or a skip anywhere. The one operator re-measurement allowed in §4.2's
  MU-STIM-POSITIVE inverse arm binds the **operator**, not the committed gate.
- **Not** committing a `-count` loop to CI. AC3/AC4 are operator acceptance commands.
- **Not** touching rows 20/21, `host/store/context_read_test.go` (it already spells expiry
  correctly), `blockingStore`/`recordingStore`, the 50 ms stimulus at `:271`, any `.ail` file,
  `verify_*.sh`, CI workflows, or `tools/launchd/*`.
- **Not** building the AST guard against future positive sub-deadline stimuli (doc §10 declines it
  with its ≥2-instance ledger).
- **Not** filing the successor rows. `w-daemon-late-read-503` and the busy_timeout-ordering
  residual (doc §11's last paragraph) are **controller/record actions**, not sprint work.

### 1.3 FROZEN — the executor must not touch these

`tools/launchd/*` and `scripts/mission_decisions.sh` (FLEET-owned, `D-WORLD-DRIVER-1`; the drift
gate is green here and must stay green) · `scripts/verify_ail.sh`, `scripts/verify_go.sh` ·
`design_docs/sketches/worlddapi.ail` (frozen; read by `sketchHTTPStatusVectors`) ·
`design_docs/planned/w-daemon-timeout-test-flake.md` (**the design doc — do not edit it**) ·
`~/.ailang/state/mission-v1*` · every file outside `host/daemon/{read_deadline_test.go,handlers.go}`.

---

## 2. P1 — doc↔plan acceptance-criteria cross-check

**A design doc and its sprint plan are two documents describing one sprint; revising either
silently rots the other. This repo has paid for that. So the two lists are written out and
compared explicitly.**

**The ACs the DOC says this milestone closes.** Doc §4 is titled *"M1 — the change (one milestone,
independently landable)"* — there is exactly one milestone, so every AC in doc §8 is closed by it.
Doc §8 enumerates, in order:

> **AC1**, **AC2**, **AC3**, **AC4**, **AC5**, **AC6**, **AC7**, **AC8** — eight criteria.

**The ACs THIS PLAN's milestone section (§3) names as closed by `M1_DETERMINISTIC_STIMULUS`:**

> **AC1, AC2, AC3, AC4, AC5, AC6, AC7, AC8** — eight criteria.

**CONFIRMED IDENTICAL.** Same count (8), same members, no AC dropped, no AC invented, no AC
deferred to a follow-on, no AC split across milestones (there is only one milestone). The sprint
JSON's `acceptance_criteria_map` carries the same eight and no others.

**Mutation cross-check, same discipline.** Doc §6 names exactly three rows — **MU-STIM-POSITIVE**,
**MU-DEADLINE-DETACH**, **MU-SITE-REVERT**. This plan prescribes exactly those three (§4.2) and
**adds none**. Where an AC has no named mutation, the plan says so rather than inventing one:

| AC | Named mutation in doc §6 | If none, what stands in |
|---|---|---|
| AC1 | none | The gates' own non-vacuity is a **standing base condition**: `verify_go.sh` FATALs without `AILANG_BIN` and on a denylisted toolchain, and the base flake itself reds it ~2.1% of the time. AC1 demonstrably **can** fail |
| AC2 | MU-STIM-POSITIVE (sign arm) **and** MU-DEADLINE-DETACH (`ctx.Err()` arm) | — both of the doc's §4.2 "two assertions, two mutations" are real, executed arms (§0.4) |
| AC3 | none | Non-vacuity is **measured, not argued**: the AC reds at base — 15/1000 and 6/1000 observed this session (§0.3) |
| AC4 | none | Same: 26/2000 and 6/1000 observed at base |
| AC5 | MU-SITE-REVERT | the only detector, per doc §6 — see §4.2 T-MU3 |
| AC6 | none (doc §6: *"the comment has no behaviour to mutate"*) | (a) carries a same-call known-positive control (`writeReadTimeout` = 7); (b) is a **declared base-green paired control** whose only reachable failure mode is this item's own edit; (c) is a compile |
| AC7 | is itself the mutation table | — |
| AC8 | none | hygiene; its instrument is `git status --porcelain` + sha256 |

---

## 3. M1 — task breakdown with per-task verification

**Milestone `M1_DETERMINISTIC_STIMULUS` — 0.25 d — closes AC1, AC2, AC3, AC4, AC5, AC6, AC7, AC8.**

Every command below is run from the worktree root with both exports of the header.

### T1.1 — the named constant, and BOTH stimulus sites (doc §4.1) → advances AC5

**Edit `host/daemon/read_deadline_test.go`:**

(a) Insert, immediately after the import block's closing `)` and before the
`// ---------- The route table under test ----------` banner (the placement the planner probed):

```go
// expiredReadDeadline is the deterministic timeout stimulus for the read-
// deadline tests: any NON-POSITIVE duration makes context.WithTimeout take
// context.WithDeadline's `dur <= 0` branch, which cancels the context
// SYNCHRONOUSLY at construction — no timer, no goroutine, no race — so the
// store read is refused at connection acquisition and the 503/Timeout branch
// runs on every route, every run. A small POSITIVE duration (the previous
// `1 * time.Nanosecond`) is a FUTURE deadline: it arms a time.AfterFunc, and a
// fast read can complete before the timer goroutine runs, answering 200 with a
// real body (measured at base: ~0.65–0.8% of runs). Do not "shrink" this back
// to a positive value; TestExpiredReadDeadlineExpiresAtConstruction reds on
// the sign, and the design doc for w-daemon-timeout-test-flake holds the
// measurements.
const expiredReadDeadline = -1 * time.Nanosecond
```

(b) Replace **both** stimulus lines — `:254` (inside `TestDaemonReadDeadline`, subtest
`real-store-expired-deadline`) and `:526` (inside `TestTimeoutStatusMirrorsSketch`) — with:

```go
d.readDeadline = expiredReadDeadline
```

keeping each line's existing indentation (`:254` is two tabs, `:526` is one). **Do not touch
`:271`** (`d.readDeadline = 50 * time.Millisecond`, the `blocking-store` subtest) — doc §8 AC5 and
D4 declare that site deliberately outside the pattern and safe.

**Verify T1.1:**
```bash
grep -n 'time.Nanosecond' host/daemon/read_deadline_test.go       # -> ONE line: the const, with a leading `-`
grep -n 'd.readDeadline = ' host/daemon/read_deadline_test.go     # -> :254-ish and :526-ish now read `= expiredReadDeadline`;
                                                                  #    the 50ms site is UNCHANGED
GOTOOLCHAIN=go1.25.6 go vet ./host/daemon                         # -> rc=0
```

### T1.2 — the mechanism pin (doc §4.2) → closes AC2

Append to the **end of** `host/daemon/read_deadline_test.go` the test below. This is the doc's §4.2
block with the one gofmt correction of §0.5(ii) already applied (`the " +`, space before `+`) —
paste **this** form:

```go
// TestExpiredReadDeadlineExpiresAtConstruction pins the property every 503
// assertion in this file now rests on: the stimulus context is DEAD BEFORE any
// store read can begin. Two assertions, two mutations:
//   - the sign check kills the "shrink it back to a positive nanosecond"
//     mutation deterministically (no timing anywhere);
//   - the ctx.Err() check goes through the production readCtx, so a readCtx
//     that ignores d.readDeadline (or re-derives from the 10s constant) reds
//     here in one run.
func TestExpiredReadDeadlineExpiresAtConstruction(t *testing.T) {
	if expiredReadDeadline >= 0 {
		t.Fatalf("expiredReadDeadline = %s, want a negative duration — a positive value arms "+
			"a timer and re-creates the 200-vs-503 race this constant exists to remove",
			expiredReadDeadline)
	}
	d := newHandlerDaemon(t)
	d.readDeadline = expiredReadDeadline
	ctx, cancel := d.readCtx(httptest.NewRequest(http.MethodGet, "/v1/head", nil))
	defer cancel()
	if ctx.Err() == nil {
		t.Fatalf("readCtx under an already-expired deadline returned a LIVE context — the " +
			"stimulus must be expired at construction, before any store read can begin")
	}
	if !errors.Is(ctx.Err(), context.DeadlineExceeded) {
		t.Fatalf("ctx.Err() = %v, want context.DeadlineExceeded", ctx.Err())
	}
}
```

**No new imports** — `context`, `errors`, `net/http`, `net/http/httptest`, `time` are all already in
the file's import block (verified, and the scratch application compiles). `newHandlerDaemon` is the
package's own constructor (`handlers_test.go:19`). **No fake anywhere**: the asserted value comes
from the `context` package through production `readCtx`.

**Verify T1.2 (this is AC2):**
```bash
GOTOOLCHAIN=go1.25.6 go test ./host/daemon -run '^TestExpiredReadDeadlineExpiresAtConstruction$' -v -count=1 \
  > /tmp/ac2.out 2>&1; echo "rc=$?"
grep -c -- '--- PASS:' /tmp/ac2.out     # -> EXACTLY 1   (base was 0 with a VACUOUS rc=0 — never read rc alone)
```

### T1.3 — the LIMITATION comment (doc §4.3) → advances AC6

In `host/daemon/handlers.go`, insert immediately **above** `func timedOut(ctx context.Context, err
error) bool {` (line 282 at base), i.e. as the tail of `timedOut`'s existing doc comment, the block
below. This is the doc's §4.3 text **after gofmt** (§0.5(ii)) — paste this form so the tree stays
gofmt-clean:

```go
// LIMITATION(w-daemon-late-read-503): this classifier is consulted ONLY on the
// error path — on a successful read err is nil and this function never runs.
// D15 verifies prompt cancellation only for one CPU-bound query using
// modernc.org/sqlite v1.54.0 on darwin/arm64. It does not establish a general
// bound for lock-blocked or other driver waits; therefore this item makes no
// general production bounded-wait claim. The residual is two things, not one:
//
//	 (i) a read that COMPLETES SUCCESSFULLY before the cancellation is
//	     observed answers 200 with the real body even though it overran the
//	     deadline;
//	(ii) a read blocked on a LOCK is bounded by busy_timeout, NOT by the
//	     request deadline (w-daemon-timeout-test-flake D18: under a 300ms
//	     context deadline a lock-blocked read returned after 2.04s, governed
//	     by busy_timeout's ~2s, the deadline exceeded 6.8x while the
//	     busy-retry loop ran — and the error still surfaces as
//	     deadline-exceeded, so the wire class is Timeout while the timing was
//	     governed by a different mechanism entirely). Today that composition
//	     is safe only because busy_timeout (2s) is shorter than the 10s
//	     request deadline — an ORDERING nothing in this code asserts, not a
//	     guarantee.
//
// Enforcing a stronger status contract means a post-read expiry check at all
// seven read sites and a decision that completed work must be discarded; that
// is the successor item named above, not a rider on a test fix. The
// read-deadline tests deliberately use an ALREADY-EXPIRED stimulus
// (expiredReadDeadline) precisely because a near-expiry one races the
// completes-before-cancellation window.
```

**ZERO behavioural bytes.** Do not touch `timedOut`'s body, `readCtx`, `writeReadTimeout`, or any
handler. The comment is deliberately free of `err.Error()`, `Nanosecond`, `Microsecond`, and the
bare token `readDeadline` — **do not "improve" it by adding any of those**; each one moves an AC.

**Verify T1.3 (this is AC6(a)+(b)+(c)):**
```bash
grep -c 'LIMITATION(w-daemon-late-read-503)' host/daemon/handlers.go   # -> 1        (base 0)
grep -c 'writeReadTimeout' host/daemon/handlers.go                     # -> 7        (control, base 7)
grep -c 'err.Error()' host/daemon/handlers.go                          # -> 5        (base 5, MUST NOT MOVE)
GOTOOLCHAIN=go1.25.6 go build ./host/daemon/; echo "rc=$?"             # -> rc=0
sed 's://.*::' host/daemon/handlers.go | grep -c 'readDeadline'        # -> 6        (instrument control, base 6)
grep -c 'readDeadline' host/daemon/handlers.go                         # -> 7        (instrument control, base 7)
```

### T1.4 — gofmt hygiene

```bash
GOTOOLCHAIN=go1.25.6 gofmt -l host/daemon/read_deadline_test.go host/daemon/handlers.go   # -> EMPTY
```
If either file is listed, run `gofmt -w` on it and **re-run T1.3's greps** (they are all
comment/token counts and were measured unmoved by the reflow — a change in any of them means the
formatter touched code, which it must not).

### T1.5 — the static sweep (doc §8 AC5) → closes AC5

```bash
sed 's://.*::' host/daemon/*.go | grep -c 'Nanosecond\|Microsecond'                          # -> 1  (base 2)
find host cmd -name '*.go' -exec sed 's://.*::' {} + | grep -c 'Nanosecond\|Microsecond'     # -> 1  (base 2)
find host cmd -name '*.go' -not -path 'host/daemon/*' -exec sed 's://.*::' {} + \
  | grep -c 'Nanosecond\|Microsecond'                                                        # -> 0  (the composition's other half)
# RECORDED, DECLARED NOT THE GATE (it cannot move — §4.1's comment quotes the retired spelling):
grep -rn 'Nanosecond\|Microsecond' host/daemon/                                              # -> 2 at base AND 2 at head
```
The single stripped hit must be the const definition line and its matched text must carry the
leading `-`:
```bash
sed 's://.*::' host/daemon/*.go | grep -n 'Nanosecond\|Microsecond'   # -> `const expiredReadDeadline = -1 * time.Nanosecond`
```

### T1.6 — the determinism proofs (doc §8 AC3/AC4) → closes AC3, AC4 · **RUN THESE LAST** (§5)

```bash
GOTOOLCHAIN=go1.25.6 go test ./host/daemon \
  -run '^TestDaemonReadDeadline$/^real-store-expired-deadline$' -count=1000 -v > /tmp/ac3.out 2>&1; echo "rc=$?"
grep -c -- '--- PASS: TestDaemonReadDeadline/real-store-expired-deadline' /tmp/ac3.out   # -> EXACTLY 1000
grep -c -- '--- FAIL' /tmp/ac3.out                                                       # -> 0

GOTOOLCHAIN=go1.25.6 go test ./host/daemon \
  -run '^TestTimeoutStatusMirrorsSketch$' -count=2000 -v > /tmp/ac4.out 2>&1; echo "rc=$?"
grep -c -- '--- PASS: TestTimeoutStatusMirrorsSketch' /tmp/ac4.out                        # -> EXACTLY 2000
grep -c -- '--- FAIL' /tmp/ac4.out                                                       # -> 0
```
Expected wall clock **~3 s** and **~5 s** (§0.5(iii)). These are not typos and a fast finish is not
a skipped run — **the verdict is the counted line, 1000 and 2000.**

### T1.7 — the mutation table → closes AC7

Execute §4 in full: T-MU1, T-MU2, T-MU3, in that order, each with its own backup/land-proof/vet/
kill/inverse/restore cycle.

### T1.8 — the full gates (doc §8 AC1) → closes AC1

```bash
./scripts/verify_ail.sh; echo "rc=$?"      # -> rc=0, 10 identities / 39 named tests / world gate 9/9   (~4 s)
./scripts/verify_go.sh;  echo "rc=$?"      # -> rc=0                                                    (~162 s)
```
**No `.ail` file is touched by this item**, so `verify_ail.sh`'s pins (10 / 39 / 9-of-9) are
unmoved by construction — any other totals mean something else changed: **STOP and report, never
adjust a pin.** At head the `status = 200, want 503` shape is **extinct**: any red in `host/daemon`
is real. A red in `host/broker` is the recorded ~18% base flake there — attribute by shape against
the base rate, never by one run.

### T1.9 — hygiene (doc §8 AC8) → closes AC8

```bash
git status --porcelain    # -> shows ONLY host/daemon/read_deadline_test.go and host/daemon/handlers.go as modified
                          #    (the executor makes no commits; the CONTROLLER commits)
shasum -a 256 host/daemon/read_deadline_test.go host/daemon/handlers.go
```
Record both post-M1 hashes — they are the values every mutation restore in §4 must return to
byte-identically, and AC8's assertion. **No stray files**: no `/tmp` artefacts inside the repo, no
`.orig`/`.bak` beside the sources (backups live **outside** the repo, §4.1).

---

## 4. Mutation discipline — first-class tasks, not a footnote

### 4.1 Protocol — every arm, no exceptions (P3)

```bash
# 0. BACKUP ONCE, before any mutation, OUTSIDE the repo.
mkdir -p /tmp/w19_bak
cp host/daemon/read_deadline_test.go /tmp/w19_bak/read_deadline_test.go
cp host/daemon/handlers.go           /tmp/w19_bak/handlers.go
shasum -a 256 host/daemon/read_deadline_test.go host/daemon/handlers.go   # PRE — the post-M1 hashes of T1.9
#    `git checkout --` IS FORBIDDEN. The worktree's work is UNCOMMITTED BY CONSTRUCTION
#    (the controller commits, not the executor) — `git checkout --` would destroy the milestone.

# 1. ANCHOR: assert the pre-edit token count/line is what the table says. A different reading is
#    INSTRUMENT FAILURE: stop, re-derive, do not proceed.

# 2. APPLY the exact edit from §4.2, identified BY ENCLOSING FUNCTION, never by a bare pattern
#    (§0.5(v): after this change the test file has THREE `d.readDeadline = expiredReadDeadline`
#    writes — two stimulus sites plus the pin test's own).

# 3. LANDED-PROOF, both halves:
shasum -a 256 <file>                       # POST != PRE, or THE EDIT DID NOT LAND
grep -n '<the mutated token>' <file>       # the mutated text is present, at the intended site
#    A mutation that never applied prints the same green as one that worked.

# 4. COMPILE the mutant:
GOTOOLCHAIN=go1.25.6 go vet ./host/daemon          # ALL mutants, test-file ones included
GOTOOLCHAIN=go1.25.6 go build ./host/daemon/       # PRODUCTION-file mutant (MU-DEADLINE-DETACH) as well
#    `go build ./...` DOES NOT COMPILE _test.go AT ALL — a build on a test-only mutant is VACUOUS.
#    A mutant that does not compile is INSTRUMENT FAILURE, never "survived".

# 5. KILL arm: the -run-scoped command from §4.2. Read the COUNTED `--- FAIL:` LINES, NEVER rc alone
#    (a -run matching nothing exits 0 with `testing: warning: no tests to run` — observed at base, §0.3).

# 6. INVERSE arm: the same package with -skip on the killer. ENUMERATE what redded and compare to
#    the stated red SET. `rc=0` is the criterion only where §4.2 says so.

# 7. RESTORE from the backup and prove it:
cp /tmp/w19_bak/<basename> host/daemon/<basename>
shasum -a 256 host/daemon/<basename>       # MUST equal PRE, byte-identical
```

### 4.2 The three arms, fully prescribed with pre-registered outputs

Every "expected" below was **measured by the planner** on the scratch head (§0.4), not predicted.

---

#### T-MU1 · MU-STIM-POSITIVE — test-file one-shot

| step | command | expected |
|---|---|---|
| anchor | `grep -n '^const expiredReadDeadline' host/daemon/read_deadline_test.go` | exactly 1 hit, `= -1 * time.Nanosecond` |
| edit | change that one line to `const expiredReadDeadline = 1 * time.Nanosecond` (drop the minus) | |
| land (a) | `shasum -a 256 host/daemon/read_deadline_test.go` | **≠ PRE** |
| land (b) | `grep -c 'const expiredReadDeadline = 1 \* time.Nanosecond' host/daemon/read_deadline_test.go` | **1** |
| vet | `go vet ./host/daemon` | **rc=0** |
| **KILL** | `go test ./host/daemon -run '^TestExpiredReadDeadlineExpiresAtConstruction$' -v -count=1 > /tmp/mu1_kill.out 2>&1; echo rc=$?; grep -c -- '--- FAIL:' /tmp/mu1_kill.out` | **rc=1**, `--- FAIL:` = **exactly 1**, and the failure text is the sign `Fatalf` (`expiredReadDeadline = 1ns, want a negative duration …`). **Zero timing dependence — this is the arm that makes a ~1% race killable at `-count=1`** |
| **INVERSE** | `go test ./host/daemon -skip '^TestExpiredReadDeadlineExpiresAtConstruction$' -count=1; echo rc=$?` | **rc=0** (observed). **Priced residual**: this mutant IS the base flake, so ~2% of inverse runs red with the `status = 200, want 503` / `answered 200` shape. **That red is the mutation's own phenotype** — attribute by shape and **re-measure once** (an operator re-measurement of a one-shot arm is not a test retry, doc §11). A red of any OTHER shape fails the arm |
| restore | `cp /tmp/w19_bak/read_deadline_test.go host/daemon/read_deadline_test.go; shasum -a 256 …` | **== PRE** |

---

#### T-MU2 · MU-DEADLINE-DETACH — **production-file** one-shot

| step | command | expected |
|---|---|---|
| anchor | `grep -n 'return context.WithTimeout(r.Context(), d.readDeadline)' host/daemon/handlers.go` | exactly 1 hit, **line 270** (inside `readCtx`, whose `func` line is 269) |
| edit | that line → `	return context.WithTimeout(r.Context(), readDeadline)` — i.e. the **field** `d.readDeadline` becomes the **10 s package constant** `readDeadline` (`daemon.go:128`) | |
| land (a) | `shasum -a 256 host/daemon/handlers.go` | **≠ PRE** |
| land (b) | `grep -c 'context.WithTimeout(r.Context(), readDeadline)' host/daemon/handlers.go` | **1** |
| **build** | `go build ./host/daemon/` | **rc=0** — the mutant compiles (D16's "both identifiers in scope", executed by the planner: build rc=0, vet rc=0) |
| vet | `go vet ./host/daemon` | **rc=0** |
| **KILL** | `go test ./host/daemon -run '^TestExpiredReadDeadlineExpiresAtConstruction$' -v -count=1 > /tmp/mu2_kill.out 2>&1; echo rc=$?; grep -c -- '--- FAIL:' /tmp/mu2_kill.out` | **rc=1**, `--- FAIL:` = **exactly 1**, failure text `readCtx under an already-expired deadline returned a LIVE context — the stimulus must be expired at construction, before any store read can begin` (the −1 ns field is ignored; the 10 s future deadline yields a live ctx) |
| **INVERSE** | `go test ./host/daemon -skip '^TestExpiredReadDeadlineExpiresAtConstruction$' -count=1 -v > /tmp/mu2_inv.out 2>&1; echo rc=$?; grep -- '--- FAIL:' /tmp/mu2_inv.out` | **NOT rc=0 — this mutant is multiply-killed and the red SET is the assertion (P5).** Expected, **measured twice by the planner** (§0.5(i)): **rc=1** with **exactly 4** `--- FAIL:` lines, enumerating<br>`--- FAIL: TestDaemonReadDeadline`<br>`    --- FAIL: TestDaemonReadDeadline/real-store-expired-deadline`<br>`    --- FAIL: TestDaemonReadDeadline/blocking-store`<br>`--- FAIL: TestTimeoutStatusMirrorsSketch`<br>and `--- PASS: TestDaemonReadDeadline/normal-deadline-answers-200`. **`blocking-store` is in the set and the doc's §6 row omits it** — it is the same phenotype (the field is ignored, so the 50 ms park is never released and the subtest's own 2 s watchdog fires), not an unexplained red. **Any red OUTSIDE these four lines is unexplained and fails the arm** |
| restore | `cp /tmp/w19_bak/handlers.go host/daemon/handlers.go; shasum -a 256 …` | **== PRE** |

---

#### T-MU3 · MU-SITE-REVERT — test-file one-shot · **NO KILLER TEST, BY DESIGN (P4)**

**Do not invent a test for this mutant.** Doc §6 states it plainly: *"no committed test can kill
this mutant — a `-count=1` run of the weakened test is green 99.35% of the time"*. The row **is**
iteration 93's lesson expressed as a mutation: the only per-commit detector for a re-introduced
sub-1% race is a **static sweep**, which is why that sweep is AC5 and its scope is the whole
package plus `host/ cmd/`.

| step | command | expected |
|---|---|---|
| anchor | `grep -n 'd.readDeadline = expiredReadDeadline' host/daemon/read_deadline_test.go` | **THREE** hits: the `real-store-expired-deadline` site, the `TestTimeoutStatusMirrorsSketch` site, and the **pin test's own** write. **Identify the target by enclosing function** (`awk '/^func TestTimeoutStatusMirrorsSketch/,/^}/'` or an editor jump) — a bare pattern edit hits the wrong one or all three (§0.5(v)) |
| edit | **only** the site inside `TestTimeoutStatusMirrorsSketch` → `	d.readDeadline = 1 * time.Nanosecond` (bypassing the constant) | |
| land (a) | `shasum -a 256 host/daemon/read_deadline_test.go` | **≠ PRE** |
| land (b) | `grep -n 'd.readDeadline = 1 \* time.Nanosecond' host/daemon/read_deadline_test.go` | **1** hit, and its line number is **inside** `TestTimeoutStatusMirrorsSketch` (verify with `awk '/^func TestTimeoutStatusMirrorsSketch/,/^}/' … \| grep -c 'time.Nanosecond'` → **1**) |
| vet | `go vet ./host/daemon` | **rc=0** |
| **DETECTOR (AC5's sweep — the only one)** | `sed 's://.*::' host/daemon/*.go \| grep -c 'Nanosecond\|Microsecond'` | **2** — the constant plus the reverted line, one of them positive — **against AC5's expected 1 → RED**. Wider scope reads **2** as well |
| non-vacuity of the "no killer" claim | `go test ./host/daemon -run '^TestExpiredReadDeadlineExpiresAtConstruction$' -v -count=1 \| grep -c -- '--- PASS:'` | **1** — the pin test **still passes** under this mutant. Measured. This is what "no committed test kills it" means, demonstrated rather than asserted |
| **INVERSE / package run** | `go test ./host/daemon -count=1; echo rc=$?` | **rc=0 green** (observed) — the mutant is the base flake at ONE site (~1%). **There is no `-skip` arm: there is no killer to skip.** A red here is the base-flake shape; attribute by shape |
| restore | `cp /tmp/w19_bak/read_deadline_test.go host/daemon/read_deadline_test.go; shasum -a 256 …` | **== PRE** |
| post-restore | `sed 's://.*::' host/daemon/*.go \| grep -c 'Nanosecond\|Microsecond'` | back to **1** |

---

### 4.3 After the table — the AC7/AC8 close-out

```bash
shasum -a 256 host/daemon/read_deadline_test.go host/daemon/handlers.go   # both == the T1.9 post-M1 hashes
git status --porcelain                                                    # ONLY the two intended modified files
ls /tmp/w19_bak                                                           # backups live OUTSIDE the repo
GOTOOLCHAIN=go1.25.6 go test ./host/daemon -count=1; echo "rc=$?"         # rc=0 — the tree is the milestone again
```

---

## 5. Acceptance commands, exactly as the executor must run them — with budgets, cheap first (P6)

```bash
export AILANG_BIN=/tmp/ailang-v0300/ailang GOTOOLCHAIN=go1.25.6   # EVERY command. Non-negotiable.
```

| order | AC | command | pass condition | base | **budget** |
|---|---|---|---|---|---|
| 1 | AC6(a) | `grep -c 'LIMITATION(w-daemon-late-read-503)' host/daemon/handlers.go`; `grep -c 'writeReadTimeout' host/daemon/handlers.go` | **1**; control **7** | 0; 7 | instant |
| 2 | AC6(b) | `grep -c 'err.Error()' host/daemon/handlers.go` | **5** | 5 | instant |
| 3 | AC6(c) | `go build ./host/daemon/` | rc=0 | rc=0 | ~1 s |
| 4 | AC5(a) | `sed 's://.*::' host/daemon/*.go \| grep -c 'Nanosecond\|Microsecond'` | **1** | **2** | instant |
| 5 | AC5(b) | `find host cmd -name '*.go' -exec sed 's://.*::' {} + \| grep -c 'Nanosecond\|Microsecond'` | **1** | **2** | instant |
| 6 | AC5 ctl | `sed 's://.*::' host/daemon/handlers.go \| grep -c readDeadline`; `grep -c readDeadline host/daemon/handlers.go` | **6**; **7** | 6; 7 | instant |
| 7 | AC2 | `go test ./host/daemon -run '^TestExpiredReadDeadlineExpiresAtConstruction$' -v -count=1` → count `--- PASS:` | **exactly 1** | **0** (vacuous rc=0) | ~1 s |
| 8 | AC7 | §4.2 in full — T-MU1, T-MU2, T-MU3 | every arm as pre-registered in §0.4/§4.2 | — | **~2 min** |
| 9 | AC8 | `git status --porcelain`; `shasum -a 256 <both>` | only the two files; hashes == post-M1 | clean | instant |
| 10 | AC1 | `./scripts/verify_ail.sh` | rc=0, **10 / 39 / 9-of-9** | rc=0 | **~4 s** |
| 11 | AC1 | `./scripts/verify_go.sh` | rc=0 | rc=0 (measured) | **~162 s (2.7 min)** |
| 12 | **AC3** | `go test ./host/daemon -run '^TestDaemonReadDeadline$/^real-store-expired-deadline$' -count=1000 -v` → count `--- PASS: …/real-store-expired-deadline` | **exactly 1000**, `--- FAIL` = **0** | rc=1, 15/1000 and 6/1000 failing | **~3 s** (planner-measured ×3; **not** the brief's 138 s — §0.5(iii)) |
| 13 | **AC4** | `go test ./host/daemon -run '^TestTimeoutStatusMirrorsSketch$' -count=2000 -v` → count `--- PASS: TestTimeoutStatusMirrorsSketch` | **exactly 2000**, `--- FAIL` = **0** | rc=1, 26/2000 failing | **~5–6 s** |

**Total acceptance wall clock ≈ 5 minutes**, dominated by `verify_go.sh`. The high-count ACs are
ordered last per P6 even though they turned out cheap on this rig — an executor on a loaded machine
should still not hit them before the cheap static ACs have already told it whether the edit is
right. **Ceiling rule**: if AC3 or AC4 exceeds **5 minutes**, or `verify_go.sh` exceeds **12
minutes** (its race leg carries a 600 s watchdog), stop and report the discrepancy — do not kill the
run and do not lower a count.

**The trap every test-running AC is written around**: a `-run` selector matching nothing exits **0**
with `testing: warning: no tests to run` — **observed live at base for AC2** (§0.3). **Never read the
exit code alone. Count the `--- PASS:` / `--- FAIL:` lines.** `grep` rc semantics: 1 = no match,
2 = no such path; counts are read from printed output, never from a pipe that eats the exit code.

---

## 6. Sizing and velocity

| Milestone | Impl LOC | Test LOC | Total changed | Days |
|---|---|---|---|---|
| M1_DETERMINISTIC_STIMULUS | **+27 / −0** (comment-only, `handlers.go`) | **+41 / −2** (`read_deadline_test.go`) | **≈70 lines, 2 files** | **0.25** |

The doc's §10 estimate of ~0.25 d **holds, and this plan agrees with it** — with the cost model
corrected downward by measurement, not upward:

- The doc prices the two high-count proof runs as *"minutes of wall clock"*; measured, they are
  **8 seconds combined** (§0.5(iii)). The real wall-clock cost centre is `verify_go.sh` at 162 s.
- The three one-shot mutations are literal reverts of measured states; all three were executed
  end-to-end by the planner on a scratch tree (§0.4), so the executor is reproducing pre-registered
  values, not discovering them.
- Total mechanical work: one constant + comment block, two one-line site edits, one ~26-line test,
  one 27-line comment. **≈70 lines against a 0.25 d band is far inside the measured repo velocity**
  (item 16 landed 335 Go insertions in a ~0.5 d band; the item-18 plan ran ≈845 LOC/day).

**The estimate's risk is not in the typing.** It is in the mutation table's discipline (three
backup/land/vet/kill/inverse/restore cycles) and in reading the counted lines correctly. Budget the
quarter-day there, not on the edit.

---

## 7. Execution protocol

- **Worktree**: `/Users/voightkampff/dev/sunholo-data/.wt-iter94-timeout-flake`, branch
  `sprint/w-daemon-timeout-test-flake`, base **`5380e6f`**. A **sibling of the repo, never `/tmp`**.
- **The executor has NO git write permission.** No `git commit`, `git add`, `git push`,
  `git checkout`, `git stash`, `git restore`, `gh pr`. **THE CONTROLLER MAKES THE COMMIT** — one
  commit for the one milestone.
- **Restores are `cp` from `/tmp/w19_bak/`.** `git checkout --` is **forbidden**: the milestone is
  uncommitted by construction and `git checkout --` would delete it.
- **`git add -A` is forbidden.**
- **Both exports on every command.** `AILANG_BIN=/tmp/ailang-v0300/ailang` (v0.30.0, `e37b370`) and
  `GOTOOLCHAIN=go1.25.6`. Without `AILANG_BIN`, 17 `host/verifygate` arms fail **loudly by design**
  — a base condition, not a regression. Without `GOTOOLCHAIN`, `verify_go.sh` FATALs on the
  denylisted local go1.26.4.
- **Never validate anything against a `-dirty` dev `ailang`.** No `.ail` file is touched by this
  item; the gate binary is the pinned one.
- **`scripts/verify_go.sh` IS this sprint's gate** (unlike item 18's plan): the driver drift gate is
  green in this worktree and was measured rc=0 at base (§0.2). If it ever FATALs on **DRIVER
  DRIFT**, that red means *"the fleet must commit"* — **never** absorb `tools/launchd/*` into this
  change. Stop and report.
- **The `Observatory: 302MB` warning on stderr is expected** and does not red the gates (measured).
- **Do not edit the design doc.** Findings go to the controller for the record commit.
- **Attribute by shape, at base, before blaming the diff.** `host/broker` carries a recorded ~18%
  base flake; `host/capsule` and `host/archive` carry rows 20/21. None of them is this diff.

---

## 8. P8 — who reports what: executor deliverables vs controller re-derivation

**Generator ≠ judge.** The controller re-runs every gate independently after the executor hands
off; the executor's report is evidence, not the verdict.

**The EXECUTOR reports (its deliverable):**

1. The two edited files, in the worktree, uncommitted.
2. **Per-AC observed output** for AC1–AC8 in §5 order: the actual counted numbers (`--- PASS:` /
   `--- FAIL:` / grep counts / rc), not "green".
3. **The full mutation ledger** for T-MU1/T-MU2/T-MU3: PRE hash, POST hash, the landed-proof grep,
   the vet (and build) rc, the kill arm's counted `--- FAIL:` lines **with the failure text**, the
   inverse arm's **enumerated** red set, and the restore hash == PRE.
4. **Post-M1 sha256 of both files** (the values AC8 and every restore compare against).
5. **Any deviation from a pre-registered value in §0.4, verbatim, with the command that produced
   it** — a deviation is a STOP-and-report, never a silently adjusted pin.
6. **Any base-shape red**, attributed at base with the re-measurement, per §7.
7. LOC actually changed, against the plan's `+41/−2` and `+27/−0`.

**The CONTROLLER re-derives independently (the executor should NOT treat these as its own pass
condition, and must not pre-empt or "pre-agree" them):**

- Both full gates (`verify_ail.sh`, `verify_go.sh`) on the **final tree**, and again on the **merge
  commit** (Gate 3b).
- AC2, AC3, AC4, AC5(a)/(b) and their controls, AC6(a)/(b)/(c) — from the committed tree.
- At least one mutation arm re-executed first-party (the doc's own precedent: a mutation table is a
  claim until someone other than its author runs a row of it).
- The doc↔plan AC identity check of §2 against the committed doc.
- Whether §0.5(i)'s red-set correction and §0.5(ii)'s gofmt reflow should be carried into the
  design doc or recorded as sprint-plan findings.
- Filing the two named follow-ups (`w-daemon-late-read-503`; the busy_timeout-vs-readDeadline
  ordering residual of doc §11) — **record actions, not sprint work**.

---

## 9. Risks

| # | Risk | Severity | Mitigation |
|---|---|---|---|
| **R1** | Executor scores MU-DEADLINE-DETACH's inverse arm as FAILED because `blocking-store` reds and the doc says any red outside its set fails the arm | **HIGH** without §0.5(i) | §0.5(i) + §4.2 T-MU2 carry the **measured** four-line red set and the mechanism that puts `blocking-store` in it |
| **R2** | A mutation "passes" because it never landed — especially T-MU3, whose target string occurs **three** times in the head file | **HIGH** | §4.1 steps 1/3 (anchor + sha256 + token grep) and T-MU3's by-enclosing-function rule; the planner's own near-miss is recorded in §0.5(v) |
| **R3** | Executor reads `rc` instead of counted lines and scores AC2 green on the **vacuous** `no tests to run` | **MEDIUM** | AC2's base reading (rc=0 with **0** PASS lines) is in §0.3; §5's trap paragraph; every AC's pass condition is a count |
| **R4** | A base-flake red (`status = 200, want 503`) at base or in a mutation inverse arm is misattributed to the diff | **MEDIUM** | ~2.1% priced in §0.2 with the measured rates; attribute-by-shape and re-measure-once rules in §4.2 T-MU1 and §7 |
| **R5** | Executor "improves" the §4.3 comment and adds `err.Error()`, `Nanosecond`, `Microsecond`, or bare `readDeadline` — silently moving AC5 or AC6(b) | **MEDIUM** | T1.3's explicit prohibition; the four controls in T1.3/§5 rows 2, 4, 6 catch each one |
| **R6** | Executor treats gofmt's reflow of the §4.3 list as a mis-paste and reverts it, or leaves the tree non-gofmt-clean | **LOW** | §0.5(ii) pre-registers the reflow with the exact landed form quoted in T1.3; T1.4 makes `gofmt -l` empty a task |
| **R7** | Executor mistakes a 3-second AC3 for a skipped run, or kills a slow `verify_go.sh` as a hang | **LOW** | §0.5(iii) budgets from four first-party samples; §5's ceiling rule says report, never kill or lower a count |
| **R8** | `git checkout --` used to clean up a mutant, destroying the uncommitted milestone | **MEDIUM** | §4.1 step 0/7 and §7: `cp` from `/tmp/w19_bak/` only; the prohibition appears three times because this is unrecoverable |
| **R9** | `verify_go.sh` FATALs on driver drift mid-sprint (a fleet commit lands, or something writes `tools/launchd/`) | **LOW** | Green at base here (6 tracked driver files, tree matches HEAD). If it reds: **stop and report**; that red means the fleet must commit, never that World absorbs it |

**No risk on this list wants a controller decision before execution.**

---

## 10. Handoff

- **Plan**: `design_docs/planned/w-daemon-timeout-test-flake-sprint-plan.md`
- **Progress JSON**: `sprint_w-daemon-timeout-test-flake.json` (repo root)
- **Milestone**: `M1_DETERMINISTIC_STIMULUS` (0.25 d) — the only one
- **Total**: 0.25 days, ≈70 changed lines across 2 files, **risk low**
- **Closes**: AC1–AC8, all eight, in one milestone (§2)
- **Findings for the record commit**: §0.5(i) the doc's MU-DEADLINE-DETACH red set omits
  `blocking-store` (measured twice); §0.5(ii) gofmt reflows §4.3 and one line of §4.2 (AC-neutral,
  landed forms quoted); §0.5(iii) the high-count ACs run in seconds, not the brief's 138 s;
  §0.5(iv) the flake is a plain-leg phenomenon (0/1000 under `-race`); the doc's own §13
  observation that `handlers.go`'s "~+18/−0" estimate is really **+27/−0**.
- **Carried forward, named, NOT this sprint**: `w-daemon-late-read-503` (the ARM B contract, with
  the §3 B.3 fixture question answered under its own quorum) and the busy_timeout-vs-readDeadline
  ordering residual (doc §11) — **both are controller/record actions**.
- **Unblocks on landing**: queue rows 20 and 21 become the next flake work with one fewer
  confounder in `go test ./...`; **no shared root cause with them is claimed** (doc §9).

**SPRINT_PLAN_PATH**: `design_docs/planned/w-daemon-timeout-test-flake-sprint-plan.md`
**SPRINT_JSON_PATH**: `sprint_w-daemon-timeout-test-flake.json`
