#!/bin/bash
# V18/V19: Git 2.54 wildmatch on multi-byte names + the [r4] widened superset pathspecs.
# One disposable repo per spelling under $S (inside the sprint worktree), fake HOME, no system/global config.
set -u
S=${S:?scratch}; export HOME=$S/home GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null GIT_OPTIONAL_LOCKS=0
mk() { # $1 name, $2 path (printf-escaped), $3 nfd-index(0/1)
  r=$S/r-$1; rm -rf "$r"; mkdir -p "$r"; cd "$r" || exit 1
  git init -q; git config user.email fixture; git config user.name t; git config core.fsmonitor false
  p=$(printf "$2")
  if [ "$3" = 1 ]; then # NFD path recorded in the index as bytes (as a Linux commit would)
    h=$(printf 'x\n' | git hash-object -w --stdin); git update-index --add --cacheinfo 100644 "$h" "$p"
  else mkdir -p "$(dirname "$p")"; printf 'x\n' > "$p"; git add -- "$p"; fi
  git commit -qm base --no-verify
  [ "$3" = 1 ] || printf 'y\n' >> "$p"
}
q() { # print pathspec -> count of ls-files hits and status entries
  for ps in "$@"; do
    a=$(git -c core.quotepath=false ls-files -z -- "$ps" | tr '\0' '\n' | grep -c .)
    b=$(git -c core.quotepath=false status --porcelain=v2 -z --untracked-files=all --ignored=matching -- "$ps" | tr '\0' '\n' | grep -c .)
    echo "  $ps  ls-files=$a status=$b"
  done
}
echo "git $(git --version)"
mk e-nfc 'pkg/\xc3\xa9_test.go' 0; echo "repo e-nfc: pkg/é_test.go (NFC, 2-byte é), modified"
q ':(top,glob,icase)**/*_test.go' ':(top,glob,icase)**/?_test.go' ':(top,glob,icase)**/??_test.go' ':(top,glob,icase)**/[^a]_test.go' ':(top,glob,icase)**/[é]_test.go' ':(top,glob,icase)**/é_test.go'
mk e-nfd 'pkg/e\xcc\x81_test.go' 1; echo "repo e-nfd: pkg/e+U+0301_test.go in the index as NFD bytes (precomposeunicode=$(git config core.precomposeunicode))"
q ':(top,glob,icase)**/é_test.go' ':(top,glob,icase)**/*_test.go' ':(top,glob,icase)**/*_te*t.go'
mk longs 'te\xc5\xbfts/x.py' 0; echo "repo longs: teſts/x.py modified; pattern tests/** widened"
q ':(top,glob,icase)tests' ':(top,glob,icase)tests/**' ':(top,glob,icase)te*t*' ':(top,glob,icase)te*t*/**'
mk kelvin 'pkg/\xe2\x84\xaa_test.go' 0; echo "repo kelvin: pkg/K(U+212A)_test.go modified; pattern *k_test.go widened"
q ':(top,glob,icase)**/*k_test.go' ':(top,glob,icase)**/*_te*t.go'
mk semi 'pkg/a\xcd\xbe_test.go' 0; echo "repo semi: pkg/a(U+037E)_test.go; pattern a;_test.go widened"
q ':(top,glob,icase)**/a;_test.go' ':(top,glob,icase)**/a*_te*t.go'
mk ign 'keep.txt' 0; mkdir -p zzign; printf 'x\n' > zzign/a_test.go; printf 'x\n' > zzign/b.txt; mkdir -p nm; printf 'x\n' > nm/c.txt
printf 'zzign/\nnm/\n' >> .git/info/exclude
echo "repo ign: zzign/a_test.go + nm/c.txt in info/exclude'd dirs"
for ps in ':(top,glob,icase)**/*_te*t.go'; do
  echo "  status matching:  $(git status --porcelain=v2 -z --untracked-files=all --ignored=matching -- "$ps" | tr '\0' ' ')"
  echo "  status no-pathspec matching: $(git status --porcelain=v2 -z --untracked-files=all --ignored=matching | tr '\0' ' ')"
  echo "  ls-files -o -i:   $(git ls-files -z -o -i --exclude-standard -- "$ps" | tr '\0' ' ')"
  echo "  ls-files -o -i no-pathspec: $(git ls-files -z -o -i --exclude-standard | tr '\0' ' ')"
done
# bare-D / replacement / top, under the widened spelling, from a subdirectory CWD
mk repl 'tests/x.py' 0; git checkout -q -- .; rm -rf tests; ln -s keep tests; mkdir -p sub; cd sub
echo "repo repl: tests/ replaced by a symlink, run from sub/"
q ':(top,glob,icase)te*t*' ':(top,glob,icase)te*t*/**' ':(glob,icase)te*t*'
