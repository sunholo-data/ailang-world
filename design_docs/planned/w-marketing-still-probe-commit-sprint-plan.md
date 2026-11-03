# Sprint plan: w-marketing-still-probe-commit (queue row 239, iteration 337)

Status: **planned**. Planner: iteration 337 (claude-opus-5-5). Executor: one run, M1 then M2, following
this file literally.
- **Design (authority):** `design_docs/planned/w-marketing-still-probe-commit.md` r3 (`8ef19b14`).
  Read §2 and §3 before starting. This plan does not re-design anything. Where it differs from the
  design's text, the difference was measured and is listed in §0.
- **Worktree:** `/Users/voightkampff/dev/sunholo-data/.wt-world-iter337-row239`,
  branch `sprint/world-iter337-still-probe-commit`.
- **Executor scratch (never inside the worktree):** `/Users/voightkampff/dev/sunholo-data/.scratch-iter337-executor/`.
- **Planner evidence:** `design_docs/verification/world-iter337/planner-*.txt` (§9 Verification Log).
- **The executor writes no Go code and no doc text by hand.** It applies two banked unified diffs,
  checks their sha256, runs the gates and the mutants, and banks small `.txt` evidence.
- **The executor makes no git write operation other than `git apply`.** No `git add`, `commit`,
  `checkout`, `restore`, `stash`, `reset`, `push`, PR or amend. Everything stays **uncommitted**;
  the controller commits.

| Milestone | What | Files | Commit |
|---|---|---|---|
| M1 | the test (design §2.2, §2.4, §2.6) | `cmd/ailang-worldd/marketing_still_test.go` (new, 142 lines), `cmd/ailang-worldd/marketing_chrome_test.go` (+5) via `plan-M1.diff.txt` | none (controller) |
| M2 | §4.4 item (1), one line (design §3 M2) | `design_docs/planned/w-marketing-capture-orchestration.md` via `plan-M2.diff.txt` | none (controller) |

Order: M1 (§3), then M2 (§4), then the full gates on the final tree (§5), then the state check and
evidence (§6). M2 lands with M1 or not at all.

## 0. What the planner measured that the design did not fix (read this first)

| Item | Design text (r3) | Measured, used by this plan |
|---|---|---|
| Sibling site of the zero-`<script` predicate (header comment, §2.6) | `workbench_still_test.go:52,82,95` | **`:53,82,95`**. Line 52 is `first[i] = body`; the predicate `strings.Count(body, "<script")` is `:53` (its `t.Fatalf` is `:54`, the line V14 reports). `:82` (tag replace) and `:95` (sleep comment; `time.Sleep` is `:98`) are as designed. The header names `render.go:203` and `workbench_still_test.go:53,82,95` (and `95-98` for the spacing) (P03) |
| `TIMING` log (§2.4: optional) | optional | **kept**, with the prototype's `routes := …` local, so mutant MT1's anchor stays unique. One doc comment is added on the test function. Otherwise the file is the prototype's lines 21–126 verbatim, with the PROTO hook and the `os` import dropped (P02) |
| `t.Errorf`/`STILL` line attribution | prototype `:130`, drill `:128` | `checkMarketingStill` calls `t.Helper()`, so every `STILL` and `P0 item (1)` line is reported at the **call site**: `marketing_still_test.go:140` standalone, `marketing_chrome_test.go:128` in the drill (P05, P07) |
| M1 step 5: `git diff --name-only 2f293570 -- . ':!design_docs'` lists the two test files | two files | `git apply` leaves the new file **untracked**, and `git diff` never lists untracked files, so that command prints only `marketing_chrome_test.go`. This plan checks scope with M1.0's `git diff --name-only 2f293570 HEAD -- . ':!design_docs'` = 0 lines (nothing committed outside `design_docs` since the base) and §6's full `git status --porcelain`, whose only non-`design_docs` lines are ` M cmd/ailang-worldd/marketing_chrome_test.go` and `?? cmd/ailang-worldd/marketing_still_test.go` (P12) |
| Full suite at the base | (not stated) | `go test ./...` at `8ef19b14`, both pins: **rc 1**, the only failure `host/verifygate TestMissionConfigDefaultFlowIntegration` (`valid_ledger`, `invalid_ledger`: "default-flow fixture timed out … context deadline exceeded"). That test gives its `verify_go.sh` fixture **10 s**; alone it takes 4.90 s / 5.42 s and passes (rc 0). Under the full suite and host load (load average 11–12, a local llama-server running), it timed out. It is a base-state, load-dependent failure in a package this row does not touch (P09, P10). See §5 for what the executor does with it |
| Drill wall | 20.00 s (V12) | 34.49 s test, 36 s wall (P07); the dry run's whole M2.3 block took 23 s (P13). Same 4 `STILL` lines and PASS both times |

The mutation table (design §2.5) reproduces exactly on the applied file, three times with byte-identical summaries (P06, P13): see §3 M1.5.

## 1. Executor rules (each one is a STOP if broken)

1. **Shell.** The Bash tool's shell is zsh. Run every block in this plan exactly as written, as a
   `bash <<'SH' … SH` heredoc. Do not retype it into zsh line by line.
2. **Runner.** Every command that builds or tests runs under `bash "$RB" <limit_s> "$OUT" <cmd…>`;
   read `rc=$?` on the next line. `$OUT` gets stdout, `$OUT.err` gets stderr. **rc 124 means the
   deadline fired: STOP and report.**
