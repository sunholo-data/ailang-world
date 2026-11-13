#!/bin/bash
# A4 sprint 2 pre-commit gates (plan §6). usage: gates.sh <tree> <logdir> <base|final>
# macOS /bin/bash 3.2: no declare -A. Every leg records its own rc; every leg runs (no set -e).
T=$1; L=$2; MODE=$3
mkdir -p "$L"; cd "$T" || exit 2
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
leg verify_ail './scripts/verify_ail.sh'
leg gofmt 'out=$(gofmt -l cmd host); echo "$out"; test -z "$out"'
leg build 'go build ./...'
leg vet 'go vet ./...'
leg fence "go test -run '^\$' ./... > '$L/fence.raw' 2>&1; rc=\$?; cat '$L/fence.raw'; if grep -E '^(FAIL|--- FAIL)' '$L/fence.raw' >/dev/null; then exit 1; fi; exit \$rc"
leg pkgs 'go test -count=1 -timeout 30m ./host/approval/... ./host/authority/ ./host/store/ ./host/sshsign/ ./host/surfacepage/ ./host/daemon/ ./cmd/ailang-worldd/'
leg broker_named 'go test -count=1 ./host/broker/ -run "^(TestReservedApprovalScope|TestDecisionRefusal|TestDecideApproval|TestPublishApproval|TestApprovalContextIdentity|TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential)" -v'
leg broker_whole_reported 'go test -count=1 -timeout 30m ./host/broker/'
leg race 'go test -race -count=1 ./host/approval/ ./host/authority/ ./host/store/ ./host/sshsign/ && go test -race -count=1 ./host/daemon/ -run "^(TestDecision|TestDecide|TestApprov|TestSurfaceApprove|TestDaemonCredentialWriterSites|TestEpisodeFinish)" && go test -race -count=1 ./cmd/ailang-worldd/ -run "^(TestApproveDecide|TestApprove|TestApprover)"'
NEW_A="TestDecideBySignatureApprover TestMintSignatureIsNotADecision TestDecisionBindsTheVerb TestSignedDecisionBinds"
NEW_AU="TestJudgeDecisionSSH"
NEW_S="TestCredentialsSchema8KeyDecisions TestKeyDecisionsRefChecks TestRecordKeyDecisionRevokedSince TestRecordKeyDecisionReplaceUnlanded TestRecordKeyDecisionGrammar"
NEW_D="TestDecideRefGrammarIsAnchored TestDecisionSubjectRefusesAnUndecodablePacket TestDecidePOSTSSHSIG TestDecidePOSTSSHSIGRefusals TestDecidePOSTBusyWriter TestDecidePOSTRetryAfterAFailedWorldWrite TestDecideRefusesReservedScopesByKey TestDecideConcurrentPOSTsDecideOnce"
NEW_C="TestApproveDecideTyped TestApproveDecideNoTerminal TestApproveDecideFlagRules TestApproveDecideRefusesAMismatchedSubject TestApproveDecideSignWith TestApproveDecideSignWithSSHKeygen"
if [ "$MODE" = base ]; then NEW_A=""; NEW_AU=""; NEW_S=""; NEW_D=""; NEW_C=""; fi
alt() { local x="$*"; x="${x// /|}"; echo "${x:+|$x}"; }
INV_APPROVAL="TestApprovalImportsOnlyTheStandardLibrary TestApprovalTestIsTestOnly TestDecisionMessageLayout TestDecisionMessageRefusesBadInput"
INV_AUTH="TestJudgeSSHSignatureIsTheCLIJudge TestJudgeAssertionIsTheCLIJudge"
INV_BROKER="TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential TestQuickstartSection10ProfilesLoad TestQuickstartProtectedContract TestApprovalContextIdentity TestReservedApprovalScope TestDecideApprovalCallerCensus"
INV_DAEMON="TestDaemonDependencyAllowlist TestFrozenV1RouteTableMatchesMux TestSurfaceAppRoutes TestApprovalRoutesRegistered TestDecisionRoutesRegistered TestReadCtxCancelledAfterHandler TestMCPMountMethodSource TestDaemonCredentialWriterSites TestQuickstartDecisionPacketSection TestDaemonConstructorCallerInventory TestDecidePOSTGuards TestDecideDiscoveryFindsAHiddenSecondPacket TestDecideDiscoveryRefusesAnOverflow"
INV_BOUNDARY="TestWorldBoundaryDependencyAllowlist TestBareNetHTTPExemptionIsPerGroup TestEveryCommandDirectoryIsAProtectedGroup TestBoundaryASTWriteGuard TestWorkbenchPackageRemainsTransportFree"
INV_STORE="TestProductionContextRoots TestProductionGoSurface TestNoNewDeadlineFreeStoreReads TestReadGettersHonorContext TestCredentialsMethodsRequireADeadline TestCredentialsSchemaMatchesCanonicalDDL TestCredentialsSchema7SurfaceApprovals"
INV_VG="TestNoParallelWriteForkPackagesOutsideVerifygate TestSubprocessSinkGate TestVerifygateTestWritesAreForkLocked"
INV_CMD="TestEveryVerbHasHelp TestQuickstartAttendedStepsBoxMatchesSection9 TestQuickstartSection9FlagsMatchTheCLI TestQuickstartSection10FlagsMatchTheCLI TestQuickstartSection11FlagsMatchTheCLI TestQuickstartServiceSectionMatchesTheCLIHelp TestRealSessionEnvSignsWithSSHKeygen TestSignWithSSHKeygenEndToEnd"
INV_RUNBOOK="TestNoCIStepOrScriptReachesThePublishEntrypoint"
leg inv "set -o pipefail
go test -v -count=1 ./host/approval/... -run '^($(echo $INV_APPROVAL | tr ' ' '|')$(alt $NEW_A))\$' &&
go test -v -count=1 ./host/authority/ -run '^($(echo $INV_AUTH | tr ' ' '|')$(alt $NEW_AU))\$' &&
go test -v -count=1 ./host/broker/ -run '^($(echo $INV_BROKER | tr ' ' '|'))\$' &&
go test -v -count=1 ./host/daemon/ -run '^($(echo $INV_DAEMON | tr ' ' '|')$(alt $NEW_D))\$' &&
go test -v -count=1 ./host/boundary/ -run '^($(echo $INV_BOUNDARY | tr ' ' '|'))\$' &&
go test -v -count=1 ./host/store/ -run '^($(echo $INV_STORE | tr ' ' '|')$(alt $NEW_S))\$' &&
go test -v -count=1 ./host/verifygate/ -run '^($(echo $INV_VG | tr ' ' '|'))\$' &&
go test -v -count=1 ./host/runbook/ -run '^($(echo $INV_RUNBOOK | tr ' ' '|'))\$' &&
go test -v -count=1 ./cmd/ailang-worldd/ -run '^($(echo $INV_CMD | tr ' ' '|')$(alt $NEW_C))\$'"
# every named test must PASS (none SKIP, none FAIL)
MISSING=""
for t in $INV_APPROVAL $INV_AUTH $INV_BROKER $INV_DAEMON $INV_BOUNDARY $INV_STORE $INV_VG $INV_CMD $INV_RUNBOOK $NEW_A $NEW_AU $NEW_S $NEW_D $NEW_C; do
  grep -qE -- "^--- PASS: $t " "$L/inv.log" || MISSING="$MISSING $t"
done
if grep -E '^\s*--- (SKIP|FAIL)' "$L/inv.log"; then FAILS="$FAILS inv_skip_or_fail"; fi
[ -z "$MISSING" ] || { echo "NOT PASSED:$MISSING" | tee -a "$L/legs.txt"; FAILS="$FAILS inv_missing"; }
echo "named PASS: $(grep -cE '^--- PASS' "$L/inv.log")" | tee -a "$L/legs.txt"
echo "GATES_DONE $(( $(date +%s) - s0 ))s FAILS:${FAILS:- none}" | tee -a "$L/legs.txt"
