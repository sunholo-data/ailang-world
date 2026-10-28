# Iteration 314 — independent evaluator r1 (sonnet, Agent tool)

Banked verbatim in substance by the controller from the evaluator's returned report. The judge
ran in its own worktree (`.wt-world-iter314-eval`, at `d7bba842`) and restored it (`git diff
--quiet` clean, re-checked by the controller: 0 porcelain lines).

**Verdict: PASS, 92/100. Zero blocking findings.**

| Dimension | Score |
|---|---|
| Completeness against acceptance criteria | 28/30 |
| Truthfulness of claims | 19/20 |
| Test quality (non-vacuous, mutation-proven) | 24/25 |
| Scope discipline | 10/10 |
| Honesty of the executor's report | 11/15 |
| **Total** | **92/100** |

## Security judgment (V46/P20), derived by the judge from the code

- A leg passes only when `judgeReverse` returns `revPass`: no child, no silent, no foreign
  dial, and a BOUND or `REFUSED ` token other than `REFUSED EADDRINUSE`.
- `classifyReverseDial` sets exactly one of child / foreign / silent on every connected dial,
  and `pollLoopback` dials at least once, so a pass requires zero connected dials.
- A genuine reach returns the child's nonce, which outranks everything in `judgeReverse`.
- A foreign process on the port is retried on a fresh port, at most 3 attempts, then refused as
  "inconclusive", never passed.

## Gates (judge's worktree)

gofmt empty; `go vet ./host/broker/` OK; compile fence OK; `-run TestExecProbe -v`: 17 PASS,
1 SKIP (`TestExecProbeArm8`, needs real srt, not new), 0 FAIL; T6 `-count=3` PASS×3. The ci.yml
step-8 `tests=` line lists all 7 new names and the PASS-guard greps each. CI-only: the real-srt
confirmation that the flake is gone (UNMEASURED by the judge; the controller later read it on
the PR head's Linux job — every new test and `TestExecSrtMCPGrantAndBudget` `--- PASS`).

## Mutation table (each against its named killer only)

| Mutant | Edit | Killer red? |
|---|---|---|
| M1 foreign counted as reach | `case seen.child` → `seen.child \|\| seen.foreign \|\| seen.silent` | yes (T1) |
| M2 no retry | `break` after `tried = append(…)` | yes (T2) |
| M3 any bytes count | nonce equality → `len(buf)>0` | yes (T2) |
| M4 nonce prefix match | → `HasPrefix(buf, "WORLD-PROBE-NONCE ")` | yes (T2 forged-nonce) |
| M5 exhaustion passes | final return → `""` | yes (T3) |
| M9 real reach never matches | `want+"X\n"` | yes (P1; and T4) |
| M13 always retry | `case revPass` → `continue` | yes (T5 clean) |
| M15 hello dropped | no hello write | yes (T2) |
| M17 silent not sticky | `case revSilent` → `continue` | yes (T7) |
| O1 first-attempt token index | `r.token(1)` → `r.token(0)` | yes (T5) |
| O2 loop bound | `<=` → `<` | yes (T3) |
| O3 drop `REFUSED ` prefix | in `judgeReverse` | yes (T2) |
| O4 garbage as silent | `len(buf)>0` → silent | yes (T2) |
| O5 listen script nonce | `q[2]` → `q[1]` | yes (T4) |
| O6 read bound off by one | `len(want)+1` → `len(want)` | yes (T4) |
| O7 EADDRINUSE passes | drop the EADDRINUSE clause | yes (T1) |
| O9 proxy argv on attempt 2 | `attempt == 1` → `== 2` | yes (T2) |
| O8 hello wait 500 ms → 5 ms | — | **SURVIVED** |
| O10 drop `c.resume()` in the child listener | — | **SURVIVED** |

Not re-run by the judge (UNMEASURED by it; the executor reported them killed): M6, M7, M8, M10–M12,
M14, M16, M18–M21.

## Findings (all NON-BLOCKING)

- **F1 (design residual).** A real reach misclassified as foreign on one attempt (e.g. a partial
  read) would be retried, and a later attempt with no connected dial would pass. Accepted by the
  design; the real child always sends the full nonce line.
- **F2 (survivor O8).** `execProbeHelloWait` is not pinned by any test.
- **F3 (survivor O10).** `c.resume()` in the child's listen script is not pinned. Its purpose is
  to avoid a RST on close with unread hello bytes; macOS loopback does not show it, Linux might, as
  a spurious silent refusal.
- **F4 (report accuracy).** The executor did not re-run G10 at base; red-at-base for T6 rests on
  the executor's and planner's accounts.

## Controller reproduction of the judge (first-party, `.wt-world-iter314-eval`)

O8 GREEN and O10 GREEN (both survivors reproduced); M4 RED and M17 RED (M17's red is the assertion
`Why = "", picker calls 2; want the sticky silent refusal on attempt 1`, not a build error). File
restored, sha-identical, 0 porcelain lines.
