#!/usr/bin/env bash
# test_check_no_personal_email.sh — 10 arms (A–J) for scripts/check_no_personal_email.sh (row 74).
#
# Rule 3k: every arm runs the REAL instrument as a process — `bash "$SCRIPT_UT"` from inside a
# throwaway `git init` repo (arms B–H) or the real repo (arm A) — and asserts rc AND the three
# contract lines (D6) via one `probe` string per arm: "rc=N hit=N ok=N floor=N", where hit = count
# of per-hit lines, ok = count of the clean line, floor = count of the empty-scope floor line.
# Every fixture address is ASSEMBLED AT RUNTIME by concatenation: a literal address in this file
# would be flagged by the very gate it tests, because this file lives in scripts/ and is in scope.
# bash 3.2 safe. Scratch lives inside the worktree (never /tmp, row 71) and is removed on exit.
set -uo pipefail
cd "$(dirname "$0")/.." || exit 1
ROOT="$(pwd)"
SCRIPT_UT="$ROOT/scripts/check_no_personal_email.sh"
pass=0; fail=0
ok()   { echo "ok $*";      pass=$((pass+1)); }
notok(){ echo "not ok $*";  fail=$((fail+1)); }
ck() { if [ "$2" = "$3" ]; then ok "$1"; else notok "$1 (got '$2' want '$3')"; fi; }

SCRATCH="$(mktemp -d "$ROOT/.cnpe_scratch.XXXXXX")"
trap 'cd "$ROOT"; rm -rf "$SCRATCH"' EXIT

# probe <instrument-path>: run it from the CURRENT directory, summarise rc + the contract lines.
probe() {
    local out rc
    out="$(bash "$1" 2>&1)"; rc=$?
    printf 'rc=%s hit=%s ok=%s floor=%s\n' "$rc" \
        "$(printf '%s\n' "$out" | /usr/bin/grep -c 'personal email in loop-written file:')" \
        "$(printf '%s\n' "$out" | /usr/bin/grep -c 'check-no-personal-email: no personal addresses in the loop-written surface')" \
        "$(printf '%s\n' "$out" | /usr/bin/grep -c 'zero in-scope files scanned')"
}
stage() { git add -A >/dev/null 2>&1; }

# Fixture addresses, assembled at runtime by splitting at the "@" (see header): neither half is
# address-shaped on its own, so the gate reads this file as clean.
PERSONAL="someone";      PERSONAL="${PERSONAL}@realdomain.com"
NOREPLY="3155884+SomeUser"; NOREPLY="${NOREPLY}@users.noreply.github.com"
RESERVED_COM="you";      RESERVED_COM="${RESERVED_COM}@example.com"
RESERVED_ORG="you";      RESERVED_ORG="${RESERVED_ORG}@example.org"
RESERVED_NET="you";      RESERVED_NET="${RESERVED_NET}@example.net"
RESERVED_INV="test";     RESERVED_INV="${RESERVED_INV}@example.invalid"
PREFIX_ONLY="a";         PREFIX_ONLY="${PREFIX_ONLY}@example.com.evil.net"
# Arm I: a personal domain that merely CONTAINS an exclusion token, or a local part that merely
# ENDS in noreply — every one must be a HIT (the round-1 judge's blocking finding).
SUBSTR_NOREPLY="attacker-noreply"; SUBSTR_NOREPLY="${SUBSTR_NOREPLY}@evil-personal-domain.com"
SUBSTR_GH="x";        SUBSTR_GH="${SUBSTR_GH}@users.noreply.github.com.evil-personal.net"
SUBSTR_GSA="y";       SUBSTR_GSA="${SUBSTR_GSA}@gserviceaccount.com.evil-personal.net"
SUBSTR_SENTRY="z";    SUBSTR_SENTRY="${SUBSTR_SENTRY}@sentry.io.evil-personal.net"
# Arm J: the legitimate machine identities the anchors must still admit (measured on both surfaces).
LEGIT_NOREPLY="noreply"; LEGIT_NOREPLY="${LEGIT_NOREPLY}@anthropic.com"
LEGIT_GSA="sa-cloudbuild"; LEGIT_GSA="${LEGIT_GSA}@some-project.iam.gserviceaccount.com"
LEGIT_SENTRY="alerts";   LEGIT_SENTRY="${LEGIT_SENTRY}@sentry.io"
LEGIT_BOT="151556158+sunholo-voight-kampff"; LEGIT_BOT="${LEGIT_BOT}@users.noreply.github.com"

# A — the real repo is clean (the gate's live assertion; this is what CI's first new step runs).
ck "A real repo clean" "$(probe "$SCRIPT_UT")" "rc=0 hit=0 ok=1 floor=0"

