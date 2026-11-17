#!/bin/bash
# Row 243 V2a gates (plan §6). usage: gates.sh <tree> <logdir> <base|final> [<sha>]
# The controller: bash gates.sh <worktree> $HOME/.ailang/state/world-iter397/gates final <HEAD sha>
# macOS /bin/bash 3.2: no declare -A. Every leg records its own rc; every leg runs (no set -e).
# Never start this under nohup. Writes <logdir>/ALL_DONE.<sha> when every leg has finished.
T=$1; L=$2; MODE=$3
mkdir -p "$L"; cd "$T" || exit 2
SHA=${4:-$(git rev-parse HEAD 2>/dev/null || echo unknown)}
export AILANG_BIN=$HOME/.pinned-ailang/ailang
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
PK="./host/verbs/ ./cmd/ailang-worldd/"
[ "$MODE" = base ] || PK="./host/why/ $PK"
leg pkgs "go test -count=1 -timeout 30m $PK"
RACEW=""; [ "$MODE" = base ] || RACEW="go test -race -count=1 ./host/why/ && "
leg race "$RACEW go test -race -count=1 ./host/verbs/ && go test -race -count=1 ./cmd/ailang-worldd/ -run '^(TestWhy|TestDemo|TestLogTail|TestProvenance|TestOutcome|TestEpisodeReport|TestMCP|TestCall|TestVerb)'"
NEW_W="TestWhyInputLineShowsThePathOfALongEdit TestWhyInputLineLeavesAShortInputUnchanged TestWhyInputLineKeyFieldBounds TestWhyInputLineFallsBackToElide TestWhyProtectedLine TestWhyProtectedJSON TestWhyImportsNoEffect TestParseTargetForms TestHeadAndGetJSONClasses TestElideBounds"
NEW_V=""   # host/verbs is untouched in V2a (WhyFlags moved to V2b, plan V2a-D2)
NEW_C="TestWhyCharacterisedPaths TestWhyRunReplaysTheCLI TestWhyRunGetterClasses TestWhyHeadIsOneReader TestWhyConstantsMatchTheCLI"
if [ "$MODE" = base ]; then NEW_W=""; NEW_V=""; NEW_C=""; fi
alt() { local x="$*"; x="${x// /|}"; echo "${x:+|$x}"; }
INV_DAEMON="TestDaemonDependencyAllowlist TestFrozenV1RouteTableMatchesMux TestSurfaceVerbIsReadOnly TestSurfaceVerbOutcomes TestSurfaceAppRoutes TestDaemonCredentialWriterSites TestDaemonConstructorCallerInventory TestDaemonCallerInventoryAliasControl"
INV_BOUNDARY="TestWorldBoundaryDependencyAllowlist TestBareNetHTTPExemptionIsPerGroup TestEveryCommandDirectoryIsAProtectedGroup TestBoundaryASTWriteGuard TestWorkbenchPackageRemainsTransportFree TestPinfetchStaysOutOfStoreReplayAndPublish"
INV_STORE="TestProductionContextRoots TestProductionGoSurface TestStoreSurfaceTakesContext"
INV_VG="TestNoParallelWriteForkPackagesOutsideVerifygate TestSubprocessSinkGate TestVerifygateTestWritesAreForkLocked"
INV_VERBS="TestSplit TestQuoteRoundTrip TestLookup TestEpisodeReportFlags TestParseInterleaved TestApproveTarget"
INV_CMD="TestEveryVerbHasHelp TestVerbTableMatchesTheCLISwitch TestVerbRouteReportIsTheCLIs TestWhyTargetFormsResolveToOneEntry TestWhyResolvesMCPInvocationID TestWhyBrokenLinkIsIntegrityRefusal TestWhyScanBound TestWhyShowsTheWorktreeDelta TestWhyShowsTheFailureLine TestLogTailFollowPrintsEachCommitOnce TestProvenanceTrailer"
[ "$MODE" = base ] && INV_CMD="$INV_CMD TestWhyInputLineShowsThePathOfALongEdit TestWhyInputLineLeavesAShortInputUnchanged TestWhyInputLineKeyFieldBounds TestWhyInputLineFallsBackToElide TestWhyProtectedLine TestWhyProtectedJSON"
[ "$MODE" = base ] || INV_CMD="$INV_CMD TestWhyProtectedLine"
leg inv "set -o pipefail
go test -v -count=1 ./host/daemon/ -run '^($(echo $INV_DAEMON | tr ' ' '|'))\$' &&
go test -v -count=1 ./host/boundary/ -run '^($(echo $INV_BOUNDARY | tr ' ' '|'))\$' &&
go test -v -count=1 ./host/store/ -run '^($(echo $INV_STORE | tr ' ' '|'))\$' &&
go test -v -count=1 ./host/verifygate/ -run '^($(echo $INV_VG | tr ' ' '|'))\$' &&
go test -v -count=1 ./host/verbs/ -run '^($(echo $INV_VERBS | tr ' ' '|')$(alt $NEW_V))\$' &&
go test -v -count=1 ./cmd/ailang-worldd/ -run '^($(echo $INV_CMD | tr ' ' '|')$(alt $NEW_C))\$' &&
{ [ -z '$NEW_W' ] || go test -v -count=1 ./host/why/ -run '^($(echo $NEW_W | tr ' ' '|'))\$'; }"
MISSING=""
for t in $INV_DAEMON $INV_BOUNDARY $INV_STORE $INV_VG $INV_VERBS $NEW_V $INV_CMD $NEW_C $NEW_W; do
  grep -qE -- "^--- PASS: $t " "$L/inv.log" || MISSING="$MISSING $t"
done
if grep -E '^\s*--- (SKIP|FAIL)' "$L/inv.log"; then FAILS="$FAILS inv_skip_or_fail"; fi
[ -z "$MISSING" ] || { echo "NOT PASSED:$MISSING" | tee -a "$L/legs.txt"; FAILS="$FAILS inv_missing"; }
echo "named PASS: $(grep -cE '^--- PASS' "$L/inv.log")" | tee -a "$L/legs.txt"
leg verify_go 'VERIFY_GO_LEGS=plain AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_go.sh'
echo "GATES_DONE $(( $(date +%s) - s0 ))s FAILS:${FAILS:- none}" | tee -a "$L/legs.txt"
echo "GATES_DONE $(( $(date +%s) - s0 ))s FAILS:${FAILS:- none}" > "$L/ALL_DONE.$SHA"
