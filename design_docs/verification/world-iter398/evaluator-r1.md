# Evaluator report r1 — iteration 398, row 243 V2b "the bar" (head b613e447)

## Verdict: PASS — 91/100. No BLOCKING findings.

| Category | Score |
|---|---|
| Correctness vs design/plan §9 (read-only, precedence, CLI-identical) | 28/30 |
| Tests: non-vacuous, kill mutants, cover operands | 27/30 |
| Security / fail-closed | 19/20 |
| Coding standards | 9/10 |
| Scope / size | 10/10 |
| Total | 93 -> reported 91 after the non-blocking deductions below are netted; see findings |

(Arithmetic as scored: 28+27+19+9+10 = 93; I withhold 2 for the unpinned AdmitWhy CLI-line quoting and the cosmetic-string survivors, giving 91.)

## BLOCKING findings
None.

## Non-blocking findings
1. `AdmitWhy`'s CLI line for `--result`/`--expect-store`/`-` refusals is built by `Quote(...)` (flags.go, `cli := Quote(append([]string{Binary, "why"}, args...))`). Mutant J-ADM-CLIQUOTE (replace with `Binary+" why "+strings.Join(args," ")`) SURVIVED the whole host/verbs package, and J-D-CLIQUOTE-DAEMON (daemon `TestSurfaceVerb|TestWhy|TestVerb`) and J-D-CLIQUOTE-CMD (cmd `TestWhy|TestVerb|TestQuickstart`) also survived: no test uses an argument needing quoting in these three refusals. Effect if wrong: a mis-quoted copy-paste line in the reason text only (the panel's copy button uses the route's `cli`). Cosmetic, but a real un-pinned operand.
2. Cosmetic page strings survive the whole surfacepage package: the placeholder suffix ` · why <target>` in app.js (J-JS-PLACEHOLDER), the keys-row text `why T` (J-JS-KEYROW), and the index.html placeholder (J-JS-HTMLPH). No test pins them. Cosmetic.
3. (a) cancelled-request reason ("the request ended before the walk did") has no CLI line in `reason`; the answer's `cli` member is still set and nobody reads the response. Acceptable (design: nobody reads this).
4. Planner-declared survivors: `W-LOGPAGE` (host/why LogPageMax 500->499) survives whole host/why and daemon narrow set; killed only by constant pin `TestWhyConstantsMatchTheCLI` (cmd) — I reproduced (W-LOGPAGE-why / -daemon SURVIVED, -cmd KILLED-by-TestWhyConstantsMatchTheCLI). `V-FLAG-USAGE` (usage string) killed only by static `TestWhyFlags` pin (reproduced). Both non-blocking (V2a-inherited / help-text only).

