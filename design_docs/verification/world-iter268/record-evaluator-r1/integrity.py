from pathlib import Path
import subprocess,re,json,hashlib,collections
ROOT=Path.cwd(); RAW=ROOT/'design_docs/verification/world-iter268/record-evaluator-r1'; BASE='420c07e7c1c70eca6043a0c4e12a20f2adcaf972'; RECORD='d6e04632851bb3ebd31c2ef5c4d32a114bb6a1ca'; REVIEWED='763c60a20358f486df395043131557c1143735ba'
def git(*args):return subprocess.check_output(['git',*args],cwd=ROOT,timeout=60)
def old(p):return git('show',f'{BASE}:{p}')
def now(p):return (ROOT/p).read_bytes()
def sha(b):return hashlib.sha256(b).hexdigest()
checks=[]
def check(name,value,detail=None):
 checks.append({'name':name,'pass':bool(value),'detail':detail});print(name, 'PASS' if value else 'FAIL',str(detail)[:200] if detail else '')
check('exact-record-head',git('rev-parse','HEAD').decode().strip()==RECORD)
changed=git('diff','--name-only',BASE,RECORD).decode().splitlines(); check('ADMIN-only-scope',all(x.startswith('design_docs/') for x in changed),{'paths':len(changed)})
check('source-reviewed-merge-trees-equal',git('rev-parse',REVIEWED+'^{tree}')==git('rev-parse',BASE+'^{tree}'))
paths=['cmd/ailang-worldd/demo.go','cmd/ailang-worldd/demo_seed.go','cmd/ailang-worldd/demo_test.go','cmd/ailang-worldd/main.go','docs/QUICKSTART.md']
check('five-product-paths-byte-equal',all(old(p)==now(p) for p in paths),{p:sha(now(p)) for p in paths})
protected=[p for p in git('ls-files').decode().splitlines() if (p.startswith('world/') and p.endswith('.ail')) or p in ['go.mod','go.sum'] or p.startswith('tools/launchd/')]
check('frozen-paths-byte-equal',all(old(p)==now(p) for p in protected),{'paths':len(protected),'files':protected})
# Whole ledger and whole-charter preservation outside declared changes.
p='design_docs/world-mission.md'; before=old(p).decode(); after=now(p).decode()
def ledger(x):return x[x.index('## Human Decision Ledger'):x.index('## STATUS (rotation rule)')]
check('ledger-whole-bytes-preserved',ledger(before)==ledger(after),{'rows':len(re.findall(r'^\| D-',ledger(after),re.M))})
check('only-D75-open',re.findall(r'^\| (D-[^|]+?) \| OPEN \|',ledger(after),re.M)==['D-WORLD-75'])
status=lambda x:re.findall(r'^## STATUS 20[^\n]*\n\n',x,re.M)
bs,as_=status(before),status(after)
check('status-three-single-lines',len(bs)==len(as_)==3)
check('status-current-268-267-266',[int(re.search(r'iteration (\d+)',x).group(1)) for x in as_]==[268,267,266])
check('old-status267266-unchanged',as_[1:]==bs[:2])
archive='design_docs/world-mission-status-archive.md'; ob,nb=old(archive),now(archive)
check('status-archive-prefix-and-exact265',nb.startswith(ob) and nb[len(ob):].decode()==bs[2],{'old_lines':len(ob.splitlines()),'new_lines':len(nb.splitlines())})
rot=before.replace(bs[2],''); at=rot.index('## STATUS (rotation rule)'); pos=rot.index('\n\n',at)+2;rot=rot[:pos]+as_[0]+rot[pos:]
check('status-before-queue-edit-arithmetic',len(rot.splitlines())==len(before.splitlines())+2-2 and status(rot)==as_,{'before':len(before.splitlines()),'after_rotation_before_queue_edits':len(rot.splitlines()),'moved':1})
def rows(x):return {int(m.group(1)):m.group(0) for m in re.finditer(r'^(\d+)\. [^\n]*(?:\n|$)',x[x.index('## Queue '):],re.M)}
br,ar=rows(before),rows(after)
check('queue-177-to179',len(br)==177 and len(ar)==179 and set(ar)-set(br)=={178,179})
check('queue-other176-lines-byte-equal',all(br[k]==ar[k] for k in br if k!=150))
sep='[v1.0.0 release row';check('row150-only-current-tag-changed',br[150][br[150].index(sep):]==ar[150][ar[150].index(sep):])
# Strip precisely the documented additions and substitutions; remaining charter must equal base.
rest=after.replace(as_[0],'',1).replace(ar[150],br[150],1)
for k in [178,179]:rest=rest.replace(ar[k],'',1)
rest=re.sub(r'^\*\*Row129 recurrence, iteration268:\*\*[^\n]*\n','',rest,flags=re.M)
# Appended additions each carried one blank separator; compare stripping only terminal blank lines.
base_no265=before.replace(bs[2],'',1)
check('charter-entire-other-content-preserved',rest.rstrip('\n')==base_no265.rstrip('\n'))
# Log records preserve bodies byte-for-byte independent of index generation.
pat=re.compile(rb'^## (?:[0-9]+ [\xe2\x80\x94]|Attended continuation|ATTENDED CONTINUATION)[^\n]*\n',re.M)
def records(data):
 matches=list(pat.finditer(data));out=[]
 for i,m in enumerate(matches):
  end=matches[i+1].start() if i+1<len(matches) else len(data)
  heading=m.group().decode().rstrip('\n'); n=re.match(r'^## (\d+) ',heading); key=int(n.group(1)) if n else 'attended'
  out.append({'key':key,'heading':heading,'body':data[m.end():end]})
 return out
