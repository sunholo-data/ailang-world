# Sprint plan — w-marketing-capture

**Iteration 266 · queue row 150 · base `3e6f6ae` · plan only**

Design: [w-marketing-capture.md](w-marketing-capture.md), and for every pending milestone (M7c0–M8b) [w-marketing-capture-orchestration.md](w-marketing-capture-orchestration.md) r2, whose §5 table governs where it and this plan disagree (amended iteration 305 per w-marketing-capture-orchestration §5). Read in full, including V1–V28, r1–r3 quorum dispositions and r4 bounded revision. R1/r2/r3 remain BLOCKED rounds; r4 uses the controller narrow-refinement carve-out, not a retroactive quorum pass. Binding: CLAUDE.md and coding-standards.md; package-first law, contracts before pure implementation, host effects, non-vacuous named gates and AILANG fluency before writing source.

Reference pair requested under planned/ is absent at base; both existing files were read under design_docs/implemented/. JSON copies its top-level structure and milestone files/tests_first/acceptance/mutations shape, adding the requested explicit commands, acceptance_criteria, dependency, estimate and status="pending" fields.

**22 executor runs / 563 estimated minutes** (amended iteration 305 per w-marketing-capture-orchestration §5): M7c0 (shot geometry, 20 min) added; was 21 / 543 at iteration 304 (M7c 28 min split into M7c1–M7c4, 126 min; was 17 / 445). M1a–M7b (14 runs) have landed (orchestration doc §2). Design M1–M8 purposes and dependency order are preserved; large stages are split below. Publish confirmation and session mint typed `y` are outside all milestones. Planner writes only this file and [sprint_w-marketing-capture.json](sprint_w-marketing-capture.json).

## Environment, gate semantics and binary discovery

```sh
export AILANG_BIN=$HOME/.pinned-ailang/ailang
export WORLD_TOOL_AILANG_BIN=$HOME/.pinned-ailang-tools/v0.54.15/ailang   # (amended iteration 305 per w-marketing-capture-orchestration §4.2); was v0.52.1
run_ail() { python3 -c 'import subprocess,sys; sys.exit(subprocess.run(sys.argv[1:],timeout=40).returncode)' "$@"; }
```

The pin was measured as v0.41.0/24ee108. The separate tool binary was v0.52.1 at iteration 266; since `4b686d8e` (row 139, 2026-10-17) it is **v0.54.15**/0384450, and seed tests fail fast against v0.52.1 with "want release AILANG v0.54.15" (amended iteration 305 per w-marketing-capture-orchestration §4.2). host/daemon/workspace.go:52 (`var ToolBinaryRelease = pins.VersionLine(pins.RoleTool)`) requires it; host/daemon/workspace_test.go:593 locates WORLD_TOOL_AILANG_BIN and skips if unset. host/daemon/setools_e2e_test.go:58-86 and cmd/ailang-worldd/cliwalk_e2e_test.go:33-44 require both env vars or explicitly skip; configured wrong versions fail. TestSeToolsE2EBinaryVersions checks both. FakeToolBin-based workspace unit tests need no real tool but cannot prove S3.

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
| M7c0 | Shot geometry: tall packet/graph/provenance heights from measured bottoms (amended iteration 305 per w-marketing-capture-orchestration §5) | 20 |
| M7c1 | Protected-file observation comparator (production) and its mutants | 20 |
| M7c2a | Python capture primitives (Chrome version, HTTP readiness, six-shot pass) and Go `routes` | 28 |
| M7c2b | Capture shell end-to-end with fake Chrome and a fake daemon | 28 |
| M7c3 | Frozen public inventory guard (real store, genuine extra commit) | 25 |
| M7c4 | Real Chrome drill (test written by the executor; RUN by the controller) | 25 |
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

- **C2 · `TestDemoSeedReadTestCheck`:** Both real pins required. Publish existing se-tools descriptors through production publisher core in test setup; real /a2a/ dispatch for demo-read/test/check. Read bytes equal stored logepoch; test argv op=test flags json:"", named passes>0 and no failures/skips/vacuity; ai_check actual non-null success. Retain actual pin, argv, stdout/stderr, invocation/effect refs; budget deltas Read=1 CLI=2 and no extra calls. Prove CLI result is v0.52.1, not proof pin. (Historical, landed M2b: the tool pin is v0.54.15 since `4b686d8e` (amended iteration 305 per w-marketing-capture-orchestration §4.2).)
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

Narrowest gate list: exactly the command block above, plus each listed mutation command and before/after C11 observation. May touch only the exact files listed for M2b and its scratch snapshots. Repository gates: focused affected-package gates now; full build/test, verify_ail and host races at M8b (base full test remains red; never skip its failures). TestDemoSeedReadTestCheck needs WORLD_TOOL_AILANG_BIN v0.52.1 (v0.54.15 since `4b686d8e` (amended iteration 305 per w-marketing-capture-orchestration §4.2)). Ordinary missing-pin run may SKIP; required milestone acceptance preflight must reject unset pin and require top-level PASS. Bad configured version is FAIL, never skip.

**no git write operations; snapshot changed files into $HOME/.ailang/state/world-iter<N>/snap/M2b/ after each milestone**. Copy changed allowlisted files with their relative paths; include boundary command rc/wall-time, named PASS evidence, mutation ledger and observation records there. Never snapshot keys/session/output credentials into tracked assets. Maximum30min; incomplete is reported, never relabelled PASS.

## M3a — HTTP proof adapters and real producer success

Design M3; 28 minutes; dependency: M2b; status **completed** (iteration 272, `1ad2743`).

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

Design M5; 29 minutes; dependency: M5a; status **completed** — built iterations 296–298, PR #353 (executor codex gpt-6.1-sol in 296; controller test bindings `66208c2a` and the r1 fix `f9503386`). Judge Agent sonnet r1 PASS 86 → r2 PASS 90, zero blocking. PROVEN's subject comes from the envelope's MAC-authenticated report (no sibling pair: the S5 producer emits four pairs and the package pins `additionalProperties:false`). It is shown only once PROVEN.

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

Design M5; 23 minutes; dependency: M5b; status **completed** — built iteration 299, PR #354 (executor codex gpt-6.1-sol `137b615d`; controller `c9fd10af` fixes M5b's TESTED rule to the real v0.54.15 tool shape, which omits `exit_code` on success; controller tests `be270368` for judge r1's survivors). Judge Agent sonnet r1 PASS 80 → r2 PASS 85, zero blocking. Deviations: `demo.go`, QUICKSTART §12 and `website/docs/reference/cli.md` edited to keep help/docs truthful; `host/daemon/workbench_packet_evidence.go` got two production fixes (resultless effect, TESTED exit_code).

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

Design M6; 25 minutes; dependency: M5c; status **completed** — built iteration 300 (executor codex gpt-6.1-sol `aeaa957d`; controller docs fix `ac1abe7f` for judge r1 N3). Judge Agent sonnet r1 PASS 93, zero blocking. Deviation: `website/docs/guides/operating-the-daemon.md` and `website/docs/reference/cli.md` also name the flag (controller-authorised, docs only). The template branch is `{{if .Still}}{{else}}…{{end}}` because the field census reads a direct bool branch.

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

Design M7; 28 minutes; dependency: M6; status **completed** — built iteration 301 (executor codex gpt-6.1-sol `38ac4167`; controller judge-r1 fix `98f48b20`: a signal now stops a hung shot (`Interrupted`) plus branch subtests; controller load-robust test bounds `dd647909`/`c70b58b9` after `verify_go.sh` measured timing reds). Judge Agent sonnet r1 PASS 86 → r2 PASS 91 → r3 PASS 91 (load-bound delta), zero blocking. Interface: `marketing_capture.py shot|daemon`, one JSON verdict; `capture_marketing.sh --unit-daemon` unit control; a real capture refuses rc 1 until M7c.

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

Design M7; 28 minutes; dependency: M7a; status **completed** — built iteration 302 (executor codex gpt-6.1-sol `fa04df00` — the run completed in the first 2026-10-22 fire whose controller then crashed at gate 3; recovered and verified first-party by the iteration-302 controller: five gates green, all eight plan mutants red at their named tests with sha-identical restores; controller fix `02e12be2`: the two verifier files enter the Go-surface allow-list — `TestProductionGoSurface` reddened in the full suite, the one finding neither the focused gates nor the round-1 judge caught). Judge pi:openrouter/minimax/minimax-m3 (generator codex ≠ judge minimax) r1 PASS 95, zero blocking, one honest survivor (M7b-truncated-chunk: neutering the guard reds at a different assertion; per the panic rule, not a kill) → r2 PASS 95 on the fix delta, zero new blocking, both guard drills red (rogue file; entries removed). Carried to M7c: N1, and from M7a N4/S6/S7/R3-1 — `TestMarketingSupervisorFakeDaemonCleanup/shell-hup` still exceeds its 10 s bound under whole-suite load (29 s measured 3/3; green isolated 4/4).

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

## M7c — split (iteration 304 planner)

The single 28-minute M7c (capture orchestration, frozen inventories and the Chrome drill) is **replaced by five children**: M7c1, M7c2a, M7c2b, M7c3 and M7c4. The designer audit `design_docs/verification/world-iter304/designer-m7c-readiness.md` recommended four (M7c1–M7c4). The planner split M7c2 into M7c2a and M7c2b because the orchestration work alone does not fit one ≤30-minute, ≤~250-line run for claude-haiku-5-5:
- Python primitives: version read, HTTP readiness and the six-shot pass loop.
- A Go `routes` subcommand.
- Shell flag parsing and diagnostics.
- A fake `ailang-worldd`.

M7c2a holds the browser-free Python/Go primitives. M7c2b holds only the shell chain end-to-end. Dependencies are linear: M7c1 → M7c2a → M7c2b → M7c3 → M7c4 → M8a. M7c3 depends on M7c2b only because both edit `scripts/marketingcheck/main.go`; it would otherwise be independent.

**Iteration 305 (amended iteration 305 per w-marketing-capture-orchestration §5).** The remaining work is governed by `w-marketing-capture-orchestration.md` r2. Its §5 table overrides this plan where they disagree; each amended spot below is marked. Its quorum disposition is the Gate-2 narrow-refinement carve-out (doc §9), not a quorum PASS. Under C13 that disposition releases M7c1–M8b. M7c0 additionally waits for R0's re-confirm V-row.
- **New M7c0** (shot geometry) depends only on R0 and shares no file with M7c1, so either can go first.
- **M7c2a depends on M7c1 and M7c0** (planner decision): its `routes` rows carry the M7c0 dimensions.
- **Chain:** M7c1 ∥ (R0 → M7c0) → M7c2a → M7c2b → M7c3 → M7c4 → M8a → M8b.
- **Controller steps:** R1 (all-six rectangles, full-path C14, threshold V-row) follows M7c4. P0/P3 and A1 are attended-release steps (doc §4.4, §5), never executor milestones.
- **Tool pin** for every pending milestone and the attended capture: `$HOME/.pinned-ailang-tools/v0.54.15/ailang` (doc §4.2).

