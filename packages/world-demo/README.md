# World demo evidence law

This domain source directory follows `packages/se-tools`: it contains standalone
AILANG source, not a registry release. M1a provides only the evidence-kind law in
`evidence.ail`. The host coordinator tests pin its proof and inline-test identities
and compare its results with the executable frozen kernel law. The repository
`verify_ail.sh` sweep does not include this directory.

The module is `evidence`, relative to this directory as source root. It has no
imports or effects. All three exported predicates have Z3 contracts and inline
tests. Consumers can stage the standalone file in their own source root; the
executable drift test demonstrates this with byte-identical kernel modules.

| API | Input | Result |
|---|---|---|
| `legacyGradeCode(kind)` | `CompilerOutput`, `HumanApproval`, `RecordedEffect` | `2` (ATTESTED) |
| `legacyGradeCode(kind)` | `TestReport` | `3` (TESTED), independent of pass/fail |
| `legacyGradeCode(kind)` | `AiReview`, `ProofReceipt` | `1` (CLAIMED) |
| `legacyGradeCode(kind)` | Any other string, including `authored-claim` and `PROVEN` | `0` (refusal) |
| `authoredClaimGradeCode(kind)` | `authored-claim` | `1` (CLAIMED) |
| `authoredClaimGradeCode(kind)` | Any other string, including `AiReview` | `0` (refusal) |
| `rosterKindAllowed(kind)` | One of the four kinds below | `true` |
| `rosterKindAllowed(kind)` | Any other string, including `AiReview`, `ProofReceipt`, `PROVEN` | `false` |

The four roster kinds admit references for later host validation:

| Roster kind | Required evidence | Grade after host validation |
|---|---|---|
| `test-invocation` | Real test invocation with named passing results and zero failures/skips | TESTED |
| `recorded-effect` | Consistent recorded effect with checked request/result hashes | ATTESTED |
| `authored-claim` | Demo-authored claim with its actual author displayed as “demo-authored” | CLAIMED |
| `proof-envelope` | Authenticated source-bound envelope, through ValidateProof → Resolve | PROVEN |

`authored-claim` is a NEW package-level kind. It never creates a kernel constructor
or casts commentary to `AiReview`. No AI review produced the demo commentary.
A raw `ProofReceipt` remains CLAIMED. The predicates never produce PROVEN;
`rosterKindAllowed("proof-envelope")` admits a kind name only. It neither checks
its reference nor assigns a grade. The host must reject mismatched or missing
objects as UNAVAILABLE. Roster input consists of kind/ref pairs, with no caller
grade. The later roster and proof adapters are outside M1a.

From the repository root, run this bounded quickstart with the released pin:

```sh
export AILANG_BIN=$HOME/.pinned-ailang/ailang
run_ail() { python3 -c 'import subprocess,sys; sys.exit(subprocess.run(sys.argv[1:],timeout=40).returncode)' "$@"; }
(cd packages/world-demo && run_ail "$AILANG_BIN" check evidence.ail)
(cd packages/world-demo && run_ail "$AILANG_BIN" test --format json evidence.ail)
(cd packages/world-demo && run_ail "$AILANG_BIN" ai-check -timeout 5s evidence.ail)
go test -run '^TestWorldDemoEvidencePinnedManifest$' ./host/coordinator -count=1 -timeout 120s -v
go test -run '^TestWorldDemoLegacyGradeDrift$' ./host/coordinator -count=1 -timeout 120s -v
```

Expect check success, 21 named passing inline cases, zero errors/skips, and the
three verified identities `legacyGradeCode`, `authoredClaimGradeCode`, and
`rosterKindAllowed`. Contract-generated properties run in addition to those inline
cases. Both Go tests fail when `AILANG_BIN` is unset or does not report exactly
`AILANG v0.41.0`. The manifest gate refuses empty results and tests. The drift gate
checks and runs a temporary probe over actual legacy constructors, including
both passing and failing TestReport values, and separately evaluates the authored
claim and AiReview refusal controls.
