"""Probe 1: World MCP per-call latency vs the same operations natively. No LLM."""
import json,os,sys,time,subprocess,statistics as st,shutil,tempfile
sys.path.insert(0,'wt/scripts/floor')
import canary,world
RIG='/Users/voightkampff/dev/sunholo-data/ailang-world/.ailang/state/row174-attended-2026-10-10/rig'
TOOL=os.path.expanduser('~/.pinned-ailang-tools/v0.52.1/ailang')
ADDR='127.0.0.1:7644'; N=10
SRC="module benchmark/solution\n\nexport func main() -> () ! {IO} = println(show(0.0 == 0.0))\n"
L=canary.layout(RIG); tok=open(L['private']+'/fl-tune-cx.session').read().strip()
res={'world':{},'native':{}}
def t(bucket,name,fn):
    t0=time.perf_counter(); r=fn(); res[bucket].setdefault(name,[]).append(time.perf_counter()-t0); return r
def call(name,args): 
    r=c.mcp(tok,'tools/call',{'name':name,'arguments':args})
    txt=r['result']['content'][0]['text']; assert '"ok":true' in txt or '"admitted":true' in txt, txt[:300]; return txt
d=canary.start_daemon(L,ADDR)
try:
    c=world.Client(ADDR)
    for i in range(N):
        t('world','initialize',lambda: c.mcp(tok,'initialize',{'protocolVersion':'2025-06-18','capabilities':{},'clientInfo':{'name':'probe1','version':'0'}}))
        t('world','tools/list',lambda: c.mcp(tok,'tools/list'))
        t('world','ailang-read',lambda: call('ailang-read',{'path':'benchmark/solution.ail'}))
        t('world','ailang-write',lambda: call('ailang-write',{'path':'benchmark/solution.ail','content':SRC}))
        t('world','ailang-check',lambda: call('ailang-check',{'path':'benchmark/solution.ail'}))
        out=t('world','ailang-run',lambda: call('ailang-run',{'path':'benchmark/solution.ail','caps':['IO']}))
finally:
    canary.stop_daemon(d)
wd=tempfile.mkdtemp(dir=os.path.dirname(os.path.abspath(__file__))+'/probe1'); os.makedirs(wd+'/benchmark')
p=wd+'/benchmark/solution.ail'; open(p,'w').write(SRC)
for i in range(N):
    t('native','read',lambda: open(p).read())
    t('native','write',lambda: open(p,'w').write(SRC))
    t('native','ailang check',lambda: subprocess.run([TOOL,'check','benchmark/solution.ail'],cwd=wd,capture_output=True,check=True))
    t('native','ailang run',lambda: subprocess.run([TOOL,'run','--caps','IO','benchmark/solution.ail'],cwd=wd,capture_output=True,check=True))
shutil.rmtree(wd)
summary={b:{k:{'median_ms':round(1000*st.median(v),1),'min_ms':round(1000*min(v),1),'max_ms':round(1000*max(v),1)} for k,v in res[b].items()} for b in res}
summary['last_run_output_head']=out[:200]
json.dump(summary,open('probe1/result.json','w'),indent=1); print(json.dumps(summary,indent=1))
