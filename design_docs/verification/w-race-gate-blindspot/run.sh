#!/usr/bin/env bash
# Re-run the go1.26 array-literal miscompilation reproduction, first-party.
#
# This script is an INSTRUMENT, so it carries its own known-positive controls and
# fails LOUDLY rather than silently passing (charter Gate-2 rule 3a: "a probe that
# came back empty is a claim, not a fact"; an empty result set must t.Fatal, never
# pass). Specifically it refuses to report success unless BOTH controls fire:
#
#   * at least one KNOWN-BAD toolchain actually reported BUG   -> proves this script
#     can still SEE the defect (if go ever fixes it, that control fails and you are
#     told, instead of the script quietly reporting "all clear")
#   * at least one KNOWN-GOOD toolchain actually reported OK   -> proves the program
#     can pass at all, i.e. a BUG line means the compiler and not the reproducer
#
# Usage: ./run.sh          (from anywhere; paths resolve relative to this script)
#
# Toolchains are fetched on demand by the go command; with no network, a missing
# toolchain is reported SKIPPED and does not count toward either control.

set -uo pipefail

cd "$(dirname "$0")/repro" || exit 1

KNOWN_BAD="go1.26.0 go1.26.3 go1.26.4 go1.26.5"
KNOWN_GOOD="go1.25.6 go1.24.9"

saw_bad=0
saw_good=0
ran=0

probe() { # $1=toolchain  $2=expectation label
	local tc="$1" expect="$2" out rc bin
	bin="$(mktemp -t w40repro)" || return 1
	if ! GOTOOLCHAIN="$tc" go build -o "$bin" . 2>/tmp/w40_build_err.txt; then
		printf '  %-10s SKIPPED (toolchain unavailable: %s)\n' \
			"$tc" "$(tr -d '\n' </tmp/w40_build_err.txt | cut -c1-70)"
		rm -f "$bin"
		return 0
	fi
	out="$("$bin" 2>&1)"
	rc=$?
	rm -f "$bin"
	ran=$((ran + 1))
	printf '  %-10s expect=%-5s got: %s (rc=%d)\n' "$tc" "$expect" "${out:-<NO OUTPUT>}" "$rc"
	# An empty reading is an instrument failure, never a pass.
	if [ -z "$out" ]; then
		echo "  !! $tc produced NO OUTPUT — treating as instrument failure, not a pass"
		return 1
	fi
	case "$out" in
	OK*) [ "$expect" = GOOD ] && saw_good=1 ;;
	BUG*) [ "$expect" = BAD ] && saw_bad=1 ;;
	esac
	return 0
}

echo "== go1.26 local-array-literal miscompilation reproduction =="
echo "host: $(go env GOOS)/$(go env GOARCH)   default toolchain: $(go version | awk '{print $3}')"
echo
echo "expected BUG (affected):"
for tc in $KNOWN_BAD; do probe "$tc" BAD || exit 1; done
echo "expected OK (unaffected):"
for tc in $KNOWN_GOOD; do probe "$tc" GOOD || exit 1; done

echo
echo "-- optimization-level control on the default toolchain --"
bin="$(mktemp -t w40reproN)"
if go build -gcflags='all=-N' -o "$bin" .; then
	printf '  %-22s got: %s\n' "-gcflags=all=-N" "$("$bin" 2>&1)"
fi
if go build -gcflags='all=-l' -o "$bin" .; then
	printf '  %-22s got: %s\n' "-gcflags=all=-l" "$("$bin" 2>&1)"
fi
rm -f "$bin"

echo
if [ "$ran" -eq 0 ]; then
	echo "INSTRUMENT FAILURE: no toolchain ran at all — this is NOT a clean result."
	exit 1
fi
if [ "$saw_bad" -eq 0 ]; then
	echo "INSTRUMENT FAILURE (or GOOD NEWS): no known-affected toolchain reproduced the"
	echo "defect. Either the toolchains were unavailable, or upstream fixed it — in which"
	echo "case re-derive the pin decision in design_docs/ rather than trusting this pass."
	exit 1
fi
if [ "$saw_good" -eq 0 ]; then
	echo "INSTRUMENT FAILURE: no known-good toolchain produced OK, so a BUG line above"
	echo "cannot be attributed to the compiler. Fix the harness before citing any result."
	exit 1
fi
echo "RESULT: reproduction confirmed, and both controls fired."
echo "  known-affected toolchain reported BUG : yes"
echo "  known-good toolchain reported OK      : yes"
