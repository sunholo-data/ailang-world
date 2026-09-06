# Fleet Residual Reporting: State the Enumeration Boundary

**Status:** Implemented — PR124 merged as `bf15c73f143651a4a32aeec956a878efb1ae328c`; merge CI GREEN3/3.
**Target:** World iteration 160, queue row 64 (clause-2)
**Priority:** P2 · **Estimated:** ~0.2d · **Class:** HARNESS
**Dependencies:** None · **Planner-Lane:** codex-ok
**Measured base:** `d81ac424460e181c2a36a2ed11894415edca8048` · **Date:** 2026-09-06

## Problem and Goals

The fleet comparison's phase-3 residual enumeration uses `tools/launchd` and the
literal file `scripts/mission_decisions.sh` (V2). Its success disclaimer currently says
“untracked fleet additions not certified,” without distinguishing enumerated residuals
from paths it never enumerated (V3). Operators cannot interpret a missing warning as
absence of an out-of-boundary addition.

Reproduced with the live base verifier and synthetic committed trees (V7): three matching
driver files give rc=0; adding fleet `scripts/mission_decisions_v2.sh` still gives rc=0,
zero filename mentions and the same success wording (apart from HEAD). Adding
`tools/launchd/new-helper.sh` gives rc=0, exactly one per-path warning and a summary of one.
The last fixture contains both additions. These are synthetic observations; this doc
makes no claim about today's real fleet inventory.

**Goal:** retain the authorized scan boundary while making its reporting precise. Every
successful comparison must state the tracked comparison scope, phase-3 enumeration
scope, and explicit required paths. Outside-boundary files are **unenumerated, not zero**.
No ownership, language, core, or scan expansion is proposed.

## Decisions and Scope

| Decision | Reason | Chosen by | Deadline | Change cost |
|---|---|---|---|---|
| Use row 64's narrow-claim option | Correct the measured ambiguity without defining a larger driver inventory | User-authorized; agent specifies wording | Design | Low |
| Keep both pathspecs and required-path mechanism | Scope disclosure must not change ownership or certification | Agent within authorization | Design | Low |
| Supersede row 54’s success/summary text pins for this follow-up | Row 54 is LANDED; current row 64 explicitly authorizes narrowing the success claim. Preserve its comparison/refusal/skip contracts; use AC1–AC3 here for replacement wording (V12–V14) | User-authorized row 64 | Design | Low |
| Add a durable synthetic Go regression | Exercise the real script with committed additions in isolated repos | Agent | Implementation | Low |

**Design freeze:** choices above resolved; no human decision pending.
**Deferred decisions:** implementer may choose helper names and assertion organization.
Broader residual discovery or a shared fleet inventory is a systemic alternative deferred
to separately authorized work; it is not a dependency or a new obligation for this row.

**Why this is not a package:** this is reporting and regression coverage for World's
repository verification harness, not a World semantic behavior or extension. A package
would add an execution/distribution boundary to a small shell diagnostic change. The
planned host addition is test-only. This respects S3 and DESIGN §14 (V9).

## Solution Design and Examples

Change only reporting and adjacent explanatory comments in `check_driver_fleet()`.
Keep the `git ls-tree` arguments, phase-1 accounting, phase-2 explicit required checks,
phase-3 eligibility/counting and return codes unchanged (V2–V3).

Output is human-readable diagnostic prose, not a structured output format or a stable
machine-parser contract. The following clauses are regression expectations for this
change (dynamic fields shown in braces):

```text
✓ fleet-comparison arm: {N} files match fleet HEAD {HEAD}; checked set: World-tracked paths under tools/launchd and exact file scripts/mission_decisions.sh, plus explicit REQUIRED_FLEET_PATHS
  explicit required paths (checked separately): tools/launchd/lib/pin-root.sh
  phase-3 residual enumeration only: tools/launchd and exact file scripts/mission_decisions.sh; files outside this boundary are unenumerated (not zero), not certified by phase 3
```

Render the required-path line from `REQUIRED_FLEET_PATHS`, not a second hard-coded list;
print `(none)` if empty. Currently its member is `tools/launchd/lib/pin-root.sh` (V2).
Required checks remain explicit even if a future member lies outside the phase-3
boundary: “unenumerated by phase 3” does not erase phase-2 checking.

