#!/usr/bin/env bash
set -euo pipefail

tmp=$(mktemp -d "${TMPDIR:-/tmp}/race-leg.XXXXXX")
cleanup() {
  local pid
  if [ -s "$tmp/group.pid" ]; then
    pid=$(cat "$tmp/group.pid")
    kill -KILL -- "-$pid" 2>/dev/null || true
  fi
  rm -rf "$tmp"
}
trap cleanup EXIT

gate=(bash scripts/verify_go.sh)
expected='race bounds: per-package 900s, start allowance 995s, leg 2000s'
printf '%s\n' "$expected" >"$tmp/expected"
order='must be at least per-package 900s + 60s backup kill + start allowance 995s (1955s)'
fail() { echo "FAIL $arm: $* (rc=$rc)" >&2; cat "$tmp/out" "$tmp/err" >&2; exit 1; }
run() {
  rc=0
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
cat > "$tmp/hang_probe.test" <<'PROBE'
#!/bin/sh
echo "$$" > "$RACE_PROBE_DIR/group.pid"
sleep 300 &
echo "$!" > "$RACE_PROBE_DIR/gc.pid"
sleep 300
PROBE
chmod +x "$tmp/hang_probe.test"
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
      rm -f "$tmp/group.pid" "$tmp/gc.pid"
      return
    fi
    sleep 0.1
  done
  fail "grandchild $pid still alive after reap poll"
}
check_timeout() {
  [ "$rc" -eq 124 ] && grep -qF 'FATAL: -race leg timed out after 2s' "$tmp/err" || fail 'expected rc 124 and timeout diagnostic'
}

arm='R6 hang'
SECONDS=0
run watchdog 24 "${gate[@]}" --race-leg-wrapper-check 2 "$tmp/hang_probe.test"
elapsed=$SECONDS
check_timeout
check_reaped
grep -F 'still running at the leg bound' "$tmp/err" | grep -qF 'hang_probe.test' || fail 'expected a process listing naming the probe'
if grep -qF -- '--race-leg-wrapper-check' "$tmp/err"; then
  fail 'listing included the wrapper outside the probe group'
fi
[ "$elapsed" -ge 2 ] && [ "$SECONDS" -lt 30 ] || fail "wall bound: wrapper ${elapsed}s, arm ${SECONDS}s"
echo "ok $arm (wrapper ${elapsed}s, arm ${SECONDS}s)"

arm='R7 stalled ps'
printf '#!/bin/sh\nexec sleep 300\n' > "$tmp/ps"
chmod +x "$tmp/ps"
SECONDS=0
run watchdog 10 env PATH="$tmp:$PATH" "${gate[@]}" --race-leg-wrapper-check 2 "$tmp/hang_probe.test"
elapsed=$SECONDS
check_timeout
grep -qF 'could not list' "$tmp/err" || fail 'expected a failed listing diagnostic'
check_reaped
[ "$elapsed" -ge 2 ] && [ "$SECONDS" -lt 10 ] || fail "wall bound: wrapper ${elapsed}s, arm ${SECONDS}s"
echo "ok $arm (wrapper ${elapsed}s, arm ${SECONDS}s)"
echo 'PASS 7/7 arms'
