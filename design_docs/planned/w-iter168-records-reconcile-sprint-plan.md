# Sprint plan — iteration 168: reconcile the five stranded mission records (162–166)

**Type**: DOCS-RECONCILIATION sprint. No Go, no `.ail`, no `tools/launchd/*`, no `scripts/*`.
The only files this sprint may create or modify are under `design_docs/`.

**Branch**: `mission/world-iter168-records` in worktree
`/Users/voightkampff/dev/sunholo-data/.wt-world-iter168-records`, based on `origin/dev` at
`9d72cfe82afcefe2adf54a5bc58dbbc37e606e87`.

**Authority**: `D-WORLD-35` = A (attended, 2026-09-07) — permit documentation-only
mission-record PRs to merge when the charter-declared ailang-code verify gate passes on the
exact head and remaining reds are demonstrated inherited against the base; *"Reconcile all
pending ledger records, including D-WORLD-34, without overwriting attended rulings."*
The blocker is a rebase/reconcile, not a ruling.

**Outcome**: one branch which, merged to `dev`, carries the iteration 162–166 records in the
CURRENT document structure, with iteration 167's landed state and every attended ledger ruling
intact. PRs #127 and #128 are then SUPERSEDED and closed unmerged.

---

## 0. Environment — do this first, every shell

```
export PATH=/opt/homebrew/bin:$PATH
export AILANG_BIN=$HOME/.pinned-ailang/ailang
cd /Users/voightkampff/dev/sunholo-data/.wt-world-iter168-records
```

`AILANG_BIN` is MANDATORY. Without it `go test ./...` is rc=1 **by design** (the pinned-binary
guard), and you will misread a passing tree as broken.

`MISSION_DOC` happens to be set to `design_docs/world-mission.md` in the mission shell, which
makes bare `bash scripts/mission_decisions.sh --check` appear to work. **Never rely on it.**
Always pass `--file design_docs/world-mission.md` explicitly — the script's built-in default is
`design_docs/v1-mission.md`, which does not exist in this repo (measured: `ls` rc=1).

### Rig hazards — measured on this machine, all of them have bitten this mission

| Hazard | Consequence | What to do instead |
|---|---|---|
| `cmd \| grep X \| head` then `echo $?` | `$?` is `head`'s, not `grep`'s — an empty result reads rc=0 | `if grep -q ... ; then ... ; fi` or capture rc immediately after the grep |
| `grep -c` with zero matches | prints `0` **and exits 1** | never chain a bare `grep -c` under `set -e`; read the printed number |
| `set -e` | RESETS `$?` | capture `rc=$?` on the line immediately after the command |
| `local out=$(cmd)` | swallows the exit code | assign, then `rc=$?` on a separate line |
| unquoted `--include=*.go` in zsh | glob-expanded/mangled | quote it: `--include='*.go'` |
| `${PIPESTATUS[0]}` in zsh | wrong index (zsh arrays are 1-based, and the var is `pipestatus`) | avoid pipelines whose rc matters |
| inline multi-line `jq` filter | mangled by zsh | `jq -r -f <file>` |
| `git checkout -- <file>` | restores to **HEAD**, not to your uncommitted baseline | mutation-drill restores use `cp` from a `/tmp` snapshot, then verify with `shasum -a 256` |
| `cat -A` | not supported on this macOS `cat` (`illegal option`) | use `od -c` or `sed -n 'Np' \| xxd` |

---

## 1. BASE-STATE MEASUREMENTS (taken by the planner on `9d72cfe`, before any work)

Every acceptance criterion below is anchored to one of these. **Re-run each of these before you
start** and confirm you reproduce it; if a base value differs, STOP and report — the branch has
moved and the criteria's non-vacuity claims no longer hold.

