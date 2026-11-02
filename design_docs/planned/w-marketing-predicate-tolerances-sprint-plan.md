# Sprint plan: w-marketing-predicate-tolerances (queue row 238, iteration 336)

Status: **planned**. Planner: iteration 336 (claude-opus-5-5). Executor: one run per milestone,
following this file literally.
- **Design (authority):** `design_docs/planned/w-marketing-predicate-tolerances.md` r3 (`b9f67a27`).
  Read §2–§5 and §9 before starting. This plan does not re-design anything. Where it differs from
  the design's numbers, the difference was measured and is listed in §0.
- **Worktree:** `/Users/voightkampff/dev/sunholo-data/.wt-world-iter336-predicate-tolerances`,
  branch `sprint/world-iter336-marketing-predicate-tolerances`.
- **Executor scratch (never inside the worktree):** `/Users/voightkampff/dev/sunholo-data/.scratch-iter336-executor/`.
- **Planner evidence:** `design_docs/verification/world-iter336/plan-*.txt` (§9 Verification Log).
- **The executor writes no Go code and no doc text by hand.** It applies two banked unified diffs,
  checks their sha256, runs the gates, and banks small `.txt` evidence.

| Milestone | What | Files | Commit |
|---|---|---|---|
| M1 | predicates + tests | `scripts/marketingcheck/{position.go,phone.go,position_test.go,phone_test.go}` via `plan-M1.diff.txt` | executor commits M1 (no push) |
| M2 | orchestration text A1–A7 | `design_docs/planned/w-marketing-capture-orchestration.md` via `plan-M2.diff.txt` | **executor does NOT commit**; the controller signs off, runs the orchestration quorum, then commits |

## 0. What the planner measured that the design did not fix (read this first)

The design (§2.3, §5) told M1 to measure arm `B7-thumb-two-equal-bridged-runs` first, and to stop if
the result differed. The planner did that measurement (P03–P06). **The result matches the design's
intent**: B7 is ACCEPT and selects the left run; the oracle REFUSEs it; the tie mutant kills it.
So B7 is in `plan-M1.diff.txt`, and these numbers change from the design text:

| Item | Design text (r3) | Measured, used by this plan |
|---|---|---|
| B7 verdict / thumb rect (render and real seed 5) | ACCEPT, X0 6, "6..139" | ACCEPT, `{6 833 139 841}` (X1 exclusive) |
| B7 oracle verdict | REFUSE (derived) | REFUSE, oracle thumb `{-1 -1 -1 -1}` (render and real) |
| Phone ledger | 3 entries (then "4 if B7 measures as derived") | **4**: `new/header-15px`, B1, B1c, B7 |
| `PARITY phone` (CI renders) | `inputs=44 diffs=0 ledger=3` | **`inputs=45 diffs=0 ledger=4`** |
| `PARITY phone` (real corpus) | `inputs=62 diffs=0 ledger=3` | **`inputs=63 diffs=0 ledger=4`** |
| `CHECK` line | `CHECK images=10 byte-equal=10` | full line `CHECK images=10 byte-equal=10; real arms=43; parity inputs=63` |
| Phone arms / CI inputs / real inputs | 42 / 44 / 62 | 43 / 45 / 63 |
| `phoneWant` | 44 rows, 14 ACCEPT | 45 rows, 15 ACCEPT |
| `M238-run-nobridge` must-red | B1, B1c | B1, B1c **and B7** |
| `M238-run-rightmost-tie` | "derived" | measured: rc 1, reds exactly `…/judge/B7-thumb-two-equal-bridged-runs` (thumb becomes `{200 833 333 841}`) |
| M2 A3 / A7 text | "phone 3" | "phone 4" (the design's §2.3 pre-authorises exactly this change) |
| M2 acceptance grep | `position 8 entries, phone 3` | `position 8 entries, phone 4` |
| Patch size | r2: 4 files +210/−39 | final: 4 files **+248/−60** (r2 + B7 + the hand-update comments §5 M1 asks for) |

The position side is unchanged from the design: `PARITY position: inputs=51 verdict-diffs=0 reason-diffs=0 ledger=8`.
The tie mutant **survives** without B7 (P06 control), so B7 is the only arm that kills it.
Because the planner already measured B7, the executor does not re-measure it before applying:
M1.2's acceptance lines and M1.6's mutant table are the measured values, and a mismatch is a STOP.
Beyond r2 and B7, the final patch only edits comments: the stale counts design §5 M1 lists, the r3-only
file headers, the `phoneParity` doc comment (it still named one ledger entry), and the row-238 code
comments' section references (§3/§4 → §2.3/§2.4/§2.5, where the design defines arms and ledgers).

## 1. Executor rules (each one is a STOP if broken)

1. **Shell.** The Bash tool's shell is zsh. Run every block in this plan exactly as written, as a
   `bash <<'SH' … SH` heredoc. Do not retype it into zsh line by line.
2. **Runner.** Every command that builds or tests runs under `bash "$RB" <limit_s> "$OUT" <cmd…>`;
   read `rc=$?` on the next line. `$OUT` gets stdout, `$OUT.err` gets stderr. **rc 124 means the
   deadline fired: STOP the milestone and report. There is no continuation to M2.**
3. **Limits (design §5):** `go vet`, `gofmt -l`, `go build ./...`: 600 s each. `go test ./...`:
   2400 s. `./scripts/verify_ail.sh`: 1200 s. The `TestMarketing(Position|Phone)` acceptance run:
   1500 s, with `-timeout 20m` inside. Each mutation run: 600 s. Each capture-day control: 300 s.
4. **Long commands.** `go test ./...` and `verify_ail.sh` can outlast one foreground tool call
   (600 s). Launch them with the Bash tool's own `run_in_background: true` and poll with the
   bounded loop given in M1.5. **Never use `nohup`** or any SIGHUP-ignoring wrapper (row 231 D-14:
   `nohup` reds `TestMarketingSupervisorFakeDaemonCleanup/shell-hup`).
5. **No pipes before an rc is read.** Redirect to a file, read rc, then grep the file.
6. **Caches.** `GOCACHE` and `GOMODCACHE` live under `$X` (set by `env.sh`). Nothing large is ever
   written under `design_docs/` (iteration 335 wrote 4.6 GB of caches there by accident). Evidence
   files under `$BANK` are small `.txt` only.
7. **Mutations run only in `$X/mut`**, never in the worktree.
8. **Do not edit** the design doc, this plan, the sprint JSON, any `plan-*.txt`, `run_bounded.sh.txt`,
   `design_docs/world-mission*.md`, `tools/launchd/*`, `.github/`, `go.mod`, fixtures, or the oracle
   sources. No push, no PR, no amend, no rebase. Never touch `~/.ailang/state/mission-v1*`.
9. **Any expected value below that does not match is a STOP.** Report the command, the expected
   value and the actual output. Do not "fix" the test, the patch or the plan.

## 2. Environment (run once, first)

`AILANG_BIN` is `$HOME/.pinned-ailang/ailang`, which reports `AILANG v0.41.0`. That is the binary
iteration 335's controller gate used (`verification/world-iter335/controller-fullgate.txt` line 2:
`AILANG_BIN=AILANG v0.41.0`), the pin `scripts/verify_go.sh` enforces by exact token (lines 12–17,
181–198), and CI's `WORLD_PKG_AILANG_BIN=$HOME/.pinned-ailang/ailang` (`.github/workflows/ci.yml:108`).
`scripts/verify_ail.sh` refuses `-dirty` and `-N-g<sha>` dev builds (lines 64–118).