ol=records(old('design_docs/world-mission-log.md'));oa=records(old('design_docs/world-mission-log-archive.md'));nl=records(now('design_docs/world-mission-log.md'));na=records(now('design_docs/world-mission-log-archive.md'))
om={r['key']:r for r in ol}; am={r['key']:r for r in oa};nm={r['key']:r for r in nl+na}
check('old-live36-numbered-plus-note',len(ol)==37 and len([x for x in ol if isinstance(x['key'],int)])==36)
check('old-live-all37-body-interiors-exact',all(nm[k]['body'].rstrip(b'\n')==r['body'].rstrip(b'\n') for k,r in om.items()))
check('old-archive-all-body-interiors-and-headings-exact',all(nm[k]['body'].rstrip(b'\n')==r['body'].rstrip(b'\n') and nm[k]['heading']==r['heading'] for k,r in am.items()),{'records':len(oa)})
check('live40-numbered',len(nl)==40 and all(isinstance(r['key'],int) for r in nl))
check('attended-note-rotated-body-exact','attended' not in {r['key'] for r in nl} and nm['attended']['body'].rstrip(b'\n')==om['attended']['body'].rstrip(b'\n') and nm['attended']['heading']==om['attended']['heading'].replace('Attended continuation','ATTENDED CONTINUATION'))
check('only265-historical-numbered-title-corrected',[k for k,r in om.items() if isinstance(k,int) and nm[k]['heading']!=r['heading']]==[265])
check('265-correction-supported-in-old-body',b'8511bad' in om[265]['body'] and b'LANDED' in om[265]['body'])
idx='design_docs/world-mission-index.md'; ids=lambda b:[int(n) for n in re.findall(rb'^\| (\d+) \|',b,re.M)]
bi,ai=ids(old(idx)),ids(now(idx));check('index257-to258-no-lost-ids',len(bi)==257 and len(ai)==258 and set(ai)-set(bi)=={268} and not set(bi)-set(ai) and len(set(ai))==258)
check('index-one-attended-row',sum(1 for l in now(idx).splitlines() if l.startswith(b'|') and b'ATTENDED CONTINUATION' in l)==1)
recovery=json.loads(now('design_docs/verification/world-iter268/orphan-marker-recovery.json'))
check('three-markers-faithful-to-old-index',len(recovery)==3 and all(r['original_index_row'].encode() in old(idx) and r['body'].encode() in nm[r['iteration']]['body'] and 'no gate, judge or outcome is invented' in r['body'] for r in recovery))
check('single-live-archive-notice',now('design_docs/world-mission-log.md').count(b'Older entries are ARCHIVED')==1)
last20=sorted([r for r in nl+na if isinstance(r['key'],int)],key=lambda r:r['key'])[-20:]
harness=[r['key'] for r in last20 if '[HARNESS]' in r['heading']]
check('last20-full-headline-primary-HARNESS0',len(last20)==20 and not harness,{'ids':[r['key'] for r in last20],'harness':harness,'positive_control_any_old_harness':any('[HARNESS]' in r['heading'] for r in nl+na)})
plan='design_docs/planned/sprint_w-marketing-capture.json'; pb=json.loads(old(plan));pa=json.loads(now(plan));bm={m['id']:m for m in pb['milestones']};mm={m['id']:m for m in pa['milestones']}
check('17-milestones3-complete14-pending',len(mm)==17 and collections.Counter(m['status'] for m in mm.values())=={'completed':3,'pending':14})
check('other16-milestone-objects-deep-equal',all(bm[k]==mm[k] for k in bm if k!='M2a'))
a2=dict(mm['M2a']);b2=bm['M2a']; extras=set(a2)-set(b2)
for k in extras:a2.pop(k)
a2['status']=b2['status'];check('M2a-no-requirement-or-gate-change',a2==b2,{'new_fields':sorted(extras)})
check('M2a-completion-exact-source',mm['M2a']['completion_evidence']['merge_commit']==BASE and mm['M2a']['completion_evidence']['reviewed_head']==REVIEWED)
md=now('design_docs/planned/w-marketing-capture-sprint-plan.md').decode()
check('17external-snapshots-in-each-plan',len(re.findall(r'snapshot changed files into \$HOME/\.ailang/state/world-iter<N>/snap/',md))==17 and now(plan).decode().count('snapshot changed files into $HOME/.ailang/state/world-iter<N>/snap/')==17)
manifest=json.loads(now('design_docs/verification/world-iter268/record-staging-manifest.json'));proot='design_docs/verification/world-iter268/'
check('130-R3-manifest-byte-hash-and-size',len(manifest)==130 and all(sha(now(proot+r['path']))==r['sha256'] and len(now(proot+r['path']))==r['bytes'] for r in manifest))
prior=[p for p in git('ls-tree','-r','--name-only',BASE,'--','design_docs/verification/world-iter268/evaluator-r1','design_docs/verification/world-iter268/evaluator-r2').decode().splitlines()]
# Prefix pathspec does not include sibling .md/.json, so explicitly enumerate all historical report/evidence paths.
prior=[p for p in git('ls-tree','-r','--name-only',BASE,'--','design_docs/verification/world-iter268').decode().splitlines() if re.search(r'/evaluator-r[12](?:[./])',p)]
check('all-prior-r1r2-report-and-raw-bytes-preserved',len(prior)>100 and all(now(p)==old(p) for p in prior),{'files':len(prior)})
check('dashboard-under40lines',len(now('design_docs/world-mission-dashboard.md').splitlines())<=40)
check('gate4-full-base-ISO-in-entry-and-status',b'base=420c07e7c1c70eca6043a0c4e12a20f2adcaf972@2026-10-15T06:18:07Z' in nm[268]['body'] and b'Base=420c07e7c1c70eca6043a0c4e12a20f2adcaf972@2026-10-15T06:18:07Z' in as_[0].encode())
newfiles=git('diff','--name-only','--diff-filter=A',BASE,RECORD).decode().splitlines()
check('all-new-candidate-artifacts-tracked',all(p in set(git('ls-files').decode().splitlines()) for p in newfiles),{'files':len(newfiles)})
msg=git('log','-1','--format=%B',RECORD).decode();kw=re.compile(r'(?:clos(?:e|es|ed)|fix(?:e|es|ed)?|resolv(?:e|es|ed))\s*:?\s*#[0-9]+',re.I)
check('record-commit-no-auto-close-keyword-with-control',not kw.search(msg) and bool(kw.search('this fixes #1')))
(RAW/'boundary-normalization.json').write_text(json.dumps([{'scope':scope,'iteration':k,'old_body_bytes_including_separator':len(r['body']),'new_body_bytes_including_separator':len(nm[k]['body']),'trailing_LF_delta':len(nm[k]['body'])-len(r['body']),'interior_exact':r['body'].rstrip(b'\n')==nm[k]['body'].rstrip(b'\n')} for scope,mapping in [('old-live',om),('old-archive',am)] for k,r in mapping.items() if r['body']!=nm[k]['body']],indent=2)+'\n')
(RAW/'integrity.json').write_text(json.dumps({'base':BASE,'record':RECORD,'reviewed_source':REVIEWED,'checks':checks,'failures':[x['name'] for x in checks if not x['pass']]},indent=2)+'\n')
