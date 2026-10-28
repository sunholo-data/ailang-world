# Ailang World Mission — iteration log

Append-only. One entry per outer-loop iteration. Newest at the bottom.

---

> **Older entries are ARCHIVED.** This file holds the newest 20. The full record of every
> iteration is in `world-mission-log-archive.md`, and a one-line index of ALL of them —
> the thing to grep before picking work, so the loop never repeats itself — is in
> `world-mission-index.md`.

## 282 — 2026-10-19 — row 156 detection redesigned [r4 → r4.2] and M1c + M1d built: Git candidates, mirror authority (PR #312; judge r1 FAIL 68 → r2 PASS 88) [PRODUCT]

**Kind:** design revision + replan + two milestones of one row (156). No release, no semantic weakening.

**Pick and why:** all seven bar clauses are MET, so D-WORLD-81's attended order governs: 156 → 150 → 151 → v1.0.0. Iteration 281's **Next** was a bounded designer revision preserving Unicode semantics, quorum, replan, then M1c. Fresh origin `6eaa997` = local dev; the main checkout sits on an attended branch (`docs/row198-episode-report`, PR #311) and was left untouched — all work ran in `.wt-world-iter282`. Iteration 281's worktrees are landed residue, not orphans.

**Clause map:**
- 1 MET — deterministic kernel.
- 2 MET — local-first daemon.
- 3 MET — explicit authority.
- 4 MET — resident-agent floor.
- 5 MET — provenance value proof.
- 6 MET — protocol boundary.
- 7 MET — controlled self-modification.

**Gate 0:** kill switch armed; gh `sunholo-voight-kampff`; billing CLEAN. #246: 0 allowlisted directives since `2026-10-18T23:25:31Z` (48 comments). Self-notice scan: the one crash notice (`#issuecomment-6136485418`, 2026-10-13) is the orphan already credited by iteration 256 (log grep 1); no new orphan. `mission-world` inbox: 0 unread. Weekly sweep (first fire after Monday 07:00 CEST): 0 orphans of 1 enumerated, list-length control agrees, negative control fired (`design_docs/verification/world-iter282/weekly-sweep.md`).

**Gates 1–2:** origin/dev `6eaa997` CI 3/3 success. Running skill (resolved symlink) and pin copy match origin on SKILL.md + all 13 resources. Index grep: rows 277–281 are this row's history. No open loop PR or iteration-282 worktree before this fire.

**Design (designer Agent opus, resolver `recipe claude:claude-opus-5-5 declared:provider-pin`):** [r4] Git becomes a candidate generator and World's mirror the authority. V17: an exhaustive scan finds exactly 4 non-ASCII scalar values whose NFC+SimpleFold key holds ASCII (U+017F→S, U+037E→;, U+1FEF→`, U+212A→K); controller reran it and checked the scan's fold is byte-for-byte `foldProtectedRune`. V18 Git wildmatch on multi-byte runes (`*` crosses them, `?` is one byte, classes miss), V19 widened pathspecs catch every counterexample with r3's top/bare-D intact, V20 a new r3 defect (`status --ignored=matching` reports every excluded directory, e.g. `! nm/`), V21 cost on a 28 713-file clone, V22 NFD index paths. **Quorum r1** (artifact `…T07-36-12Z`): BLOCKED, gemini-4-argon + oc-glm-5-3 reject, gpt6-1-sol + oc-kimi-k3 unreachable. Controller measured both premises (rule 3f): gemini's multi-rune/ligature fold premise is FALSE for World's mirror (key("ß")="ß", key("ﬀ")="ﬀ"); GLM's ambient-ignored premise is REAL (V10's 2 002 `test_*.py` in `.venv`). One protocol revision [r4.1]: fold stated as NFC-not-NFKC, one shared phase context, an episode-start ignored snapshot (V23 on the real checkouts: 3 259 / 165 / 1 240 standing entries → 0), Unicode-version skew residual. **Quorum r2** (`…T07-43-02Z`): BLOCKED, 3 present (kimi as reserve), all three objections on ONE surface — V23 shows the ignored leg alone at 406–504 ms warm / 2 114 ms cold on ailang against the 500 ms phase bound. Every objection carried a concrete reviewer fix and none disputed the direction, so the controller applied them under the narrow-refinement carve-out as [r4.2]: gemini/kimi's 2 500 ms phase bound with the recomputed child budget, GLM's "max × 1.5" construction bound (3 200 ms from V23), and all three reviewers' demand for a V24 real-worktree measurement gating M6. Doc-level round count: **5** (r1–r3 in iteration 277, two here).

**Plan (planner codex gpt-6.1-sol, `recipe codex:gpt-6.1-sol declared:planner-lane-default-pin`, probe rc0, 130,570 tok):** M1c split into M1c (pure renderer + parser) and M1d (real-Git equivalence) to keep the ≤250-LOC milestone cap; M6-pre/V24 added before M6a; M4b/M6a–c rewritten for three invocations and the snapshot. JSON valid; 21 milestones.

**Execute (executor codex gpt-6.1-sol, `recipe codex:gpt-6.1-sol declared:provider-pin`, probe rc0, 144,910 tok, ≈10 min):** M1c `gitCandidatePathspecs` (replaces `gitPathspecs`), `protectedWidenRunes`, `parseProtectedCandidates` (+140 prod / +110 test); M1d `host/broker/exec_protected_git_test.go` (+196): one disposable repo per spelling, root and subdirectory CWD, fatal without git under `CI`. Executor's 18 mutation trials all killed. Controller re-ran outside the sandbox: gofmt clean, `go vet ./...` rc0, `-run '^TestProtected'` ok (76 Git subtests PASS, 0 SKIP), `-race` ok; rebuilt one commit per milestone from snapshots (`f99ac37`, `e4810dc`; sha256-verified against the final tree, each boundary tested); own mutant (drop `sS` from the widen set) killed by `long-s` in both CWDs while `kelvin` stays green.

**Independent evaluation (Agent sonnet, `agent-tool sonnet declared:alias-pin`, own worktree `.eval-world-iter282`; generator codex/opus ≠ judge sonnet):** r1 **FAIL 68** — BLOCKING: the three over-cap parser rows lacked a trailing NUL, so they died on the NUL check and the cap check could be deleted (controller reproduced: all three PASS with the cap removed); NB stale §4.5 cost sentence; NB one-character widened patterns render `**/*`; NB missing `# branch.oid` parses as head ""; NB the bounds are the controller's application of GLM's split-budget fix. Its independent superset derivation found no hole; 16/18 own mutants killed. Controller test-only fix `ba544c0` (well-formed NUL-terminated oversize streams; cap-removed and cap-at-2× mutants now killed) + doc fixes (cost sentence, R-156-9). r2 **PASS 88**, zero blocking (re-ran three cap mutants, all killed). Reports banked in `design_docs/verification/world-iter282/`.

**Gate 3b:** PR #312 (`sprint/world-iter282-row156-detection`), record committed on the same branch so one CI cycle covers code and record; merge CI verified before Gate 5 reports LANDED.

**Routing evidence:** base=6eaa997f0bfc126c55765aa15337515b6342e633@2026-10-19T07:27:05Z (Gate 1). Controller `claude:claude-opus-5-5` (tok: not reported) · designer Agent opus (2 runs: 101,194 + 120,922 subagent tok; rotation pointer read `codex:gpt-6.1-sol`, resolver returned the driver's `claude:claude-opus-5-5` and was followed verbatim) · planner codex `gpt-6.1-sol` (130,570 tok) · executor codex `gpt-6.1-sol` (144,910 tok) · evaluator Agent sonnet (98,887 + 106,355 subagent tok). judge-independence: cross-vendor. Quorum metered **$0.47016600** (r1 $0.1355, r2 $0.3347; author vendor benched). All four required roles ran; no role silently omitted.

**Ruled out:**
- Narrowing the guarantee (refuse non-ASCII paths, accept a residual, flip fixture rows): a human policy choice, and unnecessary once V17 showed a finite widen set.
- Dropping the pathspec entirely: complete and in budget on a clean clone, but output scales with the tree (2.3 MB `ls-files -v`, 30 002 ignored entries) and would hit the cap on large ignored trees (V21).
- Keeping r3's `status --ignored=matching`: reports every excluded directory as a change (V20).
- Gemini r1's ligature/multi-rune-fold premise: measured false for World's per-rune SimpleFold after NFC.
- Content-hashing the ignored snapshot: 534 ms sha256 on sunholo-platform's 3 259 files (V23) does not fit the phase.
- Parking quorum r2 for Mark: every objection had a reviewer fix on one surface and none disputed direction — the carve-out case, not a judgment call.
- Banking the executor's MUT-CAP-SILENT "killed": the judge's counter-reading reproduced; the executor's mutant evidently did not reach the cap branch.

**Retro:** Skill lane: none (World does not edit the shared skill). Process lane: the executor's self-reported mutation table was wrong for one guard, and only the cross-vendor judge caught it — keep demanding the judge's OWN mutants, not a replay of the executor's list. Backlog lane: none new; V24 is now an explicit plan milestone (M6-pre). Harness share over the latest 20 including 282: **0/20**. Last three landings (282 product, 281 evidence, 280 product) moved no UNMET clause; all clauses are MET and no routable UNMET row exists, so no drift alarm.

**Progress**: World 1.0 **7/7 clauses MET**; **goal unmoved**. Row 156 **6/21 milestones built** (M1c split, M6-pre added); release `v1.0.0-rc.1`, with 156 → 150 → 151 before v1.0.0.

**Next:** row 156 M2a (the `protected` profile key), then M2b. DECISIONS FOR MARK: none; ledger OPEN zero.

## 283 — 2026-10-19 — row 156 M2a + M2b built: the `protected` profile key and startup grammar refusal rows (PR #314; judge r1 PASS 90 → r2 PASS 92) [PRODUCT]

**Kind:** two milestones of one row (156). No release, no semantic weakening, no design or plan change.

**Pick and why:** all seven bar clauses are MET, so D-WORLD-81's attended order governs: 156 → 150 → 151 → v1.0.0. Iteration 282's **Next** was row 156 M2a, then M2b. Both are ~25-minute milestones in the same contract (design M2), so one executor run took both, with per-milestone snapshots. Fresh origin `02bd141` (iteration 282's merge). The main checkout sits on the attended branch `docs/row198-episode-report` (PR #311) and was left untouched; all work ran in `.wt-world-iter283-row156-m2`.

**Clause map:**
- 1 MET — deterministic kernel.
- 2 MET — local-first daemon.
- 3 MET — explicit authority.
- 4 MET — resident-agent floor.
- 5 MET — provenance value proof.
- 6 MET — protocol boundary.
- 7 MET — controlled self-modification.

**Gate 0:** kill switch armed; gh `sunholo-voight-kampff`; billing CLEAN. #313 and #246: 0 allowlisted directives since `2026-10-18T23:25:31Z` (2 and 49 comments). Self-notice scan rc 0 (control 107:4 ok): the one crash notice on #246 (2026-10-13) is outside the window and already credited. `mission-world` inbox: 0 unread of 36. Weekly sweep already run this week by iteration 282.

**Gates 1–2:** origin/dev `02bd141` CI success, check-runs 3, none not-green. The running skill (resolved symlink) and the pin copy both match origin on SKILL.md and all 13 resources (per-file loop, silent). Ledger valid, 70 rows, zero OPEN. Orphan traces: no iteration-283 worktree or branch; iteration 281's worktrees are landed residue; open PRs = #311 only (attended). Quorum: the doc is r4.2, reviewed twice in iteration 282; no change since.

**Routing:** designer not spawned (doc current); planner not spawned (the plan's M2a/M2b were ready; same pattern as iteration 265). Executor resolver `recipe codex:gpt-6.1-sol declared:provider-pin`, probe rc 0; evaluator resolver `agent-tool sonnet declared:alias-pin`.

**Baseline (rule 3e / false-green 4):** the directive's gate list was first run on the pristine worktree and the broker and daemon `-race` suites came back **rc 1**. The failure lines were `AILANG_BIN is unset: … never skips` (`templates_test.go:61`, `episode_test.go:169`), so the instrument was wrong, not the base: CI sets `AILANG_BIN` to the v0.41.0 pin. Re-run with `AILANG_BIN=~/.pinned-ailang/ailang`: broker rc 0 (143 s), daemon rc 0 (216 s). The directive was amended to say so before the spawn.

**Execute (codex gpt-6.1-sol, 151,430 tok, rc 0, ≈20 min):**
- **M2a** `host/broker/exec_profile.go` +30/−1, `exec_profile_test.go` +89. `protected` joins the strict allowlist. The raw value must start with `[` (so `null`, which `encoding/json` would accept into a slice, is refused), and each element must start with `"`. The list goes through `ParseProtectedPatterns`, and `ExecProfile.Protected` / `ProtectedDigest` are stored. Tests: `TestExecProfileProtected` (`absent`, `empty`, `list`, `wrong-type` ×6, `compatibility` with whitespace-padded input so a re-marshalled digest would differ) and `TestExecProfileProtectedUnknownFormRefused`.
- **M2b** `host/daemon/exec_startup_test.go` +48/−3, `docs/QUICKSTART.md` +4. Four rows; each asserts the field is named, no daemon is returned, and there is no probe marker from an instrumented srt shim, no `archive/exec-sandbox`, no `exec-probe`, zero dispatches and zero summaries. Census pinned at 51.
- Executor mutation table: 4/4 killed (MUT-PROFILE-NO-CHECK, -TYPE, -COMPAT, MUT-STARTUP-PROTECTED-IGNORED). It labelled the broker `-race` run UNINFORMATIVE UNDER SANDBOX (denied `$HOME` writes and `ps`). The sandbox did not see `AILANG_BIN`, so the executor fetched and hash-checked the pinned interpreter into `.snap/` (733 MB, deleted before commit; nothing under it was committed).

**Controller verification (outside the sandbox, pinned `AILANG_BIN`, `WORLD_TOOL_AILANG_BIN` unset):** `go vet ./...` rc 0; changed files gofmt-clean; focused broker tests ok; `TestExecStartupRefusalTable` rc 0 (50 subtests PASS, 1 pre-existing real-srt SKIP `probe-package-only`); broker `-race` rc 0; daemon `-race` rc 0. Commits were rebuilt from the snapshots (`ebc6f0f` M2a, `eb46e92` M2b) and are sha256-identical to the executor's final tree. The M2a boundary was tested on its own (broker ok, refusal table ok). Controller mutant: accepting `null` as "no protected set" is killed by `wrong-type/null` alone. Observation: `gofmt -l host/` flags `host/store/deadline_guard_test.go` under go1.26.6; that file is unchanged from base and CI is green on it, so it was left alone.

**Independent evaluation (Agent sonnet, own worktree `.eval-world-iter283`; generator codex ≠ judge sonnet):** r1 **PASS 90** (38/25/15/12), zero blocking. It re-ran the 4 claimed kills and added 8 mutants of its own. A deferred-refusal mutant in `StartExec` is caught by the `archive/exec-sandbox` stat on all four rows. `TrimSpace(data)` hashing is caught by `compatibility`. Dropping canonicalization or the digest is caught by `list`. `len>=0` is an equivalent mutant. **NB1**: renaming the malformed-value field to `prot` survives the daemon row, because the subtest's temp path already contains the word "protected". The controller reproduced it and fixed it doc-only in `9bbcbd7` (the row now quotes `protected: must be an array of strings`), and showed the rename now goes red at `exec_startup_test.go:482`. NB2–NB4 are carried on row 156. r2 (resumed, re-checked on `9bbcbd7`) **PASS 92**. Reports: `design_docs/verification/world-iter283/`.

**Gate 3b:** PR #314 (`sprint/world-iter283-row156-m2`). The record is committed on the same branch, so one CI cycle covers code and record. Merge CI is verified before Gate 5 reports LANDED.

**Routing evidence:** base=02bd141068ce6a1b21179d3a9469e2183b83c3be@2026-10-19T12:11:22Z (Gate 4; Gate 1 read the same SHA at 11:26:18Z). Controller `claude:claude-opus-5-5` (tok: not reported) · designer not spawned (doc current) · planner not spawned (plan ready) · executor codex `gpt-6.1-sol` (151,430 tok) · evaluator Agent sonnet (66,735 + 69,566 subagent tok). judge-independence: cross-vendor. Metered **$0.00** (no quorum, no pi). Every required role for a plan-ready pick ran; no role silently omitted.

**Ruled out:**
- Reading the first baseline's `-race` rc 1 as a red base: the failure lines named the unset `AILANG_BIN`, and the pinned re-run was green.
- Re-quorum at pick: the doc has not changed since iteration 282's two quorum rounds.
- Fixing NB2 (the doubled `protected:` prefix) in this row: the prefix comes from wrapping M1a's refusal text; it is cosmetic and touches the M1a error contract, so it is carried to M7b's guidance pass.
- Reformatting `host/store/deadline_guard_test.go`: out of scope and green on CI.

**Retro:** Skill lane: none (World does not edit the shared skill). Process lane: the memory index already said full `go test` needs `AILANG_BIN`, and the first baseline still omitted it. The failure line, not the rc, told the truth: read the first FAIL line of a baseline red before calling the base red. Backlog lane: none new. Harness share over the latest 20 including 283: **0/20**. Last three landings (283, 282, 281) moved no UNMET clause; all clauses are MET and no routable UNMET row exists, so there is no drift alarm.

**Progress**: World 1.0 **7/7 clauses MET**; **goal unmoved**. Row 156 **8/21 milestones built**. Release `v1.0.0-rc.1`, with 156 → 150 → 151 before v1.0.0.

**Next:** row 156 M3a (append protected sandbox settings without changing old bytes), then M3b. DECISIONS FOR MARK: none; ledger OPEN zero.

## 284 — 2026-10-19 — row 156 M3a + M3b built: protected srt `denyWrite` and real-srt protected-write rows (PR #319; judge r1 PASS 91 → r2 PASS 93) [PRODUCT]

**Kind:** two milestones of one row (156). No release, no semantic weakening, no design or plan change.

**Pick and why:** all seven bar clauses are MET, so D-WORLD-81's attended order governs: 156 → 150 → 151 → v1.0.0. Iteration 283's **Next** was row 156 M3a, then M3b, both design M3 (sandbox rendering). One executor run took both, with per-milestone snapshots. The attended PRs open at Gate 1 (#311, #315, #316, #317) keep 156 first (189 is "NEXT after 156"); #315/#316 merged mid-fire (design docs only, no charter edit) and this branch was rebased onto `5386215` before the record. The main checkout sat on attended branches (`docs/rulings-d84-d85`, then `docs/row198-episode-report`) and was left untouched; all work ran in `.wt-world-iter284-row156-m3`.

**Clause map:**
- 1 MET — deterministic kernel.
- 2 MET — local-first daemon.
- 3 MET — explicit authority.
- 4 MET — resident-agent floor.
- 5 MET — provenance value proof.
- 6 MET — protocol boundary.
- 7 MET — controlled self-modification.

**Gate 0:** kill switch armed; gh `sunholo-voight-kampff`; billing CLEAN. #313: 0 allowlisted directives since `2026-10-18T23:25:31Z` (4 comments, all the loop's). Self-notice scan rc 0 (control 107:4 ok; the one #246 crash notice is outside the window, already credited). `ailang-world` inbox: 1 unread, the v0.54.15 release note (#1872–#1875), informational. Weekly sweep already run this week (iteration 282).

**Gates 1–2:** origin/dev `1f36947` (attended #318 merge), CI pending at read (check-runs 3); running skill (resolved symlink) and pin copy match origin on SKILL.md and all 13 resources (per-file loop). Ledger valid, 70 rows, zero OPEN. Orphan traces: none for 284 (283's worktree is landed residue). Doc unchanged since iteration 282's two quorum rounds → no re-quorum.

**Routing:** designer not spawned (doc current); planner not spawned (plan's M3a/M3b ready). Executor resolver `recipe codex:gpt-6.1-sol declared:provider-pin`, probe rc 0; evaluator resolver `agent-tool sonnet declared:alias-pin`.

**Baseline (rule 3e / false-green 4), pristine `1f36947`, pinned `AILANG_BIN`, `WORLD_TOOL_AILANG_BIN` unset:** `go vet ./...` rc 0; focused broker tests rc 0; broker `-race` rc 0 (242.7 s).

**Execute (codex gpt-6.1-sol, 114,044 tok, rc 0):**
- **M3a** `exec_settings.go` +2, `handlers_exec.go` +1/−1, `exec_settings_test.go` +142. `execSettingsSpec.Protected`; `renderExecSettings` appends M1b's `renderSrtProtected` after the srt defaults; `NewExecHandler` forwards `cfg.Profile.Protected`. Inline protected goldens (darwin, linux); existing golden files unchanged. Tests `TestExecSettingsProtectedRendered`, `TestExecSettingsGoldenDarwin`, `TestExecSettingsGoldenLinux`, `TestExecSettingsNoProtectedIsByteIdentical` (absent and `[]`, digest identical).
- **M3b** `handlers_exec_srt_test.go` +141, `handlers_exec_m2_test.go` +46. `TestExecSrtProtectedWrites`: synthetic operator home (no realHome), profile `tests/**` + `*_test.go`, 18 rows (modify, creates, delete, renames, case, symlink, hard links, parent and protected dir renames, unprotected control); a configured corrupt pin fails before any OS skip. `TestExecProtectedSettingsThroughHandler` (always on, fake srt) reads the settings file the handler passed.
- Executor mutation table: 5/5 killed in-sandbox (MUT-NO-SRT-PROTECTED, -PROTECTED-BEFORE-GITHOOKS, -LINUX-GLOB, -SRT-EMPTY-BYTES, -HANDLER-PROTECTED-DROPPED); the real-srt `/modify` kill was UNMEASURED (srt absent in the sandbox). It reported `TestExecHandlerTimeoutKillsTheGroup` red in-sandbox at base too.

**Controller verification (outside the sandbox):** real srt darwin `TestExecSrtProtectedWrites` 18/18 PASS. MUT-NO-SRT-PROTECTED against it: `/modify` FAIL (`token "WORLD-PROBE-WROTE"`, host hash changed) plus 15 more prevention rows; parent-dir-rename and the control still pass; restored green. `TestExecHandlerTimeoutKillsTheGroup` passes outside the sandbox (sandbox artifact). Commits rebuilt from snapshots; the M3a boundary was tested alone; `shasum -c` OK on all five files. `go vet ./...` rc 0, broker `-race` rc 0, daemon `-race` rc 0 (313 s). Evidence: `design_docs/verification/world-iter284/`.

**Independent evaluation (Agent sonnet, own worktree `.eval-world-iter284`; generator codex ≠ judge sonnet):** r1 **PASS 91** (38/23/15/8/7), zero blocking. It re-ran the five named kills and added five of its own (darwin `<wt>/*/G`, glob without `**`, dir pattern keeping `/**`, append before the srt defaults, the real-srt append deletion). **NB1**: the handler test built its expectation from `f.profile.Protected`, so a profile that lost its set passed both arms. The controller measured it (clearing `p.Protected` after load SURVIVED) and fixed it test-only in `1f598c6` (literal suffix in canonical order); the same mutant now fails `/true`. A first draft of the fix used the wrong literal order and failed clean — caught before commit. r2 (resumed) **PASS 93**, both mutants re-verified. NB2/NB3 carried on row 156.

**Gate 3b:** PR #319 (`sprint/world-iter284-row156-m3`), record on the same branch so one CI cycle covers code and record. Merge CI is verified before Gate 5 reports LANDED.

**Routing evidence:** base=1f369473045f0770472106c64baf1ce728b66e59@2026-10-19T15:26:39Z (Gate 1); rebased onto `5386215` (attended #315/#316, docs only) before the record. Controller `claude:claude-opus-5-5` (tok: not reported) · designer not spawned (doc current) · planner not spawned (plan ready) · executor codex `gpt-6.1-sol` (114,044 tok) · evaluator Agent sonnet (54,532 + 58,212 subagent tok). judge-independence: cross-vendor. Metered **$0.00** (no quorum, no pi). Every required role for a plan-ready pick ran; no role silently omitted.

**Ruled out:**
- Reading the executor's in-sandbox `TestExecHandlerTimeoutKillsTheGroup` red as a base red: the base race suite was green outside the sandbox and the test passes there.
- Re-quorum at pick: the doc has not changed since iteration 282.
- Re-running real srt on Linux locally: no bubblewrap on darwin; CI's Linux job is the measurement.

**Retro:** Skill lane: none (World does not edit the shared skill). Process lane: zsh again did not split `$FILES` (memory already says so); the manifest and checkout steps silently became no-ops and were recovered by per-file `cmp` against the snapshots before any commit — use arrays (`FILES=(…)`) in every multi-file command. Backlog lane: none new. Harness share over the latest 20 including 284: **0/20**. Last three landings (284, 283, 282) moved no UNMET clause; all clauses are MET, so no drift alarm.

**Progress**: World 1.0 **7/7 clauses MET**; **goal unmoved**. Row 156 **10/21 milestones built**. Release `v1.0.0-rc.1`, with 156 → 150 → 151 before v1.0.0.

**Next:** row 156 M4a (refuse configured git filters before any integrity status), then M4b. DECISIONS FOR MARK: none; ledger OPEN zero.

## 285 — 2026-10-19 — row 156 M4a + M4b built: git-filter preflight and the bounded protected integrity checker (PR #322; judge r1 FAIL 78 → r2 PASS 93) [PRODUCT]

**Kind:** two milestones of one row (156). No release, no semantic weakening, no design or plan change beyond the planned §4.5 annotation.

**Pick and why:** all seven bar clauses are MET, so D-WORLD-81's attended order governs: 156 → 150 → 151 → v1.0.0. Iteration 284's **Next** was row 156 M4a, then M4b (design M4). Attended rows 189 M1–M3 merged before this fire (#320, #321); row 189 stays "NEXT after 156", so 156 keeps the slot. All work ran in `.wt-world-iter285-row156-m4`; the main checkout was not touched.

**Clause map:**
- 1 MET — deterministic kernel.
- 2 MET — local-first daemon.
- 3 MET — explicit authority.
- 4 MET — resident-agent floor.
- 5 MET — provenance value proof.
- 6 MET — protocol boundary.
- 7 MET — controlled self-modification.

**Gate 0:** kill switch armed; gh `sunholo-voight-kampff`; billing CLEAN. #313: 0 allowlisted directives since `2026-10-18T23:25:31Z` (5 comments, all the loop's). Self-notice scan rc 0 (control 107:4 ok). `mission-world` inbox: 1 unread, the fleet's `harness-resolved` reply for `skill:design-scaffold-hardcodes-owning-repo` (ghost, fixed by ailang#1542 `b51353bcd`) → row 130 annotated HARNESS-RESOLVED, reply acknowledged. Weekly sweep already run this week (iteration 282).

**Gates 1–2:** origin/dev `fd7b325` (attended #321 merge), CI success 3/3 on it. Running skill (resolved symlink) and pin copy match origin on SKILL.md and all resources (per-file loop, after a fresh fetch of the skill repo). Ledger valid, 72 rows, zero OPEN. No open PRs; no orphan traces for 285. Doc unchanged since iteration 282's quorum rounds → no re-quorum.

**Routing:** designer not spawned (doc current); planner not spawned (plan's M4a/M4b ready). Executor resolver `recipe codex:gpt-6.1-sol declared:provider-pin`; probe from `/tmp` rc 1 (`Not inside a trusted directory` — instrument failure, as memory warns), re-probe from the worktree rc 0. Evaluator resolver `agent-tool sonnet declared:alias-pin`.

**Baseline (pristine `fd7b325`, pinned `AILANG_BIN`, `WORLD_TOOL_AILANG_BIN` unset):** `go vet ./...` rc 0; broker `-race` rc 0 (144 s).

**Execute (codex gpt-6.1-sol, 131,261 tok, rc 0):**
- **M4a** `exec_protected_check.go` (NEW) `protectedGitPreflight`: reuses `configuredGitFilters`; a configured `filter.*` key (local or via `include.path`) or a query error refuses before any status call. `TestProtectedGitFilterPreflight` (configured-clean, configured-process, configured-process-include, no-filter, query-error) drives a spawn-logging fake git and a marker the filter would create. Design §4.5 annotated: the preflight is a config query inside the same provisional 2 500 ms phase bound.
- **M4b** `checkProtected`: preflight, then `status --porcelain=v2 -z --branch …`, `ls-files -z -o -i --exclude-standard`, `ls-files -v -z`, each over `gitCandidatePathspecs`, all under one `protectedCheckTimeout` context, `cmd.Dir`=root, `gitEnv`, `core.fsmonitor=false`, 8 MiB stdout/stderr caps; streams go to M1c's `parseProtectedCandidates`; intact needs HEAD == base and no filtered change, `head-moved` reported, any failure leaves intact nil. `TestProtectedChecker` (8 real-Git rows), `TestProtectedCheckerRootAndBound` (7 fake-Git rows). No caller or result field (M6).
- Executor mutation table 8/8 killed in-sandbox. It flagged two gaps in the reused row-184 helper (96 MiB stdout cap on the config query; `configuredGitFilters`' rc-1-with-output check unreachable) and did not edit the read-only file.

**Controller verification (outside the sandbox):** M4a rebuilt from its snapshot and gated alone; both M4a mutants killed by the controller (the marker WAS created under the first). M4b as delivered: `RootAndBound/cwd` and `/argv` red — the fake git logs `$PWD`=`/private/var/…` while the test expected `t.TempDir()`'s `/var/…` (macOS symlink); the executor's worktree-local TMPDIR hid it. Fixed test-only (`EvalSymlinks`, matching the canonical root `snapshotWorktree` passes). Four controller M4b mutants killed. Snapshots `shasum -c` OK. Full `go vet ./...` rc 0, broker `-race` rc 0. Evidence: `design_docs/verification/world-iter285/`.

**Independent evaluation (Agent sonnet, own worktree `.eval-world-iter285`; generator codex ≠ judge sonnet):** r1 **FAIL 78** (36/15/9/10/8), one blocker: a `checkProtected` that ran the preflight but ignored its result passed every test — the preflight rows call the helper directly and every checker fake answers `config` rc 1. The controller reproduced the survivor. NB2 (load flake): the controller's own `go test ./...` at `167ec9a` had already failed 5 RootAndBound rows on phase timeouts. Fix `d4ff957` (test-only plus `const` → `var`): `TestProtectedChecker/filter-refused` drives the checker over a real `filter.probe.clean` repo (refusal diagnostic, intact nil, marker absent, zero status spawns); every row except `shared-phase-bound` lifts the bound to 30 s via `protectedRelaxBound` (`t.Cleanup` restores it). Survivor now killed; per-call-context mutant still killed; 6 concurrent `-race` copies green. r2 (resumed) **PASS 93**, zero blocking; its two new mutants on the relax helper killed; 8-way stress ×2 green. Its NB1 (the cleanup-restore kill depends on subtest order) and NB2 (unlocked package var) are carried on row 156.

**Gate 3b:** PR #322 (`sprint/world-iter285-row156-m4`), record on the same branch so one CI cycle covers code and record. PR-head CI: `167ec9a` 3/3 green; `d4ff957` race gate red on the row-200 flake (verify + host gates green); the final head (record + row 200, code identical to `d4ff957`) is re-run before merge. Merge CI is verified before Gate 5 reports LANDED.

**Routing evidence:** base=fd7b325fe7fcdb5204cb95debd0c5ff57d22e8c3@2026-10-19T19:26:16Z (Gate 1). Controller `claude:claude-opus-5-5` (tok: not reported) · designer not spawned (doc current) · planner not spawned (plan ready) · executor codex `gpt-6.1-sol` (131,261 tok) · evaluator Agent sonnet (78,166 + 91,203 subagent tok). judge-independence: cross-vendor. Metered **$0.00** (no quorum, no pi). Every required role for a plan-ready pick ran; no role silently omitted.

**Ruled out:**
- Reading the codex probe's rc 1 as a dead lane: it ran from `/tmp` (untrusted, non-git CWD); the same probe from the worktree is rc 0.
- Attributing the `d4ff957` full-run reds (`TestCLIRealSubprocessEpisode`, `TestWorkspaceExecEpisodeSelectsItsProfile`, `TestProducerExecTimeoutKillsTheWholeProcessGroup`) to this PR: they are in packages the PR does not touch, ran under the judge's 8-way stress, and each passes in isolation at the same head; broker + daemon `-race` is green there.
- Fixing `configuredGitFilters` in this PR: row-184 code the plan marks read-only for M4; real `git config --get-regexp` prints nothing on rc 1, so the dead branch is not exploitable today. Carried on row 156.

**Retro:** Skill lane: none (World does not edit the shared skill). Process lane: an executor's in-sandbox green with a worktree-local TMPDIR cannot see macOS's `/var` → `/private/var` symlink — every fake-binary test that compares a logged `$PWD` must canonicalise its temp root; the controller's out-of-sandbox re-run is what caught it. And a test that pins a production deadline (2 500 ms) over several fork-heavy spawns flakes under `go test ./...` load; lift the bound for rows that are not about it. Backlog lane: **row 200** filed — PR #322's CI `go race gate` at `d4ff957` went red on `TestWorktreeSnapshotRefusesAConfiguredFilter` (row 184's test, untouched here): its control step's plain `git status` sometimes skips the clean filter. Negative control at base `fd7b325`: **3/150** `-race` runs fail the same line (1/60 at `d4ff957`), so it predates this PR; it fails safe. M4 residuals carried on row 156. Harness share over the latest 20 including 285: **0/20**. Last three landings (285, 284, 283) moved no UNMET clause; all clauses are MET, so no drift alarm.

**Progress**: World 1.0 **7/7 clauses MET**; **goal unmoved**. Row 156 **12/21 milestones built**. Release `v1.0.0-rc.1`, with 156 → 150 → 151 before v1.0.0.

**Next:** row 156 M4c (wire startup probe arm 8 for prevention and Linux detection), then M5a. DECISIONS FOR MARK: none; ledger OPEN zero.

## 286 — 2026-10-20 — row 156 M4c built: startup probe arm 8, protected-write prevention and Linux detection (PR #325; judge r1 PASS 84 → r2 PASS 91) [PRODUCT]

**Kind:** one milestone of row 156 (M4c). No release, no semantic weakening, no design or plan change.

**Pick and why:** all seven bar clauses are MET, so D-WORLD-81's attended order governs: 156 → 150 → 151 → v1.0.0. Iteration 285's **Next** was row 156 M4c. Open PRs #323/#324 are attended trial branches (rows 198/189), not loop work and not touching row 156. No orphan traces for 286 (no worktree, branch or stamp). All work ran in `.wt-world-iter286-row156-m4c`; the main checkout was not touched.

**Clause map:**
- 1 MET — deterministic kernel.
- 2 MET — local-first daemon.
- 3 MET — explicit authority.
- 4 MET — resident-agent floor.
- 5 MET — provenance value proof.
- 6 MET — protocol boundary.
- 7 MET — controlled self-modification.

**Gate 0:** kill switch armed; gh `sunholo-voight-kampff`; billing CLEAN. #313: 0 allowlisted directives since `2026-10-18T23:25:31Z` (6 comments, all the loop's). Self-notice scan rc 0 (control 107:4 ok). `mission-world` inbox: 0 unread. Weekly sweep already run this week (iteration 282).

**Gates 1–2:** local dev fast-forwarded to origin/dev `c9ab49b` (#322 merge); CI 3/3 success on it. Running skill (resolved symlink) matches origin on SKILL.md and all 13 resources (per-file loop after a fresh fetch). Ledger valid, 72 rows, zero OPEN. Doc unchanged since iteration 282's quorum rounds, so no re-quorum.

**Routing:** designer not spawned (doc current); planner not spawned (plan's M4c ready). Executor resolver `recipe codex:gpt-6.1-sol declared:provider-pin`, probe from the worktree rc 0. Evaluator resolver `agent-tool sonnet declared:alias-pin`.

**Baseline (pristine `c9ab49b`, pinned `AILANG_BIN`, `WORLD_TOOL_AILANG_BIN` unset):** `go vet ./...` rc 0; broker `-race` rc 0 (194 s).

**Execute (codex gpt-6.1-sol, run 1 102,571 tok, rc 0):** `ExecProbeArmProtectedWrite = "arm8-protected-write"`; `execProbeOrder` (arms 1–7, plus arm 8 only when `protected` is non-empty; the refusal text says "of N arms"); `execProbeProtectedSample` (a specimen per pattern, found with the grammar's own matcher; unsupported constructions refused); `probeProtected`/`protectedProbe`: scratch git repo and base commit in the probe episode, planted files proven writable unsandboxed, prevention legs require an accepted refusal errno and unchanged host bytes, the darwin creation leg requires refused and absent, and the linux base-name leg requires WROTE and then `checkProtected` intact=false. `execRefusalCodes("linux")` gains EBUSY/EXDEV (V14). 5/5 planned mutants killed in-sandbox.

**Controller verification (outside the sandbox):** vet, named tests, broker `-race` rc 0, and real srt 0.0.78 on darwin: `TestExecProbeArm8` and `TestExecSrtStartupProbePasses` PASS. Two controller mutants **survived** every test: removing the `probeProtected` call from `probeExecProfile` (every arm-8 test called the arm directly, so "wire" was untested), and dropping the darwin creation leg. Fix round 1 (codex, test-only, 62,422 tok): `TestExecProbeArm8Wired` drives the production entry point (8 vs 7 arms), plus `darwin-create` and a positive control `darwin-refuses-all`. Both mutants now killed (controller re-run). Commit `885af2d`. Evidence: `design_docs/verification/world-iter286/`.

**Independent evaluation (Agent sonnet, own worktree `.eval-world-iter286`; generator codex ≠ judge sonnet):** r1 **PASS 84** (27/19/15/13/10), zero blocking. It killed J-TOTAL; 11 of its own mutants survived, 5 of them on spec'd behaviour (A linux success token, D host-deleted, E darwin created-file absent, I linux directory leg, K EBUSY/EXDEV). The controller measured `execFileRefusal` reading `hostGOOS()` (`exec_settings.go:86` = `runtime.GOOS`) rather than `probeProtected`'s goos, so K cannot be killed on darwin. Fix round 2 (codex, test-only, 61,309 tok; commit `9edf116`): rows `linux-no-token`, `linux-checker-error` (pins `WorktreeFilterRefused`), `linux-directory-write`, `darwin-create-forged`, and positive rows `linux-ebusy-passes`/`linux-exdev-passes`, which run on linux hosts (CI) and skip by name elsewhere. A, E and I killed; C single-drop equivalent (a checker error always carries nil intact; the double drop panics, not a kill); D, B, H and F equivalent. r2 (resumed) **PASS 91**, zero blocking; its N1–N4 all killed.

**Gate 3b:** PR #325 (`sprint/world-iter286-row156-m4c`), with the record on the same branch so one CI cycle covers code and record. PR-head CI: `885af2d` 3/3 green. Merge CI is verified before Gate 5 reports LANDED.

**Routing evidence:** base=c9ab49b71d5d696b4a95f5f56c4ce70bfbdab978@2026-10-19T23:26:05Z (Gate 1). Controller `claude:claude-opus-5-5` (tok: not reported) · designer not spawned (doc current) · planner not spawned (plan ready) · executor codex `gpt-6.1-sol` (3 runs: 102,571 + 62,422 + 61,309 = 226,302 tok) · evaluator Agent sonnet (63,408 + 78,444 subagent tok). judge-independence: cross-vendor. Metered **$0.00** (no quorum, no pi). Every required role for a plan-ready pick ran; no role silently omitted.

**Ruled out:**
- Treating the in-sandbox `^TestExecProbe` rc 1 as a regression: it was `bind: operation not permitted` (sockets inside the sandbox); the same selection is green outside.
- Editing `exec_probe.go` to thread goos into `execFileRefusal` in this PR: the judge rates it non-blocking, and the fix rounds were test-only by directive. Carried on row 156.
- Chasing the C/D/B/H/F survivors: each is shown equivalent by construction (judge-confirmed).

**Retro:** Skill lane: none (World does not edit the shared skill). Process lane: an executor that implements a planned test list kills exactly the planned mutants, and the plan named no mutant for its own milestone title ("Wire") or for the darwin creation leg. The controller's out-of-sandbox mutants found both before the judge saw the code. A directive should demand one production-entry-point test per milestone that *wires* something, and a negative control for every leg the design spells out, not only the plan's listed mutants. Second: a test seam (`goos` parameter) that only part of the code consults gives darwin-green rows that cannot exercise linux refusal codes; the rows are honest only because they skip by name. Backlog lane: none new; M4c residuals carried on row 156. Harness share over the latest 20 including 286: **0/20**. Last three landings (286, 285, 284) moved no UNMET clause; all clauses are MET, so no drift alarm.

**Progress**: World 1.0 **7/7 clauses MET**; **goal unmoved**. Row 156 **13/21 milestones built**. Release `v1.0.0-rc.1`, with 156 → 150 → 151 before v1.0.0.

**Next:** row 156 M5a (protected AILANG policies in every variant), then M5b. DECISIONS FOR MARK: none; ledger OPEN zero.

## 287 — 2026-10-20 — row 156 M5a built: protected AILANG policies in every variant (PR #326; judge r1 PASS 88 → r2 PASS 91) [PRODUCT]

**Kind:** one milestone of row 156 (M5a). No release, no semantic weakening, no design or plan change.

**Pick and why:** all seven bar clauses are MET, so D-WORLD-81's attended order governs: 156 → 150 → 151 → v1.0.0. Iteration 286's **Next** was row 156 M5a; `TestEpisodePolicyProtected*` absent on origin/dev and no M5a PR existed (reality check). Open PRs #323/#324 are attended trial branches (rows 198/189), not loop work. No orphan traces for 287. D-WORLD-82 (D-156-4 = A) puts the World pre-check in M5c/M5d after M5b, so M5a is the policy half only, as planned. All work ran in `.wt-world-iter287-row156-m5a`; the main checkout was not touched.

**Clause map:**
- 1 MET — deterministic kernel.
- 2 MET — local-first daemon.
- 3 MET — explicit authority.
- 4 MET — resident-agent floor.
- 5 MET — provenance value proof.
- 6 MET — protocol boundary.
- 7 MET — controlled self-modification.

**Gate 0:** kill switch armed; gh `sunholo-voight-kampff`; billing CLEAN. #313 and #246: 0 allowlisted directives since `2026-10-18T23:25:31Z` (7 and 49 comments, all the loop's). Self-notice scan rc 0 (control 107:4 ok). `mission-world` inbox: 0 unread (37 total, control). Weekly sweep already run this week (iteration 282).

**Gates 1–2:** local dev = origin/dev `0adeb70` (#325 merge); CI success on it, check-runs 3, none not-green. Running skill (resolved symlink) and the driver pin both match origin on SKILL.md and every resource (per-file loop). Ledger valid, 72 rows, zero OPEN. Design doc unchanged in substance since iteration 282's quorum, so no re-quorum.

**Routing:** designer not spawned (doc current); planner not spawned (plan's M5a ready). Executor resolver `recipe codex:gpt-6.1-sol declared:provider-pin`, probe rc 0. Evaluator resolver `agent-tool sonnet declared:alias-pin`.

**Baseline (pristine `0adeb70`, pinned `AILANG_BIN`, `WORLD_TOOL_AILANG_BIN` unset):** `go vet ./...` rc 0; broker `-race` rc 0 (213 s); daemon `-race` rc 0 (272 s).

**Execute (codex gpt-6.1-sol, 1 run, ≈96k tok, rc 0):** `RenderEpisodePolicy(root, protected)`, `RenderRunVariant(…, protected)`, `runPolicySpec.protected` appended after `episodeDenyWrite` (Env: `.git/**`, built-ins, protected) in a fresh slice; `AilangToolConfig.Protected` copied onto the handler and used by `runVariantFor` and `VerifyRunCaps`; `tomlString` escapes `\`, `"` and control runes; all callers pass nil. Four named tests, including a production-entry-point wiring test (`TestEpisodePolicyProtectedThroughHandler`, a fake policy-tool through `NewAilangToolHandler`), as last iteration's retro asked. 8/8 directed mutants killed in-sandbox; it used a JSON decoder for the deny array because BurntSushi/toml needed a go.sum edit outside its allowlist.

**Controller verification (outside the sandbox):** gofmt clean, vet, 6 + 3 named tests, broker `-race` (185 s), daemon `-race` (230 s), `./cmd/...` rc 0. Five controller mutants anchored on the diff (base short-circuit nil, raw backslash, raw control runes, protected before `.git/**`, handler aliasing the caller's slice), each killed by a named test. Committed `6350cb3`. Evidence: `design_docs/verification/world-iter287/`.

**Independent evaluation (Agent sonnet, own worktree `.eval-world-iter287`; generator codex ≠ judge sonnet):** r1 **PASS 88** (29/22/20/10/7), zero blocking. All 7 executor claims re-verified. Survivor N1: `IsControl` narrowed to `r < 0x20` passed every test, because `encoding/json` accepts a raw DEL that TOML forbids. The controller reproduced it (rc 0), then made a test-only fix `3852dbc`: `/escapes` requires `\u007F` and no raw C0/DEL rune. Both narrowed mutants are now red. r2 (resumed) **PASS 91**, zero blocking; N1 closed; N3 (invalid UTF-8 → U+FFFD) is unreachable through a profile, because `ParseProtectedPatterns` refuses invalid UTF-8 (`exec_protected.go:75`), and is carried to M5b.

**Gate 3b:** PR #326 (`sprint/world-iter287-row156-m5a`), with the record on the same branch. Merge CI is verified before Gate 5 reports LANDED.

**Routing evidence:** base=0adeb70398a1de4ff26f488eb8e384fbe866672a@2026-10-20T03:52:23Z (Gate 4; Gate 1 read the same SHA at 03:27:00Z). Controller `claude:claude-opus-5-5` (tok: not reported) · designer not spawned (doc current) · planner not spawned (plan ready) · executor codex `gpt-6.1-sol` (1 run, ≈96,017 tok) · evaluator Agent sonnet (64,674 + 73,675 subagent tok). judge-independence: cross-vendor. Metered **$0.00** (no quorum, no pi). Every required role for a plan-ready pick ran; no role silently omitted.

**Ruled out:**
- Importing BurntSushi/toml in the test: it is in the module graph but needs a go.sum edit (outside the allowlist) and would meet the dependency allowlist test; the JSON proxy plus the raw-rune scan covers the renderer's subset.
- Wiring the episode profile in the daemon now: that is M5b by plan; M5a callers pass nil.
- Treating the `runVariantFor` copy, the `\u%04x` case and the `len==0 && !nil` form as gaps: each is equivalent by construction (judge-confirmed).

**Retro:** Skill lane: none (World does not edit the shared skill). Process lane: last iteration's lesson (demand a production-entry-point test for anything wired) was put in the directive and held: the wiring test killed the handler/VerifyRunCaps drops on the first run, and no controller mutant survived. The remaining survivor came from the TEST ORACLE: a decoder that is laxer than the format under test (JSON accepts raw DEL, TOML does not) makes an escape test pass for unescaped output. Where a test round-trips through a stand-in decoder, the directive should also demand a byte-level assertion for every character class the real format forbids. Backlog lane: none new; N3 carried on row 156 for M5b. Harness share over the latest 20 including 287: **0/20**. Last three landings (287, 286, 285) moved no UNMET clause; all clauses are MET, so no drift alarm.

**Progress**: World 1.0 **7/7 clauses MET**; **goal unmoved**. Row 156 **14/21 milestones built**. Release `v1.0.0-rc.1`, with 156 → 150 → 151 before v1.0.0.

**Next:** row 156 M5b (episode profile selection and module-root policy rebasing; feed only `ParseProtectedPatterns` output), then the M5c/M5d pre-check (re-measure V7 on v0.54.15 first). DECISIONS FOR MARK: none; ledger OPEN zero.

## 288 — 2026-10-20 — row 156 M5b built: the episode profile's protected set reaches its AILANG policy (PR #328; judge r1 PASS 88 → r2 PASS 91) [PRODUCT]

**Kind:** one milestone of row 156 (M5b). No release, no semantic weakening, no design or plan change.

**Pick and why:** all seven bar clauses are MET, so D-WORLD-81's attended order governs: 156 → 150 → 151 → v1.0.0. Iteration 287's **Next** was row 156 M5b. Reality check: `RebaseProtectedForSandbox` and `workspace_protected_test.go` were absent on origin/dev `0ee1945`, the daemon still passed nil (`workspace.go:616`), no M5b branch or PR existed, and no open PR came from this loop (#323/#324 were attended trial branches and merged before this fire). The iteration-287 worktrees were clean, so there were no orphan traces. All work ran in `.wt-world-iter288-row156-m5b`; the main checkout was not touched.

**Clause map:**
- 1 MET — deterministic kernel.
- 2 MET — local-first daemon.
- 3 MET — explicit authority.
- 4 MET — resident-agent floor.
- 5 MET — provenance value proof.
- 6 MET — protocol boundary.
- 7 MET — controlled self-modification.

**Gate 0:** kill switch armed; gh `sunholo-voight-kampff`; billing CLEAN. #313 and #246: 0 allowlisted directives since `2026-10-18T23:25:31Z` (8 and 49 comments, all the loop's). Self-notice scan rc 0 (control 107:4 ok). `mission-world` inbox: 0 unread (37 total, control). The weekly sweep already ran this week (iteration 282); no thread rotation is due, because #313 was created after Monday's 07:00 local boundary.

**Gates 1–2:** local dev = origin/dev `0ee1945` (#323 merge); CI success on it, check-runs 3, none not-green. The running skill (resolved symlink) matches origin on SKILL.md and on every resource file (per-file loop). Ledger valid, 72 rows, zero OPEN. The design doc is unchanged since iteration 282's quorum, so no re-quorum.

**Routing:** designer not spawned (doc current); planner not spawned (plan's M5b ready). Executor resolver `recipe codex:gpt-6.1-sol declared:provider-pin`, probe rc 0. Evaluator resolver `agent-tool sonnet declared:alias-pin`.

**Baseline (pristine `0ee1945`, pinned `AILANG_BIN`, `WORLD_TOOL_AILANG_BIN` unset):**
- `go vet ./...` rc 0.
- Broker `-race` rc 0 (211 s).
- Daemon `-race` rc 0 (279 s).

**Directive facts the plan did not carry (controller-measured):**
- `rebaseForSandbox` is unexported, so the wrapper needed `exec_protected.go` in the allowlist.
- The AILANG handler cache and its single-flight join keyed on the sandbox only, so a binding minted after the first build would have kept the old policy.
- The plan's two broker-file tests need daemon-package fixtures (`fakeToolBin`, `realToolBin`, the registry), so they moved to package daemon.
- The tool pin is now v0.54.15, not the plan's v0.52.1.
- Semantics where the design is silent: exec off → nil; no usable profile → nil; binding read error → fail closed.

**Execute (codex gpt-6.1-sol, 1 run, ≈112k tok, rc 0):**
- `broker.RebaseProtectedForSandbox`, a wrapper over the existing rebase.
- `episodeHandler` selects the profile through `bindingFor` + `profileFor` and rebases `profile.Protected` onto `Rel(epRoot, sandbox)`.
- The rebased set reaches `RenderEpisodePolicy` and `AilangToolConfig.Protected`.
- The cache and the single-flight join compare `(sandbox, protected)`; a mismatched in-flight build is waited for, and then the call retries.
- Four named tests plus regressions. 7 of 8 directed mutants were killed in-sandbox; MUT-RAW-PATTERNS was honestly reported as unkillable, because re-parsing canonical input is idempotent.

**Controller verification (outside the sandbox):**
- gofmt clean; vet rc 0; named tests 5 PASS + 1 SKIP.
- The **real v0.54.15 tool arm PASSES**: protected write/edit is refused with `fs_deny_write` and host bytes are unchanged; an unprotected write succeeds.
- Broker `-race` (258 s), daemon `-race` (278 s) and `./cmd/...` rc 0.
- Five diff-anchored controller mutants:
  - C1 (single-flight join ignores the set) SURVIVED every executor test. The controller added `TestWorkspaceProtectedPolicyJoinRequiresSameSet`, which kills it 3/3 and passes 3/3 under `-race`.
  - C2–C4 were killed.
  - C5 (config drop) survives the real arm, as expected, and is killed by `…PolicyRequests:227`. That confirms the executor's claim.
- Committed `3f8af26`.
- Evidence: `design_docs/verification/world-iter288/`.

**Independent evaluation (Agent sonnet, own worktree `.eval-world-iter288`; generator codex ≠ judge sonnet):**
- r1 **PASS 88** (27/20/17/14/10), zero blocking, 16 own mutants.
- Survivors:
  - M9 (join ignores the module root): accepted residual. The sandbox derives from startup-static config, so a same-episode root mismatch cannot occur at runtime.
  - M10 (in-flight build does not record its set): the controller reproduced it (rc 0), then made a test-only fix `c8ab0c3`, `TestWorkspaceProtectedPolicyJoinSameSetShares`, which now kills it 3/3 at `:270`.
- r2 (resumed) **PASS 91**, zero blocking.
- Carried non-blocking:
  - The real-tool arm covers write/edit only, not CLI/FS ops.
  - A bound episode whose profile is refused (stale digest / not loaded) keeps unprotected AILANG write/edit while exec is refused. This is a policy call, recorded as a residual.
  - The plan's M5b broker gate selectors would select zero tests.

**Final tree `c8ab0c3`:** vet rc 0, daemon `-race` rc 0 (270 s). Broker production code is unchanged since the r1 run.

**Gate 3b:** PR #328 (`sprint/world-iter288-row156-m5b`), with the record on the same branch. Merge CI is verified before Gate 5 reports LANDED.

**Routing evidence:** base=0ee1945ffba795812f06dce95998b646adfac343@2026-10-20T08:05:04Z (Gate 4; Gate 1 read the same SHA at 07:26:29Z).
- Controller `claude:claude-opus-5-5` (tok: not reported).
- Designer not spawned (doc current); planner not spawned (plan ready).
- Executor codex `gpt-6.1-sol` (1 run, ≈111,731 tok).
- Evaluator Agent sonnet (67,493 + 72,586 subagent tok).
- judge-independence: cross-vendor.
- Metered **$0.00** (no quorum, no pi).
- Every required role for a plan-ready pick ran; no role was silently omitted.

**Ruled out:**
- Treating the `./cmd/...` rc 1 (`TestDemoProofNoVerifiedSetRefuses` MOD010) as a regression. It reproduces only with a non-temp-shaped scratch TMPDIR; with the default TMPDIR the same tree is rc 0.
- Treating the first final-tree `vet rc=1` as a verdict. A zsh `rm -f` glob with no match broke the `&&` chain before vet ran; re-run rc 0.
- Adding a test for M9: a runtime root mismatch for one episode is unreachable (judge agreed).
- Building the World pre-check now: that is M5c/M5d, gated on re-measuring V7 on v0.54.15.

**Retro:**
- Skill lane: none (World does not edit the shared skill).
- Process lane: both survivors this iteration were in the CONCURRENCY branch the milestone edited (single-flight join), which no executor test reached. The executor's mutant list came from the plan, and the plan did not list the cache. Where a milestone edits a cache or a single-flight path, the directive should name both join outcomes (must join, must not join) as required rows. This is rule 3n (anchor the mutation set to the diff) aimed at the directive.
- Instrument lesson: in zsh, `rm -f <glob>` with no match exits 1 and breaks an `&&` chain. Use `rm -f -- <exact paths>` or `setopt nullglob`.
- Backlog lane: none new; three non-blocking carries on row 156.
- Harness share over the latest 20 including 288: **0/20**.
- The last three landings (288, 287, 286) moved no UNMET clause; all clauses are MET, so no drift alarm.

**Progress**: World 1.0 **7/7 clauses MET**; **goal unmoved**. Row 156 **15/21 milestones built**. Release `v1.0.0-rc.1`, with 156 → 150 → 151 before v1.0.0.

**Next:** row 156 — re-measure V7 on the v0.54.15 tool binary (symlink/hard-link arms through policy-tool). If v0.54.15 already refuses them, drop M5c/M5d with the measurement; otherwise plan and build them. Then M6-pre (V24 timing) before M6a. DECISIONS FOR MARK: none; ledger OPEN zero.

## 289 — 2026-10-20 — row 156 V25 re-measure and M5e built: the v0.54.15 tool pin already refuses every alias; M5c/M5d held on D-WORLD-86 (PR #330; judge r1 PASS 90) [PRODUCT]

**Kind:** one measurement (design V25) plus one test-only milestone of row 156 (M5e); a plan amendment; one new decision for Mark. No release, no semantic weakening, no design-direction change.

**Pick and why:** all seven bar clauses are MET, so D-WORLD-81's attended order governs: 156 → 150 → 151 → v1.0.0. Iteration 288's **Next** was: re-measure V7 on the v0.54.15 tool binary, then keep or drop M5c/M5d. Reality check: no open PR from this loop, no orphan worktree with uncommitted row-156 work (the dirty `.eval-*`/`.wt-*` trees are old evaluator residue from iterations 182–281), and nothing on origin/dev touched V7 after iteration 288. Upstream ailang#1873 (filed by World from V7) was closed as fixed by `db9719bb7` in v0.54.15. That closure is a claim, so it was measured rather than taken.

**Clause map:**
- 1 MET — deterministic kernel.
- 2 MET — local-first daemon.
- 3 MET — explicit authority.
- 4 MET — resident-agent floor.
- 5 MET — provenance value proof.
- 6 MET — protocol boundary.
- 7 MET — controlled self-modification.

**Gate 0:** kill switch armed; gh `sunholo-voight-kampff`; billing CLEAN. #313 and #246: 0 allowlisted directives since `2026-10-18T23:25:31Z` (8 and 49 comments, all the loop's). Self-notice scan rc 0 (control 107:4 ok). `mission-world` inbox: 0 unread (37 total). No thread rotation is due (#313 was created 2026-10-19 10:42 local, after Monday's 07:00 boundary; 8 comments).

**Gates 1–2:** local dev = origin/dev `0ffbb9c9` (#329 merge, attended row-189 trial docs); CI success on it, check-runs 3, none not-green. The running skill matches origin on SKILL.md and every resource file (28 comparisons across the resolved symlink and the pin copy; a 30-commit-old control fires). Ledger valid, 72 rows, zero OPEN at pick time.

**V25 (controller measurement, `design_docs/verification/world-iter289/`):**
- Harness: V7's fixture and patterns, plus hard-link arms V7 never ran (D-156-4 names both), plus an unprotected hard link. Each arm runs on a fresh tree, and the verdict is the host bytes read back. The v0.54.15 digest equals the `host/pins` tool role; v0.52.1 is the control arm.
- **v0.52.1 reproduces V7:** symlink-dir (new and existing file), symlink-file and both hard-link aliases are written through into protected paths, write and edit.
- **v0.54.15 refuses all of them, host bytes unchanged.** Symlinks are refused with `(resolves to …): matches fs_deny_write`, hard links with "the file has more than one hard link".
- **Controls:**
  - Without a protected set every v0.54.15 symlink arm is written, so the refusal comes from the protected set.
  - The multi-link refusal is unconditional: it fires even with an EMPTY `fs_deny_write` (R-156-10).
  - `**` is doublestar on v0.54.15: `**/*.test.ts` refuses depth 0, 1 and 2, where v0.52.1 refused depth 1 only.
- **First-run confound, caught before banking:** my first harness hard-linked the symlink arms' targets too, so the v0.54.15 no-pattern control refused two symlink arms for the multi-link reason. I set up hard links only in hard-link arms, re-ran both binaries, and banked the clean run.

**Design and plan:**
- **Design:** V25 row, plus [V25] notes in the status line, §4.1 (the `**/x` refusal text cites V7's v0.52.1-only reading), §4.3.2, P6, D-156-4, the new R-156-10, and §10. The amendment is evidence only, so there was no designer spawn and no re-quorum (precedent: the plan names a controller-authored V24 entry).
- **Decision:** D-WORLD-86 filed rather than dropping M5c/M5d. D-156-4 = A was Mark's attended ruling, and the loop may not drop a ruled milestone on its own.
- **Planner** (codex gpt-6.1-sol, ≈141k tok):
  - Fixed the stale "D-156-4 = B while unanswered" plan line; D-156-1..4 now read answered A (D-WORLD-82). This was a plan/design divergence since 2026-10-17.
  - Added M5e and wrote M5c/M5d as HELD; their killers must use the fake tool, because the real one already refuses.
  - Made M6-pre independent of M5c/M5d/M5e.
  - Marked JSON M1c..M5b done with re-derived merge SHAs.
  - Found that the regression drill needs a temporary `ToolBinaryRelease` admission hunk, because the daemon admits only the pinned release.

**Baseline (pristine `0ffbb9c9`, pinned `AILANG_BIN`, `WORLD_TOOL_AILANG_BIN` unset):**
- `go vet ./host/daemon/` rc 0.
- Named lexical real-tool test: SKIP when unset; PASS 1.45 s on v0.54.15.
- Daemon `-race` rc 0 (229 s).

**Execute (codex gpt-6.1-sol, 1 run, ≈87k tok, rc 0):**
- `TestWorkspaceWriteProtectedAliasesPolicyTool`, 123 LOC in `host/daemon/workspace_protected_test.go` only. 14 arms through World's real `Workspace.Write` handler, each with its own context, fixture and daemon.
- Each protected arm asserts:
  - the mechanism text, plus bytes unchanged (or the target absent);
  - symlink identity, or a hard link with nlink 2 on the same inode.
- Controls: 5 no-protected symlink arms write through; `no-protected-hardlink` is refused (R-156-10); `no-protected-singlelink` writes.
- In-sandbox drills: D1, D2 ×3 and D3, all killed.
- **Finding:** the plan's isolated D3 selector was unanchored at the subtest level. The controller measured it (the second selector picked 5 subtests, not 3) and anchored both in plan and JSON (`0629962a`).

**Controller verification (outside the sandbox, `design_docs/verification/world-iter289/controller-verification.md`):**
- G1 real pin `-count=3`: 42/42 PASS, 0 SKIP.
- Pin unset: discovered SKIP.
- `go vet ./...` rc 0; gofmt clean.
- Lexical test still PASS.
- Daemon `-race` rc 0 (262 s).
- Own drills:
  - D3, the retired v0.52.1 tool (response digest `0dd70a1d…` proves it ran): the 7 alias arms and `no-protected-hardlink` FAIL, the 6 other controls PASS, as predicted.
  - D1, `RenderEpisodePolicy(sandbox, nil)`: the 5 symlink arms FAIL; hard-link arms and controls PASS.

**Independent evaluation (Agent sonnet, own worktree `.eval-world-iter289`; generator codex ≠ judge sonnet; the controller-authored evidence is same-vendor as the judge, and that is FLAGGED):**
- r1 **PASS 90** (28/21/18/13/10), zero blocking.
- It re-ran both V25 harnesses on both binaries: identical to the banked outputs, digest matched.
- Plan totals re-summed (24 rows, 5 407 LOC, 670 min); two merge SHAs re-derived.
- Own mutants:
  - Retired tool and nil-policy: as predicted.
  - `AilangToolConfig.Protected`=nil and un-rebased: both survive M5e but are killed by M5b's `TestWorkspaceProtectedPolicyRequests` / `TestEpisodePolicyProtectedRebased`.
  - One dropped pattern: killed.
  - Neutered aliases: each dies on its mechanism assertion.
- Non-blocking findings:
  - N1: R-156-10 gave the wrong reason ("because the policy is non-empty"; the refusal is unconditional).
  - N2: the pnpm/uv examples are unmeasured.
  - N3: D-WORLD-86 was not yet in the charter. It is filed in this record.
  - N4: the protected hard-link arms are no wiring mutant's sole killer.
  - N5: M5e SKIPs in CI.
- N1/N2 fixed in `5ab76bca` (doc wording). N3 is this record. N4/N5 are carried.

**Gate 3b:** PR #330 (`sprint/world-iter289-row156-v7r`). CI on the first push `1b61eba9` was all 3 jobs SUCCESS. The record rides the same PR. Merge CI is verified before Gate 5 reports LANDED.

**Routing evidence:** base=0ffbb9c9d386f83aaa45b98e4d4b7763a5c5b845@2026-10-20T12:08:21Z (Gate 4; Gate 1 read the same SHA at 11:26:15Z).
- Controller `claude:claude-opus-5-5` (tok: not reported).
- Designer not spawned: evidence-only amendment, controller-measured. The resolver said `recipe claude:claude-opus-5-5 declared:provider-pin`.
- Planner `recipe codex:gpt-6.1-sol declared:planner-lane-default-pin` (1 run, ≈140,904 tok).
- Executor `recipe codex:gpt-6.1-sol declared:provider-pin` (1 run, ≈87,314 tok).
- Evaluator `agent-tool sonnet declared:alias-pin` (100,888 subagent tok).
- judge-independence: cross-vendor for the planner and executor output; same-vendor (Anthropic) for the controller's V25 and design text, mitigated by the judge re-running V25 itself.
- Metered **$0.00** (no quorum, no pi).
- Every required role for this pick ran; the codex probe was rc 0.

**Ruled out:**
- Dropping M5c/M5d on the measurement alone: D-156-4 = A is an attended ruling (D-WORLD-86 asks).
- Building M5c/M5d now: the pin already refuses every target alias, so the work would be unkillable on the real tool and is held.
- Reading the first V25 run's no-pattern refusals as "v0.54.15 bans symlinks": the hard-linked-target confound, removed.
- Treating ailang#1873's "fixed" closure as evidence: it was re-measured as V25.
- Re-quorum: no design-direction change.

**Retro:**
- **Skill lane:** none (World does not edit the shared skill).
- **Process lane:** in a multi-arm fixture, each arm's precondition objects (here a hard link) must exist only in that arm. A global fixture contaminated the control arm in exactly the direction that reads as "the mechanism is broader than it is". The control arm caught it (rule 3d). The planner then wrote the same rule into M5e unprompted.
- **Process lane:** a planner-written drill selector is an instrument. `go test -run` matches each `/`-level separately and unanchored, so `/(x|y)$` over-selects siblings whose names end in x or y. Measure every isolated selector's subtest count before relying on its "controls stay green".
- **Backlog lane:** none new. Carried items are on row 156: the §4.1 message and R-156-3 FS-effect aliases are unmeasured on v0.54.15.
- Harness share over the latest 20 including 289: **0/20**.
- The last three landings (289, 288, 287) moved no UNMET clause; all clauses are MET, so no drift alarm.

**Progress**: World 1.0 **7/7 clauses MET**; **goal unmoved**. Row 156 **16/24 milestones built** (M5c/M5d held on D-WORLD-86). Release `v1.0.0-rc.1`, with 156 → 150 → 151 before v1.0.0.

**Next:** row 156 M6-pre (V24 real-worktree timing of the integrity legs on sunholo-platform and ailang, then size the bounds), then M6a. M5c/M5d only if D-WORLD-86 = B. DECISIONS FOR MARK: D-WORLD-86 (OPEN; default hold).

## 290 — 2026-10-20 — row 156 M6-pre built: design V24 measured on installed worktrees; protected-check bounds re-sized to 3 751 / 4 848 ms (PR #333; judge r1 PASS 78 → r2 PASS 86) [PRODUCT]

**Kind:** one measurement milestone of row 156 (M6-pre): a test-only harness, controller-run evidence, and a design and plan amendment. No production code changed. No release, no semantic weakening, no design-direction change.

**Pick and why:** all seven bar clauses are MET, so D-WORLD-81's attended order governs: 156 → 150 → 151 → v1.0.0. Iteration 289's **Next** was row 156 M6-pre. Reality check:
- No open PR from this loop.
- No orphan worktree for M6-pre; `exec_protected_v24_test.go` and `world-row156-v24/` were absent at base.
- An attended session merged the row 188 re-scope earlier today (D-WORLD-87 RESOLVED attended, provenance `Mark Edmondson` commit `4a6789ac`; design PR #331, S1 plan PR #332). Row 188's S1 is marked attended ("as World episodes"), so the loop does not take it, and it does not change the 156 → 150 → 151 order.

**Clause map:**
- 1 MET — deterministic kernel.
- 2 MET — local-first daemon.
- 3 MET — explicit authority.
- 4 MET — resident-agent floor.
- 5 MET — provenance value proof.
- 6 MET — protocol boundary.
- 7 MET — controlled self-modification.

**Gate 0:**
- Kill switch armed; gh `sunholo-voight-kampff`; billing CLEAN.
- #313: 0 allowlisted directives since `2026-10-18T23:25:31Z` (9 comments, all the loop's).
- Self-notice scan rc 0 (control 107:4 ok).
- `mission-world` inbox: 1 unread, iteration 289's own claim note.
- The weekly sweep already ran this week (iteration 282). No thread rotation is due.

**Gates 1–2:**
- Local dev = origin/dev `17c28324` (#332 merge).
- CI success on `adffe4c1`; on `17c28324` it was in progress at Gate 1, with the race and verify jobs green.
- The running skill matches origin on SKILL.md and all 13 resource files. My first per-file loop printed DIFF for every resource because it passed bare filenames as `resources/` paths. That was an instrument failure; I re-ran it with the right paths and every file read SAME.
- Ledger valid, 74 rows, one OPEN (D-WORLD-86).

**Scratch (controller, outside any sandbox, `~/.ailang/state/world-row156-v24/scratch`):**
- Disposable `git clone --shared --no-checkout` copies of sunholo-platform `fe80714f` and ailang `f0f7df204`. Sharing makes the clones read the source objects without writing to the source repositories.
- One `git worktree add --detach … HEAD` per clone.
- The real checkouts' ignored entries were APFS-clone-copied in (`cp -cRp`, 18 s and 32 s). Agent and mission state, `.env` files and one untracked-not-ignored directory were excluded. That left 51 / 206 entries and 93 732 / 107 711 ignored files.
- Manifests in `<project>.deps`.
- Patterns: V23's sets. The live profiles carry no `protected` key (0 hits across `~/.ailang-world-exec/profiles/*.json`).

**Baseline (pristine `17c28324`, pinned `AILANG_BIN`, `WORLD_TOOL_AILANG_BIN` unset):** `go vet ./host/broker/` rc 0; broker `-race` rc 0.

**Execute (codex gpt-6.1-sol, 1 run, 159,781 tok, rc 0):**
- `TestProtectedV24RealWorktree`, with three always-on synthetic arms and a `measurement` arm that SKIPs unless `WORLD_PROTECTED_V24_ROOT` is set.
- Per project: 1 discarded warm-up, 20 warm samples, and 20 cold samples. Each cold sample is a fresh `worktree add` with the dependency entries renamed in; construction runs, then one phase, then everything is restored.
- A post-snapshot dependency arm.
- The production `checkProtected` is the verdict authority at its own 2 500 ms bound.
- The executor killed the plan's three mutants. Its in-sandbox measurement was not used as evidence (UNINFORMATIVE UNDER SANDBOX), and it left the scratch restored.
- 330 lines against the plan's 200 target.

**Controller verification:**
- Read the harness.
- Own mutant: accepting a CHANGED verdict survived the always-on suite. Fixed with two `reject-null` rows; the mutant is now killed.
- Measured outside the sandbox, twice: both PASS. Run 1's slowest phase was 1 958 ms (ailang cold) and slowest construction 1 583 ms. Run 2's were 1 066 / 1 331.
- Install order: QUICKSTART tells the operator to install dependencies before the session, and the handler is built lazily on the episode's first exec call (`host/daemon/workspace.go` `execHandler`). Nothing enforces the order. The post-snapshot arm, in every run, shows a later ignored protected file reporting `intact: false`, status `ignored`, and `intact: true` again once removed (R-156-6).

**Independent evaluation (Agent sonnet, own worktree `.eval-world-iter290`; generator codex ≠ judge sonnet; the controller-authored design text and measurements are same-vendor as the judge, FLAGGED, mitigated because the judge re-ran the measurement itself):**
- **r1 PASS 78**, zero blocking.
- The judge ran two measurements:
  - Run A, with its own mutation fixtures running at the same time, FAILED closed: one production check timed out (`null`); phase max 2 501 ms; construction max 3 232 ms, above the 3 200 floor.
  - Run B PASSED.
- 29 own mutants; 13 survivors, mostly wiring:
  - The cold class was unguarded: three mutants made it identical to warm.
  - The P50-for-P95 swap, the P95-for-Max swap, a dropped phase leg, a dropped sample-count check and a 1.0 bound factor all survived.
- I verified run A from its JSON before acting on it.
- **Fix `05445bc5`:**
  - Bounds are now sized over ALL five runs. Run A is not excluded; dropping the run that disagrees would be cherry-picking.
  - `v24Bound` and six synthetic rows. All six unit mutants were killed by named subtests, with byte-identical restore.
  - Cold-class guard: a sample's root and ignored snapshot (count + digest) must match. A cold-root mutant now fails on the first cold sample, and the scratch was restored after that mid-run fatal.
  - Controller run 3 on the final harness: PASS.
  - Prose fixed: ignored share 60–95 %, R-156-8 budget arithmetic.
- **r2 PASS 86**, zero blocking.
  - Re-derived both bounds and every quoted number from `result.json`.
  - Its two doc nits are verified and fixed in the record commit: the construction max is the sample whose check then timed out, and a stale §4.5 sentence.
- Remaining survivors are cosmetic: preflight timing, parse snapshot, census, and arm `before`.

**Result (design V24, `design_docs/verification/world-row156-v24/`):**
- Runs 1–3 and B: 320/320 samples intact. Run A: 1 `null`.
- Slowest phase / construction per run: 1 958 / 1 583; 1 066 / 1 331; 954 / 1 052; **2 501 / 3 232** (A); 1 779 / 1 842 (B).
- Bounds, max × 1.5 rounded up: phase **3 751 ms** (provisional 2 500); construction **4 848 ms** (provisional 3 200).
- A protected command keeps ≈ 3.7–5.9 s of the 10 s handler deadline.
- The ignored leg is 60–95 % of a phase. On ailang it walks `docs/node_modules`, which is 2.5 GB.
- Plan: V24 row MEASURED, M6-pre done, M6a/M6b constants stated, JSON status done.

**Gate 3b:**
- PR #333 (`sprint/world-iter290-row156-m6pre`). CI on the first push `0637d1b3`: all 3 jobs SUCCESS.
- Fix and record ride the same PR. Merge CI is verified before Gate 5 reports LANDED.

**Routing evidence:** base=17c283243260fdf7ec1a5a23197f40a872e20282@2026-10-20T16:43:10Z (Gate 4; Gate 1 read the same SHA at 15:26:24Z).
- Controller `claude:claude-opus-5-5` (tok: not reported).
- Designer not spawned: the evidence-only V24 entry is controller-authored by the plan's own terms. The resolver said `recipe claude:claude-opus-5-5 declared:provider-pin`.
- Planner not spawned: the plan was ready. The resolver said `recipe codex:gpt-6.1-sol declared:planner-lane-default-pin`.
- Executor `recipe codex:gpt-6.1-sol declared:provider-pin` (1 run, 159,781 tok; probe rc 0).
- Evaluator `agent-tool sonnet declared:alias-pin` (r1 97,158 + r2 121,616 = 218,774 subagent tok).
- judge-independence: cross-vendor for the executor's harness; same-vendor (Anthropic) for the controller's measurements and design text, mitigated by the judge's own two measurement runs.
- Metered **$0.00** (no quorum, no pi).
- Every required role for this pick ran.

**Ruled out:**
- Sizing the bounds from controller runs 1–3 only (2 937 / 3 200). Run A is a real V24 observation under load, and the 1.5× margin over one unloaded run did not cover it.
- Treating "every run passed" as load-proof: run A failed closed.
- Copying the V1 checkout's agent state into the scratch: `.claude/`, `.ailang/state/` and `.eval_workspace/` were excluded.
- `git worktree add` straight off the live source checkouts: disposable shared clones were used, per the plan.
- Claiming "OS-cold": the page cache was not purged. "Cold" means a fresh worktree's first Git access.

**Retro:**
- **Skill lane:** none (World does not edit the shared skill).
- **Process lane:**
  - When a bound is sized from "observed max × k", the judge's own reproduction under different load is a sample of the same distribution. Fold it in, don't explain it away. Here the 1.5× margin over the controller's runs did not cover the judge's run.
  - A harness whose acceptance takes its bound from the same samples is partly tautological. Its real gates are the verdict checks and the floors, so mutate those first.
- **Process lane:** my resource-file drift loop passed bare filenames. Every DIFF it printed was the instrument, which the identical-SHA HEAD vs origin/dev check exposed. Build compared paths from the same root the `git show` uses.
- **Backlog lane:** none new. Carried items stay on row 156's M6-pre line.
- Harness share over the latest 20 including 290: **0/20**.
- The last three landings (290, 289, 288) moved no UNMET clause; all clauses are MET, so no drift alarm.

**Progress**: World 1.0 **7/7 clauses MET**; **goal unmoved**. Row 156 **17/24 milestones built** (M5c/M5d held on D-WORLD-86). Release `v1.0.0-rc.1`, with 156 → 150 → 151 before v1.0.0.

**Next:** row 156 M6a, applying `protectedCheckTimeout`=3 751 ms and `protectedBaselineTimeout`=4 848 ms, then M6b. M5c/M5d only if D-WORLD-86 = B. DECISIONS FOR MARK: D-WORLD-86 (OPEN; default hold).

## 291 — 2026-10-20 — row 156 M6a built: the protected baseline is captured once at NewExecHandler, with the bounded record shape (PR #338; judge r1 PASS 79 → r2 PASS 88 → r3 PASS 90) [PRODUCT]

**Kind:** one code milestone of row 156 (M6a): production plus tests in `host/broker`, and one pin in `host/store`'s context-root inventory. No release, no semantic weakening, no design-direction change. Nothing is emitted into exec results yet; that is M6b.

**Pick and why:** all seven bar clauses are MET, so D-WORLD-81's attended order governs: 156 → 150 → 151 → v1.0.0. Iteration 290's **Next** was row 156 M6a, with V24's bounds. Reality check:
- No open PR from this loop.
- No orphan worktree or branch for M6a.
- `exec_protected_record.go` was absent at base, and `git log --grep M6a` showed nothing.
- The plan was ready, so no designer or planner was spawned.
- Attended row 188 work (#334–#337, "trial" episodes) merged during the day without touching row 156's files. #337's dogfood friction #22 names a row 156 load flake, carried below.

**Clause map:** 1–7 all MET (deterministic kernel, local-first daemon, explicit authority, resident-agent floor, provenance value proof, protocol boundary, controlled self-modification).

**Gate 0:**
- Kill switch armed; gh `sunholo-voight-kampff`; billing CLEAN.
- #313: 0 allowlisted directives since `2026-10-18T23:25:31Z` (10 comments, all the loop's).
- Self-notice scan rc 0 (control ok).
- `mission-world` inbox: 0 unread (the read list returned rows, so the filter works).
- The weekly sweep already ran this week (iteration 282). No thread rotation is due.

**Gates 1–2:**
- Local dev = origin/dev `d7c647d8` (#336 merge).
- dev CI: 3/3 checks success on `d7c647d8`.
- The running skill matches origin on SKILL.md and all 13 resource files, through both the resolved symlink and the pin copy.
- Ledger valid, 74 rows, one OPEN (D-WORLD-86).

**Baseline (pristine `d7c647d8`, pinned `AILANG_BIN`, `WORLD_TOOL_AILANG_BIN` unset):**
- `go vet ./host/broker/` rc 0.
- The M6a named tests select zero tests at base (absent, not a pass).
- `go test -race ./host/broker/...` at Go's default 10-minute timeout panicked at 600 s, with `TestExecHandlerWorktreeDelta` running for 2 s. The rig was at load 37–49, and the package simply ran past 10 minutes. With `-timeout 1800s` the same package passed in 370 s on the M6a tree. CI gives that job 25 minutes.

**Execute (codex gpt-6.1-sol, 1 run, 121,065 tok, rc 0, ~12 min):**
- `captureProtectedBaseline` runs preflight, `rev-parse`, the mirror-filtered `ls-tree` digest (`<mode> <oid>\t<path>\n` in tree order) and the fold-keyed ignored snapshot (sorted keys, LF-terminated; the empty set hashes empty bytes). Any failure discards partial identities.
- `protectedRecord`/`buildProtectedRecord` produce the bounded design block: 32 changes, `boundPathJSON`, nullable intact verdicts, error precedence baseline → before → after, a 4 KiB error cap.
- `NewExecHandler` captures only for a non-empty set and keeps errors on the handler.
- Bounds are 3 751 / 4 848 ms.
- Unasked: `protectedGitPreflight` queries `git config` directly with the 8 MiB cap, closing the carried M4 judge finding. The refusal text no longer lists the filter names.
- The executor killed the plan's 6 mutants plus 2 more (unwired, no-fold).
- 552 added lines against the plan's 245 cap; recorded as a deviation.
- The fixture Git repositories are built from hand-written loose objects (the directive forbade Git writes). The judge ruled them faithful: real `git` reads them.

**Controller verification (outside any sandbox):**
- Named rows 8/12/1 PASS and the M4 neighbours 34/34. Vet and gofmt clean.
- **Found:** `TestProtectedCheckerRootAndBound/shared-phase-bound` slept 4 × 0.9 s = 3.6 s to force a timeout against a bound that is now 3 751 ms. That leaves ~150 ms of spawn overhead, so it passed here only because the rig was loaded and would likely go red on an idle runner. Fixed `123f666e`: each call sleeps 0.35 × the bound. `TestProtectedV24BoundsApplied` pins 3 751 / 4 848 ms.
- 10 own wiring mutants, 9 killed. The survivor (drop the final `ctx.Err()`) is equivalent.
- **Found:** `./scripts/verify_go.sh` on `555cc281` was red on `host/store` `TestProductionContextRoots`, because the new `context.Background()` did not feed `WithTimeout` and was not pinned. CI would have caught it; nothing in the plan or the executor's gate list did. `host/verifygate` timed out at 10 min in the same run; it does not import `host/broker`, so that is load.

**Independent evaluation (Agent sonnet, own worktree `.eval-world-iter291`; generator codex ≠ judge sonnet; the controller's fixes are same-vendor as the judge, FLAGGED, mitigated by the judge re-mutating every fix):**
- **r1 PASS 79**, zero blocking. Gates rc 0, 49 PASS. ~65 own mutants.
- r1 survivors:
  - The preflight rewrite was undefended (revert it wholesale: all green).
  - Both over-cap rows were vacuous (NUL filler fails the parse with or without the cap).
  - Removing the ignored-key sort survived 7 of 7.
  - The record's identity values and JSON tags were unasserted.
  - The ignored-leg mirror filter was never exercised.
  - The outer bound constant was untested.
  - A flat tree left `-r` untested.
- **Fix `4c8d8933`** (controller; each survivor re-run as a mutant, all killed):
  - The context root feeds `WithTimeout(protectedBaselineTimeout)` and is pinned.
  - `protectedKeysDigest` was extracted, with a sort test that builds the set in reverse 64×.
  - Non-NUL over-cap filler, plus `ErrHandlerOverflow` asserts.
  - An `info/exclude`d `x.text` that the widened pathspec selects and the mirror rejects.
  - `pkg/c_test.go` in the fixture. The judge recomputed the new literal outside Go: byte-identical.
  - `TestProtectedRecordGolden`, for glob and dirs-only.
- **r2 PASS 88**, zero blocking. Finding A: the new outer bound and the capture's inner bound masked each other, so a wrong or removed bound on either side survived.
- `verify_go.sh` on `4c8d8933`: every package passed in the non-race leg (daemon 322.9 s vs 324.2 s before, so no slowdown). The race leg hit the 8-minute per-package timeout on broker, daemon and verifygate (load ~25). In that leg the redesigned `construction-timeout` row went red: its first 0.4 s leg overran the 1 s bound.
- **Fix `88071ca8`:**
  - One bound: the capture uses the caller's root-fed deadline.
  - `construction-timeout`: config and rev-parse are instant, then two 1.2 s legs against a 2 s bound. 5/5 under `-race`.
  - Wrong-constant, 24 h and fresh-per-leg mutants are all killed.
- **r3:** **PASS 90**, zero blocking. The delta is exactly as described. `captureProtectedBaseline` has one production caller, and every leg still carries a per-leg `execTimeout`. Removing the inner bound is ruled safe. Outer-bound mutants (wrong constant, 24 h, removed, 1 ns, 300 ms) and fresh-per-leg mutants are all killed. New cosmetic finding: an expired construction ctx reports the pre-existing `HandlerTimeoutError` text ("timed out after 30s"), not the real bound.

**Result:** M6a built. 18 of 24 row-156 milestones are done (M5c/M5d held on D-WORLD-86). The sprint JSON also gets M5e set to done; it read `planned` despite landing in iteration 289.

**Gate 3b:** PR #338 (`sprint/world-iter291-row156-m6a`). CI on `88071ca8`: `ailang-code verify gate` pass; the go host and race jobs were still running when this record was written (outcome in the report and on the merge commit). The record rides the same PR. Merge CI is verified before Gate 5 reports LANDED.

**Routing evidence:** base=2284096eedfd784751d462f06f322106de1a39b6@2026-10-20T20:54:14Z (Gate 4; Gate 1 read `d7c647d806c0be8bb33f58201267d75bc448211e` at 2026-10-20T19:28:02Z; origin/dev advanced by attended #337, which touches no file in this PR).
- Controller `claude:claude-opus-5-5` (tok: not reported).
- Designer not spawned: the design is ready. The resolver said `recipe claude:claude-opus-5-5 declared:provider-pin`.
- Planner not spawned: the plan is ready. The resolver said `recipe codex:gpt-6.1-sol declared:planner-lane-default-pin`.
- Executor `recipe codex:gpt-6.1-sol declared:provider-pin` (1 run, 121,065 tok; probe rc 0).
- Evaluator `agent-tool sonnet declared:alias-pin` (r1 114,824 + r2 141,388 + r3 152,841 subagent tok).
- judge-independence: cross-vendor for the executor's code; same-vendor (Anthropic) for the controller's three fix commits, mitigated because the judge re-ran them as mutants.
- Metered **$0.00** (no quorum, no pi).
- Every required role for this pick ran.

**Ruled out:**
- Reading the base `-race` panic as a red: the package outran Go's 10-minute default under load, and it passes with an explicit bound.
- Attributing the `host/verifygate` timeouts to this diff: verifygate does not import `host/broker`.
- Attributing the race-leg broker/daemon timeouts to the new capture: the non-race leg's daemon time is unchanged (322.9 s vs 324.2 s).
- Keeping the double construction bound "for safety": each bound hid the other's mutants.
- Reverting the preflight rewrite: the judge ruled it correct and in scope, and the r1 fix now defends it. The filter names are dropped from the refusal text, which is accepted (callers match only the `WorktreeFilterRefused` prefix).

**Carried (non-blocking):**
- From the judge: capture legs' Git env unasserted (shared `gitEnv()`); `--full-tree` equivalent at the repo root; the error cut can split a rune (`json.Marshal` → U+FFFD); LOC over cap.
- From #337's dogfood friction #22: `TestExecProbeArm8FailsUnsandboxed/darwin-create-forged` times out at the 8 s handler budget under a loaded full suite. Candidate: size arm 8's probe budget for a loaded host, as V24 sized the phase bounds.

**Retro:**
- **Skill lane:** none (World does not edit the shared skill).
- **Process lane:**
  - Changing a bound constant must sweep every test whose STIMULUS was sized against the old value. A timeout test with a fixed sleep turns marginal silently when its bound grows (rule 3m).
  - Run the repo's whole `verify_go.sh` before judging, not only the plan's named gates. The plan's gate list could not see a repo-wide guard (`TestProductionContextRoots`) that any new `context.Background()` trips.
  - A timing row redesigned to kill a constant mutant must also be checked under load and `-race` before landing. My first redesign went red within the hour.
- **Process lane:** in zsh, `rm -f prefix.*` with no match broke an `&&` chain for the second time (iteration 288 lesson). Used exact paths after that.
- **Backlog lane:** row 156 gains a carried candidate (arm 8 probe budget under load, friction #22).
- Harness share over the latest 20 including 291: **0/20**.
- The last three landings (291, 290, 289) moved no UNMET clause; all clauses are MET, so no drift alarm.

**Progress**: World 1.0 **7/7 clauses MET**; **goal unmoved**. Row 156 **18/24 milestones built** (M5c/M5d held on D-WORLD-86). Release `v1.0.0-rc.1`, with 156 → 150 → 151 before v1.0.0.

**Next:** row 156 M6b: attach the before/after integrity to exec results while preserving the command deadline (the child budget becomes `min(timeout_ms, until(deadline) − protectedCheckTimeout)`), then M6c. M5c/M5d only if D-WORLD-86 = B. DECISIONS FOR MARK: D-WORLD-86 (OPEN; default hold).

## 292 — 2026-10-21 — row 156 M6b built: exec results carry the before/after protected integrity, with the after-check reserved inside the call deadline (PR #339; judge r1 PASS 80 → r2 PASS 90) [PRODUCT]

**Kind:** one code milestone of row 156 (M6b): production and tests in `host/broker`, plus one pin in `host/store`'s context-root inventory. No release, no semantic weakening, no design-direction change.

**Pick and why:** all seven bar clauses are MET, so D-WORLD-81's attended order governs: 156 → 150 → 151 → v1.0.0. Iteration 291's **Next** was row 156 M6b. Reality check:
- No open PR from this loop.
- No orphan worktree or branch for M6b; one stale `.claude/worktrees/agent-…` tree (`ci-budget`, clean) is unrelated.
- `git log origin/dev --grep M6b` matched only iteration 291's record. The four named tests selected zero tests at base.
- The plan was ready, so no designer or planner was spawned.

**Clause map:** 1–7 all MET (deterministic kernel, local-first daemon, explicit authority, resident-agent floor, provenance value proof, protocol boundary, controlled self-modification).

**Gate 0:**
- Kill switch armed; gh `sunholo-voight-kampff`; billing CLEAN.
- #313: 0 allowlisted directives since `2026-10-18T23:25:31Z` (11 comments, all the loop's).
- Self-notice scan rc 0 (control ok).
- `mission-world` inbox: 0 unread except this loop's own iteration-291 claim.
- The weekly sweep already ran this week (iteration 282). No thread rotation is due.

**Gates 1–2:**
- Local dev fast-forwarded to origin/dev `7fa8b8ee` (#338 merge).
- dev CI: 3/3 checks success on `7fa8b8ee`.
- The running skill matches origin on SKILL.md and all 13 resource files, through both the resolved symlink and the pin copy.
- Ledger valid, 74 rows, one OPEN (D-WORLD-86).

**Baseline (pristine `7fa8b8ee`, pinned `AILANG_BIN`, `WORLD_TOOL_AILANG_BIN` unset):**
- `go vet ./host/broker/` rc 0; gofmt clean.
- The four M6b named tests select zero tests (absent, not a pass).
- The M6a and neighbour regex (baseline, failure, worst case, golden, checker, settings-through-handler, worktree delta) rc 0.

**Execute (codex gpt-6.1-sol, 1 run, 112,329 tok, rc 0, ~11 min):**
- `execResult.Protected *protectedRecord` (`omitempty`), set only for a non-empty protected set.
- In `Execute`: the before-phase runs after `sandbox.Verify`/`verifyFileDigest`/`verifyExecSettings`, `projectDir` and `prepareCache`. Then `execChildBudget(timeout, deadline, hasDeadline, now, reserve)` gives `min(timeout_ms, until(deadline) − protectedCheckTimeout)`.
- The after-phase runs after stdout/stderr capture and before row 184's `probe.finish`, on `WithoutCancel(ctx)` re-capped by `WithDeadline(originalDeadline)`.
- A failed baseline skips both phases and reserves nothing. `(*ExecHandler).checkProtected` resolves `worktreeGit()` per call. `h.spawns` still counts only cache seeding and srt.
- The executor killed its 10 listed mutants. `MUT-ENFORCED-CONST` survived on darwin (Linux-only killer).
- 393 lines changed against the plan's 245 cap; recorded as a deviation. Fixture repos use real `git init/add/commit` (the directive allowed it this time).
- Profile digests hash raw JSON, so `protected: []` and absent already differ in `profile.digest`. The byte-identity test holds the digest constant to compare execution output.

**Controller verification (outside any sandbox):**
- Named tests 5/5 PASS; vet and gofmt clean.
- Two own wiring mutants survived the whole M6b suite: (A) the budget read before the before-phase, so its time is never charged to the child; (C) an after-phase context that inherits the call's cancellation. Fixed `b49c8d67`:
  - `/reserve-after-before`: before 0.6 × bound, after 0.7 × bound.
  - `/after-cancel`: cancels the call while the after-phase runs, and first asserts the cancel really happened.
- Both new rows went red once on the unmutated tree. The timestamped fake-git log showed four git spawns on a tiny repo taking more than 3 s, while `verify_go.sh` and the judge loaded the rig. They then passed 5/5 at load 17. Both now assert only the after verdict, which is what the mutants break; both mutants are killed 2/2.
- **`./scripts/verify_go.sh` on `78111156`: red twice, where the plan's gate list could not see.**
  - `host/store` `TestProductionContextRoots`: the new `WithoutCancel` root was unpinned. This is the second iteration running (291 had the same class with `Background`).
  - `TestExecProtectedCheckTimeout/after-cap` (the executor's row) under the loaded full suite: a fixed 8.25 s hold of row 184's first snapshot plus spawn overhead left the command 37 ms (`timed_out`), and `recordSuccess` rejected it.
- Fixed `b4b3624e`:
  - The pin, with its bound in the comment.
  - `/after-cap` holds the first `-C` call until an absolute `deadline − bound/4` target, keys the stall on the second `status` call, and no longer requires the command's own result.
  - Plus the judge's r1 survivors, below.
- `verify_go.sh` on `b4b3624e`: **rc 0, all legs** (broker 249.6 s non-race and 384.0 s race; daemon 268.4 / 370.0 s).
- Own drill on `b4b3624e`, each mutant against its named row, all **killed**: RESERVE-DOUBLE, BEFORE-BACKGROUND, AFTER-AFTER-FINISH, BEFORE-AFTER-BEGIN, ENFORCED-CONST (helper), AFTER-UNCAPPED (against the rewritten `/after-cap`), STALE-NOW, AFTER-CANCELLABLE. Restored byte-identically.

**Independent evaluation (Agent sonnet, own worktree `.eval-world-iter292`; generator codex ≠ judge sonnet; the controller's fixes are same-vendor as the judge, FLAGGED, mitigated by the judge re-mutating every fix):**
- **r1 PASS 80**, zero blocking; gates rc 0, `-race -count=3` on the timing tests green. 21 own mutants, 15 killed.
  - The judge found the stale-`now` survivor independently. The controller had kept its own finding out of the directive, so the agreement is a second route, not an echo.
  - Survivors: stale `now`; reserve doubled at the call site; Linux `dirs-only` untestable (no seam); before-phase on `context.Background()`; after-phase moved after `probe.finish`; failing on a before-phase check error survives the plan's named killer `TestExecProtectedNotGit`, which only produces a baseline error.
- **Fix `b4b3624e`** also holds:
  - `protectedEnforced(goos)` and `TestProtectedEnforced`.
  - A one-reserve `duration_ms` floor in `/reserve-after-before`.
  - `/before-cap`: a call deadline of bound/3 must end the call.
  - `TestExecProtectedPhaseOrder`: collapsed git-call classes must read protected, delta, protected, delta.
  - `/before-only` at a 15 s bound; the judge saw it time out the after-phase once under load.
- **r2 PASS 90**, zero blocking. The diff matched the controller's list exactly. 14 own mutants, 13 killed; the survivor (a `"darwin"` literal at the call site) is Linux-CI-only, accepted.
  - The judge ruled the after-verdict-only assertions non-vacuous: `/after-cancel` asserts the cancel happened, and the `duration_ms` floor kills the doubled reserve.
  - NB1/2/5/6 closed. NB3 closed for the mapping. NB4 and NB7 are acceptable residuals. NB8 mitigated (cause unproven).

**Result:** M6b built. 19 of 24 row-156 milestones are done (M5c/M5d held on D-WORLD-86). Plan and sprint JSON mark M6b done. The plan records that MUT-CHECK-FAILS-CALL's real killers are the `TestExecProtectedCheckTimeout` rows.

**Gate 3b:** PR #339 (`sprint/world-iter292-row156-m6b`). The record rides the same PR. Merge CI is verified before Gate 5 reports LANDED.

**Routing evidence:** base=7fa8b8ee8157f593fd352eac9da434c2c2ff4ab6@2026-10-21T00:00:21Z (Gate 4; Gate 1 read the same SHA at 2026-10-20T23:26:21Z; origin/dev did not move).
- Controller `claude:claude-opus-5-5` (tok: not reported).
- Designer not spawned: the design is ready. The resolver said `recipe claude:claude-opus-5-5 declared:provider-pin`.
- Planner not spawned: the plan is ready. The resolver said `recipe codex:gpt-6.1-sol declared:planner-lane-default-pin`.
- Executor `recipe codex:gpt-6.1-sol declared:provider-pin` (1 run, 112,329 tok; probe rc 0).
- Evaluator `agent-tool sonnet declared:alias-pin` (r1 69,059 + r2 90,651 subagent tok).
- The operator's standing request asked for Agent-tool roles "exactly as the routing table specifies". The table routes the executor through the codex recipe, not the Agent tool, so it ran there. The evaluator ran through the Agent tool as required.
- judge-independence: cross-vendor for the executor's code; same-vendor (Anthropic) for the controller's two fix commits, mitigated because the judge re-ran them as mutants.
- Metered **$0.00** (no quorum, no pi).
- Every required role for this pick ran.

**Ruled out:**
- Reading the first red of `/reserve-after-before` and `/after-cancel` on the unmutated tree as a code defect: the timestamped log showed a rig stall (four spawns over 3 s on a tiny repo), and 5/5 passed after it.
- Raising `/after-cap`'s phase bound to buy margin: the margin was set by a fixed delay relative to test start. Only an absolute target removes the load dependence.
- Adding a before-phase-error row to `TestExecProtectedNotGit`: a baseline error skips both phases, so a phase-check error cannot occur there. The real killers are recorded instead.
- Making a spent budget after a slow before-phase report-only: that is the pre-existing "budget is spent" rule, and changing it is a design question. Carried as a residual.

**Carried (non-blocking):**
- From the judge: the call-site `"darwin"` literal is Linux-CI-only; `/after-cap` depends on `/usr/bin/perl` for its absolute-time sleep; the timing rows stay scheduling-sensitive (`/before-cap` needs preparation under 2 s).
- Residual: a budget spent by a slow before-phase still returns the "budget is spent" error rather than a report-only result.
- Still carried from 291: arm 8's probe budget under load (dogfood friction #22); M6a's judge carry-overs.

**Retro:**
- **Skill lane:** none (World does not edit the shared skill).
- **Process lane:**
  - `TestProductionContextRoots` caught an unpinned production context root for the **second iteration running** (291 `Background`, 292 `WithoutCancel`). The plan's gate list and the directive's list cannot see it. This iteration's directive warned only about `context.Background()`; the inventory also counts `WithoutCancel` and `TODO`. From now on, any directive whose diff adds a production context root names `go test ./host/store/ -run '^TestProductionContextRoots$'` as a gate. Recorded in the charter's Repo Profile, under the verify gates.
  - A timing row that delays a fixed duration from test start loses its margin to load before the delay. Anchor such delays to an absolute target derived from the deadline.
  - Keeping my own surviving mutant out of the judge's directive worked: the judge found it by its own route.
- **Backlog lane:** row 156 gains the budget-spent residual and the perl note (non-blocking).
- Harness share over the latest 20 including 292: **0/20**.
- The last three landings (292, 291, 290) moved no UNMET clause; all clauses are MET, so no drift alarm.

**Progress**: World 1.0 **7/7 clauses MET**; **goal unmoved**. Row 156 **19/24 milestones built** (M5c/M5d held on D-WORLD-86). Release `v1.0.0-rc.1`, with 156 → 150 → 151 before v1.0.0.

**Next:** row 156 M6c: tamper, ignored-file, flag and moved-HEAD mutants through a pass-through fake srt, the ambient-ignored rows, and the pin-free record replay companion. Then M7a/M7b. M5c/M5d only if D-WORLD-86 = B. DECISIONS FOR MARK: D-WORLD-86 (OPEN; default hold).

## 293 — 2026-10-21 — row 156 M6c built: tamper, ignored-file, index-flag and moved-HEAD detection pinned through the handler, with byte-identical replay (PR #340; judge r1 PASS 80 → r2 PASS 87) [PRODUCT]

**Kind:** one test-only milestone of row 156 (M6c): three test files in `host/broker` and `host/daemon`. No production change, no `.ail` change, no release, no design-direction change.

**Pick and why:** all seven bar clauses are MET, so D-WORLD-81's attended order governs: 156 → 150 → 151 → v1.0.0. Iteration 292's **Next** was row 156 M6c. Reality check:
- No open PR from this loop (`gh pr list --author @me` empty).
- No orphan worktree, branch or uncommitted work for M6c. The only dirty worktree is the iteration-198 design tree (unrelated, old).
- `git log origin/dev --grep M6c` matched only iteration 292's record. `gh pr list --search 'M6c in:title'` was empty. The seven new names selected zero tests at base.
- The plan was ready, so no designer or planner was spawned.

**Clause map:** 1–7 all MET (deterministic kernel, local-first daemon, explicit authority, resident-agent floor, provenance value proof, protocol boundary, controlled self-modification).

**Gate 0:**
- Kill switch armed; gh `sunholo-voight-kampff`; billing CLEAN.
- #313: 0 allowlisted directives since `2026-10-18T23:25:31Z` (12 comments, all the loop's).
- Self-notice scan rc 0 (control 107:4 ok).
- `mission-world` inbox: 2 unread, both this loop's own claims (291, 292).
- The weekly sweep already ran this week (iteration 282). No thread rotation is due (#313 created 2026-10-19 10:42 local, after Monday 07:00; 12 comments).

**Gates 1–2:**
- Local dev = origin/dev `5e8bbb95` (#339 merge). dev CI: 3/3 checks success on `5e8bbb95`.
- The running skill matches origin on SKILL.md and all 13 resource files (resolved symlink and the pin copy).
- Ledger valid, 74 rows, one OPEN (D-WORLD-86).

**Baseline (pristine `5e8bbb95`, pinned `AILANG_BIN`, `WORLD_TOOL_AILANG_BIN` unset):**
- `go vet ./host/broker/ ./host/daemon/` rc 0; gofmt clean.
- The seven M6c names select zero tests (absent, not a pass).
- The M6b/neighbour regex rc 0 (45.9 s).
- With real srt: `TestExecSrtReplayIsByteEqual`, `TestExecSrtGitHeadIsNotWritable` and `TestExecSrtProtectedWrites` PASS.

**Execute (codex gpt-6.1-sol, 1 run, 116,969 tok, rc 0, ~12 min; cap 55 min, sized for a plan-capped milestone with a mutation drill):**
- `TestExecProtectedDetectsCreation`: creation, deletion, rename, replacement. Exact change lists, a sticky `intact_before: false` on a second call, and the command's result retained.
- `TestExecProtectedDetectsIgnoredCreation`: `/gitignore` and `/info-exclude`. Exact ignored entries; the bare directory and `b.txt` are absent (V20).
- `TestExecProtectedIndexFlag`: assume-unchanged and skip-worktree, with an in-test control that `git status` is blind (V13).
- `TestExecProtectedHeadMoved`: an empty commit gives `head-moved` in both phases.
- `TestExecProtectedAmbientIgnoredIntact`: `/standing`, `/new`, `/snapshot-mutation` (R-156-6 residual asserted), `/gitignore` and `/info-exclude`. The snapshot digest is computed independently.
- `TestExecProtectedRecordReplay`: live `NewSession` → `NewReplaySession` over one store. Byte-identical; no stub call, spawn or Git argv; `/unprotected` keeps the old shape.
- `TestExecSrtIndexWriteRefused` (real srt) on a linked worktree, `/sandbox` + `/stub-control`, plus an always-on `TestExecSrtIndexWriteStubControl`.
- The daemon `TestExecSrtReplayIsByteEqual` is split into `/unprotected` and `/protected`. `newExecSrtRig` gained an optional profile hook, so the other rig tests keep their profile.
- The executor's 8 drills were all killed. They are the plan's six, with PROTECTED-REPLAY run two ways (replay dispatching to the registry; `protected` dropped from the stored bytes) and NO-INDEX-FLAG-CHECK run on both subtests.
- +382/−44 lines against the plan's 245 cap; recorded as a deviation.

**Controller verification (outside any sandbox):**
- With real srt and `AILANG_BIN`: all 25 broker tests and subtests PASS, and daemon `/unprotected` + `/protected` PASS.
- No write outside the worktree in the diff. Real-srt fixtures read the operator HOME only through the existing `newSrtFixture`. Fixture git runs under a fake HOME.
- `verify_go.sh` on `0fe9fa6d`: **rc 0**, all legs.
- Own drill on `0fe9fa6d` (check worktree `.check-world-iter293`), each against its named rows. All **killed**, restored byte-identically:
  - lower-case flag branch dropped;
  - `S` branch dropped;
  - `?` untracked add dropped;
  - `D` → modified;
  - call-site nil snapshot (`handlers_exec.go`, distinct from the executor's in-checker one);
  - `--untracked-files=no`;
  - snapshot ignored in the parser.

**Independent evaluation (Agent sonnet, own worktree `.eval-world-iter293`; generator codex ≠ judge sonnet; the controller's fix commit is same-vendor as the judge, FLAGGED, mitigated by the judge re-mutating it):**
- **r1 PASS 80**, zero blocking. Gates rc 0 (real srt included). The plan mutants re-run and killed, ~15 own mutants killed, and 6 of 7 precondition drills die.
  - N1, survivor: `--untracked-files=normal` survives every protected test.
  - N2: replay uses the live registry; asked for an assertion that no handler is reached.
  - N3: the file is 1 081 lines (866 at base).
  - N4: the LOC overrun is unrecorded.
  - N5: the daemon replay does not check the ignored-baseline fields.
- **Reproduced N1 before acting:** confirmed (rc 0 across every `TestExecProtected*`/`TestExecWorktree*`/`TestExecHandler*`).
  - The judge's suggested fixture (`zzdir/a_test.go` under `*_test.go`) still survived the mutant.
  - Measured on real git: under a file-name pathspec, `--untracked-files=normal` still recurses into a new directory. Under a directory pathspec (`tests/**`) it reports only `tests/`.
- **Fix `7dce20d9`** (test-only):
  - `/new-dir` creates `tests/x.py` under `recordPatterns + tests/**`. The mutant now reds it: `[{tests/ untracked after}]` instead of `[{tests/x.py untracked after}]`.
  - A second replay through a registry whose handler calls `t.Error` if dispatched.
- **r2 PASS 87**, zero blocking. The diff matched the controller's list exactly.
  - N1 closed: the judge verified the git behaviour and the kill, and that the shared `tests/**` set did not weaken the other rows.
  - N2 closed: a nil-payload handler-call mutant in `invokeReplay` is caught only by the trap.
  - N3/N4/N5: accepted residuals.
  - New non-blocking: removing the parser's unmerged `u` branch survives the behavioural rows. By inspection it then falls to `default: return fail()`, so it fails closed (error, null intact), never a false intact.
- **`verify_go.sh` on `7dce20d9` (run 1): rc 1, race leg only.** It ran while judge r2 ran mutation drills on the same rig.
  - broker and verifygate hit the 8-minute per-package bound.
  - M6b's `TestExecProtectedCheckTimeout/reserve-after-before` returned at the deadline (42 s for the test).
  - The non-race leg was all `ok`.
  - Negative control, same commit, judge finished: broker `-race` rc 0 (253 s), verifygate `-race` rc 0 (264 s), `TestExecProtectedCheckTimeout -race -count=5` rc 0.
  - Run 2 of the full gate, unloaded: **rc 0, all legs** (broker 250.0 s non-race / 367.9 s race; daemon 262.1 / 348.2 s; verifygate 294.3 / 396.6 s).

**Result:** M6c built. 20 of 24 row-156 milestones are done (M5c/M5d held on D-WORLD-86). Plan and sprint JSON mark M6c done.

**Gate 3b:** PR #340 (`sprint/world-iter293-row156-m6c`). CI on `7dce20d9`: 3/3 pass on Linux, with `WORLD_EXEC_SRT_NODE_MODULES` provisioned; `host/broker` and `host/daemon` `ok`. Per-test lines are visible only for the verbose daemon rows: `TestExecSrtReplayIsByteEqual` `/unprotected` and `/protected` PASS. The record rides the same PR. Merge CI is verified before Gate 5 reports LANDED.

**Routing evidence:** base=5e8bbb954bf78d64b6f9eb3ca7f92180b84be002@2026-10-21T04:43:45Z (Gate 4; Gate 1 read the same SHA at 2026-10-21T03:26:36Z; origin/dev did not move).
- Controller `claude:claude-opus-5-5` (tok: not reported).
- Designer not spawned: the design is ready. The resolver said `recipe claude:claude-opus-5-5 declared:provider-pin`.
- Planner not spawned: the plan is ready. The resolver said `recipe codex:gpt-6.1-sol declared:planner-lane-default-pin`.
- Executor `recipe codex:gpt-6.1-sol declared:provider-pin` (1 run, 116,969 tok; probe rc 0).
- Evaluator `agent-tool sonnet declared:alias-pin` (r1 80,823 + r2 54,949 subagent tok).
- The operator's standing request asked for Agent-tool roles "exactly as the routing table specifies". The table routes the executor through the codex recipe, not the Agent tool, so it ran there. The evaluator ran through the Agent tool as required, both rounds, in the foreground.
- judge-independence: cross-vendor for the executor's code; same-vendor (Anthropic) for the controller's fix commit, mitigated because the judge re-mutated it.
- Metered **$0.00** (no quorum, no pi).
- Every required role for this pick ran.

**Ruled out:**
- The judge's suggested N1 fixture (a `*_test.go` file in a new directory): measured non-discriminating, because a file-name pathspec makes Git recurse in both modes. Re-aimed at a directory pattern instead.
- Attributing run 1's race-leg red to `7dce20d9`: the commit adds one creation row and one replay assertion. Isolated re-runs of every red package on the same commit are green.
- Splitting `exec_protected_record_test.go` now: a file outside the plan's allowlist, and the file was already over the rubric's 800 lines at base. Carried to M7.

**Carried (non-blocking):**
- From the judge:
  - the unmerged `u` branch survivor (fail-closed);
  - `exec_protected_record_test.go` 1 094 lines;
  - the daemon protected replay asserts intact flags and byte equality, not the ignored-baseline fields (the broker replay row does).
- Still carried from 292: the call-site `"darwin"` literal is Linux-CI-only; `/after-cap` needs `/usr/bin/perl`; the timing rows stay scheduling-sensitive, as run 1 showed again; a budget spent by a slow before-phase returns the "budget is spent" error.
- Still carried from 291: arm 8's probe budget under load (dogfood friction #22).

**Retro:**
- **Skill lane:** none (World does not edit the shared skill).
- **Process lane:**
  - Never run the judge's mutation drills while `verify_go.sh` runs on the same rig. Sequence them: `verify_go.sh` first, judge after, or the reverse. Run 1's race-leg red cost a third full-suite run.
  - The same co-load caused stalls in 285 (stress loop), 291 and 292 (judge plus suite), so this is the fourth instance.
  - Recorded in the charter's Repo Profile, beside the iteration-292 verify-gate fix.
- **Process lane (fixture discrimination):** a judge's *suggested* fixture is a claim like any other. Re-run the mutant against it before committing.
  - The suggested N1 row would have shipped green and vacuous.
  - Git's `--untracked-files` collapse is visible only under a directory pathspec, not under a file-name pattern.
- **Backlog lane:** row 156 gains the `u`-branch residual and the test-file split note (non-blocking).
- Harness share over the latest 20 including 293: **0/20**.
- The last three landings (293, 292, 291) moved no UNMET clause; all clauses are MET, so no drift alarm.

**Progress**: World 1.0 **7/7 clauses MET**; **goal unmoved**. Row 156 **20/24 milestones built** (M5c/M5d held on D-WORLD-86). Release `v1.0.0-rc.1`, with 156 → 150 → 151 before v1.0.0.

**Next:** row 156 M7a: `why` shows intact, changed, unchecked and unprotected runs (`cmd/ailang-worldd/why.go`), then M7b (QUICKSTART guidance). M5c/M5d only if D-WORLD-86 = B. DECISIONS FOR MARK: D-WORLD-86 (OPEN; default hold).

## 294 — 2026-10-21 — row 156 M7a built: `why` shows intact, changed, not-checked and none-declared protected exec runs (PR #341; judge r1 PASS 86 → r2 PASS 90) [PRODUCT]

**Kind:** one milestone of row 156 (M7a): `cmd/ailang-worldd/why.go` plus tests. No `.ail` change, no release, no design-direction change.

**Pick and why:** all seven bar clauses are MET, so D-WORLD-81's attended order governs: 156 → 150 → 151 → v1.0.0. Iteration 293's **Next** was row 156 M7a. Reality check:
- No open PR from this loop (`gh pr list --author @me` empty).
- No orphan worktree, branch or uncommitted work for M7a; `git log origin/dev --grep M7a` matched only iteration 293's record.
- `why.go` had no `protected` code at base; the two new selectors chose zero tests (absent, not a pass).
- The plan was ready, so no designer or planner was spawned.

**Clause map:** 1–7 all MET.

**Gate 0:**
- Kill switch armed; gh `sunholo-voight-kampff`; billing CLEAN.
- #313: 0 allowlisted directives since `2026-10-18T23:25:31Z` (13 comments, all the loop's).
- Self-notice scan rc 0 (control 107:4 ok).
- `mission-world` inbox: 2 unread, both this loop's own claims (291, 292).
- No weekly sweep or thread rotation due (#313 created after Monday 07:00).

**Gates 1–2:**
- Local dev fast-forwarded to origin/dev `782ef02e` (#340 merge). dev CI: 3/3 checks success on `782ef02e`.
- The running skill matches origin on SKILL.md and all 13 resource files (resolved symlink).
- Ledger valid, 74 rows, one OPEN (D-WORLD-86).

**Baseline (pristine `782ef02e`, pinned `AILANG_BIN`, `WORLD_TOOL_AILANG_BIN` unset):**
- `go vet ./...` rc 0.
- The two M7a selectors chose zero tests.
- `go test -race ./cmd/ailang-worldd/...` rc 0 (83.6 s).
- `go test -race ./host/broker/...` rc 0.

**Execute (codex gpt-6.1-sol, 1 run, 83,428 tok, rc 0, ~7 min; cap 45 min):**
- `resultProtected` keeps the stored block as `json.RawMessage` (absent = nil; present-but-undecodable is distinct and renders `not checked — malformed protected block`).
- `protectedLine`, in this precedence:
  - CHANGED, naming each phase whose verdict is false, the first 5 `path (status)`, `(+K more)` from `total`, and any error appended;
  - not checked, for a null verdict or an error;
  - intact, only when both verdicts are true and there is no error: `N files @ <base7>, digest <base_digest>, enforced <…>`;
  - `none declared`, for a `Workspace.Exec` result with no key.
- `buildChain` sets the block and an explicit exec identity (`d.Effect == broker.EffectWorkspaceExec`) on the link, only when the result object verified. `printChain` prints the line under the effect link.
- `whyLink.Protected` is `omitempty`: links without a block marshal byte-identically.
- Tests: `TestWhyProtectedLine` (17 renderer rows + `e2e-exec` through `whyRun` → REST → `buildChain` → `printChain`) and `TestWhyProtectedJSON` (verbatim block, literal `null` verdicts, old-link bytes pinned).
- The walk rig gains `invokeEffectResult(…, effect)` (`why_test.go` +6/−2), a controller-approved extension of the plan's two-file allowlist, because `invokeResult` hard-coded `Workspace.Read`.
- Drills: 5 of 7 were killed in-sandbox. The two wiring mutants were socket-bound, so they were labelled UNMEASURED UNDER SANDBOX rather than claimed.
- 240 changed LOC against the plan's 230 cap; recorded as a deviation.

**Controller verification (outside any sandbox):**
- `go vet ./...`, `gofmt`, both named tests (19 `--- PASS`), `-run '^TestWhy'`, and `TestProductionContextRoots`: all rc 0.
- Commit `014ac927`; own drill in check worktree `.check-world-iter294`, each against its named row, restored byte-identically:
  - killed: the two wiring mutants (as `_ = protected`; the first form did not compile and was re-run), exec identity forced false and forced true, base digest → ignored digest, full base SHA, `(+K more)` dropped, CHANGED's error suffix dropped, the after-nil guard dropped, the error test dropped from not-checked, the path cap → 32, a false before-verdict ignored, and decoding outside `resOK` in one form;
  - **survived:** dropping the `IntactBefore == nil` guard. A null before-verdict beside a true after-verdict rendered `intact`, because the only missing-verdict row nulled the after phase.
- **Fix `b008434c`** (test-only): row `missing-before`. The mutant now reds it.
- `verify_go.sh` on `b008434c`: **rc 0**, all legs (cmd/ailang-worldd 116.0 / 132.6 s; broker 224.6 / 267.0 s).

**Independent evaluation (Agent sonnet, own worktree `.eval-world-iter294`; generator codex ≠ judge sonnet; the controller's two test-only commits are same-vendor as the judge, FLAGGED, mitigated by the judge re-mutating the production hunks):**
- **r1 PASS 86**, zero blocking. Gates rc 0, the plan's four mutants re-run and killed, 30 own mutants.
  - N1, survivor: `if resOK {` → `if true {` around the protected decode survived every row. A result whose bytes fail verification would still render a verdict from unverified bytes.
  - N2: 240 LOC over the 230 cap; the `invokeEffectResult` extension was judged warranted.
  - N3: `path_truncated` and the per-change phase are not shown in text.
  - N4: plan-named subtests are renamed; an all-null-verdict record with changes is unreachable from `buildProtectedRecord`.
- **Reproduced N1 before acting:** confirmed (rc 0 across both named tests).
- **Fix `7d24aa70`** (test-only): row `e2e-tampered` serves a forged intact block for a `Workspace.Exec` record through the existing `tamperProxy`. It asserts exit 3, `✗ effect`, no `protected:` line, and no `protected` JSON key. The mutant now reds it ("a tampered result printed a protected line").
- **r2 PASS 90**, zero blocking. The diff matched the controller's exhaustive list (one file, +26).
  - N1 closed: the judge re-ran the mutant and confirmed discrimination.
  - N2/N3/N4: accepted residuals.
  - One new equivalent survivor (`raw == nil` → `len(raw) == 0`).
- **`verify_go.sh` on `7d24aa70`, run 2: rc 1.** One red: the broker race leg's `TestExecProtectedCheckTimeout/reserve-after-before` (`git ls-files` over its 3 s bound). No judge was running; a sibling's headless `godot` held ~74 % CPU (load ~5.4).
  - Negative controls on the same commit: `-race -count=5 -run '^TestExecProtectedCheckTimeout$'` rc 0; isolated `go test -race ./host/broker/` rc 0 (220.6 s).
  - `host/broker` is byte-identical to base `782ef02e`.
- **Run 3: rc 1**, now with subprocess timeouts across `host/broker` and a `TestDemoSeedPreflightRefusals/unknown-session` archive `--version` failure in `cmd/ailang-worldd`.
  - Load was 7–12 from sibling work: an `ailang test` at 234 % CPU, the headless `godot` eval, a llama-server/opencode model run, and the fleet mission's unittest discovery.
  - Controls on the same commit at load 12: `TestDemoSeedPreflightRefusals` isolated rc 0, the whole `cmd/ailang-worldd` package rc 0 (116.7 s).
  - The judge's r2 race run of `cmd/ailang-worldd` was rc 0 (98.4 s).
  - Run 1 (`b008434c`, differing by one test-only subtest) was rc 0 on all legs.

**Result:** M7a built. 21 of 24 row-156 milestones are done (M5c/M5d held on D-WORLD-86). Plan and sprint JSON mark M7a done.

**Gate 3b:** PR #341 (`sprint/world-iter294-row156-m7a`). The record rides the same PR. Merge CI is verified before Gate 5 reports LANDED.

**Routing evidence:** base=782ef02e65789a7791c94f1025e2ccce4f5fd141@2026-10-21T08:25:23Z (Gate 4; Gate 1 read the same SHA at 2026-10-21T07:26:42Z; origin/dev did not move).
- Controller `claude:claude-opus-5-5` (tok: not reported).
- Designer not spawned: the design is ready. The resolver said `recipe claude:claude-opus-5-5 declared:provider-pin`.
- Planner not spawned: the plan is ready. The resolver said `recipe codex:gpt-6.1-sol declared:planner-lane-default-pin`.
- Executor `recipe codex:gpt-6.1-sol declared:provider-pin` (1 run, 83,428 tok; probe rc 0).
- Evaluator `agent-tool sonnet declared:alias-pin` (r1 59,283 + r2 33,401 subagent tok).
- The operator's standing request asked for Agent-tool roles "exactly as the routing table specifies". The table routes the executor through the codex recipe, not the Agent tool, so it ran there. The evaluator ran through the Agent tool as required, both rounds, in the foreground.
- judge-independence: cross-vendor for the executor's code; same-vendor (Anthropic) for the controller's two test-only commits, mitigated because the judge re-mutated the production hunks.
- Metered **$0.00** (no quorum, no pi).
- Every required role for this pick ran.

**Ruled out:**
- Attributing run 2's broker red to this diff: `host/broker` is byte-identical to base, and the row passes ×5 under `-race` on the same commit.
- Attributing run 3's cross-package reds to this diff: they span untouched packages, they coincide with load 7–12 from sibling missions, and the touched package and the failing test are green isolated on the same commit.
- Surfacing `path_truncated` in the text line now: design §4.5 does not require it, `--json` carries it, and the milestone is already over its LOC cap. Carried.

**Carried (non-blocking):**
- From the judge:
  - the text line omits `path_truncated` and the per-change phase when both phases changed;
  - the plan-named subtests `not-checked`/`non-exec` exist under other names;
  - 240 LOC against the 230 cap.
- `TestExecProtectedCheckTimeout/reserve-after-before` went red in the race leg under load in 293 and 294, isolated green both times. Derive its bound from the measured stimulus (rule 3m).
- Still carried from 293: the parser's unmerged `u` branch (fails closed); `exec_protected_record_test.go` 1 094 lines; the daemon protected replay does not check the ignored-baseline fields.
- Still carried from 292: the call-site `"darwin"` literal is Linux-CI-only; `/after-cap` needs `/usr/bin/perl`; a slow before-phase still returns "budget is spent".

**Retro:**
- **Skill lane:** none (World does not edit the shared skill).
- **Process lane:** sibling missions load the rig as well as the judge does.
  - Runs 2 and 3 of `verify_go.sh` went red under load 5–12 that no part of this iteration produced, each time in a different place.
  - Record `uptime` and the top CPU consumers beside every full-suite run. Attribute a red only after the touched package re-runs green in isolation and the red package is checked for byte-identity against base.
  - Added to the charter Repo Profile, beside the iteration-293 sequencing rule (the same class, a second source of load).
- **Process lane (symmetric guards):** a renderer whose rule is "either verdict null → not checked" needs a row for EACH null. The executor's single missing-verdict row covered one phase, and the controller's drill found the other.
- **Backlog lane:** row 156 carries the `/reserve-after-before` flake (2 instances) and the `path_truncated` text residual.
- Harness share over the latest 20 including 294: **0/20**.
- The last three landings (294, 293, 292) moved no UNMET clause; all clauses are MET, so no drift alarm.

**Progress**: World 1.0 **7/7 clauses MET**; **goal unmoved**. Row 156 **21/24 milestones built** (M5c/M5d held on D-WORLD-86). Release `v1.0.0-rc.1`, with 156 → 150 → 151 before v1.0.0.

**Next:** row 156 M7b: QUICKSTART §10 protected-profile guidance, worked patterns for the three projects, and residuals (`docs/QUICKSTART.md`, `host/broker/exec_quickstart_test.go`). M5c/M5d only if D-WORLD-86 = B. DECISIONS FOR MARK: D-WORLD-86 (OPEN; default hold).

## 295 — 2026-10-21 — row 156 M7b built: QUICKSTART §10 protected-profile guidance, bound to the parser, renderers and pins (PR #345; judge r1 PASS 88 → r2 PASS 92) [PRODUCT]

**Kind:** one milestone of row 156 (M7b): `docs/QUICKSTART.md`, `host/broker/exec_quickstart_test.go`, and one refusal string in `host/broker/exec_protected.go`. No `.ail` change, no release, no design-direction change.

**Pick and why:** all seven bar clauses are MET, so D-WORLD-81's attended order governs: 156 → 150 → 151 → v1.0.0. Iteration 294's **Next** was row 156 M7b. Reality check:
- No open PR from this loop (`gh pr list --author @me` empty).
- No orphan worktree or branch for M7b. `git log origin/dev --grep M7b` matched only row 23's unrelated 2026-09-28 tranche and iteration 277's plan commit.
- At base, §10 had no `protected` key in its generator, and `TestQuickstartProtectedContract` selected zero tests (absent, not a pass).
- The plan was ready, so no designer or planner was spawned. Its line anchors were stale (§10 now starts at line 790, not 640), but its section anchor holds.

**Clause map:** 1–7 all MET.

**Gate 0:**
- Kill switch armed; gh `sunholo-voight-kampff`; billing CLEAN.
- #313: 0 allowlisted directives since `2026-10-18T23:25:31Z` (14 comments, all the loop's).
- Self-notice scan rc 0 (control 107:4 ok).
- `mission-world` inbox: 3 unread, all this loop's own claims (291, 292, 294).
- No weekly sweep or thread rotation due (#313 was created after Monday 07:00 local).

**Gates 1–2:**
- Local dev = origin/dev `ea860b0f` (the #343 merge, attended row 188). dev CI: 3/3 checks success.
- The running skill matches origin on SKILL.md and all 13 resource files (resolved symlink).
- Ledger valid, 74 rows, one OPEN (D-WORLD-86).

**Baseline (pristine `ea860b0f` in its own detached worktree `.base-world-iter295`, pinned `AILANG_BIN`, `WORLD_TOOL_AILANG_BIN` unset):**
- `go vet ./...` rc 0. `TestQuickstartSection10ProfilesLoad` PASS. `TestQuickstartProtectedContract` absent.
- `go test -race ./host/broker/...` rc 0 (258.7 s).
- A first attempt, run in the sprint worktree with `AILANG_BIN` unset, went red on `TestEpisodeLiveReplayThreeArmsAndEvidence` ("AILANG_BIN must name the pinned released interpreter"). It was also contaminated: §10's tests read `docs/QUICKSTART.md` at run time while the executor was editing it. That run was discarded (process fix below).

**Execute (codex gpt-6.1-sol, 1 run, 100,064 tok, rc 0, ~8 min; cap 45 min):**
- The generator writes the three worked protected sets: ailang `*_test.go`, `*.golden`, `tests/**`; TwilightGame `tests/**`, `*.test.ts`, `*.test.tsx`, `vitest.config.ts`; sunholo-cli `cli/tests/**`, `conftest.py`, `pyproject.toml`.
  - `*.test.tsx` was added on the controller's directive: the profile's test-file command admits `.test.tsx`, and design §4.1's own example lists it. This is a deviation from the plan's three-item list.
- A "Protect the tests" subsection covers:
  - the two forms;
  - worktree-relative paths, with the `cli/tests/**` case;
  - a marked 18-row pattern table (accepted/refused);
  - a marked per-OS table;
  - report-only and no default;
  - a marked record example;
  - the four `why` lines;
  - the alias residual on the v0.54.15 tool pin;
  - the after-check reserve, and arm 8 in the probe list ("seven arms" → "plus an eighth").
- `TestQuickstartSection10ProfilesLoad` gains `protected-go/-ts/-python`: exact canonical sets and digests.
- `TestQuickstartProtectedContract` has five rows:
  - `patterns`: each documented pattern goes through `ParseProtectedPatterns`;
  - `enforcement`: `protectedEnforced` and `renderSrtProtected` per OS;
  - `record`: strict decode into `protectedRecord`, plus a reflection key-set check;
  - `no-default`;
  - `refusal-message`.
- `exec_protected.go`: the `**/x` refusal now ends "(a protected set has exactly two forms, so every matcher reads it the same way)" instead of citing the v0.52.1 one-level reading. This is the controller-approved allowlist extension (design §4.1 [V25]).
- Drills: the plan's three mutants plus five of the directive's, all killed in-sandbox. `cmd/ailang-worldd` Quickstart was UNINFORMATIVE UNDER SANDBOX (bind denied).

**Controller verification (outside any sandbox):**
- `gofmt`, `go vet ./...` and both named tests: all PASS.
- `cmd/ailang-worldd -run Quickstart`: 5 tests selected, 5 PASS. `host/daemon` `TestExecStartupRefusalTable` PASS. `TestProductionContextRoots` and the broker `Protected` rows: rc 0.
- Commit `979b866e` (executor). Own drill, 17 mutants via `/tmp/m7b_iter295_mut/run.sh`, which asserts LANDED and BUILDS, runs only the named killer, and restores by `git checkout`:
  - killed: the CLI root, the empty pattern table, the darwin row dropped, a record key dropped or added, the Linux base-name cell saying "prevented", a code default when the key is absent, Linux rendering base-name globs, and `protectedEnforced("linux")` → glob;
  - **survived:** the doc's 64/256 limits, the 3 751 ms reserve and the v0.54.15 pin changed; `protectedMaxEntries` and `protectedCheckTimeout` changed in code; the `why`-line wording; the arm-8 sentence.
- **Fix `76ec06f5`** (test-only):
  - a `numbers` row checks the documented limits, reserve and tool pin against `protectedMaxEntries`, `protectedMaxBytes`, `protectedCheckTimeout` and `pins.ReleaseFor(RoleTool)`;
  - the no-default row gains a `[]` arm. A code default for `[]` (C17) was killed only by that arm.
  - Re-run: C7/C8/C9/C15/C16/C17 all killed.
- `verify_go.sh` on `76ec06f5`: **rc 0**, all legs, 791 s. Load was 8–10 (a sibling's headless `godot` at ~74 % CPU).

**Independent evaluation (Agent sonnet, own worktree `.eval-world-iter295`; generator codex ≠ judge sonnet; the controller's two test-only commits are same-vendor as the judge, FLAGGED, and the directive told the judge to scrutinise them equally):**
- **r1 PASS 88**, zero blocking. All gates rc 0. Broker race 253 s. 35 mutants: 27 killed, 8 survived.
- Findings, each reproduced against the code before acting:
  - two overstated sentences: the reserve applies only when the baseline was captured (`handlers_exec.go:286`), and an unprotected profile keeps its own `timeout_ms`, not a fixed 9 s;
  - the plan's "worktree-relative vs module-sandbox" note was missing;
  - the row parser skipped a malformed `|`-row (MUT-PAT-nospace);
  - rows moved outside the markers went unchecked;
  - the no-default check depended on a hard wrap;
  - unbound prose: the `why` lines, the alias and arm-8 paragraphs, the record `status` value, and "10 s".
- **Fix `2f4db287`** (doc + test):
  - exact budget wording;
  - a module-root re-basing sentence, bound to `RebaseProtectedForSandbox`;
  - every non-blank block line must be a well-formed row, with one separator allowed;
  - a table may not continue past its markers;
  - the no-default check compares whitespace-normalised text.
  - Controller re-run: J1/J2/J4/J5 killed. J3 (a reflow) survived as intended.
- On `2f4db287`, `./cmd/...` and `host/daemon` (non-race) were rc 0, and the full broker package was rc 0 (180.9 s).
- **r2 PASS 92**, zero blocking. Broker race 243 s. Round-1 findings 1, 2, 6 and 7 closed. Its 10 new mutants were killed. Survivors were harmless or prose-only: a duplicate separator, the re-base and reserve wording, and the `why` text.

**Result:** M7b built. Row 156 has 22 of 24 milestones done; M5c/M5d are held on D-WORLD-86, which also gates the row's close and the doc's move to `implemented/`. Plan and sprint JSON mark M7b done. The plan's stale M5e row now reads done (iteration 289).

**Gate 3b:** PR #345 (`sprint/world-iter295-row156-m7b`). The record rides the same PR. Merge CI is verified before Gate 5 reports LANDED. origin/dev moved during the iteration (#344, attended row 188 M3, a disjoint file set).

**Routing evidence:** base=9c3c38bd94fa5f40341d8abcd97ac03e21fc19f2@2026-10-21T12:27:49Z (Gate 4). Gate 1 read `ea860b0f` at 2026-10-21T11:27:06Z; origin/dev moved by #344 (row 188 M3, attended), which touches no file of this iteration's.
- Controller `claude:claude-opus-5-5` (tok: not reported).
- Designer not spawned: the design is ready. The resolver said `recipe claude:claude-opus-5-5 declared:provider-pin`.
- Planner not spawned: the plan is ready. The resolver said `recipe codex:gpt-6.1-sol declared:planner-lane-default-pin`.
- Executor `recipe codex:gpt-6.1-sol declared:provider-pin` (1 run, 100,064 tok; probe rc 0).
- Evaluator `agent-tool sonnet declared:alias-pin` (r1 89,250 subagent tok; r2 reported 100,653 subagent tok, resumed by `SendMessage`).
- The operator's standing request asked for Agent-tool roles "exactly as the routing table specifies". The table routes the executor through the codex recipe, so it ran there. The evaluator ran through the Agent tool as required, r1 in the foreground and r2 resumed in the same agent.
- judge-independence: cross-vendor for the executor's work; same-vendor (Anthropic) for the controller's two test-only commits, mitigated by the judge's own mutants on them.
- Metered **$0.00** (no quorum, no pi).
- Every required role for this pick ran.

**Ruled out:**
- Reading the first baseline red as a base defect: it was `AILANG_BIN` unset, in a worktree being edited. The pinned re-run in a detached base worktree was rc 0.
- Binding the `why` lines to `protectedLine` in this milestone: the renderer is in package `main` (`cmd/ailang-worldd`), outside M7b's allowlist and package. Carried.
- Treating the plan's 210-LOC budget as blocking: the overage is the marked tables and the contract rows. The judge judged each deviation warranted.

**Carried (non-blocking):**
- From the judge:
  - `why` lines, the alias and arm-8 paragraphs, the report-only wording, the record example's `status` value and "10 s handler deadline" are prose not bound to code;
  - the re-base sentence is bound only for its `cli/tests/**` example;
  - the generator rows SKIP without python3 (pre-existing).
- Still carried from 294: `path_truncated` and the per-change phase are absent from the `why` text; `/reserve-after-before` reds under load.
- Still carried from 293: the parser's `u` branch; `exec_protected_record_test.go` at 1 094 lines; the daemon replay's ignored-baseline fields.
- **[HARNESS] ticket:routing-models-check:bold-colon-heading** (filed to mission-fleet, blocking none): Gate 4's `mission_routing_models_check.sh` matches `**Routing evidence**` but World writes `**Routing evidence:**`, so it returns rc 2 (instrument failure) on entry 295 and on the control, entry 294. A manual read of entry 295 finds only `gpt-6.1-sol`, which is allowlisted.

**Retro:**
- **Skill lane:** none (World does not edit the shared skill).
- **Process lane (baselines):** take every baseline in its own `git worktree add --detach <path> <base>`, with `AILANG_BIN` exported, never in the tree a role is editing. Tests that read checked-in docs at run time see the executor's half-written files. Added to the charter Repo Profile beside the 293/294 load rules.
- **Process lane (doc milestones):** a doc-contract test should bind every FIGURE the prose quotes (limits, bounds, pins) to the constant that owns it. The executor bound the tables and none of the numbers, and five figure mutants survived until the `numbers` row. Name it in the directive the next time a milestone ships guidance.
- **Backlog lane:** none new. Row 156's remaining work is D-WORLD-86.
- Harness share over the latest 20 including 295: **0/20**.
- The last three landings (295, 294, 293) moved no UNMET clause. All clauses are MET, so there is no drift alarm.

**Progress**: World 1.0 **7/7 clauses MET**; **goal unmoved**. Row 156 **22/24 milestones built**; nothing buildable remains (M5c/M5d held on D-WORLD-86). Release `v1.0.0-rc.1`, with 156 → 150 → 151 before v1.0.0.

**Next:** row 150 M5b: TESTED/PROVEN resolution and rendering, and settling the "source subject sibling" phrase. Fold in M5a r1's survivors (`design_docs/planned/w-marketing-capture-sprint-plan.md`). DECISIONS FOR MARK: D-WORLD-86 (OPEN; default hold). It now also gates row 156's close.

## 298 — 2026-10-22 — row 150 M5b built: TESTED/PROVEN roster rows and the four-row packet evidence render, landing iteration 296's orphan (PR #353; judge r1 PASS 86 → r2 PASS 90) [PRODUCT]

**Kind:** one milestone of row 150 (M5b). Production: `host/daemon/workbench_packet_evidence.go`, `host/daemon/workbench_packet.go`, `host/workbench/render.go`. Tests, plus two tool-generated fixtures in `host/daemon/testdata/`. No `.ail` change, no release, no design-direction change.

**Orphan credit (iterations 296 and 297, no log entries of their own):**
- 296 built M5b, then was STALL-killed at gate 3 while its foreground judge ran. The driver's notice: https://github.com/sunholo-data/ailang-world/issues/313#issuecomment-6321634905 (rc=143, 2026-10-21T16:56:25Z).
- 297 rebased the branch onto `27c60f5f` and spawned a fresh judge. It was killed the same way; the driver suppressed that notice as "failure unchanged this episode".
- **Mechanism, measured from both transcripts and `tools/launchd/mission-control.sh` `_mc_stalled`:**
  - Both judges ran in the foreground, so the controller's own transcript was flat.
  - Sub-agent transcripts (`<session>/subagents/*.jsonl`) are counted by no arm.
  - `_mc_progress_bytes` reads the NEWEST `*.jsonl` in the project dir, which every session with this cwd shares. 297's kill line `prog=318821268` is a different 317 MB session, not the controller's 1.07 MB file.
  - Both judges launched their mutant batches with `( … ) &`, so the runners left the descendant tree. The judges then sat in sleep-polls: 296 a 595 s poll on `/tmp/ev296m/ledger.txt`, 297 a 235 s poll on `/tmp/ev297m/b?.log`. Each tool result is "Exit code 144" at the kill second.
  - 297's four batches had all written DONE; its report was never written.
- Filed `[HARNESS] ticket:stall:watchdog:foreground-agent-and-wrong-transcript` (mission-fleet, blocking item).
- **Workaround that held this iteration:**
  - The judge ran in the BACKGROUND while the controller kept the turn alive with bounded in-turn polls on its incremental report file.
  - The directive forbade `&`/nohup/sleep-polls and capped each Bash call at 290 s.
  - Both rounds completed: r1 26.4 min, r2 5.3 min.

**Pick and why:**
- All seven clauses are MET, and D-WORLD-81's attended order governs: 156 → 150 → 151 → v1.0.0.
- Row 156 has nothing buildable (D-WORLD-86 OPEN), so the next item is row 150 M5b. Its work already existed on the unpushed branch `sprint/world-iter296-row150-m5b`.
- Reality check:
  - No open PR of this loop's.
  - The branch was rebased onto `e8631dd8`: 2 commits, and zero file overlap with the 4 intervening merges (row 188 surface, `host/surface/*`, `cmd/ailang-worldd/surface_index.go`).
  - The attended rulings D-WORLD-88 and D-WORLD-89 since 295 do not touch row 150's order.

**Clause map:** 1–7 all MET.

**Gate 0:**
- Kill switch armed; gh `sunholo-voight-kampff`; billing CLEAN.
- #313: 0 allowlisted directives since `2026-10-18T23:25:31Z` (16 comments).
- Self-notice scan rc 1: the 296 crash, credited above (control 107:4 ok).
- `mission-world` inbox: 6 unread, all this loop's own claims.

**Gates 1–2:**
- Local dev = origin/dev `e8631dd8`.
- The running skill matches origin on SKILL.md and all 13 resource files (resolved symlink). My first loop omitted `resources/` from the path and printed 13 DIFFs; the corrected re-run printed 13 same — the iteration-290 bare-filename trap again.
- Ledger valid, 76 rows; one OPEN (D-WORLD-86).

**Controller verification (outside any sandbox, `AILANG_BIN` pinned, `WORLD_TOOL_AILANG_BIN` unset) on `66208c2a`:**
- The 9 named M5b/neighbour tests, `go test -race ./host/workbench`, gofmt, vet and `verify_ail.sh`: all rc 0. The roster test ran 70 PASS lines in 45 s.
- `verify_go.sh` rc 1 (929 s, load 5–12, a sibling's godot/perl/rig `ailang`). The only red was `host/verifygate` in the race leg, at the 8 m per-package timeout (`TestSeToolsGateArmsAreLoadBearing/plan-bytes` running).
  - The package is byte-identical to base, and its non-race leg was ok (378 s).
  - An isolated `go test -race ./host/verifygate -timeout 1500s` was rc 0 (343 s).

**Independent evaluation (Agent sonnet, background, own worktree `.eval-world-iter298`; generator codex ≠ judge sonnet; the controller's test-only commit and fix commit are same-vendor as the judge, FLAGGED, and the directive told it to scrutinise them equally):**
- **r1 PASS 86**, zero blocking; 74,842 subagent tok.
  - All gates rc 0, and the three plan mutants were killed at the named assertions.
  - Own drill of about 35 mutants. Fixtures regenerated with the v0.54.15 tool binary: structurally equal.
  - Rulings: the PROVEN-subject resolution is SOUND (the MAC covers the subject; reading it before ValidateProof only buys a bounded, typed fetch). The executor's three choices are acceptable. The `utf8.Valid` removal is equivalent (json.Unmarshal yields valid UTF-8).
  - Survivors:
    - (1) the TESTED producer could be the effect's provenance;
    - (2) no assertion on the refused-roster banner;
    - (3) a PROVEN refusal reason was never checked through `projectPacket`;
    - (4) an UNAVAILABLE proof row showed a subject read before the MAC check;
    - R8, the compiler anchor.
- **Fix `f9503386`** (controller):
  - `row.subject` is set only in the PROVEN success assignment (a 2-line production move).
  - A `result-provenance` control (result payload plus newline, distinct provenance, with an instrument check).
  - Proof refusals assert no subject/compiler/version, and store-real refusals check `view.Roster[3].Unavailable`.
  - Render: the escaped `roster-unavailable` banner and the exact subject/compiler anchors.
  - Controller mutants S1–S5 were each killed at the new assertion, with byte-identical restores.
- `verify_go.sh` on `f9503386`: rc 1 (986 s, load 13–17). Race-leg reds only: `host/verifygate` (8 m timeout) and `host/daemon` `TestAGUISlowReaderCutIsResumable` ("unexpected EOF").
  - That test is known to be timing-sensitive under `-race` (archive: the iteration on `d2d1815`), and no agui file is touched.
  - The non-race leg was all ok. Isolated `-race` re-runs at 1500 s: daemon rc 0 (362 s), verifygate rc 0 (394 s).
- **r2 PASS 90**, zero blocking, no new findings; 85,026 subagent tok. Its T1–T7 were all killed.
  - Findings 1–4 and R8 are CLOSED. Findings 5 (every link labelled `report:`) and 6 (dead utf8 guard) stay open by choice.
  - It judged the trailing-newline control synthetic but fair: real results pass through the same hash/kind/size read and strict decode.

**Result:** M5b built. Row 150 has 10 of 17 milestones. Plan and sprint JSON mark M5b completed.

**Gate 3b:** PR #353 (`sprint/world-iter296-row150-m5b`); the record rides the same PR. CI is verified on the PR head and the merge before Gate 5 reports LANDED.

**Routing evidence:** base=e8631dd8a80fd2b614cd6fd64c0a48b81ac3e98d@2026-10-22T01:04:57Z (Gate 4).
- Controller `claude:claude-opus-5-5` (tok: not reported).
- Designer not spawned: the design is ready. Planner not spawned: the plan is ready.
- Executor `recipe codex:gpt-6.1-sol declared:provider-pin`: one run in iteration 296, rc 0, 206,730 tok, about 20 min. Not re-run here; this iteration landed its output.
- Evaluator `agent-tool sonnet declared:alias-pin`: r1 74,842 and r2 85,026 subagent tok. r1 ran in the background; r2 resumed the same agent with `SendMessage`.
- The operator's standing request asked for Agent-tool roles "exactly as the routing table specifies". The evaluator ran through the Agent tool as required. The executor's table lane is the codex recipe, and its output was reused rather than re-spawned.
- Judge independence: cross-vendor for the executor's work; same-vendor for the controller's two commits, mitigated by the judge's own T-mutants on them.
- Metered **$0.00**. Every required role for this pick ran. The two lost judge runs (296, 297) are banked as slot losses on the ticket.

**Ruled out:**
- Reading either `verify_go.sh` red as this diff's: both red packages passed isolated re-runs, and agui/verifygate are untouched.
- Re-running the executor: its output was intact on the branch, and the judge measured it.
- The judge's own "loop printed a blank rc" note (r1): it used the `--- PASS`/`ok` lines instead, as did the controller.

**Carried (non-blocking):**
- M5b findings 5 and 6, as above.
- M5a r1's remaining survivors go to M5c: a non-deadline GetObject error, an empty candidate, the zero-ResultRef accept path, and a hand-typed roster not tied to the seed producer's bytes.
- M4b r2's three surviving judge mutants, for when M5c touches the validator.
- **[HARNESS] ticket:stall:watchdog:foreground-agent-and-wrong-transcript** (mission-fleet, blocking item, filed this iteration).
- **[HARNESS] ticket:routing-models-check:bold-colon-heading**: still open from 295.

**Retro:**
- **Skill lane:** none applied; World does not edit the shared skill. PROPOSE to the fleet: standing rule 7's prescription "prefer a FOREGROUND spawn" is what the STALL watchdog kills when a foreground sub-agent's own work leaves the descendant tree. Until the ticket is fixed, a long judge is safer in the background, with in-turn controller polls.
- **Process lane (judge directives):** every judge directive forbids `&`/nohup/setsid/sleep-polls, caps each Bash call at ≤290 s, and requires an incremental report file. A killed judge's partial ledger is then recoverable, as `/tmp/ev297m` was but no report was.
- **Backlog lane:** none new.
- The last three landings (298, 295, 294) moved no UNMET clause. All clauses are MET, so there is no drift alarm.

**Progress**: World 1.0 **7/7 clauses MET**; **goal unmoved**. Row 150 **10/17 milestones built**; next M5c. Release `v1.0.0-rc.1`; 156 (held on D-WORLD-86) → 150 → 151 before v1.0.0.

**Next:** row 150 M5c: seed roster validation and success manifest publication. Fold in the remaining M5a r1 survivors and M4b r2's validator mutants (`design_docs/planned/w-marketing-capture-sprint-plan.md`). DECISIONS FOR MARK: D-WORLD-86 (OPEN; default hold). It also gates row 156's close.

## 299 — 2026-10-22 — row 150 M5c built: the demo seed completes with a final roster check and one atomic success manifest; M5b's TESTED rule fixed to the real tool shape (PR #354; judge r1 PASS 80 → r2 PASS 85) [PRODUCT]

**Kind:** one milestone of row 150 (M5c). Production: `cmd/ailang-worldd/demo_seed.go`, `cmd/ailang-worldd/demo.go`, `host/daemon/workbench_packet_evidence.go`. Tests, one tool-generated fixture (`host/daemon/testdata/roster_cli_test_result_real.json`), QUICKSTART §12 and `website/docs/reference/cli.md`. No `.ail` change, no release, no design-direction change.

**Pick and why:**
- All seven clauses are MET; D-WORLD-81's attended order governs: 156 → 150 → 151 → v1.0.0. Row 156 has nothing buildable (D-WORLD-86 OPEN), so row 150's next milestone, M5c, was the pick (iteration 298's **Next**).
- Reality check: no open PR of this loop's; no M5c commit on origin (`git log origin/dev --grep M5c` = row 156's M5c only); no iteration-299 worktree or branch residue. The 2026-10-21T16:56:25Z rc=143 notice is iteration 296's, already credited by 298.

**Clause map:** 1–7 all MET.

**Gate 0:**
- Kill switch armed; gh `sunholo-voight-kampff`; billing CLEAN.
- #313: 0 allowlisted directives since `2026-10-22T01:34:42Z` (18 comments). Self-notice scan rc 1 on the already-credited 296 notice (control 107:4 ok, watermarks 2/2).
- `mission-world` inbox: 7 unread, all this loop's own claims.

**Gates 1–2:**
- Local dev = origin/dev `f4b2dfdb`; dev CI 3/3 checks green on it.
- The running skill matches origin on SKILL.md and all 13 resource files (resolved symlink). My first loop again omitted `resources/` and printed 13 DIFFs; the corrected loop printed 13 SAME (third time this trap fired: 290, 298, 299).
- Ledger valid, 76 rows; one OPEN (D-WORLD-86).

**Executor (codex gpt-6.1-sol, one run, 162,686 tok, ~14 min):**
- `demoFinalizeSeed`: re-validates the manifest, then `demoValidateSeedRoster` reads the S5 input actually committed and, per slot, requires slot order, ref = manifest ref, an object that re-hashes to its ref, and the declared semantic kind (`demoValidateRosterKind`, independent of the projection); then fetches the packet page and requires four resolved grade markers and no roster-unavailable banner (`demoValidateRosterPage`).
- `demoPublishManifest`: `CreateTemp` in the target dir, chmod 0644, write, fsync, `os.Link` (refuses an existing target), temp removed. `runDemo` returns exit 0 only after it.
- Help, QUICKSTART §12 and the CLI reference now describe a complete seed.
- Carried survivors: M4b r2's stray extra ref, why-hash count and duplicate why entry got pin-free controls; M5a r1's non-deadline read error, empty candidate and resultless effect got projection rows. The last was a **real defect**: a denied/failed recorded effect (no result) rendered ATTESTED; it now refuses `effect result absent`.
- In-sandbox: 6 gates PASS, 5 socket-bound tests UNINFORMATIVE; 6 mutants killed, 5 UNMEASURED-SANDBOX.

**Controller verification (outside any sandbox, both pins):**
- Base `f4b2dfdb`, in the sprint worktree and a separate base worktree: every directive gate green; `go test ./cmd/ailang-worldd` rc 0 (136 s).
- On the executor's tree: four real-stack tests RED — `roster page: evidence roster UNAVAILABLE`. A debug patch (restored, `cmp`) showed the reason `test exit code`, and the stored result carried no `exit_code` key.
- **Measured the pinned tool directly** (`policy-tool --policy <rendered>` from a scratch root): passing `test --json` → `{argv, ok:true, stderr, stdout}`, no `exit_code`; failing → `exit_code:1, ok:false`. M5b's `packetTestedRow` required `exit_code == 0`, and M5b's tests used a fake tool that always emitted it — so no real TESTED row could ever resolve. Only M5c's real end-to-end seed exercises that path.
- **Fix `c9fd10af`** (controller): `exit_code`, when present, must be 0; `ok` must be true. Rows `real-tool-success-shape` (fixture = the tool's verbatim stdout, checked in; instrument check that it has no `exit_code`), `exit-code-absent-not-ok`, `exit-code-null`. Mutants C1 (require exit_code again), C2 (accept null), C3 (drop ok) each killed at the new row.
- Executor's UNMEASURED-SANDBOX mutants run here: kind-finalizer (first form build-red, compiling form killed at `demo_test.go:2486`), early-write (killed `:2489`), extra ref (build-red, compiling form killed), why count (killed), duplicate entry (killed at the named assertion; another guard also refuses, so it is defence in depth).
- After the fix: all 9 named gates top-level PASS, no SKIP. `verify_ail.sh` rc 0. `verify_go.sh` rc 0 (771 s, zero FAIL lines, load 3–7) on `c9fd10af`.
- Commits: `137b615d` (executor; reconstructed by removing the controller hunks, all 8 final files SHA-matched after), `c9fd10af` (controller fix).

**Independent evaluation (Agent sonnet, background, own worktree `.eval-world-iter299`; generator codex ≠ judge sonnet; the controller's two commits are same-vendor with the judge, FLAGGED, and the directive told it to scrutinise them equally):**
- **r1 PASS 80**, zero blocking; 79,345 subagent tok; ~10.4 min. All gates PASS. It re-measured the tool shape first-party (pass, fail, and a nonexistent path → `{ok:false, refused}`) and ruled the controller fix CORRECT, with nothing else fitted to the fake tool. 11 of its mutants killed.
  - Survivors: M4 (the page check's production wiring), M5 (manifest-vs-stored-roster compare), M6 (slot order), M8 (re-hash), M9 (count `<4`), M10 (file mode), M11 (fsync), M13 (redundant validation). Notes: public-bytes scan runs on compact JSON while the file is indented; scope beyond the plan's file list (directive-allowed).
- **Fix `be270368`** (controller, test-only): real-store arms committing a genuine extra approval through the same coordinator — TESTED → S3's `ai_check` record (right kind; only the page refuses, `test op`), and the first two slots swapped (package law and projection are order-free; only the seed's slot order refuses, with an instrument check that the page still resolves); `TestDemoSeedEndToEnd` arms for a manifest naming the goal as CLAIMED and for forged CLAIMED bytes through `tamperProxy`; a 0644 mode assertion. Controller mutants F4/F5/F6/F8/F10 each killed at the new assertion with `<nil>` as the error (the guard alone catches it).
- `go test ./cmd/ailang-worldd` rc 0 with the tool pin unset (CI configuration, 55 s); `go test -race ./cmd/ailang-worldd` rc 0 with both pins (205 s).
- **r2 PASS 85**, zero blocking, no new findings; 91,250 subagent tok. Its own M4/M5/M6/M8/M10 all killed; it checked each new arm reaches the guard it claims and that the extra-approval helper uses the genuine coordinator path. M9/M11/M13 and the two notes stay open, accepted.

**Result:** M5c built. Row 150 has 11 of 17 milestones. Plan and sprint JSON mark M5c completed.

**Gate 3b:** PR #354 (`sprint/world-iter299-row150-m5c`); the record rides the same PR. CI is verified on the PR head and the merge before Gate 5 reports LANDED.

**Routing evidence:** base=f4b2dfdbcf1b2e342fce961dd4740a8bf3e0f37a@2026-10-22T03:27:47Z (Gate 3 worktree).
- Controller `claude:claude-opus-5-5` (tok: not reported).
- Designer not spawned: the design is ready. Planner not spawned: the plan is ready.
- Executor `recipe codex:gpt-6.1-sol declared:provider-pin`: one run, rc 0, 162,686 tok. Probe rc 0.
- Evaluator `agent-tool sonnet declared:alias-pin`: r1 79,345 and r2 91,250 subagent tok. r1 ran in the background; r2 resumed the same agent with `SendMessage`.
- The operator's standing request asked for Agent-tool roles exactly as the routing table specifies. The evaluator ran through the Agent tool as required; the executor's table lane is the codex recipe. No role was unspawnable.
- Judge independence: cross-vendor for the executor's work; same-vendor for the controller's `c9fd10af` and `be270368`, mitigated by the judge's own first-party tool measurement and its own mutants on both.
- Metered **$0.00**.

**Ruled out:**
- Reading the four real-stack reds as an M5c defect: the reason was M5b's TESTED rule, reproduced against the pinned tool directly.
- Treating the build-red K1/K3 mutants as kills: re-run as compiling mutants.
- Making a >4-slot roster reach the validator for M9: the package law refuses it first.

**Carried (non-blocking):**
- M5c: M9, M11, M13; the compact-vs-indented public-bytes scan.
- M5b: every roster link labelled `report:`; the dead `utf8.Valid` guard.
- **[HARNESS] ticket:stall:watchdog:foreground-agent-and-wrong-transcript** and **ticket:routing-models-check:bold-colon-heading**: still open.

**Retro:**
- **Skill lane:** none applied; World does not edit the shared skill. Gate 3's "hand-typed fixture" rule (V1 iteration 348) gained a World instance at a NEW boundary: not a parser's fixture but a FAKE TOOL's output shape, two milestones upstream of the test that could see it. Proposal to the fleet: extend that rule to fake binaries — a fake tool's response must be generated from the real tool's output, and the first milestone that runs the real tool end-to-end re-validates every fake-derived rule. One instance; recorded so a second is recognisable.
- **Process lane:** a World note in the Repo Profile: a fake-tool test row must carry a sibling row fed the real tool's captured response (generated, checked in, with the command in a comment).
- **Backlog lane:** none new.
- Harness share 0 of the last 20. The last three landings (299, 298, 295) moved no UNMET clause; all clauses are MET, so no drift alarm.

**Progress**: World 1.0 **7/7 clauses MET**; **goal unmoved**. Row 150 **11/17 milestones built**; next M6. Release `v1.0.0-rc.1`; 156 (held on D-WORLD-86) → 150 → 151 before v1.0.0.

**Next:** row 150 M6: still serve flag and stable complete server HTML (`design_docs/planned/w-marketing-capture-sprint-plan.md`). DECISIONS FOR MARK: D-WORLD-86 (OPEN; default hold). It also gates row 156's close.

## 300 — 2026-10-22 — row 150 M6 built: `serve --workbench-still` serves the complete workbench page without its one live script (PR #359; judge r1 PASS 93) [PRODUCT]

**Kind:** one milestone of row 150 (M6). Production: `cmd/ailang-worldd/main.go`, `host/daemon/daemon.go`, `host/daemon/workbench.go`, `host/workbench/render.go`. Tests: `host/daemon/workbench_still_test.go` (new), `cmd/ailang-worldd/main_test.go`, `host/workbench/render_test.go`. Docs: QUICKSTART §11, `website/docs/guides/operating-the-daemon.md`, `website/docs/reference/cli.md`. No `.ail` change, no release, no design-direction change.

**Pick and why:**
- All seven clauses are MET; D-WORLD-81's attended order governs: 156 → 150 → 151 → v1.0.0. Row 156 has nothing buildable (D-WORLD-86 OPEN), so row 150's next milestone, M6, was the pick (iteration 299's **Next**).
- Reality check: `--workbench-still`/`WorkbenchStill`/the three new test names had zero hits in `cmd/` and `host/` at base; no M6 commit on origin; no iteration-300 worktree, branch or PR. Open PRs #356/#357/#358 are attended work (D-WORLD-90, row 188 trial), not this loop's, and were left alone. The 2026-10-21T16:56:25Z rc=143 notice is iteration 296's, already credited (log count 1).
- Quorum-at-pick: not re-run — the design carries its r1–r4 record from iteration 266 and the plan has been executing since.

**Clause map:** 1–7 all MET.

**Gate 0:**
- Kill switch armed; gh `sunholo-voight-kampff`; billing CLEAN.
- #313: 0 allowlisted directives since `2026-10-22T01:34:42Z` (19 comments). Self-notice scan: one rc=143 notice in window (iteration 296's, already credited); control 107:4 ok, watermarks 2/2.
- `mission-world` inbox: 8 unread, all this loop's own claims.

**Gates 1–2:**
- Local dev = origin/dev `9aca8812`; dev CI 3/3 checks green on it.
- The running skill matches origin on SKILL.md and all 13 resource files (resolved symlink, per-file).
- Ledger valid, 76 rows; one OPEN (D-WORLD-86).

**Executor (codex gpt-6.1-sol, one run, 103,331 tok, ~6.6 min):**
- `--workbench-still` (default false) → `daemon.Config.WorkbenchStill` → `workbench.Page.Still`; the template emits the `live.js` tag inside `{{if .Still}}{{else}}…{{end}}` (the field census reads a direct bool branch; `not` is not in the template's FuncMap).
- Tests first: `TestWorkbenchStillAllRegions` (home, entry and packet-object routes on a store with a produced packet and a recorded deny: timeline entry, live-list head, graph edges, decisions row, packet refs, inspector; zero `<script`; full headers `DeepEqual` to default; still HTML equals default minus exactly the tag; repeat GETs after a real 2 s sleep byte-identical; selected head and log head unchanged), `TestServeWorkbenchStillFlag` (help, parser, Config, default one script, `=notabool` refused), `TestWorkbenchStillDefaultLiveScript` (one/zero scripts, census).
- In-sandbox: all gates PASS, no socket denial; 5 mutants killed (plan M6-script/graph/default plus threading and census).

**Controller verification (outside any sandbox, both pins):**
- Base `9aca8812` baseline: the six pre-existing focused gates PASS; the new test name selects nothing (negative control). My first baseline loop used zsh `set -- $spec`, which does not word-split, so every gate read rc 1 — an instrument failure; re-run under bash.
- On `aeaa957d`: 12 focused gates top-level PASS (the plan's 9 + `TestQuickstartSection11FlagsMatchTheCLI`, `TestWorkbenchErrorPageInert`, `TestWorkbenchFieldCensusMutations`); gofmt and vet clean.
- `verify_ail.sh` rc 0. `verify_go.sh` rc 1 (935 s): non-race leg all ok; race leg red while the judge ran concurrently (load 12.8) — `cmd/ailang-worldd` and `host/daemon` `signal: terminated` at ~218 s, and `host/broker` `TestWorktreeSnapshotRefusesAConfiguredFilter` failing its own control arm (`plain git status did not run the clean filter`). M6 touches neither broker nor git filters. Isolated `go test -race -count=1` of the three packages afterwards: rc 0 (cmd 192 s, daemon 417 s, broker 267 s).
- Diff read for out-of-worktree writes: none.
- Commits: `aeaa957d` (executor work, controller-committed), `ac1abe7f` (controller docs fix), `4f935c8a` (plan/JSON status).

**Independent evaluation (Agent sonnet, background, own worktree `.eval-world-iter300`; generator codex ≠ judge sonnet):**
- **r1 PASS 93**, zero blocking; 51,395 subagent tok, 9 tool calls, ~13.8 min. Gates PASS; full `./host/daemon` and `./cmd/ailang-worldd` UNMEASURED by it (600 s foreground cap; it ran `-run` subsets instead — the controller's isolated runs cover them).
- 10 mutants killed: plan M6-script (both packages), M6-graph, M6-default; Page threading; always-suppress; inverted condition; flag mapped to `&& false`; help text mangled; flag renamed. It re-made two non-compiling forms as compiling ones.
- Non-blocking: N1 the cmd test builds the Page itself (the cfg→daemon→Page link is killed only in the daemon test); N2 graph assertion skips the object route by index; **N3 the CLI reference's serve-flag help copy lacked the flag — reproduced (`cli.md` carries a `serve flags:` copy of the help text) and fixed `ac1abe7f`, docs only, no r2**; N4 the plan lists a mutation command twice; N5 QUICKSTART wording.

**Result:** M6 built. Row 150 has 12 of 17 milestones. Plan and sprint JSON mark M6 completed.

**Gate 3b:** PR #359 (`sprint/world-iter300-row150-m6`); the record rides the same PR. CI is verified on the PR head and the merge before Gate 5 reports LANDED.

**Routing evidence:** base=9aca881297dcc9f84a2809bacffb95878e93e4e6@2026-10-22T08:00:44Z (Gate 4; Gate 1 and Gate 3 read the same SHA).
- Controller `claude:claude-opus-5-5` (tok: not reported).
- Designer not spawned: the design is ready. Planner not spawned: the plan is ready.
- Executor `recipe codex:gpt-6.1-sol declared:provider-pin`: probe rc 0; one run, rc 0, 103,331 tok.
- Evaluator `agent-tool sonnet declared:alias-pin`: r1 51,395 subagent tok, background, controller polled in-turn.
- The operator's standing request asked for Agent-tool roles exactly as the routing table specifies. The evaluator ran through the Agent tool as required; the executor's table lane is the codex recipe (an Agent spawn for a provider-pinned role is denied by the spawn-pin hook). No role was unspawnable.
- Judge independence: cross-vendor for the executor's work; the controller's `ac1abe7f` is docs only.
- Metered **$0.00**.

**Ruled out:**
- Reading the race-leg reds as an M6 regression: two were external terminations at the same instant, the third a control-arm failure in an untouched package; all three packages are green in isolation.
- Treating my zsh baseline loop's rc 1 as a red base: the loop never selected a test.

**Carried (non-blocking):**
- M6: N1, N2 above.
- M5c: M9, M11, M13; the compact-vs-indented public-bytes scan. M5b: `report:` link labels; the dead `utf8.Valid` guard.
- **[HARNESS] ticket:stall:watchdog:foreground-agent-and-wrong-transcript** and **ticket:routing-models-check:bold-colon-heading**: still open.

**Retro:**
- **Skill lane:** none; World does not edit the shared skill.
- **Process lane:** a judge running its own full-package suites concurrently with the controller's `verify_go.sh` is the second instance of "sequence the judge and the suite" (iteration 293 recorded the first). Remedy recorded here, not yet in the Repo Profile: start `verify_go.sh` only after the judge returns, or give the judge `-run` subsets by directive. Two instances; a third makes it a charter guardrail.
- **Backlog lane:** none new.
- Harness share 0 of the last 20. The last three landings (300, 299, 298) moved no UNMET clause; all clauses are MET, so no drift alarm.

**Progress**: World 1.0 **7/7 clauses MET**; **goal unmoved**. Row 150 **12/17 milestones built**; next M7a. Release `v1.0.0-rc.1`; 156 (held on D-WORLD-86) → 150 → 151 before v1.0.0.

**Next:** row 150 M7a: bounded capture supervisor with mandatory fake Chrome companions (`design_docs/planned/w-marketing-capture-sprint-plan.md`). DECISIONS FOR MARK: D-WORLD-86 (OPEN; default hold). It also gates row 156's close.

## 301 — 2026-10-22 — row 150 M7a built: the bounded capture supervisor and its fake-Chrome/fake-daemon companions (PR #360; judge r1 PASS 86 → r2 PASS 91 → r3 PASS 91) [PRODUCT]

**Kind:** one milestone of row 150 (M7a). Production: `scripts/marketing_capture.py` (new), `scripts/capture_marketing.sh` (new). Tests: `scripts/marketingcheck/supervisor_test.go` (new, test-only package; M7b adds `main.go`). No `.ail`, no `go.mod`/`go.sum`, no `host/` change, no release, no design-direction change.

**Pick and why:**
- All seven clauses are MET; D-WORLD-81's attended order governs: 156 → 150 → 151 → v1.0.0. Row 156 has nothing buildable (D-WORLD-86 still OPEN, re-read at Gate 0), so row 150's next milestone, M7a, was the pick (iteration 300's **Next**).
- Reality check: no `marketing_capture`/`capture_marketing`/`marketingcheck` file at base, the four test names appear only in the plan and sprint JSON, no M7a commit on origin, no iteration-301 branch or worktree. No open PRs. The attended merges since iteration 300 (#356/#357/#358: D-WORLD-90 font allowance, row 188 N1) touch none of M7a's files.
- Quorum-at-pick: not re-run — the design carries its r1–r4 record from iteration 266 and the plan has been executing since.

**Clause map:** 1–7 all MET.

**Gate 0:**
- Kill switch armed; gh `sunholo-voight-kampff`; billing CLEAN.
- #313: 0 allowlisted directives since `2026-10-22T01:34:42Z` (20 comments). Self-notice scan: one rc=143 notice in window (iteration 296's, already credited, log count 1); control 107:4 ok, watermarks 2/2.
- Ledger valid, 77 rows; one OPEN (D-WORLD-86). D-WORLD-90 is new since iteration 300 and RESOLVED by an attended ruling: provenance `attended` (flip `b948cf3b`, creation `3b09253e`). It governs row 188, not this pick.
- `mission-world` inbox: only this loop's own claims.

**Gates 1–2:**
- Local dev = origin/dev `cfa6d583`; dev CI 3/3 checks green on it (the host build+test check was still pending at the first read and green at the second).
- The running skill matches origin on SKILL.md and all 13 resource files (resolved symlink, per-file).

**Executor (codex gpt-6.1-sol, one run, 102,443 tok, ~9.4 min):**
- `marketing_capture.py`: `shot` (one Chrome screenshot: the revision-notes argv with per-shot `preferredColorScheme=1|0`, a fresh absolute profile, stdin `/dev/null`, HOME inherited, own session; success = rc 0 AND a new regular nonempty PNG, any existing destination refused before launch; on the 30 s deadline a group SIGKILL, reap ≤5 s, `Timeout` even with a PNG) and `daemon` (file readiness ≤30 s, TERM → ≤5 s → group KILL → reap ≤5 s; a cleanup failure overrides shot success; hidden `--unit-cleanup-failure` seam injected after real cleanup). One JSON verdict on stdout; exit 0 only for `ok`.
- `capture_marketing.sh`: `--help` documents the full M7c configuration; EXIT/INT/TERM/HUP traps clean up a supervised daemon; a real capture exits 1 ("six-shot pass is not implemented yet (M7c)"); `--unit-daemon exit|wait` unit control.
- Four tests drive the real Python/shell with `/bin/sh` fakes in `t.TempDir()`; 10 mutants killed in-sandbox, no sandbox denials. C11 observation before/after byte-identical.

**Controller verification (outside any sandbox):**
- Diff read for out-of-worktree writes: none (fakes and outputs under `t.TempDir()`).
- On `38ac4167`: the four focused gates top-level PASS; `bash -n`, vet, gofmt clean; `--help` names all ten flags; `-race -count=3` rc 0; a full capture refuses rc 1. Plan mutants re-run by me: M7a-no-file, -child-only, -timeout-success each red at the named assertion, restore sha OK. A stray `sleep 60` in the process table was traced (ppid) to another session's PR poll, not a leak.
- `verify_ail.sh` rc 0 (51 s).
- `verify_go.sh`, three runs, each started only after the judge had returned:
  - **`98f48b20` rc 1 (487 s).** The non-race `go test ./...` failed five NEW M7a subtests on timing alone, at load ~6. A 0.3 s startup deadline lost to a slow-starting daemon, so early-exit and exits-while-ready read `Timeout`. The other three blew their bounds: escalation 4.07 s against 2 s, term-first 2.79 s against 1.5 s, and the interrupted-shot readiness wait (4 s). This is the judge's r2 risk R2-3, now measured. `TestSubprocessSinkGate` also flagged my interrupted-shot subtest's bare `bytes.Buffer` stdout sink, a real defect: it can race os/exec's copier. Fixed in `dd647909` (wider upper bounds, a per-case startup deadline, a file sink).
  - **`dd647909` rc 1 (468 s).** The daemon subtests held. `TimeoutKillsGroup` found no `shot.png` at its 3 s deadline, because the fake shell had not yet written it. Fixed in `c70b58b9` with an 8 s deadline.
  - **`c70b58b9` rc 1 (970 s).** The non-race leg is **all ok**, including `marketingcheck`. The race leg was red under load in three packages M7a does not touch: `host/daemon` hit the 8 m package timeout; `host/projection` failed `TestMCPToolsInnerBudget` on a 30 ms deadline; `host/verifygate` failed `TestMissionConfigDefaultFlowIntegration/valid_ledger` on its 10 s fixture bound. `marketingcheck` was ok in the race leg. Isolated `-race` afterwards was rc 0 for all three: daemon 409 s, projection ok, verifygate ok. The verifygate test also passed isolated in the non-race mode (8 s); it has no prior mention in the log.
- After each widening I re-applied the timing-sensitive mutants. All 11 stayed red at their named assertions: the three plan mutants, N2, S1, X18, X19, X21, X24, X26 and the TERM trap.
- Commits: `38ac4167` (executor work, controller-committed), `98f48b20` (controller r1 fix round), `dd647909` and `c70b58b9` (controller load-robust test bounds and the sink fix, test-only), `cde96472` (plan/JSON status).

**Independent evaluation (Agent sonnet, background, own worktree `.eval-world-iter301`; generator codex ≠ judge sonnet):**
- **r1 PASS 86**, zero blocking (79,502 subagent tok, ~4.8 min). Gates and the three plan mutants confirmed; 26 own mutants from the diff, 16 survived. Findings: **N1** untested branches (shell INT/HUP traps, daemon early exit / readiness timeout / exit during lifecycle, TERM-before-KILL, argument validation, marker freshness, dangling-symlink destination; KILL-reap CleanupFailed unreachable); **N2** a real defect — in `daemon --chrome`, TERM only set the stop event while `shot()` sat in `communicate()`, so the supervisor lived **29 s** after TERM (measured by the judge); **N3** the lifetime mutant reds only via the 12 s context; **N4** the shell's `$SECONDS` bound is 4–5 s.
- **Controller fix `98f48b20`** (production + tests): `shot()` waits in 50 ms slices and, when signalled, kills the Chrome group and reports `Interrupted`; subtests for every N1 survivor plus `interrupted-shot`; comments on the two unreachable branches. Each targeted survivor was re-applied on the fixed tree and went red at its named assertion (N2-no-stop, S1, S2, X4, X18, X19, X21, X24, X25, X26); X25b reds through a Python traceback (the defect's real shape); X26 first reddened only via the context kill, so the subtest got short waits and then reddened at its assertion. Plan mutants re-measured on the fixed tree: all three red, restores OK; `-race -count=3` rc 0.
- **r2 PASS 91**, zero blocking (resumed judge, 99,375 subagent tok, ~5 min). N2 CLOSED (0.63 s vs 29 s); N1 CLOSED for S1/S2/X4/X18/X19/X21/X24/X25/X26; 8 new-code mutants, 7 red. Still open, non-blocking: X23 (Python HUP handler untested — HUP reaches only the shell), X13/X14 redundant pair, S6/S7 (shell cleanup-status propagation, shell readiness bound), N3, N4, R2-1 (a 1 s deadline slack in `wait_bounded` survives), R2-3 (the new sub-second/2 s bounds are the first to flake on a loaded Linux `-race` runner).

- **r3 PASS 91**, zero blocking (resumed judge, 113,672 subagent tok, ~5.7 min), on the test-only delta `98f48b20..c70b58b9`. The diff holds exactly the listed bound, deadline and sink changes. 12 mutants re-applied: 10 red at named assertions, and X13/X14 are the known equivalent pair. No widening let a mutant survive: escalation's 8 s bound cannot hide a no-KILL mutant (X28 reds at the dead-PID check), and term-first still separates TERM from KILL-only (X24 measured 20.4 s against its 10 s bound). R2-3 is CLOSED. New, both non-blocking: R3-1, a TERM-path regression now takes 20–30 s to surface (still named); R3-2, under the S1 mutant the test cleanup SIGKILLs the supervisor and orphans the daemon (fix: kill the group; unreachable on the shipped tree). The package takes ~17 s with or without `-race`.

**Result:** M7a built. Row 150 has 13 of 17 milestones. Plan and sprint JSON mark M7a completed.

**Gate 3b:** PR #360 (`sprint/world-iter301-row150-m7a`); the record rides the same PR. CI is verified on the PR head and the merge before Gate 5 reports LANDED.

**Routing evidence:** base=cfa6d583ff769cd61b08be9c55c1d615db76753f@2026-10-22T11:26:41Z (Gate 1; Gate 3 re-read the same SHA at 11:29:56Z).
- Controller `claude:claude-opus-5-5` (tok: not reported).
- Designer not spawned: the design is ready. Planner not spawned: the plan is ready.
- Executor `recipe codex:gpt-6.1-sol declared:provider-pin`: probe rc 0; one run, rc 0, 102,443 tok.
- Evaluator `agent-tool sonnet declared:alias-pin`: r1 79,502 / r2 99,375 / r3 113,672 subagent tok (as reported per round; one resumed agent), background, controller polled in-turn (the harness refuses a bare foreground `sleep`; the poll was a bounded `until` on the transcript's size/mtime, read by `stat` only).
- The operator's standing request asked for Agent-tool roles exactly as the routing table specifies. The evaluator ran through the Agent tool as required; the executor's table lane is the codex recipe (an Agent spawn for a provider-pinned role is denied by the spawn-pin hook). No role was unspawnable.
- Judge independence: cross-vendor for the executor's work. The controller's `98f48b20` (production + tests) was judged in r2, and `dd647909`/`c70b58b9` (tests) in r3. Three rounds: the maximum.
- Metered **$0.00**.

**Ruled out:**
- Treating the stray `sleep 60` as a leak from the child-only mutant: its parent is another session's `gh pr view` poll loop.
- Counting X26's first red as a kill: it was the 12 s context, not the assertion.
- Reading the race-leg reds on `c70b58b9` as M7a's: three untouched packages, each green isolated under `-race`; `marketingcheck` green in both legs.
- Reading `TestMissionConfigDefaultFlowIntegration` as pre-existing-flaky: unproven either way. It has no prior log mention, and the two non-race reds coincided with M7a's package first running in the suite. It passed in the non-race leg once M7a's bounds were fixed; not attributed further.

**Carried (non-blocking):**
- M7a: X23, X13/X14, S6, S7, N3, N4, R2-1, R3-1, R3-2 above. M7c should take N4 (the shell bound) and R3-2 (group kill in the test cleanup).
- M6: N1, N2 (iteration 300). M5c: M9, M11, M13; the compact-vs-indented public-bytes scan. M5b: `report:` link labels; the dead `utf8.Valid` guard.
- **[HARNESS] ticket:stall:watchdog:foreground-agent-and-wrong-transcript** and **ticket:routing-models-check:bold-colon-heading**: still open.

**Retro:**
- **Skill lane:** none; World does not edit the shared skill.
- **Process lane:** "sequence the judge and the suite" (instances 293, 300) was applied this time: `verify_go.sh` started only after each judge round returned, so the watch-item stays at two. New watch-item, instance 1: an always-on companion package whose subtests carry sub-second elapsed bounds passed every focused gate, `-race -count=3` and three judge rounds, then failed five subtests under the full parallel `go test ./...`. The focused gate list never runs a package under whole-suite load. Remedy, held here: before the judge, run new timing-bounded tests under `go test ./... -count=1` once, or require bounds ≥10× the measured unloaded time. New observation, instance 1: a round-1 judge that derives its own mutants from the diff found a real production defect (N2) the executor's 10-mutant ledger missed, because the ledger was anchored on the plan's mutants plus a directive-supplied list — the directive should ask the executor for branch-complete mutants, not a named list.
- **Backlog lane:** none new.
- Harness share 0 of the last 20. The last three landings (301, 300, 299) moved no UNMET clause; all clauses are MET, so no drift alarm.

**Progress**: World 1.0 **7/7 clauses MET**; **goal unmoved**. Row 150 **13/17 milestones built**; next M7b. Release `v1.0.0-rc.1`; 156 (held on D-WORLD-86) → 150 → 151 before v1.0.0.

**Next:** row 150 M7b: PNG structure, decoded content and exact two-pass verifier (`design_docs/planned/w-marketing-capture-sprint-plan.md`). DECISIONS FOR MARK: D-WORLD-86 (OPEN; default hold). It also gates row 156's close.

## 303 — 2026-10-23 — recover iteration 302's row 150 M7b PNG verifier (PR #361; fresh independent MiniMax PASS 93/100) [PRODUCT]

**Kind:** verify-and-land recovery of inherited M7b; no new production or design change. The operator required all four roles through the Agent tool; all four native supervisors were spawned and launched the resolver's exact provider recipes.

**Pick and why:** Gate 1 found World-owned open PR #361 on its matching worktree, ahead of origin/dev and still unmerged. Its candidate `f615a9698370694664fbbc01431e4a2c549377cc` already contained executor work, two historical MiniMax PASS95 reviews and three green CI jobs, but no iteration 302 log/STATUS/index record. Recovery outranks repeating M7b or starting M7c. The first302 fire crashed at Gate 3 (rc 1 notice https://github.com/sunholo-data/ailang-world/issues/313#issuecomment-6345079873); the subsequent attempt opened the PR and stopped before landing/recording; that latter cause is UNMEASURED. Both attempts are credited here, no invented302 completion entry.

**Clause map:**
- 1 deterministic kernel: MET.
- 2 local-first daemon: MET.
- 3 explicit authority: MET.
- 4 resident-agent floor: MET.
- 5 human surface/provenance: MET.
- 6 protocol boundary: MET.
- 7 controlled self-modification: MET.
These are the ratified current states, not new measurements of the bar. D-WORLD-81 governs the remaining release order156→150→151;156 is held on unchanged OPEN D-WORLD-86, so150 is routable.

**Gate 0:** armed; gh sunholo-voight-kampff; billing CLEAN; no new allowlisted directive on #313. Ledger77 rows valid, only D-WORLD-86 OPEN; no unattended resolution. Shared-store inbox triaged. Fleet harness-resolved message inbox_1792707991628_18d60a93 reports routing-models-check:bold-colon-heading at fleet2debbb0278cc08e018be2d0a1e75befcbda55454; record receipt and ack, no World harness repair. Driver self-notice controls107:4 and 2/2 watermarks fired. Weekly issue rotation not due (#313 created after Monday07:00local;24 comments at initial read).

**Gates1–2:** local dev cfa6d583 was7 commits behind; read state from origin c73c79048daf9c3a19ea05c96e793436eb45e9c5 and PR361 branch, isolated worktree for this recovery; main untracked.claude/worktrees/ left alone. DevCI39784585567 and full SHA check-set3/3 success. Running authoritative skill and all 13 resources match cached fleet origin/dev individually; negative drift control fired. No CWD-relative skill copy is present in World.

**Quorum-at-pick:** canonical lookup returned1 (unreviewed), so fresh text quorum ran despite historical r4 narrow-refinement admission. PROCEED degraded: Gemini 4 Argon alone present/pass; GLM 5.3 and Kimi K3 unreachable, Claude Sonnet quota absent; no absent seat counted as a vote. `--author codex:gpt-6.1-sol` reflected originalr3 author; r4's bounded controller edit was Anthropic. Both author vendors were absent from the actual voting set (only Google voted). $0.02414,10,930 input/228 output. No invented broad quorumPASS and no design rewrite.

**Role execution:** designer Codex audit requires no revision; planner emitted recovery checklist/JSON; its `tools/verify_*.sh` typo was caught and corrected explicitly to actual `scripts/verify_*.sh` before execution. Executor Codex rechecked exactSHA/base, four focused gates and eight named mutants on its exclusive tree; all builds, named assertion reds, restores and originalgreen verified. No permanent production/test change. Controller additionally re-ran eight CLI mutations outside sandbox: seven assertion kills, one historical redundant-cli-reverify survivor, every build/restore/originalgreen verified. Uniform/missing drills reach direct helper assertions, not the later sibling end-to-end assertion; historical truncated-chunk/non-deterministic filesystem/redundant-guard residuals remain. See tracked role reports.

**Controller verification (darwin/arm64, outside sandbox):** pinned AILANG v0.41.0/24ee108, policy-tool v0.54.15/0384450, Go 1.26.6, Z3 4.16.0. `verify_ail.sh`rc 0(27.4s):16 identities,40 named tests,496 se-tools tests. Focused PNG/census gates rc 0. Own eight-mutant rerun BUILDS+intended assertionred+SHArestore+original green 8/8;12 protected paths unchanged across stage work. Full `verify_go.sh`rc 1(448s), plainleg; full raceleg NOT_REACHED. Two unchanged-package interpreter-diagnostic tests fail identically candidate/pristine base; raw pinned output458 bytes has LDR001 at 421 after a startup deprecation warning at 138, but archive outputHead caps 400. Queue 202 owns it; no fix smuggled in. Carried M7a fake-timeout test also failed in the full run (only childPID at 8 s); isolated candidate marketing package plain rc 0(17.9s), race rc 0(35.5s), exact base test rc 0(8.6s). These controls narrow the findings; they do not make the full suite green. Full logs/receipts retained in persistent state; decision-bearing subset tracked under verification/world-iter303/.

**Independent evaluation:** **PASS 93/100, zero blocking** ([actual report](verification/world-iter303/evaluator-recovery-report.md)). Exact reviewed candidate f615a9698370694664fbbc01431e4a2c549377cc, actual MiniMax M3 via canonical pi runner in its own worktree; session_protocol_ack returned acked=true before judge writes. The native Codex supervisor supplied transport only, no acceptance verdict. Eight independently re-run named mutations and four focused controls passed. Non-blocking rubric deductions: test file 804 lines and parser function 52 lines. Read [scope notes](verification/world-iter303/review-scope-notes.md) for the report's JSON-normalization wording, PNG-centered inventory and snapshot-CI scope..

**Result:** M7b independently accepted for verify-and-land,14/17 row 150milestones built; M7c next. No six-shot release-capture PASS, no unattended mint/publish, no fresh-seed hash-equality claim. Overall row 150 stays IN-SPRINT.

**Gate 3b:** initial exact candidate PR CI 39880506556 was3/3SUCCESS. This administrative record extends the same World PR; the controller will check the new PR head and actual merge SHA separately and report LANDED only after all three expected CI jobs are observed successful on both. Final SHA-addressed CI receipts are linked from the bookkeeping digest; no pending/skip is called success here.

**Routing evidence:** base=c73c79048daf9c3a19ea05c96e793436eb45e9c5@2026-10-23T04:01:36Z (Gate4); inherited PR candidate f615a9698370694664fbbc01431e4a2c549377cc.
- Controller codex:gpt-6.1-sol (tok:not reported); native supervisor tokens:not reported.
- Designer recipe codex:gpt-6.1-sol declared:provider-pin; probe10,285 + audit187,409 tok,rc 0,457.9s.
- Planner recipe codex:gpt-6.1-sol declared:planner-lane-default-pin; probe10,285 + run39,452 tok,rc 0,76.7s.
- Executor recipe codex:gpt-6.1-sol declared:provider-pin; probe17,392 input+5 output; real597,484 input(534,656 cached)+5,670 output,rc 0,234.9s. Cached is a subset, not added again.
- Evaluator recipe pi:openrouter/minimax/minimax-m3 declared:provider-pin; probe 1,228 input / 17 output / 128 cache-read; run 189,639 input / 28,591 output / 6,135,424 cache-read; wrapper rc 0 / ok / fenced, 446 s, cost $0.45972282. Judge independence:cross-vendor.
- All four native Agent spawns accepted. This harness exposes no native model-pin argument; each supervisor used the exact explicitly pinned provider CLI rather than claiming inherited model enforcement. No role omitted; no controller self-judge; no runtime fallback. Driver quota degradation is preserved verbatim in the digest routing note.
- Metered $0.48386282 (quorum plus MiniMax run/probe); subscription stages are quota, not dollar-priced.

**Ruled out:** redoing302's already-built milestone; treating r4's historical BLOCKED quorums as PASS; treating source diagnostic reds as introduced by M7b (unchanged bytes, identical base/candidate failures); counting a panic or early sibling failure as the named mutation's assertion; claiming full Go suite green from isolated checks; treating Agent supervisor inheritance as a model pin.

**Carried:** inherited M7a timing/signal/cleanup controls belong to M7c; M7b truncated-chunk guard survivor and filesystem read/check races retained; parent plan's earlier residual notes preserved. New202 is a product diagnostic-output backlog row below the release work, not a human decision and not a loop-harness repair. D-WORLD-86 remains HOLD.

**Retro:** no shared skill or routing-policy edit. Planner gate-path typo corrected from the real script paths; one recorded friction. Canonical quorum directory no longer contains the historical artifact, so fresh review was measured with named absences rather than relying on old filenames. Harness share 0/20 recent indexed iterations. Last three landings moved no UNMET clause; all seven clauses are already MET, so no drift alarm. Fleet routing-checker resolution recorded without local repair.

**Progress:** 1.0 clauses 7/7 MET; row 15014/17 built, advanced M7b; goal unmoved.
**Next:** row 150M7c (capture orchestration/frozen inventories/Chrome drill), then M8a/M8b;156M5c/M5d remain parked on D-WORLD-86;151 follows150.

## 304 — 2026-10-26 — row 150 M7c re-scoped, not built: the design quorum blocked twice and the remaining milestones split into an as-built doc; M7c split into five children (judge PASS 82) [PRODUCT]

**Kind:** design and plan work on the release row. No production code or test changed (`git diff --stat 039bfa77 HEAD -- '*.go' '*.py' '*.sh'` is empty). The operator's standing request had all four roles spawned through the Agent tool. Three ran. The executor was deliberately not spawned, because the design gate blocked; the reason is below.

**Pick and why:** row 150 M7c, the next milestone named by iteration 303. Re-checked against origin before picking:
- no M7c commit or PR exists;
- no open World PRs;
- no stale row-150 worktree;
- local dev equals origin/dev `039bfa77`.

**Clause map:**
- Clauses 1–7 all MET (ratified states, not re-measured). No UNMET clause, so no critical-path row exists and no drift alarm fires.
- D-WORLD-81 sets the release order 156 → 150 → 151. Row 156 is held on D-WORLD-86, which `mission_decisions.sh --open` confirms is the only OPEN row.
- So row 150 is routable.

**Gate 0:**
- Kill switch armed; gh account sunholo-voight-kampff; billing tripwire CLEAN.
- `mission_directives.sh` found 0 allowlisted directives on #313 since 2026-10-22T01:34:42Z, out of 43 comments.
- Self-notices script: rc 1. It reported two crash notices on #313 (comments `6321634905` and `6345079873`), both already credited in log entries 302/303 (one hit each). Control 107:4 fired; 2/2 watermarks read.
- Fleet `harness-resolved` reply for `stall:watchdog:foreground-agent-and-wrong-transcript`: fix `6c16b4ce`, v0.54.18. That commit is an ancestor of the pinned driver `bf22f4943`. No queue row was parked on that ticket, so nothing unparked. Acked.

**Gate 1:**
- Dev CI run 40015682267 success at `039bfa77`; full check set 3/3 success.
- The running skill and all 13 resource files are byte-identical to origin/dev; per-file cmp, all SAME. World has no CWD-relative skill copy.
- Base recorded `039bfa77…` at 07:26:38Z. At Gate 3, `snap` read the same SHA (no drift).

**Quorum at pick:**
- Reader rc 0. The artifact `w-marketing-capture-2026-10-23T03-32-39Z` reads PROCEED at N−1: GLM and Kimi absent (unreachable), Sonnet absent (quota), only Gemini voting. A pass with a named hole.
- The designer's revision made the doc REVISED, so it was re-quorumed with `--author claude:claude-opus-5-5`. This is the 6th round on this doc (r1–r4 in iterations 246–266).
- **r5** (`07-35-08Z`) BLOCKED, $0.2104:
  - Gemini 4 Argon PASS.
  - GLM 5.3 REJECT: the CLAIMED kind premise and dependency SHAs.
  - Kimi K3 REJECT: M7 fallback authority and V29 provenance.
- **One designer revision.** The controller measured GLM's premise first and handed the designer the measurement: `packages/world-demo/evidence.ail:30-42` `authoredClaimGradeCode` with its contract, plus `host/daemon/workbench_packet_evidence.go:265-267`.
- **r6** (`07-44-46Z`) BLOCKED, $0.2070:
  - Gemini PASS.
  - GLM REJECT: the doc is anchored at r4 base `c21f303` while M1a–M7b landed; re-measure V4/V6/V24 at HEAD.
  - Kimi REJECT: the authority under which milestones landed before a full passing quorum; Chrome-156 real-content evidence; the overstated §2 heading.
- **Why the narrow carve-out does not apply:**
  - Every r6 fix needs new measurement or a gating judgment.
  - Kimi's own text holds the remaining milestones until a quorum PASS.
- **Disposition: SPLIT, a controller routing call, not `needs-human-review`.** Per gate-2's rounds rule (b)/(d), the objections localised onto the as-built surface while Gemini passed.
- Both artifacts are in `verification/world-iter304/`.

**Role execution:**
- **Designer** `claude:claude-opus-5-5` (Agent `opus`; the hook allows the same-model alias). Readiness audit: verdict REVISED, 14 interface facts, a split recommendation. Then the r5 revision: V29 (Chrome 156 synthetic-page two-pass check), V30, V31, Kimi's verbatim fallback and acceptance text, and a Quorum r5 paragraph.
- **Planner** `opus`. Wrote the five-child split in plan markdown and JSON (`jq -e` ok). M7c1 is fully specified: `Observe`/`CompareObservations` plus `observe`/`compare` subcommands in `main.go`, 11 subtests including Python parity, and three production mutants. It flagged decisions for later children; those are recorded in `planner-m7c-split.md`.
- **Executor** `claude:claude-haiku-5-5`: not spawned. Its milestone is HELD behind the blocked design gate.
- **Evaluator** `sonnet`, its own worktree `.eval-world-iter304` at `00f7baf9`.

**Controller verification:**
- GLM's re-verification point, measured at HEAD (load 10): `go test -run '^(TestWorkbenchStillAllRegions|TestServeWorkbenchStillFlag)$' ./host/daemon ./cmd/ailang-worldd -count=1 -v` gave both PASS, rc 0, 9 s. These are the landed M6 tests for the still-mode premise V24 simulated.
- Dependency SHAs `38db793`/`9319e6b`/`b8b6a3f` are ancestors of `c21f303` (`git merge-base --is-ancestor`).
- The quorum artifacts were scanned for addresses and secrets; only token-count keys matched.
- No Go/Python/shell file changed, so `verify_go.sh` and `verify_ail.sh` were not re-run. Their scope is code this diff does not touch. CI on the PR is the gate.

**Independent evaluation:** r1 **PASS 82/100, zero blocking**; r2 (record + CI pin delta) **PASS 90/100, zero blocking**. r2 confirms the attribution from job logs (`Image: ubuntu-26.04`; green run `ubuntu-24.04`, every arm `[GATE PASS] … outcome=WROTE`) and that the pin weakens no gate. It asks row 203 to state an exit criterion, and that has been added. r1 detail: ([report](verification/world-iter304/evaluator-report.md)).
- Re-derived V29 version/mtime, V30 lines and V31 ancestry. It did NOT re-run V29's shots.
- 12/12 planner file:line claims confirmed (the `main.go` allow-list entry is at `host/store/context_roots_test.go:205`). The Python one-liner runs: 12 records, 0.03 s.
- Mutants reasoned, not prototyped. One doubt: a position-indexed comparator loop would make M7c1-set panic. The directive must require a path-keyed loop.
- Item 4: the disposition follows the rules (pass-with-hole at pick; r6 is the one re-quorum; carve-out inapplicable; SPLIT correct).
- Two non-blocking defects introduced by this iteration, both fixed by the controller in the reviewers' own words: the §2 heading (Kimi's downgrade text) and the stale Status/Base lines (Gemini's Base annotation plus an iteration-304 status note).
- Base gate re-measured at load 31–34: `shell-hup` red reproduces in the package run AND intermittently in isolation (rc 1/ok/rc 1).

**Result:** M7c is specified and split, not built. CI is pinned to `ubuntu-24.04` after `ubuntu-latest` moved to 26.04 and broke the srt confinement matrix (row 203). Row 150 stays at 14/17; the split adds four milestones, now M7c1–M7c4 + M8a/M8b pending. The design doc keeps the landed M1–M7b record.

**Gate 3b:**
- **PR #364 head `b4efd3e0`, docs-only.** `ailang-code verify gate` success; `go host build + test gate` FAILURE at the "Row 140 srt confinement matrix" step.
  - Every arm, control included, reported `outcome=PARTIAL token=None host=unchanged`.
  - The srt pin digest matched; self-test 16/16.
- **Attribution, measured.** The job's PROTECTED-JSON meta shows `Linux-7.0.0-1017-azure` / glibc 2.43 / bubblewrap 0.11.1. The last green dev run, 40015682267 at `039bfa77`, shows `Linux-6.17.0-1022-azure` / glibc 2.39 / bubblewrap 0.9.0. So `ubuntu-latest` moved under the gate, and this diff touched no code or workflow file.
- **Disposition (Gate 1: a time-based red is the owning mission's job).**
  - `go-verify` and `go-race-verify` pinned to `ubuntu-24.04`, the image every green run measured. The gate itself is unchanged.
  - Queue row 203 filed: root-cause the Ubuntu-26.04 PARTIAL result and unpin.
  - The re-run on the pinned head is the control arm: same gate, same srt, old image.
  - `host/pins` and `cmd/ailang-worldd` tests that read `ci.yml` pass locally.
  - **Pinned head `b1f2014c`: 3/3 success (control arm).** The "Row 140 srt confinement matrix" step passed on `ubuntu-24.04`: same gate, same srt pin, old image.
  - **Final-head race job red once:** `TestWorktreeSnapshotRefusesAConfiguredFilter` control arm at `worktree_delta_test.go:432`. This is row 200's known flake, its 3rd recorded instance. The diff from the green pinned head is docs-only. Re-run of the failed job is the arbiter; instance logged on row 200.
  - LANDED on observed success of the final head and the merge (reported on #366).

**Routing evidence:** base=039bfa77aedff084d26e3cfc71518ed0886e8dda@2026-10-26T07:26:38Z (Gate 1; `snap` at Gate 3 identical). Driver banner: `controller=claude:claude-opus-5-5 … roles: designer=claude:claude-opus-5-5 planner=opus executor=claude:claude-haiku-5-5 evaluator=sonnet`. Degraded: codex over daily ration (probe rc 75) for planner and executor.
- Controller `claude:claude-opus-5-5` (tok: not reported).
- Designer: resolver `recipe claude:claude-opus-5-5 declared:provider-pin`, spawned as Agent `opus` (hook `allow:same-model-alias`). Audit 116,852 tok / 309.7 s; revision 134,325 tok / 94.8 s. The rotation pointer read `codex:gpt-6.1-sol`, so the next entry is GLM 5.3. The driver had demoted GLM this morning (runtime bucket limit) and resolved the designer to claude, so the resolver was followed.
- Planner: resolver `agent-tool opus fail-closed:env-pin`; Agent `opus`, 167,604 tok / 540.1 s.
- Executor `claude:claude-haiku-5-5`: not spawned (design gate BLOCKED).
- Evaluator: resolver `agent-tool sonnet declared:alias-pin`; Agent `sonnet`, r1 72,450 tok / 164.3 s, r2 84,280 tok / 46.4 s. Judge independence: different model from both generators (opus), same vendor (Anthropic).
- Quorum seats: Gemini 4 Argon, GLM 5.3, Kimi K3 (author's vendor Anthropic benched).
- Metered $0.4173 (quorum r5 $0.2104 + r6 $0.2070). A first r5 attempt never launched (`timeout` is absent on macOS, rc 127, $0).

**Ruled out:**
- Running the M7c1 executor while the design quorum is BLOCKED (force-passing).
- A third revision and re-quorum (one re-quorum per round set).
- Filing `needs-human-review` for a SPLIT (rule (d)).
- Treating the 2026-10-23 degraded PROCEED as a full pass.
- Reading GLM's V24 objection as a live defect: the landed M6 still-mode tests pass at HEAD.
- Requiring exact Chrome 155: auto-update removed it, and the design only ever promised in-run equality (D-WORLD-75).
- Attributing `shell-hup` to this diff: no code changed.

**Carried:**
- Next iteration's designer authors `w-marketing-capture-orchestration` with the r6 objections verbatim.
- `shell-hup` timing flake (now red intermittently in isolation at load 31+).
- M8a must list `png_test.go` (cli/unknown); M8b needs a 300 s six-shot timeout.
- M7c4 completion semantics.
- D-WORLD-86 stays HOLD.

**Retro:**
- The pick-time artifact's named hole (two of three reviewers absent since r3) was real: the first round in which both reappeared blocked. The absent-reviewer rule exists for this.
- Friction: the quorum reader helper and my first quorum launch were both written for bash and run under zsh (`no matches found`; `timeout` missing). Both were caught by reading the rc rather than the silence.
- Harness share 0/20. No skill edit.

**Progress:** 1.0 clauses 7/7 MET; row 150 14/17 built, M7c re-scoped; goal unmoved.
**Next:** row 150, the designer authors `w-marketing-capture-orchestration` (M7c–M8, as built at HEAD, r6 objections verbatim). Quorum once; on PASS, build M7c1 (fully specified in the plan). Then 151 after 150.

## 305 — 2026-10-26 — row 150: the split doc `w-marketing-capture-orchestration` written as built and routed under the narrow-refinement carve-out (quorum r1/r2 blocked); M7c1 built (PR #370; judge PASS 81) [PRODUCT]

**Kind:** design, plan and one production milestone on the release row. All four roles were spawned through the Agent tool, per the operator's standing request.

**Pick and why:** row 150, following iteration 304's **Next**: author the split as-built doc, quorum it, then build M7c1. Re-checked before picking:
- no orchestration doc, commit or PR existed;
- the only open World PRs were #365 (website film, attended) and #368 (the attended `w-world-surface-outputs` design), neither of them this loop's;
- local dev equalled origin/dev `4e6e89dc`.

**Clause map:**
- Clauses 1–7 all MET (ratified states, not re-measured). No UNMET clause, so there is no critical-path row and no drift alarm.
- D-WORLD-81 order is 156 → 150 → 151. 156 is held on D-WORLD-86, the only OPEN row (`mission_decisions.sh --open`).

**Gate 0:**
- Kill switch armed; gh account sunholo-voight-kampff; billing tripwire CLEAN.
- First fire after the weekly rotation to #366. `mission_directives.sh`: 0 allowlisted directives on #366 (2 comments) and on #313 since 2026-10-22T01:34:42Z (44 comments).
- Self-notices rc 1: newest crash `6345079873` (rc=1, 2026-10-22) is already credited (2 hits in the log). The `mission-366-last-seen` watermark was absent (DEGRADED, single watermark) and is written this iteration.
- Inbox: `mission-world` held only this loop's own claims; nothing outranks the queue.
- **Weekly external-issue sweep:** 0 orphans of 1 enumerated open issue (#366, the bookkeeping thread: charter 0 / log 1 / status-archive 0 / dashboard 0 / log-archive 0). `open_issues_count` 3 = 1 issue + 2 PRs. Negative control fired (0 everywhere).

**Gate 1:**
- Dev CI at `4e6e89dc`: check set 3/3 success.
- Running skill plus all 13 resource files byte-identical to origin/dev (28 `ok` across the resolved symlink target and the pin copy).
- Base `4e6e89dc…` recorded at 11:26:33Z; `snap` at Gate 3 was identical. origin/dev later moved to `6ec3ee39` (#365, website only), merged into the sprint branch before the record.

**Role execution:**
- **Designer** `claude:claude-opus-5-5` (resolver `recipe … declared:provider-pin`; Agent `opus`).
  - r0: 284 lines, V32–V49. It includes landing provenance for M1a–M7b (merge SHAs; authorities D-WORLD-69 A, the iteration-266 r4 carve-out, "none owed" 267–276, "not re-run" 300/301, the single-voter PROCEED at 303), stated plainly as "no all-seats round ever passed".
  - The as-built inventory re-ran every landed milestone's pinning tests at base (V35).
  - V37–V40 re-measured parent V4/V6/V14/V24 at HEAD, including the real `--workbench-still` flag on the real seeded store: 4 routes repeat-equal, 0 scripts, still = default minus the tag.
  - **V41:** Chrome 156.0.8078.47 real content, 6/6 byte-identical across two passes.
  - **V42 (new):** packet/graph/provenance PNGs at the planned 1440×1400/1100 contain only the header and timeline; the required regions are below the fold. Controller-confirmed by viewing `packet.png`.
  - Found stale: the plan's v0.52.1 tool pin (now v0.54.15 since `4b686d8e`), and M7c4's "seed with `demoM5cRealSeed` then stop the rig", which cannot work (no daemon handle; the writer lock).
  - r1 revision (one, protocol-mandated): V50–V55, exact two-seed bottoms (V51: identical across seeds; heights 3700/2100/3000), position assertion P0, release-store containment P3, C14 full-path real-Chrome rehearsal.
- **Planner** `opus` (resolver `agent-tool opus fail-closed:env-pin`).
  - 29 amendments marked in the plan; M7c0 added.
  - M7c1 prototyped in a scratch copy (11/11 subtests, 3 mutants red, `-race` ok) and written as exact code.
  - Flagged, not decided: M8a vs R1 ordering; `main.go` in M8a; M7c2a depending on M7c0.
- **Executor** `claude:claude-haiku-5-5` (resolver `recipe … declared:provider-pin`; Agent `haiku`). Built M7c1 in `.wt-world-iter305-row150-m7c1`:
  - `main.go` 83 → 297 lines; `observation_test.go` 255 lines.
  - Gates: gofmt/vet clean; 11/11 subtests; `-race` ok 35 s; `TestProductionGoSurface` and `TestSubprocessSinkGate` (violations 0).
  - 3 named mutants each red at the named assertion and reverted byte-identical. M7c1-set reports `observation-digest` for the deleted path rather than panicking (path-keyed, per the iteration-304 judge).
- **Evaluator** `sonnet` (resolver `agent-tool sonnet declared:alias-pin`), own worktree `.eval-world-iter305` at `f354e5e5`.

**Quorum (new doc, author `claude:claude-opus-5-5`; seats Gemini 4 Argon, GLM 5.3, Kimi K3; Anthropic benched):**
- Reader: rc 1 for the new doc, control rc 0 for the parent.
- **r1** `w-marketing-capture-orchestration-2026-10-26T11-47-05Z`, $0.3335, BLOCKED: Gemini PASS; GLM REJECT (light/phone "unaffected" asserted; seed and release-store invariance); Kimi REJECT ("committed supervisor" wording; two-seed bottoms; the M7c0 test could not red its mutant).
  - The controller measured the premises first (rule 3f): light and phone meet their parent rows (phone overflows horizontally, title clipped, cosmetic).
  - The probe DID run the committed M7a supervisor `scripts/marketing_capture.py shot`; Kimi's objection was partly refuted, though `daemon` mode and the shell were not in the path.
  - `png_test.go:462` hardcodes the shot list (Gemini confirmed).
  - These measurements were handed to the designer for r1.
- **r2** `…T12-00-33Z`, $0.3266, BLOCKED: Gemini PASS (catch: tight phone rectangles); GLM REJECT (nine §3 cmd results were from the retired v0.52.1 pin; threshold-test ownership; durable artifact path; V43 wording); Kimi REJECT (C8′ claimed the layout assertion for phone, which the dump-dom instrument cannot measure).
- **Disposition: narrow-refinement carve-out (controller r2).** Every r2 objection carried a concrete `proposed_fix` and none disputed direction (Gemini passed every round). The controller applied the fixes in the reviewers' terms:
  - **V56:** GLM's exact command, rc 0, 9/9 PASS under v0.54.15.
  - **V57:** `TestMarketingPNGContent` pins all provisional thresholds as boundaries at `png_test.go:359-360/396/438-441/444-455`; M8a owns the change.
  - Durable `quorum-r1.json`/`quorum-r2.json`.
  - Kimi's C8′ text verbatim.
  - **V58:** a phone screenshot scan; both passes PASS, timeline-only crop and blank FAIL, desktop light FAIL as a cross-layout control. A first run picked the stateRoot underline as "entry 0" (index 7); fixed to the measured band index 9, and the superseded run is recorded.
  - Recorded in the doc as a carve-out, never as a PASS; C13 and the header say so.

**Controller verification:**
- Spot checks before quorum: tool pin (`ailang-releases.json` tool = v0.54.15), `supportedWorkbenchQuery` refuses `from` on object routes, 1 `id=` vs 20 `aria-label=` in render.go, V42 by image.
- Region-bottom method reviewed: a measurement-only proxy appends its script after every section, cross-checked against pixels (V52).
- M7c1 re-run outside the executor: `go build ./...` rc 0; `go test -run '^$' ./...` compile fence rc 0 (34 pkgs); vet/gofmt clean; `go test ./scripts/marketingcheck` ok; `TestProductionGoSurface` ok; `TestSubprocessSinkGate` ok; 11 subtests PASS.
- `verify_ail.sh` not run: no `.ail` changed (CI's `ailang-code verify gate` is the gate).
- New artifacts scanned for addresses: 0.

**Independent evaluation:** r1 **PASS 81/100** (A 50/60, B 31/40), zero blocking ([report](verification/world-iter305/evaluator-r1.md)).
- Re-ran all gates; non-vacuity holds (reverting `main.go` fails to compile the tests).
- About 35 diff-derived mutants, about 25 killed. No survivor is on a core comparison branch. Survivors: missing-path tolerance (low-medium); two-sided duplicate and length-only set guard; `IsRegular`; `[]` vs `null`; CLI rc/mode/newline.
- Re-derived V56/V57/V58/V42. Carve-out conditions (a) and (b) held. It flagged that "applied in the reviewers' own words" overstated (finding 5: the V58 scan, the renumbering, the C12 wording and the M8a owner/gate are controller-authored) and that the ±4 px rule was claimed, not demonstrated.
- Controller fix `dde9e0f5`: §9 now names exactly which content is controller-authored; Gemini's tight-rectangle text is in R1; **V58b** shift arms 2/4/5/10 px give rule (i) true/true/false/false, discriminating exactly at ±4. The 4 px overall FAIL is from (iii) only, an artifact of the synthetic shift moving the viewport-fixed scrollbar; the judge independently measured 2/5/10 → PASS/FAIL/FAIL.

**Result:** the orchestration doc r2 governs M7c0–M8b; M7c1 BUILT (PR #370). Row 150 has 15 of 22 milestones built (M7c0 is new). The new V42 finding means the parent's six-shot list would have shipped three release images without their subject; M7c0 fixes that before any capture.

**Judge r2 (delta `dde9e0f5` + merge `f200b9d2` + record `bb79d7f4`):** **PASS 92/100**, zero blocking. Finding 5 closed (shift arms re-run independently: rule (i) true/true/false/false). Record accurate (quorum costs $0.3334/$0.3266 vs logged $0.3335/$0.3266; no implied PASS). Rotation lossless (archive `ITERATION 301` = 1, charter 6258 → 6258, rows 150/203 present, ledger valid 77 rows). Merge brought `website/**` only. Non-blocking: R1's ">4 px fails" sentence is enforced through the ±4 px tolerance, not by rectangle tightness alone. r2 is appended to the report.

**Gate 3b:** PR #370 head `bb79d7f4`: check set 3/3 success (`go host build + test gate`, `go race gate`, `ailang-code verify gate`; bounded poll pinned to the full SHA). This commit (docs-only: log, judge report copy) is re-polled before merge, and the merge commit after it; the outcome is reported on #366.

**Routing evidence:** base=4e6e89dc44c9ee9771b415b3cfdc23fadfaadff1@2026-10-26T11:26:33Z (Gate 1; `snap` at Gate 3 identical). Driver banner: `controller=claude:claude-opus-5-5 … roles: designer=claude:claude-opus-5-5 planner=opus executor=claude:claude-haiku-5-5 evaluator=sonnet`. Degraded: codex over daily ration (probe rc 75) for planner and executor.
- Controller `claude:claude-opus-5-5` (tok: not reported).
- Designer `claude:claude-opus-5-5` (Agent `opus`): r0 278,119 tok / 1,029.7 s; r1 revision 338,359 tok / 376.4 s. The revision fired because quorum r1 blocked: one doc, one protocol-mandated revision, within the diet. The r2 revision was the controller's own (carve-out), not a third designer run.
  - The rotation pointer still reads `codex:gpt-6.1-sol` (next entry: GLM 5.3). The resolver/driver resolved the designer to the claude pin, and the resolver was followed verbatim. Pointer not advanced, since claude ran.
- Planner `opus`: 223,882 tok / 953.6 s.
- Executor `claude:claude-haiku-5-5` (Agent `haiku`): 53,444 tok / 98.4 s.
- Evaluator `sonnet`: r1 89,887 tok / 392.1 s. Judge independence: different model from every generator (opus designer/planner/controller, haiku executor), same vendor (Anthropic).
- Quorum seats: Gemini 4 Argon, GLM 5.3, Kimi K3 (author's vendor benched; reserve gpt6-1-sol unused).
- Metered **$0.66** (quorum r1 $0.3335 + r2 $0.3266).

**Ruled out:**
- A second designer revision after r2: over the diet. The carve-out applies, so the controller revised.
- Parking `needs-human-review` after r2: every objection was concrete and none disputed direction, which is the carve-out's exact case.
- A third quorum round: one re-quorum per round set.
- Kimi r1's "supervisor not in path" as stated: partly refuted by the probe source (`marketing_capture.py shot` IS the M7a supervisor).
- URL fragments for framing: blank yet byte-identical (V43).
- A still-only layout reorder: it breaks M6's still = default − tag invariant.
- Allowing `from` on object routes: a workbench grammar change outside this row.
- Building M7c0 this iteration: R0's re-confirm V-row is a precondition (C13).
- Treating the 4 px V58b overall FAIL as a tolerance defect: it is the synthetic shift's scrollbar artifact; rule (i) itself discriminates at ±4.

**Carried:**
- R0 before M7c0.
- Phone horizontal overflow (title clipped) and provenance horizontal clip (V52): cosmetic, renderer CSS outside row 150, residual for Mark if he wants different framing.
- Tall framing aesthetics (packet 3700 px) are a product choice for Mark (new row if raised).
- Planner flags: M8a vs R1 ordering (M8a reports the threshold step DEFERRED if R1's V-row is absent), `main.go` in M8a, M7c2a depending on M7c0.
- Judge M7c1 non-blocking test gaps (missing-path tolerance, two-sided duplicate, CLI rc/mode branches).
- `shell-hup` flake (green at load 2–4 this iteration).
- D-WORLD-86 stays HOLD.

**Retro:**
- Rule 3f paid twice. Measuring r1's premises before the revision turned one objection into a partial refutation and two into exact fixes. Measuring the designer's region method before the quorum avoided a round on an unverified instrument.
- The judge caught the controller overstating its own carve-out ("verbatim") — the evaluator working as designed on the controller's hands.
- The instrument-bug-in-my-own-scan (index 7 → underline) was caught by printing the band list rather than trusting a PASS.
- Harness share 0. No skill edit.

**Progress:** 1.0 clauses 7/7 MET; row 150 15/22 built; goal unmoved.
**Next:** row 150. Controller R0 (two-seed re-confirm of V51 on the then-current Chrome, `iter305-region-bottoms-probe`), then M7c0 (shot geometry, fully specified in the plan), then M7c2a. Then 151 after 150.

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
