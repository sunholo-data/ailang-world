# Sprint plan: w-marketing-checker-tolerance-pins (row 240, iteration 338)
Design authority: `design_docs/planned/w-marketing-checker-tolerance-pins.md` **r3**, especially §2–§4, §6 V13 and §7.

Status: **planned**, with the measured T11/T12 differences below requiring controller disposition before claiming the literal §4 mutation gate is satisfied. One milestone, M1. Estimate: ~0.2d, one executor run. No clause moves and no production tolerance changes.

Commits are the controller's job; the executor performs NO git write operations

Planner base: detached `572e177d741abe8bd6eddf0215284bbbdae7ca00`, worktree `/Users/voightkampff/dev/sunholo-data/.planner-wt-iter338`. The design's code base is `2c51c581`; read-only `git diff 2c51c581 -- scripts/marketingcheck/position.go scripts/marketingcheck/phone.go scripts/marketingcheck/main.go` was empty (rc 0). All §3 line addresses still apply.

| Milestone | Scope | Files | Commit |
|---|---|---|---|
| M1 | Copy the r3 reference byte-for-byte; append only the r3 comment to line 129 | `scripts/marketingcheck/tolerance_test.go`, `scripts/marketingcheck/position.go` | Controller only |

The r3 controller file supersedes §2.1's older r2 reference. Do not add arms to the banked generator tables, change any fixture or oracle, or edit any other production line. Do not edit `capture_test.go`, `phone.go`, `main.go`, `position_test.go`, `phone_test.go`, `go.mod`, `.github/`, `host/`, `.ail` files, or `tools/launchd/*`. Do not run `ailang messages`. Mutant edits to phone.go/main.go and PG's fixture occur only in scratch and are restored after each run.

## Measured differences and limits

The literal design §3 expected sets are preserved in the table below. They are authority expectations, not a claim that every r2 measurement has been repeated on r3.

- T11, scoped to `^TestMarketingPositionTolerancePins$`: rc 1, exactly parent `TestMarketingPositionTolerancePins` and leaves `gc-a3-conjunct-source-pin`, `a3-timeline-gap+3`. Design §3 says the arm is the sole killer.
- T12, with that same scope: rc 1, exactly the same parent and leaves `gc-a3-conjunct-source-pin`, `a3-timeline-gap+2`. Design §3's listed reds omit the source pin (and also require Corpus/Parity outside this measurement's scope).

The cause is the required r3 reference's exact string pin of the entire `okGap, gRes := gc >= 0 && posAbs(e) <= 2, e` clause. Changing 2 to 3 or 1 also removes that exact string. Preserve both required deliverable bytes and the design expectations; report these differences to the controller. Do not suppress a subtest or rewrite the reference to make a FAIL list match. The planner has not established the full-package T11/T12 lists. M1's literal mutation acceptance remains unresolved until the controller disposes of these measured differences.

T18's scoped measurement produced only parent `TestMarketingPositionTolerancePins` and leaf `TestMarketingPositionTolerancePins/gc-a3-conjunct-source-pin` (rc 1), consistent with V13. This is scoped confirmation; the full-package sole-killer claim remains the controller's outside-sandbox rerun. The unmutated three-function control passed: rc 0, exactly 19 `--- PASS`, 16 subtests + 3 parents, 0 FAIL and 0 SKIP. Scratch `gofmt -l scripts/marketingcheck` printed nothing (rc 0); scratch `go vet ./scripts/marketingcheck/` returned rc 0.

Every planner Go command used `perl -e 'alarm 600; exec @ARGV'` in a scratch package copy with the post-M1 files and copied seed baseline. No full package, full repository, verify_ail or capture test was run by the planner. No sandbox-uninformative result was encountered. A red mentioning bind/listen/operation not permitted, fake Chrome, or a 120 s capture deadline is **UNINFORMATIVE UNDER SANDBOX**, neither a pass nor a fail. Record the raw rc/output and let the controller rerun outside the sandbox; never use it to certify a mutation or gate.

## M1 environment and exact apply steps

Run the blocks in one Bash session from the executor worktree root. All build/test commands use the banked process-group deadline runner. rc 124 is a deadline: stop and report. Use 600 s for package tests, mutations, gofmt, vet and build; 2400 s for `go test ./...`; 1200 s for `verify_ail.sh`. Keep logs and scratch outside the tracked tree. The planner's temporary `.plan-scratch/` is removed before handoff.