| # | Command | Base value |
|---|---|---|
| B1 | `git rev-parse HEAD` | `9d72cfe82afcefe2adf54a5bc58dbbc37e606e87` |
| B2 | `grep -cE '^## ' design_docs/world-mission-log.md` | `21` |
| B3 | `grep -E '^## ' design_docs/world-mission-log.md \| grep -vcE '^## [0-9]+ — [0-9]{4}-[0-9]{2}-[0-9]{2} — .+$'` | `0` (prints 0, rc=1) |
| B4 | `grep -oE '^## [0-9]+' design_docs/world-mission-log.md \| awk '{print $2}' \| tr '\n' ' '` | `141 142 143 145 146 147 148 149 150 151 152 153 154 155 156 157 158 159 160 161 167` |
| B5 | `for n in 162 163 164 165 166; do grep -c "^## $n — " design_docs/world-mission-log.md; done` | `0 0 0 0 0` |
| B6 | `grep -c '^| [0-9]' design_docs/world-mission-index.md` | `160` |
| B7 | `python3 /tmp/i168_index_check.py` (instrument I2, §2) | rc=**1**, `headings=161 index_rows=160 errors=1`, sole error `log heading 167 has NO index row` |
| B8 | `bash scripts/mission_decisions.sh --check --file design_docs/world-mission.md` | rc=0, `decision ledger valid: 22 rows` |
| B9 | `bash scripts/mission_decisions.sh --open --file design_docs/world-mission.md \| wc -l` | `0` |
| B10 | `bash scripts/mission_decisions.sh --all --file design_docs/world-mission.md \| wc -l` | `22` |
| B11 | `grep -c '^## STATUS 2026' design_docs/world-mission.md` | `3` (iterations 167, 161, 160) |
| B12 | `grep -c '^## STATUS' design_docs/world-mission-status-archive.md` | `172`; top-of-file stamp is `(iteration 159)` at line 1 |
| B13 | `shasum -a 256 design_docs/world-mission-dashboard.md` | `4817cdbb61bffef6385e9f9c3ba3d27145be073a3f1becef53831b01bf298dec` |
| B14 | `shasum -a 256 design_docs/planned/w-de-fork-ci-ownership.md` | `5903c4e40c66c7b7550750273477fb51f646eb9e24e1901fc34cfac131590366` |
| B15 | `shasum -a 256 design_docs/world-mission.md` | `12b9947e784494c05a8084df3340ea0a363a83319e91b3855d402f6af2c62350` |
| B16 | `shasum -a 256 design_docs/world-mission-log.md` | `42f1cd58396ee09c077b42db9c3ebfce4daeffd9135b9ab86f6707fde2292538` |
| B17 | `shasum -a 256 design_docs/world-mission-index.md` | `5d8967ff9a04a5a29716cbd973a056f0072fe19042209be4701b13fae5fea81e` |
| B18 | `shasum -a 256 design_docs/world-mission-status-archive.md` | `2c9801e682cdc9bb4456b80b58879124bca813e4de036831c2aff48dfb3eb269` |
| B19 | `for f in design_docs/planned/sprint_w-defork-pin-redirects-work-repo.json design_docs/planned/w-defork-pin-redirects-work-repo{,-fleet-issue,-sprint-plan}.md; do test -e "$f"; echo "$f $?"; done` | all `1` (ABSENT) |
| B20 | `go vet ./...` | rc=**0**, no output |
| B21 | `./scripts/verify_ail.sh` | rc=**0**; `✓ verify gate PASSED: 11 required identities verified, 40 named tests pass` |
| B22 | `go test ./... -count=1` | rc=**0**, `grep -c '^ok ' = 19`, `grep -Ec '^(FAIL\|--- FAIL)' = 0` |
| B23 | `./scripts/verify_go.sh` | rc=**1** — see §1.1, this is the base red |
| B24 | `grep -rq 'NEGCTL-168-Q7X4M2' .` | rc=**1** (absent everywhere; this is the sprint's negative-control literal) |
| B25 | `grep -rq 'decision-ledger:start' .` | rc=**0** (positive control for the same instrument) |
| B26 | `git rev-parse origin/mission/world-iter163-record origin/mission/world-iter166-defork-pin` | `00d3406cf6cbd9fa99fe795b0ec08b98223f27f4` / `a529dc99cfea880393f0a3087f70d612d6c74d5b` |

### 1.1 The base `verify_go.sh` red is NOT driver drift — re-measured

The controller directive predicted a FLEET-owned driver-drift red. **That is stale.** Iteration
167 (`37a78ab`, PR #130) retired the `launchd-drivers` job and the `--driver-fleet-check`
diagnostic. `.github/workflows/ci.yml` now has exactly **2** jobs — `ailang-verify` and
`go-verify` (measured: `grep -nE '^  [a-z0-9-]+:' .github/workflows/*.yml`). Memory saying
"3 jobs" is stale.

The measured base red is a single test under the race detector:

```
── go test ./... -count=1 -race -timeout 8m
--- FAIL: TestCLIRealSubprocessEpisode (5.66s)
    cli_test.go:196: daemon announcement timed out; stderr=
FAIL	github.com/sunholo-data/ailang-world/cmd/ailang-worldd	6.413s
```

Every other stage of `verify_go.sh` passed, including the **World mission-input gate (decision
ledger)** which printed `decision ledger valid: 22 rows`. Isolated re-run is GREEN:
`go test ./cmd/ailang-worldd/ -run TestCLIRealSubprocessEpisode -race -count=2` → rc=**0**.
So it is a load/timing flake (the known `w-daemon-timeout-test-flake` class), not a
deterministic failure.

**Executor rule**: `verify_go.sh` is allowed to stay rc=1 at head ONLY if the failure set is a
subset of `{TestCLIRealSubprocessEpisode}`. Any other failing test, or any failure in a stage
*before* `── go test ./... -count=1 -race`, is a REGRESSION caused by this sprint and blocks
the milestone — because this sprint touches only `design_docs/**` and cannot legitimately move
a Go result. Record the full failing-test list at head and diff it against `{TestCLIRealSubprocessEpisode}`.

### 1.2 Two pre-existing defects found while measuring — RECORD, DO NOT FIX

Both are out of scope. Do **not** edit them; report them in the iteration-168 log entry and to
the controller.

1. **The charter's STATUS rotation rule text has been stranded in the archive.**
   `design_docs/world-mission.md` line 795 is a bare heading `## STATUS (rotation rule)` with an
   **empty body**. Its two-line body now sits at `design_docs/world-mission-status-archive.md`
   lines 18–19: *"Newest **3** STATUS stamps live here; older ones move to
   `world-mission-status-archive.md`. At Gate 4, after adding your stamp, move the now-4th stamp
   to the TOP of the archive file."* Some earlier rotation moved a stamp and dragged the rule
   with it. (This sprint follows that rule as written; it does not relocate the text.)
2. **`world-mission-index.md` has no row for iterations 64 and 144**, and neither
   `world-mission-log.md` nor `world-mission-log-archive.md` has a `## 64` / `## 144` heading.
   Those two records appear to have been lost before the rotation. Instrument I2 is therefore
   scoped to "every heading has a row and every row has a heading", not "no numeric gaps".

### 1.3 Structural facts the reconciliation depends on (measured, not assumed)

- **All three `world-mission-log.md` variants share a byte-identical 2449-line prefix.**
  `head -2449 <file> | shasum -a 256` = `906a3cd2865b3eeaa0c21b34861b9e9c1b8bef357617efbed5caff0a81ba4886`
  for `dev`, for `origin/mission/world-iter163-record`, and for
  `origin/mission/world-iter166-defork-pin`. This is what makes the splice below mechanical.
- **The index row format is derived, not free text.** For every one of the 160 base rows,
  `| <N> | <date> | <text> |` satisfies: `<N>`/`<date>` equal the log/archive heading
  `## <N> — <date> — <title>`, and `<text>` = `<title>` if `len(title) <= 150`, else
  `title[:147] + "..."`. Measured: 160/160 match, max untruncated length 149, all 110 truncated
  rows exactly 150 chars. Rows are **strictly descending** by `<N>`.
- **The charter holds exactly 3 STATUS stamps**; the displaced 4th is **prepended to the TOP**
  of `world-mission-status-archive.md`. Confirmed by `git show 9d72cfe --
  design_docs/world-mission-status-archive.md`, which added the `(iteration 159)` stamp at line 1.
- **Stamp text is stable across branches.** The `(iteration 160)` and `(iteration 161)` stamps in
  `dev`'s charter are **byte-identical** to the copies `origin/mission/world-iter163-record`
  put in its archive (656 and 675 bytes respectively). Use `dev`'s copies.
- **PR #127's charter diff contains nothing but** the `| D-WORLD-34 | OPEN |` row and three
  STATUS stamps. Verified: filtering those two patterns out of the diff leaves zero lines.
- **PR #128 has no STATUS stamp for iteration 166 anywhere** — not in its charter, not in the
  status archive (it does not touch that file). See M4.3.

---

## 2. INSTRUMENTS — build these before M1 and never edit them again

Create all three in `/tmp` (they are tools, not deliverables; they must NOT be committed).
After creating each, run the stated base check and confirm the stated base rc. **An instrument
that has not been shown to fail is not an instrument.**

### I1 — canonical-heading + order checker

```
cat > /tmp/i168_headings.sh <<'SH'
#!/bin/bash
# usage: i168_headings.sh <log.md>
f="$1"
tot=$(grep -cE '^## ' "$f"); rc_tot=$?
noncanon=$(grep -E '^## ' "$f" | grep -vcE '^## [0-9]+ — [0-9]{4}-[0-9]{2}-[0-9]{2} — .+$')
order=$(grep -oE '^## (Iteration )?[0-9]+' "$f" | grep -oE '[0-9]+$' | awk 'NR>1 && $1<=p{bad++}{p=$1}END{print bad+0}')
nolead=$(awk 'prev!="" && /^## (Iteration )?[0-9]+ —/{n++}{prev=$0}END{print n+0}' "$f")
echo "headings=$tot noncanonical=$noncanon out_of_order=$order missing_blank_before=$nolead"
[ "$noncanon" = 0 ] && [ "$order" = 0 ] && [ "$nolead" = 0 ]
SH
chmod +x /tmp/i168_headings.sh
```

**Prove it can fail (do this once, now).** The `#127` blob is a live negative control:

```
git show origin/mission/world-iter163-record:design_docs/world-mission-log.md > /tmp/i168_127.md
/tmp/i168_headings.sh /tmp/i168_127.md ; echo "rc=$?"
```
Planner-measured expected: `headings=24 noncanonical=4 out_of_order=1 missing_blank_before=4`, **rc=1**.

The `(Iteration )?` in the order/blank extractions is LOAD-BEARING. The planner's first
draft matched only `^## [0-9]+`, which silently skipped the four non-canonical headings and
reported `out_of_order=0` on `#127` — a green from an instrument that never looked. If you
edit I1, re-run this negative control and confirm all three numbers move.

Positive control on the base file:
```
/tmp/i168_headings.sh design_docs/world-mission-log.md ; echo "rc=$?"
```
Planner-measured expected: `headings=21 noncanonical=0 out_of_order=0 missing_blank_before=0`, **rc=0**.
Second positive control (a file that is canonical but is NOT the base file): the same script on
`/tmp/i168_128.md` → `headings=21 noncanonical=0 out_of_order=0 missing_blank_before=0`, rc=0.

### I2 — index↔log consistency checker

```
cat > /tmp/i168_index_check.py <<'PY'
import re, sys
base='design_docs/'
def heads(fn):
    d={}
    for l in open(base+fn):
        m=re.match(r'^## (\d+) — (\d{4}-\d\d-\d\d) — (.*?)\s*$', l)
        if m: d.setdefault(int(m.group(1)), (m.group(2), m.group(3), fn))
    return d
H={}
for fn in ('world-mission-log-archive.md','world-mission-log.md'):
    for k,v in heads(fn).items(): H[k]=v
rows=[]
for l in open(base+'world-mission-index.md'):
    m=re.match(r'^\| (\d+) \| (\S+) \| (.*) \|$', l.rstrip('\n'))
    if m: rows.append((int(m.group(1)), m.group(2), m.group(3)))
errs=[]; seen=set()
for n,d,t in rows:
    if n in seen: errs.append(f"duplicate index row {n}")
    seen.add(n)
    if n not in H: errs.append(f"index row {n} has no log heading"); continue
    hd,ht,_=H[n]
    exp = ht if len(ht)<=150 else ht[:147]+'...'
    if d!=hd: errs.append(f"row {n} date {d} != heading {hd}")
    if t!=exp: errs.append(f"row {n} text mismatch\n  idx={t!r}\n  exp={exp!r}")
for n in sorted(H):
    if n not in seen: errs.append(f"log heading {n} has NO index row")
nums=[r[0] for r in rows]
for i in range(len(nums)-1):
    if nums[i] <= nums[i+1]: errs.append(f"index not strictly descending at {nums[i]} then {nums[i+1]}")
print(f"headings={len(H)} index_rows={len(rows)} errors={len(errs)}")
for e in errs[:20]: print("  ERR:", e)
sys.exit(1 if errs else 0)
PY
python3 /tmp/i168_index_check.py ; echo "rc=$?"
```
Planner-measured base expectation (**B7**): `headings=161 index_rows=160 errors=1`, sole error
`log heading 167 has NO index row`, **rc=1**. This is the proof that I2's success criterion is
NOT vacuous: it is failing right now, and M2 is what turns it green.

### I3 — attended-ruling survival checker

```
cat > /tmp/i168_rulings.sh <<'SH'
#!/bin/bash
# The three attended rulings that MUST survive byte-identically, and the ledger fences.
c="design_docs/world-mission.md"
fail=0
check() { got=$(grep "^| $1 |" "$c" | shasum -a 256 | cut -d' ' -f1)
          if [ "$got" != "$2" ]; then echo "RULING ALTERED: $1 got=$got want=$2"; fail=1
          else echo "ok $1"; fi; }
check D-WORLD-33 c535d51b4f29c03439b9a775944bcb4bb34157f748a3048ea9138dab4ca1c97c
check D-WORLD-34 f43e1bb84a9844258480dda73e6e972b1325cd45b90a82b86c0137f5a1b95c20
check D-WORLD-35 79748252854d09d4692f81c407d890e88e1bfa9fe036ecbd5e60ae62fb30c081
s=$(grep -c '<!-- decision-ledger:start -->' "$c"); e=$(grep -c '<!-- decision-ledger:end -->' "$c")
echo "ledger_markers start=$s end=$e"
[ "$s" = 1 ] && [ "$e" = 1 ] || { echo "LEDGER MARKER DAMAGE"; fail=1; }
n=$(grep -cE '^\| D-WORLD-' "$c"); echo "ledger_rows=$n"
[ "$n" = 22 ] || { echo "LEDGER ROW COUNT CHANGED"; fail=1; }
o=$(bash scripts/mission_decisions.sh --open --file "$c" | wc -l | tr -d ' '); echo "open_rows=$o"
[ "$o" = 0 ] || { echo "A RESOLVED ROW WAS REOPENED"; fail=1; }
# iteration-167 STATUS must still be in the charter, verbatim on its headline
grep -q '^## STATUS 2026-09-07 (iteration 167) — \*\*DE-FORK CI OWNERSHIP REPAIRED;' "$c" \
  || { echo "ITERATION-167 STATUS HEADLINE LOST"; fail=1; }
exit $fail
SH
chmod +x /tmp/i168_rulings.sh
bash /tmp/i168_rulings.sh ; echo "rc=$?"
```
Planner-measured base expectation: `ok D-WORLD-33 / ok D-WORLD-34 / ok D-WORLD-35`,
`ledger_markers start=1 end=1`, `ledger_rows=22`, `open_rows=0`, **rc=0**.

### Snapshots for mutation-drill restores

`git checkout -- <file>` restores to HEAD, which is NOT your working baseline mid-sprint. Before
each drill, snapshot with `cp`, and after restoring, prove it with `shasum -a 256`:

```
mkdir -p /tmp/i168_snap
cp design_docs/world-mission.md            /tmp/i168_snap/
cp design_docs/world-mission-log.md        /tmp/i168_snap/
cp design_docs/world-mission-index.md      /tmp/i168_snap/
cp design_docs/world-mission-status-archive.md /tmp/i168_snap/
cp design_docs/world-mission-dashboard.md  /tmp/i168_snap/
```
(Re-snapshot after every milestone commit.)

---

## M0 — Preflight: reproduce the base state and stage the source blobs

**Goal**: prove the branch is where the plan thinks it is, and materialise the two stranded
branches' log files as local read-only sources.

**Files touched**: none (all output to `/tmp`).

**Steps**

1. Verify B1, B2, B3, B4, B5, B6, B8, B9, B10, B11, B12, B13, B14, B19, B24, B25, B26. Any
   mismatch → STOP and report.
2. Materialise the sources and pin them by hash:
```
git show origin/mission/world-iter163-record:design_docs/world-mission-log.md > /tmp/i168_127.md
git show origin/mission/world-iter166-defork-pin:design_docs/world-mission-log.md > /tmp/i168_128.md
shasum -a 256 /tmp/i168_127.md /tmp/i168_128.md
```
   **Required**:
   - `/tmp/i168_127.md` → `1b475319a048560d708d723ebc0382c370fa8d86258570533fd7c74ea463881e`
   - `/tmp/i168_128.md` → `62c8bb59cea86d76308cbb0113a2c11e7d5388b0f6b1a72416e49677f1f908e8`

   If either hash differs, the line numbers in M1 are invalid — STOP and report. Do not
   "find the headings yourself".
3. Confirm the shared prefix:
```
for f in design_docs/world-mission-log.md /tmp/i168_127.md /tmp/i168_128.md; do head -2449 "$f" | shasum -a 256; done
```
   All three must print `906a3cd2865b3eeaa0c21b34861b9e9c1b8bef357617efbed5caff0a81ba4886`.
4. Build I1, I2, I3 (§2) and run **all** of their stated base checks, including the
   negative-control run of I1 against `/tmp/i168_127.md` (must be rc=1 with `noncanonical=4`).
5. Take the `/tmp/i168_snap` snapshots.

**Acceptance criteria**

- **AC0.1** All 16 re-measured base values in step 1 match §1 exactly. *(Base = the criterion;
  it is a gate, not a change.)*
- **AC0.2** Both source hashes in step 2 match. Positive control: they are non-empty
  (`wc -l /tmp/i168_127.md` = `2693`, `/tmp/i168_128.md` = `2520`).
- **AC0.3** I1 rc=1 on `/tmp/i168_127.md` (with `noncanonical=4 out_of_order=1 missing_blank_before=4`) **and** rc=0 on the base log. Both directions must be
  observed and pasted into the record. An instrument observed only green is not admitted.
- **AC0.4** I2 rc=**1** at base with exactly one error, `log heading 167 has NO index row`.
- **AC0.5** I3 rc=**0** at base.

**Mutation drill (M0)** — prove I3 can fail:
```
cp design_docs/world-mission.md /tmp/i168_snap/pre_m0_drill.md
python3 - <<'PY'
p='design_docs/world-mission.md'; s=open(p).read()
i=s.index('| D-WORLD-35 |'); j=s.index('\n',i)+1
open(p,'w').write(s[:j]+s[i:j]+s[j:])
PY
bash /tmp/i168_rulings.sh ; echo "KILLER_rc=$?"
bash scripts/mission_decisions.sh --check --file design_docs/world-mission.md ; echo "LEDGER_rc=$?"
cp /tmp/i168_snap/pre_m0_drill.md design_docs/world-mission.md
shasum -a 256 design_docs/world-mission.md   # must equal B15
```
**Named killer**: `bash /tmp/i168_rulings.sh` → expected **rc=1** with `LEDGER ROW COUNT CHANGED`.
Secondary killer: `scripts/mission_decisions.sh --check` → expected **rc=1** with
`duplicate decision ID: D-WORLD-35` (planner reproduced this exact output). Restore must return
`shasum` to `12b9947e…2350`.

**Commit**: none (M0 produces no tracked change).

---

## M1 — Splice the five log entries into `world-mission-log.md`, canonical and in numeric order

**Goal**: `design_docs/world-mission-log.md` carries entries 162, 163, 164, 165, 166 in ascending
numeric order, positioned before `## 167`, each with a canonical heading
`## <N> — <date> — <title>`, with the `dev` prefix and the `## 167` entry untouched.

**Files touched**: `design_docs/world-mission-log.md` (only).

**Exact recipe** (planner ran this end-to-end; the resulting file was measured):

```
D=design_docs/world-mission-log.md; A=/tmp/i168_127.md; B=/tmp/i168_128.md
{ head -2449 "$D";
  sed -n '2451,2511p' "$A";   # entry 162
  sed -n '2513,2572p' "$A";   # entry 163
  sed -n '2634,2693p' "$A"; echo;   # entry 164 (last line of file, needs a trailing blank)
  sed -n '2574,2632p' "$A";   # entry 165
  sed -n '2450,2520p' "$B"; echo;   # entry 166 (last line of file, needs a trailing blank)
  sed -n '2450,2505p' "$D"; } > /tmp/i168_newlog.md
sed -i '' -E 's/^## Iteration ([0-9]+) — /## \1 — /' /tmp/i168_newlog.md
shasum -a 256 /tmp/i168_newlog.md
cp /tmp/i168_newlog.md "$D"
```

Notes the executor must not deviate from:
- The `sed -n` ranges **deliberately skip** the `---` separator lines at `/tmp/i168_127.md`
  lines 2450, 2512, 2573, 2633. `dev`'s log separates entries with a blank line only; `#127`
  added `---` rules that would be a new, inconsistent convention.
- `#127` stores entry **165 physically before 164**. The recipe reorders them. Do not
  "preserve source order".
- The `sed -i ''` form is BSD sed (macOS). The empty `''` argument is required.
- Only the four `## Iteration N` headings are rewritten; entry 166 is already canonical.

**Acceptance criteria**

- **AC1.1** `shasum -a 256 design_docs/world-mission-log.md` = `018cc34a76b19748cce17ae920b2a39e3d5e92bb24b0271555f416e5fd259c0b`.
  *(Base B16 = `42f1cd58…2538`. Changed → non-vacuous.)*
- **AC1.2** `/tmp/i168_headings.sh design_docs/world-mission-log.md` → rc=**0** and prints
  `headings=26 noncanonical=0 out_of_order=0 missing_blank_before=0`.
  *(Base B2 = `headings=21`. Changed.)* Negative control already banked in AC0.3.
- **AC1.3** Each of the five entries is present exactly once with the exact canonical heading:
```
for h in \
 '## 162 — 2026-09-06 — DE-FORK CI repair parks at the authority gate [HARNESS]' \
 '## 163 — 2026-09-06 — DE-FORK repair remains parked after every Agent role fails closed [HARNESS]' \
 '## 164 — 2026-09-06 — DE-FORK repair parks again on D-WORLD-34 with fallback independent judge [HARNESS]' \
 '## 165 — 2026-09-07 — DE-FORK repair remains parked; four required Agent roles confirm the authority gate [HARNESS]' \
 '## 166 — 2026-09-07 — the loaded gun in row 68 fired: the driver pin ran this mission inside the WRONG REPOSITORY, and every health instrument read green [HARNESS]' ; do
   printf '%s -> ' "${h:0:12}"; grep -Fxc "$h" design_docs/world-mission-log.md; done
```
  Expect `1` five times. *(Base B5 = 0 five times.)*
  **Positive control**: the same instrument on
  `## 167 — 2026-09-07 — the DE-FORK red was a suspended gate, not a failing check; three roles each found a different way the repair could have been vacuous [HARNESS]`
  → `1` (true at base and after; proves `grep -Fxc` is running).
  **Negative control**: `grep -Fxc '## 999 — 2026-09-07 — NEGCTL-168-Q7X4M2' design_docs/world-mission-log.md`
  → prints `0`, rc=1. *(B24 established this literal is absent repo-wide before the sprint; after
  the sprint it appears only in this plan file, never in `world-mission-log.md`, so the
  file-scoped grep stays a valid negative control.)*
- **AC1.4** No `## Iteration ` heading survives:
  `grep -c '^## Iteration ' design_docs/world-mission-log.md` → `0` (rc=1).
  Positive control that the instrument works: the same grep on `/tmp/i168_127.md` → `4` (rc=0).
- **AC1.5** The `dev` prefix and the `## 167` entry are untouched:
  `head -2449 design_docs/world-mission-log.md | shasum -a 256` = `906a3cd2…4886`, and
  `sed -n '/^## 167 — /,$p' design_docs/world-mission-log.md | shasum -a 256` equals the same
  command run against `/tmp/i168_snap/world-mission-log.md`.
- **AC1.6** `grep -c '^---$' design_docs/world-mission-log.md` → `8`. *(Base = 8; `/tmp/i168_127.md`
  = 12. This criterion is TRUE AT BASE and is therefore a **regression guard**, not a progress
  criterion — it is marked as such deliberately, and its job is to catch the `---` separators
  leaking in. It is proven capable of failing by the M1 drill below.)*

**Mutation drill (M1)** — prove AC1.2 and AC1.6 can fail:
```
cp design_docs/world-mission-log.md /tmp/i168_snap/pre_m1_drill.md
# arm A: de-canonicalise one heading (the exact regression this milestone exists to prevent)
sed -i '' 's/^## 164 — 2026-09-06 — /## Iteration 164 — 2026-09-06 — /' design_docs/world-mission-log.md
/tmp/i168_headings.sh design_docs/world-mission-log.md ; echo "KILLER_A_rc=$?"
python3 /tmp/i168_index_check.py ; echo "KILLER_A2_rc=$?"
cp /tmp/i168_snap/pre_m1_drill.md design_docs/world-mission-log.md
shasum -a 256 design_docs/world-mission-log.md      # must equal AC1.1's hash
# arm B: reinsert a `---` separator
python3 - <<'PY'
p='design_docs/world-mission-log.md'; s=open(p).read()
s=s.replace('\n## 165 — ', '\n---\n## 165 — ',1); open(p,'w').write(s)
PY
c=$(grep -c '^---$' design_docs/world-mission-log.md); echo "KILLER_B_dashes=$c (want 9, AC1.6 fails)"
cp /tmp/i168_snap/pre_m1_drill.md design_docs/world-mission-log.md
shasum -a 256 design_docs/world-mission-log.md      # must equal AC1.1's hash
```
**Named killers**: arm A → `/tmp/i168_headings.sh` **rc=1** with
`headings=26 noncanonical=1 out_of_order=0 missing_blank_before=0` (planner-measured); and, run
after M2, `python3 /tmp/i168_index_check.py` **rc=1** with `headings=165 index_rows=166 errors=1`
and the error `index row 164 has no log heading`. That second killer is the concrete proof of
guard (b): an entry whose heading is off-shape becomes invisible to the index and to any future
rotation, and the index gate says so.
Arm B → `grep -c '^---$'` prints **9**, rc=0, so AC1.6's `= 8` fails.

**Commit**:
`docs(mission): splice iterations 162-166 into the rotated log, renormalized and in order`

---

## M2 — Add index rows 162–167 to `world-mission-index.md`

**Goal**: `world-mission-index.md` covers every log heading, in the file's own derived format,
strictly descending.

**Files touched**: `design_docs/world-mission-index.md` (only).

**Why 167 is included**: it is not scope creep — I2 fails at base (**B7**) with exactly one
error, `log heading 167 has NO index row`. Iteration 167 landed its log entry without an index
row. `D-WORLD-35` = A directs reconciliation of all pending records. Six rows, not five.

**Exact recipe**: insert the block below **immediately after** the single line `|---|---|---|`
(assert it occurs exactly once first), preserving the order shown:

```
| 167 | 2026-09-07 | the DE-FORK red was a suspended gate, not a failing check; three roles each found a different way the repair could have been vacuous [HARNESS] |
| 166 | 2026-09-07 | the loaded gun in row 68 fired: the driver pin ran this mission inside the WRONG REPOSITORY, and every health instrument read green [HARNESS] |
| 165 | 2026-09-07 | DE-FORK repair remains parked; four required Agent roles confirm the authority gate [HARNESS] |
| 164 | 2026-09-06 | DE-FORK repair parks again on D-WORLD-34 with fallback independent judge [HARNESS] |
| 163 | 2026-09-06 | DE-FORK repair remains parked after every Agent role fails closed [HARNESS] |
| 162 | 2026-09-06 | DE-FORK CI repair parks at the authority gate [HARNESS] |
```

These six lines were **generated** from the log headings by the file's own truncation rule
(`title` if `len <= 150`, else `title[:147] + "..."`); none of the six exceeds 150 chars, so none
is truncated. Do not retype them — paste verbatim. Do not touch any other row: `#127`'s own
index rows for 162–165 happen to match, but its file does not have 166 or 167.

**Acceptance criteria**

- **AC2.1** `python3 /tmp/i168_index_check.py` → rc=**0**, prints `headings=166 index_rows=166 errors=0`.
  *(Base B7 = rc=1, `headings=161 index_rows=160 errors=1`. Changed → non-vacuous. Planner
  reproduced rc=0 on a dry-run assembly.)*
- **AC2.2** `grep -c '^| [0-9]' design_docs/world-mission-index.md` → `166`. *(Base B6 = 160.)*
- **AC2.3** `shasum -a 256 design_docs/world-mission-index.md` differs from B17 (`5d8967ff…a81e`).
- **AC2.4** Each of the six rows appears exactly once: run `grep -Fxc "<row>"` for each → `1`.
  **Positive control**: `grep -Fxc '| 161 | 2026-09-06 | row65 parks after its proposed guard fails its own Markdown example [HARNESS] |'` → `1` (true at base and after).
  **Negative control**: `grep -Fxc '| 168 | 2026-09-07 | NEGCTL-168-Q7X4M2 |'` → `0`, rc=1.
- **AC2.5** `grep -c '|---|---|---|' design_docs/world-mission-index.md` → `1` (the table header
  was not duplicated).

**Mutation drill (M2)** — prove AC2.1 can fail in the two ways that matter:
```
cp design_docs/world-mission-index.md /tmp/i168_snap/pre_m2_drill.md
# arm A: wrong date on one row (the silent class: row present, content wrong)
sed -i '' 's/^| 164 | 2026-09-06 |/| 164 | 2026-09-07 |/' design_docs/world-mission-index.md
python3 /tmp/i168_index_check.py ; echo "KILLER_A_rc=$?"
cp /tmp/i168_snap/pre_m2_drill.md design_docs/world-mission-index.md ; shasum -a 256 design_docs/world-mission-index.md
# arm B: rows out of descending order
python3 - <<'PY'
p='design_docs/world-mission-index.md'; L=open(p).read().split('\n')
i=[k for k,l in enumerate(L) if l.startswith('| 166 |')][0]
L[i],L[i+1]=L[i+1],L[i]; open(p,'w').write('\n'.join(L))
PY
python3 /tmp/i168_index_check.py ; echo "KILLER_B_rc=$?"
cp /tmp/i168_snap/pre_m2_drill.md design_docs/world-mission-index.md ; shasum -a 256 design_docs/world-mission-index.md
```
**Named killer**: `python3 /tmp/i168_index_check.py` → **rc=1** both arms (planner-measured):
arm A prints `errors=1` with `row 164 date 2026-09-07 != heading 2026-09-06`; arm B prints
`errors=1` with `index not strictly descending at 165 then 166`.
Both restores must return the `shasum` to AC2.3's value.

**Commit**: `docs(mission): index rows for iterations 162-167 (167 was missing too)`

---

## M3 — Land PR #128's four new design docs; leave `w-de-fork-ci-ownership.md` at `dev`'s version

**Goal**: the four documents `#128` authored are on the branch byte-identically; the one
add/add-conflicting document keeps `dev`'s (iteration 167's) newer revision.

**Files touched** (all created, none modified):
- `design_docs/planned/sprint_w-defork-pin-redirects-work-repo.json`
- `design_docs/planned/w-defork-pin-redirects-work-repo-fleet-issue.md`
- `design_docs/planned/w-defork-pin-redirects-work-repo-sprint-plan.md`
- `design_docs/planned/w-defork-pin-redirects-work-repo.md`

**Exact recipe**
```
for f in design_docs/planned/sprint_w-defork-pin-redirects-work-repo.json \
         design_docs/planned/w-defork-pin-redirects-work-repo-fleet-issue.md \
         design_docs/planned/w-defork-pin-redirects-work-repo-sprint-plan.md \
         design_docs/planned/w-defork-pin-redirects-work-repo.md ; do
  git show origin/mission/world-iter166-defork-pin:"$f" > "$f"
done
```
**Do not** run `git checkout origin/mission/world-iter166-defork-pin -- design_docs/planned/`;
that would also drag `w-de-fork-ci-ownership.md` (which `#128` does not carry, but `#127` does)
and is the exact mistake this milestone exists to avoid.

**Acceptance criteria**

- **AC3.1** All four exist with these hashes (`shasum -a 256`):
  - `sprint_w-defork-pin-redirects-work-repo.json` → `353e4f1db7d3eba9ca5bafd28e2cb48b8d751db869e3e3ba1e7c40f3058d87fc`
  - `w-defork-pin-redirects-work-repo-fleet-issue.md` → `06a6099276c125fa653f08b5e1fb4b750c4e9a69705bcc7da204962d04bcab5a`
  - `w-defork-pin-redirects-work-repo-sprint-plan.md` → `4d4547aacb94587c17ce2aef2a0cacb0ea658db5ad7114b8efe6146aca1e0778`
  - `w-defork-pin-redirects-work-repo.md` → `a3bb7aed3081dd7c19277ab143cfa265b8f5b48cc405ad8b605396c2209d5474`
  *(Base B19 = all four ABSENT. Changed → non-vacuous.)*
- **AC3.2** `shasum -a 256 design_docs/planned/w-de-fork-ci-ownership.md` = B14
  (`5903c4e4…0366`), i.e. `dev`'s version, and **not** `#127`'s
  (`23fbc9df9815a494c2a97e3b72789d29d8603f9c01f9ed79291a83c892182250`).
  Both hashes are stated so the check is a discrimination, not a tautology.
- **AC3.3** `python3 -c "import json;json.load(open('design_docs/planned/sprint_w-defork-pin-redirects-work-repo.json'))"`
  → rc=0 (the JSON parses; a truncated `git show` would fail here).
- **AC3.4** Each of the four is a NEW untracked file, and no tracked planned doc was modified:
  `for f in <the four paths>; do git status --porcelain "$f"; done` prints `?? <path>` four
  times; and `git status --porcelain design_docs/planned/ | grep -c '^ M'` → `0` (rc=1).
  Note: `grep -c '^??'` over the whole directory is **5**, not 4 — this sprint plan file is also
  untracked there. Do not write a criterion that forgets itself.

**Mutation drill (M3)**:
```
cp design_docs/planned/w-de-fork-ci-ownership.md /tmp/i168_snap/
git show origin/mission/world-iter163-record:design_docs/planned/w-de-fork-ci-ownership.md > design_docs/planned/w-de-fork-ci-ownership.md
shasum -a 256 design_docs/planned/w-de-fork-ci-ownership.md ; echo "expected the 23fbc9df... regression hash"
cp /tmp/i168_snap/w-de-fork-ci-ownership.md design_docs/planned/w-de-fork-ci-ownership.md
shasum -a 256 design_docs/planned/w-de-fork-ci-ownership.md   # must be 5903c4e4...0366
```
**Named killer**: the AC3.2 `shasum` comparison — it must print `23fbc9df…2250` under the mutant
(AC3.2 FAILS) and `5903c4e4…0366` after restore (AC3.2 passes).

**Commit**: `docs(mission): land PR #128's four defork-pin design docs, keep dev's w-de-fork-ci-ownership`

---

## M4 — Charter: refresh queue row 68, rotate STATUS to 167/166/165, touch the ledger not at all

**Files touched**: `design_docs/world-mission.md` (only).

### M4.1 — Ledger: DO NOTHING. Prove it.

`dev`'s ledger already carries **`D-WORLD-34` RESOLVED** and **`D-WORLD-35` RESOLVED** (22 rows,
zero open). PR `#127` would add `| D-WORLD-34 | OPEN |` and PR `#128` would add
`| D-WORLD-35 | OPEN |`. Applying either would create a **duplicate ID** and reopen an
attended ruling — the exact silent regression `D-WORLD-35`'s own text forbids. **Do not apply
either hunk.** The ledger block between `<!-- decision-ledger:start -->` and
`<!-- decision-ledger:end -->` must be byte-identical to base.

