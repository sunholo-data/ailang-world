#!/usr/bin/env bash
# test_tracked_binary_gate.sh — adversarial arms for scripts/tracked_binary_gate.sh.
#
# Every arm is SELF-CONTAINED: it builds a throwaway git repo, commits synthetic blobs, and runs
# the gate there with the same allowance verify_go.sh declares. No arm references a SHA or file of
# this repository. Payloads are synthesized (a Mach-O magic plus NUL-bearing bytes) so the arms are
# byte-identical on darwin and the ubuntu runner. Run by verify_go.sh before the gate itself.
#
# Mutation-tested 2026-10-14: with the detector narrowed to git's 8000-byte window, arm 12 reds;
# with the old `git diff --numstat` classifier, arms 9, 10 and 12 pass a binary — so these arms
# are what keep the two measured holes closed.
set -uo pipefail
cd "$(dirname "$0")/.." || exit 1
GATE="${TBG_GATE_UNDER_TEST:-$(pwd)/scripts/tracked_binary_gate.sh}"
ALLOW=(--allow-png 'website/static/img/v1/' --png-file-cap 1048576 --png-dir-cap 8388608)
pass=0; fail=0
ok()   { echo "ok $*";      pass=$((pass+1)); }
notok(){ echo "not ok $*";  fail=$((fail+1)); }
WORK="$(mktemp -d)"; trap 'rm -rf "$WORK"' EXIT

# mk <out> <kind> [arg] — write a synthetic blob. Kinds:
#   png [idat_bytes]  well-formed PNG (IHDR, one IDAT, IEND)
#   png-text          well-formed PNG plus a tEXt chunk
#   png-badcrc        well-formed PNG with one IDAT byte flipped after the CRC was computed
#   png-trailing      well-formed PNG with a binary payload appended after IEND
#   bin               Mach-O magic plus NUL-bearing bytes
#   padded-bin        8,100 ASCII bytes, then bin (NUL first appears past git's 8000-byte window)
mk() {
  mkdir -p "$(dirname "$1")"
  python3 -I - "$@" <<'PY'
import struct, sys, zlib
out, kind = sys.argv[1], sys.argv[2]
arg = int(sys.argv[3]) if len(sys.argv) > 3 else 64
def chunk(t, b): return struct.pack(">I", len(b)) + t + b + struct.pack(">I", zlib.crc32(t + b))
def png(extra=b"", idat=64):
    body = bytes((i * 7) % 251 for i in range(idat))
    return (b"\x89PNG\r\n\x1a\n" + chunk(b"IHDR", struct.pack(">IIBBBBB", 1, 1, 8, 2, 0, 0, 0))
            + extra + chunk(b"IDAT", body) + chunk(b"IEND", b""))
binary = b"\xcf\xfa\xed\xfe\x0c\x00\x00\x01" + bytes(range(256)) * 64
if kind == "png":          data = png(idat=arg)
elif kind == "png-text":   data = png(chunk(b"tEXt", b"Comment\x00payload"))
elif kind == "png-badcrc":
    data = bytearray(png()); i = data.index(b"IDAT") + 10; data[i] ^= 0xFF; data = bytes(data)
elif kind == "png-trailing": data = png() + binary
elif kind == "bin":        data = binary
elif kind == "padded-bin": data = b"A" * 8100 + binary
else: sys.exit(f"mk: unknown kind {kind}")
open(out, "wb").write(data)
PY
}

# arm <name> <expect: pass|reason> <setup-fn> [gate args...]
arm() {
  local name="$1" expect="$2" setup="$3"; shift 3
  local repo="$WORK/$name"
  mkdir -p "$repo" && git -C "$repo" init -q && echo seed > "$repo/README"
  ( cd "$repo" && "$setup" ) || { notok "$name — setup failed"; return; }
  git -C "$repo" add -A && git -C "$repo" -c user.name=t -c user.email=t@t commit -q --allow-empty -m arm
  local out rc
  out="$(cd "$repo" && /bin/bash "$GATE" "$@" 2>&1)"; rc=$?
  if [ "$expect" = pass ]; then
    if [ $rc -eq 0 ]; then ok "$name"; else notok "$name — expected pass, rc=$rc: $out"; fi
  elif [ $rc -ne 0 ] && printf '%s\n' "$out" | grep -qF -- "$expect"; then
    ok "$name ($expect)"
  else
    notok "$name — expected '$expect', rc=$rc: $out"
  fi
}

D=website/static/img/v1
s01() { mk $D/a.png png; }
s02() { mk website/static/img/a.png png; }
s03() { mk $D/evil.png bin; }
s04() { mk $D/a.png png-trailing; }
s05() { mk $D/a.png png-text; }
s06() { mk $D/a.png png-badcrc; }
s07() { mk $D/big.png png 1153434; }
s08() { for i in 1 2 3 4 5 6 7 8 9; do mk $D/shot-$i.png png 996147; done; }
s09() { printf '*.png diff\n' > .gitattributes; mk $D/x.png bin; }
s10() { printf '*.bin diff\n' > .gitattributes; mk tools/x.bin bin; }
s11() { rm README; }
s12() { mk tools/padded padded-bin; }
s13() { for n in workbench-light workbench-dark packet graph provenance phone; do mk $D/$n.png png 4096; done; }
s14() { mk $D/Shot.png png; mk $D/sub/a.png png; }

arm 01-valid-png-in-dir         pass                     s01 "${ALLOW[@]}"
arm 02-valid-png-outside-dir    not-allowlisted          s02 "${ALLOW[@]}"
arm 03-renamed-binary           bad-signature            s03 "${ALLOW[@]}"
arm 04-polyglot-after-iend      trailing-bytes           s04 "${ALLOW[@]}"
arm 05-text-chunk               chunk-not-allowed:tEXt   s05 "${ALLOW[@]}"
arm 06-corrupt-idat             bad-crc                  s06 "${ALLOW[@]}"
arm 07-over-file-cap            file-cap-exceeded        s07 "${ALLOW[@]}"
arm 08-over-dir-cap             dir-cap-exceeded         s08 "${ALLOW[@]}"
arm 09-gitattributes-in-dir     bad-signature            s09 "${ALLOW[@]}"
arm 10-gitattributes-tree-wide  not-allowlisted          s10 "${ALLOW[@]}"
arm 11-empty-tree               "enumerated 0 blobs"     s11 "${ALLOW[@]}"
arm 12-nul-past-8000-bytes      not-allowlisted          s12 "${ALLOW[@]}"
arm 13-row150-deliverable-names pass                     s13 "${ALLOW[@]}"
arm 14-bad-name-and-subdir      not-allowlisted          s14 "${ALLOW[@]}"
arm 15-no-allowance-no-images   not-allowlisted          s01

echo "test_tracked_binary_gate: $pass passed, $fail failed"
[ "$fail" -eq 0 ] && [ "$pass" -eq 15 ]
