#!/usr/bin/env bash
# tracked_binary_gate.sh — no binary blob is committed unless an explicit, content-verified
# allowance admits it. Called by verify_go.sh; exercised by test_tracked_binary_gate.sh.
#
#   tracked_binary_gate.sh [--allow-png <dir/> --png-file-cap <bytes> --png-dir-cap <bytes>]
#                          [--allow-font <dir/> --font-file-cap <bytes> --font-dir-cap <bytes>
#                           --font-pin <name>=<sha256> ...]
#
# With no flags NOTHING binary is allowed. The allowance is declared by the CALLER so the
# policy stays readable in verify_go.sh. Design: design_docs/implemented/w-release-png-allowance.md
# (fleet ticket ci:world:tracked-binary-rejects-release-png; fleet-authored, D-WORLD-DRIVER-1).
#
# HISTORY. SM.A (squash 13315da) committed a 15.7 MB compiled darwin/arm64 `ailang-worldd`
# Mach-O into the repository root. It passed the codex executor, a sonnet evaluator scoring
# 87/100, the controller's four-gate re-run and both CI jobs — because nothing looked. The
# first gate (#42) used git's own binary classification (`git diff --numstat` → `-\t-`).
# That classifier is STEERABLE, measured 2026-10-14:
#   - a committed `.gitattributes` line `*.png diff` makes a Mach-O count as 334 text lines;
#   - git only looks for a NUL in the first 8000 bytes, so 8,100 ASCII bytes in front of a
#     Mach-O make it count as 335 text lines, attributes pinned or not.
# Both passed the old gate. So the detector is now content-only and whole-blob: a blob is
# BINARY iff it contains a NUL byte anywhere. Every Mach-O, ELF and PNG does; the tree at
# fd321b6 had 0 of 3471. It reads COMMITTED blobs (HEAD), never the working tree.
#
# NON-VACUITY (rule 3a): an enumeration that returns nothing is indistinguishable from a clean
# repo. The blob count must be > 0, and the bytes streamed must equal the sizes git reported,
# before the binary count is believed.
#
# ALLOWANCE. A binary blob is admitted only if its path is <dir/><name>.png with name
# [a-z0-9][a-z0-9-]* (no subdirectories) AND the blob is a well-formed PNG: signature, IHDR
# first (length 13), every chunk CRC valid, chunk types limited to IMAGE_CHUNKS (no text/EXIF
# payload chunks), exactly one IEND, last, with zero trailing bytes, and size within caps.
#
# FONT ALLOWANCE (D-WORLD-90, attended 2026-10-22). A second, narrower allowance for the World
# app's OFL fonts: a binary blob is admitted only if its path is <font-dir/><name> with no
# subdirectory, <name> is PINNED by --font-pin <name>=<sha256> (exact bytes: any new or changed
# font fails until a reviewed edit to verify_go.sh updates the pin), the blob starts with a font
# signature (TrueType 00 01 00 00 or 'true', OpenType 'OTTO', WOFF2 'wOF2'), and the per-file and
# directory caps hold. Reasons: not-allowlisted, bad-signature, sha256-mismatch,
# file-cap-exceeded, dir-cap-exceeded. The PNG allowance is unchanged.
set -euo pipefail

allow_png="" file_cap="" dir_cap=""
allow_font="" font_file_cap="" font_dir_cap="" font_pins=""
while [ $# -gt 0 ]; do
  case "$1" in
    --allow-png)    allow_png="${2:?--allow-png needs a directory}"; shift 2 ;;
    --png-file-cap) file_cap="${2:?--png-file-cap needs bytes}"; shift 2 ;;
    --png-dir-cap)  dir_cap="${2:?--png-dir-cap needs bytes}"; shift 2 ;;
    --allow-font)    allow_font="${2:?--allow-font needs a directory}"; shift 2 ;;
    --font-file-cap) font_file_cap="${2:?--font-file-cap needs bytes}"; shift 2 ;;
    --font-dir-cap)  font_dir_cap="${2:?--font-dir-cap needs bytes}"; shift 2 ;;
    --font-pin)      font_pins="${font_pins}${font_pins:+,}${2:?--font-pin needs <name>=<sha256>}"; shift 2 ;;
    *) echo "tracked_binary_gate.sh: FATAL: unknown argument: $1" >&2; exit 2 ;;
  esac
