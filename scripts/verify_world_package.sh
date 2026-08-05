#!/usr/bin/env bash
# Nine-step, non-vacuous release gate for the world/core package projection.
set -uo pipefail

cd "$(dirname "$0")/.."

# This gate needs the PINNED v0.30.0 compiler, which is NOT the same binary the other legs use.
# CI job 1 installs `releases/latest` on PATH — measured 2026-08-05 at run 30993399332 to be
# **v0.33.0**, three minor versions past the pin (job 2's step log printed v0.30.0 in the same
# run, which is the control proving the difference is real and not an instrument artifact).
# So WORLD_PKG_AILANG_BIN lets this gate be pointed at the pin WITHOUT changing what legs 1-2
# verify against — that broader change is queue item 9's, and it is human-gated.
# There is deliberately NO silent fallback to a working binary: if this resolves to anything
# other than pinned v0.30.0, the version/SHA assertions below fail LOUDLY rather than skipping.
AILANG_BIN="${WORLD_PKG_AILANG_BIN:-${AILANG_BIN:-ailang}}"
readonly PACKAGE_DIR="packages/world-core"
readonly MANIFEST="$PACKAGE_DIR/ailang.toml"
readonly SMOKE="$PACKAGE_DIR/_smoke.ail"
# The compiler is pinned by EXACT BYTES, and those bytes are platform-specific: the rig runs
# darwin/arm64 (Mach-O) and CI runs linux/amd64 (ELF). A single constant here is a gate that can
# only ever pass on one of the two, so it is a per-platform table with a LOUD unknown-platform
# failure. Both values measured first-party 2026-08-05 from
# `releases/download/v0.30.0`, the darwin one from the rig's pinned install and the linux one
# by downloading the published tarball (its own published .sha256 verified `OK` as the control).
compiler_sha_for_platform() {
  case "$(uname -s)/$(uname -m)" in
    Darwin/arm64) printf '%s' 'e9746fef8570bc42b8cc52c0e88b7088468a5d2bd38bb8c42e27e5859b8f3fb5' ;;
    Linux/x86_64) printf '%s' '1e594d158dffa68834b21a192519b1ee98f86052b594f8c1c36c8cdc11d6cc50' ;;
    *) return 1 ;;
  esac
}
readonly GOLDEN="scripts/world_package_ready_packet.golden.json"
readonly MODULES=(types.ail contracts.ail transitions.ail logepoch.ail)
readonly EXPORTS=(world/types world/contracts world/transitions world/logepoch)

command -v python3 >/dev/null 2>&1 || { printf '%s\n' '✗ python3 is required' >&2; exit 1; }
command -v go >/dev/null 2>&1 || { printf '%s\n' '✗ go is required for host/pkgproj' >&2; exit 1; }

run_bounded() { # $1=timeout_s $2=out_file $3..=command
  local t="$1" out="$2"; shift 2
  python3 - "$t" "$out" "$@" <<'PY'
import os, signal, subprocess, sys
t = int(sys.argv[1]); out = sys.argv[2]; cmd = sys.argv[3:]
with open(out, "wb") as f:
    p = subprocess.Popen(cmd, stdout=f, stderr=subprocess.STDOUT, start_new_session=True)
    try:
        sys.exit(p.wait(timeout=t))
    except subprocess.TimeoutExpired:
        os.killpg(os.getpgid(p.pid), signal.SIGKILL)
        sys.stderr.write("✗ TIMEOUT after %ds: %s\n" % (t, " ".join(cmd)))
        sys.exit(124)
PY
}

