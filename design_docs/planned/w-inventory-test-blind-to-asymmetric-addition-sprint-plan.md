# Sprint plan — `w-inventory-test-blind-to-asymmetric-addition` (queue row 51, clause 2)

**Design doc:** [`design_docs/planned/w-inventory-test-blind-to-asymmetric-addition.md`](w-inventory-test-blind-to-asymmetric-addition.md)
(quorum round 1 BLOCKED → one designer revision → round 2 BLOCKED at full strength → cleared by the
bounded narrow-refinement carve-out, all three reviewers' fixes applied verbatim, each with its own
AC and named mutation)
**Design-doc base:** `f3681b0` · **Sprint base:** `dev` @ `f3681b0`, working tree clean (porcelain 0
at planning entry AND exit)
**Planner:** mission-control iteration 141, `opus`, in the ephemeral worktree
`/Users/voightkampff/dev/sunholo-data/.planner-wt-iter141`
**Estimate:** 0.1 day (≈1.0 h), 2 milestones, 2 commits · **Risk:** low
**Scope:** `host/verifygate/floor_raise_inventory_test.go` **only**. One file. No new package, no new
script, no new CI job, no `.ail`, and **no edit to `design_docs/coding-standards.md` §S8** (S8 is
ratification-class; this sprint READS it and mutates it only inside restored drill arms).

---

## 0. The one-paragraph version

The gate that keeps the floor-raise coupling inventory's two homes in sync asserts that each of six
known rows is present-and-unique in each home, and nothing else. A fabricated seventh row added to
either home alone leaves it GREEN — I reproduced that this session in both directions (§2.4). The
change is ~64 added lines in one test file: two enumerator regexes, a `canonicalSiteSet` helper, a
per-home duplicate guard, a `>= 6` cardinality floor with an instrument-failure `t.Fatalf`, and a
SYMMETRIC set-difference between the homes. **I applied the design doc's own sketch in this
worktree, ran all fifteen mutation arms against it, and restored the tree byte-identical** (§3) —
so no expected red in this plan is predicted; every one is measured. Five defects in the doc's own
drill text fell out of that (§4), including one arm (**C1**) that is a **no-op as literally
written**, and one base-state premise (**AC12 part 2**) that my re-derivation **refutes**: on this
rig, this session, `./scripts/verify_go.sh` is **rc=0 with an EMPTY failing-test set**.

---

## 1. Environment — load-bearing, on EVERY command

```sh
export PATH=/opt/homebrew/bin:$PATH
export AILANG_BIN=/tmp/ailang-v0300/ailang
```

