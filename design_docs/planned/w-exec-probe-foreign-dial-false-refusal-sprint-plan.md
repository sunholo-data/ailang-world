# Sprint plan — w-exec-probe-foreign-dial-false-refusal (iteration 359, row 250)

Approved design: [w-exec-probe-foreign-dial-false-refusal.md](w-exec-probe-foreign-dial-false-refusal.md), **r3**,
including §11's controller narrow-refinement carve-out after two blocked quorum rounds. The r3
verification requirements are binding; the r2 prototype bytes are unchanged. Planner base:
`b6a35f0414515724d539e090d313b3cbc486e628` (detached HEAD). One milestone, **M1**, approximately 0.5 day, **planned**.
The planner produces these four artifacts only. The executor applies the exact patch; the
controller runs listener tests/mutations outside the sandbox, commits, pushes and verifies Linux CI.

M1 replaces forward-leg raw accept counting with nonce-classified contiguous windows. Child and
silent accepts refuse immediately; foreign-only windows retry up to three attempts, then refuse.
Both legs use the same `readProbeLine`, `judgeLeg`, `execProbeHelloWait` and `execProbeAttempts`.
The reverse classifier retains its protocol-specific EOF/prefix rules and the loops retain their
protocol-specific work (§3.6); its verdicts, retry and reset semantics must stay unchanged.
This is Go host-boundary socket/client behavior, not a pure AILANG package or kernel expansion.

## M1 — exact five-file implementation and byte identities

| File | Required change | Pre-change sha256 | Post-change sha256 (scratch measured) |
|---|---|---|---|
| `.github/workflows/ci.yml` | Append TF/TS/TE/TC to third tests= line in existing real-srt verbose step | `41a0a58611565fe148faaac0327dc34793550f12ddef74e7cc3c33a250c4697d` | `4efb18438d736bb41a0770b3d4fd0b12e3ff5360e89dc8bc156e2ab44f67521f` |
| `host/broker/exec_probe.go` | Client dial line; classified windows; shared reader, judge, deadline and attempt bound (+237/−53) | `19de369381579387bd9d83aceae28a5912717dab00b283454ca0067d41d087bd` | `110fffc3f31a10d71a024b3e02f66b0864178ce59e44e3e5b5be2c82ac11c79f` |
| `host/broker/exec_probe_net6_dial_test.go` | New 305-line test file, TF/TS/TE/TC and deterministic dial rig | `ABSENT` | `8c7f9602c1105bdca99af2326a6e51851d7011e1ccb875f18a293f4f91838b0e` |
| `host/broker/exec_probe_net6_test.go` | Nonce-bearing connect-argument regex (+1/−1) | `400faf0ffc25a3f0a79302cc9436349da6e445b998744b5767f6d706475a64aa` | `8ec970891599257aa280436f78cfac427ebaccc1cf5cd33d404904cf60ec3b26` |
| `host/broker/exec_probe_test.go` | TP raw/proxy diagnostics and runReverseLeg string argument (+3/−3) | `bbfccede19236d17b81d40401cc079aeb9987ff356c6a08affe6bc2a159e80d0` | `ee59ba2cc5742855f6aa172cb7873484a5c07abf2f1b73d5792d8673c45e7366` |

`ABSENT` is absence, not an empty-file hash. The implementation file list has exactly five paths.
The four planning artifacts are separate from that implementation scope.

Banked diff: `design_docs/verification/world-iter359/plan-M1.diff.txt` (821 lines), sha256
`7fb2537268498db0b6db58d05b59cb0d52add1bb20178d16c3025dbfb3eeaf85`.
It is byte-identical to `designer-proto.patch.txt`; **no additional CI hunk** is allowed.
The designer patch already contains the sole ci.yml edit. Preserve the bytes: do not rewrite,
reformat, supplement or repair them. A0/A1/A5 enforce the bank and implementation identities.

## Acceptance gates (M1 remains unaccepted until A9)

All blocks run from the target tree root in Bash with `set -euo pipefail`. A0 precedes application;
A1 is applied by the executor once (if already applied, still verify A1 hashes).
A2–A9 follow on the applied tree. Use the released controller AILANG pin, never a dirty binary;
`verify_go.sh` enforces the repo's current release and runs the required `go test ./... -count=1`
with evidence-manifest checks. Node must be present; a SKIP cannot satisfy A4/A6/A9.
Run suites and mutations sequentially. Recheck post-apply hashes after mutations and before push.
The planner must not execute A3–A7's listener tests or the runner. Loopback binds are denied:
any listener-opening host/broker result here is **UNINFORMATIVE UNDER SANDBOX**, never pass/fail.

Controller setup (evidence paths are inside the tree):

