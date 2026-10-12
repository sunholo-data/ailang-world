import json,collections,statistics as st,sys
T='/Users/voightkampff/dev/sunholo-data/ailang-world/.ailang/state/row174-attended-2026-10-10/rig/m4/transcripts/iter-11/'
E='../../../design_docs/verification/world-floor-m4-2026-10-06/iter-11/runs/codex/'
rows={arm:{json.loads(l)['task']:json.loads(l) for l in open(E+arm+'/r1.jsonl')} for arm in ('shell','world')}
def tr(f):
    calls=[];msgs=0;errs=0;usage=None;turns=0;outbytes=0
    for l in open(T+f):
        try:d=json.loads(l)
        except:continue
        if d.get('type')=='turn.completed':usage=d.get('usage');turns+=1
        it=d.get('item') or {}
        if d.get('type')=='item.completed':
            t=it.get('type')
            if t=='agent_message':msgs+=1
            elif t in('mcp_tool_call','command_execution','file_change'):
                name=it.get('tool') or t
                if t=='command_execution':name='sh:'+(it.get('command','').split()[2] if len(it.get('command','').split())>2 else it.get('command',''))[:20]
                bad=bool(it.get('error')) or it.get('status')=='failed' or (t=='command_execution' and it.get('exit_code') not in (0,None))
                r=it.get('result') or it.get('aggregated_output') or ''
                outbytes+=len(json.dumps(r)) if not isinstance(r,str) else len(r)
                calls.append((name,bad));errs+=bad
    return dict(calls=calls,msgs=msgs,errs=errs,usage=usage,outbytes=outbytes)
out=[]
for t in sorted(rows['shell']):
    s,w=rows['shell'][t],rows['world'][t]
    ts,tw=tr(s['transcript_file']),tr(w['transcript_file'])
    o=w['wall_ms']/s['wall_ms']-1
    out.append(dict(task=t,o=o,ws=s['wall_ms']/1000,ww=w['wall_ms']/1000,ins=s['input_tokens'],inw=w['input_tokens'],outs=s['output_tokens'],outw=w['output_tokens'],
      ncs=len(ts['calls']),ncw=len(tw['calls']),errw=tw['errs'],errs=ts['errs'],bys=ts['outbytes'],byw=tw['outbytes'],
      wtools=collections.Counter(n for n,_ in tw['calls']),wbad=collections.Counter(n for n,b in tw['calls'] if b)))
out.sort(key=lambda r:-r['o'])
print(f"{'task':28}{'o':>7}{'wall S/W s':>13}{'in S/W k':>14}{'outtok S/W':>12}{'calls S/W':>10}{'Werr':>5}  world tools")
for r in out:
    print(f"{r['task'][:27]:28}{r['o']*100:6.0f}%{r['ws']:6.0f}/{r['ww']:<6.0f}{r['ins']/1e3:6.0f}/{r['inw']/1e3:<7.0f}{r['outs']:5}/{r['outw']:<6}{r['ncs']:4}/{r['ncw']:<5}{r['errw']:4}  {dict(r['wtools'])} {('BAD:'+str(dict(r['wbad']))) if r['wbad'] else ''}")
hi=[r for r in out if r['o']>0.25]; lo=[r for r in out if r['o']<=0.25]
for name,g in (('>25%',hi),('<=25%',lo)):
    print(name,len(g),'median in-token ratio',round(st.median(r['inw']/r['ins'] for r in g),2),'median call ratio',round(st.median(r['ncw']/max(r['ncs'],1) for r in g),2),
      'median out-token ratio',round(st.median(r['outw']/r['outs'] for r in g),2),'tool errs',sum(r['errw'] for r in g))
json.dump([{k:(dict(v) if isinstance(v,collections.Counter) else v) for k,v in r.items()} for r in out],open('analysis.json','w'),indent=1)
