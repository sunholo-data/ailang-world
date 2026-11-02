#!/bin/bash
# Controller only, outside the Codex sandbox, with the workspace-local
# cache/temp/state environment above and its existing write containment.
export AILANG_BIN="$HOME/.pinned-ailang/ailang"
"$AILANG_BIN" --version
# Require released AILANG v0.41.0; do not install/replace anything under HOME.
go test -count=1 ./host/daemon/
go test -count=1 ./host/broker/
go build ./... && go test ./...
./scripts/verify_ail.sh
