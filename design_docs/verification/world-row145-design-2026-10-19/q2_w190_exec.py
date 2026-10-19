"""Q2 (row 145 design): episode w190's workspace-exec calls, from its Claude Code transcript.
Counts calls, results that timed out (the 9 s cap) and results denied a loopback bind."""
import json
import sys

T = sys.argv[1] if len(sys.argv) > 1 else \
    '/Users/voightkampff/.ailang/state/trial-2026-10-16/evidence/ailang-world/transcripts/w190.jsonl'
calls, res = {}, {}
for line in open(T, encoding='utf-8'):
    m = (json.loads(line).get('message') or {})
    for c in m.get('content') or []:
        if not isinstance(c, dict):
            continue
        if c.get('type') == 'tool_use' and c.get('name', '').endswith('workspace-exec'):
            calls[c['id']] = c.get('input')
        elif c.get('type') == 'tool_result' and c.get('tool_use_id') in calls:
            text = json.dumps(c.get('content'))
            res[c['tool_use_id']] = ('\\"timed_out\\":true' in text, 'bind: operation not permitted' in text)
print('workspace-exec calls', len(calls), 'results', len(res))
print('timed_out', sum(r[0] for r in res.values()), 'loopback bind denied', sum(r[1] for r in res.values()))
