#!/usr/bin/env bash
# usage: commit_ms.sh <milestone> <msgfile> <test-cmd...>
# Makes tree = HEAD + snapshot(<milestone>) for every currently-changed file, runs the test cmd, commits named files.
set -u
WT=/Users/voightkampff/dev/sunholo-data/.wt-world-iter374-row201s2
SNAP=/Users/voightkampff/.ailang/state/world-iter374/snap/$1; MSG=$(cd "$(dirname "$2")" && pwd)/$(basename "$2"); shift 2
cd "$WT" || exit 2; export AILANG_BIN=$HOME/.pinned-ailang/ailang
[ -d "$SNAP" ] || { echo "no snapshot $SNAP"; exit 2; }
files=$(cd "$SNAP" && find . -type f | sed 's|^\./||' | sort)
for f in $(git status --porcelain | awk '{print $2}'); do
  if [ -f "$SNAP/$f" ]; then cp "$SNAP/$f" "$f";
  elif git cat-file -e "HEAD:$f" 2>/dev/null; then git show "HEAD:$f" > "$f";
  else rm -f "$f"; fi
done
for f in $files; do mkdir -p "$(dirname "$f")"; cp "$SNAP/$f" "$f"; done
echo "== tree for $(basename $SNAP):"; git status --porcelain
echo "== gofmt:"; gofmt -l $(echo "$files" | grep '\.go$') ; 
echo "== build+vet:"; go build ./... && go vet $(echo "$files" | grep '\.go$' | xargs -n1 dirname | sort -u | sed 's|^|./|') ; echo "build/vet rc=$?"
echo "== test: $*"; bash -c "$*"; trc=$?; echo "test rc=$trc"
[ $trc -eq 0 ] || { echo "REFUSING to commit: test rc=$trc"; exit 1; }
git add -- $files && git commit -q -F "$MSG" && git log --oneline -1
