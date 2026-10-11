from pathlib import Path
import datetime as dt
import hashlib
import json
import shutil
import sys

iteration = int(sys.argv[1])
label = sys.argv[2]
runner_rc = int(sys.argv[3])
candidate = Path('/private/tmp/ailang-world-row174-20261010')
rig = Path('/Users/voightkampff/dev/sunholo-data/ailang-world/.ailang/state/row174-attended-2026-10-10/rig')
bank = candidate/'design_docs/verification/world-row174-runtime-2026-10-11'
ev = candidate/f'design_docs/verification/world-floor-m4-2026-10-06/iter-{iteration}'
record = json.loads((ev/'iteration.json').read_text())
assert record['iteration'] == iteration
assert record['rows_reused_from'] == {}
shutil.copytree(rig/f'm4/transcripts/iter-{iteration}', bank/f'iter-{iteration}-transcripts')
shutil.copy2(f'/private/tmp/row174-codex-{label}-iter{iteration}.log', bank/f'iter-{iteration}-runner.log')
rows = {}
for arm in ['shell','world']:
    p = ev/f'runs/codex/{arm}/r1.jsonl'
    rows[arm] = [json.loads(l) for l in p.read_text().splitlines() if l.strip()] if p.exists() else []
valid = record['status'] == 'complete' and runner_rc == 0
valid = valid and all(len(rows[a]) == 6 and all(r['class'] == 'PASS' and not r['attestation'] for r in rows[a]) for a in rows)
valid = valid and not record['native_write_detected'] and all(
    r.get('provenance',{}).get('solution_provenance') == 'world' and
    r.get('provenance',{}).get('native_write_detected') is False for r in rows['world'])
summary = record['summary'].get('codex', {})
meets_target = valid and summary.get('O') is not None and summary['O'] <= .20 and summary.get('delta', -1) >= -.02
result = {'recorded_at_utc':dt.datetime.now(dt.timezone.utc).isoformat(), 'iteration':iteration,
          'label':label, 'runner_exit_code':runner_rc, 'status':record['status'], 'abort':record['abort'],
          'scope':'operator-requested Codex-only diagnostic; Claude and full four-cell M2 acceptance pending',
          'summary':summary, 'canary':record['canary'], 'row153':record['row153'],
          'codex_tokens':record['codex_tokens'], 'wall_s':record['wall_s'],
          'all_tasks_pass_and_provenance_valid':valid, 'codex_diagnostic_target_met':meets_target,
          'model':'gpt-6-luna', 'product_commit':record['commit'],
          'config_digest_smoke':record['config_digest_smoke'], 'rows_reused_from':{},
          'full_M2_acceptance':False, 'candidate_publication':'User typed 2026-10-11; local revision2',
          'task_outcomes':{a:[{k:r.get(k) for k in ['task','class','attempt','wall_ms','error_category','attestation']} for r in rows[a]] for a in rows}}
(bank/f'codex-{label}-iter{iteration}-result.json').write_text(json.dumps(result,indent=2,sort_keys=True)+'\n')
sys.path.insert(0,str(candidate/'scripts/floor'))
import tokens
known = tokens.load_tokens([rig/'private/fl-tune-cc.session',rig/'private/fl-tune-cx.session'])
for p in [bank, ev]:
    findings = tokens.scan_dir(str(p),known)
    assert not findings, findings
files = sorted(p for p in bank.rglob('*') if p.is_file() and p.name != 'evidence-index.json')
(bank/'evidence-index.json').write_text(json.dumps({'files':[
    {'path':str(p.relative_to(candidate)), 'bytes':p.stat().st_size,
     'sha256':hashlib.sha256(p.read_bytes()).hexdigest()} for p in files]},indent=2)+'\n')
print(json.dumps({'iteration':iteration, 'status':record['status'], 'summary':summary,
                  'all_tasks_pass_and_provenance_valid':valid, 'codex_diagnostic_target_met':meets_target,
                  'token_scan':'PASS', 'full_M2_acceptance':False},indent=2))
