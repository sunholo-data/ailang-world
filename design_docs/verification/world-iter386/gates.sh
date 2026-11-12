set -euo pipefail
export AILANG_BIN=$HOME/.pinned-ailang/ailang
unset WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER WORLD_PASSKEY_RECORD_DIR
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter386-a3a-s3; cd "$W"
L=$HOME/.ailang/state/world-iter386/gates; mkdir -p "$L"; MODE=final
s0=$(date +%s)
./scripts/verify_ail.sh 2>&1 | tee "$L/verify_ail.log"
out=$(gofmt -l cmd host); echo "$out" | tee "$L/gofmt.log"; test -z "$out"
go build ./... 2>&1 | tee "$L/build.log"
go vet ./... 2>&1 | tee "$L/vet.log"
go test -run '^$' ./... 2>&1 | tee "$L/fence.log"
if grep -E '^(FAIL|--- FAIL)' "$L/fence.log"; then exit 1; fi
go test -count=1 -timeout 20m ./host/surfacepage/ ./host/store/ ./host/authority/ ./host/approval/... ./cmd/ailang-worldd/ ./host/daemon/ 2>&1 | tee "$L/pkgs.log"
if grep -E '^(FAIL|--- FAIL)' "$L/pkgs.log"; then exit 1; fi
go test -race -count=1 ./host/surfacepage/ 2>&1 | tee "$L/race.log"
go test -race -count=1 ./host/daemon/ -run '^(TestApprov|TestRequest|TestDaemonCredentialWriterSites|TestSurfaceAppRoutes)' 2>&1 | tee -a "$L/race.log"
go test -race -count=1 ./cmd/ailang-worldd/ -run '^(TestApprove|TestServeSurfaceOrigin|TestSurfaceApprove|TestSurfaceSSHSIG|TestEveryVerbHasHelp)' 2>&1 | tee -a "$L/race.log"
if grep -E '^(FAIL|--- FAIL)|WARNING: DATA RACE' "$L/race.log"; then exit 1; fi
NEW_SP="TestApprovePageRendersEveryCanonicalField TestApprovePageFieldTagControls"
NEW_D="TestRequestCardPasskeyReadError TestApprovePOSTRecordsAVerifiedSSHSIG TestApprovePOSTSSHSIGViaTailnetOriginWithNoPasskey TestApprovePageRoutes"
NEW_C="TestSurfaceSSHSIGThenAwaitMints"
if [ "$MODE" = base ]; then NEW_SP=""; NEW_D=""; NEW_C=""; fi
alt() { local x="$*"; x="${x// /|}"; echo "${x:+|$x}"; }
go test -v -count=1 ./host/broker/ -run '^(TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential|TestQuickstartSection10ProfilesLoad|TestQuickstartProtectedContract)$' 2>&1 | tee "$L/inv.log"
go test -v -count=1 ./host/daemon/ -run "^(TestDaemonDependencyAllowlist|TestFrozenV1RouteTableMatchesMux|TestSurfaceAppRoutes|TestDaemonConstructorCallerInventory|TestQuickstartDecisionPacketSection|TestPasskeyPageRoutes|TestReadCtxCancelledAfterHandler|TestMCPMountMethodSource|TestDaemonCredentialWriterSites|TestApprovalRoutesRegistered|TestApprovalRoutesRefuseForeignHost|TestRequestsJSON|TestRequestCardJSON|TestApprovePOSTRefusalArms|TestApprovePOSTRecordsAVerifiedPasskey|TestApprovePOSTCrossOrigin|TestApprovalGuards$(alt $NEW_D))\$" 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./host/boundary/ -run '^(TestWorldBoundaryDependencyAllowlist|TestBareNetHTTPExemptionIsPerGroup|TestEveryCommandDirectoryIsAProtectedGroup)$' 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./host/store/ -run '^(TestProductionContextRoots|TestProductionGoSurface|TestNoNewDeadlineFreeStoreReads|TestReadGettersHonorContext|TestCredentialsMethodsRequireADeadline|TestCredentialsSchemaMatchesCanonicalDDL|TestCredentialsSchema7SurfaceApprovals|TestRecordSurfaceApprovalGuards|TestSurfaceApprovalPrunedWithItsRequest)$' 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./host/verifygate/ -run '^(TestNoParallelWriteForkPackagesOutsideVerifygate|TestSubprocessSinkGate)$' 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./host/approval/... -run '^(TestApprovalImportsOnlyTheStandardLibrary|TestApprovalTestIsTestOnly)$' 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./host/authority/ -run '^(TestJudgeAssertionIsTheCLIJudge|TestJudgeSSHSignatureIsTheCLIJudge)$' 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./cmd/ailang-worldd/ -run "^(TestEveryVerbHasHelp|TestServeProofFlags|TestQuickstartAttendedStepsBoxMatchesSection9|TestQuickstartSection9FlagsMatchTheCLI|TestQuickstartSection10FlagsMatchTheCLI|TestQuickstartSection11FlagsMatchTheCLI|TestQuickstartServiceSectionMatchesTheCLIHelp|TestSurfaceApproveThenAwaitMints|TestServeSurfaceOriginFlag$(alt $NEW_C))\$" 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./host/surfacepage/ -run "^(TestSurfaceAppPure|TestApproveModelSummaries|TestApproveModelCard$(alt $NEW_SP))\$" 2>&1 | tee -a "$L/inv.log"
for t in TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential TestQuickstartSection10ProfilesLoad TestQuickstartProtectedContract \
  TestDaemonDependencyAllowlist TestFrozenV1RouteTableMatchesMux TestSurfaceAppRoutes TestDaemonConstructorCallerInventory TestQuickstartDecisionPacketSection TestPasskeyPageRoutes TestReadCtxCancelledAfterHandler TestMCPMountMethodSource TestDaemonCredentialWriterSites \
  TestApprovalRoutesRegistered TestApprovalRoutesRefuseForeignHost TestRequestsJSON TestRequestCardJSON TestApprovePOSTRefusalArms TestApprovePOSTRecordsAVerifiedPasskey TestApprovePOSTCrossOrigin TestApprovalGuards \
  TestWorldBoundaryDependencyAllowlist TestBareNetHTTPExemptionIsPerGroup TestEveryCommandDirectoryIsAProtectedGroup \
  TestProductionContextRoots TestProductionGoSurface TestNoNewDeadlineFreeStoreReads TestReadGettersHonorContext TestCredentialsMethodsRequireADeadline TestCredentialsSchemaMatchesCanonicalDDL \
  TestCredentialsSchema7SurfaceApprovals TestRecordSurfaceApprovalGuards TestSurfaceApprovalPrunedWithItsRequest \
  TestNoParallelWriteForkPackagesOutsideVerifygate TestSubprocessSinkGate TestApprovalImportsOnlyTheStandardLibrary TestApprovalTestIsTestOnly \
  TestJudgeAssertionIsTheCLIJudge TestJudgeSSHSignatureIsTheCLIJudge \
  TestEveryVerbHasHelp TestServeProofFlags TestQuickstartAttendedStepsBoxMatchesSection9 TestQuickstartSection9FlagsMatchTheCLI TestQuickstartSection10FlagsMatchTheCLI \
  TestQuickstartSection11FlagsMatchTheCLI TestQuickstartServiceSectionMatchesTheCLIHelp TestSurfaceApproveThenAwaitMints TestServeSurfaceOriginFlag TestSurfaceAppPure TestApproveModelSummaries TestApproveModelCard $NEW_SP $NEW_D $NEW_C; do
  grep -qE -- "^--- PASS: $t " "$L/inv.log" || { echo "$t did not PASS"; exit 1; }; done
