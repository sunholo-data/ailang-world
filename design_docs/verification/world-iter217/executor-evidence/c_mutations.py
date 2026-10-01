from pathlib import Path
import subprocess,json,hashlib,re
cases=[
('MUT-RESOLV-DENY-DROPPED',{'host/projection/mcp.go':[('if out.Denied != nil {','if out.Denied != nil { return &authority.SessionBinding{EpisodeID:"denied"},nil\n'),]},'./host/projection','TestMCPDenialMatrix'),
('MUT-KEY-AS-SESSION',{'host/projection/mcp.go':[('r.Header.Get("Authorization")','r.Header.Get("X-API-Key")')]},'./host/projection','TestMCPDenialMatrix'),
('MUT-MCP-PROTECTED',{'host/daemon/daemon.go':[('r.URL.Path == "/v1/commit"','(r.URL.Path == "/v1/commit" || r.URL.Path == "/mcp/")')]},'./host/daemon','TestMCPMountedBearerDenial'),
('MUT-TOOLS-UNFILTERED',{'host/projection/projection.go':[('return req, req.Allowed(), nil','return req, req.Registry.List(), nil')]},'./host/projection','TestMCPListExactSetPerSession'),
('MUT-TOOLS-NO-DEADLINE',{'host/projection/mcp.go':[('context.WithTimeout(ctx, a.h.maxWait)','context.WithTimeout(ctx, 30*time.Second)')]},'./host/projection','TestMCPToolsInnerBudget'),
('MUT-RESOLV-NO-BUDGET',{'host/projection/mcp.go':[('context.WithTimeout(ctx, a.h.credentialWait)','context.WithTimeout(ctx, a.h.invokeWait)')]},'./host/projection','TestMCPResolverInnerBudget'),
('MUT-INVOKE-NO-CTX-CHECK',{'host/projection/mcp.go':[('if err := ctx.Err(); err != nil {\n\t\treturn protocol.InvocationResult{}, err\n\t}\n\tif a.h.coord == nil {','if a.h.coord == nil {')]},'./host/projection','TestMCPInvokeCancelAfterAdmission'),
('MUT-CB-BACKGROUND',{'host/projection/mcp.go':[('context.WithTimeout(ctx, a.h.invokeWait)','context.WithTimeout(context.Background(), a.h.invokeWait)')]},'./host/projection','TestMCPInvokeDeadlineFreesSlot'),
('MUT-INVOKE-NO-READMIT',{
'host/projection/projection.go':[('mintTask       func() (string, error)','mintTask       func() (string, error)\n mutantAdmitted transitionreg.Request')],
'host/projection/mcp.go':[('_, ds, err := a.h.allowedDescriptors(ctx, binding)','admitted, ds, err := a.h.allowedDescriptors(ctx, binding)\n a.h.mutantAdmitted=admitted'),('request, ds, err := a.h.allowedDescriptors(admissionCtx, binding)','_ = admissionCtx\n request:=a.h.mutantAdmitted\n ds:=request.Allowed()\n err = nil')]},'./host/projection','TestMCPInvokeListedThenRevoked'),
('MUT-MUX-METHOD-REMOVED',{'host/daemon/daemon.go':[('"POST /mcp/"','"/mcp/"')]},'./host/daemon','TestMCPMountMethodSource'),
('MUT-ENVELOPE-HANDROLL',{'host/projection/mcp.go':[('h.mcp.ServeHTTP(w, r.WithContext(ctx))','json.NewEncoder(w).Encode(map[string]any{})\n h.mcp.ServeHTTP(w, r.WithContext(ctx))')]},'./host/projection','TestMCPWireOwnershipSource'),
('MUT-DEADLINE-RELAX-MCP',{'host/projection/mcp.go':[('h.mcp.ServeHTTP(w, r.WithContext(ctx))','http.NewResponseController(w).SetWriteDeadline(time.Now().Add(time.Minute))\n h.mcp.ServeHTTP(w, r.WithContext(ctx))')]},'./host/projection','TestMCPWireOwnershipSource'),
('MUT-A2A-MCP-NAME-GATE',{'host/projection/projection.go':[('func (h *Handler) A2A(w http.ResponseWriter, r *http.Request) {','func (h *Handler) A2A(w http.ResponseWriter, r *http.Request) {\n protocol.CallerSurface(nil)')]},'./host/projection','TestProjection_WireOwnershipSource'),
]
cases.append(('MUT-RUNNER-UNBOUNDED',{'host/projection/projection.go':[('cfg.CallbackTimeout <= 0 || cfg.MaxCallbacks <= 0 || ','')],'host/projection/mcp.go':[('hostcall.New(cfg.CallbackTimeout, cfg.MaxCallbacks)','hostcall.New(max(cfg.CallbackTimeout,time.Second), max(cfg.MaxCallbacks,1))')]},'./host/projection','TestMCPBoundsValidation'))
for name,edits,pkg,test in cases:
 originals={f:Path(f).read_bytes() for f in edits};hashes={}
 try:
  for f,changes in edits.items():
   text=originals[f].decode()
   for old,new in changes:assert old in text,(name,old);text=text.replace(old,new,1)
   Path(f).write_text(text);hashes[f]={'mutant':hashlib.sha256(Path(f).read_bytes()).hexdigest(),'restored':hashlib.sha256(originals[f]).hexdigest()}
  args=['go','test','-p','1',pkg,'-timeout=60s']
  fence=subprocess.run(['python3','/tmp/world-iter217-executor/run.py',name+'-fence',*args,'-run','^$']);assert fence.returncode==0
  r=subprocess.run(['python3','/tmp/world-iter217-executor/run.py',name,*args,'-run','^'+test+'$','-v','-count=1']);assert r.returncode==1
  # Broad mutant red sets are measured separately, never inferred from selected tests.
  r=subprocess.run(['python3','/tmp/world-iter217-executor/run.py',name+'-broad',*args,'-v','-count=1']);assert r.returncode==1
  log=Path('/tmp/world-iter217-executor/'+name+'-broad.log').read_text();red=re.findall(r'^\s*--- FAIL: (\S+)',log,re.M)
 finally:
  for f,b in originals.items():Path(f).write_bytes(b)
 r=subprocess.run(['python3','/tmp/world-iter217-executor/run.py',name+'-restored','go','test','-p','1',pkg,'-run','^'+test+'$','-v','-count=1','-timeout=60s']);assert r.returncode==0
 Path('/tmp/world-iter217-executor/'+name+'-hash.json').write_text(json.dumps({'files':hashes,'broad_failed_tests':red,'selected_test':test})+'\n')
