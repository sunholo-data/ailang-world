import subprocess,sys,os,time,json,shutil,tempfile
sys.path.insert(0,os.path.join(os.path.dirname(__file__),'..','wt','scripts','floor'))
import canary
JS={'shell':'const T={};const tm=async(k,f)=>{const a=Date.now();await f();(T[k]=T[k]||[]).push(Date.now()-a)};for(let i=0;i<5;i++){await tm(\'cat\',()=>tools.exec_command({cmd:\'cat benchmark/solution.ail\'}));await tm(\'write\',()=>tools.exec_command({cmd:\'printf "module benchmark/solution\\\\n\\\\nexport func main() -> () ! {IO} = println(show(0.0 == 0.0))\\\\n" > benchmark/solution.ail\'}));await tm(\'check\',()=>tools.exec_command({cmd:\'ailang check benchmark/solution.ail\'}));await tm(\'run\',()=>tools.exec_command({cmd:\'ailang run --entry main --caps IO benchmark/solution.ail\'}));}text(JSON.stringify(T))','world':"const T={};const tm=async(k,f)=>{const a=Date.now();await f();(T[k]=T[k]||[]).push(Date.now()-a)};for(let i=0;i<5;i++){await tm('read',()=>tools.mcp__world__ailang_read({path:'benchmark/solution.ail'}));await tm('write',()=>tools.mcp__world__ailang_write({path:'benchmark/solution.ail',content:'module benchmark/solution\\n\\nexport func main() -> () ! {IO} = println(show(0.0 == 0.0))\\n'}));await tm('check',()=>tools.mcp__world__ailang_check({path:'benchmark/solution.ail'}));await tm('run',()=>tools.mcp__world__ailang_run({path:'benchmark/solution.ail',caps:['IO']}));}text(JSON.stringify(T))"}
RIG='/Users/voightkampff/dev/sunholo-data/ailang-world/.ailang/state/row174-attended-2026-10-10/rig'
HERE=os.path.dirname(os.path.abspath(__file__)); PORT=7699
COMMON=['codex','exec','--json','--ephemeral','--ignore-user-config','--skip-git-repo-check','-m','gpt-6.1-sol']
DIS=['--disable','plugins','--disable','apps','--disable','browser_use','--disable','computer_use','--disable','image_generation','--disable','multi_agent']
MOCK=['-c','model_provider=probe','-c','model_providers.probe.name=probe','-c',f'model_providers.probe.base_url=http://127.0.0.1:{PORT}/v1',
      '-c','model_providers.probe.env_key=PROBE_KEY','-c','model_providers.probe.wire_api=responses','-c','model_providers.probe.request_max_retries=0','-c','model_providers.probe.stream_max_retries=0']
L=canary.layout(RIG); tok=open(L['private']+'/fl-tune-cx.session').read().strip()
sd=tempfile.mkdtemp(dir=HERE); os.makedirs(sd+'/benchmark'); open(sd+'/benchmark/solution.ail','w').write('module benchmark/solution\n')
arms={'shell':COMMON+['-s','workspace-write','-C',sd]+DIS+['-c','web_search=disabled']+MOCK+['hi'],
      'world':COMMON+['-s','read-only','-C',os.path.join(L['void'],'fl-tune-cx')]+DIS+['--disable','shell_tool','--disable','unified_exec','-c','web_search=disabled',
        '-c','mcp_servers.world.url=http://127.0.0.1:7644/mcp/','-c','mcp_servers.world.bearer_token_env_var=WORLD_SESSION',
        '-c','mcp_servers.world.required=true','-c','mcp_servers.world.default_tools_approval_mode=approve']+MOCK+['hi']}
d=canary.start_daemon(L,'127.0.0.1:7644')
try:
  for arm,argv in arms.items():
    out=os.path.join(HERE,arm+'-timing'); os.makedirs(out,exist_ok=True)
    os.environ['MOCK_EXEC']=JS[arm]
    m=subprocess.Popen([sys.executable,os.path.join(HERE,'mock.py'),out,str(PORT)]); time.sleep(1)
    env=dict(os.environ,PROBE_KEY='not-a-key',MOCK_EXEC='text(JSON.stringify(ALL_TOOLS.map(t=>({name:t.name,dlen:(t.description||"").length}))))'); env.pop('OPENAI_API_KEY',None)
    if arm=='world': env['WORLD_SESSION']=tok
    t0=time.time(); r=subprocess.run(argv,env=env,capture_output=True,text=True,timeout=120,stdin=subprocess.DEVNULL)
    m.terminate(); m.wait()
    open(os.path.join(out,'codex.stdout'),'w').write(r.stdout); open(os.path.join(out,'codex.stderr'),'w').write(r.stderr.replace(tok,'<TOKEN>'))
    print(arm,'rc',r.returncode,'%.1fs'%(time.time()-t0),'requests',sorted(os.listdir(out)))
finally:
  canary.stop_daemon(d); shutil.rmtree(sd)
