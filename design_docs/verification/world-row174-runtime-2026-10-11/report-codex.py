from pathlib import Path
import datetime as dt
import hashlib
import json
import statistics as stats

root = Path('/private/tmp/ailang-world-row174-20261010')
bank = root/'design_docs/verification/world-row174-runtime-2026-10-11'
runs = []
for i, label in [(8,'Original baseline'),(9,'Candidate'),(10,'Unchanged candidate repeat')]:
    ev = root/f'design_docs/verification/world-floor-m4-2026-10-06/iter-{i}'
    r = json.loads((ev/'iteration.json').read_text())
    rows = {}
    for arm in ['shell','world']:
        rows[arm] = [json.loads(l) for l in (ev/f'runs/codex/{arm}/r1.jsonl').read_text().splitlines() if l.strip()]
        assert len(rows[arm]) == 6 and {x['task'] for x in rows[arm]} == set(r['tasks'])
        assert all(x.get('attempt',1) == 1 for x in rows[arm])
    s = r['summary']['codex']
    by = {a:{x['task']:x for x in rows[a]} for a in rows}
    ratios = {t:(by['world'][t]['wall_ms']-by['shell'][t]['wall_ms'])/by['shell'][t]['wall_ms'] for t in r['tasks']}
    median = stats.median(ratios.values())
    assert round(median,4) == s['O']
    assert s['delta'] == (s['pass_W']-s['pass_S'])/6
    faults = s['harness_faults']['shell']+s['harness_faults']['world']
    passed = r['status']=='complete' and s['O']<=.20 and s['delta']>=-.02 and faults==0 and s['pass_W']==6 and s['pass_S']==6
    runs.append({'iteration':i,'label':label,'status':r['status'],'summary':s,
                 'diagnostic_target_met':passed,'median_per_task_O_exact':median,
                 'task_overhead':ratios,'wall_s':r['wall_s'],'codex_tokens':r['codex_tokens'],
                 'cell_medians':{a:{'wall_ms':stats.median(x['wall_ms'] for x in rows[a]),
                    'input_tokens':stats.median(x['input_tokens'] for x in rows[a])} for a in rows}})
assert len({json.loads((root/f'design_docs/verification/world-floor-m4-2026-10-06/iter-{i}/iteration.json').read_text())['config_digest_smoke'] for i in [8,9,10]})==1
assert all(x['status']=='complete' for x in runs)
base = json.loads((bank/'codex-baseline-iter8-preflight/runtime.json').read_text())
cand = json.loads((bank/'codex-candidate-iter9-preflight/runtime.json').read_text())
repeat = json.loads((bank/'codex-repeat-iter10-preflight/runtime.json').read_text())
assert base['versions']==cand['versions']==repeat['versions']
for key in ['daemon','publisher','interpreter','tool','runner','floor_config']:
    assert base['sha256'][key]==cand['sha256'][key]==repeat['sha256'][key]
assert cand['sha256']==repeat['sha256'] and cand['metadata']==repeat['metadata']
decision = all(x['diagnostic_target_met'] for x in runs[1:])
record = {'recorded_at_utc':dt.datetime.now(dt.timezone.utc).isoformat(),
          'scope':'Codex-only gpt-6-luna N1 diagnostic; full four-cell row174 M2 pending',
          'runs':runs, 'candidate_and_repeat_meet_codex_target':decision,
          'full_M2_accepted':False, 'Claude':'unrun due observed weekly limit',
          'product_landed':False, 'historical_FINAL':'unchanged failed Sol O+26.57598%',
          'prior_aborts':[{'iteration':6,'cause':'Claude seven-day rate limit'},{'iteration':7,'cause':'daemon7646/agent7644 endpoint mismatch'}],
          'interpretation':'N1 and visible run-to-run shell variation do not establish causal efficacy or predict a new FINAL. No pooling or favorable selection.'}
(bank/'codex-comparison.json').write_text(json.dumps(record,indent=2,sort_keys=True)+'\n')
lines = ['# Row174 Codex-only smoke comparison — 2026-10-11', '',
         '**'+('Both candidate runs meet' if decision else 'Candidate runs do not both meet')+' the Codex diagnostic target. Full M2 remains pending.**', '',
         '| Run | Shell pass | World pass | Median per-task World overhead | Delta | Harness faults |',
         '|---|---:|---:|---:|---:|---:|']
for x in runs:
    s=x['summary']
    lines.append(f"| {x['label']} (iter{x['iteration']}) | {s['pass_S']}/6 | {s['pass_W']}/6 | {s['O']*100:+.2f}% | {s['delta']*100:+.2f}pp | {sum(s['harness_faults'].values())} |")
