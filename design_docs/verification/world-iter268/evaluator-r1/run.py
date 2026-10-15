import os,pathlib,subprocess,json,time,hashlib
root=pathlib.Path('/Users/voightkampff/dev/sunholo-data/.eval-world-iter268'); out=pathlib.Path('/tmp/world268-evaluator-r1'); env=os.environ.copy(); env.update(PATH='/opt/homebrew/bin:'+env['PATH'],AILANG_BIN=str(pathlib.Path.home()/'.pinned-ailang/ailang'),WORLD_TOOL_AILANG_BIN=str(pathlib.Path.home()/'.pinned-ailang-tools/v0.52.1/ailang'),GOTOOLCHAIN='go1.26.6')
for k in ['GOOGLE_CLOUD_PROJECT','OTEL_EXPORTER_OTLP_ENDPOINT','OTEL_RESOURCE_ATTRIBUTES','OTLP_GOOGLE_CLOUD_PROJECT']:env.pop(k,None)
files=['cmd/ailang-worldd/demo.go','cmd/ailang-worldd/demo_seed.go','cmd/ailang-worldd/main.go','cmd/ailang-worldd/demo_test.go','docs/QUICKSTART.md']; original={f:(root/f).read_bytes() for f in files}
for f,b in original.items():p=out/'original'/f;p.parent.mkdir(parents=True,exist_ok=True);p.write_bytes(b)
results=[]
def run(name,pattern,pkg='./cmd/ailang-worldd',expect=None):
 cmd=['go','test','-run',pattern,pkg,'-count=1','-timeout','120s','-v']; t=time.monotonic();p=subprocess.run(cmd,cwd=root,env=env,stdout=subprocess.PIPE,stderr=subprocess.STDOUT,timeout=150);s=p.stdout.decode();(out/(name+'.log')).write_bytes(p.stdout);r=dict(name=name,command=cmd,rc=p.returncode,seconds=round(time.monotonic()-t,3),expected_assertion=expect)
 if expect:r['compiling_assertion_kill']=p.returncode!=0 and expect in s and '--- FAIL:' in s and not any(x in s for x in ['[build failed]','panic:','SKIP','no tests to run'])
 else:r['named_green']=p.returncode==0 and '--- PASS:' in s and not any(x in s for x in ['SKIP','no tests to run'])
 results.append(r);(out/'ledger.json').write_text(json.dumps(results,indent=2));print(name,r,flush=True)
for name in ['AuthenticatedCommits','PreflightRefusals','HelpPreparation','CommitFailures']:run('green-'+name,'^TestDemoSeed'+name+'$')
run('green-mint','^TestSessionMint_RequiresControllingTerminal$');run('green-no-forgery','^TestRESTCommitRefusesDecisionPacketObjects$','./host/daemon')
mut=[]
def m(name,old,new,pat,assertion,file='cmd/ailang-worldd/demo_seed.go'):mut.append((name,file,old,new,pat,assertion))
p='^TestDemoSeedPreflightRefusals$/';h='^TestDemoSeedAuthenticatedCommits$'
m('bearer','bytes.NewReader(data), "Bearer "+token)','bytes.NewReader(data), "")',h,'SessionAbsent')
m('unselected-world','if hasWorld || hasHead || ref != "" || len(page) != 0 {','if false && (hasWorld || hasHead || ref != "" || len(page) != 0) {',p+'^world-only$','refusal changed public inventory')
m('world-scan','return worlds.Scanned > 0, nil','return worlds.Scanned > 1, nil',p+'^world-only$','refusal changed public inventory')
m('grant-check','if !broker.Allows(snap, req).Allowed {','if false && !broker.Allows(snap, req).Allowed {',p+'^missing-grant$','refusal changed public inventory')
m('spent-budget','caps[i].Budget -= spend[','caps[i].Budget -= 0 * spend[',p+'^spent-budget$','refusal changed public inventory')
m('manifest','!os.IsNotExist(err) {','false && !os.IsNotExist(err) {',p+'^manifest$','refusal changed public inventory')
m('identity','if want != got {','if false && want != got {',p+'^wrong-db$','wrong refusal')
m('pin-hash','if hashref.SumSHA256(pin).String() != h.InterpreterRef {','if false && hashref.SumSHA256(pin).String() != h.InterpreterRef {',p+'^health-interpreter$','refusal changed public inventory')
m('pin-version','if err != nil || strings.TrimSpace(string(version)) != strings.TrimSpace(h.InterpreterVersion) || !strings.HasPrefix(string(version), "AILANG v0.41.0\\n") {','if false && (err != nil || strings.TrimSpace(string(version)) != strings.TrimSpace(h.InterpreterVersion) || !strings.HasPrefix(string(version), "AILANG v0.41.0\\n")) {',p+'^health-version$','refusal changed public inventory')
m('episode','if out.Success == nil || out.Success.EpisodeID != "demo" {','if out.Success == nil {',p+'^wrong-episode$','refusal changed public inventory')
m('session-denial','if out.Denied != nil {','if false && out.Denied != nil {',p+'^expired-session$','wrong refusal')
m('epoch-hash','if !ok || hashref.SumSHA256(obj.Payload).String() != rh.Head {','if false && (!ok || hashref.SumSHA256(obj.Payload).String() != rh.Head) {',p+'^epoch-hash$','wrong refusal')
m('epoch-nomination','if candidate == "AILANG v0.41.0" {','if candidate != "" {',p+'^epoch-candidate$','refusal changed public inventory')
m('CAS','if actual != observed || has != (observed != "") {','if false && (actual != observed || has != (observed != "")) {',p+'^cas-change$','refusal changed public inventory')
m('commit-status','if status != 200 {','if false && status != 200 {','^TestDemoSeedCommitFailures$/^commit-denied$','commit refusal boundary')
m('commit-response','if json.Unmarshal(body, &result) != nil || result.SelectedHead != commit.NextWorld.Ref {','if false && (json.Unmarshal(body, &result) != nil || result.SelectedHead != commit.NextWorld.Ref) {','^TestDemoSeedCommitFailures$/^commit-response$','no retries')
m('header','SemanticsEpoch: epoch,','SemanticsEpoch: epoch + 1,',h,'frozen header identity mismatch')
m('carrying','TransitionRef: o.Hash}}','TransitionRef: interpreter}}',h,'no real carrying entry')
m('flags','if o.db == "" || o.pin == "" || o.key == "" || o.manifest == "" {','if false && (o.db == "" || o.pin == "" || o.key == "" || o.manifest == "") {',p+'^missing-key$','refusal changed public inventory','cmd/ailang-worldd/demo.go')
m('positional','if fs.NArg() != 0 {','if false && fs.NArg() != 0 {',p+'^positional$','refusal changed public inventory','cmd/ailang-worldd/demo.go')
try:
 for name,f,old,new,pattern,assertion in mut:
  source=original[f].decode();assert source.count(old)==1,(name,source.count(old));(root/f).write_text(source.replace(old,new));run('mutant-'+name,pattern,expect=assertion);(root/f).write_bytes(original[f]);assert (root/f).read_bytes()==original[f]
finally:
 for f,b in original.items():(root/f).write_bytes(b)
for name in ['AuthenticatedCommits','PreflightRefusals','HelpPreparation','CommitFailures']:run('restored-'+name,'^TestDemoSeed'+name+'$')
assert all((root/f).read_bytes()==b for f,b in original.items());print('TERMINAL RESTORED',flush=True)
