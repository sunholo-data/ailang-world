from pathlib import Path
import subprocess,json,hashlib
base=Path('/tmp/world-iter217-executor-repair-round1');p=Path('host/daemon/invoke_e2e_test.go')
cases=[('duplicate-task', 'body, err = socketCall(client, server.URL, token, "task-resend")', 'body, err = socketCall(client, server.URL, token, "task-resend-duplicate")','same task executed 2 times'),('failure-cleanup','\t// Keep the response blocked beyond the server\'s WriteTimeout, after the','\tt.Fatal("forced post-commit failure-cleanup control")\n\t// Keep the response blocked beyond the server\'s WriteTimeout, after the','forced post-commit failure-cleanup control')]
for name,old,new,message in cases:
 original=p.read_bytes()
 try:
  text=original.decode();assert old in text,(name,old);p.write_text(text.replace(old,new,1));mut=hashlib.sha256(p.read_bytes()).hexdigest()
  args=['go','test','-p','1','-timeout=30s','./host/daemon','-run','^TestA2AResendAfterWriteTimeout$','-v','-count=1']
  r=subprocess.run(['python3',str(base/'run.py'),name,*args]);assert r.returncode==1,(name,r.returncode)
  raw=(base/(name+'.log')).read_text();assert message in raw and 'panic: test timed out' not in raw and 'worker did not stop' not in raw,name
 finally:p.write_bytes(original)
 r=subprocess.run(['python3',str(base/'run.py'),name+'-restored',*args]);assert r.returncode==0,(name,'restore',r.returncode)
 assert p.read_bytes()==original
 (base/(name+'-hash.json')).write_text(json.dumps({'kind':'test-only behavioral/failure-path calibration; no kernel or production mutation','mutant_sha256':mut,'restored_sha256':hashlib.sha256(original).hexdigest(),'applied_rc':1,'restored_rc':0,'named_failure':message},indent=2)+'\n')
