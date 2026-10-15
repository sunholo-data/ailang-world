# World demo evidence law

This domain source directory follows `packages/se-tools`: it contains standalone
AILANG source, not a registry release. The evidence-kind law is in `evidence.ail`; `approve.ail` implements the
standalone convention-v2 `demo-approve` transition. The host coordinator tests pin its proof and inline-test identities
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

## Pending demo approval

`transitions.json` declares `demo-approve`: access is
`Human.Approve@worktree` at cost 0; its single declared effect costs 1.
The published capsule stages only `approve.ail`, which imports only std modules.
It repeats the roster-kind predicate; the executable drift probe compares both
sources. The transition itself is pure (`! {}`): its returned JSON describes the
host effect. The runner prints the returned string; no IO effect row is needed
inside this source.

From the repository root, this complete payload runs the plan directly under the
released pin. The temporary argument file contains ONE JSON string, matching the
existing coordinator calling convention. These sample references are labels for
this plan-only demonstration, not authenticated evidence objects.

```sh
export AILANG_BIN=$HOME/.pinned-ailang/ailang
DEMO_ARGS=$(mktemp)
python3 - "$DEMO_ARGS" <<'PY'
import json, sys
input_doc = {
    "phase": "plan",
    "args": {
        "goal": "demo goal",
        "candidates": ["candidate-a", "candidate-b"],
        "roster": [
            {"kind": "test-invocation", "ref": "test-ref"},
            {"kind": "recorded-effect", "ref": "effect-ref"},
            {"kind": "authored-claim", "ref": "claim-ref"},
            {"kind": "proof-envelope", "ref": "proof-ref"},
        ],
    },
}
with open(sys.argv[1], "w") as out:
    json.dump(json.dumps(input_doc, separators=(",", ":")), out)
PY
(cd packages/world-demo && python3 -c 'import subprocess,sys; sys.exit(subprocess.run(sys.argv[1:],timeout=40).returncode)' "$AILANG_BIN" run --quiet --caps IO --entry main --args-file "$DEMO_ARGS" approve.ail)
rm "$DEMO_ARGS"
```

Expected output:

```json
{"plan":"world/effect-plan/v1","effects":[{"id":"e1","effect":"Human.Approve","scope":"worktree","cost":1,"payload":{"requester":"demo"}}],"finish":false,"result":null}
```

The host executes that effect and records a pending approval request. Goal,
candidates and roster remain in the invocation input and its provenance chain;
the strict Human.Approve handler receives only `requester`. The plan asks for no
finish invocation (`finish=false`), and its plan result is null. A direct
`phase=finish` call refuses rather than fabricating a decision or output. Publishing
this descriptor and minting the release session remain attended operator steps.

Unknown or duplicate input/argument/roster keys, unknown or missing phases,
non-object inputs, absent or empty goal strings, candidates other than exactly
two distinct nonempty strings, and malformed rosters produce a zero-effect
`world/effect-plan/v1` refusal with `finish=false` and
`result={"ok":false,"refused":"invalid demo approval input"}`. The roster must
have exactly four entries, one for each admitted kind, with distinct nonempty
reference strings and no caller grade. The plan does not authenticate references,
resolve proof envelopes, or derive PROVEN from a kind string. Later host roster
validation owns those duties.

```sh
(cd packages/world-demo && run_ail "$AILANG_BIN" check approve.ail)
(cd packages/world-demo && run_ail "$AILANG_BIN" test --format json approve.ail)
(cd packages/world-demo && run_ail "$AILANG_BIN" ai-check -timeout 5s approve.ail)
go test -run '^TestWorldDemoApprovalExactOnePlan$' ./host/coordinator -count=1 -timeout 120s -v
```

Expect named `oneApprovalEffect` and `rosterKindAllowed` VERIFIED, 12 passing
inline identities plus generated properties, zero errors/counterexamples/skips
or vacuous skips. The Go test runs the standalone source in an empty root and
asserts effect count before proof tests or downstream coordinator checks; the
valid two-effect mutation must fail at `len(plan.effects)==1`.
