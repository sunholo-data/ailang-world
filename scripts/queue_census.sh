#!/usr/bin/env bash
# queue_census.sh — anchored leading-tag classifier for the mission queue.
#
# Reads the mission charter's queue and classifies each row by the FIRST token after
# `^<n>. `, tolerant only of a leading run of `[` and `*`. It never reads the row body.
#
# The scan anchors on the FIRST line matching `^## Queue` (the prototype the 41/89/4/44
# reading was derived from takes the first match; AC-M2-3's `^## Queue` count == 1 keeps
# the choice unobservable) and runs to end-of-file — it does NOT stop at a `## ` line,
# because a row body contains one (charter line 5512 sits inside row 89).
#
# This is a bash-3.2 script (the rig's /bin/bash is 3.2.57). It does NOT `cd` to the repo
# root (the precedent gate1_range_check.sh states the same rule): --doc defaults to
# design_docs/world-mission.md relative to cwd. Enumeration, classification, dup/gap
# detection and sorting live in ONE portable awk program (no gensub, no asort, no {n,m},
# no \b); the bash wrapper owns argument parsing, floors F0/F5/F6 and exit codes.
set -uo pipefail

DOC="design_docs/world-mission.md"
CONTROL_CLOSED=""
CONTROL_OPEN=""
GIT_BIN="git"

usage() {
  cat <<'EOF'
Usage: queue_census.sh [options]

Classify every queue row after the FIRST `^## Queue` heading by its leading tag and print
a per-row table plus a census total stamped with the charter's blob OID.

Options:
  --doc <path>          mission doc (default: design_docs/world-mission.md, relative to cwd)
  --control-closed <n>  a row number that MUST classify closed (mandatory)
  --control-open <n>    a row number that MUST classify open (mandatory)
  --git-bin <path>      git binary for `hash-object` (default: git)
  --help                show this help

Exit codes:
  0  census printed, both controls ok
  2  any floor (F0 unreadable doc, F1 no heading, F2 zero rows, F3 duplicate, F4 gap,
     F5 controls required, F6 control mismatch)
EOF
}

while [ "$#" -gt 0 ]; do
  case "$1" in
    --doc) DOC="$2"; shift 2 ;;
    --control-closed) CONTROL_CLOSED="$2"; shift 2 ;;
    --control-open) CONTROL_OPEN="$2"; shift 2 ;;
    --git-bin) GIT_BIN="$2"; shift 2 ;;
    --help) usage; exit 0 ;;
    *) echo "✗ unknown argument: $1" >&2; usage >&2; exit 2 ;;
  esac
done

# F0 — mission doc missing/unreadable
if [ ! -r "$DOC" ]; then
  echo "✗ mission doc unreadable: $DOC" >&2
  exit 2
fi

# ONE portable awk program: enumeration, classification, dup/gap detection.
AWK_PROG='
/^## Queue/ && !h {
  h=NR
  print "HEADING " NR
  next
}
h && /^[0-9]+\. / {
  n=$0; sub(/^[0-9]+\. /,"",n)
  sub(/^[\[\*]+/,"",n)
  num=$0; sub(/\..*$/,"",num)
  tag="-"; cls="UNTAGGED"
  if (n ~ /^RULED OUT([^A-Za-z]|$)/)      { tag="RULED OUT"; cls="closed" }
  else if (n ~ /^LANDED([^A-Za-z]|$)/)    { tag="LANDED";    cls="closed" }
  else if (n ~ /^ROUTED([^A-Za-z]|$)/)    { tag="ROUTED";    cls="closed" }
  else if (n ~ /^PARKED([^A-Za-z]|$)/)    { tag="PARKED";    cls="open" }
  else if (n ~ /^IN-SPRINT([^A-Za-z]|$)/) { tag="IN-SPRINT"; cls="open" }
  else if (n ~ /^NEXT([^A-Za-z]|$)/)      { tag="NEXT";      cls="open" }
  if (num in firstline) {
    print "DUP " num " " firstline[num] " " NR
  } else {
    firstline[num]=NR
  }
  nums[num]=1
  if (num+0 > max) max=num+0
  print "ROW " num "|" tag "|" cls "|" NR
  next
}
h && /^\*\*\[/ {
  print "PREAMBLE " NR
}
END {
  if (h) {
    print "MAX " max
    missing=""
    for (i=1; i<=max; i++) {
      if (!(i in nums)) {
        if (missing=="") missing=i
        else missing=missing " " i
      }
    }
    if (missing!="") print "GAP " missing
  }
}
'

out="$(awk "$AWK_PROG" "$DOC")"
arc=$?
if [ "$arc" -ne 0 ]; then
  echo "✗ awk failed (rc=$arc)" >&2
  exit 2
fi

# Parse the awk output.
heading_line=""
rows=""
preamble_count=0
preamble_line=""
dup_lines=""
max=""
gap=""

while IFS= read -r line; do
  case "$line" in
    HEADING\ *) heading_line="${line#HEADING }" ;;
    ROW\ *) rows="$rows
${line#ROW }" ;;
    PREAMBLE\ *) preamble_count=$((preamble_count+1)); preamble_line="${line#PREAMBLE }" ;;
    DUP\ *) dup_lines="$dup_lines
${line#DUP }" ;;
    MAX\ *) max="${line#MAX }" ;;
    GAP\ *) gap="${line#GAP }" ;;
  esac
