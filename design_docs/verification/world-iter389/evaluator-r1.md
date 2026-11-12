# Evaluator r1 — iteration 389, row 201 A4 sprint 1 (base 9a6b959d, head 4c2cac17)

**VERDICT: PASS — 88 / 100. No BLOCKING findings.**

| Category | Score |
|---|---|
| Correctness vs design / plan §9 | 27 / 30 |
| Tests (non-vacuous, kill mutants, per-operand) | 24 / 30 |
| Security / fail-closed | 19 / 20 |
| Coding standards | 10 / 10 |
| Scope and size | 8 / 10 |

## Size and scope
- `git diff --numstat 9a6b959d HEAD -- cmd host`: **1,982 insertions + 3 deletions = 1,985** (plan claim confirmed; margin 15 to the 2,000 cap).
- `git diff --name-status`: 19 files, all in plan §5 (plus the plan .md and sprint .json). Nothing under tools/launchd, world/, packages/, scripts/, cmd/. Test modifications: only `host/broker/context_thread_test.go` (5/4 -> 6/5, the single P3 edit; 5 lines changed).
- `host/approval/decision.go` imports only crypto/sha256, encoding/binary, errors; pure, no clock. Law takes `Now` as a fact. Added-line grep: no `t.Parallel`, no CORS header (only a test sending `Access-Control-Request-Method` to prove no answer), no exec site; `time.Now()` appears only in a test helper. Size deduction (-2): the cut leaves 15 lines of headroom, so sprint 2 starts with no slack; honestly bounded though.

## Correctness review (own reading of the production diff)
- Message: prefix(17) || 0x00 || verb || P || Q || BE uint64 = 91 bytes, verb restricted to {1,2}, negative expiry refused; matches design §2.3. Namespace constant distinct.
- Law order D1..D4a..D4b..D5..D6..D7 (SubjectRefusal) then D8, D9, D10 matches design table; D4a evaluated from Discovery alone before Carrying; fires on `!Complete || Undecodable>0`; D4b on 0 and >=2 and on packet code != 0. D10 window min(now+600, deadline) with strict lower/inclusive upper bound.
- Discovery: keyset pages, sentinel at Scanned>=2000 returns Complete=false, exactly-2000 stays complete (empty trailing page), store errors propagate, undecodable/!ok counted.
- L1: parse arm OR prefix arm OR effect-substring arm for publish; parse OR `github-pr:` for finish; test pins `worktree`/`world`/`""` as not reserved.
- L2 refuses `ssh:` and `webauthn:` decidedBy in the publish consumer; typed accepted (must-live test).
- POST: R0 (Authorization/Cookie) then R3-R6 verbatim, strict body decode (no decidedBy field, no unknown field), `humanMu` held across read-facts -> law -> consume proof -> DecideApproval; `decidedBy` is the constant `operator:terminal`; proof phase `decide-<verb>` and name `decision-<Q digest>` (verb-bound, request-bound), consumed (unlinked) before the decision. Replay after restart refused by D6 (chain read).
- Routes are `/surface/...` behind `approvalRoute` (Host check); no `/v1` change.

### P2 (idempotent identical re-decision)
Sound. Comparison is over the whole canonical payload (requestRef, verb, decidedBy, now), so a different request, verb, decider or second can never match; any non-identical second decision returns `ErrApprovalAlreadyDecided`, nothing written. The branch cannot let a *different* decision through. It is unreachable via the POST (D6 refuses first), so it only serves the existing finish path (the premise the plan measured).
### P3
The extra head/object read is the new `findApprovalDecision` walk in decideApproval; census edit is exact and justified (one more head, one more object).
### P8
Accepted: the broker's no-redecide check is unit-covered (TestDecideApprovalRefusesASecondDecision, FirstDecisionMustLive); the POST's D6 sits in front of it.
### The CLI cut
Honest and bounded (CD1–CD4 -> sprint 2 head). Without the CLI nothing writes the operator proof, so POST /surface/decide can never succeed in production: it is an **unused but guarded route** — no reachable attack surface beyond a 403 for anyone who can already pass Host/origin guards; GET is read-only and discloses stored request/packet bytes only on loopback behind the same Host guard. The only latent risk is that the proof directory is writable by the same uid (agent), unchanged from the finish route and out of this sprint's scope.