3. **Limits (design §3):** each focused `go test` 600 s (`-timeout 540s` inside); the Chrome drill
   1300 s (`-timeout 1200s` inside); each mutant 300 s (`-timeout 240s` inside); `go vet`, `gofmt -l`,
   `go build ./...` 600 s each; `go test ./...` 2400 s; `./scripts/verify_ail.sh` 1200 s.
4. **Long commands.** The mutant batch (§3 M1.5), `go test ./...` and `verify_ail.sh` can outlast one
   foreground tool call (600 s). Launch them with the Bash tool's own `run_in_background: true` and
   poll with the bounded loop given there. **Never use `nohup`** or any SIGHUP-ignoring wrapper
   (row 231 D-14: `nohup` reds `TestMarketingSupervisorFakeDaemonCleanup/shell-hup`). A background
   launch is not a wait: always poll the status file in the bounded loop until `DONE`.
5. **No pipes before an rc is read.** Redirect to a file, read rc, then grep the file.
6. **Git.** The only git write operation is `git apply` (M1.1, M2.2). Read-only git (`status`,
   `diff`, `rev-parse`, `log`) is fine. Mutants are restored with `cp` from pristine copies and
   checked by sha256 (the banked runner does this); never with `git checkout`.
7. **Do not edit** the design doc, this plan, the sprint JSON, any `plan-*.txt` or `planner-*.txt`,
   `world-iter336/run_bounded.sh.txt`, `design_docs/world-mission*.md`, `tools/launchd/*`,
   `.github/`, `go.mod`, or any `host/` or `scripts/` file (the mutants touch `host/workbench/render.go`
   only through the banked runner, which restores it). Never touch `~/.ailang/state/mission-v1*`.
8. **Do not export `WORLD_CHROME` or `WORLD_MARKETING_RECTS`** in the environment. They are set only
   on the one drill command line (M2.3), so `go test ./...` does not run the drill.
