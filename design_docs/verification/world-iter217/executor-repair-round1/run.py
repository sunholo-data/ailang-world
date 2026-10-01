import os,sys,subprocess,signal,time,json,pathlib,datetime
base=pathlib.Path('/tmp/world-iter217-executor-repair-round1');name=sys.argv[1];cmd=sys.argv[2:];env=os.environ.copy();env['AILANG_BIN']='/Users/voightkampff/.pinned-ailang/ailang';env['PATH']='/Users/voightkampff/.pinned-ailang:'+env['PATH'];start=time.monotonic();utc=datetime.datetime.now(datetime.timezone.utc).isoformat();cap=max(1,min(180,int((datetime.datetime.fromisoformat("2026-10-01T21:47:30+00:00")-datetime.datetime.now(datetime.timezone.utc)).total_seconds())))
p=subprocess.Popen(cmd,env=env,stdout=subprocess.PIPE,stderr=subprocess.STDOUT,start_new_session=True)
try:out=p.communicate(timeout=cap)[0];rc=p.returncode
except subprocess.TimeoutExpired:os.killpg(p.pid,signal.SIGKILL);out=p.communicate()[0];rc=124
(base/(name+'.log')).write_bytes(out);(base/(name+'.json')).write_text(json.dumps({'cmd':cmd,'rc':rc,'seconds':time.monotonic()-start,'started_utc':utc,'external_cap_seconds':cap,'godebug':env.get('GODEBUG'),'AILANG_BIN':env['AILANG_BIN']},indent=2)+'\n');print(out.decode(errors='replace'));print('RC',rc);sys.exit(rc)
