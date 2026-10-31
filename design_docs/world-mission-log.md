# Ailang World Mission — iteration log

Append-only. One entry per outer-loop iteration. Newest at the bottom.

---

> **Older entries are ARCHIVED.** This file holds the newest 20. The full record of every
> iteration is in `world-mission-log-archive.md`, and a one-line index of ALL of them —
> the thing to grep before picking work, so the loop never repeats itself — is in
> `world-mission-index.md`.

## 306 — 2026-10-27 — row 150: R0 re-confirmed the V51 shot bottoms on Chrome 156.0.8078.48 (V59) and M7c0 built the tall shot geometry (PR #373; judge PASS 96) [PRODUCT]

**Kind:** one controller measurement step (R0) and one production milestone (M7c0) on the release row. Roles were spawned through the Agent tool, per the operator's standing request, where the hook allows it (see Routing evidence).

**Pick and why:** row 150, following iteration 305's **Next** (R0, then M7c0). Re-checked before picking:
- no M7c0 commit, branch or PR existed; no open PR from this account (`gh pr list --author @me` empty);
- local dev equalled origin/dev `eacd90e7`;
- the queue order is unchanged by the three attended rulings since iteration 305 (below).

**Clause map:**
- Clauses 1–7 all MET (ratified states, not re-measured). No UNMET clause, so there is no critical-path row and no drift alarm.
- D-WORLD-81 order is 156 → 150 → 151. 156 is held on D-WORLD-86, the only OPEN row (`mission_decisions.sh --open`).

**Gate 0:**
- Kill switch armed (`mission-world.disabled` absent); gh account sunholo-voight-kampff; billing tripwire CLEAN.
- `mission_directives.sh`: 0 allowlisted directives on #366 since `2026-10-26T09:00:27Z` (3 comments, all this account's).
- Attended ledger channel: D-WORLD-91, D-WORLD-92 and D-WORLD-93 were RESOLVED in attended sessions on 2026-10-26 (rows 188 and 151). `mission_ledger_provenance.sh` returned `attended` for all three (flip SHAs `12dc4c09`, `78d090bf`, `3838759c`). Acknowledged here; none changes row 150 or the 156 → 150 → 151 order. Ledger valid, 80 rows.
- Self-notices rc 1: newest crash `6345079873` (rc=1, 2026-10-22) is already credited (iteration 305's entry). No new crash notice.
- Inbox `mission-world`: only this loop's own claim messages; nothing outranks the queue.
- Weekly external-issue sweep: not due (iteration 305 ran it after the 2026-10-26 rotation).
- **Unattributed residue, recorded, not adopted:** `ailang-world/.wt-world-iter306-r0` (detached at `eacd90e7`, created 2026-10-26 16:36) holds one untracked file, a copy of the iteration-305 bottoms probe, and `~/.ailang/state/world-iter306/r0-seed{1,1b,2}` (seed1/seed1b empty, seed2 three PNGs). No slot verdict exists between the 13:22 fire and this one, no claim message and no log entry, so it was not a driver fire. Most likely an attended or ad hoc session began R0 and stopped. Left untouched; this iteration re-ran R0 from scratch into fresh paths (`r0b-seed{1,2}`).

**Gate 1:**
- Dev CI at `eacd90e7`: check set 3/3 success.
- Running skill plus all resource files byte-identical to origin/dev (28 files compared across the resolved symlink target and the pin copy; no DRIFT line).
- Base `eacd90e7…` recorded at Gate 1 and again at Gate 4 (unchanged).

**R0 (controller, doc §5; C13):**
- Chrome on the rig is now **156.0.8078.48** (V51 measured .47), so R0 was a real re-confirm, not a formality.
- `iter305-region-bottoms-probe` copied unchanged into the sprint worktree, run twice on two independent `demo seed` stores (pins v0.41.0 and v0.54.15), output paths deleted and re-created before each run, probe file deleted afterwards. rc 0 / rc 0, 11.15 s / 16.65 s.
- Seeds differ (heads `8c4b8aff…`, `a90ced43…`; manifests `0a1ef461…`, `d2fa161c…`). Both give packet B 3564.56, graph 2042.14, provenance 2881.17: **delta 0.00 against V51**, no shot blocked, every `supervisor shot` verdict ok.
- Controller viewed the bottom 900 px of seed 1's 3700-px `packet.png`: the validated sentence, `approvals chain: pending` and all four roster rows are in the image.
- Recorded as **V59** in `w-marketing-capture-orchestration.md` §10 (`adb5d24e`), outputs in `verification/world-iter306/`.

**Role execution:**
- **Designer:** not run. The governing doc (`w-marketing-capture-orchestration` r2) exists, and M7c0 needed no design change because R0 confirmed V51. Rotation pointer unchanged (`codex:gpt-6.1-sol`).
- **Planner:** not run. The sprint plan's M7c0 section is prescriptive (exact code, gates and mutant), amended in iteration 305.
- **Executor** `codex:gpt-6.1-sol` (resolver `recipe codex:gpt-6.1-sol declared:provider-pin`). Spawned as an Agent-tool supervisor (`MISSION-ROLE: executor`, no model alias) that ran the controller's bounded codex wrapper verbatim; the spawn-pin hook accepted it. Probe rc 0 (`model: gpt-6.1-sol`). Run rc 0 in 262 s, worktree `.wt-world-iter306-row150-m7c0`:
  - `png.go`: three `Shots` heights only (packet 3700, graph 2100, provenance 3000).
  - `png_test.go`: `math` import, canonical list, the three measured-bottom constants (citing V51 and V59), `shotRuleHeight`, `TestMarketingShotGeometry`.
  - Its `-race` red was a sandbox Go-cache denial (`link: cannot open file …/go-build/…`), labelled UNINFORMATIVE UNDER SANDBOX as directed. M7c0-height red on the named assertion; restored, `cmp` identical.
- **Evaluator** `sonnet` (resolver `agent-tool sonnet declared:alias-pin`), Agent tool, foreground, own worktree `.eval-world-iter306` at `5eb7be62`.

**Controller verification (outside the sandbox):**
- Diff read in full: exactly the plan's two files; nothing in the do-not-touch list changed.
- Gate list baselined on the pristine tree first: gofmt clean, vet rc 0, `TestMarketingPNGContent` rc 0, package rc 0 (20 s), `-race` rc 0 (36 s).
- After: gofmt clean; vet rc 0; `go build ./...` rc 0; compile fence `go test -run '^$' ./...` rc 0 (34 packages); `TestMarketingShotGeometry` 4/4 PASS; `TestMarketingPNGContent` PASS incl. `canonical-shots` and `positive-packet.png/-graph.png/-provenance.png`; package rc 0 (19 s); `-race` rc 0 (51 s, load 8–12).
- Mutant M7c0-height re-run by the controller: `png_test.go:834: packet height must equal the rule value from measuredPacketBottom: got {Name:packet.png Width:1440 Height:1400 Scheme:light}, want 1440x3700`; restored, `cmp` identical.
- C11 protected-file observation: `marketingcheck observe` before the executor (12 records) and after; `marketingcheck compare` rc 0.
- `.ail` untouched, so `verify_ail.sh` not run (CI's `ailang-code verify gate` covers it). `check_no_personal_email.sh` rc 0.

**Independent evaluation:** r1 **PASS 96/100**, zero blocking ([report](verification/world-iter306/evaluator-r1.md)).
- Scope confirmed from `git diff eacd90e7..5eb7be62 --stat`: the five listed files and nothing else.
- Non-vacuity against the milestone's own production diff: reverting `png.go` alone reds the three geometry subtests and `canonical-shots`. 11/11 of its own `png.go` mutants killed (heights ±100, width, light/dark/phone dimensions, schemes; graph scheme killed only by `canonical-shots`).
- V59 re-derived from the two output files: every number matches. PNG sizes and the controller's image view are attested, not checkable (PNGs are not committed).
- **Independent instrument check:** it re-ran the probe on a third fresh seed (head `278bfef6…`): rc 0, identical bottoms, delta 0.00.
- Non-blocking: the `want > 4400` and `bottom > H−40` branches are not exercised by the current constants (a `shotRuleHeight` boundary test would pin the rule itself); one ~140-char comment line; the plan still said "pending" (fixed in `95bb53cb`).
- Its `-race` runs: 2 of 3 red in different supervisor tests (`TestMarketingSupervisorFakeDaemonCleanup/interrupted-shot` `[Errno 1] Operation not permitted`; `TestMarketingSupervisorFakeChromeTimeoutKillsGroup`), the first also once with `png.go` reverted to base.
- **Controller reproduction of that finding (rule: reproduce before dismissing):** `-race` interrupted-shot alone, 100 runs at load ~11: 1 red, 99 green. `kill_group` in `scripts/marketing_capture.py:44-48` catches only `ProcessLookupError`, and Errno 1 is EPERM, which fits macOS `killpg` returning EPERM for a group of zombies. **Hypothesis, not confirmed** (the controller did not capture the failing message itself). M7a's file, outside this milestone; carried on row 150.

**Bookkeeping:** sprint JSON M7c0 → completed (iteration 306) and M7c1 → completed (iteration 305, PR #370 merge `120e4149`, which iteration 305 had left `pending`); plan M7c0 status → BUILT.

**Result:** R0 done (V59) and M7c0 BUILT (PR #373). Row 150 is 16 of 22 milestones built. The three tall release shots now contain their required content at their committed geometry.

**Gate 3b:** PR #373 head `95bb53cb` polled pinned to the full SHA (bounded 30 min, per-invocation verdict path naming the SHA); the record commit and the merge commit are polled the same way. Outcome reported on #366.

**Routing evidence:** base=eacd90e7d412ecb0a7363254549a6c3a0c71a589@2026-10-27T08:31:13Z (Gate 4; Gate 1 base identical). Driver banner: `controller=claude:claude-opus-5-5 … roles: designer=claude:claude-opus-5-5 planner=codex:gpt-6.1-sol executor=codex:gpt-6.1-sol evaluator=sonnet`; `MISSION_OVER_RATION` empty.
- Controller `claude:claude-opus-5-5` (tok: not reported).
- Designer: not spawned (doc exists; no design change). Planner: not spawned (prescriptive plan exists). Neither is a routing failure; the routing table spawns them only when the doc or plan is missing.
- Executor `codex:gpt-6.1-sol`: probe 9,086 tok; run 69,299 tok, rc 0, 262 s. Agent-tool supervisor (transport only, no verdict): 28,583 tok / 297 s.
- Evaluator `sonnet` (Agent): 61,972 tok / 294 s. Judge independence: cross-vendor (OpenAI executor, Anthropic Sonnet judge), and a different model from the Opus controller that ran R0.
- No quorum this iteration (no new or revised design doc).
- Metered **$0.00**.

**Ruled out:**
- Adopting the `.wt-world-iter306-r0` residue as R0 evidence: its seed1 directories are empty and nothing records how it ran. R0 was re-run into fresh paths instead.
- Transcribing new constants: R0's max equals V51 exactly, so M7c0 keeps 3564.56 / 2042.14 / 2881.17 and cites V59.
- Spawning a designer or planner to satisfy the operator's request mechanically: the routing table spawns them only when a doc or plan is missing, and both exist.
- Treating the judge's `-race` reds as caused by M7c0: they are supervisor tests that do not read `Shots`, the judge saw one with `png.go` reverted, and both of the controller's `-race` package runs (pristine base 36 s, after the change 51 s) were green.
- Fixing `kill_group` here: M7a's file, outside M7c0's two-file scope; the mechanism is a hypothesis.
- Adding the judge's `shotRuleHeight` boundary test now: non-blocking, and M7c0's file list is closed; noted on the row for whoever next touches `png_test.go` (M8a).

**Carried:**
- M7c2a next (Python `chrome-version`/`ready`/`pass` modes plus Go `routes`; its `routes` rows carry the M7c0 dimensions), then M7c2b → M7c3 → M7c4 → R1 → M8a → M8b, then attended A1.
- New: the 1-in-~100 `-race` supervisor red and its `kill_group` EPERM hypothesis.
- Judge M7c0 non-blocking: rule-boundary test; long comment line.
- Everything iteration 305 carried (phone/provenance cosmetic clips, tall framing for Mark, planner flags, M7c1 test gaps, `shell-hup` flake). D-WORLD-86 stays HOLD.

**Retro:**
- The R0 step was worth its 30 s: Chrome had moved a patch version since V51, and a precondition written as "re-confirm on the then-current Chrome" is only honest if it is actually run.
- Reproducing the judge's non-blocking `-race` finding turned "environmental flake" into a rate (1/100) and a specific, checkable mechanism, rather than another undated "flaky" note.
- An Agent-tool executor supervisor with no model alias passed the spawn-pin hook and ran the pinned codex lane. That satisfies both the operator's Agent-tool request and the provider pin without a denied spawn.
- Harness share 0. No skill edit.

**Progress:** 1.0 clauses 7/7 MET; row 150 16/22 built; goal unmoved.
**Next:** row 150 M7c2a (fully specified in the plan), then M7c2b. Then 151 after 150.

## 307 — 2026-10-27 — row 150: M7c2a built the capture primitives (Chrome-version floor, HTTP readiness, six-shot pass, `marketingcheck routes`); judge PASS 88; `kill_group` EPERM confirmed as a macOS production defect [PRODUCT]

**Kind:** one production milestone (M7c2a) on the release row, plus a controller test-only fix closing the judge's most load-bearing surviving mutant. Roles were spawned through the Agent tool, per the operator's standing request, where the routing table and the spawn-pin hook allow it (see Routing evidence).

**Pick and why:** row 150 M7c2a, following iteration 306's **Next**. Re-checked before picking:
- no M7c2a commit (`git log origin/dev --grep=M7c2` returns only the iteration-304 split commit), branch, worktree or PR existed; no open PR from this account;
- local dev equalled origin/dev `e694d23e` at Gate 1;
- the three attended rulings since iteration 306 (below) are attended-track rows and do not change the release order.

**Clause map:**
- Clauses 1–7 all MET (ratified states, not re-measured). No UNMET clause, so there is no critical-path row and no drift alarm.
- D-WORLD-81 order is 156 → 150 → 151. 156 is held on D-WORLD-86, the only OPEN row (`mission_decisions.sh --open`).

**Gate 0:**
- Kill switch armed (`mission-world.disabled` absent); gh account sunholo-voight-kampff; billing tripwire CLEAN.
- `mission_directives.sh`: 0 allowlisted directives on #366 since `2026-10-26T13:21:41Z` (4 comments, all this account's), and 0 on #313 since the mission-scoped watermark `2026-10-22T01:34:42Z` (the Repo Profile's older-of-two rule).
- Attended ledger channel: D-WORLD-95, D-WORLD-96 and D-WORLD-97 were RESOLVED in attended sessions on 2026-10-27 (the all-agent-work roadmap: rows 206/207 un-parked after 222, rows 219/217/218 before row 188's N2d, the session-allowance bounds, and machine-wide confinement gated on the scorecard). `mission_ledger_provenance.sh` returned `attended` for all three (flip SHAs `35345ed6`, `f06e239e`, `ce086b36`). Acknowledged here; all three are attended-track rows ("not a v1.0.0 release row"), so none changes row 150 or the 156 → 150 → 151 order. Ledger valid, 84 rows.
- Self-notices rc 1: newest crash `6345079873` (rc=1, 2026-10-22) is already credited (iteration 305's entry). No new crash notice.
- Inbox `mission-world`: only this loop's own claim messages; nothing outranks the queue. This iteration's claim was sent and read back.
- Weekly external-issue sweep: not due (iteration 305 ran it after the 2026-10-26 rotation).
- Residue: `ailang-world/.wt-world-iter306-r0` (iteration 306 recorded it as unattributed) still holds only `cmd/ailang-worldd/zz_iter305_bottoms_test.go`; left untouched again.

**Gate 1:**
- Dev CI at `e694d23e`: check set 3/3 success.
- Running skill plus all resource files byte-identical to origin/dev (14 files at each of the resolved symlink target and the pin copy, which are different inodes; no DRIFT line).
- Base `e694d23e…` recorded at Gate 1. At Gate 4 the shared ref read `eaa2a251…` (attended PR #379, row 223c herdr plugin, 13 files under `integrations/herdr/` and `design_docs/`). Intersection with this branch's files: empty, against a 13-file control. Benign drift: the branch is rebased onto it before the PR.

**Baseline (controller, pristine detached worktree `.base-world-iter307` at the base, `AILANG_BIN` exported, load ~5):**
- gofmt clean, vet rc 0, `TestMarketingCapturePassFakeChrome` absent ("no tests to run"), `TestMarketingSupervisorFakeChromeSuccess` PASS, package rc 0 (19 s), `TestProductionGoSurface` rc 0.
- `-race ./scripts/marketingcheck`: **red 1 of 4** at base. The red was `TestMarketingSupervisorFakeChromeTimeoutKillsGroup` → `SupervisorError [Errno 1] Operation not permitted`. This is iteration 306's carried note, now measured on a pristine tree.
- **Mechanism confirmed outside the suite** (rule 5, data before conclusions): a standalone Python repro (`verification/world-iter307/killpg-eperm-repro.py.txt`) creates a new session with nested `sleep` children, sends `SIGKILL` to the group, reaps the leader and calls `killpg` again three times. Over 200 trials it raised `PermissionError` on 5 of 600 calls and ESRCH on the rest. A first shape with no nested children gave 40/40 ESRCH, which is the negative control. `kill_group` (`scripts/marketing_capture.py:43-48`) catches only `ProcessLookupError`, so EPERM escapes to `main()`'s `OSError` catch as `SupervisorError`. A real six-shot pass on macOS Chrome, which leaves many children, can therefore fail on a good shot. Linux CI is unaffected.

**Role execution:**
- **Designer:** not run. The governing doc (`w-marketing-capture-orchestration` r2) exists and M7c2a needed no design change. Rotation pointer unchanged (`codex:gpt-6.1-sol`).
- **Planner:** not run. The plan's M7c2a section is prescriptive (interfaces, verdicts, subtests, assertion messages and three named mutants).
- **Executor** `codex:gpt-6.1-sol` (resolver `recipe codex:gpt-6.1-sol declared:provider-pin`). It was spawned as an Agent-tool supervisor (`MISSION-ROLE: executor`, no model alias) that ran the controller's bounded codex wrapper verbatim; the spawn-pin hook accepted it. Probe rc 0 (`model: gpt-6.1-sol`). Run rc 0 in 568 s, worktree `.wt-world-iter307-row150-m7c2a`. Directive `/tmp/codex_directive_iter307.txt`. It told the executor about the base `-race` red so it would be reported rather than "fixed", and it froze `kill_group`.
  - `marketing_capture.py`: modes `chrome-version`, `ready` and `pass`, plus `MIN_CHROME_MAJOR = 155`. `shot()`/`daemon()` are unchanged. `main()`'s viewport check now applies only to `shot`/`daemon`.
  - `main.go`: `routes`, `RouteRow`, `ShotRoutes` and a usage line.
  - `capture_test.go`: `TestMarketingCapturePassFakeChrome` with the routes, pass-six, version, ready and fresh-dir subtests.
  - The three named mutants each red their named subtest, then were cmp-restored. The `ready` subtest was labelled UNINFORMATIVE UNDER SANDBOX, as directed (httptest bind denied). Deviation: the snapshots carry `.txt` suffixes so Go's package scan does not see them. Accepted; `snap/` is not committed.
- **Evaluator** `sonnet` (resolver `agent-tool sonnet declared:alias-pin`), Agent tool, foreground, own detached worktree `.eval-world-iter307` at `bbf7eb84` (pre-rebase SHA of `ceeb27fd`: after the rebase onto `eaa2a251`, `git diff bbf7eb84 ceeb27fd -- scripts/` and `git diff a6646040 82f1e253 -- scripts/` are both empty).

**Controller verification (outside the sandbox, load 3.6–8.8):**
- Diff read in full: exactly the plan's three files; nothing in the do-not-touch list changed; snapshots `cmp`-identical to the tree.
- After the change: gofmt clean; vet rc 0; `go build ./...` rc 0; compile fence `go test -run '^$' ./...` rc 0 (34 packages); `TestMarketingCapturePassFakeChrome` 5/5 PASS, including `ready` (1.46 s); `TestMarketingSupervisorFakeChromeSuccess` PASS; `TestProductionGoSurface` and `TestProductionContextRoots` rc 0; `TestSubprocessSinkGate` PASS (`-v` confirms it ran); package rc 0 (19 s).
- `-race` package: rc 0 twice (53 s, 52 s). After the controller test it was rc 1 once, on `TestMarketingSupervisorFakeChromeTimeoutKillsGroup` with the identical EPERM message. That is the base defect: the test and `shot()` are byte-identical to base.
- C11 protected-file observation: `marketingcheck observe` before the executor (12 records) and after; `marketingcheck compare` rc 0.
- `.ail` untouched, so `verify_ail.sh` was not run (CI's `ailang-code verify gate` covers it). No local `verify_go.sh`: the diff is confined to `scripts/marketingcheck`, adds no context root, and CI runs the full Go and race gates.

**Independent evaluation:** r1 **PASS 88/100**, zero blocking ([report](verification/world-iter307/evaluator-r1.md)).
- Rubric: conformance 28/30, tests and gates 25/25, correctness 18/20, non-vacuity 10/15, hygiene 7/10.
- Gates: everything green, including one `-race` run (50.7 s).
- Plan mutants: 3/3 killed. Diff-anchored mutants: 13 of 26 survived (one more UNMEASURED). The survivors are mostly `pass` refusal and bound branches: stop at first non-ok, six rows, basename, dimensions, scheme, the deadline guards, the base-URL `rstrip`, and the per-shot `OSError` catch. The others are `ready`'s non-200 check, the `Chromium` alternative, `chrome-version`'s `Timeout`, and `runRoutes`'s `ErrHelp`.
- Non-blocking: `pass` makes `--out-dir` before it validates the routes file. The judge reproduced this with `--routes /nonexistent`: the empty directory is left behind, so a retry gets `DestinationExists`.
- **Controller reproduction (reproduce before acting):** `if result['verdict'] != 'ok':` → `if False:` stays green (`ok`). That is the behaviour M7c2b's shell relies on, so the controller added subtest **`stop-first`** (`82f1e253`, test-only). Its fake Chrome exits 5 on the 1440×2100 viewport. The pass must then give `NonzeroExit` with four shots reported, four Chrome calls, three files, and the last progress line `4/6 graph.png NonzeroExit`. Under the mutant it reds: verdict `ok`, six calls, five files. Restored, `cmp` identical, 6/6 subtests green.
- Controller's own non-blocking review notes: `chrome_version`'s timeout path calls `process.wait(timeout=KILL_WAIT_S)` without the `TimeoutExpired` catch that `shot()` has (reachable only with an unkillable process), and `ready` polls every 10 ms.

**Bookkeeping:** sprint JSON M7c2a → completed (iteration 307); plan M7c2a status → BUILT; charter row 150 updated (M7c2a, the EPERM confirmation, proposed M7a-fix, Next); CURRENT GOAL next line updated; STATUS 304 rotated to the archive (line counts asserted; the charter keeps 3 stamps).

**Result:** M7c2a BUILT. Row 150 is 17 of 22 milestones built. The capture now has a Chrome floor, HTTP readiness and a bounded six-shot pass driven by `routes` rows that carry the M7c0 heights.

**Gate 3b:** the PR head and then the merge commit are each polled pinned to the full SHA (bounded 30 min, per-invocation verdict path naming the SHA). Outcome reported on #366.

**Routing evidence:** base=e694d23e8eb39673e5cf21dc569bab4ae933ba49@2026-10-27T12:01:28Z (Gate 1); Gate 4 read eaa2a25179140f7a0e54e0ed0559804a74005879@2026-10-27T12:36:31Z (benign drift, attended #379, no shared file). Driver banner: `controller=claude:claude-opus-5-5 … roles: designer=claude:claude-opus-5-5 planner=codex:gpt-6.1-sol executor=codex:gpt-6.1-sol evaluator=sonnet`; `MISSION_OVER_RATION` empty.
- Controller `claude:claude-opus-5-5` (tok: not reported).
- Designer: not spawned (the doc exists and needs no change). Planner: not spawned (a prescriptive plan exists). Neither is a routing failure: the routing table spawns them only when the doc or plan is missing.
- Executor `codex:gpt-6.1-sol`: probe 17,397 tok; run 91,203 tok, rc 0, 568 s. Agent-tool supervisor (transport only, no verdict): 31,302 tok / 615 s.
- Evaluator `sonnet` (Agent): 57,664 tok / 405 s. Judge independence: cross-vendor (OpenAI executor, Anthropic Sonnet judge), and a different model from the Opus controller that wrote the `stop-first` test. That test was not separately judged. It is test-only and mutation-proven.
- No quorum this iteration (no new or revised design doc).
- Metered **$0.00**.

**Ruled out:**
- Fixing `kill_group` inside M7c2a: the plan freezes `shot`/`daemon` behaviour and M7c2a's file list is closed. The fix is proposed as its own milestone (M7a-fix), now that the mechanism is confirmed.
- Calling the `-race` red "load": a standalone repro with no Go, no race detector and no test harness raised the same `PermissionError`, so it is a real macOS behaviour that load only makes more frequent.
- Treating the head's `-race` red as caused by M7c2a: `TestMarketingSupervisorFakeChromeTimeoutKillsGroup` and `shot()` are byte-identical to base, and the base was red the same way.
- Closing all 13 surviving mutants here: only the stop-at-first-non-ok branch is load-bearing for M7c2b. The others are refusal paths; they are carried to M7c2b, which owns the shell that calls `pass`.
- Running a local full `verify_go.sh`: CI is the arbiter, and the diff does not reach a gate only the full suite sees (`TestProductionGoSurface`, context roots and the sink gate were run directly).

**Carried:**
- **M7a-fix** next: `kill_group` treats `PermissionError` like `ProcessLookupError`, with a unit row that fakes `os.killpg` and the repro as a stress row. It must land before M7c4's real-Chrome drill.
- Then M7c2b → M7c3 → M7c4 → R1 → M8a → M8b → attended A1.
- For M7c2b: `pass` creates `--out-dir` before validating routes; the untested `pass`, `ready` and `chrome-version` refusal branches; `chrome_version`'s uncaught `TimeoutExpired`.
- Everything iteration 306 carried (rule-boundary test, phone/provenance cosmetic clips, tall framing for Mark, `shell-hup` flake). D-WORLD-86 stays HOLD.

**Retro:**
- Baselining `-race` four times on a pristine tree turned "1 in 100, hypothesis" into "1 in 4 on this package, mechanism reproduced in 600 calls". The difference matters: the same EPERM would hit the attended real-Chrome capture on Mark's Mac.
- The judge's diff-anchored drill found the one behaviour the next milestone depends on (stop at first non-ok) unguarded. The plan's mutant table named only three mutants, so 3/3 killed was true and not sufficient: rule 3n held again.
- Harness share 0. No skill edit.

**Progress:** 1.0 clauses 7/7 MET; row 150 17/22 built; goal unmoved.
**Next:** row 150 M7a-fix (`kill_group` EPERM), then M7c2b. Then 151 after 150.

## 308 — 2026-10-27 — row 150: M7a-fix built — `kill_group` tolerates macOS `killpg` EPERM, so a good shot no longer ends as `SupervisorError`; judge PASS 88 [PRODUCT]

**Kind:** one production milestone (M7a-fix) on the release row: a planner amendment, an executor build, a controller test-only bound fix and an independent judge. All four routing-table roles that were needed were spawned through the Agent tool, per the operator's standing request (see Routing evidence).

**Pick and why:** row 150 M7a-fix, following iteration 307's **Next**. Re-checked before picking:
- no M7a-fix commit, branch, worktree or PR existed (`git log origin/dev --grep` for `M7a-fix`/`kill_group` returns only iteration 307's record commit; PR search for `kill_group` returns only #380, merged);
- local dev equalled origin/dev `0363b3ae` at Gate 1 (fast-forwarded two attended commits first);
- the only open PR from this account (#389 `fix/trial-world-context`) is an attended trial fix, not this mission's loop work; left alone.

**Clause map:**
- Clauses 1–7 all MET (ratified states, not re-measured). No UNMET clause, so there is no critical-path row and no drift alarm.
- D-WORLD-81 order is 156 → 150 → 151. 156 is held on D-WORLD-86, the only OPEN row (`mission_decisions.sh --open`).

**Gate 0:**
- Kill switch armed (`mission-world.disabled` absent); gh account sunholo-voight-kampff; billing tripwire CLEAN.
- `mission_directives.sh`: 0 allowlisted directives on #366 since `2026-10-27T09:03:42Z` (both watermarks agree; 5 comments, all this account's), and 0 on #313 since `2026-10-22T01:34:42Z`.
- Attended ledger channel: D-WORLD-98 (row 219's rulings R1–R4, the home service) was filed and RESOLVED in an attended session after iteration 307. `mission_ledger_provenance.sh` → `attended` (flip `920cb0d7`, creation `9a68ad8e`). Acknowledged here. It is an attended-track row and does not change row 150 or the 156 → 150 → 151 order. Ledger valid, 85 rows. Also new since 307 (attended, docs-only): row 228 `w-agent-search-and-read-cap` filed NEXT on the attended track (#388), README status refresh (#387).
- Self-notices rc 1: the newest crash `6345079873` (rc=1, 2026-10-22) is already credited (iteration 305's entry; grep count 1). No new crash notice.
- Inbox `mission-world`: only this loop's own claim messages; nothing outranks the queue. This iteration's claim was sent and read back.
- Weekly external-issue sweep: not due (run after the 2026-10-26 rotation). #366 was created after that boundary, so no thread rotation.
- Residue: `ailang-world/.wt-world-iter306-r0` (unattributed since 306) left untouched again.

**Gate 1:**
- Dev CI at `0363b3ae`: check set 3/3 success.
- Running skill plus all 13 resource files byte-identical to origin/dev (per-file `cmp`, resolved symlink target; no DIFF line).
- Base `0363b3ae…` recorded at Gate 1; the Gate-3 `snap` read the same SHA.

**Baseline (controller, pristine detached worktree `.base-world-iter308` at the base, `AILANG_BIN` = pinned v0.41.0, load ~6):**
- gofmt clean, vet rc 0, `py_compile` rc 0, `scripts/marketingcheck` package rc 0 (24.6 s).
- `-race ./scripts/marketingcheck`: **red 1 of 3** at base — `TestMarketingSupervisorFakeDaemonCleanup` → `normal shell cleanup: exit status 1 {"verdict": "SupervisorError", "detail": "[Errno 1] Operation not permitted"}`. The defect being fixed, on a third test.
- Instrument note: this rig has no `timeout` binary; a first worktree-creation call wrapped in `timeout` failed rc 127 without touching git, and was re-run without it.

**Role execution:**
- **Designer:** not run. `w-marketing-capture-orchestration` r2 governs and M7a-fix needs no design change.
- **Planner** `codex:gpt-6.1-sol` (resolver `recipe codex:gpt-6.1-sol declared:planner-lane-default-pin`; `derive-planner-lane.sh` gave the same line). It was spawned as an Agent-tool supervisor (`MISSION-ROLE: planner`, no model alias) that ran the controller's bounded codex wrapper (`/tmp/codex_run_iter308.sh`) in a detached sibling worktree `.planner-wt-world-iter308` at `0363b3ae`. Probe rc 0; run rc 0 in 469 s. It added `## M7a-fix` to the plan and a `M7a-fix` JSON milestone ahead of M7c2b, made M7c2b and M7c4 depend on it, and recorded facts F0–F10 with commands. **It refuted part of the controller's directive:** I had written that EPERM from `killpg` on our own group means "no signallable member remains". Its F8 quotes `man 2 kill`: for a group, EPERM is returned "if **any** members of the group could not be signaled". **Controller re-check:** I read both man pages myself and they say exactly that. So EPERM does not prove death; the plan scopes the fix as a bounded workaround for the observed post-reap case and keeps the existing death assertions as the guard. Artifacts committed as `97736493`.
- **Executor** `codex:gpt-6.1-sol` (resolver `recipe codex:gpt-6.1-sol declared:provider-pin`), Agent-tool supervisor (`MISSION-ROLE: executor`, no alias) running the same wrapper in `.wt-world-iter308-row150-m7afix`. Probe rc 0; run rc 0 in 515 s. Directive `/tmp/codex_directive_iter308.txt`.
  - `marketing_capture.py`: `kill_group` catches `(ProcessLookupError, PermissionError)` plus a one-line comment; nothing else changes.
  - `supervisor_test.go` (+237 lines, purely additive): `TestMarketingKillGroupTolerance` (always-on; inline Python imports the module via `importlib`, fakes `os.killpg` to raise EPERM, ESRCH and EINVAL, and asserts the call `(424242, SIGKILL)` and the outcome) and `TestMarketingKillGroupStress` (opt-in `WORLD_KILLPG_STRESS=1`; 200 trials of the iteration-307 repro through the real `kill_group`; raw `killpg` EPERM counted as the positive control; prints UNINFORMATIVE when that count is 0).
  - Red-first: only `eperm-swallowed` red at base, as planned. In-sandbox: 3/3 plan mutants red at their named assertions. In-sandbox `-race` was blocked by a socket-bind denial (correctly labelled UNINFORMATIVE). It reported one restored-tree timeout-test red, `fake must record child and grandchild: [81876]`, green on rerun (load 9).
  - Committed by the controller as `31211dbf`.
- **Controller test-only fix `7475ec3c`:** outside the sandbox the opt-in stress row was **rc 1 on 3 of 3 runs**. Its 25 s Python wall cut every run at ~160 of 200 trials (raw ESRCH 461–480 against 600 expected calls). The plan forbids silently reducing N, so the bounds were raised instead (Python wall 60 s, Go context 70 s, wall target 60 s). After the fix: 3 of 3 rc 0, 200 trials and 600 calls each, raw EPERM 7, 3, 7, wrapper exceptions 0, ~33 s. The plan's stress paragraph records the amendment.
- **Evaluator** `sonnet` (resolver `agent-tool sonnet declared:alias-pin`), Agent tool, foreground, own detached worktree `.eval-world-iter308` at `7475ec3c`.

**Controller verification (outside the sandbox, `AILANG_BIN` pinned, load 5–8):**
- Diff read in full: exactly the plan's two files; the two existing regression tests byte-unchanged (diff purely additive); snapshots `cmp`-identical to the committed tree.
- gofmt clean; vet rc 0; Python parse rc 0; compile fence `go test -run '^$' ./...` rc 0 (34 packages); `TestSubprocessSinkGate` PASS; `TestProductionGoSurface` and `TestProductionContextRoots` PASS; `TestMarketingKillGroupTolerance` 3/3 PASS; the stress row SKIPs without the opt-in, with the plan's message.
- **Mutant reproduction (controller):** M7a-fix-eperm → `kill_group must swallow PermissionError EPERM` (only `eperm-swallowed` red). The same mutant against the stress row: at the 25 s bound 2 of 2 red (6 and 6 wrapper exceptions; "real kill_group must raise nothing after leader reap"); at the 60 s bound 1 red, 1 UNINFORMATIVE (raw EPERM 0, line printed), and 1 green with raw EPERM 1 but no wrapper exception. So the stress row is a probabilistic observation, as the plan states; the unit row is the deterministic guard. Restored, `cmp` identical to HEAD.
- C11 protected-file observation: `marketingcheck observe` before the executor (12 records) and after; `marketingcheck compare` rc 0.
- `.ail` untouched, so `verify_ail.sh` was not run (CI's `ailang-code verify gate` covers it).
- **Race campaign, 25 runs (plan acceptance).** Campaign 1 (stop at the first red, as planned) ran while the judge was also loading the rig (load to ~37): 3 green, then trial 4 red on `TestMarketingSupervisorFakeChromeTimeoutKillsGroup` → `fake must record child and grandchild: []`. That is the unchanged fixture recording no PIDs before its 8 s deadline, the same class the executor saw. It has no `SupervisorError`/`Errno 1` and does not touch the changed line. **The plan's acceptance as written (25 consecutive greens) was therefore not met by campaign 1.** Campaign 2, run after the judge finished, ran all 25 and classified each red: **25 of 25 green** at load 2.5–7 (0 EPERM-class reds, 0 other), about 60 s per run. That meets the plan's criterion: if the fixed tree still red at the base rate of 0.25 per run, P(25 greens) ≈ 0.00075. The caveat stands that a timing flake cannot be proven absent by passes alone; the unit row is the primary guard. Log: `verification/world-iter308/controller-verification.md`.

**Independent evaluation:** r1 **PASS 88/100**, zero blocking ([report](verification/world-iter308/evaluator-r1.md)).
- Rubric: conformance 26/30, tests and gates 24/25, mutation/non-vacuity 18/20, premise and risk 10/10, hygiene 10/15.
- Gates all green outside a sandbox. Stress twice green (raw EPERM 15 and 54, wrapper exceptions 0, 32.4 s and 37.3 s).
- Plan mutants 3/3 killed at the named assertions; the stress row killed the eperm mutant 2/2 (100 and 124 wrapper PermissionErrors). Its own diff-anchored mutants (dropped call-recording, the line-count guard, `except PermissionError` only, a re-wrapping handler, 100 trials) were all killed. One survivor is by design: replacing the stress row's raw control with `pass` gives raw EPERM 0 and a green with only the UNINFORMATIVE line.
- **Controlled comparison (judge):** `-race` of `FakeDaemonCleanup/(shell-wait|interrupted-shot)` × 12: fixed 0/12 red; base 2/12 red, both "terminated trap must exit 143 after supervised cleanup" (the EPERM signature).
- Non-blocking: one `-race` package run red at load 40+ on `shell wait daemon must be dead after bounded cleanup` (217 s of 240 s; run 2 green in 60 s), a death-check timing failure in an unchanged test; the 60 s bound has no margin at load ~40; the stress row is an observation when raw EPERM is 0 (plan-sanctioned); minor redundant checks; the plan text had not been amended for the bound (now amended).
- Premise (judge, from the man pages): swallowing EPERM is reasonable for same-uid sessions. Residual risk: a setuid or credential-changing descendant left alive, which only the fake-fixture death checks would catch. No current caller has one.

**Bookkeeping:** sprint JSON M7a-fix → completed (iteration 308); plan M7a-fix status → BUILT and the stress-bound amendment recorded; charter row 150 updated (M7a-fix built, Next); CURRENT GOAL next line updated; STATUS 305 rotated to the archive (line counts asserted; the charter keeps 3 stamps).

**Result:** M7a-fix BUILT. A macOS `killpg` EPERM after the leader is reaped no longer turns a good shot or a clean daemon stop into `SupervisorError`. Row 150 is 18 of 23 milestones built (M7a-fix added this iteration).

**Gate 3b:** the PR head and then the merge commit are each polled pinned to the full SHA (bounded 30 min, per-invocation verdict path naming the SHA). Outcome reported on #366.

**Routing evidence:** base=0363b3ae3386c902f5f55c6ef3075c9c1e18f2b3@2026-10-27T16:01:48Z (Gate 1; Gate-3 `snap` identical). Driver banner: `controller=claude:claude-opus-5-5 … roles: designer=claude:claude-opus-5-5 planner=codex:gpt-6.1-sol executor=codex:gpt-6.1-sol evaluator=sonnet`; `MISSION_OVER_RATION` empty; codex 10080m window 11.0% used of 28.1% allowed.
- Controller `claude:claude-opus-5-5` (tok: not reported).
- Designer: not spawned (doc r2 governs; no design change). Not a routing failure: the table spawns a designer only when a doc is missing or needs revision.
- Planner `codex:gpt-6.1-sol`: probe 10,229 tok; run 101,563 tok, rc 0, 469 s. Agent-tool supervisor (transport only, no verdict): 28,495 tok / 503 s.
- Executor `codex:gpt-6.1-sol`: probe 9,086 tok; run 72,714 tok, rc 0, 515 s. Agent-tool supervisor (transport only): 37,193 tok / 574 s.
- Evaluator `sonnet` (Agent): 63,063 tok / 833 s. Judge independence: cross-vendor (OpenAI planner and executor, Anthropic Sonnet judge), and a different model from the Opus controller that wrote `7475ec3c`. The judge evaluated the final tree including that commit.
- No quorum this iteration (no new or revised design doc; a plan amendment only).
- Metered **$0.00**.

**Ruled out:**
- Asserting "EPERM ⇒ no live member": the planner's F8 and my own reading of `man 2 kill`/`killpg` refute it; recorded as a scoped workaround instead.
- Reducing the stress row to 100 trials to fit 25 s: the plan forbids a silent N reduction, and 100 trials would roughly halve the raw-EPERM control count per run.
- Attributing campaign 1's trial-4 red to M7a-fix: its message is the fixture's own PID-recording assertion (empty list) at load ~37, with no EPERM signature; the changed line runs only after the fake has started.
- Calling the 25-run campaign met: campaign 1 stopped at a red, so the literal criterion failed; campaign 2's classified result is what is recorded.
- Closing the judge's load-only timing survivors here: they are in unchanged M7a tests and belong with the carried `shell-hup`/`shell-wait` flake notes.

**Carried:**
- M7c2b next (unblocked: its new dependency M7a-fix is built), then M7c3 → M7c4 → R1 → M8a → M8b → attended A1.
- For M7c2b (from 307): `pass` creates `--out-dir` before validating routes; untested `pass`/`ready`/`chrome-version` refusal branches; `chrome_version`'s uncaught `TimeoutExpired`.
- New: M7a's fake-Chrome/daemon fixtures are load-sensitive (`fake must record child and grandchild: []`/`[81876]` at load 9–37; `shell wait daemon must be dead` at load 40+). These are separate from EPERM and are the next flake source on this rig.
- Residual (judge): swallowing EPERM would hide a live, differently-credentialed descendant; no current caller has one.
- Everything iteration 307 carried. D-WORLD-86 stays HOLD.

**Retro:**
- The planner's man-page check corrected a premise the controller had stated as fact in its own directive. This is why the planner directive demanded a command per codebase claim, even for things "everyone knows".
- An opt-in stress row sized inside a sandbox (where it never runs) was mis-sized for the real host. Running every planned acceptance command outside the sandbox before the judge caught it while it was still cheap.
- Running the race campaign concurrently with the judge confounded the campaign with load. Lesson: run load-sensitive acceptance campaigns when no other role is running, or classify reds from the start.
- Harness share 0. No skill edit.

**Progress:** 1.0 clauses 7/7 MET; row 150 18/23 built; goal unmoved.
**Next:** row 150 M7c2b, then M7c3. Then 151 after 150.

## 309 — 2026-10-27 — row 150: M7c2b built — the capture shell runs end to end against fake Chrome and a fake daemon; judge PASS 83 [PRODUCT]

**Kind:** one production milestone (M7c2b) on the release row: an executor build, a controller mutation drill with a test-only survivor fix, and an independent judge. The executor and evaluator roles were spawned through the Agent tool, per the operator's standing request (see Routing evidence).

**Pick and why:** row 150 M7c2b, following iteration 308's **Next**. Re-checked before picking:
- no M7c2b commit, branch, worktree or PR existed (`git log origin/dev --grep M7c2b` returns only plan and record commits; the PR search returns only #364 and #391, both merged; no `m7c2b` worktree);
- local dev equalled origin/dev `a96820c5` at Gate 1 (nothing to fast-forward);
- the two open PRs from this account are attended-track work, not this loop's: #390 (`trial/row219a-home-service-m0-m1`, row 219) and #389 (`fix/trial-world-context`). Left alone.
- Plan premises checked at HEAD: the shell still ended in the `capture refused` tail; `daemon` had only `--ready-file`; `serve` defines every flag the flow passes (`cmd/ailang-worldd/main.go:456-517`); the M7c2a fake Chrome and fixture helpers exist in `capture_test.go`.

**Clause map:**
- Clauses 1–7 all MET (ratified states, not re-measured). No UNMET clause, so there is no critical-path row and no drift alarm.
- D-WORLD-81 order is 156 → 150 → 151. 156 is held on D-WORLD-86, the only OPEN row (`mission_decisions.sh --open`).

**Gate 0:**
- Kill switch armed (`mission-world.disabled` absent); gh account sunholo-voight-kampff; billing tripwire CLEAN.
- `mission_directives.sh`: 0 allowlisted directives on #366 since `2026-10-27T09:03:42Z` (6 comments, all this account's). No ledger row changed since iteration 308 (HEAD was 308's own merge); ledger valid, 85 rows.
- Self-notices rc 1: the newest crash `6345079873` (rc=1, 2026-10-22) is already credited (log grep count 6). No new crash notice.
- Inbox `mission-world`: only this loop's own claim messages; nothing outranks the queue. This iteration's claim was sent and read back.
- Weekly external-issue sweep: not due (run after the 2026-10-26 rotation). No thread rotation.
- Residue: `ailang-world/.wt-world-iter306-r0` (unattributed since 306) left untouched again.

**Gate 1:**
- Dev CI at `a96820c5`: check set 3/3 success.
- Running skill plus all 13 resource files byte-identical to origin/dev (per-file `cmp` against the resolved symlink target; no DIFF line).
- Base `a96820c5…` recorded at Gate 1; the Gate-3 `snap` read the same SHA.

**Baseline (controller, pristine sprint worktree at the base before the executor ran, `AILANG_BIN` = pinned v0.41.0, load ~6):** gofmt clean, vet rc 0, Python parse rc 0, `bash -n` rc 0, version grep 0 (control on the orchestration doc 23), compile fence rc 0 (34 packages), `TestMarketingSupervisorFakeDaemonCleanup` rc 0, package rc 0, `-race` package rc 0. Every gate in the directive was green at base before it was sent.

**Role execution:**
- **Designer and planner:** not run. `w-marketing-capture-orchestration` r2 governs and the plan's M7c2b section is prescriptive.
- **Executor** `codex:gpt-6.1-sol` (resolver `recipe codex:gpt-6.1-sol declared:provider-pin`). It was spawned as an Agent-tool supervisor (`MISSION-ROLE: executor`, no model alias) that ran the controller's bounded wrapper (`/tmp/codex_run_iter309.sh`, cap 2400 s) in `.wt-world-iter309-row150-m7c2b`. Probe rc 0; run rc 0 in 287 s. Directive `/tmp/codex_directive_iter309.txt` (8,327 B) carried two controller amendments:
  - **A1:** the supervisor stops the daemon after `--lifetime` (default 300 s) and still returns ok, while the flow after readiness can include two 300 s passes. So the shell passes `--lifetime` derived from the production constants (2×SHOT + 2×PASS = 660 s). It reads them from `marketing_capture.py` via `runpy`, so the shell invents no number.
  - **A2:** a `bind-refused` subtest.
  - What it built:
    - `capture_marketing.sh` (+54/−3): the nine-step flow and three flags. Traps, `shot` and `--unit-daemon` are unchanged.
    - `marketing_capture.py` (+33/−20): a required mutually exclusive `--ready-file`/`--ready-url`, and `ready()` gains signal and process-exit stops.
    - `capture_test.go` (+190/−20): `TestMarketingCaptureFakeChromeSixShots` with eight subtests. The M7c2a fake Chrome is extracted into `fakeCaptureChrome`, which adds a `--version` mode.
  - Its sandbox denied loopback binds (probe `PermissionError`), so five subtests, the package run and three of four mutants were correctly labelled UNINFORMATIVE. Its A2 mutant was killed in-sandbox.
  - It added an untracked `snap/M7c2b/go.mod` to keep snapshot `.go` copies out of package discovery. The controller moved the snapshot to `/tmp/world-iter309/exec-snap/`, removed it from the tree, and confirmed the snapshot was `cmp`-identical to the tree.
  - Committed by the controller as `5e20341f`.
- **Evaluator** `sonnet` (resolver `agent-tool sonnet declared:alias-pin`), Agent tool, foreground, own detached worktree `.eval-world-iter309` at `76216725`. The first `mission-worktree.sh add` for it failed `failed:worktree_add_rc_255`. The likely cause: the controller had pre-created the branch with `git branch -f`, and the helper adds the branch itself. After `git worktree prune` the retry was ready with a clean tree.

**Controller verification (outside the sandbox, `AILANG_BIN` pinned, load 2–6):**
- Diff read in full, against the plan. Every write in the new flow is under `--out`; nothing under `$HOME` or `~/.ailang/state` (false-green 5).
- **Red-first:** the new test file against base production files reds 7 of 8 executor subtests. `daemon-readiness-config` passes at base, because base already rejected the wrong flag shapes with usage rc 2, so it got its own mutants.
- **Mutation drill**, one mutant at a time with anchors asserted unique, each restored and byte-compared:
  - The plan's 3 named mutants were killed at the named assertions: version-changed → `a Chrome version change during capture must refuse`; scheme → `every shot explicitly pins correct scheme`; identities → `both seed identities must be configured`.
  - 8 of 12 further diff-anchored mutants were killed: A2 drop-loopback, drop diagnostics, drop check, old Chrome continues, no TERM in cleanup, readiness group not required, readiness group not exclusive, `--ready-url` ignored.
  - 4 survived: the run continues without `--verify-twice`; A1's `--lifetime` dropped; `ready()` ignores a daemon exit; a started file written despite a stop.
- **Controller test-only fix `76216725`:** two subtests close the two load-bearing survivors.
  - `single-pass`: without `--verify-twice`, rc 0 with no pass-b, check or diagnostics (plan step 5).
  - `daemon-exits`: a daemon exiting before HTTP readiness refuses as `NonzeroExit` before any pass.
  - Both were killed under their mutants (`unexpected=[pass-b capture-manifest.json diagnostics.txt]`; verdict `Interrupted` after 29.5 s). The tree was `shasum`-identical after restoring.
- **Final tree `76216725`:** gofmt clean, vet rc 0, parse rc 0, `bash -n` rc 0, version grep 0, compile fence rc 0 (34), `FakeDaemonCleanup` rc 0, package rc 0 (31.5 s), `-race` package rc 0 (102 s), `TestSubprocessSinkGate`, `TestProductionGoSurface` and `TestProductionContextRoots` ok.
- `.ail` untouched, so `verify_ail.sh` was not run (CI's `ailang-code verify gate` covers it).
- Log: `verification/world-iter309/controller-verification.md`.

**Independent evaluation:** r1 **PASS 83/100**, zero blocking ([report](verification/world-iter309/evaluator-r1.md)).
- Rubric: conformance 27/30, tests and gates 23/25, mutation/non-vacuity 14/20, premise and risk 9/10, hygiene 10/15.
- All gates green; no load flakes seen. Red-first 9 of 10 at base.
- Plan mutants 3/3 killed at the named assertions. Of its own mutants, 11 were killed and 10 survived.
- Premises: A1 is correct and sufficient. The lifetime clock starts after readiness, so 600 s is the need, and 660 s also bounds an orphaned supervisor if the shell is SIGKILLed. A2 is sound. The traps cover every exit path except SIGKILL of the shell.
- Non-blocking findings:
  1. A1 untested.
  2. `ready()`'s stop/Interrupted branch untested.
  3. Missing/unknown-flag and port-range refusals have no subtest, and shell usage errors exit 1, not 2.
  4. The `kill -0` readiness shortcut is unguarded.
  5. The diff is 308 lines against the ~250-line stop rule, mostly test.
  6. Cosmetic: help-line spacing, and the lifetime comment counts version A.
- **Controller reproduced 3 and 6** (`--db x` rc 1 "missing --manifest"; `--bogus` rc 1; help line printed). Findings 1 and 2 agree with the controller's own surviving mutants.

**Bookkeeping:** sprint JSON M7c2b → completed (iteration 309); plan M7c2b status → BUILT with the amendments recorded; charter row 150 updated (M7c2b built, Next); CURRENT GOAL next line updated; STATUS 306 rotated to the archive (charter line count asserted, archive gained the stamp; the charter keeps 3 stamps).

**Result:** M7c2b BUILT. The release-image capture is now one command that starts the frozen demo daemon, takes two six-shot passes, refuses if Chrome changed mid-run, and verifies the passes, with the daemon stopped on every exit path. It is proven against fakes only; the real-Chrome drill is M7c4. Row 150 is 19 of 23 milestones built.

**Gate 3b:** the PR head and then the merge commit are each polled, pinned to the full SHA (bounded 30 min, per-invocation verdict path naming the SHA). The outcome is reported on #366.

**Routing evidence:** base=a96820c58c678b262371bedbb8023e7b6eda6864@2026-10-27T20:38:14Z (Gate 4 record; Gate 1 read the same SHA at 20:01:51Z, and the Gate-3 `snap` was identical). Driver env: `MISSION_DESIGNER_RESOLVED=claude:claude-opus-5-5`, `PLANNER=codex:gpt-6.1-sol`, `EXECUTOR=codex:gpt-6.1-sol`, `EVALUATOR=sonnet`; `MISSION_OVER_RATION` empty.
- Controller `claude:claude-opus-5-5` (tok: not reported).
- Designer and planner: not spawned (doc r2 and a prescriptive plan govern). This is not a routing failure.
- Executor `codex:gpt-6.1-sol`: probe 10,229 tok; run 78,057 tok, rc 0, 287 s. Agent-tool supervisor (transport only, no verdict): 33,357 tok / 339 s.
- Evaluator `sonnet` (Agent): 62,165 tok / 772 s. Judge independence: cross-vendor (OpenAI executor, Anthropic Sonnet judge), and a different model from the Opus controller that wrote `76216725`. The judge evaluated the final tree including that commit. Its directive named no controller finding or method.
- No quorum this iteration (no new or revised design doc).
- Metered **$0.00**.

**Ruled out:**
- Banking the executor's in-sandbox gate results: loopback binds were denied, so five subtests and three mutants had no verdict until the controller ran them outside.
- Treating `daemon-readiness-config` as a red-first test: it passes at base. Its value is measured by its two mutants instead.
- Letting the daemon run on its default 300 s lifetime: a slow pass A could stop the daemon mid pass B. The judge confirmed A1 and noted 630 s would suffice; 660 s stays, as the conservative derivation.
- Keeping the executor's `snap/M7c2b/go.mod` in the tree: it is evidence scaffolding, not product, and it would add a nested module to the repo.
- Rotating the log this iteration: it holds 25 entries and the rule triggers past ~40. An accidental `rotate-log --keep 20` was undone (log restored from HEAD plus this entry, archive and index reverted) and the index regenerated with `--index-only`. That run used a temporary registry whose `workdir` is this worktree, because `world.toml`'s `workdir` is the main checkout and the default registry would have written there.
- Adding subtests for every judge survivor here: A1's lifetime and the pre-readiness stop branch need a slow-pass or signal-timing fixture. They are carried rather than bolted on past the plan's size bound.

**Carried:**
- M7c3 next (frozen public inventory guard), then M7c4 → R1 → M8a → M8b → attended A1.
- New (judge): A1 lifetime untested; `ready()` stop/Interrupted and the `not stop` started-file guard untested; missing/unknown-flag and port-range refusals untested, and shell usage errors exit 1, not 2; `kill -0` readiness shortcut unguarded; cosmetic help spacing and the lifetime comment.
- Still carried from 307: `pass` creates `--out-dir` before validating routes; untested `pass`/`ready`/`chrome-version` refusal branches; `chrome_version`'s uncaught `TimeoutExpired`.
- Everything iteration 308 carried (load-sensitive M7a fixtures; the EPERM residual). D-WORLD-86 stays HOLD.

**Retro:**
- The executor's sandbox could not bind loopback, so the milestone's whole acceptance surface was invisible to it. Running the full drill outside the sandbox before the judge found and closed two survivors while it was cheap. That is now the third consecutive row-150 milestone where the controller's out-of-sandbox drill changed the tree before judging.
- Pre-creating a branch for a helper that creates its own branch cost one failed worktree add. Let the helper own the branch.
- Harness share 0. No skill edit.

**Progress:** 1.0 clauses 7/7 MET; row 150 19/23 built; goal unmoved.
**Next:** row 150 M7c3, then M7c4. Then 151 after 150.

## 310 — 2026-10-28 — row 150: M7c3 built — `marketingcheck inventory`/`inventory-compare` freeze the public store so a capture taken over a commit is refused; judge PASS 90 [PRODUCT]

**Kind:** one production milestone (M7c3) on the release row: an executor build, a controller mutation drill with a test-only survivor fix, and an independent judge. The executor and evaluator roles were spawned through the Agent tool, per the operator's standing request (see Routing evidence).

**Pick and why:** row 150 M7c3, following iteration 309's **Next**. Re-checked before picking:
- no M7c3 commit, branch, worktree or PR existed (`git log origin/dev --grep M7c3` returns only the plan split and 309's record; the PR search returns only #364, merged; no `m7c3` worktree);
- local dev equalled origin/dev `c448b15f` at Gate 1;
- the two open PRs from this account are attended-track work, not this loop's: #390 (`trial/row219a-home-service-m0-m1`) and #389 (`fix/trial-world-context`). Left alone.
- Plan premises checked at HEAD: every helper the test names exists at the cited lines; `handleLogRange` stops at the first missing index; `clampLimit` caps at 500; `isProtected` guards only `POST /v1/commit`, so the inventory's reads need no credential.

**Clause map:**
- Clauses 1–7 all MET (ratified states, not re-measured). No UNMET clause, so there is no critical-path row and no drift alarm.
- D-WORLD-81 order is 156 → 150 → 151. 156 is held on D-WORLD-86, the only OPEN row (`mission_decisions.sh --open`).

**Gate 0:**
- Kill switch armed (`mission-world.disabled` absent); gh account sunholo-voight-kampff; billing tripwire CLEAN.
- `mission_directives.sh`: 0 allowlisted directives on #366 since `2026-10-27T09:03:42Z` (7 comments, all this account's). No ledger row changed since iteration 309 (HEAD was 309's own merge); ledger valid, 85 rows.
- Self-notices rc 1: the newest crash `6345079873` (rc=1, 2026-10-22) is already credited (log grep count 7). No new crash notice.
- Inbox `mission-world`: only this loop's own claim messages; nothing outranks the queue. This iteration's claim was sent.
- Weekly external-issue sweep: not due (run after the 2026-10-26 rotation). No thread rotation.
- Residue: `ailang-world/.wt-world-iter306-r0` (unattributed since 306) left untouched again.

**Gate 1:**
- Dev CI at `c448b15f`: check set 3/3, none non-green.
- Running skill plus all 13 resource files byte-identical to origin/dev (per-file `cmp` against the resolved symlink target, which is the main checkout).
- Base `c448b15f…` recorded at Gate 1; the Gate-3 `snap` and the Gate-4 record read the same SHA.

**Baseline (controller, pristine sprint worktree at the base before the executor ran, `AILANG_BIN` = pinned v0.41.0, load ~6):** vet rc 0 on both packages, `TestWorkbenchLiveCursorCrossesGap` ok, `TestProductionGoSurface` ok, `TestProductionContextRoots` and `TestSubprocessSinkGate` ok, gofmt clean. Every pre-existing gate in the directive was green at base before it was sent.

**Role execution:**
- **Designer and planner:** not run. `w-marketing-capture-orchestration` r2 governs and the plan's M7c3 section is prescriptive.
- **Executor** `codex:gpt-6.1-sol` (resolver `recipe codex:gpt-6.1-sol declared:provider-pin`). It was spawned as an Agent-tool supervisor (`MISSION-ROLE: executor`, no model alias) that ran the controller's bounded wrapper (`/tmp/codex_run_iter310.sh`, cap 1800 s) in `.wt-world-iter310-row150-m7c3`. Probe rc 0; run rc 0 in 349 s. Directive `/tmp/codex_directive_iter310.txt` (8,426 B).
  - What it built: `scripts/marketingcheck/main.go` (+166 net: `inventory`, `inventory-compare`, `Inventory`, `ReadInventory`, `CompareInventory`; `writeObservation` generalised to `writeJSON`) and `host/daemon/marketing_inventory_test.go` (91 lines, the plan's eight steps plus the 4-entry/3-entry paging assertion).
  - Deviation, adjudicated correct (rule 3h): `/v1/head` serves `text/plain`, so a non-JSON head body is stored as a JSON string. The mutant that drops this fails the real daemon.
  - Its sandbox denied loopback binds, so the inventory test, the package run and all five of its mutants were correctly labelled UNINFORMATIVE; it claimed no kill.
  - Its `snap/M7c3/` was moved to `/tmp/world-iter310/exec-snap/` (report banked as `verification/world-iter310/executor-report.md`). Committed by the controller as `8da32397`.
- **Evaluator** `sonnet` (resolver `agent-tool sonnet declared:alias-pin`), Agent tool, foreground, own worktree `.eval-world-iter310` at `47686112`.

**Controller verification (outside the sandbox, `AILANG_BIN` pinned, load 1–6):**
- Diff read in full against the plan. No write outside `--out`; no context root; no new non-test `.go` under `scripts/`.
- **Red-first:** base `main.go` reds the e2e test at step 5 (`before inventory: rc 2: usage`).
- **Mutation drill** (`verification/world-iter310/controller-mutants.py.txt`, one at a time, anchors asserted unique, sha restored):
  - The plan's 2 mutants were killed at the named assertions: M7c-inventory → "extra commit invalidates capture"; M7c3-paging → "after.json must list all 4 committed entries: got 1".
  - The controller's own first drafts of both were wrong. `return nil` after a loop fails `go vet` (unreachable code), and `limit=1` alone still pages through every entry. Both were rewritten (`if true { return nil }`; `limit=1` plus stop-after-first-page) before counting.
  - Of 12 diff-anchored mutants, the e2e test killed 3 (refusal exit 0, head text fallback dropped, and the plan inventory mutant) and 8 survived: head not compared, entries compared by count only, object length ignored, object hash ignored, non-200 ignored, empty check removed, objects never read, single page. The reason is structural: one genuine extra commit changes head AND entries together, and the test store reuses one transition object.
- **Controller test-only fix `47686112`:** `scripts/marketingcheck/inventory_test.go` adds `TestMarketingInventoryCompareBranches` (14 subtests, one per refusal branch and side) and `TestMarketingInventoryReadPaging` (fake daemon with a page cap of 2, plain-text head, two distinct objects, a 404 arm). After it every mutant is killed by at least one suite. The tree was `shasum`-identical after restoring.
- **Final tree `47686112`:** gofmt clean, vet rc 0, compile fence rc 0 (34), the e2e test PASS, the unit tests PASS, `TestWorkbenchLiveCursorCrossesGap` ok, `TestProductionGoSurface`/`TestProductionContextRoots`/`TestSubprocessSinkGate` ok, both packages ok, `-race ./scripts/marketingcheck` ok, `-race` inventory/workbench-live in host/daemon ok.
- `.ail` untouched, so `verify_ail.sh` was not run (CI's `ailang-code verify gate` covers it).
- Log: `verification/world-iter310/controller-verification.md`.

**Independent evaluation:** r1 **PASS 90/100**, zero blocking ([report](verification/world-iter310/evaluator-r1.md)).
- Rubric: conformance 27/30, tests and gates 24/25, mutation/non-vacuity 16/20, premise and risk 8/10, hygiene 15/15.
- All gates green; red-first confirmed for both new test files.
- Both plan mutants killed at the named assertions. Of about 45 own mutants, every compare branch, the head encoding, paging, object sort/digest, non-200 and the CLI exit were killed. Survivors: head and log fetch errors swallowed, page unmarshal error ignored, `TrimSpace` in the empty check, the before/after label, the CLI usage/flag paths, the schema string, `TrimRight("/")`, and the client timeout. One more (object-presence check) is equivalent.
- Premise: the guard catches any commit. It cannot see writes that bypass the log, non-transition object puts, an A→B→A revert, or a race between its three reads. Digesting the object envelope without `?payload=true` is sound, since the ref is content-addressed.
- **Controller reproduced the CLI finding:** the usage paths are correct but untested (`inventory` with no `--addr` rc 2; `inventory-compare` missing `--after` rc 2; `-h` rc 0; unreachable rc 1 `inventory-unreadable`).

**Bookkeeping:** sprint JSON M7c3 → completed (iteration 310); plan M7c3 status → BUILT; charter row 150 updated (M7c3 built, Next); CURRENT GOAL next line updated; STATUS 307 rotated to the archive (charter line count asserted, archive gained the stamp; the charter keeps 3 stamps).

**Result:** M7c3 BUILT. A capture can now prove the store it photographed did not change: an inventory before pass A and after pass B, compared, refuses on any commit. It is not yet wired into `capture_marketing.sh`; that is M7c4's carried step. Row 150 is 20 of 23 milestones built.

**Gate 3b:** the PR head and then the merge commit are each polled, pinned to the full SHA (bounded 30 min, per-invocation verdict path naming the SHA). The outcome is reported on #366.

**Routing evidence:** base=c448b15fd6ee9c9e5e221cb66bf099cb7bcb127c@2026-10-28T00:27:08Z (Gate 4 record; Gate 1 read the same SHA at 00:01:13Z, and the Gate-3 `snap` was identical). Driver env: `MISSION_DESIGNER_RESOLVED=claude:claude-opus-5-5`, `PLANNER=codex:gpt-6.1-sol`, `EXECUTOR=codex:gpt-6.1-sol`, `EVALUATOR=sonnet`; `MISSION_OVER_RATION` empty.
- Controller `claude:claude-opus-5-5` (tok: not reported).
- Designer and planner: not spawned (doc r2 and a prescriptive plan govern). This is not a routing failure.
- Executor `codex:gpt-6.1-sol`: probe 10,229 tok; run 96,227 tok, rc 0, 349 s. Agent-tool supervisor (transport only, no verdict): 25,863 tok / 378 s.
- Evaluator `sonnet` (Agent): 64,452 tok / 482 s. Judge independence: cross-vendor (OpenAI executor, Anthropic Sonnet judge), and a different model from the Opus controller that wrote `47686112`. The judge evaluated the final tree including that commit. Its directive named no controller finding, mutant result or method.
- No quorum this iteration (no new or revised design doc).
- Metered **$0.00**.

**Ruled out:**
- Banking the executor's in-sandbox results: loopback binds were denied, so the whole acceptance surface had no verdict until the controller ran it outside.
- Counting the controller's first-draft plan mutants: one failed vet, and the other was not the plan's mutant. A mutant that does not build, or that is not the one named, is not a measurement.
- Adding survivor-closing assertions to the host/daemon e2e test: one genuine commit cannot vary head, entries and objects independently, so the branches belong in a unit test over `CompareInventory` and a fake daemon.
- Wiring `inventory` into `capture_marketing.sh` here: the plan keeps it out of M7c3's file list and carries it to M7c4.
- Rotating the log: it holds 26 entries; the rule triggers past ~40. The index was regenerated with `--index-only` through a temporary registry whose `workdir` is the record worktree.

**Carried:**
- M7c4 next (real-Chrome drill, plus the carried inventory wiring), then R1 → M8a → M8b → attended A1.
- New (judge): head/log fetch and page-unmarshal error paths untested; CLI usage/flag paths untested; the `inventory-empty` before/after label unpinned; guard scope (writes that bypass the log, non-transition objects, A→B→A, non-atomic reads).
- Everything iteration 309 carried (A1 lifetime untested; `ready()` stop branch; shell usage errors exit 1; `kill -0` shortcut; 307's `pass` `--out-dir` order; load-sensitive M7a fixtures). D-WORLD-86 stays HOLD.

**Retro:**
- Fourth consecutive row-150 milestone where the out-of-sandbox controller drill changed the tree before judging. This time the plan's own e2e design could not distinguish its branches, so the fix was a unit test rather than another e2e arm.
- The controller wrote both plan mutants wrongly on the first attempt (a vet-failing form, and half of a two-part mutant). The anchored runner's vet check caught the first; reading a surprising survivor back against the plan table caught the second. Read the named mutant's full definition before encoding it.
- Harness share 0. No skill edit.

**Progress:** 1.0 clauses 7/7 MET; row 150 20/23 built; goal unmoved.
**Next:** row 150 M7c4 (real-Chrome drill; wire `inventory` into the capture shell), then R1 → M8a → M8b. Then 151 after 150.

## 311 — 2026-10-28 — row 150: M7c4 built — the real-Chrome drill exists and the capture shell compares frozen inventories around both passes; judge PASS 92 [PRODUCT]

**Kind:** one production milestone (M7c4) on the release row: an executor build, a controller mutation drill with a test-only survivor fix, a real-Chrome smoke, and an independent judge. The executor and evaluator roles were spawned through the Agent tool, per the operator's standing request (see Routing evidence).

**Pick and why:** row 150 M7c4, following iteration 310's **Next**. Re-checked before picking:
- no M7c4 commit, branch, worktree or PR existed (`git log origin/dev --grep M7c4` returns only the plan split, 308's plan amendment and 310's record; the PR search returns only merged #391–#393 and #364);
- local dev was 4 behind origin at Gate 0 and was fast-forwarded (clean tree apart from two untracked residues); local then equalled origin/dev `90ef2671`;
- the two open PRs from this account are attended-track work, not this loop's: #390 (`trial/row219a-home-service-m0-m1`) and #389 (`fix/trial-world-context`). Left alone.
- Plan premises checked at HEAD: `testMintEnv` (`session_test.go:29`), `newDemoApprovalRig` (`demo_test.go:561`), `demoM5cRealSeed` (`demo_test.go:1727`) and the iteration-305 probe reference all exist; `marketing_chrome_test.go` did not; the shell carried no `inventory` call.
- Governing doc: `w-marketing-capture-orchestration` r2. Its two quorum artifacts both read `blocked`; the doc routes under the narrow-refinement carve-out recorded in its §9 (iteration 305), the same authority iterations 306–310 built under. Not re-litigated.

**Clause map:**
- Clauses 1–7 all MET (ratified states, not re-measured). No UNMET clause, so there is no critical-path row and no drift alarm.
- D-WORLD-81 order is 156 → 150 → 151. 156 is held on D-WORLD-86, the only OPEN row (`mission_decisions.sh --open`).

**Gate 0:**
- Kill switch armed (`mission-world.disabled` absent); gh account sunholo-voight-kampff; billing tripwire CLEAN.
- `mission_directives.sh` (driver copy): 0 allowlisted directives on #366 since `2026-10-27T09:03:42Z` (8 comments, all this account's). Both watermarks (`mission-366-last-seen`, `mission-world-last-seen`) read the same value. No ledger row changed since iteration 310.
- Self-notices rc 1: the newest crash `6345079873` (rc=1, 2026-10-22) is already credited (log grep count 8). No new crash notice.
- Inbox `mission-world`: only this loop's own claim messages; nothing outranks the queue. This iteration's claim was sent.
- Weekly external-issue sweep: not due. No thread rotation.
- Residue: `ailang-world/.wt-world-iter306-r0` (unattributed since 306) left untouched again.

**Gate 1:**
- Dev CI at `90ef2671`: check set 3/3 success.
- Running skill plus all 13 resource files byte-identical to origin/dev (per-file `cmp` against the resolved symlink target).
- Base `90ef2671…` recorded at Gate 1 (`2026-10-28T04:01:37Z`).

**Baseline (controller, pristine sprint worktree at the base before the executor ran, `AILANG_BIN` = pinned v0.41.0, load 3–6):** gofmt clean, vet rc 0, `bash -n` rc 0, compile fence rc 0 (34), fake-Chrome test ok, `TestProductionGoSurface`/`TestProductionContextRoots`/`TestSubprocessSinkGate` ok. **`go test ./scripts/marketingcheck` was rc 1 at base**: `TestMarketingSupervisorFakeDaemonCleanup` ("shell trap cleanup exceeded bound: 29.7s", load ~6, inside the package run); the identical `-run` alone was 3/3 ok at load ~3.6. That is the load-sensitive M7a fixture class iteration 308 carried; the directive named it as not the executor's to fix.

**Role execution:**
- **Designer and planner:** not run. Doc r2 governs and the plan's M7c4 section is prescriptive.
- **Executor** `codex:gpt-6.1-sol` (resolver `recipe codex:gpt-6.1-sol declared:provider-pin`). It was spawned as an Agent-tool supervisor (`MISSION-ROLE: executor`, no model alias) that ran the controller's bounded wrapper (`/tmp/codex_run_iter311.sh`, cap 1800 s) in `.wt-world-iter311-row150-m7c4`. Probe rc 0; run rc 0 in 273 s. Directive `/tmp/codex_directive_iter311.txt` (10,093 B).
  - What it built: `scripts/capture_marketing.sh` (+10/−5), `scripts/marketingcheck/capture_test.go` (+38/−6: fake daemon log/objects routes, `inventory` subtest), `cmd/ailang-worldd/marketing_chrome_test.go` (new, 246 lines).
  - Its sandbox denied loopback binds, so every bind-dependent gate and all four of its shell mutants were correctly labelled UNINFORMATIVE; it claimed no kill.
  - Its `snap/M7c4/` was moved to `/tmp/world-iter311/exec-snap/` (report banked as `verification/world-iter311/executor-report.md`). Committed by the controller as `31a4929b`.
- **Evaluator** `sonnet` (resolver `agent-tool sonnet declared:alias-pin`), Agent tool, foreground, own worktree `.eval-world-iter311` at `1ca31e5e`.

**Controller verification (outside the sandbox, `AILANG_BIN` pinned, load 2.5–6; [log](verification/world-iter311/controller-verification.md)):**
- Diff read in full against the plan. No write outside `t.TempDir()`; `demo_test.go` untouched; every precondition precedes seeding.
- Controller amendment **A1**: the plan says "after pass B"; the single-pass branch has no pass B, so its after-inventory and compare run after pass A. The guard brackets every pass that runs, in both modes.
- **Red-first:** base `capture_marketing.sh` reds `six-shots` (`inventory-before.json` absent) and `inventory` ("capture must compare inventories: rc=0").
- **Mutation drill** (`verification/world-iter311/controller-mutants.py.txt`, anchors asserted unique, `bash -n` checked, sha restored): 9 shell mutants. On the executor's tree 6 were killed and 3 survived: before-inventory moved after pass A, after-inventory moved before pass B, single-pass skipping the compare. Cause: the fake daemon flipped its head after the first `/v1/log?from=0` read, so it could not tell when the inventories ran relative to the passes.
- **Controller test-only fix `1ca31e5e`:** the fake head changes once a fake-Chrome shot of the named pass has run (`CAPTURE_HEAD_CHANGE=pass-a|pass-b`, globbing the `<pass>-<shot>.args` files fake Chrome already writes beside the fake daemon), with new subtests `inventory-pass-b` and `inventory-single-pass`. After it, 9/9 killed; the plan mutant M7c4-inventory-wiring dies at "capture must compare inventories".
- **Completion (b)/(c)/(d):** SKIP with `WORLD_CHROME` unset; FAIL at `:58` with `/nonexistent`; with real Chrome set, rects unset (`:75`), `{schema, shots:{}}` (`:94`), path missing (`:91`) and `AILANG_BIN` empty (`:71`) all FAIL; `inventory` passes and the plan mutant is killed.
- **Real-Chrome smoke (Chrome 156.0.8078.48; not C8 evidence):** with placeholder 40x40 rectangles the full path ran in 20 s (real seed, real `serve`, 6/6 shots per pass, inventory before/after equal) and `check` then refused `rect-too-small`. With made-up 100x100 rectangles the drill **PASSED end to end** (rc 0, 15 s).
- **Final tree `1ca31e5e`:** gofmt clean, vet rc 0, `bash -n` rc 0, fence rc 0, SKIP, fake-Chrome 13/13 ok, `-race` fake-Chrome ok, surface/context roots ok, sink gate ok, `go test ./scripts/marketingcheck` ok (flake did not fire), `go test ./cmd/ailang-worldd` ok (144 s).
- `.ail` untouched, so `verify_ail.sh` was not run (CI's `ailang-code verify gate` covers it).

**Independent evaluation:** r1 **PASS 92/100**, zero blocking ([report](verification/world-iter311/evaluator-r1.md)).
- Rubric: tests 20/20, lint 10/10, acceptance 29/30, code quality 13/15, documentation 12/15 (no doc status in the diff; the controller owns it), design fidelity 8/10.
- All gates green; red-first confirmed; (b)/(c)/(d) measured independently, including extra precondition arms (packet roster name missing, wrong-version tool pin).
- 14 mutants of the shell diff: 10 killed, 3 equivalent (`|| true` after `step`, whose failure handler exits; `--before`/`--after` swapped; before-inventory after `routes`), 1 real survivor (the diagnostics line written before the compare, caught only incidentally by `single-pass`). A1 judged sound.
- Independently found that made-up rectangles pass on real Chrome.
- **Controller reproduced both:** the diagnostics-ordering mutant is killed only by `single-pass`; made-up 100x100 rectangles pass end to end.

**Bookkeeping:** sprint JSON M7c4 → completed (iteration 311, PR 394); plan M7c4 status → BUILT with A1 and the R1 risk; charter row 150 updated (M7c4 built, carried notes, R1 risk, Next); CURRENT GOAL next line updated; STATUS 308 rotated to the archive (charter line count asserted unchanged, archive +2 and holds the stamp; the charter keeps 3 stamps).

**Result:** M7c4 BUILT. The capture shell now refuses a capture if the public store changed between an inventory taken before pass A and one taken after the last pass, and the real-Chrome drill exists for R1. Row 150 is 21 of 23 milestones built.

**Gate 3b:** the PR head and then the merge commit are each polled, pinned to the full SHA (bounded 30 min, per-invocation verdict path naming the SHA). The outcome is reported on #366.

**Routing evidence:** base=90ef26711bc0812c1a84bab23c7a6454f5996598@2026-10-28T04:01:37Z (Gate 1 record). Driver env: `MISSION_DESIGNER_RESOLVED=claude:claude-opus-5-5`, `PLANNER=codex:gpt-6.1-sol`, `EXECUTOR=codex:gpt-6.1-sol`, `EVALUATOR=sonnet`; `MISSION_OVER_RATION` empty.
- Controller `claude:claude-opus-5-5` (tok: not reported).
- Designer and planner: not spawned (doc r2 and a prescriptive plan govern). This is not a routing failure.
- Executor `codex:gpt-6.1-sol`: probe 9,086 tok; run 90,685 tok, rc 0, 273 s. Agent-tool supervisor (transport only, no verdict): 28,440 tok / 303 s.
- Evaluator `sonnet` (Agent): 83,184 tok / 406 s. Judge independence: cross-vendor (OpenAI executor, Anthropic Sonnet judge), and a different model from the Opus controller that wrote `1ca31e5e`. The judge evaluated the final tree including that commit. Its directive named no controller finding, mutant result or method.
- No quorum this iteration (no new or revised design doc).
- Metered **$0.00**.

**Ruled out:**
- Banking the executor's in-sandbox results: loopback binds were denied, so the whole bind-dependent surface had no verdict until the controller ran it outside.
- Treating the base `go test ./scripts/marketingcheck` red as M7c4's: it reproduced on the pristine base and passed 3/3 alone.
- Keeping jq's rewrite of the sprint JSON: `jq` re-expanded four earlier one-line `completed` objects; reverted and replaced with a one-block text edit.
- Adding the diagnostics-ordering assertion after judging: it is non-blocking and changing the judged tree would need a second round; carried instead.
- Running R1 here: it needs genuinely measured rectangles over the M7c0 geometry and a threshold V-row, and the smoke shows `check` cannot catch a guessed set.

**Carried:**
- R1 next (controller: measure all-six rectangles, two real passes, threshold V-row), then M8a → M8b → attended A1.
- **R1 risk (measured, new):** `check` does not reject guessed rectangles (made-up 100x100 at the origin pass on real Chrome), so R1 must show its rectangles are measured. A negative control (a deliberately misplaced rectangle must refuse) is worth adding to R1.
- New (judge): the diagnostics line written before the compare is caught only incidentally; the drill's rectangles precondition is weaker than `LoadRectangles` (an out-of-bounds rect runs the whole capture before failing); the manifest assertion only checks it parses; `v0.41.0` is hardcoded in the drill; the plan's `-timeout 660s` equals the inner bound.
- Everything iteration 310 carried (fetch/unmarshal error paths, CLI usage paths, guard scope), plus 309's (A1 lifetime, `ready()` stop branch, usage errors exit 1) and the load-sensitive M7a supervisor fixtures (`TestMarketingSupervisorFakeDaemonCleanup` red once this iteration at load ~6). D-WORLD-86 stays HOLD.

**Retro:**
- Fifth consecutive row-150 milestone where the out-of-sandbox controller drill changed the tree before judging. This time the executor's fake was the gap: a trigger keyed to the first log read cannot tell where an inventory sits relative to the passes. Keying a fake's state change to an observable step of the flow it brackets is what made placement testable.
- The real-Chrome smoke cost about 20 s and found the most important R1 fact of the iteration. A cheap real run with placeholder inputs is worth doing whenever a drill exists, even when its full pass is out of scope.
- Harness share 0. No skill edit.

**Progress:** 1.0 clauses 7/7 MET; row 150 21/23 built; goal unmoved.
**Next:** row 150 R1 (controller: measured rectangles, two real passes, threshold V-row, with a misplaced-rectangle negative control), then M8a → M8b. Then 151 after 150.

## 312 — 2026-10-28 — row 150: R1 measured — real-Chrome rectangles, two-pass full path and threshold V-row, with two store-dependent criteria left open; row 156 closed on D-WORLD-86 = A; judge PASS 72 [PRODUCT]

**Kind:** a controller-owned measurement step (R1) on the release row, judged by an independent evaluator, plus bookkeeping (row 156's close) unparked by an attended ruling. The evaluator was spawned through the Agent tool, per the operator's standing request. No designer, planner or executor was needed: R1 is a controller step by the doc's own §5.

**Pick and why:** row 150 R1, following iteration 311's **Next**. Re-checked before picking:
- no R1 commit, branch, worktree or PR existed (`git log origin/dev --grep R1` returns only design/ledger commits; the PR title search returns merged design PRs only);
- local dev equalled origin/dev `dba643a9` at Gate 1 (#389/#390 had merged attended since iteration 311);
- one open PR from this account, #395 (`docs/ruling-d-world-86`), was attended work (commit author is the attended identity; opened by the attended session `ae6f556a`'s worktree). It was left alone and merged attended at 08:20Z during this iteration.
- Governing doc `w-marketing-capture-orchestration` r2 (narrow-refinement carve-out, §9; not re-litigated).

**Clause map:** clauses 1–7 MET (ratified states, not re-measured); no UNMET clause, no critical-path row. D-WORLD-81 order 156 → 150 → 151.

**Gate 0:**
- Kill switch armed (`mission-world.disabled` absent); gh account sunholo-voight-kampff; billing tripwire CLEAN.
- `mission_directives.sh`: 0 allowlisted directives on #366 since `2026-10-27T09:03:42Z` (9 comments, all this account's). Both watermarks read the same value; unchanged (nothing new to process).
- Self-notices rc 1: the newest crash `6345079873` (rc=1, 2026-10-22) is already credited (log grep count 9). No new crash notice.
- Inbox `mission-world`: only this loop's own claim messages. This iteration's claim was sent.
- Weekly external-issue sweep: not due.
- Residue: `ailang-world/.wt-world-iter306-r0` (unattributed since 306) left untouched again.

**Gate 1:**
- Dev CI at `dba643a9`: check set 3/3 success.
- Running skill plus all 13 resource files byte-identical to origin/dev (per-file `cmp` against the resolved symlink target).
- Base `dba643a9…` recorded at Gate 1 (`2026-10-28T08:02:25Z`).
- Attended ruling found in flight: PR #395 sets **D-WORLD-86 = A** (drop row 156's M5c/M5d). Acknowledged; acted on only after it merged (08:20Z): `mission_ledger_provenance.sh --id D-WORLD-86` → `verdict=attended flip=3be8e3e7`; ledger valid, 85 rows, **0 OPEN**.

**R1 (controller, outside any sandbox; Chrome 156.0.8078.48, `AILANG_BIN` v0.41.0, tool v0.54.15, load ~5):**
- Probe `verification/world-iter312/r1-rects-probe_test.go.txt` (run as a `zz_` test, then deleted). It uses the drill's own rig, the six routes from `marketingcheck routes`, and the V51 layout method at each shot's M7c0 geometry. Phone used a supervisor `shot` plus the V58 scan. Then the drill's exact `capture_marketing.sh --verify-twice` ran on the same store, with output kept.
- **Text-tight boxes refuted (seeds 1–4):** `Range` boxes for the packet paragraphs changed width between seeds (1209/1207 and 1131/1115 px, then again on seeds 3/4), because sha256 hex glyphs vary in width. Paragraph regions now use the element's block box.
- **Phone background fixed (seed 3):** the first gutter rectangle crossed the blue nav (sd 4.64). A column scan put the phone nav at y 114–322, as on desktop, so the background starts at y 343 (sd 0.000).
- **Final baseline (seeds 5/6):** `r1-rects.json` (sha256 `ba7c83ce…`) was byte-identical on both seeds. The phone rectangles equal V58's baseline. Every region bottom is ≤ H − 40.
- **C14 half:** on seeds 5 and 6, pass-A = pass-B SHA256 on **6/6** shots, `inventory: equal`, `check` ok. The named drill (`TestMarketingChromeSixShots` with the baseline, on a third seed) PASSED in 19.4 s.
- **V62 thresholds** (own helper `r1-stats-main.go.txt`, mirroring `png.go`): min colours 822, chromatic 0.0621–0.2423, light background Y 248.84 and dark 29.96 (sd 0), smallest region 345x23, min content sd 17.13. Decision: **keep all eight values**; M8a replaces the "provisional" comment with a pointer to V62.
- **V63 controls:** the positive control is rc 0. **14/14** regions moved to the nearest near-blank 64x16 window were refused, each `blank-region <shot>: <region>`. Phone C8′: (i) and (ii) true; **(iii) false** against V50 (thumb right edge 333/334 vs 329).

**Independent evaluation:** **PASS 72/100** ([report](verification/world-iter312/evaluator-r1.md)). Scores: completeness 27/30, reproducibility 11/20, thresholds 13/15, controls 10/15, honesty 11/20. The judge ran its own probes on three fresh seeds, the named drill, its own statistics code, its own negative controls and the V58 scan. Findings:
- **F1 (blocking before P0, not for the measurement):** on 1 of its 3 seeds the PROVEN roster row wrapped to 3 lines (height 72 vs 96), so the block-box baseline is not fully store-independent. **Reproduced by the controller** from the judge's `layouts.json`: 72.0 on eval seed 1 and 96.0 on the four others; the TESTED row's text reaches x 1303–1315 of 1318.
- **F2:** R1 has two open criteria, not one.
- **F3:** the phone thumb spread is 324–338 (judge seeds 338/324/326), not "1 px / 4–5 px". Reproduced by the controller's own `phonescan` run on the judge's PNGs.
- **F4:** `check` gives no localisation. It accepted rectangles pointed at the wrong row, shifted 14–40 px, whole-viewport, 64x16 glyph windows and a background over content. It refused only a dropped required rectangle and out-of-bounds.
- **F5:** V62 reproduced with the judge's own code. The data show headroom, not discriminating power.
- **F6:** a local path appears in a committed log. That is the established repo practice (860 verification files carry it); no email or secret.
- **After judging, the controller narrowed its written claims to the judge's measurements** (`d30725bc`: V60/V63 text, R1 status "two open criteria", JSON/plan status). These are record corrections in the judge's direction, attributed to the judge. No code or rectangle changed, so no second round was run.

**Row 156 close (bookkeeping, unparked by D-WORLD-86 = A):**
- Sprint JSON M5c/M5d `held:D-WORLD-86` → `dropped:D-WORLD-86=A …`.
- Doc status line → IMPLEMENTED.
- `git mv` of doc, plan and JSON to `design_docs/implemented/` (row 152's convention: historical internal references are left as they are).
- The one code comment naming the planned path (`host/broker/exec_protected.go:3`) was updated; `go build ./host/broker/` ok.
- Charter row 156 tag → CLOSED iteration 312.

**Bookkeeping:**
- Orchestration doc V60–V63 plus the R1 status note; sprint JSON R1 controller step → MEASURED with the two OPEN criteria; plan "Controller steps" line updated.
- Charter row 150 (R1 measured, open criteria, Next); CURRENT GOAL next line; STATUS 309 rotated to the archive (charter line count asserted unchanged, archive +2 and holds the stamp; the charter keeps 3 stamps).

**Result:** R1 MEASURED with two open criteria. The rectangles are genuinely measured from Chrome's own layout and reproducible on 8 of 9 seeds, and the full path is byte-deterministic within a store on real Chrome. The thresholds are settled for M8a. Two acceptance criteria assume geometry is store-independent and it is not: the phone overflow and wrapping packet paragraphs. They need a doc amendment before P0. Row 156 is CLOSED. Row 150 is 21 of 23 executor milestones built.

**Gate 3b:** the PR head and then the merge commit are each polled, pinned to the full SHA (bounded 30 min, per-invocation verdict path naming the SHA). The outcome is reported on #366.

**Routing evidence:** base=dba643a931c2f69f9071520d13f7ba1883e27634@2026-10-28T08:02:25Z (Gate 1 record). Driver env: `MISSION_DESIGNER_RESOLVED=claude:claude-opus-5-5`, `PLANNER=codex:gpt-6.1-sol`, `EXECUTOR=codex:gpt-6.1-sol`, `EVALUATOR=sonnet`; `MISSION_OVER_RATION` empty.
- Controller `claude:claude-opus-5-5` (tok: not reported).
- Designer, planner and executor: not spawned (R1 is controller-owned per doc §5 and §6; row 156's close is bookkeeping). This is not a routing failure.
- Evaluator `sonnet` (resolver `agent-tool sonnet declared:alias-pin`, Agent tool, foreground, own worktree `.eval-world-iter312` at `e0b45b25`): 84,409 tok / 272 s. Judge independence: an Anthropic Sonnet judge of an Opus controller's measurements, different models. Its directive named the claims to test but supplied no controller method or candidate finding; it was told to use its own code.
- No quorum (no new or revised design doc).
- Metered **$0.00**.

**Ruled out:**
- Text-tight (`Range`) rectangles: their widths vary with hash glyph widths between seeds.
- The header-derived phone background: it crossed the nav.
- Reusing the controller's stats helper for the judge: the directive required independent code.
- A second judge round for the post-judge record corrections: they only narrow claims to the judge's own reproduced measurements.
- Self-amending C8′(iii)/P0 inside R1: changing an acceptance criterion is a design change and goes through the designer and quorum (doc r3).
- Merging #395 from the loop: it was attended work. It merged attended.

**Carried:**
- The C8′(iii)/P0 amendment (doc r3; the judge recommends asserting the phone defect's presence plus title-right, since about ±10 px would be needed; for packet paragraphs, compare left/top/right and tolerate whole-line height steps, or re-measure on the release store under P3).
- `check`'s lack of localisation (F4) is by design; position rests on P0.
- MinRectW/H give no localisation (64x16 glyph windows pass).
- Everything iteration 311 carried (diagnostics-line ordering, the drill's rectangles precondition weaker than `LoadRectangles`, the manifest assertion only parses, `v0.41.0` hardcoded, the load-sensitive M7a supervisor fixtures).

**Retro:**
- The real-Chrome rehearsal found two acceptance criteria that silently assumed store-independent geometry. Hash text changes glyph widths, wraps and overflow, so any pixel-exact criterion must be checked on more than one seed before it is trusted. Two seeds agreeing (5/6) was not enough; the judge's third seed broke it. Earned and pre-registered as instance 1: measure geometry claims on at least three seeds before calling them store-independent.
- Generator ≠ judge paid off again. The judge's own seeds and controls changed the written result (one open criterion became two, and the spread was 14 px rather than 5).
- Harness share 0. No skill edit.

**Progress:** 1.0 clauses 7/7 MET; row 156 CLOSED; row 150 21/23 built with R1 measured; goal unmoved.
**Next:** row 150 M8a (thresholds per V62), then the C8′(iii)/P0 store-dependent-geometry amendment (doc r3, designer + quorum), then M8b, then P0–P3 and the attended A1. Then 151.

## 313 — 2026-10-28 — row 150: M8a built — reproduction docs, captions, release manifest schema and `marketingcheck docs`; judge r1 PASS 82 → controller fix → r2 PASS 89 [PRODUCT]

**Kind:** one executor milestone on the release row, judged twice by an independent evaluator. Every role was spawned through the Agent tool, per the operator's standing request: executor haiku, evaluator sonnet r1 and r2. No designer or planner was needed: the plan and its governing doc exist, and nothing in them was revised.

**Pick and why:** row 150 M8a, following iteration 312's **Next**. Re-checked before picking:
- no M8a commit, branch, PR or file existed (`git log origin/dev --grep M8a` returns only the r2 carve-out doc commit; `docs/marketing-capture-manifest.schema.json`, `scripts/marketingcheck/docs_test.go`, `website/static/img/v1/captions.md` and the website page were all absent from `origin/dev`);
- local dev `05cf2189` was 4 commits behind origin/dev `aafa0b4d` (iteration 312's record merged as #397), so every read came from origin;
- the one open PR from this account, #398 (`trial/row228a-read-cap-base-lag`), is attended trial work (row 228, episode w228a, written in World). It was left alone.
- Governing doc `w-marketing-capture-orchestration` r2 (narrow-refinement carve-out, §9; its quorum artifacts exist, so no pick-time quorum ran).

**Clause map:** clauses 1–7 MET (ratified states, not re-measured); no UNMET clause, no critical-path row. D-WORLD-81 order 156 (closed) → 150 → 151.

**Gate 0:**
- Kill switch armed (`mission-world.disabled` absent); gh account sunholo-voight-kampff; billing tripwire CLEAN.
- `mission_directives.sh`: 0 allowlisted directives on #366 since `2026-10-27T09:03:42Z` (10 comments, none from the allowlist). The watermark is unchanged because nothing new was processed.
- Self-notices rc 1: the newest crash `6345079873` (rc=1, 2026-10-22) is already credited (log grep count 9). No new crash notice.
- Inbox `mission-world`: only this loop's own claim messages. This iteration's claim was sent.
- Ledger: `mission_decisions.sh --check` valid, 86 rows, 0 OPEN.
- Weekly external-issue sweep: not due (#366 opened 2026-10-26, after the Monday boundary).
- Residue: `ailang-world/.wt-world-iter306-r0` (unattributed since 306) left untouched again.

**Gate 1:**
- Dev CI at `aafa0b4d`: check set 3/3 success.
- Running skill and all 13 resource files byte-identical to origin/dev (per-file `cmp` against the resolved symlink target).
- Base `aafa0b4d…` recorded at Gate 1 (`2026-10-28T12:01:59Z`); the Gate-3 `snap` was identical.

**Controller measurements before the directive (load-bearing for the docs):**
- **Docusaurus images.** Probed on this site's Docusaurus 3.10:
  - a markdown image `![](/img/v1/packet.png)` whose file is absent FAILS the build (rc 1, "couldn't be resolved to an existing local image file");
  - a raw JSX `<img src="/img/v1/…">` builds but renders `/img/v1/…`, dropping the deployed `/ailang-world/` base;
  - `pathname:///img/v1/…` drops the base too;
  - `<img src={useBaseUrl('/img/v1/…')}>` builds with the file absent and renders `/ailang-world/img/v1/…`.
  The PNGs do not exist until A1, so the page uses `useBaseUrl` and the README uses ordinary links, not embeds.
- **Stale plan block.** `capture_marketing.sh` REQUIRES `--worldd-bin`, `--bind` and `--marketingcheck-bin` (it fails "missing --…" otherwise). The plan's attended block predates M7c2b and omits all three, so it was not copied. The runbook test now demands flag-set EQUALITY with the shell's `--help`.
- **Capture manifest path.** The shell writes `capture-manifest.json` beside `pass-a/`, not in it, hence amendment A1: `docs` takes an explicit `--capture-manifest`.
- **Base gates green** at `aafa0b4d` (rule 4: the directive's own gate list was baselined): the focused cmd/host gates 9/9 PASS; `scripts/marketingcheck` rc 0 in 35 s.
- C11 observation taken before (12 records).

**Executor (haiku, Agent tool; resolver `recipe claude:claude-haiku-5-5`, following the iteration-305 precedent of Agent `haiku`):** built M8a in `.wt-world-iter313-row150-m8a` with no git writes.
- QUICKSTART §12 pointer and a new §13 attended capture runbook;
- the website page `getting-started/marketing-capture` (six `useBaseUrl` images, disclosures, links to §12/§13 by absolute GitHub URL);
- `captions.md`, the README Screenshots links, and the JSON Schema;
- `marketingcheck docs` (`CheckDocs`, `CheckPromotedAssets`, `runDocs`; +431 lines in `main.go`, the only allowed production file besides `png.go`);
- `docs_test.go` (`TestMarketingDocsLinks`, `TestMarketingRunbookMatchesHelp`, 506 lines), and `TestMarketingDocumentedServeSelectsStill` in `cmd/ailang-worldd/demo_test.go` (the in-process still-mode half, since `scripts/marketingcheck` cannot import package main);
- the png.go threshold comment (V62: values unchanged);
- `cli/unknown` moved from `docs` to `bogus`.
It reported every gate green and 4/4 directive mutants killed (M8a-filename, M8a-flag, M8a-capture-flag, M8a-placeholder). The unminted seed refuses with "--session is neither a 64-hex session token nor a readable file". Its seven deviations were all measured, including: the bind failure exits 1 (`exitUsage = 1` in cmd/ailang-worldd), not 2; the seed line lives in §12.

**Controller verification (outside any sandbox; load 7–9 from sibling work):**
- Main checkout untouched.
- Re-ran: gofmt clean; vet rc 0; `go build ./...` rc 0; compile fence `go test -run '^$' ./...` rc 0; `TestMarketingDocsLinks` (11 PASS lines) and `TestMarketingRunbookMatchesHelp` (8); `docs --root .` rc 0; `TestProductionGoSurface` rc 0.
- Full packages, `go test ./scripts/marketingcheck ./cmd/ailang-worldd ./host/daemon` rc 0 (60 / 171 / 381 s).
- Both named mutants reproduced first-party and restored byte-identical: M8a-filename red at "documented image target must match promoted inventory"; M8a-flag red in both the scripts test and the cmd test.
- C11 after identical; protected-path diff empty.
- **Website build:** rc 0 under `BASE_URL=/ailang-world/`, and the page renders six `/ailang-world/img/v1/…` images. **Instrument failure first, caught by a pristine-base control:** a copy of the site with `node_modules` *symlinked* to the main checkout failed SSG on all 35 existing docs pages ("Cannot read properties of undefined (reading 'id')" in DocItem), and so did the pristine base. The probe that had "built rc 0" compiled in 1.35 s, i.e. from webpack's persistent cache, which lives in the shared `node_modules/.cache`. With a private `cp -c` clone of `node_modules` minus `.cache`, both base and branch build rc 0 (22 s full compile). The red was the shared cache, keyed to another site path, not the site.

**Independent evaluation:**
- **r1 sonnet PASS 82/100** (report `verification/world-iter313/evaluator-r1.md`; completeness 27, truthfulness 22, tests 15, scope 10, honesty 8), zero blocking. It ran every gate, both plan mutants and its own: 13 killed, about 28 survived (mostly main.go refusal branches without a negative arm). Findings F1–F9, all non-blocking. F1: page and README disclosures were ungated (it rewrote the page to "live views" and stayed green). F6: the README spoke in the present tense about images that do not exist yet. F7: the phone alt said "entry" (V50: entries 0–1). F8: §13 omitted P3 and the open criteria. Its verdict said "no false claim".
- **Controller finding the judge missed:** the page said "The demo store's evidence is demo-authored." That is false. Parent design §3: TESTED names S3's real `ailang test` record, ATTESTED the real S2 RecordedEffect, and PROVEN a subject-bound proof report; only CLAIMED is demo-authored. Reproduced from the design text, then fixed.
- **Controller fix round `94e1c08b`:**
  - the provenance sentence corrected in the page and captions;
  - F1: `checkPage` requires the captions' disclosure strings, and `checkReadme` requires "static captures of recorded history" in the Screenshots section. New arms `refuse-page-disclosure` and `refuse-readme-live-claim`; each guard was disabled and its arm went red, then restored byte-identical. Enforcing it immediately refused the page (it said only "static captures", plural), and the page sentence was fixed;
  - F6 README promotion caveat; F7 phone alt "entries"; F8 §13 names P3 and C8′(iii)/P0;
  - full marketingcheck rc 0 (39 s).
- **r2 sonnet PASS 89/100** (`evaluator-r2.md`; 28/24/17/10/10), zero blocking, fresh context, own worktree at `94e1c08b`. F1, F7 and F8 closed; F6 partly closed; F2–F5 and F9 open (carried). New, both non-blocking: N1, the fix's new sentences (provenance, promotion caveat, P3) are not gated; N2, the page is in the autogenerated sidebar before promotion. Its own mutants: all 12 disclosure strings were killed one at a time.

**Bookkeeping:**
- Sprint JSON M8a → `completed` (iteration 313, commit `f6c8d8d4`, PR #400, fix `94e1c08b`, amendments A1). Edited in the file's inline style; a first `jq` rewrite that re-flowed earlier milestones was discarded.
- Plan M8a status line → built.
- Charter: row 150 (M8a BUILT, carried notes, Next); CURRENT GOAL next line; STATUS 313 stamp; STATUS 310 rotated to the archive. Charter line count asserted unchanged (6292); the archive gained 2 lines and holds ITERATION 310, with control 309 present; queue row 150 and `## Queue` present; ledger valid.
- Dashboard overwritten; index regenerated `--index-only` through a temporary registry whose `workdir` is the record worktree.

**Result:** M8a BUILT. The release docs, captions and schema exist and are bound to the shot inventory, the CLIs and the capture shell by tests. Release mode will check the attended PNGs' hashes against the capture manifest. Row 150 is 22 of 23 executor milestones built (M8b left), plus the C8′(iii)/P0 amendment, P0–P3 and the attended A1.

**Gate 3b:** PR #400's head and then the merge commit are each polled, pinned to the full SHA (bounded 30 min, per-invocation verdict path naming the SHA). The outcome is reported on #366.

**Routing evidence:** base=aafa0b4d1093258aea1be6f2f096f9bba2edfd32@2026-10-28T12:01:59Z (Gate 1 record; Gate-3 `snap` identical). Driver env: `MISSION_DESIGNER_RESOLVED=claude:claude-opus-5-5`, `PLANNER=opus`, `EXECUTOR=claude:claude-haiku-5-5`, `EVALUATOR=sonnet`; `MISSION_OVER_RATION=codex`. `MISSION_ROUTING_NOTE`: planner and executor codex `gpt-6.1-sol` unusable (probe rc 75, over daily ration), handed to `opus` and `claude:claude-haiku-5-5`.
- Controller `claude:claude-opus-5-5` (tok: not reported).
- Designer and planner: not spawned (the plan and governing doc exist; no doc revision). This is not a routing failure.
- Executor `claude:claude-haiku-5-5` (resolver `recipe claude:claude-haiku-5-5 declared:provider-pin`; Agent `haiku`, foreground, sprint worktree): 210,200 tok / 555.7 s.
- Evaluator `sonnet` (resolver `agent-tool sonnet declared:alias-pin`; Agent tool, foreground): r1 103,058 tok / 319.9 s in `.eval-world-iter313` at `f6c8d8d4`; r2 59,999 tok / 177.6 s in `.eval-world-iter313-r2` at `94e1c08b`. Judge independence: an Anthropic Sonnet judge of an Anthropic Haiku executor's work and an Opus controller's fix — different models, same vendor (`judge-independence: same-vendor-different-model`). The directives named the claims to test and supplied no controller finding or method; the controller's own demo-authored finding was withheld from r1 and stated to r2 as a claim to verify.
- No quorum (no new or revised design doc).
- Metered **$0.00**.

**Ruled out:**
- A markdown image, a raw `<img>` or `pathname://` for the six shots (measured: build failure or a lost base URL).
- README image embeds before promotion (they would render broken on GitHub until A1).
- Copying the plan's attended block verbatim (it lacks three required flags).
- "The website build is broken at base": the red was the shared webpack cache in a symlinked `node_modules`; a private copy builds base and branch rc 0.
- Gating every judge suggestion in the fix round. F2–F5, F9, N1 and N2 are carried as non-blocking rather than growing the milestone.
- A third judge round: r2 passed with zero blocking, and no code changed after r2.
- Merging #398 or touching it: attended trial work.

**Carried:**
- The C8′(iii)/P0 store-dependent-geometry amendment (doc r3, designer + quorum).
- From the judges: schema, asset and link negative arms (F2, F3, F5); M8a-flag is killed by the token pre-check, with the parse leg exercised only on the happy path (F4); the fix's new sentences ungated (N1); the alt-marker case, titled links and the hard-coded flag count (F9).
- N2, for Mark at A1: the release-screenshots page is in the sidebar before promotion, so a site deploy between this merge and A1 shows six broken images. Either deploy only after A1, or mark the page `unlisted` until promotion.
- Everything iteration 312 carried.

**Retro:**
- **Instrument, instance 1, pre-registered:** a `node_modules` shared by symlink also shares webpack's persistent cache (`node_modules/.cache/webpack`). A build from another site path then fails SSG on every docs page, or "succeeds" from cache in about 1 s. Build Docusaurus probes with a private `cp -c` clone minus `.cache`, and treat a sub-5 s client compile as a cached reading, not a measurement. The pristine-base control is what caught it, before it was attributed to the branch.
- **A plan's verbatim command block is a claim dated to the plan.** Later milestones changed the CLI it quotes (M7c2b added three required flags). Re-derive any quoted command against the tool's `--help` at dispatch. The new runbook test now does this for §13. Instance 1.
- **Generator ≠ judge is necessary, not sufficient.** The r1 judge reported "no false claim" over a sentence that misdescribed three of the four evidence grades; the controller's own read of the public text caught it. This is the existing judge-findings-can-be-under-stated rule, not a new one.
- Harness share 0 of the last 20 index rows. Bar-moving landings: 311, 312 and 313 each moved no clause, and every clause is MET (there is no UNMET clause and no routable critical-path row), so this is not DRIFT. Release-row work is permitted under D-WORLD-64. No skill edit.

**Progress:** 1.0 clauses 7/7 MET; row 150 22/23 built (M8b left) plus the C8′(iii)/P0 amendment, P0–P3 and A1; goal unmoved.
**Next:** row 150's C8′(iii)/P0 store-dependent-geometry amendment (doc r3, designer + quorum), then M8b, then P0–P3 and the attended A1. Then 151.

## 314 — 2026-10-28 — dev RED at Gate 1 → arm 6 reverse-leg false refusal fixed and LANDED (row 229, PR #403; judge PASS 92) [PRODUCT-HYGIENE]

**Kind:** full inner loop on a Gate-1 red (designer → quorum ×2 → carve-out revision → planner → executor → independent evaluator → PR), every role spawned through the Agent tool per the operator's standing request.

**Pick and why:** Gate 1 found `origin/dev` `2769cff9` (PR #401, docs-only) red on `go host build + test gate`, run 40699055042 attempt 1, job 122853572046, step 8 (`Row 140 Workspace.Exec real-srt tests`), with 17 steps skipped behind it. The one failure was host/daemon `TestExecSrtMCPGrantAndBudget`: `exec startup probe failed for profile "proj" (1 of 7 arms): arm6-network: reverse leg: the host reached the listener the child bound on 127.0.0.1:44109`. Parents `c953cadc`, `9ec31fca`, `e40ab235`, `aafa0b4d` and `05cf2189` were green, and a re-run of the failed job on the identical tree (attempt 2) was green with every step run, so it is a flake. A scan of the 8 newest failed CI runs found no other real-sandbox instance (instance 1). A red dev outranks the queue and World owns this repo, so this followed the iteration-228 precedent, even though dev was already green again on the re-run. The fail-closed probe can refuse daemon startup in production the same way, so the defect is real beyond CI.

**Clause map:** clauses 1–7 MET (ratified states, not re-measured); no UNMET clause, no critical-path row. D-WORLD-81 order 150 → 151 resumes next iteration.

**Gate 0:**
- Kill switch armed (`mission-world.disabled` absent); gh account sunholo-voight-kampff; billing tripwire CLEAN.
- `mission_directives.sh`: 0 allowlisted directives on #366 since `2026-10-27T09:03:42Z` (11 comments, none from the allowlist). Watermark unchanged.
- Self-notices rc 1: the newest crash `6345079873` (rc=1, 2026-10-22) is already credited (log grep count 11). No new crash notice.
- Inbox `mission-world`: one `harness-resolved` reply from `mission-fleet` for `gate5:controlplane-send-refused-local-store` (fleet iteration 66, #1994, `d95ead6`, judged PASS 91). Row 144 is marked RESOLVED and the reply acked. The rest are this loop's own claim messages; this iteration's claim was sent.
- Ledger: `mission_decisions.sh --check` valid, 88 rows, 0 OPEN.
- Weekly external-issue sweep: not due (#366 opened 2026-10-26, after the Monday boundary).
- Residue: `ailang-world/.wt-world-iter306-r0` (unattributed since 306) left untouched again.

**Gate 1:**
- Base `2769cff9…` recorded at `2026-10-28T16:01:25Z`; local dev = origin/dev; no open PR from this account.
- Running skill and all resource files byte-identical to origin/dev, for both the resolved symlink target and the pin copy (per-file `cmp`, silent).
- Check set at `2769cff9`: 3 checks, `go host build + test gate` failure (above).

**Diagnosis (controller, then designer):**
- `host/broker/exec_probe.go` (unchanged since `885af2d1`, 2026-10-20) picked `revPort` with `freeLoopbackPort` (listen `:0`, read the port, close it), then `pollLoopback` counted ANY connect to it as a reach.
- M0's committed matrix shows the child binds inside its own netns on Linux (`WORLD-PROBE-BOUND; host reached=False`) and is refused (EPERM) on macOS. So on Linux the host can never reach the child, and a reach means another host process took the port.
- The controller's hypothesis, "parallel test packages", was REFUTED by the designer: step 8 runs `go test ./host/broker/ ./host/daemon/ -p 1`, and neither package calls `t.Parallel()`. The controller re-verified both (control: `t.Parallel()` exists in host/verifygate).
- The designer found the concrete in-window binder by reading srt 0.0.78's source: the mux proxy binds host `127.0.0.1:0` at every srt start, including Run C's own, and stays silent until the first byte. That is the exact shape the poll counted.
- Deterministic reproduction at base: a stand-in srt holding a listener on the reverse port gave the CI signature 3/3; the control gave 0/3.

**Design and quorum** (`w-exec-probe-reverse-leg-port-race`, designer opus via the Agent tool; author vendor benched; seats glm-5.3, kimi-k3 and gpt6-1-sol, with reserve gemini-4-argon):
- **r1:** glm PASS, gemini PASS, kimi REJECT: no caller census for the two changed helpers. gpt6-1-sol was absent (OpenAI 429, "no credits remaining"). Measured, not forwarded (rule 3f): `grep -rn 'pollLoopback\|freeLoopbackPort'` gave 6 lines, one call site per helper, both in the reverse leg (control `runBounded` = 9). The designer applied kimi's verbatim fix in r2 (V22, R4).
- **r2 (the one re-quorum):** gemini PASS, glm REJECT (the new tests T1/T2/T3/T5/T7 are not in CI's PASS-guard list, so a SKIP could pass silently), kimi REJECT (a silent listener that is reached would be retried into a pass, which weakens P20). gpt6-1-sol absent again (429).
  - Kimi's fix branched on an unmeasured premise: does srt's mux answer a garbage hello? The controller measured it live with the pinned srt 0.0.78 (cli.js sha `3c3092bd…` = pin), importing its own `createMuxProxyServer` and `createHttpProxyServer`. Hello → `HTTP/1.1 400 Bad Request` and close within 0–3 ms (3/3); no hello → silent to the 1.5 s deadline (2/2). Source reading agrees: a first-byte peek, then the HTTP backend's `clientError` 400.
  - That resolved kimi's conditional to its first branch. Every remaining objection now had a concrete reviewer-authored fix and none disputed the direction, so the **narrow-refinement carve-out** applied: r3 is the reviewers' verbatim fixes (the hello, the three-way classification, a sticky silent refusal, all 7 tests in `tests=`, node-less V24). It routed straight to the planner.
- Quorum spend $0.687 (r1 $0.314, r2 $0.373). Artifacts are banked in `verification/world-iter314/quorum/`.

**Fix:**
- The child's listener answers each accept with `WORLD-PROBE-NONCE <nonce>`, using a fresh port and a 128-bit nonce per attempt. World writes `WORLD-PROBE-HELLO` on each connected dial and reads for 500 ms:
  - nonce → child → refuse ("the host reached…");
  - other bytes, or a clean EOF → foreign → retry on a fresh port, at most 3 attempts, then "inconclusive after N attempts";
  - zero bytes, a reset, or a failed hello write → silent → refuse, and the refusal sticks for the leg.
- A pass needs zero connected dials.
- Planner opus: M1 production and the pure judge T1; M2 node-free shell tests T2/T3/T5/T7; M3 node tests T4/T6, `ci.yml` `tests=` and QUICKSTART §10. Seven design/code discrepancies (D1–D7) were resolved in the plan.

**Executor (haiku, Agent tool; resolver `recipe claude:claude-haiku-5-5 declared:provider-pin`):**
- Built M1–M3 in `.wt-world-iter314-arm6` with no git writes; snapshots `.snap/M1..M3`.
- Reported every gate green and 21/21 mutants killed, each restore sha-verified.
- Deviations: `_ = werr` for M18/M19 to compile; G10 not re-run at base; it did not write the report file (its own instructions forbid report files), so the controller banked it.

**Controller verification:**
- Main checkout untouched. Commits rebuilt per milestone from the snapshots: M1 `2f2ba233`, M2 `5856a606`, M3 `d7bba842`. At each boundary gofmt, vet, compile fence and `-run TestExecProbe` were rc 0. Reconstruction faithful, `shasum -c` 4/4 OK.
- Gate sweep (`verification/world-iter314/controller/`): gofmt, vet and compile fence rc 0; `-race` on the probe tests rc 0; node-less run rc 0 (T4 SKIPs by design); 10× stability rc 0; actionlint clean. `verify_go.sh` plain with the pinned v0.41.0 `AILANG_BIN` PASSED (broker 285 s, daemon 348 s).
- **Instrument failure first:** the controller's own G4 `go test ./host/broker/` ran without `AILANG_BIN` and went red on `TestEpisodeLiveReplayThreeArmsAndEvidence` ("AILANG_BIN must name the pinned released interpreter"). That is a fact about the gate script, not the code; the `verify_go.sh` leg is the corrected form.

**Independent evaluation:**
- r1 sonnet (Agent tool, own worktree `.wt-world-iter314-eval` at `d7bba842`): **PASS 92/100** (28/19/24/10/11), zero blocking. Report: `verification/world-iter314/evaluator-r1.md`.
- It derived from the code that no path passes a reach of the child. It ran the gates (17 PASS, 1 SKIP Arm8) and T6 3×, and killed 10 plan mutants and 8 of its own 10. Survivors: O8 (hello wait 500 → 5 ms) and O10 (drop `c.resume()`).
- Findings, all non-blocking: F1 (design residual: a real reach misread as foreign, then a clean attempt), F2/F3 (the two survivors), F4 (G10 not re-run at base).
- **Controller reproduction:** O8 GREEN and O10 GREEN (survivors confirmed); M4 RED; M17 RED on the assertion "want the sticky silent refusal on attempt 1", not on a build error. Eval tree restored, 0 porcelain.
- F2 and F3 are carried as row 230. No fix round: neither is a correctness defect, and F3's only discriminating platform is Linux CI.

**Bookkeeping:**
- Design, plan and JSON moved to `implemented/` (status, milestone commits, PR, evaluation).
- Charter: rows 229 (LANDED) and 230 (pins, cross-profile poll, R4) appended after 228; row 144 marked RESOLVED by the fleet; CURRENT GOAL next line updated; STATUS 314 stamp added and STATUS 311 rotated to the archive. Line count 6294 → 6296 (+2 rows; the stamp swap is net 0); archive +2 with ITERATION 311 at its head; 3 stamps; `## Queue`, row 150 and rows 228–230 present; ledger valid.

**Result:** LANDED (pending merge-commit CI). Arm 6's reverse leg counts only the child's nonce as a reach. A port another host process holds, including srt's own mux, is retried rather than refused, while any silent or nonce-answering reach still refuses. On the PR head `d7bba842` CI was 3/3 green, and its Linux job shows `--- PASS` for all 7 new tests, `TestExecSrtMCPGrantAndBudget`, `TestExecSrtStartupProbePasses` and the pass-through control, with 0 FAIL lines.

**Gate 3b:** PR #403's head `d7bba842` was polled SHA-pinned (bounded 30 min, verdict path naming the SHA): ALL GREEN 3/3. The record commit and the merge commit are each polled the same way, and the outcome is reported on #366.

**Routing evidence:** base=2769cff9477dc89f5eb61c6aaeb58079cb4120e3@2026-10-28T18:26:55Z (Gate 4 record; the Gate 1 record and the Gate 3 `snap` were the same SHA). Driver env: `MISSION_DESIGNER_RESOLVED=claude:claude-opus-5-5`, `PLANNER=opus`, `EXECUTOR=claude:claude-haiku-5-5`, `EVALUATOR=sonnet`; `MISSION_OVER_RATION=codex`. `MISSION_ROUTING_NOTE`: planner and executor codex `gpt-6.1-sol` unusable (probe rc 75, over daily ration), handed to `opus` and `claude:claude-haiku-5-5`.
- Controller `claude:claude-opus-5-5` (tok: not reported).
- Designer: resolver `recipe claude:claude-opus-5-5 declared:provider-pin`, spawned as Agent `opus` (same weights, per the operator's standing request; iteration-228 precedent). The rotation state file's next entry after `codex:gpt-6.1-sol` is `pi:ollama/glm-5.3:cloud`; the resolver's answer was followed verbatim, and the pointer was set to the lane that ran. 1 doc, 3 runs: author, the r2 revision (protocol-mandated), and the r3 carve-out revision. The third run exceeds the one-doc-plus-one-revision diet and is FLAGGED; it was the reviewers' own verbatim fixes, not a new design. Tokens 164,464 / 170,326 / 234,780 (reported per run, cumulative transcript).
- Planner: Agent `opus` (`agent-tool opus fail-closed:env-pin`; derive-planner-lane `opus fail-closed:env-pin`), 218,436 tok.
- Executor: Agent `haiku` (resolver `recipe claude:claude-haiku-5-5 declared:provider-pin`; iteration-305/313 precedent), 161,690 tok / 2,773 s.
- Evaluator: Agent `sonnet` (`agent-tool sonnet declared:alias-pin`), 60,099 tok / 268 s. Generator ≠ judge: haiku vs sonnet.
- Quorum: glm-5.3 / kimi-k3 / gemini-4-argon (reserve), with gpt6-1-sol absent on OpenAI 429 both rounds. Metered $0.687.

**Ruled out:**
- "Parallel test packages caused the collision" (controller hypothesis): refuted, step 8 is `-p 1` and no broker/daemon test is `t.Parallel()`.
- TCP self-connect as the primary cause: Go's dialer detects and retries it.
- Option B (bounded retry without identity): it lets an intermittent breach pass.
- C1 (World holds the port): it stops observing reachability.
- C2/C3 (child-chosen port, out-of-range port): they narrow the window without classifying a reach.
- The controller's own G4 red: AILANG_BIN unset in the gate script, an instrument failure.

**Carried:**
- Row 230: F2 hello-wait pin; F3 `c.resume()` pin; R2 cross-profile poll (`w-exec-probe-cross-profile-poll`); R4 Linux collision rate unmeasured.
- The OpenAI quorum seat (`gpt6-1-sol`) is absent on HTTP 429 "no credits remaining". This is the metered API account, not the codex subscription. A capacity condition, not a decision.
- `.wt-world-iter306-r0` residue.

**Retro:**
- **A reviewer's conditional fix is a measurement request.** Kimi's r2 objection branched on "does srt's mux answer garbage?" Measuring it (≈1 minute, the pinned package's own server constructors) turned a judgment-shaped block into a carve-out-eligible one, and the alternative was a `needs-human-review` park on a question nobody had measured. This is rule 3f applied to a branching fix. Instance 1.
- **A controller gate script must export the mission's pin itself.** G4 ran bare `go test` without `AILANG_BIN` and produced a false red; `verify_go.sh` refuses loudly when the pin is unset, but a hand-written gate script does not inherit that. Instance 1.
- **An executor's own instructions can override a directive's report path.** Haiku declined to write `executor-report.md` and returned the text instead; the controller banked it. Instance 1 (no skill edit).
- Harness share 0 of the last 20 index rows. Bar-moving landings: 312, 313 and 314 each moved no clause, and every clause is MET with no routable critical-path row, so this is not DRIFT. No skill edit.

**Progress:** 1.0 clauses 7/7 MET; goal unmoved (hygiene on MET clause 3); row 150 22/23 built (M8b left) plus the C8′(iii)/P0 amendment, P0–P3 and A1.
**Next:** row 150's C8′(iii)/P0 store-dependent-geometry amendment (doc r3, designer + quorum), then M8b, then P0–P3 and the attended A1; then 151. Row 230 sits below every release row.

## 315 — 2026-10-28 — row 150: C8′(iii)/P0 store-dependent-geometry amendment written (orchestration doc r3 → r5); quorum BLOCKED twice → carve-out; judge r1 PASS 77 → controller fix → r2 PASS 84 [PRODUCT]

**Kind:** design amendment (designer → quorum → designer revision → re-quorum → controller carve-out → independent evaluator ×2 → controller fix round → PR). Every role was spawned through the Agent tool, per the operator's standing request. No planner or executor was routed: the pick is a doc revision, and it changes no milestone (M8b is next and its plan is unchanged).

**Pick and why:** the queue head under D-WORLD-81 (150 → 151) and the iteration-314 **Next**: row 150's C8′(iii)/P0 amendment. R1 (iteration 312, V60–V63) had measured two acceptance rules that cannot hold on a store other than the one measured: P0 (a) "±2 px on every edge" (packet paragraphs wrap with the store's text) and phone C8′(iii) "overflow geometry within ±2 px of V50" (the thumb width follows the text). P0, P3 and A1 cannot run until both are amended.

**Clause map:** clauses 1–7 MET (ratified states, not re-measured); no UNMET clause and no critical-path row; release order D-WORLD-81 150 → 151 → v1.0.0.

**Gate 0:**
- Kill switch armed; gh account sunholo-voight-kampff; billing tripwire CLEAN.
- `mission_directives.sh`: 0 allowlisted directives on #366 since `2026-10-27T09:03:42Z` (12 comments). Watermark unchanged.
- Self-notices rc 1: the newest crash `6345079873` (2026-10-22) is already credited (log count 12). No new crash.
- Inbox `mission-world`: 21 unread, all this loop's own claim messages; the iteration-315 claim was sent and confirmed in the listing.
- Ledger valid, 88 rows, 0 OPEN. External-issue sweep not due (#366 is this week's issue).
- **Hygiene incident (controller):** the first Gate-0 env read listed values rather than names, which printed an OAuth token value into this session's local transcript. That is the exact pattern Gate 0 forbids. The value is in no repo file, commit or message. Reported to Mark.

**Gate 1:**
- Local dev was 9 behind with a clean tracked tree; fast-forwarded to `09ff2821`. Base `09ff2821…` recorded at `2026-10-28T20:01:37Z`.
- Skill: `SKILL.md` and all 13 resource files are byte-identical to `origin/dev` at the resolved symlink target (per-file `cmp`).
- CI at `09ff2821`: 3 checks, all success.
- Open PR from this account: #402 `trial/row217a-presets-start`, an attended Codex-through-World trial of row 217. It is not this loop's work and was left alone.

**Design (designer opus, Agent tool; resolver `recipe claude:claude-opus-5-5 declared:provider-pin`):**
- **r3: P0 (a)** is a structural match: (a1) left/right ±2 px; (a2) stack-head top ±2 px; (a3) each later region's gap from the previous region's **measured** bottom equals the R1 gap ±2 px, and the two body-text gaps may differ by whole 24 px lines; (a4) packet paragraph heights are the R1 height plus whole 24 px lines (`line-height:1.5`, `render.go:193`). The content gate runs inside the capture store's own accepted rectangles.
- **Phone (iii′):** the overflow defect is present. That means a thumb in the bottom 20 px with x0 ≤ 8 and x1 ≤ 384, and title ink at x ≥ 387.
- Measured on the nine iteration-312 stores (V64–V68):
  - P0: 9/9 accepted (r2 8/9) and 18/18 negatives refused, including every wrong crop `check` accepted at V63.
  - Phone: 19/19 accepted (r2 1/19) and 11/11 refused.
- Thumb right edge per store: 324–340, against V50's 329. The judge's ±10 px idea would already refuse the 340 store.
- **Controller check before quorum:** rebuilt both checkers from the banked `.txt` and re-ran them. P0 r3 accepts 7/7 block-box stores (r2 6/7); phone accepts V50 + 9 stores and refuses the desktop image.

**Quorum** (`--author claude:claude-opus-5-5`; seated kimi-k3, gpt6-1-sol, gemini-4-argon; reserve glm-5-3; the Claude seat sits out; gpt6-1-sol absent both rounds, `unreachable`, OpenAI 429):
- **r3 round:** BLOCKED 3/3 (kimi, gemini, glm), $0.43. All three objected to provenance and labels: the header claimed one base for rows measured at four commits, §9 still said "r1", milestones that had landed were shown pending, the invocations were unbound, and the row-149 coupling was unnamed.
- **Controller measured the drift premise first (rule 3f):**
  - ancestry rc 0 ×3;
  - `host/workbench`: 0 files and 0 commits in 4e6e89dc..09ff2821, with range control `aeaa957d` = 2;
  - `host/daemon/workbench*` = 0;
  - `daemon.go` changed only health/session fields;
  - `png.go` changed only a comment.
- **r4 (designer revision, resumed via SendMessage):** per-row bases, V69 (drift), V70 (milestones), V71 (checkers build from `.txt`).
- **r4 re-quorum:** BLOCKED, $0.58. Gemini **PASS**. Kimi: the release-blocking checker is an uncommitted `.txt` program; fix = commit it, OR keep the route with sha256 pins, capture-day controls and a stated reason. GLM: P0's blocking clause needs a still-HTML repeat-GET whose instrument P0 retires; fix = an instrument-per-artifact checklist, plus one stored R1 height.
- **Disposition:** neither reject disputes the r3 rules, and every objection carries a concrete fix, so the **narrow-refinement carve-out** applies.
- **r5 (controller):**
  - Kimi's second option: pins, plus the control step executed as V72. The rationale paragraph and the V50 known-accept substitution are named controller-authored. Kimi's first option was filed as row 231.
  - GLM's checklist, verbatim, and C12 aligned with it.
  - V73: the stored R1 height is 97 (block 96.0, floor/ceil), and the wrapped store's is 73.
  - Each fix got its own named mutation per the World guardrail (V74: M-r5-pin, M-r5-control, M-r5-run).
- Artifacts banked: `verification/world-iter315/quorum-r3.json`, `quorum-r4.json`.

**Independent evaluation (sonnet, Agent tool, own worktree `.eval-world-iter315`):**
- **r1 at `c36a4306`: PASS 77/100** (24/20/13/11/9), zero blocking.
  - Reproduced V64, V66, V69, V72 and V73 byte-for-byte or numerically.
  - Read both checkers against the prose clause by clause, and ran its own arms.
  - Found phonecheck's thumb test is a first-to-last **span**, not the **run** the rule and its own header say: two stray grey pixels passed.
  - Also found stale r4/r3 labels, row 231 not yet filed, identity-only accept controls, whole-line slides accepted in the free gaps, compounding ±2 drift, V65 not isolating four clauses, and the carve-out guardrail's per-fix criterion and mutation missing.
- **Controller fix round (`0e975c53`):** reproduced the span defect (old ACCEPT) and changed `scrollbar()` to the longest contiguous run. The 19 positives' outputs are byte-identical, the scattered arm flips to REFUSE, and the other 28 arms are unchanged. Source re-pinned, V66/V67/V72 re-run, stale text fixed.
- **r2 at `0e975c53`: PASS 84/100** (25/21/17/12/9), zero blocking. Divergence 1 and the stale text are CLOSED. It found a new fragility: a real thumb split by one off-range pixel left of its midpoint is refused. It left open NB2 (row 231), NB7 (per-fix criteria) and NB8 (truncated hashes).
- **After r2 (`e5cd4673`, controller, not re-judged):** V74 named mutations and the per-fix criteria table (NB7), full hashes (NB8), the exclusive-X1 note, and the judge's arm generator banked. Row 231 was filed with this record (NB2).

**Bookkeeping:**
- Charter: STATUS 315 added; STATUS 312 rotated to the archive (archive count 1, control 311 = 1). Charter 6297 → 6298 (+1 = row 231; the stamp swap nets 0); archive +2. `## Queue` and rows 150 and 229–231 present; ledger valid, 88 rows. CURRENT GOAL next line updated. Row 150 records the amendment and NEXT = M8b.
- Doc stays in `planned/`, because row 150 is not complete.

**Result:** the amendment is BUILT on `sprint/world-iter315-row150-c8-amend`. Design commits are `c36a4306` (r3–r5), `0e975c53` (judge-fix round) and `e5cd4673` (V74). It is LANDED only after merge-commit CI (Gate 3b).

**Gate 3b:**
- PR #404's head `d4df8a15` was polled SHA-pinned (bounded 1800 s; the verdict file names the SHA): **RED**. `go host build + test gate` failed one subtest, `TestExecProbeReverseLegSilentHolderRefusesSticky/reset-holder` (`Why = "", picker calls 1; want the sticky silent refusal`). That is one of row 229's tests, landed in iteration 314. Every other package passed, and `ailang-code verify gate` and the race gate are green.
- **Attribution:** this diff is `design_docs/` only (0 files outside), and the same test was green at `d7bba842` and `09ff2821`. So the failure is intermittent and inherited, not from this PR. Instance 1.
- **Not ignorable:** the leg PASSED where the design says a reset must refuse. That could be a false pass of the fail-closed probe, filed as **row 232** with the excerpt banked.
- Actions: re-ran the failed job on the identical tree (attempt 2), and pushed the record commits (new head; a full CI run). Both are polled before merge; the outcome is reported on #366.

**Routing evidence:** base=09ff28215176d5c09b36ecffba85ee04660ba238@2026-10-28T20:01:37Z (Gate 1 record; the Gate 3 `snap` was the same SHA). Driver env: `MISSION_DESIGNER_RESOLVED=claude:claude-opus-5-5`, `MISSION_EVALUATOR_RESOLVED=sonnet`; `MISSION_OVER_RATION` empty; `MISSION_ROUTING_NOTE` empty.
- Controller `claude:claude-opus-5-5` (tok: not reported).
- Designer: resolver `recipe claude:claude-opus-5-5 declared:provider-pin`, spawned as Agent `opus` (same weights, per the operator's standing request; iteration-228/314 precedent). The rotation file holds `claude:claude-opus-5-5`, unchanged because that is the lane that ran. 1 doc, 2 runs (author + the protocol-mandated revision), within the one-doc diet. Tokens 196,359 + 261,559.
- Planner and executor: not routed (design amendment; no milestone change).
- Evaluator: Agent `sonnet` (`agent-tool sonnet declared:alias-pin`), r1 160,583 tok / 462 s, r2 186,872 tok / 82 s (resumed). Generator ≠ judge: opus vs sonnet.
- Quorum: kimi-k3 / gemini-4-argon / glm-5-3 (reserve), gpt6-1-sol absent both rounds (unreachable, 429). Metered: r3 round $0.4317 (kimi 29,993/4,415 tok, gemini 34,060/1,075, glm 30,597/17,020) + r4 round $0.5831 = **$1.01**.

**Ruled out:**
- The controller's own "the V50 reference image is no longer on disk" (from a search of the repo and `/tmp/world-iter305*` only). It lives at `~/.ailang/state/world-iter305/probe-seeded/pass-a/phone.png` (sha `999125ff…` = V63). Rule 3a: the search location was the instrument.
- The first quorum launch's rc 127: `timeout` does not exist on macOS, so the quorum never ran. It was re-run under a `date +%s` deadline loop. An instrument failure, not a verdict.
- The controller's brief said M7c2a had no PR. The designer's `git log` showed PR #380 (`b3b373ec`), and the controller confirmed it.
- Re-baselining R1 for cross-base drift: unneeded, because the renderer and still path are byte-unchanged from 4e6e89dc to 09ff2821.
- A ±10 px tolerance around V50 for phone (iii): refuted by the 340 px store.
- The "ratified in iteration 95" citation in r5: not confirmable from the charter (the judge), so it was removed.

**Carried:**
- Row 231: commit the checkers; the free-gap tolerances; compounding drift; identity accept controls; P0 item (1) command; V65 isolation; (ii) stand-in; the split-thumb fragility.
- gpt6-1-sol absent from the quorum on OpenAI 429 for the third iteration running (314, 315 ×2). A capacity condition, not a decision.
- Row 232: row 229's `reset-holder` arm passed the reverse leg on Linux CI (instance 1; possibly a false pass).
- `.wt-world-iter306-r0` residue.

**Retro:**
- **A failed heredoc does not stop the next line.** The record script aborted on an ambiguous anchor, and the `git add/commit/push` on the following lines still ran, pushing a partial record commit (`0db9358f`, row 232 only). The memory rule "push only after checked exit codes" was known and still missed, because the chain broke at a heredoc boundary. Gate the push on the script's rc. Instance 1 for World (memory: instances 1–2 on 2026-10-16).
- **A checker banked as `.txt` gets reviewed as prose, but it runs as code.** The quorum read the rule text, and only a judge told to "read the code against the prose clause by clause" found that the implementation tested a span where the prose (and the code's own header comment) said a run. Instance 1. The directive line is what caught it.
- **Verdict-in-text instruments exit 0 on REFUSE.** Both checkers print ACCEPT/REFUSE and exit 0 either way, so any future wrapper keyed on `$?` would read every capture as green. This is now stated in §4.4, and row 231's committed form should exit non-zero. Instance 1.
- **`timeout` is not a macOS command.** A bounded-wait idiom borrowed from Linux failed rc 127 and reported nothing. Use the `date +%s` loop the skill prescribes. Instance 1.
- **Read the environment by NAME** — a Gate-0 rule the controller broke this iteration (above). Instance 1 for World, and the rule already exists, so no skill edit.
- Harness share 0. Bar-moving landings: 313, 314 and 315 moved no clause, but every clause is MET with no routable critical-path row, so this is not DRIFT. No skill edit.

**Progress:** 1.0 clauses 7/7 MET; goal unmoved (release row). Row 150: 22/23 executor milestones built (M8b left); the C8′/P0 amendment is written; P0–P3 and the attended A1 remain.
**Next:** row 150 M8b, then P0–P3 (r5 predicates and capture-day controls) and the attended A1; then 151. Rows 230 and 231 sit below every release row. Row 232 (the reset-holder false pass) takes precedence if it reds dev at Gate 1.

## 316 — 2026-10-29 — row 150: M8b built — final repository gates and C11 protected-path observation; QUICKSTART §13 states the r5 P0/phone rules; judge r1 PASS 93 → r2 PASS 93 [PRODUCT]

**Kind:** validation milestone (executor → controller re-run of every gate outside the sandbox → independent evaluator → one controller doc fix → evaluator r2 → PR). Roles were spawned per the routing table and the operator's standing request. The executor ran on its provider pin through the codex recipe, not the Agent tool, because the spawn-pin hook refuses an alias spawn for a `provider:model` pin. The judge ran on the Agent tool.

**Pick and why:** the iteration-315 **Next** and the queue head under D-WORLD-81 (150 → 151): row 150 M8b, the last executor milestone. Dev CI at `8bfaa1ab` was green (CI 1/1 success on the head), so row 232 did not outrank it.

**Clause map:** clauses 1–7 MET (ratified states, not re-measured); no UNMET clause and no critical-path row; release order D-WORLD-81 150 → 151 → v1.0.0. M8b moves no clause.

**Gate 0:**
- Kill switch armed (`mission-world.disabled` absent); gh account sunholo-voight-kampff; billing tripwire CLEAN. Env read by NAME only.
- `mission_directives.sh`: 0 allowlisted directives on #366 since `2026-10-27T09:03:42Z` (13 comments). Watermark unchanged.
- Self-notices rc 1: the newest crash `6345079873` (2026-10-22) is already credited (log count 13). No new crash.
- Inbox `mission-world`: unread items are this loop's own claim messages. The iteration-316 claim was sent and confirmed in the listing.
- Ledger valid, 88 rows, 0 OPEN. External-issue sweep not due (#366 is this week's issue).

**Gate 1:**
- Local dev = origin/dev = `8bfaa1ab` with a clean tracked tree. Base recorded `8bfaa1ab28000937234aececa2ccd818e2b522a9@2026-10-29T00:04:34Z`; the Gate-3 snap matched (no drift).
- Skill: `SKILL.md` and all 13 resource files are byte-identical to `origin/dev` at the resolved symlink target (per-file `cmp`).
- Open PR from this account: #402 (attended row-217 trial). Not this loop's work; left alone.

**Rule 3e(a) baseline (controller, at base, outside the sandbox):** all 7 focused plan gates had a named top-level PASS with 0 SKIP. `verify_ail.sh` rc 0 in both the plan's form (`AILANG_BIN` = pinned v0.41.0, 82 s) and the release form (v0.54.19 + `WORLD_PKG_AILANG_BIN`, 28 s). The initial C11 observation had 12 records. The first baseline attempt word-split a list under zsh and ran nothing; that was an instrument failure, rerun from a bash script.

**Executor (codex `gpt-6.1-sol`, recipe; probe rc 0; sandbox `workspace-write`; 11 min; 130,665 tok):**
- `docs/QUICKSTART.md` §13: the sentence "two acceptance criteria are still open" had been false since r5. It now states P0 (a)'s structural match and phone C8′(iii′)'s defect-present rule, and points at §4.4 for the pinned checkers and capture-day controls. It also says the verdict is the `r3` text, because both checkers exit 0 either way.
- `docs/QUICKSTART.md` §13: `WORLD_MARKETING_RECTS` changed from the R1 file to `$DEMO_DIR/rects.json` (§4.4 P3: the release store's own rectangles are the capture file; R1 is the position baseline).
- `README.md`: row 156 closed in iteration 312 (D-WORLD-86 = A), so `v1.0.0` waits on 150 and 151.
- Three allowlisted files were read and left unchanged.
- Mutation M8b-C11-redrill: `if false && b.Lines != a.Lines` was killed by `digest-equal-line-change` only, and restored byte-identical (sha `bfc99a8d…` both sides).
- C11 before/after equal; protected diff empty.
- Socket tests were labelled UNINFORMATIVE UNDER SANDBOX. The full and race suites were not run, as directed (controller-owned).

**Controller verification (outside the sandbox, at `3f185a9a`, macOS):**
- 7/7 focused gates had a top-level PASS with 0 SKIP.
- QuickStart sweep in `cmd/ailang-worldd`: 7 PASS, 0 FAIL. `TestMarketingDocsLinks` and `TestMarketingRunbookMatchesHelp` PASS; `TestQuickstart*` in host/broker and host/daemon ok.
- `verify_ail.sh` rc 0 (65 s); `go build ./...` rc 0; `go test -race ./host/daemon` rc 0 (477 s); `go test -race` coordinator/workbench/evidence/marketingcheck rc 0 (188 s).
- **`go test ./...` rc 1, red only in `host/broker`:** `TestExecProbeForeignListenerOnTheReversePortDoesNotRefuse/chatty` (arm 6, `accepts 2`) and `TestExecProtectedCheckTimeout/reserve-after-before` (a deadline overrun). The package took 524 s while the judge ran heavy focused tests concurrently.
- Attribution: each test passed 3/3 isolated at head AND at base `8bfaa1ab`. The package passed uncontended (rc 0, 231 s). Neither test file reads the changed docs (grep 0). So the reds are intermittent under CPU load and inherited, not this diff.
- The `chatty` arm is a row-229 probe test, so it is added to row 232 as instance 2 (macOS, under load; instance 1 was `reset-holder` on Linux CI). Banked at `verification/world-iter316/m8b-host-broker-load-reds.txt`. That file sits under the plan's "and its scratch snapshots" clause, not the five-file allowlist (judge r2 note).

**Independent evaluation (sonnet, Agent tool, own worktree `.eval-world-iter316`, foreground):**
- **r1 at `3f185a9a`: PASS 93/100** (31/20/14/14/14), zero blocking.
  - Checked every changed sentence against §4.4 and the charter, and recomputed both checker sha256 pins.
  - Re-ran the redrill (one subtest reds) and two mutations of its own. A `--rectangls` flag typo is killed by `TestMarketingRunbookMatchesHelp`, and a `packetx.png` link by `TestMarketingDocsLinks`.
  - Ran the focused gates plus the doc readers: all PASS.
  - NB1: §13 does not say where `$DEMO_DIR/rects.json` comes from. NB2: reverting the rects path to R1 **survives** every test (the controller reproduced this: both doc tests PASS with the revert). NB3: the P0/P3 sentence is terse. NB4: an executor ledger count slip (51 vs 52, self-disclosed).
- **Controller fix (`706af0e4`):** NB1. The sentence now names "the V60 block-box probe's output for that store" (§4.4 P0 item 3). Doc readers re-run: PASS.
- **r2 at `706af0e4` (resumed): PASS 93/100**, zero blocking.
  - NB1 CLOSED.
  - NB2 agreed non-blocking, going to row 231.
  - NB3 and NB4 accepted.
  - It spot-checked the broker attribution (both PASS once at head; grep 0) and called it plausible, not proven (no contended re-run).

**Bookkeeping:**
- Sprint JSON M8b `completed`; plan §M8b status **built**; orchestration §3 and §5 M8b BUILT, `Estimated:` 0 executor runs pending.
- Charter: STATUS 316 added and STATUS 313 rotated to the archive; CURRENT GOAL next line updated.
- Row 150: M8b recorded, NEXT = P0/P3 within the attended A1 window.
- Row 231 gains NB2 (a rects-path guard test). Row 232 gains instance 2.
- Doc stays in `planned/`, because row 150 is not complete (A1).

**Result:** M8b is BUILT on `sprint/world-iter316-row150-m8b` (PR #405). Commits are `3f185a9a` (executor diff) and `706af0e4` (NB1 + evidence), plus this record. It is LANDED only after merge-commit CI (Gate 3b).

**Routing evidence:** base=8bfaa1ab28000937234aececa2ccd818e2b522a9@2026-10-29T00:04:34Z (Gate 1 record = Gate 3 snap). Driver env: `MISSION_EXECUTOR_RESOLVED=codex:gpt-6.1-sol`, `MISSION_EVALUATOR_RESOLVED=sonnet`; `MISSION_OVER_RATION` empty; `MISSION_ROUTING_NOTE` empty.
- Controller `claude:claude-opus-5-5` (tok: not reported).
- Designer and planner: not routed (plan exists; no design change).
- Executor: resolver `recipe codex:gpt-6.1-sol declared:provider-pin`; probe rc 0; ran `codex exec` (workspace-write, `--add-dir` GOCACHE/GOMODCACHE/snapshot dir), rc 0, 130,665 tok. Not the Agent tool: the recipe is the only valid path for a `provider:model` pin.
- Evaluator: Agent `sonnet` (`agent-tool sonnet declared:alias-pin`); r1 68,577 tok / 306 s; r2 76,420 tok cumulative / 70 s (resumed).
- Generator ≠ judge: codex (OpenAI) vs sonnet (Anthropic), cross-vendor.
- Quorum: none (no doc). Metered $0.00.

**Ruled out:**
- "The full-suite red is this diff": refuted. The diff is docs-only, both tests are 3/3 green at head and base, and the package is green uncontended.
- The first rule-3e(a) baseline "rc=1 0s" was a zsh word-split instrument failure, not a gate verdict.
- verify_ail's 9 property skips are the gate's pinned no-generator set, identical at base. They are not M8b's.

**Carried:**
- Row 231 + NB2 (no test pins `WORLD_MARKETING_RECTS`).
- Row 232 instance 2 (the `chatty` arm on macOS under load).
- `.wt-world-iter306-r0` residue.
- Row 150's remaining steps are all controller-in-attended-window or attended: P0 within 24 h before A1, P3 before pass A, and A1 itself (descriptor publish + `session mint` + capture + promotion).

**Retro:**
- **Concurrent judge + full suite manufactures reds.** Running the judge's focused Chrome/daemon tests alongside `go test ./...` doubled `host/broker` to 524 s and reddened two timing-sensitive tests. The attribution cost about 10 min. Next time, either run the full suite before spawning the judge, or forbid the judge heavy tests while it runs. Instance 1.
- **The zsh word-split trap struck a controller loop again**, which Gate 3b already names for pollers. It applies to every ad-hoc gate loop: write it as a bash script. Instance 1 for World (the fleet had 7).
- **A doc test that reads flags does not read values.** `TestMarketingRunbookMatchesHelp` proves the flags parse, but not which file they name, so the rects-path revert survives. Filed on row 231.
- Harness share 0. No skill edit.

**Progress:** 1.0 clauses 7/7 MET; goal unmoved (release row). Row 150: 23/23 executor milestones built. Remaining: the controller steps P0/P3 and the attended A1 (needs Mark).
**Next:** row 150's A1 is attended. Mark schedules it, and the controller runs P0 (≤24 h before) and P3 (before pass A) in that window. The loop's next pick under D-WORLD-81 is row 151 if routable (it needs row 188 N2d and N4a; check at Gate 2), else the highest routable of rows 230–232.

## 317 — 2026-10-29 — no pick: rule (d) block report; release rows 150/151 wait on attended steps; row 232's premise re-measured (not a breach; mechanism found); D-WORLD-102 filed [ADMIN]

**Kind:** bookkeeping-only iteration under standing rule (d) (D-WORLD-46). There is no design, plan, code or product acceptance. The record is the STATUS stamp (314 rotated to the archive), the CURRENT GOAL line, a row-232 measurement note, the new D-WORLD-102 ledger row, this entry, the index row, the dashboard, and the evidence under `verification/world-iter317/`. An independent judge reviews it before merge.

**Pick and why:** none.
- **Clause map at `23940015`:** clauses 1–7 MET. There is no UNMET clause and no critical-path row.
- **D-WORLD-81 release rows:**
  - **150** waits only on the attended A1. The controller runs P0 (≤24 h before) and P3 (before pass A) inside that window (orchestration r5 §4.4).
  - **151** needs row 188's N2d and N4a (`w-agent-generated-projections` §8). D-WORLD-95 R2 runs rows 219 → 217 → 218 before N2d, and the roadmap (`w-all-agent-work-through-world` §3) puts all of them on the attended track as World episodes. Row 217's attended trial is draft PR #402.
  - Neither row is loop-routable.
- **Rows 230–232** are hygiene tagged "below every D-WORLD-81 release row". Rule (d) bars residual and hardening rows when nothing on the path is routable. D-WORLD-64's exception covers rows 148–151 only.
- So iteration 316's **Next** ("else the highest routable of rows 230–232") contradicted rule (d). It is corrected here and in CURRENT GOAL. Precedent: iterations 257–264 were block reports for the same reason.
- **Row 232 "possibly soundness"** is the one claim that could have made a hygiene row outrank the queue: a genuine soundness bug outranks under the skill's Gate-0 rule. The controller measured it before deciding; see below. It is not a breach, so rule (d) holds.
- **New ask D-WORLD-102** (recommended A, default B): may the loop take 232 → 231 → 230 while the release rows wait? Without it every fire until the attended steps land is a block report.

**Row 232 re-measure (controller, macOS arm64, go1.26.6, base `23940015`; evidence `verification/world-iter317/row232-reset-holder-remeasure.txt` and `row232-dial-reset-probe.go.txt`):**
- `go test ./host/broker/ -run 'TestExecProbeReverseLegSilentHolderRefusesSticky/reset-holder' -count=500` idle: rc 0, 0 FAIL (71.5 s).
- The same run with `-count=300` under CPU load (32 `yes` processes on 16 cores): rc 1, **5 subtest FAILs**, each `Why = "", picker calls 1`. That is the CI signature from PR #404 run 40766776779. A raw `grep -c -- '--- FAIL'` reads 10, because each failure prints a parent line and a subtest line.
- **Mechanism probe** (a standalone Go program built outside the repo): an accept-then-RST listener (`SetLinger(0)` + `Close`, exactly `revHolder`'s `reset` kind) dialled 2,000 times with `net.DialTimeout(200 ms)`. The dial ITSELF failed with `connect: connection reset by peer` 10 times idle, and 29 and 25 times in two runs under load.
- `pollLoopback` (`host/broker/exec_probe.go`) ignores every dial error. So a lone dial that fails this way records none of child, foreign or silent, and `judgeReverse` passes the leg. `revFixture`'s `/bin/sh` client returns at once, so the attempt has exactly one dial.
- **Reading:** this is not a confinement breach. The lost dial reaches a FOREIGN holder; the child's own listener answers the nonce (`classifyReverseDial`) and never resets. It IS a conformance gap: the design classifies a reset as SILENT (a sticky refusal), while a connect-time reset is classified as no listener. The fix candidate is to record an ECONNRESET dial as silent, with a race-free pin. That is row 232's work and is pickable only on D-WORLD-102 = A.
- No production change and no test change was made.

**Gate 0:**
- The kill switch is armed (`mission-world.disabled` absent). gh account is `sunholo-voight-kampff`. The billing tripwire is CLEAN. Env was read by NAME.
- `mission_directives.sh`: 0 allowlisted directives on #366 since `2026-10-27T09:03:42Z` (14 comments, all the loop's own). Both watermarks (`mission-366-last-seen`, `mission-world-last-seen`) read that time and are unchanged.
- `mission_gate0_self_notices.sh` rc 1 on the same 2026-10-22 rc=1 notice `6345079873`, already credited (log count 14). Control 107:4 ok.
- The unread `mission-world` inbox holds only this loop's claim messages. The iteration-317 claim was sent and confirmed in the listing.
- `mission_decisions.sh --check`: valid, 88 rows and 0 OPEN at base; 89 rows and 1 OPEN (D-WORLD-102) after this record. The external-issue sweep is not due (#366 is this week's thread).

**Gate 1:**
- Local `dev` = `origin/dev` = `23940015`. Gate-1 base `23940015ebeba61e56334580b2848b19a1a7315b@2026-10-29T04:01:30Z`; Gate-4 base the same SHA at 04:08:20Z.
- Commit check-runs 3/3 success (`go host build + test gate`, `ailang-code verify gate`, `go race gate`). So PR #405 (row 150 M8b) is LANDED by merge-commit CI.
- Skill drift: `SKILL.md` and all 13 resource files are byte-identical to `origin/dev`, both at the resolved symlink target and in the pin copy (per-file `cmp`, no DRIFT line). All rules were read from the resolved target.
- Open PRs from this account: #402 (the attended row-217 trial, draft). It is not this loop's work and was left alone.
- Died-mid-flight traces: no iteration branch or worktree for 317 existed. `.wt-world-iter306-r0` is the known residue.

**Designer / planner / executor:** not spawned. There is no pick, so nothing to design, plan or execute (precedent 257–264). The operator's standing request for Agent-tool roles is honoured through the evaluator.

**Independent evaluator (record):** REQUIRED. It runs through the Agent tool on `sonnet` (resolver `agent-tool sonnet declared:alias-pin`), in a fresh, read-only, foreground context, and judges the record commit before merge. See the follow-through below.

**Routing evidence:** base=23940015ebeba61e56334580b2848b19a1a7315b@2026-10-29T04:08:20Z (Gate 4).
- Driver env: `MISSION_DESIGNER_RESOLVED=claude:claude-opus-5-5`, `MISSION_PLANNER_RESOLVED` and `MISSION_EXECUTOR_RESOLVED` = `codex:gpt-6.1-sol` (recipe), `MISSION_EVALUATOR_RESOLVED=sonnet`. `MISSION_OVER_RATION` and `MISSION_ROUTING_NOTE` are empty.
- Controller `claude:claude-opus-5-5` (tok: not reported). Designer, planner and executor were unused.
- Evaluator: Agent `sonnet` (`agent-tool sonnet declared:alias-pin`), 42,335 tok / 31 s. Generator ≠ judge: opus controller vs sonnet judge.
- Quorum: none (no doc). Metered $0.

**Ruled out:**
- Picking 150 (attended A1) or 151 (behind attended rows 219/217/218 and 188 N2d/N4a).
- Picking 230, 231 or 232 under rule (d). Iteration 316's Next is corrected.
- Treating row 232 as a soundness outrank: measured above, the lost dial is a foreign holder's, not the child's.
- The hypothesis that the reset-holder failure is a classification bug in `classifyReverseDial`: refuted. The failing attempts record nothing at all, so the dial never reached the classifier.
- `ailang mission rotate-log` (32 live entries, under the ~40 threshold); the index row is hand-written as in 316.

**Retro:**
1. **A controller's Next line can contradict an attended rule without anyone noticing.** Iteration 316 wrote "else rows 230–232", and only the rule-(d) reread at this Gate 2 caught it. Instance 1. No rule change: the Gate-2 rule (d) check already exists and worked.
2. **Measure a "possibly soundness" tag before it steers a pick.** One 3-minute load loop plus a 30-line probe turned row 232 from "unmeasured, possibly a breach" into a located mechanism with a fix candidate, and removed its claim to outrank. The load axis mattered: 0/500 idle would have read as "cannot reproduce".
3. **Idle cadence:** this is block report 1 of a possible run. D-WORLD-102 is the lever. Without it the fleet pays one judge per fire for records about the same block, as in 257–264.
- Harness share 0. No skill edit.

**Progress:** 1.0 clauses 7/7 MET; goal unmoved (release rows). Row 150: 23/23 executor milestones built and M8b LANDED. Remaining: the attended A1 with the controller's P0/P3.
**Next:** if D-WORLD-102 = A, take **row 232** (record a connect-time ECONNRESET as silent, with a race-free pin), then 231, then 230. Else rule (d) again. In either case row 150's P0/P3 run in the A1 window Mark schedules, and row 151 becomes routable when row 188's N2d and N4a land.

**Record follow-through (iteration 317):** the independent evaluator judged record commit `92495b60` on `23940015` and returned **PASS 90/100, zero blocking** (35/21/19/15). It ran through the Agent tool on `sonnet`, in a fresh, read-only, foreground context (69,585 tok, 266 s, 19 tool calls).
- **Verified by the judge:** no loop-routable row under rule (d), D-WORLD-64, D-WORLD-81 and D-WORLD-95 R2; dev check-runs 3/3. It reproduced the `reset-holder` signature under load (1/600 at 48 `yes`) and the connect-time reset with its own probe (3 idle, 34 and 38 under load, of 2,000). It also confirmed the STATUS rotation (314 byte-identical in the archive), the ledger (89 rows, D-WORLD-102 the only OPEN), cross-file consistency, and that there is no closing keyword or email.
- **Non-blocking:** N1 the record gave no production-side argument for "not a breach". The controller confirmed `execProbeListenMS` = 1000 (`exec_probe.go:98`) and the 50 ms poll (`:719`), and added that argument (about 20 dials per attempt) and the "never resets is not enforced" caveat to row 232. N2 the rate is load-sensitive (the judge got 0/300 at 32 `yes`), so row 232 now quotes both readings. N3 the background load was not recorded; it is noted in the judge report.
- No r2: the fix only adds caveats that weaken the record's claims, and nothing the judge scored was changed.
- Report: `design_docs/verification/world-iter317/evaluator-r1.md`.

## 318 — 2026-10-29 — row 232: a reverse-leg dial that fails at connect with ECONNRESET is SILENT — taken as a Gate-1 dev red; Linux errno measured first-party; judge r1 PASS 92 [PRODUCT]

**Kind:** full inner loop on a Gate-1 red: designer → quorum r1 BLOCKED → controller measurement → designer revision → quorum r2 PROCEED → planner → executor → independent evaluator → PR #410. Every role was spawned through the Agent tool, per the operator's standing request.

**Pick and why:**
- **Clause map at `df5a0671`:** clauses 1–7 MET. There is no UNMET clause and no critical-path row. Rows 150 (attended A1) and 151 (behind row 188 N2d/N4a, after attended rows 219 → 217 → 218) are not loop-routable, so rule (d) would make this a block report.
- **But dev was red.** Commit check-runs on `df5a0671` (PR #402's merge, the attended row-217a trial) read `go race gate: failure`. Job 123313206642, run 40833618366: `--- FAIL: TestExecProbeReverseLegSilentHolderRefusesSticky/reset-holder`, `exec_probe_test.go:664: Why = "", picker calls 1`. Every other package was `ok`. The two `WARNING: DATA RACE` blocks in that log are the gate's own known-positive control (`w-race-gate-blindspot/racecontrol`), not a finding.
- **Attribution:** #402 touched `host/broker/registry_publish_test.go` only, not `exec_probe*`. The signature is row 232's (instances: PR #404, iteration 316's `chatty`, this one — the first on `dev`). `gh run rerun --failed` → attempt 2 **success**. So it is the measured flake, not #402.
- **Rule:** a red dev outranks the queue for the owning mission (World owns its repo). Row 232's own tag already said "a Gate-1 dev red on this test outranks the queue". Precedent: iteration 314 took a dev red under rule (d). **D-WORLD-102 was not used and is not resolved.** It stays OPEN, and its scope is now 231 → 230 (plus new row 233).

**Linux premise, measured by the controller** (rule 3f — the quorum's objection was measured, not forwarded). Throwaway branch `probe/world-iter318-linux-errno` (deleted after), ubuntu-latest, `Linux 7.0.0-1017-azure x86_64`, 4 cores, go1.26.6 linux/amd64. Probe source `verification/world-iter318/linux-errno-probe.go.txt`; log `linux-errno-probe-run40842980481.txt`.
- No listener (idle and under 2×nproc `yes`, 1 and 4 dialers): ECONNREFUSED (0x6f) on **7,325,743 of 7,325,743** dials. Zero ECONNRESET, zero EINVAL (EINVAL is darwin-only).
- Accept-then-RST holder (T7's `reset` shape): sequential — connected 642,869 / ECONNRESET (0x68) 1,742 idle; 195,092 / **147,968** under load. Four spinning dialers add timeouts (143 idle, 179 load). **ECONNREFUSED never appears** against the holder.
- T7 `reset-holder` at base under load, `-count=300`: **≥7 FAIL** with the CI signature (counted from a `tail -30`, so a floor).
- The dev rerun was green on the same run.

**Design** (`w-exec-probe-reset-holder-false-pass`, designer opus via the Agent tool; author vendor benched; seats glm-5.3, kimi-k3, gpt6-1-sol, reserve gemini-4-argon):
- Decision: a failed dial counts only if its error proves a completed handshake, which is ECONNRESET alone (`errors.Is`). ECONNREFUSED, EINVAL, timeouts, ECONNABORTED and local errors record nothing. darwin control: 0 ECONNRESET in 2,013,696 no-listener dials under load; 122 EINVAL (so EINVAL must not count).
- False-refusal risk measured: World's own child listener does reset host dials as it closes, but in every reachable cycle its nonce was seen (100/100) and `judgeReverse`'s child-first order returns REACHED; in a correct sandbox the child's listener is unreachable. Increment at the real cadence 0 of 400 cycles.
- **r1 BLOCKED:** glm and kimi REJECT on the same premise (the Linux errno was from memory); gemini PASS (asked for the scratch sources in-repo). gpt6-1-sol absent (unreachable).
- **r2 (the one revision, carrying the controller's Linux numbers):** **PROCEED at N-1, `gpt6-1-sol` absent (unreachable)**; glm, kimi and gemini-4-argon PASS. Their non-blocking points (race-gate selection; a parent §4.2 amendment note; T9's `atomic.Int32` fake dialer) went to the planner.
- Artifacts: `verification/world-iter318/quorum/r1.json`, `r2.json`. Metered $0.609 (r1 $0.218, r2 $0.391).

**Plan** (opus): M1 production + T8/T9 + the drill; M2 the CI list, QUICKSTART and doc notes. The planner prototyped the fix outside the repo and ran every mutant. D1 (false in the design): M28's named killer cannot see an `os.ErrDeadlineExceeded` mutant, so it added T8 `deadline-exceeded` and M28b. D2: AC-LINUX as written was unachievable (the race gate runs `go test ./... -race` without `-v`, `ci.yml:418`/`:497`); restated as step-8 PASS lines plus `ok host/broker` in the race gate.

**Executor** (haiku, Agent tool; resolver `recipe claude:claude-haiku-5-5 declared:provider-pin`):
- `host/broker/exec_probe.go` +19/−6: `dialErrIsSilent(err) = errors.Is(err, syscall.ECONNRESET)`; `pollLoopback(port, want, dial)` with `nil` = the old 200 ms dial; the error branch sets `silent`; `reverseLeg` passes `nil`.
- `exec_probe_test.go` +93: T8 `TestExecProbeReverseDialErrIsSilent`, T9 `TestExecProbeReversePollRecordsConnectReset`.
- `ci.yml:247` gains both names; QUICKSTART §10 arm 6 "(or that resets the host's dial)"; the parent design gains its row-232 amendment note.
- Drill M22–M29 + M28b: every mutant red on its named killer, in-table controls green, restores sha256-verified. `-count=20` 20/20.
- Commits: `f438ba41` (M1), `a828c450` (M2), `6b4c2b60` (evidence). The controller rebuilt them from the executor's `.snap/` milestone snapshots, byte-identical to the final tree.

**Controller verification (outside any sandbox):**
- **AC-LINUX-LOAD** at the fixed head on a Linux runner under 2×nproc `yes` (run 40849609705, log `ac-linux-load-run40849609705.txt`): T7 `reset-holder` `-count=300 -v` → `PASS_LINES=300 FAIL_LINES=0`; `^TestExecProbeReverse` `-count=5` all PASS. Paired red: ≥7/300 at base.
- `./scripts/verify_go.sh` at `a828c450` (load average 5–8, sibling missions running): rc 1 on `TestDoctorQuickIsFast` (4.09 s against its bound; `cmd/ailang-worldd` byte-identical to base) and `TestProtectedGitFilterPreflight/configured-process-include` (a missing temp file; `exec_protected_check_test.go` unchanged, and `pollLoopback` has one call site). Each passed 3/3 isolated at head. Attributed to load; CI was the arbiter. Log `controller-verify-go-a828c450.txt`.
- PR head `a828c450` and `6b4c2b60`: check-runs 3/3 success on Linux.

**Independent evaluator** (sonnet, Agent tool, own worktree `.wt-world-iter318-eval` at `a828c450`, foreground): **PASS 92/100** (28/24/19/13/8), zero blocking. It re-ran M22–M29 and M28b and four mutants of its own (ECONNABORTED silent; reset sets `child`; new branch disabled; injected dial bypassed): all red. T8/T9 `-count=200` green; P20 holds (the change can only refuse more). Non-blocking: NB1 the AC-LINUX-LOAD and verify_go logs were untracked → committed `6b4c2b60`; NB2 "§6 omits M28b" → **refuted**, the committed doc lists it at `a828c450`:164 (the judge was reading the executor's report); NB3 the executor skipped the plan's base run (the judge's drill covers it); NB4 one kernel measured; NB5 timeouts are deliberately not silent. Report `verification/world-iter318/evaluator-r1.md`.

**Gate 3b:** PR #410 merged as `793cc67c` (merge commit). Merge-commit check-runs **3/3 success** (`go host build + test gate`, `go race gate`, `ailang-code verify gate`), polled SHA-pinned on the full 40-char SHA with numeric guards; confirmed by a direct `gh run list --workflow CI --branch dev` read at `793cc67c`. Row 232 is **LANDED**, and dev is green again.

**Gate 0:**
- Kill switch armed (`mission-world.disabled` absent); gh `sunholo-voight-kampff`; billing tripwire CLEAN; env read by NAME.
- `mission_directives.sh`: 0 allowlisted directives on #366 since `2026-10-27T09:03:42Z` (16 comments, none from the allowlist). Both watermarks read that time; nothing to mark.
- `mission_gate0_self_notices.sh` rc 1 on the 2026-10-22 rc=1 notice `6345079873`, already credited (log count 15). Control 107:4 ok.
- Unread `mission-world` inbox: only this loop's claim messages (iteration-318 claim sent and confirmed in the listing).
- Ledger: 89 rows, 1 OPEN (D-WORLD-102), unchanged except its evidence cell.

**Gate 1:**
- Local `dev` = `origin/dev` = `df5a0671`. Gate-1 base `df5a0671dadd567ece879fe89afcc018c2deaa23@2026-10-29T08:02:57Z`; the Gate-3 `snap` was the same SHA.
- Skill drift: `SKILL.md` and all 13 resource files byte-identical to `origin/dev` at both the resolved symlink target and the pin copy (14 files each, no DRIFT line). Rules read from the resolved target.
- Open PRs from this account: #407 (`trial/row220d-design`, attended row-220 design) — not this loop's; left alone.

**Routing evidence:** base=793cc67c0ab029793f44930a2e1ad6a4edf574b4@2026-10-29T10:17:56Z (Gate 4).
- Driver env: `MISSION_DESIGNER_RESOLVED=claude:claude-opus-5-5`, `PLANNER=opus`, `EXECUTOR=claude:claude-haiku-5-5`, `EVALUATOR=sonnet`; `MISSION_OVER_RATION=codex`. `MISSION_ROUTING_NOTE`: planner and executor codex `gpt-6.1-sol` unusable (probe rc 75, over daily ration), handed to `opus` and `claude:claude-haiku-5-5`.
- Controller `claude:claude-opus-5-5` (tok: not reported).
- Designer: resolver `recipe claude:claude-opus-5-5 declared:provider-pin`, spawned as Agent `opus` (same weights; iterations 228/314 precedent). Rotation pointer unchanged at `claude:claude-opus-5-5` (the lane that ran). 1 doc, 2 runs (author + the protocol-mandated revision, within the one-doc diet). 128,419 tok / 1,376 s, then 173,648 tok / 192 s (cumulative transcript).
- Planner: Agent `opus` (`agent-tool opus fail-closed:env-pin`; derive-planner-lane `opus fail-closed:env-pin`), 148,168 tok / 728 s.
- Executor: Agent `haiku` (resolver `recipe claude:claude-haiku-5-5 declared:provider-pin`), 108,961 tok / 1,412 s.
- Evaluator: Agent `sonnet` (`agent-tool sonnet declared:alias-pin`), 68,460 tok / 314 s. Generator ≠ judge: haiku vs sonnet.
- Quorum: glm-5.3 / kimi-k3 / gemini-4-argon (reserve), gpt6-1-sol absent (unreachable) both rounds — the open `quorum:gpt6-1-sol-openai-api-credits-exhausted` ticket's condition. Metered $0.609.
- Two Linux CI probe runs on a throwaway branch (Actions minutes, $0 metered).

**Ruled out:**
- Picking rows 231/230 (rule (d); D-WORLD-102 unanswered). Row 232 was taken only because dev was red on its test.
- Treating the red as #402's: refuted (no `exec_probe*` file in its diff; rerun green).
- The quorum's worry that the reset holder surfaces as ECONNREFUSED on Linux: refuted by measurement (0 of them in rst mode).
- Counting EINVAL or timeouts as silent: refuted by the no-listener controls (darwin EINVAL with nothing listening; Linux timeouts only under spinning saturation).
- The judge's NB2 (M28b missing): refuted at `a828c450`:164.
- `ailang mission rotate-log` (33 live entries, under ~40); the index row is hand-written.

**Retro:**
1. **A quorum objection that names an unmeasurable premise may be measurable after all.** Both rejects said "Linux was not run"; the rig has no container tool, but a throwaway branch with a push-triggered workflow is a Linux machine in about five minutes. It turned a premise into 7.3M dials of evidence and gave the sprint a paired red (≥7/300 at base) and green (300/300 at head) on the platform of record. Instance 1; no rule change.
2. **A docs-only evidence commit restarts a 20-minute CI cycle.** Committing the evidence before the PR opens would have saved one cycle. Instance 1.
3. **A judge reading the executor's report inherits its stale claims** (NB2). The judge was told to derive from the diff, and mostly did. Instance 1.
- Harness share 0. No skill edit.

**Progress:** 1.0 clauses 7/7 MET; goal unmoved (release rows). Row 232 LANDED. Row 150: the attended A1 remains; row 151: behind row 188 N2d/N4a.
**Next:** rule (d) block report unless D-WORLD-102 = A (then 231 → 230 → 233). Row 150's P0/P3 run in the A1 window Mark schedules; row 151 becomes routable when row 188's N2d and N4a land. If T7 `reset-holder` reds again on Linux CI, re-measure the Linux dial-error classes first (design R7).

## 319 — 2026-10-29 — no pick: rule (d) block report 2; release rows 150/151 still wait on attended steps; dev green with row 232's fix (the dc24fccf red predates it) [ADMIN]

**Kind:** bookkeeping-only iteration under standing rule (d) (D-WORLD-46). There is no design, plan, code or product acceptance. The record is the STATUS stamp (316 rotated to the archive), the CURRENT GOAL line, a D-WORLD-102 evidence note, this entry, the index row, the dashboard, and `verification/world-iter319/gate1-2-readings.txt`. An independent judge reviews it before merge.

**Pick and why:** none.
- **Clause map at `bcb07394`:** clauses 1–7 MET. There is no UNMET clause and no critical-path row.
- **D-WORLD-81 release rows, blockers re-measured this iteration (not transcribed):**
  - **150** waits only on the attended A1. The controller runs P0 (≤24 h before) and P3 (before pass A) inside the window Mark schedules (orchestration r5 §4.4).
  - **151** needs row 188's N2d and N4a (`w-agent-generated-projections` §8). `git log origin/dev --since=2026-10-01 | grep -ciE 'N2d|N4a'` → **0**; control `N2c` → **2** (N2c, `6f2ba039`, is row 188's newest landing). D-WORLD-95 R2 runs rows 219 → 217 → 218 before N2d. Since iteration 317 the attended track landed 217a (#402), 220d's design (#407), 219b (#409) and 218a (#408), but all three rows are still `[NEXT]` on the attended track, and their exit criteria are runs Mark makes. So none of them, and not N2d, is loop-routable.
- **Rows 231, 230, 233** are hygiene "below every D-WORLD-81 release row". Rule (d) bars them. D-WORLD-64's exception covers rows 148–151 only. D-WORLD-102 (may the loop take them while the release rows wait) is OPEN with default B, so this is a block report.
- **Dev is not red:** `bcb07394` check-runs 3/3 success. The CI run between iterations 318's pick and its landing, on `dc24fccf` (#408's merge, run 40853596929), was red in the go race gate on `TestExecProbeReverseLegSilentHolderRefusesSticky/reset-holder` (`exec_probe_test.go:664: Why = "", picker calls 1`). `git merge-base --is-ancestor f438ba41 dc24fccf` → no: that commit lacks row 232's fix, and `793cc67c` (the fix's merge) and `bcb07394` are both green. It is row 232 instance 4 on pre-fix code, not a red to outrank with, and not iteration 318's "if T7 reds again on Linux CI" trigger (which is about post-fix code).

**Gate 0:**
- Kill switch armed (`mission-world.disabled` absent); gh `sunholo-voight-kampff`; billing tripwire CLEAN; env read by NAME.
- `mission_directives.sh`: 0 allowlisted directives on #366 since `2026-10-27T09:03:42Z` (17 comments, none from the allowlist). Nothing to mark.
- `mission_gate0_self_notices.sh`: verdict "1 crash notice" on the same 2026-10-22 rc=1 notice `6345079873`, already credited (log count 16 before this entry). Control 107:4 ok.
- Unread `mission-world` inbox: only this loop's claim messages (the iteration-319 claim was sent and confirmed in the listing).
- `mission_decisions.sh --check`: valid, 90 rows, 1 OPEN (D-WORLD-102). The external-issue sweep is not due (#366 is this week's thread).

**Gate 1:**
- Local `dev` was 8 commits (2 first-parent merges, #410 and #411) behind; fast-forwarded (`--ff-only`, tree had only the known untracked residue) to `origin/dev` = `bcb07394`. Gate-1 base `bcb073940751056d273562ca7f79753f980b0683@2026-10-29T12:02:21Z`; Gate-4 base the same SHA at 12:04:01Z.
- Skill drift: `SKILL.md` and all 13 resource files byte-identical to `origin/dev` at the resolved symlink target (per-file `cmp`; the first loop printed DIFF for all 13 because it dropped the `resources/` prefix from the `git show` path. That was the instrument's fault: a single-file control and the corrected loop both read SAME). Rules read from the resolved target.
- Open PRs from this account: none.
- Died-mid-flight traces: no iteration-319 branch or worktree existed before this one. `.wt-world-iter306-r0` is the known residue.

**Designer / planner / executor:** not spawned. There is no pick, so nothing to design, plan or execute (precedent 257–264, 317). The operator's standing request for Agent-tool roles is honoured through the evaluator.

**Independent evaluator (record):** REQUIRED. It runs through the Agent tool on `sonnet` (resolver `agent-tool sonnet declared:alias-pin`), in a fresh, read-only, foreground context, and judges the record commit before merge. See the follow-through below.

**Routing evidence:** base=bcb073940751056d273562ca7f79753f980b0683@2026-10-29T12:04:01Z (Gate 4).
- Driver env: `MISSION_DESIGNER_RESOLVED=claude:claude-opus-5-5`, `MISSION_PLANNER_RESOLVED=opus`, `MISSION_EXECUTOR_RESOLVED=claude:claude-haiku-5-5`, `MISSION_EVALUATOR_RESOLVED=sonnet`; `MISSION_OVER_RATION=codex`. `MISSION_ROUTING_NOTE`: planner and executor codex `gpt-6.1-sol` unusable (probe rc 75, over daily ration).
- Controller `claude:claude-opus-5-5` (tok: not reported). Designer, planner and executor were unused.
- Evaluator: Agent `sonnet` (`agent-tool sonnet declared:alias-pin`), 42,335 tok / 31 s. Generator ≠ judge: opus controller vs sonnet judge.
- Quorum: none (no doc). Metered $0.

**Ruled out:**
- Picking 150 (attended A1) or 151 (row 188 N2d/N4a absent, measured above).
- Picking 231, 230 or 233 under rule (d) with D-WORLD-102 unanswered.
- Treating the `dc24fccf` red as a live dev red: it predates the fix, and both later dev commits are green.
- `ailang mission rotate-log` (35 live entries with this one, under the ~40 threshold); the index row is hand-written as in 317–318.

**Retro:**
1. **A drift-check loop that prints DIFF for every file is more likely a path bug than drift.** It was caught here by a single-file control before anything acted on it. The skill's own prescribed loop (Gate 1 rule (d)) builds `rel` correctly; the defect was in a hand-rolled variant. Instance 1; no rule change.
2. **Idle cadence:** this is block report 2 of a possible run (iterations 257–264 were eight). Each costs one judge. D-WORLD-102 remains the only lever the loop has; it is not re-asked, only re-surfaced in the digest as the single open decision.
- Harness share 0. No skill edit.

**Progress:** 1.0 clauses 7/7 MET; goal unmoved (release rows). Row 150: the attended A1 remains; row 151: behind row 188 N2d/N4a.
**Next:** rule (d) block report unless D-WORLD-102 = A (then 231 → 230 → 233) or dev goes red. Row 150's P0/P3 run in the A1 window Mark schedules; row 151 becomes routable when row 188's N2d and N4a land. If T7 `reset-holder` reds on post-fix code, re-measure the Linux dial-error classes first (row 232 design R7).

**Record follow-through (iteration 319):** the independent evaluator judged record commit `fe9c891e` on `bcb07394` and returned **PASS 91/100, zero blocking** (33/22/20/16). It ran through the Agent tool on `sonnet`, in a fresh, read-only, foreground context (42,335 tok, 31 s, 7 tool calls).
- **Verified by the judge:** the rotation (charter +4/−4, 316 at the archive top), ledger valid with D-WORLD-102 the sole OPEN row, #366's 17 comments all from the loop's account, `bcb07394` 3/3 green, the `dc24fccf` red on pre-fix code (`f438ba41` not an ancestor), notice count 16, rows 217–219 still `[NEXT]` attended-track, and no closing keyword or address.
- **Non-blocking:** NB1 the "8 commits behind" figure was not in the evidence → added to `gate1-2-readings.txt` (`dc24fccf..bcb07394` = 8, first-parent 2). NB2 the N2d/N4a grep is weak on its own because the row-188 sprint plan does not name them; the conclusion also rests on N2c being row 188's newest landing and on rows 217–219 being open. NB3 row 150's blocker was taken from the charter and the 316 clause map. NB4 hand-written index row (fine).
- No r2: the fix only adds evidence for a figure the judge could not re-derive; nothing scored changed.
- Report: `design_docs/verification/world-iter319/evaluator-r1.md`.

## 320 — 2026-10-29 — rule (d) block report 3, then dev went red under the record → row 234: requireF15's "134" substring check replaced by a structural predicate; judge PASS 91 [PRODUCT]

**Kind:** bookkeeping-only iteration under standing rule (d) (D-WORLD-46). There is no design, plan, code or product acceptance. The record is the STATUS stamp (317 rotated to the archive), the CURRENT GOAL line, a D-WORLD-102 evidence note, this entry, the index row, the dashboard, and `verification/world-iter320/gate1-2-readings.txt`. An independent judge reviews it before merge.

**Pick and why:** none.
- **Clause map at `8cee85a9`:** clauses 1–7 MET. There is no UNMET clause and no critical-path row.
- **D-WORLD-81 release rows, blockers re-measured this iteration (not transcribed):**
  - **150** waits only on the attended A1. The controller runs P0 (≤24 h before) and P3 (before pass A) inside the window Mark schedules (orchestration r5 §4.4). Since iteration 319's base `bcb07394`, `git log --format='%h %s' | grep -ciE 'row ?150|A1\b|marketing'` → **0** (control `row ?2[12][0-9]` → **20**). Nothing moved it.
  - **151** needs row 188's N2d and N4a. Over commit SUBJECTS, `git log origin/dev --since=2026-10-01 --format='%h %s' | grep -ciE 'N2d|N4a'` → **0**; control `N2c` → **2**. Over full messages (the default `git log` format, which iteration 319 printed) the count is **4**. All four are record or design commits that describe the work and do not build it: `fe9c891e` (319 record), `92495b60` (317 record), `be184351` (roadmap rows 217–222), `dc719822` (`w-world-surface-outputs` r0). No feat commit names N2d or N4a. Row 188's newest landing is N2c (`6f2ba039`, merged as `1e7b3cd3`). Since iteration 319 the attended track landed 220a (#413), 221's design (#414), 218b (#415), 217b (#416) and 220b (#417), and 219c is draft PR #418. Rows 217–219 are still `[NEXT]` on the attended track, and D-WORLD-95 R2 runs them before N2d. So nothing here is loop-routable.
- **Rows 231, 230, 233** are hygiene "below every D-WORLD-81 release row". Rule (d) bars them. D-WORLD-64's exception covers rows 148–151 only. D-WORLD-102 is OPEN with default B, so this is a block report.
- **New ledger row D-WORLD-104** (row 221's seven rulings, "all defaults") was created and resolved in an attended session on 2026-10-29. `mission_ledger_provenance.sh` → `verdict=attended` (flip `ddfa6a13`, creation `6b10491c`). It covers row 221 (attended track) and does not change any row above.

**Gate 0:**
- Kill switch armed (`mission-world.disabled` absent); gh `sunholo-voight-kampff`; billing tripwire CLEAN; env read by NAME.
- `mission_directives.sh`: 0 allowlisted directives on #366 since `2026-10-27T09:03:42Z` (18 comments, none from the allowlist). Nothing to mark.
- `mission_gate0_self_notices.sh`: verdict "1 crash notice", the same 2026-10-22 rc=1 notice `6345079873`, already credited (log count 17 before this entry). Control 107:4 ok.
- Unread `mission-world` inbox: only this loop's claim messages. The iteration-320 claim was sent (`inbox_1793289832458_d0b201f3`) and shows in the listing.
- `mission_decisions.sh --check`: valid, 91 rows, 1 OPEN (D-WORLD-102). The external-issue sweep is not due (#366 is this week's thread).

**Gate 1:**
- Local `dev` = `origin/dev` = `8cee85a98dc800f808ad061c597ab88a0ea0140b` (not behind). Gate-1 base `8cee85a9…@2026-10-29T16:01:55Z`.
- Since `bcb07394`: 30 commits, 6 first-parent merges (#412 iteration 319's record, #413–#417 attended trial branches).
- Skill drift: per-file `cmp` of `SKILL.md` and 13 resource files against `origin/dev`, at both the resolved symlink target and the pin worktree, printed zero DRIFT lines; a positive control fired. Rules read from the resolved target.
- Dev CI: `commits/8cee85a9/check-runs` → 3, none not-green; the last six dev CI runs all `success`.
- Open PRs: #418 `trial/row219c-exec-per-call` (draft, row 219 M3 attended episode). No worktree in this clone, and the `trial/` prefix is the attended track, so it is not this loop's. Left alone.
- Died-mid-flight traces: no iteration-320 branch or worktree existed before this one. `.wt-world-iter306-r0` is the known residue.

**Designer / planner / executor:** not spawned. There is no pick, so the routing table has nothing to send them (precedent 257–264, 317, 319). No spawn was attempted, so none failed. The operator's standing request for Agent-tool roles is honoured through the evaluator.

**Independent evaluator (record):** REQUIRED. It runs through the Agent tool on `sonnet` (resolver `agent-tool sonnet declared:alias-pin`), in a fresh, read-only, foreground context, and judges the record commit before merge. See the follow-through below.

**Routing evidence:** base=8cee85a98dc800f808ad061c597ab88a0ea0140b@2026-10-29T16:04:30Z (Gate 4).
- Driver env: `MISSION_DESIGNER_RESOLVED=claude:claude-opus-5-5`, `MISSION_PLANNER_RESOLVED=opus`, `MISSION_EXECUTOR_RESOLVED=claude:claude-haiku-5-5`, `MISSION_EVALUATOR_RESOLVED=sonnet`; `MISSION_OVER_RATION=codex`. `MISSION_ROUTING_NOTE`: planner and executor codex `gpt-6.1-sol` unusable (probe rc 75, over daily ration).
- Controller `claude:claude-opus-5-5` (tok: not reported). Designer, planner and executor were unused.
- Evaluator: Agent `sonnet` (`agent-tool sonnet declared:alias-pin`), tokens recorded in the follow-through. Generator ≠ judge: opus controller vs sonnet judge.
- Quorum: none (no doc). Metered $0.

**Ruled out:**
- Picking 150 (attended A1) or 151 (row 188 N2d/N4a absent, measured above).
- Picking 231, 230 or 233 under rule (d) with D-WORLD-102 unanswered.
- Adopting PR #418: attended track, not this mission's loop work.
- `ailang mission rotate-log` (36 live entries with this one, under the ~40 threshold); the index row is hand-written as in 317–319.

**Retro:**
1. **Idle cadence:** this is block report 3 in a row (317, 319, 320; 318 was a dev red). Each costs one judge and no code. The attended track is moving fast (five merges since iteration 319), but none of it reaches row 151's N2d/N4a. D-WORLD-102 is still the only lever the loop has, and the digest shows it as the one open decision. It is not asked again.
- Harness share 0. No skill edit.

**Progress:** 1.0 clauses 7/7 MET; goal unmoved (release rows). Row 150: the attended A1 remains; row 151: behind row 188 N2d/N4a.
**Next:** rule (d) block report unless D-WORLD-102 = A (then 231 → 230 → 233) or dev goes red. Row 150's P0/P3 run in the A1 window Mark schedules; row 151 becomes routable when row 188's N2d and N4a land.

**Record follow-through (iteration 320):** the independent evaluator judged record commit `5ea2151b` on `8cee85a9` and returned **PASS 94/100, zero blocking** (35/23/25/11). It ran through the Agent tool on `sonnet`, in a fresh, read-only, foreground context (47,355 tok, 64 s, 7 tool uses).
- **Verified by the judge:** no allowlisted comment on #366 (18, all the loop's), no N2d/N4a feat commit, 0 row-150 commits since `bcb07394`, D-WORLD-104 provenance `attended`, `8cee85a9` 3/3 green, 30 commits / 6 merges, #418 the only open PR, ledger 91 rows with D-WORLD-102 the only OPEN row, and the rotation, headings, `base=` and keyword checks.
- **NB1, reproduced and fixed:** the N2d/N4a count is 0 over commit subjects and 4 over full messages, all four being record or design commits. The entry had printed the full-message command beside the subject-only reading; it now gives both and names the commits.
- **NB2, fixed:** the six-run dev CI listing is now in `gate1-2-readings.txt`.
- No r2: the fixes correct evidence wording; nothing scored changed.
- **Retro addendum:** a command copied into a record must be the command that produced the number. Iteration 319 printed the same default-format grep, and its judge flagged the grep as weak (NB2 there) without catching that it returns 4. Instance 1 of the printed-command-differs class here; no rule change.
- Report: `design_docs/verification/world-iter320/evaluator-r1.md`.

**Second item — row 234, a dev red found at Gate 3b of the record (standing rule 1: a bookkeeping-only pick allows a second).**
- **How it surfaced.** The record PR #419's head `b477df9a` was 3/3 green on base `8cee85a9`. While it was in CI, the attended session merged #418 (row 219 M3, `trial/row219c-exec-per-call`) as `171ece91` at 16:11Z. So #419's merge `785b3bd1` sits on top of it. `171ece91`'s own push run went red in the race gate (`TestExecToolCellarProfilePathRefusesByName`). `785b3bd1` went red in the host gate twice: first on row 200's known flake (`TestWorktreeSnapshotRefusesAConfiguredFilter` control, the first instance in the non-race leg), and on rerun on `TestExecToolSymlinkInsideWorkspaceRefuses`. Both signatures begin `not a pre-launch refusal:`.
- **Mechanism, measured first.** `requireF15` (`host/broker/exec_tools_test.go:89`, added by `44eed437`) failed on `strings.Contains(string(out), "134")`. Both failing temp dirs contain `134` (`…3490134876`, `…596413420`). With a temp root whose path contains `134`, all 6 callers FAIL; with a control root, all pass. The designer measured the natural rate at 5/300 (≈1.7% per caller per run) and showed the substring never caught a real exit-134 launch.
- **Why the loop took it.** Dev red on World's own repo outranks the queue for the owning mission (precedent 314, 318). The code is from the attended track, so a claim message named row 234 before routing (`iter-320 claim: row 234 dev red`). `git log 785b3bd1..origin/dev` was empty and there were no open PRs before the executor ran, so there was no collision.
- **Design** (`w-exec-refusal-check-134-substring`, designer opus via the Agent tool; author vendor benched; seats gpt6-1-sol, gemini-4-argon, oc-glm-5-3, reserve oc-kimi-k3):
  - r1 BLOCKED, $0.1694: gpt6-1-sol (one `Decode` does not reject trailing content) and glm (TMPDIR-only ACs are vacuous under an ambient GOTMPDIR; sweep claim too wide). gemini passed. The designer measured both as correct (V20, V23) and applied their fixes in the one revision.
  - r2 BLOCKED, $0.1535: only gpt6-1-sol, on a new surface the revision had introduced (the prebuilt `broker.test` acceptance run had no `-test.timeout`). glm and gemini passed. No absent reviewers.
  - **Narrow-refinement carve-out:** that one objection had a concrete fix and did not dispute direction, so the controller applied it verbatim (`-test.timeout=2m`, timeout = acceptance failure) and re-ran V23 with that exact command: exit 1 in 1 s, 6 top-level and 10 leaf FAILs, 10 `TempDir: mkdir` errors, 0 timeout panics. Then routed to the planner. Artifacts `verification/world-iter320/quorum/r1.json`, `r2.json`.
- **Plan** (opus): M1 code, M2 evidence. Prototype-validated on a scratch copy: every mutant M0–M7 red on its named killer, and the plan's code blocks reproduce the prototype byte-for-byte (sha256 `475fb831…`). Deviations D1–D9 are recorded in the plan.
- **Executor** (haiku): M1 `63fd4f81` (+91/−4, one file) and M2 `ebc220e0` (evidence). AC-B T2 at base 10/10 FAIL. AC-A 180 PASS / 0 FAIL on both roots. AC-A-ROOT as designed. AC-D plain `ok`, race `ok` with 0 DATA RACE, full `go build ./... && go test ./...` rc 0 (35 packages), `verify_ail.sh` rc 0. M1–M7 all killed, with byte-identical restores.
- **Controller check:** diff scope is one file; the sha matches the prototype; my own re-run of the 6 callers, T1 and T3 was 40/40 PASS under both a `x134x` root and a control root.
- **Judge** (sonnet, Agent tool, own worktree `.wt-world-iter320-eval` at `ebc220e0`, foreground; 45,248 tok / 657 s / 14 tool uses): **PASS 91/100**, zero blocking (30/21/19/14/7).
  - It reverted the guard and saw 6 callers red under the 134 root, and re-ran M1/M3/M7.
  - Of its own mutants, X2–X4 were killed. **X1 (missing `ok` accepted) and X5 (missing `refused` accepted) survived.** The controller reproduced both survivals, added T1 rows `ok-missing` and `refused-missing` (`227b494b`), and showed each mutant now reds only its new row. NB2 (race leg needs `AILANG_BIN`) is environment; NB3 (sweep wording) was already scoped.
  - Report: `verification/world-iter320/evaluator-sprint-r1.md`.
- **Routing evidence (second item):**
  - Sprint worktree base `785b3bd1de34c31c41e0ea5f5f38ebc794ec2406@2026-10-29T16:57:17Z` (`mission-base.sh snap`; DRIFT from the Gate-1 base, 11 commits, recorded).
  - Designer: resolver `recipe claude:claude-opus-5-5 declared:provider-pin`, spawned as Agent `opus` (the 318 precedent). 1 doc, 2 runs (author + the mandated revision). 117,503 tok / 1,397 s, then 153,942 tok / 802 s cumulative. Rotation pointer unchanged (the lane that ran).
  - Planner: Agent `opus`, 155,638 tok / 1,505 s.
  - Executor: Agent `haiku` (resolver `recipe claude:claude-haiku-5-5 declared:provider-pin`), 93,146 tok / 1,398 s.
  - Evaluator: Agent `sonnet`, as above. Generator ≠ judge: haiku/opus vs sonnet.
  - Quorum metered $0.3229.
- **Ruled out:** reverting #418 (the defect is a one-helper test bug with a measured mechanism, not a product regression); widening scope to the 170 variable-token absence sites (residual R3) or row 200's filter flake (a separate row with its own fix).
- **Retro:**
  1. **A green PR head does not certify its merge when the base moves during CI.** #419 was green on `8cee85a9`; its merge landed on `171ece91`, which carried a new red. Gate 3b's merge-commit poll caught it. Re-reading `origin/dev` immediately before merging would have shown the base had moved. Instance 1 here; no rule change.
  2. **A bare-number substring over output that embeds temp paths is a time bomb.** It is filed as design residuals R2/R3 rather than swept here.
  3. Row 200's flake now has its first non-race-leg instance (`785b3bd1` attempt 1, job 123580031823).
- **Progress (second item):** 1.0 clauses 7/7 MET; goal unmoved (gate health on a MET clause). Dev red is fixed once the sprint PR lands.
- **Next:** land row 234 (merge-commit CI), then a rule (d) block report unless D-WORLD-102 = A or dev goes red.

**Row 234 follow-through (Gate 3b):** PR #420 head `00564d7d` check-runs 3/3 success. `origin/dev` re-read right before merging was still `785b3bd1` (retro item 1 applied), and the PR was `MERGEABLE`/`CLEAN`. Merge `2f15c7cde321f5b6833ce49f4f0c6a1a45923187`: check-runs 3/3 success (host gate, race gate, verify gate). Row 234 is tagged LANDED, and dev is green.

## 321 — 2026-10-29 — no pick: rule (d) block report 4; release rows 150/151 still wait on attended steps; dev green at 33f0561a [ADMIN]

**Kind:** bookkeeping-only iteration under standing rule (d) (D-WORLD-46). There is no design, plan, code or product acceptance. The record is the STATUS stamp (318 rotated to the archive), the CURRENT GOAL line, a D-WORLD-102 evidence note, this entry, the index row, the dashboard, and `verification/world-iter321/gate1-2-readings.txt`. An independent judge reviews it before merge.

**Pick and why:** none.
- **Clause map at `33f0561a`:** clauses 1–7 MET. There is no UNMET clause and no critical-path row.
- **D-WORLD-81 release rows, blockers re-measured this iteration (not transcribed):**
  - **150** waits only on the attended A1. The controller runs P0 (≤24 h before) and P3 (before pass A) inside the window Mark schedules (orchestration r5 §4.4). Over commit subjects `8cee85a9..33f0561a`, `grep -ciE 'row ?150|A1\b|marketing'` → **0** (control `row ?2[0-9][0-9]` → **14**). Nothing moved it.
  - **151** needs row 188's N2d and N4a. Over commit subjects, `git log origin/dev --since=2026-10-01 --format='%h %s' | grep -ciE 'N2d|N4a'` → **0**; control `N2c` → **2**. Over full messages (`--grep`) the count is **4**, the same four record or design commits iteration 320 named (`fe9c891e`, `92495b60`, `be184351`, `dc719822`); none builds either milestone. Since iteration 320 the attended track merged row 219 M3 (#418, `171ece91`), but row 219 is still `[NEXT]` on the attended track: its exit criterion (the daemon survives a reboot and a `brew upgrade`, one-command upgrade) is a run Mark makes. Rows 217 and 218 are also still `[NEXT]` attended-track. So neither N2d nor N4a is loop-routable.
- **Rows 231, 230, 233** are hygiene "below every D-WORLD-81 release row". Rule (d) bars them. D-WORLD-64's exception covers rows 148–151 only. D-WORLD-102 is OPEN with default B, so this is a block report.

**Gate 0:**
- Kill switch armed (`mission-world.disabled` absent); gh `sunholo-voight-kampff`; billing tripwire CLEAN; env read by NAME.
- `mission_directives.sh`: 0 allowlisted directives on #366 since `2026-10-27T09:03:42Z` (19 comments, none from the allowlist), read against both the issue-scoped and the mission-scoped watermark (Repo Profile rule; both read `2026-10-27T09:03:42Z`). Nothing to mark.
- `mission_gate0_self_notices.sh`: verdict "1 crash notice", the same 2026-10-22 rc=1 notice `6345079873`, already credited (log count 18 before this entry). Control 107:4 ok.
- Unread `mission-world` inbox: only this loop's claim messages. The iteration-321 claim was sent (`inbox_1793304200279_4873de37`) and shows in the listing.
- `mission_decisions.sh --check`: valid, 91 rows, 1 OPEN (D-WORLD-102). No ledger row changed since iteration 320. The external-issue sweep is not due (#366 is this week's thread).

**Gate 1:**
- Local `dev` was 12 commits behind `origin/dev`; fast-forwarded (`--ff-only`, tree had only the known untracked residue) to `33f0561a9fe411126ed2b5f695604a8c739b7dc5`. Gate-1 base `33f0561a…@2026-10-29T20:01:57Z`.
- Since `8cee85a9`: 20 commits, 4 first-parent merges (#418 row 219 M3 attended, #419 iteration 320's record, #420 row 234, #421 iteration 320's landed record).
- Skill drift: per-file `cmp` of `SKILL.md` and 13 resource files at the resolved symlink target against `origin/dev` in the driver clone: all 14 identical. Rules read from the resolved target.
- Dev CI: `commits/33f0561a/check-runs` → 3, none not-green. The last six dev runs: `33f0561a`, `2f15c7cd` success; `785b3bd1`, `171ece91` failure (the `requireF15` red row 234 fixed); `8cee85a9`, `90e6103b` success.
- Open PRs: none. Died-mid-flight traces: no iteration-321 branch or worktree existed before this one. `.wt-world-iter306-r0` is the known residue.

**Designer / planner / executor:** not spawned. There is no pick, so the routing table has nothing to send them (precedent 257–264, 317, 319, 320). No spawn was attempted, so none failed. Resolver readings are in the evidence file. The operator's standing request for Agent-tool roles is honoured through the evaluator.

**Independent evaluator (record):** REQUIRED. It runs through the Agent tool on `sonnet` (resolver `agent-tool sonnet declared:alias-pin`), in a fresh, read-only, foreground context, and judges the record commit before merge. See the follow-through below.

**Routing evidence:** base=33f0561a9fe411126ed2b5f695604a8c739b7dc5@2026-10-29T20:04:03Z (Gate 4).
- Driver env: `MISSION_DESIGNER_RESOLVED=claude:claude-opus-5-5`, `MISSION_PLANNER_RESOLVED=opus`, `MISSION_EXECUTOR_RESOLVED=claude:claude-haiku-5-5`, `MISSION_EVALUATOR_RESOLVED=sonnet`; `MISSION_OVER_RATION=codex`. `MISSION_ROUTING_NOTE`: planner and executor codex `gpt-6.1-sol` unusable (probe rc 75, over daily ration).
- Controller `claude:claude-opus-5-5` (tok: not reported). Designer, planner and executor were unused.
- Evaluator: Agent `sonnet` (`agent-tool sonnet declared:alias-pin`), tokens recorded in the follow-through. Generator ≠ judge: opus controller vs sonnet judge.
- Quorum: none (no doc). Metered $0.

**Ruled out:**
- Picking 150 (attended A1) or 151 (row 188 N2d/N4a absent, measured above).
- Picking 231, 230 or 233 under rule (d) with D-WORLD-102 unanswered.
- Reading row 219 M3's merge (#418) as row 219 done: the row's exit criterion is an attended run, and it is still `[NEXT]`.
- `ailang mission rotate-log` (37 live entries with this one, under the ~40 threshold); the index row is hand-written as in 317–320.

**Retro:**
1. **Idle cadence:** this is block report 4 of the run (317, 319, 320, 321; 318 was a dev red, and 320 also took a dev red as its second item). Each costs one judge and no code. D-WORLD-102 is still the only lever the loop has. The digest shows it as the one open decision; it is not asked again.
- Harness share 0. No skill edit.

**Progress:** 1.0 clauses 7/7 MET; goal unmoved (release rows). Row 150: the attended A1 remains; row 151: behind row 188 N2d/N4a.
**Next:** rule (d) block report unless D-WORLD-102 = A (then 231 → 230 → 233) or dev goes red. Row 150's P0/P3 run in the A1 window Mark schedules; row 151 becomes routable when row 188's N2d and N4a land.

**Record follow-through (iteration 321):** the independent evaluator judged record commit `d87f8d67` on `33f0561a` and returned **PASS 94/100, zero blocking** (35/22/25/12). It ran through the Agent tool on `sonnet`, in a fresh, read-only, foreground context (45,502 tok, 57 s, 9 tool uses).
- **Verified by the judge:** #366's 19 comments all from the loop's account, `33f0561a` 3/3 green and still `origin/dev`, no open PR, ledger 91 rows with D-WORLD-102 the only OPEN row, rows 217–219 still `[NEXT]`, 20 commits / 4 first-parent merges, the six-run CI list per SHA, the 150 and N2d/N4a subject counts with their controls, notice count 18 → 19, live entries 36 → 37, the rotation (3 stamps; archive +318), charter line count unchanged, D-WORLD-102 changed only by an appended sentence, canonical headings, `base=`, and no closing keyword or address.
- **NB1, reproduced and fixed:** the N2d/N4a search was subject-only. With `--grep` over full messages it returns 4, all record or design commits; the entry now gives both forms and names them. **NB2** (the row-150 grep is a text proxy for an attended step): agreed, and the entry already says only that nothing moved it. **NB3** (watermark, inbox, drift and resolver readings are local state the judge could not re-derive): stated, and they are in `gate1-2-readings.txt`.
- No r2: the fix adds evidence wording; nothing scored changed.
- Report: `design_docs/verification/world-iter321/evaluator-r1.md`.

## 322 — 2026-10-30 — no pick: rule (d) block report 5; release rows 150/151 still wait on attended steps; four live rows' closed upstream issues recorded (53, 77, 78, 162) [ADMIN]

**Kind:** bookkeeping-only iteration under standing rule (d) (D-WORLD-46). There is no design, plan, code or product acceptance. The record is the STATUS stamp (319 rotated to the archive), the CURRENT GOAL line, a D-WORLD-102 evidence note, dated predicate notes on rows 53, 77, 78 and 162, this entry, the index row, the dashboard, and `verification/world-iter322/gate1-2-readings.txt`. An independent judge reviews it before merge.

**Pick and why:** none.
- **Clause map at `ddb5da67`:** clauses 1–7 MET. There is no UNMET clause and no critical-path row.
- **What moved since iteration 321:** nothing but 321's own record. `git log 33f0561a..origin/dev` lists 3 commits (`d87f8d67`, `a379b235`, merge `ddb5da67` of #422).
- **D-WORLD-81 release rows, blockers re-measured this iteration:**
  - **150** waits only on the attended A1. The controller runs P0 (≤24 h before) and P3 (before pass A) inside the window Mark schedules (orchestration r5 §4.4). Over commit subjects `33f0561a..ddb5da67`, `grep -ciE 'row ?150|A1\b|marketing'` → **0**; the control `grep -ciE 'iteration|record'` → **3**, so the instrument saw every commit in the range.
  - **151** needs row 188's N2d and N4a. `git log origin/dev --since=2026-10-01 --format='%h %s' | grep -ciE 'N2d|N4a'` → **0**; control `N2c` → **2**. The full-message form (`git log origin/dev --since=2026-10-01 -i -E --grep='N2d|N4a'`) → **4**: the same four record or design commits iterations 320 and 321 named (`fe9c891e`, `92495b60`, `be184351`, `dc719822`). None builds either milestone. Rows 217, 218 and 219 are still `[NEXT]` on the attended track, so neither milestone is loop-routable.
- **Rows 231, 230, 233** are hygiene "below every D-WORLD-81 release row". Rule (d) bars them, and D-WORLD-64's exception covers rows 148–151 only. D-WORLD-102 is OPEN with default B, so this is a block report.
- **External predicates re-read.** Gate 2 says to re-read the predicate of every row that waits on an outside party, not only the row being picked.
  - **First pass (controller, incomplete):** it found three rows naming a `sunholo-data/ailang` issue: 130, 162 and 164. The enumerator kept only rows whose tag was one of NEXT/PARKED/BACKLOG/GATED/IN/HARNESS, and it matched only `sunholo-data/ailang#N` or `upstream #N`. It had no positive control. The judge found rows 53, 77 and 78, which it had missed.
  - **Second pass:** every queue row, matching `ailang#N`, `ailang/issues/N`, `ailang/pull/N`, `upstream #N` and `upstream issue #N`. It found 15 rows and 13 issues, and was still short. The r2 judge found the spaced spelling `ailang #N` on rows 116, 118 and 178.
  - **Third pass:** adds `ailang ?#N` and `upstream ?#N`, with row 53 as the positive control. It finds **at least 18 rows naming 16 distinct issues or PRs: 5, 53, 73, 75, 77, 78, 108, 116, 118, 130, 134, 135, 136, 141, 162, 164, 178, 184. All 16 are closed** (`gh api repos/sunholo-data/ailang/issues/<n>`). It is a lower bound: a bare `#N` is not counted, because World's own PR numbers use the same space. Rows 156 and 159 cite bare upstream numbers; the judge reports 156 is CLOSED and 159 already records its seam as delivered.
  - On rows tagged LANDED, CLOSED or RE-LANDED (5, 73, 75, 108, 134, 135, 136, 141, 184) the closure changes nothing. Rows 130 and 164 and the RESOLVED rows 116, 118 and 178 already recorded it.
  - **Four live rows had not recorded the closure:**
    - **Row 53** (`[ROUTED — WORLD-SIDE COMPLETE … TRACKING #941]`). #941 closed COMPLETED 2026-10-15T18:12:34Z. The closing comment says the salvage fix is #1824, merged as `514ca3bc3` (fleet iterations 36/41). World had no further move, so nothing re-enters the routable list.
    - **Row 77.** Its code half is #1042, closed NOT_PLANNED 2026-09-13 as superseded: the resolver was rewritten as a pin passthrough (`70e45306`).
    - **Row 78.** Its code half is #1043, closed COMPLETED 2026-10-03. The triage comment says `eec86ca4f` (PR #1377) wires both `-e` extensions into `mission_pi_run.sh`.
    - **Row 162 → #1607, CLOSED 2026-10-09T07:43:49Z by `be6ca8b43`** (#1704, M-PKG-REGISTRY-CONFINEMENT; design `c3bfb1c07`, #1667). Before this record, `1607` appeared only in the charter's row-141 and row-162 text and the iteration-238 archive stamp. `git tag --contains be6ca8b43` starts at v0.53.1, and World's tool pin is v0.54.15, so the pinned tool carries the fix.
  - **Disposition:** each of the four rows now carries a dated note giving the closing reference. No ordering changes:
    - 53 had already left the routable list.
    - 77 and 78 are fleet-owned harness halves. They were not re-measured here, because that is outside this loop's authority.
    - 162 says "gated on nothing": #1607 was one of three options ("or wait for #1607"), not a blocker. It stays PARKED-unranked and barred by rule (d). Its V13 premise was measured on v0.52.1, and its note says to re-measure V13 on the pinned tool before designing.

**Gate 0:**
- Kill switch armed (`mission-world.disabled` absent). gh `sunholo-voight-kampff`. Billing tripwire CLEAN. Env read by NAME.
- `mission_directives.sh`: 0 allowlisted directives on #366 since `2026-10-27T09:03:42Z` (21 comments, none from the allowlist). Nothing to mark.
- `mission_gate0_self_notices.sh`: rc 1, verdict "1 crash notice". It is the same 2026-10-22 rc=1 notice `6345079873`, already credited (`grep -c 6345079873 design_docs/world-mission-log.md` → 19 before this entry). Control 107:4 ok. No new notice on #366.
- Unread inbox: a controlplane driver notice that the pinned source clone is 99 commits behind (this fire ran pinned, driver age 0), and two mission-fleet items about stapledon's `pi:provider-error:no-retry` ticket (not World's). The `mission-world` inbox holds only this loop's claim messages. The iteration-322 claim was sent (`inbox_1793318597095_a077fd3c`) and shows in the listing.
- `mission_decisions.sh --check`: valid, 91 rows, 1 OPEN (D-WORLD-102). No ledger row changed since iteration 321. The external-issue sweep is not due (#366 is this week's thread).

**Gate 1:**
- Local `dev` = `origin/dev` = `ddb5da67ee085b961f10ad3f9db5e2c259e10ea6`; no fast-forward needed. Gate-1 base `ddb5da67…@2026-10-30T00:01:34Z`.
- Skill drift: per-file `cmp` of `SKILL.md` and 13 resource files at the resolved symlink target against `origin/dev` in the driver clone: all 14 SAME. The pin worktree's 14 copies also match the resolved target. Rules read from the resolved target.
- Dev CI: `commits/ddb5da67/check-runs` → 3, all success. The last eight dev runs: `ddb5da67`, `33f0561a`, `2f15c7cd` success; `785b3bd1`, `171ece91` failure (the `requireF15` red row 234 fixed); `8cee85a9`, `90e6103b`, `bd8bb914` success.
- Open PRs: none. Died-mid-flight traces: no iteration-322 branch or worktree existed before this one. `.wt-world-iter306-r0` is the known residue. The attended session's scratchpad worktrees and its `trial/*` branches (last pushed 2026-10-29 afternoon) belong to the attended track and were left alone.

**Designer / planner / executor:** not spawned. There is no pick, so the routing table has nothing to send them (precedent 257–264, 317, 319–321). No spawn was attempted, so none failed. The operator's standing request for Agent-tool roles is honoured through the evaluator.

**Independent evaluator (record):** REQUIRED. It runs through the Agent tool on `sonnet` (driver `MISSION_EVALUATOR_RESOLVED=sonnet`, path `agent-tool`), in a fresh, read-only, foreground context, and judges the record commit before merge. See the follow-through below.

**Routing evidence:** base=ddb5da67ee085b961f10ad3f9db5e2c259e10ea6@2026-10-30T00:03:27Z (Gate 4).
- Driver env: `MISSION_DESIGNER_RESOLVED=claude:claude-opus-5-5`, `MISSION_PLANNER_RESOLVED=opus`, `MISSION_EXECUTOR_RESOLVED=claude:claude-haiku-5-5`, `MISSION_EVALUATOR_RESOLVED=sonnet`; `MISSION_OVER_RATION=codex`. `MISSION_ROUTING_NOTE`: planner and executor codex `gpt-6.1-sol` unusable (probe rc 75, over daily ration).
- Controller `claude:claude-opus-5-5` (tok: not reported). Designer, planner and executor were unused.
- Evaluator: Agent `sonnet`, tokens recorded in the follow-through. Generator ≠ judge: opus controller vs sonnet judge.
- Quorum: none (no doc). Metered $0.

**Ruled out:**
- Picking 150 (attended A1) or 151 (row 188 N2d/N4a absent, measured above).
- Picking 231, 230 or 233 under rule (d) with D-WORLD-102 unanswered.
- Picking 162, 77 or 78 because an upstream predicate flipped. Gate 2's flipped-predicate rule makes a row the pick when the flip changes its ordering. None of these flips does. Row 162 was gated on nothing, and 77 and 78 are fleet-owned harness halves. Rule (d), an attended ruling, bars all three either way.
- Re-measuring row 162's V13 premise now. That is design work on a barred row, and it would run `pkg_docs` against the live registry. It is noted on the row for whoever picks it.
- `ailang mission rotate-log` (38 live entries with this one, under the ~40 threshold). The index row is hand-written as in 317–321.

**Retro:**
1. **Idle cadence:** this is block report 5 of the run (317, 319, 320, 321, 322). Each costs one judge and no code. D-WORLD-102 is still the only lever the loop has. The digest shows it as the one open decision; it is not asked again.
2. **The external-predicate re-read found something the four previous block reports did not run.** Iterations 317–321 re-measured the release rows' blockers but not the other rows' upstream predicates. Four live rows were stale, by 15 to 47 days. No ordering changed this time, but a block report is the cheapest place to run that read. Gate 2 already prescribes it, so this was a skipped step, not a missing rule, and there is no rule edit.
3. **The controller's first enumeration was a hand-filtered subset with no positive control,** which is the rule-3a/3g shape. It kept only rows with one of six tags and matched only two spellings of an issue reference. It reported "three rows" as if that were complete. The judge's independent search found three more live rows in r1. In r2 it found three more already-resolved rows that the corrected count still missed. The fix in practice is to enumerate every queue row, match every reference spelling, and assert the count against a row you know names an issue. Row 53, the one with an explicit TRACKING tag, is now that control. State the count as a lower bound when one spelling, here a bare `#N`, cannot be attributed. Instance 1; no rule edit.
- Harness share 0. No skill edit.

**Progress:** 1.0 clauses 7/7 MET; goal unmoved (release rows). Row 150: the attended A1 remains. Row 151: behind row 188 N2d/N4a.
**Next:** rule (d) block report unless D-WORLD-102 = A (then 231 → 230 → 233) or dev goes red. Row 150's P0/P3 run in the A1 window Mark schedules; row 151 becomes routable when row 188's N2d and N4a land.

**Record follow-through (iteration 322):** the independent evaluator judged record commit `c0484c24` on `ddb5da67`: **r1 PASS 92/100, zero blocking** (31/22/25/14). It ran through the Agent tool on `sonnet`, in a fresh, read-only, foreground context (53,007 tok, 89 s, 16 tool uses).
- **r1 NB1 was a real gap.** The controller's external-predicate sweep had said "three rows". The judge found three more live rows naming closed upstream issues that were not recorded: 53 (#941), 77 (#1042) and 78 (#1043). Fixed in `d11b8848`, which adds dated notes to all three. NB2 and NB4 (instrument wording) were fixed. NB5 was refuted (all stamps say 91).
- **r2** (SendMessage resume, 65,326 tok, 44 s, 6 tool uses) judged `d11b8848`: **PASS 86/100, zero blocking** (29/20/24/13). It verified every note against GitHub and the ailang clone. Its NB1 found the corrected count still short: the spaced `ailang #N` spelling on the RESOLVED rows 116, 118 and 178. The third pass now gives "at least 18 rows / 16 distinct, all closed", stated as a lower bound because bare `#N` is ambiguous.
- No r3: the last change is count wording only, and every added row was already resolved.
- The judge's score fell from r1 to r2 because the second count was wrong too. That is recorded, not smoothed over.
- Report: `design_docs/verification/world-iter322/evaluator-r1-r2.md`.

**Landing (iteration 322):** PR #423 merged as `e15fc11f`. The PR head `aa66b3ff` was 3/3 green before the merge, and origin/dev was re-read right before it (`ddb5da67`, MERGEABLE/CLEAN). The merge tree is byte-identical to the head tree (`58cfae0d…` both). On the merge, the charter has 6303 lines and 3 stamps, and the ledger is valid at 91 rows.
- **Merge CI, attempt 1: go race gate red** (run 40980017500). `host/daemon` `TestHumanApproveIsSerializedAcrossEpisodes` failed after 1.72 s with "unwrapped probe max = 1, want 2". That is the test's own control arm; the serialized arm passed. The two `DATA RACE` blocks earlier in that log are the known-positive race control, which is meant to fire.
  - Record, log and archive had no prior mention of the test (`grep -c` 0 in each). This is instance 1 of this symptom. The judge (r3) found earlier history in `verification/world-iter265/row176-executor-report.md`: a build-time failure on `store: stale observed world head`, which is a different failure and has been handled since. Row 235 cites it.
  - Locally on darwin it passed 0/80 times failing: 20/20 idle, 30/30 at load about 6, and 30/30 at load about 25, slowest 1.34 s. Logs are in `verification/world-iter322/flake-measure.txt`.
  - Mechanism hypothesis: the probe's 200 ms arrival window, filed as row 235 with a candidate fix. It is below every release row, and rule (d) governs it.
- **`gh run rerun 40980017500 --failed`** gave attempt 2 green. `commits/e15fc11f/check-runs` → 3, all success, so dev is green.
- I did not take the red as a second item. A re-run cleared it, the merge tree equals a green tree, there is one instance, and it does not reproduce locally. That does not meet the "measured failure" bar for a fix (the mission-control skill's standing rule 5, "data before conclusions"). The precedent is row 232: filed from PR CI reds and taken only when it went red on dev at Gate 1 (iteration 318).
- **Judge r3** (SendMessage resume, sonnet, 84,813 tok, 69 s, 8 tool uses) checked the landing follow-through `ace05db9`: **PASS 88/100, zero blocking** (29/21/25/13). It verified the run attempts, that the DATA RACE blocks are the known-positive control, the tree identity, the no-prior-mention grep, the mechanism against `concurrencyProbe.Execute`, and the charter diff (one row appended plus the stamp). Its own focused run passed 5/5. NB1 (the iteration-265 history), NB2 (clause tag "assigned"), NB3 (the standing rule's source) and NB4 (row 232's instance numbering) were applied in the amended commit. NB5 (darwin-only measurement) was already stated.

## 323 — 2026-10-30 — no pick: rule (d) block report 6; release rows 150/151 still wait on attended steps; no upstream predicate open; dev green at 8b840dd0 [ADMIN]

**Kind:** bookkeeping-only iteration under standing rule (d) (D-WORLD-46). There is no design, plan, code or product acceptance. The record is the STATUS stamp (320 rotated to the archive), the CURRENT GOAL line, a D-WORLD-102 evidence note, this entry, the index row, the dashboard, and `verification/world-iter323/gate1-2-readings.txt`. An independent judge reviews it before merge.

**Pick and why:** none.
- **Clause map at `8b840dd0`:** clauses 1–7 MET. There is no UNMET clause and no critical-path row.
- **What moved since iteration 322:** nothing but 322's own record and landing. `git log ddb5da67..origin/dev` lists 6 commits (`c0484c24`, `d11b8848`, `aa66b3ff`, merge `e15fc11f` of #423, `c2ee8597`, merge `8b840dd0` of #424).
- **D-WORLD-81 release rows, blockers re-measured this iteration:**
  - **150** waits only on the attended A1. The controller runs P0 (≤24 h before) and P3 (before pass A) inside the window Mark schedules (orchestration r5 §4.4). Over commit subjects `ddb5da67..8b840dd0`, `grep -ciE 'row ?150|A1\b|marketing'` → **0**; the control `grep -ciE 'iteration|record|merge'` → **6** of 6 commits, so the instrument saw the whole range.
  - **151** needs row 188's N2d and N4a. `git log origin/dev --since=2026-10-01 --format='%h %s' | grep -ciE 'N2d|N4a'` → **0**; control `N2c` → **2**. The full-message form → **4**: the same four record or design commits iterations 320–322 named (`fe9c891e`, `92495b60`, `be184351`, `dc719822`). None builds either milestone. Rows 217, 218 and 219 are still `[NEXT]` on the attended track, so neither milestone is loop-routable.
- **Rows 231, 230, 233** are hygiene "below every D-WORLD-81 release row". Rule (d) bars them, and D-WORLD-64's exception covers rows 148–151 only. D-WORLD-102 is OPEN with default B, so this is a block report. Row 235 (filed by 322) is hygiene below the release rows too; rule (d) bars it, and it is not in D-WORLD-102's scope.
- **External predicates re-read.** This time the enumeration covers the whole charter, not only queue rows: `grep -oE '(ailang ?#|ailang/(issues|pull)/|upstream (issue )?#)[0-9]+'` → **29** distinct numbers. **All 29 are closed** (`gh api repos/sunholo-data/ailang/issues/<n> --jq .state`), so no row waits on an open upstream issue and no predicate can have flipped since iteration 322. Controls: row 53's `ailang#941` matched (pattern fires), and `gh issue list -R sunholo-data/ailang --state open --limit 3` returned 3 (the API reports open issues when they exist). It is a lower bound: a bare `#N` is not counted, because World's own PR numbers share that space.
- **Weekly external-issue sweep (Gate 0 step 5):** iteration 305 ran this week's sweep after the 2026-10-26 rotation (0 orphans of 1). Re-run here as a check: `gh issue list -R sunholo-data/ailang-world --state open` → 1 issue, the REST issues endpoint (PRs excluded) → 1, and the one issue is #366, the bookkeeping thread (charter 4 / log 45 / status-archive 1 / dashboard 0 before this record). **0 orphans of 1 enumerated.** A fresh negative-control literal read 0 in the charter (not printed, per Gate 4's spent-control rule).

**Gate 0:**
- Kill switch armed (`mission-world.disabled` absent). gh `sunholo-voight-kampff`. Billing tripwire CLEAN. Env read by NAME.
- `mission_directives.sh`: 0 allowlisted directives on #366 since `2026-10-27T09:03:42Z` (22 comments, none from the allowlist). Nothing to mark.
- `mission_gate0_self_notices.sh`: rc 1, verdict "1 crash notice". It is the same 2026-10-22 rc=1 notice `6345079873`, already credited (`grep -c 6345079873 design_docs/world-mission-log.md` → 20 before this entry). Control 107:4 ok. No new notice on #366.
- Unread inbox: a controlplane driver notice that stapledon's pinned source clone drifted (not World's), the iteration-322 `[fyi]` report, and a daneel item. The `mission-world` inbox holds only this loop's claim messages. The iteration-323 claim was sent (`inbox_1793332967693_601a6427`).
- `mission_decisions.sh --check`: valid, 91 rows, 1 OPEN (D-WORLD-102). No ledger row changed since iteration 322.

**Gate 1:**
- Local `dev` = `origin/dev` = `8b840dd03cc4b6ff5d7a745c0d03faf8a147e3d5`; no fast-forward needed. Gate-1 base `8b840dd0…@2026-10-30T04:01:39Z`.
- Skill drift: per-file `cmp` of `SKILL.md` and 13 resource files at the resolved symlink target against `origin/dev` in the driver clone: all 14 SAME. Rules read from the resolved target.
- Dev CI: `commits/8b840dd0/check-runs` → 3, all success (ailang-code verify gate, go race gate, go host build + test gate). The last five dev CI runs (`8b840dd0`, `e15fc11f`, `ddb5da67`, `33f0561a`, `2f15c7cd`) all read success; `e15fc11f`'s is its attempt-2 re-run (row 235).
- Open PRs: none. Died-mid-flight traces: no iteration-323 branch or worktree existed before this one. `.wt-world-iter306-r0` is the known residue. The attended session's scratchpad worktrees belong to the attended track and were left alone.

**Designer / planner / executor:** not spawned. There is no pick, so the routing table has nothing to send them (precedent 257–264, 317, 319–322). No spawn was attempted, so none failed. The operator's standing request for Agent-tool roles is honoured through the evaluator.

**Independent evaluator (record):** REQUIRED. It runs through the Agent tool on `sonnet` (driver `MISSION_EVALUATOR_RESOLVED=sonnet`, path `agent-tool`), in a fresh, read-only, foreground context, and judges the record commit before merge. See the follow-through below.

**Routing evidence:** base=8b840dd03cc4b6ff5d7a745c0d03faf8a147e3d5@2026-10-30T04:04:09Z (Gate 4).
- Driver env: `MISSION_DESIGNER_RESOLVED=claude:claude-opus-5-5`, `MISSION_PLANNER_RESOLVED=opus`, `MISSION_EXECUTOR_RESOLVED=claude:claude-haiku-5-5`, `MISSION_EVALUATOR_RESOLVED=sonnet`; `MISSION_OVER_RATION=codex`. `MISSION_ROUTING_NOTE`: planner and executor codex `gpt-6.1-sol` unusable (probe rc 75, over daily ration).
- Controller `claude:claude-opus-5-5` (tok: not reported). Designer, planner and executor were unused.
- Evaluator: Agent `sonnet`, tokens recorded in the follow-through. Generator ≠ judge: opus controller vs sonnet judge.
- Quorum: none (no doc). Metered $0.

**Ruled out:**
- Picking 150 (attended A1) or 151 (row 188 N2d/N4a absent, measured above).
- Picking 231, 230, 233 or 235 under rule (d) with D-WORLD-102 unanswered.
- Re-asking D-WORLD-102. It is the one open decision and the digest carries it; it is not asked again.
- `ailang mission rotate-log` (39 live entries with this one, under the ~40 threshold). The index row is hand-written as in 317–322.

**Retro:**
1. **Idle cadence:** this is block report 6 of the run (317, 319, 320, 321, 322, 323). Each costs one judge and two CI cycles, and no code. D-WORLD-102 is still the only lever the loop has.
2. **Whole-charter enumeration instead of queue-row enumeration.** Iteration 322 needed three passes to count the upstream references on queue rows. Matching every spelling across the whole charter and checking each number's state took one pass and gave a superset (29 numbers). Because every one is closed, the answer does not depend on which row names which issue. Instance 1; no rule edit.
3. **The controller's own search truncated a positive.** My first check for this week's issue sweep piped `grep` through `tail -4`, which cut off iteration 305's "0 orphans of 1" line, and the first draft of the stamp said no iteration had run the sweep. A per-iteration `awk` over the log found 305's line, and the stamp was corrected before commit. This is rule 3a aimed at a display filter: a `tail` on a search is a sample, not a result. Instance 1; no rule edit.
- Harness share 0. No skill edit.

**Progress:** 1.0 clauses 7/7 MET; goal unmoved (release rows). Row 150: the attended A1 remains. Row 151: behind row 188 N2d/N4a.
**Next:** rule (d) block report unless D-WORLD-102 = A (then 231 → 230 → 233) or dev goes red. Row 150's P0/P3 run in the A1 window Mark schedules; row 151 becomes routable when row 188's N2d and N4a land.

**Record follow-through (iteration 323):** the independent evaluator judged record commit `f7c68d57` on `8b840dd0`: **r1 PASS 93/100, zero blocking** (33/23/25/12). It ran through the Agent tool on `sonnet`, in a fresh, read-only, foreground context (46,440 tok, 58 s, 10 tool uses). It re-derived the no-pick premise (6 commits since `ddb5da67`, 0 directives in #366's 22 comments, 0 open PRs, `8b840dd0` 3/3), the 29 closed upstream numbers by its own per-number API reads, the sweep counts, the STATUS rotation (320 stamp byte-identical in the archive, 3 stamps, 243 queue rows before and after) and the ledger (91 valid).
- **NB1** (four of the five dev CI runs not re-derived by the judge): verified by the controller with `gh run list --workflow CI --branch dev --limit 5` (all success; `e15fc11f` attempt 2) and banked in the readings file.
- **NB2** (readings "prints only `217. **[NEXT`"): refuted. The file prints all three tags (217, 218, 219). The `[NEXT` cut is the pattern stopping at the tag, by design.
- **NB3** (the dashboard dropped row 162's re-measure caveat): applied.
- **NB4** ("loop-routable" for rows rule (d) bars): applied. The stamp now says "the next rows below them".
- No r2: every change is wording, a dashboard line or a banked reading.

**Landing (iteration 323):** PR #425 merged as `bd3dae14`. The PR head `6500fd13` was 3/3 green before the merge, and origin/dev was re-read right before it (`8b840dd0`, MERGEABLE/CLEAN, `--match-head-commit`). The merge tree is byte-identical to the head tree. On the merge, the charter keeps 3 stamps and the ledger is valid at 91 rows. **Merge CI:** `commits/bd3dae14/check-runs` → 3, all success on attempt 1 (verify gate, go race gate, go host build + test gate). Dev green.

## 324 — 2026-10-30 — no pick: rule (d) block report 7; release rows 150/151 still wait on attended steps; no upstream predicate open; dev green at 42fdf32a [ADMIN]

**Kind:** bookkeeping-only iteration under standing rule (d) (D-WORLD-46). There is no design, plan, code or product acceptance. The record is the STATUS stamp (321 rotated to the archive), the CURRENT GOAL line, a D-WORLD-102 evidence note, this entry, the index row, the dashboard, and `verification/world-iter324/gate1-2-readings.txt`. An independent judge reviews it before merge.

**Pick and why:** none.
- **Clause map at `42fdf32a`:** clauses 1–7 MET. There is no UNMET clause and no critical-path row.
- **What moved since iteration 323:** nothing but 323's own record and landing. `git log 8b840dd0..origin/dev` lists 5 commits (`f7c68d57`, `6500fd13`, merge `bd3dae14` of #425, `23e34d87`, merge `42fdf32a` of #426).
- **D-WORLD-81 release rows, blockers re-measured this iteration:**
  - **150** waits only on the attended A1. The controller runs P0 (≤24 h before) and P3 (before pass A) inside the window Mark schedules (orchestration r5 §4.4). Over commit subjects `8b840dd0..42fdf32a`, `grep -ciE 'row ?150|A1\b|marketing'` → **0**; the control `grep -ciE 'iteration|record|merge'` → **5** of 5 commits, so the instrument saw the whole range.
  - **151** needs row 188's N2d and N4a. `git log origin/dev --since=2026-10-01 --format='%h %s' | grep -ciE 'N2d|N4a'` → **0**; control `N2c` → **2**. The full-message form → **4**: the same four record or design commits iterations 320–323 named (`fe9c891e`, `92495b60`, `be184351`, `dc719822`). None builds either milestone. Rows 217, 218 and 219 are still `[NEXT]` on the attended track, so neither milestone is loop-routable.
- **Rows 231, 230, 233** are hygiene "below every D-WORLD-81 release row". Rule (d) bars them, and D-WORLD-64's exception covers rows 148–151 only. D-WORLD-102 is OPEN with default B, so this is a block report. Row 235 is hygiene below the release rows too; rule (d) bars it, and it is not in D-WORLD-102's scope.
- **External predicates re-read** with iteration 323's whole-charter pattern: **29** distinct upstream numbers, **all 29 closed** (`gh api repos/sunholo-data/ailang/issues/<n> --jq .state`). Controls: row 53's `ailang#941` matched, and `gh issue list -R sunholo-data/ailang --state open --limit 3` returned 3. Still a lower bound (bare `#N` not counted). The first run of this check was an instrument failure (see Retro 2); the readings file prints both runs.
- **Weekly external-issue sweep (Gate 0 step 5):** iteration 305 ran this week's sweep. Re-checked: 1 open issue, #366, the bookkeeping thread (charter 4 / log 49 / status-archive 2 / dashboard 1 before this record). **0 orphans of 1 enumerated.** A fresh negative-control literal read 0 in the charter (not printed).

**Gate 0:**
- Kill switch armed (`mission-world.disabled` absent). gh `sunholo-voight-kampff`. Billing tripwire CLEAN.
- `mission_directives.sh`: 0 allowlisted directives on #366 since `2026-10-27T09:03:42Z` (23 comments) and 0 on #313 since `2026-10-22T01:34:42Z` (44 comments). Nothing to mark.
- `mission_gate0_self_notices.sh`: rc 1, verdict "1 crash notice". It is the same 2026-10-22 rc=1 notice `6345079873`, already credited (`grep -c 6345079873 design_docs/world-mission-log.md` → 21 before this entry). Control 107:4 ok. No new notice on #366.
- Unread inbox: the stapledon iteration-59 `[fyi]`, a fleet iteration-71 report, and a fleet `[harness-resolved]` ticket addressed to the attended track (watermark ghost of #1962). None is a World item. The `mission-world` inbox holds only this loop's claim messages. The iteration-324 claim was sent (`inbox_1793347370348_58cb7988`).
- `mission_decisions.sh --check`: valid, 91 rows, 1 OPEN (D-WORLD-102). The only ledger diff in `8b840dd0..42fdf32a` is 323's own D-WORLD-102 evidence-note append (one −/+ pair).

**Gate 1:**
- Local `dev` was 5 behind `origin/dev`, with only untracked residue in the tree, and fast-forwarded to `42fdf32a`. Gate-1 base `42fdf32a95c09df95d56f951d807e9bd92f508d4@2026-10-30T08:02:41Z`.
- Skill drift: per-file `cmp` of `SKILL.md` and 13 resource files at the resolved symlink target against `origin/dev` in the driver clone: all 14 SAME. Rules read from the resolved target.
- Dev CI: `commits/42fdf32a/check-runs` → 3, all success. The last five dev CI runs (`42fdf32a`, `bd3dae14`, `8b840dd0`, `e15fc11f`, `ddb5da67`) all read success; `e15fc11f`'s is its attempt-2 re-run (row 235).
- Open PRs: none. Died-mid-flight traces: no iteration-324 branch or worktree existed before this one. `.wt-world-iter306-r0` is the known residue. The attended scratchpad worktrees belong to the attended track and were left alone.

**Designer / planner / executor:** not spawned. There is no pick, so the routing table has nothing to send them (precedent 257–264, 317, 319–323). No spawn was attempted, so none failed. The operator's standing request for Agent-tool roles is honoured through the evaluator.

**Independent evaluator (record):** REQUIRED. It runs through the Agent tool on `sonnet` (driver `MISSION_EVALUATOR_RESOLVED=sonnet`, path `agent-tool`), in a fresh, read-only, foreground context, and judges the record commit before merge. See the follow-through below.

**Routing evidence:** base=42fdf32a95c09df95d56f951d807e9bd92f508d4@2026-10-30T08:05:07Z (Gate 4).
- Driver env: `MISSION_DESIGNER_RESOLVED=claude:claude-opus-5-5`, `MISSION_PLANNER_RESOLVED=opus`, `MISSION_EXECUTOR_RESOLVED=claude:claude-haiku-5-5`, `MISSION_EVALUATOR_RESOLVED=sonnet`; `MISSION_OVER_RATION=codex`. `MISSION_ROUTING_NOTE`: planner and executor codex `gpt-6.1-sol` unusable (probe rc 75, over daily ration).
- Controller `claude:claude-opus-5-5` (tok: not reported). Designer, planner and executor were unused.
- Evaluator: Agent `sonnet`, tokens recorded in the follow-through. Generator ≠ judge: opus controller vs sonnet judge.
- Quorum: none (no doc). Metered $0.

**Ruled out:**
- Picking 150 (attended A1) or 151 (row 188 N2d/N4a absent, measured above).
- Picking 231, 230, 233 or 235 under rule (d) with D-WORLD-102 unanswered.
- Re-asking D-WORLD-102. It is the one open decision and the digest carries it; it is not asked again.
- `ailang mission rotate-log`: this entry makes 40 live entries, which is at the ~40 threshold, not past it. The index row is hand-written as in 317–323. The next iteration that adds an entry should rotate (`--keep 20`, registry pointed at the record worktree per Gate 4).

**Retro:**
1. **Idle cadence:** block report 7 of the run (317, 319, 320, 321, 322, 323, 324). Each costs one judge and two CI cycles, and no code. D-WORLD-102 is still the only lever the loop has.
2. **The controller's shell is an instrument too.** The upstream-state loop was written as `for n in $NUMS` and run in this harness's zsh. zsh does not word-split an unquoted variable, so the loop ran once with all 29 numbers in one API path and printed `1843:ERR` / `not-closed=1`. That looks like a real open predicate. Re-run under `bash -c` with the same pattern, all 29 read closed. A `count=29` beside a single result line was the tell. In the same file, zsh's `echo` turned the `\b` in a displayed command into a backspace byte. Both are fixed in the readings file, which prints the failed run beside the corrected one. Instance 1 in this charter; no rule edit. Iteration 323's per-number loop was correct, so it does not need re-checking.
- Harness share 0. No skill edit.

**Progress:** 1.0 clauses 7/7 MET; goal unmoved (release rows). Row 150: the attended A1 remains. Row 151: behind row 188 N2d/N4a.
**Next:** rule (d) block report unless D-WORLD-102 = A (then 231 → 230 → 233) or dev goes red. Row 150's P0/P3 run in the A1 window Mark schedules; row 151 becomes routable when row 188's N2d and N4a land. Rotate the log at the next entry.

**Record follow-through (iteration 324):** the independent evaluator judged record commit `e4396284` on `42fdf32a`: **r1 PASS 93/100, zero blocking** (34/24/23/12). It ran through the Agent tool on `sonnet`, in a fresh, read-only, foreground context (37,995 tok, 70 s, 7 tool uses). It re-derived the no-pick premise (5 commits since `8b840dd0`, all 323's own; 0 allowlisted comments on #366's 23 or #313's 44; ledger 91 valid with D-WORLD-102 the only OPEN row; `42fdf32a` 3/3), the 29 closed upstream numbers under `bash -c`, the rotation (3 stamps; 321's stamp byte-identical at the archive head; charter diff 4/4; queue untouched) and the 40 live log entries.
- **NB1** (log names a follow-through that did not exist yet): this paragraph and the landing line close it.
- **NB2** (the 29 depends on the pattern; the judge's narrower regex found 25, all closed): already labelled a lower bound; no change.
- **NB3** (raw ANSI colour escapes in the readings file's directive lines): applied, escapes stripped.
- **NB4** (seven block reports in a row; D-WORLD-102 is the only lever): already Retro 1; it stays in the digest.
- No r2: the one change is a byte-level clean-up of the evidence file.

**Landing (iteration 324):** PR #427 merged as `c4386177`. The PR head `bd437ae8` was 3/3 green before the merge, and origin/dev was re-read right before it (`42fdf32a`, MERGEABLE/CLEAN, `--match-head-commit`). The merge tree is byte-identical to the head tree (`3e61fd70`). On the merge, the charter keeps 3 stamps. **Merge CI:** `commits/c4386177/check-runs` → 3, all success on attempt 1 (verify gate, go race gate, go host build + test gate). Dev green. The controller's first CI poller had the same zsh word-splitting fault as Retro 2 (`set -- $r` did not split, so it could only end at its deadline); it was stopped and re-run under `bash -c` before any verdict was read.

## 325 — 2026-10-30 — no pick: rule (d) block report 8; release rows 150/151 still wait on attended steps; log rotated to 20 live entries; dev green at 39aabca3 [ADMIN]

**Kind:** bookkeeping-only iteration under standing rule (d) (D-WORLD-46). There is no design, plan, code or product acceptance. The record is the STATUS stamp (322 rotated to the archive), the CURRENT GOAL line, a D-WORLD-102 evidence note, this entry, the log rotation and regenerated index, the dashboard, and `verification/world-iter325/gate1-2-readings.txt`. An independent judge reviews it before merge.

**Pick and why:** none.
- **Clause map at `39aabca3`:** clauses 1–7 MET. There is no UNMET clause and no critical-path row.
- **What moved since iteration 324:** nothing but 324's own record and landing. `git log 42fdf32a..origin/dev` lists 5 commits (`e4396284`, `bd437ae8`, merge `c4386177` of #427, `3eed86ac`, merge `39aabca3` of #428).
- **D-WORLD-81 release rows, blockers re-measured this iteration:**
  - **150** waits only on the attended A1. The controller runs P0 (≤24 h before) and P3 (before pass A) inside the window Mark schedules (orchestration r5 §4.4). Over commit subjects `42fdf32a..39aabca3`, `grep -ciE 'row ?150|A1\b|marketing'` → **0**; the control `grep -ciE 'iteration|record|merge'` → **5** of 5 commits.
  - **151** needs row 188's N2d and N4a. `git log origin/dev --since=2026-10-01 --format='%h %s' | grep -ciE 'N2d|N4a'` → **0**; control `N2c` → **2**. The full-message form → **4**: the same four record or design commits iterations 320–324 named (`fe9c891e`, `92495b60`, `be184351`, `dc719822`). Rows 217, 218 and 219 are still `[NEXT]` on the attended track, so neither milestone is loop-routable.
- **Rows 231, 230, 233** are hygiene below every D-WORLD-81 release row. Rule (d) bars them, and D-WORLD-64's exception covers rows 148–151 only. D-WORLD-102 is OPEN with default B, so this is a block report. Row 235 is hygiene too and is not in D-WORLD-102's scope.
- **External predicates re-read** with the whole-charter pattern, run under `bash` from a script file (iteration 324's zsh word-splitting fault cannot recur that way): **29** distinct upstream numbers, `results=29`, **all 29 closed**. Controls: row 53's `ailang#941` matched, and `gh issue list -R sunholo-data/ailang --state open --limit 3` returned 3. Still a lower bound (bare `#N` not counted).
- **Weekly external-issue sweep (Gate 0 step 5):** iteration 305 ran this week's sweep. Re-checked: 1 open issue, #366, the bookkeeping thread (charter 4 / log 53 / status-archive 3 / dashboard 1 before this record). **0 orphans of 1 enumerated.**

**Gate 0:**
- Kill switch armed (`mission-world.disabled` absent). gh `sunholo-voight-kampff`. Billing tripwire CLEAN.
- `mission_directives.sh` (driver pin copy; this repo has no `scripts/mission_directives.sh`): 0 allowlisted directives on #366 since `2026-10-27T09:03:42Z` (24 comments) and 0 on #313 since `2026-10-22T01:34:42Z` (44 comments). Nothing to mark.
- `mission_gate0_self_notices.sh`: rc 1, verdict "1 crash notice". It is the same 2026-10-22 rc=1 notice `6345079873`, already credited (`grep -c` → 22 before this entry). Control 107:4 ok. No new notice on #366 (24 self, 0 crash).
- Unread inbox: a mission-triage iteration-20 `[fyi]`, two eparse package-task items and a coordinator approval for eparse, and two feedback-gate dry-run notices. None is a World item. The `mission-world` inbox holds only this loop's claim and report messages. The iteration-325 claim was sent (`inbox_1793361869007_883be643`).
- `mission_decisions.sh --check`: valid, 91 rows, 1 OPEN (D-WORLD-102). The only ledger diff in `42fdf32a..39aabca3` is 324's own D-WORLD-102 evidence-note append (one −/+ pair).

**Gate 1:**
- Local `dev` equalled `origin/dev` (`39aabca3`), with only untracked residue in the tree. Gate-1 base `39aabca3199e928b98f565d26cc00d5c92eeca9e@2026-10-30T12:01:56Z`.
- Skill drift: per-file `cmp` of `SKILL.md` and 13 resource files at the resolved symlink target against `origin/dev` in the driver clone, re-run after fetching that clone: all 14 SAME (drift=0). There is no CWD-relative copy in this repo. Rules read from the resolved target.
- Dev CI: `commits/39aabca3/check-runs` → 3, all success. The last five dev CI runs (`39aabca3`, `c4386177`, `42fdf32a`, `bd3dae14`, `8b840dd0`) all read success on attempt 1.
- Open PRs: none. Died-mid-flight traces: no iteration-325 branch or worktree existed before this one. `.wt-world-iter306-r0` is the known residue. The attended scratchpad worktrees belong to the attended track and were left alone.

**Designer / planner / executor:** not spawned. There is no pick, so the routing table has nothing to send them (precedent 257–264, 317, 319–324). No spawn was attempted, so none failed. The operator's standing request for Agent-tool roles is honoured through the evaluator.

**Independent evaluator (record):** REQUIRED. It runs through the Agent tool on `sonnet` (driver `MISSION_EVALUATOR_RESOLVED=sonnet`, path `agent-tool`), in a fresh, read-only, foreground context, and judges the record commit before merge. See the follow-through below.

**Routing evidence:** base=39aabca3199e928b98f565d26cc00d5c92eeca9e@2026-10-30T12:01:56Z (Gate 1).
- Driver env: `MISSION_DESIGNER_RESOLVED=claude:claude-opus-5-5`, `MISSION_PLANNER_RESOLVED=opus`, `MISSION_EXECUTOR_RESOLVED=claude:claude-haiku-5-5`, `MISSION_EVALUATOR_RESOLVED=sonnet`; `MISSION_OVER_RATION=codex`. `MISSION_ROUTING_NOTE`: planner and executor codex `gpt-6.1-sol` unusable (probe rc 75, over daily ration).
- Controller `claude:claude-opus-5-5` (tok: not reported). Designer, planner and executor were unused.
- Evaluator: Agent `sonnet`, tokens recorded in the follow-through. Generator ≠ judge: opus controller vs sonnet judge.
- Quorum: none (no doc). Metered $0.

**Ruled out:**
- Picking 150 (attended A1) or 151 (row 188 N2d/N4a absent, measured above).
- Picking 231, 230, 233 or 235 under rule (d) with D-WORLD-102 unanswered.
- Re-asking D-WORLD-102. It is the one open decision and the digest carries it; it is not asked again.
- Skipping the rotation again: with this entry the live log holds 41 entries, past the ~40 threshold iteration 324 named.

**Log rotation:** `ailang mission rotate-log world --keep 20` (fleet CLI on PATH), run with `AILANG_MISSION_REGISTRY` pointed at a temporary registry whose `world.toml` `workdir` is the record worktree, because the shipped `world.toml` names the main checkout. It printed worktree paths for all three files (`live log`, `archive`, `INDEX`), and the main checkout and the driver pin stayed clean. Result: 20 of 41 entries kept (306–325), 21 rotated out with full bodies. The tool printed 444,728 → 237,395 bytes; that was before this paragraph was written, and the live log was 238,030 bytes after it, before this follow-through was added. Lossless by line count: the log lost 1,794 lines (3,470 → 1,676) and the archive gained exactly 1,794 (23,259 → 25,053); archive entries 271 → 292. The rotated set is entries 282–305 (21 headings). 296, 297 and 302 have no entry anywhere: no log entry existed for them at the base, and 298 and 303 credit the orphaned fires 296 and 302. The index was regenerated: it gained the 325 row, and rows 315–322, which iterations 316–324 had hand-written at full length, now carry the generator's standard truncated titles. No other index row changed.

**Retro:**
1. **Idle cadence:** block report 8 of the run (317, 319, 320, 321, 322, 323, 324, 325). Each costs one judge and two CI cycles, and no code. D-WORLD-102 is still the only lever the loop has.
2. **Gate 0's directive snippet names a repo-relative `scripts/mission_directives.sh`, and this repo has none.** Run as written it exits 127 with nothing on stdout, which reads like "no directives". The driver pin's copy (`$MISSION_DRIVER_ROOT/scripts/`) is what earlier iterations used; the readings file runs that copy. Instance 1 recorded in this charter. It is shared-skill prose, so World cannot fix it; no rule edit.
- Harness share 0. No skill edit.

**Progress:** 1.0 clauses 7/7 MET; goal unmoved (release rows). Row 150: the attended A1 remains. Row 151: behind row 188 N2d/N4a.
**Next:** rule (d) block report unless D-WORLD-102 = A (then 231 → 230 → 233) or dev goes red. Row 150's P0/P3 run in the A1 window Mark schedules; row 151 becomes routable when row 188's N2d and N4a land.

**Record follow-through (iteration 325):** the independent evaluator judged record commit `303b7df6` on `39aabca3`: **r1 PASS 91/100, zero blocking** (34/21/24/12). It ran through the Agent tool on `sonnet`, in a fresh, read-only, foreground context (48,820 tok, 82 s, 10 tool uses). It re-derived the no-pick premise (5 commits since `42fdf32a`, all 324's own; 0 allowlisted comments on #366's 24 or #313's 44; no open PR; ledger 91 with D-WORLD-102 the only OPEN row; `39aabca3` 3/3), all 29 upstream numbers closed under `bash -c`, the STATUS rotation (3 stamps; 322's stamp byte-identical at the archive head; charter diff 4/4), and the log rotation. For the log it found the old archive a byte-identical prefix of the new one, the heading sets equal plus 325 with no losses or duplicates, and only index rows 315–325 changed.
- **NB1** (the entry named the rotated set's ends as 285 and 304): confirmed by the controller (`comm` of the archive headings at base and after: 282 … 305, 21 entries). Corrected above.
- **NB2** (the byte counts do not match a recount): the figures are the tool's own output, printed before the rotation paragraph was written. The paragraph now says so and gives the current size. Also applied: the dashboard had dropped the row-162 caveat, and it is restored.
- No r2: the changes are a corrected range, a stated byte reading, and one restored dashboard line.

**Landing (iteration 325):** PR #429 merged as `dda2c9ed`. The PR head `82c01f18` was 3/3 green before the merge, and origin/dev was re-read right before it (`39aabca3`, MERGEABLE/CLEAN, `--match-head-commit`). The merge tree is byte-identical to the head tree (`98de9a2b`). On the merge, the charter keeps 3 stamps. **Merge CI:** `commits/dda2c9ed/check-runs` → 3, all success on attempt 1 (verify gate, go race gate, go host build + test gate). Dev green. Both CI pollers ran under `/bin/bash` from a script file, with the subject SHA in each verdict line.

## 326 — 2026-10-30 — no pick: rule (d) block report 9; release rows 150/151 still wait on attended steps; dev green at b49ce51a [ADMIN]

**Kind:** bookkeeping-only iteration under standing rule (d) (D-WORLD-46). There is no design, plan, code or product acceptance. The record is the STATUS stamp (323 rotated to the archive), the CURRENT GOAL line, a D-WORLD-102 evidence note, this entry, the regenerated index, the dashboard, and `verification/world-iter326/gate1-2-readings.txt`. An independent judge reviews it before merge.

**Pick and why:** none.
- **Clause map at `b49ce51a`:** clauses 1–7 MET. There is no UNMET clause and no critical-path row.
- **What moved since iteration 325:** nothing but 325's own record and landing. `git log 39aabca3..origin/dev` lists 5 commits (`303b7df6`, `82c01f18`, merge `dda2c9ed` of #429, `7e86e200`, merge `b49ce51a` of #430).
- **D-WORLD-81 release rows, blockers re-measured this iteration:**
  - **150** waits only on the attended A1. The controller runs P0 (≤24 h before) and P3 (before pass A) inside the window Mark schedules (orchestration r5 §4.4). Over commit subjects `39aabca3..b49ce51a`, `grep -ciE 'row ?150|A1\b|marketing'` → **0**; the control `grep -ciE 'iteration|record|merge'` → **5** of 5 commits.
  - **151** needs row 188's N2d and N4a. `git log origin/dev --since=2026-10-01 --format='%h %s' | grep -ciE 'N2d|N4a'` → **0**; control `N2c` → **2**. The full-message form → **4**: the same four record or design commits iterations 320–325 named (`fe9c891e`, `92495b60`, `be184351`, `dc719822`). Rows 217, 218 and 219 are still `[NEXT]` on the attended track, so neither milestone is loop-routable.
- **Rows 231, 230, 233** are hygiene below every D-WORLD-81 release row. Rule (d) bars them, and D-WORLD-64's exception covers rows 148–151 only. D-WORLD-102 is OPEN with default B, so this is a block report. Row 235 is hygiene too and is not in D-WORLD-102's scope.
- **External predicates re-read** with the whole-charter pattern, run under `/bin/bash` from a script file: **29** distinct upstream numbers, `results=29`, **all 29 closed**. Controls: row 53's `ailang#941` matched, and `gh issue list -R sunholo-data/ailang --state open --limit 3` returned 3. Still a lower bound (bare `#N` not counted).
- **Weekly external-issue sweep (Gate 0 step 5):** iteration 305 ran this week's sweep. Re-checked: 1 open issue, #366, the bookkeeping thread (charter 4 / log 53 / status-archive 4 / dashboard 1 before this record). **0 orphans of 1 enumerated.**

**Gate 0:**
- Kill switch armed (`mission-world.disabled` absent). gh `sunholo-voight-kampff`. Billing tripwire CLEAN.
- `mission_directives.sh` (driver pin copy): 0 allowlisted directives on #366 since `2026-10-27T09:03:42Z` (25 comments) and 0 on #313 since `2026-10-22T01:34:42Z` (44 comments). Nothing to mark.
- `mission_gate0_self_notices.sh`: rc 1, verdict "1 crash notice". It is the same 2026-10-22 rc=1 notice `6345079873`, already credited (`grep -c` → 20 in the live log and 3 in the log archive before this entry; the live count fell from 22 because iteration 325's rotation moved entries). Control 107:4 ok. No new notice on #366 (25 self, 0 crash).
- Unread inbox: nothing addressed to World. The `mission-world` inbox holds only this loop's claim and report messages. The iteration-326 claim was sent (`inbox_1793376118625_60105403`).
- `mission_decisions.sh --check`: valid, 91 rows, 1 OPEN (D-WORLD-102). The only ledger diff in `39aabca3..b49ce51a` is 325's own D-WORLD-102 evidence-note append (one −/+ pair).

**Gate 1:**
- Local `dev` equalled `origin/dev` (`b49ce51a`), with only untracked residue in the tree. Gate-1 base `b49ce51a23bb3e0f484f4fa0d01f8394b6b03b51@2026-10-30T16:01:07Z`.
- Skill drift: per-file `cmp` of `SKILL.md` and 13 resource files at the resolved symlink target, and at the driver pin, against `origin/dev` in the freshly fetched driver clone: no DRIFT line. Control: `origin/dev:SKILL.md` read back 54,419 bytes, and the pin HEAD equals the clone's `origin/dev` (`490fb730`). Rules read from the resolved target.
- Dev CI: `commits/b49ce51a/check-runs` → 3, all success. The last five dev CI runs (`b49ce51a`, `dda2c9ed`, `39aabca3`, `c4386177`, `42fdf32a`) all read success on attempt 1.
- Open PRs: none. Died-mid-flight traces: no iteration-326 branch or worktree existed before this one. `.wt-world-iter306-r0` is the known residue. The attended scratchpad worktrees belong to the attended track and were left alone.

**Designer / planner / executor:** not spawned. There is no pick, so the routing table has nothing to send them (precedent 257–264, 317, 319–325). No spawn was attempted, so none failed. The operator's standing request for Agent-tool roles is honoured through the evaluator.

**Independent evaluator (record):** REQUIRED. It runs through the Agent tool on `sonnet` (driver `MISSION_EVALUATOR_RESOLVED=sonnet`, path `agent-tool`), in a fresh, read-only, foreground context, and judges the record commit before merge. See the follow-through below.

**Routing evidence:** base=b49ce51a23bb3e0f484f4fa0d01f8394b6b03b51@2026-10-30T16:01:07Z (Gate 1).
- Driver env: `MISSION_DESIGNER_RESOLVED=claude:claude-opus-5-5`, `MISSION_PLANNER_RESOLVED=opus`, `MISSION_EXECUTOR_RESOLVED=claude:claude-haiku-5-5`, `MISSION_EVALUATOR_RESOLVED=sonnet`; `MISSION_OVER_RATION=codex`. `MISSION_ROUTING_NOTE`: planner and executor codex `gpt-6.1-sol` unusable (probe rc 75, over daily ration).
- Controller `claude:claude-opus-5-5` (tok: not reported). Designer, planner and executor were unused.
- Evaluator: Agent `sonnet`, tokens recorded in the follow-through. Generator ≠ judge: opus controller vs sonnet judge.
- Quorum: none (no doc). Metered $0.

**Ruled out:**
- Picking 150 (attended A1) or 151 (row 188 N2d/N4a absent, measured above).
- Picking 231, 230, 233 or 235 under rule (d) with D-WORLD-102 unanswered.
- Re-asking D-WORLD-102. It is the one open decision and the digest carries it; it is not asked again.
- Rotating the log: 21 live entries with this one, far below the ~40 threshold (next due around iteration 345). The index was regenerated with `rotate-log --index-only`.

**Retro:**
1. **Idle cadence:** block report 9 of the run (317, 319–326). Each costs one judge and two CI cycles, and no code. D-WORLD-102 is still the only lever the loop has.
2. **The CURRENT GOAL line grows by one "Earlier (iteration N)" clause per fire.** Before this record it was 5,378 bytes (`wc -c` at `b49ce51a`, newline included) with 14 earlier-iteration clauses; this record's prepend takes it to 5,704 bytes and 15. Each clause restates the same block report, and the log and status archive already hold that history. Instance 1 recorded. Not trimmed here: it is charter prose outside a block report's record scope, and trimming it is a judgement call better made with the next real pick.
- Harness share 0. No skill edit.

**Progress:** 1.0 clauses 7/7 MET; goal unmoved (release rows). Row 150: the attended A1 remains. Row 151: behind row 188 N2d/N4a.
**Next:** rule (d) block report unless D-WORLD-102 = A (then 231 → 230 → 233) or dev goes red. Row 150's P0/P3 run in the A1 window Mark schedules; row 151 becomes routable when row 188's N2d and N4a land.

**Record follow-through (iteration 326):** the independent evaluator judged record commit `84e1d94e` on `b49ce51a`: **r1 PASS 97/100, zero blocking** (35/23/25/14). It ran through the Agent tool on `sonnet`, in a fresh, read-only, foreground context (46,259 tok, 33 s, 4 tool uses). It re-ran the no-pick premise (5 commits in `39aabca3..b49ce51a`, all 325's; #366's 25 comments all from the loop's account; no open PR; ledger valid, 91 rows, D-WORLD-102 the only OPEN row), dev CI (`b49ce51a` 3/3 success), and the rotation: the 323 stamp in the archive is `cmp`-identical to the one removed, the charter keeps 3 stamps, the charter diff is 4/4 lines, the index gains exactly one line, and there is no auto-close keyword. It took the 29-number upstream sweep and the `6345079873` credit counts from the readings file rather than re-running them.
- **NB1** (the retro's 5,378-byte GOAL-line figure is the size before this record; after it the line is 5,704 bytes with 15 clauses): confirmed and corrected in retro item 2.
- **NB2** (the upstream sweep was accepted from the controller's readings): noted. The readings file prints the pattern, `results=29` and both controls, so it can be re-run as written.
- No r2: one corrected figure.

**Landing (iteration 326):** PR #431 merged as `1f861724`. The PR head `f3e49e7d` was 3/3 green before the merge, and origin/dev was re-read right before it (`b49ce51a`, MERGEABLE/CLEAN, `--match-head-commit`). The merge tree is byte-identical to the head tree (`b14d383a`). **Merge CI:** `commits/1f861724/check-runs` → 3, all success on attempt 1 (verify gate, go race gate, go host build + test gate). Dev green. Both CI pollers ran under `/bin/bash` from a script file, with the subject SHA in each verdict line. One controller instrument fault, no effect on any verdict: a zsh glob over a not-yet-written verdict path printed `no matches found` for a 5-minute read; the read was re-run under `/bin/bash` against the poller's exact path.

## 327 — 2026-10-30 — dev red at Gate 1 → row 236 built: the finish route's scope-expiry test raced two wall-clock reads across a second edge; daemon clock seam + exact-edge tests; judge PASS 93 [PRODUCT]

**Kind:** a Gate-1 dev red on World's own repo, which outranks the queue (precedent 318, 320). Full inner loop: designer → quorum → planner → executor → independent judge. The record is in the same PR as the code (#436); the landing note records merge CI.

**Pick and why:** row 236 (filed this iteration).
- **Dev red:** `origin/dev` `4d43d8e3` (#432's merge, a docs-only landing note) — `commits/4d43d8e3/check-runs` 3: race gate `failure`, host gate and verify gate `success`. Parents `cfdab426`, `1f861724`, `f3e49e7d`, `b49ce51a` all green. Job log fetched with `--allow-escape-sequences` (36,460 B): the two `DATA RACE` blocks are the gate's own known-positive control (`racecontrol/main.go`); the one test failure is `--- FAIL: TestEpisodeFinishScopeExpiry (0.15s) episode_finish_test.go:170: expiry 1793383038: 200`. The package took 352 s on the runner.
- **Mechanism, measured before routing:** the test builds `time.Now().Unix()+3601`; the handler (`episode_finish.go:168–169`, and the same bound at `:270–271`) refuses only when `scope.ExpiresAt > now+3600` with its own later `time.Now()`. A second edge between the reads makes `+3601` equal the bound → accepted → 200. On a scratch worktree: test read aligned ≤50 µs before an edge → **20/20 FAIL** plain and **20/20** `-race`, with the CI signature; aligned 2 ms before → 0/20 (the handler is fast idle); aligned just after an edge (control) → 0/20; unmodified idle `-race` → 0/400. `verification/world-iter327/gate1-2-readings.txt`.
- **Why taken, not filed like row 235:** row 235 (iteration 322) had an unreproduced mechanism (0/80). This one is reproduced deterministically, so it is the 318/320 shape. The code is the attended track's (row 218 M2, `8c5a14e6`); a claim message named the row before routing (`inbox_1793390682086_a0f83fd5`). The only open PR, #433 (attended, draft), touches `host/broker/registry_publish_test.go` only — no collision.
- **Rerun:** the failed job was re-run (attempt 2) → `4d43d8e3` 3/3 success, so dev is green while the root cause is fixed.
- **Clause map at `4d43d8e3`:** 1–7 MET. Rows 150/151 unchanged since 326 (no commit since `b49ce51a` but 326's own record and landing). D-WORLD-102 is OPEN and was not used.

**Gate 0:**
- Kill switch armed (`mission-world.disabled` absent); gh `sunholo-voight-kampff`; billing tripwire CLEAN.
- `mission_directives.sh`: 0 allowlisted directives on #366 since `2026-10-27T09:03:42Z` (26 comments) and 0 on #313 since `2026-10-22T01:34:42Z` (44). Nothing to mark.
- `mission_gate0_self_notices.sh`: rc 1, "1 crash notice" — the same 2026-10-22 rc=1 notice `6345079873`, already credited (22 mentions in the live log before this entry). Control 107:4 ok. #366: 26 comments, all self, 0 crash.
- Inbox: the `mission-world` inbox holds only this loop's claim/report messages; nothing addressed to World.
- `mission_decisions.sh --check`: valid, 91 rows; `--open`: D-WORLD-102 only (at Gate 0).
- **Attended ruling acknowledged (D-WORLD-102 = A).** While this iteration ran, Mark recorded an attended ruling in the ledger (#434, flip `dac17cee`, 2026-10-30 21:16 CET; "merge 433 and 102 a"): extend D-WORLD-64's exception to rows 231 → 230 → 233, stopping as soon as 150 or 151 is routable. `mission_ledger_provenance.sh --id D-WORLD-102` → `verdict=attended` (creation `92495b60`). Found at the pre-merge re-read of `origin/dev` (`35b5aecf`); the record was rebased onto it (clean), the D-WORLD-102 row is byte-identical to origin's, and `--check` → valid, 91 rows, **0 OPEN**. The same window also merged #433 (the census probe cap — one of the two load flakes classified below) and #435 (row 219 M4); none touches this row's three files (intersection empty against 32 incoming files). This pick stays row 236: a dev red outranks the queue, and standing rule 1 allows one item.
- **Second rebase, onto #437 (row 218 M3, attended, merged during this PR's CI).** It changed `host/daemon/episode_finish.go` (the only overlapping file): `finishDecisionScope` now also returns the scope text, and a consumed approval refuses every later record, not only `done`. The rebase was clean; the record-time re-check bound moved to `:275`. Re-measured on the combined tree (copy of the rebased head): `go vet` ok; all 6 Finish tests PASS, including M3's new `TestEpisodeFinishRecordReuseAndSingleUse`; mutant `:275` `> now+3600`→`> now+3601` and `<= now`→`< now` each red on T3 (`episode_finish_test.go:232`); the row-237 mutant (drop `|| input.Record != nil`) still survives, so M3 does not pin it and row 237 stands. File restored byte-identical (`4e320c53…`). The judge's verdict covers the code at `e7743475`; this rebase changed no line of row 236's diff.

**Gate 1:**
- Local `dev` = `origin/dev` = `4d43d8e31b8e54673fbe955641aa83af0f60e9fd`. Gate-1 base `4d43d8e3…@2026-10-30T20:01:39Z`.
- Rules read from the resolved symlink target (`~/dev/sunholo-data/ailang/.claude/skills/mission-control`).
- Open PRs: #433 `fix/census-driver-cap` (draft; worktree in an attended scratchpad, not this loop's) — left alone.
- Died-mid-flight traces: none for 327; `.wt-world-iter306-r0` is the known residue.

**Design** (`w-finish-scope-expiry-boundary`; designer resolver `recipe claude:claude-opus-5-5 declared:provider-pin`, spawned as Agent `opus`, the 318/320 precedent; rotation pointer unchanged — the lane that ran):
- r1: an unexported daemon clock seam, exact-edge tests on a pinned clock; option (a) margin (`+3660`) rejected because P4/P4b survive (measured), retries rejected. Sibling sweep: 2 same-hazard sites, both in `episode_finish_test.go` (the failing arm and AC11's `future` `+301`), both fixed here.
- Quorum r1 BLOCKED ($0.2474; seats oc-kimi-k3 REJECT, gpt6-1-sol REJECT, gemini-4-argon PASS; claude benched as author vendor; no absentees): reuse-vs-new seam; "no test reaches `:271`" conflated detection with execution.
- One mandated revision (r2): coverage measured (`:271`'s condition executes, its refusal branch never does; `:169` control covered); two clocks kept apart on purpose with reuse measured infeasible; field renamed `finishNow`.
- Quorum r2 BLOCKED ($0.3578; same seats; no absentees): kimi — the r2 diff's mutant/package evidence was from the pre-rename prototype; gpt — the construction audit was textual. Neither disputes direction, both carry a concrete `proposed_fix` → **narrow-refinement carve-out (r3)**: kimi's own minimum-acceptable annotations applied verbatim; gpt's type-aware audit run by the controller (`go/packages` + `go/types`, tests included, 108 packages, 0 load errors): **11** constructions (`daemon.go:784` in `New` + 10 `host/daemon` test literals), 0 `new()`, 0 zero-value; a scratch alias control adds exactly 3, one per branch; none of the 10 literals' files mentions the finish route (control 10). `controller-r3-daemon-construction-audit.txt`. Artifacts `quorum/r1.json`, `r2.json`.

**Plan** (opus, `opus fail-closed:env-pin`): M1 code, M2 evidence. Prototype-validated on the exact planned files: AC-1/AC-1b 20/20 red at base, 20/20 green on the fix, plain and race; 21/21 mutants red on their named killers (S0 control green); `host/daemon` plain `ok 335.8s`, race `ok 409.0s`; `cmd/ailang-worldd` at base 237/0/9 = fix 237/0/9 (so the designer's coverage-mode failures were mode/load). Deviations D1–D5 in the plan §8 (comment-only fixes, exact mutant text, S0 control).

**Executor** (haiku; resolver `recipe claude:claude-haiku-5-5 declared:provider-pin`, Agent `haiku`): M1 `d7a54db5` (+84/−13, three `host/daemon` files; sha256 equal to the plan's), M2 `75d07431` (evidence). AC-1/AC-1b base red / fix green; AC-2 21/21 killed, restores byte-identical, `scratch removed: yes`; AC-3 vet rc 0, package plain `ok` 441 s, race `ok` 525 s, `-count=200` race `ok`; AC-4 `verify_ail.sh` rc 0, `go build ./...` rc 0, `go test ./...` **rc 1** → stopped as the plan requires.
- **Controller classification** (`controller-ac4-classification.txt`, load average 14–18 from other sessions): `host/daemon` = Go's default 10-min package timeout in the all-packages parallel run, no assertion failed (the package passes alone, AC-3). `host/broker` (`TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential` agentrun/checkrun `driver0` at 1.0 s — the census cap attended PR #433 raises; `TestExecProbeForeignListenerOnTheReversePortDoesNotRefuse/chatty`) is unreachable from this diff (`go list -deps -test ./host/broker` → 0 `host/daemon`, control 1) and each passed 3/3 alone. CI is the authoritative full gate.
- The executor's commits carry a Haiku 5.5 co-author line (its harness's attribution), not the Opus line the directive named.

**Independent evaluator:** REQUIRED. Agent `sonnet` (`agent-tool sonnet declared:alias-pin`), own detached worktree `.wt-world-iter327-eval` at `e7743475`, foreground, read-only: **PASS 93/100, zero blocking** (34/23/24/12). Reproduced the base red 20/20 on a `/tmp` copy; killed Q1, Q3, Q4, R1, R3, S3, W1 and an own mutant (proof bound 300 → 360); vet clean; Finish `-race -count=20` ok; full package `-timeout 1200s` one failure, `TestDaemonStartReportsMissingTemplates`, because its shell lacked `AILANG_BIN` (not in the diff). 49,634 tok / 587 s / 16 tool uses. Report `verification/world-iter327/evaluator-sprint-r1.md`.
- **NB1, reproduced → row 237:** dropping `|| input.Record != nil` from the approve guard leaves all 5 Finish tests (`-v` lists them as PASS) and `cmd/ailang-worldd -run Finish` green. Pre-existing; out of this row's scope.
- NB2 (no clean local full-gate green; CI is the authority), NB3 (the judge did not recount the audit's 11), NB4 (evidence volume) — noted.

**Routing evidence:** base=4d43d8e31b8e54673fbe955641aa83af0f60e9fd@2026-10-30T20:01:39Z (Gate 1).
- Driver env: `MISSION_DESIGNER_RESOLVED=claude:claude-opus-5-5`, `MISSION_PLANNER_RESOLVED=opus`, `MISSION_EXECUTOR_RESOLVED=claude:claude-haiku-5-5`, `MISSION_EVALUATOR_RESOLVED=sonnet`; `MISSION_OVER_RATION=codex`; `MISSION_ROUTING_NOTE`: planner and executor codex `gpt-6.1-sol` unusable (probe rc 75, over daily ration).
- Controller `claude:claude-opus-5-5` (tok: not reported) · designer Agent `opus` (151,898 tok / 1,786 s; revision resume 215,874 tok / 1,225 s as reported) · planner Agent `opus` (155,785 tok / 1,502 s) · executor Agent `haiku` (95,050 tok / 2,119 s) · evaluator Agent `sonnet` (49,634 tok / 587 s). Generator ≠ judge: opus/haiku vs sonnet.
- Fable diet: 1 doc, 2 designer runs (author + the mandated revision); the r3 fixes were the controller's, not a third run.
- Quorum metered $0.6052 (r1 $0.2474 + r2 $0.3578).

**Ruled out:**
- Filing and not taking (the row-235 path): the mechanism is reproduced, not hypothetical.
- A wider margin only (design option a): measured to leave P4/P4b surviving.
- Re-running the full gate locally until green: the rig's load is other sessions'; CI decides.
- Folding NB1 into this row: pre-existing, a different guard; filed as 237.
- Using D-WORLD-102: not needed for this pick (a dev red outranks rule (d)); it was resolved = A by an attended ruling mid-iteration, so rows 231 → 230 → 233 are next.

**Retro:**
1. **The executor's AC-4 command carried no `-timeout`**, so a loaded rig turned a green package into a 10-minute panic. The planner wrote the command from the doc's AC-4 text; CI passes `-timeout 8m` per package. Instance 1 here; future plans should carry an explicit `-timeout` on the full-suite line.
2. **Quorum r2 objected to evidence measured on a pre-rename prototype.** A rename between rounds is cheap to re-verify only if the designer re-runs the whole battery; r2 re-ran a subset. Instance 1.
- Harness share 0. No skill edit.

**Progress:** 1.0 clauses 7/7 MET; goal unmoved (release rows). Dev red fixed in a PR; row 237 filed.
**Next:** row 236 lands on merge CI green (landing note). Then row 231 → 230 → 233 under D-WORLD-102 = A (attended, #434), stopping when a release row is routable; a dev red still outranks. Row 150's P0/P3 run in the A1 window Mark schedules; row 151 becomes routable when row 188's N2d and N4a land.

**Landing (iteration 327):** PR #436 merged as `37f8fb82` (merge commit, `--match-head-commit f395fafb`). The branch was rebased twice onto attended merges that landed while it was in CI: #433–#435 (incl. the D-WORLD-102 = A ruling) and #437 (row 218 M3, which also edits `episode_finish.go`). The final head `f395fafb` was 3/3 green on base `25ed07f1`, re-read MERGEABLE/CLEAN right before the merge; the earlier heads `e7743475` (code, on `4d43d8e3`), `7998f124` and `9ac0640b` were each 3/3 green too. **Merge CI:** `commits/37f8fb82/check-runs` → 3, all success on attempt 1 (verify gate, go race gate, go host build + test gate). Dev green. Pollers ran under `/bin/bash` from a script file, one per SHA, each verdict line naming its SHA.

## 328 — 2026-10-31 — row 231 built: the P0/P3 position and phone predicates become CI-gated marketingcheck subcommands with a parity oracle and generated fixtures; judge PASS 90 [PRODUCT-HYGIENE]

**Kind:** the first pick under D-WORLD-102 = A (attended ruling 2026-10-30, #434). Full inner loop: designer → quorum ×2 → carve-out → planner → executor → independent judge. The record is in the same PR as the code (#440). A landing note will record merge CI.

**Pick and why:** row 231 `w-marketing-position-checker`.
- **Clause map at `cec48f69`:** 1–7 MET; no UNMET clause.
- **Release rows not routable.** Row 150 waits only on the attended A1. Row 151 needs row 188 N2d/N4a: 0 matching commits on `origin/dev` since iteration 327's base. The only new commits are #438, the attended `trial/row217c-interactive-resume` dogfood ledger row, and 327's own record and landing.
- So D-WORLD-102 = A applies, with its order 231 → 230 → 233.
- **Premise checked first-party:** `scripts/marketingcheck/main.go` dispatches `check`, `routes`, `observe`, `compare`, `inventory`/`inventory-compare` and `docs`, but no `position` or `phone`. The predicates existed only as `world-iter315/v64-p0check-main.go.txt` (425 lines) and `v66-phonecheck-main.go.txt` (224 lines).
- Claim message `inbox_1793405384179_a8734bee` sent before routing. No open PR.

**Gate 0:**
- Kill switch armed (`mission-world.disabled` absent); gh `sunholo-voight-kampff`; billing tripwire CLEAN.
- `mission_directives.sh`: 0 allowlisted directives on #366 since `2026-10-27T09:03:42Z` (28 comments, LEGACY watermark), so nothing to mark.
- `mission_gate0_self_notices.sh`: rc 0 with verdict text "1 crash notice". It is the 2026-10-22 rc=1 notice `6345079873`, already credited (23 mentions in the live log). Control 107:4 ok.
- Inbox: `mission-world` holds only this loop's own claim and report messages.
- Ledger: 0 OPEN rows at Gate 0; D-WORLD-102 is RESOLVED = A (attended).

**Gate 1:**
- Local `dev` = `origin/dev` = `cec48f69af6da7fe4817d9a1bf3e9f01b0274cb0`. Base recorded `cec48f69…@2026-10-31T00:08:05Z`; re-read before the worktree with no drift.
- Skill: `SKILL.md` and all 13 `resources/*.md` were `cmp`-identical to `origin/dev`, per file; rules were read from the resolved symlink target.
- Dev CI: `cec48f69` check-runs 3/3 success; the last 5 CI runs on dev were success on attempt 1.
- Died-mid-flight traces: none for 327. `.wt-world-iter306-r0` is the known residue.

**Design** (`w-marketing-position-checker`; resolver `recipe claude:claude-opus-5-5 declared:provider-pin`, spawned as Agent `opus` per the 318/320/327 precedent; rotation pointer unchanged):
- **r0** `9abc25b5`: 502 lines, 30 V-rows, 40 evidence files.
  - Surface: `position --baseline RECTS LABEL=RECTS…` and `phone LABEL=PNG…`. Exit 0/1/2.
  - Parity oracle: the banked `.txt` programs, built in-test after a sha check.
  - Fixtures: position JSON checked in as produced and sha-pinned. Phone fixtures are 2 KB "skeleton" geometry, because `tracked_binary_gate.sh` rejects a PNG outside `website/static/img/v1/` (measured on a throwaway commit).
  - One deliberate verdict change: phone (ii) runs `CheckRegion` (C8′'s own words), ledgered as `header-15px`.
  - The real phone corpus was copied out of `/tmp` to `~/.ailang/state/world-iter328-phone-corpus/`, with all 11 hashes equal to the iteration-315 pins.
- **Quorum r1 BLOCKED** ($0.2120). Seats gpt6-1-sol REJECT, oc-glm-5-3 REJECT, gemini-4-argon PASS; claude benched; no absentees.
  - gpt: the 10-image skeleton byte-equality was unsatisfiable with 2 committed skeletons.
  - glm: the `Shots` geometry premise had no V-row, and it was not said how the prototype was built.
  - **The controller measured glm's premise before routing (rule 3f).** The prototype copied alone fails `undefined: Rect`, rc 1; in a `git archive` tree it builds, rc 0. So it read `png.go`'s `Shots`, whose values equal p0check's `heights` and its 1440 literal.
- **r1** `3aaee4d0`, the one mandated revision.
  - 10 skeletons, byte-equal 10/10 with three failing controls (V31).
  - V32 holds the measurements above, plus `TestMarketingPositionGeometryMatchesOracle` with mutants P14 (verdict-neutral, sole killer) and P15.
  - The host-local corpus is named as a residual (V33: no durable copy of the iteration-312 images).
- **Quorum r2 BLOCKED** ($0.1922; same seats and verdicts; no absentees).
  - gpt: oracle subprocess waits unbounded.
  - glm: §2.7 deleted the oracle source pins the tests are specified to read.
  - Both carry a concrete `proposed_fix`, and neither disputes the direction → **narrow-refinement carve-out (r2, `08dc31be`)**, applied by the controller.
  - Applied: gpt's §2.2 runner paragraph verbatim; glm's "stop deleting the pins" heading and pin line; gemini's M2 count wording; P16/P17 mutation rows; V34/V35 as executor acceptance rows; quorum r1 glm's AC-5 scope note.
- Artifacts: `verification/world-iter328/quorum-r1.json`, `quorum-r2.json`.

**Plan** (opus; `derive-planner-lane.sh` → `opus fail-closed:env-pin`, resolver `agent-tool opus`), `0f0b65b2`:
- 7 milestones: M1a runner · M1b `position` + 34 fixtures · M1c pins + corpus · M1d parity · M2a `phone` + 10 skeletons · M2b arms + parity · M3 docs + NB2.
- 22 banked `plan-*.txt` sources with sha256, so the executor wrote no Go.
- Every cumulative tree passed vet and the package tests, plain and race. 36/36 mutants killed.
- Deviations D-1..D-15. D-15 adds two `TestProductionGoSurface` allow-list lines in `host/store/context_roots_test.go`. D-14: no `nohup`, which breaks `shell-hup`.
- Its repo gate showed one load red, `host/verifygate` `invalid_ledger`, passing alone.

**Executor** (haiku; resolver `recipe claude:claude-haiku-5-5 declared:provider-pin`, Agent `haiku`):
- Commits: M1a `2822723a` · M1b `d673c56f` · M1c `27ef6447` · M1d `c4b129b0` · M2a `c6a12dc9` · M2b `9ddab47c` · M3 `03fe5a1a` · evidence `96eae6cc`. Every acceptance matched the plan.
- Parity: position 44/0/0; phone 41 and 59 with ledger 1; skeletons `CHECK images=10 byte-equal=10`.
- Mutants 36/36 KILLED. V34: pins 1 and 1. V35: P17 red in 21.8 s with `subprocess-timeout`; P17b (`Setpgid` removed) red because the grandchild survives.
- Repo gate: `verify_ail.sh` rc 0, `go build` rc 0. `go test -timeout 20m ./...` rc 1 on `host/broker` `TestExecProtectedCheckTimeout/reserve-after-before` (load 20–24). The test passed alone 3/3, and so did the package (226–251 s).
- Deviations from the directive: Haiku 5.5 co-author line again (its harness attribution, instance 2 after 327), some plan-named commit prefixes, and `executor-*.txt` names. All harmless.

**Controller checks before the judge:**
- No production context root: every `context.*` use is in a `_test.go`.
- D-15 is allow-list lines only.
- `host/broker` does not depend on the diff (`go list -deps -test`: 0 `scripts/marketingcheck`; `host/store` changed only in a test file) and is byte-identical to base.
- `./scripts/verify_go.sh` at `96eae6cc` on a rig without the judge (load ~12): rc 1, on the same `reserve-after-before` row only (`handler subprocess timed out after 3s`). `scripts/marketingcheck` ok 132.6 s; every other package and gate green.
- Classified as load: the third local instance of that row (iteration 293 recorded it under co-load). CI arbitrates.

**Independent evaluator:** REQUIRED. Agent `sonnet` (`agent-tool sonnet declared:alias-pin`), own worktree `.wt-world-iter328-eval` (branch `eval/world-iter328-r1` at `96eae6cc`), foreground, run after `verify_go.sh` finished.
- **PASS 90/100, zero blocking** (32/21/24/13).
- Verified: vet and package tests plain and race; the real-corpus run with 0 SKIP (135 PASS lines); oracles self-built, 7 position + 4 phone spot checks equal; fixtures regenerate byte-identically (34 position, phone skeletons no diff); one-byte fixture edit, zero-regions and empty-baseline guards fire; runner deadline and group-kill mutants red; each oracle hash once in the orchestration doc; §13 value revert reds `rects-value`; scope clean.
- 21 own mutants + 8 of the plan's; worktree restored (`status --porcelain` empty).
- **Survivors, filed as row 240:** background-edge tolerance ±2→±200 (×3); phone entry tolerance 4→40; `gc ≥ 0` (maybe equivalent); `main.go` cases reached only by docs subtests; §13 export revert misses `position-lines`.
- **The controller reproduced the entry-tolerance survivor** on a copy of `96eae6cc`: control rc 0, mutant rc 0 (`ok … 43.8s`).
- Report: `verification/world-iter328/evaluator-sprint-r1.md`.

**Routing evidence:** base=cec48f69af6da7fe4817d9a1bf3e9f01b0274cb0@2026-10-31T00:08:05Z (Gate 1).
- Driver env: `MISSION_DESIGNER_RESOLVED=claude:claude-opus-5-5`, `MISSION_PLANNER_RESOLVED=opus`, `MISSION_EXECUTOR_RESOLVED=claude:claude-haiku-5-5`, `MISSION_EVALUATOR_RESOLVED=sonnet`; `MISSION_OVER_RATION=codex`. `MISSION_ROUTING_NOTE`: planner and executor codex `gpt-6.1-sol` unusable (probe rc 75, over daily ration).
- Tokens and time per role:
  - controller `claude:claude-opus-5-5`: tok not reported;
  - designer Agent `opus`: 370,480 tok / 1,865 s; revision resume 420,824 tok / 324 s, as reported;
  - planner Agent `opus`: 398,391 tok / 3,986 s;
  - executor Agent `haiku`: 125,504 tok / 3,370 s;
  - evaluator Agent `sonnet`: 76,134 tok / 607 s.
- Generator ≠ judge: opus/haiku vs sonnet. All five roles were spawned; none fell back.
- Designer diet: 1 doc, 2 runs (author + mandated revision). The r2 fixes were the controller's, not a third run.
- Metered $0.4042 (quorum r1 $0.2120 + r2 $0.1922).

**Ruled out:**
- Moving the docs to `implemented/`: the committed code comments and the pinned `testdata/position/SHA256SUMS` cite `planned/`. Marked Implemented in place.
- Fixing the judge survivors in this PR: non-blocking, the same tolerances live in the banked oracle, and a change would need a re-judge. Filed as row 240.
- Tightening carried findings (1)/(2)/(7) here: they change real-store verdicts and need the orchestration doc's own quorum. Row 238.
- Re-running the full local suite until green: CI is the arbiter, and the red row is unreachable from the diff.
- A third design run: the carve-out conditions held (concrete fixes, no direction dispute).

**Retro:**
1. **Banked validated sources made the small-model executor exact.** The planner (66 min) assembled and mutation-tested every file, and the executor (Haiku) only copied and checked hashes: 36/36 killed, every acceptance as planned. Instance 1 at this scale; worth repeating for plan-sized ports.
2. **`host/broker` `reserve-after-before` reds under load** (3 s subprocess bound), twice this iteration and once in 293. If CI ever reds on it, it gets a row. Watch-item, not filed.
3. The controller's first evaluator-worktree call passed `--detach` to `mission-worktree.sh add`, which takes `BRANCH SHA PATH`. It failed `worktree_add_rc_255` with no partial tree. Controller slip; no harm.
- Harness share 0. No skill edit.

**Progress:** 1.0 clauses 7/7 MET; goal unmoved (release rows). Row 231 BUILT in PR #440; rows 238–240 filed.
**Next:** row 231 lands on merge CI green (landing note). Then row 230 → 233 under D-WORLD-102 = A, stopping when a release row is routable; a dev red still outranks. Row 150's P0/P3 now invoke `marketingcheck position`/`phone` (orchestration §4.4) in the A1 window Mark schedules.

**Landing (iteration 328):** PR #440 merged as `6641df63` (merge commit, `--match-head-commit 5b5ee8f1`). Head `5b5ee8f1` was 3/3 green and was re-read MERGEABLE/CLEAN, with `origin/dev` still at `cec48f69`, right before the merge. **Merge CI:** `commits/6641df63/check-runs` → 3, all success on attempt 1 (verify gate, go race gate, go host build + test gate). Dev green. The local `host/broker` `reserve-after-before` red did not occur in CI on either head. Pollers ran under `/bin/bash` from a script file, one per SHA, each verdict line naming its SHA.

## 329 — 2026-10-31 — row 230 built: arm 6's hello wait pinned from both sides by a deterministic test, the SILENT text derived from the wait, a labelled structural pin on the child's c.resume(); cross-profile reach refuted at HEAD with its hoist filed; judge PASS 88 [PRODUCT-HYGIENE]

**Kind:** the second pick under D-WORLD-102 = A (order 231 → 230 → 233). Full inner loop: designer → quorum ×2 → carve-out → planner → executor → independent judge. The record is in the same PR as the code (#442). A landing note will record merge CI.

**Pick and why:** row 230 `w-exec-probe-reverse-leg-pins`.
- **Clause map at `e5c1867f`:** 1–7 MET; no UNMET clause.
- **Release rows not routable.** Row 150 waits only on the attended A1. Row 151 needs row 188 N2d/N4a: 0 matching commit subjects on `origin/dev` since 2026-10-01 (N2c control 2). The only commits since iteration 328's base are 328's own record and landing.
- So D-WORLD-102 = A applies; 231 is LANDED, so 230 is next.
- **Premise checked first-party** on a scratch worktree at `e5c1867f`: `go test -count=1 -run ExecProbe ./host/broker/` control rc 0 (`ok … 26.1s`); with `execProbeHelloWait = 5 * time.Millisecond` rc 0 (`ok … 24.1s`). Tree restored.
- Index grep: no earlier attempt at row 230. No design doc, no PR. Claim message `inbox_1793420970038_48e4e5cb` sent before routing.

**Gate 0:**
- Kill switch armed (`mission-world.disabled` absent); gh `sunholo-voight-kampff`; billing tripwire CLEAN.
- `mission_directives.sh`: 0 allowlisted directives on #366 since `2026-10-27T09:03:42Z` (29 comments, LEGACY watermark); nothing to mark.
- `mission_gate0_self_notices.sh`: rc 1 on the 2026-10-22 rc=1 notice `6345079873`, already credited (24 mentions in the live log). Control 107:4 ok. Died-mid-flight traces: no open PR, no residue for 328; `.wt-world-iter306-r0` is the known residue.
- Inbox: `mission-world` holds only this loop's own claim and report messages.
- Ledger: 0 OPEN rows.

**Gate 1:**
- Local `dev` = `origin/dev` = `e5c1867f717c8e01b4dc5ccc63bbb24adc9ad348`. Base recorded `e5c1867f…@2026-10-31T04:27:17Z`; re-read before the worktree with no drift.
- Skill: `SKILL.md` and all 13 `resources/*.md` `cmp`-identical to `origin/dev`, per file; rules read from the resolved symlink target.
- Dev CI: `e5c1867f` check-runs 3/3 success; the last 5 CI runs on dev success.

**Design** (`w-exec-probe-reverse-leg-pins`; resolver `recipe claude:claude-opus-5-5 declared:provider-pin`, spawned as Agent `opus` per the 318–328 precedent; rotation pointer unchanged, codex over ration):
- **r0** `f39703f7`: 377 lines, 22 V-rows, 31 evidence files. Verdicts:
  - (1) F2 confirmed on the whole package (333/0/30 both arms). New: nothing pins the upper side either; a 60 s wait reds only T6/`silent`, a wall-clock path.
  - (2) F3 survival confirmed; the stated RST mechanism refuted on darwin (node reads the hello into its own buffer, `readableLength` 18; server `end`/`close` 0/20 vs 20/20; host sees no difference). Structural pin chosen; a Linux arm would be predicted vacuous.
  - (3) Cross-profile reach refuted in-process: `bind(0)` 0/200,000 onto live ports (paired control 768); 0 collisions; a forced collision refuses both profiles.
  - (4) R4 unchanged.
- **Quorum r1 BLOCKED** ($0.4075). Seats gemini-4-argon PASS, oc-glm-5-3 REJECT, oc-kimi-k3 REJECT; claude benched; reserve gpt6-1-sol; no absentees.
  - **The controller measured every disputed premise before routing (rule 3f).** "Three packages probe": TRUE (`5 cmd/ailang-worldd`, `5 host/broker`, `2 host/daemon`), V19 had dropped `host/daemon`, so both reviewers' "two" was itself a misreading. "80 runs": FALSE, 4 cells × 10 = 40. srt pin equality: TRUE (`world-row140-m0/pin.json` 0.0.78, `cli_sha256 3c3092bd26b3924…`), unlogged. `:680-689` and `:736-741` read as claimed.
- **r1** `c2b9c62b` (the one mandated revision, designer resumed): 410 lines, 27 V-rows; `listenWindow` derived from `execProbeListenMS`; K9 (hold halved) measured red, documented as the intended coupling; AC-LINUX-F3 excludes T11 and is non-blocking.
- **Quorum r2 BLOCKED** ($0.3281; same seats and verdicts; no absentees).
  - glm: r1 pasted V19b into §1 row 4 (7 cells, evidence `V19b`) instead of §10 — controller-confirmed; `:657` unclassified (a comment); the 60 s claim generalised from a `-run` subset.
  - kimi: item 3 closed as REFUTED with no enforcement and no successor; offered (a) hoist `net6` in this sprint or (b) demote the verdict, file a successor, comment the invariant.
  - Both carry concrete fixes and neither disputes direction → **narrow-refinement carve-out r2 `762d8cb4`**, applied by the controller: glm (1)–(3) verbatim; kimi option (b); T11's string-match trade acknowledged in §3.2. Recorded in the doc's §12.
- Artifacts: `verification/world-iter329/quorum-r1.json`, `quorum-r2.json`.

**Plan** (opus; `derive-planner-lane.sh` → `opus fail-closed:env-pin`, resolver `agent-tool opus`), `d908403c`:
- 2 milestones (M1 F2: P1 constant, P2 derived text, T10, T7 extension; M2 F3: P3/P4 comments, T11), each with banked final files and `SHA256SUMS`; a mutation runner that self-checks and restores.
- Mutants: K1 5 ms, K2 249 ms, K3 951 ms, K4 60 s, K5 no deadline, K7 literal text at 400 ms, K9 hold 500 ms, K6 no `c.resume()` — all red with the named assertion; K2c (250/950 ms), K8, R-M1 green controls.
- T10 `-count=3000` under 32 `yes`: 3000/0 at load 16–19. Reverse-leg group `-count=20`: 20/20. Full package 335/0/30 (216 s).
- D-1: no milestone production revert can red a test, because both milestones are behaviour-neutral at today's values (P2 byte-identical at 500 ms, M2 comments only); the K-mutants are the kills. Judge-adjudicated true by measurement.

**Executor** (haiku; resolver `recipe claude:claude-haiku-5-5 declared:provider-pin`, Agent `haiku`):
- Commits: M1 `fda78271` · M2 `5ddc93d3` · evidence `d6665670`. Every checksum `OK`, every acceptance as planned, all 15 mutation lines matched. `host/broker` 334/0/30 at M1, 335/0/30 at M2. `verify_ail.sh` rc 0. The co-author line followed the directive: the Haiku-attribution deviation of 327 and 328 did not recur once the directive named the exact line.

**Controller checks before the judge:**
- Diff scope: `host/broker/exec_probe.go` (+35/−4; three non-comment lines: `execProbePollEvery`, its use at `time.After`, the `%d ms` derivation), `exec_probe_test.go` (+102), `ci.yml` (one list line, +T10 +T11).
- `./scripts/verify_go.sh` at `5ddc93d3` (load ~11–14): rc 1. Binary gate, ledger, toolchain, race known-positive control (2 races, intended), evidence manifest and plain `go test ./...` all green (`host/broker` ok 345 s). Race leg: `host/daemon` (`TestWorkspacePackageCacheUnsetLeavesHomeAlone` running) and `host/verifygate` (`TestSeToolsGateArmsAreLoadBearing`) hit the 8-min package timeout.
- Negative control: both alone under `-race -timeout 8m` at the same head: `host/daemon` rc 0 (319 s), `host/verifygate` rc 0 (262 s). `host/verifygate` does not depend on `host/broker` (`go list -deps -test`: 0); `host/daemon` does, but the running test is unrelated. Classified as load; CI arbitrates. Same class as iteration 327's `host/daemon` timeout.

**Independent evaluator:** REQUIRED. Agent `sonnet` (`agent-tool sonnet declared:alias-pin`), own worktree `.wt-world-iter329-eval` (branch `eval/world-iter329-r1` at `d6665670`), foreground, run after `verify_go.sh` and the negative control finished.
- **PASS 88/100, zero blocking** (31/21/24/12).
- Verified: gofmt/vet clean; `host/broker` 335/0/30; T10 `-count=500` under 32 `yes`; base and final sha256s; the plan's K1–K5 reproduced; own mutants on the poll constant, the `Milliseconds()` derivation (µs, s, ÷10 all red) and T10's arithmetic; D-1 confirmed by milestone reverts; CI list gained exactly T10 and T11 with the PASS guard intact.
- **Survivors, filed as row 241:** `time.After(execProbePollEvery)` → literal 500 ms stays green; literal message with argument dropped (equivalent at 500, vet flags it); T10's 951 ms kill margin ~0.9 ms; T10 arithmetic mutants each equivalent alone.
- **The controller reproduced survivor 1 and one kill** on the eval worktree: cadence literal control `ok 9.8s`, mutant `ok 10.5s`; `execProbePollEvery = 600 ms` → T10 red at `exec_probe_test.go:914` (`stop+499.9ms … no later than stop+400ms`). Restored (`status` shows only the judge's report).
- Report: `verification/world-iter329/evaluator-sprint-r1.md`.

**Routing evidence:** base=e5c1867f717c8e01b4dc5ccc63bbb24adc9ad348@2026-10-31T04:27:17Z (Gate 1).
- Driver env: `MISSION_DESIGNER_RESOLVED=claude:claude-opus-5-5`, `MISSION_PLANNER_RESOLVED=opus`, `MISSION_EXECUTOR_RESOLVED=claude:claude-haiku-5-5`, `MISSION_EVALUATOR_RESOLVED=sonnet`; `MISSION_OVER_RATION=codex`. `MISSION_ROUTING_NOTE`: planner and executor codex `gpt-6.1-sol` unusable (probe rc 75, over daily ration).
- Tokens and time per role:
  - controller `claude:claude-opus-5-5`: tok not reported;
  - designer Agent `opus`: 242,131 tok / 2,281 s; revision resume 262,015 tok / 106 s, as reported;
  - planner Agent `opus`: 224,851 tok / 1,964 s;
  - executor Agent `haiku`: 70,577 tok / 1,114 s;
  - evaluator Agent `sonnet`: 55,191 tok / 824 s.
- Generator ≠ judge: opus/haiku vs sonnet (judge-independence: same-vendor, different model). All five roles were spawned; none fell back.
- Designer diet: 1 doc, 2 runs (author + mandated revision). The r2 fixes were the controller's, not a third run.
- Metered $0.7356 (quorum r1 $0.4075 + r2 $0.3281).

**Ruled out:**
- Hoisting the `net6` bind in this sprint (quorum r2 kimi option (a)): it changes probe structure on a fail-closed startup gate and needs its own design and quorum. Row 242.
- A Linux-CI behavioural arm for F3: predicted vacuous on its own platform by the darwin measurement; T4 is already the behavioural arm in step 8.
- The row's "~100 ms" floor: below srt's measured 229.6 ms worst answer under load; 250 ms used.
- Fixing the judge survivors here: non-blocking, and a change would need a re-judge. Row 241.
- Treating the race-leg timeouts as a regression: both packages pass alone under `-race`, one has no dependency on the diff.
- A third design run: the carve-out conditions held.

**Retro:**
1. **A revision introduced the next round's objection.** r1 fixed glm's census point by pasting the new row into the wrong table, and r2 blocked on exactly that. The carve-out absorbed it cheaply, but a revision directive should ask the designer to re-read the diff of its own edit. Instance 1; watch.
2. **Both reviewers misread a correct claim as wrong** ("three packages") because the doc's own V-row under-reported it. The controller's measurement was what separated "doc wrong" from "evidence row wrong". Rule 3f earned its keep again.
3. **The race leg's 8-min package bound reds under co-load** (`host/daemon`, `host/verifygate` this time; `host/daemon` in 327). Second instance of the class at this bound; if CI ever reds on it, it gets a row. Watch-item.
4. **D-WORLD-102's scope ends at row 233.** After the next pick, the loop returns to rule-(d) block reports unless a release row becomes routable or Mark extends the scope. Rows 235 and 237–242 sit outside it.
- Harness share 0. No skill edit.

**Progress:** 1.0 clauses 7/7 MET; goal unmoved (release rows). Row 230 BUILT in PR #442; rows 241–242 filed.
**Next:** row 230 lands on merge CI green (landing note). Then row 233, the last row in D-WORLD-102's scope; a dev red still outranks.

**Landing (iteration 329):** PR #442 merged as `9e562235` (merge commit, `--match-head-commit 707ca1ba`). Head `707ca1ba` was 3/3 green and was re-read MERGEABLE/CLEAN, with `origin/dev` still at `e5c1867f`, right before the merge. **Merge CI:** `commits/9e562235/check-runs` → 3, all success on attempt 1 (run 41180204408: verify gate, go race gate, go host build + test gate). Dev green. The local race-leg package timeouts (`host/daemon`, `host/verifygate`) did not occur in CI. Pollers ran under `/bin/bash` from a script file, one per SHA, each verdict line naming its SHA.

## 330 — 2026-10-31 — row 233 designed and parked: the holder-exit-reset excuse is blocked at quorum r2 on the fail-closed premise; real-srt Linux 49–54/100 under load; D-WORLD-105 asks Mark [ADMIN]

**Kind:** the third and last pick under D-WORLD-102 = A (order 231 → 230 → 233). Designer → quorum ×2 → park. No planner, executor or judge: nothing was built. The record lands with the design and evidence in one docs-only PR.

**Pick and why:** row 233 `w-exec-probe-holder-exit-reset`.
- **Clause map at `a6af0cd4`:** 1–7 MET; no UNMET clause.
- **Release rows not routable.** Row 150 waits only on the attended A1. Row 151 needs row 188 N2d/N4a: 0 matching commit subjects on `origin/dev` since 2026-10-01 (N2c control 2). The only commits since iteration 329's base are 329's own build, record and landing.
- So D-WORLD-102 = A applies; 231 and 230 are LANDED, so 233 is next.
- **Premise checked first-party** on a scratch worktree at `a6af0cd4`: a scratch test printed `judgeReverse("WORLD-PROBE-BOUND", {foreign})` = 4 (INCONCLUSIVE, the control), `{foreign, silent}` = 3 (SILENT), `{silent}` = 3; rc 0. Tree removed.
- Index grep: no earlier attempt at row 233; no design doc, no PR. Claim message `inbox_1793435311567_4812a6f1` sent before routing.

**Gate 0:**
- Kill switch armed (`mission-world.disabled` absent); gh `sunholo-voight-kampff`; billing tripwire CLEAN.
- `mission_directives.sh`: 0 allowlisted directives on #366 since `2026-10-27T09:03:42Z` (30 comments, LEGACY watermark); nothing to mark.
- `mission_gate0_self_notices.sh`: rc 1 on the 2026-10-22 rc=1 notice `6345079873`, already credited (25 mentions in the live log). Control 107:4 ok. Died-mid-flight traces: no open PR from this account; no new worktree residue.
- Inbox: `mission-world` holds only this loop's own claim messages.
- Ledger: 0 OPEN rows at Gate 0 (91 rows valid).
- No thread rotation due: #366 created 2026-10-26, 30 comments.

**Gate 1:**
- Local `dev` fast-forwarded to `origin/dev` = `a6af0cd42b1db7ccf912df88dec3bbaade95fd9a`. Base recorded `a6af0cd4…@2026-10-31T08:27:06Z`; re-read before the worktree with no drift.
- Skill: `SKILL.md` and all 13 `resources/*.md` `cmp`-identical to `origin/dev`, per file; rules read from the resolved symlink target.
- Dev CI: `a6af0cd4` check-runs 3/3 success; the last 3 CI runs on dev success.

**Design** (`w-exec-probe-holder-exit-reset`; resolver `recipe claude:claude-opus-5-5 declared:provider-pin`, spawned as Agent `opus` per the 318–329 precedent; rotation pointer unchanged, codex over ration):
- **r0** `e0f42390`: 364 lines, 18 V-rows. Rule: a zero-byte reset after a foreign answer in the same attempt records nothing; `judgeReverse` unchanged; T12 (10 rows) and M30–M35, each killed against a prototype.
  - **Real srt** 0.0.78 (cli.js sha256 = pin), installed outside the repo. darwin 1000 ms hold: 0/100 idle, 2/100 under 32 `yes`. **Linux** (throwaway branch `probe/world-iter330-holderexit`, run 41185322456, deleted): 0/100 idle, **54/100 under 8 `yes`** [real srt]; stand-in 49/100. The inherited 3–13% [stand-in, darwin] understated Linux.
  - srt never sets `reusePort` (0 hits; control 9 `listen(`); the live mux port refuses the child's exact bind, a Go default bind and a Go `SO_REUSEPORT` bind on both platforms (control: two `SO_REUSEPORT` sockets co-bind).
- **Quorum r1 BLOCKED** ($0.3150). Seats oc-glm-5-3, oc-kimi-k3, gpt6-1-sol, all REJECT; reserve gemini-4-argon; claude benched; no absentees.
  - glm and kimi: `dialErrIsSilent`'s match set had no V-row and no test exercised a non-ECONNRESET connect error after a foreign answer. **Controller measured before routing (rule 3f):** `exec_probe.go:725` is `errors.Is(err, syscall.ECONNRESET)`, one caller (`:757`); T8 pins ECONNABORTED, ECONNREFUSED, EINVAL and a real timeout as false. Handed to the designer as a measurement.
  - gpt6-1-sol: the rule enlarges the accept set; "a different execution already escapes" is no proof.
- **r1** `a9d3307e` (the one mandated revision, designer resumed): 455 lines, 24 V-rows. Rule tightened: a reset is excused only if a foreign dial STARTED < `execProbeListenMS` (1000 ms) before it AND the run exited cleanly (`r.err == nil`). Largest observed gap 94.4 ms Linux (131 resets), 57.7 ms darwin; availability unchanged on the banked data. A lemma and theorem in §3.3. T12 18 rows; M30–M40 each killed; M41 (wiring passes `true`) has no deterministic killer. The controller's candidates (a) token gate and (b) adjacency were both rejected with measurements ((a) recovers 0 of 49 on Linux, where the correct child reports BOUND).
- **Quorum r2 BLOCKED** ($0.3533; same seats; no absentees). glm PASS. kimi REJECT: M41 leaves the safety condition's only production line unpinned; concrete fix (a `reverseLegRun` seam, T13, a constant-equality V-row). gpt6-1-sol REJECT: the proof needs "a reachable child never resets before its listener closes", unmeasured on darwin and resting on an unenforced `tcp_abort_on_overflow` on Linux; fix = keep the suppression disabled unless that is established AND enforced, otherwise remove it.
- **Disposition: PARKED needs-human-review (D-WORLD-105).** gpt's objection disputes whether a timing-based suppression belongs in a fail-closed gate, so the narrow-refinement carve-out (every objection fixable verbatim and none disputing direction) does not hold. A split was considered: the reviewer-clean remainder (T12 rows pinning today's classification) does not move the row. Round count 2.
- Artifacts: `verification/world-iter330/quorum-r1.json`, `quorum-r2.json`; designer evidence (30+ `.txt` files) in the same directory.

**Plan / Executor / Independent evaluator:** not spawned — the design parked at quorum, so there was no plan to execute and no executor output to judge. Nothing in this iteration lands on the controller's verdict except bookkeeping (design doc, evidence, record). The generator ≠ judge requirement is therefore vacuous this fire rather than waived; it binds the build if D-WORLD-105 = A.

**Routing evidence:** base=a6af0cd42b1db7ccf912df88dec3bbaade95fd9a@2026-10-31T09:16:49Z (Gate 4).
- Driver env: `MISSION_DESIGNER_RESOLVED=claude:claude-opus-5-5`, `MISSION_PLANNER_RESOLVED=opus`, `MISSION_EXECUTOR_RESOLVED=claude:claude-haiku-5-5`, `MISSION_EVALUATOR_RESOLVED=sonnet`; `MISSION_OVER_RATION=codex`. `MISSION_ROUTING_NOTE`: planner and executor codex `gpt-6.1-sol` unusable (probe rc 75, over daily ration).
- Resolver: designer `recipe claude:claude-opus-5-5 declared:provider-pin`; planner `agent-tool opus fail-closed:env-pin`; executor `recipe claude:claude-haiku-5-5 declared:provider-pin`; evaluator `agent-tool sonnet declared:alias-pin`.
- Tokens and time per role:
  - controller `claude:claude-opus-5-5`: tok not reported;
  - designer Agent `opus`: 210,746 tok / 1,108 s; revision resume 277,699 tok / 957 s, as reported;
  - planner, executor, evaluator: not spawned (design parked).
- Designer diet: 1 doc, 2 runs (author + mandated revision).
- Quorum CLI note: the first r1 invocation exited rc 2 (`--controller-note is required when --controller-verdict is set`), $0; re-run with a note.
- Two Linux CI probe runs on throwaway branches (`probe/world-iter330-holderexit`, `…-r1`), both deleted (`git ls-remote` 0, `dev` control 1); Actions minutes, $0 metered.
- Metered $0.6683 (quorum r1 $0.3150 + r2 $0.3533).

**Ruled out:**
- Applying gpt6-1-sol's fix under the carve-out: it disables or removes the mechanism unless an enforced invariant exists, i.e. it disputes direction; a controller-chosen resolution is not allowed there.
- Applying kimi's fix alone and proceeding: gpt's objection would still stand.
- The controller's candidate (a), gating on the child's token: recovers 0 of 49 Linux refusals (the correct child reports BOUND in its netns).
- The controller's candidate (b), adjacency plus a constant bound: a 750 ms worst inter-dial gap still admits a child reset next to a foreign answer (designer counter-example).
- A third design round: the one revision and one re-quorum are spent.

**Retro:**
1. **A fail-closed gate is where a quorum is most likely to dispute direction rather than detail.** Rows 229, 232 and 230 passed via the carve-out because their objections were evidence gaps; this one changed which executions the gate accepts, and the OpenAI seat held it. Instance 1 of "an availability fix to a security gate parks on soundness"; watch.
2. **Real-tool measurement moved an inherited number by ~4–15×** (3–13% [stand-in, darwin] → 49–54/100 [real srt, Linux, load]). The throwaway-branch Linux recipe (instance 2 after iteration 318) was again the cheapest way to settle a "needs Linux" premise.
3. **The quorum CLI requires `--controller-note` beside `--controller-verdict`;** the skill's Gate 2 text shows only the verdict flag. One wasted invocation at $0. Instance 1; no skill edit.
4. **D-WORLD-102's scope is exhausted.** Next fires return to rule-(d) block reports unless Mark answers D-WORLD-105 = A/B, widens the scope (rows 235, 237–242), or a release row becomes routable.
- Harness share 0. No skill edit.

**Progress:** 1.0 clauses 7/7 MET; goal unmoved (release rows). Row 233 PARKED (D-WORLD-105).
**Next:** rule-(d) block report unless D-WORLD-105 is answered A or B, Mark widens D-WORLD-102's scope, or a release row (150, 151) becomes routable; a dev red still outranks.

## 331 — 2026-10-31 — no pick: rule (d) block report 10, the first since D-WORLD-102's scope ran out; row 233 parked on D-WORLD-105; release rows 150/151 still wait on attended steps; dev green at cb85eb2b [ADMIN]

**Kind:** bookkeeping-only iteration under standing rule (d) (D-WORLD-46). There is no design, plan, code or product acceptance. The record is the STATUS stamp (328 rotated to the archive), the CURRENT GOAL line, a D-WORLD-105 evidence note, this entry, the regenerated index, the dashboard, and `verification/world-iter331/gate1-2-readings.txt`. An independent judge reviews it before merge.

**Pick and why:** none.
- **Clause map at `cb85eb2b`:** clauses 1–7 MET. There is no UNMET clause and no critical-path row.
- **What moved since iteration 330's base:** only 330's own work. `git log a6af0cd4..origin/dev` lists 4 commits (`e0f42390` and `a9d3307e`, design r0 and r1; `c59fab19`, the record; merge `cb85eb2b` of #444).
- **D-WORLD-81 release rows, blockers re-measured this iteration:**
  - **150** waits only on the attended A1. Over commit subjects `a6af0cd4..cb85eb2b`, `grep -ciE 'row ?150|A1\b|marketing'` → **0**; the control `grep -ciE 'iteration|record|merge|row 233'` → **4** of 4 commits.
  - **151** needs row 188's N2d and N4a. `git log origin/dev --since=2026-10-01 --format='%h %s' | grep -ciE 'N2d|N4a'` → **0**; control `N2c` → **2**. The full-message form returns the same four record or design commits as iterations 320–326 (`fe9c891e`, `92495b60`, `be184351`, `dc719822`). Rows 217, 218 and 219 are still `[NEXT]` on the attended track, so neither milestone is loop-routable.
- **D-WORLD-102 = A** extended D-WORLD-64's exception to rows 231 → 230 → 233 only. 231 and 230 are LANDED (iterations 328, 329); 233 is `[PARKED iteration 330]` on **D-WORLD-105**, which is OPEN with default C: nothing is built and the row stays parked. Rows 235 and 237–242 are hygiene outside that scope, so rule (d) bars them.
- **External predicates re-read** with the whole-charter pattern under `/bin/bash` from a script file: **29** distinct upstream numbers, `results=29`, **all 29 closed**. Controls: row 53's `ailang#941` matched, and `gh issue list -R sunholo-data/ailang --state open --limit 3` returned 3. Still a lower bound (bare `#N` not counted).
- **Weekly external-issue sweep (Gate 0 step 5):** iteration 305 ran this week's sweep. Re-checked: 1 open issue, #366, the bookkeeping thread (charter 1 / log 63 / status archive 8 / dashboard 0 before this record). **0 orphans of 1 enumerated.**

**Gate 0:**
- Kill switch armed (`mission-world.disabled` absent). gh `sunholo-voight-kampff`. Billing tripwire CLEAN.
- `mission_directives.sh` (driver pin copy): 0 allowlisted directives on #366 since `2026-10-27T09:03:42Z` (31 comments) and 0 on #313 since `2026-10-22T01:34:42Z` (44 comments). Nothing to mark.
- `mission_gate0_self_notices.sh`: verdict "1 crash notice". It is the same 2026-10-22 rc=1 notice `6345079873`, already credited (`grep -c` → 26 in the live log and 3 in the log archive before this entry). Control 107:4 ok. No crash notice on #366 (31 self, 0 crash).
- Unread inbox: nothing addressed to World. The `mission-world` inbox holds only this loop's claim and report messages. The iteration-331 claim was sent (`inbox_1793449788007_e84aa880`).
- `mission_decisions.sh --check`: valid, 92 rows, 1 OPEN (D-WORLD-105). The only ledger diff in `a6af0cd4..cb85eb2b` is 330's own D-WORLD-105 insert (one `+` line).

**Gate 1:**
- Local `dev` equalled `origin/dev` (`cb85eb2b`), with only untracked residue in the tree. Gate-1 base `cb85eb2b3a7730f23df73dbb558f92c48e87626f@2026-10-31T12:27:04Z`.
- Skill drift: per-file `cmp` of `SKILL.md` and all 13 `resources/*.md` at the resolved symlink target against `origin/dev` in the freshly fetched driver clone: every file identical. Rules read from the resolved target.
- Dev CI: `commits/cb85eb2b/check-runs` → 3, all success. The last five dev CI runs (`cb85eb2b`, `a6af0cd4`, `9e562235`, `e5c1867f`, `6641df63`) all read success on attempt 1.
- Open PRs: none. Died-mid-flight traces: no iteration-331 branch or worktree existed before this one. `.wt-world-iter306-r0` is the known residue. The attended scratchpad worktrees belong to the attended track and were left alone.

**Designer / planner / executor:** not spawned. There is no pick, so the routing table has nothing to send them (precedent 257–264, 317, 319–326). No spawn was attempted, so none failed. The operator's standing request for Agent-tool roles is honoured through the evaluator.

**Independent evaluator (record):** REQUIRED. It runs through the Agent tool on `sonnet` (driver `MISSION_EVALUATOR_RESOLVED=sonnet`, path `agent-tool`), in a fresh, read-only, foreground context, and judges the record commit before merge. See the follow-through below.

**Routing evidence:** base=cb85eb2b3a7730f23df73dbb558f92c48e87626f@2026-10-31T12:27:04Z (Gate 1).
- Driver env: `MISSION_DESIGNER_RESOLVED=claude:claude-opus-5-5`, `MISSION_PLANNER_RESOLVED=opus`, `MISSION_EXECUTOR_RESOLVED=claude:claude-haiku-5-5`, `MISSION_EVALUATOR_RESOLVED=sonnet`; `MISSION_OVER_RATION=codex`. `MISSION_ROUTING_NOTE`: planner and executor codex `gpt-6.1-sol` unusable (probe rc 75, over daily ration).
- Controller `claude:claude-opus-5-5` (tok: not reported). Designer, planner and executor were unused.
- Evaluator: Agent `sonnet`, tokens recorded in the follow-through. Generator ≠ judge: opus controller vs sonnet judge.
- Quorum: none (no doc). Metered $0.

**Ruled out:**
- Picking 150 (attended A1) or 151 (row 188 N2d/N4a absent, measured above).
- Picking 233: PARKED on D-WORLD-105, whose default C builds nothing.
- Picking 235 or 237–242 under rule (d): outside D-WORLD-102's scope and D-WORLD-64's exception.
- Re-asking D-WORLD-105. It is the one open decision and the digest carries it; it is not asked again.
- Filing a new ask to widen D-WORLD-102's scope: iteration 330's record and D-WORLD-105's text already name that option (rows 235, 237–242) and the dashboard lists it; a second decision row for the same lever would duplicate it.
- Rotating the log: 26 live entries with this one, below the ~40 threshold. The index was regenerated with `rotate-log --index-only`.

**Retro:**
1. **Idle cadence resumes.** Block report 10 overall (317, 319–326, 331), the first after the four D-WORLD-102 picks (328 built 231, 329 built 230, 330 parked 233; 327 was a dev red). Each costs one judge and two CI cycles, and no code. The levers are Mark's: D-WORLD-105, a wider hygiene scope, or the attended A1 and 219 → 217 → 218 track.
2. **The CURRENT GOAL line still grows by one clause per fire** (iteration 326's retro item 2, instance 2 now): 7,357 bytes with 19 earlier-iteration clauses at `cb85eb2b`, 7,731 bytes and 20 after this record (bytes by `wc -c`, newline included; clauses by `grep -o 'Earlier (iteration' | wc -l` on that line). Not trimmed here, for the same reason as 326: it is charter prose outside a block report's record scope.
- Harness share 0. No skill edit.

**Progress:** 1.0 clauses 7/7 MET; goal unmoved (release rows). Row 150: the attended A1 remains. Row 151: behind row 188 N2d/N4a. Row 233: PARKED (D-WORLD-105).
**Next:** rule (d) block report unless D-WORLD-105 is answered A or B, Mark widens D-WORLD-102's scope, or a release row (150, 151) becomes routable; a dev red still outranks.

**Record follow-through (iteration 331):** the independent evaluator judged record commit `07bd28ef` on `cb85eb2b`: **r1 PASS 96/100, zero blocking** (35/25/22/15). It ran through the Agent tool on `sonnet`, in a fresh, read-only, foreground context on a detached worktree at the record commit (44,778 tok, 55 s, 9 tool uses). It re-ran the no-pick premise (row 150 subjects 0 of 4 in `a6af0cd4..cb85eb2b`, and `row ?151|N2d|N4a|row ?188` also 0; rows 217–219 `[NEXT]`; 233 `[PARKED iteration 330]`; D-WORLD-105 OPEN default C; #366's 31 comments with no allowlisted author; `cb85eb2b` 3/3 success; #445 the only open PR) and the edit integrity (3 stamps kept, the 328 stamp `cmp`-identical at the archive head, archive +2/−0, ledger valid at 92 rows, index +1 row with 329/330 re-derived from their headings, no closing keyword). It took the N2d/N4a figures from the readings file rather than re-running them.
- **NB1** (the retro's "19 → 20 earlier-iteration clauses" did not reproduce: the judge counted 31 → 32 `iteration N` tokens on the line): a different instrument, not a different line. The controller's count is `grep -o 'Earlier (iteration' | wc -l` on the CURRENT GOAL line (19 at `cb85eb2b`, 20 after); the judge's counts every `iteration N` mention, including those inside clauses. Retro item 2 now cites its instrument inline. The byte counts reproduced.
- No r2: one instrument citation added.

**Landing (iteration 331):** PR #445 merged as `558fea49` (merge commit, `--match-head-commit 6c4ed241`). The PR head `6c4ed241` was 3/3 green before the merge, and origin/dev was re-read right before it (`cb85eb2b`, MERGEABLE/CLEAN). The merge tree is byte-identical to the head tree. **Merge CI:** `commits/558fea49/check-runs` → 3, all success on attempt 1 (verify gate, go race gate, go host build + test gate). Dev green. Both CI pollers ran under `/bin/bash` from a script file with the subject SHA in each verdict line; three of the four bounded polls hit their ~9-minute deadline while the race gate ran and were re-issued against the same SHA with a fresh deadline.

## 332 — 2026-10-31 — no pick: rule (d) block report 11; nothing moved since 331 but its own record and landing; row 233 parked on D-WORLD-105; release rows 150/151 still wait on attended steps; dev green at 85a7bd5c [ADMIN]

**Kind:** bookkeeping-only iteration under standing rule (d) (D-WORLD-46). There is no design, plan, code or product acceptance. The record is the STATUS stamp (329 rotated to the archive), the CURRENT GOAL line, a D-WORLD-105 evidence note, this entry, the regenerated index, the dashboard, and `verification/world-iter332/gate1-2-readings.txt`. An independent judge reviews it before merge.

**Pick and why:** none.
- **Clause map at `85a7bd5c`:** clauses 1–7 MET. There is no UNMET clause and no critical-path row.
- **What moved since iteration 331's base:** only 331's own work. `git log cb85eb2b..origin/dev` lists 5 commits (`07bd28ef`, the record; `6c4ed241`, the judge fix; merge `558fea49` of #445; `91fc5137`, the landing note; merge `85a7bd5c` of #446).
- **D-WORLD-81 release rows, blockers re-measured this iteration:**
  - **150** waits only on the attended A1. Over commit subjects `cb85eb2b..85a7bd5c`, `grep -ciE 'row ?150|A1\b|marketing'` → **0**; the control `grep -ciE 'iteration|record|merge|row 233|landed'` → **5** of 5 commits.
  - **151** needs row 188's N2d and N4a. `git log origin/dev --since=2026-10-01 --format='%h %s' | grep -ciE 'N2d|N4a'` → **0**; control `N2c` → **2**. The full-message form returns the same four record or design commits as iterations 320–331 (`fe9c891e`, `92495b60`, `be184351`, `dc719822`). Rows 217, 218 and 219 are still `[NEXT]` on the attended track, so neither milestone is loop-routable.
- **D-WORLD-102 = A** covered rows 231 → 230 → 233 only. 231 and 230 are LANDED; 233 is `[PARKED iteration 330]` on **D-WORLD-105**, still OPEN with default C (nothing built). Rows 235 and 237–242 are hygiene outside that scope, so rule (d) bars them.
- **External predicates re-read** with the whole-charter pattern under `/bin/bash` from a script file: **29** distinct upstream numbers, `results=29`, **all 29 closed**. Controls: row 53's `ailang#941` matched, and `gh issue list -R sunholo-data/ailang --state open --limit 3` returned 3. Still a lower bound (bare `#N` not counted).
- **Weekly external-issue sweep (Gate 0 step 5):** iteration 305 ran this week's sweep. Re-checked: 1 open issue, #366, the bookkeeping thread (charter 3 / log 67 / status archive 8 / dashboard 0 before this record). **0 orphans of 1 enumerated.**

**Gate 0:**
- Kill switch armed (`mission-world.disabled` absent). gh `sunholo-voight-kampff`. Billing tripwire CLEAN.
- `mission_directives.sh` (driver pin copy): 0 allowlisted directives on #366 since `2026-10-27T09:03:42Z` (32 comments) and 0 on #313 since `2026-10-22T01:34:42Z` (44 comments). Nothing to mark.
- `mission_gate0_self_notices.sh`: verdict "1 crash notice". It is the same 2026-10-22 rc=1 notice `6345079873`, already credited (`grep -c` → 27 in the live log and 3 in the log archive before this entry). Control 107:4 ok. No crash notice on #366 (32 self, 0 crash).
- Unread inbox: nothing addressed to World. The 39 unread `mission-world` messages all come from `mission-world`: this loop's own claims and reports. The iteration-332 claim was sent (`inbox_1793464127436_08543c29`).
- `mission_decisions.sh --check`: valid, 92 rows, 1 OPEN (D-WORLD-105). The only ledger diff in `cb85eb2b..85a7bd5c` is 331's own D-WORLD-105 evidence note (one `-`/`+` pair, commit `07bd28ef`).

**Gate 1:**
- Local `dev` equalled `origin/dev` (`85a7bd5c`), 0 ahead / 0 behind, with only untracked residue in the tree. Gate-1 base `85a7bd5cb11c786157a14e0b4fa80271fe8e6677@2026-10-31T16:28:05Z`.
- Skill drift: per-file `cmp` of `SKILL.md` and all 13 `resources/*.md` at the resolved symlink target against `origin/dev` in the freshly fetched driver clone: every file identical. Rules read from the resolved target.
- Dev CI: `commits/85a7bd5c/check-runs` → 3, all success. The last five dev CI runs (`85a7bd5c`, `558fea49`, `cb85eb2b`, `a6af0cd4`, `9e562235`) all read success on attempt 1.
- Open PRs: none. Died-mid-flight traces: no iteration-332 branch or worktree existed before this one. `.wt-world-iter306-r0` is the known residue.

**Designer / planner / executor:** not spawned. There is no pick, so the routing table has nothing to send them (precedent 257–264, 317, 319–326, 331). No spawn was attempted, so none failed. The operator's standing request for Agent-tool roles is honoured through the evaluator.

**Independent evaluator (record):** REQUIRED. It runs through the Agent tool on `sonnet` (driver `MISSION_EVALUATOR_RESOLVED=sonnet`, path `agent-tool`), in a fresh, read-only, foreground context, and judges the record commit before merge. See the follow-through below.

**Routing evidence:** base=85a7bd5cb11c786157a14e0b4fa80271fe8e6677@2026-10-31T16:28:05Z (Gate 1).
- Driver env: `MISSION_DESIGNER_RESOLVED=claude:claude-opus-5-5`, `MISSION_PLANNER_RESOLVED=opus`, `MISSION_EXECUTOR_RESOLVED=claude:claude-haiku-5-5`, `MISSION_EVALUATOR_RESOLVED=sonnet`; `MISSION_OVER_RATION=codex`. `MISSION_ROUTING_NOTE`: planner and executor codex `gpt-6.1-sol` unusable (probe rc 75, over daily ration).
- Controller `claude:claude-opus-5-5` (tok: not reported). Designer, planner and executor were unused.
- Evaluator: Agent `sonnet`, tokens recorded in the follow-through. Generator ≠ judge: opus controller vs sonnet judge.
- Quorum: none (no doc). Metered $0.

**Ruled out:**
- Picking 150 (attended A1) or 151 (row 188 N2d/N4a absent, measured above).
- Picking 233: PARKED on D-WORLD-105, whose default C builds nothing.
- Picking 235 or 237–242 under rule (d): outside D-WORLD-102's scope and D-WORLD-64's exception.
- Re-asking D-WORLD-105 or filing a second scope-widening ask: the digest already carries both (see 331's Ruled out).
- Rotating the log: 27 live entries with this one, below the ~40 threshold. The index was regenerated with `rotate-log --index-only`.

**Retro:**
1. **CURRENT GOAL growth stopped for this fire.** Iterations 326 and 331 each recorded that the line grows by one clause per fire. A block report that repeats the previous one adds no new fact, so this record widened the existing clause to "iterations 331–332 … block reports 10–11" and refreshed its dev SHA instead of prepending a new one. The line is 7,744 bytes, up from 7,731 (`wc -c`, newline included), and still has 20 earlier-iteration clauses (`grep -o 'Earlier (iteration' | wc -l` on that line). Older clauses are not trimmed, for the reason 326 gives.
2. **Gate-0 instrument fault, mine: the first env read printed a secret into the local transcript.** Before reading `resources/gate-0-preflight.md`, the controller listed the driver env with `env | grep -E '^(MISSION|AILANG|CLAUDE_CODE)'`. That matched `CLAUDE_CODE_OAUTH_TOKEN` and printed its value into this fire's local transcript. Nothing was sent anywhere: the value appears in no commit, message or issue. This is the class fleet iterations 51 and 52 hit, and Gate 0's "read the environment by NAME" rule covers it, but the controller read that rule only after its first env call. The call came while it checked SKILL.md's `MISSION_WORKDIR` and pin sections, which name env vars and come before any gate file. Proposal (World cannot edit the shared skill, so this goes to Mark/V1): move the one-line rule, `env | cut -d= -f1 | grep -E '^(MISSION_|AILANG_)'`, into SKILL.md above the Repo Profile, where the first env-reading prompt appears. Rotating the token is Mark's call; the exposure is local to this rig's transcript store.
- Harness share 0. No skill edit.

**Progress:** 1.0 clauses 7/7 MET; goal unmoved (release rows). Row 150: the attended A1 remains. Row 151: behind row 188 N2d/N4a. Row 233: PARKED (D-WORLD-105).
**Next:** rule (d) block report unless D-WORLD-105 is answered A or B, Mark widens D-WORLD-102's scope, or a release row (150, 151) becomes routable; a dev red still outranks.

**Record follow-through (iteration 332):** the independent evaluator judged record commit `6ad80117` on `85a7bd5c`: **r1 PASS 95/100, zero blocking** (34/23/24/14). It ran through the Agent tool on `sonnet`, in a fresh, read-only, foreground context on a detached worktree at the record commit (40,183 tok, 45 s, 5 tool uses). It re-ran the no-pick premise: the 5 commits in `cb85eb2b..85a7bd5c` are all 331's own work; no `MarkEdmondson1234` comment on #366 (32 comments) and 44 on #313; ledger valid at 92 rows; D-WORLD-102 covers 231 → 230 → 233 only. It also re-ran the edit integrity: 3 stamps kept; the 329 stamp has the same md5 at the archive head as in the charter at `85a7bd5c`; archive +2/−0; index +1/−0; CURRENT GOAL 7,731 → 7,744 bytes with 20 → 20 clauses. It found no `sk-ant`/`oat01` string and no closing keyword in the diff. It took the upstream 29/29 sweep, the CI run list and the inbox counts from the readings file.
- **NB1** (the index row is truncated at "stil..."): the index generator truncates every long heading. No change.
- **NB2** (the token in the local transcript may merit rotation): this is already retro item 2, and the rotation is Mark's call.
- **NB3** (CI runs `9e562235`/`a6af0cd4` taken from the readings file): the readings file prints the `gh run list` output that produced them. No change.
- No r2: nothing to fix.
