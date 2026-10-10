#!/bin/bash
# User runs this in a real terminal; prompts are never answered by this script.
set -euo pipefail
# First publish the immutable parent for a fresh same-client baseline.
cd /private/tmp/ailang-world-row174-baseline-20261010
ROW174_RIG=/Users/voightkampff/dev/sunholo-data/ailang-world/.ailang/state/row174-attended-2026-10-10/rig
ROW174_PIN=/Users/voightkampff/.pinned-ailang/ailang
if ! test -t 0; then
  echo 'Run this yourself in a terminal: docs/QUICKSTART.md section 0 requires typed authority confirmation.' >&2
  exit 1
fi
if curl --connect-timeout 2 --max-time 3 -sf http://127.0.0.1:7646/v1/health >/dev/null; then
  echo 'Scratch daemon is serving on port 7646. Stop it before publishing/minting.' >&2
  exit 1
fi
for ROW174_EPISODE in fl-tune-cc fl-tune-cx; do
  if test -e "$ROW174_RIG/private/$ROW174_EPISODE.session"; then
    echo "Session file already exists for $ROW174_EPISODE; stop for review rather than overwriting it." >&2
    exit 1
  fi
done
env -u AILANG_REGISTRY_API_KEY "$ROW174_RIG/bin/world-publish" transitions \
  --store "$ROW174_RIG/store/world.db" \
  --manifest packages/se-tools/transitions.json --ailang-bin "$ROW174_PIN"
for ROW174_EPISODE in fl-tune-cc fl-tune-cx; do
  env -u AILANG_REGISTRY_API_KEY "$ROW174_RIG/bin/ailang-worldd" session new "$ROW174_EPISODE" \
    --db "$ROW174_RIG/store/world.db" --workspace-root "$ROW174_RIG/ws" \
    --grant Workspace.Read=worktree:5000 --grant Workspace.Write=worktree:5000 \
    --grant Ailang.Check=worktree:5000 --grant Ailang.Run=worktree:5000 \
    --grant Ailang.RunEnv=worktree:5000 --grant Ailang.RunNet=worktree:5000 \
    --grant Ailang.Discover=worktree:5000 --grant Ailang.CLI=worktree:5000 \
    --ttl 43200 --out "$ROW174_RIG/private/$ROW174_EPISODE.session"
done
echo 'BASELINE publication and both sessions ready. Return to Codex to run the unchanged-model baseline.'
echo 'Candidate publication is a separate typed step after the baseline; no benchmark has run yet.'