### M4.2 — Queue row 68: take `#128`'s refreshed text

Replace `dev`'s single line beginning `68. **w-driver-pin-named-world-points-at-the-wrong-repo**`
with `#128`'s. `#128` is the record of the iteration in which the row's predicted condition
actually fired, so its text is the newer measurement.

```
python3 - <<'PY'
import subprocess
new = subprocess.run(['git','show','origin/mission/world-iter166-defork-pin:design_docs/world-mission.md'],
                     capture_output=True, text=True).stdout
new = [l for l in new.split('\n') if l.startswith('68. **w-driver-pin-named-world-points-at-the-wrong-repo**')]
assert len(new)==1, new
p='design_docs/world-mission.md'; L=open(p).read().split('\n')
idx=[i for i,l in enumerate(L) if l.startswith('68. **w-driver-pin-named-world-points-at-the-wrong-repo**')]
assert len(idx)==1, idx
L[idx[0]]=new[0]; open(p,'w').write('\n'.join(L))
PY
```

### M4.3 — STATUS rotation

Target end state, top to bottom under `## STATUS (rotation rule)`:

1. `## STATUS 2026-09-07 (iteration 167) — …` — **unchanged, byte-identical to base**.
2. `## STATUS 2026-09-07 (iteration 166) — …` — **RECONSTRUCTED**, see below.
3. `## STATUS 2026-09-07 (iteration 165) — …` — byte-identical to
   `git show origin/mission/world-iter163-record:design_docs/world-mission.md`'s copy (1430 bytes).