```bash
set -euo pipefail
export AILANG_BIN="$HOME/.pinned-ailang/ailang"
export EVIDENCE="$PWD/design_docs/verification/world-iter359/controller-M1"
mkdir -p "$EVIDENCE/tmp"
export TMPDIR="$EVIDENCE/tmp"
go version
node --version
"$AILANG_BIN" --version
```

**A0 — base and bank identity, before application.**

```bash
test "$(git rev-parse b6a35f0414515724d539e090d313b3cbc486e628^{commit})" = b6a35f0414515724d539e090d313b3cbc486e628
test ! -e host/broker/exec_probe_net6_dial_test.go
shasum -a 256 -c <<'SHA'
41a0a58611565fe148faaac0327dc34793550f12ddef74e7cc3c33a250c4697d  .github/workflows/ci.yml
19de369381579387bd9d83aceae28a5912717dab00b283454ca0067d41d087bd  host/broker/exec_probe.go
400faf0ffc25a3f0a79302cc9436349da6e445b998744b5767f6d706475a64aa  host/broker/exec_probe_net6_test.go
bbfccede19236d17b81d40401cc079aeb9987ff356c6a08affe6bc2a159e80d0  host/broker/exec_probe_test.go
SHA
cmp design_docs/verification/world-iter359/designer-proto.patch.txt design_docs/verification/world-iter359/plan-M1.diff.txt
shasum -a 256 -c <<'SHA'
7fb2537268498db0b6db58d05b59cb0d52add1bb20178d16c3025dbfb3eeaf85  design_docs/verification/world-iter359/plan-M1.diff.txt
SHA
git apply --check design_docs/verification/world-iter359/plan-M1.diff.txt
```

**A1 — executor application, exactly once.**

```bash
git apply design_docs/verification/world-iter359/plan-M1.diff.txt
shasum -a 256 -c <<'SHA'
4efb18438d736bb41a0770b3d4fd0b12e3ff5360e89dc8bc156e2ab44f67521f  .github/workflows/ci.yml
110fffc3f31a10d71a024b3e02f66b0864178ce59e44e3e5b5be2c82ac11c79f  host/broker/exec_probe.go
8c7f9602c1105bdca99af2326a6e51851d7011e1ccb875f18a293f4f91838b0e  host/broker/exec_probe_net6_dial_test.go
8ec970891599257aa280436f78cfac427ebaccc1cf5cd33d404904cf60ec3b26  host/broker/exec_probe_net6_test.go
ee59ba2cc5742855f6aa172cb7873484a5c07abf2f1b73d5792d8673c45e7366  host/broker/exec_probe_test.go
SHA
```

**A2 — source formatting, test compilation, build and runner syntax.**

```bash
gofmt -l host/broker/exec_probe.go host/broker/exec_probe_net6_dial_test.go host/broker/exec_probe_net6_test.go host/broker/exec_probe_test.go > "$EVIDENCE/gofmt.log"
test ! -s "$EVIDENCE/gofmt.log"
go vet ./host/broker/
go build ./...
/bin/bash -n design_docs/verification/world-iter359/plan-mutants.sh.txt
```

**A3 — blocking r3 executor relay verification against the actual pin.**

```bash
: "${WORLD_EXEC_SRT_NODE_MODULES:?point to the pinned srt install used for this build}"
python3 -I - "$WORLD_EXEC_SRT_NODE_MODULES" "$EVIDENCE" <<'PYCODE'
from pathlib import Path
import hashlib,json,sys
pkg=Path(sys.argv[1])/'@anthropic-ai/sandbox-runtime'
e=Path(sys.argv[2])
pin=json.loads(Path('design_docs/verification/world-row140-m0/pin.json').read_text())
assert json.loads((pkg/'package.json').read_text())['version']==pin['version']
assert hashlib.sha256((pkg/'dist/cli.js').read_bytes()).hexdigest()==pin['cli_sha256']
parts=[]; hashes={}
for name in ('http-proxy.js','sandbox-manager.js','mux-proxy.js'):
    matches=list((pkg/'dist').rglob(name))
    assert len(matches)==1,(name,matches)
    p=matches[0]; hashes[str(p)]=hashlib.sha256(p.read_bytes()).hexdigest()
    # Full source with line numbers: inspect every CONNECT route and relay branch.
    parts.append('\nFILE '+str(p)+'\n'+'\n'.join('%d: %s'%(i,l) for i,l in enumerate(p.read_text().splitlines(),1)))
settings=Path('host/broker/exec_settings.go')
parts.append('\nWORLD SETTINGS\n'+settings.read_text())
# Instrument-health scan only; source review, not hit counts, discharges the relay premise.
for folder in ('host','cmd'):
    for p in Path(folder).rglob('*.go'):
        if p.name.endswith('_test.go'): continue
        for i,line in enumerate(p.read_text().splitlines(),1):
            if any(k in line for k in ('tlsTerminate','mitmProxy','parentProxy','allowedDomains')):
                parts.append('%s:%d: %s'%(p,i,line))
(e/'relay-source.txt').write_text('\n'.join(parts)+'\n')
(e/'relay-inputs.json').write_text(json.dumps({'version':pin['version'],'cli_sha256':pin['cli_sha256'],'sources':hashes},indent=2)+'\n')
PYCODE
# Executor reads relay-source.txt and writes relay-review.json with:
# {"verbatim_both_directions":true,"upstream_header_or_banner":false,
#  "world_config_enables_altering_route":false,"reviewer":"...",
#  "evidence":"handler/config paths and line numbers, explanation"}.
# These are review conclusions, never values generated from grep hit counts.
python3 -I - "$EVIDENCE" <<'PYCODE'
from pathlib import Path
import json,sys
e=Path(sys.argv[1]); r=json.loads((e/'relay-review.json').read_text())
assert r['verbatim_both_directions'] is True
assert r['upstream_header_or_banner'] is False
assert r['world_config_enables_altering_route'] is False
assert r['reviewer'].strip() and r['evidence'].strip()
with (e/'verification-log.md').open('a') as out:
    out.write('\nExecutor relay verification (blocking r3 gate):\n'+(e/'relay-inputs.json').read_text()+(e/'relay-review.json').read_text()+'\n')
PYCODE
```

