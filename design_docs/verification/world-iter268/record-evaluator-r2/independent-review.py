from pathlib import Path
import subprocess,json,hashlib,re,collections,ast,os,datetime
R=Path.cwd(); D=R/'design_docs/verification/world-iter268'; O=D/'record-evaluator-r2'; B='420c07e7c1c70eca6043a0c4e12a20f2adcaf972'; C='b20bb350cf1d60c9020b4a831be0ae26ad915246'; P='d6e04632851bb3ebd31c2ef5c4d32a114bb6a1ca'; S='763c60a20358f486df395043131557c1143735ba'
def g(*a):return subprocess.check_output(['git',*a],timeout=30)
def b(rev,p):return g('show',rev+':'+p)
def n(p):return (R/p).read_bytes()
def digest(v):return hashlib.sha256(v).hexdigest()
checks=[]
def yes(name,ok,data=None):checks.append({'name':name,'pass':bool(ok),'data':data})
yes('exact-review-head',g('rev-parse','HEAD').decode().strip()==C)
changed=g('diff','--name-only',B,C).decode().splitlines();yes('record-only-design-docs',bool(changed) and all(p.startswith('design_docs/') for p in changed),len(changed))
yes('source-reviewed-and-squash-tree-identical',g('rev-parse',S+'^{tree}')==g('rev-parse',B+'^{tree}'))
# Compare EVERY file outside design_docs against source tree, not just declared product files.
source={p for p in g('ls-tree','-r','--name-only',B).decode().splitlines() if not p.startswith('design_docs/')}
current={p for p in g('ls-files').decode().splitlines() if not p.startswith('design_docs/')}
yes('entire-product-and-frozen-surface-identical',source==current and all(b(B,p)==n(p) for p in source),{'files':len(source)})
a=b(B,'design_docs/world-mission.md');z=n('design_docs/world-mission.md')
def between(v,start,end):return v[v.index(start):v.index(end)]
led=between(z,b'## Human Decision Ledger',b'## STATUS (rotation rule)')
yes('entire-decision-ledger-exact',led==between(a,b'## Human Decision Ledger',b'## STATUS (rotation rule)'),{'rows':len(re.findall(rb'^\| D-',led,re.M))})
yes('only-D75-open',re.findall(rb'^\| (D-[^|]+) \| OPEN \|',led,re.M)==[b'D-WORLD-75'])
st=lambda x:re.findall(rb'^## STATUS 20[^\n]*\n\n',x,re.M)
x,y=st(a),st(z);yes('three-status-stamps268267266',len(x)==len(y)==3 and [int(re.search(rb'iteration (\d+)',s)[1]) for s in y]==[268,267,266]);yes('267266-status-stamps-exact',x[:2]==y[1:])
oldsa=b(B,'design_docs/world-mission-status-archive.md');newsa=n('design_docs/world-mission-status-archive.md');yes('265-stamp-move-and-archive-prefix-exact',newsa==oldsa+x[2])
rotated=a.replace(x[2],b'',1);mark=b'## STATUS (rotation rule)';j=rotated.index(b'\n\n',rotated.index(mark))+2;rotated=rotated[:j]+y[0]+rotated[j:];yes('rotation-prequeue-arithmetic',len(rotated.splitlines())==len(a.splitlines())+2-2 and st(rotated)==y,{'base_lines':len(a.splitlines()),'rotated_lines':len(rotated.splitlines())})
def q(v):return {int(m[1]):m[0] for m in re.finditer(rb'^(\d+)\. [^\n]*\n',v[v.index(b'## Queue '):],re.M)}
qa,qz=q(a),q(z);yes('queue-members177to179',len(qa)==177 and len(qz)==179 and set(qz)-set(qa)=={178,179});yes('176-other-queue-rows-exact',all(qa[k]==qz[k] for k in qa if k!=150))
marker=b'[v1.0.0 release row';yes('row150-only-current-tag',qa[150][qa[150].index(marker):]==qz[150][qz[150].index(marker):])
rest=z.replace(y[0],b'',1).replace(qz[150],qa[150],1)
for k in (178,179):rest=rest.replace(qz[k],b'',1)
rest=re.sub(rb'^\*\*Row129 recurrence, iteration268:\*\*[^\n]*\n',b'',rest,flags=re.M)
yes('charter-remainder-exact',rest.rstrip(b'\n')==a.replace(x[2],b'',1).rstrip(b'\n'))
# Independent record splitter uses all h2 headings, preserving preamble separately.
def rec(v):
 h=list(re.finditer(rb'^## [^\n]*\n',v,re.M));out={}
 for i,m in enumerate(h):
  title=m[0].decode().rstrip();num=re.match(r'## (\d+) [—-]',title)
  if num:key=int(num[1])
  elif 'continuation' in title.lower():key='attended'
  else:continue
  assert key not in out,(key,title)
  out[key]={'title':title,'body':v[m.end():h[i+1].start() if i+1<len(h) else len(v)]}
 return out
