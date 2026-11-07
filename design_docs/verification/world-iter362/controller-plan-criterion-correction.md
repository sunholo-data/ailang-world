# Controller criterion correction — iteration 362

Authority: authoritative mission-control verification-protocol.md §3e(a): a gate already red at base measures the repo, so record the result as part of the criterion and explicitly fix the broken acceptance criterion. This is no design-direction override or waiver of remote CI.

The provider-pinned plan was preserved as authored, including its generic active hold. The controller repairs only the infeasible local full-Go rc-0 criterion: run the identical full command on the candidate with the identical inherited environment; require exactly the same two named failures, at the same diagnostic-truncation layer, and no additional failures. A new/different failure or incomplete run blocks. Do not change the pin, warnings, environment, test policy or production code outside the frozen patch. All original local focused gates remain mandatory and passed at pristine baseline (10 named neighbours; 13 new tests absent expected gaps).

Baseline: pristine61a8d1c47b45fffc0f36137d29a5285c32da0145, full plain verify_go rc1, only TestTransitionsVerbRefusesImportBearingSourceUnderPinnedInterpreter and TestCheckSourceUnderPinnedInterpreter. The expected LDR001 is absent from the existing400byte retained diagnostic after a trace-exporter warning. Their code/tests are outside this patch. Full AILANG rc0. No unmeasured platform-wide diagnosis is claimed.

Require current PR and merge Linux CI: go host build + test gate, go race gate (go test -race, own job), ailang-code verify gate, all present/completed/success. Linux CI supplies the full-suite pass; the local comparison supplies regression attribution. A local equal-red pair alone never means acceptance or LANDED.

The independent MiniMax evaluator must specifically adjudicate this transparent criterion repair and the actual diff, with power to reject. All AC1–AC9 remain pending until their proper evidence closes. The authorized state copy records this limited correction; the original planner state/plan/handoff and their hashes are preserved for audit. No human decision has been created or resolved.
