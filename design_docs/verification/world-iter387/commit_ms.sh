#!/usr/bin/env bash
# usage: commit_ms.sh <milestone> <msgfile> <test-cmd...>
# tree = HEAD + snapshot(<milestone>); gofmt/build/vet/test must pass; commits the snapshot's files.
set -u
WT=/Users/voightkampff/dev/sunholo-data/.wt-world-iter387-a3a-s4
SNAP=$HOME/.ailang/state/world-iter387/snap/$1; MSG=$(cd "$(dirname "$2")" && pwd)/$(basename "$2"); shift 2
cd "$WT" || exit 2; export AILANG_BIN=$HOME/.pinned-ailang/ailang
[ -d "$SNAP" ] || { echo "no snapshot $SNAP"; exit 2; }
[ -f "$MSG" ] || { echo "no msg $MSG"; exit 2; }
files=$(cd "$SNAP" && find . -type f ! -name MANIFEST.sha256 | sed 's|^\./||' | sort)
for f in $(git status --porcelain --untracked-files=all | awk '{print $2}'); do
  if [ -f "$SNAP/$f" ]; then cp -p "$SNAP/$f" "$f";
  elif git cat-file -e "HEAD:$f" 2>/dev/null; then git show "HEAD:$f" > "$f";
  else case "$f" in design_docs/*) echo "REFUSING: untracked $f is not in the snapshot (commit the plan first)"; exit 2;; esac; rm -f "$f"; fi
done
for f in $files; do mkdir -p "$(dirname "$f")"; cp -p "$SNAP/$f" "$f"; done
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
