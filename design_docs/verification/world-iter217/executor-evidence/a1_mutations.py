from pathlib import Path
import subprocess,hashlib,json
cases=[('MUT-ABSENCE-UNTYPED','host/transitionreg/transitionreg.go','return Snapshot{}, RegistryHeadAbsentError{}','return Snapshot{}, errors.New("read transition registry: head is absent")','TestRegistryHeadAbsentErrorClassification/direct'),('MUT-ABSENCE-WRAP-LOST','host/transitionreg/bind.go','construct request: %w','construct request: %v','TestRegistryHeadAbsentErrorClassification/new_request'),('MUT-ABSENCE-POSTLOOKUP-CTX-REMOVED','host/transitionreg/transitionreg.go','\tif err := ctx.Err(); err != nil {\n\t\treturn Snapshot{}, fmt.Errorf("read transition registry: context: %w", err)\n\t}\n\tif !ok {','\tif !ok {','TestRegistryHeadAbsentAfterLookupCancellation/absent')]
for name,file,old,new,test in cases:
 p=Path(file);original=p.read_bytes();s=original.decode();assert old in s
 try:
  p.write_text(s.replace(old,new,1));mut=hashlib.sha256(p.read_bytes()).hexdigest()
  fence=subprocess.run(['python3','/tmp/world-iter217-executor/run.py',name+'-fence','go','test','-p','1','./host/transitionreg','-run','^$','-timeout=60s']);assert fence.returncode==0
  r=subprocess.run(['python3','/tmp/world-iter217-executor/run.py',name,'go','test','-p','1','./host/transitionreg','-run','^'+test+'$','-v','-count=1','-timeout=60s']);assert r.returncode==1
 finally: p.write_bytes(original)
 r=subprocess.run(['python3','/tmp/world-iter217-executor/run.py',name+'-restored','go','test','-p','1','./host/transitionreg','-run','^'+test+'$','-v','-count=1','-timeout=60s']);assert r.returncode==0
 Path('/tmp/world-iter217-executor/'+name+'-hash.json').write_text(json.dumps({'original_restored_sha256':hashlib.sha256(original).hexdigest(),'mutant_sha256':mut})+'\n')
