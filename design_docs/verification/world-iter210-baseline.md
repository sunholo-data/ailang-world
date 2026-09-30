# Iteration 210 — partial live baseline, no World value credit

Observation at frozen product commit `9b2cfac20f770811b113d62e987530f7e6299660`, Darwin/arm64, pinned AILANG v0.41.0. This is ADMIN evidence. Qualifying paired questions: **0/3**. No World walk was performed; no operational World cohort was declared for this attempt. Its walk arm is **NOT-RUN(NO-CORPUS)**, not UNANSWERABLE and not proof that capture is universally absent. Capture inventory remains UNKNOWN.

## Authentic incident and controls

Question: Why did the full local `scripts/verify_go.sh` gate fail on Darwin/arm64 at this commit while remote CI on the same commit was green? Determine the failing mechanism, distinguish instrument failure from product regression, and state what is unproved.

Both full local invocations used `AILANG_BIN=/Users/voightkampff/.pinned-ailang/ailang WORLD_PKG_AILANG_BIN=/Users/voightkampff/.pinned-ailang/ailang ./scripts/verify_go.sh`. Both returned **rc1**, at the same test `TestQueryInterfaceReturnsWhileADescendantHoldsStdout` after3.00s: `instrument: descendant pid not recorded .../pid: no such file or directory` (`host/pkgproj/iface_test.go:232`). Neither reached the full race leg. Package timings36.104s and36.809s. No product edits occurred between invocations.

The focused command `go test ./host/pkgproj -run '^TestQueryInterfaceReturnsWhileADescendantHoldsStdout$' -count=5 -timeout=90s -v` under the same pin returned **rc0,5/5**, about4.02s each: descendant PID alive after timeout, then reaped. This supports context sensitivity; it does not prove scheduling/load causation. `scripts/verify_ail.sh` under the same pin returned **rc0**: worldpkg9/9,16 required verified checks,40 named tests; PUB011 uncontracted3 remains declared. Same-SHA remote CI run36682101496 had the complete expected two checks success/success. Remote green is not a substitute for the local red.

## Fresh blinded baseline

Executor was spawned through the Agent tool with `fork_turns=none`, no charter, World DB, prior diagnoses, or other-role context. Original complete host/cmd/scripts/workflow source exported using `git archive` at the frozen SHA, with no `.git`;268 corpus files including authentic original gate output, context and exact question. The original operation.log hash is `4c3a5b6110913ee0712415ee328d06fa08030b88d2b827d05cb980ffe4a22b44`.

Attempt1: proposed deny-default sandbox aborted with **rc-6**, including its allow control. No successful question read or diagnosis occurred. Preserved as **NOT-RUN(TRANSPORT)**;32.357679s was prerequisite time, not baseline performance. Host `(allow default)` positive control rc0 isolated this to the proposed profile. This was a local observation-instrument repair, not a shared fleet defect or harness change.

Attempt2: corrected read-only Seatbelt profile denied network, writes, reads under `/Users`, `/private/tmp`, `/tmp`, with a sole corpus read exception. Before reveal the executor independently measured allowed question readability rc0, denied main `.git/config` rc1, denied `~/.ailang` rc1, denied `/private/tmp` rc1. Diagnostic reads all ran inside sandbox-exec with per-command deadlines5 or10s. The profile is preserved; runtime/system read allowance is not claimed to exclude every hypothetical contamination route.

First successful question issue: **2026-09-30T10:26:33.427670Z**, monotonic1803667432099791ns. Answer fixed: **10:27:56.366350Z**, monotonic1803750370495916ns. Executor-reported archaeology duration **82.938396125s (≈83s)**, including command errors and idle. The STOP timestamp was recorded in its report but is absent from the raw transcript; precision is not independently transcript-verifiable. It excludes later independent verification; no≤300s verified-answer claim is made. Attempt1 remains preserved separately.

# Blinded Gate 2 baseline, attempt 2
Start UTC: 2026-09-30T10:26:33.427670+00:00
Stop UTC (answer fixed): 2026-09-30T10:27:56.366350+00:00
Monotonic start: 1803667432099791 ns
Monotonic stop: 1803750370495916 ns
Elapsed: 82.938396125 seconds
Attempt 1 remains NOT-RUN(TRANSPORT); its artifact is preserved.

## Actual answer
The local gate stopped in its plain Go test leg because TestQueryInterfaceReturnsWhileADescendantHoldsStdout could not read the fake descendant's PID file, after a 3.00-second run. This is the observed test-instrument setup failure, not an observed product timeout/WaitDelay regression. operation.log:74-77 identifies the sole named failing test and the ENOENT at iface_test.go:232; operation.log:59-60 identifies the plain leg, and :88 is the final FAIL.

The test constructs a shell double that starts sleep 30 in the background, writes $! to the PID file, then waits (host/pkgproj/iface_test.go:219-224). Its configured internal timeout is three seconds, wait delay one second (:178-182). It first checks for QueryTimeoutError (:226-229) and only then reads the PID file (:230-232). Reaching :232 proves the timeout type check passed and the call returned within the ten-second runWithin guard (:184-200). The test has no readiness handshake before invoking the bounded query. The stimulus "a live descendant holding stdout" was not established by the missing PID file. Therefore this run does not prove the intended inherited-pipe scenario worked or failed.

