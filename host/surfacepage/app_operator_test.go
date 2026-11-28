package surfacepage

import "testing"

func runOperatorBar(t *testing.T, script string) {
	t.Helper()
	runOperatorScript(t, AppScript, `
const assert = require('node:assert/strict'), app = require(process.argv[1]);
const argv = ['ailang-worldd','episode','start','s8','--preset','p','--launch','--headless'];
const b = {schema:'world/surface-verb/v1',v:1,outcome:'review',argv,catalogueGeneration:7,text:'Headless codex gpt-6.2-sol'};
const review = (change={}) => app.verbOutcome(200,JSON.stringify({...b,...change}));
const id = 'oa-'+'a'.repeat(32), href = '/surface/operator-actions/'+id;
const result = {v:1, id, href, statusHref:href+'/status.json', kind:'start', episode:'s8', state:'requested', now:1, createdAt:1, expiresAt:601};
const reply = body => ({status:200,json:async()=>body});
let slot = null, posts = [], randoms = 0, current = true, response = reply(result);
const storage = { getItem:()=>slot, setItem:(key,value)=>{slot=value;}, removeItem:()=>{slot=null;} };
const io = {storage, random: x=>(randoms++, x.fill(randoms)), current:()=>current, timeout:ms=>(assert.equal(ms,15000), 'timeout'),
  fetch: async (path,opts)=>{ assert.equal(path,'/surface/operator-actions'); assert.equal(opts.signal,'timeout');
    assert.equal(opts.method,'POST'); assert.equal(opts.headers['Content-Type'],'application/json');
    assert.equal(opts.redirect,'error'); assert.equal(opts.credentials,'omit'); posts.push(JSON.parse(opts.body));
    if(!response) throw Error('offline'); return response; }};
const control = () => app.operatorRequestControl(review().review,io);
const deferred = () => { let resolve; const promise = new Promise(r=>resolve=r); return {promise,resolve}; };
let completed = false;
process.once('beforeExit', () => {
  if (!completed) { console.error('operator bar assertion did not settle'); process.exitCode=1; }
});
(async()=>{`+script+`})().then(() => { completed=true; }, e => { completed=true; console.error(e); process.exitCode=1; });
`)
}

func TestOperatorBarStartReview(t *testing.T) {
	runOperatorBar(t, `
const actions = [];
const fx = {ask:async()=>review(),act:{review:(p,live)=>actions.push([p,live()]),
  panel:()=>assert.fail('ordinary panel'),assign:()=>assert.fail('navigated'),clear:()=>assert.fail('cleared input'),say:assert.fail}};
await app.goWith(argv.join(' '),fx);
assert.equal(posts.length,0);
assert.equal(randoms,0);
assert.deepEqual(actions,[[{argv,catalogueGeneration:7,body:b.text},true]]);
assert.equal((await control().submit('')).error.includes('1–8192'),true);
assert.equal((await control().submit('é'.repeat(4097))).error.includes('1–8192'),true);
assert.equal((await control().submit('a\0b')).error.includes('without NUL'),true);
assert.equal(posts.length,0);
assert.equal((await control().submit('é'.repeat(4096))).href,href);
`)
}

func TestOperatorOtherVerbParity(t *testing.T) {
	runOperatorBar(t, `
const v = x=>app.verbOutcome(200,JSON.stringify({...b,...x}));
assert.deepEqual(v({outcome:'address'}),{address:true});
for (const cli of ['episode report e','why sha256:x']) {
  assert.deepEqual(v({outcome:'result',cli,text:'body'}),{panel:{head:'$ '+cli,body:'body',copy:''}});
}
const approve = '/surface/approve/mr-0123456789abcdef';
assert.deepEqual(v({outcome:'navigate',href:approve}),{assign:approve});
for(const cli of ['world-publish','approver enroll','approver revoke','session mint','session revoke','service']) {
  assert.deepEqual(v({outcome:'refused',cli,reason:'blocked'}),{panel:{head:'$ '+cli,body:'blocked',copy:cli}});
}
`)
}

func TestOperatorLiteralTask(t *testing.T) {
	runOperatorBar(t, `
const task = '# Task\n<script>alert(1)</script> ; $(literal) &';
assert.equal((await control().submit(task)).href,href);
assert.deepEqual(posts,[{argv,operationId:'01'.repeat(16),catalogueGeneration:7,task}]);
`)
}

