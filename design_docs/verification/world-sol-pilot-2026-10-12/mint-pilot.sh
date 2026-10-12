#!/bin/bash
# Mark runs this in a real terminal; this script never answers a prompt itself (QUICKSTART section 0).
# It mints ONE fresh codex World session for the Sol pilot on the existing row-174 rig.
# The store already holds the candidate manifest (revision 2), so there is no publish step.
set -euo pipefail
RIG=/Users/voightkampff/dev/sunholo-data/ailang-world/.ailang/state/row174-attended-2026-10-10/rig
EP=fl-tune-cx
if ! test -t 0; then
  echo 'Run this yourself in a terminal: QUICKSTART section 0 requires typed authority confirmation.' >&2; exit 1
fi
if curl --connect-timeout 2 --max-time 3 -sf http://127.0.0.1:7644/v1/health >/dev/null; then
  echo 'A daemon is serving on 7644 (it holds the store writer lock). Stop it first.' >&2; exit 1
fi
OLD="$RIG/private/$EP.session"
if test -e "$OLD"; then
  # Expired 12h token from 2026-10-11 20:30. Archived, never deleted or overwritten.
  if test -n "$(find "$OLD" -mmin -720)"; then
    echo "$OLD is younger than 12h and may still be live; stop for review." >&2; exit 1
  fi
  mkdir -p "$RIG/private/expired-2026-10-11"
  test ! -e "$RIG/private/expired-2026-10-11/$EP.session" || { echo 'archive slot occupied; stop for review' >&2; exit 1; }
  mv "$OLD" "$RIG/private/expired-2026-10-11/$EP.session"
  echo "archived expired $EP.session"
fi
env -u AILANG_REGISTRY_API_KEY "$RIG/bin/ailang-worldd" session new "$EP" \
  --db "$RIG/store/world.db" --workspace-root "$RIG/ws" \
  --grant Workspace.Read=worktree:5000 --grant Workspace.Write=worktree:5000 \
  --grant Ailang.Check=worktree:5000 --grant Ailang.Run=worktree:5000 \
  --grant Ailang.RunEnv=worktree:5000 --grant Ailang.RunNet=worktree:5000 \
  --grant Ailang.Discover=worktree:5000 --grant Ailang.CLI=worktree:5000 \
  --ttl 43200 --out "$RIG/private/$EP.session"
ls -l "$RIG/private/$EP.session"
echo 'Pilot session ready. Tell Claude to run the Sol pilot; no inference has run yet.'
