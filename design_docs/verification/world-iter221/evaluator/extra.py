exec(open('/tmp/world-iter221-product-eval-independent/mutate.py').read().split('mutants=[')[0])
results=[]
for name,file,old,new,pkg,pattern in [
('typed-absence-emission','host/transitionreg/transitionreg.go','return Snapshot{}, RegistryHeadAbsentError{}','return Snapshot{}, errors.New("read transition registry: head is absent")','./host/transitionreg','^TestRegistryHeadAbsentErrorClassification$'),
('fresh-allowed-guard','host/projection/mcp.go','if !allowed {','if false && !allowed {','./host/projection','^Test(MCPInvokeListedThenRevoked|MCPBatchAccounting|MCPListedThenTrulyAbsent)$')]:
 path=pathlib.Path(file);original=path.read_bytes();s=original.decode();assert s.count(old)==1;backup=out/(name+'.backup');backup.write_bytes(original);record={'name':name,'file':file,'old':old,'new':new,'original_sha':sha(original)}
 try:
  path.write_text(s.replace(old,new));record['mutant_sha']=sha(path.read_bytes());assert record['mutant_sha']!=record['original_sha'];(out/(name+'.diff')).write_bytes(subprocess.check_output(['git','diff','--',file]));record['build']=run(['go','build','./...'],name+'-build')
  if record['build']['rc']==0:record['test']=run(['go','test',pkg,'-run',pattern,'-count=1','-json','-timeout','150s'],name+'-test');record['failed_names']=[json.loads(l)['Test'] for l in (out/(name+'-test.log')).read_text().splitlines() if l.startswith('{') and json.loads(l).get('Action')=='fail' and json.loads(l).get('Test')]
 finally:path.write_bytes(backup.read_bytes());record['restored_sha']=sha(path.read_bytes());assert record['restored_sha']==record['original_sha']
 record['restored_test']=run(['go','test',pkg,'-run',pattern,'-count=1','-json','-timeout','150s'],name+'-restored');results.append(record);(out/'extra-mutations.json').write_text(json.dumps(results,indent=2));print(name,record.get('failed_names'),flush=True)
