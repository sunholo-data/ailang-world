# w-inventory-test-blind-to-asymmetric-addition

**Status**: **READY TO ROUTE** — quorum round 1 BLOCKED, one designer revision, round 2 BLOCKED at full strength, then resolved by a bounded **narrow-refinement carve-out** applying all three reviewers' fixes verbatim (see §Quorum Verification Log). No objection disputes the design direction; every applied fix carries its own AC and its own named mutation per the mission's iteration-98 guardrail.
**Target**: World iter-141 / queue row 51
**Priority**: P2
**Estimated**: ~0.1d
**Dependencies**: None (row 51 is gated on nothing)
**Planner-Lane**: codex-ok
**Filed by**: the `sonnet` evaluator of row 43 (iter-132), non-blocking finding

## Problem Statement

`TestFloorRaiseInventoryNamesEveryCoupledFile` (`host/verifygate/floor_raise_inventory_test.go`)
is the gate that keeps the two homes of the floor-raise coupling inventory — the bounded
`# ── FLOOR-RAISE COUPLING INVENTORY` block in `scripts/verify_ail.sh` and the `## S8` table in
`design_docs/coding-standards.md` — in sync with the six named Tier-1 sites. It is **blind to
asymmetric addition**: it asserts per-row *uniqueness* of each of the six known rows in each
home, but it never counts rows, never compares the two homes to each other, and never bounds
cardinality.

The row's finding, reproduced first-party by the controller and re-derived by me in this
worktree (V2, V3): a **fabricated seventh coupled-file row** added to the `scripts/verify_ail.sh`
home ONLY — absent from `coding-standards.md` §S8 — leaves the test **GREEN** (rc=0, `--- PASS`),
and no count moves. The same is true in the reverse direction: a fabricated seventh row added to
§S8 ONLY also leaves the test GREEN. The gate is blind in both directions.

**THE ITEM:** decide whether the two homes must agree on their SITE SET (not merely on the
presence of each known site), and if so bind it. The row names three candidate mechanisms —
count equality between the homes, a set-difference assertion, a cardinality floor with an
instrument-failure branch. It also names the trap in the obvious fix: an assertion derived from
ONE home's contents cannot detect THAT home being wrong, so the comparison must be between the
two homes, or against an independently authored list.

**ALSO IN SCOPE:** the shared bare-token needle layer (`REQUIRED_VERIFIED`, `EXACT_TOTAL_VERIFIED`
and the rest) still permits duplication — it is checked with `strings.Contains` only (presence,
never count). It is explicitly secondary to the row anchors and no attack exploited it, but it
is the residue of the very class row 43 closed, left standing one layer down.

**THE GENERALISATION:** *a gate hardened by deletion is hardened against deletion; the thing it
has never been shown to notice is the thing that was added while nobody was looking.*

## Evidence (first-party measurements, with commands)

All measurements below were re-derived by me in this worktree at `dev` = `2c698da` (pristine,
`git status --porcelain` = 0 lines) on 2026-08-31. The controller's first-party measurements
(ARM A / ARM B / restore / cardinalities) are reproduced and confirmed; the sha256 values I
observed differ from the controller's in the ARM-A case because my insertion anchor produced a
different byte layout — the *behaviour* (mutation landed, count moved, test stayed GREEN) is
identical, which is the load-bearing fact.

**Pristine control** (both arms bracketed by it):

```zsh
export AILANG_BIN=/tmp/ailang-v0300/ailang
go test ./host/verifygate/ -run '^TestFloorRaiseInventoryNamesEveryCoupledFile$' -count=1 -v
# observed: rc=0, `--- PASS: TestFloorRaiseInventoryNamesEveryCoupledFile`, `ok ... host/verifygate`
```

**ARM A — fabricated 7th row added to the `scripts/verify_ail.sh` home ONLY** (inserted
`#   7. host/store/some_new_file.go   fabricated coupled site (repro)` immediately before the
`#   6.` line, inside the bounded inventory block):

```zsh
cp scripts/verify_ail.sh /tmp/verify_ail_backup.sh
shasum -a 256 scripts/verify_ail.sh | cut -c1-8        # before: 5a1bbe89
perl -0pi -e 's/(#   6\. host\/verifygate\/module_manifest_gate_test\.go)/#   7. host\/store\/some_new_file.go   fabricated coupled site (repro)\n$1/' scripts/verify_ail.sh
shasum -a 256 scripts/verify_ail.sh | cut -c1-8        # after:  f97c4ec9  (mutation LANDED)
sed -n '/── FLOOR-RAISE COUPLING INVENTORY/,/── END FLOOR-RAISE/p' scripts/verify_ail.sh | grep -cE '^#   [0-9]+\. '
# observed: 7  (intended effect asserted against the system's own view: 6 -> 7)
bash -n scripts/verify_ail.sh; echo "rc=$?"            # rc=0 (mutant is still a valid script)
go test ./host/verifygate/ -run '^TestFloorRaiseInventoryNamesEveryCoupledFile$' -count=1 -v
# observed: rc=0, `--- PASS`. THE GATE IS BLIND.
```

**ARM B — the SYMMETRIC arm, which the row does not record and which the controller added**:
fabricated 7th row added to `design_docs/coding-standards.md` §S8 ONLY
(`| 7 | \`host/store/some_new_file.go\` | fabricated coupled site (repro) |` before the `| 6 |`
row):

```zsh
cp design_docs/coding-standards.md /tmp/standards_backup.md
shasum -a 256 design_docs/coding-standards.md | cut -c1-8   # before: b710a510
perl -0pi -e 's/(\| 6 \| `host\/verifygate\/module_manifest_gate_test\.go`)/| 7 | `host\/store\/some_new_file.go` | fabricated coupled site (repro) |\n$1/' design_docs/coding-standards.md
shasum -a 256 design_docs/coding-standards.md | cut -c1-8   # after:  63208d62  (mutation LANDED)
sed -n '/## S8/,/^## /p' design_docs/coding-standards.md | grep -cE '^\| [0-9]+ \| '
# observed: 7  (intended effect: the §S8-scoped table-row count went 6 -> 7)
go test ./host/verifygate/ -run '^TestFloorRaiseInventoryNamesEveryCoupledFile$' -count=1 -v
# observed: rc=0, `--- PASS`. THE GATE IS BLIND IN THIS DIRECTION TOO.
```

**Restore** (both files restored byte-identical, pristine control re-passed):

```zsh
cp /tmp/standards_backup.md design_docs/coding-standards.md
cp /tmp/verify_ail_backup.sh scripts/verify_ail.sh
shasum -a 256 design_docs/coding-standards.md | cut -c1-8   # b710a510 (byte-identical)
shasum -a 256 scripts/verify_ail.sh | cut -c1-8              # 5a1bbe89 (byte-identical)
git status --porcelain | wc -l                                # 0 lines
go test ./host/verifygate/ -run '^TestFloorRaiseInventoryNamesEveryCoupledFile$' -count=1
# observed: rc=0, `ok ... host/verifygate` (post-restore pristine control re-PASSED)
```

**Current cardinalities, measured:**

```zsh
sed -n '/── FLOOR-RAISE COUPLING INVENTORY/,/── END FLOOR-RAISE/p' scripts/verify_ail.sh | grep -cE '^#   [0-9]+\. '
# observed: 6  (inventory block in scripts/verify_ail.sh)
sed -n '/## S8/,/^## /p' design_docs/coding-standards.md | grep -cE '^\| [0-9]+ \| '
# observed: 6  (§S8 table in design_docs/coding-standards.md)
grep -cE '^#   [0-9]+\. ' scripts/verify_ail.sh
# observed: 6  (whole-file count == block count today, so a whole-file count is NOT a safe
#               proxy for the block count — a design that uses one must say why)
```

**Structural facts about the test** (read the file, not taken on trust):
- It bounds both homes correctly: requires the begin/end markers to appear exactly once and be
  ordered, and slices §S8 from its heading to the next `## ` (S8 is the last section in
  `coding-standards.md`, so the slice runs to EOF — V10).
- For each of the 6 known rows in each home it asserts `strings.Count(...) == 1` — a UNIQUENESS
  assertion per known row, which is what row 43 added.
