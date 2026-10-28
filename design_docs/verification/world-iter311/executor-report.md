# M7c4 executor report — iteration 311, row 150

Implementation is ready for controller review. M7c4 is **not COMPLETE**: outside-sandbox inventory evidence, PR/CI landing, and an independent judge PASS remain controller work. No real Chrome run or C8/C9/C8′/C12/C14 evidence is claimed.

Base: sprint/world-iter311-row150-m7c4 at 90ef2671. No git write operations were performed. Only the three authorized source files and artifacts under snap/M7c4/ were written. All command scratch was rooted in this worktree via TMPDIR; test scratch used t.TempDir(). Pins for all gates: AILANG_BIN=$HOME/.pinned-ailang/ailang and WORLD_TOOL_AILANG_BIN=$HOME/.pinned-ailang-tools/v0.54.15/ailang.

## Changed files and snapshots

| File | Diff | Final length | Full snapshot |
|---|---:|---:|---|
| scripts/capture_marketing.sh | +10 / -5 | 141 | scripts/capture_marketing.sh.txt |
| scripts/marketingcheck/capture_test.go | +38 / -6 | 497 | scripts/marketingcheck/capture_test.go.txt |
| cmd/ailang-worldd/marketing_chrome_test.go | +246 / -0 (new) | 246 | cmd/ailang-worldd/marketing_chrome_test.go.txt |

294 added lines total, below the approximately 300-line limit. All three final files compare byte-for-byte with their snapshots (cmp rc 0). git diff --check rc 0. The final gofmt check printed nothing, rc 0.

The shell reads inventory once after daemon readiness, before pass A, and once after the final pass (after the B version comparison in two-pass mode). Both modes compare through step, preserving verifier stderr and refusing with exit 1. Two-pass diagnostics record inventory equality after a successful comparison. The B-only steps are enclosed in a small if block so the common after-inventory/compare occurs exactly once in either mode (A1).

The fake daemon serves a fixed nonempty head, two paged log entries, known object bytes, health, and 404 for unknown routes/objects. urllib.parse handles query strings. Its optional head change is armed by the first from=0 log page, never by readiness head-request counts. The inventory subtest requires rc 1 and inventory-mutated in stderr. six-shots requires both inventory files and the diagnostics equality line.

The real drill checks Chrome absence/executability/version, both pin variables, and complete measured rectangles before seeding/building. Its local rectangle types mirror the JSON shape and enforce schema, all six shots, content names, and required packet regions. The local rig follows the iteration-305 empty-store probe, publishes tools/approval, mints only the test credential, runs demo seed, and explicitly closes the seed daemon before capture. Every subprocess introduced here uses os.Create file sinks under t.TempDir(); capture uses exec.CommandContext with 660 seconds. The shell's check verifies the capture; the drill additionally reads/parses the manifest and checks diagnostics.

## Sandbox and red-first record

The required initial command was run once:

`python3 -c "import socket; s=socket.socket(); s.bind(('127.0.0.1',0)); print('bind ok')"`

rc 1, exact error: `PermissionError: [Errno 1] Operation not permitted`.

The extended fake daemon and inventory subtest were written first, with the shell still unwired. Then only TestMarketingCaptureFakeChromeSixShots was run (logs/red-first.txt; rc 1). **UNINFORMATIVE UNDER SANDBOX**: the inventory subtest stopped at `capture_test.go:313: listen tcp 127.0.0.1:0: bind: operation not permitted`. It did not reach the intended assertion. Expected outside the sandbox at that point: unwired capture returns rc 0, triggering `capture must compare inventories: rc=0 stderr=...`. This expected outcome was NOT observed or claimed as a kill.

## Required gates

Commands below used the requested pins and worktree TMPDIR. Raw logs and rc files are in logs/gate-<name>.txt and logs/gate-<name>.rc.

| Gate (exact command) | rc | Observation / log name |
|---|---:|---|
| `gofmt -l scripts/marketingcheck cmd/ailang-worldd` | 0 | No output; format |
| `go vet ./cmd/ailang-worldd ./scripts/marketingcheck` | 0 | vet |
| `bash -n scripts/capture_marketing.sh` | 0 | shell |
| `go test -run '^$' ./...` | 0 | All packages compiled; compile |
| `env -u WORLD_CHROME go test -run '^TestMarketingChromeSixShots$' ./cmd/ailang-worldd -count=1 -timeout 120s -v` | 0 | Prints SKIP; chrome-unset; not product acceptance |
| `WORLD_CHROME=/nonexistent go test -run '^TestMarketingChromeSixShots$' ./cmd/ailang-worldd -count=1 -timeout 120s -v` | 1 | Expected FAIL: WORLD_CHROME must be an executable file; chrome-invalid |
| `go test -run '^TestMarketingCaptureFakeChromeSixShots$' ./scripts/marketingcheck -count=1 -timeout 300s -v` | 1 | **UNINFORMATIVE UNDER SANDBOX**: denied loopback binds; fake |
| `go test -run '^TestProductionGoSurface$|^TestProductionContextRoots$' ./host/store -count=1 -timeout 60s -v` | 0 | Final rerun: both PASS; surface |
| `go test -run '^TestSubprocessSinkGate$' ./host/verifygate -count=1 -timeout 120s` | 0 | sinks |
| `go test ./scripts/marketingcheck -count=1 -timeout 300s` | 1 | **UNINFORMATIVE UNDER SANDBOX**: ready subtest panics at httptest denied bind; package |

