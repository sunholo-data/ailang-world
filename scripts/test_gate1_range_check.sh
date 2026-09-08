#!/usr/bin/env bash
# test_gate1_range_check.sh — self-contained arms for scripts/gate1_range_check.sh.
#
# Every arm is SELF-CONTAINED: it builds a throwaway git repo and a SHA-keyed stub `gh`, and
# never references a SHA of this repository (both CI jobs check out shallow, fetch-depth 1).
# Every invocation of the instrument goes through an explicit `/bin/bash "$SCRIPT_UT" ...` —
# never the login shell (zsh) and never `sh`.
set -uo pipefail
cd "$(dirname "$0")/.." || exit 1
SCRIPT_UT=scripts/gate1_range_check.sh
pass=0; fail=0
ok()   { echo "ok $*";      pass=$((pass+1)); }
notok(){ echo "not ok $*";  fail=$((fail+1)); }

VERIFIED_FIXTURE="$(pwd)/scripts/testdata/gate1_check_runs_verified.json"
ZERO_FIXTURE="$(pwd)/scripts/testdata/gate1_check_runs_zero.json"

# ── helpers ──────────────────────────────────────────────────────────────────
# make_stub_gh <out> <SHA:fixture>... — writes a bash script that maps its
# `repos/<owner>/<repo>/commits/<SHA>/check-runs` argument to a fixture file.
make_stub_gh() {
  local out="$1"; shift
  {
    echo '#!/usr/bin/env bash'
    echo 'set -u'
    echo 'sha=""'
    echo 'for arg in "$@"; do'
    echo '  case "$arg" in'
    echo '    repos/*/commits/*/check-runs)'
    echo '      sha="${arg#repos/*/commits/}"'
    echo '      sha="${sha%/check-runs}"'
    echo '      ;;'
    echo '  esac'
    echo 'done'
    echo 'case "$sha" in'
    for pair in "$@"; do
      s="${pair%%:*}"
      f="${pair#*:}"
      echo "  $s) cat \"$f\" ;;"
    done
    echo '  *) echo "stub gh: no fixture for $sha" >&2; exit 1 ;;'
    echo 'esac'
  } > "$out"
  chmod +x "$out"
}

# make_throwaway_repo <dir> — git init a repo with a base commit; prints the BASE sha.
make_throwaway_repo() {
  local dir="$1"
  mkdir -p "$dir"
  git -C "$dir" init -q -b main
  git -C "$dir" -c user.name="Test" -c user.email="test@example.com" commit -q --allow-empty -m base
  git -C "$dir" rev-parse HEAD
}

# commit_in <dir> <file> <content> <msg> — add a file and commit; prints the new HEAD sha.
commit_in() {
  local dir="$1" file="$2" content="$3" msg="$4"
  printf '%s\n' "$content" > "$dir/$file"
  git -C "$dir" add "$file"
  git -C "$dir" -c user.name="Test" -c user.email="test@example.com" commit -q -m "$msg"
  git -C "$dir" rev-parse HEAD
}

# ── arms ─────────────────────────────────────────────────────────────────────
# AC-M1-2 — a linear repo whose non-HEAD commits all have check sets; each is reported :2, GREEN.
arm_check_count() {
  local dir base a b c stub out rc
  dir="$(mktemp -d)"
  base="$(make_throwaway_repo "$dir")"
  a="$(commit_in "$dir" a.txt "a" "commit a")"
  b="$(commit_in "$dir" b.txt "b" "commit b")"
  c="$(commit_in "$dir" c.txt "c" "commit c")"
  stub="$(mktemp)"
  make_stub_gh "$stub" "$a:$VERIFIED_FIXTURE" "$b:$VERIFIED_FIXTURE" "$c:$VERIFIED_FIXTURE"
  out="$(/bin/bash "$SCRIPT_UT" --repo-dir "$dir" --repo test/test --base "$base" --head "$c" --gh-bin "$stub" 2>&1)"
  rc=$?
  if [ "$rc" -eq 0 ]; then ok "check-count exits 0"; else notok "check-count rc=$rc: $out"; fi
  echo "$out" | /usr/bin/grep -q ':2' && ok "check-count reports :2" || notok "check-count no :2: $out"
  rm -rf "$dir" "$stub"
}

