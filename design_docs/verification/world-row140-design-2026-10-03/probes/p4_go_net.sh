#!/bin/bash
# P4-Go-net: a cold module cache — what network does `go test` need, and does the allowlist hold?
. "$(dirname "$0")/lib.sh"
M=$S/go-tiny; EP=$S/go-tiny-ep
rm -rf $M $EP; mkdir -p $M $EP/gocache $EP/gomod $EP/home $EP/tmp
cat > $M/go.mod <<'EOF'
module example.com/tiny

go 1.26

require github.com/google/go-cmp v0.7.0
EOF
cat > $M/tiny_test.go <<'EOF'
package tiny

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestDiff(t *testing.T) {
	if d := cmp.Diff([]int{1, 2}, []int{1, 2}); d != "" {
		t.Fatal(d)
	}
}
EOF
# go.sum lines for go-cmp v0.7.0 taken from the shared module cache's verified download
GOFLAGS= GOMODCACHE=$(go env GOMODCACHE) GOPROXY=off go -C $M mod download -json github.com/google/go-cmp >/dev/null 2>&1
cat > $S/s-gonet-none.json <<EOF
{ "filesystem": { "denyRead": [], "allowRead": [], "allowWrite": ["$M", "$EP"], "denyWrite": [] },
  "network": { "allowedDomains": [], "deniedDomains": [] } }
EOF
cat > $S/s-gonet-proxy.json <<EOF
{ "filesystem": { "denyRead": [], "allowRead": [], "allowWrite": ["$M", "$EP"], "denyWrite": [] },
  "network": { "allowedDomains": ["proxy.golang.org", "sum.golang.org"], "deniedDomains": [] } }
EOF
ENVV="env -i PATH=/usr/bin:/bin:/opt/homebrew/bin HOME=$EP/home GOCACHE=$EP/gocache GOMODCACHE=$EP/gomod GOTOOLCHAIN=local GOTELEMETRY=off CLAUDE_CODE_TMPDIR=$EP/tmp"
cd $M
echo "== N1: no network allowed, cold per-episode module cache, -mod=mod"
run $ENVV GOFLAGS=-mod=mod $SRT --settings $S/s-gonet-none.json -- go test -count=1 ./...
echo "== N2: proxy.golang.org + sum.golang.org allowed"
run $ENVV GOFLAGS=-mod=mod $SRT --settings $S/s-gonet-proxy.json -- go test -count=1 ./...
run cat go.sum
echo "== N3: GOPROXY=direct (VCS fetch from github.com) is refused by the same allowlist"
rm -rf $EP/gomod; mkdir -p $EP/gomod
run $ENVV GOFLAGS=-mod=mod GOPROXY=direct GONOSUMDB=* $SRT --settings $S/s-gonet-proxy.json -- go test -count=1 ./...
echo "== N4: the agent edits go.mod to add a dependency; -mod=readonly refuses to change go.mod/go.sum"
rm -rf $EP/gomod; mkdir -p $EP/gomod
cp go.sum go.sum.bak; : > go.sum
run $ENVV GOFLAGS=-mod=readonly $SRT --settings $S/s-gonet-proxy.json -- go test -count=1 ./...
cp go.sum.bak go.sum
echo "== N5: go-ep cache sizes"
du -sh $EP/gomod $EP/gocache
