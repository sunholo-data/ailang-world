from pathlib import Path
import subprocess,json,hashlib,re
exec(Path('/tmp/world-iter217-executor/a2_mutations.py').read_text().split('for name,old,new,test in cases[2:]:')[0])
controls=[]
for name,old,new,test in cases[:2]:
 test='TestMCPAbsentPrecheckSnapshotFailure/list/'+('direct' if name.endswith('CATCHALL') else 'same_text')
 controls.append((name,f,old,new,test))
controls.append(('MUT-MAP-NO-LEN-GUARD','host/projection/mcpname.go','if len(name) > 64 {','if false {','TestMCPSurfaceRefusalRecovers'))
for name,f,old,new,test in controls:
 orig=Path(f).read_bytes()
 try:
  text=orig.decode();assert old in text,(name,old);Path(f).write_text(text.replace(old,new,1));sha=hashlib.sha256(Path(f).read_bytes()).hexdigest()
  args=['go','test','-p','1','./host/projection','-timeout=60s']
  for suffix,run,expected in [('mcp-fence','^$',0),('mcp','^'+test+'$',1),('mcp-whole-package','.',1)]:
   r=subprocess.run(['python3','/tmp/world-iter217-executor/run.py',name+'-'+suffix,*args,'-run',run,'-v','-count=1']);assert r.returncode==expected,(name,suffix,r.returncode)
 finally:Path(f).write_bytes(orig)
 r=subprocess.run(['python3','/tmp/world-iter217-executor/run.py',name+'-mcp-restored',*args,'-run','^'+test+'$','-v','-count=1']);assert r.returncode==0
 red=re.findall(r'^\s*--- FAIL: (\S+)',Path('/tmp/world-iter217-executor/'+name+'-mcp-whole-package.log').read_text(),re.M)
 Path('/tmp/world-iter217-executor/'+name+'-mcp-hash.json').write_text(json.dumps({'files':{f:{'mutant':sha,'restored':hashlib.sha256(orig).hexdigest()}},'selected_test':test,'broad_scope':'whole projection','broad_failed_tests':red})+'\n')
