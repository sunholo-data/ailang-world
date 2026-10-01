from pathlib import Path
import subprocess,json,hashlib,re
for name,leg in [('MUT-AGGREGATE-REMOVED','single'),('MUT-AGGREGATE-PER-ITEM','batch')]:
 f=Path('host/projection/mcp.go');orig=f.read_bytes();test='TestMCPPostBudgetProductionConstants/'+leg
 try:
  old='context.WithTimeout(r.Context(), h.invokeWait)';assert old in orig.decode();f.write_text(orig.decode().replace(old,'context.WithCancel(r.Context())',1));sha=hashlib.sha256(f.read_bytes()).hexdigest()
  args=['go','test','-p','1','./host/daemon','-timeout=60s']
  for suffix,run,expected in [('accounting-fence','^$',0),('accounting','^'+test+'$',1)]:
   r=subprocess.run(['python3','/tmp/world-iter217-executor/run.py',name+'-'+suffix,*args,'-run',run,'-v','-count=1']);assert r.returncode==expected,(name,suffix,r.returncode)
  raw=Path('/tmp/world-iter217-executor/'+name+'-accounting.log').read_text();assert 'production clocks3/10/20/30' in raw and 'aggregate/wire elapsed=' in raw and raw.index('production clocks3/10/20/30')<raw.index('aggregate/wire elapsed=')
 finally:f.write_bytes(orig)
 r=subprocess.run(['python3','/tmp/world-iter217-executor/run.py',name+'-accounting-restored',*args,'-run','^'+test+'$','-v','-count=1']);assert r.returncode==0
 Path('/tmp/world-iter217-executor/'+name+'-accounting-hash.json').write_text(json.dumps({'files':{'host/projection/mcp.go':{'mutant':sha,'restored':hashlib.sha256(orig).hexdigest()}},'test_file_sha256':hashlib.sha256(Path('host/daemon/mcp_budget_test.go').read_bytes()).hexdigest(),'selected_test':test,'kind':'exact aggregate removal; real stage/receipt accounting before timing assertion'})+'\n')