lines += ['', 'Each run uses unchanged gpt-6-luna, fixed six-task smoke, N1, fizzbuzz canaries,',
          'teaching/grader/config/thresholds/current clients; both arms are fresh with no row reuse.',
          'O is the median of paired per-task ratios, independently recomputed from raw rows;',
          'it is not the ratio of the cell medians. Every recorded attempt is retained.', '',
          'Canonical exact-eight descriptors shrank13947→10301 bytes; prose7430→3794,',
          'wire14054→10368. Candidate/repeat metadata is identical. Actual daemon/publisher/',
          'interpreter/tool/runner/config hashes and CLI versions match all three runs.', '',
          'The original baseline already meets the cheap-model target. Shell timings vary,',
          'so a favorable difference does not identify descriptor prose as the cause of faster',
          'runtime. N1 results do not predict the historical Sol FINAL or establish statistical',
          'improvement. The historical Sol+26.6% failed FINAL remains unchanged.', '',
          'Iteration6 aborted on Claude weekly quota; iteration7 retained six passing shell',
          'tasks and a failed World canary due port mismatch. Both aborted rows remain in the',
          'ledger and are excluded from this complete paired comparison, never discarded.', '',
          'Claude regression and the design’s fresh all-four-cell comparison/repeat remain',
          'pending. Product remains unmerged; implementation PASS86 remains scoped to M1.',
          'New FINAL separately needs product landing, a new committed preregistration and',
          'human-typed sessions. The daemon is stopped after each runner invocation.', '',
          '## Per-task overhead', '', '| Task | Baseline | Candidate | Unchanged repeat |',
          '|---|---:|---:|---:|']
for task in sorted(runs[0]['task_overhead']):
    lines.append('| '+task+' | '+' | '.join(f"{x['task_overhead'][task]*100:+.2f}%" for x in runs)+' |')
(bank/'CODEX-COMPARISON.md').write_text('\n'.join(lines)+'\n')
(bank/'CODEX-CONTINUE.md').write_text('''# Current attended runtime handoff

The operator-requested Codex-only baseline/candidate/unchanged repeat is complete.
See CODEX-COMPARISON.md, codex-comparison.json, existing ledger and iterations8/9/10.
Candidate was published by the user into local registry revision2 on2026-10-11.
Implementation M1 remains PASS86. Product unmerged. Claude and full M2 acceptance pending.
No further Codex-only run is planned; preserve all results.

Iteration6 quota abort and iteration7 endpoint failure remain immutable. The next unused
ledger number is11 if no other run advances it. Use default127.0.0.1:7644: run.py's override
does not reach run_agent. No harness repair was made. All product code/config stayed fixed.

Claude rate event reported reset2026-10-12 07:00 Europe/Copenhagen. Scratch sessions expire
around08:30 on that date. Check capacity, client identities, expiry and budgets afresh.
The store currently holds the candidate manifest. The original design’s all-four-cell
baseline/candidate/repeat remains pending; a fresh full baseline needs the original manifest
republished by the user under QUICKSTART section0. Prepare that exact helper before asking.
Existing mint helper refuses existing token files; do not delete/overwrite or mint automatically.
Human-typed new sessions are required if the current sessions expire.

Candidate /private/tmp/ailang-world-row174-20261010, attended/row174-codex-overhead.
Baseline /private/tmp/ailang-world-row174-baseline-20261010, original product5ca1115 plus evidence.
Rig /Users/voightkampff/dev/sunholo-data/ailang-world/.ailang/state/row174-attended-2026-10-10/rig.
Capture current binary/metadata/manifest/config/client hashes before any new inference.
Commit each attempt's ledger/evidence before another, never reuse cells, overwrite attempts or select favorable outcomes.
No model substitution, billed overage, canary waiver, favorable selection or pooling of partial
runs into full M2 acceptance. New FINAL only after landing/new preregistration/human mint.
''')
files = sorted(p for p in bank.rglob('*') if p.is_file() and p.name != 'evidence-index.json')
(bank/'evidence-index.json').write_text(json.dumps({'files':[
    {'path':str(p.relative_to(root)), 'bytes':p.stat().st_size,
     'sha256':hashlib.sha256(p.read_bytes()).hexdigest()} for p in files]},indent=2)+'\n')
print(json.dumps({'runs':[{'iteration':x['iteration'],'O':x['summary']['O'],
                          'diagnostic_target_met':x['diagnostic_target_met']} for x in runs],
                  'candidate_and_repeat_meet_codex_target':decision, 'full_M2_accepted':False},indent=2))
