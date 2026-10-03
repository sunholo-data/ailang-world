#!/bin/bash
# P1: srt filesystem behaviour — writes, reads (default and deny-then-allow), children, exit codes.
. "$(dirname "$0")/lib.sh"
mkdir -p $S/ws $S/other
cat > $S/s-basic.json <<EOF
{ "filesystem": { "denyRead": [], "allowRead": [], "allowWrite": ["$S/ws"], "denyWrite": [] },
  "network": { "allowedDomains": [], "deniedDomains": [] } }
EOF
cat > $S/s-readfence.json <<EOF
{ "filesystem": { "denyRead": ["$HOME", "/private/tmp", "/tmp"], "allowRead": ["$S/ws"], "allowWrite": ["$S/ws"], "denyWrite": [] },
  "network": { "allowedDomains": [], "deniedDomains": [] } }
EOF
echo "OTHER-PROJECT-FILE" > $S/other/file.txt
cd $S/ws
echo "== version"; run node -e "console.log(require('$SRTDIR/package.json').version)"
echo "== writes (basic: allowWrite=[ws])"
run $SRT --settings $S/s-basic.json -- sh -c 'echo hi > in.txt && echo wrote-in'
run $SRT --settings $S/s-basic.json -- sh -c 'echo x > ../other/outside.txt'
run $SRT --settings $S/s-basic.json -- sh -c 'echo x > $HOME/.srt-row140-write'
run $SRT --settings $S/s-basic.json -- sh -c 'echo x > /tmp/srt-row140-tmpwrite'
run $SRT --settings $S/s-basic.json -- sh -c 'mkdir -p sub && ln -s ../../other sub/lnk && echo y > sub/lnk/viaslink.txt'
ls -la $S/other $HOME/.srt-row140-write /tmp/srt-row140-tmpwrite 2>&1
echo "== reads, default (denyRead=[])"
run $SRT --settings $S/s-basic.json -- cat $HOME/.srt-row140-decoy
run $SRT --settings $S/s-basic.json -- cat $S/other/file.txt
echo "== reads, deny-then-allow (denyRead=[HOME,/private/tmp,/tmp], allowRead=[ws])"
run $SRT --settings $S/s-readfence.json -- cat $HOME/.srt-row140-decoy
run $SRT --settings $S/s-readfence.json -- cat $S/other/file.txt
run $SRT --settings $S/s-readfence.json -- cat in.txt
run $SRT --settings $S/s-readfence.json -- ls $HOME
run $SRT --settings $S/s-readfence.json -- sh -c 'cat /etc/hosts | head -2; ls /usr/bin | head -2'
echo "== child process inheritance"
run $SRT --settings $S/s-basic.json -- sh -c 'sh -c "bash -c \"echo z > ../other/grandchild.txt\""'
run $SRT --settings $S/s-readfence.json -- sh -c 'sh -c "cat $HOME/.srt-row140-decoy"'
run $SRT --settings $S/s-basic.json -- sh -c '( nohup sh -c "sleep 2; echo late > ../other/late.txt" >/dev/null 2>&1 & ); echo parent-exits'
sleep 3; ls -la $S/other
echo "== exit codes"
run $SRT --settings $S/s-basic.json -- false
run $SRT --settings $S/s-basic.json -- sh -c 'exit 42'
run $SRT --settings $S/s-basic.json -- sh -c 'kill -9 $$'
run $SRT --settings $S/s-basic.json -- /nonexistent/binary
