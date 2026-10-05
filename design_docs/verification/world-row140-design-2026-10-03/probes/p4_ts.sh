#!/bin/bash
# P4-TS: tsc --noEmit and vitest on a scratch APFS clone of TwilightGame (no .git, dist, public) under srt.
. "$(dirname "$0")/lib.sh"
SRC=$HOME/dev/markedmondson1234/TwilightGame
T=$S/ts-twilight; EP=$S/ts-ep
if [ ! -d $T ]; then
  mkdir -p $T
  for e in $(ls -A $SRC); do
    case "$e" in .git|dist|public|output|.claude|incoming-art) continue;; esac
    cp -cR "$SRC/$e" "$T/"
  done
fi
rm -rf $EP; mkdir -p $EP/home $EP/tmp $EP/npm-cache
NODEBIN=$(dirname $(command -v node))
cat > $S/s-ts.json <<EOF
{ "filesystem": { "denyRead": [], "allowRead": [], "allowWrite": ["$T", "$EP"], "denyWrite": [] },
  "network": { "allowedDomains": [], "deniedDomains": [] } }
EOF
cat > $S/s-ts-readfence.json <<EOF
{ "filesystem": { "denyRead": ["$HOME", "/private/tmp", "/tmp"], "allowRead": ["$T", "$EP"], "allowWrite": ["$T", "$EP"], "denyWrite": [] },
  "network": { "allowedDomains": [], "deniedDomains": [] } }
EOF
cd $T
ENVV="env -i PATH=/usr/bin:/bin:$NODEBIN HOME=$EP/home npm_config_cache=$EP/npm-cache CLAUDE_CODE_TMPDIR=$EP/tmp CI=1 NO_COLOR=1"
echo "node=$(command -v node) $(node --version)"
echo "== TS1: tsc --noEmit via the local binary (no npx, no npm)"
t0=$(python3 -c 'import time;print(time.time())')
run $ENVV $SRT --settings $S/s-ts.json -- ./node_modules/.bin/tsc --noEmit
t1=$(python3 -c 'import time;print(time.time())'); python3 -c "print('TS1 wall %.1fs' % ($t1-$t0))"
echo "== TS1b: unsandboxed baseline"
t0=$(python3 -c 'import time;print(time.time())')
run $ENVV ./node_modules/.bin/tsc --noEmit
t1=$(python3 -c 'import time;print(time.time())'); python3 -c "print('TS1b wall %.1fs' % ($t1-$t0))"
echo "== TS2: vitest run (one file) via the local binary"
F=$(ls tests/*.test.ts 2>/dev/null | head -1); echo "file=$F"
t0=$(python3 -c 'import time;print(time.time())')
run $ENVV $SRT --settings $S/s-ts.json -- ./node_modules/.bin/vitest run $F
t1=$(python3 -c 'import time;print(time.time())'); python3 -c "print('TS2 wall %.1fs' % ($t1-$t0))"
echo "== TS3: npm test -- run (npm itself, per-episode npm cache)"
run $ENVV $SRT --settings $S/s-ts.json -- npm test -- run $F
echo "== TS4: the read fence (HOME denied) with node from $NODEBIN"
run $ENVV $SRT --settings $S/s-ts-readfence.json -- ./node_modules/.bin/tsc --noEmit
echo "== TS5: writes outside the worktree by the toolchain? (files newer than the marker under HOME/.npm, ~/Library/Caches)"
touch $S/marker-ts; sleep 1
run $ENVV $SRT --settings $S/s-ts.json -- ./node_modules/.bin/vitest run $F
run find $HOME/.npm $HOME/Library/Caches -maxdepth 2 -newer $S/marker-ts -print
run find $T -newer $S/marker-ts -not -path '*/node_modules/*' -print
run find $T/node_modules -maxdepth 2 -newer $S/marker-ts -print