tmp_all_ail="$(mktemp -t world_all_ail.XXXXXX)"
tmp_expected_ail="$(mktemp -t world_expected_ail.XXXXXX)"
tmp_check="$(mktemp -t world_check.XXXXXX)"
tmp_test="$(mktemp -t world_test.XXXXXX)"
tmp_smoke="$(mktemp -t world_smoke.XXXXXX)"
tmp_dry="$(mktemp -t world_dry.XXXXXX)"
tmp_proj="$(mktemp -t world_proj.XXXXXX)"
tmp_entries="$(mktemp -t world_entries.XXXXXX)"
tmp_expected_entries="$(mktemp -t world_expected_entries.XXXXXX)"
tmp_ready="$(mktemp -t world_ready.XXXXXX)"
tmp_helper="$(mktemp -t world_pkgproj.XXXXXX).go"
cleanup() {
  rm -f "$tmp_all_ail" "$tmp_expected_ail" "$tmp_check" "$tmp_test" "$tmp_smoke" \
    "$tmp_dry" "$tmp_proj" "$tmp_entries" "$tmp_expected_entries" "$tmp_ready" "$tmp_helper"
}
trap cleanup EXIT

printf '%s\n' '── World package step 1/9: required paths'
[ -d "$PACKAGE_DIR" ] || { printf '✗ missing package directory: %s\n' "$PACKAGE_DIR" >&2; exit 1; }
for path in "$MANIFEST" "$SMOKE"; do
  [ -f "$path" ] || { printf '✗ missing required file: %s\n' "$path" >&2; exit 1; }
done
required_modules=0
for module in "${MODULES[@]}"; do
  path="$PACKAGE_DIR/world/$module"
  [ -f "$path" ] || { printf '✗ missing allowlisted module: %s\n' "$path" >&2; exit 1; }
  required_modules=$((required_modules + 1))
done
[ "$required_modules" -gt 0 ] || { printf '%s\n' '✗ required module enumeration was empty' >&2; exit 1; }
printf '   ✓ package, manifest, smoke, and %d allowlisted modules exist\n' "$required_modules"

printf '%s\n' '── World package step 2/9: exact .ail allowlist'
find "$PACKAGE_DIR" -type f -name '*.ail' -print | sort > "$tmp_all_ail"
printf '%s\n' "$SMOKE" > "$tmp_expected_ail"
for module in "${MODULES[@]}"; do printf '%s\n' "$PACKAGE_DIR/world/$module"; done | sort >> "$tmp_expected_ail"
sort -o "$tmp_expected_ail" "$tmp_expected_ail"
actual_ail_count="$(wc -l < "$tmp_all_ail" | tr -d '[:space:]')"
expected_ail_count="$(wc -l < "$tmp_expected_ail" | tr -d '[:space:]')"
[ "$actual_ail_count" -gt 0 ] || { printf '%s\n' '✗ package .ail enumeration was empty' >&2; exit 1; }
[ "$expected_ail_count" -gt 0 ] || { printf '%s\n' '✗ expected .ail enumeration was empty' >&2; exit 1; }
cmp -s "$tmp_expected_ail" "$tmp_all_ail" || { printf '%s\n' '✗ package .ail set differs from exact allowlist' >&2; diff -u "$tmp_expected_ail" "$tmp_all_ail" >&2; exit 1; }
printf '   ✓ find observed exactly %s expected .ail files; no unexpected .ail file exists\n' "$actual_ail_count"

printf '%s\n' '── World package step 3/9: projection SHA-256 equality'
compared=0
for module in "${MODULES[@]}"; do
  python3 - "world/$module" "$PACKAGE_DIR/world/$module" <<'PY' || exit 1
import hashlib, sys
def digest(path):
    with open(path, "rb") as f: return hashlib.sha256(f.read()).hexdigest()
a, b = digest(sys.argv[1]), digest(sys.argv[2])
if a != b:
    sys.stderr.write("✗ projection mismatch: %s=%s %s=%s\n" % (sys.argv[1], a, sys.argv[2], b)); sys.exit(1)
PY
  compared=$((compared + 1))
done
[ "$compared" -gt 0 ] || { printf '%s\n' '✗ projection comparison was empty' >&2; exit 1; }
printf '   ✓ %d/%d projection hashes equal their canonical sources\n' "$compared" "$required_modules"

