#!/bin/bash
# Q6: can `git status` restricted to the protected pathspecs prove "protected
# paths equal the base commit"? Cost and detection arms on a SCRATCH clone of
# the ailang compiler (git clone --shared; the real checkout is never written).
set -u
SC=/Users/voightkampff/dev/sunholo-data/.scratch-iter277-row156/git
rm -rf $SC; mkdir -p $SC
t0=$(python3 -c 'import time;print(time.time())')
git clone -q --shared /Users/voightkampff/dev/sunholo-data/ailang $SC/ailang
R=$SC/ailang
echo "clone+checkout: $(python3 -c "import time;print(round(time.time()-$t0,2))") s; HEAD $(git -C $R rev-parse --short HEAD); tracked $(git -C $R ls-files | wc -l | tr -d ' ')"
G() { GIT_OPTIONAL_LOCKS=0 git -C $R -c core.fsmonitor=false -c core.quotePath=false "$@"; }
PS=(':(glob,icase)**/*_test.go' ':(glob,icase)**/*.golden' ':(glob,icase)tests/**' ':(glob,icase)examples/tests/**' ':(glob,icase)**/.gitignore')
tm() { local s=$(python3 -c 'import time;print(time.time())'); "$@" > $SC/last.out; local rc=$?; printf 'rc=%d %6.0f ms lines=%s' $rc "$(python3 -c "import time;print((time.time()-$s)*1000)")" "$(tr '\0' '\n' < $SC/last.out | grep -c .)"; }
echo "== cost (3 runs each; first warms)"
for i in 1 2 3; do
  echo "row-184 full status (ignored=no):      $(tm G status --porcelain=v1 -z --untracked-files=all --ignored=no)"
  echo "protected pathspecs, ignored=matching:  $(tm G status --porcelain=v1 -z --untracked-files=all --ignored=matching -- "${PS[@]}")"
done
arm() { # name; then the mutation is already applied
  printf '%-58s full: %s | protected: %s\n' "$1" \
    "$(G status --porcelain=v1 --untracked-files=all --ignored=no | tr '\n' ' ' | cut -c1-60)" \
    "$(G status --porcelain=v1 --untracked-files=all --ignored=matching -- "${PS[@]}" | tr '\n' ' ' | cut -c1-90)"
}
reset() { G checkout -q -- . ; G clean -q -fdx -e '!*' >/dev/null 2>&1; G clean -q -fd; rm -rf $R/zzign; }
T=$(git -C $R ls-files '*_test.go' | head -1); D=$(dirname $T)
echo "== detection arms (test file used: $T)"
reset; arm "control: clean tree"
reset; echo "// PWN" >> $R/$T; arm "modify tracked $T"
reset; rm $R/$T; arm "delete tracked $T"
reset; echo "package x" > $R/$D/zz_new_test.go; arm "create untracked $D/zz_new_test.go"
reset; echo "zz*" >> $R/.gitignore; echo "package x" > $R/$D/zz_ign_test.go; arm "edit .gitignore (zz*) + create $D/zz_ign_test.go"
reset; mkdir -p $R/zzign && echo "/zzign/" >> $R/.git/info/exclude && echo "package x" > $R/zzign/a_test.go; arm "create a_test.go in an info/exclude-ignored dir"
reset; echo "// non-protected change" >> $R/go.mod; arm "control: modify non-protected go.mod"
reset; rm -rf $SC