func TestOperatorDuplicateTap(t *testing.T) {
	runOperatorBar(t, `
const pending = deferred(), r = control();
response = pending.promise;
const first = r.submit('literal');
const second = r.submit('different');
assert.equal(posts.length,1);
assert.deepEqual(await second,{});
pending.resolve(reply(result));
assert.deepEqual(await first,{href});
assert.equal(randoms,1);
assert.equal(slot,null);
`)
}

func TestOperatorReconnect(t *testing.T) {
	runOperatorBar(t, `
response = null;
assert.match((await control().submit('frozen')).error,/Reconnect/);
const draft = app.operatorDraft(storage);
assert.equal(draft.request.task,'frozen');
assert.equal(draft.body,b.text);
response = reply(result);
assert.deepEqual(await control().submit('edited'),{href});
assert.deepEqual(posts,[draft.request,draft.request]);
assert.equal(randoms,1);
assert.equal(slot,null);
`)
}

func TestOperatorPanelLatest(t *testing.T) {
	runOperatorBar(t, `
const old = deferred(), now = deferred(), acted = [];
const fx = {ask:q=>q==='older'?old.promise:now.promise,act:{review:p=>acted.push(p.body),
  assign:h=>acted.push(h),clear:()=>{},say:assert.fail}};
const a = app.goWith('older',fx), n = app.goWith('newer',fx);
now.resolve(review()); await n;
old.resolve({assign:href}); await a;
assert.deepEqual(acted,[b.text]);
current = false;
assert.deepEqual(await control().submit('stale'),{});
assert.equal(posts.length,0);
current = true;
const pending = deferred(); response = pending.promise;
const post = control().submit('frozen'); current = false;
pending.resolve(reply(result));
assert.deepEqual(await post,{});
assert.notEqual(slot,null);
`)
}

func TestOperatorPanelNetwork(t *testing.T) {
	runOperatorBar(t, `
const calls = [];
let trigger;
const fx = {ask: raw=>app.askVerb(raw, {AbortController,setTimeout:f=>(trigger=f,1),clearTimeout:()=>{},
  fetch:()=>new Promise(()=>{})}), act:{retry:x=>calls.push(x),say:()=>assert.fail('no reconnect offered'),clear:()=>assert.fail('cleared'),assign:()=>assert.fail('navigated')}};
const ask = app.goWith('start',fx); trigger(); await ask;
assert.match(calls[0],/in time/);
await app.goWith('start', {...fx,ask:async()=>review({argv:[]})});
assert.match(calls[1],/cannot read/);
await app.goWith('start', {...fx,ask:async()=>app.verbOutcome(503,'')});
assert.match(calls[2],/HTTP 503/);
await app.goWith('start', {...fx,ask:()=>app.askVerb('start',{AbortController,
 setTimeout:()=>1,clearTimeout:()=>{},fetch:async()=>{throw Error('offline')}})});
assert.match(calls[3],/unreachable/);
for (const change of [{argv:{}},{argv:{...argv,length:8}},{argv:[...argv,'tail']},{argv:['ailang-worldd', 'episode','start',{},'--preset','p','--launch','--headless']},
  {argv:argv.slice(1)},{argv:argv.with(7,'--launch')},{argv:argv.with(0,'sh')},{argv:argv.with(1,'service')},
  {argv:argv.with(2,'finish')},{argv:argv.with(3,'bad\n')},{argv:argv.with(4,'--other')},{argv:argv.with(5,3)},
  {argv:argv.with(5,'se-tools')},{argv:argv.with(5,'A')},{argv:argv.with(5,'a'.repeat(52))},{argv:argv.with(6,'--other')},
  {catalogueGeneration:0},{catalogueGeneration:-1},{catalogueGeneration:0.5},{catalogueGeneration:'1'},{text:1},{text:''}])
  assert.match(review(change).say,/cannot read/);
for (const d of [null, {}, {body:1,request:{argv,catalogueGeneration:1,operationId:'01'.repeat(16),task:'x'}},
  {body:'',request:{argv,catalogueGeneration:1,operationId:'01'.repeat(16),task:'x'}}, {body:'review',request:null},
  {body:'review',request:{argv,catalogueGeneration:1,operationId:1,task:'x'}},
  {body:'review',request:{argv,catalogueGeneration:1,operationId:['01'.repeat(16)],task:'x'}},
  {body:'review',request:{argv,catalogueGeneration:1,operationId:'x'.repeat(32),task:'x'}},
  {body:'review',request:{argv,catalogueGeneration:1,operationId:'01'.repeat(16),task:''}}]) {
  slot=JSON.stringify(d); assert.equal(app.operatorDraft(storage),null);
}
slot=null;
const save=storage.setItem;
storage.setItem=()=>{ throw Error('disabled storage'); };
const r = control();
assert.match((await r.submit('original')).error,/could not be saved/);
assert.equal(posts.length,0);
storage.setItem=save;
response = null;
assert.match((await r.submit('original')).error,/Reconnect/);
response = reply({...result,state:['requested']});
assert.match((await r.submit('will not replace')).error,/Reconnect/);
response = {status:409,json:async()=>result};
assert.match((await r.submit('again')).error,/Reconnect/);
response = {status:409,json:async()=>({error:'Conflict',reason:'changed'})};
assert.match((await r.submit('again')).error,/Conflict: changed.*Reconnect/);
assert.equal(randoms,1);
assert(posts.every(x=>JSON.stringify(x)===JSON.stringify(posts[0])));
`)
}

