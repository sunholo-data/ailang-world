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

// CP8 (sprint 4): Approve with your passkey asks navigator.credentials.get for the message hash as the
// challenge, the origin's group's rp and credentials, UV required, A2's 240 s abort; POSTs the
// assertion file with the card's digest; shows the answer and reads the card again. No group for the
// origin: no passkey button. An aborted or failed ceremony sends nothing and the controls return.
func TestConfirmPageSubmitsAPasskey(t *testing.T) { runConfirmHarness(t, "passkey") }

// CP9 (sprint 4): a refusal (422, 409, 500) or a lost POST is shown as not confirmed, never as a
// confirmation, and the card is read again: still pending, the controls are back.
func TestConfirmPageShowsRefusals(t *testing.T) { runConfirmHarness(t, "answers") }

// CP10 (sprint 4): the pasted SSH signature is POSTed in ParseApproveBody's grammar, {v, digest, sshsig}.
func TestConfirmPageSubmitsAnSSHSIG(t *testing.T) { runConfirmHarness(t, "sshsig") }

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
  constructor(doc, tag, id, hidden) { Object.assign(this, { doc, tagName: tag, id: id || '', hidden: !!hidden, dataset: {}, attrs: {}, kids: [], text: '', on: {} }); }
  addEventListener(t, f) { (this.on[t] = this.on[t] || []).push(f); }
  click() { if (!this.disabled && !this.hidden) for (const f of this.on.click || []) f(); }
  focus() { this.doc.activeElement = this; }
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
  doc.documentElement = new El(doc, 'html');
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