done
if [ -n "$allow_png" ]; then
  case "$allow_png" in */) ;; *) echo "tracked_binary_gate.sh: FATAL: --allow-png must end in '/': $allow_png" >&2; exit 2 ;; esac
  if [ -z "$file_cap" ] || [ -z "$dir_cap" ]; then
    echo "tracked_binary_gate.sh: FATAL: --allow-png requires --png-file-cap and --png-dir-cap (no default caps)" >&2
    exit 2
  fi
fi

if [ -n "$allow_font" ]; then
  case "$allow_font" in */) ;; *) echo "tracked_binary_gate.sh: FATAL: --allow-font must end in '/': $allow_font" >&2; exit 2 ;; esac
  if [ -z "$font_file_cap" ] || [ -z "$font_dir_cap" ] || [ -z "$font_pins" ]; then
    echo "tracked_binary_gate.sh: FATAL: --allow-font requires --font-file-cap, --font-dir-cap and at least one --font-pin" >&2
    exit 2
  fi
fi

python3 -I - "$allow_png" "${file_cap:-0}" "${dir_cap:-0}" "$allow_font" "${font_file_cap:-0}" "${font_dir_cap:-0}" "$font_pins" <<'PY'
import hashlib, re, struct, subprocess, sys, zlib

allow_dir, file_cap, dir_cap = sys.argv[1], int(sys.argv[2]), int(sys.argv[3])
font_dir, font_file_cap, font_dir_cap = sys.argv[4], int(sys.argv[5]), int(sys.argv[6])
FONT_PINS = {}
for pin in filter(None, sys.argv[7].split(",")):
    name, _, sha = pin.partition("=")
    if not name or "/" in name or not re.fullmatch(r"[0-9a-f]{64}", sha):
        print(f"tracked_binary_gate.sh: FATAL: --font-pin must be <name>=<sha256>: {pin}", file=sys.stderr)
        sys.exit(2)
    FONT_PINS[name] = sha
FONT_SIGS = (b"\x00\x01\x00\x00", b"true", b"OTTO", b"wOF2")
PNG_SIG = b"\x89PNG\r\n\x1a\n"
IMAGE_CHUNKS = {b"IHDR", b"PLTE", b"IDAT", b"IEND", b"sRGB", b"gAMA", b"cHRM", b"iCCP", b"pHYs"}
NAME = re.compile(r"[a-z0-9][a-z0-9-]*\.png")

def fatal(msg):
    print(f"tracked_binary_gate.sh: FATAL: {msg}", file=sys.stderr)
    sys.exit(1)

def font_reason(name, b):
    """None if b is the pinned font <name> within the font allowance; else one stable reason token."""
    if name not in FONT_PINS:
        return "not-allowlisted"
    if len(b) > font_file_cap:
        return "file-cap-exceeded"
    if not b[:4] in FONT_SIGS:
        return "bad-signature"
    if hashlib.sha256(b).hexdigest() != FONT_PINS[name]:
        return "sha256-mismatch"
    return None

def png_reason(b):
    """None if b is a well-formed PNG within the allowance; else one stable reason token."""
    if len(b) > file_cap:
        return "file-cap-exceeded"
    if not b.startswith(PNG_SIG):
        return "bad-signature"
    i, first, seen_iend = len(PNG_SIG), True, False
    while i < len(b):
        if seen_iend:
            return "trailing-bytes"
        if i + 8 > len(b):
            return "truncated-chunk"
        n, typ = struct.unpack(">I", b[i:i + 4])[0], b[i + 4:i + 8]
        end = i + 12 + n
        if end > len(b):
            return "truncated-chunk"
        if first and (typ != b"IHDR" or n != 13):
            return "bad-ihdr"
        first = False
        if typ not in IMAGE_CHUNKS:
            return "chunk-not-allowed:" + typ.decode("latin-1")
        if zlib.crc32(b[i + 4:i + 8 + n]) != struct.unpack(">I", b[i + 8 + n:end])[0]:
            return "bad-crc"
        seen_iend = typ == b"IEND"
        i = end
    return None if seen_iend else "missing-iend"

