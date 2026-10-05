#!/bin/bash
# P2: srt network, signals, env, argv handling.
. "$(dirname "$0")/lib.sh"
cd $S/ws
cat > $S/s-net.json <<EOF
{ "filesystem": { "denyRead": [], "allowRead": [], "allowWrite": ["$S/ws"], "denyWrite": [] },
  "network": { "allowedDomains": ["proxy.golang.org"], "deniedDomains": [] } }
EOF
# a loopback server standing in for World's daemon port
python3 -m http.server 17644 --bind 127.0.0.1 >/dev/null 2>&1 &
SRV=$!
sleep 1
echo "== network"
run $SRT --settings $S/s-net.json -- curl -sS -o /dev/null -w '%{http_code}\n' --max-time 10 https://proxy.golang.org/
run $SRT --settings $S/s-net.json -- curl -sS -o /dev/null -w '%{http_code}\n' --max-time 10 https://example.com/
run $SRT --settings $S/s-net.json -- curl -sS -o /dev/null -w '%{http_code}\n' --max-time 5 http://127.0.0.1:17644/
run $SRT --settings $S/s-net.json -- curl -sS --noproxy '*' -o /dev/null -w '%{http_code}\n' --max-time 5 http://127.0.0.1:17644/
run $SRT --settings $S/s-net.json -- curl -sS --noproxy '*' -o /dev/null -w '%{http_code}\n' --max-time 5 https://example.com/
run $SRT --settings $S/s-net.json -- python3 -c 'import socket; s=socket.create_connection(("127.0.0.1",17644),3); print("connected")'
run $SRT --settings $S/s-net.json -- python3 -c 'import socket; s=socket.socket(); s.bind(("127.0.0.1",17699)); s.listen(); print("bound")'
run $SRT --settings $S/s-net.json -- python3 -c 'import socket; print(socket.getaddrinfo("example.com",443)[0][4])'
kill $SRV
echo "== env seen by the child (env -i minimal parent)"
run env -i HOME=$S/home PATH=/usr/bin:/bin:/opt/homebrew/bin LANG=C $SRT --settings $S/s-net.json -- /usr/bin/env
echo "== signals"
run $SRT --settings $S/s-net.json -- sh -c 'kill -TERM $$'
run $SRT --settings $S/s-net.json -- sh -c 'kill -INT $$'
run $SRT --settings $S/s-net.json -- sh -c 'kill -SEGV $$'
echo "== argv handling (srt quotes argv then runs via bash -c)"
run $SRT --settings $S/s-net.json -- printf '[%s]\n' 'a b' '$(echo INJECT)' '`id -un`' ';echo semi' '--' '-c'
run $SRT --settings $S/s-net.json printf '[%s]\n' -c 'echo via-c'
run $SRT --settings $S/s-net.json echo -s /nonexistent.json
echo "== missing settings file"
run $SRT --settings $S/does-not-exist.json -- true
echo "== a deny-read settings file in the sandbox is not loaded from the project"
run ls -la $HOME/.srt-settings.json
