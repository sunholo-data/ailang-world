#!/bin/bash
# Source the environment block first, then run this block in the same shell.
set -euo pipefail
P=host/daemon/episode_finish.go
T=host/daemon/episode_finish_approve_guard_test.go
N='^TestEpisodeFinishApproveRefusesRecord$'
sha() { shasum -a 256 "$1" | awk '{print $1}'; }
[ "$(sha "$P")" = "4e320c537e689620480af2bdff9bdcd9d4d6febba43b4a5dd5b3bb09971973bf" ]
[ "$(sha "$T")" = "c9e279de7c7ce87edf5fbc0e43e690253616d320b5585625062cda3d29cdec57" ]
cp "$P" "$E/production.base"
cp "$T" "$E/final-test.base"
restore() {
  cp "$E/production.base" "$P"
  cp "$E/final-test.base" "$T"
  [ "$(sha "$P")" = "4e320c537e689620480af2bdff9bdcd9d4d6febba43b4a5dd5b3bb09971973bf" ]
  [ "$(sha "$T")" = "c9e279de7c7ce87edf5fbc0e43e690253616d320b5585625062cda3d29cdec57" ]
  shasum -a 256 "$P" "$T"
}
trap restore EXIT
trap 'exit 130' INT
trap 'exit 143' TERM
# run LABEL EXPECTED_RC REQUIRED_LITERAL COMMAND [ARG ...]
# "$@" is positional arguments, never a possibly-empty Bash array.
run() {
  local label=$1 want=$2 marker=$3 rc=0
  shift 3
  "$@" > "$E/$label.txt" 2>&1 || rc=$?
  cat "$E/$label.txt"
  printf '%s rc=%s expected=%s\n' "$label" "$rc" "$want"
  # Inspect ALL other sandbox/path errors manually under the rule below.
  if grep -Eiq '(bind:.*(not permitted|permission denied)|sandbox.*deny|operation not permitted)' "$E/$label.txt"; then
    printf '%s\n' 'UNINFORMATIVE UNDER SANDBOX'
    return 1
  fi
  [ "$rc" -eq "$want" ]
  if [ -n "$marker" ]; then grep -Fq -- "$marker" "$E/$label.txt"; fi
}
probe() {
  shasum -a 256 "$P" "$T" > "$E/$1.hashes.txt"
  run "$1-vet" 0 '' go vet ./host/daemon/
  run "$1" "$2" "$3" go test -count=1 -v -run "$N" ./host/daemon/
}
m1() {
  sed 's/ || input.Record != nil {/ {/' "$P" > "$E/edited"
  cp "$E/edited" "$P"
}
m3c() {
  sed 's/if err != nil || scope.ExpiresAt <= now || scope.ExpiresAt > now+3600 || input.Record != nil {/if true || err != nil || scope.ExpiresAt <= now || scope.ExpiresAt > now+3600 || input.Record != nil {/' "$P" > "$E/edited"
  cp "$E/edited" "$P"
}
m3a() {
  m1
  sed '/^[[:space:]]*request, decision, err := broker.RecordOperatorApproval(ctx, d.store, input.Scope, "operator:episode-finish", "operator:terminal", now)$/a\
		if input.Record != nil {\
			refuse("operator route: invalid or expired finish scope")\
			return\
		}
' "$P" > "$E/edited"
  cp "$E/edited" "$P"
}
probe clean 0 '--- PASS: TestEpisodeFinishApproveRefusesRecord'
m1
probe M1 1 'approve with record: 200'
run M1-existing 0 'ok' go test -count=1 -run Finish -skip "$N" ./host/daemon/
restore
m3a
probe M3a 1 'refused approve recorded an approval: head=true'
restore
m3c
probe M3c 1 'control approve without record: 403'
restore
# M2: remove final-file lines 57..66 (control), and line 4 (json import).
sed '4d;57,66d' "$T" > "$E/edited"
cp "$E/edited" "$T"
probe M2-clean 0 '--- PASS: TestEpisodeFinishApproveRefusesRecord'
m3c
probe M2+M3c 0 '--- PASS: TestEpisodeFinishApproveRefusesRecord'
restore
# M2b: invert only the control's expected status.
sed 's/if w.Code != 200 || json.Unmarshal/if w.Code != 403 || json.Unmarshal/' "$T" > "$E/edited"
cp "$E/edited" "$T"
probe M2b 1 'control approve without record: 200'
restore
# M4: delete ONLY lines 53..55 (post-refusal no-write check).
# Keep lines 44..46 (the new pristine-store precondition).
sed '53,55d' "$T" > "$E/edited"
cp "$E/edited" "$T"
m3a
probe M4+M3a 0 '--- PASS: TestEpisodeFinishApproveRefusesRecord'
restore
trap - EXIT INT TERM