- `sharedNeedles` (8 bare tokens) are checked with `strings.Contains` only — presence, never
  count (V11). That is the "duplication still permitted" residue named in the row.
- NOTHING in the test counts rows, compares the two homes to each other, or bounds cardinality.
- `repoRoot` is a package-level helper already used by sibling gate tests in the same package.

**Gate baselines (base-state facts the design must be written around):**

```zsh
export AILANG_BIN=/tmp/ailang-v0300/ailang
./scripts/verify_ail.sh > /tmp/val 2>&1; echo "rc=$?"
# observed: rc=0  (usable rc=0 criterion)
./scripts/verify_go.sh > /tmp/vgo 2>&1; echo "rc=$?"
grep -E '^--- FAIL' /tmp/vgo | sort
grep -E '^(ok|FAIL).*host/verifygate' /tmp/vgo
# observed: rc=1, with EXACTLY ONE failing test:
#   --- FAIL: TestHandlerTimeoutKillsTheWholeProcessGroup (0.99s)   [host/broker]
#   ok  github.com/sunholo-data/ailang-world/host/verifygate
# So `scripts/verify_go.sh rc=0` is FORBIDDEN as an acceptance criterion. The criterion is a
# SET comparison: the failing-test set must not GROW relative to the base set recorded here.
```

## Solution Design

### Decision 1 — the two homes MUST agree on their SITE SET

**Yes.** The row's whole point is that the gate is blind to asymmetric addition. Requiring only
"each known site is present in each home" (the current per-row uniqueness) cannot see a site
that was *added* to one home and not the other. The two homes must agree on their site set, not
merely on the presence of each known site. This is the decision; it is not left open.

### Decision 2 — binding mechanism: set-equality between the homes, plus a cardinality floor

The row names three candidates. I choose a **composition** of four guards:

1. **Symmetric set-difference assertion (primary):** extract the canonical site set from each
   home and assert the two sets are **equal** by reporting, in BOTH directions, anything present
   in one home and absent from the other. This is the strongest of the three candidates: it
   subsumes count equality (equal sets ⟹ equal cardinality) and additionally catches a case
   count equality would miss — both homes having the same *number* of rows but *different* sites
   (e.g. a row replaced in one home). It is the direct binding of "the two homes must agree on
   their site set." The comparison is genuinely symmetric: a site added to EITHER home alone is
   reported by the loop over the OTHER home, so the divergent site is always NAMED.
2. **Duplicate-within-a-home guard (anti-evasion):** before comparing, assert no path repeats
   within EITHER home. A duplicated coupled-site row is a defect in its own right, and it is
   what makes a one-directional membership check evadable (a duplicate in one home can mask an
   asymmetric addition by keeping the counts equal). This closes the duplicate-evasion class.
3. **Cardinality floor with an instrument-failure branch (anti-vacuity):** assert
   `len(scriptSites) >= 6` AND `len(s8Sites) >= 6`. This is the guard that makes the
   set-equality non-vacuous: if the enumerator regex breaks and matches nothing, both sets are
   empty, the set-equality passes vacuously, and the floor REDS loudly. The floor is a **minimum**
   (≥ 6), never an equality (== 6), because the legitimate floor-raise path adds a 7th row to
   both homes and must PASS.

**Why this composition over the alternatives:**
- *Count equality alone* is necessary but not sufficient — it cannot see two homes with equal
  counts but divergent site sets. Rejected as the sole mechanism.
- *Symmetric set-difference alone* is sufficient for the asymmetric-addition class but is
  vacuous if the enumerator matches nothing (empty sets are equal). It needs the floor. Chosen
  as the primary, with the floor and the duplicate guard.
- *Cardinality floor alone* (a bare `== 6` or `>= 6`) cannot see a site *replaced* in one home
  (count stays 6, set changes) and, if written as `== 6`, would red on the legitimate 7th-row
  floor-raise. Rejected as the sole mechanism.

The existing per-row uniqueness assertions (the hardcoded six-site lists `scriptRows`/`s8Rows`)
are **kept byte-unchanged**. They are the independent authority for the six named Tier-1 sites
and they preserve the removal-detection behaviour.

### Decision 3 — the trap, cleared explicitly

The trap: an assertion derived from ONE home's contents cannot detect THAT home being wrong. My
design is a **comparison between the two homes** (symmetric set-equality), which is one of the
two forms the row explicitly permits. The comparison is symmetric in BOTH directions: a site
added to the script home alone is reported by the loop over the S8 set, and a site added to the
S8 home alone is reported by the loop over the script set — so whichever home is wrong, the
divergent site is NAMED. The independent authority for the six known sites is the **hardcoded
per-row list already in the test** (`scriptRows`/`s8Rows`), which pins each of the six by name
and uniqueness. The set-equality binds the *relationship* between the homes; the per-row list
binds the *identity* of the six known sites. Neither is derived from the other home's contents in
a way that would let one home's error hide itself.

**The residual this model deliberately accepts:** a *symmetric* wrongness — a fabricated 7th row
added to BOTH homes consistently — is indistinguishable from a legitimate floor-raise and passes.
This is inherent to the "two homes must agree" model and is disclosed in Declared Residuals. The
row's scope is the six named Tier-1 sites; the legitimate floor-raise path *requires* the test to
tolerate a consistent 7th addition, so the test cannot hardcode "exactly six sites."

### Decision 4 — `sharedNeedles`: presence-only is CORRECT, with a stated rationale

**Presence-only is correct for the shared bare-token needles.** Rationale, measured (V9):

- The sharedNeedles are **content markers, not structural anchors**. They are bare tokens that
  legitimately appear MORE than once in the bounded homes: `REQUIRED_VERIFIED` appears 2× in the
  script block (row 3's description "BOTH constants: REQUIRED_VERIFIED and EXACT_TOTAL_VERIFIED"
  AND the prose "a new identity in REQUIRED_VERIFIED below"); `EXACT_TOTAL_VERIFIED` appears 2× in
  the script block and 2× in §S8. A `count == 1` assertion on these would be a **false RED** on
  the pristine homes.
- The structural integrity of the inventory is enforced by the **row-anchor uniqueness
  assertions** (the six per-row `count == 1` checks) and, after this change, by the **set-equality
  between the homes**. The sharedNeedles are a coarse "the content is still present" check layered
  on top; they are not the mechanism that detects deletion or asymmetric addition.
- The duplication residue the row names is real but **harmless in the direction that matters**:
  duplicating a bare token like `interfaceHash` in prose changes nothing structural. The dangerous
  duplication class — a row anchor appearing twice — is already covered by `count == 1` per row.

So the decision is a deliberate, examined "presence is correct here," not an unexamined one. The
sharedNeedles layer is left as `strings.Contains` only, and this is recorded as a declared
residual (a future row could promote a specific needle to a count if a real attack exploits it).

### Implementation sketch (lives in `host/verifygate/floor_raise_inventory_test.go`)

Add two package-level regexes and a helper that extracts the canonical site set from a home's row
lines, then assert set-equality and the cardinality floor. The existing begin/end-marker bounds
and the §S8 heading-to-next-`##` slice are reused unchanged.

```go
var (
	// siteReScript matches `#   N. path` rows in the verify_ail.sh inventory block.
	siteReScript = regexp.MustCompile(`^#   [0-9]+\.\s+(\S+)`)
	// siteReTable matches `| N | `path` |` rows in the §S8 table; capture up to the
	// next pipe and trim backticks/space so the canonical path is the bare file path.
	siteReTable = regexp.MustCompile(`^\| [0-9]+ \| ([^|]+)`)
)

// canonicalSiteSet extracts the ordered set of coupled-site file paths from a home's
// row lines. scriptHome=true parses `#   N. path` rows; false parses `| N | `path` |`
// table rows. It returns the paths in row order; callers compare them as sets.
func canonicalSiteSet(home string, scriptHome bool) []string {
	re := siteReTable
	if scriptHome {
		re = siteReScript
	}
	var sites []string
	for _, line := range strings.Split(home, "\n") {
		m := re.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		sites = append(sites, strings.Trim(m[1], "` \t"))
	}
	return sites
}
```

In the test body, after the existing per-row uniqueness loops, add:

```go
scriptSites := canonicalSiteSet(block, true)
s8Sites := canonicalSiteSet(s8, false)