For an in-boundary addition, preserve the existing per-path warning and replace the
summary wording with:

```text
⚠ {N} unclassified fleet-only paths within the phase-3 boundary not certified (see above)
```

For the sibling addition, do not invent an enumeration or a count for that path. Emit
the scope statements even when the in-boundary residual count is zero. A zero/missing
in-boundary warning count must never be presented as a fleet-wide zero. Update the
phase-3 comment to name the boundary too. The shared function serves isolated and main
invocations (V4), so keep disclosure there, not in only one caller.

### Files to Modify/Create (implementation scope)

- `scripts/verify_go.sh` (~15 lines changed): bounded success/count wording, required-path rendering, phase-3 comment.
- `host/verifygate/driver_fleet_scope_gate_test.go` (new, ~140–180 LOC): `TestDriverFleetResidualScope` and local synthetic fixture helpers; reuse `copyGateFile` (V5).
- `design_docs/planned/w-fleet-residual-net-shares-phase-1-pathspec.md` (this document): record implementation evidence when available.

No edits to frozen `tools/launchd`, real fleet, shared skills, charter/log, or V1.
The test belongs in `host/verifygate`: its existing evidence gate executes an isolated
copy of the live script; `copyGateFile` copies from `repoRoot` into a supplied root (V5).
A scoped search found no `--driver-fleet-check` test there, with same-scope evidence
fixture hits as the positive control (V6). This is not a claim about all repository tests.

## Single Milestone M1 — Honest Boundary and Durable Regression (~0.2d)

One small implementation milestone: wording/comments, committed regression source,
focused execution and mutation evidence. “Committed regression” means the new test ships
with the later implementation; this designer neither changes code nor commits.

### Acceptance Criteria and Testing Strategy

- [x] **AC1 — Baseline disclosure.** New named test runs without skips. Copy the live
  `scripts/verify_go.sh` using `copyGateFile`; make synthetic World/fleet Git repos under
  `t.TempDir()`. Commit identical minimal `tools/launchd/control.sh`,
  `tools/launchd/lib/pin-root.sh`, and `scripts/mission_decisions.sh` in both. Invoke
  `bash scripts/verify_go.sh --driver-fleet-check` with World as cwd, explicit synthetic
  `AILANG_FLEET_REPO`, and CI/AILANG_BIN unset. Require rc=0, three compared files, all
  three disclosure statements, and the required member printed exactly once on its line.
- [x] **AC2 — Addition outside.** Commit only `scripts/mission_decisions_v2.sh` to the
  synthetic fleet. Assert `git ls-tree` sees it at HEAD (fixture health), then run the
  same command. Require rc=0, three matches, zero per-path residual warnings, and the
  complete statement that phase 3 only enumerates the two stated paths and that outside
  files are unenumerated (not zero). Assert the sibling is not individually reported.
  This success/disclosure assertion is the primary regression killer.
- [x] **AC3 — Same-scope positive control.** Keep the sibling and commit
  `tools/launchd/new-helper.sh` to fleet. Require rc=0, three matches, exactly one existing
  per-path warning naming the helper, summary count one qualified by “within the phase-3
  boundary,” and the same disclosure. Check both additions exist at HEAD. Never infer
  global zero from the silent sibling. AC2 and AC3 must run in the same test invocation.
- [x] **AC4 — Required checks remain clear.** In a separate synthetic case, omit the
  required pin-root file from World's committed tree while keeping it in fleet. Require
  rc=1, `REQUIRED fleet paths MISSING LOCALLY`, and its exact path. Reject success, skip,
  and `AILANG_BIN is unset` output. This is a compatibility guard, not the row-64 killer.
- [x] **AC5 — Production wording mutation kills the claim.** After green, mutate only
  the copied production reporting block back to the base success/summary strings and
  remove its new scope/required disclosures. Do not alter comparison logic, fixtures,
  CLI dispatch, assertions, or expected strings. Run the identical AC2 oracle against
  that copy: command rc must remain **0**, three matches must remain present, and the
  assertion must fail specifically for missing phase-3 boundary/unenumerated wording.
  Require the mutation to match exactly one reporting block and record before/after
  hashes. Prefer a committed mutation subtest that expects this oracle failure; also
  record the ordinary named test going red with this mutation applied to the production
  script in a disposable implementation checkout. Restore captured pre-drill bytes,
  verify their hash, rerun green; do not restore from HEAD over uncommitted work.
