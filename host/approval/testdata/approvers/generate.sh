#!/usr/bin/env bash
# Regenerates the approver refusal fixtures (row 201 A1b, test AP1) with REAL ssh-keygen. Test
# keys only: made in a temp dir with HOME pointed at it and no agent, deleted on exit; only the
# two public keys and SOURCE.txt are kept (text files with no NUL byte: the tracked-binary gate).
# After a run, update approverFixtureSHA256 in ../../approver_test.go IN THE SAME COMMIT (AP0).
set -euo pipefail
here=$(cd "$(dirname "$0")" && pwd)
tmp=$(mktemp -d); trap 'rm -rf "$tmp"' EXIT
export HOME="$tmp"; unset SSH_AUTH_SOCK
ssh-keygen -q -t rsa -b 3072 -N '' -C world-approver-fixture-rsa3072 -f "$tmp/rsa3072" < /dev/null
ssh-keygen -q -t ecdsa -b 384 -N '' -C world-approver-fixture-p384 -f "$tmp/p384" < /dev/null
cp "$tmp/rsa3072.pub" "$tmp/p384.pub" "$here/"
{ ssh -V 2>&1; date -u +%FT%TZ; ssh-keygen -l -f "$here/rsa3072.pub"; ssh-keygen -l -f "$here/p384.pub"; } > "$here/SOURCE.txt"
echo "regenerated; now update approverFixtureSHA256 in approver_test.go: shasum -a 256 $here/*" >&2