// Duplicate-within-a-home guard: a duplicated coupled-site row is a defect in its own
// right, and it is what makes a one-directional membership check evadable (a duplicate
// in one home can mask an asymmetric addition by keeping the counts equal). Assert no
// path repeats within EITHER home before comparing.
scriptSet := make(map[string]bool, len(scriptSites))
for _, s := range scriptSites {
	if scriptSet[s] {
		t.Fatalf("duplicate coupled-site path %q in the verify_ail.sh inventory block — a duplicated row is a defect", s)
	}
	scriptSet[s] = true
}
s8Set := make(map[string]bool, len(s8Sites))
for _, s := range s8Sites {
	if s8Set[s] {
		t.Fatalf("duplicate coupled-site path %q in the S8 table — a duplicated row is a defect", s)
	}
	s8Set[s] = true
}

// Cardinality floor with instrument-failure branch: if the enumerator matches nothing,
// both sets are empty and equal, so the set-equality below would pass vacuously. The
// floor makes that a LOUD red. It is a MINIMUM (>= 6), never an equality, because the
// legitimate floor-raise path adds a 7th row to BOTH homes and must pass.
if len(scriptSites) < 6 {
	t.Fatalf("inventory block enumerator matched %d sites, want >= 6 — instrument failure, not an empty inventory", len(scriptSites))
}
if len(s8Sites) < 6 {
	t.Fatalf("S8 enumerator matched %d sites, want >= 6 — instrument failure, not an empty inventory", len(s8Sites))
}

