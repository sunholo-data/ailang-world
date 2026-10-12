#!/usr/bin/env bash
# Fixture generator for queue row 26 (w-bounded-z3-report-producer) design evidence.
#
# Generates REAL `ai-check` JSON with the pinned released binary
# (AILANG v0.41.0, ~/.pinned-ailang/ailang) against scratch .ail modules
# written to a TEMPORARY directory outside the repo, so no new .ail file
# enters a tree swept by scripts/verify_ail.sh's LEG1_MODULES membership
# compare (ROOTS sweeps design_docs/, world/, packages/se-tools/).
#
# Five shapes (see the row-26 design doc §7 test plan):
#   1. verifiedpass  — rc=0, check.passed=true, verify.results[0].status="verified".
#   2. refuted       — rc=1, one status="counterexample" record with a Z3 model.
#   3. encerror       — rc=1 on v0.41.0 (measured), verify.errors=1,
#                       one status="error" record with a reason string.
#                      (The inherited claim "Z3 encoding error exits 0" was
#                       measured on v0.30.0, V10 of the parent doc; on the
#                       pinned v0.41.0 this shape exits 1 — recorded, not
#                       assumed.)
#   4. mixed          — rc=1: one "verified" + one "error" in the same run.
#   5. skippedshape   — rc=0, check.passed=true, verify.skipped=1,
#                       verify.verified=0: THE silent class on v0.41.0 —
#                       exit code 0 while the contract was never verified.
#   6. loaderror      — rc=1, check.passed=false, verify.available=false
#                       (a parse error in the module).
#   7. nocontracts    — rc=0, check.passed=true, every counter 0,
#                       results=[]: the CLEAN-BUT-VACUOUS shape. The
#                       anti-vacuity guard exists for exactly this fixture.
#
# Output: design_docs/verification/world-iter251/fixtures/<name>.json (+ .rc, .stderr)
set -euo pipefail
cd "$(dirname "$0")"
BIN="${AILANG_BIN:-$HOME/.pinned-ailang/ailang}"
echo "── generator using: $BIN"
"$BIN" --version 2>/dev/null | head -1

scratch="$(mktemp -d "${TMPDIR:-/tmp}/row26-aicheck.XXXXXX")"
trap 'rm -rf "$scratch"' EXIT
mkdir -p "$scratch/fixtures"

cat > "$scratch/fixtures/verifiedpass.ail" <<'AIL'
module fixtures/verifiedpass

-- Scratch fixture (row 26 design evidence): one contract Z3 can prove.
export func twicePlusOne(x: int) -> int ! {}
ensures { result == 2 * x + 1 }
tests [
  ((0), 1),
  ((5), 11)
]
{
  2 * x + 1
}
AIL

cat > "$scratch/fixtures/refuted.ail" <<'AIL'
module fixtures/refuted

-- Scratch fixture (row 26 design evidence): the body contradicts the ensures.
export func notTwice(x: int) -> int ! {}
ensures { result == 2 * x }
tests [
  ((3), 6)
]
{
  2 * x + 1
}
AIL

cat > "$scratch/fixtures/encerror.ail" <<'AIL'
module fixtures/encerror

-- Scratch fixture (row 26 design evidence): interpolation in ensures is
-- unencodable (world/logepoch.ail records the same class on v0.41.0).
export func showInt(x: int) -> string ! {}
ensures { result == "v${x}" }
tests [
  ((1), "v1")
]
{
  "v${x}"
}
AIL

cat > "$scratch/fixtures/mixed.ail" <<'AIL'
module fixtures/mixed

-- Scratch fixture (row 26 design evidence): one verifiable contract + one
-- encoding-error contract in the same module.
export func good(x: int) -> int ! {}
ensures { result == x + 1 }
tests [ ((1), 2) ]
{ x + 1 }

export func showInt(x: int) -> string ! {}
ensures { result == "v${x}" }
tests [ ((1), "v1") ]
{
  "v${x}"
}
AIL

cat > "$scratch/fixtures/skippedshape.ail" <<'AIL'
module fixtures/skippedshape

-- Scratch fixture (row 26 design evidence): ensures calls a sibling user
-- function that is not SMT-encodable in this context — the cacheKey shape
-- recorded in world/logepoch.ail. ai-check SKIPS the contract and exits 0.
export func renderRef(algo: string, digest: string) -> string { "${algo}:${digest}" }

export func usesSibling(algo: string, digest: string) -> string ! {}
ensures { result == renderRef(algo, digest) }
tests [ (("sha256", "aa"), "sha256:aa") ]
{
  renderRef(algo, digest)
}
AIL

cat > "$scratch/fixtures/loaderror.ail" <<'AIL'
module fixtures/loaderror

-- Scratch fixture (row 26 design evidence): a parse error (export helper).
export helper func bad(n: int) -> int { n + 1 }
AIL

cat > "$scratch/fixtures/nocontracts.ail" <<'AIL'
module fixtures/nocontracts

-- Scratch fixture (row 26 design evidence): a module with NO contracted
-- function: nothing to verify, so ai-check emits results=[] with rc=0.
export type Box = { v: int }

export func plain(x: int) -> int
{
  x + 1
}
AIL

mkdir -p fixtures
for name in verifiedpass refuted encerror mixed skippedshape loaderror nocontracts; do
  echo "── ai-check fixtures/$name.ail (cwd: scratch)"
  rc=0
  ( cd "$scratch" && "$BIN" ai-check -timeout 5s "fixtures/$name.ail" ) \
    > "fixtures/$name.json" 2> "fixtures/$name.stderr" || rc=$?
  # a refuted contract exits 1 by design; the rc is data, not failure
  printf 'rc=%s\n' "$rc" > "fixtures/$name.rc"
  echo "   rc=$rc stdout_bytes=$(wc -c < fixtures/$name.json | tr -d ' ') stderr_bytes=$(wc -c < fixtures/$name.stderr | tr -d ' ')"
done
echo "── done"