la=rec(b(B,'design_docs/world-mission-log.md'));aa=rec(b(B,'design_docs/world-mission-log-archive.md'));lz=rec(n('design_docs/world-mission-log.md'));az=rec(n('design_docs/world-mission-log-archive.md'));allz=lz|az
changes=[]
for scope,original in [('live',la),('archive',aa)]:
 for k,r in original.items():
  s=allz[k];interior=r['body'].rstrip(b'\n')==s['body'].rstrip(b'\n');yes('interior-'+scope+'-'+str(k),interior)
  if r['body']!=s['body']:changes.append({'scope':scope,'id':k,'lf_delta':len(s['body'])-len(r['body']),'interior_exact':interior})
yes('exact-five-declared-LF-boundaries',[(c['id'],c['lf_delta']) for c in changes]==[(247,1),('attended',-1),(257,-1),(267,1),(225,1)],changes)
yes('36-prior-numbered-plus-attended-to40-live-numbered',len(la)==37 and len(lz)==40 and 'attended' not in lz and 'attended' in az)
yes('historic-heading-changes-bounded',[k for k in la if la[k]['title']!=allz[k]['title']]==['attended',265] and all(aa[k]['title']==az[k]['title'] for k in aa))
yes('attended-heading-normalization-exact',allz['attended']['title']==la['attended']['title'].replace('Attended continuation','ATTENDED CONTINUATION'))
yes('265-title-correction-supported',b'8511bad' in la[265]['body'] and b'LANDED' in la[265]['body'])
idpat=rb'^\| (\d+) \|';ia=[int(i) for i in re.findall(idpat,b(B,'design_docs/world-mission-index.md'),re.M)];iz=[int(i) for i in re.findall(idpat,n('design_docs/world-mission-index.md'),re.M)]
yes('index257to258-only268-added',len(ia)==257 and len(iz)==258 and len(set(iz))==258 and set(iz)-set(ia)=={268} and not set(ia)-set(iz))
yes('index-one-attended-record',sum(1 for s in n('design_docs/world-mission-index.md').splitlines() if s.startswith(b'|') and b'ATTENDED CONTINUATION' in s)==1)
recovered=json.loads((D/'orphan-marker-recovery.json').read_text());yes('three-orphan-markers-from-original-index',len(recovered)==3 and all(r['original_index_row'].encode() in b(B,'design_docs/world-mission-index.md') and r['body'].encode() in lz[r['iteration']]['body'] and 'no gate, judge or outcome is invented' in r['body'] for r in recovered))
last=sorted(k for k in allz if isinstance(k,int))[-20:];yes('last20-primary-harness-zero-with-control',len(last)==20 and not any('[HARNESS]' in allz[k]['title'] for k in last) and any('[HARNESS]' in r['title'] for r in allz.values()),last)
yes('single-archive-notice',n('design_docs/world-mission-log.md').count(b'Older entries are ARCHIVED')==1)
# R1 patch is narrow and old numbered bodies never changed during this repair.
prevlog=rec(b(P,'design_docs/world-mission-log.md'));yes('R1-repair-only-current268-body',all(prevlog[k]==lz[k] for k in prevlog if k!=268))
yes('N1-qualified-in-current-record',b'interior body bytes' in lz[268]['body'] and b'five boundaries' in lz[268]['body'] and b'STATUS 265 was archived byte for byte' in lz[268]['body'])
plan='design_docs/planned/sprint_w-marketing-capture.json';pa=json.loads(b(B,plan));pz=json.loads(n(plan));ma={m['id']:m for m in pa['milestones']};mz={m['id']:m for m in pz['milestones']}
yes('plan3complete14pending',len(mz)==17 and collections.Counter(v['status'] for v in mz.values())=={'completed':3,'pending':14})
yes('plan-other16-milestones-identical',all(ma[k]==mz[k] for k in ma if k!='M2a'))
trim={k:v for k,v in mz['M2a'].items() if k in ma['M2a']};trim['status']=ma['M2a']['status'];yes('M2a-only-completion-metadata',trim==ma['M2a'])
e=mz['M2a']['completion_evidence'];yes('completion-correct-head-and-merge',e['merge_commit']==B and e['reviewed_head']==S and e['evaluator_score']==90 and e['evaluator_round']==3)
pat='snapshot changed files into $HOME/.ailang/state/world-iter<N>/snap/';yes('17plus17-external-snapshot-directives',n(plan).decode().count(pat)==17 and n('design_docs/planned/w-marketing-capture-sprint-plan.md').decode().count(pat)==17)
for filename,label,count,size in [('record-staging-manifest.json','product-R3',130,True),('record-evaluator-r1-manifest.json','record-R1',44,False)]:
 m=json.loads((D/filename).read_text());bad=[v['path'] for v in m if digest((D/v['path']).read_bytes())!=v['sha256'] or (size and len((D/v['path']).read_bytes())!=v['bytes'])];yes(label+'-manifest-intact',len(m)==count and not bad,{'files':len(m),'bad':bad})