| Pin | Measured this session | Why it is load-bearing |
|---|---|---|
| `AILANG_BIN=/tmp/ailang-v0300/ailang` | `AILANG v0.30.0`, commit `e37b370`, built `2026-07-19T09:27:00Z` (printed by `verify_go.sh`'s own header) | Without it `scripts/verify_go.sh` FATALs and a large slice of `host/verifygate` fails loudly by design. That red is a BASE CONDITION, never a regression. |
| `PATH=/opt/homebrew/bin:$PATH` | `gh`, `go`, `perl`, `shasum` all resolve | An `rc=127` in this repo is a PATH gap, not a broken toolchain. |
| **NO `GOTOOLCHAIN` pin** | ambient `go version` = **`go1.26.6 darwin/arm64`**; `go env GOTOOLCHAIN` = `auto` | `scripts/verify_go.sh:214-222` deny-lists **`go1.26.0`–`go1.26.5`** only. `go1.26.6` is NOT on that list and the full gate ran rc=0 under it this session. **Do not copy the `GOTOOLCHAIN=go1.25.6` pin from this repo's older sprint plans** — it was correct when the rig's ambient go was `go1.26.4`; it is not this rig's state today. |

---

## 2. Pristine-tree baselines — every acceptance command run by me on the unchanged tree at `f3681b0`

Rule 3e: an acceptance command is baselined on the untouched tree and the base result is recorded
**as part of the criterion**. A gate already red at base measures the repo, not the change.

**Base sha256 of the three mutation venues** (the restore oracle for the whole drill):

| Path | sha256 |
|---|---|
| `host/verifygate/floor_raise_inventory_test.go` | `a00790f1fb19f551a69478217aae7657ee12b590376b365e1d28ea9707dd113f` |
| `scripts/verify_ail.sh` | `5a1bbe8963ebc54f68c8a50f8f1a5ede4aff82ab83b7b1080e6594418e41475a` |
| `design_docs/coding-standards.md` | `b710a510be0c55d84b6ba92059a26c058d895df5b103d5dd5c2dd43c13f68329` |

### 2.1 The base-state table — every acceptance command, its pristine result, and its verdict

The test command used throughout (call it **T**):

```zsh
export PATH=/opt/homebrew/bin:$PATH
export AILANG_BIN=/tmp/ailang-v0300/ailang
go test ./host/verifygate/ -run '^TestFloorRaiseInventoryNamesEveryCoupledFile$' -count=1 -v > /tmp/t 2>&1; rc=$?
```

| AC | Command (as the executor will run it) | **PRISTINE-BASE RESULT, measured by me** | Verdict on the criterion |
|---|---|---|---|
| AC1 / AC10 | **T** on the untouched tree | **rc=0**, `=== RUN`=1, `--- PASS`=1, `--- SKIP`=0, no `no tests to run`, 0.31 s | GREEN at base **by design** — it is the control, not a tooth. Legitimate. |
| AC2 | `grep -c 'disagree on their site set' host/verifygate/floor_raise_inventory_test.go` ≥ 3 | **0** (grep rc=1) | **RED at base — correct.** Post-change I measured **3**. |
| AC3 | `grep -c 'instrument failure, not an empty inventory' <file>` = 2 | **0** (grep rc=1) | **RED at base — correct.** Post-change I measured **2**. |
| AC4 | **T** with ARM A1 landed → rc=1, `--- FAIL`=1, names `host/store/some_new_file.go` | **rc=0, `--- PASS`** — the gate is blind (§2.4) | **RED at base — this is the row's whole point.** |
| AC5 | **T** with ARM A2 landed → rc=1, `--- FAIL`=1, names the site | **rc=0, `--- PASS`** — blind in this direction too | **RED at base.** |
| AC6 | **T** with ARM A3 landed → rc=0, `--- PASS`=1 | rc=0, `--- PASS` | GREEN at base — a **regression guard**, non-discriminating alone. Its tooth is that it must STAY green after M2. |
| AC7 | **T** with ARM A4 landed → rc=0, `--- PASS`=1 | rc=0, `--- PASS` | GREEN at base — regression guard, same note. |
| AC8 | **T** with R1 / R2 landed → rc=1 | rc=1 for both (legacy per-row `count == 1` fires) | GREEN at base — existing behaviour; the AC is a non-regression guard. |
| AC9 | **T** with I1 / I2 landed → rc=1 + instrument-failure message | **not runnable at base** — `siteReScript` / `siteReTable` do not exist | RED at base by non-existence. Post-change: measured RED with the right message (§3). |
| AC11 | `gofmt -l host/verifygate/` = 0 bytes; `go vet ./host/verifygate/` rc=0; `go build ./...` rc=0 | **0 bytes**, **rc=0**, **rc=0** | GREEN at base — hygiene guard. Post-change I measured all three still green. |
| **AC12.1** | `./scripts/verify_ail.sh` → **rc=0** | **rc=0** in **7.7 s**; `✓ world package gate PASSED: 9/9 steps performed non-zero work`; `✓ verify gate PASSED: 11 required identities verified, 40 named tests pass` | **rc=0 criterion is LEGITIMATE.** Re-derived, not transcribed. |
| **AC12.2** | `./scripts/verify_go.sh` → **SET comparison, `rc=0` FORBIDDEN as a criterion** | **rc=0**, `grep -c '^--- FAIL'` = **0** — the failing-test set at base is **EMPTY**. `host/verifygate` = `ok` in **both** legs (plain 116.3 s, race 133.7 s). Whole gate ≈ **13 min** wall. | **THE DOC'S BASE PREMISE IS REFUTED — see finding F1.** The criterion stays a SET comparison; the base set is `∅`, with a named flake allowlist. `rc=0` remains forbidden as the criterion because the gate is **flaky**, not because it is always red. |
| AC12.3 | `go test ./host/verifygate/ -count=1` → rc=0 | **rc=0**, 63.8 s; with `-v`: **52 `=== RUN` / 52 `--- PASS` / 0 `--- FAIL` / 0 `--- SKIP`** | GREEN at base — this is **the narrowest gate that can actually fail for this diff**. Legitimate rc=0 criterion. |
| AC13 | **T** with D1 / D2 landed → rc=1, duplicate message, **and** `grep -c 'want exactly 1' /tmp/t` = 0 | **not runnable at base** (no duplicate guard) | RED at base by non-existence. Post-change: measured RED with `want exactly 1` = **0** — the attribution clause HOLDS (§3). |
| AC14 | `grep -rn 'siteReScript\|siteReTable\|canonicalSiteSet' --include='*.go' . \| grep -v floor_raise_inventory_test.go \| wc -l` = **0**, beside control `grep -rn 'regexp.MustCompile' --include='*.go' host/ \| wc -l` ≥ **9** | **0** and **9** — the instrument fires | **GREEN at base and STILL 0 post-change** (measured). It cannot discriminate on its own; its only tooth is mutation **C1**, which §4 F2 shows is a **no-op as the doc writes it**. |
| AC15 | `grep -c 'want exactly 1' <file>` = 2 AND `grep -c 'scriptRows :=\|s8Rows :=' <file>` = 2 | **2** and **2** | GREEN at base — it guards code this sprint does not touch. Its only tooth is mutation **C2**. |

**Nothing in this plan is gated on a command that is red at base for a reason unrelated to the
change.** The two ACs that are red at base (AC4, AC5) are red *because they are the defect*.

### 2.2 Structural facts re-derived first-party (the doc's V8–V13, re-run, not transcribed)

- Both homes hold **6** rows today: `sed -n '/── FLOOR-RAISE COUPLING INVENTORY/,/── END FLOOR-RAISE/p' scripts/verify_ail.sh | grep -cE '^#   [0-9]+\. '` → **6**; `sed -n '/## S8/,/^## /p' design_docs/coding-standards.md | grep -cE '^\| [0-9]+ \| '` → **6**.
- Whole-file `grep -cE '^#   [0-9]+\. ' scripts/verify_ail.sh` → **6** as well, so a whole-file count is **not** a safe proxy for the block count at base. (It becomes one under ARM A4 — §4 F4.)
- §S8 **is** the last `## ` section: `grep -n '^## ' design_docs/coding-standards.md | tail -3` → `73: ## S6`, `80: ## S7`, `90: ## S8`; the file is 122 lines. The test's heading-to-next-`## ` slice runs to EOF. Confirmed independently in the extraction re-run below.
- **V12 re-derived**: applying the doc's two regexes to the two real homes yields the **identical 6-element list in the same order** — `world/<module>.ail`, `packages/world-core/world/<module>.ail`, `scripts/verify_ail.sh`, `scripts/world_package_ready_packet.golden.json`, `docs/SELF_MOD_PUBLISH.md`, `host/verifygate/module_manifest_gate_test.go`. `sets equal: True`, `cardinality equal: True`.
- **V13 re-derived** by reading the file: `floor_raise_inventory_test.go:32 scriptRows := []string{`, `:40` loop, `:43 strings.Count(block, needle) != 1`, `:77 s8Rows := []string{`, `:85` loop, `:88 strings.Count(s8, needle) != 1`, and the two `strings.Contains` sharedNeedle loops at `:58` and `:93`. Both loops use `t.Errorf`, never `t.Fatalf`. Every line number in the doc's V13 is **correct**.
- `repoRoot` is defined once, at `host/verifygate/ail_binary_gate_test.go:27` (`repoRoot = findRepoRoot()`), package-level in `package verifygate`. The new helper needs no new plumbing.
- `host/verifygate` contains **zero** `net.Listen` / `httptest` sites (`grep -rn 'net.Listen\|ListenUnix\|net.Dial\|httptest' host/verifygate/*.go | wc -l` → **0**). This is the measured basis for §7.

### 2.3 `.snap/` and the two deliverables are committable

- `git check-ignore -v .snap/M1/x` → **rc=1** (NOT ignored), against the same-call control `git check-ignore -v .ailang/state/x` → **rc=0**, matching `.gitignore:3:**/.ailang/`. The instrument fires; the zero is a measurement.
- `git check-ignore -v` on both planner deliverable paths → **rc=1** (NOT ignored). The controller can commit them.

### 2.4 The defect, reproduced by me (the doc's V2 / V3, re-run)

```zsh
cp scripts/verify_ail.sh /tmp/w51_backup/verify_ail.sh
perl -0pi -e 's/(#   6\. host\/verifygate\/module_manifest_gate_test\.go)/#   7. host\/store\/some_new_file.go   fabricated coupled site (repro)\n$1/' scripts/verify_ail.sh
shasum -a 256 scripts/verify_ail.sh | cut -c1-8   # 5a1bbe89 -> f97c4ec9   (LANDED)
sed -n '/── FLOOR-RAISE COUPLING INVENTORY/,/── END FLOOR-RAISE/p' scripts/verify_ail.sh | grep -cE '^#   [0-9]+\. '   # 6 -> 7  (EFFECT)
bash -n scripts/verify_ail.sh   # rc=0
# T -> rc=0  RUN=1  PASS=1  FAIL=0.   THE GATE IS BLIND.
```

Same in the §S8 direction: sha `b710a510` → `63208d62`, S8 row count 6 → 7, **T** → rc=0 `--- PASS`.
Both venues restored by `cp`, both sha256 back to base, `git status --porcelain` → 0 lines, control
re-PASSED. **The doc's V2/V3 sha values reproduce exactly on this rig.**

---

## 3. THE PLANNER'S OWN MEASUREMENT — all fifteen arms RUN against a real implementation, not predicted

I applied the design doc's §Implementation sketch (doc lines **235–262** for the declarations,
**266–317** for the test-body additions) into `host/verifygate/floor_raise_inventory_test.go` in
this worktree, measured it, ran every arm, and restored the file byte-identical
(`a00790f1…113f`, porcelain 0, control re-PASSED).

