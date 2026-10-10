"""Additional deadline-specific retention mutants; reuse only runner definitions, never rerun the earlier six."""
from pathlib import Path
source=(Path(__file__).resolve().parent/'fixture-mutations.py').read_text()
exec(source[:source.index("prod='host/broker/handlers.go'")],globals())
records=json.loads((BANK/'fixture-mutations.json').read_text())['mutations']
prod='host/broker/handlers.go'
def reset_after_deadline(raw,stream):
    split=raw.index(b'func runHeadTail(')
    old=b'if waitErr != nil && runCtx.Err() == context.DeadlineExceeded {'
    new=old+b'\n\t\tc.'+stream.encode()+b' = newHeadTailSink(execHeadBytes, execTailBytes, c.hardCap, nil)'
    return raw[:split]+replace_once(raw[split:],old,new)
for stream in ['stdout','stderr']:
    label='deadline-'+stream+'-reset'
    mutant(label,{prod:lambda x,stream=stream:reset_after_deadline(x,stream)},'exact partial '+stream+' Text/Bytes/SHA/untruncated')
    log=(BANK/('mutation-'+label+'.log')).read_text()
    assert len(re.findall(r'no-sleep control [123] elapsed=',log))==3,log
    assert re.search(r'deadline arm elapsed=.* kills=\[\{pgid:[0-9]+ ready:true err:<nil>\}\] Outstanding=0 baseline=0',log),log
    records[-1]['three_positive_controls_and_ready_real_deadline_kill']=True
(BANK/'fixture-mutations.json').write_text(json.dumps(dict(isolated_copy=str(ROOT),mutations=records,all_eight_expected_red=True,all_eight_restored_controls='PASS'),indent=2)+'\n')
print('PASS: two additional deadline-specific sink resets fail exact stream assertion after ready real kill; restored GREEN each; eight total isolated mutants')
