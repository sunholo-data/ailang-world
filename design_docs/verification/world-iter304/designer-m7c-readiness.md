# Designer readiness audit: w-marketing-capture M7c (iteration 304)

**Verdict: REVISED: §2 now says "fixed Chrome" means fixed within one capture run (bounded `--version` before pass A and after pass B, major ≥155, both strings equal, recorded in diagnostics). Exact 155 is no longer implied. New row V29 records the measured Chrome 156 instrument check.**

Role: designer, lane claude-opus-5-5, unattended. Worktree `.wt-world-iter304`, base `039bfa77`.
No git write operations. No Chrome run against real data: the only browser runs used a synthetic `data:` page.

## 1. Chrome version: what should M7c require?

**Ruling: do not require exact 155.** The design's byte-equality promise is scoped to two passes of one capture run on one frozen store. D-WORLD-75 already excludes fresh-seed equality, and §2 already said "cross-platform equality is not promised". It also said to pin the browser in sibling text diagnostics, so the design always meant "record the version", not "require one release". The plan's "require executable Chrome155 version" (plan line 1006, JSON lines 2452 and 2479) over-specifies that. On an auto-updating rig it would make `TestMarketingChromeSixShots` unrunnable for good.

The right rule, now written into §2:
- read the explicit executable's `--version` with a bound, before pass A and after pass B;
- parse the major version and require it to be ≥155 (V13: 154 timed out; 155 and 156 are measured);
- refuse the run if the two strings differ, so an auto-update mid-run invalidates the capture;
- record the string in the sibling diagnostics.

It is not a pixel-similarity loosening. Raw SHA256 equality between passes is unchanged.

| ID | Command | Observed |
|---|---|---|
| R1 | `python3 -c '…subprocess.run(argv,timeout=20,stdin=DEVNULL)…' "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome" --version` | `rc 0 Google Chrome 156.0.8078.47` |
| R2 | `ls -la "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"` | `-rwxrwxr-x … Oct 20 22:47`: updated after iteration 266 (2026-10-14), when V13/V23/V25 measured 155.0.8059.41 |
| R3 | `grep -rn '155' design_docs/planned/w-marketing-capture{,-sprint-plan}.md design_docs/planned/sprint_w-marketing-capture.json` | Design: V13/V23/V25 instrument rows, line 20, r1-3 disposition, §2 header. Plan line 1006 "require executable Chrome155 version". JSON 2452, 2479, 3038. The design itself never stated an exact-version acceptance bar. |
| R4 | `python3 scripts/marketing_capture.py shot --chrome <R1 path> --out <fresh> --url 'data:text/html,<style>body{background:#f7f6f2}@media(prefers-color-scheme:dark){body{background:#1f1e1b}}</style><h1>probe</h1>' --width 390 --height 300 --scheme light\|dark`, two passes per scheme, scratch `/tmp/iter304-chrome156-probe.BrT7/` | All four runs returned `{"verdict": "ok", "detail": ""}`, rc 0, wall ≤1 s. `shasum -a 256`: light-a = light-b = `7e96e882ee93…c818`, dark-a = dark-b = `d1240c97216d…c0ee` |
| R5 | `/tmp/iter304-probe-tool/main.go` (Go `image/png` decode + chunk walk, outside the repo) over R4's PNGs | All `390x300`, chunks `[IHDR IDAT IDAT IEND]`. Background `#f7f6f2` (light) and `#1f1e1b` (dark), so `preferredColorScheme=1/0` still applies on 156 |
| R6 | `grep -n -i version scripts/capture_marketing.sh scripts/marketing_capture.py scripts/marketingcheck/{main,png}.go \| grep -ic version` (control: `grep -c verdict scripts/marketing_capture.py` = 30 lines) | 3 hits, all in `png.go:398-405` (the demo-manifest `Version` field and its check). None is a browser-version read. Nothing at HEAD reads or records the Chrome version. M7c must add it. |

## 2. PD1 and the rectangles: who supplies `WORLD_MARKETING_RECTS`? Is anything blocked on Mark?

