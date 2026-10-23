# World iteration303 executor recovery verification

Model requested: explicit `gpt-6.1-sol`; hidden runtime identity is not attested. Resolver declared: provider-pin. Supervisor owns actual token accounting.
Candidate/HEAD: `f615a9698370694664fbbc01431e4a2c549377cc`.
Exact base and independently observed merge-base: `c73c79048daf9c3a19ea05c96e793436eb45e9c5`.
Fresh local verification only; no implementation change, independent judging or landing verdict.

Diff inventory: seven regular files, 1,432 insertions / 22 deletions; complete diff reviewed and saved as `executor-recovery-evidence/candidate.diff`.
- `design_docs/planned/sprint_w-marketing-capture.json`
- `design_docs/planned/w-marketing-capture-sprint-plan.md`
- `host/store/context_roots_test.go`
- `scripts/marketingcheck/main.go`
- `scripts/marketingcheck/png.go`
- `scripts/marketingcheck/png_test.go`
- `scripts/marketingcheck/supervisor_test.go`

Scope: three verifier additions; supervisor test package rename; two census allow-list entries; M7b completion bookkeeping. JSON Unicode escape normalization also spans earlier/later milestones; decoded strings unchanged. No unexpected file/path/mode scope.
Boundary review: CLI writes stdout or caller-selected absent `--out` (0600, exclusive create). Output may be outside checkout, including via symlinked parents; no worktree confinement is implemented. Existing final symlink cannot be overwritten. PNG verifier otherwise reads input paths.
Tests write/remove only TempDir fixtures, create an internal temporary symlink, build CLI into TempDir and run it with inherited environment. Census reads repository. Seven changed paths resolve inside checkout. No new fixed HOME destination or HOME mutation; Go cache/module cache/temp writes are expected. Unselected inherited supervisor tests invoke Python/shell and preserve HOME.
Filesystem read/check races remain a limit; observed symlink refusal is not race-proof containment. Full boundary details: `executor-recovery-evidence/boundary-review.json`.

Initial and final SHA256s (all tracked changed files; each unchanged):
- `design_docs/planned/sprint_w-marketing-capture.json`: `bbeb187645041d25dd58f6c91a8eb539e284669bdf586ba25f6e439658dce593`
- `design_docs/planned/w-marketing-capture-sprint-plan.md`: `e78cc584a0ef6ac78b4385fe3d37d35ac339353c746539605ce5de4ef10ee21f`
- `host/store/context_roots_test.go`: `86527132157d2d2e1ffd7165ac6c3fc50f5f387b8058c2af56e9c0a53320a491`
- `scripts/marketingcheck/main.go`: `9b5cdd3a4cb25a657d684eac4067624b4ebdc5d67a4334342438cce896015c84`
- `scripts/marketingcheck/png.go`: `ea9b43df9b5e6f70757aff1798b28a258702271864a05f823ea88a760db5145e`
- `scripts/marketingcheck/png_test.go`: `cce1893446ff39ecb8366d4705f118fe143688539aa6e26602a1e41ac81e1a2c`
- `scripts/marketingcheck/supervisor_test.go`: `7cb4a1ba63ae13fb3fb0b09b9a1f25d466c1b45305da921aa3eb35c37bf9b447`

Exact fresh focused commands (complete combined stdout/stderr in named baseline logs; wall seconds):
- `go test -run '^TestMarketingPNGStructure$' ./scripts/marketingcheck -count=1 -timeout 90s -v` — rc 0, 0.614s; `--- PASS: TestMarketingPNGStructure (0.07s)`; log `baseline-TestMarketingPNGStructure.log`.
- `go test -run '^TestMarketingPNGContent$' ./scripts/marketingcheck -count=1 -timeout 90s -v` — rc 0, 0.668s; `--- PASS: TestMarketingPNGContent (0.27s)`; log `baseline-TestMarketingPNGContent.log`.
- `go test -run '^TestMarketingTwoPassInventory$' ./scripts/marketingcheck -count=1 -timeout 90s -v` — rc 0, 1.511s; `--- PASS: TestMarketingTwoPassInventory (1.12s)`; log `baseline-TestMarketingTwoPassInventory.log`.
- `go test -run '^TestProductionGoSurface$' ./host/store -count=1 -v` — rc 0, 1.213s; `--- PASS: TestProductionGoSurface (0.04s)`; log `baseline-TestProductionGoSurface.log`.

