#!/bin/bash
# Run from the executor worktree root; all writes stay under this root.
set -euo pipefail
ROOT=$(pwd -P)
E="$ROOT/verification/world-iter335/executor"
export GOTOOLCHAIN=local
module_seed=$(go env GOMODCACHE)
export GOCACHE="$E/go-build" GOMODCACHE="$E/modules" GOPATH="$E/go"
export TMPDIR="$E/tmp" GOTMPDIR="$E/tmp"
export XDG_CACHE_HOME="$E/cache" XDG_CONFIG_HOME="$E/config" XDG_STATE_HOME="$E/state"
export AILANG_STATE_DIR="$E/ailang-state" AILANG_CACHE_DIR="$E/ailang-cache"
export GOPROXY=off GOSUMDB=off GOFLAGS=-mod=readonly
mkdir -p "$E" "$GOCACHE" "$GOMODCACHE" "$GOPATH" "$TMPDIR"
mkdir -p "$XDG_CACHE_HOME" "$XDG_CONFIG_HOME" "$XDG_STATE_HOME" "$AILANG_STATE_DIR" "$AILANG_CACHE_DIR"
# Clone graph/download metadata, then only the source modules in go.mod.
# The source cache is READ ONLY; cp -cR uses macOS copy-on-write clones.
if [ ! -d "$GOMODCACHE/cache" ]; then
  cp -cR "$module_seed/cache" "$GOMODCACHE/cache"
fi
awk '$1 ~ /^(github[.]com|golang[.]org|modernc[.]org)\// {print $1, $2}' go.mod |
while read -r mod ver; do
  mkdir -p "$GOMODCACHE/$(dirname "$mod")"
  if [ ! -d "$GOMODCACHE/$mod@$ver" ]; then
    cp -cR "$module_seed/$mod@$ver" "$GOMODCACHE/$mod@$ver"
  fi
done
