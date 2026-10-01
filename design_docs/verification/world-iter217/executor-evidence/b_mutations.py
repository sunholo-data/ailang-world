from pathlib import Path
import subprocess,json,hashlib
f='host/projection/mcpname.go';s=Path(f).read_text();start=s.index('func normalizeMCPSchema');end=s.index('func mcpDescriptors',start)
floatfn='''func normalizeMCPSchema(raw []byte) (json.RawMessage,error) {
 var members map[string]any
 if err:=json.Unmarshal(raw,&members);err!=nil||members==nil{return nil,errors.New("schema")}
 if typ,ok:=members["type"];ok{if typ!="object"{return nil,errors.New("schema")};return raw,nil}
 members["type"]="object";return json.Marshal(members)
}
'''
cases=[
('MUT-MAP-IDENTITY','name := strings.NewReplacer("_", "_u", ".", "_d", "/", "_s").Replace(id)','name := id','TestMCPNameRoundTrip'),
('MUT-MAP-COLLIDE','"_", "_u"','"_", "_d"','TestMCPNameRoundTrip'),
('MUT-MAP-DECODE-INVALID','return "", errors.New("projection: invalid MCP escape")','return "accepted", nil','TestMCPNameRefusal'),
('MUT-MAP-NO-LEN-GUARD','if len(name) > 64 {','if false {','TestMCPNameRefusal'),
('MUT-SCHEMA-PASSTHRU','members["type"] = json.RawMessage(`"object"`)','return raw,nil\n members["type"] = json.RawMessage(`"object"`)','TestMCPSchemaNormalization'),
('MUT-SCHEMA-DROP-CONSTRAINT','return json.Marshal(members)','return json.RawMessage(`{"type":"object"}`),nil','TestMCPSchemaNormalization'),
('MUT-SCHEMA-FLOAT64',s[start:end],floatfn,'TestMCPSchemaNormalization')]
for name,old,new,test in cases:
 p=Path(f);b=p.read_bytes();assert old in b.decode()
 try:
  p.write_text(b.decode().replace(old,new,1));mut=hashlib.sha256(p.read_bytes()).hexdigest()
  r=subprocess.run(['python3','/tmp/world-iter217-executor/run.py',name+'-fence','go','test','-p','1','./host/projection','-run','^$','-timeout=60s']);assert r.returncode==0
  r=subprocess.run(['python3','/tmp/world-iter217-executor/run.py',name,'go','test','-p','1','./host/projection','-run','^'+test+'$','-v','-count=1','-timeout=60s']);assert r.returncode==1
 finally:p.write_bytes(b)
 r=subprocess.run(['python3','/tmp/world-iter217-executor/run.py',name+'-restored','go','test','-p','1','./host/projection','-run','^'+test+'$','-v','-count=1','-timeout=60s']);assert r.returncode==0
 Path('/tmp/world-iter217-executor/'+name+'-hash.json').write_text(json.dumps({'original_restored_sha256':hashlib.sha256(b).hexdigest(),'mutant_sha256':mut})+'\n')
