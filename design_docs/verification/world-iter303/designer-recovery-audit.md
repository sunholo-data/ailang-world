# World iteration 303 — inherited M7b recovery design audit

**Verdict: PASS — no design revision required.**

The inherited M7b can proceed through its existing verification and landing process without changing the capture design, acceptance bar, authority boundaries, or milestone scope. This is a design admissibility judgment, **not a replacement independent evaluation, fresh test verdict, or authorization to land with a red mandatory gate**. Historical reds below remain findings. No tests were run in this audit.

Assignment: unattended designer audit; requested explicit lane `codex exec --model gpt-6.1-sol`; no subagents. Candidate identity supplied by the caller: PR361, `sprint/world-iter302-row150-m7b`. Fresh local inspection found HEAD `f615a9698370694664fbbc01431e4a2c549377cc` on the recovery branch `mission/world-iter303-recover-m7b` [E1]. PR state and live CI were not queried. Model identity is not inferred from the inherited executor's model; this report does not provide independent runtime model attestation.

## Reasons and boundaries

1. **The existing authority and design remain sufficient.** D-WORLD-69 is RESOLVED A: retain the complete showcase and order 26 → 105 → 172 → 150, through normal design/quorum/plan/judge gates; 151 follows 150. Its trailing old prerequisite annotation is historical text, not a new scope ruling. The capture design r4 explicitly retains exact two-pass raw-byte equality and treats fresh-seed equality as a residual owned by Mark [E3–E5]. Recovery does not reopen that ruling or narrow the showcase.
2. **M7b implements the existing verifier contract.** Fresh source readings show the six specified names/dimensions, independent PNG structure and directory-limit checks, decoded color/chromatic/rectangle/scheme checks, and raw SHA256 equality across two validated passes. The CLI loads the supplied rectangles and public demo manifest and emits a derived capture manifest [E7–E9]. These are tooling checks. They do not establish that real rendered text occupies a measured rectangle. The design already requires raw HTML checks and measured geometry separately [E4]. The source's provisional numeric content thresholds do not constitute completed product acceptance.
3. **The inherited implementation has traceable independent review.** Evaluator r1 records PASS 95 on `fa04df00`; r2 records PASS 95 on `02e12be2`, carrying r1 forward and checking the Go-surface fix. Fresh git comparison shows that `02e12be2..HEAD` changes only the two sprint-plan artifacts. Fresh byte comparisons match all four marketingcheck source/test files to the retained executor snapshot [E1, E6, E11–E12]. Historical judgments are preserved as historical, without upgrading them into tests run by this designer.
4. **The extra census entries are an integration correction within the existing design.** The controller fix adds exactly `scripts/marketingcheck/main.go` and `scripts/marketingcheck/png.go` to the explicit production-Go census allow-list. Reading the whole guard shows that other unlisted production Go paths still trigger an error and missing listed paths also trigger an error. The inherited supervisor test changes only its package declaration. R2 records both a rogue-file red and a removed-entry red, with restored green [E6, E10–E11]. Neither change requires a different capture architecture.
5. **Recovery must preserve honest gate semantics.** All three design quorum artifacts remain BLOCKED; r4 is the recorded bounded-refinement disposition, not a retroactive quorum PASS. The design and plan already distinguish tooling verification, future real capture acceptance, and attended publish/mint. No new design or quorum cycle is proposed by this audit [E4–E5, E13].

Only this report is the authorized output. No production change, design/charter edit, new milestone, shared write, HOME seed, credential preparation, browser run, test, mutation drill, git write, message, GitHub write, commit, branch change, or push is part of this audit. Existing later capture/release obligations are cited solely to bound what M7b can claim; this report adds no work to them. It makes no AILANG language-support claim.

## Findings retained for recovery

