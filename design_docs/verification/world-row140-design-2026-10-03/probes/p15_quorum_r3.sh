#!/bin/bash
# P15 (quorum round 3 premises):
#  (a) gemini: env given to the HOST node that runs srt (e.g. NODE_OPTIONS) executes before the sandbox exists.
#  (b) the fix: host node gets a fixed minimal env; the child env is set INSIDE the sandbox by `env -i` in the trampoline.
#  (c) glm: does the child get a working TMPDIR in the exec cache? (go t.TempDir, python tempfile/pytest tmp_path, node os.tmpdir)
#  (d) kimi: srt run from an archive copy that lives under a read-DENIED state dir.
. "$(dirname "$0")/lib.sh"
H=$HOME/.srt-row140-scratch4
rm -rf $H; mkdir -p $H/root/ep1 $H/state/exec-cache/ep1/home $H/state/exec-cache/ep1/tmp $H/state/archive
cp -R $SRTDIR $H/state/archive/srt
cat > $S/s-r3.json <<EOF
{ "filesystem": { "denyRead": ["$HOME", "$H/state", "$H/root"], "allowRead": ["$H/root/ep1", "$H/state/exec-cache/ep1"],
    "allowWrite": ["$H/root/ep1", "$H/state/exec-cache/ep1"],
    "denyWrite": ["/tmp/claude", "/private/tmp/claude", "$HOME/.npm/_logs", "$HOME/.claude/debug"] },
  "network": { "allowedDomains": [], "deniedDomains": [] } }
EOF
cd $H/root/ep1
echo 'require("fs").writeFileSync(process.env.HOME_REAL + "/.srt-row140-hostnode-escape", "x"); console.error("HOST-NODE-REQUIRE-RAN")' > evil.cjs
echo "== (a) NODE_OPTIONS in the env of the host node that runs srt"
run env -i PATH=/usr/bin:/bin HOME_REAL=$HOME NODE_OPTIONS="--require ./evil.cjs" $SRT --settings $S/s-r3.json -- /usr/bin/true
run ls -la $HOME/.srt-row140-hostnode-escape
rm -f $HOME/.srt-row140-hostnode-escape
echo "== (b) host node with a minimal env; the child env set inside the sandbox by env -i"
TRAMP='exec /usr/bin/env -i "$@"'
CH="HOME=$H/state/exec-cache/ep1/home TMPDIR=$H/state/exec-cache/ep1/tmp PATH=/usr/bin:/bin:/opt/homebrew/bin LANG=C.UTF-8 NODE_OPTIONS=--require=./evil.cjs HOME_REAL=$HOME"
run env -i PATH=/usr/bin:/bin $SRT --settings $S/s-r3.json -- /bin/sh -c "$TRAMP"' ; exit $?' world-exec $CH /usr/bin/env
run env -i PATH=/usr/bin:/bin $SRT --settings $S/s-r3.json -- /bin/sh -c "$TRAMP"' ; exit $?' world-exec $CH /opt/homebrew/bin/node -e 'console.log("child node ran")'
run ls -la $HOME/.srt-row140-hostnode-escape
rm -f evil.cjs
CH2="HOME=$H/state/exec-cache/ep1/home TMPDIR=$H/state/exec-cache/ep1/tmp PATH=/usr/bin:/bin:/opt/homebrew/bin LANG=C.UTF-8"
echo "== (c) temp dirs inside the fence"
mkdir -p $H/root/ep1/gotmp && cd $H/root/ep1/gotmp
printf 'module example.com/t\n\ngo 1.26\n' > go.mod
printf 'package t\nimport ("os";"testing")\nfunc TestTmp(t *testing.T){d:=t.TempDir(); if err:=os.WriteFile(d+"/x",[]byte("x"),0o600);err!=nil{t.Fatal(err)}; t.Log("TEMPDIR",d)}\n' > t_test.go
run env -i PATH=/usr/bin:/bin $SRT --settings $S/s-r3.json -- /bin/sh -c "$TRAMP"' ; exit $?' world-exec $CH2 GOCACHE=$H/state/exec-cache/ep1/gocache GOTOOLCHAIN=local GOPROXY=off GOFLAGS=-mod=mod /opt/homebrew/bin/go test -count=1 -v ./...
cd $H/root/ep1
run env -i PATH=/usr/bin:/bin $SRT --settings $S/s-r3.json -- /bin/sh -c "$TRAMP"' ; exit $?' world-exec $CH2 /usr/bin/python3 -c 'import tempfile;f=tempfile.NamedTemporaryFile(delete=False);f.write(b"x");print("PYTMP",f.name)'
run env -i PATH=/usr/bin:/bin $SRT --settings $S/s-r3.json -- /bin/sh -c "$TRAMP"' ; exit $?' world-exec $CH2 /opt/homebrew/bin/node -e 'const fs=require("fs"),os=require("os");const d=fs.mkdtempSync(os.tmpdir()+"/n-");fs.writeFileSync(d+"/x","x");console.log("NODETMP",d)'
echo "== (d) srt run FROM an archive copy under the read-denied state dir"
run env -i PATH=/usr/bin:/bin /opt/homebrew/bin/node $H/state/archive/srt/dist/cli.js --settings $S/s-r3.json -- /bin/sh -c 'echo ran-from-archive; cat '"$H"'/state/archive/srt/package.json'
cd $S; rm -rf $H
