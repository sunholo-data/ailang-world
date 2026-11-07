# World iteration 363 — mandatory Gate 4 ADMIN record plan

Date: 2026-11-07. Actual planner: `codex:gpt-6.1-sol`.
Planning checkout: `/Users/voightkampff/dev/sunholo-data/.planner-world-iter363`.
Pinned local HEAD and local origin/dev: `565609ab51c9300924831087aa98308dcdd418b0`.
Eventual record checkout supplied by supervisor: `/Users/voightkampff/.ailang/state/world-iter363-record`.
**Plan only. No milestone has been executed or accepted. Parent supervisor owns handoff.**

## Authority and scope

Read the explicitly requested sprint-planner skill and mission-control Gate 4 and false-greens resources under `/Users/voightkampff/dev/sunholo-data/ailang/.claude/skills/`, and the local charter/log at the pinned base. Read the designer assessment, controller brief and observations at these absolute paths:

- `/Users/voightkampff/.ailang/state/world-iter363-record/design_docs/verification/world-iter363/design-gate.md`
- `/Users/voightkampff/.ailang/state/world-iter363-record/design_docs/verification/world-iter363/controller-brief.json`
- `/Users/voightkampff/.ailang/state/world-iter363-record/design_docs/verification/world-iter363/controller-observations.json`

Those inputs belong to a separate tree, which this planner does not modify. This planner creates only `design_docs/verification/world-iter363/record-plan.md` and `design_docs/verification/world-iter363/record-state.json`, exclusively, without overwriting pre-existing artifacts. No coordinator/message command, agent launch, registry change or execution handoff is part of this role.

This is ADMIN bookkeeping, with no product pick, package capability, source, tests, .ail, examples or product design. Product placement/example/syntax gates and product velocity/TDD gates do not apply. Each ADMIN milestone records registry_reuse=none with a concrete reason; no package registry search is needed. Scope estimate: three sequential record tasks, approximately 1.5 hours, zero implementation/test LOC; external receipt latency is unknown. This estimate is not historical product velocity.

Row 256 / draft PR 515 remains **PARKED needs-human-review**, independent **FAIL 85**, **LOCAL_CRITERION UNMET**. Existing PR CI cannot clear that result. D-WORLD-111 default A retains the exact-two-named-baseline-failures contract, with no additional failures, source/gate change or full rerun. B remains a proposal requiring attended authorization of subsequent design/measurement. D-WORLD-112 remains OPEN: no autonomous credential change; sanitizing stored evidence does not undo provider exposure. OPEN decisions 105/110/111/112 and their defaults stay unchanged. Clearing one of 111/112 does not clear the other. There is no automatic skip to 257/201/243; the attended order remains conditional on the parks. Release rows 150/151 retain their attended prerequisites and priority when routable. There is no new source acceptance claim.

Forbidden throughout this plan's executor work: source/tests/skills/driver/secrets edits; tests or full-suite reruns; TMPDIR/GOTMPDIR overrides; environment/key dumps or credential reads; external calls/reports; git writes, commits, stash, checkout or push; modifications outside the eventual record checkout. Do not run Gate 4's general-purpose heartbeat/fetch/normalizer/test/handoff examples under this restricted role. Controller collection/landing steps below describe future acceptance dependencies, not permission for this executor to perform them.

## Baseline actually measured now

Measured locally at `2026-11-07T08:58:14.951541+00:00`, before planner writes. These are record structural observations, not product tests, CI, independent judgement or an outside-sandbox controller verdict.