**A4 — each named test alone on the applied patch, outside sandbox.**

```bash
for name in TestExecProbePassThroughShimFailsEveryConfinementArm \
  TestExecProbeNet6ForeignDialIsRetriedNotRefused \
  TestExecProbeNet6ChildOrSilentAcceptStillFails \
  TestExecProbeNet6ForeignEveryAttemptRefuses \
  TestExecProbeNet6DialClassifier \
  TestExecProbeNet6BoundBeforeAnyPick \
  TestExecProbeNet6BindSiteIsBeforeTheFanOut \
  TestExecProbeReverseLegJudge \
  TestExecProbeReverseLegChattyForeignHolderIsRetried \
  TestExecProbeReverseLegExhaustionRefusesDistinctly \
  TestExecProbeReverseLegRealReachSurvivesAForeignAttempt \
  TestExecProbeReverseLegNoTokenIsNotRetried \
  TestExecProbeForeignListenerOnTheReversePortDoesNotRefuse \
  TestExecProbeReverseLegSilentHolderRefusesSticky \
  TestExecProbeReverseDialErrIsSilent \
  TestExecProbeReversePollRecordsConnectReset \
  TestExecProbeReverseHelloWaitBounds \
  TestExecProbeReversePollCadence \
  TestExecProbeChildListenerDrainsTheHello; do
  log="$EVIDENCE/$name.log"
  rc=0
  go test -count=1 -run "^(${name})$" ./host/broker/ -v -timeout 300s > "$log" 2>&1 || rc=$?
  cat "$log"
  test "$rc" -eq 0
  /usr/bin/grep -Fq -- "--- PASS: $name " "$log"
  if /usr/bin/grep -q -- '--- SKIP:' "$log"; then exit 1; fi
done
```

**A5 — measured mutations and HEAD-rule positive control, controller outside sandbox.**

```bash
bash design_docs/verification/world-iter359/plan-mutants.sh.txt "$PWD" > "$EVIDENCE/mutants.log" 2>&1
cat "$EVIDENCE/mutants.log"
test "$(/usr/bin/grep -c '^EXPECTATION MATCH:' "$EVIDENCE/mutants.log")" -eq 177
if /usr/bin/grep -q '^EXPECTATION MISMATCH:' "$EVIDENCE/mutants.log"; then exit 1; fi
test "$(/usr/bin/grep -c '^RESTORE cmp ok: .*tests=' "$EVIDENCE/mutants.log")" -eq 177
/usr/bin/grep -Fq 'GO FILE CENSUS unchanged; no extra .go files remain' "$EVIDENCE/mutants.log"
shasum -a 256 -c <<'SHA'
4efb18438d736bb41a0770b3d4fd0b12e3ff5360e89dc8bc156e2ab44f67521f  .github/workflows/ci.yml
110fffc3f31a10d71a024b3e02f66b0864178ce59e44e3e5b5be2c82ac11c79f  host/broker/exec_probe.go
8c7f9602c1105bdca99af2326a6e51851d7011e1ccb875f18a293f4f91838b0e  host/broker/exec_probe_net6_dial_test.go
8ec970891599257aa280436f78cfac427ebaccc1cf5cd33d404904cf60ec3b26  host/broker/exec_probe_net6_test.go
ee59ba2cc5742855f6aa172cb7873484a5c07abf2f1b73d5792d8673c45e7366  host/broker/exec_probe_test.go
SHA
```

