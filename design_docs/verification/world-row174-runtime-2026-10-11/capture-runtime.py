from pathlib import Path
import datetime as dt
import hashlib
import json
import os
import sqlite3
import subprocess
import sys

checkout = Path(sys.argv[1]).resolve()
label = sys.argv[2]
bank = Path(sys.argv[3]).resolve()
sys.path.insert(0, str(checkout / 'scripts/floor'))
import canary
import world
import arms

rig = Path('/Users/voightkampff/dev/sunholo-data/ailang-world/.ailang/state/row174-attended-2026-10-10/rig')
bank.mkdir(parents=True, exist_ok=False)
sha = lambda p: hashlib.sha256(Path(p).read_bytes()).hexdigest()
record = {'recorded_at_utc': dt.datetime.now(dt.timezone.utc).isoformat(),
          'label': label, 'checkout': str(checkout), 'rig': str(rig),
          'commit': subprocess.check_output(['git','rev-parse','HEAD'], cwd=checkout, text=True).strip(),
          'load_average': list(os.getloadavg()), 'versions': {}, 'sha256': {}, 'sessions': [],
          'serve_argv': canary.serve_argv(canary.layout(str(rig)), '127.0.0.1:7646')}
for name, cmd in {'codex':['codex','--version'], 'claude':['claude','--version'],
                 'interpreter':[canary.PIN,'--version'], 'tool':[canary.TOOL,'--version']}.items():
    record['versions'][name] = subprocess.check_output(cmd, text=True).strip()
for name, path in {'manifest':checkout/'packages/se-tools/transitions.json',
                   'floor_config':checkout/'scripts/floor/floor_config.json',
                   'runner':checkout/'scripts/floor/run.py', 'daemon':rig/'bin/ailang-worldd',
                   'publisher':rig/'bin/world-publish', 'interpreter':canary.PIN, 'tool':canary.TOOL}.items():
    record['sha256'][name] = sha(path)
db = sqlite3.connect('file:'+str(rig/'store/world.db')+'?mode=ro', uri=True)
db.row_factory = sqlite3.Row
for ep in ['fl-tune-cc','fl-tune-cx']:
    token_file = rig/'private'/f'{ep}.session'
    assert token_file.stat().st_mode & 0o077 == 0
    rows = [dict(r) for r in db.execute('SELECT credential_id,episode_id,grants_json,expires_at,created_at FROM session_credentials WHERE episode_id=?', (ep,))]
    assert len(rows) == 1, (ep, len(rows))
    r = rows[0]
    r['grants'] = json.loads(r.pop('grants_json'))
    assert len(r['grants']) == 8
    assert r['expires_at'] > int(dt.datetime.now().timestamp()) + 7200
    record['sessions'].append(r)
db.close()
L = canary.layout(str(rig))
daemon = canary.start_daemon(L, '127.0.0.1:7646')
try:
    client = world.Client('127.0.0.1:7646')
    record['health'] = client.health()
    for ep in ['fl-tune-cc','fl-tune-cx']:
        token = (rig/'private'/f'{ep}.session').read_text().strip()
        record.setdefault('preflight', {})[ep] = world.preflight(client, token)
        body = json.dumps({'jsonrpc':'2.0','id':1,'method':'tools/list'}).encode()
        status, raw, content_type = client._req('POST','/mcp/', body,
            {'Authorization':'Bearer '+token,'Content-Type':'application/json','Accept':'application/json, text/event-stream'})
        assert status == 200
        (bank/f'{ep}-tools-list.sse').write_bytes(raw)
        tools = world.parse_sse_json(raw)['result']['tools']
        manifest = json.loads((checkout/'packages/se-tools/transitions.json').read_text())
        expected = {t['id']:t for t in manifest if t['id'] != 'workspace-exec'}
        assert {t['name'] for t in tools} == set(expected)
        for t in tools:
            assert t['description'] == expected[t['name']]['description']
            assert t['inputSchema'] == expected[t['name']]['inputSchema']
        canonical = json.dumps(tools,sort_keys=True,separators=(',',':'),ensure_ascii=False).encode()
        (bank/f'{ep}-tools-list.canonical.json').write_bytes(canonical)
        record.setdefault('metadata', {})[ep] = {'wire_bytes':len(raw),'wire_sha256':hashlib.sha256(raw).hexdigest(),
            'canonical_bytes':len(canonical),'canonical_sha256':hashlib.sha256(canonical).hexdigest(),
            'description_bytes':sum(len(t['description'].encode()) for t in tools),'content_type':content_type}
finally:
    canary.stop_daemon(daemon)
(bank/'runtime.json').write_text(json.dumps(record,indent=2,sort_keys=True)+'\n')
print(json.dumps({'label':label,'commit':record['commit'],'metadata':record['metadata'],'versions':record['versions']},indent=2))