ls = subprocess.run(["git", "ls-tree", "-r", "-z", "HEAD"], capture_output=True)
if ls.returncode != 0:
    fatal("git ls-tree HEAD failed: " + ls.stderr.decode(errors="replace").strip())
tree = ls.stdout
blobs = []  # (sha, path)
for rec in tree.split(b"\0"):
    if not rec:
        continue
    meta, path = rec.split(b"\t", 1)
    mode, typ, sha = meta.split(b" ")
    if typ != b"blob" or mode == b"120000":  # submodule commits and symlink targets are not content
        continue
    blobs.append((sha, path.decode("utf-8", "surrogateescape")))
if not blobs:
    fatal("the tracked-binary detector enumerated 0 blobs; the instrument is broken, so every "
          "'no binaries' result it reports is void.")

proc = subprocess.run(["git", "cat-file", "--batch"], input=b"".join(s + b"\n" for s, _ in blobs),
                      capture_output=True, check=True)
out, pos, declared, streamed = proc.stdout, 0, 0, 0
binaries = []  # (path, bytes)
for sha, path in blobs:
    nl = out.index(b"\n", pos)
    hsha, htyp, hsize = out[pos:nl].split(b" ")
    if hsha != sha or htyp != b"blob":
        fatal(f"cat-file stream out of step at {path!r}")
    size = int(hsize)
    data = out[nl + 1:nl + 1 + size]
    declared += size
    streamed += len(data)
    pos = nl + 1 + size + 1
    if b"\0" in data:
        binaries.append((path, data))
if streamed != declared or pos != len(out):
    fatal(f"streamed {streamed} of {declared} declared bytes; the detector did not read every blob.")

rejected, allowed, allowed_bytes, fonts, font_bytes = [], 0, 0, 0, 0
for path, data in binaries:
    if font_dir and path.startswith(font_dir):
        rest = path[len(font_dir):]
        reason = "not-allowlisted" if "/" in rest else font_reason(rest, data)
        if reason:
            rejected.append((path, reason))
        else:
            fonts += 1
            font_bytes += len(data)
        continue
    if not allow_dir or not path.startswith(allow_dir) or not NAME.fullmatch(path[len(allow_dir):]):
        rejected.append((path, "not-allowlisted"))
        continue
    reason = png_reason(data)
    if reason:
        rejected.append((path, reason))
    else:
        allowed += 1
        allowed_bytes += len(data)
if allowed_bytes > dir_cap and allowed:
    rejected.append((allow_dir, f"dir-cap-exceeded ({allowed_bytes} > {dir_cap} bytes)"))
if font_bytes > font_dir_cap and fonts:
    rejected.append((font_dir, f"dir-cap-exceeded ({font_bytes} > {font_dir_cap} bytes)"))

if rejected:
    print("tracked_binary_gate.sh: FATAL: binary blob(s) rejected:", file=sys.stderr)
    for path, reason in rejected:
        print(f"    {path}  {reason}", file=sys.stderr)
    print("  Compiled/binary artifacts never belong in git: 'git rm --cached <path>' and .gitignore it.",
          file=sys.stderr)
    if allow_dir:
        print(f"  The only binary allowance is well-formed PNG under {allow_dir} (see verify_go.sh).",
              file=sys.stderr)
    if font_dir:
        print(f"  ...plus the fonts pinned by sha256 under {font_dir} (D-WORLD-90; see verify_go.sh).",
              file=sys.stderr)
    sys.exit(1)
if fonts:
    print(f"   ✓ {fonts} pinned font(s) validated, {font_bytes} bytes ({font_dir})")
if allowed:
    print(f"   ✓ {allowed} allowlisted PNG(s) validated, {allowed_bytes} bytes ({allow_dir})")
if allowed or fonts:
    print(f"   ✓ 0 other binary blobs among {len(blobs)} tracked files")
else:
    print(f"   ✓ 0 binary blobs among {len(blobs)} tracked files")
PY