| Gate | Measured result |
|---|---|
| Base/readiness | HEAD and local origin/dev both full40 `565609ab51c9300924831087aa98308dcdd418b0`; tracked diff and initial untracked list empty. No fetch performed. |
| Charter | 6,345 physical lines; 1,403,920 bytes; SHA256 `3209e7a1c245ca0172688680cf0c840086916423d9c3cc0764b58446feb9a79d`. |
| Live stamps | Exactly 3 canonical single-line stamps, each followed by one blank: 362 at line 1039, 360 at 1041, 359 at 1043. |
| Status archive controls | Bare case-insensitive ITERATION 358 token present once; 359 absent before the planned move. Existing archive is 724 lines and ends with a blank. |
| Ledger | Actual `bash scripts/mission_decisions.sh --check --file design_docs/world-mission.md`: rc0, “decision ledger valid: 99 rows”. Exact OPEN set: D-WORLD-105/110/111/112. |
| Ledger protected bytes | From its level-2 heading at line 910 up to, excluding, the next level-2 heading at line 1037: 220,998 bytes; SHA256 `302e2264c72d23ecc2f59a7322f22f90e1a9321591f9249ef99ccec8ecb1648c`. |
| Queue suffix | From the complete `## Queue (` heading at line 1807 through EOF: 992,586 bytes; SHA256 `f51d5ac1b9a8026e3ec4feb5a7dd5344cb77483c9837a23b2a64867f652190e4`. Known-present anchored rows 256 and 257 each match exactly once. |
| Canonical prior records | Entry 362 exists once in live log; canonical 360 positive control exists once. Canonical entry 361 is absent in both live and log archive; index has exactly one explicit orphan361 row, credited by 362. Do not fabricate a 361 entry or mistake the expected absence for stale state. |
| Log/dashboard | 36 canonical live entries before append; 37 prospective, below ~40 rotation threshold. Dashboard is 16 physical lines. |
| Source equality | `git diff --exit-code 565609ab51c9300924831087aa98308dcdd418b0 -- .`: rc0, empty output. This checks all tracked bytes, without exclusions at baseline. Initial `git ls-files --others --exclude-standard`: empty. |
| Whitespace | `git diff --check 565609ab51c9300924831087aa98308dcdd418b0`: rc0, empty output. |
| Artifact persistence | Both requested targets absent and not ignored (`git check-ignore -v` rc1); known ignored .ailang control fires rc0 at .gitignore:3. No staging performed. |
| Fleet CLI location | `command -v ailang` resolves to `/Users/voightkampff/go/bin/ailang`; never substitute the compiler pin `$AILANG_BIN`. |

Baseline findings, preserved rather than repaired:

1. A first exploratory `Path('missions').iterdir()` failed FileNotFoundError (command rc1): this planner checkout has no missions directory.
2. `stat -f %N /Users/voightkampff/.ailang/state/world-iter363-record/missions` failed rc1, “No such file or directory”. The eventual record checkout also lacks the required registry today. Index regeneration is **NOT RUN; PRECONDITION UNMET**. No directory, TOML, symlink or registry patch is authorized by this plan. The supervisor must establish a valid local execution context before M3; if it cannot, stop and report M3 pending.
3. The read-only fleet probe `ailang mission rotate-log --help` failed rc1: unsupported flag --help. No rotate/index command ran. The supported form below comes from the required skill and inspected fleet CLI source, not a claimed successful probe.
4. A preliminary numeric-only D-WORLD-ID regex counted 97; the validator and corrected D-identifier scan count 99, including D-WORLD-ROUTE-1 and D-WORLD-DRIVER-1. Preserve all 99; do not repair the ledger to fit a narrow matcher.