- [x] **AC6 — All tests passing (scoped); documentation updated.** Run
  `bash -n scripts/verify_go.sh`,
  `go test ./host/verifygate -run '^$'` (test compilation fence), and
  `go test ./host/verifygate -run '^TestDriverFleetResidualScope$' -count=1 -timeout=60s -v`.
  Record named RUN/PASS, no skips, and the AC5 claim failure. Update this document with
  results. Broader full tests remain controller gates outside the sandbox; no inherited
  full-suite green is evidence here. Socket denials are **UNINFORMATIVE UNDER SANDBOX**.

All source-writing fixtures and mutations must stay inside synthetic temporary repos.
Use local synthetic Git identity, disable signing/hooks, clear inherited Git repository
redirect variables, and bound subprocesses (e.g. 10 seconds each, 60 seconds for the test).
Never fall back to real HOME/fleet. Assertion expectations must be independent literals,
not parsed from production wording. The test must copy the live working-tree script,
so changing its production text actually changes the tested bytes.

### Conflict surface: measured consumers and explicit disposition

V12 independently reproduced **92 matching lines in 8 tracked files**, byte-for-byte
identical to `/tmp/world-i160-consumers.txt`. Counts are matching lines, not phrase
occurrences. The inventory below classifies **every hit**; line numbers refer to the
pre-revision tracked source, not this newly authored, untracked design. All paths in
this table are repository-relative. Grouped line lists have one common disposition.

| File | Hits and classification (all matching lines) | Disposition |
|---|---|---|
| `scripts/verify_go.sh` | **5**: 134 accounting comment; 236 missing-required diagnostic producer; 241 success producer; 243 summary producer; 267 CLI dispatch | Revise 241/243 and adjacent reporting comments (including 134’s quoted success example); retain 236 verbatim and 267’s behavior. None of these five is a text parser. Function callers at 269/373 inspect return status, not output (V4/V13). |
| `design_docs/planned/sprint_w-driver-drift-gate-compares-the-copy-to-itself.json` | **27**: 4,27,37,54,86,111,191 execution/design instructions and historical baseline; 118,128,152 historical green/count/skip findings; 207 required error pin; 211,213 success/summary pins; 227,228,229,230,231,232,233,235,236,237,238,239 historical AC1–7/9–13; 266,274 residual limitations | Historical structured sprint **document**, not a structured verifier-output parser. Success/summary pins and AC9/10/13 exact text are superseded below; AC2’s old negative needs the new positive counterpart when reused. Required error and other behavior contracts remain. Preserve artifact as history; no JSON edit. |
| `design_docs/planned/w-driver-drift-gate-compares-the-copy-to-itself-sprint-plan.md` | **25**: 33,42,60,66,122,235,252 implementation/dispatch and historical baseline; 137,145,162 historical green/count/skip findings; 222 required error pin; 226,228 success/summary pins; 266,267,268,269,270,271,272,274,275,276,277,278 historical AC1–7/9–13 | Concrete exact-text conflict: §3 says the drill is keyed to exact strings. Supersede success/summary pins and AC9/10/13 text expectations; refresh AC2’s negative if replaying the drill. Other behavior remains required. Historical file is retained unchanged. |
| `design_docs/planned/w-driver-drift-gate-compares-the-copy-to-itself.md` | **30**: 117,220,367 residual semantics; 127,296,375,385,389,391,396,408,420,424,463,464 implementation/CLI/AC instructions; 264,365 required diagnostic; 269,271,335,338,345,346,366,453 success/summary code examples and claims; 485,487,489,491 prior quorum history; 505 accounting explanation | Historical design/code snippets, not executing callers. Supersede success/summary wording and its “adopted verbatim” rule for row 64 only. Keep refusal, accounting, skip and residual behavior. Retain historical measurements and review transcript unchanged. |
| `design_docs/planned/w-canary-fence-blind-to-a-skipped-canary.md` | **1**: 539 historical V20 isolated invocation and drift failure observation | Error/rc evidence; no success/summary dependency. Preserve unchanged. |
| `design_docs/world-mission-log.md` | **2**: 16608 landed implementation narrative; 16666 historical wrong-baseline finding | Historical invocation evidence, no running parser. Preserve unchanged. |
| `design_docs/world-mission-status-archive.md` | **1**: 9 iteration-148 landed status, exact adopted wording and historical drills | Historical record, not a current exact-text acceptance gate. Preserve unchanged. |
| `design_docs/world-mission.md` | **1**: 4843 row 54 LANDED status and CLI description | Authority/status evidence. Preserve unchanged; current row 64 at 5236–5256 authorizes the follow-up. |

