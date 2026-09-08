#!/usr/bin/env bash
# test_queue_census.sh — self-contained arms for scripts/queue_census.sh.
#
# Every arm is SELF-CONTAINED: it builds a literal-heredoc fixture and runs the instrument
# against it. There is NO helper that synthesises a heading, a row number or a tag — the only
# fixture builder is make_doc, which reads literal lines from a heredoc (D4's iter-108
# defence: a floor whose input shape the harness cannot produce is a floor no arm can reach).
# Every invocation of the instrument goes through an explicit `bash "$SCRIPT_UT" ...` —
# never the login shell (zsh) and never `sh`.
set -uo pipefail
cd "$(dirname "$0")/.." || exit 1
SCRIPT_UT=scripts/queue_census.sh
pass=0; fail=0
ok()   { echo "ok $*";      pass=$((pass+1)); }
notok(){ echo "not ok $*";  fail=$((fail+1)); }

# Scratch lives in the worktree (NOT /tmp — row 71's rule), cleaned up on exit.
SCRATCH="$(mktemp -d "$(pwd)/.qc_scratch.XXXXXX")"
trap 'rm -rf "$SCRATCH"' EXIT

# ── helpers ──────────────────────────────────────────────────────────────────
# make_doc <path> — the ONLY fixture builder. Reads literal lines from a heredoc.
make_doc() {
  cat > "$1"
}

# check_row_order <out> <expected nums space-separated> — the row= lines appear in order.
check_row_order() {
  local out="$1" expected="$2" got
  got="$(echo "$out" | /usr/bin/grep '^row=' | sed -E 's/^row=([0-9]+).*/\1/' | tr '\n' ' ' | sed 's/ $//')"
  [ "$got" = "$expected" ]
}

# ── arms ─────────────────────────────────────────────────────────────────────
# AC-M1-2 — enumeration is order-independent and line-accurate.
arm_enumerate() {
  local f out rc
  f="$SCRATCH/doc.md"
  make_doc "$f" <<'EOF'
## Queue (top = next; tags: [NEXT] [IN-SPRINT] [PARKED] [LANDED] [RULED OUT])
3. [LANDED 2026-09-08] **w-three** · clause-2 ·
1. [LANDED 2026-07-24] **w-one** · clause-1 ·
2. [PARKED — DESIGN REVIEW] **w-two** · clause-2 ·
EOF
  out="$(bash "$SCRIPT_UT" --doc "$f" --control-closed 1 --control-open 2 2>&1)"
  rc=$?
  if [ "$rc" -eq 0 ]; then ok "enumerate exits 0"; else notok "enumerate rc=$rc: $out"; fi
  echo "$out" | /usr/bin/grep -q 'row=1 tag=LANDED class=closed line=3' && ok "enumerate row1 line=3" || notok "enumerate row1 wrong: $out"
  echo "$out" | /usr/bin/grep -q 'row=2 tag=PARKED class=open line=4' && ok "enumerate row2 line=4" || notok "enumerate row2 wrong: $out"
  echo "$out" | /usr/bin/grep -q 'row=3 tag=LANDED class=closed line=2' && ok "enumerate row3 line=2" || notok "enumerate row3 wrong: $out"
  check_row_order "$out" "1 2 3" && ok "enumerate sorted order" || notok "enumerate not sorted: $out"
  echo "$out" | /usr/bin/grep -q '2 closed / 3 rows (tagged-open 1, untagged 0)' && ok "enumerate census" || notok "enumerate census wrong: $out"
}

