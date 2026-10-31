# Row 231 evaluator report (iteration 328, r1)

Head 96eae6cc, worktree .wt-world-iter328-eval, final `git status --porcelain` empty.

## Score: 90 / 100, PASS
- Correctness / AC coverage: 32 / 35
- Test quality and non-vacuity: 21 / 25
- Design / plan fidelity: 24 / 25
- Hygiene / docs: 13 / 15

## BLOCKING findings
None.

## NON-BLOCKING findings
1. Partial loosening of background edge tolerances is not killed. Own mutants M1-bgY, M1-bgX0 and M1-bgbottom (each ±2 -> ±200 on one background edge) all stay green across the position tests. Only removing the whole check (P8) is pinned by `neg-background-over-content`.
2. Phone (i) entry tolerance 4 -> 40 (M2-nearentry) survives. The header tolerance mutant is killed. F8 (remove (i)) is the only entry-side killer.
3. M1-gapnonneg (`gc >= 0` removed in the non-free-gap branch) survives. Likely near-equivalent on this baseline, since R1 gaps are well above 2 px.
4. The `main.go` routing cases are exercised only by the M3 docs subtests (`position-lines`, `phone-lines`). The position and phone tests call `runPosition` and `runPhone` in-process. Swapping the routes was killed only by `TestMarketingRunbookMatchesHelp`.
5. Reverting §13's export to `r1-rects.json` reds `rects-value` only. The design (D1) also expected `position-lines` red. NB2 itself is confirmed.
6. The second `t.Skip` in `supervisor_test.go:121` is the opt-in `WORLD_KILLPG_STRESS` test from row 150. It is not new, and I count it as acceptable.

## AC verification
- AC-1: `go vet` rc=0. `go test -count=1 ./scripts/marketingcheck` rc=0 (57 s). `-race` rc=0 (166 s). Both parity tests pass. The pin-registry subtest reds when I zero the p0check hash in the orchestration doc (position parity only). Each full hash (0a9a6795..., 31576922...) appears exactly once in the doc.
- AC-2: the real-corpus run (`WORLD_PHONE_CORPUS=...`, `-run '^TestMarketing(Position|Phone)' -v`) gave rc=0, 0 `--- SKIP` lines, 135 `--- PASS` lines, 0 `--- FAIL` lines.
- AC-3: the anti-vacuity guards fire.
  - Removing the zero-regions guard reds `candidate-zero-regions`.
  - P9 (plain Unmarshal of the baseline) reds `baseline-empty-content`.
  - P10 reds `no-candidates`.
  - `t.Skip` appears only in `TestMarketingPhoneRealCorpus` and the row-150 stress test.
- AC-4: all mutants below behave as expected, apart from the findings above. The runner's own deadline test fires.
- AC-5: the docs subtests pass. Reverting §13's export value reds `rects-value`.
- AC-6: not re-run (the controller ran `verify_go.sh`).
- Scope: `host/store/context_roots_test.go` has +5 lines, all allow-list entries. Nothing touches `tools/launchd`, `.github` or `go.mod`.

## Parity spot-checks
I built the banked `p0check` and `phonecheck` myself and compared them with `marketingcheck position` and `marketingcheck phone`.
- Position: 7 inputs, covering ACCEPT and REFUSE across (a1), (a2), (a3), (b) and background. Verdicts and reason strings matched.
- Phone: 4 real corpus images, all r3 ACCEPT in both. Title-right and thumb values matched. The oracle's r2 REFUSE for seed5 and seed3 is the `header-15px` ledger-style difference.

## Fixtures
- Re-running the banked `plan-gen-position-fixtures.sh.txt` in a scratch tree reproduced all 34 position files byte-identically (`diff -r` is clean).
- Re-running the phone skeleton writer (`WORLD_PHONE_SKELETON_WRITE=1`) left no git diff.
- P13 (one byte appended to a fixture) reds `TestMarketingPositionFixtures`.

## Subprocess runner
`boundedRun` sets Setpgid, a context deadline, a group-kill Cancel, a bounded WaitDelay and a post-run group kill. `TestMarketingSubprocessDeadline` uses a hanging script with a grandchild and checks the `subprocess-timeout:` text, the elapsed bound and that the grandchild died. It is not vacuous:
- P17 (deadline removed) reds in 20 s.
- With both the group kill and the post-run kill removed, it reds with "grandchild survived".
- P17b alone (Cancel uses `Process.Kill`) is green, because the post-run group kill covers it. That is defence in depth, not a gap.

## Mutation table
| Mutant | Killer | Red? |
|---|---|---|
| M1-a1tol (a1 tolerance 2 -> 9) | PositionCorpus/iso-a1, Parity | red |
| M1-missingreg (missing-region reason dropped) | Corpus/neg-attested-dropped | red |
| M1-headtop (a2 tolerance 2 -> 20) | Corpus stores and v65 negatives | red |
| M1-schema (schema string changed) | Corpus | red |
| M1-badrect-rc (rc 1 -> 0 on bad baseline) | CLI/baseline-empty-content | red |
| M1-linestep (line-step tolerance 2 -> 12) | Corpus/neg-tested-height+14, Parity | red |
| M1-label (duplicate-label check removed) | CLI/duplicate-label | red |
| M1-inside ((b) bound h-40 -> h) | Corpus/neg-light-header-whole-view and others | red |
| M1-mainroute (position routed to `runPhone`) | RunbookMatchesHelp/position-lines | red |
| M1-bgY, M1-bgX0, M1-bgbottom (background tolerance loosened) | none | GREEN (finding 1) |
| M1-gapnonneg (`gc >= 0` removed) | none | GREEN (finding 3) |
| M2-title (title-right 387 -> 300) | Arms/v67/title-unclipped | red |
| M2-dims (height check dropped) | Arms/new/tall-845, Parity | red |
| M2-near (header tolerance 4 -> 40) | Arms/judge/title-ink-spoof, Parity | red |
| M2-decode (decode failure rc 1 -> 0) | CLI/missing-file | red |
| M2-rc (refused rc 1 -> never) | CLI/refuse-rc | red |
| M2-y0 (thumb present check removed) | Arms/v67/no-scrollbar | red |
| M2-ii ((ii) reasons dropped) | Arms/v67/timeline-only | red |
| M2-entrystart (entry union 6 -> 5 bands) | PhoneSkeleton | red |
| M2-nearentry (entry tolerance 4 -> 40) | none | GREEN (finding 2) |
| M3-D1 (QUICKSTART export reverted to `r1-rects.json`) | RunbookMatchesHelp/rects-value | red |
| M3-mainroute-phone (phone routed to `runPosition`) | RunbookMatchesHelp/phone-lines | red |
| Plan P1 | Corpus | red |
| Plan P9 | CLI/baseline-empty-content | red |
| Plan P10 | CLI/no-candidates | red |
| Plan P12 | PositionParity | red |
| Plan P13 | PositionFixtures | red |
| Plan P16 (position side) | PositionParity/pin-registry | red |
| Plan P17 | SubprocessDeadline | red |
| Plan F1, F4, F11 | PhoneArms, PhoneArms, PhoneCLI/blank | red |

## Commands (all rc=0 unless a mutant)
- `go vet ./scripts/marketingcheck`
- `go test -count=1 -timeout 20m ./scripts/marketingcheck`
- `go test -race -count=1 -timeout 20m ./scripts/marketingcheck`
- the real-corpus `-v` run
- the oracle builds with `go build`
- the fixture regeneration and diff
- the mutant harness at `/tmp/e328/mut.sh`, which restores via `git checkout`
