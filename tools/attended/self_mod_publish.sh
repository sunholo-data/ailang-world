#!/bin/bash
# Attended helper for SM.D — the irreversible first publish of world/core@0.1.0.
#
# WHAT THIS IS NOT: an automation of the publish. Every fence in
# cmd/world-publish stays exactly where it was, and this script cannot satisfy
# any of them on your behalf. You still type the confirmation phrase yourself,
# at a real terminal, and `approve` still refuses if there is not one.
#
# WHAT IT IS FOR: the two measured failures of the first attended attempt
# (2026-09-21) were both ENVIRONMENT, not decision —
#   STOP fence=store reason=unopenable                  (the runbook omitted a mkdir)
#   STOP fence=tty   reason=stdin-is-not-the-controlling-terminal
# The second is the one this script exists for. `approve` requires stdin to BE
# the controlling terminal, and a great many ordinary ways of running a command
# do not give it one: a pasted multi-line block, a command run through an editor
# or agent shell, anything with a redirect. Connecting stdin TO /dev/tty is not
# a way around that fence — it is the fence's own requirement, honestly met. The
# human still types the phrase.
#
# WHAT IT DELIBERATELY DOES NOT DO: mint and spend in one invocation. The
# runbook is explicit that "one command that did both would collapse two human
# acts into one keystroke and hide the single-use property at the very surface
# where it matters". Hence subcommands, one human act each.
#
# Runbook: docs/SELF_MOD_PUBLISH.md. Where they disagree, the runbook wins.
set -euo pipefail

# ../.. — this file lives at tools/attended/, TWO levels down. It was one level
# down under scripts/, and the move left this at ".." so every repo-relative
# path resolved under tools/. Caught by the golden read failing, which is the
# only reason it was not a silently wrong root.
REPO_ROOT="$(cd "$(dirname "$0")/../.." && pwd)"
STORE="${WORLD_STORE:-$HOME/.ailang/world/world.db}"
REGISTRY="${WORLD_REGISTRY:-https://storage.googleapis.com/ailang-registry}"
CREDENTIAL="${WORLD_CREDENTIAL:-$HOME/.config/ailang/registry.key}"
COMPILER="${WORLD_COMPILER:-$HOME/.pinned-ailang/ailang}"
BIN="${WORLD_BIN:-}"
REF_FILE="$HOME/.ailang/world/.last_approval_ref"
GOLDEN_FILE="scripts/world_package_ready_packet.golden.json"

die() { printf '\n  ✗ %s\n\n' "$*" >&2; exit 1; }

# has_tty — ATTEMPT THE OPEN, never `[ -r /dev/tty ]`.
#
# Measured 2026-09-21, in the agent shell this script was written in:
# `[ -r /dev/tty ]` returns TRUE while opening it fails with ENXIO. The test
# operator stats the path; the fence opens the device. A check that passes where
# the thing it predicts fails is worse than no check — it is the exact class
# cmd/world-publish/tty.go documents (a naive isatty admits `< /dev/null`), and
# the first draft of this script shipped it.
has_tty() { ( : < /dev/tty ) 2>/dev/null; }
ok()  { printf '  ✓ %s\n' "$*"; }

build_bin() {
  if [ -n "$BIN" ] && [ -x "$BIN" ]; then return 0; fi
  BIN="$(mktemp -d)/world-publish"
  # A BINARY, never `go run`: on go1.25.6 `go run` exits 1 for a child that
  # exited 3, and the STOP contract IS an exit code.
  ( cd "$REPO_ROOT" && go build -o "$BIN" ./cmd/world-publish )
}

