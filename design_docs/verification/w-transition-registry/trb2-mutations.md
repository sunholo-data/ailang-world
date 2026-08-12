# TR.B2 mutation transcript

All commands used `PATH=/opt/homebrew/bin:$PATH` and `GOTOOLCHAIN=go1.25.6`. Production baselines were
`bind.go` `068fc0e47a1404e53ff8e38c2e5d3d5515fbd088914e90aa8bf34e53c129746c` and
`confined.go` `c6fbdfc7178651d9ee80c5fb796f2184c5c45cb62d074b66130caabcdf962f37`.
Every mutant was copied to `/tmp/trb2_backup`, anchor-checked, built with `go build ./...`, restored
from that copy, and hash-checked. `go vet ./host/...` was rc=0 after each arm except
`MUT-PROPOSAL-FN`, where it was inadvertently delayed until the immediately following arm and was
then rc=0; this is a recorded protocol deviation, not represented as contemporaneous evidence. The scoped inverse was the
whole transitionreg package with the named killer skipped.

| Arm | Landed sha256 | Build | Required RED observation | Inverse | Restore |
|---|---|---:|---|---:|---|
| `MUT-BIND-MISSING` | `9528dfb59295` | 0 | `TestGuardedSessionRefusesUndeclaredEffect/absent_transition`: got access denial, wanted exact `TransitionAbsentError` message | 0 | match |
| `MUT-EFFECT-UNDECLARED` | `8369eb04bdda` | 0 | `.../undeclared_cost`: got nil (handler ran), wanted exact `UndeclaredEffectError` | 0 | match |
| `MUT-EFFECT-BYPASS-BROKER` | `e6df145d2fa2` | 0 | `TestGuardedSessionStillRequiresBrokerGrant/declared_but_missing_live_grant`: got nil, wanted `*DenialError denied:budget` | unsatisfiable; see below | match |
| `MUT-PROPOSAL-FN` | `a0954555c829` | 0 | `.../transition_fn_mismatch`: got nil, wanted exact mismatch message | 0 | match |
| `MUT-PROPOSAL-CAPS` | `a9513ba7896c` | 0 | `.../required_caps_mismatch`: got nil, wanted exact mismatch message | 0 | match |
| `MUT-PROPOSAL-EFFECTS` | `2fd850ac58b0` | 0 | `.../expected_effects_mismatch`: got nil, wanted exact mismatch message | 0 | match |
| `MUT-SESSION-UNION` | `2b1e2d26fb1e` | 0 | `TestTwoSessionExactOrderedSets`: session A IDs `[alpha beta gamma]`, wanted `[alpha]` | 0 | match |
| `MUT-STARTUP-CACHE` | `83d35b39053b` | 0 | `TestNextReadObservesNewHeadWithoutRestart`: second head/revision stayed at head 1/revision 1 | 0 | match |
| `MUT-CAPS-REREAD` | `c0025fbd8c7f` | 0 | `.../captured_sources_are_not_reread`: snapshot calls=2, wanted 1 | 0 | match |
| J3 `MUT-DECLARED-NAME-ONLY` | `8b3048317c9b` | 0 | `.../undeclared_scope`: got broker scope denial, wanted undeclared-effect message | 0 | match |
| J4 `MUT-DECLARED-COST-ANY` | `8369eb04bdda` | 0 | `.../undeclared_cost`: got nil, wanted undeclared-effect message | 0 | match |
| J7 `MUT-BIND-COLLAPSE-LABEL` | `a00350692acf` | 0 | scope/expired/budget label subtests got effect-name instead of exact broker labels | 0 | match |
| J8 `MUT-BIND-EMPTY-SNAPSHOT-OK` | `1490c76cfcd5` | 0 | `.../zero_snapshot`: got absent transition, wanted exact zero-snapshot refusal | 0 | match |
| J9 `MUT-PROPOSAL-INTERP` | `1550a7c85a6c` | 0 | `.../interpreter_mismatch`: got nil, wanted exact mismatch | 0 | match |
| J10 `MUT-PROPOSAL-EPOCH` | `6e382696617a` | 0 | `.../semantics_epoch_mismatch`: got nil, wanted exact mismatch | 0 | match |
| J11 `MUT-REQUEST-SWALLOW` | `78a6d4503ee1` | 0 | `.../injected_read_error`: got nil, wanted wrapped injected failure | 0 | match |
| J12 `MUT-REQUEST-ALIAS` | `253164299d9a` | 0 | `.../returned_descriptors_are_copies`: mutated schema/effect escaped | 0 | match |
| J13 `MUT-REQUEST-REORDER` | `8c0efd44c4eb` | 0 | `.../order_is_the_snapshot_order`: got `[beta alpha]`, wanted `[alpha beta]` | 0 | match |
| J14 `MUT-TARGET-BIND-SWALLOW` | `4a9fded194c3` | 0 | `.../target_bind_error`: got nil, wanted wrapped injected bind failure | 0 | match |
| `MUT-DELETE-TR-B-TEST` AC6 | `e3641e953eee` | 0 | activated inventory count was **2**, wanted 3 | n/a count gate | match (`2d06d401b5ea`) |
| `MUT-DELETE-TR-B-TEST` AC7 | `2b8406f59b31` | 0 | activated inventory count was **2**, wanted 3 | n/a count gate | match (`2d06d401b5ea`) |

## Uncertain rows and deviations

`MUT-EFFECT-BYPASS-BROKER` used the compiling direct-handler form in `broker.BoundInvoker.Request`:
look up `b.s.registry[req.Effect]`, call `Execute`, and return a zero record reference. Step 5 built.
Skipping only the named grant test remained RED because
`TestSingleRequestKeepsCapturedEpochs/captured_sources_are_not_reread` independently observes the
same broker debit call site (epoch stayed 0). Thus the prescribed inverse is **UNSATISFIABLE BY
CONSTRUCTION**; skipping both co-detectors was rc=0.

For `MUT-SESSION-UNION`, the first assertion to fire was the exact ordered-list assertion:
session A returned all three IDs instead of only `tools.alpha`. The explicit forbidden-membership
assertions are later guards and were not reached.

The plan's rule-3j cut returned zero because `bind.go` is untracked during executor work and ordinary
`git diff` omits untracked files. `git diff --no-index /dev/null host/transitionreg/bind.go` exposes
the additions. Reading that diff found J14, the previously unnamed `target.Bind` error-propagation
branch; its focused subtest and mutation are recorded above.

No arm was deferred.
