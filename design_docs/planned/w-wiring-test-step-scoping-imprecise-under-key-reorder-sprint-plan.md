# Sprint plan — w-wiring-test-step-scoping-imprecise-under-key-reorder (queue row 52)

**Design doc**: `design_docs/planned/w-wiring-test-step-scoping-imprecise-under-key-reorder.md`
**Sprint id**: `w-wiring-test-step-scoping-imprecise-under-key-reorder`
**Base**: worktree `.planner-wt-iter147` at `f6773dc` (identical to `dev`), `git status --porcelain` = 0
**Direction**: human-ratified `D-WORLD-30` — shallowest-enclosing-`steps:` line scan + Mechanism
step 1b's `expectedStepCol = 6` pin. **Not re-litigable.** This plan does not propose an
alternative locator and found no reason to want one.
**Estimate**: ~0.1d, 3 milestones. Scope: ONE function in ONE file,
`host/verifygate/toolchain_pin_gate_test.go` (func at `:492`, the two offending
`HasPrefix(strings.TrimSpace(...), "- name:")` scans at `:511` and `:524`).

---

## 0. What the planner actually ran (drill BEFORE plan)

I ran **15 mutation arms** against the landed gate in this worktree, plus a Python prototype of
the ratified locator (Mechanism steps 1, 1b, 2, 3, 4, 5, both anchor derivations) over the same
15 files. Protocol per arm: mutant generated from the doc's own recipe into
`.github/workflows/ci.yml`; **`sha256` printed and matched against the doc's pin where one
exists**; effect asserted with `ruby -ryaml` (step count, per-step key list, `continue-on-error`
value) — never against bytes; `go vet ./host/verifygate/` rc read **before** any test result;
AC1 run; `.github/workflows/ci.yml` restored from a `cp` backup and **re-asserted byte-identical
by sha256 after every single arm** (`aed8e186…`, 15/15); pristine control green before and after
the batch. No `git checkout --` was used and no git write operation was run.

**The doc's recipes are executable.** MUT-B, MUT-O, MUT-P, MUT-J, MUT-K and MUT-Q all
reproduced the doc's pinned sha256 **exactly, first try**, from the prose recipe alone:

| mutant | doc pin | planner-derived sha256 | match |
|---|---|---|---|
| base | `aed8e186…ed85` | `aed8e186fb57036eb6b03509cbb668850d577d46e6cc68b30e7a4c042108ed85` | ✅ |
| MUT-B (V23) | `c1903a86…` | `c1903a869701e47626a4fb1cd537ffb42a935d02fa1ace26b2f39f330a6cecbf` | ✅ |
| MUT-O (V25) | `72acc72d…` | `72acc72d239056d7772b769b345212d092dce715a229202dee939899415f9f31` | ✅ |
| MUT-P (V24) | `1917c413…` | `1917c413b1942f230da15dd026a432ed674e447e30e407263d2307f190d88e56` | ✅ |
| MUT-J | `96f51f64…` | `96f51f646354042b3474b0614b1bc7de35624f3effc93d1ca928458c3ec6bdf6` | ✅ |
| MUT-K | `58a2dd93…` | `58a2dd93f537acc7faab108616fb2095d0bc7d225200cfc5c61701b31db2c336` | ✅ |
| MUT-Q (V31) | `4fd378c4…` | `4fd378c4f64a9e7813586415d1cad3950faf2620db508d6b1a9fda1a69c413a6` | ✅ |
| MUT-A | `28731ce9…` (V2) | `f73024e8…` — **different construction, same verdict class** | ⚠ see §2.4 |
| MUT-D | `5f61658a…` (V17) | `8cb5ca8d…` — **different construction, same verdict (coe@175)** | ⚠ see §2.4 |
| MUT-I | `2f795e56…` (V15) | `1131db25…` — **different construction, same verdict (coe@176)** | ⚠ see §2.4 |

Every OLD (landed-gate) cell in the doc's Mutation Drill table **reproduced**, including the two
live defects the sprint exists to close:

| arm | doc's OLD cell | planner-measured OLD | NEW (prototype, shallowest) | doc's NEW cell |
|---|---|---|---|---|
| MUT-A (ARM A) | rc=1, false positive on unrelated step | **rc=1 `ci.yml:167`** = the *unrelated* step's flag line | GREEN `stepCol=6 start=175 end=179` | rc=0 ✅ |
| MUT-B (ARM B) | rc=0 PASS, **fail-open** | **rc=0 `--- PASS`** with ruby showing `step[6] coe=true` | RED `coe@181`, `start=174 end=183` | rc=1 on own flag ✅ |
| MUT-C | rc=1 | rc=1 `ci.yml:176` | RED `coe@176` | ✅ |
| MUT-D | rc=1 blaming stranger `coe@175` | **rc=1 `ci.yml:175`** (exact match to V17) | GREEN `start=177 end=181` | rc=0 ✅ |
| MUT-E | rc=1 `count(...)=2` fatal | rc=1 `count(…)=2` | fatal (pre-scan) | ✅ |
| MUT-F | rc=0 | rc=0 | GREEN | ✅ |
| MUT-G | rc=0 (OLD pins no name) | **rc=0 `--- PASS`** | FATAL InvB | ✅ |
| MUT-I | rc=1 `coe@176` incidental | rc=1 `ci.yml:176` | RED `coe@176 start=174` | ✅ |
| MUT-J | rc=1 `coe@177` | rc=1 `ci.yml:177` | RED `coe@177`, **`start == i == 174`** | ✅ |
| MUT-K | rc=0 | rc=0 | GREEN, `start == i == 174` | ✅ |
| MUT-O | rc=1 `ci.yml:182` (V28) | **rc=1 `ci.yml:182`** — exact | RED `coe@182 stepCol=6 start=174 end=184` (= V25 exactly); NEAREST anchor derives `stepCol=12` (= V25) | ✅ |
| MUT-P | rc=1 `ci.yml:181` (V24) | **rc=1 `ci.yml:181`** — exact | RED `coe@181` | ✅ |
| MUT-Q | rc=0 GREEN (V31) | **rc=0 `--- PASS`** | FATAL InvA `[96,99)`, `stepCol=6` from the OTHER JOB (= V31 exactly) | ⚠ **see Finding 1** |
| MUT-Q2 (new) | — | rc=0 (not run as AC) | FATAL InvA `[96,99)` | new |
| MUT-R (new) | — | **rc=0 GREEN** | **FATAL, step 1b PIN: `stepCol=8 != 6`** | new |

`go vet ./host/verifygate/` was rc=0 before every one of the 15 test results.

---

## 1. PLANNER FINDINGS — read these before executing

### Finding 1 (BLOCKING for AC8 as written) — **AC8 cannot fail; step 1b's pin is unreachable on MUT-Q**

AC8 is titled *"the column pin is live and loud, proven by mutation"*. It is not proven by
MUT-Q. Measured on the locator prototype, both with and without step 1b:

| mutant | ratified locator **with** step 1b | ratified locator **with step 1b neutered** | discriminating? |
|---|---|---|---|
| MUT-Q  | FATAL — **Invariant A**, block `[96,99)` | FATAL — **Invariant A**, block `[96,99)` | ❌ **no flip** |
| MUT-Q2 | FATAL — **Invariant A**, block `[96,99)` | FATAL — **Invariant A**, block `[96,99)` | ❌ **no flip** |
| MUT-R  | **FATAL — step 1b PIN, `derived step column 8 != expectedStepCol 6`** | **GREEN, rc=0** | ✅ **flips** |

The reason is structural: under the **shallowest** anchor, MUT-Q's derivation captures the
*other job's* `steps:` key (Declared Residual 8) and therefore returns `stepCol=6`, which
**equals** `expectedStepCol`. The pin's inequality branch is never taken. The refusal comes
entirely from Invariant A, which exists with or without step 1b. AC8's own neuter arm concedes
this (*"must still red — via Invariant A or B"*) — so **both halves of AC8 pass identically in a
build that has step 1b and in a build that does not.** By the mission's own MUT-N rule (iter-141:
"a `grep -c` on the invariant cannot prove it is live — only the flip does"), AC8 is a grep in
disguise.

