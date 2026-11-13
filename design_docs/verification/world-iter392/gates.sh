#!/bin/bash
# A4 sprint 4 gates (plan §6). usage: gates.sh <tree> <logdir> <base|final> [<sha>]
# The controller: bash gates.sh <worktree> $HOME/.ailang/state/world-iter392/gates final <HEAD sha>
# macOS /bin/bash 3.2: no declare -A. Every leg records its own rc; every leg runs (no set -e).
# Never start this under nohup. Writes <logdir>/ALL_DONE.<sha> when every leg has finished.
T=$1; L=$2; MODE=$3
mkdir -p "$L"; cd "$T" || exit 2
SHA=${4:-$(git rev-parse HEAD 2>/dev/null || echo unknown)}
export AILANG_BIN=$HOME/.pinned-ailang/ailang
unset WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER WORLD_PASSKEY_RECORD_DIR
: > "$L/legs.txt"
FAILS=""
s0=$(date +%s)
leg() { # leg <name> <command string>: tee'd to $L/<name>.log, rc recorded in legs.txt
  local name=$1 cmd=$2 s rc
  s=$(date +%s)
  bash -c "$cmd" 2>&1 | tee "$L/$name.log"
  rc=${PIPESTATUS[0]}
  echo "LEG $name rc=$rc $(( $(date +%s) - s ))s" | tee -a "$L/legs.txt"
  [ "$rc" -eq 0 ] || FAILS="$FAILS $name"
}
leg status 'git status --porcelain | wc -l'
leg verify_ail './scripts/verify_ail.sh'
leg gofmt 'out=$(gofmt -l cmd host); echo "$out"; test -z "$out"'
leg build 'go build ./...'
leg vet 'go vet ./...'
leg fence "go test -run '^\$' ./... > '$L/fence.raw' 2>&1; rc=\$?; cat '$L/fence.raw'; if grep -E '^(FAIL|--- FAIL)' '$L/fence.raw' >/dev/null; then exit 1; fi; exit \$rc"
leg pkgs 'go test -count=1 -timeout 30m ./host/approval/... ./host/authority/ ./host/store/ ./host/surfacepage/ ./host/coordinator/ ./host/daemon/ ./cmd/ailang-worldd/'
leg broker_named 'go test -count=1 ./host/broker/ -run "^(TestReservedApprovalScope|TestDecisionRefusal|TestDecideApproval|TestPublishApproval|TestApprovalContextIdentity|TestRecentApprovals|TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential)" -v'
leg race 'go test -race -count=1 ./host/surfacepage/ && go test -race -count=1 ./host/daemon/ -run "^(TestDecision|TestDecide|TestApprov|TestSurfaceApprove|TestSurfaceRequests|TestSurfaceApp|TestDaemonCredentialWriterSites)" && go test -race -count=1 ./cmd/ailang-worldd/ -run "^(TestApproveDecide|TestApproveList|TestApprove)"'
NEW_A="TestDecideByAssertionRefusesAMissingMessageAtW6 TestSignCountLaw"
NEW_SP="TestDecideModelListsPendingDecisionsOnly TestDecideModelListsATwoPacketRequestRefused TestDecideModelTruncation TestNeedsYouListsDecisions"
NEW_D="TestSurfaceRequestsListsDecisions TestSurfaceRequestsDecisionsTruncated TestSurfaceRequestsDecisionsUnreadable TestDecisionSubjectPasskeys TestDecisionSubjectJSON TestDecidePageFieldTable TestDecidePageAssets"
NEW_C="TestApproveDecideArtifactPostsTheSubjectsDigests TestApproveDecideArtifactReadBounds TestApproveListShowsDecisions TestApproveListDecisionLines TestQuickstartSection18CardMatchesThePage"
NEW_X="TestDecidePageApproveWithPasskey TestDecidePageFailsClosed TestDecidePageSSHSIGPaste"
if [ "$MODE" = base ]; then NEW_A=""; NEW_SP=""; NEW_D=""; NEW_C=""; NEW_X=""; fi
alt() { local x="$*"; x="${x// /|}"; echo "${x:+|$x}"; }
INV_APPROVAL="TestApprovalImportsOnlyTheStandardLibrary TestApprovalTestIsTestOnly TestDecisionMessageLayout TestDecisionMessageRefusesBadInput"
INV_BROKER="TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential TestQuickstartSection10ProfilesLoad TestQuickstartProtectedContract TestApprovalContextIdentity TestReservedApprovalScope TestDecideApprovalCallerCensus"
INV_DAEMON="TestDaemonDependencyAllowlist TestFrozenV1RouteTableMatchesMux TestSurfaceAppRoutes TestApprovalRoutesRegistered TestDecisionRoutesRegistered TestReadCtxCancelledAfterHandler TestMCPMountMethodSource TestDaemonCredentialWriterSites TestQuickstartDecisionPacketSection TestDaemonConstructorCallerInventory TestApprovePageRoutes TestPasskeyPageRoutes TestApprovalRoutesRefuseForeignHost TestGETRoutesRejectOtherMethods TestDaemonCallerInventoryAliasControl"
INV_BOUNDARY="TestWorldBoundaryDependencyAllowlist TestBareNetHTTPExemptionIsPerGroup TestEveryCommandDirectoryIsAProtectedGroup TestBoundaryASTWriteGuard TestWorkbenchPackageRemainsTransportFree"
INV_STORE="TestProductionContextRoots TestProductionGoSurface TestNoNewDeadlineFreeStoreReads TestReadGettersHonorContext TestCredentialsMethodsRequireADeadline TestCredentialsSchemaMatchesCanonicalDDL"
INV_VG="TestNoParallelWriteForkPackagesOutsideVerifygate TestSubprocessSinkGate TestVerifygateTestWritesAreForkLocked"
INV_SP="TestSurfaceAppPure TestSurfaceScriptPure TestSurfaceFontEmbed TestSurfaceScriptAssetSeams"
INV_CMD="TestEveryVerbHasHelp TestQuickstartAttendedStepsBoxMatchesSection9 TestQuickstartSection9FlagsMatchTheCLI TestQuickstartSection10FlagsMatchTheCLI TestQuickstartSection11FlagsMatchTheCLI TestQuickstartSection18FlagsMatchTheCLI TestQuickstartServiceSectionMatchesTheCLIHelp"
INV_RUNBOOK="TestNoCIStepOrScriptReachesThePublishEntrypoint"
leg inv "set -o pipefail
go test -v -count=1 ./host/approval/... -run '^($(echo $INV_APPROVAL | tr ' ' '|')$(alt $NEW_A))\$' &&
go test -v -count=1 ./host/broker/ -run '^($(echo $INV_BROKER | tr ' ' '|'))\$' &&
go test -v -count=1 ./host/daemon/ -run '^($(echo $INV_DAEMON | tr ' ' '|')$(alt $NEW_D))\$' &&
go test -v -count=1 ./host/boundary/ -run '^($(echo $INV_BOUNDARY | tr ' ' '|'))\$' &&
go test -v -count=1 ./host/store/ -run '^($(echo $INV_STORE | tr ' ' '|'))\$' &&
go test -v -count=1 ./host/verifygate/ -run '^($(echo $INV_VG | tr ' ' '|'))\$' &&
go test -v -count=1 ./host/runbook/ -run '^($(echo $INV_RUNBOOK | tr ' ' '|'))\$' &&
WORLD_EXEC_NODE=/opt/homebrew/bin/node go test -v -count=1 ./host/surfacepage/ -run '^($(echo $INV_SP | tr ' ' '|')$(alt $NEW_SP))\$' &&
go test -v -count=1 ./cmd/ailang-worldd/ -run '^($(echo $INV_CMD | tr ' ' '|')$(alt $NEW_C))\$'"
XB="TestApprovePageInChrome TestApprovePageRefusesASwappedRequest TestApprovePageSSHSIGInChrome TestApprovePageShowsAHostileValueAsText TestPasskeyPageAbortsAPendingCeremony TestPasskeyPageInChrome TestApproveDecideWebAuthnLinkInChrome"
leg chrome "WORLD_CHROME='/Applications/Google Chrome.app/Contents/MacOS/Google Chrome' WORLD_REQUIRE_PASSKEY_BROWSER=1 WORLD_REQUIRE_SSH_KEYGEN=1 go test ./cmd/ailang-worldd/ -run '^($(echo $XB | tr ' ' '|')$(alt $NEW_X))\$' -count=1 -v"
for t in $XB $NEW_X; do grep -qE -- "^--- PASS: $t " "$L/chrome.log" || { echo "CHROME NOT PASSED: $t" | tee -a "$L/legs.txt"; FAILS="$FAILS chrome_missing"; }; done
# every named test must PASS (none SKIP, none FAIL)
MISSING=""
for t in $INV_APPROVAL $NEW_A $INV_BROKER $INV_DAEMON $INV_BOUNDARY $INV_STORE $INV_VG $INV_SP $INV_CMD $INV_RUNBOOK $NEW_SP $NEW_D $NEW_C; do
  grep -qE -- "^--- PASS: $t " "$L/inv.log" || MISSING="$MISSING $t"
done
if grep -E '^\s*--- (SKIP|FAIL)' "$L/inv.log"; then FAILS="$FAILS inv_skip_or_fail"; fi
[ -z "$MISSING" ] || { echo "NOT PASSED:$MISSING" | tee -a "$L/legs.txt"; FAILS="$FAILS inv_missing"; }
echo "named PASS: $(grep -cE '^--- PASS' "$L/inv.log")" | tee -a "$L/legs.txt"
leg verify_go 'VERIFY_GO_LEGS=plain AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_go.sh'
echo "GATES_DONE $(( $(date +%s) - s0 ))s FAILS:${FAILS:- none}" | tee -a "$L/legs.txt"
echo "GATES_DONE $(( $(date +%s) - s0 ))s FAILS:${FAILS:- none}" > "$L/ALL_DONE.$SHA"
