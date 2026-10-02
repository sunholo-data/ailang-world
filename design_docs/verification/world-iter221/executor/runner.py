import os,sys,subprocess,signal,time,json,pathlib,datetime
root=pathlib.Path(__file__).parent
pin="/Users/voightkampff/.pinned-ailang/ailang"
env=os.environ.copy();env.update(AILANG_BIN=pin,WORLD_PKG_AILANG_BIN=pin)
def utc(): return datetime.datetime.now(datetime.timezone.utc).isoformat()
for tree in ["base","verify"]:
 for gate in ["ail","go"]:
  key=tree+"-"+gate;cwd="/Users/voightkampff/dev/sunholo-data/.wt-world-iter221-"+tree
  m={"command":"./scripts/verify_"+gate+".sh","cwd":cwd,"start_utc":utc(),"start_epoch":time.time(),"ceiling_s":1200,"head":subprocess.check_output(["git","rev-parse","HEAD"],cwd=cwd,text=True).strip(),"AILANG_BIN":pin,"WORLD_PKG_AILANG_BIN":pin}
  with (root/(key+".stdout.log")).open("wb") as out,(root/(key+".stderr.log")).open("wb") as err:
   p=subprocess.Popen([m["command"]],cwd=cwd,env=env,stdout=out,stderr=err,start_new_session=True);m.update(pid=p.pid,pgid=p.pid,deadline_epoch=m["start_epoch"]+1200)
   (root/(key+".json")).write_text(json.dumps(m,indent=2));print("START",key,p.pid,flush=True)
   try: rc=p.wait(timeout=max(0,m["deadline_epoch"]-time.time()))
   except subprocess.TimeoutExpired:
    os.killpg(p.pid,signal.SIGKILL);p.wait();rc=124;m["timeout"]=True
  m.update(rc=rc,end_utc=utc(),end_epoch=time.time());(root/(key+".json")).write_text(json.dumps(m,indent=2));print("END",key,rc,flush=True)
  if rc: sys.exit(rc)