# AC-M1-3 — a zero-check non-HEAD commit is reported ZERO-CHECK-CANDIDATE, exit 1.
arm_candidate() {
  local dir base a b stub out rc
  dir="$(mktemp -d)"
  base="$(make_throwaway_repo "$dir")"
  a="$(commit_in "$dir" a.txt "a" "commit a")"
  b="$(commit_in "$dir" b.txt "b" "commit b")"
  stub="$(mktemp)"
  make_stub_gh "$stub" "$a:$ZERO_FIXTURE" "$b:$VERIFIED_FIXTURE"
  out="$(/bin/bash "$SCRIPT_UT" --repo-dir "$dir" --repo test/test --base "$base" --head "$b" --gh-bin "$stub" 2>&1)"
  rc=$?
  if [ "$rc" -eq 1 ]; then ok "candidate exits 1"; else notok "candidate rc=$rc: $out"; fi
  echo "$out" | /usr/bin/grep -q 'ZERO-CHECK-CANDIDATE' && ok "candidate names ZERO-CHECK-CANDIDATE" || notok "candidate no ZERO-CHECK-CANDIDATE: $out"
  rm -rf "$dir" "$stub"
}

# AC-M1-9 — the range's head is never reported, on a LINEAR history; GREEN.
#
# The range is deliberately 2+ commits (a, then head), NOT a single commit. A one-commit
# range is degenerate: base+1 == head, so under AC-M1-10's `base..head^` mutation the range
# collapses to empty and this arm exits 3 -- i.e. it would red for a reason that has nothing
# to do with head exclusion, and AC-M1-10's independence from this arm could not be shown.
# Found by the independent evaluator (iteration 170) and reproduced first-party by the
# controller before this widening; see the design doc's AC-M1-9 note.
arm_head_excluded() {
  local dir base a b stub out rc
  dir="$(mktemp -d)"
  base="$(make_throwaway_repo "$dir")"
  a="$(commit_in "$dir" a.txt "a" "commit a")"
  b="$(commit_in "$dir" b.txt "b" "commit b")"
  stub="$(mktemp)"
  make_stub_gh "$stub" "$a:$VERIFIED_FIXTURE" "$b:$ZERO_FIXTURE"
  out="$(/bin/bash "$SCRIPT_UT" --repo-dir "$dir" --repo test/test --base "$base" --head "$b" --gh-bin "$stub" 2>&1)"
  rc=$?
  if [ "$rc" -eq 0 ]; then ok "head-excluded exits 0"; else notok "head-excluded rc=$rc: $out"; fi
  echo "$out" | /usr/bin/grep -q 'ZERO-CHECK-CANDIDATE' && notok "head-excluded reported the head as a candidate" || ok "head-excluded head not reported"
  echo "$out" | /usr/bin/grep -q "${a:0:9}" && ok "head-excluded classifies the non-head commit" || notok "head-excluded never reached the loop: $out"
  rm -rf "$dir" "$stub"
}

# AC-M1-10 — a MERGE head with a zero-check on a SECOND-parent commit: both parent histories
# examined, the exact merge head excluded, the second-parent candidate surfaced, exit 1.
arm_merge_head_coverage() {
  local dir base m1 m2 s1 s2 head stub out rc
  dir="$(mktemp -d)"
  base="$(make_throwaway_repo "$dir")"
  m1="$(commit_in "$dir" m1.txt "m1" "main 1")"
  m2="$(commit_in "$dir" m2.txt "m2" "main 2")"
  git -C "$dir" checkout -q -b side "$base"
  s1="$(commit_in "$dir" s1.txt "s1" "side 1")"
  s2="$(commit_in "$dir" s2.txt "s2" "side 2")"
  git -C "$dir" checkout -q main
  git -C "$dir" -c user.name="Test" -c user.email="test@example.com" merge -q --no-ff side -m "merge side"
  head="$(git -C "$dir" rev-parse HEAD)"
  stub="$(mktemp)"
  make_stub_gh "$stub" "$m1:$VERIFIED_FIXTURE" "$m2:$VERIFIED_FIXTURE" "$s1:$VERIFIED_FIXTURE" "$s2:$ZERO_FIXTURE" "$head:$VERIFIED_FIXTURE"
  out="$(/bin/bash "$SCRIPT_UT" --repo-dir "$dir" --repo test/test --base "$base" --head "$head" --gh-bin "$stub" 2>&1)"
  rc=$?
  if [ "$rc" -eq 1 ]; then ok "merge-head-coverage exits 1"; else notok "merge-head-coverage rc=$rc: $out"; fi
  echo "$out" | /usr/bin/grep -q 'ZERO-CHECK-CANDIDATE' && ok "merge-head-coverage surfaces second-parent candidate" || notok "merge-head-coverage no candidate: $out"
  echo "$out" | /usr/bin/grep -q "$head" && notok "merge-head-coverage reported the merge head" || ok "merge-head-coverage merge head excluded"
  rm -rf "$dir" "$stub"
}