preflight() {
  printf '\nSM.D preflight — no writes, no public request\n\n'
  [ -d "$(dirname "$STORE")" ] || { mkdir -p "$(dirname "$STORE")"; ok "created $(dirname "$STORE")"; }
  [ -d "$(dirname "$STORE")" ] && ok "store parent exists: $(dirname "$STORE")"

  [ -f "$CREDENTIAL" ] || die "credential file missing: $CREDENTIAL
    Write it WITHOUT echoing the key, e.g.
      mkdir -p \"\$(dirname \"$CREDENTIAL\")\"
      ( umask 077; printf '%s' \"\$AILANG_REGISTRY_API_KEY\" > \"$CREDENTIAL\" )"
  mode=$(stat -f '%Lp' "$CREDENTIAL")
  [ "$mode" = "600" ] || die "credential must be mode 0600, is $mode: $CREDENTIAL"
  ok "credential present, mode 0600 (contents never read by this script)"

  case "$CREDENTIAL" in
    "$REPO_ROOT"/*) die "credential is INSIDE the working tree: $CREDENTIAL" ;;
  esac
  ok "credential is outside the working tree"

  [ -x "$COMPILER" ] || die "pinned compiler missing: $COMPILER"
  # DERIVE the expected version from the golden; do not restate it.
  #
  # This script hardcoded "AILANG v0.30.0". When the pin moved to v0.41.0 the
  # other eleven files were updated and this one was not, so the helper refused
  # the very compiler the gate had just accepted — a second source of truth
  # drifting from the first, which is the defect this whole session kept finding
  # in other people's code. The golden's compilerVersion IS the authority on what
  # projected the packet, so read it there and there is nothing left to drift.
  want_compiler="$(sed -n 's/.*"compilerVersion":"\([^"]*\)".*/\1/p' "$REPO_ROOT/$GOLDEN_FILE")"
  [ -n "$want_compiler" ] || die "could not read compilerVersion from $GOLDEN_FILE"
  ver=$("$COMPILER" --version 2>/dev/null | head -1 || true)
  case "$ver" in
    "$want_compiler") ok "pinned compiler: $ver (matches the golden)" ;;
    *) die "compiler is '$ver', expected '$want_compiler' — the ready packet was projected with that exact build.
    Expected version READ FROM $GOLDEN_FILE, not hardcoded here." ;;
  esac

  if [ -n "${AILANG_REGISTRY_API_KEY:-}" ]; then
    printf '  ! AILANG_REGISTRY_API_KEY is set in this shell.\n'
    printf '    The production handler refuses to construct when it is. This script\n'
    printf '    strips it from the CHILD (env -u), so you need not unset it here —\n'
    printf '    but do not pass it to world-publish by hand.\n'
  else
    ok "AILANG_REGISTRY_API_KEY not set"
  fi

  build_bin; ok "world-publish built: $BIN"

  printf '\n  packet drift check:\n'
  ( cd "$REPO_ROOT" && env -u AILANG_REGISTRY_API_KEY "$BIN" packet ) | sed 's/^/    /'

  printf '\n  indeterminate receipts (must be 0 before a first publish):\n'
  ( cd "$REPO_ROOT" && env -u AILANG_REGISTRY_API_KEY "$BIN" reconcile \
      --store "$STORE" --registry-origin "$REGISTRY" ) | sed 's/^/    /'

  printf '\n  terminal check: '
  if has_tty; then printf 'controlling terminal OPENS — approve can run from here\n'
  else printf 'NO controlling terminal (open failed) — approve WILL refuse from here.\n                  Run it from a real shell, not an agent or editor.\n'; fi
  printf '\nPreflight done. Next:  %s approve\n\n' "$0"
}

do_approve() {
  build_bin
  has_tty || die "no controlling terminal (opening /dev/tty failed).
    Run this from a real shell — not an agent, editor task, CI job or anything
    with a redirect. This is cmd/world-publish/tty.go's fence working as designed."
  printf '\nMinting a ONE-SHOT approval. You will be asked to TYPE the phrase.\n'
  printf 'stdin is connected to /dev/tty so the fence sees the real terminal.\n\n'
  # < /dev/tty is the fence's REQUIREMENT, not a bypass: it refuses precisely
  # because stdin was not the controlling terminal. The phrase is still typed.
  ( cd "$REPO_ROOT" && env -u AILANG_REGISTRY_API_KEY "$BIN" approve \
      --store "$STORE" --registry-origin "$REGISTRY" \
      --now 1 --expires 1000 --requester "$USER" --decided-by "$USER" \
      < /dev/tty ) | tee /tmp/sm-d-approve.$$.out
  ref=$(grep -oE 'sha256:[0-9a-f]{64}' /tmp/sm-d-approve.$$.out | tail -1 || true)
  rm -f /tmp/sm-d-approve.$$.out
  [ -n "$ref" ] || die "could not read an approval ref from the output — export WORLD_APPROVAL by hand"
  ( umask 077; printf '%s' "$ref" > "$REF_FILE" )
  printf '\n  ✓ approval ref saved: %s\n' "$ref"
  printf '\nNext:  %s dry-run      (rehearses every fence, makes NO request)\n\n' "$0"
}

