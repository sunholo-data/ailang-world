package surfacepage

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
)

func runOperatorScript(t *testing.T, asset []byte, harness string) {
	t.Helper()
	node := os.Getenv("WORLD_EXEC_NODE")
	if node == "" {
		var err error
		node, err = exec.LookPath("node")
		if err != nil {
			t.Skip("node is required to execute the served script")
		}
	}
	path := filepath.Join(t.TempDir(), "asset.js")
	if err := os.WriteFile(path, asset, 0600); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, node, "-e", harness, path)
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("%v\n%s", err, out)
	}
}

func TestOperatorCardStatusDisconnected(t *testing.T) {
	const origin, rp = "http://localhost:7644", "localhost"
	digest := sha256.Sum256([]byte("operator-card fixture"))
	challenge := [32]byte{9}
	message, err := approval.OperatorActionMessage(digest, challenge, 1600)
	if err != nil {
		t.Fatal(err)
	}
	hash := sha256.Sum256(message)
	key := approvaltest.NewPasskey(approval.AlgES256, 93)
	fixture, err := json.Marshal(map[string]any{
		"digest": hex.EncodeToString(digest[:]), "message": hex.EncodeToString(message), "hash": hex.EncodeToString(hash[:]),
		"credential": base64.RawURLEncoding.EncodeToString(key.CredentialID()),
		"assertion":  json.RawMessage(key.Assert(origin, rp, hash, 2)),
	})
	if err != nil {
		t.Fatal(err)
	}
	runOperatorScript(t, OperatorScript, "const fixture = "+string(fixture)+";\n"+`
const fs = require('node:fs'), vm = require('node:vm'), assert = require('node:assert/strict');
const path = '/surface/operator-actions/oa-' + 'a'.repeat(32), origin = 'http://localhost:7644';
const nodes = {}, events = {}, sent = [];
const node = id => nodes[id] ??= { textContent: '', dataset: {state: 'requested'},
  focus() { events.focus = id; }, addEventListener(k, f) { events[id] = f; }, removeAttribute() { delete this.href; } };
const status = {v:1, id:path.split('/').pop(), episode:'e-1', kind:'start', href:path, statusHref:path+'/status.json',
  state:'requested', now:1001, createdAt:1000, expiresAt:1600};
const c = {post:path+'/confirm', digest:fixture.digest, messageHex:fixture.message, passkeys:[{origin:'http://localhost:7645',rp:'localhost',allow:['Ag']},
  {origin,rp:'not-localhost',allow:['Aw']},{origin,rp:'localhost',allow:[fixture.credential]}]};
node('operator-ceremony').textContent = JSON.stringify(c);
const get = { ...status, now:1002 };
let response = get, code = 200, poll, asked;
const context = {document:{getElementById:node}, location:{pathname:path,origin,hostname:'localhost'}, atob, AbortSignal,
  crypto:require('node:crypto').webcrypto, navigator:{credentials:{get: async options => {
    asked = options;
    return {toJSON: () => fixture.assertion};
  }}}, setTimeout: f => (poll = f, 1), clearTimeout: () => {},
  fetch: async (url, options) => { sent.push([url, options]); if (!response) throw Error('offline');
    return {status:code,json:async()=>response}; }};
vm.runInNewContext(fs.readFileSync(process.argv[1], 'utf8'), context);
(async () => {
  assert.equal(events.focus,'operator-title');
  assert.equal(sent.filter(([,o])=>o.method==='POST').length,0);
  await new Promise(setImmediate);
  assert.match(node('operator-progress').textContent, /Status at 1002/);
  assert.equal(node('operator-controls').hidden,false);
  assert.equal(node('operator-key').hidden,false);
  assert.equal(sent[0][0],path+'/status.json');
  assert.equal(sent[0][1].redirect,'error');
  assert.equal(sent[0][1].credentials,'omit');
  assert.equal(sent[0][1].cache,'no-store');
  assert(sent[0][1].signal instanceof AbortSignal);
  response = null; await poll();
  assert.match(node('operator-progress').textContent,/Disconnected.*1002/);
  assert.equal(node('operator-state').textContent,'requested');
  response = { ...get, state:'invented' }; await poll();
  assert.match(node('operator-progress').textContent,/Disconnected/);
  response = { error:'BadRequest', reason:'send a v1 digest and one key artifact' }; code = 400;
  node('operator-signature').value = 'y';
  await events['operator-sign']();
  const [u, p] = sent.find(([,o])=>o.method==='POST');
  assert.equal(u,path+'/confirm');
  assert.deepEqual(JSON.parse(p.body),{v:1,digest:c.digest,sshsig:'y'});
  assert.equal(events.focus,'operator-answer');
  assert.match(node('operator-answer').textContent,/Not confirmed \(HTTP 400\): BadRequest/);
  for (const [refused, statusCode] of [['NotApproved',422],['Expired',409]]) {
    response = {error:refused,reason:'fixture refusal'}; code = statusCode;
    await events['operator-sign']();
    assert(node('operator-answer').textContent.includes(refused));
  }
  response = {...get,state:'expired'}; code = 200; await poll();
  assert.equal(node('operator-controls').hidden,true);
  response = get; await poll();
  code = 200;
  for (const edit of [{result:'fake'},{approvedBy:[]},{approvedBy:''},{state:'invented'}]) {
    response = {...get, result:'confirmed', approvedBy:'webauthn:fixture', ...edit};
    await events['operator-sign']();
    assert.match(node('operator-answer').textContent,/Not confirmed/);
  }
  response = {...get, state:'confirmed', result:'confirmed', approvedBy:'webauthn:fixture'};
  await events['operator-key']();
  assert.equal(asked.publicKey.userVerification,'required');
  assert.equal(asked.publicKey.rpId,'localhost');
  assert.equal(Buffer.from(asked.publicKey.allowCredentials[0].id).toString('base64url'),fixture.credential);
  assert.equal(Buffer.from(asked.publicKey.challenge).toString('hex'),fixture.hash);
  assert.deepEqual(JSON.parse(sent.filter(([,o])=>o.method==='POST').at(-1)[1].body),{v:1,digest:c.digest,webauthn:fixture.assertion});
  assert.match(node('operator-answer').textContent,/Confirmed by webauthn:fixture/);
  assert.equal(node('operator-controls').hidden,true);
  response = {...get,state:'launched',reportHref:'/v1/episodes/e-1/report'}; await poll();
  assert.equal(node('operator-report').href,response.reportHref);
  assert.equal(node('operator-state').textContent,'launched');
  response = {...response,reportHref:{toString:()=>'/v1/episodes/e-1/report'}}; await poll();
  assert.equal(node('operator-report').href,undefined);
  assert.equal(node('operator-report').hidden,true);
})().catch(e=>{ console.error(e); process.exitCode=1; });
`)
}