**Row-54 acceptance supersession.** The JSON still says `planned - awaiting
sprint-executor` and all three artifacts remain under `planned/`; those are stale
artifact labels, not evidence of an open row-54 obligation. The authoritative row 54
at `world-mission.md:4843` records LANDED on 2026-09-02, iter-148, merge `14036ee`,
corroborated by the log and archive (V14). Row 64 explicitly offers “narrow the success
line so it claims only what the pathspec can see” and requires an addition-shaped test.
This design selects that already-authorized option. It deliberately supersedes the
row-54 plan §3 / JSON success and summary text pins, and the exact-string parts of
historical AC9, AC10 and AC13; they would fail if replayed unchanged after this change.
Their healthy comparison, non-fatal addition and count-one obligations survive as
AC1–AC3 here. Historical AC2’s absence of `tracked copy is current` would become vacuous:
any future reuse must reject the **new** success/checked-set disclosure on a skip, with
AC1 as its positive counterpart. Refusal strings (including required missing), skip
text and rc mapping, path-liveness/accounting, per-path warnings and exact pathspecs
remain unchanged. This is prospective, narrowly scoped supersession recorded here;
it does not rewrite the old sprint’s results or reopen its completed implementation.

**Executable consumers.** The whole-tracked-repo phrase/mode search found the five
production hits above and 87 documentation hits; it found no executable exact-success
or exact-summary assertion/parser. The supplemental caller check (V13) identifies
`.github/workflows/ci.yml:166` running the full script directly, with no pipe or output
parsing; it receives exit status and human logs. Both function call sites retain their
return-code handling. Existing evidence-mode test helpers invoke a different mode;
toolchain source checks inspect their own floor/comparator region, not these lines.
These observations are scoped to the searched tracked repository and inspected callers.
They do **not** prove universal absence of external, untracked, generated, or dynamically
constructed consumers. An external exact-text consumer could break. Row 64 intentionally
authorizes this wording interface change; preserving the old ambiguous line or adding
a versioned structured-output mode is not the selected remedy. No structured parser
compatibility is claimed. No parser/typechecker/codegen positions are touched.

**Risk:** prose refactoring causes a false red; mitigate by pinning semantic clauses and
exact boundary paths without depending on indentation or dynamic HEAD. Future pathspec
changes must revise the disclosure and test together.

## Verification Log

Executed by this designer at the measured base, with `export PATH=/opt/homebrew/bin:$PATH`.
Commands below are read-only in the real repository; synthetic fixture writes and Git
commits occur only under temporary directories, cleaned automatically.
No implementation or post-fix mutation has been executed. Earlier documents' historical
measurements are not adopted as current facts.