```bash
bash <<'SH'
set -u
X=/Users/voightkampff/dev/sunholo-data/.scratch-iter336-executor
WT=/Users/voightkampff/dev/sunholo-data/.wt-world-iter336-predicate-tolerances
BANK=$WT/design_docs/verification/world-iter336
mkdir -p "$X/gocache" "$X/gomodcache" "$X/logs" "$X/mutlogs" "$X/pristine"
cat > "$X/env.sh" <<EOF
export X=$X WT=$WT BANK=$BANK
export PATH=/opt/homebrew/bin:\$PATH
export GOCACHE=$X/gocache GOMODCACHE=$X/gomodcache GOFLAGS=-modcacherw
export AILANG_BIN=\$HOME/.pinned-ailang/ailang
export WORLD_PHONE_CORPUS=\$HOME/.ailang/state/world-iter328-phone-corpus
export RB=$X/run_bounded.sh
EOF
cp "$BANK/run_bounded.sh.txt" "$X/run_bounded.sh"
cp "$BANK/plan-mutate.py.txt" "$X/mutate.py"
cp "$BANK/plan-mutants-expected.txt" "$X/mutants-expected.txt"
source "$X/env.sh"
shasum -a 256 "$X/run_bounded.sh" "$X/mutate.py" "$X/mutants-expected.txt" "$BANK/plan-M1.diff.txt" "$BANK/plan-M2.diff.txt"
"$AILANG_BIN" --version 2>/dev/null | head -1
(cd "$WT" && go version)
(cd "$WORLD_PHONE_CORPUS" && shasum -a 256 -c "$WT/scripts/marketingcheck/testdata/phone/corpus-SHA256SUMS" > "$X/logs/corpus.txt" 2>&1; echo "corpus rc=$?"; grep -c ': OK$' "$X/logs/corpus.txt")
SH
```

Expected (every line; any difference is a STOP):
```
1635f8d7308e03a6bc02eeacd4bd6eb723ce76405c4c1d28f2c45ef93d8e2082  …/run_bounded.sh
95206757b724bfd95e86e9a253831ac840d73635a86877f4bd06f49f6754d749  …/mutate.py
fa71baedaacc7dfb9064d3ab8b86f3339cf516f61d57e2b05beeedacc147833b  …/mutants-expected.txt
7da36f3006efca03e0b49249d60e577097df799e09b8186ce9525217dc585d74  …/plan-M1.diff.txt
08e9838f84f3eaaa9f7aef6f1e6ad8fbea2fe18450c11c738806a2e2219818e9  …/plan-M2.diff.txt
AILANG v0.41.0
go version go1.26.6 darwin/arm64
corpus rc=0
11
```
`go version` must run inside the worktree: outside it the Homebrew `go` is 1.26.4, and inside it
`GOTOOLCHAIN=auto` selects go1.26.6 from `go.mod` (the first run downloads it into `$X/gomodcache`).
A missing corpus or any non-OK line is BLOCKED (design §5 acceptance 1; orchestration §4.4 step 1):
STOP and report; never regenerate the corpus.

Every later block starts with `source /Users/voightkampff/dev/sunholo-data/.scratch-iter336-executor/env.sh`.

## 3. M1: predicates and tests

### M1.0 Base check (no concurrent change since the design's base)

```bash
bash <<'SH'
source /Users/voightkampff/dev/sunholo-data/.scratch-iter336-executor/env.sh
cd "$WT"
git diff --exit-code 37100a43 -- scripts/marketingcheck design_docs/planned/w-marketing-capture-orchestration.md > /dev/null; echo "base-diff rc=$?"
git status --porcelain -- scripts/ design_docs/planned/w-marketing-capture-orchestration.md | wc -l | tr -d ' '
shasum -a 256 scripts/marketingcheck/position.go scripts/marketingcheck/phone.go scripts/marketingcheck/position_test.go scripts/marketingcheck/phone_test.go design_docs/planned/w-marketing-capture-orchestration.md
SH
```
Expected:
```
base-diff rc=0
0
7512337137e3eb08fc7930247c89ee1e6b1f392fe9014069823b22aabfaac6ac  scripts/marketingcheck/position.go
e3b8db55c53cf6d69c31c0f908b8302b81dc24277e80af6b887473358f28c4b5  scripts/marketingcheck/phone.go
8e86c2563e24ab03dc188d4818081023c7a65f73dd97d8bd188b94507e328a0d  scripts/marketingcheck/position_test.go
0bce70c4f72784d4fea1437abd636c76e90874301b2fc97453d8e006cf0f71a8  scripts/marketingcheck/phone_test.go
0b2ee85987beac8d95b82d6a36a61e28b6b8f269dc3d3b0ff100817091db0e9c  design_docs/planned/w-marketing-capture-orchestration.md
```
A difference means a concurrent change landed (design §6): STOP; the mission re-bases and re-plans.

