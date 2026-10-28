#!/usr/bin/env bash
# Capture frozen inventory around one or two supervised six-shot passes.
set -euo pipefail
script_dir=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)
supervisor_pid=

cleanup() {
    local status=$?
    trap - EXIT INT TERM HUP
    if [[ -n "$supervisor_pid" ]]; then
        kill -TERM "$supervisor_pid" 2>/dev/null || :
        # Python owns the daemon: TERM <=5s, then group KILL/reap <=5s.
        if wait "$supervisor_pid"; then :; else status=1; fi
    fi
    exit "$status"
}
trap cleanup EXIT
trap 'exit 130' INT
trap 'exit 143' TERM
trap 'exit 129' HUP

if [[ ${1:-} == --help ]]; then
    cat <<'HELP'
Usage: scripts/capture_marketing.sh [capture configuration]
  --db PATH                 Frozen demo database
  --manifest PATH           Seed manifest
  --proof-key-file PATH     Proof verification key file
  --ailang-bin PATH         Pinned kernel AILANG executable
  --workspace-root PATH     Demo workspace root
  --tool-ailang-bin PATH    Pinned toolchain AILANG executable
  --chrome PATH             Explicit Chrome executable
  --out PATH                Scratch capture output directory
  --verify-twice            Capture and compare two complete passes
  --rectangles PATH         Controller-measured content rectangles
  --worldd-bin PATH         ailang-worldd executable
  --bind 127.0.0.1:PORT      Loopback daemon address
  --marketingcheck-bin PATH Prebuilt capture verifier

No session mint, commit, approval or background agent is performed.

Supervision primitive: scripts/capture_marketing.sh shot [shot flags]
See python3 scripts/marketing_capture.py shot --help.

Unit control only (no real daemon/browser or socket):
  --unit-daemon exit|wait DAEMON READY_FILE STARTED_FILE TERM_WAIT KILL_WAIT
The daemon writes READY_FILE; Python publishes STARTED_FILE with its PID.
'exit' tests normal EXIT cleanup; 'wait' awaits a signal. Waits are seconds.
Production cleanup waits are 5s TERM + 5s KILL/reap; shot 30s; pass 300s.
HELP
    exit 0
fi

if [[ ${1:-} == shot ]]; then
    shift
    python3 "$script_dir/marketing_capture.py" shot "$@"
    exit $?
fi

