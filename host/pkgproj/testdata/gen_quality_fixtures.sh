#!/usr/bin/env bash
# Regenerates the `ailang pkg quality --json` parser fixtures for host/pkgproj
# (w-interface-hash-covers-the-interface §7). Tool output is written AS PRODUCED;
# derived negatives exist only where no tool output reaches a guard, and each is
# ONE named jq transform of the generated pristine fixture — never hand-typed.
#
#   AILANG_BIN=$HOME/.pinned-ailang/ailang host/pkgproj/testdata/gen_quality_fixtures.sh
#
# Every capture runs in a mktemp -d copy of packages/world-core, so paths that
# the tool embeds in error text are temp paths, never a home directory. The one
# nondeterministic field is `.contracts.wall_seconds` (timing); parsers ignore it.
set -euo pipefail
here="$(cd "$(dirname "$0")" && pwd)"
root="$(cd "$here/../../.." && pwd)"
: "${AILANG_BIN:?AILANG_BIN must name the pinned binary}"
ver="$("$AILANG_BIN" --version)"
case "$ver" in *'AILANG v0.41.0'*) ;; *) echo "✗ AILANG_BIN is not v0.41.0" >&2; exit 1 ;; esac

tmp="$(mktemp -d)"
trap 'rm -rf "$tmp"' EXIT
echo "gen_quality_fixtures: temp root $tmp"

fresh() { rm -rf "$tmp/pkg"; cp -R "$root/packages/world-core" "$tmp/pkg"; }

# capture NAME WANT_RC [--no-run]
capture() {
  local name="$1" want="$2"; shift 2
  local rc=0
  (cd "$tmp/pkg" && "$AILANG_BIN" pkg quality --json "$@" .) >"$here/$name" || rc=$?
  [ "$rc" -eq "$want" ] || { echo "✗ $name: rc=$rc, want $want" >&2; exit 1; }
  echo "  $name rc=$rc $(wc -c <"$here/$name" | tr -d ' ') bytes"
}

# 1. pristine, the production invocation (--no-run): rc 2, gates [PUB015 "_smoke.ail failed"]
fresh; capture quality_world_core_pristine.json 2 --no-run
# 2. pristine, RUN mode: rc 0, gates [] — the exit-0 positive control
fresh; capture quality_world_core_pristine_run.json 0
# 3. the row's stimulus: an exported ADT gains a constructor
fresh; perl -0pi -e 's/  \| ProofReceipt\(HashRef\)\n/  | ProofReceipt(HashRef)\n  | ProbeArm(HashRef)\n/ or die "anchor"' "$tmp/pkg/world/types.ail"
capture quality_world_core_probe_ctor.json 2 --no-run
# 4. a type-error export: compile fails, no hash_v2, gates [PUB000, PUB015]
fresh; printf '%s\n' 'export func probeBroken(x: int) -> int ! {} { x ++ "s" }' >>"$tmp/pkg/world/types.ail"
capture quality_world_core_v2_unbuilt.json 2 --no-run

p="$here/quality_world_core_pristine.json"
derive() { jq -c "$2" "$p" >"$here/$1"; echo "  $1 := $2"; }
derive derived_quality_v2_holds_v1.json '.interface.hash_v2 = .interface.hash_v1'
derive derived_quality_v1_absent.json   'del(.interface.hash_v1)'
derive derived_quality_schema_v2.json   '.schema = "ailang.package-quality/v2"'
derive derived_quality_extra_gate.json  '.gates += [{"code":"PUB000","level":"gate","msg":"derived: a non-spurious gate"}]'
derive derived_quality_pub015_other_msg.json '.gates[0].msg = "_smoke.ail failed: derived other message"'

if grep -l "$HOME" "$here"/*.json; then echo "✗ a fixture embeds \$HOME" >&2; exit 1; fi
echo "gen_quality_fixtures: OK"
