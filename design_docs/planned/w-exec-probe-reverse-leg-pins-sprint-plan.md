# Sprint plan — `w-exec-probe-reverse-leg-pins` (iteration 329, row 230, clause 3)

**Design doc:** [`design_docs/planned/w-exec-probe-reverse-leg-pins.md`](w-exec-probe-reverse-leg-pins.md)
(r2, `762d8cb4`: quorum r1 fixes §11, the controller's narrow-refinement carve-out after quorum r2 §12).
**Parent:** [`design_docs/implemented/w-exec-probe-reverse-leg-port-race.md`](../implemented/w-exec-probe-reverse-leg-port-race.md) (row 229).
**JSON:** [`sprint_w-exec-probe-reverse-leg-pins.json`](sprint_w-exec-probe-reverse-leg-pins.json).
**Sprint base:** worktree `/Users/voightkampff/dev/sunholo-data/.wt-world-iter329-reverse-leg-pins`, branch
`sprint/world-iter329-exec-probe-reverse-leg-pins`, HEAD `762d8cb4`. Its Go tree equals the design's
measurement base `e5c1867f` (only `design_docs/` changed since; the three files this sprint
changes hash to the base values in §2).
**Planner:** mission-control iteration 329 (planner role). **Executor:** Claude Haiku.

**How this sprint works.** The planner wrote the FINAL contents of every file the sprint changes, in
a throwaway `git worktree add --detach` copy of `762d8cb4` (a sibling of the repo, never `/tmp`,
removed afterwards). It assembled each milestone cumulatively, ran every acceptance command and
every §5 mutant, and banked the validated files as text under
`design_docs/verification/world-iter329/plan-*.txt`. **The executor writes no Go and no YAML.** For
each milestone it copies the banked files into place with `cp`, checks them with
`shasum -a 256 -c`, runs the acceptance commands, runs the banked mutation runner, saves the
output, and commits. Two milestone commits, then one evidence commit.
**Copy commands from fenced blocks, never from table cells** (in a markdown table `\|` renders as `|`).

**Estimate:** about 0.1 executor day. Wall time is dominated by two full `host/broker` runs
(about 3.5 min each) and AC-STABLE (about 1.5 min).

---

## 0. The one-paragraph version

Two mutants survive row 229's suite: the 500 ms hello wait can be set to 5 ms (judge O8), and the
child listener's `c.resume()` can be deleted (judge O10), with every test green. **M1** pins the
wait from both sides with a deterministic test, T10 `TestExecProbeReverseHelloWaitBounds`: a fake
`net.Conn` whose answer "arrives" 250 ms after the hello (floor), and a deadline that must fall
inside the child's 1000 ms hold less one 50 ms poll period (ceiling). It also names the poll
period (`execProbePollEvery`, P1) and makes the SILENT refusal text print the wait it used (P2),
with T7 extended to check that. **M2** adds T11 `TestExecProbeChildListenerDrainsTheHello`, a
structural (S6-labelled) pin of the `c.resume()` handler, and comment-only P3/P4: the wait's two
bounds, what `c.resume()` does and does not do, and the `net6`-before-Run-A ordering invariant.
Both new tests join CI step 8's PASS-guarded list.

| Milestone | Production change | Test change | CI | Mutants it owns |
|---|---|---|---|---|
| **M1** | P1 `execProbePollEvery` + its use; P2 SILENT text from `execProbeHelloWait`; the `pollLoopback` doc line (+6/−3) | `helloConn`, `firstDialReturns`, `execProbeHelloFloor`, T10, T7 extension (+87) | `:247` gains T10 | K1, K2 (+K2c controls), K3, K4, K5, K7, K9; controls K8 and R-M1 |
| **M2** | P3 and P4, comments only (+29/−1) | T11 (+15) | `:247` gains T11 | K6 (and every M1 mutant re-run at the final state) |

Cumulative: `exec_probe.go` +35/−4, `exec_probe_test.go` +102/−0, `ci.yml` ±1. No `.ail`, no
QUICKSTART edit (the user-visible text is byte-identical at 500 ms; design §7, V21).

---

## 1. Environment (on EVERY command)

```sh
cd /Users/voightkampff/dev/sunholo-data/.wt-world-iter329-reverse-leg-pins
export PATH=/opt/homebrew/bin:$PATH
export AILANG_BIN=$HOME/.pinned-ailang/ailang     # prints "AILANG v0.41.0" (measured)
BANK=design_docs/verification/world-iter329
SCR=/Users/voightkampff/dev/sunholo-data/.exec-w329   # executor scratch for BIG logs: a SIBLING of the repo, never /tmp, never inside the worktree
mkdir -p "$SCR"
```

| Pin | Measured this session | Why it matters |
|---|---|---|
| `go` (in the worktree) | `go1.26.6 darwin/arm64` | the CI pin |
| `AILANG_BIN` | `AILANG v0.41.0` | without it the daemon e2e tests SKIP and `verify_ail.sh` cannot run |
| `WORLD_EXEC_SRT_NODE_MODULES` | **unset** | the real-srt tests SKIP locally (30 SKIP in the package). Do NOT set it |
| `actionlint` | `/opt/homebrew/bin/actionlint` | G8 |

**Rules (each one is load-bearing).**
1. **Shell.** The shell is zsh. Run any loop or multi-line script with `bash -c '…'` or a `.sh`
   file. Do not `echo ======` (zsh expands `=word`); quote it.
2. **Instrument rule.** Never write `cmd | tail; echo rc=$?`. Redirect to a file, put
   `echo "rc=$?"` straight after `;`, then read the file. A `-run` that matches nothing exits 0, so
   always count the `--- PASS` lines you expect.
3. **Compile fence.** `go build ./...` does NOT compile `_test.go` files. The fence is
   `go test -run '^$' ./host/broker/` (G3).
4. **Copy, then verify.** Every banked file is copied with `cp` and checked with
   `shasum -a 256 -c "$BANK/plan-M<k>-SHA256SUMS.txt"`. Any line other than `OK` is a STOP: report
   it, do not edit anything to make it match. Never hand-edit a Go or YAML file; never run
   `gofmt -w`.
5. **Base check.** Before M1, `shasum -a 256 -c "$BANK/plan-base-SHA256SUMS.txt"` must print three
   `OK` lines. A mismatch is a STOP (the base moved).
6. **Mutants.** Only through `bash "$BANK/plan-run-mutants.sh.txt" M1|M2`. It refuses to start
   unless `exec_probe.go` equals the milestone's banked copy, restores the file by `cp` from that
   copy after every mutant, and exits 3 on any restore mismatch. Never apply a mutant by hand.
7. **Git.** One commit per milestone, then one evidence commit (§6). Stage with `git add` of the
   named files only (never `git add -A` or `git add .`). Do not push, do not open a PR, do not
   amend, do not rebase, do not `git checkout -- <file>`, do not stash.
8. **Do not edit** the design doc, this plan, the JSON, any `plan-*.txt`, `docs/QUICKSTART.md`,
   `design_docs/world-mission*.md`, `tools/launchd/*`, any skill, `go.mod`. Do not touch the main
   checkout `/Users/voightkampff/dev/sunholo-data/ailang-world` or `~/.ailang/state/mission-v1*`.
   Do not add `t.Parallel()`.
9. **Evidence** goes to new files `$BANK/executor-<M>-<what>.txt`, committed in the evidence
   commit. Big raw logs stay in `$SCR` and are not committed; leave `$SCR` in place (the
   controller removes it).
10. **Bound every wait.** Wrap every `go test` in `perl -e 'alarm N; exec @ARGV' …` with the N
    given, and give the Bash tool a timeout of 590000 ms. `rc=142` means the alarm fired: re-run
    that one command once; a second 142 is a STOP (report it as UNMEASURED for the controller).
    Never end your turn announcing a wait.
11. **One package-scale `go test` at a time.** Never run two at once, and never run G4 while a
    mutant run or AC-STABLE is going.

---

## 2. Banked sources and base pins

Base sha256 (`$BANK/plan-base-SHA256SUMS.txt`, measured at `762d8cb4`):

| File | sha256 at base |
|---|---|
| `host/broker/exec_probe.go` | `390c74e841f26a1283743c4dc95efe6d31baf3ccc44b77e64d6af8283e8a23cd` |
| `host/broker/exec_probe_test.go` | `6c4e0e7a91a4fae16ae039e69fea1ad5b5975567c5290e701b14af879d13f110` |
| `.github/workflows/ci.yml` | `38485513f306728b614e0faf3b04085d20f8c85737d813da84875810453f6e4e` |

Banked files (`$BANK/…`). The `plan-M<k>-*` files are cumulative: M2's files are the final state.

| Banked file | Milestone | Destination / role | Lines | sha256 |
|---|---|---|---|---|
| `plan-M1-exec_probe.go.txt` | M1 | `host/broker/exec_probe.go` | 885 | `6e9c8315920985df3aa53ec01a4de29eaf1b3a731844f9854593ad77f6323efe` |
| `plan-M1-exec_probe_test.go.txt` | M1 | `host/broker/exec_probe_test.go` | 918 | `008c2a8de5d23d8ac8a886eb093cb3c1907290be5777dcf38402f18a4fabfee5` |
| `plan-M1-ci.yml.txt` | M1 | `.github/workflows/ci.yml` | 503 | `f511dd8e71f2e05e18d190ced31ff0341a1303a86cfea9a63c41d7fb95815753` |
| `plan-M2-exec_probe.go.txt` | M2 | `host/broker/exec_probe.go` | 913 | `96f33172667d27defa0c52735d5cff1ad74d4e907ab00a4c056c4e5a26f698a4` |
| `plan-M2-exec_probe_test.go.txt` | M2 | `host/broker/exec_probe_test.go` | 933 | `e35784db17623df899e573096edf41ea4b205dbb4975701380594082e8181c5e` |
| `plan-M2-ci.yml.txt` | M2 | `.github/workflows/ci.yml` | 503 | `daf57b6e168f616cd99e370df8ded5c853df72207c3cc5ec248110702c0ecead` |
| `plan-base-SHA256SUMS.txt`, `plan-M1-SHA256SUMS.txt`, `plan-M2-SHA256SUMS.txt` | — | `shasum -a 256 -c` inputs (repo-relative paths) | 3 each | — |
| `plan-run-mutants.sh.txt` | M1, M2 | the mutation runner (§5) | 63 | see §10 V-P9 |
| `plan-M1.diff.txt`, `plan-M2.diff.txt` | — | **review only**: the per-milestone diffs (Appendix A). Never apply them | 156 / 91 | — |
| `plan-planner-validation.txt`, `plan-planner-mutants-M1.txt`, `plan-planner-mutants-M2.txt` | — | the planner's own runs (§8, §9) | — | — |

---

## 3. Pristine baseline (record BEFORE the first copy)

```sh
shasum -a 256 -c "$BANK/plan-base-SHA256SUMS.txt" > "$BANK/executor-M0-base.txt" 2>&1; echo "rc=$?"
gofmt -l host/broker/ >> "$BANK/executor-M0-base.txt"; go vet ./host/broker/ >> "$BANK/executor-M0-base.txt" 2>&1; echo "vet rc=$?"
cat "$BANK/executor-M0-base.txt"
```

Expected: three `OK` lines, rc 0, gofmt prints nothing, vet rc 0 (all measured by the planner).

**Base RED / GREEN.** At base T10 and T11 do not exist (T10 cannot compile without
`execProbePollEvery`). The design's mutants K1 and K6 survive the whole base package (design V3).
Every existing test is green at base and must stay green.

---

## 4. Milestones

### Milestone M1 — pin the hello wait from both sides; the SILENT text names the wait

**Files:** `host/broker/exec_probe.go`, `host/broker/exec_probe_test.go`, `.github/workflows/ci.yml`.

#### M1.1 Copy and verify

```sh
cp "$BANK/plan-M1-exec_probe.go.txt"      host/broker/exec_probe.go
cp "$BANK/plan-M1-exec_probe_test.go.txt" host/broker/exec_probe_test.go
cp "$BANK/plan-M1-ci.yml.txt"             .github/workflows/ci.yml
shasum -a 256 -c "$BANK/plan-M1-SHA256SUMS.txt"; echo "rc=$?"      # three OK, rc 0
git diff --name-only                                                 # exactly the three files
```

What M1 contains (read Appendix A for the diff; do not apply it):
- P1: `const execProbePollEvery = 50 * time.Millisecond` after `execProbeHelloWait`; `pollLoopback`'s
  `time.After(50 * time.Millisecond)` becomes `time.After(execProbePollEvery)`; its doc comment says
  "every execProbePollEvery".
- P2: the SILENT `Sprintf` prints `within %d ms;` with `execProbeHelloWait.Milliseconds()`. At 500
  the bytes are identical to today's.
- T7 `TestExecProbeReverseLegSilentHolderRefusesSticky` gains one assertion per row:
  `Why` contains `answered nothing to World's hello within <execProbeHelloWait.Milliseconds()> ms;`.
- `helloConn`, `firstDialReturns`, `execProbeHelloFloor = 250 * time.Millisecond`, and T10 with
  rows `slow-child` and `slow-mux`.
- `ci.yml:247`'s `tests=` line ends `… TestExecProbeReversePollRecordsConnectReset TestExecProbeReverseHelloWaitBounds"`.

#### M1.2 Acceptance (M1)

```sh
# G1 gofmt (want: no output)
gofmt -l host/broker/
# G2 vet (want rc 0)
go vet ./host/broker/; echo "rc=$?"
# G3 compile fence (want rc 0 and "ok ... [no tests to run]")
perl -e 'alarm 300; exec @ARGV' go test -run '^$' ./host/broker/ > "$SCR/m1-g3.txt" 2>&1; echo "rc=$?"; cat "$SCR/m1-g3.txt"
# G8 CI file (want rc 0, "yaml ok", count 1)
actionlint .github/workflows/ci.yml; echo "rc=$?"
ruby -ryaml -e 'YAML.load_file(".github/workflows/ci.yml"); puts "yaml ok"'
grep -c 'TestExecProbeReversePollRecordsConnectReset TestExecProbeReverseHelloWaitBounds"' .github/workflows/ci.yml
# focused run
perl -e 'alarm 300; exec @ARGV' go test ./host/broker/ -run '^TestExecProbeReverse|^TestExecProbeChildListener|^TestExecProbeForeignListener' -count=1 -v > "$SCR/m1-focus.txt" 2>&1; echo "rc=$?"
grep -c '^--- PASS' "$SCR/m1-focus.txt"        # want 10
grep -c '^    --- PASS' "$SCR/m1-focus.txt"    # want 34
grep -c -- '--- FAIL\|--- SKIP' "$SCR/m1-focus.txt"   # want 0
grep -- '--- PASS: TestExecProbeReverseHelloWaitBounds' "$SCR/m1-focus.txt"   # want 3 lines: parent, /slow-child, /slow-mux
```

Planner (§10 V-P4): gofmt empty, vet rc 0, fence rc 0, actionlint rc 0, count 1, focused rc 0,
**10** top-level PASS (the 9 base names `…LegJudge`, `…ChattyForeignHolderIsRetried`,
`…ExhaustionRefusesDistinctly`, `…NoTokenIsNotRetried`, `…SilentHolderRefusesSticky`,
`…RealReachSurvivesAForeignAttempt`, `TestExecProbeForeignListenerOnTheReversePortDoesNotRefuse`,
`TestExecProbeReverseDialErrIsSilent`, `TestExecProbeReversePollRecordsConnectReset`, plus T10),
**34** subtest PASS, 0 FAIL, 0 SKIP.

**G4 at M1 (full package, about 3.5 min).**

```sh
perl -e 'alarm 570; exec @ARGV' sh -c 'go build ./... && go test -count=1 -v ./host/broker/' > "$SCR/m1-g4.txt" 2>&1; echo "rc=$?"
tail -1 "$SCR/m1-g4.txt"
printf 'PASS=%s FAIL=%s SKIP=%s\n' "$(grep -c '^--- PASS' "$SCR/m1-g4.txt")" "$(grep -c '^--- FAIL' "$SCR/m1-g4.txt")" "$(grep -c '^--- SKIP' "$SCR/m1-g4.txt")"
```

Want rc 0, `ok  	github.com/sunholo-data/ailang-world/host/broker`, **PASS=334 FAIL=0 SKIP=30** (the
design's V3 baseline 333/0/30 plus T10). Planner: rc 0, 217 s, 334/0/30, and the top-level verdict
set differs from the design's V3 control summary only by `+ --- PASS: TestExecProbeReverseHelloWaitBounds`.
If the only failure is `TestExecProtectedCheckTimeout/reserve-after-before` (a known load-sensitive
3 s subprocess bound, unrelated to this diff), re-run it alone:
`perl -e 'alarm 300; exec @ARGV' go test ./host/broker/ -run '^TestExecProtectedCheckTimeout$/^reserve-after-before$' -count=3 -v`,
record both results, and continue only if the re-run is 3/3 PASS. Any other failure is a STOP.

#### M1.3 Mutation drill (M1)

```sh
perl -e 'alarm 570; exec @ARGV' bash "$BANK/plan-run-mutants.sh.txt" M1 > "$BANK/executor-M1-mutants.txt" 2>&1; echo "rc=$?"
grep -E '^== |restored|start' "$BANK/executor-M1-mutants.txt"
shasum -a 256 -c "$BANK/plan-M1-SHA256SUMS.txt"; echo "rc=$?"     # three OK: the tree is restored
```

Want rc 0 and exactly these `==` lines (planner, §8):

```
== control-T10 rc=0 fail=0 pass=3
== K1-wait-5ms rc=1 fail=3 pass=0
== K2-wait-249ms rc=1 fail=3 pass=0
== K2c-wait-250ms rc=0 fail=0 pass=3
== K2c-wait-950ms rc=0 fail=0 pass=3
== K3-wait-951ms rc=1 fail=3 pass=0
== K4-wait-60s rc=1 fail=3 pass=0
== K5-no-deadline rc=1 fail=3 pass=0
== K9-hold-500 rc=1 fail=3 pass=0
== K8-poll-literal rc=0 fail=0 pass=3
== control-T7 rc=0 fail=0 pass=4
== K7-literal-text-wait-400ms rc=1 fail=4 pass=0
== R-M1-P2-reverted-alone rc=0 fail=0 pass=4
```

and a final `restored sha256 6e9c8315…efe (start 6e9c8315…efe)`. `NOT-APPLIED`, `BUILD-FAILED` or
`RESTORE MISMATCH` anywhere is a STOP. K2c, K8 and R-M1 are **in-table controls**: they must PASS.

#### M1.4 Commit (M1)

```sh
git add host/broker/exec_probe.go host/broker/exec_probe_test.go .github/workflows/ci.yml
git commit -F - <<'EOF'
fix(row 230 M1): pin execProbeHelloWait from both sides (T10); SILENT text names the wait

P1 names the reverse leg's poll period (execProbePollEvery). P2 derives the
SILENT refusal's "within N ms" from execProbeHelloWait (byte-identical at
500 ms). T10 TestExecProbeReverseHelloWaitBounds pins the wait between a
250 ms floor (a simulated slow answer must still read child/foreign) and a
ceiling of execProbeListenMS less execProbePollEvery; T7 checks the text
names the wait. CI step 8 PASS-guards T10. Kills K1-K5, K7, K9 (judge O8).
EOF
git log --oneline -1
```

Add your session's `Co-Authored-By:` trailer as the message's last line (put it inside the heredoc
after a blank line).

---

### Milestone M2 — the structural `c.resume()` pin and the comment-only P3/P4

**Files:** the same three.

#### M2.1 Copy and verify

```sh
cp "$BANK/plan-M2-exec_probe.go.txt"      host/broker/exec_probe.go
cp "$BANK/plan-M2-exec_probe_test.go.txt" host/broker/exec_probe_test.go
cp "$BANK/plan-M2-ci.yml.txt"             .github/workflows/ci.yml
shasum -a 256 -c "$BANK/plan-M2-SHA256SUMS.txt"; echo "rc=$?"      # three OK, rc 0
git diff --stat                                                      # exec_probe.go +29/-1 (comments), exec_probe_test.go +15, ci.yml 1 +/1 -
```

What M2 contains (Appendix A): P3, the `execProbeHelloWait` doc comment stating its two bounds,
and a paragraph above `execProbeScript` on what `c.resume()` does and does not do; P4, the
ordering-invariant comments at the fan-out `go func()` and at `net6 := newProbeListener()`; T11;
and `ci.yml:247` ending `… TestExecProbeReverseHelloWaitBounds TestExecProbeChildListenerDrainsTheHello"`.
Every changed line of `exec_probe.go` in M2 is a comment (planner-checked, V-P6).

#### M2.2 Acceptance (M2)

```sh
gofmt -l host/broker/                                                       # G1: no output
go vet ./host/broker/; echo "rc=$?"                                         # G2: rc 0
perl -e 'alarm 300; exec @ARGV' go test -run '^$' ./host/broker/ > "$SCR/m2-g3.txt" 2>&1; echo "rc=$?"   # G3: rc 0
actionlint .github/workflows/ci.yml; echo "rc=$?"                           # G8: rc 0
ruby -ryaml -e 'YAML.load_file(".github/workflows/ci.yml"); puts "yaml ok"'
grep -c 'TestExecProbeReverseHelloWaitBounds TestExecProbeChildListenerDrainsTheHello"' .github/workflows/ci.yml   # want 1
perl -e 'alarm 300; exec @ARGV' go test ./host/broker/ -run '^TestExecProbeReverse|^TestExecProbeChildListener|^TestExecProbeForeignListener' -count=1 -v > "$SCR/m2-focus.txt" 2>&1; echo "rc=$?"
grep -c '^--- PASS' "$SCR/m2-focus.txt"        # want 11
grep -c '^    --- PASS' "$SCR/m2-focus.txt"    # want 34
grep -c -- '--- FAIL\|--- SKIP' "$SCR/m2-focus.txt"   # want 0
```

**G5 race (the reverse-leg family):**

```sh
perl -e 'alarm 570; exec @ARGV' go test -race ./host/broker/ -run '^TestExecProbeReverse|^TestExecProbeChildListener|^TestExecProbeForeignListener' -count=1 -v > "$SCR/m2-race.txt" 2>&1; echo "rc=$?"
grep -c 'DATA RACE' "$SCR/m2-race.txt"   # want 0
grep -c '^--- PASS' "$SCR/m2-race.txt"   # want 11
```

**AC-STABLE (G6, about 95 s):**

```sh
perl -e 'alarm 570; exec @ARGV' go test ./host/broker/ -run '^TestExecProbeReverse|^TestExecProbeChildListener' -count=20 -v > "$SCR/m2-c20.txt" 2>&1; echo "rc=$?"
bash -c 'for t in TestExecProbeReverseLegJudge TestExecProbeReverseLegChattyForeignHolderIsRetried TestExecProbeReverseLegExhaustionRefusesDistinctly TestExecProbeReverseLegNoTokenIsNotRetried TestExecProbeReverseLegSilentHolderRefusesSticky TestExecProbeReverseLegRealReachSurvivesAForeignAttempt TestExecProbeReverseDialErrIsSilent TestExecProbeReversePollRecordsConnectReset TestExecProbeReverseHelloWaitBounds TestExecProbeChildListenerDrainsTheHello; do printf "%s %s\n" "$(grep -c -- "^--- PASS: $t " "$0")" "$t"; done' "$SCR/m2-c20.txt" > "$BANK/executor-M2-stable.txt"
grep -c -- '--- FAIL' "$SCR/m2-c20.txt" >> "$BANK/executor-M2-stable.txt"
cat "$BANK/executor-M2-stable.txt"     # want ten lines starting "20 ", then a last line "0"
```

**T10 under load (AC-STABLE's second half, about 25 s).** `run-load.sh.txt` (the design's runner)
starts 32 `yes`, kills exactly those PIDs on exit, and prints `yes-alive-after=0`.

```sh
perl -e 'alarm 300; exec @ARGV' go test -c -o "$SCR/broker.test" ./host/broker/; echo "rc=$?"
( cd host/broker && perl -e 'alarm 300; exec @ARGV' bash "../../$BANK/run-load.sh.txt" "$SCR/broker.test" -test.run '^TestExecProbeReverseHelloWaitBounds$' -test.count=3000 -test.v ) > "$SCR/m2-t10load.txt" 2>&1; echo "rc=$?"
{ grep 'load averages\|load runner' "$SCR/m2-t10load.txt"; grep -c '^--- PASS: TestExecProbeReverseHelloWaitBounds ' "$SCR/m2-t10load.txt"; grep -c -- '--- FAIL' "$SCR/m2-t10load.txt"; } >> "$BANK/executor-M2-stable.txt"
tail -4 "$BANK/executor-M2-stable.txt"   # want: an uptime line, "load runner rc=0 yes-alive-after=0", 3000, 0
pgrep -x yes | wc -l                       # want 0
```

**G4 at M2 (full package, about 3.5 min), then G7:**

```sh
perl -e 'alarm 570; exec @ARGV' sh -c 'go build ./... && go test -count=1 -v ./host/broker/' > "$SCR/m2-g4.txt" 2>&1; echo "rc=$?"
tail -1 "$SCR/m2-g4.txt"
printf 'PASS=%s FAIL=%s SKIP=%s\n' "$(grep -c '^--- PASS' "$SCR/m2-g4.txt")" "$(grep -c '^--- FAIL' "$SCR/m2-g4.txt")" "$(grep -c '^--- SKIP' "$SCR/m2-g4.txt")"
perl -e 'alarm 570; exec @ARGV' ./scripts/verify_ail.sh > "$SCR/m2-g7.txt" 2>&1; echo "rc=$?"; tail -1 "$SCR/m2-g7.txt"
```

Want G4 rc 0 and **PASS=335 FAIL=0 SKIP=30** (V3's 333/0/30 plus T10 and T11); the
`reserve-after-before` rule from M1.2 applies. Want G7 rc 0 and
`✓ verify gate PASSED: 16 required identities verified, 40 named tests pass, 503 se-tools named tests pass`.

#### M2.3 Mutation drill (M2 = the final state: every K)

```sh
perl -e 'alarm 570; exec @ARGV' bash "$BANK/plan-run-mutants.sh.txt" M2 > "$BANK/executor-M2-mutants.txt" 2>&1; echo "rc=$?"
grep -E '^== |restored|start' "$BANK/executor-M2-mutants.txt"
shasum -a 256 -c "$BANK/plan-M2-SHA256SUMS.txt"; echo "rc=$?"     # three OK
```

Want rc 0, the thirteen `==` lines of M1.3 again, then:

```
== control-T11 rc=0 fail=0 pass=1
== K6-no-resume rc=1 fail=1 pass=0
```

and `restored sha256 96f33172…8a4 (start 96f33172…8a4)`.

#### M2.4 Commit (M2)

```sh
git diff --name-only    # exactly the three files
git add host/broker/exec_probe.go host/broker/exec_probe_test.go .github/workflows/ci.yml
git commit -F - <<'EOF'
test(row 230 M2): pin the child listener's c.resume() (T11); document the wait bounds and net6 ordering

T11 TestExecProbeChildListenerDrainsTheHello is a structural (S6) pin of
the listen op's handler (judge O10, K6). P3 documents execProbeHelloWait's
floor and ceiling and what c.resume() does and does not do; P4 records the
net6-before-Run-A ordering invariant at the fan-out and at the bind
(comments only). CI step 8 PASS-guards T11.
EOF
git log --oneline -2
```

(Again add your `Co-Authored-By:` trailer as the last line inside the heredoc.)

---

## 5. Mutation table (design §5; every row run by the planner at the final state)

The runner (`plan-run-mutants.sh.txt`) applies each edit with `sed` to `exec_probe.go` only, runs
only the killer, restores from the banked copy and re-checks the sha256. Killer regexes:
T10 `^TestExecProbeReverseHelloWaitBounds$`, T7 `^TestExecProbeReverseLegSilentHolderRefusesSticky$`,
T11 `^TestExecProbeChildListenerDrainsTheHello$`. Line numbers are the final-state test file's.

| ID | Mutant (exact edit) | Milestone | Killer | Planner result and the assertion that fired |
|---|---|---|---|---|
| K1 (= O8) | `const execProbeHelloWait = 5 * time.Millisecond` | M1 | T10 | **RED** 3 FAIL (both rows + parent): `exec_probe_test.go:911: an answer 250ms after the hello: seen {child:false foreign:false silent:true} (wrote "WORLD-PROBE-HELLO\n"); want {child:true foreign:false silent:false}, not silent` (and `want {child:false foreign:true …}` for `slow-mux`) |
| K2 | `= 249 * time.Millisecond` | M1 | T10 | **RED** 3 FAIL, the same `:911` assertion |
| K2c | `= 250 * time.Millisecond` / `= 950 * time.Millisecond` | M1 | T10 | **GREEN** / **GREEN** (control: the bounds are inclusive) |
| K3 | `= 951 * time.Millisecond` | M1 | T10 | **RED** 3 FAIL: `exec_probe_test.go:914: hello deadline at stop+950.905583ms (sets 1); want a deadline no later than stop+950ms (listen window less one poll period)` |
| K4 | `= 60 * time.Second` | M1 | T10 | **RED** 3 FAIL: `:914: hello deadline at stop+59.999897916s (sets 1); want … no later than stop+950ms …`. (T6/`silent`'s wall-clock red is design V4; not re-run, D-9) |
| K5 | `_ = c.SetDeadline(time.Now().Add(execProbeHelloWait))` → `// no deadline (mutant K5)` | M1 | T10 | **RED** 3 FAIL: `:914: hello deadline at stop+-2562047h47m16.854775808s (sets 0); want …` |
| K7 | P2 reverted (`within 500 ms; …", port)`) **and** `execProbeHelloWait = 400 * time.Millisecond` | M1 | T7 | **RED** 4 FAIL (all 3 rows + parent): `exec_probe_test.go:669: Why = "reverse leg: a listener on 127.0.0.1:61901 … within 500 ms; …"; want it to contain "answered nothing to World's hello within 400 ms;"` |
| K9 | `const execProbeListenMS = 500` | M1 | T10 | **RED** 3 FAIL: `:914: hello deadline at stop+499.815416ms (sets 1); want a deadline no later than stop+450ms …` (the intended coupling, design §5) |
| K8 | `time.After(execProbePollEvery)` → `time.After(50 * time.Millisecond)` | M1 | T10 | **GREEN** (control: behaviour-identical, so not a mutant; design §5) |
| R-M1 | P2 reverted alone (wait stays 500) | M1 | T7 | **GREEN** (control: P2 is byte-identical at 500 ms; D-1) |
| K6 (= O10) | `net.createServer(c=>{c.resume();c.end(` → `net.createServer(c=>{c.end(` | M2 | T11 | **RED** 1 FAIL: `exec_probe_test.go:931: the listen op's connection handler is not net.createServer(c=>{c.resume();c.end('WORLD-PROBE-NONCE '+nonce+'\n')})` |

Every RED was rc 1 with the listed `--- FAIL:` lines; every control rc 0. The file restored to its
start sha256 after every mutant (`plan-planner-mutants-M1.txt` at M1, `plan-planner-mutants-M2.txt`
at the final state, both byte-for-byte the runner the executor uses).

**Non-vacuity against each milestone's own production diff (measured; D-1).**
- **M1.** Its production diff is behaviour-identical at the shipped 500 ms by design (P1 renames a
  literal; P2 prints the same bytes; design §9). So no test can red on that diff's revert alone,
  and none does: R-M1 (P2 reverted) is GREEN and K8 (P1's use reverted) is GREEN. Reverting the
  whole M1 production hunk is a compile error, not a red: `exec_probe_test.go:896:28: undefined:
  execProbePollEvery` (`[build failed]`). The headline red is therefore measured on the lines the
  tests pin: **T10 reds on K1** (the O8 line; it is green at base's own value), and **the T7
  extension reds on K7**, the composite that makes P2's absence observable.
- **M2.** Its production diff is comments only (29 changed lines of `exec_probe.go`, 0 non-comment;
  V-P6). With M1's `exec_probe.go` and M2's test file, T11 is `--- PASS`. The headline red is
  **T11 on K6** (the O10 line).

---

## 6. Evidence commit

```sh
bash -c 'for f in m1-focus m1-g4 m2-focus m2-race m2-g4 m2-g7; do printf "== %s\n" "$f"; grep -E "^(ok|FAIL|PASS)|DATA RACE|verify gate" "'"$SCR"'/$f.txt" | tail -3; done' > "$BANK/executor-summary.txt"
printf 'M1 G4: PASS=%s FAIL=%s SKIP=%s\n' "$(grep -c '^--- PASS' "$SCR/m1-g4.txt")" "$(grep -c '^--- FAIL' "$SCR/m1-g4.txt")" "$(grep -c '^--- SKIP' "$SCR/m1-g4.txt")" >> "$BANK/executor-summary.txt"
printf 'M2 G4: PASS=%s FAIL=%s SKIP=%s\n' "$(grep -c '^--- PASS' "$SCR/m2-g4.txt")" "$(grep -c '^--- FAIL' "$SCR/m2-g4.txt")" "$(grep -c '^--- SKIP' "$SCR/m2-g4.txt")" >> "$BANK/executor-summary.txt"
cat "$BANK/executor-summary.txt"
git add "$BANK/executor-M0-base.txt" "$BANK/executor-M1-mutants.txt" "$BANK/executor-M2-mutants.txt" "$BANK/executor-M2-stable.txt" "$BANK/executor-summary.txt"
git commit -F - <<'EOF'
evidence(row 230): executor gates, mutation drill K1-K9, AC-STABLE, full package
EOF
git status --porcelain     # want no line for host/, .github/ or $BANK (other untracked trees are not yours)
```

(Trailer as before.) Then report: the three commit SHAs, every rc and count above, and anything
that differed from the planner's numbers.

---

## 7. Controller-only (not the executor)

- **Full `./scripts/verify_go.sh`** (plain and race legs; the repo-wide `go test ./...`).
- **AC-CI on the PR:** step 8 ("Row 140 Workspace.Exec real-srt tests, verbose") must print
  `--- PASS: TestExecProbeReverseHelloWaitBounds ` and `--- PASS: TestExecProbeChildListenerDrainsTheHello `
  on linux (the step's own guard at `ci.yml:254-257` fails otherwise).
- **AC-LINUX-F3** (design §3.3, non-blocking): run on `ubuntu-24.04` from a throwaway branch, with
  T11 excluded from the mutant-applied step-8 run (or the branch cut from the M1 commit, where T11
  does not exist). Record under `design_docs/verification/world-iter329/`.
- **Records:** file the successor row for the `net6` hoist (design §3.4), carry §3.4/§3.5's
  verdicts and triggers into the charter row's closing note, and move the design to
  `implemented/`. None of these is the sprint's edit.

---

## 8. Stability and full package (planner, this session)

Load is the 1-minute average from `uptime`. The host already carried other work (29.27 before any
`yes` started).

| Run | State | Load seen | Result |
|---|---|---|---|
| T10 `-count=3000 -v` | final (`b93120db`, D-11) | 29.27 (no `yes`; other host load) | rc 0, 3000 PASS, 0 FAIL |
| T10 `-count=3000` under `run-load.sh.txt` (32 `yes`), 2 s after start | final (`b93120db`) | 29.56 | 3000 PASS, 0 FAIL, `yes-alive-after=0` |
| T10 `-count=3000` under 32 `yes` after a 20 s warm-up | final (`b93120db`) | 36.76 at start, 37.50 at end | 3000 PASS (6000 subtest PASS), 0 FAIL, `yes-alive-after=0` |
| AC-STABLE `-run '^TestExecProbeReverse\|^TestExecProbeChildListener' -count=20 -v` | final (`b93120db`) | 36.41 → 13.21 | rc 0, 95 s; each of the 10 names `--- PASS` ×20; 0 FAIL, 0 SKIP |
| `-race` reverse-leg family + T6 + T11, `-count=1` | final (`b93120db`) | ~6 | rc 0, 0 `DATA RACE`, 11 PASS |
| full `go test -count=1 -v ./host/broker/` | M1 | 5.85 → 6.76 | rc 0, 217 s, **334 / 0 / 30** |
| full `go test -count=1 -v ./host/broker/` | final (`b93120db`) | 12.47 → 5.83 | rc 0, 202 s, **335 / 0 / 30** |
| full `go test -count=1 -v ./host/broker/` | **final banked (`96f33172`)** | 9.91 → 9.23 | rc 0, 216 s, **335 / 0 / 30**; `reserve-after-before` PASS 8.48 s |
| AC-STABLE `-count=20` | **final banked (`96f33172`)** | 9.23 → 11.18 | rc 0, 95 s; each of the 10 names ×20; 0 FAIL, 0 SKIP |
| T10 `-count=3000` under 32 `yes` after a 20 s warm-up | **final banked (`96f33172`)** | 19.37 | 3000 PASS, 0 FAIL, `yes-alive-after=0` |
| the executor's exact M2.2 load snippet (no warm-up), rehearsed | **final banked (`96f33172`)** | 16.15 | 3000 PASS, 0 FAIL, `yes-alive-after=0`; `pgrep -x yes` 0 afterwards |
| `-race` reverse-leg family + T6 + T11 | **final banked (`96f33172`)** | ~10 | rc 0, 0 `DATA RACE`, 11 PASS, 0 FAIL |

Versus the design's V3 baseline (333 PASS / 0 FAIL / 30 SKIP): the top-level verdict set differs
only by the added `--- PASS: TestExecProbeChildListenerDrainsTheHello` and
`--- PASS: TestExecProbeReverseHelloWaitBounds`; the 30 SKIPs are the same names.
`TestExecProtectedCheckTimeout/reserve-after-before` PASSed in all three full runs (8.54 s, 8.42 s, 8.48 s), so its known load-sensitive local red did not occur.

---

## 9. Deviations from the design

- **D-1 (non-vacuity rule vs behaviour-identical production diffs).** The mission rule asks that
  each milestone's headline test red when that milestone's production hunk alone is reverted. Both
  milestones' production diffs are behaviour-identical by the design's own terms (§3.1 P1/P2, §9
  "P2 produces the same bytes"; P3/P4 are comments), so that red cannot exist. Measured and
  recorded instead (§5): R-M1 and K8 GREEN, the whole-M1 revert a compile error, the M2 revert
  GREEN; the headline reds are T10 on K1, the T7 extension on K7 and T11 on K6, each against the
  production line the test pins. R-M1 is a planner-added control row.
- **D-2 (T10 concretised).** The design's sketch names `firstDialReturns(conn, refused)` without a
  body; it is a plain closure that returns `conn` on call 1 and the error after, used only on the
  poll's goroutine (race-clean, §8). The ceiling is a local `ceiling := listenWindow -
  execProbePollEvery`. Both checks use `t.Errorf` so a mutant that breaks both shows both (the
  prototype's form, design V8). Otherwise T10 is the design's banked prototype
  (`zz_iter329_f2_test.go.txt`) with the literal 50 ms replaced by `execProbePollEvery`.
- **D-3 (T7 extension's form).** A separate `t.Fatalf` after T7's existing check, wanting the
  longer substring `answered nothing to World's hello within %d ms;` (the design's V13 prototype
  form), so K7's failing line is the extension's (`:669`), not the existing assertion (`:665`).
- **D-4 (one extra comment line in M1).** `pollLoopback`'s doc comment said "every 50 ms"; with P1
  it says "every execProbePollEvery". Not in design §7's site list.
- **D-5 (P3 wording).** The `c.resume()` comment says it "is not what avoids a reset at close (the
  earlier reading)" and scopes the measurement to darwin, with Linux unmeasured (AC-LINUX-F3), per
  design §1 row 2 and §2. P3's `execProbeHelloWait` comment documents M1's constants but lands in
  M2, as design §6 orders it.
- **D-6 (AC-GATE scope).** The design's AC-GATE says `go build ./... && go test ./...`. The
  executor runs `go build ./...`, the full `host/broker` package, the compile fence, vet and
  `verify_ail.sh`; the repo-wide `go test ./...` is the controller's `verify_go.sh` (§7). This is
  row 232's D7.
- **D-7 (design LOC estimate).** Design §7 estimates `exec_probe.go` +12/−3 and tests +90. Banked:
  `exec_probe.go` +35/−4 (P3/P4's comments are longer than estimated; code +6/−3), tests +102/−0.
- **D-8 (executor commits).** Unlike row 232's plan (controller commits), this plan has the
  executor commit once per milestone and once for evidence, with named-file staging, per the
  iteration-329 brief.
- **D-9 (K4's T6 side effect not re-run).** Design §5 lists T6/`silent` as K4's wall-clock killer
  (V4, 132 s). It is not re-run: T10 is the deterministic pin, and V4 is the record.
- **D-10 (focused regex).** The focused acceptance runs add `^TestExecProbeForeignListener` (T6)
  to AC-STABLE's regex so the whole reverse-leg family runs once per milestone. AC-STABLE itself
  keeps the design's regex.
- **D-11 (final comment wording changed after the stability runs).** The planner re-worded P3's
  `c.resume()` paragraph after the first final-state runs (sha `b93120db…` → `96f33172…`; comment
  lines only, 0 non-comment lines differ). gofmt, vet, the fence, the focused run and the full M2
  mutation drill, `-race`, the full package, AC-STABLE, the loaded T10 and `verify_ail.sh` were
  all re-run on `96f33172` (V-P12), with the same results.

Nothing in the design was found false. These held at `762d8cb4`: `:99` `execProbeListenMS = 1000`,
`:110` the wait, `:226` the fan-out `go func() {`, `:455` `net6`, `:582` the literal SILENT text,
`:662` `SetDeadline`, `:732` the 50 ms poll; `ci.yml:247` the reverse-leg `tests=` line, which now
also carries row 232's T8/T9; T7 at `exec_probe_test.go:646`; `revDialErr`, `newRevPicker`,
`runReverseLeg` and `revFixture` exist in the test file.

## 10. Verification notes (planner, this session)

| # | Claim | Command | Observed |
|---|---|---|---|
| V-P1 | Base | `git rev-parse HEAD`; `git diff --stat e5c1867f HEAD \| tail -1`; `shasum -a 256 -c plan-base-SHA256SUMS.txt` | `762d8cb4`; 35 files, all under `design_docs/`; 3 OK |
| V-P2 | Toolchain | `go version`; `$AILANG_BIN --version`; `which actionlint ruby` | `go1.26.6 darwin/arm64`; `AILANG v0.41.0`; both present |
| V-P3 | Scratch | `git worktree add --detach /Users/voightkampff/dev/sunholo-data/.planner-scratch-iter329 762d8cb4` | created; removed with `git worktree remove --force` at the end |
| V-P4 | M1 state green | §4 M1.2 block in the scratch tree, files copied from the bank | gofmt empty; vet 0; fence 0; actionlint 0; YAML ok; anchor 1; focused 10 / 34 / 0 / 0; full package 334/0/30 in 217 s |
| V-P5 | M1 drill from the bank | `bash $BANK/plan-run-mutants.sh.txt M1` after a fresh `git checkout` of the base and `cp` of the M1 bank | the thirteen `==` lines of M1.3; restored `6e9c8315…` (`plan-planner-mutants-M1.txt`) |
| V-P6 | M2 production diff is comment-only | `diff plan-M1-exec_probe.go.txt plan-M2-exec_probe.go.txt \| grep '^[<>]' \| grep -v '^[<>][[:space:]]*//'` | no lines (29 changed lines in all, every one a comment) |
| V-P7 | M2 state green | §4 M2.2 block | gofmt empty; vet 0; fence 0; actionlint 0; YAML ok; anchor 1; focused 11 / 34 / 0 / 0; `go build ./...` 0; `-race` 0 races; `verify_ail.sh` rc 0 in 47 s (48 s on `96f33172`) |
| V-P8 | M2 drill from the bank | `bash $BANK/plan-run-mutants.sh.txt M2` | the fifteen `==` lines of M2.3; restored `96f33172…` (`plan-planner-mutants-M2.txt`) |
| V-P9 | Runner identity | `shasum -a 256 $BANK/plan-run-mutants.sh.txt` | `6cf4ff4fe720fddf174a40b2e7f39a43be034eb8aa42e24e60c70ecec5f52610` (the copy both drills ran) |
| V-P10 | Whole-M1 revert | base `exec_probe.go` + M1 test file; `go test -run '^(TestExecProbeReverseHelloWaitBounds\|TestExecProbeReverseLegSilentHolderRefusesSticky)$'` | rc 1, `exec_probe_test.go:896:28: undefined: execProbePollEvery`, `[build failed]` |
| V-P11 | M2 revert | M1 `exec_probe.go` + M2 test file; T11 | `--- PASS: TestExecProbeChildListenerDrainsTheHello` |
| V-P12 | Final sha re-runs | on `96f33172`: focused run, M2 drill, `-race`, full package, AC-STABLE `-count=20`, loaded T10 `-count=3000` (warm-up and the executor's exact snippet) | focused 11/34/0; drill as M2.3; race 0 races, 11 PASS; full 335/0/30 in 216 s; AC-STABLE 10×20, 0 FAIL; T10 3000/0 at load 19.37 and 3000/0 at 16.15; `plan-planner-validation.txt` §"Final sha re-runs" |
| V-P13 | Executor snippets rehearsed | the M2.2 AC-STABLE count loop, the load snippet and the §6 summary loop, run in zsh against the planner's logs | the AC-STABLE loop printed ten `20 …` lines and `0`; the summary loop printed `ok` lines for all six logs and the verify-gate line; `M1 G4: PASS=334 FAIL=0 SKIP=30` |

---

## Appendix A — the diffs, for review only (the executor copies files; never apply these)

`plan-M1.diff.txt` (base → M1) and `plan-M2.diff.txt` (M1 → M2) hold the full unified diffs. The
production hunks:

M1, `host/broker/exec_probe.go`:

```diff
@@ const execProbeHelloWait = 500 * time.Millisecond
+// execProbePollEvery is the reverse leg's dial cadence (pollLoopback).
+const execProbePollEvery = 50 * time.Millisecond
@@ reverseLeg, case revSilent:
-			return first, fmt.Sprintf("reverse leg: a listener on 127.0.0.1:%d accepted the host's dial and answered nothing to World's hello within 500 ms; World cannot tell it from the child's listener, so it refuses", port)
+			return first, fmt.Sprintf("reverse leg: a listener on 127.0.0.1:%d accepted the host's dial and answered nothing to World's hello within %d ms; World cannot tell it from the child's listener, so it refuses", port, execProbeHelloWait.Milliseconds())
@@ pollLoopback doc comment
-// pollLoopback dials 127.0.0.1:port every 50 ms until stopped and
+// pollLoopback dials 127.0.0.1:port every execProbePollEvery until stopped and
@@ pollLoopback select
-			case <-time.After(50 * time.Millisecond):
+			case <-time.After(execProbePollEvery):
```

M1, `host/broker/exec_probe_test.go` (T10's core; the full text is in the bank):

```go
func TestExecProbeReverseHelloWaitBounds(t *testing.T) {
	want := "WORLD-PROBE-NONCE " + strings.Repeat("a", 32)
	refused := revDialErr(syscall.ECONNREFUSED)
	listenWindow := time.Duration(execProbeListenMS) * time.Millisecond
	ceiling := listenWindow - execProbePollEvery
	// rows: slow-child (want+"\n" -> {child}), slow-mux ("HTTP/1.1 400 Bad Request\r\n" -> {foreign})
	…
			before := time.Now()
			conn := &helloConn{arrival: before.Add(execProbeHelloFloor), answer: []byte(row.answer)}
			seen := pollLoopback(40000, want, firstDialReturns(conn, refused))()
			after := time.Now()
			if seen != row.want || conn.wrote != execProbeHello { t.Errorf(…) }                                 // floor
			if conn.sets == 0 || conn.deadline.IsZero() || conn.deadline.After(after.Add(ceiling)) { t.Errorf(…) } // ceiling
```

M2 is comments in `exec_probe.go` (above `execProbeScript`, on `execProbeHelloWait`, at the
fan-out `go func()`, at `net6 := newProbeListener()`), T11 appended to the test file, and the CI
line. T11:

```go
func TestExecProbeChildListenerDrainsTheHello(t *testing.T) {
	const handler = "net.createServer(c=>{c.resume();c.end('WORLD-PROBE-NONCE '+nonce+'\\n')})"
	if n := strings.Count(execProbeScript, "net.createServer("); n != 1 {
		t.Fatalf("execProbeScript has %d net.createServer( calls, want exactly 1 (the listen op)", n)
	}
	if !strings.Contains(execProbeScript, handler) {
		t.Fatalf("the listen op's connection handler is not %s", handler)
	}
}
```
