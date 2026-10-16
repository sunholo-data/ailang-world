# Sprint plan — w-marketing-capture

**Iteration 266 · queue row 150 · base `3e6f6ae` · plan only**

Design: [w-marketing-capture.md](w-marketing-capture.md). Read in full, including V1–V28, r1–r3 quorum dispositions and r4 bounded revision. R1/r2/r3 remain BLOCKED rounds; r4 uses the controller narrow-refinement carve-out, not a retroactive quorum pass. Binding: CLAUDE.md and coding-standards.md; package-first law, contracts before pure implementation, host effects, non-vacuous named gates and AILANG fluency before writing source.

Reference pair requested under planned/ is absent at base; both existing files were read under design_docs/implemented/. JSON copies its top-level structure and milestone files/tests_first/acceptance/mutations shape, adding the requested explicit commands, acceptance_criteria, dependency, estimate and status="pending" fields.

**17 executor runs / 445 estimated minutes.** Design M1–M8 purposes and dependency order are preserved; large stages are split below. Publish confirmation and session mint typed `y` are outside all milestones. Planner writes only this file and [sprint_w-marketing-capture.json](sprint_w-marketing-capture.json).

## Environment, gate semantics and binary discovery

```sh
export AILANG_BIN=$HOME/.pinned-ailang/ailang
export WORLD_TOOL_AILANG_BIN=$HOME/.pinned-ailang-tools/v0.52.1/ailang
run_ail() { python3 -c 'import subprocess,sys; sys.exit(subprocess.run(sys.argv[1:],timeout=40).returncode)' "$@"; }
```

The pin was measured as v0.41.0/24ee108. The separate tool binary is v0.52.1 and is present here. host/daemon/workspace.go:44 requires it; host/daemon/workspace_test.go:593 locates WORLD_TOOL_AILANG_BIN and skips if unset. host/daemon/setools_e2e_test.go:58-86 and cmd/ailang-worldd/cliwalk_e2e_test.go:33-44 require both env vars or explicitly skip; configured wrong versions fail. TestSeToolsE2EBinaryVersions checks both. FakeToolBin-based workspace unit tests need no real tool but cannot prove S3.

New TestDemoSeedReadTestCheck, TestDemoSeedApprovalPacket, TestDemoSeedEndToEnd and real Chrome seeded drill need both pins because they compose real read/CLI/plan stages. TestDemoProof* adapters need v0.41.0 only. Required gate preflight must require executable pins/version and top-level named PASS; ordinary environment-gated SKIP is recorded but is not acceptance. Missing proof pin in required proof/package tests is fatal. Never substitute v0.41.0 for ToolAilangBin. All new exec probes are bounded, including version commands.

Every focused command below is literal and anchored to one test. Run with shown timeout and -v; require exactly that test executes and prints top-level PASS. Exit0/no tests/SKIP is not green. AILANG helper gives40s outer bound; ai-check also has5s solver bound. verify_ail outer supervisor limit300s; final build+test outer900s with600s per Go package; race outer540s with480s per package. Report bound exhaustion separately from an assertion verdict. No bind-denied focused result was encountered at base; any future `bind: operation not permitted` is **UNINFORMATIVE UNDER SANDBOX**, never PASS/FAIL.

## Baseline — measured in this worktree

All Go rows below used the two env exports above. AILANG commands used the exact pin and Python timeout40s; verify_ail had300s bound. Future new tests are absent at base and were not run as false-green empty selections.

| Result | Exact command | Wall seconds / detail |
|---|---|---|
| PASS | `go test -run '^TestDecisionPacketCodec$' ./host/coordinator -count=1 -timeout 45s -v` | 1.77  |
| PASS | `go test -run '^TestProducerMintsNoPacketForMultiEffectPlan$' ./host/coordinator -count=1 -timeout 45s -v` | 1.54  |
| PASS | `go test -run '^TestSessionMint_RequiresControllingTerminal$' ./cmd/ailang-worldd -count=1 -timeout 45s -v` | 2.06  |
| PASS | `go test -run '^TestSeToolsE2EBinaryVersions$' ./host/daemon -count=1 -timeout 45s -v` | 1.54  |
| PASS | `go test -run '^TestSeToolsA2AEndToEnd$' ./host/daemon -count=1 -timeout 120s -v` | 6.91  |
| PASS | `go test -run '^TestCLIAgainstSeToolsDaemon$' ./cmd/ailang-worldd -count=1 -timeout 120s -v` | 6.31  |
| PASS | `go test -run '^TestProducerEndToEndMintsEnvelopeValidatorAccepts$' ./host/evidence -count=1 -timeout 120s -v` | 2.44  |
| PASS | `go test -run '^TestRESTCommitRefusesDecisionPacketObjects$' ./host/daemon -count=1 -timeout 45s -v` | 1.26  |
| PASS | `go test -run '^TestDaemonInvocationMintsDecisionPacket$' ./host/daemon -count=1 -timeout 45s -v` | 0.68  |
| PASS | `go test -run '^TestWhyTargetFormsResolveToOneEntry$' ./cmd/ailang-worldd -count=1 -timeout 45s -v` | 0.69  |
| PASS | `go test -run '^TestWorkbenchPacketValidated$' ./host/daemon -count=1 -timeout 45s -v` | 0.75  |
| PASS | `go test -run '^TestWorkbenchPacketEvidenceOnlyThroughRowOneOhFive$' ./host/daemon -count=1 -timeout 45s -v` | 1.23  |
| PASS | `go test -run '^TestWorkbenchGradeProvenFromSubjectBoundEnvelope$' ./host/daemon -count=1 -timeout 45s -v` | 1.36  |
| PASS | `go test -run '^TestWorkbenchGraphWiring$' ./host/daemon -count=1 -timeout 45s -v` | 0.94  |
| PASS | `go test -run '^TestWorkbenchDecisionsPane$' ./host/daemon -count=1 -timeout 45s -v` | 0.71  |
| PASS | `go test -run '^TestWorkbenchLiveCursorCrossesGap$' ./host/daemon -count=1 -timeout 45s -v` | 0.76  |
| PASS | `go test -run '^TestWorkbenchLiveScriptTag$' ./host/daemon -count=1 -timeout 45s -v` | 0.76  |
| PASS | `go test -run '^TestRenderQuietFooter$' ./host/workbench -count=1 -timeout 45s -v` | 0.67  |
| PASS | `go test -run '^TestQuickstartDecisionPacketSection$' ./host/daemon -count=1 -timeout 45s -v` | 0.68  |
| SKIP (not acceptance) | `go test -run '^TestWorkbenchLiveChromeDrill$' ./host/daemon -count=1 -timeout 45s -v` | 0.66  |
| PASS | `go test -run '^TestWorkbenchViewFieldsAllRender$' ./host/workbench -count=1 -timeout 45s -v` | 0.34  |
| PASS | `go test -run '^TestWorkbenchNewViewTypesInCensus$' ./host/workbench -count=1 -timeout 45s -v` | 0.32  |
| PASS | `go test -run '^TestServeWorkspaceFlags$' ./cmd/ailang-worldd -count=1 -timeout 45s -v` | 0.69  |
| PASS | `go test -run '^TestServeProofFlags$' ./cmd/ailang-worldd -count=1 -timeout 45s -v` | 0.69  |
| PASS | `AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_ail.sh` | 27.43 16 required identities,40 named tests,444 se-tools tests; world package9/9 |
| PASS | `$HOME/.pinned-ailang/ailang check world/logepoch.ail` | 0.02 check passed |
| PASS | `$HOME/.pinned-ailang/ailang test --format json world/logepoch.ail` | 0.05 10 passed,zero failed/skipped/vacuous |
| PASS | `$HOME/.pinned-ailang/ailang ai-check -timeout 5s world/logepoch.ail` | 0.03 sameRef and servesEntry verified |
| FAIL — BASE FINDING (4 no_generator/vacuous skips) | `$HOME/.pinned-ailang/ailang test --format json world/contracts.ail` | 0.08 success=false,passed=7,skipped=4,vacuous_skips=4 |
| PASS | `go build ./...` | 0.78  |
| PASS | `go test -race ./host/daemon -count=1 -timeout 480s` | 339.09s wall,335.410s package; concurrent broad-suite load |
| PASS | `go test -count=1 -timeout 45s ./host/daemon -run '^TestWorkbench(GraphWiring\|DecisionsPane\|LiveCursorCrossesGap)$'` | 0.91 ok  	github.com/sunholo-data/ailang-world/host/daemon	0.380s  |
| FAIL — BASE FINDING | `go build ./... && go test ./... -count=1 -timeout 600s` | 279.93s; seven failing tests, listed below; no bind denial |
| PASS | `go test -race ./host/workbench -count=1 -timeout 120s` | 1.99s; ok  	github.com/sunholo-data/ailang-world/host/workbench	1.669s |
| PASS | `./scripts/tracked_binary_gate.sh` | 0.39s; ✓ 0 binary blobs among 3523 tracked files |

The full repository gate is **FAIL at base**, rc1,279.93s; build itself passed. Observed failures:

- `./cmd/ailang-worldd / TestSessionNewReadsTheAnswerFromThePtyWhenStdinIsDevNull`: open /dev/tty: operation not permitted; no controlling terminal.
- `./cmd/world-publish / TestTransitionsReadsTheConfirmationFromThePtyWhenStdinIsDevNull`: open /dev/tty: operation not permitted; STOP fence=tty.
- `./host/broker / TestOrdinaryOverflowCarriesNoESRCH`: group leader never became an unreaped zombie.
- `./host/broker / TestExecHandlerTimeoutKillsTheGroup`: specific child readiness before real kill: observation error.
- `./host/broker / TestExecSrtConfinementMatrix`: write HOME .world-m2-exec-decoy: operation not permitted.
- `./host/capsule / TestOrdinaryOverflowCarriesNoESRCH`: group leader never became an unreaped zombie.
- `./host/ptytest / TestTheTwoShapesOfControllingTerminal`: prompt did not appear within10s; empty pty transcript.

Terminal/HOME denials are concrete sandbox restrictions. They are distinct from the user’s socket-bind rule; no test logged `bind: operation not permitted`. Process readiness/zombie failures are reported as observed FAIL, with cause unestablished. The broad run overlapped the daemon race run; no claim of isolation or permanent regression is made. **Do not route around a red base gate** by excluding tests or changing protected/unrelated code. Report to controller and rerun unchanged mandatory gate in a suitable environment after findings are resolved.

V18 finding is reproduced: standalone contracts test exits1 with4 no_generator/vacuous skips. This does not make the aggregate repo AIL gate red; its named inline-test policy is different. It cannot be called successful TESTED evidence, and this sprint may not edit world/contracts.ail.

## Design defect and measured coupling findings

**PD1 — fixed shot dimensions conflict with visible required content.** Design line222 says “`workbench-light.png | /workbench | 1440×1100 | selected head, nonempty live list, graph and pending approval`”; line223 requires the same data in dark. V25 at line74 says “`the graph and decisions panes below the fold`” for that viewport. Re-read host/workbench/render.go:184,208-246 confirms one-column layout and timeline/inspector/live before graph/decisions. The design itself requires viewport-only shots and measured content rectangles (lines249-256); suppressing its single script cannot move regions into view. Tooling can proceed, but C8 product visibility remains unresolved until a controller-measured correction is specified. No silent layout redesign, changed viewport or omitted rectangle is planned. This is a successful planner finding, not a claimed acceptance pass. **RESOLVED by the controller (iteration 266, design line 222-223 amended):** the workbench-light/dark shots require only the in-viewport header and nonempty timeline; graph.png (`?from=N&entry=N`) and packet.png (`?object=<packet>`) carry the graph and the pending approval. No viewport, layout or rectangle is silently changed elsewhere; M7b/M7c measure the light/dark rectangles over the header+timeline region.

**Couplings made explicit:** verify_ail ROOTS at scripts/verify_ail.sh:165-169 omit world-demo; permanent Go pinned-manifest tests cover new sources without changing fleet/global gates. Page/roster structures must update the existing census (host/workbench/render_test.go:625-650, host/workbench/live_test.go:123-146). PNG CRC tests call the structural parser before image/png; directory8MiB cap is tested separately because six capped1MiB images cannot cross it.

## Release boundary and Chrome coverage

All named Go tests, package check/test/ai-check commands, fake Chrome scripts, real-pinned producer/coordinator integrations, observations and repository gates may run unattended in disposable test stores. Tests use only V11 credentials: testMintEnv (cmd/ailang-worldd/session_test.go:29), test-only authority.Mint provisioning (host/daemon/daemon_test.go:1053-1063), existing injected coordinator Now (host/daemon/decisionpacket_daemon_test.go:69-77). The seed itself has no Mint, injected-clock, fixture-key or fake-runner seam. Fake phase runners remain refusal/unit controls only.