| ID | Command / scope | Observed output and conclusion |
|---|---|---|
| V1 | `git rev-parse HEAD`; `git status --short` before writing | `d81ac424460e181c2a36a2ed11894415edca8048`; empty status. |
| V2 | `sed -n '111,219p' scripts/verify_go.sh` (read within `sed -n '1,270p'`) | Required array contains `tools/launchd/lib/pin-root.sh`; phase 1 and phase 3 use `HEAD -- tools/launchd scripts/mission_decisions.sh`; phase 2 loops explicit required paths; phase 3 excludes tracked and required paths before increment/warning. |
| V3 | `sed -n '220,287p' scripts/verify_go.sh` (read in two adjacent ranges) | Zero comparable files, differing/missing files return 1; success at 241 has `tracked copy is current (untracked fleet additions not certified)`; residual summary at 243 is unqualified; isolated mode precedes AILANG_BIN refusal and maps CI skip 2 to 0. |
| V4 | `git grep -n 'check_driver_fleet' -- scripts/verify_go.sh` | Definition 125, isolated call 269, main call 373. |
| V5 | `sed -n '1,100p' host/verifygate/evidence_manifest_gate_test.go`; `sed -n '195,240p'` same file; `sed -n '1,70p' host/verifygate/module_manifest_gate_test.go` | Evidence helper uses `t.TempDir`, default live script copy; named test's extra/missing/duplicate/empty cases assert outputs and rc. `copyGateFile` opens `repoRoot/rel`, creates destination exclusively, copies bytes. |
| V6 | `git grep -n -e 'driver-fleet-check' -e 'newIsolatedEvidenceGateRoot' -- 'host/verifygate/*_test.go'`; `git ls-files 'host/verifygate/*fleet*' 'host/verifygate/*evidence*'` | No driver-mode hit; positive control helper hits at 56, 204, 213, 221, 230, 239, 252. Inventory lists only `host/verifygate/evidence_manifest_gate_test.go`; same-scope positive control supports choosing a new fleet test file. |
| V7 | Python reproduction below; `cat /tmp/world-i160-repro.json`; `cat /tmp/world-i160-repro-root` | Re-executed all three states: baseline rc=0/3 matches/0 warnings; outside rc=0/3 matches/0 sibling mentions/0 warnings; inside rc=0/3 matches/0 sibling mentions/1 helper warning/summary 1. Trees confirm additions at HEAD. Existing JSON agrees in content; stdout/stderr ordering is not asserted. |
| V8 | `go test ./host/verifygate -run '^TestEvidenceNamedManifestRejectsUnpinnedTest$' -count=1 -timeout=60s` | rc=0, `ok github.com/sunholo-data/ailang-world/host/verifygate 1.057s`. Focused existing harness baseline only, not the proposed test or full suite. |
| V9 | `cat design_docs/coding-standards.md`; `sed -n '610,650p' design_docs/DESIGN.md`; `sed -n '5230,5275p' design_docs/world-mission.md` | S3 package rationale, S6 mutation discipline; §14 excludes shared mission machinery from self-modification; row 64 explicitly permits narrow success wording plus addition test, ~0.2d. |
| V10 | `git ls-files 'design_docs/planned/*fleet*' 'design_docs/implemented/*fleet*' 'design_docs/planned/*driver*' 'design_docs/implemented/*driver*'`; read first 85 lines of related design and lines 195–250 of its sprint plan | Three row-54 artifacts under planned (design, sprint plan, JSON). Plan pins old success/residual strings. Related scope is source comparison; this doc changes its reporting boundary. No neural search used. |
| V11 | `git log -4 --oneline -- scripts/verify_go.sh` | `8b600fd` row 61, `b5c303e` iter 151 record, `14036ee` row 54 driver comparison, `c13ad1f` row 48. Together with phase-1 accounting comments (V2), supports treating this as the recurring claim-versus-enumerator pattern, without expanding this item. |
| V12 | `git grep -n -e 'tracked copy is current' -e 'untracked fleet additions' -e 'unclassified fleet-only paths' -e 'driver-fleet-check' -e 'REQUIRED fleet paths MISSING LOCALLY' -- .`; Python compares complete stdout to `/tmp/world-i160-consumers.txt` and groups filename/line | Exact equality **True**, **92 hits / 8 files**, production script **5**, remaining **87** documentation. Every line classified above. Broader `untracked fleet additions` includes the reviewers’ full disclaimer query. This new untracked design is outside that tracked inventory. |
| V13 | `git grep -n -e 'verify_go.sh' -e 'fleet-comparison arm' -- .github scripts host Makefile`; read `.github/workflows/ci.yml:155–170`, `scripts/verify_go.sh:267–276,365–385`, `host/verifygate/toolchain_pin_gate_test.go:1345–1380`; evidence test helpers V5 | CI runs `./scripts/verify_go.sh` directly; isolated/main function callers handle rc, not prose. Evidence helpers use `--evidence-manifest-check`; toolchain assertions are region-scoped to floor/comparator. No executable consumer of replaced literals was found in V12; no universal absence claim. |
| V14 | Read row-54 sprint plan §3 and AC table, JSON header and matched fields; `world-mission.md:4843,5236–5256`; `world-mission-log.md:16600–16618,16660–16672`; archive line 9 | Exact pins and stale JSON planned label are real. Authoritative row 54 is LANDED, merge `14036ee`; current row 64 authorizes narrowing the claim. Success/summary pins prospectively superseded with precise AC mapping above; historical files untouched. |
| V15 | Live two-arm missing-required fixture described below, 10-second timeout per subprocess; read `scripts/verify_go.sh:180–198,219–245` | Byte-identical live script SHA256 `4614363ed51839a8b8a8a190429838b6d15315c91792ac585f08f3ab781b444a`. Healthy rc=0, **3** matches. Removing only synthetic World’s committed required member gives rc=1 and exact `REQUIRED fleet paths MISSING LOCALLY:` plus `tools/launchd/lib/pin-root.sh (REQUIRED by World, absent locally)`. All assertions passed; correct phase-2 failure, not skip/binary gate/zero-comparable/difference/missing-in-fleet/accounting failure. |
| V16 | Pre-change `go test ./host/verifygate -run '^TestDriverFleetResidualScope$' -count=1 -timeout=60s -v` after adding the regression but before changing production | rc=1. Synthetic command rc remained 0 with three old-style matches; `baseline_disclosure`, `outside_boundary_unenumerated`, `inside_boundary_positive_control`, and `empty_required_list` red because the new disclosures/qualified summary were absent. `required_path_missing` passed, preserving attribution. |
| V17 | Post-change `bash -n scripts/verify_go.sh`; compile fence; focused named test | All rc=0. The compile fence compiled `_test.go` and reported `[no tests to run]` as intended. The focused run printed the top-level RUN/PASS plus RUN/PASS for `outside_boundary_unenumerated`, `baseline_disclosure`, `inside_boundary_positive_control`, `required_path_missing`, and `empty_required_list`; no skips or no-test warning. |
| V18 | Durable self-contained fixture controls | Fresh committed World/fleet repos per case, explicit synthetic fleet and HOME, inherited `GIT_*`/CI/AILANG_BIN/HOME cleared, system/global Git config and hooks/signing disabled, and each subprocess bounded to 10 seconds. Baseline/outside/inside rc=0 with three matches; missing-required rc=1 with the exact refusal; empty array rc=0 with one `(none)` line and no nounset failure on Bash 3.2. |
| V19 | Disposable production-only wording reversion under `/private/tmp`; exact-one new reporting block replaced with the base success/unqualified summary and disclosures removed | Pre-drill hashes: verifier `2e80cca964a6b56a9f2b4fca06d59d707c236c6fe6f9865ecf85c70ec84cb33d`, test `e6e5230a39649be412dfe87d4e1d6305f7a517b5c3fa9b9f2ac1372917bd31b5`. Mutant verifier hash changed to `3f4972ab30b7e825b40d0fa182d12cbe6c48ad1771554350ee8d35f1940c2647`; old success and old summary each occurred once. Compile fence rc=0. Ordinary named test rc=1: its first failing arm was `outside_boundary_unenumerated`, specifically `missing phase-3 boundary/unenumerated disclosure`; captured output independently showed the copied script rc=0 and three old-style matches. Required-path control still passed. |
| V20 | Restore from captured pre-drill bytes, hash comparison, and restored focused rerun | Restored verifier/test hashes exactly equalled `2e80…b33d` / `e6e5…31b5`; the four new reporting literals each occurred exactly once. `bash -n`, compile fence, and focused named test returned rc=0 with all named RUN/PASS lines and no skips. No worktree production source was mutated by the drill. |

