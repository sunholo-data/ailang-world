# Ailang World Mission — iteration log

Append-only. One entry per outer-loop iteration. Newest at the bottom.

---

> **Older entries are ARCHIVED.** This file holds the newest 20. The full record of every
> iteration is in `world-mission-log-archive.md`, and a one-line index of ALL of them —
> the thing to grep before picking work, so the loop never repeats itself — is in
> `world-mission-index.md`.

## 431 — 2026-11-25 — row 243 V4 sprint 5 LANDED (PR #657, merge `ab520385`): terminal provision exclusion on `session new` and signed `approve`, a shared read-only preset-start read seam, and the S4c exact-artifact confirm replay; atomic mint split to S5b; nothing dispatches; 743 lines; judge sonnet PASS 86 [PRODUCT]

**Pick:** row 243 V4 sprint 5, the Next named by record 430.
- **Re-measured at pick** (`origin/dev` `cd06d950`): no V4 S5 plan, branch, PR or worktree (`git log origin/dev --grep S5` shows only V3 S5); product diff `fa3ad190..cd06d950` empty; quorum reader rc 0 (r2 `…T10-24-51Z` PROCEED 3/3, `absent_reviewers` = `[]`); design unchanged since `8dc91744`; ledger valid, 0 OPEN of 107. Carried premises re-measured: `provisionEpisode(` at `approve.go:398`/`session_new.go:649` without `terminalEpisodeLock`; `grep -c operator-episode docs/QUICKSTART.md` = 0.
- **Clause map:** 1.0 clauses 1–7 MET, none UNMET. Row 243 serves D-WORLD-94 and moves no clause; it is the only loop-routable build row.
- Claim sent to `controlplane` before routing.

**Gate 0:** stamped first; `gh` = `sunholo-voight-kampff`; billing CLEAN; kill switch absent. 0 allowlisted directives on #635 (8 comments). Step 6a rc 0 (0 crash, 0 transcript notices on #635/#574; `-prev` watermark absent → single-watermark read). No unread message addressed to `mission-world` as a directive (global unread: iteration 430's own report, fleet iteration 137, the AILANG v0.61.0 release note).

**Gate 1:** local `dev` `af7a610c` 22 behind `origin/dev` `cd06d950` — read state from origin; did not touch the main checkout (Mark's untracked drafts). Base `cd06d950…@2026-11-24T23:17:50Z` (a first `mission-base.sh record` ran from the driver clone's CWD and recorded an ailang SHA; re-run from the World clone — the known CWD trap). Skill drift 0 of 28 file copies. Dev check-runs 3/3 green.

**Planner** (`codex:gpt-6.2-sol`, workspace-write, background, 90-min cap, ~31 min (23:20 → 23:51 UTC), **397,781 tok**; resolver `recipe codex:gpt-6.2-sol declared:planner-lane-default-pin`, `derive-planner-lane.sh` same; probe rc 0 from the sprint worktree). Plan `w-world-address-bar-verbs-v4s5-sprint-plan.md` + JSON committed `6967f0c2`. It took design §8's named S5 split point by measurement: the compiled read seam is 743, design §8 allots 510 raw (689 calibrated) to atomic mint and recovery, and 743 + 689 = 1,432 > the 1,400 cap before the resolver and drain gates — so S5 = M1 163 / M2 730 / M3 743 (`git diff -M --numstat cd06d950 HEAD -- cmd host docs app`) and **S5b** carries the atomic minted-start proof. **79/79 named mutant kills** (the iteration-429 retro rule held a second time); 77 whole-package CLI rows UNMEASURED IN SANDBOX. Premises found false (§1.4): F2 — S4b judge finding 2 was incomplete: `service runner cancel` exists for unstarted actions, but nothing resolves an owned uncertain row, and status can report ready while that row keeps its owner; F4 — no paired-process consumer or drain exists; F5 — the replay early-return reads only the digest; F6 — the read seam exposes existing terminal preflight only, signed-input comparison is S5b.
- **Controller round trip** (`ctl_roundtrip.sh M1 M2 M3`) on a fresh `git archive cd06d950`: every `diff -r` vs `st-M*` rc 0.
- **False-green (5) pre-check:** the two new `os.WriteFile` calls write a fake agent under `t.TempDir()` and fixture files under the test's own session/task paths.

**Executor** (`codex:gpt-6.2-sol`, 45-min cap, 3 min, **64,797 tok**; resolver `recipe codex:gpt-6.2-sol declared:provider-pin`). M1→M3 each `ACCEPT_OK` (PASS 23/78/106, 0 FAIL, 0 SKIP); snapshots byte-equal to the planner's stages (7/8/10 files); `commit_ms.sh` re-ran each acceptance outside the sandbox → `613b6989` (163), `113ce3b5` (730), `2ddfb651` (743); committed tree = `st-M3`.

**Gates** on `2ddfb651` (plan §6, outside the sandbox, 1,267 s): **every leg rc 0** — `pkg-cmd-ailang-worldd` 235 s, `race-cmd-ailang-worldd` 363 s, store/authority race, inventories, routes, scope, callers, `verify_go` (plain) 535 s. Observation of the real LaunchAgents, `~/.ailang-world` and the launchctl count: identical before and after. This time the judge was told to wait on the sweep's `ALL_DONE` marker before running Go, so no load-only red appeared (iteration 430's retro instance 1).

**Evaluator r1** (Agent `sonnet`, foreground, ~44 min, 98,591 tok; resolver `agent-tool sonnet declared:alias-pin`; own worktree `.wt-world-iter431-row243-v4s5-eval`; directive posed questions — provisioning-path completeness and ordering, what `approve` consumes before refusal, the seam's parity and read-onlyness, replay ordering, the S5/S5b split, Linux portability — without controller verdicts): **PASS 86**, 0 blocking (correctness 26/30, tests 24/30, security 18/20, standards 9/10, scope 9/10). It ran 20 of its own mutants (must-die `J-LOCK-DISABLE` killed by `TestOperatorProvisionEpisodeExclusion`, must-live `J-NEUTRAL-CTX` survived) and re-measured 17 of the plan's sandbox-unmeasured rows, all killed behaviourally; per-milestone non-vacuity held (M1 lock fallback, M3 byte check, M2 terminal-path breaks). Non-blocking findings: (1) `session mint` still mints without the lock — **controller-reproduced** (`session.go` opens the store and calls `authority.Mint(mintCtx, st, *episode, …)` ~`:299`; 0 `terminalEpisodeLock` hits in `session.go`); latent while nothing can own an episode, carried to S5b/S6; (2) the parity tests partly compare the new code with itself — the unchanged pre-existing preset tests are the real regression check, and two mutants survive the whole package (`J-PRINT-AGENT` blanks the printed `agent:` line, `J-DOUBLE-OVERRIDE-APPEND` duplicates persisted overrides); (3) no test reads the QUICKSTART sidecar sentence (`J-QUICKSTART-b`); (4) the replay mutants die only in `host/authority`, as the plan discloses.

**Gate 3b:** PR #657 head `2ddfb651` attempt 1 **2/3** — `go race gate` failed on row 267's signature only: `FAIL host/daemon 480.105s` running `TestEpisodePolicyProtectedRebased/basename`, `host/verifygate` 464.0 s and `host/broker` 413.2 s in the same run (slow runner); `git diff --stat cd06d950 2ddfb651 -- host/daemon` empty; the sprint's own `cmd/ailang-worldd` race 309 s ok. `gh run rerun 45520125603 --failed` → attempt 2 **3/3 success**. `origin/dev` unmoved; PR `MERGEABLE CLEAN`; merged `--match-head-commit 2ddfb651` → `ab520385`; merge CI 3/3 success. **LANDED.**

**Routing evidence:** base=cd06d9500ec2d34e975d7566f8bbdb14b1529527@2026-11-24T23:17:50Z (Gate 1; sprint worktree created from the same SHA). Gate 4 base=ab520385a5f95b1a4d448121ce610b9f346e2883@2026-11-25T01:11:03Z (drift = this merge only).
- Controller `claude:claude-opus-5-5` (tok: not reported) · designer not spawned (plan-only fire on a reviewed design; rotation pointer unchanged) · planner `codex:gpt-6.2-sol` (397,781 tok) · executor `codex:gpt-6.2-sol` (64,797 tok) · evaluator Agent `sonnet` (98,591 tok). **Every required role spawned; none failed to spawn.**
- The run prompt asked for the Agent tool for every role "exactly as the skill's routing table specifies"; the resolver named the cross-provider codex recipe for planner and executor (the Agent tool is not a valid path for a `recipe` role), so those ran under the codex recipe and the evaluator used the Agent tool. No role ran inline; every code change was authored by the executor and judged by the evaluator.
- Judge independence: **cross-vendor** (OpenAI planner + executor, Anthropic judge).
- Metered: **$0.00** (no quorum, no metered lane).

**Progress:** 1.0 clauses 1–7 MET; goal unmoved. Row 243 V4: S1, S2, S2b, S3, S4, S4b, S5 landed. Every terminal provisioning path except `session mint` now refuses an operator-owned episode; the read stage of an operator start exists and is shared with the terminal; confirm replay is exact. S5b (atomic minted start) is next.

**Ruled out:**
- **Treating the PR's race red as the sprint's:** row 267's exact signature (8-minute `host/daemon` package timeout, slow runner by the same run's verifygate/broker times), 0 changed lines in `host/daemon`, the sprint's own package race green in the same job, rerun green.
- **A round 2 for the three survivors:** each is a test gap on output the design does not make an authority gate (a printed line, a duplicated override record, a doc sentence); carried to S5b rather than an unplanned round.
- **A round 2 for `session mint`:** design §2.2 and the S5 plan name `provisionEpisode` callers; `session mint` adds no worktree and nothing can own an episode until a runner consumer lands, so it joins S5b/S6's pre-dispatch gates.
- **Building atomic mint in this sprint:** measured over the cap (1,432 > 1,400) before the resolver and drain; design §8 names S5 as the split point.