printf '%s\n' '── World package step 4/9: frozen manifest'
python3 - "$MANIFEST" <<'PY' || exit 1
import sys, tomllib
with open(sys.argv[1], "rb") as f: got = tomllib.load(f)
want = {
  "package": {"name":"world/core", "version":"0.1.0", "edition":"1", "ailang":">=0.30.0", "module_prefix":"world", "description":"AILANG World's pure semantic core"},
  "exports": {"modules":["world/types", "world/contracts", "world/transitions", "world/logepoch"]},
  "effects": {"max":[]},
}
if got != want:
    sys.stderr.write("✗ manifest is not the exact frozen structure\ngot=%r\nwant=%r\n" % (got, want)); sys.exit(1)
if len(got["exports"]["modules"]) == 0:
    sys.stderr.write("✗ manifest export enumeration was empty\n"); sys.exit(1)
print("   ✓ exact package fields, 4 exports, and empty effects")
PY

printf '%s\n' '── World package step 5/9: package check and tests'
( cd "$PACKAGE_DIR" && run_bounded 120 "$tmp_check" "$AILANG_BIN" check --package ) || { cat "$tmp_check" >&2; exit 1; }
( cd "$PACKAGE_DIR" && run_bounded 180 "$tmp_test" "$AILANG_BIN" test world/ ) || { cat "$tmp_test" >&2; exit 1; }
test_activity="$(grep -Ec 'PASS|pass|test' "$tmp_test" || true)"
[ "$test_activity" -gt 0 ] || { printf '%s\n' '✗ package test output showed zero test activity' >&2; cat "$tmp_test" >&2; exit 1; }
printf '   ✓ package check passed and tests reported %s activity line(s)\n' "$test_activity"

printf '%s\n' '── World package step 6/9: bounded smoke execution'
( cd "$PACKAGE_DIR" && run_bounded 30 "$tmp_smoke" "$AILANG_BIN" run _smoke.ail )
smoke_rc=$?
[ "$smoke_rc" -eq 0 ] || { cat "$tmp_smoke" >&2; exit "$smoke_rc"; }
grep -Fq 'rev=1 state=sha256:bbbb log=sha256:3333 proposal=true' "$tmp_smoke" || { printf '%s\n' '✗ smoke did not produce the committed flow result' >&2; cat "$tmp_smoke" >&2; exit 1; }
printf '%s\n' '   ✓ plan → verify → commit smoke completed within 30s'

printf '%s\n' '── World package step 7/9: dry-run identity and full pkgproj hashes'
( cd "$PACKAGE_DIR" && run_bounded 60 "$tmp_dry" env -u AILANG_REGISTRY_API_KEY "$AILANG_BIN" publish --dry-run )
dry_rc=$?
[ "$dry_rc" -eq 0 ] || { cat "$tmp_dry" >&2; exit "$dry_rc"; }
python3 - "$tmp_dry" <<'PY' || exit 1
import re, sys
s = open(sys.argv[1], encoding="utf-8").read()
checks = [
 (r"^Publishing world/core@0\.1\.0\.\.\.$", "package identity"),
 (r"^  Exports: \[world/types world/contracts world/transitions world/logepoch\]$", "exact export set"),
 (r"^  Effects: \[\]$", "empty effects"),
]
for pattern, name in checks:
    if not re.search(pattern, s, re.M):
        sys.stderr.write("✗ dry-run missing %s\n" % name); sys.exit(1)
for label in ("Tarball", "Content hash", "Interface hash"):
    if not re.search(r"^  %s:" % re.escape(label), s, re.M):
        sys.stderr.write("✗ dry-run missing %s observable\n" % label); sys.exit(1)