**A6 — r3 stability on refactored bytes, outside sandbox.**

```bash
go test -count=3 -run '^(TestExecProbePassThroughShimFailsEveryConfinementArm|TestExecProbeNet6.*)$' -v ./host/broker/ > "$EVIDENCE/stability.log" 2>&1
cat "$EVIDENCE/stability.log"
python3 -I - "$EVIDENCE/stability.log" <<'PYCODE'
from pathlib import Path
import collections,re,sys
s=Path(sys.argv[1]).read_text()
names=('TestExecProbePassThroughShimFailsEveryConfinementArm','TestExecProbeNet6ForeignDialIsRetriedNotRefused','TestExecProbeNet6ChildOrSilentAcceptStillFails','TestExecProbeNet6ForeignEveryAttemptRefuses','TestExecProbeNet6DialClassifier','TestExecProbeNet6BoundBeforeAnyPick','TestExecProbeNet6BindSiteIsBeforeTheFanOut')
passes=collections.Counter(re.findall(r'^--- PASS: (\w+) ',s,re.M))
assert all(passes[n]==3 for n in names),passes
assert '--- SKIP:' not in s
PYCODE
```

**A7 — full broker, AILANG and repo gates under the released pin, outside sandbox.**

```bash
: "${AILANG_BIN:?export the controller pinned released interpreter}"
"$AILANG_BIN" --version > "$EVIDENCE/ailang-version.txt"
cat "$EVIDENCE/ailang-version.txt"
go test -count=1 -timeout 600s ./host/broker/ > "$EVIDENCE/broker.log" 2>&1
cat "$EVIDENCE/broker.log"
./scripts/verify_ail.sh > "$EVIDENCE/verify-ail.log" 2>&1
cat "$EVIDENCE/verify-ail.log"
go build ./...
go vet ./host/broker/
./scripts/verify_go.sh > "$EVIDENCE/verify-go.log" 2>&1
cat "$EVIDENCE/verify-go.log"
```

**A8 — CI list wiring, instrument health only.**

```bash
python3 -I - <<'PYCODE'
from pathlib import Path
s=Path('.github/workflows/ci.yml').read_text()
a=s.index('      - name: Row 140 Workspace.Exec real-srt tests, verbose (linux bubblewrap)')
b=s.index('      - name:',a+1); block=s[a:b]
lines=[l for l in block.splitlines() if l.lstrip().startswith('tests=')]
assert len(lines)==4
names=('TestExecProbeNet6ForeignDialIsRetriedNotRefused','TestExecProbeNet6ChildOrSilentAcceptStillFails','TestExecProbeNet6ForeignEveryAttemptRefuses','TestExecProbeNet6DialClassifier')
assert 'TestExecSrtStartupProbePasses' in lines[1]
for n in names:
    assert sum(l.count(n) for l in lines)==1,n
    assert n in lines[2],n
for text in ('WORLD_EXEC_SRT_NODE_MODULES','RUNNER_OS','go test ./host/broker/ ./host/daemon/','for t in $tests; do','^--- PASS: $t '):
    assert text in block,text
PYCODE
```

**A9 — mandatory Linux real-srt CI acceptance after controller push.**

```bash
: "${CI_RUN_ID:?controller selects the PR CI run for the pushed implementation commit}"
: "${CI_HEAD_SHA:?full pushed implementation commit SHA}"
gh run view "$CI_RUN_ID" --json headSha,conclusion,jobs,url > "$EVIDENCE/linux-ci.json"
gh run view "$CI_RUN_ID" --log > "$EVIDENCE/linux-ci.log"
python3 -I - "$EVIDENCE" "$CI_HEAD_SHA" <<'PYCODE'
from pathlib import Path
import json,sys
e=Path(sys.argv[1]); run=json.loads((e/'linux-ci.json').read_text())
assert run['headSha']==sys.argv[2]
assert run['conclusion']=='success',run['conclusion']
title='Row 140 Workspace.Exec real-srt tests, verbose (linux bubblewrap)'
steps=[step for job in run['jobs'] for step in job['steps'] if step['name']==title]
assert len(steps)==1 and steps[0]['conclusion']=='success',steps
log=(e/'linux-ci.log').read_text()
for n in ('TestExecSrtStartupProbePasses','TestExecProbeNet6ForeignDialIsRetriedNotRefused','TestExecProbeNet6ChildOrSilentAcceptStillFails','TestExecProbeNet6ForeignEveryAttemptRefuses','TestExecProbeNet6DialClassifier'):
    assert '--- PASS: '+n+' ' in log,n
with (e/'verification-log.md').open('a') as out:
    out.write('\nLinux real-srt acceptance: '+title+' SUCCESS; '+run['url']+'; commit '+run['headSha']+'; PASS per name, no SKIP substitution.\n')
print(run['url'])
PYCODE
```