// Symmetric set-difference assertion: the two homes must agree on their SITE SET, not
// merely on the presence of each known site. Report anything in EITHER home absent from
// the OTHER — both directions — so the divergent site is always NAMED. This is the
// comparison between the two homes that detects asymmetric addition.
for _, s := range s8Sites {
	if !scriptSet[s] {
		t.Errorf("S8 site %q absent from the verify_ail.sh inventory block — the two homes disagree on their site set", s)
	}
}
for _, s := range scriptSites {
	if !s8Set[s] {
		t.Errorf("verify_ail.sh site %q absent from the S8 table — the two homes disagree on their site set", s)
	}
}
if len(scriptSites) != len(s8Sites) {
	t.Errorf("site-set cardinality mismatch: verify_ail.sh has %d sites, S8 has %d — the two homes disagree on their site set", len(scriptSites), len(s8Sites))
}
```

The `t.Fatalf` on the floor is deliberate: an empty enumerator result is an instrument failure,
not a legitimate empty inventory, and must stop the test loudly. The set-equality uses `t.Errorf`
so a mismatch reports every divergent site, not just the first. The duplicate guard uses
`t.Fatalf` because a duplicated row is a defect that must stop the test before any comparison is
meaningful.

**Declared Precondition (extraction assumptions, measured — V12).** The two enumerator regexes
rest on two assumptions the reviewers named, both of which hold today by first-party measurement
at pristine `dev` = `2c698da` (V12):

- `(\S+)` in `siteReScript` assumes **no coupled-site path contains whitespace**. All six real
  paths are whitespace-free; the regex captures each to the end of the line.
- `([^|]+)` + backtick-trim in `siteReTable` assumes **the §S8 path column is backtick-wrapped**.
  All six real rows are `` `path` ``; the capture runs to the next pipe and the trim strips the
  backticks.

Both hold today by measurement. A future site that violates either would produce a **FALSE RED on
pristine** — the enumerator would fail to extract that site, the two sets would diverge, and the
gate would red loudly. That is a loud, correct-to-notice failure, not a silent one: it forces the
regex to be updated in the same commit that adds the site, which is exactly the coupling the gate
is meant to enforce.

**Why this stays in `floor_raise_inventory_test.go`:** the fix is a test-only edit to the host
verify surface, touching no kernel type, no transition law, and exporting no API. There is no
package-shaped surface here. No new package, no new script, no new CI job — all out of scope for
a ~0.1d row.

## Alternatives Considered (and why rejected)

| Alternative | Why rejected |
|---|---|
| **Count equality between the homes only** (`len(scriptSites) == len(s8Sites)`) | Necessary but not sufficient: cannot see two homes with equal counts but divergent site sets (a row *replaced* in one home). Chosen as a secondary check inside the set-equality, not as the primary mechanism. |
| **Cardinality floor only** (bare `== 6` or `>= 6`) | A bare `== 6` reds on the legitimate 7th-row floor-raise (unusable — S8's own recipe requires exactly that). A bare `>= 6` cannot see a site *replaced* in one home (count stays 6, set changes). Rejected as the sole mechanism. |
| **Set-difference against an independently authored hardcoded list of the full set** | The legitimate floor-raise path adds a 7th site, so the full set is not fixed; a hardcoded "exactly these six" would red on the legitimate path. The hardcoded per-row list is kept as the authority for the six *known* sites, and the set-equality binds the *relationship* between the homes. |
| **Promote `sharedNeedles` to `count == 1`** | Measured false RED: `REQUIRED_VERIFIED` appears 2× in the script block and `EXACT_TOTAL_VERIFIED` 2× in both homes (V9). The needles are content markers, not structural anchors; the structural integrity is enforced by the row-anchor uniqueness + set-equality. Rejected; presence-only retained with rationale. |
| **A new package / script / CI job** | Out of scope for ~0.1d. The fix is a test-only edit to one file. |

## Acceptance Criteria

Every AC names a command and an observable that can FAIL. Run-existence is asserted on `=== RUN`
/ `--- PASS` counts, never exit codes alone (an rc=0 `[no tests to run]` is a FALSE PASS). Shell
recipes follow the rig rules: capture with `cmd > /tmp/out 2>&1; rc=$?` (zsh `${PIPESTATUS[0]}`
is empty), never `|| echo 0` inside `$(...)`, `export PATH=/opt/homebrew/bin:$PATH` first.
**Doc-wide convention:** every command in this document is copy-paste runnable under zsh exactly
as printed, which is why runnable commands live in fenced code blocks and never in markdown
table cells — a bare `|` breaks a table cell, and the `\|` escape that fixes the rendering
silently corrupts the command.

The test command used throughout (call it **T**):

```zsh
export AILANG_BIN=/tmp/ailang-v0300/ailang
go test ./host/verifygate/ -run '^TestFloorRaiseInventoryNamesEveryCoupledFile$' -count=1 -v > /tmp/t 2>&1; rc=$?
```

- **AC1 (run-existence floor).** Run **T** on the pristine tree → rc=0 AND
  `grep -c '=== RUN' /tmp/t` ≥ 1 AND `grep -c -- '--- PASS' /tmp/t` ≥ 1. A `[no tests to run]`
  line anywhere in `/tmp/t` FAILS this AC. A green exit code over zero executed tests is the
  false pass this repo has been burned by.
- **AC2 (set-equality assertion is present and load-bearing).** `grep -c 'disagree on their site set' host/verifygate/floor_raise_inventory_test.go` ≥ 3 (the S8-absent Errorf, the script-absent Errorf, and the cardinality-mismatch Errorf). Load-bearing proof is mutation arm N1: with the set-equality neutered (both membership Errorfs AND the cardinality-mismatch Errorf) and ARM A1 landed, the test goes GREEN (see Mutation Drill).
- **AC3 (cardinality floor is present and load-bearing).** `grep -c 'instrument failure, not an empty inventory' host/verifygate/floor_raise_inventory_test.go` = 2 (one per home). Load-bearing proof is the N2 pair: with the floor neutered AND both enumerators neutered (I1+I2), the test goes GREEN (the vacuous pass the floor exists to prevent); with the SAME both-enumerators-neutered mutant and the floor PRESENT, the test goes RED with the `instrument failure, not an empty inventory` Fatalf.
- **AC4 (asymmetric addition to the script home reds — the row's point).** With ARM A1 landed (protocol below): run **T** → rc=1 AND `grep -c -- '--- FAIL' /tmp/t` = 1 AND the failure message names the divergent site (`host/store/some_new_file.go`). A1 and A2 are now SYMMETRIC: both red, both naming the divergent site.
- **AC5 (asymmetric addition to the S8 home reds).** With ARM A2 landed: run **T** → rc=1 AND `grep -c -- '--- FAIL' /tmp/t` = 1 AND the failure message names the divergent site. Symmetric to AC4: the reverse-direction loop names the site added to the S8 home alone.
- **AC6 (consistent addition to BOTH homes passes — the legitimate floor-raise path).** With ARM A3 landed: run **T** → rc=0 AND `grep -c -- '--- PASS' /tmp/t` = 1. A design that reds here is unusable, because S8's own recipe requires exactly this.
- **AC7 (a row outside the bounded block / outside §S8 does NOT affect the verdict).** With ARM A4 landed: run **T** → rc=0 AND `grep -c -- '--- PASS' /tmp/t` = 1.
- **AC8 (removal arms still red — existing behaviour must not regress).** With ARM R1 landed (delete a row from the script home): run **T** → rc=1. With ARM R2 landed (delete a row from §S8): run **T** → rc=1. The per-row uniqueness assertions must still fire.
- **AC9 (instrument-failure arms fail LOUDLY, never pass silently).** With ARM I1 landed (neuter the script enumerator so it matches nothing): run **T** → rc=1 AND the message is the `instrument failure, not an empty inventory` Fatalf. With ARM I2 landed (neuter the S8 enumerator): run **T** → rc=1 with the same message.
- **AC10 (pristine control).** On the untouched tree (sha256 equal to the V1/V4 baselines): run **T** → rc=0 AND `grep -c -- '--- PASS' /tmp/t` = 1, before AND after the drill.
- **AC11 (hygiene).** `gofmt -l host/verifygate/` → 0 bytes; `go vet ./host/verifygate/` rc=0; `go build ./...` rc=0.
- **AC12 (gates).** Three parts, under the pinned `AILANG_BIN=/tmp/ailang-v0300/ailang`:
  1. `./scripts/verify_ail.sh` rc=0 — green at pristine base (V5), so the exit-code criterion stands.
  2. `./scripts/verify_go.sh` — an **rc=0 criterion is FORBIDDEN** (the gate is RED at pristine base on this rig, V6). The criterion is a SET comparison: the set of `--- FAIL` test names after the change is **identical to the base set**, AND `host/verifygate` reports `ok` in the same run. Base set (rig-local, measured at pristine `dev` = `2c698da`): `TestHandlerTimeoutKillsTheWholeProcessGroup` (host/broker). Any name appearing or disappearing fails this AC.
  3. The narrowest gate that can actually fail for this diff: `go test ./host/verifygate/` rc=0.
- **AC13 (duplicate-within-a-home reds, names the duplicate, and is ATTRIBUTED to the new guard).** With ARM D1 landed: run **T** → rc=1 AND `/tmp/t` contains `duplicate coupled-site path "host/store/some_new_file.go" in the verify_ail.sh inventory block` AND — this is the attribution half — `grep -c 'want exactly 1' /tmp/t` = **0**, proving the legacy `scriptRows` per-row loop did not fire. With ARM D2 landed: the same three conditions against the S8 message. **The duplicated row MUST be a fabricated 7th path, never a duplicate of one of the six known rows**: the pre-existing per-row `strings.Count(...) == 1` loops (`floor_raise_inventory_test.go:40-43` and `:85-88`, published at V13) intercept a duplicated KNOWN row, so that arm would red for the legacy reason and prove nothing about the new guard (`gemini-3-1-pro`, round 2, applied verbatim under the narrow-refinement carve-out). This is the anti-evasion half of the comparator: a duplicated coupled-site row is a defect in its own right, and it is what makes a one-directional membership check evadable.
- **AC14 (the extraction symbols remain collision-free — guards the applied Conflict-Surface fix).**
  `grep -rn 'siteReScript\|siteReTable\|canonicalSiteSet' --include='*.go' . | grep -v floor_raise_inventory_test.go | wc -l` = **0**, run beside the same-scope known-positive control
  `grep -rn 'regexp.MustCompile' --include='*.go' host/ | wc -l` ≥ **9** (a zero there means the instrument is broken, not that the repo is clean — rule 3a(i-d)). **Named mutation C1:** declare a second `siteReScript` in any other `host/` file → the first count must move 0→≥1 and the package must stop building. The AC that fails if the Conflict-Surface fix is absent is this one; the mutation that reds if it is neutered is C1.
- **AC15 (the per-row uniqueness authority is still present and still load-bearing — guards the applied Decision-3 verification).**
  `grep -c 'want exactly 1' host/verifygate/floor_raise_inventory_test.go` = **2** (one loop per home) AND `grep -c 'scriptRows :=\|s8Rows :=' host/verifygate/floor_raise_inventory_test.go` = **2**. **Named mutation C2:** delete one entry from the hardcoded `scriptRows` list → ARM R1 (delete that same row from the script block) must stop redding on the per-row message. Decision 3 names these lists as the independent authority for the six known sites; without this AC nothing re-derives that they exist (`oc-glm-5-2`, round 2, applied under the carve-out).

## Mutation Drill

Per the doc convention, the table cells reference named arms; the commands are in the fenced
blocks below, runnable as printed. **Drill protocol (binding on the sprint), for EVERY arm:**
(1) assert the mutant LANDED by sha256 AND by an intended-effect query against the system's own
view (a `grep -c` that must MOVE — per-arm gates below), never against the file's bytes alone;
(2) assert the tree still builds (`go vet ./host/verifygate/` rc=0) BEFORE reading any test
result; (3) restore byte-identical from a `cp` BACKUP (never `git checkout --`), verify by sha256,
and re-run the pristine control (AC10) before AND after every arm. The verdict of every arm is
the NAME of the failing test and the message that fired, never an exit code alone.

| Arm | Mutant (target) | Landed/effect gate | Expected post-fix | Kills which mutation / proves what |
|---|---|---|---|---|
| A1 | **ADDITION**: 7th row `#   7. host/store/some_new_file.go   fabricated coupled site (repro)` in the `scripts/verify_ail.sh` block ONLY | G1 sha moves; G2 block count 6→7; G3 rc=0 | **RED**, `--- FAIL`, message names `host/store/some_new_file.go` (AC4) | The row's point — the reverse-direction loop names the site added to the script home alone |
| A2 | **ADDITION**: 7th row `\| 7 \| \`host/store/some_new_file.go\` \| fabricated coupled site (repro) \|` in §S8 ONLY | G4 sha moves; G5 S8 count 6→7; G3 rc=0 | **RED**, `--- FAIL`, message names `host/store/some_new_file.go` (AC5) | The symmetric direction — the forward loop names the site added to the S8 home alone; proves the comparison is between the two homes, not one home against itself |
| A3 | **ADDITION**: 7th row added to BOTH homes consistently | G1 AND G4 sha move; G2 AND G5 counts 6→7; G3 rc=0 | **PASS**, rc=0 `--- PASS` (AC6) | The legitimate floor-raise path — a design that reds here is unusable |
| A4 | **ADDITION**: a `#   7. ...` row added OUTSIDE the bounded block in `scripts/verify_ail.sh` (e.g. after the END marker) | G1 sha moves; G2 block count STAYS 6; G3 rc=0 | **PASS**, rc=0 (AC7) | Proves the extraction is bounded to the block — a row outside must not affect the verdict |
| D1 | **DUPLICATE**: insert `#   7. host/store/some_new_file.go   fabricated coupled site (repro)` **TWICE** inside the script block. It MUST be a fabricated 7th path, never a duplicate of a known row: the pre-existing `scriptRows` loop asserts `strings.Count(block, needle) == 1` for each of the six KNOWN rows (`floor_raise_inventory_test.go:40-43`), so duplicating a known row is intercepted by the legacy assertion and the arm reds for the wrong reason (`gemini-3-1-pro`, round 2, verbatim `proposed_fix`) | G1 sha moves; G2 block count 6→8; G3 rc=0 | **RED**, `duplicate coupled-site path "host/store/some_new_file.go" in the verify_ail.sh inventory block` Fatalf (AC13). The legacy per-row loop must NOT appear in the output — its absence is what proves the new guard fired autonomously | The anti-evasion half — a duplicated row in one home is a defect in its own right and is what makes a one-directional membership check evadable. The fabricated path bypasses the hardcoded `scriptRows` checks, so the arm proves the NEW dynamic `scriptSet` guard, uninterfered with |
| D2 | **DUPLICATE**: insert the fabricated §S8 row for `host/store/some_new_file.go` **TWICE** inside §S8 (numbered 7 and 8). Fabricated for the same reason as D1: the pre-existing `s8Rows` loop asserts `strings.Count(s8, needle) == 1` per known row (`floor_raise_inventory_test.go:85-88`) | G4 sha moves; G5 S8 count 6→8; G3 rc=0 | **RED**, `duplicate coupled-site path "host/store/some_new_file.go" in the S8 table` Fatalf (AC13). The legacy per-row loop must NOT appear in the output | Same guard on the S8 home — the duplicate-within-a-home red is symmetric, and the fabricated path proves the new guard rather than the legacy one |
| R1 | **REMOVAL**: delete one row (e.g. `#   3. scripts/verify_ail.sh`) from the script block | G1 sha moves; G2 block count 6→5; G3 rc=0 | **RED** (AC8) | Existing removal-detection must not regress |
| R2 | **REMOVAL**: delete one row (e.g. `\| 3 \| \`scripts/verify_ail.sh\` \|`) from §S8 | G4 sha moves; G5 S8 count 6→5; G3 rc=0 | **RED** (AC8) | Existing removal-detection must not regress |
| I1 | **INSTRUMENT FAILURE**: neuter `siteReScript` (Go-side) so it matches nothing | test-file sha moves; G3 rc=0 | **RED**, `instrument failure, not an empty inventory` Fatalf (AC9) | Proves the floor catches an empty enumerator result — never a silent pass |
| I2 | **INSTRUMENT FAILURE**: neuter `siteReTable` (Go-side) so it matches nothing | test-file sha moves; G3 rc=0 | **RED**, `instrument failure, not an empty inventory` Fatalf (AC9) | Same, for the S8 enumerator |
| N1 | **NEUTER the set-equality** (delete the two `disagree on their site set` membership Errorfs AND the cardinality-mismatch Errorf) | test-file sha moves; G3 rc=0 | With ARM A1 simultaneously landed: **GREEN** — the previously-expected `--- FAIL` is GONE | Proves the set-equality is load-bearing, not shadowed by the floor. The cardinality Errorf must be neutered too: under A1 the counts differ (7 vs 6), so it would otherwise fire and mask the proof |
| N2a | **NEUTER the cardinality floor** (delete the two `instrument failure` Fatalfs) AND land I1 AND I2 together (both enumerators match nothing) | test-file sha moves; G3 rc=0 | **GREEN** — both sets empty, every membership loop vacuous, cardinality `0 == 0` | Proves the vacuous pass the floor exists to prevent: a single arm cannot isolate the floor, because with only ONE enumerator neutered the other home's six sites still fire the membership loop |
| N2b | **PAIRED CONTROL**: land I1 AND I2 together (both enumerators match nothing) WITH the floor PRESENT | test-file sha moves; G3 rc=0 | **RED**, `instrument failure, not an empty inventory` Fatalf | The paired control that isolates the floor: the SAME both-enumerators-neutered mutant, floor present, goes RED — proving the floor is what makes the instrument-failure arm loud |
| C1 | **CARVE-OUT GUARD (collision)**: declare a second `siteReScript` in another `host/` file | new file sha exists; `grep -rn 'siteReScript' --include='*.go' . | wc -l` moves 1→2 | **RED** — the package stops building (`go vet ./host/verifygate/` rc≠0), and AC14's count moves off 0 | Binds the applied §Conflict-Surface fix (`gpt5-6-sol` R2-1): proves AC14's zero is a live measurement and not a claim about a search nobody can break |
| C2 | **CARVE-OUT GUARD (authority)**: delete one entry from the hardcoded `scriptRows` list, then land ARM R1 (delete that same row from the script block) | test-file sha moves; `grep -c 'want exactly 1' host/verifygate/floor_raise_inventory_test.go` stays 2 but the list shrinks 6→5; `go vet` rc=0 | R1 must **stop redding on the per-row message** — the removal is then caught only by the cardinality floor, which is a weaker verdict | Binds the applied Decision-3 verification (`oc-glm-5-2` R2-3): proves the hardcoded per-row lists really are the independent authority for the six known sites, rather than a described one |

