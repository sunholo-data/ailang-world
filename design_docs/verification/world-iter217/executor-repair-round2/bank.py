from pathlib import Path
import json,hashlib,subprocess
src=Path('/tmp/world-iter217-repair-round2');out=Path('design_docs/verification/world-iter217/executor-repair-round2');out.mkdir(exist_ok=True)
rows=json.loads((src/'production-before.json').read_text())
for row in rows:assert hashlib.sha256(Path(row['path']).read_bytes()).hexdigest()==row['sha256'],row['path']
(src/'production-after.json').write_text(json.dumps(rows,indent=2)+'\n')
assert Path('host/projection/mcp_bounds_test.go').read_bytes()==(src/'repaired.go').read_bytes()
results={p.stem:json.loads(p.read_text()) for p in src.glob('*.json') if 'rc' in json.loads(p.read_text())}
assert all(x['rc']==x['expected'] for x in results.values())
assert all(name in results for name in ['final-normal','final-race','cb-background-compile','cb-background-selected','cb-background-restored','aggregate-removed-compile','aggregate-removed-selected','aggregate-removed-restored'])
manifest={'base':'362f361d00a2eddd15ac45c6a4784fbdd4b4e738','scope':'test-only phase-isolation plus truthful design/plan status; no production/module/kernel/contract changes','checks':results,'production_hash_scope_count':len(rows),'return_signal_semantics':'Runner close/counter records cooperative return path before actual function return; successful bounded followup through same one-slot hostcall Runner demonstrates slot release','mutations':'actual CB-BACKGROUND and aggregate wrapper removal; selected test only; compile0/semantic1/exact byte restore0','pending':['final normal AIL then Go profiles','fresh Round3 independent judge','remote CI'],'native_agent_tokens':'not reported; not zero'}
(src/'repair-manifest.json').write_text(json.dumps(manifest,indent=2)+'\n')
raw=[]
for p in sorted(src.iterdir()):
 if not p.is_file():continue
 data=p.read_bytes();data.decode();assert b'\0' not in data
 name='repaired-source.txt' if p.name=='repaired.go' else p.name
 (out/name).write_bytes(data);raw.append({'path':name,'bytes':len(data),'sha256':hashlib.sha256(data).hexdigest()})
(out/'raw-file-manifest.json').write_text(json.dumps(raw,indent=2)+'\n')
(out/'README.md').write_text('''# Gate3 round2 repair evidence\n\nNormal and race named tests plus actual production semantic controls use pinned AILANG_BIN and no GODEBUG override. All commands have finite external bounds. Initial-only resolver delay and aggregate stimulus are separated from same-handler one-slot recovery. Return-path counters precede actual function return; second successful execution proves released capacity. No timing error is swallowed or retried. Production/module hashes match the controller base.\n\nRound2 original full Go failure remains preserved in root-owned evidence; resolver calibration does not establish its unique cause. Final normal profiles, fresh Round3 and remote CI remain pending. D passes remains null; no product acceptance/release claim.\n''')
p=Path('design_docs/verification/world-iter217/sprint_w-mcp-dispatch-projection-tranche-a.json');plan=json.loads(p.read_text());plan['repair_round2']={'scope':manifest['scope'],'evidence':str(out),'checks':results,'status':'local repair verified; final normal profiles/fresh Round3/remote CI pending','native_agent_tokens':'not reported'};plan['current_design_sha256']=hashlib.sha256(Path('design_docs/planned/w-mcp-dispatch-projection.md').read_bytes()).hexdigest();plan['current_design_sha256_note']='Current canonical metadata/status includes round2 local fixture repair; immutable original design authorization/snapshot hashes retained; contract unchanged.'
for f in plan['features']:
 if f['id'].startswith('M108D'):assert f['passes'] is None
p.write_text(json.dumps(plan,indent=2)+'\n');print('banked',len(raw),'text files; production hash count',len(rows))
