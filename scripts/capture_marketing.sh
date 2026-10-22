#!/usr/bin/env bash
# Capture entry point. M7a ships supervision; M7c supplies the six-shot pass.
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

The six-shot pass is not implemented yet (M7c); capture exits nonzero.
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

echo 'six-shot pass is not implemented yet (M7c); capture refused' >&2
exit 1