- **Historical planner baseline RED:** the companion records `go build ./... && go test ./... -count=1 -timeout 600s` failing at its iteration-266 baseline, rc 1, 279.93 s. Its seven failures concern two controlling-terminal tests, two broker process/zombie-readiness tests, broker HOME confinement, capsule zombie readiness, and the PTY harness. It separately records the standalone `world/contracts.ail` test red with four vacuous/no-generator skips. These are inherited reports, not fresh measurements at the M7b parent or current HEAD. The standalone result is not a declaration that the aggregate AIL gate failed [E14]. Do not exclude tests, edit protected files, or seed real HOME to bypass them.
- **Historical iteration-302 broad Go reds:** `verify_go_iter302c.log` records `TestMissionConfigDefaultFlowIntegration/invalid_ledger` timing out and `TestMarketingSupervisorFakeDaemonCleanup/shell-hup` failing: signal killed, cleanup 28.896924125 s, and daemon-death assertion. The separate race log records an 8-minute daemon-package timeout while `TestWorkbenchPacketRosterTestedAndProven` was running, and another shell-hup failure with cleanup 29.684643292 s [E15]. These remain recorded failures; their cause is not freshly diagnosed here. The plan's isolated-green statements do not erase whole-suite reds [E5].
- **Historical positive evidence has a narrower scope:** the retained executor gate JSON records the three focused M7b tests, the package race command, and vet/format success; the AIL log ends in `verify gate PASSED` with 16 identities, 40 named tests and 496 se-tools tests. Two saved exact-head CI verdict files progress from `in_progress:none` to `completed:success`. The latter are local historical summaries, without a fresh API check or complete job inventory in this audit [E12, E16]. They do not certify current landing readiness.
- **Mutation limits remain explicit:** evaluator N1 records the neutered truncated-chunk guard producing a panic, which is not an assertion kill. The executor ledger additionally retains redundant-condition and filesystem-race survivors. R1's forbidden-chunk mutation reports `unknown-chunk` after removing the specific forbidden entry; that is discrimination of the typed refusal, not proof that the mutated parser accepted metadata. The executor ledger also records the stronger final mutation that bypasses both chunk maps for tEXt and reaches `got <nil>`. This audit does not claim universal mutation completeness [E11–E12].
- **Sandbox evidence rule remains binding:** socket bind denials are UNINFORMATIVE UNDER SANDBOX. A SKIP, empty test selection, panic, build failure, timeout, or unrelated earlier refusal is not acceptance or a mutation kill. No sandbox test result was generated here. The read-only process probe for runtime ancestry was itself denied (`Operation not permitted: 'ps'`); that supplies neither model attestation nor a test verdict [E2, E4–E5].

**Landing boundary:** the recovery controller retains responsibility for the normal independent-evaluator and unchanged mandatory-gate requirements at the exact candidate SHA, with the charter's isolation/load discipline. A fresh red is a finding to stop and report through the existing process, never permission to weaken a gate. The design-only PASS does not discharge that responsibility [E2–E3].

## Evidence register — exact commands and observed results

All paths below are relative to `/Users/voightkampff/dev/sunholo-data/.wt-world-iter303-recover-m7b` unless absolute. A read performed now is fresh verification of file contents only. Test/mutation/CI statements inside those files remain **UNVERIFIED historical execution evidence, inherited from the identified iteration**, because this audit did not rerun them. Summaries below describe observed command output; they do not imply omitted output was empty.

### E1 — fresh identity and commit lineage

- `git status --short --branch` → `## mission/world-iter303-recover-m7b`, with the initial working tree clean.
- `git rev-parse HEAD` → `f615a9698370694664fbbc01431e4a2c549377cc`.
- `git log -5 --format='%H %s'` → candidate plan completion `f615a969…`, controller census fix `02e12be22143bcad7f79ccf2b6e4cb76bab268ce`, verifier `fa04df0070c84b68a49b0e57b7e13330e36edd47`, M7a merge `c73c79048daf9c3a19ea05c96e793436eb45e9c5`, prior record `46af40ea…`.
- `git diff --stat c73c7904..HEAD` → seven files: two plan artifacts, census test, main.go, png.go, png_test.go, supervisor_test.go; 1,432 insertions, 22 deletions.
- `git diff --name-only 02e12be2..HEAD` → exactly `design_docs/planned/sprint_w-marketing-capture.json` and `design_docs/planned/w-marketing-capture-sprint-plan.md`. Positive output enumerates the post-review delta rather than relying on an empty command.

