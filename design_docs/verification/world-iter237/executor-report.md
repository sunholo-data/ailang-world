# Row 136 executor report (iteration 237)

Worktree ~/.ailang/state/world-wt/iter236-row136, branch sprint/row136-mcp-typed-errors. Not pushed. Tree clean.

## Commits
- M1 e480ddc feat(row136): M1 — MCP refusal line labelled mcp, refs tested
- M2 7a0ede1 feat(row136): M2 — mcp: invocation ids and coordinator:mcp provenance
- M3 18a4e9c docs(row136): M3 — ... upstream ailang#1602 (real URL, no placeholder; grep ISSUE-URL = 0)

M1 = proto_m1.diff + AC1.1 D3 id-in-cause assertion + AC1.1a (2-arg InvocationID, a2a: prefix at M1). M2 = the rest of proto_final.diff (tree equals it) with AC1.1 prefix flipped to mcp:.

## Gates (each commit; gate.sh, logs in row136-executor/{m1,m2,m3}_*.log)
| commit | verify_ail | go vet | go test -race ./... | tool-binary e2e |
| M1 | rc=0 | rc=0 | rc=1: 2607 RUN, 2551 PASS, 1 FAIL, 55 SKIP (flake, below) | rc=0, 4/4 named tests --- PASS |
| M2 | rc=0 | rc=0 | rc=0: 2612 RUN, 2556 PASS, 0 FAIL, 56 SKIP | rc=0, 5/5 --- PASS (adds TestSeToolsSurfaceTaggedInLog) |
| M3 | rc=0 | rc=0 | rc=0: 2612 RUN, 2556 PASS, 0 FAIL, 56 SKIP | rc=0, 5/5 --- PASS |
Tool e2e names: TestSeToolsMCPPostEffectFailureIsLogged, TestSeToolsA2APostEffectFailureIsMapped, TestSeToolsMCPEndToEnd, TestCLIAgainstSeToolsDaemon, (M2+) TestSeToolsSurfaceTaggedInLog; WORLD_TOOL_AILANG_BIN = v0.52.1 tool pin, AILANG_BIN = v0.41.0 gate pin.

## Flake
M1 full run: host/broker TestRunBoundedHeadTailTimeoutKeepsPartialOutput (handlers_capture_test.go:122, partial stdout = ""). Package untouched by row 136. `go test -race -count=1 ./host/broker/` alone 3/3 ok (~129s each). Same class as the planner's P4 broker flake. Not fixed.

## Sweeps
- After M1: grep "a2a refusal|mcp invoke|mcp tools" website docs: 7 hits = 3 pi `/mcp tools` slash lines (coding-tools.md:225, connecting.md:79, QUICKSTART.md:574), A2A lines (operating-the-daemon.md:99, QUICKSTART.md:93, troubleshooting.md:187 which names both), mcp-and-a2a.md:174. No `mcp invoke`.
- After M2: grep 'coordinator:a2a|a2a:<' website docs README.md: 10 hits, each names both surfaces or A2A explicitly (coding-tools:262, provenance:45, provenance-walks:79, sessions-and-episodes:16, cli.md:37,254, AGENTS.md:165, README:240 and a2a-only example lines). Call{ grep = 6 hits, 5 literals carry Surface:.

## Mutants (final tree, applied via perl, restored, git status clean after each)
All KILLED, 0 survived, 0 invalid builds (mutants.txt in row136-executor/):
- MUT-MCP-LOG-NOCAUSE: ci TestMCPRefusalLineNamesSurfaceAndCause KILLED; tool TestSeToolsMCPPostEffectFailureIsLogged KILLED (:826)
- MUT-MCP-LOG-LABEL: ci KILLED (mcp_conformance_test.go:403); tool KILLED (:821)
- MUT-A2A-LABEL: TestA2ADispatchWire KILLED (a2a_wire_test.go:212)
- MUT-UNRECORDED-ORDER: TestA2ADispatch KILLED (projection_test.go:847)
- MUT-MCP-SURFACE-a (Call+log): AC2.1 KILLED (:167); AC1.1 KILLED (:821)
- MUT-MCP-SURFACE-b (Call literal only): AC2.1 KILLED; AC1.1 KILLED (:826, D3 assertion)
- MUT-SURFACE-DEFAULT: TestDispatchRefuses/R1_empty_surface KILLED
- MUT-WRITTENBY-CONST: TestSeToolsSurfaceTaggedInLog KILLED (:897); TestPlanWrittenByNamesSurface KILLED in CI mode (:208); object Provenance variant KILLED (:212)
- MUT-WHY-PREFIX: TestWhyResolvesMCPInvocationID KILLED
- MUT-REST-ACCEPTS-MCP: TestCommitInvocationReceiptIdentity/foreign_namespaces_are_refused KILLED
- MUT-A2A-ID-FORMAT: TestNewRefusesMissingSeamsAndCaps + TestReconcileRefusesDamagedRecord KILLED
- MUT-CLIWALK (regex left at a2a:): TestCLIAgainstSeToolsDaemon KILLED

## Not done / notes
- Mutant regexes used `|` and each run reported nonzero RUN counts. Nothing pushed; no PRs touched. se-tools e2e still SKIPs in CI (plan D1, controller's follow-up row).