**Confirmed coherent. Nothing in M7c is blocked on Mark.**
- PD1 is resolved, not open. The controller amended design §2 rows `workbench-light/dark` in iteration 266, so they now require only the header and the nonempty timeline. The graph and the pending approval move to graph.png and packet.png (plan line 83, "RESOLVED by the controller").
- The plan's M7c phrase "PD1 unresolved visibility must red" is **stale** wording (planner note N9).
- The rectangles are controller-owned. Design §2 says "baseline rectangles are measured in the controller rehearsal, not guessed". JSON `controller_owned` lists "Measured shot rectangles…; real Chrome155 product rehearsal".
- No rectangles file exists at HEAD. Only the verifier and its tests know the schema `world/demo/capture-rectangles/v1` (R8).
- The executor therefore records the real drill DEFERRED or SKIP. That is coherent with the acceptance criteria for three reasons:
  - V21 is explicitly "PENDING (M7 acceptance)".
  - The plan's Final disposition keeps release C8 pixels and byte equality controller-owned.
  - C8/C9 product evidence stays **unmet** until the controller rehearsal. Landing M7c tooling must not be recorded as C8 met.
- The only Mark-gated steps are descriptor publish and the typed-`y` session mint for the release store, which come after M8. M7c and the rehearsal can use V11 test-minted stores.
- Mark question: none.

| ID | Command | Observed |
|---|---|---|
| R7 | `sed -n 81,86p design_docs/planned/w-marketing-capture-sprint-plan.md`; design §2 table rows | PD1 "RESOLVED by the controller (iteration 266, design line 222-223 amended)". The light/dark row text matches. |
| R8 | `grep -rln 'capture-rectangles/v1' .` (positive control: `png.go` must appear) | Only `scripts/marketingcheck/png.go` and `png_test.go` |
| R9 | `sed -n 3036,3043p design_docs/planned/sprint_w-marketing-capture.json` | `controller_owned` includes rectangles/PD1/rehearsal and the release seed/capture/promotion |
| R10 | `grep -n 'D-WORLD-75' design_docs/world-mission.md` | RESOLVED, answered A by Mark on 2026-10-16: frozen-store byte equality, fresh-seed equality not claimed |

## 3. Consistency of the M7c tests and mutants with the M7a/M7b interfaces at HEAD

Baseline: `go test ./scripts/marketingcheck -count=1 -timeout 170s` gives `ok … 18.737s`, rc 0, 20 s wall at load average ~9. `go test -list '.*'` shows exactly 7 tests: PNGStructure, PNGContent, TwoPassInventory and the 4 Supervisor tests. None of the four M7c names exists yet.

