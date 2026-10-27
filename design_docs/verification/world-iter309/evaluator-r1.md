# Evaluator r1 - World iter 309, row 150 M7c2b (a96820c5..76216725)

SCORE: 83/100  VERDICT: PASS   (blocking findings: none)

Breakdown: plan conformance 27/30, tests and gates 23/25, mutation/non-vacuity 14/20, premise and risk 9/10, hygiene 10/15.

## Gates (own worktree, AILANG_BIN=$HOME/.pinned-ailang/ailang)
| gate | rc / result |
|---|---|
| gofmt -l scripts/marketingcheck | rc 0, no output |
| go vet ./scripts/marketingcheck | rc 0 |
| bash -n capture_marketing.sh | rc 0 |
| python ast.parse marketing_capture.py | rc 0 |
| grep -cE 'v0\.[0-9]+\.[0-9]+' capture_marketing.sh | prints 0 (grep rc 1 is just "no match") |
| go test -run '^$' ./... | rc 0 |
| FakeChromeSixShots -v | PASS 9.8s; 9 subtests + daemon-readiness-config all PASS |
| FakeDaemonCleanup -v | PASS 5.7s, all 13 subtests incl. shell-hup |
| go test ./scripts/marketingcheck | ok 33.0s |
| go test -race ./scripts/marketingcheck | ok 98.4s |
| TestSubprocessSinkGate; TestProductionGoSurface/ContextRoots | ok, ok |
No load flakes seen. Worktree restored: `git diff --exit-code` clean, `git status --short` empty.

## Plan conformance
Met: three new flags + help text, all required; fixed flow order (version A, daemon with --ready-url, routes, pass A, early rc 0 without --verify-twice, pass B, version B, compare, diagnostics, check); identities and tool-bin forwarded, no version literal; traps and --unit-daemon untouched; Python --ready-url is mutually exclusive with --ready-file and required one-of. All plan subtests present, plus bind-refused (A2), single-pass, daemon-exits, daemon-readiness-config.
Base red-first: restoring a96820c5 production files with the new test file reds 9/9 shell subtests (daemon-readiness-config legitimately passes at base because --ready-file was already required, rc 2).
Minor deviations: the test file grew ~240 lines and the diff is 308 lines against the plan's ~250-line stop rule (almost all test; the executor note says stop and report, not shortcut). The existing M7c2a pass test's fake Chrome was refactored into fakeCaptureChrome and its args filename changed to "pass-<name>.args" (necessary to separate pass A/B; the M7c2a test still passes). The shell prints no help for the earlier-documented flags beyond what existed (unchanged). Help line for --bind has a stray extra space (cosmetic).

## Controller amendments
A1 (lifetime 660 s): sound and sufficient. In marketing_capture.py the lifetime clock starts after readiness (line ~182), so the daemon needs only pass A + pass B (<= 600 s worst case); 660 s covers that with margin, and version A/B and routes/check do not need the daemon. It also bounds an orphaned supervisor if the shell is SIGKILLed. The "version A + ... + version B" comment arithmetic is loose (versions run outside the daemon window) but harmless. The value is derived from the production constants via runpy, no new deadlines invented. Weakness: nothing tests it (see mutants). If the lifetime lapses the supervisor exits rc 0 "ok" silently and the shell then fails later with a confusing step error.
A2 (bind-refused): sound. Refusal precedes any daemon spawn; test asserts rc != 0, no argv file written, stderr names the bind.
Exit paths: EXIT/INT/TERM/HUP traps TERM the supervisor which owns the daemon group; the fail() paths, early exit 0, and daemon-exits all pass through cleanup (cleanup subtest asserts the daemon PID is dead). Only SIGKILL of the shell leaves the supervisor, bounded by the 660 s lifetime. Writes: only under --out (daemon-started.json, daemon-verdict.json, routes.json, passes, diagnostics, manifest) plus Chrome's pre-existing temp user-data-dir from M7a. Missing/unknown flag refusal is sensible (exit 1 with "capture: missing --x"/"unknown flag"), though exit code 1 rather than the verifier's usage code 2.

## Plan mutants (all killed with the named assertion)
| mutant | killed by |
|---|---|
| M7c2b-version-changed (comparison `|| 1`) | version-changed: "a Chrome version change during capture must refuse: rc=0" |
| M7c-scheme re-drill (pin forced to '1' / '0') | scheme: "every shot explicitly pins correct scheme: pass-a/workbench-dark.png argv=..." (also six-shots/serve-argv/cleanup via rc 1) |
| M7c2b-identities (drop servesEntry) | serve-argv: "both seed identities must be configured: [...]" |

## Own mutants (K=killed, S=survived)
| mutant | result |
|---|---|
| scheme flag dropped entirely | K scheme |
| bind regex widened to any IP | K bind-refused |
| early exit removed (`:`) | K single-pass ("stops after pass A ... unexpected=[pass-b capture-manifest.json diagnostics.txt]") |
| diagnostics version label / bind value changed | K six-shots (diagnostics must record version and bind) |
| check step replaced with true | K six-shots (capture-manifest.json missing) |
| step() ignores failures | K old-chrome |
| started-file wait loop disabled | K daemon-exits |
| ready() process.poll NonzeroExit branch removed | K daemon-exits (verdict Interrupted, 30 s) |
| `and not stop.is_set()` guard dropped to `ok` only | K nothing in full package (S) |
| `result ok` dropped (only not stop) | K daemon-exits (NonzeroExit during lifecycle) |
| ready() stop-Interrupted branch removed | S (full package) |
| poll thread stop check removed | S (full package) |
| --ready-url ignored | K six-shots/serve-argv/version-changed |
| mutually exclusive group made optional | K daemon-readiness-config |
| `lifetime` flag dropped (default 300 s) | S |
| lifetime derivation shortened to PASS only | S |
| port > 0 check removed / port <= 65535 widened | S / S |
| kill -0 supervisor check in readiness wait removed | S |
| missing-flag check removed | S |
| unknown-flag refusal replaced by shift | S |
| `routes ... || true` | S (set -e/step semantics make it equivalent-ish; routes.json errors still fail later) |
| readyfile-guard rewrite | S (mutant was near-equivalent) |

## Findings
BLOCKING: none.
NON-BLOCKING:
1. A1 is untested: dropping --lifetime or shrinking the derivation survives the whole package (unit test would need a lifetime-visible hook, e.g. a fake daemon logging the supervisor argv or a short-constant override). The premise is correct by reading, but the guarantee is unguarded.
2. ready()'s new stop/Interrupted behaviour (stop param branch, poll-thread stop check, `not stop.is_set()` guard) survives; no test signals the supervisor during HTTP readiness. Only the process-exit branch is covered.
3. Shell argument validation (missing flag, unknown flag, port 0 / >65535, leading-zero port) has no subtest; only the non-loopback address is covered. Exit code for usage errors is 1, not the 2 used elsewhere.
4. `kill -0 $supervisor_pid` early-exit in the readiness wait is unguarded (daemon-exits is killed by the Python NonzeroExit verdict first, so the shell bound is redundant in tests).
5. Diff size 308 lines exceeds plan's ~250-line stop rule; this is almost entirely test code and the plan said to report rather than shortcut, so no penalty beyond hygiene.
6. Cosmetic: extra space in the --bind help line; the lifetime comment arithmetic (includes version A/B, which run outside the daemon window).