**Fix, measured, in this plan**: AC8 is replaced by **AC8'** built on a new arm **MUT-R** —
re-indent **both** job bodies two columns right, so the shallowest `steps:` in the file sits at
indent 6 and its first item at indent 8. Measured:
- sha256 `575e8788210be5499fe22aeb106c0849a5e469e1aae59ec3993a2190f28bb953`
- ruby oracle: **valid YAML, `ailang-verify: steps=7`, `go-verify: steps=10`** — a legitimate
  reformat, semantics preserved
- **landed gate today: rc=0 `--- PASS`** — the OLD gate silently absorbs it, so the arm is
  non-vacuous
- ratified locator: **FATAL via step 1b's pin** (`stepCol=8`); with 1b neutered: **rc=0 GREEN**.
  That rc=1 → rc=0 flip is the liveness proof AC8 was supposed to carry.
MUT-Q is **kept** as AC8'-b (a second, independent loud-refusal arm) but re-labelled honestly:
it proves *Invariant A / Residual 8*, not the pin.

### Finding 2 (REFUTES a doc claim) — **the MUT-Q row's landed-assertion is false**

The MUT-Q row instructs the executor to *"assert LANDED by sha-differs from `aed8e186…` and by
the ruby view still parsing to 10 steps."* Measured on the mutant that reproduces the doc's own
pinned V31 sha `4fd378c4…`:

```
ailang-verify: steps=7
go-verify: steps=NIL          <-- not 10; the job has NO steps key at all
```

Shifting only the `steps:` block two columns right moves the key from indent 4 (a job key) to
indent 6 — where it becomes a **member of the `env:` mapping**:

```
102|    env:
103|      GOTOOLCHAIN: go1.26.6
104|      steps:                <-- now an env var, not a job key
```

The document still parses (so "valid YAML" is true), but the go-verify job is destroyed. **The
executor cannot satisfy the row's ruby assertion, and must not weaken it to make the arm pass.**
The corrected legitimate re-indent is **MUT-Q2** — shift the *whole go-verify job body* (line
100 onward) two columns right — sha256
`e50029d14bb1acb8374c279310594186c763b06643a601902d5313a535dd7311`, ruby view
`ailang-verify: steps=7 / go-verify: steps=10`. **Doc line to correct: the MUT-Q row's
landed-assertion clause in §Mutation Drill (the `| MUT-Q |` table row).** The V31 *measurements*
are all correct and reproduced exactly; only the assertion prose is wrong.

### Finding 3 (STALE) — **§Milestones predates the round-3 carve-out**

The doc's MS1/MS2/MS3 split is **structurally right** (locator / drill / write-up) and I keep it.
Its contents are stale by one revision:
- **MS1 does not mention Mechanism step 1b** (`expectedStepCol`) — the whole round-3 `gpt5-6-sol`
  fix is absent from the milestone that has to implement it.
- **MS2 says "MUT-A through MUT-P"** — MUT-Q, added at round 3, is owned by no milestone.
- **MS2 maps to "AC2-AC5, AC7"** — **AC8 is orphaned**; no milestone claims it.
- MS2 also claims AC7 (byte-identity at landing), which is a landing check, not a drill check.
§3 below fixes all four.

### Finding 4 (minor, non-blocking) — three mutant pins are absent, and three prose numbers drift

- MUT-A, MUT-D and MUT-I carry **no pinned sha256** (only MUT-B/O/P/J/K/Q do). My constructions
  are legitimate members of each class and produce the doc's stated verdicts exactly (MUT-D:
  `coe@175`, matching V17; MUT-I: `coe@176`, matching V15), but they are not byte-identical to
  whatever V2/V15/V17 built. The executor must pin its own and record them — the doc's
  "assert LANDED by MATCHING the pinned hash" rule cannot be applied to these three.
