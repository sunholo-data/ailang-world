from run import *
B=pathlib.Path.home()/'.ailang/state/world-iter268/evaluator-r2-backups'; originals={f:(B/f).read_bytes() for f in json.loads((E/'start.json').read_text())['files']}
mut=[]
def m(name,old,new,sub,assertion,file='cmd/ailang-worldd/demo_seed.go',test='TestDemoSeedPreflightRefusals'):
 mut.append(dict(name=name,file=file,old=old,new=new,pattern='^'+test+'$'+('/^'+sub+'$' if sub else ''),expected=assertion))
def guard(name,cond,sub,expected='refusal changed public inventory',file='cmd/ailang-worldd/demo_seed.go'):
 m(name,'if '+cond+' {','if false && ('+cond+') {',sub,expected,file)
m('manifest','!os.IsNotExist(err) {','false && !os.IsNotExist(err) {','manifest','refusal changed public inventory')
guard('db-identity','want != got','wrong-db','wrong refusal')
guard('pin-hash','pinRef.String() != h.InterpreterRef','health-interpreter')
guard('archive-absence','err != nil','archive-absent','wrong refusal') if False else None
m('archive-absence','if err != nil {\n\t\treturn h, 0, fmt.Errorf("interpreter archive manifest unavailable")','if false && err != nil {\n\t\treturn h, 0, fmt.Errorf("interpreter archive manifest unavailable")','archive-absent','wrong refusal')
m('archive-hash','manifest.Hash != pinRef.String() || manifest.Size != int64(len(pin))','false && manifest.Hash != pinRef.String() || manifest.Size != int64(len(pin))','archive-hash','refusal changed public inventory')
m('archive-size','manifest.Hash != pinRef.String() || manifest.Size != int64(len(pin))','manifest.Hash != pinRef.String() || false && manifest.Size != int64(len(pin))','archive-size','refusal changed public inventory')
m('archive-health-version','manifest.Version != h.InterpreterVersion || !strings.HasPrefix','false && manifest.Version != h.InterpreterVersion || !strings.HasPrefix','health-version','refusal changed public inventory')
m('archive-release','|| !strings.HasPrefix(manifest.Version, "AILANG v0.41.0\\n")','|| false && !strings.HasPrefix(manifest.Version, "AILANG v0.41.0\\n")','archive-release','refusal changed public inventory')
guard('nonempty','hasWorld || hasHead || ref != "" || len(page) != 0','world-only')
guard('epoch-integrity','!ok || hashref.SumSHA256(obj.Payload).String() != rh.Head','epoch-hash','wrong refusal')
m('epoch-decode','if err := json.Unmarshal(obj.Payload, &reg); err != nil {','if err := json.Unmarshal(obj.Payload, &reg); false && err != nil {','epoch-decode','wrong refusal')
m('epoch-nomination','candidate == "AILANG v0.41.0"','candidate != ""','epoch-candidate','refusal changed public inventory')
guard('session-denial','out.Denied != nil','expired-session','wrong refusal')
m('episode','out.Success == nil || out.Success.EpisodeID != "demo"','out.Success == nil','wrong-episode','refusal changed public inventory')
guard('grant','!broker.Allows(snap, req).Allowed','missing-grant')
m('spend','caps[i].Budget -= spend[','caps[i].Budget -= 0 * spend[','spent-budget','refusal changed public inventory')
m('world-row-scan','worlds.Scanned > 0','worlds.Scanned > 1','world-only','refusal changed public inventory')
guard('CAS','actual != observed || has != (observed != "")','cas-change')
guard('required-GET','!ok','health-absent','wrong refusal')
m('GET-read-error','if err != nil {\n\t\treturn err\n\t}\n\tif !ok {','if false && err != nil {\n\t\treturn err\n\t}\n\tif !ok {','health-down','wrong refusal')
m('pin-read','if err != nil {\n\t\treturn h, 0, fmt.Errorf("cannot read interpreter pin")','if false && err != nil {\n\t\treturn h, 0, fmt.Errorf("cannot read interpreter pin")','pin-absent','wrong refusal')
m('source-read','if err != nil {\n\t\treturn fmt.Errorf("prepare worktree','if false && err != nil {\n\t\treturn fmt.Errorf("prepare worktree','source-absent','refusal changed public inventory')
m('readonly-open','if err != nil {\n\t\treturn false, fmt.Errorf("cannot open read-only authority view")','if err != nil {\n\t\treturn false, fmt.Errorf("incorrect read-only refusal")','readonly-db','wrong refusal')
guard('required-flags','o.db == "" || o.pin == "" || o.key == "" || o.manifest == ""','missing-key',file='cmd/ailang-worldd/demo.go')
guard('positional','fs.NArg() != 0','positional',file='cmd/ailang-worldd/demo.go')
m('subcommand','len(args) == 0 || args[0] != "seed"','len(args) == 0','unknown-verb','refusal changed public inventory',file='cmd/ailang-worldd/demo.go')
for name,old,new,sub,expected in [
 ('commit-status','if status != 200 {','if false && status != 200 {','commit-denied','commit refusal boundary'),
 ('commit-response','if json.Unmarshal(body, &result) != nil || result.SelectedHead != commit.NextWorld.Ref {','if false && (json.Unmarshal(body, &result) != nil || result.SelectedHead != commit.NextWorld.Ref) {','commit-response','no retries'),
 ('commit-transport','if err != nil {\n\t\t\treturn fmt.Errorf("commit outcome uncertain','if false && err != nil {\n\t\t\treturn fmt.Errorf("commit outcome uncertain','commit-uncertain','commit refusal boundary')]:m(name,old,new,sub,expected,test='TestDemoSeedCommitFailures')
