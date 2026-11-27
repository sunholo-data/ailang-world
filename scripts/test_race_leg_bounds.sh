#!/usr/bin/env bash
set -euo pipefail

tmp=$(mktemp -d "${TMPDIR:-/tmp}/race-leg.XXXXXX")
cleanup_probe() {
  local pid
  if [ -s "$tmp/group.pid" ]; then
    pid=$(cat "$tmp/group.pid")
    kill -KILL -- "-$pid" 2>/dev/null || true
  fi
  rm -f "$tmp/group.pid" "$tmp/gc.pid" "$tmp/ready"
}
cleanup() {
  cleanup_probe
  rm -rf "$tmp"
}
trap cleanup EXIT

gate=(bash scripts/verify_go.sh)
expected='race bounds: per-package 900s, start allowance 995s, leg 2000s'
printf '%s\n' "$expected" >"$tmp/expected"
order='must be at least per-package 900s + 60s backup kill + start allowance 995s (1955s)'
fail() { echo "FAIL $arm: $* (rc=$rc, elapsed=${SECONDS}s)" >&2; cat "$tmp/out" "$tmp/err" >&2; exit 1; }
run() {
  rc=0
  SECONDS=0
  "$@" >"$tmp/out" 2>"$tmp/err" || rc=$?
}
ok() { echo "ok $arm"; }

arm='R1 default'
run env -u VERIFY_GO_RACE_TIMEOUT_S "${gate[@]}" --race-bounds-check
[ "$rc" -eq 0 ] && cmp -s "$tmp/out" "$tmp/expected" || fail "$expected"
ok
arm='R2 boundary accept'
run env VERIFY_GO_RACE_TIMEOUT_S=1955 "${gate[@]}" --race-bounds-check
[ "$rc" -eq 0 ] || fail 'expected rc 0 at 1955s'
ok
arm='R3 boundary refuse'
run env VERIFY_GO_RACE_TIMEOUT_S=1954 "${gate[@]}" --race-bounds-check
[ "$rc" -eq 2 ] && grep -qF "$order" "$tmp/err" || fail 'expected rc 2 and the bound calculation'
ok
arm='R4 malformed'
run env VERIFY_GO_RACE_TIMEOUT_S=15m "${gate[@]}" --race-bounds-check
[ "$rc" -eq 2 ] && grep -qF 'must be a positive integer' "$tmp/err" || fail 'expected rc 2 and the integer diagnostic'
ok
arm='R5 passthrough'
run "${gate[@]}" --race-leg-wrapper-check 30 sh -c 'exit 3'
[ "$rc" -eq 3 ] && ! grep -qF FATAL "$tmp/err" || fail 'expected child rc 3 without FATAL'
ok

export RACE_PROBE_DIR="$tmp"
# Keep the command on one line so ps lists it with the shell's $0 marker.
printf -v hang_probe '%s ' \
  "trap '' TERM;" \
  'echo "$$" > "$RACE_PROBE_DIR/group.pid";' \
  'sleep 300 &' \
  'echo "$!" > "$RACE_PROBE_DIR/gc.pid";' \
  ': > "$RACE_PROBE_DIR/ready";' \
  'wait'
watchdog() {
  python3 - "$@" <<'PYTHON'
import os, signal, subprocess, sys
bound = int(sys.argv[1])
p = subprocess.Popen(sys.argv[2:], start_new_session=True)
try:
    sys.exit(p.wait(timeout=bound))
except subprocess.TimeoutExpired:
    os.killpg(p.pid, signal.SIGKILL)
    p.wait()
    sys.exit(f"watchdog fired after {bound}s")
PYTHON
}
check_reaped() {
  [ -s "$tmp/gc.pid" ] || fail 'grandchild PID control missing'
  local pid i
  pid=$(cat "$tmp/gc.pid")
  for ((i=0; i<50; i++)); do
    if ! kill -0 "$pid" 2>/dev/null; then
      rm -f "$tmp/group.pid" "$tmp/gc.pid" "$tmp/ready"
      return
    fi
    sleep 0.1
  done
  fail "grandchild $pid still alive after reap poll"
}
check_timeout() {
  [ "$rc" -eq 124 ] && grep -qF "FATAL: -race leg timed out after ${hang_bound}s" "$tmp/err" || fail 'expected rc 124 and timeout diagnostic'
  [ "$elapsed" -ge "$hang_bound" ] && [ "$elapsed" -lt "$wall_upper_bound" ] ||
    fail "wrapper elapsed ${elapsed}s, expected ${hang_bound}s to under ${wall_upper_bound}s"
}

# The probe writes ready after recording its child PID. The TERM trap makes
# child survival detectable if the wrapper sends SIGTERM.
hang_bound=10
# The exclusive upper wall bound adds 2s of ps grace and 8s of scheduling time.
readonly wall_upper_bound=$((hang_bound + 2 + 8))
watchdog_bound=$((hang_bound + 40))

# The wrapper starts its clock at launch, so retry only when ready is absent,
# up to three attempts. Every product assertion runs on a ready attempt, and
# any assertion failure exits immediately. A mutant that prevents readiness
# still fails after three attempts.
run_hang() {
  local attempt
  for attempt in 1 2 3; do
    run watchdog "$watchdog_bound" "$@" "${gate[@]}" --race-leg-wrapper-check "$hang_bound" \
      /bin/sh -c "$hang_probe" hang_probe.test
    elapsed=$SECONDS
    [ -f "$tmp/ready" ] && return 0
    echo "retry $arm: probe not ready before the ${hang_bound}s bound (attempt $attempt/3, rc=$rc, elapsed=${elapsed}s)" >&2
    cat "$tmp/err" >&2
    cleanup_probe
  done
  fail 'probe never signaled readiness'
}

arm='R6 hang'
run_hang
check_timeout
check_reaped
grep -F 'still running at the leg bound' "$tmp/err" | grep -qF 'hang_probe.test' || fail 'expected a process listing naming the probe'
if grep -qF -- '--race-leg-wrapper-check' "$tmp/err"; then
  fail 'listing included the wrapper outside the probe group'
fi
echo "ok $arm (wrapper ${elapsed}s, arm ${SECONDS}s)"

arm='R7 stalled ps'
printf '#!/bin/sh\nexec sleep 300\n' > "$tmp/ps"
chmod +x "$tmp/ps"
run_hang env PATH="$tmp:$PATH"
check_timeout
grep -qF 'could not list' "$tmp/err" || fail 'expected a failed listing diagnostic'
check_reaped
echo "ok $arm (wrapper ${elapsed}s, arm ${SECONDS}s)"
echo 'PASS 7/7 arms'
