# Sprint plan — w-exec-protected-test-paths

**Iteration 282 · queue row 156 · detached base `451dc3d` · affected milestones planned · plan only**

Approved design: [w-exec-protected-test-paths.md](w-exec-protected-test-paths.md), including [r4], [r4.1], [r4.2] and quorum §12. Iteration 281 refuted plain Git pathspec equivalence (V16); Git now generates widened candidates and World's mirror filters them. Quorum round 2 was BLOCKED on the phase-budget surface; the controller applied the reviewer-authored fixes under the narrow-refinement carve-out, without a third-pass claim. M1c/M1d may proceed; M6a waits for measured V24 and resized bounds. Unicode semantics are unchanged. JSON retains the existing schema and done milestone objects unchanged; M0a/M0b/M1a/M1b are done in JSON, their original Markdown implementation sections below are historical planning records. No completion evidence is fabricated.

21 executor/controller runs, 590 minutes; 4757 changed LOC total including production/tests/docs. <=250 LOC and <=30min each; CI/measurement waiting separate from active work. M1c <=140 production +110 tests; M1d <=240 tests. Every run has a separate reviewable PR/evidence boundary; normal CI must be green before the next dependency. Total changed LOC includes production, tests and docs. The planner edits only this file and [sprint_w-exec-protected-test-paths.json](sprint_w-exec-protected-test-paths.json); no Git writes, fixture repos or network. Existing default-decision/precheck material outside the affected detection milestones is retained, not re-ratified by this replan.

## Decisions and boundaries

- **D-156-1 = A while unanswered**; affects M4b, M6a, M6b, M6c, M7a, M7b. Report changed paths/error without failing or refusing command or next call.
- **D-156-2 = A while unanswered**; affects M1a, M2a, M3a, M4c, M5a, M5b, M6b, M7a, M7b. No defaults; absent/[] means exact old policy/settings/result bytes, why none declared.
- **D-156-3 = A while unanswered**; affects M0a, M0b, M1b, M3a, M3b, M4c, M6b, M7b. Linux basename detected only, no literal expansion production path, enforced dirs-only.
- **D-156-4 = B while unanswered**; affects M5a, M5b, M6c, M7b. Policy half only. Do not implement precheck on unanswered decision. If A arrives, replan separate <=250LOC/30min M5c/M5d after M5b, before later integration: resolution guard killed by MUT-PRECHECK-LEXICAL-ONLY via TestWorkspaceWriteProtectedPrecheck/symlink-dir, and Nlink guard killed by MUT-NO-NLINK via TestWorkspaceWriteProtectedPrecheck/hardlink; lexical/file-alias rows and spawn counter0. Existing anchors handlers_ailang.go:203 resolveExistingPrefix, :406 Execute, :242 AilangToolConfig. No current milestone/JSON entry for this conditional unapproved work.

No `.ail` changes. This is host/operator configuration across external enforcement engines, not a package/kernel law (§4.7–4.8); no `world/`, `se-tools`, InterfaceHash or wire-version change. Therefore **no milestone runs `./scripts/verify_ail.sh`**. A later `.ail` scope addition requires S5 fluency/pin validation, that exact gate and a fresh LOC/time split.

## Environment and honest gates

- WORLD_TOOL_AILANG_BIN MUST be unset for every normal gate (CI does not provision it); optional real v0.52.1 arm skips cleanly when unset.
- No .ail changes: do NOT run scripts/verify_ail.sh for this plan. If scope later changes to .ail, first perform S5 reference/builtin/pin checks and add ./scripts/verify_ail.sh to that milestone; re-estimate/split, do not assume it compiles.
- Real srt uses WORLD_EXEC_SRT_NODE_MODULES and WORLD_EXEC_NODE, pin @anthropic-ai/sandbox-runtime@0.0.78 / cli sha256 3c3092bd26b3924046f38d793c716b513dde619cf792ec50b181e2c7cd40d96e. Linux also requires bwrap; missing dependency skips before fixture setup; configured corrupt pin fails.
- Git and Python3 are CI prerequisites, not optional skips for the real pathspec/matrix companion tests.
- Observe hashes before/after for existing mission/core/driver/state guards; absence stays absence. Never seed or alter a real guard to prove it unchanged. Fake HOME means a fixture subdirectory in the worktree, never the actual $HOME or ~/.ailang/state. Set TMPDIR/GOTMPDIR to worktree-owned scratch so t.TempDir stays contained. For any full legacy package suite, require a contained runner first (PF6); otherwise keep its gate UNMEASURED, never silently omit it or report green.
- Every daemon/REST t.Run owns context.WithTimeout inside its closure plus cancellation/cleanup. Never one shared context across subtests, even if they currently finish quickly.

Exact package gates below are run with `WORLD_TOOL_AILANG_BIN` **unset**; do not copy the marketing plan’s export of that variable. The command strings assume that environment. Separate optional provisioned runs may set it to a read-only v0.52.1 path. Missing real engines skip by name before setup; a configured corrupt pin fails. Require all nonoptional named tests discovered and PASS, never a zero-selection/skip as acceptance. The real-engine arm alone may SKIP; its always-on companion must PASS.

**Full package baseline limitation:** required unfiltered package race gates include existing off-worktree probe controls. This planner did not run them and cannot call them green. A controller-provided contained runner is a prerequisite for executing those exact gates under the user’s worktree-only rule. Do not seed a real HOME/state guard, edit unrelated tests, omit mandatory gates, or turn a sandbox denial into a pass. Normal CI must still be green for every PR. This is a reported prerequisite, not an invented baseline result. All new named acceptance fixtures and mutation stores are worktree-local; real guards are hash-before/hash-after observations only.

## Measured baseline at 451dc3d (iteration 282)

| Command | Result |
|---|---|
| `go vet ./host/broker/` | GREEN, rc 0; 0.53 s tool wall time |
| `go test -count=1 ./host/broker/ -run '^TestProtected' -timeout 300s` | GREEN, rc 0; package time 0.495 s; optional pinned arm skips with binary unset |
| `jq -e . design_docs/planned/sprint_w-exec-protected-test-paths.json` | GREEN, rc 0; JSON structural validation |
| `go test -race -count=1 ./host/broker/...` | NOT MEASURED (PF6: legacy tests write real HOME); no unfiltered broker run |
| V24 real installed-worktree measurements | NOT MEASURED; controller/executor hard prerequisite before M6a |

Both Go gates used `env -u WORLD_TOOL_AILANG_BIN` and TMPDIR/GOTMPDIR under `.scratch-planner282/`, removed after checks. No socket checks, Git fixture repositories or Git writes ran. Socket-touching checks are **UNINFORMATIVE UNDER SANDBOX**. Candidate tests are **absent at base; zero tests is not a pass**. Existing `TestProtectedRenderers/git` passed the old golden only. Exact new/changed standalone gates and all mutation drills are unmeasured unless explicitly listed GREEN here.

## Historical iteration-277 baseline (d919701; not re-measured at 451dc3d)

The following inherited measurements and per-gate historical baseline statements in untouched sections refer to d919701, never the current base. In particular, repository-wide `go vet ./...` and old focused race selectors were NOT MEASURED in iteration 282.

| Command | Result |
|---|---|
| `go vet ./...` | GREEN with WORLD_TOOL_AILANG_BIN unset; rc 0, 1.46 s |
| `go test -race -count=1 ./host/broker/... -run '^(TestRenderExecSettings.*\|TestWriteAndVerifyExecSettings\|TestParseExecProfile.*\|TestRenderEpisodePolicyHasExactlyTheDesignKeySet\|TestQuickstartSection10ProfilesLoad)$' -timeout 90s` | GREEN focused broker race baseline (no sockets); not a full package pass; rc 0, 4.52 s |
| `go test -count=1 ./host/broker/ -run '^(TestRenderExecSettings.*\|TestWriteAndVerifyExecSettings\|TestParseExecProfile.*\|TestRenderEpisodePolicyHasExactlyTheDesignKeySet\|TestQuickstartSection10ProfilesLoad)$' -timeout 90s -v` | Named existing tests PASS, no pin, no sockets; rc 0, 2.1 s |
| `go test -race -count=1 ./host/broker/...` | UNKNOWN; NOT RUN: legacy off-worktree writes and/or sockets, not safe within planner scope. Any socket result is UNINFORMATIVE UNDER SANDBOX. |
| `go test -race -count=1 ./host/daemon/...` | UNKNOWN; NOT RUN: legacy off-worktree writes and/or sockets, not safe within planner scope. Any socket result is UNINFORMATIVE UNDER SANDBOX. |
| `go test -race -count=1 ./cmd/ailang-worldd/...` | UNKNOWN; NOT RUN: legacy off-worktree writes and/or sockets, not safe within planner scope. Any socket result is UNINFORMATIVE UNDER SANDBOX. |

The measured commands used worktree-local temporary directories and unset `WORLD_TOOL_AILANG_BIN`; fixtures were removed after each command. No socket test ran. **Any result touching sockets is UNINFORMATIVE UNDER SANDBOX**, never a reported pass/fail. Each milestone repeats baseline state for each literal gate; future tests are absent at base, not green. Linux/bubblewrap results can be measured only in CI M0.

## Verified paths and planning findings

Existing paths below were checked with `ls`, function `rg -n`, and `Path.is_file()`. Future source/test files are explicitly **NEW/absent at base** with existing parent directories verified; no fictitious base line is assigned to a future function. Function/line references are from this checkout, not the design’s older `b7bae72` anchors. Read-only fixture reuse does not authorize edits to the referenced helper file.

**PF1 — M0 extend existing matrix step.** The current matrix creates BASE under real HOME, seeds a HOME decoy and writes /tmp/claude; the existing CI call does not pass --json. It cannot be invoked verbatim as a worktree-only acceptance drill. Command: `rg -n 'HOME|mkdtemp|DECOYS|m0-probe|--json' design_docs/verification/world-row140-m0/m0_matrix.py; sed -n '195,216p' .github/workflows/ci.yml`. Disposition: M0a/b add early contained protected-only mode and explicit JSON output, preserving legacy row140 measurement separately. Do not run legacy mode locally.

**PF2 — Each milestone <=250 LOC /30 min (body estimates exclude tests).** M4 needs an actual integrity checker before M6 introduces it, and M5 changes a required public signature referenced in eight files, not only its listed implementation files. The stated tiny per-stage counts are not safe total changed-LOC budgets. Command: `rg -n 'RenderEpisodePolicy\(|RenderRunVariant\(' --glob '*.go'; sed -n '189,207p' design_docs/planned/w-exec-protected-test-paths.md`. Disposition: 19 bounded runs, including tests/docs; shared checker introduced M4b, later result wiring M6. All signature caller files listed.

**PF3 — V5 16 arms portable to Linux; NONE all arms written.** The banked macOS script creates pkg, never PKG, then appends PKG/A_TEST.GO and reads back pkg/a_test.go. On case-sensitive Linux that control would fail with a missing parent or inspect the wrong file. This is a fixture portability defect, not a measured Linux confinement result. Command: `sed -n '12,17p' design_docs/verification/world-row156-design-2026-10-17/q3_srt_denywrite.py; sed -n '29,46p' design_docs/verification/world-row156-design-2026-10-17/q3_srt_denywrite.py`. Disposition: M0b creates all parents and records actual case-spelled host targets; keep original macOS evidence unchanged.

**PF4 — execDecode learns the new field.** execDecode already decodes directly to execResult with DisallowUnknownFields. Adding a field to execResult automatically updates its schema; no separate decoder field exists. Command: `sed -n '162,171p' host/broker/handlers_exec_m2_test.go`. Disposition: M6b adds a compatibility assertion; do not invent a duplicate schema edit.

**PF5 — Grammar accepted strings can be appended directly to policy.** tomlStrings currently wraps strings in quotes without escaping; the new basename grammar permits a quote character. Direct interpolation would produce invalid TOML. This is a coupling omitted by the design, not justification to narrow the grammar silently. Command: `sed -n '211,223p' host/broker/handlers_ailang_run.go; sed -n '123,144p' design_docs/planned/w-exec-protected-test-paths.md`. Disposition: M5a safe TOML string encoding with exact legacy byte compatibility and quote negative control.

**PF6 — Broad package tests are safe to run with the tool pin unset.** TestExecSrtConfinementMatrix seeds real HOME before calling the helper that skips absent srt; shim probe code also uses /tmp/claude controls. Unsetting WORLD_TOOL_AILANG_BIN does not make an unfiltered broker baseline a read-only or contained action. Command: `sed -n '264,282p' host/broker/handlers_exec_srt_test.go; sed -n '365,381p' host/broker/exec_probe.go`. Disposition: Full mandatory package baseline is UNKNOWN here, not green/red. Do not execute without a controller-provided contained runner; do not add new acceptance cases that seed these paths. This runner prerequisite is reported openly.

**PF7 — NFC pure matcher needs no dependency file changes.** No direct x/text import or checksum exists in the base requirements, but Unicode NFC needs golang.org/x/text/unicode/norm. The module graph already selects v0.42.0 through ailang v0.47.2. Command: `cat go.mod; rg -n golang.org/x/text go.sum; go list -mod=readonly -m golang.org/x/text`. Disposition: M1a adds existing-version direct requirement and checksums; include go.mod/go.sum in its file allowlist and total LOC. No graph-version upgrade.

**PF8 — Two git invocations plus gitEnv are sufficient for no filters.** gitEnv disables global/system configuration, not local filter.<name>.clean/process. Existing gitSnapshot calls configuredGitFilters and refuses before status, a third invocation. The design's two-call description omits this mandatory host execution guard. Command: `sed -n '274,290p' host/broker/worktree_delta.go; sed -n '297,355p' host/broker/worktree_delta.go`. Disposition [r4.x]: retain M4a preflight; M4b includes query plus THREE candidate data calls in one provisional 2500ms phase bound. V24 measures all overhead and resizes before M6a; no host filter execution allowed.

## M0 exact CI change

Current `.github/workflows/ci.yml:118` job is **go-verify**, `runs-on: ubuntu-latest`; there is no strategy.matrix in this job. Its “Row 140 srt confinement matrix (linux bubblewrap)” step (`:195`, timeout eight minutes) installs bubblewrap/socat/ripgrep, adjusts the userns sysctl, reads version+SHA from `world-row140-m0/pin.json`, npm-installs srt under `$RUNNER_TEMP/srt`, verifies cli.js, and invokes (`:210`):

```sh
python3 design_docs/verification/world-row140-m0/m0_matrix.py \
  --srt-node-modules "$RUNNER_TEMP/srt/node_modules" \
  --node "$(command -v node)" --seccomp-read allowread
```

M0a/b retain that legacy invocation and its later WORLD_EXEC_* exports. They add an early **protected-only** dispatch in the same Python file, leaving legacy mode separate, and a second invocation in the same Linux step. For the added path, npm install root and npm cache are `$GITHUB_WORKSPACE/.row156-ci/srt` and `.row156-ci/npm-cache`; use the same version/SHA assertion, no new pin or tool binary. Add worktree-local scratch/output, a self-test invocation and the command below. Keep the eight-minute timeout initially; measure actual duration and split the step before increasing it if necessary. No matrix job expansion, no tool-pin provisioning, no unrelated workflow edits. Existing machine provisioning is not a new acceptance drill.

```sh
python3 design_docs/verification/world-row140-m0/m0_matrix.py --protected-self-test
python3 design_docs/verification/world-row140-m0/m0_matrix.py --srt-node-modules "$GITHUB_WORKSPACE/.row156-ci/srt/node_modules" --node "$(command -v node)" --seccomp-read allowread --protected-only --scratch-root "$GITHUB_WORKSPACE/.row156-ci/matrix" --json "$GITHUB_WORKSPACE/.row156-ci/protected-linux.json"
```

The new mode exits before legacy HOME/decoy/archive setup. Test 16 operation arms × NONE/DIRS_ONLY/LITERAL = 48 independently snapshotted rows: ordinary control, modify, new nested/root tests, new directory test, conftest, TS test, delete, rename onto/away, protected-dir rename, parent-dir rename, symlink, hardlink, existing/new case variants. NONE must show all intended operations written, with actual case-sensitive targets and all parents present. Record tokens, rc and host snapshots independently; lack of a client token cannot masquerade as refusal. Real Linux output becomes new V rows in the design, never a planner-authored outcome. Disagreement with §4.4/P5 blocks M3 until those statements are amended with evidence.

## Milestones

