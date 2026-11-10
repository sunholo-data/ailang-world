#!/usr/bin/env bash
WT=/Users/voightkampff/dev/sunholo-data/.wt-world-iter377-row201-a1b-s3
H=$(git -C $WT rev-parse --short=8 HEAD); OUT=$HOME/.ailang/state/world-iter377/gates-$H.log
cd $WT || exit 2; export AILANG_BIN=$HOME/.pinned-ailang/ailang
run(){ echo "=== $1"; shift; s=$(date +%s); bash -c "$*"; rc=$?; echo "=== rc=$rc ($(( $(date +%s)-s ))s)"; }
{
echo "HEAD $(git rev-parse HEAD)"; $AILANG_BIN --version | head -1
run verify_ail ./scripts/verify_ail.sh
run gofmt 'test -z "$(gofmt -l cmd/ailang-worldd host/approval host/store host/authority)"'
run build_vet "go build ./... && go vet ./..."
run compile_fence "go test -run '^\$' ./... >/dev/null"
run touched "go test ./cmd/ailang-worldd/ ./host/approval/... ./host/store/ ./host/authority/ -count=1"
run roots "go test ./host/store/ -run '^TestProductionContextRoots\$' -count=1"
run allowlist "go test ./host/daemon/ -run '^TestDaemonDependencyAllowlist\$' -count=1"
run race_pkgs "go test -race ./host/approval/... ./host/store/ ./host/authority/ -count=1"
run race_cmd "go test -race ./cmd/ailang-worldd/ -run '^(TestApproveBySignature|TestApproveRevokeRace|TestApproveWriteMessage|TestApproveFlagRules|TestApproveTypedPathUnchanged|TestApproveRefusesANonHexChallenge|TestApproverVerbsUseTheNamedStores|TestApproveListUsesTheNamedStores)' -count=3"
run noskip "go test ./cmd/ailang-worldd/ -count=1 -v -run '^(TestApprove|TestApprover|TestEpisodeStartRequest|TestSessionMint_|TestSessionNew|TestEveryVerbHasHelp|TestHelpIsNeverBlockedByTheRegistryKeyGuard)' > /tmp/w377-cl.log 2>&1; rc=\$?; tail -2 /tmp/w377-cl.log; if grep -E '^\s*--- (SKIP|FAIL)' /tmp/w377-cl.log; then exit 1; fi; exit \$rc"
run verify_go_plain "VERIFY_GO_LEGS=plain ./scripts/verify_go.sh"
run status "test -z \"\$(git status --porcelain)\""
echo "ALL GATES DONE for $H"
} > "$OUT" 2>&1
echo "$OUT"