func TestOperatorLinkAllowlist(t *testing.T) {
	runOperatorBar(t, `
const pick = b => app.operatorCardHref(b,'s8');
assert.equal(pick(result),href);
for (const bad of [href+'x',href+'a',href.slice(0,-1),'x'+href,href+'\n',href.toUpperCase(),href+'#f',href+'?q',href+'/status.json',
  'https://localhost'+href,'http://localhost:8'+href,'http://user@localhost'+href,'//localhost'+href,{toString:()=>href}]) {
  assert.equal(pick({...result,href:bad}),'',JSON.stringify(bad));
  assert.equal(pick({...result,statusHref:bad===result.statusHref ? result.href : bad}),'',JSON.stringify(bad));
}
for (const bad of [id+'x','x'+id,id+'\n','x\n'+id,id+'\nx',id.toUpperCase(),id+'a',id.slice(0,-1),{toString:()=>id}]) {
  const h = '/surface/operator-actions/'+bad;
  assert.equal(pick({...result,id:bad,href:h,statusHref:h+'/status.json'}),'');
}
for (const [k,v] of Object.entries({v:'1',episode:'other',kind:'finish',state:'invented',now:'1',createdAt:[],expiresAt:null}))
  assert.equal(pick({...result,[k]:v}),'',k);
`)
}

func TestOperatorFinishBarReview(t *testing.T) {
	runOperatorBar(t, `
const finishArgv = ['ailang-worldd','episode','finish','s8','--preset','p','--confirm-on-surface'];
for(const bad of [finishArgv.slice(0,6), [...finishArgv.slice(0,6), "--dry-run"], finishArgv.concat('--title'), finishArgv.concat('--ready','--ready')])
  assert.equal(review({argv:bad}).review, undefined);
const admitted = review({argv:finishArgv.concat('--ready')}).review;
assert.deepEqual(admitted.argv, finishArgv.concat('--ready'));
assert.equal(posts.length,0);
const startCase = app.operatorRequestControl(admitted,io);
assert.ok((await startCase.submit('not a task')).error);
assert.equal(posts.length,0);
assert.ok((await startCase.submit('')).error); // A start card answers a finish request.
assert.deepEqual(posts[0].argv, finishArgv.concat('--ready'));
assert.equal(posts[0].task,'');
response = reply({...result, kind:'finish'});
const reconnected = app.operatorRequestControl(admitted,io);
assert.equal((await reconnected.submit('')).href,href);
assert.deepEqual(posts[0],posts[1]);
assert.equal(slot,null);
const mounted = [], clicked = [];
let focus, link;
const el = { isConnected:true, querySelector:()=>({replaceWith(){}}) };
const fakeDoc = {createElement:tag=>({tag, append(...children){ mounted.push(...children); },
  setAttribute(){}, focus(){ focus=this; }, addEventListener(name,fn){ if(name==='click') clicked.push(fn); }})};
app.showOperatorReview(admitted, () => true, {...io, doc:fakeDoc, show:()=>el, assign:x=>{link=x;}});
assert.equal(mounted.some(n=>n.tag==='textarea'||n.tag==='label'),false);
assert.equal(focus.textContent,'Request approval');
assert.equal(posts.length,2);
await clicked[0]();
assert.equal(posts.length,3);
assert.equal(posts[2].task,'');
assert.equal(focus.id,'operator-request');
assert.equal(link,href);
`)
}
