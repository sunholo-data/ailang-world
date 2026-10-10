#!/bin/bash
# Run ONLY after the fresh baseline is banked; user types the publication phrase.
set -euo pipefail
cd /private/tmp/ailang-world-row174-20261010
ROW174_RIG=/Users/voightkampff/dev/sunholo-data/ailang-world/.ailang/state/row174-attended-2026-10-10/rig
if ! test -t 0; then
  echo 'Run this yourself in a terminal: docs/QUICKSTART.md section 0 requires typed authority confirmation.' >&2
  exit 1
fi
if curl --connect-timeout 2 --max-time 3 -sf http://127.0.0.1:7646/v1/health >/dev/null; then
  echo 'Scratch daemon is serving on port 7646. Stop it before publishing.' >&2
  exit 1
fi
env -u AILANG_REGISTRY_API_KEY "$ROW174_RIG/bin/world-publish" transitions \
  --store "$ROW174_RIG/store/world.db" \
  --manifest packages/se-tools/transitions.json \
  --ailang-bin /Users/voightkampff/.pinned-ailang/ailang
echo 'Candidate published. Return to Codex to run candidate smoke and the no-change repeat.'
