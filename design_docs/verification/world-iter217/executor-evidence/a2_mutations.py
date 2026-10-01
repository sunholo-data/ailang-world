from pathlib import Path
import subprocess,json,hashlib
f='host/projection/projection.go';cases=[
('MUT-ABSENT-PRECHECK-CATCHALL','if !hasHead && errors.As(rerr, &target) {','if !hasHead { _ = target','TestAgentCard_AbsentPrecheckSnapshotFailure/store_error'),
('MUT-ABSENCE-STRING-MATCH','if !hasHead && errors.As(rerr, &target) {','if !hasHead && rerr.Error() == "transition registry: construct request: read transition registry: head is absent" { _ = target','TestProjectionAbsenceRaceControls/same_text_genuine'),
('MUT-PRESENT-ABSENCE-EMPTY','if !hasHead && errors.As(rerr, &target) {','if errors.As(rerr, &target) { _ = hasHead','TestProjectionAbsenceRaceControls/present_then_absent'),
('MUT-PRECHECK-ABSENCE-EARLY-RETURN','\treq, rerr := transitionreg.NewRequest','\tif !hasHead { return transitionreg.Request{}, nil, nil }\n\treq, rerr := transitionreg.NewRequest','TestAgentCard_HeadRaceAbsentThenSucceeds'),
('MUT-PRECHECK-ERROR-IGNORED','if err != nil {\n\t\treturn transitionreg.Request{}, nil, fmt.Errorf("projection: registry head check: %w", err)','if false {\n\t\treturn transitionreg.Request{}, nil, fmt.Errorf("projection: registry head check: %w", err)','TestProjectionAbsenceRaceControls/precheck_error')]
for name,old,new,test in cases[2:]:
 p=Path(f);b=p.read_bytes();assert old in b.decode()
 try:
  p.write_text(b.decode().replace(old,new,1));mut=hashlib.sha256(p.read_bytes()).hexdigest()
  r=subprocess.run(['python3','/tmp/world-iter217-executor/run.py',name+'-fence','go','test','-p','1','./host/projection','-run','^$','-timeout=60s']);assert r.returncode==0
  r=subprocess.run(['python3','/tmp/world-iter217-executor/run.py',name,'go','test','-p','1','./host/projection','-run','^'+test+'$','-v','-count=1','-timeout=60s']);assert r.returncode==1
 finally:p.write_bytes(b)
 r=subprocess.run(['python3','/tmp/world-iter217-executor/run.py',name+'-restored','go','test','-p','1','./host/projection','-run','^'+test+'$','-v','-count=1','-timeout=60s']);assert r.returncode==0
 Path('/tmp/world-iter217-executor/'+name+'-hash.json').write_text(json.dumps({'original_restored_sha256':hashlib.sha256(b).hexdigest(),'mutant_sha256':mut})+'\n')
