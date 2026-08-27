# Sprint plan — `w-toolchain-pin-normalizer-accepts-malformed-gotoolchain` (queue row 45)

**Design doc**: `design_docs/planned/w-toolchain-pin-normalizer-accepts-malformed-gotoolchain.md`
(663 lines, committed at `f3790e4`; two quorum rounds, closed under the ratified
narrow-refinement carve-out).
**Planner**: mission-control iteration 134, `claude-opus-5` sprint-planner.
**Base commit**: `f3790e4`, `git status --porcelain` = **0 lines** at planning entry AND exit.
**Worktree**: the isolated worktree the controller assigns. The executor performs **NO git
write operations**; the controller commits.
**Size**: ~0.15 d, 3 milestones, ONE file modified
(`host/verifygate/toolchain_pin_gate_test.go`, ~+55/−20). Zero files created, zero `.ail`
touched. `ci.yml` and `run.sh` ship **byte-untouched** — they are mutation venues only.
**Milestone count**: **3** (M1 / M2 / M3), unchanged from the design doc's decomposition. No
reason to re-cut was found.

---

## 0. Rig rules — read before the first command

1. **Every** Bash call starts `export PATH=/opt/homebrew/bin:$PATH`, or `go`/`gh` fail rc=127
   and you will misread that as a broken toolchain.
2. The shell is **zsh**. Quote glob-shaped flag values (`--include='*.go'`). Never `echo` to
   inspect bytes — use `printf '%s'` or `od -c`. zsh does **not** word-split an unquoted
   variable.
3. **`PIPESTATUS` DOES NOT EXIST IN ZSH.** Planner-measured this session:
   `... | tail -6; echo ${PIPESTATUS[0]}` → `(eval):15: PIPESTATUS[0]: parameter not set`, and
   the exit code is silently lost. Capture rc directly:
   `go test … > /tmp/out 2>&1; rc=$?` then read `/tmp/out`. zsh's own array is
   `$pipestatus[1]` (lowercase, 1-indexed) if you must pipe.
4. **`AILANG_BIN=/tmp/ailang-v0300/ailang`** on every whole-package Go invocation. See finding
   **D1** — the design doc's AC7 omits it and is red at base without it.
5. **Restores use `cp` from a `/tmp` backup, never `git checkout --`.** Milestone work is
   uncommitted by construction; `git checkout -- host/verifygate/toolchain_pin_gate_test.go`
   would delete M1/M2/M3 with no undo. Verify every restore by **sha256 equality**, not by eye.
6. `run.sh` is mode `-rwxr-xr-x`. `cp` preserves content but you must re-assert the exec bit if
   you ever restore across filesystems (`ls -l …/run.sh | cut -c1-12` → `-rwxr-xr-x`); Test B
   asserts `info.Mode()&0o111 != 0`.
7. Every `ailang` invocation prints `Observatory: NNNMB (warn threshold: 200MB)` on stderr with
   rc=0. **Not** a regression. `./scripts/verify_go.sh` legitimately emits exactly 2
   `WARNING: DATA RACE` lines from its own known-positive control. **Not** a finding.
8. **`git status --porcelain` will NOT be 0 during the sprint** (the deliverable is
   uncommitted). Integrity of the *mutation venues* is checked by sha256 against the table in
   §2, never by porcelain. Porcelain-0 applies only to `ci.yml` / `run.sh` / `go.mod`, which
   must be byte-identical to base at every milestone boundary.

**Base sha256 ledger (planner-measured at `f3790e4`, this session):**

| Path | sha256 | Role |
|---|---|---|
| `host/verifygate/toolchain_pin_gate_test.go` | `cb576cfb654fefbd5bda7fc8a2a84220964b5c662d7a64d41b4a087040ed6031` | THE deliverable — this one changes |
| `.github/workflows/ci.yml` | `9c16e64ca28c58f889a8be122a1b2f48c7d240f2a559fba6faf312ec5166239c` | mutation venue — must match at every boundary |
| `design_docs/verification/w-race-gate-blindspot/run.sh` | `b80109aa57882a3d261757fd268888cf84ea0326180bb34969ce2ce9d149d0bb` | mutation venue — must match at every boundary |
| `go.mod` | `7a2983617bb9fc33747f664564fe8d8ab54fc3a177ec4dfb8c61b29ba79a7e52` | read-only |

The `run.sh` sha256 is byte-identical to the one the design doc records at P23
(`b80109aa…d0bb`) — an independent corroboration that this plan and the doc measured the same
file.

Backup once, at sprint start:

```bash
export PATH=/opt/homebrew/bin:$PATH
mkdir -p /tmp/p45-backup
cp .github/workflows/ci.yml /tmp/p45-backup/ci.yml
cp design_docs/verification/w-race-gate-blindspot/run.sh /tmp/p45-backup/run.sh
shasum -a 256 /tmp/p45-backup/ci.yml /tmp/p45-backup/run.sh   # must match the table above
```

---

## 1. Planner findings — where this plan departs from the design doc, and why

