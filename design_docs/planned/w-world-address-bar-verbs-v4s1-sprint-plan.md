# Sprint plan — w-world-address-bar-verbs-v4, sprint S1 "the pins and the admitted subset" (iteration 407, row 243)

> **THE SPRINT — row 243 phase V4, sprint S1, the first V4 sprint (design §0, §1, §3.1 (the action-id grammar only), §4.1,
> §5.1 (the grammar and identifier rows), §5.6, §6, §7, §8 row S1, §9, §12 W23/W29, §13); D-WORLD-119 and D-WORLD-120
> OPEN — S1 is answer-independent by design §8 ("Waits for: None") and by D-WORLD-119's default:**
> **M1** — the three judge-405 pins as node harnesses over the daemon's own answers: `TestNeedsYouConfirmationEmbeddedNewline`
> (`J-JS-RE-NL`, plus each anchor and the `i` flag on its own), `TestNeedsYouConfirmationsNonArray` (`J-JS-CS-NOGUARD`;
> object, array-like, string, null, number, boolean, a thing with a `map`) and `TestVerbDecidePrefixBoundary`
> (`J-JS-ALLOW-DECIDE`; foreign origin, contained prefix, no colon, case, and the approve/confirm siblings). **Each survivor
> reproduced SURVIVED on the base first (C9); no production defect found, so `app.js` is unchanged.**
> **M2** — NEW `host/verbs/operator.go`: the pure, **unregistered** admitted CLI-surface subset `AdmitOperator(argv,
> OperatorPolicy) OperatorAnswer` over the shared `Lookup` and `ParseInterleaved` (fed `Split`'s words), with a mirror of the
> CLI's two flag sets, the episode grammar, the preset cap, design §3.1's NEW action-id grammar, and the ledger answers as a
> pure input whose zero value admits nothing; the grammar tests (Case, Arity, ExactWords, Canonical, Duplicates, Options).
> **M3** — the policy, never-dispatch, identifier and other-verb tests (Policy, NeverDispatch, IdentifierAnchors/Flags/Length,
> OtherVerbParity). **M4** — the join to the CLI: a 10-line test seam `episodeFlagSetSeen` in `runEpisodeStart` and
> `runEpisodeFinishWith` (the CLI's flag sets as built, before any parse) and `TestOperatorFlagParity`,
> `TestOperatorParsesLikeTheCLI`, `TestOperatorIdentifierGrammarsMatchTheCLI`. **No schema change, no daemon route, no `.ail`,
> `surfaceCSP` unchanged, no served-behaviour change, no CLI behaviour change; the parser has no production caller.**
> - **Measured: 976 changed lines** — `git diff -M --numstat d92983d5 HEAD -- cmd host docs app` (insertions + deletions,
>   renames on; no rename occurs; **the same 976 with no path filter**; there is no `app/` at the repo root; nothing lands
>   outside `cmd` and `host`), after a commit of each of the four scratch stages: **237 production** (`operator.go` 227,
>   the seam 7 + 3), **739 test**. **Cap 1,400: margin 424.**
> - **No cut (S1-D1).** Design 620 raw → **976 compiled (1.57×)**, inside the directive's ~750–1,100; nothing parked.
> - **Four sub-milestones, each green at its commit** (§1.2 C14): **M1** pins (127) → **M2** parser + grammar tests (570) →
>   **M3** policy/identifier tests (778) → **M4** the CLI join (976).
> - **Every line ships as a compiled patch** (Appendix P), extracted from this file and round-tripped onto a fresh
>   `git archive d92983d5` export, byte-identical to the stage tree (§1.2 C16). The executor applies, never retypes.
> - **Pre-existing test files changed: 0; deleted: 0.** Pre-existing production files changed: 2, the seam only
>   (`cmd/ailang-worldd/episode.go` +7, `episode_finish.go` +3; §1.4 F2).

---

## 0. Binding constraints (each is restated where it governs)

- **D-WORLD-119 and D-WORLD-120 OPEN** (controller, iteration 407: `mission_decisions.sh --open` → both OPEN; `--check` valid,
  107 rows). D-WORLD-119's default: "build S1 and S2 against fixtures only. Leave the schema, service installation, browser
  writes, mint and launch paths unchanged." S1 registers nothing: the answers enter the parser as `OperatorPolicy`, a value the
  caller passes, tested for every answer (`TestOperatorGrammarPolicy`, 384 cases); **the zero value admits nothing**, and an
  answer other than exactly `A`, `B` or `C` is unanswered (S1-D5).
- **No daemon route, no schema change, no `.ail`, `surfaceCSP` unchanged, no served-behaviour change, no CLI behaviour change**:
  the gates' `UNTOUCHED_VIOLATIONS` lists every changed file that is not a `_test.go`, `host/verbs/operator.go` or the two seam
  files, plus any non-test file under `host/daemon`, `host/store`, `host/surfacepage/app`, `world`, `packages`, `scripts`,
  `tools`, `plugins` or any `.ail`, and must say `none`; `DAEMON_ROUTES` base = final, `SURFACE_POSTS: 3`; `UNCHANGED:` for
  `app.js`, `surface_page.go` (the CSP literal) and `host/store/credentials.go` (the schema); `CMD_SEAM_LINES: 10 0`;
  `PARSER_PRODUCTION_CALLERS: none` with a positive control (`verbs.AdmitWhy`'s production callers ≥ 1).
- **D-WORLD-109**: no store is opened by any S1 test; the CLI's flag sets are read from a `--help` run under `t.Setenv("HOME",
  t.TempDir())` that returns before any effect; `loadLaunchPreset` reads only a `t.TempDir()`; the gates' Chrome smoke is the
  existing tests on the driver's temporary profile; no real key, `~/.ssh`, `~/.ailang/world`, `/dev/tty`, `gh`, mint or publish.
- Coding standards: the subset is pure (data in, a typed answer out; no I/O, clock or randomness; stdlib imports within host/verbs'
  pinned pure set, S1-D3); effects stay at the host boundary (none here).
- **Not touched**: `tools/launchd/*`, any skill, `world/*.ail`, `packages/`, `scripts/`, `cmd/world-publish/*`, `plugins/`,
  `app.js`, `surface_page.go`, `daemon.go`, every `host/store` file, every pre-existing test file.

---

## 1. Pre-flight

### 1.1 Baseline at `d92983d5` (measured by the planner, 2026-11-18)

`tools/gates.sh <base-tree> base-gates base d92983d5` (the §6 script in `base` mode, on a clone of the repo at `d92983d5`;
`GATES_SKIP_VERIFY_GO=1`), run beside the planner's compile and mutant work (load ≈ 17–27):

| leg | rc | time | note |
|---|---|---|---|
| `TOUCHED_PACKAGES` | — | — | *(empty at the base)* |
| `status` | 0 | 0 s |  |
| `verify_ail` | 0 | 40 s | `✓ verify gate PASSED: 16 required identities verified, 40 named tests pass, 503 se-tools named tests pass` — the final must match |
| `gofmt`, `build`, `vet` | 0 | 0 / 4 / 2 s |  |
| `fence` | 0 | 27 s | `go test -run '^$' ./...` is the test-file compile fence |
| `pkgs` | **1** | 1,327 s | `host/verbs` 0.1 s, `host/surfacepage` 1.3 s, `cmd/ailang-worldd` 460 s **FAIL** `TestDoctorQuickIsFast` (2.57 s against its own budget), `host/daemon` 862 s **FAIL** `TestSeToolsRunBenchmarkTasksEndToEnd` (directive-listed), `TestServeStartsWhenTheProbeSpendsTheStartupClock`, `TestRegistryPhaseHasItsOwnClock` — **all four load-sensitive and green alone at the base** (`logs/base-flaky-alone.txt`: 0.32 s, 15.8 s, 1.94 s, 5.49 s at load ≈ 17) |
| `race` | 0 | 223 s | `host/verbs`, `host/surfacepage` whole; the cmd subset; the daemon verb-route subset |
| `node` | 0 | 5 s | the thirteen named harnesses `--- PASS`, `WORLD_EXEC_NODE` set |
| `inv` | 0 | 24 s | **39 named `--- PASS`**, 0 SKIP, 0 FAIL |
| `DAEMON_ROUTES` / `SURFACE_POSTS` | — | — | **47** (base 47) / **3** |
| `UNTOUCHED_VIOLATIONS` / `CMD_SEAM_LINES` / `PARSER_PRODUCTION_CALLERS` | — | — | none / `0 0` / none (control 1) |
| `UNCHANGED` | — | — | app.js, surface_page.go, host/store/credentials.go |
| `chrome` | 0 | 12 s | **5 `--- PASS`** (the Needs-you pair, XC2, XD1, the bar's report test) |
| **`GATES_DONE 1665s FAILS: pkgs`** | | | the four load-sensitive reds above, each green alone; the final tree is compared leg by leg (C17) |

### 1.2 Codebase claims checked (at `d92983d5` unless "work" — the planner's stage tree, tag `M4`)

| # | claim | command | observed |
|---|---|---|---|
| C1 | **Input 1 — the design's cites, re-measured at `d92983d5` (design measured at `41331b4d`): none moved** | `git show 41331b4d:host/surfacepage/app/app.js \| grep -n …` vs the same at HEAD; `grep -n` each | app.js `:434` `const cs = r && Array.isArray(r.confirmations) ? r.confirmations : [];`, `:456` `function confirmLineHTML(c)`, `:457` `/^vr-[0-9a-f]{16}$/.test(c.id)`, `:473` `notifications`, `:499` `verbOutcome`, `:507` the three-prefix navigate allowlist (`startsWith('/surface/decide/sha256:')`), `:526` the fetch — **identical line numbers at both commits** (PR #604 at `65108353` changed 30 app.js lines, none of these; `git diff --stat 41331b4d d92983d5 -- host/surfacepage/app/` = app.css +29, app.js +25 −5, index.html +1). `split.go:35` `Split`; `lookup.go:36–40` the word sets, `:43–49` `neverWhy`, `:67` `Lookup`, `:90–95` the start/finish rows, `:119` the address fallback, `:129` `ApproveTarget`; `flags.go:14` `ParseInterleaved`; `daemon.go:1154` `GET /surface/verb.json`; `surface_verb.go:14` `VerbSchema`, `:38` the handler; `doctor.go:461` `workspaceEpisodePattern = ^[a-z0-9][a-z0-9-]{0,63}$`; `preset.go:150` `maxPresetNameLen = 51`; `host/daemon/workspace.go:62` `episodeIDPattern` (the same expression). `git diff --stat 41331b4d d92983d5 -- host/verbs cmd/ailang-worldd/{episode,episode_finish,doctor,preset}.go host/daemon/{daemon,surface_verb}.go` → *(empty)* |
| C2 | the CLI's flag sets (the parity's source of truth) | `sed -n '190,218p' episode.go`; `sed -n '63,77p' episode_finish.go` | `runEpisodeStart` builds 22 flags inline (`db credentials-db workspace-root repo branch profile exec-profiles-dir task agent model ttl preset presets-dir resume hosts-dir launch headless request grant budget sessions-dir addr`; bools `launch headless request`; `grant` a repeatable `multiFlag`), `fs.Usage = func() {}`, then `verbs.ParseInterleaved` at `:218`; `--help` prints the static `episodeHelp`. `runEpisodeFinishWith` builds 9 (`preset presets-dir title include-unattributed ready exec-profiles-dir dry-run confirm-on-surface allow-failing`; bools `include-unattributed ready dry-run confirm-on-surface`), parse at `:77`, `--help` prints `episodeFinishHelp` and returns before `env.openTerminal` (F1 at `:88`). **No function returns either flag set** (§1.4 F2) |
| C3 | **Input 2 — the closest landed analogues, and what was mirrored** | `sed -n` of each | **Mirrored**: CP7 `TestNeedsYouListsConfirmations` (`app_confirm_test.go`: `AppScript` + a harness in `t.TempDir()`, `WORLD_EXEC_NODE` non-skippable, a 30 s context, a PASS manifest, over `testdata/requests-confirmations.golden.json` written by host/daemon `TestRequestsJSONConfirmationsGolden`) and `TestSurfaceAppVerbs` (`app_verb_test.go`, over `verb-route.golden.json`) → M1's `runPinHarness` and its three harnesses; V3 S5/S6's `ApproveTarget` `vr-` arm and `TestApproveTarget`/`TestLookup` (table rows over `Split`'s words) → M2/M3's `admitLine` tables; the V1 shared-definition pattern `EpisodeReportFlags`/`WhyFlags` + `TestEpisodeReportFlags`/`TestWhyFlags` (`fs.VisitAll` over the real flag set) and `TestVerbTableMatchesTheCLISwitch` (the CLI's switch arms against host/verbs' word sets, with controls) → M4's parity. **Differs**: the start/finish flag sets cannot move into host/verbs as V1's did (their defaults are cmd's: `defaultWorldPaths`, `defaultPresetsDir`, …), so the join is the CLI's real flag set seen through a seam (S1-D7) and compared with a host/verbs mirror, and the parse itself is compared (`TestOperatorParsesLikeTheCLI`); M1's harnesses call `verbOutcome`/`needsYouHTML` directly (no fetch rig) |
| C4 | **no existing start/finish flag-parity test** (input 2) | `grep -rn 'VisitAll' cmd/ailang-worldd/*_test.go host/verbs/*_test.go`; control the same in host/verbs | cmd: *(empty)*; control: `flags_test.go:18` (`TestEpisodeReportFlags`), `why_admit_test.go:17` (`TestWhyFlags`). `TestEveryVerbHasHelp` (`help_test.go:20`) and `TestVerbTableMatchesTheCLISwitch` (`verbs_parity_test.go:66`) pin verbs and help paths, never a flag |
| C5 | nothing named like S1's code exists | `grep -rnE '\b(AdmitOperator\|OperatorPolicy\|ValidOperatorActionID\|episodeFlagSetSeen\|runPinHarness)\b' cmd host \| wc -l`; control `grep -rn 'func AdmitWhy' host/verbs` | `0`; control `flags.go:136` |
| C6 | every type, field and function the patches name | `grep -n` at HEAD | **existing**: `Split`, `Quote`, `Lookup`, `Verb{Words, Class, Reason, CLI}`, `Never`, `Refused`, `Binary`, `OperatorWords`/`NeverWords`/`OtherWords`, `ParseInterleaved` (host/verbs); `runEpisodeStart(args, stdout, stderr, env)`, `runEpisodeFinishWith(args, stdout, stderr, env, record)`, `sessionEnv`, `exitOK`, `workspaceEpisodePattern`, `maxPresetNameLen`, `loadLaunchPreset(ctx, name, dir, profilesDir, workspaceRoot)` (cmd); `AppScript`, `needsYouHTML`, `esc`, `verbOutcome` (app.js `module.exports`). **NEW**: `OperatorFlag{Name, Bool, Admit}`, `EpisodeStartFlags`, `EpisodeFinishFlags`, `EpisodeIDPattern`, `OperatorActionIDPattern`, `MaxPresetName`, `ValidOperatorActionID`, `Answer`, `OperatorPolicy{Start, Finish, HandoffProven}`, `OperatorClass` (`OperatorReview`, `OperatorUsage`, `OperatorTerminal`, `OperatorNever`, `OperatorOther`), `OperatorIntent{Kind, Episode, Preset, Headless, ConfirmOnSurface, Ready, Argv}`, `OperatorAnswer{Class, Intent, Field, Reason, CLI}`, `AdmitOperator`, unexported `boolWords`, `operatorValue` (host/verbs); `episodeFlagSetSeen` (cmd); every test named in §3 |
| C7 | **Input 3 — fixtures from the real tool** | the test files | M1: `requests-confirmations.golden.json` (written by `TestRequestsJSONConfirmationsGolden` over the real handler, `WORLD_UPDATE_CONFIRMATIONS_GOLDEN=1 go test ./host/daemon/ -run '^TestRequestsJSONConfirmationsGolden$'`) and `verb-route.golden.json` (checked byte for byte by `TestSurfaceVerbOutcomes`, `WORLD_UPDATE_VERB_GOLDEN=1 go test ./host/daemon/ -run '^TestSurfaceVerbOutcomes$'`); each hostile row is the real row with one field changed. M2/M3: every bar line is cut by the real `Split`; Lookup's own answer is the expectation for every non-operator line; `TestOperatorOtherVerbParity` walks every key of `verb-route.golden.json`. M4: the CLI's flag sets are the ones `runEpisodeStart`/`runEpisodeFinishWith` build (seam), parsed by the same `verbs.ParseInterleaved`; the preset cap is where `loadLaunchPreset` itself refuses. No hand-typed bytes or JSON |
| C8 | **Input 4 — tests that enumerate packages, imports, symbols, routes, app.js exports, harnesses, verbs, flags, help** | `grep -rn` by name; each run on the final tree (C17 `inv`) | **unchanged and green**: `TestSurfaceVerbIsReadOnly` (host/daemon: pins host/verbs' production imports to `errors flag fmt io regexp strings time unicode/utf8` — **it turned the first draft red**, §1.4 F3), `TestDaemonDependencyAllowlist`, `TestWorldBoundaryDependencyAllowlist`, `TestProductionContextRoots`, `TestProductionGoSurface` (host/verbs is under `host/`), `TestEveryVerbHasHelp`, `TestVerbTableMatchesTheCLISwitch` (no new verb), `TestFrozenV1RouteTableMatchesMux`, `TestConfirmRoutesRegistered`, `TestDaemonCredentialWriterSites`, `TestCLICredentialWriterSites` (the seam opens no store), `TestSurfaceAppPure` (`module.exports` unchanged: app.js unchanged), `TestAddressHintsArePinned`, `TestSurfaceAppMatchesWorkingCopy`, `TestDoctorEpisodeGrammarMatchesTheDaemon`, `TestEpisodeFinishHelp`, `TestQuickstartSection19MatchesTheBar`, `TestQuickstartSection20MatchesTheFinishFlag`. **Changed: none. New inventories: none** (M4's parity tests are behavioural joins) |
| C9 | **the three judge-405 survivors reproduced on the base first** | `tools/mut.py <clone at d92983d5> <id> … '.' --tag=base-whole` | `J-JS-RE-NL` (`/^vr-…$/m`), `J-JS-CS-NOGUARD` (drop `Array.isArray`), `J-JS-ALLOW-DECIDE` (`includes`): **SURVIVED** `go test -count=1 ./host/surfacepage/ -run '.'` (2–3 s each); their siblings `P3-DECIDE-NOCOLON`, `-INDEXOF`, `-CASE`, `P3-APPROVE-NODASH`, `-CASE`, `P2-CS-TRUTHY-OBJ` **SURVIVED**; `P1-RE-CARET`, `-DOLLAR`, `-I` **KILLED** by the pre-existing CP7 (§4) |
| C10 | **is any pinned code wrong for a real input?** | work: the three harnesses on the base `app.js` | all `--- PASS` before any change: a JS `$` without `m` matches only at the end of input (no trailing-newline allowance), `Array.isArray` and `startsWith` are correct. **No production defect; app.js unchanged** (§1.4 F8 records the one unpinned shape found) |
| C11 | **the parser has no production caller** | the gates' caller search (S1-D10) with a positive control in the same call | `PARSER_PRODUCTION_CALLERS: none (control: verbs.AdmitWhy production callers 1)` |
| C12 | **Input 5 — no asynchronous wait** | the test files | no S1 test waits on asynchronous work: the node harnesses are synchronous scripts under a 30 s `exec.CommandContext` (the CP7 bound); the parser tests are pure calls; M4's `--help` runs return synchronously. No sleep, tick or clock anywhere in S1's tests (directive 12 not engaged) |
| C13 | the new tests, work | `go test -count=1 -v` by name (`accept.sh`) | all `--- PASS`, none SKIP (node set) |
| C14 | **every sub-milestone is green when committed** | `tools/accept_all.sh` (a worktree of `work` at each stage tag running §3's `accept.sh`) | on the final stage trees (`logs/accept_all.txt`, `logs/accept/*.txt`; worktrees `st-M1`…`st-M4` at the stage tags; load ≈ 17–27 beside the gates and the mutant shards): `M1 fast rc=0 ACCEPT_OK` (18 s), `M2 fast rc=0` (33 s), `M3 fast rc=0` (32 s), `M4 fast rc=0` (60 s), `M4 cmd1 rc=0` (313 s), `M4 cmd2 rc=0` (51 s), `M4 cmd3 rc=0` (74 s). **Measured and changed before shipping**: (1) M2's first acceptance went **red** on `TestSurfaceVerbIsReadOnly` (`../verbs/operator.go imports strconv`, F3) — fixed in `operator.go`, all stages rebuilt; (2) the whole-cmd halves first measured `^Test[A-E]` **607 s** at load ≈ 20 (over one 600 s call), so M4's whole package is three parts |
| C15 | the snapshot block (§7) on every sub-milestone | `tools/snap_rehearse.sh` (a clone at `d92983d5`, the patches in order, §7's block verbatim) | `M1 1/0/1`, `M2 3/0/3`, `M3 4/0/4`, `M4 7/0/7` (`STATUS` lines / `DELETED` lines / MANIFEST lines), each `MANIFEST_OK`, rc 0; dirs `host` for M1–M3 and `cmd host` for M4 (the `find` guard holds both ways) — `logs/snap.txt`, re-run on this final plan file |
| C16 | **the round trip** | `tools/roundtrip.sh <this plan>` | `CMP_EQUAL` for all four (each block extracted from this file with §2's awk equals the generated patch byte for byte); each applied with `/usr/bin/patch -p1 -E` after a `--dry-run` to a fresh `git archive d92983d5` export, no `.orig`/`.rej`, committed and measured: `NUMSTAT M1 127`, `M2 570`, `M3 778`, `M4 976` (each `unfiltered=` the same); then `diff -r -x .git` against the stage tree (`git archive M4`) → **`ROUNDTRIP_IDENTICAL_TO_STAGE_TREE`** (`logs/roundtrip.txt`, re-run on this final plan file); `accept.sh` and `gates.sh` extracted from this file are byte-equal to the scripts that produced C14 and C17 |
| C17 | **the gates on the final tree** | `tools/gates.sh <final-tree> final-gates final <sha>` | `tools/gates.sh <final-tree> final-gates final 86b8b79a` (`final-tree` = a worktree of the planner's clone at the stage tag `M4`, the four stage commits on `d92983d5`; `diff -r` equal to the round-trip tree, C16; `GATES_SKIP_VERIFY_GO=1`; load ≈ 7–20): **`TOUCHED_PACKAGES: cmd/ailang-worldd host/surfacepage host/verbs`** (= the plan's list); `status` 0; **`verify_ail` rc 0** (`✓ verify gate PASSED: 16 required identities verified, 40 named tests pass, 503 se-tools named tests pass` — the base's line, byte for byte); `gofmt`, `build`, `vet`, `fence` rc 0; **`pkgs` rc 1** (1,084 s: `host/verbs` 0.2 s, `host/surfacepage` 8.0 s, **`cmd/ailang-worldd` ok 545 s**, `host/daemon` 518 s with one red, **`TestSeToolsRunBenchmarkTasksEndToEnd`** — directive-listed, red at the base too — **green alone on the final tree** (`logs/final-flaky-alone.txt`, 6.3 s); the base's other three load reds passed here); **`race` rc 0** (69 s); **`node` rc 0** (the sixteen named harnesses incl. the three pins); **`inv` rc 0: 54 named `--- PASS`** (the base's 39 + the twelve host/verbs and three cmd `TestOperator*`), 0 SKIP, 0 FAIL, no `NOT PASSED`; **`DAEMON_ROUTES: 47 (base 47)`, `SURFACE_POSTS: 3`, `UNTOUCHED_VIOLATIONS: none`, `CMD_SEAM_LINES: 10 0`, `PARSER_PRODUCTION_CALLERS: none (control: verbs.AdmitWhy production callers 1)`, `UNCHANGED:` app.js, surface_page.go, host/store/credentials.go**; **`chrome` rc 0** (12 s: **5 `--- PASS`**, none SKIP — the Needs-you pair 2.30 s / 1.92 s, XC2 2.58 s, XD1 2.34 s, the bar 1.66 s); **`GATES_DONE 1425s FAILS: pkgs`** (the one directive-listed load red, green alone). `verify_go` is the controller's |
| C18 | node, the pinned binaries, Chrome | `/opt/homebrew/bin/node --version`; `$HOME/.pinned-ailang/ailang --version` | `v26.13.0`; `AILANG v0.41.0`; `go1.26.6 darwin/arm64`; `Google Chrome 157.0.8103.79` (headless, the driver's temporary profile) |
| C19 | the CLI's bool spellings (the subset's must equal them) | `go doc strconv.ParseBool`; `go doc flag` (bool flags) | `1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False`; any other value → `invalid boolean value "<v>" for -<name>: parse error` — the subset's `boolWords` and error text are these (S1-D3, F3); `TestOperatorGrammarOptions` compares 18 spellings against `strconv.ParseBool` and `TestOperatorParsesLikeTheCLI` compares the error words with the CLI's own parse |

### 1.3 Planner decisions where the design or the directive leaves a choice

- **S1-D1 — no cut, by measurement.** M1 127, M2 +443, M3 +208, M4 +198: **976 ≤ 1,400**. Every §5.1 grammar and
  identifier row and all three §5.6 pins are in. No patch is parked under `plan-scratch/patches/`.
- **S1-D2 — four sub-milestones**, each independently green (C14): **M1** test-only, host/surfacepage; **M2** the parser with
  its first six grammar rows (its production lands with behavioural tests); **M3** the remaining rows, test-only; **M4** the seam
  and the three CLI joins. M2 is 443 lines (one production file and its first test file); splitting the parser from its tests
  would land untested production, so M2 stays one unit.
- **S1-D3 — the package is `host/verbs`, by measurement.** The shared seams are there (`Split`, `Lookup`, `ParseInterleaved`,
  `ApproveTarget`); host/daemon already imports it (`surface_verb.go`), so S7's route calls `AdmitOperator` without a new
  import edge, and `cmd/ailang-worldd` imports it, so the parity test can sit beside the CLI. `TestSurfaceVerbIsReadOnly` pins
  host/verbs' production imports to the pure set: `operator.go` imports `errors flag fmt io regexp` only (the first draft's
  `strconv` was refused by that test, F3, so the bool words are a table checked against `strconv.ParseBool` in the test). Design
  §6 places validators in `host/approval` — that is for the canonical action envelope (S2); the bar's grammar is host/verbs'.
- **S1-D4 — the API, designed so S7 need not change it.** `AdmitOperator(argv []string, p OperatorPolicy) OperatorAnswer`
  takes `Split`'s words (the route already splits); `OperatorAnswer{Class, Intent, Field, Reason, CLI}` with `Intent` set
  **exactly** when `Class == OperatorReview` (every test helper asserts it). The classes map one-to-one onto S7's outcomes:
  `review` → the NEW `review` outcome; `terminal` → the existing `terminal` panel (Reason ends `; at a terminal: <CLI>`);
  `usage` → `refused` with the CLI line; `never` → the never refusal; `other` → the existing Lookup path, unchanged (Reason and
  CLI are Lookup's, byte for byte). `Field` names the word a refusal is about (`<episode>`, `--preset`, `--headless`, …) for the
  page to point at. `Intent.Argv` is the canonical line (`ailang-worldd episode start <ep> --preset <p> --launch
  [--headless]`, `… finish <ep> --preset <p> [--confirm-on-surface] [--ready]`).
- **S1-D5 — the answers as an input.** `OperatorPolicy{Start, Finish Answer; HandoffProven bool}`: headless start iff `Start`
  ∈ {A, B}; interactive start (`--launch` without `--headless`) iff `Start == A && HandoffProven` (design §4.1 "under 119 A when
  E4 passes"); finish iff `Start` ∈ {A, B} and `Finish` ∈ {A, B} (design §10: S9–S11 wait for both); `--ready` iff additionally
  `Finish == A`. Only the exact strings `A`, `B`, `C` are answers (`a`, `D`, `A ` are unanswered). The zero value admits nothing.
- **S1-D6 — refusal classes per design §5.1 row.** Case of a flag or id, arity, a missing or empty value, an undefined or
  prefix/substring flag, a duplicate (even of mixed spellings), `--headless` without `--launch` (the CLI's own refusal,
  `episode.go:253`), `se-tools` as a launch preset, an invalid bool, `-h`/`--help` → **usage**. An unsupported CLI option, an
  explicit false, a start without `--launch` (an attended mint), a mode the answers deny → **terminal** with the CLI line. A
  case-mismatched or prefix/substring verb never reaches the subset: Lookup's classification stands (**other**, F5). The never
  families and `world-publish` → **never**.
- **S1-D7 — the join to the CLI's flag sets is a test seam, not a refactor or an AST read.** No function returns the CLI's flag
  sets and `--help` prints static text (C2, F2). Options measured: (a) move the 31 flag definitions into a function — a ≈ 60-line
  rewrite of two long CLI functions' locals; (b) read `fs.String("…")` calls by AST — static, and blind to flags registered
  by other means; (c) **a nil-by-default hook called with the built flag set just before `ParseInterleaved`** — 10 lines, no
  behaviour change when nil, and the test reads the real flag set (`VisitAll`, each flag's real `Value` kind) from a `--help`
  run that returns before any effect. (c) taken; the precedent for a package-level test hook in the CLI is `finishAfterStep`
  (`episode_finish.go:57`).
- **S1-D8 — the subset is stricter than the CLI where the design says so**: duplicates (the CLI's `flag` keeps the last) and
  explicit false are refused; everything the subset admits, the CLI parses to the same values (`TestOperatorParsesLikeTheCLI`).
- **S1-D9 — design §3.1's action id is defined, not used**: `OperatorActionIDPattern = ^oa-[0-9a-f]{32}$` (35 characters) and
  `ValidOperatorActionID`, pinned at every boundary class; no route, no `ApproveTarget` arm (that would change the served
  `approve` answer — S7/S8).
- **S1-D10 — "unregistered" is proven by a caller search with a positive control** in the gates (C11), plus the route count, the
  untouched-file list and `CMD_SEAM_LINES` (S1-D7's seam is the only cmd production change).
- **S1-D11 — no new Chrome test.** Design §5.6's "S1 pins queue rendering with one live smoke": S1 changes no `app.js` byte
  (C10), so the queue's live rendering is the existing `TestNeedsYouListsPendingRequests`, `TestNeedsYouNotifiesOnlyWhenAskedOnce`,
  XC2 `TestConfirmPageConfirmsWithAPasskeyInChrome`, `TestDecidePageApproveWithPasskey` and the bar's
  `TestVerbBarEpisodeReportInChrome`, run by the gates' `chrome` leg (`--- PASS` per name required). No sub-milestone's acceptance
  launches Chrome.
- **S1-D12 — mutants run against the named test AND the WHOLE package of the mutated file**; a host/verbs mutant that only the
  CLI join can see (a mirror entry the CLI lacks) also runs against `cmd/ailang-worldd` named and whole (§4).
- **S1-D13 — no docs change**: S1 adds no command, flag, output line or route; QUICKSTART changes when S7/S8 land the route and
  the card.

### 1.4 Premises found false or incomplete at HEAD (each routed around minimally)

| # | premise | measured | routing |
|---|---|---|---|
| F1 | directive: re-measure the design's `41331b4d` line numbers because #604 changed `app.js` | C1: none of 434, 456–458, 507 (nor 473, 499, 526) moved | none needed; cites stand |
| F2 | design §4.1: the subset is "joined to the CLI's flag definitions" | C2: the definitions are inline locals of two long functions with `fs.Usage = func() {}` and static help; nothing returns them | S1-D7: a 10-line nil-by-default seam in the two CLI functions — **the only production change outside host/verbs**, outside design §7's S1 row ("test-only shell changes and unregistered parser") |
| F3 | a host/verbs parser may use any stdlib package | `TestSurfaceVerbIsReadOnly` (host/daemon) refuses `strconv` in host/verbs' production files; the first draft (`strconv.ParseBool`) went red in M2's acceptance | `boolWords`, the flag package's twelve spellings; checked against `strconv.ParseBool` in the test (C19); the inventory is unedited |
| F4 | design §4.1 "the existing 1–64 ASCII grammar" | the grammar is `^[a-z0-9][a-z0-9-]{0,63}$` (lowercase letters, digits and `-`, first not `-`), not ASCII at large | the subset uses that exact expression; `TestOperatorIdentifierGrammarsMatchTheCLI` compares it with the CLI's on 20 boundary inputs |
| F5 | design §5.1 `TestOperatorGrammarCase` "refused, fixed grammar shown" for a verb's case | a case-mismatched verb never reaches the subset: `Episode …` is an address (Lookup step 5) and `episode Start …` is Lookup's `episode needs a subcommand: start \| launch \| finish \| report` (the fixed grammar) | `OperatorOther` carrying Lookup's reason; the test asserts Lookup's text |
| F6 | design §5.1 "standalone launch → terminal guidance" | `--headless` without `--launch` is refused by the CLI itself (`episode.go:253`); a start without `--launch` is the attended mint | the former is `usage`, the latter `terminal` (S1-D6) |
| F7 | design §5.1 "duplicate flag → usage refusal" reads as CLI-equivalent | the CLI accepts a repeated flag (last wins) | the subset is stricter (S1-D8); parity covers the spellings the subset admits |
| F8 | judge-405: the queue is safe for non-array `confirmations` | it is (C10); but an array **element** that is `null` throws in `confirmLineHTML`/`decisionLineHTML` and the pending map (`q.refused` of null) — measured: `node -e` over the base `app.js` with `{confirmations: [null]}`, `{decisions: [null]}`, `{pending: [null]}` → each `TypeError: Cannot read properties of null` (`logs/f8-null-element.txt`); unreachable from the daemon (typed slices of structs) | not routed in S1 (no app.js change; design §8 S1 is test-only); Appendix C carries it to S8's `TestOperatorWireTypes` |
| F9 | design §5.6 "S1 pins queue rendering with one live smoke" implies a new Chrome test | no app.js byte changes | S1-D11: the existing live tests in the gates' chrome leg |
| F10 | design §4.1: `world-publish` "falls outside those CLI word sets at HEAD. Add an explicit V4 refusal" | Lookup answers `world-publish …` as an address (step 5) | the subset answers `never` (pure; the served bar is unchanged until S7 renders the class) |

**Files outside design §7**: none outside its paths — §7 names `cmd/ailang-worldd` and `host/{…,verbs,…}`; S1 adds
`host/surfacepage/app_pins_test.go` (the page tests, as V3 S6's).

### 1.5 The size and the race-mode cost

| stage | production | test | changed lines | cumulative (numstat after commit) |
|---|---|---|---|---|
| M1 | 0 | 127 (`app_pins_test.go`) | 127 | **127** |
| M2 | 227 (`operator.go`) | 216 (`operator_test.go`) | 443 | **570** |
| M3 | 0 | 208 (`operator_policy_test.go`) | 208 | **778** |
| M4 | 10 (`episode.go` +7, `episode_finish.go` +3) | 188 (`operator_parity_test.go`) | 198 | **976** |

Design 620 raw (pins 210, subset and parity 410) → compiled **976 (1.57×)**: the pins came in under (127 for 210: no
production change was needed), the subset over (849 for 410: the boundary-class matrix the directive requires).

**Race-mode wall time of the new tests** (directive 8; work tree at load ≈ 27, `go test -race -count=1 -v`,
`logs/race-new.txt`): `host/verbs` — the twelve `TestOperator*` 0.07 s together (package 1.7 s); `host/surfacepage` — the three
pins 0.26 s (node subprocesses); `cmd/ailang-worldd` — the three joins 0.01 s. **≈ 0.34 s added in all**; CI's 8-minute bound
per package is untouched (host/daemon and host/broker gain nothing: no S1 test there). No sleep anywhere in S1's tests.

---

## 2. How the executor applies a sub-milestone (the patch protocol)

Each sub-milestone `<M>` is the block fenced by a line `` ````diff <M> `` (four backticks) and a line of four backticks,
in Appendix P. Paths are `a/<repo path>` and `b/<repo path>`; a new file's patch has `--- /dev/null`; **no
sub-milestone deletes a file**. Apply it **only** this way, in one `bash` shell (the controller's shell is zsh):

```bash
bash <<'EOF2'
set -euo pipefail
M=M1   # <- the sub-milestone, exactly as in the Appendix P heading, in order: M1 M2 M3 M4
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter407-row243-v4s1
PLAN="$W/design_docs/planned/w-world-address-bar-verbs-v4s1-sprint-plan.md"
T="$HOME/.ailang/state/world-iter407/exec-patches"; mkdir -p "$T"
awk -v m="$M" '$0 == "````diff " m {f=1; next} f && $0 == "````" {exit} f' "$PLAN" > "$T/$M.patch"
test -s "$T/$M.patch"
cd "$W"
/usr/bin/patch -p1 -E --dry-run < "$T/$M.patch"
/usr/bin/patch -p1 -E < "$T/$M.patch"
if find cmd host \( -name '*.orig' -o -name '*.rej' \) 2>/dev/null | grep .; then echo "STOP: patch left .orig/.rej"; exit 1; fi
git status --short --untracked-files=all
EOF2
```

- If `--dry-run` fails, **stop and report** (no hand edits, no fuzz options). It can only fail if an earlier
  sub-milestone was skipped, reordered or changed.
- Then run that sub-milestone's **Acceptance** (§3): extract `accept.sh` from this plan (`awk '$0 == "````bash
  accept.sh" {f=1; next} f && $0 == "````" {exit} f' "$PLAN" > "$T/accept.sh"`) and run **each part §3 lists for M, one
  per Bash call**: `bash "$T/accept.sh" "$W" <M> <part> 2>&1 | tee
  "$HOME/.ailang/state/world-iter407/exec-logs/accept_<M>_<part>.txt"` (never only through `tail`). **Never `nohup`,
  never background a test leg; pass a 600000 ms timeout to every long Bash call.** No part launches Chrome.
- Then snapshot (§7). The controller commits and compares the numstat with §3's expected value.

---

## 3. Milestones

**Acceptance.** Every part runs the same `[base]` checks first (gofmt, build, vet of the three touched packages, the
test-file compile fence `go test -run '^$' ./...`). `v` requires every named test to `--- PASS` and none to SKIP or FAIL;
`p` is a whole-package run. Parts: `fast` (≈ 0.5–1.5 min) for every sub-milestone; `cmd1` + `cmd2` + `cmd3`
(`cmd/ailang-worldd` whole, in three parts by test name, `^Test[A-D]`, `^TestE`, `^Test[^A-E]`) for M4 only — measured in
C14, each one 600 s call (two halves were measured first: `^Test[A-E]` took 607 s at load ≈ 20, over the call).
`host/verbs` and `host/surfacepage` run whole in `fast` (≈ 1–3 s each).

````bash accept.sh
#!/bin/bash
# accept.sh <tree> <M> [fast|cmd1|cmd2|cmd3]: exactly plan §3's acceptance for sub-milestone M (one bash shell;
# rc 0 = ACCEPT_OK). Each part fits one 600 s call; run every part §3 lists for M, in order. No part
# launches Chrome (S1 adds no Chrome test; the gates' chrome leg is the controller's).
T=$1; M=$2; P=${3:-fast}
cd "$T" || exit 2
export AILANG_BIN=$HOME/.pinned-ailang/ailang
export WORLD_EXEC_NODE=/opt/homebrew/bin/node
export WORLD_TOOL_AILANG_BIN=$HOME/.pinned-ailang-tools/v0.54.15/ailang
unset WORLD_UPDATE_VERB_GOLDEN WORLD_UPDATE_CONFIRM_GOLDEN WORLD_UPDATE_CONFIRMATIONS_GOLDEN WORLD_UPDATE_FINISH_PROMPT_GOLDEN WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER
set -o pipefail
chk() { tee /dev/stderr | { out=$(cat); echo "$out" | grep -qE -- '--- (FAIL|SKIP)|^FAIL' && exit 1; echo "$out" | grep -q -- '--- PASS' || exit 1; }; }
v() { go test -count=1 -v "$@" 2>&1 | chk; }   # -v: every named test PASS, none SKIP/FAIL
p() { go test -count=1 "$@"; }
step() { echo "== $*"; "$@" || { echo "ACCEPT_RED $M $P: $*"; exit 1; }; }
# [base]
step test -z "$(gofmt -l cmd host)"
step go build ./...
step go vet ./cmd/ailang-worldd/ ./host/surfacepage/ ./host/verbs/
if go test -run '^$' ./... 2>&1 | grep -E '^(FAIL|--- FAIL)'; then echo "ACCEPT_RED $M $P: FENCE_RED"; exit 1; fi
PINS='TestNeedsYouConfirmationEmbeddedNewline|TestNeedsYouConfirmationsNonArray|TestVerbDecidePrefixBoundary'
NODE='TestSurfaceAppPure|TestSurfaceAppVerbs|TestNeedsYouListsDecisions|TestNeedsYouListsConfirmations|TestAddressHintsArePinned|TestSurfaceAppMatchesWorkingCopy'
VOLD='TestLookup|TestApproveTarget|TestSplit|TestParseInterleaved|TestEpisodeReportFlags|TestAdmitWhy|TestNeverAndOtherReasonsEndWithTheCLILine'
G2='TestOperatorGrammarCase|TestOperatorGrammarArity|TestOperatorGrammarExactWords|TestOperatorGrammarCanonical|TestOperatorGrammarDuplicates|TestOperatorGrammarOptions'
G3='TestOperatorGrammarPolicy|TestOperatorNeverDispatch|TestOperatorIdentifierAnchors|TestOperatorIdentifierFlags|TestOperatorIdentifierLength|TestOperatorOtherVerbParity'
C4='TestOperatorFlagParity|TestOperatorParsesLikeTheCLI|TestOperatorIdentifierGrammarsMatchTheCLI'
INV='TestEveryVerbHasHelp|TestVerbTableMatchesTheCLISwitch|TestDoctorEpisodeGrammarMatchesTheDaemon|TestCLICredentialWriterSites|TestEpisodeFinishHelp'
case "$M:$P" in
M1:fast) step v ./host/surfacepage/ -run "^($PINS|$NODE)\$"
         step p ./host/surfacepage/ ;;
M2:fast) step v ./host/verbs/ -run "^($G2|$VOLD)\$"
         step p ./host/verbs/
         step v ./host/daemon/ -run '^(TestSurfaceVerbOutcomes|TestSurfaceVerbRouteRegistered|TestSurfaceVerbIsReadOnly|TestDaemonDependencyAllowlist)$'
         step v ./cmd/ailang-worldd/ -run "^($INV)\$" ;;
M3:fast) step v ./host/verbs/ -run "^($G2|$G3|$VOLD)\$"
         step p ./host/verbs/ ;;
M4:fast) step v ./cmd/ailang-worldd/ -run "^($C4|$INV)\$"
         step v ./host/verbs/ -run "^($G2|$G3)\$"
         step p ./host/verbs/ ./host/surfacepage/ ;;
M4:cmd1) step p -timeout 30m ./cmd/ailang-worldd/ -run '^Test[A-D]' ;;    # the three thirds are the whole package
M4:cmd2) step p -timeout 30m ./cmd/ailang-worldd/ -run '^TestE' ;;
M4:cmd3) step p -timeout 30m ./cmd/ailang-worldd/ -run '^Test[^A-E]' ;;
*) echo "unknown $M $P"; exit 2 ;;
esac
echo "ACCEPT_OK $M $P"
````

### M1 — the judge-405 pins · 0 prod + 127 test · part: `fast`

- `host/surfacepage/app_pins_test.go` (NEW): `runPinHarness` (CP7's form: `AppScript`, the two goldens, `WORLD_EXEC_NODE`
  non-skippable, 30 s context, a manifest line) and a shared prelude (the handler's summarised draft row `vr-00000000000000c4`,
  a real-shaped decision and its exact link line).
  - `TestNeedsYouConfirmationEmbeddedNewline`: the valid row and the decision both link (control); then the id with `\n`, `\r`,
    `\r\n`, U+2028 or U+2029 after or before it, upper-cased whole, `VR`-prefixed, or upper-case hex → exactly the decision line
    plus `Confirm · <esc(id)> · refused: not a confirmation id`, and no `/surface/confirm/` anywhere (kills `J-JS-RE-NL`,
    `P1-RE-CARET`, `P1-RE-DOLLAR`, `P1-RE-I`, `P1-RE-MI` named; `P1-RE-S` named equivalent).
  - `TestNeedsYouConfirmationsNonArray`: the mint request + decision answer with `confirmations: []` is the expected HTML; then
    `{}`, `{0: row, length: 1}`, a string id, `''`, `null`, `0`, `1`, `true`, `{map: () => [a link]}` and the key absent each
    render exactly that, and the real array still counts 5 (kills `J-JS-CS-NOGUARD`, `P2-CS-*`).
  - `TestVerbDecidePrefixBoundary`: the daemon's own `approve sha256:…`, `approve mr-…`, `approve vr-…` answers navigate
    (control); 15 hrefs refused: `//evil/…`, `https://evil/…`, `/x/…`, `/surface/approve/x/surface/decide/…`, a leading space,
    `sha256` without the colon (two forms), `sha256/`, `SHA256:`, `/Surface/`, `mr` without the dash (two forms), `MR-`,
    `/x/surface/approve/…`, `/SURFACE/confirm/…` (kills `J-JS-ALLOW-DECIDE`, `P3-*`).
- Expected numstat after commit **127**.

### M2 — the admitted subset and its grammar rows · 227 prod + 216 test · part: `fast`

- `host/verbs/operator.go` (NEW, S1-D3 … D9): the mirror (`EpisodeStartFlags` 22, `EpisodeFinishFlags` 9, each `{Name, Bool,
  Admit}`, the CLI's order), the grammars, `Answer`/`OperatorPolicy` and its four predicates, the classes, `OperatorIntent`,
  `OperatorAnswer`, `operatorValue` (counts each flag's occurrences; a bool parses through `boolWords` with the flag package's
  `parse error`), and `AdmitOperator`: `world-publish` → never; `Lookup`; not start/finish → never/other with Lookup's answer;
  a fresh `flag.FlagSet` of counting values over the mirror; `ParseInterleaved`; help → usage; parse error → usage (the flag
  package's words); duplicates → usage; an unadmitted flag or an explicit false → terminal; exactly one episode in
  `EpisodeIDPattern`; a non-empty, non-`se-tools` preset in the grammar and ≤ 51; start: `--headless` without `--launch` →
  usage, no `--launch` → terminal, the mode by the policy → terminal; finish: the policy, `--ready` by the policy; the canonical
  argv.
- `host/verbs/operator_test.go` (NEW): the helpers (`admitLine` over `Split`, `admitArgv` asserting Intent ⇔ review,
  `refused` asserting class, field, reason **and no tuple**) and `TestOperatorGrammarCase`, `…Arity`, `…ExactWords`,
  `…Canonical` (8 + 2 + 3 + 2 + 2 + 1 spellings to six exact tuples; the canonical argv re-admits to the same tuple; `CLI` is
  Lookup's), `…Duplicates`, `…Options` (every unadmitted flag of each verb, in three spellings, by an independent literal list;
  six explicit falses; standalone forms; `se-tools`; `-h`/`--help`/`-help`; 18 bool spellings against `strconv.ParseBool`).
- Expected numstat after commit **570**.

### M3 — policy, never-dispatch, identifiers, the other verbs · 0 prod + 208 test · part: `fast`

- `host/verbs/operator_policy_test.go` (NEW): `TestOperatorGrammarPolicy` (8 Start × 6 Finish answers × HandoffProven × 4
  lines = 384 cases against `policyOracle`, the design's table written out: 42 admitted; every refusal terminal with its field;
  the zero value refuses all four), `TestOperatorNeverDispatch` (12 never-family lines incl. `world-publish` → never; 10 other
  lines keep Lookup's reason; shell metacharacters never reach the parser — `Split` answers `ErrShell`),
  `TestOperatorIdentifierAnchors`, `…Flags`, `…Length` (episode, preset and action id: prefix/suffix junk, line breaks either
  side, case; episode 1/64 in, 0/65 out; preset 51 in, 52/64/65 out; action id 32 hex in, 30/31/33/34 and `oa-` alone out), and
  `TestOperatorOtherVerbParity` (every `verb-route.golden.json` line and every word bare and with an argument: ≥ 70 lines keep
  Lookup's class, reason and CLI; the golden's own `episode finish w1 --preset p` is a review under all-A).
- Expected numstat after commit **778**.

### M4 — the join to the CLI · 10 prod + 188 test · parts: `fast`, `cmd1`, `cmd2`, `cmd3`

- `cmd/ailang-worldd/episode.go`: `var episodeFlagSetSeen func(fs *flag.FlagSet)` (with its comment, before `runEpisodeStart`)
  and `if episodeFlagSetSeen != nil { episodeFlagSetSeen(fs) }` just before `verbs.ParseInterleaved` (+7);
  `episode_finish.go`: the same call before its `ParseInterleaved` (+3). Nil in production: no behaviour change.
- `cmd/ailang-worldd/operator_parity_test.go` (NEW): `cliEpisodeFlagSet` (sets the seam, runs `episode <verb> --help` with a
  temp HOME, requires rc 0, exactly one flag set seen, help on stdout, nothing on stderr); **`TestOperatorFlagParity`** (the CLI's
  `VisitAll` names and bool kinds equal host/verbs' mirror exactly, 22 and 9; the admitted sets are design §4.1's by literal; four
  controls — one more, one fewer, renamed, kind flipped — each differ); **`TestOperatorParsesLikeTheCLI`** (seven admitted
  spellings and their canonical lines, parsed by a fresh CLI flag set with `ParseInterleaved`, give the tuple's episode and
  values; five lines the CLI's parse rejects are usage refusals ending with the CLI's own error); 
  **`TestOperatorIdentifierGrammarsMatchTheCLI`** (`EpisodeIDPattern` decides as `workspaceEpisodePattern` on 20 boundary
  inputs, 7 matching; `MaxPresetName == maxPresetNameLen`; `loadLaunchPreset` refuses 52 characters as invalid and reads past the
  grammar at 51, and the subset agrees).
- Expected numstat after commit **976**.

### 3.1 Every design §5.1/§5.6 row S1 drives, traced (directive retro 4)

| design row | branch | test (and the no-tuple assertion) |
|---|---|---|
| §5.6 embedded newline | `confirmLineHTML`'s regex, `m`/`i`/anchors | `TestNeedsYouConfirmationEmbeddedNewline` (no `/surface/confirm/` in the HTML) |
| §5.6 non-array confirmations | `needsYouHTML`'s `Array.isArray` | `TestNeedsYouConfirmationsNonArray` (HTML equals the no-confirmations answer) |
| §5.6 decide prefix | `verbOutcome`'s navigate allowlist | `TestVerbDecidePrefixBoundary` (`{say: 'refused: unexpected target'}`, no `assign`) |
| §5.1 existing address, report, why, approve | default arm → Lookup's | `TestOperatorOtherVerbParity` (class other/never, Lookup's reason and CLI, Intent nil) |
| §5.1 valid supported start / finish (the tuple only; the panel is S7/S8) | review | `TestOperatorGrammarCanonical`, `TestOperatorGrammarPolicy` |
| case mismatch (verb, flag, id) | Lookup / parse error / grammar | `TestOperatorGrammarCase` |
| extra positional, missing episode/preset/value | arity, empty preset, parse error | `TestOperatorGrammarArity` |
| prefix or substring of a root or flag | Lookup's exact rows; flag's exact names | `TestOperatorGrammarExactWords` |
| quoted/equals/interleaved | `Split` + `ParseInterleaved` → one tuple | `TestOperatorGrammarCanonical`, `TestOperatorParsesLikeTheCLI` |
| duplicate flag, mixed spellings | `operatorValue.n > 1` | `TestOperatorGrammarDuplicates` |
| unsupported option, explicit false, standalone launch | unadmitted / false / mode | `TestOperatorGrammarOptions` |
| denied authority mode (119 B interactive, 120 B `--ready`, unanswered) | the four predicates | `TestOperatorGrammarPolicy` |
| enrol/revoke/mint/service/publish, arbitrary shell | never arm, world-publish, `Split` | `TestOperatorNeverDispatch` |
| id with prefix/suffix junk | the anchors of three grammars | `TestOperatorIdentifierAnchors` |
| multiline id or case mismatch | `(?m)`, `(?i)` | `TestOperatorIdentifierFlags` |
| identifier just outside its length | `{0,63}`, the 51 cap, `{32}` | `TestOperatorIdentifierLength` |
| the CLI adds, renames, removes a flag; the subset names one the CLI lacks | the mirror vs the real flag set | `TestOperatorFlagParity` |
| the subset parses an admitted line differently from the CLI | `operatorValue`, `boolWords` | `TestOperatorParsesLikeTheCLI` |
| the grammars drift from the CLI's | `EpisodeIDPattern`, `MaxPresetName` | `TestOperatorIdentifierGrammarsMatchTheCLI` |

**Not S1** (Appendix C names each owner): route, body, Origin/Host, media type, method, runner readiness, catalogue, admission,
duplicate tap, reconnect, panel network/latest, literal task, strict JSON, wire types, request limits, link allowlist.

---

## 4. Mutation table

Enumerated **from the diff** (`tools/mutants.py`, the ids below are the ledger's): for `app.js`, the three pinned branches at every boundary class (each regex anchor, `m`, `i`, `s`, `mi`; input type — object, length-only, has-`map`, not-undefined; prefix vs substring — `includes`, `indexOf`, `endsWith`; the colon delimiter; case; a trim), and the approve and confirm siblings; for `operator.go`, **every branch of the diff**: each unadmitted flag admitted (25), each admitted flag dropped (5), kind flips, a mirror entry dropped, renamed or added, the counting value (count, bool, string, IsBoolFlag ×2, the bool words ×4, the error words), help, parse error, duplicates, the unadmitted and explicit-false checks, arity ×2, the episode check, the preset checks (empty, `se-tools`, grammar, `>=`, `+1`, cap 52/50), both grammars at every boundary class (anchors, `(?m)`, `(?i)`, `(?s)`, length up/down, lower bound, leading dash, the delimiter), `world-publish` ×2, each never family inserted into the start arm (10), the never class, the other arm's reason and CLI, the start/finish mirrors swapped, the terminal suffix, the reason prefix, the four mode cases, every answer comparison of the four policy predicates (9), the finish policy ×3, and the canonical tuple and argv (7); for the CLI, the seam (off ×2, after the parse), a flag added to either verb, removed, renamed, the preset cap and the episode grammar. Each is one exact replacement (`tools/mut.py`; text not found exactly once, or a mutant that does not compile, is **NM**). Each runs against its **named** test, then against the **WHOLE package** of the mutated file (`-run '.'`: every test, static included); a host/verbs mutant whose kill belongs to the CLI join also runs against `cmd/ailang-worldd` (named / whole). No S1 test is static (no inventory, AST or string-grep test is new), so no kill needs a `w/o static` leg; the pre-existing inventories that killed a mutant are named in the cell. Ledger: `$HOME/.ailang/state/world-iter407/plan-scratch/mutants/ledger.txt` (one line per run; the table shows each id's last run per leg; earlier drafts' runs — e.g. `MIR-KIND-ttl`, `VAL-BOOL-ERRTEXT` and `REFUSE-PREFIX` SURVIVED the first `operator_test.go`, and `P1-RE-I` the first newline harness, each killed after a row was added — stay in the ledger). The cmd-side `whole` legs of `SEAM-*`/`CLI-*` ran on the M4 draft whose only difference from the final tree is host/verbs' bool-word table and test files (cmd's files and tests byte-identical); their named legs were re-run on the final tree.

**Base reproductions (C9; tag `base-whole`, on a clone at `d92983d5`)**:

| id | file | mutation | command | verdict |
|---|---|---|---|---|
| `J-JS-RE-NL` | `app.js` | `/^vr-[0-9a-f]{16}$/.test(c.id) -> /^vr-[0-9a-f]{16}$/m.test(c.id)` | `go test -count=1 -timeout 900s ./host/surfacepage/ -run '.'` | **SURVIVED**  |
| `P1-RE-CARET` | `app.js` | `/^vr-[0-9a-f]{16}$/.test(c.id) -> /vr-[0-9a-f]{16}$/.test(c.id)` | `go test -count=1 -timeout 900s ./host/surfacepage/ -run '.'` | **KILLED** TestNeedsYouListsConfirmations |
| `P1-RE-DOLLAR` | `app.js` | `/^vr-[0-9a-f]{16}$/.test(c.id) -> /^vr-[0-9a-f]{16}/.test(c.id)` | `go test -count=1 -timeout 900s ./host/surfacepage/ -run '.'` | **KILLED** TestNeedsYouListsConfirmations |
| `P1-RE-I` | `app.js` | `/^vr-[0-9a-f]{16}$/.test(c.id) -> /^vr-[0-9a-f]{16}$/i.test(c.id)` | `go test -count=1 -timeout 900s ./host/surfacepage/ -run '.'` | **KILLED** TestNeedsYouListsConfirmations |
| `P1-RE-S` | `app.js` | `/^vr-[0-9a-f]{16}$/.test(c.id) -> /^vr-[0-9a-f]{16}$/s.test(c.id)` | `go test -count=1 -timeout 900s ./host/surfacepage/ -run '.'` | **SURVIVED**  |
| `J-JS-CS-NOGUARD` | `app.js` | `const cs = r && Array.isArray(r.confirmations) ? r.confirmations : []; -> const cs = r && r.confirmations ? r.` | `go test -count=1 -timeout 900s ./host/surfacepage/ -run '.'` | **SURVIVED**  |
| `J-JS-ALLOW-DECIDE` | `app.js` | `b.href.startsWith('/surface/decide/sha256:') -> b.href.includes('/surface/decide/sha256:')` | `go test -count=1 -timeout 900s ./host/surfacepage/ -run '.'` | **SURVIVED**  |
| `P3-DECIDE-INDEXOF` | `app.js` | `b.href.startsWith('/surface/decide/sha256:') -> b.href.indexOf('/surface/decide/sha256:') >= 0` | `go test -count=1 -timeout 900s ./host/surfacepage/ -run '.'` | **SURVIVED**  |
| `P3-DECIDE-NOCOLON` | `app.js` | `b.href.startsWith('/surface/decide/sha256:') -> b.href.startsWith('/surface/decide/sha256')` | `go test -count=1 -timeout 900s ./host/surfacepage/ -run '.'` | **SURVIVED**  |
| `P3-DECIDE-CASE` | `app.js` | `b.href.startsWith('/surface/decide/sha256:') -> b.href.toLowerCase().startsWith('/surface/decide/sha256:')` | `go test -count=1 -timeout 900s ./host/surfacepage/ -run '.'` | **SURVIVED**  |
| `P3-APPROVE-NODASH` | `app.js` | `b.href.startsWith('/surface/approve/mr-') -> b.href.startsWith('/surface/approve/mr')` | `go test -count=1 -timeout 900s ./host/surfacepage/ -run '.'` | **SURVIVED**  |
| `P3-APPROVE-CASE` | `app.js` | `b.href.startsWith('/surface/approve/mr-') -> b.href.toLowerCase().startsWith('/surface/approve/mr-')` | `go test -count=1 -timeout 900s ./host/surfacepage/ -run '.'` | **SURVIVED**  |

| id | milestone · file | mutation | named test | named | whole package (killing tests) | cmd join named / whole | class |
|---|---|---|---|---|---|---|---|
| `J-JS-RE-NL` | M1 · `app.js` | `/^vr-[0-9a-f]{16}$/.test(c.id) -> /^vr-[0-9a-f]{16}$/m.test(c.id)` | `^TestNeedsYouConfirmationEmbeddedNewline$` | killed (TestNeedsYouConfirmationEmbeddedNewline) | killed (TestNeedsYouConfirmationEmbeddedNewline) | — | K |
| `P1-RE-CARET` | M1 · `app.js` | `/^vr-[0-9a-f]{16}$/.test(c.id) -> /vr-[0-9a-f]{16}$/.test(c.id)` | `^TestNeedsYouConfirmationEmbeddedNewline$` | killed (TestNeedsYouConfirmationEmbeddedNewline) | killed (TestNeedsYouConfirmationEmbeddedNewline,TestNeedsYouListsConfirmations) | — | K |
| `P1-RE-DOLLAR` | M1 · `app.js` | `/^vr-[0-9a-f]{16}$/.test(c.id) -> /^vr-[0-9a-f]{16}/.test(c.id)` | `^TestNeedsYouConfirmationEmbeddedNewline$` | killed (TestNeedsYouConfirmationEmbeddedNewline) | killed (TestNeedsYouConfirmationEmbeddedNewline,TestNeedsYouListsConfirmations) | — | K |
| `P1-RE-I` | M1 · `app.js` | `/^vr-[0-9a-f]{16}$/.test(c.id) -> /^vr-[0-9a-f]{16}$/i.test(c.id)` | `^TestNeedsYouConfirmationEmbeddedNewline$` | killed (TestNeedsYouConfirmationEmbeddedNewline) | killed (TestNeedsYouConfirmationEmbeddedNewline,TestNeedsYouListsConfirmations) | — | K |
| `P1-RE-S` | M1 · `app.js` | `/^vr-[0-9a-f]{16}$/.test(c.id) -> /^vr-[0-9a-f]{16}$/s.test(c.id)` | `^TestNeedsYouConfirmationEmbeddedNewline$` | survived (-) | survived (-) | — | E — the expression has no `.`; `s` changes only what `.` matches |
| `P1-RE-MI` | M1 · `app.js` | `/^vr-[0-9a-f]{16}$/.test(c.id) -> /^vr-[0-9a-f]{16}$/mi.test(c.id)` | `^TestNeedsYouConfirmationEmbeddedNewline$` | killed (TestNeedsYouConfirmationEmbeddedNewline) | killed (TestNeedsYouConfirmationEmbeddedNewline,TestNeedsYouListsConfirmations) | — | K |
| `J-JS-CS-NOGUARD` | M1 · `app.js` | `const cs = r && Array.isArray(r.confirmations) ? r.confirmations : []; -> const cs = r && r.confirmations ? r.confirmati` | `^TestNeedsYouConfirmationsNonArray$` | killed (TestNeedsYouConfirmationsNonArray) | killed (TestNeedsYouConfirmationsNonArray) | — | K |
| `P2-CS-OBJECT` | M1 · `app.js` | `const cs = r && Array.isArray(r.confirmations) ? r.confirmations : []; -> const cs = r && typeof r.confirmations === "ob` | `^TestNeedsYouConfirmationsNonArray$` | killed (TestNeedsYouConfirmationsNonArray) | killed (TestNeedsYouConfirmationsNonArray) | — | K |
| `P2-CS-NOTUNDEF` | M1 · `app.js` | `const cs = r && Array.isArray(r.confirmations) ? r.confirmations : []; -> const cs = r && r.confirmations !== undefined ` | `^TestNeedsYouConfirmationsNonArray$` | killed (TestNeedsYouConfirmationsNonArray) | killed (TestNeedsYouConfirmationsNonArray) | — | K |
| `P2-CS-LENGTH` | M1 · `app.js` | `const cs = r && Array.isArray(r.confirmations) ? r.confirmations : []; -> const cs = r && r.confirmations && r.confirmat` | `^TestNeedsYouConfirmationsNonArray$` | killed (TestNeedsYouConfirmationsNonArray) | killed (TestNeedsYouConfirmationsNonArray) | — | K |
| `P2-CS-HASMAP` | M1 · `app.js` | `const cs = r && Array.isArray(r.confirmations) ? r.confirmations : []; -> const cs = r && r.confirmations && typeof r.co` | `^TestNeedsYouConfirmationsNonArray$` | killed (TestNeedsYouConfirmationsNonArray) | killed (TestNeedsYouConfirmationsNonArray) | — | K |
| `J-JS-ALLOW-DECIDE` | M1 · `app.js` | `b.href.startsWith('/surface/decide/sha256:') -> b.href.includes('/surface/decide/sha256:')` | `^TestVerbDecidePrefixBoundary$` | killed (TestVerbDecidePrefixBoundary) | killed (TestVerbDecidePrefixBoundary) | — | K |
| `P3-DECIDE-INDEXOF` | M1 · `app.js` | `b.href.startsWith('/surface/decide/sha256:') -> b.href.indexOf('/surface/decide/sha256:') >= 0` | `^TestVerbDecidePrefixBoundary$` | killed (TestVerbDecidePrefixBoundary) | killed (TestVerbDecidePrefixBoundary) | — | K |
| `P3-DECIDE-NOCOLON` | M1 · `app.js` | `b.href.startsWith('/surface/decide/sha256:') -> b.href.startsWith('/surface/decide/sha256')` | `^TestVerbDecidePrefixBoundary$` | killed (TestVerbDecidePrefixBoundary) | killed (TestVerbDecidePrefixBoundary) | — | K |
| `P3-DECIDE-CASE` | M1 · `app.js` | `b.href.startsWith('/surface/decide/sha256:') -> b.href.toLowerCase().startsWith('/surface/decide/sha256:')` | `^TestVerbDecidePrefixBoundary$` | killed (TestVerbDecidePrefixBoundary) | killed (TestVerbDecidePrefixBoundary) | — | K |
| `P3-DECIDE-TRIM` | M1 · `app.js` | `b.href.startsWith('/surface/decide/sha256:') -> b.href.trim().startsWith('/surface/decide/sha256:')` | `^TestVerbDecidePrefixBoundary$` | killed (TestVerbDecidePrefixBoundary) | killed (TestVerbDecidePrefixBoundary) | — | K |
| `P3-DECIDE-ENDS` | M1 · `app.js` | `b.href.startsWith('/surface/decide/sha256:') -> b.href.endsWith(b.href.slice(-64)) && b.href.includes('decide/sha256:')` | `^TestVerbDecidePrefixBoundary$` | killed (TestVerbDecidePrefixBoundary) | killed (TestVerbDecidePrefixBoundary) | — | K |
| `P3-APPROVE-NODASH` | M1 · `app.js` | `b.href.startsWith('/surface/approve/mr-') -> b.href.startsWith('/surface/approve/mr')` | `^TestVerbDecidePrefixBoundary$` | killed (TestVerbDecidePrefixBoundary) | killed (TestVerbDecidePrefixBoundary) | — | K |
| `P3-APPROVE-CASE` | M1 · `app.js` | `b.href.startsWith('/surface/approve/mr-') -> b.href.toLowerCase().startsWith('/surface/approve/mr-')` | `^TestVerbDecidePrefixBoundary$` | killed (TestVerbDecidePrefixBoundary) | killed (TestVerbDecidePrefixBoundary) | — | K |
| `P3-APPROVE-INCLUDES` | M1 · `app.js` | `b.href.startsWith('/surface/approve/mr-') -> b.href.includes('/surface/approve/mr-')` | `^TestVerbDecidePrefixBoundary$` | killed (TestVerbDecidePrefixBoundary) | killed (TestSurfaceAppVerbs,TestVerbDecidePrefixBoundary) | — | K |
| `P3-CONFIRM-CASE` | M1 · `app.js` | `b.href.startsWith('/surface/confirm/vr-') -> b.href.toLowerCase().startsWith('/surface/confirm/vr-')` | `^TestVerbDecidePrefixBoundary$` | killed (TestVerbDecidePrefixBoundary) | killed (TestSurfaceAppVerbs,TestVerbDecidePrefixBoundary) | — | K |
| `OPT-S-db` | M2 · `operator.go` | `{"db", false, false} -> {"db", false, true}` | `^TestOperatorGrammarOptions$` | killed (TestOperatorGrammarOptions) | killed (TestOperatorGrammarOptions) | — | K |
| `OPT-S-credentials-db` | M2 · `operator.go` | `{"credentials-db", false, false} -> {"credentials-db", false, true}` | `^TestOperatorGrammarOptions$` | killed (TestOperatorGrammarOptions) | killed (TestOperatorGrammarOptions) | — | K |
| `OPT-S-workspace-root` | M2 · `operator.go` | `{"workspace-root", false, false} -> {"workspace-root", false, true}` | `^TestOperatorGrammarOptions$` | killed (TestOperatorGrammarOptions) | killed (TestOperatorGrammarOptions) | — | K |
| `OPT-S-repo` | M2 · `operator.go` | `{"repo", false, false} -> {"repo", false, true}` | `^TestOperatorGrammarOptions$` | killed (TestOperatorGrammarOptions) | killed (TestOperatorGrammarOptions) | — | K |
| `OPT-S-branch` | M2 · `operator.go` | `{"branch", false, false} -> {"branch", false, true}` | `^TestOperatorGrammarOptions$` | killed (TestOperatorGrammarOptions) | killed (TestOperatorGrammarOptions) | — | K |
| `OPT-S-profile` | M2 · `operator.go` | `{"profile", false, false} -> {"profile", false, true}` | `^TestOperatorGrammarOptions$` | killed (TestOperatorGrammarOptions) | killed (TestOperatorGrammarOptions) | — | K |
| `OPT-S-task` | M2 · `operator.go` | `{"task", false, false} -> {"task", false, true}` | `^TestOperatorGrammarOptions$` | killed (TestOperatorGrammarOptions) | killed (TestOperatorGrammarOptions) | — | K |
| `OPT-S-agent` | M2 · `operator.go` | `{"agent", false, false} -> {"agent", false, true}` | `^TestOperatorGrammarOptions$` | killed (TestOperatorGrammarOptions) | killed (TestOperatorGrammarOptions) | — | K |
| `OPT-S-model` | M2 · `operator.go` | `{"model", false, false} -> {"model", false, true}` | `^TestOperatorGrammarOptions$` | killed (TestOperatorGrammarOptions) | killed (TestOperatorGrammarOptions) | — | K |
| `OPT-S-ttl` | M2 · `operator.go` | `{"ttl", false, false} -> {"ttl", false, true}` | `^TestOperatorGrammarOptions$` | killed (TestOperatorGrammarOptions) | killed (TestOperatorGrammarOptions) | — | K |
| `OPT-S-resume` | M2 · `operator.go` | `{"resume", false, false} -> {"resume", false, true}` | `^TestOperatorGrammarOptions$` | killed (TestOperatorGrammarOptions) | killed (TestOperatorGrammarOptions) | — | K |
| `OPT-S-hosts-dir` | M2 · `operator.go` | `{"hosts-dir", false, false} -> {"hosts-dir", false, true}` | `^TestOperatorGrammarOptions$` | killed (TestOperatorGrammarOptions) | killed (TestOperatorGrammarOptions) | — | K |
| `OPT-S-grant` | M2 · `operator.go` | `{"grant", false, false} -> {"grant", false, true}` | `^TestOperatorGrammarOptions$` | killed (TestOperatorGrammarOptions) | killed (TestOperatorGrammarOptions) | — | K |
| `OPT-S-budget` | M2 · `operator.go` | `{"budget", false, false} -> {"budget", false, true}` | `^TestOperatorGrammarOptions$` | killed (TestOperatorGrammarOptions) | killed (TestOperatorGrammarOptions) | — | K |
| `OPT-S-sessions-dir` | M2 · `operator.go` | `{"sessions-dir", false, false} -> {"sessions-dir", false, true}` | `^TestOperatorGrammarOptions$` | killed (TestOperatorGrammarOptions) | killed (TestOperatorGrammarOptions) | — | K |
| `OPT-S-addr` | M2 · `operator.go` | `{"addr", false, false} -> {"addr", false, true}` | `^TestOperatorGrammarOptions$` | killed (TestOperatorGrammarOptions) | killed (TestOperatorGrammarOptions) | — | K |
| `OPT-S-request` | M2 · `operator.go` | `{"request", true, false} -> {"request", true, true}` | `^TestOperatorGrammarOptions$` | killed (TestOperatorGrammarOptions) | killed (TestOperatorGrammarOptions) | — | K |
| `OPT-S-exec-profiles-dir` | M2 · `operator.go` | `{"exec-profiles-dir", false, false}, {"task" -> {"exec-profiles-dir", false, true}, {"task"` | `^TestOperatorGrammarOptions$` | killed (TestOperatorGrammarOptions) | killed (TestOperatorGrammarOptions) | — | K |
| `OPT-S-presets-dir` | M2 · `operator.go` | `{"presets-dir", false, false}, {"resume" -> {"presets-dir", false, true}, {"resume"` | `^TestOperatorGrammarOptions$` | killed (TestOperatorGrammarOptions) | killed (TestOperatorGrammarOptions) | — | K |
| `OPT-F-presets-dir` | M2 · `operator.go` | `{"presets-dir", false, false}, {"title" -> {"presets-dir", false, true}, {"title"` | `^TestOperatorGrammarOptions$` | killed (TestOperatorGrammarOptions) | killed (TestOperatorGrammarOptions) | — | K |
| `OPT-F-title` | M2 · `operator.go` | `{"title", false, false} -> {"title", false, true}` | `^TestOperatorGrammarOptions$` | killed (TestOperatorGrammarOptions) | killed (TestOperatorGrammarOptions) | — | K |
| `OPT-F-include-unattributed` | M2 · `operator.go` | `{"include-unattributed", true, false} -> {"include-unattributed", true, true}` | `^TestOperatorGrammarOptions$` | killed (TestOperatorGrammarOptions) | killed (TestOperatorGrammarOptions) | — | K |
| `OPT-F-exec-profiles-dir` | M2 · `operator.go` | `{"exec-profiles-dir", false, false}, {"dry-run" -> {"exec-profiles-dir", false, true}, {"dry-run"` | `^TestOperatorGrammarOptions$` | killed (TestOperatorGrammarOptions) | killed (TestOperatorGrammarOptions) | — | K |
| `OPT-F-dry-run` | M2 · `operator.go` | `{"dry-run", true, false} -> {"dry-run", true, true}` | `^TestOperatorGrammarOptions$` | killed (TestOperatorGrammarOptions) | killed (TestOperatorGrammarOptions) | — | K |
| `OPT-F-allow-failing` | M2 · `operator.go` | `{"allow-failing", false, false} -> {"allow-failing", false, true}` | `^TestOperatorGrammarOptions$` | killed (TestOperatorGrammarOptions) | killed (TestOperatorGrammarOptions) | — | K |
| `MIR-S-preset-OUT` | M2 · `operator.go` | `{"ttl", false, false}, {"preset", false, true} -> {"ttl", false, false}, {"preset", false, false}` | `^TestOperatorGrammarCanonical$` | killed (TestOperatorGrammarCanonical) | killed (TestOperatorGrammarArity,TestOperatorGrammarCanonical,TestOperatorGrammarCase,TestOperatorGrammarOptions,TestOperatorGrammarPolicy,TestOperatorIdentifierAnchors,TestOperatorIdentifierFlags,TestOperatorIdentifierLength) | — | K |
| `MIR-S-launch-OUT` | M2 · `operator.go` | `{"launch", true, true} -> {"launch", true, false}` | `^TestOperatorGrammarCanonical$` | killed (TestOperatorGrammarCanonical) | killed (TestOperatorGrammarArity,TestOperatorGrammarCanonical,TestOperatorGrammarCase,TestOperatorGrammarOptions,TestOperatorGrammarPolicy,TestOperatorIdentifierAnchors,TestOperatorIdentifierFlags,TestOperatorIdentifierLength) | — | K |
| `MIR-S-headless-OUT` | M2 · `operator.go` | `{"headless", true, true} -> {"headless", true, false}` | `^TestOperatorGrammarCanonical$` | killed (TestOperatorGrammarCanonical) | killed (TestOperatorGrammarArity,TestOperatorGrammarCanonical,TestOperatorGrammarCase,TestOperatorGrammarOptions,TestOperatorGrammarPolicy,TestOperatorIdentifierAnchors,TestOperatorIdentifierFlags,TestOperatorIdentifierLength) | — | K |
| `MIR-F-ready-OUT` | M2 · `operator.go` | `{"ready", true, true} -> {"ready", true, false}` | `^TestOperatorGrammarCanonical$` | killed (TestOperatorGrammarCanonical) | killed (TestOperatorGrammarArity,TestOperatorGrammarCanonical,TestOperatorGrammarOptions,TestOperatorGrammarPolicy) | — | K |
| `MIR-F-cos-OUT` | M2 · `operator.go` | `{"confirm-on-surface", true, true} -> {"confirm-on-surface", true, false}` | `^TestOperatorGrammarCanonical$` | killed (TestOperatorGrammarCanonical) | killed (TestOperatorGrammarCanonical,TestOperatorGrammarOptions,TestOperatorGrammarPolicy) | — | K |
| `MIR-KIND-launch` | M2 · `operator.go` | `{"launch", true, true} -> {"launch", false, true}` | `^TestOperatorGrammarCanonical$` | killed (TestOperatorGrammarCanonical) | killed (TestOperatorGrammarArity,TestOperatorGrammarCanonical,TestOperatorGrammarCase,TestOperatorGrammarDuplicates,TestOperatorGrammarOptions,TestOperatorGrammarPolicy,TestOperatorIdentifierAnchors,TestOperatorIdentifierLength) | — | K |
| `MIR-KIND-request` | M2 · `operator.go` | `{"request", true, false} -> {"request", false, false}` | `^TestOperatorGrammarOptions$` | killed (TestOperatorGrammarOptions) | killed (TestOperatorGrammarOptions) | — | K |
| `MIR-KIND-ttl` | M2 · `operator.go` | `{"ttl", false, false} -> {"ttl", true, false}` | `^TestOperatorGrammarOptions$` | killed (TestOperatorGrammarOptions) | killed (TestOperatorGrammarOptions) | — | K |
| `MIR-DROP-hosts-dir` | M2 · `operator.go` | `{"hosts-dir", false, false},  -> ` | `^TestOperatorGrammarOptions$` | killed (TestOperatorGrammarOptions) | killed (TestOperatorGrammarOptions) | — | K |
| `MIR-RENAME-budget` | M2 · `operator.go` | `{"budget", false, false} -> {"budgets", false, false}` | `^TestOperatorGrammarOptions$` | killed (TestOperatorGrammarOptions) | killed (TestOperatorGrammarOptions) | — | K |
| `MIR-ADD-zzz` | M2 · `operator.go` | `{"sessions-dir", false, false}, {"addr", false, false}, -> {"sessions-dir", false, false}, {"addr", false, false}, {"zzz` | `^TestOperator` | survived (-) | survived (-) | killed (TestOperatorFlagParity) / killed (TestDoctorQuickIsFast,TestOperatorFlagParity,TestWorldContextAndStatusline) | K-join — a mirror entry the CLI lacks is invisible to host/verbs by construction; the join is the CLI parity test |
| `VAL-COUNT` | M2 · `operator.go` | ` v.n++\n ->  v.n += 0\n` | `^TestOperatorGrammarDuplicates$` | killed (TestOperatorGrammarDuplicates) | killed (TestOperatorGrammarDuplicates,TestOperatorGrammarOptions) | — | K |
| `VAL-BOOL-ALWAYS` | M2 · `operator.go` | `  v.on = b\n ->   v.on = true \|\| b\n` | `^TestOperatorGrammarOptions$` | killed (TestOperatorGrammarOptions) | killed (TestOperatorGrammarOptions) | — | K |
| `VAL-STRING-LOST` | M2 · `operator.go` | ` v.s = s\n ->  v.s = ""\n` | `^TestOperatorGrammarCanonical$` | killed (TestOperatorGrammarCanonical) | killed (TestOperatorGrammarCanonical,TestOperatorGrammarCase,TestOperatorGrammarOptions,TestOperatorGrammarPolicy,TestOperatorIdentifierAnchors,TestOperatorIdentifierFlags,TestOperatorIdentifierLength,TestOperatorOtherVerbParity) | — | K |
| `VAL-ISBOOL-TRUE` | M2 · `operator.go` | `{ return v.isBool } -> { return true }` | `^TestOperatorGrammarCanonical$` | killed (TestOperatorGrammarCanonical) | killed (TestOperatorGrammarArity,TestOperatorGrammarCanonical,TestOperatorGrammarCase,TestOperatorGrammarOptions,TestOperatorGrammarPolicy,TestOperatorIdentifierAnchors,TestOperatorIdentifierFlags,TestOperatorIdentifierLength,TestOperatorOtherVerbParity) | — | K |
| `VAL-ISBOOL-FALSE` | M2 · `operator.go` | `{ return v.isBool } -> { return false }` | `^TestOperatorGrammarCanonical$` | killed (TestOperatorGrammarCanonical) | killed (TestOperatorGrammarArity,TestOperatorGrammarCanonical,TestOperatorGrammarCase,TestOperatorGrammarDuplicates,TestOperatorGrammarExactWords,TestOperatorGrammarOptions,TestOperatorGrammarPolicy,TestOperatorIdentifierAnchors,TestOperatorIdentifierFlags,TestOperatorIdentifierLength) | — | K |
| `VAL-BOOL-ERRTEXT` | M2 · `operator.go` | `return errors.New("parse error") -> return errors.New("not a boolean")` | `^TestOperatorGrammarOptions$` | killed (TestOperatorGrammarOptions) | killed (TestOperatorGrammarOptions) | killed (TestOperatorParsesLikeTheCLI) / killed (TestOperatorParsesLikeTheCLI) | K |
| `BOOL-DROP-T` | M2 · `operator.go` | `"t": true, "T": true, -> "t": true,` | `^TestOperatorGrammarOptions$` | killed (TestOperatorGrammarOptions) | killed (TestOperatorGrammarCanonical,TestOperatorGrammarOptions) | — | K |
| `BOOL-FLIP-F` | M2 · `operator.go` | `"f": false, "F": false, -> "f": false, "F": true,` | `^TestOperatorGrammarOptions$` | killed (TestOperatorGrammarOptions) | killed (TestOperatorGrammarOptions) | — | K |
| `BOOL-ADD-YES` | M2 · `operator.go` | `"True": true,\n -> "True": true, "yes": true,\n` | `^TestOperatorGrammarOptions$` | killed (TestOperatorGrammarOptions) | killed (TestOperatorGrammarOptions) | — | K |
| `BOOL-DROP-EMPTYCHK` | M2 · `operator.go` | `b, ok := boolWords[s]\n  if !ok { -> b, ok := boolWords[s]\n  if !ok && s != "" {` | `^TestOperatorGrammarOptions$` | killed (TestOperatorGrammarOptions) | killed (TestOperatorGrammarOptions) | — | K |
| `HELP-ARM` | M2 · `operator.go` | `if errors.Is(err, flag.ErrHelp) { -> if false && errors.Is(err, flag.ErrHelp) {` | `^TestOperatorGrammarOptions$` | killed (TestOperatorGrammarOptions) | killed (TestOperatorGrammarOptions) | — | K |
| `PARSE-ERR` | M2 · `operator.go` | `if err != nil {\n  return refuse(OperatorUsage, "", err.Error()) -> if false && err != nil {\n  return refuse(OperatorUs` | `^TestOperatorGrammarCase$` | killed (TestOperatorGrammarCase) | killed (TestOperatorGrammarArity,TestOperatorGrammarCase,TestOperatorGrammarExactWords,TestOperatorGrammarOptions) | — | K |
| `DUP-GT2` | M2 · `operator.go` | `if got[f.Name].n > 1 { -> if got[f.Name].n > 2 {` | `^TestOperatorGrammarDuplicates$` | killed (TestOperatorGrammarDuplicates) | killed (TestOperatorGrammarDuplicates) | — | K |
| `NOTADMIT-OFF` | M2 · `operator.go` | `g.n == 1 && !f.Admit { -> g.n == 1 && !f.Admit && false {` | `^TestOperatorGrammarOptions$` | killed (TestOperatorGrammarOptions) | killed (TestOperatorGrammarOptions) | — | K |
| `EQ-NOTADMIT-GT0` | M2 · `operator.go` | `g.n == 1 && !f.Admit { -> g.n > 0 && !f.Admit {` | `^TestOperatorGrammarOptions$` | survived (-) | survived (-) | — | E — a count above one was refused by the duplicate loop just before, so `n == 1` and `n > 0` agree here |
| `FALSE-OFF` | M2 · `operator.go` | `f.Bool && !g.on { -> f.Bool && !g.on && false {` | `^TestOperatorGrammarOptions$` | killed (TestOperatorGrammarOptions) | killed (TestOperatorGrammarOptions) | — | K |
| `FALSE-NOBOOL` | M2 · `operator.go` | `f.Bool && !g.on { -> !g.on {` | `^TestOperatorGrammarOptions$` | killed (TestOperatorGrammarOptions) | killed (TestOperatorGrammarArity,TestOperatorGrammarCanonical,TestOperatorGrammarCase,TestOperatorGrammarOptions,TestOperatorGrammarPolicy,TestOperatorIdentifierAnchors,TestOperatorIdentifierFlags,TestOperatorIdentifierLength,TestOperatorOtherVerbParity) | — | K |
| `ARITY-GT1` | M2 · `operator.go` | `if len(pos) != 1 { -> if len(pos) > 1 {` | `^TestOperatorGrammarArity$` | killed (TestOperatorGrammarArity) | killed (TestOperatorGrammarArity) | — | K |
| `ARITY-LT1` | M2 · `operator.go` | `if len(pos) != 1 { -> if len(pos) < 1 {` | `^TestOperatorGrammarArity$` | killed (TestOperatorGrammarArity) | killed (TestOperatorGrammarArity) | — | K |
| `EP-CHECK-OFF` | M2 · `operator.go` | `if !EpisodeIDPattern.MatchString(pos[0]) { -> if false && !EpisodeIDPattern.MatchString(pos[0]) {` | `^TestOperatorGrammarCase$` | killed (TestOperatorGrammarCase) | killed (TestOperatorGrammarArity,TestOperatorGrammarCase,TestOperatorIdentifierAnchors,TestOperatorIdentifierFlags,TestOperatorIdentifierLength) | — | K |
| `PRESET-EMPTY-OFF` | M2 · `operator.go` | `if preset == "" { -> if false {` | `^TestOperatorGrammarArity$` | killed (TestOperatorGrammarArity) | killed (TestOperatorGrammarArity) | — | K |
| `PRESET-SETOOLS-OFF` | M2 · `operator.go` | `preset == "se-tools" \|\| len -> len` | `^TestOperatorGrammarOptions$` | killed (TestOperatorGrammarOptions) | killed (TestOperatorGrammarOptions) | — | K |
| `PRESET-GRAMMAR-OFF` | M2 · `operator.go` | ` \|\| !EpisodeIDPattern.MatchString(preset) { ->  {` | `^TestOperatorGrammarCase$` | killed (TestOperatorGrammarCase) | killed (TestOperatorGrammarCase,TestOperatorIdentifierAnchors,TestOperatorIdentifierFlags) | — | K |
| `PRESET-LEN-GE` | M2 · `operator.go` | `len(preset) > MaxPresetName -> len(preset) >= MaxPresetName` | `^TestOperatorIdentifierLength$` | killed (TestOperatorIdentifierLength) | killed (TestOperatorIdentifierLength) | — | K |
| `PRESET-LEN-PLUS1` | M2 · `operator.go` | `len(preset) > MaxPresetName -> len(preset) > MaxPresetName+1` | `^TestOperatorIdentifierLength$` | killed (TestOperatorIdentifierLength) | killed (TestOperatorIdentifierLength) | — | K |
| `PRESET-CAP-52` | M2 · `operator.go` | `const MaxPresetName = 51 -> const MaxPresetName = 52` | `^TestOperatorIdentifierLength$` | killed (TestOperatorIdentifierLength) | killed (TestOperatorIdentifierLength) | — | K |
| `PRESET-CAP-50` | M2 · `operator.go` | `const MaxPresetName = 51 -> const MaxPresetName = 50` | `^TestOperatorIdentifierLength$` | killed (TestOperatorIdentifierLength) | killed (TestOperatorIdentifierLength) | — | K |
| `EP-RE-CARET` | M2 · `operator.go` | `'^[a-z0-9][a-z0-9-]{0,63}$' -> '[a-z0-9][a-z0-9-]{0,63}$'` | `^TestOperatorIdentifierAnchors$` | killed (TestOperatorIdentifierAnchors) | killed (TestOperatorGrammarCase,TestOperatorIdentifierAnchors,TestOperatorIdentifierFlags,TestOperatorIdentifierLength) | — | K |
| `EP-RE-DOLLAR` | M2 · `operator.go` | `'^[a-z0-9][a-z0-9-]{0,63}$' -> '^[a-z0-9][a-z0-9-]{0,63}'` | `^TestOperatorIdentifierAnchors$` | killed (TestOperatorIdentifierAnchors) | killed (TestOperatorGrammarCase,TestOperatorIdentifierAnchors,TestOperatorIdentifierFlags,TestOperatorIdentifierLength) | — | K |
| `EP-RE-M` | M2 · `operator.go` | `'^[a-z0-9][a-z0-9-]{0,63}$' -> '(?m)^[a-z0-9][a-z0-9-]{0,63}$'` | `^TestOperatorIdentifierFlags$` | killed (TestOperatorIdentifierFlags) | killed (TestOperatorIdentifierFlags) | — | K |
| `EP-RE-I` | M2 · `operator.go` | `'^[a-z0-9][a-z0-9-]{0,63}$' -> '(?i)^[a-z0-9][a-z0-9-]{0,63}$'` | `^TestOperatorIdentifierFlags$` | killed (TestOperatorIdentifierFlags) | killed (TestOperatorGrammarCase,TestOperatorIdentifierFlags) | — | K |
| `EP-RE-S` | M2 · `operator.go` | `'^[a-z0-9][a-z0-9-]{0,63}$' -> '(?s)^[a-z0-9][a-z0-9-]{0,63}$'` | `^TestOperatorIdentifierFlags$` | survived (-) | survived (-) | — | E — as `P1-RE-S` (no `.` in the episode grammar) |
| `EP-RE-UP64` | M2 · `operator.go` | `'^[a-z0-9][a-z0-9-]{0,63}$' -> '^[a-z0-9][a-z0-9-]{0,64}$'` | `^TestOperatorIdentifierLength$` | killed (TestOperatorIdentifierLength) | killed (TestOperatorIdentifierLength) | — | K |
| `EP-RE-UP62` | M2 · `operator.go` | `'^[a-z0-9][a-z0-9-]{0,63}$' -> '^[a-z0-9][a-z0-9-]{0,62}$'` | `^TestOperatorIdentifierLength$` | killed (TestOperatorIdentifierLength) | killed (TestOperatorIdentifierLength) | — | K |
| `EP-RE-LOW0` | M2 · `operator.go` | `'^[a-z0-9][a-z0-9-]{0,63}$' -> '^[a-z0-9]?[a-z0-9-]{0,63}$'` | `^TestOperatorIdentifierLength$` | killed (TestOperatorIdentifierLength) | killed (TestOperatorGrammarArity,TestOperatorIdentifierLength) | — | K |
| `EP-RE-LEADDASH` | M2 · `operator.go` | `'^[a-z0-9][a-z0-9-]{0,63}$' -> '^[a-z0-9-][a-z0-9-]{0,63}$'` | `^TestOperatorIdentifierLength$` | killed (TestOperatorIdentifierLength) | killed (TestOperatorIdentifierLength) | — | K |
| `OA-RE-CARET` | M2 · `operator.go` | `'^oa-[0-9a-f]{32}$' -> 'oa-[0-9a-f]{32}$'` | `^TestOperatorIdentifierAnchors$` | killed (TestOperatorIdentifierAnchors) | killed (TestOperatorIdentifierAnchors,TestOperatorIdentifierFlags) | — | K |
| `OA-RE-DOLLAR` | M2 · `operator.go` | `'^oa-[0-9a-f]{32}$' -> '^oa-[0-9a-f]{32}'` | `^TestOperatorIdentifierAnchors$` | killed (TestOperatorIdentifierAnchors) | killed (TestOperatorIdentifierAnchors,TestOperatorIdentifierFlags,TestOperatorIdentifierLength) | — | K |
| `OA-RE-M` | M2 · `operator.go` | `'^oa-[0-9a-f]{32}$' -> '(?m)^oa-[0-9a-f]{32}$'` | `^TestOperatorIdentifierFlags$` | killed (TestOperatorIdentifierFlags) | killed (TestOperatorIdentifierFlags) | — | K |
| `OA-RE-I` | M2 · `operator.go` | `'^oa-[0-9a-f]{32}$' -> '(?i)^oa-[0-9a-f]{32}$'` | `^TestOperatorIdentifierFlags$` | killed (TestOperatorIdentifierFlags) | killed (TestOperatorIdentifierFlags) | — | K |
| `OA-RE-S` | M2 · `operator.go` | `'^oa-[0-9a-f]{32}$' -> '(?s)^oa-[0-9a-f]{32}$'` | `^TestOperatorIdentifierFlags$` | survived (-) | survived (-) | — | E — as `P1-RE-S` (no `.` in the action-id grammar) |
| `OA-RE-31` | M2 · `operator.go` | `'^oa-[0-9a-f]{32}$' -> '^oa-[0-9a-f]{31,32}$'` | `^TestOperatorIdentifierLength$` | killed (TestOperatorIdentifierLength) | killed (TestOperatorIdentifierLength) | — | K |
| `OA-RE-33` | M2 · `operator.go` | `'^oa-[0-9a-f]{32}$' -> '^oa-[0-9a-f]{32,33}$'` | `^TestOperatorIdentifierLength$` | killed (TestOperatorIdentifierLength) | killed (TestOperatorIdentifierLength) | — | K |
| `OA-RE-NODASH` | M2 · `operator.go` | `'^oa-[0-9a-f]{32}$' -> '^oa-?[0-9a-f]{32}$'` | `^TestOperatorIdentifierAnchors$` | killed (TestOperatorIdentifierAnchors) | killed (TestOperatorIdentifierAnchors) | — | K |
| `OA-RE-ANYSEP` | M2 · `operator.go` | `'^oa-[0-9a-f]{32}$' -> '^oa.[0-9a-f]{32}$'` | `^TestOperatorIdentifierAnchors$` | killed (TestOperatorIdentifierAnchors) | killed (TestOperatorIdentifierAnchors) | — | K |
| `OA-RE-LEN16` | M2 · `operator.go` | `'^oa-[0-9a-f]{32}$' -> '^oa-[0-9a-f]{16}$'` | `^TestOperatorIdentifierLength$` | killed (TestOperatorIdentifierLength) | killed (TestOperatorIdentifierAnchors,TestOperatorIdentifierLength) | — | K |
| `WP-OFF` | M2 · `operator.go` | `argv[0] == "world-publish" { -> argv[0] == "world-publishx" {` | `^TestOperatorNeverDispatch$` | killed (TestOperatorNeverDispatch) | killed (TestOperatorNeverDispatch) | — | K |
| `WP-ARITY` | M2 · `operator.go` | `len(argv) > 0 && argv[0] == "world-publish" -> len(argv) > 1 && argv[0] == "world-publish"` | `^TestOperatorNeverDispatch$` | killed (TestOperatorNeverDispatch) | killed (TestOperatorNeverDispatch) | — | K |
| `NEV-approver` | M2 · `operator.go` | ` case "episode start":\n ->  case "episode start", "approver":\n` | `^(TestOperatorNeverDispatch\|TestOperatorOtherVerbParity)$` | killed (TestOperatorNeverDispatch,TestOperatorOtherVerbParity) | killed (TestOperatorNeverDispatch,TestOperatorOtherVerbParity) | — | K |
| `NEV-session` | M2 · `operator.go` | ` case "episode start":\n ->  case "episode start", "session":\n` | `^(TestOperatorNeverDispatch\|TestOperatorOtherVerbParity)$` | killed (TestOperatorNeverDispatch,TestOperatorOtherVerbParity) | killed (TestOperatorNeverDispatch,TestOperatorOtherVerbParity) | — | K |
| `NEV-service` | M2 · `operator.go` | ` case "episode start":\n ->  case "episode start", "service":\n` | `^(TestOperatorNeverDispatch\|TestOperatorOtherVerbParity)$` | killed (TestOperatorNeverDispatch,TestOperatorOtherVerbParity) | killed (TestOperatorNeverDispatch,TestOperatorOtherVerbParity) | — | K |
| `NEV-serve` | M2 · `operator.go` | ` case "episode start":\n ->  case "episode start", "serve":\n` | `^(TestOperatorNeverDispatch\|TestOperatorOtherVerbParity)$` | killed (TestOperatorNeverDispatch,TestOperatorOtherVerbParity) | killed (TestOperatorNeverDispatch,TestOperatorOtherVerbParity) | — | K |
| `NEV-setup` | M2 · `operator.go` | ` case "episode start":\n ->  case "episode start", "setup":\n` | `^(TestOperatorNeverDispatch\|TestOperatorOtherVerbParity)$` | killed (TestOperatorNeverDispatch,TestOperatorOtherVerbParity) | killed (TestOperatorNeverDispatch,TestOperatorOtherVerbParity) | — | K |
| `NEV-episode-launch` | M2 · `operator.go` | ` case "episode start":\n ->  case "episode start", "episode launch":\n` | `^(TestOperatorNeverDispatch\|TestOperatorOtherVerbParity)$` | killed (TestOperatorNeverDispatch) | killed (TestOperatorNeverDispatch) | — | K |
| `NEV-episode-report` | M2 · `operator.go` | ` case "episode start":\n ->  case "episode start", "episode report":\n` | `^(TestOperatorNeverDispatch\|TestOperatorOtherVerbParity)$` | killed (TestOperatorNeverDispatch,TestOperatorOtherVerbParity) | killed (TestOperatorGrammarExactWords,TestOperatorNeverDispatch,TestOperatorOtherVerbParity) | — | K |
| `NEV-why` | M2 · `operator.go` | ` case "episode start":\n ->  case "episode start", "why":\n` | `^(TestOperatorNeverDispatch\|TestOperatorOtherVerbParity)$` | killed (TestOperatorNeverDispatch,TestOperatorOtherVerbParity) | killed (TestOperatorNeverDispatch,TestOperatorOtherVerbParity) | — | K |
| `NEV-approve` | M2 · `operator.go` | ` case "episode start":\n ->  case "episode start", "approve":\n` | `^(TestOperatorNeverDispatch\|TestOperatorOtherVerbParity)$` | killed (TestOperatorNeverDispatch,TestOperatorOtherVerbParity) | killed (TestOperatorNeverDispatch,TestOperatorOtherVerbParity) | — | K |
| `NEV-episode` | M2 · `operator.go` | ` case "episode start":\n ->  case "episode start", "episode":\n` | `^(TestOperatorNeverDispatch\|TestOperatorOtherVerbParity)$` | killed (TestOperatorOtherVerbParity) | killed (TestOperatorGrammarCase,TestOperatorGrammarExactWords,TestOperatorOtherVerbParity) | — | K |
| `NEV-CLASS` | M2 · `operator.go` | `if v.Class == Never { -> if v.Class == Refused {` | `^TestOperatorNeverDispatch$` | killed (TestOperatorNeverDispatch) | killed (TestOperatorGrammarCase,TestOperatorGrammarExactWords,TestOperatorNeverDispatch,TestOperatorOtherVerbParity) | — | K |
| `OTHER-REASON` | M2 · `operator.go` | `return OperatorAnswer{Class: OperatorOther, Reason: v.Reason, CLI: v.CLI} -> return OperatorAnswer{Class: OperatorOther,` | `^TestOperatorOtherVerbParity$` | killed (TestOperatorOtherVerbParity) | killed (TestOperatorGrammarCase,TestOperatorGrammarExactWords,TestOperatorNeverDispatch,TestOperatorOtherVerbParity) | — | K |
| `OTHER-CLI` | M2 · `operator.go` | `return OperatorAnswer{Class: OperatorOther, Reason: v.Reason, CLI: v.CLI} -> return OperatorAnswer{Class: OperatorOther,` | `^TestOperatorOtherVerbParity$` | killed (TestOperatorOtherVerbParity) | killed (TestOperatorOtherVerbParity) | — | K |
| `SWAP-START-FINISH` | M2 · `operator.go` | `mirror = EpisodeStartFlags\n case "episode finish":\n  mirror = EpisodeFinishFlags -> mirror = EpisodeFinishFlags\n case` | `^TestOperatorGrammarCanonical$` | killed (TestOperatorGrammarCanonical) | killed (TestOperatorGrammarPolicy,TestOperatorIdentifierAnchors,TestOperatorIdentifierFlags,TestOperatorIdentifierLength) | — | K |
| `AT-EMPTY` | M2 · `operator.go` | `at := "; at a terminal: " + v.CLI -> at := ""` | `^TestOperatorGrammarOptions$` | killed (TestOperatorGrammarOptions) | killed (TestOperatorGrammarOptions,TestOperatorGrammarPolicy) | — | K |
| `REFUSE-PREFIX` | M2 · `operator.go` | `Reason: Binary + " " + v.Words + ": " + why -> Reason: why` | `^TestOperatorGrammarOptions$` | killed (TestOperatorGrammarOptions) | killed (TestOperatorGrammarOptions) | killed (TestOperatorParsesLikeTheCLI) / killed (TestOperatorParsesLikeTheCLI) | K |
| `MODE-HL-NOLAUNCH` | M2 · `operator.go` | `case headless && !launch: -> case false:` | `^TestOperatorGrammarOptions$` | killed (TestOperatorGrammarOptions) | killed (TestOperatorGrammarOptions) | — | K |
| `MODE-NOLAUNCH` | M2 · `operator.go` | `case !launch: -> case false:` | `^TestOperatorGrammarOptions$` | killed (TestOperatorGrammarOptions) | killed (TestOperatorGrammarOptions) | — | K |
| `MODE-HL-POLICY` | M2 · `operator.go` | `case headless && !p.headless(): -> case false:` | `^TestOperatorGrammarPolicy$` | killed (TestOperatorGrammarPolicy) | killed (TestOperatorGrammarPolicy) | — | K |
| `MODE-IA-POLICY` | M2 · `operator.go` | `case !headless && !p.interactive(): -> case false:` | `^TestOperatorGrammarPolicy$` | killed (TestOperatorGrammarPolicy) | killed (TestOperatorGrammarPolicy) | — | K |
| `POL-HL-NOB` | M2 · `operator.go` | `return p.Start == "A" \|\| p.Start == "B" } -> return p.Start == "A" }` | `^TestOperatorGrammarPolicy$` | killed (TestOperatorGrammarPolicy) | killed (TestOperatorGrammarPolicy) | — | K |
| `POL-HL-NOTC` | M2 · `operator.go` | `return p.Start == "A" \|\| p.Start == "B" } -> return p.Start != "C" }` | `^TestOperatorGrammarPolicy$` | killed (TestOperatorGrammarPolicy) | killed (TestOperatorGrammarPolicy) | — | K |
| `POL-HL-FOLD` | M2 · `operator.go` | `return p.Start == "A" \|\| p.Start == "B" } -> return p.Start == "A" \|\| p.Start == "B" \|\| p.Start == "a" }` | `^TestOperatorGrammarPolicy$` | killed (TestOperatorGrammarPolicy) | killed (TestOperatorGrammarPolicy) | — | K |
| `POL-IA-NOHANDOFF` | M2 · `operator.go` | `return p.Start == "A" && p.HandoffProven } -> return p.Start == "A" }` | `^TestOperatorGrammarPolicy$` | killed (TestOperatorGrammarPolicy) | killed (TestOperatorGrammarPolicy) | — | K |
| `POL-IA-B` | M2 · `operator.go` | `return p.Start == "A" && p.HandoffProven } -> return p.headless() && p.HandoffProven }` | `^TestOperatorGrammarPolicy$` | killed (TestOperatorGrammarPolicy) | killed (TestOperatorGrammarPolicy) | — | K |
| `POL-FIN-NOSTART` | M2 · `operator.go` | `return p.headless() && (p.Finish == "A" \|\| p.Finish == "B") } -> return (p.Finish == "A" \|\| p.Finish == "B") }` | `^TestOperatorGrammarPolicy$` | killed (TestOperatorGrammarPolicy) | killed (TestOperatorGrammarPolicy) | — | K |
| `POL-FIN-NOB` | M2 · `operator.go` | `return p.headless() && (p.Finish == "A" \|\| p.Finish == "B") } -> return p.headless() && p.Finish == "A" }` | `^TestOperatorGrammarPolicy$` | killed (TestOperatorGrammarPolicy) | killed (TestOperatorGrammarPolicy) | — | K |
| `POL-FIN-NOTC` | M2 · `operator.go` | `return p.headless() && (p.Finish == "A" \|\| p.Finish == "B") } -> return p.headless() && p.Finish != "C" }` | `^TestOperatorGrammarPolicy$` | killed (TestOperatorGrammarPolicy) | killed (TestOperatorGrammarPolicy) | — | K |
| `POL-READY-B` | M2 · `operator.go` | `return p.Finish == "A" } -> return p.Finish == "A" \|\| p.Finish == "B" }` | `^TestOperatorGrammarPolicy$` | killed (TestOperatorGrammarPolicy) | killed (TestOperatorGrammarPolicy) | — | K |
| `FIN-POLICY-OFF` | M2 · `operator.go` | `if !p.finish() { -> if false {` | `^TestOperatorGrammarPolicy$` | killed (TestOperatorGrammarPolicy) | killed (TestOperatorGrammarPolicy) | — | K |
| `FIN-READY-OFF` | M2 · `operator.go` | `if in.Ready && !p.ready() { -> if false && in.Ready && !p.ready() {` | `^TestOperatorGrammarPolicy$` | killed (TestOperatorGrammarPolicy) | killed (TestOperatorGrammarPolicy) | — | K |
| `FIN-READY-ALL` | M2 · `operator.go` | `if in.Ready && !p.ready() { -> if !p.ready() {` | `^TestOperatorGrammarPolicy$` | killed (TestOperatorGrammarPolicy) | killed (TestOperatorGrammarPolicy) | — | K |
| `CAN-HEADLESS-FIELD` | M2 · `operator.go` | `in.Kind, in.Headless = "start", headless -> in.Kind, in.Headless = "start", true` | `^TestOperatorGrammarCanonical$` | killed (TestOperatorGrammarCanonical) | killed (TestOperatorGrammarCanonical) | — | K |
| `CAN-HEADLESS-ARGV` | M2 · `operator.go` | `in.Argv = append(in.Argv, "--headless") -> in.Argv = in.Argv[:len(in.Argv):len(in.Argv)]` | `^TestOperatorGrammarCanonical$` | killed (TestOperatorGrammarCanonical) | killed (TestOperatorGrammarCanonical) | — | K |
| `CAN-COS-ARGV` | M2 · `operator.go` | `in.Argv = append(in.Argv, "--confirm-on-surface") -> in.Argv = in.Argv[:len(in.Argv):len(in.Argv)]` | `^TestOperatorGrammarCanonical$` | killed (TestOperatorGrammarCanonical) | killed (TestOperatorGrammarCanonical) | — | K |
| `CAN-READY-ARGV` | M2 · `operator.go` | `in.Argv = append(in.Argv, "--ready") -> in.Argv = in.Argv[:len(in.Argv):len(in.Argv)]` | `^TestOperatorGrammarCanonical$` | killed (TestOperatorGrammarCanonical) | killed (TestOperatorGrammarCanonical) | — | K |
| `CAN-FIN-SWAP` | M2 · `operator.go` | `"finish", got["confirm-on-surface"].on, got["ready"].on -> "finish", got["ready"].on, got["confirm-on-surface"].on` | `^TestOperatorGrammarCanonical$` | killed (TestOperatorGrammarCanonical) | killed (TestOperatorGrammarCanonical,TestOperatorGrammarPolicy) | — | K |
| `CAN-START-CLI` | M2 · `operator.go` | `return OperatorAnswer{Class: OperatorReview, Intent: in, CLI: v.CLI}\n } -> return OperatorAnswer{Class: OperatorReview,` | `^TestOperatorGrammarCanonical$` | killed (TestOperatorGrammarCanonical) | killed (TestOperatorGrammarCanonical) | — | K |
| `CAN-START-EPISODE` | M2 · `operator.go` | `"start", in.Episode, "--preset" -> "start", preset, "--preset"` | `^TestOperatorGrammarCanonical$` | killed (TestOperatorGrammarCanonical) | killed (TestOperatorGrammarCanonical) | — | K |
| `SEAM-START-OFF` | M4 · `episode.go` | ` if episodeFlagSetSeen != nil {\n  episodeFlagSetSeen(fs)\n }\n positional, err ->  positional, err` | `^TestOperator` | killed (TestOperatorFlagParity,TestOperatorParsesLikeTheCLI) | killed (TestOperatorFlagParity,TestOperatorParsesLikeTheCLI) | — | K |
| `SEAM-FINISH-OFF` | M4 · `episode_finish.go` | ` if episodeFlagSetSeen != nil {\n  episodeFlagSetSeen(fs)\n }\n pos, err ->  pos, err` | `^TestOperator` | killed (TestOperatorFlagParity,TestOperatorParsesLikeTheCLI) | killed (TestOperatorFlagParity,TestOperatorParsesLikeTheCLI) | — | K |
| `CLI-ADD-FLAG` | M4 · `episode.go` | ` fs.String("hosts-dir", defaultHostsDir(), "hosts init output directory")\n ->  fs.String("hosts-dir", defaultHostsDir()` | `^TestOperator` | killed (TestOperatorFlagParity) | killed (TestOperatorFlagParity) | — | K |
| `CLI-REMOVE-FLAG` | M4 · `episode.go` | ` fs.String("hosts-dir", defaultHostsDir(), "hosts init output directory")\n -> ` | `^TestOperator` | killed (TestOperatorFlagParity) | killed (TestStartLaunchHeadlessMintsAndLaunchesOnce) | — | K; whole again with the panicking `TestStartLaunchHeadlessMintsAndLaunchesOnce` skipped: killed (TestOperatorFlagParity) |
| `CLI-RENAME-FLAG` | M4 · `episode_finish.go` | `fs.String("title", "", "commit title") -> fs.String("subject", "", "commit title")` | `^TestOperator` | killed (TestOperatorFlagParity) | killed (TestEpisodeFinishAttributionCommitDryRunAC5AC6,TestEpisodeFinishDispatchAC9AC11AC12,TestOperatorFlagParity) | — | K |
| `CLI-ADD-FINISH-FLAG` | M4 · `episode_finish.go` | ` fs.Var(&allowed, "allow-failing", "allow this failing check")\n ->  fs.Var(&allowed, "allow-failing", "allow this faili` | `^TestOperator` | killed (TestOperatorFlagParity) | killed (TestOperatorFlagParity) | — | K |
| `CLI-PRESETCAP` | M4 · `preset.go` | `const maxPresetNameLen = 51 -> const maxPresetNameLen = 52` | `^TestOperator` | killed (TestOperatorIdentifierGrammarsMatchTheCLI) | killed (TestOperatorIdentifierGrammarsMatchTheCLI,TestPresetRefusalsBeforeDaemon) | — | K |
| `CLI-EPGRAMMAR` | M4 · `doctor.go` | `'^[a-z0-9][a-z0-9-]{0,63}$' -> '^[a-z0-9][a-z0-9-]{0,64}$'` | `^TestOperator` | killed (TestOperatorIdentifierGrammarsMatchTheCLI) | killed (TestDoctorEpisodeGrammarMatchesTheDaemon,TestOperatorIdentifierGrammarsMatchTheCLI) | — | K |
| `SEAM-AFTER-PARSE` | M4 · `episode.go` | ` if episodeFlagSetSeen != nil {\n  episodeFlagSetSeen(fs)\n }\n positional, err := verbs.ParseInterleaved(fs, args)\n ->` | `^TestOperator` | survived (-) | survived (-) | — | E — a parse never changes a flag set's definitions, so the hook sees the same set before or after it (start's hook precedes its error check either way) |

**Totals: 153 mutants (NM excluded), 148 killed (1 by the CLI join only, by construction), 5 named equivalents, 0 survivors, 0 NM, 0 pending.** Every `K` is killed by its named test and by its whole package.

---

## 5. Files (final `git diff --name-status d92983d5 HEAD`)

```
M	cmd/ailang-worldd/episode.go
M	cmd/ailang-worldd/episode_finish.go
A	cmd/ailang-worldd/operator_parity_test.go
A	host/surfacepage/app_pins_test.go
A	host/verbs/operator.go
A	host/verbs/operator_policy_test.go
A	host/verbs/operator_test.go
```

Two pre-existing production files (the seam, F2); no pre-existing test file; nothing deleted; no `app.js`, `daemon.go`,
`surface_page.go`, `host/store`, docs, `.ail`, `world`, `packages`, `scripts`, `tools` or `plugins` file.

---

## 6. Sprint gates (controller, after the last sub-milestone's commit; each measured at the base, §1.1)

`bash gates.sh <worktree> $HOME/.ailang/state/world-iter407/gates final <HEAD sha>` (one bash shell; `GATES_SKIP_VERIFY_GO=1`
skips `verify_go` when the controller runs it separately). Each leg is its own subshell with its own rc, `tee`'d to
`<logdir>/<leg>.txt`; no leg is chained to another; `ALL_DONE.<sha>` holds `FAILS: …`. `TOUCHED_PACKAGES` is derived from
`git diff --name-only` and must equal `cmd/ailang-worldd host/surfacepage host/verbs`. Never start it under `nohup`.

````bash gates.sh
#!/bin/bash
# Row 243 V4 S1 gates (plan §6). usage: gates.sh <tree> <logdir> <base|final> [<sha>]
# The controller: bash gates.sh <worktree> $HOME/.ailang/state/world-iter407/gates final <HEAD sha>
# macOS /bin/bash 3.2: no declare -A. Every leg records its own rc; every leg runs (no set -e, no leg
# chained to another with &&); every leg is tee'd to <logdir>/<leg>.txt. GATES_SKIP_VERIFY_GO=1 skips
# verify_go (the executor sets it; the controller runs it). Never start this under nohup. Writes
# <logdir>/ALL_DONE.<sha> with FAILS: when every leg has finished.
T=$1; L=$2; MODE=$3
mkdir -p "$L"; cd "$T" || exit 2
SHA=${4:-$(git rev-parse HEAD 2>/dev/null || echo unknown)}
B=d92983d59f887905251df5b99e825c299136760f
export AILANG_BIN=$HOME/.pinned-ailang/ailang
export WORLD_EXEC_NODE=/opt/homebrew/bin/node
export WORLD_TOOL_AILANG_BIN=$HOME/.pinned-ailang-tools/v0.54.15/ailang
unset WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER WORLD_REQUIRE_SSH_KEYGEN WORLD_PASSKEY_RECORD_DIR WORLD_UPDATE_VERB_GOLDEN WORLD_UPDATE_WHY_GOLDEN WORLD_UPDATE_FINISH_PROMPT_GOLDEN WORLD_UPDATE_CONFIRM_GOLDEN WORLD_UPDATE_CONFIRMATIONS_GOLDEN
: > "$L/legs.txt"
FAILS=""
s0=$(date +%s)
leg() { # leg <name> <command string>: tee'd to $L/<name>.txt, rc recorded in legs.txt
  local name=$1 cmd=$2 s rc
  s=$(date +%s)
  bash -c "$cmd" 2>&1 | tee "$L/$name.txt"
  rc=${PIPESTATUS[0]}
  echo "LEG $name rc=$rc $(( $(date +%s) - s ))s" | tee -a "$L/legs.txt"
  [ "$rc" -eq 0 ] || FAILS="$FAILS $name"
}
need() { # need <leg> <test id>...: each must show --- PASS in <leg>.txt
  local lg=$1 t; shift
  for t in "$@"; do grep -qE -- "^ *--- PASS: $t " "$L/$lg.txt" || { echo "NOT PASSED ($lg): $t" | tee -a "$L/legs.txt"; FAILS="$FAILS ${lg}_missing"; }; done
}
alt() { local x="$*"; x="${x// /|}"; echo "${x:+|$x}"; }
bar() { echo "$*" | tr ' ' '|'; }
# The touched Go packages, derived from the diff (printed so the controller can diff it). The plan's list is
# cmd/ailang-worldd host/surfacepage host/verbs; a difference FAILS.
TOUCHED=$(git diff --name-only "$B" HEAD -- cmd host 2>/dev/null | grep '\.go$' | xargs -n1 dirname 2>/dev/null | sort -u | tr '\n' ' ')
echo "TOUCHED_PACKAGES: $TOUCHED" | tee -a "$L/legs.txt"
if [ "$MODE" = final ] && [ "$TOUCHED" != "cmd/ailang-worldd host/surfacepage host/verbs " ]; then FAILS="$FAILS touched_packages_differ"; fi
leg status 'git status --porcelain | wc -l'
leg verify_ail './scripts/verify_ail.sh'
leg gofmt 'out=$(gofmt -l cmd host); echo "$out"; test -z "$out"'
leg build 'go build ./...'
leg vet 'go vet ./...'
leg fence "go test -run '^\$' ./... > '$L/fence.raw' 2>&1; rc=\$?; cat '$L/fence.raw'; if grep -E '^(FAIL|--- FAIL)' '$L/fence.raw' >/dev/null; then exit 1; fi; exit \$rc"
# the three touched packages whole, and host/daemon (it imports host/verbs), each its own rc
leg pkgs "r=0; for p in ./host/verbs/ ./host/surfacepage/ ./cmd/ailang-worldd/ ./host/daemon/; do go test -count=1 -timeout 30m \$p; r=\$(( r | \$? )); done; exit \$r"
# -race on every touched package: host/verbs, host/surfacepage whole; cmd/ailang-worldd on the operator, episode
# start/finish, help and verb-table tests; host/daemon (a dependant) on the verb-route tests
CRACE='^(TestOperator|TestEpisodeStart|TestEpisodeFinish|TestEveryVerbHasHelp|TestVerbTableMatchesTheCLISwitch|TestDoctorEpisodeGrammarMatchesTheDaemon|TestCLICredentialWriterSites)'
DRACE='^(TestSurfaceVerb|TestDaemonCredentialWriterSites)'
leg race "r=0; go test -race -count=1 ./host/verbs/ ./host/surfacepage/; r=\$(( r | \$? )); go test -race -count=1 -v ./cmd/ailang-worldd/ -run '$CRACE' > '$L/race_cmd.raw' 2>&1; r=\$(( r | \$? )); grep -E '^(--- |ok|FAIL)' '$L/race_cmd.raw'; go test -race -count=1 -v ./host/daemon/ -run '$DRACE' > '$L/race_daemon.raw' 2>&1; r=\$(( r | \$? )); grep -E '^(--- |ok|FAIL)' '$L/race_daemon.raw'; exit \$r"
# node: every page harness, non-skippable (WORLD_EXEC_NODE is set above; a harness that cannot run FAILS)
NODE="TestSurfaceAppPure TestSurfaceScriptPure TestSurfaceAppVerbs TestNeedsYouListsDecisions TestNeedsYouListsConfirmations TestAddressHintsArePinned TestConfirmPageFailsClosed TestConfirmPageRendersText TestConfirmPageMessage TestConfirmPageStates TestConfirmPageSubmitsAPasskey TestConfirmPageShowsRefusals TestConfirmPageSubmitsAnSSHSIG"
[ "$MODE" = final ] && NODE="$NODE TestNeedsYouConfirmationEmbeddedNewline TestNeedsYouConfirmationsNonArray TestVerbDecidePrefixBoundary"
leg node "set -o pipefail; go test -v -count=1 ./host/surfacepage/ -run '^($(bar $NODE))\$'"
need node $NODE
if grep -E '^\s*--- (SKIP|FAIL)' "$L/node.txt"; then FAILS="$FAILS node_skip_or_fail"; fi
NEW_CMD="TestOperatorFlagParity TestOperatorParsesLikeTheCLI TestOperatorIdentifierGrammarsMatchTheCLI"
NEW_V="TestOperatorGrammarCase TestOperatorGrammarArity TestOperatorGrammarExactWords TestOperatorGrammarCanonical TestOperatorGrammarDuplicates TestOperatorGrammarOptions TestOperatorGrammarPolicy TestOperatorNeverDispatch TestOperatorIdentifierAnchors TestOperatorIdentifierFlags TestOperatorIdentifierLength TestOperatorOtherVerbParity"
if [ "$MODE" = base ]; then NEW_CMD=""; NEW_V=""; fi
INV_DAEMON="TestDaemonDependencyAllowlist TestFrozenV1RouteTableMatchesMux TestSurfaceAppRoutes TestApprovalRoutesRegistered TestDecisionRoutesRegistered TestConfirmRoutesRegistered TestDaemonCredentialWriterSites TestSurfaceVerbRouteRegistered TestSurfaceVerbIsReadOnly TestSurfaceVerbOutcomes TestDecidePageAssets TestConfirmPageAssetsServed TestRequestsJSONConfirmationsGolden"
INV_BOUNDARY="TestWorldBoundaryDependencyAllowlist TestBareNetHTTPExemptionIsPerGroup TestEveryCommandDirectoryIsAProtectedGroup TestBoundaryASTWriteGuard"
INV_STORE="TestProductionContextRoots TestProductionGoSurface TestStoreSurfaceTakesContext TestCredentialsSchemaMatchesCanonicalDDL TestCredentialsSchema9VerbRequests TestCredentialsNewerThanNineIsRefused"
INV_SP="TestSurfaceAppMatchesWorkingCopy TestSurfaceAppShellAssembly TestConfirmCardGolden"
INV_VERBS="TestLookup TestApproveTarget TestEpisodeReportFlags TestParseInterleaved TestSplit TestNeverAndOtherReasonsEndWithTheCLILine"
INV_CMD="TestEveryVerbHasHelp TestVerbTableMatchesTheCLISwitch TestEpisodeFinishHelp TestDoctorEpisodeGrammarMatchesTheDaemon TestCLICredentialWriterSites TestQuickstartSection19MatchesTheBar TestQuickstartSection20MatchesTheFinishFlag"
leg inv "set -o pipefail; r=0
go test -v -count=1 ./host/daemon/ -run '^($(bar $INV_DAEMON))\$'; r=\$(( r | \$? ))
go test -v -count=1 ./host/boundary/ -run '^($(bar $INV_BOUNDARY))\$'; r=\$(( r | \$? ))
go test -v -count=1 ./host/store/ -run '^($(bar $INV_STORE))\$'; r=\$(( r | \$? ))
go test -v -count=1 ./host/surfacepage/ -run '^($(bar $INV_SP))\$'; r=\$(( r | \$? ))
go test -v -count=1 ./host/verbs/ -run '^($(bar $INV_VERBS)$(alt $NEW_V))\$'; r=\$(( r | \$? ))
go test -v -count=1 ./cmd/ailang-worldd/ -run '^($(bar $INV_CMD)$(alt $NEW_CMD))\$'; r=\$(( r | \$? ))
exit \$r"
need inv $INV_DAEMON $INV_BOUNDARY $INV_STORE $INV_SP $INV_VERBS $NEW_V $INV_CMD $NEW_CMD
if grep -E '^\s*--- (SKIP|FAIL)' "$L/inv.txt"; then FAILS="$FAILS inv_skip_or_fail"; fi
echo "named PASS: $(grep -cE '^--- PASS' "$L/inv.txt")" | tee -a "$L/legs.txt"
# No new daemon route: every mux registration in daemon.go, counted (base = final), and the /surface POSTs (3)
ROUTES=$(grep -cE 'mux\.(HandleFunc|Handle)\(' host/daemon/daemon.go); POSTS=$(grep -cE 'mux\.HandleFunc\("POST /surface/' host/daemon/daemon.go)
BROUTES=$(git show "$B:host/daemon/daemon.go" | grep -cE 'mux\.(HandleFunc|Handle)\(')
echo "DAEMON_ROUTES: $ROUTES (base $BROUTES) SURFACE_POSTS: $POSTS" | tee -a "$L/legs.txt"
if [ "$ROUTES" != "$BROUTES" ] || [ "$POSTS" != 3 ]; then FAILS="$FAILS routes_changed"; fi
# Production bytes: S1's only production files are host/verbs/operator.go (NEW) and the two-site test seam in
# cmd/ailang-worldd/episode.go and episode_finish.go. Every other changed file is a _test.go; nothing under
# host/daemon, host/store, host/surfacepage/app, world, packages, scripts, tools, plugins; no .ail.
NOPROD=$(git diff --name-only "$B" HEAD | grep -vE '(_test\.go$|^host/verbs/operator\.go$|^cmd/ailang-worldd/episode(_finish)?\.go$|^design_docs/planned/(w-world-address-bar-verbs-v4s1-sprint-plan\.md|sprint_w-world-address-bar-verbs-v4s1\.json)$)' | tr '\n' ' ')
NOPROD="$NOPROD$(git diff --name-only "$B" HEAD -- host/daemon host/store host/surfacepage/app world packages scripts tools plugins '*.ail' | grep -v '_test\.go$' | tr '\n' ' ')"
echo "UNTOUCHED_VIOLATIONS: ${NOPROD:-none}" | tee -a "$L/legs.txt"
if [ "$MODE" = final ] && [ -n "$NOPROD" ]; then FAILS="$FAILS untouched_violated"; fi
# The seam is the only production change in cmd: exactly 10 added lines, none removed.
SEAM=$(git diff --numstat "$B" HEAD -- cmd/ailang-worldd/episode.go cmd/ailang-worldd/episode_finish.go | awk '{a+=$1; d+=$2} END {print a+0, d+0}')
echo "CMD_SEAM_LINES: $SEAM" | tee -a "$L/legs.txt"
if [ "$MODE" = final ] && [ "$SEAM" != "10 0" ]; then FAILS="$FAILS seam_lines"; fi
# The parser has no production caller (a caller search with a positive control in the same call).
CALLERS=$(grep -rnE '\b(AdmitOperator|ValidOperatorActionID|OperatorPolicy|EpisodeStartFlags|EpisodeFinishFlags|OperatorActionIDPattern|EpisodeIDPattern|MaxPresetName)\b' --include='*.go' cmd host | grep -v '_test\.go:' | grep -v '^host/verbs/operator\.go:' | tr '\n' ' ')
CONTROL=$(grep -rnE '\bverbs\.AdmitWhy\(' --include='*.go' cmd host | grep -v '_test\.go:' | wc -l | tr -d ' ')
echo "PARSER_PRODUCTION_CALLERS: ${CALLERS:-none} (control: verbs.AdmitWhy production callers $CONTROL)" | tee -a "$L/legs.txt"
if [ "$MODE" = final ] && { [ -n "$CALLERS" ] || [ "$CONTROL" -lt 1 ]; }; then FAILS="$FAILS parser_has_a_caller"; fi
# Served bytes unchanged: app.js, the CSP literal's file and the store's schema
for f in host/surfacepage/app/app.js host/daemon/surface_page.go host/store/credentials.go; do
  if git diff --quiet "$B" HEAD -- "$f"; then echo "UNCHANGED: $f" | tee -a "$L/legs.txt"; else echo "CHANGED: $f" | tee -a "$L/legs.txt"; FAILS="$FAILS served_bytes_changed"; fi
done
# live Chrome smoke: app.js is unchanged, so the Needs-you pair, XC2 and the bar stay green on it
XB="TestNeedsYouListsPendingRequests TestNeedsYouNotifiesOnlyWhenAskedOnce TestConfirmPageConfirmsWithAPasskeyInChrome TestVerbBarEpisodeReportInChrome TestDecidePageApproveWithPasskey"
CH="'/Applications/Google Chrome.app/Contents/MacOS/Google Chrome'"
leg chrome "set -o pipefail; WORLD_CHROME=$CH WORLD_REQUIRE_PASSKEY_BROWSER=1 go test ./cmd/ailang-worldd/ -run '^($(bar $XB))\$' -count=1 -v"
need chrome $XB
if grep -E '^\s*--- (SKIP|FAIL)' "$L/chrome.txt"; then FAILS="$FAILS chrome_skip_or_fail"; fi
if [ "${GATES_SKIP_VERIFY_GO:-0}" = 1 ]; then echo "LEG verify_go SKIPPED (GATES_SKIP_VERIFY_GO=1)" | tee -a "$L/legs.txt"
else leg verify_go 'VERIFY_GO_LEGS=plain AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_go.sh'; fi
echo "GATES_DONE $(( $(date +%s) - s0 ))s FAILS:${FAILS:- none}" | tee -a "$L/legs.txt"
echo "GATES_DONE $(( $(date +%s) - s0 ))s FAILS:${FAILS:- none}" > "$L/ALL_DONE.$SHA"
````

---

## 7. Executor directive notes

- **No git writes of any kind** (no add, rm, mv, commit, stash, checkout, reset, branch). **`patch -p1 -E` (§2) is the
  only writer** of the tree; no hand edit, no `gofmt -w`, no golden regeneration (S1 has none), no fuzz. A failed
  `--dry-run` or a red acceptance → **stop and report** with the log; do not repair.
- **Never `nohup`, never background a test leg; pass a 600000 ms timeout to every long Bash call.** Each acceptance part is
  one call (C14 measured each). **The executor never launches Chrome** and never sets `WORLD_CHROME`.
- **Runs**: **one run, M1 → M2 → M3 → M4** (each patch assumes the previous tree; the controller may commit between
  sub-milestones). Parts, in order: **M1** `fast`; **M2** `fast`; **M3** `fast`; **M4** `fast`, `cmd1`, `cmd2`, `cmd3` (the three are the
  package whole). `accept.sh`
  exports `AILANG_BIN=$HOME/.pinned-ailang/ailang`, `WORLD_EXEC_NODE=/opt/homebrew/bin/node` (the node harnesses cannot skip)
  and `WORLD_TOOL_AILANG_BIN=$HOME/.pinned-ailang-tools/v0.54.15/ailang`, and unsets `WORLD_CHROME`,
  `WORLD_REQUIRE_PASSKEY_BROWSER` and every golden-update variable.
- **Commit-time test commands** (`commit_ms.sh <M> <msgfile> <test-cmd…>`; each carries its env through `accept.sh`; with
  `T=$HOME/.ailang/state/world-iter407/exec-patches W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter407-row243-v4s1`):
  **M1** `bash $T/accept.sh $W M1 fast`; **M2** `bash $T/accept.sh $W M2 fast`; **M3** `bash $T/accept.sh $W M3 fast`; **M4**
  `bash $T/accept.sh $W M4 fast`. The `cmd1`/`cmd2`/`cmd3` parts ran at acceptance and are not repeated at commit. `commit_ms.sh` vets
  the directories of the changed `.go` files only; every S1 file is `.go`. No commit command launches Chrome (S1-D11); the gates'
  `chrome` leg runs the existing live tests with both env vars and requires `--- PASS` per name.
- After each sub-milestone's acceptance, snapshot **cumulatively against the base** to
  `$HOME/.ailang/state/world-iter407/snap/<M>/`: every file added or changed since `d92983d5` (repo-relative paths), a
  **`DELETED`** file (one repo path per line; **empty for every S1 stage**), a **`STATUS`** file, and `MANIFEST.sha256`
  (the two plan files, `STATUS` and `DELETED` are not in it). **Never an `rsync` of whole directories.** The `find` takes
  only the top-level dirs that exist (M1–M3 have only `host/`; M4 has `cmd/` and `host/`; no stage has `docs/`), rehearsed on
  all four (C15):
  ```bash
  bash <<'EOF2'
  set -euo pipefail
  M=M1   # the sub-milestone just accepted
  B=d92983d59f887905251df5b99e825c299136760f
  W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter407-row243-v4s1
  D="$HOME/.ailang/state/world-iter407/snap/$M"; rm -rf "${D:?}"; mkdir -p "$D"
  cd "$W"
  { git diff --name-status --no-renames "$B" -- cmd host docs; git ls-files --others --exclude-standard -- cmd host docs | sed 's/^/A\t/'; } | sort -u -t "$(printf '\t')" -k2 > "$D/STATUS"
  awk -F'\t' '$1 != "D" {print $2}' "$D/STATUS" | while read -r p; do mkdir -p "$D/$(dirname "$p")"; cp -p "$p" "$D/$p"; done
  awk -F'\t' '$1 == "D" {print $2}' "$D/STATUS" > "$D/DELETED"
  (cd "$D" && find $(for d in cmd host docs; do [ -d "$d" ] && echo "$d"; done) -type f -print0 | sort -z | xargs -0 shasum -a 256) > "$D/MANIFEST.sha256"
  cat "$D/STATUS"; echo "DELETED: $(wc -l < "$D/DELETED")"; echo "MANIFEST: $(wc -l < "$D/MANIFEST.sha256")"
  EOF2
  ```
  Expected (`STATUS` lines / `DELETED` / MANIFEST): **M1** 1 / 0 / 1; **M2** 3 / 0 / 3; **M3** 4 / 0 / 4; **M4** 7 / 0 / 7. The
  controller copies each snapshot over the worktree, commits it as one commit, and checks the numstat (`git diff -M --numstat
  d92983d5 HEAD -- cmd host docs app | awk '{s+=$1+$2} END {print s}'`): **M1 127, M2 570, M3 778, M4 976**. A different
  number → stop.
- **Stop conditions**: a dry-run failure; any `--- FAIL`, `--- SKIP` or `FAIL` line in an acceptance (`ACCEPT_RED`); gofmt
  output; a numstat mismatch; any file outside §5 changed (`git status --short` after each patch must equal the cumulative list
  of §3's paragraphs); a numstat above 1,400.
- **Known load-sensitive tests** (the directive's list) are not this sprint's: a red there in `cmd1`/`cmd2`/`cmd3` is re-run once
  alone (`TestCLIRealSubprocessEpisode`, `TestCapsuleTemplateIsCopiedNotShared`, `TestDoctorQuickIsFast` — red at the base under
  load and green alone, §1.1 — `TestMarketingChromeSixShots` needs
  `WORLD_TOOL_AILANG_BIN`, which `accept.sh` exports).
- D-WORLD-109: nothing touches `~/.ailang/world`, `~/.ssh`, a real `/dev/tty` or a real key; no S1 test opens a store.

---

## 8. Risks

- **R1 — the seam is package state.** `episodeFlagSetSeen` is read by every `runEpisodeStart`/`runEpisodeFinishWith` call; the
  parity tests set it only inside `cliEpisodeFlagSet` (no `t.Parallel`, reset by `defer`), and Go runs sequential top-level
  tests before any paused parallel one, so no other test's CLI call can see it set. Under `-race` the three tests ran green
  (§1.5).
- **R2 — the mirror is a second copy of the CLI's flag names.** That is the point of the parity: a CLI flag change fails
  `TestOperatorFlagParity` by name (CLI-ADD/REMOVE/RENAME mutants, §4), and the fix is one mirror line.
- **R3 — the subset admits less than the CLI parses** (duplicates, explicit false, every override). Each is a refusal with a
  terminal line, never a different execution; S7 renders them.
- **R4 — S7 may want a field the answer lacks** (e.g. the resolved preset's display). `OperatorAnswer` is a struct; adding a
  field is additive, and the classes and `Intent` are fixed by design §4.1 (S1-D4).

---

## 9. Acceptance criteria (the judge scores against these)

1. **No new daemon route, no schema change, no `.ail`, `surfaceCSP` unchanged, no served-behaviour change**: `DAEMON_ROUTES` base
   = final, `SURFACE_POSTS: 3`, `UNTOUCHED_VIOLATIONS: none`, `UNCHANGED:` app.js / surface_page.go / host/store/credentials.go,
   `verify_ail` identical to the base's; `TestFrozenV1RouteTableMatchesMux`, `TestConfirmRoutesRegistered`,
   `TestDaemonCredentialWriterSites`, `TestSurfaceVerbIsReadOnly`, `TestSurfaceAppPure` green unedited.
2. **The parser has no production caller**: `PARSER_PRODUCTION_CALLERS: none` with the positive control ≥ 1; the only cmd
   production change is the seam (`CMD_SEAM_LINES: 10 0`), nil in production.
3. **Each judge-405 survivor is killed by a named test, whole package**: `J-JS-RE-NL` (`TestNeedsYouConfirmationEmbeddedNewline`),
   `J-JS-CS-NOGUARD` (`TestNeedsYouConfirmationsNonArray`), `J-JS-ALLOW-DECIDE` (`TestVerbDecidePrefixBoundary`); each reproduced
   SURVIVED on the base first (C9); every boundary sibling killed but the named equivalent `P1-RE-S`.
4. **The CLI-flag parity fails when a flag is added, renamed or removed on either side**: `CLI-ADD-FLAG`, `CLI-ADD-FINISH-FLAG`,
   `CLI-REMOVE-FLAG`, `CLI-RENAME-FLAG`, `SEAM-START-OFF`, `SEAM-FINISH-OFF` (cmd; `SEAM-AFTER-PARSE` is a named equivalent) and `MIR-ADD-zzz`,
   `MIR-DROP-hosts-dir`, `MIR-RENAME-budget`, the `MIR-KIND-*` and `OPT-*` rows (host/verbs) killed, each by a named test.
5. **Every grammar gets a mutant per boundary class and each is killed** (§4): verb words, flag names, flag values (incl. the
   twelve bool spellings), episode id, preset name and cap, action id — arity, prefix vs substring, delimiter, case, length
   (exact, one under, one over), each anchor, `(?m)` and `(?i)` each, input type; the `(?s)` mutants are named equivalents.
6. **The answers are an input**: `TestOperatorGrammarPolicy` (384 cases, 42 admitted) kills every `POL-*`, `MODE-*`, `FIN-*`
   mutant; the zero `OperatorPolicy` admits nothing.
7. **Every refusal asserts its class and that no tuple is admitted** (the `refused` helper; `admitArgv`'s Intent ⇔ review check).
8. **Every §3.1 row is driven by a named test.**
9. **Size**: `git diff -M --numstat d92983d5 HEAD -- cmd host docs app` sums to **976** (≤ 1,400; the same unfiltered);
   each sub-milestone's numstat as §7.
10. **Exactly these pre-existing files change**: `cmd/ailang-worldd/episode.go` (+7) and `episode_finish.go` (+3), the seam.
    **No pre-existing test file changes or is deleted.**
11. **Mutants**: every `K` row killed by its named test and by the whole package; the `E` rows named with reasons; the plan's
    ids are the ledger's.
12. **Gates**: §6's script `FAILS: none` (or only a directive-listed load-sensitive test, green alone), the `chrome` leg's five
    names `--- PASS`, `verify_go` rc 0.

---

## Appendix C — after this sprint (S1's close; what S2–S12 own)

- **Moved out: nothing** (S1-D1: 976 ≤ 1,400). No patch is parked.
- **Closed here**: judge-405's three survivors (`J-JS-RE-NL`, `J-JS-CS-NOGUARD`, `J-JS-ALLOW-DECIDE`) and their siblings; design
  §5.6 (all three pins; the live smoke is the existing tests, S1-D11); design §5.1's grammar and identifier rows
  (`TestOperatorGrammar{Case,Arity,ExactWords,Canonical,Duplicates,Options,Policy}`, `TestOperatorNeverDispatch`,
  `TestOperatorIdentifier{Anchors,Flags,Length}`, `TestOperatorOtherVerbParity` as far as a pure function goes); the CLI join
  (design §4.1 first paragraph).
- **§5.1/§5.6 rows S2–S12 own** (by design §8): **S2** the canonical action envelope and typed SSHSIG/passkey proof against fakes
  (§3.1 message/namespace/limits; `TestOperatorStrictJSON` and the envelope half of `TestOperatorRequestLimits`); **S3** schema 10
  `operator_actions` and its receipts (`TestOperatorAdmissionBoundaries`, `TestOperatorDuplicateTap`, `TestOperatorReconnect`'s
  durable dedup, `TestOperatorTerminalExclusion`); **S4** runner install and heartbeat, the inert catalogue fixture
  (`TestOperatorRunnerUnavailable`, `TestOperatorCatalogueAdmission`); **S7** the inert request and guarded confirm routes —
  `GET /surface/verb.json`'s NEW `review` outcome rendering `AdmitOperator`'s classes (S1-D4), `TestOperatorOtherVerbParity` at
  the route, `TestOperatorBarStartReview`, `TestOperatorBarFinishReview`, `TestOperatorMisdirected`, `TestOperatorCrossOrigin`,
  `TestOperatorFetchSite`, `TestOperatorMediaType`, `TestOperatorBodyBound`, `TestOperatorUnsupportedMethod`,
  `TestOperatorRequestLimits` (body/task), `TestOperatorLiteralTask`; **S8** the card and the app bar
  (`TestOperatorPanelNetwork`, `TestOperatorPanelLatest`, `TestOperatorLinkAllowlist`, `TestOperatorWireTypes` — and F8's
  array-element type confusion, which S8 must pin when it touches the shared app shell; `TestOperatorReconnect`'s browser half;
  the first Chrome slice); **S9–S11** finish (§5.4); **S12** the herdr hand-off (`HandoffProven` becomes true only there).
- **What S2 inherits from S1**: the admitted tuple (`OperatorIntent{Kind, Episode, Preset, Headless, ConfirmOnSurface, Ready,
  Argv}`) is the "command words and admitted options" the canonical action carries (design §3.1) — S2 canonicalizes it, it does
  not re-parse the bar line; `OperatorActionIDPattern` (`^oa-[0-9a-f]{32}$`) is the id S2/S3 mint and S7 routes; the policy stays
  an input (S2 is also answer-independent); the mirror and its parity test are the place a new admitted option is added (one
  line each side, design §4.1 first).
- **Residuals carried (audit iteration 406, for later sprints — not acted on in S1)**: an uncertain action has no attended
  resolver; a headless child outliving the runner is unowned; S3's schema 10 must stay start-only unless D-WORLD-120 is
  answered; the second credentials writer vs the `LOCK_NB` writer lock and the writer-site inventories
  (`TestDaemonCredentialWriterSites`, `TestCLICredentialWriterSites`) must move together; proof types in `host/authority` with
  unexported constructors; the herdr plugin hard-codes `codex=gpt-6.1-sol` (`world-config-init:96`).
- **r2 quorum notes (non-blocking; for the design's next revision, the controller records)**: §2.1 should say W33 ran in the
  driver's process tree, not the runner's; "non-blocking with a ≤ 5 s bound" should read `LOCK_NB`, with the 5 s bounding only
  the store read for the refusal message; glm's six symbols, verified by the controller at 406 — `runPresetEpisodeStart`
  (`episode_preset.go:53`), `launchEpisode` (`episode_launch.go:258`), `pollVerbRequest` (`episode_finish_surface.go:214`),
  `RecheckVerb` (`host/authority/judge_verb.go:107`), `systemServiceController` (`service_launchctl.go:24`), `serviceEnv.confirm`
  (`service.go:167`) — matter to S4–S11.
- **The V4a / V4b / herdr split (planner's view, one paragraph).** Agree, by measurement of this sprint: S1 shows the start and
  finish halves already separate cleanly at the grammar — the finish policy predicate is the only finish-specific line, and
  every finish row could be dormant without touching start's. V4a (S1–S8, start/launch, schema 10 start-only) is deliverable on
  D-WORLD-119 alone and ends with a browser-proven headless start; V4b (S9–S11) needs D-WORLD-120 and E2's attended
  confirmation and carries the largest recovery surface (remote reconciliation), so planning it after V4a's runner has run
  real (fake-child) restarts is cheaper than planning it now; the herdr child (S12) depends on E4 and on 119 = A only, so it
  should not block either. The split also bounds the S3 schema to what is authorized. The one cost is a second schema step if
  120 is answered A/B later; S3 should reserve the finish columns' names (not their DDL) so V4b's migration is additive.
- **Controller decisions requested**: none for S1. (D-WORLD-119/120 remain the ledger's; S2 is answer-independent too.)
- The controller records; this plan moves no design doc and edits no charter.

---

## Appendix P — the patches (apply with §2, in order; generated from the compiled scratch stages; do not edit)

### Patch M1 — the judge-405 pins

````diff M1
diff --git a/host/surfacepage/app_pins_test.go b/host/surfacepage/app_pins_test.go
new file mode 100644
index 00000000..2ccfa346
--- /dev/null
+++ b/host/surfacepage/app_pins_test.go
@@ -0,0 +1,127 @@
+package surfacepage
+
+// Row 243 V4 S1 (design §5.6): the three judge-405 pins on app.js's Needs you and the bar's navigate
+// allowlist, each a node harness over the daemon's own answers (requests-confirmations.golden.json,
+// written by host/daemon TestRequestsJSONConfirmationsGolden; verb-route.golden.json, checked by
+// TestSurfaceVerbOutcomes). Each harness pairs its hostile rows with the real row they are cut from.
+
+import (
+	"context"
+	"os"
+	"os/exec"
+	"path/filepath"
+	"strings"
+	"testing"
+	"time"
+)
+
+// runPinHarness runs one harness under node beside app.js and the two goldens, and requires its
+// manifest line. It skips only when no node exists; with WORLD_EXEC_NODE set, a harness that cannot
+// run fails.
+func runPinHarness(t *testing.T, harness, manifest string) {
+	t.Helper()
+	node := os.Getenv("WORLD_EXEC_NODE")
+	if node == "" {
+		var err error
+		node, err = exec.LookPath("node")
+		if err != nil {
+			t.Skip("node unavailable locally")
+		}
+	}
+	dir := t.TempDir()
+	files := map[string][]byte{"app.js": AppScript, "harness.js": []byte(pinPrelude + harness)}
+	for name, src := range map[string]string{"confirm.json": "requests-confirmations.golden.json", "verb.json": "verb-route.golden.json"} {
+		b, err := os.ReadFile(filepath.Join("testdata", src))
+		if err != nil {
+			t.Fatal(err)
+		}
+		files[name] = b
+	}
+	for name, data := range files {
+		if err := os.WriteFile(filepath.Join(dir, name), data, 0o600); err != nil {
+			t.Fatal(err)
+		}
+	}
+	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
+	defer cancel()
+	cmd := exec.CommandContext(ctx, node, filepath.Join(dir, "harness.js"))
+	cmd.Dir = dir
+	out, err := cmd.CombinedOutput()
+	if err != nil {
+		t.Fatalf("node harness: %v\n%s", err, out)
+	}
+	if !strings.Contains(string(out), manifest) {
+		t.Fatalf("harness manifest missing: %s", out)
+	}
+}
+
+const pinPrelude = `
+const assert = require('node:assert/strict');
+const fs = require('node:fs');
+const app = require('./app.js');
+const G = JSON.parse(fs.readFileSync('confirm.json', 'utf8')), V = JSON.parse(fs.readFileSync('verb.json', 'utf8'));
+const HEAD = (n) => '<div class="nsec" data-needs><div class="nh" style="color:var(--' + (n ? 'youShell' : 'sub') + ')">Needs you · ' + n + '</div>';
+const real = G.v9.confirmations[1];   // the handler's summarised draft row, vr-00000000000000c4
+const Q = 'ab'.repeat(32), dec = {request: 'sha256:' + Q, scope: 'worktree:1', requester: 'agent-1', deadline: 50520, reserved: '', refusal: ''};
+const decLink = '<div class="nl"><a id="decide-' + Q + '" href="/surface/decide/sha256:' + Q + '">Decide · worktree:1 · requested by agent-1 · deadline 14:02</a></div>';
+`
+
+// TestNeedsYouConfirmationEmbeddedNewline (judge-405 J-JS-RE-NL): a confirmation id that is a valid
+// id plus a line break and more, either side, is not a confirmation id — no link; nor is the id in
+// another case (each of the regex's flags and anchors on its own); the valid row and a valid decision
+// beside it still link.
+func TestNeedsYouConfirmationEmbeddedNewline(t *testing.T) {
+	runPinHarness(t, `
+assert.equal(real.id, 'vr-00000000000000c4');
+const ok = app.needsYouHTML({confirmations: [real], decisions: [dec]}, 0);
+assert(ok.includes(decLink) && ok.includes('<a id="confirm-' + real.id + '" href="/surface/confirm/' + real.id + '">'), ok);
+for (const id of [real.id + '\n', real.id + '\njunk', 'junk\n' + real.id, '\n' + real.id, real.id + '\r', real.id + '\r\nx',
+  real.id + '\u2028x', 'x\u2029' + real.id, real.id.toUpperCase(), 'VR' + real.id.slice(2), real.id.slice(0, 17) + 'C4']) {
+  const h = app.needsYouHTML({confirmations: [{...real, id}], decisions: [dec]}, 0);
+  assert.equal(h, HEAD(2) + decLink + '<div class="nl sub">Confirm · ' + app.esc(id) + ' · refused: not a confirmation id</div></div>', JSON.stringify(id));
+  assert(!h.includes('/surface/confirm/') && !h.includes('confirm-vr'), JSON.stringify(id));
+}
+console.log('NEEDS YOU NEWLINE PASS');
+`, "NEEDS YOU NEWLINE PASS")
+}
+
+// TestNeedsYouConfirmationsNonArray (judge-405 J-JS-CS-NOGUARD): a confirmations field that is not an
+// array — an object, an array-like, a string, null, a number, a boolean, a thing with a map — lists
+// nothing and throws nothing; the mint request and the decision of the same answer still render, exactly
+// as with no confirmations; the real array still lists.
+func TestNeedsYouConfirmationsNonArray(t *testing.T) {
+	runPinHarness(t, `
+const pending = [{id: 'mr-1', episode: 'e', agent: 'a', grants: 1, expires: 60}];
+const base = {...G.v9, pending, decisions: [dec]};
+const none = app.needsYouHTML({...base, confirmations: []}, 0);
+assert(none.startsWith(HEAD(2) + '<div class="nl"><a id="needs-mr-1"') && none.includes(decLink), none);
+assert(app.needsYouHTML(base, 0).startsWith(HEAD(5)), 'the real array lists');
+for (const bad of [{}, {0: real, length: 1}, real.id, '', null, 0, 1, true, {map: () => ['<a href="/surface/confirm/x">']}]) {
+  assert.equal(app.needsYouHTML({...base, confirmations: bad}, 0), none, JSON.stringify(bad));
+}
+const { confirmations: _, ...absent } = base;
+assert.equal(app.needsYouHTML(absent, 0), none);
+console.log('NEEDS YOU NON-ARRAY PASS');
+`, "NEEDS YOU NON-ARRAY PASS")
+}
+
+// TestVerbDecidePrefixBoundary (judge-405 J-JS-ALLOW-DECIDE): a navigate answer opens a decision card
+// only when its href STARTS with /surface/decide/sha256: — a foreign origin, a path that contains the
+// prefix, the prefix without its colon, or another case is refused; the same boundaries hold for the
+// approve and confirm prefixes. The daemon's own decide answer navigates.
+func TestVerbDecidePrefixBoundary(t *testing.T) {
+	runPinHarness(t, `
+const REF = 'approve sha256:' + Q, MR = 'approve mr-0123456789abcdef', VR = 'approve vr-0123456789abcdef';
+for (const q of [REF, MR, VR]) assert.deepEqual(app.verbOutcome(200, V[q]), {assign: JSON.parse(V[q]).href}, q);
+const b = JSON.parse(V[REF]), H = 'sha256:' + Q;
+assert.equal(b.href, '/surface/decide/' + H);
+for (const href of ['//evil/surface/decide/' + H, 'https://evil/surface/decide/' + H, '/x/surface/decide/' + H,
+  '/surface/approve/x/surface/decide/' + H, ' /surface/decide/' + H, '/surface/decide/sha256' + Q, '/surface/decide/sha256x' + Q,
+  '/surface/decide/sha256/' + Q, '/surface/decide/SHA256:' + Q, '/Surface/decide/' + H,
+  '/surface/approve/mr0123456789abcdef', '/surface/approve/mrx0123456789abcdef', '/surface/approve/MR-0123456789abcdef',
+  '/x/surface/approve/mr-0123456789abcdef', '/SURFACE/confirm/vr-0123456789abcdef']) {
+  assert.deepEqual(app.verbOutcome(200, JSON.stringify({...b, href})), {say: 'refused: unexpected target'}, href);
+}
+console.log('DECIDE PREFIX PASS');
+`, "DECIDE PREFIX PASS")
+}
````

### Patch M2 — the admitted subset and its grammar rows

````diff M2
diff --git a/host/verbs/operator.go b/host/verbs/operator.go
new file mode 100644
index 00000000..76ffac2d
--- /dev/null
+++ b/host/verbs/operator.go
@@ -0,0 +1,227 @@
+package verbs
+
+// Row 243 V4 S1 (design §4.1, §5.1): the bar's admitted operator subset — the episode start and
+// episode finish lines the app may later show a review for (S7). It is pure and, in S1, has no caller:
+// it reads Lookup's rows and runs ParseInterleaved over a flag set that mirrors the CLI's own
+// (cmd/ailang-worldd TestOperatorFlagParity holds the mirror to the CLI's real flag sets). The ledger
+// answers it depends on arrive as an input (OperatorPolicy); nothing unanswered is admitted.
+
+import (
+	"errors"
+	"flag"
+	"fmt"
+	"io"
+	"regexp"
+)
+
+// OperatorFlag is one flag of the CLI's episode start or episode finish flag set: its name, whether
+// it is a boolean (it takes no value word), and whether the app's subset admits it.
+type OperatorFlag struct {
+	Name  string
+	Bool  bool
+	Admit bool
+}
+
+// EpisodeStartFlags and EpisodeFinishFlags are the CLI's two flag sets, every flag in the CLI's order
+// (runEpisodeStart, runEpisodeFinishWith); Admit marks design §4.1's subset.
+var (
+	EpisodeStartFlags = []OperatorFlag{
+		{"db", false, false}, {"credentials-db", false, false}, {"workspace-root", false, false}, {"repo", false, false},
+		{"branch", false, false}, {"profile", false, false}, {"exec-profiles-dir", false, false}, {"task", false, false},
+		{"agent", false, false}, {"model", false, false}, {"ttl", false, false}, {"preset", false, true},
+		{"presets-dir", false, false}, {"resume", false, false}, {"hosts-dir", false, false}, {"launch", true, true},
+		{"headless", true, true}, {"request", true, false}, {"grant", false, false}, {"budget", false, false},
+		{"sessions-dir", false, false}, {"addr", false, false},
+	}
+	EpisodeFinishFlags = []OperatorFlag{
+		{"preset", false, true}, {"presets-dir", false, false}, {"title", false, false}, {"include-unattributed", true, false},
+		{"ready", true, true}, {"exec-profiles-dir", false, false}, {"dry-run", true, false}, {"confirm-on-surface", true, true},
+		{"allow-failing", false, false},
+	}
+)
+
+// EpisodeIDPattern is the CLI's and the daemon's episode grammar (cmd workspaceEpisodePattern, host/daemon
+// episodeIDPattern): 1 to 64 characters. MaxPresetName is the CLI's preset-name cap (cmd maxPresetNameLen).
+// OperatorActionIDPattern is design §3.1's NEW action id, oa- and 32 lowercase hex digits; S1 only defines it.
+var (
+	EpisodeIDPattern        = regexp.MustCompile(`^[a-z0-9][a-z0-9-]{0,63}$`)
+	OperatorActionIDPattern = regexp.MustCompile(`^oa-[0-9a-f]{32}$`)
+)
+
+const MaxPresetName = 51
+
+// ValidOperatorActionID reports whether id is an operator action id.
+func ValidOperatorActionID(id string) bool { return OperatorActionIDPattern.MatchString(id) }
+
+// Answer is the ledger's answer to one decision; only "A", "B" and "C" are answers.
+type Answer string
+
+// OperatorPolicy is the caller's answers, a pure input: Start is D-WORLD-119's, Finish D-WORLD-120's,
+// HandoffProven whether design §11 E4 has passed. The zero value — nothing answered — admits nothing.
+type OperatorPolicy struct {
+	Start, Finish Answer
+	HandoffProven bool
+}
+
+func (p OperatorPolicy) headless() bool    { return p.Start == "A" || p.Start == "B" }
+func (p OperatorPolicy) interactive() bool { return p.Start == "A" && p.HandoffProven }
+func (p OperatorPolicy) finish() bool      { return p.headless() && (p.Finish == "A" || p.Finish == "B") }
+func (p OperatorPolicy) ready() bool       { return p.Finish == "A" } // asked only after finish()
+
+// OperatorClass is AdmitOperator's classification, the typed outcome a later GET /surface/verb.json
+// renders (S7): review a supported tuple; a usage refusal; a terminal line (outside the subset, or a mode
+// the answers do not allow); a never family; or another line whose Lookup classification stands.
+type OperatorClass string
+
+const (
+	OperatorReview   OperatorClass = "review"
+	OperatorUsage    OperatorClass = "usage"
+	OperatorTerminal OperatorClass = "terminal"
+	OperatorNever    OperatorClass = "never"
+	OperatorOther    OperatorClass = "other"
+)
+
+// OperatorIntent is one admitted tuple; Argv is its canonical CLI line's words.
+type OperatorIntent struct {
+	Kind                    string // "start" or "finish"
+	Episode, Preset         string
+	Headless                bool // start: --launch --headless, else --launch (interactive)
+	ConfirmOnSurface, Ready bool // finish
+	Argv                    []string
+}
+
+// OperatorAnswer is AdmitOperator's answer. Intent is set exactly when Class is OperatorReview; Field
+// names the word a refusal is about ("" for none), Reason says why, CLI is the line as typed.
+type OperatorAnswer struct {
+	Class              OperatorClass
+	Intent             *OperatorIntent
+	Field, Reason, CLI string
+}
+
+// boolWords are the flag package's boolean spellings (strconv.ParseBool's), each with its value; host/verbs
+// imports no strconv (host/daemon TestSurfaceVerbIsReadOnly pins its imports).
+var boolWords = map[string]bool{"1": true, "t": true, "T": true, "TRUE": true, "true": true, "True": true,
+	"0": false, "f": false, "F": false, "FALSE": false, "false": false, "False": false}
+
+// operatorValue counts a flag's occurrences and keeps its last value (a bool's as on).
+type operatorValue struct {
+	isBool, on bool
+	n          int
+	s          string
+}
+
+func (v *operatorValue) String() string   { return v.s }
+func (v *operatorValue) IsBoolFlag() bool { return v.isBool }
+func (v *operatorValue) Set(s string) error {
+	if v.isBool {
+		b, ok := boolWords[s]
+		if !ok {
+			return errors.New("parse error") // the flag package's own bool error, word for word
+		}
+		v.on = b
+	}
+	v.n++
+	v.s = s
+	return nil
+}
+
+// AdmitOperator classifies one bar line's words (Split's) for the operator subset under policy p
+// (design §4.1). It never runs, reads or writes anything.
+func AdmitOperator(argv []string, p OperatorPolicy) OperatorAnswer {
+	if len(argv) > 0 && argv[0] == "world-publish" { // a separate binary, outside Lookup's word sets
+		return OperatorAnswer{Class: OperatorNever, Field: "world-publish", CLI: Quote(argv),
+			Reason: "world-publish is never run from the app: publishing is typed at a terminal; at a terminal: " + Quote(argv)}
+	}
+	v, rest := Lookup(argv)
+	var mirror []OperatorFlag
+	switch v.Words {
+	case "episode start":
+		mirror = EpisodeStartFlags
+	case "episode finish":
+		mirror = EpisodeFinishFlags
+	default:
+		if v.Class == Never {
+			return OperatorAnswer{Class: OperatorNever, Reason: v.Reason, CLI: v.CLI}
+		}
+		return OperatorAnswer{Class: OperatorOther, Reason: v.Reason, CLI: v.CLI}
+	}
+	refuse := func(c OperatorClass, field, why string) OperatorAnswer {
+		return OperatorAnswer{Class: c, Field: field, Reason: Binary + " " + v.Words + ": " + why, CLI: v.CLI}
+	}
+	at := "; at a terminal: " + v.CLI
+	fs := flag.NewFlagSet(Binary+" "+v.Words, flag.ContinueOnError)
+	fs.SetOutput(io.Discard)
+	got := map[string]*operatorValue{}
+	for _, f := range mirror {
+		got[f.Name] = &operatorValue{isBool: f.Bool}
+		fs.Var(got[f.Name], f.Name, "")
+	}
+	pos, err := ParseInterleaved(fs, rest)
+	if errors.Is(err, flag.ErrHelp) {
+		return refuse(OperatorUsage, "--help", "see "+Binary+" "+v.Words+" --help at a terminal")
+	}
+	if err != nil {
+		return refuse(OperatorUsage, "", err.Error())
+	}
+	for _, f := range mirror {
+		if got[f.Name].n > 1 {
+			return refuse(OperatorUsage, "--"+f.Name, "--"+f.Name+" is given more than once")
+		}
+	}
+	for _, f := range mirror {
+		if g := got[f.Name]; g.n == 1 && !f.Admit {
+			return refuse(OperatorTerminal, "--"+f.Name, "--"+f.Name+" is not run from the app"+at)
+		}
+		if g := got[f.Name]; g.n == 1 && f.Bool && !g.on {
+			return refuse(OperatorTerminal, "--"+f.Name, "--"+f.Name+" is given false; the app runs only its true form"+at)
+		}
+	}
+	if len(pos) != 1 {
+		return refuse(OperatorUsage, "<episode>", "needs exactly one <episode>")
+	}
+	if !EpisodeIDPattern.MatchString(pos[0]) {
+		return refuse(OperatorUsage, "<episode>", fmt.Sprintf("episode %q does not match %s", pos[0], EpisodeIDPattern))
+	}
+	preset := got["preset"].s
+	if preset == "" {
+		return refuse(OperatorUsage, "--preset", "--preset <name> is required")
+	}
+	if preset == "se-tools" || len(preset) > MaxPresetName || !EpisodeIDPattern.MatchString(preset) {
+		return refuse(OperatorUsage, "--preset", fmt.Sprintf("preset %q is not a launch preset name (the episode grammar, at most %d characters)", preset, MaxPresetName))
+	}
+	in := &OperatorIntent{Episode: pos[0], Preset: preset}
+	if v.Words == "episode start" {
+		launch, headless := got["launch"].on, got["headless"].on
+		switch {
+		case headless && !launch:
+			return refuse(OperatorUsage, "--headless", "--headless requires --launch")
+		case !launch:
+			return refuse(OperatorTerminal, "--launch", "the app starts an episode only with --launch"+at)
+		case headless && !p.headless():
+			return refuse(OperatorTerminal, "--headless", "starting from the app waits for D-WORLD-119 (A or B)"+at)
+		case !headless && !p.interactive():
+			return refuse(OperatorTerminal, "--launch", "an interactive launch from the app needs D-WORLD-119 A and the herdr hand-off proven (design §11 E4)"+at)
+		}
+		in.Kind, in.Headless = "start", headless
+		in.Argv = []string{Binary, "episode", "start", in.Episode, "--preset", preset, "--launch"}
+		if headless {
+			in.Argv = append(in.Argv, "--headless")
+		}
+		return OperatorAnswer{Class: OperatorReview, Intent: in, CLI: v.CLI}
+	}
+	in.Kind, in.ConfirmOnSurface, in.Ready = "finish", got["confirm-on-surface"].on, got["ready"].on
+	if !p.finish() {
+		return refuse(OperatorTerminal, "", "finishing from the app waits for D-WORLD-119 and D-WORLD-120 (A or B)"+at)
+	}
+	if in.Ready && !p.ready() {
+		return refuse(OperatorTerminal, "--ready", "--ready from the app needs D-WORLD-120 A"+at)
+	}
+	in.Argv = []string{Binary, "episode", "finish", in.Episode, "--preset", preset}
+	if in.ConfirmOnSurface {
+		in.Argv = append(in.Argv, "--confirm-on-surface")
+	}
+	if in.Ready {
+		in.Argv = append(in.Argv, "--ready")
+	}
+	return OperatorAnswer{Class: OperatorReview, Intent: in, CLI: v.CLI}
+}
diff --git a/host/verbs/operator_test.go b/host/verbs/operator_test.go
new file mode 100644
index 00000000..393d8f89
--- /dev/null
+++ b/host/verbs/operator_test.go
@@ -0,0 +1,216 @@
+package verbs
+
+// Row 243 V4 S1 (design §5.1's grammar and identifier rows): the admitted operator subset. Every line
+// is cut by Split, as the bar cuts it; every refusal asserts its class AND that no tuple is admitted.
+// The expectations are this file's literals, never the production tables.
+
+import (
+	"reflect"
+	"strconv"
+	"strings"
+	"testing"
+)
+
+var allYes = OperatorPolicy{Start: "A", Finish: "A", HandoffProven: true}
+
+// admitLine splits line and admits it; Intent must be set exactly for a review.
+func admitLine(t *testing.T, line string, p OperatorPolicy) OperatorAnswer {
+	t.Helper()
+	argv, err := Split(line)
+	if err != nil {
+		t.Fatalf("Split(%q): %v", line, err)
+	}
+	return admitArgv(t, argv, p)
+}
+
+func admitArgv(t *testing.T, argv []string, p OperatorPolicy) OperatorAnswer {
+	t.Helper()
+	a := AdmitOperator(argv, p)
+	if (a.Class == OperatorReview) != (a.Intent != nil) {
+		t.Fatalf("AdmitOperator(%q) = class %q with intent %+v", argv, a.Class, a.Intent)
+	}
+	return a
+}
+
+// refused requires class c, field (when not "-") and a reason containing has, with no tuple.
+func refused(t *testing.T, a OperatorAnswer, what string, c OperatorClass, field, has string) {
+	t.Helper()
+	if a.Class != c || a.Intent != nil || (field != "-" && a.Field != field) || !strings.Contains(a.Reason, has) {
+		t.Errorf("%s: got class %q field %q reason %q intent %+v; want %q %q containing %q, no intent", what, a.Class, a.Field, a.Reason, a.Intent, c, field, has)
+	}
+}
+
+const (
+	start  = "episode start w1 --preset p --launch --headless"
+	finish = "episode finish w1 --preset p"
+)
+
+func TestOperatorGrammarCase(t *testing.T) {
+	if a := admitLine(t, start, allYes); a.Class != OperatorReview {
+		t.Fatalf("control: %q is %q", start, a.Class)
+	}
+	for _, line := range []string{"Episode start w1 --preset p --launch --headless", "episode Start w1 --preset p --launch --headless",
+		"EPISODE START w1 --preset p --launch --headless", "episode FINISH w1 --preset p", "episode Finish w1 --preset p"} {
+		argv, _ := Split(line)
+		v, _ := Lookup(argv)
+		refused(t, admitLine(t, line, allYes), line, OperatorOther, "", v.Reason)
+	}
+	for _, line := range []string{"episode start w1 --Preset p --launch --headless", "episode start w1 --preset p --LAUNCH --headless",
+		"episode start w1 --preset p --launch --Headless", "episode start w1 --PRESET=p --launch --headless",
+		"episode finish w1 --preset p --Ready", "episode finish w1 --preset p --Confirm-On-Surface"} {
+		refused(t, admitLine(t, line, allYes), line, OperatorUsage, "", "flag provided but not defined")
+	}
+	refused(t, admitLine(t, "episode start W1 --preset p --launch --headless", allYes), "W1", OperatorUsage, "<episode>", "does not match")
+	refused(t, admitLine(t, "episode start w1 --preset P --launch --headless", allYes), "P", OperatorUsage, "--preset", "not a launch preset name")
+	refused(t, admitLine(t, "episode finish wA --preset p", allYes), "wA", OperatorUsage, "<episode>", "does not match")
+}
+
+func TestOperatorGrammarArity(t *testing.T) {
+	for _, line := range []string{"episode start --preset p --launch --headless", "episode start w1 w2 --preset p --launch --headless",
+		"episode start w1 --preset p --launch --headless extra", "episode start --preset p w1 --launch w2 --headless",
+		"episode finish --preset p", "episode finish w1 w2 --preset p", "episode finish w1 --preset p --ready x"} {
+		refused(t, admitLine(t, line, allYes), line, OperatorUsage, "<episode>", "needs exactly one <episode>")
+	}
+	for _, line := range []string{"episode start w1 --launch --headless", "episode start w1 --preset '' --launch --headless",
+		"episode finish w1", "episode finish w1 --preset= --ready"} {
+		refused(t, admitLine(t, line, allYes), line, OperatorUsage, "--preset", "--preset <name> is required")
+	}
+	for _, line := range []string{"episode start w1 --launch --headless --preset", "episode finish w1 --preset", "episode finish w1 --ready --preset"} {
+		refused(t, admitLine(t, line, allYes), line, OperatorUsage, "", "flag needs an argument")
+	}
+	refused(t, admitLine(t, "episode start '' --preset p --launch --headless", allYes), "empty episode", OperatorUsage, "<episode>", "does not match")
+}
+
+func TestOperatorGrammarExactWords(t *testing.T) {
+	for _, line := range []string{"episode starts w1 --preset p --launch --headless", "episode star w1 --preset p --launch --headless",
+		"episode start-x w1 --preset p --launch --headless", "episodes start w1 --preset p --launch --headless",
+		"xepisode start w1 --preset p --launch --headless", "episode finis w1 --preset p", "episode finished w1 --preset p",
+		"episode launch w1 --preset p --headless", "episode report w1", "episode", "start w1 --preset p --launch --headless"} {
+		argv, _ := Split(line)
+		v, _ := Lookup(argv)
+		refused(t, admitLine(t, line, allYes), line, OperatorOther, "", v.Reason)
+	}
+	for _, x := range []string{"--pres p", "--presets p", "--preset-dir p", "--head", "--headles", "--headlessx", "--launc", "--launched",
+		"--read", "--readyy", "--confirm", "--confirm-on-surface-x", "--dry", "--ttl-x 1"} {
+		line := start + " " + x
+		refused(t, admitLine(t, line, allYes), line, OperatorUsage, "", "flag provided but not defined")
+	}
+	refused(t, admitLine(t, "episode start w1 ---preset p --launch --headless", allYes), "---preset", OperatorUsage, "", "bad flag syntax")
+}
+
+func TestOperatorGrammarCanonical(t *testing.T) {
+	head := []string{Binary, "episode", "start", "w1", "--preset", "p", "--launch", "--headless"}
+	cases := []struct {
+		lines []string
+		want  OperatorIntent
+	}{
+		{[]string{start, "ailang-worldd " + start, "episode start --preset=p --launch --headless w1", "episode start -preset p -launch -headless w1",
+			"episode start --headless --launch w1 --preset 'p'", `episode start "w1" --preset="p" --launch=true --headless=1`,
+			"episode start --launch --preset p w1 --headless", "episode start w1 --preset p --launch=T --headless=TRUE"},
+			OperatorIntent{Kind: "start", Episode: "w1", Preset: "p", Headless: true, Argv: head}},
+		{[]string{"episode start w1 --preset p --launch", "episode start --launch=true w1 -preset=p"},
+			OperatorIntent{Kind: "start", Episode: "w1", Preset: "p", Argv: head[:7]}},
+		{[]string{finish, "episode finish --preset=p w1", "ailang-worldd episode finish -preset p w1"},
+			OperatorIntent{Kind: "finish", Episode: "w1", Preset: "p", Argv: []string{Binary, "episode", "finish", "w1", "--preset", "p"}}},
+		{[]string{"episode finish --confirm-on-surface w1 --preset p", "episode finish w1 --preset p --confirm-on-surface=true"},
+			OperatorIntent{Kind: "finish", Episode: "w1", Preset: "p", ConfirmOnSurface: true, Argv: []string{Binary, "episode", "finish", "w1", "--preset", "p", "--confirm-on-surface"}}},
+		{[]string{"episode finish w1 --ready --preset p --confirm-on-surface", "episode finish --confirm-on-surface=1 --ready=true w1 --preset=p"},
+			OperatorIntent{Kind: "finish", Episode: "w1", Preset: "p", ConfirmOnSurface: true, Ready: true, Argv: []string{Binary, "episode", "finish", "w1", "--preset", "p", "--confirm-on-surface", "--ready"}}},
+		{[]string{"episode finish w1 --ready --preset my-preset-2"},
+			OperatorIntent{Kind: "finish", Episode: "w1", Preset: "my-preset-2", Ready: true, Argv: []string{Binary, "episode", "finish", "w1", "--preset", "my-preset-2", "--ready"}}},
+	}
+	for _, c := range cases {
+		for _, line := range c.lines {
+			a := admitLine(t, line, allYes)
+			if a.Class != OperatorReview || !reflect.DeepEqual(*a.Intent, c.want) {
+				t.Errorf("%q: class %q intent %+v, want %+v", line, a.Class, a.Intent, c.want)
+				continue
+			}
+			argv, _ := Split(line)
+			if v, _ := Lookup(argv); a.CLI != v.CLI || v.CLI == "" {
+				t.Errorf("%q: CLI %q", line, a.CLI)
+			}
+			if again := admitArgv(t, a.Intent.Argv, allYes); again.Class != OperatorReview || !reflect.DeepEqual(*again.Intent, c.want) {
+				t.Errorf("%q: the canonical line %q re-admits as %+v", line, a.Intent.Argv, again.Intent)
+			}
+		}
+	}
+}
+
+func TestOperatorGrammarDuplicates(t *testing.T) {
+	for _, c := range []struct{ line, field string }{
+		{"episode start w1 --preset p --preset p --launch --headless", "--preset"},
+		{"episode start w1 --preset p -preset=q --launch --headless", "--preset"},
+		{"episode start w1 --preset=p --launch --launch --headless", "--launch"},
+		{"episode start w1 --preset p --launch=true --launch --headless", "--launch"},
+		{"episode start w1 --preset p --launch --headless -headless", "--headless"},
+		{"episode finish w1 --preset p --ready --ready", "--ready"},
+		{"episode finish w1 --preset p --confirm-on-surface --confirm-on-surface=true", "--confirm-on-surface"},
+		{"episode finish w1 --preset p --title a --title b", "--title"},
+	} {
+		refused(t, admitLine(t, c.line, allYes), c.line, OperatorUsage, c.field, "is given more than once")
+	}
+}
+
+func TestOperatorGrammarOptions(t *testing.T) {
+	startOut := map[string]bool{"db": false, "credentials-db": false, "workspace-root": false, "repo": false, "branch": false, "profile": false,
+		"exec-profiles-dir": false, "task": false, "agent": false, "model": false, "ttl": false, "presets-dir": false, "resume": false,
+		"hosts-dir": false, "request": true, "grant": false, "budget": false, "sessions-dir": false, "addr": false}
+	finishOut := map[string]bool{"presets-dir": false, "title": false, "include-unattributed": true, "exec-profiles-dir": false,
+		"dry-run": true, "allow-failing": false}
+	for _, set := range []struct {
+		base string
+		out  map[string]bool
+	}{{start, startOut}, {finish, finishOut}} {
+		for name, isBool := range set.out {
+			x, eq := "--"+name+" v", "--"+name+"=v" // =v is not a boolean: a value flag mirrored as a bool fails here
+			if isBool {
+				x, eq = "--"+name, "--"+name+"=true"
+			}
+			for _, line := range []string{set.base + " " + x, strings.Replace(set.base, " w1 ", " "+x+" w1 ", 1), set.base + " " + eq} {
+				a := admitLine(t, line, allYes)
+				refused(t, a, line, OperatorTerminal, "--"+name, "--"+name+" is not run from the app; at a terminal: "+a.CLI)
+				if verb := strings.Fields(set.base)[1]; !strings.HasPrefix(a.Reason, "ailang-worldd episode "+verb+": --"+name+" ") {
+					t.Errorf("%q: reason %q does not lead with the command", line, a.Reason)
+				}
+			}
+		}
+	}
+	for _, c := range []struct{ line, field string }{
+		{"episode start w1 --preset p --launch=false --headless", "--launch"},
+		{"episode start w1 --preset p --launch --headless=false", "--headless"},
+		{"episode start w1 --preset p --launch=f --headless=0", "--launch"},
+		{"episode finish w1 --preset p --ready=false", "--ready"},
+		{"episode finish w1 --preset p --confirm-on-surface=FALSE", "--confirm-on-surface"},
+		{"episode finish w1 --preset p --dry-run=false", "--dry-run"},
+	} {
+		has := "is given false"
+		if c.field == "--dry-run" {
+			has = "is not run from the app"
+		}
+		refused(t, admitLine(t, c.line, allYes), c.line, OperatorTerminal, c.field, has)
+	}
+	refused(t, admitLine(t, "episode start w1 --preset p", allYes), "no --launch", OperatorTerminal, "--launch", "only with --launch; at a terminal: ")
+	refused(t, admitLine(t, "episode start w1 --preset p --headless", allYes), "standalone --headless", OperatorUsage, "--headless", "--headless requires --launch")
+	refused(t, admitLine(t, "episode start w1 --preset p --launch=maybe", allYes), "not a bool", OperatorUsage, "",
+		`ailang-worldd episode start: invalid boolean value "maybe" for -launch: parse error`) // the flag package's words
+	refused(t, admitLine(t, "episode start w1 --preset se-tools --launch --headless", allYes), "se-tools", OperatorUsage, "--preset", "not a launch preset name")
+	// Every boolean spelling decides as strconv.ParseBool (the flag package) does: true admits, false is
+	// refused as false, anything else is a parse error.
+	for _, w := range []string{"1", "t", "T", "TRUE", "true", "True", "0", "f", "F", "FALSE", "false", "False", "", "2", "yes", "tRUE", "true ", "on"} {
+		a, line := admitLine(t, "episode start w1 --preset p --launch --headless='"+w+"'", allYes), "--headless="+w
+		switch b, err := strconv.ParseBool(w); {
+		case err != nil:
+			refused(t, a, line, OperatorUsage, "", "invalid boolean value")
+		case b:
+			if a.Class != OperatorReview {
+				t.Errorf("%s: %q %q", line, a.Class, a.Reason)
+			}
+		default:
+			refused(t, a, line, OperatorTerminal, "--headless", "is given false")
+		}
+	}
+	for _, h := range []string{"-h", "--help", "-help"} {
+		refused(t, admitLine(t, start+" "+h, allYes), h, OperatorUsage, "--help", "--help at a terminal")
+	}
+}
````

### Patch M3 — policy, never-dispatch, identifiers, the other verbs

````diff M3
diff --git a/host/verbs/operator_policy_test.go b/host/verbs/operator_policy_test.go
new file mode 100644
index 00000000..f935eff8
--- /dev/null
+++ b/host/verbs/operator_policy_test.go
@@ -0,0 +1,208 @@
+package verbs
+
+// Row 243 V4 S1 (design §5.1's policy, never-dispatch, identifier and other-verb rows), on operator_test.go's
+// helpers: the answers as an input, the families never admitted, the three identifier grammars at each
+// boundary class, and every other line keeping Lookup's classification.
+
+import (
+	"encoding/json"
+	"fmt"
+	"os"
+	"path/filepath"
+	"strings"
+	"testing"
+)
+
+// policyOracle is design §4.1/§10's answer table, written out here (not read from the code).
+func policyOracle(p OperatorPolicy, what string) bool {
+	ab := func(a Answer) bool { return a == "A" || a == "B" }
+	switch what {
+	case "headless":
+		return ab(p.Start)
+	case "interactive":
+		return p.Start == "A" && p.HandoffProven
+	case "finish":
+		return ab(p.Start) && ab(p.Finish)
+	}
+	return ab(p.Start) && p.Finish == "A" // ready
+}
+
+func TestOperatorGrammarPolicy(t *testing.T) {
+	lines := []struct{ what, line, field string }{
+		{"headless", start, "--headless"},
+		{"interactive", "episode start w1 --preset p --launch", "--launch"},
+		{"finish", finish + " --confirm-on-surface", ""},
+		{"ready", finish + " --ready", ""},
+	}
+	n, admitted := 0, 0
+	for _, s := range []Answer{"", "A", "B", "C", "a", "b", "D", "A "} {
+		for _, f := range []Answer{"", "A", "B", "C", "a", "D"} {
+			for _, h := range []bool{false, true} {
+				p := OperatorPolicy{Start: s, Finish: f, HandoffProven: h}
+				for _, l := range lines {
+					a, want := admitLine(t, l.line, p), policyOracle(p, l.what)
+					n++
+					if want {
+						admitted++
+						if a.Class != OperatorReview {
+							t.Errorf("%+v %s: %q %q, want review", p, l.what, a.Class, a.Reason)
+						}
+						continue
+					}
+					field := l.field
+					if l.what == "ready" && policyOracle(p, "finish") {
+						field = "--ready"
+					}
+					refused(t, a, l.what+" "+string(s)+"/"+string(f), OperatorTerminal, field, "; at a terminal: ")
+				}
+			}
+		}
+	}
+	if n != 384 || admitted != 42 {
+		t.Fatalf("%d cases, %d admitted; want 384 and 42", n, admitted)
+	}
+	for _, l := range lines { // the zero value answers nothing and admits nothing
+		refused(t, admitLine(t, l.line, OperatorPolicy{}), "zero policy "+l.what, OperatorTerminal, "-", "D-WORLD-1")
+	}
+}
+
+func TestOperatorNeverDispatch(t *testing.T) {
+	for _, line := range []string{"approver enroll --name mark-key", "approver revoke mark-key", "session mint w1", "session revoke w1",
+		"service install", "service", "serve", "setup", "ailang-worldd approver enroll", "ailang-worldd service uninstall",
+		"world-publish approve sha256:x", "world-publish"} {
+		a := admitLine(t, line, allYes)
+		refused(t, a, line, OperatorNever, "-", "is never run from the app")
+		if !strings.Contains(a.Reason, "; at a terminal: ") || a.CLI == "" {
+			t.Errorf("%q: reason %q CLI %q", line, a.Reason, a.CLI)
+		}
+	}
+	for _, line := range []string{"sh -c ls", "bash", "episode launch w1 --preset p --headless", "episode report w1", "call x", "log", "why head",
+		"approve mr-0123456789abcdef", "world-publishx", "commit"} {
+		argv, _ := Split(line)
+		v, _ := Lookup(argv)
+		refused(t, admitLine(t, line, allYes), line, OperatorOther, "", v.Reason)
+	}
+	for _, line := range []string{start + "; rm -rf x", start + " $(id)", start + " | sh", "episode start w1 --preset `p` --launch --headless"} {
+		if argv, err := Split(line); err != ErrShell {
+			t.Errorf("Split(%q) = %q, %v; want ErrShell (the bar runs no shell)", line, argv, err)
+		}
+	}
+}
+
+func TestOperatorIdentifierAnchors(t *testing.T) {
+	ep := func(id string) OperatorAnswer {
+		return admitArgv(t, []string{"episode", "start", id, "--preset", "p", "--launch", "--headless"}, allYes)
+	}
+	pr := func(id string) OperatorAnswer {
+		return admitArgv(t, []string{"episode", "finish", "w1", "--preset", id}, allYes)
+	}
+	if ep("w1").Class != OperatorReview || pr("p").Class != OperatorReview || !ValidOperatorActionID("oa-"+strings.Repeat("0a", 16)) {
+		t.Fatal("controls: w1, p and oa-0a… must be admitted")
+	}
+	for _, id := range []string{"_w1", "w1_", ".w1", "w1.", "w1/", "/w1", " w1", "w1 ", "é1", "w1é"} {
+		refused(t, ep(id), "episode "+id, OperatorUsage, "<episode>", "does not match")
+		refused(t, pr(id), "preset "+id, OperatorUsage, "--preset", "not a launch preset name")
+	}
+	hex := strings.Repeat("0a", 16)
+	for _, id := range []string{"x" + "oa-" + hex, "oa-" + hex + "x", " oa-" + hex, "oa-" + hex + " ", "-oa-" + hex, "oa-" + hex + "-", "oa" + hex, "oa_" + hex, "oa-" + "-" + hex[1:]} {
+		if ValidOperatorActionID(id) {
+			t.Errorf("action id %q admitted", id)
+		}
+	}
+}
+
+func TestOperatorIdentifierFlags(t *testing.T) {
+	hex := strings.Repeat("0a", 16)
+	for _, id := range []string{"w1\nzz", "zz\nw1", "w1\n", "\nw1", "W1", "wQ"} {
+		refused(t, admitArgv(t, []string{"episode", "start", id, "--preset", "p", "--launch", "--headless"}, allYes), "episode "+id, OperatorUsage, "<episode>", "does not match")
+		refused(t, admitArgv(t, []string{"episode", "finish", "w1", "--preset", id}, allYes), "preset "+id, OperatorUsage, "--preset", "not a launch preset name")
+	}
+	for _, id := range []string{"oa-" + hex + "\nx", "x\noa-" + hex, "oa-" + hex + "\n", "oa-" + strings.ToUpper(hex), "OA-" + hex, "Oa-" + hex} {
+		if ValidOperatorActionID(id) {
+			t.Errorf("action id %q admitted", id)
+		}
+	}
+}
+
+func TestOperatorIdentifierLength(t *testing.T) {
+	ep := func(id string) OperatorAnswer {
+		return admitArgv(t, []string{"episode", "start", id, "--preset", "p", "--launch", "--headless"}, allYes)
+	}
+	pr := func(id string) OperatorAnswer {
+		return admitArgv(t, []string{"episode", "finish", "w1", "--preset", id}, allYes)
+	}
+	for _, id := range []string{"a", "0", "a" + strings.Repeat("-", 63), strings.Repeat("a", 64)} {
+		if a := ep(id); a.Class != OperatorReview || a.Intent.Episode != id {
+			t.Errorf("episode of %d characters: %q %q", len(id), a.Class, a.Reason)
+		}
+	}
+	for _, id := range []string{"", strings.Repeat("a", 65), "a" + strings.Repeat("-", 64), "-"} {
+		refused(t, ep(id), fmt.Sprintf("episode of %d characters", len(id)), OperatorUsage, "<episode>", "does not match")
+	}
+	for _, id := range []string{"p", strings.Repeat("p", 51)} {
+		if a := pr(id); a.Class != OperatorReview || a.Intent.Preset != id {
+			t.Errorf("preset of %d characters: %q %q", len(id), a.Class, a.Reason)
+		}
+	}
+	for _, id := range []string{strings.Repeat("p", 52), strings.Repeat("p", 64), strings.Repeat("p", 65)} {
+		refused(t, pr(id), fmt.Sprintf("preset of %d characters", len(id)), OperatorUsage, "--preset", "at most 51 characters")
+	}
+	hex := strings.Repeat("0a", 16)
+	if !ValidOperatorActionID("oa-"+hex) || len("oa-"+hex) != 35 {
+		t.Fatal("control: oa- and 32 hex digits")
+	}
+	for _, id := range []string{"oa-" + hex[:31], "oa-" + hex + "0", "oa-", "oa-" + hex[:30], "oa-" + hex + "00"} {
+		if ValidOperatorActionID(id) {
+			t.Errorf("action id of %d hex digits admitted", len(id)-3)
+		}
+	}
+}
+
+// TestOperatorOtherVerbParity: every line of GET /surface/verb.json's golden corpus that is not episode
+// start or episode finish keeps Lookup's classification — never a review, the same reason and CLI line.
+func TestOperatorOtherVerbParity(t *testing.T) {
+	b, err := os.ReadFile(filepath.Join("..", "surfacepage", "testdata", "verb-route.golden.json"))
+	if err != nil {
+		t.Fatal(err)
+	}
+	var golden map[string]string
+	if err := json.Unmarshal(b, &golden); err != nil {
+		t.Fatal(err)
+	}
+	lines := []string{}
+	for q := range golden {
+		lines = append(lines, q)
+	}
+	for _, w := range append(append(append([]string{}, OperatorWords...), NeverWords...), OtherWords...) {
+		lines = append(lines, w, w+" x")
+	}
+	n := 0
+	for _, line := range lines {
+		argv, err := Split(line)
+		if err != nil {
+			continue
+		}
+		v, _ := Lookup(argv)
+		a := admitArgv(t, argv, allYes)
+		if v.Words == "episode start" || v.Words == "episode finish" {
+			if a.Class == OperatorOther || a.Class == OperatorNever {
+				t.Errorf("%q: %q is an operator line", line, a.Class)
+			}
+			continue
+		}
+		want := OperatorOther
+		if v.Class == Never {
+			want = OperatorNever
+		}
+		if a.Class != want || a.Reason != v.Reason || a.CLI != v.CLI || a.Field != "" {
+			t.Errorf("%q: %+v, want class %q with Lookup's reason %q and CLI %q", line, a, want, v.Reason, v.CLI)
+		}
+		n++
+	}
+	if n < 70 {
+		t.Fatalf("%d lines kept Lookup's classification; want at least 70", n)
+	}
+	if a := admitLine(t, "episode finish w1 --preset p", allYes); a.Class != OperatorReview {
+		t.Errorf("the golden's finish line under every answer A: %q", a.Class)
+	}
+}
````

### Patch M4 — the join to the CLI

````diff M4
diff --git a/cmd/ailang-worldd/episode.go b/cmd/ailang-worldd/episode.go
index 2f803a6b..d6fa2a7a 100644
--- a/cmd/ailang-worldd/episode.go
+++ b/cmd/ailang-worldd/episode.go
@@ -187,6 +187,10 @@ func runEpisode(addr string, args []string, stdout, stderr io.Writer) int {
 	}
 }
 
+// episodeFlagSetSeen, when set (tests only), is shown each episode start and finish flag set as built,
+// before any parse: row 243 V4 S1's TestOperatorFlagParity holds host/verbs' mirror to them.
+var episodeFlagSetSeen func(fs *flag.FlagSet)
+
 func runEpisodeStart(args []string, stdout, stderr io.Writer, env sessionEnv) int {
 	def := defaultWorldPaths()
 	fs := flag.NewFlagSet("ailang-worldd episode start", flag.ContinueOnError)
@@ -215,6 +219,9 @@ func runEpisodeStart(args []string, stdout, stderr io.Writer, env sessionEnv) in
 	budget := fs.Int64("budget", 50, "budget of each preset grant")
 	sessionsDir := fs.String("sessions-dir", defaultSessionsDir(), "where the session file and the MCP config are written")
 	addr := fs.String("addr", defaultServiceBind(), "the serving daemon, host:port")
+	if episodeFlagSetSeen != nil {
+		episodeFlagSetSeen(fs)
+	}
 	positional, err := verbs.ParseInterleaved(fs, args)
 	if err != nil {
 		if errors.Is(err, flag.ErrHelp) {
diff --git a/cmd/ailang-worldd/episode_finish.go b/cmd/ailang-worldd/episode_finish.go
index 076ed573..da0d65f7 100644
--- a/cmd/ailang-worldd/episode_finish.go
+++ b/cmd/ailang-worldd/episode_finish.go
@@ -74,6 +74,9 @@ func runEpisodeFinishWith(args []string, stdout, stderr io.Writer, env sessionEn
 	onSurface := fs.Bool("confirm-on-surface", false, "confirm the pull request with an enrolled key on the surface")
 	var allowed multiFlag
 	fs.Var(&allowed, "allow-failing", "allow this failing check")
+	if episodeFlagSetSeen != nil {
+		episodeFlagSetSeen(fs)
+	}
 	pos, err := verbs.ParseInterleaved(fs, args)
 	if errors.Is(err, flag.ErrHelp) {
 		fmt.Fprint(stdout, episodeFinishHelp)
diff --git a/cmd/ailang-worldd/operator_parity_test.go b/cmd/ailang-worldd/operator_parity_test.go
new file mode 100644
index 00000000..20fdb960
--- /dev/null
+++ b/cmd/ailang-worldd/operator_parity_test.go
@@ -0,0 +1,188 @@
+package main
+
+// Row 243 V4 S1 (design §4.1, W23, W29): host/verbs' admitted operator subset is joined to the CLI's
+// own flag sets and identifier grammars. The flag sets are the ones runEpisodeStart and
+// runEpisodeFinishWith build, seen through episodeFlagSetSeen on a --help run, which returns before
+// any effect.
+
+import (
+	"bytes"
+	"context"
+	"flag"
+	"reflect"
+	"sort"
+	"strings"
+	"testing"
+
+	"github.com/sunholo-data/ailang-world/host/verbs"
+)
+
+// cliEpisodeFlagSet runs the CLI's own `episode <verb> --help` and returns the flag set it built.
+func cliEpisodeFlagSet(t *testing.T, verb string) *flag.FlagSet {
+	t.Helper()
+	t.Setenv("HOME", t.TempDir())
+	var seen []*flag.FlagSet
+	episodeFlagSetSeen = func(fs *flag.FlagSet) { seen = append(seen, fs) }
+	defer func() { episodeFlagSetSeen = nil }()
+	var out, errOut bytes.Buffer
+	rc := -1
+	switch verb {
+	case "start":
+		rc = runEpisodeStart([]string{"--help"}, &out, &errOut, sessionEnv{})
+	case "finish":
+		rc = runEpisodeFinishWith([]string{"--help"}, &out, &errOut, sessionEnv{}, nil)
+	}
+	if rc != exitOK || len(seen) != 1 || out.Len() == 0 || errOut.Len() != 0 {
+		t.Fatalf("episode %s --help: rc %d, %d flag sets seen, stderr %q", verb, rc, len(seen), errOut.String())
+	}
+	return seen[0]
+}
+
+// flagKinds is a flag set's names, each true when the flag is a boolean.
+func flagKinds(fs *flag.FlagSet) map[string]bool {
+	m := map[string]bool{}
+	fs.VisitAll(func(f *flag.Flag) {
+		b, ok := f.Value.(interface{ IsBoolFlag() bool })
+		m[f.Name] = ok && b.IsBoolFlag()
+	})
+	return m
+}
+
+func mirrorKinds(fl []verbs.OperatorFlag) (m map[string]bool, admitted []string) {
+	m = map[string]bool{}
+	for _, f := range fl {
+		m[f.Name] = f.Bool
+		if f.Admit {
+			admitted = append(admitted, f.Name)
+		}
+	}
+	sort.Strings(admitted)
+	return m, admitted
+}
+
+// TestOperatorFlagParity: host/verbs' mirror is exactly the CLI's flag set — every name, each kind —
+// so the CLI adding, renaming or removing a flag, or the subset naming one the CLI lacks, fails here.
+func TestOperatorFlagParity(t *testing.T) {
+	for _, c := range []struct {
+		verb   string
+		mirror []verbs.OperatorFlag
+		n      int
+		admit  []string
+	}{
+		{"start", verbs.EpisodeStartFlags, 22, []string{"headless", "launch", "preset"}},
+		{"finish", verbs.EpisodeFinishFlags, 9, []string{"confirm-on-surface", "preset", "ready"}},
+	} {
+		cli := flagKinds(cliEpisodeFlagSet(t, c.verb))
+		mir, admitted := mirrorKinds(c.mirror)
+		if len(cli) != c.n || len(c.mirror) != len(mir) || !reflect.DeepEqual(cli, mir) {
+			t.Errorf("episode %s: the CLI's %d flags %v; host/verbs' %d (%d entries) %v", c.verb, len(cli), cli, len(mir), len(c.mirror), mir)
+		}
+		if !reflect.DeepEqual(admitted, c.admit) {
+			t.Errorf("episode %s: the subset admits %v, want %v (design §4.1)", c.verb, admitted, c.admit)
+		}
+		// Controls: one flag more, one fewer, one renamed, one kind flipped — each differs.
+		for i, edit := range []func(m map[string]bool){
+			func(m map[string]bool) { m["zzz"] = true },
+			func(m map[string]bool) { delete(m, "preset") },
+			func(m map[string]bool) { delete(m, "preset"); m["presets"] = false },
+			func(m map[string]bool) { m["preset"] = true },
+		} {
+			m := map[string]bool{}
+			for k, v := range cli {
+				m[k] = v
+			}
+			if edit(m); reflect.DeepEqual(m, mir) {
+				t.Errorf("episode %s: control %d passed the parity check", c.verb, i)
+			}
+		}
+	}
+}
+
+// TestOperatorParsesLikeTheCLI: each admitted spelling, parsed by the CLI's own flag set with the same
+// ParseInterleaved, gives the episode and flag values of the admitted tuple, and so does the tuple's
+// canonical line; a line the CLI's parse rejects is refused with the CLI's own parse error.
+func TestOperatorParsesLikeTheCLI(t *testing.T) {
+	yes := verbs.OperatorPolicy{Start: "A", Finish: "A", HandoffProven: true}
+	cliParse := func(verb string, words []string) (map[string]string, []string, error) {
+		fs := cliEpisodeFlagSet(t, verb)
+		pos, err := verbs.ParseInterleaved(fs, words)
+		vals := map[string]string{}
+		for _, n := range []string{"preset", "launch", "headless", "ready", "confirm-on-surface"} {
+			if f := fs.Lookup(n); f != nil {
+				vals[n] = f.Value.String()
+			}
+		}
+		return vals, pos, err
+	}
+	for _, line := range []string{
+		"episode start w1 --preset p --launch --headless", "episode start --headless --launch w1 --preset=p",
+		`episode start -preset "p" -launch=true -headless=1 w1`, "episode start w1 --preset p --launch",
+		"episode finish w1 --preset p", "episode finish --confirm-on-surface --preset=p w1 --ready",
+		"episode finish w1 --ready=T --preset 'p'",
+	} {
+		argv, err := verbs.Split(line)
+		if err != nil {
+			t.Fatal(err)
+		}
+		a := verbs.AdmitOperator(argv, yes)
+		if a.Class != verbs.OperatorReview {
+			t.Fatalf("%q: %q %q", line, a.Class, a.Reason)
+		}
+		in := a.Intent
+		want := map[string]string{"preset": in.Preset}
+		if in.Kind == "start" {
+			want["launch"], want["headless"] = "true", map[bool]string{true: "true", false: "false"}[in.Headless]
+		} else {
+			want["ready"], want["confirm-on-surface"] = map[bool]string{true: "true", false: "false"}[in.Ready], map[bool]string{true: "true", false: "false"}[in.ConfirmOnSurface]
+		}
+		for _, words := range [][]string{argv[2:], in.Argv[3:]} {
+			vals, pos, err := cliParse(in.Kind, words)
+			if err != nil || !reflect.DeepEqual(pos, []string{in.Episode}) || !reflect.DeepEqual(vals, want) {
+				t.Errorf("%q: the CLI's parse of %q is %v %q (%v), the subset's tuple %v %q", line, words, vals, pos, err, want, in.Episode)
+			}
+		}
+	}
+	for _, line := range []string{"episode start w1 --Preset p --launch", "episode start w1 --launch --preset", "episode start w1 --preset p --launch=maybe",
+		"episode finish w1 --preset p --ready=x", "episode finish w1 --dryrun --preset p"} {
+		argv, _ := verbs.Split(line)
+		_, _, cliErr := cliParse(argv[1], argv[2:])
+		a := verbs.AdmitOperator(argv, yes)
+		if cliErr == nil || a.Class != verbs.OperatorUsage || a.Intent != nil || !strings.HasSuffix(a.Reason, ": "+cliErr.Error()) {
+			t.Errorf("%q: the CLI's parse says %v; the subset %q %q", line, cliErr, a.Class, a.Reason)
+		}
+	}
+}
+
+// TestOperatorIdentifierGrammarsMatchTheCLI: the subset's episode grammar decides exactly as the CLI's
+// on every boundary class, and its preset cap is the CLI's — where loadLaunchPreset itself draws it.
+func TestOperatorIdentifierGrammarsMatchTheCLI(t *testing.T) {
+	n := 0
+	for _, id := range []string{"", "a", "0", "-a", "a-", "A", "aB", "a_b", "a.b", "é", "a b", "a\n", "\na", "a\nb", strings.Repeat("a", 51),
+		strings.Repeat("a", 52), strings.Repeat("a", 64), strings.Repeat("a", 65), "a" + strings.Repeat("-", 63), "a" + strings.Repeat("-", 64)} {
+		if verbs.EpisodeIDPattern.MatchString(id) != workspaceEpisodePattern.MatchString(id) {
+			t.Errorf("%q: host/verbs %v, the CLI %v", id, verbs.EpisodeIDPattern.MatchString(id), workspaceEpisodePattern.MatchString(id))
+		}
+		if workspaceEpisodePattern.MatchString(id) {
+			n++
+		}
+	}
+	if n != 7 {
+		t.Errorf("%d of the corpus match the CLI's grammar, want 7", n)
+	}
+	if verbs.MaxPresetName != maxPresetNameLen {
+		t.Errorf("host/verbs' preset cap %d, the CLI's %d", verbs.MaxPresetName, maxPresetNameLen)
+	}
+	dir := t.TempDir()
+	yes := verbs.OperatorPolicy{Start: "A", Finish: "A"}
+	for _, c := range []struct {
+		n       int
+		invalid bool
+	}{{maxPresetNameLen, false}, {maxPresetNameLen + 1, true}, {51, false}, {52, true}} {
+		name := strings.Repeat("p", c.n)
+		_, err := loadLaunchPreset(context.Background(), name, dir, dir, dir)
+		a := verbs.AdmitOperator([]string{"episode", "finish", "w1", "--preset", name}, yes)
+		if err == nil || strings.Contains(err.Error(), "invalid preset name") != c.invalid || (a.Class == verbs.OperatorReview) == c.invalid {
+			t.Errorf("a %d-character preset: the CLI %v; the subset %q %q", c.n, err, a.Class, a.Reason)
+		}
+	}
+}
````

