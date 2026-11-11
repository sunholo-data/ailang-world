#!/bin/sh
# Run as: sh generate.sh <an empty directory>, with WORLD_CHROME set to Chrome's binary and node on
# PATH. Re-records chrome-es256.json and chrome-eddsa.json through the passkey page (row 201 A2 M7)
# with Chrome's virtual authenticator (a test key, D-WORLD-109): cmd/ailang-worldd's
# TestPasskeyPageRecordsTheFixtures, driving cmd/ailang-worldd/testdata/passkey_cdp.mjs. It writes
# only into that directory. Copying the files over this directory's is a re-pin: XW0's SHA-256s and
# every golden value taken from them change in the same commit, and SOURCE.txt names the recording.
set -eu
if [ "$#" -ne 1 ]; then
	echo "usage: sh generate.sh <an empty directory>" >&2
	exit 2
fi
if [ ! -d "$1" ] || [ -n "$(ls -A "$1")" ]; then
	echo "generate.sh: $1 is not an empty directory" >&2
	exit 2
fi
if [ -z "${WORLD_CHROME:-}" ]; then
	echo "generate.sh: set WORLD_CHROME to Chrome's binary (only its virtual authenticator is used)" >&2
	exit 2
fi
out=$(cd "$1" && pwd)
cd "$(dirname "$0")/../../../.."
WORLD_PASSKEY_RECORD_DIR="$out" WORLD_REQUIRE_PASSKEY_BROWSER=1 go test ./cmd/ailang-worldd/ -run '^TestPasskeyPageRecordsTheFixtures$' -count=1 -v
ls "$out/chrome-es256.json" "$out/chrome-eddsa.json"