- ARM A's blamed line is **167** in my construction, 166 in V2, 164 in V21. Same class, three
  different constructions. The AC must therefore say "**a line belonging to another step as the
  ruby oracle sees it**", never a fixed line number.
- The MUT-O prose says the in-scalar `steps:` decoy sits at "col 12"; measured, `steps:` is at
  **col 10** and the decoy dash at **col 12** (which is what makes the nearest anchor derive
  `stepCol=12`, as V25 records). Cosmetic; V25's numbers are right.

### Finding 5 (no action; recorded) — **the ratified mechanism is sound on every arm I ran**

The prototype was GREEN/RED/FATAL exactly as the doc's NEW column predicts on all 15 arms, and
every V21/V23/V25/V31 number reproduced to the line (`start=174 end=183` for MUT-B;
`stepCol=6 start=174 end=184` for MUT-O; `stepCol=12` for the nearest anchor on MUT-O;
`stepCol=6 start=96 end=99` for MUT-Q). **I found no reason to dispute `D-WORLD-30`.** Residual 8
(cross-job capture) is real and I reproduced it, and it fails loud in every arm I ran.

---

## 2. Base measurements for every acceptance command

Run from the repo root with `export PATH=/opt/homebrew/bin:$PATH`. **Measured by me in this
worktree at `f6773dc`, this iteration.**

| # | command | measured at base |
|---|---|---|
| B1 | `AILANG_BIN=/tmp/ailang-v0300/ailang go test ./host/verifygate/ -run '^TestMiscompileInstrumentStepIsGatedInCI$' -count=1 -v` | **rc=0**, `=== RUN` present, `--- PASS` |
| B2 | `go vet ./host/verifygate/` | **rc=0** |
| B3 | `gofmt -l host/verifygate/` | **rc=0, zero lines** |
| B4 | `AILANG_BIN=/tmp/ailang-v0300/ailang go test ./host/verifygate/ -count=1` | **rc=0** (~67s) |
| B5 | `shasum -a 256 .github/workflows/ci.yml` | **`aed8e186fb57036eb6b03509cbb668850d577d46e6cc68b30e7a4c042108ed85`** |
| B6 | `AILANG_BIN=/tmp/ailang-v0300/ailang ./scripts/verify_ail.sh` | **rc=0** (11 identities / 40 named tests / 9-of-9 package steps) |
| B7 | `git status --porcelain \| wc -l` | **0** |

**Commands that are RED at base and are therefore FORBIDDEN as acceptance criteria** (all
re-measured by me, all matching the doc):
- `go build ./host/verifygate/` → **rc=1**, "no non-test Go files" (test-only package). Use B2.
- bare `go test ./host/verifygate/ -count=1` → **rc=1, 17 FAIL**, all "AILANG_BIN is unset"
  instrument-failure refusals. Deliberate anti-false-green guard. Always carry `AILANG_BIN`.
- bare `./scripts/verify_ail.sh` → **rc=1**, `AILANG_BIN refused [DEV_BUILD] … v0.34.0-247-…-dirty`.
- `./scripts/verify_go.sh` → **FORBIDDEN** (queue row 58 flake: 4 runs / 2 red / 2 green / 3
  different failing sets). No AC in this sprint cites it.

---

## 3. Milestones (revised — see Finding 3)

### MS1 (~0.05d) — rewrite the locator, **including Mechanism step 1b**
`host/verifygate/toolchain_pin_gate_test.go`, inside `TestMiscompileInstrumentStepIsGatedInCI`
only. ~40 changed lines + one helper + two consts. **No other file changes; `.github/workflows/ci.yml`
is not touched by MS1.**

1. `func indentOf(s string) int { return len(s) - len(strings.TrimLeft(s, " ")) }`
2. `const miscompileStepName = "Measure compiler reproducer (platform-conditional, gated)"`
3. `const expectedStepCol = 6`
4. **step 1** — `stepCol` from the **SHALLOWEST** enclosing `steps:` above `i`, tie broken at the
   minimal indent toward the nearest such line (total rule); `stepCol` = indent of the first
   following line whose trimmed text has prefix `- `. Either lookup failing ⇒
   `t.Fatalf("instrument failure: …")`.
