# Evaluator r1 - row 240 (w-marketing-checker-tolerance-pins), M1 at e091d27d

Score 91/100, PASS (>= 70). Judge worktree detached at e091d27d; the only change left is this file.
Env: AILANG_BIN=$HOME/.pinned-ailang/ailang, WORLD_TOOL_AILANG_BIN=$HOME/.pinned-ailang-tools/v0.54.15/ailang.

| Category | Score |
|---|---|
| Correctness vs design (35) | 33 |
| Test strength / non-vacuity (20) | 17 |
| Scope & conflict discipline (15) | 15 |
| Evidence integrity (15) | 13 |
| Process claims (15) | 13 |

## Acceptance, re-derived
- Three-function `-v` run: rc 0, 19 `--- PASS` (3 top-level + 14 position + 2 phone subtests), 0 FAIL, 0 SKIP.
- Full package `go test -count=1 ./scripts/marketingcheck/`: ok, 40.6 s. The row-247 flake did not occur.
- `gofmt -l` empty. `go vet ./scripts/marketingcheck/` rc 0.
- Production diff: `git diff --name-only 2c51c581 -- . ':!design_docs'` lists exactly position.go and tolerance_test.go. position.go is `--numstat 1 1`, and the only difference is the trailing comment on :129. phone.go and main.go are unchanged (`git diff --quiet` rc 0).
- Not re-run: `go build ./... && go test ./...` and `verify_ail.sh`. I took the controller's out-of-sandbox round on trust for these, as the task allowed. No `.ail` file changed.

## Mutants run by me (full package each; every file restored and checked with cmp)
| ID | Replacement | rc | FAIL list |
|---|---|---|---|
| T1 | :150 X0 `<=2` -> `<=200` | 1 | PositionTolerancePins and /bg-x0+3 only |
| T2 | :150 Y -> 200 | 1 | /bg-y+3 only |
| T3 | :150 right -> 200 | 1 | /bg-right-3 only |
| T4 | :150 bottom -> 200 | 1 | /bg-bottom-3 only |
| T9 | :101 `dh>2` -> `>3` | 1 | /a4-graph-svg-h+3 only |
| T11 | :129 `e<=2` -> `<=3` | 1 | /a3-timeline-gap+3 **and** /gc-a3-conjunct-source-pin |
| T12 | :129 `e<=2` -> `<=1` | 1 | /a3-timeline-gap+2, /gc-a3-conjunct-source-pin, Corpus/new/perm-roster-drift-2px, Parity |
| T13 | phone.go:168 `phoneR1Entry, 4)` -> `40)` | 1 | PhoneEntryTolerancePins/entry+5 only |
| T15 | main.go:51 position case renamed | 1 | RunRoutesPositionPhone, RunbookMatchesHelp/position-lines |
| T18 | :129 drop `gc >= 0 &&` | 1 | /gc-a3-conjunct-source-pin only |
| C | :150 Y `<= -1` (positive control) | 1 | about 50 Corpus subtests, all 12 new arms, Parity, CLI/accept-rc, position-lines, RunRoutes |

These match the design's V08/V13 claims. The only difference is that T11 and T12 also red the source pin (below).

### Mutants I derived from the diff (all in tolerance_test.go unless stated)
| ID | Mutation | Result |
|---|---|---|
| D1 | premise helper `gb < 2` -> `gb < 1` | **survives** |
| D2 | premise helper: stop skipping free gaps | **survives** |
| D3 | premise helper: delete the `n == 0` non-vacuity guard | **survives** |
| D4 | drop `regions != 17` from the position assertion | **survives** |
| D5 | source pin `n != 1` -> `n < 1` | **survives** (it still covers deletion) |
| D6 | routes test: phone `rc != 0` -> `!= 99` | killed (RunRoutesPositionPhone) |
| D7 | a3 +2 arm: remove the header `H -= 1` compensation | survives. It is harmless: the +2 arm only needs the gap to stay <= 2 |
| D7b | the same removal on the +3 arm | killed (a3-timeline-gap+3) |
| D8 | phone assertion: drop the `sc.Header` check | **survives** |
| E1 | position.go: `&& true` added to the :129 clause (keeps the semantics, changes the text) | **red** (source pin). This is the brittleness case. |