func TestOperatorCardLinkAllowlist(t *testing.T) {
	runOperatorScript(t, OperatorScript, `
const assert = require('node:assert/strict'), {pathID,statusFor,ceremonyFor,passkeyProof} = require(process.argv[1]);
const id = 'oa-'+'a'.repeat(32), p = '/surface/operator-actions/'+id;
const s = {v:1,id,episode:'e-1',kind:'start',state:'launched',now:1,createdAt:1,expiresAt:601,href:p,statusHref:p+'/status.json'};
const h = '/v1/episodes/e-1/report';
assert.equal(pathID(p),id);
assert.equal(statusFor({...s,reportHref:h}, p).reportHref,h);
for (const bad of [p+'x','x'+p,p+'\n', 'x\n'+p, p+'\nx',p.toUpperCase(),p+'A',p.slice(0,-1),p+'a',p+'?x',p+'#x',p+'/status.json',
  '//localhost'+p,'https://user@localhost:7644'+p,{toString:()=>p}]) assert.equal(pathID(bad),undefined,JSON.stringify(bad));
for (const bad of [h+'x','x'+h,h+'\n',h.toUpperCase(),h+'?q',h+'#f','//localhost'+h,'https://localhost:2'+h,
  'http://user@localhost'+h,'/surface/confirm/vr-0123456789abcdef','https://github.com/sunholo-data/ailang-world/pull/604','/v1/episodes/other/report',{toString:()=>h}]) assert.equal(statusFor({...s,reportHref:bad},p).reportHref,'');
for (const episode of ['A', '_a', '-a', 'a_b', '', 'a'.repeat(65), 'a\n','!a','!\na']) assert.equal(statusFor({...s,episode},p),null);
for (const [k,v] of Object.entries({v:'1',episode:['e-1'],id:id+'x',href:p+'#',statusHref:{href:p},state:'ok',now:'1',expiresAt:null,createdAt:[],runnerSeenAt:'1',terminalAt:'2',reason:1,kind:'unknown'}))
  assert.equal(statusFor({...s,[k]:v},p),null,k);
const c = {post:p+'/confirm',digest:'b'.repeat(64),messageHex:'abcd'};
assert.equal(ceremonyFor(c,p),c);
assert.equal(ceremonyFor(null,p),null);
assert.equal(statusFor(null,p),null);
(async()=>{
 const origin='http://localhost:7644', rp='localhost';
 const get=()=>assert.fail('credential requested without enrolled fixture');
 for (const passkeys of [null,{},[],[{origin,rp:'other',allow:['AQ']}],
   [{origin:'http://localhost:1',rp,allow:['AQ']}],[{origin,rp,allow:[]}],[{origin,rp,allow:{}}]]) {
  await assert.rejects(passkeyProof({...c,passkeys},origin,rp,require('node:crypto').webcrypto,{get}),/No passkey/);
 }
})().catch(e=>{console.error(e);process.exitCode=1;});
for (const [k,v] of [['post',c.post+'#'],['post',{toString:()=>c.post}],['digest','A'.repeat(64)],['digest',1],['digest',[c.digest]],['messageHex','a'],['messageHex',12],['messageHex',[c.messageHex]]])
  assert.equal(ceremonyFor({...c,[k]:v},p),null,k);
`)
}
