# Evaluator r1, iter 335, row 237: PASS 92/100
Rubric: sprint-evaluator 100 pts. Correctness 33/35, non-vacuity 24/25, scope 15/15, evidence 12/15, standards 8/10.
Test asserts 403 + exact text, no approvals head/request/decision objects, pristine-store precondition, and a same-clock/scope control (200, objects match response).
Scope: one new _test.go; episode_finish.go sha256 = base 4e320c53...; no writes outside repo seen.
Stability: -count=20 ok; -race -count=5 ok.
| Mutant | Result | Line |
| M1 row mutant (drop arm) | KILLED | :51 got 200 |
| M3a write-then-refuse | KILLED | :54 head/objects recorded |
| M3c always refuse | KILLED | :61 control 403 |
| M6 different refusal text | KILLED | :51 |
| M7 arm `&& false` | KILLED | :51 |
| M9 refuse only if Scope=="" | KILLED | :51 |
| M5 refuse only if Record.Episode!="" | SURVIVED | test uses Episode "ep1" only |
| M8 record-phase `input.Scope != ""` arm dropped | SURVIVED (Finish tests and full host/daemon pkg, 312s) | adjacent, unpinned |
Non-blocking: M5 is a weak-guard survivor (behaviour-narrowing; a record with empty Episode would be accepted in approve); M8 is a new-row candidate. Both outside the row's claim.
