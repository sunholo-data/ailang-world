#!/bin/bash
# P6: argument / config injection — what each vector buys INSIDE the fence (writes confined; reads fenced or not).
. "$(dirname "$0")/lib.sh"
M=$S/go-tiny; EP=$S/go-tiny-ep; T=$S/ts-twilight; Y=$S/py-cli
GOENVV="env -i PATH=/usr/bin:/bin:/opt/homebrew/bin HOME=$EP/home GOCACHE=$EP/gocache GOMODCACHE=$EP/gomod GOTOOLCHAIN=local GOTELEMETRY=off CLAUDE_CODE_TMPDIR=$EP/tmp GOFLAGS=-mod=readonly GOPROXY=off"
cat > $S/s-inj-go.json <<EOF
{ "filesystem": { "denyRead": [], "allowRead": [], "allowWrite": ["$M", "$EP"], "denyWrite": [] },
  "network": { "allowedDomains": [], "deniedDomains": [] } }
EOF
cat > $S/s-inj-go-fence.json <<EOF
{ "filesystem": { "denyRead": ["$HOME"], "allowRead": ["$M", "$EP"], "allowWrite": ["$M", "$EP"], "denyWrite": [] },
  "network": { "allowedDomains": [], "deniedDomains": [] } }
EOF
cd $M
echo "== I1: go test -exec runs an arbitrary program (here sh) as the test runner"
run $GOENVV $SRT --settings $S/s-inj-go.json -- go test -count=1 -exec "sh -c 'cat $HOME/.srt-row140-decoy; echo x > /tmp/srt-row140-i1; echo; exec \"\$0\" \"\$@\"'" ./...
run $GOENVV $SRT --settings $S/s-inj-go-fence.json -- go test -count=1 -exec "sh -c 'cat $HOME/.srt-row140-decoy; exec \"\$0\" \"\$@\"'" ./...
echo "== I2: go test -o / -coverprofile to an absolute path outside the worktree"
run $GOENVV $SRT --settings $S/s-inj-go.json -- go test -count=1 -coverprofile=/tmp/srt-row140-cover.out ./...
echo "== I3: a go.mod toolchain line asking for a newer toolchain, GOTOOLCHAIN pinned by the operator env vs auto"
cp go.mod go.mod.bak; printf '\ntoolchain go1.99.0\n' >> go.mod
run $GOENVV $SRT --settings $S/s-inj-go.json -- go test -count=1 ./...
run env -i PATH=/usr/bin:/bin:/opt/homebrew/bin HOME=$EP/home GOCACHE=$EP/gocache GOMODCACHE=$EP/gomod GOTOOLCHAIN=go1.26.6 GOPROXY=off $SRT --settings $S/s-inj-go.json -- go version
cp go.mod.bak go.mod; rm go.mod.bak
echo "== I4: npm --prefix points npm at ANOTHER project (its scripts run; writes still fenced)"
mkdir -p $S/other-npm; cat > $S/other-npm/package.json <<'EOF'
{ "name": "other", "version": "0.0.0", "scripts": { "test": "echo OTHER-PROJECT-SCRIPT-RAN && echo x > ./wrote-by-other.txt" } }
EOF
cat > $S/s-inj-ts.json <<EOF
{ "filesystem": { "denyRead": [], "allowRead": [], "allowWrite": ["$T", "$S/ts-ep"], "denyWrite": [] },
  "network": { "allowedDomains": [], "deniedDomains": [] } }
EOF
cd $T
run env -i PATH=/usr/bin:/bin:/opt/homebrew/bin HOME=$S/ts-ep/home npm_config_cache=$S/ts-ep/npm-cache $SRT --settings $S/s-inj-ts.json -- npm --prefix $S/other-npm test
echo "== I5: a worktree .npmrc (project config) injecting node-options --require (code runs inside the fence)"
echo 'console.log("NPMRC-REQUIRE-RAN"); try { require("fs").writeFileSync("/tmp/srt-row140-npmrc","x"); console.log("WROTE") } catch (e) { console.log("REFUSED", e.code) }' > $T/zz-evil.cjs
echo "node-options=--require ./zz-evil.cjs" > $T/.npmrc
run env -i PATH=/usr/bin:/bin:/opt/homebrew/bin HOME=$S/ts-ep/home npm_config_cache=$S/ts-ep/npm-cache $SRT --settings $S/s-inj-ts.json -- npm test -- run tests/activityDiscovery.test.ts
rm -f $T/.npmrc $T/zz-evil.cjs
echo "== I6: pytest -p loads a plugin module named on the command line; -c points at another config; --basetemp outside"
cd $Y
run env -i PATH=/usr/bin:/bin LANG=C.UTF-8 HOME=$S/py-ep/home $SRT --settings $S/s-py.json -- .venv/bin/python -m pytest -q -p no:cacheprovider -p this_plugin_does_not_exist tests/test_cli_access.py
run env -i PATH=/usr/bin:/bin LANG=C.UTF-8 HOME=$S/py-ep/home $SRT --settings $S/s-py.json -- .venv/bin/python -m pytest -q -p no:cacheprovider --basetemp=/tmp/srt-row140-basetemp tests/test_cli_access.py
ls -la /tmp/srt-row140-i1 /tmp/srt-row140-cover.out /tmp/srt-row140-npmrc /tmp/srt-row140-basetemp $S/other-npm/wrote-by-other.txt 2>&1