| # | Mismatch or gap | Evidence | Planner must |
|---|---|---|---|
| X1 | **The real drill cannot seed from `./scripts/marketingcheck`.** The only real-seed harness, `demoM5cRealSeed` and `newDemoApprovalRig` (in-process rig, V11 mint), is unexported test code in `package main` `cmd/ailang-worldd`. `testMintEnv` is also unexported, in `cmd/ailang-worldd/session_test.go:29`. | `sed -n 1722,1760p cmd/ailang-worldd/demo_test.go`; `grep -rn 'testMintEnv(' cmd host scripts` hits only in `cmd/ailang-worldd` | Either (a) put `TestMarketingChromeSixShots` in `cmd/ailang-worldd`: seed with `demoM5cRealSeed`, stop the rig, then drive `scripts/capture_marketing.sh` on `rig.db`; or (b) have the drill take a controller-prepared frozen store, manifest and key from env. Either way the file list changes. (a) keeps "test-minted through V11" literal. |
| X2 | **Capture has no daemon-binary input.** The attended block builds `$WORLD_CAPTURE_BIN`, but `capture_marketing.sh --help` lists no flag for it. | `sed -n 22,36p scripts/capture_marketing.sh`; plan lines 99–103 | Add `--worldd-bin PATH` (help + docs) or build into scratch. Pick one, prescriptively. |
| X3 | **Serve needs more than the 10 listed flags:** `--workbench-still`, `--bind 127.0.0.1:<port>`, and with `--proof-key-file` at least one `--proof-required-identity`. The seed's identities are `sameRef` and `servesEntry`. | `cmd/ailang-worldd/main.go:250-260` ("Needs --ailang-bin and at least one --proof-required-identity"); `docs/QUICKSTART.md:1553-1557,1576-1577` | Hard-code the two identities and `--workbench-still` in the script, and choose a bind or port strategy (free-port probe or flag). Without the identities PROVEN won't render, and serve refuses. |
| X4 | **Readiness is file-based at HEAD.** `daemon` mode waits for a `--ready-file` that `serve` never writes. | `grep -rn 'ready-file\|readyFile' cmd/ailang-worldd host/daemon` (non-test) = none; control: 3 in `marketing_capture.py` | Add HTTP readiness polling `GET /v1/health` and `/v1/head` (≤10 s per request, ≤30 s startup). Keep the file mode as the unit control. |
| X5 | **`daemon` mode takes at most ONE shot** (`if args.chrome: result = shot(args, stop)`) and prints one JSON verdict. | `scripts/marketing_capture.py` `daemon()` | Add a six-shot pass mode (e.g. `pass --routes FILE`) that emits one JSON verdict with per-shot verdicts and progress on stderr, or loop `shot` from the shell while a background `daemon` holds the server (the `--unit-daemon` pattern). Pass budget is ≤300 s (`PASS_TIMEOUT_S`). |
| X6 | **Routes have one source of truth.** `BuildCaptureManifest` derives routes: packet from role `packet`, graph `?from=N&entry=N` with N = the `approval` why entry, provenance from the `invocation` ref at N plus `&payload=1`. Sizes and schemes come from `Shots`. | `scripts/marketingcheck/png.go:43-53,462-496`; seed roles at `cmd/ailang-worldd/demo_seed.go:703,719,730,739` match | Do not re-derive routes in Python or shell. Add a `marketingcheck routes --manifest M` subcommand to `main.go` that reuses `Shots` and `BuildCaptureManifest`, or the capture-manifest route check will drift. |
| X7 | **No inventory code exists, and there is no list-all-objects REST route.** Routes: `/v1/head`, `/v1/log`, `/v1/log/{index}`, `/v1/objects/{ref}`, `/v1/objects/by-semantic-id/…`. | `grep -n HandleFunc host/daemon/daemon.go` lines 1033–1059 | Define the inventory as head + full `/v1/log` range + objects named by the log entries (or manifest refs) and say so explicitly. Put the guard in production code (`main.go` or the Python script) so mutant M7c-inventory edits production code, not the test. |
| X8 | **The observation comparator must be production code.** M7c-observation and M7c-line-observation neuter "the comparator", but at HEAD it exists only as the plan's inline Python. | plan lines 121–129; `main.go` has only `check` | Put the comparator in `main.go` (e.g. `marketingcheck observe` / `compare-observation`) and test it from `observation_test.go`. A test-local comparator would make both mutants test edits, which the rules forbid. |
| X9 | **Go-surface guard.** Any NEW non-test `.go` file under `scripts/marketingcheck/` reds `TestProductionGoSurface` (`host/store/context_roots_test.go:198-207`), the M7b controller-fix lesson. | allow-list holds only `main.go` and `png.go` | Keep new logic in `main.go` (in the file list), or add `host/store/context_roots_test.go` to the file list with the allow-list entry. |
| X10 | Fake companion PNGs: shell fakes can't synthesize valid content. `png_test.go` `fixture(t, Shot)` and `fixtureSet` already produce scheme-correct, content-passing PNGs plus rectangles in the same package. | `sed -n 145,190p scripts/marketingcheck/png_test.go` | Have the Go test write the six fixtures to disk, keyed by `--window-size` and the `preferredColorScheme` value. The fake Chrome copies the matching one to `--screenshot=`. For M7c-scheme, an argv with no scheme flag must yield a non-dark fixture so `CheckScheme` reds. |
| X11 | Fake companion needs a loopback HTTP fixture daemon (health/head/log/workbench HTML). Socket tests may hit sandbox bind denial. | design C11 note; plan line 23 | Record `bind: operation not permitted` as UNINFORMATIVE UNDER SANDBOX, never PASS. |
| X12 | Verifier CLI (matches the plan): `check --pass-a DIR --pass-b DIR --manifest DEMO --rectangles RECTS [--out NEW]`. Exit 0 ok, 1 refusal, 2 usage. stderr is `marketingcheck: <code> <name>: <detail>`. Output is capture-manifest JSON (`world/demo/capture-manifest/v1`, fields name/sha256/width/height/route, `frozen-head`), to stdout or an O_EXCL `--out`. | `scripts/marketingcheck/main.go` | Use as-is. The Chrome-version stamp goes in a sibling diagnostics text file, not in this schema (png.go is not in M7c's file list). |
| X13 | Supervisor JSON (matches): one line `{verdict, detail[, shot_verdict]}`. Exit 0 only for `ok`, 1 for runtime refusals, 2 for argparse. Verdicts: `ok`, `MissingPNG`, `NotRegularFile`, `EmptyPNG`, `DestinationExists`, `Timeout`, `Interrupted`, `NonzeroExit`, `CleanupFailed`, `InvalidConfiguration`, `SupervisorError`. | `scripts/marketing_capture.py` | Parse by `verdict`. Never treat rc 0 without `ok` as success. |
| X14 | Capture flag list in the plan (`--db --manifest --proof-key-file --ailang-bin --workspace-root --tool-ailang-bin --chrome --out --verify-twice --rectangles`) matches `--help` exactly. `--rectangles` is consumed only by `marketingcheck check`. | `scripts/capture_marketing.sh:22-34` | For a measurement rehearsal the controller needs the pass dirs even when `check` refuses on missing or provisional rects. Keep pass-a/pass-b on disk and report the check verdict separately. |

## 4. Sizing: is M7c one run for claude-haiku-5-5?

**No. Split it.** M7c bundles five artefact changes across five files plus a real seeded Chrome drill:
- the six-shot orchestrator (X2–X6);
- HTTP readiness;
- a REST inventory guard (X7);
- a production observation comparator (X8);
- version stamping (Q1);
- four mutants;
- unresolved placement decisions (X1, X2, X3, X6, X7, X8).

At 28 estimated minutes it is the largest of the M7 stages, and it is not prescriptive enough for a smaller model. Recommended split by acceptance test (the planner writes it):

1. **M7c1, observation comparator (C11).** `TestMarketingProtectedObservationComparator` plus mutants M7c-observation and M7c-line-observation. Files: `main.go` and `observation_test.go`. Self-contained, no sockets, no browser.
2. **M7c2, orchestration with fake Chrome (C7/C8/C9).** `TestMarketingCaptureFakeChromeSixShots` plus mutant M7c-scheme. Covers X2–X6 and X10–X11, the version read/compare/record against the fake Chrome's `--version`, and the `routes` subcommand. Files: shell, Python, `main.go`, `capture_test.go`.
3. **M7c3, frozen inventory (C8/C11).** `TestMarketingCaptureInventoryMutation` plus mutant M7c-inventory (X7). It needs a genuine commit in a disposable store, which needs a session. Place it where V11 seams exist (see X1) or use an importable `host/authority.Mint` + `host/store` path.
4. **M7c4, real drill (controller-run).** `TestMarketingChromeSixShots` with the X1 placement and Chrome rule Q1. The executor may write it, but it SKIPs/DEFERs without `WORLD_CHROME` and `WORLD_MARKETING_RECTS`. Running it is the controller's rehearsal.

## Doc edits made (design only; plan untouched)

`design_docs/planned/w-marketing-capture.md`:
1. Header `**Revised:**`: appended `· iteration 304 designer, 2026-10-26 (§2 Chrome-version rule, V29)`.
2. Verification log: added row **V29** (R1, R2, R4, R5 evidence; synthetic instrument check, not product evidence).
3. §2 heading: `Bounded Chrome 155 capture…` → `Bounded system-Chrome capture (measured on 155 and 156)…`. Verified no inbound anchors: `grep -rn 'bounded-chrome-155\|Bounded Chrome 155'` hit only that heading.
4. §2 body, after "fixed Chrome/OS/fonts.": the run-scoped version rule from Q1 (bounded `--version` before A and after B, major ≥155, equal strings, recorded in diagnostics, equality never across versions).

No viewport, shot list, rectangle, threshold or equality bar changed.

## Notes for the planner

- N1: Replace plan line 1006 and JSON 2452/2479 "require executable Chrome155 version" with design §2's rule. Suggested wording: "bounded `--version` before A and after B, major ≥155, identical, recorded in sibling diagnostics; set-but-invalid, <155 or changed fails". Update JSON 3038 "Chrome155" similarly.
- N2: Apply the X1–X14 decisions prescriptively in the directives: exact flag names, subcommand names, where each guard lives. Haiku should not choose placements.
- N3: Fake Chrome must answer `--version` (e.g. `Google Chrome 155.0.0.0`). Add a control where pass B reports a different version, which must refuse. That makes the version guard non-vacuous without a real browser.
- N4: Production deadlines are already constants in `marketing_capture.py`: shot 30 s, pass 300 s, startup 30 s, TERM 5 s, KILL 5 s. Reuse them; do not invent new ones.
- N5: Carried flake: `TestMarketingSupervisorFakeDaemonCleanup/shell-hup` exceeds its 10 s bound under whole-suite load (29 s, 3/3). Any `capture_marketing.sh` trap change should rerun it isolated and under load.
- N6: The seed manifest carries `proof_version` and `tool_version` (`demo_seed.go:695`), precedent for stamping instrument versions. The capture-manifest schema does not; keep the Chrome string in diagnostics.
- N7: The serve identities `sameRef` and `servesEntry` are fixed by the seed (QUICKSTART 1576–1577). Hard-code them, don't expose them as a flag.
- N8: Thresholds in `png.go` are labelled "provisional until M7c measures real renders". Measuring them is controller rehearsal work, not executor work, and png.go is not in M7c's file list.
- N9: Drop the stale "PD1 unresolved visibility must red" clause (PD1 was resolved iteration 266).
- N10: C8/C9 product evidence remains PENDING after any M7c split lands. The row record must not mark C8 met from fake-Chrome success.

## Quorum r5 revision (iteration 304, the one permitted revision run)

Quorum on the §2 revision was BLOCKED: Gemini 4 Argon PASS, GLM 5.3 REJECT, Kimi K3 REJECT. The four bounded fixes were applied to `design_docs/planned/w-marketing-capture.md`. No shot list, viewport, rectangle or equality bar changed.

| Fix | Doc line(s) after edit | Evidence I re-ran |
|---|---|---|
| 1. V29 supervisor provenance (Gemini+Kimi) | V29 row (line 78): now names the COMMITTED M7a supervisor `38ac4167` + controller fix `98f48b20`, PR #360 merge `c73c7904`, invoked read-only on a synthetic `data:` page. Line 40 ("design_docs only") and line 286 ("future work") get dated annotations; nothing deleted. M7 row (line 298) annotated "M7a built (`38ac4167`, PR #360), M7b built (`fa04df00`/`02e12be2`, PR #361), M7c pending". | `git merge-base --is-ancestor 38ac4167 c73c7904` rc 0 and the same for `98f48b20`. `git log --format='%h %s' c73c7904^1..c73c7904^2 -- scripts/marketing_capture.py` lists exactly those two commits. `git diff --quiet HEAD -- scripts/marketing_capture.py` rc 0. `git log --merges` shows `c73c7904 Merge pull request #360 … row150-m7a` and `72dda61d … #361 … m7b`. |
| 2. M7 fallback authority (Kimi, verbatim) | M7 row fallback clause (line 298) replaced with the reviewer's text. The acceptance sentence (line 308) "These discharge 150(a)–(c)." became "C1–C11 with two byte-identical passes across all six shots discharge 150(a)–(c); if the M7 fallback triggers, 150(b) remains open and the row is reported blocked, not passed." | n/a (wording fix) |
| 3. CLAIMED premise (GLM): measured, not re-argued | New row V30 (line 79). Roster paragraph (lines 191–193) now says the fourth kind's CLAIMED code is pinned by the world-demo package contract `authoredClaimGradeCode`, not by gradeOf, which stays closed. | evidence.ail:4 codes line, :30 `authoredClaimGradeCode`, :32 `ensures { result == (if kind == "authored-claim" then 1 else 0),`. workbench_packet_evidence.go:265-267 `case "authored-claim": row.label = "CLAIMED"; row.producer = authored.Author`. demo_seed.go:274/275/1017. Commits `28ff4db3` (M1a) and `bcd2d7cd` (M5a) exist. **Correction to the brief:** `TestWorldDemoLegacyGradeDrift` (host/coordinator/world_demo_test.go:221) covers SIX legacy constructors in seven cases (CompilerOutput, TestReport true/false, HumanApproval, AiReview, RecordedEffect, ProofReceipt), not three. It also pins `authoredClaimGradeCode("authored-claim")` separately. V30 states the measured list. |
| 4. Dependency SHAs (GLM) | New row V31 (line 80). | `git show -s` printed the three subjects exactly as the brief gives them (full 9319e6bf subject: "…(complete 408 or zero bytes; >=2 s stuck guards) (#233)"). `--is-ancestor` of `c21f303…`: rc 0 for all three. Control: HEAD `039bfa77` gives rc 1. |
| Review status | Line 356: "Quorum r5 (iteration 304)" paragraph with all three verdicts and per-objection dispositions (fixed ×3; GLM CLAIMED = measured-at-HEAD). | — |

Line 6, the `**Revised:**` stamp from the first pass, still applies. Plan and JSON are untouched; notes N1–N10 above stand.
