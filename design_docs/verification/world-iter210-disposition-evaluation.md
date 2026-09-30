# Supplemental evaluation: ADMIN disposition of PARKED row-114 draft plus observational baseline (iter 210), independent Sonnet judge

Base: worktree HEAD = 9b2cfac20f770811b113d62e987530f7e6299660 (clean). Designer commit b702dad exists in the repo (`docs: preserve parked row114 protocol and independent quorums`).

## Verdict: PASS (narrow administrative banking only), 86/100, 0 blocking findings

Bankable ONLY as two things:
1. A PARKED, REJECTED r2 design draft, preserved unmerged as a record. **Design NOT APPROVED.**
2. One observational, filesystem-sandboxed baseline (see `eval-baseline-sonnet.md`, PASS 84 narrow).

This review does NOT pass the design. There is no implementation, no walk, and no sprint or plan under it. **0/3 qualifying paired questions.** The walk arm is NOT-RUN(NO-CORPUS). Capture inventory is UNKNOWN. **No bar credit for clause 5, and no World walk credit.**

## Inputs read (all under `.ailang/state/iter210-supplement-inputs/`)
- `design_docs/planned/w-provenance-teeth-unseeded-baselined-walk.md` (241 lines, read in full)
- `design_docs/verification/world-iter210-quorum/quorum-r1.json` and `quorum-r2.json`
- `world-iter210-baseline.md` and `world-iter210-baseline-evaluation.md` (the proposed bookkeeping)
- My earlier baseline report and the baseline inputs, re-consulted where noted.

## R2 blockers checked against source at 9b2cfac

**1. Method enforcement (GLM blocker): CONFIRMED, and worse than stated.**
- `host/daemon/daemon.go:713-729` registers `POST /v1/commit` (:721), `GET /workbench` and `POST /a2a/` (:729).
- `isProtected` (:738-740) is `r.Method == POST && r.URL.Path == "/v1/commit"` only. The comment at :733-737 says the GET routes are an "unauthenticated declared residual R1".
- The comment at :724-727 says the A2A routes are **NOT in isProtected**.
- The draft's own P9 admits `isProtected` is "a routing fact, NOT a confinement mechanism". §P3-W allows "network except loopback scratch port" with no method filter.
- So W's sanctioned channel reaches the unauthenticated `POST /a2a/`. That is a transition-invocation surface, and row 106 (`27f2574`) has landed it. The PR deny-controls (§5 PR) contain no POST deny-leg.
- The blocker is real. The draft's §16 V3 quotes `isProtected` but never claims it enforces anything.

**2. Ancestral bundle (GLM blocker): CONFIRMED in substance, with an imprecise example.**
- `29e1336` is present locally but is NOT an ancestor of 9b2cfac. It lives on a sprint branch, and its content was squashed as `e437c8a`.
- The diagnosis text is nevertheless ancestral. `aee19c1` (record(206)) reads "CI red after local green (test budget leak, 29e1336)". The charter row at world-mission.md:1891 also carries it. A bundle truncated at the pin therefore contains diagnosis-bearing commit messages and files.
- The P3-B bundle proofs (`rev-list --all --count` equals pin count; `cat-file -e <diagnosis-sha>` fails) cannot detect this. A sha-only check is defeated by squash. A count check is trivially true at the pin.
- The controller's phrasing "diagnosis commits already ancestral to its pin" is right for the record commits and imprecise for the literal `29e1336` sha. That is non-blocking for banking.

**3. Timer (controller blocker): CONFIRMED.**
- Clause 5 (world-mission.md:1000-1006) reads "a provenance walk yields the verified answer in ≤5 minutes each". It bounds the WALK only.
- Draft §5 Timer, "Envelope": "question-issued(B) → verified-result(last arm) total WALL ... **All ≤300 s claims use the envelope**". §6 `PASS` requires "envelope ≤300 s". So B's archaeology time (about 83 s in the C3 baseline) is charged against the walk bar. That changes the clause.
- This is an attended-amendment matter, not a verbatim fix.