```bash
set -eu
WT=$(pwd -P)
X=$(mktemp -d /tmp/world-iter338-executor.XXXXXX)
export WT X
mkdir -p "$X/logs" "$X/pristine"
cp design_docs/verification/world-iter336/run_bounded.sh.txt "$X/run_bounded.sh"
RB=$X/run_bounded.sh
export RB
export AILANG_BIN=$HOME/.pinned-ailang/ailang WORLD_TOOL_AILANG_BIN=$HOME/.pinned-ailang-tools/v0.54.15/ailang
export WORLD_PHONE_CORPUS=$HOME/.ailang/state/world-iter328-phone-corpus
# Use the existing writable Go caches (no cache regeneration under design_docs).
# Expected environment: go1.26.6 darwin/arm64; AILANG v0.41.0; tool AILANG v0.54.15.
shasum -a 256 -c <<'SHA'
59ae0dc3bf666cd92fb84e9cd971e2357cfb15cb47b076f2ae8a5609c6d4aab1  design_docs/planned/w-marketing-checker-tolerance-pins.md
c6e5fb0ba113886c0e1f05015425552cdecc122fb17c6911fc409c66e51b27ed  design_docs/verification/world-iter338/controller-r3-tolerance_test.go.txt
9029dbabd574761050399864cdd45f399ae9a3682ff1523e547479e75f8c669d  scripts/marketingcheck/position.go
9540d55a323640edab4933e0a4d8f8445be5046a5f8d6c1b0596035b1f0c0770  scripts/marketingcheck/phone.go
24448bb78265726f66085833b8b61ec5b28f16f9743cf6926a5c12e4cd189268  scripts/marketingcheck/main.go
733b40d11d7be70cdf5a0dfbade0a79e46ea7652fa9929390d733f2efc2201e0  design_docs/verification/world-iter338/plan-M1.diff.txt
SHA
test ! -e scripts/marketingcheck/tolerance_test.go
git apply --check design_docs/verification/world-iter338/plan-M1.diff.txt
# --check only: do not apply through git.
python3 - <<'PY'
import os
from pathlib import Path
p = Path('scripts/marketingcheck/position.go')
lines = p.read_bytes().splitlines(keepends=True)
assert lines[128] == b'\t\t\tokGap, gRes := gc >= 0 && posAbs(e) <= 2, e\n'
lines[128] = lines[128][:-1] + b' // gc >= 0 is banked as equivalent on R1 (premise pinned by TestMarketingPositionTolerancePins/gc-premise-nonfree-gaps-ge-2); editing this clause re-opens T18\n'
Path(os.environ['X'], 'position.post.go').write_bytes(b''.join(lines))
PY
cp design_docs/verification/world-iter338/controller-r3-tolerance_test.go.txt scripts/marketingcheck/tolerance_test.go
cp "$X/position.post.go" scripts/marketingcheck/position.go
cmp design_docs/verification/world-iter338/controller-r3-tolerance_test.go.txt scripts/marketingcheck/tolerance_test.go
shasum -a 256 -c <<'SHA'
c6e5fb0ba113886c0e1f05015425552cdecc122fb17c6911fc409c66e51b27ed  scripts/marketingcheck/tolerance_test.go
af0429fd70cdce9843c04bba18f800b7216f4a778234a7bbb2a158b4ef1fdb3a  scripts/marketingcheck/position.go
SHA
```

The resulting `position.go:129` is exactly the following line, with three leading tabs and the existing newline:

```go
			okGap, gRes := gc >= 0 && posAbs(e) <= 2, e // gc >= 0 is banked as equivalent on R1 (premise pinned by TestMarketingPositionTolerancePins/gc-premise-nonfree-gaps-ge-2); editing this clause re-opens T18
```

Do not gofmt-write the reference or production file. `gofmt -l` is the gate. The copied test has 14 position subtests, 2 phone entry subtests, and the routes top-level test. The gap +2/+3 arms retain the header-height compensation that prevents (a5) masking; phone bands 9–14 shift independently of the header; the routes exercise `run` itself.

## Full design §3 mutation table