**Gate commands** referenced by the table (repo root, zsh, runnable as printed):

```zsh
shasum -a 256 scripts/verify_ail.sh | cut -c1-8        # G1 — script-home landed/restore assertion
sed -n '/── FLOOR-RAISE COUPLING INVENTORY/,/── END FLOOR-RAISE/p' scripts/verify_ail.sh | grep -cE '^#   [0-9]+\. '   # G2 — script block row count (must MOVE)
bash -n scripts/verify_ail.sh                          # G3 — mutant is valid shell (rc=0); for Go-side arms, `go vet ./host/verifygate/` rc=0 instead
shasum -a 256 design_docs/coding-standards.md | cut -c1-8   # G4 — S8-home landed/restore assertion
sed -n '/## S8/,/^## /p' design_docs/coding-standards.md | grep -cE '^\| [0-9]+ \| '   # G5 — §S8 table row count (must MOVE)
```

A `grep -c` gate that legitimately counts 0 exits rc=1 — read the COUNT, never the exit code,
and never wrap in `|| echo 0`.

## Conflict Surface (reuse audit — `gpt5-6-sol` round 2, applied verbatim)

The reviewer's objection was that the doc dismisses new packages and asserts the logic belongs in
the existing test **without ever searching for existing machinery to extend**. The audit had
genuinely never been run. It is run here, by the controller, in the repo at pristine `dev` =
`2c698da`; its **conclusion** ("there may be machinery to reuse") is REFUTED, and the fix is
therefore to PUBLISH the audit rather than to change the design. Every count below is paired with
a same-scope known-positive control, and the audit itself is bound by **AC14**.

