# Sprint plan — `P43` (queue row 43, `w-floor-raise-coupling-inventory`)

**Design doc**: [`w-floor-raise-coupling-inventory.md`](w-floor-raise-coupling-inventory.md)
(621 lines, reviewed; round 1 BLOCKED → revision 1; round 2 BLOCKED → narrow-refinement
carve-out applied by the controller). **THE DESIGN DOC WINS** over this plan wherever the two
differ, except at the three places named in §1, each of which carries the planner's own
measurement of why the doc's literal text would red a correct landing.

**Planner**: `claude-opus-5`, mission-control iteration 132.
**Worktree**: `/Users/voightkampff/dev/sunholo-data/.wt-world-iter132-design`, base `50c3b91`.
**Estimate**: ~0.1 day. Four files touched (1 created, 3 modified), ~95 changed lines, **zero
executable gate lines, zero thresholds moved, zero existing assertions edited**.
**Risk**: LOW on the typing, MEDIUM on the *reading*: two of this sprint's gates are green at
base while measuring nothing (AC1's naive `-run` form, AC5's diff-by-identity). Budget the
time in the acceptance sweep and the mutation drill, not in the edits.

---

## 0. Rules of engagement for the executor (read first — these override habit)

1. **NO git write operations.** No `add`, `commit`, `stash`, `checkout`, `branch`, `worktree`,
   `mv`, `restore`, `reset`. Git **reads** are required and allowed: `git show`, `git status`,
   `git diff`. The controller commits; the sprint hands off an **uncommitted** tree.
2. **`git status --porcelain` is NEVER `0` during this sprint** and that is correct by
   construction. Do not treat a non-zero porcelain as a defect and do not "clean" it.
3. **Restores during the mutation drill use `cp` from `/tmp/p43-backup/`, NEVER
   `git checkout --`.** The entire deliverable is uncommitted; `git checkout -- <path>` would
   delete this sprint's work on `scripts/verify_ail.sh` and `design_docs/coding-standards.md`
   and there would be no way back. §5.1 is the only restore recipe.
4. **Snapshots.** After each milestone that changes files (M1, M2, M3), copy the **full current
   content** of every deliverable path that exists into `.snap/M<k>/`, mirroring the repo-relative
   path. Snapshots are **cumulative** (M3's snapshot contains M1's and M2's files too, at their
   M3 state) so the controller can reconstruct commits in order. Snapshot **only** the four
   deliverable paths — nothing else. `.snap/` is safe: Go's `./...` ignores directories whose
   name begins with `.`, and every repo scanner in play globs a fixed non-recursive directory
   (`scripts/*.sh`, `host/verifygate/*.go`), so a snapshot copy is invisible to all of them.
   (Planner-verified, §1 D2 note.)
5. **Never read `rc` alone on a `-run`-scoped Go test.** A `-run` selector that matches nothing
   exits **0** with `[no tests to run]`. This is not hypothetical here — it is AC1's measured
   base state (§0.2). Every scoped test in this plan is judged by **counting `=== RUN` and
   `--- PASS` / `--- FAIL` lines from `-v` output**.
6. **Environment.** `export AILANG_BIN=/tmp/ailang-v0300/ailang
   WORLD_PKG_AILANG_BIN=/tmp/ailang-v0300/ailang` for the pinned gate (AC5, M6-drill). The
   inventory test itself takes **no** `AILANG_BIN` and is run with `env -u AILANG_BIN` exactly
   as AC1 specifies. Every `ailang` invocation prints `Observatory: NNNMB (warn threshold:
   200MB)` on stderr; the gate is rc=0 with it present. **Not a regression, do not "fix" it.**
7. Do not touch `packages/world-core/**`, `world/*.ail`, the golden, the marker at
   `module_manifest_gate_test.go:128`, `verify_world_package.sh`, `build_world_package.sh`, or
   `ci.yml`. Do not run `./scripts/build_world_package.sh` — this sprint changes no `.ail` byte.

---

## 0.1 Base readings taken as GIVEN from the controller, and re-derived by the planner

All controller readings were **reproduced first-party by the planner in this worktree at
`50c3b91`** and agreed exactly. `VERIFIED BY CONTROLLER` + `RE-DERIVED BY PLANNER` on all seven.

