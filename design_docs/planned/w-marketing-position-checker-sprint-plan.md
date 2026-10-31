# Sprint plan — `w-marketing-position-checker` (iteration 328, queue row 231)

**Design doc:** [`design_docs/planned/w-marketing-position-checker.md`](w-marketing-position-checker.md)
r2 (`08dc31be`; the controller's narrow-refinement carve-out after quorum r2, design §10).
**Sprint base:** worktree `/Users/voightkampff/dev/sunholo-data/.wt-world-iter328-row231`, branch
`sprint/world-iter328-marketing-position-checker`, HEAD `08dc31be`. Its Go tree equals the design's
measurement base `cec48f69` (only `design_docs/` changed since).
**Planner:** mission-control iteration 328 (planner role). **Executor:** Claude Haiku 5.5.
**JSON:** [`sprint_w-marketing-position-checker.json`](sprint_w-marketing-position-checker.json).

**How this sprint works.** The planner wrote every source file of the row, assembled it in throwaway
copies of the base tree (`git archive HEAD`, siblings of the repo, never `/tmp`), ran every
acceptance command and every §4 mutant, and banked the validated sources as text under
`design_docs/verification/world-iter328/plan-*.txt`. **The executor writes no Go code.** For each
milestone it copies the banked files into place with `cp`, checks their sha256, runs the
generators and the gate, banks the output, and commits: seven milestone commits, each independently
green, plus one evidence commit.
**Copy commands from fenced blocks, never from table cells.**

**Estimate:** about 0.5 executor day; wall time is dominated by the package's `-race` runs
(about 2 min each, 7 milestones) and the repo gate.

---

## 0. The one-paragraph version

`p0check` and `phonecheck` become `marketingcheck position` and `marketingcheck phone` (exit codes
0/1/2), with the iteration-315 corpora as CI fixtures. The banked `.txt` programs stay as the
**parity oracle**: CI builds them from their pinned sources through a deadline-enforcing subprocess
runner and compares verdicts and reasons label by label. The one deliberate verdict change is
phone (ii) := `CheckRegion` (ledger: `header-15px`). M3 rewrites orchestration §4.4 and QUICKSTART
§13 to the binary and pins which file `WORLD_MARKETING_RECTS` names (NB2).

| Plan milestone | Design milestone | What lands | Prod LOC | Test LOC |
|---|---|---|---|---|
| **M1a** | M1 (runner, §2.2 [r2]) | deadline runner + `TestMarketingSubprocessDeadline` | 0 | 131 |
| **M1b** | M1 | `position.go`, `main.go` `position` case, 34 fixtures, CLI + geometry tests, surface allow-list line | 229 + 4 | 150 + 3 |
| **M1c** | M1 | `TestMarketingPositionFixtures`, `TestMarketingPositionCorpus` (44 rows, 10 in-test arms) | 0 | 314 |
| **M1d** | M1 | `TestMarketingPositionParity` (oracle, pin registry) | 0 | 108 |
| **M2a** | M2 | `phone.go`, `main.go` `phone` case, 10 skeletons, skeleton/CLI/real-corpus tests, surface allow-list line | 224 + 3 | 341 + 2 |
| **M2b** | M2 | 39 arm ports, `TestMarketingPhoneArms`, `TestMarketingPhoneParity`, real-corpus parity | 0 | 404 |
| **M3** | M3 | QUICKSTART §13, orchestration §4.4/§9, 3 docs subtests | docs ~45 | 168 |

M1b and M2a exceed the ~250-LOC guidance: `position.go` (229) and `phone.go` (224) are ports of
single banked programs and cannot be split usefully, and each lands with the smallest test set that
exercises it (D-1). One file outside design §6 is touched: `host/store/context_roots_test.go`
gains two allow-list lines, because its `TestProductionGoSurface` census reds on any new non-test
Go file outside `host/`/`cmd/` (measured on the prototype: `non-test Go file outside context-root
scan and explicit reproducer allow-list: scripts/marketingcheck/position.go`, D-15).

---

## 1. Environment (on EVERY command)

```sh
cd /Users/voightkampff/dev/sunholo-data/.wt-world-iter328-row231
export PATH=/opt/homebrew/bin:$PATH
export AILANG_BIN=$HOME/.pinned-ailang/ailang     # prints "AILANG v0.41.0" (measured)
BANK=design_docs/verification/world-iter328
MC=scripts/marketingcheck
SCR=/Users/voightkampff/dev/sunholo-data/.exec-w328   # executor scratch: a SIBLING of the repo, never /tmp, never inside the worktree
CORPUS=$HOME/.ailang/state/world-iter328-phone-corpus
mkdir -p "$SCR"
```

| Pin | Measured | Why |
|---|---|---|
| `go` (in the repo) | `go1.26.6 darwin/arm64` (GOTOOLCHAIN auto) | the CI pin |
| `go` (outside a module) | `go1.26.4` | builds the oracles in file mode in `t.TempDir()`; both build them (design V26) |
| `AILANG_BIN` | `AILANG v0.41.0` | doc tests and `verify_ail.sh` |
| phone corpus | 11 PNGs + `SHA256SUMS` in `$CORPUS` | M2a skeleton generation and the real-corpus acceptance |

**Rules (each one is load-bearing).**
1. **Shell.** The shell is zsh. Run loops and multi-line scripts with `bash`.
2. **Instrument rule.** Never `cmd | tail; echo rc=$?`. Redirect to a file, put `echo "rc=$?"`
   after `;` on the same line, then read the file. A `-run` that matches nothing exits 0, so always
   count the `--- PASS` lines you expect.
3. **No `nohup`, no SIGHUP-ignoring wrapper.** Measured: run under `nohup`, the existing
   `TestMarketingSupervisorFakeDaemonCleanup/shell-hup` fails (`hangup trap must exit 129 … signal:
   killed`, 35 s), because a shell cannot trap a signal ignored at entry. Run tests in the
   foreground, or with the Bash tool's own background mode.
4. **One package-scale `go test` at a time.** Never run two at once.
5. **Copy, then verify.** Every banked file is copied with `cp` and its sha256 compared with §2's
   table. A mismatch is a STOP: report it, do not edit anything to make it match. Never hand-edit a
   Go file, a fixture or a doc; never run `gofmt -w` on banked code.
6. **Base check.** Before M1b, M2a and M3 overwrite an existing file, its sha256 must equal the base
   value in §2. A mismatch is a STOP (the base moved).
7. **Git.** Commit once per milestone with `git add` of the named files only (never `git add -A`
   or `.`). Do not push, do not open a PR, do not amend, do not rebase.
8. **Do not edit** `design_docs/world-mission*.md`, the design doc, this plan, the JSON, any
   `plan-*.txt`, `tools/launchd/*`, `.github/`, `go.mod`; do not move any doc; do not touch the
   main checkout, `~/.ailang/state/mission-v1*`, or the corpus directory (read only).
9. **Evidence** goes to `$BANK/executor-<milestone>-<what>.txt` (new files only), committed with
   the milestone it belongs to.
10. **Timeouts.** Every `go test` line below carries `-timeout 20m` (iteration 327 lost a run to
    Go's default 10-minute package timeout under rig load).

---

## 2. Banked sources and base pins

Base sha256 (the files M1b/M2a/M3 overwrite or edit):

| File | sha256 at `08dc31be` |
|---|---|
| `scripts/marketingcheck/main.go` | `bfc99a8da3e0bd97a38b1d30193ecc2fa58b5319b433deaffd7fb77f7fa8feef` |
| `scripts/marketingcheck/docs_test.go` | `17c9b520fd3c99f61a7a656814b7732c780528eb100ab78a967dde730022fed8` |
| `docs/QUICKSTART.md` | `64b05ab526c939671b5c33920fa31ef2d296ea2b09a6d9526beb988b7c3ca386` |
| `design_docs/planned/w-marketing-capture-orchestration.md` | `a57995e638f108ea503ad40fe2583e80632527cdb6f7d7673c5bf2abfd3769c1` |
| `host/store/context_roots_test.go` | `8a63b34104171a88417f2d19441c1e4ff0df53ebff9dc6d572881f774f41c6b3` |

Banked sources (`$BANK/…`), each validated green in its milestone's cumulative tree (§9):

| Banked file | Milestone | Destination / role | Lines | sha256 |
|---|---|---|---|---|
| `plan-M1a-position_test.go.txt` | M1a | `scripts/marketingcheck/position_test.go` | 131 | `ef839701fcc9cc6489182eafe34c59b4440aa67288a2460b097d432411d7cacc` |
| `plan-M1b-main.go.txt` | M1b | `scripts/marketingcheck/main.go` | 962 | `2f9864e0776b0829b832f5fa72e83f4983fef58f6fd9c71143510c8b363992b9` |
| `plan-M1b-position.go.txt` | M1b | `scripts/marketingcheck/position.go` | 229 | `7512337137e3eb08fc7930247c89ee1e6b1f392fe9014069823b22aabfaac6ac` |
| `plan-M1b-position_test.go.txt` | M1b | `scripts/marketingcheck/position_test.go` | 281 | `6002ade7f023e13355675d0e13b1a266d6a795277fde47e3836ae2bbd2b65ab3` |
| `plan-M1b-context_roots_test.go.txt` | M1b | `host/store/context_roots_test.go` | 261 | `6365fa4f98f6977bb002e41b3400cf695a4c0748d642f11d003771eaf11227b8` |
| `plan-gen-position-fixtures.sh.txt` | M1b | generator (run, not copied) | 37 | `41ed23322bf58acaa7bb430b737f05968a722a93471dbb9e30bdb37f7cd9b4c7` |
| `plan-expected-position-SHA256SUMS.txt` | M1b | expected `testdata/position/SHA256SUMS` (cmp) | 41 | `d8963a3c2245d491c18883da68c1f74eb1f0dad9c79503af324fa100352d300a` |
| `plan-M1c-position_test.go.txt` | M1c | `scripts/marketingcheck/position_test.go` | 595 | `33126f459e979ea1220e3dac60937425259bc840d06eeb029903021f63a956c9` |
| `plan-M1d-position_test.go.txt` | M1d | `scripts/marketingcheck/position_test.go` (final) | 703 | `8e86c2563e24ab03dc188d4818081023c7a65f73dd97d8bd188b94507e328a0d` |
| `plan-M2a-main.go.txt` | M2a | `scripts/marketingcheck/main.go` (final) | 965 | `24448bb78265726f66085833b8b61ec5b28f16f9743cf6926a5c12e4cd189268` |
| `plan-M2a-phone.go.txt` | M2a | `scripts/marketingcheck/phone.go` | 224 | `e3b8db55c53cf6d69c31c0f908b8302b81dc24277e80af6b887473358f28c4b5` |
| `plan-M2a-phone_test.go.txt` | M2a | `scripts/marketingcheck/phone_test.go` | 341 | `e91ec24c18bd9416d91339e288f1e6da5f34a7a33e77b8f8e23c0b5e0e4b47a5` |
| `plan-M2a-context_roots_test.go.txt` | M2a | `host/store/context_roots_test.go` (final) | 263 | `b7ed3067a862ecc0402eb791804e138443a04e8c3b6af5711548f30b029106cc` |
| `plan-expected-phone-sha256.txt` | M2a | expected skeleton + corpus-SHA256SUMS sha list (diff) | 11 | `7eb9e00945d500a66d274fd927143c41f739ab37a5076264f9aa6de9d75f8b20` |
| `plan-M2b-phone_test.go.txt` | M2b | `scripts/marketingcheck/phone_test.go` (final) | 745 | `0bce70c4f72784d4fea1437abd636c76e90874301b2fc97453d8e006cf0f71a8` |
| `plan-m3-apply-docs.py.txt` | M3 | doc-edit script (run, not copied) | 129 | `0681ce7ed4e26a40be7a4768c0b309fcf33daf6a971dfb2bd989973a27761529` |
| `plan-M3-docs_test.go.txt` | M3 | `scripts/marketingcheck/docs_test.go` | 684 | `b7420a6d8dfca7652efda94caf8d846bac3b89f28dbc159ebd93aafa407200a7` |
| `plan-m3-s7-runbook.sh.txt` | M3 | §13 lines verbatim (run, not copied) | 18 | `f7935fd680a753519b5f95db3b5831eee1e5ff9052af4faa174468fa854dfcd8` |
| `plan-mutants.py.txt` | §8 | AC-4 drill (run in a throwaway copy) | 105 | `1f9b6072a15bf9d53f5310fdebdd4b6c38cdc8a1b3adcf707bb549bf2430e8cf` |
| `plan-mutants-planner-run.txt` | §8 | the planner's drill output (reference) | 72 | `f19eaeeb88a257f1dc5323353e2eb259265aa174967a9aec32286e1c56fe70fb` |
| `plan-fidelity-test.go.txt` | §9 | planner instrument (not copied) | 55 | `d6ebed50668e7cb77a195a21390b060231d30e49e7434b9b65d73be747c56a5c` |
| `plan-validation-run.txt` | §9 | the planner's validation logs (reference) | 221 | `1533167c5de5191d35c16291d63c9907e4a027687a866b69a640baf3bf8bc9d5` |

`plan-mutants.py.txt` and `plan-fidelity-test.go.txt` are instruments: the first is the AC-4 drill
(§8), the second is the planner's one-off proof that the in-test phone arm ports are pixel-equal to
the banked generators (§9.4). Neither is copied into the package.

---

## 3. Milestone M1a — the deadline-enforcing subprocess runner

**Files:** `scripts/marketingcheck/position_test.go` (new; the runner part only).

What it is (design §2.2 [r2]): `boundedRun(limit, dir, env, name, args...)` runs a child with
`exec.CommandContext` under `context.WithTimeout`, in its own process group (`Setpgid`), kills the
group on expiry (`cmd.Cancel`), bounds output draining with `cmd.WaitDelay = 5s`, reaps the group
after every run, and returns `subprocess-timeout: <command> exceeded <limit>`. Limits:
`subprocBuildLimit = 120s`, `subprocRunLimit = 30s`, `subprocDrainLimit = 5s`. `mustRun` fails the
test on any such error; it never skips or substitutes a result. `TestMarketingSubprocessDeadline`
is the runner's own control (P17): a fake checker (`hang.sh`: `sleep 20 &`, `wait`) under a 2 s
limit must fail with `subprocess-timeout` naming the script, return within 2 + 5 + 3 s, and leave
its grandchild dead (checked with `kill -0` and `ps -o stat=`, zombies count as dead).

```sh
cp $BANK/plan-M1a-position_test.go.txt $MC/position_test.go
shasum -a 256 $MC/position_test.go     # must equal the M1a row of §2
```

Acceptance (expected output measured on the prototype, §9):

```sh
go test -count=1 -timeout 20m -v -run '^TestMarketingSubprocessDeadline$' ./$MC > $SCR/m1a-acc.log 2>&1; echo "rc=$?"
grep -c '^--- PASS: TestMarketingSubprocessDeadline' $SCR/m1a-acc.log      # 1
grep 'subprocess-timeout after' $SCR/m1a-acc.log                           # "subprocess-timeout after 2.00xs; grandchild <pid> reaped"
```

Then the milestone gate (§10.1). Mutants owned (run in the final drill, §8): **P17**, **P17b**.

Commit:

```sh
cp $SCR/m1a-acc.log $BANK/executor-M1a-acceptance.txt
git add $MC/position_test.go $BANK/executor-M1a-acceptance.txt
git commit -m "test(row 231 M1a): deadline-enforcing subprocess runner for marketingcheck oracle tests"
```

---

## 4. Milestone M1b — `marketingcheck position`, fixtures, CLI and geometry tests

**Files:** `scripts/marketingcheck/position.go` (new), `scripts/marketingcheck/main.go` (the
`position` case and usage line), `scripts/marketingcheck/position_test.go` (grows),
`scripts/marketingcheck/testdata/position/**` (34 JSON + `SHA256SUMS`, generated),
`host/store/context_roots_test.go` (allow-list line for `position.go`, D-15).

Surface (design §2.1): `marketingcheck position --baseline RECTS LABEL=RECTS...`. Functions:
`PositionJudge(base, cand Rectangles) (reasons []string, regions int)` (the oracle's r3 `judge`,
reason strings byte-identical), `loadPositionCandidate(path) (Rectangles, error)`
(`bad-rectangles` / `zero-regions`), `parseLabelled(args) ([][2]string, error)`,
`printRefusal(stderr, err)`, `runPosition(args []string, stdout, stderr io.Writer) int`. The
baseline goes through `LoadRectangles`. Output: `%-38s ACCEPT|REFUSE regions=17`, one indented
line per reason, `TOTAL candidates=N accept=A refuse=R`. rc 0 all ACCEPT, 1 any REFUSE or a refused
input (stderr `marketingcheck: <code> <name>: <detail>`), 2 usage, `--help` 0.

```sh
test "$(shasum -a 256 $MC/main.go | cut -d' ' -f1)" = bfc99a8da3e0bd97a38b1d30193ecc2fa58b5319b433deaffd7fb77f7fa8feef && echo base-ok
test "$(shasum -a 256 host/store/context_roots_test.go | cut -d' ' -f1)" = 8a63b34104171a88417f2d19441c1e4ff0df53ebff9dc6d572881f774f41c6b3 && echo store-base-ok
cp $BANK/plan-M1b-main.go.txt $MC/main.go
cp $BANK/plan-M1b-position.go.txt $MC/position.go
cp $BANK/plan-M1b-position_test.go.txt $MC/position_test.go
cp $BANK/plan-M1b-context_roots_test.go.txt host/store/context_roots_test.go
shasum -a 256 $MC/main.go $MC/position.go $MC/position_test.go host/store/context_roots_test.go    # must equal the M1b rows of §2
```

### M1b.1 Fixtures — generated, never typed (design §2.3)

```sh
bash $BANK/plan-gen-position-fixtures.sh.txt "$PWD" "$SCR/gen" > $SCR/m1b-gen.log 2>&1; echo "rc=$?"   # rc=0, log "files=34"
cmp $MC/testdata/position/SHA256SUMS $BANK/plan-expected-position-SHA256SUMS.txt; echo "cmp rc=$?"    # cmp rc=0
```

The generator (read it; do not edit it) checks the oracle's pin, builds the banked `p0check` in
`$SCR/gen`, and writes:

| Files | Produced by | Pin |
|---|---|---|
| `stores/<store>.json` ×9 (`r1-seed1`, `r1-seed2`, `r1b-seed3`, `r1b-seed4`, `r1c-seed5`, `r1c-seed6`, `eval-p1`, `eval-p2`, `eval-p3`) | `cp world-iter315/v64-rects-<store>.json` | each sha equals the store's `rects.json` line of `world-iter315/v64-inputs-sha256.txt` (asserted by M1c) |
| `stores/r1-seed{1,2}-blockconverted.json` | `p0check fromlayouts v64-layouts-<s>.json v64-rects-<s>.json <out>`, then `cmp` to the banked files | byte-equal to `world-iter315/v64-rects-<s>-blockconverted.json` |
| `v65/<arm>.json` ×19 | `p0check mutate v64-rects-r1c-seed5.json v65/` | all 19 sha equal the V23 seed-5 list (measured) |
| `v65-p1/<arm>.json` ×4 (`neg-tested-at-attested-row`, `neg-packet-all+14`, `syn-pos-tested-wraps-3-lines`, `neg-proven-grows-6-lines-below-fold`) | `p0check mutate v64-rects-eval-p1.json <scratch>`, 4 copied | equal the V23 eval-p1 list (measured) |
| `SHA256SUMS` | the generator: 7 `#` comment lines naming the commands, then `shasum -a 256` of the 34 files in `LC_ALL=C` order | sha256 `d8963a3c2245d491c18883da68c1f74eb1f0dad9c79503af324fa100352d300a` |

### M1b.2 Acceptance

```sh
go test -count=1 -timeout 20m -v -run '^TestMarketingPosition' ./$MC > $SCR/m1b-acc.log 2>&1; echo "rc=$?"
grep -c '^--- PASS' $SCR/m1b-acc.log     # 2 (TestMarketingPositionCLI, TestMarketingPositionGeometryMatchesOracle)
grep -c -- '--- FAIL' $SCR/m1b-acc.log   # 0
go run ./$MC position --help > $SCR/m1b-help.log 2>&1; echo "rc=$?"     # rc=0; first line "usage: marketingcheck position --baseline RECTS LABEL=RECTS..."
```

`TestMarketingPositionCLI` subtests (16): `accept-rc`, `refuse-rc`, `baseline-empty-content`,
`candidate-zero-regions`, `candidate-wrong-schema`, `candidate-unparsable`, `candidate-missing-file`,
`no-candidates`, `missing-baseline`, `bad-token`, `duplicate-label`, `empty-label`, `empty-path`,
`bad-flag`, `help`, `dispatch` (through `run`, rc 2). These are design §2.4's position cases.
`TestMarketingPositionGeometryMatchesOracle` is the design's V32 prototype unchanged.

Gate §10.1 (including its `host/store` line). Mutants owned: **P9, P10, P11, P14, P15, G1**
(killers in §8). Commit:

```sh
cp $SCR/m1b-acc.log $BANK/executor-M1b-acceptance.txt
git add $MC/main.go $MC/position.go $MC/position_test.go $MC/testdata/position host/store/context_roots_test.go $BANK/executor-M1b-acceptance.txt
git commit -m "feat(row 231 M1b): marketingcheck position (P0/P3 r3 predicate) with V64/V65 fixtures"
```

---

## 5. Milestone M1c — fixture pins and the 44-row corpus

**Files:** `scripts/marketingcheck/position_test.go` (grows).

- `TestMarketingPositionFixtures`: every file under `testdata/position` is listed in
  `SHA256SUMS` with its sha (34 lines, 34 files); the 9 stores' sha appear in
  `v64-inputs-sha256.txt` as `<sha>  /tmp/world-iter312/<store>/rects.json`; the two block-converted
  files are byte-equal to the banked ones; the baseline is `ba7c83ce…3333`; and (D-6) the ten in-test
  arms are byte-equal to the banked V13 generator's output (hardcoded V13 sha256).
- `TestMarketingPositionCorpus`: one subtest per label over `positionInputs` (44 =
  V64 11 + V65 19 + V65-p1 4 + new 10), each asserting `regions=17` and the verdict string
  `ACCEPT` or `REFUSE <sorted clause tags>` against the table `positionWant` (19 ACCEPT). The 10
  new arms are `positionNewArms`, a Go port of `world-iter328/v13-isoarms-main.go.txt`, applied to
  `stores/r1c-seed5.json` and written to `t.TempDir()`.

```sh
cp $BANK/plan-M1c-position_test.go.txt $MC/position_test.go
shasum -a 256 $MC/position_test.go
go test -count=1 -timeout 20m -v -run '^TestMarketingPosition' ./$MC > $SCR/m1c-acc.log 2>&1; echo "rc=$?"
grep -c '^--- PASS' $SCR/m1c-acc.log                                  # 4
grep -c '^    --- PASS: TestMarketingPositionCorpus/' $SCR/m1c-acc.log # 44
grep -c -- '--- FAIL' $SCR/m1c-acc.log                                # 0
```

Gate §10.1. Mutants owned: **P1–P8, P13**. Commit:

```sh
cp $SCR/m1c-acc.log $BANK/executor-M1c-acceptance.txt
git add $MC/position_test.go $BANK/executor-M1c-acceptance.txt
git commit -m "test(row 231 M1c): position fixture pins and the 44-row V64/V65/iso-arm corpus"
```

---

## 6. Milestone M1d — position parity with the banked oracle

**Files:** `scripts/marketingcheck/position_test.go` (final).

`TestMarketingPositionParity`: subtest `pin-registry` asserts the orchestration doc carries the
hardcoded pin `0a9a6795…927c` exactly once (P16); `buildOracle` reads
`world-iter315/v64-p0check-main.go.txt`, checks its sha256 against the hardcoded pin, writes it to
`t.TempDir()/oracle.go` and runs `go build -o oracle oracle.go` there under the 120 s deadline; the
oracle runs `check <R1> label=path…` over the 44 inputs under the 30 s deadline; per label the r3
verdict and the full r3 reason list must equal `PositionJudge`'s. Ledger: empty. It asserts 44
labels compared and logs `PARITY position: inputs=44 verdict-diffs=0 reason-diffs=0`.

```sh
cp $BANK/plan-M1d-position_test.go.txt $MC/position_test.go
shasum -a 256 $MC/position_test.go
go test -count=1 -timeout 20m -v -run '^TestMarketingPosition' ./$MC > $SCR/m1d-acc.log 2>&1; echo "rc=$?"
grep -c '^--- PASS' $SCR/m1d-acc.log            # 5  (design §3 M1 acceptance)
grep -c '^--- SKIP' $SCR/m1d-acc.log            # 0
grep -c -- '--- FAIL' $SCR/m1d-acc.log          # 0
grep 'PARITY position' $SCR/m1d-acc.log         # PARITY position: inputs=44 verdict-diffs=0 reason-diffs=0
go test -count=1 -timeout 20m -v -run '^TestMarketingSubprocessDeadline$' ./$MC > $SCR/m1d-runner.log 2>&1; echo "rc=$?"   # 1 PASS
```

Gate §10.1. Mutants owned: **P12, P16 (position half), O1**. Commit:

```sh
cat $SCR/m1d-acc.log $SCR/m1d-runner.log > $BANK/executor-M1d-acceptance.txt
git add $MC/position_test.go $BANK/executor-M1d-acceptance.txt
git commit -m "test(row 231 M1d): position parity against the pinned banked p0check oracle"
```

---

## 7. Milestones M2a, M2b, M3

### 7.1 M2a — `marketingcheck phone`, the 10 skeletons, skeleton/CLI/real-corpus tests

**Files:** `scripts/marketingcheck/phone.go` (new), `scripts/marketingcheck/main.go` (final:
adds the `phone` case and usage line), `scripts/marketingcheck/phone_test.go` (new, M2a part),
`scripts/marketingcheck/testdata/phone/corpus-SHA256SUMS` and `skeleton-*.json` ×10 (generated),
`host/store/context_roots_test.go` (allow-list line for `phone.go`, D-15).

`phone.go`: the V58 scan functions copied unchanged (`phoneLuma`, `phoneBands`, `phoneScrollbar`,
`phoneTitleRight`), `ScanPhone(im) PhoneScan`, `PhoneJudge(im, s) []string` ((i) ±4 vs R1,
(ii) = `CheckRegion` on header and first entry, (iii′) dims/thumb/title-right, plus
`zero-bands: no ink rows read`), `runPhone(args, stdout, stderr) int` (`decode-failed` rc 1).
`phone_test.go` (M2a part): `phoneSkeleton` (schema `world/marketingcheck/phone-skeleton/v1`),
`renderPhoneSkeleton`, `derivePhoneSkeleton`, `TestMarketingPhoneSkeleton` (10 skeletons: render
sha, lossless PNG round trip, scan == recorded geometry, header/entry == R1, ACCEPT),
`TestMarketingPhoneCLI` (12 subtests: `accept-rc`, `refuse-rc`, `blank` (zero-bands), `decode-failed`,
`missing-file`, `no-candidates`, `bad-token`, `duplicate-label`, `empty-label`, `bad-flag`, `help`,
`dispatch`), and `TestMarketingPhoneRealCorpus` without its arms/parity tail (skips without
`WORLD_PHONE_CORPUS`; with it: all 11 hashes vs `corpus-SHA256SUMS` or `BLOCKED`, the 10 fresh
skeletons byte-equal to the committed ones, each image ACCEPT with header/entry == R1 and its thumb
X1, the desktop shot REFUSE).

```sh
test "$(shasum -a 256 $MC/main.go | cut -d' ' -f1)" = 2f9864e0776b0829b832f5fa72e83f4983fef58f6fd9c71143510c8b363992b9 && echo m1b-main-ok
cp $BANK/plan-M2a-main.go.txt $MC/main.go
cp $BANK/plan-M2a-phone.go.txt $MC/phone.go
cp $BANK/plan-M2a-phone_test.go.txt $MC/phone_test.go
cp $BANK/plan-M2a-context_roots_test.go.txt host/store/context_roots_test.go
shasum -a 256 $MC/main.go $MC/phone.go $MC/phone_test.go host/store/context_roots_test.go     # must equal the M2a rows of §2
```

Fixtures (design §2.3: written by the named command, checked in as written):

```sh
(cd "$CORPUS" && shasum -a 256 -c SHA256SUMS) > $SCR/m2a-corpus.log 2>&1; echo "rc=$?"   # rc=0, 11 lines "…: OK"
mkdir -p $MC/testdata/phone
cp -p "$CORPUS/SHA256SUMS" $MC/testdata/phone/corpus-SHA256SUMS
WORLD_PHONE_CORPUS=$HOME/.ailang/state/world-iter328-phone-corpus WORLD_PHONE_SKELETON_WRITE=1 go test -count=1 -timeout 20m -run '^TestMarketingPhoneRealCorpus$' ./scripts/marketingcheck > $SCR/m2a-write.log 2>&1; echo "rc=$?"   # rc=0
(cd $MC/testdata/phone && shasum -a 256 skeleton-*.json corpus-SHA256SUMS) > $SCR/m2a-phone-sha.txt; diff $SCR/m2a-phone-sha.txt $BANK/plan-expected-phone-sha256.txt; echo "diff rc=$?"   # diff rc=0
```

If the corpus is missing or any hash differs: **STOP** (design §2.5 step 1, BLOCKED; do not
regenerate anything). Expected `corpus-SHA256SUMS` content (11 lines, sha256 of the file
`476d971ab75bd33ccb7450405c43c5dac2c38cddeb08b7ac18f0236e489b2ddf`):

```text
746a020e6ad5e6a009a3149fe5b2f39de14d1bf18adbd5bd285718a47cb3bb1b  eval-p1-pass-a-phone.png
3793ef8b914a92030adbeb3b0fb548bb123accfb1b13f996769df06b7ff0b0fb  eval-p2-pass-a-phone.png
aa471bfa73da37f82870047b1d7ff963befa9055648f909708f64f0d46ddbbba  eval-p3-pass-a-phone.png
676d93d3bde3668fdeef2380763ff7714dba6ba14ec45665a6a0c079582f7bf4  r1-seed1-pass-a-phone.png
ff61fd2d3265825299c84f99028159065dc9787289730601279f5d892031fb78  r1-seed2-pass-a-phone.png
20a8699be3a156b133f18dfed968f9a5782c1313d4783742641bf8a3a457d52f  r1b-seed3-pass-a-phone.png
ef0f5680df915d4c2bc4738e0c68bcebd939e4a6b3244f40e6801426751edeb1  r1b-seed4-pass-a-phone.png
94b14d0f403130ac812a5cbebd22919d6d68a0548dbd0078f735b7008ae8cd1c  r1c-seed5-pass-a-phone.png
8c3c7493bb37a4f18561bb18307a3b9bd20826d3ec9b446280e817a171ed2594  r1c-seed5-pass-a-workbench-light.png
07195b928cdc6b840bdac5b1ab4fe1fb412d26918a8ed220e7e812a072d18065  r1c-seed6-pass-a-phone.png
999125ff43bdbe5857fc053aa35b55fdfa1a6440447830ee539a43e7f6994be8  v50-reference-phone.png
```

The 10 skeletons total 23649 bytes (2363–2366 each, equal to design V31). Their sha256 are in
`plan-expected-phone-sha256.txt`.

Acceptance:

```sh
go test -count=1 -timeout 20m -v -run '^TestMarketingPhone' ./$MC > $SCR/m2a-acc.log 2>&1; echo "rc=$?"
grep -c '^--- PASS' $SCR/m2a-acc.log    # 2 (Skeleton, CLI)
grep -c '^--- SKIP' $SCR/m2a-acc.log    # 1 (TestMarketingPhoneRealCorpus)
WORLD_PHONE_CORPUS=$CORPUS go test -count=1 -timeout 20m -v -run '^TestMarketingPhone' ./$MC > $SCR/m2a-acc-corpus.log 2>&1; echo "rc=$?"
grep -c '^--- PASS' $SCR/m2a-acc-corpus.log   # 3
grep -c '^--- SKIP' $SCR/m2a-acc-corpus.log   # 0
go run ./$MC phone --help > $SCR/m2a-help.log 2>&1; echo "rc=$?"     # rc=0, "usage: marketingcheck phone LABEL=PNG..."
```

Gate §10.1 (including its `host/store` line). Mutants owned: **F10, F11, G2** (and P10's phone
half). Commit:

```sh
cat $SCR/m2a-write.log $SCR/m2a-acc.log $SCR/m2a-acc-corpus.log > $BANK/executor-M2a-acceptance.txt
git add $MC/main.go $MC/phone.go $MC/phone_test.go $MC/testdata/phone host/store/context_roots_test.go $BANK/executor-M2a-acceptance.txt
git commit -m "feat(row 231 M2a): marketingcheck phone with 10 real-geometry skeleton fixtures"
```

### 7.2 M2b — 39 arms, phone parity, real-corpus parity

**Files:** `scripts/marketingcheck/phone_test.go` (final).

Adds `phoneArms` (Go ports: 12 of `world-iter315/v67-phoneneg-main.go.txt`, 23 of
`v67b-judge-phone-arms-main.go.txt`, and `header-solid-ink`, `entry-solid-ink`, `tall-845`
(`world-iter328/v14-newarms-main.go.txt`) and `header-15px` (`v11-h15-main.go.txt`)),
`phoneArmByName`, `phoneArmInputs`, `phoneRenderInputs` (seed-5 and seed-2 renders + 39 arms of the
seed-5 render = 41 PNGs), `phoneWant` (41 rows, 12 ACCEPT), `TestMarketingPhoneArms`,
`phoneParity` (oracle r3 verdict and its five clause booleans vs the phone reasons; ledger exactly
`{new/header-15px: oracle ACCEPT, phone REFUSE with "(ii) rect-too-small phone.png: header"}`; any
other difference red, and the ledger entry red if it stops differing), `TestMarketingPhoneParity`
(`pin-registry` subtest, 41 inputs), and the real-corpus tail (39 arms of the real seed-5 image
with `phoneWant`'s verdicts; parity over 59 inputs: V50, 9 stores × A/B, desktop, 39 arms).

```sh
cp $BANK/plan-M2b-phone_test.go.txt $MC/phone_test.go
shasum -a 256 $MC/phone_test.go
go test -count=1 -timeout 20m -v -run '^TestMarketingPhone' ./$MC > $SCR/m2b-acc.log 2>&1; echo "rc=$?"
grep -c '^--- PASS' $SCR/m2b-acc.log    # 4
grep -c '^--- SKIP' $SCR/m2b-acc.log    # 1 (TestMarketingPhoneRealCorpus)
grep 'PARITY phone' $SCR/m2b-acc.log    # PARITY phone: inputs=41 diffs=0 ledger=1
WORLD_PHONE_CORPUS=$CORPUS go test -count=1 -timeout 20m -v -run '^TestMarketingPhone' ./$MC > $SCR/m2b-acc-corpus.log 2>&1; echo "rc=$?"
grep -c '^--- PASS' $SCR/m2b-acc-corpus.log   # 5
grep -c '^--- SKIP' $SCR/m2b-acc-corpus.log   # 0
grep -E 'PARITY phone|CHECK images' $SCR/m2b-acc-corpus.log
#   PARITY phone: inputs=41 diffs=0 ledger=1
#   PARITY phone: inputs=59 diffs=0 ledger=1
#   CHECK images=10 byte-equal=10; real arms=39; parity inputs=59
```

The design asks the executor to bank the with-corpus output: it is `executor-M2b-acceptance.txt`.
Gate §10.1. Mutants owned: **F1–F9, P16 (phone half), O2**. Commit:

```sh
cat $SCR/m2b-acc.log $SCR/m2b-acc-corpus.log > $BANK/executor-M2b-acceptance.txt
git add $MC/phone_test.go $BANK/executor-M2b-acceptance.txt
git commit -m "test(row 231 M2b): phone arms, oracle parity (ledger header-15px) and real-corpus parity"
```

### 7.3 M3 — documents and NB2

**Files:** `docs/QUICKSTART.md` (§13 only), `design_docs/planned/w-marketing-capture-orchestration.md`
(§4.4 and §9), `scripts/marketingcheck/docs_test.go`.

The doc edits are design §2.7, applied by a script whose every replacement must match exactly once
(else it exits 1 and writes nothing). Read it before running it.

```sh
test "$(shasum -a 256 docs/QUICKSTART.md | cut -d' ' -f1)" = 64b05ab526c939671b5c33920fa31ef2d296ea2b09a6d9526beb988b7c3ca386 && echo q-base-ok
test "$(shasum -a 256 design_docs/planned/w-marketing-capture-orchestration.md | cut -d' ' -f1)" = a57995e638f108ea503ad40fe2583e80632527cdb6f7d7673c5bf2abfd3769c1 && echo o-base-ok
test "$(shasum -a 256 $MC/docs_test.go | cut -d' ' -f1)" = 17c9b520fd3c99f61a7a656814b7732c780528eb100ab78a967dde730022fed8 && echo d-base-ok
python3 -I $BANK/plan-m3-apply-docs.py.txt "$PWD" > $SCR/m3-apply.log 2>&1; echo "rc=$?"   # rc=0, "applied: QUICKSTART 3 edits, orchestration 9 edits"
cp $BANK/plan-M3-docs_test.go.txt $MC/docs_test.go
shasum -a 256 docs/QUICKSTART.md design_docs/planned/w-marketing-capture-orchestration.md $MC/docs_test.go
#   92378af24c73f46107d1f5ae7a42667b74d7fb3d312eea8a528857609749133e  docs/QUICKSTART.md
#   0b2ee85987beac8d95b82d6a36a61e28b6b8f269dc3d3b0ff100817091db0e9c  design_docs/planned/w-marketing-capture-orchestration.md
#   b7420a6d8dfca7652efda94caf8d846bac3b89f28dbc159ebd93aafa407200a7  scripts/marketingcheck/docs_test.go
```

What the script changes (design §2.7): orchestration `:220` (P0 item (1) gets the probe command
and the 4/0 reading rule), `:223` (names), `:235` (final sentence), `:236`, `:237`, `:238` (the
five-step capture-day control step, V50 and seed-5 hashes in full, ending with the
`Parity-oracle sources (row 231): …` pin line), `:239` (`Committed (row 231).`), `:248` (P3 lines),
and a new `### Row 231` in §9. QUICKSTART §13: a P3 paragraph and bash block with three `position`
lines after the inputs block, a phone paragraph and bash block with two `phone` lines after the
capture block, and the `.txt`-checker sentence at old :1720–1723 replaced with the exit-code
sentence plus one sentence on `go run`'s exit-code collapse (D-11). `usageText` already carries both usage lines (M1b, M2a).

The three new subtests of `TestMarketingRunbookMatchesHelp`:
- `rects-value`: exactly one `export WORLD_MARKETING_RECTS=$DEMO_DIR/rects.json`; the capture line
  passes `--rectangles "$WORLD_MARKETING_RECTS"`; the one `position … release=` line passes
  `release="$WORLD_MARKETING_RECTS"` and its `--baseline` is the R1 path, which is not the candidate.
- `position-lines`: builds `marketingcheck` once (120 s deadline), runs the 3 documented position
  lines from the repo root (`../..`, 30 s each); `accept*` → ACCEPT, `refuse*` → REFUSE with the
  matching rc; the release line runs twice with `$WORLD_MARKETING_RECTS` (expanded from §13's
  export, which must lie under `$DEMO_DIR`, D-8) holding `stores/eval-p1.json` (rc 0, ACCEPT) and
  then `v65/neg-attested+14.json` (rc 1, REFUSE).
- `phone-lines`: writes the seed-5 render as `$CAPTURE_OUT/pass-a/phone.png` and the `tall-845` arm
  as `$CAPTURE_OUT/pass-a/workbench-light.png`, runs the 2 documented phone lines: `refuse` →
  REFUSE rc 1, `release` → ACCEPT rc 0.

Acceptance:

```sh
go test -count=1 -timeout 20m -v -run 'TestMarketingRunbookMatchesHelp|TestMarketingDocsLinks' ./$MC > $SCR/m3-acc.log 2>&1; echo "rc=$?"
grep -c -- '--- FAIL' $SCR/m3-acc.log    # 0
grep -E '^    --- PASS: TestMarketingRunbookMatchesHelp/(rects-value|position-lines|phone-lines)' $SCR/m3-acc.log   # 3 lines
grep -c 0a9a6795d9d711f99d77030e0491cd1c5e254c61fc7434898a80b4875404927c design_docs/planned/w-marketing-capture-orchestration.md   # 1  (V34)
grep -c 31576922bc566ff9d6c2c40b4cfa41e71f234fc706d013377f657a0ffd010b4c design_docs/planned/w-marketing-capture-orchestration.md   # 1  (V34)
go test -count=1 -timeout 20m -v -run '^TestMarketing(Position|Phone)Parity$' ./$MC > $SCR/m3-parity.log 2>&1; echo "rc=$?"   # rc=0, 2 PASS (V34: parity on the edited doc)
```

Run §13's new lines verbatim (S7), from the repo root, and bank the output. The script holds the
five documented lines byte for byte and sets only the variables §13 sets; it feeds the real
seed-5 phone and desktop shots as the capture.

```sh
bash $BANK/plan-m3-s7-runbook.sh.txt "$SCR" "$CORPUS" > $SCR/m3-s7.log 2>&1; echo "rc=$?"
grep '^rc=' $SCR/m3-s7.log      # rc=0, rc=1, rc=0, rc=1, rc=0  (in that order)
```

Order: gate §10.1 → the M3 commit below → repo gate §10.2 → mutation drill §8 (it archives the
committed HEAD) → the evidence commit. Mutants owned: **D1–D3, P16**. M3 commit:

```sh
cp $SCR/m3-acc.log $BANK/executor-M3-acceptance.txt
cat $SCR/m3-parity.log > $BANK/executor-V34-pins-after-M3.txt
grep -c 0a9a6795d9d711f99d77030e0491cd1c5e254c61fc7434898a80b4875404927c design_docs/planned/w-marketing-capture-orchestration.md >> $BANK/executor-V34-pins-after-M3.txt
grep -c 31576922bc566ff9d6c2c40b4cfa41e71f234fc706d013377f657a0ffd010b4c design_docs/planned/w-marketing-capture-orchestration.md >> $BANK/executor-V34-pins-after-M3.txt
cp $SCR/m3-s7.log $BANK/executor-M3-s7-runbook.txt
git add docs/QUICKSTART.md design_docs/planned/w-marketing-capture-orchestration.md $MC/docs_test.go $BANK/executor-M3-acceptance.txt $BANK/executor-V34-pins-after-M3.txt $BANK/executor-M3-s7-runbook.txt
git commit -m "docs(row 231 M3): orchestration §4.4 and QUICKSTART §13 switch to marketingcheck position/phone; NB2 pinned"
```

Then run §10.2 and §8, and commit their evidence as one more commit:

```sh
git add $BANK/executor-AC4-mutants.txt $BANK/executor-V35-hanging-control.txt $BANK/executor-repo-gate.txt $BANK/executor-capture-day-step1.txt
git commit -m "evidence(row 231): AC-4 mutation drill, V35 hanging control, repo gate"
```

---

## 8. AC-4 mutation drill (run once, after the M3 commit, in a throwaway copy)

`plan-mutants.py.txt` applies each mutant alone to a copy of the committed tree, runs the named
tests with `-timeout 20m`, records every red leaf (sub)test, restores the file and asserts it is
byte-equal. It never touches the worktree.

```sh
rm -rf $SCR/mut && mkdir -p $SCR/mut && git archive HEAD | tar -x -C $SCR/mut
python3 -I $BANK/plan-mutants.py.txt $SCR/mut > $SCR/mutants.txt 2>&1; echo "rc=$?"
grep -c ': KILLED ' $SCR/mutants.txt                          # 36
grep -cE ': (SURVIVED|BUILD-FAILED|SETUP-ERROR)' $SCR/mutants.txt   # 0
cp $SCR/mutants.txt $BANK/executor-AC4-mutants.txt
grep -A1 '^P17' $SCR/mutants.txt > $BANK/executor-V35-hanging-control.txt
rm -rf $SCR/mut
```

Expected (measured on the planner's prototype of the same final tree; "red" counts leaf
(sub)tests; P14/P15 also list the pre-existing `Shots` tests where they red):

| ID | Mutant | Result | Red | Time | Red (sub)tests (`Pos…` = `TestMarketingPosition…`, `Ph…` = `TestMarketingPhone…`) |
|---|---|---|---|---|---|
| P1 | `(a1)` check removed | KILLED | 6 | 4.8 s | PosCorpus/stores/r1-seed1, PosCorpus/stores/r1-seed2, PosCorpus/v65/neg-light-header-whole-view, PosCorpus/v65/neg-tested-64x16-glyphs, PosCorpus/new/iso-a1-tested-w-3, PosParity |
| P2 | `c.H >= 24` removed | KILLED | 2 | 4.6 s | PosCorpus/new/iso-minh-approvals-h-1, PosParity |
| P3 | non-paragraph `abs(dh) > 2` removed | KILLED | 5 | 4.6 s | PosCorpus/v65/neg-graph-svg-at-timeline, PosCorpus/v65/neg-light-header-whole-view, PosCorpus/new/iso-nonpara-graph-svg-h+24, PosCorpus/new/iso-nonpara-payload-h+24, PosParity |
| P4 | free-gap `gc >= 0` removed | KILLED | 2 | 4.6 s | PosCorpus/new/iso-freegap-validation-h+240, PosParity |
| P5 | `(a2)` removed | KILLED | 10 | 4.6 s | PosCorpus/stores/r1-seed1, PosCorpus/stores/r1-seed2, PosCorpus/v65/neg-graph-svg-at-timeline, PosCorpus/v65/neg-light-header+40, PosCorpus/v65/neg-light-header-whole-view, … (+5) |
| P6 | `(a3)` removed | KILLED | 14 | 5.4 s | PosCorpus/stores/r1-seed1, PosCorpus/stores/r1-seed2, PosCorpus/v65/neg-attested+14, PosCorpus/v65/neg-dark-timeline-first+24, PosCorpus/v65/neg-light-header+40, … (+9) |
| P7 | content `(b)` removed | KILLED | 4 | 5.0 s | PosCorpus/v65/neg-light-header-whole-view, PosCorpus/v65/neg-packet-stack-down-6-lines, PosCorpus/v65/neg-proven-grows-6-lines-below-fold, PosParity |
| P8 | background edge check removed | KILLED | 2 | 4.8 s | PosCorpus/v65/neg-background-over-content, PosParity |
| P9 | baseline via plain `json.Unmarshal` | KILLED | 1 | 5.5 s | PosCLI/baseline-empty-content |
| P10 | zero-candidate guard removed | KILLED | 3 | 4.9 s | PhCLI/no-candidates, PhCLI/dispatch, PosCLI/no-candidates |
| P11 | REFUSE `return 1` → `return 0` | KILLED | 1 | 4.7 s | PosCLI/refuse-rc |
| P12 | `(a3) gap after` → `(a3) gap` | KILLED | 2 | 4.6 s | PosParity, PosCLI/refuse-rc |
| P13 | one space appended to `v65/neg-attested+14.json` | KILLED | 1 | 4.6 s | PosFixtures |
| P14 | `Shots` graph 2100 → 2101 | KILLED | 3 | 5.0 s | TestMarketingPNGContent/canonical-shots, TestMarketingShotGeometry/graph.png, PosGeometryMatchesOracle |
| P15 | `Shots` packet 3700 → 3660 | KILLED | 43 | 4.6 s | PosCorpus/stores/r1-seed1, PosCorpus/stores/r1-seed2, PosCorpus/stores/r1-seed1-blockconverted, PosCorpus/stores/r1-seed2-blockconverted, PosCorpus/stores/r1b-seed3, … (+38) |
| P16 | orchestration pin line deleted | KILLED | 2 | 1.8 s | PhParity/pin-registry, PosParity/pin-registry |
| P17 | runner deadline removed (`WithCancel`) | KILLED | 1 | 21.0 s | SubprocessDeadline |
| P17b | runner process group removed (`Setpgid` line) | KILLED | 1 | 10.8 s | SubprocessDeadline |
| O1 | one byte appended to `v64-p0check-main.go.txt` | KILLED | 1 | 0.4 s | PosParity |
| O2 | one byte appended to `v66-phonecheck-main.go.txt` | KILLED | 1 | 0.4 s | PhParity |
| F1 | (ii) removed | KILLED | 5 | 4.7 s | PhArms/v67/timeline-only, PhArms/v67/blank, PhArms/new/header-solid-ink, PhArms/new/header-15px, PhParity |
| F2 | dims removed | KILLED | 3 | 4.6 s | PhArms/v67/crop-370, PhArms/new/tall-845, PhParity |
| F3 | (ii) back to the std-dev stand-in | KILLED | 2 | 4.6 s | PhArms/new/header-15px, PhParity |
| F4 | `X0 <= 8` removed | KILLED | 6 | 4.6 s | PhArms/judge/B1-thumb-split-px150, PhArms/judge/B1c-thumb-split-px8, PhArms/judge/B5-two-pieces-rows-differ, PhArms/judge/short-thumb-x20, PhArms/judge/thumb-x9, PhParity |
| F5 | thumb-present removed | KILLED | 10 | 4.6 s | PhArms/v67/no-scrollbar, PhArms/v67/no-overflow, PhArms/v67/timeline-only, PhArms/v67/blank, PhArms/judge/B4-run-100px, … (+5) |
| F6 | `X1 <= 384` removed | KILLED | 4 | 4.6 s | PhArms/v67/full-thumb, PhArms/judge/B3b-grey-bar-full-width, PhArms/judge/thumb-end385, PhParity |
| F7 | clip removed | KILLED | 9 | 4.6 s | PhArms/v67/title-unclipped, PhArms/v67/no-overflow, PhArms/v67/crop-370, PhArms/v67/pad-370, PhArms/v67/timeline-only, … (+4) |
| F8 | (i) removed | KILLED | 7 | 4.6 s | PhArms/v67/shift-5, PhArms/v67/shift-10, PhArms/v67/timeline-only, PhArms/v67/blank, PhArms/judge/title-ink-spoof-pixel-388, PhArms/new/entry-solid-ink, PhParity |
| F9 | longest run → first-to-last span | KILLED | 4 | 5.0 s | PhArms/judge/B1-thumb-split-px150, PhArms/judge/B1c-thumb-split-px8, PhArms/judge/scattered-pixels-no-thumb, PhParity |
| F10 | renderer stripe period 2 → 3 | KILLED | 11 | 5.2 s | PhSkeleton/skeleton-eval-p1-pass-a-phone, PhSkeleton/skeleton-eval-p2-pass-a-phone, PhSkeleton/skeleton-eval-p3-pass-a-phone, PhSkeleton/skeleton-r1-seed1-pass-a-phone, PhSkeleton/skeleton-r1-seed2-pass-a-phone, … (+6) |
| F11 | `zero-bands` reason removed | KILLED | 2 | 4.7 s | PhArms/v67/blank, PhCLI/blank |
| D1 | §13 export reverted to the R1 file | KILLED | 2 | 3.4 s | Runbook/rects-value, Runbook/position-lines |
| D2 | §13 `release=` names the R1 path | KILLED | 2 | 1.8 s | Runbook/rects-value, Runbook/position-lines |
| D3 | §13 refuse path mistyped | KILLED | 1 | 1.8 s | Runbook/position-lines |
| G1 | `position.go` allow-list line removed | KILLED | 1 | 1.2 s | TestProductionGoSurface |
| G2 | `phone.go` allow-list line removed | KILLED | 1 | 1.6 s | TestProductionGoSurface |

The verdict-level properties the design claims are intact: P1–P4 flip no V64/V65 *verdict* (their extra
Corpus reds are tag changes, D-4); P14 is verdict-neutral (no Corpus or Parity red; only the
geometry pin and the two pre-existing `Shots` tests); P15 reds parity; F3 is killed only by
`header-15px` and the ledger. Full output: `plan-mutants-planner-run.txt`.

**V35** (AC-4, P17): unmutated, `TestMarketingSubprocessDeadline` passes in about 2.0 s with the
grandchild reaped; with the deadline removed it fails after about 20 s (the fake's own sleep); with
the process group removed (P17b) it fails because the grandchild survives. The drill output for P17
and P17b is the V35 record.

---

## 9. Planner validation (what was measured, and on what)

Trees: `.plan-w328/proto` (development), `.plan-w328/val` (fresh `git archive HEAD` of `08dc31be`,
banked files applied milestone by milestone exactly as §3–§7 say), `.plan-w328/mut` (the final
tree, for the drill), `.plan-w328/base` (the untouched base). All siblings of the repo; removed
after the run.

### 9.1 Per-milestone cumulative trees (`val`, banked files applied in order)

Every row: `gofmt -l` printed nothing, `go vet` rc 0, compile fence `go test -run '^$' ./...` rc 0.
Package times are wall seconds for `go test -count=1 [-race] -timeout 20m ./scripts/marketingcheck`;
`uptime` load averages are the 1-minute value at the milestone's start (other missions share the
rig). Base package at CI: 26.9 s plain, 157 s race (design V19).

| Milestone | Load | Plain | Race | Acceptance (measured) |
|---|---|---|---|---|
| M1a | 3.81 | rc 0, 38 s | rc 0, 89 s | `^TestMarketingSubprocessDeadline$`: 1 PASS (2.00 s, grandchild reaped) |
| M1b | 5.66 | rc 0, 41 s | rc 0, 105 s | generator rc 0, `files=34`, `SHA256SUMS` cmp-equal (`d8963a3c…300a`); `^TestMarketingPosition`: 2 PASS, 16 subtests PASS, 0 FAIL; `position --help` rc 0 |
| M1c | 4.33 | rc 0, 46 s | rc 0, 95 s | `^TestMarketingPosition`: 4 PASS, 60 subtests PASS (44 Corpus + 16 CLI), 0 FAIL |
| M1d | 5.42 | rc 0, 38 s | rc 0, 93 s | `^TestMarketingPosition`: **5 PASS**, 0 SKIP, 0 FAIL, 61 subtests; `PARITY position: inputs=44 verdict-diffs=0 reason-diffs=0` |
| M2a | 6.14 | rc 0, 37 s | rc 0, 99 s | skeleton write rc 0; 10 skeletons + `corpus-SHA256SUMS` sha-equal to the prototype's; `^TestMarketingPhone`: 2 PASS + 1 SKIP; with the corpus 3 PASS, 0 SKIP; `phone --help` rc 0 |
| M2b | 4.34 | rc 0, 37 s | rc 0, 117 s | `^TestMarketingPhone`: **4 PASS + 1 SKIP**; with the corpus **5 PASS, 0 SKIP**; `PARITY phone: inputs=41 diffs=0 ledger=1`, `inputs=59 diffs=0 ledger=1`, `CHECK images=10 byte-equal=10` |
| M3 | 5.22 | rc 0, 38 s | rc 0, 118 s | apply rc 0; `TestMarketingRunbookMatchesHelp\|TestMarketingDocsLinks`: 2 PASS, 22 subtests, the 3 new ones PASS; V34 `grep -c` = 1 and 1; capture-day step 1 (`^TestMarketing(Position\|Phone)` with the corpus): 10 PASS, 0 SKIP, `--- PASS: TestMarketingPhoneRealCorpus` 1 |

Two files changed after this run and were re-validated separately: the `host/store` allow-list
snapshots (D-15; `TestProductionGoSurface` PASS on the base tree with the M1b files, `reproducers=5`,
and with the M2a files, `reproducers=6`) and the D-11 sentence in the doc-edit script (re-run
inside the repo gate's `scripts/marketingcheck` package and by D1–D3 in the drill).

### 9.2 Repo gate on the final tree

- **Run 1** (`git archive` copy, before D-15 and D-11; started 01:38Z, load 3.98): `verify_ail.sh`
  rc 0; `go build ./...` rc 0; `go test -count=1 -timeout 20m ./...` rc 1 in 367 s with two reds:
  `host/store TestProductionGoSurface` (the real gap, D-15) and `host/pins
  TestReleaseLiteralsExplainedByManifest` (`git ls-files: exit status 128`: the copy is not a git
  checkout).
- **Run 2** (a `git clone` of the worktree with the final tree committed; 01:46Z, load 4.97 → 10.59):
  `verify_ail.sh` rc 0; `go build ./...` rc 0; `go test ./...` rc 1 in 410 s, 34 `ok`, one red:
  `host/verifygate TestMissionConfigDefaultFlowIntegration/invalid_ledger` (`default-flow fixture
  timed out … context deadline exceeded`, a 10 s fixture bound) at load ≈10. Re-run alone at load
  8.87: rc 0, `invalid_ledger` 4.15 s. It is a timing bound in a package this row does not touch;
  `host/pins`, `host/store` and `scripts/marketingcheck` (128.9 s under the parallel run) are `ok`.
  **Executor rule:** a red in a package this row does not touch is re-run alone once with the same
  `-run`; green alone = record it as a load flake in the report; red alone = STOP.

### 9.3 Mutation drill

`plan-mutants.py.txt` on a `git archive` of the final commit (01:54Z, load 8.40):
**36 of 36 KILLED**, 0 survived, 0 build failures, every file restored byte-equal (§8 table).

### 9.4 Fidelity checks

- Position fixtures: the 19 seed-5 and 4 eval-p1 arm hashes equal design V23's lists; both
  block-converted files `cmp`-equal the banked ones; the 9 stores equal `v64-inputs-sha256.txt`.
- The ten in-test position arms are byte-equal to the banked V13 generator's output (10/10).
- Phone arm ports: `plan-fidelity-test.go.txt` against the banked generators (phoneneg, v67b judge,
  newarms, h15) built from their `.txt`: `FIDELITY render: arms=39 pixel-equal=39`,
  `FIDELITY real: arms=39 pixel-equal=39` (seed-5 render and real seed-5 base).
- Skeletons: 10 files, 23649 bytes, seed 5 2366 B and seed 2 2365 B (= design V15/V31).
- `phoneWant` equals design V15 column 3 on all 39 arms plus the two renders (12 ACCEPT).

### 9.5 S7

`plan-m3-s7-runbook.sh.txt` on the final tree: rc lines `0, 1, 0, 1, 0`; the refuse control prints
`packet.png/roster-attested: (a3) gap after roster-tested 29 vs R1 15`; the phone refuse on the
real desktop shot prints `(i)`, `(iii′) dims 1440x1100`, `(iii′) thumb {-1 -1 -1 -1}`; the phone
release on the real seed-5 shot prints `ACCEPT 390x844 bands=21 … thumb={6 833 333 841}
titleRight=389`. Under `go run` the non-zero exits print `exit status 1` (D-11).

Logs: `plan-validation-run.txt`, `plan-mutants-planner-run.txt`.

---

## 10. Gates (state them exactly)

### 10.1 Per-milestone gate (every milestone, before its commit)

```sh
gofmt -l scripts/marketingcheck > $SCR/gofmt.log; echo "rc=$?"; cat $SCR/gofmt.log      # prints nothing
go vet ./scripts/marketingcheck && go test -count=1 -timeout 20m ./scripts/marketingcheck > $SCR/pkg-plain.log 2>&1; echo "rc=$?"    # rc=0, "ok  …/scripts/marketingcheck"
go test -count=1 -race -timeout 20m ./scripts/marketingcheck > $SCR/pkg-race.log 2>&1; echo "rc=$?"                                  # rc=0
go test -count=1 -timeout 20m -run '^$' ./... > $SCR/compile.log 2>&1; echo "rc=$?"     # rc=0: the compile fence for _test.go files (go build does not compile tests)
go test -count=1 -timeout 20m -run '^TestProductionGoSurface$' ./host/store > $SCR/surface.log 2>&1; echo "rc=$?"   # rc=0 (the census of non-test Go files, D-15)
```

### 10.2 Repo gate (after M3; CI runs the same)

```sh
./scripts/verify_ail.sh > $SCR/verify_ail.log 2>&1; echo "rc=$?"                       # rc=0
go build ./... > $SCR/build.log 2>&1; echo "rc=$?"                                      # rc=0
go test -count=1 -timeout 20m ./... > $SCR/test-all.log 2>&1; echo "rc=$?"              # rc=0, no "FAIL" line
grep -c '^FAIL' $SCR/test-all.log                                                       # 0
WORLD_PHONE_CORPUS=$CORPUS go test -count=1 -timeout 20m -v -run '^TestMarketing(Position|Phone)' ./scripts/marketingcheck > $SCR/capday.log 2>&1; echo "rc=$?"   # rc=0
grep -c '^--- SKIP' $SCR/capday.log                                   # 0
grep -c '^--- PASS: TestMarketingPhoneRealCorpus' $SCR/capday.log     # 1
cat $SCR/verify_ail.log $SCR/build.log $SCR/test-all.log > $BANK/executor-repo-gate.txt
cp $SCR/capday.log $BANK/executor-capture-day-step1.txt
```

`host/pins` `TestReleaseLiteralsExplainedByManifest` runs `git ls-files`, so the repo gate must
run in a git checkout (the worktree is one); in a `git archive` copy it fails with `git ls-files:
exit status 128` (measured; an instrument artefact, not a defect).

**`./scripts/verify_go.sh` is the controller's gate** (what CI runs: plain leg, race leg at
`-timeout 8m` per package, the binary gate and the race-detector control). It needs
`AILANG_BIN=$HOME/.pinned-ailang/ailang` and takes tens of minutes under the rig's load; the
executor does not run it. The controller runs `AILANG_BIN=$HOME/.pinned-ailang/ailang
./scripts/verify_go.sh` and reads CI on the PR. The new phone and position work stays inside the
race leg's per-package 8-minute bound (measured package `-race` totals in §9).

---

## 11. Deviations from the design (each inside the design's direction)

- **D-1 (milestone split).** The design's three milestones are seven commits here: M1 →
  M1a–M1d, M2 → M2a–M2b, M3 → M3, for the ~250-LOC-per-milestone guidance. M1b (≈383 lines) and
  M2a (≈568) exceed it: each carries a 229/224-line port of one banked program plus the smallest
  test set that exercises it (M2a's CLI test needs a skeleton render, and the skeletons are written
  only by `TestMarketingPhoneRealCorpus`). Every intermediate tree was validated green (§9.1).
- **D-2 (runner test name).** The runner's own control is a top-level
  `TestMarketingSubprocessDeadline` (in `position_test.go`), so the design's M1 acceptance count
  (`^TestMarketingPosition` → 5 PASS) holds. It is outside capture-day step 1's
  `^TestMarketing(Position|Phone)` pattern; CI runs it. A second runner mutant, **P17b** (process
  group removed: `Setpgid` line deleted), is added beside P17; it is killed by the
  surviving-grandchild check.
- **D-3 (in-process CLI).** `runPosition`/`runPhone` take `(args []string, stdout, stderr
  io.Writer) int`; `run` passes `os.Stdout`/`os.Stderr`. The CLI tests call them in process (fast,
  race-covered). Subprocesses are only the oracle builds and runs, `ps` in the runner test, and
  M3's documented-line runs of the built binary; all go through the deadline runner.
- **D-4 (tags make the corpus stricter).** `TestMarketingPositionCorpus` and
  `TestMarketingPhoneArms` assert the verdict **and the sorted set of clause tags** (design §3 M1:
  "the failing clause's tag"). So most §4 "Kills" cells are wider than the design wrote: P1 reds 5
  Corpus subtests (r1-seed1/2, neg-light-header-whole-view and neg-tested-64x16-glyphs also carry an
  (a1) reason), not only `iso-a1-tested-w-3`. The verdict-level claim (V12: P1–P4 flip no V65
  verdict) is untouched; §8 lists the measured red sets.
- **D-5 (P12/P6 also red the CLI).** `TestMarketingPositionCLI/refuse-rc` pins the documented
  `(a3) gap after roster-tested 29 vs R1 15` line verbatim (the line §2.5 step 3 reads), so P12 and
  P6 red it as well as parity. The design said only parity reds P12.
- **D-6 (V13 byte pin).** `TestMarketingPositionFixtures` adds a fifth check: the ten in-test
  arms, marshalled with the package's own types, are byte-equal to the banked V13 generator's
  output (its ten sha256, hardcoded). Measured 10/10 equal.
- **D-7 (59 real parity inputs).** The corpus holds pass A only. The real-corpus parity reaches
  the design's 59 V11 inputs by labelling each store's pass-A image as both `<store>/A` and
  `<store>/B`, because pass B is byte-equal to pass A for every store (V08).
- **D-8 (D1 via a path guard).** `position-lines` requires `$WORLD_MARKETING_RECTS`, expanded
  from §13's export, to lie under `$DEMO_DIR` and fails otherwise, instead of running the release
  line on the reverted path: writing the fixture to the reverted R1 path would overwrite a tracked
  repo file. Measured: D1 reds `rects-value` and `position-lines`, as the design says, by this
  route.
- **D-9 (phone CLI uses the blank page).** `TestMarketingPhoneCLI/refuse-rc` uses an inline
  v67/blank image (no arm generators needed in M2a); `tall-845` is exercised by
  `TestMarketingPhoneArms` and M3's `phone-lines`.
- **D-10 (oracle reference image in CI).** The phone oracle needs a reference PNG argument; it
  only feeds its r2 field and `CONTROL` line, which parity ignores. CI passes the seed-5 render;
  the real-corpus parity passes V50.
- **D-11 (`go run` exit codes).** Measured: `go run ./scripts/marketingcheck position` with no
  arguments exits **1** and prints `exit status 2`; the built binary exits 2. So under §13's and
  §4.4's `go run` lines a usage error reads as rc 1 (and prints no label line, so the "read rc and
  the label line together" rule still BLOCKs it). M3 appends one sentence to §13's exit-code
  paragraph saying so. M3's `position-lines`/`phone-lines` run the built binary, so they assert the
  binary's codes.
- **D-12 (mutant forms).** P13 appends a space (JSON still parses, so only the sha pin can see it);
  P6/P8 use `&& false` forms (a bare `if false` leaves a variable unused and does not compile);
  F3 restores the oracle's own std-dev (0 for an empty rectangle), not `rectStats` (NaN). Two
  oracle-pin mutants are added: **O1/O2** (one byte appended to each banked source) red the parity
  test's sha check.
- **D-13 (doc wording).** The QUICKSTART P3 and phone paragraphs around the design's commands are
  planner-authored (design §2.7 gave the commands and placement only); the orchestration
  replacements use the design's text, with step 1–5 wording taken from design §2.5.
- **D-14 (instrument: nohup).** Running the package under `nohup` reds the pre-existing
  `TestMarketingSupervisorFakeDaemonCleanup/shell-hup` (measured twice, `signal: killed` after
  35 s; 3/3 green without `nohup` at base and on the final tree). Rule 3 of §1 forbids it.
- **D-15 (a file outside design §6).** `host/store/context_roots_test.go`'s
  `TestProductionGoSurface` allow-lists every non-test Go file outside `host/` and `cmd/`
  (`scripts/marketingcheck/main.go` and `png.go` are listed for row 150). The planner's first
  repo-gate run on the prototype went red there: `non-test Go file outside context-root scan and
  explicit reproducer allow-list: scripts/marketingcheck/position.go` (and `phone.go`). M1b and
  M2a each add one allow-list line with a row-231 comment; G1/G2 (each line removed) are in the
  drill. AC-5's file list therefore gains this one test file.

**No STOP.** Nothing measured contradicts the design's direction.

---

## 12. Not measured by the planner

- **`./scripts/verify_go.sh` was not run by the planner** (the controller's gate; the CI
  workflow is unchanged). Its race leg's per-package bound is 8 min; the planner measured the
  package's local `-race` run at 89–118 s per milestone under load average 4–6, against CI's 157 s
  at base (design V19). The CI race time of the final package is unmeasured.
- **CI never sees a real Chrome phone image** (design §8, unchanged): the 10 skeletons stand in;
  their fidelity is measured on 10 images and 39 arms, pixel-equal port fidelity on 39 + 39.
- **Linux behaviour of the runner test** (`Setpgid`, `ps -o stat=`, zombie handling) was
  measured on darwin only. It uses only POSIX calls, and CI runs ubuntu; the first CI run is the
  measurement.
- **The design's capture-day step 4/5** (phone accepts on V50 and seed 2, phone refuse on the
  capture day's own desktop shot) is documented, not executed: it is capture-day work. The S7
  script runs the analogous lines on the real seed-5 phone and desktop shots (rc 1 / rc 0).
- **`go run` vs binary exit codes** were measured for the usage case only (D-11).