### M1.1 Apply the patch

```bash
bash <<'SH'
source /Users/voightkampff/dev/sunholo-data/.scratch-iter336-executor/env.sh
cd "$WT"
git apply --check "$BANK/plan-M1.diff.txt"; echo "check rc=$?"
git apply "$BANK/plan-M1.diff.txt"; echo "apply rc=$?"
git diff --stat -- scripts/marketingcheck
shasum -a 256 scripts/marketingcheck/position.go scripts/marketingcheck/phone.go scripts/marketingcheck/position_test.go scripts/marketingcheck/phone_test.go
SH
```
Expected:
```
check rc=0
apply rc=0
 scripts/marketingcheck/phone.go         |  31 ++++---
 scripts/marketingcheck/phone_test.go    |  76 ++++++++++++-----
 scripts/marketingcheck/position.go      |  57 +++++++++++--
 scripts/marketingcheck/position_test.go | 144 +++++++++++++++++++++++++++-----
 4 files changed, 248 insertions(+), 60 deletions(-)
9029dbabd574761050399864cdd45f399ae9a3682ff1523e547479e75f8c669d  scripts/marketingcheck/position.go
9540d55a323640edab4933e0a4d8f8445be5046a5f8d6c1b0596035b1f0c0770  scripts/marketingcheck/phone.go
7e9b1cc6e91aacb785ca3f32533d0dbf45f966b34d6e2920dd5b1e95c05cc35c  scripts/marketingcheck/position_test.go
e777985d087d882f3d4bd6ace1d6ed0580f8efd41a5914ba19a1db10798db2a0  scripts/marketingcheck/phone_test.go
```
The patch already contains every count and comment update that design §5 M1 asks for by hand
("the 44 position inputs", "44-row … 19 ACCEPT", "41-row … 12 ACCEPT", "39 arms", and the r3-only
file headers of `position.go` and `phone.go`). Do not edit anything by hand.

### M1.2 Acceptance 1: the marketingcheck run (1500 s)

```bash
bash <<'SH'
source /Users/voightkampff/dev/sunholo-data/.scratch-iter336-executor/env.sh
cd "$WT"
LOG="$X/logs/m1-accept.log"
bash "$RB" 1500 "$LOG" go test -count=1 -timeout 20m -v -run '^TestMarketing(Position|Phone)' ./scripts/marketingcheck; echo "rc=$?"
echo "skip=$(grep -c '^--- SKIP' "$LOG")"
echo "fail=$(grep -c -- '--- FAIL' "$LOG")"
echo "realcorpus=$(grep -c '^--- PASS: TestMarketingPhoneRealCorpus' "$LOG")"
echo "b7=$(grep -c -- '--- PASS: TestMarketingPhoneArms/judge/B7-thumb-two-equal-bridged-runs' "$LOG")"
for l in 'PARITY position: inputs=51 verdict-diffs=0 reason-diffs=0 ledger=8' 'PARITY phone: inputs=45 diffs=0 ledger=4' 'PARITY phone: inputs=63 diffs=0 ledger=4' 'CHECK images=10 byte-equal=10; real arms=43; parity inputs=63'; do echo "$(grep -cF -- "$l" "$LOG") $l"; done
grep -E '^(ok|FAIL)' "$LOG"
SH
```
Expected:
```
rc=0
skip=0
fail=0
realcorpus=1
b7=1
1 PARITY position: inputs=51 verdict-diffs=0 reason-diffs=0 ledger=8
1 PARITY phone: inputs=45 diffs=0 ledger=4
1 PARITY phone: inputs=63 diffs=0 ledger=4
1 CHECK images=10 byte-equal=10; real arms=43; parity inputs=63
ok  	github.com/sunholo-data/ailang-world/scripts/marketingcheck	<time>
```

### M1.3 vet, gofmt, build (600 s each)

```bash
bash <<'SH'
source /Users/voightkampff/dev/sunholo-data/.scratch-iter336-executor/env.sh
cd "$WT"
bash "$RB" 600 "$X/logs/vet.txt" go vet ./scripts/marketingcheck; echo "vet rc=$?"
bash "$RB" 600 "$X/logs/gofmt.txt" gofmt -l scripts/marketingcheck; echo "gofmt rc=$? lines=$(wc -l < "$X/logs/gofmt.txt" | tr -d ' ')"
bash "$RB" 600 "$X/logs/build.txt" go build ./...; echo "build rc=$?"
git diff --exit-code 37100a43 -- design_docs/verification/world-iter315/ scripts/marketingcheck/testdata/ > /dev/null; echo "fixtures rc=$?"
SH
```
Expected: `vet rc=0`, `gofmt rc=0 lines=0`, `build rc=0`, `fixtures rc=0` (design §5 acceptance 2 and 3).

### M1.4 Capture-day controls, steps 2–5 (300 s each; design §3 / V14)