Stamps for iterations **164, 163, 162, 161, 160** leave the charter (M5 prepends them to the
archive). Take 164 and 163 verbatim from `#127`'s charter (1352 and 1231 bytes); take 162
verbatim from `#127`'s status archive (1223 bytes); take 161 and 160 verbatim from the charter
you already have (they are byte-identical to `#127`'s archive copies — planner verified).

**Do not "correct" any historical stamp.** Stamps 162–165 say things like *"Ledger 21 rows /
3 OPEN"* and *"D-WORLD-32/33 remain OPEN"*. Those were true when written. They are records, not
current-state claims. Rewriting them would be exactly the overwrite `D-WORLD-35` forbids. The
discrepancy is recorded in the iteration-168 log entry instead.

**The iteration-166 stamp does not exist.** PR `#128` wrote a log entry for 166 and no charter
stamp (planner verified: no `(iteration 166)` STATUS line in `#128`'s charter, and `#128` does
not touch `world-mission-status-archive.md`). Rather than leave a hole or invent evidence, insert
this stamp **verbatim** — it is a summary of the landed log entry 166 and is explicitly labelled
as a reconstruction, on one line:

```
## STATUS 2026-09-07 (iteration 166) — **THE LOADED GUN IN ROW 68 FIRED: THE DRIVER PIN RAN THIS MISSION INSIDE THE WRONG REPOSITORY, AND EVERY HEALTH INSTRUMENT READ GREEN [HARNESS]**: **RECONSTRUCTED AT ITERATION 168** — PR #128 landed a log entry for iteration 166 but never wrote a charter stamp; every claim below is a quotation of the landed log entry 166, not a measurement taken at iteration 168. Gate 1 read GREEN while measuring the WRONG REPOSITORY: `git fetch` + `rev-parse` agreed, `mission-base.sh record gate1` banked `878939117`, and the running-skill-vs-`origin/dev` `cmp` was byte-identical — every one of those readings was taken inside `~/.ailang-driver-pin/world`, a worktree of `sunholo-data/ailang` whose `git-common-dir` is `~/dev/sunholo-data/ailang/.git` and which has no `design_docs/world-mission.md`. Mechanism: the DE-FORK (`e92594c`) deleted World's driver copy, the regenerated plist runs the FLEET driver, and `tools/launchd/lib/pin-root.sh` sets `MISSION_WORKDIR` to `pwd`, so the driver's own comment claiming `$REPO` survives the re-exec is FALSE. Pick was not the queue head: row 68 `w-driver-pin-named-world-points-at-the-wrong-repo`, surfaced iteration 150 as a prediction, had fired. Design `planned/w-defork-pin-redirects-work-repo.md` (588 lines); quorum r1 BLOCKED (gemini + glm reject, astra ABSENT on a $0.1224 > $0.10 budget refusal, $0.04441749); quorum r2 BLOCKED on a real, previously-unseen defect — the helper was called inside a command substitution, so `_pin_stale` ran in a subshell and could never reach the driver; narrow-refinement carve-out applied with gemini's verbatim fix. Planner (codex) repaired the design rather than restating it: AC3 could not red on its own named mutation, AC4's `grep -c … >= 1` passes on log HISTORY, and row 68 already existed. M1 mitigation LANDED OUT OF SANDBOX: `AILANG_DRIVER_PIN=0` appended to `~/.config/ailang/mission-world.env` behind a pre-edit backup, idempotent, AC-C1..AC-C4 measured, mutation drill run and restored byte-identical by sha256. The durable M2 guard is FLEET-owned and was handed off as a prepared issue; no edit to `tools/launchd/*` in either repository. Independent evaluator `sonnet` PASS 94/100, zero blocking, and found three defects of ours first-hand: AC-F8 is VACUOUS (the specified implementation and its own named mutant produce byte-identical output), AC-E3 has an expiring precondition, and **`D-WORLD-34` did not exist in the committed ledger** — it lived only on unmerged PR #127, which is what `D-WORLD-35` was filed to answer. AC-C5/AC-C6 are fire-dependent and DEFERRED. Metered $0.183783 of $5. Seven-clause 1.0 bar: goal unmoved; no product code landed, by design. Next: the next fire must run AC-C5 and AC-C6 first.
```

**Acceptance criteria**

- **AC4.1** `bash /tmp/i168_rulings.sh` → rc=**0**, all three `ok D-WORLD-3x`, `ledger_rows=22`,
  `open_rows=0`, `ledger_markers start=1 end=1`, iteration-167 headline present.
  *(True at base — a **regression guard**, and the M0 drill already proved it can fail.)*
- **AC4.2** The ledger block is byte-identical to base:
```
sed -n '/<!-- decision-ledger:start -->/,/<!-- decision-ledger:end -->/p' design_docs/world-mission.md | shasum -a 256
sed -n '/<!-- decision-ledger:start -->/,/<!-- decision-ledger:end -->/p' /tmp/i168_snap/world-mission.md | shasum -a 256
```
  The two hashes must be equal. Positive control: both commands print a non-empty 64-hex hash and
  `sed -n '…/p' … | wc -l` is `26`.
- **AC4.3** `grep -c '| D-WORLD-34 | OPEN |' design_docs/world-mission.md` → `0` (rc=1) and
  `grep -c '| D-WORLD-35 | OPEN |' design_docs/world-mission.md` → `0` (rc=1).
  Positive control for the same instrument:
  `grep -c '| D-WORLD-35 | RESOLVED |' design_docs/world-mission.md` → `1`.
- **AC4.4** `bash scripts/mission_decisions.sh --check --file design_docs/world-mission.md` →
  rc=0, `decision ledger valid: 22 rows`. *(= B8, regression guard; M0's drill proved it fails on
  a duplicate ID.)*
- **AC4.5** `grep -c '^## STATUS 2026' design_docs/world-mission.md` → `3`, and
  `grep -oE '^## STATUS [0-9-]+ \(iteration [0-9]+\)' design_docs/world-mission.md | grep -oE '[0-9]+\)' | tr -d ')' | tr '\n' ' '`
  → `167 166 165`. *(Base B11 = 3 stamps, `167 161 160`. The count is unchanged by design; the
  **identities** change, which is the non-vacuous part — assert the identities, not the count.)*
- **AC4.6** `grep -c 'RECONSTRUCTED AT ITERATION 168' design_docs/world-mission.md` → `1`.
  *(Base = 0. This is the guard that a reconstructed stamp can never be mistaken for an attended
  contemporaneous one.)*
- **AC4.7** Row 68 was refreshed:
  `grep -c '^68\. \*\*w-driver-pin-named-world-points-at-the-wrong-repo\*\* · clause-2 · \*\*THE LATENT CONDITION FIRED\.\*\*' design_docs/world-mission.md` → `1`;
  `grep -c 'A PIN WORKTREE NAMED FOR THIS MISSION HOLDS THE \*OTHER\* MISSION' design_docs/world-mission.md` → `0` (rc=1);
  and `grep '^68\. \*\*w-driver-pin' design_docs/world-mission.md | shasum -a 256 | cut -d' ' -f1`
  = `7379b864d857900e7e11d29cd876d41b72a87a7a70b203f9490ec7fcb3357c99`.
  *(Base: the first grep = 0, the second = 1, so both directions change.)*
- **AC4.8** `grep -c '^68\. ' design_docs/world-mission.md` → `1` (the row was replaced, not duplicated).
- **AC4.9** Stamps 165/164/163/162 are byte-identical to their `#127` sources. For each N in
  165 (charter) and 164, 163 (charter), 162 (`#127` archive), compare
  `grep "^## STATUS .*(iteration N)" <target> | shasum -a 256` against the same grep over the
  `#127` blob. All four must match.

**Mutation drill (M4)** — the guard-(a) drill, run it explicitly:
```
cp design_docs/world-mission.md /tmp/i168_snap/pre_m4_drill.md
# arm A: reopen an attended ruling (what applying #127's hunk would effectively do)
sed -i '' 's/^| D-WORLD-35 | RESOLVED |/| D-WORLD-35 | OPEN |/' design_docs/world-mission.md
bash /tmp/i168_rulings.sh ; echo "KILLER_A_rc=$?"
bash scripts/mission_decisions.sh --open --file design_docs/world-mission.md | wc -l
cp /tmp/i168_snap/pre_m4_drill.md design_docs/world-mission.md ; shasum -a 256 design_docs/world-mission.md
# arm B: overwrite iteration 167's STATUS with a stranded branch's older stamp
python3 - <<'PY'
p='design_docs/world-mission.md'; L=open(p).read().split('\n')
i=[k for k,l in enumerate(L) if l.startswith('## STATUS 2026-09-07 (iteration 167)')][0]
L[i]='## STATUS 2026-09-07 (iteration 165) — **STALE OVERWRITE**: x'
open(p,'w').write('\n'.join(L))
PY
bash /tmp/i168_rulings.sh ; echo "KILLER_B_rc=$?"
cp /tmp/i168_snap/pre_m4_drill.md design_docs/world-mission.md ; shasum -a 256 design_docs/world-mission.md
```
**Named killers**: arm A → `bash /tmp/i168_rulings.sh` **rc=1** with `RULING ALTERED: D-WORLD-35`
and `A RESOLVED ROW WAS REOPENED`, and `--open | wc -l` → `1`. Arm B → `bash /tmp/i168_rulings.sh`
**rc=1** with `ITERATION-167 STATUS HEADLINE LOST`. Both restores verified by `shasum`.

**Commit**: `docs(mission): charter — refresh queue row 68, rotate STATUS to 167/166/165, ledger untouched`

---

## M5 — Prepend the five displaced STATUS stamps to `world-mission-status-archive.md`

**Goal**: stamps 164, 163, 162, 161, 160 land at the TOP of the archive, in that order, above the
existing `(iteration 159)` stamp — which is what the rotation rule (quoted in §1.2) prescribes
when it is applied five times in a row.

**Files touched**: `design_docs/world-mission-status-archive.md` (only).

**End state**, archive lines 1–11 (odd lines are stamps, even lines blank):

```
 1  ## STATUS 2026-09-06 (iteration 164) — ...
 3  ## STATUS 2026-09-06 (iteration 163) — ...
 5  ## STATUS 2026-09-06 (iteration 162) — ...
 7  ## STATUS 2026-09-06 (iteration 161) — ...
 9  ## STATUS 2026-09-06 (iteration 160) — ...
11  ## STATUS 2026-09-06 (iteration 159) — ...   <- previously line 1
```

Each stamp is followed by exactly one blank line, matching the existing file's shape.

**Acceptance criteria**

- **AC5.1** `grep -c '^## STATUS' design_docs/world-mission-status-archive.md` → `177`.
  *(Base B12 = 172. Changed.)*
- **AC5.2** `grep -oE '^## STATUS [0-9-]+ \(iteration [0-9]+\)' design_docs/world-mission-status-archive.md | head -6 | grep -oE '[0-9]+\)' | tr -d ')' | tr '\n' ' '`
  → `164 163 162 161 160 159`. *(Base = `159` alone at the head.)*
- **AC5.3** No stamp is in BOTH files. For each N in 160..164:
  `grep -c "(iteration $N)" design_docs/world-mission.md` → `0` (rc=1) and
  `grep -c "(iteration $N)" design_docs/world-mission-status-archive.md` → `1`.
  Positive control: for N=167, charter → `1`, archive → `0`.
- **AC5.4** For each N in 160..164, the archived stamp is byte-identical to its source:
  `grep "(iteration $N)" design_docs/world-mission-status-archive.md | shasum -a 256` equals the
  same grep against `/tmp/i168_snap/world-mission.md` (for 160, 161) or against the `#127` blobs
  (for 162, 163, 164).
- **AC5.5** The tail of the archive is untouched:
  `tail -n +12 design_docs/world-mission-status-archive.md | shasum -a 256` equals
  `tail -n +2 /tmp/i168_snap/world-mission-status-archive.md | shasum -a 256`.
- **AC5.6** No stamp appears twice: for each N in 160..164,
  `grep -c "(iteration $N)" design_docs/world-mission-status-archive.md` → exactly `1`.

**Mutation drill (M5)**:
```
cp design_docs/world-mission-status-archive.md /tmp/i168_snap/pre_m5_drill.md
# arm A: duplicate a stamp (the "moved but not removed" failure)
python3 - <<'PY'
p='design_docs/world-mission-status-archive.md'; L=open(p).read().split('\n')
L.insert(1, ''); L.insert(1, L[0]); open(p,'w').write('\n'.join(L))
PY
for N in 160 161 162 163 164; do printf "%s:" $N; grep -c "(iteration $N)" design_docs/world-mission-status-archive.md; done
cp /tmp/i168_snap/pre_m5_drill.md design_docs/world-mission-status-archive.md ; shasum -a 256 design_docs/world-mission-status-archive.md
# arm B: truncate the tail (proves AC5.5 is load-bearing)
head -100 design_docs/world-mission-status-archive.md > /tmp/i168_trunc && cp /tmp/i168_trunc design_docs/world-mission-status-archive.md
tail -n +12 design_docs/world-mission-status-archive.md | shasum -a 256 ; echo "must differ from the AC5.5 reference"
cp /tmp/i168_snap/pre_m5_drill.md design_docs/world-mission-status-archive.md ; shasum -a 256 design_docs/world-mission-status-archive.md
```
**Named killers**: arm A → the AC5.6 loop prints `164:2` (rc=0) instead of `164:1`; arm B → the
AC5.5 `shasum` differs from the reference. Both restores verified by `shasum` against the
pre-drill snapshot.

**Commit**: `docs(mission): rotate STATUS stamps 160-164 into the archive, newest at top`

---

## M6 — Dashboard: forward-update three regions, regress nothing

**Goal**: `world-mission-dashboard.md` reflects that the reconciliation happened, while every
line of iteration 167's snapshot survives verbatim. It must equal **neither** `dev`'s version
**nor** either stranded branch's older version.

**Files touched**: `design_docs/world-mission-dashboard.md` (only).

**Exactly three edits. Nothing else.**

1. Line 3 — replace
   `Snapshot: 2026-09-07, iteration 167. History: world-mission-log.md.`
   with
   `Snapshot: 2026-09-07, iteration 168 (records reconciliation). Iteration 167's product findings below are unchanged. History: world-mission-log.md.`
2. Replace the whole `## Blocked (work, not a decision)` section body (its three bullet lines)
   with:
```
- **Nothing.** Iteration 168 reconciled the five stranded records (162–166) from PRs #127/#128
  onto the rotated log (`37a10f5`) and canonical headings (`9166de0`) under `D-WORLD-35` = A.
  `world-mission-index.md` now covers 162–167 — the 167 row was missing too — and the charter
  STATUS series is gapless through 167. PRs #127/#128 are SUPERSEDED; close them unmerged.
```
3. In `## Next`, replace
   `- Rebase + land #127/#128, then rows 66, 68–78, 81, 82, 83, then 39.`
   with
   `- Rows 66, 68–78, 81, 82, 83, then 39. (#127/#128 superseded by iteration 168's reconciliation.)`

**Acceptance criteria**

- **AC6.1** `shasum -a 256 design_docs/world-mission-dashboard.md` differs from B13
  (`4817cdbb…8dec`) — the file changed.
- **AC6.2** It also differs from both stranded versions:
  `git show origin/mission/world-iter163-record:design_docs/world-mission-dashboard.md | shasum -a 256`
  and the `#128` equivalent. Three distinct hashes. *(This is the anti-regression discrimination:
  a hash equal to either branch's means the older snapshot won.)*
- **AC6.3** `diff /tmp/i168_snap/world-mission-dashboard.md design_docs/world-mission-dashboard.md | grep -c '^[<>]'` → `11`
  — and, split out, `grep -c '^<'` → `5` and `grep -c '^>'` → `6`. (Edit 1: 1 out, 1 in. Edit 2:
  3 out, 4 in. Edit 3: 1 out, 1 in.) Planner measured this on a dry run of the exact three edits.
  Any other number means an unplanned edit.
- **AC6.4** Every iteration-167 fact survives verbatim — each `grep -Fc` → `1`:
  - `` `verify_go.sh` exited 127 before `go build` ``
  - `Landed: PR #130 — retire the `launchd-drivers` job`
  - `` Authority: `D-WORLD-34` = A (attended, 2026-09-07) ``
  - `Independent evaluation PASS 95/100, zero blocking`
  - `Decision ledger: **22 rows, ZERO OPEN**`
  **Negative control**: `grep -Fc 'NEGCTL-168-Q7X4M2' design_docs/world-mission-dashboard.md` → `0` (rc=1).
- **AC6.5** `grep -Fc 'PRs #127 and #128 are CONFLICTING/DIRTY' design_docs/world-mission-dashboard.md` → `0` (rc=1).
  *(Base = 1. Changed.)*

**Mutation drill (M6)**:
```
cp design_docs/world-mission-dashboard.md /tmp/i168_snap/pre_m6_drill.md
git show origin/mission/world-iter166-defork-pin:design_docs/world-mission-dashboard.md > design_docs/world-mission-dashboard.md
shasum -a 256 design_docs/world-mission-dashboard.md
grep -Fc 'Independent evaluation PASS 95/100, zero blocking' design_docs/world-mission-dashboard.md ; echo "KILLER_rc=$?"
cp /tmp/i168_snap/pre_m6_drill.md design_docs/world-mission-dashboard.md ; shasum -a 256 design_docs/world-mission-dashboard.md
```
**Named killer**: with `#128`'s older dashboard in place, AC6.4's
`grep -Fc 'Independent evaluation PASS 95/100, zero blocking'` prints **`0` with rc=1** — i.e.
regressing to a stranded snapshot is caught. Restore verified by `shasum` against the pre-drill copy.

**Commit**: `docs(mission): dashboard — records reconciled, #127/#128 superseded`

---

## M7 — Verify gate, scope proof, and the iteration-168 log entry

**Goal**: prove the branch is green to the charter-declared standard, prove the sprint touched
nothing outside `design_docs/`, and record the iteration.

**Files touched**: `design_docs/world-mission-log.md`, `design_docs/world-mission-index.md`,
`design_docs/world-mission.md`, `design_docs/world-mission-status-archive.md` (the iteration-168
record itself, added at the end of M7 following the same rules M1/M2/M4/M5 established).

**Gate — run in this order, and read `go vet` BEFORE any test verdict.** A tree that does not
compile reads exactly like a green.

```
export AILANG_BIN=$HOME/.pinned-ailang/ailang
go vet ./... ; echo "vet_rc=$?"
./scripts/verify_ail.sh ; echo "ail_rc=$?"
go test ./... -count=1 > /tmp/i168_head_test.txt 2>&1 ; echo "test_rc=$?"
grep -c '^ok ' /tmp/i168_head_test.txt
grep -Ec '^(FAIL|--- FAIL)' /tmp/i168_head_test.txt
./scripts/verify_go.sh > /tmp/i168_head_vgo.txt 2>&1 ; echo "vgo_rc=$?"
grep -E '^--- FAIL' /tmp/i168_head_vgo.txt
bash scripts/mission_decisions.sh --check --file design_docs/world-mission.md ; echo "dec_rc=$?"
```

**Acceptance criteria**

- **AC7.1** `vet_rc=0` (= B20).
- **AC7.2** `ail_rc=0` and the output ends `✓ verify gate PASSED: 11 required identities verified, 40 named tests pass` (= B21).
- **AC7.3** `test_rc=0`, `^ok ` count = `19`, FAIL count = `0` (= B22).
- **AC7.4** `vgo_rc` is `0`, **or** `1` with the failing-test set a subset of
  `{TestCLIRealSubprocessEpisode}` — i.e. `grep -E '^--- FAIL' /tmp/i168_head_vgo.txt` names no
  other test. Base was rc=1 with exactly that one test (B23, §1.1). If any other test or any
  earlier stage fails, this is a REGRESSION and M7 blocks. If `vgo_rc=1`, re-run the isolated
  control `go test ./cmd/ailang-worldd/ -run TestCLIRealSubprocessEpisode -race -count=2` and
  record its rc (planner measured rc=0), so the flake attribution is evidence, not assertion.
- **AC7.5** `dec_rc=0`, `decision ledger valid: 22 rows` (= B8).
- **AC7.6** **Scope proof**: `git diff --name-only origin/dev...HEAD | grep -vc '^design_docs/'`
  → `0` (rc=1). Positive control on the same instrument:
  `git diff --name-only origin/dev...HEAD | grep -c '^design_docs/'` → the actual file count,
  which must be `≥ 9` and must include all four M3 files. *(A rc=0 on the first command would
  mean a Go/`.ail`/script file leaked in.)*
- **AC7.7** No instrument was committed: `git diff --name-only origin/dev...HEAD | grep -c 'i168'` → `0` (rc=1),
  and `git status --porcelain | grep -c '^??'` → `0` (rc=1) after staging.
- **AC7.8** Full-branch re-run of all three instruments at head: `/tmp/i168_headings.sh
  design_docs/world-mission-log.md` rc=0, `python3 /tmp/i168_index_check.py` rc=0,
  `bash /tmp/i168_rulings.sh` rc=0. All three must be re-run AFTER the iteration-168 record is
  added, not before.
- **AC7.9** The iteration-168 record is itself consistent: its log heading is canonical, its index
  row exists and matches (I2 rc=0 covers both), and its STATUS stamp is in the charter with the
  now-4th stamp (iteration 165) rotated into the archive, leaving `grep -c '^## STATUS 2026'
  design_docs/world-mission.md` = `3` with identities `168 167 166`.

**Mutation drill (M7)** — prove AC7.6's scope proof can fail:
```
touch host/NEGCTL_168_scope_probe.go && git add host/NEGCTL_168_scope_probe.go
git diff --name-only --cached origin/dev | grep -vc '^design_docs/' ; echo "KILLER_rc=$?"
git rm --cached -f host/NEGCTL_168_scope_probe.go && rm host/NEGCTL_168_scope_probe.go
git status --porcelain | grep -c 'NEGCTL_168_scope_probe' ; echo "restore_check (want 0, rc=1)"
```
**Named killer**: `grep -vc '^design_docs/'` prints **`1` with rc=0** while the probe is staged
(AC7.6 FAILS), and `0` with rc=1 after removal. Note the empty file would also break `go vet`
if left in place — a second, independent killer.

**Commit**: `docs(mission): World iteration 168 — reconcile the five stranded records (162-166)`

---

## 3. Landing

Open ONE PR from `mission/world-iter168-records` to `dev`, titled
`docs(mission): World iteration 168 — reconcile stranded records 162-166 (#127/#128 superseded)`.
The body must state:

- `D-WORLD-35` = A is the authority, quoted.
- The head's `verify_ail.sh` / `go vet` / `go test` results and the `verify_go.sh` rc with its
  failing-test set, alongside the **base** values from §1 — so "inherited" is demonstrated,
  not asserted, exactly as `D-WORLD-35` requires.
- That PRs **#127 and #128 are superseded and must be closed unmerged**, with the reason: their
  ledger hunks would reopen `D-WORLD-34` and `D-WORLD-35` as duplicate `OPEN` rows.
- The two out-of-scope defects from §1.2, as candidate queue rows.

Do not merge without the CI result. Do not close #127/#128 before this PR merges.

---

## 4. Risks / what could still be vacuous

1. **AC1.1's exact hash is the strongest and most brittle criterion.** It was produced by the
   planner running the M1 recipe verbatim. If the executor deviates by even one blank line the
   hash breaks — that is the point, but it means a hash mismatch needs diagnosis (`diff` against
   `/tmp/i168_newlog.md`), not a silent re-derivation of a new "expected" hash. **Never update an
   expected hash to match what you produced.**
2. **AC4.5 asserts a count that is 3 at base and 3 at head.** The count alone is vacuous; only the
   *identity* assertion (`167 166 165`) carries information. Kept because a count of 4 is a real
   failure mode (stamp added without rotating). Do not report AC4.5 as passing on the count alone.
3. **AC1.6 (`--- count = 8`) and AC4.1/AC4.4 are regression guards, true at base.** They are
   flagged as such. Their non-vacuity rests entirely on their mutation drills (M1 arm B, M0). If a
   drill is skipped, treat the criterion as unproven and say so.
4. **The reconstructed iteration-166 STATUS stamp is the sprint's one authored artefact.** Every
   other byte is moved, not written. It is a summary of a landed record and is labelled
   `RECONSTRUCTED AT ITERATION 168` (AC4.6), but no criterion can check that its summary is
   *faithful* — only a human or an evaluator reading log entry 166 can. Flag it explicitly for
   the evaluator.
5. **Historical stamps 162–165 assert "Ledger 21 rows / 3 OPEN" and "D-WORLD-32/33 remain
   OPEN".** Those are now false as current state and true as history. Nothing in the gate
   distinguishes the two. A future reader could mistake a rotated stamp for live state. The
   iteration-168 log entry must say so; consider a queue row for a "stamps are history" banner.
6. **`verify_go.sh`'s base red is a flake, so AC7.4 is a judgement call under load.** If the
   flake does not reproduce at head, AC7.4 passes trivially (rc=0) — which is fine but proves
   nothing about the sprint. If a *different* test fails, do not reclassify it as "another
   flake" without the isolated `-count=2` control that the plan names for
   `TestCLIRealSubprocessEpisode`.
7. **`check_mission_config()` in `verify_go.sh` ends with an `echo`**, so its return value is the
   echo's. It only fails loudly because the script runs under `set -euo pipefail` (line 17). That
   is a fragile coupling: anyone who wraps that call in an `if` or a `||` would silently disarm
   the ledger gate. Not this sprint's to fix; worth a queue row.
8. **Nothing in CI or `scripts/` reads `world-mission-log.md`, `world-mission-index.md`, or
   `world-mission-dashboard.md`** (measured: `grep -rq 'world-mission-log\|world-mission-index'
   .github/ scripts/ host/` → rc=1, with `grep -rq 'world-mission.md' …` → rc=0 as the positive
   control). So M1/M2/M6's correctness rests **entirely** on instruments I1 and I2, which live in
   `/tmp` and are not committed. If those are not run, those milestones have no gate at all. The
   durable fix — promoting I1+I2 into `scripts/` and wiring them into `verify_go.sh` — is a
   candidate queue row this sprint deliberately does not take, because it would move the sprint
   out of `design_docs/**` and void AC7.6.
9. **This plan file itself contains lines that look like mission records.** The reconstructed
   iteration-166 stamp is quoted verbatim inside a fenced block, so a repo-wide
   `grep -c '^## STATUS 2026'` or `grep -rc 'RECONSTRUCTED AT ITERATION 168'` will match it. Every
   criterion above is deliberately **file-scoped** to the one file it is about. Never widen one to
   `-r` over `design_docs/` — the plan would then satisfy its own criteria.
10. **AC7.6's "≥ 9 files" is a floor, not an identity.** The expected change set is exactly ten
   paths: the four M3 documents, `world-mission-log.md`, `world-mission-index.md`,
   `world-mission.md`, `world-mission-status-archive.md`, `world-mission-dashboard.md`, and this
   plan file. If the count is anything other than 10, list the diff and explain each extra path
   before proceeding.
