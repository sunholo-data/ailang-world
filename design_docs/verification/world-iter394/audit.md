# Audit: w-world-address-bar-verbs r2 (c062bf69), iteration 394

CARVE-OUT VALID (with one required clean-up before planning)

Score: 84/100 (carve-out correctness 34/40, internal consistency 17/25, evidence reproduction 19/20, buildability 14/15)

## Q1 Carve-out validity
Round-2 objections (artifact 10-04-11Z), all 3 REJECT, all carry proposed_fix, none disputes direction:
- glm: bare `episode` falls to step 5 (fix: head-entry/refuse, or remove from RK1); unbounded fetch (fix: AbortSignal timeout, "did not answer in time").
- kimi: same Lookup hole (fix: step 4b, verbatim text); §14 header/V36 (fix: V36 + header correction).
- sol: unbounded wait (fix: named finite budget, AbortController over fetch AND body, finally cleanup, timeout refusal, view/input preserved, never addressTarget, VP2 hang cases, readCtx log row).
Carve-out preconditions (a) and (b) hold.

Fidelity of r2 (git diff 15070038 c062bf69, +93/-8, one file):
- Step 4b (doc:255-262): kimi's text is quoted verbatim. FAITHFUL. Choosing kimi over glm's alternative (class=mutate head entry): glm's fix is an either/or ("or, if bare episode should navigate, remove from RK1"); picking the variant that keeps RK1/T8 is what glm itself calls "the resolution consistent with RK1/T8/VB4" -> faithful.
- The added sentence on `world` (doc:260-262, "applies only to the shapes steps 3 and 3b leave over") is a controller INSTANTIATION of kimi's "or `world`" -- needed, but it is the source of the defect in finding B1 below.
- 15 s budget (`VERB_ROUTE_BUDGET_MS = 15000`, doc:343-352): sol demanded "a named finite budget" without a value; a value must be chosen to apply the fix. NEEDED INSTANTIATION, and justified by V37 (readCtx 10 s; reproduced). Not a direction change.
- glm's message string `the verb route did not answer in time`: used verbatim. FAITHFUL.
- "A second submit while one is pending aborts the first (one controller per bar)" (doc:351-352): NOT in any reviewer text. INVENTION (small, behavioural; neither objection asked for it). Not a direction change, but it is controller design and it creates a spec gap (B2). Recommend deleting or flagging for the sprint planner rather than treating as decided.
- V36/V37/V38 and §14 header: kimi/sol/glm asked for them; applied. Header says V36-V38 ran at 15070038, and V31-V35 at ad144327 (kimi asked "V1-V32 at base, V33-V35 on the r1 branch"; r2 says V31-V35 on r0 commit, which is accurate to what the log shows). FAITHFUL.
- VP2/VP6/VB3/§2.4 rows: all requested by the reviewers. FAITHFUL.
Verdict: the carve-out is valid; one invention (second-submit abort) and one instantiation defect (world clause) are the only departures. Neither overrides an objection nor changes direction.

## Q2 Consistency
BLOCKING (small, doc-only, in-scope for the planner but should be fixed before planning):
- B1. doc:236-237 (step 3b) still says "Any other `world` input of three or more words falls to step 5 and gets `parseAddress`'s own error, as today". r2's step 4b (doc:260-262) says those inputs are SEIZED and refused, and step 5 (doc:263-264) says nothing with a CLI-arm first word reaches it. Direct contradiction in the same list; a builder cannot tell whether `world a b c` is class=address (parseAddress error) or refused. Also 4b gives no refusal text for it: kimi's "the row's own refusal text applies" has no referent because `world` has no step-4 table row (only 3b's `world get` text). The text/class for `world a b c` is therefore undefined, and §2.4 has no row and VB3/VP6 no case for it. (My read: simplest consistent fix is to make step 4b exclude `world` and keep r1's step-3b sentence, since parseAddress already errors; but that is a decision for the designer/Mark, not the evaluator.)
- B2 (NON-BLOCKING in severity, listed here because it is a spec hole). doc:343-352 / §2.4 rows: the abort path says "timeout refusal", yet "a second submit aborts the first" also fires AbortError on the first request; the doc does not say the page must distinguish the superseded abort (silent) from the budget abort (message). As written the first request would print "did not answer in time" over the second's result. Also the budget abort vs network "unreachable" classification is not stated (AbortError must map to the timeout message).

Other consistency checks (pass):
- 4b vs steps 3/3b: `world get <ref>` (3+ words) hits 3b; `world get` (2 words) hits 3 (navigates brane `get`); `world NAME`, `back` step 3. OK.
- §2.4 rows `episode`, `episode frobnicate`, `episode w1` -> 4b; `session frob` -> table row (never) by first word, so the row label "4 (never)" is right but 4b's "no row matched" wording is then vacuous for never/other words (their rows match on first word). Harmless but 4b's "for any other seized first word" can only ever apply to `world`, i.e. exactly the undefined case in B1.
- RK1 (doc:615-619), T8 (doc:430), VB3 (doc:471-475), VP6 (doc:539-543): updated and agree with 4b. 27-word count unchanged and consistent (1+3+5+19=28; 27 shadow).
- Stale phrasing check: "anything else" at step 5 restated; r1 "falls to step 5" survives only at doc:237 (B1). `class=refused` for steps 2/4b vs §2.2 item 3 list (doc:317) which only maps `never, other` to refused: step 2/4b's class=refused has no §2.2 mapping line (NON-BLOCKING; the answer shape for class "refused" should be listed).
- §2.4 trace row for the stalled-body case reads "pass | pass" for R2/step1 with the outcome as timeout refusal: consistent.

