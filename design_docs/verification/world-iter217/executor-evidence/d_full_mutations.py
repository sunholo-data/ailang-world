from pathlib import Path
import subprocess,json,hashlib,re
# Reuse exact recorded edits, but execute whole packages rather than scoped MCP sets.
exec(Path('/tmp/world-iter217-executor/d_mutations.py').read_text().split('for name,edits,pkg,test,kind in cases:')[0])
cases[0]=('MUT-TASK-PER-POST',{'host/projection/mcp.go':[
 ('type mcpAdapter struct{ h *Handler }','type mutantTaskKey struct{}\n type mcpAdapter struct{ h *Handler }'),
 ('defer cancel()\n\th.mcp.ServeHTTP','defer cancel()\n ctx=context.WithValue(ctx,mutantTaskKey{},new(string))\n\th.mcp.ServeHTTP'),
 ('task, err := a.h.mintTask()','slot:=ctx.Value(mutantTaskKey{}).(*string)\n task:=*slot\n if task=="" {task,err=a.h.mintTask();*slot=task}')
]},'./host/projection','TestMCPBatchAccounting/repeated_rpc_id','production lazy task cache scoped to one POST context')
cases.extend([
 ('MUT-AGGREGATE-REMOVED',{'host/projection/mcp.go':[('context.WithTimeout(r.Context(), h.invokeWait)','context.WithCancel(r.Context())')]},'./host/daemon','TestMCPPostBudgetProductionConstants/single','production exact aggregate removal'),
 ('MUT-AGGREGATE-PER-ITEM',{'host/projection/mcp.go':[('context.WithTimeout(r.Context(), h.invokeWait)','context.WithCancel(r.Context())')]},'./host/daemon','TestMCPPostBudgetProductionConstants/batch','production only fresh per-item/Runner20s contexts remain; aggregate removed'),
])
for name,edits,pkg,test,kind in cases:
 originals={f:Path(f).read_bytes() for f in edits};hashes={}
 try:
  for f,changes in edits.items():
   text=originals[f].decode()
   for old,new in changes:assert old in text,(name,old);text=text.replace(old,new,1)
   Path(f).write_text(text);hashes[f]={'mutant':hashlib.sha256(Path(f).read_bytes()).hexdigest(),'restored':hashlib.sha256(originals[f]).hexdigest()}
  args=['go','test','-p','1',pkg,'-timeout=120s']
  # Exact task and aggregate shapes get fresh selected controls too.
  if name in ('MUT-TASK-PER-POST','MUT-AGGREGATE-REMOVED','MUT-AGGREGATE-PER-ITEM'):
   for suffix,run,expected in [('fence','^$',0),('', '^'+test+'$',1)]:
    r=subprocess.run(['python3','/tmp/world-iter217-executor/run.py',name+('-'+suffix if suffix else ''),*args,'-run',run,'-v','-count=1']);assert r.returncode==expected,(name,suffix,r.returncode)
    if expected==1 and 'AGGREGATE' in name:assert 'aggregate/wire elapsed=' in Path('/tmp/world-iter217-executor/'+name+'.log').read_text(),'fixture failure is not a kill'
  r=subprocess.run(['python3','/tmp/world-iter217-executor/run.py',name+'-whole-package',*args,'-v','-count=1']);assert r.returncode==1,(name,'broad',r.returncode)
  log=Path('/tmp/world-iter217-executor/'+name+'-whole-package.log').read_text();red=re.findall(r'^\s*--- FAIL: (\S+)',log,re.M)
  assert red
 finally:
  for f,b in originals.items():Path(f).write_bytes(b)
 r=subprocess.run(['python3','/tmp/world-iter217-executor/run.py',name+'-restored',*args,'-run','^'+test+'$','-v','-count=1']);assert r.returncode==0
 Path('/tmp/world-iter217-executor/'+name+'-hash.json').write_text(json.dumps({'files':hashes,'selected_test':test,'kind':kind,'broad_scope':'whole '+pkg,'broad_failed_tests':red})+'\n')
