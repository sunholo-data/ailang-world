#!/bin/bash
cd /Users/voightkampff/dev/sunholo-data/.wt-world-iter314-arm6 || { echo "rc=99 cd" > /tmp/w314/ctl_gates.done; exit 1; }
L=/tmp/w314/ctl_gates_d7bba842.log; : > $L
run(){ echo "=== $1" >> $L; shift; ( "$@" ) >> $L 2>&1; echo "rc=$?" >> $L; }
run G1-gofmt sh -c 'test -z "$(gofmt -l host/broker/)"'
run G2-vet go vet ./host/broker/
run G3-fence go test -count=1 -run '^$' ./...
run G4-broker go test ./host/broker/ -count=1
run G5-daemon go test ./host/daemon/ -count=1 -run 'TestExecStartup|TestExecSrt' -v
run G7-race go test -race ./host/broker/ -count=1 -run 'TestExecProbe'
run G8-nodeless env -i HOME=$HOME PATH=/usr/bin:/bin:/opt/homebrew/bin/go-dummy WORLD_EXEC_SRT_NODE_MODULES=/nonexistent GOCACHE=$(go env GOCACHE) GOMODCACHE=$(go env GOMODCACHE) GOPATH=$(go env GOPATH) /opt/homebrew/bin/go test ./host/broker -count=1 -v -run 'TestExecProbeReverseLeg'
run G9-stab go test ./host/broker/ -count=10 -run 'TestExecProbeReverseLeg|TestExecProbeForeignListener'
run G11-actionlint sh -c 'command -v actionlint >/dev/null && actionlint .github/workflows/ci.yml || echo "actionlint absent"'
run G10-verify-go env AILANG_BIN=$HOME/.pinned-ailang/ailang VERIFY_GO_LEGS=plain ./scripts/verify_go.sh
echo "ALL DONE for d7bba842" >> $L; echo done > /tmp/w314/ctl_gates.done
