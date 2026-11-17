package surfacepage

// Row 243 V3 sprint 2 M5 (design §2.5, §2.9 P4, §6 CP1-CP6): confirm.js under node, on a stand-in DOM
// built from index.html's own ids, against the cards of testdata/confirm-card.golden.jsonl (written by
// TestConfirmCardGolden from the real tools). Each test skips only when no node exists; with
// WORLD_EXEC_NODE set, a harness that cannot run fails.

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// runConfirmHarness runs section of confirmHarness under node and requires its pass line.
func runConfirmHarness(t *testing.T, section string) {
	t.Helper()
	node := os.Getenv("WORLD_EXEC_NODE")
	if node == "" {
		var err error
		node, err = exec.LookPath("node")
		if err != nil {
			t.Skip("node unavailable locally")
		}
	}
	golden, err := os.ReadFile(confirmGolden)
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	for name, data := range map[string][]byte{"confirm.js": ConfirmScript, "index.html": ConfirmPage, "golden.jsonl": golden, "harness.js": []byte(confirmHarness)} {
		if err := os.WriteFile(filepath.Join(dir, name), data, 0o600); err != nil {
			t.Fatal(err)
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, node, filepath.Join(dir, "harness.js"), section)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("node harness %s: %v\n%s", section, err, out)
	}
	if !strings.Contains(string(out), "CONFIRM "+section+" PASS") {
		t.Fatalf("harness manifest missing: %s", out)
	}
}

// CP1, CP2: a fetch that fails, times out (the 10 s abort) or answers non-2xx or non-JSON; an answer for
// another id, or with a malformed digest, challenge or expiry; bytes that are not base64url, do not hash
// to the digest or are not UTF-8 JSON; an unlisted, missing or mistyped field; a path that is not a card's
// → a refusal with its reason, nothing shown, no controls, no message made. Control: the card is ready.
func TestConfirmPageFailsClosed(t *testing.T) { runConfirmHarness(t, "closed") }

// CP3: hostile values (markup, quotes, a NUL in the title and the body) are shown as text, whole, in every
// field, the argv line, the prompt and the body; no HTML sink is written.
func TestConfirmPageRendersText(t *testing.T) { runConfirmHarness(t, "text") }

// CP4, CP6: the page's message is Go's VerbMessage (86 bytes) and its hash Go's VerbMessageHash; Save
// the message is exactly those bytes as world-verb-<id>.msg; the signing line names world-verb@sunholo.
func TestConfirmPageMessage(t *testing.T) { runConfirmHarness(t, "message") }

// CP5: an approved, superseded, expired or unknown state shows that state and the fields, and no controls.
func TestConfirmPageStates(t *testing.T) { runConfirmHarness(t, "states") }

