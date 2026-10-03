#!/bin/bash
# P4-Go-tls: why Go's module download fails TLS under srt on macOS, and what fixes it.
. "$(dirname "$0")/lib.sh"
M=$S/go-tiny; EP=$S/go-tiny-ep
cat > $S/s-gonet-weak.json <<EOF
{ "filesystem": { "denyRead": [], "allowRead": [], "allowWrite": ["$M", "$EP"], "denyWrite": [] },
  "network": { "allowedDomains": ["proxy.golang.org", "sum.golang.org"], "deniedDomains": [] },
  "enableWeakerNetworkIsolation": true }
EOF
ENVV="env -i PATH=/usr/bin:/bin:/opt/homebrew/bin HOME=$EP/home GOCACHE=$EP/gocache GOMODCACHE=$EP/gomod GOTOOLCHAIN=local GOTELEMETRY=off CLAUDE_CODE_TMPDIR=$EP/tmp GOFLAGS=-mod=mod"
cd $M
echo "== T1: SSL_CERT_FILE=/etc/ssl/cert.pem (does Go on darwin honour it?)"
rm -rf $EP/gomod; mkdir -p $EP/gomod
run $ENVV SSL_CERT_FILE=/etc/ssl/cert.pem $SRT --settings $S/s-gonet-proxy.json -- go test -count=1 ./...
echo "== T2: enableWeakerNetworkIsolation (trustd reachable)"
rm -rf $EP/gomod; mkdir -p $EP/gomod
run $ENVV $SRT --settings $S/s-gonet-weak.json -- go test -count=1 ./...
echo "== T3: curl (SecureTransport/LibreSSL) through the same proxy"
run $SRT --settings $S/s-gonet-proxy.json -- curl -sS -o /dev/null -w '%{http_code}\n' https://proxy.golang.org/github.com/google/go-cmp/@v/list
echo "== T4: offline: populate the module cache OUTSIDE the sandbox (operator step), then run offline inside"
rm -rf $EP/gomod; mkdir -p $EP/gomod
run env GOMODCACHE=$EP/gomod GOFLAGS=-mod=mod go -C $M mod download
run $ENVV GOFLAGS=-mod=readonly GOPROXY=off $SRT --settings $S/s-gonet-none.json -- go test -count=1 ./...