done <<< "$out"

# F1 — no "## Queue" heading
if [ -z "$heading_line" ]; then
  echo "✗ no \"## Queue\" heading in $DOC" >&2
  exit 2
fi

# F2 — zero rows enumerated after the heading
if [ -z "$rows" ]; then
  echo "✗ zero queue rows enumerated after line $heading_line" >&2
  exit 2
fi

# F3 — duplicate row number
if [ -n "$dup_lines" ]; then
  while IFS= read -r d; do
    [ -z "$d" ] && continue
    set -- $d
    echo "✗ duplicate row number: $1 (lines $2, $3)" >&2
  done <<< "$dup_lines"
  exit 2
fi

# F4 — numbering gap
if [ -n "$gap" ]; then
  echo "✗ numbering gap: expected 1..$max, missing $gap" >&2
  exit 2
fi

# F5 — controls required
if [ -z "$CONTROL_CLOSED" ] || [ -z "$CONTROL_OPEN" ]; then
  echo "✗ controls required: --control-closed <n> --control-open <n>" >&2
  exit 2
fi

# Print the header line.
echo "queue_census: doc=$DOC heading_line=$heading_line"

# Print the table, sorted by row number regardless of file order.
printf '%s\n' "$rows" | /usr/bin/sort -t'|' -k1,1n | while IFS= read -r r; do
  [ -z "$r" ] && continue
  IFS='|' read -r num tag cls line <<< "$r"
  echo "row=$num tag=$tag class=$cls line=$line"
done

# Print the preamble line (always, so the denominator decision is visible).
if [ "$preamble_count" -gt 0 ]; then
  echo "preamble: $preamble_count unnumbered closed block (line $preamble_line) excluded from rows"
else
  echo "preamble: 0 unnumbered closed block (line -) excluded from rows"
fi

# Compute census totals.
closed=0; open=0; untagged=0
landed=0; routed=0; ruledout=0; parked=0; in_sprint=0; next=0
total=0
while IFS= read -r r; do
  [ -z "$r" ] && continue
  IFS='|' read -r num tag cls line <<< "$r"
  total=$((total+1))
  case "$cls" in
    closed) closed=$((closed+1)) ;;
    open) open=$((open+1)) ;;
    UNTAGGED) untagged=$((untagged+1)) ;;
  esac
  case "$tag" in
    LANDED) landed=$((landed+1)) ;;
    ROUTED) routed=$((routed+1)) ;;
    "RULED OUT") ruledout=$((ruledout+1)) ;;
    PARKED) parked=$((parked+1)) ;;
    IN-SPRINT) in_sprint=$((in_sprint+1)) ;;
    NEXT) next=$((next+1)) ;;
  esac
done <<< "$rows"

# F6 — control adjudication (the table is already printed above).
lookup_row() {  # $1=num ; prints "tag cls" or "ABSENT"
  local num="$1" r rnum tag cls rline
  while IFS= read -r r; do
    [ -z "$r" ] && continue
    IFS='|' read -r rnum tag cls rline <<< "$r"
    if [ "$rnum" = "$num" ]; then
      echo "$tag $cls"
      return 0
    fi
  done <<< "$rows"
  echo "ABSENT"
  return 1
}

mismatch=0
res="$(lookup_row "$CONTROL_CLOSED")"
if [ "$res" = "ABSENT" ]; then
  echo "✗ CONTROL MISMATCH: row $CONTROL_CLOSED expect=closed got=ABSENT"
  mismatch=1
else
  set -- $res
  ctag="$1"; ccls="$2"
  if [ "$ccls" = "closed" ]; then
    echo "control: row=$CONTROL_CLOSED expect=closed got=$ctag ok"
  else
    echo "✗ CONTROL MISMATCH: row $CONTROL_CLOSED expect=closed got=$ctag"
    mismatch=1
  fi
fi

res="$(lookup_row "$CONTROL_OPEN")"
if [ "$res" = "ABSENT" ]; then
  echo "✗ CONTROL MISMATCH: row $CONTROL_OPEN expect=open got=ABSENT"
  mismatch=1
else
  set -- $res
  ctag="$1"; ccls="$2"
  if [ "$ccls" = "open" ]; then
    echo "control: row=$CONTROL_OPEN expect=open got=$ctag ok"
  else
    echo "✗ CONTROL MISMATCH: row $CONTROL_OPEN expect=open got=$ctag"
    mismatch=1
  fi
fi

if [ "$mismatch" -ne 0 ]; then
  exit 2
fi

# charter-blob: the blob OID of the exact bytes read. Loud `unavailable` token on failure.
blob="$("$GIT_BIN" hash-object "$DOC" 2>/dev/null)"
if [ -n "$blob" ]; then
  blob_field="charter-blob $blob"
else
  blob_field="charter-blob unavailable"
  echo "✗ charter-blob unavailable: $GIT_BIN hash-object $DOC failed" >&2
fi

# The total line — the only line a stamp may quote.
echo "census: $closed closed / $total rows (tagged-open $open, untagged $untagged) [LANDED $landed, ROUTED $routed, RULED OUT $ruledout; PARKED $parked, IN-SPRINT $in_sprint, NEXT $next] [queue_census.sh, $blob_field, controls $CONTROL_CLOSED/$CONTROL_OPEN]"
exit 0
