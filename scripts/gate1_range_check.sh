#!/usr/bin/env bash
# gate1_range_check.sh — enumerate last-gate1..origin/dev, drop the head by SHA equality,
# read one check set per remaining commit, and report every total_count==0 commit as a
# ZERO-CHECK-CANDIDATE (exit 1).
#
# Exit 1 means "a controller must look", NEVER "a commit is unverified". The instrument
# surfaces candidates; it does not adjudicate push tips (that is queue row 88, out of scope).
#
# This is a bash-3.2 script (the rig's /bin/bash is 3.2.57). It must NOT `cd "$(dirname
# "$0")/.."` the way verify_ail.sh does — that would pin it to the repo root and make every
# throwaway-repo arm untestable. It operates on --repo-dir (default `.`).
set -uo pipefail

# ── defaults ──────────────────────────────────────────────────────────────────
REPO_DIR="."
GH_BIN="gh"
GH_TIMEOUT="5"
BASE=""
HEAD=""
REPO=""
CHECK_RUNS_FILE=""

usage() {
  cat <<'EOF'
Usage: gate1_range_check.sh [options]

Enumerate the range <base>..<head>, drop the head by SHA equality (never $head^), read one
check set per remaining commit, and report every total_count==0 commit as a
ZERO-CHECK-CANDIDATE. Exit 1 means "a controller must look", never "a commit is unverified".

Options:
  --base <sha>          lower bound (default: last base-gate1 row in ~/.ailang/state/mission-world-base)
  --head <sha>          upper bound (default: origin/dev)
  --repo <owner/name>   repo slug for the check-runs API (default: derived from the origin remote)
  --repo-dir <path>     git repo directory (default: .)
  --gh-bin <path>       gh binary (default: gh)
  --gh-timeout <secs>   per-call timeout in seconds (default: 5)
  --check-runs-file <f> PARSER-ONLY mode: parse one check-runs response, print total_count, exit 0/2
  --help                show this help

Exit codes:
  0  GREEN      every non-HEAD commit is VERIFIED
  1  CANDIDATE  one or more ZERO-CHECK-CANDIDATEs found (a controller must look)
  2  LOUD-FAIL  malformed/missing total_count, unreachable API, gh timeout, or missing lower bound
  3  EMPTY      the range is empty (no new commits since the previous fire)
EOF
}

# ── bounded subprocess wrapper (bash 3.2; no timeout(1)/gtimeout on this rig) ──
# Background the child in a subshell, start a watchdog that sleeps then TERM then KILL,
# wait the child, and map a signal-kill (rc >= 128) to 124 (TIMEOUT). The subshell + the
# `wait ... 2>/dev/null` suppress the job-control "Terminated: 15" line that would otherwise
# pollute the report. Measured outcomes: sleep-under-bound -> 124; echo -> 0; /nonexistent/gh
# -> 127; `sh -c 'exit 7'` -> 7.
run_bounded() {  # $1=timeout_s  $2=out_file  $3..=cmd ; returns 124 on timeout, child rc otherwise
  local t="$1" out="$2"; shift 2
  local child rc watchdog
  ( "$@" >"$out" 2>"$out.err" ) &
  child=$!
  ( /bin/sleep "$t"; kill -TERM "$child" 2>/dev/null; /bin/sleep 0.2; kill -KILL "$child" 2>/dev/null ) &
  watchdog=$!
  wait "$child" 2>/dev/null; rc=$?
  kill "$watchdog" 2>/dev/null
  wait "$watchdog" 2>/dev/null
  if [ "$rc" -ge 128 ]; then return 124; fi
  return "$rc"
}

# ── python3 total_count parser (house pattern: verify_ail.sh parses JSON with python3) ──
# Prints `total_count=<n>` on a valid parse, exits 2 on an invalid one. Valid zero cardinality
# ({"total_count":0,"check_runs":[]}) is ACCEPTED — it is a true zero, not a parse failure.
# The consistency rule is one-directional (the endpoint paginates at 30): fail only when
# (total_count==0 AND len(check_runs)>0) or (total_count>0 AND len(check_runs)==0).
parse_check_runs() {  # $1=file
  python3 - "$1" <<'PY'
import json, sys
try:
    with open(sys.argv[1]) as fh:
        d = json.load(fh)
except Exception as e:
    sys.stderr.write("missing or invalid total_count: could not parse JSON (%s)\n" % e)
    sys.exit(2)

tc = d.get("total_count")
if tc is None:
    sys.stderr.write("missing or invalid total_count: no total_count key\n")
    sys.exit(2)
if isinstance(tc, bool) or not isinstance(tc, int):
    sys.stderr.write("missing or invalid total_count: total_count is not an integer\n")
    sys.exit(2)
if tc < 0:
    sys.stderr.write("missing or invalid total_count: total_count is negative\n")
    sys.exit(2)

if "check_runs" in d:
    n = len(d["check_runs"])
    if (tc == 0 and n > 0) or (tc > 0 and n == 0):
        sys.stderr.write("missing or invalid total_count: inconsistent with check_runs array\n")
        sys.exit(2)

print("total_count=%d" % tc)
sys.exit(0)
PY
}