The command column is runnable macOS BSD sed, one addressed replacement per mutant. T1–T8 are at :150, never the similar :90 clauses. T13/T14 match only `phoneR1Entry, 4)`, leaving the header tolerance intact. T17/T17b use BRE escapes for the literal brackets/dots. All ordinary rows run the full package; PG alone is scratch-copy-only and `-run`-limited. All rows must return rc 1. Design abbreviations below retain their §3 meaning: add `TestMarketing` to top-level names, resolve bare position arms under `TestMarketingPositionTolerancePins/`, bare `entry+4` under `TestMarketingPhoneEntryTolerancePins/`, `Corpus`/position `Parity` under `TestMarketingPosition`, and phone `Parity` under `TestMarketingPhone`. `phone-lines`/`position-lines` are under `TestMarketingRunbookMatchesHelp/`. FAIL parents are included alongside their failing children. A sole killer means one leaf plus its parent, not one raw FAIL line.

| ID | One-line mutation command | Must red, exactly design §3 | Design §3 measured red set / qualification |
|---|---|---|---|
| T1 | `sed -i '' '150s/posAbs(c.X-b.X) <= 2/posAbs(c.X-b.X) <= 200/' scripts/marketingcheck/position.go` | PositionTolerancePins/bg-x0+3 | sole killer |
| T2 | `sed -i '' '150s/posAbs(c.Y-b.Y) <= 2/posAbs(c.Y-b.Y) <= 200/' scripts/marketingcheck/position.go` | bg-y+3 | sole |
| T3 | `sed -i '' '150s/posAbs(c.X+c.W-b.X-b.W) <= 2/posAbs(c.X+c.W-b.X-b.W) <= 200/' scripts/marketingcheck/position.go` | bg-right-3 | sole |
| T4 | `sed -i '' '150s/posAbs(c.Y+c.H-b.Y-b.H) <= 2/posAbs(c.Y+c.H-b.Y-b.H) <= 200/' scripts/marketingcheck/position.go` | bg-bottom-3 | sole |
| T5 | `sed -i '' '150s/posAbs(c.X-b.X) <= 2/posAbs(c.X-b.X) <= 1/' scripts/marketingcheck/position.go` | bg-x0+2 | each sole |
| T6 | `sed -i '' '150s/posAbs(c.Y-b.Y) <= 2/posAbs(c.Y-b.Y) <= 1/' scripts/marketingcheck/position.go` | bg-y+2 | each sole |
| T7 | `sed -i '' '150s/posAbs(c.X+c.W-b.X-b.W) <= 2/posAbs(c.X+c.W-b.X-b.W) <= 1/' scripts/marketingcheck/position.go` | bg-right-2 | each sole |
| T8 | `sed -i '' '150s/posAbs(c.Y+c.H-b.Y-b.H) <= 2/posAbs(c.Y+c.H-b.Y-b.H) <= 1/' scripts/marketingcheck/position.go` | bg-bottom-2 | each sole |
| T9 | `sed -i '' '101s/posAbs(dh) > 2/posAbs(dh) > 3/' scripts/marketingcheck/position.go` | a4-graph-svg-h+3 | sole |
| T10 | `sed -i '' '101s/posAbs(dh) > 2/posAbs(dh) > 1/' scripts/marketingcheck/position.go` | a4-graph-svg-h+2 | also `Corpus/new/bnd-nonpara-drift-2/-3`, Parity |
| T11 | `sed -i '' '129s/posAbs(e) <= 2/posAbs(e) <= 3/' scripts/marketingcheck/position.go` | a3-timeline-gap+3 | sole |
| T12 | `sed -i '' '129s/posAbs(e) <= 2/posAbs(e) <= 1/' scripts/marketingcheck/position.go` | a3-timeline-gap+2 | also `Corpus/new/perm-roster-drift-2px`, Parity |
| T13 | `sed -i '' '168s/phoneR1Entry, 4)/phoneR1Entry, 40)/' scripts/marketingcheck/phone.go` | PhoneEntryTolerancePins/entry+5 | sole |
| T14 | `sed -i '' '168s/phoneR1Entry, 4)/phoneR1Entry, 3)/' scripts/marketingcheck/phone.go` | entry+4 | also `PhoneArms/v67/shift-4`, Parity |
| T15 | `sed -i '' '51s/case "position":/case "position-unrouted":/' scripts/marketingcheck/main.go` | RunRoutesPositionPhone | also `RunbookMatchesHelp/position-lines` |
| T16 | `sed -i '' '53s/case "phone":/case "phone-unrouted":/' scripts/marketingcheck/main.go` | RunRoutesPositionPhone | also `phone-lines` |
| T17 | `sed -i '' '52s/runPosition(args\[1:\], os\.Stdout, os\.Stderr)/runPhone(args[1:], os.Stdout, os.Stderr)/' scripts/marketingcheck/main.go` | RunRoutesPositionPhone | also `position-lines` |
| T17b | `sed -i '' '54s/runPhone(args\[1:\], os\.Stdout, os\.Stderr)/runPosition(args[1:], os.Stdout, os.Stderr)/' scripts/marketingcheck/main.go` | RunRoutesPositionPhone | also `phone-lines` |
| T18 | `sed -i '' '129s/gc >= 0 && posAbs(e) <= 2/posAbs(e) <= 2/' scripts/marketingcheck/position.go` | PositionTolerancePins/gc-a3-conjunct-source-pin | sole killer (V13). Behaviourally still equivalent on R1 (V09/V09b); premise guarded by `gc-premise-nonfree-gaps-ge-2` |
| PG | `sed -i '' '162s/"y": 423,/"y": 324,/' design_docs/verification/world-iter312/r1-rects.json` | PositionTolerancePins/gc-premise-nonfree-gaps-ge-2 (`non-free baseline gap 1 < 2`) | red (V12). `"y": 325,` (gb 2) stays green |
| C | `sed -i '' '150s/posAbs(c.Y-b.Y) <= 2/posAbs(c.Y-b.Y) <= -1/' scripts/marketingcheck/position.go` | 50 Corpus subtests, all 12 new position arms, and more | rc 1 |