for name,old,new,expected in [
 ('bearer','bytes.NewReader(data), "Bearer "+token)','bytes.NewReader(data), "")','SessionAbsent'),
 ('header-interpreter','Interpreter: interpreter,','Interpreter: o.Hash,','frozen header identity mismatch'),
 ('header-epoch','SemanticsEpoch: epoch,','SemanticsEpoch: epoch+1,','frozen header identity mismatch'),
 ('carrying','TransitionRef: o.Hash}}','TransitionRef: interpreter}}','no real carrying entry')]:m(name,old,new,'',expected,test='TestDemoSeedAuthenticatedCommits')
m('doc-command','"$DEMO_BIN" demo seed --addr','"$OTHER_BIN" demo seed --addr','','demo artifact command wiring missing',file='docs/QUICKSTART.md',test='TestDemoSeedHelpPreparation')
m('timeout-coverage','budgetContext(90 * time.Second)','budgetContext(24 * time.Hour)','','',file='cmd/ailang-worldd/demo.go',test='TestDemoSeed')
mut[-1]['pattern']='^TestDemoSeed'
results=[]
try:
 for x in mut:
  f=W/x['file']; original=originals[x['file']]; s=original.decode(); assert s.count(x['old'])==1,(x['name'],s.count(x['old']))
  changed=s.replace(x['old'],x['new'],1).encode();f.write_bytes(changed);assert f.read_bytes()==changed
  start=time.time(); exe=B/'mutant.test';buildargv=['go','test','-c','-o',str(exe),'./cmd/ailang-worldd']
  b=subprocess.run(buildargv,cwd=W,env=ENV,stdout=subprocess.PIPE,stderr=subprocess.STDOUT,timeout=90);(E/('build-'+x['name']+'.log')).write_bytes(b.stdout)
  x.update(build_argv=buildargv,build_rc=b.returncode,mutant_sha256=hashlib.sha256(changed).hexdigest())
  if b.returncode==0:
   argv=[str(exe),'-test.run='+x['pattern'],'-test.v','-test.count=1','-test.timeout=120s'];r=subprocess.run(argv,cwd=W/'cmd/ailang-worldd',env=ENV,stdout=subprocess.PIPE,stderr=subprocess.STDOUT,timeout=150)
   (E/('mutant-'+x['name']+'.log')).write_bytes(r.stdout);log=r.stdout.decode();x.update(argv=argv,rc=r.returncode,seconds=round(time.time()-start,3),assertion_kill=r.returncode!=0 and x['expected'] in log and '--- FAIL: TestDemoSeed' in log and not any(t in log for t in ['panic:','SKIP','no tests to run']))
  f.write_bytes(original);assert f.read_bytes()==original;x['restored_sha256']=hashlib.sha256(original).hexdigest();results.append(x);(E/'attacks.json').write_text(json.dumps(results,indent=2));print(x['name'],x['build_rc'],x.get('rc'),x.get('assertion_kill'),flush=True)
finally:
 for f,v in originals.items():
  assert (W/f).read_bytes()==v,(f,'unrestored')
