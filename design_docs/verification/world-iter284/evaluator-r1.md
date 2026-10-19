# Judge report, iter 284, row 156 M3a+M3b (38ed7e7, 7f6c762) — evaluator r1

Score: 91/100 — PASS. Correctness 38/40, tests/mutation 23/25, scope 15/15, quality 8/10, honesty 7/10... (sum 91 rounded: 38+23+15+8+7=91)

## Blocking findings
None.

## Verified
- Focused tests all PASS on darwin with real srt 0.0.78: TestExecSettings* (ProtectedRendered, GoldenDarwin, GoldenLinux, NoProtectedIsByteIdentical x4), TestExecProtectedSettingsThroughHandler (true/false), TestExecSrtProtectedWrites 18/18 (parent-dir-rename logs the detection residual). `go vet ./...` rc 0; `go test -race -count=1 ./host/broker/...` ok (382 s).
- Scope: 5 files only, all within the allowlists. LOC M3a 146 (<=210), M3b 187 (<=235). No .ail, no testdata/golden change, no realHome/UserHomeDir/HOME use; new srt test uses canonicalTempDir home + worktree fixtures; a corrupt configured pin fails (t.Fatalf) before the bwrap/OS skip.
- Production change is minimal (2 lines): append after deny names and srt defaults; handler forwards cfg.Profile.Protected.

## Mutation table (one at a time, restored byte-identical; scratch removed; git status clean)
| Mutant | Killer | Observed red | Restored |
|---|---|---|---|
| MUT-NO-SRT-PROTECTED (delete append) | TestExecSettingsProtectedRendered | FAIL darwin+linux subtests | green |
| MUT-PROTECTED-BEFORE-GITHOOKS (append before deny names) | TestExecSettingsGoldenLinux | FAIL | green |
| MUT-LINUX-GLOB (glob on linux) | TestExecSettingsGoldenLinux | FAIL | green |
| MUT-SRT-EMPTY-BYTES (extra entry when Protected empty) | TestExecSettingsNoProtectedIsByteIdentical | FAIL all 4 subtests | green |
| MUT-HANDLER-PROTECTED-DROPPED (drop Protected: at handlers_exec.go) | TestExecProtectedSettingsThroughHandler | FAIL /true only | green |
| Extra: darwin glob `<wt>/*/G` | GoldenDarwin + ProtectedRendered (fast); real srt test also FAIL create-nested, create-root | red | green |
| Extra: darwin glob missing `**` (`<wt>/G`) | GoldenDarwin + ProtectedRendered | FAIL | green |
| Extra: dir pattern keeps `/**` suffix | GoldenLinux | FAIL | green |
| Extra: append before srt defaults | GoldenDarwin | FAIL | green |
| Extra: delete append, real srt | TestExecSrtProtectedWrites | FAIL modify, create-nested, create-root, ... (many) | green |
Each killer was a real assertion failure, not compile error/skip/timeout.

## Non-blocking findings
1. TestExecProtectedSettingsThroughHandler derives `want` from f.profile.Protected and never asserts it is non-empty for the true arm. If the fixture/profile parsing silently dropped "protected", both arms would pass vacuously. Mitigated: the profile goes through ParseExecProfile and the other tests pin the non-empty golden. Suggest `len(f.profile.Protected) > 0` for protected=true and `==0` for false.
2. The srt rows share one worktree sequentially; I found no masking (refused arms leave bytes unchanged; the only mutating protected-dir rows are refused on both OSes; parent-dir-rename runs late and only moves parent/), but order-independence is by inspection only.
3. Darwin PKG/pkg share one directory (case-insensitive FS); the case-variant arm uses distinct files (PKG/A_TEST.GO vs pkg/B_TEST.GO), so it is not vacuous, but it does not test a case-folding difference of the directory itself.
4. The "allowed" branch is strong (exact post-state hashes and before!=after), not just "bytes changed".

## UNMEASURED-LOCAL
Linux behaviour (bubblewrap, errnos, directory-only denyWrite, linuxAllowed flags for create-nested/create-root/case-*/basename-* rows, tests/modify.py modify row) — reasoned from design V14 only; CI is the measurement. Settings-golden Linux bytes are measured (pure rendering).
