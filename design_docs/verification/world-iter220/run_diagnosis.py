import os,subprocess,time,datetime,json,hashlib,signal,base64
from pathlib import Path
out=Path(__file__).resolve().parent
candidate=out.parents[2]
base=Path('/Users/voightkampff/dev/sunholo-data/.wt-world-iter220-base')
pin='/Users/voightkampff/.pinned-ailang/ailang'
env=dict(os.environ,AILANG_BIN=pin,WORLD_PKG_AILANG_BIN=pin)
def utc(): return datetime.datetime.now(datetime.timezone.utc).isoformat()
def read(cmd,cwd): return subprocess.check_output(cmd,cwd=cwd,env=env,stderr=subprocess.STDOUT).decode()
meta={'start_utc':utc(),'pin':pin,'pin_sha256':hashlib.sha256(Path(pin).read_bytes()).hexdigest(),'pin_version':read([pin,'--version'],candidate),'go_version':read(['go','version'],candidate),'go_env':json.loads(read(['go','env','-json'],candidate)),'explicit_environment':{k:env.get(k) for k in ['AILANG_BIN','WORLD_PKG_AILANG_BIN','GOFLAGS','GOTOOLCHAIN','GOMAXPROCS','GODEBUG']},'trees':{}}
for label,cwd in [('base',base),('candidate',candidate)]:
 meta['trees'][label]={'cwd':str(cwd),'head':read(['git','rev-parse','HEAD'],cwd).strip(),'tracked_status':read(['git','status','--short','--untracked-files=no'],cwd),'scripts':{s:hashlib.sha256((cwd/s).read_bytes()).hexdigest() for s in ['scripts/verify_ail.sh','scripts/verify_go.sh']}}
(out/'execution-environment.json').write_text(json.dumps(meta,indent=2)+'\n')
commands=["go test ./cmd/ailang-worldd -run '^TestCLIRealSubprocessEpisode$' -count=1", "go test ./host/capsule -run '^TestF5WallClockTimeoutHasElapsedBound$' -count=1", "go test ./host/coordinator -run '^TestDispatchDurableDeadline$/^Commit$' -count=1", "go test ./host/daemon -run '^TestCommitBudgetAndUncertainReconcile$/^ii_A_landed_uncertain_then_B_on_top$' -count=1"]
def run(label,cmd,cwd):
 started=utc(); t=time.monotonic(); raw=out/(label+'.raw.tmp'); timeout=False
 with raw.open('wb') as f:
  p=subprocess.Popen(['/bin/sh','-c',cmd],cwd=cwd,env=env,stdout=f,stderr=subprocess.STDOUT,start_new_session=True)
  try: rc=p.wait(timeout=1200)
  except subprocess.TimeoutExpired:
   timeout=True; os.killpg(p.pid,signal.SIGKILL); p.wait(); rc=124
 data=raw.read_bytes(); sha=hashlib.sha256(data).hexdigest()
 try:
  text=data.decode('utf-8'); valid=all(ord(c)>=32 or c in '\n\r\t' for c in text)
 except UnicodeDecodeError: valid=False
 if valid and not any(x.rstrip('\r').endswith((' ','\t')) for x in text.split('\n')):
  dest=out/(label+'.log'); dest.write_bytes(data)
 else:
  dest=out/(label+'.log.b64'); dest.write_text(base64.b64encode(data).decode()+'\n'); (out/(label+'.display.log')).write_text(data.decode('utf-8',errors='replace').replace('\x00','<NUL>').rstrip()+'\n')
 raw.unlink()
 result={'command':cmd,'cwd':str(cwd),'start_utc':started,'end_utc':utc(),'duration_seconds':time.monotonic()-t,'rc':rc,'terminal_state':'TIMED_OUT' if timeout else 'COMPLETED','raw_sha256':sha,'raw_bytes':len(data),'raw_artifact':dest.name,'process_group':p.pid,'timeout_cleanup':'SIGKILL process group' if timeout else 'not needed','head':read(['git','rev-parse','HEAD'],cwd).strip(),'pins':{'AILANG_BIN':pin,'WORLD_PKG_AILANG_BIN':pin}}
 (out/(label+'.json')).write_text(json.dumps(result,indent=2)+'\n'); print(json.dumps(result),flush=True)
for label,cwd in [('base',base),('candidate',candidate)]:
 for i,cmd in enumerate(commands,1): run(f'{label}-compare-{i}',cmd,cwd)
(out/'comparator-complete.json').write_text(json.dumps({'completed_utc':utc()})+'\n')
run('verify-ail','./scripts/verify_ail.sh',candidate)
run('verify-go','./scripts/verify_go.sh',candidate)
