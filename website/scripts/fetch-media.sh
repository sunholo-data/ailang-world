#!/bin/sh
# Fetch the homepage film and its poster into static/media/ before a build, refusing any byte
# that does not match its sha256. The repo's tracked-binary gate keeps video out of git, so the
# files live as assets on a GitHub release and are copied into the built site at build time;
# the deployed site then serves them from its own origin (GitHub Pages, with range requests).
#
# Source: films/world_intro.py in the ailang-explainers repo (narration: Gemini Sulafat en-GB,
# voiced through AILANG std/ai.callSpeech by films/narrate/explainer-gemini).
# A new cut: upload it to a release, then change the URL and sha256 below in the same PR.
set -eu
cd "$(dirname "$0")/../static"
mkdir -p media
cd media
BASE=https://github.com/sunholo-data/ailang-world/releases/download/v1.0.0-rc.1
fetch() { # file sha256
  if [ -f "$1" ] && [ "$(shasum -a 256 "$1" | cut -d' ' -f1)" = "$2" ]; then
    echo "fetch-media: $1 ok (cached)"; return
  fi
  curl -sfL -o "$1.part" "$BASE/$1"
  got=$(shasum -a 256 "$1.part" | cut -d' ' -f1)
  [ "$got" = "$2" ] || { rm -f "$1.part"; echo "fetch-media: $1 sha256 $got, want $2" >&2; exit 1; }
  mv "$1.part" "$1"; echo "fetch-media: $1 ok"
}
fetch ailang-world-intro.mp4 25c986ae4d3e280e1b5e10503863989aed720e4f3aa41f1c658fd73ed90e1dfb
fetch ailang-world-intro-poster.jpg 149a226f5b5d98ff73bab29622e1a64eee7a9a7bec16491d0ef6d486313e9684
