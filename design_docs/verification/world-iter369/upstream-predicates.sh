#!/bin/bash
# World iteration 369 Gate-2 re-read: upstream predicates.
# Union of current charter ailang#N references and iteration 368's enumerated list.
# All must be state=closed. Controls: known-closed 885; an open-issue enumeration
# (gh issue list --state open) proving the check can see OPEN issues.
set -u
CHARTER="${1:-design_docs/world-mission.md}"
PRIOR="476 477 486 493 498 517 611 633 764 885 941 984 1036 1037 1042 1043 1096 1160 1424 1542 1547 1557 1558 1559 1580 1602 1607 1635 1843 2053 2127 2132 2160"
charter_nums=$(grep -oE 'ailang#[0-9]+' "$CHARTER" | sed 's/ailang#//' | sort -n -u)
nums=$(printf '%s\n%s\n' "$charter_nums" "$PRIOR" | sort -n -u)
total=0; closed=0; open_list=""
for n in $nums; do
  total=$((total+1))
  st=$(gh api "repos/sunholo-data/ailang/issues/$n" --jq '.state' 2>/dev/null)
  if [ "$st" = "closed" ]; then closed=$((closed+1)); else open_list="$open_list $n($st)"; fi
done
echo "charter ailang#N distinct: $(printf '%s\n' "$charter_nums" | wc -l | tr -d ' ')"
echo "union distinct: $total"
echo "closed: $closed"
echo "open:${open_list:- none}"
echo "control closed-885: $(gh api repos/sunholo-data/ailang/issues/885 --jq '.state' 2>/dev/null) (expect closed)"
echo "control open-list (expect >=1, proves the enumeration sees OPEN issues):"
gh issue list --repo sunholo-data/ailang --state open --limit 3 --json number --jq '.[] | .number'