| ID | Title | Changed LOC total | Active minutes | Dependency |
|---|---|---:|---:|---|
| M0a | Contain the new Linux measurement mode and validate its oracle | 210 | 25 | — |
| M0b | Measure all protected-write arms on Linux and bank evidence | 245 | 30 | M0a |
| M1a | Parse, fold, canonicalize and digest the protected grammar | 242 | 30 | M0b |
| M1b | Render the four consumers and rebase module roots | 230 | 25 | M1a |
| M1c | Pure candidate renderer and mirror filter | 250 | 30 | M1b |
| M1d | Real-Git candidate and mirror equivalence | 240 | 30 | M1c |
| M2a | Load the additive protected profile key | 200 | 25 | M1d |
| M2b | Add startup grammar refusal rows with matching usage table | 180 | 25 | M2a |
| M3a | Append protected sandbox settings without changing old bytes | 210 | 25 | M2b |
| M3b | Exercise real srt protected writes with pin-free companions | 235 | 30 | M3a |
| M4a | Refuse configured git filters before any integrity status | 170 | 20 | M3b |
| M4b | Build the bounded integrity checker needed by probe arm 8 | 245 | 30 | M4a |
| M4c | Wire startup arm 8 for prevention and Linux detection | 245 | 30 | M4b |
| M5a | Render protected AILANG policies in every variant | 245 | 30 | M4c |
| M5b | Wire episode profile selection and module-root policy rebasing | 235 | 30 | M5a |
| M6-pre | Measure V24 and size phase/construction bounds | 200 | 30 | M5b, M1d |
| M6a | Capture the immutable episode baseline and bounded record shape | 245 | 30 | M6-pre |
| M6b | Attach before/after integrity while preserving the command deadline | 245 | 30 | M6a |
| M6c | Kill tamper, ignored-file, flag and moved-HEAD mutants; pin replay | 245 | 30 | M6b |
| M7a | Show intact, changed, unchecked and unprotected runs in why | 230 | 30 | M6c |
| M7b | Ship executable protected-profile guidance and residuals | 210 | 25 | M7a |

### M0a — Contain the new Linux measurement mode and validate its oracle

Design M0; **210 LOC total / 25 min**, own PR, depends on nothing. Status **planned**.

**Files touched and base functions:**

- `design_docs/verification/world-row140-m0/m0_matrix.py:18 argument parser; :28 HOME; :32 rec; :37 BASE; :102 settings; :123 srt; :146 unsbx; :371 summary (existing top-level blocks, no main function)`
- `.github/workflows/ci.yml:118 go-verify; :195 Row 140 srt confinement matrix; :210 existing invocation`

Allowlist: `design_docs/verification/world-row140-m0/m0_matrix.py`, `.github/workflows/ci.yml`.

**Implementation and landable boundary.** Add an early --protected-only/--protected-self-test dispatch before legacy HOME, BASE, decoy and archive side effects. Mode flags are parsed before enforcing legacy required node/module arguments, so --protected-self-test needs no binary. New runProtectedMatrix and validateProtectedRows use an explicit canonical --scratch-root contained by GITHUB_WORKSPACE; create a synthetic home/state/cache beneath it. Self-test uses in-memory tokens and no real srt. Keep legacy row-140 invocation and pin unchanged. Prepare worktree-local npm install/cache for the added invocation; no production Go code. All returned rows require a nonempty exact arm/rendering census, expected positive operation token and independent host snapshot, not rc alone.

M0 runs on Linux CI only; the Python oracle itself has no pin. Legacy M0 is not a local gate.

**Named tests and assertions:**

- `ProtectedMatrixSelfTest`: In-memory exact census, empty result, missing client token, forged host readback and out-of-root scratch controls. No archive/install or network.

**One named negative control per guard (independent mutation runs):**

- `ProtectedMatrixSelfTest/empty` kills **MUT-M0-EMPTY** (Remove exact nonempty row-census rejection); assertion: empty matrix rejected.
- `ProtectedMatrixSelfTest/missing-client` kills **MUT-M0-RC-ONLY** (Accept nonzero rc without a positive client token); assertion: missing client rejected.
- `ProtectedMatrixSelfTest/outside-root` kills **MUT-M0-ROOT-ESCAPE** (Permit a scratch root outside worktree); assertion: root escape rejected before any mkdir or write.

**Exact gates** (tool pin unset; mandatory package scopes below):

```sh
python3 design_docs/verification/world-row140-m0/m0_matrix.py --protected-self-test
go test -race -count=1 ./host/broker/...
go vet ./...
```

**Baseline for every gate:**

- `python3 design_docs/verification/world-row140-m0/m0_matrix.py --protected-self-test` — NOT MEASURED / not present at base; no green claim; Linux CI only for matrix, unfiltered broker unsafe locally.
- `go test -race -count=1 ./host/broker/...` — NOT MEASURED / not present at base; no green claim; Linux CI only for matrix, unfiltered broker unsafe locally.
- `go vet ./...` — GREEN: rc0, 1.46s

No `.ail` change; no AIL verification command. New fixture writes are worktree-local; no seeded external state. Any daemon/REST table deadline is created inside its own `t.Run`.

### M0b — Measure all protected-write arms on Linux and bank evidence

Design M0; **245 LOC total / 30 min**, own PR, depends on M0a. Status **planned**.

**Files touched and base functions:**

- `design_docs/verification/world-row140-m0/m0_matrix.py:102 settings, :123 srt, :134 readp, :137 writep, :146 unsbx, :371 summary`
- `.github/workflows/ci.yml:195 step; :203 pin; :210 invocation`
- `design_docs/planned/w-exec-protected-test-paths.md:100 Verification Log, :185 enforcement table, :454 Quorum log`
- `host/broker/exec_settings_test.go:129 TestRenderExecSettingsGolden`

Allowlist: `design_docs/verification/world-row140-m0/m0_matrix.py`, `.github/workflows/ci.yml`, `design_docs/planned/w-exec-protected-test-paths.md`, `host/broker/exec_settings_test.go`.

**Implementation and landable boundary.** Add the 16 V5 arms from q3_srt_denywrite.py under NONE, DIRS_ONLY and LITERAL (48 observations, fresh tree per row). Make all parent dirs exist in controls, including PKG on case-sensitive Linux; inspect the actual case-spelled target, not the macOS folded alias. Bound every subprocess. Add ProtectedMatrixSelfTest/none-not-written and TestProtectedMatrixModeContract (Python self-test wrapper, no pin). Add the second protected-only invocation to the existing Linux step, require all row tokens/census and write worktree-local JSON. Record actual CI output as new V rows in the design; revise §4.4/P5 before M3 if observations disagree. No fabricated Linux result.

Existing step is go-verify on ubuntu-latest, NOT a GitHub strategy.matrix job. Linux matrix means the Python arm table. Retain apt/sysctl provisioning and old invocation; added npm install/cache/output paths are under GITHUB_WORKSPACE/.row156-ci. New mode never executes legacy HOME or /tmp/claude arms. Self-tests do not replace the real 48-row measurement. Preserve INFO observations separately from enforcement promises; NONE all-written and expected protected arms are gating.

**Named tests and assertions:**

- `TestProtectedMatrixModeContract`: Invoke --protected-self-test; require named positive and negative case census, not zero tests.
- `ProtectedMatrixLinux`: Actual Linux srt run discovers all48 rows; NONE written controls, DIRS_ONLY prevention and LITERAL limits use independent host readback. This is not replaced by synthetic self-test outcomes.

**One named negative control per guard (independent mutation runs):**

- `ProtectedMatrixSelfTest/none-not-written` kills **MUT-M0-NONE-REFUSES** (Force one NONE rendering operation to refuse); assertion: NONE control requires every intended operation written.
- `ProtectedMatrixLinux/DIRS_ONLY/dir-modify` kills **MUT-M0-NO-DIR-DENY** (Omit directory deny from Linux candidate settings); assertion: actual directory-form modify emits WROTE/changes host bytes instead of refused+unchanged, so the real Linux matrix fails. Its command is the protected-only Linux invocation below, not a synthetic self-test.
- `ProtectedMatrixSelfTest/host-changed` kills **MUT-M0-FORGED-REFUSAL** (Report a refusal token while host bytes changed); assertion: host snapshot must remain unchanged for a denied operation.

**Exact gates** (tool pin unset; mandatory package scopes below):

```sh
python3 design_docs/verification/world-row140-m0/m0_matrix.py --protected-self-test
python3 design_docs/verification/world-row140-m0/m0_matrix.py --srt-node-modules "$GITHUB_WORKSPACE/.row156-ci/srt/node_modules" --node "$(command -v node)" --seccomp-read allowread --protected-only --scratch-root "$GITHUB_WORKSPACE/.row156-ci/matrix" --json "$GITHUB_WORKSPACE/.row156-ci/protected-linux.json"
go test -count=1 ./host/broker/... -run '^TestProtectedMatrixModeContract$' -timeout 120s -v
go test -race -count=1 ./host/broker/...
go vet ./...
```

**Baseline for every gate:**

- `python3 design_docs/verification/world-row140-m0/m0_matrix.py --protected-self-test` — NOT PRESENT at base; Linux measurement is CI ONLY, not runnable on this Darwin sandbox. No green claim.
- `python3 design_docs/verification/world-row140-m0/m0_matrix.py --srt-node-modules "$GITHUB_WORKSPACE/.row156-ci/srt/node_modules" --node "$(command -v node)" --seccomp-read allowread --protected-only --scratch-root "$GITHUB_WORKSPACE/.row156-ci/matrix" --json "$GITHUB_WORKSPACE/.row156-ci/protected-linux.json"` — NOT PRESENT at base; Linux measurement is CI ONLY, not runnable on this Darwin sandbox. No green claim.
- `go test -count=1 ./host/broker/... -run '^TestProtectedMatrixModeContract$' -timeout 120s -v` — NEW TEST/ROW absent at base; NOT RUN. Zero tests would not be a baseline pass.
- `go test -race -count=1 ./host/broker/...` — NOT MEASURED; green on this base is UNKNOWN. Unfiltered suites include legacy off-worktree writes and/or sockets; do not execute here or infer green from the focused baseline. Any socket result is UNINFORMATIVE UNDER SANDBOX.
- `go vet ./...` — GREEN: rc 0, 1.46 s, measured on this base with WORLD_TOOL_AILANG_BIN unset.

No `.ail` change; no AIL verification command. New fixture writes are worktree-local; no seeded external state. Any daemon/REST table deadline is created inside its own `t.Run`.

### M1a — Parse, fold, canonicalize and digest the protected grammar

Design M1; **242 LOC total / 30 min**, own PR, depends on M0b. Status **planned**.

**Files touched and base functions:**

- `host/broker/exec_protected.go NEW (ls: absent); add ParseProtectedPatterns, canonicalProtected, protectedDigest, foldProtectedKey, matchProtected; existing integration anchor host/broker/exec_profile.go:127 ParseExecProfile`
- `host/broker/exec_protected_test.go NEW (parent host/broker exists); reuse fixture conventions at exec_profile_test.go:36 profileJSON; measured upstream matcher at design_docs/verification/world-row156-design-2026-10-17/q0_anchors.out.txt:44`

Allowlist: `host/broker/exec_protected.go`, `host/broker/exec_protected_test.go`, `go.mod`, `go.sum`.

Dependency anchors: existing `go.mod:1` requirements and `go.sum:1` checksums. Add direct `golang.org/x/text v0.42.0` for `unicode/norm`; `go list -mod=readonly -m golang.org/x/text` confirms this version is already selected. Add only required checksums; no graph upgrade or handwritten NFC tables.

**Implementation and landable boundary.** Implement §4.1 with no profile wiring yet. Preserve optional empty behavior, NFC/full SimpleFold orbits matching fileguard, canonical fold-key ordering/dedup and domain-separated sha256 bytes. Do not mistake strings.ToLower for Unicode orbit folding. Golden expected digest literal is independently computed. Parameterize refusal rows so each guard has its own negative control. Do not add .ail or a third language matcher.

Refusal controls are table rows, each mutates only its named validation branch. **Iteration-282 correction to historical M1a guidance:** V16 refuted the ASCII-pattern fallback; do not drop Unicode rows / do not narrow the guarantee. Use the approved candidate-generator plus mirror-filter design in M1c/M1d.

**Named tests and assertions:**

- `TestProtectedGrammar`: Positive directory/basename forms and exact indexed refusal reasons.
- `TestProtectedGrammarRefusesDoubleStarPrefix`: **/x must refuse with basename fix.
- `TestProtectedDigestGolden`: Pin exact domain/newlines/canonical sequence and fixed hash; duplicates and order stable.
- `TestProtectedMatchFixture`: Measured V7 table includes pkg/A_TEST.GO, ſ/K orbits, NFC equivalent spellings, root/deep hits and non-hit controls.

**One named negative control per guard (independent mutation runs):**

- `TestProtectedGrammarRefusesDoubleStarPrefix` kills **MUT-DOUBLESTAR-ACCEPTED** (Accept **/x); assertion: unsupported **/x rejected.
- `TestProtectedMatchFixture/pkg-A_TEST.GO` kills **MUT-NO-FOLD** (Drop folding); assertion: case-folded path must match.
- `TestProtectedMatchFixture/nfc` kills **MUT-NO-NFC** (Drop NFC normalization); assertion: canonically equivalent path must match.
- `TestProtectedDigestGolden` kills **MUT-DIGEST-TEXT** (Remove canonical sorting/dedup/domain/newline rule); assertion: exact independent bytes and digest.
- `TestProtectedGrammar/count-65` kills **MUT-GRAMMAR-COUNT** (Remove entry count cap); assertion: 65 entries rejected.
- `TestProtectedGrammar/empty-and-257` kills **MUT-GRAMMAR-LENGTH** (Remove byte length checks); assertion: empty and overlong entries rejected.
- `TestProtectedGrammar/characters` kills **MUT-GRAMMAR-CHARACTERS** (Allow invalid Unicode/control or backslash); assertion: NUL/control/backslash/invalid UTF-8 refused.
- `TestProtectedGrammar/relative` kills **MUT-GRAMMAR-RELATIVE** (Allow absolute/dot/dotdot/nonclean directory); assertion: unclean directory refused.
- `TestProtectedGrammar/slash-glob` kills **MUT-GRAMMAR-SLASH** (Allow path slash-globs); assertion: src/*.ts refused.
- `TestProtectedGrammar/bracket` kills **MUT-GRAMMAR-GLOB** (Allow malformed bracket glob); assertion: bad path.Match rejected.
- `TestProtectedGrammar/everything` kills **MUT-GRAMMAR-ALL** (Allow blanket *, **, or ./**); assertion: whole-tree patterns rejected.
- `TestProtectedGrammar/unsupported-meta` kills **MUT-GRAMMAR-META** (Allow braces or leading !); assertion: unsupported grammar rejected.
- `TestProtectedGrammar/already-denied` kills **MUT-GRAMMAR-DENY** (Allow .git or fixed deny directory descendants); assertion: already denied directory form rejected.

**Exact gates** (tool pin unset; mandatory package scopes below):

```sh
go test -count=1 ./host/broker/... -run '^TestProtectedGrammar$' -timeout 120s -v
go test -count=1 ./host/broker/... -run '^TestProtectedGrammarRefusesDoubleStarPrefix$' -timeout 120s -v
go test -count=1 ./host/broker/... -run '^TestProtectedDigestGolden$' -timeout 120s -v
go test -count=1 ./host/broker/... -run '^TestProtectedMatchFixture$' -timeout 120s -v
go test -race -count=1 ./host/broker/...
go vet ./...
```

**Baseline for every gate:**

- `go test -count=1 ./host/broker/... -run '^TestProtectedGrammar$' -timeout 120s -v` — NEW TEST/ROW absent at base; NOT RUN. Zero tests would not be a baseline pass.
- `go test -count=1 ./host/broker/... -run '^TestProtectedGrammarRefusesDoubleStarPrefix$' -timeout 120s -v` — NEW TEST/ROW absent at base; NOT RUN. Zero tests would not be a baseline pass.
- `go test -count=1 ./host/broker/... -run '^TestProtectedDigestGolden$' -timeout 120s -v` — NEW TEST/ROW absent at base; NOT RUN. Zero tests would not be a baseline pass.
- `go test -count=1 ./host/broker/... -run '^TestProtectedMatchFixture$' -timeout 120s -v` — NEW TEST/ROW absent at base; NOT RUN. Zero tests would not be a baseline pass.
- `go test -race -count=1 ./host/broker/...` — NOT MEASURED; green on this base is UNKNOWN. Unfiltered suites include legacy off-worktree writes and/or sockets; do not execute here or infer green from the focused baseline. Any socket result is UNINFORMATIVE UNDER SANDBOX.
- `go vet ./...` — GREEN: rc 0, 1.46 s, measured on this base with WORLD_TOOL_AILANG_BIN unset.

No `.ail` change; no AIL verification command. New fixture writes are worktree-local; no seeded external state. Any daemon/REST table deadline is created inside its own `t.Run`.

### M1b — Render the four consumers and rebase module roots

Design M1; **230 LOC total / 25 min**, own PR, depends on M1a. Status **planned**.

**Files touched and base functions:**

- `host/broker/exec_protected.go NEW at base; add rebaseForSandbox, renderSrtProtected, gitPathspecs; integration anchors host/broker/exec_settings.go:91 renderExecSettings and host/daemon/workspace_layout.go:85 sandboxRoot`
- `host/broker/exec_protected_test.go NEW; fixture anchor host/broker/handlers_ailang_run_test.go:32 policyLines`

Allowlist: `host/broker/exec_protected.go`, `host/broker/exec_protected_test.go`.

**Implementation and landable boundary.** Add pure renderers. Emit directory git pathspecs as BOTH :(top,glob,icase)D and :(top,glob,icase)D/**, basename as :(top,glob,icase)**/G. Darwin uses literal dirs plus recursive basename globs; Linux only literal dirs. Rebase at module sandbox, rendering ancestor/equal directories as *, omitting disjoint dirs and leaving basenames. Return fresh slices; preserve empty nil semantics.

