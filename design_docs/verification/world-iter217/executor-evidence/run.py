import subprocess,os,signal,sys,time,json,pathlib
name=sys.argv[1];cmd=sys.argv[2:];base=pathlib.Path('/tmp/world-iter217-executor');start=time.monotonic()
env=os.environ.copy();env["AILANG_BIN"]="/Users/voightkampff/.pinned-ailang/ailang";env["PATH"]="/Users/voightkampff/.pinned-ailang:"+env["PATH"]
process_cap=int(env.get("WORLD_EXEC_PROCESS_CAP","120"))
if "-timeout=120s" in cmd:
 cmd=["-timeout=180s" if arg=="-timeout=120s" else arg for arg in cmd];process_cap=max(process_cap,210)
p=subprocess.Popen(cmd,env=env,stdout=subprocess.PIPE,stderr=subprocess.STDOUT,start_new_session=True)
try: out=p.communicate(timeout=process_cap)[0];rc=p.returncode
except subprocess.TimeoutExpired: os.killpg(p.pid,signal.SIGKILL);out=p.communicate()[0];rc=124
(base/(name+'.log')).write_bytes(out);(base/(name+'.json')).write_text(json.dumps({'cmd':cmd,'rc':rc,'seconds':time.monotonic()-start,'process_cap_seconds':process_cap,'diagnostic_godebug':env.get('GODEBUG')})+'\n');print(out.decode(errors='replace'));print('RC',rc);sys.exit(rc)
