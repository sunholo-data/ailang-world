package workbench

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func liveScriptTestAsset(t *testing.T) []byte {
	t.Helper()
	b := LiveScript
	if len(b) == 0 {
		t.Fatal("live asset empty")
	}
	return b
}
func TestLiveScriptAssetNonEmpty(t *testing.T) {
	b := liveScriptTestAsset(t)
	for _, name := range []string{"splitFrames", "applyDelta", "nextBackoff", "applySwap", "renderStatus"} {
		if !strings.Contains(string(b), name) {
			t.Fatalf("missing seam %s", name)
		}
	}
}
func TestLiveScriptPure(t *testing.T) {
	b := liveScriptTestAsset(t)
	node := os.Getenv("WORLD_EXEC_NODE")
	if node == "" {
		var err error
		node, err = exec.LookPath("node")
		if err != nil {
			t.Skip("node unavailable locally")
		}
	}
	golden, err := os.ReadFile("../agui/testdata/stream_fixture.golden")
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	for name, data := range map[string][]byte{"live.js": b, "golden": golden, "harness.js": []byte(liveHarness)} {
		if err := os.WriteFile(filepath.Join(dir, name), data, 0600); err != nil {
			t.Fatal(err)
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, node, filepath.Join(dir, "harness.js"))
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("node harness: %v\n%s", err, out)
	}
	if !strings.Contains(string(out), "LIVE HARNESS PASS") {
		t.Fatalf("harness manifest missing: %s", out)
	}
}

