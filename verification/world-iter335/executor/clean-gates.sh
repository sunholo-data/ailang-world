#!/bin/bash
# Same shell/environment as above, AFTER all mutations have been restored.
go vet ./host/daemon/
gofmt -l host/daemon/ > "$E/gofmt.txt"
cat "$E/gofmt.txt"
test ! -s "$E/gofmt.txt"
go test -count=1 -run '^TestEpisodeFinishApproveRefusesRecord$' ./host/daemon/
go test -count=20 -run '^TestEpisodeFinishApproveRefusesRecord$' ./host/daemon/
go test -race -count=5 -run '^TestEpisodeFinishApproveRefusesRecord$' ./host/daemon/
go test -count=1 -run Finish ./host/daemon/
go test -count=1 -v -run Finish ./host/daemon/ > "$E/finish-verbose.txt" 2>&1
cat "$E/finish-verbose.txt"
test "$(grep -c '^--- PASS:' "$E/finish-verbose.txt")" -eq 7
shasum -a 256 host/daemon/episode_finish.go host/daemon/episode_finish_test.go host/daemon/episode_finish_approve_guard_test.go
wc -l host/daemon/episode_finish.go host/daemon/episode_finish_test.go host/daemon/episode_finish_approve_guard_test.go
git diff --exit-code -- host/daemon/episode_finish.go host/daemon/episode_finish_test.go
git status --short
