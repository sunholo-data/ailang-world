# Ailang World Mission — iteration log

Append-only. One entry per outer-loop iteration. Newest at the bottom.

---

> **Older entries are ARCHIVED.** This file holds the newest 40. The full record of every
> iteration is in `world-mission-log-archive.md`, and a one-line index of ALL of them —
> the thing to grep before picking work, so the loop never repeats itself — is in
> `world-mission-index.md`.

> **Older entries are ARCHIVED.** This file holds the newest 40. The full record of every
> iteration is in `world-mission-log-archive.md`, and a one-line index of ALL of them —
> the thing to grep before picking work, so the loop never repeats itself — is in
> `world-mission-index.md`.

> **Older entries are ARCHIVED.** This file holds the newest 40. The full record of every
> iteration is in `world-mission-log-archive.md`, and a one-line index of ALL of them —
> the thing to grep before picking work, so the loop never repeats itself — is in
> `world-mission-index.md`.

> **Older entries are ARCHIVED.** This file holds the newest 20. The full record of every
> iteration is in `world-mission-log-archive.md`, and a one-line index of ALL of them —
> the thing to grep before picking work, so the loop never repeats itself — is in
> `world-mission-index.md`.

> **Older entries are ARCHIVED.** This file holds the newest 20. The full record of every
> iteration is in `world-mission-log-archive.md`, and a one-line index of ALL of them —
> the thing to grep before picking work, so the loop never repeats itself — is in
> `world-mission-index.md`.

> **Older entries are ARCHIVED.** This file holds the newest 20. The full record of every
> iteration is in `world-mission-log-archive.md`, and a one-line index of ALL of them —
> the thing to grep before picking work, so the loop never repeats itself — is in
> `world-mission-index.md`.

## 274 — 2026-10-16 — row 150 M4a seed S5 real approval transition + one pending packet + exact budget delta LANDED `7ed1152` (PR #283; judge r1 PASS 91) [PRODUCT]

**Kind:** one product milestone on the D-WORLD-69 release path. Seed stages S1–S5 now run. S6 (packet listing, `why` walk) and the public manifest come in M4b, and the seed still refuses completion.

**Gate 0:** kill switch armed (`mission-world.disabled` absent); gh `sunholo-voight-kampff`; billing CLEAN. `mission_directives.sh` → 0 allowlisted directives on #246 since `2026-10-16T04:02:59Z` (31 comments). `mission_gate0_self_notices.sh` rc 1 on the same 2026-10-13T02:01:58Z rc=143 notice (orphan 255, already credited by 256); control 107:4 ok, watermarks 2/2. Unread `mission-world` inbox: 0. Weekly rotation not due (#246 created 2026-10-12T12:08Z).

**Gate 1:** local dev = origin/dev = `206ce2e` (Gate-1 base `206ce2e2822acc60f6bcb613d3f06bf04238106e@2026-10-16T11:25:48Z`); dev CI success on `206ce2e`, check-runs 2, none not-green. Skill directory: SKILL.md and all 13 `resources/*.md` at the resolved symlink target match V1 origin/dev (per-file check, zero DIFF lines). Open PRs by this account: #282 (row 139 M1, `sprint/row139-m1-release-manifest`) and #281 (row 183 trial runner). Both come from attended `.claude/worktrees/agent-*` worktrees, and neither branch carries a later log entry. Also #166 (row 114's preserved old draft). Ledger: D-WORLD-75 the only OPEN row.

