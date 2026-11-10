set -euo pipefail
export AILANG_BIN=$HOME/.pinned-ailang/ailang
./scripts/verify_ail.sh                                                    # no .ail change: exactly green
test -z "$(gofmt -l cmd/ailang-worldd host/approval host/store host/authority)"
go build ./... && go vet ./...
go test -run '^$' ./...                                                    # compile fence for every _test.go
go test ./cmd/ailang-worldd/ ./host/approval/... ./host/store/ ./host/authority/ -count=1
go test ./host/store/ -run '^TestProductionContextRoots$' -count=1         # the two edited pins
go test ./host/daemon/ -run '^TestDaemonDependencyAllowlist$' -count=1
go test -race ./host/approval/... ./host/store/ ./host/authority/ -count=1
go test -race ./cmd/ailang-worldd/ -run '^(TestApproveDecidedList|TestApproverRevokeListsWhatItApproved|TestApproveSignWith|TestSignWithSSHKeygen|TestCappedBuffer|TestApproveBySignatureJudgesTheApproverAtNow|TestApproveFlagRules|TestApproverRevoke)' -count=3
WORLD_REQUIRE_SSH_KEYGEN=1 go test ./host/approval/... ./cmd/ailang-worldd/ -count=1 -v \
  -run '^(TestSSHSIGCrossCheckWithSSHKeygen|TestTestSignerAgreesWithSSHKeygen|TestSignWithSSHKeygenEndToEnd)$' 2>&1 | tee /tmp/w378-xs.log | tail -3
for t in TestSSHSIGCrossCheckWithSSHKeygen TestTestSignerAgreesWithSSHKeygen TestSignWithSSHKeygenEndToEnd/ed25519 TestSignWithSSHKeygenEndToEnd/ecdsa-p256; do
  grep -qE -- "^ *--- PASS: $t " /tmp/w378-xs.log || { echo "$t did not PASS"; exit 1; }; done
go test ./cmd/ailang-worldd/ -count=1 -v -run '^(TestApprove|TestApprover|TestEpisodeStartRequest|TestSessionMint_|TestSessionNew|TestEveryVerbHasHelp|TestHelpIsNeverBlockedByTheRegistryKeyGuard|TestSignWith|TestCappedBuffer|TestRealSessionEnv)' 2>&1 | tee /tmp/w378-cl.log | tail -3
if grep -E '^\s*--- (SKIP|FAIL)' /tmp/w378-cl.log; then exit 1; fi
# AFTER the controller commits (the tracked-binary gate reads committed blobs):
VERIFY_GO_LEGS=plain ./scripts/verify_go.sh
git status --short                                                         # only the §5 files
