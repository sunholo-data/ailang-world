#!/usr/bin/env bash
set -u

if [ "$#" -eq 1 ] && [ "$1" = "--version" ]; then
  [ -n "${AILANG_SHIM_VERSION_LINE:-}" ] || exit 97
  printf '%s\n' "$AILANG_SHIM_VERSION_LINE"
  printf '%s\n' "AILANG version shim fixture"
  exit 0
fi

[ -n "${AILANG_SHIM_DELEGATE:-}" ] || exit 97
exec "$AILANG_SHIM_DELEGATE" "$@"