# valid_ref — a minted ref is sha256: plus 64 LOWERCASE hex. The fence checks
# this too, but by then the operator has already typed YES to an irreversible
# write and is reading a STOP instead of a result.
valid_ref() {
  case "$1" in
    sha256:*) : ;;
    *) return 1 ;;
  esac
  printf '%s' "${1#sha256:}" | grep -qE '^[0-9a-f]{64}$'
}

approval_ref() {
  # VALIDATE an exported WORLD_APPROVAL before trusting it over the saved ref.
  #
  # Measured 2026-09-21: the hand-written instructions that preceded this script
  # said `export WORLD_APPROVAL="sha256:<the ref that was printed>"`, the
  # operator pasted the block verbatim — entirely reasonably — and the literal
  # placeholder then SHADOWED the good ref this script had saved. The result was
  # a STOP at fence=approval reason=malformed, after the YES prompt, on a run
  # that was otherwise ready. Failing here, by name, costs one command instead.
  if [ -n "${WORLD_APPROVAL:-}" ]; then
    if valid_ref "$WORLD_APPROVAL"; then printf '%s' "$WORLD_APPROVAL"; return; fi
    saved=""
    [ -f "$REF_FILE" ] && saved="$(cat "$REF_FILE")"
    if [ -n "$saved" ] && valid_ref "$saved"; then
      die "WORLD_APPROVAL is set but is not a valid ref:
      $WORLD_APPROVAL
    A VALID minted ref IS saved at $REF_FILE:
      $saved
    The exported value wins by default, so clear it and re-run:
      unset WORLD_APPROVAL && $0 ${mode:+${mode#--}}"
    fi
    die "WORLD_APPROVAL is set but is not a valid ref: $WORLD_APPROVAL
    Expected sha256: followed by 64 lowercase hex characters."
  fi
  [ -f "$REF_FILE" ] || die "no approval ref; run '$0 approve' first, or export a valid WORLD_APPROVAL"
  ref="$(cat "$REF_FILE")"
  valid_ref "$ref" || die "saved ref at $REF_FILE is malformed: $ref"
  printf '%s' "$ref"
}

# print_live_command — `live` PRINTS, it does not invoke.
#
# host/runbook's AC30 guard forbids any automated surface from reaching the
# publish entrypoint, and its reasoning is right: "a CI step that TRIES is a
# design error, not a near miss". The tty fence would refuse automation anyway;
# the guard is about what the repository is SHAPED to allow.
#
# This helper therefore stops at the boundary. It prepares everything, rehearses
# everything, and hands the irreversible command to a human to paste — which is
# also the runbook's own principle that minting and spending are two human acts.
# The printed command carries `< /dev/tty` because the fence compares stdin to
# the controlling terminal with os.SameFile, and an interactive shell's stdin is
# the pty, a different file.
print_live_command() {
  ref="$(approval_ref)"
  cat <<BANNER

  ⚠ THE IRREVERSIBLE PUBLIC WRITE IS NOT RUN BY THIS SCRIPT, BY DESIGN.
    world/core@0.1.0 · the registry is immutable (409 on re-publish).

    Everything above is ready. Paste this, in this terminal:

BANNER
  printf '      %s publish --live \\
' "$BIN"
  printf '        --store %s \\
' "$STORE"
  printf '        --registry-origin %s \\
' "$REGISTRY"
  printf '        --publisher %s \\
' "$COMPILER"
  printf '        --credential-file %s \\
' "$CREDENTIAL"
  printf '        --approval-ref %s \\
' "$ref"
  printf '        --now 2 --expires 1000 < /dev/tty
'
  cat <<'BANNER'

    Read the outcome:
      exit 0  PUBLISHED      done, no reconciliation
      exit 1  FAILED         over; the cause now prints beneath the FAILED line
      exit 1  INDETERMINATE  DO NOT RETRY — run this helper's `reconcile`
BANNER
}

do_publish() {
  mode="$1"   # --dry-run only; --live prints instead (see print_live_command)
  build_bin
  has_tty || die "no controlling terminal (opening /dev/tty failed).
    publish runs the same tty fence as approve. Run this from a real shell."
  ref="$(approval_ref)"
  set +e
  # < /dev/tty on BOTH publish modes, not just approve.
  #
  # `publish` runs the SAME tty fence as `approve` — "runs every fence" in the
  # runbook means every fence. The first version of this script wired the
  # terminal for approve only, so the mint succeeded and the rehearse died at
  # `fence=tty reason=stdin-is-not-the-controlling-terminal`.
  #
  # And the fence needs the redirect even at a REAL terminal: it compares stdin
  # to an opened /dev/tty with os.SameFile, and an interactive shell's stdin is
  # the pty (/dev/ttysNNN), which is a different file from /dev/tty. So a
  # genuine attended operator fails the SameFile check unless stdin IS /dev/tty.
  # Measured 2026-09-21 on a normal zsh session. Conservative and fail-closed,
  # but it is the reason this helper has to exist at all.
  ( cd "$REPO_ROOT" && env -u AILANG_REGISTRY_API_KEY "$BIN" publish "$mode" \
      --store "$STORE" --registry-origin "$REGISTRY" --publisher "$COMPILER" \
      --credential-file "$CREDENTIAL" --approval-ref "$ref" --now 2 --expires 1000 \
      < /dev/tty )
  rc=$?
  set -e
  # MODE-AWARE. The first version printed "PUBLISHED — done" on a successful
  # DRY-RUN, directly under the child's own "REHEARSAL — no request of any kind
  # was made". A wrapper that tells the operator they published when they
  # rehearsed is worse than no wrapper, and on an irreversible procedure it is
  # the single most dangerous line in the script. Measured 2026-09-21.
  printf '\n  exit=%d  ' "$rc"
  if [ "$mode" = "--dry-run" ]; then
    case "$rc" in
      0) printf '(REHEARSED — every fence passed, NOTHING was sent. Next: %s live)\n' "$0" ;;
      3) printf '(STOP — a fence refused during the rehearsal; it names itself above)\n' ;;
      *) printf '(rehearsal did not complete — read the line above)\n' ;;
    esac
  else
    case "$rc" in
      0) printf '(PUBLISHED — the public write landed. Done, no reconciliation)\n' ;;
      1) printf '(FAILED or INDETERMINATE — read the line above. If INDETERMINATE, DO NOT RETRY; run "%s reconcile")\n' "$0" ;;
      3) printf '(STOP — a fence refused BEFORE any network reach; nothing was sent)\n' ;;
      *) printf '\n' ;;
    esac
  fi
  return "$rc"
}

case "${1:-}" in
  preflight) preflight ;;
  approve)   do_approve ;;
  dry-run)   do_publish --dry-run ;;
  live)      build_bin; print_live_command ;;
  reconcile) build_bin; ( cd "$REPO_ROOT" && env -u AILANG_REGISTRY_API_KEY "$BIN" reconcile --store "$STORE" --registry-origin "$REGISTRY" --probe ) ;;
  *) cat <<USAGE
SM.D attended publish helper — one human act per invocation, by design.

  $0 preflight    checks + packet drift + receipts. No writes.
  $0 approve      mints the ONE-SHOT approval. You type the phrase at /dev/tty.
  $0 dry-run      rehearses every fence. Makes no request.
  $0 live         PRINTS the irreversible command for you to paste. Runs nothing.
  $0 reconcile    read-only: resolve an INDETERMINATE attempt.

Runbook: docs/SELF_MOD_PUBLISH.md (it wins on any disagreement).
USAGE
     exit 2 ;;
esac
