from pathlib import Path
import subprocess,json,hashlib
for name in ['MUT-PIN-REGRESS','MUT-FACADE-IMPORT-MCP']:
 files=['go.mod','go.sum','host/projection/mcp.go'];originals={f:Path(f).read_bytes() for f in files}
 try:
  if name=='MUT-PIN-REGRESS':Path('go.mod').write_text(originals['go.mod'].decode().replace('ailang v0.47.2','ailang v0.33.2'))
  else:Path('host/projection/mcp.go').write_text(originals['host/projection/mcp.go'].decode().replace('"context"','"context"\n _ "github.com/sunholo-data/ailang/serveapi"'))
  if name=='MUT-FACADE-IMPORT-MCP':
   fence=subprocess.run(['python3','/tmp/world-iter217-executor/run.py',name+'-fence','go','test','-mod=mod','-p','1','./host/daemon','-run','^$','-timeout=60s']);assert fence.returncode==0,('facade compile fence',fence.returncode)
  r=subprocess.run(['python3','/tmp/world-iter217-executor/run.py',name,'go','test',('-mod=readonly' if name=='MUT-PIN-REGRESS' else '-mod=mod'),'-p','1','./host/daemon','-run','^TestDaemonDependencyAllowlist$','-v','-count=1','-timeout=60s'])
  verdict=r.returncode
  sha={f:{'mutant':hashlib.sha256(Path(f).read_bytes()).hexdigest(),'restored':hashlib.sha256(b).hexdigest()} for f,b in originals.items()}
 finally:
  for f,b in originals.items():Path(f).write_bytes(b)
 r=subprocess.run(['python3','/tmp/world-iter217-executor/run.py',name+'-restored','go','test','-p','1','./host/daemon','-run','^TestDaemonDependencyAllowlist$','-v','-count=1','-timeout=60s'])
 Path('/tmp/world-iter217-executor/'+name+'-hash.json').write_text(json.dumps({'files':sha,'applied_rc':verdict,'restored_rc':r.returncode,'kind':'dependency boundary; may be expected compile failure, classify actual raw output'})+'\n')
 assert verdict==1 and r.returncode==0,(name,verdict,r.returncode)