A3 is a human/source-inspection gate as well as runnable collection and receipt checks: a
zero grep count cannot certify byte fidelity. Record the install identity, inspected source hashes,
handler line references, both answers and World's config routes in `controller-M1/verification-log.md`.
V21 is prior controller evidence, not a substitute for the executor inspecting the actual pin.

**Executor verification (blocking; kimi's §7 text, verbatim):** with the pinned srt install,
grep/inspect the proxy's CONNECT handler and record in the sprint Verification Log (a) whether
bytes after a 200 are relayed verbatim in both directions, (b) whether any header/banner is
prepended or appended upstream. If any alteration exists, the design reverts to
fail-at-first-foreign (no retry on windows containing byte-bearing accepts; retry only nothing),
i.e. reclassify any accept carrying bytes that begin neither expected line as silent. The
controller's own reading of the pinned 0.0.78 source (V21) found no alteration; the executor
confirms it against the pin the patch is built with and records the result.

If this check finds alteration, stop acceptance and report the finding to the controller. That
fallback changes approved bytes and therefore needs a revised bank/acceptance direction; M1 does
not silently edit the quorum-approved diff. See finding F2 below.

A6 carries gemini's exact `-count=3` command on the applied r2 bytes. V12 was r1; V22 is historical
r2 controller evidence and does not replace this executor rerun. A4 separately runs all twelve
reverse tests. MUT-RAW-COUNT/TF in A5 is the r3 carve-out's HEAD-rule positive control on r2's tests
(V19/V22): expected RED, then pristine CONTROL/TF GREEN. Do not substitute the old whole HEAD
source file into r2 tests: TC requires new identifiers, and the reverse helper signature changed.
Historical parent measurement V8 is retained; MUT-RAW-COUNT reproduces HEAD's rule in the new structure.

A9 is the controller's **post-push acceptance gate**, not just a wiring grep. The existing Linux
step must run with bwrap, the pinned real srt install and AILANG_BIN. ci.yml:246 already includes
`TestExecSrtStartupProbePasses` and TP; :247 includes the twelve reverse tests, row-242 tests and
the four new tests. The one combined list is run via `go test ./host/broker/ ./host/daemon/
-count=1 -p 1 -timeout 600s -v -run ...`; every listed name must print `--- PASS`. A SKIP fails
the step. The executor/evaluator and controller record the step conclusion, pushed SHA and run
URL in the sprint Verification Log. No local green, V21 static read or A8 wiring check substitutes
for this first dynamic Linux real-srt measurement of the new client. CI is UNMEASURED by the planner.

## Test identities and mutation expectations (design §5)

Each runner row selects **one** exact top-level test with `go test -count=1 -run '^(<name>)$'
./host/broker/ -v -timeout 300s`. Rows are independent: restore → one mutation → one named test →
compare expected rc and assertion evidence → restore → cmp. RED means expected rc 1 with the
named test's FAIL and assertion text; GREEN means rc 0 with its PASS. Compilation errors, no
matched test, SKIP, crash or timeout never count as a killed mutant. TP is expected RED only
where the designer actually measured it RED.

| Alias | Exact test name | Load-bearing assertion |
|---|---|---|
| TP | `TestExecProbePassThroughShimFailsEveryConfinementArm` | Each raw/proxy diagnostic contains `child 1, silent 0` through the real node client |
| TF | `TestExecProbeNet6ForeignDialIsRetriedNotRefused` | Arm 6 has no refusal; each op count is at least planned retries and at most 3 |
| TS | `TestExecProbeNet6ChildOrSilentAcceptStillFails` | Exact attempt/class prefix is present and op runs equal expected N; child/silent never retry |
| TE | `TestExecProbeNet6ForeignEveryAttemptRefuses` | Inconclusive after 3 attempts and exactly 3 connect runs |
| TC | `TestExecProbeNet6DialClassifier` | Window wait exceeds line deadline; each byte shape has exact class; elapsed time is bounded |

**Rev** is exactly these twelve named tests, each run separately by A4 and individually by the
runner for CONTROL and the shared-piece/window scenarios:

- `TestExecProbeReverseLegJudge`
- `TestExecProbeReverseLegChattyForeignHolderIsRetried`
- `TestExecProbeReverseLegExhaustionRefusesDistinctly`
- `TestExecProbeReverseLegRealReachSurvivesAForeignAttempt`
- `TestExecProbeReverseLegNoTokenIsNotRetried`
- `TestExecProbeForeignListenerOnTheReversePortDoesNotRefuse`
- `TestExecProbeReverseLegSilentHolderRefusesSticky`
- `TestExecProbeReverseDialErrIsSilent`
- `TestExecProbeReversePollRecordsConnectReset`
- `TestExecProbeReverseHelloWaitBounds`
- `TestExecProbeReversePollCadence`
- `TestExecProbeChildListenerDrainsTheHello`

The following table mirrors the complete design §5 table, including the historical HEAD row and
both surviving mutants. R = RED/rc 1, G = GREEN/rc 0; — = not measured/not a runner expectation.
Rev reds name every measured failing test; other Rev tests are GREEN where the column was measured.
The runner's 177 rows cover CONTROL + all twenty measured r2 mutants: five forward tests per
scenario (105), plus twelve reverse tests for CONTROL and the five shared/window scenarios (72).
The historical HEAD row is not an additional runner scenario. No unmeasured extra mutant is added.

| Mutant | TP | TF | TS | TE | TC | Rev expected (named red tests; others green) | Edit/assertion that fires |
|---|---|---|---|---|---|---|---|
| CONTROL | G | G | G | G | G | G: all 12 | None; all 17 identities PASS |
| HEAD (historical positive control, V8) | — | R | R (texts only; all shapes refused) | R | — | — | Parent raw-count rule: TF raw/proxy/raw-twice-proxy false refusals (`accepts 2`); represented by MUT-RAW-COUNT on r2 tests |
| MUT-SHARED-READ-SHORT | R | R | R | G | R | `TestExecProbeReverseLegRealReachSurvivesAForeignAttempt`<br>`TestExecProbeReverseHelloWaitBounds` | One byte short: TF control never whole; TS 4 child rows; TC control/child/child-then-more; TP child diagnostic absent; reverse slow-child has foreign instead of child |
| MUT-SHARED-JUDGE-SILENT-RETRIED | G | G | R | G | G | `TestExecProbeReverseLegJudge`<br>`TestExecProbeForeignListenerOnTheReversePortDoesNotRefuse`<br>`TestExecProbeReverseLegSilentHolderRefusesSticky` | Silent → inconclusive: TS raw-silent/raw-scan/proxy-silent lose sticky refusal; reverse judge verdict and sticky-silent assertions fire |
| MUT-SHARED-ATTEMPTS-2 | G | R | G | R | G | `TestExecProbeReverseLegExhaustionRefusesDistinctly` | Bound 2: TF raw-twice-proxy refuses; TE and reverse exhaustion require 3 attempts/runs |
| MUT-SHARED-WAIT-1500MS | G | G | G | G | G | `TestExecProbeForeignListenerOnTheReversePortDoesNotRefuse`<br>`TestExecProbeReverseHelloWaitBounds` | Deadline 1500 ms: reverse HelloWaitBounds ceiling 950 ms fires; reverse silent row loses silent text; forward bound tests scale |
| MUT-WINDOW-WAIT-NOT-ABOVE | G | G | G | G | R | G: all 12 | Window wait = hello wait: TC first assertion requires strict > |
| MUT-RAW-COUNT | G | R | R | R | G | — | Any accept refuses (HEAD rule): TF false-refusal assertion; TS raw-child-on-retry expects attempt 2; TE requires exhaustion/3 runs |
| MUT-SILENT-AS-FOREIGN | G | G | R | G | R | — | Silent tail → foreign: TS 3 silent rows lose refusal; TC connect-close/held-silent/prefix rows get wrong class |
| MUT-PREFIX-AS-FOREIGN | G | G | G | G | R | — | Nonempty prefix → foreign: TC two prefix class assertions |
| MUT-CHILD-AS-FOREIGN | R | G | R | G | R | — | Child → foreign: TS 4 child rows lose reach; TC child class and TP child diagnostics fire |
| MUT-CHILD-AS-CONTROL | R | G | R | G | R | — | Child → control: window maps it silent, still fail-closed; TS/TP exact child diagnostic and TC class fire |
| MUT-TOKEN-ONLY | R | R | R | R | G | — | Ignore window: TS all 7 lose refusal; TF required retries absent; TE 3-run/exhaustion assertion; TP child diagnostics absent |
| MUT-RETRY-CHILD | G | G | R | G | G | — | Child retried: TS 4 child rows false pass, missing exact attempt/class prefix |
| MUT-RETRY-SILENT | G | G | R | G | G | — | Silent retried: TS 3 silent rows lose sticky refusal |
| MUT-EXHAUST-PASSES | G | G | G | R | G | — | Third foreign window passes: TE inconclusive/3-runs assertion |
| MUT-UNBOUNDED | G | G | G | R | G | — | Remove bound: fourth clean run passes, TE requires 3 and exhaustion (not a timeout kill) |
| MUT-FROM-NOT-ADVANCED | R | R | G | G | G | — | Keep from: TF sees control again and refuses; TP proxy window includes earlier reach, exact child diagnostic fails |
| MUT-WINDOW-AT-JUDGE | R | R | R | R | G | — | Begin at judgement: all 7 TS cases lose refusal; TF retries missing; TE loses exhaustion; TP child diagnostics absent |
| MUT-CLIENT-NO-LINE | R | G | G | G | G | — | connect destroys: TP must say child 1, silent 0 on both legs |
| MUT-PENDING-AS-FOREIGN | G | G | G | G | G | — | SURVIVES: pending branch not reached by measured tests; normal classifier deadline 500 ms vs window 1.5 s; see F3 |
| MUT-CONTROL-ANY-ACCEPT | G | G | G | G | G | — | SURVIVES: any accept satisfies control; no test coordinates a foreign pre-control race; TC/TP independently cover classification, not control linkage |

All R/G cells above are designer measurements (`designer-r2-mutants.txt`, V19), not planner
runtime results. In particular shared read/judge/attempt-bound mutations each have separate
forward and reverse expected-red rows. Deadline value 1500 ms is pinned only by the reverse
ceiling; the forward tests scale with it. MUT-WINDOW-WAIT-NOT-ABOVE separately pins the forward
strict inequality. Their expected patterns deliberately differ.

The script takes the tree root as `$1`, canonicalizes it, checks all five applied hashes, and
stores logs, pristine source, Go temporary files, modules, toolchain and build caches under
`$1/.plan-mutants.XXXXXX`. It leaves Go's toolchain selection at its default (go.mod is go 1.26.6).
Pristine is `exec_probe.go.pristine`, never a `.go` file. Each row restores source and verifies
`cmp`; EXIT/signal traps restore again. Unexpected test results accumulate mismatches; runner
returns 1 for any mismatch and 3 for restoration/census failure. At exit it deletes pristine
and the whole runtime/cache tree (downloaded modules/toolchains include .go files), and compares
before/after inventories of every .go path. Logs remain for review. No helper .go file is created.
Bash syntax and mutation anchors were checked without executing Go tests. Finding F1 states the
existing fixed-path exception that a portable shell runner cannot redirect.

## Residuals carried from design §6

- **R1 — remaining false refusals.** No bytes within 500 ms, a stalled probe-line prefix, or
  foreign dials in all three windows still refuse closed. The known start page sends GET within
  2 ms every ~31.8 s; its expected contribution becomes zero by the period/window argument,
  not by a measured startup census. Other hosts' silent dialers are UNMEASURED.
- **R2 — quantified weakening.** A byte-altering and intermittent escape with reach probability
  p per attempt passes with probability 1 − p³, versus HEAD's 1 − p (assuming independent,
  stationary attempts). A deterministic altering reach still exhausts. V21's static pinned
  0.0.78 source inspection found direct CONNECT `head` forwarding and `pipe` in both directions;
  TLS termination, MITM and parent proxy need config keys absent from World's four-field network
  struct and non-test Go callers. Thus no instance is known under this pin/config. A3 must
  verify that premise on the actual install; future pin/config changes reopen it. This is a
  static source premise, not dynamic Linux evidence; A9 stays required.