Descriptor publish confirmation and production session mint typed y are attended release steps, never executor milestones. After the operator has done them, normal seed/capture/verify can run unattended on that frozen release store with the real key; this planner run has no such session. Controller later executes actual six-shot C8/C9 and actual C10 image links, promotes raw verified bytes and generates the sibling capture-manifest.json. Captions/docs/schema can ship before that; no placeholder actual manifest or release image passes a gate. Fleet committed-blob gate runs only after controller commit.

TestMarketingChromeSixShots explicitly checks WORLD_CHROME; unset ordinary-suite run SKIPs, invalid configured executable/version FAILs. It also requires both pins and controller-measured WORLD_MARKETING_RECTS. Its always-on companion TestMarketingCaptureFakeChromeSixShots exercises full supervisor/orchestration with a fake Chrome script in CI. The individual supervisor success/no-file/timeout/daemon-cleanup tests are all fake-script, browser-free and always-on. Existing TestWorkbenchLiveChromeDrill remains opt-in at host/daemon/workbench_live_chrome_test.go:136-145; its baseline SKIP is disclosed. Recorder HTML tests need no Chrome. No product pixel pass is inferred from a fake PNG.

Attended release acceptance (environment variables are operator-provided private scratch paths):

```sh
run_ail "$HOME/.pinned-ailang/ailang" --version
run_ail "$WORLD_TOOL_AILANG_BIN" --version
go build -o "$WORLD_CAPTURE_BIN" ./cmd/ailang-worldd
"$WORLD_CAPTURE_BIN" demo seed --db "$DEMO_DB" --session "$DEMO_SESSION" --ailang-bin "$HOME/.pinned-ailang/ailang" --proof-key-file "$DEMO_PROOF_KEY" --addr "$DEMO_ADDR" --manifest "$DEMO_MANIFEST"
scripts/capture_marketing.sh --db "$DEMO_DB" --manifest "$DEMO_MANIFEST" --proof-key-file "$DEMO_PROOF_KEY" --ailang-bin "$HOME/.pinned-ailang/ailang" --workspace-root "$DEMO_WORKSPACES" --tool-ailang-bin "$WORLD_TOOL_AILANG_BIN" --chrome "$WORLD_CHROME" --rectangles "$WORLD_MARKETING_RECTS" --out "$CAPTURE_OUT" --verify-twice
go run ./scripts/marketingcheck check --pass-a "$CAPTURE_OUT/pass-a" --pass-b "$CAPTURE_OUT/pass-b" --manifest "$DEMO_MANIFEST" --rectangles "$WORLD_MARKETING_RECTS"
go run ./scripts/marketingcheck docs --root . --assets "$CAPTURE_OUT/pass-a" --require-assets
./scripts/tracked_binary_gate.sh # controller AFTER committed promotion
```

Seed whole orchestration is supervised with300s bound; source-proof leg90s. Capture shell supervises both≤300s passes and reports progress per shot; use660s outer guard. Publisher/mint commands from design lines147-152 are preparation requiring real attended confirmation and are deliberately not in this acceptance block. Unminted TestDemoSeedPreflightRefusals must refuse.

## C11 observations and mutation discipline

**Snapshot-path correction (iteration267):** charter snapshot guardrail applies to every milestone directive below, in Markdown and JSON. Use $HOME/.ailang/state/world-iter<N>/snap/M<k>/ outside the checkout for source copies and OBS_OUT; never create Go snapshots inside the worktree. Location only changes; gates/protected-file policy do not.

Protected OBSERVATION uses SHA256 plus newline count on both sides, never writes protected originals. Include recursive tools/launchd subtree (12 paths measured, including testdata). At each run, capture observation before and after, compare sorted path/hash/line triples and scoped read-only diff. Preserve initial sprint observation for final comparison; any inherited difference is a finding, not permission to restore. The copied-input comparator mutation lives in TestMarketingProtectedObservationComparator; inventory mutation makes a genuine extra commit in a disposable test store.

Exact per-boundary observation command (run before/after with OBS_OUT naming a file under the external snapshot directory):

```sh
python3 - <<'PY' > "$OBS_OUT"
import pathlib, hashlib, json
p = pathlib.Path('.')
paths = sorted(set([q for pat in ('world/*.ail','go.mod','go.sum','design_docs/world-mission.md') for q in p.glob(pat) if q.is_file()] + [q for q in p.joinpath('tools/launchd').rglob('*') if q.is_file()]))
print(json.dumps([{'path':str(q),'sha256':hashlib.sha256(q.read_bytes()).hexdigest(),'lines':q.read_bytes().count(b'\n')} for q in paths],sort_keys=True))
PY
cmp "$OBS_BEFORE" "$OBS_AFTER"
git diff --exit-code -- world go.mod go.sum tools/launchd design_docs/world-mission.md
```

| Protected path | Before SHA256 = after SHA256 | Before lines = after lines |
|---|---|---|
| `design_docs/world-mission.md` | `5977231a7a5eca3c8de47e2fe4cf00c540dbd87a540a38fcbc9d78d725b05a47` = `5977231a7a5eca3c8de47e2fe4cf00c540dbd87a540a38fcbc9d78d725b05a47` | 6169 = 6169 |
| `go.mod` | `b3427f9c52e6240ada240666f6bb70ab6dcb4df22c4f7383d44cb4dfa46c9653` = `b3427f9c52e6240ada240666f6bb70ab6dcb4df22c4f7383d44cb4dfa46c9653` | 20 = 20 |
| `go.sum` | `0efaabe8522a52a6603eb11b56716b9518b8ce1e093ea2a8a8d42fc729ca4a26` = `0efaabe8522a52a6603eb11b56716b9518b8ce1e093ea2a8a8d42fc729ca4a26` | 57 = 57 |
| `tools/launchd/README.md` | `aeb663f2e95b4f1e168ad7a1bb687fa81ea516b4dcad78863bed7fea1b304793` = `aeb663f2e95b4f1e168ad7a1bb687fa81ea516b4dcad78863bed7fea1b304793` | 48 = 48 |
| `tools/launchd/derive-planner-lane.sh` | `a7c52f4852026dccd7e26d5cb058957ca7d8496073d43a87c9c1a4bc91837451` = `a7c52f4852026dccd7e26d5cb058957ca7d8496073d43a87c9c1a4bc91837451` | 109 = 109 |
| `tools/launchd/mission-template.plist` | `aaf57d459ae5917e8acdf69cee287586189469090da830dbba27729d09cd7dc0` = `aaf57d459ae5917e8acdf69cee287586189469090da830dbba27729d09cd7dc0` | 33 = 33 |
| `tools/launchd/test_mission_stall.sh` | `5e8c04072015274e54486f2d2df8d638d1e1c5f46fd55e07e89dc30b68daae33` = `5e8c04072015274e54486f2d2df8d638d1e1c5f46fd55e07e89dc30b68daae33` | 221 = 221 |
| `tools/launchd/testdata/planner-lane/n-no-backtick-bullet.md` | `833205678140b5b47079994eabd145e1c414e2d8a890a280e54b35fc17876006` = `833205678140b5b47079994eabd145e1c414e2d8a890a280e54b35fc17876006` | 12 = 12 |
| `world/contracts.ail` | `157574f614ec0d57b2c5302534984a84dfbcbaa6591eea560b4275932ab5f071` = `157574f614ec0d57b2c5302534984a84dfbcbaa6591eea560b4275932ab5f071` | 168 = 168 |
| `world/logepoch.ail` | `477ed4ce7ee223d032b12ebe998331de38fbe638a09b660a578f296d6ef611b5` = `477ed4ce7ee223d032b12ebe998331de38fbe638a09b660a578f296d6ef611b5` | 139 = 139 |
| `world/transitions.ail` | `3bc711b44ce0e06dc2d7630502968db8609c64d31982a07d3236e1eb79dead36` = `3bc711b44ce0e06dc2d7630502968db8609c64d31982a07d3236e1eb79dead36` | 104 = 104 |
| `world/types.ail` | `323cb6d40317e3100d8367b7904a772069798ea9a282fa08064a668f57efed6f` = `323cb6d40317e3100d8367b7904a772069798ea9a282fa08064a668f57efed6f` | 283 = 283 |

Each mutant is a temporary compiling guard/data edit and must reach its specified assertion. Mutate one at a time, save/restore only allowed files using scratch copies/cmp (no git), rerun original test green, record mutation/assertion/rc/restore result. No build-red, panic, watchdog expiration, skip, earlier guard failure or empty test selection counts as kill. The test source itself is not weakened to make a mutation succeed. No protected original mutation/restore.

## Ordered milestones

| ID | Title | Estimated minutes |
|---|---|---|
| M1a | Package evidence law and non-vacuous pinned manifest | 22 |
| M1b | Convention-v2 approval plan, descriptor and package docs | 25 |
| M2a | Demo CLI, preflight and authenticated domain commits | 25 |
| M2b | Real read and tool test/check invocation stages | 28 |
| M3a | HTTP proof adapters and real producer success | 28 |
| M3b | Proof/key refusals and seed S4 integration | 24 |
| M4a | Real approval transition, packet anchoring and budget delta | 28 |
| M4b | Complete why validation and public manifest assembly | 26 |
| M5a | Bounded anchored roster, ATTESTED and authored CLAIMED | 28 |
| M5b | TESTED/PROVEN resolution and truthful roster rendering | 29 |
| M5c | Seed roster validation and success manifest publication | 23 |
| M6 | Still serve flag and stable complete server HTML | 25 |
| M7a | Bounded supervisor with mandatory fake Chrome companions | 28 |
| M7b | PNG structure, decoded content and exact two-pass verifier | 28 |
| M7c | Capture orchestration, frozen inventories and Chrome drill | 28 |
| M8a | Reproduction docs, captions and release manifest schema | 25 |
| M8b | Final repository gates and protected-path observation | 25 |

## M1a — Package evidence law and non-vacuous pinned manifest

Design M1; 22 minutes; dependency: base; status **completed** — landed iteration 266 in `28ff4db371056929b9ddfa0c1dd9e87f95c8cfe6` (PASS 88/100, zero blocking; commit evidence). Iteration 267 independently reran both named evidence/legacy-drift tests under v0.41.0: PASS. At initial M1a metadata reconciliation, other milestones were pending; M1b completion is recorded below.

Exact executor files:

- `packages/world-demo/evidence.ail`
- `packages/world-demo/README.md`
- `host/coordinator/world_demo_test.go`