The design doc is unusually good: every conflict claim it makes is TRUE (§1.7), and its two
headline mutation claims reproduce exactly. Six defects were found, four of them in the
acceptance criteria and mutation table. **Where this plan and the doc disagree, this plan's
measured numbers win**; everywhere else the doc wins, and its `## Quorum round 2` section wins
over its own earlier sections.

### D1 — HIGH — AC7's package gate is RED at base without `AILANG_BIN`

AC7 asserts `go test ./host/verifygate/ -count=1` → **rc=0 on the unmutated tree**.

**Measured at `f3790e4`, tree clean:** rc=**1**, **17** `--- FAIL` lines, every message
reading `AILANG_BIN is unset — the shim arms need the pinned released delegate to run the real
gate` (`ail_binary_gate_test.go:186/237/311/317/321/325/347/348/409/456/506/618` and
`module_manifest_gate_test.go:172/215/241/257/286`). With
`AILANG_BIN=/tmp/ailang-v0300/ailang` the same command is rc=**0** in 30.5 s.

**Why it matters:** an executor running AC7 verbatim reports a 17-failure RED against a correct
landing, and the plausible next move — teaching those tests to skip — is precisely the
false-green class this whole mission line exists to kill (they `t.Fatalf` on purpose).
**This is the SECOND consecutive iteration in which a design doc has shipped this exact
error** (iteration 133's planner recorded it as its own D1 for row 44). It is now a standing
rig fact, not a one-off.

**Resolution:** every whole-package invocation in this plan carries the export. The scoped
text-scanning runs (`-run 'TestGoToolchainPins…|TestMiscompileInstrument…'`) need **no**
`AILANG_BIN` — planner-measured rc=0 unexported — so they stay unexported and their
lane-independence stays visible.

### D2 — HIGH — mutation **M8 cannot fail for the milestone it is attached to**

The doc's M8 (append a second `#!/usr/bin/env bash`) predicts *"Test B RED via the shebang
floor, now `Fatalf`: `exact shebang count=2, want 1`, single failure message, zero downstream —
the sharpened prediction that distinguishes the shipped `Fatalf` from today's measured
cascade"*.

**Measured at base (Errorf, i.e. M2 NOT applied):** rc=1, **exactly ONE** message —
`toolchain_pin_gate_test.go:209: instrument failure: <run.sh> exact shebang count=2, want 1` —
and **zero** downstream messages. The predicted post-M2 output is **byte-identical to the
pre-M2 output**. Nothing downstream of `:209` depends on shebang count, so `Errorf` and
`Fatalf` are indistinguishable on this input. M8 as written is a vacuous discriminator: it
would pass whether or not M2 landed.

**Resolution:** M8 is retained as a *floor-existence* arm (it proves the shebang floor fires),
explicitly labelled as such and **not** as evidence of fatality. A new arm **M8′** is added,
measured at base by the planner, that *does* discriminate — see §4.

### D3 — HIGH — AC5 asserts the absence of two messages that never fire in that arm

AC5 requires the M7 arm's output to contain the control floor once and **"ZERO of Test B's
downstream assertion messages (assignment-count, PINNED-equality, exec-bit, guard-message)"**.

**Measured at base** (M7 = `sed 's/PINNED="/PINHOLE="/'`, landed-proof `grep -c 'PINNED='` 1→0
with same-call control `grep -c 'KNOWN_BAD='` = 1): rc=1, **exactly 3** messages —

```
toolchain_pin_gate_test.go:199: instrument failure: <run.sh> does not contain known-positive control "PINNED="
toolchain_pin_gate_test.go:222: <run.sh>: PINNED assignment count=0, want 1
toolchain_pin_gate_test.go:247: PINNED="", want go.mod floor "go1.26.6"
```

The **exec-bit** and **guard-message** assertions do not fire in this arm and never could: the
file keeps its mode and keeps `INSTRUMENT FAILURE: the PINNED toolchain`. Asserting their zero
is a dead grep — the mission's own "an absence you did not make observable is a claim, not a
fact" rule, applied to the doc that states that rule.

**Resolution:** AC5 becomes a **measured delta**: message count `3 → 1`, with the *two
observable* downstream needles (`PINNED assignment count=`, `, want go.mod floor `) asserted to
go `1 → 0` each, and the base 3-message run standing as the same-session known-positive control
that proves those needles ARE observable. Exec-bit and guard-message are dropped from the
assertion.

### D4 — MEDIUM — M4 and M9(b) are RED at base and stay RED; neither is evidence M1 landed

Both are correct arms, but the mutation table's framing invites an executor to read their RED
as proof of the new code.

- **M4** (`s/GOTOOLCHAIN:/GO_TOOLCHAIN:/g`): measured at base rc=1, **1** message,
  `toolchain_pin_gate_test.go:108: ci.yml: GOTOOLCHAIN keyed-line count=0, want 2 (one per
  enumerated expected job)`. Identical before and after. It is a **regression arm** — AC4's own
  wording ("the zero-extraction floor stays *reachable across the extraction split*") is the
  right reading; the table row is not.
- **M9(b)** (`go1.26.6` → `go1.25.6` at the first site): measured at base rc=1, **2** messages,
  both target strings present (`toolchain pins disagree: GOTOOLCHAIN=[go1.25.6 go1.26.6] …` and
  `go.mod floor="go1.26.6" disagrees with ci.yml toolchain pin="go1.25.6"`). It is a
  **known-positive control** for M9(a)'s and AC3's absence greps, and its output is unchanged by
  this sprint. Correct as designed; labelled here so its RED is not misread.

**Resolution:** §4's table carries a `base` column with these measured values, and each arm is
typed `NEW TEETH` / `REGRESSION` / `CONTROL`.

### D5 — MEDIUM — needle residue the doc's "Needle discipline" note does not close

The note is **correct and sufficient for the two validator messages**: it names the shipped
literals `is not an allowed standard Go toolchain pin` and `is a toolchain-selection mode, not a
pin`, both of which match the Decision sketch byte-for-byte. Verified.

What it does **not** close is the reverse direction. The doc's own §Thesis and §"The finding in
one paragraph" still carry the round-1 framing — *"a value the Go runtime itself refuses"*,
*"the Go runtime itself refuses with `go: invalid GOTOOLCHAIN`"* — as **premise prose** (true;
P3 measures it). Round-2 R1's whole point was to keep that claim OUT of the shipped message.
An executor who writes the `Fatalf` text from the thesis rather than the sketch re-introduces
exactly the over-claim the revision removed, and the quorum-round-2 record would then not match
the code.

**Resolution — the shipped message text is FROZEN.** M1's `requireToolchainNamePin` must emit,
byte-exactly:

```
%s: %s=%q is a toolchain-selection mode, not a pin; only a bare toolchain name (e.g. go1.26.6) pins
%s: %s=%q is not an allowed standard Go toolchain pin; this repository requires a bare standard toolchain version accepted by go/version.IsValid (for example go1.26.6)
```

and the file must contain **zero** occurrences of `the Go runtime itself refuses`, `is not a
valid Go toolchain name`, and `the Go runtime` in any `t.Fatalf`/`t.Errorf` string. AC7 gains
that check with a same-call known-positive control.

### D6 — LOW — every line number in the doc is correct today and every one of them MOVES

Verified at `f3790e4`: `:13` `normalizeToolchainPin`, `:25` `pinValues`, `:30`/`:45`/`:244` its
three call sites, `:78`/`:128` Test A's fatal floors, `:108` keyed-count, `:120` agreement,
`:150` floor comparison, `:199`/`:209` Test B's two `t.Errorf` instrument floors, `:263`
`saw_pinned_ok` count, `:318`–`:407` row 44's wiring test. All match.

M1 inserts ~35 lines above Test B, so by the time M2 runs, `:199` is near `:233` and `:209`
near `:243`. **Never `sed -i '' '199s/…'`.** Every edit and every assertion in this plan is
content-anchored. The only positional `sed` commands appear in the mutation drill, and they act
on `ci.yml` / `run.sh`, which this sprint never edits — their line numbers are stable by
construction, and each still carries a content assertion.

### 1.7 — the design's conflict claims, re-derived first-party (all four HOLD)

| Doc claim | Planner's first-party check | Result |
|---|---|---|
| Disjoint from row 44's wiring test at the file end | `TestMiscompileInstrumentStepIsGatedInCI` occupies `:318`–`:407`; this item edits `:13`–`:34`, the `:104` region, `:199`, `:209`, `:244`, and inserts before `:268`. Read in full. | **HOLDS** — disjoint |
| `TestReproModuleFloorStaysBelowKnownBadToolchains` unaffected | It calls `moduleGoFloor(repro/go.mod)` and `shellAssignmentValues(…, "KNOWN_BAD")` only. `grep -n '^go ' …/repro/go.mod` → `11:go 1.22`; `canonicalizeVersionPin("1.22")` = `"go1.22"` = `normalizeToolchainPin("1.22")`. Byte-identical. | **HOLDS** |
| Row 50's `shellAssignmentValues` byte-untouched | Defined `:167`–`:179`; no milestone touches it. | **HOLDS** |
| `TestCanaryDeclaresPositiveArmOnly`'s `GOTOOLCHAIN` fence is over a DIFFERENT file | `grep -c 'GOTOOLCHAIN' host/store/toolchain_canary_test.go` → **0** with same-call control `grep -c 'func Test'` → **1** (proves the file was read). The fence's subject is `host/store/toolchain_canary_test.go`; M1 adds `GOTOOLCHAIN` tokens only to `host/verifygate/toolchain_pin_gate_test.go`. | **HOLDS** |

**One check the doc did not make, added here.** Is anything in the repo a *fence over the
deliverable file itself* (a count/enumeration that a new function would trip)?
`grep -rn 'toolchain_pin_gate_test' --include='*.go' --include='*.sh' --include='*.yml' .`
excluding `design_docs/` → **0 hits**, with same-call known-positive control
`grep -rn 'ail_binary_gate_test' …` → **1 hit** (`toolchain_pin_gate_test.go:56`, the doc-comment
cross-reference), proving the grep is live. Neither `evidence_manifest_gate_test.go` nor
`floor_raise_inventory_test.go` enumerates `host/verifygate/*.go` or `design_docs/planned/`
(both read first-party). **No fence. Adding three functions and two artifact files is safe.**

---

## 2. Base-green ledger — every gate command in this plan, measured on the PRISTINE tree

All run in the planner's detached worktree at `f3790e4`, `git status --porcelain` = 0 before
and after, `go version go1.26.6 darwin/arm64`, `AILANG v0.30.0` at `/tmp/ailang-v0300/ailang`.

| # | Command | Base result |
|---|---|---|
| G1 | `go build ./...` | **rc=0** |
| G2 | `gofmt -l host/verifygate/` | **empty output, rc=0** |
| G3 | `go vet ./host/verifygate/` | **rc=0**, no output |
| G4 | `go test ./host/verifygate/ -run 'TestGoToolchainPinsAgreeAndMatchJobList\|TestMiscompileInstrumentProbesPinnedToolchain' -count=1 -v` (unexported) | **rc=0**, 2 `=== RUN`, 2 `--- PASS`, `ok … 0.325s` |
| G5 | `go test ./host/verifygate/ -run '^TestNoSuchPinGateZZZ$' -count=1 -v` (AC1's nonsense control) | rc=0, `testing: warning: no tests to run`, `ok … [no tests to run]`, **0** `=== RUN` |
| G6 | `AILANG_BIN=/tmp/ailang-v0300/ailang go test ./host/verifygate/ -count=1` | **rc=0**, `ok … 30.5s` — **and rc=1 with 17 FAIL without the export (D1)** |
| G7 | `AILANG_BIN=/tmp/ailang-v0300/ailang ./scripts/verify_go.sh` | **rc=0**, 116 lines, all 13 host packages `ok`, banner `✓ go gate PASSED: build clean, plain and race tests pass with pinned AILANG_BIN` |
| G8 | `AILANG_BIN=/tmp/ailang-v0300/ailang ./scripts/verify_ail.sh` | **rc=0**, `✓ verify gate PASSED: 11 required identities verified, 40 named tests pass` — **the unchanged floor**; this sprint touches zero `.ail`, so any movement here is a finding, not a result |
| G9 | `git check-ignore -v design_docs/planned/<plan>.md` and `…/<sprint>.json` | **rc=1 (NOT ignored)** for both, with same-call control `git check-ignore -v .ailang/state/x` → **rc=0**, `.gitignore:3:**/.ailang/`. The instrument fires; the two artifacts are commit-able where they are written |

**No gate command in this plan is red at base.** G6 is red at base only in the form the design
doc wrote it (D1), and this plan does not use that form.

---

## 3. Milestones

Order is free except **M3 reads Test B's `lines` variable**, which M1 and M2 both leave in
place. Recommended order M1 → M2 → M3 (largest first). Each milestone is independently
committable and leaves the whole gate green.

### M1 — split the normalizer; validate names strictly (~0.08 d, ~+45/−12)

**File touched:** `host/verifygate/toolchain_pin_gate_test.go` **only**.

**What lands** (the Decision sketch is authoritative for semantics; identifiers may be adjusted,
message text may **not** — D5):

1. `stripPinQuotes(value string) string` — trim + symmetric quote strip, lifted verbatim from
   the first half of today's `normalizeToolchainPin`.
2. `canonicalizeVersionPin(value string) string` — `stripPinQuotes` + the `go`-prefix
   prepend. Serves **only** setup-go `go-version:` and go.mod `go ` directives.
3. `requireToolchainNamePin(t *testing.T, source, key, raw string) string` — `t.Helper()`,
   `stripPinQuotes`, then two `t.Fatalf` arms with the **frozen** text of D5. Arm order:
   selection-mode first (`auto`/`local`/`path`/contains `+`), then `!version.IsValid(value)`.
4. `pinValues` → `keyedValues`, returning **quote-stripped RAW** values (the `normalizeToolchainPin`
   call at today's `:30` is deleted, replaced by `stripPinQuotes`).
5. Call-site disposition:
   - Test A: `goToolchains` maps **per collected value** through
     `requireToolchainNamePin(t, "ci.yml", "GOTOOLCHAIN", raw)`; `goVersions` maps through
     `canonicalizeVersionPin`.
   - `moduleGoFloor` (today `:45`) → `canonicalizeVersionPin`.
   - Test B (today `:244`) → `requireToolchainNamePin(t, scriptPath, "PINNED", pinnedAssignments[0])`,
     inside the existing `if len(pinnedAssignments) == 1` guard.
6. `normalizeToolchainPin` **deleted**.

**Two implementation constraints the doc leaves implicit — both load-bearing:**

- **Validation happens at the Test A call site, over the returned slice — NEVER inside
  `keyedValues`.** `keyedValues(lines, "GOTOOLCHAIN")` is also called for its *length only*
  inside the job-enumeration `t.Errorf` (today `:100`–`:101`). Putting the `Fatalf` inside the
  extractor would make that count call fatal and would convert a job-enumeration mismatch into
  a pin-validity message — the misattribution this item exists to remove, reintroduced one
  layer down.
- **`moduleGoFloor` must NOT get the strict validator.** `design_docs/verification/w-race-gate-blindspot/repro/go.mod:11`
  is `go 1.22`; `version.IsValid("1.22")` is **false**. Routing the go.mod floor through
  `requireToolchainNamePin` reds `TestReproModuleFloorStaysBelowKnownBadToolchains` and Test B
  on a correct tree.

**Import hygiene:** `go/version` stays used (by `requireToolchainNamePin` and by
`TestReproModuleFloorStaysBelowKnownBadToolchains`); `strings` stays used. No import block
change is needed and none should be made.

**Gates for M1** (in this order; stop at the first red):

```bash
export PATH=/opt/homebrew/bin:$PATH
gofmt -l host/verifygate/                                   # must be EMPTY
go vet ./host/verifygate/                                   # rc=0
go build ./...                                              # rc=0
go test ./host/verifygate/ \
  -run 'TestGoToolchainPinsAgreeAndMatchJobList|TestMiscompileInstrumentProbesPinnedToolchain' \
  -count=1 -v > /tmp/m1.out 2>&1; echo rc=$?                # rc=0, 2 RUN / 2 PASS
AILANG_BIN=/tmp/ailang-v0300/ailang go test ./host/verifygate/ -count=1 > /tmp/m1p.out 2>&1; echo rc=$?
grep -c 'AILANG_BIN is unset' /tmp/m1p.out                  # must be 0 — a forgotten export is self-diagnosing
# the name is dead, with a live-grep control:
grep -c 'normalizeToolchainPin' host/verifygate/toolchain_pin_gate_test.go   # 0  (base: 4)
grep -c 'canonicalizeVersionPin' host/verifygate/toolchain_pin_gate_test.go  # >= 3 (base: 0)
grep -c 'requireToolchainNamePin' host/verifygate/toolchain_pin_gate_test.go # >= 3 (base: 0)
# D5 message freeze, with a known-positive control in the same call:
grep -c 'is not an allowed standard Go toolchain pin' host/verifygate/toolchain_pin_gate_test.go  # 1
grep -c 'is a toolchain-selection mode, not a pin' host/verifygate/toolchain_pin_gate_test.go     # 1
grep -c 'the Go runtime itself refuses' host/verifygate/toolchain_pin_gate_test.go                # 0
grep -c 'is not a valid Go toolchain name' host/verifygate/toolchain_pin_gate_test.go             # 0
# venue integrity:
shasum -a 256 .github/workflows/ci.yml design_docs/verification/w-race-gate-blindspot/run.sh
```

The paired 0/≥3 grep counts are the non-vacuity control: the zeros are only meaningful because
the same greps return non-zero for the strings that DID land.

**Mutation arms rehearsed at M1:** M1, M2, M3, M4, M5, M9(a), M9(b) — see §4.

### M2 — Test B's instrument floors go fatal (~0.03 d, ~+2/−2)

**File touched:** `host/verifygate/toolchain_pin_gate_test.go` only.

Two tokens, token-for-token otherwise unchanged: the `t.Errorf` at today's `:199`
(`instrument failure: %s does not contain known-positive control %q`) and today's `:209`
(`instrument failure: %s exact shebang count=%d, want 1`) become `t.Fatalf`.
**Content-anchored, not positional** — after M1 these sit near `:233` and `:243`.

**Gates:** G2, G3, G4, G6 as above, plus:

```bash
grep -c 'instrument failure' host/verifygate/toolchain_pin_gate_test.go     # 11 (unchanged; base 11)
# Test B's two floors are now fatal — count Errorf-class instrument floors:
grep -n 'instrument failure' host/verifygate/toolchain_pin_gate_test.go | grep -c 't.Errorf'  # 0 (base: 2)
grep -n 'instrument failure' host/verifygate/toolchain_pin_gate_test.go | grep -c 't.Fatalf'  # 11 (base: 9)
```

The `11` total with a `2 → 0` / `9 → 11` split is the landed-proof: a count that MOVED, plus a
count that did not.

**Mutation arms rehearsed at M2:** M7, M8 (floor-existence only), **M8′** (the actual fatality
discriminator) — see §4 and D2.

### M3 — pinned-OK guard direction pin (~0.04 d, ~+16/−0)

**File touched:** `host/verifygate/toolchain_pin_gate_test.go` only.

The `guardLines` block from the Decision sketch lands **inside Test B**, after the existing
`saw_pinned_ok` site-count check (today `:263`–`:265`) and before Test B's closing brace. It
reuses Test B's existing `lines` variable. It strips at the first `#` so a comment may quote the
guard literal freely (the iteration-133 namespace lesson), and its `!= 1` predicate is two-sided
(0 = flipped/deleted, ≥2 = duplicated), so it cannot pass vacuously on its own null case (S6).

Shipped message (frozen, as the sketch):

```
%s: executable pinned-OK guard-line count=%d, want exactly 1 — the floor must test ABSENCE of the OK flag (`-eq 0`); a flipped or duplicated guard is a different instrument
```

Note this one is a `t.Errorf`, deliberately — it is not an instrument floor, and Test B's other
non-floor assertions are Errorf.

**Gates:** G2, G3, G4, G6, plus:

```bash
grep -c 'executable pinned-OK guard-line count' host/verifygate/toolchain_pin_gate_test.go  # 1 (base: 0)
grep -c 'saw_pinned_ok' host/verifygate/toolchain_pin_gate_test.go                          # >= 4 (base: 3)
```

**Mutation arms rehearsed at M3:** M6, plus AC6(b)'s must-stay-green comment control — see §4.

---

## 4. Mutation drill — every arm with its measured base, landed-proof, and predicted RED text

**House recipe, every arm, no exceptions:**

```
1. assert the venue's sha256 equals the §0 table
2. apply the mutation
3. print the LANDED-PROOF counts (before/after) — a count that MOVED, or exact line content
4. run the scoped test, capturing rc DIRECTLY (no PIPESTATUS — see §0 rule 3)
5. count messages:  grep -c 'toolchain_pin_gate_test.go:' /tmp/arm.out
6. grep the SHIPPED LITERAL, never a description of it (D5)
7. restore:  cp /tmp/p45-backup/<file> <path>
8. assert sha256 equality with the §0 table   <-- the restore proof
```

**Message counting is by `grep -c 'toolchain_pin_gate_test.go:'` on the `-v` output.** Planner
-verified: it counts exactly the assertion messages and nothing else, on every arm below.

**Scoped test names:**
- Test A = `go test ./host/verifygate/ -run 'TestGoToolchainPinsAgreeAndMatchJobList' -count=1 -v`
- Test B = `go test ./host/verifygate/ -run 'TestMiscompileInstrumentProbesPinnedToolchain' -count=1 -v`

Neither needs `AILANG_BIN`.

| # | Type | Mutation (venue) | Landed-proof (planner-measured) | **BASE result (planner-measured at `f3790e4`)** | **Required POST-milestone result** | Milestone |
|---|---|---|---|---|---|---|
| **M1** | **NEW TEETH — the headline** | `sed -i '' '21s/GOTOOLCHAIN: go1\.26\.6/GOTOOLCHAIN: 1.26.6/'` `ci.yml`; also assert `sed -n '21p'` → `      GOTOOLCHAIN: 1.26.6` | `grep -c 'GOTOOLCHAIN: go1.26.6'` **2→1**; `grep -c 'GOTOOLCHAIN: 1\.26\.6'` **0→1** | **rc=0, `--- PASS`, 0 messages** — the gate certifies a value that kills every `go` command in the job | Test A rc≠0; `grep -c 'toolchain_pin_gate_test.go:'` = **1**; that message contains `is not an allowed standard Go toolchain pin` **and** `"1.26.6"` | M1 |
| **M2** | **NEW TEETH — per-site, not per-line** | same flip at line **102** (the OTHER collected site); assert `sed -n '102p'` | as M1, on the other keyed line | **rc=0, `--- PASS`, 0 messages** (planner-measured; the doc did not measure this arm) | Test A rc≠0; **1** message, same literal, carrying the second site's value | M1 |
| **M3** | **NEW TEETH — pin vs mode** | `sed -i '' '21s/GOTOOLCHAIN: go1\.26\.6/GOTOOLCHAIN: auto/'` `ci.yml` | `grep -c 'GOTOOLCHAIN: auto'` **0→1** | **rc=1, 2 messages, BOTH misattributing**: `:120 ci.yml: toolchain pins disagree: GOTOOLCHAIN=[goauto go1.26.6] go-version=[go1.26.6 go1.26.6]` and `:150 go.mod floor="go1.26.6" disagrees with ci.yml toolchain pin="goauto"` | Test A rc≠0; **1** message containing `is a toolchain-selection mode, not a pin`; `grep -c 'toolchain pins disagree'` = **0**; `grep -c 'disagrees with ci.yml toolchain pin'` = **0** | M1 |
| **M4** | **REGRESSION — red before and after** | `sed -i '' 's/GOTOOLCHAIN:/GO_TOOLCHAIN:/g'` `ci.yml` | `grep -c 'GOTOOLCHAIN:'` **2→0** with same-call control `grep -c 'GO_TOOLCHAIN:'` **0→2** | **rc=1, 1 message**: `:108 ci.yml: GOTOOLCHAIN keyed-line count=0, want 2 (one per enumerated expected job)` | **identical**: rc≠0, 1 message, same literal. The S6 floor must survive the extraction split — this arm proves nothing about M1's new code and must not be read as if it did | M1 |
| **M5** | **NEW TEETH — the static/runtime mispolarity** | `sed -i '' '26s/PINNED="go1\.26\.6"/PINNED="1.26.6"/'` `run.sh` | `grep -c 'PINNED="1\.26\.6"'` **0→1** | **rc=0, `--- PASS`, 0 messages** (planner-measured; the doc had this as code-derived only). Auto-correction makes it statically green while `run.sh:87`'s byte compare `[ "$tc" = "$PINNED" ]` could never set `saw_pinned_ok` | Test B rc≠0; **1** message containing `is not an allowed standard Go toolchain pin` and key `PINNED` | M1 |
| **M6** | **NEW TEETH — direction blindness** | `sed -i '' '140s/-eq 0/-ne 0/'` `run.sh` | `grep -c 'saw_pinned_ok" -eq 0'` **1→0**; `grep -c 'saw_pinned_ok" -ne 0'` **0→1**; `grep -c 'saw_pinned_ok'` **3→3** (the count that must NOT move — that immobility IS the defect) | **rc=0, `--- PASS`, 0 messages** (planner-measured; third independent agreeing run, after the doc's P14 and P23) | Test B rc≠0; output contains `guard-line count=0, want exactly 1` | M3 |
| **M7** | **NEW TEETH — fatality, the real discriminator** | `sed -i '' 's/PINNED="/PINHOLE="/'` `run.sh` | `grep -c 'PINNED='` **1→0** with same-call control `grep -c 'KNOWN_BAD='` = **1** | **rc=1, 3 messages**: `:199` control floor, `:222 PINNED assignment count=0, want 1`, `:247 PINNED="", want go.mod floor "go1.26.6"` | Test B rc≠0; message count **3→1**; the surviving message contains `does not contain known-positive control` and `"PINNED="`; `grep -c 'PINNED assignment count='` = **0**; `grep -c ', want go.mod floor '` = **0**. The base run is the same-session known-positive control proving both needles are observable (**D3**) | M2 |
| **M8** | **FLOOR-EXISTENCE ONLY — NOT a fatality proof (D2)** | `printf '%s\n' '#!/usr/bin/env bash' >> run.sh` | `grep -c '^#!/usr/bin/env bash$'` **1→2** | **rc=1, 1 message**: `:209 instrument failure: <run.sh> exact shebang count=2, want 1`, zero downstream | **identical** — rc≠0, 1 message, same literal. Nothing downstream of the shebang floor depends on shebang count, so this arm cannot distinguish `Errorf` from `Fatalf`. Record it as floor-existence; do NOT report it as evidence M2 landed | M2 |
| **M8′** | **NEW — the shebang floor's actual fatality discriminator (planner-authored, base-measured)** | `printf '%s\n' '#!/usr/bin/env bash' >> run.sh` **and** append a duplicate of `run.sh:25`'s `KNOWN_GOOD="…"` line: `grep '^KNOWN_GOOD="' run.sh \| head -1 >> run.sh` | `grep -c '^#!/usr/bin/env bash$'` **1→2** **and** `grep -c '^KNOWN_GOOD="'` **1→2** | **rc=1, 5 messages**: `:209` shebang floor, `:216 KNOWN_GOOD assignment count=2, want 1`, `:232 KNOWN_GOOD must contain at least one toolchain`, `:240 KNOWN_GOOD=[] does not probe the pinned toolchain go1.26.6 from go.mod`, `:250 PINNED="go1.26.6" is absent from KNOWN_GOOD=[]…` | Test B rc≠0; message count **5→1**; the single surviving message contains `exact shebang count=2, want 1`. A `5→1` delta is what "the floor is fatal" actually looks like | M2 |
| **M9(a)** | **NEW TEETH — single attribution** | M1's flip at site 21 | as M1 | rc=0 (identical to M1's base) | Test A rc≠0; `grep -c 'is not an allowed standard Go toolchain pin'` = **1**; `grep -c 'toolchain pins disagree'` = **0**; `grep -c 'disagrees with ci.yml toolchain pin'` = **0**; total messages = **1** | M1 |
| **M9(b)** | **CONTROL — proves M9(a)'s zeros are measured absences** | restore, then `sed -i '' '21s/GOTOOLCHAIN: go1\.26\.6/GOTOOLCHAIN: go1.25.6/'` `ci.yml` — a VALID name that genuinely disagrees | `grep -c 'GOTOOLCHAIN: go1\.25\.6'` **0→1**; `grep -c 'GOTOOLCHAIN: go1\.26\.6'` **2→1** | **rc=1, 2 messages**, BOTH strings present: `toolchain pins disagree: GOTOOLCHAIN=[go1.25.6 go1.26.6] …` and `go.mod floor="go1.26.6" disagrees with ci.yml toolchain pin="go1.25.6"` | **identical** (the validator passes `go1.25.6`). Must run in the SAME probe session as M9(a); its two non-zero counts are what make M9(a)'s two zeros a measurement rather than a dead grep | M1 |
| **AC6(b)** | **MUST-STAY-GREEN control** | `printf '%s\n' '# guard direction note: [ "$saw_pinned_ok" -eq 0 ] is the floor' >> run.sh` | `grep -c 'saw_pinned_ok'` **3→4** | **rc=0, `--- PASS`, 0 messages** (planner-measured) | **rc=0 still.** The `#`-stripping must let a comment quote the guard literal freely. If this arm reds after M3, M3 shipped the iteration-133 defect: a guard and the prose explaining it competing for one namespace | M3 |

**Reading rule for the whole drill.** M4, M8, M9(b) and AC6(b) are red-before-and-red-after or
green-before-and-green-after **by design**. The arms that must FLIP polarity because of this
sprint are exactly: **M1, M2, M3, M5, M6 (green → red)** and **M7, M8′ (red-with-N-messages →
red-with-1-message)**. If any of those seven does not flip, the milestone did not land — do not
report a pass. If M4, M8, M9(b) or AC6(b) changes behavior, that is a regression, not a win.

---

## 5. Acceptance criteria — as amended, mapped to milestones

| AC | Amended form | Milestone | Base (measured) |
|---|---|---|---|
| AC1 | G4 rc=0 with 2 `=== RUN` + 2 `--- PASS`, paired with G5's nonsense control printing `[no tests to run]` and **0** `=== RUN`. Teeth are §4, not this green. | all | rc=0 (G4), control rc=0/0-RUN (G5) |
| AC2 | Derive the site set at execution time: `grep -n 'GOTOOLCHAIN:' .github/workflows/ci.yml` → **2** lines, equal to `len(wantJobs)`=2 (same-call control: non-zero). For **each** site in turn, run arm M1/M2. **No line number is asserted — the sed targets are derived from the grep.** | M1 | both arms rc=0 today |
| AC3 | Arm M3's exact counts (1 / 0 / 0), with M9(b) as the same-session control. | M1 | rc=1, 2 misattributing messages |
| AC4 | Arm M4 — the S6 vacuous-input floor survives the split. Regression arm (D4). | M1 | rc=1, 1 message |
| AC5 | **Amended per D3**: arm M7's message count **3 → 1**, with the two *observable* downstream needles going 1→0 each. Exec-bit and guard-message assertions dropped (they never fire in this arm). | M2 | rc=1, 3 messages |
| AC5′ | **Added per D2**: arm M8′'s message count **5 → 1**. | M2 | rc=1, 5 messages |
| AC6 | (a) arm M6 flips green→red; (b) AC6(b) control stays green. | M3 | (a) rc=0, (b) rc=0 |
| AC7 | `gofmt -l host/verifygate/` empty; `go vet ./host/verifygate/` rc=0; **`AILANG_BIN=/tmp/ailang-v0300/ailang go test ./host/verifygate/ -count=1` rc=0** (D1); `grep -c 'normalizeToolchainPin' <file>` = 0 with controls `canonicalizeVersionPin` ≥ 3 and `requireToolchainNamePin` ≥ 3; **D5's frozen-message greps (2× count-1, 2× count-0)**; `ci.yml` and `run.sh` sha256 equal to §0 after every rehearsal. | all | all green (§2) |

---

## 6. Risks

| Risk | Mitigation |
|---|---|
| Executor routes `moduleGoFloor` through the strict validator | Named in M1's constraints with the measured counter-example: `repro/go.mod:11` is `go 1.22`, `IsValid("1.22")=false`. Gate G4 catches it immediately (Test B and `TestReproModuleFloorStaysBelowKnownBadToolchains` both red). |
| Executor puts validation inside `keyedValues` | Named in M1's constraints. Symptom: a job-enumeration mismatch reports a pin-validity message. Not caught by any base gate — this is a review item. |
| Executor writes the `Fatalf` text from the doc's thesis prose | D5 freezes both strings; AC7 greps for the two banned phrases with a live-grep control. |
| Executor sed's a line number in the deliverable file | D6: all deliverable-file edits are content-anchored; the only positional seds act on `ci.yml`/`run.sh`, which never change, and each still asserts line content. |
| `git checkout --` used to restore | §0 rule 5. It would delete uncommitted milestone work. |
| Whole-package gate run without `AILANG_BIN` | D1; every invocation in this plan carries it, and M1's gate block greps `AILANG_BIN is unset` and asserts 0. |
| M8 reported as fatality evidence | D2; the table types it `FLOOR-EXISTENCE ONLY` and M8′ carries the real 5→1 delta. |

## 7. Non-goals

- `ci.yml`, `run.sh`, `go.mod`, `scripts/verify_go.sh`, `scripts/verify_ail.sh`,
  `host/store/toolchain_canary_test.go`, `tools/launchd/*` — all **byte-untouched**.
- No `.ail` file is touched; `./scripts/build_world_package.sh` must not run. G8's floor
  (11 identities / 40 named tests) must not move.
- Structural YAML / actionlint adoption, availability probing, `verify_go.sh` GOTOOLCHAIN
  strictness, and generalizing the direction pin to run.sh's other floors are all Deferred
  Scope in the design doc and stay out.
- Do **not** widen any acceptance gate to `go test ./...`. The narrowest gate that can fail for
  this diff is `go test ./host/verifygate/ -count=1`.

---

**SPRINT_PLAN_PATH**: `design_docs/planned/w-toolchain-pin-normalizer-accepts-malformed-gotoolchain-sprint-plan.md`
**SPRINT_JSON_PATH**: `design_docs/planned/sprint_w-toolchain-pin-normalizer-accepts-malformed-gotoolchain.json`
