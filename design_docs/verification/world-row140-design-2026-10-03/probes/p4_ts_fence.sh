#!/bin/bash
# P4-TS-fence: why node aborts under the read fence, and the configuration that works.
. "$(dirname "$0")/lib.sh"
T=$S/ts-twilight; EP=$S/ts-ep
H=$HOME/.srt-row140-scratch          # a scratch workspace UNDER HOME (as World's are), removed at the end
rm -rf $H; mkdir -p $H/ws $H/ep/tmp $H/ep/home
echo 'console.log("node-ok", process.cwd())' > $H/ws/a.js
cat > $S/s-f1.json <<EOF
{ "filesystem": { "denyRead": ["$HOME"], "allowRead": ["$H/ws", "$H/ep"], "allowWrite": ["$H/ws", "$H/ep"], "denyWrite": [] },
  "network": { "allowedDomains": [], "deniedDomains": [] } }
EOF
cat > $S/s-f2.json <<EOF
{ "filesystem": { "denyRead": ["/private/tmp", "/tmp"], "allowRead": ["$T", "$EP"], "allowWrite": ["$T", "$EP"], "denyWrite": [] },
  "network": { "allowedDomains": [], "deniedDomains": [] } }
EOF
cat > $S/s-f3.json <<EOF
{ "filesystem": { "denyRead": ["$HOME"], "allowRead": ["$T", "$EP"], "allowWrite": ["$T", "$EP"], "denyWrite": [] },
  "network": { "allowedDomains": [], "deniedDomains": [] } }
EOF
ENVV="env -i PATH=/usr/bin:/bin:/opt/homebrew/bin LANG=C"
echo "== F1: workspace under HOME, HOME denied, workspace re-allowed: node"
cd $H/ws
run $ENVV HOME=$H/ep/home CLAUDE_CODE_TMPDIR=$H/ep/tmp $SRT --settings $S/s-f1.json -- node a.js
run $ENVV HOME=$H/ep/home CLAUDE_CODE_TMPDIR=$H/ep/tmp $SRT --settings $S/s-f1.json -- cat $HOME/.srt-row140-decoy
run $ENVV HOME=$H/ep/home CLAUDE_CODE_TMPDIR=$H/ep/tmp $SRT --settings $S/s-f1.json -- ls $HOME
run $ENVV HOME=$H/ep/home CLAUDE_CODE_TMPDIR=$H/ep/tmp $SRT --settings $S/s-f1.json -- /bin/pwd -P
echo "== F2: workspace under /private/tmp, /tmp denied, workspace re-allowed: node"
cd $T
run $ENVV HOME=$EP/home $SRT --settings $S/s-f2.json -- node -e 'console.log("node-ok")'
run $ENVV HOME=$EP/home $SRT --settings $S/s-f2.json -- /bin/pwd -P
echo "== F3: HOME denied only; workspace under /private/tmp: tsc"
run $ENVV HOME=$EP/home CLAUDE_CODE_TMPDIR=$EP/tmp $SRT --settings $S/s-f3.json -- ./node_modules/.bin/tsc --noEmit
run $ENVV HOME=$EP/home CLAUDE_CODE_TMPDIR=$EP/tmp $SRT --settings $S/s-f3.json -- cat $HOME/.srt-row140-decoy
cd $S; rm -rf $H