| AC | Command (verbatim) | Base observation |
|---|---|---|
| AC1 | `env -u AILANG_BIN go test ./host/verifygate/ -run '^TestFloorRaiseInventoryNamesEveryCoupledFile$' -count=1` | **rc=0**, `ok … 0.276s [no tests to run]` — **green while measuring nothing** |
| AC2 | `grep -c 'FLOOR-RAISE COUPLING INVENTORY' scripts/verify_ail.sh` / KP `grep -c 'Leg 1' scripts/verify_ail.sh` | `0` (rc=1) / KP `6` firing |
| AC3 | `grep -c '^## S8' design_docs/coding-standards.md` / KP `grep -c '^## S6' …` | `0` (rc=1) / KP `1` firing |
| AC4 | `grep -c 'interfaceHash' scripts/verify_ail.sh` and `… design_docs/coding-standards.md` / KP `grep -c 'EXACT_TOTAL_VERIFIED=11' scripts/verify_ail.sh` | `0` and `0` / KP `1` firing |
| AC5 | `diff <(grep -v '^[[:space:]]*#' scripts/verify_ail.sh) <(git show 476069d:scripts/verify_ail.sh \| grep -v '^[[:space:]]*#') \| wc -l` | `0` — **green by identity at base** |
| AC6 | `grep -c 'at :224' docs/SELF_MOD_PUBLISH.md` / `grep -c 'its final leg invokes verify_world_package.sh' …` | `1` (red at base, as required) / `0` |
| AC7 | `go vet ./host/verifygate/` ; `gofmt -l host/verifygate/ \| wc -l` | rc=0 ; `0` |

## 0.2 What the planner measured on top, and why each is load-bearing

- **The gate is 6 s, not minutes.** `AILANG_BIN=… WORLD_PKG_AILANG_BIN=… ./scripts/verify_ail.sh`
  → **rc=0 in 6.0 s wall**, banner `✓ 11/11 required world/ identities verified across 11
  module(s)` / `✓ all 40 required named tests pass (failed_tests=0)` / `✓ verify gate PASSED`,
  `git status --porcelain | wc -l` → `0` after. AC5's gate clause and the M6 drill arm are
  cheap; run them, do not skip them.
- **`go build ./...` → rc=0 in 1.5 s. `AILANG_BIN=… go test ./host/runbook/ -count=1` → `ok …
  1.364s`.** AC6's third clause is cheap.
- **AC5's diff base `476069d` is valid from `50c3b91`**: `git diff --stat 476069d HEAD --
  scripts/verify_ail.sh` → empty (the only delta between the two commits is the design doc
  itself, +621 lines). The AC5 command is runnable **verbatim** from this worktree.
- **All nine needles are present in both proposed homes, as the doc writes them.** Extracting
  §(a)'s fenced block and §(b)'s fenced §S8 from the design doc and grepping each needle:
  `packages/world-core/world/` 1/1, `REQUIRED_VERIFIED` 2/1, `EXACT_TOTAL_VERIFIED` 2/2,
  `world_package_ready_packet.golden.json` 1/1, `SELF_MOD_PUBLISH.md` 1/1,
  `module_manifest_gate_test.go` 1/1, `interfaceHash` 1/1, `does not move for` 1/1, plus
  site 1's row literal 1/1 in each. **The doc's proposed text satisfies its own test.**
- **The vacuity witness is real in both homes**: bare `world/<module>.ail` occurs **2×** in the
  script block and **2×** in §S8 (site 1's row and site 2's row), so deleting site 1 leaves the
  bare needle firing at 1. That is the M7/M8 drill's whole point and it is confirmed here.
- **Neither home contains the token `world-publish`** (0/0), so AC30's fatal scanner
  (`runbook_stageb_test.go:341-381`) cannot fire on the landing. AC30's known-positive control
  counts `verify_world_package.sh` in `verify_ail.sh` with a `< 1` bound; base count is `1` and
  the new block adds none of that token, so the control still fires.
- **The block does not collide with any count-exactly-1 mutation anchor** used by the four
  source-reading tests over `scripts/verify_ail.sh`: `v0.30.0|OK`, `mods+=("${f#./}")`,
  `(-dirty$|-[0-9]+-g[0-9a-f]+)`, `LEG1_MODULES=(` — **0 occurrences of each inside the block**.
- **`design_docs/coding-standards.md` is bound by no Go test** (`grep -rn 'coding-standards'
  --include='*.go' --include='*.sh' --include='*.yml' .` → one unrelated comment). §S8's only
  enforcement will be the new inventory test.
- **`docs/SELF_MOD_PUBLISH.md`'s Stage-A control survives the `:39` edit**:
  `runbook_commands_test.go:147` requires the literal `./scripts/verify_world_package.sh` in
  Stage A; it stands at `:25` and `## Stage B` is at `:50`. The `:39` edit is inside Stage A but
  removes no such literal.
- **`repoRoot` is available**: it is a package-level var at
  `host/verifygate/ail_binary_gate_test.go:27` in `package verifygate`; the new file joins that
  package and reuses it, exactly as §(c) requires.

---

## 1. Design-doc defects the planner found — three, none fatal, two would have shipped RED

### D1 — **AC2's "exactly the 25-line block" is STALE; the doc's own §(a) block is 28 lines** · severity HIGH

The design doc's AC2 reads: *"emits exactly the 25-line block (first line the begin marker,
last the END marker, V28)"*. Measured, first-party, by extracting §(a)'s fenced block from the
design doc itself:

```
awk '/^# ── FLOOR-RAISE COUPLING INVENTORY/,/^# ── END FLOOR-RAISE COUPLING INVENTORY/' \
  design_docs/planned/w-floor-raise-coupling-inventory.md | wc -l     →  28
```

`25` was correct when V28 measured it — against the **pre-revision** block. Revision 1 and the
round-2 carve-out then grew the block by three lines (the Tier-2 forward-reference sentence and
the softened `interfaceHash` wording). **The number in AC2 was transcribed, not re-derived, and
that is the exact defect this whole item exists to name.** An executor asserting `= 25` against
a correct landing reds AC2.

**Resolution — the doc's intent survives, its number does not.** AC2 in §4 asserts the
properties V28 actually established (the range is *bounded*: first line is the begin marker,
last line is the END marker, and the count is far below the unanchored form's runaway) and
**re-derives** the expected count from §(a) in the same call rather than hardcoding it. The
executor never types a line count.

### D2 — **`TestNoRigAbsolutePaths` scans every `*.go` in `host/verifygate/`, and the doc's §Conflict Surface does not name it** · severity MEDIUM

`host/verifygate/ail_binary_gate_test.go:552` globs `host/verifygate/*.go`, reads each file, and
`t.Errorf`s if the source contains `/tmp/ailang`, `/Users/`, or `/home/runner/` (needles
assembled at runtime so the scanner does not match itself). The new file
`floor_raise_inventory_test.go` **will be scanned**. The design doc's §Conflict Surface names
the five module-manifest tests, AC30, AC28/AC29 and `coding-standards.md` — but not this one.

**Concrete failure this avoids**: §(c)'s prose cites the design-time probes and the pinned
binary. An executor writing a provenance comment such as `// prototyped at /tmp/ailang-v0300`
into the new test lands a **red package** for a reason unrelated to the design.

**Resolution**: M1 forbids any absolute path in the new file — in code *or* comments. Cite
provenance as `V20/V26/V27` and the doc path, never a rig path. §4 AC7 gains the whole-package
run that catches it. (This same class is why `.snap/` uses a leading dot, per §0 rule 4.)

### D3 — **the pointer comment lands INSIDE a Python heredoc; §(a)'s phrasing hides that** · severity LOW

§(a) says the pointer line goes *"directly above the `REQUIRED_VERIFIED = {` heredoc at `:274`"*.
Measured: the heredoc opens at `:270` (`python3 - "$mod" "$tmp_json" <<'PY'`) and
`REQUIRED_VERIFIED = {` is at `:274`, i.e. **inside** it. The pointer is therefore Python source,
not shell source.

**Resolution** (no change to the doc's intent — the line text is unchanged): place it at
**column 0**, matching the existing Python comment at `:273`, so Python tokenizes it as a
comment. AC5's `grep -v '^[[:space:]]*#'` filter strips it from the non-comment diff either way,
so AC5 stays empty. Verified downstream by the AC5 clause and by the full pinned gate re-run —
if the heredoc were broken, Leg 1 would die loudly.

### Not a defect, recorded so nobody re-derives it

- **AC5 is green at base by identity** (`0` lines of diff). That is by design; its teeth are the
  M6 drill arm, which must make it **non-empty**. Do not report AC5's base green as a finding.
- **AC1's base rc=0 is the sprint's sharpest fact**, not an anomaly. See §0 rule 5.
- **`.snap/` is invisible to `go build ./...` / `go vet ./...`** because Go ignores `.`-prefixed
  directories in `./...` patterns. This is why rule 4 mandates that exact directory name.

---

## 2. Milestone → acceptance-criterion mapping

**The design doc defines no milestones** — it defines four artifacts (§Files to Create/Modify
(a)–(d)) and seven acceptance criteria. The mapping below is planner-derived by reading each
AC's own subject back to the artifact that produces it. **Diffed against the doc's AC list:
every AC1–AC7 is claimed by exactly one milestone, and no milestone claims an AC the doc
attaches to a different artifact. No disagreement with the doc found.**

| Milestone | Artifact (doc §) | ACs it closes | ACs it re-confirms |
|---|---|---|---|
| **M1** — the enforcement test, red-before-green | §(c) `host/verifygate/floor_raise_inventory_test.go` | **AC7** | — (AC1's `=== RUN` clause becomes satisfiable here; AC1 is *closed* by M2) |
| **M2** — the two durable homes | §(a) `scripts/verify_ail.sh` block + pointer; §(b) `coding-standards.md` §S8 | **AC1, AC2, AC3, AC4, AC5** | AC7 |
| **M3** — the stale-anchor repair | §(d) `docs/SELF_MOD_PUBLISH.md:39` | **AC6** | — |
| **M4** — full acceptance sweep | (no edits) | — | **AC1–AC7**, all seven in one recorded pass |
| **M5** — mutation drill + hand-off | (no net edits; every arm restored) | — | drill arms M1–M8 + green control |

**Why M1 precedes M2 (red-before-green).** The test is written first and observed **RED** —
loudly, through its own instrument fatal (`begin marker count=0, want 1`), because at base the
script carries no markers at all. That red is the proof the test is wired to the artifact. V20
arm (a) measured exactly this shape at design time against the real `scripts/verify_ail.sh`.
Landing the homes first would make the test green on its first-ever run, which proves nothing.

---

## 3. The exact file touch set

| # | Path | Action | Approx. changed lines | Milestone |
|---|---|---|---|---|
| 1 | `host/verifygate/floor_raise_inventory_test.go` | **CREATE** | ~60 | M1 |
| 2 | `scripts/verify_ail.sh` | MODIFY — 28-line comment block after `:29`, + 1 pointer comment at column 0 above `:274` | +29, −0 | M2 |
| 3 | `design_docs/coding-standards.md` | MODIFY — §S8 appended after §S7, before the `---` ratification footer | +29, −0 | M2 |
| 4 | `docs/SELF_MOD_PUBLISH.md` | MODIFY — line `:39` only | +1, −1 | M3 |

Nothing else. Plus the untracked, uncommitted `.snap/M1/`, `.snap/M2/`, `.snap/M3/`.

---

## 4. Milestones, ordered

### M1 — `host/verifygate/floor_raise_inventory_test.go`, and its RED

**Write** one file, `package verifygate`, one exported test
`TestFloorRaiseInventoryNamesEveryCoupledFile(t *testing.T)`, implementing §(c) verbatim:

1. Read `filepath.Join(repoRoot, "scripts", "verify_ail.sh")`.
2. Locate the block by the **styled** literals `# ── FLOOR-RAISE COUPLING INVENTORY` and
   `# ── END FLOOR-RAISE COUPLING INVENTORY`. **`t.Fatalf` unless each occurs EXACTLY ONCE and
   begin precedes end** — this is the known-positive control: an empty, vanished, duplicated or
   misordered enumeration must fail **loudly**, never pass as zero. Never anchor on the bare
   phrase `FLOOR-RAISE COUPLING INVENTORY`: the pointer line repeats it by design (V28).
3. Assert the bounded block contains all **nine** script-home needles:
   `packages/world-core/world/`, `REQUIRED_VERIFIED`, `EXACT_TOTAL_VERIFIED`,
   `world_package_ready_packet.golden.json`, `SELF_MOD_PUBLISH.md`,
   `module_manifest_gate_test.go`, `interfaceHash`, `does not move for`, and the
   enumerated-row literal `#   1. world/<module>.ail`.
   Failure message shape: `inventory block omits %q`.
4. Read `filepath.Join(repoRoot, "design_docs", "coding-standards.md")`. `t.Fatalf` unless the
   `## S8` heading exists. Extract **only** the text from that heading to the next `##` heading
   or EOF, and assert within **that bounded extract** the eight shared needles plus §S8's site-1
   row literal `` | 1 | `world/<module>.ail` ``. Failure message shape: `S8 omits %q`.
   (The bound is round-2's non-blocking catch: a future §S9 reusing a needle term must not
   satisfy §S8's assertion.)

**Forbidden in this file** (D2): any absolute path — `/tmp/…`, `/Users/…`, `/home/runner/…` —
in code or comments. No `requirePinned`, no `AILANG_BIN`, no subprocess, no `exec.Command`. It
must run in any lane. Reuse `repoRoot` only.

**Gates for M1** — run and record all four:

```bash
gofmt -l host/verifygate/ | wc -l                      # expect 0        (base: 0)
go vet ./host/verifygate/                              # expect rc=0     (base: rc=0)
go build ./...                                         # expect rc=0     (base: rc=0, 1.5s)
env -u AILANG_BIN go test ./host/verifygate/ \
  -run '^TestFloorRaiseInventoryNamesEveryCoupledFile$' -count=1 -v 2>&1 | tee /tmp/p43-m1.txt
grep -c '^=== RUN' /tmp/p43-m1.txt                     # expect 1        (base: 0)
grep -c '^--- FAIL' /tmp/p43-m1.txt                    # expect 1  ← THE RED
grep -c 'begin marker count=0' /tmp/p43-m1.txt         # expect 1  ← through the INSTRUMENT fatal
```

**M1 is DONE when the test RUNS (1 `=== RUN`) and FAILS through its marker fatal.** A `--- PASS`
here means the extractor found markers that cannot exist yet — stop and diagnose. rc will be 1;
that rc is *expected*, and it is still not what you are reading. Read the counted lines.

**Closes: AC7.** Snapshot → `.snap/M1/host/verifygate/floor_raise_inventory_test.go`.

### M2 — the two durable homes

**(a) `scripts/verify_ail.sh`** — insert the design doc's §(a) fenced block **verbatim, byte for
byte** after line 29 (after the `NOTE (test-name collisions, V22)` header comment, before
`set -uo pipefail`). Do not reflow, do not retype: extract it from the design doc so the
box-drawing characters and column alignment survive.

```bash
awk '/^# ── FLOOR-RAISE COUPLING INVENTORY/,/^# ── END FLOOR-RAISE COUPLING INVENTORY/' \
  design_docs/planned/w-floor-raise-coupling-inventory.md > /tmp/p43-block.txt
wc -l < /tmp/p43-block.txt        # 28 today — re-derive it, never transcribe it (D1)
```

Then insert the pointer line at **column 0**, directly above `REQUIRED_VERIFIED = {` (currently
`:274`, **inside** the `<<'PY'` heredoc that opens at `:270` — D3), matching the existing
column-0 Python comment at `:273`:

```
# Before adding an identity here: read the FLOOR-RAISE COUPLING INVENTORY at the head of this file.
```

**(b) `design_docs/coding-standards.md`** — append the design doc's §(b) fenced §S8 **verbatim**
after §S7's last paragraph and **before** the `---` ratification footer. Renumber nothing.

> `coding-standards.md` is ratification-class. The authorization is the ratified queue row 43
> itself, which names this file as a durable home (`world-mission.md:4403`). The §S8 text is
> proposed verbatim in the design doc precisely so the PR ratifies exactly what lands — so land
> exactly it, with no editorial improvements.

**Gates for M2**, each with its base observation:

```bash
# AC2 — bounded, and the count is re-derived, not typed (D1)
awk '/^# ── FLOOR-RAISE COUPLING INVENTORY/,/^# ── END FLOOR-RAISE COUPLING INVENTORY/' \
  scripts/verify_ail.sh > /tmp/p43-landed.txt
head -1 /tmp/p43-landed.txt | grep -c '^# ── FLOOR-RAISE COUPLING INVENTORY'      # 1  (base: 0)
tail -1 /tmp/p43-landed.txt | grep -c '^# ── END FLOOR-RAISE COUPLING INVENTORY'  # 1  (base: 0)
diff /tmp/p43-block.txt /tmp/p43-landed.txt && echo BLOCK_IS_VERBATIM             # empty
grep -c 'FLOOR-RAISE COUPLING INVENTORY' scripts/verify_ail.sh                    # 3  (base: 0, rc=1)
grep -c 'Leg 1' scripts/verify_ail.sh                                             # 6  KP, firing
```
The `3` is begin marker + END marker + pointer line — **3 precisely because the pointer exists**.
The awk range MUST keep the `# ── ` styling: unanchored on the bare phrase it re-opens at the
pointer and runs to EOF (V28 measured 32 lines ending in the script tail vs the bounded block).

```bash
# AC3   base: 0 (rc=1); KP `^## S6` → 1
grep -c '^## S8' design_docs/coding-standards.md            # 1
grep -c '^## S6' design_docs/coding-standards.md            # 1   KP, firing
# AC4   base: 0 in BOTH files; KP EXACT_TOTAL_VERIFIED=11 → 1
grep -c 'interfaceHash' scripts/verify_ail.sh               # >= 1
grep -c 'interfaceHash' design_docs/coding-standards.md     # >= 1
grep -c 'EXACT_TOTAL_VERIFIED=11' scripts/verify_ail.sh     # 1   KP, firing
# AC5 clause 1 — gate strength untouched.   base: 0 (empty, by identity)
diff <(grep -v '^[[:space:]]*#' scripts/verify_ail.sh) \
     <(git show 476069d:scripts/verify_ail.sh | grep -v '^[[:space:]]*#') | wc -l   # 0
# AC5 clause 2 — the gate still passes with the identical banner.   base: rc=0, 6.0s
AILANG_BIN=/tmp/ailang-v0300/ailang WORLD_PKG_AILANG_BIN=/tmp/ailang-v0300/ailang \
  ./scripts/verify_ail.sh > /tmp/p43-gate.txt 2>&1; echo "rc=$?"        # rc=0
grep -c '✓ 11/11 required world/ identities verified across 11 module(s)' /tmp/p43-gate.txt  # 1
grep -c '✓ all 40 required named tests pass' /tmp/p43-gate.txt                               # 1
# AC1 — the test now RUNS and PASSES.   base: rc=0 with `[no tests to run]`, 0 `=== RUN`
env -u AILANG_BIN go test ./host/verifygate/ \
  -run '^TestFloorRaiseInventoryNamesEveryCoupledFile$' -count=1 -v 2>&1 | tee /tmp/p43-m2.txt
grep -c '^=== RUN' /tmp/p43-m2.txt      # 1
grep -c '^--- PASS' /tmp/p43-m2.txt     # 1
grep -c '^--- FAIL' /tmp/p43-m2.txt     # 0
# AC1's paired nonsense control, same call shape — proves the counter can read a zero
env -u AILANG_BIN go test ./host/verifygate/ -run 'TestNoSuchInventoryZZZ' -count=1 -v 2>&1 \
  | grep -c 'no tests to run'           # 1   ← and its rc is 0. THIS is why rc is never the gate.
```

**Closes: AC1, AC2, AC3, AC4, AC5.** Snapshot → `.snap/M2/` (test file + both edited files).

### M3 — `docs/SELF_MOD_PUBLISH.md:39`, the stale-anchor repair

Replace, on line 39 only:

```
./scripts/verify_ail.sh   # invokes verify_world_package.sh at :224
```
with
```
./scripts/verify_ail.sh   # its final leg invokes verify_world_package.sh
```

The number was two refactors stale (the call site is `:403`). **Remove the rotting number, do
not refresh it** — refreshing mints the next stale anchor (§Alternatives rejected 4).

```bash
# AC6.  base: 1 / 0 / rc=0 `ok … 1.364s`
grep -c 'at :224' docs/SELF_MOD_PUBLISH.md                                     # 0   (base: 1)
grep -c 'its final leg invokes verify_world_package.sh' docs/SELF_MOD_PUBLISH.md  # 1 (base: 0)
grep -c '^\./scripts/verify_world_package\.sh$' docs/SELF_MOD_PUBLISH.md       # 1   KP: the
#   Stage-A control literal runbook_commands_test.go:147 requires, at :25, untouched by this edit
AILANG_BIN=/tmp/ailang-v0300/ailang go test ./host/runbook/ -count=1; echo "rc=$?"  # rc=0
```

**Closes: AC6.** Snapshot → `.snap/M3/` (all four files).

### M4 — the full acceptance sweep, recorded in one pass

Re-run **AC1 through AC7 in order**, from the commands above plus:

```bash
go vet ./host/verifygate/ ; echo "vet rc=$?"                    # rc=0
gofmt -l host/verifygate/ | wc -l                               # 0
go build ./... ; echo "build rc=$?"                             # rc=0
# AC7, widened per D2: the whole package, so TestNoRigAbsolutePaths actually scans the new file
AILANG_BIN=/tmp/ailang-v0300/ailang go test ./host/verifygate/ -count=1 -v 2>&1 \
  | tee /tmp/p43-pkg.txt
grep -c '^--- FAIL' /tmp/p43-pkg.txt                            # 0
grep -c 'rig-absolute path' /tmp/p43-pkg.txt                    # 0   ← D2's guard
grep -c '^--- PASS' /tmp/p43-pkg.txt                            # >= 30 (29 pre-existing + the new one)
```

Record every observed value next to its expectation. **This is the "green control for all arms"
the design doc's §Non-Vacuity requires** — capture it before touching anything in M5.

### M5 — the mutation drill (§5) and hand-off

Run §5 in full, then leave the tree **uncommitted** with `.snap/M1..M3/` present and report:
the four touched paths, the M4 sweep table, the eight drill arms with their observed RED lines,
and the post-drill green control.

---

## 5. The mutation drill — 8 named RED arms + 1 green control

Every arm mutates the **production side** (the two published documents, or the gate constant).
**No arm edits the test.** Each arm must show three things, in this order:

1. **The mutation LANDED** — proved by a *query against the system's own view* (an occurrence
   count moving, a heading vanishing), **never** by the file's bytes or a diff of them.
2. **The tree still BUILDS** — `go build ./... && bash -n scripts/verify_ail.sh` both rc=0.
   A mutation that breaks the build reds every test for the wrong reason, and the arm proves
   nothing. (Recall iteration 131's lesson: *a mutation that did not apply is a false green.*)
3. **The expected single test REDS**, judged by counted `--- FAIL` lines and the named message —
   never by rc.

### 5.1 The house restore recipe — read before arm 1

**Take the backup AFTER M3, from the POST-SPRINT tree.** Backing up from `476069d` or from
base would restore the sprint away.

```bash
mkdir -p /tmp/p43-backup/scripts /tmp/p43-backup/design_docs
cp scripts/verify_ail.sh                 /tmp/p43-backup/scripts/verify_ail.sh
cp design_docs/coding-standards.md       /tmp/p43-backup/design_docs/coding-standards.md
shasum -a 256 /tmp/p43-backup/scripts/verify_ail.sh \
              /tmp/p43-backup/design_docs/coding-standards.md > /tmp/p43-backup/POST_SPRINT.sha
cat /tmp/p43-backup/POST_SPRINT.sha      # record these two hashes in the report
```

**Restore after EVERY arm — `cp`, never `git checkout --`:**

```bash
cp /tmp/p43-backup/scripts/verify_ail.sh           scripts/verify_ail.sh
cp /tmp/p43-backup/design_docs/coding-standards.md design_docs/coding-standards.md
shasum -a 256 scripts/verify_ail.sh design_docs/coding-standards.md   # compare to POST_SPRINT.sha
```

`git checkout -- scripts/verify_ail.sh` would **delete this sprint's block**, because nothing is
committed. There is no undo. Do not type it.

`git status --porcelain` is non-zero throughout and that is correct (§0 rule 2). The integrity
check is the sha256 comparison above, plus `git diff --stat` showing **only** the four expected
paths.

### 5.2 The arms

Shorthand for the scoped run, used by every arm:

```bash
T() { env -u AILANG_BIN go test ./host/verifygate/ \
        -run '^TestFloorRaiseInventoryNamesEveryCoupledFile$' -count=1 -v 2>&1; }
```

| # | Exact edit | Land-proof (system's own view) | Build proof | Expected RED — the ONE test |
|---|---|---|---|---|
| **M1** | delete the `#   5. docs/SELF_MOD_PUBLISH.md …` line from the script block | `grep -c 'SELF_MOD_PUBLISH.md' scripts/verify_ail.sh` **1 → 0** | `go build ./...` rc=0; `bash -n scripts/verify_ail.sh` rc=0 | `T \| grep -c '^--- FAIL'` = 1 and output contains `inventory block omits "docs/SELF_MOD_PUBLISH.md"` |
| **M2** | delete **only** the END marker line | `grep -c '^# ── END FLOOR-RAISE COUPLING INVENTORY' …` **1 → 0**; `grep -c 'FLOOR-RAISE COUPLING INVENTORY' …` **3 → 2** | same | `T` fatals through the **instrument**: `END marker count=0, want 1`. It must **not** pass as a silent zero (V20 arm d, V27 arm g) |
| **M3** | delete §S8 from `coding-standards.md` | `grep -c '^## S8' design_docs/coding-standards.md` **1 → 0** | `go build ./...` rc=0 | `T` fatals on the `## S8` heading assertion |
| **M4** | in the script block, `does not move for` → `does move for` | `grep -c 'does not move for' scripts/verify_ail.sh` **1 → 0** AND `grep -c 'does move for' …` **0 → 1** | same | `inventory block omits "does not move for"` |
| **M5** | delete site 6's row from **§S8 only**, script block intact | `grep -c 'module_manifest_gate_test.go' design_docs/coding-standards.md` **1 → 0**, while `grep -c 'module_manifest_gate_test.go' scripts/verify_ail.sh` stays **1** | same | `S8 omits "module_manifest_gate_test.go"` — **the cross-home clause: the two hand-authored homes may not drift apart** |
| **M6** | `EXACT_TOTAL_VERIFIED=11` → `EXACT_TOTAL_VERIFIED=12` (**AC5's teeth**) | `grep -c 'EXACT_TOTAL_VERIFIED=12' scripts/verify_ail.sh` **0 → 1** | `bash -n scripts/verify_ail.sh` rc=0 | **two reds, not the inventory test**: (i) AC5's non-comment diff → **non-zero** line count; (ii) `AILANG_BIN=… WORLD_PKG_AILANG_BIN=… ./scripts/verify_ail.sh` → **rc=1** with `✗ expected exactly 12 proven world/ contracts, got 11` (the long-established `:324` refusal). ~6 s |
| **M7** | delete **only** `#   1. world/<module>.ail` from the script block; site 2's row intact | `grep -c '^#   1\. world/<module>\.ail' scripts/verify_ail.sh` **1 → 0**; **and the vacuity witness**: `grep -c 'world/<module>\.ail' scripts/verify_ail.sh` **2 → 1** — the bare needle still fires | `bash -n` rc=0 | `inventory block omits "#   1. world/<module>.ail"`. **Record the witness**: a naive bare-path needle would be GREEN under this exact arm. That green is the measured vacuity the reviewer's literal fix would have shipped (V26, V27 arms b/c) |
| **M8** | delete **only** site 1's table row from §S8; site 2's row intact | ``grep -c '^| 1 | `world/<module>\.ail`' design_docs/coding-standards.md`` **1 → 0**; witness: `grep -c 'world/<module>\.ail' design_docs/coding-standards.md` **2 → 1** | `go build ./...` rc=0 | ``S8 omits "| 1 | `world/<module>.ail`"``; same naive-needle witness in the §S8 home (V27 arms e/f) |

**Green control**: after the last restore, re-run the whole of §M4. All seven ACs green, the
package run clean, and the two sha256s equal to `POST_SPRINT.sha`. **Run the control before arm
1 as well** (M4 already provides it) — a drill whose control is only run at the end cannot
distinguish "the arms restored correctly" from "the arms never applied".

### 5.3 Not mutated, and declared

- **The test file itself.** No arm edits it; every red is produced from the production side.
- **`docs/SELF_MOD_PUBLISH.md`.** Per §Open Decisions OD-2 the design doc deliberately declines
  to pin the runbook repair with a needle: a prose needle on a doc line no behaviour depends on
  is ceremony. AC6's greps are the check; AC28/AC30 guard what matters in that file. **Declared
  residual, doc-sanctioned — do not invent an arm for it.**
- **The marker at `module_manifest_gate_test.go:128`, the golden, `packages/world-core/**`.**
  Out of scope; a mutation there tests row 42's work, not this one.

---

## 6. What this sprint explicitly does NOT do

- **No derivation.** The verifygate control's expected marker is **not** computed from
  `EXACT_TOTAL_VERIFIED`, `REQUIRED_VERIFIED`'s cardinality, or any value the gate produces. A
  control whose expectation is computed from the value it checks cannot fail when that value is
  wrong. If the executor's instinct is "why not just derive it" — that instinct was the
  controller's at `P6.V`, the evaluator refuted it, and the refutation is **binding here**
  (§Non-Goals 1). This sprint adds **zero** derivation anywhere; the test compares, never generates.
- **No executable gate changes.** No threshold, no assertion, no strength change. AC5 pins it.
- **No auto-generated inventory.** A generator's sweep scope silently becomes the map — the same
  defect one level up.
- **No Tier-2 publication.** The §Tier 2 list stays in the design doc, unpublished, pending the
  rehearsal-gated follow-up row `w-floor-raise-tier2-inventory` (the **controller** files it,
  not the executor). Both homes carry only the one-line forward reference.
- **No repair of `w-mcp-projection.md`'s historical table**, no `run.sh` work (row 44), no
  token-counting-control work (row 49).
- **No `./scripts/build_world_package.sh` run.** No `.ail` byte changes; the projection is
  already byte-identical (V15).

---

## 7. Velocity and risk

**Velocity.** ~95 changed lines across 4 files. Measured cost centres in this worktree:
`verify_ail.sh` **6.0 s**, `go build ./...` **1.5 s**, `go test ./host/runbook/` **1.4 s**,
scoped verifygate test **~0.3 s**. The full acceptance sweep is well under a minute; the
eight-arm drill (one of which runs the 6 s gate) is a few minutes. **The estimate's risk is not
in the typing.** It is in (a) transcribing the two verbatim blocks byte-exactly, and (b) reading
counted lines instead of exit codes. Budget the 0.1 day there. Reference band: row 42 (`58c8f7f`)
landed a comparable documentation-plus-one-test item inside 0.1 d.

**Risks.**

| Risk | Mitigation |
|---|---|
| **AC1 read as rc** — the naive `-run` form is rc=0 at base with zero tests run | Every scoped run in this plan counts `=== RUN` / `--- PASS` / `--- FAIL` from `-v`. AC1 ships a paired nonsense-pattern control to prove the counter can read a zero |
| **AC2 asserted at 25 lines** (D1) | The plan never types a count; it diffs the landed block against the block extracted from the design doc |
| **A rig-absolute path in the new test file** (D2) | Forbidden in M1; caught by the whole-package run in M4 (`grep -c 'rig-absolute path'`) |
| **Box-drawing/alignment corruption when retyping the block** | Extract with `awk` from the design doc, never retype; `diff` the landed block against the extract |
| **`git checkout --` reflex during the drill** | §5.1 states the consequence plainly: the deliverable is uncommitted, so checkout deletes it. `cp` from `/tmp/p43-backup` only |
| **Pointer line breaks the Python heredoc** (D3) | Column-0 `#`, matching `:273`; the full pinned gate re-run in M2 would die loudly on Leg 1 if it did not |
| **Blanket `go test ./...` reds on an unrelated known base flake** | This sprint's gates are **narrow by design**: `./host/verifygate/` and `./host/runbook/` only. Do not widen to `./...`; a red there is not this diff |

---

## 8. Open questions for the human

Both are the design doc's own Open Decisions, carried forward at their defaults. Neither blocks
execution; if the reviewer overrules either, the text moves unchanged.

- **OD-1** — §S8 as a top-level section (default) vs an S6 subsection.
- **OD-2** — should the test also pin the runbook repair? **Default NO**, per §5.3.

---

**SPRINT_PLAN_PATH**: `design_docs/planned/w-floor-raise-coupling-inventory-sprint-plan.md`
**SPRINT_JSON_PATH**: `design_docs/planned/sprint_w-floor-raise-coupling-inventory.json`