## Judgement of the design decisions
- **`gc >= 0` equivalence.** The bound is sound. Removing the conjunct can only turn a refusal into an accept. That needs `gc < 0` and `|gc - gb| <= 2`, so `gb <= 1`. R1's non-free gaps are 100/100/15/15/15. The premise is pinned by a subtest that passes on R1 and, per the designer's V12 scratch run, reds at gb 0 and 1. I did not re-run V12.
- **Premise helper weakness (D1-D3).** `positionGapPremise` is only ever called on the real baseline. Its own boundary (`< 2`), its free-gap skip and its non-vacuity guard are therefore never exercised in-suite. The guard is proven only by scratch runs recorded in design docs. A 10-line in-memory test of the helper (a baseline copy with one gap set to 1, expecting an error) would close D1-D3 without touching a tracked fixture. This is non-blocking, because the helper is a tripwire on a fixed baseline and not a judge.
- **r3 source pin versus kimi's proposed_fix.** The r2 artifact's proposed_fix (1) offers a zero-production-diff alternative: "a source-level check that position.go still contains the `gc >= 0 &&` conjunct at the (a3) site". Fix (2) is the extended re-validation wording plus the :129 comment. r3 implements all of it faithfully. The pin is slightly stricter than asked: it matches the whole clause text `okGap, gRes := gc >= 0 && posAbs(e) <= 2, e`, with a count of exactly 1.
- **Pin quality.** It is a good tripwire but brittle by design.
  - It is conservative: any edit of the clause, including a semantically neutral one (E1), reds it, and the failure message says to re-derive §2.5, so a false red costs only a pin update.
  - It does not check the semantics of the gap walk. Judge-logic drift elsewhere (for example in `posFreeGap`) is covered only by the prose trigger.
  - It is blind to a duplicated clause (D5 survives, though D5's count check is `== 1` in the shipped code, so duplication does red).
- **Controller disposition that T11/T12 also red the pin.** Confirmed, and benign. The tolerance edits sit inside the pinned text, so the design's "sole killer" language holds for T18 only. T11/T12 keep their dedicated behavioural killers (a3-timeline-gap+3 and +2), so no coverage is lost, but the pin's failure message ("T18 re-derive") is misleading for those two mutants. I did not rate this blocking.
- **Scope cuts.** §2.6 (the QUICKSTART export observation no longer reproduces) rests on designer V03 (D-EXP reds both rects-value and position-lines). I did not re-run it. It is consistent with the guard at docs_test.go:569 and is plausibly correct. §2.7 (sign-loss mutants, free-gap `gc >= 0`) is a legitimate cut: row 246 owns the drift sign, and the free-gap conjunct is explicitly out of the row's list.

## Process claims
- Quorum r1 artifact: oc-kimi-k3 reject, gpt6-1-sol reject, gemini-4-argon pass; absent_reviewers []. r2 artifact: kimi reject, sol pass, argon pass; absent []. Both match §7 of the design.
- The r3 carve-out stayed within its rule. The only remaining objection was kimi's, it carried a reviewer-authored `proposed_fix`, and it did not dispute direction (it was a tighten-the-guard objection). r3 applied that fix without any new design choice, and sol and argon had passed r2 without objection. The r3 status line is honest that r3 itself was not re-quorumed. That is the carve-out's purpose and I judge it within bounds, but the pin string and comment wording are controller-authored and went unreviewed.
- Banked executor evidence is labelled UNINFORMATIVE where the sandbox denied sockets. I re-derived everything above and did not rely on it.

## Findings
Blocking: none.

Non-blocking:
1. Premise helper is untested at its own boundary and guards (D1-D3). Suggest an in-memory synthetic baseline test.
2. Phone header assertion (D8) and the `regions == 17` assertion (D4) are not load-bearing; harmless.
3. The source pin also reds T11/T12 with a message pointing at T18, which may confuse a future debugger.
4. The source pin is textual and will red on formatting-only edits (E1). This is acceptable and arguably intended.
5. `go build ./... && go test ./...` and `verify_ail.sh` were not re-run here (see the acceptance section); the controller's round on 5a672ea4 is the evidence.

## Queue row 240
Every candidate item is closed or explicitly scoped out. Background edges: four pairs (T1-T8, with T1-T4 re-run and killed). Phone entry: entry+4/+5 (T13 killed). The two iteration-336 addendum items (dh > 2 and a3 e <= 2) are pinned (T9, T11). main.go routes: pinned (T15). `gc >= 0`: settled as equivalent on R1, with the premise pinned and the clause text pinned. The §13 export observation is scoped out because it no longer reproduces.
