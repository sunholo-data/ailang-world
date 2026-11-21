package surfacepage

import (
	"context"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

func TestNodeTrajectoryLoadsLazilyAndVerifiesLocalTranscript(t *testing.T) {
	node, err := exec.LookPath("node")
	if err != nil {
		t.Skip("node is required for the trajectory interaction test")
	}
	script, err := filepath.Abs(filepath.Join("app", "node-trajectory.js"))
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	out, err := exec.CommandContext(ctx, node, "-e", nodeTrajectoryHarness, script).CombinedOutput()
	if err != nil {
		t.Fatalf("node trajectory: %v\n%s", err, out)
	}
}

const nodeTrajectoryHarness = `
const assert = require('node:assert/strict');
const {webcrypto, createHash} = require('node:crypto');
const {trajectoryPath, rawReceipt, verifyTranscript, transcriptPages} = require(process.argv[1]);
const {mount, clear} = globalThis.WorldNodeTrajectory;
const transcript = '{"kind":"<private>"}\r\n' + 'x'.repeat(11975) + '💡'.repeat(200) + '\n';
const bytes = new TextEncoder().encode(transcript);
const receipt = {filename:'verified.transcript.jsonl', sha256:createHash('sha256').update(bytes).digest('hex')};
const file = {name:receipt.filename, size:bytes.length, arrayBuffer:async () => bytes.buffer};
const config = {id:'world:episode:verified', brane:'world:brane:open', run:{episode:'verified'}};
const endpoint = '/surface/trajectory.json?brane=world%3Abrane%3Aopen&at=world%3Aepisode%3Averified';
assert.equal(trajectoryPath(config.brane, config.id), endpoint);
assert(rawReceipt(receipt, 'verified'));
assert(!rawReceipt({...receipt, filename:'../../private'}, 'verified'));
assert(!rawReceipt({...receipt, sha256:'unverified'}, 'verified'));

class Node {
  constructor(tag) {
    this.tag = tag; this.ownerDocument = document; this.children = []; this.events = {}; this.textContent = ''; this.open = false; this.isConnected = true;
  }
  append(...children) { this.children.push(...children); }
  replaceChildren() { this.children = []; }
  setAttribute() {}
  addEventListener(name, fn) { this.events[name] = fn; }
  fire(name, data) { return this.events[name]?.(data); }
  querySelector(selector) { return this.matches?.[selector] ?? this.find(selector.split(' ').at(-1)); }
  find(tag) {
    for (const child of this.children) {
      if (child instanceof Node) {
        if (child.tag === tag) return child;
        const nested = child.find(tag);
        if (nested) return nested;
      }
    }
  }
  text() { return this.textContent + this.children.map(child => typeof child === 'string' ? child : child.text()).join(' '); }
}
const document = {createElement: tag => new Node(tag)};
function disclosure() {
  const front = new Node('article'), details = new Node('details'), body = new Node('div');
  front.matches = {'[data-trajectory-body]': body};
  details.matches = {'[data-trajectory-body]': body};
  return {front, details, body};
}

(async () => {
  assert.equal(await verifyTranscript(file, receipt, webcrypto), transcript);
  const pages = transcriptPages(transcript);
  assert.equal(pages.join(''), transcript);
  assert(pages.length > 1 && pages.every(page => page.length <= 12000 && !/[\uD800-\uDBFF]$/.test(page)));
  await assert.rejects(verifyTranscript({...file, name:'wrong.jsonl'}, receipt, webcrypto), /Choose verified/);
  await assert.rejects(verifyTranscript({...file, size:9*1024*1024}, receipt, webcrypto), /8 MiB/);
  await assert.rejects(verifyTranscript(file, {...receipt, sha256:'0'.repeat(64)}, webcrypto), /does not match/);

  const requests = [];
  const getJSON = (url, signal) => new Promise(resolve => requests.push({url, signal, resolve}));
  const first = disclosure();
  mount(first.front, config, getJSON, false);
  assert.equal(requests.length, 0, 'closed disclosure should not fetch');
  mount(first.front, config, getJSON, true);
  assert.equal(requests.length, 1);
  assert.equal(requests[0].url, endpoint);

  const refresh = disclosure();
  mount(refresh.front, config, getJSON, true);
  assert.equal(requests.length, 1, 'same selection refresh reuses inflight request');
  const result = {ok:true, body:{schema:'world/surface-trajectory/v1', episode:'verified',
    world:{available:true, calls:[{tool:'read_checkpoint', surface:'mcp', entry:7}], truncated:true},
    harness:{available:true, events:[{launch:1, kind:'turn-completed'}]}, rawTranscript:receipt}};
  requests[0].resolve(result);
  await new Promise(setImmediate);
  assert.match(refresh.body.text(), /World tool calls.*read_checkpoint.*entry 7/);
  assert.match(refresh.body.text(), /Only the first recorded entries/);
  assert.match(refresh.body.text(), /Harness events.*Turn completed/);
  assert.match(refresh.body.text(), /~\/.ailang\/world\/sessions\/verified\.transcript\.jsonl/);
  assert(!refresh.body.text().includes('<private>'), 'raw text must not come from a fetch');
  if (!globalThis.crypto) Object.defineProperty(globalThis, 'crypto', {value:webcrypto});
  const raw = refresh.body.find('details');
  raw.open = true;
  raw.fire('toggle');
  const input = raw.find('input');
  input.files = [file];
  await input.fire('change');
  assert.match(refresh.body.find('pre').textContent, /<private>/, 'local match reveals verbatim text');
  assert.equal(requests.length, 1, 'raw file selection has no network request');
  const localRefresh = disclosure();
  mount(localRefresh.front, config, getJSON, true);
  assert(localRefresh.body.find('details').open, 'raw disclosure survives same-node refresh');
  assert.equal(localRefresh.body.find('pre').textContent, pages[0]);

  const second = disclosure();
  mount(second.front, {...config, id:'world:episode:next', run:{episode:'next'}}, getJSON, true);
  assert.equal(requests.length, 2);
  clear();
  assert(requests[1].signal.aborted, 'leaving a selection should cancel its request');
  requests[1].resolve(result);
  await new Promise(setImmediate);
  assert(!second.body.text().includes('read_checkpoint'), 'stale evidence must not render');
})().catch(error => { console.error(error); process.exitCode = 1; });
`
