# Iteration 443 — controller readings (row 267, race-job timeout headroom)

Fire `world-1795751901-5784`, attempt 1. Controller `claude-opus-5-5` (session). Base `origin/dev` `6420426e` (= S5f4 merge `c3504c56` + record #690).

## Gates 0–2
- Gate 0: kill switch absent (`mission-world.disabled`); `gh` = `sunholo-voight-kampff`; billing CLEAN; `mission_directives.sh` on #635: 0 allowlisted directives (25 comments; read against both watermarks, the older `2026-08-17T06:15:12Z` governing); step 6a: rc 0, 0 crash notices on #635/#574 (control 107:4), 1 transcript-scan notice (fire `world-1795730668-41608`, 2 key/file hits — key names only in the private controlplane notice); `mission-world` inbox 0 unread inbound.
- Gate 1: local `dev` was 26 behind `origin/dev`; fast-forwarded (untracked files only, no conflict). Base `6420426e…@2026-11-27T04:01:52Z`. Skill drift: resolved symlink and pin copy byte-identical to ailang `origin/dev` on all 14 files. **Dev CI at `6420426e`: race job FAILED at attempt 1** (run 45936239736: `FAIL cmd/ailang-worldd 480.085s`, `FAIL host/daemon 480.072s`, each `panic: test timed out after 8m0s`; `host/verifygate` 424.143 s, `host/broker` 333.040 s, `host/store` 97.766 s). Re-run requested (`gh run rerun --failed`) → attempt 2 green.
- Gate 2: clause map 1.0 clauses 1–7 MET (no UNMET clause). Ledger 1 OPEN: D-WORLD-121; its default B says row 267 runs when Gate 1 finds `dev`'s tip red — it did, so the pick is row 267 (the standing red-dev rule; this mission owns the repo). No row-267 design, plan, branch, PR or worktree existed.

## Evidence the design used (controller)
- `race_table.txt` (from `race_table.py` over `racelogs/`): one row per race-job attempt for the last 40 CI workflow runs (53 rows = 51 complete attempts + 1 in progress + 1 duplicate). 18 complete attempts hit the 480 s cap (`F480.xxx`), every one on a slow runner (`host/store` ≥ ~81 s against ~58–68 s normal).

## Routing
- Resolver, followed verbatim: designer `recipe claude:claude-opus-5-5 declared:provider-pin`, planner `recipe codex:gpt-6.2-sol declared:planner-lane-default-pin`, executor `recipe codex:gpt-6.2-sol declared:provider-pin`, evaluator `agent-tool sonnet declared:alias-pin`. Probes: claude-sub rc 0, codex rc 0.
- The operator's standing request asked for every role through the Agent tool "exactly as the routing table specifies". The table returned `recipe` for the designer, planner and executor, which makes the Agent tool not a valid path for them; they ran as `claude-sub -p` / `codex exec` recipes. The REQUIRED judge ran on the Agent tool (`sonnet`, cross-vendor vs the codex executor). No role failed to spawn.

## Designer
- Launch 1 failed instantly (rc 1): `claude -p … --add-dir "$S" "$PROMPT"` — `--add-dir` is variadic and swallowed the prompt. Controller instrument error; fixed with `--`; relaunched.
- r0 (opus, 32 turns, ~12 min) → `8256a83e`: option A (fixed per-package 900 s, leg 1800, step 40) + ordering guard + loud wrapper.
- Quorum r1 (author benched; gemini-4-carbon, oc-glm-5-3, oc-kimi-k3; no absent seat): **BLOCKED 3/3** ($0.34 metered). Controller measured C1–C4 (readers of verify_go.sh; plain leg has no `-timeout`; `ps` did not truncate a 146-char argv under a pipe on darwin; start allowance) → r1 (resumed session, 12 turns) `4c76f7c0`: A 813, guard L ≥ P + 60 + A, L 1900, `ps -ww`, bounded R6 poll.
- Quorum r2 (oc-glm-5-3, oc-kimi-k3, codex:gpt-6.2-astra; reserve gemini-4-carbon — FLAG: not the same three seats as round 1): glm PASS, kimi REJECT (censored daemon slow tail → A 995, L 2000), astra REJECT (unbounded `ps` before the kill). Both with concrete fixes, neither disputing the direction → controller r2 `c50264d6` under the narrow-refinement carve-out (verbatim fixes; §0 + §11 of the design).

## Planner
- `codex exec`, ~60 min, **398,347 tok**, rc 0. Plan + JSON `915ba9c0`. One milestone M1, 351 changed lines (`-- scripts host .github docs`); 24 named kills, 0 observed survivors; R6 UNINFORMATIVE UNDER SANDBOX (real `ps` denied). Daemon per-k across the 9 uncapped mid-class attempts falls with k (Pearson −0.583), so the r2 values are conservative.
- Controller round trip `ctl_roundtrip.sh M1`: `diff -r` rc 0 on `scripts`, `host`, `.github`, `docs`.

## Baseline (controller, pristine `6420426e`, outside the sandbox, co-loaded with the planner; load ~25–37)
- `verifygate-plain` 1 fail (`TestMissionConfigDefaultFlowIntegration`, 10 s fixture deadline, load), 513 s; `pins-plain` 0; `all-plain` (`go test ./... -count=1 -timeout 45m`) rc 1, 2,626 s, **49 fails** (cmd/ailang-worldd 20, host/daemon 19, host/broker 7, host/coordinator 1, cmd/world-publish 1, host/verifygate 1).

## Executor (round 1)
- `codex exec` (**55,613 tok**), rc 0. Patch applied clean; `accept.sh` stopped at R6 — UNINFORMATIVE UNDER SANDBOX (`could not list race leg processes: [Errno 1] Operation not permitted: 'ps'`); snapshot 6/0/8, 6/6 cmp-equal to `st-M1`.
- `commit_ms.sh M1` → **`234ac3f6`**, NUMSTAT 351; commit-time test outside the sandbox: six named pins PASS, `PASS 7/7 arms` (R6 2 s, R7 4 s), whole `host/verifygate` ok 556.693 s.

## Controller gates on `234ac3f6` (outside the sandbox)
- `accept.sh`: `ACCEPT_OK` (AC1–AC5 incl. both AC3 environments and the named mutant replays).
- `gates.sh`: gofmt, build, vet, Linux vet, both `bash -n`, `verify_ail` rc 0; `verifygate-plain` and `pins-plain` **no new failures**; `all-plain` 37 fails, 11 "new" vs base — all in packages M1 does not touch (`cmd/ailang-worldd`, `host/broker`, `host/daemon`, `scripts/marketingcheck`; `git diff --stat 6420426e 234ac3f6` over them is empty; no file there reads `verify_go.sh`/`ci.yml`/the arms script — the two broker hits are a protected-path string and a comment) while 23 base fails did not recur. Re-run alone at head: 10/11 pass; `TestOperatorFilesReconcile` fails alone too (git-fixture `ep.task.md` missing; on the charter's load-sensitive list; byte-identical package).
- PR #691 CI on `234ac3f6`: 3/3 at attempt 1 (race job 12m16s, build+test 13m16s).

## Judge round 1
- Agent `sonnet`, foreground, ~20 min, **96,298 tok**: **FAIL 74** (30/14/18/12). **B1:** R6 load-flaky — 3/14 runs with a concurrent `go test ./host/verifygate/` (`grandchild PID control missing (rc=124)`); R7 near its bound. Controller reproduced: 2/10 with 12 `yes` spinners (`expected rc 124 and timeout diagnostic (rc=1)` = the arm's own 24 s watchdog firing). Non-blocking survivors G5 (SIGKILL→SIGTERM), G14 (`|| true` on the arms step), G21 (doubled wait bound); X07/X08 ledger rows mis-attributed. Numbers re-derived from the raw logs; ordering and margins hold.
- **Side effect (disclosed by the judge):** it SIGKILLed PIDs 96794 and 96825 — the divergence watchdog and the rig-lock heartbeat subshells of an unrelated `ailang/tools/finetune/pilot5_run.sh` run (parent 340, training 96824) — believing them its own leaks. Training kept running (`train.log` advancing). Controller restored both loops verbatim from `pilot5_run.sh:47–50` and `:176` (`restore_pilot5_loops.sh`, detached, new PIDs 47384/47385), so the GPU rig lock (`/Users/Shared/ailang/rig.lock.d`, 6 h steal threshold) keeps its heartbeat.

## Round 2
- Codex (**85,755 tok**) in a separate worktree (`.wt-world-iter443-row267-r2`, so the controller's running gates were not disturbed): 39 lines over the arms script + the product-sequence pin — readiness file, 10 s bound, 50 s watchdogs, relative wall window, `trap '' TERM`, whole-line arms-step needle.
- Controller outside the sandbox: gates rc 0; G5, G14, G21 KILLED; **contention 2/20 FAILED** (`probe never signaled readiness` at 10 s and 13 s; the listing shows the probe shell alive with no `sleep` child). Discriminating measurement: exec of a freshly written script median 0.191 s / max 0.244 s vs inline `/bin/sh -c` 0.006 s / 0.023 s (12 spinners, 25 samples). → follow-up: run the probe inline.
- Follow-up 2b (codex, **76,961 tok**): inline probe. Controller: G5/G14/G21 killed; contention 1/40 (`probe never signaled readiness`, R6 at 10 s). Fresh judge round 2 (Agent `sonnet`, **58,002 tok**): **FAIL 77** (31/17/17/12): B1 not closed, 2/32 (`probe never signaled readiness`); G5, G14, G21 CLOSED; own mutants N1–N4, N8–N10 killed, N5/N6/N7/N11 survive (test-side). Its suggested fix: retry the arm on never-ready.
- Follow-up 2c (codex, **80,169 tok**): `run_hang` retries ONLY an attempt that ended before `ready` (≤ 3, each logged); every product assertion runs on a ready attempt and fails at once; named `wall_upper_bound = bound + 2 + 8`; `controller.sh notready`. Controller outside the sandbox: gates rc 0; G5/G14/G21 KILLED; `notready` FAILS after 3 logged retries; **contention 0/40, 0 retries**. Round 2 total: 72 changed lines over the arms script and the product-sequence pin.
- Committed `b2b44792` (named pins PASS, `PASS 7/7 arms`); sprint total **379** changed lines (`-- scripts host .github docs`).

## Judge round 3
- The round-2 instance resumed by name (`SendMessage`), **73,942 tok** cumulative for rounds 2–3: **PASS 79** (31/18/18/12; its report's score line is garbled — "= 85? NO" — and the corrected total is 79), 0 blocking. Contention 1/37 loaded runs (12 spinners 1/22, concurrent `go test` loop 0/12, quiet 0/3), 0 retries needed; the one failure is the outer wall window (`wrapper elapsed 26s, expected 10s to under 20s` at load ~31), not the readiness race. Retry cannot mask a defect: G5, G21, N3, N8 fail on a ready attempt with no retry; `notready` fails after 3; a wrong retry condition (retry on rc ≠ 124) plus G5 or plus `exit 124 → 0` still fails. Survivor: retry attempts 3 → 1 (a tolerance parameter). Non-blocking residual carried: compare the wrapper's own reported elapsed (`00:NN` in the listing) instead of the arm's outer wall clock.

## Gate 3b
- PR #691 on `b2b44792` 3/3 at attempt 1 (race `cmd/ailang-worldd` 284.279 s, `host/daemon` 359.767 s; `PASS 7/7 arms` with Linux procps — AC7); merged `--match-head-commit` → `6447e776`; merge CI 3/3 at attempt 1 (run 45963259258) on a slow-class runner (`host/store` 86.649 s): `cmd/ailang-worldd` 422.919 s, `host/daemon` 452.341 s, `host/verifygate` 393.993 s — AC8 observation 1/10, green.

## Post-judge controller check
- `ps` after the round-1 judge: the pilot5 parent (340) and training (96824) alive; the two restored loops alive (PIDs 47384/47385, reparented to launchd); `train.log` advancing.
