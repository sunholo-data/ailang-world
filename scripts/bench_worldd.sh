#!/usr/bin/env bash
# Non-vacuous benchmark smoke gate for the world daemon.
set -euo pipefail
cd "$(dirname "$0")/.."

usage() {
  echo "usage: $0 --smoke" >&2
  echo "       $0 --record-pair --variant <dir> --control <dir>" >&2
  echo "       $0 --check-claims" >&2
}

REC_PROBE_TIMEOUT_S=120
REC_UTIL_TIMEOUT_S=20
REC_PREBUILD_TIMEOUT_S=600
REC_LEG_TIMEOUT_S=120

# Mirrored byte-identically from verify_ail.sh:61-74 (V26, P23/P30).
run_bounded() {  # $1=timeout_s  $2=out_file  $3..=cmd ;  exit 124 + named msg on expiry
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
        sys.stderr.write("✗ TIMEOUT after %ds: %s\n" % (t, " ".join(cmd))); sys.exit(124)
PY
}

record_pair() {
  local variant_dir="" control_dir="" result rc
  local record_tmp="${TMPDIR:-/tmp}/bench_worldd.record.$$"

  shift
  while [ "$#" -gt 0 ]; do
    case "$1" in
      --variant)
        [ "$#" -ge 2 ] || { usage; return 2; }
        variant_dir="$2"
        shift 2
        ;;
      --control)
        [ "$#" -ge 2 ] || { usage; return 2; }
        control_dir="$2"
        shift 2
        ;;
      *)
        usage
        return 2
        ;;
    esac
  done
  [ -n "$variant_dir" ] && [ -n "$control_dir" ] || { usage; return 2; }

  local role dir
  for role in variant control; do
    if [ "$role" = variant ]; then dir="$variant_dir"; else dir="$control_dir"; fi
    if [ ! -d "$dir" ]; then
      echo "✗ record-pair FAILED: $role directory does not exist: $dir" >&2
      return 1
    fi
    if run_bounded "$REC_UTIL_TIMEOUT_S" "$record_tmp" git -C "$dir" rev-parse --git-dir; then
      result="$(<"$record_tmp")"
    else
      rc=$?
      echo "✗ record-pair FAILED: $role directory is not a git worktree: $dir" >&2
      return 1
    fi
    if [[ ! "$result" =~ [^[:space:]] ]]; then
      echo "✗ record-pair FAILED: $role directory is not a git worktree: $dir" >&2
      return 1
    fi
  done

  probe_capture() {
    local timeout_s="$1" label="$2" target="$3"; shift 3
    if run_bounded "$timeout_s" "$record_tmp" "$@"; then
      result="$(<"$record_tmp")"
    else
      rc=$?
      echo "✗ probe FAILED: $label" >&2
      return 1
    fi
    if [[ ! "$result" =~ [^[:space:]] ]]; then
      echo "✗ probe FAILED: $label" >&2
      return 1
    fi
    printf -v "$target" '%s' "$result"
  }

  local ncpu hw_model rig_load rig_ps rig_utc ailang_version
  probe_capture "$REC_UTIL_TIMEOUT_S" "sysctl -n hw.ncpu" ncpu sysctl -n hw.ncpu || return 1
  probe_capture "$REC_UTIL_TIMEOUT_S" "sysctl -n hw.model" hw_model sysctl -n hw.model || return 1
  probe_capture "$REC_UTIL_TIMEOUT_S" "sysctl -n vm.loadavg" rig_load sysctl -n vm.loadavg || return 1
  probe_capture "$REC_UTIL_TIMEOUT_S" "ps -Ao pid=,ppid=,pcpu=,comm=" rig_ps ps -Ao pid=,ppid=,pcpu=,comm= || return 1
  probe_capture "$REC_UTIL_TIMEOUT_S" "date -u" rig_utc date -u +%Y-%m-%dT%H:%M:%SZ || return 1
  if [ -z "${AILANG_BIN:-}" ] || [ ! -x "$AILANG_BIN" ]; then
    echo "✗ probe FAILED: AILANG_BIN set and executable" >&2
    return 1
  fi
  probe_capture "$REC_PROBE_TIMEOUT_S" "AILANG_BIN --version" ailang_version "$AILANG_BIN" --version || return 1
  # Schema is ONE line (doc D1). `ailang --version` prints a banner; keep only the
  # version line so a line-oriented parser cannot mistake the banner for extra keys.
  ailang_version="$(printf '%s' "$ailang_version" | head -1)"

  local session_utc pair_nonce
  probe_capture "$REC_UTIL_TIMEOUT_S" "date -u (session stamp)" session_utc date -u +%Y-%m-%dT%H:%M:%SZ || return 1
  probe_capture "$REC_UTIL_TIMEOUT_S" "python3 pair nonce" pair_nonce python3 -c 'import secrets; print(secrets.token_hex(16))' || return 1
  if [[ ! "$pair_nonce" =~ ^[0-9a-f]{32}$ ]]; then
    echo "✗ session nonce FAILED" >&2
    return 1
  fi

  local variant_go control_go
  if run_bounded "$REC_PROBE_TIMEOUT_S" "$record_tmp" go -C "$variant_dir" env GOVERSION GOOS GOARCH; then
    variant_go="$(<"$record_tmp")"
  else
    rc=$?
    echo "✗ tree probe FAILED: go -C $variant_dir env GOVERSION GOOS GOARCH" >&2
    return 1
  fi
  if [[ ! "$variant_go" =~ [^[:space:]] ]]; then
    echo "✗ tree probe FAILED: go -C $variant_dir env GOVERSION GOOS GOARCH" >&2
    return 1
  fi
  if run_bounded "$REC_PROBE_TIMEOUT_S" "$record_tmp" go -C "$control_dir" env GOVERSION GOOS GOARCH; then
    control_go="$(<"$record_tmp")"
  else
    rc=$?
    echo "✗ tree probe FAILED: go -C $control_dir env GOVERSION GOOS GOARCH" >&2
    return 1
  fi
  if [[ ! "$control_go" =~ [^[:space:]] ]]; then
    echo "✗ tree probe FAILED: go -C $control_dir env GOVERSION GOOS GOARCH" >&2
    return 1
  fi

  local variant_commit variant_parent control_commit control_parent git_state
  for role in variant control; do
    if [ "$role" = variant ]; then dir="$variant_dir"; else dir="$control_dir"; fi
    if run_bounded "$REC_UTIL_TIMEOUT_S" "$record_tmp" git -C "$dir" status --porcelain; then
      git_state="$(<"$record_tmp")"
    else
      rc=$?
      echo "✗ git state FAILED: git -C $dir status --porcelain" >&2
      return 1
    fi
    if [ -n "$git_state" ]; then
      echo "✗ record-pair REFUSED: $role tree is dirty: $dir" >&2
      printf '%s\n' "$git_state" >&2
      return 1
    fi
    if run_bounded "$REC_UTIL_TIMEOUT_S" "$record_tmp" git -C "$dir" rev-parse HEAD; then
      result="$(<"$record_tmp")"
    else
      rc=$?
      echo "✗ git state FAILED: git -C $dir rev-parse HEAD" >&2
      return 1
    fi
    printf -v "${role}_commit" '%s' "$result"
    if run_bounded "$REC_UTIL_TIMEOUT_S" "$record_tmp" git -C "$dir" rev-parse 'HEAD^'; then
      result="$(<"$record_tmp")"
    else
      rc=$?
      echo "✗ git state FAILED: git -C $dir rev-parse HEAD^" >&2
      return 1
    fi
    printf -v "${role}_parent" '%s' "$result"
  done
  if [ "$control_commit" != "$variant_parent" ]; then
    echo "✗ record-pair REFUSED: control commit is not the variant parent" >&2
    return 1
  fi

  local pair_id
  if run_bounded "$REC_UTIL_TIMEOUT_S" "$record_tmp" python3 -c 'import hashlib, sys; print(hashlib.sha256(("bench-pair/2\n" + "\n".join(sys.argv[1:]) + "\n").encode()).hexdigest())' "$session_utc" "$variant_commit" "$control_commit" "$pair_nonce"; then
    pair_id="$(<"$record_tmp")"
  else
    rc=$?
    echo "✗ pair ID derivation FAILED" >&2
    return 1
  fi
  if [[ ! "$pair_id" =~ ^[0-9a-f]{64}$ ]]; then
    echo "✗ pair ID derivation FAILED" >&2
    return 1
  fi

  local variant_pair_variant_commit="$variant_commit"
  local variant_pair_control_commit="$control_commit"
  local control_pair_variant_commit="$variant_commit"
  local control_pair_control_commit="$control_commit"
  : "$ncpu" "$hw_model" "$rig_load" "$rig_ps" "$rig_utc" "$ailang_version"
  : "$variant_go" "$control_go" "$variant_parent" "$control_parent" "$pair_id"
  : "$variant_pair_variant_commit" "$variant_pair_control_commit"
  : "$control_pair_variant_commit" "$control_pair_control_commit"

  local session_dir
  if run_bounded "$REC_UTIL_TIMEOUT_S" "$record_tmp" mktemp -d -t bench_worldd.session.XXXXXX; then
    session_dir="$(<"$record_tmp")"
  else
    rc=$?
    echo "✗ session staging FAILED" >&2
    return 1
  fi
  if [[ ! "$session_dir" =~ [^[:space:]] ]] || [ ! -d "$session_dir" ]; then
    echo "✗ session staging FAILED" >&2
    return 1
  fi

  local started_s ended_s binary_file
  local control_prebuild_elapsed_s variant_prebuild_elapsed_s
  local control_binary_sha256 variant_binary_sha256
  for role in control variant; do
    if [ "$role" = variant ]; then dir="$variant_dir"; else dir="$control_dir"; fi
    binary_file="$session_dir/bench.$role"
    probe_capture "$REC_UTIL_TIMEOUT_S" "date +%s (prebuild $role start)" started_s date +%s || return 1
    if run_bounded "$REC_PREBUILD_TIMEOUT_S" "$session_dir/prebuild.$role.out" \
      go -C "$dir" test -c -o "$binary_file" ./host/daemon/; then
      :
    else
      rc=$?
      if [ "$rc" -eq 124 ]; then
        echo "✗ prebuild $role TIMEOUT (>${REC_PREBUILD_TIMEOUT_S}s) — no conditions emitted" >&2
      else
        echo "✗ prebuild $role FAILED — no conditions emitted" >&2
      fi
      printf '%s' "$(<"$session_dir/prebuild.$role.out")" >&2
      return 1
    fi
    probe_capture "$REC_UTIL_TIMEOUT_S" "date +%s (prebuild $role end)" ended_s date +%s || return 1
    printf -v "${role}_prebuild_elapsed_s" '%s' "$((ended_s - started_s))"
    if run_bounded "$REC_UTIL_TIMEOUT_S" "$record_tmp" python3 -c \
      'import hashlib, pathlib, sys; print(hashlib.sha256(pathlib.Path(sys.argv[1]).read_bytes()).hexdigest())' \
      "$binary_file"; then
      result="$(<"$record_tmp")"
    else
      rc=$?
      echo "✗ binary hash FAILED: $role" >&2
      return 1
    fi
    if [[ ! "$result" =~ ^[0-9a-f]{64}$ ]]; then
      echo "✗ binary hash FAILED: $role" >&2
      return 1
    fi
    printf -v "${role}_binary_sha256" '%s' "$result"
  done

  local -a leg_roles=(control variant variant control)
  local leg_num leg_start_utc leg_end_utc leg_load leg_competing leg_started_s leg_ended_s
  local leg_output leg_hash ps_file leg_command
  for leg_num in 1 2 3 4; do
    role="${leg_roles[$((leg_num - 1))]}"
    binary_file="$session_dir/bench.$role"
    leg_output="$session_dir/leg$leg_num.out"
    ps_file="$session_dir/leg$leg_num.ps"
    probe_capture "$REC_UTIL_TIMEOUT_S" "date -u (leg $leg_num start)" leg_start_utc date -u +%Y-%m-%dT%H:%M:%SZ || return 1
    probe_capture "$REC_UTIL_TIMEOUT_S" "date +%s (leg $leg_num start)" leg_started_s date +%s || return 1
    probe_capture "$REC_UTIL_TIMEOUT_S" "sysctl -n vm.loadavg (leg $leg_num)" leg_load sysctl -n vm.loadavg || return 1
    if run_bounded "$REC_UTIL_TIMEOUT_S" "$ps_file" ps -Ao pid=,ppid=,pcpu=,comm=; then
      :
    else
      rc=$?
      echo "✗ probe FAILED: ps -Ao pid=,ppid=,pcpu=,comm= (leg $leg_num)" >&2
      return 1
    fi
    if run_bounded "$REC_UTIL_TIMEOUT_S" "$record_tmp" python3 -c '
