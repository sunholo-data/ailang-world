bash <<'EOF'
set -euo pipefail
export AILANG_BIN=$HOME/.pinned-ailang/ailang
export PATH="/opt/homebrew/bin:$PATH"     # node v26 for the XB leg
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter382-row201-a2s4; cd "$W"
L=$HOME/.ailang/state/world-iter382/gates; mkdir -p "$L"
CHROME="/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
./scripts/verify_ail.sh 2>&1 | tee "$L/verify_ail.log"                      # no .ail change: exactly green
out=$(gofmt -l cmd host); echo "$out" | tee "$L/gofmt.log"; test -z "$out"
go build ./... 2>&1 | tee "$L/build.log"
go vet ./... 2>&1 | tee "$L/vet.log"
go test -run '^$' ./... 2>&1 | tee "$L/fence.log"                           # compile fence for every _test.go
go test -count=1 -timeout 20m ./host/surfacepage/... ./host/approval/... ./cmd/ailang-worldd/ ./host/daemon/ 2>&1 | tee "$L/pkgs.log"
go test -race -count=1 ./host/surfacepage/ 2>&1 | tee "$L/race.log"
go test -race -count=1 ./host/daemon/ -run '^(TestPasskeyPageRoutes|TestSurfaceAppRoutes|TestSurfaceScriptRoute)$' 2>&1 | tee -a "$L/race.log"
go test -race -count=1 ./cmd/ailang-worldd/ -run '^(TestPasskey|TestApproveByPasskey|TestApprovePasskey|TestApprover|TestEveryVerbHasHelp)' 2>&1 | tee -a "$L/race.log"
# the repo-wide inventories (design §7), the reflective, help and quickstart tests, by name, -v:
go test -v -count=1 ./host/broker/ -run '^TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential$' 2>&1 | tee "$L/inv.log"
go test -v -count=1 ./host/daemon/ -run '^(TestDaemonDependencyAllowlist|TestFrozenV1RouteTableMatchesMux|TestSurfaceAppRoutes|TestDaemonConstructorCallerInventory|TestQuickstartDecisionPacketSection|TestPasskeyPageRoutes)$' 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./host/boundary/ -run '^(TestWorldBoundaryDependencyAllowlist|TestBareNetHTTPExemptionIsPerGroup|TestEveryCommandDirectoryIsAProtectedGroup)$' 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./host/store/ -run '^(TestProductionContextRoots|TestProductionGoSurface|TestNoNewDeadlineFreeStoreReads|TestReadGettersHonorContext|TestCredentialsMethodsRequireADeadline)$' 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./host/verifygate/ -run '^(TestNoParallelWriteForkPackagesOutsideVerifygate|TestSubprocessSinkGate)$' 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./host/approval/... -run '^(TestApprovalImportsOnlyTheStandardLibrary|TestApprovalTestIsTestOnly|TestPasskeyFixturesAreTheGeneratedBytes)$' 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./cmd/ailang-worldd/ -run '^(TestEveryVerbHasHelp|TestQuickstartAttendedStepsBoxMatchesSection9|TestQuickstartSection9FlagsMatchTheCLI|TestQuickstartSection10FlagsMatchTheCLI|TestQuickstartSection11FlagsMatchTheCLI|TestQuickstartServiceSectionMatchesTheCLIHelp|TestPasskeyFixtureGeneratorRefuses)$' 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./host/broker/ -run '^(TestQuickstartSection10ProfilesLoad|TestQuickstartProtectedContract)$' 2>&1 | tee -a "$L/inv.log"
for t in TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential TestDaemonDependencyAllowlist TestFrozenV1RouteTableMatchesMux \
  TestSurfaceAppRoutes TestDaemonConstructorCallerInventory TestQuickstartDecisionPacketSection TestPasskeyPageRoutes \
  TestWorldBoundaryDependencyAllowlist TestBareNetHTTPExemptionIsPerGroup TestEveryCommandDirectoryIsAProtectedGroup \
  TestProductionContextRoots TestProductionGoSurface TestNoNewDeadlineFreeStoreReads TestReadGettersHonorContext \
  TestCredentialsMethodsRequireADeadline TestNoParallelWriteForkPackagesOutsideVerifygate TestSubprocessSinkGate TestApprovalImportsOnlyTheStandardLibrary \
  TestApprovalTestIsTestOnly TestPasskeyFixturesAreTheGeneratedBytes TestEveryVerbHasHelp TestQuickstartAttendedStepsBoxMatchesSection9 \
  TestQuickstartSection9FlagsMatchTheCLI TestQuickstartSection10FlagsMatchTheCLI TestQuickstartSection11FlagsMatchTheCLI \
  TestQuickstartServiceSectionMatchesTheCLIHelp TestPasskeyFixtureGeneratorRefuses TestQuickstartSection10ProfilesLoad TestQuickstartProtectedContract; do
  grep -qE -- "^--- PASS: $t " "$L/inv.log" || { echo "$t did not PASS"; exit 1; }; done
if grep -E '^\s*--- (SKIP|FAIL)' "$L/inv.log"; then exit 1; fi
# XB1 and XB2, required, three consecutive runs (Chrome headless, temp profile, virtual authenticator):
for i in 1 2 3; do
  s=$(date +%s)
  WORLD_CHROME="$CHROME" WORLD_REQUIRE_PASSKEY_BROWSER=1 go test -v -count=1 ./cmd/ailang-worldd/ -run '^(TestPasskeyPageInChrome|TestPasskeyPageAbortsAPendingCeremony)$' 2>&1 | tee "$L/xb_$i.log"
  echo "XB run $i: $(( $(date +%s) - s )) s" | tee -a "$L/xb_times.log"
  for t in TestPasskeyPageInChrome TestPasskeyPageAbortsAPendingCeremony; do grep -qE -- "^--- PASS: $t " "$L/xb_$i.log" || { echo "$t did not PASS (run $i)"; exit 1; }; done
  if grep -E '^\s*--- (SKIP|FAIL)' "$L/xb_$i.log"; then exit 1; fi
  test "$(grep -cE '^    --- PASS: TestPasskeyPage' "$L/xb_$i.log")" = 7                   # 3 XB1 + 4 XB2 subtests
done
echo PRE_COMMIT_GATES_GREEN | tee "$L/pre_commit.ok"
EOF
