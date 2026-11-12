set -uo pipefail
V=$HOME/.ailang/state/world-iter389/gates/verdicts_4c2cac17.txt; : > "$V"
leg(){ echo "$1 rc=$2" | tee -a "$V"; }
export AILANG_BIN=$HOME/.pinned-ailang/ailang
unset WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER WORLD_PASSKEY_RECORD_DIR
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter389-a4-s1; cd "$W"
L=$HOME/.ailang/state/world-iter389/gates; mkdir -p "$L"; MODE=final
s0=$(date +%s)
./scripts/verify_ail.sh 2>&1 | tee "$L/verify_ail.log"; leg verify_ail ${PIPESTATUS[0]}
out=$(gofmt -l cmd host); echo "$out" | tee "$L/gofmt.log"; leg gofmt $([ -z "$out" ] && echo 0 || echo 1)
go build ./... 2>&1 | tee "$L/build.log"; leg build ${PIPESTATUS[0]}
go vet ./... 2>&1 | tee "$L/vet.log"; leg vet ${PIPESTATUS[0]}
go test -run '^$' ./... 2>&1 | tee "$L/fence.log"; leg fence ${PIPESTATUS[0]}

go test -count=1 -timeout 30m ./host/approval/... ./host/coordinator/ ./host/surfacepage/ ./host/daemon/ 2>&1 | tee "$L/pkgs.log"; leg pkgs ${PIPESTATUS[0]}
grep -E "^(FAIL|--- FAIL)" "$L/pkgs.log" | tee -a "$V"
go test -count=1 -timeout 30m ./host/broker/ 2>&1 | tee "$L/broker.log"; leg broker ${PIPESTATUS[0]}   # reported, not gating (§6)
if grep -E '^(FAIL|--- FAIL)' "$L/broker.log"; then echo "BROKER_RED (re-run each red alone; known-list reds measure the repo)"; fi
go test -race -count=1 ./host/approval/ ./host/coordinator/ 2>&1 | tee "$L/race.log"
go test -race -count=1 ./host/broker/ -run '^(TestDecideApproval|TestPublishApproval|TestDecisionRefusal|TestReservedApprovalScope|TestApprovalContextIdentity)' 2>&1 | tee -a "$L/race.log"
go test -race -count=1 ./host/daemon/ -run '^(TestDecision|TestDecide|TestEpisodeFinish|TestApprov|TestRequest|TestDaemonCredentialWriterSites|TestSurfaceAppRoutes)' 2>&1 | tee -a "$L/race.log"
grep -E "^(FAIL|--- FAIL)|WARNING: DATA RACE" "$L/race.log" | tee -a "$V"; leg race_grep_hits $(grep -cE "^(FAIL|--- FAIL)|WARNING: DATA RACE" "$L/race.log")
NEW_A="TestDecisionMessageLayout TestDecisionMessageRefusesBadInput"
NEW_B="TestDecisionRefusalD1 TestDecisionRefusalD2 TestDecisionRefusalD3 TestDecisionRefusalD4a TestDecisionRefusalD4b TestDecisionRefusalD5 TestDecisionRefusalD6 TestDecisionRefusalD7 TestDecisionRefusalD8 TestDecisionRefusalD9 TestDecisionRefusalD10 TestDecisionRefusalOrder TestDecisionRefusalMustLive TestReservedApprovalScope TestDecideApprovalRefusesASecondDecision TestDecideApprovalFirstDecisionMustLive TestPublishApprovalRefusesAKeyDecision TestPublishApprovalTypedStampMustLive TestDecideApprovalCallerCensus"
NEW_C="TestDiscoverRequestPacketsCompleteness"
NEW_D="TestDecisionRoutesRegistered TestDecisionSubjectJSON TestDecidePOSTGuards TestDecidePOSTLawArms TestDecideTypedLands TestDecideTypedProofBindsTheVerb TestDecidePOSTStaleDigests TestDecideReplayAfterRestart TestDecideRefusesReservedScopesTyped TestDecideRefusesASharedRequest TestDecideDiscoveryFindsAHiddenSecondPacket TestDecideDiscoveryRefusesAnOverflow TestDecideDiscoveryFindsTheOnlyPacketBeyondThePrefix"
if [ "$MODE" = base ]; then NEW_A=""; NEW_B=""; NEW_C=""; NEW_D=""; fi
alt() { local x="$*"; x="${x// /|}"; echo "${x:+|$x}"; }
go test -v -count=1 ./host/approval/... -run "^(TestApprovalImportsOnlyTheStandardLibrary|TestApprovalTestIsTestOnly$(alt $NEW_A))\$" 2>&1 | tee "$L/inv.log"
go test -v -count=1 ./host/broker/ -run "^(TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential|TestQuickstartSection10ProfilesLoad|TestQuickstartProtectedContract|TestApprovalContextIdentity$(alt $NEW_B))\$" 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./host/coordinator/ -run "^(TestValidatePacketBytesCheckNeedsThePacketLabel$(alt $NEW_C))\$" 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./host/daemon/ -run "^(TestDaemonDependencyAllowlist|TestFrozenV1RouteTableMatchesMux|TestSurfaceAppRoutes|TestApprovalRoutesRegistered|TestReadCtxCancelledAfterHandler|TestMCPMountMethodSource|TestDaemonCredentialWriterSites|TestQuickstartDecisionPacketSection|TestDaemonConstructorCallerInventory|TestEpisodeFinishOperatorProofExpiryBoundary$(alt $NEW_D))\$" 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./host/boundary/ -run '^(TestWorldBoundaryDependencyAllowlist|TestBareNetHTTPExemptionIsPerGroup|TestEveryCommandDirectoryIsAProtectedGroup)$' 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./host/store/ -run '^(TestProductionContextRoots|TestProductionGoSurface|TestNoNewDeadlineFreeStoreReads|TestReadGettersHonorContext|TestCredentialsMethodsRequireADeadline)$' 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./host/verifygate/ -run '^(TestNoParallelWriteForkPackagesOutsideVerifygate|TestSubprocessSinkGate)$' 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./cmd/ailang-worldd/ -run '^(TestEveryVerbHasHelp|TestQuickstartAttendedStepsBoxMatchesSection9|TestQuickstartSection9FlagsMatchTheCLI|TestQuickstartSection10FlagsMatchTheCLI|TestQuickstartSection11FlagsMatchTheCLI|TestQuickstartServiceSectionMatchesTheCLIHelp)$' 2>&1 | tee -a "$L/inv.log"
for t in TestApprovalImportsOnlyTheStandardLibrary TestApprovalTestIsTestOnly TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential TestQuickstartSection10ProfilesLoad TestQuickstartProtectedContract TestApprovalContextIdentity \
  TestValidatePacketBytesCheckNeedsThePacketLabel TestDaemonDependencyAllowlist TestFrozenV1RouteTableMatchesMux TestSurfaceAppRoutes TestApprovalRoutesRegistered TestReadCtxCancelledAfterHandler TestMCPMountMethodSource \
  TestDaemonCredentialWriterSites TestQuickstartDecisionPacketSection TestDaemonConstructorCallerInventory TestEpisodeFinishOperatorProofExpiryBoundary \
  TestWorldBoundaryDependencyAllowlist TestBareNetHTTPExemptionIsPerGroup TestEveryCommandDirectoryIsAProtectedGroup \
  TestProductionContextRoots TestProductionGoSurface TestNoNewDeadlineFreeStoreReads TestReadGettersHonorContext TestCredentialsMethodsRequireADeadline \
  TestNoParallelWriteForkPackagesOutsideVerifygate TestSubprocessSinkGate TestEveryVerbHasHelp TestQuickstartAttendedStepsBoxMatchesSection9 \
  TestQuickstartSection9FlagsMatchTheCLI TestQuickstartSection10FlagsMatchTheCLI TestQuickstartSection11FlagsMatchTheCLI TestQuickstartServiceSectionMatchesTheCLIHelp \
  $NEW_A $NEW_B $NEW_C $NEW_D; do
  grep -qE -- "^--- PASS: $t " "$L/inv.log" || echo "NOTPASS $t" | tee -a "$V"; done
grep -E "^\s*--- (SKIP|FAIL)" "$L/inv.log" | tee -a "$V"
echo "named PASS: $(grep -cE '^--- PASS' "$L/inv.log")" | tee -a "$L/inv.log"
echo "GATES_DONE for 4c2cac17 $(( $(date +%s) - s0 ))s" | tee -a "$V"
