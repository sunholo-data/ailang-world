import subprocess,os,time,json,signal,datetime,pathlib,sys
root=pathlib.Path('/tmp/world-iter217-repair-round2')
def run(name,args,expected=0):
 env=os.environ.copy();env['AILANG_BIN']='/Users/voightkampff/.pinned-ailang/ailang';env['PATH']='/Users/voightkampff/.pinned-ailang:'+env['PATH'];env.pop('GODEBUG',None)
 start=time.monotonic();utc=datetime.datetime.now(datetime.timezone.utc).isoformat()
 with (root/(name+'.log')).open('wb') as f:
  p=subprocess.Popen(args,stdout=f,stderr=subprocess.STDOUT,env=env,start_new_session=True)
  try:rc=p.wait(timeout=max(1,min(180,(datetime.datetime(2026,10,1,22,27,0,tzinfo=datetime.timezone.utc)-datetime.datetime.now(datetime.timezone.utc)).total_seconds())))
  except subprocess.TimeoutExpired:
   os.killpg(p.pid,signal.SIGKILL);p.wait();rc=124
 row={'command':args,'rc':rc,'expected':expected,'seconds':time.monotonic()-start,'start_utc':utc,'external_cap_s':180,'GODEBUG':None,'AILANG_BIN':env['AILANG_BIN']}
 (root/(name+'.json')).write_text(json.dumps(row,indent=2)+'\n');print(name,rc,flush=True);assert rc==expected,row
if __name__=='__main__':
 run(sys.argv[1],['go','test','-p','1','-timeout=60s']+(['-race'] if 'race' in sys.argv[1] else [])+['./host/projection','-run','^TestMCPInvokeDeadlineFreesSlot$','-v','-count=3'])
