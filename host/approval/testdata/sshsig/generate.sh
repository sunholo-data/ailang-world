#!/usr/bin/env bash
# Regenerates the SSHSIG fixtures with REAL ssh-keygen (row 201 A1a). Test keys only:
# they are made in a temp dir and deleted on exit; only public keys and signatures are kept.
# The 90-byte message is kept as message.hex (text), never as a binary file: the tracked-binary
# gate (scripts/tracked_binary_gate.sh) rejects any committed blob that holds a NUL byte.
# After a run, update fixtureSHA256 in ../../sshsig_test.go IN THE SAME COMMIT (test SG1).
set -euo pipefail
here=$(cd "$(dirname "$0")" && pwd)
hex=$(sed -n 's/^const fixtureMessageHex = "\([0-9a-f]*\)"$/\1/p' "$here/../../sshsig_test.go")
[ ${#hex} -eq 180 ] || { echo "fixtureMessageHex not found" >&2; exit 1; }
tmp=$(mktemp -d); trap 'rm -rf "$tmp"' EXIT
printf '%s\n' "$hex" > "$here/message.hex"
printf '%s' "$hex" | xxd -r -p > "$tmp/message.bin"
ssh-keygen -q -t ed25519 -N '' -C world-approval-fixture-ed25519 -f "$tmp/ed25519"
ssh-keygen -q -t ecdsa -b 256 -N '' -C world-approval-fixture-p256 -f "$tmp/p256"
ssh-keygen -Y sign -n world-approval@sunholo -f "$tmp/ed25519" < "$tmp/message.bin" > "$here/ed25519.sig"
ssh-keygen -Y sign -n world-approval@sunholo -f "$tmp/p256" < "$tmp/message.bin" > "$here/p256.sig"
ssh-keygen -Y sign -n git -f "$tmp/ed25519" < "$tmp/message.bin" > "$here/git-ns-ed25519.sig"
cp "$tmp/ed25519.pub" "$tmp/p256.pub" "$here/"
{ ssh -V 2>&1; date -u +%FT%TZ; ssh-keygen -l -f "$here/ed25519.pub"; ssh-keygen -l -f "$here/p256.pub"; } > "$here/SOURCE.txt"
echo "regenerated; now update fixtureSHA256 in sshsig_test.go: shasum -a 256 $here/*" >&2