## Read-only / safety evidence
- Production diff has no new /v1 route, no POST, no store write, no new `Daemon` field, no exec, no `.ail`, no new `context.Background/TODO/WithoutCancel` (grep over the non-test diff: empty). The private mux is `whyRoutes()`' seven `GET` patterns, built per walk, with `mux.Handler(req)` pattern=="" -> "unrouted in-process read" error (no fallback).
- Mutants: J-D-MUX8 (eighth handler `GET /v1/episodes`) KILLED-by-TestSurfaceVerbWhyIsReadOnly,TestWhyGetterTransportParity; J-D-MUXPOST (`POST /v1/head`) KILLED-by-TestSurfaceVerbWhyIsReadOnly,TestWhyGetterTransportParity; J-D-MUXDROPRCPT KILLED-by-TestSurfaceVerbWhyOutcomes; D-MUX-REG KILLED. Behavioural, not only a pin set (outcome/parity tests also die).
- Outcome precedence: (a) cancelled walk: J-D-PRECA-OFF and J-D-PRECA-TOBC (route to a normal result) both KILLED-by-TestSurfaceVerbWhyDeadline — the cancelled-request arm is real (cancelled ctx, handler invoked, asserts the exact reason; without it the walk returned a normal result). (b) deadline: J-D-DEADCTXOP / J-D-DEADTALLYOP (each operand separately) and J-D-DEADREASON KILLED-by-TestSurfaceVerbWhyDeadline (asserts no "not found" and the CLI line). (c) read failure: J-D-PRECC-OFF, J-D-FAILREASON KILLED-by-TestSurfaceVerbWhyReadFailure. (d) shapes: J-D-RESBOTH, J-D-RES-ERREXIT KILLED-by-TestWhyResult; J-D-RES-OUTEXIT KILLED. (e) cap: see below. Unknown target / exit 3 appear as results with stderr/exit and `cli` (golden rows `why 9`, `why 2`).
- Cap: raw-byte comparison (`len(Text)+len(JSON)+len(Stderr)`); boundary tests at 1048576/1048577 across members; `TestVerbAnswerCap` also encodes 1 MiB of 0x01 and asserts body > cap and <= 6*cap+1024 (the 6x claim is tested and true: \u0001 = 6 bytes). J-D-CAP-CMP (`>=`), J-D-CAPREASON, J-D-CAPOUTCOME, J-D-CAPWIRE2 (cap applied after write) all KILLED-by-TestVerbAnswerCap/TestVerbAnswerCapEndToEnd.
- Stale reason: J-LOOK-REASON and L-REASON (old reason restored) KILLED-by-TestLookup, TestSurfaceVerbOutcomes, TestSurfaceVerbWhyOutcomes (and TestVerbRouteWhyIsTheCLIs in plan run).
- Admission: A-ORDER, J-ADM-DASH, J-ADM-RESULT, J-ADM-ES, J-ADM-HELPMSG, J-ADM-TRAILPFX, J-ADM-REASONTXT, J-SCAN-HI, J-TRAIL-NARG, J-LEAD-LEN all KILLED-by-TestAdmitWhy/TestWhyFlags (direct unit behaviour). `C-SCAN` (CLI `--scan` bound disabled via `if false`) SURVIVES the old V2a tests but is KILLED-by-TestWhyBarAdmissionIsTheCLIs when the new test is in the run (C-SCAN-b): the CLI's scan-bound was newly pinned.
- CLI unchanged: `git diff c09e0cdb HEAD -- cmd/ailang-worldd/testdata/` is empty (VC1 golden unedited); TestWhyCharacterisedPaths green.
- Bar == CLI independently: TestVerbRouteWhyIsTheCLIs compares the route to (i) the real CLI path (`whyRun` -> `runWhy`) on the same daemon, and (ii) `testdata/why-characterise.golden.jsonl` recorded from the base CLI (not from `why.Run`); TestWhyBarAdmissionIsTheCLIs compares the usage line to the CLI's first stderr line. J-VE1-BYTE (text + " ") and J-VE1-BYTE-JSON both KILLED-by-TestVerbRouteWhyIsTheCLIs. The route golden `verb-route.golden.json` is a generated artefact checked by the existing golden test (not hand-typed by me to verify further; its `why` rows match the CLI outputs the tests assert).
- Page: P-GO-TURN, P-FINALLY, J-JS-ASKCURRENT, J-JS-ADDRCURRENT, J-JS-DROPVERB, J-JS-TRIM, J-JS-REF-TYPEOF, J-JS-NOEXIT1, J-JS-NOEXITTXT, J-JS-STDERRIFF, J-JS-PANELCOPY all KILLED-by-TestSurfaceAppVerbs (node-run VP7/VP8/VP9/VP10/VP11). Escape-closes-panel and live behaviour: XV2 passes in Chrome.

## Mutant table (ledger: $HOME/.ailang/state/world-iter398/judge/mutants.tsv, 72 lines; harness judge/bin/mut.py asserts exactly-one-match, non-empty `git diff --stat`, and `go vet` of the package before reading; restored after each; `git status --short` clean at end)
Controls: J-CTL-LIVE (comment-only edit) SURVIVED (must-live OK); J-CTL-DIE (AdmitWhy returns asJSON=false) KILLED-by-TestAdmitWhy (must-die OK).
Own mutants (53 ids, `comm -23` of my J-* ids vs the plan ledger ids = 53 not in plan): J-ADM-{DASH,RESULT,ES,HELPMSG,TRAILPFX,REASONTXT,CLIQUOTE}, J-LEAD-LEN, J-TRAIL-NARG, J-SCAN-HI, J-LOOK-{CLASS,REASON}, J-D-{MUX8,MUXPOST,MUXDROPRCPT,CAPREASON,CAPOUTCOME,DEADREASON,DEADCTXOP,DEADTALLYOP,FAILREASON,ADMITREASON,PARSEPFX,RESBOTH,RES-OUTEXIT,RES-ERREXIT,GET-OVERERR,GET-NOROUTE,TALLY-TIMEOUT,HDR,ROUTE-WHY,CAPWIRE2,PRECA-OFF,PRECA-TOBC,PRECC-OFF,CLIQUOTE-DAEMON,CLIQUOTE-CMD}, J-VE1-BYTE(-JSON), J-JS-{REF-TYPEOF,NOEXIT1,NOEXITTXT,TRIM,DROPVERB,ASKCURRENT,ADDRCURRENT,STDERRIFF,PLACEHOLDER,KEYROW,HTMLPH,PANELCOPY}. Result: 47 killed, 4 survived (J-ADM-CLIQUOTE, J-D-CLIQUOTE-DAEMON/-CMD, J-JS-PLACEHOLDER/KEYROW/HTMLPH = 6 survivor lines across 4 distinct edits + controls; see non-blocking 1-2), 0 invalid in the own set; every kill named a unit/behavioural test (TestAdmitWhy, TestWhyFlags, TestLookup, TestWhyResult, TestSurfaceVerbWhy{Deadline,ReadFailure,Outcomes,IsReadOnly}, TestVerbAnswerCap(EndToEnd), TestWhyGetterUnroutedAndOversize, TestVerbRouteWhyIsTheCLIs, TestSurfaceAppVerbs) — none killed only by an inventory/AST/golden test except the two planner-declared pins (V-FLAG-USAGE -> TestWhyFlags usage pin; W-LOGPAGE -> TestWhyConstantsMatchTheCLI).
Plan-table sample re-run (13 ids, M4..M6e incl. JS): A-ORDER, C-SCAN(+b), D-CAP-CMP, D-READCTX (KILLED-by-TestSurfaceVerbWhyDeadline,TestSurfaceVerbWhyIsReadOnly), D-MUX-REG, D-GET-CAP, L-REASON, P-GO-TURN, P-FINALLY, W-MAXRESP (KILLED-by-TestWhyGetterUnroutedAndOversize), V-FLAG-USAGE, W-LOGPAGE: all reproduce the plan's verdicts (W-LOGPAGE survives as declared; C-SCAN as above). First attempt of W-MAXRESP/W-LOGPAGE was INVALID (tab/space mismatch in my edit), re-run valid; both ledgered.

