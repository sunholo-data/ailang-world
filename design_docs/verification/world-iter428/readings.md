# Iteration 428 — readings (row 243 V4 sprint 3)

All times UTC, 2026-11-24. Base (Gate 1) `92e1c44ebd8dbbdb2b962ef54aed0dd9a0bf2456` @ 11:17:30Z.

## Gate 0–2
- Kill switch `mission-world.disabled` absent; `gh` = `sunholo-voight-kampff`; billing CLEAN.
- Directives: `mission_directives.sh --issue 635` → 0 allowlisted directives (5 comments, all public feedback); watermark file LEGACY (bare timestamp, unchanged by this fire — nothing to `--mark`).
- Step 6a: `mission_gate0_self_notices.sh` rc 0 — `issue: 635 comments=5 self=5 crash=0 transcript=0`, `issue: 574 comments=44 self=44 crash=0 transcript=0`, control `107:4 ok`; `-prev` watermark absent (DEGRADED single-watermark read, as before).
- Inbox: `ailang messages list --unread --inbox mission-world` → no messages.
- Attended rulings since 427: `82f05004` (Mark, 2026-11-24 09:20 +0100) flips **D-WORLD-119 → RESOLVED A** and **D-WORLD-120 → RESOLVED A**; `mission_ledger_provenance.sh` → `verdict=attended flip=82f05004… creation=933a4dab…` for both. Same commit audits rows 217/218/219/223/228 (attended status edits; no loop action).
- Driver: `inputs changed since the last no-pick fire (e4920be1… head=ea4a6f0 ledger=107 open=0 …) — proceeding with a full fire`.
- Dev CI at `92e1c44e`: `checks=3`, all success.
- Skill drift: 28 files (SKILL.md + 13 resources, resolved symlink target and driver pin) `cmp` against the driver clone's `origin/dev` `16d168f19` → 0 drift; negative control fired.
- Quorum reader `mc_quorum_gate2 w-world-address-bar-verbs-v4` rc 0 (`w-world-address-bar-verbs-v4-2026-11-18T10-13-56Z.json`); no re-quorum (r2 PROCEED 3/3, nothing changed).
- Died-mid-flight: no worktree or remote branch named iter428 before this fire.
- Collision surface at pick: open PR #641 (attended agy) and #617 (attended UI) — no file shared with the S3 plan's §5 list.
- Codex window: `state=ok used=41.0% allowed=56.3%` (driver, 11:15:58Z). Codex probe rc 0 (`ok`).

## Routing (resolver output, verbatim)
- planner: `recipe codex:gpt-6.2-sol declared:planner-lane-default-pin`; `derive-planner-lane.sh` → `codex:gpt-6.2-sol declared:planner-lane-default-pin`.
- executor: `recipe codex:gpt-6.2-sol declared:provider-pin`.
- evaluator: `agent-tool sonnet declared:alias-pin`.

## Planner (codex gpt-6.2-sol, workspace-write, 11:22–12:20Z, rc 0, 479,146 tok)
- Plan 241 / 738 / 1,396 (cap 1,400); S3b 132 lines carried. 138 mutants: 129 killed (named test + whole package), 9 named equivalents, 0 survivors, 0 unmeasured-in-sandbox.
- Controller round trip (`ctl_roundtrip.sh`, loop M1 M2 M3, fresh `git archive 92e1c44e`): every `diff -r` vs `st-M*` rc 0 (0-byte diffs). The two `.orig` files listed are tracked at the base under `design_docs/verification/world-iter362/`.

## Executor (codex gpt-6.2-sol, workspace-write, 12:21–12:24Z, rc 0, 80,255 tok)
- M1/M2/M3 `ACCEPT_OK`; `--- PASS` 6/18/55, 0 FAIL, 0 SKIP, 0 sandbox bind denials. Snapshots 19/0/19, 21/0/21, 28/0/28; each MANIFEST `shasum -c` ok; snapshot M3 `cmp`-equal to planner `st-M3` on 28/28 files.
- `commit_ms.sh`: `7281c090` (NUMSTAT 241), `1e088583` (738), `786f6743` (1,396); 0 refusals; committed tree `cmp`-equal to snapshot M3 (0 differ).
- Out-of-worktree write scan of the diff: the only file create is under `t.TempDir()`.

## Gates (plan §6 `gates.sh` final on `786f6743`, 2,080 s, outside the sandbox)
- gofmt, build, vet, fence, verify_ail, plain + race on host/approval, host/approval/approvaltest, host/authority, host/store, cmd/ailang-worldd, host/daemon, and all 7 inventory legs: rc 0.
- `FAILS: routes_changed caller_fence verify_go` adjudicated:
  - `routes_changed`, `caller_fence`: INSTRUMENT FAILURE — `rg` is not on bash's PATH (`command -v rg` empty), so `DAEMON_ROUTES` and the control printed empty/0. Re-measured with grep: routes head 50 = base 50, `POST /surface/` 3; new production callers none; control `authority.RecheckVerb(` 1 file; fresh negative literal 0.
  - `verify_go` rc 1: 63 `ok`, 1 `FAIL` = `host/broker` `TestExecProtectedCheckTimeout/reserve-after-before` (package byte-identical to base); re-run alone rc 0 (34.7 s).

## Evaluator (Agent `sonnet`, foreground, 55.7 min, 136,633 tok) — PASS 86, 0 blocking
- Report: `evaluator-r1.md` (this directory). Ledger `~/.ailang/state/world-iter428/judge/mutants.tsv` (55 lines; 45 own ids, `comm -12` vs plan ids = 0).
- Controller-verified before recording: replay branch `host/store/operator_cas.go:58` returns the confirmed row after only the digest check (non-blocking, no state change); `go.mod:17` `golang.org/x/sys v0.48.0 // indirect` while `host/authority/judge_operator_silent_test.go:17` imports it directly.

## Gate 3b
- PR #643 head `786f6743a34ceac886ad9497e6e2ff7bbc858622`: 3/3 check-runs success (13:36Z).
- `origin/dev` at merge time `58e47002` (attended #641 + #642 landed mid-sprint; 0 files shared with S3). PR `MERGEABLE CLEAN`; merged `--match-head-commit` → `e83d171754ab036af5218d62f98e6f8206fd7868`.
