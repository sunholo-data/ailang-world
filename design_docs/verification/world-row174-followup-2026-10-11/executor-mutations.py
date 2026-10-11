"""Four disclosed compiled mutants, in disposable copy only; exact restore and real GREEN after each."""
from pathlib import Path
import difflib,hashlib,json,re,subprocess
ROOT=Path('/private/tmp/ailang-world-exec-timeout-mutants-20261011')
BANK=Path(__file__).resolve().parent
FILES=['host/broker/handlers_exec_m2_test.go','host/broker/handlers.go']
ORIGINAL={f:(ROOT/f).read_bytes() for f in FILES}
SHAS={f:hashlib.sha256(v).hexdigest() for f,v in ORIGINAL.items()}
COMMAND=['go','test','./host/broker','-run','^TestExecHandlerTimeoutKillsTheGroup$','-v','-count=1','-timeout=600s']
records=[]
def replace(raw,old,new):
    assert raw.count(old)==1,(old,raw.count(old))
    return raw.replace(old,new)
def invoke(label):
    r=subprocess.run(COMMAND,cwd=ROOT,capture_output=True,text=True)
    log=r.stdout+r.stderr
    (BANK/(label+'.log')).write_text(log)
    assert '--- SKIP:' not in log,log
    return r.returncode,log

def mutant(label,file,fn,identity):
    try:
        changed=fn(ORIGINAL[file]); (ROOT/file).write_bytes(changed)
        diff=''.join(difflib.unified_diff(ORIGINAL[file].decode().splitlines(True),changed.decode().splitlines(True),fromfile=file,tofile=file+' (isolated mutant)'))
        (BANK/('mutation-'+label+'.diff')).write_text(diff)
        sha=hashlib.sha256(changed).hexdigest()
        rc,log=invoke('mutation-'+label)
        assert rc!=0 and identity in log,(label,rc,log)
        assert len(re.findall(r'short control [123] complete=',log))==3,log
        assert 'Outstanding=0 baseline=0' in log and 'cleanup actual reap baseline=0;' in log,log
        if label=='deadline-disabled-natural':
            assert re.search(r'deadline arm elapsed=.*timed_out=false exit=0 observations=\[\]',log),log
            assert '0 completed callbacks joined before seam restore' in log,log
            assert log.count('confirmed not live')==2,log
        else:
            assert re.search(r'deadline arm elapsed=.*timed_out=true exit=null observations=\[\{leader:[0-9]+ pids:\[[0-9]+ [0-9]+\] states:\[alive alive\] groups:\[[0-9]+ [0-9]+\] err:<nil> killErr:<nil>',log),log
        if label=='leader-only-kill':
            assert log.count('cleanup revalidated fixture child pid=')==2,log
            assert log.count('confirmed not live')==2,log
        records.append(dict(mutant=label,file=file,mutated_sha256=sha,diff='mutation-'+label+'.diff',command=COMMAND,returncode=rc,red_assertion=identity,three_full_trampoline_controls='PASS',specific_child_readiness_and_live_group='PASS' if label!='deadline-disabled-natural' else 'natural exit/no kill',finite_cleanup='PASS'))
    finally:
        for f,b in ORIGINAL.items(): (ROOT/f).write_bytes(b)
        assert {f:hashlib.sha256((ROOT/f).read_bytes()).hexdigest() for f in FILES}==SHAS
    rc,log=invoke('restored-'+label)
    assert rc==0 and '--- PASS: TestExecHandlerTimeoutKillsTheGroup' in log,(label,rc,log)
    records[-1]['restored_sha256']=SHAS
    records[-1]['restored_control']='PASS'
    (BANK/'executor-mutations-progress.json').write_text(json.dumps(records,indent=2)+'\n')

test='host/broker/handlers_exec_m2_test.go'; prod='host/broker/handlers.go'
mutant('leader-only-kill',test,lambda raw:replace(raw,b'o.killErr = originalKill(leader)',b'o.killErr = proctest.Signal(leader, syscall.SIGKILL, syscall.Kill)'), 'specific fixture child alive after deadline')
def lose_stream(raw,stream):
    split=raw.index(b'func runHeadTail(')
    old=b'if waitErr != nil && runCtx.Err() == context.DeadlineExceeded {'
    new=old+b'\n\t\tc.'+stream.encode()+b' = newHeadTailSink(execHeadBytes, execTailBytes, c.hardCap, nil)'
    return raw[:split]+replace(raw[split:],old,new)
for stream in ['stdout','stderr']:
    mutant('deadline-'+stream+'-loss',prod,lambda raw,stream=stream:lose_stream(raw,stream),'exact deadline '+stream+' retention')
def remove_deadline(raw):
    raw=replace(raw,b'runCtx, cancel := context.WithTimeout(ctx, bounds.execTimeout)',b'runCtx, cancel := context.WithCancel(ctx)')
    split=raw.index(b'func runHeadTail(')
    old=b'cleanupDeadline := time.Now().Add(bounds.execTimeout + 2*pipeCloseGrace)'
    new=b'cleanupDeadline := func() time.Time { d, ok := runCtx.Deadline(); if !ok { panic("mutation requires finite caller deadline") }; return d.Add(-pipeCloseGrace) }()'
    return raw[:split]+replace(raw[split:],old,new)
mutant('deadline-disabled-natural',prod,remove_deadline,'real exact-profile timeout required B=')
(BANK/'executor-mutations.json').write_text(json.dumps(dict(isolated_copy=str(ROOT),all_four_compiled_expected_red=True,all_four_byte_exact_restored_green=True,mutations=records),indent=2)+'\n')
print('PASS: four intended RED identities; readiness/control scope checked; four exact restores and fresh GREEN controls')