**The sketch as written compiles and is clean:** `gofmt -l host/verifygate/` → **0 bytes**,
`go vet ./host/verifygate/` → **rc=0**, **T** on pristine homes → **rc=0, RUN=1, PASS=1**.
The file goes **97 → 161 lines** (+64 added, 0 removed). Post-change AC greps measured:
AC2 = **3**, AC3 = **2**, duplicate message = **2**, AC14 filtered = **0**, AC15 = **2** / **2**.

| Arm | Landed + effect gate (measured) | **Result, RUN** | **Which assertion fired** — an arm that reds on a different one is a false reproduction |
|---|---|---|---|
| **A1** script-home-only 7th row | sha `5a1bbe89`→`f97c4ec9`; block 6→7; `bash -n` rc=0 | **RED** rc=1, `--- FAIL`=1 | `verify_ail.sh site "host/store/some_new_file.go" absent from the S8 table — the two homes disagree on their site set` **plus** the cardinality Errorf (`7 sites, S8 has 6`). `want exactly 1` = **0**. |
| **A2** §S8-only 7th row | sha `b710a510`→`63208d62`; S8 6→7 | **RED** rc=1, `--- FAIL`=1 | `S8 site "host/store/some_new_file.go" absent from the verify_ail.sh inventory block — …` **plus** the cardinality Errorf (`6 sites, S8 has 7`). |
| **A3** 7th row in BOTH homes | both shas move; both counts 6→7 | **PASS** rc=0, `--- PASS`=1 | none — the legitimate floor-raise path. A design that reds here is unusable. |
| **A4** 7th row OUTSIDE the block | sha `5a1bbe89`→`75a447ab`; **whole-file 6→7** while **block stays 6**; `bash -n` rc=0 | **PASS** rc=0 | none — the extraction really is bounded. |
| **D1** fabricated row TWICE in the script block | sha →`8ad8923c`; block 6→**8**; `bash -n` rc=0 | **RED** rc=1 | `duplicate coupled-site path "host/store/some_new_file.go" in the verify_ail.sh inventory block — a duplicated row is a defect` (`t.Fatalf`). **`want exactly 1` = 0** → AC13's attribution clause HOLDS: the legacy loop did not fire. |
| **D2** fabricated §S8 row TWICE (7 and 8) | sha →`2ee0000a`; S8 6→**8** | **RED** rc=1 | `duplicate coupled-site path … in the S8 table — a duplicated row is a defect`. **`want exactly 1` = 0.** |
| **R1** delete script row 3 | block 6→5; `bash -n` rc=0 | **RED** rc=1 | **TWO** messages: legacy `inventory block site 3 row "#   3. scripts/verify_ail.sh" count=0, want exactly 1` **AND** the new floor `inventory block enumerator matched 5 sites, want >= 6 — instrument failure…`. |
| **R2** delete an §S8 row | S8 6→5 | **RED** rc=1 | Legacy `S8 site N row … count=0, want exactly 1` **AND** the new floor. **Use row 1, not row 3** — see finding F5. |
| **I1** neuter `siteReScript` | test-file sha moves; `go vet` **rc=0** | **RED** rc=1 | `inventory block enumerator matched 0 sites, want >= 6 — instrument failure, not an empty inventory`. |
| **I2** neuter `siteReTable` | test-file sha moves; `go vet` **rc=0** | **RED** rc=1 | `S8 enumerator matched 0 sites, want >= 6 — instrument failure, not an empty inventory`. |
| **N1** neuter set-equality (both membership Errorfs **and** the cardinality Errorf) **+ A1 landed** | `grep -c 'disagree on their site set' <file>` 3→**0**; `go vet` rc=0; script sha →`f97c4ec9` | **GREEN** rc=0, `--- PASS`=1 | none — **the previously-expected `--- FAIL` is GONE.** This is the proof that the set-equality, not the floor, is what makes A1 red. |
| **N2a** neuter the floor **+ I1 + I2 together** | `grep -c 'instrument failure…'` 2→**0** AND `grep -c ZZNEVERMATCHZZ` 0→**2**; `go vet` rc=0 | **GREEN** rc=0 | none — **the vacuous pass the floor exists to prevent**: both sets empty, every membership loop vacuous, cardinality `0 == 0`. |
| **N2b** PAIRED CONTROL: I1 + I2 with the floor **PRESENT** | `grep -c ZZNEVERMATCHZZ` 0→**2** AND `grep -c 'instrument failure…'` stays **2**; `go vet` rc=0 | **RED** rc=1 | `inventory block enumerator matched 0 sites, want >= 6 — instrument failure…`. The SAME mutant as N2a, floor restored → the floor is isolated. |
| **C1** second `siteReScript` in **`host/verifygate/zz_c1_collision_test.go`** (same package) | file exists; **AC14's filtered count 0→1**; `go vet ./host/verifygate/` **rc=1** | **RED — the package stops building**: `vet: host/verifygate/zz_c1_collision_test.go:5:5: siteReScript redeclared in this block` | build failure, before any test result. **The doc's venue is wrong — see F2.** |
| **C2** delete `"#   3. scripts/verify_ail.sh"` from `scriptRows`, then land R1 | `scriptRows` entries 6→**5** while `grep -c 'want exactly 1' <file>` stays **2**; block 6→5; `go vet` rc=0 | **R1 STOPS redding on the per-row message**: `want exactly 1` in the output goes 1→**0** | only the floor fires (`matched 5 sites, want >= 6`) — a strictly weaker verdict. The hardcoded lists really are the independent authority. |

**A live demonstration of why the landed-proof is mandatory.** My first attempt at N2a/N2b used a
patch whose regex escaping was wrong. The mutation **silently did not land** (`grep -c
ZZNEVERMATCHZZ` = **0**, `grep -c 'instrument failure…'` still **2**) and the run reported
`N2b → GREEN` and `N2a → RED` — both **backwards**. The landed-proof caught it; the test result did
not. Every arm below carries a MOVING `grep -c` for exactly this reason.

---

## 4. Planner findings — defects in the design doc's own text, found by running it

A refutation is a success. Five are real, one is a base-premise refutation, and none of them touch
the design DIRECTION (which the quorum settled and this plan does not re-litigate).

### F1 — AC12 part 2's base premise is REFUTED on this rig, this session. (HIGH)

The doc's **V6** and Declared Residual 5 record `./scripts/verify_go.sh` as **rc=1 at pristine base**
with exactly one failing test, `TestHandlerTimeoutKillsTheWholeProcessGroup` (host/broker).
**I re-derived it and measured the opposite:**