## Pre-existing test edits audit (per file)
- lookup_test.go: 2 rows' reason "why reads World" -> "" (the stale reason is gone; expectation tightened to the new truth). Not weakened.
- surface_verb_outcomes_test.go: `why` row terminal -> refused "a target is required"; 4 `why` lines added to the shape list. Strengthened.
- surface_verb_test.go: outcome terminal->refused with `Reason != ErrTooLong` guard; `verbWhy` added to the call-set pin. Not weakened.
- app_verb_test.go: `why` moved from the terminal loop to VP11 panel refusals; VP6 split keeping `world a b c` and TooLong as `say` cases (and asserting `cli` undefined); VP7-VP10 added. Assertions kept or tightened (32 -> 45 asserts).
- app_test.go: export list += goWith, addressGo. verb_bar_chrome_test.go: why/log asserted in the panel with their exact CLI line; `world a b c` keeps the header-say path (log's old say arm preserved via a different input). quickstart_verbs_test.go: adds `why` to the want map; `seen` keyed by Words (stricter). why_seam_test.go: adds pin of verbs.WhyScan* == why.Scan*. verb-route.golden.json: `why` row regenerated; 3 `why` rows + `why --result f` added.
- Test-function counts base -> head (unique `func Test*` names): host/verbs 6->10, host/daemon 382->392, host/surfacepage 56->56, cmd/ailang-worldd 425->428, host/why 10->10; names removed: none in any package.

## Planner deviations
No separate verb-why.golden.json (answers pinned by the base-recorded VC1 golden plus the real CLI on the same daemon — better, independent of `why.Run`); daemon tests split in two files (new files, D14); §1.4 premises #7 (exit-3 parity on the report seed) and #8 (Escape dispatched in-page) are honest and acceptable. Accepted.

## Scope and size
Measured `git diff -M --numstat c09e0cdb HEAD -- cmd host docs` insertions+deletions = **1,255** (cap 1,400; matches plan's final). `git diff --name-status -M` = 25 paths, all in plan §5 plus the two plan files; nothing under tools/, world/, packages/, scripts/, cmd/world-publish, host/why, skills, no `.ail`; no t.Parallel added (only plan text mentions it); no new exec site or context root.

## Test runs (rc; outputs under judge/)
- `go test -count=1 ./host/verbs/ ./host/why/ ./host/surfacepage/` rc=0 (ok x3). `-race` verbs+why rc=0. `go vet ./host/... ./cmd/...` clean (rc=0). `go test -run '^$' ./...` rc=0.
- `go test -count=1 -timeout 30m ./host/daemon/` rc=0 (258.6 s).
- cmd `^(TestWhy|TestVerb|TestQuickstart|TestDemoSeedWhy)` rc=0.
- Real stack with WORLD_TOOL_AILANG_BIN=v0.54.15: `TestDemoSeedWhyNamedLinks` PASS incl. `real-stack` subtest (5.22 s, RAN, not SKIP) and `TestCLIAgainstSeToolsDaemon` PASS (4.59 s); rc=0.
- Live Chrome (WORLD_REQUIRE_PASSKEY_BROWSER=1): `TestVerbBarWhyInChrome` PASS (1.07 s), `TestVerbBarEpisodeReportInChrome` PASS (1.77 s); rc=0.
- Per-commit spot-check of M5g (55d7c197) was not run separately: the controller's gate sweep and the planner's accept_M5g.log (plan-scratch/logs) cover it; I did not re-run it, so this is not independently verified by me.
- Sequencing: ALL_DONE.b613e447... existed before I ran any go command (second poll).

## Caveat
Daemon mutants were run against the narrow `^(TestSurfaceVerb|TestWhy|TestVerb)` set (kills are valid at that scope); the only daemon-package survivors (CLIQUOTE-DAEMON, W-LOGPAGE-daemon) were confirmed in the narrow set only, not the full 258 s package, but the full-package unit survivors (J-ADM-CLIQUOTE in host/verbs, J-JS-* in surfacepage) were whole-package. Score margin is wide enough that this does not change the verdict.
