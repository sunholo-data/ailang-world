#!/bin/bash
# Extracted plan blocks run in this one macOS Bash process.
# Foreground tool sessions will be polled to completion, bounded to 120 polls
# of at most 15 seconds each (30 minutes); no background shell processes.
set -euo pipefail
evidence="$PWD/verification/world-iter335/executor"
. "$evidence/environment.sh" > "$evidence/environment.txt" 2>&1
printf 'environment rc=0\n'
{
  /bin/bash --version
  go version
  go env GOTOOLCHAIN GOCACHE GOMODCACHE GOPATH TMPDIR GOTMPDIR GOPROXY GOSUMDB GOFLAGS
  shasum -a 256 host/daemon/episode_finish_approve_guard_test.go
  wc -l host/daemon/episode_finish_approve_guard_test.go
  gofmt -l host/daemon/episode_finish_approve_guard_test.go
} > "$E/instruments.txt" 2>&1
. "$E/mutation-drills.sh"
printf 'mutation-block rc=0\n'
# Record each exact clean-block command and its actual status independently,
# as the plan requires, even if an earlier clean gate is red.
set +e
gate_number=0
while IFS= read -r gate_command; do
  case "$gate_command" in ''|'#'*) continue ;; esac
  gate_number=$((gate_number + 1))
  printf -v gate_label '%02d' "$gate_number"
  eval "$gate_command" > "$E/clean-$gate_label.txt" 2>&1
  gate_rc=$?
  printf '%s\t%s\t%s\n' "$gate_label" "$gate_rc" "$gate_command" >> "$E/clean-exit-codes.tsv"
  cat "$E/clean-$gate_label.txt"
  printf 'clean-%s rc=%s command=%s\n' "$gate_label" "$gate_rc" "$gate_command"
done < "$E/clean-gates.sh"