import pathlib, sys
rows = []
for line in pathlib.Path(sys.argv[1]).read_text(errors="replace").splitlines():
    fields = line.strip().split(None, 3)
    if len(fields) != 4:
        continue
    try:
        pid, ppid, pcpu = int(fields[0]), int(fields[1]), float(fields[2])
    except ValueError:
        continue
    rows.append((pid, ppid, pcpu, fields[3]))
parents = {pid: comm for pid, _, _, comm in rows}
hot = sorted((row for row in rows if row[2] >= 25.0), key=lambda row: row[2], reverse=True)[:8]
if not hot:
    print("none>=25%")
else:
    for pid, ppid, pcpu, comm in hot:
        print(f"pcpu={pcpu:.1f} pid={pid} comm={comm} parent={parents.get(ppid, chr(63))}")
' "$ps_file"
    then
      leg_competing="$(<"$record_tmp")"
    else
      rc=$?
      echo "✗ competing-process capture FAILED: leg $leg_num" >&2
      return 1
    fi
    if [[ ! "$leg_competing" =~ [^[:space:]] ]]; then
      echo "✗ competing-process capture FAILED: leg $leg_num" >&2
      return 1
    fi

    leg_command="$binary_file -test.bench . -test.benchtime 200x -test.run '^$'"
    if (cd "$session_dir" && run_bounded "$REC_LEG_TIMEOUT_S" "$leg_output" \
      "$binary_file" -test.bench . -test.benchtime 200x -test.run '^$'); then
      :
    else
      rc=$?
      if [ "$rc" -eq 124 ]; then
        echo "✗ bench leg $leg_num/4 TIMEOUT (>${REC_LEG_TIMEOUT_S}s) — no conditions emitted" >&2
      else
        echo "✗ bench leg $leg_num/4 FAILED — no conditions emitted" >&2
      fi
      printf '%s' "$(<"$leg_output")" >&2
      return 1
    fi
    probe_capture "$REC_UTIL_TIMEOUT_S" "date -u (leg $leg_num end)" leg_end_utc date -u +%Y-%m-%dT%H:%M:%SZ || return 1
    probe_capture "$REC_UTIL_TIMEOUT_S" "date +%s (leg $leg_num end)" leg_ended_s date +%s || return 1
    if run_bounded "$REC_UTIL_TIMEOUT_S" "$record_tmp" python3 -c \
      'import hashlib, pathlib, sys; data=pathlib.Path(sys.argv[1]).read_bytes(); data=data.rstrip(b"\n")+b"\n"; print(hashlib.sha256(data).hexdigest())' \
      "$leg_output"; then
      leg_hash="$(<"$record_tmp")"
    else
      rc=$?
      echo "✗ leg output hash FAILED: leg $leg_num" >&2
      return 1
    fi
    if [[ ! "$leg_hash" =~ ^[0-9a-f]{64}$ ]]; then
      echo "✗ leg output hash FAILED: leg $leg_num" >&2
      return 1
    fi
    printf -v "leg${leg_num}_start_utc" '%s' "$leg_start_utc"
    printf -v "leg${leg_num}_end_utc" '%s' "$leg_end_utc"
    printf -v "leg${leg_num}_elapsed_s" '%s' "$((leg_ended_s - leg_started_s))"
    leg_load="${leg_load#\{}"; leg_load="${leg_load%\}}"
    leg_load="${leg_load# }"; leg_load="${leg_load% }"
    printf -v "leg${leg_num}_load" '%s' "$leg_load"
    printf -v "leg${leg_num}_competing" '%s' "$leg_competing"
    printf -v "leg${leg_num}_output_sha256" '%s' "$leg_hash"
  done

  # Assembly begins only after leg 4 and all hashes have succeeded.
  local emission_stage="$session_dir/emission.stage" emission_file conditions_file conditions_hash
  : > "$emission_stage"
  emit_role() {
    local emit_role="$1" first="$2" second="$3" go_values commit parent prebuild binary_hash
    local first_competing second_competing
    if [ "$emit_role" = variant ]; then
      go_values="$variant_go"; commit="$variant_commit"; parent="$variant_parent"
      prebuild="$variant_prebuild_elapsed_s"; binary_hash="$variant_binary_sha256"
    else
      go_values="$control_go"; commit="$control_commit"; parent="$control_parent"
      prebuild="$control_prebuild_elapsed_s"; binary_hash="$control_binary_sha256"
    fi
    local goversion goos goarch
    goversion="${go_values%%$'\n'*}"
    go_values="${go_values#*$'\n'}"; goos="${go_values%%$'\n'*}"; goarch="${go_values##*$'\n'}"
    eval "first_competing=\$leg${first}_competing"
    eval "second_competing=\$leg${second}_competing"
    first_competing="${first_competing//$'\n'/$'\nleg1_competing: '}"
    second_competing="${second_competing//$'\n'/$'\nleg2_competing: '}"
    conditions_file="$session_dir/conditions.$emit_role"
    : > "$conditions_file"
    printf '%s\n' \
      "schema: bench-conditions/2" "role: $emit_role" "pair_id: $pair_id" \
      "pair_nonce: $pair_nonce" "session_utc: $session_utc" \
      "pair_variant_commit: $variant_commit" "pair_control_commit: $control_commit" \
      "commit: $commit" "parent: $parent" "tree: clean" "goversion: $goversion" \
      "goos_goarch: $goos/$goarch" "ncpu: $ncpu" "hw_model: $hw_model" \
      "ailang_pin: $ailang_version via \$AILANG_BIN" "prebuild_elapsed_s: $prebuild" \
      "binary_sha256: $binary_hash" \
      "invocation: $session_dir/bench.$emit_role -test.bench . -test.benchtime 200x -test.run '^$'" \
      "leg_order: control,variant,variant,control" \
      "leg1_seq: $first/4" "leg1_start_utc: $(eval printf '%s' \"\$leg${first}_start_utc\")" \
      "leg1_end_utc: $(eval printf '%s' \"\$leg${first}_end_utc\")" \
      "leg1_elapsed_s: $(eval printf '%s' \"\$leg${first}_elapsed_s\")" \
      "leg1_load: $(eval printf '%s' \"\$leg${first}_load\")" \
      "leg1_competing: $first_competing" \
      "leg1_output_sha256: $(eval printf '%s' \"\$leg${first}_output_sha256\")" \
      "leg2_seq: $second/4" "leg2_start_utc: $(eval printf '%s' \"\$leg${second}_start_utc\")" \
      "leg2_end_utc: $(eval printf '%s' \"\$leg${second}_end_utc\")" \
      "leg2_elapsed_s: $(eval printf '%s' \"\$leg${second}_elapsed_s\")" \
      "leg2_load: $(eval printf '%s' \"\$leg${second}_load\")" \
      "leg2_competing: $second_competing" \
      "leg2_output_sha256: $(eval printf '%s' \"\$leg${second}_output_sha256\")" >> "$conditions_file"
    if run_bounded "$REC_UTIL_TIMEOUT_S" "$record_tmp" python3 -c \
      'import hashlib, pathlib, sys; print(hashlib.sha256(pathlib.Path(sys.argv[1]).read_bytes()).hexdigest())' \
      "$conditions_file"; then
      conditions_hash="$(<"$record_tmp")"
    else
      rc=$?
      echo "✗ conditions hash FAILED: $emit_role" >&2
      return 1
    fi
    printf '%s\n' '```bench-conditions' "$(<"$conditions_file")" \
      "conditions_sha256: $conditions_hash" '```' '```text' >> "$emission_stage"
    printf '%s\n' "$(<"$session_dir/leg$first.out")" '```' '```text' \
      "$(<"$session_dir/leg$second.out")" '```' >> "$emission_stage"
  }
  emit_role variant 2 3 || return 1
  emit_role control 1 4 || return 1

  if run_bounded "$REC_UTIL_TIMEOUT_S" "$record_tmp" mktemp -t bench_worldd.pair.XXXXXX; then
    emission_file="$(<"$record_tmp")"
  else
    rc=$?
    echo "✗ emission file creation FAILED" >&2
    return 1
  fi
  if run_bounded "$REC_UTIL_TIMEOUT_S" "$record_tmp" cp "$emission_stage" "$emission_file"; then
    :
  else
    rc=$?
    echo "✗ emission publish FAILED" >&2
    return 1
  fi
  printf '%s\n' "$(<"$emission_file")"
  printf 'paste-ready pair written to: %s\n' "$emission_file"
}

