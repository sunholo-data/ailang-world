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

REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
STORE="${WORLD_STORE:-$HOME/.ailang/world/world.db}"
REGISTRY="${WORLD_REGISTRY:-https://storage.googleapis.com/ailang-registry}"
CREDENTIAL="${WORLD_CREDENTIAL:-$HOME/.config/ailang/registry.key}"
COMPILER="${WORLD_COMPILER:-$HOME/.pinned-ailang/ailang}"
BIN="${WORLD_BIN:-}"
REF_FILE="$HOME/.ailang/world/.last_approval_ref"

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
  ver=$("$COMPILER" --version 2>/dev/null | head -1 || true)
  case "$ver" in
    "AILANG v0.30.0") ok "pinned compiler: $ver" ;;
    *) die "compiler is '$ver', expected 'AILANG v0.30.0' — the ready packet was projected with that exact build" ;;
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

approval_ref() {
  if [ -n "${WORLD_APPROVAL:-}" ]; then printf '%s' "$WORLD_APPROVAL"; return; fi
  [ -f "$REF_FILE" ] || die "no approval ref; run '$0 approve' first, or export WORLD_APPROVAL"
  cat "$REF_FILE"
}

do_publish() {
  mode="$1"   # --dry-run | --live
  build_bin
  ref="$(approval_ref)"
  if [ "$mode" = "--live" ]; then
    printf '\n  ⚠ THIS IS THE IRREVERSIBLE PUBLIC WRITE.\n'
    printf '    world/core@0.1.0 · the registry is immutable (409 on re-publish).\n'
    printf '    On INDETERMINATE: DO NOT RETRY — run "%s reconcile".\n\n' "$0"
    printf '    Type YES to proceed: '
    read -r go < /dev/tty
    [ "$go" = "YES" ] || die "aborted (nothing was sent)"
  fi
  set +e
  ( cd "$REPO_ROOT" && env -u AILANG_REGISTRY_API_KEY "$BIN" publish "$mode" \
      --store "$STORE" --registry-origin "$REGISTRY" --publisher "$COMPILER" \
      --credential-file "$CREDENTIAL" --approval-ref "$ref" --now 2 --expires 1000 )
  rc=$?
  set -e
  printf '\n  exit=%d  ' "$rc"
  case "$rc" in
    0) printf '(PUBLISHED — done, no reconciliation)\n' ;;
    1) printf '(FAILED or INDETERMINATE — read the line above. If INDETERMINATE, DO NOT RETRY)\n' ;;
    3) printf '(STOP — a fence refused; it names itself above)\n' ;;
    *) printf '\n' ;;
  esac
  return "$rc"
}

case "${1:-}" in
  preflight) preflight ;;
  approve)   do_approve ;;
  dry-run)   do_publish --dry-run ;;
  live)      do_publish --live ;;
  reconcile) build_bin; ( cd "$REPO_ROOT" && env -u AILANG_REGISTRY_API_KEY "$BIN" reconcile --store "$STORE" --registry-origin "$REGISTRY" --probe ) ;;
  *) cat <<USAGE
SM.D attended publish helper — one human act per invocation, by design.

  $0 preflight    checks + packet drift + receipts. No writes.
  $0 approve      mints the ONE-SHOT approval. You type the phrase at /dev/tty.
  $0 dry-run      rehearses every fence. Makes no request.
  $0 live         the irreversible public write. Asks YES first.
  $0 reconcile    read-only: resolve an INDETERMINATE attempt.

Runbook: docs/SELF_MOD_PUBLISH.md (it wins on any disagreement).
USAGE
     exit 2 ;;
esac