prior=[p for p in g('ls-tree','-r','--name-only',B,'--','design_docs/verification/world-iter268').decode().splitlines() if re.search(r'/evaluator-r[12](?:[./])',p)];yes('prior-product-FAIL-reports-allbytes-identical',len(prior)==174 and all(n(p)==b(B,p) for p in prior),len(prior))
rr=json.loads((D/'record-evaluator-r1.json').read_text());pr=json.loads((D/'evaluator-r3.json').read_text());yes('prior-verdicts-retained',rr['result']=='FAIL' and rr['score']==76 and 'PASS' in str(pr).upper() and '90' in str(pr))
route=json.loads((D/'routing.json').read_text());yes('routing-actual-R1-and-productR3',route['record_evaluator']['rounds'][0]['verdict']=='FAIL' and route['evaluator']['rounds'][-1]['score']==90 and 'same-model-fresh-context' in route['judge_independence'])
yes('gate4-base-full-SHA-and-ISO',('base='+B+'@2026-10-15T06:18:07Z').encode() in lz[268]['body'] and ('Base='+B+'@2026-10-15T06:18:07Z').encode() in y[0])
dash=n('design_docs/world-mission-dashboard.md');yes('dashboard14lines-and-goal-unmoved',len(dash.splitlines())==14 and b'Goal unmoved' in dash and b'3/17' in dash and b'M2b' in dash)
kw=re.compile(r'(clos(e|es|ed)|fix(e|es|ed)?|resolv(e|es|ed))\s*:?\s*#[0-9]+',re.I);yes('commit-autoclose-control',bool(kw.search('this fixes #1')) and not kw.search(g('show','-s','--format=%B',C).decode()))
# Execute the actual consumer function definitions; no import-time shared side effects.
f=Path('/Users/voightkampff/dev/sunholo-data/ailang/tools/mission-weekly-report.py');tree=ast.parse(f.read_text());nodes=[t for t in tree.body if isinstance(t,ast.FunctionDef) and t.name in ('clean','goal_distance') or isinstance(t,ast.Assign) and any(isinstance(v,ast.Name) and v.id=='PROG_RE' for v in t.targets)];env={'re':re,'os':os};exec(compile(ast.Module(body=nodes,type_ignores=[]),str(f),'exec'),env)
log=R/'design_docs/world-mission-log.md';text=log.read_text();current_line=next(l for l in lz[268]['body'].decode().splitlines() if l.startswith('**Progress'));old_line=current_line.replace('**Progress**:','**Progress:**',1)
control=O/'regressed-progress-log.txt';control.write_text(text.replace(current_line,old_line,1));actual=env['goal_distance'](log);regressed=env['goal_distance'](control);expected=current_line.split(' ',1)[1]
pp={'consumer':str(f),'consumer_sha256':digest(f.read_bytes()),'current_line':current_line,'current_value':actual,'regressed_value':regressed,'expected':expected,'pattern':env['PROG_RE'].pattern,'candidate_matches':bool(env['PROG_RE'].match(current_line)),'regressed_matches':bool(env['PROG_RE'].match(old_line)),'current_heading':'268','token_usage':'unknown'}
(O/'consumer-progress.json').write_text(json.dumps(pp,indent=2)+'\n');yes('B1-fixed-actual-consumer-with-regression-control',actual==expected and regressed!=actual and pp['candidate_matches'] and not pp['regressed_matches'],pp)
# Independent live CI identity and terminal-success checks.
ci=[]
for label,rid,event,head in [('source-pr',38687777190,'pull_request',S),('source-merge',38689782525,'push',B)]:
 run=json.loads((O/(label+'-run.json')).read_text());jobs=json.loads((O/(label+'-jobs.json')).read_text());names={'ailang-code verify gate','go host build + test gate'};ok=run['id']==rid and run['event']==event and run['head_sha']==head and run['status']=='completed' and run['conclusion']=='success' and jobs['total_count']==2 and {j['name'] for j in jobs['jobs']}==names and all(j['status']=='completed' and j['conclusion']=='success' and j['head_sha']==head for j in jobs['jobs']);yes(label+'-exact-real-CI2of2',ok);ci.append({'id':rid,'event':run['event'],'head':run['head_sha'],'status':run['status'],'conclusion':run['conclusion'],'jobs':[{'name':j['name'],'head':j['head_sha'],'status':j['status'],'conclusion':j['conclusion']} for j in jobs['jobs']]})
api=json.loads((O/'source-pr.json').read_text());yes('source-PR267-merged-correct-head-base-squash',api['merged'] and api['merge_commit_sha']==B and api['head']['sha']==S and api['base']['sha']=='39066a46c1ac13c994bacb162b8d4b754d5ec083')
yes('latest-release-remains-none',json.loads((O/'releases.json').read_text())==[])
(O/'source-ci-validation.json').write_text(json.dumps(ci,indent=2)+'\n')
new=g('diff','--name-only','--diff-filter=A',B,C).decode().splitlines();tracked=set(g('ls-files').decode().splitlines());yes('all-new-record-artifacts-tracked',bool(new) and all(p in tracked for p in new),len(new))
result={'record_head':C,'source_merge':B,'checks':checks,'count':len(checks),'failures':[x for x in checks if not x['pass']],'boundary_normalizations':changes,'utc':datetime.datetime.now(datetime.timezone.utc).isoformat()};(O/'independent-review.json').write_text(json.dumps(result,indent=2)+'\n');print(json.dumps({'count':len(checks),'failures':result['failures'],'boundaries':changes,'consumer':pp,'CI':ci},indent=2))
