#!/usr/bin/env bash
# check_no_personal_email.sh — refuse a personal email address in the surface the LOOP writes.
#
# WHY (the shared mission-control skill's ATTENDED LEDGER EDITS contract, Mark, attended
# 2026-09-02): ledger rows are pasted VERBATIM into the public bookkeeping issue by every mission
# report, so an address written into an evidence cell is published. The rule says compare, never
# record; this gate is what makes that stick, because the generator is automated and prose review
# is not. World carried a personal address in 7 doc locations at iter-152 (row 74); the redaction
# was a one-time sweep with nothing behind it until this gate.
#
# SCOPE IS DELIBERATELY NARROW: the loop-written surface only (World's own SCOPE_RE, row 74 D2):
#   * design_docs/*mission*.md — the six mission docs the loop writes (charter, log, archives,
#     dashboard, index)
#   * scripts/**               — World-owned instruments the loop writes and CI runs
# It does NOT police design_docs/planned|implemented|verification (sprint artifacts; a full-tree
# scan reads 7 non-address false positives there today — SSH remote strings, Go toolchain
# module ids), Go/.ail sources, workflows, go.mod/go.sum, docs/**, README.md or CLAUDE.md
# (deliberate public contact points — a governance decision, not a lint). Widening this without
# ruling on those first would just fail the build; it is a separate queue row.
#
# Allowed: GitHub noreply addresses (attributable, expose nothing), example.com/org/net and
# .invalid/.test/.localhost placeholders (RFC 2606/6761 reserved, cannot be real), and machine
# identities such as GCP service accounts — none of which is a person. The reserved-TLD exclusion
# is ANCHORED (`@example\.(com|org|net)$`): an address at example.com.evil.net is a HIT. World's one
# deliberate deviation from the fleet precedent (sunholo-data/ailang f0d44915d), pinned by arm H.
#
# Exit codes: 0 clean; 1 at least one address hit; 2 instrument failure (zero in-scope files —
# a gate that scans nothing must never print ✓).
# bash 3.2 safe: set -u, no associative arrays, no mapfile, no ${var,,}. Colour only on a TTY so
# captured output is the bare contract lines.

set -u
if [ -t 1 ]; then RED=$'\033[0;31m'; GREEN=$'\033[0;32m'; RESET=$'\033[0m'; else RED=''; GREEN=''; RESET=''; fi
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT" || exit 1

# The loop writes here (World's own scope, row 74 D2 — no .claude/skills/: World has no such tree).
SCOPE_RE='^(design_docs/[^/]*mission[^/]*\.md|scripts/.*)$'
PAT='[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}'

files="$(git ls-files | /usr/bin/grep -E "$SCOPE_RE")"
n_scope=$(printf '%s\n' "$files" | /usr/bin/grep -c .)
if [ "$n_scope" -eq 0 ]; then
    echo "${RED}✗ check-no-personal-email: zero in-scope files scanned — the gate must never print ✓ over an empty surface${RESET}"
    exit 2
fi

hits=0
while IFS= read -r f; do
    [ -n "$f" ] || continue
    case "$f" in
        *.png|*.jpg|*.gif|*.pdf|*.zip) continue ;;
    esac
    [ -f "$f" ] || continue
    found=$(LC_ALL=C /usr/bin/grep -oE "$PAT" "$f" 2>/dev/null \
        | /usr/bin/grep -vE 'users\.noreply\.github\.com|noreply@|@example\.(com|org|net)$|\.(invalid|test|localhost)$|gserviceaccount\.com|@sentry\.io' \
        | sort -u)
    if [ -n "$found" ]; then
        while IFS= read -r addr; do
            [ -n "$addr" ] || continue
            echo "${RED}✗ personal email in loop-written file:${RESET} $f  ->  $addr"
            hits=$((hits+1))
        done <<< "$found"
    fi
done <<< "$files"

if [ "$hits" -gt 0 ]; then
    echo ""
    echo "${RED}check-no-personal-email FAILED${RESET} ($hits occurrence(s) in $n_scope in-scope files)."
    echo "Ledger rows are pasted verbatim into the public bookkeeping issue — an address here is published."
    echo "Record the provenance VERDICT (\"attended\" / \"fleet\"), never the address."
    echo "Use a GitHub noreply address where an identity is genuinely needed."
    exit 1
fi
echo "${GREEN}✓ check-no-personal-email: no personal addresses in the loop-written surface${RESET}"
exit 0
