"""Six isolated real-process mutation arms, each followed by exact restore and one pristine control."""
from pathlib import Path
import hashlib,json,os,re,subprocess
ROOT=Path('/private/tmp/ailang-world-row166-mutations-20261010')
BANK=Path(__file__).resolve().parent
FILES=['host/broker/handlers.go','host/broker/handlers_capture_test.go']
ORIGINAL={f:(ROOT/f).read_bytes() for f in FILES}
SHAS={f:hashlib.sha256(b).hexdigest() for f,b in ORIGINAL.items()}
COMMAND=['go','test','./host/broker','-run','^TestRunBoundedHeadTailTimeoutKeepsPartialOutput$','-v','-count=1','-timeout=600s']
records=[]
def replace_once(raw,old,new):
    assert raw.count(old)==1,(old,raw.count(old))
    return raw.replace(old,new)
def invoke(label):
    r=subprocess.run(COMMAND,cwd=ROOT,capture_output=True,text=True)
    log=r.stdout+r.stderr
    (BANK/(label+'.log')).write_text(log)
    assert '--- SKIP:' not in log,log
    return r.returncode,log

def mutant(label,changes,identity):
    try:
        for f,fn in changes.items(): (ROOT/f).write_bytes(fn(ORIGINAL[f]))
        rc,log=invoke('mutation-'+label)
        assert rc!=0 and identity in log,(label,rc,log)
        entry=dict(mutant=label,red_assertion=identity,returncode=rc)
        if label=='deadline-disabled-natural':
            assert re.search(r'deadline arm elapsed=.* err=<nil> kills=\[\] Outstanding=0 baseline=0',log),log
            assert 'deadline snapshots stdout={Text:PARTIAL-OUT Bytes:11' in log,log
            assert 'stderr={Text:PARTIAL-ERR Bytes:11' in log,log
            entry['natural_exit_reap_before_parent_guard']=True
            entry['compound_mutation']='WithTimeout -> WithCancel; runHeadTail cleanupDeadline -> required caller deadline minus existing pipeCloseGrace (two changed production lines only)'
        else:
            entry['exercised_deadline_arm_reached']='deadline arm elapsed=' in log
        records.append(entry)
    finally:
        for f,b in ORIGINAL.items(): (ROOT/f).write_bytes(b)
        assert {f:hashlib.sha256((ROOT/f).read_bytes()).hexdigest() for f in FILES}==SHAS
    rc,log=invoke('restored-'+label)
    assert rc==0 and '--- PASS: TestRunBoundedHeadTailTimeoutKeepsPartialOutput' in log,(label,rc,log)
    records[-1]['byte_exact_restored_sha256']=SHAS
    records[-1]['restored_control']='PASS'
    (BANK/'fixture-mutations-progress.json').write_text(json.dumps(records,indent=2)+'\n')
prod='host/broker/handlers.go'; test='host/broker/handlers_capture_test.go'
mutant('stdout-discard',{prod:lambda x:replace_once(x,b'cmd.Stdout = c.stdout',b'cmd.Stdout = io.Discard')},'exact partial stdout Text/Bytes/SHA/untruncated')
mutant('stderr-discard',{prod:lambda x:replace_once(x,b'cmd.Stderr = c.stderr',b'cmd.Stderr = io.Discard')},'exact partial stderr Text/Bytes/SHA/untruncated')
mutant('stdout-printf-removed',{test:lambda x:replace_once(x,b"printf 'PARTIAL-OUT'; printf 'PARTIAL-ERR' >&2",b"printf 'PARTIAL-ERR' >&2")},'exact partial stdout Text/Bytes/SHA/untruncated')
mutant('stderr-printf-removed',{test:lambda x:replace_once(x,b"printf 'PARTIAL-OUT'; printf 'PARTIAL-ERR' >&2",b"printf 'PARTIAL-OUT'")},'exact partial stderr Text/Bytes/SHA/untruncated')
mutant('readiness-removed',{test:lambda x:replace_once(x,b': > "$READY"',b': # readiness creation removed')},'control readiness after both printf')
def disable_deadline(raw):
    raw=replace_once(raw,b'runCtx, cancel := context.WithTimeout(ctx, bounds.execTimeout)',b'runCtx, cancel := context.WithCancel(ctx)')
    split=raw.index(b'func runHeadTail(')
    suffix=replace_once(raw[split:],b'cleanupDeadline := time.Now().Add(bounds.execTimeout + 2*pipeCloseGrace)',b'cleanupDeadline := func() time.Time { d, ok := runCtx.Deadline(); if !ok { panic("mutation requires caller deadline") }; return d.Add(-pipeCloseGrace) }()')
    return raw[:split]+suffix
mutant('deadline-disabled-natural',{prod:disable_deadline},'exact timeout B=')
(BANK/'fixture-mutations.json').write_text(json.dumps(dict(isolated_copy=str(ROOT),mutations=records,all_six_expected_red=True,all_six_restored_controls='PASS'),indent=2)+'\n')
print('PASS: six intended RED identities; six byte-exact restores and fresh GREEN controls; compound deadline mutant naturally reaped with no kill')
