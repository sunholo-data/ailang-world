#!/bin/bash
# P7: seeding a per-episode Go build cache from a read-only, operator-prewarmed project cache by APFS clonefile.
. "$(dirname "$0")/lib.sh"
G=$S/go-ailang; SEED=$S/go-ep/gocache; EP2=$S/go-ep2
rm -rf $EP2; mkdir -p $EP2/home $EP2/tmp
run du -sh $SEED
run find $SEED -type f -print -quit
N=$(find $SEED -type f | wc -l); echo "seed files: $N"
t0=$(python3 -c 'import time;print(time.time())')
cp -cR $SEED $EP2/gocache
t1=$(python3 -c 'import time;print(time.time())'); python3 -c "print('clone seed (cp -cR) wall %.2fs' % ($t1-$t0))"
rm -rf $EP2/gocache-copy
t0=$(python3 -c 'import time;print(time.time())')
cp -R $SEED $EP2/gocache-copy
t1=$(python3 -c 'import time;print(time.time())'); python3 -c "print('full copy (cp -R) wall %.2fs' % ($t1-$t0))"
rm -rf $EP2/gocache-copy
cat > $S/s-seed.json <<EOF
{ "filesystem": { "denyRead": [], "allowRead": [], "allowWrite": ["$G", "$EP2"], "denyWrite": [] },
  "network": { "allowedDomains": [], "deniedDomains": [] } }
EOF
cd $G
ENVV="env -i PATH=/usr/bin:/bin:/opt/homebrew/bin HOME=$EP2/home GOCACHE=$EP2/gocache GOMODCACHE=$(go env GOMODCACHE) GOTOOLCHAIN=go1.26.6 GOFLAGS=-mod=readonly GOPROXY=off GOTELEMETRY=off CLAUDE_CODE_TMPDIR=$EP2/tmp"
t0=$(python3 -c 'import time;print(time.time())')
run $ENVV $SRT --settings $S/s-seed.json -- go test -count=1 ./internal/lexer/ ./internal/parser/
t1=$(python3 -c 'import time;print(time.time())'); python3 -c "print('go test on the cloned seed wall %.1fs' % ($t1-$t0))"
echo "== the seed stays byte-unchanged (the episode wrote only its clone)"
run find $SEED -newer $S/marker-ts -type f -print -quit
