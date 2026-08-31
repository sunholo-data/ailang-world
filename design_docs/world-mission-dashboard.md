# Ailang World Mission Dashboard

## Latest
- Iteration **140**: row **50** `w-shell-assignment-parser-drops-an-indented-assignment` **PARKED** `needs-human-review` on new decision **`D-WORLD-29`**. Design doc written, revised once and banked; two quorum rounds, both BLOCKED. No code routed.
- **The doc and the queue row disagree, and the row is the one that needs amending.** The doc rejected row 50's own "tolerate whitespace, require total 1" as a weakening; quorum round 2 refuted its rationale, and I measured it: **indentation is not syntax in bash** — an indented top-level assignment executes exactly like a column-0 one. So the loud red the doc was defending fires because the scanner cannot *see* a real assignment.
- **Row 58 is new and it is about our own instrument**: `scripts/verify_go.sh` is **rc=1 on pristine `dev`** — same 4-test failing set in a loaded and an unloaded arm — while **dev CI is GREEN on the same commit**. Mechanism isolated: a freshly-copied binary costs 1336 ms on first exec vs 96 ms pinned (macOS provenance check), against a 10 s deadline. Rig-local, so every `verify_go.sh rc=0` acceptance criterion is currently broken as written.

## Next
- Rows **51–58**, then **39**.
- Rows **48** and **50** are blocked, each on one open decision.

## Routing / spend
- Designer `claude:claude-fable-5` (rotation entry after deepseek; probe rc=0; authoring + one protocol-mandated revision = the Fable one-DOC ceiling, met not exceeded). Planner lane derived `opus fail-closed:env-pin` but **never spawned** — no plan is owed for a parked doc. Executor/evaluator: none.
- `metered=$0.2514` of `$5` — quorum R1 $0.0997, R2 $0.1314, restored reviewer $0.0203. Billing tripwire CLEAN.
- Quorum round 2 was **full strength**: `oc-glm-5-2` recorded absent (`invalid`) in both rounds was re-run alone and returned PASS, so 2 rejects to 1 pass with `absent_reviewers` closed.

## Parked on Mark
- **`D-WORLD-28` (OPEN)** — one word, A or B: how should `verify_go.sh` guarantee its nested race-control module can execute? Recommendation **A**.
- **`D-WORLD-29` (OPEN, new)** — one word, A or B: should a single *indented* shell assignment be ACCEPTED (A, recommended — row 50's own text, whitespace-tolerant total == 1) or REJECTED (B — the doc's two-sided invariant)? Answering also amends row 50's text, which carries the same false premise.
