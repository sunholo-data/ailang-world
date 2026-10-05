#!/bin/bash
D=$HOME/dev/sunholo-data/ailang/tools/pi-extensions/sandbox/node_modules/@anthropic-ai/sandbox-runtime/dist
grep -n "tmp/claude\|TMPDIR\|getDefaultWritePaths\|DANGEROUS\|dangerous\|hooks\|bashrc\|mandatory" $D/sandbox/sandbox-utils.js $D/sandbox/macos-sandbox-utils.js $D/sandbox/sandbox-manager.js | head -60
