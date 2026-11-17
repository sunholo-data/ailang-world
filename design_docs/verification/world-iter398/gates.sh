#!/bin/bash
# Row 243 V2b gates (plan §6). usage: gates.sh <tree> <logdir> <base|final> [<sha>]
# The controller: bash gates.sh <worktree> $HOME/.ailang/state/world-iter398/gates final <HEAD sha>
# macOS /bin/bash 3.2: no declare -A. Every leg records its own rc; every leg runs (no set -e); every leg is tee'd
# to <logdir>/<leg>.txt. GATES_SKIP_VERIFY_GO=1 skips verify_go (the executor sets it; the controller runs it).
# Never start this under nohup. Writes <logdir>/ALL_DONE.<sha> with FAILS: when every leg has finished.
T=$1; L=$2; MODE=$3
mkdir -p "$L"; cd "$T" || exit 2
SHA=${4:-$(git rev-parse HEAD 2>/dev/null || echo unknown)}
export AILANG_BIN=$HOME/.pinned-ailang/ailang
export WORLD_EXEC_NODE=/opt/homebrew/bin/node
unset WORLD_CHROME WORLD_REQUIRE_PASSKEY_BROWSER WORLD_PASSKEY_RECORD_DIR WORLD_UPDATE_VERB_GOLDEN WORLD_UPDATE_WHY_GOLDEN WORLD_TOOL_AILANG_BIN
: > "$L/legs.txt"
FAILS=""
s0=$(date +%s)
leg() { # leg <name> <command string>: tee'd to $L/<name>.txt, rc recorded in legs.txt
  local name=$1 cmd=$2 s rc
  s=$(date +%s)
  bash -c "$cmd" 2>&1 | tee "$L/$name.txt"
  rc=${PIPESTATUS[0]}
  echo "LEG $name rc=$rc $(( $(date +%s) - s ))s" | tee -a "$L/legs.txt"
  [ "$rc" -eq 0 ] || FAILS="$FAILS $name"
}
need() { # need <leg> <test id>...: each must show --- PASS in <leg>.txt
  local lg=$1 t; shift
  for t in "$@"; do grep -qE -- "^ *--- PASS: $t " "$L/$lg.txt" || { echo "NOT PASSED ($lg): $t" | tee -a "$L/legs.txt"; FAILS="$FAILS ${lg}_missing"; }; done
}
alt() { local x="$*"; x="${x// /|}"; echo "${x:+|$x}"; }
leg status 'git status --porcelain | wc -l'
leg verify_ail './scripts/verify_ail.sh'
leg gofmt 'out=$(gofmt -l cmd host); echo "$out"; test -z "$out"'
leg build 'go build ./...'
leg vet 'go vet ./...'
leg fence "go test -run '^\$' ./... > '$L/fence.raw' 2>&1; rc=\$?; cat '$L/fence.raw'; if grep -E '^(FAIL|--- FAIL)' '$L/fence.raw' >/dev/null; then exit 1; fi; exit \$rc"
leg pkgs "go test -count=1 -timeout 30m ./host/verbs/ ./host/why/ ./host/surfacepage/ ./host/daemon/ ./cmd/ailang-worldd/"
leg race "go test -race -count=1 ./host/verbs/ ./host/why/ ./host/surfacepage/ && go test -race -count=1 ./host/daemon/ -run '^(TestSurface|TestWhy|TestVerb|TestEpisodeReport)' && go test -race -count=1 ./cmd/ailang-worldd/ -run '^(TestWhy|TestVerb|TestQuickstartSection19)'"
NEW_V="TestWhyFlags TestAdmitWhy TestNeverAndOtherReasonsEndWithTheCLILine TestSplitBoundsAreLiterals"
NEW_D="TestSurfaceVerbWhyOutcomes TestSurfaceVerbWhyDeadline TestSurfaceVerbWhyReadFailure TestWhyGetterTransportParity TestWhyGetterUnroutedAndOversize TestVerbAnswerCap TestVerbAnswerCapEndToEnd TestWhyHandlersReadNoMiddlewareState TestSurfaceVerbWhyIsReadOnly TestWhyResult"
NEW_C="TestWhyBarAdmissionIsTheCLIs TestVerbRouteWhyIsTheCLIs"
if [ "$MODE" = base ]; then NEW_V=""; NEW_D=""; NEW_C=""; fi
INV_DAEMON="TestDaemonDependencyAllowlist TestFrozenV1RouteTableMatchesMux TestSurfaceAppRoutes TestApprovalRoutesRegistered TestDecisionRoutesRegistered TestReadCtxCancelledAfterHandler TestDaemonCredentialWriterSites TestDaemonConstructorCallerInventory TestGETRoutesRejectOtherMethods TestDaemonCallerInventoryAliasControl TestEpisodeReportRESTRoute TestSurfaceVerbRouteRegistered TestSurfaceVerbOutcomes TestSurfaceVerbBoundsAndErrors TestSurfaceVerbIsReadOnly"
INV_BOUNDARY="TestWorldBoundaryDependencyAllowlist TestBareNetHTTPExemptionIsPerGroup TestEveryCommandDirectoryIsAProtectedGroup TestBoundaryASTWriteGuard TestWorkbenchPackageRemainsTransportFree TestPinfetchStaysOutOfStoreReplayAndPublish"
INV_STORE="TestProductionContextRoots TestProductionGoSurface TestNoNewDeadlineFreeStoreReads TestReadGettersHonorContext"
INV_VG="TestNoParallelWriteForkPackagesOutsideVerifygate TestSubprocessSinkGate TestVerifygateTestWritesAreForkLocked"
INV_SP="TestSurfaceAppPure TestSurfaceScriptPure TestSurfaceAppMatchesWorkingCopy TestSurfaceAppShellAssembly TestSurfaceAppVerbs"
INV_VERBS="TestSplit TestQuoteRoundTrip TestLookup TestEpisodeReportFlags TestParseInterleaved TestApproveTarget"
INV_WHY="TestWhyImportsNoEffect TestParseTargetForms TestHeadAndGetJSONClasses TestElideBounds"
INV_CMD="TestEveryVerbHasHelp TestVerbTableMatchesTheCLISwitch TestVerbRouteReportIsTheCLIs TestQuickstartSection19MatchesTheBar TestWhyCharacterisedPaths TestWhyRunReplaysTheCLI TestWhyRunGetterClasses TestWhyConstantsMatchTheCLI TestWhyHeadIsOneReader TestWhyTargetFormsResolveToOneEntry TestWhyScanBound TestWhyBrokenLinkIsIntegrityRefusal"
leg inv "set -o pipefail
go test -v -count=1 ./host/daemon/ -run '^($(echo $INV_DAEMON | tr ' ' '|')$(alt $NEW_D))\$' &&
go test -v -count=1 ./host/boundary/ -run '^($(echo $INV_BOUNDARY | tr ' ' '|'))\$' &&
go test -v -count=1 ./host/store/ -run '^($(echo $INV_STORE | tr ' ' '|'))\$' &&
go test -v -count=1 ./host/verifygate/ -run '^($(echo $INV_VG | tr ' ' '|'))\$' &&
go test -v -count=1 ./host/surfacepage/ -run '^($(echo $INV_SP | tr ' ' '|'))\$' &&
go test -v -count=1 ./host/verbs/ -run '^($(echo $INV_VERBS | tr ' ' '|')$(alt $NEW_V))\$' &&
go test -v -count=1 ./host/why/ -run '^($(echo $INV_WHY | tr ' ' '|'))\$' &&
go test -v -count=1 ./cmd/ailang-worldd/ -run '^($(echo $INV_CMD | tr ' ' '|')$(alt $NEW_C))\$'"
need inv $INV_DAEMON $NEW_D $INV_BOUNDARY $INV_STORE $INV_VG $INV_SP $INV_VERBS $NEW_V $INV_WHY $INV_CMD $NEW_C
if grep -E '^\s*--- (SKIP|FAIL)' "$L/inv.txt"; then FAILS="$FAILS inv_skip_or_fail"; fi
echo "named PASS: $(grep -cE '^--- PASS' "$L/inv.txt")" | tee -a "$L/legs.txt"
# the real-stack why leg: CI skips these (WORLD_TOOL_AILANG_BIN unset); here they must PASS, not SKIP
leg realstack "set -o pipefail; WORLD_TOOL_AILANG_BIN=\$HOME/.pinned-ailang-tools/v0.54.15/ailang go test -count=1 -v -timeout 20m ./cmd/ailang-worldd/ -run '^(TestDemoSeedWhyNamedLinks|TestCLIAgainstSeToolsDaemon)\$'"
need realstack 'TestDemoSeedWhyNamedLinks/real-stack' TestDemoSeedWhyNamedLinks TestCLIAgainstSeToolsDaemon
if grep -E -- '--- (SKIP|FAIL)' "$L/realstack.txt"; then FAILS="$FAILS realstack_skip_or_fail"; fi
XB="TestApprovePageInChrome TestApprovePageRefusesASwappedRequest TestApprovePageSSHSIGInChrome TestApprovePageShowsAHostileValueAsText TestPasskeyPageAbortsAPendingCeremony TestPasskeyPageInChrome TestApproveDecideWebAuthnLinkInChrome TestDecidePageApproveWithPasskey TestDecidePageFailsClosed TestDecidePageSSHSIGPaste TestNeedsYouListsPendingRequests TestNeedsYouNotifiesOnlyWhenAskedOnce TestVerbBarEpisodeReportInChrome"
NEW_X="TestVerbBarWhyInChrome"; [ "$MODE" = base ] && NEW_X=""
XD="TestSurfaceLiveChrome TestSurfaceFollowChrome"
CH="'/Applications/Google Chrome.app/Contents/MacOS/Google Chrome'"
leg chrome "set -o pipefail; WORLD_CHROME=$CH WORLD_REQUIRE_PASSKEY_BROWSER=1 WORLD_REQUIRE_SSH_KEYGEN=1 go test ./cmd/ailang-worldd/ -run '^($(echo $XB | tr ' ' '|')$(alt $NEW_X))\$' -count=1 -v && WORLD_CHROME=$CH go test ./host/daemon/ -run '^($(echo $XD | tr ' ' '|'))\$' -count=1 -v"
need chrome $XB $NEW_X $XD
if [ "${GATES_SKIP_VERIFY_GO:-0}" = 1 ]; then echo "LEG verify_go SKIPPED (GATES_SKIP_VERIFY_GO=1)" | tee -a "$L/legs.txt"
else leg verify_go 'VERIFY_GO_LEGS=plain AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_go.sh'; fi
echo "GATES_DONE $(( $(date +%s) - s0 ))s FAILS:${FAILS:- none}" | tee -a "$L/legs.txt"
echo "GATES_DONE $(( $(date +%s) - s0 ))s FAILS:${FAILS:- none}" > "$L/ALL_DONE.$SHA"
