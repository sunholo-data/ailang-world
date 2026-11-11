bash <<'EOF'
set -euo pipefail
export AILANG_BIN=$HOME/.pinned-ailang/ailang
unset WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER WORLD_PASSKEY_RECORD_DIR
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter384-a3a; cd "$W"
L=$HOME/.ailang/state/world-iter384/gates; mkdir -p "$L"
s0=$(date +%s)
./scripts/verify_ail.sh 2>&1 | tee "$L/verify_ail.log"                      # no .ail change: exactly green
out=$(gofmt -l cmd host); echo "$out" | tee "$L/gofmt.log"; test -z "$out"
go build ./... 2>&1 | tee "$L/build.log"
go vet ./... 2>&1 | tee "$L/vet.log"
go test -run '^$' ./... 2>&1 | tee "$L/fence.log"                           # compile fence for every _test.go
go test -count=1 -timeout 20m ./host/store/ ./host/authority/ ./host/approval/... ./cmd/ailang-worldd/ ./host/daemon/ 2>&1 | tee "$L/pkgs.log"
go test -race -count=1 ./host/store/ -run 'Credentials|SurfaceApproval|MintRequest|Passkey|Approver|ContextRoots' 2>&1 | tee "$L/race.log"
go test -race -count=1 ./host/authority/ 2>&1 | tee -a "$L/race.log"
go test -race -count=1 ./cmd/ailang-worldd/ -run '^(TestApprove|TestPasskeyEnroll|TestEveryVerbHasHelp)' 2>&1 | tee -a "$L/race.log"
# the repo-wide inventories (design §7), the schema pins, the reflective, help and quickstart tests and the new tests, by name, -v:
go test -v -count=1 ./host/broker/ -run '^(TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential|TestQuickstartSection10ProfilesLoad|TestQuickstartProtectedContract)$' 2>&1 | tee "$L/inv.log"
go test -v -count=1 ./host/daemon/ -run '^(TestDaemonDependencyAllowlist|TestFrozenV1RouteTableMatchesMux|TestSurfaceAppRoutes|TestDaemonConstructorCallerInventory|TestQuickstartDecisionPacketSection|TestPasskeyPageRoutes)$' 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./host/boundary/ -run '^(TestWorldBoundaryDependencyAllowlist|TestBareNetHTTPExemptionIsPerGroup|TestEveryCommandDirectoryIsAProtectedGroup)$' 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./host/store/ -run '^(TestProductionContextRoots|TestProductionGoSurface|TestNoNewDeadlineFreeStoreReads|TestReadGettersHonorContext|TestCredentialsMethodsRequireADeadline|TestCredentialsSchemaMatchesCanonicalDDL|TestCredentialsSchema7SurfaceApprovals|TestRecordSurfaceApprovalGuards|TestSurfaceApprovalPrunedWithItsRequest)$' 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./host/verifygate/ -run '^(TestNoParallelWriteForkPackagesOutsideVerifygate|TestSubprocessSinkGate)$' 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./host/approval/... -run '^(TestApprovalImportsOnlyTheStandardLibrary|TestApprovalTestIsTestOnly)$' 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./cmd/ailang-worldd/ -run '^(TestEveryVerbHasHelp|TestQuickstartAttendedStepsBoxMatchesSection9|TestQuickstartSection9FlagsMatchTheCLI|TestQuickstartSection10FlagsMatchTheCLI|TestQuickstartSection11FlagsMatchTheCLI|TestQuickstartServiceSectionMatchesTheCLIHelp)$' 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./host/authority/ -run '^(TestJudgeAssertionIsTheCLIJudge|TestJudgeSSHSignatureIsTheCLIJudge)$' 2>&1 | tee -a "$L/inv.log"
for t in TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential TestQuickstartSection10ProfilesLoad TestQuickstartProtectedContract \
  TestDaemonDependencyAllowlist TestFrozenV1RouteTableMatchesMux TestSurfaceAppRoutes TestDaemonConstructorCallerInventory TestQuickstartDecisionPacketSection TestPasskeyPageRoutes \
  TestWorldBoundaryDependencyAllowlist TestBareNetHTTPExemptionIsPerGroup TestEveryCommandDirectoryIsAProtectedGroup \
  TestProductionContextRoots TestProductionGoSurface TestNoNewDeadlineFreeStoreReads TestReadGettersHonorContext TestCredentialsMethodsRequireADeadline TestCredentialsSchemaMatchesCanonicalDDL \
  TestNoParallelWriteForkPackagesOutsideVerifygate TestSubprocessSinkGate TestApprovalImportsOnlyTheStandardLibrary TestApprovalTestIsTestOnly \
  TestEveryVerbHasHelp TestQuickstartAttendedStepsBoxMatchesSection9 TestQuickstartSection9FlagsMatchTheCLI TestQuickstartSection10FlagsMatchTheCLI \
  TestQuickstartSection11FlagsMatchTheCLI TestQuickstartServiceSectionMatchesTheCLIHelp \
  TestCredentialsSchema7SurfaceApprovals TestRecordSurfaceApprovalGuards TestSurfaceApprovalPrunedWithItsRequest \
  TestJudgeAssertionIsTheCLIJudge TestJudgeSSHSignatureIsTheCLIJudge; do
  grep -qE -- "^--- PASS: $t " "$L/inv.log" || { echo "$t did not PASS"; exit 1; }; done
if grep -E '^\s*--- (SKIP|FAIL)' "$L/inv.log"; then exit 1; fi
go test -v -count=1 ./cmd/ailang-worldd/ -run '^TestApproveAwait' 2>&1 | tee "$L/aw.log"
test "$(grep -cE '^--- PASS: TestApproveAwait' "$L/aw.log")" = 8
if grep -E '^\s*--- (SKIP|FAIL)' "$L/aw.log"; then exit 1; fi
echo "PRE_COMMIT_GATES_GREEN $(( $(date +%s) - s0 ))s" | tee "$L/pre_commit.ok"
EOF
