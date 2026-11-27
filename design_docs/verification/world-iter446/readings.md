# World iteration 446 — readings (row 243 V4 S7b)

All readings are first-party controller measurements unless marked otherwise. Paths under `~/.ailang/state/world-iter446/` are
rig-local scratch and are not published.

## Gate 0–2
- Kill switch absent; `gh` = `sunholo-voight-kampff`; billing tripwire CLEAN.
- `mission_directives.sh --issue 635`: 0 allowlisted directives (31 comments; LEGACY watermark). `#574` and `#1`: 0
  `MarkEdmondson1234` comments since the mission-scoped watermark `2026-10-27T09:03:42Z`.
- Step 6a `mission_gate0_self_notices.sh` rc 0: 0 crash notices (control 107:4 ok). 3 transcript-scan notices in the window, fires
  `world-1795730668-41608` (2 hits), `world-1795769251-75944` (1 hit) and `world-1795786215-30397` (3 hits). Key names are in the
  private controlplane notice only.
- `mission-world` inbox: 0 unread.
- Local `dev` 8 behind `origin/dev`, fast-forwarded (untracked files only in the main checkout) → `88fea1e9`.
- Skill drift: the pin worktree matched origin on all 14 skill files. The symlink-resolved copy (V1 main checkout) differed from origin
  on `resources/gate-5-retro.md` only (one line: origin adds two idle-fingerprint inputs). Gate 5 read origin's version.
- Ledger valid, 108 rows, 0 OPEN. Clause map: 1.0 clauses 1–7 MET, none UNMET.
- Pick re-measured at `88fea1e9`: `grep -rn 'ConfirmOperatorAction(' --include='*.go' host cmd | grep -v _test.go` → 3 lines (the
  two definitions and `host/authority/judge_operator.go:66`), i.e. no production caller of the store CAS; no S7b branch, PR or
  worktree. Quorum reader rc 0: V4 r2 `…T10-24-51Z` PROCEED 3/3, `absent_reviewers` `[]`.

## Routing
- `resolve-role-spawn.sh`: planner `recipe codex:gpt-6.2-sol declared:planner-lane-default-pin`; executor `recipe codex:gpt-6.2-sol
  declared:provider-pin`; evaluator `agent-tool sonnet declared:alias-pin`. `derive-planner-lane.sh` same as the planner line.
- Codex probe rc 0 (`ok`).

## Baseline (pristine `88fea1e9`, outside the sandbox, isolated `AILANG_STATE_DIR`, while the planner ran)
store-plain 0 · store-race 0 · authority-plain 0 · authority-race 0 · agentrun-plain 0 · agentrun-race 0 · verifygate-plain 1
(`TestMissionConfigDefaultFlowIntegration`) · verbs-plain 0 · broker-plain 3 · cmd-plain 15 · daemon-plain 20 · daemon-race 8.

## Planner (`codex:gpt-6.2-sol`, ~78 min, 625,032 tok, rc 0)
- One milestone M1, +1,191/−72 = 1,263 (cap 1,400). Re-based the S7 plan's Appendix P M3 (517 lines, compiled against the
  pre-round-2 stage): `operator_catalogue_test.go` hunk 2 rejected (landed catalogue subtests at the old EOF; the literal-task pin moved
  into `TestOperatorStartBrowser`), `surface_operator.go` offset −4, `surface_operator_test.go` offset +62.
- Caller census: `JudgeOperatorAction` 0 → 1 and store `ConfirmOperatorAction` 0 → 1, both in `handleOperatorConfirm`
  (`host/daemon/surface_operator_confirm.go:117`, `:141`).
- Own ledger 120 ids (80 diff-anchored, 40 activation dependencies), 0 non-equivalent named survivors; carried J06/J22/J02/J10/J18/J19
  killed by their exact edits (planner-reported; route rows on an in-process recorder, the real listener being denied in the sandbox).
- Controller round trip (`ctl_roundtrip.sh M1`, fresh `git archive 88fea1e9`): `diff -r` rc 0 for cmd, host, docs.