```zsh
# A1 — do the three proposed symbols already exist anywhere?  EXPECT 0
grep -rn 'siteReScript\|siteReTable\|canonicalSiteSet' --include='*.go' . | wc -l
# observed: 0
# A1-control — same scope, a pattern that MUST hit (proves the instrument sees a positive)
grep -rn 'regexp.MustCompile' --include='*.go' host/ | wc -l
# observed: 9

# A2 — bounded-section / row-parsing machinery inside the target package
grep -rn 'regexp.MustCompile\|FindStringSubmatch' --include='*.go' host/verifygate/
# observed, 4 hits, all unsuitable:
#   ail_binary_gate_test.go:373  dev   := regexp.MustCompile(`(-dirty$|-[0-9]+-g[0-9a-f]+)`)
#   ail_binary_gate_test.go:374  shape := regexp.MustCompile(`^v[0-9]+\.[0-9]+\.[0-9]+(-[A-Za-z0-9.]+)?$`)
#   toolchain_pin_gate_test.go:124 jobLine := regexp.MustCompile(`^  ([a-z0-9-]+):$`)
#   toolchain_pin_gate_test.go:133   if match := jobLine.FindStringSubmatch(line); match != nil {

# A3 — any markdown table-row parsing anywhere in the repo
grep -rn 'strings.Split(.*"|"' --include='*.go' . | head
# observed, 1 hit: host/broker/invoke_boundary_test.go:179  parts := strings.Split(line, "|")

# A4 — existing "canonical"/"inventory" helpers
grep -rniE 'func .*(canonical|inventory)' --include='*.go' . | head
# observed: host/verifygate/toolchain_pin_gate_test.go:29 canonicalizeVersionPin(value string)
#           host/transitionreg/codec.go:296 canonicalSchema(raw []byte)
#           (plus test names and unrelated store/evidence helpers)

# A5 — third-party parsing libraries
sed -n '/^require/,/^)/p' go.mod
# observed: ONE direct dependency, modernc.org/sqlite v1.54.0; every other entry is `// indirect`.
#           Zero markdown parsers, zero shell parsers, zero table parsers.
```

**Verdict, hit by hit.**

| Hit | Reuse? | Why |
|---|---|---|
| `ail_binary_gate_test.go:373-374` | **No** | Version-string *shape* regexes over a `--version` token. Different grammar, no bounded section, no row enumeration, nothing to extend. |
| `toolchain_pin_gate_test.go:124/133` (`jobLine`) | **No** — and it is the closest analogue, so it is named rather than omitted | A line-anchored `FindStringSubmatch` over YAML job headers. Same *shape* (anchor a line, capture one group, walk `strings.Split(s,"\n")`), entirely different *grammar*, and it is an unexported local inside another test function — there is no helper to call. Extending it would mean generalising a two-line local into a parser, which is more machinery than the ~0.1d fix it would serve. |
| `host/broker/invoke_boundary_test.go:179` | **No** | The repo's only markdown-pipe split, in an unrelated boundary test in another package, splitting a protocol line rather than a table row. Not exported, not a table parser. |
| `canonicalizeVersionPin` / `canonicalSchema` | **No** | Name collision only. One normalises a *version string*, the other canonicalises *JSON*; neither touches file paths. |
| `go.mod` | **No** | One direct dependency. Adding a markdown or shell parsing library to enumerate six rows in two files is disproportionate and would enlarge the daemon dependency closure the repo deliberately bounds. |

**Conclusion:** there is nothing to extend; the two regexes and one helper are new, local, unexported
and collision-free (A1 = 0 against a control of 9). The route-to-extension bias is satisfied by
having looked, and the looking is now a committed AC (AC14) rather than a claim.

## Declared Residuals

1. **Symmetric wrongness passes.** A fabricated 7th row added to BOTH homes consistently is
   indistinguishable from a legitimate floor-raise and passes. This is inherent to the "two
   homes must agree" model and is the price of AC6 (the legitimate floor-raise path must pass).
   The row's scope is the six named Tier-1 sites; the hardcoded per-row list pins those six, and
   the set-equality binds the relationship. A future row that wants to bound the *absolute* set
   would need an independent authority for the full set — out of scope here. Note that the
   duplicate-evasion variant of this class (a duplicate in ONE home masking an asymmetric
   addition) is NOT a residual: it is closed by the duplicate-within-a-home guard (AC13, arms
   D1/D2).
2. **`sharedNeedles` duplication remains permitted.** The 8 bare tokens are still checked with
   `strings.Contains` only. This is a deliberate, examined decision (Decision 4): they are content
   markers that legitimately appear >1 times, and the structural integrity is enforced by the
   row-anchor uniqueness + set-equality. A future row could promote a specific needle to a count
   if a real attack exploits it.
3. **A row *replaced* in one home with a different path of the same count** is caught by the
   set-equality (the sets differ), but a row *renamed identically in both homes* is not — that is
   the same symmetric-wrongness class as residual 1.
4. **The §S8 slice runs to EOF** (S8 is the last section in `coding-standards.md`, V10). If a new
   `## ` section is ever added after S8, the slice would shrink and the S8 table rows would still
   be inside it — the existing heading-to-next-`##` logic already handles that; no change needed,
   but the assumption is recorded.
5. **`verify_go.sh` is RED at pristine base on this rig** (V6: exactly one failing test,
   `TestHandlerTimeoutKillsTheWholeProcessGroup` in host/broker, a 100 ms process-group-kill
   deadline under parallel-package load; dev CI is GREEN on the same commit). This is a rig-local
   environmental cost, not a code defect. AC12 is written around it (set comparison, never exit
   code); fixing the rig cost is a queue item, not this sprint's scope.

## Quorum Verification Log

