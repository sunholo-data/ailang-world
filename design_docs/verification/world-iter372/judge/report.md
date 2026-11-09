# Judge report: row 257 w-exec-cache-seed-at-bind (63bbde90..33b42294)

Verdict: PASS, score 91/100.

## Named tests (isolated copy, pinned ailang)
All pass: broker Seed/LegacyPartialCache/FailureClasses/NoProcessLifecycle; daemon+cmd ExecCache/BindSeed/ReportRacesFirstDispatch; -race ReportRacesFirstDispatch; TestProductionContextRoots (row runSeedAttempt|Background:1 present). gofmt clean, go vet clean, no seedMu/execSeedBudget left, no files outside host/ and cmd/.

## Production review vs section 4
- No call runs against a missing/partial cache: prepareCache ends in awaitSeed, which only returns nil when seedCacheVerified (dir + exact marker, <=64 bytes) holds.
- Copy never on a call's clock: only runBounded for the seed is in seedOne, under runSeedAttempt's context.WithTimeout(Background, seedBudget); calls only select on done/ctx.Done and leave the copy running.
- Cap 3 via execSeedNext under seed.mu; give-up spawns nothing. One copy at a time: running!=nil gate under seed.mu.
- Unique staging per incarnation+attempt, exclusive os.Mkdir, marker written after cp exit 0, publish by os.Rename; error paths RemoveAll staging and do not set ok.
- Handler map read only via boundExecHandler under execMu; execMu is released before SeedStatus takes seed.mu (no lock nesting). execHandler holds execMu while calling StartSeed, which takes seed.mu briefly and does Lstat I/O: order is execMu -> seed.mu only, no inversion.
- No goroutine leak: one goroutine per attempt, bounded by seedBudget.

## Blocking findings
None.

## Non-blocking notes
1. StartSeed idempotence is untested: mutant X-NOT-IDEMPOTENT (guard `if false`) survives. Production calls it once per handler under the cached-handler path, so no safety claim is hit.
2. seedIncarnation is process-wide, not per handler. If execHandler rebinds the same episode with a changed root/profile, the old handler's goroutine and the new handler share incarnation+k naming; the new handler's cleanupStaging could delete the old one's in-flight staging (the old attempt then fails and retries), or an exclusive Mkdir collision fails an attempt. Failures are visible and safe (never publishes unverified), but it costs an attempt. Edge case; worth a follow-up row.
3. StartSeed does Lstat/ReadFile of the caches while execMu is held (bind path), a small latency cost only.
4. MUT-CALLER-CTX not run (needs a structural refactor); the property is covered structurally by TestProductionContextRoots plus TestSeedOutlivesTheCallThatWaits.
5. X-RM-STAGING-OTHERS-OWN (drop the `own` skip in cleanupStaging) survives: equivalent, own staging does not exist at step 3 in normal flow.

## Mutation tally (isolated copy, each restored, git status clean at end)
Plan-table mutants run: 19 (START-NOOP, SKIP-CHECK, NO-CAP, SWAP, BUDGET, MARKER-PREFIX, LSTAT-READY, NO-SIZECAP, FIXED-STAGING, NO-CLEANUP, UNBOUNDED, NO-STEP1, REMOVE-OK, NO-REPORT, NO-BIND, NO-STEP0, NO-LINE, UNLOCKED) = 18 listed killed; plus own mutants 8: X-MARKER-BEFORE-CP, X-GIVEUP-NOTGUARD, X-READY-ON-FAIL, X-ATTEMPT-NOINC, X-LAST-NOTRETURNED killed; X-NOT-IDEMPOTENT survived; X-RM-STAGING-OTHERS-OWN survived (equivalent).
Total 26 run: 24 killed, 2 survived (neither on a section 4 safety claim). (First NO-STEP0 attempt failed to compile; redone with a valid mutant and killed.)
Not run: MUT-CALLER-CTX, MUT-ALWAYS-COPY, MUT-DIRECT-DST, MUT-PROMOTE.
