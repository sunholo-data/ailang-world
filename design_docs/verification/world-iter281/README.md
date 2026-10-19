# Iteration 281: measured Git/mirror Unicode discrepancy

**Evidence only. Measurement succeeded; equivalence failed; M1c remains incomplete.**
This is a standalone characterization, not a permanent failing CI acceptance test or
completion of `TestProtectedGitPathspecEquivalence`.

Reproduce from this checkout with Git, Python 3 and the project's Go toolchain/dependencies
already available:

```sh
export PATH=/opt/homebrew/bin:$PATH
python3 design_docs/verification/world-iter281/repro.py
```

`repro.py` uses Python's standard library. It removes stale `result.json` before
measurement and writes a new result only after instrument controls and mirror tests
succeed. Exit 0 means valid measurement, not equivalence. Missing exact unfiltered
paths/status, missing ASCII selected controls, failed commands, timeouts or absent
named mirror PASS output cause failure. `result.json` contains actual argv, cwd, rc,
stderr/stdout as UTF-8 plus lossless base64 and NUL-field base64, Git version,
platform and checkout base. Recorded temporary cwd paths are removed after the run.

| Pattern | Path | Unfiltered ls-files/status | Selected ls-files/status | Mirror |
|---|---|---|---|---|
| `tests/**` | `tests/x.py` | present / modified | present / modified | expected true (ASCII control) |
| `tests/**` | `teſts/x.py` | present / modified | empty / empty | named test PASS, expected true |
| `*k_test.go` | `pkg/k_test.go` | present / modified | present / modified | expected true (ASCII control) |
| `*k_test.go` | `pkg/K_test.go` | present / modified | empty / empty | named test PASS, expected true |

Each row has its own disposable repository, avoiding APFS aliases between ASCII and
Unicode spellings. Both selection commands use `-z`; same-repository unfiltered
controls prove the path exists and its change is visible. Git config is isolated
with fake HOME/XDG config, no system/global config, cleared inherited Git variables,
`core.ignorecase=false`, `core.fsmonitor=false`, and commit hooks disabled. All fixture
repositories and temporary files live under the worktree and are removed. No checkout
index/HEAD writes occur. Existing Go toolchain/module/build caches are retained per
controller instruction; no real HOME/state guard is seeded. Go module downloads are
disabled (`GOPROXY=off`), with checksum verification retained. Git subprocesses have a
10-second timeout; the controller explicitly allowed the focused Go test 120 seconds
(internal test timeout 110 seconds), superseding the original blanket 10-second limit.

**Source literal used:** the two directory pathspecs are
`:(top,glob,icase)tests` and `:(top,glob,icase)tests/**`; the basename pathspec is
`:(top,glob,icase)**/*k_test.go`. These were separately checked by reading
`host/broker/exec_protected.go:327–337`: directory rendering appends bare D and D/**;
basename rendering prefixes `:(top,glob,icase)**/`. The script does not parse source
or call the production renderer. Existing `TestProtectedRenderers/git` independently
pins that construction, but was not run by this script. Existing fixture rows at
`host/broker/exec_protected_test.go:196–197` assert true; the recorded exact Go command
runs those existing long-s/Kelvin rows, including pattern parsing and mirror matching.

This establishes a counterexample for Git 2.54.0 on the recorded macOS host. It does
not measure every Git/platform/locale, the full M1c matrix, r3 root/subdirectory or
replacement mutation controls, other engines, full broker tests, CI or integrity
correctness. Broad broker gates remain unmeasured due to PF6 real-HOME writes.
No production, tests, design, sprint status, driver or skill files were changed.
Preserve the Unicode semantics and route detection redesign through bounded design
revision/review/replanning; do not weaken fixtures or call this equivalence green.
