#!/bin/bash
# Q8 (quorum r1): (a) can sandboxed code set index flags or move HEAD in a
# `git worktree add` episode under World's srt rendering? (b) which per-call
# check sees an assume-unchanged / skip-worktree tamper, and at what cost on the
# ailang repo (scratch clone)? Scratch is a SIBLING dir, never /tmp.
set -u
SC=/Users/voightkampff/dev/sunholo-data/.scratch-iter277-row156-r2; rm -rf $SC; mkdir -p $SC
SRT=$HOME/.ailang-world-exec/srt/node_modules/@anthropic-ai/sandbox-runtime/dist/cli.js
mk() { # fresh source repo + linked worktree $SC/root/$1
  rm -rf $SC/src $SC/root $SC/state; mkdir -p $SC/src/pkg $SC/root $SC/state/exec-cache/$1/home
  ( cd $SC/src && git init -q -b main && git config user.email t@x && git config user.name t &&
    printf 'package pkg\n' > pkg/a.go && printf 'package pkg\n// ORIG\n' > pkg/a_test.go && git add . && git commit -qm base &&
    git worktree add -q --detach $SC/root/$1 )
  WT=$(cd $SC/root/$1 && pwd -P); CACHE=$(cd $SC/state/exec-cache/$1 && pwd -P)
}
settings() { python3 - "$WT" "$CACHE" "$HOME" > $SC/s.json <<'PY'
import json,sys
wt,cache,home=sys.argv[1:4]
deny=[f"{wt}/{n}" for n in [".git/hooks",".git/config",".git",".github",".pi",".claude",".ailang",".gitmodules",".gitattributes"]]
deny+=["/tmp/claude","/private/tmp/claude",f"{home}/.npm/_logs",f"{home}/.claude/debug"]
print(json.dumps({"filesystem":{"denyRead":[],"allowRead":[],"allowWrite":[wt,cache],"denyWrite":deny},
 "network":{"allowedDomains":[],"deniedDomains":[],"allowLocalBinding":False,"allowAllUnixSockets":False},
 "enableWeakerNetworkIsolation":False,"allowAppleEvents":False}))
PY
chmod 600 $SC/s.json; }
ENVV=(/usr/bin/env -i HOME=$SC/home PATH=/opt/homebrew/bin:/usr/bin:/bin GIT_CONFIG_NOSYSTEM=1 GIT_AUTHOR_NAME=a GIT_AUTHOR_EMAIL=a@x GIT_COMMITTER_NAME=a GIT_COMMITTER_EMAIL=a@x LANG=C)
mkdir -p $SC/home
run() { # mode cmd...
  if [ "$1" = srt ]; then shift; (cd $WT && node $SRT --settings $SC/s.json -- "${ENVV[@]}" "$@" 2>&1 | tail -1); else shift; (cd $WT && "${ENVV[@]}" "$@" 2>&1 | tail -1); fi
}
echo "# (a) index flags / HEAD under srt 0.0.78 ($(shasum -a 256 $SRT | cut -c1-12)), World's allowWrite=[worktree,cache], denyWrite incl <wt>/.git; denyRead EMPTY (more permissive than World's fence)"
for mode in plain srt; do
  for arm in skip-worktree assume-unchanged commit; do
    mk ep1; settings
    echo "worktree .git is a $( [ -f $WT/.git ] && echo 'POINTER FILE' || echo dir ): $(cat $WT/.git 2>/dev/null | head -1)"
    HEAD0=$(git -C $WT rev-parse HEAD)
    case $arm in
      skip-worktree) out=$(run $mode git update-index --skip-worktree pkg/a_test.go);;
      assume-unchanged) out=$(run $mode git update-index --assume-unchanged pkg/a_test.go);;
      commit) echo "// edit" >> $WT/pkg/a.go; out=$(run $mode git commit -qam x);;
    esac
    printf '[%-5s] %-17s rc-tail: %-70s host: ls-files -v=%s HEAD %s\n' "$mode" "$arm" "${out:-<no output>}" \
      "$(git -C $WT ls-files -v pkg/a_test.go | cut -c1)" "$( [ "$(git -C $WT rev-parse HEAD)" = "$HEAD0" ] && echo unchanged || echo MOVED)"
  done
done
echo
echo "# (b) detection: set assume-unchanged / skip-worktree UNSANDBOXED (an attacker who got there), then modify the file"
mk ep2
PS=(':(glob,icase)**/*_test.go')
for flag in assume-unchanged skip-worktree; do
  git -C $WT checkout -q -- . ; git -C $WT update-index --no-assume-unchanged --no-skip-worktree pkg/a_test.go
  git -C $WT update-index --$flag pkg/a_test.go; echo "// PWN" >> $WT/pkg/a_test.go
  echo "-- --$flag + modified pkg/a_test.go"
  echo "   status --ignored=matching : [$(git -C $WT status --porcelain=v1 --untracked-files=all --ignored=matching -- "${PS[@]}" | tr '\n' ' ')]"
  echo "   diff-index --cached base  : [$(git -C $WT diff-index --cached --name-status HEAD -- "${PS[@]}" | tr '\n\t' '  ')]"
  echo "   diff-files                : [$(git -C $WT diff-files --name-status -- "${PS[@]}" | tr '\n\t' '  ')]"
  echo "   ls-files -v               : [$(git -C $WT ls-files -v -- "${PS[@]}" | tr '\n' ' ')]"
done
echo
echo "# (c) cost on a scratch clone of ailang (28 699 tracked), protected pathspecs as V9 (warm, 3 runs)"
git clone -q --shared /Users/voightkampff/dev/sunholo-data/ailang $SC/ail
P2=(':(glob,icase)**/*_test.go' ':(glob,icase)**/*.golden' ':(glob,icase)tests/**' ':(glob,icase)examples/tests/**')
B=$(git -C $SC/ail rev-parse HEAD)
G() { GIT_OPTIONAL_LOCKS=0 git -C $SC/ail -c core.fsmonitor=false "$@"; }
tm() { local s=$(python3 -c 'import time;print(time.time())'); "$@" >/dev/null; printf '%5.0f ms' "$(python3 -c "import time;print((time.time()-$s)*1000)")"; }
G status >/dev/null 2>&1
for i in 1 2 3; do
  echo "status v2 --ignored=matching: $(tm G status --porcelain=v2 -z --branch --untracked-files=all --ignored=matching -- "${P2[@]}") | ls-files -v: $(tm G ls-files -v -z -- "${P2[@]}") | diff-index --cached base + diff-files: $(tm G diff-index --cached -z $B -- "${P2[@]}") + $(tm G diff-files -z -- "${P2[@]}")"
done
rm -rf $SC
