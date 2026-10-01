from pathlib import Path
import subprocess,json,hashlib,re
cases=[
('MUT-PLAIN-JSON','host/projection/mcp_conformance_test.go','actual := postMCP(t, h, "Bearer "+tok, body)','actual := postMCP(t, h, "Bearer "+tok, body)\n actual.Header().Set("Content-Type","application/json"); actual.Body.Reset(); actual.Body.WriteString(`{"jsonrpc":"2.0","id":7,"result":{}}`)','./host/projection','TestMCPWireConformance','test response decorator'),
('MUT-AGGREGATE-REMOVED','host/projection/mcp.go','context.WithTimeout(r.Context(), h.invokeWait)','context.WithTimeout(r.Context(), time.Minute)','./host/daemon','TestMCPPostBudgetProductionConstants/single','production aggregate extended'),
('MUT-AGGREGATE-PER-ITEM','host/projection/mcp.go','context.WithTimeout(r.Context(), h.invokeWait)','context.WithTimeout(r.Context(), time.Minute)','./host/daemon','TestMCPPostBudgetProductionConstants/batch','production aggregate extended; each callback retains own20s'),
]
for name,f,old,new,pkg,test,kind in cases:
 if name != "MUT-AGGREGATE-PER-ITEM": continue
 orig=Path(f).read_bytes()
 try:
  text=orig.decode();assert old in text;Path(f).write_text(text.replace(old,new,1));mut=hashlib.sha256(Path(f).read_bytes()).hexdigest()
  args=['go','test','-p','1',pkg,'-timeout=60s']
  for suffix,run,expected in [('fence','^$',0),('', '^'+test+'$',1)]:
   r=subprocess.run(['python3','/tmp/world-iter217-executor/run.py',name+('-'+suffix if suffix else ''),*args,'-run',run,'-v','-count=1']);assert r.returncode==expected,(name,suffix,r.returncode)
   if expected==1: assert 'aggregate/wire elapsed=' in Path('/tmp/world-iter217-executor/'+name+'.log').read_text(), 'setup failure is not a kill'
  # Timing controls enumerate exactly the selected runtime red set; no invented full-package red set.
  red=re.findall(r'^\s*--- FAIL: (\S+)',Path('/tmp/world-iter217-executor/'+name+'.log').read_text(),re.M)
 finally:Path(f).write_bytes(orig)
 r=subprocess.run(['python3','/tmp/world-iter217-executor/run.py',name+'-restored',*args,'-run','^'+test+'$','-v','-count=1']);assert r.returncode==0
 Path('/tmp/world-iter217-executor/'+name+'-hash.json').write_text(json.dumps({'files':{f:{'mutant':mut,'restored':hashlib.sha256(orig).hexdigest()}},'selected_test':test,'selected_failed_tests':red,'kind':kind,'broad_execution':False})+'\n')
