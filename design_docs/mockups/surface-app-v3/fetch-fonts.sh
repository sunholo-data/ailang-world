#!/bin/sh
# Fetch the working copy's OFL fonts from google/fonts at pinned commits, refusing any byte
# that does not match its sha256. The repo's tracked-binary gate keeps fonts out of git.
set -eu
cd "$(dirname "$0")/dist/fonts"
fetch() { # url sha256 file
  curl -sfL -o "$3.part" "https://raw.githubusercontent.com/google/fonts/$1"
  got=$(shasum -a 256 "$3.part" | cut -d' ' -f1)
  [ "$got" = "$2" ] || { rm -f "$3.part"; echo "fetch-fonts: $3 sha256 $got, want $2" >&2; exit 1; }
  mv "$3.part" "$3"; echo "ok $3"
}
fetch 76fca9fd0bb4ea46583f92e978660f3984ab9442/ofl/montserrat/Montserrat%5Bwght%5D.ttf 0f7b311b2f3279e4eef9b2f968bcdbab6e28f4daeb1f049f4f278a902bcd82f7 Montserrat-Variable.ttf
fetch 633f3200539c52ee0aba2dfd7f46921417a81877/ofl/ibmplexmono/IBMPlexMono-Regular.ttf 6a3412f058c7d8dfd9170c41e85ade48e5156ecb89356110ca57a0a27734af46 IBMPlexMono-Regular.ttf
fetch 633f3200539c52ee0aba2dfd7f46921417a81877/ofl/ibmplexmono/IBMPlexMono-Medium.ttf a9b4c49bb299e05b5f6c481e7fb5e78943d2793249a0c8874ab574a2d1ea6755 IBMPlexMono-Medium.ttf
