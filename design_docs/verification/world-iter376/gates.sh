#!/usr/bin/env bash
WT=/Users/voightkampff/dev/sunholo-data/.wt-world-iter376-row201-a1b-s2
OUT=$HOME/.ailang/state/world-iter376/gates-$(git -C $WT rev-parse --short=8 HEAD)-r2.log
cd $WT || exit 2; export AILANG_BIN=$HOME/.pinned-ailang/ailang
run(){ echo "=== $1"; shift; s=$(date +%s); bash -c "$*"; rc=$?; echo "=== rc=$rc ($(( $(date +%s)-s ))s)"; }
{
echo "HEAD $(git rev-parse HEAD)"; $AILANG_BIN --version
run verify_ail ./scripts/verify_ail.sh
run verify_go_plain "VERIFY_GO_LEGS=plain ./scripts/verify_go.sh"
run vet "go build ./... && go vet ./..."
run roots "go test ./host/store/ -run '^TestProductionContextRoots\$' -count=1"
run allowlist "go test ./host/daemon/ -run '^TestDaemonDependencyAllowlist\$' -count=1"
run race "go test -race ./host/store/ ./host/approval/... ./host/authority/ -count=1 && go test -race ./cmd/ailang-worldd/ -run '^(TestApprover|TestReadTermLineErrWithin|TestConfirmApproverAct|TestReadRegularFileCapped|TestReadCappedRefusesGrowth)' -count=1"
run cmd_full "go test ./cmd/ailang-worldd/ -count=1"
echo "ALL GATES DONE for $(git rev-parse --short=8 HEAD)"
} > "$OUT" 2>&1
echo "$OUT"