# AC-M1-3 — the closed vocabulary and every decoration shape classify closed.
arm_closed_vocab() {
  local f out rc nclosed
  f="$SCRATCH/doc.md"
  make_doc "$f" <<'EOF'
## Queue (top = next; tags: [NEXT] [IN-SPRINT] [PARKED] [LANDED] [RULED OUT])
1. [LANDED 2026-09-08] **w1** · clause-2 ·
2. [**LANDED 2026-09-08] **w2** · clause-2 ·
3. [**[LANDED 2026-09-08] **w3** · clause-2 ·
4. **[LANDED 2026-09-08] **w4** · clause-2 ·
5. **[ROUTED 2026-09-08] **w5** · clause-2 ·
6. **[RULED OUT 2026-09-08] **w6** · clause-2 ·
7. [PARKED] control
EOF
  out="$(bash "$SCRIPT_UT" --doc "$f" --control-closed 1 --control-open 7 2>&1)"
  rc=$?
  if [ "$rc" -eq 0 ]; then ok "closed-vocab exits 0"; else notok "closed-vocab rc=$rc: $out"; fi
  nclosed="$(echo "$out" | /usr/bin/grep -c 'class=closed')"
  if [ "$nclosed" -eq 6 ]; then ok "closed-vocab six closed"; else notok "closed-vocab $nclosed closed"; fi
  echo "$out" | /usr/bin/grep -q '6 closed / 7 rows (tagged-open 1, untagged 0)' && ok "closed-vocab census" || notok "closed-vocab census wrong: $out"
}

# AC-M1-4 — the open vocabulary classifies open, not closed.
arm_open_vocab() {
  local f out rc nopen
  f="$SCRATCH/doc.md"
  make_doc "$f" <<'EOF'
## Queue (top = next; tags: [NEXT] [IN-SPRINT] [PARKED] [LANDED] [RULED OUT])
1. [PARKED — DESIGN REVIEW] **w1** · clause-2 ·
2. **[NEXT] **w2** · clause-2 ·
3. [**[IN-SPRINT] **w3** · clause-2 ·
4. [LANDED] control
EOF
  out="$(bash "$SCRIPT_UT" --doc "$f" --control-closed 4 --control-open 1 2>&1)"
  rc=$?
  if [ "$rc" -eq 0 ]; then ok "open-vocab exits 0"; else notok "open-vocab rc=$rc: $out"; fi
  nopen="$(echo "$out" | /usr/bin/grep -c 'class=open')"
  if [ "$nopen" -eq 3 ]; then ok "open-vocab three open"; else notok "open-vocab $nopen open"; fi
  echo "$out" | /usr/bin/grep -q '1 closed / 4 rows (tagged-open 3, untagged 0)' && ok "open-vocab census" || notok "open-vocab census wrong: $out"
}

# AC-M1-5 — an untagged row is UNTAGGED, counted in rows, not in closed.
arm_untagged() {
  local f out rc
  f="$SCRATCH/doc.md"
  make_doc "$f" <<'EOF'
## Queue (top = next; tags: [NEXT] [IN-SPRINT] [PARKED] [LANDED] [RULED OUT])
1. [LANDED 2026-09-08] **w1** · clause-2 ·
2. [PARKED — DESIGN REVIEW] **w2** · clause-2 ·
3. **w-slug** · clause-2 · prose
EOF
  out="$(bash "$SCRIPT_UT" --doc "$f" --control-closed 1 --control-open 2 2>&1)"
  rc=$?
  if [ "$rc" -eq 0 ]; then ok "untagged exits 0"; else notok "untagged rc=$rc: $out"; fi
  echo "$out" | /usr/bin/grep -q 'row=3 tag=- class=UNTAGGED' && ok "untagged row 3 UNTAGGED" || notok "untagged row 3 not UNTAGGED: $out"
  echo "$out" | /usr/bin/grep -q '1 closed / 3 rows (tagged-open 1, untagged 1)' && ok "untagged census" || notok "untagged census wrong: $out"
}

# AC-M1-6 — prose is never a tag (the row-57 defect).
arm_prose_not_tag() {
  local f out rc
  f="$SCRATCH/doc.md"
  make_doc "$f" <<'EOF'
## Queue (top = next; tags: [NEXT] [IN-SPRINT] [PARKED] [LANDED] [RULED OUT])
1. [LANDED] control
2. [PARKED] control
3. **w-prose** · clause-2 · the body contains [LANDED] [ROUTED] [RULED OUT] [PARKED] [NEXT] and it landed correctly and LANDED 2026-09-08 (iter-1)
4. **w-slug** · **[LANDED 2026-09-08]** · clause-2 ·
EOF
  out="$(bash "$SCRIPT_UT" --doc "$f" --control-closed 1 --control-open 2 2>&1)"
  rc=$?
  if [ "$rc" -eq 0 ]; then ok "prose-not-tag exits 0"; else notok "prose-not-tag rc=$rc: $out"; fi
  echo "$out" | /usr/bin/grep -q 'row=3 tag=- class=UNTAGGED' && ok "prose-not-tag row 3 UNTAGGED" || notok "prose-not-tag row 3 not UNTAGGED: $out"
  echo "$out" | /usr/bin/grep -q 'row=4 tag=- class=UNTAGGED' && ok "prose-not-tag row 4 UNTAGGED" || notok "prose-not-tag row 4 not UNTAGGED: $out"
  echo "$out" | /usr/bin/grep -q '1 closed / 4 rows (tagged-open 1, untagged 2)' && ok "prose-not-tag census" || notok "prose-not-tag census wrong: $out"
}

# AC-M1-7 — a `## ` line inside a row body does not end the enumeration.
arm_body_heading() {
  local f out rc
  f="$SCRATCH/doc.md"
  make_doc "$f" <<'EOF'
## Queue (top = next; tags: [NEXT] [IN-SPRINT] [PARKED] [LANDED] [RULED OUT])
1. [LANDED 2026-09-08] **w1** · clause-2 ·
## Premise Verification Log
2. [PARKED — DESIGN REVIEW] **w2** · clause-2 ·
3. [LANDED 2026-09-08] **w3** · clause-2 ·
EOF
  out="$(bash "$SCRIPT_UT" --doc "$f" --control-closed 1 --control-open 2 2>&1)"
  rc=$?
  if [ "$rc" -eq 0 ]; then ok "body-heading exits 0"; else notok "body-heading rc=$rc: $out"; fi
  echo "$out" | /usr/bin/grep -q '2 closed / 3 rows (tagged-open 1, untagged 0)' && ok "body-heading census" || notok "body-heading census wrong: $out"
}

# AC-M1-8 — the unnumbered preamble block is reported and excluded (two fixtures).
arm_preamble() {
  local fA fB out rc
  fA="$SCRATCH/docA.md"
  make_doc "$fA" <<'EOF'
## Queue (top = next; tags: [NEXT] [IN-SPRINT] [PARKED] [LANDED] [RULED OUT])
**[LANDED 2026-07-24 (iter-13)] w-m1-ailang-hardening** · clause-1 ·
1. [LANDED 2026-09-08] **w1** · clause-2 ·
2. [PARKED — DESIGN REVIEW] **w2** · clause-2 ·
EOF
  out="$(bash "$SCRIPT_UT" --doc "$fA" --control-closed 1 --control-open 2 2>&1)"
  rc=$?
  if [ "$rc" -eq 0 ]; then ok "preamble A exits 0"; else notok "preamble A rc=$rc: $out"; fi
  echo "$out" | /usr/bin/grep -q 'preamble: 1 unnumbered closed block (line 2) excluded from rows' && ok "preamble A reports 1 block line 2" || notok "preamble A wrong: $out"
  echo "$out" | /usr/bin/grep -q '1 closed / 2 rows (tagged-open 1, untagged 0)' && ok "preamble A census" || notok "preamble A census wrong: $out"
  fB="$SCRATCH/docB.md"
  make_doc "$fB" <<'EOF'
## Queue (top = next; tags: [NEXT] [IN-SPRINT] [PARKED] [LANDED] [RULED OUT])
1. [LANDED 2026-09-08] **w1** · clause-2 ·
2. [PARKED — DESIGN REVIEW] **w2** · clause-2 ·
EOF
  out="$(bash "$SCRIPT_UT" --doc "$fB" --control-closed 1 --control-open 2 2>&1)"
  rc=$?
  if [ "$rc" -eq 0 ]; then ok "preamble B exits 0"; else notok "preamble B rc=$rc: $out"; fi
  echo "$out" | /usr/bin/grep -q 'preamble: 0 ' && ok "preamble B reports 0" || notok "preamble B no 0: $out"
}

# AC-M1-9 — F0: mission doc missing/unreadable.
arm_no_doc() {
  local out rc
  out="$(bash "$SCRIPT_UT" --doc /nonexistent/x.md 2>&1)"
  rc=$?
  if [ "$rc" -eq 2 ]; then ok "no-doc exits 2"; else notok "no-doc rc=$rc: $out"; fi
  echo "$out" | /usr/bin/grep -q '✗ mission doc unreadable: /nonexistent/x.md' && ok "no-doc message" || notok "no-doc no message: $out"
}

# AC-M1-9 — F1: no line matches `^## Queue`.
arm_no_heading() {
  local f out rc
  f="$SCRATCH/doc.md"
  make_doc "$f" <<'EOF'
1. [LANDED 2026-09-08] **w1** · clause-2 ·
2. [PARKED — DESIGN REVIEW] **w2** · clause-2 ·
EOF
  out="$(bash "$SCRIPT_UT" --doc "$f" 2>&1)"
  rc=$?
  if [ "$rc" -eq 2 ]; then ok "no-heading exits 2"; else notok "no-heading rc=$rc: $out"; fi
  echo "$out" | /usr/bin/grep -q "✗ no \"## Queue\" heading in $f" && ok "no-heading message" || notok "no-heading no message: $out"
}

# AC-M1-9 — F2: zero rows enumerated after the heading.
arm_zero_rows() {
  local f out rc
  f="$SCRATCH/doc.md"
  make_doc "$f" <<'EOF'
## Queue (top = next; tags: [NEXT] [IN-SPRINT] [PARKED] [LANDED] [RULED OUT])
Some prose that is not a queue row.
EOF
  out="$(bash "$SCRIPT_UT" --doc "$f" 2>&1)"
  rc=$?
  if [ "$rc" -eq 2 ]; then ok "zero-rows exits 2"; else notok "zero-rows rc=$rc: $out"; fi
  echo "$out" | /usr/bin/grep -q '✗ zero queue rows enumerated after line 1' && ok "zero-rows message" || notok "zero-rows no message: $out"
}

# AC-M1-9 — F3: a row number appears twice.
arm_dup_number() {
  local f out rc
  f="$SCRATCH/doc.md"
  make_doc "$f" <<'EOF'
## Queue (top = next; tags: [NEXT] [IN-SPRINT] [PARKED] [LANDED] [RULED OUT])
1. [LANDED 2026-09-08] **w1** · clause-2 ·
2. [PARKED — DESIGN REVIEW] **w2** · clause-2 ·
2. [LANDED 2026-09-08] **w2b** · clause-2 ·
EOF
  out="$(bash "$SCRIPT_UT" --doc "$f" 2>&1)"
  rc=$?
  if [ "$rc" -eq 2 ]; then ok "dup-number exits 2"; else notok "dup-number rc=$rc: $out"; fi
  echo "$out" | /usr/bin/grep -q '✗ duplicate row number: 2 (lines' && ok "dup-number message" || notok "dup-number no message: $out"
}

# AC-M1-9 — F4: numbering is not 1..max.
arm_gap() {
  local f out rc
  f="$SCRATCH/doc.md"
  make_doc "$f" <<'EOF'
## Queue (top = next; tags: [NEXT] [IN-SPRINT] [PARKED] [LANDED] [RULED OUT])
1. [LANDED 2026-09-08] **w1** · clause-2 ·
3. [LANDED 2026-09-08] **w3** · clause-2 ·
EOF
  out="$(bash "$SCRIPT_UT" --doc "$f" 2>&1)"
  rc=$?
  if [ "$rc" -eq 2 ]; then ok "gap exits 2"; else notok "gap rc=$rc: $out"; fi
  echo "$out" | /usr/bin/grep -q '✗ numbering gap: expected 1..3, missing 2' && ok "gap message" || notok "gap no message: $out"
}

# AC-M1-10 — F5: controls are mandatory.
arm_controls_required() {
  local f out rc
  f="$SCRATCH/doc.md"
  make_doc "$f" <<'EOF'
## Queue (top = next; tags: [NEXT] [IN-SPRINT] [PARKED] [LANDED] [RULED OUT])
1. [LANDED 2026-09-08] **w1** · clause-2 ·
2. [PARKED — DESIGN REVIEW] **w2** · clause-2 ·
EOF
  out="$(bash "$SCRIPT_UT" --doc "$f" 2>&1)"
  rc=$?
  if [ "$rc" -eq 2 ]; then ok "controls-required exits 2"; else notok "controls-required rc=$rc: $out"; fi
  echo "$out" | /usr/bin/grep -q '✗ controls required: --control-closed <n> --control-open <n>' && ok "controls-required message" || notok "controls-required no message: $out"
}

# AC-M1-10 — F6: a control row with the wrong class, table still printed above.
arm_control_mismatch() {
  local f out rc
  f="$SCRATCH/doc.md"
  make_doc "$f" <<'EOF'
## Queue (top = next; tags: [NEXT] [IN-SPRINT] [PARKED] [LANDED] [RULED OUT])
1. [LANDED 2026-09-08] **w1** · clause-2 ·
2. [PARKED — DESIGN REVIEW] **w2** · clause-2 ·
3. **w-slug** · clause-2 · prose
EOF
  out="$(bash "$SCRIPT_UT" --doc "$f" --control-closed 2 --control-open 2 2>&1)"
  rc=$?
  if [ "$rc" -eq 2 ]; then ok "control-mismatch run1 exits 2"; else notok "control-mismatch run1 rc=$rc: $out"; fi
  echo "$out" | /usr/bin/grep -q '✗ CONTROL MISMATCH: row 2 expect=closed got=PARKED' && ok "control-mismatch run1 names mismatch" || notok "control-mismatch run1 no mismatch: $out"
  echo "$out" | /usr/bin/grep -q '^row=1 ' && ok "control-mismatch run1 table printed" || notok "control-mismatch run1 table suppressed"
  out="$(bash "$SCRIPT_UT" --doc "$f" --control-closed 1 --control-open 3 2>&1)"
  rc=$?
  if [ "$rc" -eq 2 ]; then ok "control-mismatch run2 exits 2"; else notok "control-mismatch run2 rc=$rc: $out"; fi
  echo "$out" | /usr/bin/grep -q '✗ CONTROL MISMATCH: row 3 expect=open got=-' && ok "control-mismatch run2 names mismatch" || notok "control-mismatch run2 no mismatch: $out"
  echo "$out" | /usr/bin/grep -q '^row=1 ' && ok "control-mismatch run2 table printed" || notok "control-mismatch run2 table suppressed"
}

# AC-M1-11 — the report format (D5).
arm_report() {
  local f out rc last_line
  f="$SCRATCH/doc.md"
  make_doc "$f" <<'EOF'
## Queue (top = next; tags: [NEXT] [IN-SPRINT] [PARKED] [LANDED] [RULED OUT])
1. [LANDED 2026-09-08] **w1** · clause-2 ·
2. **[ROUTED 2026-09-08] **w2** · clause-2 ·
3. [PARKED — DESIGN REVIEW] **w3** · clause-2 ·
4. **w-slug** · clause-2 · prose
EOF
  out="$(bash "$SCRIPT_UT" --doc "$f" --control-closed 1 --control-open 3 2>&1)"
  rc=$?
  if [ "$rc" -eq 0 ]; then ok "report exits 0"; else notok "report rc=$rc: $out"; fi
  last_line="$(echo "$out" | tail -1)"
  echo "$last_line" | /usr/bin/grep -E '^census: 2 closed / 4 rows \(tagged-open 1, untagged 1\) \[LANDED 1, ROUTED 1, RULED OUT 0; PARKED 1, IN-SPRINT 0, NEXT 0\] \[queue_census\.sh, charter-blob [0-9a-f]{40}, controls 1/3\]$' && ok "report census line exact" || notok "report census line: $last_line"
  echo "$out" | /usr/bin/grep -q 'control: row=1 expect=closed got=LANDED ok' && ok "report control closed ok" || notok "report no closed control ok"
  echo "$out" | /usr/bin/grep -q 'control: row=3 expect=open got=PARKED ok' && ok "report control open ok" || notok "report no open control ok"
  echo "$out" | /usr/bin/grep -q 'carried' && notok "report contains carried" || ok "report no carried"
  echo "$out" | /usr/bin/grep -q 'of 4 rows closed' && notok "report contains 'of 4 rows closed'" || ok "report no carried-prose form"
}

# AC-M1-15 — the printed blob OID is the input file's blob.
arm_blob_oid() {
  local f out rc expected
  f="$SCRATCH/doc.md"
  make_doc "$f" <<'EOF'
## Queue (top = next; tags: [NEXT] [IN-SPRINT] [PARKED] [LANDED] [RULED OUT])
1. [LANDED 2026-09-08] **w1** · clause-2 ·
2. [PARKED — DESIGN REVIEW] **w2** · clause-2 ·
EOF
  out="$(bash "$SCRIPT_UT" --doc "$f" --control-closed 1 --control-open 2 2>&1)"
  rc=$?
  if [ "$rc" -eq 0 ]; then ok "blob-oid exits 0"; else notok "blob-oid rc=$rc: $out"; fi
  expected="$(git hash-object "$f")"
  echo "$out" | /usr/bin/grep -q "charter-blob $expected" && ok "blob-oid matches git hash-object" || notok "blob-oid mismatch: expected $expected"
}

# AC-M1-15 — a missing git binary is loud, still exit 0.
arm_blob_unavailable() {
  local f out rc
  f="$SCRATCH/doc.md"
  make_doc "$f" <<'EOF'
## Queue (top = next; tags: [NEXT] [IN-SPRINT] [PARKED] [LANDED] [RULED OUT])
1. [LANDED 2026-09-08] **w1** · clause-2 ·
2. [PARKED — DESIGN REVIEW] **w2** · clause-2 ·
EOF
  out="$(bash "$SCRIPT_UT" --doc "$f" --control-closed 1 --control-open 2 --git-bin /nonexistent/git 2>&1)"
  rc=$?
  if [ "$rc" -eq 0 ]; then ok "blob-unavailable exits 0"; else notok "blob-unavailable rc=$rc: $out"; fi
  echo "$out" | /usr/bin/grep -q 'charter-blob unavailable' && ok "blob-unavailable token present" || notok "blob-unavailable no token: $out"
  echo "$out" | /usr/bin/grep -q '✗ charter-blob unavailable' && ok "blob-unavailable stderr names it" || notok "blob-unavailable no stderr line"
}

# AC-M1-12 — the heading anchor is the `^## Queue` prefix (post-M2 heading text).
arm_heading_tail() {
  local f out rc
  f="$SCRATCH/doc.md"
  make_doc "$f" <<'EOF'
## Queue (top = next; tags: [NEXT] [IN-SPRINT] [PARKED] [LANDED] [RULED OUT] [ROUTED])
1. [LANDED 2026-09-08] **w1** · clause-2 ·
2. [PARKED — DESIGN REVIEW] **w2** · clause-2 ·
EOF
  out="$(bash "$SCRIPT_UT" --doc "$f" --control-closed 1 --control-open 2 2>&1)"
  rc=$?
  if [ "$rc" -eq 0 ]; then ok "heading-tail exits 0"; else notok "heading-tail rc=$rc: $out"; fi
  echo "$out" | /usr/bin/grep -q 'heading_line=1' && ok "heading-tail reports heading_line" || notok "heading-tail no heading_line: $out"
  echo "$out" | /usr/bin/grep -q '1 closed / 2 rows (tagged-open 1, untagged 0)' && ok "heading-tail census" || notok "heading-tail census wrong: $out"
}

# AC-M1-13 — the frozen charter snapshot reproduces the prototype.
arm_snapshot() {
  local out rc n
  out="$(bash "$SCRIPT_UT" --doc scripts/testdata/queue_census_rowstarts_8b15153.md --control-closed 1 --control-open 79 2>&1)"
  rc=$?
  if [ "$rc" -eq 0 ]; then ok "snapshot exits 0"; else notok "snapshot rc=$rc: $out"; fi
  echo "$out" | /usr/bin/grep -q 'census: 41 closed / 89 rows (tagged-open 4, untagged 44)' && ok "snapshot reproduces 41/89/4/44" || notok "snapshot wrong census: $out"
  echo "$out" | /usr/bin/grep -q 'preamble: 1 ' && ok "snapshot preamble 1" || notok "snapshot no preamble"
  echo "$out" | /usr/bin/grep -q 'control: row=1 expect=closed got=LANDED ok' && ok "snapshot control closed ok" || notok "snapshot no closed control"
  echo "$out" | /usr/bin/grep -q 'control: row=79 expect=open got=PARKED ok' && ok "snapshot control open ok" || notok "snapshot no open control"
  for n in 16 18 66 68 70; do
    echo "$out" | /usr/bin/grep -q "row=$n tag=- class=UNTAGGED" && ok "snapshot row $n UNTAGGED" || notok "snapshot row $n not UNTAGGED"
  done
}

# ── dispatcher ───────────────────────────────────────────────────────────────
run_arm() {
  case "$1" in
    enumerate) arm_enumerate ;;
    closed-vocab) arm_closed_vocab ;;
    open-vocab) arm_open_vocab ;;
    untagged) arm_untagged ;;
    prose-not-tag) arm_prose_not_tag ;;
    body-heading) arm_body_heading ;;
    preamble) arm_preamble ;;
    no-doc) arm_no_doc ;;
    no-heading) arm_no_heading ;;
    zero-rows) arm_zero_rows ;;
    dup-number) arm_dup_number ;;
    gap) arm_gap ;;
    controls-required) arm_controls_required ;;
    control-mismatch) arm_control_mismatch ;;
    report) arm_report ;;
    blob-oid) arm_blob_oid ;;
    blob-unavailable) arm_blob_unavailable ;;
    heading-tail) arm_heading_tail ;;
    snapshot) arm_snapshot ;;
    *) echo "unknown arm: $1" >&2; exit 2 ;;
  esac
}

ONLY=""
if [ "$#" -gt 0 ]; then
  if [ "$1" = "--only" ]; then
    ONLY="$2"
  else
    echo "unknown argument: $1" >&2
    exit 2
  fi
fi

if [ -n "$ONLY" ]; then
  run_arm "$ONLY"
else
  for arm in enumerate closed-vocab open-vocab untagged prose-not-tag body-heading preamble no-doc no-heading zero-rows dup-number gap controls-required control-mismatch report blob-oid blob-unavailable heading-tail snapshot; do
    run_arm "$arm"
  done
fi

echo "---"
echo "$pass passed, $fail failed"
[ "$fail" -eq 0 ]