### E2 — instruction discovery, skill and process limits

`rg --files --hidden --no-ignore -g '!.git' -g 'AGENTS.md' -g 'CLAUDE.md' .` → `./CLAUDE.md`, `./website/static/AGENTS.md`. Positive controls are the real CLAUDE and AGENTS files on the same repository enumeration surface; the found static subtree guide does not govern a root report.

The ancestor check paired both filenames on the identical ancestor surface:

```sh
python3 - <<'PY'
from pathlib import Path
p=Path.cwd()
for d in reversed([p,*p.parents]):
 for name in ['AGENTS.md','CLAUDE.md']:
  f=d/name
  print(str(f), 'FILE' if f.is_file() else 'absent')
PY
```

Observed: AGENTS.md absent at each enumerated ancestor from `/` through this root; same-surface positive control is this root's `CLAUDE.md FILE`. This establishes the checked ancestor files, not arbitrary instructions elsewhere on the machine.

- `cat CLAUDE.md design_docs/coding-standards.md`, followed by `sed -n '1,100p' CLAUDE.md` and `sed -n '1,100p' design_docs/coding-standards.md` → binding charter/standards, nothing lands red, frozen machinery, package-first and non-vacuous evidence rules. The initial combined output was truncated; the relevant instruction bodies were read again separately.
- `cat /Users/voightkampff/dev/sunholo-data/ailang/.claude/skills/design-doc-creator/SKILL.md`, plus `sed -n '125,170p'` on that same absolute path → requested skill read; live language-claim verification, premise/test-body/negative-evidence discipline, and bounded quorum workflow. The assignment's report-only scope governs over creation/commit examples.
- `sed -n '27,48p' design_docs/DESIGN.md` and `sed -n '610,641p' design_docs/DESIGN.md` → defining thesis and normative self-modification boundaries read; no language capability assertion adopted into this report.

The attempted process probe was:

```sh
python3 - <<'PY'
import os, subprocess
pid=os.getpid()
for _ in range(6):
 r=subprocess.run(['ps','-p',str(pid),'-o','pid=,ppid=,comm='],capture_output=True,text=True)
 print(r.stdout.strip() or r.stderr.strip())
 if r.returncode or not r.stdout.strip():break
 pid=int(r.stdout.split()[1])
 if pid<=1:break
PY
```

Observed rc 1: `PermissionError: [Errno 1] Operation not permitted: 'ps'`. Runtime model identity was not established by it.

### E3 — charter body and D69

- `sed -n '1,98p' design_docs/world-mission.md` → gate profile and explicit prohibition on concurrent judge mutations/full-suite measurements; isolated baselines and load controls.
- `sed -n '1139,1310p' design_docs/world-mission.md`, `sed -n '1311,1489p' design_docs/world-mission.md`, `sed -n '1490,1672p' design_docs/world-mission.md` → complete Guardrails section read: positive enumeration, controlled negative evidence, narrow-refinement AC/mutations, historical-record verification, gate-list completeness, frozen kernel/shared machinery and valid mutation requirements.
- `rg -n '^\| D-WORLD-69 \|' design_docs/world-mission.md` → line 988, RESOLVED, attended A, complete-showcase/order/normal-gates ruling quoted in the reasons above. The same row still contains its old trailing prerequisite annotation; this audit does not rewrite the ledger.

### E4 — design body

`cat design_docs/planned/w-marketing-capture.md`, then `sed -n '134,211p' design_docs/planned/w-marketing-capture.md` and `sed -n '212,345p' design_docs/planned/w-marketing-capture.md` → seed/authority boundary, all six routes and dimensions, exact equality, structure/content thresholds, DOM plus measured-rectangle obligations, C11, fallback and review disposition read. The separate body reads cover the approach/review portions after larger batched output was truncated. V21 is explicitly pending product evidence; r4 is bounded refinement after blocked r3.

### E5 — binding M7b companion