TestProtectedMatchFixture is the always-on companion for the optional binary arm. Future production files are still absent from the base; do not invent base line numbers for them.

**Named tests and assertions:**

- `TestProtectedRebase`: Root ., child/ancestor/equal/disjoint directories, basename and case rows.
- `TestProtectedRenderers`: Exact AILANG/Darwin/Linux/git string table, directory itself included.
- `TestProtectedMatchFixturePinned`: Real v0.52.1 fileguard policy-tool equivalence with same V7 matrix; skip if WORLD_TOOL_AILANG_BIN unset; no home writes.

**One named negative control per guard (independent mutation runs):**

- `TestProtectedRebase/under-sub` kills **MUT-NO-REBASE** (Remove rebasing); assertion: sub/tests/** becomes tests/**.
- `TestProtectedRebase/ancestor` kills **MUT-REBASE-ANCESTOR** (Drop ancestor entire-sandbox protection); assertion: ancestor renders *.
- `TestProtectedRenderers/linux` kills **MUT-LINUX-GLOB** (Render Linux basename glob); assertion: no Linux basename glob.
- `TestProtectedRenderers/darwin` kills **MUT-DARWIN-DEPTH** (Render Darwin basename only one level); assertion: recursive **/G exact string.

**Exact gates** (tool pin unset; mandatory package scopes below):

```sh
go test -count=1 ./host/broker/... -run '^TestProtectedRebase$' -timeout 120s -v
go test -count=1 ./host/broker/... -run '^TestProtectedRenderers$' -timeout 120s -v
go test -count=1 ./host/broker/... -run '^TestProtectedMatchFixturePinned$' -timeout 120s -v
go test -race -count=1 ./host/broker/...
go vet ./...
```

**Baseline for every gate:**

- `go test -count=1 ./host/broker/... -run '^TestProtectedRebase$' -timeout 120s -v` — NEW TEST/ROW absent at base; NOT RUN. Zero tests would not be a baseline pass.
- `go test -count=1 ./host/broker/... -run '^TestProtectedRenderers$' -timeout 120s -v` — NEW TEST/ROW absent at base; NOT RUN. Zero tests would not be a baseline pass.
- `go test -count=1 ./host/broker/... -run '^TestProtectedMatchFixturePinned$' -timeout 120s -v` — NEW TEST/ROW absent at base; NOT RUN. Zero tests would not be a baseline pass.
- `go test -race -count=1 ./host/broker/...` — NOT MEASURED; green on this base is UNKNOWN. Unfiltered suites include legacy off-worktree writes and/or sockets; do not execute here or infer green from the focused baseline. Any socket result is UNINFORMATIVE UNDER SANDBOX.
- `go vet ./...` — GREEN: rc 0, 1.46 s, measured on this base with WORLD_TOOL_AILANG_BIN unset.

No `.ail` change; no AIL verification command. New fixture writes are worktree-local; no seeded external state. Any daemon/REST table deadline is created inside its own `t.Run`.

### M1c — Pure candidate renderer and mirror filter

Design M1; **250 LOC total / 30 min**, own PR/evidence boundary, depends on M1b. Status **planned**.

**Files + allowlist:** `host/broker/exec_protected.go`, `host/broker/exec_protected_test.go`. No other files. Future files/tests are absent at base unless explicitly stated.

**Implementation and assertions.** [r4/r4.1/r4.2] Replace gitPathspecs with gitCandidatePathspecs in host/broker/exec_protected.go:317; keep top,glob,icase and bare W(D) plus W(D)/**. W replaces s S k K ; `, every non-ASCII rune, every ? and each whole bracket class with *, collapses adjacent stars within a segment, and preserves / and the structural /**. protectedWidenRunes is the closed widen set, checked against a full scalar scan. Git is only a candidate generator; matchProtected remains authoritative. Re-derive M1b TestProtectedRenderers/git's fixed golden for *.test.ts, conftest.py, src/gen/**, tests/** as **/*.te*t.t*, **/confte*t.py, *rc/gen plus *rc/gen/**, te*t* plus te*t*/**, each with :(top,glob,icase); call the replacement renderer. No production I/O. parseProtectedCandidates receives canonical patterns, statusOut/ignoredOut/flagsOut []byte and the fold-keyed episode-start ignored set, returns head, changes, error. Parse all three NUL-terminated streams, preserve raw path bytes for matchProtected's NFC fold, keep branch.oid unfiltered; map ordinary changed/deleted/type-change and ? records, ignored files, and lower-case/S index flags. Validate complete record fields/types/tags before filtering; malformed or >8 MiB in ANY input returns an error, never a clean empty success. Parse v2 2-record extra NUL path correctly even with --no-renames; preserve whitespace/newlines inside NUL-delimited paths. Snapshot subtraction applies ONLY to ignored entries; deduplicate exact raw path/status pairs without hiding case-distinct paths or distinct statuses; only ignored snapshot membership uses fold keys.  Signature contract: parseProtectedCandidates(canonical []string, statusOut, ignoredOut, flagsOut []byte, ignoredBaseline map[string]struct{}) (head string, changes []protectedChange, err error); the exact set representation may follow local style, but neither canonical patterns nor the ignored baseline may be omitted. Introduce the minimal protectedChange data type consumed by the parser here; M6a reuses it rather than declaring a second type. Budget <=140 changed production LOC + <=110 changed test LOC, <=250 total and <=30 minutes; if this cannot fit, split again before coding, never omit rows.

Do not drop Unicode rows / do not narrow the guarantee. V16 refuted the old premise; the r4.2 narrow-refinement disposition now authorizes this replacement. Both files exist at 451dc3d; candidate functions/tests are absent. Existing TestProtectedRenderers/git passed only its OLD golden in the focused baseline, not the new candidate golden.

**Named tests and assertions:**

