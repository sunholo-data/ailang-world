#!/bin/bash
# V22: an NFD path in the index (core.precomposeunicode=false, as a Linux-made commit), NFC pattern é_test.go.
set -u; S=${S:?}; export HOME=$S/home GIT_CONFIG_NOSYSTEM=1 GIT_CONFIG_GLOBAL=/dev/null GIT_OPTIONAL_LOCKS=0
r=$S/r-nfd2; rm -rf "$r"; mkdir -p "$r"; cd "$r"; git init -q; git config core.precomposeunicode false
git config user.email fixture; git config user.name t
h=$(printf 'x\n' | git hash-object -w --stdin); git update-index --add --cacheinfo 100644 "$h" "$(printf 'pkg/e\xcc\x81_test.go')"
git commit -qm b --no-verify
echo "index bytes: $(git ls-files -z | od -An -c | tr -s ' ' | head -2 | tr '\n' ' ')"
for pc in false true; do for ps in ':(top,glob,icase)**/é_test.go' ':(top,glob,icase)**/*_te*t.go'; do
  echo "precompose=$pc $ps ls-files=$(git -c core.precomposeunicode=$pc ls-files -z -- "$ps" | tr '\0' '\n' | grep -c .)"
done; done
cd "$S"; rm -rf "$r"
