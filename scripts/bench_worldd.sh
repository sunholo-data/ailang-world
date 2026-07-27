#!/usr/bin/env bash
# Non-vacuous benchmark smoke gate for the world daemon.
set -euo pipefail
cd "$(dirname "$0")/.."

usage() {
  echo "usage: $0 --smoke" >&2
}

if [ "$#" -ne 1 ] || [ "$1" != "--smoke" ]; then
  usage
  exit 2
fi

expected=(
  BenchmarkStoreCommit
  BenchmarkHeadRead
  BenchmarkHealth
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
