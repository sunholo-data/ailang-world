#!/bin/bash
# poll_pr.sh <pr> <full-sha> <verdict-file>: wait for all check-runs on sha to complete; 30-min cap.
PR=$1; SHA=$2; V=$3; R=sunholo-data/ailang-world
case "$SHA" in ????????????????????????????????????????) ;; *) echo "INSTRUMENT: sha not 40 chars" > "$V"; exit 2;; esac
deadline=$(( $(date +%s) + 1800 ))
while :; do
  j=$(gh api "repos/$R/commits/$SHA/check-runs" 2>/dev/null)
  tot=$(echo "$j" | jq -r '.total_count // empty'); pend=$(echo "$j" | jq -r '[.check_runs[]|select(.status!="completed")]|length')
  case "$tot" in ''|*[!0-9]*) tot=x;; esac
  if [ "$tot" != x ] && [ "$tot" -ge 3 ] && [ "$pend" = 0 ]; then
    echo "ALL COMPLETE for $SHA total=$tot" > "$V"; echo "$j" | jq -r '.check_runs[]|"  \(.name): \(.conclusion)"' >> "$V"; exit 0; fi
  [ "$(date +%s)" -ge "$deadline" ] && { echo "TIMEOUT for $SHA after 30m total=$tot pending=$pend — PARK" > "$V"; exit 1; }
  sleep 30
done
