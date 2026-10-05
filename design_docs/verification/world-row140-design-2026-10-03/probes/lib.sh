# shared probe helpers (row 140 design measurements)
S=${S:?set S to a scratch directory}
SRTDIR=$HOME/dev/sunholo-data/ailang/tools/pi-extensions/sandbox/node_modules/@anthropic-ai/sandbox-runtime
SRT="/opt/homebrew/bin/node $SRTDIR/dist/cli.js"
# run: print the command, run it, print rc
run() { echo "\$ $*"; "$@" 2>&1; echo "[rc=$?]"; echo; }