if [ "$#" -ge 1 ] && [ "$1" = "--record-pair" ]; then
  record_pair "$@"
  exit $?
fi

if [ "$#" -ne 1 ] || [ "$1" != "--smoke" ]; then
  usage
  exit 2
fi

expected=(
  BenchmarkStoreCommit
  BenchmarkJournalAppend
  BenchmarkCommitWithReceipt
  BenchmarkHeadRead
  BenchmarkHealth
  BenchmarkRESTCommit
  BenchmarkLogRange/limit_100
  BenchmarkLogRange/limit_500
  BenchmarkBrokerDecide
  BenchmarkBrokerFSRead
)

out="$(mktemp -t bench_worldd.XXXXXX)"
trap 'rm -f "$out"' EXIT

echo "── worldd benchmark smoke: go test -bench . -benchtime 1x -run '^$' ./host/daemon/"
if ! go test -bench . -benchtime 1x -run '^$' ./host/daemon/ 2>&1 | tee "$out"; then
  echo "✗ worldd benchmark smoke FAILED: underlying go test failed" >&2
  exit 1
fi

missing=()
for name in "${expected[@]}"; do
  if ! grep -Eq "^${name}(-[0-9]+)?[[:space:]]" "$out"; then
    missing+=("$name")
  fi
done

if [ "${#missing[@]}" -ne 0 ]; then
  echo "✗ worldd benchmark smoke FAILED: missing expected benchmark(s): ${missing[*]}" >&2
  exit 1
fi

echo "✓ worldd benchmark smoke PASSED: ${expected[*]}"