T11/T12's extra r3 pin red is a measured deviation above, not silently added to these authority sets. For C, §3 says “and more”; V08 additionally names `TestMarketingPositionParity`, `TestMarketingPositionCLI/accept-rc`, `TestMarketingRunbookMatchesHelp/position-lines`, and `TestMarketingRunRoutesPositionPhone`. Do not invent a sole-killer requirement for C. T18 now reds **only** `TestMarketingPositionTolerancePins/gc-a3-conjunct-source-pin` (and its parent); it must not be treated as a surviving mutant. Its behavioural equivalence on R1 is a separate, baseline-specific classification.

Use a full scratch copy for the mutation round so docs, oracle and subprocess dependencies remain available. Only archive the committed base with read-only git commands, then copy in the two post-M1 files; exclude the launchd tree. Copy a pre-mutation backup before any sed operation, restore with `cp`, and verify with `cmp` even after a red. Never use a git restore/checkout/reset/apply to mutate or restore.

```bash
git archive HEAD > "$X/base.tar"
mkdir -p "$X/mut"
tar -xf "$X/base.tar" -C "$X/mut" --exclude='tools/launchd' --exclude='tools/launchd/*'
cp "$WT/scripts/marketingcheck/tolerance_test.go" "$X/mut/scripts/marketingcheck/tolerance_test.go"
cp "$WT/scripts/marketingcheck/position.go" "$X/mut/scripts/marketingcheck/position.go"
for f in scripts/marketingcheck/position.go scripts/marketingcheck/phone.go scripts/marketingcheck/main.go design_docs/verification/world-iter312/r1-rects.json; do
  mkdir -p "$X/pristine/$(dirname "$f")"
  cp "$X/mut/$f" "$X/pristine/$f"
done
cd "$X/mut"

# Run once before any mutation; require rc 0 outside the sandbox.
set +e
bash "$RB" 600 "$X/logs/mutation-control.log" go test -count=1 ./scripts/marketingcheck/
control_rc=$?
set -e
printf 'mutation control rc=%s\n' "$control_rc"
test "$control_rc" -eq 0

prepare_mutant() {
  ID=$1
  F=$2
  export ID F
  cmp "$X/pristine/$F" "$F"
  cp "$F" "$X/pre-mutation"
}
finish_mutant() {
  set +e
  if [ "$ID" = PG ]; then
    bash "$RB" 600 "$X/logs/$ID.log" go test -count=1 -v -run '^TestMarketingPositionTolerancePins$/^gc-premise-nonfree-gaps-ge-2$' ./scripts/marketingcheck/
  else
    bash "$RB" 600 "$X/logs/$ID.log" go test -count=1 -v ./scripts/marketingcheck/
  fi
  mutant_rc=$?
  set -e
  # Restore BEFORE interpreting the test status or matching the red set.
  cp "$X/pre-mutation" "$F"
  cmp "$X/pre-mutation" "$F"
  cmp "$X/pristine/$F" "$F"
  printf '%s rc=%s (restored byte-identical)\n' "$ID" "$mutant_rc"
  cat "$X/logs/$ID.log" "$X/logs/$ID.log.err"
  test "$mutant_rc" -eq 1
}

# For EACH table row, run prepare_mutant ID FILE, its literal sed command,
# then finish_mutant. Example T18 (all other rows follow the same steps):
prepare_mutant T18 scripts/marketingcheck/position.go
sed -i '' '129s/gc >= 0 && posAbs(e) <= 2/posAbs(e) <= 2/' scripts/marketingcheck/position.go
finish_mutant
# Read both logs and compare the complete FAIL list to the authority table.
# Do not declare acceptance from rc 1 alone. Run every row independently.
```