The supplied controller observations separately report base CI 3/3 success at `565609ab51c9300924831087aa98308dcdd418b0`, run [42349328842](https://github.com/sunholo-data/ailang-world/actions/runs/42349328842), and product-head CI 3/3 at `7655388936b80bea9c4843f0a2e0abc085e52c0f`, run [42344789832](https://github.com/sunholo-data/ailang-world/actions/runs/42344789832). Those are attributed historical input receipts, not fresh planner verification or iteration363 record PR/merge CI. The supplied product PR observation is draft, OPEN and CONFLICTING; this plan does not resolve it.

## ADMIN-363-M1 — dashboard and complete canonical log

Dependencies: none. Estimate: 30 minutes, approximately 90 record lines, zero code/test LOC.
Registry reuse: **none** — mission bookkeeping only; no package capability, dependency or implementation.

Paths: `design_docs/world-mission-dashboard.md`, `design_docs/world-mission-log.md` in the eventual record checkout.

1. Reconfirm the pinned baseline and prior362/orphan361 controls read-only. Baseline drift/failure is a finding to report; do not repair unrelated state.
2. Overwrite the namespaced dashboard, at most 40 physical lines. Include latest release pin AILANG v0.41.0, in-flight256 park, conditional attended order and release predicates, cadence/routing, parked-on-Mark105/110/111/112, and receipt-based quota posture. Weekly sweep/thread boundary remains 2026-11-09 07:00 Europe/Copenhagen. Leave the unnamespaced dashboard and siblings untouched.
3. Append exactly one new entry headed:
   `## 363 — 2026-11-07 — ADMIN record only; row 256 remains parked [ADMIN]`
   Preserve every pre-existing log byte as an exact prefix. Do not rewrite362 or synthesize361.
4. Populate every template section, using “none”, “not performed” or explicit pending state where appropriate:
   **Kind**; **Context / preflight**; **Pick and why**; **Gate 0**; **Gate 1**; **Work done**; **Designer / planner / executor**; **Independent evaluator (record)**; **Routing evidence**; **Metered ledger**; **Ruled out**; **Parked for human**; **Retro**; **Progress**; **Next**; **Judge follow-through (iteration 363)**; **Landing (iteration 363)**.
   Work done refers only to actual ADMIN edits. Gate0/Gate1 and release facts inherited from the brief/observations must be attributed; do not claim those commands were rerun.
5. Routing evidence must include actual, separately attributed controller/designer/planner/executor/evaluator provider:model identities, driver-declared routes, actual reroutes, per-role provider token receipts and explicit unknowns. Planner is actual `codex:gpt-6.1-sol`; its final CLI token report is collected by supervisor, currently `tok: not reported`. Controller/designer/executor declarations use Codex Sol in the supplied brief; evaluator is declared `pi:openrouter/minimax/minimax-m3`. A declaration is not evidence of a completed execution/evaluation. Preserve the supplied Anthropic quota-failure reroutes verbatim where relevant. Do not invent controller identity, role completion, token totals or costs. Separate subscription quota, metered dollars and any evidenced imputed cost; do not copy362's spend into363. No product quorum for this ADMIN work.
6. Record `base=<full40>@<UTC>` from the controller's actual Gate4 base receipt when available. Supplied Gate1 receipt is `565609ab51c9300924831087aa98308dcdd418b0@2026-11-07T08:49:44Z`; label it Gate1, never Gate4. Until controller supplies Gate4 receipt, write “Gate4 receipt pending”, not a fabricated time.
7. Ruled out must explicitly cover product merge on green CI despite FAIL85/UNMET; silent attribution-gate amendment; extra full retries/unrelated test or pin/env fixes; automatic257/201/243 pick; independent resolution of105/110/111/112; autonomous credential action; new hygiene/source/skill work; log rotation at37; treating inherited or sandbox evidence as fresh controller acceptance.

Acceptance: one canonical ADMIN363 entry, all template sections present, old log prefix exact, dashboard <=40 lines, correct product/decision defaults, honest routing/cost/token receipts and pending follow-through. No new source acceptance claim.

## ADMIN-363-M2 — safe3 stamp and independent CURRENT GOAL correction

Depends on M1. Estimate: 30 minutes, approximately 12 record lines affected, zero code/test LOC.
Registry reuse: **none** — charter record maintenance only, without reusable code or .ail syntax.

Paths: `design_docs/world-mission.md`, `design_docs/world-mission-status-archive.md`.

Use this exact canonical one-line stamp, followed by one blank:

```text
## STATUS 2026-11-07 — ITERATION 363: ADMIN only; row 256 PARKED needs-human-review, independent FAIL 85 and LOCAL_CRITERION UNMET; D-WORLD-111/112 await attended decisions; no automatic 257/201/243 route.
```

Perform two independently checked operations:

1. **Stamp-only operation.** Read bytes and split with line endings retained. Require current stamps362/360/359, each one line plus one blank; capture359's two-line bytes before any write. Prepend363 at the top of the dated-stamp area and remove only359's two lines. Never scan from a header to the next header or EOF to infer a stamp extent.
   Calculate in memory and assert **after = before + 2 - 2*moved = 6345 + 2 - 2*1 = 6345** before writing. Require exactly3 dated canonical stamps,363/362/360. If baseline differs, stop/report, without a self-heal.
   Append359's exact captured two-line bytes to the status archive, retaining its entire old prefix. Destination gains exactly2 lines (724→726). Moved stamp+blank SHA256 must be `6867bd184b9d21d03d55093efc6ab09305e7961d785aea718a210e557495258c`. Require its exact bytes in the destination, bare case-insensitive ITERATION359 token >=1, and known-present ITERATION358 control >=1. Verify the destination's resolved path is exactly under the eventual record checkout.
2. **CURRENT GOAL operation, after recording the stamp arithmetic.** Replace only the physical line currently beginning “When this line and the Queue disagree” and containing “Current next (after iteration 360)”, within CURRENT GOAL, with this one physical line:

```text
    When this line and the Queue disagree, the Queue governs; the Human Decision Ledger governs current authority. Current route (iteration 363): ADMIN record only. Row 256 remains PARKED needs-human-review, draft PR #515, independent FAIL 85 and LOCAL_CRITERION UNMET. D-WORLD-111 default A retains the exact-two-failures contract; no full reruns or source/gate changes. D-WORLD-112 stays OPEN with no autonomous credential changes; sanitization does not undo exposure. No automatic skip to 257, 201 or 243. The attended order 250 → 255 → 256 → 257 → 201 → 243 remains conditional on both parks being explicitly cleared; release rows 150/151 retain priority when routable. D-WORLD-105 and D-WORLD-110 remain OPEN with their existing defaults. History is in the log and status archive.
```

This replacement has independent prose line delta0; expected final charter line count6345. Record stamp-before, stamp-after, moved count and prose-before/prose-after separately. Do not absorb any prose delta into the stamp assertion. Preserve every charter byte outside the dated-stamp edit and this exact stale line.

Acceptance: exact3 live canonical stamps;359 survives byte-exact in the archive with358 positive control; correct arithmetic before the prose edit; current route says ADMIN/park rather than offering later rows as next. Entire ledger-heading section remains byte-exact,99 validated data rows and OPEN105/110/111/112 unchanged. Entire queue suffix equals pinned-base bytes, not merely a search hit; row256/257 controls still fire. Inspect diff/stat before considering this milestone complete.

## ADMIN-363-M3 — pinned fleet index and record acceptance evidence

Depends on M1/M2. Estimate: 30 minutes plus unknown controller/evaluator/CI latency, approximately40 evidence/index lines, zero code/test LOC.
Registry reuse: **none** — existing fleet mission CLI is reused as a bookkeeping tool, without creating a package or changing a registry.

Paths: `design_docs/world-mission-index.md` and only new curated evidence under `design_docs/verification/world-iter363/`.

**Unmet execution precondition:** a valid existing `<eventual-record-root>/missions` registry must be supplied by the supervisor. The current absolute candidate is absent. Validate its resolved origin and World output target read-only before invoking a writer. CLI source inspected at `/Users/voightkampff/dev/sunholo-data/ailang/cmd/ailang/mission_rotate.go:63` uses registry root for shared-repo missions but configured World Workdir otherwise. Therefore the explicit registry pin alone is insufficient: World must resolve its output Workdir to the intended record checkout too. If either target differs, STOP/report before invoking the CLI. Do not create/repair registry files, patch the driver, or fall back to a main/pin checkout. If supervisor changes the eventual record checkout, use that absolute root consistently and reconfirm baseline; never substitute this planner checkout silently.

Once that precondition is met, the future executor runs, from the confirmed record root:

```sh
record_root=/Users/voightkampff/.ailang/state/world-iter363-record
AILANG_MISSION_REGISTRY="$record_root/missions" ailang mission rotate-log world --index-only
```

Use fleet `ailang` on PATH, whose measured planner location is `/Users/voightkampff/go/bin/ailang`; do not use `$AILANG_BIN`. Never combine --index-only and --keep. There are36 existing live entries and37 after append; no full-log rotation is due. The log archive `design_docs/world-mission-log-archive.md` stays byte-identical to base.

Capture rc/stdout/stderr in new curated `design_docs/verification/world-iter363/record-checks.json` only when execution occurs; no raw provider transcripts, environment dumps or credentials. An existing evidence path must not be overwritten. Require emitted INDEX path to resolve exactly to `<record-root>/design_docs/world-mission-index.md`. Index-only normally emits no archive path: record that fact rather than inventing one. Resolve/check both expected archive destinations (`world-mission-log-archive.md`, `world-mission-status-archive.md`) under that same record root; if the CLI emits any archive path, assert that path too, and stop on unexpected rotation. Prevalidation is mandatory because a post-write path assertion cannot undo a wrong-root write.

Acceptance: regenerated index contains exactly one363 row, preserves exact orphan361 row and all prior numbered iterations, and reflects new ADMIN title; both source streams remain unchanged by index-only. No log rotation. Run the structural gates below; any red baseline is reported without unrelated repairs.

## Structural validation procedure for the future executor/controller

These are prospective post-edit checks, not tests run now. Read pinned-base bytes with `git -C <record-root> show 565609ab51c9300924831087aa98308dcdd418b0:<path>`; do not rely on snapshots written over real artifacts.

- Canonical new heading regex: `^## STATUS [0-9]{4}-[0-9]{2}-[0-9]{2} — ITERATION [0-9]+: [^[:space:]]`. Validate exactly3 dated live headings and each following blank. Validate canonical log363 once,362 once; preserve index361 orphan once. Never require a nonexistent canonical361 entry.
- Compare ledger bytes from the exact line-anchored Human Decision Ledger heading through, excluding, its next level2 heading. Compare the full bytes to base and independently SHA256/length against the measured values above. Run `bash scripts/mission_decisions.sh --check --file design_docs/world-mission.md` and require99. Read OPEN IDs structurally without printing decision contents or environment.
- Compare full queue suffix from its exact line-anchored level2 Queue heading through EOF to base bytes; require identical bytes, hash and length. Then verify known row256/257 anchored positive controls. This protects the Premise Verification Log and trailing rows too.
- Verify old log is an exact prefix of final log and only one363 entry was appended. Verify status archive old prefix and exact moved359 bytes, line delta+2, archive358 positive control, and that no independent prose was archived. Verify log archive remains wholly unchanged.
- Validate dashboard <=40 physical lines. Review diff/stat for only the intended stamp, one CURRENT GOAL line, log append, dashboard overwrite, index regeneration and new evidence.
- Run `git diff --check 565609ab51c9300924831087aa98308dcdd418b0` and require rc0.
- Require source equality with **only** these explicit mission-record exclusions:
  - `design_docs/world-mission-dashboard.md`
  - `design_docs/world-mission-log.md`
  - `design_docs/world-mission.md`
  - `design_docs/world-mission-status-archive.md`
  - `design_docs/world-mission-index.md`
  Additionally allow **new** evidence files under `design_docs/verification/world-iter363/`, only where no such file existed at the pinned base. Excluding all design_docs or verification is forbidden. Existing evidence outside that exact iteration and every source/test/skill/driver/secret path must have zero diff.

Concrete tracked-byte source gate (run read-only in the record root):

```sh
git diff --exit-code 565609ab51c9300924831087aa98308dcdd418b0 -- . \
  ':(exclude)design_docs/world-mission-dashboard.md' \
  ':(exclude)design_docs/world-mission-log.md' \
  ':(exclude)design_docs/world-mission.md' \
  ':(exclude)design_docs/world-mission-status-archive.md' \
  ':(exclude)design_docs/world-mission-index.md' \
  ':(exclude)design_docs/verification/world-iter363/**'
```

Pair that with `git ls-files --others --exclude-standard -z` and a path audit: every untracked file must be new iteration363 evidence; no new directory/path elsewhere is allowed. Independently enumerate all changed paths relative to base: for each excluded iteration363 path, `git cat-file -e <base>:<path>` must show it did not exist at base; otherwise exclusion would hide a modified pre-existing artifact and must fail. Use `git ls-files --others --ignored --exclude-standard -z` as well to detect any newly ignored artifact rather than treating it as source equality; initial planner tree was clean. Do not force-add ignored artifacts. This combination must yield zero differences everywhere else, including deletions, staged changes and untracked source. It is not enough to compare only .go or only source directories.

The expected post-plan planner-only state is exactly the two new requested artifacts and zero tracked changes. Final controller checks repeat the same preservation assertions **outside the sandbox**, inspecting any potential outside-tree writer before running it. Observe protected state without seeded writes to real HOME/logs. Socket results inside workspace-write are **UNINFORMATIVE UNDER SANDBOX**; no socket or product tests are authorized for this ADMIN task.

## Prospective independent/controller/CI acceptance

All of the following remain pending/unavailable until the controller collects actual receipts:

1. Controller outside-sandbox structural check result and path to its commands/results.
2. Actual independent record evaluator **PASS** on the final reviewed record bytes before any record PR merge. A product FAIL85 and prior record362 PASS98 remain separate historical verdicts. Do not manufacture363 evaluator/model/protocol/token/verdict receipts. Material later edits require evaluation coverage of the resulting record.
3. Exact full40 **record PR-head SHA**, all3 completed/success GitHub checks (go host build + test gate, go race gate, ailang-code verify gate), and run URLs bound to that SHA. No moving ref, base run, product PR run, skipped/pending check or prior head is sufficient.
4. After controller-owned landing, exact full40 **record merge-head SHA**, the same3 completed/success checks and run URLs bound to that merge SHA. Until collected, landing/merge CI is pending; never report LANDED from PR CI alone.
5. Actual Gate4 base@UTC and per-role CLI token/cost receipts join the log when supplied. Supervisor captures this planner's real terminal token report. Unknown is explicit, never zero.

This planner does not run tests, collect external CI, launch an evaluator, merge a PR, report externally, change credentials or perform handoff. If missing registry/baseline/receipt conditions remain, report the exact gate as pending or failed and stop; do not expand the scope.

