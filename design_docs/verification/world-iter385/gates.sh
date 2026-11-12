#!/usr/bin/env bash
set -euo pipefail
export AILANG_BIN=$HOME/.pinned-ailang/ailang
unset WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER WORLD_PASSKEY_RECORD_DIR
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter385-a3a-s2; cd "$W"
L=$HOME/.ailang/state/world-iter385/gates; mkdir -p "$L"
s0=$(date +%s)
./scripts/verify_ail.sh 2>&1 | tee "$L/verify_ail.log"                      # no .ail change: exactly green
out=$(gofmt -l cmd host); echo "$out" | tee "$L/gofmt.log"; test -z "$out"
go build ./... 2>&1 | tee "$L/build.log"
go vet ./... 2>&1 | tee "$L/vet.log"
go test -run '^$' ./... 2>&1 | tee "$L/fence.log"                           # compile fence for every _test.go
go test -count=1 -timeout 20m ./host/surfacepage/ ./host/store/ ./host/authority/ ./host/approval/... ./cmd/ailang-worldd/ ./host/daemon/ 2>&1 | tee "$L/pkgs.log"
if grep -E '^(FAIL|--- FAIL)' "$L/pkgs.log"; then exit 1; fi
go test -race -count=1 ./host/surfacepage/ 2>&1 | tee "$L/race.log"
go test -race -count=1 ./host/daemon/ -run '^(TestApprov|TestRequest|TestDaemonCredentialWriterSites|TestSurfaceAppRoutes)' 2>&1 | tee -a "$L/race.log"
go test -race -count=1 ./cmd/ailang-worldd/ -run '^(TestApprove|TestServeSurfaceOrigin|TestSurfaceApprove|TestEveryVerbHasHelp)' 2>&1 | tee -a "$L/race.log"
if grep -E '^(FAIL|--- FAIL)|WARNING: DATA RACE' "$L/race.log"; then exit 1; fi
NEW_SP="TestApproveModelSummaries TestApproveModelRefusesWhatDoesNotDecode TestApproveModelCard TestApproveModelBodyGrammar TestApproveModelAdmission"
NEW_D="TestApprovalGuards TestApprovalRoutesRegistered TestApprovalRoutesRefuseForeignHost TestRequestsJSON TestRequestCardJSON TestApprovePOSTRefusalArms TestApprovePOSTRecordsAVerifiedPasskey TestApprovePOSTCrossOrigin TestApprovePOSTReplay TestApprovePOSTWriterBusy TestApprovePOSTNoSessionNoCookie TestDaemonCredentialWriterSites"
NEW_C="TestApproveAwaitJudgesTheApproverAtNow TestServeSurfaceOriginFlag TestSurfaceApproveThenAwaitMints"
alt() { local x="$*"; echo "${x// /|}"; }
# the repo-wide inventories (design §7), the schema pins, the reflective, help and quickstart tests and the 20 new tests, by name, -v:
go test -v -count=1 ./host/broker/ -run '^(TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential|TestQuickstartSection10ProfilesLoad|TestQuickstartProtectedContract)$' 2>&1 | tee "$L/inv.log"
go test -v -count=1 ./host/daemon/ -run "^(TestDaemonDependencyAllowlist|TestFrozenV1RouteTableMatchesMux|TestSurfaceAppRoutes|TestDaemonConstructorCallerInventory|TestQuickstartDecisionPacketSection|TestPasskeyPageRoutes|TestReadCtxCancelledAfterHandler|$(alt $NEW_D))\$" 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./host/boundary/ -run '^(TestWorldBoundaryDependencyAllowlist|TestBareNetHTTPExemptionIsPerGroup|TestEveryCommandDirectoryIsAProtectedGroup)$' 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./host/store/ -run '^(TestProductionContextRoots|TestProductionGoSurface|TestNoNewDeadlineFreeStoreReads|TestReadGettersHonorContext|TestCredentialsMethodsRequireADeadline|TestCredentialsSchemaMatchesCanonicalDDL|TestCredentialsSchema7SurfaceApprovals|TestRecordSurfaceApprovalGuards|TestSurfaceApprovalPrunedWithItsRequest)$' 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./host/verifygate/ -run '^(TestNoParallelWriteForkPackagesOutsideVerifygate|TestSubprocessSinkGate)$' 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./host/approval/... -run '^(TestApprovalImportsOnlyTheStandardLibrary|TestApprovalTestIsTestOnly)$' 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./host/authority/ -run '^(TestJudgeAssertionIsTheCLIJudge|TestJudgeSSHSignatureIsTheCLIJudge)$' 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./cmd/ailang-worldd/ -run "^(TestEveryVerbHasHelp|TestServeProofFlags|TestQuickstartAttendedStepsBoxMatchesSection9|TestQuickstartSection9FlagsMatchTheCLI|TestQuickstartSection10FlagsMatchTheCLI|TestQuickstartSection11FlagsMatchTheCLI|TestQuickstartServiceSectionMatchesTheCLIHelp|$(alt $NEW_C))\$" 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./host/surfacepage/ -run "^($(alt $NEW_SP))\$" 2>&1 | tee -a "$L/inv.log"
for t in TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential TestQuickstartSection10ProfilesLoad TestQuickstartProtectedContract \
  TestDaemonDependencyAllowlist TestFrozenV1RouteTableMatchesMux TestSurfaceAppRoutes TestDaemonConstructorCallerInventory TestQuickstartDecisionPacketSection TestPasskeyPageRoutes TestReadCtxCancelledAfterHandler \
  TestWorldBoundaryDependencyAllowlist TestBareNetHTTPExemptionIsPerGroup TestEveryCommandDirectoryIsAProtectedGroup \
  TestProductionContextRoots TestProductionGoSurface TestNoNewDeadlineFreeStoreReads TestReadGettersHonorContext TestCredentialsMethodsRequireADeadline TestCredentialsSchemaMatchesCanonicalDDL \
  TestCredentialsSchema7SurfaceApprovals TestRecordSurfaceApprovalGuards TestSurfaceApprovalPrunedWithItsRequest \
  TestNoParallelWriteForkPackagesOutsideVerifygate TestSubprocessSinkGate TestApprovalImportsOnlyTheStandardLibrary TestApprovalTestIsTestOnly \
  TestJudgeAssertionIsTheCLIJudge TestJudgeSSHSignatureIsTheCLIJudge \
  TestEveryVerbHasHelp TestServeProofFlags TestQuickstartAttendedStepsBoxMatchesSection9 TestQuickstartSection9FlagsMatchTheCLI TestQuickstartSection10FlagsMatchTheCLI \
  TestQuickstartSection11FlagsMatchTheCLI TestQuickstartServiceSectionMatchesTheCLIHelp $NEW_SP $NEW_D $NEW_C; do
  grep -qE -- "^--- PASS: $t " "$L/inv.log" || { echo "$t did not PASS"; exit 1; }; done
if grep -E '^\s*--- (SKIP|FAIL)' "$L/inv.log"; then exit 1; fi
go test -v -count=1 ./cmd/ailang-worldd/ -run '^TestApproveAwait' 2>&1 | tee "$L/aw.log"
test "$(grep -cE '^--- PASS: TestApproveAwait' "$L/aw.log")" = 9                       # sprint 1's 8 + AW-NOW
if grep -E '^\s*--- (SKIP|FAIL)' "$L/aw.log"; then exit 1; fi
echo "PRE_COMMIT_GATES_GREEN $(( $(date +%s) - s0 ))s" | tee "$L/pre_commit.ok"