### Missing-required live fixture (revision evidence)

Executed with Python `tempfile.TemporaryDirectory(prefix='world-i160-required-',
dir='/tmp')`, making fresh `world`, `fleet`, and `home` directories. Environment removes
all inherited `GIT_*`, `CI`, and `AILANG_BIN`; sets HOME to the synthetic home,
`GIT_CONFIG_NOSYSTEM=1`, `GIT_CONFIG_GLOBAL=/dev/null`, and explicit synthetic
`AILANG_FLEET_REPO`. Every subprocess has `timeout=10`. Git uses
`-c core.hooksPath=/dev/null -c commit.gpgsign=false` and local synthetic name/email.
Both repos commit identical `#!/bin/sh\n: synthetic\n` files at
`tools/launchd/control.sh`, `tools/launchd/lib/pin-root.sh`, and
`scripts/mission_decisions.sh`. Copy the live verifier into synthetic World and assert
byte equality before running `bash scripts/verify_go.sh --driver-fleet-check` there.
Healthy assertions: rc=0, `3 tracked frozen-core files match fleet HEAD`, and
`tracked copy is current`. Measured fleet HEAD was
`b18ed9ce934006eee313f78e80bfd6de26af1e69` (fixture-generated, not a reusable dependency).
Then `git rm -- tools/launchd/lib/pin-root.sh` and commit **only in synthetic World**;
assert `git ls-tree -r --name-only HEAD` omits the path in World and retains it in fleet.
Re-run the identical invocation and require rc=1 plus these measured literal lines:

```text
verify_go.sh: FATAL: DRIVER DRIFT vs FLEET (D-WORLD-DRIVER-1) — REQUIRED fleet paths MISSING LOCALLY:

  tools/launchd/lib/pin-root.sh (REQUIRED by World, absent locally)
  The driver is fleet-owned; land the required file as a fleet-authored commit. World's controller must not edit or absorb it.
```

Assert absence of `tracked copy is current`, `SKIPPED`, `AILANG_BIN is unset`,
`0 comparable driver files`, `MISSING IN FLEET`, `committed copy differs`, and
`PHASE-1 ACCOUNTING BROKEN`. The two remaining equal non-required paths keep the
comparison live, and the fleet still contains the required file: the rc=1 is attributable
to the intended missing-local-required branch. TemporaryDirectory cleaned all fixture
writes on exit. This proves AC4’s base diagnostic; it is not implementation or mutation
proof for the proposed reporting change.

### Reproduction command (actually executed)

The controller's pointer file resolves to synthetic repos. Copies below avoid modifying
those supplied artifacts. No new commits were made: reset selects their existing synthetic
committed baseline and additions. The future durable test must create its own fixtures,
not depend on this pointer or these hashes.

```bash
export PATH=/opt/homebrew/bin:$PATH
python3 - <<'PY'
import os, pathlib, shutil, subprocess, tempfile, json
source=pathlib.Path('/tmp/world-i160-repro-root').read_text().strip()
with tempfile.TemporaryDirectory(prefix='world-i160-designer-') as temp:
    root=pathlib.Path(temp)
    for repo in ('world','fleet'):
        shutil.copytree(pathlib.Path(source)/repo, root/repo)
    shutil.copyfile('scripts/verify_go.sh', root/'world/scripts/verify_go.sh')
    env=os.environ.copy()
    env.pop('CI',None); env.pop('AILANG_BIN',None)
    env['AILANG_FLEET_REPO']=str(root/'fleet')
    for label,rev in [('baseline','9703e22'),('outside','4242d1e'),('inside','debcc52')]:
        subprocess.run(['git','-C',str(root/'fleet'),'reset','--hard',rev],check=True,capture_output=True)
        p=subprocess.run(['bash','scripts/verify_go.sh','--driver-fleet-check'],cwd=root/'world',env=env,text=True,stdout=subprocess.PIPE,stderr=subprocess.STDOUT,timeout=30)
        print(json.dumps({'label':label,'rc':p.returncode,'outside_mentions':p.stdout.count('scripts/mission_decisions_v2.sh'),'inside_warning_count':p.stdout.count('unclassified fleet-only path (not tracked, not required): tools/launchd/new-helper.sh'),'output':p.stdout}))
        tree=subprocess.check_output(['git','-C',str(root/'fleet'),'ls-tree','-r','--name-only','HEAD'],text=True)
        print(label+' fleet tree: '+tree.replace('\n',', '))
PY
```

## Axiom Compliance

Scored against the skill's [design structure reference](../../../ailang/.claude/skills/design-doc-creator/resources/design_doc_structure.md).