const liveHarness = `
const assert = require('node:assert/strict');
const fs = require('node:fs');
const vm = require('node:vm');
const api = require('./live.js');
assert.deepEqual(Object.keys(api).sort(), ['applyDelta','applySwap','nextBackoff','renderStatus','splitFrames']);
const bytes = fs.readFileSync('golden');
// Independent stock-client splitter; never include the unterminated tail.
function stock(text) {return text.split('\n\n').slice(0,-1).flatMap(f=>f.split('\n').filter(l=>l.startsWith('data: ')).map(l=>JSON.parse(l.slice(6))));}
const expected=stock(bytes.toString());assert(expected.length>0);
function stateFor(events) {let s={lastIndex:-1};for(const e of events) {if(e.type==='STATE_SNAPSHOT')s={...e.snapshot};if(e.type==='STATE_DELTA')api.applyDelta(s,e.delta);}return s;}
for(let cut=0;cut<=bytes.length;cut++) {
 const decoder=new TextDecoder();
 const first=decoder.decode(bytes.subarray(0,cut),{stream:true});
 const a=api.splitFrames('',first);
 assert.deepEqual(a.events,stock(first),'complete frames at byte '+cut);
 assert.equal(stateFor(a.events).lastIndex,stateFor(stock(first)).lastIndex,'severed cursor '+cut);
 const second=decoder.decode(bytes.subarray(cut),{stream:true})+decoder.decode();
 const b=api.splitFrames(a.tail,second);
 assert.deepEqual([...a.events,...b.events],expected,'loss/duplication at byte '+cut);
 assert.equal(b.tail,'');
}
assert.throws(()=>api.applyDelta({lastIndex:0},[{op:'add',path:'/lastIndex',value:1}]));
let back=0;for(const want of [1,2,4,8,16,30,30]) {back=api.nextBackoff('1',back);assert.equal(back,want);}
assert.equal(api.nextBackoff(undefined),1);assert.equal(api.nextBackoff('x'),1);
const regions=['nav[aria-label="world browser"]','section[aria-label="timeline"]','section[aria-label="live"]','section[aria-label="world graph"]','section[aria-label="decisions"]'];
const footer='footer[aria-label="live status"]';
function doc(missing) {
 const replacements=[];const span={textContent:'Live updates: off. Reload to refresh.'};
 const nodes=new Map([...regions,footer].filter(s=>s!==missing).map(s=>[s,{replaceWith(n){replacements.push(s);},querySelectorAll(){return [];}}]));
 return {replacements,span,visibilityState:'visible',querySelector(s){if(s==='[data-live-status]')return span;return nodes.get(s)||null;},importNode(n){return n;},addEventListener(n,f){this.listener=f;}};
}
for(const state of [{mode:'live',cursor:6,checked:'12:00:00',count:0},{mode:'paused',reason:'network',wait:2,cursor:6}]) {
 const old=doc(),fresh=doc();old.span.textContent=api.renderStatus(state);
 assert.equal(api.applySwap(old,fresh,state),true);
 assert.equal(old.span.textContent,api.renderStatus(state));assert.equal(old.replacements.length,5);assert(!old.replacements.includes(footer));
 old.span.textContent='Live updates: off. Reload to refresh.';
 api.applySwap(old,doc(),state);assert.equal(old.span.textContent,api.renderStatus(state),'status reassert');
}
for(const side of ['old','fresh']) {
 const old=doc(side==='old'?regions[4]:undefined),fresh=doc(side==='fresh'?regions[4]:undefined),s={mode:'live',cursor:6};
 assert.doesNotThrow(()=>assert.equal(api.applySwap(old,fresh,s),false));
 assert.equal(old.replacements.length,4);assert.equal(old.span.textContent,'paused: page layout changed, reload required');
}
const terminal='data: {"type":"RUN_FINISHED","result":{"lastIndex":42}}\n\n';
// Fresh VM per lifecycle stimulus: fake time, fetch, readers and aborts.
function rig() {
 const d=doc();d.querySelector(regions[2]).dataset={liveCursor:'6'};
 let next=0;const timers=new Map(),calls=[],pending=[];
 class AC {constructor(){this.signal={aborted:false};}abort(){this.signal.aborted=true;if(this.signal.cancel)this.signal.cancel();}}
 const context={document:d,location:{href:'http://local/workbench?from=0&entry=0'},DOMParser:class{parseFromString(){return doc();}},TextDecoder,Date,AbortController:AC,
 setTimeout(fn,ms){const id=++next;timers.set(id,{fn,ms});return id;},clearTimeout(id){timers.delete(id);},
 fetch(url,options){calls.push({url,options});return new Promise((resolve,reject)=>pending.push({resolve,reject}));}};
 vm.runInNewContext(fs.readFileSync('live.js','utf8'),context);
 return {d,timers,calls,pending,fire(ms){const item=[...timers].find(([id,t])=>t.ms===ms);assert(item,'missing timer '+ms);timers.delete(item[0]);item[1].fn();}};
}
async function flush(){for(let i=0;i<12;i++)await Promise.resolve();}
function stream() {const queue=[],wait=[];return {read(){return new Promise(resolve=>queue.length?resolve(queue.shift()):wait.push(resolve));},push(text,done=false){const v={value:new TextEncoder().encode(text),done};if(wait.length)wait.shift()(v);else queue.push(v);}};}
function response(reader){return {ok:true,status:200,body:{getReader(){return reader;}},headers:{get(){return null;}}};}
(async()=>{
 let r=rig();assert.equal(r.calls.length,1);assert.equal(r.calls[0].url,'/agui/');assert.equal(JSON.parse(r.calls[0].options.body).state.lastIndex,6);
 r.d.visibilityState='hidden';r.d.listener();assert(r.calls[0].options.signal.aborted);r.pending.shift().reject(Error('abort'));await flush();assert.equal(r.calls.length,1);
 r.d.visibilityState='visible';r.d.listener();await flush();assert.equal(r.calls.length,2,'visible resume');
 let reader=stream();r.pending.shift().resolve(response(reader));await flush();reader.push(terminal);await flush();assert.equal(r.calls.length,3,'terminal immediate repost');assert.equal(JSON.parse(r.calls[2].options.body).state.lastIndex,42);
 r=rig();reader=stream();r.pending.shift().resolve(response(reader));await flush();
 const custom='data: {"type":"CUSTOM","name":"world.entry.committed","value":{}}\n\n';
 reader.push(custom+custom);await flush();assert.equal([...r.timers.values()].filter(t=>t.ms===300).length,1,'burst debounce');
 r.fire(300);await flush();assert.equal(r.calls.filter(c=>c.url.startsWith('http')).length,1);
 reader.push(custom);await flush();assert.equal([...r.timers.values()].filter(t=>t.ms===300).length,0,'one refresh in flight');
 const refresh=r.pending.shift();refresh.resolve({ok:false,status:503,statusText:'Timeout'});await flush();assert.equal(r.d.replacements.length,0);assert.match(r.d.span.textContent,/refresh failed.*503/);
 reader.push(custom);await flush();r.fire(300);await flush();r.pending.shift().resolve({ok:true,text:async()=>''});await flush();assert.equal(r.d.replacements.length,5);
 for(const failure of ['503','network','RUN_ERROR']) {
  r=rig();
  if(failure==='503')r.pending.shift().resolve({ok:false,status:503,headers:{get(){return '1';}}});
  if(failure==='network')r.pending.shift().reject(Error('network'));
  if(failure==='RUN_ERROR'){reader=stream();r.pending.shift().resolve(response(reader));await flush();reader.push('data: {"type":"RUN_ERROR","code":"Internal"}\n\n');}
  await flush();assert.match(r.d.span.textContent,/paused:/);assert.equal(r.calls.length,1);r.fire(1000);await flush();assert.equal(r.calls.length,2);
 }
 r=rig();reader=stream();r.pending.shift().resolve(response(reader));await flush();reader.push(custom);await flush();r.fire(300);await flush();
 // Old-side missing region is discovered at the real lifecycle boundary.
 const original=r.d.querySelector.bind(r.d);r.d.querySelector=s=>s===regions[4]?null:original(s);
 r.pending.shift().resolve({ok:true,text:async()=>''});await flush();assert.equal(r.d.replacements.length,4);assert.equal(r.d.span.textContent,'paused: page layout changed, reload required');
 reader.push(terminal);await flush();assert.equal(r.calls.length,2,'layout pause must stop repost');
 console.log('LIVE HARNESS PASS');
})().catch(e=>{console.error(e);process.exitCode=1;});
`