PG's sole permitted test command is the premise-only `-run` above. It must report `non-free baseline gap 1 < 2`; never run the full package on the edited baseline because its SHA pin would add unrelated failures. After PG's restore, also measure the boundary control in scratch:

```bash
prepare_mutant PG-boundary design_docs/verification/world-iter312/r1-rects.json
sed -i '' '162s/"y": 423,/"y": 325,/' design_docs/verification/world-iter312/r1-rects.json
set +e
bash "$RB" 600 "$X/logs/PG-boundary.log" go test -count=1 -v -run '^TestMarketingPositionTolerancePins$/^gc-premise-nonfree-gaps-ge-2$' ./scripts/marketingcheck/
boundary_rc=$?
set -e
cp "$X/pre-mutation" "$F"
cmp "$X/pre-mutation" "$F"
cmp "$X/pristine/$F" "$F"
test "$boundary_rc" -eq 0
cd "$WT"
```

Any premise failure, baseline or baseline-selection change, or change to the (a3) clause / gap-walk / `posFreeGap` rule requires re-deriving the non-free gaps and rerunning V09b before keeping T18's equivalence classification. R1 non-free gaps are 100, 100, 15, 15, 15; the exact bound is every non-free gap ≥2, with at least one non-free gap. A distinguishing production-baseline case requires a regression pin and a revised behavioural expectation for T18.

## Design §4 acceptance list (verbatim)

1. `go test -count=1 -v -run '^TestMarketing(PositionTolerancePins|PhoneEntryTolerancePins|RunRoutesPositionPhone)$' ./scripts/marketingcheck/`
   → rc 0; 19 `--- PASS` lines: 16 subtest lines (14 + 2), 3 top-level `--- PASS: TestMarketing…` lines.
   0 `--- SKIP` and 0 `--- FAIL` (V13; V07b lists the r2 18).
2. `go test -count=1 ./scripts/marketingcheck/` → rc 0.
3. Every §3 row reproduces as listed. T18 reds only `gc-a3-conjunct-source-pin`. PG reds only the premise subtest (scratch copy).
   The positive control C is rc 1.
4. `gofmt -l scripts/marketingcheck` prints nothing, and `go vet ./scripts/marketingcheck/` gives rc 0.
5. `go build ./... && go test ./...` → rc 0; `./scripts/verify_ail.sh` → rc 0.
6. `git diff --name-only 2c51c581 -- . ':!design_docs'` lists exactly
   `scripts/marketingcheck/position.go` and `scripts/marketingcheck/tolerance_test.go`.
   `git diff 2c51c581 -- scripts/marketingcheck/position.go` is exactly one changed line, `:129`, whose
   only difference is the §2.5 trailing comment (`--numstat` `1 1`; the line with the comment stripped
   equals the base line). `git diff --quiet 2c51c581 -- scripts/marketingcheck/phone.go scripts/marketingcheck/main.go`
   → rc 0.

**Flake fence.** A single re-run is allowed only when **all** of these hold:
- the FAIL list contains exactly `TestMarketingCaptureFakeChromeSixShots` and
  `TestMarketingCaptureFakeChromeSixShots/version-changed`, and nothing else;
- the subtest duration is `(120.0…s)`;
- its message carries `rc=-1`.