An initial surface run returned rc 1 because concurrent Go commands were writing generated _testmain.go files under snap/M7c4/tmp, which the surface scan correctly rejected. This was an executor scratch placement error, not a base finding. Original output is retained in logs/gate-surface-scratch-collision.txt/.rc. After those commands completed and scratch was moved into snap/M7c4/.snap (the gate's existing explicit scratch exclusion), the same unmodified gate returned rc 0. No gate or test was weakened. The package run stopped at denied binds, so it cannot adjudicate the known supervisor flake; supervisor files were untouched.

## Mutation ledger

For each production mutant: apply to scripts/capture_marketing.sh; run ONLY `go test -run '^TestMarketingCaptureFakeChromeSixShots$' ./scripts/marketingcheck -count=1 -timeout 300s -v`; restore from scripts/capture_marketing.sh.txt; run cmp against that snapshot. Tests and mutations were sequenced after the gate processes completed.

| Mutant | Applied edit | test rc | Outcome | restore cmp rc |
|---|---|---:|---|---:|
| M7c4-inventory-wiring | Delete the inventory-compare step | 1 | **UNINFORMATIVE UNDER SANDBOX** | 0 |
| before-after-pass-A (own) | Move before-inventory from readiness to immediately after pass A | 1 | **UNINFORMATIVE UNDER SANDBOX** | 0 |
| compare-rc-ignored (own) | Replace compare step with direct verifier call followed by `|| true` | 1 | **UNINFORMATIVE UNDER SANDBOX** | 0 |
| single-pass-skips-compare (own) | Run the compare step only inside `if [[ $verify_twice == yes ]]` | 1 | **UNINFORMATIVE UNDER SANDBOX** | 0 |

Exact observed failure text for **each** mutant: `capture_test.go:313: listen tcp 127.0.0.1:0: bind: operation not permitted`. No mutant reached `capture must compare inventories: rc=%d stderr=%q`; none can honestly be called killed or survived. Raw per-mutant logs are logs/mutant-<name>.txt; machine-readable ledger: logs/mutants.json. The named wiring kill remains for the controller outside the sandbox.

Review finding: the new head-change fixture detects omission/ignored comparison but does not itself assert that before-inventory precedes pass A, and the inventory refusal subtest runs in two-pass mode. Those are potential coverage gaps for the before-after-pass-A and single-pass-skips-compare mutants; no observed survival is claimed under denied binds.

## M7c4-rects-required review check (test-code mutation)

The mutation changed only the missing-WORLD_MARKETING_RECTS t.Fatal to t.Skip. It was applied, exercised, restored, and cmp against cmd/ailang-worldd/marketing_chrome_test.go.txt returned 0. This is a review check, not a production kill.

Each run used `go test -run '^TestMarketingChromeSixShots$' ./cmd/ailang-worldd -count=1 -timeout 120s -v`, pins set, WORLD_MARKETING_RECTS unset unless stated. No seeding/building or real Chrome occurred.

| WORLD_CHROME / rectangles | Original | Mutant |
|---|---|---|
| /usr/bin/true / unset | rc 1, FAIL at the earlier chrome-version subprocess (`exit status 1`) | rc 1, same earlier version failure; cannot reach rectangle guard |
| worktree fixtures/chrome-version / unset | rc 1, FAIL: `WORLD_MARKETING_RECTS is required; never guess rectangles` | rc 0, SKIP with the same text |
| worktree fixtures/chrome-version / nonexistent file | rc 1, FAIL: `WORLD_MARKETING_RECTS: open .../fixtures/does-not-exist: no such file or directory` | Not needed |

The executable fixture only answers `--version` with Google Chrome 155.0.8059.41 and refuses every other invocation. It permits the required precondition order to reach the rectangle guard without running real Chrome. /usr/bin/true is an executable but is not a valid Chrome-version instrument, so the additional fixture was necessary for a meaningful rectangles review. Exact outputs: logs/rects-*.txt and logs/rects-review.json.

## Deviations and outstanding completion conditions

- Red-first and all production mutation runs were attempted as requested but cannot reach their assertions under the sandbox bind prohibition. Controller reruns outside the sandbox are required; no acceptance is inferred from rc 1.
- The small conditional around the existing B-only steps implements controller amendment A1 while keeping one before, one after, and one compare call.
- Scratch placement was corrected after the initial surface artifact collision; original and corrected evidence are retained.
- The version-only fixture supplements the literal any-executable rectangles review because version validation must precede rectangle validation.
- No real Chrome was run, no rectangles were invented, no forbidden source files were changed, and no git write operations occurred. PR landing/CI and the independent judge remain outstanding; the controller owns those and real-Chrome rehearsal R1.