const confirmHarness = `
const assert = require('node:assert/strict');
const fs = require('node:fs');
const vm = require('node:vm');
const nodeCrypto = require('node:crypto');
const SRC = fs.readFileSync('confirm.js', 'utf8'), HTML = fs.readFileSync('index.html', 'utf8');
const G = Object.fromEntries(fs.readFileSync('golden.jsonl', 'utf8').trim().split('\n').map((l) => JSON.parse(l)).map((x) => [x.name, x]));
const sha = (b) => nodeCrypto.createHash('sha256').update(b).digest('hex');
const bytesOf = (card) => Buffer.from(card.request, 'base64url');
const reqOf = (card) => JSON.parse(bytesOf(card).toString('utf8'));
// withBytes is card carrying bytes b, its digest made to match.
const withBytes = (card, b) => Object.assign({}, card, { request: Buffer.from(b).toString('base64url'), digest: sha(b) });

// A stand-in DOM: index.html's elements by id (with their hidden attribute), and created ones.
class El {
  constructor(doc, tag, id, hidden) { Object.assign(this, { doc, tagName: tag, id: id || '', hidden: !!hidden, dataset: {}, attrs: {}, kids: [], text: '' }); }
  set textContent(v) { this.text = String(v); this.kids = []; }
  get textContent() { return this.text + this.kids.map((k) => k.textContent).join(''); }
  set innerHTML(v) { this.doc.sinks.push(String(v)); }
  set outerHTML(v) { this.doc.sinks.push(String(v)); }
  insertAdjacentHTML(_, v) { this.doc.sinks.push(String(v)); }
  append(...k) { this.kids.push(...k); }
  toggleAttribute(n, on) { if (on) this.attrs[n] = ''; else delete this.attrs[n]; }
  hasAttribute(n) { return n in this.attrs; }
}
function makeDoc() {
  const doc = { sinks: [], byId: {} };
  doc.body = new El(doc, 'body');
  doc.body.dataset.state = HTML.match(/<body data-state="([a-z]+)">/)[1];
  for (const m of HTML.matchAll(/<([a-z0-9]+) id="([a-z-]+)"( hidden)?/g)) doc.byId[m[2]] = new El(doc, m[1], m[2], !!m[3]);
  const find = (els, id) => { for (const e of els) { if (e.id === id) return e; const k = find(e.kids, id); if (k) return k; } return null; };
  doc.getElementById = (id) => doc.byId[id] || find(Object.values(doc.byId), id);
  doc.createElement = (tag) => new El(doc, tag);
  return doc;
}
const answer = (status, body) => () => Promise.resolve({ status, text: async () => (typeof body === 'string' ? body : JSON.stringify(body)) });
const hang = (init) => new Promise((_, reject) => init.signal.addEventListener('abort', () => reject(new DOMException('The operation was aborted.', 'AbortError'))));

// open runs confirm.js at path against fetch; fire runs the page's pending timers once it has asked.
async function open(path, fetch, fire) {
  const doc = makeDoc(), timers = [], blobs = [], fetches = [];
  vm.runInNewContext(SRC, { document: doc, location: { pathname: path, origin: 'http://localhost:1' }, crypto: globalThis.crypto, TextEncoder, TextDecoder, atob, btoa, Blob, AbortController, console,
    URL: { createObjectURL: (b) => { blobs.push(b); return 'blob:http://localhost:1/' + blobs.length; } },
    setTimeout: (f, ms) => timers.push({ f, ms, cleared: false }), clearTimeout: (n) => { if (timers[n - 1]) timers[n - 1].cleared = true; },
    fetch: (url, init) => { fetches.push({ url, init }); return fetch(init); } });
  for (let i = 0; i < 400 && doc.body.dataset.state === 'loading'; i++) {
    if (fire && i === 20) for (const t of timers) if (!t.cleared) t.f();
    await new Promise((r) => setImmediate(r));
  }
  const $ = (id) => doc.getElementById(id);
  return { doc, $, timers, blobs, fetches, state: doc.body.dataset.state, what: $('what').textContent };
}
const card = (c) => '/surface/confirm/' + c.id;
// refused asserts a fail-closed page: state, the reason, nothing shown, no controls, no message made.
function refused(p, state, reason, name) {
  assert.equal(p.state, state, name + ': ' + p.what);
  assert(p.what.includes(reason), name + ': ' + p.what);
  assert(p.$('what').hasAttribute('data-error'), name);
  assert(p.$('shown').hidden && p.$('controls').hidden && p.blobs.length === 0 && p.$('sshline').textContent === '' && p.$('prompt').textContent === '', name);
}

const sections = {
  async closed() {
    const c = G.draft.card, other = G.hostile.card;
    const ok = await open(card(c), answer(200, c));
    assert.equal(ok.state, 'ready', ok.what);
    assert.equal(ok.fetches.length, 1);
    assert.equal(ok.fetches[0].url, '/surface/confirmations/' + c.id);
    assert(ok.fetches[0].init.cache === 'no-store' && ok.fetches[0].init.credentials === 'omit' && ok.fetches[0].init.signal);
    assert(ok.timers.length === 1 && ok.timers[0].ms === 10000 && ok.timers[0].cleared, JSON.stringify(ok.timers));
    refused(await open(card(c), () => Promise.reject(new TypeError('Failed to fetch'))), 'failed', 'The request could not be read (Failed to fetch).', 'network');
    const slow = await open(card(c), hang, true);
    refused(slow, 'failed', 'The request could not be read (no answer within 10 s).', 'timeout');
    assert(slow.timers.length === 1 && slow.timers[0].ms === 10000);
    refused(await open(card(c), answer(404, { reason: 'no such verb request' })), 'failed', 'The daemon did not give the request (HTTP 404: no such verb request).', '404');
    refused(await open(card(c), answer(503, 'down')), 'failed', '(HTTP 503)', '503');
    refused(await open(card(c), answer(200, '<html>')), 'failed', '(HTTP 200)', 'not JSON');
    const malformed = 'The daemon\'s answer is malformed; nothing to confirm.';
    for (const [name, edit] of [['another id', { id: other.id }], ['digest case', { digest: c.digest.toUpperCase() }], ['short challenge', { challenge: c.challenge.slice(1) }],
      ['negative expiry', { expires: -1 }], ['fractional expiry', { expires: c.expires + 0.5 }], ['no state', { state: 1 }]]) {
      refused(await open(card(c), answer(200, Object.assign({}, c, edit))), 'refused', malformed, name);
    }
    refused(await open(card(c), answer(200, Object.assign({}, c, { request: '*' + c.request }))), 'refused', 'is not base64url', 'base64url');
    const digest = 'The request\'s bytes do not match its digest; nothing to confirm.';
    refused(await open(card(c), answer(200, Object.assign({}, c, { digest: sha('other') }))), 'refused', digest, 'digest');
    refused(await open(card(c), answer(200, Object.assign({}, c, { request: other.request }))), 'refused', digest, 'swapped bytes');
    const padded = other.request + '='.repeat((4 - other.request.length % 4) % 4); // standard base64's padding: atob takes it
    assert(padded !== other.request && Buffer.from(padded, 'base64').equals(bytesOf(other)));
    refused(await open(card(other), answer(200, Object.assign({}, other, { request: padded }))), 'refused', 'is not base64url', 'padded base64');
    const latin = Buffer.from(JSON.stringify(Object.assign({}, reqOf(c), { title: 'TITLE' })).replace('TITLE', '\\u00ff'), 'utf8');
    refused(await open(card(c), answer(200, withBytes(c, Buffer.concat([latin.subarray(0, latin.indexOf('\\u00ff')), Buffer.from([0xff]), latin.subarray(latin.indexOf('\\u00ff') + 6)])))),
      'refused', 'is not UTF-8 JSON', 'a lone 0xff in the title');
    const r = reqOf(c), unshown = 'This request has a field the page cannot show (';
    for (const [name, obj, key] of [['unlisted', Object.assign({}, r, { extra: 'x' }), 'extra'], ['missing', Object.assign({}, r, { committer: undefined }), 'committer'],
      ['argv a string', Object.assign({}, r, { argv: 'ailang-worldd episode finish w1' }), 'argv'], ['protected a number', Object.assign({}, r, { protected: [1] }), 'protected'],
      ['title a number', Object.assign({}, r, { title: 7 }), 'title'], ['body null', Object.assign({}, r, { body: null }), 'body'], ['an array', [r], 'the request']]) {
      refused(await open(card(c), answer(200, withBytes(c, JSON.stringify(obj)))), 'refused', unshown + key + '); nothing to confirm here; confirm at the terminal.', name);
    }
    const lost = await open('/surface/confirm/vr-0123', answer(200, c));
    refused(lost, 'refused', 'This is not a confirmation card\'s address; nothing to confirm.', 'path');
    assert.equal(lost.fetches.length, 0);
  },
  async text() {
    for (const name of ['hostile', 'draft']) {
      const c = G[name].card, r = reqOf(c), p = await open(card(c), answer(200, c));
      assert.equal(p.state, 'ready', p.what);
      assert.equal(p.$('argv').textContent, G[name].argvLine);
      assert.equal(p.$('prompt').textContent, r.prompt);
      for (const k of ['schema', 'verb', 'episode', 'prompt', 'scope', 'title', 'committer']) assert.equal(p.$('f-' + k).textContent, r[k], k);
      for (const k of ['argv', 'protected']) assert.equal(p.$('f-' + k).textContent, r[k].join('\n'), k);
      assert.equal(p.$('f-body').textContent, r.body);
      assert.equal(p.$('f-body-size').textContent, Buffer.byteLength(r.body) + ' bytes' + 'the full body' + r.body);
      assert.equal(p.$('m-digest').textContent, c.digest);
      assert.equal(p.$('m-challenge').textContent, c.challenge);
      assert.equal(p.$('m-expires').textContent, new Date(c.expires * 1000).toISOString() + ' (' + c.expires + ')');
      assert.deepEqual(p.doc.sinks, []);
    }
    const r = reqOf(G.hostile.card), p = await open(card(G.hostile.card), answer(200, G.hostile.card));
    // The NUL rows (S1 Appendix C, P6): the text after a NUL is shown, in the title, the prompt and the body.
    assert(r.title.includes('a\u0000b & ') && r.body.includes('NUL \u0000 here &amp; after'));
    assert(p.$('f-title').textContent.endsWith('a\u0000b & \'q\''));
    assert(p.$('prompt').textContent.includes('title: <img src=x onerror=alert(1)> a\u0000b & \'q\'\n'));
    assert(p.$('f-body').textContent.includes('</pre><script>alert(2)</script>\nNUL \u0000 here &amp; after — ünï\n'));
    assert(Buffer.byteLength(r.body) > r.body.length); // the size line counts bytes, not characters
    assert(p.$('argv').textContent.endsWith(' --title \'it\'"\'"\'s <b>\' --ready'), p.$('argv').textContent);
  },
  async message() {
    for (const name of ['draft', 'hostile']) {
      const g = G[name], c = g.card, p = await open(card(c), answer(200, c));
      assert.equal(p.state, 'ready', p.what);
      assert.equal(p.$('m-hash').textContent, g.hash);
      assert(p.what.endsWith('; message hash ' + g.hash), p.what);
      assert.equal(p.blobs.length, 1);
      assert.equal(p.blobs[0].type, 'application/octet-stream');
      const m = Buffer.from(await p.blobs[0].arrayBuffer());
      assert.equal(m.length, 86);
      assert.equal(m.toString('hex'), g.message);
      assert.equal(sha(m), g.hash);
      assert.equal(p.$('save').href, 'blob:http://localhost:1/1');
      assert.equal(p.$('save').download, 'world-verb-' + c.id + '.msg');
      assert.equal(p.$('sshline').textContent, 'ssh-keygen -Y sign -n world-verb@sunholo -f <key> world-verb-' + c.id + '.msg');
      assert(!p.$('controls').hidden && !p.$('shown').hidden);
    }
  },
  async states() {
    const ready = await open(card(G.draft.card), answer(200, G.draft.card));
    assert(ready.state === 'ready' && !ready.$('controls').hidden && ready.blobs.length === 1);
    for (const [c, status] of [[G['draft-approved'].card, 'Already confirmed by webauthn:mark-mac: decision sha256:' + '5e'.repeat(32) + '; nothing to confirm here.'],
      [G['draft-superseded'].card, 'Superseded by a newer finish of this episode; nothing to confirm here.'],
      [G['draft-expired'].card, 'Expired; nothing to confirm here.'],
      [Object.assign({}, G.draft.card, { state: 'odd' }), 'This request is odd; nothing to confirm here.']]) {
      const p = await open(card(c), answer(200, c));
      assert.equal(p.state, 'closed', c.state + ': ' + p.what);
      assert.equal(p.$('status').textContent, status);
      assert(p.$('controls').hidden && p.blobs.length === 0 && p.$('sshline').textContent === '' && p.$('save').download === undefined, c.state);
      assert(!p.$('shown').hidden && p.$('f-title').textContent === reqOf(c).title, c.state);
    }
    assert.deepEqual([G['draft-approved'].card.state, G['draft-superseded'].card.state, G['draft-expired'].card.state], ['approved', 'superseded', 'expired']);
  },
};
const which = process.argv[2];
sections[which]().then(() => console.log('CONFIRM ' + which + ' PASS'), (e) => { console.error(e); process.exit(1); });
`
