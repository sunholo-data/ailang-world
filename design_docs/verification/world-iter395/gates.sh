#!/bin/bash
# Row 243 V1 gates (plan §6). usage: gates.sh <tree> <logdir> <base|final> [<sha>]
# The controller: bash gates.sh <worktree> $HOME/.ailang/state/world-iter395/gates final <HEAD sha>
# macOS /bin/bash 3.2: no declare -A. Every leg records its own rc; every leg runs (no set -e).
# Never start this under nohup. Writes <logdir>/ALL_DONE.<sha> when every leg has finished.
T=$1; L=$2; MODE=$3
mkdir -p "$L"; cd "$T" || exit 2
SHA=${4:-$(git rev-parse HEAD 2>/dev/null || echo unknown)}
export AILANG_BIN=$HOME/.pinned-ailang/ailang
export WORLD_EXEC_NODE=/opt/homebrew/bin/node
unset WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER WORLD_PASSKEY_RECORD_DIR WORLD_UPDATE_VERB_GOLDEN
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
PK="./host/surfacepage/ ./host/daemon/ ./cmd/ailang-worldd/"
[ "$MODE" = base ] || PK="./host/verbs/ $PK"
leg pkgs "go test -count=1 -timeout 30m $PK"
RACEV=""; [ "$MODE" = base ] || RACEV="go test -race -count=1 ./host/verbs/ && "
leg race "$RACEV go test -race -count=1 ./host/surfacepage/ && go test -race -count=1 ./host/daemon/ -run '^(TestSurface|TestEpisodeReport|TestApprov|TestDaemonCredentialWriterSites)' && go test -race -count=1 ./cmd/ailang-worldd/ -run '^(TestEpisodeReport|TestVerb|TestSession|TestApprove|TestPreset|TestEpisode)'"
NEW_V="TestSplit TestQuoteRoundTrip TestLookup TestEpisodeReportFlags TestParseInterleaved TestApproveTarget"
NEW_D="TestSurfaceVerbRouteRegistered TestSurfaceVerbOutcomes TestSurfaceVerbBoundsAndErrors TestSurfaceVerbIsReadOnly"
NEW_SP="TestSurfaceAppVerbs"
NEW_C="TestVerbTableMatchesTheCLISwitch TestVerbRouteReportIsTheCLIs TestQuickstartSection19MatchesTheBar"
NEW_X="TestVerbBarEpisodeReportInChrome"
if [ "$MODE" = base ]; then NEW_V=""; NEW_D=""; NEW_SP=""; NEW_C=""; NEW_X=""; fi
alt() { local x="$*"; x="${x// /|}"; echo "${x:+|$x}"; }
INV_DAEMON="TestDaemonDependencyAllowlist TestFrozenV1RouteTableMatchesMux TestSurfaceAppRoutes TestApprovalRoutesRegistered TestDecisionRoutesRegistered TestReadCtxCancelledAfterHandler TestMCPMountMethodSource TestDaemonCredentialWriterSites TestQuickstartDecisionPacketSection TestDaemonConstructorCallerInventory TestApprovePageRoutes TestPasskeyPageRoutes TestApprovalRoutesRefuseForeignHost TestGETRoutesRejectOtherMethods TestDaemonCallerInventoryAliasControl TestEpisodeReportRESTRoute"
INV_BOUNDARY="TestWorldBoundaryDependencyAllowlist TestBareNetHTTPExemptionIsPerGroup TestEveryCommandDirectoryIsAProtectedGroup TestBoundaryASTWriteGuard TestWorkbenchPackageRemainsTransportFree"
INV_STORE="TestProductionContextRoots TestProductionGoSurface TestNoNewDeadlineFreeStoreReads TestReadGettersHonorContext TestCredentialsMethodsRequireADeadline"
INV_VG="TestNoParallelWriteForkPackagesOutsideVerifygate TestSubprocessSinkGate TestVerifygateTestWritesAreForkLocked"
INV_SP="TestSurfaceAppPure TestSurfaceScriptPure TestSurfaceScriptAssetSeams TestSurfaceAppMatchesWorkingCopy TestSurfaceAppShellAssembly TestNeedsYouListsDecisions"
INV_CMD="TestEveryVerbHasHelp TestEpisodeReportUsage TestEpisodeReportTablePrintsEveryClass TestEpisodeReportTablePrintsTheReadsLine TestEpisodeReportJSONAndTableLoopback TestEpisodeReportFollowReprintsOnChangeLoopback TestEpisodeReportPrintsExecCache TestEpisodeGlobalAddrRefusedForStartOnly TestGlobalBareAddrIsNormalisedForEveryClientVerb TestQuickstartAttendedStepsBoxMatchesSection9 TestQuickstartSection9FlagsMatchTheCLI TestQuickstartSection10FlagsMatchTheCLI TestQuickstartSection11FlagsMatchTheCLI TestQuickstartSection18FlagsMatchTheCLI TestQuickstartSection18CardMatchesThePage TestQuickstartServiceSectionMatchesTheCLIHelp"
leg inv "set -o pipefail
go test -v -count=1 ./host/daemon/ -run '^($(echo $INV_DAEMON | tr ' ' '|')$(alt $NEW_D))\$' &&
go test -v -count=1 ./host/boundary/ -run '^($(echo $INV_BOUNDARY | tr ' ' '|'))\$' &&
go test -v -count=1 ./host/store/ -run '^($(echo $INV_STORE | tr ' ' '|'))\$' &&
go test -v -count=1 ./host/verifygate/ -run '^($(echo $INV_VG | tr ' ' '|'))\$' &&
go test -v -count=1 ./host/surfacepage/ -run '^($(echo $INV_SP | tr ' ' '|')$(alt $NEW_SP))\$' &&
go test -v -count=1 ./cmd/ailang-worldd/ -run '^($(echo $INV_CMD | tr ' ' '|')$(alt $NEW_C))\$' &&
{ [ -z '$NEW_V' ] || go test -v -count=1 ./host/verbs/ -run '^($(echo $NEW_V | tr ' ' '|'))\$'; }"
XB="TestApprovePageInChrome TestApprovePageRefusesASwappedRequest TestApprovePageSSHSIGInChrome TestApprovePageShowsAHostileValueAsText TestPasskeyPageAbortsAPendingCeremony TestPasskeyPageInChrome TestApproveDecideWebAuthnLinkInChrome TestDecidePageApproveWithPasskey TestDecidePageFailsClosed TestDecidePageSSHSIGPaste TestNeedsYouListsPendingRequests TestNeedsYouNotifiesOnlyWhenAskedOnce"
XD="TestSurfaceLiveChrome TestSurfaceFollowChrome"
CH="'/Applications/Google Chrome.app/Contents/MacOS/Google Chrome'"
leg chrome "set -o pipefail; WORLD_CHROME=$CH WORLD_REQUIRE_PASSKEY_BROWSER=1 WORLD_REQUIRE_SSH_KEYGEN=1 go test ./cmd/ailang-worldd/ -run '^($(echo $XB | tr ' ' '|')$(alt $NEW_X))\$' -count=1 -v && WORLD_CHROME=$CH go test ./host/daemon/ -run '^($(echo $XD | tr ' ' '|'))\$' -count=1 -v"
for t in $XB $NEW_X $XD; do grep -qE -- "^--- PASS: $t " "$L/chrome.log" || { echo "CHROME NOT PASSED: $t" | tee -a "$L/legs.txt"; FAILS="$FAILS chrome_missing"; }; done
MISSING=""
for t in $INV_DAEMON $NEW_D $INV_BOUNDARY $INV_STORE $INV_VG $INV_SP $NEW_SP $INV_CMD $NEW_C $NEW_V; do
  grep -qE -- "^--- PASS: $t " "$L/inv.log" || MISSING="$MISSING $t"
done
if grep -E '^\s*--- (SKIP|FAIL)' "$L/inv.log"; then FAILS="$FAILS inv_skip_or_fail"; fi
[ -z "$MISSING" ] || { echo "NOT PASSED:$MISSING" | tee -a "$L/legs.txt"; FAILS="$FAILS inv_missing"; }
echo "named PASS: $(grep -cE '^--- PASS' "$L/inv.log")" | tee -a "$L/legs.txt"
leg verify_go 'VERIFY_GO_LEGS=plain AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_go.sh'
echo "GATES_DONE $(( $(date +%s) - s0 ))s FAILS:${FAILS:- none}" | tee -a "$L/legs.txt"
echo "GATES_DONE $(( $(date +%s) - s0 ))s FAILS:${FAILS:- none}" > "$L/ALL_DONE.$SHA"