## Q3 Evidence (re-run in the worktree, HEAD c062bf69, clean)
- V36: `git diff --stat 48c92a52 HEAD -- cmd host | wc -l` -> 0; control `-- design_docs | tail -1` -> `1 file changed, 787 insertions(+)` (doc says 702 at r1; the control fires, count is just newer; harmless). REPRODUCES.
- V33: arm list identical to the doc's 28-word list, `wc -l` 28; `1118: POST /v1/commit`, `1153: POST /mcp/`; demo.go:12 usage line. REPRODUCES.
- V34: `158: if len(args) != 2 || args[0] != "get"`, surface.js 84-86 as quoted. REPRODUCES.
- V37: readCtx = `context.WithTimeout(r.Context(), d.readDeadline)`, `144: readDeadline = 10 * time.Second`, handler passes ctx to buildEpisodeReport/readEpisodeStart/readEpisodeFinish, `timedOut(ctx, err)`; ReadTimeout present at daemon.go:93/1234. REPRODUCES.
- V38: `episode_report.go:702 const maxReportEpisodeBytes = 256`, `:709`, `episode_finish.go:146`. REPRODUCES.
- V31 (empty-grep with control): `grep -rn Access-Control host cmd` -> 7 lines (comments/tests); control `grep -rn Content-Type host/daemon | grep -v _test | wc -l` -> 28, so the pattern is live; handler files only mention it in comments. REPRODUCES.
- V35: DOES NOT REPRODUCE EXACTLY: doc cites `main.go:374-378` and `374: if text, ok := genericVerbHelp(rest)`; actual is line 375 (call at 375-378). Off by one (the line quoted at 374 is the comment). Substance holds (call precedes the switch). NON-BLOCKING; fix the citation (also doc:200-201 §1.14).
- V5/V12/V18/V20 (sample): `ADDRESS_HELP` at surface.js:11 (doc says 11); why.go 1050 lines, `349: func runWhy`; `doctor.go:461` workspaceEpisodePattern; sprint shortstats 1982+/3- and 1929+/54-. REPRODUCE.
Score driver: 1 off-by-one line cite in 11 rows.

## Q4 Buildability of V1
Largely tight: package/function names, route, schema fields, refusal strings, moved functions (printReport/handleEpisodeReport extraction), test IDs VB/VR/VP/XV, and prior sprint sizes (~1.5-2.0k changed lines) support the ~2,000-line envelope. Decisions a planner would still have to make:
1. `world a b c` class and text (B1).
2. Superseded-vs-budget abort handling and the AbortError-to-message mapping (B2).
3. Where the §2.2 response shape lists class=refused with no `cli` for global-flag refusals (step 2) vs `cli` for 4b.
4. Controlled-clock mechanics for VP2 in the node harness (fake timers vs injected `setTimeout`); not specified, but a routine planner choice.
Size risk: moving printReport + CLI refactor + page + 3 test layers is near the cap; the doc's phase split gives the planner room to cut XV1.

## Q5 Safety
- The only new route is `GET /surface/verb.json` behind `approvalRoute` (R2 Host check; V9-style mux lines 1138-1142 confirm the pattern); no CORS (V31 reproduced).
- V1 runs exactly one verb: `episode report`, an in-process read under `d.readCtx`; everything else is `terminal`/`refused`/`address`/`navigate`. `commit`, `outcome`, `call`, `demo`, `episode start/finish/launch`, never-words are refused or handed back as a copyable CLI line; no POST, store write, exec, or credential open. `navigate` goes only to `/surface/approve/mr-...` or `/surface/decide/sha256:...` (the card's own key flow) and the page re-checks the prefix (doc:327-329); `result`/panel text uses textContent.
- Row 201 rulings not reopened: approve is navigation only (never a bar decision; "any flag or other shape refused"); enrolment typed-only (never row); typed y is untouched (D-WORLD-107 quoted, V29); no approval is created or bound from the page.
- Step 4b only increases refusal; it cannot widen what runs.

## BLOCKING findings
1. doc:236-237 vs doc:260-264 -- step 3b says 3+-word `world ...` falls to step 5; step 4b seizes it; 4b's refusal text for `world` is undefined and no §2.4/VB3/VP6 row pins it. Fix before planning (designer or Mark).

## NON-BLOCKING findings
1. doc:351-352 -- "second submit aborts the first" is a controller invention not in any reviewer text; also unspecified how the superseded abort differs from the timeout message (B2).
2. doc:200-201 and V35 -- `main.go:374` should be 375.
3. doc:317 (§2.2 step 3) -- class=refused (steps 2, 4b) has no mapping line.
4. §14 V36 control "702 insertions" is stale relative to the head (787); say "r1 value".
5. 15 s budget is a controller instantiation (justified by V37); mark it as a tunable the planner may keep.