That is row 247's signature, seen once here (V10). Record both runs. Any other red, or this one
twice, stops M1. Do not edit `capture_test.go`.

This fence applies to the outside-sandbox acceptance gate. In this sandbox the specified infrastructure reds are UNINFORMATIVE UNDER SANDBOX; the fence cannot turn them into a pass/fail or authorize broader reruns. Record both runs if the exact outside-sandbox signature permits one rerun. A mutation's expected red list plus that flake is not the fence's “exactly” list. Any other red or recurrence stops M1. Never edit `capture_test.go`.

Run acceptance commands through the bounded runner as follows. The mutation round above implements item 3; it is the only required source of intended reds. Save both stdout and stderr and the rc for each gate; a nonzero unexpected rc stops M1 subject only to the exact fence. No stdout pipe may hide a command's rc.

```bash
cd "$WT"
gate() {
  gate_name=$1
  seconds=$2
  shift 2
  set +e
  bash "$RB" "$seconds" "$X/logs/$gate_name.log" "$@"
  gate_rc=$?
  set -e
  printf '%s rc=%s\n' "$gate_name" "$gate_rc"
  cat "$X/logs/$gate_name.log" "$X/logs/$gate_name.log.err"
  test "$gate_rc" -eq 0
}
gate tolerance-pins 600 go test -count=1 -v -run '^TestMarketing(PositionTolerancePins|PhoneEntryTolerancePins|RunRoutesPositionPhone)$' ./scripts/marketingcheck/
python3 - <<'PY'
import os
from pathlib import Path
text = Path(os.environ['X'], 'logs/tolerance-pins.log').read_text()
assert text.count('--- PASS:') == 19
assert sum(line.startswith('--- PASS: TestMarketing') for line in text.splitlines()) == 3
assert sum(line.lstrip().startswith('--- PASS: TestMarketingPositionTolerancePins/') for line in text.splitlines()) == 14
assert sum(line.lstrip().startswith('--- PASS: TestMarketingPhoneEntryTolerancePins/') for line in text.splitlines()) == 2
assert '--- SKIP' not in text and '--- FAIL' not in text
PY
gate marketingcheck 600 go test -count=1 ./scripts/marketingcheck/
# Run and evaluate ALL mutation rows, including PG and C, as above (item 3).
gate gofmt 600 gofmt -l scripts/marketingcheck
test ! -s "$X/logs/gofmt.log"
gate vet 600 go vet ./scripts/marketingcheck/
# Item 5's && is implemented by the gate's stop-on-nonzero behaviour.
gate build 600 go build ./...
gate all-tests 2400 go test ./...
gate verify-ail 1200 ./scripts/verify_ail.sh

# Item 6, literal read-only diff commands:
git diff --name-only 2c51c581 -- . ':!design_docs'
git diff 2c51c581 -- scripts/marketingcheck/position.go
git diff --numstat 2c51c581 -- scripts/marketingcheck/position.go
git diff --quiet 2c51c581 -- scripts/marketingcheck/phone.go scripts/marketingcheck/main.go
python3 - <<'PY'
from pathlib import Path
import subprocess
base = subprocess.check_output(['git', 'show', '2c51c581:scripts/marketingcheck/position.go'])
post = Path('scripts/marketingcheck/position.go').read_bytes()
b, p = base.splitlines(keepends=True), post.splitlines(keepends=True)
assert len(b) == len(p)
assert [i + 1 for i, (x, y) in enumerate(zip(b, p)) if x != y] == [129]
comment = b' // gc >= 0 is banked as equivalent on R1 (premise pinned by TestMarketingPositionTolerancePins/gc-premise-nonfree-gaps-ge-2); editing this clause re-opens T18'
assert p[128] == b[128][:-1] + comment + b'\n'
p[128] = p[128].replace(comment, b'', 1)
assert b''.join(p) == base
numstat = subprocess.check_output(['git', 'diff', '--numstat', '2c51c581', '--', 'scripts/marketingcheck/position.go']).decode()
assert numstat == '1\t1\tscripts/marketingcheck/position.go\n'
# git diff excludes an untracked new file. Validate executor scope using
# tracked differences plus the explicit untracked M1 code file instead.
tracked = subprocess.check_output(['git', 'diff', '--name-only', '2c51c581', '--', '.', ':!design_docs']).decode().splitlines()
untracked = subprocess.check_output(['git', 'ls-files', '--others', '--exclude-standard', '--', '.', ':!design_docs']).decode().splitlines()
assert set(tracked + untracked) == {'scripts/marketingcheck/position.go', 'scripts/marketingcheck/tolerance_test.go'}
PY
cmp design_docs/verification/world-iter338/controller-r3-tolerance_test.go.txt scripts/marketingcheck/tolerance_test.go
shasum -a 256 -c <<'SHA'
c6e5fb0ba113886c0e1f05015425552cdecc122fb17c6911fc409c66e51b27ed  scripts/marketingcheck/tolerance_test.go
af0429fd70cdce9843c04bba18f800b7216f4a778234a7bbb2a158b4ef1fdb3a  scripts/marketingcheck/position.go
SHA
```