- **R3 — latency.** Empty windows return immediately; foreign bytes classify at the first
  divergence. Each retry adds one small run + 200 ms settle (~0.3 s measured by designer).
  Silent accepts cost up to the 500 ms read deadline, then refuse.
- **R4 — Linux/real srt.** UNMEASURED locally. Expected raw/proxy refusals leave empty windows;
  the new `s.end(line)` runs only on successful connects. The mandatory post-push real-srt step
  must PASS the startup probe and four new tests by name; no SKIP substitution.
- **R5 — nonce in argv.** A local process reading the child nonce can cause refusal by sending
  it; it cannot create a pass. The control nonce never leaves World/appears in argv.

## M1 must NOT touch

- Any implementation path outside the exact five listed above; no extra CI hunk, helper .go
  file, seam, `t.Parallel()`, new dependency, public API or generic attempt-loop driver.
- Reverse-leg verdict/retry/reset semantics (especially parked row 233, D-WORLD-105), hello
  wait values, listen duration, poll cadence or token precedence; only the banked shared-code
  extraction/constant rename and helper argument adaptation are authorized.
- Row-242 bind hoist/lifetime, arm order, other arms, control/reverse diagnostics, the control's
  2 s dial/wait bounds, settings rendering, sandbox pin or srt installation/configuration.