```zsh
export PATH=/opt/homebrew/bin:$PATH; export AILANG_BIN=/tmp/ailang-v0300/ailang
./scripts/verify_go.sh > /tmp/vgo 2>&1; echo "rc=$?"        # rc=0
grep -cE '^--- FAIL' /tmp/vgo                                # 0
grep -E '^(ok|FAIL).*host/(verifygate|broker)' /tmp/vgo
# ok  …/host/verifygate  116.349s     (plain leg)
# ok  …/host/broker       87.916s     (plain leg)
# ok  …/host/verifygate  133.660s     (race leg)
# ok  …/host/broker      164.331s     (race leg)
# ✓ go gate PASSED: build clean, plain and race tests pass with pinned AILANG_BIN
```

So the base failing-test set is **empty**, not `{TestHandlerTimeoutKillsTheWholeProcessGroup}`. The
right reading is that this test is a **rig-local FLAKE** (a 100 ms process-group-kill deadline under
parallel-package load), red in the designer's session and green in mine — **not** a deterministic
base red.

**Repair (this plan's criterion supersedes AC12.2's literal base set):**

> **AC12.2′.** Run `./scripts/verify_go.sh` once, after the last milestone, **outside the sandbox**
> (§7). Let `S_after` = the sorted set of test names on `^--- FAIL` lines. The criterion is
> `S_after \ {TestHandlerTimeoutKillsTheWholeProcessGroup}` **== ∅** (the base set I measured this
> session), **AND** `host/verifygate` reports `ok` in **both** the plain and the race leg. **Reading
> `rc` is still FORBIDDEN as the criterion** — not because the gate is always red, but because it is
> flaky: an rc=0 requirement would fail a correct landing whenever the known flake recurs, and an
> rc=1 expectation would fail a correct landing whenever it does not. If any *other* name appears,
> that is a real regression and the sprint stops.

### F2 — mutation C1 is a NO-OP as literally written. (HIGH)

AC14 says "declare a second `siteReScript` in any **other `host/` file**"; the C1 row says "in
another `host/` file", and both predict "the package stops building". **Measured, both placements:**

| C1 venue | `go vet ./host/verifygate/` | `go vet ./host/broker/` | AC14 filtered count |
|---|---|---|---|
| `host/verifygate/zz_c1_collision_test.go` (**same package**) | **rc=1** — `siteReScript redeclared in this block` | n/a | 0 → **1** |
| `host/broker/zz_c1_collision_test.go` (a *different* package) | **rc=0** | **rc=0** | 0 → **1** |

Go redeclaration is package-scoped. A second `siteReScript` "in any other `host/` file" that happens
to be in another package breaks **nothing**, and the arm would report "the collision guard has no
teeth" when the guard is fine. **Repair: C1 MUST create its collision file inside
`host/verifygate/` (package `verifygate`).** The plan's arm table (§6) carries the corrected venue.

### F3 — C1's landed gate ("count moves 1→2") is wrong on both possible readings. (MEDIUM)

Measured post-change, before C1: the **unfiltered** `grep -rn 'siteReScript' --include='*.go' . | wc -l`
is **3** (doc comment, declaration, and the `re = siteReScript` assignment), and AC14's **filtered**
count is **0**. Neither is 1. **Repair: C1's landed gate is AC14's own filtered count moving
`0 → 1` (measured), paired with `go vet ./host/verifygate/` moving `rc=0 → rc=1`.**

### F4 — ARM A4's effect gate must NOT MOVE, which breaks the drill's own protocol. (MEDIUM)

The drill protocol (binding, doc §Mutation Drill) requires every arm to assert its effect with "a
`grep -c` that must **MOVE**". A4's row instead specifies "G2 block count **STAYS** 6" — a gate that
must not move cannot distinguish "the row landed outside the block" from "the `perl` did nothing".
**Repair, measured:** A4's landed/effect gate is the **whole-file** count
`grep -cE '^#   [0-9]+\. ' scripts/verify_ail.sh` moving **6 → 7** (I measured exactly this), *paired
with* the block count staying **6**. The moving gate proves the edit landed; the static one proves it
landed outside the bounds. This is also the operational payoff of the doc's own V8 note.

### F5 — ARM R2's named row is CONFOUNDED by a sharedNeedle. (MEDIUM)

