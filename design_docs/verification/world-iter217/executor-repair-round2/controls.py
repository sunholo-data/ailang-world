from run import run,root
from pathlib import Path
import hashlib,json
p=Path('host/projection/mcp.go');original=p.read_bytes();before=hashlib.sha256(original).hexdigest()
rows=[]
for name,old,new in [('cb-background','context.WithTimeout(ctx, a.h.invokeWait)','context.WithTimeout(context.Background(), a.h.invokeWait)'),('aggregate-removed','context.WithTimeout(r.Context(), h.invokeWait)','context.WithCancel(r.Context())')]:
 try:
  s=original.decode();assert s.count(old)==1;p.write_text(s.replace(old,new))
  mutant=hashlib.sha256(p.read_bytes()).hexdigest()
  run(name+'-compile',['go','test','-p','1','-timeout=60s','./host/projection','-run','^$'])
  run(name+'-selected',['go','test','-p','1','-timeout=60s','-race','./host/projection','-run','^TestMCPInvokeDeadlineFreesSlot$','-v','-count=1'],1)
  assert 'callback lost aggregate context:' in (root/(name+'-selected.log')).read_text()
 finally:p.write_bytes(original)
 assert hashlib.sha256(p.read_bytes()).hexdigest()==before
 run(name+'-restored',['go','test','-p','1','-timeout=60s','-race','./host/projection','-run','^TestMCPInvokeDeadlineFreesSlot$','-v','-count=1'])
 rows.append({'name':name,'path':str(p),'old':old,'new':new,'original_sha256':before,'mutant_sha256':mutant,'restored_sha256':hashlib.sha256(p.read_bytes()).hexdigest(),'scope':'selected named test only; no broad package red-set claim'})
 (root/'mutation-manifest.json').write_text(json.dumps(rows,indent=2)+'\n')
