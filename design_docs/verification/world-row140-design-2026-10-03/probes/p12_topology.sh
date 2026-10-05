#!/bin/bash
# P12: the World topology — a workspace root with two episode worktrees (one with a .git POINTER FILE);
# the workspace root denied for reads, only the running episode re-allowed; .git pointer file write-denied.
. "$(dirname "$0")/lib.sh"
H=$HOME/.srt-row140-scratch2
rm -rf $H; mkdir -p $H/root/ep1 $H/root/ep2 $H/state/exec-cache/ep1/home $H/state/exec-cache/ep1/tmp
echo "gitdir: /somewhere/.git/worktrees/ep1" > $H/root/ep1/.git
echo "EP2-PRIVATE" > $H/root/ep2/notes.txt
echo "STATE-PRIVATE" > $H/state/world.db
cat > $S/s-topo.json <<EOF
{ "filesystem": {
    "denyRead": ["$HOME", "$H/state", "$H/root"],
    "allowRead": ["$H/root/ep1", "$H/state/exec-cache/ep1"],
    "allowWrite": ["$H/root/ep1", "$H/state/exec-cache/ep1"],
    "denyWrite": ["$H/root/ep1/.git", "$H/root/ep1/.github", "/tmp/claude", "/private/tmp/claude", "$HOME/.npm/_logs", "$HOME/.claude/debug"] },
  "network": { "allowedDomains": [], "deniedDomains": [] } }
EOF
cd $H/root/ep1
ENVV="env -i PATH=/usr/bin:/bin LANG=C HOME=$H/state/exec-cache/ep1/home CLAUDE_CODE_TMPDIR=$H/state/exec-cache/ep1/tmp"
run $ENVV $SRT --settings $S/s-topo.json -- sh -c 'echo ok > a.txt && cat a.txt'
run $ENVV $SRT --settings $S/s-topo.json -- sh -c 'echo x > .git'
run $ENVV $SRT --settings $S/s-topo.json -- sh -c 'echo x > .GIT'
run $ENVV $SRT --settings $S/s-topo.json -- cat ../ep2/notes.txt
run $ENVV $SRT --settings $S/s-topo.json -- ls ..
run $ENVV $SRT --settings $S/s-topo.json -- cat $H/state/world.db
run $ENVV $SRT --settings $S/s-topo.json -- sh -c 'echo x > $HOME/x && echo wrote-ep-home'
run $ENVV $SRT --settings $S/s-topo.json -- sh -c 'echo x > ../ep2/pwn'
run cat $H/root/ep1/.git
cd $S; rm -rf $H
