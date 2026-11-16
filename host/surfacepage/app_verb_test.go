package surfacepage

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestSurfaceAppVerbs (row 243 V1, design §2.3, §6 VP1–VP6) runs the bar's verb part of app.js under
// node against GET /surface/verb.json's real answers (testdata/verb-route.golden.json, which
// host/daemon's TestSurfaceVerbOutcomes checks byte for byte): only a schema-valid address answer
// reaches the address path; every failure, timeout and superseded ask leaves the view and the input.
// It skips only when no node exists; with WORLD_EXEC_NODE set, a harness that cannot run fails.
func TestSurfaceAppVerbs(t *testing.T) {
	node := os.Getenv("WORLD_EXEC_NODE")
	if node == "" {
		var err error
		node, err = exec.LookPath("node")
		if err != nil {
			t.Skip("node unavailable locally")
		}
	}
	dir := t.TempDir()
	files := map[string][]byte{"app.js": AppScript, "surface.js": Script, "harness.js": []byte(verbHarness)}
	for name, src := range map[string]string{"golden.json": filepath.Join("testdata", "verb-route.golden.json"),
		"index.json": filepath.Join("..", "surface", "testdata", "index-rig.json")} {
		b, err := os.ReadFile(src)
		if err != nil {
			t.Fatal(err)
		}
		files[name] = b
	}
	for name, data := range files {
		if err := os.WriteFile(filepath.Join(dir, name), data, 0o600); err != nil {
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
	if !strings.Contains(string(out), "APP VERBS PASS") {
		t.Fatalf("harness manifest missing: %s", out)
	}
}

const verbHarness = `
const assert = require('node:assert/strict');
const fs = require('node:fs');
const app = require('./app.js');
const seams = require('./surface.js');
const index = JSON.parse(fs.readFileSync('index.json', 'utf8'));
const G = JSON.parse(fs.readFileSync('golden.json', 'utf8'));
const B = q => JSON.parse(G[q]);
const LONG = 'x'.repeat(1025), MR = 'approve mr-0123456789abcdef', REF = 'approve sha256:' + 'ab'.repeat(32);
const UNREAD = 'the verb route answered something this page cannot read', LATE = 'the verb route did not answer in time';
// io: a fetch stub and a controlled clock (setTimeout records, nothing is scheduled).
function rig(fetchImpl) {
  const io = { calls: [], timers: [], cleared: [], AbortController };
  io.fetch = (u, o) => { io.calls.push({ u, o }); return fetchImpl(u, o); };
  io.setTimeout = (fn, ms) => { io.timers.push({ fn, ms }); return io.timers.length; };
  io.clearTimeout = id => { io.cleared.push(id); };
  return io;
}
const reply = (status, text) => () => Promise.resolve({ status, text: () => Promise.resolve(text) });
const golden = q => reply(200, G[q]);
function spies() { const fx = { log: [] }; for (const k of ['address', 'assign', 'panel', 'say', 'clear']) fx[k] = (...a) => { fx.log.push([k, ...a]); }; return fx; }
async function run(raw, fetchImpl) {
  const io = rig(fetchImpl), act = await app.askVerb(raw, io), fx = spies();
  app.verbAct(act, raw, fx);
  assert.equal(io.calls.length, 1, raw); // the page never short-circuits on its own word list (VP4)
  assert.equal(io.calls[0].u, '/surface/verb.json?q=' + encodeURIComponent(raw));
  assert(io.calls[0].o.signal instanceof AbortSignal, raw);
  assert.deepEqual(io.timers.map(x => x.ms), [app.VERB_ROUTE_BUDGET_MS], raw);
  assert.deepEqual(io.cleared, [1], 'the timer is cleared: ' + raw);
  return { act, log: fx.log };
}
(async () => {
  assert.equal(app.VERB_ROUTE_BUDGET_MS, 15000);
  // VP1, VP4, VP6: the address answers reach the address path with the raw line, after clearing the bar.
  for (const q of ['row-189@6', 'world shared@0ffbb9c', 'back', 'world why', 'world session', 'world episode']) {
    const { act, log } = await run(q, golden(q));
    assert.deepEqual(act, { address: true }, q); assert.deepEqual(log, [['clear'], ['address', q]], q);
  }
  // Must-live: that path gives today's answers.
  const at = raw => app.addressTarget(raw, index, seams.parseAddress);
  assert.deepEqual(at('row-189@6'), { open: { brane: 'row-189', at: '652c341d5072c5f49fefe29f9409c93a5119a5be' } });
  assert.deepEqual(at('world shared@0ffbb9c'), { open: { brane: 'shared', at: '0ffbb9c9d386f83aaa45b98e4d4b7763a5c5b845' } });
  assert.deepEqual(at('back'), { home: true });
  // Results: the CLI's text, its JSON, its not-found line with exit 1.
  for (const [q, body] of [['episode report w1', B('episode report w1').text], ['episode report w1 --json', B('episode report w1 --json').json],
    ['episode report nobody', B('episode report nobody').stderr + 'exit 1']]) {
    const { act, log } = await run(q, golden(q));
    assert.deepEqual(act, { panel: { head: '$ ' + B(q).cli, body, copy: '' } }, q); assert.deepEqual(log, [['clear'], ['panel', act.panel]], q);
  }
  assert(B('episode report nobody').stderr.endsWith('has no committed calls\n'));
  // Terminal: the CLI line, with a copy.
  for (const q of ['why', 'episode finish w1 --preset p']) {
    const { act, log } = await run(q, golden(q)), b = B(q);
    assert.deepEqual(act, { panel: { head: '$ ' + b.cli, body: b.reason + '\nAt a terminal: ' + b.cli, copy: b.cli } }, q);
    assert.deepEqual(log, [['clear'], ['panel', act.panel]], q);
  }
  // Navigate: the two card prefixes (VP5's must-live).
  for (const [q, href] of [[MR, '/surface/approve/mr-0123456789abcdef'], [REF, '/surface/decide/sha256:' + 'ab'.repeat(32)]]) {
    const { act, log } = await run(q, golden(q));
    assert.equal(B(q).href, href); assert.deepEqual(act, { assign: href }); assert.deepEqual(log, [['clear'], ['assign', href]]);
  }
  // VP5: any other target is refused.
  for (const href of ['https://evil/x', '/surface/model.json', '//evil/surface/approve/mr-']) {
    const { log } = await run(MR, reply(200, JSON.stringify({ ...B(MR), href })));
    assert.deepEqual(log, [['say', 'refused: unexpected target']], href);
  }
  // Refused (VP6, the shadow set; TooLong): the reason is said; nothing else happens.
  for (const [q, has] of [['log', 'at a terminal: ailang-worldd log'], ['ailang-worldd world get world-1', 'ailang-worldd world get world-1'],
    ['episode', 'needs a subcommand'], ['world a b c', 'unknown address: try "world"'], [LONG, 'TooLong']]) {
    const { act, log } = await run(q, golden(q));
    assert.deepEqual(act, { say: B(q).reason }, q); assert.deepEqual(log, [['say', B(q).reason]], q); assert(B(q).reason.includes(has), q);
  }
  assert(B(LONG).reason.startsWith('TooLong'));
  // VP2: no silent fallback. Every failure says its message and calls nothing else, for verbs and addresses alike.
  for (const q of ['why', 'log', 'row-189@6', 'world shared@0ffbb9c']) {
    const b = B(q), drop = k => { const c = { ...b }; delete c[k]; return JSON.stringify(c); };
    const failures = [404, 400, 421, 500].map(s => [reply(s, G[q]), 'the verb route did not answer: HTTP ' + s]).concat([
      [() => Promise.reject(new TypeError('Failed to fetch')), 'the verb route is unreachable: Failed to fetch'],
      [reply(200, '<html>not json</html>'), UNREAD], [reply(200, drop('schema')), UNREAD], [reply(200, JSON.stringify({ ...b, outcome: 'run' })), UNREAD], [reply(200, JSON.stringify({ ...b, outcome: 'constructor' })), UNREAD],
      [reply(200, JSON.stringify({ ...b, v: 2 })), UNREAD], [reply(200, 'null'), UNREAD], [reply(200, JSON.stringify([b])), UNREAD]]);
    for (const [f, msg] of failures) { const { log } = await run(q, f); assert.deepEqual(log, [['say', msg]], q + ': ' + msg); }
  }
  // An answer missing its outcome's own members is unreadable too.
  for (const body of [{ ...B('why'), cli: '' }, { ...B('why'), reason: undefined }, { ...B(MR), href: 7 }, { ...B('log'), reason: '' },
    { ...B('episode report w1'), json: '{}' }, { ...B('episode report nobody'), exit: 0 }, { ...B('episode report w1'), text: '' }]) {
    const { log } = await run('why', reply(200, JSON.stringify(body)));
    assert.deepEqual(log, [['say', UNREAD]], JSON.stringify(body).slice(0, 80));
  }
  // VP2 hangs, on a controlled clock: stalled headers, then a stalled body. Firing the budget's timer aborts the ask.
  for (const fetchImpl of [() => new Promise(() => {}), () => Promise.resolve({ status: 200, text: () => new Promise(() => {}) })]) {
    const io = rig(fetchImpl), p = app.askVerb('why', io);
    await new Promise(r => setImmediate(r));
    assert.deepEqual(io.timers.map(x => x.ms), [15000]); assert.deepEqual(io.cleared, []); assert(!io.calls[0].o.signal.aborted);
    io.timers[0].fn();
    const act = await p, fx = spies(); app.verbAct(act, 'why', fx);
    assert.deepEqual(fx.log, [['say', LATE]]); assert.deepEqual(io.cleared, [1]); assert(io.calls[0].o.signal.aborted);
  }
  // Superseded (planner S-D1): a second ask aborts the first, whose answer is dropped (no say, no navigation),
  // whether it was still waiting (A) or had already arrived (B); the second answers for the bar.
  let release; const io1 = rig(() => new Promise(r => { release = r; }));
  const p1 = app.askVerb('row-189@6', io1), p2 = app.askVerb('log', rig(golden('log')));
  assert(io1.calls[0].o.signal.aborted);
  release({ status: 200, text: () => Promise.resolve(G['row-189@6']) });
  const a1 = await p1, a2 = await p2;
  assert.deepEqual(a1, { superseded: true }); assert.deepEqual(io1.cleared, [1]); assert.deepEqual(a2, { say: B('log').reason });
  { const fx = spies(); app.verbAct(a1, 'row-189@6', fx); assert.deepEqual(fx.log, []); }
  const io3 = rig(golden('row-189@6')), p3 = app.askVerb('row-189@6', io3), p4 = app.askVerb('why', rig(golden('why')));
  assert.deepEqual(await p3, { superseded: true }); assert.deepEqual(io3.cleared, [1]); assert(Boolean((await p4).panel));
  // VP3: the panel is built with textContent only; a hostile report or reason stays text.
  const fakeDoc = () => { const doc = { made: [], createElement(tag) {
    const el = { tag, children: [], dataset: {}, style: {}, attrs: {}, textContent: '', setAttribute(k, v) { this.attrs[k] = v; }, append(...c) { this.children.push(...c); } };
    Object.defineProperty(el, 'innerHTML', { get() { throw new Error('innerHTML read'); }, set() { throw new Error('innerHTML set'); } });
    doc.made.push(el); return el; } }; return doc; };
  { const doc = fakeDoc(), el = app.verbPanel(doc, app.verbOutcome(200, G['episode report w1']).panel);
    assert.deepEqual(doc.made.map(e => e.tag), ['div', 'div', 'pre']); assert.equal(el.id, 'verb'); assert.equal(el.children[0].textContent, '$ ailang-worldd episode report w1');
    assert(el.children[1].textContent.includes('#2 <img src=x onerror=alert(1)>: <b>hostile</b> "reason"')); assert.equal(el.children[1].textContent, B('episode report w1').text); }
  { const doc = fakeDoc(), el = app.verbPanel(doc, app.verbOutcome(200, G.why).panel);
    assert.deepEqual(doc.made.map(e => e.tag), ['div', 'div', 'pre', 'button']); assert.equal(el.children[2].dataset.copy, 'ailang-worldd why');
    assert.equal(el.children[2].textContent, 'Copy the line'); }
  { const doc = fakeDoc(), el = app.verbPanel(doc, { head: '$ <b>x</b>', body: '<img src=x onerror=1>', copy: '' });
    assert.equal(el.children[1].textContent, '<img src=x onerror=1>'); assert.equal(doc.made.length, 3); }
  console.log('APP VERBS PASS');
})().catch(e => { console.error(e); process.exit(1); });
`