// open runs confirm.js at path against fetch and waits, on the wall clock (at most 10 s), until the page leaves
// 'loading': its SHA-256 digests resolve off the event loop, so no count of ticks bounds them (iteration 401 M7);
// fire runs the page's stubbed timers once, as soon as it has asked (its fetch made and its abort timer set).
async function open(path, fetch, fire, navigator, hash = '') {
  const doc = makeDoc(), timers = [], blobs = [], fetches = [];
  vm.runInNewContext(SRC, { document: doc, URLSearchParams, localStorage: {getItem: () => JSON.stringify({theme: 'slate'})}, location: { pathname: path, hash, origin: 'http://localhost:1' }, crypto: globalThis.crypto, TextEncoder, TextDecoder, atob, btoa, Blob, AbortController, console, navigator,
    URL: { createObjectURL: (b) => { blobs.push(b); return 'blob:http://localhost:1/' + blobs.length; } },
    setTimeout: (f, ms) => timers.push({ f, ms, cleared: false }), clearTimeout: (n) => { if (timers[n - 1]) timers[n - 1].cleared = true; },
    fetch: (url, init) => { fetches.push({ url, init }); return fetch(init); } });
  for (const end = Date.now() + 10000; doc.body.dataset.state === 'loading' && Date.now() < end;) {
    if (fire && fetches.length && timers.length) { fire = false; for (const t of timers) if (!t.cleared) t.f(); }
    await new Promise((r) => setTimeout(r, 2)); // node's own timer: the page's are the stubs above
  }
  const $ = (id) => doc.getElementById(id);
  return { doc, $, timers, blobs, fetches, state: doc.body.dataset.state, what: $('what').textContent };
}
const card = (c) => '/surface/confirm/' + c.id;
// until waits, on the wall clock (at most 10 s), for ok() after a click: the page's work is promises.
async function until(ok, name) {
  for (const end = Date.now() + 10000; !ok();) {
    assert(Date.now() < end, 'timed out: ' + name);
    await new Promise((r) => setTimeout(r, 2));
  }
}
// daemon answers the card's GETs with cards in turn (the last repeats) and its POSTs with post.
const daemon = (cards, post) => (init) => (init.method === 'POST' ? post(init) : answer(200, cards.length > 1 ? cards.shift() : cards[0])());
const posts = (p) => p.fetches.filter((f) => f.init.method === 'POST');
const CRED = { id: 'Bwc', rawId: new Uint8Array([7, 7]).buffer, type: 'public-key',
  response: { clientDataJSON: Buffer.from('cd'), authenticatorData: Buffer.from('ad'), signature: Buffer.from('sg'), userHandle: null } };
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
    assert.equal(ok.doc.documentElement.dataset.theme, 'slate');
    for (const [fragment, target] of [['#activity', '/surface#activity'],
      ['#brane=row-218&at=step%2F1&return=https://elsewhere.invalid', '/surface?brane=row-218&at=step%2F1'],
      ['#brane=world%3Abrane%3Ashared&at=world%3Aepisode%3Aw1&view=heads%2Ftrajectory', '/surface?brane=world%3Abrane%3Ashared&at=world%3Aepisode%3Aw1#heads/trajectory']]) {
      const contextual = await open(card(c), answer(200, c), false, undefined, fragment);
      assert.equal(contextual.$('back').href, target);
      assert.equal(contextual.fetches[0].url, '/surface/confirmations/' + c.id);
    }
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
    for (const scope of [r.scope + '\n', r.scope.replace('&draft=true', '&draft=false&draft=true')]) {
      const bad = withBytes(c, JSON.stringify({ ...r, scope }));
      refused(await open(card(bad), answer(200, bad)), 'refused', 'does not describe a recognizable pull request', 'scope');
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
      assert.equal(p.$('f-body-size').textContent, '(' + Buffer.byteLength(r.body) + ' bytes)');
      assert.equal(p.$('review-title').textContent, r.title);
      assert.equal(p.$('episode').textContent, r.episode);
      assert.equal(p.$('repository').textContent, 'sunholo-data/ailang-world');
      assert.equal(p.$('head').textContent, 'world/w1');
      assert.equal(p.$('head-commit').textContent, 'ab'.repeat(20));
      assert.equal(p.$('base').textContent, 'dev');
      assert(p.$('review-facts').open && p.$('receipt-actions').hidden);
      assert.equal(p.$('action').textContent, name === 'draft' ? 'Open a draft pull request' : 'Open a pull request');
      assert.equal(p.$('pr-state').textContent, name === 'draft' ? 'Draft' : 'Ready for review');
      assert.equal(p.$('label').textContent, name === 'draft' ? 'ailang-world' : 'ailang-world (will be created)');
      assert.deepEqual(p.$('review-protected').kids.map((el) => el.textContent), r.protected);
      assert.equal(p.$('no-protected').hidden, r.protected.length !== 0);
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
      assert.equal(p.$('status').textContent, 'Awaiting approval');
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
    for (const [c, status] of [[G['draft-approved'].card, 'Approval recorded'],
      [G['draft-superseded'].card, 'Superseded by a newer finish'],
      [G['draft-expired'].card, 'Request expired'],
      [Object.assign({}, G.draft.card, { state: 'odd' }), 'Request closed']]) {
      const p = await open(card(c), answer(200, c));
      assert.equal(p.state, 'closed', c.state + ': ' + p.what);
      assert.equal(p.$('status').textContent, status);
      if (c.state === 'approved') {
        assert.equal(p.what, 'mark-mac approved with a passkey. The finish command reports the pull request outcome.');
        assert.equal(p.$('approval-decision').textContent, c.decision);
        assert.equal(p.$('approval-signer').textContent, c.approvedBy);
        assert.equal(p.$('action').textContent, 'Requested a draft pull request');
        assert.equal(p.$('receipt-back').href, '/surface');
        assert.equal(p.$('receipt-back').textContent, 'Return to World');
        assert(!p.$('receipt-actions').hidden && p.$('receipt-pr').hidden);
      } else assert(p.$('receipt-actions').hidden);
      assert(!p.$('review-facts').open);
      assert(p.$('controls').hidden && p.blobs.length === 0 && p.$('sshline').textContent === '' && p.$('save').download === undefined, c.state);
      assert(!p.$('shown').hidden && p.$('f-title').textContent === reqOf(c).title, c.state);
    }
    assert.deepEqual([G['draft-approved'].card.state, G['draft-superseded'].card.state, G['draft-expired'].card.state], ['approved', 'superseded', 'expired']);
    const approved = G['draft-approved'].card;
    const episodeHref = '/surface?at=world%3Aepisode%3Aw1&brane=world%3Abrane%3Ashared';
    const prHref = 'https://github.com/sunholo-data/ailang-world/pull/123';
    const linked = await open(card(approved), answer(200, {...approved, episodeHref, prHref}));
    assert.equal(linked.$('receipt-back').href, '/surface?brane=world%3Abrane%3Ashared&at=world%3Aepisode%3Aw1');
    assert.equal(linked.$('receipt-back').textContent, 'Return to episode');
    assert(!linked.$('receipt-pr').hidden && linked.$('receipt-pr').href === prHref);
    assert.equal(linked.$('receipt-pr').textContent, 'View PR #123 on GitHub');
    assert.equal(linked.$('status').textContent, 'Pull request opened');
    assert.equal(linked.$('action').textContent, 'Draft pull request');
    assert.equal(linked.what, 'mark-mac approved with a passkey.');
    const ssh = await open(card(approved), answer(200, {...approved, approvedBy:'ssh:mark', episodeHref, prHref}));
    assert.equal(ssh.what, 'mark approved with an SSH signature.');
    assert.equal(ssh.$('approval-signer').textContent, 'ssh:mark');
    const pendingWithLinks = await open(card(G.draft.card), answer(200, {...G.draft.card, episodeHref, prHref}));
    assert(pendingWithLinks.$('receipt-actions').hidden && pendingWithLinks.$('receipt-pr').hidden);
    for (const mismatched of [episodeHref.replace('%3Aw1', '%3Aother'), episodeHref.replace('world%3Aepisode', 'other%3Aepisode'),
      episodeHref + '&at=world%3Aepisode%3Aw1', 'https://elsewhere.invalid' + episodeHref]) {
      const refusedLinks = await open(card(approved), answer(200, {...approved, episodeHref:mismatched, prHref:prHref.replace('ailang-world', 'another-repo')}));
      assert.equal(refusedLinks.$('receipt-back').href, '/surface');
      assert(refusedLinks.$('receipt-pr').hidden);
    }

  },
  async passkey() {
    const g = G.draft, c = g.card, approved = G['draft-approved'].card, asks = [];
    const nav = { credentials: { get: (o) => { asks.push(o); return Promise.resolve(CRED); } } };
    const ok = { v: 1, confirmed: true, request: c.id, approvedBy: approved.approvedBy, decision: approved.decision, next: 'the waiting episode finish opens the pull request' };
    const p = await open(card(c), daemon([c, approved], answer(200, ok)), false, nav);
    assert(p.state === 'ready' && !p.$('go').hidden && !p.$('go').disabled, p.what);
    p.$('go').click();
    await until(() => p.state === 'closed' || p.doc.body.dataset.state === 'closed', 'passkey');
    assert.equal(asks.length, 1);
    const pk = asks[0].publicKey;
    assert.equal(Buffer.from(pk.challenge).toString('hex'), g.hash); // the message hash, never the mint hash
    assert.deepEqual([pk.rpId, pk.userVerification, pk.timeout, pk.allowCredentials.length, pk.allowCredentials[0].type, Buffer.from(pk.allowCredentials[0].id).toString('hex')],
      ['localhost', 'required', 240000, 1, 'public-key', '0707']);
    assert(asks[0].signal && p.timers.some((t) => t.ms === 240000 && t.cleared), 'the ceremony\'s abort timer is cleared');
    const sent = posts(p);
    assert.equal(sent.length, 1);
    assert.equal(sent[0].url, '/surface/confirm/' + c.id);
    assert.equal(sent[0].init.headers['Content-Type'], 'application/json');
    assert.deepEqual(JSON.parse(sent[0].init.body), { v: 1, digest: c.digest, webauthn: { id: 'Bwc', rawId: 'Bwc', type: 'public-key',
      response: { clientDataJSON: 'Y2Q', authenticatorData: 'YWQ', signature: 'c2c', userHandle: null } } });
    assert.equal(p.$('approval-decision').textContent, approved.decision);
    assert.equal(p.$('answer').textContent, '');
    assert(p.$('controls').hidden && p.fetches.length === 3 && p.$('status').textContent === 'Approval recorded', p.$('status').textContent);
    assert.equal(p.doc.activeElement, p.$('status'));
    const elsewhere = [{ origin: 'https://other.tail0.ts.net', rp: 'other.tail0.ts.net', allow: ['Bwc'] }]; // a group for another origin only
    const none = await open(card(c), answer(200, Object.assign({}, c, { passkeys: elsewhere })), false, nav);
    assert(none.state === 'ready' && none.$('go').hidden && none.$('go').disabled && !none.$('post-sig').disabled);
    assert(none.$('ssh-fallback').open && none.what.includes('without --confirm-on-surface'));
    for (const [name, get, fire] of [['aborted at 240 s', (o) => new Promise((_, reject) => o.signal.addEventListener('abort', () => reject(new DOMException('The operation was aborted.', 'AbortError')))), true],
      ['refused by the browser', () => Promise.reject(new DOMException('The operation is not allowed.', 'NotAllowedError')), false]]) {
      const q = await open(card(c), daemon([c], answer(500, {})), false, { credentials: { get } });
      q.$('go').click();
      if (fire) {
        await until(() => q.timers.some((t) => t.ms === 240000), name + ' timer');
        assert.equal(q.$('status').textContent, 'Waiting for your passkey…');
        assert(q.doc.body.dataset.state === 'waiting' && q.$('go').disabled && q.$('post-sig').disabled, name + ': waiting for the passkey');
        for (const t of q.timers) if (t.ms === 240000 && !t.cleared) t.f();
      }
      await until(() => q.$('answer').textContent !== '' && q.doc.body.dataset.state === 'ready', name);
      assert(q.$('answer').textContent.startsWith('Not confirmed: the browser made no passkey response ('), q.$('answer').textContent);
      assert(posts(q).length === 0 && !q.$('controls').hidden && !q.$('go').disabled, name);
    }
  },
  async answers() {
    const c = G.draft.card;
    for (const [name, post, text] of [
      ['422', answer(422, { v: 1, error: 'NotApproved', reason: 'the SSH signature does not confirm request' }), 'Not confirmed (HTTP 422 NotApproved): the SSH signature does not confirm request'],
      ['409', answer(409, { v: 1, error: 'Conflict', reason: 'the request changed since the card was shown; reload it' }), 'Not confirmed (HTTP 409 Conflict): the request changed since the card was shown; reload it'],
      ['500', answer(500, { v: 1, error: 'Internal', reason: 'the confirmation did not complete; approve again' }), 'Not confirmed (HTTP 500 Internal): the confirmation did not complete; approve again'],
      ['200 not confirmed', answer(200, { v: 1, confirmed: false }), 'Not confirmed (HTTP 200 ): no reason given'],
      ['500 that says confirmed', answer(500, { v: 1, confirmed: true, approvedBy: 'ssh:x', decision: 'd', next: 'n', error: 'Internal', reason: 'r' }), 'Not confirmed (HTTP 500 Internal): r'],
      ['lost', () => Promise.reject(new TypeError('Failed to fetch')), 'Confirmation outcome unknown: Failed to fetch. Reading the request again.']]) {
      const p = await open(card(c), daemon([c], post));
      p.$('sig').value = 'sig';
      p.$('post-sig').click();
      await until(() => posts(p).length === 1 && p.$('answer').textContent !== '' && p.doc.body.dataset.state === 'ready', name);
      assert.equal(p.$('answer').textContent, text, name);
      assert.equal(p.doc.activeElement, p.$('status'), name);
      assert(!p.$('controls').hidden && !p.$('post-sig').disabled && p.fetches.length === 3, name);
    }
  },
  async sshsig() {
    const c = G.draft.card, approved = G['draft-approved'].card;
    let release;
    const held = new Promise((r) => { release = r; });
    const p = await open(card(c), daemon([c, approved], () => held.then(answer(200, { v: 1, confirmed: true, request: c.id, approvedBy: 'ssh:mark', decision: 'sha256:' + '6f'.repeat(32), next: 'n' }))));
    const armored = '-----BEGIN SSH SIGNATURE-----\nU1NIU0lH\n-----END SSH SIGNATURE-----\n';
    p.$('sig').value = armored;
    p.$('post-sig').click();
    await until(() => posts(p).length === 1, 'sshsig sent');
    assert(p.doc.body.dataset.state === 'posting' && p.$('post-sig').disabled && p.$('go').disabled, 'no second submit while one is in flight');
    assert.equal(p.$('status').textContent, 'Submitting approval…');
    release();
    await until(() => p.doc.body.dataset.state === 'closed', 'sshsig');
    const sent = posts(p);
    assert(sent.length === 1 && sent[0].url === '/surface/confirm/' + c.id && sent[0].init.credentials === 'omit');
    assert.deepEqual(JSON.parse(sent[0].init.body), { v: 1, digest: c.digest, sshsig: armored });
    assert.equal(p.$('answer').textContent, 'Approval recorded by ssh:mark. Decision sha256:' + '6f'.repeat(32) + '. Check the episode finish at the terminal for the PR outcome.');
  },
};
const which = process.argv[2];
sections[which]().then(() => console.log('CONFIRM ' + which + ' PASS'), (e) => { console.error(e); process.exit(1); });
`