- `sed -n '919,986p' design_docs/planned/w-marketing-capture-sprint-plan.md` → completed M7b, three named focused tests, package race gate, eight named mutations, controller census fix, PASS95 → PASS95 history, N1 and shell-hup carry.
- `sed -n '81,145p' design_docs/planned/w-marketing-capture-sprint-plan.md` → PD1 resolved by distributing required regions across the existing shots; attended publish/mint; real geometry and product capture separate from fake fixtures; C11 and invalid-kill rules.
- The exact JSON inspection was:

```sh
python3 - <<'PY'
import json
from pathlib import Path
p=Path('design_docs/planned/sprint_w-marketing-capture.json')
d=json.loads(p.read_text())
print('top-level keys:',list(d))
for k in ['gates_binding_on_every_milestone','verify_commands']:
 if k in d:print(k,d[k])
for m in d.get('milestones',[]):
 if m.get('id')=='M7b':
  print(json.dumps(m,indent=2))
PY
```

Observed M7b `status: completed`, M7a dependency, same three files/tests, four acceptance commands, eight production-guard mutation definitions, and completed executor/judge/controller-fix metadata. No test command in this output was executed by this audit.

### E6 — fresh source scope and controlled protected-path comparison

- `git diff c73c7904..HEAD -- scripts/marketingcheck/supervisor_test.go host/store/context_roots_test.go` → only supervisor `package marketingcheck` → `package main`, and four census insertions including exactly two verifier filenames.
- `git diff --name-only c73c7904..HEAD -- world go.mod go.sum tools/launchd design_docs/world-mission.md scripts/marketingcheck/png.go` → only `scripts/marketingcheck/png.go`. The known changed verifier is the positive control in the **same diff surface** for the protected-path non-change finding.

### E7 — fresh production body reads

- `cat scripts/marketingcheck/main.go` → CLI `check --pass-a --pass-b --manifest --rectangles [--out]`; calls LoadRectangles → CheckTwoPass → LoadDemoManifest → BuildCaptureManifest → VerifyCaptureManifest; exclusive creation of optional output.
- `sed -n '1,320p' scripts/marketingcheck/png.go` → exact six Shot entries; per-file 1,048,576, total 8,388,608; signature/chunk/CRC/IEND guards; rectangle schema, five required packet region names, minimum region size, color/chromatic/scheme/variance checks; provisional threshold comment.
- `sed -n '321,512p' scripts/marketingcheck/png.go` → directory inventory and regular-file requirement, all six presences, structure/content checks before hashing, both passes checked, per-name SHA256 comparison, public-manifest reference resolution and derived capture manifest. This establishes control flow by reading the mechanism, not by inferring it from a historical output.

### E8 — fresh structure/content test-body reads

`sed -n '1,230p' scripts/marketingcheck/png_test.go` and `sed -n '230,545p' scripts/marketingcheck/png_test.go` → synthetic fixture generation and exact-sized valid PNG construction; direct structural CRC refusal, positive/at-cap controls, cap+1 decoding/content checks, independent eight-under/ninth-over directory test, six content positives, uniform/mostly-gray/scheme/dimension/region guards, boundary values and packet-region requirements. These are read assertions, not fresh PASS results or real browser evidence.

### E9 — fresh inventory/CLI test-body read

`sed -n '546,804p' scripts/marketingcheck/png_test.go` → equal-copy manifest control, missing/extra/symlink refusals, two valid altered PNG variants checked before expected pass-mismatch, both hashes required in the error, manifest/reference refusal variants, CLI build/invocation assertions, and existing-output preservation.

### E10 — fresh census mechanism read

`sed -n '180,250p' host/store/context_roots_test.go` → entire TestProductionGoSurface body: explicit four-path allow-list, host/cmd recognition, unlisted production .go error, and missing allow-listed path error. These static observations support the scope of the controller correction; r2 supplies historical dynamic controls.

### E11 — independent evaluator artifacts, historical runs

Discovery: `rg --files --hidden --no-ignore --maxdepth 1 /Users/voightkampff/.ailang/state/world-iter302` → both evaluator reports. Body read: `cat /Users/voightkampff/.ailang/state/world-iter302/evaluator-r1.md /Users/voightkampff/.ailang/state/world-iter302/evaluator-r2.md`.