## Executor (`codex:gpt-6.2-sol`)
- Round 1 (74,383 tok) stopped before applying anything, correctly: the plan's §7 block ran `grab "````bash $tool"`, and inside
  double quotes bash treats the four backticks as command substitutions, so the extractor wrote a zero-byte `accept.sh`. The controller
  verified a single-quoted form (`grab '````bash '"$tool"`: accept.sh 69, gates.sh 140, snapshot.sh 40 lines, each `cmp`-equal to the
  planner's scratch copy; the patch `cmp`-equal to the round-trip extract) and re-sent the block with only those two lines changed.
- Round 2 (74,396 tok): patch applied clean; `ACCEPT_OK M1 fast` (257 PASS, 0 FAIL, 0 SKIP); snapshot 22 STATUS / 1 DELETED / 21
  manifest rows, all 21 `cmp`-equal to the planner's stage.

## Commit gate (`commit_ms.sh`, outside the sandbox)
- `accept.sh … outside` → `ACCEPT_OK M1 outside` on all three attempts.
- The planner's `commit-guards.sh` (whole verifygate + broker, isolated state) refused twice on broker fails absent from the base set,
  a different set each time — attempt 1: `TestAilangRunCapsOnTheToolBinary`, `TestAilangRunOutcomeShapesOnTheToolBinary`,
  `TestAilangToolConfinementMatrix`, `TestExecProbeArm8FailsUnsandboxed`, `TestExecSrtIndexWriteStubControl`,
  `TestToolchainCallsWriteOnlyInTheWorktreeAndEpisodeCache` (subprocess timeouts 8/10/30 s, `signal: killed`); attempt 2:
  `TestAilangRunCapsOnTheToolBinary`, `TestExecProtectedCheckTimeout`, `TestExecProtectedIntact`. `host/broker` imports `host/store`
  (`go list -deps`), which M1 changes, so the package argument alone does not clear it.
- Interleaved alone reruns: base FAILED `TestAilangRunCapsOnTheToolBinary` once (round 1, 21.97 s) and passed once; head failed both
  tool-binary tests once and passed once; all three attempt-2 names passed alone on both sides, then head failed
  `TestExecProtectedCheckTimeout` once (42.64 s vs 34.2–34.7 s passes). Negative control fired → attributed to rig load, not M1.
- Committed on `accept.sh … outside` alone (iteration 445 precedent): **M1 `1aa6badc`**, `NUMSTAT 1263`; committed tree = snapshot
  (21 manifest rows, 0 mismatches); `operator_process_linux.go` deleted as planned.

## Controller outside gates on the stage (`gates.sh … outside`)
- OK: gofmt, build, vet, linux vet, inv-store/boundary/cli/daemon/daemon-headers/daemon-read/verbs/verifygate/broker, verify_ail
  (271 s), callers, t1777. `inv-exec-kernel` SKIPPED (Linux runtime control).
- Fail sets vs base: store/agentrun plain+race 0 → 0; verifygate 1 → 0; broker 3 → 1 new (`TestAilangToolConfinementMatrix`, the
  attempt-1 class); cmd-plain 15 → 12, 1 new (`TestPresetCommitterValidationAndPrompt`); daemon-plain 20 → 12, 4 new
  (`TestSeToolsExamplesSearchWithoutCorpusRefuses`, `TestSeToolsWorkspaceExecWithoutProfileRefuses`, `TestSessionSelfGit`,
  `TestWorkspaceRegistryRefusesEpisodesOutsideTheGrammarAndRoot`). All five new cmd/daemon names passed alone at head and at base.
- The accept legs failed on git-fixture worktree tests (`TestOperatorStartOnce/worktree`, `TestOperatorStartCollisions/runner`,
  `TestOperatorWorktreeWindow/head_changed`, `TestOperatorWorktreeDispositionCAS/action_state`, `TestOperatorStartChanged/profile`
  — `git … commit … timed out after 30s`) while the judge was mutating in parallel. Alone: all passed except `TestOperatorStartCollisions`
  at head (248 s; subtests 2–57 s against base's 0.4–0.5 s). Interleaved A/B on it: base 8.00 s PASS, head 7.90 s PASS, **base 104.74 s
  FAIL**, head 76.75 s PASS — base-red alone under rig I/O contention (an unrelated `ailang` at 102% CPU, `godot`, `find`) → not M1.
- `cmd-race` truncated by construction on this rig: pristine whole `cmd/ailang-worldd -race` hit `panic: test timed out after 15m0s`
  (the 11 known-red real-binary e2e tests). `daemon-race` was cut off when the controller's background sweep reached the 2 h tool
  ceiling. Both are UNMEASURED locally; CI's race job is the measurement.
- Added `host/daemon` race time, alone, `-count=2`: `TestOperatorConfirmHTTP` 3.74/3.52 s, `TestOperatorConfirmProjection`
  0.56/0.55 s, `TestOperatorConfirmationAdmissionNeighbours` 0.36/0.25 s → ~4.7 s per run (plan budget 12 s). The planner's
  `outside.sh` read 17.47 s with three heavy legs running concurrently.

## Rig event, attributed
- The real `~/.ailang-world/service.json` was replaced at 20:46:25 (the planner saw it in its read-only before/after observation). The
  writer was an attended `ailang-worldd service upgrade --binary /tmp/ailang-worldd-candidate` (pid 5650, started 21:00:44 from a
  herdr-pane `zsh` started 20:43:33; `~/.ailang-world/bin/ailang-worldd.failed` written 20:45). Not this loop and not any test; left
  alone.

## Judge (Agent `sonnet`, foreground, ~100 min, 84,346 tok)
- **PASS 76** (correctness 26, tests 22, security 17, standards 9, scope 8), 0 blocking. Report
  `~/.ailang/state/world-iter446/evaluator-r1.md`.
- Reachability sound (census matches the planner's); refusals asserted by read-back snapshots; rollback test uses real SQLite triggers;
  `humanMu` held across read/judge/write (removal killed).
- Own mutants: 29 defined, 18 valid; killed E01–E03, E05, E08, E09, E11–E14, E21–E23, E25, E26 (daemon kills on the real
  `httptest` listener); survived E10, E15 (store `BackupEligible` drift check untested), E16, E19 (darwin zombie check), E24, E27, E29
  (several equivalent or shadowed); must-live control E28 survived. Did NOT re-run the six carried ids or the planner's recorder-judged
  route ids.
- Its sequencing file never appeared (the controller sweep's tool ceiling); it waited ~80 min, then measured on a quiet rig.
