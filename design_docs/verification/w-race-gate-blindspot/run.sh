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

# Bounded fail-loud parser for the data-only fixture. Accepts ONLY the three
# names KNOWN_BAD, KNOWN_GOOD, PINNED; values must be inert toolchain tokens
# (character class [A-Za-z0-9._+:/ -]) plus spaces. Refuses loudly on anything
# else. This is a parser, not a source: it never evaluates the file as code.
conf="$(dirname "$0")/toolchain_pins.conf"
seen_bad=0; seen_good=0; seen_pinned=0
while IFS= read -r line || [ -n "$line" ]; do
	case "$line" in
		''|\#*) continue ;;
	esac
	# Both regexes must live in variables: their inline forms are syntax errors
	# under bash 3.2, which this repository's launchd lane pins.
	RE='^([A-Z_]+)="([^"]*)"([[:space:]]+#.*)?$'
	BADCHARS='[^A-Za-z0-9._+:/ -]'
	if [[ "$line" =~ $RE ]]; then
		name="${BASH_REMATCH[1]}"; value="${BASH_REMATCH[2]}"
		case "$name" in
			KNOWN_BAD)  seen_bad=$((seen_bad+1)) ;;
			KNOWN_GOOD) seen_good=$((seen_good+1)) ;;
			PINNED)     seen_pinned=$((seen_pinned+1)) ;;
			*) echo "toolchain_pins.conf: unknown name '$name' (only KNOWN_BAD, KNOWN_GOOD, PINNED allowed)" >&2; exit 1 ;;
		esac
		if [[ "$value" =~ $BADCHARS ]]; then
			echo "toolchain_pins.conf: value for '$name' contains a disallowed character (only toolchain tokens and spaces allowed)" >&2; exit 1
		fi
		printf -v "$name" '%s' "$value"
	else
		echo "toolchain_pins.conf: malformed line (must be NAME=\"value\" at column 0): $line" >&2; exit 1
	fi
done < "$conf"
if [ "$seen_bad" -ne 1 ] || [ "$seen_good" -ne 1 ] || [ "$seen_pinned" -ne 1 ]; then
	echo "toolchain_pins.conf: expected exactly one each of KNOWN_BAD, KNOWN_GOOD, PINNED (got $seen_bad/$seen_good/$seen_pinned)" >&2; exit 1
fi
# toolchain_pins.conf parser ends here.

cd "$(dirname "$0")/repro" || exit 1

saw_bad=0
saw_good=0
saw_pinned_ok=0
ran=0
ran_bad=0       # KNOWN-BAD probes that actually built and ran (row 44)
bad_expected=0  # KNOWN-BAD probes CONFIGURED, counted at probe entry — tracks the list

# Platform probe (row 44; quorum round-1 R3): read the kernel, never `go env` —
# `go env GOOS` honours a $GOOS environment variable (measured: design doc P16), so an
# overridable variable must not be the polarity's sole input. `uname` has no such
# override channel. Normalized to Go's naming so the deny-set comparisons stay in one
# vocabulary; unknown values fall into the refuse arm below.
case "$(uname -s)" in
	Darwin) host_os=darwin ;;
	Linux)  host_os=linux ;;
	*)      host_os=unknown ;;
esac
case "$(uname -m)" in
	arm64|aarch64) host_arch=arm64 ;;
	x86_64|amd64)  host_arch=amd64 ;;
	*)             host_arch=unknown ;;
esac
host_pair="$host_os/$host_arch"

# Platform expectation (row 44; quorum round-1 R1): the verified contract set is
# EXACTLY two pairs (design doc P18). darwin/arm64: the deny-list's measured set — a
# KNOWN-BAD probe MUST report BUG or this script can no longer see the defect.
# linux/amd64: measured NOT affected (iteration-46 AC6) — no KNOWN-BAD probe may report
# BUG on the one platform CI builds on. Any other host has no verified contract and
# refuses; extend the measurement set before trusting it, do not inherit a polarity.
case "$host_pair" in
	darwin/arm64) expect_defect=1 ;;
	linux/amd64) expect_defect=0 ;;
	*) echo "INSTRUMENT FAILURE: no verified platform contract for $host_pair"; exit 1 ;;
