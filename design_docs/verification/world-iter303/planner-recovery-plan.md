# Iteration303 M7b recovery handoff — pending

PR361 row150 M7b: orphan bookkeeping recovery; controller/executor verification only.
Candidate `f615a9698370694664fbbc01431e4a2c549377cc`; exact base `c73c79048daf9c3a19ea05c96e793436eb45e9c5`.
Read-only local verification: HEAD equals candidate; merge-base equals exact base; three commits above base.
PR361 association is controller supplied; live PR head and CI have not been verified here.
Diff inventory: seven files, 1,432 insertions, 22 deletions:
- design_docs/planned/sprint_w-marketing-capture.json
- design_docs/planned/w-marketing-capture-sprint-plan.md
- host/store/context_roots_test.go
- scripts/marketingcheck/main.go
- scripts/marketingcheck/png.go
- scripts/marketingcheck/png_test.go
- scripts/marketingcheck/supervisor_test.go

Authority: existing design_docs/planned/w-marketing-capture.md, M7b sprint-plan section,
and sprint_w-marketing-capture.json remain authoritative. Their completion bookkeeping
and designer-recovery-audit.md introductory PASS are inherited claims, not fresh gate verdicts.
Read scope: audit first65 lines, relevant M7b plan/JSON, narrow source and identity lookups.
No runtime model identity is attested from inherited executor or reviewer labels.
Controller supplied: /tmp/world-mc-20261023/quorum303.json proceeds flagged:
Gemini pass only; GLM/Kimi unreachable; Sonnet quota absent. No revision/new quorum needed.

Recovery checklist (all verification/landing statuses pending):
- [ ] Preserve exact candidate/base; recheck local HEAD, PR head and scope before each boundary.
- [ ] Before controller gates, review full base-to-candidate diff for writes outside worktree,
  including runtime destinations, subprocesses, symlinks and HOME writes; report unexpected scope.
- [ ] Controller runs the following outside sandbox under existing isolation/load discipline:
```sh
go test -run '^TestMarketingPNGStructure$' ./scripts/marketingcheck -count=1 -timeout 90s -v
go test -run '^TestMarketingPNGContent$' ./scripts/marketingcheck -count=1 -timeout 90s -v
go test -run '^TestMarketingTwoPassInventory$' ./scripts/marketingcheck -count=1 -timeout 90s -v
go test -race ./scripts/marketingcheck -count=1 -timeout 180s
go test -run '^TestProductionGoSurface$' ./host/store -count=1 -v
tools/verify_ail.sh
tools/verify_go.sh
```
- [ ] Resolve full-gate entrypoint paths before execution: narrow tools listing found neither
  requested script; existing M7b JSON references scripts/verify_ail.sh. Record actual controller
  entrypoints and pinned interpreter; do not silently substitute or omit either full gate.
- [ ] Capture actual commands, SHA/base, environment/pins, complete stdout/stderr, exit status,
  duration, selected tests and named results. PNG names were located at png_test.go:202,307,548;
  census at host/store/context_roots_test.go:198. Require nonempty selected tests and genuine execution.
- [ ] Collect controller's already-running verify_go.sh result before further broad execution.
  Controller reports cmd/world-publish import-bearing-source and host/archive pinned-interpreter
  LDR001 reds; output is truncated behind a deprecation warning. Attribution is UNVERIFIED.
- [ ] Preserve all raw red output; controller isolates exact failing tests/commands at candidate
  and compares identical commands/pins/environment against pristine exact base. Recover complete
  LDR001 diagnostics. This planner neither diagnoses, fixes nor downgrades either red.
- [ ] For unchanged-package suite reds, retain isolated plus pristine-base controls. If these
  demonstrate inherited instrument/environment failure and exact PR/merge CI is green, apply
  existing skill semantics; do not invent a stronger rejection rule or silently waive the red.
  Record actual evidence, attribution and disposition explicitly; unresolved mandatory reds park.
- [ ] Obtain fresh independent MiniMax review of this exact candidate before landing; preserve
  exact reviewed SHA, verdict and findings. Historical reviews/design PASS do not satisfy it.
- [ ] Query SHA-pinned PR361 head and full CI job inventory per existing skill; record every
  required job's run/job IDs, SHA and conclusion. Required jobs must be present and non-skipped.
- [ ] Controller lands only through existing gates, then records merge SHA and verifies full
  merge-candidate CI inventory per skill. Merge CI remains pending until observed on that SHA.

No sandbox gate result is evidence; socket denials are UNINFORMATIVE UNDER SANDBOX.
Skips, empty selection, timeouts, panic and build errors are not success.
Mutation audit limits persist: truncated-chunk panic is not an assertion kill; redundant-condition
and filesystem-race survivors remain; unknown-chunk refusal alone does not prove acceptance after
forbidden-entry removal. Retain recorded stronger tEXt bypass evidence without claiming completeness.
No new mutation drills are requested.

Park for unresolved mandatory red; failed/missing independent judge; unexpected SHA/scope;
missing or skipped required CI; or human-bound mint/publish decisions. No unattended mint/publish.
D86 is unrelated and stays parked. Registry reuse inherits existing design: tooling-only recovery
adds no package capability; no registry operations. Real capture acceptance remains with existing design.
No new acceptance criteria, thresholds, sprint, gate weakening, production/design/charter edits,
protected writes, git writes, messages or tests are authorized for this planner.
Disposition: pending controller/executor recovery verification; no landing authorization asserted.
