#!/usr/bin/env bash
# usage: commit_ms.sh <milestone> <msgfile> <test-cmd...>
# tree = HEAD + snapshot(<milestone>) - DELETED; gofmt/build/vet/test must pass; commits the snapshot's files.
set -u
WT=/Users/voightkampff/dev/sunholo-data/.wt-world-iter398-row243-v2b
SNAP=$HOME/.ailang/state/world-iter398/snap/$1; MSG=$(cd "$(dirname "$2")" && pwd)/$(basename "$2"); shift 2
cd "$WT" || exit 2; export AILANG_BIN=$HOME/.pinned-ailang/ailang
[ -d "$SNAP" ] || { echo "no snapshot $SNAP"; exit 2; }
[ -f "$MSG" ] || { echo "no msg $MSG"; exit 2; }
[ -f "$SNAP/DELETED" ] || { echo "no DELETED list in $SNAP"; exit 2; }
(cd "$SNAP" && shasum -a 256 -c MANIFEST.sha256 >/dev/null) || { echo "REFUSING: snapshot MANIFEST mismatch"; exit 2; }
files=$(cd "$SNAP" && find . -type f ! -path ./MANIFEST.sha256 ! -path ./STATUS ! -path ./DELETED | sed 's|^\./||' | sort)
dels=$(grep -v '^$' "$SNAP/DELETED" || true)
# reset every working-tree change not in this snapshot back to HEAD
git status --porcelain --untracked-files=all | while read -r st f; do
  if [ -f "$SNAP/$f" ]; then :;
  elif echo "$dels" | grep -qxF "$f"; then :;
  elif git cat-file -e "HEAD:$f" 2>/dev/null; then mkdir -p "$(dirname "$f")"; git show "HEAD:$f" > "$f";
  else case "$f" in design_docs/*) echo "REFUSING: untracked $f is not in the snapshot"; exit 2;; esac; rm -f "$f"; fi
done || exit 2
for f in $files; do mkdir -p "$(dirname "$f")"; cp -p "$SNAP/$f" "$f"; done
for f in $dels; do if git cat-file -e "HEAD:$f" 2>/dev/null; then git rm -q --cached -- "$f"; fi; rm -f "$f"; done
echo "== tree for $(basename "$SNAP"):"; git status --porcelain --untracked-files=all
gofiles=$(echo "$files" | grep '\.go$' || true)
if [ -n "$gofiles" ]; then
  bad=$(gofmt -l $gofiles); [ -z "$bad" ] || { echo "REFUSING: gofmt: $bad"; exit 1; }
  go build ./... || { echo "REFUSING: build"; exit 1; }
  go vet $(echo "$gofiles" | xargs -n1 dirname | sort -u | sed 's|^|./|') || { echo "REFUSING: vet"; exit 1; }
fi
echo "== test: $*"; bash -c "$*"; trc=$?; echo "test rc=$trc"
[ $trc -eq 0 ] || { echo "REFUSING to commit: test rc=$trc"; exit 1; }
git add -- $files && git commit -q -F "$MSG" && git log --oneline -1
echo "NUMSTAT $(git diff -M --numstat c09e0cdb HEAD -- cmd host docs | awk '{s+=$1+$2} END {print s}')"
git diff -M --numstat c09e0cdb HEAD -- cmd host docs | grep '=>' || true