5. **step 1b** — immediately after: `if stepCol != expectedStepCol { t.Fatalf("instrument
   failure: derived step column %d; update expectedStepCol after an intentional ci.yml
   re-indent", stepCol) }`.
6. **step 2** — `start`: walk back from `j := i` (**inclusive**, so a `- run:` dash-key
   identifying line matches itself) to the nearest line with trimmed prefix `- ` **and** indent
   **exactly** `stepCol`.
7. **step 3** — `end`: from `start+1`, first dash line at exactly `stepCol`, or first non-blank
   non-comment line at indent **strictly less than** `stepCol`; EOF otherwise. Blank and `#`
   lines never terminate.
8. **Invariant A** — `start <= i < end` else `t.Fatalf`.
9. **Invariant B** — the block must contain a line whose trimmed text is exactly
   `"- name: "+miscompileStepName` or `"name: "+miscompileStepName`, else `t.Fatalf`.
10. The existing `continue-on-error` loop over `[start,end)` and everything after it is
    **unchanged**.
11. Update the function's doc-comment: the `DECLARED RESIDUAL` block's scoping sentence and its
    reference to the **row-44** doc's V19 now describe the new mechanism, and the comment names
    ruling `D-WORLD-30`.

**Independently CI-green**: yes — verified by prototype, base ci.yml yields
`stepCol=6 start=174 end=178 GREEN`, identical to the landed gate's verdict.
**Gates**: AC1, AC6.

### MS2 (~0.04d) — mutation drill, **MUT-A … MUT-R** (16 arms)
No production change. Every arm: generate → **sha256 pin match** (MUT-B/O/P/J/K/Q/Q2/R) or
**pin-and-record** (MUT-A/C/D/E/F/G/I) → `ruby -ryaml` effect assertion → `go vet` rc read
**first** → AC1 command → record rc **and message text** → `cp` restore → **sha256 back to
`aed8e186…`** → pristine re-control after the batch. **Never `git checkout --`** (the tree holds
uncommitted work). Backup lives **outside** `/tmp` and outside the repo.
**Gates**: AC2, AC3, AC4, AC5, **AC8'** (new), AC8'-b.

### MS3 (~0.01d) — landing hygiene + evidence write-up
Byte-identity of `.github/workflows/ci.yml`, clean porcelain modulo the intended test-file and
doc changes, `verify_ail.sh` with the pin, sprint notes carrying every arm's rc + message +
sha256, and the doc moved to `implemented/` per repo flow. **Gates**: AC6, AC7.

---

## 4. Acceptance criteria (as revised)

- **AC1 (base green, scope-reach proven)** — B1 → **rc=0** with `=== RUN
  TestMiscompileInstrumentStepIsGatedInCI` present in `-v` output. *Base: rc=0 with the RUN
  line (measured).*
- **AC2 (fail-open family dies)** — MUT-B (`c1903a86…`, ruby: `step[6] coe=true`) → **rc=1**
  naming a line inside the miscompile step's own block *as the ruby oracle sees it*; MUT-I →
  **rc=1** on the step's own flag line. Restore to `aed8e186…` after each.
  *OLD measured: MUT-B rc=0 PASS (the live defect); MUT-I rc=1 `ci.yml:176`.*
- **AC3 (false-positive / key-order family closes)** — MUT-A (flag on the unrelated `go build +
  test gate` step + miscompile key reorder; ruby: `step[5] coe=true, step[6] coe=nil`) →
  **rc=0**; MUT-K (`58a2dd93…`) → **rc=0**. *OLD measured: MUT-A rc=1 blaming the unrelated
  step's own flag line (the live defect); MUT-K rc=0.* **Do not pin a line number here** —
  Finding 4.