PY
cat > "$tmp_helper" <<'GO'
package main
import (
  "archive/tar"
  "bytes"
  "compress/gzip"
  "fmt"
  "io"
  "os"
  "github.com/sunholo-data/ailang-world/host/pkgproj"
)
func main() {
  m := pkgproj.Manifest{Package: pkgproj.Package{Name:"world/core", Edition:"1", AILANG:">=0.30.0"}, Exports:pkgproj.Exports{Modules:[]string{"world/types","world/contracts","world/transitions","world/logepoch"}}, Effects:pkgproj.Effects{Max:[]string{}}}
  r, err := pkgproj.CrossCheck("packages/world-core", m, os.Args[1]); if err != nil { panic(err) }
  data, err := pkgproj.CreateTarball("packages/world-core"); if err != nil { panic(err) }
  fmt.Printf("contentHash=%s\ninterfaceHash=%s\ntarballSHA256=%s\ntarballBytes=%d\n", r.Local.Content, r.Local.Interface, r.Local.Tarball, r.Local.TarballBytes)
  zr, err := gzip.NewReader(bytes.NewReader(data)); if err != nil { panic(err) }; tr := tar.NewReader(zr)
  for { h, err := tr.Next(); if err == io.EOF { break }; if err != nil { panic(err) }; fmt.Printf("entry=%s\n", h.Name) }
}
GO
run_bounded 120 "$tmp_proj" env -u AILANG_REGISTRY_API_KEY go run "$tmp_helper" "$AILANG_BIN" || { cat "$tmp_proj" >&2; exit 1; }
for key in contentHash interfaceHash tarballSHA256 tarballBytes; do
  count="$(grep -Ec "^${key}=" "$tmp_proj" || true)"
  [ "$count" -eq 1 ] || { printf '✗ pkgproj emitted %s %s times\n' "$key" "$count" >&2; cat "$tmp_proj" >&2; exit 1; }
done
printf '%s\n' '   ✓ dry-run identity/exports/effects and 3 full hashes agree; tarball byte length agrees'

printf '%s\n' '── World package step 8/9: exact tar entry allowlist'
if [ -d "$PACKAGE_DIR/.ailang" ]; then
  cache_ail_count="$(find "$PACKAGE_DIR/.ailang" -type f -name '*.ail' -print | wc -l | tr -d '[:space:]')"
  cache_control_count="$(find "$PACKAGE_DIR/.ailang" -type f ! -name '*.ail' -print | wc -l | tr -d '[:space:]')"
  [ "$cache_control_count" -gt 0 ] || { printf '%s\n' '✗ .ailang cache control found no non-.ail files; zero assertion would be vacuous' >&2; exit 1; }
  [ "$cache_ail_count" -eq 0 ] || { printf '✗ .ailang cache contains %s .ail file(s)\n' "$cache_ail_count" >&2; exit 1; }
  printf '   cache control fired on %s non-.ail file(s); observed zero cached .ail files\n' "$cache_control_count"
else
  [ ! -e "$PACKAGE_DIR/.ailang" ] || { printf '%s\n' '✗ .ailang is not a directory' >&2; exit 1; }
  printf '%s\n' '   .ailang directory absence asserted explicitly'
fi
grep '^entry=' "$tmp_proj" | sed 's/^entry=//' > "$tmp_entries"
printf '%s\n' ailang.toml _smoke.ail world/types.ail world/contracts.ail world/transitions.ail world/logepoch.ail | sort > "$tmp_expected_entries"
sort -o "$tmp_entries" "$tmp_entries"
entry_count="$(wc -l < "$tmp_entries" | tr -d '[:space:]')"
expected_entry_count="$(wc -l < "$tmp_expected_entries" | tr -d '[:space:]')"
[ "$entry_count" -gt 0 ] || { printf '%s\n' '✗ tar entry enumeration was empty' >&2; exit 1; }
[ "$expected_entry_count" -gt 0 ] || { printf '%s\n' '✗ expected tar enumeration was empty' >&2; exit 1; }
cmp -s "$tmp_expected_entries" "$tmp_entries" || { printf '%s\n' '✗ tar entries differ from exact allowlist' >&2; diff -u "$tmp_expected_entries" "$tmp_entries" >&2; exit 1; }
printf '   ✓ tar contains exactly %s allowlisted entries\n' "$entry_count"