esac

probe() { # $1=toolchain  $2=expectation label
	local tc="$1" expect="$2" out rc bin
	[ "$expect" = "BAD" ] && bad_expected=$((bad_expected + 1))
	bin="$(mktemp "${TMPDIR:-/tmp}/w40repro.XXXXXX")" || return 1
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
	[ "$expect" = "BAD" ] && ran_bad=$((ran_bad + 1))
	printf '  %-10s expect=%-5s got: %s (rc=%d)\n' "$tc" "$expect" "${out:-<NO OUTPUT>}" "$rc"
	# An empty reading is an instrument failure, never a pass.
	if [ -z "$out" ]; then
		echo "  !! $tc produced NO OUTPUT — treating as instrument failure, not a pass"
		return 1
	fi
	case "$out" in
	OK*) [ "$expect" = GOOD ] && saw_good=1; [ "$tc" = "$PINNED" ] && saw_pinned_ok=1 ;;
	BUG*) [ "$expect" = BAD ] && saw_bad=1 ;;
	esac
	return 0
}

echo "== go1.26 local-array-literal miscompilation reproduction =="
echo "host: $host_pair   default toolchain: $(go version | awk '{print $3}')"
echo
echo "expected BUG (affected):"
for tc in $KNOWN_BAD; do probe "$tc" BAD || exit 1; done
echo "expected OK (unaffected):"
for tc in $KNOWN_GOOD; do probe "$tc" GOOD || exit 1; done

echo
echo "-- optimization-level control on the default toolchain --"
bin="$(mktemp "${TMPDIR:-/tmp}/w40reproN.XXXXXX")"
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
if [ "$bad_expected" -eq 0 ] || [ "$ran_bad" -ne "$bad_expected" ]; then
	echo "INSTRUMENT FAILURE: $ran_bad of $bad_expected KNOWN-BAD probes completed —"
	echo "a partial (or hollowed-out) negative arm cannot certify $host_pair. Every"
	echo "configured KNOWN-BAD entry must run; a SKIP on this arm is a refusal, not a pass."
	exit 1
fi
if [ "$expect_defect" -eq 1 ] && [ "$saw_bad" -eq 0 ]; then
	echo "INSTRUMENT FAILURE (or GOOD NEWS): no known-affected toolchain reproduced the"
	echo "defect. Either the toolchains were unavailable, or upstream fixed it — in which"
	echo "case re-derive the pin decision in design_docs/ rather than trusting this pass."
	exit 1
fi
if [ "$expect_defect" -eq 0 ] && [ "$saw_bad" -ne 0 ]; then
	echo "INSTRUMENT FAILURE (PLATFORM ALARM): a KNOWN-BAD toolchain reported BUG on"
	echo "$host_pair — measured NOT affected (iteration-46 AC6). The defect escaped its"
	echo "measured set: treat as a toolchain incident and re-derive the pin."
	exit 1
fi
if [ "$saw_good" -eq 0 ]; then
	echo "INSTRUMENT FAILURE: no known-good toolchain produced OK, so a BUG line above"
	echo "cannot be attributed to the compiler. Fix the harness before citing any result."
	exit 1
fi
if [ "$saw_pinned_ok" -eq 0 ]; then
	echo "INSTRUMENT FAILURE: the PINNED toolchain ($PINNED) never reported OK — it was"
	echo "SKIPPED (unfetchable) or is absent from the probe lists. A success banner that"
	echo "never probed the pin is a false clean; refusing to print it."
	exit 1
fi
if [ "$expect_defect" -eq 1 ]; then
	echo "RESULT: reproduction confirmed, and both controls fired."
	echo "  known-affected toolchain reported BUG : yes"
	echo "  known-good toolchain reported OK      : yes"
	echo "  pinned toolchain ($PINNED) reported OK  : yes"
else
	echo "RESULT: $host_pair clean — no KNOWN-BAD toolchain reproduced here, matching"
	echo "the iteration-46 AC6 measurement; all $bad_expected known-bad and $ran total"
	echo "probes ran, and the known-good and pinned ($PINNED) toolchains both reported OK."
fi