Shared rules for every child:
- **Interfaces at HEAD.** Verifier CLI: `check --pass-a --pass-b --manifest --rectangles [--out]`. Exit codes: 0 ok, 1 refusal, 2 usage. stderr: `marketingcheck: <code> <name>: <detail>` (`scripts/marketingcheck/main.go:13-42`). Supervisor: one JSON line `{verdict, detail[, shot_verdict]}`, exit 0 only for `ok` (`scripts/marketing_capture.py:236-237`). Production deadlines are the constants at `marketing_capture.py:32-36` (shot 30 s, pass 300 s, startup 30 s, TERM 5 s, KILL 5 s). Reuse them; never invent new production deadlines.
- **Go-surface guard.** `TestProductionGoSurface` (`host/store/context_roots_test.go:198-245`) reds on any new non-test `.go` file under `scripts/`. Its allow-list holds only `scripts/marketingcheck/main.go` and `png.go` (lines 204-206). Every child below therefore puts production Go in `main.go`, so **no allow-list edit is needed**. If an executor nonetheless creates a new non-test `.go` file there, it must add that path to the allow-list map in the same change. That is the M7b controller-fix lesson (`02e12be2`).
- **Subprocess sinks.** A Go test that sets `cmd.Stdout` or `cmd.Stderr` on a started subprocess uses an `os.Create`d file under `t.TempDir()`, never a bare `bytes.Buffer` (`TestSubprocessSinkGate`, `host/verifygate/subprocess_sink_gate_test.go:223`; precedent `supervisor_test.go:388`).
- **Timing.** Every test time bound is ≥10× its unloaded time. The package runs at load 7–25 in the full suite.
- **Foreign-format fixtures.** A fixture that models another tool's output is generated by running that tool inside the test, never typed by hand.
- **Carried flake.** `TestMarketingSupervisorFakeDaemonCleanup/shell-hup` has a load flake and is not owned by any child; see the planner verification file. An executor that sees it red in a whole-package run reports it as a carried base finding, and must not edit `supervisor_test.go` to make it pass.
- **Not product evidence.** No child makes C8/C9 product evidence met. Only the controller rehearsal on real Chrome with measured rectangles can do that.

## M7c0 — Shot geometry (new, iteration 305)

(added iteration 305 per w-marketing-capture-orchestration §4.1/§5). Design M7 (C8′); 20 minutes; dependency: M7b (landed) and **R0's re-confirm V-row** (controller, doc §5 R0, C13; not a milestone id, so the JSON records it in `implementation`). It does not depend on M7c1; the two share no file, and their relative order is free. Status **pending**: do not dispatch it before the R0 V-row exists. If R0's max differs from V51, transcribe R0's values and cite its row. If any B + 40 > 4,400, the shot is BLOCKED and escalated to Mark; M7c0 does not run and the height is never clamped.

Why: at HEAD the packet, graph and provenance shots do not contain their required region at 1440×1400/1100 (doc V42). The heights come from the doc's rule H = ⌈(B + 40)/100⌉ × 100, where B is the maximum across two seeds of the region's `getBoundingClientRect().bottom` (doc V51: packet 3564.56 → 3700, graph 2042.14 → 2100, provenance 2881.17 → 3000).

### Files (exactly these two)

- `scripts/marketingcheck/png.go`: in `Shots` (`png.go:49-53`), change **only** the three heights: packet `1400` → `3700`, graph `1100` → `2100`, provenance `1400` → `3000`. Nothing else in `png.go` changes. Light/dark stay 1440×1100 and phone stays 390×844.
- `scripts/marketingcheck/png_test.go`:
  1. Add `"math"` to the import block (`png_test.go:3-21`).
  2. Update the canonical list at `png_test.go:462` (subtest `TestMarketingPNGContent/canonical-shots`) to `{"packet.png", 1440, 3700, "light"}`, `{"graph.png", 1440, 2100, "light"}`, `{"provenance.png", 1440, 3000, "light"}`. The other three entries are unchanged.
  3. Append the constants and `TestMarketingShotGeometry` below.

**Do NOT touch:**
- The generic 1440×1100 cases at `png_test.go:321-337`, `:366` and `:383`. They are bound to `workbench-light.png`, not to the tall shots.
- The x-coordinate `1400` at `:357`.
- `fixture()`/`fixtureSet()`: they already size from `Shots` (`png_test.go:147`, `:164`).
- `main.go`, `supervisor_test.go`, the scripts, or any protected path.

```go
// Measured bottom edge (CSS px, device scale 1, getBoundingClientRect().bottom at
// scroll 0) of each tall shot's last required region, max across two independently
// seeded stores. Source: w-marketing-capture-orchestration §10 V51 (iteration 305)
// and the R0 re-confirm V-row <cite the R0 row id here>. A Go test cannot read a
// V-row; these constants are the transcription.
const (
	measuredPacketBottom     = 3564.56 // section[aria-label="decision packet"]
	measuredGraphBottom      = 2042.14 // section[aria-label="world graph"]
	measuredProvenanceBottom = 2881.17 // pre.payload
)

func shotRuleHeight(bottom float64) int { return int(math.Ceil((bottom+40)/100)) * 100 }

func TestMarketingShotGeometry(t *testing.T) {
	for _, tc := range []struct {
		name      string
		bottom    float64
		assertion string
	}{
		{"packet.png", measuredPacketBottom, "packet height must equal the rule value from measuredPacketBottom"},
		{"graph.png", measuredGraphBottom, "graph height must equal the rule value from measuredGraphBottom"},
		{"provenance.png", measuredProvenanceBottom, "provenance height must equal the rule value from measuredProvenanceBottom"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			s, ok := shotByName(tc.name)
			want := shotRuleHeight(tc.bottom)
			if !ok || s.Width != 1440 || s.Height != want {
				t.Fatalf("%s: got %+v, want 1440x%d", tc.assertion, s, want)
			}
			if want > 4400 {
				t.Fatalf("%s needs %d px > 4,400: BLOCKED, escalate; never clamp", tc.name, want)
			}
			if tc.bottom > float64(s.Height-40) {
				t.Fatalf("%s: required bottom %.2f must be <= H-40 = %d", tc.name, tc.bottom, s.Height-40)
			}
		})
	}
	t.Run("unchanged", func(t *testing.T) {
		for _, want := range []Shot{{"workbench-light.png", 1440, 1100, "light"}, {"workbench-dark.png", 1440, 1100, "dark"}, {"phone.png", 390, 844, "light"}} {
			if got, ok := shotByName(want.Name); !ok || got != want {
				t.Fatalf("light, dark and phone geometry must stay unchanged: got %+v want %+v", got, want)
			}
		}
	})
}
```

### Acceptance commands (expected results)

```sh
gofmt -l scripts/marketingcheck                                                              # no output
go vet ./scripts/marketingcheck                                                              # rc 0
go test -run '^TestMarketingShotGeometry$' ./scripts/marketingcheck -count=1 -timeout 60s -v  # PASS: packet.png, graph.png, provenance.png, unchanged
go test -run '^TestMarketingPNGContent$' ./scripts/marketingcheck -count=1 -timeout 90s -v    # PASS incl. canonical-shots and positive-packet.png/-graph.png/-provenance.png at the new sizes
go test ./scripts/marketingcheck -count=1 -timeout 240s -v                                    # all tests PASS (planner scratch probe at load 1.7: 17.5 s, every test PASS with the three heights changed)
go test -race ./scripts/marketingcheck -count=1 -timeout 240s
```

| Mutant | Valid edit (production `png.go`) | Command | Must red with |
|---|---|---|---|
| M7c0-height | `{"packet.png", 1440, 3700, "light"}` → `{"packet.png", 1440, 1400, "light"}` | `go test -run '^TestMarketingShotGeometry$/^packet.png$' ./scripts/marketingcheck -count=1 -timeout 60s -v` | `packet height must equal the rule value from measuredPacketBottom: got {Name:packet.png Width:1440 Height:1400 Scheme:light}, want 1440x3700`. `canonical-shots` also reds (`canonical shot packet.png: got … want …`), but the named kill is the equality assertion, not the fixture or cap checks. Planner scratch probe: both observed. |

Restore by scratch copy and `cmp` (no git), then re-run green. The C11 observation runs before and after, as for every milestone.

**Executor notes:** stay within the two files. A changed PNG fixture size is expected: `fixture()` follows `Shots`, and every PNG stays under the 1 MiB cap. Snapshot to `snap/M7c0/`. Maximum 30 min. No git write operations.

## M7c1 — Protected-file observation comparator (production) and its mutants