- `TestProtectedFoldToASCIIClosed` (AC1c.1): /closed-set scans U+0080..U+10FFFF excluding surrogates: exactly U+017F, U+037E, U+1FEF, U+212A have ASCII in foldProtectedKey; both ASCII cases of every key must be in protectedWidenRunes. /no-multi-rune asserts keys of ß, ﬀ, ﬁ stay themselves (NFC, not NFKC).
- `TestProtectedCandidatePathspecsGolden` (AC1c.2): Independent fixed strings: /directory tests/** -> te*t* and te*t*/**; /kelvin *k_test.go -> **/*_te*t.go; /non-ascii é_test.go and /qmark ?_test.go -> **/*_te*t.go; /class [ab]_x.go -> **/*_x.go. Every string carries :(top,glob,icase); /collapse-stars checks collapse and slash preservation.
- `TestProtectedCandidateParseFailsClosed` (AC1c.4): Synthetic bytes only: /over-cap-status, /over-cap-ignored, /over-cap-flags each exceeds 8 MiB and must error; /missing-nul, /unknown-status, /bad-flag and /bad-fields must error even for mirror-false paths; /valid-three-streams asserts head and exact filtered path/status set, including type-change, ignored and index flags.
- `TestProtectedIgnoredBaseline` (AC1c.5): Synthetic inputs: /standing suppresses snapshot entry; /new reports new path as ignored; /case-variant and /nfc-variant suppress equivalent fold keys; /tracked-not-exempt ensures snapshot cannot mask a tracked modification.
- `TestProtectedRenderers` (M1b golden re-derived): Existing /git calls gitCandidatePathspecs and compares the widened fixed golden; other existing consumers retain their prior expectations.

**Mutation → exact killer (one compiling mutant at a time):**

- **MUT-WIDEN-SET-STALE** → `TestProtectedFoldToASCIIClosed/closed-set`: Remove ; from protectedWidenRunes; assertion: Scan still finds U+037E -> ; and rejects missing coverage.
- **MUT-CAP-SILENT** → `TestProtectedCandidateParseFailsClosed/over-cap-ignored`: Truncate excess ignored output to cap and return success; assertion: An oversized ignored stream returns non-nil error.
- **MUT-NO-IGNORED-BASELINE** → `TestProtectedIgnoredBaseline/standing`: Ignore the supplied episode-start snapshot; assertion: Standing ignored protected path produces zero changes, with no error.
- **MUT-PARSE-MALFORMED-OPEN** → `TestProtectedCandidateParseFailsClosed/missing-nul`: Accept unterminated status record; assertion: Malformed stream returns non-nil error.
- **MUT-PARSE-STATUS-OPEN** → `TestProtectedCandidateParseFailsClosed/unknown-status`: Ignore unknown porcelain record type; assertion: Unknown type returns non-nil error.
- **MUT-PARSE-FLAG-OPEN** → `TestProtectedCandidateParseFailsClosed/bad-flag`: Accept invalid ls-files tag; assertion: Invalid tag returns non-nil error.

- **MUT-PARSE-FIELDS-OPEN** → `TestProtectedCandidateParseFailsClosed/bad-fields`: Accept malformed porcelain field count; assertion: Bad required fields return non-nil error even on a mirror-false path.

Each mutation uses its exact `go test -count=1 ./host/broker/ -run '^<full killer above>$' -timeout 300s -v` selector (literal commands in JSON); these mutation gates are **NOT MEASURED: named additions absent at base; zero tests is not a pass**. One valid guard/data mutant at a time, exact named killer/subtest and assertion. Save and restore allowed bytes in worktree-local scratch without git. Require named red, then restoration green. No compile error, panic, timeout, skip, missing executable, socket denial, unrelated refusal or grep-only assertion is a kill. Fixed independent expectation must be reachable after the changed guard; delete/move masking cannot satisfy r3 bare-D control.

**Exact gates** (WORLD_TOOL_AILANG_BIN unset; TMPDIR/GOTMPDIR worktree-local):

```sh
go test -count=1 ./host/broker/ -run '^TestProtectedFoldToASCIIClosed$' -timeout 300s -v
go test -count=1 ./host/broker/ -run '^TestProtectedCandidatePathspecsGolden$' -timeout 300s -v
go test -count=1 ./host/broker/ -run '^TestProtectedCandidateParseFailsClosed$' -timeout 300s -v
go test -count=1 ./host/broker/ -run '^TestProtectedIgnoredBaseline$' -timeout 300s -v
go test -count=1 ./host/broker/ -run '^TestProtectedRenderers$' -timeout 300s -v
go test -race -count=1 ./host/broker/...
go vet ./host/broker/
```

**Baseline for every gate at 451dc3d:**

- `go test -count=1 ./host/broker/ -run '^TestProtectedFoldToASCIIClosed$' -timeout 300s -v` — NOT MEASURED: absent at base; zero tests is not a pass. No fixture Git repositories were created by planner.
- `go test -count=1 ./host/broker/ -run '^TestProtectedCandidatePathspecsGolden$' -timeout 300s -v` — NOT MEASURED: absent at base; zero tests is not a pass. No fixture Git repositories were created by planner.
- `go test -count=1 ./host/broker/ -run '^TestProtectedCandidateParseFailsClosed$' -timeout 300s -v` — NOT MEASURED: absent at base; zero tests is not a pass. No fixture Git repositories were created by planner.
- `go test -count=1 ./host/broker/ -run '^TestProtectedIgnoredBaseline$' -timeout 300s -v` — NOT MEASURED: absent at base; zero tests is not a pass. No fixture Git repositories were created by planner.
- `go test -count=1 ./host/broker/ -run '^TestProtectedRenderers$' -timeout 300s -v` — Exact standalone gate NOT MEASURED. Existing old golden passed within focused TestProtected baseline (rc 0); widened assertions absent at base; zero tests is not a pass.
- `go test -race -count=1 ./host/broker/...` — NOT MEASURED (PF6: legacy tests write real HOME); requires controller-contained runner. Socket-touching results are UNINFORMATIVE UNDER SANDBOX.
- `go vet ./host/broker/` — GREEN rc 0 at 451dc3d, iteration 282; WORLD_TOOL_AILANG_BIN unset, TMPDIR/GOTMPDIR worktree-local. Package-only vet measured; no repository-wide vet claim.

No `.ail` changes; no AIL verification. Full package acceptance requires contained controller runner; sockets are UNINFORMATIVE UNDER SANDBOX.

### M1d — Real-Git candidate and mirror equivalence

Design M1; **240 LOC total / 30 min**, own PR/evidence boundary, depends on M1c. Status **planned**.

**Files + allowlist:** `host/broker/exec_protected_git_test.go` (persistent tests), `host/broker/exec_protected.go` (temporary mutation edits only, restored byte-for-byte). No other files. Future files/tests are absent at base unless explicitly stated.

**Implementation and assertions.** Implement AC1c.3 only in NEW host/broker/exec_protected_git_test.go (parent host/broker exists); read-only reuse of protectedFixtureRows and matcher helpers. <=240 changed test LOC, no production edits; M1c+M1d totals <=140 production + <=350 tests. One disposable repo per spelling/row under t.TempDir(), with fake HOME inside fixture, GIT_CONFIG_NOSYSTEM=1, GIT_CONFIG_GLOBAL=/dev/null, cleared inherited Git overrides, deterministic local user config, no hooks/filters, GIT_OPTIONAL_LOCKS=0, core.fsmonitor=false. Git absence: host/broker test search found no CI absence-detection convention (CI=1 occurs only in child environments); use os.Getenv("CI") != "" to FAIL when absent on CI, SKIP only locally. Never convert a Git setup/command error into a skip. Every protectedFixtureRows row is retained unchanged: pkg-A_TEST.GO, root, deep, dir-root, dir-deep, long-s, kelvin, non-hit-sibling, non-hit-nested-dir, non-hit-base. Add greek-question-mark path pkg/a;_test.go (U+037E; pattern a;_test.go), greek-varia (U+1FEF vs ASCII backtick pattern), NFD pkg/e+U+0301_test.go indexed with core.precomposeunicode=false (assert index raw bytes; use cacheinfo if needed), qmark-utf8 pkg/é_test.go vs ?_test.go, class-utf8 pkg/é_test.go vs [^a]_test.go, excluded zzign/a_test.go beside excluded nm/c.txt, directory-replaced-by-file and directory-replaced-by-symlink, deleted tests/x.py, assume-unchanged and skip-worktree. Add mirror-false controls sharp-s (pkg/a_teßt.go vs *_test.go, assert candidate is actually selected before filtering) and ligature-ff (buﬀ/x.py vs buff/**); all mirror-false rows absent after filter. For EVERY row run root-CWD and subdirectory-CWD subtests, executing the three design §4.5 invocations with gitCandidatePathspecs then parseProtectedCandidates and snapshot. Require exact expected path/status set, no missing or extra entries, expected mirror boolean independently checked against fixture bank; branch.oid equals known commit. Root/subdir CWD must exercise actual selector anchoring, not force -C root for both; normalize any CWD-relative ls-files output to root-relative test input without changing selection. Production always uses cmd.Dir=h.worktree. The subdirectory-CWD control includes a root-level matching path outside sub/. Replacements include an initially untracked bare tests path, so deleting tracked descendants cannot mask losing W(D). Ignored nm/c.txt row asserts successful empty filtered output and command argv --ignored=no plus a per-file ignored stream; replacing the ignored leg with status --ignored=matching must fail that exact assertion, even if parser would otherwise discard ! entries.

Do not drop Unicode rows / do not narrow the guarantee. No old TestProtectedGitPathspecEquivalence acceptance remains. Git fixture operations are executor-only and confined to t.TempDir(); planner creates none. Mutation-only allowlist: host/broker/exec_protected.go may be temporarily changed for renderer/filter mutants, then restored byte-for-byte without Git; its persistent diff must be empty. Only exec_protected_git_test.go gains persistent test code.

For each real-Git row capture the ignored baseline AFTER initial fixture setup but BEFORE the row mutation; ignored-creation arms create zzign/a_test.go only after that snapshot, while nm/c.txt remains a mirror-false control. Assert nonempty raw candidate selection where the row is intended to exercise filtering; candidate absence cannot prove a filter kill.

**Named tests and assertions:**

- `TestProtectedGitCandidateEquivalence` (AC1c.3): Full unchanged fixture bank plus all r4/r4.1/r4.2 rows listed below; root and subdirectory CWD for every spelling; candidate+filter output equals independent mirror-true changed path/status set, and every mirror-false row is absent.

**Mutation → exact killer (one compiling mutant at a time):**

- **MUT-ASCII-ONLY-SUPERSET** → `TestProtectedGitCandidateEquivalence/long-s/root-CWD`: Restore r3 unwidened pathspecs; assertion: teſts/x.py is selected and reported as modified.
- **MUT-ASCII-ONLY-SUPERSET** → `TestProtectedGitCandidateEquivalence/kelvin/root-CWD`: Restore r3 unwidened pathspecs; assertion: pkg/K_test.go is selected and reported.
- **MUT-ASCII-ONLY-SUPERSET** → `TestProtectedGitCandidateEquivalence/greek-question-mark/root-CWD`: Restore r3 unwidened pathspecs; assertion: U+037E spelling is reported.
- **MUT-QMARK-KEPT** → `TestProtectedGitCandidateEquivalence/qmark-utf8/root-CWD`: Leave ? and bracket classes unwidened; assertion: pkg/é_test.go is reported under ?_test.go.
- **MUT-QMARK-KEPT** → `TestProtectedGitCandidateEquivalence/class-utf8/root-CWD`: Leave bracket classes unwidened; assertion: Multi-byte é matches mirror [^a] and is reported.
- **MUT-NO-MIRROR-FILTER** → `TestProtectedGitCandidateEquivalence/non-hit-sibling/root-CWD`: Report every Git candidate; assertion: tests2/x.py is a candidate but absent from reported changes.
- **MUT-IGNORED-MATCHING** → `TestProtectedGitCandidateEquivalence/ignored-nm-c-txt/root-CWD`: Use status --ignored=matching as ignored-file source; assertion: Excluded nm/c.txt yields successful empty result; no ! directory record is accepted and argv remains per-file ignored enumeration.
- **MUT-NO-TOP** → `TestProtectedGitCandidateEquivalence/root/subdirectory-CWD`: Remove top from all pathspecs; assertion: Root c_test.go outside sub/ is still selected and reported.
- **MUT-NO-BARE-D** → `TestProtectedGitCandidateEquivalence/directory-replaced-by-file/subdirectory-CWD`: Drop bare W(D), keep W(D)/**; assertion: Untracked bare tests file is reported, independently of descendant deletions.
- **MUT-NO-BARE-D** → `TestProtectedGitCandidateEquivalence/directory-replaced-by-symlink/subdirectory-CWD`: Drop bare W(D), keep W(D)/**; assertion: Bare tests symlink is reported without descendant masking.
- **MUT-PATHSPEC-NO-ICASE** → `TestProtectedGitCandidateEquivalence/pkg-A_TEST.GO/root-CWD`: Drop icase; assertion: Case-varied basename is selected on case-sensitive Git index.

Each mutation uses its exact `go test -count=1 ./host/broker/ -run '^<full killer above>$' -timeout 300s -v` selector (literal commands in JSON); these mutation gates are **NOT MEASURED: named additions absent at base; zero tests is not a pass**. One valid guard/data mutant at a time, exact named killer/subtest and assertion. Save and restore allowed bytes in worktree-local scratch without git. Require named red, then restoration green. No compile error, panic, timeout, skip, missing executable, socket denial, unrelated refusal or grep-only assertion is a kill. Fixed independent expectation must be reachable after the changed guard; delete/move masking cannot satisfy r3 bare-D control.

**Exact gates** (WORLD_TOOL_AILANG_BIN unset; TMPDIR/GOTMPDIR worktree-local):

```sh
go test -count=1 ./host/broker/ -run '^TestProtectedGitCandidateEquivalence$' -timeout 300s -v
go test -race -count=1 ./host/broker/...
go vet ./host/broker/
```

**Baseline for every gate at 451dc3d:**

- `go test -count=1 ./host/broker/ -run '^TestProtectedGitCandidateEquivalence$' -timeout 300s -v` — NOT MEASURED: absent at base; zero tests is not a pass. No fixture Git repositories were created by planner.
- `go test -race -count=1 ./host/broker/...` — NOT MEASURED (PF6: legacy tests write real HOME); requires controller-contained runner. Socket-touching results are UNINFORMATIVE UNDER SANDBOX.
- `go vet ./host/broker/` — GREEN rc 0 at 451dc3d, iteration 282; WORLD_TOOL_AILANG_BIN unset, TMPDIR/GOTMPDIR worktree-local. Package-only vet measured; no repository-wide vet claim.

No `.ail` changes; no AIL verification. Full package acceptance requires contained controller runner; sockets are UNINFORMATIVE UNDER SANDBOX.

### M2a — Load the additive protected profile key

Design M2; **200 LOC total / 25 min**, own PR, depends on M1d. Status **planned**.

**Files touched and base functions:**

- `host/broker/exec_profile.go:56 ExecProfile; :127 ParseExecProfile; :128 strictObject allowlist`
- `host/broker/exec_profile_test.go:36 profileJSON; :64 TestParseExecProfileDesignExample; :90 TestParseExecProfileRefusals`

Allowlist: `host/broker/exec_profile.go`, `host/broker/exec_profile_test.go`.

**Implementation and landable boundary.** Add Protected canonical list and allow protected in strictObject. Decode optional array strictly and route through ParseProtectedPatterns. Preserve ExecProfileVersion and profile file-byte digest; distinguish absent/[] compatibility from malformed values, and name protected[i] in errors. No host behavior change until consumers are wired.

**Named tests and assertions:**

- `TestExecProfileProtected`: Optional/empty/list canonicalization and existing version/file digest behavior.
- `TestExecProfileProtectedUnknownFormRefused`: src/*.ts fails load with indexed reason.

**One named negative control per guard (independent mutation runs):**

- `TestExecProfileProtectedUnknownFormRefused` kills **MUT-PROFILE-NO-CHECK** (Remove protected grammar validation at profile load); assertion: protected[0] names refused src/*.ts.
- `TestExecProfileProtected/wrong-type` kills **MUT-PROFILE-TYPE** (Accept a string/null/object instead of declared pattern list); assertion: wrong field type refused.
- `TestExecProfileProtected/compatibility` kills **MUT-PROFILE-COMPAT** (Change schema version or hash canonical JSON instead of input bytes); assertion: v1 and raw-byte digest unchanged for absent field.

**Exact gates** (tool pin unset; mandatory package scopes below):

```sh
go test -count=1 ./host/broker/... -run '^TestExecProfileProtected$' -timeout 120s -v
go test -count=1 ./host/broker/... -run '^TestExecProfileProtectedUnknownFormRefused$' -timeout 120s -v
go test -race -count=1 ./host/broker/...
go vet ./...
```

**Baseline for every gate:**

- `go test -count=1 ./host/broker/... -run '^TestExecProfileProtected$' -timeout 120s -v` — NEW TEST/ROW absent at base; NOT RUN. Zero tests would not be a baseline pass.
- `go test -count=1 ./host/broker/... -run '^TestExecProfileProtectedUnknownFormRefused$' -timeout 120s -v` — NEW TEST/ROW absent at base; NOT RUN. Zero tests would not be a baseline pass.
- `go test -race -count=1 ./host/broker/...` — NOT MEASURED; green on this base is UNKNOWN. Unfiltered suites include legacy off-worktree writes and/or sockets; do not execute here or infer green from the focused baseline. Any socket result is UNINFORMATIVE UNDER SANDBOX.
- `go vet ./...` — GREEN: rc 0, 1.46 s, measured on this base with WORLD_TOOL_AILANG_BIN unset.

No `.ail` change; no AIL verification command. New fixture writes are worktree-local; no seeded external state. Any daemon/REST table deadline is created inside its own `t.Run`.

### M2b — Add startup grammar refusal rows with matching usage table

Design M2; **180 LOC total / 25 min**, own PR, depends on M2a. Status **planned**.

**Files touched and base functions:**

- `host/daemon/exec_startup_test.go:421 TestExecStartupRefusalTable; :121 newExecRig; :159 execRig.start`
- `docs/QUICKSTART.md:640 section 10; startup-refusal table within that section`

Allowlist: `host/daemon/exec_startup_test.go`, `docs/QUICKSTART.md`.

**Implementation and landable boundary.** Extend the existing exact-census startup refusal table and QUICKSTART rows together. Each grammar failure stops at StageConfig before any listener/probe/subprocess. Use fake tool/node fixtures, no pinned tool. Each t.Run creates its own context.WithTimeout and cleanup; never share a deadline across cases. This early doc edit is only refusal rows; worked protected examples arrive in M7b.

The existing test exists at base, but the protected subrows are new. Existing full daemon gate includes sockets: any local socket verdict is UNINFORMATIVE UNDER SANDBOX. No report of PASS/FAIL for such a run.

**Named tests and assertions:**

- `TestExecStartupRefusalTable`: protected-double-star, protected-slash-glob, protected-over-limit and malformed-value rows name their field and refuse before serve; existing row census stays exact.

**One named negative control per guard (independent mutation runs):**

- `TestExecStartupRefusalTable/protected-slash-glob` kills **MUT-STARTUP-PROTECTED-IGNORED** (Suppress grammar error in startup load); assertion: StageConfig error names protected[0] and no listener starts.

**Exact gates** (tool pin unset; mandatory package scopes below):

```sh
go test -count=1 ./host/daemon/... -run '^TestExecStartupRefusalTable$' -timeout 120s -v
go test -race -count=1 ./host/daemon/...
go vet ./...
```

**Baseline for every gate:**

- `go test -count=1 ./host/daemon/... -run '^TestExecStartupRefusalTable$' -timeout 120s -v` — Existing test; new protected subrows absent. Exact command NOT RUN; daemon/socket activity is UNINFORMATIVE UNDER SANDBOX. No green claim.
- `go test -race -count=1 ./host/daemon/...` — NOT MEASURED; green on this base is UNKNOWN. Unfiltered suites include legacy off-worktree writes and/or sockets; do not execute here or infer green from the focused baseline. Any socket result is UNINFORMATIVE UNDER SANDBOX.
- `go vet ./...` — GREEN: rc 0, 1.46 s, measured on this base with WORLD_TOOL_AILANG_BIN unset.

No `.ail` change; no AIL verification command. New fixture writes are worktree-local; no seeded external state. Any daemon/REST table deadline is created inside its own `t.Run`.

### M3a — Append protected sandbox settings without changing old bytes

Design M3; **210 LOC total / 25 min**, own PR, depends on M2b. Status **planned**.

**Files touched and base functions:**

- `host/broker/exec_settings.go:25 execSettingsSpec; :91 renderExecSettings`
- `host/broker/exec_settings_test.go:129 TestRenderExecSettingsGolden`
- `host/broker/handlers_exec.go:123 NewExecHandler; :156 execSettingsSpec literal`

Allowlist: `host/broker/exec_settings.go`, `host/broker/exec_settings_test.go`, `host/broker/handlers_exec.go`.

**Implementation and landable boundary.** Add Protected to spec, pass profile.Protected from NewExecHandler and append OS-specific entries after existing deny names and srt defaults. Keep old unprotected golden files unchanged; add inline protected goldens so no new golden files are needed. Settings digest naturally tracks changed bytes. M0 CI evidence must agree before this production PR.

**Named tests and assertions:**

- `TestExecSettingsProtectedRendered`: Exact order and OS strings, including .git/hooks and .git/config before .git and before protection.
- `TestExecSettingsGoldenDarwin`: Protected Darwin bytes pinned independently.
- `TestExecSettingsGoldenLinux`: Protected Linux directory-only bytes pinned independently.
- `TestExecSettingsNoProtectedIsByteIdentical`: Compare absent/empty to both existing golden files; settings digest identical.

**One named negative control per guard (independent mutation runs):**

- `TestExecSettingsProtectedRendered` kills **MUT-NO-SRT-PROTECTED** (Delete append); assertion: expected protected denyWrite suffix present.
- `TestExecSettingsGoldenLinux` kills **MUT-PROTECTED-BEFORE-GITHOOKS** (Insert protected entries before .git/hooks); assertion: existing prefix order remains exact.
- `TestExecSettingsGoldenLinux` kills **MUT-LINUX-GLOB** (Render Linux basename glob); assertion: Linux golden omits all basename globs.
- `TestExecSettingsNoProtectedIsByteIdentical` kills **MUT-SRT-EMPTY-BYTES** (Emit a changed empty settings rendering); assertion: unprotected bytes equal old golden.

**Exact gates** (tool pin unset; mandatory package scopes below):

```sh
go test -count=1 ./host/broker/... -run '^TestExecSettingsProtectedRendered$' -timeout 120s -v
go test -count=1 ./host/broker/... -run '^TestExecSettingsGoldenDarwin$' -timeout 120s -v
go test -count=1 ./host/broker/... -run '^TestExecSettingsGoldenLinux$' -timeout 120s -v
go test -count=1 ./host/broker/... -run '^TestExecSettingsNoProtectedIsByteIdentical$' -timeout 120s -v
go test -race -count=1 ./host/broker/...
go vet ./...
```

**Baseline for every gate:**

- `go test -count=1 ./host/broker/... -run '^TestExecSettingsProtectedRendered$' -timeout 120s -v` — NEW TEST/ROW absent at base; NOT RUN. Zero tests would not be a baseline pass.
- `go test -count=1 ./host/broker/... -run '^TestExecSettingsGoldenDarwin$' -timeout 120s -v` — NEW TEST/ROW absent at base; NOT RUN. Zero tests would not be a baseline pass.
- `go test -count=1 ./host/broker/... -run '^TestExecSettingsGoldenLinux$' -timeout 120s -v` — NEW TEST/ROW absent at base; NOT RUN. Zero tests would not be a baseline pass.
- `go test -count=1 ./host/broker/... -run '^TestExecSettingsNoProtectedIsByteIdentical$' -timeout 120s -v` — NEW TEST/ROW absent at base; NOT RUN. Zero tests would not be a baseline pass.
- `go test -race -count=1 ./host/broker/...` — NOT MEASURED; green on this base is UNKNOWN. Unfiltered suites include legacy off-worktree writes and/or sockets; do not execute here or infer green from the focused baseline. Any socket result is UNINFORMATIVE UNDER SANDBOX.
- `go vet ./...` — GREEN: rc 0, 1.46 s, measured on this base with WORLD_TOOL_AILANG_BIN unset.

No `.ail` change; no AIL verification command. New fixture writes are worktree-local; no seeded external state. Any daemon/REST table deadline is created inside its own `t.Run`.

### M3b — Exercise real srt protected writes with pin-free companions

Design M3; **235 LOC total / 30 min**, own PR, depends on M3a. Status **planned**.

**Files touched and base functions:**

- `host/broker/handlers_exec_srt_test.go:39 realSrt; :75 newSrtFixture; :264 TestExecSrtConfinementMatrix (do NOT reuse its external seeded decoy)`
- `host/broker/handlers_exec_m2_test.go:73 defaultExecProfile; :162 execDecode`

Allowlist: `host/broker/handlers_exec_srt_test.go`, `host/broker/handlers_exec_m2_test.go`.

**Implementation and landable boundary.** Add TestExecSrtProtectedWrites using a synthetic operator home and fixture paths only, skip before setup when srt/node/bwrap absent. Run V5 modify/create/root/deep/delete/rename/alias/case/parent-dir operations through actual handler. Assert explicit refusal token plus host snapshot for prevention rows; record allowed Linux basename/darwin-parent-rename rows as expected detection residuals, not false prevention. Pin-free fake-srt handler test checks settings consumption and arguments; it does not claim kernel confinement.

Absent real engine: optional test SKIPs, always-on TestExecSettingsProtectedRendered and TestExecProtectedSettingsThroughHandler PASS. Configured corrupt srt pin fails; only dependency absence skips. Do not use realHome or call the existing HOME-seeding confinement matrix as a new acceptance drill.

**Named tests and assertions:**

- `TestExecSrtProtectedWrites`: Real srt V5 rows, with explicit per-OS expectation and unsandboxed controls solely within fixture.
- `TestExecProtectedSettingsThroughHandler`: Always-on fake srt records the exact protected settings path/bytes passed by handler.

**One named negative control per guard (independent mutation runs):**

- `TestExecSrtProtectedWrites/modify` kills **MUT-NO-SRT-PROTECTED** (Delete protected deny append); assertion: modify protected host bytes must be refused on promised OS/form.
- `TestExecProtectedSettingsThroughHandler` kills **MUT-HANDLER-PROTECTED-DROPPED** (Stop forwarding profile.Protected into settings spec); assertion: fake launch observes protected settings suffix.

**Exact gates** (tool pin unset; mandatory package scopes below):

```sh
go test -count=1 ./host/broker/... -run '^TestExecSrtProtectedWrites$' -timeout 120s -v
go test -count=1 ./host/broker/... -run '^TestExecProtectedSettingsThroughHandler$' -timeout 120s -v
go test -race -count=1 ./host/broker/...
go vet ./...
```

**Baseline for every gate:**

- `go test -count=1 ./host/broker/... -run '^TestExecSrtProtectedWrites$' -timeout 120s -v` — NEW TEST/ROW absent at base; NOT RUN. Zero tests would not be a baseline pass.
- `go test -count=1 ./host/broker/... -run '^TestExecProtectedSettingsThroughHandler$' -timeout 120s -v` — NEW TEST/ROW absent at base; NOT RUN. Zero tests would not be a baseline pass.
- `go test -race -count=1 ./host/broker/...` — NOT MEASURED; green on this base is UNKNOWN. Unfiltered suites include legacy off-worktree writes and/or sockets; do not execute here or infer green from the focused baseline. Any socket result is UNINFORMATIVE UNDER SANDBOX.
- `go vet ./...` — GREEN: rc 0, 1.46 s, measured on this base with WORLD_TOOL_AILANG_BIN unset.

No `.ail` change; no AIL verification command. New fixture writes are worktree-local; no seeded external state. Any daemon/REST table deadline is created inside its own `t.Run`.

### M4a — Refuse configured git filters before any integrity status

Design M4; **170 LOC total / 20 min**, own PR, depends on M3b. Status **planned**.

**Files touched and base functions:**

- `host/broker/exec_protected_check.go` and `exec_protected_check_test.go`: NEW/absent at base; add `protectedGitPreflight`. Existing parent verified. Reuse `configuredGitFilters` at `host/broker/worktree_delta.go:297` and refusal behavior at `:326` (read only).
- `design_docs/planned/w-exec-protected-test-paths.md:204`: integrity prose; annotate the third no-filter preflight invocation.

Allowlist: `host/broker/exec_protected_check.go`, `host/broker/exec_protected_check_test.go`, `design_docs/planned/w-exec-protected-test-paths.md`.

**Implementation and landable boundary.** Verified base anchors: host/broker/exec_protected_check.go and exec_protected_check_test.go NEW/absent with host/broker parent verified; add protectedGitPreflight; reuse existing host/broker/worktree_delta.go:297 configuredGitFilters and :326 gitSnapshot no-filter refusal; design_docs/planned/w-exec-protected-test-paths.md:202 integrity-check prose. gitEnv disables system/global config, but repository filter.<name>.clean/process remains active. Add bounded filter preflight before each phase, refusing configured filters or preflight errors without running status. Reuse row184 query semantics. The later M4b checker includes this query and THREE candidate data commands in the SAME provisional 2500ms phase bound: one config query plus three data invocations. V24 includes the preflight overhead before M6a; preserve the no-host-filter-execution requirement. Annotate §4.5 with this measured coupling before later code lands. No result field yet.

**Named tests and assertions:**

- `TestProtectedGitFilterPreflight`: Temp repo local/included clean and process filter commands would create a marker inside fixture. Guard returns diagnostic and marker stays absent, no status spawn. No-filter control allows progression; query error refuses progression. Marker is created only by a killed mutant inside disposable worktree-local fixture, never as a seeded external guard.

**One named negative control per guard:**

- `TestProtectedGitFilterPreflight/configured-clean` kills **MUT-PROTECTED-FILTER-PREFLIGHT** (Remove configured-filter refusal before status.); assertion: host-side status never executes the configured clean filter; fixture marker must remain absent.
- `TestProtectedGitFilterPreflight/query-error` kills **MUT-PROTECTED-FILTER-ERROR-OPEN** (Treat failed filter configuration query as no filters.); assertion: status spawn count is zero and check diagnostic is present on query failure.

**Exact gates** (tool pin unset):

```sh
go test -count=1 ./host/broker/... -run '^TestProtectedGitFilterPreflight$' -timeout 120s -v
go test -race -count=1 ./host/broker/...
go vet ./host/broker/
```

**Baseline for every gate:**

- `go test -count=1 ./host/broker/... -run '^TestProtectedGitFilterPreflight$' -timeout 120s -v` — NEW TEST absent at base; NOT RUN; no green claim.
- `go test -race -count=1 ./host/broker/...` — NOT MEASURED; legacy external writes require contained runner. Any socket result is UNINFORMATIVE UNDER SANDBOX.
- `go vet ./host/broker/` — GREEN rc 0 at 451dc3d (iteration 282); package-only vet, tool pin unset, scratch contained.

No `.ail` change. Every fixture/sentinel is under the worktree; no real HOME/state seeding. Full package gate requires the contained runner; socket results are UNINFORMATIVE UNDER SANDBOX.

### M4b — Build the bounded integrity checker needed by probe arm 8

Design M4; **245 LOC total / 30 min**, own PR/evidence boundary, depends on M4a. Status **planned**.

**Files + allowlist:** `host/broker/exec_protected_check.go`, `host/broker/exec_protected_check_test.go`. No other files. Future files/tests are absent at base unless explicitly stated.

**Implementation and assertions.** [r4.x] Add checkProtected in NEW exec_protected_check.go; reuse gitEnv/runBounded and M4a configured-filter preflight without editing worktree_delta.go. Inputs include root/base/canonical patterns AND fold-key ignored snapshot. After preflight run exactly three data invocations with gitCandidatePathspecs: status --porcelain=v2 -z --branch --untracked-files=all --ignored=no --no-renames --ignore-submodules=all -- <candidates>; ls-files -z -o -i --exclude-standard -- <candidates>; ls-files -v -z -- <candidates>. All data legs and preflight share ONE protectedCheckTimeout=2500ms context, cmd.Dir=root, minimal gitEnv, optional locks off and core.fsmonitor=false. The preflight is an additional config query, not a fourth candidate stream; do not drop this no-host-filter-execution guard. Feed all three capped outputs to M1c parseProtectedCandidates, never duplicate parsers. HEAD==base and zero filtered changes are required for intact; moved HEAD emits head-moved even with empty paths. Each invocation enforces 8MiB, timeout/nonzero/malformed yields error with intact null, report-only. Probe fixtures capture their ignored snapshot before the probe mutation; M6a owns episode construction. No production result field yet. V24 will measure complete preflight+three-leg cost and resize before M6a.

One shared 2500ms provisional phase bound, including filter preflight plus THREE candidate data calls; later re-sized from V24. M1c pure parser owns syntax and filtering; no duplicate parseProtectedStatus/parseProtectedFlags.

**Named tests and assertions:**

- `TestProtectedChecker` (M4): Independent status/flags fixtures: clean, modified/deleted/untracked/ignored, lower-case/S flags, HEAD moved, malformed output.
- `TestProtectedCheckerRootAndBound` (M4): Fake Git records root, minimal environment, widened argv and all three data legs plus filter preflight under one shared 2500ms deadline; output/nonzero failures produce diagnostic/null.

**Mutation → exact killer (one compiling mutant at a time):**

- **MUT-NO-INDEX-FLAG-CHECK** → `TestProtectedChecker/index-flag`: Remove ls-files -v leg; assertion: lower-case and S tags yield index-flag false
- **MUT-HEAD-UNCHECKED** → `TestProtectedChecker/head-moved`: Ignore branch oid; assertion: HEAD change produces head-moved false even with zero paths
- **MUT-IGNORED-NO** → `TestProtectedChecker/ignored`: Remove ls-files -o -i ignored-file leg (status --ignored=no is intentional); assertion: ignored protected creation listed
- **MUT-CHECK-PARSE-OPEN** → `TestProtectedChecker/malformed`: Accept malformed git output; assertion: error set and intact null
- **MUT-CHECK-CWD** → `TestProtectedCheckerRootAndBound/cwd`: Omit root Dir / minimal git environment; assertion: cmd.Dir is worktree and optional locks/filters disabled
- **MUT-CHECK-DOUBLE-BUDGET** → `TestProtectedCheckerRootAndBound/shared-phase-bound`: Give each subprocess a fresh 2500ms deadline; assertion: Preflight plus all THREE data calls share one phase deadline, rather than one per command.

Each mutation uses its exact `go test -count=1 ./host/broker/ -run '^<full killer above>$' -timeout 300s -v` selector (literal commands in JSON); these mutation gates are **NOT MEASURED: named additions absent at base; zero tests is not a pass**. One valid guard/data mutant at a time, exact named killer/subtest and assertion. Save and restore allowed bytes in worktree-local scratch without git. Require named red, then restoration green. No compile error, panic, timeout, skip, missing executable, socket denial, unrelated refusal or grep-only assertion is a kill. Fixed independent expectation must be reachable after the changed guard; delete/move masking cannot satisfy r3 bare-D control.

**Exact gates** (WORLD_TOOL_AILANG_BIN unset; TMPDIR/GOTMPDIR worktree-local):

```sh
go test -count=1 ./host/broker/... -run '^TestProtectedChecker$' -timeout 120s -v
go test -count=1 ./host/broker/... -run '^TestProtectedCheckerRootAndBound$' -timeout 120s -v
go test -race -count=1 ./host/broker/...
go vet ./host/broker/
```

**Baseline for every gate at 451dc3d:**

- `go test -count=1 ./host/broker/... -run '^TestProtectedChecker$' -timeout 120s -v` — NEW TEST/ROW absent at base; NOT RUN. Zero tests would not be a baseline pass.
- `go test -count=1 ./host/broker/... -run '^TestProtectedCheckerRootAndBound$' -timeout 120s -v` — NEW TEST/ROW absent at base; NOT RUN. Zero tests would not be a baseline pass.
- `go test -race -count=1 ./host/broker/...` — NOT MEASURED; green on this base is UNKNOWN. Unfiltered suites include legacy off-worktree writes and/or sockets; do not execute here or infer green from the focused baseline. Any socket result is UNINFORMATIVE UNDER SANDBOX.
- `go vet ./host/broker/` — GREEN rc 0 at 451dc3d (iteration 282); package-only vet, tool pin unset, scratch contained.

No `.ail` changes; no AIL verification. Full package acceptance requires contained controller runner; sockets are UNINFORMATIVE UNDER SANDBOX.

### M4c — Wire startup arm 8 for prevention and Linux detection

Design M4; **245 LOC total / 30 min**, own PR, depends on M4b. Status **planned**.

**Files touched and base functions:**

- `host/broker/exec_probe.go:44 arm constants; :294 probeExecProfile; :365 execProber.probe; :199 execRefusalCodes; :252 execProber.run`
- `host/broker/exec_probe_test.go:84 newProbeFixture; :167 TestExecProbeArmNames; :178 passthrough test; :242 missing-client test`
- `host/broker/exec_probe_srt_test.go:22 TestExecSrtStartupProbePasses`

Allowlist: `host/broker/exec_probe.go`, `host/broker/exec_probe_test.go`, `host/broker/exec_probe_srt_test.go`.

**Implementation and landable boundary.** Add arm8-protected-write only for declared sets; create sample matching names from first directory/basename patterns in its scratch worktree, initialize the dedicated scratch git repo before base capture and invoke M4b checker. Arbitrary valid globs need an explicit matching sample generator tested with ?, classes and mixed case; no hardcoded *_test.go only. Darwin requires refusal token+unchanged/absent host data; Linux basename write requires success token and intact=false. Pin-free arm-specific harness must avoid arms 1–7 external /tmp controls and listeners. Add explicit legacy seven-arm order/no-protected assertions.

Always-on companions directly drive arm8 oracle/runner using shell/fake tokens; they must not require node, real srt, tool pin or network. Real srt integration is a separately named optional arm. Keep production probe lifecycle through runBounded.

**Named tests and assertions:**

- `TestExecProbeArm8`: Real srt: directory refusal, Darwin basename creation refusal, Linux basename allowed then detected. Skip absent srt/bwrap.
- `TestExecProbeArm8FailsUnsandboxed`: Pin-free pass-through and missing-client fail specifically arm8; no raw rc accepted.
- `TestExecProbeArm8Absent`: No protected set runs exactly arms1–7 and never invokes arm8.
- `TestExecProbeProtectedSample`: Each valid pattern class creates a matching probe specimen; unsupported construction fails explicitly.

**One named negative control per guard (independent mutation runs):**

- `TestExecProbeArm8FailsUnsandboxed/missing-client` kills **MUT-ARM8-RC-ONLY** (Accept any nonzero rc without refusal token); assertion: missing client cannot satisfy arm8.
- `TestExecProbeArm8FailsUnsandboxed/forged-token` kills **MUT-ARM8-NO-HOST-CHECK** (Omit host byte comparison); assertion: changed host bytes fail despite refusal text.
- `TestExecProbeArm8FailsUnsandboxed/linux-detect` kills **MUT-ARM8-LINUX-NO-DETECT** (Do not run Linux checker after allowed basename write); assertion: Linux detection leg must return false.
- `TestExecProbeArm8Absent` kills **MUT-ARM8-ALWAYS** (Run arm8 for empty set); assertion: no arm8 for absent set.
- `TestExecProbeProtectedSample/class` kills **MUT-ARM8-SAMPLE** (Hardcode a specimen that does not match ?/class glob); assertion: probe specimen must match declared pattern.

**Exact gates** (tool pin unset; mandatory package scopes below):

```sh
go test -count=1 ./host/broker/... -run '^TestExecProbeArm8$' -timeout 120s -v
go test -count=1 ./host/broker/... -run '^TestExecProbeArm8FailsUnsandboxed$' -timeout 120s -v
go test -count=1 ./host/broker/... -run '^TestExecProbeArm8Absent$' -timeout 120s -v
go test -count=1 ./host/broker/... -run '^TestExecProbeProtectedSample$' -timeout 120s -v
go test -race -count=1 ./host/broker/...
go vet ./...
```

**Baseline for every gate:**

- `go test -count=1 ./host/broker/... -run '^TestExecProbeArm8$' -timeout 120s -v` — NEW TEST/ROW absent at base; NOT RUN. Zero tests would not be a baseline pass.
- `go test -count=1 ./host/broker/... -run '^TestExecProbeArm8FailsUnsandboxed$' -timeout 120s -v` — NEW TEST/ROW absent at base; NOT RUN. Zero tests would not be a baseline pass.
- `go test -count=1 ./host/broker/... -run '^TestExecProbeArm8Absent$' -timeout 120s -v` — NEW TEST/ROW absent at base; NOT RUN. Zero tests would not be a baseline pass.
- `go test -count=1 ./host/broker/... -run '^TestExecProbeProtectedSample$' -timeout 120s -v` — NEW TEST/ROW absent at base; NOT RUN. Zero tests would not be a baseline pass.
- `go test -race -count=1 ./host/broker/...` — NOT MEASURED; green on this base is UNKNOWN. Unfiltered suites include legacy off-worktree writes and/or sockets; do not execute here or infer green from the focused baseline. Any socket result is UNINFORMATIVE UNDER SANDBOX.
- `go vet ./...` — GREEN: rc 0, 1.46 s, measured on this base with WORLD_TOOL_AILANG_BIN unset.

No `.ail` change; no AIL verification command. New fixture writes are worktree-local; no seeded external state. Any daemon/REST table deadline is created inside its own `t.Run`.

### M5a — Render protected AILANG policies in every variant

Design M5; **245 LOC total / 30 min**, own PR, depends on M4c. Status **planned**.

**Files touched and base functions:**

- `host/broker/handlers_ailang.go:165 RenderEpisodePolicy; :242 AilangToolConfig; :287 NewAilangToolHandler`
- `host/broker/handlers_ailang_run.go:165 runPolicySpec; :178 renderPolicy; :215 tomlStrings; :258 RenderRunVariant; :320 runVariantFor; :430 VerifyRunCaps`
- `host/broker/handlers_ailang_test.go:39 TestRenderEpisodePolicyHasExactlyTheDesignKeySet; :585 fixture call`
- `host/broker/handlers_ailang_run_test.go:61 TestRenderRunVariantKeySetsPerClass`
- `host/daemon/workspace.go:507 buildAilangHandler (mechanical nil argument here, actual wiring M5b)`
- `host/daemon/workspace_test.go:302 TestWorkspacePolicyRenderedPerEpisode`
- `host/daemon/workspace_layout.go:40 checkModuleRoot`
- `host/daemon/workspace_layout_test.go:104 TestWorkspaceModuleRootIsTheSandbox; :545 TestWorkspacePackageCacheLinkedIntoEpisodeHome`

Allowlist: `host/broker/handlers_ailang.go`, `host/broker/handlers_ailang_run.go`, `host/broker/handlers_ailang_test.go`, `host/broker/handlers_ailang_run_test.go`, `host/daemon/workspace.go`, `host/daemon/workspace_test.go`, `host/daemon/workspace_layout.go`, `host/daemon/workspace_layout_test.go`.

**Implementation and landable boundary.** Change RenderEpisodePolicy(root, protected); update ALL caller sites in the same compiling PR with nil where not yet wired. Add protected plumbing to AilangToolConfig and run-policy spec; append fresh copied slices to every base/Declassify/Net/Env render, cached runVariantFor and VerifyRunCaps. Extend RenderRunVariant consistently and update its callers. Encode TOML quoted pattern strings safely: grammar permits quote characters, while current tomlStrings interpolates unescaped text. Preserve existing bytes for old safe strings and empty protection. Pre-rebased patterns are passed, not full worktree paths.

D-156-4=B: no World pre-check, Nlink guard or new typed refusal. D-156-4=A branch is deferred, not silently implemented. Mechanical caller changes across all eight existing files are included in LOC budget to keep this PR compiling.

**Named tests and assertions:**

- `TestEpisodePolicyProtectedGolden`: Independent exact expected TOML including safe encoding of a quote-bearing basename; no golden self-computation.
- `TestEpisodePolicyProtectedAllVariants`: base, Declassify, Net, Env all contain exact protected append and retain prior fields.
- `TestEpisodePolicyNoProtectedIsByteIdentical`: nil and [] produce existing base/variant bytes.

**One named negative control per guard (independent mutation runs):**

- `TestEpisodePolicyProtectedGolden` kills **MUT-NO-POLICY-APPEND** (Remove policy append); assertion: fs_deny_write contains exact protected list.
- `TestEpisodePolicyProtectedAllVariants/Env` kills **MUT-VARIANT-MISSED** (Omit Env or Net append); assertion: Env policy retains protected entries.
- `TestEpisodePolicyProtectedGolden/quote` kills **MUT-POLICY-QUOTE** (Fail to escape TOML quote); assertion: quoted pattern decodes as original string in valid TOML.
- `TestEpisodePolicyNoProtectedIsByteIdentical` kills **MUT-POLICY-EMPTY-BYTES** (Change old policy bytes with nil patterns); assertion: base and all variants match prior bytes.

**Exact gates** (tool pin unset; mandatory package scopes below):

```sh
go test -count=1 ./host/broker/... -run '^TestEpisodePolicyProtectedGolden$' -timeout 120s -v
go test -count=1 ./host/broker/... -run '^TestEpisodePolicyProtectedAllVariants$' -timeout 120s -v
go test -count=1 ./host/broker/... -run '^TestEpisodePolicyNoProtectedIsByteIdentical$' -timeout 120s -v
go test -race -count=1 ./host/broker/...
go test -race -count=1 ./host/daemon/...
go vet ./...
```

**Baseline for every gate:**

- `go test -count=1 ./host/broker/... -run '^TestEpisodePolicyProtectedGolden$' -timeout 120s -v` — NEW TEST/ROW absent at base; NOT RUN. Zero tests would not be a baseline pass.
- `go test -count=1 ./host/broker/... -run '^TestEpisodePolicyProtectedAllVariants$' -timeout 120s -v` — NEW TEST/ROW absent at base; NOT RUN. Zero tests would not be a baseline pass.
- `go test -count=1 ./host/broker/... -run '^TestEpisodePolicyNoProtectedIsByteIdentical$' -timeout 120s -v` — NEW TEST/ROW absent at base; NOT RUN. Zero tests would not be a baseline pass.
- `go test -race -count=1 ./host/broker/...` — NOT MEASURED; green on this base is UNKNOWN. Unfiltered suites include legacy off-worktree writes and/or sockets; do not execute here or infer green from the focused baseline. Any socket result is UNINFORMATIVE UNDER SANDBOX.
- `go test -race -count=1 ./host/daemon/...` — NOT MEASURED; green on this base is UNKNOWN. Unfiltered suites include legacy off-worktree writes and/or sockets; do not execute here or infer green from the focused baseline. Any socket result is UNINFORMATIVE UNDER SANDBOX.
- `go vet ./...` — GREEN: rc 0, 1.46 s, measured on this base with WORLD_TOOL_AILANG_BIN unset.

No `.ail` change; no AIL verification command. New fixture writes are worktree-local; no seeded external state. Any daemon/REST table deadline is created inside its own `t.Run`.

### M5b — Wire episode profile selection and module-root policy rebasing

Design M5; **235 LOC total / 30 min**, own PR, depends on M5a. Status **planned**.

**Files touched and base functions:**

- `host/daemon/workspace.go:154 profileFor; :507 buildAilangHandler; :522 policy construction`
- `host/daemon/workspace_protected_test.go NEW; existing fake binary host/daemon/workspace_test.go:95 fakeToolBin; module root host/daemon/workspace_layout.go:85 sandboxRoot`
- `host/broker/exec_protected_tools_test.go NEW; fixture anchors host/broker/handlers_ailang_test.go:39 and host/daemon/workspace_test.go:593 realToolBin`

Allowlist: `host/daemon/workspace.go`, `host/daemon/workspace_protected_test.go`, `host/broker/exec_protected_tools_test.go`.

**Implementation and landable boundary.** Select episode profile, rebase protected list from worktree to AILANG sandbox, pass to base policy and config for cached variants. Fake-tool tests assert policy bytes and spawn transcript directly without sockets. Optional pinned policy-tool replay uses v0.52.1 and V7 lexical write/edit/CLI/FS operations; symlink aliases remain a documented upstream residual under B. Run all REST/daemon table rows with contexts made inside t.Run.

Policy half only under D-156-4=B. Optional TestWorkspaceWriteProtectedPolicyTool has pin-free TestWorkspaceProtectedPolicyRequests plus M5a golden companion. A fake policy-tool refusal is not proof of the upstream matcher.

**Named tests and assertions:**

- `TestEpisodePolicyProtectedRebased`: Module sub rebase, ancestor->*, disjoint omitted, basename unchanged in selected episode; no cross-episode leak.
- `TestWorkspaceProtectedPolicyWiring`: Fake tool captures selected base and variant policy; no profile/empty stays byte-identical.
- `TestWorkspaceWriteProtectedPolicyTool`: Real tool refuses V7 lexical write/edit/CLI and FS operations; unset WORLD_TOOL_AILANG_BIN skips before setup.
- `TestWorkspaceProtectedPolicyRequests`: Always-on fake-tool request/variant policy transcript has exact append; this is policy wiring evidence only.

**One named negative control per guard (independent mutation runs):**

- `TestEpisodePolicyProtectedRebased/under-sub` kills **MUT-NO-REBASE** (Skip rebase at daemon boundary); assertion: sandbox policy contains tests/** instead of sub/tests/**.
- `TestWorkspaceProtectedPolicyWiring/episode-selection` kills **MUT-POLICY-WRONG-EPISODE** (Use default profile instead of episode selection); assertion: second episode receives its own protected set.
- `TestWorkspaceProtectedPolicyRequests` kills **MUT-NO-POLICY-APPEND** (Remove append while requests still forward); assertion: fake tool sees expected deny patterns.

**Exact gates** (tool pin unset; mandatory package scopes below):

```sh
go test -count=1 ./host/daemon/... -run '^TestEpisodePolicyProtectedRebased$' -timeout 120s -v
go test -count=1 ./host/daemon/... -run '^TestWorkspaceProtectedPolicyWiring$' -timeout 120s -v
go test -count=1 ./host/broker/... -run '^TestWorkspaceWriteProtectedPolicyTool$' -timeout 120s -v
go test -count=1 ./host/broker/... -run '^TestWorkspaceProtectedPolicyRequests$' -timeout 120s -v
go test -race -count=1 ./host/broker/...
go test -race -count=1 ./host/daemon/...
go vet ./...
```

**Baseline for every gate:**

- `go test -count=1 ./host/daemon/... -run '^TestEpisodePolicyProtectedRebased$' -timeout 120s -v` — NEW TEST/ROW absent at base; NOT RUN. Zero tests would not be a baseline pass.
- `go test -count=1 ./host/daemon/... -run '^TestWorkspaceProtectedPolicyWiring$' -timeout 120s -v` — NEW TEST/ROW absent at base; NOT RUN. Zero tests would not be a baseline pass.
- `go test -count=1 ./host/broker/... -run '^TestWorkspaceWriteProtectedPolicyTool$' -timeout 120s -v` — NEW TEST/ROW absent at base; NOT RUN. Zero tests would not be a baseline pass.
- `go test -count=1 ./host/broker/... -run '^TestWorkspaceProtectedPolicyRequests$' -timeout 120s -v` — NEW TEST/ROW absent at base; NOT RUN. Zero tests would not be a baseline pass.
- `go test -race -count=1 ./host/broker/...` — NOT MEASURED; green on this base is UNKNOWN. Unfiltered suites include legacy off-worktree writes and/or sockets; do not execute here or infer green from the focused baseline. Any socket result is UNINFORMATIVE UNDER SANDBOX.
- `go test -race -count=1 ./host/daemon/...` — NOT MEASURED; green on this base is UNKNOWN. Unfiltered suites include legacy off-worktree writes and/or sockets; do not execute here or infer green from the focused baseline. Any socket result is UNINFORMATIVE UNDER SANDBOX.
- `go vet ./...` — GREEN: rc 0, 1.46 s, measured on this base with WORLD_TOOL_AILANG_BIN unset.

No `.ail` change; no AIL verification command. New fixture writes are worktree-local; no seeded external state. Any daemon/REST table deadline is created inside its own `t.Run`.

### M6-pre — Measure V24 and size phase/construction bounds

Design M6; **200 LOC total / 30 min**, own PR/evidence boundary, depends on M5b, M1d. Status **planned**.

**Files + allowlist:** `host/broker/exec_protected_v24_test.go`, `design_docs/verification/world-row156-v24/README.md`, `design_docs/verification/world-row156-v24/result.json`, `design_docs/planned/w-exec-protected-test-paths.md`. No other files. Future files/tests are absent at base unless explicitly stated.

**Implementation and assertions.** Controller/executor measurement prerequisite, NOT performed by this planner. M6a may not start until design V24 is measured and the bounds re-sized from it. Design §6 M6 [r4.2] requires git worktree add of sunholo-platform and ailang with .venv/node_modules installed, measuring all three §4.5 legs plus construction-time snapshot together, warm and cold, using the REAL matchProtected filter, never V23's ASCII stand-in. Prepare authorized disposable source copies and V24_SCRATCH under the controller-owned worktree/scratch; the two worktree-add commands below operate only on those copies, never live source metadata. Populate dependency trees from existing local caches/offline copies before construction; no network required by this plan. Harness TestProtectedV24RealWorktree in new host/broker/exec_protected_v24_test.go requires WORLD_PROTECTED_V24_ROOT, SKIP with an explicit missing-measurement diagnostic if unset; a skipped gate never satisfies V24. Require both installed trees, current project profiles and real Git; fail on incomplete input. Record revs, toolchain/Git versions, installed dependency census, 20 samples per project per warm/cold class (first fresh-worktree access labelled cold; describe cache preparation, never claim OS-cache cold without evidence), per-leg and combined phase p95/max, output bytes and selected/filtered counts. Use existing M4b checker (filter preflight included) plus M1c's renderer/parser. Time construction as rev-parse + full ls-tree mirror filter + ignored enumeration and fold-key snapshot/digest; no M6a implementation dependency. Log exact argv: status --porcelain=v2 -z --branch --untracked-files=all --ignored=no --no-renames --ignore-submodules=all -- <candidates>; ls-files -z -o -i --exclude-standard -- <candidates>; ls-files -v -z -- <candidates>. Use gitEnv, GIT_OPTIONAL_LOCKS=0 and core.fsmonitor=false; preflight must precede status and be included in measured phase total. Provisional protectedCheckTimeout=2500ms and protectedBaselineTimeout=3200ms are NOT acceptance evidence: size phase and construction bounds separately from their V24 observed max ×1.5 (round upward in ms); retain at least the stated provisional bounds, increase if needed, amend design V24 and plan constants before M6a. Require every clean phase's intact verdict non-null/true and p95 below chosen phase bound, construction successful within its chosen bound; record null/error/CHANGED counts explicitly. Verify from M4 episode construction flow/logs whether ALL dependency installation precedes NewExecHandler. If not, demonstrate a new ignored post-snapshot dependency reports CHANGED and disclose R-156-6, never silently whitelist it. Evidence belongs in design_docs/verification/world-row156-v24/README.md and result.json, with controller-authored V24 design entry; no measurement or design amendment is fabricated here. <=200 changed LOC artifacts/harness, <=30 active minutes, split if needed. Run reject-null, reject-over-budget and required-legs synthetic subtests BEFORE the optional /measurement subtest checks WORLD_PROTECTED_V24_ROOT. V24 completion additionally requires /measurement to PASS for BOTH projects with nonempty recorded samples; parent PASS with /measurement SKIP is not acceptance.

Controller/executor only; requires separately authorized contained real-worktree measurement environment. Design V24 remains unmeasured in iteration 282. This is a hard dependency, not a timing estimate.

**Named tests and assertions:**

- `TestProtectedV24RealWorktree` (V24): Both installed worktrees measured with real mirror; clean intact non-null, p95 < resized bound; constructor max ×1.5 and install ordering recorded. Synthetic reject-null/reject-over-budget/required-legs controls are always-on.

**Mutation → exact killer (one compiling mutant at a time):**

- **MUT-V24-NULL-PASS** → `TestProtectedV24RealWorktree/reject-null`: Treat timed-out/null phase as measured intact; assertion: Injected null sample rejects acceptance.
- **MUT-V24-BUDGET-PASS** → `TestProtectedV24RealWorktree/reject-over-budget`: Accept p95 >= selected phase bound; assertion: Synthetic over-budget sample rejects acceptance.
- **MUT-V24-NO-IGNORED** → `TestProtectedV24RealWorktree/required-legs`: Drop ignored leg from combined timing; assertion: Missing ignored timing/census rejects acceptance.

Each mutation uses its exact `go test -count=1 ./host/broker/ -run '^<full killer above>$' -timeout 300s -v` selector (literal commands in JSON); these mutation gates are **NOT MEASURED: named additions absent at base; zero tests is not a pass**. One valid guard/data mutant at a time, exact named killer/subtest and assertion. Save and restore allowed bytes in worktree-local scratch without git. Require named red, then restoration green. No compile error, panic, timeout, skip, missing executable, socket denial, unrelated refusal or grep-only assertion is a kill. Fixed independent expectation must be reachable after the changed guard; delete/move masking cannot satisfy r3 bare-D control.

**Exact gates** (WORLD_TOOL_AILANG_BIN unset; TMPDIR/GOTMPDIR worktree-local):

```sh
git -C "$V24_SUNHOLO_SOURCE" worktree add --detach "$V24_SCRATCH/sunholo-platform" HEAD
git -C "$V24_AILANG_SOURCE" worktree add --detach "$V24_SCRATCH/ailang" HEAD
WORLD_PROTECTED_V24_ROOT="$V24_SCRATCH" go test -count=1 ./host/broker/ -run '^TestProtectedV24RealWorktree$' -timeout 300s -v
go test -race -count=1 ./host/broker/...
go vet ./host/broker/
```

**Baseline for every gate at 451dc3d:**

- `git -C "$V24_SUNHOLO_SOURCE" worktree add --detach "$V24_SCRATCH/sunholo-platform" HEAD` — NOT MEASURED: controller/executor V24 worktree measurement; no Git writes/fixtures permitted for planner; harness absent at base; zero tests is not a pass.
- `git -C "$V24_AILANG_SOURCE" worktree add --detach "$V24_SCRATCH/ailang" HEAD` — NOT MEASURED: controller/executor V24 worktree measurement; no Git writes/fixtures permitted for planner; harness absent at base; zero tests is not a pass.
- `WORLD_PROTECTED_V24_ROOT="$V24_SCRATCH" go test -count=1 ./host/broker/ -run '^TestProtectedV24RealWorktree$' -timeout 300s -v` — NOT MEASURED: controller/executor V24 worktree measurement; no Git writes/fixtures permitted for planner; harness absent at base; zero tests is not a pass.
- `go test -race -count=1 ./host/broker/...` — NOT MEASURED (PF6: legacy tests write real HOME); requires controller-contained runner. Socket-touching results are UNINFORMATIVE UNDER SANDBOX.
- `go vet ./host/broker/` — GREEN rc 0 at 451dc3d, iteration 282; WORLD_TOOL_AILANG_BIN unset, TMPDIR/GOTMPDIR worktree-local. Package-only vet measured; no repository-wide vet claim.

No `.ail` changes; no AIL verification. Full package acceptance requires contained controller runner; sockets are UNINFORMATIVE UNDER SANDBOX.

### M6a — Capture the immutable episode baseline and bounded record shape

Design M6; **245 LOC total / 30 min**, own PR/evidence boundary, depends on M6-pre. Status **planned**.

**Files + allowlist:** `host/broker/exec_protected_record.go`, `host/broker/exec_protected_record_test.go`, `host/broker/handlers_exec.go`, `host/broker/exec_protected_check.go`. No other files. Future files/tests are absent at base unless explicitly stated.

**Implementation and assertions.** M6a may not start until design V24 is measured and the bounds re-sized from it. Apply V24 bounds (at least protectedCheckTimeout=2500ms, protectedBaselineTimeout=3200ms), not pre-r4 constants. NewExecHandler once per episode/project, only nonempty patterns: one shared named protectedBaselineTimeout for rev-parse HEAD, ls-tree -r -z --full-tree <base> filtered by matchProtected (digest mode SP oid TAB path LF in tree order, base_files), AND candidate ls-files -z -o -i --exclude-standard followed by mirror filtering to an immutable fold-key set. Store base_ignored_digest=sha256 over sorted unique fold keys, each LF terminated, and base_ignored_files=count, alongside base/base_digest/base_files; reuse the same snapshot for every phase, never refresh it after commands. Empty snapshot hashes the empty byte string, not an absent field. Use gitEnv/runBounded and 8MiB stdout caps with no-filter preflight included under construction bound. Missing Git/non-Git, any malformed/oversized stream, nonzero, snapshot failure or construction timeout is retained episode error/null, never construction failure or intact. Reuse M1c protectedChange; do not redeclare it. Result shape includes nullable phase verdicts, ignored metadata, max32 changes, total/truncated and bounded path JSON; no result emission until M6b. All fixtures worktree-local; no sockets.

V24 is a hard prerequisite. <=245 total LOC /30 minutes includes construction snapshot tests and resized constants; split before exceeding cap.

**Named tests and assertions:**

- `TestExecProtectedBaseline` (M6): Independent committed-tree digest plus /ignored-snapshot sorted fold-key digest/count (case/NFC dedup, empty set), /snapshot-immutable across calls; committed digest ignores dirty bytes.
- `TestExecProtectedBaselineFailure` (M6): Missing Git/non-Git, /snapshot-error, /snapshot-over-cap and /construction-timeout: handler constructed, retained error, both intact verdicts null; one named baseline bound covers all construction legs.
- `TestExecProtectedWorstCaseFitsOutputCap` (M6): 32 maximally escaped path entries, both phases, errors and long baseline metadata stay below existing 1MiB cap; total reflects omitted paths.

**Mutation → exact killer (one compiling mutant at a time):**

- **MUT-BASE-WORKTREE** → `TestExecProtectedBaseline/dirty-worktree`: Hash worktree bytes instead of committed modes/oids; assertion: base digest stays committed-tree identity
- **MUT-BASE-DIGEST** → `TestExecProtectedBaseline/independent-digest`: Omit mode or change ordering/newline; assertion: independent digest equality
- **MUT-BASE-FAILS-HANDLER** → `TestExecProtectedBaselineFailure`: Fail construction on unavailable baseline; assertion: handler constructed with diagnostic
- **MUT-PROTECTED-UNBOUNDED** → `TestExecProtectedWorstCaseFitsOutputCap`: Do not cap or JSON-bound changed paths; assertion: at most32 entries and bounded encoded size
- **MUT-NO-IGNORED-SNAPSHOT** → `TestExecProtectedBaseline/ignored-snapshot`: Omit construction snapshot or digest/count; assertion: Independent ignored fold-key digest/count equal expected.
- **MUT-BASELINE-FRESH-BUDGET** → `TestExecProtectedBaselineFailure/construction-timeout`: Give each construction subprocess a fresh bound; assertion: Total construction shares one protectedBaselineTimeout and records error on exhaustion.

Each mutation uses its exact `go test -count=1 ./host/broker/ -run '^<full killer above>$' -timeout 300s -v` selector (literal commands in JSON); these mutation gates are **NOT MEASURED: named additions absent at base; zero tests is not a pass**. One valid guard/data mutant at a time, exact named killer/subtest and assertion. Save and restore allowed bytes in worktree-local scratch without git. Require named red, then restoration green. No compile error, panic, timeout, skip, missing executable, socket denial, unrelated refusal or grep-only assertion is a kill. Fixed independent expectation must be reachable after the changed guard; delete/move masking cannot satisfy r3 bare-D control.

**Exact gates** (WORLD_TOOL_AILANG_BIN unset; TMPDIR/GOTMPDIR worktree-local):

```sh
go test -count=1 ./host/broker/... -run '^TestExecProtectedBaseline$' -timeout 120s -v
go test -count=1 ./host/broker/... -run '^TestExecProtectedBaselineFailure$' -timeout 120s -v
go test -count=1 ./host/broker/... -run '^TestExecProtectedWorstCaseFitsOutputCap$' -timeout 120s -v
go test -race -count=1 ./host/broker/...
go vet ./host/broker/
```

**Baseline for every gate at 451dc3d:**

- `go test -count=1 ./host/broker/... -run '^TestExecProtectedBaseline$' -timeout 120s -v` — NOT MEASURED: absent at base; zero tests is not a pass. No fixture Git repositories were created by planner.
- `go test -count=1 ./host/broker/... -run '^TestExecProtectedBaselineFailure$' -timeout 120s -v` — NOT MEASURED: absent at base; zero tests is not a pass. No fixture Git repositories were created by planner.
- `go test -count=1 ./host/broker/... -run '^TestExecProtectedWorstCaseFitsOutputCap$' -timeout 120s -v` — NOT MEASURED: absent at base; zero tests is not a pass. No fixture Git repositories were created by planner.
- `go test -race -count=1 ./host/broker/...` — NOT MEASURED (PF6: legacy tests write real HOME); requires controller-contained runner. Socket-touching results are UNINFORMATIVE UNDER SANDBOX.
- `go vet ./host/broker/` — GREEN rc 0 at 451dc3d, iteration 282; WORLD_TOOL_AILANG_BIN unset, TMPDIR/GOTMPDIR worktree-local. Package-only vet measured; no repository-wide vet claim.

No `.ail` changes; no AIL verification. Full package acceptance requires contained controller runner; sockets are UNINFORMATIVE UNDER SANDBOX.

### M6b — Attach before/after integrity while preserving the command deadline

Design M6; **245 LOC total / 30 min**, own PR/evidence boundary, depends on M6a. Status **planned**.

**Files + allowlist:** `host/broker/handlers_exec.go`, `host/broker/exec_protected_record.go`, `host/broker/exec_protected_record_test.go`, `host/broker/handlers_exec_m2_test.go`. No other files. Future files/tests are absent at base unless explicitly stated.

**Implementation and assertions.** [r4.x] Emit protected omitempty only for nonempty profile, including base_ignored_digest/base_ignored_files and the immutable construction snapshot. Before-phase after stack verification, after-phase before worktree-delta finish: each uses M4b checkProtected with THREE data outputs via parseProtectedCandidates and stored ignored set, all within one shared protectedCheckTimeout=2500ms (or V24 resized bound), including no-filter preflight. Recompute child budget after preparation AND before-phase as min(timeout_ms, time.Until(originalDeadline)-protectedCheckTimeout). With 9s command cap/10s handler deadline the protected command gets about 7.3s after a warm before-phase (~0.2–0.65s; ~6.85–7.3s across that range), about 5s after cold 2.5s, less other elapsed preparation; without protected keep 9s unchanged. After-phase uses cancellation-independent bounded context capped by ORIGINAL handler deadline; never extend deadline or grant each subprocess a full phase. A failed construction snapshot makes both phase verdicts null for episode. Preserve before/after errors without overwrite, nonzero Git/timeout/over-cap/malformed all error/null; command rc/stdout still returned under report-only D-156-1. Enforced=glob/dirs-only; execDecode already uses execResult, no duplicate schema. Compare normalized legacy bytes, never fixed nondeterministic duration_ms.

V24-sized one-phase reserve, three candidate legs and preflight per phase. 2500ms/3200ms are provisional floors; use recorded V24 choices. Unprotected profile has no new Git calls and retains 9000ms command cap.

**Named tests and assertions:**

- `TestExecProtectedIntact` (M6): Shim command touches nothing: both true and independent baseline fields, per-OS enforced.
- `TestExecProtectedCheckTimeout` (M6): Fake Git sleeps beyond 2500ms (or V24 bound); /reserve asserts exact min(timeout_ms, remaining-after-before - phaseBound) including 9s/10s cases; /shared-phase asserts single deadline for all three legs; /null and /before-only assert error/null retention and command still succeeds within original deadline.
- `TestExecProtectedNotGit` (M6): Command succeeds, bounded error present and intact null.
- `TestExecNoProtectedIsByteIdentical` (M6): Absent/[] omit protected completely, execute no extra git and preserve normalized legacy result bytes.

**Mutation → exact killer (one compiling mutant at a time):**

- **MUT-NO-RESERVE** → `TestExecProtectedCheckTimeout/reserve`: Give child the full remaining deadline; assertion: after-check time remains inside original handler deadline
- **MUT-CHECK-FAILS-CALL** → `TestExecProtectedNotGit`: Fail the effect on check error; assertion: command exit0/stdout retained and handler succeeds
- **MUT-FIELD-ALWAYS** → `TestExecNoProtectedIsByteIdentical`: Emit field with empty set; assertion: protected key absent and no protected git invocations
- **MUT-CHECK-ERROR-TRUE** → `TestExecProtectedCheckTimeout/null`: Report a timeout as intact=true; assertion: timed-out phase intact is null
- **MUT-CHECK-ERROR-OVERWRITE** → `TestExecProtectedCheckTimeout/before-only`: Lose the before-phase error when after succeeds; assertion: before error remains visible and before intact null
- **MUT-PHASE-FRESH-BUDGET** → `TestExecProtectedCheckTimeout/shared-phase`: Reset context for each data leg; assertion: Combined phase times out at one shared bound; no threefold allowance.

Each mutation uses its exact `go test -count=1 ./host/broker/ -run '^<full killer above>$' -timeout 300s -v` selector (literal commands in JSON); these mutation gates are **NOT MEASURED: named additions absent at base; zero tests is not a pass**. One valid guard/data mutant at a time, exact named killer/subtest and assertion. Save and restore allowed bytes in worktree-local scratch without git. Require named red, then restoration green. No compile error, panic, timeout, skip, missing executable, socket denial, unrelated refusal or grep-only assertion is a kill. Fixed independent expectation must be reachable after the changed guard; delete/move masking cannot satisfy r3 bare-D control.

**Exact gates** (WORLD_TOOL_AILANG_BIN unset; TMPDIR/GOTMPDIR worktree-local):

```sh
go test -count=1 ./host/broker/... -run '^TestExecProtectedIntact$' -timeout 120s -v
go test -count=1 ./host/broker/... -run '^TestExecProtectedCheckTimeout$' -timeout 120s -v
go test -count=1 ./host/broker/... -run '^TestExecProtectedNotGit$' -timeout 120s -v
go test -count=1 ./host/broker/... -run '^TestExecNoProtectedIsByteIdentical$' -timeout 120s -v
go test -race -count=1 ./host/broker/...
go vet ./host/broker/
```

**Baseline for every gate at 451dc3d:**

- `go test -count=1 ./host/broker/... -run '^TestExecProtectedIntact$' -timeout 120s -v` — NOT MEASURED: absent at base; zero tests is not a pass. No fixture Git repositories were created by planner.
- `go test -count=1 ./host/broker/... -run '^TestExecProtectedCheckTimeout$' -timeout 120s -v` — NOT MEASURED: absent at base; zero tests is not a pass. No fixture Git repositories were created by planner.
- `go test -count=1 ./host/broker/... -run '^TestExecProtectedNotGit$' -timeout 120s -v` — NOT MEASURED: absent at base; zero tests is not a pass. No fixture Git repositories were created by planner.
- `go test -count=1 ./host/broker/... -run '^TestExecNoProtectedIsByteIdentical$' -timeout 120s -v` — NOT MEASURED: absent at base; zero tests is not a pass. No fixture Git repositories were created by planner.
- `go test -race -count=1 ./host/broker/...` — NOT MEASURED (PF6: legacy tests write real HOME); requires controller-contained runner. Socket-touching results are UNINFORMATIVE UNDER SANDBOX.
- `go vet ./host/broker/` — GREEN rc 0 at 451dc3d, iteration 282; WORLD_TOOL_AILANG_BIN unset, TMPDIR/GOTMPDIR worktree-local. Package-only vet measured; no repository-wide vet claim.

No `.ail` changes; no AIL verification. Full package acceptance requires contained controller runner; sockets are UNINFORMATIVE UNDER SANDBOX.

### M6c — Kill tamper, ignored-file, flag and moved-HEAD mutants; pin replay

Design M6; **245 LOC total / 30 min**, own PR/evidence boundary, depends on M6b. Status **planned**.

**Files + allowlist:** `host/broker/exec_protected_record_test.go`, `host/broker/handlers_exec_srt_test.go`, `host/daemon/exec_e2e_srt_test.go`. No other files. Future files/tests are absent at base unless explicitly stated.

**Implementation and assertions.** Verified base anchors: host/broker/exec_protected_record_test.go NEW; reuse host/broker/worktree_delta_test.go:670 TestExecHandlerWorktreeDelta fixture style; host/broker/handlers_exec_srt_test.go:359 TestExecSrtGitHeadIsNotWritable; host/daemon/exec_e2e_srt_test.go:118 execRequest; :249 TestExecSrtReplayIsByteEqual. Through a pass-through fake srt, test creation, ignored creation via .gitignore and info/exclude, deletion/rename/replacement, sticky before=false on next call, host-set assume-unchanged/skip-worktree and host commit moving HEAD. Mutations and git commits occur ONLY in disposable fixture repositories under worktree. Real-srt linked-worktree regression refuses index writes and commit and observes flag/HEAD unchanged; skip absent engine. Add pin-free broker record replay companion and extend optional daemon replay bytes without reexecution. Old records have no protected field and remain readable. [r4.1/r4.2] Every phase exercises all three candidate legs and parseProtectedCandidates with immutable episode-start ignored set. Add TestExecProtectedAmbientIgnoredIntact with protected test_*.py: ignored .venv/lib/site-packages/pkg/test_a.py exists BEFORE construction; /standing gets true/true with snapshot digest/count; /new creates test_b.py AFTER construction and gets after=false, ignored. /snapshot-mutation documents existing snapshot file content edits remain exempt (R-156-6); dependency installs after construction report CHANGED, never whitelist/refresh the snapshot. Repeat .gitignore and info/exclude creation controls against snapshot. Replay asserts ignored baseline fields stored byte-for-byte; no re-enumeration on replay.

Always-on flag and head tests are the pin-free companions for real-srt index prevention, proving detection rather than prevention. TestExecProtectedRecordReplay is the pin-free companion for optional daemon real-srt replay. Real E2E requires WORLD_TOOL_AILANG_BIN; absence must skip before any setup. Any REST table has its own context per subtest.

**Named tests and assertions:**

- `TestExecProtectedDetectsCreation` (M6): New zz_test.go yields after=false/untracked; next call before=false; rc0 retained.
- `TestExecProtectedDetectsIgnoredCreation` (M6): Both .gitignore and info/exclude-hidden creations appear as ignored.
- `TestExecProtectedIndexFlag` (M6): Host-set assume-unchanged and skip-worktree plus modification produce before=false/index-flag when status alone empty.
- `TestExecProtectedHeadMoved` (M6): Host commit changes HEAD but status empty; before=false and HEAD/head-moved.
- `TestExecSrtIndexWriteRefused` (M6): Real srt: linked worktree index flag and HEAD unchanged after update-index/commit attempts.
- `TestExecProtectedRecordReplay` (M6): Pin-free broker normal record/recovery path returns exact saved protected bytes and no exec spawn; old unprotected result also replays.
- `TestExecSrtReplayIsByteEqual` (M6): Extend existing optional srt/tool E2E with protected field; exact output bytes and no reexecution.
- `TestExecProtectedAmbientIgnoredIntact` (AC1c.5/M6): Standing .venv protected file snapshot gives true/true; newly created sibling gives false/ignored, stored baseline unchanged.

**Mutation → exact killer (one compiling mutant at a time):**

- **MUT-ALWAYS-INTACT** → `TestExecProtectedDetectsCreation`: Return always intact; assertion: intact_after false and untracked zz_test.go
- **MUT-IGNORED-NO** → `TestExecProtectedDetectsIgnoredCreation/info-exclude`: Remove ls-files -o -i ignored leg; status --ignored=no is correct by design; assertion: ignored new protected file must be named
- **MUT-NO-INDEX-FLAG-CHECK** → `TestExecProtectedIndexFlag/skip-worktree`: Remove ls-files check; assertion: status-empty modified flagged file must be false/index-flag
- **MUT-HEAD-UNCHECKED** → `TestExecProtectedHeadMoved`: Ignore branch.oid; assertion: empty status with moved HEAD still false
- **MUT-PROTECTED-REPLAY** → `TestExecProtectedRecordReplay`: Drop saved protected output or reexecute during replay; assertion: saved result bytes equal and spawn counter unchanged
- **MUT-NO-IGNORED-BASELINE** → `TestExecProtectedAmbientIgnoredIntact/standing`: Ignore episode snapshot; assertion: Existing site-packages test_a.py does not cause CHANGED; both phase verdicts true.

Each mutation uses its exact `go test -count=1 ./host/broker/ -run '^<full killer above>$' -timeout 300s -v` selector (literal commands in JSON); these mutation gates are **NOT MEASURED: named additions absent at base; zero tests is not a pass**. One valid guard/data mutant at a time, exact named killer/subtest and assertion. Save and restore allowed bytes in worktree-local scratch without git. Require named red, then restoration green. No compile error, panic, timeout, skip, missing executable, socket denial, unrelated refusal or grep-only assertion is a kill. Fixed independent expectation must be reachable after the changed guard; delete/move masking cannot satisfy r3 bare-D control.

**Exact gates** (WORLD_TOOL_AILANG_BIN unset; TMPDIR/GOTMPDIR worktree-local):

```sh
go test -count=1 ./host/broker/... -run '^TestExecProtectedDetectsCreation$' -timeout 120s -v
go test -count=1 ./host/broker/... -run '^TestExecProtectedDetectsIgnoredCreation$' -timeout 120s -v
go test -count=1 ./host/broker/... -run '^TestExecProtectedIndexFlag$' -timeout 120s -v
go test -count=1 ./host/broker/... -run '^TestExecProtectedHeadMoved$' -timeout 120s -v
go test -count=1 ./host/broker/... -run '^TestExecSrtIndexWriteRefused$' -timeout 120s -v
go test -count=1 ./host/broker/... -run '^TestExecProtectedRecordReplay$' -timeout 120s -v
go test -count=1 ./host/daemon/... -run '^TestExecSrtReplayIsByteEqual$' -timeout 120s -v
go test -count=1 ./host/broker/ -run '^TestExecProtectedAmbientIgnoredIntact$' -timeout 300s -v
go test -race -count=1 ./host/broker/...
go test -race -count=1 ./host/daemon/...
go vet ./host/broker/
```

**Baseline for every gate at 451dc3d:**

- `go test -count=1 ./host/broker/... -run '^TestExecProtectedDetectsCreation$' -timeout 120s -v` — NOT MEASURED: absent at base; zero tests is not a pass. No fixture Git repositories were created by planner.
- `go test -count=1 ./host/broker/... -run '^TestExecProtectedDetectsIgnoredCreation$' -timeout 120s -v` — NOT MEASURED: absent at base; zero tests is not a pass. No fixture Git repositories were created by planner.
- `go test -count=1 ./host/broker/... -run '^TestExecProtectedIndexFlag$' -timeout 120s -v` — NOT MEASURED: absent at base; zero tests is not a pass. No fixture Git repositories were created by planner.
- `go test -count=1 ./host/broker/... -run '^TestExecProtectedHeadMoved$' -timeout 120s -v` — NOT MEASURED: absent at base; zero tests is not a pass. No fixture Git repositories were created by planner.
- `go test -count=1 ./host/broker/... -run '^TestExecSrtIndexWriteRefused$' -timeout 120s -v` — NOT MEASURED: absent at base; zero tests is not a pass. No fixture Git repositories were created by planner.
- `go test -count=1 ./host/broker/... -run '^TestExecProtectedRecordReplay$' -timeout 120s -v` — NOT MEASURED: absent at base; zero tests is not a pass. No fixture Git repositories were created by planner.
- `go test -count=1 ./host/daemon/... -run '^TestExecSrtReplayIsByteEqual$' -timeout 120s -v` — NOT MEASURED: existing optional daemon E2E, protected assertions absent; socket-touching results UNINFORMATIVE UNDER SANDBOX.
- `go test -count=1 ./host/broker/ -run '^TestExecProtectedAmbientIgnoredIntact$' -timeout 300s -v` — NOT MEASURED: absent at base; zero tests is not a pass. No fixture Git repositories were created by planner.
- `go test -race -count=1 ./host/broker/...` — NOT MEASURED (PF6: legacy tests write real HOME); requires controller-contained runner. Socket-touching results are UNINFORMATIVE UNDER SANDBOX.
- `go test -race -count=1 ./host/daemon/...` — NOT MEASURED (PF6: legacy tests write real HOME); requires controller-contained runner. Socket-touching results are UNINFORMATIVE UNDER SANDBOX.
- `go vet ./host/broker/` — GREEN rc 0 at 451dc3d, iteration 282; WORLD_TOOL_AILANG_BIN unset, TMPDIR/GOTMPDIR worktree-local. Package-only vet measured; no repository-wide vet claim.

No `.ail` changes; no AIL verification. Full package acceptance requires contained controller runner; sockets are UNINFORMATIVE UNDER SANDBOX.

### M7a — Show intact, changed, unchecked and unprotected runs in why

Design M7; **230 LOC total / 30 min**, own PR, depends on M6c. Status **planned**.

**Files touched and base functions:**

- `cmd/ailang-worldd/why.go:571 whyLink; :585 resultWorktreeDelta pattern; :701 buildChain; :817 lazy output decode; :865 printChain`
- `cmd/ailang-worldd/why_protected_test.go NEW; existing fixture anchor cmd/ailang-worldd/why_test.go:370 TestWhyShowsTheWorktreeDelta`

Allowlist: `cmd/ailang-worldd/why.go`, `cmd/ailang-worldd/why_protected_test.go`.

**Implementation and landable boundary.** Add lazy resultProtected, protectedLine and optional JSON block to whyLink for exec results. Carry explicit exec-result identity so none-declared line does not print for arbitrary non-exec records. Text reports intact with base/count/digest/enforcement, CHANGED with named status/phase paths, not checked with error, or none declared for old/no-set exec results. --json retains the complete result block including truncation and null intact. Pure renderer tests need no daemon or pin; any end-to-end REST rows create their deadline inside t.Run.

The host/broker race gate is included for the result type; cmd/ailang-worldd race gate is additionally required because changed CLI code lives outside host. No .ail changes and no verify_ail gate.

**Named tests and assertions:**

- `TestWhyProtectedLine`: intact, changed-before/after, index-flag, head-moved, error, old-record/none-declared, and non-exec row.
- `TestWhyProtectedJSON`: Whole block carried with patterns, digest/base, enforced, nulls and truncation; old record decodes.

**One named negative control per guard (independent mutation runs):**

- `TestWhyProtectedLine/none-declared` kills **MUT-WHY-NONE-SILENT** (Omit absent-field exec line); assertion: protected: none declared must print.
- `TestWhyProtectedLine/not-checked` kills **MUT-WHY-STATE** (Render false/null as intact); assertion: error/null renders not checked, never intact.
- `TestWhyProtectedJSON` kills **MUT-WHY-JSON** (Drop block from JSON); assertion: entire protected payload survives JSON.
- `TestWhyProtectedLine/non-exec` kills **MUT-WHY-NONEXEC** (Print none declared for every old non-exec result); assertion: unrelated result gets no protected line.

**Exact gates** (tool pin unset; mandatory package scopes below):

```sh
go test -count=1 ./cmd/ailang-worldd/... -run '^TestWhyProtectedLine$' -timeout 120s -v
go test -count=1 ./cmd/ailang-worldd/... -run '^TestWhyProtectedJSON$' -timeout 120s -v
go test -race -count=1 ./host/broker/...
go test -race -count=1 ./cmd/ailang-worldd/...
go vet ./...
```

**Baseline for every gate:**

- `go test -count=1 ./cmd/ailang-worldd/... -run '^TestWhyProtectedLine$' -timeout 120s -v` — NEW TEST/ROW absent at base; NOT RUN. Zero tests would not be a baseline pass.
- `go test -count=1 ./cmd/ailang-worldd/... -run '^TestWhyProtectedJSON$' -timeout 120s -v` — NEW TEST/ROW absent at base; NOT RUN. Zero tests would not be a baseline pass.
- `go test -race -count=1 ./host/broker/...` — NOT MEASURED; green on this base is UNKNOWN. Unfiltered suites include legacy off-worktree writes and/or sockets; do not execute here or infer green from the focused baseline. Any socket result is UNINFORMATIVE UNDER SANDBOX.
- `go test -race -count=1 ./cmd/ailang-worldd/...` — NOT MEASURED; green on this base is UNKNOWN. Unfiltered suites include legacy off-worktree writes and/or sockets; do not execute here or infer green from the focused baseline. Any socket result is UNINFORMATIVE UNDER SANDBOX.
- `go vet ./...` — GREEN: rc 0, 1.46 s, measured on this base with WORLD_TOOL_AILANG_BIN unset.

No `.ail` change; no AIL verification command. New fixture writes are worktree-local; no seeded external state. Any daemon/REST table deadline is created inside its own `t.Run`.

### M7b — Ship executable protected-profile guidance and residuals

Design M7; **210 LOC total / 25 min**, own PR, depends on M7a. Status **planned**.

**Files touched and base functions:**

- `docs/QUICKSTART.md:640 section 10 profile generator; :914 section 11 module roots`
- `host/broker/exec_quickstart_test.go:22 quickstartSection10; :35 TestQuickstartSection10ProfilesLoad`

Allowlist: `docs/QUICKSTART.md`, `host/broker/exec_quickstart_test.go`.

**Implementation and landable boundary.** Update the existing Python profile generator and prose together with grammar/rendering table, Linux dirs-only detection, record/why examples and all four open decisions/defaults. Worked sets: ailang *_test.go, *.golden, tests/**; TwilightGame tests/**, *.test.ts, vitest.config.ts; sunholo-cli cli/tests/**, conftest.py, pyproject.toml (basename matches every such config). Explicitly distinguish worktree-relative cli/tests from module sandbox tests. Existing profile loader test runs generator verbatim with EXEC/HOME/tool fake paths wholly inside fixture and asserts exact protected sets. No new .ail snippet or package changes.

Use hash-before/hash-after observations only for real guard files/state. No seeding under real HOME, no live daemon drill, no publisher/mint operation. These docs complete the new user-facing surface in S7; earlier policy/profile PRs include their local refusal examples.

**Named tests and assertions:**

- `TestQuickstartSection10ProfilesLoad`: Execute checked-in generator, load all three profiles, assert exact canonical protected sets and existing flags/probes/read-roots.
- `TestQuickstartProtectedContract`: Bind dirs-only/report-only/empty-default/alias residual wording to executable outputs and examples, with parsed profile rows, not grep alone.

**One named negative control per guard (independent mutation runs):**

- `TestQuickstartSection10ProfilesLoad/protected-go` kills **MUT-DOC-GO-PATTERN** (Drop Go basename test pattern from runnable generator); assertion: Go profile includes *_test.go.
- `TestQuickstartSection10ProfilesLoad/protected-python` kills **MUT-DOC-ROOT** (Use tests/** for CLI worktree pattern instead of cli/tests/**); assertion: CLI profile worktree-relative protected set exact.
- `TestQuickstartProtectedContract` kills **MUT-DOC-PROMISE** (Claim Linux basename prevention / introduce protected default); assertion: documented behavior agrees with Linux renderer and no-set result.

**Exact gates** (tool pin unset; mandatory package scopes below):

```sh
go test -count=1 ./host/broker/... -run '^TestQuickstartSection10ProfilesLoad$' -timeout 120s -v
go test -count=1 ./host/broker/... -run '^TestQuickstartProtectedContract$' -timeout 120s -v
go test -race -count=1 ./host/broker/...
go vet ./...
```

**Baseline for every gate:**

- `go test -count=1 ./host/broker/... -run '^TestQuickstartSection10ProfilesLoad$' -timeout 120s -v` — Existing test PASS in measured focused broker baseline (also under race), with tool pin unset. Exact standalone command not run. New protected-set assertions are absent at base.
- `go test -count=1 ./host/broker/... -run '^TestQuickstartProtectedContract$' -timeout 120s -v` — NEW TEST/ROW absent at base; NOT RUN. Zero tests would not be a baseline pass.
- `go test -race -count=1 ./host/broker/...` — NOT MEASURED; green on this base is UNKNOWN. Unfiltered suites include legacy off-worktree writes and/or sockets; do not execute here or infer green from the focused baseline. Any socket result is UNINFORMATIVE UNDER SANDBOX.
- `go vet ./...` — GREEN: rc 0, 1.46 s, measured on this base with WORLD_TOOL_AILANG_BIN unset.

No `.ail` change; no AIL verification command. New fixture writes are worktree-local; no seeded external state. Any daemon/REST table deadline is created inside its own `t.Run`.

## Mutation protocol and residual acceptance

One valid guard/data mutant at a time, exact named killer/subtest and assertion. Save and restore allowed bytes in worktree-local scratch without git. Require named red, then restoration green. No compile error, panic, timeout, skip, missing executable, socket denial, unrelated refusal or grep-only assertion is a kill. Fixed independent expectation must be reachable after the changed guard; delete/move masking cannot satisfy r3 bare-D control.

All design §7 active-policy mutants are assigned above, including r2 flags/HEAD and r3 bare-D/top controls. **MUT-PRECHECK-LEXICAL-ONLY** → `TestWorkspaceWriteProtectedPrecheck/symlink-dir` and **MUT-NO-NLINK** → `TestWorkspaceWriteProtectedPrecheck/hardlink` are explicitly conditional on D-156-4=A, not missing implemented guards under B. The optional A work must be freshly split into M5c/M5d and add lexical/file-alias/hardlink spawn-counter-zero controls before it can land. The current JSON has no speculative milestones for that unanswered decision.

Real guard observations compare before/after SHA-256 (or stable absence) without planting decoys. Output artifacts, fake homes, temporary repos and restored mutants remain inside the worktree; remove them after verification. No work under `~/.ailang/state`, no `$HOME` snapshot directive, no source-index/HEAD writes, no live release action. The final controller collects CI evidence, resolves the legacy containment prerequisite and owns external upstream filing/mission-control note at landing.

The unchecked/changed record does not fail a successful command under D-156-1=A. Linux basename creation is detected only; aliases remain residual under D-156-4=B; no semantic “every test ran” proof is claimed. A cold git timeout is visible error/null, never intact. Each required gate must actually run green on the contained runner before its PR is complete; current full baseline remains unknown.
