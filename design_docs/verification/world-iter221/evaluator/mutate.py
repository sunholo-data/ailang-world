import pathlib,subprocess,hashlib,json,time,os
out=pathlib.Path('/tmp/world-iter221-product-eval-independent'); env=os.environ.copy(); env.update(AILANG_BIN='/Users/voightkampff/.pinned-ailang/ailang',WORLD_PKG_AILANG_BIN='/Users/voightkampff/.pinned-ailang/ailang')
def run(cmd,name,cap=180):
 t=time.time()
 with (out/(name+'.log')).open('w') as f:
  try:r=subprocess.run(cmd,env=env,stdout=f,stderr=subprocess.STDOUT,timeout=cap); rc=r.returncode
  except subprocess.TimeoutExpired:rc=124
 return {'command':cmd,'rc':rc,'seconds':time.time()-t,'cap':cap}
def sha(b):return hashlib.sha256(b).hexdigest()
mutants=[
('absence-catchall','host/projection/projection.go','if !hasHead && errors.As(rerr, &target) {','if !hasHead && (errors.As(rerr, &target) || true) {','./host/projection','^Test(MCPAbsentPrecheckSnapshotFailure|AgentCard_AbsentPrecheckSnapshotFailure|A2A_AbsentPrecheckSnapshotFailure)$'),
('schema-type','host/projection/mcpname.go',"members[\"type\"] = json.RawMessage(`\"object\"`)","members[\"type\"] = json.RawMessage(`\"array\"`)",'./host/projection','^TestMCPSchemaNormalization$'),
('map-length','host/projection/mcpname.go','if len(name) > 64 {','if false && len(name) > 64 {','./host/projection','^Test(MCPNameRefusal|MCPSurfaceRefusalRecovers)$'),
('route-method','host/daemon/daemon.go','mux.HandleFunc(\"POST /mcp/\", d.projection.MCP)','mux.HandleFunc(\"/mcp/\", d.projection.MCP)','./host/daemon','^Test(MCPMountMethodSource|MCPGetRefused)$'),
('late-cancel','host/projection/mcp.go','if err := ctx.Err(); err != nil {\n\t\treturn protocol.InvocationResult{}, err\n\t}\n\tif a.h.coord == nil {','if err := ctx.Err(); false && err != nil {\n\t\treturn protocol.InvocationResult{}, err\n\t}\n\tif a.h.coord == nil {','./host/projection','^TestMCPInvokeCancelAfterAdmission$'),
('capacity-survivor','host/projection/mcp.go','hostcall.New(cfg.CallbackTimeout, cfg.MaxCallbacks)','hostcall.New(cfg.CallbackTimeout, 800)','./host/projection','.'),
('nil-coordinator-survivor','host/projection/mcp.go','return protocol.InvocationResult{}, errors.New(\"projection: invocation coordinator is unavailable\")','return protocol.InvocationResult{}, nil','./host/projection','.'),
('aggregate','host/projection/mcp.go','context.WithTimeout(r.Context(), h.invokeWait)','context.WithCancel(r.Context())','./host/daemon','^TestMCPPostBudgetProductionConstants$')]
results=[]
for name,file,old,new,pkg,pattern in mutants:
 path=pathlib.Path(file); original=path.read_bytes(); s=original.decode(); assert s.count(old)==1,(name,s.count(old)); backup=out/(name+'.backup');backup.write_bytes(original); record={'name':name,'file':file,'old':old,'new':new,'original_sha':sha(original)}
 try:
  path.write_text(s.replace(old,new)); record['mutant_sha']=sha(path.read_bytes()); assert record['mutant_sha']!=record['original_sha']; (out/(name+'.diff')).write_bytes(subprocess.check_output(['git','diff','--',file])); record['build']=run(['go','build','./...'],name+'-build');
  if record['build']['rc']==0:record['test']=run(['go','test',pkg,'-run',pattern,'-count=1','-json','-timeout','150s'],name+'-test'); record['failed_names']=[json.loads(l)['Test'] for l in (out/(name+'-test.log')).read_text().splitlines() if l.startswith('{') and json.loads(l).get('Action')=='fail' and json.loads(l).get('Test')]
 finally:path.write_bytes(backup.read_bytes()); record['restored_sha']=sha(path.read_bytes()); assert record['restored_sha']==record['original_sha']
 record['restored_test']=run(['go','test',pkg,'-run',pattern,'-count=1','-json','-timeout','150s'],name+'-restored');results.append(record);(out/'mutations.json').write_text(json.dumps(results,indent=2));print(name,record.get('failed_names'),flush=True)
print(run(['go','vet','./...'],'vet',180)); print(subprocess.check_output(['git','status','--porcelain'],text=True))
