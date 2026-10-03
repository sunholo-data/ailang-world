#!/bin/bash
# P4-Go: `go test` on a scratch copy of the ailang compiler (go.mod, go.sum, internal/, cmd/) under srt.
. "$(dirname "$0")/lib.sh"
A=$HOME/dev/sunholo-data/ailang
G=$S/go-ailang
if [ ! -d $G ]; then
  mkdir -p $G && cp -c $A/go.mod $A/go.sum $G/ && cp -cR $A/internal $A/cmd $G/
fi
EP=$S/go-ep            # per-episode cache root (outside the worktree, World-owned)
rm -rf $EP; mkdir -p $EP/gocache $EP/tmp $EP/home
GOROOT_REAL=$(go env GOROOT)
MODCACHE=$(go env GOMODCACHE)
PKGS="./internal/lexer/ ./internal/parser/"
cd $G
cat > $S/s-go-wsonly.json <<EOF
{ "filesystem": { "denyRead": [], "allowRead": [], "allowWrite": ["$G"], "denyWrite": [] },
  "network": { "allowedDomains": [], "deniedDomains": [] } }
EOF
cat > $S/s-go-ep.json <<EOF
{ "filesystem": { "denyRead": [], "allowRead": [], "allowWrite": ["$G", "$EP"], "denyWrite": [] },
  "network": { "allowedDomains": [], "deniedDomains": [] } }
EOF
cat > $S/s-go-ep-readfence.json <<EOF
{ "filesystem": { "denyRead": ["$HOME", "/private/tmp", "/tmp"], "allowRead": ["$G", "$EP", "$GOROOT_REAL", "$MODCACHE"], "allowWrite": ["$G", "$EP"], "denyWrite": [] },
  "network": { "allowedDomains": [], "deniedDomains": [] } }
EOF
echo "GOROOT=$GOROOT_REAL"; echo "GOMODCACHE=$MODCACHE"; echo "go on PATH: $(command -v go)"
echo "== G1: writes = worktree only, the operator's default env (GOCACHE ~/Library/Caches/go-build)"
run $SRT --settings $S/s-go-wsonly.json -- go test -count=1 $PKGS
echo "== G2: per-episode GOCACHE/TMPDIR/HOME outside the worktree, shared module cache READ-ONLY, offline"
ENVV="env -i PATH=/usr/bin:/bin:/opt/homebrew/bin HOME=$EP/home GOCACHE=$EP/gocache GOMODCACHE=$MODCACHE GOTOOLCHAIN=go1.26.6 GOFLAGS=-mod=readonly GOPROXY=off GOTELEMETRY=off CLAUDE_CODE_TMPDIR=$EP/tmp"
t0=$(python3 -c 'import time;print(time.time())')
run $ENVV $SRT --settings $S/s-go-ep.json -- go test -count=1 $PKGS
t1=$(python3 -c 'import time;print(time.time())'); python3 -c "print('G2 cold wall %.1fs' % ($t1-$t0))"
t0=$(python3 -c 'import time;print(time.time())')
run $ENVV $SRT --settings $S/s-go-ep.json -- go test -count=1 $PKGS
t1=$(python3 -c 'import time;print(time.time())'); python3 -c "print('G2 warm wall %.1fs' % ($t1-$t0))"
du -sh $EP/gocache
echo "== G2b: same, unsandboxed, for the overhead baseline (warm cache)"
t0=$(python3 -c 'import time;print(time.time())')
run env -i PATH=/usr/bin:/bin:/opt/homebrew/bin HOME=$EP/home GOCACHE=$EP/gocache GOMODCACHE=$MODCACHE GOTOOLCHAIN=go1.26.6 GOFLAGS=-mod=readonly GOPROXY=off GOTELEMETRY=off go test -count=1 $PKGS
t1=$(python3 -c 'import time;print(time.time())'); python3 -c "print('G2b unsandboxed warm wall %.1fs' % ($t1-$t0))"
echo "== G3: G2 plus the read fence (deny HOME and /tmp; re-allow worktree, episode cache, GOROOT, module cache)"
run $ENVV $SRT --settings $S/s-go-ep-readfence.json -- go test -count=1 $PKGS
run $ENVV $SRT --settings $S/s-go-ep-readfence.json -- cat $HOME/.srt-row140-decoy
echo "== G4: does the toolchain write to the read-only shared module cache? (find newer files)"
touch $S/marker; sleep 1
run $ENVV $SRT --settings $S/s-go-ep.json -- go test -count=1 $PKGS
run find $MODCACHE -maxdepth 3 -newer $S/marker -print
echo "== G5: GOTOOLCHAIN not pinned (auto) and GOPROXY off: what does the go wrapper need?"
run env -i PATH=/usr/bin:/bin:/opt/homebrew/bin HOME=$EP/home GOCACHE=$EP/gocache GOPROXY=off GOTELEMETRY=off $SRT --settings $S/s-go-ep.json -- go version
run env -i PATH=/usr/bin:/bin:/opt/homebrew/bin HOME=$EP/home GOCACHE=$EP/gocache GOMODCACHE=$MODCACHE GOPROXY=off GOTELEMETRY=off $SRT --settings $S/s-go-ep.json -- go version