**Gate 2 — clause map:** 1/2/3/4/6/7 MET. 5 UNMET: row 139 is NEXT as an attended build (D-WORLD-78, attended PR #282 open), row 183 is IN-SPRINT ATTENDED (PR #281 open), and row 114 waits on 183's stores. No loop-routable UNMET-clause row → **pick row 150 M4a** (iteration 273's Next). Reality check: no M4a commit on origin/dev (`git log --grep M4a` shows only unrelated rows), no M4a branch, worktree or index line; sprint JSON M4a `pending`. Quorum: none owed (plan r4 already routed in 266).

**Execute:** resolver `recipe codex:gpt-6.1-sol declared:provider-pin`; probe rc 0. Before writing the directive, the controller baselined the gate list on pristine `206ce2e` outside the sandbox: 18/18 `TestDemo*` PASS, plus `TestDaemonInvocationMintsDecisionPacket`, `TestRESTCommitRefusesDecisionPacketObjects` and `TestDecisionPacketCodec`. The directive carried the plan's M4a contract and iteration 273's lessons: skip on an unset tool pin like `cliSeBins`, pin-free controls for every pure guard, one named negative control per guard plus a run-it mutant, REST read-back, and no git writes. The executor ran about 19 min (219,425 tok) and changed only the 2 allowlisted files. S5 sends one `demo-approval` A2A call to `demo-approve`. Its roster holds the S3 test invocation record, the S2 read effect record, an S1 authored claim and the S4 envelope. The goal and candidates come from the S1 committed objects. After S5 the seed refuses with "S1-S5 recorded; S6 incomplete; no success manifest written". New tests: `TestDemoSeedApprovalPacket` (real checked-in approve.ail on the pinned capsule, coordinator and HumanHandler), `TestDemoSeedApprovalControls` and `TestDemoSeedApprovalBudgetControl` (pin-free). The budget delta is read from REST effect records against the fixture's minted grant and cross-checked against the broker record's BudgetBefore/After. The ledger has 13 mutants, each red at a named assertion, including the plan's M4a-extra-effect (the coordinator's L1 refusal is visible in the diagnostic). The controller re-ran outside the sandbox and scanned the diff for writes outside the worktree first (all `os.WriteFile` sites are under `t.TempDir()`; `TMPDIR` is set to `t.TempDir()`). Results: vet clean; gofmt clean on the diff (the one gofmt hit, `host/store/deadline_guard_test.go`, is pre-existing on base, row 101); 21/21 `TestDemo*` PASS with both pins; with the tool pin unset rc 0, 18 PASS and 3 SKIP (ReadTestCheck, ProofStage, ApprovalPacket); `-race` on the three new tests clean; the two daemon packet tests PASS. Committed as `ded5f0b` (draft PR #283).

**Independent evaluation:** Agent `sonnet` (`agent-tool sonnet declared:alias-pin`; generator codex ≠ judge sonnet), in its own detached worktree `.eval-world-iter274`. r1 **PASS 91/100, zero blocking.** It reproduced every gate. It found the budget observation not circular. It ran 11 mutants (the plan's M4a-extra-effect, roster ref swaps, a forced `EntryIndex=0`, a replaced goal, a duplicate S5 call, plan cost, output status, published access cost 1, and approvals-chain pending→resolved), all 11 killed at named assertions. The controller checked the eight cited assertion lines in the file. Non-blocking: a live double spend is covered only by the pure predicate plus the broker cross-check; the seed ignores `json.Unmarshal` errors on its own fixed S1 payloads; a missing `entry_index` would decode to 0 (the test catches it); CI without the tool pin skips the only real S5 test (row 160). Reports: `design_docs/verification/world-iter274/evaluator/r1.md`, `…/executor/REPORT.md`. The raw mutant logs are banked under `~/.ailang/state/world-iter274/snap/M4a/evidence/`.

**Gate 3b:** PR #283. ailang-code verify success on both heads. The go gate **failed** on `ded5f0b` and on `cd4d0e8` (reports-only delta). Each log was about 100 KB with 84 `--- PASS` and 0 FAIL, ending `verify_go.sh: FATAL: -race leg timed out after 600s`. That is row 181's flake, instances 7 and 8; run 38897729178 attempts 3 and 4 were instances 9 and 10, same signature. Per-package race times show the diff is not the cause: `cmd/ailang-worldd` took 83.1 s and 81.7 s, against 78.6–81.5 s on recent greens, and the new S5 test skips in CI. `host/broker` took 344 s and 338 s, `host/daemon` 310 s and 302 s. PR CI 2/2 on `cd4d0e8` at run 38897729178 attempt 5. Re-fetch before merge: origin/dev unchanged at `206ce2e`; `mergeable` MERGEABLE/CLEAN; closing-keyword scan of the PR title and body clean (control fired). Squash-merged `7ed1152`. `git diff cd4d0e8 7ed1152` is empty, so the merge tree is the green PR head. The merge-push CI verdict is in the record PR. Filed **D-WORLD-79** (promote row 181? default B = keep the order).

**Record:** no separate record judge. Gate 4 does not require one, and the record carries no new claim beyond the judged milestone, the CI logs and D-WORLD-79.

**Routing evidence:** base=206ce2e2822acc60f6bcb613d3f06bf04238106e@2026-10-16T11:25:48Z. Controller `claude-opus-5-5` (tok: not reported) · designer not spawned (no doc needed) · planner not spawned (plan exists, M4a ready-made) · executor `codex:gpt-6.1-sol` via recipe for resolver `recipe codex:gpt-6.1-sol declared:provider-pin` (219,425 tok) · evaluator Agent `sonnet` (66,102 tok). judge-independence: cross-vendor. Metered $0.

**Ruled out:**
- Row 139 or 183 as the pick: both are attended builds with open attended PRs (#282, #281) from attended agent worktrees. The loop must not take an attended row.
- A fix round on r1's non-blocking findings: none blocks, and the residuals (S1 unmarshal errors, `entry_index` default) belong with M4b, which edits the same seed and validates the full why set.
- Attributing the go-gate reds to the diff: 0 FAIL in both logs, the S5 integration test skips in CI, and `cmd/ailang-worldd` race time is unchanged within runner variance.

**Retro:** Skill lane: none. Process lane: carrying iteration 273's lesson into the directive (skip on an unset tool pin, pin-free controls) worked. The executor shipped no CI-only defect this time, and r1 passed with zero survivors among the judge's own mutants. The remaining CI cost is row 181: each landing now spends 2–3 go-gate attempts on the 600 s race cap. Backlog lane: row 181 gained instances 7–10 and is now D-WORLD-79 (promote it? default keep order); row 160 gained one (TestDemoSeedApprovalPacket is rig-only).

**Progress**: World 1.0: **6 of 7 clauses MET**; clause 5 is the last UNMET, and its rows (139, 183, then 114) are attended. Release row 150 advanced to 7/17 milestones.

**Next:** re-read the row 139 and 183 tags first. If either has become loop-routable, it outranks release work under Gate 2 rule (b). If D-WORLD-79 is answered A, take row 181. Otherwise row 150 M4b (complete why validation and public manifest assembly). Fold in r1's NB2/NB3 (S1 unmarshal errors, `entry_index` presence).

## 275 — 2026-10-16 — row 150 M4b seed S6 why validation + one anchored packet + public manifest assembly (fail-closed) LANDED `acb3f54` (PR #289; judge r1 FAIL 66 → r2 PASS 86 → r3 PASS 86) [PRODUCT]

**Kind:** one product milestone on the D-WORLD-69 release path. Seed stages S1–S6 now run. The seed still refuses success: the roster-grade checks and manifest publication are M5a–M5c.

**Gate 0:** kill switch armed (`mission-world.disabled` absent); gh `sunholo-voight-kampff`; billing CLEAN. `mission_directives.sh` → 0 allowlisted directives on #246 since `2026-10-16T04:02:59Z` (32 comments). `mission_gate0_self_notices.sh` rc 1 on the same 2026-10-13T02:01:58Z rc=143 notice (orphan 255, already credited by 256); control 107:4 ok, watermarks 2/2. Weekly rotation not due (#246 created 2026-10-12T12:08Z, 32 comments). Ledger: D-WORLD-75 and D-WORLD-79 OPEN, neither answered.

**Gate 1:** local dev fast-forwarded 2 commits to origin/dev `157adc5` (Gate-1 base `157adc557095e44206d2b55e74bcd931f0c1582a@2026-10-16T15:26:09Z`). Dev CI success on `157adc5`, 3 check-runs (verify, go host, and the new separate race job from attended `38e7859`), none not-green. Skill directory: SKILL.md and all 13 `resources/*.md` at the resolved symlink target match V1 origin/dev (per-file, 14/14 `same`), and so does the pin copy. Open PRs by this account: #286 (row 114 Q3), #287 (row 184) and #288 (row 139 M2), all from attended worktrees under `.claude/worktrees/agent-*` (two of them held uncommitted edits, left alone), plus #166 (row 114's preserved old draft). None carries a later log entry.

**Gate 2 — clause map:** 1/2/3/4/6/7 MET. 5 UNMET: row 139 is an attended build (M1 landed `157adc5`, M2 PR #288 open at pick time), rows 183 and 184 are IN-SPRINT ATTENDED, and row 114 waits on 183's stores. Row 181 (D-WORLD-79) was already fixed attended by `38e7859` (race leg in its own job, 1500 s bound), so its question no longer gates the order. No loop-routable UNMET-clause row → **pick row 150 M4b** (iteration 274's Next). Reality check: no M4b commit on a fresh origin/dev (`git log --grep M4b` shows only unrelated rows), no M4b branch or worktree, sprint JSON M4b `pending`. Quorum: none owed (plan r4 routed in 266).

**Execute:** resolver `recipe codex:gpt-6.1-sol declared:provider-pin`; probe rc 0. The controller baselined on pristine `157adc5` outside the sandbox: 21/21 `TestDemo*` PASS (269 s) and 5/5 why/CLI tests. The directive carried the plan's M4b contract (reuse the real `worldReader.locate`/`buildChain`, coordinator checked first, fail-closed success), the M4a r1 NB2/NB3 fold-ins, tool-pin skip gating with pin-free controls, and one named negative control per guard. r1 executor (184k tok, about 25 min) changed only `demo_seed.go` and `demo_test.go`. It reported 26 mutants and declared two strictness gaps honestly: binary-pin bytes cannot be fetched over REST, and not every transport error branch was mutated. Controller re-run outside the sandbox: vet and gofmt clean; 31/31 named PASS with both pins; whole package rc 0 with the tool pin unset; `-race` clean on the new tests. Committed `e0b4f09`, draft PR #289.

**Independent evaluation:** Agent `sonnet` (`agent-tool sonnet declared:alias-pin`; generator codex ≠ judge sonnet), in its own detached worktree `.eval-world-iter275`. **r1 FAIL 66/100, 1 blocking:** no test asserted what the manifest CONTAINS. Deleting the packet/request/proposal refs or the four evidence refs left every test green. The controller reproduced it first-party: with the triple emptied, `TestDemoSeedManifestReferences` and `TestDemoSeedWhyNamedLinks` both PASS. Non-blocking: two packets were accepted (last one silently wins), a forced packet anchor survived, the seed's fail-closed return was pin-gated only, and one import was misplaced. **Fix round r2** (codex, 133k tok): an exact-cardinality role check in `demoValidateManifest`, table-driven pin-free remove/duplicate controls per role, a real-stack role multiset, two-packet and wrong/missing-anchor proxy controls, and a pin-free drive of `demoSeedApproval`. O6/O10/O5/O8/O7 were each red pin-free at the named assertion. The controller re-ran the gates (31/31 PASS, unset package rc 0, race clean) and its own O6 (killed pin-free at `demo_test.go:1973`). **r2 PASS 86/100, zero blocking:** all five r1 survivors killed with the pin unset; the closed-daemon `store.PutObject` packet fixture under `t.TempDir()` judged an honest seam. Three new judge mutants survive (stray extra ref, loosened why-hash count, duplicate-entry check), all non-blocking.

**Gate 3b, and the one red:** PR #289's race job FAILED on `cc39b9f` (r2 code plus reports). 43 `TestDemoSeedManifestReferences/controls/completeness/*` subtests hit `context deadline exceeded`. Cause: the controls shared the parent's single 30 s `boundedTestContext` reader across about 120 subtests, and each re-reads the whole manifest over REST, so under `-race` on the 2-core runner the budget ran out about 31 s in. The Mac finishes the block in 16 s, which is why every local run passed. Controller fix `d0fc657` (test-only): a per-subtest bounded reader at 12 call sites. Discriminating drill with the bound cut to 10 s under `-race` and `GOMAXPROCS=2`: old controls 49 failures, new 0; the judge re-ran it and got 46 versus 0. **r3 (delta review) PASS 86, zero blocking:** only the deadline changes, the `mismatched-content` proxy path is intact, and O6 is still killed. PR CI 3/3 on `d0fc657`. Re-fetch before merge: origin/dev had moved 3 attended commits (`44c2008` row 184 changes `why.go` to add a `Worktree` field to links; `4c40690` row 139 M2; `a08ca14` charter). Rebased onto them and re-ran: 32/32 named PASS with both pins, package rc 0 with the tool pin unset, vet clean. PR CI 3/3 on rebased head `adc0c72`; `mergeable` MERGEABLE/CLEAN; closing-keyword scan clean (control fired). Squash-merged `acb3f54`; `git diff adc0c72 acb3f54` is empty. Merge-push CI verdict is in the record PR.

**Record:** no separate record judge; the record carries no new claim beyond the judged milestone, the CI logs and the drill.

**Routing evidence:** base=157adc557095e44206d2b55e74bcd931f0c1582a@2026-10-16T15:26:09Z. Controller `claude-opus-5-5` (tok: not reported) · designer not spawned (no doc needed) · planner not spawned (plan exists, M4b ready-made) · executor `codex:gpt-6.1-sol` via recipe for resolver `recipe codex:gpt-6.1-sol declared:provider-pin` (r1 184,108 tok, r2 132,745 tok) · evaluator Agent `sonnet` (`agent-tool sonnet declared:alias-pin`; one transcript resumed twice, 109,991 tok cumulative over r1–r3) · controller-authored test fix `d0fc657`, reviewed by the same independent judge (r3). judge-independence: cross-vendor. Metered $0.

**Ruled out:**
- Row 181 / D-WORLD-79 as the pick: the race-leg fix already landed attended (`38e7859`). This PR's race job ran in its own job and passed on 3 of 4 heads; the one red was a genuine test-deadline defect, not the cap.
- Row 185 as the pick: it was "NEXT after row 184", and row 184 was still IN-SPRINT ATTENDED at pick time (it landed `44c2008` mid-fire).
- Re-running the race red instead of diagnosing it: 43 failures carried one `context deadline exceeded` signature on one shared context, and the 10 s drill showed the defect on the rig too.
- A second fix round for r2's three surviving judge mutants: all non-blocking, and M5c, which validates the success manifest, edits the same validator.

**Retro:** Skill lane: none. Process lane: (1) a table-driven control that shares one bounded context across many REST-heavy subtests is a CI-only time bomb. Local speed hides it, and a `-race` run on 2 cores shows it. The cheap discriminating check is to cut the bound, and future directives should demand per-subtest deadlines for table-driven network controls. (2) Re-fetching before merge caught an attended `why.go` change to the very producer S6 consumes. Rebase-and-retest cost one CI cycle and is the right price. Backlog lane: row 160 gains an instance (the M4b real-stack arms skip in CI); row 181 gains a note (fixed attended, D-WORLD-79 likely moot).

**Progress**: World 1.0: **6 of 7 clauses MET**; clause 5 is the last UNMET, and its rows (139, 183, then 114) are attended. Release row 150 advanced to 8/17 milestones.

**Next:** re-read the 139, 183 and 185 tags first. Row 185 (`w-episode-requires-git`, clause 3 MET +5) became open when row 184 landed. It is an attended "NEXT after 184" ordering, so if no attended session holds it, it outranks release work: take it (design first). Otherwise row 150 M5a (bounded anchored roster, ATTESTED and authored CLAIMED). Fold in r2's three surviving judge mutants (stray extra ref, why-hash count, duplicate entry) when M5c touches the validator. Every table-driven network control gets a per-subtest deadline.

## 276 — 2026-10-16 — row 150 M5a bounded anchored evidence roster (ATTESTED + demo-authored CLAIMED) LANDED `bcd2d7c` (PR #295; judge r1 PASS 90) [PRODUCT]

**Kind:** one product milestone on the release path (D-WORLD-69; order now set by D-WORLD-81). The validated decision packet gains an internal evidence roster. TESTED/PROVEN resolution and all rendering are M5b, and seed roster validation is M5c.

**Gate 0:** kill switch armed (`mission-world.disabled` absent); gh `sunholo-voight-kampff`; billing CLEAN. `mission_directives.sh` → 0 allowlisted directives on #246 since `2026-10-16T13:53:04Z` (33 comments). `mission_gate0_self_notices.sh` rc 1 on the same 2026-10-13T02:01:58Z rc=143 notice (orphan 255, already credited by 256); control 107:4 ok, watermarks 2/2. `mission-world` unread inbox empty (`-unread -limit 500 -json` → null). Weekly rotation not due (#246 created 2026-10-12T12:08Z).

**Gate 1:** local dev == origin/dev `24eb2a3` (Gate-1 base `24eb2a3102d2c8b9cf69604ce600d6685f11dbf4@2026-10-16T19:25:43Z`). Dev CI 3 check-runs, none not-green. Skill directory: SKILL.md and all 13 `resources/*.md` at the resolved symlink target match V1 origin/dev (per-file, 14/14 `same`), and so does the pin copy (14/14). Open PRs by this account: #291 (row 185 built as attended row-183 trial episode w185, draft for Mark), #286 (row 114 Q3, attended) and #166 (row 114's preserved old draft). None carries a later log entry.

**Gate 2 — clause map at pick time:** 1/2/3/4/6/7 MET. 5 UNMET: rows 139 and 183 attended, 114 waiting on 183. Row 185 ("NEXT after 184") was held: the attended row-183 trial had built it as draft PR #291 for Mark's review. No loop-routable UNMET-clause row → **pick row 150 M5a** (iteration 275's Next). Reality check: no M5a commit on fresh origin/dev, no M5a branch or worktree, sprint JSON M5a `pending`, and the new test absent at base (`no tests to run`, which is the control for the instrument). Quorum: none owed (plan r4 routed in 266).

**Execute:** resolver `recipe codex:gpt-6.1-sol declared:provider-pin`; probe rc 0. The controller baselined the gate list on pristine `24eb2a3` outside the sandbox: `TestWorldDemoLegacyGradeDrift`, `TestWorkbenchPacketValidated`, `…EvidenceOnlyThroughRowOneOhFive`, `…Refusals` and `TestWorkbenchGradeProvenFromSubjectBoundEnvelope` all PASS. The directive carried the plan's M5a contract, one named negative control per guard, a per-subtest deadline for every table row (iteration 275's lesson), and a hard requirement to PASS with `WORLD_TOOL_AILANG_BIN` unset, since CI has no tool pin. The executor (206,304 tok, about 23 min) changed `workbench_packet.go` and added `workbench_packet_evidence.go` and `…_test.go`, with no change to `host/workbench`, packet v1 bytes or row 105. Its one declared deviation was that the plan's "source subject sibling" was not implemented: the real S5 producer emits four kind/ref pairs only, and the package descriptor pins `additionalProperties:false`. Controller re-run outside the sandbox: 7/7 named PASS (the new test has 68 subtests), PASS with the tool pin unset, vet and gofmt clean, `host/daemon` + `host/coordinator` packages ok with the pin unset (`-p 1`), and `-race` on `TestWorkbenchPacket*` under `GOMAXPROCS=2` clean (89 PASS, 0 races). No write outside the worktree in the diff. Committed `3bb1ea5`.

**Independent evaluation:** Agent `sonnet` (`agent-tool sonnet declared:alias-pin`; generator codex ≠ judge sonnet), in its own detached worktree `.eval-world-iter276`. **r1 PASS 90/100, zero blocking.** All three plan mutants (record-hash, claim-kind, limit) red at the named assertion. Of its own ~40 mutants, survivors are non-blocking: no control for a non-deadline `GetObject` error, no empty-candidate control, the zero-ResultRef (denied/failed record) accept path untested, the ATTESTED row's producer unasserted, and a few equivalent mutants. It also noted that the test roster is hand-typed to the S1/S5 shape because the seed helpers are unexported (M5c should tie it to the real producer). It adjudicated the sibling deviation as a correct reading. The controller reproduced the M5a-limit kill first-party with the tool pin unset (`input-cap-plus-one` red at `workbench_packet_evidence_test.go:472`; restore byte-identical).

**Gate 3b:** re-fetch before PR: origin/dev had moved 18 attended commits (row 114 K5, the `why` input-fields fix, row 185, and the D-WORLD-75/79/80/81 rulings). None touched the three M5a files. Rebased to `a9c0771` and re-ran: 4/4 named PASS and vet clean, then `host/daemon`, `host/coordinator` and `host/workbench` packages ok with the pin unset. PR #295 CI 3/3 on `a9c0771`. A second re-fetch found 3 more attended docs-only commits (rows 187–189 filed). `mergeable` MERGEABLE/CLEAN on that head; closing-keyword scan clean (control fired). Squash-merged `bcd2d7c`, whose code tree is identical to `a9c0771` (only the attended charter differs). The merge-push CI verdict is in the record PR.

**Mid-fire attended landings (acknowledged):** D-WORLD-80 = A, so **row 114 LANDED and clause 5 is MET, putting all 7 clauses MET**. D-WORLD-81 = A: `v1.0.0-rc.1` is tagged as a GitHub pre-release, and the build order is now 156 → 150 → 151 → `v1.0.0`. D-WORLD-75 = A (frozen-store scope). D-WORLD-79 is closed moot. Row 185 merged (PR #291). Rows 186–189 were filed. M5a was picked before D-WORLD-81 existed. It sits on that ordered path (150 after 156), so landing it reorders nothing. The next pick follows the ruling.

**Record:** no separate record judge. The record makes no new claim beyond the judged milestone and its CI.

**Routing evidence:** base=bcd2d7cd9993283d1154de52a465f263c1633e68@2026-10-16T20:24:49Z (Gate 4). Controller `claude-opus-5-5` (tok: not reported) · designer not spawned (no doc needed) · planner not spawned (plan exists, M5a ready-made) · executor `codex:gpt-6.1-sol` via recipe for resolver `recipe codex:gpt-6.1-sol declared:provider-pin` (206,304 tok) · evaluator Agent `sonnet` (`agent-tool sonnet declared:alias-pin`; 83,470 tok, one round). judge-independence: cross-vendor. Metered $0.

**Ruled out:**
- Row 185 as the pick: held by the attended row-183 trial (draft PR #291, built through World), and it merged mid-fire.
- Holding M5a after D-WORLD-81 put row 156 first: the ruling arrived after the pick, M5a is on the same ordered path, and it was already built and judged. Next follows the ruling.
- Implementing the plan's "source subject sibling" in M5a: the real producer and the package schema both reject it. PROVEN's subject is M5b's job via the envelope.
- A fix round for r1's survivors: all non-blocking, and M5b/M5c edit the same validator and tests.

**Retro:** Skill lane: none. Process lane: (1) the per-subtest deadline and pin-unset PASS requirements were written into the directive up front, and CI went green first time. Keep them. (2) A plan phrase that names a schema field the producer does not have ("source subject sibling") was caught by the executor rather than by the plan. Executor directives should keep asking executors to report such mismatches instead of inventing schema. Backlog lane: row 150 carries M5a's survivors forward to M5b/M5c.

**Progress**: World 1.0: **7 of 7 clauses MET** (clause 5 by attended row 114 this evening); `v1.0.0-rc.1` pre-released. This iteration moved no clause (goal unmoved by this landing). It advanced release row 150 to 9/17 on the path to `v1.0.0` (156 → 150 → 151).

**Next:** row 156 (`w-exec-protected-test-paths`), NEXT per D-WORLD-81 = A. No design doc exists yet, so start with the rotation designer (measure first: how the three M4 projects lay out tests, and whether a test run writes under those paths). Then row 150 M5b, which must also settle the "source subject sibling" phrase and fold in M5a's survivors.

## 277 — 2026-10-17 — row 156 `w-exec-protected-test-paths` designed (V1–V13), planned (19 milestones), M0a matrix oracle built (PR #298; judge r1 PASS 82) [PRODUCT]

**Kind:** the first product iteration on row 156, the release row D-WORLD-81 = A puts before row 150's remaining milestones. Design, quorum, plan, and one milestone (M0a, CI-only Python).

**Gate 0:** kill switch armed (`mission-world.disabled` absent); gh `sunholo-voight-kampff`; billing CLEAN. `mission_directives.sh` → 0 allowlisted directives on #246 since `2026-10-16T13:53:04Z` (34 comments); watermark moved to `2026-10-16T20:38:37Z`. `mission_gate0_self_notices.sh` rc 1 on the same 2026-10-13T02:01:58Z rc=143 notice (orphan 255, already credited by 256); control 107:4 ok, watermarks 2/2. `mission-world` unread inbox empty (`-unread -limit 500 -json`). Weekly rotation not due (#246 created 2026-10-12T12:08Z).

**Gate 1:** local dev == origin/dev `b7bae72` (Gate-1 base `b7bae72709ce46f4d3c236314d653ad2b6ff2190@2026-10-16T23:26:27Z`). Dev CI 3 check-runs, all success. Skill directory: SKILL.md and all 13 `resources/*.md` at the resolved symlink target match V1 origin/dev (14/14 `same`), and the pin copy matches too. Open PRs: #297 only (attended row-187 dogfood docs sweep, not this loop's; left alone).

**Gate 2 — clause map:** 1–7 MET (as of iteration 276). No UNMET clause, so the release order D-WORLD-81 = A governs: **pick row 156** (`[NEXT]`). Reality check: no row-156 commit beyond the attended filing (`1bf36bd`, `47cae3f`), no design doc, no branch or worktree; ledger valid, zero OPEN. The worktrees with dirty files are all consumed eval/record trees from iterations ≤ 268.

**Design:** the rotation pointer read `codex:gpt-6.1-sol` (last used), so the next entry was `pi:ollama/glm-5.3:cloud`. Its probe returned rc 0 with `stopReason: error`, `429 … Pro weekly limit`, recorded with `mission-lane-dead.sh`. `kimi-k3:cloud` shares the Ollama bucket and the OpenRouter twins are over ration, so the rotation fell to `claude:claude-opus-5-5`, which is also what the resolver printed (`recipe claude:claude-opus-5-5 declared:provider-pin`). Spawned as Agent `opus` (the operator prompt asks for Agent-tool roles; same model). It ran ~15 min (220,707 tok) and banked 12 V-rows with raw outputs under `design_docs/verification/world-row156-design-2026-10-17/`. Findings against the charter row: (1) `**/*.test.ts` protects exactly one directory level in AILANG v0.52.1 `fileguard` (depth 0 and 2 written, V7), so the grammar allows only `D/**` and base-name globs; (2) render-time expansion to literals protects existing files but 6/6 creation arms were written on macOS (V5), so prevention is backed by a per-call integrity check; (3) AILANG's deny check is lexical and a symlink alias writes into a protected path (V7; upstream ticket at landing); (4) "Linux srt has no globs" is partly stale: only `denyWrite`/`allowWrite` are literal-only on srt 0.0.78 (README). Targeted test runs in all three M4 projects wrote nothing under test paths (V3).

**Quorum:** author `claude:claude-opus-5-5`, seats `gpt6-1-sol`, `gemini-4-argon`, `oc-glm-5-3` (reserve `oc-kimi-k3`). Every round only `gemini-4-argon` was present (others unreachable), so each verdict rests on one external seat. r1 BLOCKED: `git status` against HEAD honours assume-unchanged/skip-worktree, and a moved HEAD hides changes against base. The controller measured first that episode worktrees are `git worktree add` checkouts (`host/worktree/worktree.go:57`) whose index lives outside srt's `allowWrite`; the designer's r2 then measured it (V13: sandboxed `update-index`/`commit` fail on `index.lock: Operation not permitted`) and found the reviewer's `diff-index`/`diff-files` alternative is blind to the flags too, so it added a `git ls-files -v` check (36–42 ms) and the reviewer's `head-moved` rule. r2 BLOCKED on a new, narrow point (CWD-relative pathspecs; `D/**` misses `D` replaced by a file/symlink) with a verbatim fix and no dispute of direction → applied verbatim as r3 under the narrow-refinement carve-out (§12 of the doc). Cost $0.076 metered.

**Plan:** `derive-planner-lane.sh` → `codex:gpt-6.1-sol declared:planner-lane-default-pin`; probe rc 0 (a first probe from `/tmp` was rc 1 "Not inside a trusted directory", an instrument failure, re-run from the repo). Codex recipe 2a in a detached sibling worktree (191,237 tok, ~17 min): 19 milestones ≤ 250 LOC each, 78 named negative controls, planning findings PF1–PF8 (legacy M0 matrix writes under real HOME; TOML quoting for patterns; `golang.org/x/text` needed for NFC; a git-filter preflight before the integrity check). Artifacts validated (`jq -e`, no placeholders), copied, planner worktree removed.

**Execute (M0a):** resolver `recipe codex:gpt-6.1-sol declared:provider-pin`. Directive: Python + `ci.yml` only, no git writes, one named negative control per guard, no legacy matrix or broker Go tests (they write under HOME). Executor 49,881 tok: +116 lines in `m0_matrix.py`, +1 in `ci.yml`. Controller re-ran outside the sandbox: self-test 5/5 rc 0, `--help` rc 0, AST parse ok, workflow YAML parses (ruby), and a scratch root outside the worktree is refused rc 2 with nothing created. Committed `e03aee9`; draft PR #298 pushed before judging.

**Independent evaluation:** Agent `sonnet` (`agent-tool sonnet declared:alias-pin`; generator codex ≠ judge sonnet) in its own detached worktree `.eval-world-iter277`. **r1 PASS 82/100, zero blocking** (34/40, 21/30, 18/20, 9/10). All four named mutants killed by their named cases; containment shown with HOME pointed at an empty in-worktree dir; legacy code below the dispatch byte-identical. Non-blocking survivors: a WROTE row with an unchanged host snapshot is accepted; invalid-outcome, snapshot/digest format, `root == allowed` and symlink-escape mutants survive; the self-test's own exit code is unguarded; `--help` omits the new modes. The controller reproduced the exit-code survivor first-party (with the census check also disabled, `self-test: 5 cases, 4 passed` and rc 0; file restored byte-identical). All carried to M0b in the row-156 tag.

**Routing evidence:** base=b7bae72709ce46f4d3c236314d653ad2b6ff2190@2026-10-17T00:05:18Z (Gate 4). Controller `claude-opus-5-5` (tok: not reported) · designer Agent `opus` for resolver `recipe claude:claude-opus-5-5 declared:provider-pin` after rotation entry `pi:ollama/glm-5.3:cloud` was dead (220,707 tok + r2 revision; one doc, one revision, within the diet) · planner `codex:gpt-6.1-sol` via recipe (191,237 tok) · executor `codex:gpt-6.1-sol` via recipe (49,881 tok) · evaluator Agent `sonnet` (49,614 tok, one round). judge-independence: cross-vendor. Quorum: one external seat present in both rounds. Metered $0.076 (quorum).

**Ruled out:**
- Following the rotation to `pi:ollama/glm-5.3:cloud` or `kimi-k3:cloud`: Ollama weekly limit (429); OpenRouter twins over ration.
- Treating the r1 objection's attack as live: V13 shows sandboxed code cannot write the episode index or move HEAD; the flag check stays as defence in depth.
- The reviewer's `diff-index`/`diff-files` fix: measured blind to assume-unchanged/skip-worktree, like `status`.
- A third quorum round for the r2 objection: narrow, verbatim fix, direction undisputed (carve-out).
- Running the legacy M0 matrix or the full broker package locally: both write under the real HOME (plan PF1/PF6); CI runs them.

**Retro:** Skill lane: none. Process lane: (1) the quorum ran on a single external seat in both rounds (OpenAI seat and both Ollama seats unreachable); the doc's verdicts should be quoted that way. (2) Measuring a reviewer's premise before routing the revision (the linked-worktree index location) again let the designer answer with a V-row instead of a re-assertion. Backlog lane: the upstream `fileguard` symlink ticket is owed at row 156's landing (design §10).

**Progress**: World 1.0: **7 of 7 clauses MET**; `v1.0.0-rc.1` pre-released. Goal unmoved by this iteration. Row 156 (first of D-WORLD-81's three release rows): designed, planned, 1 of 19 milestones built.

**Next:** row 156 **M0b** (extend the protected-only mode with V5's 16 arms × NONE/DIRS_ONLY/LITERAL on Linux CI, bank the real outputs as V-rows, amend §4.4/P5 if Linux differs), folding in M0a's judge survivors. Then M1a. D-WORLD-82 is open with defaults A, A, A, B.

## 278 — 2026-10-17 — row 156 M0b Linux protected-write matrix (48 rows) measured on CI, banked as V14/V15 (PR #299; judge r1 PASS 85) [PRODUCT]

**Kind:** the second product milestone on row 156 (D-WORLD-81 release order 156 → 150 → 151). A measurement milestone: CI-only Python, a Go contract test, and design V-rows written from real CI output.

**Gate 0:** kill switch armed (`mission-world.disabled` absent); gh `sunholo-voight-kampff`; billing CLEAN. `mission_directives.sh` → 0 allowlisted directives on #246 since `2026-10-17T00:34:31Z` (35 comments); watermark unchanged (no newer comment). `mission_gate0_self_notices.sh` rc 1 on the same 2026-10-13T02:01:58Z rc=143 notice (orphan 255, credited by 256); control 107:4 ok, watermarks 2/2. `mission-world` unread inbox empty. Weekly rotation not due (#246 created 2026-10-12).

**Gate 1:** local dev fast-forwarded to origin/dev `8d2ee68` (Gate-1 base `8d2ee689f1f385c87d3d1c6c632ea5956abb22ab@2026-10-17T03:26:07Z`). Dev CI 3 check-runs, all success. Skill directory: SKILL.md and all 13 `resources/*.md` at the resolved symlink target match V1 origin/dev, and the pin copy matches too (14/14). Open PRs: #297 only (attended row-187 dogfood docs sweep; left alone).

**Gate 2 — clause map:** 1–7 MET. No UNMET clause, so D-WORLD-81 = A governs: row 156 (`[IN-SPRINT]`, NEXT) at its plan's next milestone, **M0b**. Reality check: no commit or branch for M0b past `8d2ee68`; no iteration-278 worktree; ledger valid, D-WORLD-82 the only OPEN row. Ready-made plan milestone → no designer or planner.

**Execute:** resolver `recipe codex:gpt-6.1-sol declared:provider-pin`; probe from the sprint worktree rc 0. The controller baselined the directive's gates on pristine `8d2ee68` first (self-test 5/5, `go vet` rc 0, golden test PASS, all with `WORLD_TOOL_AILANG_BIN` unset). Codex recipe, 92,573 tok, ~4 min: `runProtectedMatrix` (16 V5 arms × NONE / DIRS_ONLY / LITERAL, fresh tree per row, unsandboxed host sha256 snapshots, outcome from the host as WROTE / REFUSED / PARTIAL, bounded subprocesses, `PROTECTED-JSON` line); self-test 5 → 12 named cases (the M0a judge survivors); `--help` lists the modes; `TestProtectedMatrixModeContract` (rc, census line, every named PASS); a second invocation in the existing Linux CI step against the already pinned and verified srt install (no second npm install: a plan deviation, recorded). Executor's own mutation table: 16 self-test mutants and 4 contract-test mutants killed.

**Controller re-run outside the sandbox:** self-test 12/12, `go vet` rc 0, contract test PASS, and a real macOS `--protected-only` run (rig srt 0.0.78 at the pin, scratch root inside the worktree): rc 0 in 8 s, matching V5. Real-engine control MUT-M0-NO-DIR-DENY (DIRS_ONLY without the `tests` deny): rc 1 with exactly `DIRS_ONLY/dir-modify` failing. Committed `4c3b878`, draft PR #299.

**Linux measurement (CI run 38999371033):** all 48 rows, every gate PASS. Then the job went **red** in three repo-wide Go scans (`TestRegistryDispatchBindingBoundary`, `TestProductionGoSurface`, `TestSubprocessSinkGate`): the matrix's fixture trees, with package-less `.go` files, were still under `$GITHUB_WORKSPACE/.row156-ci`. Controller fix `4c12f2c`: remove each row's tree right after its snapshot, add a `fixtures-removed` gate, and assert the directory is gone in CI; control: with the `rmtree` removed the macOS run exits 1 on `fixtures-removed` only. Re-run 39000212431 green 3/3 with **identical rows** (arm, rendering, outcome, token, host change).

**What Linux showed:** NONE wrote all 16. DIRS_ONLY refused only the directory form: create `tests/new.py` `EROFS`, `mv tests tests2` `EBUSY`, write through a new `lnk→tests` symlink `EROFS` (the link itself is created); every base-name arm written. LITERAL refused the existing literals (`EROFS` modify, `EBUSY` rm / mv away / parent rename, `EXDEV` hard link) and wrote the same six creation arms as V5. macOS with the same harness is identical outcome for outcome (`EPERM` everywhere). Design: V14 (Linux) + V15 (macOS) with `m0b_linux_ci.json` (re-banked from the green run) and `m0b_darwin.json`; §4.4's Linux column is now "measured", with a Linux-errno bullet (refusal checks from M3b on must accept `EROFS`/`EBUSY`/`EXDEV`) and the symlink row measured; P5 gains "Linux confirms". No §4.4/P5 statement was contradicted, so M3 is not blocked.

**Independent evaluation:** Agent `sonnet` (`agent-tool sonnet declared:alias-pin`; generator codex ≠ judge sonnet) in its own detached worktree `.eval-world-iter278` at `4c12f2c`, foreground, 61,570 tok, ~3 min. **r1 PASS 85/100, zero blocking** (35/40, 22/30, 19/20, 9/10). It fetched the Linux job log itself (160 KB), matched `PROTECTED-JSON` to the banked JSON and to V14, found the red's cause in the log and judged the fix, and re-ran the macOS matrix (identical to both JSONs). Its 21 mutants: 13 killed; 8 survived (validator `rc` check, before/after key equality, PARTIAL-token rule, WROTE-with-REFUSED-token, the LITERAL/modify gate, a vacuous-but-named contract case, and two it judged equivalent). The controller reproduced the `rc` survivor (12/12 with the guard removed). Non-blocking doc nits fixed in the record commit: V14 now cites the green run and LITERAL's PARTIAL symlink row, and §4.4's parent-rename row names the rendering. The remaining survivors are carried to M1a in the row-156 tag.

**Routing evidence:** base=8d2ee689f1f385c87d3d1c6c632ea5956abb22ab@2026-10-17T04:03:14Z (Gate 4). Controller `claude-opus-5-5` (tok: not reported) · designer: none (ready-made plan milestone) · planner: none · executor `codex:gpt-6.1-sol` via recipe (92,573 tok) · evaluator Agent `sonnet` (61,570 tok, one round). judge-independence: cross-vendor. Metered $0.

**Ruled out:**
- A second npm install for the new mode (plan text): the step already pins and sha256-verifies one install; a second adds time inside the 8-minute timeout and no evidence.
- Writing the Linux V-rows before CI ran: the plan forbids a planner- or executor-authored outcome; they were written from the CI JSON.
- Attributing the first red to the measurement: every matrix gate passed; the failures were later repo scans parsing the leftover fixtures.
- A fix round for the judge's survivors before landing: all non-blocking and in a measurement harness that has already produced its evidence; carried to M1a.

**Retro:** Skill lane: none. Process lane: (1) a measurement that writes fixtures inside `$GITHUB_WORKSPACE` must clean them before later steps, because this repo's Go guards scan the whole tree; ask for that in any CI-fixture directive. (2) Pushing the draft PR straight after the controller re-run let the Linux step finish in ~1 min of job time, but job logs appear only when the whole job ends (~10 min). Backlog lane: none new.

**Progress**: World 1.0: **7 of 7 clauses MET**; `v1.0.0-rc.1` pre-released. Goal unmoved by this iteration. Row 156: 2 of 19 milestones built.

**Next:** row 156 **M1a** (parse, fold, canonicalize and digest the protected grammar), folding in the M0b judge survivors. D-WORLD-82 open with defaults A, A, A, B.

## 279 — 2026-10-17 — row 156 M1a protected grammar (parse, fold, canonicalize, digest, match) built with the M0b judge carry-over (PR #300; judge r1 PASS 88 → r2 PASS 94) [PRODUCT]

**Kind:** the third product milestone on row 156 (D-WORLD-81 release order 156 → 150 → 151). Pure Go, no wiring: the grammar and digest core that M1b–M7b consume.

**Gate 0:** kill switch armed (`mission-world.disabled` absent); gh `sunholo-voight-kampff`; billing CLEAN. `mission_directives.sh` → 0 allowlisted directives on #246 since `2026-10-17T04:30:18Z` (36 comments); watermark unchanged. `mission_gate0_self_notices.sh` rc 1 on the same 2026-10-13T02:01:58Z rc=143 notice (orphan 255, credited by 256); control 107:4 ok, watermarks 2/2. `mission-world` unread inbox empty (`-unread -limit 500 -json` → null). Weekly rotation not due (#246 created 2026-10-12).

**Gate 1:** local dev = origin/dev `1a073ce` (Gate-1 base `1a073cedb49c8b64696a447bfd70c3a8a75fd374@2026-10-17T07:27:03Z`). Dev CI 3 check-runs, all success. Skill directory: SKILL.md and all 13 `resources/*.md` at the resolved symlink target, and the pin copy, match V1 origin/dev (no DRIFT line). Open PRs: #297 only (attended row-187 docs sweep; left alone — it merged mid-fire as `6731d87`, touching no charter/log file, and the sprint branch was rebased onto it before the record).

**Gate 2 — clause map:** 1–7 MET. No UNMET clause, so D-WORLD-81 = A governs: row 156 (`[IN-SPRINT]`, NEXT) at its plan's next milestone, **M1a**. Reality check: `host/broker/exec_protected.go` absent at `1a073ce`, no M1a commit or branch on origin, no iteration-279 worktree, quorum artifacts present (`w-exec-protected-test-paths-2026-10-16T23-43-50Z.json`, `…T23-46-12Z.json`), ledger valid (69 rows), D-WORLD-82 the only OPEN row. Ready-made plan milestone → no designer or planner.

**Execute:** resolver `recipe claude:claude-haiku-5-5 declared:provider-pin` (codex over daily ration, `$MISSION_ROUTING_NOTE`); run as Agent `haiku`, same weights, per the operator's standing Agent-tool request and the iteration-271/273 precedent. The controller baselined first on pristine `1a073ce` (`go vet` rc 0; matrix self-test 12/12, both with `WORLD_TOOL_AILANG_BIN` unset) and measured the facts the directive carried: `golang.org/x/text` absent from go.mod/go.sum but selected at v0.42.0 and cached; fileguard v0.52.1's `foldKey`/`foldRune` source; World's deny names (`handlers_ailang.go:57-58`). Executor (105,944 tok, ~5 min): `exec_protected.go` (`ParseProtectedPatterns`, `canonicalProtected`, `protectedDigest`, `foldProtectedKey`, `matchProtected`; reuses `execDenyNames()`), `exec_protected_test.go` (the four plan-named tests, one refusal subtest per guard), x/text v0.42.0 direct (it hand-removed two stray x/mod and x/tools go.sum lines that `go get` added). Part B, the M0b judge carry-over: four validator self-test cases (`invalid-rc`, `snapshot-keys`, `partial-token`, `wrote-refused-token`; self-test 12 → 16) and `TestProtectedMatrixModeContract` now fails instead of skipping without python3. Its mutation table: 19 mutants, all killed by the named subtests.

**Controller re-run (own tree, tool pin unset):** the five named gates PASS (11/1/1/12/1), self-test 16/16, `go vet` rc 0, gofmt clean, `go mod verify` ok. Test bytes checked with python: the NFD row really is `é`, the long-s row U+017F, the kelvin row U+212A. Golden digest re-derived with `hashlib`. Five controller mutants, all killed by their named subtest: no-NFC fold (`nfc`), ToLower fold (`long-s`), deny prefix arm removed (`already-denied`), no dedupe (golden), dir prefix without `/` (`non-hit-sibling`). Committed `b128e0f`, draft PR #300.

**Independent evaluation:** Agent `sonnet` (`agent-tool sonnet declared:alias-pin`; generator haiku ≠ judge sonnet, same vendor), in its own detached worktree `.eval-world-iter279`, foreground. **r1 PASS 88/100, zero blocking** (35/40, 24/30, 19/20, 10/10): about 45 mutants of its own; findings N1 padded entries accepted while fileguard `TrimSpace`s patterns; N2 the digest depended on which spelling of a fold-key duplicate came first (`["tests/**","TESTS/**"]` vs the reverse: two digests); N3 the dir-part glob guard, `matchProtected`'s `path.Clean` and the slash-free `**` message untested; N4 "the kelvin row is ASCII"; N5 `sub/.git/**` accepted. The controller reproduced N1, N2, N3 and N5 with a probe test, and refuted N4's premise (the row's path is U+212A; it cannot kill a ToLower fold because `ToLower(U+212A) == 'k'`, and `long-s` carries that). **Fix round** (same executor, resumed; 120,949 tok): padded entries refused (`entry-characters`), byte-smallest NFC spelling per fold key (golden recomputed outside Go: `bbe9239e…`), `dir-glob` and `clean-path` subtests plus slash-free `**` rows, a `.git` segment refused at any depth; 5/5 new mutants killed; the controller re-derived the digest and re-ran the first-seen mutant. A controller edit to design §4.1 records the canonical-spelling, padding and any-depth `.git` rules. **r2 PASS 94/100, zero blocking** (38/40, 27/30, 19/20, 10/10): 15 new-branch mutants, one survivor (ASCII-only whitespace check), and the judge withdrew N4 after dumping the bytes. Its nits N6 (non-ASCII whitespace rows) and N7 (stale golden comment) were applied by the controller (test-only, commit 3); the ASCII-only mutant now fails `TestProtectedGrammar/characters`.

**CI (PR #300):** the first run on the record head `c8c4bfd` was **red** on two jobs, only in `TestBrokerDependencyAllowlist` and `TestDaemonDependencyAllowlist/daemon_core_build_graph`: the zero-cloud allowlists rejected `golang.org/x/text/transform` and `golang.org/x/text/unicode/norm`, which the plan (PF7) mandated without naming these guards, and which the directive's gate list never ran. Controller fix (commit 5): the two package paths (exactly what `go list -deps` reaches), never the module root, in both lists with a justification comment; both guards plus their null-case and intruder subtests PASS locally; the first CI run is the negative control. Reviewed by the same judge as an r3 delta (see below).

**Routing evidence:** base=6731d87e745e9bbc7a88bdd12b5b442d9844c2cf@2026-10-17T07:54:05Z (Gate 4; Gate 1 read `1a073ce`, attended #297 merged between). Controller `claude-opus-5-5` (tok: not reported) · designer: none (ready-made plan milestone) · planner: none · executor Agent `haiku` for pin `claude:claude-haiku-5-5` (105,944 + 120,949 tok) · evaluator Agent `sonnet` (79,775 + 92,956 tok, two rounds). judge-independence: same-vendor, different model (codex lane over ration). Metered $0.

**Ruled out:**
- Judge r1 N4 ("kelvin row uses ASCII K"): refuted by a byte dump (U+212A); the judge agreed in r2.
- Refusing a base-name pattern `.git`: fileguard denies any `.git` component before the protected check, so the grammar only refuses the directory form.
- Running the unfiltered `go test -race ./host/broker/...` locally: it seeds the real HOME (plan PF6); CI runs it.
- Admitting the whole `golang.org/x/text` module to the allowlists: two package paths suffice and keep the guard narrow.
- A second executor round for N6/N7: two test rows and a comment; applied by the controller with a mutation proof.

**Retro:** Skill lane: none. Process lane: (1) asking the executor for a digest computed outside Go, with the command written in the test, paid off twice — the judge and the controller each re-derived both goldens in one line. (2) "Keep the first spelling" is order-dependent whenever a dedupe is case-insensitive; a canonical form needs a total order on spellings, not first-seen. Ask for that in any future canonicalization directive. (3) A new third-party import must name the repo's dependency-allowlist guards in the directive's gate list (rule 3g): the plan said "add x/text", and neither the plan nor the directive ran `TestBrokerDependencyAllowlist`/`TestDaemonDependencyAllowlist`, so CI found it. Backlog lane: none new.

**Progress**: World 1.0: **7 of 7 clauses MET**; `v1.0.0-rc.1` pre-released. Goal unmoved by this iteration. Row 156: 3 of 19 milestones built.

**Next:** row 156 **M1b** (render the four consumers and rebase module roots).

**Mid-fire attended activity (Gate 3b):** Mark merged PR #300 himself (`ee773d1`, 08:14Z, tree identical to the green head `b7666b0`) together with #301–#303, and resolved **D-WORLD-82 = A, A, A, A** in-session (`167f5f9`; `mission_ledger_provenance.sh` verdict `attended`). D-156-4 = A changes the plan: M5c/M5d (World's own symlink/hard-link pre-check, precedent `cliWriteFlags`) must be split and planned after M5b. D-WORLD-83 = A (attended) filed evidence rows 190–196 for after the v1.0.0 tag. A correction commit updated this entry, the 279 stamp, row 156 and the dashboard, which had been written while D-WORLD-82 was still OPEN.

## 280 — 2026-10-17 — row 156 M1b protected renderers (AILANG rebase, srt darwin/linux, git pathspecs) + pinned policy-tool fixture (PR #307; judge r1 PASS 80 → r2 PASS 84) [PRODUCT]

**Kind:** the fourth product milestone on row 156 (D-WORLD-81 release order 156 → 150 → 151). Pure Go, no wiring: the four renderings that M3a (srt), M4b/M6 (git) and M5a/M5b (AILANG policy) consume.

**Gate 0:** kill switch armed (`mission-world.disabled` absent); gh `sunholo-voight-kampff`; billing CLEAN. `mission_directives.sh` → 0 allowlisted directives on #246 since `2026-10-17T04:30:18Z` (37 comments); watermark unchanged. `mission_gate0_self_notices.sh` rc 1 on the same 2026-10-13T02:01:58Z rc=143 notice (orphan 255, credited by 256); control 107:4 ok, watermarks 2/2. `mission-world` unread inbox empty. The one unread fleet-wide item for this repo (v0.54.15 in prod, #1872–#1875) was already actioned attended in `38adaee`. Weekly rotation not due (#246 created 2026-10-12).

**Gate 1:** local dev = origin/dev `498683c` (Gate-1 base `498683c4e81f29f23a1e9885716a5a7cdbeb4595@2026-10-17T11:26:37Z`). Dev CI 3 check-runs, all success. Skill directory: SKILL.md and all 13 `resources/*.md` at the resolved symlink target, and the pin copy, match V1 origin/dev (no DRIFT line). Open PRs: none. Since record 279: attended `38adaee` (v0.54.15 notes: re-measure before M5c/M5d; rows 190/191/139) and `05efc0c` (trial ledger). Neither touches M1b.

**Gate 2 — clause map:** 1–7 MET. No UNMET clause, so D-WORLD-81 = A governs: row 156 (`[IN-SPRINT]`, NEXT) at its plan's next milestone, **M1b**. Reality check: `rebaseForSandbox`/`renderSrtProtected`/`gitPathspecs` absent at `498683c`, no M1b commit or branch on origin, no iteration-280 worktree, quorum reader rc 0 (`w-exec-protected-test-paths-2026-10-16T23-43-50Z.json`), ledger valid (70 rows, zero OPEN). Ready-made plan milestone → no designer or planner.

**Execute:** resolver `recipe claude:claude-haiku-5-5 declared:provider-pin` (codex over daily ration, `$MISSION_ROUTING_NOTE`); run as Agent `haiku`, same weights, per the operator's standing Agent-tool request (iterations 271/279 precedent). The controller baselined on pristine `498683c` with `WORLD_TOOL_AILANG_BIN` unset: `go vet` rc 0; the M1a protected tests, `TestSubprocessSinkGate` and `TestProductionGoSurface` PASS; the three new test names select zero tests. Measured facts handed to the executor: the V7 policy-tool request/policy/env shape (`q4_policytool_denywrite.py`), the refusal wording, the missing-parent refusal trap, `toolBinary(t)` at `handlers_ailang_test.go:475`, and the sink-gate pattern to avoid. Executor (110,618 tok, ~4 min) delivered the three renderers, `TestProtectedRebase` (11 subtests), `TestProtectedRenderers` (literal tables) and `TestProtectedMatchFixturePinned`. The pinned test drives the real policy-tool over the shared V7 table (moved to `protectedFixtureRows`), with empty-deny positive controls and two module-root rows. Only a `matches fs_deny_write` refusal counts as a hit. Its 9 named mutants were all killed, plus one pinned mutant.

**Controller re-run (own tree, tool pin unset):** Rebase 12 PASS, Renderers 6 PASS, MatchFixture 13 PASS, the pinned test SKIPs cleanly; focused `-race` set (incl. `TestBrokerDependencyAllowlist`, `TestRegistryDispatchBindingBoundary`) ok; sink gate, production surface and `TestDaemonDependencyAllowlist` ok; `go vet` ok; gofmt clean; `go test -run '^$' ./...` compiles. With the pin: **23 PASS, 0 SKIP**. Two controller mutants (no dedupe → `dedupe-star`; no `icase` on base names → `git`) were killed, and the file was restored byte-identical. Committed `0dfbbe5`, draft PR #307.

**Independent evaluation:** Agent `sonnet` (`agent-tool sonnet declared:alias-pin`; generator haiku ≠ judge sonnet, same vendor), in its own detached worktree `.eval-world-iter280`, foreground. **r1 PASS 80/100, zero blocking** (37/40, 22/30, 12/20, 9/10). The plan's four named controls held. It ran 15 more diff-anchored mutants and 12 were killed. It confirmed the pinned run executes (23/0/0) and cannot pass vacuously. Findings: (1) the diff is about 390 lines against the plan's 250 cap, tests included; (2) dropping the fold check in the strictly-under branch survives (`["other/tests/**"]` at root `sub` would render `tests/**`); (3) `path.Clean(moduleRoot)` untested; (4) `filepath.Join` vs concatenation in the srt renderer (only differs for a non-clean worktree). The controller reproduced (2) and (3) and fixed them in `d55b7d5`: three test rows (`disjoint-deep`, `module-dot-slash`, `module-trailing-slash`), each red under its mutant and green after restore. **r2 (resumed, delta review) PASS 84**: both survivors killed, (1) and (4) unchanged. The judge's gofmt remark on `host/store/deadline_guard_test.go` is pre-existing at base and outside this diff. Reports: `~/.ailang/state/world-iter280/evaluator-r{1,2}.md`.

**Routing evidence:** base=498683c4e81f29f23a1e9885716a5a7cdbeb4595@2026-10-17T11:26:37Z (Gate 1; no attended commits since). Controller `claude-opus-5-5` (tok: not reported) · designer: none (ready-made plan milestone) · planner: none · executor Agent `haiku` for pin `claude:claude-haiku-5-5` (110,618 tok) · evaluator Agent `sonnet` (54,279 + 57,253 tok, two rounds). judge-independence: same-vendor, different model (codex lane over ration). Metered $0.

**Ruled out:**
- Trimming the test file to meet the 250-line cap: the overrun is the pinned harness and the shared-table move, both of which the plan's own named test requires. Recorded as a deviation; the judge rated it non-blocking.
- Fixing judge finding (4): it only differs for a non-clean worktree, and `renderExecSettings` refuses non-clean paths before any renderer runs.
- Running the unfiltered `go test -race ./host/broker/...` locally: it seeds the real HOME (plan PF6); CI runs it.

**Retro:** Skill lane: none. Process lane: (1) zsh does not word-split `$E` in `E="env -u X"; $E go test`, so the first controller re-run measured nothing ("command not found"). It was caught because the PASS counts were zero. Unset the variable instead. (2) Handing the executor V7's exact policy-tool request, policy and env (measured), plus the "missing parent dir also refuses" trap, gave a pinned test that the judge could not make pass vacuously on r1. Backlog lane: none new.

**Progress**: World 1.0: **7 of 7 clauses MET**; `v1.0.0-rc.1` pre-released. Goal unmoved by this iteration. Row 156: 4 of 19 milestones built.

**Next:** row 156 **M1c** (root-anchored git pathspec equivalence with real git, including the r3 bare-`D` and subdirectory-CWD rows).

## 281 — 2026-10-19 — row 156 M1c Unicode equivalence premise REFUTED; evidence banked (PR #309 / `228645a`, independent PASS 83) [REFUTATION]

**Kind:** evidence-only refutation of the next product milestone, not completion of M1c. One row (156); no second backlog item, production change, release or semantic weakening.

**Pick and why:** all seven bar clauses are MET, so D-WORLD-81’s attended release order governs: 156 → 150 → 151 → v1.0.0. Fresh origin ended at iteration 280 / M1b; M1c was not built. Local dev was `498683c`, behind origin `a0c73053f10092fd511ad22fcf6b95d6a3f896aa`; worked in isolated sibling worktrees instead of changing the shared checkout. The attended row139 tool-pin update to v0.54.15 was accounted for, not repeated.

**Clause map:**
- 1 MET — deterministic kernel.
- 2 MET — local-first daemon.
- 3 MET — explicit authority.
- 4 MET — resident-agent floor.
- 5 MET — provenance value proof.
- 6 MET — protocol boundary.
- 7 MET — controlled self-modification.

**Gate 0:** kill switch armed; gh `sunholo-voight-kampff`; billing CLEAN. Allowlisted directives on #246: zero since the conservative older watermark `2026-10-16T13:53:04Z` (47 comments); ledger valid, 70 rows / zero OPEN. Canonical GCP mission-world inbox unread null, with cloud store header independently confirmed. Old crash notice `https://github.com/sunholo-data/ailang-world/issues/246#issuecomment-6136485418` is the orphan already credited by 256; no duplicate orphan credit. Weekly sweep: 0 orphans of 1 open issue enumerated, independent list-length control and fresh negative control passed; per-issue counts are banked. Rotation is not due before Monday 07:00 CEST; #246 had 47 comments at the check. Initial controller heartbeat used the wrong root and failed; immediately corrected to the valid absolute World driver pin, with no driver edit.

**Gates 1–2:** origin `a0c73053f10092fd511ad22fcf6b95d6a3f896aa` CI was 3/3 success. Running resolved-symlink skill and pin copy both match upstream origin on SKILL.md plus all 13 resources (14/14 each). Index checked; no open PR or M1c worktree/orphan. Checked quorum reader rc0, both retained artifacts read: each had only Gemini present, three named unreachable seats; original R1/R2 BLOCKED and existing r3 applied the reviewer’s exact top+bare-D fix under the narrow-refinement carve-out. No new design was authored or quorum silently passed. Baseline outside sandbox on darwin/arm64: vet rc0, 46 focused existing protected tests PASS; the new M1c test is absent with existing-test positive controls. Broad broker baseline remains UNKNOWN under PF6 because legacy tests seed real HOME; it was not run locally.

**Design and planning:** operator explicitly required all four Agent roles despite the ready-made plan. Designer freshness assessment found the Git/Unicode mismatch; controller freshly reproduced it before acting. Planner confirmed unchanged test-only M1c cannot pass, and that ASCII-pattern restriction cannot fix non-ASCII paths matched by ASCII patterns. This is a technical design/decomposition blocker, not a manufactured human decision: preserve semantics, then bounded designer revision, quorum and replan. Reports banked under `design_docs/verification/world-iter281/`.

**Execute:** pinned native Agent Sol 6.1 produced only `repro.py`, real `result.json` and README. Four separate disposable repos avoid APFS aliases; fake HOME/config isolation, bounded subprocesses, exact unfiltered index/status controls, ASCII selected controls and NUL-safe raw/base64 outputs. Git 2.54.0 selects `tests/x.py` and `pkg/k_test.go`, but omits `teſts/x.py` for `tests/**` and `pkg/K_test.go` for `*k_test.go`; both mirror fixture rows PASS with expected true. Controller reran the committed script: valid measurement, equivalence false, M1c incomplete. Runtime fixture paths removed. No checkout Git writes by executor. Controller annotations add V16 and mark the plan’s M1c incomplete without approving a replacement algorithm. The script’s literal pathspec source is explicitly disclosed; renderer construction was independently checked.

**Independent evaluation:** configured MiniMax probe rc0, then exact provider recipe through pinned `mission_pi_run.sh`, required session acknowledgement and isolated tree. Nested fixture Git metadata writes were denied (rc128); the model moved fixtures to `/private/tmp/claude/eval281-fixtures` against the assigned locality. Controller stopped the attempt: runner rc14 `launch_failed`, 102s, 63 tool calls, no report; outside-tree observations excluded. Sonnet4.6 and Opus5.5 `claude-sub` probes both rc1 / HTTP429 weekly quota. Dead lanes recorded for this fire. Dispatcher spawned a separate fresh Sol judge in a second clean tree at `d8e17680b642b34736ce1e57c74c73aa54c07d60`, D-WORLD-48 fallback FLAGGED, generator separate from judge. **PASS 83/100, evidence only**: independent hash-object/update-index fixtures without commits, porcelain v2 and ls-files-v reproduce both misses; exact literal Unicode recovery and health mutants confirm instruments. The two low nits: README aggregate scope wording (clarified in this record), and >50-line probe main / non-exact selected-output classifier (retained as a disclosed nonblocking maintainability limitation; stored Unicode outputs are empty). Judge broad suite/lint/CI UNMEASURED, never claimed green. Full report is banked; native usage not exposed.

**Gate 3b:** PR #309’s pull_request run `39195289460` and exact merge `228645a358d24e9bc12a891467e8cea46c935029` push run `39196307642`: expected 3/3 success on each, recorded via full-SHA bounded polls and direct per-workflow/check reads. Merge and judged trees equal `60b7ad410db7af0add81dc3f18921676eab46ef4`. Full remote gates are Ubuntu CI, not a claim of Unicode equivalence on Linux or a local full-suite green. Bookkeeping issue #246 remains OPEN. This Gate-4 record is a separate bookkeeping PR; its own required checks and merge are verified before Gate 5 reports completion.

**Gate 4:** 3 canonical live STATUS stamps, moved 278 preserved in the archive with source/destination arithmetic; dashboard 15 lines; ledger valid, zero OPEN. Queue census 92 closed / 197 rows with controls 1/79, row156 IN-SPRINT. Fleet CLI regenerated the intended record-worktree index only: 272 rows include one unnumbered attended continuation; all 271 numbered iteration IDs are present, prior IDs preserved, latest281 unique, and live/archive log hashes unchanged by regeneration. Canonical mission test, whitespace, personal-email and queue/ledger guards passed.

**Routing evidence:** base=228645a358d24e9bc12a891467e8cea46c935029@2026-10-19T04:01:27Z (Gate 4; Gate1 base was `a0c73053f10092fd511ad22fcf6b95d6a3f896aa@2026-10-19T03:26:19Z`). Controller native Sol 6.1 (tok: not reported) · designer native Agent Sol 6.1 (`recipe codex:gpt-6.1-sol declared:provider-pin`, operator transport override; tok: not reported) · planner native Agent Sol 6.1 (`recipe codex:gpt-6.1-sol declared:planner-lane-default-pin`; tok: not reported) · executor native Agent Sol 6.1 (`recipe codex:gpt-6.1-sol declared:provider-pin`; tok: not reported) · evaluator native Agent dispatcher, configured `recipe pi:openrouter/minimax/minimax-m3 declared:provider-pin` failed as above, final separate fresh native Agent Sol 6.1 (tok: not reported). judge-independence: **same-model-fresh-context FLAGGED**, D-WORLD-48; never controller self-verdict. MiniMax 1,982,634 tokens including cache/probe, metered **$0.14741610**; Anthropic probes $0. Provider usage unavailable for native Agents is unknown, not zero. All four required roles ran; no role silently omitted. Machine-readable routing and usage receipts are banked.

**Ruled out:**
- Repeating M1b from stale local dev: fresh origin has iteration 280 / PR #307.
- Calling original quorum a pass: both artifacts are blocked with one present seat; r3’s recorded verbatim remedy is the operative disposition.
- Restricting pattern bytes to ASCII: both counterexample patterns already are ASCII.
- Dropping Unicode fixture rows, changing their expected true, or calling a valid characterization M1c completion: each weakens or misstates the current contract.
- Accepting MiniMax outside-tree observations or a rc0 probe as a judge verdict: no admissible report existed; independent fresh-context fallback supplied the report.
- Asking Mark to choose merely because a faithful redesign is needed: existing semantics already authorize preservation; no OPEN ledger row was invented.
- Treating focused local tests or green CI as a completed M1c matrix: original r3, ignored/untracked and mutation requirements remain owed.

**Retro:** Skill lane: none (World does not edit shared skills/drivers). Process lane: carry fixture locality as a hard stop in provider directives; use a fresh native evaluator when sandbox transport cannot honor it. The existing instruction was violated, so this is not a routing-policy change. Scope claims explicitly distinguish executor/script actions from controller annotations. Backlog lane: the existing row156 NEXT is a bounded Unicode-preserving detection design revision and replan; no new item or human ask. Harness share over latest20 including 281: **0/20**. Last three landings (281 evidence, 280/279 product) moved no UNMET clause; all clauses are MET and no routable UNMET row exists, so no drift alarm.

**Progress**: World 1.0 **7/7 clauses MET**; **goal unmoved**. Row156 remains **4/19 milestones built**; only its failed M1c premise and reproducible evidence moved. Release remains `v1.0.0-rc.1`, with 156 → 150 → 151 before v1.0.0.

**Next:** row156 bounded designer revision preserving NFC/SimpleFold semantics; quorum and replan detection, then resume M1c. Retain top/bare-D, ignored/untracked coverage, index flags/HEAD and deadlines. DECISIONS FOR MARK: none; ledger OPEN zero. Row150 and row151 retain their attended order after156.

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