**Round 1 — BLOCKED at full strength.** Three external reviewers plus the controller; zero
absent reviewers; cost $0.1368. All four objections land on ONE surface: the `## Solution
Design` implementation sketch's comparator, and the ACs / drill rows that depend on it. Nobody
disputes the design direction — set-equality between the two homes plus a cardinality floor is
accepted. Everything else (Problem Statement, Evidence, Decision 1, Decision 4, Alternatives,
AC12's base-red set-comparison, Declared Residuals) is ACCEPTED and preserved.

| # | Objection surface | Raised by | Verdict | Resolution in this revision |
|---|---|---|---|---|
| 1a | The comparator is not set equality — evadable under duplicates (a duplicate in one home masks an asymmetric addition) | gpt5-6-sol, gemini-3-1-pro, controller | **CONFIRMED** (controller-measured) | Comparator made a genuine symmetric set-difference in BOTH directions, plus a duplicate-within-a-home guard (AC13, arms D1/D2) |
| 1b | The comparator is not set equality — makes AC4 unsatisfiable (with ARM A1 landed, only the cardinality Errorf fires, which names counts, never the site) | gpt5-6-sol, gemini-3-1-pro, controller | **CONFIRMED** (controller-measured) | Reverse-direction loop added so ARM A1 reds with the site NAMED; AC4/AC5 both read "names the divergent site"; A1 and A2 are now symmetric |
| 2 | The claimed N2 result is impossible (neuter floor + land I1 alone → RED, not GREEN) | gpt5-6-sol | **CONFIRMED** (controller-measured) | N2 redesigned as an explicit pair N2a/N2b (both enumerators neutered, floor absent → GREEN; same mutant, floor present → RED); N1 re-checked and adjusted to neuter the cardinality Errorf too |
| 3 | The extraction regexes are never verified against the real files (AC1 rests on an unverified premise) | oc-glm-5-2 | **CONFIRMED-AS-DOC-GAP-BUT-CONCLUSION-REFUTED** | The reviewer's conclusion ("non-viable as written") is refuted by measurement — both homes extract to the identical 6-element set. The fix is to PUBLISH the verification: new V12 row re-derived first-party, plus a Declared Precondition paragraph recording the two assumptions |

### Round 2 — BLOCKED at FULL STRENGTH, then resolved under the NARROW-REFINEMENT CARVE-OUT

3 external reviewers + controller. Controller verdict **PASS**; all three external reviewers
**REJECT**. `absent_reviewers` initially carried `oc-glm-5-2` (reason `invalid` — a malformed-JSON
response whose raw text already read `"verdict":"reject"`); it was **RESTORED** by a single-reviewer
re-run at a raised cap (`ailang design-review --reviewer oc-glm-5-2 --max-cost-usd 0.30`, $0.0417,
verdict **reject**), so round 2 is read at full strength with the hole closed, exactly as the shared
skill requires. Round-2 quorum cost $0.1642.

| # | Objection | Raised by | Controller measurement | Applied fix |
|---|---|---|---|---|
| R2-1 | No verified conflict-surface / reuse analysis for the new regexes and helper | gpt5-6-sol | **CONCLUSION REFUTED, PREMISE UPHELD.** The audit had genuinely never been run; run first-party it finds nothing to extend (V14). This is the same shape World iteration 140 recorded from the same reviewer — *the audit is owed even when its conclusion is wrong* | New **§Conflict Surface** with the reviewer's own prescribed searches, every hit adjudicated, plus a Go-form V12 and the new **AC14** + mutation **C1** |
| R2-2 | D1/D2 duplicate a KNOWN row, so the pre-existing per-row `count == 1` loop intercepts the arm | gemini-3-1-pro | **CONFIRMED.** `floor_raise_inventory_test.go:40-43` and `:85-88` do assert `strings.Count(...) == 1` per known row. The reviewer's *"if it uses `t.Fatalf` the test aborts"* branch is FALSE — both loops use `t.Errorf` — but its actual claim, that the drill is confounded and fails to isolate the new machinery, is TRUE | D1/D2 rewritten to the reviewer's verbatim `proposed_fix`: duplicate a **fabricated** 7th row. AC13 gains the attribution assertion `grep -c 'want exactly 1' /tmp/t` = 0 |
| R2-3 | Decision 3's structural claims about the test file are asserted, never verified (only V11 touches it) | oc-glm-5-2 (restored) | **CONFIRMED as a doc gap; the underlying claim is TRUE.** Both hardcoded lists and both per-row loops exist exactly as described | New **V13** publishing the greps, plus **AC15** + mutation **C2** binding the per-row authority |

**Disposition — the narrow-refinement carve-out, not a park.** The bounded
one-revision-one-re-quorum allowance is spent. Every remaining blocking objection is a
**completeness / verification-publication / drill-attribution** defect; not one disputes the design
DIRECTION (all three accept two-homes-must-agree + a cardinality floor), and each carries a concrete
remedy — two as verbatim `proposed_fix` text, the third named explicitly in its `catch`. That is the
carve-out's precondition, met. Contrast World iteration 140, which parked on the same allowance
being spent: there the surviving objection attacked the direction and the ratified queue row behind
it. Here the fixes are measurements the controller had already taken or could take in one command.

Per this mission's own guardrail (iteration 98) — *a fix applied under the carve-out must acquire
its own acceptance criterion and its own named mutation, because it enters the document after the
round that reviewed it* — every applied fix is bound: R2-1 → **AC14** / mutation **C1**;
R2-2 → **AC13**'s attribution clause / arms **D1**, **D2**; R2-3 → **AC15** / mutation **C2**.
No reviewer objection was overridden, re-litigated, or resolved by controller invention.

## Verification Log

Rows V1–V11 were measured by me in this worktree at `dev` = `2c698da` (pristine, porcelain 0) on
2026-08-31. Row V12 was re-derived first-party in the revision pass (objection 3) at the same
pristine `dev` = `2c698da` — it applies the two proposed regexes to the two real homes and prints
both extracted lists plus the set-equality verdict. Per the doc-wide convention, **no runnable command sits in a table cell**: each row's
command is in the matching fenced block under "Commands", byte-for-byte as executed.

| # | Claim | Observed (command in block V-x below) |
|---|---|---|
| V1 | Pristine control: the test passes on the untouched tree | rc=0, `--- PASS: TestFloorRaiseInventoryNamesEveryCoupledFile`, `ok ... host/verifygate` |
| V2 | **ARM A (script-home-only 7th row) leaves the test GREEN — the gate is blind** | mutation LANDED sha `5a1bbe89` → `f97c4ec9`; block count 6→7; `bash -n` rc=0; test rc=0 `--- PASS` |
| V3 | **ARM B (S8-only 7th row) leaves the test GREEN — blind in this direction too** | mutation LANDED sha `b710a510` → `63208d62`; S8 count 6→7; test rc=0 `--- PASS` |
| V4 | Restore verified byte-identical; pristine control re-passes | shas back to `5a1bbe89` / `b710a510`; porcelain 0 lines; test rc=0 `ok` |
| V5 | `scripts/verify_ail.sh` is rc=0 at pristine base | rc=0, `verify gate PASSED: 11 required identities verified, 40 named tests pass` |
| V6 | `scripts/verify_go.sh` is rc=1 at pristine base with EXACTLY ONE failing test; `host/verifygate` is `ok` | rc=1; `--- FAIL: TestHandlerTimeoutKillsTheWholeProcessGroup (0.99s)`; `ok github.com/sunholo-data/ailang-world/host/verifygate` |
| V7 | Hygiene baselines | `go vet ./host/verifygate/` rc=0; `gofmt -l host/verifygate/` empty (0 bytes) |
| V8 | Whole-file count == block count today (6 == 6) — a whole-file count is NOT a safe proxy for the block count | `grep -cE '^#   [0-9]+\. ' scripts/verify_ail.sh` = 6; block-scoped count = 6 |
| V9 | `sharedNeedles` legitimately appear >1 times in the bounded homes — presence-only is correct | `REQUIRED_VERIFIED` block=2, s8=1; `EXACT_TOTAL_VERIFIED` block=2, s8=2; all others =1 in both |
| V10 | §S8 is the last `## ` section in `coding-standards.md` — the test's slice runs to EOF | `awk 'NR>=90 && /^## /'` prints only line 90 (`## S8 ...`) |
| V11 | `sharedNeedles` are checked with `strings.Contains` only, never `strings.Count` | `grep -n 'strings.Contains(block\|strings.Contains(s8'` → lines 58 and 93 only; no `strings.Count` on sharedNeedles |
| V12 | **The two extraction regexes produce IDENTICAL site sets on the pristine homes** — AC1's premise is verified, not assumed. Re-run in **Go**, using `regexp.MustCompile` and the exact proposed `canonicalSiteSet` semantics against the real bounded strings (`gpt5-6-sol` round 2 asked for the Go form; the Python form is kept as the cross-check) | both homes extract to the identical 6-element set in the same order: `world/<module>.ail`, `packages/world-core/world/<module>.ail`, `scripts/verify_ail.sh`, `scripts/world_package_ready_packet.golden.json`, `docs/SELF_MOD_PUBLISH.md`, `host/verifygate/module_manifest_gate_test.go`; `sets equal: true ; cardinality equal: true (6, 6)` |
| V13 | **Decision 3's structural claims about the existing test are PUBLISHED, not asserted** (`oc-glm-5-2` round 2: the only row touching the test file was V11, which greps `strings.Contains` only) | `grep -n 'scriptRows\|s8Rows\|strings.Count\|strings.Contains' host/verifygate/floor_raise_inventory_test.go` → `:32 scriptRows := []string{` · `:40 for site, needle := range scriptRows` · `:43 if count := strings.Count(block, needle); count != 1` · `:77 s8Rows := []string{` · `:85 for site, needle := range s8Rows` · `:88 if count := strings.Count(s8, needle); count != 1` · `:58` and `:93` the two `strings.Contains` sharedNeedle loops. Both hardcoded six-element lists and both per-row `count == 1` loops exist exactly as Decision 3 describes, and both use `t.Errorf` (not `t.Fatalf`), which is why a duplicated KNOWN row confounds rather than aborts the D1/D2 arms. Bound by **AC15** |
| V14 | **Reuse / conflict-surface audit** (`gpt5-6-sol` round 2, run verbatim) | proposed symbols repo-wide = **0** against a same-scope known-positive control of **9** `regexp.MustCompile` in `host/`; 4 regexp hits in `host/verifygate`, closest being the unexported `jobLine` YAML matcher; 1 markdown-pipe split repo-wide, in another package's unrelated test; `go.mod` has ONE direct dependency and zero parsing libraries. Nothing to extend — full table in **§Conflict Surface**, bound by **AC14** |

### Commands (runnable as printed)

Repo-root-relative, zsh, `export PATH=/opt/homebrew/bin:$PATH` first. Observed outputs are
restated as `#` comments so each block stays copy-paste safe.

```zsh
# V1 — pristine control (also the test command of V2/V3/V4 and ACs 1,4-10)
export AILANG_BIN=/tmp/ailang-v0300/ailang
go test ./host/verifygate/ -run '^TestFloorRaiseInventoryNamesEveryCoupledFile$' -count=1 -v > /tmp/v1 2>&1; rc=$?
grep -c '=== RUN' /tmp/v1; grep -c -- '--- PASS' /tmp/v1
# observed: rc=0, RUN=1, PASS=1
```

```zsh
# V2 — ARM A: land the script-home-only 7th row, assert landed + effect, then run the V1 block
cp scripts/verify_ail.sh /tmp/verify_ail_backup.sh
shasum -a 256 scripts/verify_ail.sh | cut -c1-8        # before: 5a1bbe89
perl -0pi -e 's/(#   6\. host\/verifygate\/module_manifest_gate_test\.go)/#   7. host\/store\/some_new_file.go   fabricated coupled site (repro)\n$1/' scripts/verify_ail.sh
shasum -a 256 scripts/verify_ail.sh | cut -c1-8        # after:  f97c4ec9  (LANDED)
sed -n '/── FLOOR-RAISE COUPLING INVENTORY/,/── END FLOOR-RAISE/p' scripts/verify_ail.sh | grep -cE '^#   [0-9]+\. '   # 7 (6->7)
bash -n scripts/verify_ail.sh; echo "rc=$?"            # rc=0
# then the V1 block -> observed: rc=0, `--- PASS`. THE GATE IS BLIND.
```

```zsh
# V3 — ARM B: land the S8-only 7th row, assert landed + effect, then run the V1 block
cp design_docs/coding-standards.md /tmp/standards_backup.md
shasum -a 256 design_docs/coding-standards.md | cut -c1-8   # before: b710a510
perl -0pi -e 's/(\| 6 \| `host\/verifygate\/module_manifest_gate_test\.go`)/| 7 | `host\/store\/some_new_file.go` | fabricated coupled site (repro) |\n$1/' design_docs/coding-standards.md
shasum -a 256 design_docs/coding-standards.md | cut -c1-8   # after:  63208d62  (LANDED)
sed -n '/## S8/,/^## /p' design_docs/coding-standards.md | grep -cE '^\| [0-9]+ \| '   # 7 (6->7)
# then the V1 block -> observed: rc=0, `--- PASS`. THE GATE IS BLIND IN THIS DIRECTION TOO.
```

```zsh
# V4 — restore both byte-identical, verify, re-run pristine control
cp /tmp/standards_backup.md design_docs/coding-standards.md
cp /tmp/verify_ail_backup.sh scripts/verify_ail.sh
shasum -a 256 design_docs/coding-standards.md | cut -c1-8   # b710a510
shasum -a 256 scripts/verify_ail.sh | cut -c1-8              # 5a1bbe89
git status --porcelain | wc -l                                # 0
# then the V1 block -> observed: rc=0, `ok ... host/verifygate`
```

```zsh
# V5 — verify_ail.sh baseline
export AILANG_BIN=/tmp/ailang-v0300/ailang
./scripts/verify_ail.sh > /tmp/v5 2>&1; echo "rc=$?"
tail -3 /tmp/v5
# observed: rc=0; `✓ verify gate PASSED: 11 required identities verified, 40 named tests pass`
```

```zsh
# V6 — verify_go.sh baseline (rc=0 is FORBIDDEN as a criterion; record the base FAIL set)
export AILANG_BIN=/tmp/ailang-v0300/ailang
./scripts/verify_go.sh > /tmp/v6 2>&1; echo "rc=$?"
grep -E '^--- FAIL' /tmp/v6 | sort
grep -E '^(ok|FAIL).*host/verifygate' /tmp/v6
# observed: rc=1; `--- FAIL: TestHandlerTimeoutKillsTheWholeProcessGroup (0.99s)`;
# `ok github.com/sunholo-data/ailang-world/host/verifygate`
```

```zsh
# V7 — hygiene baselines
go vet ./host/verifygate/; echo "vet rc=$?"
gofmt -l host/verifygate/; echo "gofmt rc=$? (empty=clean)"
# observed: vet rc=0; gofmt empty (0 bytes)
```

```zsh
# V8 — whole-file vs block count (the scope-control fact)
grep -cE '^#   [0-9]+\. ' scripts/verify_ail.sh
sed -n '/── FLOOR-RAISE COUPLING INVENTORY/,/── END FLOOR-RAISE/p' scripts/verify_ail.sh | grep -cE '^#   [0-9]+\. '
# observed: 6 and 6 — equal today, so a whole-file count is NOT a safe proxy for the block count
```

```zsh
# V9 — sharedNeedles multiplicity (justifies presence-only)
BLOCK=$(sed -n '/── FLOOR-RAISE COUPLING INVENTORY/,/── END FLOOR-RAISE/p' scripts/verify_ail.sh)
S8=$(sed -n '/## S8/,/^## /p' design_docs/coding-standards.md)
for n in 'packages/world-core/world/' 'REQUIRED_VERIFIED' 'EXACT_TOTAL_VERIFIED' 'world_package_ready_packet.golden.json' 'SELF_MOD_PUBLISH.md' 'module_manifest_gate_test.go' 'interfaceHash' 'does not move for'; do
  printf '%s  block=%s  s8=%s\n' "$n" "$(printf '%s' "$BLOCK" | grep -oF "$n" | wc -l | tr -d ' ')" "$(printf '%s' "$S8" | grep -oF "$n" | wc -l | tr -d ' ')"
done
# observed: REQUIRED_VERIFIED block=2 s8=1; EXACT_TOTAL_VERIFIED block=2 s8=2; all others =1 in both
```

```zsh
# V10 — §S8 is the last `## ` section (the test's slice runs to EOF)
awk 'NR>=90 && /^## /{print NR": "$0}' design_docs/coding-standards.md
# observed: 90: ## S8 — The floor-raise coupling inventory (added 2026-08-27, row 43)  (only hit)
```

```zsh
# V11 — sharedNeedles are presence-only (strings.Contains, never strings.Count)
grep -n 'strings.Contains(block\|strings.Contains(s8' host/verifygate/floor_raise_inventory_test.go
# observed: 58: if !strings.Contains(block, needle) {  and  93: if !strings.Contains(s8, needle) {
```

```zsh
# V12 — the two extraction regexes applied to the two real homes (objection 3, re-derived first-party)
# Applies siteReScript (`^#   [0-9]+\.\s+(\S+)`) to the verify_ail.sh block and siteReTable
# (`^\| [0-9]+ \| ([^|]+)`) to the §S8 slice (heading to EOF), prints both lists and the verdict.
python3 - <<'PY'
import re

def block(path, start, end):
    txt = open(path).read()
    i = txt.index(start)
    j = txt.index(end, i) + len(end)
    return txt[i:j]

script_home = block('scripts/verify_ail.sh', '── FLOOR-RAISE COUPLING INVENTORY', '── END FLOOR-RAISE')
s8_home = open('design_docs/coding-standards.md').read()
s8_home = s8_home[s8_home.index('## S8'):]  # slice to EOF (S8 is the last section)

siteReScript = re.compile(r'^#   [0-9]+\.\s+(\S+)')
siteReTable  = re.compile(r'^\| [0-9]+ \| ([^|]+)')

def extract(home, re_):
    out = []
    for line in home.split('\n'):
        m = re_.match(line)
        if m:
            out.append(m.group(1).strip('` \t'))
    return out

script = extract(script_home, siteReScript)
s8 = extract(s8_home, siteReTable)
print('script home (%d):' % len(script))
for s in script: print('  ', s)
print('S8 home (%d):' % len(s8))
for s in s8: print('  ', s)
print('sets equal:', set(script) == set(s8))
print('cardinality equal:', len(script) == len(s8))
PY
# observed: both homes print the identical 6-element set in the same order:
#   world/<module>.ail
#   packages/world-core/world/<module>.ail
#   scripts/verify_ail.sh
#   scripts/world_package_ready_packet.golden.json
#   docs/SELF_MOD_PUBLISH.md
#   host/verifygate/module_manifest_gate_test.go
# sets equal: True; cardinality equal: True (6 and 6)
```

## Out of Scope

- **No change to `design_docs/coding-standards.md` §S8's own text.** §S8 is ratification-class
  (human gate). This design READS §S8 and REQUIRES it to stay in sync with the script home (via
  the set-equality), but it does not propose editing §S8's text as part of this sprint. If a
  future row decides the two homes should be *merged* or that §S8 should carry an explicit
  cardinality statement, that is a parked human decision — not designed around here.
- **No change to `scripts/verify_ail.sh`** (the fixture is mutated only inside the drill,
  restored byte-identical).
- **No new package, no new script, no new CI job.** The fix is a test-only edit to
  `host/verifygate/floor_raise_inventory_test.go`. A design that required any of these would be
  out of scope for a ~0.1d row.
- **No `.ail` code** in scope (this row proposes no AILANG code).
- **No promotion of `sharedNeedles` to counts** (Decision 4; recorded as a declared residual).
- **No fix for the rig-local `verify_go.sh` base-red** (`TestHandlerTimeoutKillsTheWholeProcessGroup`).
  That is a queue item, not this sprint's scope; AC12 is written around it.
- If the sprint uncovers a genuinely separate defect, it is FILED as a "for the queue, not this
  sprint" note in the sprint log — not absorbed here.

## Related Documents

- `design_docs/coding-standards.md` §S8 — the second home of the inventory; the set-equality
  binds it to the script home. Ratification-class; not edited by this sprint.
- `scripts/verify_ail.sh` (lines 30-57) — the first home of the inventory; the bounded block.
- `host/verifygate/floor_raise_inventory_test.go` — the gate this row hardens; the only file
  modified.
- `design_docs/planned/w-shell-assignment-parser-drops-an-indented-assignment.md` — house style
  for the doc conventions (fenced-block commands, run-existence floors, mutation drill protocol,
  declared residuals, set-comparison gate criterion).