# B–E, G, H run in ONE throwaway repo whose scripts/ holds a copy of the instrument (in scope,
# clean by construction — the instrument's own text carries no address-shaped token).
W="$SCRATCH/repo"; mkdir -p "$W/design_docs" "$W/scripts"
cp "$SCRIPT_UT" "$W/scripts/check_no_personal_email.sh"
cd "$W" || exit 1
git init -q .
git -c user.email=t@t -c user.name=t commit -q --allow-empty -m init >/dev/null 2>&1
stage

# B — a personal address in a mission doc must FAIL (the whole point).
echo "provenance is the commit author $PERSONAL" > design_docs/world-mission.md; stage
ck "B personal addr in mission doc -> rc 1" "$(probe scripts/check_no_personal_email.sh)" "rc=1 hit=1 ok=0 floor=0"

# C — a GitHub noreply is ALLOWED (identities may still be recorded).
echo "author $NOREPLY" > design_docs/world-mission.md; stage
ck "C noreply allowed" "$(probe scripts/check_no_personal_email.sh)" "rc=0 hit=0 ok=1 floor=0"

# D — reserved-TLD placeholders are allowed (RFC 2606/6761): .com/.org/.net/.invalid.
echo "to: $RESERVED_COM $RESERVED_ORG $RESERVED_NET and $RESERVED_INV" > design_docs/world-mission.md; stage
ck "D reserved placeholders allowed" "$(probe scripts/check_no_personal_email.sh)" "rc=0 hit=0 ok=1 floor=0"

# E — OUT OF SCOPE by design: governance/code/planned-doc files are not policed. A clean mission
# doc stays in place so the scope is non-empty under every mutation (the floor is arm F's, not E's).
echo "# charter (clean)" > design_docs/world-mission.md
mkdir -p design_docs/planned
echo "maintainer $PERSONAL" > SECURITY.md
echo "const owner = \"$PERSONAL\"" > access_control.go
echo "reviewer $PERSONAL" > design_docs/planned/some-design.md
stage
ck "E governance/code/planned out of scope" "$(probe scripts/check_no_personal_email.sh)" "rc=0 hit=0 ok=1 floor=0"

# G — a hit in scripts/ (not only in a mission doc) is caught: scripts/.* is a live alternative.
rm -f SECURITY.md access_control.go design_docs/planned/some-design.md
echo "ATT_EMAIL=\"\${MISSION_ATTENDED_EMAIL:-$PERSONAL}\"" > scripts/some_tool.sh; stage
ck "G hit in scripts/ caught" "$(probe scripts/check_no_personal_email.sh)" "rc=1 hit=1 ok=0 floor=0"

# H — the reserved-TLD exclusion is ANCHORED: example.com as a mere PREFIX of the domain is a HIT.
rm -f scripts/some_tool.sh
echo "contact $PREFIX_ONLY" > design_docs/world-mission.md; stage
ck "H example.com-as-prefix is a HIT (anchored)" "$(probe scripts/check_no_personal_email.sh)" "rc=1 hit=1 ok=0 floor=0"

# I — a domain that merely CONTAINS an exclusion token, or a local part merely ENDING in noreply,
# is a HIT: every exclusion clause is anchored to the whole token. Four tokens, four hit lines.
printf 'a %s\nb %s\nc %s\nd %s\n' "$SUBSTR_NOREPLY" "$SUBSTR_GH" "$SUBSTR_GSA" "$SUBSTR_SENTRY" > design_docs/world-mission.md; stage
ck "I substring-of-exclusion domains -> rc 1, 4 hits" "$(probe scripts/check_no_personal_email.sh)" "rc=1 hit=4 ok=0 floor=0"

# J — the anchors must still ADMIT every legitimate machine identity (an anchor that is too tight
# reds the real surface, which carries the bot's noreply identity and a service account).
printf 'a %s\nb %s\nc %s\nd %s\n' "$LEGIT_NOREPLY" "$LEGIT_GSA" "$LEGIT_SENTRY" "$LEGIT_BOT" > design_docs/world-mission.md; stage
ck "J legitimate machine identities still allowed" "$(probe scripts/check_no_personal_email.sh)" "rc=0 hit=0 ok=1 floor=0"

# F — the anti-vacuity floor: a repo with ZERO in-scope tracked files reads rc 2, no ✓ line.
# The instrument lives OUTSIDE scripts/ here (tools/ is not in SCOPE_RE) so the scope is empty.
W2="$SCRATCH/empty"; mkdir -p "$W2/tools"
cp "$SCRIPT_UT" "$W2/tools/check_no_personal_email.sh"
cd "$W2" || exit 1
git init -q .
echo "# nothing the loop writes" > README.md; stage
ck "F empty scope -> rc 2, floor line, no clean line" "$(probe tools/check_no_personal_email.sh)" "rc=2 hit=0 ok=0 floor=1"

cd "$ROOT" || exit 1
echo "---"
echo "$pass passed, $fail failed"
[ "$fail" -eq 0 ]
