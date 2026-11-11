#!/usr/bin/env bash
cd /Users/voightkampff/dev/sunholo-data/.wt-world-iter379-row201-a2 || exit 2
set -euo pipefail
export AILANG_BIN=$HOME/.pinned-ailang/ailang
L=$HOME/.ailang/state/world-iter379/gates; mkdir -p "$L"
./scripts/verify_ail.sh 2>&1 | tee "$L/verify_ail.log"                     # no .ail change: exactly green
test -z "$(gofmt -l cmd/ailang-worldd host/approval host/sshsign)"
go build ./... && go vet ./... 2>&1 | tee "$L/vet.log"
go test -run '^$' ./... 2>&1 | tee "$L/fence.log"                          # compile fence for every _test.go
go test ./cmd/ailang-worldd/ ./host/approval/... ./host/sshsign/ -count=1 2>&1 | tee "$L/pkgs.log"
go test -race ./host/approval/... ./host/sshsign/ -count=1 2>&1 | tee "$L/race.log"
go test -race ./cmd/ailang-worldd/ -run '^(TestApproverRevoke|TestApproverEnroll)' -count=1 2>&1 | tee "$L/race-cmd.log"
# the repo-wide inventories (design §7, iteration 378 retro (a)), by name, -v:
go test -v -count=1 ./host/broker/ -run '^TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential$' 2>&1 | tee "$L/inv.log"
go test -v -count=1 ./host/daemon/ -run '^(TestDaemonDependencyAllowlist|TestFrozenV1RouteTableMatchesMux|TestSurfaceAppRoutes|TestDaemonConstructorCallerInventory)$' 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./host/boundary/ -run '^(TestWorldBoundaryDependencyAllowlist|TestBareNetHTTPExemptionIsPerGroup|TestEveryCommandDirectoryIsAProtectedGroup)$' 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./host/store/ -run '^(TestProductionContextRoots|TestProductionGoSurface|TestNoNewDeadlineFreeStoreReads|TestReadGettersHonorContext|TestCredentialsMethodsRequireADeadline)$' 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./host/verifygate/ -run '^TestNoParallelWriteForkPackagesOutsideVerifygate$' 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./host/approval/... -run '^(TestApprovalImportsOnlyTheStandardLibrary|TestApprovalTestIsTestOnly)$' 2>&1 | tee -a "$L/inv.log"
for t in TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential TestDaemonDependencyAllowlist TestFrozenV1RouteTableMatchesMux \
  TestSurfaceAppRoutes TestDaemonConstructorCallerInventory TestWorldBoundaryDependencyAllowlist TestBareNetHTTPExemptionIsPerGroup \
  TestEveryCommandDirectoryIsAProtectedGroup TestProductionContextRoots TestProductionGoSurface TestNoNewDeadlineFreeStoreReads \
  TestReadGettersHonorContext TestCredentialsMethodsRequireADeadline TestNoParallelWriteForkPackagesOutsideVerifygate \
  TestApprovalImportsOnlyTheStandardLibrary TestApprovalTestIsTestOnly; do
  grep -qE -- "^--- PASS: $t " "$L/inv.log" || { echo "$t did not PASS"; exit 1; }; done
if grep -E '^\s*--- (SKIP|FAIL)' "$L/inv.log"; then exit 1; fi
# AFTER the controller commits (the tracked-binary gate reads committed blobs; this leg runs go test ./...):
VERIFY_GO_LEGS=plain ./scripts/verify_go.sh 2>&1 | tee "$L/verify_go.log"
git status --short                                                         # only the §5 files