- `.ail`, kernel/package projection, go.mod/go.sum, QUICKSTART, startup callers/daemon, the
  Sunholo start page, silent-dialer mitigation or arbitrary Linux descriptor refinements.
- The approved design, designer evidence/tools, prototype patch, plan diff, plan/JSON or mutation
  runner during execution; record executor/controller/evaluator evidence in new log paths.
- Mission charter/status/dashboard/index/log machinery, `tools/launchd/*`, shared skills, V1
  checkout/state, other worktrees, Git index/history/metadata. Controller owns commits and push.
- The r3 fallback as an unreviewed patch alteration. If relay verification contradicts the
  premise, report it and hold acceptance for controller disposition (F2).

## Planner Verification Log and findings

The planner read CLAUDE.md, coding standards, approved design in full (including §§3.6, 4–7,
10–11), prototype, r2 tools/results, HEAD sources and CI step, prior iteration-340 plan/runner.
The scratch procedure, from this detached worktree root, was:

```bash
mkdir .plan-scratch
bash -o pipefail -c 'git archive HEAD | tar -x -C .plan-scratch/'
cp design_docs/verification/world-iter359/designer-proto.patch.txt design_docs/verification/world-iter359/plan-M1.diff.txt
cmp design_docs/verification/world-iter359/designer-proto.patch.txt design_docs/verification/world-iter359/plan-M1.diff.txt
(
  cd .plan-scratch
  export GIT_CEILING_DIRECTORIES="$(dirname "$PWD")"
  git apply --check --verbose ../design_docs/verification/world-iter359/plan-M1.diff.txt
  git apply --verbose ../design_docs/verification/world-iter359/plan-M1.diff.txt
)
```

