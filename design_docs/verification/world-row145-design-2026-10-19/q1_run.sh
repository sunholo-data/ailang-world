#!/bin/sh
# Q1 (row 145 design): srt 0.0.78 loopback behaviour, allowLocalBinding false vs true.
# Settings mirror World's rendering (exec_settings.go) on a scratch dir; only the flag differs.
set -u
D=$(cd "$(dirname "$0")" && pwd)
S="$D/scratch"
mkdir -p "$S"
CLI=/Users/voightkampff/.ailang-world-exec/srt/node_modules/@anthropic-ai/sandbox-runtime/dist/cli.js
for bind in false true; do
  cat > "$S/settings-$bind.json" <<EOF
{"filesystem":{"denyRead":[],"allowRead":["$D"],"allowWrite":["$S"],"denyWrite":[]},
 "network":{"allowedDomains":[],"deniedDomains":[],"allowLocalBinding":$bind,"allowAllUnixSockets":false},
 "enableWeakerNetworkIsolation":false,"allowAppleEvents":false}
EOF
  echo "== allowLocalBinding=$bind"
  (cd "$S" && /opt/homebrew/bin/node "$CLI" --settings "$S/settings-$bind.json" -- /usr/bin/python3 "$D/q1_bind.py" 2>&1)
done
rm -rf "$S"
