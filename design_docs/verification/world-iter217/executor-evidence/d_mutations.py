from pathlib import Path
import subprocess,json,hashlib,re
p='host/projection/projection.go';m='host/projection/mcp.go';t='host/projection/mcp_test.go';d='host/projection/mcp_conformance_test.go'
cached={p:[('mintTask       func() (string, error)','mintTask       func() (string, error)\n mutantAdmitted transitionreg.Request')],m:[('_, ds, err := a.h.allowedDescriptors(ctx, binding)','admitted, ds, err := a.h.allowedDescriptors(ctx, binding)\n a.h.mutantAdmitted=admitted'),('request, ds, err := a.h.allowedDescriptors(admissionCtx, binding)','_ = admissionCtx\n request:=a.h.mutantAdmitted\n ds:=request.Allowed()\n err = nil')]}
cases=[
('MUT-TASK-PER-POST',{p:[('mintTask       func() (string, error)','mintTask       func() (string, error)\n mutantTask string')],m:[('task, err := a.h.mintTask()','task:=a.h.mutantTask\n if task=="" {task,err=a.h.mintTask();a.h.mutantTask=task}')]},'./host/projection','TestMCPBatchAccounting/repeated_rpc_id','production'),
('MUT-SNAPSHOT-PER-POST',cached,'./host/projection','TestMCPBatchAccounting/repeated_rpc_id','production'),
('MUT-BATCH-NO-READMIT',cached,'./host/projection','TestMCPBatchAccounting/removed_between_items','production'),
('MUT-SURFACE-ERROR-CACHED',{p:[('mintTask       func() (string, error)','mintTask       func() (string, error)\n mutantSurfaceError error')],m:[('tools, err := mcpDescriptors(ds)','if a.h.mutantSurfaceError!=nil {return nil,a.h.mutantSurfaceError}\n tools, err := mcpDescriptors(ds)\n if err!=nil {a.h.mutantSurfaceError=err}')]},'./host/projection','TestMCPSurfaceRefusalRecovers','production'),
('MUT-BATCH-AS-SINGLE',{d:[('r.Header.Set("MCP-Protocol-Version", version)','r.Header.Set("MCP-Protocol-Version", "2025-06-18")\n _ =version')]},'./host/projection','TestMCPBatchConformance/old_versions_and_item_errors','test decorator'),
('MUT-BATCH-PARTIAL',{d:[('var wire struct {\n\t\t\tID    any','w.Header().Set("Content-Type","text/event-stream")\n w.Body.Reset(); w.Body.WriteString("event: message\\ndata: [{\\\"jsonrpc\\\":\\\"2.0\\\",\\\"id\\\":7,\\\"result\\\":{}}]\\n\\n")\n var wire struct {\n\t\t\tID    any')]},'./host/projection','TestMCPBatchConformance/whole_host_failure','test response decorator'),
('MUT-PLAIN-JSON',{d:[('if !strings.HasPrefix(got.Body.String(),','if !strings.HasPrefix(got.Body.String(),')]},'./host/projection','TestMCPWireConformance','test response decorator'),
('MUT-CROSS-SURFACE-DIVERGE',{m:[('tools, err := mcpDescriptors(ds)','ds=ds[:0]\n tools, err := mcpDescriptors(ds)')]},'./host/daemon','TestMCPCrossSurfaceExactSet','production'),
('MUT-AMBIENT-HARDCODED',{m:[('return tools, err','tools=append(tools,protocol.ToolDescriptor{Name:"exit",InputSchema:json.RawMessage(`{"type":"object"}`)})\n return tools, err')]},'./host/daemon','TestMCPAmbientExportsAbsent','production'),
]
# Wire calibration transforms observed response, without inventing a World codec.
s=Path(d).read_text();needle='got := httptest.NewRecorder()';print('plain needle?',needle in s)
# defer plain separately until exact source inspected
cases=[c for c in cases if c[0]!='MUT-PLAIN-JSON']
for name,edits,pkg,test,kind in cases:
 originals={f:Path(f).read_bytes() for f in edits};hashes={}
 try:
  for f,changes in edits.items():
   text=originals[f].decode()
   for old,new in changes:assert old in text,(name,old);text=text.replace(old,new,1)
   Path(f).write_text(text);hashes[f]={'mutant':hashlib.sha256(Path(f).read_bytes()).hexdigest(),'restored':hashlib.sha256(originals[f]).hexdigest()}
  args=['go','test','-p','1',pkg,'-timeout=60s']
  for suffix,run,expected in [('fence','^$',0),('', '^'+test+'$',1)]:
   r=subprocess.run(['python3','/tmp/world-iter217-executor/run.py',name+('-'+suffix if suffix else ''),*args,'-run',run,'-v','-count=1']);assert r.returncode==expected,(name,suffix,r.returncode)
  # Production-budget test is excluded from broad diagnostic here and recorded explicitly.
  run='TestMCP' if pkg.endswith('projection') else 'TestMCP(Cross|Ambient|Mount)'
  r=subprocess.run(['python3','/tmp/world-iter217-executor/run.py',name+'-broad',*args,'-run',run,'-v','-count=1']);assert r.returncode==1
  red=re.findall(r'^\s*--- FAIL: (\S+)',Path('/tmp/world-iter217-executor/'+name+'-broad.log').read_text(),re.M)
 finally:
  for f,b in originals.items():Path(f).write_bytes(b)
 r=subprocess.run(['python3','/tmp/world-iter217-executor/run.py',name+'-restored',*args,'-run','^'+test+'$','-v','-count=1']);assert r.returncode==0
 Path('/tmp/world-iter217-executor/'+name+'-hash.json').write_text(json.dumps({'files':hashes,'selected_test':test,'kind':kind,'broad_scope':run,'broad_failed_tests':red})+'\n')