Both final apply commands returned **rc 0**, explicitly checked/applied all five paths. The
post-change hashes in the M1 table were measured after this application, not copied from V13.
The initial unisolated Git invocation returned 0 while skipping every patch: Git discovered the
parent worktree and filtered paths relative to `.plan-scratch`. The missing new file exposed that
false-green instrument. Setting the discovery ceiling to the **parent of scratch** made Git
apply outside a repository; no Git index/history/worktree metadata was written. The implementation
files in the detached planner worktree remain unchanged.

| Permitted planner check | Pristine archive rc | Applied archive rc |
|---|---|---|
| `go vet ./host/broker/` | 0 | 0 |
| `go build ./...` | 0 | 0 |
| `gofmt -l` on the four patched Go paths | new file ABSENT | 0, empty output |
| `/bin/bash -n plan-mutants.sh.txt` | artifact only | 0 (Bash 3.2 compatible) |

Toolchain: go1.26.6 darwin/arm64 (matches go.mod). TMPDIR was set to the already-created
`.plan-scratch/.planner-tmp`; allowed Go caches used their existing locations. No package suite,
listener test or mutation runner was run by the planner; such results would be
**UNINFORMATIVE UNDER SANDBOX**. The controller measured base formatting/vet/build rc 0 too.
JSON validation, every acceptance block's Bash syntax, embedded Python syntax, all twenty
mutation anchors against applied scratch source, runner/table expectations against the banked
r2 results, cmp identity and five-path diff inventory were checked statically. Scratch was
removed after final validation. Runtime/mutation/stability/full-suite and Linux gates remain
controller work; historical V8/V17/V19/V21/V22 are attributed designer/controller evidence.

Findings (no implementation fix in this plan):

1. **F1 — literal all-process write confinement is incompatible with the approved probe.**
   `probe()` writes `/tmp/claude/world-exec-probe-...` and its control directly; `plantExecProbe()`
   may mkdir `/tmp/claude`, independently of TMPDIR. Runner-owned files/caches stay under `$1`,
   but a portable Bash runner cannot guarantee its Go-test children write nowhere else while
   testing the exact approved source. Full confinement needs controller-supplied OS path
   isolation; no new overlay, source rewrite, external symlink or unapproved implementation
   workaround is added. The runner calls this out explicitly and cannot certify the literal
   stronger property. It does ensure no extra .go files remain under the root.
2. **F2 — fallback versus byte-exact scope.** Kimi's relay-alteration fallback requires changing
   the approved classifier and its expected outcomes. The instructions also require the unchanged
   r2 patch. V21 supports the expected no-alteration path, but if executor inspection disagrees,
   the controller must resolve this scope conflict before M1 can be accepted. A missing/negative
   relay receipt blocks acceptance; it never licenses silent patch drift.
3. **F3 — pending is not mathematically unreachable.** The strict inequality between two
   durations does not bound goroutine scheduling delay. A classifier starved longer than the
   extra second can still be pending when window expires. The approved code safely maps that
   case to silent; the survivor is a coverage residual, not proof that the branch is impossible.
   TC pins the inequality, not scheduler progress; no stall seam is authorized.
4. **F4 — security-claim scope/probability premise.** §4 Claim A's "no new false pass" applies
   per empty attempt window, not to the whole multi-attempt probe; the absolute Claim B heading
   needs §6 R2's byte-altering/intermittent qualification. The formulas 1 − p³ versus 1 − p also
   assume independent, identically distributed reach attempts. Static relay verification is
   the approved premise check, not a general proof against every future relay/configuration.