# AC-M1-5 twin — an unreachable gh binary is a LOUD-FAIL (exit 2) naming the API.
arm_gh_missing() {
  local dir base a b out rc
  dir="$(mktemp -d)"
  base="$(make_throwaway_repo "$dir")"
  a="$(commit_in "$dir" a.txt "a" "commit a")"
  b="$(commit_in "$dir" b.txt "b" "commit b")"
  out="$(/bin/bash "$SCRIPT_UT" --repo-dir "$dir" --repo test/test --base "$base" --head "$b" --gh-bin /nonexistent/gh 2>&1)"
  rc=$?
  if [ "$rc" -eq 2 ]; then ok "gh-missing exits 2"; else notok "gh-missing rc=$rc: $out"; fi
  echo "$out" | /usr/bin/grep -q 'unreachable' && ok "gh-missing names unreachable" || notok "gh-missing no unreachable: $out"
  rm -rf "$dir"
}

# AC-M1-6 twin — a hung gh call is bounded and maps to a LOUD-FAIL (exit 2) naming a timeout.
arm_gh_timeout() {
  local dir base a b stub out rc
  dir="$(mktemp -d)"
  base="$(make_throwaway_repo "$dir")"
  a="$(commit_in "$dir" a.txt "a" "commit a")"
  b="$(commit_in "$dir" b.txt "b" "commit b")"
  stub="$(mktemp)"
  printf '#!/usr/bin/env bash\n/bin/sleep 3\n' > "$stub"
  chmod +x "$stub"
  out="$(/bin/bash "$SCRIPT_UT" --repo-dir "$dir" --repo test/test --base "$base" --head "$b" --gh-bin "$stub" --gh-timeout 0.5 2>&1)"
  rc=$?
  if [ "$rc" -eq 2 ]; then ok "gh-timeout exits 2"; else notok "gh-timeout rc=$rc: $out"; fi
  echo "$out" | /usr/bin/grep -q 'timed out' && ok "gh-timeout names timeout" || notok "gh-timeout no timeout msg: $out"
  rm -rf "$dir" "$stub"
}

# AC-M2-2 — the report format: the per-commit line names ZERO-CHECK-CANDIDATE, the exit-1
# summary reads "a controller must look: N zero-check candidate(s)", and the output does NOT
# contain the string "unverified" (the instrument surfaces; it does not adjudicate).
arm_report() {
  local dir base a b stub out rc
  dir="$(mktemp -d)"
  base="$(make_throwaway_repo "$dir")"
  a="$(commit_in "$dir" a.txt "a" "commit a")"
  b="$(commit_in "$dir" b.txt "b" "commit b")"
  stub="$(mktemp)"
  make_stub_gh "$stub" "$a:$ZERO_FIXTURE" "$b:$VERIFIED_FIXTURE"
  out="$(/bin/bash "$SCRIPT_UT" --repo-dir "$dir" --repo test/test --base "$base" --head "$b" --gh-bin "$stub" 2>&1)"
  rc=$?
  if [ "$rc" -eq 1 ]; then ok "report exits 1"; else notok "report rc=$rc: $out"; fi
  echo "$out" | /usr/bin/grep -q 'ZERO-CHECK-CANDIDATE' && ok "report per-commit line names ZERO-CHECK-CANDIDATE" || notok "report no ZERO-CHECK-CANDIDATE: $out"
  echo "$out" | /usr/bin/grep -q 'a controller must look: 1 zero-check candidate(s)' && ok "report summary reads 'a controller must look'" || notok "report summary wrong: $out"
  echo "$out" | /usr/bin/grep -q 'unverified' && notok "report contains 'unverified'" || ok "report does not claim 'unverified'"
  rm -rf "$dir" "$stub"
}

# ── dispatcher ───────────────────────────────────────────────────────────────
run_arm() {
  case "$1" in
    check-count) arm_check_count ;;
    candidate) arm_candidate ;;
    head-excluded) arm_head_excluded ;;
    merge-head-coverage) arm_merge_head_coverage ;;
    gh-missing) arm_gh_missing ;;
    gh-timeout) arm_gh_timeout ;;
    report) arm_report ;;
    *) echo "unknown arm: $1" >&2; exit 2 ;;
  esac
}

ONLY=""
if [ "$#" -gt 0 ]; then
  if [ "$1" = "--only" ]; then
    ONLY="$2"
  else
    echo "unknown argument: $1" >&2
    exit 2
  fi
fi

if [ -n "$ONLY" ]; then
  run_arm "$ONLY"
else
  for arm in check-count candidate head-excluded merge-head-coverage gh-missing gh-timeout report; do
    run_arm "$arm"
  done
fi

echo "---"
echo "$pass passed, $fail failed"
[ "$fail" -eq 0 ]