Define demo evidence-kind predicates in the package. Pin the named function/test manifest independently in Go test code, which runs in go test ./...; verify_ail ROOTS does not sweep world-demo, so its green result alone cannot discharge C1. Load ailang prompt/MCP language reference and inspect builtins before any .ail authoring. Standalone publisher modules may import only std/*; the package is a domain source directory, matching se-tools, not an invented registry release.

Acceptance criteria and named assertions:

- **C1 · `TestWorldDemoEvidencePinnedManifest`:** Require check.passed=true; named Z3 results legacyGradeCode, authoredClaimGradeCode, rosterKindAllowed all verified, errors/counterexamples/skips=0. Require named inline tests for six legacy constructors, authored-claim refusal/control and all four roster kinds plus unknown; fail empty tests/results. Scalar/string predicates carry requires/ensures first; no ADT generator skip can satisfy this gate.
- **C1/C5 · `TestWorldDemoLegacyGradeDrift`:** Execute package predicates and world/types.gradeOf on actual inputs under the pin (scratch runtime probe, checked before run). Compare all six legacy mappings, raw ProofReceipt=CLAIMED, authored demo claim=CLAIMED. No new world constructor, no AiReview impersonation, no mapping can produce PROVEN.

Exact focused and repository gate commands (environment/helper above applies):

```sh
go test -run '^TestWorldDemoEvidencePinnedManifest$' ./host/coordinator -count=1 -timeout 120s -v
go test -run '^TestWorldDemoLegacyGradeDrift$' ./host/coordinator -count=1 -timeout 120s -v
go test -run '^TestDecisionPacketCodec$' ./host/coordinator -count=1 -timeout 45s -v
(cd packages/world-demo && run_ail "$HOME/.pinned-ailang/ailang" check evidence.ail)
(cd packages/world-demo && run_ail "$HOME/.pinned-ailang/ailang" test --format json evidence.ail)
(cd packages/world-demo && run_ail "$HOME/.pinned-ailang/ailang" ai-check -timeout 5s evidence.ail)
AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_ail.sh
```

Mutations; each must reach this assertion:

| Mutant | Valid edit | Test / assertion that must red |
|---|---|---|
| M1a-law | Change the package legacy RecordedEffect mapping from ATTESTED to CLAIMED; keep valid return type. | `TestWorldDemoLegacyGradeDrift`: RecordedEffect grade differs from executable kernel law |
| M1a-empty | In the test gate parser remove the explicit empty-results rejection; its null-control uses empty results, not a missing executable. | `TestWorldDemoEvidencePinnedManifest`: null report must be refused |

Exact mutation commands:

```sh
go test -run '^TestWorldDemoLegacyGradeDrift$' ./host/coordinator -count=1 -timeout 120s -v
go test -run '^TestWorldDemoEvidencePinnedManifest$' ./host/coordinator -count=1 -timeout 120s -v
```

Base focused gates for this design milestone:

- **PASS** `go test -run '^TestDecisionPacketCodec$' ./host/coordinator -count=1 -timeout 45s -v` (1.77s).
- **PASS** `go test -run '^TestProducerMintsNoPacketForMultiEffectPlan$' ./host/coordinator -count=1 -timeout 45s -v` (1.54s).

New named tests/source modules are future implementation, absent at base; no base PASS is claimed for them.

### Executor directive notes

Narrowest gate list: exactly the command block above, plus each listed mutation command and before/after C11 observation. May touch only the exact files listed for M1a and its scratch snapshots. Repository gates: AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_ail.sh. 

**no git write operations; snapshot changed files into $HOME/.ailang/state/world-iter<N>/snap/M1a/ after each milestone**. Copy changed allowlisted files with their relative paths; include boundary command rc/wall-time, named PASS evidence, mutation ledger and observation records there. Never snapshot keys/session/output credentials into tracked assets. Maximum30min; incomplete is reported, never relabelled PASS.

## M1b — Convention-v2 approval plan, descriptor and package docs

Design M1; 25 minutes; dependency: M1a; status **completed** — iteration267 merge52ae7ff912fa65fcbedeb6f1b7c64e592b4b894b, PR #266; REQUIRED fresh-context evaluator PASS90/100, zero blocking; exact PR/merge CI2/2GREEN. Supporting output-schema gap carried as charter177. At that iteration267 boundary, later milestones were pending; current statuses are recorded below.

Exact executor files:

- `packages/world-demo/approve.ail`
- `packages/world-demo/transitions.json`
- `packages/world-demo/README.md`
- `host/coordinator/world_demo_test.go`

Publishable self-contained approval module; never import sibling evidence.ail into the capsule source (publisher stages one module). If law must be repeated, extend executable drift test to compare the two sources. Ship descriptor payload construction and plan/finish conventions with README.

Acceptance criteria and named assertions:

- **C1 · `TestWorldDemoApprovalExactOnePlan`:** Run real pinned approve module with convention-v2 input. Decode plan; exactly one effect e1 Human.Approve@worktree cost=1 requester=demo, finish=false/result=null; descriptor access cost=0 and declared cost=1. Unknown keys/phases and malformed roster refuse with zero effects. Require named oneApprovalEffect Z3 identity and named happy/refusal inline cases.

Exact focused and repository gate commands (environment/helper above applies):

```sh
go test -run '^TestWorldDemoApprovalExactOnePlan$' ./host/coordinator -count=1 -timeout 120s -v
go test -run '^TestProducerMintsNoPacketForMultiEffectPlan$' ./host/coordinator -count=1 -timeout 45s -v
(cd packages/world-demo && run_ail "$HOME/.pinned-ailang/ailang" check approve.ail)
(cd packages/world-demo && run_ail "$HOME/.pinned-ailang/ailang" test --format json approve.ail)
(cd packages/world-demo && run_ail "$HOME/.pinned-ailang/ailang" ai-check -timeout 5s approve.ail)
go test -run '^TestWorldDemoLegacyGradeDrift$' ./host/coordinator -count=1 -timeout 120s -v
AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_ail.sh
```

Mutations; each must reach this assertion:

| Mutant | Valid edit | Test / assertion that must red |
|---|---|---|
| M1b-two-effects | Append a second valid effect e2 to the approve module plan array while retaining all types and JSON validity. | `TestWorldDemoApprovalExactOnePlan`: len(plan.effects)==1 |

Exact mutation commands:

```sh
go test -run '^TestWorldDemoApprovalExactOnePlan$' ./host/coordinator -count=1 -timeout 120s -v
```

Base focused gates for this design milestone:

- **PASS** `go test -run '^TestDecisionPacketCodec$' ./host/coordinator -count=1 -timeout 45s -v` (1.77s).
- **PASS** `go test -run '^TestProducerMintsNoPacketForMultiEffectPlan$' ./host/coordinator -count=1 -timeout 45s -v` (1.54s).

New named tests/source modules are future implementation, absent at base; no base PASS is claimed for them.

### Executor directive notes

Narrowest gate list: exactly the command block above, plus each listed mutation command and before/after C11 observation. May touch only the exact files listed for M1b and its scratch snapshots. Repository gates: AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_ail.sh. 

**no git write operations; snapshot changed files into $HOME/.ailang/state/world-iter<N>/snap/M1b/ after each milestone**. Copy changed allowlisted files with their relative paths; include boundary command rc/wall-time, named PASS evidence, mutation ledger and observation records there. Never snapshot keys/session/output credentials into tracked assets. Maximum30min; incomplete is reported, never relabelled PASS.

## M2a — Demo CLI, preflight and authenticated domain commits

Design M2; 25 minutes; dependency: M1b; status **completed** — iteration268 merge 420c07e7c1c70eca6043a0c4e12a20f2adcaf972, PR #267; REQUIRED independent fresh-context evaluator round3 PASS90/100, zero blocking; exact PR/merge CI2/2 GREEN. Three nonblocking supporting-hunk gaps carried as charter179; later milestones pending.

Exact executor files:

- `cmd/ailang-worldd/demo.go`
- `cmd/ailang-worldd/demo_seed.go`
- `cmd/ailang-worldd/demo_test.go`
- `cmd/ailang-worldd/main.go`
- `docs/QUICKSTART.md`

Add dispatch/help and the preflight/REST stage only; intermediate CLI must explicitly refuse incomplete S2-S6, with no success manifest. No stub claims a complete seed. Use normal public commit headers and sole-writer HTTP API. Keep daemon configuration operator-owned; seed --db first checks daemon health identity and never opens another writable store. **Measured M2a clarification (iteration 268):** before any domain commit, the CLI may reuse the existing bounded read-only session path (`store.OpenReadOnly`, `authority.New(ro).ResolveContext`, existing broker grant predicates and `EffectSpend`) against that matched DB to validate all required grants/expiry/budget. The same read-only preflight may query existing `ScanUnreadableWorlds(ctx,"",1).Scanned` to refuse any stored World row, including an unselected World with no public head/log; this is not an all-object inventory claim. Close the read-only handle before commits. The REST commit endpoint resolves a live session but does not expose or enforce these grant amounts, while the session file contains only an opaque token; authenticated tools/list cannot prove the world.apply grant. This clarification preserves the design's sole-writer prohibition and missing/expired-grant acceptance, adds no endpoint or Mint path, and is a named evaluator target. Require a missing-world.apply arm with the other grants valid and no public-inventory change. Pin hash/version validation reuses the daemon-startup archive manifest through existing read-only archive.ReadManifest, matching its Hash/Size/Version to local pin bytes and health; no new subprocess or registry environment inheritance.

Acceptance criteria and named assertions:

- **C2 · `TestDemoSeedAuthenticatedCommits`:** Provision test session only through V11 testMintEnv/authority.Mint test setup before sole-writer daemon starts. Execute real demo commit stage with valid session; assert goal, two candidates, unchanged production logepoch source and explicitly demo-authored commentary have exact hashes, real carrying entries and CAS; all requests carry bearer, no Mint from seed.
- **C2 · `TestDemoSeedPreflightRefusals`:** Existing manifest, nonempty public world/log, absent/expired grants, wrong db/health interpreter, pin mismatch and unexpected CAS head must refuse before commits and leave no manifest; bootstrap registry and session rows alone are allowed. Unminted CLI refuses without opening /dev/tty or printing credentials.
- **C2/C10 · `TestDemoSeedHelpPreparation`:** demo seed --help accepts/documents all designed flags, reports separate pins, and QUICKSTART describes a disposable worktree with world/logepoch.ail plus attended descriptor publish then typed-y mint; no hidden git operation.

Exact focused and repository gate commands (environment/helper above applies):

```sh
go test -run '^TestDemoSeedAuthenticatedCommits$' ./cmd/ailang-worldd -count=1 -timeout 120s -v
go test -run '^TestDemoSeedPreflightRefusals$' ./cmd/ailang-worldd -count=1 -timeout 120s -v
go test -run '^TestDemoSeedHelpPreparation$' ./cmd/ailang-worldd -count=1 -timeout 120s -v
go test -run '^TestSessionMint_RequiresControllingTerminal$' ./cmd/ailang-worldd -count=1 -timeout 45s -v
go test -run '^TestRESTCommitRefusesDecisionPacketObjects$' ./host/daemon -count=1 -timeout 45s -v
```

Mutations; each must reach this assertion:

| Mutant | Valid edit | Test / assertion that must red |
|---|---|---|
| M2a-bearer | Omit Authorization header from a valid production demo REST commit request. | `TestDemoSeedAuthenticatedCommits`: successful seed commit stage returns success; denial detail is SessionAbsent |
| M2a-nonempty | Neuter the empty-public-world/log guard only. Use a real seeded disposable store and a fresh absent output filename; other preconditions valid. | `TestDemoSeedPreflightRefusals`: nonempty store must refuse before any public inventory change |

Exact mutation commands:

```sh
go test -run '^TestDemoSeedAuthenticatedCommits$' ./cmd/ailang-worldd -count=1 -timeout 120s -v
go test -run '^TestDemoSeedPreflightRefusals$' ./cmd/ailang-worldd -count=1 -timeout 120s -v
```

Base focused gates for this design milestone:

- **PASS** `go test -run '^TestSessionMint_RequiresControllingTerminal$' ./cmd/ailang-worldd -count=1 -timeout 45s -v` (2.06s).
- **PASS** `go test -run '^TestSeToolsE2EBinaryVersions$' ./host/daemon -count=1 -timeout 45s -v` (1.54s).
- **PASS** `go test -run '^TestSeToolsA2AEndToEnd$' ./host/daemon -count=1 -timeout 120s -v` (6.91s).
- **PASS** `go test -run '^TestCLIAgainstSeToolsDaemon$' ./cmd/ailang-worldd -count=1 -timeout 120s -v` (6.31s).

New named tests/source modules are future implementation, absent at base; no base PASS is claimed for them.

### Executor directive notes

Narrowest gate list: exactly the command block above, plus each listed mutation command and before/after C11 observation. May touch only the exact files listed for M2a and its scratch snapshots. Repository gates: focused affected-package gates now; full build/test, verify_ail and host races at M8b (base full test remains red; never skip its failures). 

**no git write operations; snapshot changed files into $HOME/.ailang/state/world-iter<N>/snap/M2a/ after each milestone**. Copy changed allowlisted files with their relative paths; include boundary command rc/wall-time, named PASS evidence, mutation ledger and observation records there. Never snapshot keys/session/output credentials into tracked assets. Maximum30min; incomplete is reported, never relabelled PASS.

## M2b — Real read and tool test/check invocation stages

Design M2; 28 minutes; dependency: M2a; status **completed** (iteration 271, `07d1d9f`).

Exact executor files:

- `cmd/ailang-worldd/demo_seed.go`
- `cmd/ailang-worldd/demo_test.go`
- `docs/QUICKSTART.md`

Resume stage orchestration through S3, still refuse CLI completion until later stages exist. Reuse cliSeBins pin discovery; provision real sessions using only V11 credentials. Do not create a release seed in the planner/executor. Positive test fixtures use the unmodified repository logepoch bytes.

Acceptance criteria and named assertions:

- **C2 · `TestDemoSeedReadTestCheck`:** Both real pins required. Publish existing se-tools descriptors through production publisher core in test setup; real /a2a/ dispatch for demo-read/test/check. Read bytes equal stored logepoch; test argv op=test flags json:"", named passes>0 and no failures/skips/vacuity; ai_check actual non-null success. Retain actual pin, argv, stdout/stderr, invocation/effect refs; budget deltas Read=1 CLI=2 and no extra calls. Prove CLI result is v0.52.1, not proof pin.
- **C2 · `TestDemoSeedRejectsToolReportRefusals`:** Use actual malformed/non-test/no-generator input files in disposable episode, and explicit null-report unit control; all must refuse, never publish a successful manifest or use fake phase output as positive evidence.

Exact focused and repository gate commands (environment/helper above applies):

```sh
go test -run '^TestDemoSeedReadTestCheck$' ./cmd/ailang-worldd -count=1 -timeout 120s -v
go test -run '^TestDemoSeedRejectsToolReportRefusals$' ./cmd/ailang-worldd -count=1 -timeout 120s -v
go test -run '^TestSeToolsE2EBinaryVersions$' ./host/daemon -count=1 -timeout 45s -v
go test -run '^TestSeToolsA2AEndToEnd$' ./host/daemon -count=1 -timeout 120s -v
go test -run '^TestCLIAgainstSeToolsDaemon$' ./cmd/ailang-worldd -count=1 -timeout 120s -v
```

Mutations; each must reach this assertion:

| Mutant | Valid edit | Test / assertion that must red |
|---|---|---|
| M2b-null | Neuter non-null tool-report guard in the report-validation helper; pass the same syntactically valid null report unit control. | `TestDemoSeedRejectsToolReportRefusals`: null report is refused |
| M2b-cli-op | Change only S3 test invocation op to ai_check (valid tool call). | `TestDemoSeedReadTestCheck`: retained TESTED invocation argv[op] is test |

Exact mutation commands:

```sh
go test -run '^TestDemoSeedRejectsToolReportRefusals$' ./cmd/ailang-worldd -count=1 -timeout 120s -v
go test -run '^TestDemoSeedReadTestCheck$' ./cmd/ailang-worldd -count=1 -timeout 120s -v
```

Base focused gates for this design milestone:

- **PASS** `go test -run '^TestSessionMint_RequiresControllingTerminal$' ./cmd/ailang-worldd -count=1 -timeout 45s -v` (2.06s).
- **PASS** `go test -run '^TestSeToolsE2EBinaryVersions$' ./host/daemon -count=1 -timeout 45s -v` (1.54s).
- **PASS** `go test -run '^TestSeToolsA2AEndToEnd$' ./host/daemon -count=1 -timeout 120s -v` (6.91s).
- **PASS** `go test -run '^TestCLIAgainstSeToolsDaemon$' ./cmd/ailang-worldd -count=1 -timeout 120s -v` (6.31s).

New named tests/source modules are future implementation, absent at base; no base PASS is claimed for them.

### Executor directive notes

Narrowest gate list: exactly the command block above, plus each listed mutation command and before/after C11 observation. May touch only the exact files listed for M2b and its scratch snapshots. Repository gates: focused affected-package gates now; full build/test, verify_ail and host races at M8b (base full test remains red; never skip its failures). TestDemoSeedReadTestCheck needs WORLD_TOOL_AILANG_BIN v0.52.1. Ordinary missing-pin run may SKIP; required milestone acceptance preflight must reject unset pin and require top-level PASS. Bad configured version is FAIL, never skip.

**no git write operations; snapshot changed files into $HOME/.ailang/state/world-iter<N>/snap/M2b/ after each milestone**. Copy changed allowlisted files with their relative paths; include boundary command rc/wall-time, named PASS evidence, mutation ledger and observation records there. Never snapshot keys/session/output credentials into tracked assets. Maximum30min; incomplete is reported, never relabelled PASS.

## M3a — HTTP proof adapters and real producer success

Design M3; 28 minutes; dependency: M2b; status **pending**.

Exact executor files:

- `cmd/ailang-worldd/demo_proof.go`
- `cmd/ailang-worldd/demo_proof_test.go`

First production NewProducer caller, unexported HTTP ObjectReader/ObjectWriter adapters. Pin checker from health archived interpreter and hash actual executable bytes. Preserve store.UncertainError across HTTP uncertain-outcome translation; never retry unknown commit. Tests only may wrap adapters; no new production injection seam.

Acceptance criteria and named assertions:

- **C3 · `TestDemoProofHTTPProducerRoundTrip`:** Real v0.41.0 producer on committed production world/logepoch.ail, new private random test key. GET payload=true validates source hash and cap; writer POST commit bearer/CAS preserves supplied envelope bytes/hash/interface/semantic ID. ValidateProof then Resolve yields PROVEN only for exact subject; REST carrying entry inventory includes envelope. Assert no competing writable store. Reader/writer BusyTimeout=0, HTTP 10s < read/write12s, exec30s/max1MiB/whole90s.
- **C3 · `TestDemoProofSourceHashRefusal`:** Forwarding test transport reads real stored source, flips one returned payload byte; valid HTTP response reaches producer hash guard. Assert hash-mismatch detail, zero ref and no envelope writer call.

Exact focused and repository gate commands (environment/helper above applies):

```sh
go test -run '^TestDemoProofHTTPProducerRoundTrip$' ./cmd/ailang-worldd -count=1 -timeout 120s -v
go test -run '^TestDemoProofSourceHashRefusal$' ./cmd/ailang-worldd -count=1 -timeout 120s -v
go test -run '^TestProducerEndToEndMintsEnvelopeValidatorAccepts$' ./host/evidence -count=1 -timeout 120s -v
go test -run '^TestWorkbenchGradeProvenFromSubjectBoundEnvelope$' ./host/daemon -count=1 -timeout 45s -v
```

Mutations; each must reach this assertion:

| Mutant | Valid edit | Test / assertion that must red |
|---|---|---|
| M3a-source | Flip one response source byte in the forwarding real-HTTP read arm; no malformed JSON. | `TestDemoProofHTTPProducerRoundTrip`: positive producer result must mint; producer reports source bytes hash mismatch before writer |
| M3a-no-rest | For the test-only forwarding writer, store supplied object through real store.PutObject in the test fixture instead of its HTTP commit; return nil. Production still has no store handle. | `TestDemoProofHTTPProducerRoundTrip`: envelope must have a REST carrying entry (object exists, commit missing) |

Exact mutation commands:

```sh
go test -run '^TestDemoProofHTTPProducerRoundTrip$' ./cmd/ailang-worldd -count=1 -timeout 120s -v
go test -run '^TestDemoProofHTTPProducerRoundTrip$' ./cmd/ailang-worldd -count=1 -timeout 120s -v
```

Base focused gates for this design milestone:

- **PASS** `go test -run '^TestProducerEndToEndMintsEnvelopeValidatorAccepts$' ./host/evidence -count=1 -timeout 120s -v` (2.44s).
- **PASS** `go test -run '^TestWorkbenchGradeProvenFromSubjectBoundEnvelope$' ./host/daemon -count=1 -timeout 45s -v` (1.36s).

New named tests/source modules are future implementation, absent at base; no base PASS is claimed for them.

### Executor directive notes

Narrowest gate list: exactly the command block above, plus each listed mutation command and before/after C11 observation. May touch only the exact files listed for M3a and its scratch snapshots. Repository gates: focused affected-package gates now; full build/test, verify_ail and host races at M8b (base full test remains red; never skip its failures). 

**no git write operations; snapshot changed files into $HOME/.ailang/state/world-iter<N>/snap/M3a/ after each milestone**. Copy changed allowlisted files with their relative paths; include boundary command rc/wall-time, named PASS evidence, mutation ledger and observation records there. Never snapshot keys/session/output credentials into tracked assets. Maximum30min; incomplete is reported, never relabelled PASS.

## M3b — Proof/key refusals and seed S4 integration

Design M3; 24 minutes; dependency: M3a; status **pending**.

Exact executor files:

- `cmd/ailang-worldd/demo_proof.go`
- `cmd/ailang-worldd/demo_proof_test.go`
- `cmd/ailang-worldd/demo_seed.go`
- `cmd/ailang-worldd/demo_test.go`

Compose S4 with S1-S3; keep success manifest disabled until anchored approval/why/roster validation. Fresh random key t.TempDir 0600, never fixture key, stdout, manifest or tracked bytes. No edits to producer/validator/core codecs.

Acceptance criteria and named assertions:

- **C3 · `TestDemoProofWrongKeyCannotProve`:** After real successful HTTP mint, validator configured with independently generated different key returns UNAVAILABLE, never PROVEN. Happy arm uses same key and exact subject.
- **C3 · `TestDemoProofNoVerifiedSetRefuses`:** Read a real committed standalone no-contract scratch source under the real pin, not hand-authored ai-check JSON: producer returns ErrVacuousCheck, zero envelope ref, zero writer calls; successful production logepoch control requires sameRef and servesEntry.
- **C3 · `TestDemoProofUncertainCommitStops`:** Test transport forwards real commit then loses the acknowledgement; demo surfaces uncertain outcome and performs exactly one write, no retry/reseed/manifest success. Both safe error and lost-response after commit have assertions.

Exact focused and repository gate commands (environment/helper above applies):

```sh
go test -run '^TestDemoProofWrongKeyCannotProve$' ./cmd/ailang-worldd -count=1 -timeout 120s -v
go test -run '^TestDemoProofNoVerifiedSetRefuses$' ./cmd/ailang-worldd -count=1 -timeout 120s -v
go test -run '^TestDemoProofUncertainCommitStops$' ./cmd/ailang-worldd -count=1 -timeout 120s -v
go test -run '^TestProducerEndToEndMintsEnvelopeValidatorAccepts$' ./host/evidence -count=1 -timeout 120s -v
```

Mutations; each must reach this assertion:

| Mutant | Valid edit | Test / assertion that must red |
|---|---|---|
| M3b-key | Change the validator key in the real positive arm to a different valid 32-byte private key. | `TestDemoProofWrongKeyCannotProve`: same-key positive arm must Resolve PROVEN |
| M3b-null-verified | Change the positive source/subject in the disposable test store to a valid no-contract module (consistent hash); retain the real checker. | `TestDemoProofNoVerifiedSetRefuses`: positive proof-control must mint required identities; observed ErrVacuousCheck |
| M3b-retry | Retry the HTTP proof writer once after the test transport uncertain outcome, retaining the same valid object. | `TestDemoProofUncertainCommitStops`: write attempt count must equal one |

Exact mutation commands:

```sh
go test -run '^TestDemoProofWrongKeyCannotProve$' ./cmd/ailang-worldd -count=1 -timeout 120s -v
go test -run '^TestDemoProofNoVerifiedSetRefuses$' ./cmd/ailang-worldd -count=1 -timeout 120s -v
go test -run '^TestDemoProofUncertainCommitStops$' ./cmd/ailang-worldd -count=1 -timeout 120s -v
```

Base focused gates for this design milestone:

- **PASS** `go test -run '^TestProducerEndToEndMintsEnvelopeValidatorAccepts$' ./host/evidence -count=1 -timeout 120s -v` (2.44s).
- **PASS** `go test -run '^TestWorkbenchGradeProvenFromSubjectBoundEnvelope$' ./host/daemon -count=1 -timeout 45s -v` (1.36s).

New named tests/source modules are future implementation, absent at base; no base PASS is claimed for them.

### Executor directive notes

Narrowest gate list: exactly the command block above, plus each listed mutation command and before/after C11 observation. May touch only the exact files listed for M3b and its scratch snapshots. Repository gates: focused affected-package gates now; full build/test, verify_ail and host races at M8b (base full test remains red; never skip its failures). 

**no git write operations; snapshot changed files into $HOME/.ailang/state/world-iter<N>/snap/M3b/ after each milestone**. Copy changed allowlisted files with their relative paths; include boundary command rc/wall-time, named PASS evidence, mutation ledger and observation records there. Never snapshot keys/session/output credentials into tracked assets. Maximum30min; incomplete is reported, never relabelled PASS.

## M4a — Real approval transition, packet anchoring and budget delta

Design M4; 28 minutes; dependency: M3b; status **pending**.

Exact executor files:

- `cmd/ailang-worldd/demo_seed.go`
- `cmd/ailang-worldd/demo_test.go`

S5 joins preexisting producer/coordinator paths. Tests may inject existing coordinator Now solely as V11; use real approve module for positive arm, never humanPlanRunner. Do not serialize complete manifest yet.

Acceptance criteria and named assertions:

- **C4 · `TestDemoSeedApprovalPacket`:** Run checked-in demo-approve through real pinned capsule/coordinator and daemon HumanHandler. Exactly one S5 carrying entry/one packet, pending approval, effect e1/request/proposal/input/output linked by content hash; Human.Approve remaining budget decreases by exactly1; descriptor access cost0. No decision verb, duplicate invocation or fabricated packet. Four evidence refs come from preceding real stages.

Exact focused and repository gate commands (environment/helper above applies):

```sh
go test -run '^TestDemoSeedApprovalPacket$' ./cmd/ailang-worldd -count=1 -timeout 120s -v
go test -run '^TestDaemonInvocationMintsDecisionPacket$' ./host/daemon -count=1 -timeout 45s -v
go test -run '^TestRESTCommitRefusesDecisionPacketObjects$' ./host/daemon -count=1 -timeout 45s -v
go test -run '^TestDecisionPacketCodec$' ./host/coordinator -count=1 -timeout 45s -v
```

Mutations; each must reach this assertion:

| Mutant | Valid edit | Test / assertion that must red |
|---|---|---|
| M4a-extra-effect | Return a valid two-effect Human.Approve plan from a copied test module (distinct e1/e2), preserving schema and declared capability. | `TestDemoSeedApprovalPacket`: successful single-approval stage with exactly one effect/packet (coordinator refusal must be visible in failure diagnostic) |

Exact mutation commands:

```sh
go test -run '^TestDemoSeedApprovalPacket$' ./cmd/ailang-worldd -count=1 -timeout 120s -v
```

Base focused gates for this design milestone:

- **PASS** `go test -run '^TestDecisionPacketCodec$' ./host/coordinator -count=1 -timeout 45s -v` (1.77s).
- **PASS** `go test -run '^TestCLIAgainstSeToolsDaemon$' ./cmd/ailang-worldd -count=1 -timeout 120s -v` (6.31s).
- **PASS** `go test -run '^TestRESTCommitRefusesDecisionPacketObjects$' ./host/daemon -count=1 -timeout 45s -v` (1.26s).
- **PASS** `go test -run '^TestDaemonInvocationMintsDecisionPacket$' ./host/daemon -count=1 -timeout 45s -v` (0.68s).
- **PASS** `go test -run '^TestWhyTargetFormsResolveToOneEntry$' ./cmd/ailang-worldd -count=1 -timeout 45s -v` (0.69s).

New named tests/source modules are future implementation, absent at base; no base PASS is claimed for them.

### Executor directive notes

Narrowest gate list: exactly the command block above, plus each listed mutation command and before/after C11 observation. May touch only the exact files listed for M4a and its scratch snapshots. Repository gates: focused affected-package gates now; full build/test, verify_ail and host races at M8b (base full test remains red; never skip its failures). Real approve plan needs AILANG_BIN; preceding read/test/check stages also require WORLD_TOOL_AILANG_BIN. Required gate cannot accept helper SKIP.

**no git write operations; snapshot changed files into $HOME/.ailang/state/world-iter<N>/snap/M4a/ after each milestone**. Copy changed allowlisted files with their relative paths; include boundary command rc/wall-time, named PASS evidence, mutation ledger and observation records there. Never snapshot keys/session/output credentials into tracked assets. Maximum30min; incomplete is reported, never relabelled PASS.

## M4b — Complete why validation and public manifest assembly

Design M4; 26 minutes; dependency: M4a; status **pending**.

Exact executor files:

- `cmd/ailang-worldd/demo_seed.go`
- `cmd/ailang-worldd/demo_test.go`
- `cmd/ailang-worldd/why_test.go`

Reuse actual CLI why implementation, not a substitute provenance walker. Complete manifest assembly helper; success serialization stays fail-closed until M5c roster checks. Unit omission controls derive from real reports.

Acceptance criteria and named assertions:

- **C4 · `TestDemoSeedWhyNamedLinks`:** Why for each real read/test/check/approval entry: coordinator=true, ok=true, at least one effect and exact expected world/entry/record/input/plan/effect/output link names all true; zero links must refuse. Assert coordinator before packet count in replacement control, so plain REST cannot mask the intended assertion.
- **C4 · `TestDemoSeedManifestReferences`:** Assemble (do not publish success yet) schema/version, frozen head, ordered actual entries/objects, pins/source/goal/candidates, packet/request/proposal, four evidence refs, invocation/effects and why raw JSON hashes. Follow each ref/hash, reject missing/mismatched refs; no credential/key/absolute path/capture time in public bytes. Scratch assembly is not release acceptance.

Exact focused and repository gate commands (environment/helper above applies):

```sh
go test -run '^TestDemoSeedWhyNamedLinks$' ./cmd/ailang-worldd -count=1 -timeout 120s -v
go test -run '^TestDemoSeedManifestReferences$' ./cmd/ailang-worldd -count=1 -timeout 120s -v
go test -run '^TestWhyTargetFormsResolveToOneEntry$' ./cmd/ailang-worldd -count=1 -timeout 45s -v
go test -run '^TestCLIAgainstSeToolsDaemon$' ./cmd/ailang-worldd -count=1 -timeout 120s -v
```

Mutations; each must reach this assertion:

| Mutant | Valid edit | Test / assertion that must red |
|---|---|---|
| M4b-plain-rest | Replace S5 test stimulus with a valid ordinary REST demo object commit (not forbidden packet semantic ID), then call the actual why validator on that entry. | `TestDemoSeedWhyNamedLinks`: coordinator must be true (evaluate before packet count) |
| M4b-link | Remove only the effect link from a decoded real why report in the validator unit control, retaining coordinator/ok and other links. | `TestDemoSeedWhyNamedLinks`: expected effect link must exist and be true |
| M4b-ref | Replace one manifest public ref with another valid absent SHA256 ref. | `TestDemoSeedManifestReferences`: every declared manifest reference resolves to matching content |

Exact mutation commands:

```sh
go test -run '^TestDemoSeedWhyNamedLinks$' ./cmd/ailang-worldd -count=1 -timeout 120s -v
go test -run '^TestDemoSeedWhyNamedLinks$' ./cmd/ailang-worldd -count=1 -timeout 120s -v
go test -run '^TestDemoSeedManifestReferences$' ./cmd/ailang-worldd -count=1 -timeout 120s -v
```

Base focused gates for this design milestone:

- **PASS** `go test -run '^TestDecisionPacketCodec$' ./host/coordinator -count=1 -timeout 45s -v` (1.77s).
- **PASS** `go test -run '^TestCLIAgainstSeToolsDaemon$' ./cmd/ailang-worldd -count=1 -timeout 120s -v` (6.31s).
- **PASS** `go test -run '^TestRESTCommitRefusesDecisionPacketObjects$' ./host/daemon -count=1 -timeout 45s -v` (1.26s).
- **PASS** `go test -run '^TestDaemonInvocationMintsDecisionPacket$' ./host/daemon -count=1 -timeout 45s -v` (0.68s).
- **PASS** `go test -run '^TestWhyTargetFormsResolveToOneEntry$' ./cmd/ailang-worldd -count=1 -timeout 45s -v` (0.69s).

New named tests/source modules are future implementation, absent at base; no base PASS is claimed for them.

### Executor directive notes

Narrowest gate list: exactly the command block above, plus each listed mutation command and before/after C11 observation. May touch only the exact files listed for M4b and its scratch snapshots. Repository gates: focused affected-package gates now; full build/test, verify_ail and host races at M8b (base full test remains red; never skip its failures). 

**no git write operations; snapshot changed files into $HOME/.ailang/state/world-iter<N>/snap/M4b/ after each milestone**. Copy changed allowlisted files with their relative paths; include boundary command rc/wall-time, named PASS evidence, mutation ledger and observation records there. Never snapshot keys/session/output credentials into tracked assets. Maximum30min; incomplete is reported, never relabelled PASS.

## M5a — Bounded anchored roster, ATTESTED and authored CLAIMED

Design M5; 28 minutes; dependency: M4b; status **pending**.

Exact executor files:

- `host/daemon/workbench_packet_evidence.go`
- `host/daemon/workbench_packet_evidence_test.go`
- `host/daemon/workbench_packet.go`
- `host/coordinator/world_demo_test.go`

Implement transport/projection validation only; host performs bounded I/O, package owns kind law. Expose no new authority API. No changes to packet v1 bytes or object-grade resolution. Add hash/kind refusals before any label. Re-run executable package law drift.

Acceptance criteria and named assertions:

- **C5 · `TestWorkbenchPacketRosterAnchoredBounds`:** Follow validated packet carrying record.input only, verify content hash/semantic ID and schema. Exactly four unique kind/ref pairs and source subject sibling; unknown caller grade/extra keys/schema never decorate a packet. 64KiB input, 1MiB result/read,4KiB claim,1MiB decoded stdout,256 named tests, shared10s deadline; boundary controls then cap+1 refusal. Decode RecordedEffect with RecordConsistent and request/result hashes for ATTESTED. Authored claim kind resolves only matching demo claim object, producer/author=demo-authored, never AiReview. Unavailable on bad kind/hash/overflow; no arbitrary semantic-ID scanning.

Exact focused and repository gate commands (environment/helper above applies):

```sh
go test -run '^TestWorkbenchPacketRosterAnchoredBounds$' ./host/daemon -count=1 -timeout 120s -v
go test -run '^TestWorldDemoLegacyGradeDrift$' ./host/coordinator -count=1 -timeout 120s -v
go test -run '^TestWorkbenchPacketValidated$' ./host/daemon -count=1 -timeout 45s -v
go test -run '^TestWorkbenchPacketEvidenceOnlyThroughRowOneOhFive$' ./host/daemon -count=1 -timeout 45s -v
```

Mutations; each must reach this assertion:

| Mutant | Valid edit | Test / assertion that must red |
|---|---|---|
| M5a-record-hash | Flip one returned record payload byte in test read wrapper after genuine store read; keep ref/meta intact. | `TestWorkbenchPacketRosterAnchoredBounds`: bad record hash yields UNAVAILABLE |
| M5a-claim-kind | Point CLAIMED pair at a real object of another declared semantic kind; retain valid object hash. | `TestWorkbenchPacketRosterAnchoredBounds`: resolved kind must equal declared kind; otherwise UNAVAILABLE |
| M5a-limit | Neuter decoded input-size cap guard; feed a valid 64KiB+1 input with all refs resolving. | `TestWorkbenchPacketRosterAnchoredBounds`: oversize input yields UNAVAILABLE |

Exact mutation commands:

```sh
go test -run '^TestWorkbenchPacketRosterAnchoredBounds$' ./host/daemon -count=1 -timeout 120s -v
go test -run '^TestWorkbenchPacketRosterAnchoredBounds$' ./host/daemon -count=1 -timeout 120s -v
go test -run '^TestWorkbenchPacketRosterAnchoredBounds$' ./host/daemon -count=1 -timeout 120s -v
```

Base focused gates for this design milestone:

- **PASS** `go test -run '^TestWorkbenchPacketValidated$' ./host/daemon -count=1 -timeout 45s -v` (0.75s).
- **PASS** `go test -run '^TestWorkbenchPacketEvidenceOnlyThroughRowOneOhFive$' ./host/daemon -count=1 -timeout 45s -v` (1.23s).
- **PASS** `go test -run '^TestWorkbenchGradeProvenFromSubjectBoundEnvelope$' ./host/daemon -count=1 -timeout 45s -v` (1.36s).
- **PASS** `go test -run '^TestWorkbenchViewFieldsAllRender$' ./host/workbench -count=1 -timeout 45s -v` (0.34s).
- **PASS** `go test -run '^TestWorkbenchNewViewTypesInCensus$' ./host/workbench -count=1 -timeout 45s -v` (0.32s).

New named tests/source modules are future implementation, absent at base; no base PASS is claimed for them.

### Executor directive notes

Narrowest gate list: exactly the command block above, plus each listed mutation command and before/after C11 observation. May touch only the exact files listed for M5a and its scratch snapshots. Repository gates: focused affected-package gates now; full build/test, verify_ail and host races at M8b (base full test remains red; never skip its failures). 

**no git write operations; snapshot changed files into $HOME/.ailang/state/world-iter<N>/snap/M5a/ after each milestone**. Copy changed allowlisted files with their relative paths; include boundary command rc/wall-time, named PASS evidence, mutation ledger and observation records there. Never snapshot keys/session/output credentials into tracked assets. Maximum30min; incomplete is reported, never relabelled PASS.

## M5b — TESTED/PROVEN resolution and truthful roster rendering

Design M5; 29 minutes; dependency: M5a; status **pending**.

Exact executor files:

- `host/daemon/workbench_packet_evidence.go`
- `host/daemon/workbench_packet_evidence_test.go`
- `host/daemon/workbench_packet.go`
- `host/workbench/render.go`
- `host/workbench/render_test.go`
- `host/workbench/live_test.go`

Finish grades and renderer; no callers supply authoritative grade. Preserve row105 authority. Extend workbenchViewTypes census with roster structures in the same run; control StillMode must render through the template branch when added next.

Acceptance criteria and named assertions:

- **C5 · `TestWorkbenchPacketRosterTestedAndProven`:** Real TESTED S3 invocation: carrying commit, exactly one Ailang.CLI effect, actual request/result hashes, argv test, strict success JSON and nonempty named pass set, zero failed/skipped/vacuous; reject arbitrary prose. PROVEN follows only existing ValidateProof→Resolve with actual operator key/exact source, missing/wrong key/subject refuse. Demo claim presented as PROVEN input kind must not promote. Existing two subject-grade lines preserve prior values.
- **C5 · `TestRenderPacketEvidenceFourRows`:** Render all four grade labels, TESTED pass verdict, author/producer on every row, demo-authored caption, report/source/pin links and existing honest subject-grade lines; escape HTML. No raw grade constant confers authority. Add new view structs to existing field census; every exported field renders.

Exact focused and repository gate commands (environment/helper above applies):

```sh
go test -run '^TestWorkbenchPacketRosterTestedAndProven$' ./host/daemon -count=1 -timeout 120s -v
go test -run '^TestRenderPacketEvidenceFourRows$' ./host/workbench -count=1 -timeout 120s -v
go test -run '^TestWorkbenchGradeProvenFromSubjectBoundEnvelope$' ./host/daemon -count=1 -timeout 45s -v
go test -run '^TestWorkbenchPacketEvidenceOnlyThroughRowOneOhFive$' ./host/daemon -count=1 -timeout 45s -v
go test -run '^TestWorkbenchViewFieldsAllRender$' ./host/workbench -count=1 -timeout 45s -v
go test -run '^TestWorkbenchNewViewTypesInCensus$' ./host/workbench -count=1 -timeout 45s -v
go test -race ./host/workbench -count=1 -timeout 120s
```

Mutations; each must reach this assertion:

| Mutant | Valid edit | Test / assertion that must red |
|---|---|---|
| M5b-promote-claim | Change only authored claim requested kind to PROVEN in a valid roster control. | `TestWorkbenchPacketRosterTestedAndProven`: claim cannot promote; row remains UNAVAILABLE |
| M5b-test-prose | Substitute arbitrary UTF-8 prose as stdout in a real test result via forwarding reader control; rebuild surrounding result hash so JSON validation is the reached guard. | `TestWorkbenchPacketRosterTestedAndProven`: TESTED validator rejects non-test-report text |
| M5b-test-guard | Neuter failed/skipped/vacuous test-report guard on otherwise structurally valid test JSON control. | `TestWorkbenchPacketRosterTestedAndProven`: nonpassing test report yields UNAVAILABLE |

Exact mutation commands:

```sh
go test -run '^TestWorkbenchPacketRosterTestedAndProven$' ./host/daemon -count=1 -timeout 120s -v
go test -run '^TestWorkbenchPacketRosterTestedAndProven$' ./host/daemon -count=1 -timeout 120s -v
go test -run '^TestWorkbenchPacketRosterTestedAndProven$' ./host/daemon -count=1 -timeout 120s -v
```

Base focused gates for this design milestone:

- **PASS** `go test -run '^TestWorkbenchPacketValidated$' ./host/daemon -count=1 -timeout 45s -v` (0.75s).
- **PASS** `go test -run '^TestWorkbenchPacketEvidenceOnlyThroughRowOneOhFive$' ./host/daemon -count=1 -timeout 45s -v` (1.23s).
- **PASS** `go test -run '^TestWorkbenchGradeProvenFromSubjectBoundEnvelope$' ./host/daemon -count=1 -timeout 45s -v` (1.36s).
- **PASS** `go test -run '^TestWorkbenchViewFieldsAllRender$' ./host/workbench -count=1 -timeout 45s -v` (0.34s).
- **PASS** `go test -run '^TestWorkbenchNewViewTypesInCensus$' ./host/workbench -count=1 -timeout 45s -v` (0.32s).

- **PASS** `go test -race ./host/workbench -count=1 -timeout 120s` (1.99s).

New named tests/source modules are future implementation, absent at base; no base PASS is claimed for them.

### Executor directive notes

Narrowest gate list: exactly the command block above, plus each listed mutation command and before/after C11 observation. May touch only the exact files listed for M5b and its scratch snapshots. Repository gates: go test -race ./host/workbench -count=1 -timeout 120s. 

**no git write operations; snapshot changed files into $HOME/.ailang/state/world-iter<N>/snap/M5b/ after each milestone**. Copy changed allowlisted files with their relative paths; include boundary command rc/wall-time, named PASS evidence, mutation ledger and observation records there. Never snapshot keys/session/output credentials into tracked assets. Maximum30min; incomplete is reported, never relabelled PASS.

## M5c — Seed roster validation and success manifest publication

Design M5; 23 minutes; dependency: M5b; status **pending**.

Exact executor files:

- `cmd/ailang-worldd/demo_seed.go`
- `cmd/ailang-worldd/demo_test.go`
- `host/daemon/workbench_packet_evidence_test.go`

Enable CLI success now that the complete roster validator exists. Tests exercise real validation helper with actual stored packet/input and derived mismatches; no test-only bypass in release seed. Serialize public manifest atomically with absent-output guard; uncertain outcome never reseeds.

Acceptance criteria and named assertions:

- **C2/C3/C4/C5 · `TestDemoSeedEndToEnd`:** Real seeded temporary store from empty public state: S1-S6 all occur through REST/coordinator, real pins/source/private key, valid packet roster and all why links. Serialize success manifest once only after checks; no missing/mismatched evidence. No release credentials, fake positive phase runner or new production injection seam.
- **C5 · `TestDemoSeedRosterKindMismatchRefuses`:** Genuine real object substituted for CLAIMED ref with another semantic kind: workbench shows UNAVAILABLE and production seed final validator refuses with no success manifest; positive test reaches final roster-validation helper using same real seed data.

Exact focused and repository gate commands (environment/helper above applies):

```sh
go test -run '^TestDemoSeedEndToEnd$' ./cmd/ailang-worldd -count=1 -timeout 120s -v
go test -run '^TestDemoSeedRosterKindMismatchRefuses$' ./cmd/ailang-worldd -count=1 -timeout 120s -v
go test -run '^TestDemoSeedWhyNamedLinks$' ./cmd/ailang-worldd -count=1 -timeout 120s -v
go test -run '^TestWorkbenchPacketValidated$' ./host/daemon -count=1 -timeout 45s -v
```

Mutations; each must reach this assertion:

| Mutant | Valid edit | Test / assertion that must red |
|---|---|---|
| M5c-kind-finalizer | Neuter only final seed roster-kind refusal; leave real projection returning UNAVAILABLE. | `TestDemoSeedRosterKindMismatchRefuses`: seed must refuse kind mismatch and write no success manifest |

Exact mutation commands:

```sh
go test -run '^TestDemoSeedRosterKindMismatchRefuses$' ./cmd/ailang-worldd -count=1 -timeout 120s -v
```

Base focused gates for this design milestone:

- **PASS** `go test -run '^TestWorkbenchPacketValidated$' ./host/daemon -count=1 -timeout 45s -v` (0.75s).
- **PASS** `go test -run '^TestWorkbenchPacketEvidenceOnlyThroughRowOneOhFive$' ./host/daemon -count=1 -timeout 45s -v` (1.23s).
- **PASS** `go test -run '^TestWorkbenchGradeProvenFromSubjectBoundEnvelope$' ./host/daemon -count=1 -timeout 45s -v` (1.36s).
- **PASS** `go test -run '^TestWorkbenchViewFieldsAllRender$' ./host/workbench -count=1 -timeout 45s -v` (0.34s).
- **PASS** `go test -run '^TestWorkbenchNewViewTypesInCensus$' ./host/workbench -count=1 -timeout 45s -v` (0.32s).

New named tests/source modules are future implementation, absent at base; no base PASS is claimed for them.

### Executor directive notes

Narrowest gate list: exactly the command block above, plus each listed mutation command and before/after C11 observation. May touch only the exact files listed for M5c and its scratch snapshots. Repository gates: focused affected-package gates now; full build/test, verify_ail and host races at M8b (base full test remains red; never skip its failures). 

**no git write operations; snapshot changed files into $HOME/.ailang/state/world-iter<N>/snap/M5c/ after each milestone**. Copy changed allowlisted files with their relative paths; include boundary command rc/wall-time, named PASS evidence, mutation ledger and observation records there. Never snapshot keys/session/output credentials into tracked assets. Maximum30min; incomplete is reported, never relabelled PASS.

## M6 — Still serve flag and stable complete server HTML

Design M6; 25 minutes; dependency: M5c; status **pending**.

Exact executor files:

- `cmd/ailang-worldd/main.go`
- `cmd/ailang-worldd/main_test.go`
- `host/daemon/daemon.go`
- `host/daemon/workbench.go`
- `host/daemon/workbench_still_test.go`
- `host/workbench/render.go`
- `host/workbench/render_test.go`
- `docs/QUICKSTART.md`

D2 suppresses one script only; no CSS compacting/layout redesign to hide PD1. If seeded repeated GET varies report element/bytes and stop product acceptance. Do not broaden flag to freeze authority. No real Chrome needed for these tests.

Acceptance criteria and named assertions:

- **C6 · `TestWorkbenchStillAllRegions`:** Recorder-based real route: still mode includes timeline/live list/graph/decisions/packet refs/static footer, zero script tags. Repeat GET bytes on same store with request times different by2s equal; no nonce/time rewrite. Preserve all CSP/headers/read-only semantics.
- **C6 · `TestServeWorkbenchStillFlag`:** Run serve help/flag parser and renderer; explicit --workbench-still maps to Config/Page flag; default absent leaves exactly one live script. Invalid flag usage still refuses.
- **C6 · `TestWorkbenchStillDefaultLiveScript`:** Default configuration still emits /workbench/live.js; true suppresses it. Page flag is referenced in the template branch so field census remains non-vacuous.

Exact focused and repository gate commands (environment/helper above applies):

```sh
go test -run '^TestWorkbenchStillAllRegions$' ./host/daemon -count=1 -timeout 120s -v
go test -run '^TestServeWorkbenchStillFlag$' ./cmd/ailang-worldd -count=1 -timeout 120s -v
go test -run '^TestWorkbenchStillDefaultLiveScript$' ./host/workbench -count=1 -timeout 120s -v
go test -run '^TestWorkbenchGraphWiring$' ./host/daemon -count=1 -timeout 45s -v
go test -run '^TestWorkbenchDecisionsPane$' ./host/daemon -count=1 -timeout 45s -v
go test -run '^TestWorkbenchLiveCursorCrossesGap$' ./host/daemon -count=1 -timeout 45s -v
go test -run '^TestWorkbenchLiveScriptTag$' ./host/daemon -count=1 -timeout 45s -v
go test -run '^TestRenderQuietFooter$' ./host/workbench -count=1 -timeout 45s -v
go test -run '^TestWorkbenchViewFieldsAllRender$' ./host/workbench -count=1 -timeout 45s -v
```

Mutations; each must reach this assertion:

| Mutant | Valid edit | Test / assertion that must red |
|---|---|---|
| M6-script | Replace conditional script emission with unconditional valid template script tag. | `TestWorkbenchStillAllRegions`: still page has zero script tags |
| M6-graph | Remove graph section emission only in still template branch, preserving valid template syntax. | `TestWorkbenchStillAllRegions`: world graph required region must exist |
| M6-default | Change serve flag default from false to true. | `TestServeWorkbenchStillFlag`: default configuration must emit live script |

Exact mutation commands:

```sh
go test -run '^TestWorkbenchStillAllRegions$' ./host/daemon -count=1 -timeout 120s -v
go test -run '^TestWorkbenchStillAllRegions$' ./host/daemon -count=1 -timeout 120s -v
go test -run '^TestServeWorkbenchStillFlag$' ./cmd/ailang-worldd -count=1 -timeout 120s -v
```

Base focused gates for this design milestone:

- **PASS** `go test -run '^TestWorkbenchGraphWiring$' ./host/daemon -count=1 -timeout 45s -v` (0.94s).
- **PASS** `go test -run '^TestWorkbenchDecisionsPane$' ./host/daemon -count=1 -timeout 45s -v` (0.71s).
- **PASS** `go test -run '^TestWorkbenchLiveCursorCrossesGap$' ./host/daemon -count=1 -timeout 45s -v` (0.76s).
- **PASS** `go test -run '^TestWorkbenchLiveScriptTag$' ./host/daemon -count=1 -timeout 45s -v` (0.76s).
- **PASS** `go test -run '^TestRenderQuietFooter$' ./host/workbench -count=1 -timeout 45s -v` (0.67s).

New named tests/source modules are future implementation, absent at base; no base PASS is claimed for them.

### Executor directive notes

Narrowest gate list: exactly the command block above, plus each listed mutation command and before/after C11 observation. May touch only the exact files listed for M6 and its scratch snapshots. Repository gates: focused affected-package gates now; full build/test, verify_ail and host races at M8b (base full test remains red; never skip its failures). 

**no git write operations; snapshot changed files into $HOME/.ailang/state/world-iter<N>/snap/M6/ after each milestone**. Copy changed allowlisted files with their relative paths; include boundary command rc/wall-time, named PASS evidence, mutation ledger and observation records there. Never snapshot keys/session/output credentials into tracked assets. Maximum30min; incomplete is reported, never relabelled PASS.

## M7a — Bounded supervisor with mandatory fake Chrome companions

Design M7; 28 minutes; dependency: M6; status **pending**.

Exact executor files:

- `scripts/capture_marketing.sh`
- `scripts/marketing_capture.py`
- `scripts/marketingcheck/supervisor_test.go`

Python stdlib supervises only; Go tests drive real Python code with fake scripts. No image decoder, installation or Go dependency. Keep HOME intact, DEVNULL stdin and bounded group cleanup. Shell --help documents full config; capture must not mint/commit/approve.

Acceptance criteria and named assertions:

- **C7 · `TestMarketingSupervisorFakeChromeSuccess`:** Always-on fake Chrome shell script records complete argv, preserved HOME, stdin EOF, unique absolute profile/destination, scheme flag/viewport/scale/language. Writes fresh regular nonempty PNG and exits0; real Python supervisor accepts. Existing stale PNG cannot satisfy a new destination.
- **C7 · `TestMarketingSupervisorFakeChromeNoFile`:** Always-on fake Chrome exits0 but writes no file; supervisor returns typed MissingPNG. Controls also reject nonzero exit even with PNG, empty/nonregular file. Test does not need browser/socket.
- **C7 · `TestMarketingSupervisorFakeChromeTimeoutKillsGroup`:** Always-on fake Chrome starts child/grandchild, writes a PNG then hangs. Short test deadline via subprocess CLI timeout parameter (unit control only); Timeout verdict despite PNG, SIGKILL process group and reap≤5s; checked PIDs child+grandchild dead. Production deadline30s and entire pass≤300s.
- **C7 · `TestMarketingSupervisorFakeDaemonCleanup`:** Always-on fake daemon ignores TERM then valid KILL/reap;≤5s each production wait. Inject cleanup failure unit control must override earlier shot success. Shell normal/signal traps exercised, no unsupervised orphan.

Exact focused and repository gate commands (environment/helper above applies):

```sh
go test -run '^TestMarketingSupervisorFakeChromeSuccess$' ./scripts/marketingcheck -count=1 -timeout 90s -v
go test -run '^TestMarketingSupervisorFakeChromeNoFile$' ./scripts/marketingcheck -count=1 -timeout 90s -v
go test -run '^TestMarketingSupervisorFakeChromeTimeoutKillsGroup$' ./scripts/marketingcheck -count=1 -timeout 90s -v
go test -run '^TestMarketingSupervisorFakeDaemonCleanup$' ./scripts/marketingcheck -count=1 -timeout 90s -v
bash -n scripts/capture_marketing.sh
```

Mutations; each must reach this assertion:

| Mutant | Valid edit | Test / assertion that must red |
|---|---|---|
| M7a-no-file | Neuter regular-nonempty-new-file guard after rc=0 in Python supervisor. | `TestMarketingSupervisorFakeChromeNoFile`: rc0 with no file must be MissingPNG |
| M7a-child-only | Replace group kill with child-only kill; test cleanup always kills/checks its own retained grandchild to avoid leak. | `TestMarketingSupervisorFakeChromeTimeoutKillsGroup`: grandchild must be dead after bounded cleanup |
| M7a-timeout-success | On timeout with existing PNG return success after normal kill/reap. | `TestMarketingSupervisorFakeChromeTimeoutKillsGroup`: timeout must fail even when PNG exists |

Exact mutation commands:

```sh
go test -run '^TestMarketingSupervisorFakeChromeNoFile$' ./scripts/marketingcheck -count=1 -timeout 90s -v
go test -run '^TestMarketingSupervisorFakeChromeTimeoutKillsGroup$' ./scripts/marketingcheck -count=1 -timeout 90s -v
go test -run '^TestMarketingSupervisorFakeChromeTimeoutKillsGroup$' ./scripts/marketingcheck -count=1 -timeout 90s -v
```

Base focused gates for this design milestone:

- **PASS** `go test -run '^TestWorkbenchGraphWiring$' ./host/daemon -count=1 -timeout 45s -v` (0.94s).
- **PASS** `go test -run '^TestWorkbenchDecisionsPane$' ./host/daemon -count=1 -timeout 45s -v` (0.71s).
- **PASS** `go test -run '^TestWorkbenchLiveCursorCrossesGap$' ./host/daemon -count=1 -timeout 45s -v` (0.76s).
- **SKIP (not acceptance)** `go test -run '^TestWorkbenchLiveChromeDrill$' ./host/daemon -count=1 -timeout 45s -v` (0.66s).

New named tests/source modules are future implementation, absent at base; no base PASS is claimed for them.

### Executor directive notes

Narrowest gate list: exactly the command block above, plus each listed mutation command and before/after C11 observation. May touch only the exact files listed for M7a and its scratch snapshots. Repository gates: focused affected-package gates now; full build/test, verify_ail and host races at M8b (base full test remains red; never skip its failures). All four tests are binary-free companions, always run in CI. Any future real-browser supervisor probe requires WORLD_CHROME and these exact companions; fake success proves supervisor only, never visible product content.

**no git write operations; snapshot changed files into $HOME/.ailang/state/world-iter<N>/snap/M7a/ after each milestone**. Copy changed allowlisted files with their relative paths; include boundary command rc/wall-time, named PASS evidence, mutation ledger and observation records there. Never snapshot keys/session/output credentials into tracked assets. Maximum30min; incomplete is reported, never relabelled PASS.

## M7b — PNG structure, decoded content and exact two-pass verifier

Design M7; 28 minutes; dependency: M7a; status **pending**.

Exact executor files:

- `scripts/marketingcheck/main.go`
- `scripts/marketingcheck/png.go`
- `scripts/marketingcheck/png_test.go`

Separate structural parser/limit functions from image decoder and six-set gate for non-masked mutation controls. Per-file cap makes six-shot aggregate8MiB unreachable (six≤6MiB); test directory-total function independently with other PNGs. These are validity controls not candidate release assets. CRC fixture direct parser test required because image/png also rejects CRC. Verifier CLI check --pass-a --pass-b --manifest --rectangles; no image mutation.

Acceptance criteria and named assertions:

- **C8/C9 · `TestMarketingPNGStructure`:** Go stdlib image/png+crc32: six exact names/path rule/signature, first13-byte IHDR, valid CRCs, allowed chunks, one terminal IEND/no tail, reject text/exif, per-file1,048,576 and directory-total8,388,608 (counts ALL PNGs). Cap/CRC/unknown-chunk fixtures remain valid far enough to reach specific guard.
- **C8 · `TestMarketingPNGContent`:** Exact viewport dimensions, >32 RGB colors, chromatic fraction≥1%, variance in each measured content rectangle incl packet validation+four rows, per-shot light/dark background luminance. Reject uniform PNG and separately mostly-gray colorful test. Rectangles in bounds; decoded content gate cannot accept tiny arbitrary crop as visible text proof.
- **C8/C9 · `TestMarketingTwoPassInventory`:** Exactly six names each pass, SHA256 raw PNG bytes equal, hashes match manifest; no recompress/metadata strip/pixel similarity. Missing shot and altered valid PNG (different ancillary allowed sRGB/chunk encoding with rebuilt CRC or one actual pixel, then valid encode) refuse specifically equality. Equal-copy control passes.

Exact focused and repository gate commands (environment/helper above applies):

```sh
go test -run '^TestMarketingPNGStructure$' ./scripts/marketingcheck -count=1 -timeout 90s -v
go test -run '^TestMarketingPNGContent$' ./scripts/marketingcheck -count=1 -timeout 90s -v
go test -run '^TestMarketingTwoPassInventory$' ./scripts/marketingcheck -count=1 -timeout 90s -v
go test -race ./scripts/marketingcheck -count=1 -timeout 180s
```

Mutations; each must reach this assertion:

| Mutant | Valid edit | Test / assertion that must red |
|---|---|---|
| M7b-uniform | Neuter color/variance guard; valid exact-dimension uniform decoded PNG fixture reaches gate. | `TestMarketingPNGContent`: uniform PNG must be rejected |
| M7b-missing | Neuter exact six-name set guard; five valid matching pairs and manifest control. | `TestMarketingTwoPassInventory`: missing shot must be rejected |
| M7b-altered | Neuter SHA256 pair equality check; use two valid nonuniform dimension-correct files with different bytes. | `TestMarketingTwoPassInventory`: altered raw byte must be rejected |
| M7b-forbidden | Allow tEXt chunk explicitly; valid CRC/control PNG containing tEXt. | `TestMarketingPNGStructure`: forbidden tEXt must be rejected |
| M7b-crc | Neuter CRC comparison; call chunk structure validator directly before image/png so decoder redundancy cannot mask guard. | `TestMarketingPNGStructure`: bad CRC typed structure refusal |
| M7b-tail | Neuter trailing bytes guard after terminal IEND. | `TestMarketingPNGStructure`: one trailing byte must be rejected |
| M7b-cap | Raise per-file cap by1 on cap+1 valid chunk stream; CRCs/dimensions/content remain valid. | `TestMarketingPNGStructure`: cap+1 must be rejected |
| M7b-total | Neuter directory-total cap; use at least9 otherwise-allowed PNGs to reach standalone fleet-limit function before six-shot exact-set check. | `TestMarketingPNGStructure`: all-PNG aggregate >8MiB must be rejected |

Exact mutation commands:

```sh
go test -run '^TestMarketingPNGContent$' ./scripts/marketingcheck -count=1 -timeout 90s -v
go test -run '^TestMarketingTwoPassInventory$' ./scripts/marketingcheck -count=1 -timeout 90s -v
go test -run '^TestMarketingTwoPassInventory$' ./scripts/marketingcheck -count=1 -timeout 90s -v
go test -run '^TestMarketingPNGStructure$' ./scripts/marketingcheck -count=1 -timeout 90s -v
go test -run '^TestMarketingPNGStructure$' ./scripts/marketingcheck -count=1 -timeout 90s -v
go test -run '^TestMarketingPNGStructure$' ./scripts/marketingcheck -count=1 -timeout 90s -v
go test -run '^TestMarketingPNGStructure$' ./scripts/marketingcheck -count=1 -timeout 90s -v
go test -run '^TestMarketingPNGStructure$' ./scripts/marketingcheck -count=1 -timeout 90s -v
```

Base focused gates for this design milestone:

- **PASS** `go test -run '^TestWorkbenchGraphWiring$' ./host/daemon -count=1 -timeout 45s -v` (0.94s).
- **PASS** `go test -run '^TestWorkbenchDecisionsPane$' ./host/daemon -count=1 -timeout 45s -v` (0.71s).
- **PASS** `go test -run '^TestWorkbenchLiveCursorCrossesGap$' ./host/daemon -count=1 -timeout 45s -v` (0.76s).
- **SKIP (not acceptance)** `go test -run '^TestWorkbenchLiveChromeDrill$' ./host/daemon -count=1 -timeout 45s -v` (0.66s).

New named tests/source modules are future implementation, absent at base; no base PASS is claimed for them.

### Executor directive notes

Narrowest gate list: exactly the command block above, plus each listed mutation command and before/after C11 observation. May touch only the exact files listed for M7b and its scratch snapshots. Repository gates: go test -race ./scripts/marketingcheck -count=1 -timeout 180s. 

**no git write operations; snapshot changed files into $HOME/.ailang/state/world-iter<N>/snap/M7b/ after each milestone**. Copy changed allowlisted files with their relative paths; include boundary command rc/wall-time, named PASS evidence, mutation ledger and observation records there. Never snapshot keys/session/output credentials into tracked assets. Maximum30min; incomplete is reported, never relabelled PASS.

## M7c — Capture orchestration, frozen inventories and Chrome drill

Design M7; 28 minutes; dependency: M7b; status **pending**.

Exact executor files:

- `scripts/capture_marketing.sh`
- `scripts/marketing_capture.py`
- `scripts/marketingcheck/main.go`
- `scripts/marketingcheck/capture_test.go`
- `scripts/marketingcheck/observation_test.go`

Full capture flags include --db --manifest --proof-key-file --ailang-bin --workspace-root --tool-ailang-bin --chrome --out --verify-twice --rectangles.≤10s HTTP/≤30s startup, two≤300s passes, progress between shots. Scratch output only. If raw equality fails outside scheme/script causes report varying element+both hashes as M7 residual; criterion remains unmet, never pixel-similarity success.

Acceptance criteria and named assertions:

- **C7/C8/C9 · `TestMarketingCaptureFakeChromeSixShots`:** Always-on binary-free companion exercises full shell/Python/Go chain with fake Chrome plus fixture daemon responses (refusal/unit control, not release data). Both six-name passes, exact argv routes/viewport/scheme, HTML required refs/no scripts, bounded startup/cleanup and manifest matching; no commit/mint/decision/background call.
- **C8/C11 · `TestMarketingCaptureInventoryMutation`:** Read actual head/entries/objects before/after test-store capture. Inject genuine extra normal commit into disposable test store; frozen inventory equality refuses. Baseline no-write pass equal; registries/session diagnostics excluded from public inventory only if explicitly enumerated.
- **C11 · `TestMarketingProtectedObservationComparator`:** Observe protected files by SHA256 AND newline count before/after. Mutate only a t.TempDir copied input, prove digest comparator reds even for same line count and line comparator reds independently using controlled observations. Never write/restore protected source.
- **C7/C8/C9 · `TestMarketingChromeSixShots`:** WORLD_CHROME explicit gate; require executable Chrome155 version, both pins, measured rectangle config WORLD_MARKETING_RECTS. Real demo seed only test-minted through V11, frozen store, all six real shots in two passes, raw GET required text plus decoded pixels and exact hashes. Skip only if WORLD_CHROME unset in ordinary suite; set-but-invalid or missing rect config fails. This is test-store drill, not attended release acceptance. PD1 unresolved visibility must red and be reported, never guessed rectangles.

Exact focused and repository gate commands (environment/helper above applies):

```sh
go test -run '^TestMarketingCaptureFakeChromeSixShots$' ./scripts/marketingcheck -count=1 -timeout 120s -v
go test -run '^TestMarketingCaptureInventoryMutation$' ./scripts/marketingcheck -count=1 -timeout 120s -v
go test -run '^TestMarketingProtectedObservationComparator$' ./scripts/marketingcheck -count=1 -timeout 90s -v
go test -run '^TestMarketingChromeSixShots$' ./scripts/marketingcheck -count=1 -timeout 660s -v
go test -run '^TestWorkbenchGraphWiring$' ./host/daemon -count=1 -timeout 45s -v
go test -run '^TestWorkbenchDecisionsPane$' ./host/daemon -count=1 -timeout 45s -v
go test -run '^TestWorkbenchLiveCursorCrossesGap$' ./host/daemon -count=1 -timeout 45s -v
```

Mutations; each must reach this assertion:

| Mutant | Valid edit | Test / assertion that must red |
|---|---|---|
| M7c-inventory | Neuter frozen public inventory equality guard; real extra-commit control runs. | `TestMarketingCaptureInventoryMutation`: extra commit invalidates capture |
| M7c-observation | Neuter SHA256 equality half of observation comparator; same-newline-count copied input differs. | `TestMarketingProtectedObservationComparator`: same-line-count changed copy must red |
| M7c-line-observation | Neuter line-count equality guard in copied observation comparator; keep digest-equal instrument records with different counts | `TestMarketingProtectedObservationComparator`: line-count mismatch must refuse |
| M7c-scheme | Drop preferredColorScheme flag from dark shot argv while fake Chrome records args and fake pixels enforce expected scheme. | `TestMarketingCaptureFakeChromeSixShots`: every shot explicitly pins correct scheme |

Exact mutation commands:

```sh
go test -run '^TestMarketingCaptureInventoryMutation$' ./scripts/marketingcheck -count=1 -timeout 120s -v
go test -run '^TestMarketingProtectedObservationComparator$' ./scripts/marketingcheck -count=1 -timeout 90s -v
go test -run '^TestMarketingCaptureFakeChromeSixShots$' ./scripts/marketingcheck -count=1 -timeout 120s -v
```

Base focused gates for this design milestone:

- **PASS** `go test -run '^TestWorkbenchGraphWiring$' ./host/daemon -count=1 -timeout 45s -v` (0.94s).
- **PASS** `go test -run '^TestWorkbenchDecisionsPane$' ./host/daemon -count=1 -timeout 45s -v` (0.71s).
- **PASS** `go test -run '^TestWorkbenchLiveCursorCrossesGap$' ./host/daemon -count=1 -timeout 45s -v` (0.76s).
- **SKIP (not acceptance)** `go test -run '^TestWorkbenchLiveChromeDrill$' ./host/daemon -count=1 -timeout 45s -v` (0.66s).

New named tests/source modules are future implementation, absent at base; no base PASS is claimed for them.

### Executor directive notes

Narrowest gate list: exactly the command block above, plus each listed mutation command and before/after C11 observation. May touch only the exact files listed for M7c and its scratch snapshots. Repository gates: focused affected-package gates now; full build/test, verify_ail and host races at M8b (base full test remains red; never skip its failures). TestMarketingChromeSixShots needs WORLD_CHROME; TestMarketingCaptureFakeChromeSixShots is its mandatory always-on fake-script companion. Acceptance list includes the real drill but unattended executor runs it only with explicitly supplied browser+measured rectangles; otherwise record DEFERRED/SKIP, never C8 product PASS. Controller owns measured geometry and PD1 correction; do not spend executor time inventing capture layout.

**no git write operations; snapshot changed files into $HOME/.ailang/state/world-iter<N>/snap/M7c/ after each milestone**. Copy changed allowlisted files with their relative paths; include boundary command rc/wall-time, named PASS evidence, mutation ledger and observation records there. Never snapshot keys/session/output credentials into tracked assets. Maximum30min; incomplete is reported, never relabelled PASS.

## M8a — Reproduction docs, captions and release manifest schema

Design M8; 25 minutes; dependency: M7c; status **pending**.

Exact executor files:

- `docs/QUICKSTART.md`
- `README.md`
- `website/docs/getting-started/marketing-capture.md`
- `website/static/img/v1/captions.md`
- `docs/marketing-capture-manifest.schema.json`
- `cmd/ailang-worldd/demo_test.go`
- `scripts/marketingcheck/docs_test.go`

Ship concrete reproduction sequence and all payload construction. Captions truthful before release; release sibling capture-manifest.json is generated from attended actual bytes only, not checked-in dummy hashes. Manifest schema/example distinguished from public capture record. Docs link test fixture mode discharges wiring only; actual asset resolution deferred until promotion.

Acceptance criteria and named assertions:

- **C10 · `TestMarketingDocsLinks`:** Resolve existing reproduction-doc links and exact six image targets against a disposable simulated promoted directory before release; runtime release mode requires actual six files and hashes, never fixture placeholders. README light/packet links and website /img/v1 names/alt captions match inventory. Disclose static capture, demo-authored claim, evidence limits and separate pins.
- **C10 · `TestMarketingRunbookMatchesHelp`:** Parse documented flags/payload preparation then execute CLI/script help and test credential recipe through existing V11 setup; all seed/capture/proof/still/tool flags match. No unattended typed-y confirmation or production publisher confirmation. Unminted seed command refuses.

Exact focused and repository gate commands (environment/helper above applies):

```sh
go test -run '^TestMarketingDocsLinks$' ./scripts/marketingcheck -count=1 -timeout 90s -v
go test -run '^TestMarketingRunbookMatchesHelp$' ./scripts/marketingcheck -count=1 -timeout 120s -v
go test -run '^TestQuickstartDecisionPacketSection$' ./host/daemon -count=1 -timeout 45s -v
go test -run '^TestServeWorkspaceFlags$' ./cmd/ailang-worldd -count=1 -timeout 45s -v
go test -run '^TestServeProofFlags$' ./cmd/ailang-worldd -count=1 -timeout 45s -v
```

Mutations; each must reach this assertion:

| Mutant | Valid edit | Test / assertion that must red |
|---|---|---|
| M8a-filename | Change one doc image filename to a valid different nonexistent name in copied doc input. | `TestMarketingDocsLinks`: documented image target must match promoted inventory |
| M8a-flag | Change documented --workbench-still to unknown --workbench-static in copied recipe; invoke actual parser. | `TestMarketingRunbookMatchesHelp`: documented serve argv must parse and select still mode |

Exact mutation commands:

```sh
go test -run '^TestMarketingDocsLinks$' ./scripts/marketingcheck -count=1 -timeout 90s -v
go test -run '^TestMarketingRunbookMatchesHelp$' ./scripts/marketingcheck -count=1 -timeout 120s -v
```

Base focused gates for this design milestone:

- **PASS** `go test -run '^TestQuickstartDecisionPacketSection$' ./host/daemon -count=1 -timeout 45s -v` (0.68s).
- **PASS** `go test -run '^TestServeWorkspaceFlags$' ./cmd/ailang-worldd -count=1 -timeout 45s -v` (0.69s).
- **PASS** `go test -run '^TestServeProofFlags$' ./cmd/ailang-worldd -count=1 -timeout 45s -v` (0.69s).

New named tests/source modules are future implementation, absent at base; no base PASS is claimed for them.

### Executor directive notes

Narrowest gate list: exactly the command block above, plus each listed mutation command and before/after C11 observation. May touch only the exact files listed for M8a and its scratch snapshots. Repository gates: focused affected-package gates now; full build/test, verify_ail and host races at M8b (base full test remains red; never skip its failures). 

**no git write operations; snapshot changed files into $HOME/.ailang/state/world-iter<N>/snap/M8a/ after each milestone**. Copy changed allowlisted files with their relative paths; include boundary command rc/wall-time, named PASS evidence, mutation ledger and observation records there. Never snapshot keys/session/output credentials into tracked assets. Maximum30min; incomplete is reported, never relabelled PASS.

## M8b — Final repository gates and protected-path observation

Design M8; 25 minutes; dependency: M8a; status **pending**.

Exact executor files:

- `docs/QUICKSTART.md`
- `README.md`
- `website/docs/getting-started/marketing-capture.md`
- `website/static/img/v1/captions.md`
- `docs/marketing-capture-manifest.schema.json`

Final validation run, no new product scope. Re-run the narrow cross-boundary tests and mandatory repo gates. Use observation snapshots, never restore protected files. Report every gate rc/wall time and all unresolved product criteria. Attended capture/promotion is outside milestones.

Acceptance criteria and named assertions:

- **C11 · `TestMarketingProtectedObservationComparator`:** Re-drill copied-input comparator; final actual protected inventory equals initial sprint observation in both digest and newline count. Scoped read-only diff must be empty for all protected paths.

Exact focused and repository gate commands (environment/helper above applies):

```sh
go test -run '^TestMarketingProtectedObservationComparator$' ./scripts/marketingcheck -count=1 -timeout 90s -v
go test -run '^TestDemoSeedEndToEnd$' ./cmd/ailang-worldd -count=1 -timeout 120s -v
go test -run '^TestDemoSeedWhyNamedLinks$' ./cmd/ailang-worldd -count=1 -timeout 120s -v
go test -run '^TestWorkbenchPacketEvidenceOnlyThroughRowOneOhFive$' ./host/daemon -count=1 -timeout 45s -v
go test -run '^TestMarketingCaptureFakeChromeSixShots$' ./scripts/marketingcheck -count=1 -timeout 120s -v
go test -run '^TestMarketingPNGStructure$' ./scripts/marketingcheck -count=1 -timeout 90s -v
AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_ail.sh
go build ./... && go test ./... -count=1 -timeout 600s
go test -race ./host/daemon -count=1 -timeout 480s
go test -race ./host/coordinator ./host/workbench ./host/evidence ./scripts/marketingcheck -count=1 -timeout 480s
```

Mutations; each must reach this assertion:

| Mutant | Valid edit | Test / assertion that must red |
|---|---|---|
| M8b-C11-redrill | Re-drill M7c line-count guard mutant on copied observations | `TestMarketingProtectedObservationComparator`: line-count mismatch must refuse |

Exact mutation command: `go test -run '^TestMarketingProtectedObservationComparator$' ./scripts/marketingcheck -count=1 -timeout 90s -v`.

Base focused gates for this design milestone:

- **PASS** `go test -run '^TestQuickstartDecisionPacketSection$' ./host/daemon -count=1 -timeout 45s -v` (0.68s).
- **PASS** `go test -run '^TestServeWorkspaceFlags$' ./cmd/ailang-worldd -count=1 -timeout 45s -v` (0.69s).
- **PASS** `go test -run '^TestServeProofFlags$' ./cmd/ailang-worldd -count=1 -timeout 45s -v` (0.69s).

New named tests/source modules are future implementation, absent at base; no base PASS is claimed for them.

### Executor directive notes

Narrowest gate list: exactly the command block above, plus each listed mutation command and before/after C11 observation. May touch only the exact files listed for M8b and its scratch snapshots. Repository gates: AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_ail.sh; go build ./... && go test ./... -count=1 -timeout 600s; go test -race ./host/daemon -count=1 -timeout 480s; go test -race ./host/coordinator ./host/workbench ./host/evidence ./scripts/marketingcheck -count=1 -timeout 480s. Full mandatory gates are paid once here, plus focused/race gates at affected earlier boundaries. CI runs host race checks; daemon alone is ~200s, not a 45s test. Total controller outer bound900s for full build/test,540s per race command. If this boundary exceeds30min report incomplete and split validation from doc corrections without weakening gates.

**no git write operations; snapshot changed files into $HOME/.ailang/state/world-iter<N>/snap/M8b/ after each milestone**. Copy changed allowlisted files with their relative paths; include boundary command rc/wall-time, named PASS evidence, mutation ledger and observation records there. Never snapshot keys/session/output credentials into tracked assets. Maximum30min; incomplete is reported, never relabelled PASS.

## Final disposition

Unattended implementation and test gates can be completed without release confirmations. Release C8 pixels/byte equality, actual image/manifest promotion and C10 actual asset resolution remain controller-owned attended-release checks. PD1 and any M7 fallback residual must appear in the boundary report; no claim of full showcase completion or clause movement follows from tool/test success. All JSON milestone statuses remain pending.