Design M7 (C11); 20 minutes; dependency: M7b (landed, PR #361); status **pending**; **dispatchable now** (amended iteration 305 per w-marketing-capture-orchestration §5). The iteration-304 HOLD is lifted: `w-marketing-capture-orchestration` r2 routed to the planner under the Gate-2 narrow-refinement carve-out (doc §9; C13). That is not a quorum PASS. M7c1 shares no file with M7c0 and does not wait for R0.

**Re-verified at worktree HEAD `ff16bf65` (iteration 305 planner).** Every file:line below was read with `sed -n`/`awk` at that commit, and the whole section was built once in a scratch copy outside the repository (`/tmp/m7c1-proto-iter305`). There, the test went 11/11 PASS, each of the three mutants went red on its named assertion and was restored green, and `go test -race ./scripts/marketingcheck` went `ok` in 33 s at load 1.5. Receipt: `design_docs/verification/world-iter305/planner-receipt.md`. The code blocks below are that prototype. Type them in as given. If any line will not compile or behave as stated, stop and report it; do not improvise a different design.

### Files (exactly these two; nothing else)

| File | Action | Size |
|---|---|---|
| `scripts/marketingcheck/main.go` | modify: replace original lines 1-16 (`package main` through the closing `}` of the `args[0] != "check"` block) with block A, keeping original line 17 (`fs := flag.NewFlagSet("check", …)`) as the first statement of `runCheck`; in the `check` body change `fs.Parse(args[1:])` (line 23) to `fs.Parse(args)`; append block B at the end of the file | 83 → ~297 lines |
| `scripts/marketingcheck/observation_test.go` | create, exactly block C | ~255 lines |

**Do NOT touch:**
- `scripts/marketingcheck/png.go`, `png_test.go` or `supervisor_test.go`;
- `scripts/capture_marketing.sh` or `scripts/marketing_capture.py`;
- `host/store/context_roots_test.go`: no allow-list edit, because all production Go goes into the already-allowed `main.go`, and a new non-test `.go` file under `scripts/` reds `TestProductionGoSurface`;
- any protected path: `world/*.ail`, `go.mod`, `go.sum`, `design_docs/world-mission.md`, `tools/launchd/**`. `go.mod` being protected also means no new module dependency; stdlib only;
- this plan, the design docs, and the plan's C11 Python heredoc (section "C11 observations and mutation discipline"). Block C copies the heredoc verbatim as a Go raw string; the parity subtest pins the Go observer to it. A mismatch is a finding to report, never a reason to edit either side.

### Facts this section relies on (verified at `ff16bf65`)

| Cite | What is there |
|---|---|
| `main.go:13` | `if len(args) == 0 \|\| args[0] != "check" {` (the single-command gate block A replaces) |
| `main.go:23` | `if err := fs.Parse(args[1:]); err != nil {` |
| `main.go:36` | `fmt.Fprintf(os.Stderr, "marketingcheck: %s %s: %s\n", r.Code, r.Name, r.Detail)`, the refusal format reused by `compare` |
| `main.go:73` | `os.OpenFile(out, os.O_WRONLY\|os.O_CREATE\|os.O_EXCL, 0600)`, the O_EXCL pattern reused by `observe --out` |
| `png.go:19-22` | `type Refusal struct{ Code, Name, Detail string }`; `func refuse(code, name, detail string) error` |
| `png_test.go:23-29` | `func wantCode(t *testing.T, err error, code, assertion string)` fails with `"<assertion>: want Code=<code>, got <err>"` |
| `png_test.go:30` | `func pngOK(t *testing.T, err error)` (t.Fatal on non-nil) |
| `png_test.go:734-741` | test-time `go build -o bin .` with a 60 s context (pattern reused by `cli`) |
| `png_test.go:800` | `TestMarketingTwoPassInventory/cli/unknown` passes `docs` and expects rc 2, so do **not** add a `docs` subcommand (M8a owns that) |
| `supervisor_test.go:25-30` | `python3` is mandatory for this package (`t.Fatal("python3 is mandatory: ", err)`) |
| `supervisor_test.go:388` / `host/verifygate/subprocess_sink_gate_test.go:223` | file sinks for subprocess output; `TestSubprocessSinkGate` flags bare buffers on a `Start()`ed command |
| `host/store/context_roots_test.go:198-245`, allow-list `:204-206` | only `scripts/marketingcheck/main.go` and `png.go` are allowed non-test Go under `scripts/` |
| plan "C11 observations and mutation discipline", lines 119-124 | the Python heredoc copied into block C |
| package identifiers | `go vet` at base is clean; none of the new names (`usageText`, `runCheck`, `runObserve`, `runCompare`, `ObservationRecord`, `Observe`, `observationPaths`, `CompareObservations`, `writeObservation`, `readObservations`, `observationPython`, `protectedFixture`, `protectedPaths`, `protectedTree`, `flipProtectedByte`, `mustObserve`, `pythonObserve`, `runObservationCLI`) exists in the package |

Real protected set at HEAD (plan Python heredoc run at the worktree root): **12 paths**: `design_docs/world-mission.md`, `go.mod`, `go.sum`, five `tools/launchd/` files (including `testdata/planner-lane/n-no-backtick-bullet.md`), and `world/{contracts,logepoch,transitions,types}.ail`. There are no symlinks under `tools/launchd` or `world`. Python 3.14.5; go1.26.6.

### Interface (normative)

```go
type ObservationRecord struct {
	Path   string `json:"path"`   // slash-separated, relative to the observed root
	SHA256 string `json:"sha256"` // lowercase hex of the raw file bytes
	Lines  int    `json:"lines"`  // bytes.Count(raw, []byte{'\n'})
}
func Observe(root string) ([]ObservationRecord, error)
func CompareObservations(before, after []ObservationRecord) error
```

- **`Observe`** reads and never writes. It covers: `world/*.ail` (non-recursive; `os.ReadDir` + `filepath.Match("*.ail", name)`), `go.mod`, `go.sum`, `design_docs/world-mission.md`, and every non-directory entry under `tools/launchd/` (recursive `filepath.WalkDir`). It keeps a candidate only if `os.Stat` (which follows symlinks, like Python's `is_file()`) says regular file. A missing path or directory contributes nothing. It returns a non-nil slice sorted by `Path` (`[]` when empty, so `observe` prints `[]`, never `null`).
- **`CompareObservations`** returns refusal codes in this order:
  1. `observation-empty`: `len(before)==0`.
  2. `observation-set`: ONE guard. The sorted path lists differ, or `before` repeats a path. The refusal detail must reference both `beforePaths` and `afterPaths`, so the M7c1-set mutant still compiles.
  3. Loop over `before`, **looking each record up by path in `afterByPath := map[string]ObservationRecord` built from `after`. Never index by position, and never `continue` on a miss** (a miss yields the zero record): (amended iteration 305 per w-marketing-capture-orchestration §5)
     - `observation-digest`: `Name` = path, when the SHA256 differs.
     - `observation-lines`: `Name` = path, when the Lines differ. Digest and lines are two **separate** `if` statements, so each can be neutered alone.

  The path-keyed lookup is why the M7c1-set mutant reaches `observation-digest` for the deleted record instead of panicking with index-out-of-range.

### CLI (normative; argv and output)

| argv | stdout | stderr | rc |
|---|---|---|---|
| *(none)* or unknown first arg (e.g. `docs`) | — | `usageText` (three usage lines) | 2 |
| `check …` | unchanged from M7b | unchanged | unchanged (0/1/2) |
| `observe --root DIR` | `json.MarshalIndent(records, "", "  ")` + `\n` | — | 0 |
| `observe --root DIR --out FILE` | — (FILE created `O_CREATE\|O_EXCL`, mode 0600, same bytes) | — | 0 |
| `observe` with FILE existing, or a read error | — | `marketingcheck: observation-io <root>: <err>` | 1 |
| `observe` without `--root`, an unknown flag, or a positional arg | — | flag usage | 2 |
| `observe --help` | — | flag usage | 0 |
| `compare --before F --after F` (equal) | — | — | 0 |
| `compare …` refusal | — | `marketingcheck: <code> <name>: <detail>` (format of `main.go:36`) | 1 |
| `compare` with an unreadable file or unparsable JSON | — | `marketingcheck: observation-json <file>: <err>` | 1 |
| `compare` missing `--before`/`--after`, unknown flag, or positional | — | flag usage | 2 |

`compare` accepts both `observe` output and the plan's Python heredoc output (`json.dumps(..., sort_keys=True)`: the same three keys).

### Block A — replaces `main.go` lines 1-16 (imports, usage, dispatch); it ends with the `runCheck` header, and original line 17 onward follows unchanged

```go
package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

const usageText = "usage: marketingcheck check --pass-a DIR --pass-b DIR --manifest DEMO_MANIFEST --rectangles RECTS [--out CAPTURE_MANIFEST]\n" +
	"       marketingcheck observe --root DIR [--out FILE]\n" +
	"       marketingcheck compare --before FILE --after FILE"

func main() { os.Exit(run(os.Args[1:])) }
func run(args []string) int {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, usageText)
		return 2
	}
	switch args[0] {
	case "check":
		return runCheck(args[1:])
	case "observe":
		return runObserve(args[1:])
	case "compare":
		return runCompare(args[1:])
	}
	fmt.Fprintln(os.Stderr, usageText)
	return 2
}
func runCheck(args []string) int {
```

Then change line 23 (now inside `runCheck`) from `if err := fs.Parse(args[1:]); err != nil {` to `if err := fs.Parse(args); err != nil {`. Nothing else in `runCheck` or `check()` changes.

### Block B — append to the end of `main.go`

```go
// ObservationRecord is one C11 protected-path observation.
type ObservationRecord struct {
	Path   string `json:"path"`   // slash-separated, relative to the observed root
	SHA256 string `json:"sha256"` // lowercase hex of the raw file bytes
	Lines  int    `json:"lines"`  // bytes.Count(raw, []byte{'\n'})
}

// Observe reads (never writes) the C11 protected set under root.
func Observe(root string) ([]ObservationRecord, error) {
	var paths []string
	world, err := os.ReadDir(filepath.Join(root, "world"))
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		return nil, err
	}
	for _, e := range world {
		if ok, _ := filepath.Match("*.ail", e.Name()); ok {
			paths = append(paths, "world/"+e.Name())
		}
	}
	paths = append(paths, "go.mod", "go.sum", "design_docs/world-mission.md")
	err = filepath.WalkDir(filepath.Join(root, "tools", "launchd"), func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				return nil
			}
			return err
		}
		if d.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(root, p)
		if err != nil {
			return err
		}
		paths = append(paths, filepath.ToSlash(rel))
		return nil
	})
	if err != nil {
		return nil, err
	}
	records := []ObservationRecord{}
	for _, rel := range paths {
		full := filepath.Join(root, filepath.FromSlash(rel))
		info, err := os.Stat(full)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				continue
			}
			return nil, err
		}
		if !info.Mode().IsRegular() {
			continue
		}
		raw, err := os.ReadFile(full)
		if err != nil {
			return nil, err
		}
		sum := sha256.Sum256(raw)
		records = append(records, ObservationRecord{Path: rel, SHA256: hex.EncodeToString(sum[:]), Lines: bytes.Count(raw, []byte{'\n'})})
	}
	slices.SortFunc(records, func(a, b ObservationRecord) int { return strings.Compare(a.Path, b.Path) })
	return records, nil
}

func observationPaths(records []ObservationRecord) []string {
	paths := make([]string, 0, len(records))
	for _, r := range records {
		paths = append(paths, r.Path)
	}
	slices.Sort(paths)
	return paths
}

// CompareObservations is the C11 equality comparator.
func CompareObservations(before, after []ObservationRecord) error {
	if len(before) == 0 {
		return refuse("observation-empty", "before", "an empty observation proves nothing")
	}
	beforePaths, afterPaths := observationPaths(before), observationPaths(after)
	if !slices.Equal(beforePaths, afterPaths) || len(slices.Compact(slices.Clone(beforePaths))) != len(beforePaths) {
		return refuse("observation-set", "paths", fmt.Sprintf("before has %d paths, after has %d; the protected set changed or repeats a path", len(beforePaths), len(afterPaths)))
	}
	afterByPath := make(map[string]ObservationRecord, len(after))
	for _, a := range after {
		afterByPath[a.Path] = a
	}
	for _, b := range before {
		a := afterByPath[b.Path] // zero value when absent; never skip a miss
		if b.SHA256 != a.SHA256 {
			return refuse("observation-digest", b.Path, fmt.Sprintf("sha256 %s became %s", b.SHA256, a.SHA256))
		}
		if b.Lines != a.Lines {
			return refuse("observation-lines", b.Path, fmt.Sprintf("lines %d became %d", b.Lines, a.Lines))
		}
	}
	return nil
}

func runObserve(args []string) int {
	fs := flag.NewFlagSet("observe", flag.ContinueOnError)
	root := fs.String("root", "", "repository root to observe (read-only)")
	out := fs.String("out", "", "new absent output path (default stdout)")
	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 2
	}
	if fs.NArg() != 0 || *root == "" {
		fs.Usage()
		return 2
	}
	records, err := Observe(*root)
	if err == nil {
		err = writeObservation(records, *out)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "marketingcheck: observation-io %s: %s\n", *root, err)
		return 1
	}
	return 0
}

func writeObservation(records []ObservationRecord, out string) error {
	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	if out == "" {
		_, err = os.Stdout.Write(data)
		return err
	}
	f, err := os.OpenFile(out, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	closeErr := f.Close()
	if err != nil {
		return err
	}
	return closeErr
}

func readObservations(path string) ([]ObservationRecord, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, refuse("observation-json", path, err.Error())
	}
	var records []ObservationRecord
	if err := json.Unmarshal(raw, &records); err != nil {
		return nil, refuse("observation-json", path, err.Error())
	}
	return records, nil
}

func runCompare(args []string) int {
	fs := flag.NewFlagSet("compare", flag.ContinueOnError)
	beforePath := fs.String("before", "", "earlier observation JSON")
	afterPath := fs.String("after", "", "later observation JSON")
	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 2
	}
	if fs.NArg() != 0 || *beforePath == "" || *afterPath == "" {
		fs.Usage()
		return 2
	}
	before, err := readObservations(*beforePath)
	if err == nil {
		var after []ObservationRecord
		after, err = readObservations(*afterPath)
		if err == nil {
			err = CompareObservations(before, after)
		}
	}
	if err != nil {
		var r *Refusal
		if errors.As(err, &r) {
			fmt.Fprintf(os.Stderr, "marketingcheck: %s %s: %s\n", r.Code, r.Name, r.Detail)
		} else {
			fmt.Fprintf(os.Stderr, "marketingcheck: observation-json compare: %s\n", err)
		}
		return 1
	}
	return 0
}
```

### Block C — `scripts/marketingcheck/observation_test.go` (create)

```go
package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"
	"time"
)

// observationPython is the plan's C11 boundary command, copied verbatim from
// w-marketing-capture-sprint-plan.md "C11 observations and mutation discipline".
const observationPython = `import pathlib, hashlib, json
p = pathlib.Path('.')
paths = sorted(set([q for pat in ('world/*.ail','go.mod','go.sum','design_docs/world-mission.md') for q in p.glob(pat) if q.is_file()] + [q for q in p.joinpath('tools/launchd').rglob('*') if q.is_file()]))
print(json.dumps([{'path':str(q),'sha256':hashlib.sha256(q.read_bytes()).hexdigest(),'lines':q.read_bytes().count(b'\n')} for q in paths],sort_keys=True))
`

var protectedFixture = map[string]string{
	"world/a.ail":                      "alpha one\nalpha two\n",
	"world/b.ail":                      "beta one\nbeta two\nbeta three\n",
	"world/sub/c.ail":                  "nested one\nnested two\n",
	"go.mod":                           "module fixture\n\ngo 1.26\n",
	"go.sum":                           "sum one\nsum two\n",
	"design_docs/world-mission.md":     "# charter\nrow one\n",
	"tools/launchd/x.sh":               "#!/bin/sh\necho x\n",
	"tools/launchd/testdata/deep/y.md": "deep one\ndeep two\n",
	"README.md":                        "readme one\nreadme two\n",
}

// protectedPaths is the sorted C11 subset of protectedFixture.
var protectedPaths = []string{"design_docs/world-mission.md", "go.mod", "go.sum", "tools/launchd/testdata/deep/y.md", "tools/launchd/x.sh", "world/a.ail", "world/b.ail"}

func protectedTree(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	for rel, body := range protectedFixture {
		p := filepath.Join(root, filepath.FromSlash(rel))
		pngOK(t, os.MkdirAll(filepath.Dir(p), 0o755))
		pngOK(t, os.WriteFile(p, []byte(body), 0o644))
	}
	return root
}

// flipProtectedByte turns "alpha one" into "blpha one": same length, same newline count.
func flipProtectedByte(t *testing.T, root string) {
	t.Helper()
	p := filepath.Join(root, "world", "a.ail")
	raw, err := os.ReadFile(p)
	pngOK(t, err)
	raw[0] = 'b'
	pngOK(t, os.WriteFile(p, raw, 0o644))
}

func mustObserve(t *testing.T, root string) []ObservationRecord {
	t.Helper()
	records, err := Observe(root)
	pngOK(t, err)
	return records
}

// pythonObserve runs observationPython in root and returns the path of its stdout file.
func pythonObserve(t *testing.T, root string) string {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	dir := t.TempDir()
	outPath, errPath := filepath.Join(dir, "python-observation.json"), filepath.Join(dir, "python-stderr.txt")
	out, err := os.Create(outPath)
	pngOK(t, err)
	defer out.Close()
	errFile, err := os.Create(errPath)
	pngOK(t, err)
	defer errFile.Close()
	cmd := exec.CommandContext(ctx, "python3", "-I", "-c", observationPython)
	cmd.Dir = root
	cmd.Stdout = out
	cmd.Stderr = errFile
	if err := cmd.Run(); err != nil {
		msg, _ := os.ReadFile(errPath)
		t.Fatalf("python observation: %v %s", err, msg)
	}
	return outPath
}

// runObservationCLI runs bin with file sinks and returns the exit code and stderr.
func runObservationCLI(t *testing.T, bin string, args ...string) (int, string) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	dir := t.TempDir()
	stdout, err := os.Create(filepath.Join(dir, "stdout"))
	pngOK(t, err)
	defer stdout.Close()
	stderrPath := filepath.Join(dir, "stderr")
	stderr, err := os.Create(stderrPath)
	pngOK(t, err)
	defer stderr.Close()
	cmd := exec.CommandContext(ctx, bin, args...)
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	rc := 0
	if err := cmd.Run(); err != nil {
		var ee *exec.ExitError
		if !errors.As(err, &ee) {
			t.Fatal(err)
		}
		rc = ee.ExitCode()
	}
	text, err := os.ReadFile(stderrPath)
	pngOK(t, err)
	return rc, string(text)
}

func TestMarketingProtectedObservationComparator(t *testing.T) {
	t.Run("set-coverage", func(t *testing.T) {
		got := mustObserve(t, protectedTree(t))
		var want []ObservationRecord
		for _, rel := range protectedPaths {
			sum := sha256.Sum256([]byte(protectedFixture[rel]))
			want = append(want, ObservationRecord{Path: rel, SHA256: hex.EncodeToString(sum[:]), Lines: strings.Count(protectedFixture[rel], "\n")})
		}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("Observe must cover exactly the C11 protected set: got %+v want %+v", got, want)
		}
	})
	t.Run("equal", func(t *testing.T) {
		root := protectedTree(t)
		if err := CompareObservations(mustObserve(t, root), mustObserve(t, root)); err != nil {
			t.Fatalf("unchanged copy must compare equal: %v", err)
		}
	})
	t.Run("same-line-count-change", func(t *testing.T) {
		root := protectedTree(t)
		before := mustObserve(t, root)
		flipProtectedByte(t, root)
		after := mustObserve(t, root)
		if before[5].Path != "world/a.ail" || after[5].Path != "world/a.ail" || before[5].Lines != after[5].Lines {
			t.Fatalf("fixture must be a genuine same-line-count change: %+v %+v", before[5], after[5])
		}
		err := CompareObservations(before, after)
		wantCode(t, err, "observation-digest", "same-line-count changed copy must red")
		var r *Refusal
		if !errors.As(err, &r) || r.Name != "world/a.ail" {
			t.Fatalf("same-line-count change must name world/a.ail: %v", err)
		}
	})
	t.Run("digest-equal-line-change", func(t *testing.T) {
		before := mustObserve(t, protectedTree(t))
		after := slices.Clone(before)
		after[0].Lines++
		wantCode(t, CompareObservations(before, after), "observation-lines", "line-count mismatch must refuse")
	})
	t.Run("deleted", func(t *testing.T) {
		root := protectedTree(t)
		before := mustObserve(t, root)
		pngOK(t, os.Remove(filepath.Join(root, "tools", "launchd", "testdata", "deep", "y.md")))
		wantCode(t, CompareObservations(before, mustObserve(t, root)), "observation-set", "deleted protected copy must red")
	})
	t.Run("added", func(t *testing.T) {
		root := protectedTree(t)
		before := mustObserve(t, root)
		pngOK(t, os.WriteFile(filepath.Join(root, "world", "z.ail"), []byte("zeta\n"), 0o644))
		wantCode(t, CompareObservations(before, mustObserve(t, root)), "observation-set", "added protected copy must red")
	})
	t.Run("duplicate", func(t *testing.T) {
		root := protectedTree(t)
		before := mustObserve(t, root)
		dup := append(slices.Clone(before), before[0])
		wantCode(t, CompareObservations(dup, mustObserve(t, root)), "observation-set", "duplicate path must red")
	})
	t.Run("empty", func(t *testing.T) {
		wantCode(t, CompareObservations(nil, nil), "observation-empty", "empty observation must not compare equal")
	})
	t.Run("python-parity", func(t *testing.T) {
		root := protectedTree(t)
		raw, err := os.ReadFile(pythonObserve(t, root))
		pngOK(t, err)
		var py []ObservationRecord
		pngOK(t, json.Unmarshal(raw, &py))
		slices.SortFunc(py, func(a, b ObservationRecord) int { return strings.Compare(a.Path, b.Path) })
		if got := mustObserve(t, root); !reflect.DeepEqual(py, got) {
			t.Fatalf("Go observer must equal the plan's Python boundary command: python %+v go %+v", py, got)
		}
	})
	t.Run("cli", func(t *testing.T) {
		dir := t.TempDir()
		bin := filepath.Join(dir, "marketingcheck")
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		if out, err := exec.CommandContext(ctx, "go", "build", "-o", bin, ".").CombinedOutput(); err != nil {
			t.Fatalf("CLI build: %v %s", err, out)
		}
		root := protectedTree(t)
		f1, f2 := filepath.Join(dir, "f1.json"), filepath.Join(dir, "f2.json")
		expect := func(label string, rc int, stderr string, wantRC int, wantText string) {
			t.Helper()
			if rc != wantRC || !strings.Contains(stderr, wantText) {
				t.Fatalf("%s: want rc %d stderr containing %q, got rc %d stderr %q", label, wantRC, wantText, rc, stderr)
			}
		}
		rc, se := runObservationCLI(t, bin, "observe", "--root", root, "--out", f1)
		expect("(i) observe --out", rc, se, 0, "")
		first, err := os.ReadFile(f1)
		pngOK(t, err)
		var parsed []ObservationRecord
		pngOK(t, json.Unmarshal(first, &parsed))
		if len(parsed) != len(protectedPaths) {
			t.Fatalf("(i) observe output must hold %d records, got %d", len(protectedPaths), len(parsed))
		}
		rc, se = runObservationCLI(t, bin, "observe", "--root", root, "--out", f1)
		expect("(ii) existing --out", rc, se, 1, "marketingcheck:")
		again, err := os.ReadFile(f1)
		pngOK(t, err)
		if string(again) != string(first) {
			t.Fatal("(ii) existing --out must remain unchanged")
		}
		rc, se = runObservationCLI(t, bin, "compare", "--before", f1, "--after", f1)
		expect("(iv) compare equal", rc, se, 0, "")
		rc, se = runObservationCLI(t, bin, "compare", "--before", pythonObserve(t, root), "--after", f1)
		expect("(vi) python output vs observe output", rc, se, 0, "")
		rc, se = runObservationCLI(t, bin, "compare", "--before", f1)
		expect("(v) compare missing --after", rc, se, 2, "")
		bad := filepath.Join(dir, "bad.json")
		pngOK(t, os.WriteFile(bad, []byte("{"), 0o600))
		rc, se = runObservationCLI(t, bin, "compare", "--before", f1, "--after", bad)
		expect("(vii) unparsable JSON", rc, se, 1, "marketingcheck: observation-json")
		rc, se = runObservationCLI(t, bin, "observe")
		expect("(viii) observe missing --root", rc, se, 2, "")
		flipProtectedByte(t, root)
		rc, se = runObservationCLI(t, bin, "observe", "--root", root, "--out", f2)
		expect("(iii) observe after flip", rc, se, 0, "")
		rc, se = runObservationCLI(t, bin, "compare", "--before", f1, "--after", f2)
		expect("(iii) compare after flip", rc, se, 1, "marketingcheck: observation-digest world/a.ail:")
	})
	t.Run("repo-read-only", func(t *testing.T) {
		first, second := mustObserve(t, "../.."), mustObserve(t, "../..")
		if err := CompareObservations(first, second); err != nil {
			t.Fatalf("real protected set observable without writes: %v", err)
		}
		paths := observationPaths(first)
		hasLaunchd := slices.ContainsFunc(paths, func(p string) bool { return strings.HasPrefix(p, "tools/launchd/") })
		if len(first) < 10 || !slices.Contains(paths, "go.mod") || !slices.Contains(paths, "design_docs/world-mission.md") || !hasLaunchd {
			t.Fatalf("real protected set observable without writes: %d records %v", len(first), paths)
		}
	})
}
```

### Tests and assertions (C11 · `TestMarketingProtectedObservationComparator`)

Every subtest that touches files builds its own fresh `protectedTree(t)` under `t.TempDir()`. No subtest depends on another subtest's writes; the `cli` subtest runs the Python heredoc itself through `pythonObserve`.

| Subtest | Setup | Assertion (exact failure text) |
|---|---|---|
| `set-coverage` | fresh tree (9 files: 7 protected; `world/sub/c.ail` and `README.md` are not observed) | `Observe` equals the 7 `protectedPaths` records, sorted, with SHA256/Lines computed from `protectedFixture`: "Observe must cover exactly the C11 protected set" |
| `equal` | `Observe` twice, no writes | nil: "unchanged copy must compare equal" |
| `same-line-count-change` | flip byte 0 of `world/a.ail` (`alpha one` → `blpha one`) | records[5] is `world/a.ail` with equal Lines both sides ("fixture must be a genuine same-line-count change"); `wantCode … "observation-digest", "same-line-count changed copy must red"`; `Refusal.Name == "world/a.ail"` |
| `digest-equal-line-change` | `after := slices.Clone(before); after[0].Lines++` (a controlled record, not a file) | `wantCode … "observation-lines", "line-count mismatch must refuse"` |
| `deleted` | remove `tools/launchd/testdata/deep/y.md` | `wantCode … "observation-set", "deleted protected copy must red"` |
| `added` | write `world/z.ail` | `wantCode … "observation-set", "added protected copy must red"` |
| `duplicate` | `append(slices.Clone(before), before[0])` vs a fresh `Observe` | `wantCode … "observation-set", "duplicate path must red"` |
| `empty` | `CompareObservations(nil, nil)` | `wantCode … "observation-empty", "empty observation must not compare equal"` |
| `python-parity` | `python3 -I -c observationPython`, `cmd.Dir = root`, stdout/stderr to `os.Create`d files, 30 s context | Python records sorted by Path `reflect.DeepEqual` `Observe(root)`: "Go observer must equal the plan's Python boundary command" |
| `cli` | `go build -o bin .` (60 s context); each run uses file sinks and a 60 s context | (i) `observe --root R --out F1` rc 0, F1 parses to 7 records; (ii) the same again: rc 1, stderr has `marketingcheck:`, F1 byte-identical; (iv) `compare F1 F1` rc 0; (vi) `compare --before <python output of R> --after F1` rc 0 (format interop); (v) `compare --before F1` rc 2; (vii) `--after` a file holding `{` gives rc 1 with stderr `marketingcheck: observation-json`; (viii) `observe` alone rc 2; (iii) flip byte, `observe --out F2` rc 0, then `compare F1 F2` gives rc 1 with stderr `marketingcheck: observation-digest world/a.ail:` |
| `repo-read-only` | `Observe("../..")` twice (the repository root; read-only) | compare nil; ≥10 records; has `go.mod`, `design_docs/world-mission.md` and a `tools/launchd/` path: "real protected set observable without writes" (12 at HEAD) |

Write order: write block C first and run the focused command. A build failure for missing `Observe`/`CompareObservations` is expected development red, **not** a mutant kill. Then write blocks A and B.

### Gates — exact commands and expected results

Run from the worktree root. Record rc, wall seconds and `uptime` for each.

```sh
gofmt -l scripts/marketingcheck                                                                      # expect: no output
go vet ./scripts/marketingcheck                                                                      # expect: rc 0, no output (compile fence for _test.go; go build ./... is not one)
go test -run '^TestMarketingProtectedObservationComparator$' ./scripts/marketingcheck -count=1 -timeout 120s -v   # expect: --- PASS top-level + 11 subtests PASS (prototype 0.5 s)
go test -run '^TestMarketingTwoPassInventory$' ./scripts/marketingcheck -count=1 -timeout 120s -v   # expect: PASS incl. cli/unknown, cli/no-command, cli/missing-flags, cli/help (dispatch regression)
go test -run '^TestProductionGoSurface$' ./host/store -count=1 -timeout 60s -v                      # expect: PASS, log "allow-listed reproducers=4" (base: host/cmd=165, reproducers=4)
go test -run '^TestSubprocessSinkGate$' ./host/verifygate -count=1 -timeout 120s -v                 # expect: PASS, "violations=0" (base: filesParsed=476)
go test -race ./scripts/marketingcheck -count=1 -timeout 240s                                       # expect: ok (prototype 33 s at load 1.5)
```

An empty selection ("no tests to run"), rc 0 with no top-level `--- PASS`, or a SKIP is not green. If the race or whole-package run is red **only** in `TestMarketingSupervisorFakeDaemonCleanup/shell-*`, that is the carried load flake. Re-run `go test -run '^TestMarketingSupervisorFakeDaemonCleanup$' ./scripts/marketingcheck -count=1 -timeout 120s -v` in isolation, report both results as a carried base finding, and never edit `supervisor_test.go`.

### Mutations — production `main.go` only, one at a time

Procedure for each mutant:
1. `cp scripts/marketingcheck/main.go "$SNAP/main.go.orig"`.
2. Make the one-line edit.
3. Run the command and require the exact red line.
4. `cp "$SNAP/main.go.orig" scripts/marketingcheck/main.go && cmp "$SNAP/main.go.orig" scripts/marketingcheck/main.go`.
5. Re-run the same command; it must PASS.

A build failure, panic, timeout, SKIP or empty selection is not a kill.

| Mutant | Exact edit in `CompareObservations` | Command | Must red with (observed in the prototype) |
|---|---|---|---|
| M7c-observation | `		if b.SHA256 != a.SHA256 {` → `		if false {` | `go test -run '^TestMarketingProtectedObservationComparator$/^same-line-count-change$' ./scripts/marketingcheck -count=1 -timeout 120s -v` | `same-line-count changed copy must red: want Code=observation-digest, got <nil>` |
| M7c-line-observation | `		if b.Lines != a.Lines {` → `		if false {` | `go test -run '^TestMarketingProtectedObservationComparator$/^digest-equal-line-change$' ./scripts/marketingcheck -count=1 -timeout 120s -v` | `line-count mismatch must refuse: want Code=observation-lines, got <nil>` |
| M7c1-set | `	if !slices.Equal(beforePaths, afterPaths) \|\| len(slices.Compact(slices.Clone(beforePaths))) != len(beforePaths) {` → `	if false {` | `go test -run '^TestMarketingProtectedObservationComparator$/^deleted$' ./scripts/marketingcheck -count=1 -timeout 120s -v` | `deleted protected copy must red: want Code=observation-set, got observation-digest tools/launchd/testdata/deep/y.md: sha256 … became ` (the path-keyed lookup returns the zero record; no panic) |

### C11 boundary observation (read-only)

Run the plan's Python heredoc (section "C11 observations and mutation discipline") at the worktree root before the first edit and after the last restore, with `OBS_OUT="$SNAP/obs-before.json"` / `"$SNAP/obs-after.json"`. `cmp` the two files; they must be identical. You may additionally run `go run ./scripts/marketingcheck compare --before "$SNAP/obs-before.json" --after "$SNAP/obs-after.json"` (rc 0) once the CLI exists. Never write a protected path.

### Executor directive notes

- **Narrowest gate list:** the seven gate commands, the three mutation commands with their restores, and the C11 observation before and after.
- **Allowed edits:** the two files above, plus scratch snapshots outside the checkout.
- **Reuse:** `refuse`/`Refusal` (png.go), and `wantCode`/`pngOK` (png_test.go). Add no third-party import.
- **Base measurement (iteration 305 planner):**
  - `TestMarketingTwoPassInventory` PASS, 1.3 s, load 8.2.
  - `go vet ./scripts/marketingcheck` rc 0.
  - `TestProductionGoSurface` PASS.
  - `TestSubprocessSinkGate` PASS.
  - `TestMarketingProtectedObservationComparator` is absent at base: `git grep` gives 0 matches, against a control of 1 for `TestMarketingTwoPassInventory`. Its first run must be the implementation, not a false-green empty selection.
  - Iteration 304 recorded a whole-package rc 1 from the carried `shell-hup` flake at load 8.2.

**no git write operations; snapshot changed files into $HOME/.ailang/state/world-iter<N>/snap/M7c1/ after each milestone** (`SNAP` above).
- Copy the changed allowlisted files with their relative paths.
- Include each gate's rc and wall time, the named PASS evidence, the mutation ledger and the observation records.
- Maximum 30 min. If incomplete, report it as incomplete; never relabel it PASS.

## M7c2a — Python capture primitives (Chrome version, HTTP readiness, six-shot pass) and Go `routes`

Design M7 (C7/C8/C9); 28 minutes; dependency: M7c1 **and M7c0** (amended iteration 305 per w-marketing-capture-orchestration §5); status **pending**.

Exact executor files:

- `scripts/marketing_capture.py`: add three modes, `chrome-version`, `ready` and `pass`. Existing `shot` and `daemon` behaviour is unchanged.
- `scripts/marketingcheck/main.go`: add a `routes` subcommand.
- `scripts/marketingcheck/capture_test.go`: create it, holding `TestMarketingCapturePassFakeChrome`.

The interfaces are prescribed; do not choose others.

- **`marketingcheck routes --manifest DEMO`** prints a JSON array of `{"name","width","height","scheme","route"}` in `Shots` order. Implementation: `func ShotRoutes(d DemoManifest) ([]RouteRow, error)`. It calls `BuildCaptureManifest(d, map[string]string{})` (`png.go:462`) for the route and takes width, height and scheme from `Shots` (`png.go:49-53`). Routes are never re-derived anywhere else (designer X6). The rows therefore carry the **M7c0** dimensions: packet 1440×3700, graph 1440×2100, provenance 1440×3000; light/dark 1440×1100; phone 390×844 (amended iteration 305 per w-marketing-capture-orchestration §5). Exit codes: 1 on `bad-manifest`, 2 on usage.
- **`marketing_capture.py chrome-version --chrome PATH`** runs `[PATH, '--version']` with stdin DEVNULL, output to a temp file, and `SHOT_TIMEOUT_S`. It parses `^(?:Google Chrome|Chromium) (\d+)\.\S+` from the stripped output. On success it prints `{"verdict":"ok","detail":"","version":"<full stripped string>"}`. Failures: `ChromeVersion` refusal for unparsable output or major `< 155` (module constant `MIN_CHROME_MAJOR = 155`, design §2 / V13), `NonzeroExit`, or `Timeout`.
- **`marketing_capture.py ready --base-url URL`** polls `GET URL/v1/health` and then `GET URL/v1/head` with `urllib.request.urlopen(…, timeout=10)` until both return 200, with a `STARTUP_TIMEOUT_S` overall deadline. The verdict is `ok` or `Timeout`.
- **`marketing_capture.py pass --routes FILE --base-url URL --chrome PATH --out-dir DIR`** handles the routes array from `routes`. The `DIR` must not exist; the mode creates it and refuses with `DestinationExists` otherwise. For each row it calls the existing `shot()` with `out=DIR/name`, `url=base+route`, and width, height and scheme from the row. After each shot it prints `progress: <i>/6 <name> <verdict>` to **stderr**. It stops at the first non-ok shot. The whole pass is bounded by `PASS_TIMEOUT_S`. The single stdout verdict is `{"verdict","detail","shots":[{"name","verdict"}…]}`.

Acceptance criteria and named assertions, in **`TestMarketingCapturePassFakeChrome`**. All of it is binary-free, socket-free except `ready`, and always-on.
- `routes`: build the CLI and run `routes` on `demoFixture()` (`png_test.go:198`) written by `putJSON`. Exactly six rows in `Shots` order. The routes equal `BuildCaptureManifest(demoFixture(), …)` row routes, for example `graph.png` → `/workbench?from=19&entry=19` and `provenance.png` → `/workbench?object=approval-invocation&payload=1`. Each row's width/height/scheme equals its `Shots` entry, and the test also asserts the literal M7c0 heights `packet.png` 3700, `graph.png` 2100, `provenance.png` 3000: *"routes must carry the M7c0 shot dimensions"* (amended iteration 305 per w-marketing-capture-orchestration §5).
- `pass-six`: the test writes `fixtureSet(t)` PNGs (`png_test.go:179`) to a fixture dir. The fake Chrome is a `fakeScript` (`supervisor_test.go:56`). It records argv per call, picks a fixture by `--window-size=` and by `preferredColorScheme=0` (dark) versus anything else (light), and copies it to `--screenshot=`. The `pass` verdict must be `ok`, with exactly six files. `CheckPass(dir, rects)` (`png.go:323`) must return nil. Each recorded argv carries the scheme flag matching its row: *"every shot explicitly pins correct scheme"*.
- `version`: the fake Chrome answers `--version` with each of `Google Chrome 156.0.8078.47`, `Google Chrome 155.0.8059.41`, `Google Chrome 154.0.1.2`, `garbage` and an exit-3 case. Expected verdicts: `ok` (full string recorded), `ok`, `ChromeVersion`, `ChromeVersion`, `NonzeroExit`. The two accepted strings are the real outputs recorded at V13 and V29, copied from the design rows that recorded them; the subtest keeps both (`156.0.8078.47` is also orchestration V41/V47) (amended iteration 305 per w-marketing-capture-orchestration §5).
- `ready`: a Go `httptest.NewServer` serves 503 on `/v1/health` for the first 2 calls and 200 afterwards, and 200 on `/v1/head`. Expect `ok`. A server that always returns 503 with `--startup-timeout 1` (unit-control override, same pattern as `--term-wait`) gives `Timeout`. Record `bind: operation not permitted` as UNINFORMATIVE UNDER SANDBOX.
- `fresh-dir`: an existing `--out-dir` gives `DestinationExists`.

Focused gates:

```sh
go test -run '^TestMarketingCapturePassFakeChrome$' ./scripts/marketingcheck -count=1 -timeout 180s -v
go test -run '^TestMarketingSupervisorFakeChromeSuccess$' ./scripts/marketingcheck -count=1 -timeout 90s -v
go vet ./scripts/marketingcheck
go test -race ./scripts/marketingcheck -count=1 -timeout 240s
```

| Mutant | Valid edit (production) | Test / assertion that must red |
|---|---|---|
| M7c-scheme | `marketing_capture.py` `shot()`: emit the `--blink-settings=preferredColorScheme=…` argument only when `args.scheme == 'light'` | `TestMarketingCapturePassFakeChrome/pass-six`: "every shot explicitly pins correct scheme". Both the argv check and `CheckPass` → `wrong-scheme workbench-dark.png` must red. |
| M7c2a-version-floor | `MIN_CHROME_MAJOR = 154` | `TestMarketingCapturePassFakeChrome/version`: "Chrome 154 must be refused" |
| M7c2a-routes | In `ShotRoutes`, use `"/workbench"` for every row instead of the `BuildCaptureManifest` route | `TestMarketingCapturePassFakeChrome/routes`: "routes must equal the capture-manifest routes" |

Executor notes:
- Allowed edits: only the three files.
- No real Chrome. No real daemon.
- `shot()` keeps its argv order (`marketing_capture.py:91-97`); only the new modes are added.
- Snapshot to `snap/M7c2a/`.
- Maximum 30 min.

## M7c2b — Capture shell end-to-end with fake Chrome and a fake daemon

Design M7 (C7/C8/C9); 28 minutes; dependency: M7c2a; status **pending**.

Exact executor files:

- `scripts/capture_marketing.sh`: replace the `capture refused` tail at lines 84-85 with the real flow.
- `scripts/marketing_capture.py`: `daemon` mode gains `--ready-url URL` as an alternative to `--ready-file`, using the `ready` poll from M7c2a. Exactly one of the two is required.
- `scripts/marketingcheck/capture_test.go`: add `TestMarketingCaptureFakeChromeSixShots`.

Shell flags. Keep the ten documented flags at `capture_marketing.sh:25-34` and add exactly three (X2/X3):
- `--worldd-bin PATH` (required; the `ailang-worldd` executable);
- `--bind 127.0.0.1:PORT` (required; loopback only, refused otherwise);
- `--marketingcheck-bin PATH` (required; prebuilt verifier).

All three are documented in `--help`.

Fixed flow, in order:
1. `chrome-version` → A.
2. Start `marketing_capture.py daemon --ready-url http://BIND --started-file … -- WORLDD serve --db DB --bind BIND --workbench-still --ailang-bin AILANG --workspace-root WS --tool-ailang-bin TOOL --proof-key-file KEY --proof-required-identity sameRef --proof-required-identity servesEntry`. The identities are hard-coded (designer X3/N7; `cmd/ailang-worldd/main.go:250-260`; QUICKSTART 1576-1577).
3. `MARKETINGCHECK routes --manifest M > OUT/routes.json`.
4. `pass` → `OUT/pass-a`.
5. `pass` → `OUT/pass-b` (only with `--verify-twice`; without it the run stops after pass A with rc 0 and no check).
6. `chrome-version` → B.
7. Refuse with rc 1 and `chrome version changed during capture` if A ≠ B.
8. Write `OUT/diagnostics.txt` with the lines `chrome-version: <A>` and `bind: <BIND>`.
9. Run `MARKETINGCHECK check --pass-a … --pass-b … --manifest M --rectangles R --out OUT/capture-manifest.json`. A check refusal leaves `pass-a` and `pass-b` on disk and exits 1 (designer X14).

The existing EXIT/INT/TERM/HUP traps stay, and `--unit-daemon` stays untouched.

**`TestMarketingCaptureFakeChromeSixShots`** sets up:
- prebuilt `marketingcheck`;
- the M7c2a fake Chrome;
- a fake `ailang-worldd`. This is a `python3` script that records argv to a file and serves `/v1/health` and `/v1/head` with 200 on the `--bind` address via `http.server`. The test obtains a free loopback port with `net.Listen("tcp","127.0.0.1:0")`, then closes it.

Subtests:
- `six-shots`: rc 0; `pass-a` and `pass-b` hold six names each and are byte-equal; `capture-manifest.json` passes `VerifyCaptureManifest`; `diagnostics.txt` contains the version.
- `serve-argv`: the recorded daemon argv contains `serve`, `--workbench-still`, both `--proof-required-identity` values, `--tool-ailang-bin` followed by exactly the path the test passed to the shell's `--tool-ailang-bin`, and **no** `commit`, `session`, `mint` or `approve`. The v0.54.15 tool pin is the operator's input. The shell hard-codes no AILANG version: `grep -cE 'v0\.[0-9]+\.[0-9]+' scripts/capture_marketing.sh` must print 0 (0 at base; positive control: the same grep on `design_docs/planned/w-marketing-capture-orchestration.md` is non-zero) (amended iteration 305 per w-marketing-capture-orchestration §5).
- `version-changed`: the fake Chrome answers 155 then 156 (a counter file). rc 1, with stderr naming the change.
- `old-chrome`: 154 → rc 1 before any shot.
- `cleanup`: after rc, the fake daemon PID (recorded) is dead (`assertDead`, `supervisor_test.go:212`).
- `scheme`: re-asserts the M7c2a scheme argv across the full chain.

Bounds: 120 s per shell run, against a ≈ few-second unloaded run.

Focused gates:

```sh
go test -run '^TestMarketingCaptureFakeChromeSixShots$' ./scripts/marketingcheck -count=1 -timeout 300s -v
go test -run '^TestMarketingSupervisorFakeDaemonCleanup$' ./scripts/marketingcheck -count=1 -timeout 120s -v   # trap regression, isolated
bash -n scripts/capture_marketing.sh
go test -race ./scripts/marketingcheck -count=1 -timeout 300s
```

| Mutant | Valid edit (production) | Test / assertion that must red |
|---|---|---|
| M7c2b-version-changed | In `capture_marketing.sh`, make the A≠B comparison never refuse | `TestMarketingCaptureFakeChromeSixShots/version-changed`: "a Chrome version change during capture must refuse" |
| M7c-scheme (re-drill) | Same Python edit as M7c2a | `TestMarketingCaptureFakeChromeSixShots/scheme`: "every shot explicitly pins correct scheme" |
| M7c2b-identities | Drop one `--proof-required-identity` from the serve argv in the shell | `TestMarketingCaptureFakeChromeSixShots/serve-argv`: "both seed identities must be configured" |

Executor notes:
- Rerun the shell-hup subtest in isolation after any trap change (designer N5).
- Fake success is never C8 product evidence.
- Snapshot to `snap/M7c2b/`.
- Maximum 30 min. If the shell plus test exceeds ~250 lines, stop and report rather than shortcutting a subtest.

## M7c3 — Frozen public inventory guard (real store, genuine extra commit)

Design M7 (C8/C11); 25 minutes; dependency: M7c2b; status **pending**.

Exact executor files:

- `scripts/marketingcheck/main.go`: add `inventory --addr URL [--out FILE]` and `inventory-compare --before FILE --after FILE`.
- `host/daemon/marketing_inventory_test.go`: create it in `package daemon`, so it can use the existing in-package genuine-commit seams.

**Inventory definition** (designer X7: there is no list-all-objects route; `host/daemon/daemon.go:1033-1059`). The inventory is `{"schema":"world/demo/capture-inventory/v1","head":<GET /v1/head body>,"entries":[{"entryHash","transitionRef"}…],"objects":[{"ref","sha256"}…]}`:
- `entries` comes from paging `GET /v1/log?from=N&limit=500` from 0 until a page returns zero items. `items` is the response key; `entryHash` and `transitionRef` are fields of `logEntryResponse`, `host/daemon/handlers.go:65-73`.
- `objects` holds, for each distinct `transitionRef`, the SHA256 of the raw `GET /v1/objects/{ref}` body.
- Each request uses a 10 s client timeout.

Functions: `func ReadInventory(base string) (Inventory, error)` and `func CompareInventory(before, after Inventory) error`. The refusal code is `inventory-mutated` for a differing head, entry list or object set, and `inventory-empty` for an empty `head` or no entries. The CLI exit and stderr format match `check`.

**`TestMarketingCaptureInventoryMutation`** (host/daemon) runs in this order:
1. `d := newHandlerDaemon(t)` (`handlers_test.go:21`).
2. `c := workbenchLiveREST(t, d, []int64{0,1,2})` (`workbench_live_test.go:74`), which makes genuine authenticated REST commits.
3. `srv := httptest.NewServer(d.Handler())`.
4. Build `../../scripts/marketingcheck` with `go build -o` (60 s context).
5. `inventory --addr srv.URL --out before.json`.
6. Again to `same.json`, then `inventory-compare` → rc 0: *"no-write capture inventory must compare equal"*.
7. One extra genuine commit: `testCommit(c.NextWorld, 3, "workbench-live")` posted with `authHeader(t, d)` through `requestRecorderAuth`, the exact pattern at `workbench_live_test.go:77-84`.
8. `inventory --out after.json`, then `inventory-compare --before before.json --after after.json` → rc 1, with stderr containing `marketingcheck: inventory-mutated`: *"extra commit invalidates capture"*.

All subprocess sinks are files. Indexes must be contiguous: `handleLogRange` stops at the first missing index (`handlers.go:656-657`). Otherwise unchanged (amended iteration 305 per w-marketing-capture-orchestration §5). The seeded demo store's log is contiguous (indexes 0–9, orchestration V40), so stop-at-gap does not truncate the real capture inventory.

Focused gates:

```sh
go test -run '^TestMarketingCaptureInventoryMutation$' ./host/daemon -count=1 -timeout 180s -v
go test -run '^TestWorkbenchLiveCursorCrossesGap$' ./host/daemon -count=1 -timeout 90s -v
go test -run '^TestProductionGoSurface$' ./host/store -count=1 -timeout 60s -v
go vet ./scripts/marketingcheck ./host/daemon
```

| Mutant | Valid edit (production `main.go`) | Test / assertion that must red |
|---|---|---|
| M7c-inventory | `CompareInventory` returns nil after the empty check (neuter the equality comparison) | `TestMarketingCaptureInventoryMutation`: "extra commit invalidates capture" |
| M7c3-paging | `ReadInventory` stops after the first page and requests `limit=1` | `TestMarketingCaptureInventoryMutation`: "after.json must list all 4 committed entries" (the executor adds this assertion on the parsed `after.json`; the head comparison alone would still red, so the count is the named kill) |

Executor notes:
- The controller may wire `inventory` before and after both passes into `capture_marketing.sh` as a one-line follow-up. That is **not** in M7c3's file list, to keep it within one run, and is recorded as a carried integration step for M7c4.
- Snapshot to `snap/M7c3/`.

## M7c4 — Real Chrome drill (test written by the executor; RUN by the controller)

Design M7 (C7/C8/C9); 25 minutes; dependency: M7c3; status **pending**.

Exact executor files:

- `cmd/ailang-worldd/marketing_chrome_test.go`: create `TestMarketingChromeSixShots` in `package main`, **plus a local rig helper in this same file** (amended iteration 305 per w-marketing-capture-orchestration §5). The V11 seams live in this package: `testMintEnv` (`session_test.go:29`), `newDemoApprovalRig` (`demo_test.go:561`) and `demoM5cRealSeed` (`demo_test.go:1727`). This is designer X1, option (a).
  - `demoM5cRealSeed` cannot be reused. It returns `(worldReader, demoPublicManifest, string)` (`return r, m, rig.token`, `demo_test.go:1830`), exposing neither the db path nor the daemon. `demoSeRig` (`demo_test.go:548`) holds no daemon; the daemon closes only at test cleanup. A second daemon on the same db fails on the writer lock (`host/daemon/daemon.go:744-751`; orchestration §4.3, V45).
  - The helper therefore composes `newDemoApprovalRig`'s steps with the daemon handle retained, runs `demo seed`, closes the seed daemon, and returns the db path, manifest path and pins. The capture shell's `serve` can then open the db.
  - The reference is the iteration-305 probe `design_docs/verification/world-iter305/iter305-seeded-still-chrome-probe_test.go.txt`: rig at `:64` ("newDemoApprovalRig's body (demo_test.go:561), keeping the daemon handle"), `demo seed` at `:139`, daemon close at `:143-144`.
  - **`cmd/ailang-worldd/demo_test.go` is NOT edited.
- `scripts/capture_marketing.sh`: wire `marketingcheck inventory` before pass A and after pass B, plus `inventory-compare`. This is the M7c3 carried step: one call each, refused with rc 1 on mismatch.
- `scripts/marketingcheck/capture_test.go`: add subtest `inventory` to `TestMarketingCaptureFakeChromeSixShots`. Extend its fake daemon to also serve `/v1/log?from=&limit=` (`{"items":[…]}`) and `/v1/objects/{ref}`. The `inventory` subtest makes the fake daemon return a different `/v1/head` after pass A; rc must be 1.

The test's preconditions, in order:
- `WORLD_CHROME` unset → `t.Skip` (ordinary suite).
- `WORLD_CHROME` set but not executable, or `chrome-version` not `ok` → `t.Fatal`.
- Both pins are required as described in the plan environment; a missing pin → `t.Fatal`.
- `WORLD_MARKETING_RECTS` unset or unloadable by `LoadRectangles` semantics → `t.Fatal`, never guessed rectangles.

The test's run:
1. Seed with the local rig helper (above), then close the seed daemon (amended iteration 305 per w-marketing-capture-orchestration §5).
2. `go build` `ailang-worldd` and `marketingcheck` into `t.TempDir()`.
3. Run `scripts/capture_marketing.sh` with all 13 flags and `--verify-twice` on the rig's db. The outer bound is 660 s.
4. Require rc 0, `capture-manifest.json` verified, `diagnostics.txt` with an equal pass-A/pass-B Chrome version, and an inventory-equal record.

Chrome rule: design §2 / V29. A bounded `--version` runs before A and after B; the major version must be ≥155; the two strings must be identical; the string is recorded in sibling diagnostics. "Exact 155" is no longer required. The stale "PD1 unresolved visibility must red" clause is dropped, because PD1 was resolved in iteration 266 (designer N9).

Executor acceptance (unattended):

```sh
export AILANG_BIN=$HOME/.pinned-ailang/ailang WORLD_TOOL_AILANG_BIN=$HOME/.pinned-ailang-tools/v0.54.15/ailang   # (amended iteration 305 per w-marketing-capture-orchestration §4.2)
go test -run '^TestMarketingChromeSixShots$' ./cmd/ailang-worldd -count=1 -timeout 120s -v                 # WORLD_CHROME unset: must print SKIP (recorded, not acceptance)
WORLD_CHROME=/nonexistent go test -run '^TestMarketingChromeSixShots$' ./cmd/ailang-worldd -count=1 -timeout 120s -v   # must FAIL: set-but-invalid
go vet ./cmd/ailang-worldd
bash -n scripts/capture_marketing.sh
go test -run '^TestMarketingCaptureFakeChromeSixShots$' ./scripts/marketingcheck -count=1 -timeout 300s -v
```

Controller rehearsal. This is never executor acceptance and never unattended C8 evidence. It needs measured rectangles first:

```sh
WORLD_CHROME="/Applications/Google Chrome.app/Contents/MacOS/Google Chrome" WORLD_MARKETING_RECTS=<controller-measured file> \
  go test -run '^TestMarketingChromeSixShots$' ./cmd/ailang-worldd -count=1 -timeout 660s -v
```

| Mutant | Valid edit (production) | Test / assertion that must red |
|---|---|---|
| M7c4-inventory-wiring | Remove the `inventory-compare` call from `capture_marketing.sh` | `TestMarketingCaptureFakeChromeSixShots/inventory`: "capture must compare inventories". The executor adds this subtest: the fake daemon changes `/v1/head` between calls, and rc must be 1. |
| M7c4-rects-required | In the test's precondition helper, which lives in the test file, make a missing `WORLD_MARKETING_RECTS` a Skip | Not a production mutant. It is listed as a **review check**, not a kill. The guard lives in test code by design, because the drill is a test. |

**M7c4 completion semantics** (amended iteration 305 per w-marketing-capture-orchestration §5). M7c4 is COMPLETE when all of the following hold:
- (a) its PR lands with PR CI green;
- (b) with `WORLD_CHROME` unset the drill prints SKIP, and with `WORLD_CHROME=/nonexistent` it FAILs;
- (c) an unset or unloadable `WORLD_MARKETING_RECTS` FAILs with Chrome set (a review check, since that guard is test code);
- (d) `TestMarketingCaptureFakeChromeSixShots/inventory` passes, and mutant M7c4-inventory-wiring is killed;
- (e) an independent judge PASSes.

Completion does **not** include a real-Chrome pass, and it never marks C8/C9 (or C8′/C12/C14) met. M8a depends on M7c4 COMPLETE, not on R1. The real-Chrome run is R1 (doc §5): a controller step that measures all-six rectangles over the M7c0 geometry, writes `WORLD_MARKETING_RECTS`, runs the full `capture_marketing.sh` → `marketing_capture.py daemon` → `pass` → `shot` path twice (C14), and records the provisional `png.go:27` thresholds measured on the real renders as a V-row.

Executor notes:
- A real run without supplied rectangles is recorded as DEFERRED/SKIP.
- Measuring the rectangles and running the rehearsal are controller-owned (JSON `controller_owned`).
- Landing M7c4 never marks C8/C9 met (designer N10).
- Snapshot to `snap/M7c4/`.

## M8a — Reproduction docs, captions and release manifest schema

Design M8; 25 minutes; dependency: M7c4 **COMPLETE** (the completion semantics in §M7c4; not R1) (amended iteration 305 per w-marketing-capture-orchestration §5); status **pending**.

Exact executor files:

- `docs/QUICKSTART.md`
- `README.md`
- `website/docs/getting-started/marketing-capture.md`
- `website/static/img/v1/captions.md`
- `docs/marketing-capture-manifest.schema.json`
- `cmd/ailang-worldd/demo_test.go`
- `scripts/marketingcheck/docs_test.go`
- `scripts/marketingcheck/png_test.go` (amended iteration 305 per w-marketing-capture-orchestration §5). Its `TestMarketingTwoPassInventory/cli/unknown` case passes `docs` and expects rc 2 (`png_test.go:800`). M8a's `docs` subcommand changes that, so the case moves to another unknown word (e.g. `bogus`). Also move the `TestMarketingPNGContent` threshold boundary literals (below).
- `scripts/marketingcheck/png.go` (amended iteration 305 per w-marketing-capture-orchestration §5): the threshold constants at `png.go:27-35`.
- `scripts/marketingcheck/main.go` (planner addition, iteration 305): the `docs` subcommand used by the release block (`go run ./scripts/marketingcheck docs --root . --assets … --require-assets`). `main.go` is the only allowed home for production Go under `scripts/` (`host/store/context_roots_test.go:204-206`).

**Environment** (amended iteration 305 per w-marketing-capture-orchestration §4.2): `AILANG_BIN=$HOME/.pinned-ailang/ailang`, `WORLD_TOOL_AILANG_BIN=$HOME/.pinned-ailang-tools/v0.54.15/ailang`. Every documented reproduction recipe names the v0.54.15 tool pin, never v0.52.1.

**Captions** (amended iteration 305 per w-marketing-capture-orchestration §5): say "static capture of recorded history" and disclose the tall framing (packet 1440×3700, graph 1440×2100, provenance 1440×3000). The provenance caption uses the narrowed claim (orchestration §4.1): "the stored invocation record, opened from the approval entry, with its leading fields legible; the complete reference set is the accompanying `why` JSON". Also disclose the cosmetic risks in orchestration §8 (scrollbars, phone overflow, provenance horizontal clip, truncated graph labels).

**Threshold ownership (orchestration §5 R1 / V57)** (amended iteration 305 per w-marketing-capture-orchestration §5):
- R1 (controller) measures the provisional content thresholds (`png.go:27` "provisional until M7c measures real renders": MinColors 32, ChromaticDelta 12, ChromaticFraction 0.01, LightMinY 160, DarkMaxY 96, MinRectW 64, MinRectH 16, MinRectStdDev 4.0) on the real renders and records a V-row. **The values are not the executor's to choose.**
- M8a applies that V-row to `png.go` and moves the `TestMarketingPNGContent` boundary literals with them. The literals are at `png_test.go:359-360` (MinRectW/H: 63×16 and 64×15 rejected), `:396` (MinColors), `:438-441` (LightMinY/DarkMaxY: 160/96 accepted, 159/97 rejected) and `:444-455` (MinRectStdDev: delta 7 rejected, 8 accepted).
- **Gate:** `go test ./scripts/marketingcheck -count=1 -timeout 240s -v` is green, with each boundary arm still failing one step past the new value.
- **Flagged for the controller (planner, iteration 305):** the doc says M8a "does not depend on R1" yet applies R1's V-row. If no R1 threshold V-row exists when M8a is dispatched, the threshold step is reported DEFERRED (not invented and not skipped silently). The controller then holds it for a follow-up run; the docs part of M8a proceeds.

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

Final validation run, no new product scope. Re-run the narrow cross-boundary tests and mandatory repo gates. Go commands run with `AILANG_BIN=$HOME/.pinned-ailang/ailang` and `WORLD_TOOL_AILANG_BIN=$HOME/.pinned-ailang-tools/v0.54.15/ailang` (amended iteration 305 per w-marketing-capture-orchestration §4.2). Use observation snapshots, never restore protected files. Report every gate rc/wall time and all unresolved product criteria. Attended capture/promotion is outside milestones.

Acceptance criteria and named assertions:

- **C11 · `TestMarketingProtectedObservationComparator`:** Re-drill copied-input comparator; final actual protected inventory equals initial sprint observation in both digest and newline count. Scoped read-only diff must be empty for all protected paths.

Exact focused and repository gate commands (environment/helper above applies):

```sh
go test -run '^TestMarketingProtectedObservationComparator$' ./scripts/marketingcheck -count=1 -timeout 90s -v
go test -run '^TestDemoSeedEndToEnd$' ./cmd/ailang-worldd -count=1 -timeout 120s -v
go test -run '^TestDemoSeedWhyNamedLinks$' ./cmd/ailang-worldd -count=1 -timeout 120s -v
go test -run '^TestWorkbenchPacketEvidenceOnlyThroughRowOneOhFive$' ./host/daemon -count=1 -timeout 45s -v
go test -run '^TestMarketingCaptureFakeChromeSixShots$' ./scripts/marketingcheck -count=1 -timeout 300s -v   # (amended iteration 305 per w-marketing-capture-orchestration §5): spans two <=120 s shell runs
go test -run '^TestMarketingShotGeometry$' ./scripts/marketingcheck -count=1 -timeout 60s -v                 # M7c0 geometry regression (iteration 305)
go test -run '^TestMarketingPNGStructure$' ./scripts/marketingcheck -count=1 -timeout 90s -v
AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_ail.sh
go build ./... && go test ./... -count=1 -timeout 600s
go test -race ./host/daemon -count=1 -timeout 480s
go test -race ./host/coordinator ./host/workbench ./host/evidence ./scripts/marketingcheck -count=1 -timeout 480s
```

Mutations; each must reach this assertion:

| Mutant | Valid edit | Test / assertion that must red |
|---|---|---|
| M8b-C11-redrill | Re-drill M7c1 line-count guard mutant (M7c-line-observation) on copied observations | `TestMarketingProtectedObservationComparator`: line-count mismatch must refuse |

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

Unattended implementation and test gates can be completed without release confirmations. Release C8 pixels/byte equality, actual image/manifest promotion and C10 actual asset resolution remain controller-owned attended-release checks. PD1 and any M7 fallback residual must appear in the boundary report; no claim of full showcase completion or clause movement follows from tool/test success. JSON statuses (amended iteration 305 per w-marketing-capture-orchestration §5): M1a–M7b `completed` (PRs #264, #266, #267, #273, #275, #277, #283, #289, #295, #353, #354, #359, #360, #361; orchestration §2, V33); M7c0–M8b `pending`. C8′, C12, C13 and C14 (orchestration §7) are met only by the controller/attended steps R0, R1, P0, P3 and A1, never by an executor milestone.