## Findings (all non-blocking)
1. **Weak test: D1 grammar anchors/case are not pinned at the daemon.** Mutants J31 (drop `$`) and J32 (accept uppercase hex) SURVIVE the daemon sprint tests. Production is correct, but a regression there would let `hashref.MustParse(ref)` (surface_decide.go, after D1) panic on a 65-hex ref. Suggest one POST/GET with a 65-hex and an uppercase ref.
2. **Weak test: L1 effect arm** — J07 (`Contains(scope,"effect=Registry.Publish")` -> `HasSuffix`) survives; the test's only effect-arm case has the effect as a substring at the end. Add a case with the effect mid-string.
3. **GET D4a mapping untested** — J33 (D4a dropped from the GET's refusal-case list) survives; GET would then decode an empty packet and show a bogus 200 subject. Display-only, POST law still refuses.
4. J06, J11, J38 are equivalent mutants (documented below), not defects.
5. Load-sensitive reds in the full host/daemon run (below) all green on re-run, except `TestServeStartsWhenTheProbeSpendsTheStartupClock` which is base-red under load (known list, environmental).

## Mutants I ran (ledger: `$HOME/.ailang/state/world-iter389/judge/mutants.tsv`, 46 lines; every run a line; worktree restored after each, `git status --short` empty)
`comm -12` of my ids vs the plan ledger's ids (field 4): **0 shared** (my J-ids are new; P- rows are re-runs of plan-table rows under my own ids/commands). 41 J/P rows.

| id | edit | result (killing test) |
|---|---|---|
| J01 | DM `expiry<0` -> `<1` | KILLED TestDecisionMessageRefusesBadInput |
| J02 | namespace typo | KILLED TestDecisionMessageLayout |
| J03 | MaxDecisionWindow 601 | KILLED TestDecisionRefusalD10 |
| J04 | `Undecodable>0` -> `>1` | KILLED TestDecisionRefusalD4a |
| J05 | PacketCode arm never | KILLED TestDecisionRefusalD4b |
| J06 | `Carrying != 1` -> `> 1` | SURVIVED — equivalent (0 caught by the earlier arm) — **must-live control** |
| J07 | effect Contains -> HasSuffix | SURVIVED — finding 2 |
| J08 | D10 high `>` -> `>=` | KILLED TestDecisionRefusalD10, D7, MustLive |
| J09 | D10 low `<=` -> `<` | KILLED TestDecisionRefusalD10 |
| J10 | D8 `||` -> `&&` | KILLED TestDecisionRefusalD8, Order |
| J11 | D9 extra operand (equivalent) | SURVIVED — equivalent |
| J12 | window drops `deadline` min | KILLED TestDecisionRefusalD10 |
| J13 | D4a drops Undecodable operand | KILLED TestDecisionRefusalD4a |
| J14 | D5 disabled | KILLED TestDecisionRefusalD5, Order |
| J15 | ChainState undecided returns onChain=false | KILLED TestDecideApprovalFirstDecisionMustLive |
| J16 | L2 `||` -> `&&` | KILLED TestPublishApprovalRefusesAKeyDecision |
| J17 | already-decided sentinel swapped | KILLED TestDecideApprovalRefusesASecondDecision |
| J18 | idempotent path returns zero ref | KILLED TestDecideApprovalRefusesASecondDecision |
| J19 | discovery `carrying = o` (no payload) | KILLED TestDiscoverRequestPacketsCompleteness |
| J20 | short-page `<` -> `<=` | KILLED TestDiscoverRequestPacketsCompleteness |
| J21 | drop `!ok` operand | COMPILEERR (unused var; not a mutant) |
| J22 | cursor = page[0] | KILLED TestDiscoverRequestPacketsCompleteness |
| J23 | drop Scanned++ | KILLED TestDiscoverRequestPacketsCompleteness |
| J30 | typedDecider changed | KILLED TestDecideTypedLands |
| J31 | D1 regex drop `$` | SURVIVED — finding 1 |
| J32 | D1 regex accepts uppercase | SURVIVED — finding 1 |
| J33 | GET drops D4a from refusal list | SURVIVED — finding 3 |
| J34 | window.Min = now | KILLED TestDecisionSubjectJSON |
| J35 | Entry = 0 | KILLED TestDecisionSubjectJSON |
| J36 | proof phase fixed to decide-approve | KILLED TestDecideTypedLands, ReplayAfterRestart, ...BeyondThePrefix (TestDecideTypedProofBindsTheVerb also in the P-PROOF row) |
| J37 | Proposal = packet hash | KILLED TestDecisionSubjectJSON |
| J38 | facts gate on Carrying only | SURVIVED — equivalent (D4a refuses first) |
| J39 | lawClass drops 400 | KILLED TestDecidePOSTLawArms, TestDecisionSubjectJSON |
| J40 | OnChain forced true | COMPILEERR |
| J41 | strict decode (DisallowUnknownFields) removed | KILLED TestDecidePOSTGuards |
| P-H-R0-AUTH | R0 Authorization disabled | KILLED TestDecidePOSTGuards |
| P-H-PROOF-SKIPPED | proof consumption skipped | KILLED TestDecideTypedLands, TestDecideTypedProofBindsTheVerb |
| P-H-DECIDEDBY | decidedBy `operator:x` | KILLED TestDecideTypedLands, ReplayAfterRestart |
| P-H-LAW | law skipped in POST | KILLED TestDecidePOSTLawArms, Discovery* tests |
| P-LAW-D6 / P-LAW-D1 | D6 / D1 branch deleted | KILLED TestDecisionRefusalD6 / D1 (+Order) |
| P-DM-FIELD-ORDER | swap P and Q | KILLED TestDecisionMessageLayout |
| P-DIS-SENTINEL-EDGE | `>=` -> `>` (2,001 treated complete) | KILLED TestDiscoverRequestPacketsCompleteness |
| P-L2-SSH | L2 ignores `ssh:` | KILLED TestPublishApprovalRefusesAKeyDecision |
| P-BR1-NO-IDEMPOTENCE | idempotent branch removed | KILLED TestDecideApprovalRefusesASecondDecision |
| P-RES-FINISH-PREFIX | `github-pr:` arm removed | KILLED TestDecisionRefusalD5, TestReservedApprovalScope |

Result: 36 compile-checked mutants executed (+2 COMPILEERR); 29 killed, 7 survived of which 3 are equivalent (J06, J11, J38) and 4 real test gaps (J07, J31, J32, J33). Every kill is by a behavioural test (no inventory/static test was the sole killer, so no exclusion re-runs were needed). Must-die controls: P-H-LAW, P-H-PROOF-SKIPPED (killed). Must-live control: J06 (survived, as expected).

## Non-vacuity (run via the rows above)
Delete D1 / D6 branch -> DLn test fails; swap P/Q -> TestDecisionMessageLayout fails; 2,001 treated complete -> discovery test fails; no-redecide check dropped -> TestDecideApprovalRefusesASecondDecision fails; L2 accepts ssh -> TestPublishApprovalRefusesAKeyDecision fails; POST law skipped / R0 skipped -> guard/law tests fail; proof consumption skipped -> TestDecideTypedLands/ProofBindsTheVerb fail; decidedBy altered -> TestDecideTypedLands fails; strict-body removal -> TestDecidePOSTGuards fails. (I did not run a separate "read decidedBy from the body" edit; the strict decode mutant J41 covers the supply path.)

## Test runs (AILANG_BIN=$HOME/.pinned-ailang/ailang, own worktree, head 4c2cac17)
- `go test -count=1 -timeout 30m ./host/approval/... ./host/coordinator/`: ok approval 0.6s, approvaltest 0.9s, coordinator 105s (rc 0). Log t1.log.
- `go test -count=1 ./host/broker/ -run '^(TestDecide|TestDecision|TestReserved|TestPublishApproval|TestApprovalContextIdentity|TestApprove)'`: ok 2.3s (rc 0). t2.log.
- `go test -count=1 -timeout 30m ./host/daemon/`: rc 1, 591s; reds: TestMCPPostBudgetProductionConstants, TestSessionSelfGit, TestServeStartsWhenTheProbeSpendsTheStartupClock, TestRegistryPhaseHasItsOwnClock. Re-run alone of the first three (TestSessionSelfGit and RegistryPhase added): **ok 51.9s**; TestServeStarts... is the known base-red load-sensitive test (not re-run separately; not touched by this sprint). All sprint tests passed in the full run. t3.log, t3b.log.
- `go test -race -count=1 ./host/approval/ ./host/coordinator/`: ok (2.2s, 89s).
- `go vet ./host/... ./cmd/...`: clean. `go test -run '^$' ./...`: no compile failure.
- Not run (per directive): verify_go.sh, whole `go test ./...`.