Item 6's literal `git diff --name-only` can list the new test only after the controller stages or commits it. The executor leaves it untracked because git writes are forbidden. Its read-only scope assertion above includes that untracked file; the controller must rerun the literal item 6 after its own git operation. This is a git reporting constraint, not authorization for executor staging. The planner never applies M1 in the tracked worktree.

## Banked patch, hashes, and planner verification

Patch: `design_docs/verification/world-iter338/plan-M1.diff.txt`, SHA256 `733b40d11d7be70cdf5a0dfbade0a79e46ea7652fa9929390d733f2efc2201e0`. It is concatenated `diff -u` output from scratch base/post files, labelled `a/scripts/marketingcheck/position.go` and `b/scripts/marketingcheck/position.go`, plus `a/scripts/marketingcheck/tolerance_test.go` → `b/scripts/marketingcheck/tolerance_test.go` for the new file. Its empty scratch base is represented by the `-0,0` hunk range. There are no other hunks.

Planner ran `git apply --check design_docs/verification/world-iter338/plan-M1.diff.txt` against this worktree: **rc 0**. `jq -e . design_docs/planned/sprint_w-marketing-checker-tolerance-pins.json` returned rc 0; top-level and M1 key sets match the requested schema. All four Bash blocks passed `bash -n`, embedded Python compiled, the literal apply Python produced the exact scratch post image, and a scratch `patch -p1` reconstruction matched both post hashes. No git write operation was used. Scratch post images were hashed before generating the diff; the required test file equals the reference bytes. All 21 table sed commands were actually applied to pristine scratch files, verified to change exactly their addressed line, restored by `cp`, and checked by `cmp`.

| Post-M1 file | SHA256 |
|---|---|
| `scripts/marketingcheck/tolerance_test.go` | `c6e5fb0ba113886c0e1f05015425552cdecc122fb17c6911fc409c66e51b27ed` |
| `scripts/marketingcheck/position.go` | `af0429fd70cdce9843c04bba18f800b7216f4a778234a7bbb2a158b4ef1fdb3a` |

| Banked source | SHA256 |
|---|---|
| `design_docs/planned/w-marketing-checker-tolerance-pins.md` | `59ae0dc3bf666cd92fb84e9cd971e2357cfb15cb47b076f2ae8a5609c6d4aab1` |
| `design_docs/verification/world-iter338/controller-r3-tolerance_test.go.txt` | `c6e5fb0ba113886c0e1f05015425552cdecc122fb17c6911fc409c66e51b27ed` |

Planner measurement command sequence (all in scratch): unmutated three-function command from item 1 under `perl -e 'alarm 600; exec @ARGV'`; `gofmt -l` and `go vet` under that same bound; then T11/T12/T18 one at a time using the exact table sed commands and `perl -e 'alarm 600; exec @ARGV' go test -count=1 -v -run '^TestMarketingPositionTolerancePins$' ./scripts/marketingcheck/`. Each mutation was restored by `cp` from a pre-mutation backup and `cmp` returned rc 0. The results are recorded above and in JSON `measured_deviations_from_design`. Full acceptance and all full-package mutation claims remain controller verification work.

Deliverables are only this plan, `design_docs/planned/sprint_w-marketing-checker-tolerance-pins.json`, and the banked diff. The planner validates JSON with `jq -e .`, removes `.plan-scratch/`, verifies tracked files remain unchanged, and checks that final `git status --porcelain` lists exactly these three deliverables as untracked.
