#!/bin/bash
# Run ONLY after Codex-only iteration8 baseline is complete and committed; user types the publication phrase.
set -euo pipefail
cd /private/tmp/ailang-world-row174-20261010
ROW174_RIG=/Users/voightkampff/dev/sunholo-data/ailang-world/.ailang/state/row174-attended-2026-10-10/rig
if ! test -t 0; then
  echo 'Run this yourself in a terminal: docs/QUICKSTART.md section 0 requires typed authority confirmation.' >&2
  exit 1
fi
for ROW174_PORT in 7644 7646; do
  if curl --connect-timeout 2 --max-time 3 -sf http://127.0.0.1:$ROW174_PORT/v1/health >/dev/null; then
    echo "A daemon is serving on port $ROW174_PORT. Stop it before publishing." >&2
    exit 1
  fi
done
python3 - <<'CHECK'
from pathlib import Path
import json, subprocess
p=Path('design_docs/verification/world-floor-m4-2026-10-06/iter-8/iteration.json')
subprocess.run(['git','ls-files','--error-unmatch',str(p)],check=True,capture_output=True)
subprocess.run(['git','diff','--exit-code','HEAD','--',str(p)],check=True,capture_output=True)
r=json.loads(p.read_text())
assert r['status']=='complete' and r['summary']['codex']['T']==6
assert len(r['canary'])==2 and all(c['class']=='PASS' for c in r['canary'])
for arm in ['shell','world']:
    q=p.parent/'runs/codex'/arm/'r1.jsonl'
    rows=[json.loads(l) for l in q.read_text().splitlines() if l.strip()]
    assert len(rows)==6 and not any(x.get('attestation') for x in rows)
print('Committed Codex-only baseline verified; full Claude acceptance remains pending.')
CHECK
env -u AILANG_REGISTRY_API_KEY "$ROW174_RIG/bin/world-publish" transitions \
  --store "$ROW174_RIG/store/world.db" \
  --manifest packages/se-tools/transitions.json \
  --ailang-bin /Users/voightkampff/.pinned-ailang/ailang
echo 'Candidate published. Return to Codex for the Codex-only candidate and no-change repeat; Claude validation remains pending.'