if [[ ${1:-} == --unit-daemon ]]; then
    if [[ $# != 7 || ( $2 != exit && $2 != wait ) ]]; then
        echo 'unit control requires: --unit-daemon exit|wait DAEMON READY STARTED TERM_WAIT KILL_WAIT' >&2
        exit 2
    fi
    mode=$2
    started_file=$5
    python3 "$script_dir/marketing_capture.py" daemon \
        --ready-file "$4" --started-file "$started_file" \
        --term-wait "$6" --kill-wait "$7" -- "$3" &
    supervisor_pid=$!
    # Bound shell readiness too. The supervisor separately bounds daemon readiness.
    deadline=$((SECONDS + 5))
    while [[ ! -s "$started_file" ]]; do
        if ! kill -0 "$supervisor_pid" 2>/dev/null || (( SECONDS >= deadline )); then
            echo 'unit daemon failed to become ready within shell bound' >&2
            exit 1
        fi
        sleep 0.01
    done
    if [[ $mode == wait ]]; then
        wait "$supervisor_pid"
        supervisor_pid=
    fi
    exit 0
fi

fail() { echo "capture: $*" >&2; exit 1; }
db= manifest= proof_key_file= ailang_bin= workspace_root= tool_ailang_bin=
chrome= out= rectangles= worldd_bin= bind= marketingcheck_bin= verify_twice=
while (( $# )); do
    case $1 in
        --verify-twice) verify_twice=yes; shift ;;
        --db|--manifest|--proof-key-file|--ailang-bin|--workspace-root|--tool-ailang-bin|--chrome|--out|--rectangles|--worldd-bin|--bind|--marketingcheck-bin)
            (( $# >= 2 )) || fail "missing value for $1"
            flag=${1#--}; flag=${flag//-/_}; printf -v "$flag" '%s' "$2"; shift 2 ;;
        *) fail "unknown flag: $1" ;;
    esac
done
for flag in db manifest proof_key_file ailang_bin workspace_root tool_ailang_bin chrome out rectangles worldd_bin bind marketingcheck_bin; do
    [[ -n ${!flag} ]] || fail "missing --${flag//_/-}"
done
[[ $bind =~ ^127\.0\.0\.1:([0-9]+)$ ]] || fail "invalid loopback bind: $bind"
port=${BASH_REMATCH[1]}
[[ ${#port} -le 5 ]] && (( 10#$port > 0 && 10#$port <= 65535 )) || fail "invalid loopback bind: $bind"
mkdir -p -- "$out" || fail 'output directory'
supervisor="$script_dir/marketing_capture.py"
# Read the production constants rather than introducing shell-specific deadlines.
read -r startup lifetime < <(python3 -B -c 'import runpy, sys; m=runpy.run_path(sys.argv[1]); print(int(m["STARTUP_TIMEOUT_S"]), int(2*m["SHOT_TIMEOUT_S"] + 2*m["PASS_TIMEOUT_S"]))' "$supervisor")
# Lifetime = version A + pass A + pass B + version B = 2*30 + 2*300 = 660 s (A1).
step() {
    local name=$1; shift
    "$@" || fail "$name failed"
}
version_a=$(step 'chrome-version A' python3 "$supervisor" chrome-version --chrome "$chrome") || exit 1
started_file="$out/daemon-started.json"
python3 "$supervisor" daemon --ready-url "http://$bind" --started-file "$started_file" --lifetime "$lifetime" -- \
    "$worldd_bin" serve --db "$db" --bind "$bind" --workbench-still \
    --ailang-bin "$ailang_bin" --workspace-root "$workspace_root" --tool-ailang-bin "$tool_ailang_bin" \
    --proof-key-file "$proof_key_file" --proof-required-identity sameRef --proof-required-identity servesEntry \
    > "$out/daemon-verdict.json" &
supervisor_pid=$!
deadline=$((SECONDS + startup))
while [[ ! -s "$started_file" ]]; do
    if ! kill -0 "$supervisor_pid" 2>/dev/null || (( SECONDS >= deadline )); then
        fail 'daemon readiness failed within shell bound'
    fi
    sleep 0.01
done
step 'inventory before' "$marketingcheck_bin" inventory --addr "http://$bind" --out "$out/inventory-before.json"
step routes "$marketingcheck_bin" routes --manifest "$manifest" > "$out/routes.json"
step 'pass A' python3 "$supervisor" pass --routes "$out/routes.json" --base-url "http://$bind" --chrome "$chrome" --out-dir "$out/pass-a"
if [[ $verify_twice == yes ]]; then
    step 'pass B' python3 "$supervisor" pass --routes "$out/routes.json" --base-url "http://$bind" --chrome "$chrome" --out-dir "$out/pass-b"
    version_b=$(step 'chrome-version B' python3 "$supervisor" chrome-version --chrome "$chrome") || exit 1
    [[ $version_a == "$version_b" ]] || fail 'chrome version changed during capture'
fi
step 'inventory after' "$marketingcheck_bin" inventory --addr "http://$bind" --out "$out/inventory-after.json"
step 'inventory compare' "$marketingcheck_bin" inventory-compare --before "$out/inventory-before.json" --after "$out/inventory-after.json"
[[ $verify_twice == yes ]] || exit 0
version=$(python3 -c 'import json,sys; print(json.loads(sys.argv[1])["version"])' "$version_a")
printf 'chrome-version: %s\nbind: %s\ninventory: equal\n' "$version" "$bind" > "$out/diagnostics.txt"
step check "$marketingcheck_bin" check --pass-a "$out/pass-a" --pass-b "$out/pass-b" --manifest "$manifest" --rectangles "$rectangles" --out "$out/capture-manifest.json"