Observed r1: own `.eval-world-iter302` tree, fa04df00/base c73c7904, PASS95, named focused/race/vet-format results, eight required mutation rows, N1 panic survivor, and protected-path observation. Its forbidden row reports `unknown-chunk` rather than nil after removing tEXt from forbiddenChunks; derived mutation totals are not independently reconciled by this audit. Observed r2: own `.eval-world-iter302-r2` tree, 02e12be2, PASS95, exact census delta, guard positive/negative drills, restored hash, N1 unchanged. Preserve their verdicts and limits; this designer supplies no replacement score.

Transport body read: `cat /private/tmp/pi_run_iter302_eval.verdict.json /private/tmp/pi_run_iter302_eval_r2.verdict.json /private/tmp/ci_iter302_headf615a969.verdict /private/tmp/ci_iter302_headf615a969_r2.verdict` → both judge wrappers `verdict: ok`, rc0, minimax-m3; evaluated SHAs fa04df00 and 02e12be2 respectively. Wrapper success is evidence of completed transport, distinct from the evaluators' substantive PASS verdicts.

### E12 — retained executor artifacts and fresh byte comparisons

Discovery: `rg --files --hidden --no-ignore /Users/voightkampff/.ailang/state/world-iter302 /Users/voightkampff/.ailang/state/world-iter303 design_docs/verification/world-row150-design/quorum` enumerated the executor snapshot and reports; the iter303 path lookup returned an error, so it supplies **no global absence finding**. Broad temporary-directory discovery also encountered permission errors and is not used to establish any absence.

`cat /Users/voightkampff/.ailang/state/world-iter302/snap/M7b/ledger.md /Users/voightkampff/.ailang/state/world-iter302/snap/M7b/gates.json` → historical final gate table (focused rc0 at 0.437/0.552/1.270 s; race rc0 at 41.868 s; vet/format rc0), final eight mutation rows with restores, C11 equality, provisional geometry and residual survivors. This large output was truncated; findings here are confined to the visible initial gate/survivor and final mutation/C11 sections, not a claim to have inspected all 97 attempts. The stronger final tEXt mutation bypasses both forbiddenChunks and allowedChunks and records the specific assertion `got <nil>`.

Fresh comparison command:

```sh
python3 - <<'PY'
import hashlib
from pathlib import Path
snap=Path('/Users/voightkampff/.ailang/state/world-iter302/snap/M7b')
for name in ['main.go','png.go','png_test.go','supervisor_test.go']:
 p=Path('scripts/marketingcheck')/name
 q=snap/p
 print(p,hashlib.sha256(p.read_bytes()).hexdigest(),'snapshot_equal='+str(p.read_bytes()==q.read_bytes()))
a=Path('/Users/voightkampff/.ailang/state/world-iter302/snap/obs-exec-before.json')
b=a.with_name('obs-exec-after.json')
print('historical_observation_equal', a.read_bytes()==b.read_bytes())
PY
```

Observed all four `snapshot_equal=True`; hashes respectively:

```text
main.go            9b5cdd3a4cb25a657d684eac4067624b4ebdc5d67a4334342438cce896015c84
png.go             ea9b43df9b5e6f70757aff1798b28a258702271864a05f823ea88a760db5145e
png_test.go        cce1893446ff39ecb8366d4705f118fe143688539aa6e26602a1e41ac81e1a2c
supervisor_test.go 7cb4a1ba63ae13fb3fb0b09b9a1f25d466c1b45305da921aa3eb35c37bf9b447
historical_observation_equal True
```

### E13 — historical quorum artifact bodies

The exact focused body inspection was:

```sh
python3 - <<'PY'
import json
from pathlib import Path
for p in sorted(Path('design_docs/verification/world-row150-design/quorum').glob('*.json')):
 d=json.loads(p.read_text())
 print(p, 'synthesis='+d['synthesis']['verdict'])
 for r in d['reviewers']:
  print(r['model'], 'present='+str(r.get('present')), 'verdict='+str(r.get('result',{}).get('verdict')))
  if p.name.endswith('21-13-09Z.json'):
   result=r.get('result',{})
   print('strongest_objection:',result.get('strongest_objection'))
   print('proposed_fix:',result.get('proposed_fix'))
PY
```

