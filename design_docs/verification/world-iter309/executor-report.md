# M7c2b executor report

Implemented in the authorized worktree only. No git writes, no edits to supervisor_test.go, no new production Go files, no agents, and no product Chrome evidence. Socket-dependent acceptance remains for the controller outside the sandbox.

## Files and size

| File | Diff | Final lines |
| --- | --- | --- |
| scripts/capture_marketing.sh | +54 / -3 | 136 |
| scripts/marketing_capture.py | +33 / -20 | 367 |
| scripts/marketingcheck/capture_test.go | +190 / -20 | 434 |

Shell plus test: 244 added lines, 23 removed, net +221 (267 total diff lines, approximately the 250-line bound). All seven required subtests are present; the existing M7c2a fake Chrome is extracted into a shared helper. Both/neither readiness-source usage errors are also tested. Full final source snapshots preserve relative paths under scripts/ here; base copies are under base/scripts/. All three final cmp checks return 0 (restoration.json).

## Behavior and amendments

The shell implements the fixed nine-step order, both hard-coded proof identities, operator-supplied tool path, optional single-pass early success, version consistency refusal, diagnostics, and verifier invocation. Failed verifier output remains on disk. The only added capture flags are --worldd-bin, --bind and --marketingcheck-bin; help documents them. There is no AILANG version string. EXIT/INT/TERM/HUP traps, shot passthrough, and the full --unit-daemon block are byte-identical to base, verified by direct substring equality.

A1: --lifetime is computed from the imported Python constants as 2*SHOT_TIMEOUT_S + 2*PASS_TIMEOUT_S = 660 seconds, with the requested derivation comment. Shell readiness uses STARTUP_TIMEOUT_S from the same source. No new production deadline is introduced. The specified lifetime is conservative for the named operations because version A precedes daemon startup; routes/check are not separately deadline-bounded by this milestone's contract. No substitute lifetime was chosen.

A2: non-127.0.0.1 addresses and invalid ports are refused before Chrome or daemon startup; stderr names the supplied bind. bind-refused checks both the stderr and absence of the daemon argv file. The mutant below proves that pre-start assertion.

Python daemon mode requires exactly one readiness source through argparse's required mutually exclusive group. HTTP mode calls ready(), whose existing health-then-head 200 poll is extended with process-exit and signal cancellation; started-file behavior remains unchanged. A direct socket-free probe verified health/head order, NonzeroExit for an exited process, and Interrupted for a signalled event (rc 0).

## Sandbox and red-first

The prescribed socket probe returned rc 1: PermissionError: [Errno 1] Operation not permitted. Socket-dependent results below are UNINFORMATIVE UNDER SANDBOX, never acceptance passes or failures.

Tests were written and run before production edits (red-first.log, rc 1). old-chrome failed with `old Chrome must refuse before any shot: rc=1 stderr="six-shot pass is not implemented yet (M7c); capture refused\n"`. bind-refused failed with `non-loopback bind must refuse before daemon starts: rc=1 stderr="six-shot pass is not implemented yet (M7c); capture refused\n"` and argv absent. six-shots, serve-argv, version-changed, cleanup and scheme stopped at `listen tcp 127.0.0.1:0: bind: operation not permitted`; their base behavior is UNINFORMATIVE UNDER SANDBOX. Inspection confirms base's unconditional refusal, but the socket tests could not reach it. daemon-readiness-config returned usage errors at base as well (missing required --ready-file / unknown --ready-url); it is a compatibility check, not a claimed red-first kill.

After implementation and restoration, old-chrome, bind-refused and daemon-readiness-config pass together (rc 0, gate-restored-refusals.log).

## Gates

Exact commands and initial rc values are in gates.json; raw logs are gate-*.log.

