bash <<'EOF'
set -euo pipefail
export AILANG_BIN=$HOME/.pinned-ailang/ailang
W=/Users/voightkampff/dev/sunholo-data/.wt-world-iter380-row201-a2s2; cd "$W"
L=$HOME/.ailang/state/world-iter380/gates; mkdir -p "$L"
./scripts/verify_ail.sh 2>&1 | tee "$L/verify_ail.log"                      # no .ail change: exactly green
out=$(gofmt -l cmd host); echo "$out" | tee "$L/gofmt.log"; test -z "$out"
go build ./... 2>&1 | tee "$L/build.log"
go vet ./... 2>&1 | tee "$L/vet.log"
go test -run '^$' ./... 2>&1 | tee "$L/fence.log"                           # compile fence for every _test.go
go test -count=1 ./host/approval/... ./host/store/ ./host/authority/ ./cmd/ailang-worldd/ ./host/daemon/ 2>&1 | tee "$L/pkgs.log"
go test -race -count=1 ./host/approval/... ./host/store/ ./host/authority/ 2>&1 | tee "$L/race.log"
# the repo-wide inventories (design §7) and the reflective tests, by name, -v:
go test -v -count=1 ./host/broker/ -run '^TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential$' 2>&1 | tee "$L/inv.log"
go test -v -count=1 ./host/daemon/ -run '^(TestDaemonDependencyAllowlist|TestFrozenV1RouteTableMatchesMux|TestSurfaceAppRoutes|TestDaemonConstructorCallerInventory)$' 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./host/boundary/ -run '^(TestWorldBoundaryDependencyAllowlist|TestBareNetHTTPExemptionIsPerGroup|TestEveryCommandDirectoryIsAProtectedGroup)$' 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./host/store/ -run '^(TestProductionContextRoots|TestProductionGoSurface|TestNoNewDeadlineFreeStoreReads|TestReadGettersHonorContext|TestCredentialsMethodsRequireADeadline)$' 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./host/verifygate/ -run '^TestNoParallelWriteForkPackagesOutsideVerifygate$' 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./host/approval/... -run '^(TestApprovalImportsOnlyTheStandardLibrary|TestApprovalTestIsTestOnly)$' 2>&1 | tee -a "$L/inv.log"
go test -v -count=1 ./host/approval/ ./host/store/ ./host/authority/ -run '^(TestSignedApprovalFieldsAreAllBound|TestEnrolmentPromptShowsEveryPersistedField|TestRequestConsumeCompleteness|TestApproverRowMatchesTheTable|TestPasskeyRowMatchesTheTable|TestRequestMintFieldsAreSourced)$' 2>&1 | tee -a "$L/inv.log"
for t in TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential TestDaemonDependencyAllowlist TestFrozenV1RouteTableMatchesMux \
  TestSurfaceAppRoutes TestDaemonConstructorCallerInventory TestWorldBoundaryDependencyAllowlist TestBareNetHTTPExemptionIsPerGroup \
  TestEveryCommandDirectoryIsAProtectedGroup TestProductionContextRoots TestProductionGoSurface TestNoNewDeadlineFreeStoreReads \
  TestReadGettersHonorContext TestCredentialsMethodsRequireADeadline TestNoParallelWriteForkPackagesOutsideVerifygate \
  TestApprovalImportsOnlyTheStandardLibrary TestApprovalTestIsTestOnly TestSignedApprovalFieldsAreAllBound \
  TestEnrolmentPromptShowsEveryPersistedField TestRequestConsumeCompleteness TestApproverRowMatchesTheTable \
  TestPasskeyRowMatchesTheTable TestRequestMintFieldsAreSourced; do
  grep -qE -- "^--- PASS: $t " "$L/inv.log" || { echo "$t did not PASS"; exit 1; }; done
if grep -E '^\s*--- (SKIP|FAIL)' "$L/inv.log"; then exit 1; fi
echo PRE_COMMIT_GATES_GREEN | tee "$L/pre_commit.ok"
EOF
# AFTER the controller commits (the tracked-binary gate reads committed blobs; this leg runs go test ./...):
bash -c 'set -o pipefail; export AILANG_BIN=$HOME/.pinned-ailang/ailang; cd /Users/voightkampff/dev/sunholo-data/.wt-world-iter380-row201-a2s2 && VERIFY_GO_LEGS=plain ./scripts/verify_go.sh 2>&1 | tee $HOME/.ailang/state/world-iter380/gates/verify_go.log; echo "verify_go rc=$?" | tee -a $HOME/.ailang/state/world-iter380/gates/verify_go.log'
git -C /Users/voightkampff/dev/sunholo-data/.wt-world-iter380-row201-a2s2 status --short   # only the two plan files, until the controller commits them