Mutation ledger: only the eight requested IDs, sequential; recipes imported solely as edits from historical JSON. Every replacement uniquely matched, changed SHA differs from original, build rc0, named assertion rc1 without panic/skip/build-red/timeout/empty selection. Each finally restored original bytes then reran exact named test with top-level PASS.
Build artifacts/byte backup are outside checkout in Go cache (path in `executor-recovery-evidence/scratch.json`). No `.go` backup in checkout.
All restored SHA256 values below equal `ea9b43df9b5e6f70757aff1798b28a258702271864a05f823ea88a760db5145e` (S). Logs are `<id>-build.log`, `<id>-red.log`, `<id>-originalgreen.log`; full per-command rc/walltime/named evidence and recipes are in `results.json`.

| Mutant | Changed SHA256 | Build rc / s | Assertion red (rc1 / s) | Restored SHA / originalgreen rc / s |
|---|---|---|---|---|
| M7b-uniform | `428f9d5e4bd38bbc20d15c07874672288a8fe3e4315ee4d13fe6612dea6fc51d` | 0 / 0.248 | uniform PNG rejected by color guard / 0.811 | S / 0 / 0.604 |
| M7b-missing | `e99e0700ad305abd2af491760ca341d5df1aa701e13520c2221297062fbed106` | 0 / 0.228 | missing shot rejected / 1.684 | S / 0 / 1.326 |
| M7b-altered | `81ca1a00da2a3142241dbc41d5dd37e8b043e6f13ca32ab58766654f4368e6b9` | 0 / 0.198 | altered bytes rejected / 1.718 | S / 0 / 1.339 |
| M7b-forbidden | `6aa7d277b041bf6364059b0cf5d3f6c207b1540dad4a34d8bfb81bd4c5cb0605` | 0 / 0.176 | forbidden tEXt rejected / 0.599 | S / 0 / 0.419 |
| M7b-crc | `ba5ae08e8669e7daf535fd3854af16139859bac8f5fcb95af475aada48ec45a9` | 0 / 0.331 | bad-CRC typed refusal / 0.593 | S / 0 / 0.417 |
| M7b-tail | `a85426cb9184e464b67860f143a4dcf21191446a0a3d604f39aa724342f9e3e0` | 0 / 0.207 | one trailing byte rejected / 0.621 | S / 0 / 0.489 |
| M7b-cap | `ece5230071127da88937510592a8dd1cc3f2712f79b5d77b4c481f8a85efe390` | 0 / 0.234 | cap+1 rejected / 0.613 | S / 0 / 0.442 |
| M7b-total | `cd2b94312c813730d9198ea8feca0b818bdd99df71c289ef866590dd8735b171` | 0 / 0.223 | aggregate > DirCap rejected / 0.659 | S / 0 / 0.368 |

Fresh selection: 8/8 assertion kills, 8/8 successful builds, 8/8 SHA-identical restores and originalgreen reruns. No fresh survivor among these eight.
Assertion reach limits: uniform red occurs at direct CheckColors assertion (png_test.go:332), before subsequent CheckContent assertion; missing red occurs at CheckPass assertion (:583), before subsequent CheckTwoPass assertion. This does not freshly prove those later assertion paths executed under mutation. Altered reaches both ancillary and pixel pair mismatch assertions and CLI mismatch.
Historical residual survivors persist, not freshly drilled or counted as kills: truncated-chunk panic; redundant-condition masking; filesystem race. Unknown-chunk refusal alone is not forbidden-chunk bypass proof; this drill explicitly bypassed both tEXt policy guards and failed the tEXt assertion.
No full suites or race command run. Controller owns `bash scripts/verify_ail.sh`, `bash scripts/verify_go.sh` and full non-sandbox gates; controller reports full Go rc1. No controller/base-control evidence is relabelled as fresh executor execution.
All results here are scoped sandbox observations, not controller Gate4 evidence. Socket denials are UNINFORMATIVE UNDER SANDBOX; none appeared in these selected runs. Real capture acceptance and independent MiniMax judging remain outside this task.
Final boundary: HEAD/base still exact; `git diff --exit-code` rc0; all seven initial/final hashes identical. Only report/evidence untracked additions. No git writes, test edits, permanent source changes, protected writes, browser, messages, publishing or agents.
Runner elapsed 18.565s; complete work stayed within 600s cap. All spawned test/build processes completed and runner exited rc0; nothing left active.
Concrete deliverable: this report plus complete raw logs, identity, SHA boundaries and mutation receipts in `executor-recovery-evidence/`.

Supervisor terminal receipt: probe rc0 / 9.376s; real run rc0 / 234.899s; neither capped. Actual CLI-reported usage: probe input 17392, cached 7168, output 5; real input 597484, cached 534656, output 5670, reasoning 397. Full supervisor receipt: `/tmp/world-mc-20261023/executor-303-summary.json`.
