# Evaluator r1 - row 150 M7c2a (commit bbf7eb84, base e694d23e)

VERDICT: PASS, 88/100. Zero blocking findings.

Rubric (approximate): conformance to prescribed interfaces and file scope 28/30; tests and gates green 25/25; correctness and robustness of the code 18/20; non-vacuity 10/15 (the 3 plan mutants and 9 of the other mutants are killed, but 13 new branches are unguarded); standards and hygiene 7/10.

## Gates (AILANG_BIN=$HOME/.pinned-ailang/ailang; uptime load 2.6-2.8)
- gofmt -l: empty. go vet ./scripts/marketingcheck: rc 0.
- TestMarketingCapturePassFakeChrome -v: PASS, 3.24 s. All five subtests (routes, pass-six, version, ready, fresh-dir) PASS.
- TestProductionGoSurface: ok. TestSubprocessSinkGate: ok.
- go test ./scripts/marketingcheck: ok, 19.8 s.
- go test -race ./scripts/marketingcheck: ok, 50.7 s (one run; the base flakes in supervisor tests did not appear).

## Conformance
- Only the three allowed files changed. The shot() argv order is untouched; the diff shows no edit to shot() or daemon().
- routes, chrome-version, ready and pass match the plan: MIN_CHROME_MAJOR=155, regex as prescribed, health and then head, urlopen timeout=10, DestinationExists, progress on stderr, stop at the first non-ok shot, one JSON line on stdout.
- The Go `routes` subcommand: exit 2 on usage (unknown flag, missing --manifest, extra argument), exit 1 on a refusal with the stderr form `marketingcheck: bad-manifest demo: ...`. Production Go stays in main.go.
- Hygiene: the test sinks are os.Create'd files. The Chrome fakes come from fakeScript. The PNG fixtures come from fixtureSet. The two accepted version strings are the real outputs. Time bounds are generous.
- Manual probes:
  - Missing routes file: `InvalidConfiguration`, rc 1.
  - `ready` on a dead port with `--startup-timeout 1`: `Timeout`, rc 1.
  - Nonexistent Chrome: `SupervisorError`, rc 1.
  - Missing `pass` args: argparse usage, rc 2.

## NON-BLOCKING findings
1. scripts/marketing_capture.py:~260-283 (capture_pass): many new refusal and bound branches have no test. Mutants SURVIVED (table below):
   - stop at first non-ok shot;
   - the six-row check;
   - the name basename check;
   - the dimension check;
   - the scheme check;
   - the remaining<=0 guard;
   - the min(SHOT_TIMEOUT_S, remaining) timeout;
   - the post-loop deadline check;
   - the base_url rstrip;
   - the per-shot OSError catch;
   - the outer InvalidConfiguration catch.
   The first of these means the plan's "stops at first non-ok shot" claim is unguarded. Reproduce: `go test -run ^TestMarketingCapturePassFakeChrome$` with `if result['verdict'] != 'ok':` changed to `if False:`. It stays green.
2. scripts/marketing_capture.py:~214 (`ready`): the `status != 200` break (the 2xx-but-not-200 case) is untested, and so is the `Chromium` regex alternative and the chrome-version `Timeout` branch (the "no timeout verdict" mutant survived; shrinking the deadline to 0.0001 s was killed only via the ok cases).
3. capture_pass creates the out-dir before validating the routes file. A bad routes file or a row count other than six leaves an empty directory, so a retry with the same --out-dir gets DestinationExists. Observed: with `--routes /nonexistent`, rc 1 and `o1` was created. This is arguably acceptable, but it is worth noting for M7c2b's shell.
4. The 'go non-refusal err' branch (the non-Refusal error message in runRoutes) was not mutated; my edit matched two sites and was not applied. UNMEASURED.

## Mutant table (test = TestMarketingCapturePassFakeChrome)
| Mutant | Edit | Result |
|---|---|---|
| M7c-scheme (plan) | emit scheme flag only when light | KILLED: pass-six "every shot explicitly pins correct scheme: workbench-dark.png" and "CheckPass ... wrong-scheme workbench-dark.png" |
| M7c2a-version-floor (plan) | MIN=154 | KILLED: version "Chrome 154 must be refused" |
| M7c2a-routes (plan) | every route "/workbench" (a first compile-failing version was redone as a valid edit) | KILLED: routes "routes must equal the capture-manifest routes: packet.png got /workbench" |
| floor `<=` | `< MIN` to `<= MIN` | KILLED: "Chrome 155 is accepted" |
| regex without Chromium | drop alternative | SURVIVED |
| chrome-version Timeout verdict renamed | | SURVIVED |
| chrome-version deadline 0.0001 s | | KILLED (version/ok cases) |
| ready skips head | | KILLED: head=0 |
| ready skips health | | KILLED: health=0 |
| ready accepts non-200 | remove status check | SURVIVED |
| ready ignores deadline | wait(1000) | KILLED (the test hung, 121 s) |
| pass: no stop on fail | | SURVIVED |
| pass: six-row check removed | | SURVIVED |
| pass: basename check removed | | SURVIVED |
| pass: dims check removed | | SURVIVED |
| pass: scheme check removed | | SURVIVED |
| pass: mkdir exist_ok | | KILLED: fresh-dir |
| pass: remaining<=0 off | | SURVIVED |
| pass: shot timeout not min(remaining) | | SURVIVED |
| pass: final deadline check removed | | SURVIVED |
| pass: progress to stdout | | KILLED: one JSON line |
| pass: url rstrip removed | | SURVIVED (the test base URL has no trailing slash) |
| pass: per-shot OSError catch altered | | SURVIVED |
| pass: outer catch narrowed | | SURVIVED |
| main dispatch guard (`args.mode in (shot,daemon)`) removed | | KILLED |
| go: missing --manifest not usage | | KILLED |
| go: ErrHelp returns 2 | | SURVIVED (benign; -h is not specified) |
| go: bad-manifest rc 3 | | KILLED |
| go: ShotRoutes reversed order | | KILLED |
| go: height forced 1100 | | KILLED |

## Restoration
cmp against the saved copies of scripts/marketing_capture.py and scripts/marketingcheck/main.go: identical. `git diff --quiet` is clean and `git -C /Users/voightkampff/dev/sunholo-data/.eval-world-iter307 status --porcelain` is empty.