# ── lower bound: last base-gate1 row in the mission-world-base TSV, else mission-base.sh ──
resolve_base() {
  local f="$HOME/.ailang/state/mission-world-base"
  if [ -f "$f" ]; then
    local sha
    sha="$(/usr/bin/awk -F'\t' '$3=="base-gate1" {last=$5} END {print last}' "$f")"
    if [ -n "$sha" ]; then
      printf '%s\n' "$sha"
      return 0
    fi
  fi
  if [ -x /Users/voightkampff/dev/sunholo-data/ailang/tools/launchd/mission-base.sh ]; then
    local s
    s="$(bash /Users/voightkampff/dev/sunholo-data/ailang/tools/launchd/mission-base.sh last gate1 2>/dev/null)"
    if [ -n "$s" ]; then
      printf '%s\n' "$s"
      return 0
    fi
  fi
  return 1
}

# ── repo slug: --repo, else derived from the origin remote ──
derive_repo() {
  local origin slug
  origin="$(git -C "$REPO_DIR" remote get-url origin 2>/dev/null)" || return 1
  [ -n "$origin" ] || return 1
  slug="$(printf '%s\n' "$origin" | sed -E 's#^.*[:/]([^/]+/[^/]+)(\.git)?$#\1#; s#\.git$##')"
  [ -n "$slug" ] || return 1
  printf '%s\n' "$slug"
}

# ── argument parsing ──────────────────────────────────────────────────────────
while [ "$#" -gt 0 ]; do
  case "$1" in
    --base) BASE="$2"; shift 2 ;;
    --head) HEAD="$2"; shift 2 ;;
    --repo) REPO="$2"; shift 2 ;;
    --repo-dir) REPO_DIR="$2"; shift 2 ;;
    --gh-bin) GH_BIN="$2"; shift 2 ;;
    --gh-timeout) GH_TIMEOUT="$2"; shift 2 ;;
    --check-runs-file) CHECK_RUNS_FILE="$2"; shift 2 ;;
    --help) usage; exit 0 ;;
    *) echo "✗ unknown argument: $1" >&2; usage >&2; exit 2 ;;
  esac
done

# ── PARSER-ONLY mode ──────────────────────────────────────────────────────────
# Parses one check-runs response, prints total_count=<n>, exits 0 on a valid parse and 2 on
# an invalid one. Performs NO classification and NO enumeration (planner resolution P2).
if [ -n "$CHECK_RUNS_FILE" ]; then
  parse_check_runs "$CHECK_RUNS_FILE"
  exit $?
fi

# ── resolve base and head ─────────────────────────────────────────────────────
if [ -z "$BASE" ]; then
  BASE="$(resolve_base)" || {
    echo "✗ no lower bound: no base-gate1 row in mission-world-base and mission-base.sh unavailable" >&2
    exit 2
  }
fi
if [ -z "$HEAD" ]; then
  HEAD="$(git -C "$REPO_DIR" rev-parse origin/dev 2>/dev/null)" || {
    echo "✗ cannot resolve origin/dev in $REPO_DIR" >&2
    exit 2
  }
fi

# ── resolve repo slug ────────────────────────────────────────────────────────
if [ -z "$REPO" ]; then
  REPO="$(derive_repo)" || {
    echo "✗ no --repo and no origin remote — cannot derive the repo slug" >&2
    exit 2
  }
fi

# ── enumerate the range ──────────────────────────────────────────────────────
commits="$(git -C "$REPO_DIR" rev-list "$BASE..$HEAD" 2>/dev/null)"
if [ $? -ne 0 ]; then
  echo "✗ invalid revision range $BASE..$HEAD" >&2
  exit 2
fi

# ── EMPTY range (exit 3) ────────────────────────────────────────────────────
if [ -z "$commits" ]; then
  echo "no new commits since $BASE; nothing to check"
  exit 3
fi

# ── classify each commit, excluding the head by SHA equality (never $head^) ──
candidates=0
for sha in $commits; do
  if [ "$sha" = "$HEAD" ]; then
    continue
  fi
  out="$(mktemp)"
  run_bounded "$GH_TIMEOUT" "$out" "$GH_BIN" api "repos/$REPO/commits/$sha/check-runs"
  rc=$?
  if [ "$rc" -eq 124 ]; then
    echo "✗ gh api timed out after ${GH_TIMEOUT}s for $sha" >&2
    rm -f "$out" "$out.err"
    exit 2
  fi
  if [ "$rc" -eq 127 ]; then
    echo "✗ gh binary unreachable: $GH_BIN" >&2
    rm -f "$out" "$out.err"
    exit 2
  fi
  if [ "$rc" -ne 0 ]; then
    echo "✗ gh api failed (rc=$rc) for $sha" >&2
    rm -f "$out" "$out.err"
    exit 2
  fi
  tc_line="$(parse_check_runs "$out")"
  prc=$?
  rm -f "$out" "$out.err"
  if [ "$prc" -ne 0 ]; then
    echo "✗ malformed check-runs response for $sha" >&2
    exit 2
  fi
  tc="${tc_line#total_count=}"
  if [ "$tc" -gt 0 ]; then
    echo "$sha VERIFIED :$tc"
  else
    echo "$sha ZERO-CHECK-CANDIDATE :0"
    candidates=$((candidates + 1))
  fi
done

# ── verdict ─────────────────────────────────────────────────────────────────
if [ "$candidates" -gt 0 ]; then
  echo "a controller must look: $candidates zero-check candidate(s)"
  exit 1
fi
echo "GREEN: all non-HEAD commits verified"
exit 0