```bash
bash <<'SH'
source /Users/voightkampff/dev/sunholo-data/.scratch-iter336-executor/env.sh
cd "$WT"
C=$X/logs/control
bash "$RB" 300 "$C-2.txt" go run ./scripts/marketingcheck position --baseline design_docs/verification/world-iter312/r1-rects.json accept-seed5=scripts/marketingcheck/testdata/position/stores/r1c-seed5.json accept-eval-p1=scripts/marketingcheck/testdata/position/stores/eval-p1.json; echo "step2 rc=$?"; cat "$C-2.txt"
bash "$RB" 300 "$C-3.txt" go run ./scripts/marketingcheck position --baseline design_docs/verification/world-iter312/r1-rects.json refuse=scripts/marketingcheck/testdata/position/v65/neg-attested+14.json; echo "step3 rc=$?"; cat "$C-3.txt"
bash "$RB" 300 "$C-4.txt" go run ./scripts/marketingcheck phone accept-v50=$HOME/.ailang/state/world-iter305/probe-seeded/pass-a/phone.png accept-seed2=$HOME/.ailang/state/world-iter328-phone-corpus/r1-seed2-pass-a-phone.png; echo "step4 rc=$?"; cat "$C-4.txt"
bash "$RB" 300 "$C-5.txt" go run ./scripts/marketingcheck phone refuse=$HOME/.ailang/state/world-iter328-phone-corpus/r1c-seed5-pass-a-workbench-light.png; echo "step5 rc=$?"; cat "$C-5.txt"
SH
```
Expected (planner-measured, P09; identical to designer V14):
```
step2 rc=0
accept-seed5                           ACCEPT regions=17
accept-eval-p1                         ACCEPT regions=17
TOTAL candidates=2 accept=2 refuse=0
step3 rc=1
refuse                                 REFUSE regions=17
    packet.png/roster-attested: (a3) gap after roster-tested 29 vs R1 15
    packet.png/roster-claimed: (a3) gap after roster-attested 1 vs R1 15
TOTAL candidates=1 accept=0 refuse=1
step4 rc=0
accept-v50                   ACCEPT 390x844 bands=21 header={25 45 370 68} entry={41 428 364 610} thumb={6 833 329 841} titleRight=389
accept-seed2                 ACCEPT 390x844 bands=21 header={25 45 370 68} entry={41 428 364 610} thumb={6 833 340 841} titleRight=389
TOTAL candidates=2 accept=2 refuse=0
step5 rc=1
refuse                       REFUSE 1440x1100 bands=28 header={25 45 370 68} entry={90 428 370 610} thumb={-1 -1 -1 -1} titleRight=413
    (i) header {25 45 370 68} entry {90 428 370 610} vs R1 {25 45 370 68} {41 428 364 610} ±4
    (iii′) dims 1440x1100
    (iii′) thumb {-1 -1 -1 -1}
TOTAL candidates=1 accept=0 refuse=1
```
(Step 5 uses the corpus's seed-5 desktop shot as the stand-in for "capture day's
`pass-a/workbench-light.png`", as V14 did. `go run` prints `exit status 1` into `$C-5.txt.err`; that is expected.)

### M1.5 Full gates: `go test ./...` (2400 s) and `verify_ail.sh` (1200 s)

These run in the **worktree** (a git checkout). `go test ./...` in a `git archive` copy reds
`host/pins TestReleaseLiteralsExplainedByManifest` (`git ls-files: exit status 128`, P10): that is
an instrument artifact of the copy, never a reason to run the suite outside the worktree.

Launch with the Bash tool parameter `run_in_background: true` (not `nohup`, not `&`):
```bash
bash <<'SH'
source /Users/voightkampff/dev/sunholo-data/.scratch-iter336-executor/env.sh
cd "$WT"
rm -f "$X/logs/fullgate.status"
echo "start $(date -u +%FT%TZ) AILANG_BIN=$AILANG_BIN ($("$AILANG_BIN" --version 2>/dev/null | head -1))" > "$X/logs/fullgate.status"
bash "$RB" 2400 "$X/logs/gotest-all.txt" go test ./...; echo "go test ./... rc=$?" >> "$X/logs/fullgate.status"
bash "$RB" 1200 "$X/logs/verify_ail.txt" ./scripts/verify_ail.sh; echo "verify_ail.sh rc=$?" >> "$X/logs/fullgate.status"
echo "DONE $(date -u +%FT%TZ)" >> "$X/logs/fullgate.status"
SH
```
Poll (foreground, Bash tool timeout 600000 ms; repeat this call until `DONE` appears; the gate's
own deadlines bound the total at 3600 s):
```bash
bash <<'SH'
F=/Users/voightkampff/dev/sunholo-data/.scratch-iter336-executor/logs/fullgate.status
end=$(( $(date +%s) + 540 ))
while [ "$(date +%s)" -lt "$end" ]; do grep -q '^DONE' "$F" && break; sleep 15; done
cat "$F"
SH
```
Then:
```bash
bash <<'SH'
source /Users/voightkampff/dev/sunholo-data/.scratch-iter336-executor/env.sh
echo "ok=$(grep -c '^ok' "$X/logs/gotest-all.txt") fail=$(grep -c '^FAIL\|--- FAIL' "$X/logs/gotest-all.txt")"
grep -F 'scripts/marketingcheck' "$X/logs/gotest-all.txt"
tail -1 "$X/logs/verify_ail.txt"
SH
```
Expected: the status file shows `AILANG v0.41.0`, `go test ./... rc=0`, `verify_ail.sh rc=0`;
`fail=0`; an `ok … scripts/marketingcheck` line; and the last `verify_ail.txt` line
`✓ verify gate PASSED: 16 required identities verified, 40 named tests pass, 503 se-tools named tests pass`.
(Planner timing, P10/P11: `go test ./...` 376 s, `verify_ail.sh` under 6.5 min, on this host.)
A failure in a package outside `scripts/marketingcheck`: STOP and report the package, the test
and its output. Do not re-run to "get a green".

### M1.6 Mutations (each alone, 600 s each, in `$X/mut`)

The mutation tree is a fresh `git archive` of the worktree's HEAD plus the same patch, so it equals
the M1 tree byte for byte (sha check below). The driver `mutate.py` makes one exact string
replacement per mutant and refuses (exit 2, file untouched) if the old string is not unique.
The expected `--- FAIL:` names for each mutant (parents included) are in `mutants-expected.txt`;
the run must produce **exactly** that set (`match=EXACT`).

```bash
bash <<'SH'
source /Users/voightkampff/dev/sunholo-data/.scratch-iter336-executor/env.sh
rm -rf "$X/mut" && mkdir -p "$X/mut"
git -C "$WT" archive HEAD | tar -x -C "$X/mut"
cd "$X/mut" && git apply "$BANK/plan-M1.diff.txt"; echo "mut apply rc=$?"
cp scripts/marketingcheck/position.go scripts/marketingcheck/phone.go "$X/pristine/"
shasum -a 256 scripts/marketingcheck/position.go scripts/marketingcheck/phone.go scripts/marketingcheck/position_test.go scripts/marketingcheck/phone_test.go
SH
```
Expected: `mut apply rc=0` and the four M1.1 sha256 values (`9029dbab…`, `9540d55a…`, `7e9b1cc6…`, `e777985d…`).
If `git apply` fails here because HEAD already contains M1 (you committed first), use
`git -C "$WT" archive HEAD` alone and skip the apply; the sha check must still match.

Then run all 14 (Bash tool timeout 600000 ms; planner total 43 s):
```bash
bash <<'SH'
source /Users/voightkampff/dev/sunholo-data/.scratch-iter336-executor/env.sh
cd "$X/mut"
OUT=$X/logs/m1-mutants.txt; : > "$OUT"
for ID in M238-a3-unbounded M238-a3-bound2 M238-a3-bound0 M238-a5-off M238-a5-tol3 M238-a5-tol1 M238-a5-nonpara-zero M238-a5-accumulate-failed M238-run-nobridge M238-run-gap2px M238-run-chain-bridges M238-run-rightmost-tie F4 F9; do
  python3 -I "$X/mutate.py" "$ID" "$X/mut" >> "$OUT" 2>&1 || { echo "MUTANT $ID NOT-APPLIED: STOP" >> "$OUT"; break; }
  bash "$RB" 600 "$X/mutlogs/$ID.log" go test -count=1 -timeout 300s -v -run '^TestMarketing(Position|Phone)(Corpus|Arms|Parity)$' ./scripts/marketingcheck; rc=$?
  got=$(grep -E '^ *--- FAIL: ' "$X/mutlogs/$ID.log" | sed -E 's/^ *--- FAIL: ([^ ]+) .*/\1/' | sort)
  want=$(grep "^$ID " "$X/mutants-expected.txt" | cut -d' ' -f2 | sort)
  if [ "$got" = "$want" ]; then m=EXACT; else m=MISMATCH; fi
  echo "MUTANT $ID rc=$rc match=$m" >> "$OUT"
  printf '%s\n' "$got" | sed 's/^/    --- FAIL: /' >> "$OUT"
  [ "$m" = MISMATCH ] && diff <(printf '%s\n' "$want") <(printf '%s\n' "$got") | sed 's/^/    DIFF /' >> "$OUT"
  cp "$X/pristine/position.go" "$X/pristine/phone.go" scripts/marketingcheck/
  echo "    restore: $(shasum -a 256 scripts/marketingcheck/position.go | cut -c1-16) $(shasum -a 256 scripts/marketingcheck/phone.go | cut -c1-16)" >> "$OUT"
  [ "$rc" = 124 ] && { echo "TIMEOUT: STOP" >> "$OUT"; break; }
done
grep -E '^MUTANT|restore|DIFF|STOP' "$OUT"
SH
```
Expected: 14 `MUTANT <id> rc=1 match=EXACT` lines, each followed by
`    restore: 9029dbabd5747610 9540d55a323640ed` (byte-identical restore). Any `rc=0`, `rc=124`,
`MISMATCH`, `NOT-APPLIED` or other restore hash is a STOP.

The mutants and their must-red subtests (from design §5, with the two planner-measured changes marked **bold**;
each also reds its parent `TestMarketingPositionCorpus`/`TestMarketingPhoneArms`):

| Id | Mutation (old → new, exact, in `mutate.py`) | Must-red subtests |
|---|---|---|
| M238-a3-unbounded | position.go `posLineStep(e) && posAbs(posLines(e)) <= positionFreeGapLines` → `posLineStep(e)` | PositionCorpus `new/perm-approvals-y+48`, `new/perm-approvals-y-48`, `new/perm-validation-h+216`, `new/bnd-createdAt-wraps-2-lines`; `TestMarketingPositionParity` |
| M238-a3-bound2 | `positionFreeGapLines = 1` → `= 2` | `new/perm-approvals-y+48`, `new/perm-approvals-y-48`, `new/bnd-createdAt-wraps-2-lines`; PositionParity |
| M238-a3-bound0 | `positionFreeGapLines = 1` → `= 0` | `v65/syn-pos-createdAt-wraps`, `new/iso-minh-approvals-h-1`, `new/perm-approvals-y+24`, `new/perm-approvals-y-24`; PositionParity |
| M238-a5-off | `if chain {` (before `drift += prevRes`) → `if chain && false {` | `new/perm-roster-drift-2px`, `new/bnd-drift-sum-3`, `new/bnd-nonpara-drift-3`; PositionParity |
| M238-a5-tol3 | `positionDriftTol     = 2` → `= 3` | `new/bnd-drift-sum-3`, `new/bnd-nonpara-drift-3`; PositionParity |
| M238-a5-tol1 | `positionDriftTol     = 2` → `= 1` | `new/bnd-drift-sum-2`, `new/bnd-nonpara-drift-2`; PositionParity |
| M238-a5-nonpara-zero | `hOK, hRes := true, dh` → `hOK, hRes := true, 0` | `new/bnd-nonpara-drift-3`; PositionParity |
| M238-a5-accumulate-failed | `if okGap { drift += gRes }` → `drift += gRes` | `stores/r1-seed1`, `stores/r1-seed2`, `v65/neg-attested+14`, `v65/neg-dark-timeline-first+24`, `v65/neg-light-header+40`, `v65/neg-light-header-whole-view`, `v65/neg-rows-reordered-tested-claimed`, `v65/neg-tested-64x16-glyphs`, `v65/neg-tested-at-attested-row`, `v65/neg-tested-at-proven-row`, `v65-p1/neg-tested-at-attested-row` (11); PositionParity |
| M238-run-nobridge | phone.go `segs[i+1][0]-sg[1] == 2` → `== -1` | PhoneArms `judge/B1-thumb-split-px150`, `judge/B1c-thumb-split-px8`, **`judge/B7-thumb-two-equal-bridged-runs`**; `TestMarketingPhoneParity` (3 ledger entries) |
| M238-run-gap2px | `segs[i+1][0]-sg[1] == 2` → `<= 3` | `judge/B1e-thumb-gap-2px-at-150`; PhoneParity |
| M238-run-chain-bridges | the `if i+1 < len(segs) && … == 2 { end = segs[i+1][1] }` block → a `for j := i; …` loop | `judge/B1d-thumb-two-splits-px100-px200`, `judge/B6-alternating-pixels-no-thumb`; PhoneParity |
| **M238-run-rightmost-tie** | `if x0 < 0 \|\| end-sg[0] > x1-x0 {` → `>=` | `judge/B7-thumb-two-equal-bridged-runs`; PhoneParity (measured: thumb `{200 833 333 841}`) |
| F4 (row 231) | `s.Scrollbar.Y0 >= 0 && s.Scrollbar.X0 <= 8 && ` → `s.Scrollbar.Y0 >= 0 && ` | `judge/B5-two-pieces-rows-differ`, `judge/short-thumb-x20`, `judge/thumb-x9`, `judge/B1d-…`, `judge/B1e-…`; PhoneParity |
| F9 (row 231) | collapse `segs` to one first-to-last span before the run loop | `judge/scattered-pixels-no-thumb`, `judge/B1d-…`, `judge/B1e-…`, `judge/B6-…`; PhoneParity |

Each must-red list is a load-bearing claim discharged by its mutant (coding standards S6).

### M1.7 Bank evidence and commit M1

```bash
bash <<'SH'
source /Users/voightkampff/dev/sunholo-data/.scratch-iter336-executor/env.sh
cd "$WT"
{ echo "# executor M1 acceptance, $(date -u +%FT%TZ), HEAD $(git rev-parse --short HEAD) + plan-M1.diff.txt"
  grep -E 'PARITY|CHECK |^(ok|FAIL)|^--- (PASS|FAIL|SKIP)' "$X/logs/m1-accept.log"
  echo "vet/gofmt/build:"; cat "$X/logs/vet.txt" "$X/logs/gofmt.txt" "$X/logs/build.txt"
  for i in 2 3 4 5; do echo "control step $i:"; cat "$X/logs/control-$i.txt"; done
  echo "full gate:"; cat "$X/logs/fullgate.status"; grep -v '^ok\|no test files' "$X/logs/gotest-all.txt"; tail -1 "$X/logs/verify_ail.txt"
} > "$BANK/executor-M1-acceptance.txt"
cp "$X/logs/m1-mutants.txt" "$BANK/executor-M1-mutants.txt"
wc -c "$BANK/executor-M1-acceptance.txt" "$BANK/executor-M1-mutants.txt"
git add scripts/marketingcheck/position.go scripts/marketingcheck/phone.go scripts/marketingcheck/position_test.go scripts/marketingcheck/phone_test.go "$BANK/executor-M1-acceptance.txt" "$BANK/executor-M1-mutants.txt"
git status --porcelain
SH
```
Expected: both evidence files are under 40 kB; `git status --porcelain` shows exactly the six
staged paths as `M `/`A ` (plus any untracked files the controller left, which you do not add).
Commit with the message `feat(row 238 M1): bound free gaps to one line, add (a5) drift, bridge a 1 px thumb split`
followed by a blank line, one line `Plan: design_docs/planned/w-marketing-predicate-tolerances-sprint-plan.md M1`,
and your session's attribution lines. Do not push.

## 4. M2: orchestration text (after M1 is committed)

### M2.1 Pre-paste anchor check (hard stop; design §4)

```bash
bash <<'SH'
source /Users/voightkampff/dev/sunholo-data/.scratch-iter336-executor/env.sh
cd "$WT"; D=design_docs/planned/w-marketing-capture-orchestration.md
for a in 'the gap may instead differ from the R1 gap by a whole number of 24 px body lines, ±2 px;' "Every other region's height is within ±2 px of R1;" 'or any verdict difference outside the one-entry ledger, is a red.' 'moved by anything other than whole lines in the two body-text gaps' '(a luma 120–200 run of at least 100 px)' 'and packet paragraph heights and the two body-text gaps in whole 24 px lines ±2 px (§4.4, r3)' '## 10. Verification Log'; do printf '%s\t%s\n' "$(grep -cF -- "$a" "$D")" "$a"; done
printf '%s\t%s\n' "$(grep -cF -- '(a3)' "$D")" 'control (a3)'
printf '%s\t%s\n' "$(grep -cF -- 'the gap may instead differ from the R1 gap by a whole number of 25 px' "$D")" 'control misspelled'
shasum -a 256 "$D"
SH
```
Expected: `1` for each of the seven anchors, `6` for the `(a3)` control, `0` for the misspelled
control, and sha256 `0b2ee85987beac8d95b82d6a36a61e28b6b8f269dc3d3b0ff100817091db0e9c`. Anything
else: STOP; the mission re-bases, re-runs V20 and asks the controller. Never paste by line number.

### M2.2 Apply

```bash
bash <<'SH'
source /Users/voightkampff/dev/sunholo-data/.scratch-iter336-executor/env.sh
cd "$WT"; D=design_docs/planned/w-marketing-capture-orchestration.md
git apply --check "$BANK/plan-M2.diff.txt"; echo "check rc=$?"
git apply "$BANK/plan-M2.diff.txt"; echo "apply rc=$?"
git diff --stat -- "$D"
shasum -a 256 "$D"
SH
```
Expected: `check rc=0`, `apply rc=0`, `1 file changed, 10 insertions(+), 5 deletions(-)`, sha256
`e66293d11e303218b3e1259872064b311b9cef4917e7f67e064cc59754b4ccc2`.

The diff pastes A1–A7 of design §4 verbatim, with exactly two measured deltas: A3 reads
`(row 238: position 8 entries, phone 4; …)` and A7 reads `(position 8, phone 4)`, because B7 is the
fourth phone ledger entry (§0). Design §2.3 pre-authorises this change.

### M2.3 Acceptance

```bash
bash <<'SH'
source /Users/voightkampff/dev/sunholo-data/.scratch-iter336-executor/env.sh
cd "$WT"; D=design_docs/planned/w-marketing-capture-orchestration.md; OUT=$X/logs/m2.txt
bash "$RB" 600 "$X/logs/m2-parity.log" go test -count=1 -timeout 300s -v -run '^TestMarketing(Position|Phone)Parity$' ./scripts/marketingcheck; echo "parity rc=$?" > "$OUT"
grep -E -- '--- (PASS|FAIL)|PARITY' "$X/logs/m2-parity.log" | sed -E 's/ \([0-9.]+s\)$//' >> "$OUT"
for a in '(a5) (row 238) cumulative drift' 'exactly one off-range pixel' '### Row 238 (iteration 336' 'position 8 entries, phone 4' 'outside the one-entry ledger' '0a9a6795d9d711f99d77030e0491cd1c5e254c61fc7434898a80b4875404927c' '31576922bc566ff9d6c2c40b4cfa41e71f234fc706d013377f657a0ffd010b4c'; do printf '%s\t%s\n' "$(grep -cF -- "$a" "$D")" "$a" >> "$OUT"; done
cat "$OUT"
SH
```
Expected:
```
parity rc=0
    phone_test.go:589: PARITY phone: inputs=45 diffs=0 ledger=4
--- PASS: TestMarketingPhoneParity
    --- PASS: TestMarketingPhoneParity/pin-registry
    position_test.go:653: PARITY position: inputs=51 verdict-diffs=0 reason-diffs=0 ledger=8
--- PASS: TestMarketingPositionParity
    --- PASS: TestMarketingPositionParity/pin-registry
1	(a5) (row 238) cumulative drift
1	exactly one off-range pixel
1	### Row 238 (iteration 336
1	position 8 entries, phone 4
0	outside the one-entry ledger
1	0a9a6795d9d711f99d77030e0491cd1c5e254c61fc7434898a80b4875404927c
1	31576922bc566ff9d6c2c40b4cfa41e71f234fc706d013377f657a0ffd010b4c
```

### M2.4 Mutation M238-doc-pin-dup (both pins, each alone, 600 s)

```bash
bash <<'SH'
source /Users/voightkampff/dev/sunholo-data/.scratch-iter336-executor/env.sh
cd "$WT"; D=design_docs/planned/w-marketing-capture-orchestration.md; OUT=$X/logs/m2-pindup.txt; : > "$OUT"
cp "$D" "$X/pristine/orchestration.md"
for pin in 0a9a6795d9d711f99d77030e0491cd1c5e254c61fc7434898a80b4875404927c 31576922bc566ff9d6c2c40b4cfa41e71f234fc706d013377f657a0ffd010b4c; do
  python3 -I - "$D" "$pin" <<'PY'
import sys
p, pin = sys.argv[1], sys.argv[2]
s = open(p).read()
old = "This amendment was reviewed as that design's §4, as row 231's was."
assert s.count(old) == 1, s.count(old)
open(p, "w").write(s.replace(old, old + " Oracle pin `" + pin + "`."))
PY
  bash "$RB" 600 "$X/logs/m2-pindup-${pin:0:8}.log" go test -count=1 -timeout 300s -v -run '^TestMarketing(Position|Phone)Parity$' ./scripts/marketingcheck; echo "pin-dup ${pin:0:8} rc=$?" >> "$OUT"
  grep -E -- '--- FAIL' "$X/logs/m2-pindup-${pin:0:8}.log" | sed -E 's/ \([0-9.]+s\)$//' >> "$OUT"
  cp "$X/pristine/orchestration.md" "$D"
  echo "    restore: $(shasum -a 256 "$D" | cut -c1-16)" >> "$OUT"
done
cat "$OUT"
SH
```
Expected:
```
pin-dup 0a9a6795 rc=1
--- FAIL: TestMarketingPositionParity
    --- FAIL: TestMarketingPositionParity/pin-registry
    restore: e66293d11e303218
pin-dup 31576922 rc=1
--- FAIL: TestMarketingPhoneParity
    --- FAIL: TestMarketingPhoneParity/pin-registry
    restore: e66293d11e303218
```

### M2.5 Bank evidence and STOP (do not commit M2)

```bash
bash <<'SH'
source /Users/voightkampff/dev/sunholo-data/.scratch-iter336-executor/env.sh
cd "$WT"
cat "$X/logs/m2.txt" "$X/logs/m2-pindup.txt" > "$BANK/executor-M2-acceptance.txt"
git status --porcelain
SH
```
Expected: ` M design_docs/planned/w-marketing-capture-orchestration.md` and
`?? design_docs/verification/world-iter336/executor-M2-acceptance.txt` (plus untracked files the
controller left). Leave both uncommitted. Report to the controller and end the run.

## 5. Controller-only steps (not the executor)

1. **M2 sign-off (design §5 M2, kimi r2):** diff each pasted A1–A7 block in the orchestration doc
   against design §4 and bank `world-iter336/m2-signoff.txt`. The two deltas `phone 3` → `phone 4`
   (A3, A7) are the measured, design-§2.3-authorised ones (§0); any other delta is a hard stop.
2. **Orchestration quorum:** `ailang design-quorum design_docs/planned/w-marketing-capture-orchestration.md --author claude:claude-opus-5-5`.
   A BLOCKED round with a direction objection means neither M1 nor M2 lands, and the row parks.
3. Commit M2 with the sign-off and quorum artifacts; push, PR, CI (all jobs green), merge, and the
   charter/log/STATUS entries for row 238. M1 and M2 land together or not at all.
4. Optional cross-check: `AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_go.sh` (CI's gate).

## 6. Risks

- **Concurrent change** to `scripts/marketingcheck` or the orchestration doc: M1.0 and M2.1 catch it (STOP).
- **Corpus drift** (`~/.ailang/state/world-iter328-phone-corpus` is host-local): §2 catches it (BLOCKED).
- **Full-suite flake outside marketingcheck**: report, do not re-run for green; the controller decides.
- **Shell traps**: zsh does not word-split unquoted variables and `echo` eats `\b`; every block is a
  `bash` heredoc for that reason.

## 9. Verification Log (planner, iteration 336; files under `design_docs/verification/world-iter336/`)

Scratch: `/Users/voightkampff/dev/sunholo-data/.scratch-iter336-planner/` (`git archive b9f67a27`
copies; GOCACHE/GOMODCACHE there). Go `go1.26.6 darwin/arm64`; `AILANG_BIN=$HOME/.pinned-ailang/ailang` (`AILANG v0.41.0`).

| P | Command | Observed |
|---|---|---|
| P01 | `git apply` designer `designer-r2-patch.diff.txt` on a `git archive HEAD` copy; `bash run_bounded.sh 1500 LOG go test -count=1 -timeout 20m -v -run '^TestMarketing(Position\|Phone)' ./scripts/marketingcheck` (runner sha `1635f8d7…2082`) | rc 0, 0 SKIP; `PARITY position: inputs=51 verdict-diffs=0 reason-diffs=0 ledger=8`; `PARITY phone: inputs=44 diffs=0 ledger=3`; `PARITY phone: inputs=62 diffs=0 ledger=3`; `CHECK images=10 byte-equal=10; real arms=42; parity inputs=62`. Matches design §5 (V12) |
| P02 | Added B7 to `phoneArms` (`phoneFill` luma 250 at x 61, x 139..199, x 260, rows 833..840, mirroring B1) and to `phoneWant` as ACCEPT | arm builds; `TestMarketingPhoneArms/judge/B7-thumb-two-equal-bridged-runs` PASS |
| P03 | Scratch-only probe test (`ScanPhone` on `phoneArmByName(...).Make(base)`), render and real seed 5 (`plan-M1-acceptance.txt` §B7) | B7 render and real: thumb `{6 833 139 841}`, ACCEPT. Bases `{6 833 333 841}`; B1/B1c/B1b `{6 833 333 841}` ACCEPT |
| P04 | Oracle built from `world-iter315/v66-phonecheck-main.go.txt` (sha `31576922…0b4c`), run on the B7 PNGs with the V50 reference | `b7-render … thumb={-1 -1 -1 -1} … -> r3 REFUSE`; `b7-real … -> r3 REFUSE`. So B7 is a REFUSE→ACCEPT parity difference: phone ledger = 4 |
| P05 | B7 + ledger entry + counts (45, 63, 45/15) + design §5 M1 comment updates; acceptance as P01 | rc 0, 0 SKIP; `PARITY position: inputs=51 verdict-diffs=0 reason-diffs=0 ledger=8`; `PARITY phone: inputs=45 diffs=0 ledger=4`; `PARITY phone: inputs=63 diffs=0 ledger=4`; `CHECK images=10 byte-equal=10; real arms=43; parity inputs=63` |
| P06 | `M238-run-rightmost-tie` on the P05 tree, and (control) on the P01 tree without B7 | P05: rc 1, FAIL exactly `TestMarketingPhoneArms`, `…/judge/B7-thumb-two-equal-bridged-runs`, `TestMarketingPhoneParity` (`phone REFUSE ["(iii′) thumb {200 833 333 841}"]`). Control: rc 0, 0 FAIL (survives) |
| P07 | Final diff = scratch git repo diff of the 4 files; `git apply --check` in the worktree; apply to a fresh `git archive` copy; `cmp` with the P05 files | `plan-M1.diff.txt` sha256 `7da36f3006efca03e0b49249d60e577097df799e09b8186ce9525217dc585d74`, 4 files +248/−60; check rc 0; all 4 files identical; post-apply sha256 as M1.1 |
| P08 | 14 mutants on the final tree, each alone under the runner (600 s) (`plan-M1-mutants.txt`) | every mutant rc 1, `match=EXACT` against `plan-mutants-expected.txt`; every restore byte-identical; 43 s total |
| P09 | `go vet ./scripts/marketingcheck`, `gofmt -l`, `go build ./...`, fixture `diff -r` vs base, capture-day steps 2–5 (`plan-M1-acceptance.txt`) | vet/gofmt/build rc 0, gofmt empty, fixtures identical; steps rc 0/1/0/1 with the outputs in M1.4 |
| P10 | `go test ./...` (2400 s) in a `git archive` copy; then in a local `git clone` at `b9f67a27` + `plan-M1.diff.txt` (`plan-M1-acceptance.txt`) | archive copy: rc 1, only `host/pins TestReleaseLiteralsExplainedByManifest` (`git ls-files: exit status 128`; no git checkout), 34 ok. Clone (uncommitted M1, like the executor's M1.5): `go build ./...` rc 0; `go test ./...` **rc 0**, 35 ok, 0 FAIL, 376 s; `ok … scripts/marketingcheck 123.674s` |
| P11 | `./scripts/verify_ail.sh` (1200 s) on the final tree (archive copy; finished 20:49:11Z, under 6.5 min) | rc 0; `✓ verify gate PASSED: 16 required identities verified, 40 named tests pass, 503 se-tools named tests pass` |
| P12 | M2: anchors at the base, `plan-M2.diff.txt` built in a scratch git repo, `git apply --check` in the worktree, acceptance greps, parity, M238-doc-pin-dup × 2 (`plan-M2-acceptance.txt`) | anchors 1×7, controls 6/0; `plan-M2.diff.txt` sha256 `08e9838f84f3eaaa9f7aef6f1e6ad8fbea2fe18450c11c738806a2e2219818e9`, +10/−5; doc sha `0b2ee859…` → `e66293d1…`; greps 1/1/1/1, 0; pins 1/1; parity rc 0 with both `pin-registry` PASS; pin-dup rc 1 / rc 1 reds the matching `pin-registry` subtest; restores byte-identical |
| P13 | `ls verification/world-iter335/`; `head -3 controller-fullgate.txt`; `scripts/verify_go.sh` lines 12–17, 181–198; `ci.yml:108`; `~/.pinned-ailang/ailang --version` | iteration 335 used `AILANG v0.41.0` (its evidence lives at repo-root `verification/world-iter335/`, not under `design_docs/`); the pinned binary reports `AILANG v0.41.0`, `Commit: 24ee108` |
| P14 | Executor dry run: the bash blocks of §2, M1.0–M1.4, M1.6 and M2.1–M2.4 extracted verbatim from this file and run against a scratch `git clone` at `b9f67a27` (only `X` and `WT` substituted) (`plan-M1-acceptance.txt` §P14) | every block printed exactly its Expected output; M1.6: 14 × `rc=1 match=EXACT`, restores `9029dbabd5747610 9540d55a323640ed`, 42 s; M2.4 restores `e66293d11e303218` |