R2 names `| 3 | \`scripts/verify_ail.sh\`` as the row to delete from §S8. That cell is the **only**
place `REQUIRED_VERIFIED` appears in the §S8 slice (the doc's own V9 records `s8=1`). Measured, R2 on
row 3 reds on **three** assertions:

```
:114  S8 site 3 row "| 3 | `scripts/verify_ail.sh`" count=0, want exactly 1
:119  S8 omits "REQUIRED_VERIFIED"                     <-- the confound
:145  S8 enumerator matched 5 sites, want >= 6 …
```

**Repair, measured:** R2 deletes **row 1** (`| 1 | \`world/<module>.ail\``), whose "What moves" cell
contains no sharedNeedle. Result: exactly the two intended messages (`:114` per-row, `:145` floor),
no `S8 omits` line. R1 on the script side is already clean — `REQUIRED_VERIFIED` appears **twice** in
the script block (V9), so deleting row 3 there leaves one and the sharedNeedle loop stays quiet.

### F6 — AC8's removal arms now produce a DUAL verdict, and the AC does not say so. (LOW)

AC8 asks only for `rc=1`. Measured post-change, R1 and R2 red on the legacy per-row `want exactly 1`
**and** on the new `>= 6` floor (5 < 6). That is fine — but a reader who sees only the floor message
could mistake a **lost** per-row authority for a healthy gate. **Repair: AC8 asserts BOTH messages
present** (`grep -c 'want exactly 1' /tmp/t` = 1 **and** `grep -c 'instrument failure' /tmp/t` = 1).
That is precisely the separation mutation **C2** exists to prove, so the two now line up.

### F7 — AC14 and AC15 are GREEN AT BASE and cannot discriminate on their own. (INFO, not a defect)

Both were measured at their required values on the pristine tree (`0` / `9`, and `2` / `2`) and I
measured AC14 **still 0** after the change. They are correctly-designed *non-regression* criteria
whose only teeth are C1 and C2 — which is exactly why the carve-out bound each fix to its own named
mutation. Recorded so the evaluator does not read "green at base" as "vacuous".

### F8 — a `-v`-less package run cannot count skips. (INFO)

`go test ./host/verifygate/ -count=1` prints no `--- SKIP` lines at all, so `grep -c -- '--- SKIP'`
= 0 on a run with skips. AC12.3's SKIP clause is only meaningful with `-v`. Measured with `-v`:
**52 RUN / 52 PASS / 0 FAIL / 0 SKIP**, 63.2 s.

**Authority.** THE DESIGN DOC WINS over this plan, **except** on F1, F2, F3, F4, F5 and F6, where my
first-party measurement showed the doc's literal text would either fail a correct landing (F1),
prove nothing (F2, F3, F4), or red for the wrong reason (F5, F6). On those six, this plan wins.

---

## 5. Milestone plan — 2 milestones, 2 commits, each green at its boundary

Standing rule: **one commit per milestone**, and `go test ./host/verifygate/` must pass at every
boundary. Two milestones is the right granularity for a ~64-line single-file change; a third
milestone would have no diff of its own and is rejected below.

### Milestone M1 — the enumerators, the duplicate guard, and the cardinality floor

**File:** `host/verifygate/floor_raise_inventory_test.go` (the only file) · **LOC: +51, −0** (measured
on my scratch application: 97 → 148 lines)

**Work**

1. Add `"regexp"` to the import block (between `"path/filepath"` and `"strings"` — gofmt order).
2. Add the package-level `var (siteReScript, siteReTable)` block and the `canonicalSiteSet` helper
   **byte-verbatim from the design doc's sketch, doc lines 235–262**, comments included. The comments
   are the declared precondition (`(\S+)` assumes whitespace-free paths; `([^|]+)` + backtick-trim
   assumes a backtick-wrapped §S8 path column) — do not trim them.
3. At the **end** of the test body (after the §S8 sharedNeedle loop at `:92-96`, where both `block`
   and `s8` are in scope), add, from doc lines 266–317: `scriptSites` / `s8Sites` extraction, the two
   duplicate-within-a-home guards (`t.Fatalf`), and the two `>= 6` cardinality floors (`t.Fatalf`).
   **Stop there** — the set-difference is M2.

**Acceptance at the M1 boundary — measured by me on the M1 subset, not predicted**

| Check | Measured on the M1 subset |
|---|---|
| **T** on pristine homes | **rc=0, RUN=1, PASS=1** — the boundary is green |
| `gofmt -l host/verifygate/` | **0 bytes** |
| `go vet ./host/verifygate/` | **rc=0** |
| `grep -c 'instrument failure, not an empty inventory' <file>` (AC3) | **2** |
| `grep -c 'duplicate coupled-site path' <file>` | **2** |
| `grep -c 'disagree on their site set' <file>` (AC2) | **0** — correctly still 0; M2 raises it to 3 |
| **Honest non-claim:** **T** with ARM A1 landed | **rc=0, `--- PASS`** — **M1 does NOT yet close the row's named gap.** M1 buys the anti-vacuity and anti-evasion guards; M2 buys the detection. Stating this is the point of a bisectable boundary. |
| **T** with ARM D1 landed | **rc=1**, duplicate `t.Fatalf`, `want exactly 1` = 0 — the guard is already live at M1 |

**Drill arms run and recorded in M1:** **D1, D2, I1, I2, N2a, N2b, C1**.

**Gates:** AC11 (gofmt / vet / build), AC12.1 `./scripts/verify_ail.sh` rc=0, AC12.3
`go test ./host/verifygate/ -count=1` rc=0. **Snapshot to `.snap/M1/`. → Commit 1.**

---

### Milestone M2 — the symmetric set-difference (the row's point)

**File:** the same one file · **LOC: +13, −0** (148 → 161 lines)

**Work**

1. Append, from doc lines 266–317, the two membership loops (`t.Errorf`, one per direction) and the
   cardinality-mismatch `t.Errorf`. Three `Errorf`s, three occurrences of the frozen substring
   `disagree on their site set`.
2. Nothing else. `scriptRows` / `s8Rows` and both per-row `count == 1` loops stay **byte-unchanged**
   (AC15); `sharedNeedles` stays `strings.Contains`-only (Decision 4, declared residual 2).

**Acceptance at the M2 boundary — measured**

| Check | Measured |
|---|---|
| **T** on pristine homes (AC1 / AC10) | **rc=0, RUN=1, PASS=1, SKIP=0** |
| AC2 `grep -c 'disagree on their site set' <file>` | **3** |
| AC4 — **T** with A1 | **rc=1**, `--- FAIL`=1, names `host/store/some_new_file.go` |
| AC5 — **T** with A2 | **rc=1**, `--- FAIL`=1, names the site |
| AC6 — **T** with A3 | **rc=0**, `--- PASS`=1 |
| AC7 — **T** with A4 | **rc=0**, `--- PASS`=1 |
| AC11 | gofmt 0 bytes, vet rc=0, `go build ./...` rc=0 |
| AC14 / AC15 | filtered count **0** (still), `want exactly 1` **2**, `scriptRows :=|s8Rows :=` **2** |

**Drill arms run and recorded in M2:** **A1, A2, A3, A4, R1, R2, N1, C2** — plus a re-run of D1 and
I1 as carry-forward controls (the M1 guards must still fire with M2 layered on: measured, they do).

**Gates:** AC11, AC12.1, AC12.3, **and** AC12.2′ (the whole-repo Go gate, **controller-run**, §7).
**Snapshot to `.snap/M2/` (cumulative). → Commit 2.**

---

### Seams chosen and seams rejected

| Seam | Verdict | Reason (measured where it matters) |
|---|---|---|
| M1 (floor + dup guard) → M2 (set-difference) | **CHOSEN** | The floor must exist **before** the set-equality, or commit 1 would ship a comparator that passes vacuously the moment an enumerator breaks — exactly the N2a failure. Both boundaries measured green. |
| M2 first, M1 second | **REJECTED** | Opens a one-commit window in which the vacuous-pass class is live. The doc's own Decision 2 names the floor as what makes the set-equality non-vacuous. |
| One single commit for the whole item | **REJECTED** | The standing rule is one commit per milestone with a bisectable boundary, and the two halves have genuinely different teeth (D/I/N2 vs A/R/N1). |
| A third "drill + gates" milestone | **REJECTED** | It would have **no diff**. An empty commit is not a milestone; the drill is distributed across M1 and M2 at the point where each arm's machinery first exists. |
| Splitting M1 into "enumerators" / "dup guard" / "floor" | **REJECTED — inflation** | The enumerators alone are dead code with no assertion; that boundary has no tooth. |
| Touching `design_docs/coding-standards.md` §S8 | **FORBIDDEN** | §S8 is ratification-class (human gate). This sprint READS it and mutates it only inside drill arms that are restored byte-identical. |

---

## 6. The mutation drill — the deliverable, not a nicety

**Per-arm protocol, binding, no shortcuts.** For every one of the fifteen arms, in order:

1. **Restore first.** `cp` the venue back from `/tmp/w51_backup/` so the previous mutant cannot leak
   into this one. **NEVER `git checkout -- <path>`** — it would delete the executor's uncommitted
   deliverable, which is the whole sprint (§7).
2. **Run the pristine control** (**T** → rc=0, `--- PASS`=1) *before* the arm.
3. **Apply the mutant.**
4. **Prove it LANDED, two ways:** (a) `shasum -a 256 <venue> | cut -c1-8` differs from the base
   prefix in the table; **and** (b) the arm's **intended-effect query** — a `grep -c` against the
   *system's own view* that must **MOVE** to the value in the table. Never bytes alone.
5. **Prove the tree still builds:** `go vet ./host/verifygate/` **rc=0** — read this **BEFORE** any
   test result. (The one exception is **C1**, whose entire verdict *is* `go vet` rc≠0.) For
   shell-venue arms also `bash -n scripts/verify_ail.sh` rc=0.
6. **Run T** and record **the message that fired**, not the exit code. An arm that reds on a
   different assertion than the one named in the "must fire" column is a **FALSE REPRODUCTION** and
   must be reported as such, not counted as a kill.
7. **Restore by `cp`**, assert `shasum -a 256` equals the base value, and **re-run the pristine
   control** (rc=0, `--- PASS`=1) *after* the arm.

A `grep -c` that legitimately counts 0 exits rc=1. **Read the COUNT, never the exit code**, and never
wrap it in `|| echo 0`.

### Backups — taken once, before the first arm

```zsh
mkdir -p /tmp/w51_backup
cp host/verifygate/floor_raise_inventory_test.go /tmp/w51_backup/   # the DELIVERABLE — cp only, never git checkout
cp scripts/verify_ail.sh                        /tmp/w51_backup/
cp design_docs/coding-standards.md              /tmp/w51_backup/
shasum -a 256 /tmp/w51_backup/*
```

### The fifteen arms

Venue shorthand: **S** = `scripts/verify_ail.sh` (base sha `5a1bbe89`), **D** = `design_docs/coding-standards.md`
(base sha `b710a510`), **G** = `host/verifygate/floor_raise_inventory_test.go` (base sha `a00790f1`,
post-M2 sha `c84f87bf` on my reference implementation — the executor's will differ if a comment
differs; assert *movement*, never a literal post-change sha).

Effect queries: `BLK` = `sed -n '/── FLOOR-RAISE COUPLING INVENTORY/,/── END FLOOR-RAISE/p' scripts/verify_ail.sh | grep -cE '^#   [0-9]+\. '`;
`S8C` = `sed -n '/## S8/,/^## /p' design_docs/coding-standards.md | grep -cE '^\| [0-9]+ \| '`;
`WF` = `grep -cE '^#   [0-9]+\. ' scripts/verify_ail.sh`.

| # | Arm | Venue | Effect query — MUST MOVE | Milestone | Expected | **THE assertion that must fire** |
|---|---|---|---|---|---|---|
| 1 | **D1** insert the fabricated 7th row `#   7. host/store/some_new_file.go   fabricated coupled site (repro)` **twice** in the block (as `#   7.` and `#   8.`) | S | `BLK` **6→8** | M1 | RED rc=1 | `duplicate coupled-site path "host/store/some_new_file.go" in the verify_ail.sh inventory block` **AND** `grep -c 'want exactly 1' /tmp/t` = **0** (AC13 attribution) |
| 2 | **D2** insert the fabricated §S8 row twice, numbered 7 and 8 | D | `S8C` **6→8** | M1 | RED rc=1 | `duplicate coupled-site path "host/store/some_new_file.go" in the S8 table` **AND** `want exactly 1` = **0** |
| 3 | **I1** neuter `siteReScript` (e.g. `^#   [0-9]+` → `^ZZNEVERMATCHZZ[0-9]+`) | G | `grep -c ZZNEVERMATCHZZ <G>` **0→1**; `go vet` rc=0 | M1 | RED rc=1 | `inventory block enumerator matched 0 sites, want >= 6 — instrument failure, not an empty inventory` |
| 4 | **I2** neuter `siteReTable` the same way | G | `grep -c ZZNEVERMATCHZZ <G>` **0→1**; `go vet` rc=0 | M1 | RED rc=1 | `S8 enumerator matched 0 sites, want >= 6 — instrument failure, not an empty inventory` |
| 5 | **N2b** PAIRED CONTROL — I1 **and** I2 together, floor **present** | G | `grep -c ZZNEVERMATCHZZ` **0→2** AND `grep -c 'instrument failure, not an empty inventory' <G>` stays **2** | M1 | RED rc=1 | the `inventory block enumerator matched 0 sites` Fatalf. **Run N2b BEFORE N2a** — it is the control that gives N2a its meaning. |
| 6 | **N2a** I1 **and** I2 together, **and** delete both floor `t.Fatalf` blocks | G | `grep -c ZZNEVERMATCHZZ` **0→2** AND `grep -c 'instrument failure, not an empty inventory' <G>` **2→0** | M1 | **GREEN rc=0** | **none** — the vacuous pass. The proof is the *absence* of N2b's red under the identical enumerator mutant. |
| 7 | **C1** create `host/verifygate/zz_c1_collision_test.go` (**package `verifygate`** — F2) declaring a second `var siteReScript = regexp.MustCompile(\`^collision\`)` | new file | **AC14's filtered count 0→1** (F3), i.e. `grep -rn 'siteReScript\|siteReTable\|canonicalSiteSet' --include='*.go' . \| grep -v floor_raise_inventory_test.go \| wc -l` | M1 | RED — **build failure** | `go vet ./host/verifygate/` **rc=1** with `siteReScript redeclared in this block`. Clean up with **`rm`** (the file is untracked; `git checkout` cannot remove it). |
| 8 | **A1** the fabricated 7th row in the script block ONLY | S | `BLK` **6→7**; `bash -n` rc=0 | M2 | RED rc=1, `--- FAIL`=1 | `verify_ail.sh site "host/store/some_new_file.go" absent from the S8 table — the two homes disagree on their site set` (the cardinality Errorf also fires; the **named** one is the membership Errorf) |
| 9 | **A2** the fabricated 7th row in §S8 ONLY | D | `S8C` **6→7** | M2 | RED rc=1, `--- FAIL`=1 | `S8 site "host/store/some_new_file.go" absent from the verify_ail.sh inventory block — the two homes disagree on their site set` |
| 10 | **A3** the fabricated 7th row in BOTH homes, consistently | S + D | `BLK` **6→7** AND `S8C` **6→7** | M2 | **PASS rc=0**, `--- PASS`=1 | **none.** A red here means the legitimate floor-raise path is broken and the design is unusable. |
| 11 | **A4** a `#   7. …` row **after** the END marker | S | **`WF` 6→7** (the MOVING gate — F4) **paired with** `BLK` staying **6**; `bash -n` rc=0 | M2 | **PASS rc=0** | **none** — the extraction is bounded. |
| 12 | **R1** delete `#   3. scripts/verify_ail.sh …` from the block | S | `BLK` **6→5**; `bash -n` rc=0 | M2 | RED rc=1 | **BOTH** (F6): `inventory block site 3 row "#   3. scripts/verify_ail.sh" count=0, want exactly 1` **AND** `inventory block enumerator matched 5 sites, want >= 6` |
| 13 | **R2** delete **`\| 1 \| \`world/<module>.ail\``** from §S8 — **row 1, not row 3** (F5) | D | `S8C` **6→5** | M2 | RED rc=1 | **BOTH**: `S8 site 1 row "\| 1 \| \`world/<module>.ail\`" count=0, want exactly 1` **AND** `S8 enumerator matched 5 sites, want >= 6`. **`grep -c 'S8 omits' /tmp/t` must be 0** — a non-zero there means the confounded row was used. |
| 14 | **N1** delete the two membership `t.Errorf`s **and** the cardinality-mismatch `t.Errorf`, **and** land A1 | G + S | `grep -c 'disagree on their site set' <G>` **3→0** AND `BLK` **6→7**; `go vet` rc=0 | M2 | **GREEN rc=0**, `--- PASS`=1 | **none** — A1's expected `--- FAIL` is GONE. The cardinality Errorf **must** be neutered too: under A1 the counts differ (7 vs 6), so leaving it in would mask the proof. |
| 15 | **C2** delete `"#   3. scripts/verify_ail.sh",` from the `scriptRows` literal, then land **R1** | G + S | `scriptRows` entries **6→5** while `grep -c 'want exactly 1' <G>` stays **2**; `BLK` 6→5; `go vet` rc=0 | M2 | RED rc=1, but with a **weaker verdict** | `grep -c 'want exactly 1' /tmp/t` goes **1→0** — R1 stops redding on the per-row message and is caught only by the floor. That delta is the proof that the hardcoded lists are the independent authority. |

**Every one of these fifteen results was measured by me in this worktree (§3).** The executor's job
is to reproduce them, not to discover them; a divergence is a finding to report, not a number to
adjust.

---

## 7. Sandbox posture — which gates the executor may run, and which it may not

The executor runs under **`codex exec --sandbox workspace-write`**, which **denies loopback socket
binds**. I asked of every gate in this plan whether it binds a socket or needs the network.

| Gate | Binds a socket / needs network? | Verdict |
|---|---|---|
| **T** (`go test ./host/verifygate/ -run '^TestFloorRaise…$'`) | **No.** Measured: `grep -rn 'net.Listen\|ListenUnix\|net.Dial\|httptest' host/verifygate/*.go \| wc -l` → **0**. The test reads two files off disk. | **EXECUTOR GATE.** Sandbox-safe; the executor's result is of record. |
| `go test ./host/verifygate/ -count=1` (AC12.3) | **No** — same package, no socket site anywhere in it. | **EXECUTOR GATE.** Sandbox-safe. |
| `gofmt -l`, `go vet ./host/verifygate/`, `go build ./...` (AC11) | No. | **EXECUTOR GATE.** |
| every `grep -c` / `sed` / `shasum` / `bash -n` in the drill | No. | **EXECUTOR GATE.** |
| `./scripts/verify_ail.sh` (AC12.1) | No socket bind; it drives the pinned `ailang` binary against local files. Measured rc=0 in 7.7 s. | **EXECUTOR GATE**, but see the note below on the Observatory stderr line. |
| **`./scripts/verify_go.sh` (AC12.2′)** | **YES.** It runs `go test ./...` (plain **and** `-race`), which includes `host/daemon` — and `host/daemon/daemon.go:634` is `ln, err := net.Listen("tcp", addr)`. `host/broker` and `host/daemon` tests both bind. | **`UNINFORMATIVE UNDER SANDBOX` — a CONTROLLER gate, never an executor obligation.** If the executor runs it anyway, the result **MUST** be labelled `UNINFORMATIVE UNDER SANDBOX` in the sprint log: reported, never silently dropped, and never read as a pass or a fail. **The controller re-runs it outside the sandbox; the controller's numbers are the ones of record.** Budget ≈13 min. |

**No acceptance criterion in this sprint needs `gh` or the network.**

**Observatory noise, not a regression:** every `ailang` invocation prints
`Observatory: 513MB (DB=513MB WAL=0MB) — running retention cleanup` on stderr. `verify_ail.sh` is
rc=0 with it present (measured this session). Do not "fix" it.

---

## 8. Executor constraints — the executor CANNOT commit

- **NO git write operations.** No `add`, `commit`, `stash`, `checkout`, `restore`, `reset`, `branch`,
  `worktree`, `git mv`. A linked worktree's `.git` is excluded by the sandbox, so a git write will
  fail — and `git checkout -- host/verifygate/floor_raise_inventory_test.go` would **destroy the
  entire deliverable**, which is uncommitted by construction. Git **reads** (`status`, `diff`,
  `show`, `log`, `check-ignore`, `ls-files`) are required and allowed. **The controller builds the
  commits.**
- **Per-milestone cumulative snapshots.** After M1: `mkdir -p .snap/M1/host/verifygate && cp
  host/verifygate/floor_raise_inventory_test.go .snap/M1/host/verifygate/`. After M2: the same into
  `.snap/M2/` (cumulative — M2's snapshot is the M2 state of the same one file). This is how the
  controller reconstructs **one commit per milestone** from an uncommitted tree. Measured: `.snap/`
  is **not** gitignored (`git check-ignore -v .snap/M1/x` → rc=1, control `git check-ignore -v
  .ailang/state/x` → rc=0), and dot-prefixed directories are skipped by both the Go tool and
  `gofmt`, so snapshots cannot red a gate.
- **All restores use `cp` from `/tmp/w51_backup/`**, verified by `shasum -a 256`. C1's collision file
  is removed with `rm`.
- **Files this sprint MAY touch:** `host/verifygate/floor_raise_inventory_test.go` (MODIFY, +64 / −0)
  — and, transiently and only inside restored drill arms, `scripts/verify_ail.sh`,
  `design_docs/coding-standards.md`, and `host/verifygate/zz_c1_collision_test.go` (CREATE then `rm`,
  arm C1 only).
- **Files this sprint MUST NOT touch, at any point that survives an arm's restore:**
  `design_docs/coding-standards.md` §S8's own text (**ratification-class, human gate** — the drill
  mutates and restores it; it is never edited as a deliverable); `scripts/verify_ail.sh` as a
  deliverable; `scripts/verify_go.sh`; any sibling `host/verifygate/*_test.go`; `go.mod` / `go.sum`;
  **`tools/launchd/*` (frozen core, FLEET-owned)**; any `.ail` file.
- **Message text is FROZEN.** The drill is keyed on these exact substrings; a re-wording silently
  voids the arms that assert them:
  `duplicate coupled-site path %q in the verify_ail.sh inventory block — a duplicated row is a defect` ·
  `duplicate coupled-site path %q in the S8 table — a duplicated row is a defect` ·
  `inventory block enumerator matched %d sites, want >= 6 — instrument failure, not an empty inventory` ·
  `S8 enumerator matched %d sites, want >= 6 — instrument failure, not an empty inventory` ·
  `absent from the verify_ail.sh inventory block — the two homes disagree on their site set` ·
  `absent from the S8 table — the two homes disagree on their site set` ·
  `site-set cardinality mismatch: verify_ail.sh has %d sites, S8 has %d — the two homes disagree on their site set`.
  Do not rename `siteReScript`, `siteReTable` or `canonicalSiteSet` — **AC14** and **C1** are keyed on
  those identifiers.
- **Assert on MESSAGE SUBSTRINGS, never on line numbers.** The line numbers in §3 (`:129`, `:136`,
  `:142`, `:145`, `:150`, `:155`, `:159`) are from **my** reference implementation; a one-line comment
  difference moves all of them.
- **Do not widen any gate to `go test ./...`.** The narrowest gate that can fail for this diff is
  `go test ./host/verifygate/ -count=1`.
- **Authority:** the design doc wins over this plan **except** on F1–F6 (§4).

---

## 9. Task breakdown, velocity and cost centres

Recent comparable single-file test-hardening rows in this repo: `P49` (row 49, canary fence),
`P47` (row 47, ~113 new lines in one `host/verifygate` test), `P43` (row 43 — the very inventory
this row hardens). The measured band for a one-file `host/verifygate` test change with a full drill
is **0.1–0.25 d**, dominated by the drill and the gates, not the typing.

| Task | Milestone | Est. | Output |
|---|---|---|---|
| T1 — import, the two regexes, `canonicalSiteSet`, dup guards, floors (doc lines 235–262 + the first half of 266–317) | M1 | 0.15 h | +51 lines |
| T2 — drill arms **D1, D2, I1, I2, N2b, N2a, C1** with full landed/vet/restore protocol | M1 | 0.30 h | 7 recorded arms |
| T3 — M1 gates (gofmt, vet, build, `verify_ail.sh`, `go test ./host/verifygate/`) + `.snap/M1/` | M1 | 0.15 h | boundary green |
| T4 — the three set-difference `t.Errorf`s (rest of doc lines 266–317) | M2 | 0.05 h | +13 lines |
| T5 — drill arms **A1, A2, A3, A4, R1, R2, N1, C2** + D1/I1 carry-forward controls | M2 | 0.35 h | 10 recorded arms |
| T6 — M2 gates + `.snap/M2/` + the AC1–AC15 sweep with base results attached | M2 | 0.15 h | boundary green |
| Contingency (a false reproduction to diagnose, per §3's own live example) | — | 0.20 h | — |
| **Total** | | **≈1.35 h ≈ 0.15 d** | matches the doc's ~0.1 d, with the drill priced honestly |

**Cost centres, measured this session:**

| Command | Wall clock |
|---|---|
| **T** (scoped test) | **0.31 s** — the drill is cheap; 15 arms ≈ 1 min of test time |
| `go test ./host/verifygate/ -count=1` | **63.8 s** (63.2 s with `-v`) — run once per milestone boundary, not per arm |
| `./scripts/verify_ail.sh` | **7.7 s** |
| `go vet ./host/verifygate/` | < 1 s |
| **`./scripts/verify_go.sh`** | **≈13 min** (`host/verifygate` alone is 116 s plain + 134 s race). **Run ONCE, by the controller, outside the sandbox.** |

---

## 10. Risks

| Risk | Likelihood | Mitigation |
|---|---|---|
| A mutant silently fails to land and the arm reports a **backwards** result | **measured-real — it happened to me this session** (§3, N2a/N2b) | The step-4 moving `grep -c` is mandatory and is what caught it. Never read a test result before the landed-proof. |
| The executor runs C1 in `host/broker` per the doc's literal wording and reports "the guard has no teeth" | **measured-real** (F2) | C1's venue is pinned to `host/verifygate/` in §6 arm 7, with the exact filename. |
| `git checkout --` used to restore, destroying the uncommitted deliverable | recurring class in this repo | §8: `cp`-from-`/tmp/w51_backup/` only, sha-verified. |
| `./scripts/verify_go.sh` reds on `TestHandlerTimeoutKillsTheWholeProcessGroup` and the run is read as a regression | **measured-real flake** — red in the designer's session, **green in mine** | AC12.2′ (F1): a named single-test allowlist, set comparison, never `rc`. Any *other* name stops the sprint. |
| The executor "also fixes" §S8 to make a drill arm green | low, but §S8 is one `cp` away from the code being mutated | §8: §S8 is ratification-class. Its base sha `b710a510…8329` is asserted after every arm. |
| The executor promotes `sharedNeedles` to `count == 1` as a "bonus" | low | Decision 4 measured it a **false RED** (`REQUIRED_VERIFIED` 2× in the block, `EXACT_TOTAL_VERIFIED` 2× in both). Declared residual 2. Out of scope. |
| A `-run` selector matching nothing exits **0** with `no tests to run` — a false pass | recurring class | AC1's run-existence floor: assert `=== RUN` ≥ 1 and `--- PASS` ≥ 1 in the same call, and that `no tests to run` is absent. |
| Ambient toolchain drifts back onto `verify_go.sh`'s deny-list (`go1.26.0`–`go1.26.5`) | low — ambient is `go1.26.6` today | The gate FATALs loudly and names the fix. Do not pre-emptively pin `GOTOOLCHAIN=go1.25.6` from older plans (§1). |
| `AILANG_BIN` unset → a large slice of `host/verifygate` fails loudly | recurring | It is a **base condition**, not a regression. Export it on every command. |

---

## 11. Out of scope / owed elsewhere

- **No edit to `design_docs/coding-standards.md` §S8's text.** Ratification-class (human gate). If a
  future row wants the two homes merged, or §S8 to carry an explicit cardinality statement, that is a
  parked human decision.
- **No new package, no new script, no new CI job, no `.ail`, no Z3 contract, no kernel change.** The
  §Conflict Surface audit (doc §Conflict Surface, bound by AC14) found nothing to extend; I
  re-derived its two headline numbers (**0** proposed symbols repo-wide against a same-scope control
  of **9** `regexp.MustCompile` in `host/`).
- **No promotion of `sharedNeedles` to counts** (Decision 4, declared residual 2).
- **The symmetric-wrongness residual stands**: a fabricated 7th row added to BOTH homes consistently
  passes (arm A3, measured PASS) and is indistinguishable from a legitimate floor-raise. Inherent to
  the two-homes-must-agree model; it is the price of AC6.
- **No fix for `TestHandlerTimeoutKillsTheWholeProcessGroup`.** F1 upgrades what is known about it
  (flaky, not deterministically red) — that is a **queue note**, not this sprint's work.
- **S7 is not engaged**: nothing user-facing is added, so no `docs/QUICKSTART.md` change.
- If the sprint uncovers a genuinely separate defect, **FILE it** as a "for the queue, not this
  sprint" note in the sprint log. Do not absorb it.

---

## 12. Definition of done

1. **Two commits** on `dev`, reconstructed by the controller from `.snap/M1/` and `.snap/M2/`, each
   individually green on `go test ./host/verifygate/ -count=1`, `gofmt -l host/verifygate/`,
   `go vet ./host/verifygate/`, `go build ./...` and `./scripts/verify_ail.sh`.
2. **AC1–AC15 all recorded with their measured base result from §2.1 attached**, and AC12.2 replaced
   by **AC12.2′** (F1).
3. **All fifteen mutation arms recorded**, each with: the moving landed-proof, the `go vet` rc read
   *before* the test result, **the message that fired**, and a sha-verified `cp` restore with the
   pristine control re-run on both sides.
4. `git status --porcelain` shows **only** `host/verifygate/floor_raise_inventory_test.go` modified
   plus the two `.snap/` trees. Both drill venues back at their base sha256
   (`5a1bbe89…475a`, `b710a510…8329`); `host/verifygate/zz_c1_collision_test.go` **gone**.
5. `git status --porcelain -- tools/launchd/` **empty** (no driver drift — a red there means "the
   fleet must commit", never "absorb it").
6. Findings **F1–F8** carried into the iteration log, with **F1** flagged as an owed **doc
   amendment** to the design doc's V6 / Declared Residual 5 rather than silently absorbed.