9. **Any expected value below that does not match is a STOP**, except where a step names its own
   handling (§5's known base failure). Report the command, the expected value and the actual output.
   Do not "fix" the test, the patch or the plan.

## 2. Environment (run once, first)

`AILANG_BIN` is `$HOME/.pinned-ailang/ailang` (`AILANG v0.41.0`, the pin `scripts/verify_go.sh`
enforces and CI uses). `WORLD_TOOL_AILANG_BIN` is `$HOME/.pinned-ailang-tools/v0.54.15/ailang`
(`AILANG v0.54.15` = `daemon.ToolBinaryRelease`). The test needs both; `cliSeBins` skips without the
second (design §2.1).

```bash
bash <<'SH'
set -u
X=/Users/voightkampff/dev/sunholo-data/.scratch-iter337-executor
WT=/Users/voightkampff/dev/sunholo-data/.wt-world-iter337-row239
BANK=$WT/design_docs/verification/world-iter337
mkdir -p "$X/logs" "$X/mutlogs" "$X/pristine"
cat > "$X/env.sh" <<EOF
export X=$X WT=$WT BANK=$BANK
export PATH=/opt/homebrew/bin:\$PATH
export AILANG_BIN=\$HOME/.pinned-ailang/ailang
export WORLD_TOOL_AILANG_BIN=\$HOME/.pinned-ailang-tools/v0.54.15/ailang
CHROME="/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
RECTS=$WT/design_docs/verification/world-iter312/r1-rects.json
export RB=$X/run_bounded.sh
EOF
cp "$WT/design_docs/verification/world-iter336/run_bounded.sh.txt" "$X/run_bounded.sh"
cp "$BANK/plan-mutate.py.txt" "$X/mutate.py"
cp "$BANK/plan-run_mutants.sh.txt" "$X/run_mutants.sh"
cp "$BANK/plan-mutants-expected.txt" "$X/mutants-expected.txt"
source "$X/env.sh"
shasum -a 256 "$X/run_bounded.sh" "$X/mutate.py" "$X/run_mutants.sh" "$X/mutants-expected.txt" "$BANK/plan-M1.diff.txt" "$BANK/plan-M2.diff.txt" | sed -E 's|  .*/|  |'
"$AILANG_BIN" --version 2>/dev/null | head -1
"$WORLD_TOOL_AILANG_BIN" --version 2>/dev/null | head -1
"$CHROME" --version
(cd "$WT" && go version)
shasum -a 256 "$RECTS" | cut -c1-16
SH
```

Expected (every line; any difference is a STOP):
```
1635f8d7308e03a6bc02eeacd4bd6eb723ce76405c4c1d28f2c45ef93d8e2082  run_bounded.sh
71fb8f14f9d3ea25cfc64276b1ffeecb93fcdae48320414197ee26ac0c81c304  mutate.py
506ba8b9ee1bcb002633aec233a2823d03637afaf8a61cbf9cb864232fdda2e0  run_mutants.sh
ce3964b200a11d235005602174203d968badfacd18c65bcbb86152d8e7b145ba  mutants-expected.txt
08542390e97ddd9adfde675449a211aa7355155540aac46f0163320aca8d29e0  plan-M1.diff.txt
29d81e47bc8914c5c82d0020361ed46c376cada3a3f662bf339ac5df0be43243  plan-M2.diff.txt
AILANG v0.41.0
AILANG v0.54.15
Google Chrome 156.0.8078.80 
go version go1.26.6 darwin/arm64
ba7c83ce4e31c86a
```
(`Chrome --version` prints a trailing space after the version.) `go version` must run inside the worktree (`GOTOOLCHAIN=auto` selects go1.26.6 from `go.mod`). A
different Chrome version is not a STOP by itself: report it, and the drill (M2.3) decides.

Every later block starts with `source /Users/voightkampff/dev/sunholo-data/.scratch-iter337-executor/env.sh`.

## 3. M1: the test

### M1.0 Base check (no concurrent change since the design's base)

```bash
bash <<'SH'
source /Users/voightkampff/dev/sunholo-data/.scratch-iter337-executor/env.sh
cd "$WT"
git diff --name-only 2f293570 HEAD -- . ':!design_docs' | wc -l | tr -d ' '
git status --porcelain -- cmd host scripts design_docs/planned/w-marketing-capture-orchestration.md | wc -l | tr -d ' '
ls cmd/ailang-worldd/marketing_still_test.go 2>&1 | grep -c 'No such file'
shasum -a 256 cmd/ailang-worldd/marketing_chrome_test.go host/workbench/render.go design_docs/planned/w-marketing-capture-orchestration.md
SH
```
Expected:
```
0
0
1
04ab4310b2cb89779feed8e438a90c340551ce5b83e72e84e3685b3220e28259  cmd/ailang-worldd/marketing_chrome_test.go
d0a108cf829b34348860cd85782c1cfd2db2675c1dedb6d3bbe16033952c7152  host/workbench/render.go
aae1bdcdb5b6c3ba8a17df4a2a6579f0013818077365dbb022e8304b48761a27  design_docs/planned/w-marketing-capture-orchestration.md
```
A difference means a concurrent change landed (design §4): STOP; the mission re-bases and re-plans.

### M1.1 Apply the patch

```bash
bash <<'SH'
source /Users/voightkampff/dev/sunholo-data/.scratch-iter337-executor/env.sh
cd "$WT"
git apply --check "$BANK/plan-M1.diff.txt"; echo "check rc=$?"
git apply "$BANK/plan-M1.diff.txt"; echo "apply rc=$?"
git diff --numstat -- cmd/ailang-worldd
wc -l < cmd/ailang-worldd/marketing_still_test.go | tr -d ' '
shasum -a 256 cmd/ailang-worldd/marketing_still_test.go cmd/ailang-worldd/marketing_chrome_test.go
SH
```
Expected:
```
check rc=0
apply rc=0
5	0	cmd/ailang-worldd/marketing_chrome_test.go
142
e38f4ccac76ae66857af1a597b6f8bc1c8e25ab9ed4f16536347da46058a71ff  cmd/ailang-worldd/marketing_still_test.go
70ef39226db6a71f9647e5a4ce52eb43d01975ff1889eb8c6a01dbc5b4eefb9f  cmd/ailang-worldd/marketing_chrome_test.go
```
(The new file is untracked, so `git diff --numstat` lists only the chrome test.) The patch is the
whole of M1: the new file (prototype lines 21–126, test renamed `TestMarketingStillRoutes`, PROTO hook
and `os` import dropped, the §2.6 header comment) and the five §2.2 lines after the `check` build in
`TestMarketingChromeSixShots`. The chrome hunk is byte-identical to the designer's V12 patch (P04).

### M1.2 Acceptance steps 1 and 2 (600 s each)

```bash
bash <<'SH'
source /Users/voightkampff/dev/sunholo-data/.scratch-iter337-executor/env.sh
cd "$WT"
LOG=$X/logs/m1-step1.log
bash "$RB" 600 "$LOG" go test -count=1 -timeout 540s -v -run '^TestMarketingStillRoutes$' ./cmd/ailang-worldd; echo "step1 rc=$?"
echo "pass=$(grep -c '^--- PASS: TestMarketingStillRoutes' "$LOG")"
echo "skip=$(grep -c -- '--- SKIP' "$LOG")"
echo "still4=$(grep -cE 'STILL .* repeat_equal=true script_tags=0 equals_default_minus_tag=true default_live_tags=1 ' "$LOG")"
echo "p0err=$(grep -c 'P0 item (1)' "$LOG")"
LOG2=$X/logs/m1-step2.log
env -u WORLD_TOOL_AILANG_BIN bash "$RB" 600 "$LOG2" go test -count=1 -timeout 540s -v -run '^TestMarketingStillRoutes$' ./cmd/ailang-worldd; echo "step2 rc=$?"
echo "skip2=$(grep -c -- '--- SKIP: TestMarketingStillRoutes' "$LOG2") toolset=$(grep -c 'WORLD_TOOL_AILANG_BIN set=false' "$LOG2")"
SH
```
Expected (P05):
```
step1 rc=0
pass=1
skip=0
still4=4
p0err=0
step2 rc=0
skip2=1 toolset=1
```
Step 2 is the CI-shaped control (design §2.1, V09): `AILANG_BIN` only, so `cliSeBins` skips.

### M1.3 gofmt, vet, build (600 s each)

```bash
bash <<'SH'
source /Users/voightkampff/dev/sunholo-data/.scratch-iter337-executor/env.sh
cd "$WT"
bash "$RB" 600 "$X/logs/gofmt.txt" gofmt -l cmd/ailang-worldd; echo "gofmt rc=$? lines=$(wc -l < "$X/logs/gofmt.txt" | tr -d ' ')"
bash "$RB" 600 "$X/logs/vet.txt" go vet ./cmd/ailang-worldd; echo "vet rc=$?"
bash "$RB" 600 "$X/logs/build.txt" go build ./...; echo "build rc=$?"
SH
```
Expected: `gofmt rc=0 lines=0`, `vet rc=0`, `build rc=0` (P05).

### M1.4 Pristine record for the mutants

The banked runner snapshots `host/workbench/render.go` and `cmd/ailang-worldd/marketing_still_test.go`
into `$X/pristine/` at its start and restores both by `cp` after every run, checking the sha256.
Nothing to do here except to confirm the two files are the M1.0 / M1.1 values:
```bash
bash <<'SH'
source /Users/voightkampff/dev/sunholo-data/.scratch-iter337-executor/env.sh
cd "$WT"; shasum -a 256 host/workbench/render.go cmd/ailang-worldd/marketing_still_test.go | cut -c1-16
SH
```
Expected: `d0a108cf829b3434` and `e38f4ccac76ae668`.

### M1.5 Mutants (design §2.5; 11 runs, 300 s each, in the worktree, restored by `cp`)

`mutate.py` makes one exact string replacement per mutant (two for MS4) and exits 2 with the file
untouched if an anchor does not occur exactly once. `run_mutants.sh` runs, in order: `healthy`,
`MS1`…`MS5`, `MT1`, `MT2`, `MT3+MS3`, `MT4+MS4`, `MT5+MS5`; each under the runner with
`go test -count=1 -timeout 240s -v -run '^TestMarketingStillRoutes$' ./cmd/ailang-worldd`. It writes a
normalized summary (`sha256:<h>` for route hashes, timings dropped, `LC_ALL=C` sort) whose bytes must
equal `plan-mutants-expected.txt`. A **build failure is not a kill**: the runner prints
`!! BUILD FAILED: not a kill`, and that line is not in the expected file, so `cmp` fails.

Launch with `run_in_background: true`:
```bash
bash <<'SH'
source /Users/voightkampff/dev/sunholo-data/.scratch-iter337-executor/env.sh
rm -f "$X/logs/mutants.status"
bash "$X/run_mutants.sh" "$WT" "$X" "$RB" "$X/mutate.py" > "$X/logs/mutants-stdout.txt" 2>&1; echo "run_mutants rc=$?" > "$X/logs/mutants.status"
echo DONE >> "$X/logs/mutants.status"
SH
```
Poll (foreground, Bash tool timeout 600000 ms; repeat until `DONE`; planner wall about 160 s):
```bash
bash <<'SH'
F=/Users/voightkampff/dev/sunholo-data/.scratch-iter337-executor/logs/mutants.status
end=$(( $(date +%s) + 540 ))
while [ "$(date +%s)" -lt "$end" ]; do grep -q '^DONE' "$F" 2>/dev/null && break; sleep 10; done
cat "$F"
SH
```
Then:
```bash
bash <<'SH'
source /Users/voightkampff/dev/sunholo-data/.scratch-iter337-executor/env.sh
cd "$WT"
cmp "$X/logs/mutants-summary.txt" "$X/mutants-expected.txt"; echo "mutants cmp rc=$?"
grep -c 'restore=identical' "$X/logs/mutants-summary.txt"
grep -E '^== ' "$X/logs/mutants-summary.txt"
grep -E 'NOT-APPLIED|ANCHOR|unknown' "$X/logs/mutants-raw.txt" | wc -l | tr -d ' '
shasum -a 256 host/workbench/render.go cmd/ailang-worldd/marketing_still_test.go | cut -c1-16
SH
```
Expected:
```
mutants cmp rc=0
11
== healthy rc=0 restore=identical
== MS1 rc=1 restore=identical
== MS2 rc=1 restore=identical
== MS3 rc=1 restore=identical
== MS4 rc=1 restore=identical
== MS5 rc=1 restore=identical
== MT1 rc=1 restore=identical
== MT2 rc=1 restore=identical
== MT3+MS3 rc=0 restore=identical
== MT4+MS4 rc=0 restore=identical
== MT5+MS5 rc=0 restore=identical
0
d0a108cf829b3434
e38f4ccac76ae668
```
On `cmp rc=1`, run `diff "$X/mutants-expected.txt" "$X/logs/mutants-summary.txt"` and STOP with it.

What each run must show (the full lines are in `plan-mutants-expected.txt`; every line is
`marketing_still_test.go:140: P0 item (1) …`, once per route, 4 routes):

| Id | Mutation | Must fail with | rc |
|---|---|---|---|
| healthy | none | `--- PASS` | 0 |
| MS1 | `render.go:203` → bare live.js tag | `script_tags=1, want 0` ×4 and `still HTML is not the default minus exactly the live.js tag` ×4 (A2+A3) | 1 |
| MS2 | append `<script type="application/json">{}</script>` after the `:203` construct | `script_tags=1, want 0` ×4 (A2 only) | 1 |
| MS3 | `{{if .Still}}<meta name="still" content="1">{{else}}…` | `still HTML is not the default minus exactly the live.js tag` ×4 (A3 only) | 1 |
| MS4 | `Render` appends `" "` to `p.Title` on still pages once 1 s has passed since the first still render | `repeat_equal=false (still HTML changed across 2 s)` ×4 (A1 only) | 1 |
| MS5 | `:203` construct → empty | `control …: default HTML has 0 live.js tags, want 1` ×4 (A0 only) | 1 |
| MT1 | `routes := marketingStillRoutes(…)[:0]` | `P0 item (1): want 4 distinct capture routes, got 0: []` ×1 (A4) | 1 |
| MT2 | `if !repeat {` → `if repeat {` | `repeat_equal=false …` ×4 on healthy code (A1) | 1 |
| MT3+MS3 | A3 disabled + MS3 | `--- PASS` (A3 is MS3's sole killer) | 0 |
| MT4+MS4 | the 2 s sleep removed + MS4 | `--- PASS` (the sleep is load-bearing) | 0 |
| MT5+MS5 | A0 disabled + MS5 | `--- PASS` (A0 is MS5's sole killer) | 0 |

MI1 (delete the §2.2 call from the drill) stays derived, as in the design: the M2 pass rule on the
MS1 log gives 0 (M2.3 control).

## 4. M2: §4.4 item (1) text

### M2.1 Pre-paste anchor check (hard stop; design §3 M2)

```bash
bash <<'SH'
source /Users/voightkampff/dev/sunholo-data/.scratch-iter337-executor/env.sh
cd "$WT"; D=design_docs/planned/w-marketing-capture-orchestration.md
printf '%s\t%s\n' "$(grep -cF -- '  1. **Route-level still determinism.**' "$D")" 'anchor'
printf '%s\t%s\n' "$(grep -cF -- '  1. **Route-level still determinsm.**' "$D")" 'control misspelled'
for a in checkMarketingStill zz_iter305_probe_test.go TestIter305SeededStillCaptureProbe; do printf '%s\t%s\n' "$(grep -cF -- "$a" "$D")" "$a"; done
shasum -a 256 "$D"
SH
```
Expected: `1 anchor`, `0 control misspelled`, `0 checkMarketingStill`, `2 zz_iter305_probe_test.go`,
`2 TestIter305SeededStillCaptureProbe`, and sha256
`aae1bdcdb5b6c3ba8a17df4a2a6579f0013818077365dbb022e8304b48761a27`. Anything else: STOP. Never paste
by line number.

### M2.2 Apply

```bash
bash <<'SH'
source /Users/voightkampff/dev/sunholo-data/.scratch-iter337-executor/env.sh
cd "$WT"; D=design_docs/planned/w-marketing-capture-orchestration.md
git apply --check "$BANK/plan-M2.diff.txt"; echo "check rc=$?"
git apply "$BANK/plan-M2.diff.txt"; echo "apply rc=$?"
git diff --numstat -- "$D"
shasum -a 256 "$D"
SH
```
Expected: `check rc=0`, `apply rc=0`, `1	1	design_docs/planned/w-marketing-capture-orchestration.md`,
sha256 `50d8ef34143e930cdf188aaaa27bf4baef3c08c8542c79c844db3785ec6507ee`.

The patch replaces the one anchor line (221) with the design's fenced line, byte-exact (P11:
`cmp` of the pasted line 221 with design line 226 is equal). Items (2)–(4), the rest of §4.4, §7 and
§9/§10 are untouched.

### M2.3 Acceptance, including the drill (1300 s)

```bash
bash <<'SH'
source /Users/voightkampff/dev/sunholo-data/.scratch-iter337-executor/env.sh
cd "$WT"; D=design_docs/planned/w-marketing-capture-orchestration.md; OUT=$X/logs/m2.txt
git diff --numstat -- "$D" > "$OUT"
for a in checkMarketingStill zz_iter305_probe_test.go TestIter305SeededStillCaptureProbe '  1. **Route-level still determinism.**' '  2. **Six-shot two-pass.**'; do printf '%s\t%s\n' "$(grep -cF -- "$a" "$D")" "$a" >> "$OUT"; done
bash "$RB" 600 "$X/logs/m2-parity.log" go test -count=1 -timeout 300s -v -run '^TestMarketing(Position|Phone)Parity$' ./scripts/marketingcheck; echo "parity rc=$?" >> "$OUT"
grep -E -- '--- (PASS|FAIL)' "$X/logs/m2-parity.log" | sed -E 's/ \([0-9.]+s\)$//' >> "$OUT"
LOG=$X/logs/m2-drill.log
WORLD_CHROME="$CHROME" WORLD_MARKETING_RECTS="$RECTS" bash "$RB" 1300 "$LOG" go test -count=1 -timeout 1200s -v -run '^TestMarketingChromeSixShots$' ./cmd/ailang-worldd; echo "drill rc=$?" >> "$OUT"
echo "drill pass=$(grep -c '^--- PASS: TestMarketingChromeSixShots' "$LOG") still=$(grep -cE 'STILL .* repeat_equal=true script_tags=0 equals_default_minus_tag=true default_live_tags=1 ' "$LOG") p0err=$(grep -c 'P0 item (1)' "$LOG")" >> "$OUT"
echo "control MS1 log still=$(grep -cE 'STILL .* repeat_equal=true script_tags=0 equals_default_minus_tag=true default_live_tags=1 ' "$X/mutlogs/MS1.log")" >> "$OUT"
cat "$OUT"
SH
```
Expected (P07, P08, P13):
```
1	1	design_docs/planned/w-marketing-capture-orchestration.md
1	checkMarketingStill
1	zz_iter305_probe_test.go
1	TestIter305SeededStillCaptureProbe
1	  1. **Route-level still determinism.**
1	  2. **Six-shot two-pass.**
parity rc=0
--- PASS: TestMarketingPhoneParity
    --- PASS: TestMarketingPhoneParity/pin-registry
--- PASS: TestMarketingPositionParity
    --- PASS: TestMarketingPositionParity/pin-registry
drill rc=0
drill pass=1 still=4 p0err=0
control MS1 log still=0
```
The drill line is the pasted M2 command itself (`-timeout 1200s`, `-run '^TestMarketingChromeSixShots$'`),
with V12's rig env. Its 4 `STILL` lines come from `marketing_chrome_test.go:128`. A drill failure is
a STOP: report the log's `--- FAIL`/`P0 item (1)`/`Fatal` lines and `$LOG.err`.

## 5. Full gates on the final tree (M1 + M2): `go test ./...` (2400 s) and `verify_ail.sh` (1200 s)

Design §3 M1 step 4. They run once, after M2, so they cover the final tree. Launch with
`run_in_background: true` (not `nohup`, not `&`):
```bash
bash <<'SH'
source /Users/voightkampff/dev/sunholo-data/.scratch-iter337-executor/env.sh
cd "$WT"
F=$X/logs/fullgate.status
echo "start $(date -u +%FT%TZ) AILANG_BIN=$("$AILANG_BIN" --version 2>/dev/null | head -1) TOOL=$("$WORLD_TOOL_AILANG_BIN" --version 2>/dev/null | head -1) WORLD_CHROME=${WORLD_CHROME:-unset}" > "$F"
bash "$RB" 600 "$X/logs/gobuild-all.txt" go build ./...; echo "go build ./... rc=$?" >> "$F"
bash "$RB" 2400 "$X/logs/gotest-all.txt" go test ./...; echo "go test ./... rc=$?" >> "$F"
bash "$RB" 1200 "$X/logs/verify_ail.txt" ./scripts/verify_ail.sh; echo "verify_ail.sh rc=$?" >> "$F"
echo "DONE $(date -u +%FT%TZ)" >> "$F"
SH
```
Poll (foreground, Bash tool timeout 600000 ms; repeat until `DONE`; the gates' own deadlines bound
the total at 4200 s; planner wall about 430 s (`go test ./...` 386 s, `verify_ail.sh` 39 s)):
```bash
bash <<'SH'
F=/Users/voightkampff/dev/sunholo-data/.scratch-iter337-executor/logs/fullgate.status
end=$(( $(date +%s) + 540 ))
while [ "$(date +%s)" -lt "$end" ]; do grep -q '^DONE' "$F" 2>/dev/null && break; sleep 15; done
cat "$F"
SH
```
Then:
```bash
bash <<'SH'
source /Users/voightkampff/dev/sunholo-data/.scratch-iter337-executor/env.sh
L=$X/logs/gotest-all.txt
echo "ok=$(grep -c '^ok' "$L") notest=$(grep -c 'no test files' "$L") FAILpkg=$(grep -c '^FAIL	' "$L") failtests=$(grep -c -- '--- FAIL' "$L")"
grep -E -- '^FAIL	|--- FAIL' "$L"
grep -F 'cmd/ailang-worldd' "$L"
tail -1 "$X/logs/verify_ail.txt"
SH
```
Expected: the status file shows `AILANG v0.41.0`, `AILANG v0.54.15`, `WORLD_CHROME=unset`,
`go build ./... rc=0`, `go test ./... rc=0`, `verify_ail.sh rc=0`; `ok=35 notest=2 FAILpkg=0 failtests=0`;
an `ok … cmd/ailang-worldd` line; the last `verify_ail.txt` line
`✓ verify gate PASSED: 16 required identities verified, 40 named tests pass, 503 se-tools named tests pass`.

**The one known base failure.** If `go test ./...` is rc 1 and the **only** failures are
`host/verifygate` `TestMissionConfigDefaultFlowIntegration` (and its `valid_ledger`/`invalid_ledger`
subtests) with `default-flow fixture timed out`, that is the base-state load failure (§0, P09/P10).
Then run this one diagnostic, once, and report both results; do not re-run the suite for a green.
The controller decides.
```bash
bash <<'SH'
source /Users/voightkampff/dev/sunholo-data/.scratch-iter337-executor/env.sh
cd "$WT"
bash "$RB" 600 "$X/logs/verifygate-isolated.txt" go test -count=1 -v -run '^TestMissionConfigDefaultFlowIntegration$' ./host/verifygate; echo "isolated rc=$?"
grep -E -- '^--- |^    --- |^(ok|FAIL)' "$X/logs/verifygate-isolated.txt"
uptime
SH
```
Any other failing package or test is a STOP: report the package, the test and its output.

## 6. State check and evidence (no commit)

```bash
bash <<'SH'
source /Users/voightkampff/dev/sunholo-data/.scratch-iter337-executor/env.sh
cd "$WT"
{ echo "# executor M1 acceptance, $(date -u +%FT%TZ), HEAD $(git rev-parse --short HEAD) + plan-M1.diff.txt"
  grep -E 'STILL |TIMING|^--- |^(ok|FAIL)' "$X/logs/m1-step1.log" | sed -E 's/html_sha256=[0-9a-f]{64}/html_sha256=<h>/; s/sha256:[0-9a-f]{64}/sha256:<h>/g'
  echo "step 2 (AILANG_BIN only):"; grep -E -- '--- (PASS|FAIL|SKIP)|set=' "$X/logs/m1-step2.log"
  echo "gofmt/vet/build:"; cat "$X/logs/gofmt.txt" "$X/logs/vet.txt" "$X/logs/vet.txt.err" "$X/logs/build.txt"
} > "$BANK/executor-M1-acceptance.txt"
{ cat "$X/logs/mutants-summary.txt"; echo "# raw"; cat "$X/logs/mutants-raw.txt"; } > "$BANK/executor-M1-mutants.txt"
{ cat "$X/logs/m2.txt"; echo "# drill"; grep -E 'STILL |^--- |^(ok|FAIL)' "$X/logs/m2-drill.log" | sed -E 's/html_sha256=[0-9a-f]{64}/html_sha256=<h>/; s/sha256:[0-9a-f]{64}/sha256:<h>/g'; } > "$BANK/executor-M2-acceptance.txt"
{ cat "$X/logs/fullgate.status"; grep -v '^ok\|no test files' "$X/logs/gotest-all.txt"; tail -1 "$X/logs/verify_ail.txt"; [ -f "$X/logs/verifygate-isolated.txt" ] && { echo "# verifygate isolated"; grep -E -- '^--- |^    --- |^(ok|FAIL)' "$X/logs/verifygate-isolated.txt"; }; } > "$BANK/executor-fullgate.txt"
wc -c "$BANK"/executor-*.txt | tail -1
git status --porcelain
shasum -a 256 cmd/ailang-worldd/marketing_still_test.go cmd/ailang-worldd/marketing_chrome_test.go host/workbench/render.go design_docs/planned/w-marketing-capture-orchestration.md | cut -c1-16
SH
```
Expected: the evidence total is under 60 kB; `git status --porcelain` is exactly (plus any untracked
paths the controller left, e.g. `.claude/`, which you leave alone):
```
 M cmd/ailang-worldd/marketing_chrome_test.go
 M design_docs/planned/w-marketing-capture-orchestration.md
?? cmd/ailang-worldd/marketing_still_test.go
?? design_docs/verification/world-iter337/executor-M1-acceptance.txt
?? design_docs/verification/world-iter337/executor-M1-mutants.txt
?? design_docs/verification/world-iter337/executor-M2-acceptance.txt
?? design_docs/verification/world-iter337/executor-fullgate.txt
```
and the sha256 prefixes `e38f4ccac76ae668`, `70ef39226db6a71f`, `d0a108cf829b3434` (render.go
unchanged), `50d8ef34143e930c`. **Do not commit.** Report to the controller: every rc and count from
§3–§5, the mutant `cmp` result, the drill result, and the full-gate counts; then end the run.

## 7. Controller-only steps (not the executor)

1. Review and commit M1 + M2 together with the `executor-*.txt` evidence; push, PR, CI (all jobs
   green), merge; charter/log/STATUS entries for row 239.
2. Whether the orchestration quorum reviews the M2 line is the controller's call (design §3 M2, row 238 V27).
3. Optional cross-check: `AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_go.sh` (CI's gate).

## 8. Risks and UNMEASURED

- **Concurrent change** to the three touched files or `render.go`: M1.0 and M2.1 catch it (STOP).
- **R1 (design), UNMEASURED.** GitHub CI does not export `WORLD_TOOL_AILANG_BIN`, so
  `TestMarketingStillRoutes` SKIPs there (M1.2 step 2 is that shape). Whether `demo seed` runs on the
  Linux runner is not measured.
- **Flake rate, UNMEASURED.** The test ran healthy 5× standalone (P05, P06 ×2, P13 ×2)
  and the drill 2× (P07, P13); with the designer's runs that is still not a rate. MS4/MT4 depend on wall-clock spacing (1 s mutant
  threshold against a 2 s sleep).
- **`host/verifygate` load failure** at the base (§0, §5): a 10 s fixture deadline under full-suite
  load; unrelated to this row; its rate is not measured.
- **Shell traps**: zsh does not word-split unquoted variables and `echo` eats `\b`; every block is a
  `bash` heredoc for that reason.

## 9. Verification Log (planner, iteration 337; files under `design_docs/verification/world-iter337/`)

Planner scratch: `/Users/voightkampff/dev/sunholo-data/.scratch-iter337-planner/` (a `git clone` at
`8ef19b14` for authoring; logs). Every proof below ran **in the worktree** with the patches applied,
then the worktree was returned to `8ef19b14` (P14). Go `go1.26.6 darwin/arm64`; `AILANG v0.41.0`
(`Commit: 24ee108`); tool `AILANG v0.54.15`; Chrome `156.0.8078.80`; runner sha `1635f8d7…2082`.

| P | Command | Observed |
|---|---|---|
| P01 | `git rev-parse HEAD`; `git diff --name-only 2f293570 HEAD`; `go version` (in the worktree); both `ailang --version`; `Chrome --version` | HEAD `8ef19b14`; since the design's base `2f293570` only `design_docs/` changed (3 design commits); `go1.26.6 darwin/arm64`; `AILANG v0.41.0` / `AILANG v0.54.15`; `Google Chrome 156.0.8078.80` |
| P02 | Authoring in a scratch `git clone` at `8ef19b14`: the new file = §2.6 header + prototype lines 21–114 verbatim + the test (prototype 115–126 renamed, with a doc comment; hook lines 127–129 and the `os` import dropped; `TIMING` kept); `diff` against the prototype; `gofmt -l cmd/ailang-worldd` | `diff` shows only the rename, the doc comment and the dropped hook; `gofmt -l` empty, rc 0; 142 lines |
| P03 | `grep -n '<script\|time.Sleep\|Sleep after' host/daemon/workbench_still_test.go`; `sed -n 203p host/workbench/render.go` | `53` `strings.Count(body, "<script")`, `82` the tag `strings.Replace`, `95` the sleep comment, `98` `time.Sleep(2 * time.Second)`; `render.go:203` is the `{{if .Still}}{{else}}<script src="/workbench/live.js" defer></script>{{end}}` line. So the design's `:52` is `:53` (§0) |
| P04 | The chrome insertion by an anchored python edit (anchor count asserted = 1); `git add -N` + `git diff` in the clone → `plan-M1.diff.txt`; `diff` of its chrome hunk against the designer's V12 hunk (`designer-sixshots-with-still.txt`); M2: the anchor line replaced by design line 226 (count asserted = 1) → `plan-M2.diff.txt` | `plan-M1.diff.txt` sha256 `08542390e97ddd9adfde675449a211aa7355155540aac46f0163320aca8d29e0`, 2 files, +147; chrome hunk `diff` rc 0 (identical to V12); `plan-M2.diff.txt` sha256 `29d81e47bc8914c5c82d0020361ed46c376cada3a3f662bf339ac5df0be43243`, numstat `1 1`, replaced line 221 |
| P05 | In the worktree: `git apply --check` + `git apply` of both patches; M1 steps 1, 2; `gofmt -l`, `go vet ./cmd/ailang-worldd`, `go build ./...` (`planner-M1-acceptance.txt`) | check/apply rc 0 ×4; post sha256 `e38f4cca…` / `70ef3922…` / `50d8ef34…`; step 1 rc 0, `--- PASS: TestMarketingStillRoutes (12.23s)`, 0 SKIP, 4 `STILL … repeat_equal=true script_tags=0 equals_default_minus_tag=true default_live_tags=1` at `marketing_still_test.go:140`, `TIMING seed=8.78s build=0.50s still=2.91s total=12.19s`; step 2 rc 0, `--- SKIP` with `WORLD_TOOL_AILANG_BIN set=false`; gofmt/vet/build rc 0, no output |
| P06 | `run_mutants.sh` + `mutate.py` on the applied worktree file, twice (run 2 with `LC_ALL=C` pinned) (`planner-mutants.txt`) | the §3 M1.5 table exactly: healthy rc 0; MS1 A2+A3 ×4; MS2 A2 ×4; MS3 A3 ×4; MS4 A1 ×4; MS5 A0 ×4; MT1 A4 `got 0: []`; MT2 A1 ×4; MT3+MS3, MT4+MS4, MT5+MS5 rc 0 PASS; no build failure; 11/11 restores identical; run 1 vs run 2 `cmp` rc 0; run 2's summary is `plan-mutants-expected.txt` (sha256 `ce3964b2…`) |
| P07 | Drill (the pasted M2 command) with `WORLD_CHROME`, `WORLD_MARKETING_RECTS=world-iter312/r1-rects.json`, both pins, runner 1300 s | rc 0, wall 36 s; `--- PASS: TestMarketingChromeSixShots (34.49s)`; 4 `STILL` pass lines at `marketing_chrome_test.go:128`; 0 `P0 item (1)` |
| P08 | M2 acceptance in the worktree (`planner-M2-acceptance.txt`) | numstat `1 1`; `grep -cF` 1 for `checkMarketingStill`, `zz_iter305_probe_test.go`, `TestIter305SeededStillCaptureProbe`, both item anchors; parity rc 0, both `pin-registry` PASS; pass rule on the P07 log: `ok` 1, PASS 1, STILL 4; on the MS1 log: 0 |
| P09 | Baseline on the **unpatched** worktree (porcelain 0), both pins: `go build ./...` (600 s), `go test ./...` (2400 s), `./scripts/verify_ail.sh` (1200 s) (`planner-baseline.txt`) | build rc 0; `go test ./...` **rc 1**, 542 s: ok 34, no-test 2, FAIL 1 package (`host/verifygate`, 3 `--- FAIL` lines: `TestMissionConfigDefaultFlowIntegration` and its 2 subtests, `default-flow fixture timed out`); `verify_ail.sh` rc 0, 54 s, `✓ verify gate PASSED: 16 required identities verified, 40 named tests pass, 503 se-tools named tests pass` |
| P10 | The P09 failure alone, once, same tree: `go test -count=1 -v -run '^TestMissionConfigDefaultFlowIntegration$' ./host/verifygate` (600 s); `runMissionCommand` source; `uptime`; `ps` | rc 0, 11 s (`valid_ledger` 4.90 s, `invalid_ledger` 5.42 s); the fixture deadline is `context.WithTimeout(…, 10*time.Second)` at `mission_config_gate_test.go:38`; load averages 11–12; a local llama-server at ~80% CPU |
| P11 | Full gate on the patched worktree (M1 + M2), both pins, after the P06 mutants (`planner-M1-acceptance.txt`) | build rc 0; `go test ./...` **rc 0**, 386 s: ok 35, no-test 2, FAIL 0; `ok … cmd/ailang-worldd 341.094s`, `ok … host/verifygate 365.419s`; `verify_ail.sh` rc 0, 39 s, same PASSED line as P09 |
| P12 | With M1 applied: `git diff --name-only 2f293570 -- . ':!design_docs'`; `git status --porcelain -- . ':!design_docs'`; `git diff --name-only 2f293570 HEAD -- . ':!design_docs' \| wc -l` | only `cmd/ailang-worldd/marketing_chrome_test.go` (the untracked new file is not listed); ` M …marketing_chrome_test.go` and `?? …marketing_still_test.go`; `0` |
| P13 | Executor dry run: this file's bash blocks extracted verbatim and run in the worktree in order (§2, M1.0–M1.5, M2.1–M2.3, §6), only the scratch path replaced (`planner-dryrun.txt`). §5 not dry-run (P09/P11 measured it directly) | every block printed its Expected output (the §6 block's three missing-§5-log errors aside); mutant summary `cmp` rc 0 (third identical run), 96 s for M1.5; drill rc 0, PASS, STILL 4, MS1 control 0; §6 `git status` and sha prefixes as listed |
| P14 | Restore: remove the new file and the dry run's `executor-*.txt`; `git checkout --` the two patched files; sha256 | worktree back at `8ef19b14` (`04ab4310…`, `d0a108cf…`, `aae1bdcd…`); only the planner's `plan-*`/`planner-*` files, this plan and the sprint JSON are added, in the planner commit |