**4. Gemini blocker (V6/V7 ellipses): CONFIRMED.** Draft lines 215-216 contain literal `…` in the command cells. The commands are not reproducible as written. V5 (:214) also has "…" inside its quoted result, though that is a quotation.

## Quorum artifacts
- r1: gemini reject ($0.022216), glm reject ($0.057019), Sonnet absent (quota), GPT6.1 absent (unreachable), controller reject; synthesis blocked, total $0.079235.
- r2: gemini reject ($0.01967), glm reject ($0.0486066), same two absent seats, controller reject; synthesis blocked, total $0.0682766.
- **Independent sum: 0.079235 + 0.0682766 = $0.1475116.** This matches the controller's figure. Per-seat sums also check out.
- In r2, one of the three listed blocking objections is the controller's own in-session note, not an external seat. Two external seats blocked. The disposition's "Gemini and GLM present both rounds" is accurate. "No N-1 proceed claimed" is consistent with 2 of 4 seats present.
- Both r1 and r2 recorded costs are metered API only. Pi/Kimi flat-rate authorship and subscription transports are unmetered. The bookkeeping does not claim otherwise.

## Bookkeeping claims vs artifacts
- `world-iter210-baseline.md`: the ≈83 s figure, START stamp and monotonic values agree with the transcript. I recomputed 82.938396125 s earlier. **STOP is absent from the raw transcript** (grep count 0 for "STOP" in `baseline-transcript-*.txt`). The bookkeeping discloses this in three places and makes no ≤300 s verified-answer claim. Honest.
- The cognitive-fence limitation (filesystem fencing, not cognitive blinding; the diagnosing agent and controller had seen the red) is disclosed via my baseline evaluation. It could be stated more explicitly in the baseline.md body. Minor.
- The draft's status line, header and controller appendix (lines 3, 230-240) correctly say PARKED, needs-human-review, no plan, no approved protocol. The body still reads "Net +10 ✅ proceed to quorum" (§15, :202) and describes C3 as "RUNNING" (§8). The appendix expressly supersedes both as historic snapshot. Acceptable but confusing to a future reader.
- **D-WORLD-44 does not exist in the charter at 9b2cfac.** `grep` finds only D-WORLD-43 (row 108). It is a proposal in the appendix (A: one further rotation-designer revision plus full quorum; B: defer for attended regroom; default parked). It must not be treated as ratified. D-WORLD-43 is unchanged.
- Planner derive/resolver `opus fail-closed:planner-lane-field-invalid`: I could NOT verify this. No planner artifact is in the inputs. The draft header line 9 does put a model string (`codex:gpt-6.1-sol`) in the Planner-Lane field, which is consistent with the reported cause. Treated as controller-reported, not independently confirmed. No sprint plan crosses the gate.
- Baseline hash table, 268-file corpus, `fork_turns=none`, run 36682101496 and PUB011 counts: not re-verified here. The first is covered by my prior baseline review, the rest by controller report.

## Findings
- **B0 (blocking): none.** Preserving a labelled-rejected draft and the observational baseline is honest and safe, given the labels in place.
- N1 (moderate). The draft file sits under `planned/`. A future planner could mistake it for approved design. The PARKED banner and appendix mitigate this; a filename or index marker would help. This is not blocking.
- N2 (minor). The correction to record: the GLM "C1/C2 diagnosis commits are ancestors" example is imprecise for `29e1336` (see blocker 2).
- N3 (minor). §16 has "…" in V5's result too.
- N4 (note). The baseline claims stay weak. Two full reds vs 0/10 isolated runs supports context sensitivity only. There is still no product-regression proof and no exclusion.

## Distinction
| Item | Status |
|---|---|
| r2 protocol design | REJECTED, PARKED, NOT APPROVED, unmerged |
| Implementation, sprint, plan | none; none authorized |
| Walk arm | NOT-RUN(NO-CORPUS); capture inventory UNKNOWN |
| Baseline | one observational ≈83 s attempt (n=1, STOP untranscribed, fence not cognitive); attempt 1 NOT-RUN(TRANSPORT) |
| Paired qualifying questions | 0/3 |
| Clause 5 / World bar credit | NONE |
| D-WORLD-44 | proposal only; default remain parked |