Observed all three syntheses `blocked`. Files: `w-marketing-capture-2026-10-08T13-05-07Z.json`, `w-marketing-capture-2026-10-08T13-11-05Z.json`, `w-marketing-capture-2026-10-14T21-13-09Z.json`. R1/r2 have absent sonnet, never an approval. R3 has three present rejects: truthful CLAIMED sourcing/kind checks, repeat-GET/still-mode premises, and content-bearing raw-byte stability. R3 proposed fixes were read in full; current r4's dispositions, measurement rows and fallback were read in E4. No new quorum was run.

### E14 — historical baseline red body

`sed -n '25,80p' design_docs/planned/w-marketing-capture-sprint-plan.md` → exact baseline command/results and seven failure bodies summarized above, plus standalone contracts test rc1/four vacuous skips. The text explicitly says not to route around red base gates and distinguishes standalone evidence failure from the aggregate AIL policy. Not rerun.

### E15 — historical broad-suite red bodies

- `sed -n '118,157p' /private/tmp/verify_go_iter302c.log` → verifygate invalid_ledger deadline failure, marketingcheck shell-hup killed/28.896924125 s/death-assertion failures.
- `sed -n '97,112p' /private/tmp/verify_go_iter302_race.log` → `panic: test timed out after 8m0s`, running TestWorkbenchPacketRosterTestedAndProven/test-result-missing.
- `sed -n '247,254p' /private/tmp/verify_go_iter302_race.log` → shell-hup killed/29.684643292 s/death-assertion failure, package FAIL.

These read commands establish recorded failure contents, not current reproduction or environmental attribution. Source-body non-change of the inherited supervisor test is controlled by E6's positive package-declaration diff.

### E16 — historical positive AIL/CI artifacts

- `tail -4 /private/tmp/verify_ail_iter302.log` → package 9/9, PUB011 ratchet, aggregate `verify gate PASSED: 16 required identities verified, 40 named tests pass, 496 se-tools named tests pass`.
- The `cat` command in E11 also read both exact-head CI verdict files → `CI for f615a9698370694664fbbc01431e4a2c549377cc: in_progress:none`, then `completed:success`. Historical local summaries only; no current remote check.

### E17 — fresh comparison against all 12 retained protected observations

```sh
python3 - <<'PY'
import hashlib,json
from pathlib import Path
p=Path('/Users/voightkampff/.ailang/state/world-iter302/snap/obs-exec-before.json')
d=json.loads(p.read_text())
print(type(d).__name__,len(d))
print(d[0])
for row in d:
 q=Path(row['path']);b=q.read_bytes()
 print(row['path'],'sha_equal='+str(hashlib.sha256(b).hexdigest()==row['sha256']),'lines_equal='+str(b.count(b'\n')==row['lines']))
PY
```

Observed `list 12`; every enumerated row printed `sha_equal=True lines_equal=True`: charter, go.mod, go.sum, five tools/launchd files including its testdata, and four world/*.ail files. The charter observation is 6,255 lines, hash `3c3a16bf49b31d7d501bed465468d030fcbf2e8ff71d5dd6f162df42bd498e79`. This is a fresh equality measurement against a positive enumerated twelve-file control, not proof of historical execution or a full filesystem audit. Originals were only read.

### E18 — report write-boundary verification

After writing this report, `git status --short --branch` printed the same recovery branch and only `?? designer-recovery-audit.md`; the report itself is the same-surface positive control for the absence of other status entries. `git rev-parse HEAD` again printed `f615a9698370694664fbbc01431e4a2c549377cc`. `wc -c designer-recovery-audit.md` measured 24,251 bytes before adding this verification paragraph, exceeding the required 200 bytes. `git diff --check` returned rc0 with empty output; that command checks tracked diffs, not the untracked report's Markdown or any test behavior.