| Axiom | Score | Reason |
|---|---:|---|
| A1 Determinism | 0 | Comparison algorithm unchanged. |
| A2 Replayability | +1 | Durable committed synthetic additions reproduce the reporting condition. |
| A3 Effect Legibility | 0 | No new production effects. |
| A4 Explicit Authority | 0 | Existing ownership and required checks retained. |
| A5 Bounded Verification | +1 | Isolated named regression and wording mutant, bounded subprocesses. |
| A6 Safe Concurrency | 0 | No concurrency behavior change. |
| A7 Machines First | +1 | Explicit enumeration domain and executable regression oracle prevent treating silence as zero; output remains human-readable prose. |
| A8 Minimal Syntax | 0 | No language syntax change. |
| A9 Cost Visibility | 0 | No cost semantics change. |
| A10 Composability | 0 | CLI and rc contracts retained; human-readable success/summary wording intentionally changes under row 64. |
| A11 Structured Failure | 0 | Refusal behavior unchanged. |
| A12 System Boundary | +1 | Reporting explicitly distinguishes measured and unenumerated paths. |

**Net +4. Hard checks:** A1 deterministic; A3 effects unchanged; A4 authority unchanged;
A7 scope is explicit and regression-tested by the planned oracle; no structured-output
claim. No hard violations.

## Related Documents and Review Handoff

- [Row-54 design](../planned/w-driver-drift-gate-compares-the-copy-to-itself.md) and
  [sprint plan](../planned/w-driver-drift-gate-compares-the-copy-to-itself-sprint-plan.md): textual
  coverage check read these; retain their comparison mechanism, revise the ambiguous
  residual claim. This is row 64's separately authorized follow-up, not a duplicate
  source-comparison feature. Their historical measurements/status are not current evidence.
- [Coding standards](../coding-standards.md): S3, S6 and usage disclosure guide this change.

Unattended quorum trigger applies. Per the mission instruction, the controller will run
Sol/Gemini/GLM; this designer does not run quorum or self-review as Astra. No unresolved
scope judgment remains. Implementation/test/mutation proof and controller review are pending.

## Quorum Round 1 — revision disposition

Source: `/tmp/world-i160-quorum-output.json`, timestamp `2026-09-06T05:52:19Z`.
**Present 3/3; blocked/reject 3; absent 0; cost USD 0.05905092.** Synthesis: blocked;
controller in-session verdict: pass. This is the ONE permitted designer revision.

| Reviewer | Round-1 objection | Resolution in this revision |
|---|---|---|
| Sol (`gpt5-6-sol`) | Unmeasured output consumer conflict; human versus machine output unclear | V12 independently reproduces all 92 hits and classifies each; V13 checks actual callers. Human-readable prose stated explicitly. External/untracked consumers not ruled out; row 64 authorizes the intentional wording change, so the suggested preservation/versioning alternative is not adopted. |
| Gemini (`gemini-3-1-pro`) | AC4’s exact missing-required literal unverified | V15 measures healthy rc=0 then intended missing-required rc=1 using the byte-identical live script; exact literal/path and rejection of wrong failure causes asserted. |
| GLM (`oc-glm-5-2`) | Concrete consumer inventory absent; row-54 pinned strings unresolved | V12 inventory plus V14 authoritative LANDED/current-row evidence; prospective supersession of §3/JSON and AC9/10/13 wording, with AC2 non-vacuity disposition, is explicit. Historical artifacts remain unchanged. |

All three objections are addressed in the artifact; this does not assert reviewer
acceptance. The controller will perform the single re-quorum. No production, tests,
charter/log, real Git metadata, fleet, messages, or shared skills were changed by this
revision. Implementation and post-fix tests/mutation evidence remain pending.

## Quorum Round 2 — controller record

2026-09-06T05:58:06Z: **PROCEED, 3/3 present and pass; absent_reviewers=[]**.
Sol substitutes for author Astra in the OpenAI seat. Round cost $0.09106965;
iteration quorum total $0.15012057. No carve-out used. Implementation handoff:
make fixtures self-contained, safely render an empty required array, and record
focused tests plus production-wording mutant/red/restoration. External-consumer
checks remain explicitly scoped, never a universal absence claim.
