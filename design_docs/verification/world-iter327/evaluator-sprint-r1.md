# Evaluator r1 — iter 327 (PR #436, head e7743475)
VERDICT: PASS — score 93/100 (correctness 34/35, tests 23/25, evidence 24/25, scope/standards 12/15)

## Reproduced
- Defect mechanism: CI message `expiry 1793383038: 200` = test_now+3601 arm returning 200; base test reads wall clock for expiry then handler reads it again (episode_finish.go:169/:271 `ExpiresAt > now+3600`). A second tick between reads makes +3601 become +3600 -> accepted. Explains it fully; no other cause found (fix pins the clock, and the same flake class exists at the :117 proof bound, now covered).
- AC-1 at base 4d43d8e3 (/tmp copy, spin to Nanosecond>=999_950_000, +3601 arm): 20/20 `expiry N: 200` FAIL. At head (unmodified): `go test -run Finish` ok.
- Mutants on head copy (sha256 e429375dc3e0 / daemon.go e6b70d7e3f3f restored, verified): Q1 killed (Rechecks:232 `200 want 403`), Q4 killed, Q3 killed, R1 killed (ProofExpiryBoundary:204), R3 killed, S3 (wall clock at :269) killed, W1 (New wires `return 0`) killed (AC11:106). Own mutant OWN2 (proof bound 300->360) killed by arm now+301.
- Own mutant OWN1 (drop `|| input.Record != nil` in approve guard) SURVIVES the Finish tests. Pre-existing, not touched by this diff (NB1).
- `go vet ./host/daemon/` clean. `go test ./host/daemon -run Finish -race -count=20`: ok (6.2s).
- Full package plain, -timeout 1200s (216s): one FAIL, TestDaemonStartReportsMissingTemplates: "AILANG_BIN is unset" — environmental in my shell, file not in diff, unrelated. No assertion failure from finish tests; no timeout.
- Evidence tallies: executor-ac2-mutants.txt has 21 rc=1 mutant lines (P1-6,Q1-4,6,R1-4,6,S1-3,W1,TM1) = recorded 21/21; S0 control rc=0. AC-1 base 20/20 matches my 20/20. Controller audit "constructions=11" stated in the file; my grep of `Daemon{` literals found 10+ test literals plus New (not independently reconciled to 11, see NB3). AC-4 classification (broker 0 deps on host/daemon, 3/3 PASS alone, daemon 10m default-timeout panic under load) is consistent with the files; I did not rerun broker.

## Code review
- Production unchanged: New sets finishNow to time.Now().Unix(); bounds 3600/300/3600 identical; `time` import dropped only from episode_finish.go.
- Nil-func risk: none reachable. Test `&Daemon{...}` literals (startup_budget, home_service, guarded_chains, binding, health_exec) — none reference the finish route (grep empty). Real construction path is New only. No silent fallback added (good under coding-standards); a future literal Daemon hitting the route would panic loudly, not mask.
- Race: field written only in New and by tests before synchronous handler calls; read under humanMu on request path. -race x20 clean.
- Test change AC11 `future` +301 -> +3600 removes a second-boundary flake; arm still > +300 so still refused; the exact +301 edge is now owned by the new T2 test.

## BLOCKING
None.

## NON-BLOCKING
- NB1: `input.Record != nil` in the approve guard (episode_finish.go:168) is unkilled by any Finish test (OWN1 survives). Pre-existing gap, out of this row's scope; worth a follow-up row.
- NB2: Full `go test ./...` in executor's AC-4 timed out (host/daemon 10m default, broker 2 failures) under load ~15; classified as load by the controller with alone-reruns, not a clean full-gate green. Authority is CI.
- NB3: construction-audit count (11) not independently reconciled by me; no literal reaches the finish route, which is the property that matters.
- NB4: Scope: diff = 3 code files plus evidence docs; nothing extra in code. Evidence directory is large (42 files) but within docs.

## CI (`gh pr checks 436`, at time of check)
verify gate: pass (35s). go host build+test gate: pending. go race gate: pending. Not waited on.
