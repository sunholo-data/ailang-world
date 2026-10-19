#!/bin/sh
# Q4 (row 145 design): how long an APFS copy-on-write clone (cp -c) of a real episode
# worktree takes, as the premise for running a long command on a snapshot.
# Source is read only; the clone goes to a scratch dir that is removed afterwards.
set -eu
SRC=${1:-/Users/voightkampff/.ailang/state/trial-2026-10-16/ws/w190}
D=$(cd "$(dirname "$0")" && pwd)
S="$D/scratch-clone"
rm -rf "$S"; mkdir -p "$S"
echo "source: $(du -sh "$SRC" | cut -f1), $(find "$SRC" -type f | wc -l | tr -d ' ') files"
for i in 1 2 3; do
  /usr/bin/time -p cp -c -R "$SRC" "$S/run$i" 2>&1 | grep real
done
echo "clone files: $(find "$S/run1" -type f | wc -l | tr -d ' ')"
rm -rf "$S"