if grep -E '^\s*--- (SKIP|FAIL)' "$L/inv.log"; then exit 1; fi
go test -v -count=1 ./cmd/ailang-worldd/ -run '^TestApproveAwait' 2>&1 | tee "$L/aw.log"
test "$(grep -cE '^--- PASS: TestApproveAwait' "$L/aw.log")" = 9
if grep -E '^\s*--- (SKIP|FAIL)' "$L/aw.log"; then exit 1; fi
# the live-Chrome leg (this machine only: CI has no Chrome step)
CHROME_RUN='^(TestPasskeyPageInChrome|TestPasskeyPageAbortsAPendingCeremony|TestApprovePage|TestApproveCard|TestNeedsYou)'
WORLD_CHROME="/Applications/Google Chrome.app/Contents/MacOS/Google Chrome" WORLD_REQUIRE_PASSKEY_BROWSER=1 go test -v -count=1 -timeout 600s ./cmd/ailang-worldd/ -run "$CHROME_RUN" 2>&1 | tee "$L/chrome.log"
if grep -E '^\s*--- (SKIP|FAIL)|^FAIL' "$L/chrome.log"; then exit 1; fi
WANT_CHROME="TestPasskeyPageInChrome TestPasskeyPageAbortsAPendingCeremony"
if [ "$MODE" != base ]; then WANT_CHROME="$WANT_CHROME TestApprovePageInChrome TestApprovePageRefusesASwappedRequest TestApprovePageSSHSIGInChrome TestApproveCardNetworkDeadlines TestNeedsYouListsPendingRequests"; fi
for t in $WANT_CHROME; do grep -qE -- "^--- PASS: $t " "$L/chrome.log" || { echo "$t did not PASS in Chrome"; exit 1; }; done
echo "PRE_COMMIT_GATES_GREEN $(( $(date +%s) - s0 ))s" | tee "$L/pre_commit.ok"