printf '%s\n' '── World package step 9/9: canonical ready-packet golden'
compiler_version="$($AILANG_BIN --version | sed -n '1p')"
[ "$compiler_version" = 'AILANG v0.30.0' ] || { printf '✗ wrong compiler version: %s\n' "$compiler_version" >&2; exit 1; }
measured_compiler_sha="$(python3 - "$AILANG_BIN" <<'PY'
import hashlib, sys
with open(sys.argv[1], "rb") as f: print(hashlib.sha256(f.read()).hexdigest())
PY
)"
expected_compiler_sha="$(compiler_sha_for_platform)" || {
  printf '✗ no pinned v0.30.0 SHA-256 recorded for platform %s/%s — refusing to pass an unpinned compiler\n' "$(uname -s)" "$(uname -m)" >&2; exit 1; }
[ -n "$expected_compiler_sha" ] || { printf '%s\n' '✗ platform SHA table returned empty' >&2; exit 1; }
[ "$measured_compiler_sha" = "$expected_compiler_sha" ] || { printf '✗ compiler SHA-256 mismatch on %s/%s: measured=%s expected=%s (binary=%s)\n' "$(uname -s)" "$(uname -m)" "$measured_compiler_sha" "$expected_compiler_sha" "$AILANG_BIN" >&2; exit 1; }
printf '   ✓ compiler pinned by exact bytes: %s on %s/%s\n' "$compiler_version" "$(uname -s)" "$(uname -m)"
# The ready packet is the ARTIFACT's identity, and every field in it was measured this run to be
# toolchain-independent (content/interface are sha256 over bytes; the tarball reproduced
# byte-identically across go1.25.6 and the go1.26.5-built CLI). compilerSHA256 is the one field
# that is NOT — it is provenance about the machine, not about the package — so it is asserted
# above against the platform table and deliberately kept OUT of the byte-compared golden. Putting
# it in would make the golden pass on exactly one platform, which is a gate that cannot run in CI.
python3 - "$tmp_proj" "$tmp_ready" "$compiler_version" "$measured_compiler_sha" <<'PY' || exit 1
import json, sys
values = {}
for line in open(sys.argv[1], encoding="utf-8"):
    if "=" in line:
        k, v = line.rstrip("\n").split("=", 1); values[k] = v
required = ("contentHash", "interfaceHash", "tarballSHA256", "tarballBytes")
if any(not values.get(k) for k in required):
    sys.stderr.write("✗ ready packet input is incomplete\n"); sys.exit(1)
packet = {"compilerVersion":sys.argv[3], "contentHash":values["contentHash"], "effects":[], "exports":["world/types","world/contracts","world/transitions","world/logepoch"], "interfaceHash":values["interfaceHash"], "package":"world/core", "tarballBytes":int(values["tarballBytes"]), "tarballSHA256":values["tarballSHA256"], "version":"0.1.0"}
with open(sys.argv[2], "wb") as f: f.write((json.dumps(packet, sort_keys=True, separators=(",", ":")) + "\n").encode())
PY
[ -s "$tmp_ready" ] || { printf '%s\n' '✗ ready packet enumeration was zero-length' >&2; exit 1; }
[ -f "$GOLDEN" ] || { printf '✗ missing committed ready-packet golden: %s\n' "$GOLDEN" >&2; cat "$tmp_ready" >&2; exit 1; }
cmp -s "$GOLDEN" "$tmp_ready" || { printf '%s\n' '✗ ready packet differs byte-for-byte from golden' >&2; diff -u "$GOLDEN" "$tmp_ready" >&2; exit 1; }
printf '%s\n' '   ✓ canonical JSON equals committed golden byte-for-byte'
printf '%s\n' '✓ world package gate PASSED: 9/9 steps performed non-zero work'