| Gate | rc | Result |
| --- | --- | --- |
| gofmt -l scripts/marketingcheck | 0 | Empty output |
| go vet ./scripts/marketingcheck | 0 | Green |
| python3 -c "import ast; ast.parse(open('scripts/marketing_capture.py').read())" | 0 | Green |
| bash -n scripts/capture_marketing.sh | 0 | Green |
| grep -cE 'v0\.[0-9]+\.[0-9]+' scripts/capture_marketing.sh | 1 | Prints 0; grep's expected no-match rc. Positive control on orchestration document prints 23, rc 0 |
| go test -run '^$' ./... | initially 1, rerun 0 | Initial source snapshots were discovered as incomplete Go test packages; snapshot-only module boundary fixes discovery; gate-compile-restored.log is green |
| go test -run '^TestMarketingCaptureFakeChromeSixShots$' ./scripts/marketingcheck -count=1 -timeout 300s -v | 1 | UNINFORMATIVE UNDER SANDBOX: five denied net.Listen calls; old-chrome, bind-refused and readiness-config pass |
| go test -run '^TestMarketingSupervisorFakeDaemonCleanup$' ./scripts/marketingcheck -count=1 -timeout 120s -v | 0 | All subtests pass, including shell-hup |
| go test ./scripts/marketingcheck -count=1 -timeout 300s | 1 | UNINFORMATIVE UNDER SANDBOX: pre-existing ready subtest panics when httptest cannot listen |
| go test -race ./scripts/marketingcheck -count=1 -timeout 300s | 1 | UNINFORMATIVE UNDER SANDBOX: same httptest bind denial |
| go test -run '^TestSubprocessSinkGate$' ./host/verifygate -count=1 -timeout 120s | 0 | File sink gate green |

No carried shell-hup flake observed in the isolated gate. Package/race runs stop at the bind denial, so later carried fixtures were not fully exercised.

## Mutants

Each valid production mutation was applied alone, only its named subtest was run, production was restored, and cmp was run against the final source snapshot. Commands/rc/cmp are in mutants.json; logs use the mutant name. No concurrent gate ran during mutation.

| Mutant | Named subtest | rc | Status | cmp rc |
| --- | --- | --- | --- | --- |
| M7c2b-version-changed: replace version equality refusal with true | version-changed | 1 | UNINFORMATIVE UNDER SANDBOX; neither killed nor survived established | 0 |
| M7c-scheme re-drill: emit scheme argv only for light shots | scheme | 1 | UNINFORMATIVE UNDER SANDBOX; neither killed nor survived established | 0 |
| M7c2b-identities: remove servesEntry requirement | serve-argv | 1 | UNINFORMATIVE UNDER SANDBOX; neither killed nor survived established | 0 |
| A2-drop-loopback: remove address/port validation | bind-refused | 1 | Killed by pre-start assertion: daemon argv exists | 0 |

The first three all stop at this exact diagnostic, before their named assertions:

    capture_test.go:307: listen tcp 127.0.0.1:0: bind: operation not permitted

Their intended assertions (`a Chrome version change during capture must refuse`, `every shot explicitly pins correct scheme`, `both seed identities must be configured`) were not reached. No named-assertion failure is claimed for them.

A2's exact failing assertion:

    capture_test.go:335: non-loopback bind must refuse before daemon starts: rc=1 stderr="capture: daemon readiness failed within shell bound\n" argv stat=<nil>

This kill does not depend on socket success: the forbidden daemon was demonstrably launched and wrote its argv before its bind was denied. After restoration, bind-refused passes and all final source snapshots compare equal.

## Deviations / limitations

- Added snap/M7c2b/go.mod solely to exclude full .go evidence copies from repository package discovery; this retains the requested exact snapshot paths without modifying repository build configuration. The initial compile-fence finding and corrected rerun are retained.
- The three planned mutant verdicts and socket-dependent gates require the controller's unsandboxed rerun; tests were not weakened to avoid the restriction. A2, old-Chrome refusal, argument exclusivity, syntax/static gates and cleanup regression have local evidence.
- No real-Chrome C8/C9 evidence is claimed. The milestone is implemented, with environment-limited verification still outstanding.
