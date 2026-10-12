import subprocess,sys,os,time,json
HERE=os.path.dirname(os.path.abspath(__file__))
sys.path.insert(0,os.path.join(HERE,'..','wt','scripts','floor'))
import canary,run_task
RIG='/Users/voightkampff/dev/sunholo-data/ailang-world/.ailang/state/row174-attended-2026-10-10/rig'
L=canary.layout(RIG); tok=open(L['private']+'/fl-tune-cx.session').read().strip(); PORT=7699
argv=['codex','exec','--json','--ephemeral','--ignore-user-config','--skip-git-repo-check','-m','gpt-6.1-sol','-s','read-only','-C',os.path.join(L['void'],'fl-tune-cx'),
 '--disable','plugins','--disable','apps','--disable','browser_use','--disable','computer_use','--disable','image_generation','--disable','multi_agent','--disable','shell_tool','--disable','unified_exec','-c','web_search=disabled',
 '-c','mcp_servers.world.url=http://127.0.0.1:7644/mcp/','-c','mcp_servers.world.bearer_token_env_var=WORLD_SESSION','-c','mcp_servers.world.required=true','-c','mcp_servers.world.default_tools_approval_mode=approve',
 '-c','model_provider=probe','-c','model_providers.probe.name=probe','-c',f'model_providers.probe.base_url=http://127.0.0.1:{PORT}/v1','-c','model_providers.probe.env_key=PROBE_KEY','-c','model_providers.probe.wire_api=responses','-c','model_providers.probe.request_max_retries=0','-c','model_providers.probe.stream_max_retries=0','hi']
out=os.path.join(HERE,'spawn-timing'); os.makedirs(out,exist_ok=True)
menv=dict(os.environ,MOCK_DELAY='1.5',MOCK_EXEC='const r = await tools.mcp__world__ailang_read({path: "benchmark/solution.ail"}); text("ok")')
d=canary.start_daemon(L,'127.0.0.1:7644')
try:
    m=subprocess.Popen([sys.executable,os.path.join(HERE,'mock_delay.py'),out,str(PORT)],env=menv); time.sleep(1)
    tr=os.path.join(out,'t.jsonl')
    r=run_task.spawn(argv,cwd=os.path.join(L['void'],'fl-tune-cx'),env=dict(os.environ,PROBE_KEY='x',WORLD_SESSION=tok),transcript_path=tr,deadline_s=60)
    m.terminate(); m.wait()
finally:
    canary.stop_daemon(d)
times=json.load(open(tr+'.times')); text=open(tr).read()
for ln,t in zip([l for l in text.split('\n') if l.strip()],times): print(t, ln[:110])
print(r); print(run_task.event_timing(text,times,r['wall_ms']))
