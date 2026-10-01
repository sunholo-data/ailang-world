from pathlib import Path
import subprocess,hashlib,json
import ast
source=Path('/tmp/world-iter217-executor/d_full_mutations.py').read_text();tree=ast.parse(source);loop=next(n for n in tree.body if isinstance(n,ast.For));exec('\n'.join(source.splitlines()[:loop.lineno-1]))
name,edits,pkg,test,kind=cases[0];originals={f:Path(f).read_bytes() for f in edits};hashes={}
try:
 for f,changes in edits.items():
  s=originals[f].decode()
  for old,new in changes:assert old in s;s=s.replace(old,new,1)
  Path(f).write_text(s);hashes[f]={'mutant':hashlib.sha256(Path(f).read_bytes()).hexdigest(),'restored':hashlib.sha256(originals[f]).hexdigest()}
 args=['go','test','-p','1',pkg,'-timeout=60s']
 for suffix,run,rc in [('current-fence','^$',0),('current','^'+test+'$',1)]:
  r=subprocess.run(['python3','/tmp/world-iter217-executor/run.py',name+'-'+suffix,*args,'-run',run,'-v','-count=1']);assert r.returncode==rc
finally:
 for f,b in originals.items():Path(f).write_bytes(b)
r=subprocess.run(['python3','/tmp/world-iter217-executor/run.py',name+'-current-restored',*args,'-run','^'+test+'$','-v','-count=1']);assert r.returncode==0
Path('/tmp/world-iter217-executor/'+name+'-current-hash.json').write_text(json.dumps({'files':hashes,'selected_test':test,'test_file_sha256':hashlib.sha256(Path('host/projection/mcp_conformance_test.go').read_bytes()).hexdigest(),'kind':kind})+'\n')
