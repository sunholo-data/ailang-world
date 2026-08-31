# Ailang World Mission Dashboard

## Latest
- Iteration **139**: row **49** `w-canary-fence-passes-a-gutted-canary` **LANDED** — PR #106 → squash `7ab42aa`, Gate 3b green on the merge commit.
- The judge measured the counterfactual: **22 of 23 mutations survived the fence this item replaced**. The old `strings.Count(src,"stateRoot") >= 2` guard was very nearly no guard at all.
- Evaluator `sonnet` **96/100 PASS, zero blocking**; 23/23 shape-gutting mutations killed (13 from the doc + 10 the judge invented).

## Next
- Rows **50–56**, then **39**. Row **56** is new: the fence is blind to a `t.Skip()`-ed canary (disclosed residual, reproduced first-party, cheap to close in the same AST pass).
- Row **48** waits on `D-WORLD-28` — the only human ask.

## Routing / spend
- Designer: none needed (inherited doc). Planner `opus`. Executor `codex:gpt-5.6-sol` (iteration 138). Evaluator `sonnet`.
- generator≠judge held three ways (Codex executor / Sonnet judge / Opus controller) — the reason iteration 138's Codex controller correctly refused to land this and this one could.
- `metered=$0.0000` of `$5`; billing tripwire CLEAN.

## Parked on Mark
- **`D-WORLD-28` (OPEN)** — one word, A or B: how should `verify_go.sh` guarantee its nested race-control module can execute? Recommendation **A**.
- Row 49 had **no** human ask; its `PARKED-ON-LANE` predicate fired on schedule and resolved itself.
