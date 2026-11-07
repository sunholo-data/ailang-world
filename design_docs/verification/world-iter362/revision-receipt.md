# World row256 — iteration362 designer revision receipt

One r2 DESIGN revision by `codex:gpt-6.1-sol`; all three iteration361 r1 objections accepted.
No designer quorum. Controller owns commits, outside-sandbox gates and quorum with
`--author codex:gpt-6.1-sol` (OpenAI vendor benched).

Deliverables:

- `design_docs/planned/w-daemon-startup-under-load.md` (existing document revised).
- `design_docs/planned/w-daemon-startup-under-load.patch` (existing measured patch revised).
- Evidence and reproducible commands in `design_docs/verification/world-iter362/`.

Baseline SHA: `61a8d1c47b45fffc0f36137d29a5285c32da0145`.
Patch SHA-256: `90c6df45c6b4dcf09fd241d76eb87294f5b84d688eb9f2ff90d0eca1de09db16`.
Patch: **11 files, 881 insertions, 20 deletions, 1088 lines**. Four new test files,
747 lines in total; six existing Go files and QUICKSTART changed. Per-file stats and
all prototype/applied hashes: [fresh-apply.txt](fresh-apply.txt),
[prototype-manifest.json](prototype-manifest.json).
No production patch application; only independent `git archive HEAD` scratch repositories
with `git init`. Fresh `git apply --check`, `git apply`, all three required test-file
presence assertions and all 11 file equivalence assertions passed. Scratch was removed.

The flag certifies only all failed arms recorded at/after the startup deadline. Logs,
health comments, doctor, help, QUICKSTART and tests use observation-only text: timeout may
have contributed; review reasons and restart if appropriate. Both ambiguous directions and
the absence of any established causal/scheduling bound are explicit. A deterministic late
genuine-refusal fixture is excluded with its reason intact even with OutOfTime=true.
False-flag advice also asks for review and conditional correction/restart. No retries or
served-set widening. Split clocks, parent cancellation, writer release and lifecycle ownership
are preserved; the successful socket lifecycle awaits the controller.

AST census: **60 New references in test files**, **1 baseline production reference**,
`Run → startBounded(..., New)`; patched **60 test / 0 production New references**.
The 60 include benchmark/test-helper call expressions. Baseline has 622 Go files;
the test-only inventory fixture makes its scan 623. Patched scan has 626. Renamed imports,
function-value aliases and nested selector expressions have positive controls; local shadows
and unrelated fields survive. Three compiling caller mutants fail the forward guard.
Separate V20/V21/V22 rows bank commitBudget=3s, additive HealthResponse and the actual
240-byte reason cap (ellipsis included; project/file/fixed text are outside that cap).

| Gate | Baseline | Revised/final | Evidence |
|---|---|---|---|
| Build, vet, gofmt | rc 0 | rc 0 | baseline-build-vet.txt; fresh-build-vet.txt |
| Fresh independent apply + three presence guards | original measured patch applied in prototype | rc 0 | fresh-apply.txt |
| Non-socket focused race + neighbours | existing daemon/broker/store neighbours green in baseline run | **rc 0, 14 top-level PASS**, no race warning | fresh-focused-race.txt |
| All 13 new tests requested under race | new tests absent at baseline; neighbours baselined | **rc 1**, 10 top-level PASS, 2 FAIL, T5b masked by doctor panic | fresh-new-race.txt |
| Revision mutants | green 8-test unmutated control, rc 0 | **14/14 killed**, each rc 1, harness rc 0; restored hashes equal final prototype | final-mutation-run.txt; mutations.json |
| Plain verify_go in archive | **rc 1** | **rc 1**, same missing-HEAD failure | baseline-verify-go.txt; fresh-verify-go.txt |
| Race-detector positive control | expected rc 1, WARNING: DATA RACE | detector armed; this is separate from product tests | race-detector-control.txt |

Socket results are **UNINFORMATIVE UNDER SANDBOX even when green**. Denials retained in full:

- Baseline `TestDoctorDaemonVsForeignListener`: `httptest` IPv6 `[::1]:0` bind denied,
  panicking and preventing later doctor tests in that package.
- Prototype and final T1: probe arm6 control/reverse loopback binds denied and Run's
  `127.0.0.1:0` listener denied before announcement.
- Prototype and final T5: `httptest` IPv6 `[::1]:0` bind denied; T5b then ran separately
  without sockets. T2 green is also uninformative because its probe uses sockets.

`verify_go.sh` first passed its 24 tracked-binary controls, then failed `git ls-tree HEAD`
in both empty initialized archives. No history was seeded and no gate was bypassed.
The controller must baseline and run the real-checkout gate outside sandbox.

Intermediate evidence is retained and explicitly superseded where applicable: the first
caller instrument omitted the existing Run→startBounded allowlist edge (counts unchanged
when corrected), and the first new cancellation fixture canceled during exec setup,
producing a config-stage error rather than isolating the fresh registry phase. Final T7
pairs a live startup context with a canceled phase parent and proves registry cancellation
and writer reopening; the detached-parent mutant kills it. Historical V10 causal advice
is withdrawn; r1 V16/V17 are superseded by the fresh r2 receipt. `pre-*` evidence files
are earlier internal prototype states, not evidence for the final patch digest.

Every gate log records date/deadline (≤1800 seconds), uptime, Go version, exact command
and rc. PATH starts `/opt/homebrew/bin`; AILANG_BIN is `$HOME/.pinned-ailang/ailang`;
TMPDIR/GOTMPDIR untouched. No stress spinners/load experiment, external artifact seeding,
secret reads or shared HOME writes. The initial two read-only discovery calls lacked
explicit dated command deadlines; later commands and all validation gates recorded bounds.

[final-audit.txt](final-audit.txt) proves unchanged HEAD, zero tracked diff and unchanged
charter/standards/thesis/mission-log hashes. Row233 remains parked on D105; row257, logs,
dashboard, frozen machinery and production source remain untouched. No source commits.
The process census command returned **rc 3**, `sysmond service not found / Cannot get process
list`; it is UNINFORMATIVE, not a zero-process assertion. All bounded test runners had
returned; scratch removal and byte restoration were independently asserted.

Remaining concerns: successful Run/doctor socket lifecycle and related historical clock
mutants require controller reruns outside sandbox; full real-checkout verify_go and Linux CI
remain pending; the inherited template-clock mutant gap remains. Context-free archives and
scheduling have no established startup wall-time ceiling. The flag supplies no retry policy
or causal diagnosis. These are design validation limits, not passed acceptance gates.