Production query code starts the internal deadline before running the child (host/pkgproj/iface.go:168-178), configures WaitDelay (:174), and reports QueryTimeoutError if that internal deadline expires (:189-190). Production bounds are 30 seconds and two seconds (:63-73), distinct from the test's three/one-second bounds. The local failure is consistent with the fake shell not reaching its PID write before cancellation; that is an inference, not proof of scheduling latency. Original source comments (:178-181) describe prior one-second spawn latency under full-suite load, but do not measure this invocation.

The script uses set -euo pipefail (scripts/verify_go.sh:20), runs the plain leg at :310-312, then would run the full race leg at :316 onward. The log ends at plain FAIL, so the full race leg and final success line were not reached. The earlier WARNING: DATA RACE is intentional detector-control evidence, not the failure: operation.log:14-54 shows the known-positive control, followed by build and the successful 37-test evidence manifest (:55-58).

context.json records the frozen SHA, Darwin/arm64 and successful remote checks. It establishes that CI was green on that commit, but provides no remote per-test timings or platform details. Different success/failure of a readiness-dependent test is compatible with that result; an exact platform/load explanation cannot be proved from this corpus.

## Proved / unproved
Proved: named failure, missing PID file, accepted QueryTimeoutError, return within guard, plain-leg stopping point, race-control success, passed build/evidence stages, source mechanism and absence of a test readiness handshake.
Unproved: why the PID file was absent; actual shell/descendant scheduling, filesystem or process events; descendant existence or inherited-pipe state; whether load or Darwin behavior caused it; reproducibility; production behavior under that intended stimulus; full local race-leg outcome; remote per-test behavior. No product regression is demonstrated by this red. No test runs or source changes were made.

## Isolation and commands
Profile read before timer: /Users/voightkampff/.ailang/state/mission-world-iter210/baseline.sb.
Controls independently repeated before reveal: sandboxed /bin/test -r <corpus>/question.txt rc0; sandboxed /bin/cat /Users/voightkampff/dev/sunholo-data/ailang-world/.git/config rc1 Operation not permitted; sandboxed /bin/ls /Users/voightkampff/.ailang rc1 Operation not permitted; sandboxed /bin/ls /private/tmp rc1 Operation not permitted. Each subprocess timeout was five seconds.
Every diagnostic command ran through /usr/bin/sandbox-exec -f <profile>; subprocess deadlines were five seconds for question read, ten seconds for all subsequent reads/searches. The exact Python command text and returned outputs are in the accompanying transcript. Two command-path mistakes were visible: /opt/homebrew/bin/rg nonexistent (rc71), and host/childenv/env.go nonexistent (rc1); neither result supports a diagnosis. rg was corrected to the already-known Codex vendor executable. No network, .git inspection, World DB access, tests, or external user-data diagnosis reads occurred.
Token usage unavailable.


## Evidence location and integrity

Decision-bearing summary is tracked here. Full raw transcript, original exported corpus, sandbox controls and complete gate outputs remain at the mission-namespaced runtime location `/Users/voightkampff/.ailang/state/mission-world-iter210/`; they are local artifacts, not GitHub-hosted evidence. Original source is reproducible from the full frozen SHA. No raw bulk is claimed banked in git. Hashes below identify the exact reviewed files.

| Runtime file | SHA256 |
|---|---|
| `base-verify-ail.log` | `03b2848e7fb136b11d0225707977b0b774426d13975796e9d8178b3b4401726d` |
| `base-verify-go.log` | `4c3a5b6110913ee0712415ee328d06fa08030b88d2b827d05cb980ffe4a22b44` |
| `base-verify-go-2.log` | `b7f0282fb1cf6e001f39a2ed8caeb1978e8a8da63dd24774aa1fbaa23e2688e7` |
| `pkgproj-focused.log` | `e943a9d435432cc3d5340fdbe5348d4f072ec0b3009df7a0ab20ad5a59832e69` |
| `baseline-report.md` | `9f99f8a5cbb0f06496cb714f653275b12cf48db48eaadaa07494d1f2389290ac` |
| `baseline-report-2.md` | `cd69616fa63665cfe80fa7a1ea9ed0b2f4f1c210bdbfa24a0b7f1b23ad2187bf` |
| `baseline-transcript-20260930T102633Z.txt` | `030c8854cfc6694779ba1c98e6afb9bed3e130571691498377b99b0cb29b7e21` |
| `baseline.sb` | `1513539ab111d4caab681ab6d1dba3f77d97f9591b23198315e855450076ca02` |
| `sandbox-controls.json` | `36d953e361ee050957401046aadccdb5a51f8e907eeb0aed3d0d8892ebeee20d` |
| `sandbox-controls-2.json` | `7ea144c4923f7ab9f419ac245bd47ae99ce1af19792a587ad5c39db5a8f4f527` |

Independent Sonnet judge (`claude-sonnet-5-5`) **PASS84/100**, zero blocking, solely for observational baseline banking. It independently verified source/export fidelity and both gate logs, then ran its own focused five:5/5green. Combined isolated runs0/10failures vs full gates2/2PID-instrument reds. Review limitations: filesystem fencing does not prove cognitive blinding; stop stamp untranscribed; controller-focused log lacks command/SHA header. No product/regression-exclusion/paired-walk claim follows. Full independent report: [world-iter210-baseline-evaluation.md](world-iter210-baseline-evaluation.md).

Supplemental independentSonnet review: **PASS86/100**, zero blocking, for honest ADMIN disposition only; design remains REJECTED/PARKED. [Disposition evaluation](world-iter210-disposition-evaluation.md).
