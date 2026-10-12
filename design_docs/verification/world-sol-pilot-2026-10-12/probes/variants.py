import subprocess,sys,os,time,json
sys.path.insert(0,os.path.join(os.path.dirname(os.path.abspath(__file__)),'..','wt','scripts','floor'))
import canary
RIG='/Users/voightkampff/dev/sunholo-data/ailang-world/.ailang/state/row174-attended-2026-10-10/rig'
HERE=os.path.dirname(os.path.abspath(__file__)); PORT=7699
L=canary.layout(RIG); tok=open(L['private']+'/fl-tune-cx.session').read().strip()
BASE=['codex','exec','--json','--ephemeral','--ignore-user-config','--skip-git-repo-check','-m','gpt-6.1-sol','-s','read-only','-C',os.path.join(L['void'],'fl-tune-cx'),
 '--disable','plugins','--disable','apps','--disable','browser_use','--disable','computer_use','--disable','image_generation','--disable','multi_agent','--disable','shell_tool','--disable','unified_exec','-c','web_search=disabled',
 '-c','mcp_servers.world.url=http://127.0.0.1:7644/mcp/','-c','mcp_servers.world.bearer_token_env_var=WORLD_SESSION','-c','mcp_servers.world.required=true','-c','mcp_servers.world.default_tools_approval_mode=approve',
 '-c','model_provider=probe','-c','model_providers.probe.name=probe','-c',f'model_providers.probe.base_url=http://127.0.0.1:{PORT}/v1','-c','model_providers.probe.env_key=PROBE_KEY','-c','model_providers.probe.wire_api=responses','-c','model_providers.probe.request_max_retries=0','-c','model_providers.probe.stream_max_retries=0']
V={'no-code-mode-host':['--disable','code_mode_host'],'no-defer':['--disable','tool_search_always_defer_mcp_tools'],'both':['--disable','code_mode_host','--disable','tool_search_always_defer_mcp_tools']}
d=canary.start_daemon(L,'127.0.0.1:7644')
try:
  for name,extra in V.items():
    out=os.path.join(HERE,'v-'+name); os.makedirs(out,exist_ok=True)
    m=subprocess.Popen([sys.executable,os.path.join(HERE,'mock.py'),out,str(PORT)],env={k:v for k,v in os.environ.items() if k!='MOCK_EXEC'}); time.sleep(1)
    env=dict(os.environ,PROBE_KEY='x',WORLD_SESSION=tok)
    r=subprocess.run(BASE[:-0 or None]+extra+['hi'],env=env,capture_output=True,text=True,timeout=60,stdin=subprocess.DEVNULL)
    m.terminate(); m.wait()
    open(out+'/codex.stderr','w').write(r.stderr.replace(tok,'<TOKEN>'))
    f=out+'/req-0_v1_responses.json'
    if not os.path.exists(f): print(name,'rc',r.returncode,'NO REQUEST',r.stderr[-300:]); continue
    b=json.load(open(f)); s=json.dumps(b)
    names=[]
    for it in b['input']:
      for ns in it.get('tools') or []:
        names+= [ (ns.get('name','')+'.'+t.get('name','')) for t in ns.get('tools',[])] or [ns.get('name')]
    for t in b.get('tools') or []: names.append('top.'+str(t.get('name') or t.get('type')))
    print(name,'rc',r.returncode,'bytes',len(s),'world-tool schemas visible:',sum('ailang' in n for n in names),names)
finally:
  canary.stop_daemon(d)