- **AC4 (the original kill still fires, in every layout)** — MUT-C → **rc=1**; MUT-J
  (`96f51f64…`) → **rc=1** naming the step's own flag line. *OLD measured: rc=1 `ci.yml:176`
  and rc=1 `ci.yml:177`.*
- **AC5 (Invariant B live, proven by flip)** — MUT-G → **rc=1** with an `instrument failure`
  from Invariant B; then with Invariant B neutered in the test source via `if false && …`
  (neutered, **not deleted**, so a compile error cannot masquerade as the guard firing) the same
  MUT-G → **rc=0**. Revert both. *OLD measured: MUT-G rc=0 `--- PASS` — the OLD scan pins no
  name, so the arm is non-vacuous.*
- **AC6 (package hygiene + package-wide green)** — B2 rc=0, B3 zero lines, B4 rc=0.
  *Base: rc=0 / 0 lines / rc=0 (all measured).*
- **AC7 (no collateral)** — B5 = `aed8e186…` at landing; B7 shows only the intended test-file
  and doc changes; B6 rc=0. *Base: all three measured green.*
- **AC8' (the column pin is live, proven by a FLIP — replaces AC8)** — **MUT-R**
  (`575e8788210be5499fe22aeb106c0849a5e469e1aae59ec3993a2190f28bb953`; ruby oracle must show
  `ailang-verify: steps=7` **and** `go-verify: steps=10`) → **rc=1** whose message is step 1b's
  `derived step column 8`; then with **step 1b's `t.Fatalf` neutered via `if false && …`** the
  same MUT-R → **rc=0**. Revert both, sha back to `aed8e186…`. *OLD measured: MUT-R rc=0 `--- PASS`
  — the landed gate silently absorbs it, so the arm is non-vacuous. Prototype measured: FATAL
  (pin) with 1b, GREEN without. This is the flip AC8 lacked.*
- **AC8'-b (Residual 8 fails loud, not silent)** — **MUT-Q2**
  (`e50029d14bb1acb8374c279310594186c763b06643a601902d5313a535dd7311`; ruby oracle
  `ailang-verify: steps=7`, `go-verify: steps=10`) → **rc=1** carrying an `instrument failure`,
  **and the record must state WHICH refusal fired** (measured on the prototype: **Invariant A**,
  block `[96,99)` — the cross-job capture of Declared Residual 8, *not* the pin). MUT-Q
  (`4fd378c4…`) may be run as a third arm **only with the corrected landed-assertion of
  Finding 2** (`go-verify: steps=NIL`), never with "still parsing to 10 steps".
  *OLD measured: MUT-Q rc=0, MUT-Q2 rc=0.*

**No AC in this sprint cites `./scripts/verify_go.sh`** (row-58 flake). **No arm carries a
`-skip <test>` rc=0 criterion**: every verdict is read from the targeted AC1 command, per the
doc's standing drill rule 2. The package-wide red set for the ci.yml mutants was not measured
and no AC depends on it.

## 5. Risks

- **Low.** One test function; production code untouched; the ratified mechanism was prototyped
  against 15 arms with zero surprises.
- The one real risk is **executing AC8 as the doc writes it** and recording a green that proves
  nothing (Finding 1). AC8' closes it.
- `AILANG_BIN=/tmp/ailang-v0300/ailang` must be exported on every test command; three of the
  four obvious gate commands are red at base without it.

## 6. Out of scope (deferred, named)

- Bounding the backward `steps:` scan to the enclosing job (would close Residual 8 outright) —
  a mechanism change nobody has ratified. Named as follow-up, not taken.
- Any YAML parser / new `go.mod` dependency — **rejected by `D-WORLD-30`**.
- `actionlint` in CI — rig-dependent, and would not catch either arm (both mutants parse).

**SPRINT_PLAN_PATH**: `design_docs/planned/w-wiring-test-step-scoping-imprecise-under-key-reorder-sprint-plan.md`
**SPRINT_JSON_PATH**: `design_docs/planned/sprint_w-wiring-test-step-scoping-imprecise-under-key-reorder.json`