**Retro:**
- **The planner retro rule held again** (79/79 named kills; 0 blocking). Instance 2 of the rule working; keep it.
- **Sequencing the judge behind the gate sweep worked** (iteration 430's retro, instance 1): no load-only red in `verify_go` this time.
- **Row 267 hits sprint PRs whose changes exclude `host/daemon`** (instances 2 and 5, #651 and #657; five instances in ~9 h of run creation): every rerun costs ~13 min of CI wall time. It stays a rule-(d) hygiene row; the count is the signal for Mark.

**Next:** V4 S5b planner — design §5.3's atomic minted-start rows per the S5 plan's Appendix C (signed revalidation, complete collision checks, atomic mint, token/task/MCP file recovery, attended resolver before any new uncertain outcome, paired-process lock consumer/drain), carrying `session mint`'s lock gap, the judge's three survivors, S3b and J-OWN-UID.

**Gate 5 bookkeeping:** rotation not due (#635 opened 2026-11-23T07:34Z; next 2026-11-30). Weekly external-issue sweep not due (424 ran it for this week). Idle fingerprint NOT recorded (this fire picked and landed). Harness share last 20: 0/20 (20 tagged). Last 3 landings (431, 430, 429) move no clause; no clause is UNMET → no DRIFT line. Log not rotated (21 entries, under the ~40 threshold); index regenerated.

**Record judge** (Agent `sonnet`, ~1.5 min, 74,050 tok): **FAIL** on two facts in the first draft (`53b08d0f`), both in my row 267 text — "five instances in about thirty hours" (the span is ~9 h: run 45425824922 created 2026-11-24T15:28Z, run 45520125603 2026-11-25T00:42Z) and "instance 5 the first on a sprint whose package list excludes `host/daemon`" (PR #651, instance 2, also changed only `cmd/ailang-worldd`, `host/store` and docs). Both corrected in the follow-up commit, plus the dashboard's "lock gap closed" scoped to `session new`/`approve` and row 267's `TestWorkspacePackageCache*` sentence scoped to instances 1–3 (it was already inaccurate for instance 4).

## 432 — 2026-11-25 — row 243 V4 sprint 5b LANDED (PR #659, merge `5d3b0581`): the terminal `session mint` takes the per-episode lock, an atomic operator mint (sealed runner-recheck proof; start, credential and observed receipt in one transaction) with no production caller, and the attended resolver `service runner reconcile-mint`; worktree, files and the paired-process consumer split to S5c; nothing dispatches; 1,400 lines; judge sonnet PASS 84 [PRODUCT]

**Pick:** row 243 V4 sprint 5b, the Next named by record 431.
- **Re-measured at pick** (`origin/dev` `f8d8e3ad`): no V4 S5b plan, branch, PR or worktree; product diff `ab520385..f8d8e3ad` empty; quorum reader rc 0 (r2 `…T10-24-51Z` PROCEED 3/3, `absent_reviewers` = `[]`); design unchanged since `8dc91744`; ledger valid, 0 OPEN of 107. Carried premise re-measured: `session.go:299` calls `authority.Mint` with 0 `terminalEpisodeLock` hits in `session.go`.
- **Clause map:** 1.0 clauses 1–7 MET, none UNMET. Row 243 serves D-WORLD-94 and moves no clause; it is the only loop-routable build row.
- Claim sent to `controlplane` before routing.

**Gate 0:** stamped first; `gh` = `sunholo-voight-kampff`; billing CLEAN; kill switch absent. 0 allowlisted directives on #635 (9 comments). Step 6a rc 0 (0 crash, 0 transcript notices on #635/#574; `-prev` watermark absent → single-watermark read). `mission-world` inbox: 0 unread.

**Gate 1:** local `dev` `af7a610c` behind `origin/dev` `f8d8e3ad` — read state from origin; did not touch the main checkout (Mark's untracked drafts). Base `f8d8e3ad…@2026-11-25T03:18:33Z`, recorded from the World clone with the driver script's absolute path. Skill drift 0 (SKILL.md + 13 resources). Dev check-runs 3/3 green.

**Planner** (`codex:gpt-6.2-sol`, workspace-write, background, 90-min cap, ~63 min (03:21 → 04:24 UTC), **479,066 tok**; resolver `recipe codex:gpt-6.2-sol declared:planner-lane-default-pin`, `derive-planner-lane.sh` same; probe rc 0 from the sprint worktree). Plan `w-world-address-bar-verbs-v4s5b-sprint-plan.md` + JSON committed `8f7764ec`. M1 95 / M2 1,400 cumulative (M2 alone 1,305; `git diff -M --numstat f8d8e3ad HEAD -- cmd host docs app`; exactly the cap). It split again by measurement: worktree reservation, no-overwrite task/token/MCP publication, full live-input reconstruction and the paired-process consumer/drain go to **S5c** (measured floor 378 lines including the 132-line S3b bank; full size unmeasured). **123/123 named mutant kills** (the iteration-429 retro rule held a third time); 44 whole-package CLI rows and 3 legacy rows UNMEASURED IN SANDBOX. Premises (§1.4): F1 a legacy opaque episode cannot become an operator action; F2 re-measured on a base export — `service runner cancel` leaves an owned uncertain row unchanged; F7 S5 supplies `sha256:`-tagged profile and task digests.
- **Controller round trip** (`ctl_roundtrip.sh M1 M2`) on a fresh `git archive f8d8e3ad`: every `diff -r` vs `st-M*` rc 0.
- **False-green (5) pre-check:** every new `os.WriteFile`/`os.Remove`/`os.Chmod` is in a `_test.go` under a fixture root; the production token read is `O_RDONLY|O_NOFOLLOW|O_NONBLOCK`; the gate script's `observe()` only reads.

**Executor** (`codex:gpt-6.2-sol`, 45-min cap, 2 min, **59,843 tok**; resolver `recipe codex:gpt-6.2-sol declared:provider-pin`). M1, M2 each `ACCEPT_OK` (PASS 16/132, 0 FAIL, 0 SKIP); snapshots byte-equal to the planner's stages (2/15 files); `commit_ms.sh` re-ran each acceptance outside the sandbox → `b01125ee` (95), `6e354749` (1,400); committed tree = `st-M2`.

**Gates** on `6e354749` (plan §6, outside the sandbox, 1,259 s): every leg rc 0 — `pkg-cmd-ailang-worldd` 240 s, `race-cmd-ailang-worldd` 347 s, store/authority/approval race, inventories, routes, `mint-quickstart`, scope, callers — **except `verify_go` rc 1** (545 s) on `TestSeToolsRunBenchmarkTasksEndToEnd` in `host/daemon` ("the mock saw [GET, POST] … want exactly the one POST"). `git diff --stat f8d8e3ad HEAD -- host/daemon` is empty, the test passes alone at head (5.52 s), and it is on the directive's load-sensitive list (red in iteration 430's sweep too). Observation of the real LaunchAgents, `~/.ailang-world`, the real credentials store and the launchctl count: identical before and after.

**Evaluator r1** (Agent `sonnet`, foreground, ~40 min, 114,085 tok; resolver `agent-tool sonnet declared:alias-pin`; own worktree `.wt-world-iter432-row243-v4s5b-eval`; directive posed questions — `session mint` ordering, transaction atomicity, proof forgeability, binding of every signed field, the resolver's revoke safety on today's installations, inertness, the S5b/S5c split, Linux portability — without controller verdicts; sequenced behind the sweep's `ALL_DONE` marker): **PASS 84**, 0 blocking (correctness 26/30, tests 24/30, security 18/20, standards 9/10, scope 7/10). It re-measured 17 of the plan's UNMEASURED CLI rows, all killed behaviourally; must-die L-FENCE killed, must-live `0600`→`0o600` survived; M1 non-vacuous against its own `session.go` hunk. Non-blocking findings: (1) untested defence-in-depth guards — the controller audited its 76-line ledger: the first batch ran mint mutants against `./host/store` while the tests live in `host/authority` (judge-flagged); all were re-run correctly and killed except `E-TX-STORE-ID`, `E-TX-KIND` and `E-RES-STORE-ID` (whole-package SURVIVED) and `E-RES-EPISODE-CONFLICT`, `E-RES-REREAD`, `E-RES-REGULAR`, `E-RES-TOKENPATH-BASE` (named/broad runs only). **Controller-reproduced `E-TX-KIND`:** dropping `|| row.Kind != "start"` at `host/store/operator_mint.go:115` leaves `go test ./host/authority/ ./host/store/` green, and neither `proof.Bind` nor `decodeOperatorMintAction` requires `start` — so it is the transaction's only start-kind gate, correct and untested. All seven are latent while `CommitOperatorMint`'s only caller, `authority.MintOperatorStart`, has no production caller; (2) no explicit test of `reconcile-mint` refusing on a machine with no installed runner (the judge read it refusing before any write); (3) Appendix C carriage of J-PRINT-AGENT, J-DOUBLE-OVERRIDE-APPEND and J-OWN-UID unconfirmed — **closed by the controller** (Appendix C bullets "Legacy parity survivors, 8 measured fixture lines" and "J-OWN-UID"); (4) the resolver holds the episode lock across the human prompt (intended; blocks other terminal operations on that episode meanwhile).

**Gate 3b:** PR #659 head `6e354749` attempt 1 **3/3 success** (no row 267 hit). `origin/dev` unmoved; PR `MERGEABLE CLEAN`; merged `--match-head-commit 6e354749` → `5d3b0581`; merge CI 3/3 success. **LANDED.**

**Routing evidence:** base=f8d8e3ad27b8442e1cd62674052b0392f03089a6@2026-11-25T03:18:33Z (Gate 1; sprint worktree created from the same SHA, snap = gate1). Gate 4 base=5d3b0581d32c4a086594c67f9afda0f95cf52f9b@2026-11-25T05:25:31Z (drift = this merge only).
- Controller `claude:claude-opus-5-5` (tok: not reported) · designer not spawned (plan-only fire on a reviewed design; rotation pointer unchanged) · planner `codex:gpt-6.2-sol` (479,066 tok) · executor `codex:gpt-6.2-sol` (59,843 tok) · evaluator Agent `sonnet` (114,085 tok). **Every required role spawned; none failed to spawn.**
- The run prompt asked for the Agent tool for every role "exactly as the skill's routing table specifies"; the resolver named the cross-provider codex recipe for planner and executor (the Agent tool is not a valid path for a `recipe` role), so those ran under the codex recipe and the evaluator used the Agent tool. No role ran inline; every code change was authored by the executor from the planner's compiled patches and judged by the evaluator.
- Judge independence: **cross-vendor** (OpenAI planner + executor, Anthropic judge).
- Metered: **$0.00** (no quorum, no metered lane).

**Progress:** 1.0 clauses 1–7 MET; goal unmoved. Row 243 V4: S1, S2, S2b, S3, S4, S4b, S5, S5b landed. Every terminal provisioning and mint path now refuses an operator-owned episode; the atomic mint and its attended recovery exist and are unreachable from production. S5c (worktree, files, paired-process consumer) is next.

**Ruled out:**
- **A round 2 for the seven judge survivors:** each guard is present and correct (one reproduced first-party), and none is reachable — `CommitOperatorMint` has no production caller, so the resolver can find no receipt on a real installation either. Killing tests are owed before S5c gives the transaction a caller, and row 243 carries them by name. Iteration 431 carried its three survivors the same way.
- **Treating the `verify_go` red as the sprint's:** untouched package, passes alone at head, on the load-sensitive list, red in iteration 430's sweep as well; PR and merge CI green.
- **Row 268's first draft "iterations 430 and 431":** iteration 431's ticket evidence names iteration 424 as instance 1; corrected before commit.

**Retro:**
- **The planner retro rule held a third time** (123/123 named kills; 0 blocking). Keep it.
- **The judge's first mutant batch targeted the wrong package** (the mint tests live in `host/authority`, the code in `host/store`), so half its ledger read SURVIVED until re-run. Instance 1; for S5c, name the package each named test lives in beside every production file in the judge directive.
- **Row 267 did not hit this PR** (3/3 on attempt 1); its count stays 5.

**Next:** V4 S5c planner — the S5b plan's Appendix C (worktree reservation and attended disposition, no-overwrite task/token/MCP publication, full live-input reconstruction, paired-process consumer/drain with its four named process pins, S3b), carrying judge-432's seven survivors as must-kill tests before any production caller of `CommitOperatorMint`, S5's two parity survivors, J-OWN-UID, Z-ABSENT-ANY and Z-WAIT-READBUDGET.

**Gate 5 bookkeeping:** rotation not due (#635 opened 2026-11-23T07:34Z; next 2026-11-30). Weekly external-issue sweep not due (424 ran it for this week). Idle fingerprint NOT recorded (this fire picked and landed). Harness share last 20: 0/20 (20 tagged). Last 3 landings (432, 431, 430) move no clause; no clause is UNMET → no DRIFT line. Row 268 records iteration 431's fleet ticket in the queue.

**Record judge** (Agent `sonnet`, ~42 s, 56,023 tok): **PASS**, no false facts; one wording caution (M2's 1,400 is cumulative, not per-milestone) applied before push.

## 433 — 2026-11-25 — row 243 V4 sprint 5c LANDED (PR #661, merge `3a6e98a7`): killing tests for five of the seven judge-432 survivors and both S5 parity survivors, and the live-input reconstruction of a signed headless start (`reconstructOperatorStart`, no production caller); worktree reservation, file publication and the paired-process consumer split to S5d/S5e/S5f; nothing dispatches; 401 lines; judge sonnet PASS 87 [PRODUCT]

**Pick:** row 243 V4 sprint 5c, the Next named by record 432.
- **Re-measured at pick** (`origin/dev` `606f6e2c`): no V4 S5c plan, branch, PR or worktree; product diff `5d3b0581..606f6e2c` empty; quorum reader rc 0 (r2 `…T10-24-51Z` PROCEED 3/3, `absent_reviewers` = `[]`); design unchanged since `8dc91744`; ledger valid, 0 OPEN of 107. Carried premise re-measured: `MintOperatorStart` has no production caller (`git grep` rc 1 outside its definition).
- **Clause map:** 1.0 clauses 1–7 MET, none UNMET. Row 243 serves D-WORLD-94 and moves no clause; it is the only loop-routable build row.
- Claim sent to `controlplane` before routing.

**Gate 0:** stamped first; `gh` = `sunholo-voight-kampff`; billing CLEAN; kill switch absent. 0 allowlisted directives on #635 (10 comments). Step 6a rc 0 (0 crash, 0 transcript notices on #635/#574; `-prev` watermark absent → single-watermark read). `mission-world` inbox: 0 unread.

**Gate 1:** local `dev` `af7a610c` behind `origin/dev` `606f6e2c` — read state from origin; did not touch the main checkout (Mark's untracked drafts). Base `606f6e2c…@2026-11-25T07:20:23Z`. Skill drift 0 (14 files) at both the resolved symlink target and the pin worktree. Dev check-runs 3/3 green.

**Routing:** the driver's codex probe timed out (rc 124 after 120 s) for planner and executor, so the resolver returned planner `agent-tool opus fail-closed:env-pin` (derive-planner-lane same), executor `recipe claude:claude-haiku-5-5 declared:provider-pin`, evaluator `agent-tool sonnet declared:alias-pin`. Followed verbatim; no codex re-probe.

**Planner** (Agent `opus`, foreground, ~27 min, **272,409 tok**). Plan `w-world-address-bar-verbs-v4s5c-sprint-plan.md` + JSON committed `987cfaf6`. M1 91 (tests only) / M2 401 cumulative (`git diff -M --numstat 606f6e2c HEAD -- cmd host docs app`). It took a small, complete, effect-free sprint and cut the rest by measurement: worktree reservation with its attended disposition (S5d; 59-line probe plus the 132-line S3b bank), file publication (S5e; 53), and the consumer/drain with the four process pins (S5f/S6; 72) — 316 measured probe lines, full sizes unmeasured. **40/40 named mutant kills** (four first survived and were strengthened or retargeted); `E-RES-EPISODE-CONFLICT` and `E-RES-REREAD` need an interleaving between the resolver's first read and its lock, so they are carried to the paired-process sprint. Premises F1–F5 (§1.4): `E-TX-KIND` is killable only by refusal order; a profile-less exec preset is already refused by the seam; a paired daemon+runner upgrade passes `checkDaemon` and is caught by the new `served_build` fact; nothing yet produces the signed resolution facts, so M2 defines the list the catalogue must sign.
- **Controller round trip** (`ctl_roundtrip.sh M1 M2`) on a fresh `git archive 606f6e2c`: every `diff -r` vs `st-M*` rc 0.

**Executor** (`claude-haiku-5-5` via `claude-role`/`claude-sub`, probe rc 0, 9 turns; tokens in 18 / out 3,934 / cache read 325,864 / cache create 29,529). M1, M2 each `ACCEPT_OK` (PASS 117 / 137, 0 FAIL); snapshots byte-equal to the planner's stages; `commit_ms.sh` → `5e884ac5` (91), `272d8591` (401).

**Gates** on `272d8591` (plan §6, 1,257 s): **all 31 legs rc 0** — incl. `pkg-cmd-ailang-worldd` 263 s, `race-cmd-ailang-worldd` 389 s, `verify_go` 488 s, inventories, routes, `start-input`, scope and callers. Observation of the real LaunchAgents, `~/.ailang-world`, the credentials store and the launchctl count: identical before and after.

**Evaluator r1** (Agent `sonnet`, foreground, ~35 min, 99,087 tok; own worktree `.wt-world-iter433-row243-v4s5c-eval`; the directive named the package holding each production file's killing tests, per iteration 432's retro, and posed questions without controller verdicts): **PASS 87**, 0 blocking (correctness 26/30, tests 24/30, security 19/20, standards 9/10, scope 9/10). 45-line ledger: 31 own-id mutants and 11 plan rows; must-die `N-LOOP-DISABLED2` killed, must-live `N-SWAP-EQ-OPERANDS-LIVE` survived. All seven M1 kills (five judge-432 survivors, two parity survivors) confirmed behaviourally in the right package; F1 confirmed; no existing assertion weakened; the FIFO/directory refusals are portable. Non-blocking: (1) no test pins that `workspace`/`hosts` reach the seam — `workspaceRoot: ""` at `operator_start_input.go:20` survived the whole `cmd/ailang-worldd` package; (2) five more named-only survivors (`N-STORE-RAW-CRED`, `N-WORKTREE-DERIVED`, `N-ATTEMPT-ZERO`, `N-CLOCK-CONST`, `N-RES-ABS`) judged equivalent or fixture-bound, none a signed-field comparison; (3) `ID`/`Argv`/`Options`/`Review` are derived from signed fields, not re-read (acceptable); (4) `-race` `cmd/ailang-worldd` 389 s against the 480 s bound; (5) the planner's `NV-M2` does not vet — the judge's `N-LOOP-DISABLED2` replaces it.
- **Controller whole-package replay** of the planner's 38 CLI rows: **38/38 KILLED** (`mutate.py --whole-only --jobs 2` on the planner's `st-M1`/`st-M2` copies, 08:16–09:49 UTC, ~92 min, finishing after the PR merge at 08:46 UTC, ledger `plan-scratch/whole.1w2fc5ww/ledger.jsonl`); each row's failing tests are its own named family (`TestOperatorStartChanged`/`TestOperatorStartFacts`, `TestOperatorMintRecovery`/`TestOperatorMintTokenReadback`, `TestPresetStartRead*`), no unrelated failure. `NV-M2` compiles under `go test` (it fails only `go vet`'s unreachable-code check, as the judge found) and is killed.

**Gate 3b:** PR #661 head `272d8591` attempt 1 **3/3 success**. `origin/dev` had moved to `f8fb3ac1` (attended #650, row 228; 0 shared files; the PR run tested the merge ref). `MERGEABLE CLEAN`; merged `--match-head-commit 272d8591` → `3a6e98a7`. Merge CI attempt 1 red on the race gate: `FAIL host/daemon 480.136s` and `FAIL host/verifygate 480.038s` (8-minute timeouts; neither package touched; parent `f8fb3ac1` green) — row 267 instance 6; attempt 2 **3/3 success**. **LANDED.**

**Routing evidence:** base=606f6e2c57c79e98975153f055e213df143638d0@2026-11-25T07:20:23Z (Gate 1; sprint worktree from the same SHA, snap = gate1). Gate 4 base=3a6e98a7fbeaf0382458cd83e056b3aa666f5a65@2026-11-25T08:46:57Z (drift = #650 + this merge).
- Controller `claude:claude-opus-5-5` (tok: not reported) · designer not spawned (plan-only fire on a reviewed design; rotation pointer unchanged) · planner Agent `opus` (272,409 tok) · executor `claude:claude-haiku-5-5` (3,952 in+out; 355,393 cache) · evaluator Agent `sonnet` (99,087 tok). **Every required role spawned; none failed to spawn.**
- The run prompt asked for the Agent tool for every role "exactly as the skill's routing table specifies"; the resolver named the Agent tool for planner and evaluator and the `claude:` CLI recipe for the executor, so the executor ran as a nested `claude-sub -p` and the other two as Agent sub-agents. No role ran inline; every code change was authored by the planner as compiled patches, applied by the executor and judged by the evaluator.
- Judge independence: **model-level only, FLAGGED** — planner opus, executor haiku, judge sonnet are all Anthropic, because the codex lane was down at the driver (D-WORLD-48: prefer a different provider; independence is a nice-to-have when none has quota).
- Metered: **$0.00** (no quorum, no metered lane; the haiku run is subscription OAuth with API keys stripped).

**Progress:** 1.0 clauses 1–7 MET; goal unmoved. Row 243 V4: S1, S2, S2b, S3, S4, S4b, S5, S5b, S5c landed. The mint transaction's defence-in-depth guards now have killing tests (5 of 7), and a signed headless start can be checked against its live inputs before any effect; both remain unreachable from production. S5d (worktree reservation with its attended disposition) is next.

**Ruled out:**
- **Folding worktree reservation into S5c to use the cap's headroom (999 lines):** the planner measured only probes for the reservation, publication and consumer, each needing injected failures at every step boundary and its own attended remedy; landing a new uncertain outcome with unmeasured recovery breaks the binding resolver gate.
- **A round 2 for judge-433's `workspace`/`hosts` forwarding gap:** non-blocking, inert while `reconstructOperatorStart` has no caller; carried as an S5d pre-caller must-kill.
- **Treating the merge's race red as the sprint's:** neither `host/daemon` nor `host/verifygate` touched; PR merge-ref run and parent push run green; attempt 2 green.

**Retro:**
- **The planner retro rule held a fourth time** (40/40 named kills; 0 blocking), now on an Opus planner rather than codex. Keep it.
- **Naming the test package beside each production file in the judge directive worked** (iteration 432's retro): the judge ran every survivor mutant in the package holding its test, and no ledger re-run was needed. Keep it.
- **`rotate-log --index-only` targets the registry's `workdir` (the main checkout), and Gate 4's `$PWD/missions` form does not apply to a repo with no `missions/`:** regenerated through a scratch registry copy instead; filed as fleet ticket `skill:rotate-log-index-only-targets-registry-workdir` (row 270, blocking none).
- **Timestamps: `ls -lT` prints local time.** The record judge caught a mixed local/UTC window in the first draft; read artifact times with `date -u -r <file>`.
- **Row 267 hit a merge again, this time with `host/verifygate` timing out in the same run** — instance 6, the first with two packages over 480 s at once.

**Next:** V4 S5d planner — worktree reservation with its attended worktree disposition in the same patch, and the S3b bank; carrying judge-433's `workspace`/`hosts` forwarding test and named-only survivors, `E-RES-EPISODE-CONFLICT`/`E-RES-REREAD` (paired-process sprint, before any production caller of `MintOperatorStart`), J-OWN-UID, Z-ABSENT-ANY and Z-WAIT-READBUDGET.

**Gate 5 bookkeeping:** rotation not due (#635 opened 2026-11-23T07:34Z; next 2026-11-30). Weekly external-issue sweep not due (424 ran it for this week). Idle fingerprint NOT recorded (this fire picked and landed). Harness share last 20: 0/20 (20 tagged). Last 3 landings (433, 432, 431) move no clause; no clause is UNMET → no DRIFT line.

**Record judge** (Agent `sonnet`, ~2 min, 89,065 tok): **FAIL on three minor facts**, all corrected before push — the executor's 40 s is wall time (19 s API), the mutant replay window is 08:16–09:49 UTC (it finished after the merge), and the judge's finding 2 lists five equivalent survivors, not eight; plus row 267's "five instances" sentence relabelled "instances 1–5". Every SHA, size, count, score and CI conclusion checked matched.

## 434 — 2026-11-25 — row 243 V4 sprint 5d LANDED (PR #666, merge `bb214516`): the S3b step bank, the start-input forwarding pins, and a signed start's worktree reservation with an independent Git readback and the attended `service runner reconcile-worktree` in the same patch (`reserveOperatorStartWorktree`, no production caller); file publication and the paired-process consumer stay in S5e/S5f; nothing dispatches; 1,397 lines; judge sonnet PASS 84 [PRODUCT]

**Pick:** row 243 V4 sprint 5d, the Next named by record 433.
- **Re-measured at pick** (`origin/dev` `68c36735`): no S5d plan, branch, PR or worktree; quorum reader rc 0 (r2 `…T10-24-51Z` PROCEED 3/3, `absent_reviewers` = `[]`); design unchanged since `8dc91744`; ledger valid, 0 OPEN of 107. Carried premises re-measured: no production caller of `MintOperatorStart` or `reconstructOperatorStart` (positive control `authority.RecheckVerb` present); `host/store/operator_steps.go` absent. Product delta since S5c = attended row-269 Phase 1 (#664, #665, 11 files): it shares one file with S5d, `cmd/ailang-worldd/main.go` (row 269's hunk and S5d's one-line help edit are separate, per the planner's preflight); the rest is outside S5d's files.
- **Clause map:** 1.0 clauses 1–7 MET, none UNMET. Row 243 serves D-WORLD-94 and moves no clause; it is the only loop-routable build row.
- Claim sent to `controlplane` before routing.

**Gate 0:** stamped first; `gh` = `sunholo-voight-kampff`; billing CLEAN; kill switch absent. 0 allowlisted directives on #635 (12 comments). Step 6a rc 0 (0 crash, 0 transcript notices on #635/#574; `-prev` watermark absent → single-watermark read). `mission-world` inbox: 0 unread.

**Gate 1:** local `dev` = `origin/dev` = `68c36735`. Base `68c36735…@2026-11-25T11:18:33Z`. Skill drift 0 (14 files) at both the resolved symlink target and the pin worktree. **Dev check-runs: race gate red** at `68c36735` and at `d6bc362f` (both attended row-269 commits), each `FAIL host/daemon` at 480 s on a slow runner (`host/verifygate` 425–442 s vs ~300–345 s on green neighbours) — row 267 instances 7 and 8; parent `2fbf6966` green; `gh run rerun 45600281525 --failed` → 3/3 green. Kept the pick (the red is row 267's known class, not this mission's code).

**Routing:** resolver planner `recipe codex:gpt-6.2-sol declared:planner-lane-default-pin`, executor `recipe codex:gpt-6.2-sol declared:provider-pin`, evaluator `agent-tool sonnet declared:alias-pin`; codex probe rc 0. Followed verbatim. Designer not spawned (reviewed design, no new doc).

**Planner** (`codex exec`, sandboxed, ~48 min, **435,121 tok**). Plan `w-world-address-bar-verbs-v4s5d-sprint-plan.md` + JSON committed `5d170ecd`. M1 280 / M2 1,397 cumulative (`git diff -M --numstat 68c36735 HEAD -- cmd host docs app`; the cap is 1,400). It kept reservation, Git effect, readback and the attended disposition in one milestone (the binding resolver gate) and put the S3b bank plus the judge-433 forwarding pins in M1. 222-row ledger: 220 named kills (32 `host/store` rows also whole-package; 188 CLI whole-package rows sandbox-uninformative), 2 inherited equivalents.
- **Controller round trip** (`ctl_roundtrip.sh M1 M2`) on a fresh `git archive 68c36735`: every `diff -r` vs `st-M*` rc 0.
- **Controller out-of-sandbox replay** of the 188 CLI named rows (`mutate.py --named-only --jobs 3`, ~4 min): **188/188 KILLED**. `mutate.py` calls a row KILLED only on a `--- FAIL: <named test>` line with no sandbox bind-denial text, so a whole-package re-run of a named kill is redundant; the 188-row whole-package replay the plan's §7 lists (`whole-mutants.sh`; controller estimate ~4.5 h at `--jobs 3`, from the 257 s `pkg-cmd-ailang-worldd` leg) was not run. The 2 equivalents' whole-package runs in the planner's git-less archive copies were uninformative (the baseline red below); the judge decided them.

**Executor** (`codex exec`, sandboxed, ~3 min, 72,886 tok). M1, M2 each `ACCEPT_OK` (PASS 41 / 121, 0 FAIL); snapshots 5/0/5 and 14/0/14, byte-equal to the planner's stages; `commit_ms.sh` → `2ee5468e` (280), `875fcce2` (1,397).

**Gates** on `875fcce2` (plan §6, 1,256 s): every S5d leg rc 0 — `pkg-host-store`, `race-host-store`, `s5d-accept`, `added-cli-race` (35 s), `start-input`, `mint-quickstart`, inventories, routes, scope, callers. Observation of the real LaunchAgents, `~/.ailang-world`, the credentials store and the launchctl count: identical before and after. **Red, pre-existing:** `pkg-cmd-ailang-worldd`, `race-cmd-ailang-worldd` and `verify_go` on 14 real-binary e2e tests that still expect 10 se-tools after attended #650 (row 228) made it 12 — reproduced at base `68c36735` in a pristine worktree; iteration 433's `verify_go` (pre-#650) was rc 0; CI skips these tests without both real binaries. Filed as **row 271**.

**Evaluator r1** (Agent `sonnet`, foreground, ~69 min, 120,246 tok; own worktree `.wt-world-iter434-row243-v4s5d-eval`; sequenced on the gate sweep's ALL_DONE marker): **PASS 84**, 0 blocking (correctness 26/30, tests 22/30, security 18/20, standards 9/10, scope 9/10). 86-line ledger; 12 plan rows re-measured (all killed behaviourally), 58 own ids; must-die `CTL-DIE-ORDER` killed, must-live `CTL-LIVE-MSG` survived. Confirmed: the lock fd is held across every effect (`TestOperatorStartOnce`'s `LOCK_NB` probe kills the early-release and reservation-after-add mutants); the reservation commits and its writer closes before Git; readback is independent; every `uncertain` outcome exits through `reconcile-worktree`; the remedy text runs verbatim; nothing removes Git evidence; the two inherited survivors are equivalent on the whole package. Non-blocking: (1) store reservation guards (`host/store/operator_worktree.go:30-41`: kind, headless, store path, approver, challenge, claimed state) survive their named tests — kind on the whole package; (2) readback common-dir, `.git` link shape and registration-count checks have no killing test; (3) reconcile tamper checks survive named runs; (4) no reconcile test on an already-created step; (5) receipt-size tested only at max+1; (6) race cost ~38 s against the plan's 28.9 s.
- **Controller reproduction** of two judge survivors (edit asserted by `git diff --stat` + `go vet`, restored clean): `J-ST-KIND` survives `go test ./host/store/` and the CLI worktree set; `J-W-STORECHK` (`operator_worktree.go:37`) survives the CLI worktree set. Both inert; carried to S5e M1.

**Gate 3b:** PR #666 head `875fcce2` attempt 1 **3/3 success** (run 45625696797). `origin/dev` unchanged; `MERGEABLE CLEAN`; merged `--match-head-commit 875fcce2` → `bb214516`. Merge CI run 45628783011 attempt 1 red on one test only: `FAIL host/broker` `TestSeedRestartIgnoresASurvivingOldCopy/survivor` (a known load-sensitive test; `host/broker` unchanged since base; the merge tree equals the green PR head's tree); `gh run rerun --failed` → attempt 2 **3/3 success**. **LANDED.**

**Routing evidence:** base=68c3673579b53326a2c5d4d7bc62b9d3bb7a1105@2026-11-25T11:18:33Z (Gate 1; sprint worktree from the same SHA, snap = gate1). Gate 4 base=bb214516b20865a8f59df46113011a0df1bd0c33 (drift = this merge only).
- Controller `claude:claude-opus-5-5` (tok: not reported) · designer not spawned (plan-only fire on a reviewed design; rotation pointer unchanged) · planner `codex:gpt-6.2-sol` (435,121 tok) · executor `codex:gpt-6.2-sol` (72,886 tok) · evaluator Agent `sonnet` (120,246 tok) · record judge Agent `sonnet` (101,021 tok). **Every required role spawned; none failed to spawn.**
- The run prompt asked for the Agent tool for every role "exactly as the skill's routing table specifies"; the resolver named the codex recipe for planner and executor and the Agent tool for the evaluator, so planner and executor ran as `codex exec` and the judge as an Agent sub-agent. No role ran inline; every code change was authored by the planner as compiled patches, applied by the executor and judged by the evaluator.
- Judge independence: **cross-vendor** (OpenAI planner/executor, Anthropic judge).
- Metered: **$0.00** (no quorum, no metered lane; codex runs on the OpenAI subscription plan).

**Progress:** 1.0 clauses 1–7 MET; goal unmoved. Row 243 V4: S1, S2, S2b, S3, S4, S4b, S5, S5b, S5c, S5d landed. A signed start can now reserve and create its worktree under the episode lock with an attended way out of every uncertain result; it remains unreachable from production. S5e (judge-434 guard tests, then file publication with its attended disposition) is next.

**Ruled out:**
- **Running the plan's 188-row whole-package CLI replay (controller estimate ~4.5 h at `--jobs 3`):** each row's named kill was a real `--- FAIL` of the named test, which the whole package contains; the controller re-ran all 188 named rows outside the sandbox instead (188/188).
- **A round 2 for the judge's guard survivors:** the sprint sits at 1,397 of 1,400 and the survivors are on a path with no production caller; S5b → S5c set the precedent of a survivor-pin M1 in the next sprint.
- **Treating the 14 e2e reds as this sprint's:** identical at base in a pristine worktree, `host/daemon` byte-identical to base, and iteration 433's pre-#650 `verify_go` green.
- **Treating dev's race reds as a pick-changer:** row 267's signature (8-minute `host/daemon` timeout on a slow runner), parent green, attempt 2 green.

**Retro:**
- **The planner retro rule held a fifth time** (220/220 named kills; 0 blocking), back on codex.
- **A whole-package mutant run in a `git archive` copy is uninformative for `cmd/ailang-worldd` now** — the 14 real-binary e2e tests red there anyway (row 271), so every mutant reads "REVIEW FAILURE". Until row 271 lands, decide CLI whole-package survivors in a git worktree, against a recorded baseline fail set (the judge did exactly that).
- **`mission-worktree.sh add` returns before the checkout exists.** I wrote into the record path straight after `add` printed rc 0; my `mkdir -p` made `git worktree add` fail rc 128 and the helper cleaned up, deleting the copied file. Re-ran `add`, then `wait` until `ready`, then wrote. Standing rule 7's launched-not-done clause, met first-hand; no new rule needed — call `wait` every time.
- **A stale count in a CI-skipped test is invisible to CI by construction.** #650 passed CI with 14 real-binary tests broken on the rig; only a local `verify_go` sees them. Worth knowing for attended merges that touch `packages/se-tools`.
- **Row 267 hit dev twice on attended commits** (instances 7 and 8, ~20 h, 8 instances).

**Next:** V4 S5e planner — M1 killing tests for judge-434's reservation-path survivors (store guards, readback shape/common-dir, reconcile tamper checks, receipt-size boundary) before any production caller of `reserveOperatorStartWorktree`; then no-overwrite task/token/MCP publication with its attended file disposition; carry `E-RES-EPISODE-CONFLICT`/`E-RES-REREAD`, J-OWN-UID, Z-ABSENT-ANY, Z-WAIT-READBUDGET and the four process pins to S5f/S6. Row 271 (a test-expectation fix plus one design read) is small and keeps the loop's local `verify_go` usable.

**Gate 5 bookkeeping:** rotation not due (#635 opened 2026-11-23T07:34Z; next 2026-11-30). Weekly external-issue sweep not due (424 ran it for this week). Idle fingerprint NOT recorded (this fire picked and landed). Harness share last 20: 0/20 (20 tagged). Last 3 landings (434, 433, 432) move no clause; no clause is UNMET → no DRIFT line.

**Record judge** (Agent `sonnet`, ~3 min, 101,021 tok): **FAIL on two facts**, both corrected before push — the judge report was missing from `design_docs/verification/world-iter434/` (the first copy was lost when the record worktree's first `add` failed), and "product delta … outside S5d's files" was false at file granularity (both touch `main.go`). Two imprecisions also corrected: which commit added which `host/daemon` test lines (row 267), and the 188-row replay's duration is a controller estimate. One finding refuted by measurement: `OperatorStep` exists as a method (`host/store/operator_steps.go:72`). Every SHA, size, token count, mutant count, score, timing and CI conclusion it checked matched.

## 435 — 2026-11-25 — row 243 V4 sprint 5e LANDED (PR #669, merge `2df9dc80`): reservation-path pre-caller pins — killing tests for the judge-434 survivors (store reservation and disposition guards, receipt at the limit, Git readback shape and common dir, recovery tamper and observed-step refusals, a locked re-read seam) before anything calls `reserveOperatorStartWorktree`; file publication measured at 1,462 cumulative lines and kept whole for the next sprint; nothing dispatches; 495 lines; judge sonnet PASS 89 [PRODUCT]

**Pick:** row 243 V4 sprint 5e, the Next named by record 434.
- **Re-measured at pick** (`origin/dev` `ec4565ec`): no S5e plan, branch, PR or worktree; quorum reader rc 0 (r2 `…T10-24-51Z` PROCEED 3/3, `absent_reviewers` = `[]`); design unchanged since `8dc91744`; ledger 107 rows, 0 OPEN. No product change since S5d (`git diff --name-only bb214516 ec4565ec -- cmd host docs app` = 0). Production callers of `reserveOperatorStartWorktree`/`reconstructOperatorStart`/`MintOperatorStart`: definitions plus the one internal call (positive control `authority.RecheckVerb` present).
- **Clause map:** 1.0 clauses 1–7 MET, none UNMET. Row 243 serves D-WORLD-94 and moves no clause; it is the loop-routable build row.
- Claim sent to `controlplane` before routing.

**Gate 0:** stamped first; `gh` = `sunholo-voight-kampff`; billing CLEAN; kill switch absent. 0 allowlisted directives on #635 (13 comments). Step 6a rc 0 (0 crash, 0 transcript notices on #635/#574). `mission-world` inbox: 0 unread.

**Gate 1:** local `dev` 6 behind `origin/dev` (iteration 434's sprint and record), fast-forwarded on a clean tracked tree; `dev` = `origin/dev` = `ec4565ec`. Base `ec4565ec…@2026-11-25T15:19:03Z`. **Skill drift:** 13 of 14 files identical to `origin/dev`; the running `resources/gate-2-pick.md` carries an uncommitted 2-line reviewer-pool edit in the ailang main checkout (astra/carbon/fable seats, dated 2026-11-25) — followed as the running copy; no quorum ran, so it did not bind. Dev CI 3/3 green.

**Routing:** resolver planner `recipe codex:gpt-6.2-sol declared:planner-lane-default-pin`, executor `recipe codex:gpt-6.2-sol declared:provider-pin`, evaluator `agent-tool sonnet declared:alias-pin`; codex probe rc 0. Followed verbatim. Designer not spawned (reviewed design, no new doc).

**Planner** (`codex exec`, sandboxed, ~41 min, **363,768 tok**). Plan `w-world-address-bar-verbs-v4s5e-sprint-plan.md` + JSON committed `0d3e1020`. It built and measured both halves, then shipped **M1 alone (495 lines)**: the file-publication milestone (task/token/MCP, private no-replace install, independent readback, attended `reconcile-files`) prototyped green on its focused tests at **967 incremental / 1,462 cumulative**, 62 over the cap, and stays whole in Appendix C rather than split from its disposition. All 30 judge-434 survivor rows reproduced at base; at M1 31/31 named kills (10/10 store rows whole-package too), 4 equivalents with stated invariants.
- **Controller round trip** on a fresh path and a fresh `git archive ec4565ec`: `diff -r` vs `st-M1` rc 0 (cmd, host, docs).

**Executor** (`codex exec`, sandboxed, ~2–3 min, 50,092 tok): patch clean, `ACCEPT_OK` (87 PASS / 0 FAIL), snapshot 5/0/5, every file byte-equal to the planner's stage; `commit_ms.sh` → `fd17279f` (495).

**Controller baseline + gates.** Pristine detached base checkout at `ec4565ec`: plain and race `cmd/ailang-worldd` fail sets = row 271's 11 CLI tests exactly. **The plan's `verify_go` leg was an instrument failure** — it runs `GOFLAGS=-json ./scripts/verify_go.sh`, which turns `go env` into JSON and kills `verify_go.sh` at its toolchain check before any test; the JSON parser then read an empty fail set from a 1.9 KB log. Ran a controller copy of `gates.sh` (2 lines: plain `verify_go`, a text fail-set parser with a firing control). Gates on `fd17279f` (1,497 s): 32 legs rc 0, CLI plain/race fail sets equal to base; `verify_go` = row 271's 14 plus three charter-listed load-sensitive tests (`TestSeToolsRunBenchmarkTasksEndToEnd`, `TestSeToolsRunReplayIsByteEqual`, `TestMissionConfigDefaultFlowIntegration`), each rc 0 alone at head in untouched packages; real-path observation unchanged.

**Evaluator r1** (Agent `sonnet`, foreground, ~55 min, 111,756 tok; own worktree; sequenced on `ALL_DONE`): **PASS 89**, 0 blocking (correctness 27/30, tests 25/30, security 19/20, standards 9/10, scope 9/10). 58-line ledger, 25 own ids; must-die controls killed, must-live survived; 14 plan rows re-measured, all killed behaviourally; all four equivalents accepted with the masking invariant checked in code; per-hunk non-vacuity confirmed for all three production hunks; the start-kind move is behaviour-preserving (decode already enforced `a.Kind == row.Kind`); the reader seam is reachable only with the real reader; no existing test edited; numstat 495; new CLI tests +6.85 s under `-race`; agreed the gate copy's `verify_go` verdict is sound. Non-blocking: (1) `row.Owner != at.Runner.Nonce` (`host/store/operator_worktree.go:31`) unpinned — survives whole `host/store` and whole `cmd/ailang-worldd` (fail set = base); (2) dropping `addErr != nil` at `cmd/ailang-worldd/operator_worktree.go:101` survives whole cmd (a successful add with an absent readback would read "stopped: absent" instead of uncertain); (3) store empty identity, expiry `>=`, read-state filter, `count = 1`, dispose SQL `kind = 'start'`, and `readRunnerWorktree` signed-kind/repo-abs/common-abs/re-read-error survive named runs (store ones whole `host/store` too); (4) `exact != 1` → `< 1` and `inspect`'s `IsAbs(raw)` probably equivalent.
- **Controller reproduction** of finding 1 (judge worktree, edit asserted, restored clean): `go test -count=1 ./host/store/` rc 0 with the owner guard disabled. Inert today; carried.

**Gate 3b:** PR #669 head `fd17279f` 3/3 success. `origin/dev` had moved to `e6606ffe` (attended #667, episode GC; no file overlap); `MERGEABLE CLEAN`; merged `--match-head-commit fd17279f` → `2df9dc80`. Merge CI run 45665379388 attempt 1: race gate red, `FAIL host/daemon 480.080s` running `TestWorkbenchTimelinePaging`, `host/verifygate` 414.9 s — row 267's signature (instance 9); parent `e6606ffe` race green (`host/daemon` 346.5 s); S5e touches nothing in `host/daemon` (the PR-head run's `host/daemon` already took 445.7 s); `gh run rerun --failed` → attempt 2 **3/3 success** (run 45665379388 attempt 2). **LANDED.**.

**#667 × row 243 (found this iteration, code reading):** attended #667's `RunEpisodeGC` (inside `serve` 2 min after boot, then every 6 h) removes any directory under the workspace root (default `~/.ailang/world-ws`, the episode worktree root) that is ≥72 h old, has no live session and is clean with no unmerged commits. It reads no operator action or step row. A reserved or `uncertain` operator worktree — the evidence `reconcile-worktree` reads back — would qualify (inferred, not run: a fresh operator worktree's HEAD is the base it was cut from, so `checkWorktreeUnsafe`'s merged-into-`dev`/`main`/`@{u}` test passes). Latent while `reserveOperatorStartWorktree` has no production caller; carried as a row-243 pre-caller gate.

**Routing evidence:** base=ec4565ec5eedfed693f9ea5f2c78b89b63bf432e@2026-11-25T15:19:03Z (Gate 1; sprint worktree from the same SHA). Gate 4 base=2df9dc802ecd1e69dfa45f38cb404a92cb027c57@2026-11-25T17:30:07Z (drift = attended #667 + this merge).
- Controller `claude:claude-opus-5-5` (tok: not reported) · designer not spawned (plan-only fire on a reviewed design; rotation pointer unchanged) · planner `codex:gpt-6.2-sol` (363,768 tok) · executor `codex:gpt-6.2-sol` (50,092 tok) · evaluator Agent `sonnet` (111,756 tok) · record judge Agent `sonnet` (84,764 tok; PASS, 0 errors, 6 imprecisions corrected or noted). **Every required role spawned; none failed to spawn.**
- The run prompt asked for the Agent tool for every role "exactly as the skill's routing table specifies"; the resolver named the codex recipe for planner and executor and the Agent tool for the evaluator, so planner and executor ran as `codex exec` and the judge as an Agent sub-agent. No role ran inline.
- Judge independence: **cross-vendor** (OpenAI planner/executor, Anthropic judge).
- Metered: **$0.00** (no quorum, no metered lane).

**Progress:** 1.0 clauses 1–7 MET; goal unmoved. Row 243 V4: S1–S5d and S5e (M1) landed. The worktree reservation path now has killing tests for its judge-434 survivors, still with no production caller. File publication (sprint S5e2, measured 967 lines on top of M1) is next.

**Ruled out:**
- **The plan's 25-row whole-package CLI replay (~25 × 4–6 min):** every named kill was a real `--- FAIL` of the named test outside a bind denial, and the named tests need no socket; the judge re-ran 14 rows and the equivalents itself. Whole-package runs stayed for survivors only.
- **A round 2 for judge-435's unpinned guards:** all sit on a path with no production caller; the first caller is S5f's consumer, and the next sprint's M1 takes them (the S5b→S5c, S5d→S5e precedent).
- **Splitting file publication to fit the cap** (e.g. task + token now, MCP later): the planner's measured full shape needs the admission matrix and linked-directory refusal, and the reviewed design gives no approval to split Claude's three-file bundle.
- **Treating the extra `verify_go` reds or the merge race red as this sprint's:** all in packages S5e does not touch, each passing alone or on the parent.

**Retro:**
- **A planner's gate script can be wrong in a way that reads green.** `GOFLAGS=-json` broke `verify_go.sh` before any test ran, and an empty fail set looked like a clean baseline. Caught only because the leg's rc (1) disagreed with its empty fail set. Next directive: name `verify_go.sh`'s text output as the thing to parse, and require each fail-set parser to show a firing control.
- **The planner reused the controller's round-trip path for its rehearsal**, so the controller's run refused ("export dir already exists"). Fresh path per invocation (standing rule 7's stale-artifact clause); the directive should give the planner its own path.
- **A split decided by measurement held up:** the planner prototyped the whole publication slice, found it 62 lines over, and kept the disposition with its effect rather than squeezing it. The judge called that defensible.
- **Attended merges can change a loop row's safety premises without touching its files:** #667's GC shares the episode worktree root with operator reservations. Worth a pre-caller gate, not an alarm.
- **Row 267 hit a merge again** (instance 9; the PR-head run's `host/daemon` was 445.7 s of 480).

**Next:** V4 S5e2 planner — M1 killing tests for judge-435's survivors (owner guard `host/store/operator_worktree.go:31`, `addErr` absent disposition `cmd/ailang-worldd/operator_worktree.go:101`, store identity/expiry/read-state/count/dispose-kind, `readRunnerWorktree` kind/abs/re-read-error); then file publication with its attended `reconcile-files` disposition in one milestone (planner-measured 967 lines on M1; prototype in `~/.ailang/state/world-iter435/plan-scratch/st-M2`). Before S5f's consumer (the first production caller): episode GC must not remove a worktree with a non-terminal operator step. Row 271 stays small and keeps local `verify_go` readable.

**Gate 5 bookkeeping:** rotation not due (#635 opened 2026-11-23T07:34Z; next 2026-11-30 07:00 local). Weekly external-issue sweep not due. Idle fingerprint NOT recorded (this fire picked and landed). Harness share last 20: 0/20 (20 tagged). Last 3 landings (435, 434, 433) move no clause; no clause is UNMET → no DRIFT line.

## 436 — 2026-11-25 — row 243 V4 sprint 5e2 LANDED (PR #672, merge `d9c24a3b`): killing tests for the judge-435 reservation-path survivors (owner, empty identity and exact expiry at reservation; uncertain read state and duplicate-row count; the dispose kind CAS; success-then-absent stays uncertain with its executed remedy; signed kind, absolute repo/common dir, locked re-read error and relative registration in the attended resolver); tests only, no production change; file publication re-measured at 1,259 incremental / 1,464 with M1 and kept whole a second time; #667's GC measured to qualify reserved and published operator worktrees; 205 lines; judge sonnet PASS 84 [PRODUCT]

**Pick:** row 243 V4 sprint 5e2, the Next named by record 435.
- **Re-measured at pick** (`origin/dev` `5cac9316`): no S5e2 plan, branch, PR or worktree; quorum reader rc 0 (r2 `…T10-24-51Z` PROCEED, `absent_reviewers` = `[]`); design unchanged since `8dc91744`; ledger 107 rows, 0 OPEN. No product change since S5e (`git diff --name-only 2df9dc80 5cac9316 -- cmd host docs app` = 0). Attended #667 overlaps the S5e publication prototype in `cmd/ailang-worldd/episode.go` (directive told the planner to re-base, not apply blind).
- **Clause map:** 1.0 clauses 1–7 MET, none UNMET. Row 243 serves D-WORLD-94 and moves no clause; it is the loop-routable build row.
- Claim sent to `controlplane` before routing.

**Gate 0:** stamped first; `gh` = `sunholo-voight-kampff`; billing CLEAN; kill switch absent. 0 allowlisted directives on #635 (14 comments). Step 6a rc 0 (0 crash, 0 transcript notices on #635/#574). `mission-world` inbox: 0 unread.

**Gate 1:** `dev` = `origin/dev` = `5cac9316`. Base `5cac9316…@2026-11-25T19:17:43Z`. **Skill drift:** pin copy identical to `origin/dev` on all 14 files; the running copy differs only in `resources/gate-2-pick.md` — the same uncommitted 2-line reviewer-pool edit iteration 435 recorded; followed as the running copy; no quorum ran, so it did not bind. Dev CI 3/3 green.

**Routing:** resolver planner `recipe codex:gpt-6.2-sol declared:planner-lane-default-pin`, executor `recipe codex:gpt-6.2-sol declared:provider-pin`, evaluator `agent-tool sonnet declared:alias-pin`; codex probe rc 0. Followed verbatim. Designer not spawned (reviewed design, no new doc).

**Planner** (`codex exec`, sandboxed, ~30 min, **322,946 tok**). Plan `w-world-address-bar-verbs-v4s5e2-sprint-plan.md` + JSON committed `81ebfd7e`. **M1 alone (205 lines, tests only, three existing test files).** M2 — file publication with its attended `reconcile-files` — was re-based onto `5cac9316` (pure rebase 967 / 1,172 cumulative) and then hardened to close S5e Appendix C's listed gaps: **1,259 incremental / 1,464 cumulative**, 64 over the cap, kept whole in Appendix C rather than shipped without its disposition. All 13 judge-435 survivors reproduced at base; at M1 11 KILLED by named behavioural tests (all six store rows also whole `host/store`, per the planner's ledger), 2 measured equivalents with stated masking invariants. A scratch probe (`TestDeferredFilePublicationGCProbe`) found #667's GC treats both a reserved and a fully published operator worktree as a removal candidate at an injected 73 h age.
- **Controller round trip** on a fresh path and a fresh `git archive 5cac9316`: `diff -r` vs `st-M1` rc 0 (cmd, host, docs).

**Controller baseline** (plan §7 `baseline-controller.sh`, on the sprint worktree at the plan-only commit; load ~15–16 during the run, from sibling missions): plain and race `cmd/ailang-worldd` fail sets = row 271's 11 CLI tests exactly. `verify_go` (plain, text-parsed, with a firing parser control — iteration 435's `GOFLAGS=-json` lesson applied by the planner) = row 271's 14 plus `TestWorldContextAndStatusline` and `TestMissionConfigDefaultFlowIntegration`, each rc 0 alone at base. The script's `set -e` stopped at that diff before writing the base files, so the controller copied the three measured sets by hand.

**Executor** (`codex exec`, sandboxed, ~2 min, 57,392 tok): patch clean, `ACCEPT_OK` (101 PASS / 0 FAIL), snapshot 3/0/3, every file byte-equal to the planner's stage; `commit_ms.sh` → `3795994c` (205).

**Gates on `3795994c`** (plan `gates.sh`, 1,500 s): 33 of 34 legs rc 0 — CLI plain/race fail sets equal to base, `linux-vet`, inventories, routes, scope, callers; real-path observation before = after. `verify_go` = base set plus `host/daemon TestSeToolsRunBenchmarkTasksEndToEnd`. **Attributed by negative control, not by its charter "load-sensitive" label:** `host/daemon` is byte-identical base..head; run alone and interleaved, it failed 3/4 in the sprint worktree, 1/3 in the judge's worktree (same head) and 1/4 in a pristine base worktree, always `the mock saw ["GET  " "POST value123 application/json"]` — an extra empty GET reaching the test's mock server. A rig flake that exists at base.

**Evaluator r1** (Agent `sonnet`, foreground, ~57.5 min, 113,874 tok; own worktree; sequenced on the gate sweep's rc file): **PASS 84**, 0 blocking (correctness 26/30, tests 25/30, security 19/20, standards 10/10, scope 4/10). 52-line ledger, 39 ids not in the plan. All 13 judge-435 survivors re-run: 11 KILLED by named runtime tests against the real store, real Git and the real CLI; `Q-RB-EXACT` and `Q-R-SIGNEDKIND` SURVIVED the whole package with fail set = base, and both masking invariants checked in code (`exact` increments only inside the delimiter-appending branch; a signed finish is never headless). Must-die controls killed, must-live survived. Zero production `.go` changed; no existing assertion weakened; race cost ≈ +0.08 s (the planner's +3.5 s did not reproduce). Non-blocking (judge's numbering): (1) **publication deferred a second time although M2 alone (1,259) fits the cap** — "a third deferral should be refused"; (2) five pre-existing store guards survive whole `host/store`: heartbeat lower bound (`operator_worktree.go:59`), `episode_starts`/`session_credentials` NOT EXISTS (`:56`), read `phase = 'reserved'` (`:89`), read other-kind NOT EXISTS (`:90`), dispose other-kind NOT EXISTS (`:138`); (3) `signed_kind` is a non-discriminating regression pin, since its mutant is equivalent; (4) the race-delta mismatch; (5) two of its own equivalents lack a whole-package run. On #667: the GC claim holds (caveats: the probe stubs `IsActive` false; the published arm needs the deferred M2 code).
- **Controller reproduction** of finding 2 (judge worktree, edit asserted, restored clean): `ReadOperatorWorktree` without `AND phase = 'reserved'` → `go test -count=1 ./host/store/` rc 0. Inert today; carried.

**Gate 3b:** PR #672 head `3795994c` 3/3 success. `origin/dev` had moved to `90e97b1f` (attended #671, row 224 standing session allowances; no file overlap with S5e2; its own dev push CI run 45697364842 was red on row 267's signature — race job `FAIL host/daemon 480.082s`); `MERGEABLE CLEAN`; merged `--match-head-commit 3795994c` → `d9c24a3b`. Merge CI run 45705517790 attempt 1: **3/3 success**. **LANDED.**

**Routing evidence:** base=5cac9316b755bf4b5b33263448366ec877926c8b@2026-11-25T19:17:43Z (Gate 1; sprint worktree from the same SHA). Gate 4 base=d9c24a3bf3b94674f57e305096fb5c04665f3525@2026-11-25T21:34:19Z (drift = attended #671 + this merge).
- Controller `claude:claude-opus-5-5` (tok: not reported) · designer not spawned (plan-only fire on a reviewed design; rotation pointer unchanged) · planner `codex:gpt-6.2-sol` (322,946 tok) · executor `codex:gpt-6.2-sol` (57,392 tok) · evaluator Agent `sonnet` (113,874 tok) · record judge Agent `sonnet` (88,394 tok; PASS, 0 errors, 6 imprecisions corrected: load figure, the intermediate dev red, a run count, the judge's finding numbering, one non-killing test name in row 243, and an unattributed planner claim). **Every required role spawned; none failed to spawn.**
- The run prompt asked for the Agent tool for every role "exactly as the skill's routing table specifies"; the resolver named the codex recipe for planner and executor and the Agent tool for the evaluator, so planner and executor ran as `codex exec` and the judge as an Agent sub-agent. No role ran inline.
- Judge independence: **cross-vendor** (OpenAI planner/executor, Anthropic judge).
- Metered: **$0.00** (no quorum, no metered lane).

**Progress:** 1.0 clauses 1–7 MET; goal unmoved. Row 243 V4: S1–S5e and S5e2 (M1) landed. The worktree reservation path now has killing tests for every judge-435 survivor (two measured equivalents), still with no production caller. No publication code has landed.

**Ruled out:**
- **Attributing `TestSeToolsRunBenchmarkTasksEndToEnd` to this sprint, or dismissing it on its charter label alone:** the label was not evidence; the interleaved base runs were (it fails at base too, with the same stray-GET signature).
- **A round 2 for judge-436's store survivors:** all five sit on a path with no production caller and pre-date S5e2; the next sprint takes them before the first caller (the S5b→S5c … S5e→S5e2 precedent).
- **Re-planning this fire to ship M2 alone:** the directive set M1-first; the planner measured both and the M1 pins are what M2 builds on. The judge's point stands for the NEXT sprint, not as a reason to discard a measured, green M1.
- **Trusting the baseline script's absence of base files as "no base":** its `set -e` stopped at a known load-sensitive diff; the measured sets were copied by hand and named as such.

**Retro:**
- **Two consecutive sprints deferred publication by measuring it together with a pin milestone.** Each sprint's M1 (495, then 205) was judged correct, and each time M1 + M2 crossed the cap while M2 alone would have fit. Next directive: **M2 is the sprint — publication + `reconcile-files` alone (planner-measured 1,259)**, with any new pins only if they fit beside it, never ahead of it.
- **A charter "load-sensitive" label is a claim.** This time the listed test failed ALONE, unloaded, 3 of 4 runs at head; only a base control showed it was not ours. The signature (an extra empty GET on a loopback mock) suggests something on the rig dialling ephemeral ports — worth a row-267-style tracking note rather than a blanket label.
- **A planner's baseline script with `set -e` loses its own output on the first expected diff.** Ask for "record, then compare" ordering in the next directive.
- **I backgrounded a CI poller with `&` inside a foreground call** (rule 7's double-background); caught at once, killed, relaunched under a fresh verdict path. No harm; noting it so the habit stays visible.

**Next:** V4 S5e3 planner — **file publication + attended `reconcile-files` as the sprint's whole scope** (start from `~/.ailang/state/world-iter436/plan-scratch/st-M2` and `M2-deferred.patch`, 1,259 incremental on `5cac9316`; re-base onto the current head — attended #671 merged since), closing its Appendix C matrix (flag bits, timing, refusal branches, Linux behaviour). Judge-436's five store survivors and #667's GC skip for non-terminal operator steps stay pre-caller gates before S5f's first production caller. Row 271 stays small and keeps local `verify_go` readable.

**Gate 5 bookkeeping:** rotation not due (#635 opened 2026-11-23T07:34Z; next 2026-11-30 07:00 local). Weekly external-issue sweep not due. Idle fingerprint NOT recorded (this fire picked and landed). Last 3 landings (436, 435, 434) move no clause; no clause is UNMET → no DRIFT line.

## 437 — 2026-11-26 — row 243 V4 sprint 5e3 LANDED (PR #674, merge `22bb7804`): signed start file publication with independent readback and the attended `service runner reconcile-files` as one milestone — the task, in-memory mint token and Claude MCP bundle installed privately with exclusive no-replace under the episode lock after one S3b reservation, read back before the step is observed; #671's allowance origin refused; no production caller; 1,650 lines under the directive's safe-unit exception; judge sonnet PASS 79 [PRODUCT]

**Pick:** row 243 V4 sprint 5e3, the Next that record 436 named.
- **Re-measured at pick** (`origin/dev` `eaa202b1`):
  - No S5e3 plan, branch, PR or worktree existed.
  - Quorum reader rc 0: r2 `…T10-24-51Z` PROCEED, `absent_reviewers` = `[]`. Design unchanged since `8dc91744`.
  - Ledger 107 rows, 0 OPEN.
  - Attended #671 (`90e97b1f`) shares 4 files with the deferred publication patch (`comm -12`, positive control firing). The directive told the planner to re-base rather than apply blind, and to decide the allowance credential origin with a test.
- **Clause map:** 1.0 clauses 1–7 MET, none UNMET. Row 243 serves D-WORLD-94 and moves no clause; it is the loop-routable build row.
- Claim sent to `controlplane` before routing.

**Gate 0:**
- Stamped first. `gh` = `sunholo-voight-kampff`. Billing CLEAN. Kill switch absent.
- 0 allowlisted directives on #635 (16 comments).
- Step 6a rc 0: 0 crash and 0 transcript notices on #635/#574.
- `mission-world` inbox: 0 unread.

**Gate 1:**
- Local `dev` was 5 behind (S5e2 plan, M1, #672 merge, record 436 and the #673 merge). Fast-forwarded; `dev` = `origin/dev` = `eaa202b1`. Base `eaa202b1…@2026-11-25T23:18:25Z`.
- **Skill drift:** the running copy and the pin copy are each identical to `origin/dev` on all 14 files.
- Dev CI 3/3 green.

**Routing:**
- Resolver: planner `recipe codex:gpt-6.2-sol declared:planner-lane-default-pin`, executor `recipe codex:gpt-6.2-sol declared:provider-pin`, evaluator `agent-tool sonnet declared:alias-pin`. Followed verbatim.
- Codex probe: rc 0 from the worktree. A first probe from `/tmp` returned rc 1 on codex's trusted-directory check — the wrong CWD, an instrument error, not a lane failure.
- Designer not spawned (reviewed design, no new doc).

**Directive (controller's call, stated in it):** publication and its disposition are the WHOLE sprint and may not be deferred a third time — the S5e2 judge wrote "a third deferral should be refused". The minimum safe unit may exceed the 1,400 cap up to 1,650 rather than split publication from `reconcile-files`. Optional store pins go AFTER it, only if they fit.

**Planner** (`codex exec`, sandboxed, ~47 min, **430,131 tok**). Plan `w-world-address-bar-verbs-v4s5e3-sprint-plan.md` + JSON committed `d34c358b`.
- **M1 alone: 1,650 lines** (1,625+/25−, 16 files).
  - The pure rebase over #671 reproduced the inherited 1,259; hardening the S5e2 Appendix C gates added 391.
  - The #671 origin is refused before reservation through the real `GrantAllowance`/`MintFromAllowance` producer (`TestOperatorFileAllowanceOrigin`).
- Mutants: 95 of 98 named KILLED, 3 MASKED (`F-CREATE-NOFOLLOW`, `R-POST-PARENT`, `S-WORKTREE-PHASE`). The planner's own diff census lists 122 AUDIT lines with no independent mutation measurement.
- The optional store pins were moved to Appendix C (zero capacity left).
- **Controller round trip** on a fresh `git archive eaa202b1`: `diff -r` against `st-M1` rc 0 (cmd, host, docs).

**Executor** (`codex exec`, sandboxed, ~2.5 min, 77,895 tok):
- Patch clean, `ACCEPT_OK` (320 PASS / 0 FAIL). Snapshot 16/0/16, every file byte-equal to the planner's stage.
- `commit_ms.sh` reran acceptance outside the sandbox and committed `de608652` (1,650).

**Controller baseline** (`baseline-controller.sh`, in a pristine detached base worktree per the iter-295 rule):
- CLI plain and race fail sets = row 271's 11 exactly.
- **`verify_go` came back with 60 failures — an INSTRUMENT FAILURE, not a base.** The plan's `gate-bin/git` PATH wrapper ran `exec "$WORLD_GATE_GIT"`, and tests that scrub their child env drop that variable. Every `git` they spawned died `exec: : not found` (43 extra `host/broker` failures; the CLI failures were identical in both runs).
- Fixed in scratch with a real-git fallback: the original gives rc 127 under `env -i`, the fixed copy rc 0. The contaminated files are kept.
- Re-recorded the base `verify_go` alone: 18 = row 271's 14 plus four tests on the known load-sensitive list (two of them also charter-listed).

**Gates on `de608652`** (plan `gates.sh`, 1,622 s, fixed wrapper, rig observation on):
- Every leg rc 0 except `compare-verify_go` (below). CLI plain/race fail sets = base. Linux vet ok. Inventories, callers and file race (48 s) ok. Real-path observation before = after.
- `compare-verify_go` rc 1 only because head had one FEWER failure (`host/broker TestExecProtectedCheckTimeout` passed). No new failure.

**Evaluator r1** (Agent `sonnet`, background, own worktree, sequenced on `gates-final.rc`, ~119 min, 133,078 tok): **PASS 79**, 0 blocking.
- Scores: correctness 26/30, tests 20/30, security 18/20, standards 8/10, scope 7/10.
- What holds:
  - No write without a reservation, no overwrite or symlink follow, no observe without readback.
  - The only reachable production path is the attended verb, and its remedy text is executed verbatim.
  - The allowance refusal uses the real producer.
  - The 3 masks are genuine on the whole package (fail set = base).
  - 11 of 11 sampled planner KILLs confirmed. The planner's `S-LATE-APPROVAL` fails `go vet`; a vet-clean variant is killed.
  - The 1,650 exception was earned.
- What does not: of 32 valid mutants of its own, 9 were killed and 23 survive their named tests (about two thirds). Non-blocking, ranked:
  1. The N04/N12 directory-swap windows. Both fail safe to `uncertain`.
  2. N33/N34 mask each other. No test shows `reconcile-mint` refusing once a files reservation exists.
  3. The store admission predicates in `host/store/operator_files.go` (N22–N31) have no `host/store` test and survive whole `host/store`.
  4. N15.
  5. Whole-package CLI race at 95% of the 480 s bound locally.
  6. The `gate-bin/git` wrapper should not exist as a gate input.
- **Controller reproduction** of finding 3 (judge worktree, edit asserted, restored clean): dropping `now >= row.ExpiresAt ||` at `host/store/operator_files.go:191` gives `go test ./host/store/` rc 0 and the CLI file tests rc 0. Inert today (no caller); carried.

**Gate 3b:**
- PR #674 head `de608652`, run 45731365295. Attempt 1 race red: `FAIL host/daemon 480.093s` (row 267) with `cmd/ailang-worldd` at **477.565 s**. Re-run failed → attempt 2 3/3 (`cmd/ailang-worldd` 262.975 s).
- `origin/dev` unchanged, `MERGEABLE CLEAN`. Merged `--match-head-commit de608652` → `22bb7804`, tree-identical to the head.
- Merge run 45743770171, attempt 1, red on two jobs:
  - `host/broker TestSeedRestartIgnoresASurvivingOldCopy` (charter-listed load-sensitive);
  - `FAIL host/daemon 480.125s` (row 267), with `cmd/ailang-worldd` at **468.977 s**.
  - Both packages are byte-identical to the base.
- Re-run failed → **attempt 2 3/3 success** (`cmd/ailang-worldd` 259.477 s). **LANDED.**

**Routing evidence:** base=eaa202b1d23f601ce597f063034a061b5bf5dea3@2026-11-25T23:18:25Z (Gate 1; sprint worktree from the same SHA). Gate 4 base=22bb78042d52dde832e2a8a2c4f2a652364db702@2026-11-26T03:11:45Z (drift = this merge).
- Roles and tokens: controller `claude:claude-opus-5-5` (tok: not reported) · designer not spawned (plan-only fire on a reviewed design; rotation pointer unchanged) · planner `codex:gpt-6.2-sol` (430,131 tok) · executor `codex:gpt-6.2-sol` (77,895 tok) · evaluator Agent `sonnet` (133,078 tok) · record judge Agent `sonnet` (102,710 tok; FAIL on imprecisions only — no wrong SHA, run id, token or line count; corrected before commit: the contaminated-set arithmetic, "every leg rc 0", which load-sensitive tests are charter-listed, the judge's survivor count (23 of 32, not ~22), the Gate-1 behind count, a line count and a two-run "typical" range). **Every required role spawned; none failed to spawn.**
- The run prompt asked for the Agent tool for every role "exactly as the skill's routing table specifies". The resolver named the codex recipe for planner and executor and the Agent tool for the evaluator, so planner and executor ran as `codex exec` and the judge as an Agent sub-agent. No role ran inline.
- Judge independence: **cross-vendor** (OpenAI planner/executor, Anthropic judge).
- Metered: **$0.00** (no quorum, no metered lane).

**Progress:** 1.0 clauses 1–7 MET; goal unmoved. Row 243 V4: S1–S5e3 landed.
- A signed start's files can now be published, read back and, when uncertain, resolved by attended review — all still dormant behind the S5f/S6 and #667 GC gates.
- Publication, deferred twice, has landed.

**Ruled out:**
- **Banking the first base `verify_go` set (60).** The 44 `host/broker` `exec: : not found` lines were the gate's own wrapper, not the repo. A comparison against it would have hidden any real `git`-spawning regression at head.
- **Attributing either CI red to S5e3.** Neither `host/daemon` nor `host/broker` changed. The merge tree equals the green PR head, and both re-runs went 3/3.
- **Calling `cmd/ailang-worldd`'s 477.6/469.0 s pure noise.** The same head ran 263/259 s on normal runners, so the slow runner explains most of it. But this sprint adds ~51 s locally, and the package is now within 11 s of the bound on a slow runner. Recorded on row 267 as a new package at risk.
- **A round 2 for judge-437's survivors.** All sit on paths with no production caller, and the S5b→S5e2 precedent takes such pins in the next sprint, before the first caller.

**Retro:**
- **A gate script that puts a wrapper on `PATH` is an instrument that the code under test can defeat.** Tests that scrub their child env silently dropped the wrapper's variable. The 60-line base set was internally consistent and would have made a head comparison look "equal". It was caught only because 60 against a 16–18 history did not look like a baseline. Next directive: no PATH shims in checkout-mode gates — the archive bridge only where an archive is the input.
- **The planner's 95/98 overstated the coverage; the judge's own sample came back about 2/3 surviving.** A planner ledger enumerates the mutants it thought of. The 122 AUDIT lines it disclosed were where the judge found the survivors. Next directive: every AUDIT line in a diff census needs a measured disposition or goes in Appendix C as a named pin.
- **CI race headroom for `cmd/ailang-worldd` is now a constraint on row 243 itself.** New pins go in `host/store`, where the judge also says the store predicates belong.

**Next:** V4 S5e4 planner — pre-caller pins before S5f's first production caller:
- judge-437's store admission predicates as `host/store` tests (N22–N31; N31 reproduced);
- the N33/N34 `reconcile-mint`-after-files-reservation refusal;
- judge-436's five `operator_worktree.go` store survivors;
- N04/N12/N15 if a seam is cheap;
- the #667 GC operator-step-aware skip, measured.

Keep `cmd/ailang-worldd` race time flat. Row 271 stays small.

**Gate 5 bookkeeping:**
- Rotation not due (#635 opened 2026-11-23T07:34Z; next 2026-11-30 07:00 local). Weekly external-issue sweep not due.
- Idle fingerprint NOT recorded (this fire picked and landed).
- Harness share 0/20. The last 3 landings (437, 436, 435) move no clause; no clause is UNMET → no DRIFT line.

## 438 — 2026-11-26 — row 243 V4 sprint 5e4 LANDED (PR #676, merge `e554db6b`): the pre-caller pins — `host/store` behavioural tests kill all sixteen judge-437/436 store survivor edits (eight file-admission predicates N22–N31, N33/N34 — each now with a sole killer — and judge-436's six worktree-store edits), and #667's episode GC skips any worktree whose episode has an unresolved operator step on the daemon pass and the `episode gc` CLI, failing safe; no production caller; 929 lines; judge sonnet PASS 88 [PRODUCT]

**Pick:** row 243 V4 sprint 5e4, the Next that record 437 named.
- **Re-measured at pick** (`origin/dev` `684069f2`): no S5e4 plan, branch, PR or worktree existed; design r1 reviewed (r2 PROCEED 3/3); ledger 107 rows, 0 OPEN.
- **Clause map:** 1.0 clauses 1–7 MET, none UNMET. Row 243 serves D-WORLD-94 and moves no clause; it is the loop-routable build row.
- Claim sent to `controlplane` before routing (`inbox_1795665451823_88ce66ef`).

**Gate 0:** stamped first; `gh` = `sunholo-voight-kampff`; billing CLEAN; kill switch absent; 0 allowlisted directives on #635 (17 comments); step 6a rc 0 (0 crash, 0 transcript notices; control 107:4).
- **Inbox:** two fleet `harness-resolved` replies. Row 268 (`skill:mission-base-relative-path-records-driver-repo`, fleet #2466 `fcc3d3481`) and row 270 (`skill:rotate-log-index-only-targets-registry-workdir`, fleet #2471 `92bc214ba`) are marked RESOLVED upstream. Row 270's flag is not yet reachable: the rig's `ailang` v0.62.0-dirty still rejects `--workdir`, so this record used the scratch-registry recipe again.

**Gate 1:** local `dev` 5 behind (record #675), only untracked files in the main checkout, fast-forwarded; base `684069f2…@2026-11-26T03:55:33Z`. Skill drift: resolved symlink and pin copy each byte-identical to `origin/dev` on all 14 files. Dev CI 3/3 green.

**Routing:**
- Resolver (followed verbatim): planner `recipe codex:gpt-6.2-sol declared:planner-lane-default-pin`, executor `recipe codex:gpt-6.2-sol declared:provider-pin`, evaluator `agent-tool sonnet declared:alias-pin`. Codex probe rc 0 from the sprint worktree. Designer not spawned (reviewed design, no new doc).
- Controller `claude-opus-5-5` (session) · planner `codex:gpt-6.2-sol` (465,015 tok) · executor `codex:gpt-6.2-sol` (59,934 tok) · evaluator Agent `sonnet` (112,169 tok). Planner and executor are cross-provider recipes, not Agent-tool spawns, because the resolver returned `recipe` for both; the evaluator is the Agent tool. judge-independence: cross-vendor.

**Directive (controller's):** M1 store pins, M2 the files fence, M3 the GC skip, M4 optional; ≤ 1,400 lines; pins in `host/store` because CI race for `cmd/ailang-worldd` (477.6 s on a slow runner) and `host/daemon` (480.0x s) sit at the 480 s bound (figures carried from row 267 and record 437, not measured this fire); at most one pure in-process GC test in `host/daemon`; no PATH shims (iteration 437's `gate-bin/git` lesson); mutation counts on both the diff-anchored set and every named judge survivor.

**Planner** (~68 min). Plan + JSON `7f7ff262`.
- M1 320 (`host/store/operator_s5e4_test.go`), M2 52 (`operator_mint_s5e4_test.go`), M3 557 (new `host/store/operator_episode_gc.go` + test, `host/daemon/auto_gc.go` + `auto_gc_s5e4_test.go`, `cmd/ailang-worldd/episode_gc.go` + test). **929** changed lines (`-- cmd host docs app`).
- M4 (N04/N12/N15) prototyped and killed in scratch, 82 lines and +4.4 s focused CLI race: parked in Appendix C.
- Ledgers: 37 diff-anchored mutants killed in named tests; all 15 required judge edits killed in named tests and whole `host/store`. Survivors: N01 (mask — no observed mint receipt) and the parked M4 ids.
- Premise corrected: N33/N34 pass `fileStage=false` at base; each mutant flips its own call to `true`.

**Controller round trip:** `ctl_roundtrip.sh M1 M2 M3` → `diff -r` rc 0 for cmd, host and docs at every milestone.

**Executor** (~5 min): three clean patches; `ACCEPT_OK` M1 20/0, M2 24/0, M3 55/0; snapshots 1/0/1, 2/0/2, 9/0/9, all `cmp`-equal to the stages. Controller commits M1 `4f8ff2be`, M2 `54cb1b96`, M3 `1b78d779` (NUMSTAT 320 → 372 → 929).

**Controller baseline** (pristine `684069f2` worktree with the planner running; load 4–6, spiking to 17 during the CLI race leg and 9 during `verify_go`): `host/store` race 0 fails; `host/daemon` race the 3 row-271 tests; `cmd/ailang-worldd` plain the 11 row-271 tests. **Instrument failure, discarded:** base `cmd/ailang-worldd` race and plain `verify_go.sh` hit `panic: test timed out after 10m0s` (verify_go in three packages), so those fail sets were truncated. CI arbitrated whole-package race and `verify_go`.

**Gates** on `1b78d779` (no PATH shims, `-timeout 20m`): build, vet, Linux vet, test compile, nine inventories, `verify_ail.sh`, GC probe (reserved and published worktrees now skipped as `non-terminal operator step`), caller census = base — all rc 0. `host/store` plain/race 0 fails; CLI plain fail set = base; daemon race = base + `TestSeToolsRunBenchmarkTasksEndToEnd` (charter-listed load-sensitive; passes alone at head, 8.3 s). `gofmt` flags `host/daemon/daemon.go`, untouched and flagged at base too. Whole-package mutant replay (28 ids) not run: every id has a named killer, the survivors are N01 and the parked M4 ids, and the replay would take about four hours (UNMEASURED).

**Evaluator** (Agent `sonnet`, foreground, ~34 min, own worktree at `1b78d779`, clean at the end): **PASS 88**, 0 blocking (28/25/19/9/7).
- All 16 judge-437/436 survivor edits (N22–N25, N27–N29, N31, N33, N34, and judge-436's six) killed by named `host/store` tests and whole `host/store`. N01 confirmed as a mask.
- GC policy judged correct: a reserved step, or any non-`operation` step under an action with no `terminal_at`, protects; `operation` rows are alias identity rows (schema CHECK: observed, no owner, no receipt). The CLI cannot be steered to an unrelated store; the daemon uses `d.credentialsPath`; no new context root.
- 31 own M3 mutants: all daemon and CLI arms killed; store S6 a must-live equivalent; **S9 (drop `files`) and S10 (drop `mint`) survive whole `host/store`** — the lookup's test table holds only `worktree`-kind rows.
- Finding 2: an absent, pre-v10, locked or unresolvable credentials store makes GC skip every worktree (fail-safe, reason reported) — a doc line would help.
- **Controller reproduction of S9:** `s.kind NOT IN ('operation','files')` at `host/store/operator_episode_gc.go:23` → `go test -count=1 ./host/store/` ok 12.171 s. Survives; carried to S5f.

**Gate 3b:** PR #676 CI 3/3 on attempt 1 (race job: `cmd/ailang-worldd` 262.9 s, `host/daemon` 355.8 s, `host/broker` 280.8 s, `host/store` 62.8 s). `origin/dev` unchanged; `MERGEABLE CLEAN`; merged `--match-head-commit` → `e554db6b` (tree-identical). Merge CI on `e554db6b` (run 45765797616): 3/3 success on attempt 1. **LANDED.**

**Progress:** 1.0: clauses 1–7 MET; goal unmoved. Row 243 V4: S1–S5e4 landed; every pre-caller store survivor edit the S5e2/S5e3 judges carried (16) is now killed, and GC no longer removes a worktree an unresolved operator effect needs. Still dormant: no production caller.

**Ruled out:** a third executor round to add the S9/S10 lookup rows before landing — the production SQL is correct today, the gap is a test-table one, and it is carried as a named pin before S5f's first caller rather than buying a re-judge; running the 28-id whole-package replay (about four hours, every id already has a named killer).

**Next:** V4 S5f — the paired-process consumer and drain with the first production caller, carrying the S9/S10 lookup-table pins, N04/N12/N15, `E-RES-EPISODE-CONFLICT`, `E-RES-REREAD`, J-OWN-UID, Z-ABSENT-ANY, Z-WAIT-READBUDGET and the four process pins.

**Gate 5 bookkeeping:**
- Rotation not due (#635 opened 2026-11-23T07:34Z; next 2026-11-30 07:00 local). Weekly external-issue sweep not due.
- Idle fingerprint NOT recorded (this fire picked and landed).
- Harness share 0/20. The last 3 landings (438, 437, 436) move no clause; no clause is UNMET → no DRIFT line.

## 439 — 2026-11-26 — row 243 V4 sprint 5f LANDED (PR #681, merge `886648ab`): the staged seams and pins the first production caller needs — judge-438's S9/S10 GC kind rows, the attended mint resolver's injected lock (E-RES-EPISODE-CONFLICT and E-RES-REREAD killed), a dormant start selector with a lock-scoped candidate visit, a signed `launch_binary_sha256` start fact, and a heartbeat that no longer advances `catalogue_generation` on an identical payload; the caller moves to S5f2; 578 lines; judge sonnet PASS 86 [PRODUCT]

**Pick:** row 243 V4 sprint 5f, the Next that record 438 named.
- **Re-measured at pick** (`origin/dev` `9610f5de`): no S5f plan, branch, PR or worktree existed; design r1 reviewed (r2 PROCEED 3/3); ledger 107 rows, 0 OPEN.
- **Clause map:** 1.0 clauses 1–7 MET, none UNMET. Row 243 serves D-WORLD-94 and moves no clause; it is the loop-routable build row.
- Claim sent to `controlplane` before routing (`inbox_1795679938770_f891b6f0`).

**Gate 0:** stamped first; `gh` = `sunholo-voight-kampff`; billing CLEAN; kill switch absent; 0 allowlisted directives on #635 (18 comments); step 6a: 0 crash, 0 transcript notices (control 107:4); `mission-world` inbox 0 unread.

**Gate 1:** local `dev` = `origin/dev` = `9610f5de`; only untracked files in the main checkout; base `9610f5de…@2026-11-26T07:55:39Z`. Skill drift: resolved symlink and pin copy each byte-identical to ailang `origin/dev` on every file. Dev CI 3/3 green.

**Routing evidence:**
- Resolver (followed verbatim): planner `recipe codex:gpt-6.2-sol declared:planner-lane-default-pin`, executor `recipe codex:gpt-6.2-sol declared:provider-pin`, evaluator `agent-tool sonnet declared:alias-pin`. Codex probe rc 0 from the sprint worktree. Designer not spawned (reviewed design, no new doc).
- Controller `claude-opus-5-5` (session; tok: not reported) · planner `codex:gpt-6.2-sol` (382,611 tok; **rc=1**, see below) · executor `codex:gpt-6.2-sol` (65,180 tok) · evaluator Agent `sonnet` (108,428 tok). Planner and executor are cross-provider recipes because the resolver returned `recipe` for both; the evaluator is the Agent tool. judge-independence: cross-vendor. metered=$0.00.

**Directive (controller's):** S5f = the paired-process consumer and drain with the first production caller; M1 S9/S10, M2 E-RES-*, M3 caller + drain + four process pins and the carried J-OWN-UID/Z-*/`session revoke`/GC exclusion, M4 optional; ≤ 1,400 lines with a safe-unit exception to 1,650; "a production caller without its four process pins is not a safe unit"; if it could not fit, a pinned dormant seam with the caller in Appendix C at a MEASURED size; total added `cmd/ailang-worldd` race ≤ 12 s (row 267).

**Planner** (~36 min). **Lane outcome rc=1:** codex printed `WROTE …sprint-plan.md and sprint_….json`, then `ERROR: This content was flagged for possible cybersecurity risk` (twice) and ended with no final message.
- **Disposition:** a measured partial, not a fallback. The JSON parses, no placeholder text, the planner's own stage checks are rc 0 for all three stages, and the controller round trip reproduces every stage. A kimi fallback would have restarted planning from nothing. Plan + JSON committed `fab82ec3`.
- Scope: M1 20, M2 98, M3 **578** cumulative (`-- cmd host docs app`). The caller is S5f2 in Appendix C. §1.4 measured three prerequisites absent at HEAD: the upgrade path refuses while a runner is installed and cannot drain a pair; the dormant worktree/files helpers open their own lock descriptor; and every heartbeat advanced `catalogue_generation`, which invalidated a signed mint/worktree CAS on the next tick.
- Left "pending" in the plan (§1.5 race receipts, §4 mutation table, §6, §7) — measured by the controller and the judge, below. Appendix C names no S5f2 size; the controller measured the prototype `proto-S5f2-current` at **1,082** cumulative lines without the four process pins.

**Controller round trip:** `ctl_roundtrip.sh M1 M2 M3` → `diff -r` rc 0 for M1 and M2; the M3 patch is byte-identical to `st-M3-refined`, not the earlier `st-M3` draft, which the plan's §4.1 still names.

**Executor** (~3 min): three clean patches; `ACCEPT_OK` M1 22/0, M2 46/0, M3 83/0; snapshots 1/0/1, 4/0/4, 17/0/17, all `cmp`-equal to the stages. Controller commits M1 `c5ab7d18`, M2 `cabe66c7`, M3 `0ed2ab53` (NUMSTAT 20 → 98 → 578).

**Controller mutants:** all 45 `mutants.json` specs run outside the sandbox against their named tests → 44 KILLED; D-SELECT-LIMIT (`LIMIT 1` dropped) SURVIVED, pre-declared EQUIVALENT. No whole-package replay.

**Gates** on `0ed2ab53` (no PATH shims, `-timeout 20m`; base fail sets reused from iteration 438's `1b78d779` head, whose code tree equals `9610f5de`): build, vet, Linux vet, test compile, inventories, `verify_ail.sh`, GC probe — rc 0. Caller census = base once the planner script got its missing 4th argument; negative control fires. `host/store` plain/race 0 fails; CLI plain = base 11; daemon race = base minus the load-sensitive `TestSeToolsRunBenchmarkTasksEndToEnd`. `gofmt` flags untouched `host/daemon/daemon.go`, as at base.

**Evaluator** (Agent `sonnet`, foreground, ~39 min, own worktree at `0ed2ab53`, clean at the end): **PASS 86**, 0 blocking (27/26/18/8/7).
- S9, S10, E-RES-EPISODE-CONFLICT and E-RES-REREAD each killed by a named behavioural test; 11 controller-KILLED ids re-run, all killed.
- The catalogue CASE terms read the pre-update row; `catalogue_generation` is read by the mint and worktree CAS and `ReadOperatorHeartbeat`, and nothing relied on it advancing every tick.
- 21 own mutants: J02 must-die killed; J07/J19/J21 equivalents; **J01 (`catalogue IS ?` → `= ?`) survives whole `host/store`**; J05 (drop the visit's `defer db.Close()`) survives the named selectors.
- Non-blocking: Appendix C has no measured size and the plan still reads "pending"; S5f2 needs a compatibility rule for starts signed before `launch_binary_sha256`; the dormant selector is a justified staged seam.
- **Controller reproduction of J01:** `catalogue IS ?` → `catalogue = ?` in `host/store/operator_runner.go` → `go test -count=1 ./host/store/` ok 12.439 s. Survives; production SQL is correct; carried to S5f2.

**Gate 3b:** PR #681 attempt 1 RED on `go host build + test` only: `host/broker TestSeedRestartIgnoresASurvivingOldCopy/survivor` ("the old staging dir survived the restart"). It is charter-listed load-sensitive; `host/broker` is untouched by this sprint (`git diff --stat 9610f5de 0ed2ab53 -- host/broker` empty); iteration 437 hit the same test on merge CI. Attempt 2 3/3 (race job: `cmd/ailang-worldd` 269.4 s, `host/daemon` 357.1 s, `host/verifygate` 328.9 s, `host/broker` 275.0 s, `host/store` 63.4 s). `origin/dev` had moved to `419cca9d` (attended #678, #679; no file overlap); `MERGEABLE CLEAN`; merged `--match-head-commit` → `886648ab`. Merge CI on `886648ab` (run 45796678393): attempt 1 RED on the race job only — `panic: test timed out after 8m0s` in `host/daemon` (480.140 s; untouched by this sprint), with `cmd/ailang-worldd` at 464.950 s on that slow runner; attempt 2 **3/3 success** (`cmd/ailang-worldd` 265.2 s, `host/daemon` 355.8 s, `host/verifygate` 324.9 s, `host/broker` 272.2 s). **LANDED.**

**Found:** dev went red mid-iteration on the two attended merges — `6ac6d966` (#678: race `FAIL host/daemon 480.077s`) and `419cca9d` (#679: race `FAIL cmd/ailang-worldd 480.232s`, `host/daemon 480.103s`, `host/verifygate 480.049s`, `panic: test timed out after 8m0s`). `419cca9d`'s race job attempt 2 was green (`6ac6d966` was not re-run), so it is row 267's slow-runner class, not a regression; recorded on row 267. With this merge's attempt 1, and attended #680's `0033ca11` after it (`FAIL host/daemon 480.110s`), it is four consecutive `dev` merges red at attempt 1 on the race bound; the broker flake (row 259) also hit the docs-only record PR #682; filed **D-WORLD-121** (take row 267 before S5f2? A recommended; default B = the attended order).

**Progress:** 1.0: clauses 1–7 MET; goal unmoved. Row 243 V4: S1–S5f landed. Every pin record 438 owed before the caller except the process pins is now killed. The catalogue-generation churn that would have refused every signed mint is fixed. Still dormant: no production caller.

**Ruled out:** falling back to the kimi planner on the codex rc=1 — the plan existed, compiled and round-tripped, so a restart would only have re-bought it on a weaker lane; a third executor round for J01/J05 — the production SQL is correct, the gaps are a test row and a handle close, carried as named S5f2 pins rather than buying a re-judge; re-recording base fail sets — iteration 438's head is code-identical to `9610f5de`.

**Next:** V4 S5f2 — the first production caller with its four process pins, the episode fd borrowed across the worktree/files helpers, the attended paired drain, J01's NULL-catalogue row, J05, a compatibility rule for pre-digest signed starts, J-OWN-UID, terminal `session revoke` exclusion, the activation-vs-GC exclusion and N04/N12/N15. Row 267 goes first only if Mark answers D-WORLD-121 = A; either way S5f2's directive must report its added `cmd/ailang-worldd` race seconds against the 480 s bound.

**Gate 5 bookkeeping:**
- Rotation not due (#635 opened 2026-11-23T07:34Z; next 2026-11-30 07:00 local). Weekly external-issue sweep not due.
- Idle fingerprint NOT recorded (this fire picked and landed).
- Harness share 0/20 (20 rows, 20 tagged). The last 3 landings (439, 438, 437) move no clause; no clause is UNMET → no DRIFT line.

## 440 — 2026-11-26 — row 243 V4 sprint 5f2 LANDED (PR #685, merge `150d0d96`): the activation prerequisites — judge-439 J01/J05 pins, the runner's borrowed per-episode lock description across the worktree and files helpers, activation-vs-GC exclusion in either lock order (a live GC change, fail-safe), N04/N12/N15 publication-window pins, and a round-2 fix for three existing refusal arms the borrow contract had broken; activation moves to S5f3; 1,476 lines; judge sonnet round 1 FAIL 67, round 2 PASS 77 [PRODUCT]

**Pick:** row 243 V4 sprint 5f2, the Next that record 439 named.
- **D-WORLD-121** (row 267 before S5f2?) is still OPEN and unanswered. Its default is "B, from the next fire" (this fire), so the groomed order stood.
- **Re-measured at pick** (`origin/dev` `b5f41458`): no S5f2 plan, branch, PR or worktree existed. The design r1 was reviewed (quorum r2 PROCEED 3/3). Ledger 108 rows, 1 OPEN.
- **Clause map:** 1.0 clauses 1–7 MET, none UNMET. Row 243 serves D-WORLD-94.
- Claim `inbox_1795694365193_283d84da`.

**Gate 0:** stamped first; `gh` = `sunholo-voight-kampff`; billing CLEAN; kill switch absent; 0 allowlisted directives on #635 (19 comments); step 6a: 0 crash, 0 transcript notices (control 107:4); `mission-world` inbox 0 unread.

**Gate 1:** local `dev` was 3 behind `origin/dev`, so the record was read from origin. Base `b5f41458…@2026-11-26T11:56:44Z`. Skill drift: resolved symlink and pin copy byte-identical to ailang `origin/dev` on all 14 files. Dev CI 3/3 green at `b5f41458`.

**Routing evidence:**
- Resolver (followed verbatim): planner `recipe codex:gpt-6.2-sol declared:planner-lane-default-pin`, executor `recipe codex:gpt-6.2-sol declared:provider-pin`, evaluator `agent-tool sonnet declared:alias-pin`. Designer not spawned (reviewed design). Codex probe rc 0.
- The operator's standing request was to spawn all roles with the Agent tool "exactly as the routing table specifies". The table returned `recipe` for the planner and executor, which makes the Agent tool not a valid path for them, so they ran as `codex exec` recipes. The REQUIRED judge ran on the Agent tool both rounds. No role failed to spawn.
- Controller `claude-opus-5-5` (session; tok: not reported) · planner `codex:gpt-6.2-sol` (914,431 tok) · executor `codex:gpt-6.2-sol` (run 1 ~75k tok, stopped on the controller's own directive defect; run 2 75,210 tok; round 2 104,364 tok) · evaluator Agent `sonnet` (r1 128,239 tok; r2 64,235 tok, fresh instance). judge-independence: cross-vendor. metered=$0.00.

**Directive (controller's):** S5f2 = activation, i.e. the first production caller in `operator-runner` mode, with its prerequisites, in priority order:
- M1: J01/J05.
- M2: the borrowed episode fd.
- M3: the activation-vs-GC exclusion.
- M4: the caller, its drain, the four process pins, J-OWN-UID, terminal `session revoke`, and N04/N12/N15 if they fit.
- Cap 1,400, with the 1,650 exception only for the complete activation unit. "A production caller without its four process pins and its drain is NOT a safe unit."
- Race budget: `cmd/ailang-worldd` ≤ 12 s, `host/daemon` ≤ 5 s.

**Planner** (~68 min, rc 0). Plan + JSON `4dd16433`.
- M1 73, M2 647, M3 1,292, and M4p (N04/N12/N15) at **1,400** cumulative.
- **Activation moves to S5f3** (Appendix C). The compiled partial prototype is 1,984 cumulative, and its 410-line core plus M1–M3 alone is 1,702, over 1,650.
- §1.4 premises found false: the upgrade path refuses a loaded companion; two prototype durable states were refused by all three attended remedies; the CLI GC override bypassed workspace binding; terminal revoke lacked episode exclusion.
- Focused race deltas: +3.28 s CLI, +1.12 s daemon.
- Controller round trip `diff -r` rc 0 for all four milestones.

**Executor:** run 1 applied M1 and stopped correctly. My directive demanded an `ACCEPT_OK` marker carried over from iteration 439's `accept.sh`; this plan's `accept.sh` signals by exit code. That is false-green rule (4): the gate list I typed was never baselined. I reset the worktree, corrected the directive, and re-ran. Run 2 delivered four clean milestones, with `accept.sh` exit 0 each time and snapshots `cmp`-equal to the stages. Controller commits: M1 `9094d9db`, M2 `1bd64507`, M3 `651077c2`, M4p `c4fb23f6`.

**Baseline instrument note:** the first base recording ran co-loaded with the planner, and `daemon-race`, `cmd-race` and `verify_go` hit Go's default 10-minute timeout. I discarded them and re-recorded `daemon-race` plus `host/worktree` and `host/store` plain with `-timeout 20m` (`base2/`). `cmd-race` and `verify_go` were left to CI.

**Gates** on `c4fb23f6` (no PATH shims):
- Clean legs: build, vet, Linux vet, test compile, inventories, `verify_ail`, GC probe, caller census (four start definitions, zero production callers) and lock census.
- `host/store` and `host/worktree` plain and race: 0 fails.
- `gofmt` red only on untouched `host/daemon/daemon.go` (red at base too).
- **`cmd-plain` added `TestOperatorStartCollisions` and `TestOperatorFilePublicationAdmission`.**
- The `daemon-race` drift (five tests) passes alone at head and at base: load.
- Named-mutant replay of all 102 specs on each milestone's own stage: **101 KILLED**. `D-LOCK-NO-NONBLOCK` SURVIVED and is a declared equivalent. The 70 whole-package comparisons were not run.

**Evaluator round 1** (Agent `sonnet`, ~70 min, own worktree): **FAIL 67** (23/14/16/8/6).
- **B1:** M2 changed the start/files test helpers to borrow the lock before the function under test. That made `TestOperatorStartCollisions/held` and `/missing_lock` and `TestOperatorFilePublicationAdmission/lock_held` red or dead. Controller-reproduced: they FAIL at `c4fb23f6` and PASS at `b5f41458`.
- What held: M2's verification, M3's non-blocking acquire, both lock orders, fail-safe skips, and start-path inertness.
- Non-blocking: N1 `J-RR-FORCE` survives; N2 `episode gc --credentials-db` for a non-default workspace skips all; N3 GC leaves persistent lock sidecars; N4 tests pin error wording.

**Round 2** (codex, test files only, 80 lines) → `c7c08e5b`.
- The three arms again drive the production helpers: a non-holder description, a sidecar unlinked after borrow, and another episode's description. Each kills a named `Verify` mutation.
- `TestRemoveRegisteredRefusesLockedOrDirty` kills `--force`.
- Controller gates: `cmd-plain` = base plus two load-flaked git-fixture tests, which pass alone. The changed cmd tests pass alone under `-race` (57.3 s).
- **Evaluator round 2** (fresh Agent `sonnet`, ~29 min): **PASS 77** (23/25/16/8/5), 0 blocking. B1, N1 and N6 CLOSED; N2 and N3 accepted; its own five mutants were each killed.

**Gate 3b:**
- PR #685 CI on `c7c08e5b` was 3/3 at attempt 1. Race: `cmd/ailang-worldd` 273.2 s, `host/daemon` 360.3 s.
- `origin/dev` had moved to `91d1cabf` (attended #683/#684, no overlap). `MERGEABLE CLEAN`; merged `--match-head-commit` → `150d0d96`.
- Merge CI (run 45852708959): **3/3 at attempt 1**. Race: `cmd/ailang-worldd` 351.0 s, `host/daemon` 409.3 s. **LANDED.**

**Found:**
- The sandboxed planner could not run whole `cmd/ailang-worldd` (socket denials), and its 70 "whole comparisons" were too slow for the controller to replay. The only thing that caught a regression of existing tests was the controller's whole-package fail-set comparison. Keep it in every gate sweep.
- GC now creates a persistent per-episode lock sidecar on every installation (N3).
- This is the third consecutive sprint to defer the caller: S5f → S5f2 → S5f3.

**Progress:** 1.0: clauses 1–7 MET; goal unmoved. Row 243 V4: S1–S5f2 landed.
- Every activation prerequisite the S5f record listed except the drain, recovery matrix, revoke exclusion and J-OWN-UID is now in.
- The GC path is exclusion-safe against a future runner.
- Still dormant: no production caller.

**Ruled out:**
- A third executor round for N2/N3: both are accepted, documented behaviour, carried as S5f3 notes.
- Whole-package replay of the 70 daemon/cmd mutants: at ~5–10 min each it would take hours. The named replay plus two independent judges' own mutants covered the diff.
- Treating the five daemon-race additions as regressions: all pass alone at head and at base.
- Dropping M4p to stay under the cap for the round-2 fix: M4p was judged and green, and the 76-line overage is a test-only repair the judge accepted.

**Next:** V4 **S5f3 — activation**. The first production caller in `operator-runner` mode carries:
- the attended paired drain and rollback, and the reconcile/disposition matrix (plan Appendix C);
- signed-executable compatibility for starts signed before `launch_binary_sha256`;
- terminal `session revoke` exclusion and J-OWN-UID;
- the four process pins;
- judge N2 (CLI `--credentials-db` scope) and N3 (lock-sidecar lifecycle) as design notes.

The judge measured about 584 lines over this head, so it fits one sprint. Its directive must report added `cmd/ailang-worldd` race seconds against the 480 s bound (D-WORLD-121 still open).

**Gate 5 bookkeeping:**
- Rotation not due (#635 opened 2026-11-23; next 2026-11-30 07:00 local). Weekly external-issue sweep not due.
- Idle fingerprint NOT recorded (this fire picked and landed).
- The last 3 landings (440, 439, 438) move no clause. No clause is UNMET, so there is no DRIFT line.

## 441 — 2026-11-26 — row 243 V4 sprint 5f3 LANDED (PR #687, merge `a50b0d72`): the last dormant prerequisites for the first production caller — terminal `session revoke` exclusion with J-OWN-UID, the seven durable-window dispositions with pre-digest signed-start compatibility, and the attended paired upgrade/drain/rollback (a live change to `service upgrade`); outcome B ends the deferral chain with a measured schedule: the caller is the single remaining sprint S5f4 (917 lines); 1,527 lines; judge sonnet round 1 FAIL 71, round 2 PASS 83 [PRODUCT]

**Pick:** row 243 V4 sprint 5f3, the Next that record 440 named.
- **D-WORLD-121** (row 267 first?) is still OPEN. Its default B applied, so the groomed order stood.
- **Re-measured at pick** (`origin/dev` `2f573488`): no S5f3 plan, branch, PR or worktree existed. Four start-path declarations and zero callers. Zero process-pin test names, with positive control `operator_upgrade_test.go:13`.
- **Clause map:** 1.0 clauses 1–7 MET, none UNMET. Row 243 serves D-WORLD-94.

**Gate 0:** stamped first; `gh` = `sunholo-voight-kampff`; billing CLEAN; kill switch absent; 0 allowlisted directives on #635 (21 comments, LEGACY watermark); step 6a: 0 crash and 0 transcript notices (control 107:4); `mission-world` inbox 0 unread.

**Gate 1:** local `dev` was 9 behind `origin/dev`, so the record was read from origin. Base `2f573488…@2026-11-26T16:26:44Z`. Skill drift: resolved symlink and pin copy byte-identical to ailang `origin/dev` on all 14 files. Dev CI at `2f573488` was in progress at Gate 1; parent `150d0d96` 3/3 green.

**Routing evidence:**
- Resolver (followed verbatim): planner `recipe codex:gpt-6.2-sol declared:planner-lane-default-pin`, executor `recipe codex:gpt-6.2-sol declared:provider-pin`, evaluator `agent-tool sonnet declared:alias-pin`. Designer not spawned (reviewed design). Codex probe rc 0.
- The operator's standing request was to spawn every role with the Agent tool "exactly as the routing table specifies". The table returned `recipe` for the planner and executor, which makes the Agent tool not a valid path for them, so they ran as `codex exec` recipes. The REQUIRED judge ran on the Agent tool both rounds, and was resumed by name for a round-3 addendum. No role failed to spawn.
- Controller `claude-opus-5-5` (session; tok: not reported) · planner `codex:gpt-6.2-sol` (599,391 tok) · executor `codex:gpt-6.2-sol` (run 1 73,830 tok; round 2 157,191 tok; round 3 68,031 tok) · evaluator Agent `sonnet` (r1 149,830 tok; r2 fresh instance 63,389 tok; r3 addendum 66,062 tok). judge-independence: cross-vendor. metered=$0.00.

**Directive (controller's):** S5f3 = activation, but this time the plan had to END the deferral. Either (A) the complete safe unit lands within 1,650, or (B), only if the measured unit exceeds 1,650, the dormant prerequisites land in a fixed priority order with a measured schedule whose last sprint is the caller and fits ≤ 1,400. Priority: (1) terminal revoke exclusion + J-OWN-UID, (2) the durable-window matrix + pre-digest compatibility, (3) the paired drain/upgrade, (4) the caller. It also required a §5 table for every modified existing test (the iteration-440 B1 lesson). Race budget: `cmd/ailang-worldd` ≤ 12 s, `host/daemon` ≤ 5 s.

**Planner** (~78 min, rc 0). Plan + JSON `ce74776b`.
- **Outcome B.** The compiled activation prototype `st-M4` measured 2,302 cumulative, over 1,650.
- S5f3 = M1 289, M2 686, M3 1,385 cumulative. Appendix C: S5f4 = the caller, 917 on M3.
- Race: S5f3 +3.88 s CLI, +0 s daemon. S5f4 +10.96 s CLI.
- 65/65 diff-anchored and 3/3 carried mutants killed. The only existing test edited was `TestOperatorServiceUpgrade` (one added arm).
- Controller round trip `diff -r` rc 0 for M1, M2 and M3.

**Executor:**
- Run 1 applied M1 and M2 with `ACCEPT_OK`, then stopped correctly when M3's in-sandbox acceptance went red on `TestOperatorServiceStartup/gate/*` (`context deadline exceeded`).
- Outside the sandbox: that test passed `-count=3` at base and at head, and `accept.sh M3 fast` printed `ACCEPT_OK` (60 PASS). So the sandbox invented the red (false-green 3). I took the M3 snapshot with the plan's own `snapshot.sh` (18/18, cmp-equal to the stage).
- `commit_ms.sh` attempt 1 refused M1 on the existing `TestSessionRevokeByEpisode`: its setup mint hit a deadline at 6.15 s. `-count=5` passed at base (0.10–3.53 s) and at M1, so that was load. The retry committed M1 `7b7acdf6`, M2 `b56285dc` and M3 `211524ae` (NUMSTAT 289 → 686 → 1,385, exactly the plan's).

**Baseline instrument note:** base `daemon-race` timed out at 20 minutes twice (co-loaded with the planner, then alone after it). There was no valid local daemon base, and CI arbitrated it. The sprint has no `host/daemon` diff.

**Gates** on `211524ae` (no PATH shims):
- Clean legs: build, vet, Linux vet, test compile, five inventories, `verify_ail`, GC probe, caller census (zero non-declaration hits) and `accept.sh M3 race`.
- `host/store` and `host/authority` plain and race: 0 fails.
- **`cmd-plain` added `TestSessionRootBudgets`**, plus four git-fixture tests that pass alone (load).
- `host/verifygate` was not in the sweep.

**Evaluator round 1** (Agent `sonnet`, ~87 min, own worktree): **FAIL 71** (22/15/17/8/9).
- **B1:** M1 renamed the revoke deadline root, so the existing source guard `TestSessionRootBudgets` went red deterministically. Controller-reproduced: PASS at base, FAIL at head. The plan's §1.2 had claimed the root "stays".
- Non-blocking:
  - a refused upgrade left a loaded companion stopped, silently;
  - `service runner stop|uninstall` refuses without a service record;
  - the process census matches on name only;
  - 8 survivors on the live upgrade path and 2 whole-store survivors.
- 8 of 8 planner kills re-confirmed.

**Round 2** (codex, 169 lines, 37 of them production) → `81789324`.
- `runSessionRevoke` names its root `revokeCtx` again; the guard is unchanged.
- `upgradeFailureNext` appends the `service runner start` remedy when a previously loaded companion is left stopped.
- New tests kill six of the judge's exact edits: J-ROLLBACK-SCHEMA, J-RESTART-LOADED, J-PROMPT-HASHES, J-POST-CENSUS, J-DISPOSE-OWNER and J-RECOVERY-SINGLE-STEP.
- Controller round-2 gates: every leg green. The `cmd-plain` set equals the base count, and the one swap (`TestOperatorWorktreeWindow`) passes alone at base and head (7.66 / 7.67 s).
- **Evaluator round 2** (fresh Agent `sonnet`, ~9 min): **PASS 83** (27/25/18/8/5), 0 blocking. B1 and N1 CLOSED: the guard re-detects the round-1 form. All six kills reproduced. Its O-DROP-LOADED survives.

**Round 3** (CI-found):
- PR #687 CI on `81789324` was red in both Go jobs on one test: `host/verifygate TestSubprocessSinkGate`. M3's `operator_pair_upgrade_test.go:81` used a bare `bytes.Buffer` as a subprocess stderr sink. Locally PASS at base, FAIL at head.
- Codex fix, test-only, 2 lines (`syncBuffer`) → `ed135376`, 1,527 lines.
- The round-2 judge, resumed by name, confirmed the diff and the gate. PASS 83 stands.

**Gate 3b:**
- PR #687 on `ed135376`: attempt 1 red only on the race job's 480 s bound (`host/daemon` 480.134 s, no diff; `cmd/ailang-worldd` 478.231 s). Attempt 2 was 3/3 (275.1 s / 358.2 s).
- `MERGEABLE CLEAN`; merged `--match-head-commit` → `a50b0d72`.
- Merge CI attempt 1: red only on the race bound again (`host/daemon` 480.095 s, `panic: test timed out after 8m0s`; `cmd/ailang-worldd` 443.9 s). Attempt 2: **3/3** (`cmd/ailang-worldd` 273.9 s, `host/daemon` 358.8 s). **LANDED.**

**Found:**
- **Two reds in one sprint came from whole-repo gates my sweep did not run.** One was a whole-package `cmd` guard (B1); the other was a static `host/verifygate` scan over every test file. The planner and executor ran named selectors only, and my gate sweep covered touched packages only. Rule 3g says to derive the CI command list instead of remembering it. Next time, add `host/verifygate` (or the whole CI list) to the controller sweep.
- The race bound fired on both the PR head and the merge, each time on the untouched `host/daemon`. That makes six instances.
- Outcome B gave the first measured end to the deferral chain: S5f4 is the caller, at 917 lines.

**Progress:** 1.0: clauses 1–7 MET; goal unmoved. Row 243 V4: S1–S5f3 landed.
- Every activation prerequisite is now in.
- The caller is planned and measured as one sprint (S5f4).
- Still dormant: no production caller.

**Ruled out:**
- Fixing the four remaining judge survivors (J-FROM-BINARY, J-BUDGET-EXPIRED, J-RUNNER-STOP-CANON, J-REVOKE-MISSINGSTAT) in round 2. They would have kept growing a sprint already over the cap; they are carried to S5f4.
- Treating `TestOperatorServiceStartup`, `TestSessionRevokeByEpisode`, `TestOperatorWorktreeWindow` or the daemon-race drift as regressions. Each passes alone at base and head, and none has a diff on its path.
- Editing the `TestSessionRootBudgets` guard to match: the production root was restored instead.

**Next:** V4 **S5f4 — activation, the first production caller**:
- `runOperatorRunner` takes one confirmed start per tick through `operator_preparation.go`, ending at prepared-awaiting-S6;
- restart adoption and the attended gap remedies;
- `TestOperatorFirstManagedEpisodeProcess` and `TestOperatorCurrentOwnerProcessRecovery`.

The planner must re-measure `world-iter441/plan-scratch/M4.patch` on `a50b0d72` first. Its gates must include whole `cmd/ailang-worldd` and `host/verifygate`, and the added CLI race (~11 s) must be reported against the bound. **Row 267 first** if Gate 1 finds `dev` red: under D-WORLD-121's default B, a red tip triggers it.

**Gate 5 bookkeeping:**
- Rotation not due (#635 opened 2026-11-23; next 2026-11-30 07:00 local). Weekly external-issue sweep not due.
- Idle fingerprint NOT recorded (this fire picked and landed).
- The last 3 landings (441, 440, 439) move no clause. No clause is UNMET, so there is no DRIFT line.

## 442 — 2026-11-27 — row 243 V4 sprint 5f4 LANDED (PR #689, merge `c3504c56`): ACTIVATION — the first production caller of the operator start path: the installed `operator-runner` now visits at most one confirmed start per bounded tick through `operator_preparation.go` (kernel lock → episode lock → candidate re-read → lifecycle exclusion across stop check and claim → worktree reservation + Git readback → atomic start/session/mint → file publication) and stops at durable prepared-awaiting-S6, with read-only restart adoption of the seven durable windows, the attended `service runner reconcile-start` remedy and the partial-files orphan revoke; the deferral chain is closed; 1,617 lines; judge sonnet round 1 FAIL 68, round 2 PASS 85 [PRODUCT]

**Pick:** row 243 V4 sprint 5f4, the Next that record 441 named.
- **D-WORLD-121** (row 267 first?) is still OPEN. Its default B applied: dev's tip `ad1a482c` finished 3/3 green at attempt 1, so the red-tip trigger did not fire.
- **Re-measured at pick** (`origin/dev` `ad1a482c`): no S5f4 plan, branch, PR or worktree existed. Four start-path declarations and zero callers. Iteration 441's `M4.patch` applied to a fresh export except 1 of 1 hunk in `host/store/operator_start_recovery_test.go`.
- **Clause map:** 1.0 clauses 1–7 MET, none UNMET. Row 243 serves D-WORLD-94.

**Gate 0:** stamped first; `gh` = `sunholo-voight-kampff`; billing CLEAN; kill switch absent; 0 allowlisted directives on #635 (22 comments, LEGACY watermark); step 6a: 0 crash and 0 transcript notices (control 107:4); `mission-world` inbox 0 unread.

**Gate 1:** local `dev` was 11 behind `origin/dev`, so the record was read from origin. Base `ad1a482c…@2026-11-26T22:07:18Z`. Skill drift: resolved symlink and pin copy byte-identical to ailang `origin/dev` on all 14 files. Dev CI at `ad1a482c` 3/3 green at attempt 1.

**Routing evidence:**
- Resolver (followed verbatim): planner `recipe codex:gpt-6.2-sol declared:planner-lane-default-pin`, executor `recipe codex:gpt-6.2-sol declared:provider-pin`, evaluator `agent-tool sonnet declared:alias-pin`. Designer not spawned (reviewed design). Codex probe rc 0.
- The operator's standing request was to spawn every role with the Agent tool "exactly as the routing table specifies". The table returned `recipe` for the planner and executor, which makes the Agent tool not a valid path for them, so they ran as `codex exec` recipes. The REQUIRED judge ran on the Agent tool: round 1, a fresh round-2 instance, and that instance resumed by name for a round-3 addendum. No role failed to spawn.
- Controller `claude-opus-5-5` (session; tok: not reported) · planner `codex:gpt-6.2-sol` (609,798 tok) · executor `codex:gpt-6.2-sol` (run 1 70,401 tok; round 2 231,842 tok; round 3 tok: not captured) · evaluator Agent `sonnet` (r1 91,556 tok; r2 fresh instance 74,419 tok; r3 addendum 83,213 tok). judge-independence: cross-vendor. metered=$0.00.

**Directive (controller's):** the caller lands in S5f4. (A) the complete unit within 1,400 (safe-unit exception to 1,650); if the rebased unit measured over 1,650, only the attended remedies could move to a follow-up, never the caller. It carried judge-441's notes to decide, required a §5 table for every edited existing test, a §1.2 list of every whole-repo guard the diff touches, whole `host/verifygate` on every stage, and whole `cmd/ailang-worldd` by fail set. Race budget: `cmd/ailang-worldd` ≤ 12 s, `host/daemon` ≤ 5 s.

**Planner** (~84 min, rc 0). Plan + JSON `48b4a9ec`.
- **Outcome A**, one milestone: the rebased prototype compiled at 917; M1 measured **1,391**, including both attended remedies. Zero activation lines deferred.
- 136 diff-anchored kills + 1 equivalent (`D-OWNED-LIMIT`); all five carried judge-441 ids killed by named tests.
- §1.4: the prototype read a catalogue generation the live loop never fills, so M1 re-reads the durable heartbeat.
- Race: lowest focused sample 12.73 s, 0.73 s over budget, with samples to 58 s under load. Daemon +0.
- Controller round trip `diff -r` rc 0 for `cmd`, `host`, `docs`.

**Executor:** codex applied M1 cleanly, `ACCEPT_OK` (125 PASS), snapshot 14/0/14 cmp-equal to the stage.
- `commit_ms.sh` attempt 1 REFUSED: the NEW `TestOperatorFirstManagedEpisodeProcess` failed at 24.99 s (`production tick failed: 0 operator start blocked`) beside existing git-fixture 30 s timeouts. Alone it passed `-count=3` (3.3–6.1 s). The retry committed **M1 `1602bc5b`**. I recorded it as a candidate load-sensitivity for the judge rather than retrying past it.

**Baseline instrument note:** base `cmd-plain` and `daemon-race` at `-timeout 20m` both timed out under planner co-load (the third iteration running). A second base `cmd-plain` at `-timeout 45m` completed (696 s, 13 fails) and was the comparator.

**Gates** on `1602bc5b` (no PATH shims): gofmt, build, vet, Linux vet, test compile, `verify_ail`, caller census (four call sites in `operator_preparation.go`, one entry in `operator_runner.go`), `host/store` and `host/verbs` plain/race all clean; whole `cmd/ailang-worldd` 0 new fails against base; `accept.sh M1 race` red only on existing `TestOperatorBorrowedWorktreeGitProcess/during_git` (load).

**Evaluator round 1** (Agent `sonnet`, ~65 min, own worktree): **FAIL 68** (25/13/16/8/6).
- **B1:** the new test ran the real loop under a real `budgetContext(22 * time.Second)`. With `GOMAXPROCS=1 -race` and 32 busy loops it failed at 43.70 s with the controller's exact message. A new test must not be load-sensitive.
- Non-blocking survivors: J08 (post-prompt `before != after` in `reconcileOperatorStart`) and J09 (companion-unloaded check), both past the broad selector; J04, J10, J19 past the five new tests; J18 killed only by existing load-sensitive tests.
- No double-effect hole across the seven windows; lock order holds; nothing of S6 reachable; 11 sampled planner ids re-confirmed.

**Round 2** (codex, 250 lines over 15 files) → `c1f848de`.
- The tick budget is the named constant `operatorPreparationBudget` (22 s). A test-only `serviceEnv.operatorBudget` (zero in production) raises the tick, its child budgets and the per-git cap together, via a context key and a new exported `host/worktree.WithCommandTimeout`.
- New pins kill J04/J19 (`TestOperatorStartMintRecoveryRemedies`), J08/J09 (`TestOperatorAttendedStartDisposition`), J10 (`TestOperatorFirstManagedEpisodeProcess`) and J18 (`TestOperatorPreparationRechecksHeartbeatAfterWorktree`).
- Controller round-2 gates: every leg rc 0 except `accept-race` (existing `TestOperatorTerminalRevokeExclusion/durable_owner`, a 4 s mint-fixture deadline under load); whole `cmd/ailang-worldd` 0 new; `host/worktree` 0 fails.
- **Evaluator round 2** (fresh Agent `sonnet`, ~37 min): **PASS 85** (28/25/18/8/6), 0 blocking. B1 CLOSED: the contended reproduction passed 79.24 / 76.59 / 28.88 s. All six ids killed.
- Its non-blocking notes: the production budgets are pinned only by constant equality (its J21, 22 s → 1 s, and J22, +1 h on the zero-override path, both survive); `WithCommandTimeout` is an unclamped exported seam with no production caller; one override value replaces every child budget; 1,601 lines is over the 1,400 cap.

**Round 3** (CI-found):
- PR #689 CI on `c1f848de` was red in both Go jobs on ubuntu-24.04, on the five new activation tests, each under a second: `private session directory: path writable by others`. `runnerPath` walks ancestors up to the injected home and refuses `perm & 0022`. The fixture's sessions directory and fake agent sat outside the runner home, so on Linux the walk reached `/tmp` (1777); macOS's private temp root hid it.
- Reproduced on macOS with `TMPDIR` under `/private/tmp` (1777): both messages; PASS without it.
- Codex, test-only, 16 lines: the shared start-rig fixture creates the agent directory and sessions under the runner home → `f250e9f8`, 1,617 lines. The round-2 judge, resumed by name, re-read the diff (no assertion weakened, no refusal-proving test moved), re-ran the seven tests with and without the override and re-killed J08 and J18: PASS 85 stands.

**Gate 3b:**
- PR #689 on `f250e9f8`: **3/3 at attempt 1** (race `cmd/ailang-worldd` 283.984 s, `host/daemon` 358.156 s). `MERGEABLE CLEAN`; merged `--match-head-commit` → `c3504c56`.
- Merge `c3504c56` CI: **3/3 at attempt 1** (run 45931571443; race `cmd/ailang-worldd` 374.780 s, `host/daemon` 421.681 s — a slower runner, under the 480 s bound). **LANDED.**

**Found:**
- **Two of three reds this sprint came from wall-clock or platform assumptions in NEW tests, not from the code.** The 22 s real tick budget made the first activation test load-sensitive; macOS's private temp root hid a Linux-only refusal. My local gates cannot see the second at all; a `TMPDIR` under a 1777 parent reproduces it on macOS in under a second. That reproduction belongs in the controller sweep for any test that drives `runnerPath`.
- The race bound did NOT fire this landing (PR race 3/3 at attempt 1); `cmd/ailang-worldd` rose about 9 s against S5f3's normal-runner PR reading (275.1 → 284.0 s).
- The deferral chain is closed: the runner prepares an approved start on its own and stops at prepared-awaiting-S6.

**Progress:** 1.0: clauses 1–7 MET; goal unmoved. Row 243 V4: S1–S5f4 landed. The installed runner now prepares a confirmed start (worktree, session credential, private files) and stops at prepared-awaiting-S6. Still unbuilt: the agent launch (S6) and the browser route (S7).

**Ruled out:**
- Retrying the first commit attempt's red past the judge: the failing test was NEW, so "load" was a finding against the test, not a flake to wave through.
- Changing `runnerPath` or any production permission check for the Linux red: the refusal is correct; the fixture was wrong.
- Pinning the production budgets behaviourally (judge-442 J21/J22) inside round 3: carried as a note, below.

**Next:** V4 **S6 — launch**: the agent spawn reservation, the launch receipt and launch recovery from prepared-awaiting-S6, per reviewed design §8 (reviewed 890 raw / 1,202 calibrated). Carry judge-442's notes: behavioural pins for the production tick/heartbeat/reconcile budgets (J21/J22 survive), clamp or narrow `worktree.WithCommandTimeout`, and the stuck-owned-action tick alternation. Gates must include a `TMPDIR`-under-1777 run of every new test that drives `runnerPath`, whole `cmd/ailang-worldd` and `host/verifygate`. **Row 267 first** if Gate 1 finds `dev` red (D-WORLD-121 default B).

**Gate 5 bookkeeping:**
- Rotation not due (#635 opened 2026-11-23; next 2026-11-30 07:00 local). Weekly external-issue sweep not due.
- Idle fingerprint NOT recorded (this fire picked and landed).
- Harness share over the last 20 log headlines: 0/20 tagged HARNESS (20 tagged).
- The last 3 landings (442, 441, 440) move no clause. No clause is UNMET, so there is no DRIFT line.

## 443 — 2026-11-27 — row 267 LANDED (PR #691, merge `6447e776`): the CI race leg's bounds are derived from the measured slow runners — per-package 900 s (was 8m), start allowance 995 s, an ordering guard (leg ≥ per-package + Go's 60 s backup kill + allowance = 1955 s), leg 2000 s, step 40 min — and the leg wrapper names what it kills and kills the whole process group; arms R1–R7 and `TestRaceBoundsAreOrdered`; 379 lines; judge sonnet round 1 FAIL 74, round 2 FAIL 77, round 3 PASS 79 [HARNESS]

**Pick:** row 267, not the record-442 Next (row 243 V4 S6).
- Dev's tip `6420426e` (the docs-only record merge #690) went red at attempt 1 on the race job's 480 s per-package bound: `cmd/ailang-worldd` 480.085 s and `host/daemon` 480.072 s, each `panic: test timed out after 8m0s` (run 45936239736). The attempt-2 re-run was green.
- **D-WORLD-121** is still OPEN. Its default B says row 267 runs when Gate 1 finds `dev`'s tip red, under the standing red-dev rule, and that happened. This mission owns the repo.
- **Clause map:** 1.0 clauses 1–7 MET, none UNMET. Row 267 is clause-2 hygiene: World's own CI, not the fleet's loop harness.

**Gate 0:** stamped first; `gh` = `sunholo-voight-kampff`; billing CLEAN; kill switch absent.
- 0 allowlisted directives on #635 (25 comments), read against both watermarks; the older one governs.
- Step 6a rc 0: 0 crash notices (control 107:4) and 1 transcript-scan notice, for fire `world-1795730668-41608` with 2 key/file hits. The key names are only in the private controlplane notice.
- `mission-world` inbound inbox: 0 unread.

**Gate 1:** local `dev` was 26 behind and was fast-forwarded (untracked files only). Base `6420426e…@2026-11-27T04:01:52Z`. Skill drift: both copies byte-identical to ailang `origin/dev` on 14/14 files.

**Routing evidence:**
- The resolver was followed verbatim: designer `recipe claude:claude-opus-5-5 declared:provider-pin`, planner `recipe codex:gpt-6.2-sol declared:planner-lane-default-pin`, executor `recipe codex:gpt-6.2-sol declared:provider-pin`, evaluator `agent-tool sonnet declared:alias-pin`. Probes: claude-sub rc 0, codex rc 0.
- The operator's standing request asked for every role through the Agent tool "exactly as the routing table specifies". The table returned `recipe` for the designer, planner and executor, so they ran as `claude-sub -p` and `codex exec` recipes. The REQUIRED judge ran on the Agent tool: round 1, a fresh round-2 instance, and that instance resumed by name for round 3. No role failed to spawn.
- Token use:
  - Controller `claude-opus-5-5` (session; tok: not reported).
  - Designer `claude:claude-opus-5-5`: r0 64,345 output tokens, ~3.0M cache-read; r1 20,352 output.
  - Planner `codex:gpt-6.2-sol`: 398,347 tok.
  - Executor `codex:gpt-6.2-sol`: 55,613, then 85,755 (round 2a), 76,961 (2b) and 80,169 (2c).
  - Evaluator Agent `sonnet`: r1 96,298 tok; r2 fresh instance 58,002; r3 resumed, 73,942 cumulative.
- judge-independence: cross-vendor. Metered spend: quorum r1 $0.34 and r2 $0.40, so $0.74 against the $5 cap.

**Evidence:** `race_table.txt` lists one row per race-job attempt over the last 40 CI runs: 51 complete attempts, 18 of which hit the 480 s cap, all on slow runners. The runner speed index k reaches 1.585.

**Design** (`w-daemon-race-package-timeout-headroom.md`). The designer's first launch failed in 1 s because `--add-dir` is variadic and swallowed the prompt; the fix was `--`.
- r0 `8256a83e`: option A, a fixed larger per-package bound with strict layer ordering.
- Quorum r1 (gemini-4-carbon, oc-glm-5-3, oc-kimi-k3) BLOCKED 3/3. Before routing the revision I measured C1–C4:
  - C1: the readers of `verify_go.sh`.
  - C2: the plain leg passes no `-timeout`.
  - C3: `ps` does not truncate a 146-character argv under a pipe on darwin.
  - C4: the start allowance.
- r1 `4c76f7c0`.
- Quorum r2: glm PASS, kimi REJECT (the censored `host/daemon` slow tail; fix A 995, L 2000), astra REJECT (an unbounded `ps` before the kill). Both gave concrete fixes and neither disputed the direction, so r2 `c50264d6` is the controller's application of their verbatim fixes under the narrow-refinement carve-out.
- **FLAG:** round 2 seated astra in place of round 1's gemini, because I passed a four-model `--reviewers` pool. The rule is the same three seats every round.

**Planner** (~60 min, rc 0). Plan + JSON `915ba9c0`.
- One milestone, 351 lines, 24 named kills, 0 observed survivors.
- Daemon per-k falls with k across the uncapped mid-class runs (Pearson −0.583), so the r2 values are conservative.
- R6 was UNINFORMATIVE under the sandbox (real `ps` denied).
- Controller round trip `diff -r` rc 0 on all four directories.

**Executor:** codex applied the patch verbatim; snapshot 6/0/8, cmp-equal to the stage.
- `commit_ms.sh` → **M1 `234ac3f6`**. The commit-time test ran outside the sandbox: `PASS 7/7 arms` and whole `host/verifygate` ok (556.7 s).
- Controller gates (outside the sandbox):
  - `accept.sh` `ACCEPT_OK`.
  - `verifygate`/`pins`: no new failures.
  - Whole-repo plain: 37 fails against base 49. The 11 "new" failures are all in packages M1 does not touch and that read none of its files. Alone, 10 of 11 pass; `TestOperatorFilesReconcile` (charter load list) fails alone too.
- PR CI on `234ac3f6`: 3/3 at attempt 1.

**Judge round 1** (~20 min): **FAIL 74**.
- **B1:** R6 is load-flaky, 3/14 under a concurrent `go test` (`grandchild PID control missing`). I reproduced it at 2/10 with 12 spinners.
- Non-blocking survivors: G5 (SIGKILL→SIGTERM), G14 (`|| true` on the arms step), G21 (a doubled wait); the X07/X08 ledger rows are mis-attributed.
- The numbers re-derive from the raw logs.
- **The judge SIGKILLed PIDs 96794/96825.** They were the divergence watchdog and the rig-lock heartbeat subshells of an unrelated GPU finetune run (`ailang/tools/finetune/pilot5_run.sh`, parent 340), and it took them for its own leaks. Training (96824) kept running. I restored both loops verbatim from `pilot5_run.sh:47–50` and `:176`, detached, so the GPU rig lock keeps its heartbeat (the steal threshold is 6 h). The round-2 and round-3 directives forbid killing any process the judge did not start.

**Round 2:** codex, test-only, in a separate worktree so my running gates were not disturbed.
- 2a, readiness + 10 s bound + pins: G5/G14/G21 killed, but contention failed 2/20 (`probe never signaled readiness`). Measured: a freshly written script costs 0.191 s median to exec under load, against 0.006 s for inline `sh -c`.
- 2b, inline probe: 1/40. The fresh judge r2 gave **FAIL 77** (2/32).
- 2c: retry ONLY an attempt that never became ready (≤ 3, each logged), with every product assertion failing at once. Results: 0/40 with 0 retries; `notready` fails after 3; 72 lines.
- Committed → **`b2b44792`**, 379 lines in total.

**Judge round 3** (resumed): **PASS 79** (31/18/18/12), 0 blocking.
- 1 failure in 37 loaded runs, and it was not the readiness race: R6's outer wall window overran at load ~31 (26 s against < 20 s).
- The retry cannot mask a defect: G5, G21, N3 and N8 each fail on a ready attempt with no retry.

**Gate 3b:**
- PR #691 on `b2b44792`: **3/3 at attempt 1**. The race log shows `── legs: race (race bounds: per-package 900s, leg 2000s)`, `PASS 7/7 arms` (R6 passing with Linux procps `ps`, which is AC7), and `-race -timeout 900s (leg bound 2000s)`. `MERGEABLE CLEAN`; merged `--match-head-commit` → `6447e776`.
- Merge `6447e776` CI: **3/3 at attempt 1** (run 45963259258). The runner was slow-class (`host/store` 86.649 s): `cmd/ailang-worldd` 422.919 s, `host/daemon` 452.341 s, `host/verifygate` 393.993 s, with `PASS 7/7 arms`. That is AC8 observation 1/10, green, though no package passed 480 s this time. **LANDED.**

**Found:**
- **A judge that runs contention loops kills by pattern unless told not to,** and on a shared rig the pattern matched someone else's job. The fix is one line in the directive (kill only PIDs you started). It is now in this iteration's round-2/3 directives and in the Gate 5 process fix.
- **Two of the first fix's three instruments were wrong in the same direction:** a fixed wall bound and a freshly written probe file both depended on the rig's speed. Only a precondition-aware arm was load-independent.

**Progress:** 1.0: clauses 1–7 MET; goal unmoved. Row 267 is landed: the race job no longer pays a re-run on the slow-runner class (AC8, the first 10 post-merge attempts, stays observational).

**Ruled out:**
- Retrying the race leg or dropping packages from `-race` (design constraint 1).
- An in-run calibration bound (design §3 B): `host/store` finishes too late and does not order the slow tail.
- A fixed larger R6 bound alone, which still failed 1/40 at 10 s (2b).

**Next:** row 243 **V4 S6 — launch** (spawn reservation, launch receipt, launch recovery from prepared-awaiting-S6; design §8), carrying judge-442's notes (J21/J22 budget pins, the `worktree.WithCommandTimeout` clamp). Carry the judge-443 residual: compare R6/R7 against the wrapper's own reported elapsed time rather than the arm's outer clock (a test-only follow-up, logged here rather than filed as a row, per D-WORLD-109 (2)).

**Gate 5 bookkeeping:** Rotation not due (#635 opened 2026-11-23; next 2026-11-30 07:00 local). Weekly external-issue sweep not due. Idle fingerprint NOT recorded (this fire picked and landed). Harness share over the last 20 log headlines: 1/20 tagged HARNESS (20 tagged), this iteration's row 267 (World's own CI; picked under the red-dev rule, not loop-harness work). The last 3 landings (443, 442, 441) move no clause; no clause is UNMET, so there is no DRIFT line. One process fix (charter Repo Profile: roles kill only PIDs they started); no skill edit (one instance).

## 444 — 2026-11-27 — row 243 V4 sprint 6 LANDED (PR #695, merge `00b4ec55`): LAUNCH — the installed operator runner resumes a durable prepared-awaiting-S6 start under the episode lock, reserves the launch while holding lifecycle exclusion, rechecks session, approver, signed binary and launch dir at dispatch, and runs one supervised confined headless agent with a birth-identity process receipt, a bounded cap and scrubbed history; the worker spawn is `host/agentrun.StartWorker` with a scrubbed env; 1,888 lines; judge sonnet round 1 PASS 73, round 2 PASS 80 [PRODUCT]

**Pick:** row 243 V4 S6, the record-443 Next.
- **Clause map:** 1.0 clauses 1–7 MET, none UNMET.
- **Ledger:** 1 OPEN, D-WORLD-121, moot since row 267 landed (the loop cannot resolve it).
- **Design:** V4 r1 is reviewed (quorum r2 PROCEED 3/3, no absent seat).
- **Prior attempts:** no S6 commit, PR or orphan worktree existed.

**Gate 0:** stamped first; `gh` = `sunholo-voight-kampff`; billing CLEAN; kill switch absent.
- 0 allowlisted directives on #635 (26 comments) and #574 (44 comments).
- Step 6a at the default `--gh-timeout 5` returned rc 2 (`gh timed out after 5s`; a direct read took 1.2 s), which is an instrument failure, not a verdict. Re-run at 30 s: rc 0.
- That re-run found 0 crash notices (control 107:4) and 1 transcript-scan notice, fire `world-1795730668-41608`, 2 hits. The key names are only in the private controlplane notice.

**Gate 1:** local `dev` == `origin/dev` == `b2b0ebab`. Skill drift: 14/14 byte-identical, and the positive control fires. Dev CI on `b2b0ebab` was 3/3.

**Routing evidence:**
- The resolver was followed verbatim:
  - planner `recipe codex:gpt-6.2-sol declared:planner-lane-default-pin`;
  - executor `recipe codex:gpt-6.2-sol declared:provider-pin`;
  - evaluator `agent-tool sonnet declared:alias-pin`;
  - no designer.
- Codex probe rc 0.
- The operator's standing request asked for every role through the Agent tool "exactly as the routing table specifies". The table returned `recipe` for the planner and executor, so they ran as `codex exec` recipes. The REQUIRED judge ran on the Agent tool: round 1, then resumed by name for round 2. No role failed to spawn.
- Token use:
  - Controller `claude-opus-5-5` (session; tok: not reported).
  - Planner `codex:gpt-6.2-sol`: 656,730 tok.
  - Executor `codex:gpt-6.2-sol`: 76,995, then 155,961 (round 2).
  - Evaluator Agent `sonnet`: r1 161,452 tok; r2 resumed, 179,347 cumulative.
- judge-independence: cross-vendor. Metered: $0.00.

**Planner** (~87 min, rc 0). Plan + JSON `4b567e9a`.
- Outcome A: M1 512, M2 1,650. That is the safe-unit exception, exactly at its ceiling.
- Controller re-measure: cmd 1,054 + host 582 + docs 14 = 1,650.
- Round trip: `diff -r` rc 0 on both milestones.
- Planner ledger: 49 literal + 42 conditional + 8 carried kills, 0 named survivors. It reported 153 source anchors without a completed mutation result.
- **Ruled out by measurement:** the controller's directive cited an E3 raw path that no longer exists (`~/.ailang/state/world-iter429/e3.out`). The tracked reading is `design_docs/verification/world-iter429/readings.md`.

**Executor:** codex applied both patches verbatim.
- `ACCEPT_OK M1 fast` (29 PASS) and `ACCEPT_OK M2 fast` (40 PASS). Snapshots are 20/34 paths, cmp-equal to the planner stages.
- `commit_ms.sh` produced **M1 `6398a4a2`** and **M2 `808ce248`**. M2's commit gate ran whole `host/verifygate`: `ok` 510 s.

**Controller gates (outside the sandbox):**
- Clean: gofmt, build, vet, linux vet, test-compile, `verify_ail`.
- New tests with `TMPDIR` under a 1777 parent and at `-count=3`: 0 fails.
- store, verbs, worktree, agentrun, authority: 0 new against base.
- `accept.sh M2 race`: first run rc 1 on `TestOperatorStartOnce/worktree`, a git 30 s fixture deadline under load. Re-run: `ACCEPT_OK`.
- Whole `cmd/ailang-worldd` plain on M1/M2 had extra fails: `TestDemoSeedCommitFailures`, `TestDemoSeedPreflightRefusals`, `TestOperatorWorktreeWindow`, `TestOperatorFilePublicationAdmission`. All pass alone at head, so they are load.
- Local whole-package race legs truncated at 15 min for base and head alike; CI arbitrates.
- Whole-package mutants: 9 of the planner's 25 (the non-cmd ids) replayed, 9 KILLED.

**Judge round 1** (~57 min): **PASS 73** (25/17/16/8/7), 0 blocking.
- Launch safety holds: there is no second child, restart never adopts a live process, the pre-exec rechecks are present, and scrubbing was tested through real child output.
- Survivors:
  - T1: the real worker/supervisor path is never executed (J16, J17).
  - T2: J09 (signed-binary hash recheck), J10 (launch dir) and J22 (`Before` error).
  - T3: J25 (`Started` error leaves the group alive).
- I reproduced J09 (`if false && …`; launch selector `ok` 5.6 s) and T1 (0 test references, against production references at `operator_launch.go:212,215,268` and `main.go:396`).
- Reachability: `store.Credentials.ConfirmOperatorAction` has no non-test caller, so the launch is dormant until S7. On that basis I chose to land.

**Gate 3b, round 1:** PR #695 on `808ce248` was **RED** on both legs, on one test only: `host/broker TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential` ("14 files … but 13 have drivers").
- The worker re-exec `exec.Command` at `operator_launch.go:217` was undriven and set no `cmd.Env`, so the worker inherited the runner's whole environment.
- Dev `b2b0ebab` was 3/3 green, so the red belongs to this sprint.

**Round 2:** codex wrote 250 lines → **`b40e2f48`**.
- `agentrun.StartWorker` sets `cmd.Env = childenv.Scrubbed(os.Environ())` and Setpgid. `cmd/` now has 0 exec sites.
- `driveAgentrunWorker` was added to the guard's driver list.
- New tests:
  - `TestOperatorLaunchDispatchRechecks/{binary_digest,nonempty_directory}`;
  - `TestRunObservedBeforeRefusalPreventsStart`;
  - `TestRunObservedStartedRefusalKillsGroup`;
  - `TestStartWorkerEnvironmentAndFiles`.
- Controller re-run outside the sandbox: full guard `ok` 24.2 s.

**Judge round 2** (resumed): **PASS 80** (26/24/17/9/4), 0 blocking.
- All five pin mutants were re-run, and each fails its named test. The new tests take 0.00–0.01 s, with no wall-time dependence.
- Scope scored 4/10: round 2 took the sprint to **1,888**, past the 1,650 exception. **FLAG.**

**Gate 3b, round 2:**
- PR on `b40e2f48`: **3/3 at attempt 1** (run 45999595851). Race times: `cmd/ailang-worldd` 287.9 s, `host/daemon` 356.3 s, `host/verifygate` 325.5 s, `host/broker` 278.3 s.
- Merged `--match-head-commit` → **`00b4ec55`**.
- Merge CI: **3/3 at attempt 1** (run 46001578504; race `cmd/ailang-worldd` 291.0 s, `host/daemon` 358.6 s, `host/verifygate` 323.0 s, `host/broker` 275.5 s). **LANDED.**

**Found:**
- **Rule 3g, third consecutive instance:** a whole-repo guard outside the touched packages caught this sprint.
  - 441: `TestSessionRootBudgets` / `TestSubprocessSinkGate`.
  - 442: the Linux fixture.
  - 444: the `host/broker` AC10(a) guard.
  - My controller sweep ran only the packages the diff touched. A launch sprint adds an exec site, and the repo has a guard for exactly that.
- **The guard and the judge found the same hole independently.** The judge's T1 ("the worker path is never executed") and the CI red ("an undriven subprocess site") pointed at the same line, `operator_launch.go:217`.
- **My own instrument error:** a Python one-liner `open(a,"w").write(open(a).read()+…)` truncated `world-mission-status-archive.md` before reading it, which is an 874-line mass deletion. `git diff --stat` caught it before commit, and the file was restored from git. A shell heredoc with unquoted backticks also corrupted the first draft of the row-243 text.

**Progress:** 1.0: clauses 1–7 MET; goal unmoved. Row 243 V4 now runs a start through to a supervised agent launch. It stays dormant until S7's confirmation routes.

**Ruled out:**
- Landing on judge r1 PASS 73 without round 2: CI made round 2 mandatory.
- Splitting round 2 into a separate PR: the red blocked the merge.
- Replaying the 16 `cmd/ailang-worldd` whole-package mutants locally: about 15 min each, so the judge sampled them by named selectors instead (8 killed).

**Next:** row 243 **V4 S7**: the inert request and guarded action confirmation routes, plus the status/Needs-you projection (design §8). It carries the judge-444 notes:
- `runOperatorLaunchWorker` and the supervisor cap (J17) are unexecuted by tests. Pin them in S7, before the path becomes reachable.
- J06.
- No attended settle verb for a stuck launch row (a worker death leaves the episode fenced).
- The agent child inherits lock fds 3/4.
- S7's gates add whole `host/broker`.

**Gate 5 bookkeeping:**
- Rotation is not due (#635 opened 2026-11-23; next 2026-11-30 07:00 local). The weekly external-issue sweep is not due.
- Harness share over the last 20 headlines: 1/20 HARNESS (row 267, iteration 443).
- The last 3 landings (444, 443, 442) move no clause, and no clause is UNMET, so there is no DRIFT line.
- One process fix: the charter guardrail "gate sweeps include whole `host/broker` when the diff adds an exec site". No skill edit.

## 445 — 2026-11-27 — row 243 V4 sprint 7 LANDED (PR #700, merge `6836c9a9`): INERT OPERATOR REQUESTS — outcome B, because the whole compiled S7 measured 1,892 lines (> 1,650). This sprint lands the launch-worker pins (J06/J17/J18, T1), the `review` outcome on `GET /surface/verb.json`, the inert `POST /surface/operator-actions` behind every `approvalRoute` guard, the escaped card, `status.json` and the Needs-you count. The confirm URL stays 404 until S7b (517 lines, compiled). 1,559 lines; judge sonnet round 1 FAIL 68 (B1), round 2 PASS 85 [PRODUCT]

**Pick:** row 243 V4 S7, the record-444 Next.
- **Clause map:** 1.0 clauses 1–7 MET, none UNMET.
- **Ledger:** 0 OPEN. Attended groom #697 (`216c4c42`) resolved D-WORLD-121 as moot mid-iteration, and filed rows 272/273 "after row 243 V4".
- **Design:** V4 r1 is reviewed (quorum r2 PROCEED 3/3).
- **Prior attempts:** no S7 commit, branch, PR or orphan worktree existed.

**Gate 0:** stamped first; `gh` = `sunholo-voight-kampff`; billing CLEAN; kill switch absent.
- 0 allowlisted directives on #635 (28 comments; LEGACY watermark).
- Step 6a rc 0: 0 crash notices (control 107:4).
- 2 transcript-scan notices in window: fire `world-1795730668-41608` (2 hits, already reported by 444) and fire `world-1795769251-75944` (1 hit, whole-file fallback). The key names are only in the private controlplane notice.

**Gate 1:**
- Local `dev` was 2 behind; `origin/dev` = `a83d40a6`, read from origin.
- Skill drift: 14/14 byte-identical.
- Dev CI on `a83d40a6` was in progress; parent `2b3ba7af` green.
- Only own-account PR: #617 (`ui/live-branch-heads`, attended, not this loop's).

**Routing evidence:**
- The resolver was followed verbatim:
  - planner `recipe codex:gpt-6.2-sol declared:planner-lane-default-pin`;
  - executor `recipe codex:gpt-6.2-sol declared:provider-pin`;
  - evaluator `agent-tool sonnet declared:alias-pin`;
  - no designer.
- Codex probe rc 0.
- The standing request asked for every role through the Agent tool "exactly as the routing table specifies". The table returned `recipe` for the planner and executor, so they ran as `codex exec` recipes. The REQUIRED judge ran on the Agent tool: round 1, then resumed by name for round 2. No role failed to spawn.
- Token use:
  - Controller `claude-opus-5-5` (session; tok: not reported).
  - Planner `codex:gpt-6.2-sol`: 995,656 tok.
  - Executor `codex:gpt-6.2-sol`: 73,587, then 177,288 (round 2).
  - Evaluator Agent `sonnet`: r1 150,245 tok; r2 resumed, 173,014 cumulative.
- judge-independence: cross-vendor. Metered: $0.00.

**Baseline** (controller, pristine `a83d40a6`, outside the sandbox, while the planner ran): store/authority/worktree/verbs 0 fails; verifygate 1; broker 4; cmd 17; daemon-race 16 (load-inflated).

**Planner** (~95 min, rc 0). Plan + JSON `d29214c3`.
- **Outcome B:** the complete compiled candidate measured 1,892 lines.
- M1 173: the real `runOperatorLaunchWorker` e2e with a fake child, the supervisor cap, the birth guards, and a QUICKSTART `uncertain` note.
- M2 1,226 (1,399 cumulative): the inert routes.
- S7b: the compiled M3, 517 lines, in Appendix P (not applied).
- Controller round trip: `diff -r` rc 0 for M1 and M2.
- Planner corrected the inherited `OperatorAction` read census (six → seven).

**Executor:** codex applied M1/M2 verbatim.
- `ACCEPT_OK M1/M2 fast`; snapshots 5/29 paths, cmp-equal to the stages.
- Planner `outside.sh` on st-M2 outside the sandbox: `ACCEPT_OK M2 outside` + `outside-race`.
- Commits: **M1 `80af7f36`**, **M2 `e27c93ef`**.

**Controller gates and two commit-gate reds:**
- **M2 attempt 1** refused on `TestOperatorAdmissionHTTP` (`context deadline exceeded` `:332`, 11.02 s). It passes 10/10 alone. I reproduced it 1/5 under `GOMAXPROCS=2 -race` beside whole `host/store` (`:357`), which makes it a load-sensitive NEW test. Routed to the judge as O2.
- **M2 attempt 2** refused on 4 guard fails absent from the base set. All are attributed, none to the diff:
  - verifygate `TestKnownPositiveDelegates` / `TestReleaseChangeNotice`: also red ALONE at base, and green at head with an isolated `AILANG_STATE_DIR`. **My wrapper (`commit-guards.sh`) did not isolate ailang state; the baseline did.**
  - broker `TestExecProtectedCheckTimeout`: red alone at base too.
  - broker `TestExecProtectedDetectsCreation`: 2/2 alone at head.
  - M2 was then committed on `accept M2 outside`.
- **Sweep on `e27c93ef`:** clean on gofmt, build, vet, linux vet, test-compile, `verify_ail`, the callers census and accept outside-race.
  - store/verbs/approval/authority/verifygate: 0 new.
  - broker: 1 new (`TestExecProtectedCheckTimeout`, base-red alone).
  - **`daemon-new-repeat` RED: `TestSurfaceVerbOutcomes`** (my selector `TestSurfaceVerb.*` caught it; the judge found it independently).

**Judge round 1** (~74 min): **FAIL 68** (22/14/17/8/7).
- **B1 BLOCKING:** M2's `operatorVerb` changed the public answers for plain `episode start w1` (terminal → `refused`, "--preset is required", though the CLI accepts a plain start) and for `episode finish`. That broke the EXISTING `TestSurfaceVerbOutcomes` + `verb-route.golden.json`; the plan never named the test. Reproduced first-party: `ok` at base, FAIL at head.
- O2 was diagnosed as the fixture sharing one 5 s writer context across 15 fsynced commits.
- Catalogue guards J11–J15 untested (survived).
- J29 (cap SIGTERM) and dedup J23/J24, J07/J08 survived.
- Reachability SOUND: request 0 → 1 (`handleOperatorRequest`); confirm/judge 0; claim runner-only; confirm URL 404.
- All 16 planner PENDING listener-bound mutants killed outside the sandbox.

**Round 2:** codex wrote 226 lines → **`eb762a2e`**.
- `operatorVerb` owns only the supported headless start and the legacy `never` tuples; everything else falls through to the old switch. The test and golden are unedited.
- The pending-cap fixture now does one write per row.
- New pins:
  - `TestOperatorReviewRefusesForgedCatalogue/{Store,Agent}`;
  - `TestOperatorLaunchWorker/catalogue/{BuildDigest,Agent,RemoteCredentials}`;
  - SIGTERM asserted;
  - `TestOperatorReplayConflicts/{Argv,Generation}`;
  - `TestOperatorAdmissionBoundaries/PendingProjection`.
- The executor reported 11/11 mutation drills killed.
- Controller gates outside the sandbox: all clean.
  - Daemon selector `-count=3` and `-race -count=2` rc 0.
  - F1 contention 5/5 PASS.
  - cmd new tests `-count=3` and 1777 rc 0; store 0 new; surfacepage/verbs rc 0.
  - Whole `host/daemon` plain: 2 not in the base set, `TestA2AInvokeDeadlineClosesSocket` and `TestRegistryPhaseHasItsOwnClock`. Both pass 3/3 alone at head AND at base, so they are load.

**Judge round 2** (resumed): **PASS 85** (28/26/18/8/5), 0 blocking.
- B1, O2, finding 2 and J29 CLOSED. The judge's own edits for J11–J15, J29, J23, J24, J07 and J08 were all killed.
- Still surviving, non-blocking: J01 (equivalent), J02, J10, J06, J18, J19, J22. None reaches a confirm/claim/mint/launch.
- Scope 5/10: 1,559 > 1,400, under the 1,650 exception that the measured 1,892 complete unit justifies. **FLAG.**

**Gate 3b:** PR #700 on `eb762a2e`: **3/3 at attempt 1** (run 46044386173; race `cmd/ailang-worldd` 503.4 s, `host/daemon` 513.2 s, `host/verifygate` 486.6 s, `host/broker` 404.6 s — a slow runner, all under the 900 s bound). `mergeable` re-read MERGEABLE/CLEAN on that head; merged `--match-head-commit` → **`6836c9a9`**. Merge CI: **3/3 at attempt 1** (run 46046770928; race `cmd/ailang-worldd` 284.9 s, `host/daemon` 363.2 s, `host/verifygate` 324.9 s, `host/broker` 279.0 s). **LANDED.**

**Found:**
- **Rule 3g at its fourth consecutive instance, and the first one caught IN-SWEEP:** the existing whole-package test `TestSurfaceVerbOutcomes` in the touched package. The planner's §5 modified-tests table listed `TestSurfaceVerbIsReadOnly` and missed its sibling. My sweep caught it only because I widened the selector to `TestSurfaceVerb.*`. Together with 444's process fix, that argues for running every test in each touched package whose name shares the touched file's stem.
- **Instrument error, mine:** `commit-guards.sh` (the planner's wrapper, run by me) did not isolate `AILANG_STATE_DIR`, so two verifygate tests went red at head and base alike. The baseline and the gate sweep both isolate it. One wasted guard cycle (~15 min).
- **A design-grounded answer the plan changed silently:** a CLI-valid plain `episode start` became `refused` on the bar. Design §4.1 says existing non-operator commands keep their classifications.

**Progress:** 1.0: clauses 1–7 MET; goal unmoved. Row 243 V4 now has a browser-reachable INERT request path. The S6 launch is still dormant: no production caller of `ConfirmOperatorAction`.

**Ruled out:**
- Landing on the M2 commit's own gate after a single retry without diagnosis: the flake was reproduced and fixed in round 2.
- Re-deciding the plain-start answer: design §4.1 fixes it.
- Running all 25 whole-package mutants: the judge resolved 16 listener-bound ids plus its own 17.

**Next:** row 243 **V4 S7b**: apply the compiled M3 from the S7 plan's Appendix P (517 lines, re-measure on the S7 merge). It registers `POST …/confirm` with its key-ceremony, replay and counter pins and live process readback, making the S6 launch reachable. Carry:
- the judge-445 survivors J06 (launch NoSpawn projection) and J22 (card digest recheck), to pin when confirm lands; J02/J10/J18/J19 noted;
- the planner's S7b gate commands (`gates.sh st-M3`, `outside.sh … M3`, `whole-mutants.sh … mutant-follow-up-ids.txt`);
- every gate sweep runs every test in each touched package whose file stem matches a touched file, plus whole `host/broker` and `host/verifygate` with an isolated `AILANG_STATE_DIR`.

**Gate 5 bookkeeping:**
- Rotation is not due (#635 opened 2026-11-23; next 2026-11-30 07:00 local). The weekly external-issue sweep is not due.
- Harness share over the last 20 headlines: 1/20 HARNESS (row 267, iteration 443).
- The last 3 landings move no clause, and no clause is UNMET, so there is no DRIFT line.

## 446 — 2026-11-27 — row 243 V4 sprint 7b LANDED (PR #703, merge `26b6947f`): CONFIRMATION ACTIVATION — `POST /surface/operator-actions/{id}/confirm` behind every `approvalRoute` guard, the exact key ceremony on the card, every design §5.2 refusal with read-back no-change pins, exact-artifact replay with one counter advance, one transaction for counter use and the state flip, a confirmation lock, and live process birth readback for the status projection; the first browser POST that can make one confined headless start claimable by the installed runner; one milestone, 1,263 lines; judge sonnet PASS 76 [PRODUCT]

**Pick:** row 243 V4 S7b, the record-445 Next.
- **Clause map:** 1.0 clauses 1–7 MET, none UNMET.
- **Ledger:** valid, 108 rows, 0 OPEN.
- **Design:** V4 r1 reviewed (quorum reader rc 0, r2 `…T10-24-51Z` PROCEED 3/3, `absent_reviewers` `[]`).
- **Re-measured at `88fea1e9`:** the store `ConfirmOperatorAction` had no production caller (3-line grep: the two definitions and `judge_operator.go:66`); no S7b plan, branch, PR or worktree.

**Gate 0:** stamped first; `gh` = `sunholo-voight-kampff`; billing CLEAN; kill switch absent.
- 0 allowlisted directives on #635 (31 comments; LEGACY watermark); 0 on #574 and #1 since the mission-scoped watermark.
- Step 6a rc 0: 0 crash notices (control 107:4). 3 transcript-scan notices in the window: fires `world-1795730668-41608` (2 hits) and `world-1795769251-75944` (1 hit), both reported before, and `world-1795786215-30397` (3 hits, new). Key names are only in the private controlplane notice.
- `mission-world` inbox: 0 unread.

**Gate 1:**
- Local `dev` was 8 behind (the 445 record); fast-forwarded to `88fea1e9` (the main checkout holds only untracked files).
- Skill drift: the pin worktree matched origin on all 14 files. The symlinked copy differed on `resources/gate-5-retro.md` only (one line, origin adds two idle-fingerprint inputs); Gate 5 followed origin's version.
- Dev CI on `88fea1e9` (record-only) was in progress; the code parent `6836c9a9` was 3/3.
- No open own-account PRs.

**Routing evidence:**
- Resolver followed verbatim: planner `recipe codex:gpt-6.2-sol declared:planner-lane-default-pin` (`derive-planner-lane.sh` the same); executor `recipe codex:gpt-6.2-sol declared:provider-pin`; evaluator `agent-tool sonnet declared:alias-pin`; no designer. Codex probe rc 0.
- The standing request asked for every role through the Agent tool "exactly as the routing table specifies". The table returned `recipe` for the planner and executor, so they ran as `codex exec` recipes; the REQUIRED judge ran on the Agent tool, foreground. No role failed to spawn.
- Token use: controller `claude-opus-5-5` (session; tok: not reported) · planner `codex:gpt-6.2-sol` 625,032 · executor `codex:gpt-6.2-sol` 74,383 (r1) + 74,396 (r2) · evaluator Agent `sonnet` 84,346.
- Judge independence: cross-vendor. Metered: $0.00.
- Base `88fea1e9a8422b12a059eb1ecfc9c10c5b7b4e7d` (Gate 1 and Gate 3 snap agree, no drift). The sprint and base worktrees used a plain `git worktree add` rather than `mission-worktree.sh`; completeness was asserted after the add (7,474 files each = `git ls-tree` count, `status` clean).

**Baseline** (controller, pristine `88fea1e9`, outside the sandbox, isolated `AILANG_STATE_DIR`, while the planner ran): store/authority/agentrun plain+race and verbs 0 fails; verifygate 1; broker 3; cmd 15; daemon-plain 20; daemon-race 8.

**Planner** (~78 min, rc 0). Plan + JSON `81813d6d`.
- One milestone **M1**, +1,191/−72 = **1,263** (cap 1,400).
- Re-based the stale S7 Appendix P M3 prototype (517 lines, compiled against the pre-round-2 stage): `operator_catalogue_test.go` hunk 2 rejected because the landed catalogue subtests occupy the old EOF, so the literal-task pin moved into `TestOperatorStartBrowser`; `surface_operator.go` −4 and `surface_operator_test.go` +62 offsets.
- Planner census: `JudgeOperatorAction` 0 → 1 and store `ConfirmOperatorAction` 0 → 1, both in `handleOperatorConfirm`.
- Own ledger 120 ids, 0 non-equivalent named survivors; carried J06/J22 (`TestOperatorConfirmProjection`) and J02/J10/J18/J19 (`TestOperatorConfirmationAdmissionNeighbours`) killed by exact edit. These are planner-reported, with route rows judged on an in-process recorder because the sandbox denies loopback.
- Controller round trip on a fresh `git archive 88fea1e9`: `diff -r` rc 0 for cmd, host, docs.

**Executor:**
- **Round 1** stopped before applying anything, correctly. The plan's own §7 block ran `grab "````bash $tool"`; inside double quotes bash treats the backticks as command substitutions, so `accept.sh` extracted to 0 bytes.
- The controller verified a single-quoted form against the planner's scratch copies (3 × `cmp`, patch `cmp`-equal) and re-sent the block with only those two lines changed.
- **Round 2:** patch clean; `ACCEPT_OK M1 fast` (257 PASS, 0 FAIL); snapshot 22 STATUS / 1 DELETED / 21 manifest, 21/21 `cmp`-equal to the planner's stage.

**Commit gate and attributed reds:**
- `ACCEPT_OK M1 outside` on all three attempts.
- The planner's `commit-guards.sh` (whole verifygate + broker) refused twice on broker fails absent from the base set, a different set each time: six tool-binary subprocess timeouts, then three names. `host/broker` imports `host/store`, which M1 changes, so the package argument does not clear them.
- Interleaved alone reruns: the base FAILED `TestAilangRunCapsOnTheToolBinary` once and every name passed on both sides at least once; head failed `TestExecProtectedCheckTimeout` once at 42.6 s against 34.2–34.7 s passes. Negative control fired, so these are load.
- Committed on `accept … outside` alone (the 445 precedent): **M1 `1aa6badc`**, `NUMSTAT 1263`; committed tree = snapshot.

**Controller gates on the stage, outside the sandbox:**
- OK: gofmt, build, vet, linux vet, nine route/header/writer/dependency/subprocess inventories, `verify_ail`, callers, t1777.
- Fail sets vs base: store/agentrun 0 → 0; verifygate 1 → 0; broker 3 → 1 new (`TestAilangToolConfinementMatrix`, the attempt-1 class); cmd 15 → 12 (1 new); daemon-plain 20 → 12 (4 new). All five new cmd/daemon names passed alone at head AND at base.
- The accept legs failed on git-fixture worktree tests (`git … commit … timed out after 30s`) while the judge mutated in parallel. Alone, only `TestOperatorStartCollisions` failed at head. The A/B then had base 8.00 s PASS, head 7.90 s PASS, **base 104.74 s FAIL**, head 76.75 s PASS: rig I/O contention, not M1.
- UNMEASURED locally: whole `cmd/ailang-worldd -race` (the pristine leg itself hit Go's 15 m timeout on the 11 known-red e2e tests) and whole `host/daemon -race` (the controller's background sweep hit the 2 h tool ceiling). CI's race job measured both.
- The three added daemon tests under `-race`, alone, `-count=2`: ~4.7 s per run (plan budget 12 s).

**Judge** (Agent `sonnet`, foreground, ~100 min): **PASS 76** (26/22/17/9/8), 0 blocking.
- Reachability sound: one production caller each of the judge and the store CAS, in `handleOperatorConfirm`; the handler owns no claim, mint, worktree or spawn.
- Refusals asserted by read-back snapshots (row, counters, credentials file, world DB, zero effects/starts/sessions); rollback uses real SQLite triggers; `humanMu` held across read, judge and write (removal killed).
- Own mutants: 18 valid, 15 killed (daemon kills on the real `httptest` listener), 7 survived: E15 the store `BackupEligible` drift check untested; E16/E29 shadowed checks; E10 the judge-time lifetime check shadowed; E19 the darwin zombie check; E24/E27 near-equivalent. Must-live control survived.
- One load flake in the PRE-EXISTING S7 test `TestOperatorAdmissionHTTP` (`surface_operator_test.go:410`, the 5 s per-row write); controller read the diff: the default budget is unchanged and the new tests inject their own.
- Not done by the judge: the six carried ids by exact edit and the planner's recorder-judged route ids on the real listener. Its sequencing file never appeared (the controller sweep's tool ceiling), so it waited ~80 min and then measured on a quiet rig.

**Gate 3b:** PR #703 on `1aa6badc`: **3/3 at attempt 1** (run 46080358867; race `cmd/ailang-worldd` 289.7 s, `host/daemon` 367.3 s, `host/verifygate` 321.6 s, `host/broker` 280.0 s). `origin/dev` had moved to `8805fd7a` (#702, attended: `episode finish` infers `--preset`, a herdr finish key, a surfacefollow fix): 0 of M1's 22 files shared. `mergeable` re-read MERGEABLE/CLEAN on that head; merged `--match-head-commit` → **`26b6947f`**. Merge CI: 3/3 at attempt 1 (run 46081726459, a slow-class runner; race `cmd/ailang-worldd` 511.1 s, `host/daemon` 539.8 s, `host/verifygate` 491.9 s, `host/broker` 422.1 s, all under the 900 s bound). **LANDED.**

**Found:**
- **The planner's own §7 command block was never executed as written.** It rehearsed through its `roundtrip.sh`, not through the block it handed the executor, so a quoting defect that extracted nothing reached the executor. Rule 3e(a) aimed at the plan's apply block; instance 1. The executor stopping instead of improvising was correct.
- **A whole-package local race comparison for `cmd/ailang-worldd` cannot complete on this rig.** The pristine leg itself hit Go's 15 m timeout, so the planner's `gates.sh … outside` comparison is truncated by construction. A local race sweep should `-skip` the 11 known-red e2e tests.
- **Gate legs and the judge's mutation runs contend.** Every broker/cmd/daemon red this iteration came from running two heavy sweeps at once, and every one cleared alone, with base-side failures as controls. Instance 2 of the 430 retro note; instance 3 is a process change.
- **An attended `service upgrade` rewrote the real `~/.ailang-world/service.json` at 20:46:25** while the planner's read-only observation ran. Attributed by process tree to a herdr pane; nothing in the loop touched it.

**Progress:** 1.0: clauses 1–7 MET; goal unmoved. Row 243 V4 now has a browser path from review to an approved, claimable start: the S6 launch is reachable for the first time. S8 (the app bar bindings and Chrome) is next; Mark's no-terminal exit episode stays attended.

**Ruled out:**
- Attributing the broker/cmd/daemon whole-package reds to M1: each name passed alone at head; the base failed the same names under the same load.
- A round 2 for the judge's survivors: none is blocking and none can make an action claimable; E15 and E19 are carried.
- Re-planning the cmd-race comparison mid-iteration: CI's race job measured it (289.7 s).

**Next:** row 243 **V4 S8**, the app bar bindings and Chrome (design §8, 1,161 calibrated). Carry:
- the judge-446 survivors E15 (`BackupEligible` drift) and E19 (darwin zombie birth), plus re-running the six carried ids on the real listener;
- the plan apply block is EXECUTED verbatim by the planner before handing it over;
- local race sweeps `-skip` the 11 known-red e2e tests, and gates never run while the judge mutates.

**Gate 5 bookkeeping:**
- Rotation is not due (#635 opened 2026-11-23; next 2026-11-30 07:00 local). The weekly external-issue sweep is not due.
- Harness share over the last 20 headlines: 1/20 HARNESS (row 267, iteration 443), 20/20 tagged.
- The last 3 landings move no clause, and no clause is UNMET, so there is no DRIFT line.

## 447 — 2026-11-28 — row 243 V4 sprint 8 LANDED (PR #705, merge `85ed2234`): THE FIRST BROWSER SLICE — the standalone operator-action card (key gestures, read-only status reconnect, typed links) and the app bar bindings (read-only Enter review, explicit request, frozen-request reconnect, link allowlist), with a live-Chrome fake Sol start at 390 × 844 that runs the QUICKSTART walkthrough verbatim; the walkthrough's first real run found six defects the sandboxes could not see, fixed in rounds 2–3; 1,268 lines; judge sonnet PASS 85 [PRODUCT]

**Pick:** row 243 V4 S8, the record-446 Next.
- **Clause map:** 1.0 clauses 1–7 MET, none UNMET.
- **Ledger:** valid, 108 rows, 0 OPEN at pick; this record adds D-WORLD-122 (OPEN).
- **Design:** V4 r1 reviewed (quorum reader rc 0, r2 `…T10-24-51Z` PROCEED 3/3, `absent_reviewers` `[]`).
- **S8's design gates (§8 "119 A/B, S7, E5 and #604 sequencing"), re-measured:** 119 = A; S7/S7b landed; PR #604 MERGED 2026-11-18 (`65108353`), so sequencing is resolved against post-merge dev; **E5 needs Mark's chosen preset and has never run.** The loop cannot choose it, so S8 ran on a FAKE Sol preset with E5's loop-runnable half (W25/W32) in planning, and D-WORLD-122 asks Mark for the preset before S11.
- **Re-measured at `0e32a8d7`:** `app.js` had 0 `operator-actions` references; no S8 plan, branch, PR or worktree.

**Gate 0:** stamped first; `gh` = `sunholo-voight-kampff`; billing CLEAN; kill switch absent.
- 0 allowlisted directives on #635 (33 comments; LEGACY watermark) and none since the mission-scoped watermark.
- Step 6a rc 0: 0 crash notices (control 107:4). The 3 transcript-scan notices in the window are the three iteration 446 reported (fires `world-1795730668-41608`, `world-1795769251-75944`, `world-1795786215-30397`); none new.
- `mission-world` inbox: nothing addressed to this mission. Weekly sweep not due (424 ran it); no rotation (#635 created after Monday's boundary, 33 comments).

**Gate 1:** local `dev` = `origin/dev` = `0e32a8d7`. Skill drift: the pin matched origin on all 14 files; the symlinked copy differed on `resources/gate-5-retro.md` only (one line, the idle-fingerprint input list) — Gate 5 followed origin. Dev CI 3/3 green. No open own-account PRs.

**Routing evidence:**
- Resolver followed verbatim: planner `recipe codex:gpt-6.2-sol declared:planner-lane-default-pin` (`derive-planner-lane.sh` the same); executor `recipe codex:gpt-6.2-sol declared:provider-pin`; evaluator `agent-tool sonnet declared:alias-pin`; no designer. Codex probe rc 0; codex 14.0% used / 25.9% allowed.
- The standing request asked for every role through the Agent tool "exactly as the routing table specifies". The table returned `recipe` for the planner and executor, so they ran as `codex exec` recipes; the REQUIRED judge ran on the Agent tool (`sonnet`, background, polled in-turn). No role failed to spawn.
- Token use: controller `claude-opus-5-5` (session; tok: not reported) · planner `codex:gpt-6.2-sol` 727,578 · executor `codex:gpt-6.2-sol` 48,945 (r1) + 90,445 (r2) + 44,953 (r3) · evaluator Agent `sonnet` 124,233.
- Judge independence: cross-vendor. Metered: $0.00.
- Base `0e32a8d7de0bb524678a3c443400cf05d6ec47eb` (Gate 1 and Gate 3 snap agree). Sprint and base worktrees via `mission-worktree.sh` (7,485 files each = `git ls-tree`, clean); the judge's detached worktree via a backgrounded plain add (7,495 files, clean).

**Baseline** (controller, pristine `0e32a8d7`, outside the sandbox, isolated `AILANG_STATE_DIR`, while the planner ran): store/authority/agentrun plain+race, verifygate and verbs 0 fails; broker 4; cmd 14; daemon-plain 14; daemon-race 10.

**Planner** (~77 min, rc 0). Plan + JSON `cc16e985`, outcome A: M1 510, M1+M2 1,269 (cap 1,400).
- It **executed its own §7 apply block verbatim** against a fresh archive (`APPLY_ACCEPT_SNAPSHOT_OK M1 fast` / `M2 fast`, six `ROUNDTRIP_OK`), the iteration-446 carry.
- Planner-reported (not first-party): census 1 → 1 → 1 for the four operator CAS/judge entry points; 171 diff mutants (145 named kills, 26 outside-only); E15/E19 killed; every route and Chrome result UNINFORMATIVE UNDER SANDBOX.
- Controller round trip on a fresh `git archive 0e32a8d7`: `diff -r` rc 0 for cmd, host and docs at both stages.

**Executor round 1** (~6 min): the §7 block verbatim with `TOOLS=$STATE/exec-tools` (the block's own parameter, so planner scratch stayed untouched: 17/17 hashes unchanged). `ACCEPT_OK M1 fast`, `ACCEPT_OK M2 fast`; snapshots 11/0/11 and 16/0/16, all `cmp`-equal to the planner's stages.

**Commit gate:** M1 `ACCEPT_OK M1 outside` → **M1 `86f40c81`** (510). M2 attempt 1 REFUSED: `TestOperatorStartWalkthroughInChrome` `operator_browser_chrome_test.go:193: EOF`, `driver: signal: killed` at 120 s, reproduced alone on a quiet rig. The existing Chrome controls passed at base and on the M2 tree, so the harness and page load were fine.

**Diagnosis (controller, scratch copy, outside the sandbox):** the walkthrough had never run anywhere — both sandboxes deny Chrome.
- **F1 (harness):** the test set `HOME` to an empty temp dir BEFORE starting Chrome; CDP `Page.navigate` then never returned (traced), while the server answered `/surface` in 1.5 ms.
- **F2 (product):** after a lost POST and reload the restore worked (Reconnect, task read-only, request in sessionStorage), but `boot()`'s last line (`app.js:4022`) took focus to `#root`.
- **F3 (product):** the card's external stylesheet was blocked by the inherited `surfaceCSP` (`style-src 'unsafe-inline'` only; violation recorded); the card was 1,025 px wide at 390 px. Fix as the app shell does: inline one `<style>`; the CSS route removed; `surfaceCSP` unchanged.
- **F4/F5 (harness):** two card checks wanted `scrollWidth===390`; a vertical scrollbar gives 375 with nothing overflowing → `<=innerWidth`.
- **F6 (harness):** the driver JSON-stringifies a Promise to `{}`; read the fetch result in two steps.
- With F1–F6 the walkthrough PASSED in scratch (15.2 s).

**Executor rounds 2–3:** round 2 applied F1–F6 as specified (the controller's diff of the snapshot vs the plan's M2: exactly those six, 60 lines, 5 files). **The controller's own F2 spec was wrong** — `#verb` is absent on an ordinary load, so `$('verb').contains` would throw at the end of every `boot()` — caught on review; round 3 changed that one token to `$('verb')?.contains`.

**M2 final gate:** `ACCEPT_OK M2 outside`, `CHROME operator-chrome RAN=1 SKIP=0`, `CHROME chrome-controls RAN=2 SKIP=0`; `commit-guards.sh` refused on two verifygate names (`TestFixtureDiscrimination`, `TestMissionConfigDefaultFlowIntegration`; broker no new fails). Verifygate imports no other World package and S8 touches no script; alone, head PASS (57 s) and base FAIL 3 (80 s) — load. Committed on the outside receipts (the 445/446 precedent): **M2 `8e341d9b`**, NUMSTAT 1,268.

**Controller sweep** (fresh archive of `8e341d9b`; whole-package record legs skipped for time, left to CI and `commit-guards`): every cheap leg and inventory OK, verify_ail OK; `accept` rc 1 only because Chrome did not START within the driver's 30 s wait (load ≈4.8), while `accept-race` and the 1777 leg (same walkthrough) passed. That failed start orphaned a Chrome, which the controller identified by profile path and stopped. E15/E19 killed named and whole; the rest of `whole-mutants.sh` stopped (each id runs whole `host/daemon`, ~19 min).

**Judge** (Agent `sonnet`, ~54 min): **PASS 85** (27/22/18/9/9), 0 blocking.
- Reachability sound: no page path creates or confirms an action without an explicit click; the census adds no caller; every new navigation passes the anchored allowlist.
- 24 of its own mutants killed (JS only by executing it), 9 survived; 6/6 planner outside-only ids killed on the walkthrough; E15/E19 and J02/J06/J10/J18/J19 killed on the real listener (J22's original edit no longer compiles; variant J22b killed). Chrome PASS twice, 0 SKIP.
- Survivors / notes (non-blocking): a DELAYED auto-click on reload or on the card (C-09/C-10/C-10b) survives the walkthrough — no such product code, a test gap; Chrome-only DOM wiring (C-01 liveness argument, C-02 draft body); the card's `state !== 'requested'` guard (S8-14) is defence in depth; an outstanding draft overrides a freshly reviewed panel.

**Gate 3b:** PR #705 on `8e341d9b`:
- Attempts 1 and 2 RED on `host/daemon TestA2AInvokeDeadlineClosesSocket` only (`invoke_e2e_test.go:412`); it passed 3/3 alone locally.
- No S8 path: the only daemon production change is a route and an approval-clock seam used only by `surface_approve.go`, and S8's daemon tests sort after the A2A file.
- **CI negative control:** `gh workflow run CI --ref dev` at base `0e32a8d7` and head attempt 3, same window — both host jobs green. Base 2/2, head 1/3; no other failure in the last 40 runs. Recorded as a flake, instances named below.
- 3/3 at attempt 3 → MERGEABLE/CLEAN, `origin/dev` unmoved → merged `--match-head-commit` → **`85ed2234`**. Merge CI: 3/3 at attempt 1 (run 46106233909) → **LANDED**.

**Found:**
- **A browser gate that only runs outside the sandbox is first executed at the commit gate, so a plan can be internally consistent and still fail on first contact.** The planner's 171-mutant ledger, its own verbatim §7 run and the executor's acceptance were all green, while the walkthrough held three product defects (focus, CSP-blocked CSS, overflow) and three harness defects. Process fix below.
- **The shared CDP driver leaks Chrome when its 30 s start wait expires** (it never learns the PID); a 16-day-old orphan of `TestPasskeyPageRecordsTheFixtures` is the same class. Not ours to kill; reported.
- `TestA2AInvokeDeadlineClosesSocket` flake instances: 445 (local whole-package, load), 447 CI attempts 1 and 2 at one SHA. Logged per D-WORLD-109 (2); one more instance makes it a row under 258.
- The judge's own daemon `-race` selector saw `TestWorkspaceProtectedPolicyRequests` fail once (passes alone 4.6 s; no S8 code on its path).

**Progress:** 1.0: clauses 1–7 MET; goal unmoved. Row 243 V4 now has its first browser slice: from the address bar, through review, request and key approval on the card, to one confined fake launch observed in a real Chrome. S9 (verified private F1 entry and auth readiness) is next. D-WORLD-122 (Mark's Sol preset) is owed before S11.

**Ruled out:**
- Attributing the A2A CI red to S8: no production path; base green twice in the same window; head green at attempt 3.
- Attributing the verifygate commit-guard fails to S8: no dependency, base fails alone.
- Landing M1 alone after M2's first red: M1's card page had no Chrome coverage of its own, and the design gates every page change on live Chrome.
- A CSP change for the card stylesheet: inlining matches the app shell and leaves `surfaceCSP` unchanged.

**Next:** row 243 **V4 S9** — verified private F1 entry and auth readiness, frozen manifest, sandboxed checks and local commit journal (design §8, 1,215 calibrated; waits on 120 = A ✓, E2 read-level PASS (W33), S8 ✓). Carry:
- judge-447 C-09/C-10/C-10b (delayed auto-click) and C-01/C-02;
- the process fix (the controller runs the plan's outside-only browser/listener legs on the planner's stage BEFORE the executor);
- D-WORLD-122 before S11.

## 448 — 2026-11-28 — row 243 V4 sprint 9 LANDED (PR #707, merge `c1869819`), outcome B: THE FIRST FINISH SLICE — public finish keeps F1 first and its typed tty path unchanged, a sealed prep-proof preflight guards the NEW private finish entry, the runner's attended install gains the E2-equivalent finish readiness probe, and the frozen finish manifest, no-follow worktree digest, repo-config fence and exact local-commit readback land with no caller yet; S9b (1,400 lines) compiled in the plan; 1,399 lines; judge sonnet PASS 74 [PRODUCT]

**Pick:** row 243 V4 S9, the record-447 Next.
- **Clause map:** 1.0 clauses 1–7 MET, none UNMET.
- **Ledger:** valid, 109 rows, 1 OPEN (D-WORLD-122, Mark's Sol preset for E5 — an S11 input, not S9's).
- **Design:** V4 r1 reviewed (quorum reader rc 0, r2 `…T10-24-51Z` PROCEED 3/3, `absent_reviewers` `[]`).
- **S9's design gates (§8 "119 A/B and 120 A/B, E2, S8"), re-measured:** D-WORLD-119 and D-WORLD-120 RESOLVED (A); E2 read-level PASS (W33); S8 landed (`85ed2234`). No S9 plan, branch, PR or worktree existed.
- **Verified at `9a40fe68`:** `runEpisodeFinishWith` at `episode_finish.go:63`, the F1 comment at `:92`; `cmd host docs app` byte-identical to `85ed2234`.

**Gate 0:** stamped first; `gh` = `sunholo-voight-kampff`; billing CLEAN; kill switch absent.
- 0 allowlisted directives on #635 (34 comments; LEGACY watermark), none since the mission-scoped watermark.
- Step 6a rc 0: 0 crash notices (control 107:4). The 3 transcript-scan notices in the window are the three iterations 446 and 447 reported; none new.
- `mission-world` inbox: nothing unread. Weekly sweep not due; no rotation (next Monday 2026-11-30).

**Gate 1:** local `dev` = `origin/dev` = `9a40fe68`. Skill drift: the pin matched origin; the symlinked copy differed on `resources/gate-5-retro.md` only (the idle-fingerprint input list) — Gate 5 followed origin. Dev CI 3/3 green. No open own-account PRs.

**Routing evidence:**
- Resolver followed verbatim: planner `recipe codex:gpt-6.2-sol declared:planner-lane-default-pin` (`derive-planner-lane.sh` the same); executor `recipe codex:gpt-6.2-sol declared:provider-pin`; evaluator `agent-tool sonnet declared:alias-pin`; no designer. Codex probe rc 0; codex 15.0% used / 25.9% allowed.
- The standing request asked for every role through the Agent tool "exactly as the routing table specifies". The table returned `recipe` for the planner and executor, so they ran as `codex exec` recipes; the REQUIRED judge ran on the Agent tool (`sonnet`, foreground). No role failed to spawn.
- Token use: controller `claude-opus-5-5` (session; tok: not reported) · planner `codex:gpt-6.2-sol` 819,311 · executor `codex:gpt-6.2-sol` 133,203 · evaluator Agent `sonnet` 120,903.
- Judge independence: cross-vendor. Metered: $0.00.
- Base `9a40fe68e45cd9b5c787aecd716b047327bddb22` (Gate 1 and Gate 3 snap agree). Sprint, base, judge and record worktrees via `mission-worktree.sh`.

**Baseline** (controller, pristine `9a40fe68`, outside the sandbox, isolated `AILANG_STATE_DIR`, while the planner ran): store/authority/agentrun plain+race, verifygate and verbs 0 fails; broker 6; cmd 13; daemon-plain 4; daemon-race 14. Base-extra: checkrun/worktree/surfacepage plain+race 0; whole `cmd -race` timed out (the known local limit).

**Planner** (~96 min, rc 0). Plan + JSON `a1a08bae`, **outcome B**: compiled whole 2,799 lines; S9 = M1 621 → M2 1,399 (cap 1,400); M3 = S9b, 1,400 incremental, compiled in Appendix P with its activation gates listed.
- Executed its own §7 block verbatim on a fresh export (`APPLY_ACCEPT_SNAPSHOT_OK` M1/M2/M3, `FINAL_APPLY_RC 0`).
- S9/S10 boundary: S9b ends at a `local-commit-done` receipt with an empty `v3_id`; S10 issues the V3 `VerbRequest`.
- Controller round trip on a fresh `git archive 9a40fe68`: `diff -r` rc 0 for cmd, host and docs at both stages. Destructive-write scan of the plan scripts: clean.

**Outside legs on the planner's stage, before the executor (the iteration-447 process fix, first use):** M2 `outside`, `outside-race` (cmd added 8.35 s, cap 15 s) and `outside-1777` all `ACCEPT_OK`. M1: one fail, `TestOperatorFinishAuthInstall` in the 1777 leg (3.30 s, "could not bootstrap the companion"); alone 3/3 and in its leg 2/2 PASS. Not reproduced, so the executor started and the judge got it as a mandatory question.

**Executor** (~9 min): the §7 block verbatim with `TOOLS=$STATE/exec-tools`, M1 and M2 only. `ACCEPT_OK` ×2; snapshots 9/0/9 and 15/0/15, `cmp`-equal to the planner's stages.

**Commit gates:**
- **M1 `9a04073e`** on `ACCEPT_OK M1 outside` (621).
- M2 attempt 1 REFUSED: `INSTRUMENT FAILURE cmd-finish-stems truncated` at the 240 s bound in an EXISTING finish test, while the controller's own whole `cmd -race` base leg ran beside it. The controller stopped its own sweep driver.
- M2 attempt 2 (quiet rig): `ACCEPT_OK M2 outside`, then `commit-guards.sh` REFUSED `truncated commit guard` — it runs whole verifygate+broker with no `-timeout`, so broker hit Go's 10 m default (pristine 730 s). Verifygate's two fails (`TestMissionConfigDefaultFlowIntegration`, `TestProtectedGitCandidateEquivalence`) are unreachable from S9: `go list -deps ./host/verifygate` has none of the touched packages. Whole broker with `-timeout 20m` at the stage: 3 fails outside the base 6, each PASS 2/2 alone at head and at base.
- **M2 `c9cabdfd`** on `ACCEPT_OK M2 outside` (the 445–447 precedent), NUMSTAT 1,399.

**Controller sweep** on `c9cabdfd` (outside the sandbox): gofmt, build, vet, linux vet, `go test -run '^$' ./...`, every guard inventory (store, CLI writer sites/session roots/production deadlines, boundary, daemon, `TestSubprocessSinkGate`, the broker subprocess-site guard), `CALLER_CENSUS_OK`, `verify_ail.sh`, checkrun/worktree/store race (store 44.6 s vs base 43 s) — all rc 0.

**Judge** (Agent `sonnet`, ~66 min): **PASS 74** (25/18/16/6/9), 0 blocking.
- Reachability sound: no production caller of the private entry or the local-commit primitive; push/PR/label/V3 put unchanged. F1 now precedes flag parsing (exact `--help`/`-h` exempt): a bad flag with no tty now gets the F1 refusal — stricter; no existing assertion weakened.
- 66 own mutants + 9 planner ids; all 8 re-run planner mutants died; no static-only kill.
- Survivors / notes (non-blocking): ~290 M2 lines have no test or caller until S9b; the install pinning-skip scope (skipping `git` or everything when `finish_remotes` is set survives); the private-entry `Kind=="finish"` and `ConfirmOnSurface` checks (controller-reproduced: removing the kind check, `cmp`-confirmed, leaves `TestOperatorFinishNoProof|F1|TTYParity` ok); digest content, walk/name/byte bounds, `config.worktree`/alternates fence; commit branch check, post-commit recheck, update-ref CAS, author line.
- O1: the new test passed 25/25 under 32 CPU hogs; the EXISTING `TestOperatorServiceStartup` failed 8/25 under the same load with the identical message — a 1 s real-time heartbeat-write budget in `runnerTick` shared by every `start` test. O3: store race base 109.2 s vs head 40.3 s — contention, no S9 delta.

**Gate 3b:** PR #707 on `c9cabdfd`: CI 3/3 at attempt 1 → MERGEABLE/CLEAN, `origin/dev` unmoved → merged `--match-head-commit` → **`c1869819`**. Merge CI: 3/3 at attempt 1 (run 46129639435) → **LANDED**.

**Found:**
- **The iteration-447 process fix worked as designed on first use:** the outside legs ran on the planner's stage before any executor round; the only red was attributed before the executor started, at no round cost.
- **The controller's own sweeps contended with its commit gate** (instance 3 after 446/447), and the plan's `commit-guards.sh` lacks `-timeout` for a package that needs 12 minutes here. Process fix below.
- **A shared 1 s real-time heartbeat budget in `runnerTick`** makes every `start`-path service test load-sensitive (judge: `TestOperatorServiceStartup` 8/25 under load). Pre-existing; the S9 test inherits it. Candidate row under 258 if it reaches CI.
- `TestA2AInvokeDeadlineClosesSocket` did not recur this iteration.

**Progress:** 1.0: clauses 1–7 MET; goal unmoved. Row 243 V4 now has its first finish slice: F1 is pinned for every public finish, the runner can tell an operator at install that finish is blocked on GitHub auth while start still works, and the evidence primitives a no-terminal local commit needs exist. S9b, compiled, activates them.

**Ruled out:**
- Attributing the M2 commit-gate refusals to S9: attempt 1 was the controller's own contention (passes on a quiet rig); attempt 2 was a missing `-timeout`, verifygate cannot reach S9, and the three broker fails pass alone at head and at base.
- A planner round 2 for `TestOperatorFinishAuthInstall`: not reproduced in 5 alone runs and 2 leg runs; the judge attributed it to a pre-existing shared budget.
- Landing M3 (S9b) in this sprint: it would take the sprint to 2,799 lines, and the plan lists open activation gates.

**Next:** row 243 **V4 S9b** — apply the plan's compiled Appendix P M3 (re-measure on `c1869819`): prep admission and claim, pinned sandboxed checks with reservations, commit reservation and reconcile, redacted status and the taskless prep review. Carry:
- judge-448 survivors as named pins (private-entry `Kind=="finish"`/`ConfirmOnSurface`, install pinning-skip scope, digest content and bounds, commit branch/CAS/author);
- give `commit-guards.sh` and every whole-package plan leg an explicit `-timeout` (process fix);
- D-WORLD-122 before S11.

## 449 — 2026-11-28 — row 243 V4 sprint 9b LANDED (PR #709, merge `d9635675`), outcome B: the finish proof scope, the install tool pins and the local Git boundary are pinned, and the private check and commit journal runs end to end in fixtures to a durable `local-commit-done` receipt — still with no production caller; S9c (1,243 lines) compiled in the plan; 1,218 lines; judge sonnet PASS 78 [PRODUCT]

**Pick:** row 243 V4 S9b, the record-448 Next.
- **Clause map:** 1.0 clauses 1–7 MET, none UNMET.
- **Ledger:** valid, 109 rows, 1 OPEN (D-WORLD-122, Mark's Sol preset for E5 — an S11 input; default B, fakes).
- **Design:** V4 r1 reviewed (quorum reader rc 0, r2 `…T10-24-51Z` PROCEED 3/3, `absent_reviewers` `[]`).
- **Verified at `77dca5ed`:** `cmd host docs app` byte-identical to `c1869819`; the iteration-448 `M3.patch` dry-run applies rc 0 with no offset/fuzz/reject; `frozen.Kind != "finish"` at `operator_finish_preflight.go:23`. No S9b plan, branch, PR or worktree existed.

**Gate 0:** stamped first; `gh` = `sunholo-voight-kampff`; billing CLEAN; kill switch absent.
- 0 allowlisted directives on #635 (36 comments; LEGACY watermark).
- Step 6a rc 0: 0 crash notices (control 107:4). 4 transcript-scan notices in the window: the three iterations 446–448 reported, plus one new for iteration 448's fire `world-1795842267-52962` (3 hits; key names in the private controlplane notice only).
- `mission-world` inbox: nothing unread. Weekly sweep not due; no rotation (next Monday 2026-11-30).

**Gate 1:** local `dev` `9a40fe68` 2 behind `origin/dev` `77dca5ed` (448's record); read from origin. Skill drift: the symlinked copy differed on `resources/gate-5-retro.md` only (as in 448) — Gate 5 followed origin. Dev CI 3/3 green. No open own-account PRs.

**Routing evidence:**
- Resolver followed verbatim: planner `recipe codex:gpt-6.2-sol declared:planner-lane-default-pin` against the V4 doc (a doc-less resolver call first printed `agent-tool opus fail-closed:no-doc`; re-run with the doc, and `derive-planner-lane.sh` agreed — no opus spawn attempted); executor `recipe codex:gpt-6.2-sol declared:provider-pin`; evaluator `agent-tool sonnet declared:alias-pin`; no designer. Codex probe rc 0; codex 17.0% used / 25.9% allowed.
- The standing request asked for every role through the Agent tool "exactly as the routing table specifies". The table returned `recipe` for the planner and executor, so they ran as `codex exec` recipes; the REQUIRED judge ran on the Agent tool (`sonnet`, foreground). No role failed to spawn.
- Token use: controller `claude-opus-5-5` (session; tok: not reported) · planner `codex:gpt-6.2-sol` 538,978 · executor `codex:gpt-6.2-sol` 118,719 · evaluator Agent `sonnet` 105,809.
- Judge independence: cross-vendor. Metered: $0.00.
- base=d9635675ec6ea746542a435ad658baf212cbd346@2026-11-28T13:22:40Z (Gate 4 record). Gate 1 and Gate 3 snap `77dca5edc8206a1791743ebb3c64512f051c52f3` agree. Sprint, base, judge and record worktrees via `mission-worktree.sh`.

**Baseline** (controller, pristine `77dca5ed`, outside the sandbox, isolated `AILANG_STATE_DIR`, while the planner ran): store/authority/agentrun plain+race and verbs 0 fails; verifygate 1; broker 4; cmd 16; daemon-plain 8; daemon-race 10 (974 s).

**Planner** (~67 min, rc 0). Plan + JSON `3d8c0bd5`, **outcome B**: rebased whole 2,461 lines; S9b = M1 401 → M2 1,218; M3 = S9c activation, 1,243 incremental, compiled in Appendix P with five activation gates (real-listener admission and route-to-worker, the integrated eligibility/drift matrix, two-way terminal/runner exclusion, privacy absence, whole-package gates and ledgers).
- Executed its own §7 block verbatim (`APPLY_ACCEPT_SNAPSHOT_OK` M1/M2/M3, `FINAL_APPLY_RC 0`).
- Corrected five inherited premises: a phantom path counted against the 257-name cap; an unverified author header in the commit primitive; a loose content-addressed object accepted as check evidence without a matching World log entry; invalid SQL bindings in judge-448's `e-store-commit-once`; an empty workspace in the inherited boot catalogue (S9c).
- Controller round trip on a fresh `git archive 77dca5ed`: `diff -r` rc 0 for cmd, host and docs at both stages. Destructive-write scan of the plan scripts: clean.

**Outside legs on the planner's stage, before the executor (process fix 447, second use):** M1 `outside`, `outside-race` (`added=3.84s cap=15s`) and `outside-1777` all `ACCEPT_OK`. M2 `outside` and `outside-1777` `ACCEPT_OK`; `outside-race` every assertion PASS but `RACE_BUDGET added=111.61s cap=15s` (load ~5, `corespotlightd` ~140% CPU, a `godot` job). Re-run alone 3/3: 6.52 s, 6.43 s, 6.71 s, each `ACCEPT_OK M2 outside-race`. Handed to the judge as O1; the executor started on that basis.

**Executor** (~5 min): the §7 block verbatim with `TOOLS=$STATE/exec-tools`, M1 and M2 only. `ACCEPT_OK` ×2; snapshots 5/0/5 and 10/0/10, `cmp`-equal to the planner's stages; all 17 tools `cmp`-equal.

**Commit gates (quiet rig — the 448 process fix held; no controller sweep beside either gate):**
- **M1 `5717be11`** on `ACCEPT_OK M1 outside` (401), first attempt.
- **M2 `f11d09c7`** on `ACCEPT_OK M2 outside` + `COMMIT_GUARDS_OK` (whole verifygate + broker, `-timeout 20m`, no new fail vs base), first attempt, NUMSTAT 1,218.

**Controller sweep** on `f11d09c7` (outside the sandbox, finished before the judge ran anything): gofmt, build, vet, linux vet, `go test -run '^$' ./...`, every guard inventory, `TestSubprocessSinkGate`, the broker subprocess-site guard, `CALLER_CENSUS_OK` (local-commit 0; push/label/view-PR/create-PR/V3-put 1 each), `verify_ail.sh`, checkrun/worktree/store race (store 42 s) — all rc 0.

**Judge** (Agent `sonnet`, ~55 min): **PASS 78** (25/21/17/7/8), 0 blocking.
- Reachability sound: no production caller of the private worker, private entry or local-commit seam; claim/push/PR/label/V3-put unchanged; nothing new exported. A signed start and a finish without confirm-on-surface are refused with a no-change read-back (j01/j02 killed); the ID-vs-CAS check is redundant (j03 survived), as judge-448 said.
- 35 own mutants; all 19 judge-448 survivor edits the plan pinned are now killed, including `e-store-commit-once`. Survivors (non-blocking): j14 resume from a reserved-but-unobserved commit (controller-reproduced: replacing `return w.reconcileCommit(j.CommitKey)` at `operator_finish_local.go:101` with an error leaves every `TestOperatorFinish*`/`TestOperatorCommit*` ok); j08 earlier check-receipt World-evidence recheck on resume; j13 `commit-failed` recorded without an input recheck first; j16 reconcile config re-read; j19 evidence SemanticID; j12, j17, j25.
- O1: not reproduced on a quiet rig; the new tests inject a 3-minute budget, a clock and a fake check runner, so load slows them but does not change their verdicts. Whole `cmd` race delta and scratch hostile-repo cases UNMEASURED by the judge.
- Landing: safe; the second consecutive sprint shipping production code with no production caller — accepted as dead weight owned by S9c.

**Gate 3b:** PR #709 on `f11d09c7`: CI 3/3 at attempt 1 (run 46148073561) → MERGEABLE/CLEAN, `origin/dev` unmoved → merged `--match-head-commit` → **`d9635675`**. Merge CI: 3/3 at attempt 1 (run 46149287254; race `cmd/ailang-worldd` 292.3 s, `host/daemon` 366.1 s, all under the 900 s bound) → **LANDED**.

**Found:**
- **A single per-test race-budget reading is load-dominated on this rig** (instance 2 after 448's planner-reported 98.11 s contended M3 timing; this iteration's planner also saw 24/176/142 s): 111.61 s vs 6.4–6.7 s quiet for the same bytes, assertions green every time. The plan's `RACE_BUDGET` gate fails the leg on one reading, which reads as a red that is not a code defect. Process fix below.
- **Two consecutive slices (S9, S9b) shipped production code with no production caller.** Each is tested in fixtures and owned by the next slice, and the judge accepted it twice; S9c must activate it, not split again behind another uncalled layer.
- The iteration-448 process fix (no controller sweep beside a commit gate; explicit `-timeout` on whole-package legs) worked on first use: both commit gates passed first time.

**Progress:** 1.0: clauses 1–7 MET; goal unmoved. Row 243 V4's no-terminal finish now has its whole local worker — proof scope, pinned tools, bounded capture, sandboxed checks with reservations and an exact, reconciled local commit — exercised end to end in fixtures. S9c, compiled, puts it behind the bar and the card.

**Ruled out:**
- Attributing the M2 `outside-race` overshoot to S9b: 3/3 quiet re-runs at 6.4–6.7 s on the same stage bytes, and the judge found no wall-time dependence in the new tests.
- A planner round 2 for the race overshoot: not a plan defect.
- Landing M3 (S9c) in this sprint: 2,461 lines cumulative, and its activation gates are open.

**Next:** row 243 **V4 S9c** — apply the plan's compiled Appendix P M3 (re-measure on `d9635675`): prep admission over the real listener, runtime readiness and eligibility before claim, the claim and the worker caller, the taskless review, `checks`/`localCommit`/`finish-status`, with its five activation gates. Carry:
- judge-449 survivors as named pins (j14 reserved-commit resume, j08 earlier-receipt evidence recheck, j13 input recheck before `commit-failed`, j16, j19);
- the race-budget process fix below;
- S9c must ACTIVATE the S9/S9b worker; if it cannot fit, cut so the caller and every gate it needs land together rather than another uncalled layer;
- D-WORLD-122 before S11.

## 450 — 2026-11-28 — row 243 V4 sprint 9c LANDED (PR #716, merge `9b11c617`), outcome B: THE NO-TERMINAL LOCAL FINISH IS SWITCHED ON — bar review → Request approval → signed card → enrolled key → runner proof, Git readiness, flock and eligibility → claim → the S9/S9b check-and-commit journal → typed `local-commit-done` receipt, no remote effect; 1,399 lines; judge sonnet PASS 84 → re-ruled 77 non-blocking on 41 controller-measured named survivors, S10 gated on a test-only S9d [PRODUCT]

**Pick:** row 243 V4 S9c, the record-449 Next.
- **Clause map:** 1.0 clauses 1–7 MET, none UNMET.
- **Ledger:** valid, 109 rows at pick, 1 OPEN (D-WORLD-122, default B, fakes). After the mid-iteration fetch: 113 rows, 2 OPEN (D-WORLD-122; D-WORLD-129, row 274, filed by an attended session).
- **Verified at `91befdab`:** `cmd host docs app` byte-identical to `d9635675`; the iteration-449 `M3.patch` dry-run rc 0, no offset/fuzz/reject.

**Gate 0:** stamped first; `gh` = `sunholo-voight-kampff`; billing CLEAN; kill switch absent.
- 0 allowlisted directives on #635 (37 comments; LEGACY watermark).
- Step 6a rc 0: 0 crash notices (control 107:4). 4 transcript-scan notices in the window, the same four fires iterations 446–449 reported (`world-1795730668-41608`, `world-1795769251-75944`, `world-1795786215-30397`, `world-1795842267-52962`); none new. Key names are in the private controlplane notice only.
- No rotation (#635 created 2026-11-23; next Monday 2026-11-30).

**Gate 1:** local `dev` `9a40fe68` 2 behind `origin/dev` `91befdab`; read from origin. Skill drift: the symlinked copy differed on `resources/gate-5-retro.md` only (as in 448/449); followed origin. Dev CI 3/3 green. No open own-account PRs.

**Routing evidence:**
- Resolver followed verbatim against the V4 doc: planner `recipe codex:gpt-6.2-sol declared:planner-lane-default-pin` (`derive-planner-lane.sh` agreed); executor `recipe codex:gpt-6.2-sol declared:provider-pin`; evaluator `agent-tool sonnet declared:alias-pin`; no designer. Codex probe rc 0.
- The standing request asked for every role through the Agent tool "exactly as the routing table specifies". The table returned `recipe` for the planner and executor, so they ran as `codex exec` recipes. The REQUIRED judge ran on the Agent tool (`sonnet`, foreground), and was resumed twice by `SendMessage` with new controller evidence. No role failed to spawn.
- Token use: controller `claude-opus-5-5` (session; tok: not reported) · planner `codex:gpt-6.2-sol` 494,248 · executor `codex:gpt-6.2-sol` 57,489 · evaluator Agent `sonnet` 97,258 (r1), re-rulings as reported 106,952 and 111,761.
- Judge independence: cross-vendor. Metered: $0.00.
- The size cap was 1,400 with a directive-granted extension to ≤ 1,800 only if no activated cut fit; it was not needed.
- Worktrees via `mission-worktree.sh` (sprint, base, judge `eval/world-iter450-r1`, record).

**Baseline** (pristine `91befdab`, outside the sandbox, isolated `AILANG_STATE_DIR`, during the planner): store/authority/agentrun plain+race, verbs, broker, surfacepage plain+race, checkrun/worktree race 0; verifygate 1; cmd-plain 11 (the row-271 e2e set); daemon-plain 3; daemon-race 5.

**Planner** (~50 min, rc 0). **Outcome B:**
- S9c = M1, the activation, 1,399 lines.
- S9d = M2, the read-only `service runner finish-status ACTION-ID` mirror, 57 lines, compiled but not applied.
- §7 executed verbatim (`APPLY_ACCEPT_SNAPSHOT_OK M1/M2`, `FINAL_APPLY_RC 0`); controller round trip `diff -r` rc 0 for cmd, host and docs at both stages; destructive-write scan of the plan scripts clean.
- It corrected three inherited-draft defects: already-committed paths missing from attribution, the reserved restart bypassing the worker, and check-evidence locators exposed in status.
- It reported zero named survivors, with its listener rows "compiled, verdicts pending".

**Outside legs on the planner's stage before the executor (447 fix, third use):**
- Focused: `ACCEPT_OK M1 outside` / `outside-race` / `outside-1777`; `RACE_BUDGET cmd min=12.47s cap=15s`, `daemon 0.60s cap=12s`.
- Chrome: `TestOperatorStartWalkthroughInChrome` PASS.

**Executor** (~8 min): `FINAL_APPLY_RC 0`; snapshot 21/0/21, manifest verified, 21/21 `cmp`-equal to the stage; tools 22/22 `cmp`-equal.

**Commit gate (quiet rig, no controller sweep beside it):** **M1 `e69298f4`** on `ACCEPT_OK M1 outside` + `COMMIT_GUARDS_OK`, first attempt, NUMSTAT 1,399.

**Controller sweep** on `e69298f4`, sequential:
- The pristine skip-11 cmd race base-extra: 427 s.
- The planner's whole `gates.sh … M1 outside` (2,824 s): gofmt, build, vet, linux vet, test-compile, every guard inventory, subprocess sink/sites, js-syntax, verify_ail, the caller census, and accept outside/race/1777 all OK.
- Whole cmd, daemon, verifygate, broker and surfacepage, plain and race: all 8 compares OK against base fail sets, `GATE_RESULT issue=0 pending=0`.
- Whole cmd race 415.5 s pristine vs 425.7 s head (two distinct runs, checked: different files and loads).
- One controller error, corrected: the first sweep also queued the named mutants beside the judge's window. It was stopped one minute in and restarted without them, and its orphaned base-extra child was awaited.

**Judge** (Agent `sonnet`, ~70 min): **PASS 84** (26/24/17/9/8), 0 blocking.
- **Reachability is exactly the runner path:** confirmed rows → proof judge → Git auth → non-blocking flock → `finishInputs` → claim → worker. Claim 1 → 2, private entry 0 → 1, local commit 0 → 1. No new push/PR/label/V3-put caller; the route test sees `verb_requests`=0, no finish record and an unchanged bare remote. The public CLI was checked by grep only (UNMEASURED by fuzzing).
- **QUICKSTART:** `TestOperatorFinishAdmission` drives the printed command verbatim.
- **Judge-449 pins:** j08, j13, j14, j16 and j19 re-applied and killed by their named tests.
- **Its own mutants:** 14 of 21 survived (e01, e04, e06–e09, e12–e14, e16, e17, e20–e24).
- **Race:** cmd ≈ 8.8 s and daemon 0.54 s, single readings at load about 6.
- No existing test was weakened.

**Named mutants, first outside run** (controller, after the judge's verdict, while PR CI ran remotely; overlay-based, tree untouched): 124 ids → **80 killed, 41 SURVIVED, 1 live control, 2 S9d skips**.
- The plan's `named-mutants.sh` first crashed (`args[@]: unbound variable`, bash 3.2), and was re-run as its exact `mutate.py` call.
- Survivors:
  - 16 `runner-elig-*`
  - 6 `runner-drift-*`
  - `runner-effect-{clock,inputs,key}`
  - `runner-observe-after-expiry`, `runner-old-observation`
  - `runner-manifest-protected`, `runner-catalogue-{fresh,approve-size,bound}`, `runner-heartbeat-{owner,ready}`, `runner-life-{quiet,pair}`
  - `route-finish-store`, `route-protected-and-unattributed`, `route-operation-generation-conflict`, `route-confirm-without-surface`
  - `status-definite-failure`, `status-finish-auth`
- The controller reproduced judge e12 first-party: removing `|| len(m.Unattributed) != 0` at `operator_finish_runner.go:278` leaves every `TestOperatorFinish*`/`TestOperatorCommit*` ok (21.0 s); the tree was restored byte-identical. The worker guard at `operator_finish_local.go:93` and the route guard at `surface_operator.go:56` remain.

**Judge re-rulings** (resumed with that evidence; code trace, execution UNMEASURED):
- **RE-RULING 1: PASS 78, non-blocking.** Each traced survivor drops a re-check, not the only guard: `prep.claimed()`, `boundManifest`, the signed-snapshot compare in `w.inputs()` before every check and the commit, the `git.commit` config check and the store CAS still refuse an unapproved commit. What is lost is post-claim drift and expiry detection inside one claimed local run. The plan's gates 2 and 5 were not met as written.
- **RE-RULING 2: PASS 77, non-blocking.** `route-operation-generation-conflict` returns the first card instead of a 409, with no second action. `route-finish-store` is still refused at confirm and in the runner. `status-definite-failure` shows a misleading message on a failed local commit, but `localCommit` stays empty.
- **Hard precondition:** S10 (the first remote effect) must not start until a test-only follow-up kills these survivors.

**Gate 3b:**
- PR #716 on `e69298f4` got 0 runs; `mergeable` → **CONFLICTING/DIRTY**. Cause: attended row-274 merges #711–#713 to `dev` (`91befdab..4bac6fe1`) also appended to `docs/QUICKSTART.md` §21. Read before any dropped-event lever; none was fired.
- Rebased: resolved as base + the S9c `###` subsection + the attended `## 22`. Product diff byte-identical to `e69298f4` outside QUICKSTART, QUICKSTART added lines identical, NUMSTAT 1,399.
- Re-checked the rebased head `f45da2a9` outside the sandbox: build, vet, linux vet, test-compile, js-syntax, `ACCEPT_OK M1 outside`, QUICKSTART/verb/parity/app tests in cmd and daemon, verbs/appmanifest/surfacepage, verify_ail — all rc 0.
- Force-pushed with lease; the stale poller on `e69298f4` was killed.
- PR CI 3/3 on `f45da2a9` (event `pull_request`) → `MERGEABLE CLEAN`, `origin/dev` unmoved → merged `--match-head-commit` → **`9b11c617`**. Merge CI: merge CI 3/3 at attempt 1 (run 46169801018; race `cmd/ailang-worldd` 528.9 s, `host/daemon` 532.0 s, `host/verifygate` 534.5 s, `host/broker` 422.3 s, all under the 900 s bound — broker and verifygate, untouched by S9c, rose by similar amounts, i.e. runner variance). → **LANDED**.

**Found:**
- **Named-mutant rows a sandboxed planner cannot run were unmeasured, and nothing ran them before the verdict.** The 447 fix moved acceptance legs to the planner's stage but not mutant rows, so a "zero survivors" ledger with 0 of 74 outside rows executed reached the executor, the commit gate and the judge. The first outside run found 41, and the plan's own gate 5 was unmet at landing. Process fix below.
- The plan's `named-mutants.sh` crashes on macOS bash 3.2 (empty-array expansion under `set -u`) — an instrument failure that reads as "nothing ran".
- Attended sessions merged to `dev` mid-iteration (rows 274/275, new D-WORLD-123..129). The record and rebase absorbed them; the `mergeable`-first rule caught the conflict at the first missing-runs reading.

**Progress:** 1.0: clauses 1–7 MET; goal unmoved. Row 243 V4 now has its first operator-visible no-terminal finish: a reviewed, key-approved local commit made by the runner, with checks and the commit SHA in status. Pushing, the PR and the finish record remain S10/S11.

**Ruled out:**
- Holding the merge for a round-2 test milestone. The judge, given every survivor, ruled it non-blocking with a precondition on S10; the effect is local, opt-in and fake-preset only.
- A dropped `pull_request` event (`mergeable` read first: CONFLICTING).
- Treating the identical 427 s cmd race readings as one run (415.5 s vs 425.7 s, distinct files).
- The 1,800 cap extension (M1 fit at 1,399).

**Next:** row 243 **V4 S9d** — a TEST-ONLY milestone that kills the 41 named survivors and judge-450 e12. Use a table-driven route-worker test with one perturbation before the claim (expect no claim and zero checks or outcomes) and one between the claim and the first check (expect refusal before the check, which needs a post-claim hook). Add the compiled Appendix P M2 `finish-status` mirror. Re-run the named set outside on the planner's stage under the process fix below. **S10 must not start until S9d lands.** D-WORLD-122 before S11.
