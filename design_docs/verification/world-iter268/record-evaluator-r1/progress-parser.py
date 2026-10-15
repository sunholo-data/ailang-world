"""Read-only reproduction using authoritative weekly-report functions; no module side effects."""
from pathlib import Path
import ast,json,os,re,hashlib,subprocess
ROOT=Path.cwd();RAW=ROOT/'design_docs/verification/world-iter268/record-evaluator-r1'
p=Path('/Users/voightkampff/dev/sunholo-data/ailang/tools/mission-weekly-report.py')
tree=ast.parse(p.read_text()); selected=[x for x in tree.body if (isinstance(x,ast.FunctionDef) and x.name in {'clean','goal_distance'}) or (isinstance(x,ast.Assign) and any(isinstance(t,ast.Name) and t.id=='PROG_RE' for t in x.targets))]
ns={'re':re,'os':os};exec(compile(ast.Module(body=selected,type_ignores=[]),str(p),'exec'),ns)
log=ROOT/'design_docs/world-mission-log.md';text=log.read_text();entry=text[text.index('## 268 —'):];line=next(l for l in entry.splitlines() if l.startswith('**Progress'))
control=line.replace('**Progress:**','**Progress**:',1);matches=[];head=None
for i,l in enumerate(text.splitlines(),1):
 if l.startswith('## '):head=l
 if ns['PROG_RE'].match(l):matches.append({'line':i,'heading':head,'value':ns['PROG_RE'].match(l).group(1)})
control_path=RAW/'progress-label-control.txt';control_path.write_text(control+'\n')
result={'head':subprocess.check_output(['git','rev-parse','HEAD'],text=True,timeout=15).strip(),'authored_line':line,'current_entry_line':next(i for i,l in enumerate(text.splitlines(),1) if l==line),'pattern':ns['PROG_RE'].pattern,'authored_matches':bool(ns['PROG_RE'].match(line)),'minimal_corrected_control':control,'control_matches':bool(ns['PROG_RE'].match(control)),'value_unchanged':line.split(' ',1)[1]==control.split(' ',1)[1],'actual_unmodified_goal_distance_function_result':ns['goal_distance'](log),'actual_latest_matching_record':matches[-1],'control_goal_distance_result':ns['goal_distance'](control_path),'function_source_sha256':hashlib.sha256(p.read_bytes()).hexdigest()}
(RAW/'progress-parser.json').write_text(json.dumps(result,indent=2)+'\n');print(json.dumps(result,indent=2))
