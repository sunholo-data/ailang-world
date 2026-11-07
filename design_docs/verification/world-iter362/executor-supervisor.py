#!/usr/bin/env python3
import datetime, hashlib, json, os, re, signal, subprocess, sys, time
from pathlib import Path
wt=Path('/Users/voightkampff/dev/sunholo-data/.exec-world-iter362-row256')
v=wt/'design_docs/verification/world-iter362'
mode=sys.argv[1]; assert mode in ('probe','run')
name='executor-probe' if mode=='probe' else 'executor'
directive=v/(name+'-directive.txt'); prompt=directive.read_text()
assert prompt.startswith('MISSION-ROLE: executor') and len(prompt.encode())>=200
seconds=120 if mode=='probe' else 1800
cache=subprocess.check_output(['go','env','GOCACHE'],text=True).strip()
modcache=subprocess.check_output(['go','env','GOMODCACHE'],text=True).strip()
cmd=['/opt/homebrew/bin/codex','exec','--model','gpt-6.1-sol','-c','features.multi_agent_v2.expose_spawn_agent_model_overrides=false','--sandbox','workspace-write','--add-dir',cache,'--add-dir',modcache,'-C',str(wt),'-o',str(v/(name+'-final.txt')),prompt]
env=dict(os.environ); env.pop('OPENAI_API_KEY',None)
start=time.monotonic(); utc=datetime.datetime.now(datetime.timezone.utc).isoformat(); timeout=False
log=v/(name+'-codex.txt')
with log.open('w') as out:
 p=subprocess.Popen(cmd,cwd=wt,env=env,stdin=subprocess.DEVNULL,stdout=out,stderr=subprocess.STDOUT,start_new_session=True)
 (v/(name+'-pid.json')).write_text(json.dumps({'pid':p.pid,'pgid':p.pid,'deadline_seconds':seconds,'utc':utc})+'\n')
 try: rc=p.wait(timeout=seconds-5)
 except subprocess.TimeoutExpired:
  timeout=True
  try: os.killpg(p.pid,signal.SIGTERM)
  except ProcessLookupError: pass
  try: p.wait(timeout=2)
  except subprocess.TimeoutExpired:
   try: os.killpg(p.pid,signal.SIGKILL)
   except ProcessLookupError: pass
   p.wait(timeout=2)
  rc=124
 teardown=[]
 for sig in (signal.SIGTERM,signal.SIGKILL):
  try: os.killpg(p.pid,sig); teardown.append(sig.name+' sent')
  except ProcessLookupError: teardown.append(sig.name+' group absent')
raw=log.read_text(); m=re.search(r'^model:\s*(.+)$',raw,re.M)
t=re.search(r'tokens used\s*\n([0-9.,]+)',raw)
receipt={'utc':utc,'requested_model':'gpt-6.1-sol','model_header':m.group(1).strip() if m else None,'rc':rc,'timeout':timeout,'elapsed_seconds':round(time.monotonic()-start,2),'directive_path':str(directive.relative_to(wt)),'directive_sha256':hashlib.sha256(directive.read_bytes()).hexdigest(),'directive_bytes':len(prompt.encode()),'stdin':'devnull','api_key_stripped':True,'tokens_used_cli_raw':t.group(1) if t else None,'tokens_used':int(re.sub(r'[^0-9]','',t.group(1))) if t else None,'process_group':p.pid,'process_group_teardown':teardown,'log_path':str(log.relative_to(wt)),'add_dirs':[cache,modcache],'command_prefix':cmd[:-1]}
(v/(name+'-cli-receipt.json')).write_text(json.dumps(receipt,indent=2)+'\n')
(v/(name+'-codex.rc')).write_text(str(rc)+'\n')
print(json.dumps(receipt,indent=2),flush=True)
sys.exit(rc)
